# Chapter 39: SOLID Principles & Design Patterns in Go

Low-Level Design (LLD) interviews test whether you write code that can be maintained, extended, and tested by other senior engineers. SOLID principles and design patterns are the vocabulary for these conversations.

## Table of Contents

1. [SOLID Principles in Go](#1-solid-principles-in-go)
2. [Creational Patterns](#2-creational-patterns)
3. [Structural Patterns](#3-structural-patterns)
4. [Behavioral Patterns](#4-behavioral-patterns)
5. [Patterns Used in Real Systems](#5-patterns-used-in-real-systems)
6. [Interview Questions & Model Answers](#6-interview-questions--model-answers)
7. [Summary](#summary)

---

## 1. SOLID Principles in Go

### S — Single Responsibility Principle

One struct/function = one reason to change.

```go
// BAD: UserService does too many things
type UserService struct{}
func (s *UserService) CreateUser(name, email string) error { /* insert into DB */ }
func (s *UserService) SendWelcomeEmail(email string) error { /* call SMTP server */ }
func (s *UserService) GenerateReport(userID int) error { /* generate PDF */ }

// GOOD: separate each responsibility
type UserService struct{ repo UserRepository }
func (s *UserService) CreateUser(name, email string) (*User, error) { /* DB only */ }

type EmailService struct{ smtp SMTPClient }
func (s *EmailService) SendWelcomeEmail(email string) error { /* SMTP only */ }

type ReportService struct{ userRepo UserRepository }
func (s *ReportService) GenerateUserReport(userID int) ([]byte, error) { /* PDF only */ }

// Now each service can change independently. Email provider changes: only EmailService changes.
```

### O — Open/Closed Principle

Open for extension, closed for modification. Add new behavior without changing existing code.

```go
// BAD: add new payment method = modify existing code
func processPayment(method string, amount float64) error {
    switch method {
    case "credit_card":
        // ...
    case "paypal":
        // ...
    // Adding "crypto" requires modifying this function
    }
}

// GOOD: define an interface; new methods implement it without touching existing code
type PaymentProcessor interface {
    Charge(amount float64, currency string) error
    Refund(transactionID string, amount float64) error
}

type CreditCardProcessor struct { /* Stripe credentials */ }
func (c *CreditCardProcessor) Charge(amount float64, currency string) error { /* ... */ }
func (c *CreditCardProcessor) Refund(txID string, amount float64) error { /* ... */ }

type PayPalProcessor struct { /* PayPal credentials */ }
func (p *PayPalProcessor) Charge(amount float64, currency string) error { /* ... */ }
func (p *PayPalProcessor) Refund(txID string, amount float64) error { /* ... */ }

// Adding CryptoProcessor: create new struct, implement interface. Zero existing code changes.

type PaymentService struct{ processor PaymentProcessor }
func (s *PaymentService) Pay(amount float64) error {
    return s.processor.Charge(amount, "USD")
}
```

### L — Liskov Substitution Principle

Subtypes must be substitutable for their base types without altering correctness.

```go
// In Go: any type implementing the interface is a valid substitute

type Storage interface {
    Get(key string) ([]byte, error)
    Set(key string, value []byte) error
}

// Both implement Storage — callers don't care which is used
type MemoryStorage struct{ data map[string][]byte }
type RedisStorage struct{ client *redis.Client }
type S3Storage struct{ bucket string }

// LSP violation (if it were possible in Go):
// A ReadOnlyStorage that panics on Set() would break any code expecting Storage
// Don't return errors or panic for operations the interface promises to support
```

### I — Interface Segregation Principle

Many specific interfaces are better than one large interface.

```go
// BAD: large interface forces unnecessary methods
type Worker interface {
    GetTask() Task
    ProcessTask(t Task) error
    GenerateReport() string
    NotifyCompletion() error
}

// GOOD: small, focused interfaces
type TaskGetter interface { GetTask() Task }
type TaskProcessor interface { ProcessTask(t Task) error }
type Reporter interface { GenerateReport() string }
type Notifier interface { NotifyCompletion() error }

// Compose as needed:
type TaskWorker interface {
    TaskGetter
    TaskProcessor
}

// Functions only ask for what they need:
func processAllTasks(getter TaskGetter, processor TaskProcessor) error { /* ... */ }
func sendReport(r Reporter, n Notifier) error { /* ... */ }
```

### D — Dependency Inversion Principle

Depend on abstractions, not concretions.

```go
// BAD: high-level module depends on concrete low-level module
type OrderService struct {
    db *PostgresDB  // concrete dependency
}

// GOOD: depend on interface (abstraction)
type OrderRepository interface {
    GetOrder(ctx context.Context, id string) (*Order, error)
    SaveOrder(ctx context.Context, order *Order) error
}

type OrderService struct {
    repo OrderRepository  // abstraction
}

// In tests: use a mock that implements OrderRepository
// In production: use *PostgresOrderRepository
// OrderService doesn't care which — it depends on the abstraction
```

---

## 2. Creational Patterns

### Builder Pattern (for complex object construction)

```go
// Used when an object has many optional fields
type HTTPClientBuilder struct {
    timeout       time.Duration
    maxRetries    int
    baseURL       string
    headers       map[string]string
    transport     http.RoundTripper
}

func NewHTTPClientBuilder(baseURL string) *HTTPClientBuilder {
    return &HTTPClientBuilder{
        timeout:    30 * time.Second,
        maxRetries: 3,
        baseURL:    baseURL,
        headers:    make(map[string]string),
    }
}

func (b *HTTPClientBuilder) WithTimeout(d time.Duration) *HTTPClientBuilder {
    b.timeout = d
    return b
}

func (b *HTTPClientBuilder) WithHeader(key, value string) *HTTPClientBuilder {
    b.headers[key] = value
    return b
}

func (b *HTTPClientBuilder) Build() *HTTPClient {
    return &HTTPClient{/* ... */}
}

// Usage:
client := NewHTTPClientBuilder("https://api.example.com").
    WithTimeout(10 * time.Second).
    WithHeader("Authorization", "Bearer token").
    Build()
```

### Factory Pattern

```go
type NotificationSender interface {
    Send(ctx context.Context, to, message string) error
}

// Factory function returns the right implementation based on type
func NewNotificationSender(typ string, cfg Config) NotificationSender {
    switch typ {
    case "email":
        return &EmailSender{smtp: cfg.SMTPConfig}
    case "sms":
        return &SMSSender{twilioKey: cfg.TwilioKey}
    case "push":
        return &PushSender{fcmKey: cfg.FCMKey}
    default:
        panic(fmt.Sprintf("unknown sender type: %s", typ))
    }
}
```

### Singleton (with sync.Once)

```go
// Singleton: only one instance across the application
var (
    dbInstance *sql.DB
    dbOnce     sync.Once
)

func GetDB() *sql.DB {
    dbOnce.Do(func() {
        // This runs exactly once, even with concurrent callers
        db, err := sql.Open("postgres", os.Getenv("DATABASE_URL"))
        if err != nil { panic(err) }
        dbInstance = db
    })
    return dbInstance
}
```

---

## 3. Structural Patterns

### Adapter Pattern

```go
// Adapt an existing type to match a required interface
type OldPaymentGateway struct{}
func (g *OldPaymentGateway) MakePayment(amount int, cardNum string) bool { /* legacy */ }

// New interface our system expects:
type PaymentProcessor interface {
    Charge(amount float64, currency string) error
}

// Adapter wraps the old gateway and implements the new interface:
type OldGatewayAdapter struct{ gateway *OldPaymentGateway }

func (a *OldGatewayAdapter) Charge(amount float64, currency string) error {
    cents := int(amount * 100)
    success := a.gateway.MakePayment(cents, "stored")
    if !success {
        return errors.New("payment failed")
    }
    return nil
}

// Now OldPaymentGateway can be used wherever PaymentProcessor is expected
```

### Decorator Pattern (Middleware in Go)

```go
// Add behavior to a function/interface without modifying it
type Handler func(ctx context.Context, req Request) (Response, error)

// Decorator adds logging:
func WithLogging(h Handler) Handler {
    return func(ctx context.Context, req Request) (Response, error) {
        start := time.Now()
        resp, err := h(ctx, req)
        log.Printf("request took %v, err=%v", time.Since(start), err)
        return resp, err
    }
}

// Decorator adds retries:
func WithRetry(h Handler, maxRetries int) Handler {
    return func(ctx context.Context, req Request) (Response, error) {
        for i := 0; i < maxRetries; i++ {
            resp, err := h(ctx, req)
            if err == nil { return resp, nil }
        }
        return Response{}, errors.New("max retries exceeded")
    }
}

// Chain decorators:
handler := WithLogging(WithRetry(baseHandler, 3))
```

---

## 4. Behavioral Patterns

### Strategy Pattern

```go
// Interchangeable algorithms
type SortStrategy interface {
    Sort(data []int) []int
}

type QuickSort struct{}
func (q *QuickSort) Sort(data []int) []int { /* quicksort */ }

type MergeSort struct{}
func (m *MergeSort) Sort(data []int) []int { /* mergesort */ }

type Sorter struct{ strategy SortStrategy }
func (s *Sorter) Sort(data []int) []int { return s.strategy.Sort(data) }

// Change strategy at runtime:
sorter := &Sorter{strategy: &QuickSort{}}
if len(data) < 10 {
    sorter.strategy = &InsertionSort{} // better for small arrays
}
```

### Observer Pattern (Event System)

```go
// Observers are notified when state changes
type EventType string

type Event struct {
    Type    EventType
    Payload any
}

type Handler func(e Event)

type EventBus struct {
    mu       sync.RWMutex
    handlers map[EventType][]Handler
}

func (b *EventBus) Subscribe(eventType EventType, h Handler) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.handlers[eventType] = append(b.handlers[eventType], h)
}

func (b *EventBus) Publish(e Event) {
    b.mu.RLock()
    handlers := b.handlers[e.Type]
    b.mu.RUnlock()
    
    for _, h := range handlers {
        go h(e) // notify asynchronously
    }
}

// Usage:
bus := &EventBus{handlers: make(map[EventType][]Handler)}

bus.Subscribe("order.created", func(e Event) {
    order := e.Payload.(*Order)
    emailService.SendConfirmation(order)
})

bus.Subscribe("order.created", func(e Event) {
    order := e.Payload.(*Order)
    inventoryService.ReserveItems(order)
})

bus.Publish(Event{Type: "order.created", Payload: newOrder})
```

---

## 5. Patterns Used in Real Systems

```
Functional Options (ubiquitous in Go standard library and popular libraries):
  func NewServer(opts ...Option) *Server
  Used for: optional configuration without breaking callers

Pipeline (data transformation):
  stage1 → stage2 → stage3
  Used for: ETL pipelines, request processing

Command Pattern (encapsulate requests as objects):
  Used for: undo/redo, task queues, request logging

Repository Pattern (data access abstraction):
  type UserRepository interface { Find, Save, Delete }
  Separates business logic from data access — enables testing with mocks
```

---

## 6. Interview Questions & Model Answers

**Q: What is the Open/Closed Principle and how does it apply in Go?**

"The Open/Closed Principle says code should be open for extension but closed for modification. In Go, this is achieved through interfaces. When you need to add new behavior, you create a new type that implements an existing interface rather than modifying a switch statement. For example, a PaymentProcessor interface with Charge() and Refund() methods. When adding a new payment method, I implement the interface in a new struct — zero existing code changes. The calling code doesn't change because it depends on the interface, not the concrete implementation."

**Q: Explain Dependency Injection and why it matters for testing.**

"Dependency Injection means passing dependencies (like a database connection or HTTP client) into a struct rather than creating them inside. Instead of `NewOrderService()` that creates its own `*sql.DB`, I write `NewOrderService(repo OrderRepository)` where OrderRepository is an interface. In production, I pass `*PostgresOrderRepository`. In tests, I pass a mock that implements the same interface but returns predefined data. This makes tests fast (no real database), isolated (tests don't fail due to external dependencies), and reliable (deterministic behavior). It's one of the most impactful practices for maintaining a testable codebase."

---

## Summary

- **S** — Single Responsibility: one reason to change. Separate DB access from email from business logic.
- **O** — Open/Closed: interfaces over switch statements. New behavior = new implementation, not modified code.
- **L** — Liskov Substitution: any implementation can replace another. Don't violate interface contracts.
- **I** — Interface Segregation: small focused interfaces. Callers only depend on what they use.
- **D** — Dependency Inversion: depend on abstractions. Inject dependencies via interfaces for testability.
- **Builder:** complex objects with optional fields.
- **Factory:** create the right implementation based on configuration.
- **Adapter:** wrap legacy code to implement a modern interface.
- **Decorator (middleware):** add cross-cutting concerns (logging, retry, auth) without modifying handlers.
- **Observer (event bus):** decouple publishers from subscribers.
- **Strategy:** interchangeable algorithms selected at runtime.
