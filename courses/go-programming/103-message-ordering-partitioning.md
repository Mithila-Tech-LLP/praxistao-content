# Chapter 103: Message Ordering and Partitioning

Out-of-order message processing is a silent killer in event-driven systems. Unlike crashes — which produce obvious errors — ordering violations produce subtle state corruption that can take hours to detect. This chapter covers why order matters, how partitioning enforces it, and how to handle the cases where you can't guarantee it at the broker level.

## Table of Contents

1. [Why Ordering Matters](#1-why-ordering-matters)
2. [Total Order vs Partial Order](#2-total-order-vs-partial-order)
3. [Kafka Partitioning](#3-kafka-partitioning)
4. [Partition Key Selection](#4-partition-key-selection)
5. [The Cardinality Problem](#5-the-cardinality-problem)
6. [FIFO Queues and Group IDs](#6-fifo-queues-and-group-ids)
7. [Reordering at the Consumer](#7-reordering-at-the-consumer)
8. [Ordering vs Throughput Trade-Off](#8-ordering-vs-throughput-trade-off)
9. [Handling Out-of-Order Messages](#9-handling-out-of-order-messages)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Ordering Matters

Consider an order lifecycle. Events arrive in this logical sequence:

```
OrderCreated → OrderConfirmed → OrderShipped → OrderCancelled
```

Now imagine the consumer receives them out of order:

```
Consumer receives: OrderCancelled, then OrderCreated

State after OrderCancelled:   order.status = "cancelled"
State after OrderCreated:     order.status = "pending"   ← overwrites!

Result: a cancelled order appears as pending.
        The customer might get a second shipment.
```

More subtle: an inventory service consuming these events:

```go
// Correct order
OnOrderConfirmed:  inventory.Reserve(sku, qty)   // reserve 5 units
OnOrderCancelled:  inventory.Release(sku, qty)   // return 5 units
// Final: inventory is back to where it was. ✓

// Out-of-order
OnOrderCancelled:  inventory.Release(sku, qty)   // return 5 units (never reserved!)
OnOrderConfirmed:  inventory.Reserve(sku, qty)   // reserve 5 units
// Final: inventory is 5 units lower than reality. ✗
//        Stock level is wrong — phantom reservation never gets released.
```

The damage compounds over time: each ordering violation leaves the system in a subtly wrong state, and diagnosing the root cause requires replaying the event history.

---

## 2. Total Order vs Partial Order

```
Total order: all messages across all producers are ordered globally.
             msg1 < msg2 < msg3 < ... (one global sequence)

Partial order: messages are ordered within a group (per entity).
               Order A: A1 < A2 < A3
               Order B: B1 < B2           (A and B events interleave freely)
               Order C: C1 < C2 < C3
```

Total order is expensive: it requires a single writer or a consensus protocol (like Raft) that all writes flow through. This becomes a throughput bottleneck.

Partial order is almost always sufficient. Your business logic cares that events for **one entity** are ordered — it doesn't care about the relative ordering of events from different entities:

```
Payment for order-A processed before payment for order-B? Irrelevant.
Both payments must be processed after their respective order was confirmed. This is partial order.
```

Kafka implements partial order naturally: ordering is guaranteed within a partition, and you assign all events for the same entity to the same partition.

---

## 3. Kafka Partitioning

Kafka distributes messages across partitions. Each partition is an ordered, append-only log consumed sequentially by one consumer in a group.

```
Topic: orders  (6 partitions)

Producer writes order-A events:
  msg{key:"order-A", ...} → hash("order-A") % 6 = partition 2

Producer writes order-B events:
  msg{key:"order-B", ...} → hash("order-B") % 6 = partition 0

Partition 0: [B1, B2, B3, B7]           ← only order-B events
Partition 2: [A1, A2, A3, A4, A5]       ← only order-A events
Partition 4: [C1, C2, D1, D2, C3]       ← order-C and order-D events

Consumer reading partition 2 sees: A1, A2, A3, A4, A5 — in order. ✓
```

The partition key is the mechanism. In Go with `kafka-go`:

```go
import "github.com/segmentio/kafka-go"

type OrderEventProducer struct {
    writer *kafka.Writer
}

type OrderEvent struct {
    ID        string          `json:"id"`
    Type      string          `json:"type"` // OrderCreated, OrderConfirmed, etc.
    OrderID   string          `json:"order_id"`
    Payload   json.RawMessage `json:"payload"`
    Version   int             `json:"version"` // monotonically increasing per order
    OccurredAt time.Time      `json:"occurred_at"`
}

func (p *OrderEventProducer) Publish(ctx context.Context, event OrderEvent) error {
    data, err := json.Marshal(event)
    if err != nil {
        return fmt.Errorf("marshal: %w", err)
    }

    return p.writer.WriteMessages(ctx, kafka.Message{
        // Key determines partition assignment.
        // All events for the same order go to the same partition → ordered.
        Key:   []byte(event.OrderID),
        Value: data,
        Headers: []kafka.Header{
            {Key: "event-type", Value: []byte(event.Type)},
            {Key: "version", Value: []byte(strconv.Itoa(event.Version))},
        },
    })
}

// Publishing a sequence of events for order-123:
//   All three use Key: []byte("order-123")
//   → all land on the same partition
//   → consumed in the order they were produced
func (p *OrderEventProducer) PublishOrderLifecycle(ctx context.Context, orderID string) error {
    events := []OrderEvent{
        {ID: newUUID(), Type: "OrderCreated",   OrderID: orderID, Version: 1},
        {ID: newUUID(), Type: "OrderConfirmed", OrderID: orderID, Version: 2},
        {ID: newUUID(), Type: "OrderShipped",   OrderID: orderID, Version: 3},
    }
    for _, e := range events {
        if err := p.Publish(ctx, e); err != nil {
            return err
        }
    }
    return nil
}
```

**Why this works**: `kafka-go` uses `hash(key) % numPartitions` to route messages. The same key always hashes to the same partition. Within a partition, messages are appended sequentially and consumed in that sequence.

---

## 4. Partition Key Selection

Choosing the wrong key is a common mistake. The rule: **key on the entity whose events must be ordered**.

```
Scenario: Order processing system
  Events: OrderCreated, OrderConfirmed, OrderShipped, OrderCancelled

  Must be ordered per: order (not per user, not per product)
  → Key: orderID ✓

Scenario: User activity stream
  Events: UserSignedUp, UserUpdatedProfile, UserDeleted
  Must be ordered per: user
  → Key: userID ✓

Scenario: Inventory updates
  Events: StockAdded, StockReserved, StockReleased
  Must be ordered per: SKU (not per order, not per warehouse)
  → Key: sku ✓

Scenario: Chat messages
  Events: MessageSent, MessageEdited, MessageDeleted
  Must be ordered per: conversation
  → Key: conversationID ✓  (not messageID — that would scatter messages randomly)
```

Common mistakes:

```go
// WRONG: random key — ordering is completely lost
msg := kafka.Message{
    Key:   []byte(uuid.New().String()), // random per message
    Value: data,
}

// WRONG: no key — kafka-go uses round-robin → no ordering guarantee
msg := kafka.Message{
    Value: data, // Key is nil
}

// WRONG: too coarse a key — all events on one partition = hot spot
msg := kafka.Message{
    Key:   []byte("orders"), // all orders → one partition
    Value: data,
}

// CORRECT: key on the entity
msg := kafka.Message{
    Key:   []byte(event.OrderID),
    Value: data,
}
```

---

## 5. The Cardinality Problem

More partitions = more parallelism = more throughput. But partition count is fixed at topic creation (changing it re-hashes keys and breaks ordering).

```
Throughput calculation:

  Peak load: 10,000 order events/second
  Target per-partition throughput: 2,000 msg/sec (conservative — leave headroom)
  
  Minimum partitions = 10,000 / 2,000 = 5
  With 2x headroom: 10 partitions
  Round up to a number divisible by consumer count: 12 partitions

  Consumer group has 4 pods → 3 partitions per pod
```

Hot spot problem: if your key cardinality is too low, some partitions get much more traffic than others.

```
Bad: 3 partition keys for 12 partitions
  Key "region:us"  → partition 1 → 90% of traffic
  Key "region:eu"  → partition 5 → 8% of traffic
  Key "region:ap"  → partition 9 → 2% of traffic

  Partition 1 consumer is overwhelmed.
  Partitions 2-12 are idle.
  ← This is a hot spot.

Good: high-cardinality keys
  orderID, userID, sessionID — millions of distinct values
  → traffic distributed roughly evenly across all partitions
```

When you must use a low-cardinality key, add a sub-key to distribute load while preserving ordering within the group:

```go
// Composite key: guarantees order within a category, distributes across partitions
func partitionKey(category, entityID string) []byte {
    // All events for the same entityID go to the same partition
    return []byte(fmt.Sprintf("%s:%s", category, entityID))
}
```

---

## 6. FIFO Queues and Group IDs

Amazon SQS FIFO queues use a **Message Group ID** — equivalent to Kafka's partition key. All messages in the same group are delivered in order to a single consumer at a time.

```go
import "github.com/aws/aws-sdk-go-v2/service/sqs"

type SQSFIFOProducer struct {
    client   *sqs.Client
    queueURL string
}

func (p *SQSFIFOProducer) Send(ctx context.Context, event OrderEvent) error {
    data, err := json.Marshal(event)
    if err != nil { return err }

    _, err = p.client.SendMessage(ctx, &sqs.SendMessageInput{
        QueueUrl:    &p.queueURL,
        MessageBody: aws.String(string(data)),

        // MessageGroupId is the FIFO equivalent of Kafka's partition key.
        // All messages with the same group ID are delivered in order.
        MessageGroupId: aws.String(event.OrderID),

        // Deduplication ID: SQS deduplicates messages with the same ID
        // within a 5-minute window (at-least-once + dedup = effectively exactly-once)
        MessageDeduplicationId: aws.String(event.ID),
    })
    return err
}
```

The conceptual mapping:

```
Kafka                   SQS FIFO
──────────────────      ─────────────────────
Partition key           MessageGroupId
Partition               Message Group
Consumer group member   In-flight lock on the group
```

Both enforce: within a group/partition, one consumer processes messages in order, and no other consumer takes from that group/partition while one is active.

---

## 7. Reordering at the Consumer

Sometimes ordering at the broker level isn't available (RabbitMQ, standard SQS) or messages genuinely arrive out of order despite a partition key (two producers writing to the same key concurrently). In these cases, buffer and sort at the consumer.

```go
// SequenceBuffer collects messages for an entity until they can be
// processed in order, identified by a monotonically increasing sequence number.
type SequenceBuffer struct {
    mu       sync.Mutex
    buffers  map[string]*entityBuffer // entityID → buffer
    handler  func(ctx context.Context, msg OrderEvent) error
}

type entityBuffer struct {
    nextExpected int
    pending      map[int]OrderEvent // sequence → event
    lastActivity time.Time
}

func (b *SequenceBuffer) Push(ctx context.Context, event OrderEvent) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    buf, ok := b.buffers[event.OrderID]
    if !ok {
        buf = &entityBuffer{
            nextExpected: 1,
            pending:      make(map[int]OrderEvent),
        }
        b.buffers[event.OrderID] = buf
    }

    buf.lastActivity = time.Now()

    if event.Version < buf.nextExpected {
        // Already processed — duplicate or old message. Skip.
        return nil
    }

    if event.Version == buf.nextExpected {
        // In-order delivery — process immediately
        if err := b.handler(ctx, event); err != nil {
            return err
        }
        buf.nextExpected++

        // Drain any buffered messages that are now in sequence
        for {
            next, ok := buf.pending[buf.nextExpected]
            if !ok { break }
            if err := b.handler(ctx, next); err != nil {
                return err
            }
            delete(buf.pending, buf.nextExpected)
            buf.nextExpected++
        }
    } else {
        // Out-of-order — buffer it until the gap fills
        buf.pending[event.Version] = event
    }

    return nil
}

// Flush processes all buffered messages for an entity regardless of gaps.
// Use on a timeout when you suspect a sequence number was permanently lost.
func (b *SequenceBuffer) Flush(ctx context.Context, orderID string) error {
    b.mu.Lock()
    defer b.mu.Unlock()

    buf, ok := b.buffers[orderID]
    if !ok { return nil }

    // Sort pending keys and process in order
    versions := make([]int, 0, len(buf.pending))
    for v := range buf.pending { versions = append(versions, v) }
    sort.Ints(versions)

    for _, v := range versions {
        if err := b.handler(ctx, buf.pending[v]); err != nil {
            return err
        }
    }

    delete(b.buffers, orderID)
    return nil
}
```

---

## 8. Ordering vs Throughput Trade-Off

More partitions = more parallelism = more throughput. But ordering scope shrinks as partition count grows.

```
                    Partitions
                    ┌────────────────────────────────────────┐
                    │  1      2      6      12     100        │
Throughput          │  low    ──────────────────► high       │
                    │                                        │
Order guarantee     │  global ──────────────────► per-entity │
                    │                                        │
Consumer parallelism│  1      2      6      12     100       │
                    │                                        │
Rebalance impact    │  low    ──────────────────► high       │
                    └────────────────────────────────────────┘

With 1 partition:
  - All messages globally ordered
  - Only 1 consumer can process at a time
  - Throughput bottleneck

With 100 partitions:
  - Messages ordered per-entity (order-123 events always land on partition 37)
  - 100 consumers process in parallel
  - High throughput
  - Rebalancing when a consumer dies reassigns ~50 partitions → brief pause
```

Recommendation for most systems: start with 12–24 partitions (even if you have fewer consumers today), using your entity ID as the key. You can scale consumers up to partition count without losing ordering.

---

## 9. Handling Out-of-Order Messages

When a message arrives that is older than the current entity version, reject or ignore it rather than letting it corrupt state.

### Version check pattern

```go
type OrderProjector struct {
    db     *sqlx.DB
    logger *slog.Logger
}

func (p *OrderProjector) Apply(ctx context.Context, event OrderEvent) error {
    // Fetch current version of the order
    var currentVersion int
    err := p.db.QueryRowContext(ctx,
        "SELECT version FROM orders WHERE id = $1",
        event.OrderID,
    ).Scan(&currentVersion)
    if err != nil && !errors.Is(err, sql.ErrNoRows) {
        return fmt.Errorf("get version: %w", err)
    }

    // Reject stale events
    if event.Version <= currentVersion {
        p.logger.Info("stale event rejected",
            "order_id", event.OrderID,
            "event_version", event.Version,
            "current_version", currentVersion,
        )
        return nil // ack the message — it's genuinely old, not an error
    }

    // Reject events with a gap (missing predecessor)
    if event.Version > currentVersion+1 {
        p.logger.Warn("version gap — event buffered or rejected",
            "order_id", event.OrderID,
            "expected_version", currentVersion+1,
            "received_version", event.Version,
        )
        // Options:
        //   1. Buffer and wait for the missing event (ch 7 approach)
        //   2. Return error — let broker redeliver — hope gap fills
        //   3. Accept and record as a gap for later reconciliation
        return fmt.Errorf("version gap: expected %d, got %d", currentVersion+1, event.Version)
    }

    // Version is exactly currentVersion+1 — apply the event
    return p.applyEvent(ctx, event)
}

func (p *OrderProjector) applyEvent(ctx context.Context, event OrderEvent) error {
    _, err := p.db.ExecContext(ctx, `
        UPDATE orders
        SET status = $1, version = $2, updated_at = NOW()
        WHERE id = $3 AND version = $4`,  // optimistic lock: only update if version still matches
        eventTypeToStatus(event.Type),
        event.Version,
        event.OrderID,
        event.Version-1, // must still be on the previous version
    )
    return err
}
```

### Event store with expected version

Event sourcing systems use an `expected_version` parameter on appends — a strict concurrency control that prevents out-of-order writes:

```go
type EventStore struct{ db *sqlx.DB }

func (s *EventStore) Append(ctx context.Context, streamID string, event any, expectedVersion int) error {
    data, err := json.Marshal(event)
    if err != nil { return err }

    result, err := s.db.ExecContext(ctx, `
        INSERT INTO event_store (stream_id, version, payload, created_at)
        SELECT $1, $2, $3, NOW()
        WHERE NOT EXISTS (
            SELECT 1 FROM event_store
            WHERE stream_id = $1 AND version >= $2
        )`,
        streamID,
        expectedVersion+1,
        data,
    )
    if err != nil { return err }

    rows, _ := result.RowsAffected()
    if rows == 0 {
        return fmt.Errorf("optimistic concurrency conflict: stream %s expected version %d",
            streamID, expectedVersion)
    }
    return nil
}
```

---

## Summary

- **Out-of-order events corrupt state**: a `Cancelled` event processed before `Created` leaves entities in wrong states
- **Partial order is sufficient**: events need to be ordered per entity, not globally
- **Kafka partition key**: `Key: []byte(entityID)` routes all events for one entity to the same partition, where they are consumed in order
- **Key selection**: key on the entity whose events must be ordered — `orderID`, `userID`, `conversationID`, not a random ID or a constant
- **Cardinality**: use high-cardinality keys (entity IDs) to avoid hot spots; estimate `peak_msg_per_sec / target_per_partition * 2` for partition count
- **SQS FIFO**: `MessageGroupId` is the equivalent of Kafka's partition key
- **Consumer-side reordering**: if broker ordering is unavailable, use a sequence number and buffer out-of-order events until gaps fill
- **Version checks**: reject `event.Version <= currentVersion` (stale) and gap-check `event.Version == currentVersion+1` before applying
- **Throughput trade-off**: more partitions = more parallelism, but ordering remains per-partition

---

## Exercises

### Easy
1. Create a Kafka topic `order-events` with 6 partitions. Write a producer that publishes `OrderCreated`, `OrderConfirmed`, and `OrderShipped` events for 10 different orders, using `orderID` as the partition key. Verify with `kafka-go`'s reader that events for each order land on the same partition.
2. Write a consumer that reads from `order-events` and logs the partition, offset, and event type for each message. Confirm that events for the same order always come from the same partition in sequence.
3. Intentionally use `nil` as the message key (round-robin assignment). Observe that events for the same order now arrive on different partitions and are interleaved. Document what ordering properties are lost.

### Medium
4. Implement `OrderProjector.Apply` from section 9 with a real PostgreSQL database. Write a test that sends events version 1, 3, 2 (out of order) and verifies: version 3 returns an error (gap), version 2 succeeds, then version 3 is reprocessed and also succeeds.
5. Build a `SequenceBuffer` (section 7) that buffers out-of-order events and drains them in sequence. Write a test that pushes events in order 1, 3, 5, 2, 4 and verifies they are processed in order 1, 2, 3, 4, 5.
6. Design a partition key strategy for a chat application. The requirements: messages in a conversation must be ordered; users can be in up to 1,000 conversations; peak load is 50,000 messages/second. Calculate the partition count, justify the key choice, and explain how you'd handle a single viral conversation that becomes a hot partition.

### Hard
7. Implement a **consumer group rebalance stress test**: start 6 consumers reading from a 6-partition topic. Each consumer tracks which partition it owns and the last offset processed. Kill one consumer mid-stream. Verify that the surviving consumers pick up the killed consumer's partitions and continue from the last committed offset without gaps or duplicates.
8. Build an **out-of-order reconciliation job**: a background process that scans the `event_store` table for version gaps (missing sequence numbers within a stream) and logs them as alerts. For gaps older than 5 minutes, emit a `MissingEventAlert` with the stream ID and the missing version range. This simulates an on-call alert for permanently lost events.
