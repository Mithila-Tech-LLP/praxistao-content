# Chapter 12: Channels — Patterns, Buffered vs Unbuffered, Select

Channels are Go's mechanism for goroutines to communicate and synchronize. Understanding channels deeply — their semantics, patterns, and when NOT to use them — is essential for senior Go interviews.

## Table of Contents

1. [Channel Fundamentals](#1-channel-fundamentals)
2. [Buffered vs Unbuffered Channels](#2-buffered-vs-unbuffered-channels)
3. [Channel Direction Types](#3-channel-direction-types)
4. [The Select Statement](#4-the-select-statement)
5. [Core Channel Patterns](#5-core-channel-patterns)
6. [When NOT to Use Channels](#6-when-not-to-use-channels)
7. [Channel Axioms — What Every Gopher Must Know](#7-channel-axioms--what-every-gopher-must-know)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Channel Fundamentals

A channel is a typed conduit through which goroutines communicate. Channels provide both communication AND synchronization — sending and receiving are synchronization points.

```go
// Create a channel of ints (unbuffered)
ch := make(chan int)

// Create a buffered channel of strings with capacity 5
bufCh := make(chan string, 5)

// Send a value (blocks until receiver is ready, for unbuffered)
ch <- 42

// Receive a value (blocks until sender sends, for unbuffered)
val := <-ch

// Receive with ok to detect if channel is closed
val, ok := <-ch
if !ok {
    fmt.Println("channel closed")
}

// Close a channel: signals no more values will be sent
close(ch)

// Range over a channel: receives until closed
for v := range ch {
    fmt.Println(v) // exits when ch is closed and drained
}
```

---

## 2. Buffered vs Unbuffered Channels

### Unbuffered Channel — Synchronous Rendezvous

```go
ch := make(chan int) // capacity = 0

// Send blocks until a receiver is ready.
// Receive blocks until a sender sends.
// They must meet at the same time — a "rendezvous".

go func() {
    ch <- 42 // blocks here until main receives
}()
fmt.Println(<-ch) // unblocks the sender when we receive
```

**Guarantee:** When a receive from an unbuffered channel returns, you know the sender has already sent. The goroutines synchronized at that point.

### Buffered Channel — Asynchronous up to Capacity

```go
ch := make(chan int, 3) // capacity = 3

// Send does NOT block until the buffer is full
ch <- 1  // non-blocking (buffer: [1])
ch <- 2  // non-blocking (buffer: [1, 2])
ch <- 3  // non-blocking (buffer: [1, 2, 3])
ch <- 4  // BLOCKS: buffer full, wait for a receiver

// Receive does NOT block until the buffer is empty
v := <-ch // non-blocking if buffer has items (v=1, buffer: [2, 3])
```

### Choosing Buffer Size

```go
// Capacity 1: useful for non-blocking signal with memory of one value
done := make(chan struct{}, 1)
// A goroutine can send "done" without blocking, even if receiver isn't ready yet

// Capacity N: allows producer to be N steps ahead of consumer (bounded buffer)
// Use when: you know the maximum burst size and want to prevent backpressure up to that size
workQueue := make(chan Job, 100)

// Common mistake: using a large buffer to "avoid blocking"
// This hides backpressure problems — the buffer fills up and you block anyway,
// but now you also have 1000 unprocessed items in memory
```

---

## 3. Channel Direction Types

You can restrict a channel to send-only or receive-only. Use this in function signatures to document intent and prevent misuse.

```go
// Bidirectional channel (created with make)
ch := make(chan int)

// Send-only channel: can only send to it
func producer(ch chan<- int) {
    ch <- 42     // OK
    // <-ch      // COMPILE ERROR: cannot receive from send-only channel
}

// Receive-only channel: can only receive from it
func consumer(ch <-chan int) {
    v := <-ch    // OK
    // ch <- 42  // COMPILE ERROR: cannot send on receive-only channel
}

// Bidirectional channels are implicitly convertible to directional
var sendOnly chan<- int = ch
var recvOnly <-chan int = ch

// Good pattern: return receive-only to callers
func startWorker() <-chan Result {
    ch := make(chan Result) // internally bidirectional
    go func() {
        // do work
        ch <- Result{...}
        close(ch)
    }()
    return ch // caller can only receive — cannot accidentally send or close
}
```

---

## 4. The Select Statement

`select` is Go's way to wait on multiple channel operations simultaneously. It's a multiplexer for channels.

```go
// Basic select
select {
case v := <-ch1:
    fmt.Println("received from ch1:", v)
case v := <-ch2:
    fmt.Println("received from ch2:", v)
case ch3 <- 42:
    fmt.Println("sent to ch3")
}
// If multiple cases are ready, select picks one at RANDOM.

// Non-blocking with default
select {
case v := <-ch:
    fmt.Println("got value:", v)
default:
    fmt.Println("no value ready") // executes immediately if ch has nothing
}

// Timeout pattern
select {
case result := <-workCh:
    fmt.Println("got result:", result)
case <-time.After(5 * time.Second):
    fmt.Println("timed out")
}
```

### Select with Done Channel (Cancellation)

```go
func worker(done <-chan struct{}, jobs <-chan Job) {
    for {
        select {
        case job := <-jobs:
            process(job)
        case <-done:
            return // stop when done signal received
        }
    }
}

// Usage:
done := make(chan struct{})
go worker(done, jobs)
// ...later...
close(done) // broadcast cancellation to all goroutines listening on done
```

---

## 5. Core Channel Patterns

### Pattern 1: Generator — Produce Values on Demand

```go
func fibonacci() <-chan int {
    ch := make(chan int)
    go func() {
        defer close(ch)
        a, b := 0, 1
        for {
            ch <- a
            a, b = b, a+b
        }
    }()
    return ch
}

// Usage:
gen := fibonacci()
for i := 0; i < 10; i++ {
    fmt.Println(<-gen)
}
```

### Pattern 2: Pipeline — Chain Processing Stages

```go
// Each stage: receives from input channel, sends to output channel
func generate(nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums { out <- n }
    }()
    return out
}

func square(in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in { out <- n * n }
    }()
    return out
}

// Chain: generate → square → print
c := generate(2, 3, 4)
sq := square(c)
for v := range sq { fmt.Println(v) } // 4, 9, 16
```

### Pattern 3: Fan-Out — Distribute Work

```go
func fanOut(in <-chan Job, n int) []<-chan Result {
    outputs := make([]<-chan Result, n)
    for i := 0; i < n; i++ {
        outputs[i] = worker(in) // all workers read from same input
    }
    return outputs
}

// Multiple workers consume from the same channel
// Go's channel is already safe for concurrent readers
```

### Pattern 4: Fan-In — Merge Multiple Channels

```go
func fanIn(channels ...<-chan int) <-chan int {
    merged := make(chan int)
    var wg sync.WaitGroup

    output := func(c <-chan int) {
        defer wg.Done()
        for v := range c { merged <- v }
    }

    wg.Add(len(channels))
    for _, ch := range channels { go output(ch) }

    go func() {
        wg.Wait()
        close(merged) // close when all inputs are done
    }()

    return merged
}
```

### Pattern 5: Done Channel for Cancellation

```go
// Propagate cancellation down the call tree with a done channel
// (Prefer context.Context in production code — same concept, more features)
func doWork(done <-chan struct{}) <-chan Result {
    results := make(chan Result)
    go func() {
        defer close(results)
        for {
            select {
            case <-done:
                return
            default:
                result := compute() // may take a while
                select {
                case results <- result:
                case <-done:
                    return
                }
            }
        }
    }()
    return results
}
```

### Pattern 6: Semaphore via Buffered Channel

```go
// Limit to maxN concurrent operations using a buffered channel
func withSemaphore(maxN int, work func()) {
    sem := make(chan struct{}, maxN)
    for _, task := range tasks {
        sem <- struct{}{}  // acquire: blocks when maxN goroutines are active
        go func(t Task) {
            defer func() { <-sem }() // release
            work(t)
        }(task)
    }
}
```

---

## 6. When NOT to Use Channels

Rob Pike's guideline: "Don't communicate by sharing memory; share memory by communicating." But the converse is also important: **don't use channels when a mutex is simpler**.

```go
// BAD: using a channel when a mutex is clearer
type Counter struct {
    ch chan int
    n  int
}
func (c *Counter) Increment() {
    c.ch <- 1 // awkward channel ping to trigger increment
}

// GOOD: use a mutex for shared state
type Counter struct {
    mu sync.Mutex
    n  int
}
func (c *Counter) Increment() {
    c.mu.Lock()
    c.n++
    c.mu.Unlock()
}

// USE CHANNELS WHEN:
// - Passing ownership of data between goroutines
// - Distributing work units across goroutines
// - Signaling events (done, cancel, timer)
// - Building pipelines

// USE MUTEX WHEN:
// - Protecting shared state that multiple goroutines read/write
// - Coordinating access to a cache or pool
// - Simple counters and flags
```

---

## 7. Channel Axioms — What Every Gopher Must Know

These rules determine goroutine deadlocks. Know them cold.

```go
// A nil channel blocks forever on both send and receive
var ch chan int  // nil
ch <- 1          // blocks forever
<-ch             // blocks forever
// Useful: close a channel by setting it to nil in select to disable that case

// A closed channel:
close(ch)
<-ch   // returns zero value immediately, ok=false
ch <- 1 // PANIC: send on closed channel

// RULE: Only the sender should close a channel.
// Multiple senders: use sync.WaitGroup + a closer goroutine.

// Receiving from a buffered closed channel drains remaining values:
ch := make(chan int, 3)
ch <- 1; ch <- 2; ch <- 3
close(ch)
for v := range ch { fmt.Println(v) } // 1, 2, 3 (then exits)
```

| Operation | Nil channel | Open channel | Closed channel |
|---|---|---|---|
| Send | Blocks forever | Sends or blocks | **PANIC** |
| Receive | Blocks forever | Receives or blocks | Returns zero value |
| Close | **PANIC** | Closes | **PANIC** |

---

## 8. Interview Questions & Model Answers

**Q: What is the difference between a buffered and unbuffered channel?**

"An unbuffered channel is synchronous — a send blocks until a receiver is ready, and a receive blocks until a sender sends. They meet at the channel. A buffered channel decouples sender and receiver up to the buffer capacity — a sender can send without a receiver being ready, as long as the buffer isn't full. Unbuffered channels provide a stronger synchronization guarantee: when the receive returns, you know the sender has already executed the send."

**Q: What happens if you send to a closed channel?**

"Panic. That's why only the sender should close a channel — the sender knows when there are no more values. If multiple goroutines send on the same channel, you coordinate with a sync.WaitGroup: all senders Done(), and a separate goroutine closes the channel after Wait()."

**Q: When would you use select with a default case?**

"When you want non-blocking channel operations. If no channel is ready, instead of blocking, the default case executes immediately. This is useful for polling patterns, implementing try-send/try-receive, or for worker goroutines that should do something else if no work is available rather than waiting."

**Q: Channels vs mutex — when do you use which?**

"Channels when transferring data ownership between goroutines, distributing work, or signaling events — the data moves from one goroutine to another. Mutex when protecting shared mutable state that multiple goroutines access — the data stays in place, you just control who touches it at a time. Overusing channels for state protection leads to more complex code than a simple mutex would."

---

## Summary

- Unbuffered channels synchronize sender and receiver — both must be present simultaneously.
- Buffered channels allow sends up to capacity before blocking — decouples sender and receiver.
- Channel direction types (`chan<-`, `<-chan`) document intent and prevent misuse at compile time.
- `select` multiplexes channel operations — picks randomly among ready cases.
- Core patterns: generator, pipeline, fan-out, fan-in, done channel, semaphore.
- Nil channel blocks forever (useful in select to disable a case). Closed channel receive returns zero value; closed channel send panics.
- Only senders close channels. Only use channels when ownership transfer or signaling is the intent; use mutex for protecting shared state.
