# Chapter 24: Profiling and Performance

Fast code starts with measurement. Go has world-class built-in profiling tools — `pprof`, built-in benchmarks, and trace — that let you find exactly where time and memory are going. This chapter teaches you how to profile, what to look for, and how to optimize based on evidence rather than guesswork.

## Table of Contents

1. [Benchmarking — Measuring Before Optimizing](#1-benchmarking--measuring-before-optimizing)
2. [CPU Profiling with pprof](#2-cpu-profiling-with-pprof)
3. [Memory Profiling](#3-memory-profiling)
4. [Goroutine and Block Profiling](#4-goroutine-and-block-profiling)
5. [The go tool trace](#5-the-go-tool-trace)
6. [HTTP Server Profiling](#6-http-server-profiling)
7. [Optimization Techniques Recap](#7-optimization-techniques-recap)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Benchmarking — Measuring Before Optimizing

**Writing a useful benchmark:**
```go
func BenchmarkProcessRequest(b *testing.B) {
    // Setup outside the loop — not counted in timing:
    data := generateTestData(1000)
    
    b.ResetTimer()  // Reset after setup
    b.ReportAllocs()  // Show memory allocation stats
    
    for i := 0; i < b.N; i++ {
        _ = ProcessRequest(data)  // _ prevents dead-code elimination
    }
}
```

**Running benchmarks:**
```bash
go test -bench=BenchmarkProcessRequest -benchmem -count=5 -benchtime=3s
```

**Reading benchmark output:**
```
BenchmarkProcessRequest-8    500000    2543 ns/op    1024 B/op    12 allocs/op
```
- `-8`: GOMAXPROCS=8 (running on 8 cores)
- `500000`: how many times the function ran (`b.N`)
- `2543 ns/op`: nanoseconds per call
- `1024 B/op`: bytes allocated per call
- `12 allocs/op`: heap allocations per call

**Benchstat — comparing before and after:**
```bash
go install golang.org/x/perf/cmd/benchstat@latest

# Run old version:
go test -bench=. -count=10 | tee old.txt

# After your optimization:
go test -bench=. -count=10 | tee new.txt

benchstat old.txt new.txt
```
```
name                   old time/op   new time/op   delta
ProcessRequest-8       2.54µs ± 2%  1.12µs ± 1%  -55.90%  (p=0.000 n=10+10)
name                   old alloc/op  new alloc/op  delta
ProcessRequest-8       1.02kB ± 0%  0.00kB ± 0%  -100.00% (p=0.000 n=10+10)
```

**Benchmarking sub-operations:**
```go
func BenchmarkSearchAlgorithms(b *testing.B) {
    data := sortedSlice(10000)
    target := data[5000]
    
    b.Run("linear", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            linearSearch(data, target)
        }
    })
    
    b.Run("binary", func(b *testing.B) {
        for i := 0; i < b.N; i++ {
            binarySearch(data, target)
        }
    })
}
```

### Quick Check
> 1. Why is `b.ResetTimer()` called after setup code?
> 2. What does `-benchmem` add to the output?
> 3. What does `_ = result` in a benchmark prevent?

---

## 2. CPU Profiling with pprof

**Step 1: Generate a CPU profile:**
```go
import (
    "os"
    "runtime/pprof"
)

func main() {
    // Start CPU profiling:
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()
    
    // Your code here:
    runWorkload()
}
```

**Or directly from benchmarks (easiest):**
```bash
go test -bench=BenchmarkProcessRequest -cpuprofile=cpu.prof
```

**Step 2: Analyze the profile:**
```bash
go tool pprof cpu.prof
```

**Common pprof commands:**
```
(pprof) top10           # Top 10 functions by CPU time
(pprof) top10 -cum      # Top 10 by cumulative (includes callees)
(pprof) list FuncName   # Show annotated source for FuncName
(pprof) web             # Open flame graph in browser (requires graphviz)
(pprof) svg > out.svg   # Save call graph as SVG
```

**Sample top output:**
```
(pprof) top10
Showing nodes accounting for 4.52s, 89.70% of 5.04s total
      flat  flat%   sum%        cum   cum%
     1.85s 36.71% 36.71%      1.85s 36.71%  runtime.memmove
     0.93s 18.45% 55.16%      0.93s 18.45%  encoding/json.(*encodeState).string
     0.72s 14.29% 69.44%      2.54s 50.40%  encoding/json.marshalJSON
     0.41s  8.13% 77.58%      0.41s  8.13%  sync.(*RWMutex).Lock
```

- **flat**: time spent IN this function (not its callees)
- **cum**: cumulative time including callees
- **flat%**: percentage of total CPU time

**Reading the flame graph:**
```
A flame graph shows the call stack horizontally.
Wide bars = more CPU time in that function/path.
Tall stacks = deep call trees.
Look for wide bars near the BOTTOM — those are the hot spots.

                  ┌──────────┐
              ┌───┤ json.Marshal ├────────┐
         ┌────┤               ─────────────────┐
    ┌────┤ handleRequest ─────────────────────────────┐
────┴─────────────────────────────────────────────────────
```

**Benchmark with CPU profile in one command:**
```bash
# Profile benchmark, open in browser:
go test -bench=. -cpuprofile=cpu.prof && go tool pprof -http=:8080 cpu.prof
```

### Quick Check
> 1. What is the difference between `flat` and `cum` in pprof output?
> 2. What does a wide bar in a flame graph mean?
> 3. What flag generates a CPU profile from a benchmark?

---

## 3. Memory Profiling

**Heap profile — where are allocations happening?**
```go
import (
    "os"
    "runtime"
    "runtime/pprof"
)

func main() {
    runWorkload()
    
    // Write heap profile at end of program:
    runtime.GC()  // Run GC first to get accurate live object stats
    f, _ := os.Create("mem.prof")
    defer f.Close()
    pprof.WriteHeapProfile(f)
}
```

**From benchmarks:**
```bash
go test -bench=. -memprofile=mem.prof
```

**Analyzing memory profiles:**
```bash
go tool pprof -alloc_space mem.prof    # Where was memory allocated?
go tool pprof -inuse_space mem.prof    # What's currently live?
```

**Key pprof memory commands:**
```
(pprof) top10                   # Top allocating functions
(pprof) top10 -cum              # Include callees
(pprof) list functionName       # Annotated source with allocation lines
(pprof) web                     # Visual flame graph
```

**Allocation profile example:**
```
(pprof) top5 -alloc_space
Showing nodes accounting for 1.2GB, 87% of 1.38GB total
      flat  flat%   sum%        cum   cum%
   512.00MB 37.14% 37.14%   512.00MB 37.14%  bytes.makeSlice
   256.00MB 18.57% 55.71%   768.00MB 55.71%  encoding/json.Marshal
   192.00MB 13.91% 69.62%   192.00MB 13.91%  strings.(*Builder).grow
```

**Finding specific allocation lines:**
```
(pprof) list ProcessRequest
Total: 1.38GB
ROUTINE ======================== main.ProcessRequest
     0     768MB (flat, cum) 55.71% of Total
         .          .     12:func ProcessRequest(data []byte) string {
     0     256MB     13:    result, _ := json.Marshal(parseData(data))  // ← HERE
     0     512MB     14:    return string(result)  // ← and HERE (string copy)
         .          .     15:}
```

### Quick Check
> 1. What is the difference between `-alloc_space` and `-inuse_space`?
> 2. Why run `runtime.GC()` before writing a heap profile?
> 3. What does the `list` command show in pprof?

---

## 4. Goroutine and Block Profiling

**Goroutine profile — stuck goroutines:**
```bash
go test -bench=. -cpuprofile=cpu.prof

# Or from HTTP endpoint (see Section 6):
curl http://localhost:6060/debug/pprof/goroutine?debug=1
```

```go
// In code:
import "runtime/pprof"

f, _ := os.Create("goroutine.prof")
pprof.Lookup("goroutine").WriteTo(f, 0)
```

**Block profile — where goroutines block:**
```go
import "runtime"

// Enable block profiling (not enabled by default — has overhead):
runtime.SetBlockProfileRate(1)  // Sample every blocking event

// ... run workload ...

f, _ := os.Create("block.prof")
pprof.Lookup("block").WriteTo(f, 0)
```

**Mutex profile — lock contention:**
```go
runtime.SetMutexProfileFraction(1)  // Sample every mutex event

// ... run workload ...

f, _ := os.Create("mutex.prof")
pprof.Lookup("mutex").WriteTo(f, 0)
```

**Analyzing goroutine dumps:**
```
goroutine 1 [running]:
main.main()
	/tmp/main.go:24 +0x...

goroutine 18 [chan receive, 1 minutes]:  ← STUCK for 1 minute!
main.worker(0x1400012e0c0)
	/tmp/main.go:15 +0x...
```

### Quick Check
> 1. What does a block profile show?
> 2. Why is block profiling not enabled by default?
> 3. What does a goroutine marked `[chan receive, 5 minutes]` mean?

---

## 5. The go tool trace

`trace` captures a detailed timeline of everything happening in your program — goroutine scheduling, GC, syscalls:

```go
import "runtime/trace"

func main() {
    f, _ := os.Create("trace.out")
    trace.Start(f)
    defer trace.Stop()
    
    runWorkload()
}
```

**From tests:**
```bash
go test -bench=. -trace=trace.out
go tool trace trace.out  # Opens browser UI
```

**What trace shows:**
```
Timeline view:
  Goroutine 1: [running][blocked on chan][running][GC assist]
  Goroutine 2:          [blocked]        [running]
  GC:                                    [mark][sweep]
  
Useful for:
  - Why is there a GC pause affecting all goroutines?
  - Why are goroutines blocking on each other?
  - Is GOMAXPROCS being used effectively?
  - Where are the scheduling delays?
```

**Custom trace annotations:**
```go
import "runtime/trace"

ctx, task := trace.NewTask(context.Background(), "processOrder")
defer task.End()

trace.WithRegion(ctx, "fetchUser", func() {
    user = fetchUser(userID)
})

trace.WithRegion(ctx, "validateOrder", func() {
    validate(order, user)
})
```

### Quick Check
> 1. What does `go tool trace` show that pprof doesn't?
> 2. How do you add custom labels to a trace?

---

## 6. HTTP Server Profiling

Add `net/http/pprof` to any HTTP server for live profiling:

```go
import (
    _ "net/http/pprof"  // Side-effect import: registers /debug/pprof/* routes
    "net/http"
)

func main() {
    // Start pprof server on a different port (don't expose on production port!):
    go func() {
        log.Println(http.ListenAndServe("localhost:6060", nil))
    }()
    
    // Your main server:
    http.ListenAndServe(":8080", yourRouter)
}
```

**Available endpoints:**
```
http://localhost:6060/debug/pprof/           # Index
http://localhost:6060/debug/pprof/goroutine  # Goroutine dump
http://localhost:6060/debug/pprof/heap       # Heap profile
http://localhost:6060/debug/pprof/profile    # 30s CPU profile
http://localhost:6060/debug/pprof/trace      # 5s execution trace
http://localhost:6060/debug/pprof/block      # Block profile
http://localhost:6060/debug/pprof/mutex      # Mutex contention
```

**Profile a running server:**
```bash
# 30-second CPU profile of running server:
go tool pprof http://localhost:6060/debug/pprof/profile

# Heap profile:
go tool pprof http://localhost:6060/debug/pprof/heap

# Open directly in browser with flame graph:
go tool pprof -http=:8081 http://localhost:6060/debug/pprof/profile
```

**Security warning:** Never expose pprof endpoints publicly — they reveal internal code structure and can be used to DoS the server.

### Quick Check
> 1. What import enables the pprof HTTP endpoints?
> 2. Why should pprof run on a different port from your main server?
> 3. How do you profile a running production server?

---

## 7. Optimization Techniques Recap

**Profiling-guided optimization checklist:**
```
1. Profile first — don't guess
2. Optimize the hot path (top 3 functions in pprof)
3. Reduce allocations (allocs/op → 0 where possible)
4. Reduce copying (use pointers for large structs)
5. Improve cache locality (use arrays over pointer chains)
6. Reduce lock contention (use RWMutex, sharding, or lock-free)
7. Parallelize CPU-bound work (goroutines + WaitGroup)
8. Use sync.Pool for hot-path temporary objects
9. Pre-allocate slices/maps when size is known
10. Profile again — verify the improvement
```

**Quick wins reference:**
```go
// 1. Append with known size:
s := make([]T, 0, n)  // Pre-allocate

// 2. String build:
var sb strings.Builder
sb.Grow(estimatedSize)

// 3. Buffer reuse:
var pool = sync.Pool{New: func() any { return new(bytes.Buffer) }}

// 4. Avoid repeated map lookups:
if v, ok := m[k]; ok {
    use(v)
}

// 5. Integer division vs multiplication:
x := n / 2    // Slower
x := n >> 1   // Faster (for power-of-2 divisors)

// 6. Use slices over linked lists:
type SliceList []Item  // Better cache locality than *Node chains

// 7. Interface satisfaction check at compile time:
var _ io.Writer = (*MyWriter)(nil)  // Fail fast if interface breaks

// 8. Avoid defer in tight loops:
// BAD:
for _, f := range files {
    f.Open()
    defer f.Close()  // Accumulates, doesn't run til function end
}
// GOOD: call Close directly or use helper function
```

---

## Summary

- **Benchmark first**: `go test -bench=. -benchmem -count=10` — measure before optimizing
- **benchstat**: compare before/after with statistical significance
- **CPU profile**: `go test -cpuprofile=cpu.prof` + `go tool pprof -http=:8080 cpu.prof`
- **flat vs cum**: flat = time in function; cum = including callees
- **Flame graph**: wide bar = hot function; read from bottom-up
- **Memory profile**: `-memprofile=mem.prof`; `-alloc_space` = total allocs; `-inuse_space` = live
- **Block/mutex profile**: needs explicit enabling with `SetBlockProfileRate` / `SetMutexProfileFraction`
- **trace**: goroutine timeline, GC events, scheduling — use for latency investigations
- **HTTP pprof**: `import _ "net/http/pprof"` + separate port — profile live servers
- **Rule**: profile → find hotspot → optimize → benchmark → verify → repeat

---

## Exercises

### Easy
1. Write a benchmark for three ways to concatenate 100 strings: (a) `+=` in loop, (b) `strings.Join`, (c) `strings.Builder`. Show allocations with `-benchmem`. What is the allocation count for each?
2. Add the pprof HTTP handler to a simple HTTP server. Write a load generator that sends 1000 requests/second for 30 seconds. Capture a CPU profile and find the top 3 functions by flat CPU time.
3. Write two versions of sorting 1M integers: (a) `sort.Ints`, (b) `sort.Slice`. Benchmark both. Which is faster? Profile to see why.

### Medium
4. Memory leak detector: Write a server that has a deliberate memory leak (e.g., a cache that grows forever). Use pprof heap profiling (`-inuse_space`) to identify the leaking function. Fix the leak (add eviction). Show before/after memory profiles confirming the fix.
5. Lock contention analysis: Build a concurrent counter backed by a single mutex. Load test it with 100 goroutines incrementing 10,000 times each. Profile with mutex profiling enabled. Then rewrite using 16 sharded counters (each with its own mutex). Compare: total time, mutex contention (from profile), and correctness.
6. Allocation-free JSON path: Profile a JSON encoding hot path. Current code: `json.Marshal(struct)` per request. Using the profile, identify the allocation sources. Rewrite using a pre-allocated buffer from `sync.Pool` + manual JSON construction with `encoding/json.Encoder` writing to the pooled buffer. Verify allocs/op drops from 10+ to 1 or 0.

### Hard
7. Complete profiling workflow: Take a provided slow HTTP handler (processes a CSV, does aggregation, returns JSON). Use ALL profiling tools: (a) CPU profile — find the computational bottleneck, (b) Memory profile — find the allocation hotspot, (c) Block profile — find the synchronization bottleneck, (d) Trace — find scheduling inefficiencies. Fix each issue. Document before/after for each metric. Target: 10× throughput improvement.
8. Continuous profiling integration: Build a `ProfileCollector` that automatically captures 60-second CPU and heap profiles every 10 minutes and stores them (in-memory ring buffer, last 24 hours). Expose a `/profiles` HTTP endpoint listing available profiles and `/profiles/{id}` to download a specific profile. Add automatic alerting: if any function increases its CPU share by >20% between consecutive profiles, log a warning. Test by artificially slowing a function and verifying the alert fires.
