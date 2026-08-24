# Chapter 15: Go Runtime — Garbage Collector, Escape Analysis & Memory Management

Understanding how Go manages memory explains why certain code patterns are faster, why GC pauses happen, and how to write Go code that is kind to the garbage collector. This knowledge is consistently tested in senior interviews.

## Table of Contents

1. [Stack vs Heap in Go](#1-stack-vs-heap-in-go)
2. [Escape Analysis](#2-escape-analysis)
3. [Go's Garbage Collector](#3-gos-garbage-collector)
4. [Tri-Color Mark-and-Sweep](#4-tri-color-mark-and-sweep)
5. [Write Barriers](#5-write-barriers)
6. [GC Tuning — GOGC and GOMEMLIMIT](#6-gc-tuning)
7. [Reducing GC Pressure](#7-reducing-gc-pressure)
8. [Memory Profiling in Practice](#8-memory-profiling)
9. [Interview Questions & Model Answers](#9-interview-questions--model-answers)
10. [Summary](#summary)

---

## 1. Stack vs Heap in Go

Every goroutine has its own stack — a region of memory for local variables, function arguments, and return values. The heap is a shared pool for longer-lived allocations.

```
STACK: fast, managed automatically
  - Local variables that don't escape the function
  - Function arguments and return values
  - Small structs and arrays
  - Allocation: just move stack pointer (near-zero cost)
  - Deallocation: just move stack pointer back (zero cost)
  - Each goroutine has its own stack (starts 2KB, grows as needed)

HEAP: slower, requires garbage collection
  - Objects whose lifetime exceeds the function that created them
  - Objects whose address is taken and might outlive the function
  - Large allocations
  - Allocation: memory management overhead
  - Deallocation: garbage collector (periodic)
```

### Why It Matters

Stack allocations are essentially free — they're just pointer arithmetic. Heap allocations involve the allocator, increase GC pressure, and can cause GC pauses. Writing code that minimizes heap allocations is a key Go performance skill.

---

## 2. Escape Analysis

The Go compiler performs escape analysis at compile time to determine whether a variable can stay on the stack or must be moved ("escape") to the heap.

### Seeing Escape Analysis Results

```bash
go build -gcflags="-m" ./...
# OR for verbose output:
go build -gcflags="-m -m" ./...
```

```go
package main

func stackAlloc() int {
    n := 42        // stays on stack: doesn't escape
    return n
}

func heapAlloc() *int {
    n := 42        // escapes to heap: pointer to n is returned
    return &n      // "n escapes to heap" in compiler output
}

func interfaceEscape(v interface{}) {
    // storing in interface causes escape (type info + value pointer needed)
    _ = v
}
```

### Common Escape Patterns

```go
// 1. RETURNING A POINTER to a local variable — always escapes
func createUser() *User {
    u := User{Name: "Alice"} // escapes to heap
    return &u
}

// 2. STORING IN AN INTERFACE — escapes
var w io.Writer = &bytes.Buffer{} // Buffer escapes

// 3. GOROUTINE CLOSURE captures a variable — escapes
func startWorker() {
    count := 0  // escapes: referenced by goroutine
    go func() {
        count++ // goroutine may outlive the function
    }()
}

// 4. LARGE ALLOCATIONS — compiler may put large arrays on heap
func bigArray() [1000000]int { // may escape if too large for stack
    return [1000000]int{}
}

// STAYS ON STACK:
func small() {
    // Small values used only within the function
    x := 42
    y := x * 2
    _ = y
}
```

### Avoiding Unnecessary Heap Allocations

```go
// BAD: allocates a new User on the heap for every call
func processUser() *User {
    u := &User{Name: "Alice"} // heap allocation
    doSomething(u)
    return u
}

// GOOD: caller owns the User, pass by pointer to avoid copy
func processUser(u *User) {
    doSomething(u) // u was allocated by caller, may be on stack
}

// BAD: string formatting allocates
func logRequest(id string) {
    log.Printf("Processing request: " + id) // string concatenation allocates
}

// GOOD: use format verbs, let Printf handle it
func logRequest(id string) {
    log.Printf("Processing request: %s", id)
}
```

---

## 3. Go's Garbage Collector

Go uses a concurrent, tri-color mark-and-sweep garbage collector. "Concurrent" means it runs mostly while your program is running — not in a stop-the-world pause.

### GC Phases

```
1. MARK SETUP (stop-the-world, very brief ~0.1ms)
   - Turn on write barriers
   - Stack scanning setup

2. MARKING (concurrent with application)
   - Scan roots: goroutine stacks, global variables
   - Traverse the object graph, marking reachable objects
   - Application goroutines continue running
   - Write barriers track pointer mutations during this phase

3. MARK TERMINATION (stop-the-world, brief ~0.1ms)
   - Final scan of stacks that mutated during marking
   - Turn off write barriers
   - Calculate next GC target

4. SWEEPING (concurrent)
   - Free unmarked objects
   - Return memory to the OS eventually
```

### When GC Runs

The GC runs based on the **heap growth ratio** set by `GOGC`:

```
GC runs when: heap_size >= GOGC/100 * heap_size_after_last_GC

With GOGC=100 (default):
  After GC cleans up: heap = 10MB
  Next GC triggers when: heap grows to 20MB (100% growth)

With GOGC=200:
  Next GC at: 30MB (200% growth)
  → Less frequent GC, higher peak memory

With GOGC=50:
  Next GC at: 15MB (50% growth)
  → More frequent GC, lower peak memory, more CPU overhead
```

---

## 4. Tri-Color Mark-and-Sweep

The GC uses three colors for objects:

```
WHITE: not yet visited — might be garbage
GREY:  discovered but children not yet scanned
BLACK: fully scanned, all children are grey or black

Algorithm:
1. Start: all objects are white
2. Mark roots as grey (globals, stacks)
3. While grey set is not empty:
   - Pick a grey object
   - Mark all its white references as grey
   - Mark the object black
4. End: white objects have no references — they are garbage
5. Sweep: free all white (garbage) objects
```

```
Initial state:     root → A → B → C
All white:         □ → □ → □ → □

Mark root grey:    ⊡ → □ → □ → □

Process root:      ■ → ⊡ → □ → □
(root black, A grey)

Process A:         ■ → ■ → ⊡ → □
(A black, B grey)

Process B:         ■ → ■ → ■ → ⊡
(B black, C grey)

Process C:         ■ → ■ → ■ → ■

Orphan D (unreachable): □ ← still white, will be freed
```

---

## 5. Write Barriers

During concurrent marking, the application may modify the object graph. Write barriers ensure that newly created pointers are handled correctly.

```go
// When the application does:
a.ptr = b  // a now points to b

// The write barrier fires and marks b as grey (if it was white)
// This ensures b won't be freed even though the GC may have
// already scanned a before this pointer assignment
```

Write barriers have a small performance cost — this is why GC turns them off between GC cycles.

---

## 6. GC Tuning

### GOGC

```bash
# Default: 100 (GC runs when heap doubles)
GOGC=100 ./myapp

# Reduce GC frequency at cost of higher memory
GOGC=400 ./myapp  # GC only when heap grows 4x

# Disable GC entirely (dangerous — memory grows unbounded)
GOGC=off ./myapp

# In code:
import "runtime/debug"
debug.SetGCPercent(200) // equivalent to GOGC=200
```

### GOMEMLIMIT (Go 1.19+)

A hard limit on Go runtime memory. When the heap approaches this limit, the GC runs more aggressively.

```bash
# Set soft memory limit to 500MB
GOMEMLIMIT=500MiB ./myapp

# In code:
import "runtime/debug"
debug.SetMemoryLimit(500 * 1024 * 1024) // 500MB
```

**Best practice for containers:** Set `GOMEMLIMIT` to ~90% of the container memory limit. Without it, Go may allocate more memory than the container allows, causing OOM kills.

---

## 7. Reducing GC Pressure

### Technique 1: Reuse Objects with sync.Pool

```go
var bufPool = sync.Pool{
    New: func() interface{} { return make([]byte, 0, 4096) },
}

func processRequest(data []byte) {
    buf := bufPool.Get().([]byte)
    buf = buf[:0] // reset without reallocating
    defer bufPool.Put(buf)
    
    buf = append(buf, data...)
    process(buf)
}
// Without pool: every request allocates a new buffer → GC pressure
// With pool: buffers are reused across requests → less GC
```

### Technique 2: Pre-allocate Slices

```go
// BAD: many small reallocations as slice grows
result := []int{}
for i := 0; i < 1000000; i++ {
    result = append(result, i)
}

// GOOD: one allocation
result := make([]int, 0, 1000000)
for i := 0; i < 1000000; i++ {
    result = append(result, i)
}
```

### Technique 3: Use Value Types in Slices

```go
// BAD: slice of pointers — each element is a heap allocation
users := []*User{
    {Name: "Alice"},
    {Name: "Bob"},
}

// GOOD: slice of values — single allocation
users := []User{
    {Name: "Alice"},
    {Name: "Bob"},
}
// Better cache locality + less GC pressure
```

### Technique 4: Avoid Interface Allocations in Hot Paths

```go
// BAD: fmt.Sprintf always allocates a string
func log(id int) {
    message := fmt.Sprintf("processing id: %d", id) // allocation
    logger.Log(message)
}

// BETTER: use strconv for numbers to avoid fmt allocation
func log(id int) {
    logger.Log("processing id: " + strconv.Itoa(id))
}

// BEST: structured logging avoids string building entirely
slog.Info("processing", "id", id) // slog avoids most allocations
```

---

## 8. Memory Profiling

```go
import (
    "os"
    "runtime/pprof"
)

// Write heap profile to file
f, _ := os.Create("heap.prof")
defer f.Close()
pprof.WriteHeapProfile(f)

// Or use the HTTP endpoint (in dev):
import _ "net/http/pprof"
// Then: go tool pprof http://localhost:6060/debug/pprof/heap
```

```bash
# Analyze a heap profile
go tool pprof heap.prof
(pprof) top          # top allocating functions
(pprof) list myFunc  # show line-by-line allocations in myFunc
(pprof) web          # open flame graph in browser

# View allocations during a benchmark
go test -bench=. -benchmem ./...
# Output:
# BenchmarkProcess-8   1000000   1234 ns/op   256 B/op   4 allocs/op
# B/op: bytes per operation
# allocs/op: heap allocations per operation
```

---

## 9. Interview Questions & Model Answers

**Q: How does Go's garbage collector work?**

"Go uses a concurrent tri-color mark-and-sweep GC. It has two short stop-the-world pauses — one to set up marking (turn on write barriers, ~0.1ms) and one to finalize marking (~0.1ms). The actual marking — traversing the object graph and identifying live objects — runs concurrently with the application. Write barriers ensure that any pointer changes during concurrent marking are tracked so no live objects are missed. After marking, sweeping (freeing garbage) also runs concurrently. Typical GC pause times in modern Go are under 1 millisecond."

**Q: What is escape analysis and why does it matter?**

"Escape analysis is a compiler optimization that determines whether a variable can be allocated on the goroutine's stack rather than the heap. Stack allocation is essentially free — just a pointer increment — and requires no GC. Heap allocation adds GC pressure and allocation overhead. A variable 'escapes' to the heap when its address outlives the function that created it — for example, when you return a pointer to a local variable, or when a goroutine closure captures a variable. You can see escape analysis decisions with `go build -gcflags='-m'`."

**Q: How would you reduce GC pressure in a high-throughput Go service?**

"Several techniques: First, use sync.Pool to reuse frequently allocated objects like buffers. Second, pre-allocate slices with the right capacity to avoid repeated reallocations. Third, prefer value types over pointer types in data structures to improve locality and reduce pointer scanning work for the GC. Fourth, set GOMEMLIMIT appropriate to your container's memory limit to prevent OOM kills. Fifth, profile first with pprof — don't optimize blindly. Focus on the allocating paths shown in `top allocs` in the heap profile."

---

## Summary

- Go allocates small, short-lived variables on the stack (fast, free). Long-lived or shared objects go on the heap (GC overhead).
- **Escape analysis:** compiler decides stack vs heap. Use `go build -gcflags="-m"` to see decisions.
- **GC:** concurrent tri-color mark-and-sweep. Two short stop-the-world pauses (~0.1ms each). Marking runs concurrently.
- **Write barriers:** ensure pointer updates during concurrent marking are tracked.
- **GOGC:** controls GC frequency. Default 100 = GC when heap doubles. Higher = less frequent but higher peak memory.
- **GOMEMLIMIT:** hard memory limit. Essential for containerized applications.
- **Reduce GC pressure:** sync.Pool for reuse, pre-allocate slices, prefer value types, avoid interface conversions in hot paths.
