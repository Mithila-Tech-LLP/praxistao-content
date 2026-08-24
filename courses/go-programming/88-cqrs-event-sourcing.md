# Chapter 88: CQRS and Event Sourcing

CQRS (Command Query Responsibility Segregation) splits your model into two: a write model that handles commands and a read model optimized for queries. Event Sourcing stores every state change as an immutable event instead of updating rows in place — the current state is derived by replaying events.

They're often used together but are independent patterns. You can do CQRS without Event Sourcing, and Event Sourcing without CQRS.

## Table of Contents

1. [CQRS Fundamentals](#1-cqrs-fundamentals)
2. [Commands and Queries](#2-commands-and-queries)
3. [Read Model Projections](#3-read-model-projections)
4. [Event Sourcing](#4-event-sourcing)
5. [Event Store in Go](#5-event-store-in-go)
6. [When to Use](#6-when-to-use)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. CQRS Fundamentals

Traditional CRUD: one model handles reads and writes. The `User` struct is used for both `UpdateUser(user)` and `GetUserProfile()`.

Problems at scale:
- Read and write load differ — you might need 50 read replicas but only 3 write nodes
- Read queries need data joined across many tables; write transactions need to be fast
- Optimizing for one often hurts the other

CQRS solution: separate models.

```
                  ┌─────────────┐
                  │   Client    │
                  └──────┬──────┘
                         │
             ┌───────────┴────────────┐
             ▼                        ▼
      Commands (writes)          Queries (reads)
             │                        │
     ┌───────▼────────┐     ┌─────────▼──────────┐
     │  Write Model   │     │    Read Model       │
     │  (aggregate)   │─────│  (denormalized      │
     │                │event│   projections)      │
     └───────┬────────┘     └─────────────────────┘
             │
       ┌─────▼──────┐
       │  Event     │
       │  Store     │
       └────────────┘
```

---

## 2. Commands and Queries

### Commands: intent to change state

```go
// A command is a request to do something. It may be rejected.
type Command interface{ commandName() string }

type CreateOrderCommand struct {
    CustomerID string
    Items      []OrderItemInput
}

type SubmitOrderCommand struct {
    OrderID string
}

type CancelOrderCommand struct {
    OrderID string
    Reason  string
}

// Command handler: validates, loads aggregate, calls domain method, saves
type OrderCommandHandler struct {
    orders OrderRepository // write model repository
    bus    EventBus
}

func (h *OrderCommandHandler) HandleSubmitOrder(ctx context.Context, cmd SubmitOrderCommand) error {
    order, err := h.orders.Load(ctx, cmd.OrderID)
    if err != nil { return err }
    
    if err := order.Submit(); err != nil { return err } // domain validation
    
    if err := h.orders.Save(ctx, order); err != nil { return err }
    
    for _, event := range order.PopEvents() {
        h.bus.Publish(ctx, event)
    }
    return nil
}
```

### Queries: read-only, no side effects

```go
// A query just returns data. It never changes state.
type GetOrderQuery struct {
    OrderID string
}

type OrderSummary struct {
    OrderID    string
    CustomerID string
    Status     string
    Total      string
    ItemCount  int
    CreatedAt  time.Time
}

// Query handler reads from the optimized read model
type OrderQueryHandler struct {
    readDB *sqlx.DB // could be a separate read replica
}

func (h *OrderQueryHandler) GetOrder(ctx context.Context, q GetOrderQuery) (*OrderSummary, error) {
    var summary OrderSummary
    err := h.readDB.QueryRowxContext(ctx, `
        SELECT order_id, customer_id, status, total, item_count, created_at
        FROM order_summaries
        WHERE order_id = $1`,
        q.OrderID,
    ).StructScan(&summary)
    return &summary, err
}

// Read model can be denormalized for the exact query pattern
type CustomerOrdersView struct {
    CustomerID    string
    TotalOrders   int
    TotalSpent    float64
    LastOrderAt   time.Time
    RecentOrders  []OrderSummary
}
```

---

## 3. Read Model Projections

A projection listens to domain events and updates the read model. The read model is derived data — it can always be rebuilt from events.

```go
type OrderProjection struct {
    db *sqlx.DB
}

func (p *OrderProjection) OnOrderCreated(ctx context.Context, event OrderCreatedEvent) error {
    _, err := p.db.ExecContext(ctx, `
        INSERT INTO order_summaries (order_id, customer_id, status, total, item_count, created_at)
        VALUES ($1, $2, 'draft', 0, 0, $3)
        ON CONFLICT (order_id) DO NOTHING`,
        event.OrderID, event.CustomerID, event.OccurredAt(),
    )
    return err
}

func (p *OrderProjection) OnItemAdded(ctx context.Context, event OrderItemAddedEvent) error {
    _, err := p.db.ExecContext(ctx, `
        UPDATE order_summaries
        SET total = total + $1,
            item_count = item_count + $2,
            updated_at = $3
        WHERE order_id = $4`,
        event.ItemPrice*float64(event.Quantity),
        event.Quantity,
        event.OccurredAt(),
        event.OrderID,
    )
    return err
}

func (p *OrderProjection) OnOrderSubmitted(ctx context.Context, event OrderSubmittedEvent) error {
    _, err := p.db.ExecContext(ctx,
        "UPDATE order_summaries SET status = 'pending', submitted_at = $1 WHERE order_id = $2",
        event.OccurredAt(), event.OrderID,
    )
    return err
}
```

---

## 4. Event Sourcing

Instead of storing the current state, store the sequence of events that led to it. The current state is always computed by replaying events.

```go
// Events are the source of truth
type OrderCreatedEvent struct {
    OrderID    string    `json:"order_id"`
    CustomerID string    `json:"customer_id"`
    At         time.Time `json:"at"`
}

type OrderItemAddedEvent struct {
    OrderID   string    `json:"order_id"`
    ProductID string    `json:"product_id"`
    Name      string    `json:"name"`
    Quantity  int       `json:"quantity"`
    UnitPrice float64   `json:"unit_price"`
    At        time.Time `json:"at"`
}

type OrderSubmittedEvent struct {
    OrderID string    `json:"order_id"`
    Total   float64   `json:"total"`
    At      time.Time `json:"at"`
}

// The aggregate knows how to apply each event
type Order struct {
    id         string
    customerID string
    items      []OrderItem
    status     OrderStatus
    total      float64
    version    int  // optimistic concurrency
    changes    []StoredEvent // uncommitted events
}

// Apply replays a single event to update the aggregate state
func (o *Order) Apply(event StoredEvent) error {
    switch e := event.Data.(type) {
    case OrderCreatedEvent:
        o.id = e.OrderID
        o.customerID = e.CustomerID
        o.status = OrderDraft
    case OrderItemAddedEvent:
        o.items = append(o.items, OrderItem{
            ProductID: e.ProductID,
            Quantity:  e.Quantity,
            UnitPrice: e.UnitPrice,
        })
        o.total += e.UnitPrice * float64(e.Quantity)
    case OrderSubmittedEvent:
        o.status = OrderPending
    default:
        return fmt.Errorf("unknown event type: %T", event.Data)
    }
    o.version++
    return nil
}

// Reconstruct from event stream
func LoadOrder(events []StoredEvent) (*Order, error) {
    order := &Order{}
    for _, event := range events {
        if err := order.Apply(event); err != nil { return nil, err }
    }
    return order, nil
}

// Business method: validate, record event
func (o *Order) Submit() error {
    if o.status != OrderDraft   { return errors.New("not a draft") }
    if len(o.items) == 0        { return errors.New("empty order") }
    
    o.record(OrderSubmittedEvent{
        OrderID: o.id,
        Total:   o.total,
        At:      time.Now(),
    })
    return nil
}

func (o *Order) record(event any) {
    stored := StoredEvent{
        StreamID:  o.id,
        EventType: typeName(event),
        Data:      event,
        Version:   o.version + len(o.changes) + 1,
        OccurredAt: time.Now(),
    }
    o.changes = append(o.changes, stored)
}
```

---

## 5. Event Store in Go

```go
type StoredEvent struct {
    ID          int64
    StreamID    string     // e.g., "order:abc123"
    EventType   string     // "OrderSubmittedEvent"
    Data        any        // deserialized event
    RawData     []byte     // JSON bytes stored in DB
    Version     int        // position in this stream
    OccurredAt  time.Time
}

type EventStore interface {
    Append(ctx context.Context, streamID string, events []StoredEvent, expectedVersion int) error
    Load(ctx context.Context, streamID string) ([]StoredEvent, error)
    LoadFrom(ctx context.Context, streamID string, version int) ([]StoredEvent, error)
}

// PostgreSQL event store
type PostgresEventStore struct {
    db      *sqlx.DB
    registry EventRegistry
}

func (s *PostgresEventStore) Append(ctx context.Context, streamID string, events []StoredEvent, expected int) error {
    tx, err := s.db.BeginTxx(ctx, nil)
    if err != nil { return err }
    defer tx.Rollback()
    
    // Optimistic concurrency check
    var current int
    tx.QueryRowContext(ctx,
        "SELECT COALESCE(MAX(version), 0) FROM events WHERE stream_id = $1", streamID,
    ).Scan(&current)
    
    if current != expected {
        return fmt.Errorf("optimistic concurrency conflict: expected version %d, got %d", expected, current)
    }
    
    for _, event := range events {
        data, err := json.Marshal(event.Data)
        if err != nil { return err }
        
        _, err = tx.ExecContext(ctx, `
            INSERT INTO events (stream_id, event_type, data, version, occurred_at)
            VALUES ($1, $2, $3, $4, $5)`,
            streamID, event.EventType, data, event.Version, event.OccurredAt,
        )
        if err != nil { return err }
    }
    return tx.Commit()
}

func (s *PostgresEventStore) Load(ctx context.Context, streamID string) ([]StoredEvent, error) {
    rows, err := s.db.QueryxContext(ctx, `
        SELECT id, stream_id, event_type, data, version, occurred_at
        FROM events
        WHERE stream_id = $1
        ORDER BY version ASC`,
        streamID,
    )
    if err != nil { return nil, err }
    defer rows.Close()
    
    var events []StoredEvent
    for rows.Next() {
        var row struct {
            ID         int64     `db:"id"`
            StreamID   string    `db:"stream_id"`
            EventType  string    `db:"event_type"`
            Data       []byte    `db:"data"`
            Version    int       `db:"version"`
            OccurredAt time.Time `db:"occurred_at"`
        }
        if err := rows.StructScan(&row); err != nil { return nil, err }
        
        data, err := s.registry.Deserialize(row.EventType, row.Data)
        if err != nil { return nil, err }
        
        events = append(events, StoredEvent{
            ID:         row.ID,
            StreamID:   row.StreamID,
            EventType:  row.EventType,
            Data:       data,
            RawData:    row.Data,
            Version:    row.Version,
            OccurredAt: row.OccurredAt,
        })
    }
    return events, rows.Err()
}

// Event registry: maps type names to deserializers
type EventRegistry struct {
    deserializers map[string]func([]byte) (any, error)
}

func (r *EventRegistry) Register(name string, fn func([]byte) (any, error)) {
    r.deserializers[name] = fn
}

func (r *EventRegistry) Deserialize(name string, data []byte) (any, error) {
    fn, ok := r.deserializers[name]
    if !ok { return nil, fmt.Errorf("unknown event type: %s", name) }
    return fn(data)
}
```

---

## 6. When to Use

**Use CQRS when:**
- Read and write loads differ significantly
- Read models need to join many tables but writes need to be fast
- You have complex, reporting-heavy read requirements

**Use Event Sourcing when:**
- Complete audit trail is required (financial, healthcare, legal)
- You need to replay history to answer "what was the state at time T?"
- You need temporal queries or undo/redo
- Multiple downstream systems need the event stream

**Don't use Event Sourcing when:**
- Simple CRUD is sufficient
- Your team isn't familiar with it (steep learning curve)
- Data storage constraints are tight (events accumulate forever without snapshotting)

---

## Summary

- **CQRS**: separate write model (commands → aggregates) from read model (queries → projections)
- **Command**: intent to change state, may be rejected; goes through the domain model
- **Query**: read-only, no side effects; hits the optimized read model directly
- **Projection**: event handler that maintains the read model
- **Event Sourcing**: state = replay of events; current state is never stored directly
- **Optimistic concurrency**: `expectedVersion` prevents lost updates when two writers race

## Exercises

### Easy
1. Implement a `BankAccountProjection` that maintains a `balances` table: `{account_id, balance, last_updated}`. Subscribe to `MoneyDepositedEvent` and `MoneyWithdrawnEvent` and update the balance accordingly.
2. Write a function `Replay(events []StoredEvent, account *BankAccount)` that rebuilds a bank account's state from its event stream. Test it with a known sequence of events.
3. Implement `GetBalanceQuery` that reads from the `balances` projection table. Verify that after replaying events, the query returns the correct balance.

### Medium
4. Add **snapshots** to the event-sourced `Order`: after 50 events, take a snapshot of the current state and store it. On load, find the latest snapshot and replay only events after it.
5. Build a **time travel query**: given an `OrderID` and a timestamp, replay events up to that timestamp and return the order's state at that moment. This is only possible with event sourcing.
6. Implement a **competing consumer projection**: multiple instances of the projection worker process events. Use PostgreSQL advisory locks to ensure each event is processed exactly once.

### Hard
7. Build a complete CQRS+ES order management system: commands (`CreateOrder`, `AddItem`, `Submit`, `Cancel`), event store (PostgreSQL), projections (`order_summaries`, `customer_stats`), and query handlers. Wire everything together.
8. Implement **event upcasting**: your events have versioned schemas. `OrderCreatedEvent` v1 has `customer_id string`; v2 has `customer_id string` + `customer_email string`. Write an `Upcast` function that converts old event bytes to the latest version before they're deserialized, so old data works with new code without a migration.
