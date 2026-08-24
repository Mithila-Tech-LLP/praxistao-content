# Chapter 15: Defer, Panic, and Recover

Go provides three special mechanisms for managing unusual control flow: `defer` for cleanup that must always run, `panic` for truly unrecoverable situations, and `recover` for catching panics at boundaries. Used correctly, these three work together to make Go programs robust — used incorrectly, they introduce subtle bugs. This chapter covers both the mechanics and the idioms.

## Table of Contents

1. [Defer — Guaranteed Cleanup](#1-defer--guaranteed-cleanup)
2. [Defer Execution Order and Gotchas](#2-defer-execution-order-and-gotchas)
3. [Panic — The Emergency Stop](#3-panic--the-emergency-stop)
4. [Recover — Catching Panics](#4-recover--catching-panics)
5. [Defer + Recover Pattern](#5-defer--recover-pattern)
6. [Real-World Usage](#6-real-world-usage)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Defer — Guaranteed Cleanup

`defer` schedules a function call to run **when the surrounding function returns**, regardless of whether the function returns normally or panics:

```go
func readFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()  // Will run when readFile returns — no matter what
    
    data, err := io.ReadAll(f)
    if err != nil {
        return "", err  // f.Close() still runs here
    }
    
    return string(data), nil  // f.Close() runs here too
}
```

Without `defer`, you'd have to call `f.Close()` on every return path — and it's easy to forget one:
```go
// Without defer — error-prone:
func readFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    
    data, err := io.ReadAll(f)
    if err != nil {
        f.Close()  // Easy to forget!
        return "", err
    }
    
    f.Close()  // Must remember here too
    return string(data), nil
}
```

**Common defer uses:**
```go
// 1. File handling:
f, _ := os.Create("output.txt")
defer f.Close()

// 2. Mutex unlock:
mu.Lock()
defer mu.Unlock()

// 3. Database transaction:
tx, _ := db.Begin()
defer tx.Rollback()  // No-op if tx.Commit() was called first
// ... do work ...
tx.Commit()

// 4. WaitGroup Done:
wg.Add(1)
go func() {
    defer wg.Done()
    // ... do work ...
}()

// 5. Timing:
start := time.Now()
defer func() {
    fmt.Println("elapsed:", time.Since(start))
}()
```

### Quick Check
> 1. When does a deferred function run?
> 2. Does defer still run if the function returns an error?
> 3. Name two common uses of defer.

---

## 2. Defer Execution Order and Gotchas

**LIFO order** — defers run in last-in, first-out order:
```go
func main() {
    defer fmt.Println("first defer")
    defer fmt.Println("second defer")
    defer fmt.Println("third defer")
    
    fmt.Println("main body")
}
// Output:
// main body
// third defer
// second defer
// first defer
```

This mirrors resource acquisition: open A, open B, open C → close C, close B, close A.

**Arguments are evaluated immediately:**
```go
x := 10
defer fmt.Println(x)  // x is captured as 10 RIGHT NOW
x = 20
fmt.Println(x)  // 20

// Output:
// 20
// 10  ← NOT 20, even though x changed
```

**Closures capture by reference:**
```go
x := 10
defer func() {
    fmt.Println(x)  // Closure captures x by reference
}()
x = 20
fmt.Println(x)  // 20

// Output:
// 20
// 20  ← The closure sees the current value of x when it runs
```

**Named return values + defer — advanced:**
```go
// Named returns allow defer to modify the return value!
func withCleanup() (result string, err error) {
    defer func() {
        if err != nil {
            result = "error occurred"  // Modify named return!
        }
    }()
    
    // ... do work ...
    return "success", nil
}
```

**Defer in a loop — common bug:**
```go
// BUG: files only close AFTER the entire function returns (not after each iteration)
func processFiles(paths []string) error {
    for _, path := range paths {
        f, err := os.Open(path)
        if err != nil {
            return err
        }
        defer f.Close()  // All defers queue up, none run until function exits
        
        process(f)
    }
    return nil
    // All files close here — holding many open file handles simultaneously!
}

// FIX: wrap in a helper function so defer runs per iteration:
func processFiles(paths []string) error {
    for _, path := range paths {
        if err := processOne(path); err != nil {
            return err
        }
    }
    return nil
}

func processOne(path string) error {
    f, err := os.Open(path)
    if err != nil {
        return err
    }
    defer f.Close()  // Runs when processOne returns — correct!
    return process(f)
}
```

### Quick Check
> 1. What is the execution order of multiple defers in one function?
> 2. If `defer fmt.Println(x)` is called when `x=5`, and then `x` changes to 10, what prints?
> 3. Why is defer-in-a-loop dangerous for file handles?

---

## 3. Panic — The Emergency Stop

`panic` stops normal execution, runs all deferred functions in the current goroutine, and crashes the program (unless recovered):

```go
func mustNotBeNegative(n int) {
    if n < 0 {
        panic(fmt.Sprintf("expected non-negative, got %d", n))
    }
}

func main() {
    mustNotBeNegative(-1)
    fmt.Println("this never prints")
}

// Output:
// goroutine 1 [running]:
// main.mustNotBeNegative(...)
//         /tmp/main.go:4
// main.main()
//         /tmp/main.go:9 +0x...
// exit status 2
```

**When Go panics automatically:**
```go
// 1. Nil pointer dereference:
var p *int
_ = *p  // panic: runtime error: invalid memory address or nil pointer dereference

// 2. Out of bounds slice/array access:
s := []int{1, 2, 3}
_ = s[10]  // panic: runtime error: index out of range [10] with length 3

// 3. Type assertion failure:
var i interface{} = "hello"
n := i.(int)  // panic: interface conversion: interface {} is string, not int
// Use comma-ok form to avoid: n, ok := i.(int)

// 4. Send on closed channel:
ch := make(chan int)
close(ch)
ch <- 1  // panic: send on closed channel

// 5. Division by zero (integers only):
x := 5
y := 0
_ = x / y  // panic: runtime error: integer divide by zero
```

**When to use panic:**
```go
// 1. Programmer errors — conditions that should never happen:
func NewServer(port int) *Server {
    if port < 1 || port > 65535 {
        panic(fmt.Sprintf("invalid port %d: must be 1-65535", port))
    }
    return &Server{port: port}
}

// 2. Impossible states — code that should be unreachable:
switch direction {
case North, South, East, West:
    move(direction)
default:
    panic(fmt.Sprintf("unknown direction: %v", direction))
}

// 3. Initialization failures:
var db *sql.DB
func init() {
    var err error
    db, err = sql.Open("postgres", connStr)
    if err != nil {
        panic("failed to connect to database: " + err.Error())
    }
}
```

**When NOT to use panic:**
```go
// Don't panic for expected errors — use the (value, error) pattern:

// Bad:
func divide(a, b float64) float64 {
    if b == 0 {
        panic("division by zero")  // Caller can't handle this!
    }
    return a / b
}

// Good:
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, errors.New("division by zero")
    }
    return a / b, nil
}
```

### Quick Check
> 1. What happens to deferred functions when a panic occurs?
> 2. Name three situations where Go panics automatically.
> 3. Should you use panic for a database connection failure during a web request?

---

## 4. Recover — Catching Panics

`recover` stops a panic and returns the value that was passed to `panic`. It only works inside a deferred function:

```go
func safeDiv(a, b int) (result int, err error) {
    defer func() {
        if r := recover(); r != nil {
            err = fmt.Errorf("recovered from panic: %v", r)
        }
    }()
    
    result = a / b  // Panics if b == 0
    return result, nil
}

n, err := safeDiv(10, 0)
fmt.Println(n, err)  // 0 recovered from panic: runtime error: integer divide by zero

n, err = safeDiv(10, 2)
fmt.Println(n, err)  // 5 <nil>
```

**Key rules for recover:**
1. `recover()` returns `nil` if there is no panic (or if called outside a defer)
2. `recover()` ONLY works inside a `defer` function — not just any function called from a defer
3. After recovery, normal execution continues in the caller of the panicking function

```go
// recover only works directly inside defer:
func wrong() {
    handlePanic()  // BUG: recover inside handlePanic won't catch THIS function's panics
    panic("oops")
}

func handlePanic() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("recovered:", r)
        }
    }()
}

// Correct:
func correct() {
    defer func() {
        if r := recover(); r != nil {  // recover is DIRECTLY inside the defer
            fmt.Println("recovered:", r)
        }
    }()
    panic("oops")
}
```

### Quick Check
> 1. What does `recover()` return when there is no panic?
> 2. Must `recover()` be called directly inside a `defer`, or can it be in a function called by a defer?
> 3. After a successful recovery, where does execution continue?

---

## 5. Defer + Recover Pattern

The canonical pattern for preventing panics from crashing a server:

```go
// HTTP middleware that recovers from panics:
func recoveryMiddleware(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        defer func() {
            if r := recover(); r != nil {
                // Log the panic with stack trace:
                buf := make([]byte, 4096)
                n := runtime.Stack(buf, false)
                log.Printf("panic: %v\n%s", r, buf[:n])
                
                // Return 500 to the client (don't leak panic details):
                http.Error(w, "Internal Server Error", http.StatusInternalServerError)
            }
        }()
        
        next.ServeHTTP(w, r)
    })
}

// Usage:
mux := http.NewServeMux()
mux.HandleFunc("/", handler)
http.ListenAndServe(":8080", recoveryMiddleware(mux))
```

**Converting panic to error** — useful when calling third-party code:
```go
func safeCall(fn func()) (err error) {
    defer func() {
        if r := recover(); r != nil {
            switch v := r.(type) {
            case error:
                err = v  // Re-use the error if panic value is already an error
            default:
                err = fmt.Errorf("panic: %v", v)
            }
        }
    }()
    fn()
    return nil
}

// Useful for calling untrusted/panicky code:
err := safeCall(func() {
    riskyThirdPartyCode()
})
if err != nil {
    log.Printf("call failed: %v", err)
}
```

**Re-panicking** — recover, inspect, re-panic if not your problem:
```go
func (s *Server) handleRequest(w http.ResponseWriter, r *http.Request) {
    defer func() {
        if rec := recover(); rec != nil {
            // Only handle our own panics; re-panic on unknown ones
            if err, ok := rec.(*AppError); ok {
                http.Error(w, err.Error(), err.StatusCode)
                return
            }
            panic(rec)  // Re-panic — not our problem
        }
    }()
    // ... handle request ...
}
```

### Quick Check
> 1. Why should an HTTP server use a recovery middleware?
> 2. How do you re-panic after catching a panic in recover?
> 3. How do you get a stack trace at the point of a panic?

---

## 6. Real-World Usage

**Cleanup with error-aware defer:**
```go
func transferFunds(from, to *Account, amount float64) (err error) {
    tx, err := db.Begin()
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    
    // If any error occurs after this point, rollback:
    defer func() {
        if err != nil {
            tx.Rollback()  // Rollback on error
        }
    }()
    
    if err = debit(tx, from, amount); err != nil {
        return fmt.Errorf("debit from account: %w", err)
    }
    
    if err = credit(tx, to, amount); err != nil {
        return fmt.Errorf("credit to account: %w", err)
    }
    
    if err = tx.Commit(); err != nil {
        return fmt.Errorf("commit transaction: %w", err)
    }
    
    return nil  // err is nil here — no rollback
}
```

**Tracing with defer:**
```go
func trace(name string) func() {
    start := time.Now()
    log.Printf("START %s", name)
    return func() {
        log.Printf("END %s took %s", name, time.Since(start))
    }
}

func doSomething() {
    defer trace("doSomething")()  // Call trace, then immediately defer the returned func
    // ... work ...
}
```

**Assert pattern for tests:**
```go
func mustNoError(t *testing.T, err error) {
    t.Helper()
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
}

// In tests:
func TestCreateUser(t *testing.T) {
    user, err := store.CreateUser(User{Name: "Alice"})
    mustNoError(t, err)
    // ...
}
```

### Quick Check
> 1. In the `transferFunds` example, why does the defer check `if err != nil`?
> 2. What does `defer trace("f")()` do — why the extra `()`?
> 3. How does the named return `err` enable the defer to see the return error?

---

## Summary

- **`defer`**: schedules a call for when the function returns — runs even on panic
- **LIFO order**: last deferred runs first — mirrors nested open/close patterns
- **Arguments evaluated immediately**: `defer f(x)` captures `x` NOW, not when it runs
- **Closures capture by reference**: use a closure if you need the final value of a variable
- **Don't defer in loops**: use a helper function instead — otherwise descriptors accumulate
- **`panic`**: for programmer errors and impossible states — not expected errors
- **Auto-panics**: nil deref, out-of-bounds, failed type assertion, send on closed channel
- **`recover`**: must be directly inside a `defer` — returns nil if no panic
- **Recovery pattern**: HTTP middleware, converting panics to errors, re-panicking unknown ones
- **Named returns + defer**: powerful pattern for guaranteed cleanup with error awareness

---

## Exercises

### Easy
1. Write a function `cleanup(files []*os.File)` that defers a close on each file. Verify that all files are closed even if processing fails midway (hint: wrap in a helper).
2. Write a `timed(name string, fn func())` function that uses `defer` to print how long `fn` took to run.
3. Write a function that panics with a custom string. Call it from `main` using a `defer`+`recover` in `main` to catch the panic and print a friendly message.

### Medium
4. Transaction helper: Write `func WithTransaction(db *sql.DB, fn func(*sql.Tx) error) error` that: begins a transaction, calls `fn` with it, commits if `fn` returns nil, rolls back if `fn` returns an error. Use `defer` for the rollback. Test with a function that succeeds and one that fails mid-way.
5. Safe goroutine launcher: Write `func Go(fn func()) <-chan error` that launches `fn` in a goroutine, recovering from any panic, and sends any error (or recovered panic as error) to the returned channel. If `fn` completes without panic, send nil. Test with: a goroutine that succeeds, one that returns an error, and one that panics.
6. Defer argument evaluation quiz: Without running them, predict the output of each of these programs. Then run them to verify. Write a brief explanation for each: (a) `x:=1; defer fmt.Println(x); x=2`, (b) `x:=1; defer func(){fmt.Println(x)}(); x=2`, (c) multiple defers in reverse order, (d) named return modified by defer.

### Hard
7. Panic-safe plugin system: Design a `PluginRunner` that loads and runs plugins (represented as `func() error`). Each plugin runs in isolation — a panic in one plugin should not crash the runner. After all plugins run, return a summary: how many succeeded, how many failed normally (returned error), how many panicked (and with what). Use goroutines so plugins run concurrently (max 4 at a time using a semaphore channel). Test with 10 plugins where some succeed, some return errors, and some panic.
8. Context-aware defer: Implement a `DeferStack` — a stack of cleanup functions that can be committed or rolled back. `Push(name string, cleanup func() error)` adds a cleanup. `Commit()` runs all cleanups in LIFO order, stopping at the first error. `Rollback()` runs ALL cleanups in LIFO order, collecting errors. This models a multi-step transaction where each step has a corresponding undo. Test with a sequence of steps where step 3 fails — verify Rollback runs cleanup for steps 2 and 1 in order.
