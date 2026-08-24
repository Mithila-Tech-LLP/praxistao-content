# Chapter 18: Goroutines — Concurrency in Go

Goroutines are Go's answer to the question: "how do you do many things at once without the overhead and complexity of OS threads?" A goroutine is a lightweight, cooperatively scheduled function that runs concurrently with other goroutines. You can run millions of them simultaneously. This chapter explains how goroutines work, how to launch them safely, and the gotchas that trip up beginners.

## Table of Contents

1. [What Is a Goroutine](#1-what-is-a-goroutine)
2. [Launching Goroutines](#2-launching-goroutines)
3. [Goroutines vs OS Threads](#3-goroutines-vs-os-threads)
4. [The Race Condition Problem](#4-the-race-condition-problem)
5. [WaitGroups — Waiting for Goroutines](#5-waitgroups--waiting-for-goroutines)
6. [The Go Scheduler](#6-the-go-scheduler)
7. [Goroutine Leaks](#7-goroutine-leaks)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is a Goroutine

A goroutine is a function that runs **concurrently** with other functions. It starts with the keyword `go`:

```go
func sayHello(name string) {
    fmt.Printf("Hello, %s!\n", name)
}

func main() {
    go sayHello("Alice")  // Launch as goroutine — runs concurrently
    go sayHello("Bob")    // Another goroutine
    
    fmt.Println("main continuing...")
    time.Sleep(10 * time.Millisecond)  // Wait for goroutines to finish
}
// Possible output (order is not guaranteed):
// main continuing...
// Hello, Alice!
// Hello, Bob!
```

**Key points:**
- `go f()` starts `f` as a goroutine and returns **immediately** — it doesn't wait for `f` to finish
- Goroutines run concurrently (potentially in parallel on multiple cores)
- The execution ORDER of goroutines is not guaranteed
- When `main` returns, ALL goroutines are killed — even if they're still running

### Why not just use threads?
OS threads are expensive: ~1MB of stack space each, slow to create (~100µs), expensive context switches. Go's goroutines start with ~2KB of stack (growing as needed), are cheap to create (~1µs), and are scheduled by Go's own scheduler — not the OS. You can have millions of goroutines where you could only have thousands of threads.

### Quick Check
> 1. What keyword launches a goroutine?
> 2. Does `go f()` wait for `f` to finish?
> 3. What happens to all goroutines when `main()` returns?

---

## 2. Launching Goroutines

**Launching a named function:**
```go
func processItem(item string) {
    fmt.Println("Processing:", item)
    time.Sleep(100 * time.Millisecond)
}

for _, item := range items {
    go processItem(item)
}
```

**Launching an anonymous function (goroutine closure):**
```go
for _, item := range items {
    item := item  // Capture! (See Section 4 — loop variable bug)
    go func() {
        fmt.Println("Processing:", item)
    }()
}

// Passing values explicitly (cleaner):
for _, item := range items {
    go func(i string) {
        fmt.Println("Processing:", i)
    }(item)  // item is passed as argument — evaluated NOW
}
```

**The loop variable bug:**
```go
// BUG: all goroutines may print the SAME (last) value of item
for _, item := range items {
    go func() {
        fmt.Println(item)  // Captures variable, not value — ALL see last item
    }()
}

// FIX 1: shadow the variable
for _, item := range items {
    item := item  // New variable per iteration
    go func() {
        fmt.Println(item)  // Each closure has its own item
    }()
}

// FIX 2: pass as argument (preferred — explicit)
for _, item := range items {
    go func(i string) {
        fmt.Println(i)
    }(item)
}
```

**Using goroutines for I/O-bound parallelism:**
```go
func fetchAll(urls []string) []string {
    results := make([]string, len(urls))
    
    for i, url := range urls {
        i, url := i, url  // Capture
        go func() {
            resp, err := http.Get(url)
            if err == nil {
                body, _ := io.ReadAll(resp.Body)
                resp.Body.Close()
                results[i] = string(body)
            }
        }()
    }
    
    time.Sleep(5 * time.Second)  // Bad! Don't do this — use WaitGroup instead
    return results
}
```

### Quick Check
> 1. Why must you capture loop variables when using goroutines in loops?
> 2. What are two ways to fix the loop variable bug?
> 3. Why is `time.Sleep` a bad way to wait for goroutines?

---

## 3. Goroutines vs OS Threads

Understanding the difference helps you reason about performance and behavior:

```
OS Thread:
  - Created by the OS kernel
  - ~1MB default stack
  - ~1-10µs to create
  - Context switch handled by OS (~1-10µs)
  - Hard limit: typically thousands per process

Goroutine:
  - Created by the Go runtime
  - ~2KB initial stack (grows/shrinks as needed, up to 1GB)
  - ~1µs to create
  - Context switch handled by Go scheduler (~0.1µs)
  - Can have millions per process
```

**The M:N threading model:**
```
Goroutines (G): your code — can have millions
      ↓  scheduled by
OS Threads (M): kernel threads — typically = number of CPUs
      ↓  running on
CPU Cores (P): actual hardware parallelism
```

Go's scheduler multiplexes goroutines onto OS threads automatically. When a goroutine blocks (network I/O, channel, mutex), the scheduler moves other goroutines onto that OS thread. This is why 10,000 goroutines doing network I/O don't need 10,000 OS threads.

**GOMAXPROCS:**
```go
import "runtime"

// How many OS threads can run Go code simultaneously (default: number of CPUs):
fmt.Println(runtime.GOMAXPROCS(0))  // 0 = query, don't change

// Set it (rarely needed):
runtime.GOMAXPROCS(4)  // Use 4 OS threads

// Number of currently running goroutines:
fmt.Println(runtime.NumGoroutine())
```

### Quick Check
> 1. How much stack does a goroutine start with vs an OS thread?
> 2. What is GOMAXPROCS?
> 3. What happens when a goroutine blocks on I/O?

---

## 4. The Race Condition Problem

A **race condition** occurs when two goroutines access the same data simultaneously and at least one is writing:

```go
// BUG: multiple goroutines writing to `count` simultaneously
count := 0
var wg sync.WaitGroup

for i := 0; i < 1000; i++ {
    wg.Add(1)
    go func() {
        defer wg.Done()
        count++  // READ-MODIFY-WRITE — not atomic!
    }()
}

wg.Wait()
fmt.Println(count)  // Likely NOT 1000 — could be 987, 951, anything
```

**Why `count++` is not safe:**
```
count++ is actually three operations:
  1. temp = count    (read)
  2. temp = temp + 1 (compute)
  3. count = temp    (write)

If two goroutines both read count=5 simultaneously:
  G1: read 5, compute 6, write 6
  G2: read 5, compute 6, write 6
  Result: count=6 (lost one increment!)
```

**Detecting races — use the race detector:**
```bash
go test -race ./...       # Run tests with race detector
go run -race main.go      # Run with race detector
go build -race .          # Build with race detector enabled
```

The race detector instruments your code and reports data races:
```
WARNING: DATA RACE
Write at 0x00c0000b4010 by goroutine 7:
  main.main.func1()
        /tmp/main.go:12 +0x4c
Read at 0x00c0000b4010 by goroutine 8:
  main.main.func1()
        /tmp/main.go:12 +0x3c
```

**Solutions (covered in depth in Ch 19 — Channels and Ch 20 — Sync Package):**
```go
// Solution 1: Mutex (mutual exclusion)
var mu sync.Mutex
count := 0

go func() {
    mu.Lock()
    count++
    mu.Unlock()
}()

// Solution 2: Channel (coordination via message passing)
ch := make(chan int, 1)
ch <- 0
go func() {
    count := <-ch
    ch <- count + 1
}()

// Solution 3: Atomic operations (for simple integers)
var count atomic.Int64
go func() {
    count.Add(1)
}()
```

### Quick Check
> 1. What is a race condition?
> 2. Why is `count++` not safe across goroutines?
> 3. How do you enable the race detector?

---

## 5. WaitGroups — Waiting for Goroutines

`sync.WaitGroup` is the standard way to wait for a group of goroutines to finish:

```go
import "sync"

func main() {
    var wg sync.WaitGroup
    
    items := []string{"a", "b", "c", "d"}
    
    for _, item := range items {
        wg.Add(1)  // Increment BEFORE launching the goroutine
        item := item
        go func() {
            defer wg.Done()  // Decrement when goroutine finishes
            process(item)
        }()
    }
    
    wg.Wait()  // Block until count reaches 0
    fmt.Println("All done!")
}
```

**The three methods:**
- `wg.Add(n)` — increment the counter by n (call BEFORE `go`)
- `wg.Done()` — decrement by 1 (call inside the goroutine, always via `defer`)
- `wg.Wait()` — block until counter reaches 0

**Common mistake — Add inside the goroutine:**
```go
// BUG: race condition — main might reach wg.Wait() before goroutines call wg.Add(1)
for _, item := range items {
    go func() {
        wg.Add(1)  // Too late! Main might have already passed Wait()
        defer wg.Done()
        process(item)
    }()
}
wg.Wait()

// CORRECT: Add before go
for _, item := range items {
    wg.Add(1)     // Increment BEFORE launching
    go func() {
        defer wg.Done()
        process(item)
    }()
}
wg.Wait()
```

**Collecting results from goroutines safely:**
```go
func processAll(items []string) []Result {
    results := make([]Result, len(items))  // Pre-allocate — each goroutine writes to its own index
    var wg sync.WaitGroup
    
    for i, item := range items {
        i, item := i, item
        wg.Add(1)
        go func() {
            defer wg.Done()
            results[i] = processItem(item)  // Writing to different indices — safe!
        }()
    }
    
    wg.Wait()
    return results
}
// Safe because each goroutine writes to a unique index — no overlap
```

### Quick Check
> 1. What are the three WaitGroup methods?
> 2. Why must `wg.Add(1)` be called BEFORE `go func()`?
> 3. When is it safe for multiple goroutines to write to different slice indices?

---

## 6. The Go Scheduler

Understanding the scheduler helps you write better concurrent code:

**Goroutine states:**
```
Running   → actually executing on an OS thread
Runnable  → ready to run, waiting for an OS thread
Blocked   → waiting for: channel, mutex, I/O, syscall, sleep
Dead      → finished
```

**When a goroutine is preempted (gives up its OS thread):**
1. Makes a blocking call (I/O, channel operation, `time.Sleep`)
2. Calls `runtime.Gosched()` (explicit yield)
3. Makes a function call (Go 1.14+: goroutines are asynchronously preemptible)

**Explicit yield:**
```go
// Rarely needed, but useful for tight CPU loops:
for {
    doExpensiveWork()
    runtime.Gosched()  // Yield to other goroutines
}
```

**Goroutine inspection:**
```go
fmt.Println("Goroutines:", runtime.NumGoroutine())

// Print a stack trace of ALL goroutines (useful for debugging leaks):
buf := make([]byte, 1<<20)
n := runtime.Stack(buf, true)  // true = all goroutines
fmt.Printf("%s", buf[:n])
```

### Quick Check
> 1. Name three situations when a goroutine gives up its OS thread.
> 2. What does `runtime.NumGoroutine()` return?
> 3. What does `runtime.Stack(buf, true)` do?

---

## 7. Goroutine Leaks

A **goroutine leak** occurs when a goroutine is launched but never terminates — it runs forever, consuming memory and CPU:

```go
// LEAK: goroutine blocks on channel forever if nobody sends
func processRequest(req Request) {
    results := make(chan Result)
    go func() {
        results <- compute(req)  // Blocks if nobody receives!
    }()
    
    // If we return early (timeout, error), the goroutine is stuck forever
    select {
    case r := <-results:
        use(r)
    case <-time.After(5 * time.Second):
        return  // Goroutine still alive — stuck on `results <- ...`
    }
}

// FIX: buffered channel or context cancellation (Ch 20):
func processRequest(req Request) {
    results := make(chan Result, 1)  // Buffer of 1 — goroutine won't block
    go func() {
        results <- compute(req)  // Can always send without blocking
    }()
    
    select {
    case r := <-results:
        use(r)
    case <-time.After(5 * time.Second):
        return  // Goroutine will send to buffered channel and finish
    }
}
```

**Common leak patterns:**
```go
// 1. Goroutine waiting on a channel that nobody will send to:
go func() {
    data := <-abandonedChannel  // Nobody sends here anymore
    process(data)
}()

// 2. Goroutine in infinite loop with no exit condition:
go func() {
    for {
        doWork()
        // No break condition, no channel to signal stop
    }
}()

// 3. HTTP handler launching goroutines without proper lifecycle:
func handler(w http.ResponseWriter, r *http.Request) {
    go expensiveBackground()  // Runs after handler returns — can accumulate
}
```

**Detecting leaks:**
```go
// In tests — use goleak:
// go get go.uber.org/goleak

func TestMyFunc(t *testing.T) {
    defer goleak.VerifyNone(t)
    // ... test code ...
    // goleak fails the test if any goroutines leaked
}

// Manually — check goroutine count:
before := runtime.NumGoroutine()
myFunction()
time.Sleep(time.Millisecond)
after := runtime.NumGoroutine()
if after > before {
    fmt.Printf("Possible leak: %d goroutines added\n", after-before)
}
```

**Rule: every goroutine needs an exit path.** Use `context.Context` for cancellation (Ch 20).

### Quick Check
> 1. What is a goroutine leak?
> 2. How does a buffered channel help prevent leaks?
> 3. What tool can detect goroutine leaks in tests?

---

## Summary

- **`go f()`**: launches `f` as a goroutine — doesn't wait, returns immediately
- **Goroutines vs threads**: ~2KB stack, ~1µs to create, millions possible; Go scheduler handles OS thread multiplexing
- **Loop variable bug**: always shadow loop variables (`i := i`) or pass as function arguments before launching goroutines
- **Race conditions**: concurrent unsynchronized access to shared data — use mutex, channels, or atomics
- **Race detector**: `go test -race` / `go run -race` — always enable in development and CI
- **WaitGroup**: `Add(1)` before `go`, `defer Done()` inside, `Wait()` to block until all finish
- **GOMAXPROCS**: number of OS threads that can run Go code in parallel (default = CPU count)
- **Goroutine leaks**: goroutines that never finish — every goroutine needs an exit condition

---

## Exercises

### Easy
1. Launch 5 goroutines, each printing a message with their index. Use `sync.WaitGroup` to wait for all. Run 3 times and note how the order changes.
2. Write `parallelSum(nums []int) int` that splits a slice in half, computes the sum of each half in a separate goroutine, and returns the total. Use WaitGroup.
3. Write a program that launches 100 goroutines, each incrementing a counter using `sync.Mutex`. Verify the final count is exactly 100. Then run with `-race` to confirm no data races.

### Medium
4. Fan-out processor: Write `func ProcessConcurrently[T, R any](items []T, maxWorkers int, fn func(T) R) []R` that processes items concurrently with at most `maxWorkers` goroutines running simultaneously. Results must be in the same order as input. Use WaitGroup and a semaphore (buffered channel of size `maxWorkers`).
5. Goroutine pool: Implement a fixed-size worker pool with `N` goroutines that read from a jobs channel and write results to a results channel. `Pool.Submit(job Job)` sends a job. `Pool.Results() <-chan Result` returns the results channel. `Pool.Stop()` shuts down all workers gracefully. Test with 100 jobs and 5 workers.
6. Parallel web scraper: Write a function that takes a list of URLs and fetches them all concurrently (max 10 at a time). Return a map of URL → response body (or error string). Handle: HTTP errors, timeouts (5 second per URL), and URL fetch panics safely with recover.

### Hard
7. Goroutine lifecycle manager: Build a `Manager` that tracks goroutines it launched. It must: limit concurrency (configurable max), report live goroutine count, support graceful shutdown (wait for in-progress goroutines to finish, reject new ones), and detect and report leaks (goroutines that run longer than a configurable max duration). Use atomics for the counter. Write tests that: verify max concurrency is respected, verify shutdown behavior, simulate a "slow" goroutine and verify leak detection.
8. Benchmark goroutine overhead: Write benchmarks comparing: (a) function call, (b) goroutine launch + WaitGroup wait, (c) goroutine from a pool (pre-allocated), (d) goroutine via `sync.Pool` reuse. Measure: time per operation, allocations per operation, scalability from 1 to GOMAXPROCS to 100× GOMAXPROCS goroutines. Write a markdown analysis of the results — what is the minimum overhead worth paying a goroutine for?
