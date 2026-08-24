# Chapter 32: Messaging at Scale — Kafka Internals, Consumer Groups & Guarantees

Kafka is the de-facto standard for event streaming in large-scale systems. Understanding Kafka at the depth that Uber, Netflix, and Stripe engineers use it is essential for senior system design interviews.

## Table of Contents

1. [Why Kafka Exists — Problems It Solves](#1-why-kafka-exists--problems-it-solves)
2. [Kafka Architecture](#2-kafka-architecture)
3. [Topics, Partitions & Offsets](#3-topics-partitions--offsets)
4. [Producer Internals](#4-producer-internals)
5. [Consumer Groups & Rebalancing](#5-consumer-groups--rebalancing)
6. [Delivery Guarantees](#6-delivery-guarantees)
7. [Kafka Patterns in Go](#7-kafka-patterns-in-go)
8. [Kafka vs RabbitMQ vs SQS](#8-kafka-vs-rabbitmq-vs-sqs)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. Why Kafka Exists — Problems It Solves

**Point-to-point integration problem:**
```
Before Kafka:
  Service A ──HTTP──▶ Service B
  Service A ──HTTP──▶ Service C
  Service A ──HTTP──▶ Service D
  
  Problem: A is coupled to B, C, D. If C is down, A must handle it. If new D is added, A changes.

With Kafka:
  Service A ──Kafka──▶ Topic "orders"
  Service B ────────────────────────────▶ reads "orders"
  Service C ────────────────────────────▶ reads "orders"
  Service D ────────────────────────────▶ reads "orders"
  
  A only knows about the topic. B, C, D subscribe independently.
  Adding D doesn't require changing A.
```

**Log replay:** Kafka stores messages on disk for configurable retention (days, weeks, forever). New consumers can replay the entire history. This is impossible with traditional message queues.

---

## 2. Kafka Architecture

```
Cluster:
  Multiple Brokers (servers), one is the Controller (leader for metadata)
  ZooKeeper (old) or KRaft (new) for cluster metadata and leader election

Broker:
  Stores partitions on disk as append-only log segments
  Handles producer writes and consumer reads

Replication:
  Each partition has 1 leader and N-1 followers (replicas)
  Producers write to leader, leader replicates to followers
  replication.factor = 3 means 2 followers for each partition leader
  
  ISR (In-Sync Replicas): followers that are caught up with the leader
  If leader fails, a new leader is elected from the ISR
```

---

## 3. Topics, Partitions & Offsets

```
Topic: "user-events"
  Partition 0: [msg0, msg1, msg2, msg3, ...]  ← append only, immutable
  Partition 1: [msg0, msg1, msg2, ...]
  Partition 2: [msg0, msg1, msg2, msg3, msg4, ...]

Offset: the position of a message within a partition (monotonically increasing)
  Partition 0, offset 3 → message msg3 in partition 0

Ordering:
  Kafka ONLY guarantees order within a partition
  Messages across partitions have NO ordering guarantee
  
Partition assignment (routing):
  Default: round-robin across partitions
  With key: hash(key) % num_partitions
    All messages with the same key go to the same partition → ordered per key
    Example: all events for user_id=123 → same partition → ordered for that user

Topic creation:
  kafka-topics.sh --create --topic user-events --partitions 6 --replication-factor 3
```

---

## 4. Producer Internals

```
Producer workflow:
  1. Producer calls send(record)
  2. Record goes into a per-partition buffer (accumulator)
  3. Sender thread batches records and sends to broker
  4. Broker appends to log, replicates to ISR
  5. Broker sends ack back to producer

Key configurations:
  acks=0:   fire-and-forget (fastest, can lose messages)
  acks=1:   leader acknowledges (fast, can lose if leader fails before replication)
  acks=all: all ISR must acknowledge (slowest, most durable — use for critical data)
  
  retries=3:     retry failed sends
  linger.ms=5:   wait up to 5ms to batch more records (reduces network requests)
  batch.size=16KB: max batch size per partition
  
  idempotent producer (enable.idempotence=true):
    Each message gets a sequence number
    Broker deduplicates retries → exactly-once semantics on broker side
```

---

## 5. Consumer Groups & Rebalancing

```
Consumer Group:
  Multiple consumer instances sharing a group.id
  Each partition is consumed by exactly one consumer in the group
  
  Topic: "orders" with 6 partitions
  Consumer Group "order-processor" with 3 consumers:
    Consumer 1 → Partition 0, 1
    Consumer 2 → Partition 2, 3
    Consumer 3 → Partition 4, 5
  
  Scaling up (add consumer 4):
    Rebalance: each consumer now handles ~1.5 partitions
    Consumer 1 → Partition 0
    Consumer 2 → Partition 2
    Consumer 3 → Partition 4
    Consumer 4 → Partition 1, 3, 5
  
  Rule: num_consumers > num_partitions → excess consumers are idle
  Rule: num_consumers < num_partitions → some consumers handle multiple partitions

Offset management:
  Consumers commit their current offset to Kafka (or ZooKeeper in old versions)
  On restart/rebalance, consumer resumes from last committed offset
  
  auto.commit.enable=true:  commits every auto.commit.interval.ms (5s)
  Manual commit (recommended): commit only after processing is complete
```

---

## 6. Delivery Guarantees

```
At-most-once:
  Commit offset before processing
  If processing fails, message is LOST (offset already advanced)
  Use when: high throughput, loss is OK (metrics, analytics)

At-least-once:
  Process message, then commit offset
  If commit fails, message is reprocessed (duplicate processing)
  Use when: you can handle duplicates, data loss is not OK
  Consumers must be IDEMPOTENT

Exactly-once:
  Kafka's transactional API + idempotent producer
  Producer gets a transaction ID, brokers deduplicate
  Consumer reads only committed offsets (isolation.level=read_committed)
  Use when: financial transactions, inventory updates, deduplicated writes are required
  Cost: significant throughput overhead (~30-50% reduction)
```

---

## 7. Kafka Patterns in Go

```go
import "github.com/segmentio/kafka-go"

// Producer
func publishEvent(ctx context.Context, topic string, key, value []byte) error {
    writer := &kafka.Writer{
        Addr:                   kafka.TCP("localhost:9092"),
        Topic:                  topic,
        Balancer:               &kafka.Hash{},       // consistent hashing by key
        RequiredAcks:           kafka.RequireAll,    // acks=all
        AllowAutoTopicCreation: false,
    }
    defer writer.Close()
    
    return writer.WriteMessages(ctx, kafka.Message{
        Key:   key,
        Value: value,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte("order.created")},
        },
    })
}

// Consumer with manual offset commit
func consumeOrders(ctx context.Context) error {
    reader := kafka.NewReader(kafka.ReaderConfig{
        Brokers:     []string{"localhost:9092"},
        Topic:       "orders",
        GroupID:     "order-processor",
        StartOffset: kafka.LastOffset,     // start from latest
        MinBytes:    10e3,                 // 10KB min fetch
        MaxBytes:    10e6,                 // 10MB max fetch
        MaxWait:     1 * time.Second,
    })
    defer reader.Close()
    
    for {
        // FetchMessage does NOT commit offset
        msg, err := reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil { return nil }
            return err
        }
        
        // Process the message
        if err := processOrder(ctx, msg.Value); err != nil {
            // Handle error: dead-letter queue, retry, etc.
            log.Printf("failed to process order: %v", err)
            // Don't commit — this message will be reprocessed
            continue
        }
        
        // Commit after successful processing (at-least-once)
        if err := reader.CommitMessages(ctx, msg); err != nil {
            log.Printf("failed to commit offset: %v", err)
        }
    }
}

// Consumer with graceful shutdown:
func runConsumer(ctx context.Context) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()
    
    // Handle OS signals
    sigChan := make(chan os.Signal, 1)
    signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
    go func() {
        <-sigChan
        cancel()
    }()
    
    if err := consumeOrders(ctx); err != nil && err != context.Canceled {
        log.Fatal(err)
    }
}
```

---

## 8. Kafka vs RabbitMQ vs SQS

| Feature | Kafka | RabbitMQ | SQS |
|---|---|---|---|
| **Model** | Distributed log (pull) | Message broker (push/pull) | Managed queue (pull) |
| **Ordering** | Per partition | Per queue | FIFO queues (with limits) |
| **Retention** | Configurable (days/weeks/forever) | Until consumed | Max 14 days |
| **Replay** | Yes — seek to any offset | No | No |
| **Throughput** | Very high (millions/sec) | High (100K+/sec) | High (managed) |
| **Consumer groups** | Multiple groups, each gets all messages | Competing consumers (one gets each) | Competing consumers |
| **Best for** | Event streaming, log aggregation, ETL | Task queues, RPC, complex routing | Simple queues on AWS |
| **Operational overhead** | High (cluster management) | Medium | None (managed) |

---

## 9. Interview Questions & Model Answers

**Q: How does Kafka ensure message ordering?**

"Kafka guarantees message ordering only within a single partition. When producing messages, you assign a partition key (e.g., user_id). All messages with the same key are consistently hashed to the same partition, ensuring they're ordered relative to each other. Consumers within a consumer group receive messages from their assigned partitions in order. The trade-off: if you need strict global ordering across all messages, you need only one partition — which eliminates parallelism."

**Q: Explain the difference between at-least-once and exactly-once delivery in Kafka.**

"At-least-once is the most common pattern: the consumer commits its offset only after successfully processing a message. If the consumer crashes between processing and committing, it reprocesses on restart. This means duplicates are possible — consumers must be idempotent to handle them. Exactly-once uses Kafka's transactional API: the producer gets a transactional ID and the broker deduplicates retries. The consumer uses `isolation.level=read_committed` to only read messages from committed transactions. This is stronger but has ~30-50% throughput overhead, so it's used primarily for financial and inventory updates where duplicates have business impact."

---

## Summary

- **Kafka:** distributed append-only log. Topics divided into partitions. Messages within a partition are ordered.
- **Partition key:** ensures all messages with the same key go to the same partition (ordering per entity).
- **Consumer group:** one consumer per partition. Adding consumers = parallel processing. More consumers than partitions = idle consumers.
- **acks=all + idempotent producer:** durable writes. At-least-once by default; exactly-once with transactions.
- **Offset commit after processing:** at-least-once. Before processing: at-most-once.
- **Replay:** unique to Kafka — consumers can seek to any historical offset and replay.
- Use Kafka for: event sourcing, activity tracking, log aggregation, decoupling services.
