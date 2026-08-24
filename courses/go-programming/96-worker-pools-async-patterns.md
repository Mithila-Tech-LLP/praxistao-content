# Chapter 92: Worker Pools and Async Patterns

Concurrency in Go is cheap, but uncontrolled concurrency is dangerous. Spawning one goroutine per incoming request or per database row can exhaust memory, connections, and file descriptors. Worker pools give you bounded concurrency: a fixed number of workers consuming from a shared job channel.

## Table of Contents

1. [Unbounded vs Bounded Concurrency](#1-unbounded-vs-bounded-concurrency)
2. [Fixed Worker Pool](#2-fixed-worker-pool)
3. [Dynamic Worker Pool](#3-dynamic-worker-pool)
4. [Pipeline Pattern](#4-pipeline-pattern)
5. [Fan-Out, Fan-In](#5-fan-out-fan-in)
6. [Backpressure](#6-backpressure)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Unbounded vs Bounded Concurrency

```go
// BAD: one goroutine per item — may spawn millions
func processAll(urls []string) {
    var wg sync.WaitGroup
    for _, url := range urls {
        wg.Add(1)
        go func(u string) {
            defer wg.Done()
            fetch(u) // may timeout, but first you run out of memory
        }(url)
    }
    wg.Wait()
}

// GOOD: fixed pool of workers
func processAll(urls []string) {
    jobs := make(chan string, len(urls))
    for _, url := range urls { jobs <- url }
    close(jobs)
    
    var wg sync.WaitGroup
    for range 10 { // exactly 10 concurrent workers
        wg.Add(1)
        go func() {
            defer wg.Done()
            for url := range jobs { fetch(url) }
        }()
    }
    wg.Wait()
}
```

---

## 2. Fixed Worker Pool

A reusable worker pool with jobs, results, and graceful shutdown:

```go
type Job[I, O any] struct {
    Input I
    Index int
}

type Result[I, O any] struct {
    Input  I
    Output O
    Err    error
    Index  int
}

type Pool[I, O any] struct {
    workers int
    process func(context.Context, I) (O, error)
}

func NewPool[I, O any](workers int, fn func(context.Context, I) (O, error)) *Pool[I, O] {
    return &Pool[I, O]{workers: workers, process: fn}
}

// Run processes all inputs concurrently and returns results in any order.
// Results are returned in the order they complete, not the order of inputs.
func (p *Pool[I, O]) Run(ctx context.Context, inputs []I) []Result[I, O] {
    jobs := make(chan Job[I, O], len(inputs))
    results := make(chan Result[I, O], len(inputs))
    
    // Launch workers
    var wg sync.WaitGroup
    for range p.workers {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for job := range jobs {
                out, err := p.process(ctx, job.Input)
                results <- Result[I, O]{
                    Input:  job.Input,
                    Output: out,
                    Err:    err,
                    Index:  job.Index,
                }
            }
        }()
    }
    
    // Send jobs
    for i, input := range inputs {
        jobs <- Job[I, O]{Input: input, Index: i}
    }
    close(jobs)
    
    // Close results when all workers finish
    go func() { wg.Wait(); close(results) }()
    
    // Collect results
    out := make([]Result[I, O], 0, len(inputs))
    for r := range results { out = append(out, r) }
    return out
}

// Usage
pool := NewPool(10, func(ctx context.Context, url string) ([]byte, error) {
    return fetch(ctx, url)
})
results := pool.Run(ctx, urls)
for _, r := range results {
    if r.Err != nil { log.Printf("failed %s: %v", r.Input, r.Err); continue }
    process(r.Output)
}
```

---

## 3. Dynamic Worker Pool

Scale workers up and down based on queue depth:

```go
type DynamicPool struct {
    minWorkers  int
    maxWorkers  int
    jobs        chan func(context.Context)
    activeCount atomic.Int32
    mu          sync.Mutex
    cancel      context.CancelFunc
    wg          sync.WaitGroup
}

func NewDynamicPool(min, max int) *DynamicPool {
    ctx, cancel := context.WithCancel(context.Background())
    p := &DynamicPool{
        minWorkers: min,
        maxWorkers: max,
        jobs:       make(chan func(context.Context), max*10),
        cancel:     cancel,
    }
    
    // Start minimum workers
    for range min {
        p.startWorker(ctx)
    }
    
    // Scale-up monitor
    go p.monitor(ctx)
    return p
}

func (p *DynamicPool) Submit(fn func(context.Context)) bool {
    select {
    case p.jobs <- fn:
        return true
    default:
        return false // queue full
    }
}

func (p *DynamicPool) startWorker(ctx context.Context) {
    p.wg.Add(1)
    p.activeCount.Add(1)
    go func() {
        defer func() { p.wg.Done(); p.activeCount.Add(-1) }()
        for {
            select {
            case fn, ok := <-p.jobs:
                if !ok { return }
                fn(ctx)
            case <-ctx.Done():
                return
            }
        }
    }()
}

func (p *DynamicPool) monitor(ctx context.Context) {
    ticker := time.NewTicker(100 * time.Millisecond)
    defer ticker.Stop()
    for {
        select {
        case <-ctx.Done(): return
        case <-ticker.C:
            queueDepth := len(p.jobs)
            active := int(p.activeCount.Load())
            // Scale up if queue growing and we have capacity
            if queueDepth > active*2 && active < p.maxWorkers {
                p.startWorker(ctx)
            }
        }
    }
}

func (p *DynamicPool) Shutdown() {
    p.cancel()
    p.wg.Wait()
}
```

---

## 4. Pipeline Pattern

Process data through stages, each running concurrently:

```go
// Each stage takes a channel in and returns a channel out
func generateURLs(urls []string) <-chan string {
    ch := make(chan string)
    go func() {
        defer close(ch)
        for _, url := range urls { ch <- url }
    }()
    return ch
}

func fetchPages(ctx context.Context, in <-chan string) <-chan []byte {
    out := make(chan []byte, 10)
    go func() {
        defer close(out)
        for url := range in {
            data, err := fetch(ctx, url)
            if err != nil { continue }
            out <- data
        }
    }()
    return out
}

func parseProducts(in <-chan []byte) <-chan Product {
    out := make(chan Product, 10)
    go func() {
        defer close(out)
        for data := range in {
            for _, p := range parse(data) { out <- p }
        }
    }()
    return out
}

// Wire the pipeline
func scrape(ctx context.Context, urls []string) []Product {
    urls_ch  := generateURLs(urls)
    pages_ch := fetchPages(ctx, urls_ch)
    prods_ch := parseProducts(pages_ch)
    
    var products []Product
    for p := range prods_ch { products = append(products, p) }
    return products
}

// Parallel stage: run N workers for the slow fetch stage
func fetchPagesConcurrent(ctx context.Context, in <-chan string, workers int) <-chan []byte {
    out := make(chan []byte, workers*2)
    
    var wg sync.WaitGroup
    for range workers {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for url := range in {
                data, err := fetch(ctx, url)
                if err != nil { continue }
                out <- data
            }
        }()
    }
    go func() { wg.Wait(); close(out) }()
    return out
}
```

---

## 5. Fan-Out, Fan-In

```go
// Fan-out: one input channel → multiple output channels
func fanOut[T any](in <-chan T, n int) []<-chan T {
    outs := make([]chan T, n)
    for i := range outs { outs[i] = make(chan T, 1) }
    
    go func() {
        defer func() {
            for _, ch := range outs { close(ch) }
        }()
        i := 0
        for v := range in {
            outs[i%n] <- v
            i++
        }
    }()
    
    result := make([]<-chan T, n)
    for i, ch := range outs { result[i] = ch }
    return result
}

// Fan-in: multiple input channels → one output channel
func fanIn[T any](inputs ...<-chan T) <-chan T {
    out := make(chan T, len(inputs))
    
    var wg sync.WaitGroup
    for _, in := range inputs {
        wg.Add(1)
        go func(ch <-chan T) {
            defer wg.Done()
            for v := range ch { out <- v }
        }(in)
    }
    go func() { wg.Wait(); close(out) }()
    return out
}
```

---

## 6. Backpressure

Backpressure prevents fast producers from overwhelming slow consumers:

```go
// Semaphore: limit concurrent operations
type Semaphore struct{ ch chan struct{} }

func NewSemaphore(n int) *Semaphore {
    return &Semaphore{ch: make(chan struct{}, n)}
}

func (s *Semaphore) Acquire(ctx context.Context) error {
    select {
    case s.ch <- struct{}{}: return nil
    case <-ctx.Done():       return ctx.Err()
    }
}

func (s *Semaphore) Release() { <-s.ch }

// Usage: at most 10 concurrent DB queries
sem := NewSemaphore(10)
for _, id := range ids {
    if err := sem.Acquire(ctx); err != nil { break }
    go func(id int64) {
        defer sem.Release()
        doDBQuery(ctx, id)
    }(id)
}

// Rate limiter as backpressure
type RateLimiter struct{ ticker *time.Ticker }

func NewRateLimiter(rps int) *RateLimiter {
    return &RateLimiter{ticker: time.NewTicker(time.Second / time.Duration(rps))}
}

func (r *RateLimiter) Wait(ctx context.Context) error {
    select {
    case <-r.ticker.C: return nil
    case <-ctx.Done(): return ctx.Err()
    }
}
```

---

## Summary

- **Fixed pool**: N workers reading from one buffered channel; `close(jobs)` signals done, `wg.Wait()` drains
- **Results channel**: collect results via a separate channel; close it after all workers finish
- **Pipeline**: chain channels between stages; each stage is a goroutine transforming its input channel
- **Fan-out**: distribute work across N parallel workers for the same stage
- **Fan-in**: merge N channels into one for the next stage
- **Backpressure**: semaphore or rate limiter prevents the producer from flooding the consumer
- Always handle `ctx.Done()` in workers — pools must be context-aware for graceful shutdown

## Exercises

### Easy
1. Write a `Batch[T, R any](items []T, batchSize, workers int, fn func([]T) []R) []R` function that splits `items` into batches of `batchSize`, processes each batch with `fn` using `workers` goroutines, and returns all results.
2. Implement a `ThrottledFetcher` that limits concurrent HTTP requests to at most 5 at a time using a semaphore. Fetch 100 URLs and print the status code of each.
3. Build a pipeline: generate numbers 1-100 → square each → filter odd squares → sum. Each stage is a goroutine communicating via channels.

### Medium
4. Implement a **priority queue worker pool**: jobs have a priority (high/medium/low). High-priority jobs should always be processed before lower-priority ones when workers are available. Use Go's `container/heap` for the internal queue.
5. Build a **rate-limited pipeline stage**: the fetch stage should make at most 10 requests per second regardless of how many workers are running. Use a `time.Ticker` shared across all fetch workers.
6. Implement **circuit-breaker protection** on a worker pool: if more than 50% of the last 10 jobs failed, the pool opens its circuit and rejects new jobs for 30 seconds before retrying.

### Hard
7. Build a **back-pressure-aware streaming processor**: read records from a Kafka consumer (or a channel simulating it), process each through a 5-stage pipeline, and write results to a database. If the DB write stage falls behind, slow down the Kafka consumer using pause/resume so memory doesn't grow unboundedly.
8. Implement a **work-stealing scheduler**: N workers each have a local queue. A worker that finishes its queue steals from the largest non-empty queue. Compare throughput of work-stealing vs a single shared queue under skewed workloads (some jobs take 100× longer than others).
