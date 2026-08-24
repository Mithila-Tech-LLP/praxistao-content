# Chapter 21: The sync Package — Mutexes, Once, Pool, and More

Channels handle communication between goroutines; the `sync` package handles **coordination and shared state protection**. When you need multiple goroutines to safely read and write the same data without passing it through a channel, the sync package is your toolkit. This chapter covers every synchronization primitive you'll use in production Go.

## Table of Contents

1. [Mutex — Mutual Exclusion](#1-mutex--mutual-exclusion)
2. [RWMutex — Reader-Writer Lock](#2-rwmutex--reader-writer-lock)
3. [Once — Run Exactly Once](#3-once--run-exactly-once)
4. [WaitGroup (Deeper Look)](#4-waitgroup-deeper-look)
5. [Cond — Condition Variable](#5-cond--condition-variable)
6. [Map — Concurrent Map](#6-map--concurrent-map)
7. [Pool — Object Reuse](#7-pool--object-reuse)
8. [Atomic Operations](#8-atomic-operations)
9. [Choosing the Right Primitive](#9-choosing-the-right-primitive)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Mutex — Mutual Exclusion

A mutex ensures only one goroutine can execute a critical section at a time:

```go
import "sync"

type Counter struct {
    mu    sync.Mutex  // Protects value
    value int
}

func (c *Counter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.value++
}

func (c *Counter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.value
}
```

**Critical section rules:**
- Lock before accessing protected data, unlock immediately after
- Always `defer c.mu.Unlock()` right after `c.mu.Lock()` — prevents forgetting to unlock
- Keep critical sections short — holding a mutex blocks all other goroutines

**What a mutex protects:**
```go
type Cache struct {
    mu    sync.Mutex
    items map[string]Item   // Shared data — must be protected
}

func (c *Cache) Get(key string) (Item, bool) {
    c.mu.Lock()
    defer c.mu.Unlock()
    item, ok := c.items[key]
    return item, ok
}

func (c *Cache) Set(key string, item Item) {
    c.mu.Lock()
    defer c.mu.Unlock()
    c.items[key] = item
}

func (c *Cache) Delete(key string) {
    c.mu.Lock()
    defer c.mu.Unlock()
    delete(c.items, key)
}
```

**Lock correctly — protect the DATA, not the operation:**
```go
// WRONG: lock/unlock around individual operations leaves gaps
c.mu.Lock()
val := c.items[key]
c.mu.Unlock()

if val == nil {
    // Another goroutine might have set key between unlock and this check!
    c.mu.Lock()
    c.items[key] = newVal
    c.mu.Unlock()
}

// CORRECT: hold lock for the entire check-then-act sequence
c.mu.Lock()
defer c.mu.Unlock()
if c.items[key] == nil {
    c.items[key] = newVal  // This is atomic — no gap
}
```

**Mutex must not be copied:**
```go
// BAD: copying a mutex breaks it (copy contains lock state)
func badFunction(m sync.Mutex) {  // Copies the mutex!
    m.Lock()
}

// GOOD: pass pointer
func goodFunction(m *sync.Mutex) {
    m.Lock()
    defer m.Unlock()
}

// When embedding mutex in struct, always use pointer to struct:
cache := &Cache{}  // Use pointer, not value, so mutex isn't copied
```

### Quick Check
> 1. What does a mutex prevent?
> 2. Why use `defer mu.Unlock()` right after `mu.Lock()`?
> 3. Can you copy a `sync.Mutex`?

---

## 2. RWMutex — Reader-Writer Lock

`sync.RWMutex` allows multiple simultaneous readers OR one writer — not both at the same time. Use when reads are much more frequent than writes:

```go
type Config struct {
    mu     sync.RWMutex
    values map[string]string
}

// Multiple goroutines can Read() simultaneously:
func (c *Config) Get(key string) string {
    c.mu.RLock()          // Read lock — allows other readers
    defer c.mu.RUnlock()
    return c.values[key]
}

// Write() is exclusive — blocks all readers and writers:
func (c *Config) Set(key, value string) {
    c.mu.Lock()           // Write lock — exclusive
    defer c.mu.Unlock()
    c.values[key] = value
}
```

**When to use RWMutex vs Mutex:**
```go
// Use Mutex when reads and writes are roughly equal frequency
// Use RWMutex when reads >> writes (e.g., 100:1 ratio)

// A goroutine count: many reads, rare writes → RWMutex
type Registry struct {
    mu       sync.RWMutex
    handlers map[string]Handler
}

func (r *Registry) Get(name string) (Handler, bool) {
    r.mu.RLock()
    defer r.mu.RUnlock()
    h, ok := r.handlers[name]
    return h, ok
}

func (r *Registry) Register(name string, h Handler) {
    r.mu.Lock()
    defer r.mu.Unlock()
    r.handlers[name] = h
}
```

**RLock/RUnlock vs Lock/Unlock:**
```go
mu.RLock()    // Acquire read lock
mu.RUnlock()  // Release read lock

mu.Lock()     // Acquire write lock (exclusive)
mu.Unlock()   // Release write lock
```

### Quick Check
> 1. How many goroutines can hold a read lock simultaneously?
> 2. When a goroutine holds a write lock, can other goroutines read?
> 3. Should you use `RWMutex` if reads and writes are 50/50?

---

## 3. Once — Run Exactly Once

`sync.Once` guarantees a function runs exactly once, even if called from multiple goroutines simultaneously:

```go
type Singleton struct {
    // ...
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{
            // Expensive initialization...
        }
    })
    return instance
}

// Thread-safe: only ONE goroutine runs the init, others wait for it to finish
g1 := go GetInstance()
g2 := go GetInstance()
g3 := go GetInstance()
// All three get the same instance — init runs only once
```

**Lazy initialization of a connection pool:**
```go
type DB struct {
    once sync.Once
    pool *sql.DB
}

func (db *DB) getPool() *sql.DB {
    db.once.Do(func() {
        var err error
        db.pool, err = sql.Open("postgres", connStr)
        if err != nil {
            panic("failed to open DB: " + err.Error())
        }
    })
    return db.pool
}
```

**Important**: `once.Do` always calls the function at most once, even if the first call panics:
```go
var once sync.Once
count := 0

// If f panics, once.Do will never call f again:
safe := func() {
    once.Do(func() {
        count++
        panic("oops")
    })
}

// First call: panics (and count is incremented)
// Second call: does nothing — once is "spent"
```

### Quick Check
> 1. What does `sync.Once.Do` guarantee?
> 2. If the function passed to `Do` panics, will `Do` call it again?
> 3. What is a common use case for `sync.Once`?

---

## 4. WaitGroup (Deeper Look)

We covered WaitGroup in Chapter 18; here's some deeper usage:

```go
// Pattern: Add outside, Done inside with defer
var wg sync.WaitGroup
for i, item := range items {
    wg.Add(1)
    go func(i int, item Item) {
        defer wg.Done()
        results[i] = process(item)
    }(i, item)
}
wg.Wait()

// Pattern: batch work with controlled concurrency
var wg sync.WaitGroup
sem := make(chan struct{}, 10)  // Semaphore: max 10 concurrent

for _, item := range items {
    wg.Add(1)
    go func(item Item) {
        defer wg.Done()
        sem <- struct{}{}        // Acquire semaphore slot
        defer func() { <-sem }() // Release on return
        process(item)
    }(item)
}
wg.Wait()
```

**WaitGroup for goroutine lifetime management:**
```go
type Server struct {
    wg sync.WaitGroup
}

func (s *Server) Start() {
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.runBackgroundWorker()
    }()
    
    s.wg.Add(1)
    go func() {
        defer s.wg.Done()
        s.runMetricsCollector()
    }()
}

func (s *Server) Stop() {
    // Signal goroutines to stop (via context or done channel)
    s.wg.Wait()  // Wait for all goroutines to exit
    log.Println("all goroutines stopped")
}
```

### Quick Check
> 1. What is a semaphore channel and how does it control concurrency?
> 2. Why is `wg.Add(1)` placed before `go func()`?
> 3. In `Server.Stop()`, what ensures the goroutines actually exit?

---

## 5. Cond — Condition Variable

`sync.Cond` lets goroutines wait for a condition to become true without busy-waiting:

```go
type Queue struct {
    mu    sync.Mutex
    cond  *sync.Cond
    items []int
}

func NewQueue() *Queue {
    q := &Queue{}
    q.cond = sync.NewCond(&q.mu)
    return q
}

func (q *Queue) Enqueue(v int) {
    q.mu.Lock()
    q.items = append(q.items, v)
    q.mu.Unlock()
    q.cond.Signal()  // Wake up one waiting goroutine
}

func (q *Queue) Dequeue() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    for len(q.items) == 0 {
        q.cond.Wait()  // Release lock + wait for signal, then re-acquire lock
    }
    
    v := q.items[0]
    q.items = q.items[1:]
    return v
}

func (q *Queue) Broadcast() {
    q.cond.Broadcast()  // Wake ALL waiting goroutines
}
```

**`cond.Wait()` releases the lock, sleeps, re-acquires the lock when woken.**

`sync.Cond` is rarely used in modern Go — channels and context usually work better. But it's useful for cases where many goroutines wait for a state change that's expensive to re-broadcast.

### Quick Check
> 1. What does `cond.Wait()` do to the mutex it holds?
> 2. What is the difference between `Signal()` and `Broadcast()`?
> 3. Why must `cond.Wait()` be inside a `for` loop?

---

## 6. Map — Concurrent Map

`sync.Map` is a map designed for concurrent use without external locking:

```go
var m sync.Map

// Store:
m.Store("key1", "value1")
m.Store("key2", 42)

// Load:
val, ok := m.Load("key1")
if ok {
    fmt.Println(val.(string))  // "value1"
}

// LoadOrStore — atomic: returns existing or stores new:
actual, loaded := m.LoadOrStore("key1", "default")
// loaded = true if key existed, false if we stored the new value

// Delete:
m.Delete("key2")

// Range — iterate (no guaranteed order):
m.Range(func(key, value any) bool {
    fmt.Println(key, value)
    return true  // Return false to stop iteration
})

// LoadAndDelete — atomic load + delete:
val, loaded := m.LoadAndDelete("key1")
```

**When to use `sync.Map` vs `map` + `RWMutex`:**

| | `sync.Map` | `map` + `RWMutex` |
|--|-----------|-------------------|
| Type safety | No (uses `any`) | Yes (generic map) |
| Performance | Better for stable key sets (append-only, few deletes) | Better for general use |
| API | Load/Store/Range | Direct map operations |
| Best for | Caches, registries, read-heavy stable data | General concurrent maps |

```go
// sync.Map is best when:
// 1. Keys are written once and read many times (like a registry)
// 2. Multiple goroutines modify disjoint sets of keys

// Regular map + mutex is best when:
// 1. You need type safety
// 2. Access patterns are complex (transactions across multiple keys)
```

### Quick Check
> 1. What is the return type of `sync.Map.Load()`?
> 2. What does `LoadOrStore` do atomically?
> 3. When is `sync.Map` faster than a regular map + RWMutex?

---

## 7. Pool — Object Reuse

`sync.Pool` is a cache of temporary objects that can be reused to reduce GC pressure:

```go
var bufPool = sync.Pool{
    New: func() any {
        return bytes.NewBuffer(make([]byte, 0, 256))
    },
}

func processRequest(data []byte) string {
    buf := bufPool.Get().(*bytes.Buffer)  // Get from pool (or create new)
    buf.Reset()                            // Reset before use!
    defer bufPool.Put(buf)                 // Return to pool when done
    
    buf.Write(data)
    buf.WriteString(" processed")
    return buf.String()
}
```

**Why it helps performance:**
```go
// Without pool: each request allocates a new buffer
for _, req := range requests {
    buf := bytes.NewBuffer(...)  // 1000 allocations = 1000 GC targets
    process(req, buf)
}

// With pool: buffers are reused
for _, req := range requests {
    buf := bufPool.Get().(*bytes.Buffer)  // Reuses from pool if available
    buf.Reset()
    defer bufPool.Put(buf)
    process(req, buf)
    // 1000 requests, but GC only sees ~numCPU buffers
}
```

**Rules for `sync.Pool`:**
1. Always `Reset()` the object before use — it might have data from the previous user
2. Never store state in pooled objects that must persist across calls
3. Don't use for resources that must be explicitly closed (like file handles)
4. Pool objects can be freed by the GC at any time — don't rely on them persisting

**Common pool patterns:**
```go
// Buffer pool:
var bufPool = sync.Pool{
    New: func() any { return new(bytes.Buffer) },
}

// JSON encoder pool:
var encoderPool = sync.Pool{
    New: func() any { return json.NewEncoder(io.Discard) },
}

// Struct pool:
type Request struct{ /* fields */ }
var reqPool = sync.Pool{
    New: func() any { return &Request{} },
}
```

### Quick Check
> 1. What happens if `sync.Pool.Get()` is called when the pool is empty?
> 2. Why must you reset pooled objects before use?
> 3. Can you use `sync.Pool` for database connections?

---

## 8. Atomic Operations

The `sync/atomic` package provides low-level atomic operations for simple values — faster than a mutex for a single variable:

```go
import "sync/atomic"

// Typed atomic types (Go 1.19+) — preferred:
var counter atomic.Int64
var flag    atomic.Bool
var ptr     atomic.Pointer[User]
var value   atomic.Value  // For arbitrary types

// Operations:
counter.Add(1)           // Atomic increment
counter.Add(-1)          // Atomic decrement
n := counter.Load()      // Atomic read
counter.Store(100)       // Atomic write
old := counter.Swap(42)  // Atomic swap, returns old value

// Compare-and-swap (CAS) — the fundamental atomic operation:
// "Set to newVal only if current value is still expectedVal"
swapped := counter.CompareAndSwap(expected, newVal)
if swapped {
    fmt.Println("successfully updated")
} else {
    fmt.Println("value changed between read and CAS — retry")
}
```

**Real-world use: atomic flag for shutdown:**
```go
type Worker struct {
    running atomic.Bool
}

func (w *Worker) Start() {
    w.running.Store(true)
    go func() {
        for w.running.Load() {
            doWork()
        }
    }()
}

func (w *Worker) Stop() {
    w.running.Store(false)
}
```

**Real-world use: atomic counter for metrics:**
```go
type Metrics struct {
    requests    atomic.Int64
    errors      atomic.Int64
    latencySum  atomic.Int64
}

func (m *Metrics) RecordRequest(latencyMs int64, err error) {
    m.requests.Add(1)
    m.latencySum.Add(latencyMs)
    if err != nil {
        m.errors.Add(1)
    }
}

func (m *Metrics) Stats() (reqs, errs, avgLatency int64) {
    reqs = m.requests.Load()
    errs = m.errors.Load()
    if reqs > 0 {
        avgLatency = m.latencySum.Load() / reqs
    }
    return
}
```

**`atomic.Value` for arbitrary types:**
```go
var config atomic.Value

// Writer:
config.Store(&Config{MaxConns: 100, Timeout: 30})

// Reader (hot path — no lock):
cfg := config.Load().(*Config)
doWork(cfg)

// Atomic config reload: 
// New config is written atomically — readers always see a complete config
config.Store(&Config{MaxConns: 200, Timeout: 60})
```

### Quick Check
> 1. When should you use atomic operations instead of a mutex?
> 2. What does `CompareAndSwap` do?
> 3. When would you use `atomic.Value`?

---

## 9. Choosing the Right Primitive

```
Question 1: Are you passing DATA between goroutines?
  YES → Use CHANNEL

Question 2: Are you protecting SHARED STATE?
  YES, simple counter/flag → Use ATOMIC
  YES, one variable, read-heavy → Use RWMUTEX
  YES, multi-field struct → Use MUTEX
  YES, concurrent map → Use SYNC.MAP or MAP+RWMUTEX

Question 3: Do you need to WAIT for goroutines?
  Finite set of goroutines → Use WAITGROUP
  Waiting for a condition → Use COND (rarely) or CHANNEL

Question 4: Do you need to RUN SOMETHING ONCE?
  YES → Use ONCE

Question 5: Do you need to REUSE OBJECTS?
  YES → Use POOL
```

**Performance comparison (rough order, fastest to slowest):**
```
Atomic operations  → ns range, no contention
sync.Mutex         → ~20ns uncontended, increases with contention
sync.RWMutex       → ~25ns for reads uncontended, better under high read load
Channel (buffered) → ~50ns per operation
Channel (unbuffered) → requires goroutine synchronization
```

---

## Summary

- **Mutex**: exclusive access to shared data — `Lock()`/`Unlock()` (always defer)
- **RWMutex**: multiple simultaneous readers, exclusive writers — use when reads >> writes
- **Once**: run initialization exactly once, thread-safe — ideal for singletons and lazy init
- **WaitGroup**: wait for N goroutines to finish — `Add` before `go`, `Done` via defer
- **Cond**: wait for a condition — rarely needed, channels usually better
- **sync.Map**: concurrent map without external locking — best for append-only/stable key sets
- **Pool**: reuse objects to reduce GC pressure — always reset before use
- **Atomic**: lock-free operations for single values — fastest, but limited to simple operations
- **Never copy mutexes** — pass pointers to structs containing mutexes

---

## Exercises

### Easy
1. Build a thread-safe `Set[T comparable]` using `sync.Mutex`. Methods: `Add(v T)`, `Remove(v T)`, `Contains(v T) bool`, `Size() int`, `Slice() []T`. Verify with 100 concurrent goroutines each adding a unique integer.
2. Implement `LazyLoader[T any]` using `sync.Once` — it takes an `init func() T` and calls it at most once (lazily on first access). `Load() T` returns the value. Verify from 10 concurrent goroutines that init is called exactly once.
3. Create a `HitCounter` using `atomic.Int64`. `Hit()` increments. `Count() int64` returns current count. `Reset() int64` atomically resets to 0 and returns the old value. Test with 100 concurrent goroutines hitting 1000 times each — verify final count is exactly 100,000.

### Medium
4. Read-heavy cache: Implement `Cache[K comparable, V any]` with `sync.RWMutex`. Methods: `Get(k K) (V, bool)`, `Set(k K, v V)`, `Delete(k K)`, `Len() int`. Add a `TTL` field — entries expire after a given duration. Write a goroutine that periodically evicts expired entries (every 30s). Benchmark with 10 readers and 1 writer vs `sync.Map` for the same workload.
5. Object pool with stats: Extend `sync.Pool` into a `TrackedPool[T any]` that tracks: total `Get()` calls, pool hits (object was reused), pool misses (new object was created), current pool size estimate. Methods: `Get() *T`, `Put(*T)`, `Stats() PoolStats`. Benchmark a JSON serialization workload showing allocation reduction.
6. Concurrent pipeline with backpressure: Build a pipeline with bounded stages. `Stage[T, R any](input <-chan T, workers int, fn func(T) R, bufSize int) <-chan R` runs `workers` goroutines processing from `input`. `bufSize` controls the output buffer (backpressure). Chain 3 stages: parse → transform → encode. Test that slow downstream stages cause upstream to pause (backpressure is working).

### Hard
7. Lock-free stack: Implement a lock-free stack using `atomic.Pointer` and Compare-and-Swap. `Push(v T)` adds to the top. `Pop() (T, bool)` removes from the top. Handle the ABA problem (hint: use a version counter in the node pointer). Benchmark against `sync.Mutex`-based stack under 1, 4, 8, 16, 32 concurrent goroutines. At what concurrency does the lock-free version win?
8. Rate limiter with burst: Implement a token bucket rate limiter. `NewRateLimiter(rate float64, burst int)` creates one allowing `rate` tokens/second with `burst` max bucket size. `Wait(ctx context.Context) error` blocks until a token is available or ctx is cancelled. `TryAcquire() bool` is non-blocking. Use `sync/atomic` for the token count and `sync.Cond` to wake waiters when tokens become available. Test: verify exact rate limiting under concurrent load, burst behavior, and context cancellation during wait.
