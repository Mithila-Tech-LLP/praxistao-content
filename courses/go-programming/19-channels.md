# Chapter 19: Channels — Communication Between Goroutines

Channels are the pipes through which goroutines communicate. Go's philosophy is: **"Don't communicate by sharing memory; share memory by communicating."** Instead of letting goroutines fight over shared variables with locks, channels pass data directly between goroutines. This makes concurrent programs easier to reason about and less prone to data races.

## Table of Contents

1. [What Is a Channel](#1-what-is-a-channel)
2. [Buffered vs Unbuffered Channels](#2-buffered-vs-unbuffered-channels)
3. [Channel Direction and Closing](#3-channel-direction-and-closing)
4. [Range Over Channel](#4-range-over-channel)
5. [Select Statement](#5-select-statement)
6. [Pipeline Pattern](#6-pipeline-pattern)
7. [Fan-Out and Fan-In](#7-fan-out-and-fan-in)
8. [Done Channel Pattern](#8-done-channel-pattern)
9. [Channel Pitfalls](#9-channel-pitfalls)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Is a Channel

A channel is a typed conduit for passing values between goroutines. Think of it as a thread-safe queue:

```go
// Create a channel:
ch := make(chan int)        // Unbuffered int channel
ch := make(chan string, 5)  // Buffered string channel with capacity 5
ch := make(chan struct{})   // Signal channel (no data, just signals)

// Send a value (blocks until a receiver is ready — for unbuffered):
ch <- 42

// Receive a value (blocks until a sender sends):
val := <-ch

// Receive with ok check (ok=false means channel is closed and empty):
val, ok := <-ch
```

**A simple example:**
```go
func main() {
    ch := make(chan string)
    
    go func() {
        ch <- "hello from goroutine"  // Send — blocks until main receives
    }()
    
    msg := <-ch  // Receive — blocks until goroutine sends
    fmt.Println(msg)  // "hello from goroutine"
}
```

The channel **synchronizes** the two goroutines — the sender waits until the receiver is ready, and vice versa. This is the unbuffered channel's guarantee.

**Channels are first-class values** — you can pass them to functions, return them, store in structs:
```go
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)
}

func consumer(ch <-chan int) {
    for v := range ch {
        fmt.Println(v)
    }
}

func main() {
    ch := make(chan int)
    go producer(ch)
    consumer(ch)
}
```

### Quick Check
> 1. What does `make(chan int)` create?
> 2. What happens when you send to an unbuffered channel with no receiver?
> 3. How do you check if a channel is closed after receiving?

---

## 2. Buffered vs Unbuffered Channels

**Unbuffered channel** — synchronous, sender and receiver must both be ready:
```go
ch := make(chan int)  // Capacity = 0

// This DEADLOCKS — sends block until someone receives:
ch <- 1   // Block... waiting for receiver
ch <- 2   // Never reached

// This works — goroutine receives while main sends:
go func() {
    fmt.Println(<-ch)  // Receive first
    fmt.Println(<-ch)
}()
ch <- 1   // Now main can send
ch <- 2
```

**Buffered channel** — asynchronous up to capacity:
```go
ch := make(chan int, 3)  // Capacity = 3

// These DON'T block (buffer has space):
ch <- 1
ch <- 2
ch <- 3

// This BLOCKS — buffer is full:
ch <- 4  // Waits until someone receives

// Receive values:
fmt.Println(<-ch)  // 1
fmt.Println(<-ch)  // 2
fmt.Println(<-ch)  // 3
```

**Choosing capacity:**
```go
// 1. Unbuffered (0): for synchronization/rendezvous
results := make(chan int)   // Sender waits for receiver

// 2. Buffer of 1: prevent a goroutine from blocking when the receiver may be slow
done := make(chan struct{}, 1)  // Common for "signal done" channels

// 3. Fixed buffer: for bounded work queues
jobs := make(chan Job, 100)  // Buffered queue of up to 100 jobs

// 4. Buffer = goroutine count: so all goroutines can send without waiting
results := make(chan Result, len(items))  // Each goroutine sends once
```

**`len` and `cap` of channels:**
```go
ch := make(chan int, 5)
ch <- 1
ch <- 2
fmt.Println(len(ch))  // 2 — items currently in buffer
fmt.Println(cap(ch))  // 5 — total buffer capacity
```

### Quick Check
> 1. When does sending to a buffered channel block?
> 2. What does `make(chan int, 1)` create?
> 3. What does `len(ch)` return for a channel?

---

## 3. Channel Direction and Closing

**Directional channels** restrict which operations are allowed:
```go
chan int      // Bidirectional: can send and receive
chan<- int    // Send-only: can only send
<-chan int    // Receive-only: can only receive
```

Use directional channels to make APIs explicit and prevent misuse:
```go
// producer only sends — compiler prevents accidental receive
func producer(ch chan<- int) {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // Can close a send-only channel
}

// consumer only receives — compiler prevents accidental send
func consumer(ch <-chan int) {
    for v := range ch {
        fmt.Println(v)
    }
    // close(ch)  // COMPILE ERROR: cannot close receive-only channel
}

// main creates bidirectional channel, passes typed versions to functions
func main() {
    ch := make(chan int)  // Bidirectional
    go producer(ch)      // Implicitly converts to chan<- int
    consumer(ch)         // Implicitly converts to <-chan int
}
```

**Closing channels:**
```go
ch := make(chan int, 3)
ch <- 1
ch <- 2
close(ch)  // Signal "no more values will be sent"

// Receiving from closed channel:
v, ok := <-ch
fmt.Println(v, ok)  // 1 true  — still buffered values

v, ok = <-ch
fmt.Println(v, ok)  // 2 true

v, ok = <-ch
fmt.Println(v, ok)  // 0 false  — channel closed and empty, zero value returned

// RULES for close:
// 1. Only the SENDER closes (not the receiver)
// 2. Close only once — closing a closed channel PANICS
// 3. Sending to a closed channel PANICS
// 4. You can receive from a closed channel (drains buffer, then zero + false)
```

**The nil channel:**
```go
var ch chan int  // nil channel

// Sending to nil channel: blocks forever
ch <- 1  // Deadlock

// Receiving from nil channel: blocks forever
v := <-ch  // Deadlock

// Useful in select: a nil case is never selected (effectively disabled)
var optionalCh chan int
select {
case v := <-optionalCh:  // Never selected when optionalCh is nil
    fmt.Println(v)
case v := <-otherCh:
    fmt.Println(v)
}
```

### Quick Check
> 1. What is the difference between `chan<- int` and `<-chan int`?
> 2. Who should close a channel — the sender or the receiver?
> 3. What happens when you receive from a closed, empty channel?

---

## 4. Range Over Channel

`for range` over a channel receives values until the channel is closed:

```go
ch := make(chan int)

go func() {
    for i := 0; i < 5; i++ {
        ch <- i
    }
    close(ch)  // MUST close — otherwise for range loops forever
}()

for v := range ch {
    fmt.Println(v)
}
// Prints: 0 1 2 3 4
// Loop exits when ch is closed and empty
```

**Important: if you don't close the channel, `for range` blocks forever (deadlock).**

```go
// Common pattern: generator function returns receive-only channel
func naturals(max int) <-chan int {
    ch := make(chan int)
    go func() {
        for i := 0; i < max; i++ {
            ch <- i
        }
        close(ch)
    }()
    return ch
}

func main() {
    for n := range naturals(5) {
        fmt.Println(n)
    }
}
```

### Quick Check
> 1. What causes a `for range` channel loop to exit?
> 2. What happens if you `for range` a channel that's never closed?
> 3. What is a "generator function" in terms of channels?

---

## 5. Select Statement

`select` waits on multiple channel operations simultaneously — it's like a switch for channels:

```go
select {
case v := <-ch1:
    fmt.Println("received from ch1:", v)
case v := <-ch2:
    fmt.Println("received from ch2:", v)
case ch3 <- 42:
    fmt.Println("sent to ch3")
}
// Blocks until at least ONE case is ready, then executes that case
// If multiple are ready simultaneously: one is chosen randomly
```

**`select` with default** — non-blocking:
```go
select {
case v := <-ch:
    fmt.Println("received:", v)
default:
    fmt.Println("no value ready — continue without waiting")
}
```

**Timeout pattern:**
```go
select {
case result := <-compute():
    fmt.Println("result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("timed out after 5s")
}
```

**Done channel for cancellation:**
```go
func processWithTimeout(input <-chan int, done <-chan struct{}) {
    for {
        select {
        case v, ok := <-input:
            if !ok {
                return  // Input channel closed
            }
            fmt.Println("processing:", v)
        case <-done:
            fmt.Println("cancelled")
            return
        }
    }
}
```

**Merging two channels with select:**
```go
func merge(ch1, ch2 <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for ch1 != nil || ch2 != nil {
            select {
            case v, ok := <-ch1:
                if !ok {
                    ch1 = nil  // Disable this case by setting to nil
                    continue
                }
                out <- v
            case v, ok := <-ch2:
                if !ok {
                    ch2 = nil
                    continue
                }
                out <- v
            }
        }
    }()
    return out
}
```

### Quick Check
> 1. What happens if multiple select cases are ready at the same time?
> 2. How do you make a non-blocking channel receive?
> 3. How do you disable a select case dynamically?

---

## 6. Pipeline Pattern

A pipeline is a series of stages connected by channels — each stage reads from an input channel, does some work, and sends to an output channel:

```go
// Stage 1: Generate numbers
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        for _, n := range nums {
            out <- n
        }
        close(out)
    }()
    return out
}

// Stage 2: Square each number
func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            out <- n * n
        }
        close(out)
    }()
    return out
}

// Stage 3: Filter even numbers only
func filterEven(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        for n := range in {
            if n%2 == 0 {
                out <- n
            }
        }
        close(out)
    }()
    return out
}

func main() {
    // Chain the stages:
    nums := generate(1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    squares := square(nums)
    evens := filterEven(squares)
    
    for v := range evens {
        fmt.Println(v)  // 4 16 36 64 100
    }
}
```

**The pipeline runs concurrently** — while `filterEven` is processing result 1, `square` is already computing result 2, and `generate` is already preparing result 3. All stages run simultaneously.

### Quick Check
> 1. What does each stage in a pipeline do?
> 2. Are pipeline stages sequential or concurrent?
> 3. When does a pipeline stage's goroutine exit?

---

## 7. Fan-Out and Fan-In

**Fan-out** — one channel feeds multiple workers:
```go
func fanOut(input <-chan Job, numWorkers int) []<-chan Result {
    outputs := make([]<-chan Result, numWorkers)
    for i := 0; i < numWorkers; i++ {
        outputs[i] = worker(input)  // All workers read from SAME input channel
    }
    return outputs
}

func worker(jobs <-chan Job) <-chan Result {
    out := make(chan Result)
    go func() {
        defer close(out)
        for job := range jobs {
            out <- process(job)
        }
    }()
    return out
}
```

**Fan-in** — merge multiple channels into one:
```go
func fanIn(channels ...<-chan Result) <-chan Result {
    var wg sync.WaitGroup
    merged := make(chan Result)
    
    // Start a goroutine for each input channel:
    forward := func(ch <-chan Result) {
        defer wg.Done()
        for v := range ch {
            merged <- v
        }
    }
    
    wg.Add(len(channels))
    for _, ch := range channels {
        go forward(ch)
    }
    
    // Close merged when all forwarders are done:
    go func() {
        wg.Wait()
        close(merged)
    }()
    
    return merged
}
```

**Complete fan-out + fan-in example:**
```go
func processURLs(urls []string) []string {
    // Create a jobs channel:
    jobs := make(chan string, len(urls))
    for _, url := range urls {
        jobs <- url
    }
    close(jobs)
    
    // Fan-out: 5 workers
    numWorkers := 5
    workerChans := make([]<-chan string, numWorkers)
    for i := 0; i < numWorkers; i++ {
        workerChans[i] = fetch(jobs)
    }
    
    // Fan-in: collect results
    var results []string
    for result := range fanIn(workerChans...) {
        results = append(results, result)
    }
    return results
}
```

### Quick Check
> 1. In fan-out, do workers each get their own channel or share one?
> 2. Why do you need fan-in after fan-out?
> 3. Who closes the merged channel in a fan-in function?

---

## 8. Done Channel Pattern

The "done channel" signals goroutines to stop — it's the predecessor to `context.Context`:

```go
func work(done <-chan struct{}) {
    for {
        select {
        case <-done:
            fmt.Println("stopping work")
            return
        default:
            doOneUnitOfWork()
        }
    }
}

func main() {
    done := make(chan struct{})
    
    go work(done)
    go work(done)
    go work(done)
    
    time.Sleep(5 * time.Second)
    close(done)  // Close broadcasts to ALL goroutines waiting on done
    
    time.Sleep(100 * time.Millisecond)  // Wait for goroutines to see the close
    fmt.Println("all workers stopped")
}
```

**Why close broadcasts but send doesn't:**
```go
done := make(chan struct{})

// Sending only unblocks ONE receiver:
done <- struct{}{}  // Only one goroutine sees this

// Closing unblocks ALL receivers:
close(done)  // ALL goroutines blocked on <-done are unblocked simultaneously
```

This is why done channels use `close()` for signaling and `chan struct{}` (zero-size type, no allocation):
```go
// struct{}{} has no data — it's just a signal
done := make(chan struct{})
close(done)             // Signal "stop"
<-done                  // Receive the signal — works even after close
```

**Today, prefer `context.Context`** for cancellation (Chapter 20) — done channels are the manual version. But you'll see done channels in older code, and understanding them helps you understand context.

### Quick Check
> 1. Why use `close(done)` instead of `done <- struct{}{}`?
> 2. What happens when you receive from a closed channel?
> 3. Why use `chan struct{}` instead of `chan bool` for done channels?

---

## 9. Channel Pitfalls

**Pitfall 1: Deadlock — all goroutines are waiting:**
```go
// DEADLOCK: main blocks sending, nobody receives
ch := make(chan int)
ch <- 1  // Blocks forever — no goroutine is receiving
<-ch     // Never reached

// Fix: use goroutine or buffered channel
go func() { ch <- 1 }()
v := <-ch
```

**Pitfall 2: Goroutine leak from abandoned channel:**
```go
// LEAK: if timeout fires, the goroutine blocks on send forever
func fetch(url string) <-chan string {
    ch := make(chan string)
    go func() {
        result := httpGet(url)
        ch <- result  // BLOCKS if nobody receives
    }()
    return ch
}

// If caller uses a timeout and stops receiving, goroutine is stuck:
select {
case r := <-fetch(url):
    use(r)
case <-time.After(1 * time.Second):
    return  // fetch goroutine now stuck on send forever
}

// Fix: use buffered channel of 1 so goroutine can always send:
func fetch(url string) <-chan string {
    ch := make(chan string, 1)  // Buffer of 1
    go func() {
        ch <- httpGet(url)  // Never blocks — buffer accepts it
    }()
    return ch
}
```

**Pitfall 3: Closing from the receiver side:**
```go
// WRONG: closing from receiver (the channel's sender might still send!)
go func() {
    for v := range producerChannel {
        process(v)
    }
    close(producerChannel)  // BUG: producer might panic on next send
}()
```

**Pitfall 4: Panic on send to closed channel:**
```go
ch := make(chan int)
close(ch)
ch <- 1  // PANIC: send on closed channel

// Use sync.Once to close exactly once:
var once sync.Once
closeOnce := func() { once.Do(func() { close(ch) }) }
go func() { closeOnce() }()
go func() { closeOnce() }()  // Safe — only closes once
```

---

## Summary

- **Channel**: typed conduit between goroutines; `make(chan T)` / `make(chan T, n)`
- **Unbuffered**: sender blocks until receiver is ready (synchronous rendezvous)
- **Buffered**: sender blocks only when buffer is full (capacity = `n`)
- **Directional**: `chan<- T` (send-only), `<-chan T` (receive-only) — enforce ownership
- **Close**: sender closes; receiving from closed empty channel returns zero value + false
- **`for range`**: reads until channel is closed — always close the channel
- **`select`**: wait on multiple channels; random if multiple ready; `default` = non-blocking
- **Nil channel**: blocks forever — use to disable `select` cases
- **Pipeline**: chain stages via channels — concurrent processing
- **Fan-out**: one input → multiple workers sharing the channel
- **Fan-in**: multiple channels → one merged output channel
- **Done channel**: `close(done)` broadcasts stop to all goroutines
- **Common pitfalls**: deadlock, goroutine leak, closing from receiver, double-close panic

---

## Exercises

### Easy
1. Write a ping-pong program: two goroutines take turns sending "ping" and "pong" to each other via channels. Do 10 rounds, then stop. Use two channels, one for each direction.
2. Build a `timeout` function: `func timeout[T any](fn func() T, d time.Duration) (T, bool)` that runs `fn` in a goroutine and returns its result if it finishes within `d`. Return `(zero, false)` if it times out.
3. Write a producer that generates Fibonacci numbers and sends them to a channel. A consumer reads from the channel and prints each number. Stop after 20 numbers. Use `close` to signal completion.

### Medium
4. Concurrent merge sort: Implement merge sort where the split phase uses goroutines. `MergeSort(data []int) []int` — for slices larger than a threshold (e.g., 1000), sort the two halves concurrently using goroutines and channels. Below the threshold, sort sequentially. Benchmark against sequential sort for slices of 10K, 100K, 1M elements.
5. Rate limiter: Implement a rate limiter using a ticker channel. `NewRateLimiter(rps int) *RateLimiter` creates one that allows `rps` requests per second. `Allow() bool` returns true if a request is allowed now, false if rate limit is exceeded. `Tokens() int` returns current available tokens (leaky bucket model). Test: verify exactly `rps` requests are allowed per second, and excess requests are denied.
6. Event bus: Build a simple pub/sub event bus. `Subscribe(topic string) <-chan Event` returns a channel that receives events on that topic. `Publish(topic string, event Event)` sends an event to all subscribers of that topic. `Unsubscribe(topic string, ch <-chan Event)` removes a subscriber. Handle slow subscribers with a buffer so publishing never blocks. Use a mutex for the subscriber map.

### Hard
7. Circuit breaker channel: Implement a circuit breaker using channels. State machine: Closed → Open → Half-Open → Closed. Track failures in a sliding window. `Execute(fn func() error) error` wraps function calls. When failure rate exceeds threshold in the window, enter Open state (reject calls immediately). After a timeout, enter Half-Open and allow one probe. On success, go Closed; on failure, go Open again. Test all state transitions with concurrent goroutines.
8. Reactive stream: Build a simple reactive stream library (like RxGo but minimal). `Observable[T]` wraps a `<-chan T`. Operations: `Map[T, R any](obs Observable[T], fn func(T) R) Observable[R]`, `Filter[T any](obs Observable[T], pred func(T) bool) Observable[T]`, `Buffer[T any](obs Observable[T], size int) Observable[[]T]` (collect N items then emit as batch), `Merge[T any](a, b Observable[T]) Observable[T]`. All operations must propagate cancellation via a shared done channel. Write a demo pipeline: generate integers 1-100 → filter evens → map to squares → buffer in groups of 5 → print each group.
