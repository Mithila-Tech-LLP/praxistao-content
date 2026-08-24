# Chapter 20: Select, Timeouts, and Non-Blocking Operations

In the last chapter you met `select` briefly — the switch statement for channels. This chapter makes you fluent in it. `select` is how real Go programs deal with time: timing out slow operations, running periodic work, detecting stalled goroutines, and doing "try it, but don't wait" channel operations. Almost every production Go service has a `select` loop at its heart, and almost every subtle concurrency bug lives inside one. Master this chapter and the `sync` package (next chapter) becomes the easy part.

## Table of Contents

1. [Select — The Complete Rules](#1-select--the-complete-rules)
2. [Non-Blocking Operations with default](#2-non-blocking-operations-with-default)
3. [Timers — time.After and time.NewTimer](#3-timers--timeafter-and-timenewtimer)
4. [Tickers — Periodic Work](#4-tickers--periodic-work)
5. [Timeout Patterns](#5-timeout-patterns)
6. [Heartbeats — Detecting Stalled Goroutines](#6-heartbeats--detecting-stalled-goroutines)
7. [The Select Loop — Heart of a Server](#7-the-select-loop--heart-of-a-server)
8. [Common Mistakes](#8-common-mistakes)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Select — The Complete Rules

`select` blocks until one of its cases can proceed, then executes that case. Here are ALL the rules, precisely:

```go
select {
case v := <-ch1:        // Case A: receive
    fmt.Println("got", v)
case ch2 <- 42:         // Case B: send
    fmt.Println("sent")
case v, ok := <-ch3:    // Case C: receive with closed-check
    fmt.Println(v, ok)
default:                // Optional: runs if NO case is ready
    fmt.Println("nothing ready")
}
```

**Rule 1 — All channel expressions are evaluated first, top to bottom.** Before blocking, Go evaluates every channel operand and every value being sent:

```go
select {
case ch1 <- computeValue():  // computeValue() runs even if this case never fires!
    // ...
case v := <-getChannel():    // getChannel() runs even if this case never fires!
    _ = v
}
```

Keep case expressions cheap and side-effect free — put expensive work inside the case body, not in the case header.

**Rule 2 — If multiple cases are ready, one is chosen uniformly at random.** This prevents starvation:

```go
// Both channels have values waiting:
ch1 := make(chan string, 1)
ch2 := make(chan string, 1)
ch1 <- "one"
ch2 <- "two"

select {
case v := <-ch1:
    fmt.Println(v)  // ~50% of runs
case v := <-ch2:
    fmt.Println(v)  // ~50% of runs
}
```

If Go always picked the first ready case, a busy `ch1` could starve `ch2` forever. Random choice guarantees fairness. **Never rely on case order for priority** — we'll see the correct priority pattern in section 7.

**Rule 3 — With no ready case and no `default`, select blocks.** With zero cases at all, it blocks forever:

```go
select {}  // Blocks forever — sometimes used to park main() in demos
```

**Rule 4 — A closed channel is always ready to receive.** Receiving from a closed channel returns the zero value immediately, so a closed channel makes its case fire instantly, every time:

```go
done := make(chan struct{})
close(done)

select {
case <-done:
    fmt.Println("always selected — done is closed")  // Fires immediately
case <-time.After(time.Second):
    fmt.Println("never happens")
}
```

This is exactly why `close(done)` works as a broadcast: every `select` waiting on `<-done` unblocks at once.

**Rule 5 — A nil channel case is never ready.** Sending or receiving on `nil` blocks forever, so a nil channel silently disables its case:

```go
var maybe chan int  // nil — this case is switched OFF
select {
case v := <-maybe:  // Never fires
    fmt.Println(v)
case v := <-realCh:
    fmt.Println(v)
}
```

Rules 4 and 5 together give you a switch you can flip at runtime: set a channel variable to nil to disable a case, restore it to re-enable. You saw this in the `merge` function last chapter; section 7 pushes it further.

### Quick Check
> 1. If two cases are ready at the same moment, which one runs?
> 2. Why should you avoid calling expensive functions in a `select` case header?
> 3. What is the difference between a closed channel and a nil channel inside `select`?

---

## 2. Non-Blocking Operations with default

Adding `default` turns `select` from "wait for a channel" into "check a channel" — the operation happens only if it can proceed *right now*.

**Try-receive:**
```go
select {
case v := <-ch:
    fmt.Println("got:", v)
default:
    fmt.Println("channel empty — moving on")
}
```

**Try-send:**
```go
select {
case ch <- v:
    // Delivered
default:
    // Channel full (or no receiver ready) — drop, log, or fall back
    metrics.Dropped.Add(1)
}
```

**Why try-send matters: never block the hot path.** Imagine an HTTP handler that publishes an event after each request. If the event channel fills up, do you want every user request to hang? Usually not — dropping an event beats stalling the service:

```go
type EventBus struct {
    events chan Event
}

// Publish never blocks. Returns false if the event was dropped.
func (b *EventBus) Publish(e Event) bool {
    select {
    case b.events <- e:
        return true
    default:
        return false  // Buffer full — caller decides: drop, log, count
    }
}
```

**Drop-oldest instead of drop-newest.** Sometimes the newest value matters most (e.g., latest sensor reading). Evict the oldest to make room:

```go
// PublishLatest keeps the freshest values when the buffer is full.
func (b *EventBus) PublishLatest(e Event) {
    for {
        select {
        case b.events <- e:
            return
        default:
            // Buffer full: discard one old event, then retry the send.
            select {
            case <-b.events:
            default:
            }
        }
    }
}
```

**Checking if work was cancelled without stopping:**
```go
func worker(done <-chan struct{}) {
    for {
        select {
        case <-done:
            return  // Stop requested
        default:
            // Not cancelled — do one unit of work, then check again
        }
        doOneUnitOfWork()
    }
}
```

**The busy-wait trap.** `default` makes `select` return instantly — so a bare polling loop spins the CPU at 100%:

```go
// BAD: burns a full CPU core doing nothing
for {
    select {
    case v := <-ch:
        process(v)
    default:
        // Runs millions of times per second!
    }
}

// GOOD: just block — that's what channels are for
for v := range ch {
    process(v)
}
```

Use `default` when you have *something else to do* if the channel isn't ready. If the answer to "what do I do instead?" is "nothing, just check again" — remove the `default` and block.

### Quick Check
> 1. What does `default` change about `select`'s behavior?
> 2. When is dropping a message better than blocking?
> 3. Why is `for { select { ... default: } }` with an empty default a bug?

---

## 3. Timers — time.After and time.NewTimer

A timer is a channel that receives exactly one value after a duration. The quick version:

```go
// time.After returns <-chan Time that fires once after d:
select {
case v := <-ch:
    fmt.Println("got:", v)
case <-time.After(2 * time.Second):
    fmt.Println("gave up after 2s")
}
```

`time.After` is perfect for one-shot timeouts. But it hides a `time.Timer` you can't stop — and that used to matter a lot:

```go
// CAUTION (pre-Go 1.23): a NEW timer is created on EVERY loop iteration.
// If messages arrive every millisecond, you allocate 60,000 one-minute
// timers per minute, and none can be freed until it fires.
for {
    select {
    case msg := <-messages:
        process(msg)
    case <-time.After(time.Minute):
        return  // Idle timeout
    }
}
```

Since **Go 1.23**, unstopped timers are garbage-collected as soon as nothing references them, so this pattern no longer leaks memory — but it still allocates a fresh timer per iteration. For a hot loop, reuse one timer:

```go
idle := time.NewTimer(time.Minute)
defer idle.Stop()

for {
    select {
    case msg := <-messages:
        process(msg)
        // Rearm the timer for another minute of idle allowance:
        if !idle.Stop() {
            select { // Drain a value that may have already fired (pre-1.23 safety)
            case <-idle.C:
            default:
            }
        }
        idle.Reset(time.Minute)
    case <-idle.C:
        fmt.Println("no messages for 1 minute — shutting down")
        return
    }
}
```

**The `time.Timer` API:**
```go
t := time.NewTimer(5 * time.Second)

<-t.C            // Wait for it to fire (channel receives the fire time)

t.Stop()         // Cancel — returns false if it already fired or was stopped
t.Reset(2 * time.Second)  // Rearm with a new duration

// Fire a callback instead of a channel:
t2 := time.AfterFunc(5*time.Second, func() {
    fmt.Println("runs in its own goroutine after 5s")
})
t2.Stop()  // Cancel the callback if it hasn't run yet
```

**Go version note.** In Go 1.23+ the `Stop`-then-drain dance became unnecessary: `Reset` and `Stop` now guarantee no stale value is left in the channel. The drain pattern above is still what you'll see in most existing codebases (and it's harmless on new versions), so learn to read it — but on Go 1.23+ a plain `idle.Reset(d)` is correct.

**`time.Sleep` vs a timer:** `time.Sleep(d)` blocks the goroutine unconditionally — nothing can interrupt it. A timer inside `select` is an *interruptible* sleep:

```go
// Uninterruptible — cancel signal is ignored for the full 10 minutes:
time.Sleep(10 * time.Minute)

// Interruptible — wakes early if done closes:
select {
case <-time.After(10 * time.Minute):
case <-done:
}
```

### Quick Check
> 1. What does the channel returned by `time.After(d)` deliver, and how many times?
> 2. Why was `time.After` inside a loop a problem before Go 1.23?
> 3. How do you make a "sleep" that a cancellation signal can interrupt?

---

## 4. Tickers — Periodic Work

A ticker fires repeatedly at an interval — the tool for heartbeats, metrics flushing, cache cleanup, and polling:

```go
ticker := time.NewTicker(1 * time.Second)
defer ticker.Stop()  // ALWAYS stop tickers — they never fire "once and done"

done := make(chan struct{})
go func() {
    time.Sleep(5 * time.Second)
    close(done)
}()

for {
    select {
    case t := <-ticker.C:
        fmt.Println("tick at", t.Format("15:04:05"))
    case <-done:
        fmt.Println("stopping")
        return
    }
}
```

**Tickers drop ticks when you're slow.** The ticker's channel has a buffer of 1. If your work takes longer than the interval, missed ticks are discarded — the ticker does *not* queue them up and burst-fire later:

```go
ticker := time.NewTicker(100 * time.Millisecond)
defer ticker.Stop()

for range ticker.C {
    time.Sleep(250 * time.Millisecond)  // Slower than the interval
    // You get a tick roughly every 250ms, not a backlog of 100ms ticks.
    // This "skip when behind" behavior is usually exactly what you want.
}
```

**A background maintenance goroutine — the classic shape:**
```go
type Cache struct {
    mu    sync.Mutex
    items map[string]entry
    stop  chan struct{}
}

func (c *Cache) StartJanitor(interval time.Duration) {
    go func() {
        ticker := time.NewTicker(interval)
        defer ticker.Stop()
        for {
            select {
            case <-ticker.C:
                c.evictExpired()
            case <-c.stop:
                return
            }
        }
    }()
}

func (c *Cache) Close() {
    close(c.stop)  // Janitor exits; ticker is stopped by its defer
}
```

**Avoid `time.Tick` in long-lived code you can't clean up.** `time.Tick(d)` is shorthand that gives you only the channel, with no way to call `Stop()`. Before Go 1.23 that meant the ticker ran (and consumed resources) forever. Post-1.23 it's collectible, but `NewTicker` + `defer Stop()` remains the habit worth building.

**Ticker vs sleeping in a loop:**
```go
// Sleep loop: interval measured AFTER each job → drifts by the job's duration
for {
    doJob()                      // Takes 2s
    time.Sleep(10 * time.Second) // Next job starts 12s after the previous one
}

// Ticker: fires on a fixed schedule regardless of job duration (skipping if behind)
ticker := time.NewTicker(10 * time.Second)
defer ticker.Stop()
for range ticker.C {
    doJob()  // Starts every ~10s as long as doJob takes < 10s
}
```

### Quick Check
> 1. What happens to ticks you don't receive in time?
> 2. Why must you call `ticker.Stop()`?
> 3. When does a sleep-loop schedule drift compared to a ticker?

---

## 5. Timeout Patterns

### Pattern 1: Timeout a single operation

```go
func fetchWithTimeout(url string, timeout time.Duration) (string, error) {
    result := make(chan string, 1)  // Buffered! See note below.
    go func() {
        result <- slowFetch(url)
    }()

    select {
    case r := <-result:
        return r, nil
    case <-time.After(timeout):
        return "", fmt.Errorf("fetch %s: timed out after %v", url, timeout)
    }
}
```

**The buffer of 1 is load-bearing.** If the timeout fires first, nobody will ever receive from `result`. With an unbuffered channel the goroutine would block on its send *forever* — a goroutine leak (Pitfall 2 from last chapter). With a buffer of 1 the goroutine deposits its result and exits cleanly; the abandoned value is garbage-collected.

**A timeout does not stop the work.** The goroutine keeps running `slowFetch` to completion — you've stopped *waiting*, not stopped *working*. Truly cancelling in-flight work requires cooperation from the work itself, which is exactly what `context.Context` provides (two chapters from now).

### Pattern 2: A reusable generic timeout wrapper

```go
var ErrTimeout = errors.New("operation timed out")

func WithTimeout[T any](d time.Duration, fn func() (T, error)) (T, error) {
    type outcome struct {
        val T
        err error
    }
    ch := make(chan outcome, 1)
    go func() {
        v, err := fn()
        ch <- outcome{v, err}
    }()

    select {
    case o := <-ch:
        return o.val, o.err
    case <-time.After(d):
        var zero T
        return zero, ErrTimeout
    }
}

// Usage:
user, err := WithTimeout(2*time.Second, func() (User, error) {
    return db.LoadUser(42)
})
```

### Pattern 3: Total deadline across multiple operations

A per-operation timeout lets three 4-second calls take 12 seconds total. For a whole-job budget, create the timer once and share it:

```go
func processAll(items []Item, budget time.Duration) error {
    deadline := time.After(budget)  // ONE timer for the entire batch

    for _, item := range items {
        result := make(chan error, 1)
        go func() { result <- process(item) }()

        select {
        case err := <-result:
            if err != nil {
                return err
            }
        case <-deadline:
            return fmt.Errorf("budget of %v exhausted", budget)
        }
    }
    return nil
}
```

### Pattern 4: First response wins

Query several equivalent sources and take whichever answers first — a classic latency-reduction trick (Google calls these "hedged requests"):

```go
func fastest(replicas ...func() string) string {
    // Buffer = len(replicas): every loser can deposit its result
    // and exit. No goroutine leaks.
    ch := make(chan string, len(replicas))
    for _, replica := range replicas {
        go func() { ch <- replica() }()
    }
    return <-ch  // First answer wins; the rest land in the buffer and get GC'd
}

// Usage:
answer := fastest(
    func() string { return queryMirror("eu-west") },
    func() string { return queryMirror("us-east") },
    func() string { return queryMirror("ap-south") },
)
```

### Quick Check
> 1. Why must the result channel in a timeout pattern be buffered?
> 2. Does timing out an operation stop the underlying goroutine?
> 3. In "first response wins", why is the buffer sized to the number of replicas?

---

## 6. Heartbeats — Detecting Stalled Goroutines

A timeout tells you an operation is slow. A **heartbeat** tells you a goroutine is *alive* — even when it legitimately has nothing to report. Long-running workers send a pulse at an interval; a monitor treats silence as failure:

```go
// worker processes jobs and emits a heartbeat every interval.
func worker(jobs <-chan int, interval time.Duration) (<-chan int, <-chan struct{}) {
    results := make(chan int)
    heartbeat := make(chan struct{}, 1)  // Buffer 1: never block on the pulse

    go func() {
        defer close(results)
        pulse := time.NewTicker(interval)
        defer pulse.Stop()

        for {
            select {
            case j, ok := <-jobs:
                if !ok {
                    return
                }
                results <- j * j  // The "work"
            case <-pulse.C:
                select {
                case heartbeat <- struct{}{}:
                default:  // Monitor hasn't consumed the last pulse — skip, don't block
                }
            }
        }
    }()
    return results, heartbeat
}
```

The monitor waits for results OR pulses, and declares the worker dead after prolonged silence:

```go
func monitor(results <-chan int, heartbeat <-chan struct{}, interval time.Duration) {
    // Allow two missed pulses before declaring death:
    timeout := 2 * interval

    for {
        select {
        case r, ok := <-results:
            if !ok {
                fmt.Println("worker finished cleanly")
                return
            }
            fmt.Println("result:", r)
        case <-heartbeat:
            // Worker is alive — nothing to do, just reset the wait
        case <-time.After(timeout):
            fmt.Println("worker stalled! restarting...")
            return
        }
    }
}
```

Why this works: every arm of the monitor's `select` restarts the `time.After(timeout)` on the next loop iteration. As long as *something* — a result or a pulse — arrives within `timeout`, the deadline never fires. If the worker deadlocks, gets stuck on a blocking call, or its goroutine dies, the pulses stop and the monitor notices within two intervals.

Note the deliberate design in the worker: the heartbeat send uses try-send (`select`/`default`) so a slow monitor can never back-pressure the worker, and the pulse is skipped while the worker is blocked sending a result — meaning a worker stuck on a downstream consumer *also* reads as stalled, which is usually what you want.

This trio — ticker for pulses, try-send for delivery, `time.After` for silence detection — is the standard Go liveness pattern. You'll meet it again scaled up as Kubernetes liveness probes and consumer-group session timeouts in Kafka.

### Quick Check
> 1. What is the difference between a timeout and a heartbeat?
> 2. Why does the worker use a non-blocking send for its pulse?
> 3. Why is the monitor's timeout set to a multiple of the pulse interval?

---

## 7. The Select Loop — Heart of a Server

Put everything together and you get the **select loop** (also called an event loop or actor loop): one goroutine that owns some state and serially handles commands, ticks, and shutdown. No mutexes needed — the loop is the only code touching the state:

```go
type Counter struct {
    incr chan int          // Commands in
    read chan chan int     // Requests carrying a reply channel
    quit chan struct{}
}

func NewCounter() *Counter {
    c := &Counter{
        incr: make(chan int),
        read: make(chan chan int),
        quit: make(chan struct{}),
    }
    go c.loop()
    return c
}

func (c *Counter) loop() {
    total := 0                                // Owned by this goroutine ONLY
    ticker := time.NewTicker(10 * time.Second)
    defer ticker.Stop()

    for {
        select {
        case n := <-c.incr:
            total += n
        case reply := <-c.read:
            reply <- total                    // Answer on the caller's channel
        case <-ticker.C:
            fmt.Println("periodic snapshot:", total)
        case <-c.quit:
            fmt.Println("final total:", total)
            return
        }
    }
}

func (c *Counter) Increment(n int) { c.incr <- n }

func (c *Counter) Value() int {
    reply := make(chan int)
    c.read <- reply
    return <-reply
}

func (c *Counter) Close() { close(c.quit) }
```

**Priority between cases.** `select` picks randomly among ready cases, so how do you say "always check shutdown first"? Nest a non-blocking check:

```go
for {
    // High priority: check quit before anything else
    select {
    case <-c.quit:
        return
    default:
    }

    // Normal priority: block on all events (including quit, so we
    // still wake up if it closes while we're idle)
    select {
    case n := <-c.incr:
        total += n
    case <-c.quit:
        return
    }
}
```

**Dynamic cases with nil channels.** A select loop that forwards buffered values only *when it has any* — enabling and disabling its send case on the fly:

```go
// bridge receives from in without ever blocking the sender,
// buffering internally, and forwards to out when out is ready.
func bridge(in <-chan int, out chan<- int) {
    var buffer []int
    for in != nil || len(buffer) > 0 {
        // Enable the send case only when there's something to send:
        var sendCh chan<- int  // nil → send case disabled
        var next int
        if len(buffer) > 0 {
            sendCh = out
            next = buffer[0]
        }

        select {
        case v, ok := <-in:
            if !ok {
                in = nil  // Disable the receive case; drain the buffer
                continue
            }
            buffer = append(buffer, v)
        case sendCh <- next:
            buffer = buffer[1:]
        }
    }
    close(out)
}
```

Trace it: when `buffer` is empty, `sendCh` is nil, so the loop can only receive. When `in` is closed, we set it to nil, so the loop can only send. Two channel variables act as switches that reshape the `select` every iteration. This is the most advanced idiom in this chapter — read it twice; it's worth it.

**Graceful drain on shutdown.** When quitting, handle whatever is already queued, then leave:

```go
case <-c.quit:
    for {
        select {
        case n := <-c.incr:
            total += n       // Handle stragglers already in flight
        default:
            return           // Queue empty — now we can exit
        }
    }
```

### Quick Check
> 1. Why does a select loop not need a mutex to protect its state?
> 2. How does a caller get a value OUT of a select loop?
> 3. How do you give one case priority over the others?

---

## 8. Common Mistakes

**Mistake 1: `default` that busy-waits** — covered in section 2. If the `default` body doesn't do real work, delete it and block.

**Mistake 2: Timeout that resets when you meant a deadline:**
```go
// BUG: each iteration gets a FRESH 5 seconds — the loop can run forever
for {
    select {
    case v := <-ch:
        process(v)
    case <-time.After(5 * time.Second):
        return
    }
}

// FIX (if you wanted a total budget): create the timer once
deadline := time.After(5 * time.Second)
for {
    select {
    case v := <-ch:
        process(v)
    case <-deadline:
        return
    }
}
```
Both versions are legitimate — one is an *idle* timeout, the other a *total* deadline. The bug is not knowing which one you wrote.

**Mistake 3: Forgetting that a closed channel spins the loop:**
```go
// BUG: once ch closes, this case is ALWAYS ready — infinite hot loop
for {
    select {
    case v := <-ch:
        process(v)  // After close: processes zero values at full CPU speed
    case <-done:
        return
    }
}

// FIX: detect the close and disable the case
for {
    select {
    case v, ok := <-ch:
        if !ok {
            ch = nil  // Nil channel: case never fires again
            continue
        }
        process(v)
    case <-done:
        return
    }
}
```

**Mistake 4: Sending in a case header with side effects:**
```go
// BUG: pop() runs (and consumes the item!) even if the send case loses
select {
case out <- queue.pop():
case <-done:
    return  // The popped item is silently lost
}

// FIX: peek first, pop only after the send succeeds
next := queue.peek()
select {
case out <- next:
    queue.pop()
case <-done:
    return
}
```

**Mistake 5: Assuming select gives you ordering.** If ten values are queued across two channels, `select` gives no guarantee about interleaving — only per-channel FIFO order. If cross-channel ordering matters, use one channel.

---

## Summary

- **Evaluation**: all case headers evaluate first (keep them cheap); ready cases are chosen **uniformly at random** — never rely on case order
- **Closed channel**: always ready — the basis of broadcast cancellation; guard with `v, ok` and nil the channel to avoid spin loops
- **Nil channel**: never ready — a runtime switch to disable/enable cases dynamically
- **`default`**: makes select non-blocking — try-send/try-receive; beware busy-wait loops
- **`time.After`**: one-shot timeout channel; fine for single use, allocates per call in loops (reuse `time.NewTimer` + `Reset` in hot paths; Go 1.23+ made stale timers GC-able and `Reset` drain-free)
- **`time.NewTicker`**: periodic firing, drops ticks when you fall behind — always `defer ticker.Stop()`
- **Timeout patterns**: buffered result channel (size 1) prevents goroutine leaks; timeouts stop *waiting*, not *working*; one shared timer = total deadline, timer-per-iteration = idle timeout
- **First-response-wins**: N goroutines, buffer of N, take the first receive
- **Heartbeats**: ticker pulse + try-send + `time.After` silence detection = liveness monitoring
- **Select loop**: one goroutine owning state, fed by command channels — concurrency without mutexes; reply channels get answers out; nested select gives priority

Next chapter: sometimes you *do* want shared memory instead of message passing — mutexes, wait groups, atomic operations, and the rest of the `sync` package toolkit.

---

## Exercises

### Easy
1. Write `TrySend[T any](ch chan T, v T) bool` and `TryReceive[T any](ch chan T) (T, bool)` — non-blocking send and receive helpers using `select` with `default`. Demonstrate both on a full and an empty buffered channel.
2. Write a `countdown(n int)` function that prints n, n-1, ... 1 once per second using a `time.NewTicker`, but stops early (printing "aborted!") if the user presses Enter. Hint: run `fmt.Scanln` in a goroutine that closes an `abort` channel.
3. Build `SleepInterruptible(d time.Duration, cancel <-chan struct{}) bool` that waits for `d` but returns early (with `false`) if `cancel` closes. Return `true` if the full duration elapsed. Verify both paths with short durations.

### Medium
4. Idle-shutdown worker: build a worker goroutine that processes strings from a channel and shuts itself down after 3 seconds without input, using a single reused `time.Timer` with `Stop`/`Reset` (not `time.After`). Feed it bursts of input separated by pauses and log when the idle shutdown triggers.
5. Hedged requests with latency stats: implement `Hedged(primary, backup func() string, hedgeAfter time.Duration) string` — call `primary` immediately; if it hasn't answered within `hedgeAfter`, ALSO start `backup` and return whichever finishes first. Ensure no goroutine leaks (buffered channels). Simulate with random latencies and measure how often the backup wins.
6. Heartbeat supervisor: extend section 6 into a `Supervise(work func(jobs <-chan int) (<-chan int, <-chan struct{}))` function that monitors the heartbeat and **restarts** the worker (up to 3 times) when it stalls. Simulate a stall by making the worker randomly stop pulsing. Log each restart.

### Hard
7. Priority job dispatcher: build a dispatcher goroutine with two input channels, `high chan Job` and `low chan Job`, that always processes ALL queued high-priority jobs before any low-priority job (nested-select priority pattern), supports graceful shutdown that drains both queues (high first), and emits a per-second throughput report from a ticker. Prove with a test that under continuous load, no low job ever runs while a high job is queued.
8. Batching writer: implement `NewBatcher(flush func([]Record), maxSize int, maxDelay time.Duration) chan<- Record` — a select loop that accumulates records and calls `flush` when EITHER the batch reaches `maxSize` OR `maxDelay` has elapsed since the first record of the current batch (not since the last flush — you'll need to arm/disarm a timer, or use the nil-channel trick with a timer channel). Closing the input channel flushes the remainder and stops the loop. Test all three flush triggers: size, delay, and close.
