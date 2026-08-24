# Chapter 93: Kafka — Fundamentals and Go Integration

Apache Kafka is a distributed event streaming platform. Unlike message queues (where messages are deleted after consumption), Kafka retains messages for a configurable period. Multiple independent consumer groups can each read the entire stream. This makes Kafka ideal for event-driven architectures where multiple services need the same events.

## Table of Contents

1. [Kafka Concepts](#1-kafka-concepts)
2. [Producer in Go](#2-producer-in-go)
3. [Consumer in Go](#3-consumer-in-go)
4. [Consumer Groups](#4-consumer-groups)
5. [Idempotent Producer and Exactly-Once](#5-idempotent-producer-and-exactly-once)
6. [Schema Registry and Avro](#6-dead-letter-queue)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Kafka Concepts

```
Topic: orders
  ├── Partition 0: [msg0, msg4, msg8, ...]
  ├── Partition 1: [msg1, msg5, msg9, ...]
  └── Partition 2: [msg2, msg6, msg10, ...]
```

| Concept | Description |
|---------|-------------|
| **Topic** | Named stream of messages (like a table, but append-only) |
| **Partition** | Ordered, immutable log within a topic. Parallelism unit. |
| **Offset** | Position of a message within a partition |
| **Producer** | Writes messages to topics |
| **Consumer** | Reads messages from topics |
| **Consumer Group** | Multiple consumers sharing partitions; each partition assigned to one consumer |
| **Broker** | A Kafka server node |
| **Retention** | How long messages are kept (default 7 days) |

### Key properties

- **Ordering**: guaranteed within a partition, not across partitions
- **Replayability**: consumer groups can reset their offset and re-process old messages
- **Multiple consumers**: two consumer groups reading the same topic each get all messages independently
- **Throughput**: millions of messages/second on commodity hardware

---

## 2. Producer in Go

Using `github.com/segmentio/kafka-go`:

```go
import "github.com/segmentio/kafka-go"

// Create a writer (producer)
func newProducer(brokers []string, topic string) *kafka.Writer {
    return &kafka.Writer{
        Addr:                   kafka.TCP(brokers...),
        Topic:                  topic,
        Balancer:               &kafka.LeastBytes{},    // partition selection
        RequiredAcks:           kafka.RequireAll,       // wait for all replicas
        Async:                  false,                  // synchronous writes
        BatchSize:              100,                    // batch up to 100 messages
        BatchTimeout:           10 * time.Millisecond,
        WriteTimeout:           10 * time.Second,
        AllowAutoTopicCreation: true,
    }
}

type OrderEventProducer struct {
    writer *kafka.Writer
}

type OrderEvent struct {
    EventType string          `json:"event_type"`
    OrderID   string          `json:"order_id"`
    Payload   json.RawMessage `json:"payload"`
    OccuredAt time.Time       `json:"occurred_at"`
}

func (p *OrderEventProducer) Publish(ctx context.Context, event OrderEvent) error {
    data, err := json.Marshal(event)
    if err != nil { return fmt.Errorf("marshal: %w", err) }
    
    return p.writer.WriteMessages(ctx, kafka.Message{
        Key:   []byte(event.OrderID), // same key → same partition → ordered per order
        Value: data,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte(event.EventType)},
            {Key: "content-type", Value: []byte("application/json")},
        },
    })
}

// Batch publish
func (p *OrderEventProducer) PublishBatch(ctx context.Context, events []OrderEvent) error {
    msgs := make([]kafka.Message, len(events))
    for i, event := range events {
        data, _ := json.Marshal(event)
        msgs[i] = kafka.Message{
            Key:   []byte(event.OrderID),
            Value: data,
        }
    }
    return p.writer.WriteMessages(ctx, msgs...)
}

func (p *OrderEventProducer) Close() error {
    return p.writer.Close()
}
```

### Partition key strategy

```go
// Use a meaningful partition key to ensure ordering per entity
// Same order_id → same partition → events arrive in order

// Hash-based: kafka-go's default — hash(key) % numPartitions
msg := kafka.Message{Key: []byte(orderID), Value: data}

// Manual partition assignment:
msg := kafka.Message{
    Partition: int(fnv32(orderID) % numPartitions),
    Key:       []byte(orderID),
    Value:     data,
}
```

---

## 3. Consumer in Go

```go
func newConsumer(brokers []string, topic, groupID string) *kafka.Reader {
    return kafka.NewReader(kafka.ReaderConfig{
        Brokers:         brokers,
        Topic:           topic,
        GroupID:         groupID,           // consumer group ID
        MinBytes:        10e3,              // 10KB — fetch when at least this much available
        MaxBytes:        10e6,              // 10MB max per fetch
        MaxWait:         500 * time.Millisecond,
        CommitInterval:  1 * time.Second,   // auto-commit offset every second
        StartOffset:     kafka.LastOffset,  // start from newest (use FirstOffset to replay all)
    })
}

type OrderEventConsumer struct {
    reader  *kafka.Reader
    handler EventHandler
    logger  *slog.Logger
}

func (c *OrderEventConsumer) Run(ctx context.Context) error {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil { return nil } // graceful shutdown
            return fmt.Errorf("fetch: %w", err)
        }
        
        if err := c.processMessage(ctx, msg); err != nil {
            c.logger.Error("process message failed",
                "topic", msg.Topic,
                "partition", msg.Partition,
                "offset", msg.Offset,
                "err", err,
            )
            // Don't commit — message will be re-delivered
            continue
        }
        
        // Commit offset AFTER successful processing
        if err := c.reader.CommitMessages(ctx, msg); err != nil {
            c.logger.Warn("commit failed", "err", err)
        }
    }
}

func (c *OrderEventConsumer) processMessage(ctx context.Context, msg kafka.Message) error {
    var event OrderEvent
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }
    
    return c.handler.Handle(ctx, event)
}

func (c *OrderEventConsumer) Close() error {
    return c.reader.Close()
}
```

---

## 4. Consumer Groups

Each partition is assigned to exactly one consumer in a group. When consumers join or leave, Kafka **rebalances** partitions.

```
Topic "orders" with 6 partitions, Consumer Group "inventory-service":

  3 consumers:
  Consumer A → partitions [0, 1]
  Consumer B → partitions [2, 3]
  Consumer C → partitions [4, 5]

  If Consumer B dies:
  Consumer A → partitions [0, 1, 2]
  Consumer C → partitions [3, 4, 5]

  Two consumer groups reading the same topic:
  "inventory-service" reads all 6 partitions independently
  "notification-service" also reads all 6 partitions independently
```

```go
// Multiple consumers in the same group — process in parallel
func startConsumerGroup(ctx context.Context, brokers []string, topic, groupID string, workers int) {
    var wg sync.WaitGroup
    for i := range workers {
        wg.Add(1)
        go func(workerID int) {
            defer wg.Done()
            
            reader := kafka.NewReader(kafka.ReaderConfig{
                Brokers: brokers,
                Topic:   topic,
                GroupID: groupID, // same group ID — Kafka distributes partitions among these readers
            })
            defer reader.Close()
            
            consumer := &OrderEventConsumer{reader: reader}
            if err := consumer.Run(ctx); err != nil {
                log.Printf("worker %d error: %v", workerID, err)
            }
        }(i)
    }
    wg.Wait()
}
```

---

## 5. Idempotent Producer and Exactly-Once

```go
// Idempotent producer: Kafka deduplicates retries using producer ID + sequence number
writer := &kafka.Writer{
    Addr:         kafka.TCP("localhost:9092"),
    Topic:        "orders",
    RequiredAcks: kafka.RequireAll,
    // EnableIdempotency: kafka-go enables this by default with RequireAll
}

// Idempotent CONSUMER: handle re-delivery gracefully
type IdempotentOrderHandler struct {
    processed *sync.Map // in production: use Redis or a DB table
    delegate  EventHandler
}

func (h *IdempotentOrderHandler) Handle(ctx context.Context, event OrderEvent) error {
    // Check if we've already processed this event
    key := fmt.Sprintf("%s:%s", event.EventType, event.OrderID)
    if _, loaded := h.processed.LoadOrStore(key, true); loaded {
        return nil // already processed, skip
    }
    return h.delegate.Handle(ctx, event)
}

// In production: use a DB table for deduplication
func processIdempotently(ctx context.Context, db *sqlx.DB, eventID string, fn func() error) error {
    // Try to insert the event ID; if it already exists, skip processing
    _, err := db.ExecContext(ctx, `
        INSERT INTO processed_events (event_id, processed_at)
        VALUES ($1, NOW())
        ON CONFLICT (event_id) DO NOTHING`,
        eventID,
    )
    if err != nil { return err }
    
    // Check if we inserted (new) or skipped (duplicate)
    var count int
    db.QueryRowContext(ctx,
        "SELECT COUNT(*) FROM processed_events WHERE event_id=$1 AND processed_at > NOW()-INTERVAL '1s'",
        eventID,
    ).Scan(&count)
    
    if count == 0 { return nil } // duplicate, skip
    return fn()
}
```

---

## 6. Dead Letter Queue

When a message can't be processed after N retries, send it to a DLQ for manual investigation:

```go
type DLQConsumer struct {
    reader   *kafka.Reader
    dlqWriter *kafka.Writer
    handler  EventHandler
    maxRetries int
}

func (c *DLQConsumer) processWithRetry(ctx context.Context, msg kafka.Message) error {
    var lastErr error
    for attempt := 1; attempt <= c.maxRetries; attempt++ {
        if err := c.handler.Handle(ctx, msg); err != nil {
            lastErr = err
            time.Sleep(time.Duration(attempt) * 500 * time.Millisecond)
            continue
        }
        return nil
    }
    
    // Send to DLQ
    return c.dlqWriter.WriteMessages(ctx, kafka.Message{
        Key: msg.Key,
        Value: msg.Value,
        Headers: append(msg.Headers,
            kafka.Header{Key: "dlq-reason", Value: []byte(lastErr.Error())},
            kafka.Header{Key: "dlq-retries", Value: []byte(strconv.Itoa(c.maxRetries))},
            kafka.Header{Key: "original-topic", Value: []byte(msg.Topic)},
        ),
    })
}
```

---

## Summary

- Kafka retains messages (default 7 days) — multiple consumer groups read independently
- **Partition key**: same key → same partition → guaranteed ordering per entity
- **Consumer group**: partitions are assigned across consumers; Kafka rebalances on topology change
- Commit offset **after** successful processing — not before — to avoid message loss
- **Idempotent processing**: deduplicate via a `processed_events` table or Redis set
- **Dead Letter Queue**: after N failures, send to a DLQ topic for investigation
- `kafka-go` is simpler than `sarama` and sufficient for most production uses

## Exercises

### Easy
1. Create a Kafka topic `user-events` with 3 partitions. Write a producer that sends 100 events with user IDs as partition keys. Verify that events for the same user land on the same partition.
2. Write a consumer that reads from `user-events`, prints each message, and commits the offset manually after processing.
3. Start two consumers in the same consumer group and verify that they split the 3 partitions between them (not both reading all 3).

### Medium
4. Implement a **retry consumer**: if `handler.Handle()` returns an error, retry up to 3 times with exponential backoff. If all retries fail, publish the message to `user-events.dlq` with an error header.
5. Build a **Kafka to PostgreSQL sink**: consume events from `order-events`, parse each, and upsert into `order_summaries`. Use the idempotency pattern to safely handle re-delivery.
6. Implement a **lag monitor**: read consumer group offsets and log partition lag (latest offset - committed offset) every 30 seconds. Alert if any partition lag exceeds 1000.

### Hard
7. Implement an **exactly-once processor** using Kafka transactions: begin a transaction, consume a message, write to PostgreSQL, produce a result message, commit the transaction atomically. If anything fails, abort — no partial writes.
8. Build a **multi-topic fan-in consumer**: read from `orders`, `inventory`, and `payments` topics simultaneously. For each order ID, collect all three events, then join them into a single `OrderFullView` record. Use a `sync.Map` keyed by order ID with a TTL to expire incomplete joins.
