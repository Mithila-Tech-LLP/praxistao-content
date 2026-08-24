# Chapter 16: Concurrency Patterns — Worker Pools, Pipelines & Fan-out/Fan-in

Production Go services use a set of recurring concurrency patterns. This chapter teaches you to recognize and implement each one correctly, with edge cases handled — goroutine leaks prevented, cancellation propagated, backpressure applied.

## Table of Contents

1. [Worker Pool](#1-worker-pool)
2. [Pipeline Pattern](#2-pipeline-pattern)
3. [Fan-out / Fan-in](#3-fan-out--fan-in)
4. [Bounded Parallelism](#4-bounded-parallelism)
5. [Error Group Pattern](#5-error-group-pattern)
6. [Pub/Sub with Channels](#6-pubsub-with-channels)
7. [Rate Limiting with a Token Bucket](#7-rate-limiting)
8. [Summary](#summary)

---

## 1. Worker Pool

The worker pool pattern runs a fixed number of goroutines that process work from a shared queue. It provides bounded concurrency — you don't spawn an unlimited number of goroutines.

```go
func workerPool(jobs <-chan Job, results chan<- Result, numWorkers int) {
    var wg sync.WaitGroup
    wg.Add(numWorkers)

    for i := 0; i < numWorkers; i++ {
        go func(id int) {
            defer wg.Done()
            for job := range jobs { // worker exits when jobs channel is closed
                result := processJob(job)
                results <- result
            }
        }(i)
    }

    // Close results when all workers are done
    go func() {
        wg.Wait()
        close(results)
    }()
}

// Usage:
func main() {
    const numJobs = 100
    const numWorkers = 5

    jobs := make(chan Job, numJobs)
    results := make(chan Result, numJobs)

    workerPool(jobs, results, numWorkers)

    // Send jobs
    for i := 0; i < numJobs; i++ {
        jobs <- Job{ID: i}
    }
    close(jobs) // signal no more jobs

    // Collect results
    for result := range results {
        fmt.Println(result)
    }
}
```

### Worker Pool with Cancellation

```go
func workerPoolWithCancel(ctx context.Context, jobs <-chan Job, numWorkers int) <-chan Result {
    results := make(chan Result, numWorkers)
    var wg sync.WaitGroup
    wg.Add(numWorkers)

    for i := 0; i < numWorkers; i++ {
        go func() {
            defer wg.Done()
            for {
                select {
                case job, ok := <-jobs:
                    if !ok { return } // jobs channel closed
                    
                    select {
                    case results <- processJob(job):
                    case <-ctx.Done(): // stop if context cancelled
                        return
                    }
                    
                case <-ctx.Done():
                    return
                }
            }
        }()
    }

    go func() {
        wg.Wait()
        close(results)
    }()

    return results
}
```

---

## 2. Pipeline Pattern

A pipeline is a series of processing stages connected by channels. Each stage takes input from a channel, transforms it, and sends output to another channel.

```go
// Stage 1: Generate numbers
func generate(ctx context.Context, nums ...int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for _, n := range nums {
            select {
            case out <- n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Stage 2: Square each number
func square(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            select {
            case out <- n * n:
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Stage 3: Filter evens
func filterEven(ctx context.Context, in <-chan int) <-chan int {
    out := make(chan int)
    go func() {
        defer close(out)
        for n := range in {
            if n%2 == 0 {
                select {
                case out <- n:
                case <-ctx.Done():
                    return
                }
            }
        }
    }()
    return out
}

// Chain the stages:
func main() {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    gen := generate(ctx, 1, 2, 3, 4, 5, 6, 7, 8, 9, 10)
    sq := square(ctx, gen)
    filtered := filterEven(ctx, sq)

    for v := range filtered {
        fmt.Println(v) // 4, 16, 36, 64, 100
    }
}
```

---

## 3. Fan-out / Fan-in

**Fan-out:** Distribute work across multiple goroutines reading from the same channel.
**Fan-in:** Merge multiple channels into one.

```go
// Fan-out: multiple workers reading from the same input
func fanOut(ctx context.Context, in <-chan Job, n int) []<-chan Result {
    channels := make([]<-chan Result, n)
    for i := 0; i < n; i++ {
        channels[i] = worker(ctx, in) // all workers share the same input
    }
    return channels
}

func worker(ctx context.Context, jobs <-chan Job) <-chan Result {
    out := make(chan Result)
    go func() {
        defer close(out)
        for job := range jobs {
            select {
            case out <- processJob(job):
            case <-ctx.Done():
                return
            }
        }
    }()
    return out
}

// Fan-in: merge multiple output channels into one
func fanIn(ctx context.Context, channels ...<-chan Result) <-chan Result {
    merged := make(chan Result)
    var wg sync.WaitGroup

    // Forward from each channel to merged
    forward := func(ch <-chan Result) {
        defer wg.Done()
        for result := range ch {
            select {
            case merged <- result:
            case <-ctx.Done():
                return
            }
        }
    }

    wg.Add(len(channels))
    for _, ch := range channels {
        go forward(ch)
    }

    // Close merged when all forwards are done
    go func() {
        wg.Wait()
        close(merged)
    }()

    return merged
}

// Full usage:
func processWithFanOutFanIn(jobs []Job, numWorkers int) []Result {
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    jobCh := make(chan Job, len(jobs))
    for _, j := range jobs { jobCh <- j }
    close(jobCh)

    workerOutputs := fanOut(ctx, jobCh, numWorkers)
    merged := fanIn(ctx, workerOutputs...)

    results := []Result{}
    for r := range merged {
        results = append(results, r)
    }
    return results
}
```

---

## 4. Bounded Parallelism

Process N items with at most maxConcurrent goroutines at a time — without a pre-built worker pool.

```go
// Semaphore pattern using buffered channel
func processBounded(items []Item, maxConcurrent int) []Result {
    sem := make(chan struct{}, maxConcurrent) // semaphore
    var wg sync.WaitGroup
    var mu sync.Mutex
    results := make([]Result, 0, len(items))

    for _, item := range items {
        wg.Add(1)
        sem <- struct{}{} // acquire semaphore (blocks at maxConcurrent)
        
        go func(i Item) {
            defer func() {
                <-sem // release semaphore
                wg.Done()
            }()
            
            result := process(i)
            mu.Lock()
            results = append(results, result)
            mu.Unlock()
        }(item)
    }

    wg.Wait()
    return results
}
```

---

## 5. Error Group Pattern

`golang.org/x/sync/errgroup` provides a WaitGroup that propagates the first error and optionally cancels on error.

```go
import "golang.org/x/sync/errgroup"

func fetchAll(ctx context.Context, urls []string) ([][]byte, error) {
    g, ctx := errgroup.WithContext(ctx)
    results := make([][]byte, len(urls))

    for i, url := range urls {
        i, url := i, url // capture loop variables
        g.Go(func() error {
            data, err := fetch(ctx, url)
            if err != nil { return err }
            results[i] = data // safe: each goroutine writes to its own index
            return nil
        })
    }

    // Wait for all goroutines. If any returned an error, ctx is cancelled.
    if err := g.Wait(); err != nil {
        return nil, err
    }
    return results, nil
}
```

### Manual Error Group (for interviews)

```go
// Implementing the same behavior without the package:
func fetchAllManual(ctx context.Context, urls []string) ([][]byte, error) {
    ctx, cancel := context.WithCancel(ctx)
    defer cancel()

    results := make([][]byte, len(urls))
    errs := make(chan error, len(urls)) // buffered to not block
    var wg sync.WaitGroup

    for i, url := range urls {
        i, url := i, url
        wg.Add(1)
        go func() {
            defer wg.Done()
            data, err := fetch(ctx, url)
            if err != nil {
                cancel()     // cancel other goroutines
                errs <- err  // report error
                return
            }
            results[i] = data
        }()
    }

    wg.Wait()
    close(errs)

    for err := range errs {
        if err != nil { return nil, err } // return first error
    }
    return results, nil
}
```

---

## 6. Pub/Sub with Channels

A publish/subscribe system lets multiple subscribers receive the same events.

```go
type PubSub struct {
    mu          sync.RWMutex
    subscribers map[string][]chan Event
}

func NewPubSub() *PubSub {
    return &PubSub{subscribers: make(map[string][]chan Event)}
}

func (ps *PubSub) Subscribe(topic string) <-chan Event {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    ch := make(chan Event, 10) // buffered to avoid blocking the publisher
    ps.subscribers[topic] = append(ps.subscribers[topic], ch)
    return ch
}

func (ps *PubSub) Publish(topic string, event Event) {
    ps.mu.RLock()
    defer ps.mu.RUnlock()
    for _, ch := range ps.subscribers[topic] {
        select {
        case ch <- event:
        default: // drop if subscriber is slow (non-blocking)
        }
    }
}

func (ps *PubSub) Unsubscribe(topic string, ch <-chan Event) {
    ps.mu.Lock()
    defer ps.mu.Unlock()
    subs := ps.subscribers[topic]
    for i, sub := range subs {
        if sub == ch {
            ps.subscribers[topic] = append(subs[:i], subs[i+1:]...)
            close(sub)
            return
        }
    }
}
```

---

## 7. Rate Limiting

A token bucket rate limiter using `time.Ticker`:

```go
// Simple rate limiter: N requests per second
type RateLimiter struct {
    ticker  *time.Ticker
    tokens  chan struct{}
    quit    chan struct{}
}

func NewRateLimiter(rps int) *RateLimiter {
    rl := &RateLimiter{
        ticker: time.NewTicker(time.Second / time.Duration(rps)),
        tokens: make(chan struct{}, rps),
        quit:   make(chan struct{}),
    }
    go rl.refill()
    return rl
}

func (rl *RateLimiter) refill() {
    for {
        select {
        case <-rl.ticker.C:
            select {
            case rl.tokens <- struct{}{}: // add token
            default: // bucket full, drop
            }
        case <-rl.quit:
            return
        }
    }
}

func (rl *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-rl.tokens: // consume a token
        return nil
    case <-ctx.Done():
        return ctx.Err()
    }
}

// In production: use golang.org/x/time/rate which has a proper token bucket:
import "golang.org/x/time/rate"
limiter := rate.NewLimiter(rate.Limit(100), 10) // 100 r/s, burst 10
if err := limiter.Wait(ctx); err != nil { ... }
```

---

## Summary

- **Worker pool:** fixed number of goroutines, shared input channel. Always close the jobs channel to signal completion.
- **Pipeline:** chain of stages connected by channels. Each stage closes its output when input is exhausted.
- **Fan-out:** multiple goroutines reading from the same input channel — Go channels are safe for concurrent receivers.
- **Fan-in:** merge multiple output channels into one using a WaitGroup + a closer goroutine.
- **Bounded parallelism:** semaphore via buffered channel. Acquire before launching goroutine, release when done.
- **errgroup:** WaitGroup that propagates errors and can cancel context on first error. Essential for parallel API calls.
- All patterns should handle context cancellation to avoid goroutine leaks.
