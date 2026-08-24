# Chapter 55: Advanced Kafka Patterns in Go

Basic Kafka is produce and consume. Advanced Kafka is what separates toy systems from production ones: exactly-once semantics, schema evolution, Kafka Streams processing, and monitoring.

## Table of Contents

1. Schema Registry and Avro
2. Exactly-Once Semantics (Transactions)
3. Kafka Streams-Style Processing in Go
4. Monitoring and Lag Tracking
5. Consumer Error Handling and Retry
6. Production Configuration
7. Exercises

---

## 1. Schema Registry and Avro

**The problem:** Your producer publishes JSON. Three months later, you add a field. Now all your consumers break because they don't expect the new field. How do you evolve schemas without breaking consumers?

**Schema Registry:** A central service that stores message schemas. Producers register schemas; consumers look them up. Schema evolution rules enforce backward/forward compatibility.

```
Producer registers schema:
  {
    "type": "record",
    "name": "OrderPlaced",
    "fields": [
      {"name": "order_id", "type": "string"},
      {"name": "amount",   "type": "double"}
    ]
  }

3 months later, producer adds a field:
  {
    "fields": [
      {"name": "order_id", "type": "string"},
      {"name": "amount",   "type": "double"},
      {"name": "currency", "type": "string", "default": "USD"}  ← new field with default
    ]
  }
```

Old consumers (that don't know about `currency`) still work because the new field has a default. This is **backward compatibility**.

```bash
# Start Confluent Platform (Kafka + Schema Registry)
docker run -d --name schema-registry \
  -p 8081:8081 \
  -e SCHEMA_REGISTRY_HOST_NAME=schema-registry \
  -e SCHEMA_REGISTRY_KAFKASTORE_BOOTSTRAP_SERVERS=kafka:9092 \
  confluentinc/cp-schema-registry:7.5.0
```

```go
// Simplified schema-aware producer (without full Avro library)
type SchemaAwareProducer struct {
    writer    *kafka.Writer
    schemaID  int
    schemaReg string
}

// In production: use github.com/riferrei/srclient for full Schema Registry support
func (p *SchemaAwareProducer) Produce(data interface{}) error {
    payload, _ := json.Marshal(data)

    // Confluent wire format: 1 magic byte + 4 byte schema ID + payload
    msg := make([]byte, 5+len(payload))
    msg[0] = 0 // magic byte
    binary.BigEndian.PutUint32(msg[1:5], uint32(p.schemaID))
    copy(msg[5:], payload)

    return p.writer.WriteMessages(context.Background(), kafka.Message{Value: msg})
}
```

---

## 2. Exactly-Once Semantics (Transactions)

Kafka transactions enable atomic writes across multiple topics:

```go
import "github.com/twmb/franz-go/pkg/kgo"

func exactlyOnceProcessor(ctx context.Context) {
    // Transactional producer
    client, err := kgo.NewClient(
        kgo.SeedBrokers("localhost:9092"),
        kgo.TransactionalID("processor-1"), // unique ID per process instance
        kgo.RequiredAcks(kgo.AllISRAcks()),
    )
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()

    // Read from input topic, process, write to output topic — atomically
    for {
        fetches := client.PollFetches(ctx)
        if errs := fetches.Errors(); len(errs) > 0 {
            log.Println(errs)
            continue
        }

        // Begin transaction
        if err := client.BeginTransaction(); err != nil {
            log.Println("begin txn:", err)
            continue
        }

        var msgs []kgo.Record

        fetches.EachRecord(func(rec *kgo.Record) {
            // Process the message
            processed := transform(rec.Value)

            // Write to output topic (within transaction)
            msgs = append(msgs, kgo.Record{
                Topic: "orders-processed",
                Key:   rec.Key,
                Value: processed,
            })
        })

        // Write all results
        client.ProduceSync(ctx, msgs...).FirstErr()

        // Commit the transaction AND the consumer offsets atomically
        if err := client.EndTransaction(ctx, kgo.TryCommit); err != nil {
            client.EndTransaction(ctx, kgo.TryAbort)
            log.Println("txn failed, aborted:", err)
        }
    }
}

func transform(data []byte) []byte {
    // Your processing logic here
    return data
}
```

**What exactly-once means here:**
- If the process crashes after writing to output but before committing, the transaction is rolled back.
- No duplicates in the output topic.
- No messages lost.
- The consumer offset and the output write happen atomically.

---

## 3. Kafka Streams-Style Processing in Go

Kafka Streams is a Java library for stateful stream processing. In Go, we implement similar patterns manually.

**Windowed aggregation: count orders per minute**

```go
type WindowAggregator struct {
    mu      sync.Mutex
    windows map[int64]int64  // minute → count
}

func NewWindowAggregator() *WindowAggregator {
    return &WindowAggregator{windows: make(map[int64]int64)}
}

func (a *WindowAggregator) Add(ts time.Time) {
    a.mu.Lock()
    defer a.mu.Unlock()
    // Round to minute
    minute := ts.Truncate(time.Minute).Unix()
    a.windows[minute]++
}

func (a *WindowAggregator) GetAndClear(olderThan time.Time) map[int64]int64 {
    a.mu.Lock()
    defer a.mu.Unlock()

    threshold := olderThan.Unix()
    result := make(map[int64]int64)
    for ts, count := range a.windows {
        if ts < threshold {
            result[ts] = count
            delete(a.windows, ts)
        }
    }
    return result
}

// Stream processor
func windowedOrderCounter(ctx context.Context) {
    r := kafka.NewReader(kafka.ReaderConfig{
        Brokers: []string{"localhost:9092"},
        Topic:   "orders",
        GroupID: "order-counter",
    })
    defer r.Close()

    agg := NewWindowAggregator()

    // Flush completed windows every 10 seconds
    go func() {
        ticker := time.NewTicker(10 * time.Second)
        for range ticker.C {
            completed := agg.GetAndClear(time.Now().Add(-2 * time.Minute))
            for minute, count := range completed {
                ts := time.Unix(minute, 0).Format("15:04")
                fmt.Printf("Window %s: %d orders\n", ts, count)
                // Write to "order-metrics" topic or a database
            }
        }
    }()

    for {
        msg, err := r.ReadMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            continue
        }

        var order struct {
            Timestamp time.Time `json:"timestamp"`
        }
        if err := json.Unmarshal(msg.Value, &order); err != nil {
            continue
        }
        agg.Add(order.Timestamp)
    }
}
```

---

## 4. Monitoring and Lag Tracking

**Consumer lag** is the most important metric: how far behind is your consumer?

```go
type LagMonitor struct {
    admin  *kafka.Client
    group  string
    topics []string
}

func (m *LagMonitor) GetLag(ctx context.Context) (map[string]int64, error) {
    // Get the latest offsets for each partition
    topicOffsets, err := m.admin.ListOffsets(ctx,
        &kafka.ListOffsetsRequest{
            Topics: map[string][]kafka.OffsetRequest{
                "orders": {{Partition: 0}, {Partition: 1}, {Partition: 2}},
            },
        },
    )
    if err != nil {
        return nil, err
    }

    // Get consumer group offsets
    groupOffsets, err := m.admin.OffsetFetchRequest(ctx, m.group,
        []kafka.Topic{{Name: "orders", Partitions: []int{0, 1, 2}}},
    )
    if err != nil {
        return nil, err
    }

    lag := make(map[string]int64)
    for topic, partitions := range topicOffsets.Topics {
        for _, p := range partitions {
            consumerOffset := groupOffsets.Topics[topic][p.Partition]
            partitionLag := p.LastOffset - consumerOffset
            lag[fmt.Sprintf("%s[%d]", topic, p.Partition)] = partitionLag
        }
    }
    return lag, nil
}

// In production: expose lag as a Prometheus metric
// and alert when lag exceeds a threshold (e.g., > 100K messages)
```

---

## 5. Consumer Error Handling and Retry

```go
type RetryConsumer struct {
    reader      *kafka.Reader
    retryWriter *kafka.Writer
    dlqWriter   *kafka.Writer
    maxRetries  int
}

type MessageWithMetadata struct {
    OriginalTopic string    `json:"original_topic"`
    RetryCount    int       `json:"retry_count"`
    Error         string    `json:"error"`
    Payload       []byte    `json:"payload"`
    FailedAt      time.Time `json:"failed_at"`
}

func (c *RetryConsumer) Run(ctx context.Context, process func([]byte) error) {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil {
                return
            }
            continue
        }

        err = process(msg.Value)
        if err == nil {
            c.reader.CommitMessages(ctx, msg)
            continue
        }

        // Get retry count from header
        retryCount := getRetryCount(msg.Headers)

        if retryCount >= c.maxRetries {
            // Send to DLQ
            meta := MessageWithMetadata{
                OriginalTopic: msg.Topic,
                RetryCount:    retryCount,
                Error:         err.Error(),
                Payload:       msg.Value,
                FailedAt:      time.Now(),
            }
            dlqData, _ := json.Marshal(meta)
            c.dlqWriter.WriteMessages(ctx, kafka.Message{
                Key:   msg.Key,
                Value: dlqData,
            })
            log.Printf("Message sent to DLQ after %d retries: %v", retryCount, err)
        } else {
            // Send to retry topic with incremented count
            headers := incrementRetryCount(msg.Headers, retryCount+1)
            c.retryWriter.WriteMessages(ctx, kafka.Message{
                Key:     msg.Key,
                Value:   msg.Value,
                Headers: headers,
            })
        }

        c.reader.CommitMessages(ctx, msg)
    }
}

func getRetryCount(headers []kafka.Header) int {
    for _, h := range headers {
        if h.Key == "retry-count" {
            count, _ := strconv.Atoi(string(h.Value))
            return count
        }
    }
    return 0
}

func incrementRetryCount(headers []kafka.Header, count int) []kafka.Header {
    for i, h := range headers {
        if h.Key == "retry-count" {
            headers[i].Value = []byte(strconv.Itoa(count))
            return headers
        }
    }
    return append(headers, kafka.Header{
        Key:   "retry-count",
        Value: []byte(strconv.Itoa(count)),
    })
}
```

---

## 6. Production Configuration

```go
// Producer: high-throughput, durable
writer := &kafka.Writer{
    Addr:  kafka.TCP("kafka1:9092", "kafka2:9092", "kafka3:9092"),
    Topic: "orders",

    // Durability
    RequiredAcks: kafka.RequireAll, // all ISR replicas

    // Throughput
    BatchSize:    1000,              // batch up to 1000 messages
    BatchTimeout: 5 * time.Millisecond, // or flush every 5ms

    // Reliability
    MaxAttempts:  10,
    WriteTimeout: 30 * time.Second,

    // Compression (lz4 is fastest, zstd has best compression ratio)
    Compression: kafka.Lz4,

    // Transactions
    // TransactionalID: "my-producer-1", // for exactly-once
}

// Consumer: reliable processing
reader := kafka.NewReader(kafka.ReaderConfig{
    Brokers:        []string{"kafka1:9092", "kafka2:9092", "kafka3:9092"},
    Topic:          "orders",
    GroupID:        "order-processor",
    MinBytes:       1,
    MaxBytes:       50 * 1024 * 1024, // 50 MB
    MaxWait:        500 * time.Millisecond,
    ReadBackoffMin: 100 * time.Millisecond, // backoff on errors
    ReadBackoffMax: 10 * time.Second,
    CommitInterval: 0, // disable auto-commit, use manual CommitMessages
    StartOffset:    kafka.LastOffset, // start from current end (new consumer)
    ErrorLogger:    kafka.LoggerFunc(func(msg string, args ...interface{}) {
        log.Printf("KAFKA ERROR: "+msg, args...)
    }),
})
```

---

## Summary

- Schema Registry prevents breaking changes when evolving message schemas. Always add default values to new fields.
- Exactly-once semantics require `TransactionalID` on the producer and atomic offset+message commits.
- Windowed aggregation: bucket by time window, flush completed windows periodically.
- Consumer lag is the critical metric: `latestOffset - consumerOffset`. Alert when lag grows.
- Retry with exponential backoff + DLQ (Dead Letter Queue) for unprocessable messages.

### Exercises

**Easy:** Add a consumer lag monitor to your order processing consumer. Print the lag for each partition every 5 seconds.

**Medium:** Implement windowed aggregation: for a stream of `{user_id, action, timestamp}` events, compute the count of actions per user per hour. Emit a "user-hourly-stats" message when a window closes.

**Hard:** Implement exactly-once semantics for an order status pipeline: read from "orders" topic, compute order status updates, write to "order-status" topic — atomically. Verify by killing the process mid-processing and confirming no duplicates in the output topic.
