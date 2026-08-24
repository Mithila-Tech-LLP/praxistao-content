# Chapter 52: Kafka — The Distributed Log

Kafka processes over a trillion messages per day at LinkedIn. It's the backbone of Uber's real-time pricing, Netflix's event pipeline, and thousands of other high-scale systems. This chapter covers Kafka from first principles to production patterns.

## Table of Contents

1. Kafka's Design Philosophy
2. Core Concepts: Topics, Partitions, Offsets
3. Kafka Architecture
4. Docker Setup
5. Building with Kafka in Go
6. Producer Patterns
7. Consumer Patterns
8. Mini Project: Order Event Pipeline
9. Exercises

---

## 1. Kafka's Design Philosophy

Kafka was built at LinkedIn in 2011 to solve one problem: move massive amounts of data reliably between services.

**Key design decisions:**

**1. Append-only log:** Writes are always appends. Reads are sequential. This is 100-1000x faster than random read/write patterns.

**2. Disk as primary storage:** Unlike Redis (RAM-first), Kafka is designed for disk. 100 GB on disk is cheap. 100 GB in RAM is expensive. Combined with the OS page cache, disk performance approaches RAM for sequential reads.

**3. Zero-copy transfer:** Kafka uses Linux's `sendfile()` syscall to copy data from disk directly to the network socket, bypassing user space. This dramatically reduces CPU usage.

**4. Pull-based consumption:** Consumers pull messages at their own rate. This means a slow consumer won't crash the broker — it just falls behind.

**5. Batching everywhere:** Producers batch before sending. Brokers batch before writing to disk. Consumers read in batches.

---

## 2. Core Concepts: Topics, Partitions, Offsets

```
Topic: "user-events"
  Partition 0: msg(offset=0) → msg(offset=1) → msg(offset=2) → ...
  Partition 1: msg(offset=0) → msg(offset=1) → msg(offset=2) → ...
  Partition 2: msg(offset=0) → msg(offset=1) → msg(offset=2) → ...
```

**Message structure:**
- **Key:** Optional bytes. Used for partitioning (same key → same partition).
- **Value:** The actual message payload (usually JSON, Avro, or Protobuf).
- **Headers:** Key-value metadata (tracing IDs, content type, etc.).
- **Timestamp:** When the message was produced.
- **Offset:** Position in the partition (assigned by broker).

**Partition assignment:**
- No key: round-robin across partitions.
- With key: `murmur2(key) % partitions`. Same key always → same partition.
- Custom: implement a `Partitioner` interface.

**Why keys matter for ordering:**
```
Without key: order events for the same user might land in different partitions → different consumers → out-of-order processing
With key = userID: all events for user 123 go to the same partition → same consumer → guaranteed order per user
```

---

## 3. Kafka Architecture

```
Cluster (3 brokers):
  ┌─────────────────────────────────────────────────────────────┐
  │                                                             │
  │  Broker 1 (leader for p0, p2)    Partition 0  [leader]     │
  │  Broker 2 (leader for p1)        Partition 1  [leader]     │
  │  Broker 3 (follower for all)     Partition 2  [leader]     │
  │                                                             │
  │  ZooKeeper / KRaft: cluster metadata, leader election      │
  └─────────────────────────────────────────────────────────────┘
```

**KRaft (Kafka without ZooKeeper):** Since Kafka 3.3, ZooKeeper is optional. KRaft mode uses an internal Raft consensus protocol, making Kafka self-contained. Production deployments are moving to KRaft.

**Controller:** One broker acts as the cluster controller, managing partition leadership and broker registration. In KRaft, this is a dedicated set of controller nodes.

---

## 4. Docker Setup

```yaml
# docker-compose.yml
services:
  kafka:
    image: apache/kafka:3.7.0
    ports:
      - "9092:9092"
    environment:
      KAFKA_NODE_ID: 1
      KAFKA_PROCESS_ROLES: broker,controller
      KAFKA_LISTENERS: PLAINTEXT://0.0.0.0:9092,CONTROLLER://0.0.0.0:9093
      KAFKA_ADVERTISED_LISTENERS: PLAINTEXT://localhost:9092
      KAFKA_CONTROLLER_LISTENER_NAMES: CONTROLLER
      KAFKA_LISTENER_SECURITY_PROTOCOL_MAP: CONTROLLER:PLAINTEXT,PLAINTEXT:PLAINTEXT
      KAFKA_CONTROLLER_QUORUM_VOTERS: 1@localhost:9093
      KAFKA_OFFSETS_TOPIC_REPLICATION_FACTOR: 1
      KAFKA_LOG_DIRS: /var/lib/kafka/data
      CLUSTER_ID: MkU3OEVBNTcwNTJENDM2Qg==
```

```bash
docker compose up -d

# Create a topic
docker exec kafka /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --create --topic user-events \
  --partitions 3 \
  --replication-factor 1

# List topics
docker exec kafka /opt/kafka/bin/kafka-topics.sh \
  --bootstrap-server localhost:9092 \
  --list

# Tail messages
docker exec kafka /opt/kafka/bin/kafka-console-consumer.sh \
  --bootstrap-server localhost:9092 \
  --topic user-events \
  --from-beginning
```

---

## 5. Building with Kafka in Go

```bash
go get github.com/segmentio/kafka-go
```

```go
package main

import (
    "context"
    "fmt"
    "log"
    "time"

    kafka "github.com/segmentio/kafka-go"
)

func main() {
    ctx := context.Background()

    // Auto-create topic if needed
    conn, _ := kafka.DialLeader(ctx, "tcp", "localhost:9092", "user-events", 0)
    conn.CreateTopics(kafka.TopicConfig{
        Topic:             "user-events",
        NumPartitions:     3,
        ReplicationFactor: 1,
    })
    conn.Close()

    // Start producer and consumer
    go producer(ctx)
    consumer(ctx)
}
```

---

## 6. Producer Patterns

```go
func producer(ctx context.Context) {
    w := &kafka.Writer{
        Addr:     kafka.TCP("localhost:9092"),
        Topic:    "user-events",
        Balancer: &kafka.Hash{}, // hash(key) to pick partition

        // Batching settings
        BatchSize:    100,
        BatchTimeout: 10 * time.Millisecond,

        // Reliability settings
        RequiredAcks: kafka.RequireAll, // wait for all replicas
        MaxAttempts:  3,
    }
    defer w.Close()

    events := []struct {
        UserID string
        Action string
    }{
        {"user-1", "login"},
        {"user-2", "purchase"},
        {"user-1", "logout"},
        {"user-3", "signup"},
    }

    var msgs []kafka.Message
    for _, e := range events {
        payload := fmt.Sprintf(`{"user_id":"%s","action":"%s","ts":"%s"}`,
            e.UserID, e.Action, time.Now().Format(time.RFC3339))
        msgs = append(msgs, kafka.Message{
            Key:   []byte(e.UserID),  // same user → same partition → ordered
            Value: []byte(payload),
        })
    }

    if err := w.WriteMessages(ctx, msgs...); err != nil {
        log.Fatal("write:", err)
    }
    fmt.Printf("Produced %d messages\n", len(msgs))
}
```

**Async producer with error handling:**

```go
type AsyncProducer struct {
    writer *kafka.Writer
    errCh  chan error
}

func NewAsyncProducer(brokers []string, topic string) *AsyncProducer {
    w := &kafka.Writer{
        Addr:                   kafka.TCP(brokers...),
        Topic:                  topic,
        Async:                  true,  // don't wait for ACK
        Completion: func(msgs []kafka.Message, err error) {
            if err != nil {
                // Log and potentially retry
                log.Printf("async write error: %v (%d messages)", err, len(msgs))
            }
        },
    }
    return &AsyncProducer{writer: w}
}

func (p *AsyncProducer) Send(key, value []byte) {
    p.writer.WriteMessages(context.Background(), kafka.Message{
        Key: key, Value: value,
    })
    // Returns immediately — completion callback fires later
}
```

---

## 7. Consumer Patterns

```go
func consumer(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers:        []string{"localhost:9092"},
        Topic:          "user-events",
        GroupID:        "analytics-service",   // consumer group
        MinBytes:       1,                     // fetch when at least 1 byte available
        MaxBytes:       10e6,                  // 10 MB max per fetch
        MaxWait:        500 * time.Millisecond, // wait up to 500ms for batches
        StartOffset:    kafka.FirstOffset,      // start from beginning if new group
        CommitInterval: time.Second,            // auto-commit every second
    })
    defer r.Close()

    fmt.Println("Consumer started, waiting for messages...")

    for {
        msg, err := r.ReadMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return // context cancelled
            }
            log.Printf("read error: %v", err)
            continue
        }

        fmt.Printf("Partition %d | Offset %d | Key: %s | Value: %s\n",
            msg.Partition, msg.Offset, msg.Key, msg.Value)
    }
}
```

**Manual commit for at-least-once processing:**

```go
func consumerManualCommit(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "orders",
        GroupID: "order-processor",
    })
    defer r.Close()

    for {
        msg, err := r.FetchMessage(ctx) // fetch but don't commit
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            continue
        }

        // Process the message
        if err := processOrder(msg.Value); err != nil {
            log.Printf("process error: %v — will retry", err)
            // Don't commit — message will be redelivered on restart
            continue
        }

        // Commit only after successful processing
        if err := r.CommitMessages(ctx, msg); err != nil {
            log.Printf("commit error: %v", err)
        }
    }
}

func processOrder(data []byte) error {
    var order struct {
        ID     string  `json:"id"`
        Amount float64 `json:"amount"`
    }
    if err := json.Unmarshal(data, &order); err != nil {
        return fmt.Errorf("bad message: %w", err) // don't retry bad messages
    }
    // ... do processing ...
    return nil
}
```

---

## 8. Mini Project: Order Event Pipeline

A complete pipeline: order service publishes events → multiple services consume them.

```go
package main

import (
    "context"
    "encoding/json"
    "fmt"
    "log"
    "math/rand"
    "sync"
    "time"

    kafka "github.com/segmentio/kafka-go"
)

type OrderEvent struct {
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    Amount    float64   `json:"amount"`
    Status    string    `json:"status"`
    OccurredAt time.Time `json:"occurred_at"`
}

func main() {
    ctx, cancel := context.WithCancel(context.Background())
    var wg sync.WaitGroup

    // Producer: simulate orders
    wg.Add(1)
    go func() {
        defer wg.Done()
        orderProducer(ctx)
    }()

    // Consumer 1: email service
    wg.Add(1)
    go func() {
        defer wg.Done()
        emailConsumer(ctx)
    }()

    // Consumer 2: analytics service
    wg.Add(1)
    go func() {
        defer wg.Done()
        analyticsConsumer(ctx)
    }()

    // Run for 10 seconds
    time.Sleep(10 * time.Second)
    cancel()
    wg.Wait()
}

func orderProducer(ctx context.Context) {
    w := &kafka.Writer{
        Addr:  kafka.TCP("localhost:9092"),
        Topic: "orders",
    }
    defer w.Close()

    ticker := time.NewTicker(200 * time.Millisecond)
    i := 0
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            event := OrderEvent{
                OrderID:    fmt.Sprintf("order-%d", i),
                UserID:     fmt.Sprintf("user-%d", rand.Intn(100)),
                Amount:     float64(10 + rand.Intn(990)),
                Status:     "placed",
                OccurredAt: time.Now(),
            }
            data, _ := json.Marshal(event)
            w.WriteMessages(ctx, kafka.Message{
                Key:   []byte(event.UserID),
                Value: data,
            })
            i++
        }
    }
}

func emailConsumer(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "orders",
        GroupID: "email-service",
    })
    defer r.Close()

    for {
        msg, err := r.ReadMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            continue
        }
        var event OrderEvent
        json.Unmarshal(msg.Value, &event)
        fmt.Printf("[Email] Sending confirmation to user %s for order %s ($%.2f)\n",
            event.UserID, event.OrderID, event.Amount)
    }
}

func analyticsConsumer(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "orders",
        GroupID: "analytics-service",
    })
    defer r.Close()

    var total float64
    count := 0
    for {
        msg, err := r.ReadMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                log.Printf("[Analytics] Processed %d orders, total $%.2f\n", count, total)
                return
            }
            continue
        }
        var event OrderEvent
        json.Unmarshal(msg.Value, &event)
        total += event.Amount
        count++
        if count%10 == 0 {
            fmt.Printf("[Analytics] %d orders, $%.2f total revenue\n", count, total)
        }
    }
}
```

---

## Summary

- Kafka is an append-only distributed log. Writes are always sequential, making them extremely fast.
- Topics split into partitions. Same key → same partition → ordered delivery per key.
- Use `RequiredAcks = RequireAll` for durability. Use `BatchSize` and `BatchTimeout` for throughput.
- `FetchMessage` + `CommitMessages` for manual at-least-once delivery.
- Multiple consumer groups each get all messages independently.

### Exercises

**Easy:** Create a `notifications` topic with 3 partitions. Write a producer that sends 100 notification events with random user IDs as keys. Write a consumer that prints each notification. Verify all 100 are received.

**Medium:** Implement a retry queue: if `processMessage` returns an error, write the message to a `notifications.retry` topic with an incremented `retry_count` header. A separate consumer reads from the retry topic with a 5-second delay before re-processing.

**Hard:** Implement exactly-once semantics using Kafka transactions: a producer that reads from `orders` topic, transforms each message, and writes to `order-totals` topic — exactly once, with transactional atomicity across both operations.
