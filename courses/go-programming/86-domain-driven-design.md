# Chapter 86: Domain-Driven Design in Go

Domain-Driven Design (DDD) is a way of structuring software around the business problem it solves. The core idea: your code vocabulary should match the business vocabulary. When an order is "submitted", the code calls `order.Submit()`, not `db.SetStatus("pending")`.

DDD gives us concrete building blocks: entities, value objects, aggregates, domain events, and bounded contexts. Go's type system makes most of these patterns natural.

## Table of Contents

1. [Building Blocks](#1-building-blocks)
2. [Entities and Value Objects](#2-entities-and-value-objects)
3. [Aggregates and Aggregate Roots](#3-aggregates-and-aggregate-roots)
4. [Domain Events](#4-domain-events)
5. [Bounded Contexts](#5-bounded-contexts)
6. [Ubiquitous Language](#6-ubiquitous-language)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Building Blocks

| Concept | Definition | Go implementation |
|---------|------------|-------------------|
| Entity | Has identity; identity persists across time | Struct with `ID` field |
| Value Object | Defined by its value, not identity | Immutable struct, no ID |
| Aggregate | Cluster of entities with one root | Root struct containing child entities |
| Domain Event | Something that happened in the domain | Struct describing a past event |
| Repository | Persistence abstraction for aggregates | Interface in domain package |
| Domain Service | Logic that doesn't belong on any entity | Function/struct in domain package |
| Bounded Context | A cohesive sub-domain with its own language | Go module or package subtree |

---

## 2. Entities and Value Objects

### Entity — identified by ID, mutable over time

```go
// Order is an entity: it has an identity (ID) that persists even as it changes
type Order struct {
    id        OrderID
    customerID CustomerID
    items      []OrderItem
    status    OrderStatus
    createdAt time.Time
    events    []DomainEvent // uncommitted events
}

// Exported fields are accessible; state changes go through methods
func (o *Order) ID() OrderID       { return o.id }
func (o *Order) Status() OrderStatus { return o.status }
func (o *Order) Items() []OrderItem { return append([]OrderItem{}, o.items...) } // defensive copy

// Business logic lives on the entity — not in the service layer
func (o *Order) Submit() error {
    if o.status != OrderDraft {
        return fmt.Errorf("can only submit a draft order, got %s", o.status)
    }
    if len(o.items) == 0 {
        return errors.New("cannot submit an empty order")
    }
    o.status = OrderPending
    o.events = append(o.events, OrderSubmittedEvent{
        OrderID:    o.id,
        CustomerID: o.customerID,
        Total:      o.Total(),
        At:         time.Now(),
    })
    return nil
}
```

### Value Object — defined entirely by its value

```go
// Money is a value object: two Money{100, "USD"} values are equal
type Money struct {
    amount   int64  // cents — never use float for money
    currency string
}

func NewMoney(amount int64, currency string) (Money, error) {
    if amount < 0 { return Money{}, errors.New("amount cannot be negative") }
    if len(currency) != 3 { return Money{}, errors.New("currency must be 3-letter ISO code") }
    return Money{amount: amount, currency: strings.ToUpper(currency)}, nil
}

func (m Money) Add(other Money) (Money, error) {
    if m.currency != other.currency {
        return Money{}, fmt.Errorf("cannot add %s and %s", m.currency, other.currency)
    }
    return Money{amount: m.amount + other.amount, currency: m.currency}, nil
}

func (m Money) IsZero() bool { return m.amount == 0 }
func (m Money) Cents() int64 { return m.amount }
func (m Money) String() string {
    return fmt.Sprintf("%s %.2f", m.currency, float64(m.amount)/100)
}
// No setters — value objects are immutable. All operations return new values.

// Email as a value object (Ch 85)
type Email struct{ value string }
type OrderID struct{ value int64 }
type CustomerID struct{ value int64 }

func NewOrderID(v int64) OrderID       { return OrderID{v} }
func (o OrderID) Value() int64         { return o.value }
func (o OrderID) String() string       { return strconv.FormatInt(o.value, 10) }
```

---

## 3. Aggregates and Aggregate Roots

An aggregate is a cluster of related objects treated as a single unit for data changes. The aggregate root is the only entry point — all invariants must hold within the aggregate boundary.

```go
// Order is the aggregate root
// OrderItem and ShippingAddress are part of the Order aggregate
type Order struct {
    id         OrderID
    customerID CustomerID
    items      []OrderItem
    shipping   ShippingAddress // value object
    status     OrderStatus
    total      Money
    createdAt  time.Time
    events     []DomainEvent
}

type OrderItem struct {
    ProductID  ProductID
    Name       string
    Quantity   int
    UnitPrice  Money
}

func (i OrderItem) Subtotal() Money {
    return Money{amount: i.UnitPrice.Cents() * int64(i.Quantity), currency: i.UnitPrice.currency}
}

// Aggregate root enforces all invariants
func (o *Order) AddItem(productID ProductID, name string, qty int, price Money) error {
    if o.status != OrderDraft { return errors.New("cannot modify a non-draft order") }
    if qty <= 0              { return errors.New("quantity must be positive") }
    
    // Check if item already in order
    for i, item := range o.items {
        if item.ProductID == productID {
            o.items[i].Quantity += qty
            o.recalculateTotal()
            return nil
        }
    }
    
    o.items = append(o.items, OrderItem{
        ProductID: productID,
        Name:      name,
        Quantity:  qty,
        UnitPrice: price,
    })
    o.recalculateTotal()
    return nil
}

func (o *Order) recalculateTotal() {
    total := Money{currency: "USD"}
    for _, item := range o.items {
        sub := item.Subtotal()
        total.amount += sub.Cents()
    }
    o.total = total
}

func (o *Order) Total() Money { return o.total }

// Factory function — ensures invariants hold on creation
func NewOrder(id OrderID, customerID CustomerID, shipping ShippingAddress) (*Order, error) {
    if customerID.Value() == 0 { return nil, errors.New("customer ID required") }
    return &Order{
        id:         id,
        customerID: customerID,
        shipping:   shipping,
        status:     OrderDraft,
        total:      Money{currency: "USD"},
        createdAt:  time.Now(),
    }, nil
}
```

### Repository per aggregate root (not per entity)

```go
// Repository is for the whole aggregate, not individual OrderItems
type OrderRepository interface {
    Save(ctx context.Context, order *Order) error
    FindByID(ctx context.Context, id OrderID) (*Order, error)
    FindByCustomer(ctx context.Context, customerID CustomerID) ([]*Order, error)
}

// You never have an OrderItemRepository — items are loaded with the order
```

---

## 4. Domain Events

Domain events record something that happened. They're past tense, immutable, and carry enough context to act on them.

```go
// domain/events.go
type DomainEvent interface {
    EventName() string
    OccurredAt() time.Time
}

type OrderSubmittedEvent struct {
    OrderID    OrderID
    CustomerID CustomerID
    Total      Money
    At         time.Time
}

func (e OrderSubmittedEvent) EventName() string    { return "order.submitted" }
func (e OrderSubmittedEvent) OccurredAt() time.Time { return e.At }

type OrderCancelledEvent struct {
    OrderID CustomerID
    Reason  string
    At      time.Time
}

func (e OrderCancelledEvent) EventName() string    { return "order.cancelled" }
func (e OrderCancelledEvent) OccurredAt() time.Time { return e.At }

// Aggregate collects events; they're dispatched AFTER successful persistence
func (o *Order) PopEvents() []DomainEvent {
    events := o.events
    o.events = nil
    return events
}

// In the use case: save → dispatch events
func (s *OrderService) Submit(ctx context.Context, orderID OrderID) error {
    order, err := s.orders.FindByID(ctx, orderID)
    if err != nil { return err }
    
    if err := order.Submit(); err != nil { return err }
    
    if err := s.orders.Save(ctx, order); err != nil { return err }
    
    // Dispatch after successful save
    for _, event := range order.PopEvents() {
        s.bus.Publish(ctx, event)
    }
    return nil
}
```

### Event Bus

```go
type EventBus interface {
    Publish(ctx context.Context, event DomainEvent)
    Subscribe(eventName string, handler func(context.Context, DomainEvent))
}

// Synchronous in-process bus (for testing and simple setups)
type InMemoryBus struct {
    mu       sync.RWMutex
    handlers map[string][]func(context.Context, DomainEvent)
}

func (b *InMemoryBus) Subscribe(name string, h func(context.Context, DomainEvent)) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[name] = append(b.handlers[name], h)
}

func (b *InMemoryBus) Publish(ctx context.Context, e DomainEvent) {
    b.mu.RLock()
    handlers := b.handlers[e.EventName()]
    b.mu.RUnlock()
    
    for _, h := range handlers {
        h(ctx, e)
    }
}
```

---

## 5. Bounded Contexts

A bounded context defines the boundary within which a specific domain model applies. Different contexts can use the same word to mean different things.

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Order Context │    │ Catalog Context  │    │  Shipping Context│
│                 │    │                  │    │                 │
│ Product = line  │    │ Product = item   │    │ Package = box   │
│ item in order   │    │ with photos/SEO  │    │ with dimensions │
└────────┬────────┘    └────────┬─────────┘    └────────┬────────┘
         │                      │                        │
         └──────────────────────┴────────────────────────┘
                    Anti-Corruption Layer / Events
```

In Go, bounded contexts map to Go module subtrees or separate packages:

```
orders/domain/   ← Order context's view of a product (just ID + price)
catalog/domain/  ← Catalog context's full product (photos, SEO, reviews)
shipping/domain/ ← Shipping context's view (weight, dimensions)
```

When the Order context needs catalog data, it doesn't import `catalog/domain` directly — that would couple the contexts. Instead, it queries its own read model (maintained via events from the catalog context).

---

## 6. Ubiquitous Language

DDD's most important practice: **use the domain expert's exact words in your code**.

```go
// BAD: translation layer between domain and code
func updateOrderPaymentStatus(db *sql.DB, orderID int, isPaid bool) error

// GOOD: ubiquitous language — domain expert says "record payment"
func (o *Order) RecordPayment(payment Payment) error

// BAD: set status
order.status = "shipped"

// GOOD: past tense event, triggered by domain action
order.Ship(carrier, trackingNumber) // sets status, emits OrderShippedEvent
```

Maintain a glossary:
- **Order**: a request to purchase items, transitions from Draft → Pending → Confirmed → Shipped → Delivered
- **Customer**: a registered user who can place orders
- **Cart**: temporary container for items before checkout
- **Checkout**: the process of converting a Cart to a confirmed Order
- **Fulfillment**: the process of picking, packing, and shipping an Order

---

## Summary

- **Entity**: has identity that persists over time; use struct with typed ID
- **Value Object**: defined by its value; make it immutable (no setters, operations return new values)
- **Aggregate**: cluster of related objects; only the root is accessible from outside
- **Domain Events**: past-tense facts ("OrderSubmitted"); collected on the aggregate, dispatched after successful save
- **Bounded Context**: isolated sub-domain with its own model and language
- **Ubiquitous Language**: code vocabulary = domain expert vocabulary; the most important DDD practice

## Exercises

### Easy
1. Model a `BankAccount` aggregate: `Balance`, `AccountNumber`, `Status`. Add methods `Deposit(amount Money)` and `Withdraw(amount Money) error` — withdrawal fails if it would take balance below zero. Emit `MoneyDepositedEvent` and `MoneyWithdrawnEvent`.
2. Create a `PhoneNumber` value object. It must be E.164 format (`+1234567890`). Two `PhoneNumber` values are equal if their strings are equal.
3. Write unit tests for `Order.Submit()`: test that submitting a draft order succeeds and transitions to Pending, that submitting an empty order fails, and that submitting a non-draft order fails.

### Medium
4. Design an `Inventory` aggregate for the catalog bounded context. `Inventory` tracks quantity per `(ProductID, WarehouseID)` pair. Add `Reserve(qty int) error` and `Release(qty int)` methods. Emit `InventoryReservedEvent` when reservation succeeds.
5. Implement a **Saga** sketch: when `OrderSubmittedEvent` fires, the saga reserves inventory (`Inventory.Reserve`). If reservation fails, the saga emits `OrderCancelledEvent`. Wire this using the `InMemoryBus`.
6. Build a **domain model test** for the full order lifecycle: Create Draft → AddItems → Submit → Confirm → Ship → Deliver. Assert that each transition produces the correct domain event and that invalid transitions are rejected.

### Hard
7. Implement an **Anti-Corruption Layer** between the Order and Catalog bounded contexts: define `catalog.ProductView` in the order context (ID, name, price only). Build a `CatalogACL` that queries the catalog HTTP API and translates its response into `catalog.ProductView`. The order use case depends on this ACL interface, not on the catalog API directly.
8. Build a **DDD conformance checker** for your codebase: write a Go test that uses `go/ast` to assert that (a) no entity struct has public fields (must use methods), (b) no value object has any mutating methods (methods that take pointer receiver and modify state), (c) no repository is called from the domain layer.
