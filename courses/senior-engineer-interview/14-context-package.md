# Chapter 14: The context Package — Cancellation, Timeouts & Values

The `context` package is one of the most important packages in the Go standard library for building production services. Every senior Go engineer must understand it deeply — not just how to pass it around, but why it exists, what guarantees it provides, and what mistakes to avoid.

## Table of Contents

1. [Why context Exists](#1-why-context-exists)
2. [The context.Context Interface](#2-the-contextcontext-interface)
3. [Creating Contexts](#3-creating-contexts)
4. [Cancellation](#4-cancellation)
5. [Timeouts and Deadlines](#5-timeouts-and-deadlines)
6. [Context Values — The Right and Wrong Uses](#6-context-values--the-right-and-wrong-uses)
7. [Threading Context Through a Service](#7-threading-context-through-a-service)
8. [Common Mistakes](#8-common-mistakes)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. Why context Exists

Before `context`, cancelling in-flight operations was a manual, error-prone affair. Each library invented its own stop channel. There was no standard way to propagate cancellation through a call chain.

`context` solves this: it provides a standard mechanism to carry deadlines, cancellation signals, and request-scoped values across API boundaries and between goroutines in a call tree.

**The concrete problem it solves:** A user sends an HTTP request. Your handler calls service A, which calls service B, which queries a database. The user cancels their request (closes the browser). Without context, all that work continues pointlessly. With context, cancellation propagates automatically down the entire call tree.

---

## 2. The context.Context Interface

```go
type Context interface {
    // Deadline returns the time when the context will be cancelled (if set)
    Deadline() (deadline time.Time, ok bool)

    // Done returns a channel that's closed when the context is cancelled
    Done() <-chan struct{}

    // Err returns the reason the context was cancelled:
    // context.Canceled if cancelled explicitly
    // context.DeadlineExceeded if deadline passed
    Err() error

    // Value returns the value associated with key, or nil
    Value(key interface{}) interface{}
}
```

The interface is read-only. You cannot cancel a context you received — only the creator can cancel it (via the CancelFunc they hold).

---

## 3. Creating Contexts

```go
// context.Background(): the root context. Use at the top of main(), tests, or request handlers.
ctx := context.Background()

// context.TODO(): placeholder when you're not sure which context to use yet.
// Signals work in progress — should not appear in production code.
ctx := context.TODO()

// Create a cancellable child context
ctx, cancel := context.WithCancel(parent)
defer cancel() // ALWAYS call cancel to release resources, even if context is cancelled first

// Create a context that auto-cancels after a duration
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel()

// Create a context that cancels at a specific time
deadline := time.Now().Add(10 * time.Second)
ctx, cancel := context.WithDeadline(parent, deadline)
defer cancel()

// Create a context with a value
ctx = context.WithValue(parent, keyType("userID"), "user123")
```

---

## 4. Cancellation

### The Cancel Function

```go
func doWork(ctx context.Context) error {
    // Check if context is already cancelled before doing work
    select {
    case <-ctx.Done():
        return ctx.Err() // context.Canceled or context.DeadlineExceeded
    default:
    }

    // Do the work
    result, err := fetchData()
    if err != nil { return err }

    // Check again after potentially slow operation
    select {
    case <-ctx.Done():
        return ctx.Err()
    default:
    }

    return process(result)
}
```

### Propagating Cancellation Down the Call Tree

```go
// Top-level handler creates the context
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context() // HTTP handlers get context from the request
    
    // If user disconnects, r.Context() is cancelled automatically
    result, err := businessLogic(ctx)
    // ...
}

func businessLogic(ctx context.Context) (Result, error) {
    // Pass ctx down — each function checks it
    data, err := fetchFromDB(ctx)
    if err != nil { return Result{}, err }
    
    return processData(ctx, data)
}

func fetchFromDB(ctx context.Context) (Data, error) {
    // Pass ctx to database driver — it respects cancellation
    rows, err := db.QueryContext(ctx, "SELECT ...")
    if err != nil { return Data{}, err }
    // ...
}
```

### Cancelling a Goroutine

```go
func startWorker(ctx context.Context) {
    go func() {
        for {
            select {
            case <-ctx.Done():
                fmt.Println("worker stopped:", ctx.Err())
                return
            default:
                doWork()
            }
        }
    }()
}

// Cancel the worker from the parent:
ctx, cancel := context.WithCancel(context.Background())
startWorker(ctx)
// ...
cancel() // stop the worker
```

---

## 5. Timeouts and Deadlines

```go
// WithTimeout: relative duration from now
func callExternalService(data []byte) (Response, error) {
    ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
    defer cancel() // critical: releases timer resources

    req, _ := http.NewRequestWithContext(ctx, "POST", serviceURL, bytes.NewReader(data))
    
    resp, err := http.DefaultClient.Do(req)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            return Response{}, fmt.Errorf("service call timed out after 3s")
        }
        return Response{}, err
    }
    defer resp.Body.Close()
    // ...
}
```

### Checking Remaining Time

```go
func longOperation(ctx context.Context) error {
    // Check how much time is left before we start something expensive
    deadline, ok := ctx.Deadline()
    if ok {
        timeLeft := time.Until(deadline)
        if timeLeft < 500*time.Millisecond {
            return fmt.Errorf("not enough time to complete operation: %v remaining", timeLeft)
        }
    }
    // proceed with operation
    return nil
}
```

### Timeout vs Deadline

```go
// WithTimeout: relative — "cancel in 5 seconds from now"
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)

// WithDeadline: absolute — "cancel at 15:30:00"
deadline := time.Date(2024, 1, 1, 15, 30, 0, 0, time.UTC)
ctx, cancel := context.WithDeadline(ctx, deadline)

// Use WithTimeout for API call timeouts.
// Use WithDeadline for SLA guarantees: "this entire request must complete by time X."
```

---

## 6. Context Values — The Right and Wrong Uses

Context values are for request-scoped data that crosses API boundaries: request IDs, auth tokens, trace IDs. They are NOT for optional parameters.

```go
// CORRECT uses of context values:
type contextKey string
const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
    traceIDKey   contextKey = "traceID"
)

func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}

// Middleware sets request-scoped values:
func loggingMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ctx := context.WithValue(r.Context(), requestIDKey, generateID())
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}
```

### Why Context Values Should Use Custom Key Types

```go
// BAD: string keys can collide across packages
ctx = context.WithValue(ctx, "userID", "123")
// Two packages both use "userID" — one overwrites the other!

// GOOD: unexported custom key type prevents collisions
type authKey struct{} // unexported empty struct

ctx = context.WithValue(ctx, authKey{}, user)
// authKey{} from package A ≠ authKey{} from package B (different types)
```

### Wrong Uses of Context Values

```go
// WRONG: using context for optional function parameters
func processOrder(ctx context.Context, orderID string) {
    discount := ctx.Value("discount").(float64) // this is a function parameter!
    // ...
}
// CORRECT: pass as explicit parameter
func processOrder(ctx context.Context, orderID string, discount float64) {}

// WRONG: storing mutable state in context
func handler(ctx context.Context) {
    counter := 0
    ctx = context.WithValue(ctx, "counter", &counter) // counter can be mutated
    // This is a hidden mutable dependency — very bad for testing
}
```

---

## 7. Threading Context Through a Service

The first parameter of every function that does I/O, calls another service, or runs for a significant duration should be `ctx context.Context`.

```go
// The pattern: ctx is ALWAYS the first argument
func (s *OrderService) CreateOrder(ctx context.Context, req CreateOrderRequest) (*Order, error) {
    // Validate
    if err := s.validator.Validate(ctx, req); err != nil {
        return nil, err
    }
    
    // Check inventory
    available, err := s.inventoryClient.CheckStock(ctx, req.ProductID, req.Quantity)
    if err != nil {
        return nil, fmt.Errorf("checking stock: %w", err)
    }
    if !available {
        return nil, ErrOutOfStock
    }
    
    // Create in database
    order, err := s.repo.Create(ctx, req)
    if err != nil {
        return nil, fmt.Errorf("creating order: %w", err)
    }
    
    // Publish event
    if err := s.eventBus.Publish(ctx, OrderCreatedEvent{OrderID: order.ID}); err != nil {
        // Log but don't fail — order was created successfully
        s.log.Error("failed to publish order event", "err", err)
    }
    
    return order, nil
}
```

### Context in Tests

```go
// For tests, use context.Background() unless you're specifically testing timeout behavior
func TestCreateOrder(t *testing.T) {
    ctx := context.Background()
    result, err := service.CreateOrder(ctx, req)
    // ...
}

// For testing timeout handling:
func TestCreateOrderTimeout(t *testing.T) {
    ctx, cancel := context.WithTimeout(context.Background(), 1*time.Millisecond)
    defer cancel()
    
    _, err := service.CreateOrder(ctx, req)
    if !errors.Is(err, context.DeadlineExceeded) {
        t.Fatalf("expected deadline exceeded, got: %v", err)
    }
}
```

---

## 8. Common Mistakes

### Mistake 1: Not Calling Cancel

```go
// WRONG: timer goroutine leaks if you never call cancel
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
// ... forgot defer cancel()
result, _ := doWork(ctx)

// CORRECT:
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel() // even if doWork finishes before timeout, this cleans up
result, _ := doWork(ctx)
```

### Mistake 2: Storing context in a struct

```go
// WRONG: context should not be stored in a struct
type Service struct {
    ctx context.Context // WRONG
}

// CORRECT: pass context as a function parameter
func (s *Service) DoWork(ctx context.Context) error {}
```

### Mistake 3: Passing nil context

```go
// WRONG: panics in many stdlib functions
doWork(nil)

// CORRECT: if you're not sure what context to use, use context.TODO()
doWork(context.TODO())
```

### Mistake 4: Creating context from the wrong parent

```go
// WRONG: creating context from context.Background() inside a handler
// loses the parent's deadline/cancellation
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // This creates a NEW root context — doesn't inherit request deadline!
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    doWork(ctx)
}

// CORRECT: derive from request context
func handleRequest(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
    defer cancel()
    doWork(ctx)
}
```

---

## 9. Interview Questions & Model Answers

**Q: What is the purpose of the context package?**

"Context provides a standard way to carry cancellation signals, deadlines, and request-scoped values across API boundaries. Before context, libraries invented their own cancellation mechanisms and there was no standard way to propagate them. Context solves two main problems: cancellation (if the caller is done, stop the work) and timeout propagation (the request must complete within X seconds, regardless of how many nested calls it makes)."

**Q: What's the difference between context.WithTimeout and context.WithDeadline?**

"WithTimeout takes a duration — 'cancel in 5 seconds from now.' WithDeadline takes an absolute time — 'cancel at 15:30:00.' Internally, WithTimeout calls WithDeadline with time.Now().Add(timeout). Use WithTimeout for API call budgets. Use WithDeadline when you have a hard SLA — 'this request, which started at 15:29:55, must complete by 15:30:00 — that's only 5 seconds total, regardless of when suboperations start.'"

**Q: Why should context values use custom key types instead of strings?**

"String keys are global. Two different packages using 'userID' as a key would collide — the second Write overwrites the first. Custom unexported types prevent this: `type myKey struct{}` in package A is a different type than `type myKey struct{}` in package B, so their keys never collide. The type system enforces namespacing."

**Q: When should you NOT propagate a context?**

"When you're starting background work that should outlive the request. For example, if a request triggers a long-running background job, that job should not use the request's context — when the request completes (and its context is cancelled), the job should continue. Instead, create a new context with context.Background() for the background job, possibly with its own timeout."

---

## Summary

- `context.Context` is an immutable tree of request-scoped data: deadline, cancellation, and values.
- Root contexts: `Background()` for real code, `TODO()` for work in progress.
- Always `defer cancel()` after `WithCancel`, `WithTimeout`, or `WithDeadline` — prevents goroutine leaks.
- Cancellation propagates automatically: cancelling a parent cancels all children.
- Check `ctx.Done()` in long loops and before expensive operations.
- Context values: for request-scoped metadata (trace ID, user ID) crossing API boundaries. Use custom unexported key types to prevent collisions.
- Pass ctx as the first parameter of every function that does I/O or runs for a significant time.
- Never store context in a struct; never pass nil context; always derive from the appropriate parent.
