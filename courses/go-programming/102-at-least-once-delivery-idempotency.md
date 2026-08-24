# Chapter 102: At-Least-Once Delivery and Idempotency

In distributed systems, messages can be delivered more than once. Network failures, broker restarts, consumer crashes, and offset commit failures all create situations where a message is re-delivered. The standard contract for production messaging systems is **at-least-once delivery** — every message is guaranteed to be processed, but it may be processed more than once. Your code must handle duplicates.

## Table of Contents

1. [Why At-Least-Once?](#1-why-at-least-once)
2. [The Duplicate Problem](#2-the-duplicate-problem)
3. [Idempotency Keys](#3-idempotency-keys)
4. [Database Deduplication](#4-database-deduplication)
5. [Redis Deduplication](#5-redis-deduplication)
6. [Naturally Idempotent Operations](#6-naturally-idempotent-operations)
7. [Designing Idempotent Handlers](#7-designing-idempotent-handlers)
8. [The Consumer Pattern](#8-the-consumer-pattern)
9. [Exactly-Once: Why It's a Trap](#9-exactly-once-why-its-a-trap)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why At-Least-Once?

Message delivery has three theoretical guarantees:

```
At-most-once:   message may be lost, never duplicated
                (fire and forget — fast but unreliable)

At-least-once:  message will be delivered, may be duplicated
                (the practical standard for production systems)

Exactly-once:   message delivered exactly once, no duplicates
                (requires distributed coordination — very expensive)
```

At-least-once happens naturally whenever a system acknowledges delivery only after successful processing:

```
Consumer crashes after processing but BEFORE committing offset:
  Broker:   [msg1✓] [msg2✓] [msg3?] [msg4] ...
                                 ↑
                       Consumer processed this, then crashed.
                       No commit → broker thinks msg3 undelivered.

On restart:
  Broker re-delivers msg3.
  Consumer processes msg3 again.
  ← This is the duplicate.
```

The scenarios that trigger redelivery:

- Consumer process crashes mid-flight
- Network timeout between consumer and broker
- Broker leader election (brief unavailability)
- Manual offset reset for replay
- Consumer rebalance during processing

---

## 2. The Duplicate Problem

Duplicates are harmless for reads but dangerous for writes. The classic example is a payment:

```go
// PaymentHandler processes charge events
func (h *PaymentHandler) HandleChargeEvent(ctx context.Context, msg kafka.Message) error {
    var event ChargeEvent
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        return fmt.Errorf("unmarshal: %w", err)
    }

    // DANGER: if this message is delivered twice, the customer is charged twice
    if err := h.stripe.Charge(ctx, event.CustomerID, event.AmountCents); err != nil {
        return fmt.Errorf("charge: %w", err)
    }

    // If we crash HERE, before committing the offset:
    //   - payment was charged ✓
    //   - offset not committed ✗
    //   - on restart: broker re-delivers → charged AGAIN

    return nil // commit offset
}
```

What duplicates look like in practice:

```
Timeline:
  t=0  Consumer receives charge_event(id=abc123, amount=5000)
  t=1  Stripe.Charge() succeeds — customer charged $50
  t=2  Consumer crashes (power cut, OOM kill, deploy)
  t=3  Consumer restarts
  t=4  Broker re-delivers charge_event(id=abc123, amount=5000)
  t=5  Stripe.Charge() succeeds — customer charged $50 AGAIN

Result: customer sees two $50 charges for one order.
```

The fix is not to prevent redelivery — it's to make your handler safe when it happens.

---

## 3. Idempotency Keys

An **idempotency key** is a unique identifier for an operation. The server stores it after the first successful execution and returns the same result for any repeat with the same key — without re-executing the side effect.

### HTTP-level idempotency

```go
// Client sends an Idempotency-Key header with each request
func (c *PaymentClient) ChargeCard(ctx context.Context, req ChargeRequest) (*ChargeResponse, error) {
    // Generate a stable key from the operation — not random on retry
    idempotencyKey := fmt.Sprintf("charge:%s:%d", req.OrderID, req.AmountCents)

    httpReq, _ := http.NewRequestWithContext(ctx, "POST", "/v1/charges", encodeJSON(req))
    httpReq.Header.Set("Idempotency-Key", idempotencyKey)
    httpReq.Header.Set("Content-Type", "application/json")

    resp, err := c.http.Do(httpReq)
    if err != nil { return nil, err }
    defer resp.Body.Close()

    var result ChargeResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result, nil
}

// Server side: store and replay the response
func (h *ChargeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    idempKey := r.Header.Get("Idempotency-Key")

    if idempKey != "" {
        if cached, ok := h.store.Get(idempKey); ok {
            // Return the previously computed response without re-charging
            w.Header().Set("X-Idempotency-Replayed", "true")
            writeJSON(w, cached)
            return
        }
    }

    result, err := h.stripe.Charge(r.Context(), parseRequest(r))
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }

    if idempKey != "" {
        // Store result for 24 hours — future retries return this
        h.store.Set(idempKey, result, 24*time.Hour)
    }

    writeJSON(w, result)
}
```

Stripe, PayPal, and most payment APIs implement this exact pattern. The key must be:
- **Deterministic**: generated from the operation, not random — so a retry uses the same key
- **Scoped**: include user ID and operation type to prevent cross-user collisions
- **Bounded TTL**: idempotency windows are typically 24 hours

---

## 4. Database Deduplication

For event consumers, the cleanest deduplication strategy is a `processed_events` table. The `ON CONFLICT DO NOTHING` pattern is atomic and works correctly under concurrent consumers.

```sql
-- One-time setup
CREATE TABLE processed_events (
    id           TEXT PRIMARY KEY,
    processed_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
```

```go
// processIfNew inserts the event ID; returns (true, nil) if new, (false, nil) if duplicate.
func processIfNew(ctx context.Context, db *sqlx.DB, eventID string) (bool, error) {
    result, err := db.ExecContext(ctx, `
        INSERT INTO processed_events (id, processed_at)
        VALUES ($1, NOW())
        ON CONFLICT (id) DO NOTHING`,
        eventID,
    )
    if err != nil {
        return false, fmt.Errorf("dedup insert: %w", err)
    }

    rowsAffected, err := result.RowsAffected()
    if err != nil {
        return false, err
    }

    // 1 row inserted = new event. 0 rows = conflict = already processed.
    return rowsAffected == 1, nil
}

// Example handler using the deduplication function
type PaymentEventHandler struct {
    db     *sqlx.DB
    stripe StripeClient
    logger *slog.Logger
}

func (h *PaymentEventHandler) Handle(ctx context.Context, event ChargeEvent) error {
    isNew, err := processIfNew(ctx, h.db, event.ID)
    if err != nil {
        return fmt.Errorf("dedup check: %w", err)
    }
    if !isNew {
        h.logger.Info("duplicate event, skipping", "event_id", event.ID)
        return nil // safe to ack
    }

    // Only reaches here on the first delivery
    if err := h.stripe.Charge(ctx, event.CustomerID, event.AmountCents); err != nil {
        // The dedup row was inserted but charge failed.
        // Delete the dedup row so a retry can re-attempt.
        h.db.ExecContext(ctx, "DELETE FROM processed_events WHERE id = $1", event.ID)
        return fmt.Errorf("charge: %w", err)
    }

    return nil
}
```

### Combining dedup + business logic in one transaction

For maximum safety, wrap the dedup insert and the business write in a single transaction:

```go
func (h *PaymentEventHandler) HandleTransactional(ctx context.Context, event ChargeEvent) error {
    return h.db.WithTx(ctx, func(tx *sqlx.Tx) error {
        // Dedup insert — if this conflicts, roll back and skip
        result, err := tx.ExecContext(ctx, `
            INSERT INTO processed_events (id) VALUES ($1)
            ON CONFLICT (id) DO NOTHING`,
            event.ID,
        )
        if err != nil { return err }

        rows, _ := result.RowsAffected()
        if rows == 0 {
            return nil // duplicate — rollback is a no-op
        }

        // Business write in the same transaction
        _, err = tx.ExecContext(ctx, `
            INSERT INTO payment_ledger (order_id, customer_id, amount_cents, created_at)
            VALUES ($1, $2, $3, NOW())`,
            event.OrderID, event.CustomerID, event.AmountCents,
        )
        return err
        // Commit makes both rows appear atomically.
        // If the process crashes and retries, ON CONFLICT catches it.
    })
}
```

---

## 5. Redis Deduplication

When your consumer doesn't use a relational database, Redis provides fast deduplication via `SETNX` (Set if Not Exists). Use a TTL to bound memory growth.

```go
import "github.com/redis/go-redis/v9"

type RedisDeduplicator struct {
    rdb *redis.Client
    ttl time.Duration
}

// MarkSeen returns true if the ID is new (first time seen), false if duplicate.
func (d *RedisDeduplicator) MarkSeen(ctx context.Context, id string) (bool, error) {
    key := fmt.Sprintf("event:%s", id)

    // SET key 1 EX ttl NX
    // NX = only set if the key does not exist
    // Returns true if the key was set (new), false if it already existed (duplicate)
    set, err := d.rdb.SetNX(ctx, key, 1, d.ttl).Result()
    if err != nil {
        return false, fmt.Errorf("redis setnx: %w", err)
    }
    return set, nil
}

// Example usage in a consumer
type NotificationHandler struct {
    dedup  *RedisDeduplicator
    mailer Mailer
    logger *slog.Logger
}

func (h *NotificationHandler) Handle(ctx context.Context, event OrderPlacedEvent) error {
    isNew, err := h.dedup.MarkSeen(ctx, event.ID)
    if err != nil {
        // Redis is down — decide: fail safe (block) or fail open (allow through)
        // For notifications, fail open is usually acceptable
        h.logger.Warn("dedup unavailable, processing anyway", "event_id", event.ID, "err", err)
        isNew = true
    }

    if !isNew {
        h.logger.Info("duplicate notification skipped", "event_id", event.ID)
        return nil
    }

    return h.mailer.SendOrderConfirmation(ctx, event.CustomerID, event.OrderID)
}
```

### Redis vs PostgreSQL deduplication

```
                Redis SETNX          PostgreSQL ON CONFLICT
Latency         ~0.5ms               ~2-5ms
Durability      Optional (AOF/RDB)   Always (WAL)
TTL support     Built-in             Manual cleanup job
Transactions    No (standalone)      Yes — wrap with business write
Memory          Proportional to TTL  Proportional to event volume
Best for        High-throughput,     Payments, financial writes,
                ephemeral events     must-not-duplicate operations
```

---

## 6. Naturally Idempotent Operations

Some operations are safe to repeat by definition. Prefer these when designing your APIs:

```go
// PUT: overwrite is safe — calling twice yields the same state
//   PUT /users/123/preferences  { "theme": "dark" }
//   Running twice: user still has theme=dark. No harm.

// DELETE: deleting something already deleted is fine
//   DELETE /sessions/abc
//   First call: session removed. Second call: 404 or no-op. State is the same.

// UPSERT: insert or update — idempotent by construction
_, err = db.ExecContext(ctx, `
    INSERT INTO user_preferences (user_id, theme, updated_at)
    VALUES ($1, $2, NOW())
    ON CONFLICT (user_id) DO UPDATE
    SET theme = EXCLUDED.theme, updated_at = EXCLUDED.updated_at`,
    userID, theme,
)

// Conditional UPDATE: only apply if state hasn't advanced
_, err = db.ExecContext(ctx, `
    UPDATE orders
    SET status = 'confirmed', confirmed_at = NOW()
    WHERE id = $1 AND status = 'pending'`,  // guard: only if still pending
    orderID,
)
// Running this twice on a 'confirmed' order: no rows updated. Safe.

// NOT idempotent — avoid in event handlers:
//   INSERT without ON CONFLICT         (will fail on duplicate)
//   Increment a counter unconditionally (counter grows with each retry)
//   Append to a log without dedup      (creates duplicate entries)
```

---

## 7. Designing Idempotent Handlers

The principle: **make the side effect conditional on current state, not on receiving the message**.

```
Fragile (side effect is unconditional):
  receive event → charge customer → ack

Idempotent (side effect is conditional on state):
  receive event → check if already charged → if not, charge → ack
```

```go
// Fragile: triggers the side effect every time
func (h *Handler) FulfillOrder(ctx context.Context, event OrderPaidEvent) error {
    return h.warehouse.CreateShipment(ctx, event.OrderID) // runs even on re-delivery
}

// Idempotent: checks current state first
func (h *Handler) FulfillOrderIdempotent(ctx context.Context, event OrderPaidEvent) error {
    order, err := h.orders.Get(ctx, event.OrderID)
    if err != nil {
        return fmt.Errorf("get order: %w", err)
    }

    // Guard: only create shipment if not already fulfilled
    if order.Status == "fulfilled" || order.Status == "shipped" {
        h.logger.Info("order already fulfilled, skipping", "order_id", event.OrderID)
        return nil
    }

    if err := h.warehouse.CreateShipment(ctx, event.OrderID); err != nil {
        return fmt.Errorf("create shipment: %w", err)
    }

    return h.orders.UpdateStatus(ctx, event.OrderID, "fulfilled")
}
```

State machine guards are a powerful idempotency tool. If the transition is only valid from certain states, duplicate events that arrive after the state has advanced are silently ignored:

```go
// Valid transitions: pending → confirmed → fulfilled → shipped
var validTransitions = map[string]string{
    "pending":   "confirmed",
    "confirmed": "fulfilled",
    "fulfilled": "shipped",
}

func (h *Handler) TransitionOrder(ctx context.Context, orderID, newStatus string) error {
    order, err := h.orders.Get(ctx, orderID)
    if err != nil { return err }

    expectedCurrent := reverseTransition(newStatus)
    if order.Status != expectedCurrent {
        // Duplicate or out-of-order event — already past this state
        h.logger.Info("transition not applicable",
            "order_id", orderID,
            "current", order.Status,
            "requested", newStatus,
        )
        return nil // not an error — just idempotent no-op
    }

    return h.orders.UpdateStatus(ctx, orderID, newStatus)
}
```

---

## 8. The Consumer Pattern

Putting it all together: the safe consumer loop that handles at-least-once delivery correctly.

```go
type IdempotentConsumer struct {
    reader *kafka.Reader
    dedup  Deduplicator
    handler EventHandler
    logger *slog.Logger
}

// Deduplicator abstracts over Redis or PostgreSQL backends
type Deduplicator interface {
    MarkSeen(ctx context.Context, id string) (isNew bool, err error)
}

func (c *IdempotentConsumer) Run(ctx context.Context) error {
    for {
        msg, err := c.reader.FetchMessage(ctx)
        if err != nil {
            if ctx.Err() != nil { return nil } // graceful shutdown
            return fmt.Errorf("fetch: %w", err)
        }

        if err := c.processMessage(ctx, msg); err != nil {
            c.logger.Error("processing failed — will retry on re-delivery",
                "offset", msg.Offset,
                "partition", msg.Partition,
                "err", err,
            )
            // Do NOT commit the offset — let the broker redeliver
            continue
        }

        // Only commit after confirmed success
        if err := c.reader.CommitMessages(ctx, msg); err != nil {
            c.logger.Warn("commit failed", "err", err)
        }
    }
}

func (c *IdempotentConsumer) processMessage(ctx context.Context, msg kafka.Message) error {
    var event DomainEvent
    if err := json.Unmarshal(msg.Value, &event); err != nil {
        // Malformed message — ack it to avoid poison-pill loop
        c.logger.Error("unmarshal failed — skipping malformed message", "err", err)
        c.reader.CommitMessages(ctx, msg)
        return nil
    }

    // Step 1: check for duplicate
    isNew, err := c.dedup.MarkSeen(ctx, event.ID)
    if err != nil {
        return fmt.Errorf("dedup: %w", err) // treat as transient — retry
    }
    if !isNew {
        c.logger.Info("duplicate event, acking without processing", "event_id", event.ID)
        return nil // return nil so the caller commits the offset
    }

    // Step 2: process (first delivery only)
    if err := c.handler.Handle(ctx, event); err != nil {
        // Undo the dedup mark so a retry can re-attempt
        c.dedup.Delete(ctx, event.ID)
        return fmt.Errorf("handle: %w", err)
    }

    // Step 3: offset committed by the caller after this returns nil
    return nil
}
```

The flow visualised:

```
Message arrives
      │
      ▼
 dedup.MarkSeen(id)
      │
   ┌──┴──────────────────┐
already seen?           new?
      │                   │
   ack & skip         handler.Handle()
                           │
                      ┌────┴────┐
                   error      success
                      │           │
               delete dedup     commit offset
               return error
               (broker retries)
```

---

## 9. Exactly-Once: Why It's a Trap

Exactly-once delivery sounds ideal but requires cross-system distributed coordination:

```
To guarantee exactly-once across Kafka + PostgreSQL:
  - Read message from Kafka (within a Kafka transaction)
  - Write to PostgreSQL (within a DB transaction)
  - Commit the Kafka offset
  - Commit the DB transaction
  ← These two commits must be atomic. They're different systems.
    Making them atomic requires 2-Phase Commit (2PC).
```

The problems with 2PC:

```
2PC requires a coordinator. If the coordinator crashes after
"prepare" but before "commit", all participants block indefinitely.

         Coordinator
              │
     ┌────────┴────────┐
     │                 │
  Kafka DB           Postgres
  "prepared"         "prepared"
              │
        coordinator CRASHES
              │
  Both systems hold locks.
  Neither knows to commit or abort.
  System is stuck until coordinator recovers.
```

Kafka's exactly-once (transactions API) works within Kafka itself — produce + consume atomically. But the moment you write to Postgres, you cross a system boundary, and the atomic guarantee breaks.

**The practical answer**: at-least-once + idempotency achieves the same business outcome:

```
Exactly-once delivery:
  Complex, expensive, breaks across system boundaries

At-least-once + idempotent handler:
  Simple to implement, composes across systems, resilient to failures
  First delivery: side effect happens
  Duplicate delivery: side effect is skipped (already done)
  Net result: side effect happened exactly once ✓
```

---

## Summary

- **At-least-once** is the practical delivery guarantee: every message is delivered, but duplicates can occur on crash/restart
- **Duplicate sources**: consumer crash before offset commit, broker restart, rebalance, manual replay
- **Idempotency key**: a stable unique ID per operation; the server deduplicates by storing it; use `Idempotency-Key` HTTP header for APIs
- **DB deduplication**: `INSERT INTO processed_events (id) ... ON CONFLICT DO NOTHING` — check `rowsAffected == 1` to detect new vs duplicate
- **Redis deduplication**: `SETNX event:{id} 1 EX 86400` — fast, ephemeral, no transaction support
- **Natural idempotency**: `PUT`, `DELETE`, `UPSERT`, and state-guarded `UPDATE` are safe to repeat
- **Handler design**: make the side effect conditional on current state, not on message receipt
- **Consumer pattern**: check dedup → process → ack; on error, undo dedup and don't commit offset
- **Exactly-once** requires 2PC across system boundaries — fragile; at-least-once + idempotency is the correct answer

---

## Exercises

### Easy
1. Write a `processIfNew` function that uses `ON CONFLICT DO NOTHING` and returns `(bool, error)`. Test it by calling it twice with the same event ID — the second call should return `false, nil` without running any business logic.
2. Implement a `RedisDeduplicator` with a configurable TTL. Write a test that verifies `MarkSeen` returns `true` the first time and `false` on repeat calls within the TTL window.
3. Identify which of these operations are idempotent: `INSERT INTO users(email)`, `UPDATE users SET last_login=NOW() WHERE id=$1`, `DELETE FROM sessions WHERE token=$1`, `UPDATE wallets SET balance = balance - $1 WHERE id=$2`. Explain why.

### Medium
4. Build a full `IdempotentConsumer` that reads from a Kafka topic (or channel-based mock), deduplicates via PostgreSQL, and processes `ChargeEvent` messages. Write an integration test that sends the same event 5 times and asserts it was only processed once.
5. Implement the HTTP `Idempotency-Key` server-side pattern: store `(key → response)` in Redis. A `POST /payments` with a repeated key must return the cached response without calling the payment provider again. Add a test that verifies this by checking the mock provider was called exactly once.
6. Design and implement a state-machine guard for an `Order` with states `[pending, confirmed, fulfilled, shipped, cancelled]`. Write a `Transition(orderID, newStatus string)` function that is idempotent — applying the same transition twice is a no-op, not an error.

### Hard
7. Implement a **transactional deduplication pattern**: in a single PostgreSQL transaction, insert into `processed_events` AND into your business table. The dedup insert uses `ON CONFLICT DO NOTHING`; if 0 rows are inserted, skip the business write and commit (no-op). Write a test using `pgx` with a real database that verifies this under concurrent goroutines — 10 goroutines all send the same event, and only one business row appears.
8. Build a **dead-letter-aware idempotent consumer**: after 3 delivery attempts for the same event ID (tracked in Redis with a counter), route the message to a dead-letter topic instead of retrying indefinitely. Implement `ShouldDeadLetter(ctx, eventID string) bool` using Redis `INCR` + `EXPIRE`, and wire it into the consumer loop.
