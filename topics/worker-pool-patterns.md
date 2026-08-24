---
title: Worker Pool Patterns in Go
category: Software & Programming
tags: [Go, Concurrency]
duration: 8 min read
relatedCourses: [go-programming]
relatedProjects: [rest-api-server]
relatedTopics: [channels-vs-mutexes, goroutine-leaks, context-cancellation-patterns]
---

## TL;DR

- A worker pool is a fixed number of goroutines pulling work from a shared channel, instead of spawning one goroutine per task — it bounds concurrency so you don't accidentally launch 100,000 goroutines for 100,000 incoming jobs.
- The core shape is always: a jobs channel, N worker goroutines ranging over it, a results channel (if you need output), and a `sync.WaitGroup` to know when everyone's done.
- The number of workers should usually be a deliberate, tunable constant — not "however many jobs arrive."
- Getting shutdown right (closing channels in the right order, respecting cancellation) is the part people actually get wrong.

## Why Not Just `go` a Goroutine Per Task?

```go
for _, job := range jobs {
    go process(job) // one goroutine per job — fine for 10 jobs, dangerous for 10 million
}
```

Goroutines are cheap, but not free — each one still needs stack space, and more importantly, spawning unboundedly many of them means you have no control over how much concurrent work is actually happening. If `process` hits a shared resource (a database connection pool, a rate-limited external API), an unbounded number of simultaneous goroutines will just overwhelm it. A worker pool exists specifically to put a ceiling on "how many of these run at once," independent of how many jobs show up.

## The Standard Shape

```go
func runPool(numWorkers int, jobs []int) []int {
    jobCh := make(chan int)
    resultCh := make(chan int, len(jobs)) // buffered so workers never block on send

    var wg sync.WaitGroup
    for w := 0; w < numWorkers; w++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobCh { // exits when jobCh is closed and drained
                resultCh <- process(job)
            }
        }()
    }

    go func() {
        for _, j := range jobs {
            jobCh <- j
        }
        close(jobCh) // signals workers: no more jobs coming
    }()

    go func() {
        wg.Wait()      // wait for all workers to finish draining jobCh
        close(resultCh) // now safe to close — nobody will send to it again
    }()

    results := make([]int, 0, len(jobs))
    for r := range resultCh { // reads until resultCh is closed
        results = append(results, r)
    }
    return results
}
```

Walk through the shutdown sequence carefully, because it's the part that's easy to get subtly wrong:

1. The job producer closes `jobCh` once every job has been sent — `for job := range jobCh` in each worker exits cleanly when a channel is both closed *and* drained.
2. A **separate** goroutine calls `wg.Wait()` and only then closes `resultCh`. This has to happen in its own goroutine, not inline after starting the workers, because `wg.Wait()` blocks — if you called it directly in `runPool` before reading from `resultCh`, you'd deadlock (workers can't finish because they're blocked sending to an unbuffered/full `resultCh` that nobody's draining, and the main goroutine can't drain it because it's stuck on `wg.Wait()` first).
3. The final `for r := range resultCh` in the caller drains results as they arrive and exits when `resultCh` is closed.

## Choosing the Number of Workers

There's no universal right number — it depends entirely on what the work is bottlenecked by:

- **CPU-bound work** (hashing, image processing, compression): `runtime.NumCPU()` is a reasonable starting point — more workers than cores just adds scheduling overhead without more real parallelism.
- **I/O-bound work** (HTTP calls, database queries): the right number is usually much higher than `NumCPU()`, since workers spend most of their time blocked waiting on a network round-trip, not consuming CPU. The actual ceiling here is usually an external constraint — how many concurrent connections your database or the API you're calling can handle.
- **Rate-limited external APIs**: sometimes the "worker pool" size should just directly match whatever concurrency limit the external API imposes, full stop.

```go
numWorkers := runtime.NumCPU() // CPU-bound default; override deliberately for I/O-bound work
```

## Adding Context Cancellation

A production worker pool almost always needs to stop early — a client disconnects, a shutdown signal arrives. Thread a `context.Context` through:

```go
func worker(ctx context.Context, jobCh <-chan int, resultCh chan<- int) {
    for {
        select {
        case job, ok := <-jobCh:
            if !ok {
                return
            }
            select {
            case resultCh <- process(job):
            case <-ctx.Done():
                return
            }
        case <-ctx.Done():
            return
        }
    }
}
```

Every blocking operation — receiving a job, sending a result — is raced against `ctx.Done()` so cancellation actually takes effect promptly instead of waiting for the current job (or the whole remaining queue) to finish first.

## Bounding the Jobs Channel Itself

The example above uses an unbuffered `jobCh`, which means the producer blocks until a worker is ready to take the next job — a natural form of backpressure. If instead you want the producer to be able to queue up a batch of work without waiting on workers, use a buffered channel with a deliberate size:

```go
jobCh := make(chan int, 100) // producer can queue up to 100 jobs ahead of the workers
```

Making this buffer unbounded (or very large "just in case") reintroduces the same problem worker pools exist to prevent — an unbounded queue can grow without limit if producers outpace workers, just moved from "too many goroutines" to "too much memory in a channel buffer."

## Common Pitfalls

- **Calling `wg.Wait()` before starting to read results, in the same goroutine** — the deadlock described above. `wg.Wait()` and draining the results channel must happen concurrently, not sequentially, unless the results channel is buffered large enough to hold every result without blocking (which only works if you know the total count in advance).
- **Forgetting `defer wg.Done()`** — if a worker panics or hits an early `return` and skips `wg.Done()`, `wg.Wait()` blocks forever. Always pair `Add`/`Done` with `defer` right at the top of the goroutine.
- **Sharing mutable state across workers without synchronization** — a worker pool doesn't automatically make concurrent access to shared state safe. If workers write to a shared map or slice, you still need a mutex around that specific state (see Channels vs Mutexes in Related).
- **Not bounding the number of workers at all** — spawning `len(jobs)` workers instead of a fixed pool size defeats the entire purpose; you're back to one-goroutine-per-task with extra steps.
