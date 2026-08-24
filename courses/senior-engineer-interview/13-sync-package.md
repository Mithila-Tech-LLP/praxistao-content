# Chapter 13: The sync Package — Mutex, RWMutex, WaitGroup, Once & Atomic

The `sync` package is the foundation of concurrent Go programs. Every senior Go interview will touch these primitives. Understanding not just how to use them, but when and why, is what separates senior engineers from mid-level engineers.

## Table of Contents

1. [sync.Mutex — Exclusive Access](#1-syncmutex--exclusive-access)
2. [sync.RWMutex — Read/Write Locking](#2-syncrwmutex--readwrite-locking)
3. [sync.WaitGroup — Coordinating Goroutines](#3-syncwaitgroup--coordinating-goroutines)
4. [sync.Once — Safe Initialization](#4-synconce--safe-initialization)
5. [sync.Cond — Condition Variables](#5-synccond--condition-variables)
6. [sync.Map — Concurrent Map](#6-syncmap--concurrent-map)
7. [sync.Pool — Object Reuse](#7-syncpool--object-reuse)
8. [Atomic Operations](#8-atomic-operations)
9. [Common Mistakes](#9-common-mistakes)
10. [Interview Questions & Model Answers](#10-interview-questions--model-answers)
11. [Summary](#summary)

---

## 1. sync.Mutex — Exclusive Access

A mutex (mutual exclusion lock) ensures only one goroutine at a time can access a critical section.

```go
type SafeCounter struct {
    mu sync.Mutex
    count int
}

func (c *SafeCounter) Increment() {
    c.mu.Lock()
    defer c.mu.Unlock() // always defer unlock to handle panics correctly
    c.count++
}

func (c *SafeCounter) Value() int {
    c.mu.Lock()
    defer c.mu.Unlock()
    return c.count
}
```

### Critical Rules for Mutexes

```go
// RULE 1: Always use defer unlock to prevent lock leaks from panics
func badExample(mu *sync.Mutex) {
    mu.Lock()
    // if this panics, mu is never unlocked — deadlock!
    doSomethingRisky()
    mu.Unlock() // might not execute
}

func goodExample(mu *sync.Mutex) {
    mu.Lock()
    defer mu.Unlock() // always executes, even on panic
    doSomethingRisky()
}

// RULE 2: Never copy a mutex after first use
type Bad struct {
    mu sync.Mutex
    data int
}
func process(b Bad) { // BAD: b is a copy, mu is copied too!
    b.mu.Lock()
    // ...
}
func processSafe(b *Bad) { // GOOD: pointer, no copy
    b.mu.Lock()
    // ...
}

// RULE 3: Keep critical sections as short as possible
func goodCriticalSection(c *SafeCounter, expensive func() int) {
    result := expensive() // do expensive work OUTSIDE the lock
    c.mu.Lock()
    c.count += result     // only update shared state under lock
    c.mu.Unlock()
}
```

---

## 2. sync.RWMutex — Read/Write Locking

An RWMutex allows concurrent reads but exclusive writes. Use when reads are much more frequent than writes.

```go
type Cache struct {
    mu    sync.RWMutex
    items map[string]string
}

// Multiple goroutines can call Get simultaneously
func (c *Cache) Get(key string) (string, bool) {
    c.mu.RLock()         // acquire read lock
    defer c.mu.RUnlock()
    v, ok := c.items[key]
    return v, ok
}

// Only one goroutine can call Set at a time; blocks all readers
func (c *Cache) Set(key, value string) {
    c.mu.Lock()         // acquire write lock
    defer c.mu.Unlock()
    c.items[key] = value
}
```

### When RWMutex is Faster Than Mutex

```go
// With sync.Mutex: readers queue up even though they don't conflict
// With sync.RWMutex: multiple readers proceed in parallel

// Rule of thumb: use RWMutex when reads >> writes (e.g., read 95%, write 5%)
// Benchmark: RWMutex can be 2-10x faster for read-heavy workloads
```

---

## 3. sync.WaitGroup — Coordinating Goroutines

WaitGroup waits for a collection of goroutines to finish. The standard pattern for fan-out and fan-in.

```go
func processAll(items []Item) {
    var wg sync.WaitGroup

    for _, item := range items {
        wg.Add(1)  // MUST call Add before starting the goroutine
        go func(i Item) {
            defer wg.Done() // signal completion when goroutine exits
            process(i)
        }(item) // pass item as argument to avoid closure variable capture bug
    }

    wg.Wait() // blocks until all goroutines call Done()
}
```

### Common WaitGroup Mistakes

```go
// MISTAKE 1: Calling Add inside the goroutine (race condition)
for _, item := range items {
    go func(i Item) {
        wg.Add(1) // WRONG: wg.Wait() might return before this runs
        defer wg.Done()
        process(i)
    }(item)
}

// MISTAKE 2: Not passing item by value (closure captures loop variable)
for _, item := range items {
    wg.Add(1)
    go func() {
        defer wg.Done()
        process(item) // WRONG: 'item' is the loop variable, may have changed
    }()
}

// MISTAKE 3: Reusing WaitGroup before Wait() returns
wg.Wait()
// Still safe to reuse after Wait() returns:
wg.Add(1) // OK: Wait has returned
```

### WaitGroup with Error Collection

```go
func processAllWithErrors(items []Item) ([]Result, error) {
    var wg sync.WaitGroup
    var mu sync.Mutex
    var results []Result
    var firstErr error

    for _, item := range items {
        wg.Add(1)
        go func(i Item) {
            defer wg.Done()
            result, err := process(i)
            mu.Lock()
            if err != nil && firstErr == nil {
                firstErr = err
            }
            if err == nil {
                results = append(results, result)
            }
            mu.Unlock()
        }(item)
    }
    wg.Wait()
    return results, firstErr
}
```

---

## 4. sync.Once — Safe Initialization

`sync.Once` guarantees a function is called exactly once, even from multiple goroutines. Essential for lazy initialization.

```go
type Singleton struct {
    value string
}

var (
    instance *Singleton
    once     sync.Once
)

func GetInstance() *Singleton {
    once.Do(func() {
        instance = &Singleton{value: "initialized"}
        // This runs exactly once, even if GetInstance is called concurrently
    })
    return instance
}
```

### sync.Once for Expensive One-Time Operations

```go
type Config struct {
    once   sync.Once
    loaded map[string]string
    err    error
}

func (c *Config) Load(path string) (map[string]string, error) {
    c.once.Do(func() {
        // This expensive operation runs only once
        c.loaded, c.err = loadFromFile(path)
    })
    return c.loaded, c.err
}

// Limitation: once.Do doesn't retry if the function panics or errors.
// If you need retry-on-error, use atomic.Value with double-checked locking instead.
```

---

## 5. sync.Cond — Condition Variables

`sync.Cond` allows goroutines to wait for a condition to become true. Less common but important for certain patterns.

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

func (q *Queue) Push(item int) {
    q.mu.Lock()
    q.items = append(q.items, item)
    q.cond.Signal()  // wake up one waiting goroutine
    q.mu.Unlock()
}

func (q *Queue) Pop() int {
    q.mu.Lock()
    defer q.mu.Unlock()
    
    // MUST use a loop — spurious wakeups are possible
    for len(q.items) == 0 {
        q.cond.Wait() // releases lock, waits for Signal/Broadcast, reacquires lock
    }
    
    item := q.items[0]
    q.items = q.items[1:]
    return item
}

// Broadcast wakes all waiting goroutines (use when condition change affects all)
// Signal wakes one (use when one waiter can make progress)
```

---

## 6. sync.Map — Concurrent Map

`sync.Map` is a concurrent map optimized for two access patterns: when each key is written once and read many times, or when many goroutines write to disjoint keys.

```go
var m sync.Map

// Store
m.Store("key", "value")

// Load
val, ok := m.Load("key")
if ok {
    fmt.Println(val.(string))
}

// LoadOrStore: atomically load if exists, store if not
actual, loaded := m.LoadOrStore("key", "new_value")
// loaded=true if key already existed, actual=existing value
// loaded=false if stored new value, actual=new_value

// Delete
m.Delete("key")

// Range: iterate (not safe to modify during Range)
m.Range(func(key, value interface{}) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true // return false to stop iteration
})
```

### sync.Map vs map + RWMutex

```go
// sync.Map: better when:
// - Write once, read many (e.g., cache populated at startup)
// - Many goroutines write to disjoint keys (no single lock bottleneck)

// map + RWMutex: better when:
// - You need iteration with consistent snapshot
// - Keys are frequently written by same goroutines
// - You need additional operations (len, copy)
// - Simpler code is more important than marginal performance gains
```

---

## 7. sync.Pool — Object Reuse

`sync.Pool` is a set of temporary objects that can be reused to reduce GC pressure. Used for objects that are expensive to allocate and frequently created/discarded.

```go
var bufferPool = sync.Pool{
    New: func() interface{} {
        buf := make([]byte, 0, 1024)
        return &buf
    },
}

func processRequest(data []byte) {
    // Get a buffer from the pool (or create new if pool is empty)
    bufPtr := bufferPool.Get().(*[]byte)
    buf := (*bufPtr)[:0] // reset without reallocating

    defer func() {
        *bufPtr = buf
        bufferPool.Put(bufPtr) // return to pool when done
    }()

    // Use buf...
    buf = append(buf, data...)
    process(buf)
}

// IMPORTANT: Pool objects may be garbage collected at any time.
// Never store anything in a Pool that has destructors or finalizers.
// Never store goroutines in a Pool.
// sync.Pool is for reducing allocation pressure, not for connection pools.
```

---

## 8. Atomic Operations

`sync/atomic` provides low-level atomic operations on integers and pointers. Faster than a mutex for simple counter operations.

```go
import "sync/atomic"

var counter int64

// Atomic increment
atomic.AddInt64(&counter, 1)

// Atomic read (safe without a mutex)
val := atomic.LoadInt64(&counter)

// Atomic write
atomic.StoreInt64(&counter, 0)

// Compare-and-swap: if *addr == old, set *addr = new, return true
swapped := atomic.CompareAndSwapInt64(&counter, old, new)

// atomic.Value: for storing any type atomically
var config atomic.Value
config.Store(Config{maxConns: 100}) // store any value
cfg := config.Load().(Config)       // load atomically
```

### Atomic vs Mutex

```go
// Use atomic when: single value (counter, flag, pointer)
// Use mutex when: multiple related values that must change together

// ATOMIC for a counter: clean and fast
var totalRequests int64
atomic.AddInt64(&totalRequests, 1)

// MUTEX for related state: both must update together
type Stats struct {
    mu           sync.Mutex
    totalRequests int64
    totalErrors   int64
}
func (s *Stats) RecordError() {
    s.mu.Lock()
    s.totalRequests++
    s.totalErrors++ // these two must be consistent
    s.mu.Unlock()
}
```

---

## 9. Common Mistakes

### Mistake 1: Forgetting to Pass Counter to WaitGroup Before goroutine

```go
// WRONG: race condition — main may exit before Add(1) executes
for _, item := range items {
    go func(i Item) {
        wg.Add(1) // too late
        defer wg.Done()
    }(item)
}
wg.Wait()

// CORRECT:
for _, item := range items {
    wg.Add(1) // before go
    go func(i Item) {
        defer wg.Done()
    }(item)
}
```

### Mistake 2: Copying a Mutex

```go
type Counter struct {
    mu    sync.Mutex
    value int
}

func process(c Counter) { // WRONG: copies the mutex!
    c.mu.Lock()
}

func process(c *Counter) { // CORRECT: pointer
    c.mu.Lock()
}
```

### Mistake 3: Unlocking in a Different Goroutine

```go
var mu sync.Mutex
mu.Lock()
go func() {
    // do work
    mu.Unlock() // WRONG: unlock from different goroutine — undefined behavior
}()
```

---

## 10. Interview Questions & Model Answers

**Q: When would you use sync.RWMutex instead of sync.Mutex?**

"When reads significantly outnumber writes. An RWMutex allows multiple concurrent readers as long as there are no writers. A regular Mutex only allows one goroutine at a time regardless of whether it's reading or writing. For a cache that's read thousands of times per second but updated once per minute, RWMutex can give 10-100x better read throughput. The trade-off: RWMutex has higher overhead per operation than Mutex, so if write frequency is high or the critical section is very short, the simpler Mutex may actually be faster."

**Q: What's the difference between sync.Once and an init function?**

"init() runs once per package automatically at program startup, before main(). sync.Once runs exactly once when first called, lazily — not at program startup. sync.Once is for lazy initialization (initialize expensive resources only when first needed) and for initialization in library code where you don't control startup order. Also, sync.Once can be called from concurrent goroutines safely."

**Q: When should you use atomic operations vs a mutex?**

"Atomic for simple, single-value operations — incrementing a counter, setting a flag, doing a compare-and-swap on a pointer. They're lock-free and faster. Mutex for operations that need to atomically update multiple values that must remain consistent with each other. If you have a struct where two fields must always be updated together, a mutex ensures they're never observed in a partially-updated state. Atomics don't give you that guarantee across multiple variables."

---

## Summary

- `sync.Mutex`: exclusive access. Always defer Unlock. Never copy. Keep critical sections short.
- `sync.RWMutex`: multiple concurrent readers OR one writer. Use for read-heavy shared state.
- `sync.WaitGroup`: coordinate goroutine completion. Call `Add` before starting goroutines, not inside.
- `sync.Once`: guaranteed single execution. Essential for lazy initialization.
- `sync.Cond`: condition variables. Always use `Wait()` in a loop — spurious wakeups are possible.
- `sync.Map`: concurrent map for write-once/read-many or disjoint-key access patterns.
- `sync.Pool`: object reuse to reduce GC pressure. Objects can be collected at any time.
- `sync/atomic`: lock-free single-value operations. Use for counters and flags; use mutex for related multi-value state.
