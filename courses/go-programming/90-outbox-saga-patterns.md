# Chapter 89: Outbox Pattern and Saga Pattern

Two distributed systems patterns that appear together in microservices: the **Outbox Pattern** ensures events are published reliably alongside database writes (no lost events, no phantom events). The **Saga Pattern** coordinates multi-step transactions across services without a distributed lock.

## Table of Contents

1. [The Problem: Dual Write](#1-the-problem-dual-write)
2. [Outbox Pattern](#2-outbox-pattern)
3. [Outbox Relay in Go](#3-outbox-relay-in-go)
4. [Saga Pattern](#4-saga-pattern)
5. [Choreography vs Orchestration](#5-choreography-vs-orchestration)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. The Problem: Dual Write

```go
// The naive approach: two independent writes
func (s *OrderService) PlaceOrder(ctx context.Context, order *Order) error {
    // Step 1: save to database
    if err := s.orders.Save(ctx, order); err != nil { return err }
    
    // Step 2: publish event to Kafka
    if err := s.bus.Publish(ctx, OrderPlacedEvent{OrderID: order.ID}); err != nil {
        // PROBLEM: order is saved but event is not published
        // Order exists in DB but inventory service never gets notified
        return err
    }
    return nil
}

// If the process crashes between step 1 and step 2:
// - Database has the order ✓
// - Kafka has no event ✗
// - Inventory is never reserved ✗
// Result: ghost order that will never be fulfilled
```

You cannot make a database write and a Kafka write atomically — they're two different systems with no shared transaction.

---

## 2. Outbox Pattern

Solution: write the event to a database table (`outbox`) in the **same transaction** as the business data. A separate relay process reads the outbox and publishes to the message broker.

```sql
-- The outbox table
CREATE TABLE outbox (
    id          BIGSERIAL PRIMARY KEY,
    topic       TEXT NOT NULL,
    key         TEXT,
    payload     JSONB NOT NULL,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ  -- NULL = not yet published
);
```

### Writing to the outbox in the same transaction

```go
func (s *OrderService) PlaceOrder(ctx context.Context, order *Order) error {
    return s.db.WithTx(ctx, func(tx *sqlx.Tx) error {
        // Step 1: save order (in transaction)
        if err := s.orders.SaveTx(ctx, tx, order); err != nil { return err }
        
        // Step 2: write event to outbox (same transaction!)
        event := OrderPlacedEvent{
            OrderID:    order.ID,
            CustomerID: order.CustomerID,
            Total:      order.Total,
            At:         time.Now(),
        }
        payload, err := json.Marshal(event)
        if err != nil { return err }
        
        _, err = tx.ExecContext(ctx, `
            INSERT INTO outbox (topic, key, payload)
            VALUES ($1, $2, $3)`,
            "orders.placed", order.ID, payload,
        )
        return err
        // If this commit succeeds: both order AND outbox entry exist
        // If this fails: neither exists (atomicity)
    })
}
```

---

## 3. Outbox Relay in Go

A relay process polls the outbox, publishes messages, and marks them as published.

```go
type OutboxRelay struct {
    db     *sqlx.DB
    kafka  *kafka.Writer  // or any message broker
    logger *slog.Logger
}

type OutboxEntry struct {
    ID        int64           `db:"id"`
    Topic     string          `db:"topic"`
    Key       string          `db:"key"`
    Payload   []byte          `db:"payload"`
    CreatedAt time.Time       `db:"created_at"`
}

func (r *OutboxRelay) Run(ctx context.Context) {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop()
    
    for {
        select {
        case <-ctx.Done():
            return
        case <-ticker.C:
            if err := r.processOnce(ctx); err != nil {
                r.logger.Error("outbox relay error", "err", err)
            }
        }
    }
}

func (r *OutboxRelay) processOnce(ctx context.Context) error {
    // Lock unpublished entries for this relay instance (prevents double-publish with multiple replicas)
    rows, err := r.db.QueryxContext(ctx, `
        SELECT id, topic, key, payload, created_at
        FROM outbox
        WHERE published_at IS NULL
        ORDER BY id ASC
        LIMIT 100
        FOR UPDATE SKIP LOCKED`)
    if err != nil { return err }
    defer rows.Close()
    
    var entries []OutboxEntry
    for rows.Next() {
        var e OutboxEntry
        if err := rows.StructScan(&e); err != nil { return err }
        entries = append(entries, e)
    }
    if len(entries) == 0 { return nil }
    
    // Publish to Kafka
    msgs := make([]kafka.Message, len(entries))
    for i, e := range entries {
        msgs[i] = kafka.Message{
            Topic: e.Topic,
            Key:   []byte(e.Key),
            Value: e.Payload,
        }
    }
    if err := r.kafka.WriteMessages(ctx, msgs...); err != nil {
        return fmt.Errorf("publish: %w", err)
    }
    
    // Mark as published
    ids := make([]int64, len(entries))
    for i, e := range entries { ids[i] = e.ID }
    
    _, err = r.db.ExecContext(ctx, `
        UPDATE outbox SET published_at = NOW()
        WHERE id = ANY($1)`,
        pq.Array(ids),
    )
    return err
}

// Cleanup old published entries (run periodically)
func (r *OutboxRelay) Cleanup(ctx context.Context, olderThan time.Duration) error {
    _, err := r.db.ExecContext(ctx, `
        DELETE FROM outbox
        WHERE published_at IS NOT NULL
          AND published_at < $1`,
        time.Now().Add(-olderThan),
    )
    return err
}
```

### PostgreSQL LISTEN/NOTIFY for faster relay

Instead of polling every 500ms, use PostgreSQL's `NOTIFY` to wake the relay immediately:

```sql
-- Trigger on outbox INSERT
CREATE OR REPLACE FUNCTION notify_outbox()
RETURNS TRIGGER AS $$
BEGIN
    PERFORM pg_notify('outbox_insert', NEW.id::text);
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER outbox_notify
AFTER INSERT ON outbox
FOR EACH ROW EXECUTE FUNCTION notify_outbox();
```

```go
func (r *OutboxRelay) RunWithNotify(ctx context.Context) {
    conn, _ := r.db.Conn(ctx)
    conn.ExecContext(ctx, "LISTEN outbox_insert")
    
    for {
        // Wait for notification (with timeout fallback)
        // pq supports WaitForNotification
        notif, err := waitForNotify(ctx, conn, 5*time.Second)
        if err != nil { continue } // timeout or context cancelled
        _ = notif
        r.processOnce(ctx)
    }
}
```

---

## 4. Saga Pattern

A saga breaks a long-running transaction into a sequence of local transactions, each with a compensating transaction that reverses it on failure.

```
Order Saga:
  1. ReserveInventory  → compensate: ReleaseInventory
  2. ChargePAyment     → compensate: RefundPayment
  3. ConfirmOrder      → compensate: CancelOrder
  4. NotifyWarehouse   → (no compensation needed — idempotent)

If step 3 fails:
  - Undo step 2: RefundPayment
  - Undo step 1: ReleaseInventory
```

### Orchestration Saga

The orchestrator drives the saga, calling each service and handling failures:

```go
type OrderSaga struct {
    inventory InventoryService
    payments  PaymentService
    orders    OrderRepository
    notifier  NotificationService
    logger    *slog.Logger
}

type SagaState struct {
    OrderID            string
    InventoryReserved  bool
    PaymentCharged     bool
    Status             string
    Error              string
}

func (s *OrderSaga) Execute(ctx context.Context, order *Order) error {
    state := &SagaState{OrderID: order.ID}
    
    // Step 1: Reserve inventory
    if err := s.inventory.Reserve(ctx, order.ID, order.Items); err != nil {
        state.Error = fmt.Sprintf("inventory: %v", err)
        return s.compensate(ctx, state)
    }
    state.InventoryReserved = true
    s.logger.Info("saga: inventory reserved", "order_id", order.ID)
    
    // Step 2: Charge payment
    if err := s.payments.Charge(ctx, order.CustomerID, order.Total); err != nil {
        state.Error = fmt.Sprintf("payment: %v", err)
        return s.compensate(ctx, state)
    }
    state.PaymentCharged = true
    s.logger.Info("saga: payment charged", "order_id", order.ID)
    
    // Step 3: Confirm order
    if err := s.orders.Confirm(ctx, order.ID); err != nil {
        state.Error = fmt.Sprintf("confirm: %v", err)
        return s.compensate(ctx, state)
    }
    
    // Step 4: Notify warehouse (idempotent, best-effort)
    s.notifier.NotifyWarehouse(ctx, order.ID)
    
    s.logger.Info("saga: completed successfully", "order_id", order.ID)
    return nil
}

func (s *OrderSaga) compensate(ctx context.Context, state *SagaState) error {
    s.logger.Warn("saga: compensating", "order_id", state.OrderID, "error", state.Error)
    
    if state.PaymentCharged {
        if err := s.payments.Refund(ctx, state.OrderID); err != nil {
            s.logger.Error("compensation failed: refund", "err", err)
            // Compensation failures need manual intervention or a retry queue
        }
    }
    
    if state.InventoryReserved {
        if err := s.inventory.Release(ctx, state.OrderID); err != nil {
            s.logger.Error("compensation failed: release inventory", "err", err)
        }
    }
    
    return fmt.Errorf("saga failed: %s", state.Error)
}
```

---

## 5. Choreography vs Orchestration

### Choreography (event-driven)

Each service reacts to events autonomously. No central coordinator.

```
OrderService publishes OrderPlaced
  → InventoryService listens, reserves stock, publishes InventoryReserved
      → PaymentService listens, charges card, publishes PaymentCharged
          → OrderService listens, confirms order, publishes OrderConfirmed
              → WarehouseService listens, creates shipment

On failure:
  PaymentService publishes PaymentFailed
    → InventoryService listens, releases reservation
    → OrderService listens, marks order as failed
```

```go
// Each service just handles its events
type InventoryEventHandler struct{ inventory *InventoryService }

func (h *InventoryEventHandler) OnOrderPlaced(ctx context.Context, e OrderPlacedEvent) error {
    if err := h.inventory.Reserve(ctx, e.OrderID, e.Items); err != nil {
        // Publish compensating event
        h.bus.Publish(ctx, InventoryReservationFailedEvent{
            OrderID: e.OrderID,
            Reason:  err.Error(),
        })
        return nil
    }
    h.bus.Publish(ctx, InventoryReservedEvent{OrderID: e.OrderID})
    return nil
}
```

| | Orchestration | Choreography |
|---|---|---|
| Control | Central orchestrator | Each service is autonomous |
| Visibility | Easy — one place shows state | Harder — distributed across services |
| Coupling | Services coupled to orchestrator | Services coupled to event schema |
| Failure handling | Explicit compensation logic | Compensating events |
| Best for | Complex workflows with many decisions | Simple sequential flows |

---

## Summary

- **Dual write problem**: you cannot atomically write to both a database and a message broker
- **Outbox Pattern**: write events to an `outbox` table in the same DB transaction; relay them to the broker asynchronously
- **FOR UPDATE SKIP LOCKED**: allows multiple relay instances without double-publishing
- **Saga**: replace a distributed ACID transaction with a sequence of local transactions + compensations
- **Orchestration**: central saga coordinator calls each service step-by-step
- **Choreography**: each service reacts to events autonomously; compensating events undo failures

## Exercises

### Easy
1. Create the `outbox` table in PostgreSQL. Write a function `WriteToOutbox(tx *sqlx.Tx, topic, key string, payload any) error` that serializes the payload and inserts it into the outbox within a transaction.
2. Implement a simple polling relay that reads 10 entries from the outbox, logs their payloads, and marks them as published. Run it every 1 second.
3. Write a unit test that verifies that if the database commit succeeds, the outbox entry exists; and if the commit is rolled back, the outbox entry does not exist.

### Medium
4. Implement the full outbox relay with `FOR UPDATE SKIP LOCKED`. Start three relay goroutines simultaneously. Verify that each outbox entry is published exactly once using an in-memory counter.
5. Build an orchestrated `AccountTransferSaga`: `DebitSource → CreditDestination`. If `CreditDestination` fails, compensate by `CreditSource`. Test with mock services that fail on demand.
6. Implement a **saga state machine**: store the current saga step in a `sagas` table (`saga_id`, `order_id`, `current_step`, `status`). Resume a saga from where it left off if the process crashes mid-execution.

### Hard
7. Implement a **choreography saga** for order processing: 3 services (inventory, payment, shipping) each run in their own goroutine and communicate only via channels (simulating event buses). Wire the compensation flow so that a payment failure automatically releases inventory.
8. Build an **idempotent relay**: the outbox relay should be safe to run multiple times for the same entry (in case of relay crash and restart). Use the Kafka `transactional.id` producer feature or a deduplication table to guarantee exactly-once delivery.
