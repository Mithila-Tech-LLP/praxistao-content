# Chapter 91: Major Project 1 — Food Delivery Backend

This project brings together everything from Vol 7 (Clean Architecture): the domain model, use cases, repository pattern, CQRS read models, the outbox pattern for event publishing, config management, and dependency injection.

You'll build the core of a food delivery platform: customers browse restaurants, place orders, and track delivery status. Restaurants receive and fulfill orders. Drivers pick up and deliver.

---

## Architecture Overview

```
┌─────────────────────────────────────────────────────────┐
│                   Clean Architecture                     │
├─────────────┬──────────────┬───────────┬────────────────┤
│   Domain    │  Use Cases   │ Adapters  │ Infrastructure │
│             │              │  (HTTP)   │                │
│ Restaurant  │ PlaceOrder   │           │  PostgreSQL     │
│ Order       │ ConfirmOrder │ REST API  │  Redis cache   │
│ Driver      │ AssignDriver │           │  Outbox relay  │
│ Customer    │ DeliverOrder │           │  Event bus     │
└─────────────┴──────────────┴───────────┴────────────────┘
```

---

## Project Structure

```
food-delivery/
├── domain/
│   ├── order.go          ← Order aggregate with state machine
│   ├── restaurant.go     ← Restaurant entity
│   ├── driver.go         ← Driver entity
│   └── repository.go     ← Repository interfaces
├── usecase/
│   ├── order_service.go
│   ├── restaurant_service.go
│   └── driver_service.go
├── adapter/
│   └── http/
│       ├── order_handler.go
│       ├── restaurant_handler.go
│       └── driver_handler.go
├── infrastructure/
│   ├── postgres/
│   │   ├── db.go
│   │   ├── order_repo.go
│   │   └── outbox.go
│   └── redis/
│       └── order_cache.go
├── projection/
│   └── order_summary.go  ← CQRS read model
├── config/
│   └── config.go
└── cmd/
    ├── api/main.go
    └── relay/main.go     ← Outbox relay process
```

---

## Domain Model

### Order State Machine

```go
// domain/order.go
package domain

import (
    "errors"
    "fmt"
    "time"
)

type OrderStatus string

const (
    OrderStatusPending    OrderStatus = "pending"    // placed, awaiting restaurant
    OrderStatusConfirmed  OrderStatus = "confirmed"  // restaurant accepted
    OrderStatusPreparing  OrderStatus = "preparing"  // kitchen is cooking
    OrderStatusReady      OrderStatus = "ready"      // waiting for driver pickup
    OrderStatusPickedUp   OrderStatus = "picked_up"  // driver has it
    OrderStatusDelivered  OrderStatus = "delivered"  // done
    OrderStatusCancelled  OrderStatus = "cancelled"  // cancelled at any stage
)

type Order struct {
    id           OrderID
    customerID   CustomerID
    restaurantID RestaurantID
    driverID     *DriverID // nil until assigned
    items        []OrderItem
    status       OrderStatus
    total        Money
    placedAt     time.Time
    confirmedAt  *time.Time
    pickedUpAt   *time.Time
    deliveredAt  *time.Time
    events       []DomainEvent
}

// Allowed transitions:
var transitions = map[OrderStatus][]OrderStatus{
    OrderStatusPending:   {OrderStatusConfirmed, OrderStatusCancelled},
    OrderStatusConfirmed: {OrderStatusPreparing, OrderStatusCancelled},
    OrderStatusPreparing: {OrderStatusReady},
    OrderStatusReady:     {OrderStatusPickedUp},
    OrderStatusPickedUp:  {OrderStatusDelivered},
}

func (o *Order) transitionTo(next OrderStatus) error {
    allowed := transitions[o.status]
    for _, s := range allowed {
        if s == next { o.status = next; return nil }
    }
    return fmt.Errorf("invalid transition %s → %s", o.status, next)
}

func (o *Order) Confirm(restaurantID RestaurantID) error {
    if o.restaurantID != restaurantID {
        return errors.New("wrong restaurant")
    }
    if err := o.transitionTo(OrderStatusConfirmed); err != nil { return err }
    
    now := time.Now()
    o.confirmedAt = &now
    o.emit(OrderConfirmedEvent{OrderID: o.id, RestaurantID: restaurantID, At: now})
    return nil
}

func (o *Order) AssignDriver(driverID DriverID) error {
    if o.status != OrderStatusReady {
        return fmt.Errorf("can only assign driver when order is ready, got %s", o.status)
    }
    o.driverID = &driverID
    return nil
}

func (o *Order) MarkPickedUp(driverID DriverID) error {
    if o.driverID == nil || *o.driverID != driverID {
        return errors.New("driver not assigned to this order")
    }
    if err := o.transitionTo(OrderStatusPickedUp); err != nil { return err }
    now := time.Now()
    o.pickedUpAt = &now
    o.emit(OrderPickedUpEvent{OrderID: o.id, DriverID: driverID, At: now})
    return nil
}

func (o *Order) Deliver(driverID DriverID) error {
    if o.driverID == nil || *o.driverID != driverID {
        return errors.New("driver not assigned to this order")
    }
    if err := o.transitionTo(OrderStatusDelivered); err != nil { return err }
    now := time.Now()
    o.deliveredAt = &now
    o.emit(OrderDeliveredEvent{OrderID: o.id, At: now})
    return nil
}

func (o *Order) Cancel(reason string) error {
    if _, ok := transitions[o.status]; !ok { return errors.New("cannot cancel at this stage") }
    if err := o.transitionTo(OrderStatusCancelled); err != nil { return err }
    o.emit(OrderCancelledEvent{OrderID: o.id, Reason: reason, At: time.Now()})
    return nil
}

func (o *Order) emit(e DomainEvent) { o.events = append(o.events, e) }
func (o *Order) PopEvents() []DomainEvent {
    evts := o.events; o.events = nil; return evts
}

// Domain events
type OrderPlacedEvent    struct { OrderID OrderID; CustomerID CustomerID; RestaurantID RestaurantID; Total Money; At time.Time }
type OrderConfirmedEvent  struct { OrderID OrderID; RestaurantID RestaurantID; At time.Time }
type OrderPickedUpEvent   struct { OrderID OrderID; DriverID DriverID; At time.Time }
type OrderDeliveredEvent  struct { OrderID OrderID; At time.Time }
type OrderCancelledEvent  struct { OrderID OrderID; Reason string; At time.Time }

func (e OrderPlacedEvent) EventName() string    { return "order.placed" }
func (e OrderConfirmedEvent) EventName() string  { return "order.confirmed" }
func (e OrderPickedUpEvent) EventName() string   { return "order.picked_up" }
func (e OrderDeliveredEvent) EventName() string  { return "order.delivered" }
func (e OrderCancelledEvent) EventName() string  { return "order.cancelled" }
```

---

## Use Case: Place Order

```go
// usecase/order_service.go
package usecase

type PlaceOrderInput struct {
    CustomerID   string
    RestaurantID string
    Items        []OrderItemInput
}

type OrderItemInput struct {
    MenuItemID string
    Quantity   int
}

type PlaceOrderOutput struct {
    OrderID string
    Total   float64
    Status  string
}

type OrderService struct {
    orders      domain.OrderRepository
    restaurants domain.RestaurantRepository
    outbox      domain.OutboxWriter
    idgen       IDGenerator
}

func (s *OrderService) PlaceOrder(ctx context.Context, in PlaceOrderInput) (*PlaceOrderOutput, error) {
    // Load restaurant and validate it accepts orders
    restaurant, err := s.restaurants.GetByID(ctx, domain.RestaurantID(in.RestaurantID))
    if err != nil { return nil, fmt.Errorf("restaurant not found: %w", err) }
    if !restaurant.AcceptsOrders() {
        return nil, errors.New("restaurant is not accepting orders right now")
    }
    
    // Resolve menu items and calculate total
    items := make([]domain.OrderItem, 0, len(in.Items))
    total := domain.Money{Currency: "USD"}
    for _, i := range in.Items {
        menuItem, err := restaurant.GetMenuItem(domain.MenuItemID(i.MenuItemID))
        if err != nil { return nil, fmt.Errorf("menu item %s not found: %w", i.MenuItemID, err) }
        if !menuItem.Available { return nil, fmt.Errorf("item %s is not available", menuItem.Name) }
        
        items = append(items, domain.OrderItem{
            MenuItemID: menuItem.ID,
            Name:       menuItem.Name,
            Quantity:   i.Quantity,
            UnitPrice:  menuItem.Price,
        })
        total.Amount += menuItem.Price.Amount * int64(i.Quantity)
    }
    
    // Create order aggregate
    orderID := domain.OrderID(s.idgen.New())
    order, err := domain.NewOrder(orderID,
        domain.CustomerID(in.CustomerID),
        domain.RestaurantID(in.RestaurantID),
        items, total)
    if err != nil { return nil, err }
    
    // Persist order and outbox events in one transaction
    err = s.orders.WithTx(ctx, func(tx domain.Tx) error {
        if err := s.orders.SaveTx(ctx, tx, order); err != nil { return err }
        
        for _, event := range order.PopEvents() {
            if err := s.outbox.WriteTx(ctx, tx, event); err != nil { return err }
        }
        return nil
    })
    if err != nil { return nil, err }
    
    return &PlaceOrderOutput{
        OrderID: string(orderID),
        Total:   float64(total.Amount) / 100,
        Status:  string(domain.OrderStatusPending),
    }, nil
}
```

---

## CQRS: Order Summary Projection

```go
// projection/order_summary.go
package projection

type OrderSummary struct {
    OrderID      string    `db:"order_id" json:"order_id"`
    CustomerID   string    `db:"customer_id" json:"customer_id"`
    RestaurantName string  `db:"restaurant_name" json:"restaurant_name"`
    Status       string    `db:"status" json:"status"`
    Total        float64   `db:"total" json:"total"`
    ItemCount    int       `db:"item_count" json:"item_count"`
    PlacedAt     time.Time `db:"placed_at" json:"placed_at"`
    DeliveredAt  *time.Time `db:"delivered_at" json:"delivered_at,omitempty"`
}

type OrderProjection struct {
    db *sqlx.DB
}

func (p *OrderProjection) OnOrderPlaced(ctx context.Context, e domain.OrderPlacedEvent) error {
    _, err := p.db.ExecContext(ctx, `
        INSERT INTO order_summaries
            (order_id, customer_id, restaurant_id, status, total, item_count, placed_at)
        VALUES ($1, $2, $3, 'pending', $4, $5, $6)
        ON CONFLICT (order_id) DO NOTHING`,
        e.OrderID, e.CustomerID, e.RestaurantID,
        float64(e.Total.Amount)/100,
        e.ItemCount,
        e.At,
    )
    return err
}

func (p *OrderProjection) OnOrderConfirmed(ctx context.Context, e domain.OrderConfirmedEvent) error {
    _, err := p.db.ExecContext(ctx,
        "UPDATE order_summaries SET status='confirmed', confirmed_at=$1 WHERE order_id=$2",
        e.At, e.OrderID)
    return err
}

func (p *OrderProjection) OnOrderDelivered(ctx context.Context, e domain.OrderDeliveredEvent) error {
    _, err := p.db.ExecContext(ctx,
        "UPDATE order_summaries SET status='delivered', delivered_at=$1 WHERE order_id=$2",
        e.At, e.OrderID)
    return err
}
```

---

## HTTP API

```
POST   /orders                        → PlaceOrder
GET    /orders/:id                    → GetOrder (CQRS query)
GET    /customers/:id/orders          → ListCustomerOrders (CQRS query)

POST   /restaurants/:id/orders/:oid/confirm  → ConfirmOrder
POST   /restaurants/:id/orders/:oid/ready    → MarkOrderReady

POST   /drivers/:id/orders/:oid/pickup   → MarkPickedUp
POST   /drivers/:id/orders/:oid/deliver  → MarkDelivered
```

```go
// adapter/http/order_handler.go (abbreviated)
func (h *OrderHandler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
    var req PlaceOrderRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        writeError(w, 400, "invalid JSON")
        return
    }
    
    customerID := auth.CustomerIDFromContext(r.Context())
    out, err := h.orders.PlaceOrder(r.Context(), usecase.PlaceOrderInput{
        CustomerID:   string(customerID),
        RestaurantID: req.RestaurantID,
        Items:        toUseCaseItems(req.Items),
    })
    if err != nil {
        switch {
        case errors.Is(err, domain.ErrNotFound): writeError(w, 404, err.Error())
        default: writeError(w, 500, "order placement failed")
        }
        return
    }
    writeJSON(w, 201, out)
}
```

---

## Database Schema

```sql
-- orders table (write model — normalized)
CREATE TABLE orders (
    id            TEXT PRIMARY KEY,
    customer_id   TEXT NOT NULL,
    restaurant_id TEXT NOT NULL,
    driver_id     TEXT,
    status        TEXT NOT NULL DEFAULT 'pending',
    total_cents   BIGINT NOT NULL,
    placed_at     TIMESTAMPTZ NOT NULL,
    confirmed_at  TIMESTAMPTZ,
    picked_up_at  TIMESTAMPTZ,
    delivered_at  TIMESTAMPTZ,
    version       INT NOT NULL DEFAULT 1
);

CREATE TABLE order_items (
    order_id     TEXT REFERENCES orders(id) ON DELETE CASCADE,
    menu_item_id TEXT NOT NULL,
    name         TEXT NOT NULL,
    quantity     INT NOT NULL,
    unit_price   BIGINT NOT NULL,
    PRIMARY KEY (order_id, menu_item_id)
);

-- order_summaries table (read model — denormalized)
CREATE TABLE order_summaries (
    order_id        TEXT PRIMARY KEY,
    customer_id     TEXT NOT NULL,
    restaurant_id   TEXT NOT NULL,
    restaurant_name TEXT,
    status          TEXT NOT NULL,
    total           NUMERIC(10,2) NOT NULL,
    item_count      INT NOT NULL,
    placed_at       TIMESTAMPTZ NOT NULL,
    confirmed_at    TIMESTAMPTZ,
    delivered_at    TIMESTAMPTZ
);

CREATE INDEX idx_order_summaries_customer ON order_summaries (customer_id, placed_at DESC);
CREATE INDEX idx_order_summaries_restaurant ON order_summaries (restaurant_id, status);

-- outbox
CREATE TABLE outbox (
    id           BIGSERIAL PRIMARY KEY,
    stream_id    TEXT NOT NULL,
    event_type   TEXT NOT NULL,
    payload      JSONB NOT NULL,
    created_at   TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    published_at TIMESTAMPTZ
);

CREATE INDEX idx_outbox_unpublished ON outbox (id) WHERE published_at IS NULL;
```

---

## What You Built

This project demonstrates:

1. **Domain model** (`order.go`): state machine enforced by domain methods, events emitted on transitions
2. **Use cases** (`order_service.go`): orchestrate domain + repos + outbox in a single transaction
3. **Outbox pattern** (`outbox.go`): events written atomically with business data
4. **CQRS** (`order_summaries`): separate read model maintained by event projections
5. **Clean architecture**: domain ← use cases ← adapters ← infrastructure; all dependencies point inward
6. **Repository pattern**: `domain.OrderRepository` interface satisfied by `postgres.OrderRepository`
7. **Config management**: Viper-loaded config bound to typed struct, validated at startup

---

## Extension Challenges

1. **Add real-time tracking**: when an order status changes, push an update via WebSocket to the customer's browser. Use Redis Pub/Sub: the HTTP handler publishes the status change, a WebSocket hub subscribes and forwards to connected clients.

2. **Driver matching**: when an order is ready, a background worker finds the nearest available driver using a Redis sorted set (`driver_locations` with score = geo-hash). Assign the first match.

3. **Restaurant analytics**: build a projection `restaurant_daily_stats(restaurant_id, date, total_orders, total_revenue, avg_delivery_time_minutes)`. Update it when `OrderDeliveredEvent` arrives. Expose as `GET /restaurants/:id/stats?from=2026-01-01&to=2026-06-30`.

4. **Retry logic for the outbox relay**: if publishing to Kafka fails, implement exponential backoff with a dead-letter entry after 5 failures. Dead-letter entries are logged and alerted.
