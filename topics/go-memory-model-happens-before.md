---
title: The Go Memory Model — Happens-Before, Explained
category: Software & Programming
tags: [Go, Concurrency]
duration: 8 min read
relatedCourses: [go-programming, senior-engineer-interview]
relatedProjects: []
relatedTopics: [channels-vs-mutexes, goroutine-leaks]
---

## TL;DR

- The Go memory model answers one question: *given two goroutines, when is it guaranteed that a write in one is visible to a read in the other?* That guarantee is called "happens-before."
- Without a happens-before relationship, the compiler and CPU are free to reorder, cache, or otherwise not synchronize memory between goroutines — even if the code "looks" sequential.
- Channels, mutexes, `sync.WaitGroup`, and `sync.Once` all establish happens-before relationships. A plain shared variable with no synchronization does not, no matter how the code reads.
- This is exactly what the `-race` detector checks for, and exactly why "it worked in my testing" is not evidence of correctness for concurrent code.

## Why This Needs a Formal Model at All

```go
var done bool
var result int

go func() {
    result = 42
    done = true
}()

for !done {
} // busy-wait
fmt.Println(result)
```

This looks obviously correct: the goroutine sets `result` before `done`, so by the time the main goroutine sees `done == true`, surely `result` is already 42. **This is not guaranteed by Go**, and treating it as guaranteed is a real, if subtle, bug — one that might run "correctly" in a test a thousand times and then fail unpredictably in production once the compiler optimizes differently, or on different hardware.

The reason: without an explicit synchronization primitive, the compiler is allowed to reorder those two writes as far as any other goroutine can observe (it just has to look sequential *from within the same goroutine* — that's the only guarantee "program order" gives you). The reading goroutine is also allowed to cache `done` in a register and never see the write at all, spinning forever. Neither of these are exotic edge cases — they're exactly the kind of transformation an optimizing compiler is designed to make when it doesn't know two goroutines are involved.

## What "Happens-Before" Actually Means

The Go memory model defines a partial order over all memory operations across all goroutines, called happens-before. If event A happens-before event B, then A's effects (its writes) are guaranteed visible to B. Critically: **within a single goroutine, happens-before follows normal program order.** The relationship that actually needs establishing is *across* goroutines — and that only happens through specific synchronization operations, not by accident of timing.

## The Primitives That Establish Happens-Before

### Channels

> A send on a channel happens-before the corresponding receive from that channel completes.

```go
var result int
done := make(chan struct{})

go func() {
    result = 42       // (1)
    done <- struct{}{} // (2) send
}()

<-done                // (3) receive — happens-after (2)
fmt.Println(result)   // guaranteed to see 42, because (1) happens-before (2), and (2) happens-before (3)
```

This is the fixed version of the earlier busy-wait example. The channel send/receive pair is the synchronization point that makes the write to `result` visible.

> The closing of a channel happens-before a receive that returns because the channel is closed.

This is why `close(ch)` is a reliable broadcast mechanism — every goroutine that later reads from a closed channel is guaranteed to see everything written before `close` was called.

### Mutexes

> For a `sync.Mutex` or `sync.RWMutex`, the *n*th call to `Unlock` happens-before the *m*th call to `Lock` returns, for any *m* > *n*.

In plain terms: whatever a goroutine wrote while holding the lock is guaranteed visible to the next goroutine that acquires that same lock. This is the entire reason a mutex is sufficient for protecting shared state — it's not just about preventing simultaneous access, it's specifically about establishing this visibility guarantee too.

### sync.WaitGroup

> A call to `Add` happens-before the goroutine it's tracking starts (when used correctly), and the corresponding `Done` happens-before the matching `Wait` returns.

This is why code after `wg.Wait()` can safely read values written by the goroutines it waited for — the `Wait()` call itself is the synchronization point.

### sync.Once

> A single call to `f` inside `once.Do(f)` happens-before any call to `once.Do` returns, for that `once`.

This is what makes `sync.Once` safe for lazy initialization shared across goroutines — every caller of `once.Do` is guaranteed to see the fully-initialized result, not a partially-constructed one.

## What Does *Not* Establish Happens-Before

- **A plain `bool`/`int`/pointer read and write with no synchronization primitive at all** — as in the busy-wait example. This is a *data race*, and Go's memory model explicitly says the behavior in this case is undefined, not just "probably fine."
- **Starting a goroutine with `go f()`** — the goroutine's *start* happens-before anything inside it runs, but nothing in the calling goroutine after the `go` statement is guaranteed to happen-before anything inside the new goroutine. You need an actual synchronization point (a channel, a `WaitGroup`) to coordinate them.
- **Comments, code layout, or "it looks sequential"** — none of these are memory model guarantees. Only the specific documented primitives are.

## The Race Detector Is How You Actually Check This

You cannot reliably eyeball concurrent code and know it's race-free — the busy-wait example above looks correct to most readers. Go ships a runtime instrumentation for exactly this:

```
go test -race ./...
go run -race main.go
```

The race detector doesn't prove correctness (it only catches races that actually occur during that particular execution), but it's dramatically better than nothing, and it's standard practice to run the full test suite with `-race` in CI for any codebase with real concurrency.

## Common Pitfalls

- **Assuming `volatile`-like semantics exist in Go** — they don't. There's no keyword that makes a plain variable "safe" to share across goroutines; you need an actual synchronization primitive.
- **Double-checked locking without `sync.Once`** — hand-rolling a "check a bool, then lock and check again" pattern for lazy init is a classic way to reintroduce exactly the race `sync.Once` was built to prevent. Just use `sync.Once`.
- **Assuming atomics alone give you a full memory barrier for unrelated variables** — `sync/atomic` operations do establish happens-before for the atomic variable itself, but don't assume that extends to protecting other, non-atomic fields in the same struct unless you've read the specific guarantee you're relying on.
