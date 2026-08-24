# Chapter 20: Context and Cancellation

Every real Go server must handle three signals: a client disconnects mid-request, a deadline expires, or a parent operation is cancelled. Go's `context.Context` package is the standard mechanism for carrying all three. Once you understand context, your Go code will be cancellation-aware throughout its entire call chain — from HTTP handler to database query to downstream API call.

## Table of Contents

1. [What Is Context](#1-what-is-context)
2. [Creating Contexts](#2-creating-contexts)
3. [Cancellation](#3-cancellation)
4. [Deadlines and Timeouts](#4-deadlines-and-timeouts)
5. [Context Values](#5-context-values)
6. [Propagating Context Through a Call Chain](#6-propagating-context-through-a-call-chain)
7. [Context in the Standard Library](#7-context-in-the-standard-library)
8. [Context Best Practices](#8-context-best-practices)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What Is Context

`context.Context` is an interface that carries:
1. **Cancellation signal** — when to stop
2. **Deadline/timeout** — the hard time limit
3. **Values** — request-scoped data (request ID, auth token)

```go
type Context interface {
    Done() <-chan struct{}           // Closed when context is cancelled/deadline reached
    Err() error                     // Why it was cancelled (nil if not yet cancelled)
    Deadline() (deadline time.Time, ok bool) // When it will expire (ok=false if no deadline)
    Value(key any) any              // Get a value stored in the context
}
```

The key insight: **context flows DOWN the call stack.** A parent context cancels all child contexts. If an HTTP request is cancelled, all goroutines doing work for that request should stop.

```
HTTP Request arrives
    └── ctx, cancel := context.WithTimeout(r.Context(), 5*time.Second)
           ├── fetchUser(ctx, userID)
           │      └── db.QueryRowContext(ctx, ...)  ← cancellation propagates here
           ├── fetchOrders(ctx, userID)
           │      └── http.NewRequestWithContext(ctx, ...) ← and here
           └── fetchRecommendations(ctx, userID)
                  └── client.DoWithContext(ctx, ...) ← and here
```

When the 5-second deadline expires, ALL of these calls receive the cancellation signal simultaneously.

### Quick Check
> 1. What three things does a `context.Context` carry?
> 2. How does cancellation propagate — upward or downward?
> 3. What does `ctx.Done()` return?

---

## 2. Creating Contexts

**`context.Background()`** — the root context, never cancelled:
```go
// Use as the starting point for all contexts:
ctx := context.Background()

// In main() or top-level goroutines:
ctx := context.Background()
result, err := doWork(ctx)
```

**`context.TODO()`** — placeholder when you're not sure yet:
```go
// Use when you know you should pass a context but haven't decided which:
func someFunc() {
    ctx := context.TODO()  // Signals "this needs a real context later"
    doSomething(ctx)
}
```

**`context.WithCancel`** — create a cancellable context:
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()  // ALWAYS call cancel — even if done naturally (frees resources)

go doWork(ctx)
go doMoreWork(ctx)

// Cancel when done:
cancel()  // All goroutines using ctx will see ctx.Done() close
```

**`context.WithTimeout`** — auto-cancels after a duration:
```go
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()  // Still call cancel to release resources early if done before timeout

result, err := slowOperation(ctx)
if err != nil {
    if errors.Is(err, context.DeadlineExceeded) {
        fmt.Println("operation timed out after 5s")
    }
}
```

**`context.WithDeadline`** — auto-cancels at a specific time:
```go
deadline := time.Now().Add(10 * time.Minute)
ctx, cancel := context.WithDeadline(context.Background(), deadline)
defer cancel()

result, err := processData(ctx)
```

**`context.WithValue`** — attach a value to the context:
```go
type contextKey string
const requestIDKey contextKey = "requestID"

ctx := context.WithValue(context.Background(), requestIDKey, "req-12345")
// Retrieve:
reqID := ctx.Value(requestIDKey).(string)
```

### Quick Check
> 1. When should you use `context.Background()` vs `context.TODO()`?
> 2. Why must you always call `cancel()` even if the work finishes successfully?
> 3. What is the difference between `WithTimeout` and `WithDeadline`?

---

## 3. Cancellation

**Checking for cancellation in your code:**
```go
func processItems(ctx context.Context, items []Item) error {
    for _, item := range items {
        // Check cancellation at each iteration:
        select {
        case <-ctx.Done():
            return ctx.Err()  // context.Canceled or context.DeadlineExceeded
        default:
            // Not cancelled — continue
        }
        
        if err := processItem(ctx, item); err != nil {
            return err
        }
    }
    return nil
}
```

**`ctx.Err()` after cancellation:**
```go
// ctx.Err() returns nil if not cancelled, otherwise:
context.Canceled         // Parent called cancel()
context.DeadlineExceeded // Deadline/timeout expired

if err := ctx.Err(); err != nil {
    switch {
    case errors.Is(err, context.Canceled):
        log.Println("request was cancelled by client")
    case errors.Is(err, context.DeadlineExceeded):
        log.Println("request timed out")
    }
}
```

**Cancelling a long-running loop:**
```go
func crawler(ctx context.Context, seed string) error {
    queue := []string{seed}
    visited := make(map[string]bool)
    
    for len(queue) > 0 {
        // Check for cancellation at the start of each iteration:
        if err := ctx.Err(); err != nil {
            return err
        }
        
        url := queue[0]
        queue = queue[1:]
        
        if visited[url] {
            continue
        }
        visited[url] = true
        
        links, err := fetch(ctx, url)  // fetch also respects context
        if err != nil {
            continue  // Don't stop on single fetch error
        }
        queue = append(queue, links...)
    }
    return nil
}
```

**Parent cancels children but not vice versa:**
```go
parent, parentCancel := context.WithCancel(context.Background())
child, childCancel := context.WithCancel(parent)
defer parentCancel()
defer childCancel()

// Cancelling child does NOT cancel parent:
childCancel()
fmt.Println(child.Err())   // context.Canceled
fmt.Println(parent.Err())  // nil

// Cancelling parent DOES cancel child:
parentCancel()
fmt.Println(parent.Err())  // context.Canceled
fmt.Println(child.Err())   // context.Canceled (already was, but also would be)
```

### Quick Check
> 1. What does `ctx.Err()` return when a deadline expires?
> 2. If you cancel a child context, does the parent context get cancelled?
> 3. How do you check for cancellation in a tight loop without blocking?

---

## 4. Deadlines and Timeouts

**Setting timeouts at different levels:**
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Overall request timeout: 10 seconds
    ctx, cancel := context.WithTimeout(r.Context(), 10*time.Second)
    defer cancel()
    
    // Database query timeout: 2 seconds (subset of overall)
    dbCtx, dbCancel := context.WithTimeout(ctx, 2*time.Second)
    defer dbCancel()
    user, err := db.GetUser(dbCtx, userID)
    if err != nil {
        // Was it DB timeout or overall timeout?
        if errors.Is(err, context.DeadlineExceeded) {
            http.Error(w, "service timeout", http.StatusGatewayTimeout)
        }
        return
    }
    
    // External API call: 5 seconds
    apiCtx, apiCancel := context.WithTimeout(ctx, 5*time.Second)
    defer apiCancel()
    recommendations, _ := externalAPI.Get(apiCtx, user.ID)
    
    respond(w, user, recommendations)
}
```

**Checking remaining deadline:**
```go
func callWithDeadlineCheck(ctx context.Context, minTime time.Duration) error {
    if deadline, ok := ctx.Deadline(); ok {
        remaining := time.Until(deadline)
        if remaining < minTime {
            return fmt.Errorf("not enough time: need %s, have %s", minTime, remaining)
        }
    }
    return doExpensiveWork(ctx)
}
```

### Quick Check
> 1. Can a child context have a longer deadline than its parent?
> 2. How do you check how much time is remaining on a context?
> 3. What error does a deadline-exceeded context return from `ctx.Err()`?

---

## 5. Context Values

Context values carry request-scoped data through your call chain:

```go
// Define typed keys to avoid collisions:
type contextKey string

const (
    requestIDKey contextKey = "requestID"
    userIDKey    contextKey = "userID"
    loggerKey    contextKey = "logger"
)

// Helper functions (idiomatic):
func WithRequestID(ctx context.Context, id string) context.Context {
    return context.WithValue(ctx, requestIDKey, id)
}

func RequestIDFromContext(ctx context.Context) (string, bool) {
    id, ok := ctx.Value(requestIDKey).(string)
    return id, ok
}
```

**Middleware adding values to context:**
```go
func requestIDMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        requestID := r.Header.Get("X-Request-ID")
        if requestID == "" {
            requestID = generateID()
        }
        
        ctx := WithRequestID(r.Context(), requestID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

func authMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        token := r.Header.Get("Authorization")
        userID, err := validateToken(token)
        if err != nil {
            http.Error(w, "unauthorized", 401)
            return
        }
        
        ctx := context.WithValue(r.Context(), userIDKey, userID)
        next.ServeHTTP(w, r.WithContext(ctx))
    })
}

// Handler uses values from context:
func profileHandler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    
    reqID, _ := RequestIDFromContext(ctx)
    userID, ok := ctx.Value(userIDKey).(int)
    if !ok {
        http.Error(w, "unauthorized", 401)
        return
    }
    
    log.Printf("[%s] fetching profile for user %d", reqID, userID)
    // ...
}
```

**Context value rules:**
1. Only use for **request-scoped** data, not function parameters (don't abuse it)
2. Keys must be **unexported custom types** — not strings — to avoid package collisions
3. Values should be **immutable** — don't store mutable state in context
4. Store only: request ID, auth info, trace IDs, deadline propagation data

### Quick Check
> 1. Why should context keys be custom types instead of plain strings?
> 2. What kind of data should you put in context values?
> 3. What is wrong with using `context.WithValue(ctx, "userID", id)`?

---

## 6. Propagating Context Through a Call Chain

The golden rule: **the first parameter of any function that does I/O should be `ctx context.Context`**.

```go
// Full stack example — context flows from HTTP handler to database:

// Layer 1: HTTP Handler
func getUserHandler(w http.ResponseWriter, r *http.Request) {
    id := chi.URLParam(r, "id")
    user, err := userService.GetUser(r.Context(), id)
    if err != nil {
        http.Error(w, err.Error(), 500)
        return
    }
    json.NewEncoder(w).Encode(user)
}

// Layer 2: Service
func (s *UserService) GetUser(ctx context.Context, id string) (*User, error) {
    // Add service-level timeout:
    ctx, cancel := context.WithTimeout(ctx, 3*time.Second)
    defer cancel()
    
    user, err := s.repo.FindByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("getting user %s: %w", id, err)
    }
    
    // Fetch profile data concurrently:
    g, ctx := errgroup.WithContext(ctx)
    var profile *Profile
    g.Go(func() error {
        var err error
        profile, err = s.profileSvc.Get(ctx, user.ProfileID)
        return err
    })
    
    if err := g.Wait(); err != nil {
        return nil, err
    }
    user.Profile = profile
    return user, nil
}

// Layer 3: Repository
func (r *UserRepository) FindByID(ctx context.Context, id string) (*User, error) {
    var user User
    err := r.db.QueryRowContext(ctx,
        "SELECT id, name, email FROM users WHERE id = $1", id,
    ).Scan(&user.ID, &user.Name, &user.Email)
    if err != nil {
        return nil, fmt.Errorf("query: %w", err)
    }
    return &user, nil
}
```

**Cancellation propagates automatically** — if the HTTP client disconnects, `r.Context()` is cancelled, which cancels the service timeout context, which cancels the DB query. The database driver sees the cancelled context and aborts the query immediately.

### Quick Check
> 1. Where should `ctx context.Context` appear in a function signature?
> 2. What happens to a database query when the context is cancelled?
> 3. How does `errgroup.WithContext` help with concurrent operations?

---

## 7. Context in the Standard Library

Most standard library functions that do I/O have context-aware variants:

```go
// database/sql:
db.QueryContext(ctx, query, args...)
db.QueryRowContext(ctx, query, args...)
db.ExecContext(ctx, query, args...)
tx.QueryContext(ctx, query, args...)

// net/http — client:
req, _ := http.NewRequestWithContext(ctx, "GET", url, nil)
resp, err := client.Do(req)

// net/http — server (r.Context() is automatically cancelled when client disconnects):
func handler(w http.ResponseWriter, r *http.Request) {
    ctx := r.Context()
    // ctx is cancelled when:
    // 1. Client closes connection
    // 2. Request handler returns
    // 3. Server is shutting down
}

// os/exec:
cmd := exec.CommandContext(ctx, "sleep", "10")
cmd.Run()  // Kills the process if ctx is cancelled

// grpc:
conn, _ := grpc.Dial(addr, grpc.WithInsecure())
client := pb.NewGreeterClient(conn)
reply, err := client.SayHello(ctx, &pb.HelloRequest{Name: "Alice"})

// redis (go-redis):
rdb.Get(ctx, "key")
rdb.Set(ctx, "key", value, 0)
```

### Quick Check
> 1. What is `http.NewRequestWithContext` used for?
> 2. When is `r.Context()` automatically cancelled in an HTTP handler?
> 3. What happens to an `exec.CommandContext` command if its context is cancelled?

---

## 8. Context Best Practices

```go
// 1. Always pass ctx as the FIRST parameter:
func process(ctx context.Context, data Data) error { ... }  // Correct
func process(data Data, ctx context.Context) error { ... }  // Wrong

// 2. Never store context in a struct:
type Service struct {
    ctx context.Context  // Wrong! Context belongs to a request, not a service
    db  *sql.DB
}

// Correct: pass ctx to each method
type Service struct {
    db *sql.DB
}
func (s *Service) Process(ctx context.Context, data Data) error { ... }

// 3. Never pass nil context — use context.Background() or context.TODO():
func bad(ctx context.Context) {}
bad(nil)  // Causes nil pointer panics in any code that uses ctx

// 4. Always call cancel() — wrap it in defer immediately:
ctx, cancel := context.WithTimeout(ctx, 5*time.Second)
defer cancel()  // Do this immediately, before any other code

// 5. Don't use context for optional parameters:
// Wrong:
ctx = context.WithValue(ctx, "debug", true)  // Use function parameter instead
// Right:
func doWork(ctx context.Context, debug bool) {}

// 6. Log the context's request ID in errors:
func (s *Service) Process(ctx context.Context, id string) error {
    reqID, _ := RequestIDFromContext(ctx)
    log.Printf("[%s] processing %s", reqID, id)
    
    if err := s.db.Exec(ctx, ...); err != nil {
        return fmt.Errorf("[%s] db exec: %w", reqID, err)
    }
    return nil
}
```

**Detecting cancellation quickly:**
```go
// For computationally intensive loops:
func heavyCompute(ctx context.Context, data []int) (int, error) {
    result := 0
    for i, v := range data {
        if i%1000 == 0 {  // Check every 1000 iterations, not every single one
            select {
            case <-ctx.Done():
                return 0, ctx.Err()
            default:
            }
        }
        result += compute(v)
    }
    return result, nil
}
```

### Quick Check
> 1. Should you store `context.Context` in a struct?
> 2. What should you pass if you don't have a context yet?
> 3. How often should you check for context cancellation in a heavy loop?

---

## Summary

- **`context.Context`**: carries cancellation signal, deadline, and request values
- **`context.Background()`**: root context — use in `main`, tests, top-level goroutines
- **`context.TODO()`**: placeholder — marks "context needed here later"
- **`WithCancel`**: manual cancellation; **`WithTimeout`**: duration; **`WithDeadline`**: absolute time
- **Always `defer cancel()`** — even if the operation finishes before the deadline
- **`ctx.Done()`**: channel closed when cancelled; **`ctx.Err()`**: `Canceled` or `DeadlineExceeded`
- **Parent cancels children** — but cancelling a child does NOT affect the parent
- **Context values**: use typed keys, only for request-scoped data (IDs, auth)
- **Propagate `ctx` everywhere** — first parameter of any function doing I/O
- **Don't store context in structs** — it belongs to a request/operation, not a type
- **Standard library**: `QueryContext`, `NewRequestWithContext`, `CommandContext`, etc.

---

## Exercises

### Easy
1. Write a `fetchWithTimeout(url string, timeout time.Duration) (string, error)` function that uses `context.WithTimeout` and `http.NewRequestWithContext`. Return `context.DeadlineExceeded` (wrapped) if the timeout expires.
2. Build a long-running loop that processes integers 1 to 1,000,000. Accept a `ctx context.Context`. Check for cancellation every 1,000 iterations. Cancel it from main after 100ms and verify the loop exits cleanly with the correct error.
3. Create a middleware that adds a unique request ID (UUID or random hex string) to every request's context. Log the request ID at the start and end of each request. Verify it's available in handlers via `ctx.Value`.

### Medium
4. Cascading timeouts: Write a service with three layers — Handler → Service → Repository. Each layer applies its own timeout (Handler: 10s, Service: 5s, Repository: 2s). In the repository, simulate a slow query with `time.Sleep`. Verify that the repository respects its context, and that when the service timeout fires it cancels the repository operation. Write a test that verifies each timeout level works independently.
5. Context propagation in errgroup: Write a function `fetchDashboard(ctx context.Context, userID int) (*Dashboard, error)` that concurrently fetches user, orders, and recommendations using `errgroup.WithContext`. If any fetch fails, cancel the others. Add timeouts to each individual fetch. Test: verify that if one fetch fails, the others are cancelled; verify the overall timeout works.
6. Request tracing: Build a tracing system using context values. Each incoming request gets a `TraceID` and `SpanID`. When a function calls a sub-function, it creates a child span (new SpanID, same TraceID). Implement `StartSpan(ctx, name) (context.Context, func())` that stores span data and returns a finish function. The finish function logs `"[traceID] spanName: Xms"`. Chain 3 spans and verify the trace output.

### Hard
7. Graceful shutdown: Build an HTTP server that handles graceful shutdown. When `SIGTERM` is received: (1) stop accepting new connections, (2) wait for in-flight requests to finish (up to 30 seconds), (3) cancel any requests that haven't finished within the grace period. Use `context.WithCancel` for the server's lifetime context and `server.Shutdown(ctx)` for graceful shutdown. Test by sending 10 slow requests (each taking 5s), then sending SIGTERM — verify requests in progress complete but no new requests are accepted.
8. Distributed tracing: Implement W3C TraceContext propagation (https://www.w3.org/TR/trace-context/). `InjectHeaders(ctx, headers http.Header)` — write trace-parent header from ctx. `ExtractHeaders(ctx, headers http.Header) context.Context` — read trace-parent header and store in ctx. `NewRootSpan(ctx) context.Context` — create a new root trace. `NewChildSpan(ctx) context.Context` — create a child span from parent in ctx. Test by simulating a 3-service call chain (A → B → C), injecting and extracting headers between each, and verifying the trace IDs propagate correctly.
