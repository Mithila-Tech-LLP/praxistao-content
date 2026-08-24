# Chapter 56: Building Resilient Consumers in Go

A consumer that works in development fails in production in a hundred ways: the broker goes down, the database is slow, the message is malformed, the service it calls is unavailable. Resilient consumers handle all of this gracefully.

## Table of Contents

1. The Resilient Consumer Framework
2. Graceful Shutdown
3. Circuit Breakers
4. Backpressure
5. Observability: Metrics and Tracing
6. The Complete Resilient Consumer
7. Exercises

---

## 1. The Resilient Consumer Framework

A production consumer needs:

```
Read message
    ↓
Deserialize (handle bad messages)
    ↓
Validate (handle schema violations)
    ↓
Process (with timeout, circuit breaker)
    ↓
Commit (at-least-once guarantee)
    ↓
On failure: retry? DLQ? alert?
```

```go
package consumer

import (
    "context"
    "fmt"
    "log"
    "time"

    kafka "github.com/segmentio/kafka-go"
)

type Config struct {
    Brokers        []string
    Topic          string
    GroupID        string
    MaxRetries     int
    ProcessTimeout time.Duration
    DLQTopic       string
}

type Handler[T any] func(ctx context.Context, msg T) error

type Consumer[T any] struct {
    cfg     Config
    reader  *kafka.Reader
    dlq     *kafka.Writer
    handler Handler[T]
    metrics *Metrics
}

func New[T any](cfg Config, handler Handler[T]) *Consumer[T] {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:        cfg.Brokers,
        Topic:          cfg.Topic,
        GroupID:        cfg.GroupID,
        MinBytes:       1,
        MaxBytes:       10e6,
        MaxWait:        500 * time.Millisecond,
        CommitInterval: 0, // manual commit
    })

    var dlq *kafka.Writer
    if cfg.DLQTopic != "" {
        dlq = &kafka.Writer{Addr: kafka.TCP(cfg.Brokers...), Topic: cfg.DLQTopic}
    }

    return &Consumer[T]{
        cfg:     cfg,
        reader:  reader,
        dlq:     dlq,
        handler: handler,
        metrics: newMetrics(cfg.GroupID),
    }
}
```

---

## 2. Graceful Shutdown

```go
func (c *Consumer[T]) Run(ctx context.Context) error {
    log.Printf("Consumer %s starting on %s", c.cfg.GroupID, c.cfg.Topic)
    defer log.Printf("Consumer %s stopped", c.cfg.GroupID)
    defer c.reader.Close()

    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return nil // normal shutdown
            }
            c.metrics.ReadErrors.Inc()
            log.Printf("fetch error: %v — backing off", err)
            select {
            case <-time.After(time.Second):
            case <-ctx.Done():
                return nil
            }
            continue
        }

        c.metrics.MessagesRead.Inc()
        start := time.Now()

        if err := c.processWithRetry(ctx, msg); err != nil {
            c.metrics.ProcessErrors.Inc()
            c.sendToDLQ(ctx, msg, err)
        } else {
            c.metrics.MessagesProcessed.Inc()
        }

        c.metrics.ProcessDuration.Observe(time.Since(start).Seconds())

        if err := c.reader.CommitMessages(ctx, msg); err != nil {
            if ctx.Err() != nil {
                return nil
            }
            log.Printf("commit error: %v", err)
        }
    }
}

func (c *Consumer[T]) processWithRetry(ctx context.Context, raw kafka.Message) error {
    var msg T
    if err := json.Unmarshal(raw.Value, &msg); err != nil {
        // Malformed message: don't retry, send to DLQ immediately
        return fmt.Errorf("deserialize: %w", err)
    }

    var lastErr error
    for attempt := 0; attempt <= c.cfg.MaxRetries; attempt++ {
        if attempt > 0 {
            // Exponential backoff: 1s, 2s, 4s, 8s...
            delay := time.Duration(1<<attempt) * time.Second
            if delay > 30*time.Second {
                delay = 30 * time.Second
            }
            select {
            case <-time.After(delay):
            case <-ctx.Done():
                return ctx.Err()
            }
            log.Printf("retry %d/%d for message %s", attempt, c.cfg.MaxRetries, raw.Key)
        }

        processCtx, cancel := context.WithTimeout(ctx, c.cfg.ProcessTimeout)
        lastErr = c.handler(processCtx, msg)
        cancel()

        if lastErr == nil {
            return nil
        }

        // Don't retry if it's a non-retryable error
        var nonRetryable *NonRetryableError
        if errors.As(lastErr, &nonRetryable) {
            return lastErr
        }
    }
    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

// NonRetryableError signals that a message should go directly to DLQ
type NonRetryableError struct {
    Err error
}
func (e *NonRetryableError) Error() string { return e.Err.Error() }
func (e *NonRetryableError) Unwrap() error { return e.Err }

func (c *Consumer[T]) sendToDLQ(ctx context.Context, msg kafka.Message, err error) {
    if c.dlq == nil {
        log.Printf("DLQ not configured, dropping message: %v", err)
        return
    }

    dlqMsg := map[string]interface{}{
        "original_topic": msg.Topic,
        "original_key":   string(msg.Key),
        "payload":        string(msg.Value),
        "error":          err.Error(),
        "failed_at":      time.Now(),
    }
    data, _ := json.Marshal(dlqMsg)

    if err := c.dlq.WriteMessages(ctx, kafka.Message{
        Key:   msg.Key,
        Value: data,
    }); err != nil {
        log.Printf("DLQ write failed: %v — message dropped", err)
    } else {
        log.Printf("Message sent to DLQ: %s", msg.Key)
    }
}
```

---

## 3. Circuit Breakers

If the downstream service (database, API) is down, stop hammering it. A circuit breaker pauses calls temporarily.

```go
type CircuitBreaker struct {
    mu           sync.Mutex
    failures     int
    lastFailure  time.Time
    state        string // "closed", "open", "half-open"
    maxFailures  int
    resetTimeout time.Duration
}

func NewCircuitBreaker(maxFailures int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:        "closed",
        maxFailures:  maxFailures,
        resetTimeout: resetTimeout,
    }
}

func (cb *CircuitBreaker) Call(fn func() error) error {
    cb.mu.Lock()
    state := cb.state
    if state == "open" {
        if time.Since(cb.lastFailure) > cb.resetTimeout {
            cb.state = "half-open"
            state = "half-open"
        }
    }
    cb.mu.Unlock()

    if state == "open" {
        return fmt.Errorf("circuit breaker open: service unavailable")
    }

    err := fn()

    cb.mu.Lock()
    defer cb.mu.Unlock()

    if err != nil {
        cb.failures++
        cb.lastFailure = time.Now()
        if cb.failures >= cb.maxFailures {
            cb.state = "open"
            log.Printf("Circuit breaker opened after %d failures", cb.failures)
        }
    } else {
        // Success
        cb.failures = 0
        cb.state = "closed"
    }

    return err
}

// Usage in message handler
var dbBreaker = NewCircuitBreaker(5, 30*time.Second)

func orderHandler(ctx context.Context, order Order) error {
    return dbBreaker.Call(func() error {
        return saveOrderToDB(ctx, order)
    })
}
```

---

## 4. Backpressure

If the consumer is processing slowly, stop reading more messages until caught up:

```go
type BackpressuredConsumer struct {
    reader    *kafka.Reader
    semaphore chan struct{} // controls concurrency
}

func NewBackpressuredConsumer(reader *kafka.Reader, maxConcurrent int) *BackpressuredConsumer {
    return &BackpressuredConsumer{
        reader:    reader,
        semaphore: make(chan struct{}, maxConcurrent),
    }
}

func (c *BackpressuredConsumer) Run(ctx context.Context, handler func([]byte) error) {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            continue
        }

        // Wait for a slot (backpressure: max maxConcurrent in-flight)
        c.semaphore <- struct{}{}

        go func(m kafka.Message) {
            defer func() { <-c.semaphore }()

            if err := handler(m.Value); err != nil {
                log.Printf("handler error: %v", err)
                return
            }
            c.reader.CommitMessages(ctx, m)
        }(msg)
    }
}
```

---

## 5. Observability: Metrics and Tracing

```go
// Simple metrics (in production: use prometheus/client_golang)
type Metrics struct {
    MessagesRead      *Counter
    MessagesProcessed *Counter
    ProcessErrors     *Counter
    ReadErrors        *Counter
    ProcessDuration   *Histogram
}

type Counter struct {
    mu    sync.Mutex
    value int64
}

func (c *Counter) Inc() {
    c.mu.Lock()
    c.value++
    c.mu.Unlock()
}

func (c *Counter) Get() int64 {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}

type Histogram struct {
    mu     sync.Mutex
    values []float64
}

func (h *Histogram) Observe(v float64) {
    h.mu.Lock()
    h.values = append(h.values, v)
    if len(h.values) > 10000 {
        h.values = h.values[1:]
    }
    h.mu.Unlock()
}

func newMetrics(groupID string) *Metrics {
    return &Metrics{
        MessagesRead:      &Counter{},
        MessagesProcessed: &Counter{},
        ProcessErrors:     &Counter{},
        ReadErrors:        &Counter{},
        ProcessDuration:   &Histogram{},
    }
}

// Expose metrics endpoint
func (c *Consumer[T]) MetricsHandler() http.HandlerFunc {
    return func(w http.ResponseWriter, r *http.Request) {
        json.NewEncoder(w).Encode(map[string]interface{}{
            "messages_read":      c.metrics.MessagesRead.Get(),
            "messages_processed": c.metrics.MessagesProcessed.Get(),
            "process_errors":     c.metrics.ProcessErrors.Get(),
            "read_errors":        c.metrics.ReadErrors.Get(),
        })
    }
}
```

---

## 6. The Complete Resilient Consumer

```go
package main

import (
    "context"
    "fmt"
    "log"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type Order struct {
    ID     string  `json:"id"`
    UserID string  `json:"user_id"`
    Amount float64 `json:"amount"`
}

func main() {
    cfg := consumer.Config{
        Brokers:        []string{"localhost:9092"},
        Topic:          "orders",
        GroupID:        "order-processor",
        MaxRetries:     3,
        ProcessTimeout: 10 * time.Second,
        DLQTopic:       "orders.dead-letter",
    }

    c := consumer.New[Order](cfg, handleOrder)

    // Expose metrics
    http.HandleFunc("/metrics", c.MetricsHandler())
    go http.ListenAndServe(":9090", nil)

    // Graceful shutdown
    ctx, stop := signal.NotifyContext(context.Background(),
        os.Interrupt, syscall.SIGTERM)
    defer stop()

    log.Println("Starting order processor...")
    if err := c.Run(ctx); err != nil {
        log.Fatal("consumer error:", err)
    }
    log.Println("Gracefully shut down")
}

func handleOrder(ctx context.Context, order Order) error {
    // Simulate processing
    if order.Amount < 0 {
        return &consumer.NonRetryableError{
            Err: fmt.Errorf("invalid negative amount: %f", order.Amount),
        }
    }

    // Simulate slow database
    time.Sleep(50 * time.Millisecond)

    log.Printf("Processed order %s for user %s: $%.2f",
        order.ID, order.UserID, order.Amount)
    return nil
}
```

---

## Summary

- `signal.NotifyContext` + `ctx.Done()` for graceful shutdown — finish in-flight messages before stopping.
- Exponential backoff for retries: 1s, 2s, 4s... cap at 30s.
- `NonRetryableError` distinguishes "try again later" from "this message is permanently bad."
- Circuit breakers prevent cascading failures when downstream services are down.
- Semaphore-based backpressure prevents unbounded in-flight message accumulation.
- Expose a `/metrics` HTTP endpoint for monitoring.

### Exercises

**Easy:** Add a `NonRetryableError` check to your existing order consumer. Any message with `amount < 0` should go directly to the DLQ without retrying.

**Medium:** Implement a circuit breaker for a database call in a consumer. Test it by intentionally failing the database 5 times and verifying the circuit opens. Then restore the DB and verify the circuit closes after the reset timeout.

**Hard:** Implement a consumer that processes messages concurrently (up to 10 goroutines) with backpressure. Each handler sleeps a random duration (10-500ms). Verify that at most 10 messages are in-flight at any time, and that the consumer correctly commits offsets even for out-of-order completions.
