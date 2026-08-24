# Chapter 49: Why Async Systems Exist

When a user places an order on Amazon, does Amazon make you wait while it: charges your card, reserves inventory, sends a confirmation email, notifies the warehouse, updates analytics, and generates a shipping label? No — most of that happens after the order is placed, in parallel, asynchronously. This chapter explains why that pattern exists and why every serious system uses it.

## Table of Contents

1. The Synchronous Problem
2. What Async Systems Actually Are
3. The Three Models: Queues, Pub/Sub, Streams
4. When NOT to Use Async
5. CAP Theorem for Async Systems
6. Exercises

---

## 1. The Synchronous Problem

**Scenario: User places an order**

Synchronous approach (everything in one request):

```
POST /orders
├── Validate payment details (100ms)
├── Charge credit card via Stripe API (300ms)
├── Reserve inventory in warehouse system (200ms)
├── Send confirmation email via SendGrid (150ms)
├── Notify analytics system (50ms)
├── Update recommendation engine (400ms)
└── Response: "Order placed" (total: 1200ms)
```

Problems:
1. **Slow:** 1.2 seconds to place an order. Users leave if > 300ms.
2. **Fragile:** If the email service is down, the whole order fails.
3. **Coupled:** Order service needs to know about email service, analytics, recommendations.
4. **Can't retry:** If the notification to the warehouse fails, we charged the card but didn't notify the warehouse.

**Async approach:**

```
POST /orders
├── Validate payment details (100ms)
├── Charge credit card via Stripe API (300ms)
├── Reserve inventory in warehouse system (200ms)
└── Response: "Order placed" (total: 600ms)

Meanwhile, in background:
├── Message broker receives "order_placed" event
├── Email service reads the event → sends email
├── Analytics service reads the event → updates dashboard
└── Recommendation engine reads the event → updates model
```

Benefits:
1. **Fast:** Response in 600ms — only critical work is synchronous.
2. **Resilient:** Email service down? Message queues up, email sent when service recovers.
3. **Decoupled:** Order service doesn't know about email, analytics, or recommendations. It just fires an event.
4. **Retryable:** Each consumer handles its own retries independently.

---

## 2. What Async Systems Actually Are

An async system has three components:

```
Producer                    Broker                  Consumer
(Order Service)         (Kafka/RabbitMQ)        (Email Service)
     │                        │                        │
     │── "order_placed" ──────►│                        │
     │                        │── "order_placed" ─────►│
     │                        │                        │ (processes)
     │                        │                        │
     │   The producer          │   Stores the message   │   The consumer
     │   continues immediately │   durably              │   processes at its
     │   without waiting       │                        │   own pace
```

**Key insight:** The producer and consumer are decoupled in time. The producer doesn't wait for the consumer. They communicate through a shared, durable message store.

Think of it like a mailbox:
- You (producer) put a letter in the mailbox.
- You don't wait for the recipient to read it.
- The letter stays in the mailbox until the recipient picks it up.
- If the recipient is away, the letter waits safely.

---

## 3. The Three Models: Queues, Pub/Sub, Streams

### Model 1: Message Queue (Point-to-Point)

```
Producer → [Queue] → Consumer
```

- Each message is processed by **exactly one consumer**.
- Multiple consumers compete for messages (work distribution).
- Messages are removed after processing.

**Use case:** Background job processing. "Resize this image." "Send this email." "Process this payment."

```
Order Service ──► [email_jobs queue] ──► Email Worker 1
                                     └──► Email Worker 2
                                     └──► Email Worker 3
```

Three workers, each processes different jobs. Total throughput = sum of all workers.

### Model 2: Pub/Sub (One-to-Many)

```
Publisher → [Topic] → Subscriber A
                    → Subscriber B
                    → Subscriber C
```

- Each message is delivered to **every subscriber**.
- Adding a new subscriber doesn't change the publisher.
- Messages may be ephemeral (not stored after delivery).

**Use case:** Event broadcasting. "User logged in" → notify auth service, analytics, fraud detection simultaneously.

### Model 3: Stream (Replay-able Log)

```
Producer → [Stream/Log] ← Consumer A (at its own position)
                        ← Consumer B (at its own position)
                        ← Consumer C (at its own position)
```

- Messages are stored in order (an append-only log).
- Multiple consumers, each tracking their own position in the log.
- **Old messages can be re-read.** Consumer A can replay from the beginning.
- Messages are retained for a configurable time (hours, days, years).

**Use case:** Event sourcing, audit logs, replay for new services. Kafka is the canonical stream.

| | Queue | Pub/Sub | Stream |
|---|---|---|---|
| Consumers | One consumer gets each message | All consumers get each message | All consumers get each message (at their own pace) |
| History | No replay | No replay | Replay from any point |
| Durability | Yes | Sometimes | Yes |
| Examples | RabbitMQ | Redis Pub/Sub | Kafka, NATS JetStream |

---

## 4. When NOT to Use Async

Async is not always better. Use synchronous when:

**1. You need the result immediately**
```
// Synchronous is correct here:
user := db.GetUser(userID) // must have user to render the page
```

**2. You need strong consistency**
```
// If deducting balance is async, two requests could both deduct from the same $100
// This must be synchronous with proper locking:
db.Transaction(func() {
    balance := db.GetBalance(userID)
    if balance < amount { return error }
    db.SetBalance(userID, balance - amount)
})
```

**3. The operation is fast and simple**
```
// Don't add a message broker for this:
db.UpdateLastSeen(userID, time.Now())
```

**The rule:** Use async when work can be deferred, when work should be done by a different service, or when work can fail and retry independently.

---

## 5. CAP Theorem for Async Systems

Async systems are distributed systems. CAP applies:

- **Consistency:** All consumers see messages in the same order.
- **Availability:** The broker is always writable/readable, even if some nodes are down.
- **Partition Tolerance:** The system continues when network splits occur.

| System | CAP |
|--------|-----|
| Kafka | CP (prefers consistency over availability during splits) |
| RabbitMQ | CP by default, AP with specific configs |
| NATS | AP (focuses on availability) |
| Redis Pub/Sub | AP (fire-and-forget, no guarantees) |

**Delivery guarantees:**

- **At-most-once:** Message may be lost, never duplicated. Fire-and-forget. Fastest.
- **At-least-once:** Message is delivered at least once. May duplicate on retry. Most common.
- **Exactly-once:** Message is delivered exactly once. Hardest. Kafka transactions provide this.

Most production systems use **at-least-once** and make consumers **idempotent** (safe to process the same message twice).

```go
// Idempotent consumer: safe to process the same message twice
func processOrderEvent(orderID string) error {
    // Use INSERT OR IGNORE / ON CONFLICT DO NOTHING
    _, err := db.Exec("INSERT OR IGNORE INTO processed_orders (order_id) VALUES (?)", orderID)
    if err != nil {
        return err
    }
    // ... do work ...
}
```

---

## Summary

- Async systems decouple producers and consumers: producers don't wait for consumers to finish.
- Three models: queues (one consumer), pub/sub (all consumers), streams (replayable log per consumer).
- Benefits: faster responses, fault tolerance, service decoupling, independent retries.
- Use synchronous for: operations that require immediate results, financial consistency, simple fast operations.
- At-least-once delivery + idempotent consumers is the production standard.

### Exercises

**Easy:** For each of these operations, decide: sync or async? Justify your answer. (1) Charging a credit card. (2) Sending a welcome email. (3) Updating a user's profile picture. (4) Logging a page view for analytics.

**Medium:** Draw the architecture for a food delivery app where a customer places an order. Identify which operations should be synchronous (in the request) and which should be async (via events). Name the events and which services consume them.

**Hard:** Research the "two generals problem" and explain why exactly-once delivery is so hard. How does Kafka's transactional API solve it (or partially solve it)?
