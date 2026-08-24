# Chapter 95: Async Systems — Why, When, and Trade-offs

Every service you build does some work that takes time: sending an email, charging a credit card, resizing an uploaded image, generating a report. The question is: does the user have to wait for it? If yes, you have a synchronous system. If no, you have an asynchronous one. This chapter is about understanding when async is the right choice, what guarantees it actually provides, and which tool to reach for.

## Table of Contents

1. [Sync vs Async — The Fundamental Difference](#1-sync-vs-async--the-fundamental-difference)
2. [When to Choose Async](#2-when-to-choose-async)
3. [When NOT to Use Async](#3-when-not-to-use-async)
4. [Async Patterns](#4-async-patterns)
5. [Delivery Guarantees](#5-delivery-guarantees)
6. [The Trade-off Table](#6-the-trade-off-table)
7. [Tools in This Volume](#7-tools-in-this-volume)
8. [How to Choose the Right Tool](#8-how-to-choose-the-right-tool)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Sync vs Async — The Fundamental Difference

### The analogy

Picture a restaurant. A customer orders food. What does the waiter do next?

```
Synchronous waiter                  Asynchronous waiter
──────────────────                  ───────────────────
Customer: "I'll have a steak."      Customer: "I'll have a steak."
Waiter walks to kitchen.            Waiter writes the order on a ticket,
Waiter stands there.                clips it to the kitchen rail,
Waiter waits for the steak.         then takes orders from 5 other tables.
Steak arrives.
Waiter walks back.                  Kitchen finishes the steak.
Waiter serves the steak.            Kitchen rings a bell.
Waiter takes next order.            Waiter delivers the steak.
```

The synchronous waiter serves one table at a time. The asynchronous waiter can handle many tables concurrently because they decouple "taking the order" from "delivering the food."

### In code

```go
// Synchronous: the HTTP handler does everything before responding
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    user := createUser(r)
    db.Save(user)
    sendWelcomeEmail(user)      // ← 300ms email service call
    sendSlackNotification(user) // ← 150ms slack call
    generateAvatarThumbnail(user) // ← 800ms image processing
    // Total response time: 1250ms for something the user doesn't see
    w.WriteHeader(http.StatusCreated)
}

// Asynchronous: handler does only what the user needs, queues the rest
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    user := createUser(r)
    db.Save(user)
    h.queue.Publish("user.registered", user.ID) // ← ~1ms
    // Total response time: ~20ms
    w.WriteHeader(http.StatusCreated)
    // Email, Slack, thumbnail happen in the background
}
```

The user gets a response in 20ms instead of 1250ms. The background work still happens — just not in the hot path.

### What changes under the hood

```
Synchronous request lifecycle:
  HTTP Request → Handler → [work A] → [work B] → [work C] → HTTP Response
                             ↑                        ↑
                          blocking                 blocking

Asynchronous request lifecycle:
  HTTP Request → Handler → [save to DB] → publish event → HTTP Response
                                                   ↓
                                            Message Queue
                                                   ↓
                                  Worker: [work A] [work B] [work C]
                                  (independently, in background)
```

---

## 2. When to Choose Async

### Tasks that take time and the user does not need to see the result immediately

```go
// These are all good candidates for async processing:
// - Sending email or SMS
// - Charging a payment (you acknowledge "we received your request", then confirm later)
// - Resizing uploaded images to multiple sizes
// - Generating a PDF report
// - Syncing data to a third-party CRM
// - Sending a push notification
// - Indexing a new document for search

type UserRegisteredEvent struct {
    UserID    int64     `json:"user_id"`
    Email     string    `json:"email"`
    CreatedAt time.Time `json:"created_at"`
}

// Publisher: the HTTP handler
func (h *Handler) Register(w http.ResponseWriter, r *http.Request) {
    user := h.createUser(r)
    event := UserRegisteredEvent{
        UserID:    user.ID,
        Email:     user.Email,
        CreatedAt: time.Now(),
    }
    h.queue.Publish(ctx, "user.registered", event)
    json.NewEncoder(w).Encode(user)
}

// Subscribers: independent workers, each doing one thing
// Worker 1: sends welcome email
// Worker 2: creates Stripe customer
// Worker 3: sends Slack notification to #new-signups
// Worker 4: generates avatar thumbnail
```

### Decoupling services

Async messaging decouples services in time and space. The publisher does not need to know which services consume its events, and it does not need them to be alive.

```
Synchronous coupling:
  Order Service ──HTTP──> Email Service   (Email Service down? Order fails)
  Order Service ──HTTP──> Billing Service (Billing down? Order fails)
  Order Service ──HTTP──> Warehouse       (Warehouse down? Order fails)

Asynchronous decoupling:
  Order Service ──publishes──> "order.placed" event
                                    ↓
                     ┌─────────────┼──────────────┐
                     ↓             ↓               ↓
               Email Worker  Billing Worker  Warehouse Worker
               (retries on failure independently)
```

If the Email Worker is down, the Billing Worker still processes the payment. When the Email Worker recovers, it processes the backlog. The Order Service never knew anything was wrong.

### Handling traffic spikes

A queue acts as a buffer. Your background workers can process at a sustainable rate even when the front-end receives 10x normal traffic.

```
Normal traffic:   100 req/s  → queue depth: ~0     → workers process in real-time
Traffic spike:   1000 req/s  → queue depth: grows  → workers drain at 100/s
After spike:      100 req/s  → queue depth: shrinks → back to normal in minutes
```

Without a queue, the spike hits your database and downstream services directly. With a queue, they only ever see 100 req/s.

### Fan-out to multiple consumers

One event, many reactions. This is where async shines most clearly.

```go
// One event published once
event := OrderPlacedEvent{OrderID: 42, CustomerID: 7, Total: 149.99}
queue.Publish("order.placed", event)

// Multiple independent consumers, each with their own offset/cursor:
// Consumer Group A (billing):    charges the card
// Consumer Group B (warehouse):  reserves the items
// Consumer Group C (email):      sends confirmation
// Consumer Group D (analytics):  increments daily revenue counter
// Consumer Group E (fraud):      scores the transaction
```

In a synchronous system you would have to call all five services in sequence. Each adds latency and each can fail the whole operation.

---

## 3. When NOT to Use Async

Async adds real complexity. Do not reach for it when you do not need it.

### When the user needs the result immediately

```
Synchronous: User submits a form.
             Server validates and saves.
             Server responds: "saved" or "email already taken".
             User sees the result in real-time.

Async gone wrong: User submits a form.
                  Server responds: "processing..."
                  Validation happens in a worker 200ms later.
                  User gets an email saying "email already taken".
                  User is confused.
```

Validation, authentication, and anything that directly determines the response content must be synchronous.

### When strict ordering matters and you cannot partition

Some workflows are inherently sequential:

```
Step 1: Charge the card
Step 2: Create the order (only if step 1 succeeded)
Step 3: Send confirmation (only if step 2 succeeded)
```

If you split these into separate async messages, you need to handle step 2 failing after step 1 succeeded (partial execution). This is solvable — with sagas, outbox patterns, or idempotency keys — but it is significantly more complex than a simple sequential function.

### When debugging complexity is too high for the value

Async systems are harder to debug:
- A bug may not surface until the worker processes the message, which could be minutes later.
- The stack trace in the worker has no trace of the original HTTP request.
- Distributed tracing (with `trace-id` propagated through the message) is required to connect the dots.
- Failed messages can silently rot in a dead-letter queue.

If the task is simple, fast, and rarely fails — just do it synchronously.

### The rule of thumb

```
Ask: "If this background task fails permanently, what happens?"

  → "The user never gets their confirmation email"
    → Async is fine. Email is recoverable, not critical-path.

  → "The user paid but never got the product"
    → You need strong guarantees. Consider sync + saga, or
      at-least-once delivery with idempotency.

  → "A counter is slightly off for a few minutes"
    → Async is fine. Eventual consistency is acceptable here.
```

---

## 4. Async Patterns

### Fire-and-Forget

Publish an event and move on. No response is expected. The simplest pattern — and the least reliable.

```go
func (h *Handler) UpdateProfile(w http.ResponseWriter, r *http.Request) {
    user := updateUser(r)
    db.Save(user)

    // Fire-and-forget: invalidate CDN cache, sync to CRM, etc.
    // If this fails, it fails silently. Use only for non-critical side effects.
    go func() {
        h.cdn.InvalidateCache(user.ID)
        h.crm.SyncUser(user)
    }()

    json.NewEncoder(w).Encode(user)
}
```

### Task Queues

A job is enqueued and a worker picks it up. The producer gets an acknowledgement that the job was queued (not that it ran). Workers can retry on failure.

```
Producer                Queue              Worker
────────                ─────              ──────
Enqueue(SendEmail{...}) → [job1, job2, ...]  ← Dequeue
                                             Process(job1)
                                             Ack  ← message removed from queue
```

### Event Streaming

An ordered log of events that consumers read at their own pace. Unlike a queue, events are not deleted after consumption — each consumer group tracks its own position.

```
Event Log (Kafka / Redis Streams)
  offset: 1  2  3  4  5  6  7  8 ...
  events: A  B  C  D  E  F  G  H

Consumer Group "billing":   at offset 7, processing 8
Consumer Group "email":     at offset 4, processing 5 (fell behind)
Consumer Group "analytics": at offset 8, keeping up
```

### Pub/Sub

Publish to a topic; all current subscribers receive it. Messages are not persisted — if a subscriber is offline, it misses the message.

```
Publisher: topic="user.updated", payload={id: 1, name: "Alice"}

Subscriber A (cache invalidation): receives immediately
Subscriber B (search index):       receives immediately
Subscriber C (was offline):        misses the message entirely
```

Pub/sub is appropriate for real-time notifications where missing a message is acceptable (e.g., live dashboard updates, chat typing indicators).

---

## 5. Delivery Guarantees

This is the most important thing to understand about async systems. Not all queues give the same promises.

### At-Most-Once (Best-Effort)

The message is sent once. If the consumer crashes before processing it, the message is lost. No retries.

```
Producer → Queue → Consumer crashes → message gone forever

Use for: metrics, analytics events, non-critical notifications
         where a small % of loss is acceptable and retries would
         cause more harm than good (e.g., spam prevention)
```

### At-Least-Once (The Practical Default)

The message is delivered at least once. If the consumer crashes after processing but before acknowledging, the message is redelivered. The consumer may process the same message more than once.

```
Producer → Queue → Consumer processes → crashes before ACK
                 → Queue redelivers    → Consumer processes again

This is what Redis Streams, Kafka, RabbitMQ, and asynq all give you
by default (with the right configuration).
```

The solution is **idempotency**: make processing the same message twice have the same effect as processing it once.

```go
// Idempotent job handler: safe to run multiple times
func ProcessOrder(ctx context.Context, orderID int64) error {
    // Check if already processed
    var processed bool
    db.QueryRowContext(ctx,
        "SELECT processed FROM orders WHERE id = $1", orderID,
    ).Scan(&processed)
    if processed {
        return nil // already done, safe to return success
    }

    // Do the work
    if err := chargeCard(ctx, orderID); err != nil {
        return err
    }

    // Mark as processed
    _, err := db.ExecContext(ctx,
        "UPDATE orders SET processed = TRUE WHERE id = $1 AND processed = FALSE",
        orderID,
    )
    return err
}
```

### Exactly-Once (Extremely Hard)

The message is processed exactly once, even across crashes. Truly exactly-once requires distributed transactions between the queue and the database, which most systems do not support.

```
"Exactly-once" in practice means:
  at-least-once delivery + idempotent consumer = effectively-once processing

Kafka Transactions: producer can write to multiple partitions atomically.
  Still doesn't prevent the consumer from processing twice if it crashes
  after processing but before committing the offset.

The real solution: don't try to achieve exactly-once in the infrastructure.
  Make your consumer idempotent instead.
```

```go
// Idempotency key pattern: deduplicate by a unique key per event
func HandlePayment(ctx context.Context, event PaymentEvent) error {
    // Use idempotency key to prevent double-charging
    _, err := db.ExecContext(ctx, `
        INSERT INTO payment_events (idempotency_key, order_id, amount, status)
        VALUES ($1, $2, $3, 'pending')
        ON CONFLICT (idempotency_key) DO NOTHING
    `, event.IdempotencyKey, event.OrderID, event.Amount)

    // If ON CONFLICT fired, this event was already processed
    // Check if the INSERT actually happened
    var exists bool
    db.QueryRowContext(ctx,
        "SELECT EXISTS(SELECT 1 FROM payment_events WHERE idempotency_key = $1 AND status = 'done')",
        event.IdempotencyKey,
    ).Scan(&exists)
    if exists {
        return nil // already charged, skip
    }

    return chargeCard(ctx, event)
}
```

---

## 6. The Trade-off Table

No tool or pattern is free. Here is what you are trading:

```
Dimension          Sync                    Async
─────────────────────────────────────────────────────────────────────────
Latency            Low (immediate result)  Higher (enqueue + process)
Throughput         Limited by slowest step Each stage scales independently
Complexity         Low                     High: queues, workers, retries,
                                           dead-letter queues, monitoring
Observability      Easy: one stack trace   Hard: need distributed tracing,
                                           correlate request → event → worker
Error handling     Immediate, in-process   Delayed, in worker logs
Ordering           Natural (sequential)    Requires explicit partitioning
Consistency        Strong (same process)   Eventual (may be seconds behind)
Failure mode       User sees error         Silent failure in background
Testing            Straightforward         Requires test harness for queues
Deployment         Single service          Workers + queue infra to run
```

### The practical rule

Start synchronous. Move tasks to async when you have a specific problem:
- Response time is too high because of slow side effects.
- One slow dependency is blocking unrelated work.
- Traffic spikes are overwhelming a downstream service.
- You need the same event processed by multiple independent consumers.

Do not add async complexity speculatively.

---

## 7. Tools in This Volume

This volume covers async tools in order from simple to powerful:

```
Chapter 92: Worker Pools
  What:   In-process goroutine pools, no external infrastructure
  Use:    CPU-bound work, bounded concurrency within a service
  Infra:  None

Chapter 94: asynq (Redis-backed task queue)
  What:   Persistent job queue backed by Redis
  Use:    Background jobs with retries, scheduling, priorities
  Infra:  Redis

Chapter 93: Kafka
  What:   Distributed event streaming log
  Use:    High throughput, multiple consumers, event replay
  Infra:  Kafka cluster (Zookeeper or KRaft)

Redis Streams (Chapter 92 extended)
  What:   Lightweight persistent streams in Redis
  Use:    Simpler alternative to Kafka for moderate scale
  Infra:  Redis

RabbitMQ / NATS (covered in microservices chapters)
  What:   Traditional message broker / lightweight pub-sub
  Use:    Complex routing, fan-out, request-reply
  Infra:  RabbitMQ or NATS server

Watermill (abstraction layer)
  What:   Go library that wraps Kafka, RabbitMQ, NATS, etc.
  Use:    Event-driven architectures without locking into one broker
  Infra:  Depends on backend
```

---

## 8. How to Choose the Right Tool

```
                    START HERE
                        │
                        ▼
            Do you need persistence?
            (survive restarts, replay)
               │               │
              YES               NO
               │               │
               ▼               ▼
    Do you need        Use goroutine
    multiple           channels or
    consumers?         in-process
    (fan-out)          worker pool
       │       │       (Chapter 92)
      YES       NO
       │        │
       ▼        ▼
   How many    Do you need
   msg/sec?    scheduling /
               priorities?
               │          │
              YES          NO
               │           │
               ▼           ▼
             asynq       asynq or
             (Chapter    postgres
               94)        queue
                          (Chapter 78
                           FOR UPDATE
                           SKIP LOCKED)
       │
  < 100k/s        > 100k/s
       │               │
       ▼               ▼
  Redis Streams    Kafka
  or RabbitMQ    (Chapter 93)
```

### Decision in plain words

| Situation | Recommendation |
|-----------|----------------|
| Background jobs with retries (emails, notifications) | asynq (Redis) |
| Job queue with no extra infra | PostgreSQL + `FOR UPDATE SKIP LOCKED` |
| High throughput event log, multiple consumers, replay | Kafka |
| Simple pub/sub, no persistence needed | NATS or Redis Pub/Sub |
| CPU-bound tasks, same process | Worker pool (goroutines) |
| Complex routing rules between services | RabbitMQ |
| You want to swap brokers later | Watermill |

---

## Summary

- **Sync**: the caller waits; simple, low-latency for fast operations, breaks under slow side effects
- **Async**: decouple fast acceptance from slow processing; the waiter takes more orders while the kitchen cooks
- **Use async when**: side effects take time, services need decoupling, traffic spikes need buffering, fan-out to multiple consumers is required
- **Avoid async when**: the user needs an immediate result, strict ordering is required without partitioning, the debugging overhead is not worth the gain
- **Patterns**: fire-and-forget, task queues, event streaming, pub/sub
- **Delivery guarantees**: at-most-once (lossy), at-least-once (most practical), exactly-once (use idempotency instead)
- **Trade-offs**: async buys throughput and resilience at the cost of latency, complexity, and observability
- **Tool selection**: PostgreSQL queue → asynq → Redis Streams → Kafka, scaled to your throughput and fan-out needs

## Exercises

### Easy
1. Draw (in ASCII or on paper) the lifecycle of a user registration request in two versions: synchronous (everything in the HTTP handler) and asynchronous (handler saves user, queues an event, workers handle email/Slack/thumbnail). Label where failures can occur in each.
2. List five operations from a real-world app you know (a bank, an e-commerce site, a social network). For each, decide: sync or async? Write one sentence justifying each decision.
3. Write a `FireAndForget(fn func())` helper in Go that runs `fn` in a goroutine, recovers from panics, and logs them. Use it to trigger a cache invalidation after saving a user. Why is this appropriate for cache invalidation but not for billing?

### Medium
4. Implement a **simple in-process task queue** using a buffered channel and a worker pool. The queue should support: `Enqueue(task Task) error` (returns error if queue is full), `Start(workers int)`, and `Shutdown(ctx context.Context) error` (drains pending tasks then stops). Write a test that enqueues 1000 tasks and verifies all complete.
5. Build an **idempotent event handler** for a `PaymentProcessed` event. The handler inserts into a `processed_events` table using `ON CONFLICT (event_id) DO NOTHING`, then charges a card only if the insert succeeded. Write a test that delivers the same event 5 times and verifies the card is charged exactly once.
6. Implement a **dead-letter queue** simulation: a `Queue` struct that retries failed jobs up to 3 times with exponential backoff (100ms, 200ms, 400ms), then moves them to a separate `deadLetterQueue` channel. Write a test that submits 10 jobs where 3 always fail, and verifies those 3 end up in the dead-letter queue after all retries are exhausted.

### Hard
7. Design a **multi-consumer event bus** in Go: one `Publish(topic string, payload any)` method fans out to all registered subscribers for that topic. Each subscriber has its own buffered channel and independent goroutine. If a subscriber's channel is full, the message is dropped for that subscriber (backpressure). Write a benchmark comparing throughput at 1, 5, and 20 subscribers.
8. Build a **saga coordinator** for a three-step checkout flow: (1) reserve inventory, (2) charge card, (3) create shipment. Each step is a separate async task. If step 2 fails, compensate step 1 (release inventory). If step 3 fails, compensate steps 1 and 2 (release inventory + refund card). Implement using a state machine stored in PostgreSQL so the saga survives service restarts.
