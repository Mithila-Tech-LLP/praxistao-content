# Chapter 47: VaultDB — Performance Benchmarking and Tuning

You've built a working database. Now: how fast is it? Where is time being spent? What can we improve? This chapter covers profiling, benchmarking, and the key optimizations that make databases fast.

## Table of Contents

1. Measuring Performance
2. CPU Profiling
3. Memory Profiling
4. Key Performance Bottlenecks
5. Optimization: Write Batching
6. Optimization: Buffer Pool Tuning
7. Optimization: Sequential I/O for Scans
8. Benchmarks and Target Numbers
9. Exercises

---

## 1. Measuring Performance

Go has excellent built-in benchmarking. A benchmark function starts with `Benchmark`:

```go
// benchmarks_test.go
package vaultdb_test

import (
    "fmt"
    "os"
    "testing"
    "github.com/yourname/vaultdb/query"
    "github.com/yourname/vaultdb/storage"
    "github.com/yourname/vaultdb/wal"
)

func setupDB(b *testing.B) (*storage.DiskManager, *storage.BufferPool, *wal.WAL, *storage.Catalog) {
    b.Helper()
    dm, _ := storage.NewDiskManager(b.TempDir() + "/bench.vault")
    bp := storage.NewBufferPool(dm, 1024)
    w, _ := wal.Open(b.TempDir() + "/bench.wal")
    catalog := &storage.Catalog{FreePage: storage.InvalidPageID}
    return dm, bp, w, catalog
}

func BenchmarkInsert(b *testing.B) {
    dm, bp, w, catalog := setupDB(b)
    exec := query.NewExecutor(dm, bp, w, catalog)
    exec.Execute(query.MustParse("CREATE TABLE bench (id INT, val VARCHAR)"))

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        sql := fmt.Sprintf("INSERT INTO bench (id, val) VALUES (%d, 'value%d')", i, i)
        exec.Execute(query.MustParse(sql))
    }
}

func BenchmarkSelectAll(b *testing.B) {
    dm, bp, w, catalog := setupDB(b)
    exec := query.NewExecutor(dm, bp, w, catalog)
    exec.Execute(query.MustParse("CREATE TABLE bench (id INT, val VARCHAR)"))

    // Pre-populate 10,000 rows
    for i := 0; i < 10000; i++ {
        exec.Execute(query.MustParse(fmt.Sprintf(
            "INSERT INTO bench VALUES (%d, 'value%d')", i, i)))
    }

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        exec.Execute(query.MustParse("SELECT * FROM bench"))
    }
}
```

Run with:
```bash
go test -bench=. -benchmem -benchtime=5s ./...
```

Output example:
```
BenchmarkInsert-8         50000    32451 ns/op    2048 B/op    23 allocs/op
BenchmarkSelectAll-8        100  12453210 ns/op    512KB B/op  1234 allocs/op
```

- `ns/op`: nanoseconds per operation. Lower is better.
- `B/op`: bytes allocated per operation. Lower is better.
- `allocs/op`: heap allocations per operation. Lower is better.

---

## 2. CPU Profiling

Find where VaultDB spends its CPU time:

```go
// cmd/profile/main.go
package main

import (
    "fmt"
    "os"
    "runtime/pprof"
)

func main() {
    // Start CPU profiling
    f, _ := os.Create("cpu.prof")
    pprof.StartCPUProfile(f)
    defer pprof.StopCPUProfile()

    // Run workload
    runInsertBenchmark(100000)
    runSelectBenchmark(1000)
    
    fmt.Println("Profile written to cpu.prof")
}

// Then analyze:
// go tool pprof cpu.prof
// (pprof) top20
// (pprof) web  (opens flame graph in browser)
```

**What to look for:**

```
(pprof) top10
      flat  flat%   sum%        cum   cum%
    450ms  45.0%  45.0%      490ms  49.0%  os.(*File).Write
    200ms  20.0%  65.0%      200ms  20.0%  runtime.memmove
    150ms  15.0%  80.0%      200ms  20.0%  encoding/binary.BigEndian.PutUint64
    100ms  10.0%  90.0%      100ms  10.0%  sync.(*Mutex).Lock
```

If `os.Write` dominates → you're calling write too often. Batch writes!
If `sync.Mutex.Lock` dominates → too much lock contention. Reduce critical section size.
If `memmove` dominates → too many slice copies. Pre-allocate slices.

---

## 3. Memory Profiling

Find memory leaks and excessive allocations:

```go
// At end of benchmark, write heap profile
func writeMemProfile() {
    f, _ := os.Create("mem.prof")
    defer f.Close()
    runtime.GC() // force GC before profiling
    pprof.WriteHeapProfile(f)
}

// Analyze:
// go tool pprof -alloc_space mem.prof
// (pprof) top20
```

**Common memory issues in databases:**
- Allocating a new `[]byte` for every row encoding → use `sync.Pool` for buffers
- Keeping all rows in memory during a scan → use iterators
- Not releasing pinned buffer pool frames → memory grows until eviction fails

---

## 4. Key Performance Bottlenecks

### Bottleneck 1: fsync per transaction

Every `COMMIT` calls `file.Sync()` which flushes the OS write buffer to physical disk. This takes 1-10ms.

**Solution: Group commit**

Instead of syncing per transaction, batch multiple commits and sync once:

```go
type WAL struct {
    // ...
    commitBatch chan commitRequest
}

type commitRequest struct {
    txnID uint64
    done  chan error
}

func (w *WAL) startGroupCommitter() {
    go func() {
        var pending []commitRequest
        ticker := time.NewTicker(1 * time.Millisecond)

        for {
            select {
            case req := <-w.commitBatch:
                pending = append(pending, req)
            case <-ticker.C:
                if len(pending) == 0 {
                    continue
                }
                // One fsync for all pending commits
                err := w.file.Sync()
                for _, req := range pending {
                    req.done <- err
                }
                pending = pending[:0]
            }
        }
    }()
}

func (w *WAL) CommitSync(txnID uint64) error {
    done := make(chan error, 1)
    w.commitBatch <- commitRequest{txnID: txnID, done: done}
    return <-done
}
```

Result: 1000 transactions/second → 50,000+ transactions/second.

### Bottleneck 2: Lock contention on the buffer pool

Every page access acquires the buffer pool's global mutex. With many goroutines, this serializes all I/O.

**Solution: Partition the buffer pool** into N sub-pools, each with its own mutex. Page i goes to sub-pool i % N.

```go
type ShardedBufferPool struct {
    shards [16]*BufferPool
}

func (bp *ShardedBufferPool) shard(pageID PageID) *BufferPool {
    return bp.shards[pageID%16]
}

func (bp *ShardedBufferPool) FetchPage(pageID PageID) (*Page, error) {
    return bp.shard(pageID).FetchPage(pageID)
}
```

### Bottleneck 3: Allocations in the hot path

Row encoding allocates a new `[]byte` on every insert. GC pressure increases with insert rate.

**Solution: Buffer pool with sync.Pool**

```go
var encodePool = sync.Pool{
    New: func() interface{} {
        buf := make([]byte, 0, 256)
        return &buf
    },
}

func EncodeRowPooled(row Row) []byte {
    bufPtr := encodePool.Get().(*[]byte)
    buf := (*bufPtr)[:0]
    // ... encode into buf ...
    result := make([]byte, len(buf))
    copy(result, buf)
    encodePool.Put(bufPtr)
    return result
}
```

---

## 5. Optimization: Write Batching

Instead of writing one page per insert:

```go
type WriteBatcher struct {
    dm      *DiskManager
    pending map[PageID]*Page
    mu      sync.Mutex
}

func (b *WriteBatcher) SetDirty(id PageID, page *Page) {
    b.mu.Lock()
    b.pending[id] = page
    b.mu.Unlock()
}

func (b *WriteBatcher) Flush() error {
    b.mu.Lock()
    pages := b.pending
    b.pending = make(map[PageID]*Page)
    b.mu.Unlock()

    // Write all dirty pages in sorted order (minimizes disk seeks)
    ids := sortedPageIDs(pages)
    for _, id := range ids {
        if err := b.dm.WritePage(id, pages[id]); err != nil {
            return err
        }
    }
    return nil
}
```

---

## 6. Benchmarks and Target Numbers

After optimizations, VaultDB should achieve:

| Operation | Target | Notes |
|-----------|--------|-------|
| Single INSERT (no fsync) | < 5,000 ns | Pure memory write |
| Single INSERT (with fsync) | < 2 ms | Disk I/O per commit |
| Batch INSERT (group commit) | > 50,000/s | Amortized fsync |
| SELECT (1K rows, no WHERE) | < 5 ms | Sequential scan |
| SELECT (1K rows, B-Tree index) | < 100 µs | Index lookup |
| Full scan (100K rows) | < 100 ms | ~250 pages × 0.4 ms |

Compare to PostgreSQL on same hardware:
- INSERT: 5,000-10,000/sec (with fsync, single row)
- SELECT by primary key: < 1ms
- Table scan 100K rows: < 50ms

VaultDB won't match PostgreSQL (which has 30 years of optimizations) but it teaches the same principles.

---

## Summary

- Use `go test -bench=. -benchmem` to measure performance.
- CPU profiling shows where time is spent. Memory profiling shows allocation hotspots.
- The three biggest bottlenecks: fsync per commit (solution: group commit), lock contention (solution: sharding), allocations (solution: sync.Pool).
- Sequential writes are 10x faster than random writes on spinning disks. Sort page writes by page ID.
- Always measure before optimizing — intuition about what's slow is often wrong.

### Exercises

**Easy:** Run `BenchmarkInsert` with `-benchmem`. Note the allocs/op. Add `sync.Pool` for the row encoding buffer and compare.

**Medium:** Implement group commit and benchmark the difference in commit throughput: solo commit (one fsync per commit) vs group commit (one fsync per batch). Aim for > 10x improvement.

**Hard:** Implement a benchmark that measures the full read path (INSERT 1M rows, then SELECT with B-Tree index). Profile the result and identify the top 3 hotspots by CPU time. Fix at least one.
