# Chapter 20: Select, Timeouts, and Non-Blocking Operations

In the last chapter you met channels and saw a first taste of `select`. Channels alone let a goroutine wait on **one** conversation at a time. But real programs juggle many conversations at once: "give me the next job, *or* shut down if someone cancels, *or* give up if this takes longer than two seconds." The tool for waiting on several channels at once — and for deciding what to do when *none* of them are ready — is `select`. This chapter makes `select` your everyday tool.

## Table of Contents

1. [The select Statement, Properly](#1-the-select-statement-properly)
2. [Random Choice Among Ready Cases](#2-random-choice-among-ready-cases)
3. [Non-Blocking Operations with default](#3-non-blocking-operations-with-default)
4. [Timeouts with time.After](#4-timeouts-with-timeafter)
5. [Tickers and Periodic Work](#5-tickers-and-periodic-work)
6. [The for-select Loop](#6-the-for-select-loop)
7. [Combining Cancellation and Timeout](#7-combining-cancellation-and-timeout)
8. [Nil Channels: Disabling a Case](#8-nil-channels-disabling-a-case)
9. [Common Pitfalls](#9-common-pitfalls)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The select Statement, Properly

A `select` blocks until **one** of its communication cases can proceed, then runs that case. It is like a `switch`, but every case is a channel send or receive:

```go
select {
case v := <-ch1:
    fmt.Println("got from ch1:", v)
case ch2 <- 99:
    fmt.Println("sent 99 to ch2")
case v, ok := <-ch3:
    if !ok {
        fmt.Println("ch3 is closed")
    } else {
        fmt.Println("got from ch3:", v)
    }
}
```

Each case is evaluated only for readiness — a receive case is ready when a value (or a close) is available, a send case is ready when a receiver is waiting (or the buffer has room). `select` picks a ready case and runs it. If nothing is ready, `select` **blocks** until something becomes ready (unless there is a `default` — see §3).

---

## 2. Random Choice Among Ready Cases

If **several** cases are ready at the same moment, `select` chooses one **at random**. This is deliberate: it prevents starvation, where one busy channel would always win and the others would never be served.

```go
func main() {
    a := make(chan string, 1)
    b := make(chan string, 1)
    a <- "from a"
    b <- "from b"

    // Both cases are ready. Over many runs, each wins roughly half the time.
    select {
    case v := <-a:
        fmt.Println(v)
    case v := <-b:
        fmt.Println(v)
    }
}
```

Do not rely on any ordering between cases. If you need priority, you must express it explicitly (for example, try the high-priority channel first with a non-blocking `select`, then fall through to a blocking one).

---

## 3. Non-Blocking Operations with default

Add a `default` case and `select` never blocks: if no other case is ready, `default` runs immediately.

```go
// Non-blocking receive:
select {
case v := <-ch:
    fmt.Println("received", v)
default:
    fmt.Println("nothing available right now")
}

// Non-blocking send (won't block if the buffer is full / no receiver):
select {
case ch <- value:
    fmt.Println("sent")
default:
    fmt.Println("channel full, dropped the value")
}
```

This is how you build things like "drop a metric if the reporting channel is backed up" instead of stalling the whole program.

> ⚠️ A `select` with a `default` inside a tight `for` loop with nothing ready becomes a **busy-wait** — it spins at 100% CPU. Use `default` for genuinely non-blocking single checks, not as the body of a hot loop.

---

## 4. Timeouts with time.After

`time.After(d)` returns a channel that receives the current time after duration `d`. Put it in a `select` and you have a timeout:

```go
func fetch(result <-chan string) {
    select {
    case r := <-result:
        fmt.Println("result:", r)
    case <-time.After(2 * time.Second):
        fmt.Println("timed out after 2s")
    }
}
```

Whichever channel fires first wins. If `result` arrives within two seconds, you print it; otherwise the timeout case runs.

> ⚠️ `time.After` creates a new timer each time it is evaluated, and that timer is not garbage-collected until it fires. In a high-frequency `for-select` loop this leaks timers. For loops, create one `time.NewTimer` (or `time.NewTicker`) outside the loop and reuse it, or on Go 1.23+ note that the timer becomes eligible for GC as soon as it is unreachable. When in doubt, use an explicit timer:

```go
timer := time.NewTimer(2 * time.Second)
defer timer.Stop()

select {
case r := <-result:
    if !timer.Stop() {
        <-timer.C // drain if it already fired
    }
    fmt.Println("result:", r)
case <-timer.C:
    fmt.Println("timed out")
}
```

---

## 5. Tickers and Periodic Work

A `time.Ticker` delivers a value on its channel at a fixed interval — perfect for "do this every N seconds":

```go
func main() {
    ticker := time.NewTicker(500 * time.Millisecond)
    defer ticker.Stop() // ALWAYS stop a ticker, or it leaks

    done := make(chan bool)
    go func() {
        time.Sleep(2 * time.Second)
        done <- true
    }()

    for {
        select {
        case t := <-ticker.C:
            fmt.Println("tick at", t.Format("15:04:05.000"))
        case <-done:
            fmt.Println("stopping")
            return
        }
    }
}
```

Use `time.Tick` (the convenience version) only in short-lived programs — it has no `Stop`, so its underlying ticker leaks for the life of the process.

---

## 6. The for-select Loop

The dominant concurrency shape in Go is a goroutine that loops forever, reacting to whatever channel fires:

```go
func worker(jobs <-chan int, results chan<- int, quit <-chan struct{}) {
    for {
        select {
        case job, ok := <-jobs:
            if !ok {
                return // jobs channel closed — no more work
            }
            results <- job * 2
        case <-quit:
            return // asked to stop
        }
    }
}
```

This one pattern — receive work, or exit on a signal — underlies worker pools, event loops, and background daemons.

---

## 7. Combining Cancellation and Timeout

In real services you rarely hand-roll a `quit` channel; you use `context.Context` (Chapter 22). But the mechanism is the same `select`. A context's `Done()` method returns a channel that closes when the context is cancelled or times out:

```go
func doWork(ctx context.Context, out chan<- int) error {
    for i := 0; ; i++ {
        select {
        case out <- i:
            // sent one unit of work
        case <-ctx.Done():
            return ctx.Err() // context.Canceled or context.DeadlineExceeded
        }
    }
}
```

`ctx.Done()` unifies "the caller cancelled" and "we ran out of time" into a single channel you can select on. This is the idiomatic way to make any blocking operation cancellable.

---

## 8. Nil Channels: Disabling a Case

A send or receive on a `nil` channel blocks **forever**. That sounds useless, but it is a precise tool: setting a channel variable to `nil` **disables** its `select` case, because a case that can never proceed is never chosen.

```go
func merge(a, b <-chan int, out chan<- int) {
    for a != nil || b != nil {
        select {
        case v, ok := <-a:
            if !ok {
                a = nil // stop selecting on a once it's drained
                continue
            }
            out <- v
        case v, ok := <-b:
            if !ok {
                b = nil
                continue
            }
            out <- v
        }
    }
    close(out)
}
```

Without the `nil` trick, a closed channel would keep returning `(zero, false)` instantly, spinning the loop. Setting it to `nil` cleanly removes it from consideration.

---

## 9. Common Pitfalls

- **`select {}` with no cases blocks forever.** Occasionally used as `select {}` to park the main goroutine, but usually a bug.
- **`default` in a hot loop** turns `select` into a CPU-burning spin. Add a timeout or ticker case instead.
- **Leaking timers/tickers.** Always `defer timer.Stop()` / `ticker.Stop()`. Reuse timers inside loops.
- **Forgetting the `ok` on receive.** A closed channel is always "ready" and yields the zero value; without checking `ok` you can loop forever on a dead channel.
- **Assuming case order matters.** It does not. Ready cases are chosen at random.

---

## Summary

- `select` waits on multiple channel operations and runs one ready case; if several are ready it picks at random.
- `default` makes `select` non-blocking — great for "try, else move on," dangerous in tight loops.
- `time.After`, `time.NewTimer`, and `time.NewTicker` bring time into `select` for timeouts and periodic work; always stop tickers and reused timers.
- The `for-select` loop is the backbone of workers and event loops; `ctx.Done()` is the idiomatic cancellation/timeout signal.
- A `nil` channel disables its case — a clean way to retire a drained input in a merge.

Next chapter covers the `sync` package — mutexes, wait groups, and the tools you reach for when channels are not the right fit.

---

## Exercises

1. Write a function that reads from a channel but gives up and returns an error if no value arrives within 100ms. Test it with both a fast and a slow producer.
2. Build a rate limiter: using a ticker, allow at most one "request" to proceed every 200ms, printing each as it goes through.
3. Implement the `merge` function from §8 and prove it terminates when both inputs are closed. Then break it by removing the `a = nil` line and observe the CPU spin.
4. Write a non-blocking metrics sink: a function `report(ch chan<- Metric, m Metric)` that drops `m` (and increments a dropped counter) if the channel is full instead of blocking.
