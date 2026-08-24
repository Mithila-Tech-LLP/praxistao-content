# Chapter 18: Go Performance — pprof, Benchmarks & Profiling

Performance engineering is a senior-level skill. You're expected to know not just how to write fast code, but how to measure what's actually slow, identify where time is spent, and make targeted improvements. This chapter covers the full Go profiling workflow.

## Table of Contents

1. [The Performance Engineering Mindset](#1-the-mindset)
2. [Writing Benchmarks](#2-writing-benchmarks)
3. [CPU Profiling](#3-cpu-profiling)
4. [Memory Profiling](#4-memory-profiling)
5. [Tracing](#5-tracing)
6. [Practical Optimizations](#6-practical-optimizations)
7. [Interview Questions & Model Answers](#7-interview-questions--model-answers)
8. [Summary](#summary)

---

## 1. The Mindset

**Don't optimize blind.** The first rule of performance optimization is: measure first, optimize second. Premature optimization wastes time on the wrong things.

**The profiling workflow:**
1. Write correct code
2. Measure under realistic load
3. Profile to find the actual bottleneck
4. Optimize the bottleneck
5. Measure again to verify improvement
6. Repeat from step 2

**Where time actually goes in Go services:**
- I/O wait (database, network) — often 90%+ of wall time
- Memory allocation and GC — often the #1 CPU bottleneck
- Lock contention — serializes concurrent code
- Compute (CPU-bound) — actual computation

---

## 2. Writing Benchmarks

Go's `testing` package includes benchmarking support. Always write benchmarks before optimizing.

```go
// Benchmark function names start with Benchmark
func BenchmarkMyFunction(b *testing.B) {
    // Setup code here (runs once, not timed)
    input := generateInput(1000)
    
    b.ResetTimer() // reset timer after setup
    
    for i := 0; i < b.N; i++ { // b.N is determined automatically for stable results
        myFunction(input)
    }
}

// Run benchmarks:
// go test -bench=. -benchtime=5s ./...
// go test -bench=BenchmarkMyFunction -benchmem ./...
```

### Reading Benchmark Output

```
BenchmarkMyFunction-8     1000000     1234 ns/op     256 B/op     4 allocs/op
                ↑                ↑          ↑             ↑              ↑
           GOMAXPROCS     iterations   nanoseconds    bytes            heap
                                        per op       allocated      allocations
                                                      per op          per op
```

### Benchmark Tips

```go
// Prevent compiler from optimizing away results (dead code elimination)
var result int
func BenchmarkMyFunction(b *testing.B) {
    var r int
    for i := 0; i < b.N; i++ {
        r = myFunction(input)
    }
    result = r // use result to prevent elimination
}

// Sub-benchmarks to compare approaches
func BenchmarkAlgorithms(b *testing.B) {
    b.Run("naive", func(b *testing.B) {
        for i := 0; i < b.N; i++ { naiveAlgo(input) }
    })
    b.Run("optimized", func(b *testing.B) {
        for i := 0; i < b.N; i++ { optimizedAlgo(input) }
    })
}

// Parallel benchmark
func BenchmarkConcurrent(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            myFunction(input)
        }
    })
}
```

---

## 3. CPU Profiling

CPU profiling tells you where your code spends CPU time. The Go runtime periodically (every 10ms by default) stops all goroutines and records the current stack trace.

### Enable CPU Profiling

```go
// Option 1: In a test or benchmark
go test -bench=. -cpuprofile=cpu.prof ./...

// Option 2: In a long-running service (HTTP handler)
import _ "net/http/pprof"
// Then curl http://localhost:6060/debug/pprof/profile?seconds=30 > cpu.prof

// Option 3: Programmatic
import "runtime/pprof"
f, _ := os.Create("cpu.prof")
pprof.StartCPUProfile(f)
defer pprof.StopCPUProfile()
// ... do work ...
```

### Analyzing CPU Profiles

```bash
go tool pprof cpu.prof

# Inside pprof interactive mode:
(pprof) top         # top functions by CPU time
(pprof) top -cum    # top functions by cumulative time (includes children)
(pprof) list myFunc # annotated source with time per line
(pprof) web         # open flame graph in browser (requires graphviz)

# Flame graph directly:
go tool pprof -http=:8080 cpu.prof
# Opens browser with interactive flame graph
```

### Reading Flame Graphs

```
A flame graph shows the call stack on the Y axis (root at bottom)
and CPU time on the X axis (wider = more time).

Wide flat boxes at the top = hot spots (spend most time here)

      ┌──────────────────────────────────────────────────────┐
      │                   json.Marshal (30%)                 │  ← hot!
      ├─────────────────────────┬────────────────────────────┤
      │  json.reflectValue (15%)│    json.encodeString (15%) │
      ├─────────────────────────┴────────────────────────────┤
      │              json.Marshal (50%)                      │
      ├──────────────────────────────────────────────────────┤
      │              handleRequest (80%)                     │
      ├──────────────────────────────────────────────────────┤
      │              http.(*ServeMux).ServeHTTP (100%)       │
      └──────────────────────────────────────────────────────┘
```

---

## 4. Memory Profiling

Memory profiling shows where heap allocations are happening.

```bash
# Generate a memory profile during a benchmark
go test -bench=. -memprofile=mem.prof ./...

# From a running service:
curl http://localhost:6060/debug/pprof/heap > mem.prof

# Analyze
go tool pprof mem.prof
(pprof) top          # functions allocating most memory
(pprof) list myFunc  # line-by-line allocation view
```

### Key pprof Flags for Memory

```bash
# Total bytes allocated (default)
go tool pprof -alloc_space mem.prof

# Number of objects allocated
go tool pprof -alloc_objects mem.prof

# Currently live bytes (what's in heap now)
go tool pprof -inuse_space mem.prof
```

---

## 5. Tracing

The execution tracer provides a detailed timeline of events: goroutine scheduling, GC, system calls, and more.

```bash
# Generate a trace
go test -trace=trace.out ./...

# From a service:
curl http://localhost:6060/debug/pprof/trace?seconds=5 > trace.out

# View the trace
go tool trace trace.out
```

The trace viewer shows:
- Goroutine timelines (when each goroutine was running, blocked, waiting)
- GC events (when GC ran, how long it took)
- Network and syscall events
- Heap growth

---

## 6. Practical Optimizations

### Optimization 1: Avoid Allocations in Hot Paths

```go
// Benchmark shows: 4 allocs/op — let's find and fix them

// BEFORE:
func formatResponse(data []Item) string {
    parts := []string{} // allocation 1: slice header
    for _, item := range data {
        parts = append(parts, fmt.Sprintf("%s: %d", item.Name, item.Value)) // alloc per item!
    }
    return strings.Join(parts, ", ") // allocation: join result
}

// AFTER: use strings.Builder, avoid fmt.Sprintf
func formatResponseFast(data []Item) string {
    var sb strings.Builder
    sb.Grow(len(data) * 20) // pre-allocate estimated size
    for i, item := range data {
        if i > 0 { sb.WriteString(", ") }
        sb.WriteString(item.Name)
        sb.WriteString(": ")
        sb.WriteString(strconv.Itoa(item.Value)) // no allocation!
    }
    return sb.String()
}
```

### Optimization 2: Inlining

The Go compiler inlines small functions to eliminate call overhead. A function is inlined if it's simple enough (default budget: 80 nodes).

```go
// Check inlining decisions:
go build -gcflags="-m" ./...
// Output: "can inline myFunc" or "myFunc too complex to inline"

// Force non-inlineable:
//go:noinline
func mustNotInline() {}
```

### Optimization 3: String vs []byte

```go
// Converting string ↔ []byte allocates! Avoid in hot paths.

// BAD:
func processRequest(body []byte) error {
    s := string(body)     // allocation!
    if strings.HasPrefix(s, "{") { ... }
    return nil
}

// GOOD: work with []byte directly
func processRequestFast(body []byte) error {
    if len(body) > 0 && body[0] == '{' { ... }
    return nil
}

// For JSON: use []byte directly with encoding/json
json.Unmarshal(body, &result) // takes []byte, no conversion needed
```

### Optimization 4: Struct Field Ordering for Cache Locality

```go
// Prefer fields accessed together to be adjacent in memory.
// The CPU fetches 64-byte cache lines. If hot fields are scattered,
// you need more cache lines → more cache misses → slower.

type HotPath struct {
    Count  int64  // hot: accessed on every request
    Total  int64  // hot: accessed on every request
    // ------ cache line boundary ------
    Name   string // cold: only accessed occasionally
    Config Config // cold: only accessed occasionally
}
```

### Optimization 5: Avoid Reflection

```go
// Reflection (reflect package) is much slower than direct code.
// json.Marshal uses reflection internally.

// For maximum performance: use code generation
// Tools: easyjson, ffjson generate marshal/unmarshal without reflection

//go:generate easyjson -all models.go

// Generated code is 3-10x faster than reflection-based encoding/json
```

---

## 7. Interview Questions & Model Answers

**Q: How do you diagnose a slow Go service?**

"I start with metrics to understand what kind of slowness it is: high latency, high CPU, or high memory. If it's CPU-bound: take a CPU profile with pprof and look at the flame graph — where are we spending most time? Usually it's either a hot loop or excessive memory allocation causing GC pressure. If it's memory: take a heap profile and look at top allocating functions. If it's concurrency-related: look at goroutine count, check for lock contention. For request latency issues: use distributed tracing to find which service or DB query is slow."

**Q: What is the overhead of running the race detector?**

"The race detector increases memory usage by ~5-10x and slows execution by ~2-20x. It instruments every memory access at compile time and uses shadow memory to track access timing and goroutines. It's too expensive for production but should run in CI on every test. We run `go test -race ./...` on every pull request."

**Q: Describe a real performance optimization you made in Go.**

"[Tailor to your experience, or use this template] In our API service, the heap profile showed that JSON marshaling was our #1 allocator — 60% of all allocations. The hot path was marshaling 100-200 response objects per request. We switched from `encoding/json` to `easyjson` (code-generated marshalers), which reduced allocations by 80% and improved P99 latency by 40%. The generated code avoids reflection and uses pre-allocated byte slices."

---

## Summary

- Measure before optimizing. Profile, don't guess.
- **Benchmarks:** `func BenchmarkX(b *testing.B)`. Use `-benchmem` to see allocations.
- **CPU profiling:** `go test -cpuprofile` or `pprof/profile` endpoint. Flame graphs show hot spots.
- **Memory profiling:** `go test -memprofile` or `pprof/heap`. Identifies top allocating code.
- **Tracing:** detailed goroutine and GC timeline. Useful for latency spikes.
- **Common wins:** eliminate allocations in hot paths (strings.Builder, pre-allocated slices, sync.Pool), fix struct field ordering, avoid reflection, check for lock contention.
- Always run the race detector in CI: `go test -race ./...`.
