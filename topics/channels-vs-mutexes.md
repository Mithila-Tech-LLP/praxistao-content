---
title: Channels vs Mutexes — Choosing the Right Tool
category: Software & Programming
tags: [Go, Concurrency]
duration: 7 min read
relatedCourses: [go-programming, senior-engineer-interview]
relatedProjects: []
relatedTopics: [goroutine-leaks, worker-pool-patterns]
---

## TL;DR

- A mutex protects **shared state** — use it when multiple goroutines need to read/modify the same in-memory data safely.
- A channel coordinates **shared work or events** — use it when goroutines need to hand something to each other, signal completion, or communicate across a pipeline.
- Go's own proverb — "share memory by communicating, don't communicate by sharing memory" — is a real design preference, not just a slogan, but it's not an absolute rule. Plenty of correct, idiomatic Go uses `sync.Mutex` directly.
- The decision usually comes down to: is this a *value* that needs protecting, or an *event/message* that needs delivering?

## The Core Difference

A `sync.Mutex` is a lock. It doesn't move data anywhere — it just makes sure only one goroutine at a time can touch a piece of shared state, so concurrent reads/writes don't race.

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func (c *Counter) Inc() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}
```

A channel, by contrast, *transfers ownership* of a value from one goroutine to another. Once you send a value on a channel, the idiom is that the sender stops touching it — the receiver now owns it. There's no shared state to protect because, conceptually, only one goroutine has the data at a time.

```go
func generate(out chan<- int) {
    for i := 0; i < 10; i++ {
        out <- i // ownership of `i` passes to whoever receives it
    }
    close(out)
}
```

## When a Mutex Is the Right Call

If you have one piece of state — a counter, a cache, an in-memory map — that many goroutines need to read and write concurrently, a mutex is usually simpler and faster than building a channel-based "actor" to own that state instead.

```go
type Cache struct {
    mu   sync.RWMutex
    data map[string]string
}

func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()
    defer c.mu.RUnlock()
    v, ok := c.data[key]
    return v, ok
}

func (c *Cache) Set(key, value string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.data[key] = value
}
```

`sync.RWMutex` specifically is worth knowing here: it allows many concurrent readers, or one exclusive writer, which matters a lot for a cache that's read far more often than it's written. A plain `sync.Mutex` would serialize reads unnecessarily.

You could implement this cache with a goroutine that owns the map and a channel of "get"/"set" requests instead — but that's strictly more code, an extra goroutine to manage, and no real benefit for a case this simple. Reach for the mutex.

## When a Channel Is the Right Call

Channels shine when the problem is really about *coordination* between independent units of work, not shared state:

- **Pipelines** — stage A produces values, stage B consumes and transforms them, stage C consumes those. Each stage is a goroutine; channels are the conveyor belts between them.
- **Fan-out/fan-in** — one goroutine distributes work across N worker goroutines via a channel, and their results are gathered back through another channel.
- **Signaling** — a `done chan struct{}` that's closed to broadcast "stop" to any number of listening goroutines is a channel doing a job a mutex simply can't do (a mutex protects state; it doesn't broadcast events).

```go
func fanIn(a, b <-chan int) <-chan int {
    out := make(chan int)
    var wg sync.WaitGroup
    wg.Add(2)
    merge := func(c <-chan int) {
        defer wg.Done()
        for v := range c {
            out <- v
        }
    }
    go merge(a)
    go merge(b)
    go func() {
        wg.Wait()
        close(out)
    }()
    return out
}
```

Notice the mutex/channel pairing here isn't either/or in a strict sense — the `sync.WaitGroup` above is doing a *coordination* job (wait for both merge goroutines to finish) that's neither a mutex nor a channel, but shows up constantly alongside both.

## A Concrete Test: "What Am I Actually Protecting?"

Ask: is there a **value that must never be read half-written**, or is there a **thing that needs to happen after another thing finishes**?

- "Two goroutines might increment this same integer at once" → shared state → mutex.
- "This goroutine needs to know when that other goroutine is done" → event → channel (or `sync.WaitGroup`, `context.Context`, depending on shape).
- "I need exactly one goroutine to ever run this initialization code" → `sync.Once`, a third primitive worth knowing for this specific job.

## Common Pitfalls

- **Using a channel purely to protect a variable, with only one goroutine reading and one writing "turns" through it** — this works but is usually more code and harder to read than just using a mutex. Channels aren't inherently more idiomatic; they're the right tool for a specific shape of problem.
- **Using a mutex to build a queue** — passing work items through a mutex-protected slice with manual signaling reimplements, badly, what a buffered channel already does correctly and more readably.
- **Holding a mutex while doing something slow (a network call, a channel send)** — this turns a fast, short-lived lock into a long-held one, and every other goroutine wanting that lock now queues up behind the slow operation. Keep the critical section as small as possible.
- **Copying a `sync.Mutex` by value** — a struct containing a mutex must always be passed by pointer. Copying it copies the lock's state, silently breaking the mutual exclusion it was supposed to provide. `go vet` catches this one, at least.
