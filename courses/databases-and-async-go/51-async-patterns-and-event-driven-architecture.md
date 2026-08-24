# Chapter 51: Async Patterns and Event-Driven Architecture

Knowing the tools is one thing. Knowing the patterns is another. This chapter covers the recurring design patterns that appear in every async system — the building blocks you'll use over and over.

## Table of Contents

1. Event-Driven Architecture (EDA)
2. Pattern: Event Notification
3. Pattern: Event-Carried State Transfer
4. Pattern: Event Sourcing
5. Pattern: CQRS (Command Query Responsibility Segregation)
6. Pattern: Saga — Distributed Transactions
7. Pattern: Inbox/Outbox for Reliability
8. Anti-Patterns to Avoid
9. Exercises

---

## 1. Event-Driven Architecture (EDA)

Event-Driven Architecture means services communicate by emitting and reacting to events rather than calling each other directly.

**Service-to-service calls (synchronous):**
```
Order Service → calls → Email Service → sends email
Order Service → calls → Analytics Service → updates dashboard
Order Service → calls → Warehouse Service → reserves stock
```
Problems: Order Service knows about every downstream service. Tight coupling.

**Event-driven (async):**
```
Order Service → emits → "OrderPlaced" event

Email Service         → reacts to "OrderPlaced" → sends email
Analytics Service     → reacts to "OrderPlaced" → updates dashboard
Warehouse Service     → reacts to "OrderPlaced" → reserves stock
```
Benefits: Order Service doesn't know about any downstream service. Loose coupling.

---

## 2. Pattern: Event Notification

The simplest pattern. Emit an event to say "something happened." Include just enough data for consumers to know what happened. If they need more detail, they fetch it.

```go
type OrderPlacedEvent struct {
    OrderID   string    `json:"order_id"`
    UserID    string    `json:"user_id"`
    OccuredAt time.Time `json:"occurred_at"`
    // Note: no order items, no amount, no address
    // Consumers fetch these from the Order Service if needed
}
```

**When to use:** Decoupling services. When consumers don't need the full state, just a notification that something changed.

**Downside:** Consumers must make additional API calls to get full state. Creates chattiness.

---

## 3. Pattern: Event-Carried State Transfer

Include all the data consumers need in the event itself. Consumers don't need to call back.

```go
type OrderPlacedEvent struct {
    OrderID     string    `json:"order_id"`
    UserID      string    `json:"user_id"`
    UserEmail   string    `json:"user_email"`   // needed by email service
    Items       []Item    `json:"items"`         // needed by warehouse
    TotalAmount float64   `json:"total_amount"`  // needed by analytics
    ShippingAddr string   `json:"shipping_addr"` // needed by warehouse
    OccurredAt  time.Time `json:"occurred_at"`
}
```

**When to use:** When you want consumers to be fully autonomous. Good for fan-out patterns where many services need the same data.

**Downside:** Events become large. If the schema changes, all consumers must update.

---

## 4. Pattern: Event Sourcing

Instead of storing current state, store every event that ever happened. The current state is derived by replaying all events.

```
Normal DB approach:
  users table: {id: 1, name: "Alice", balance: 150}

Event Sourcing approach:
  events: [
    {type: "UserCreated", userID: 1, name: "Alice"},
    {type: "MoneyDeposited", userID: 1, amount: 200},
    {type: "MoneyWithdrawn", userID: 1, amount: 50},
  ]
  
  Current balance = 0 + 200 - 50 = 150 (derived by replaying)
```

**Benefits:**
- Complete audit trail. You can always answer "what was the balance on January 15th at 3pm?"
- Replay events to rebuild any read model or fix bugs.
- Easy integration with downstream systems — they just consume the event stream.

**Downsides:**
- Snapshots needed for large event histories (don't replay from the beginning every time).
- More complex than simple CRUD.

```go
// Event store interface
type EventStore interface {
    Append(aggregateID string, events []Event, expectedVersion int) error
    Load(aggregateID string, fromVersion int) ([]Event, error)
}

// Snapshot to avoid replaying from beginning
type Snapshot struct {
    AggregateID string
    Version     int
    State       []byte
}
```

---

## 5. Pattern: CQRS (Command Query Responsibility Segregation)

Separate the "write model" (command side) from the "read model" (query side).

```
Write Side (Command):               Read Side (Query):
  Order Service                       Order Read Model
  ↓                                   ↓
  Orders DB (normalized,              Orders Read DB
  optimized for writes)               (denormalized,
  ↓                                    optimized for queries)
  "OrderPlaced" event ─────────────►  Consumes events,
                                       rebuilds read model
```

**Why:** Write models need ACID transactions and normalization. Read models need denormalized, fast queries. CQRS lets each side be optimized independently.

```go
// Command side: processes writes
type PlaceOrderCommand struct {
    UserID string
    Items  []OrderItem
}

func (s *OrderService) PlaceOrder(cmd PlaceOrderCommand) (string, error) {
    order := NewOrder(cmd)
    s.ordersDB.Save(order) // normalized write
    s.eventBus.Publish(OrderPlacedEvent{...}) // notify read side
    return order.ID, nil
}

// Query side: serves reads
type OrderReadService struct {
    db *sql.DB // denormalized, pre-joined view
}

func (s *OrderReadService) GetOrdersForUser(userID string) ([]OrderSummary, error) {
    // Single query against a denormalized view — fast!
    return s.db.Query("SELECT * FROM order_summaries WHERE user_id = ?", userID)
}
```

---

## 6. Pattern: Saga — Distributed Transactions

You need to update multiple services atomically. Sagas coordinate this without a distributed transaction manager.

**Example: Place order across 3 services**

```
1. Reserve inventory (Inventory Service)
2. Charge payment (Payment Service)
3. Create shipment (Shipping Service)

If step 3 fails, we need to undo steps 1 and 2 (compensating transactions)
```

**Choreography-based saga (event-driven):**

```
1. Order Service → emits "OrderStarted"
2. Inventory Service → listens to "OrderStarted" → reserves stock → emits "StockReserved"
3. Payment Service → listens to "StockReserved" → charges card → emits "PaymentDone"
4. Shipping Service → listens to "PaymentDone" → creates shipment → emits "OrderCompleted"

If any step fails:
4. Shipping Service → emits "ShipmentFailed"
3. Payment Service → listens to "ShipmentFailed" → refunds → emits "PaymentRefunded"
2. Inventory Service → listens to "PaymentRefunded" → unreserves → emits "StockUnreserved"
```

**Orchestration-based saga (central coordinator):**

```go
type OrderSaga struct {
    orderID string
    step    int
}

func (s *OrderSaga) Execute(ctx context.Context) error {
    // Step 1
    if err := inventorySvc.Reserve(s.orderID); err != nil {
        return err // nothing to compensate
    }

    // Step 2
    if err := paymentSvc.Charge(s.orderID); err != nil {
        inventorySvc.Unreserve(s.orderID) // compensate step 1
        return err
    }

    // Step 3
    if err := shippingSvc.Create(s.orderID); err != nil {
        paymentSvc.Refund(s.orderID)   // compensate step 2
        inventorySvc.Unreserve(s.orderID) // compensate step 1
        return err
    }

    return nil
}
```

---

## 7. Pattern: Outbox — Guaranteed Event Delivery

**The problem:** How do you guarantee that after saving to the database, the event is also published to Kafka? What if the process crashes between the DB write and the Kafka publish?

```go
// WRONG: Can crash between DB write and Kafka publish
db.Save(order)      // succeeds
kafka.Publish(event) // CRASH — event never published!
// Now DB has the order but Kafka never got the event
```

**The Outbox Pattern:** Write the event to the database in the same transaction as the business data. A separate process (the outbox poller) reads from the outbox table and publishes to Kafka.

```sql
-- Part of your main DB schema
CREATE TABLE outbox (
    id         BIGSERIAL PRIMARY KEY,
    topic      TEXT NOT NULL,
    key        TEXT,
    payload    JSONB NOT NULL,
    created_at TIMESTAMPTZ DEFAULT NOW(),
    sent_at    TIMESTAMPTZ  -- NULL until published
);
```

```go
// CORRECT: Both operations in the same transaction
tx, _ := db.Begin()

// Save the order
tx.Exec("INSERT INTO orders (...) VALUES (?)", ...)

// Write the event to outbox (same transaction!)
tx.Exec("INSERT INTO outbox (topic, key, payload) VALUES ('orders', ?, ?)",
    orderID, eventJSON)

tx.Commit() // Both committed atomically

// Separate goroutine: poll outbox and publish to Kafka
func outboxPublisher(db *sql.DB, kafka KafkaProducer) {
    for {
        rows, _ := db.Query(
            "SELECT id, topic, key, payload FROM outbox WHERE sent_at IS NULL LIMIT 100")
        for rows.Next() {
            var id int64; var topic, key, payload string
            rows.Scan(&id, &topic, &key, &payload)
            kafka.Publish(topic, key, payload) // may retry on failure
            db.Exec("UPDATE outbox SET sent_at = NOW() WHERE id = ?", id)
        }
        time.Sleep(100 * time.Millisecond)
    }
}
```

This guarantees: either both the business data AND the event are committed, or neither is.

---

## 8. Anti-Patterns to Avoid

**Anti-pattern 1: Chatty events**
Don't emit an event for every DB row change. Emit meaningful business events.

```go
// BAD: technical event, not business event
emit("users_table_row_updated", {...})

// GOOD: business event
emit("UserEmailVerified", {...})
emit("UserUpgradedToPremium", {...})
```

**Anti-pattern 2: Ordered events across topics**
Kafka only guarantees order within a partition, not across topics or even across partitions.

```
// BAD: assuming events across two topics arrive in order
Topic A: "payment-charged"
Topic B: "inventory-reserved"
// These may arrive in any order to consumers
```

**Anti-pattern 3: Ignoring poison messages**
A "poison message" causes a consumer to crash every time it processes it. Without handling, the consumer infinitely retries and can't move forward.

```go
// GOOD: dead-letter queue for unprocessable messages
func processMessage(msg *kafka.Message) error {
    if err := doProcess(msg); err != nil {
        if isRetryable(err) {
            return err // retry
        }
        // Not retryable: send to dead-letter topic
        kafka.Produce("orders.dead-letter", msg.Key, msg.Value)
        return nil // don't retry
    }
    return nil
}
```

---

## Summary

- Event notification = "something happened" (minimal data). Event-carried state = full data in event.
- Event sourcing: store all events, derive current state by replay. Full audit trail + time-travel queries.
- CQRS: separate write model (normalized) from read model (denormalized). Each optimized independently.
- Saga: coordinate multi-service operations with compensating transactions on failure.
- Outbox pattern: write events to DB in same transaction as business data, then publish asynchronously. Guarantees no lost events.

### Exercises

**Easy:** For a bank transfer (debit account A, credit account B), design a saga. What are the steps? What are the compensating transactions if each step fails?

**Medium:** Implement the outbox pattern in Go: a function `SaveOrderWithEvent(db, order, event)` that writes both atomically, and an `OutboxPoller` goroutine that reads unsent events and publishes to a channel (simulate Kafka with a Go channel).

**Hard:** Design a CQRS system for an e-commerce order listing. Write model: normalized PostgreSQL tables. Read model: a denormalized Redis sorted set (orders sorted by creation time per user). Implement the event handler that updates the Redis read model when `OrderPlaced` events arrive.
