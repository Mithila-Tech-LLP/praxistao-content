# Chapter 31: Reliability Patterns — Retry, Idempotency, Circuit Breakers & Sagas

Building reliable distributed services requires patterns that handle partial failures gracefully. These patterns are what separate systems that survive production from those that don't.

## Table of Contents

1. [Retry with Exponential Backoff & Jitter](#1-retry-with-exponential-backoff--jitter)
2. [Idempotency — Safe Retries](#2-idempotency--safe-retries)
3. [Circuit Breaker Pattern](#3-circuit-breaker-pattern)
4. [Bulkhead Pattern](#4-bulkhead-pattern)
5. [Sagas — Distributed Transactions Without 2PC](#5-sagas--distributed-transactions-without-2pc)
6. [Outbox Pattern — At-Least-Once Delivery](#6-outbox-pattern--at-least-once-delivery)
7. [Timeout Strategy](#7-timeout-strategy)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Retry with Exponential Backoff & Jitter

Retrying failed requests is the first defense against transient failures. But naive retry (retry immediately) can cause a "thundering herd" — all clients hammering a recovering service at the same time.

```go
// Naive retry — DANGEROUS:
for i := 0; i < 3; i++ {
    err := callService()
    if err == nil { break }
    // all clients retry at exactly the same time!
}

// Exponential backoff: wait doubles each attempt
// Jitter: add randomness to spread retries across time
func withRetry(ctx context.Context, maxAttempts int, fn func() error) error {
    var lastErr error
    
    for attempt := 0; attempt < maxAttempts; attempt++ {
        if err := fn(); err == nil {
            return nil // success
        } else {
            lastErr = err
            
            // Don't retry non-retryable errors (404, 403, etc.)
            if !isRetryable(err) {
                return err
            }
        }
        
        // Exponential backoff: 100ms, 200ms, 400ms, 800ms...
        baseDelay := 100 * time.Millisecond * (1 << attempt)
        // Add jitter: randomize between 50-100% of base delay
        jitter := time.Duration(rand.Int63n(int64(baseDelay / 2)))
        delay := baseDelay/2 + jitter
        
        // Cap at 30 seconds
        if delay > 30*time.Second {
            delay = 30 * time.Second
        }
        
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(delay):
        }
    }
    
    return fmt.Errorf("max retries exceeded: %w", lastErr)
}

func isRetryable(err error) bool {
    // Retry on network errors, 5xx errors
    // Don't retry on 4xx (client error — won't get better)
    var apiErr *APIError
    if errors.As(err, &apiErr) {
        return apiErr.StatusCode >= 500
    }
    return true // network errors are retryable
}
```

---

## 2. Idempotency — Safe Retries

Before retrying, your operation must be idempotent — calling it multiple times must have the same effect as calling it once.

```
Non-idempotent: POST /payments (each call charges the card)
Idempotent:     POST /payments with idempotency-key header
                Same key = same outcome (charge card once)

POST is not idempotent by default → make it idempotent with an idempotency key
GET, PUT, DELETE are naturally idempotent
```

### Implementing Idempotency Keys

```go
// Client sends a unique idempotency key with every request
// Server uses it to detect duplicate requests

// Database schema:
// idempotency_keys(key, response_code, response_body, created_at)

func (h *PaymentHandler) ChargeCard(w http.ResponseWriter, r *http.Request) {
    idempotencyKey := r.Header.Get("Idempotency-Key")
    if idempotencyKey == "" {
        http.Error(w, "Idempotency-Key header required", http.StatusBadRequest)
        return
    }
    
    // Check if we've seen this key before
    existing, err := h.store.GetIdempotencyRecord(r.Context(), idempotencyKey)
    if err == nil && existing != nil {
        // Replay the stored response
        w.WriteHeader(existing.ResponseCode)
        w.Write([]byte(existing.ResponseBody))
        return
    }
    
    // Process the payment
    result, err := h.paymentSvc.Charge(r.Context(), parseBody(r))
    statusCode := 200
    if err != nil { statusCode = 500 }
    
    responseBody, _ := json.Marshal(result)
    
    // Store the result for future duplicate detection
    h.store.SaveIdempotencyRecord(r.Context(), &IdempotencyRecord{
        Key:          idempotencyKey,
        ResponseCode: statusCode,
        ResponseBody: string(responseBody),
    })
    
    w.WriteHeader(statusCode)
    w.Write(responseBody)
}
```

---

## 3. Circuit Breaker Pattern

The circuit breaker prevents repeated calls to a failing service, giving it time to recover and protecting your service from cascading failures.

```
States:
  CLOSED (normal):   requests flow through. Count failures.
  OPEN (tripped):    requests fail immediately (no network call). Reset timer running.
  HALF-OPEN (probe): allow a few requests through. If successful, go CLOSED. Else, OPEN again.

Thresholds:
  Trip to OPEN when: 50% of last 100 requests failed
  Try HALF-OPEN after: 30 seconds in OPEN state
  Return to CLOSED after: 5 consecutive successes in HALF-OPEN
```

```go
type CircuitBreaker struct {
    mu           sync.Mutex
    state        string        // "closed", "open", "half-open"
    failureCount int
    successCount int
    threshold    int
    resetTimeout time.Duration
    lastFailTime time.Time
}

func NewCircuitBreaker(threshold int, resetTimeout time.Duration) *CircuitBreaker {
    return &CircuitBreaker{
        state:        "closed",
        threshold:    threshold,
        resetTimeout: resetTimeout,
    }
}

func (cb *CircuitBreaker) Execute(fn func() error) error {
    cb.mu.Lock()
    
    switch cb.state {
    case "open":
        if time.Since(cb.lastFailTime) > cb.resetTimeout {
            cb.state = "half-open"
            cb.successCount = 0
        } else {
            cb.mu.Unlock()
            return errors.New("circuit breaker open: service unavailable")
        }
    case "closed", "half-open":
        // allow through
    }
    cb.mu.Unlock()
    
    err := fn()
    
    cb.mu.Lock()
    defer cb.mu.Unlock()
    
    if err != nil {
        cb.failureCount++
        cb.successCount = 0
        cb.lastFailTime = time.Now()
        if cb.failureCount >= cb.threshold {
            cb.state = "open"
        }
        return err
    }
    
    // Success
    cb.failureCount = 0
    if cb.state == "half-open" {
        cb.successCount++
        if cb.successCount >= 5 {
            cb.state = "closed"
        }
    }
    return nil
}

// Usage:
cb := NewCircuitBreaker(5, 30*time.Second)

func callPaymentService() error {
    return cb.Execute(func() error {
        return paymentSvc.Charge(ctx, amount)
    })
}
```

---

## 4. Bulkhead Pattern

Isolate different parts of your system so failures in one don't consume resources from another.

```go
// Without bulkhead: one goroutine pool for all operations
// If the payment service is slow, all goroutines pile up waiting for it,
// and no goroutines are left to serve other requests

// With bulkhead: separate pools for different operations
type BulkheadPool struct {
    sem chan struct{}
}

func NewBulkheadPool(size int) *BulkheadPool {
    return &BulkheadPool{sem: make(chan struct{}, size)}
}

func (b *BulkheadPool) Execute(ctx context.Context, fn func() error) error {
    select {
    case b.sem <- struct{}{}: // acquire a slot
        defer func() { <-b.sem }() // release slot when done
        return fn()
    case <-ctx.Done():
        return ctx.Err()
    default:
        return errors.New("bulkhead full: too many concurrent requests")
    }
}

// Separate pools for different dependencies:
var (
    paymentBulkhead = NewBulkheadPool(20)   // max 20 concurrent payment calls
    emailBulkhead   = NewBulkheadPool(10)   // max 10 concurrent email calls
    dbBulkhead      = NewBulkheadPool(50)   // max 50 concurrent db calls
)
```

---

## 5. Sagas — Distributed Transactions Without 2PC

When a business operation spans multiple services, you need to coordinate across them without a distributed transaction (which is fragile). Sagas coordinate local transactions with compensating actions for rollback.

### Choreography-Based Saga

Each service listens for events and decides what to do next:

```
Order Service    Payment Service    Inventory Service    Notification Service

1. OrderPlaced event ──▶
                    2. PaymentProcessed event ──▶
                                           3. InventoryReserved event ──▶
                                                                    4. OrderConfirmedEmail

Failure: if Inventory fails:
InventoryFailed event ──▶ Payment Service compensates (RefundPayment)
PaymentRefunded event ──▶ Order Service compensates (CancelOrder)
```

### Orchestration-Based Saga

A central orchestrator drives the workflow:

```go
// Saga orchestrator — controls the workflow
type OrderSaga struct {
    orderSvc     OrderService
    paymentSvc   PaymentService
    inventorySvc InventoryService
    emailSvc     EmailService
}

type SagaState struct {
    OrderID   string
    PaymentID string
    Step      string
    Failed    bool
}

func (s *OrderSaga) Execute(ctx context.Context, req CreateOrderRequest) error {
    state := &SagaState{Step: "start"}
    
    // Step 1: Create order
    orderID, err := s.orderSvc.CreateOrder(ctx, req)
    if err != nil { return err }
    state.OrderID = orderID
    state.Step = "order_created"
    
    // Step 2: Process payment — if this fails, cancel the order
    paymentID, err := s.paymentSvc.Charge(ctx, req.Amount, req.CardToken)
    if err != nil {
        // Compensate: cancel the order
        s.orderSvc.CancelOrder(ctx, orderID)
        return fmt.Errorf("payment failed: %w", err)
    }
    state.PaymentID = paymentID
    state.Step = "payment_complete"
    
    // Step 3: Reserve inventory — if this fails, refund payment and cancel order
    err = s.inventorySvc.Reserve(ctx, req.ProductID, req.Quantity)
    if err != nil {
        s.paymentSvc.Refund(ctx, paymentID)     // compensate payment
        s.orderSvc.CancelOrder(ctx, orderID)    // compensate order
        return fmt.Errorf("inventory reservation failed: %w", err)
    }
    state.Step = "inventory_reserved"
    
    // Step 4: Send confirmation email (best-effort, don't fail the saga)
    s.emailSvc.SendConfirmation(ctx, req.Email, orderID)
    
    return nil
}
```

---

## 6. Outbox Pattern — At-Least-Once Delivery

Reliably publish events to a message broker in the same transaction as a database write:

```
Problem:
  tx.Begin()
  UPDATE orders SET status = 'confirmed'
  tx.Commit()
  kafka.Publish("order.confirmed", ...) ← crash here! Order confirmed but event not published.

Solution: Transactional Outbox
  tx.Begin()
  UPDATE orders SET status = 'confirmed'
  INSERT INTO outbox (event_type, payload) VALUES ('order.confirmed', '...')
  tx.Commit()
  
  Separate background worker reads from outbox table and publishes to Kafka
  Marks outbox records as processed
  If worker crashes: restart, re-publish. Consumers must be idempotent!
```

```go
// Write to outbox atomically with the domain change
func (r *OrderRepo) ConfirmOrder(ctx context.Context, tx *sql.Tx, orderID string) error {
    _, err := tx.ExecContext(ctx, "UPDATE orders SET status = 'confirmed' WHERE id = $1", orderID)
    if err != nil { return err }
    
    payload, _ := json.Marshal(map[string]string{"order_id": orderID})
    _, err = tx.ExecContext(ctx, 
        "INSERT INTO outbox (event_type, payload, created_at) VALUES ($1, $2, NOW())",
        "order.confirmed", payload)
    return err
}

// Background worker that reads outbox and publishes to Kafka
func (w *OutboxWorker) Run(ctx context.Context) {
    ticker := time.NewTicker(100 * time.Millisecond)
    for {
        select {
        case <-ticker.C:
            w.processOutbox(ctx)
        case <-ctx.Done():
            return
        }
    }
}

func (w *OutboxWorker) processOutbox(ctx context.Context) {
    rows, _ := w.db.QueryContext(ctx,
        "SELECT id, event_type, payload FROM outbox WHERE published = false ORDER BY created_at LIMIT 100")
    
    for rows.Next() {
        var id int64
        var eventType, payload string
        rows.Scan(&id, &eventType, &payload)
        
        if err := w.kafka.Publish(ctx, eventType, payload); err != nil {
            continue // retry next tick
        }
        
        w.db.ExecContext(ctx, "UPDATE outbox SET published = true WHERE id = $1", id)
    }
}
```

---

## 7. Timeout Strategy

```go
// Every external call needs a timeout — never call without one
func callWithTimeout(ctx context.Context, timeout time.Duration, fn func(context.Context) error) error {
    ctx, cancel := context.WithTimeout(ctx, timeout)
    defer cancel()
    return fn(ctx)
}

// Timeout hierarchy: outer timeouts must be larger than inner
// HTTP handler context (30s) > gRPC call timeout (5s) > DB query timeout (2s)
// Otherwise the inner call succeeds but the outer context is already cancelled

// Never use sleep to implement timeout — always use context
// Bad:
go func() {
    time.Sleep(5 * time.Second)
    cancel() // too late — goroutine may already be stuck
}()

// Good: context deadline propagates to all downstream calls
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()
resp, err := httpClient.Do(req.WithContext(ctx))
```

---

## 8. Interview Questions & Model Answers

**Q: How would you implement idempotent payment processing?**

"The key is the idempotency key: clients generate a UUID per payment attempt and send it in a header. The server stores this key along with the response in a database (with a unique constraint on the key). Before processing, the server checks if it's seen the key before — if yes, return the stored response without reprocessing. If no, process the payment, store the result atomically, then return. This ensures clients can safely retry on network failures without double-charging. The idempotency table should have a TTL (e.g., 24 hours) to avoid unlimited storage growth."

**Q: What is the difference between saga choreography and orchestration?**

"In choreography, services communicate via events and each service decides its next action based on events it receives. There's no central coordinator — it's decentralized. This is simpler to start with but harder to debug and trace as workflows grow complex. In orchestration, a central saga orchestrator explicitly calls each service in sequence and handles compensation on failure. It's easier to understand the entire workflow (it's in one place) and easier to add retry/compensation logic. The trade-off is a central component that needs to be fault-tolerant itself. Most production systems at scale prefer orchestration for complex workflows because it's easier to reason about."

---

## Summary

- **Retry:** use exponential backoff with jitter. Only retry on retryable errors (5xx, network). Don't retry 4xx.
- **Idempotency:** client-generated keys. Server stores result for duplicate detection. Required for any non-idempotent operation that can be retried.
- **Circuit breaker:** CLOSED → OPEN (on threshold) → HALF-OPEN (after timeout) → CLOSED (on success). Prevents cascading failures.
- **Bulkhead:** separate thread/goroutine pools for different dependencies. One slow service can't starve all resources.
- **Saga:** coordinate multi-service operations. Each step has a compensating action for rollback.
- **Outbox:** write domain change + event in same DB transaction. Background worker publishes events. Consumers must be idempotent.
- **Timeouts:** every external call needs one. Context cancellation propagates downstream. Timeout hierarchy must be respected.
