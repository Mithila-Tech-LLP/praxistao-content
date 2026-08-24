# Chapter 14: Error Handling

Go's approach to errors is refreshingly direct: errors are values, not exceptions. There's no `try/catch` block that changes control flow invisibly. Every function that can fail returns an error as its last return value, and the caller decides what to do with it. This explicitness is Go's most defining characteristic — and once you embrace it, you'll appreciate how much clearer code becomes.

## Table of Contents

1. [The error Interface](#1-the-error-interface)
2. [Creating Errors](#2-creating-errors)
3. [Handling Errors](#3-handling-errors)
4. [Error Wrapping and Unwrapping](#4-error-wrapping-and-unwrapping)
5. [Custom Error Types](#5-custom-error-types)
6. [Sentinel Errors](#6-sentinel-errors)
7. [Error Handling Patterns](#7-error-handling-patterns)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The error Interface

`error` is a built-in interface:

```go
type error interface {
    Error() string
}
```

That's it. Any type that implements `Error() string` is an error. This simplicity means error handling integrates with Go's type system naturally.

```go
// Functions that can fail return (result, error):
func ReadFile(name string) ([]byte, error) { ... }
func ParseJSON(data []byte, v any) error   { ... }
func Connect(addr string) (*Conn, error)   { ... }

// Functions that always succeed return just the value:
func Add(a, b int) int { ... }
func ToUpper(s string) string { ... }
```

The convention: **last return value is error, nil means success**.

```go
content, err := os.ReadFile("config.json")
if err != nil {
    // Something went wrong
    log.Fatal(err)
}
// err is nil here — content is valid
```

### Quick Check
> 1. What is the `error` interface's only method?
> 2. What does `nil` mean when returned as an error?
> 3. Where does the error appear in a function's return values?

---

## 2. Creating Errors

**`errors.New`** — for simple static messages:
```go
import "errors"

var ErrNotFound = errors.New("not found")

func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, errors.New("invalid user id")
    }
    // ...
    return nil, errors.New("user not found")
}
```

**`fmt.Errorf`** — for formatted messages with context:
```go
import "fmt"

func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, fmt.Errorf("invalid user id: %d", id)
    }
    
    user, err := db.Query(id)
    if err != nil {
        return nil, fmt.Errorf("querying user %d: %w", id, err)  // %w wraps the error
    }
    return user, nil
}
```

The `%w` verb (Go 1.13+) wraps an existing error — preserving the original while adding context. This is how you build error chains.

**Implementing the `error` interface directly:**
```go
type ValidationError struct {
    Field   string
    Message string
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %q: %s", e.Field, e.Message)
}

func validateAge(age int) error {
    if age < 0 {
        return &ValidationError{Field: "age", Message: "must be non-negative"}
    }
    if age > 150 {
        return &ValidationError{Field: "age", Message: "unreasonably large"}
    }
    return nil
}
```

### Quick Check
> 1. What is the difference between `errors.New` and `fmt.Errorf`?
> 2. What does `%w` do in `fmt.Errorf`?
> 3. How do you make a custom type satisfy the `error` interface?

---

## 3. Handling Errors

**The standard pattern — check immediately:**
```go
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}
```

**Never ignore errors (except when truly intentional):**
```go
// Bad — silently swallows errors:
data, _ := os.ReadFile("config.json")  // _ discards the error
process(data)                           // May crash or corrupt data

// Good — handle or explicitly propagate:
data, err := os.ReadFile("config.json")
if err != nil {
    return fmt.Errorf("reading config: %w", err)
}

// Intentional ignore (add a comment to document why):
_ = os.Remove(tempFile)  // Best-effort cleanup; not critical if it fails
```

**Multiple operations — check each:**
```go
func processFile(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return fmt.Errorf("opening file: %w", err)
    }
    defer f.Close()
    
    data, err := io.ReadAll(f)
    if err != nil {
        return fmt.Errorf("reading file: %w", err)
    }
    
    result, err := parseData(data)
    if err != nil {
        return fmt.Errorf("parsing data: %w", err)
    }
    
    if err := saveResult(result); err != nil {
        return fmt.Errorf("saving result: %w", err)
    }
    
    return nil
}
```

**Where to handle vs propagate:**
```go
// Handler (boundary layer) — log, respond, recover:
func handleRequest(w http.ResponseWriter, r *http.Request) {
    user, err := getUser(r.Context(), userID)
    if err != nil {
        // At the boundary: log the full error, respond with sanitized message
        log.Printf("getUser: %v", err)
        http.Error(w, "user not found", http.StatusNotFound)
        return
    }
    // ...
}

// Library/business logic — wrap and propagate:
func getUser(ctx context.Context, id int) (*User, error) {
    var u User
    err := db.QueryRowContext(ctx, "SELECT ...", id).Scan(&u.ID, &u.Name)
    if err != nil {
        return nil, fmt.Errorf("querying user %d: %w", id, err)  // Don't log here
    }
    return &u, nil
}
```

### Quick Check
> 1. Should you log AND return an error, or just one of them?
> 2. Is `data, _ := os.ReadFile(...)` ever acceptable? When?
> 3. What context information should you add when wrapping an error?

---

## 4. Error Wrapping and Unwrapping

**Wrapping** adds context to an error while preserving the original:
```go
_, err := os.Open("config.json")
// err might be: &PathError{Op: "open", Path: "config.json", Err: syscall.ENOENT}

wrapped := fmt.Errorf("loading config: %w", err)
// wrapped.Error() == "loading config: open config.json: no such file or directory"
// But the original error is still accessible!
```

**`errors.Is`** — checks if any error in the chain matches a target:
```go
var ErrNotFound = errors.New("not found")

func findUser(id int) error {
    return fmt.Errorf("finding user: %w", ErrNotFound)  // Wrap it
}

err := findUser(42)
if errors.Is(err, ErrNotFound) {
    fmt.Println("user does not exist")  // This prints!
}
// errors.Is checks the ENTIRE chain — it finds ErrNotFound even though it's wrapped
if err == ErrNotFound {
    // This does NOT print — direct equality fails on wrapped errors
}
```

**`errors.As`** — extracts a specific type from the error chain:
```go
type NotFoundError struct {
    Resource string
    ID       int
}
func (e *NotFoundError) Error() string {
    return fmt.Sprintf("%s with ID %d not found", e.Resource, e.ID)
}

func getUser(id int) error {
    return fmt.Errorf("getting user: %w", &NotFoundError{Resource: "user", ID: id})
}

err := getUser(42)

var notFound *NotFoundError
if errors.As(err, &notFound) {
    // notFound is now the *NotFoundError from inside the chain
    fmt.Printf("Looking for %s #%d\n", notFound.Resource, notFound.ID)
    // "Looking for user #42"
}
```

**`errors.Unwrap`** — get the next error in the chain:
```go
err1 := errors.New("root cause")
err2 := fmt.Errorf("middle: %w", err1)
err3 := fmt.Errorf("outer: %w", err2)

fmt.Println(errors.Unwrap(err3))  // "middle: root cause"
fmt.Println(errors.Unwrap(err2))  // "root cause"
fmt.Println(errors.Unwrap(err1))  // nil (no more)
```

**Multiple errors (Go 1.20+):**
```go
// Join multiple errors together:
err1 := errors.New("error 1")
err2 := errors.New("error 2")
combined := errors.Join(err1, err2)
fmt.Println(combined)         // "error 1\nerror 2"
fmt.Println(errors.Is(combined, err1))  // true
fmt.Println(errors.Is(combined, err2))  // true
```

### Quick Check
> 1. What is the difference between `errors.Is` and `==` when comparing errors?
> 2. When would you use `errors.As` instead of `errors.Is`?
> 3. What does `%w` preserve when wrapping an error?

---

## 5. Custom Error Types

Custom error types let callers inspect and react to specific error conditions:

```go
// HTTP-style error with status code:
type APIError struct {
    StatusCode int
    Message    string
    Err        error  // Underlying cause
}

func (e *APIError) Error() string {
    if e.Err != nil {
        return fmt.Sprintf("API error %d: %s: %v", e.StatusCode, e.Message, e.Err)
    }
    return fmt.Sprintf("API error %d: %s", e.StatusCode, e.Message)
}

func (e *APIError) Unwrap() error { return e.Err }  // Enable errors.Is/As chain walk

// Constructor helpers:
func NotFound(msg string) *APIError {
    return &APIError{StatusCode: 404, Message: msg}
}
func Internal(msg string, err error) *APIError {
    return &APIError{StatusCode: 500, Message: msg, Err: err}
}

// Usage:
func getUser(id int) (*User, error) {
    if id <= 0 {
        return nil, NotFound("user not found")
    }
    user, err := db.Query(id)
    if err != nil {
        return nil, Internal("database query failed", err)
    }
    return user, nil
}

// Handling:
user, err := getUser(0)
var apiErr *APIError
if errors.As(err, &apiErr) {
    switch apiErr.StatusCode {
    case 404:
        fmt.Println("Resource does not exist")
    case 500:
        fmt.Println("Server error:", apiErr.Err)
    }
}
```

**Validation errors — collecting multiple issues:**
```go
type FieldError struct {
    Field   string
    Message string
}

type ValidationErrors []FieldError

func (ve ValidationErrors) Error() string {
    msgs := make([]string, len(ve))
    for i, e := range ve {
        msgs[i] = fmt.Sprintf("%s: %s", e.Field, e.Message)
    }
    return strings.Join(msgs, "; ")
}

func validateUser(u User) error {
    var errs ValidationErrors
    
    if u.Name == "" {
        errs = append(errs, FieldError{"name", "is required"})
    }
    if u.Age < 0 || u.Age > 150 {
        errs = append(errs, FieldError{"age", "must be between 0 and 150"})
    }
    if !strings.Contains(u.Email, "@") {
        errs = append(errs, FieldError{"email", "is invalid"})
    }
    
    if len(errs) > 0 {
        return errs
    }
    return nil
}

// Usage:
err := validateUser(User{Age: -5, Email: "notanemail"})
var validationErrs ValidationErrors
if errors.As(err, &validationErrs) {
    for _, e := range validationErrs {
        fmt.Printf("Field %q: %s\n", e.Field, e.Message)
    }
}
```

### Quick Check
> 1. What method must a custom error type implement?
> 2. Why would you add an `Unwrap() error` method to a custom error type?
> 3. How would you return multiple validation errors from a single function?

---

## 6. Sentinel Errors

Sentinel errors are package-level error variables that represent known, named conditions:

```go
package store

import "errors"

// Sentinel errors — exported, starts with Err by convention:
var (
    ErrNotFound     = errors.New("not found")
    ErrUnauthorized = errors.New("unauthorized")
    ErrConflict     = errors.New("resource already exists")
    ErrInvalidInput = errors.New("invalid input")
)

// Usage in functions:
func (s *Store) GetUser(id int) (*User, error) {
    user, ok := s.users[id]
    if !ok {
        return nil, ErrNotFound  // Return the sentinel directly
    }
    return user, nil
}

func (s *Store) CreateUser(u User) error {
    if _, ok := s.users[u.ID]; ok {
        return ErrConflict  // Return the sentinel directly
    }
    s.users[u.ID] = u
    return nil
}
```

**Checking sentinels:**
```go
user, err := store.GetUser(42)
if errors.Is(err, store.ErrNotFound) {
    // Create the user or return 404
}
if errors.Is(err, store.ErrUnauthorized) {
    // Redirect to login
}
```

**Standard library sentinels you should know:**
```go
io.EOF              // End of file (not really an "error" — normal termination)
io.ErrUnexpectedEOF // EOF when more data was expected
os.ErrNotExist      // File or directory does not exist (tested by os.IsNotExist)
os.ErrExist         // File already exists
os.ErrPermission    // Permission denied
context.Canceled    // Context was canceled
context.DeadlineExceeded // Context deadline exceeded
sql.ErrNoRows       // database/sql query returned no rows

// Examples:
if errors.Is(err, io.EOF) {
    break  // End of input stream — stop reading
}
if errors.Is(err, os.ErrNotExist) {
    // File doesn't exist — create it
}
if errors.Is(err, context.DeadlineExceeded) {
    // Request took too long — client probably gave up
}
```

**When to use sentinels vs custom types:**
- **Sentinel**: caller only needs to know WHICH condition occurred, not details about it
- **Custom type**: caller needs additional data about the error (e.g., which field failed, HTTP status code, retry-after duration)

### Quick Check
> 1. What is the naming convention for sentinel errors?
> 2. Why use `errors.Is(err, io.EOF)` instead of `err == io.EOF`?
> 3. When should you prefer a sentinel error over a custom error type?

---

## 7. Error Handling Patterns

### 7.1 Inline Init with if

```go
// Use init statement to scope err to the if block:
if user, err := getUser(id); err != nil {
    return fmt.Errorf("getting user: %w", err)
} else {
    process(user)  // user is accessible here
}
// user and err are NOT accessible here — reduces scope
```

### 7.2 Reduce Nesting with Early Return

```go
// Nested (hard to read):
func processOrder(orderID int) error {
    order, err := getOrder(orderID)
    if err == nil {
        user, err := getUser(order.UserID)
        if err == nil {
            if err := validateOrder(order, user); err == nil {
                return saveOrder(order)
            } else {
                return err
            }
        } else {
            return err
        }
    } else {
        return err
    }
}

// Early return (clean):
func processOrder(orderID int) error {
    order, err := getOrder(orderID)
    if err != nil {
        return fmt.Errorf("getting order %d: %w", orderID, err)
    }

    user, err := getUser(order.UserID)
    if err != nil {
        return fmt.Errorf("getting user %d: %w", order.UserID, err)
    }

    if err := validateOrder(order, user); err != nil {
        return fmt.Errorf("validating order: %w", err)
    }

    return saveOrder(order)
}
```

### 7.3 Errgroup — Concurrent Operations

```go
import "golang.org/x/sync/errgroup"

func loadDashboard(userID int) (*Dashboard, error) {
    g, ctx := errgroup.WithContext(context.Background())
    
    var user *User
    var orders []Order
    var recommendations []Product
    
    g.Go(func() error {
        var err error
        user, err = getUser(ctx, userID)
        return err
    })
    
    g.Go(func() error {
        var err error
        orders, err = getOrders(ctx, userID)
        return err
    })
    
    g.Go(func() error {
        var err error
        recommendations, err = getRecommendations(ctx, userID)
        return err
    })
    
    // Wait for all — returns first non-nil error
    if err := g.Wait(); err != nil {
        return nil, fmt.Errorf("loading dashboard: %w", err)
    }
    
    return &Dashboard{User: user, Orders: orders, Recs: recommendations}, nil
}
```

### 7.4 Retry with Backoff

```go
func withRetry(attempts int, fn func() error) error {
    var err error
    for i := 0; i < attempts; i++ {
        if err = fn(); err == nil {
            return nil  // Success
        }
        
        // Don't retry on non-retriable errors:
        var apiErr *APIError
        if errors.As(err, &apiErr) && apiErr.StatusCode < 500 {
            return err  // 4xx errors are client errors — don't retry
        }
        
        if i < attempts-1 {
            time.Sleep(time.Duration(i+1) * 100 * time.Millisecond)  // Linear backoff
        }
    }
    return fmt.Errorf("after %d attempts: %w", attempts, err)
}

// Usage:
err := withRetry(3, func() error {
    return callExternalAPI()
})
```

### 7.5 Must Pattern (Initialization Only)

```go
// Must panics if err != nil — use ONLY in initialization code:
func Must[T any](v T, err error) T {
    if err != nil {
        panic(err)
    }
    return v
}

// Safe at startup:
var db = Must(sql.Open("postgres", connStr))
var tmpl = Must(template.ParseGlob("templates/*.html"))

// NEVER use Must in request handlers — panics crash the server
```

### Quick Check
> 1. What is the benefit of early returns in error handling?
> 2. When is it appropriate to use the `Must` pattern?
> 3. What does `errgroup.Wait()` return?

---

## Summary

- **error is a value**: it's an interface `{ Error() string }`, not magic
- **`nil` = success**, non-nil = failure — check immediately after each call
- **`errors.New`**: static message; **`fmt.Errorf`**: formatted message; **`%w`**: wraps
- **`errors.Is`**: walks the chain to check identity; **`errors.As`**: walks to extract type
- **Sentinel errors**: package-level `var ErrXxx = errors.New(...)` for known conditions
- **Custom error types**: implement `Error() string` (and optionally `Unwrap() error`)
- **Early return**: keep the happy path unindented; handle errors and return immediately
- **Don't log AND return**: log at the boundary, wrap and propagate in libraries
- **`Must` pattern**: only in program initialization, never in request handlers

---

## Exercises

### Easy
1. Write a `divide(a, b float64) (float64, error)` function. Return an error if `b` is 0. Call it with valid and invalid inputs, handling both cases.
2. Create a sentinel error `ErrNegativeNumber`. Write `sqrt(n float64) (float64, error)` that returns `ErrNegativeNumber` for negative input. Use `errors.Is` to check in the caller.
3. Write a function that reads a file, parses each line as an integer, and returns all parsed integers. If a line fails to parse, wrap the error with the line number.

### Medium
4. HTTP-style errors: Create an `HTTPError` struct with `Code int` and `Message string`. Implement the `error` interface. Add `Unwrap() error`. Write a middleware `func wrapHTTPError(handler func() error) error` that wraps any returned error in a 500 `HTTPError`. Test that `errors.As` can extract the `HTTPError` from a wrapped error.
5. Retry decorator: Write `Retry(maxAttempts int, backoff time.Duration, fn func() error) error` that retries `fn` up to `maxAttempts` times with exponential backoff (`backoff * 2^attempt`). It should NOT retry if the error is a `*ValidationError`. Log each retry attempt with the attempt number and error. Test with: a function that fails twice then succeeds, a function that always fails, and a function that fails with a ValidationError immediately.
6. Error aggregator: Implement a type `MultiError` that can collect multiple errors. `func (m *MultiError) Add(err error)` adds an error. `func (m *MultiError) Err() error` returns nil if no errors were added, otherwise returns the MultiError itself. The `Error()` method returns all messages formatted nicely. Implement `errors.Is` and `errors.As` support. Use it to collect all errors from processing a batch of items without stopping at the first failure.

### Hard
7. Structured error chain: Build an error system for a payment processing service. Define these error types: `NetworkError` (retryable), `AuthError` (not retryable, has token info), `InsufficientFundsError` (not retryable, has balance/amount), `ProcessingError` (wraps another error with transaction ID). Write `ProcessPayment(req PaymentRequest) error` that can return any of these. Write `handlePaymentError(err error)` that reacts differently to each type. Write a test that verifies each error type is detectable through `errors.As` even when wrapped in a `ProcessingError`.
8. Error budget: Implement `ErrorBudget` — a sliding window error tracker for SRE-style error budgets. It tracks success/failure events in the last N minutes. Methods: `Record(success bool)`, `ErrorRate() float64`, `WithinBudget(targetReliability float64) bool`. If `ErrorRate()` exceeds `1 - targetReliability`, return a `BudgetExhaustedError` with the current error rate and target. Make it thread-safe. Write a benchmark comparing performance with mutex vs atomic operations.
