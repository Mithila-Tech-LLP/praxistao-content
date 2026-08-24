# Chapter 17: Goroutine Leaks, Race Conditions & Deadlocks

These three problems kill production Go services. Every senior Go interview will ask about at least one of them. This chapter teaches you to identify, prevent, and debug each one.

## Table of Contents

1. [Goroutine Leaks](#1-goroutine-leaks)
2. [Race Conditions](#2-race-conditions)
3. [Deadlocks](#3-deadlocks)
4. [Detection Tools](#4-detection-tools)
5. [Interview Questions & Model Answers](#5-interview-questions--model-answers)
6. [Summary](#summary)

---

## 1. Goroutine Leaks

A goroutine leaks when it is started but never terminates. Leaked goroutines accumulate over time, consuming memory and CPU. A service with goroutine leaks will slowly die.

### The Seven Causes of Goroutine Leaks

**Cause 1: Sending to a channel with no receiver**

```go
// LEAK: if nobody receives from ch, this goroutine blocks forever
func startWork() {
    ch := make(chan Result)
    go func() {
        result := doWork()
        ch <- result // blocks if caller doesn't receive
    }()
    // caller returns without receiving — goroutine leaks!
}

// FIX: Use buffered channel OR ensure the caller always receives
func startWorkFixed() <-chan Result {
    ch := make(chan Result, 1) // buffered: send won't block even if caller is slow
    go func() {
        ch <- doWork()
    }()
    return ch
}
```

**Cause 2: Receiving from a channel that's never closed or written to**

```go
// LEAK: goroutine waits forever on an empty channel
func leak() {
    ch := make(chan int)
    go func() {
        val := <-ch // blocks forever — nobody sends
        fmt.Println(val)
    }()
}

// FIX: close the channel or use a done signal
func fixed(ctx context.Context) {
    ch := make(chan int)
    go func() {
        select {
        case val := <-ch:
            fmt.Println(val)
        case <-ctx.Done(): // exit when context cancelled
            return
        }
    }()
}
```

**Cause 3: No done signal for infinite loops**

```go
// LEAK: goroutine runs forever with no way to stop it
func startPoller() {
    go func() {
        for {
            poll() // runs forever even after the service shuts down
            time.Sleep(time.Second)
        }
    }()
}

// FIX: always provide a way to stop
func startPollerFixed(ctx context.Context) {
    go func() {
        ticker := time.NewTicker(time.Second)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                poll()
            case <-ctx.Done():
                return // stop when context cancelled
            }
        }
    }()
}
```

**Cause 4: HTTP handler goroutines with no timeout**

```go
// LEAK: if external service never responds, this goroutine leaks
func handleRequest(w http.ResponseWriter, r *http.Request) {
    go func() {
        result := callExternalService() // might block forever
        respond(w, result)
    }()
}

// FIX: use context with timeout, and the request's context
func handleRequestFixed(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 30*time.Second)
    defer cancel()
    result := callExternalServiceCtx(ctx) // respects cancellation
    respond(w, result)
}
```

**Cause 5: WaitGroup with negative counter**

```go
// PANIC/DEADLOCK: Done called more than Add
var wg sync.WaitGroup
wg.Add(1)
wg.Done()
wg.Done() // panic: negative WaitGroup counter
```

**Cause 6: Goroutine blocked on mutex**

```go
// LEAK: mutex never unlocked, goroutines block forever
func processLocked(mu *sync.Mutex) {
    mu.Lock()
    // panic here — mu is never unlocked!
    doRiskyWork() // panics
    mu.Unlock()   // never reached
}

// FIX: always use defer
func processLockedFixed(mu *sync.Mutex) {
    mu.Lock()
    defer mu.Unlock() // runs even on panic
    doRiskyWork()
}
```

**Cause 7: Fan-in where one sender is slow**

```go
// LEAK: if one goroutine sends slowly, others block in the merge function
func merge(channels ...<-chan int) <-chan int {
    out := make(chan int)
    for _, ch := range channels {
        go func(c <-chan int) {
            for v := range c {
                out <- v // blocks if nobody reads out fast enough
            }
        }(ch)
    }
    // If receiver stops reading, all sending goroutines leak!
    return out
}

// FIX: use context cancellation
func mergeWithCancel(ctx context.Context, channels ...<-chan int) <-chan int {
    out := make(chan int)
    for _, ch := range channels {
        go func(c <-chan int) {
            for {
                select {
                case v, ok := <-c:
                    if !ok { return }
                    select {
                    case out <- v:
                    case <-ctx.Done(): return
                    }
                case <-ctx.Done(): return
                }
            }
        }(ch)
    }
    return out
}
```

---

## 2. Race Conditions

A race condition occurs when multiple goroutines access shared mutable state without proper synchronization, and at least one access is a write.

### Common Race Patterns

**Race 1: Incrementing a shared counter**

```go
var count int

func racyIncrement() {
    for i := 0; i < 1000; i++ {
        go func() {
            count++ // RACE: read-modify-write is not atomic
        }()
    }
}

// FIX A: atomic
var count int64
atomic.AddInt64(&count, 1)

// FIX B: mutex
var mu sync.Mutex
var count int
mu.Lock()
count++
mu.Unlock()
```

**Race 2: Loop variable capture**

```go
// RACE: all goroutines share the same 'i' variable
for i := 0; i < 5; i++ {
    go func() {
        fmt.Println(i) // prints 5, 5, 5, 5, 5 (or worse, crashes)
    }()
}

// FIX: pass as argument
for i := 0; i < 5; i++ {
    go func(n int) {
        fmt.Println(n) // each goroutine has its own copy
    }(i)
}
```

**Race 3: Map concurrent access**

```go
m := map[string]int{}
var wg sync.WaitGroup

// RACE: concurrent map read and write
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        m["key"] = n // concurrent writes → race condition → panic in Go 1.6+
    }(i)
}

// FIX A: mutex
var mu sync.Mutex
var wg sync.WaitGroup
for i := 0; i < 100; i++ {
    wg.Add(1)
    go func(n int) {
        defer wg.Done()
        mu.Lock()
        m["key"] = n
        mu.Unlock()
    }(i)
}

// FIX B: sync.Map for concurrent access
var sm sync.Map
sm.Store("key", n)
```

**Race 4: Check-then-act on shared state**

```go
// RACE: another goroutine may change `active` between the check and the action
if !service.active {
    service.active = true // gap between check and set!
    service.Start()
}

// FIX: use sync.Once or CAS
var started sync.Once
started.Do(func() {
    service.active = true
    service.Start()
})
```

---

## 3. Deadlocks

A deadlock occurs when a set of goroutines are all waiting for each other, and none can make progress.

### The Four Coffman Conditions (from OS theory)

A deadlock requires all four:
1. **Mutual exclusion:** resources cannot be shared (mutex)
2. **Hold and wait:** a goroutine holds one resource while waiting for another
3. **No preemption:** resources can only be released voluntarily
4. **Circular wait:** G1 waits for G2, G2 waits for G1

Break any one condition to prevent deadlock.

### Go Deadlock Patterns

**Deadlock 1: Classic two-goroutine lock ordering**

```go
var mu1, mu2 sync.Mutex

// Goroutine A: acquires mu1 then mu2
go func() {
    mu1.Lock(); defer mu1.Unlock()
    time.Sleep(time.Millisecond) // gives time for goroutine B to acquire mu2
    mu2.Lock(); defer mu2.Unlock()
}()

// Goroutine B: acquires mu2 then mu1 — opposite order!
go func() {
    mu2.Lock(); defer mu2.Unlock()
    mu1.Lock(); defer mu1.Unlock() // blocks forever: mu1 held by goroutine A
}()
// DEADLOCK: A holds mu1 waiting for mu2. B holds mu2 waiting for mu1.

// FIX: always acquire locks in the same global order
// Both goroutines: acquire mu1 then mu2
```

**Deadlock 2: Unbuffered channel with no receiver**

```go
ch := make(chan int)
ch <- 42  // DEADLOCK: no goroutine to receive — main goroutine blocks forever
// Go runtime detects this: "all goroutines are asleep - deadlock!"
```

**Deadlock 3: Goroutine waiting for itself**

```go
var mu sync.Mutex
mu.Lock()
defer mu.Unlock()
// ... calls a function that also tries to Lock mu
mu.Lock() // DEADLOCK: same goroutine trying to acquire lock it already holds
// sync.Mutex is NOT reentrant — use sync.RWMutex or refactor
```

**Deadlock 4: WaitGroup.Wait() before all goroutines start**

```go
var wg sync.WaitGroup
go func() {
    time.Sleep(time.Second)
    wg.Add(1) // Add called AFTER Wait might have already returned
    defer wg.Done()
}()
wg.Wait() // might return before the goroutine even calls Add!
```

---

## 4. Detection Tools

### The Race Detector

```bash
# Build with race detector
go build -race ./...

# Run tests with race detector (ALWAYS do this in CI)
go test -race ./...

# Run with race detector
go run -race main.go
```

The race detector instruments every memory access and reports:

```
==================
WARNING: DATA RACE
Write at 0x00c000124010 by goroutine 7:
  main.main.func1()
      /tmp/main.go:12 +0x3c

Previous write at 0x00c000124010 by goroutine 6:
  main.main.func1()
      /tmp/main.go:12 +0x3c
==================
```

**Performance overhead:** The race detector increases memory usage ~5-10x and slows execution ~2-20x. Use in tests and staging, not production.

### Goroutine Leak Detection with pprof

```go
// Check goroutine count
numGoroutines := runtime.NumGoroutine()
fmt.Println("goroutines:", numGoroutines)

// Get goroutine stack traces
pprof.Lookup("goroutine").WriteTo(os.Stdout, 1)
```

```bash
# Get goroutine profile from running service (HTTP pprof enabled)
go tool pprof http://localhost:6060/debug/pprof/goroutine

# In pprof:
(pprof) top     # show goroutine counts by stack
(pprof) list    # show which lines goroutines are stuck at
```

### goleak for Tests

```go
import "go.uber.org/goleak"

func TestMyFunction(t *testing.T) {
    defer goleak.VerifyNone(t) // fails if any goroutines leaked during the test
    
    // run your function
    doWork()
}
```

---

## 5. Interview Questions & Model Answers

**Q: What causes a goroutine leak and how do you detect it?**

"A goroutine leaks when it blocks indefinitely — usually on a channel send or receive with no counterpart, or on a mutex that's never released. In production services, they accumulate over time: high-memory services that never recover, goroutine count monotonically increasing. To detect them: use `runtime.NumGoroutine()` as a metric — a monotonically increasing count is a leak. Profile with `pprof goroutine` to see stack traces of all goroutines and identify the blocked ones. In tests, use the `goleak` library to fail tests when goroutines are left behind."

**Q: How does Go's race detector work?**

"The race detector instruments every memory read and write at compile time (via `-race` flag). It uses a shadow memory to track the 'happens-before' timestamp of the last access to each memory location and the goroutine that made it. When a new access happens, it checks whether a proper synchronization event establishes a happens-before relationship with the previous access. If not, it reports a race. The overhead is 5-10x memory and 2-20x slowdown, so it's used in testing and staging, not production."

**Q: Describe a deadlock you've seen in production and how you fixed it.**

"[Use a real example if you have one, or this template] In a service that managed database connection pools, two components each held a lock on their own pool and tried to acquire the other's lock when negotiating failover. The fix was to establish a consistent lock ordering: always acquire pool A's lock before pool B's, regardless of which component was doing the work. Alternatively, we could have used a single top-level lock for cross-pool operations. The fix was identified by examining the goroutine stacks in pprof, which clearly showed both goroutines blocked on each other's mutex."

---

## Summary

- **Goroutine leaks:** 7 common causes — all involve a goroutine blocked with no way to escape. Prevention: always provide a done/cancel signal. Detection: monitor goroutine count, use goleak in tests, pprof goroutine profile.
- **Race conditions:** concurrent access to shared mutable state without synchronization. Prevention: mutex, atomic, channel ownership transfer, sync.Map. Detection: `go test -race` always in CI.
- **Deadlocks:** circular waiting on resources. Prevention: consistent lock ordering, avoid nested locks, use defer for unlock. Detection: Go runtime detects total deadlock, pprof finds partial deadlocks.
- **Rule:** always run the race detector in tests. Always give goroutines a way to stop. Always defer unlock.
