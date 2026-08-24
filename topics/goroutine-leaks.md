---
title: Goroutine Leaks and How to Prevent Them
category: Software & Programming
tags: [Go, Concurrency]
duration: 8 min read
relatedCourses: [go-programming, senior-engineer-interview]
relatedProjects: [rest-api-server]
relatedTopics: [context-cancellation-patterns, worker-pool-patterns]
---

## TL;DR

- A goroutine leaks when it's still alive but will never finish — usually because it's blocked forever on a channel send/receive nobody will ever complete.
- The Go runtime never garbage-collects a blocked goroutine, no matter how "obviously dead" it looks to you — it's still a reachable stack, waiting.
- The fix is almost always the same shape: give every goroutine a way to be told "stop," usually a `context.Context` or a dedicated `done` channel, and make sure something actually cancels it.
- `go build`/`go vet` will not catch this. You find it in production as slowly climbing memory and goroutine counts, or you catch it in tests with `go.uber.org/goleak`.

## The Problem

A goroutine is cheap to start — a few KB of stack, no OS thread required. That cheapness is exactly what makes leaks easy to introduce: nothing forces you to think about how a goroutine *ends*, the way a function call obviously ends when it returns.

A goroutine leak isn't a goroutine that crashed or panicked — it's one that's still technically alive, doing nothing useful, forever. The most common shape:

```go
func process(jobs <-chan int) {
    results := make(chan int)
    go func() {
        for j := range jobs {
            results <- j * 2 // sends here
        }
    }()

    // caller only reads the first result, then moves on
    first := <-results
    fmt.Println(first)
    // the goroutine above is now stuck forever on `results <- j * 2`
    // for the second job, because nobody is reading `results` anymore
}
```

The goroutine inside `process` is blocked on an unbuffered channel send. Nothing will ever read from `results` again once the caller stops after the first value — so that goroutine sits there, holding its stack, its captured variables, and anything referenced through them, for the lifetime of the program. Do this inside a request handler that runs a thousand times a minute, and you have a slow, silent memory leak that looks like nothing in particular until the process falls over.

## Why the Runtime Can't Save You

Go's garbage collector reclaims memory that's *unreachable*. A goroutine blocked on a channel operation is very much reachable — the runtime's scheduler still has a pointer to it, waiting for the channel to become ready. There is no timeout, no "this has been blocked for 10 minutes, let's kill it." The goroutine is doing exactly what you told it to do: wait. Forever, if that's how long it takes.

This is different from, say, a dangling pointer in a manually-memory-managed language. In Go, the leaked resource isn't memory you forgot to free — it's a whole execution context (stack + scheduling metadata) that will never be freed *because it's still logically "in use,"* just uselessly so.

## How to Actually Prevent It

The pattern that fixes almost every goroutine leak is: **never start a goroutine without deciding, at the same time, how it stops.**

### 1. Give it a cancellation signal

```go
func process(ctx context.Context, jobs <-chan int) <-chan int {
    results := make(chan int)
    go func() {
        defer close(results)
        for {
            select {
            case j, ok := <-jobs:
                if !ok {
                    return
                }
                select {
                case results <- j * 2:
                case <-ctx.Done():
                    return
                }
            case <-ctx.Done():
                return
            }
        }
    }()
    return results
}
```

Now the goroutine has two ways out: the `jobs` channel closing (normal completion) or `ctx` being cancelled (the caller giving up early). Either way, it returns instead of blocking forever.

### 2. Use buffered channels where the producer must not block

If a goroutine only ever needs to send *one* result and the caller might not read it, a buffered channel of size 1 sidesteps the leak entirely — the send succeeds immediately into the buffer even if nobody reads it later, and the goroutine exits normally instead of blocking.

```go
results := make(chan int, 1) // buffered — send never blocks here
```

This only works because you know there's exactly one send. It's not a general substitute for cancellation — it just removes the blocking send as the specific cause of the leak in this narrow shape.

### 3. Always pair `context.WithCancel`/`WithTimeout` with `defer cancel()`

```go
ctx, cancel := context.WithTimeout(parent, 5*time.Second)
defer cancel() // releases the context's internal resources even on the happy path
```

Forgetting `cancel()` doesn't leak a goroutine by itself, but it does leak the small internal timer goroutine `context` uses for `WithTimeout`/`WithDeadline` until the timeout actually fires — `defer cancel()` is what makes that deterministic instead of "eventually."

## Catching Leaks Before Production

`go vet` and the compiler have no concept of "this channel operation might never happen" — it's not a static property they can check. Two practical detection strategies:

- **In tests**: [`go.uber.org/goleak`](https://github.com/uber-go/goleak) snapshots running goroutines at the start and end of a test and fails if any unexpected ones are still around. Run it in `TestMain` for a whole package.
- **In production**: expose `runtime.NumGoroutine()` on a metrics endpoint and alert on it trending upward with no corresponding traffic increase. A leak rarely shows up as a spike — it shows up as a slope.

## Common Pitfalls

- **Fire-and-forget goroutines with no channel at all** — if a goroutine's only job is a side effect (writing a log line, say) it can't leak in the blocking sense, but if it holds a lock or a DB connection and panics without recovering, you can still leave things in a bad state. Not every goroutine problem is a classic leak.
- **`for { select { ... } }` without a `ctx.Done()` case** — the most common actual leak in real codebases. If you write an infinite `select` loop, one of its cases must be a way out.
- **Leaking through a `WaitGroup`** — calling `wg.Wait()` when one of the goroutines it's waiting for is itself leaked means `Wait()` blocks forever too. The leak propagates to its caller.
