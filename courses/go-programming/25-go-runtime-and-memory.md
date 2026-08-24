# Chapter 23: The Go Runtime and Memory Model

Go programs run inside the Go runtime — a sophisticated piece of software that manages goroutines, memory, and the garbage collector. Understanding how the runtime works doesn't just satisfy curiosity; it explains why certain code is fast or slow, why goroutines are cheap, and how to write code the GC can handle efficiently.

## Table of Contents

1. [The Go Runtime](#1-the-go-runtime)
2. [Memory Layout — Stack and Heap](#2-memory-layout--stack-and-heap)
3. [Escape Analysis](#3-escape-analysis)
4. [The Garbage Collector](#4-the-garbage-collector)
5. [The Go Memory Model](#5-the-go-memory-model)
6. [Optimizing for the GC](#6-optimizing-for-the-gc)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. The Go Runtime

The Go runtime is linked into every Go binary. It provides:

```
Go Runtime
├── Goroutine scheduler (G-M-P model)
├── Memory allocator (tcmalloc-inspired, size classes)
├── Garbage collector (concurrent tri-color mark-and-sweep)
├── Stack management (growable stacks)
├── Channel and sync primitives
├── Signal handling
└── Runtime introspection (runtime package)
```

**The G-M-P model (Goroutines, Machine threads, Processors):**
```
G (Goroutine): the unit of work — your code
M (Machine): OS thread — there are GOMAXPROCS of these running Go
P (Processor): virtual CPU — has a local run queue of Gs

Each P has a local run queue.
Global run queue holds Gs not yet assigned to a P.
Idle Ms wait for work. Work-stealing: idle P steals from busy P's queue.

G1 →                  P1 → M1 (running on CPU core 1)
G2, G3 → local queue ↗
G4 →                  P2 → M2 (running on CPU core 2)
G5, G6 → local queue ↗
```

**Runtime introspection:**
```go
import "runtime"

// CPU/Goroutine info:
fmt.Println(runtime.GOMAXPROCS(0))     // Number of OS threads
fmt.Println(runtime.NumCPU())           // Number of CPU cores
fmt.Println(runtime.NumGoroutine())     // Current goroutine count

// Memory stats:
var m runtime.MemStats
runtime.ReadMemStats(&m)
fmt.Printf("Heap alloc: %d bytes\n", m.HeapAlloc)
fmt.Printf("Total alloc: %d bytes\n", m.TotalAlloc)
fmt.Printf("Sys: %d bytes\n", m.Sys)
fmt.Printf("GC cycles: %d\n", m.NumGC)
fmt.Printf("GC pause total: %s\n", time.Duration(m.PauseTotalNs))

// Force GC (only for testing/benchmarks):
runtime.GC()

// Print stack trace of all goroutines:
buf := make([]byte, 1<<20)
n := runtime.Stack(buf, true)
fmt.Printf("%s", buf[:n])
```

### Quick Check
> 1. What are G, M, and P in Go's scheduler?
> 2. What does `runtime.GOMAXPROCS(0)` return?
> 3. What does `runtime.NumGoroutine()` tell you?

---

## 2. Memory Layout — Stack and Heap

**Stack:**
- Each goroutine has its own stack
- Starts at 2KB, grows dynamically (up to 1GB by default)
- Fast allocation (just move stack pointer)
- Automatically freed when the function returns

**Heap:**
- Shared across all goroutines
- Allocation is slower (GC must track it)
- Lives until the GC collects it

```go
func stackExample() {
    x := 42         // x lives on the stack — allocated and freed instantly
    y := "hello"    // y (the string header) lives on the stack
                    // The string data lives on the heap
}

func heapExample() *int {
    x := 42
    return &x  // x MUST live on the heap — it outlives the function
               // Go's escape analysis detects this and allocates x on heap
}
```

**Growing stacks:**
```go
// Go stacks grow automatically — you never overflow a goroutine's stack
// (unlike C where stack overflow is undefined behavior)

func recursive(n int) {
    if n == 0 {
        return
    }
    // Each call adds a frame. Go doubles the stack when needed.
    // Small goroutines can grow to large stacks and shrink back.
    recursive(n - 1)
}
recursive(1000000)  // Fine! Stack grows as needed.
```

### Quick Check
> 1. How large does a goroutine's stack start?
> 2. Where does a local variable live if you return a pointer to it?
> 3. What happens when a goroutine's stack runs out of space?

---

## 3. Escape Analysis

**Escape analysis** determines whether a variable lives on the stack or heap at compile time. Variables "escape" to the heap when they outlive the function that created them:

```go
// Escapes to heap (returned as pointer — outlives function):
func newUser() *User {
    u := User{Name: "Alice"}  // u escapes to heap
    return &u
}

// Stays on stack (value is copied out, original doesn't need to persist):
func getUser() User {
    u := User{Name: "Alice"}  // u stays on stack
    return u  // Copy goes to caller
}

// Escapes via interface (concrete type unknown at runtime):
func printAnything(v any) {
    fmt.Println(v)
}
x := 42
printAnything(x)  // x escapes to heap — boxing for interface
```

**Viewing escape analysis:**
```bash
go build -gcflags='-m' main.go
# Output:
# ./main.go:5:2: moved to heap: u
# ./main.go:15:14: x escapes to heap
# ./main.go:15:14: main ... argument does not escape
```

**Common escape causes:**
```go
// 1. Returning a pointer:
func newFoo() *Foo { f := Foo{}; return &f }  // f escapes

// 2. Assigning to interface:
var i interface{} = 42  // 42 escapes (boxing)

// 3. Sending to channel:
ch <- someStruct  // someStruct may escape

// 4. Closure captures:
x := 42
f := func() { fmt.Println(x) }  // x may escape (closure outlives x's scope)

// 5. Slice/map backed by heap:
s := make([]int, 1000)  // slice backing array on heap (too large for stack)
```

**Does escape analysis matter for correctness?** No — Go handles it automatically. It matters for performance: heap allocations are slower and create GC work.

### Quick Check
> 1. What does it mean for a variable to "escape to the heap"?
> 2. What compiler flag shows escape analysis output?
> 3. Give two reasons why a local variable might escape to the heap.

---

## 4. The Garbage Collector

Go uses a **concurrent tri-color mark-and-sweep** garbage collector. Since Go 1.5, GC pauses are typically under 1ms.

**How it works:**
```
Tri-color marking:
  White: not yet visited (potential garbage)
  Grey:  visited, but children not yet scanned
  Black: visited, all children scanned (definitely live)

1. Mark (concurrent with program): 
   - Start from roots (goroutine stacks, global vars)
   - Turn white → grey → black as objects are scanned
   - Runs concurrently — your program keeps running!

2. Stop The World (STW) — very brief:
   - Scan any grey objects missed during concurrent mark
   - Typical STW: < 1ms

3. Sweep (concurrent):
   - Reclaim white (unreachable) objects
   - Runs concurrently

4. Scavenge:
   - Return unused memory to OS
```

**GC pacing — the GOGC variable:**
```bash
GOGC=100   # Default: run GC when heap doubles (100% growth)
GOGC=200   # Less frequent GC (more memory used, less CPU overhead)
GOGC=50    # More frequent GC (less memory, more CPU overhead)
GOGC=off   # Disable GC (careful!)
```

**Controlling GC in code:**
```go
import "runtime"
import "runtime/debug"

// Adjust GC target:
debug.SetGCPercent(200)  // Same as GOGC=200

// Set max total heap size (Go 1.19+):
debug.SetMemoryLimit(1 << 30)  // 1GB max — GC runs more aggressively near limit

// Get GC statistics:
var stats debug.GCStats
debug.ReadGCStats(&stats)
fmt.Printf("Last GC pause: %s\n", stats.PauseQuantiles[len(stats.PauseQuantiles)-1])
```

**What triggers a GC cycle:**
1. Heap has grown to `GOGC`% of the previous live heap size
2. 2 minutes have passed since the last GC
3. You call `runtime.GC()` explicitly

### Quick Check
> 1. What does "concurrent" GC mean?
> 2. What does `GOGC=200` do compared to the default `GOGC=100`?
> 3. What are the three main phases of Go's GC?

---

## 5. The Go Memory Model

The **Go Memory Model** defines when writes in one goroutine are visible to reads in another. This is critical for writing correct concurrent code.

**The golden rule: "If you read a variable in one goroutine that another goroutine might be writing, protect access with a mutex, channel, or atomic."**

**Happens-before relationship:**
```go
// Single goroutine: order is guaranteed
x := 0
x = 1     // Happens before
x = 2     // Happens before
fmt.Println(x)  // Sees x=2

// Multiple goroutines WITHOUT synchronization: NO guarantee
var x int
go func() { x = 1 }()    // Might not be visible immediately
fmt.Println(x)            // Might see 0 or 1 — UNDEFINED
```

**Synchronization guarantees (happens-before):**
```go
// 1. Channel send happens before receive:
ch := make(chan int, 1)
x := 0
go func() {
    x = 42       // (1)
    ch <- 1      // (2) — send happens after (1)
}()
<-ch             // (3) — receive happens after (2)
fmt.Println(x)   // SAFE: sees x=42 — (1) happens before (3)

// 2. Mutex unlock happens before next lock:
var mu sync.Mutex
x := 0
mu.Lock()
x = 42
mu.Unlock()       // Unlock happens before...
// ... another goroutine's mu.Lock() — and that goroutine sees x=42

// 3. Once.Do completion happens before any subsequent call returns:
var once sync.Once
x := 0
once.Do(func() { x = 42 })  // All future once.Do calls see x=42

// 4. WaitGroup.Wait returns after all Done calls:
var wg sync.WaitGroup
x := 0
wg.Add(1)
go func() {
    x = 42
    wg.Done()  // Happens before wg.Wait() returns
}()
wg.Wait()
fmt.Println(x)  // SAFE: sees x=42
```

**Data race example:**
```go
var count int

// BAD: two goroutines read and write count concurrently
go func() { count++ }()
go func() { count++ }()
// Result: undefined. Could be 1 or 2, or memory corruption.
```

### Quick Check
> 1. What is "happens-before"?
> 2. Is reading a variable written by another goroutine safe without synchronization?
> 3. What does "a channel send happens before the receive" guarantee?

---

## 6. Optimizing for the GC

**Reduce allocations — profile first:**
```bash
go test -bench=. -benchmem          # See allocs/op
go tool pprof -alloc_space          # See where allocations happen
```

**Technique 1: Reuse with sync.Pool:**
```go
var bufPool = sync.Pool{
    New: func() any { return make([]byte, 0, 4096) },
}

func processRequest(data []byte) {
    buf := bufPool.Get().([]byte)
    buf = buf[:0]        // Reset length, keep capacity
    defer bufPool.Put(buf)
    
    buf = append(buf, data...)
    // ... process buf ...
}
```

**Technique 2: Pre-allocate slices:**
```go
// BAD: appends cause multiple reallocations
var result []Item
for _, v := range items {
    result = append(result, process(v))
}

// GOOD: pre-allocate if you know the size
result := make([]Item, 0, len(items))
for _, v := range items {
    result = append(result, process(v))
}
```

**Technique 3: Avoid interface boxing:**
```go
// BAD: boxing int into interface on every call
func printNum(v interface{}) { fmt.Println(v) }
for i := 0; i < 1000000; i++ {
    printNum(i)  // Each call boxes i — heap allocation
}

// GOOD: use concrete type
func printNum(v int) { fmt.Println(v) }
for i := 0; i < 1000000; i++ {
    printNum(i)  // No boxing
}
```

**Technique 4: Strings.Builder for string concatenation:**
```go
// BAD: each + creates a new string (O(n²) total):
s := ""
for i := 0; i < 1000; i++ {
    s += strconv.Itoa(i)
}

// GOOD: strings.Builder reuses one buffer:
var sb strings.Builder
for i := 0; i < 1000; i++ {
    sb.WriteString(strconv.Itoa(i))
}
s := sb.String()
```

**Technique 5: Avoid unnecessary pointer indirection:**
```go
// If a struct is small and doesn't need to be shared/modified through pointer,
// pass by value — stays on stack, no GC involvement:
type Point struct{ X, Y float64 }

func distance(a, b Point) float64 {  // Pass by value — no heap allocation
    dx, dy := a.X-b.X, a.Y-b.Y
    return math.Sqrt(dx*dx + dy*dy)
}
```

---

## Summary

- **Runtime**: linked into every binary; manages goroutines (G-M-P model), memory, GC
- **G-M-P**: G=goroutine, M=OS thread, P=processor with local run queue; work stealing
- **Stack**: per-goroutine, 2KB start, grows automatically, fast
- **Heap**: shared, GC-managed, slower; variables escape when they outlive their function
- **Escape analysis**: `go build -gcflags='-m'`; heap escapes when: pointer returned, interface assigned, closure captures, large allocation
- **GC**: concurrent tri-color mark-and-sweep; STW < 1ms; triggered at `GOGC`% heap growth
- **`GOGC`**: controls GC frequency; `GOMEMLIMIT` (Go 1.19) caps total heap
- **Memory model**: writes visible to other goroutines only through synchronization (mutex, channel, atomic, WaitGroup)
- **Optimizations**: sync.Pool, pre-allocate slices, avoid interface boxing, strings.Builder, pass small structs by value

---

## Exercises

### Easy
1. Write a program that creates 100,000 goroutines, each sleeping for 100ms. Print `runtime.NumGoroutine()` before, during, and after. Then measure memory usage with `runtime.ReadMemStats`. How much does each goroutine cost?
2. Use `go build -gcflags='-m'` to examine escape analysis on a file with 5 functions: returning a value, returning a pointer, assigning to interface, sending to channel, and a closure capture. Record which variables escape and explain why each does or doesn't.
3. Benchmark two implementations of string joining: (a) `+=` in a loop, (b) `strings.Builder`. Use `b.N` and `-benchmem`. How many allocations does each approach make for 1000 joins?

### Medium
4. GC pressure profiler: Write a program that simulates three workloads: (a) many small short-lived allocations, (b) few large long-lived allocations, (c) reusing objects with sync.Pool. For each workload, measure: `NumGC`, `PauseTotalNs`, `HeapAlloc` after 5 seconds. Write a comparison report showing how each pattern affects GC behavior.
5. Escape-free hot path: Take this slow function (provided) that processes 1M events: it allocates a `*Event` per call and passes it through an `EventHandler interface`. Rewrite it to eliminate all allocations on the hot path using: value types instead of pointers, concrete types instead of interfaces, sync.Pool for any unavoidable allocations. Benchmark before and after — show allocs/op drops to 0.
6. Memory model verifier: Write a test that demonstrates a data race WITHOUT the race detector catching it (using carefully timed operations). Then write the same test WITH proper synchronization (mutex, channel, or atomic) and show the race disappears. Include a section explaining WHY each synchronization mechanism establishes the happens-before guarantee.

### Hard
7. Custom allocator: Implement a slab allocator `SlabAllocator[T any]` that pre-allocates a pool of objects in one large slice. `Alloc() *T` returns a pointer into the slab. `Free(*T)` marks the slot as available. Use a free-list (linked list through the slab). Benchmark against `new(T)` for 1M alloc/free cycles. Show GC statistics before and after — does the custom allocator reduce GC pauses?
8. GC tuning study: Build a web server that handles 10,000 requests/second, each request allocating ~100KB of temporary data. Benchmark under three configurations: (a) default `GOGC=100`, (b) `GOGC=400`, (c) `GOGC=100` + `GOMEMLIMIT=512MB`. Measure: p50/p95/p99 latency, GC pause time, total memory usage. Write a recommendation for which configuration is best and why.
