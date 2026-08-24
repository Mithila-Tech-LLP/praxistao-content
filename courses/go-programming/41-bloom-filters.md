# Chapter 41: Bloom Filters and Probabilistic Data Structures

Some problems don't require exact answers — they need fast, memory-efficient approximations. A Bloom filter can tell you "definitely not in the set" or "probably in the set" using kilobytes of memory where a hash set would use megabytes.

## Table of Contents

1. [Bloom Filter](#1-bloom-filter)
2. [HyperLogLog — Counting Uniques](#2-hyperloglog--counting-uniques)
3. [Count-Min Sketch — Frequency Estimation](#3-count-min-sketch--frequency-estimation)
4. [When to Use Probabilistic DS](#4-when-to-use-probabilistic-ds)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. Bloom Filter

### Concept

A Bloom filter is a bit array of size `m` and `k` hash functions. To **insert** an element, hash it `k` times and set those bit positions to 1. To **query**, hash it `k` times and check all bit positions — if any are 0, the element is **definitely not** in the set. If all are 1, it's **probably** in the set (false positive possible).

**Properties:**
- No false negatives — if an element was inserted, the query always returns true
- False positive rate is controllable — depends on m (bit array size) and k (hash functions)
- Cannot delete elements (use Counting Bloom Filter for that)
- O(k) time for insert and query — independent of number of elements

**Optimal parameters:**
```
Given n expected elements and desired FPR p:
  m = -n * ln(p) / (ln(2))²      (bits)
  k = (m/n) * ln(2)              (hash functions)
```

For n=1,000,000 and p=0.01 (1% false positive rate):
- m ≈ 9,585,058 bits ≈ 1.14 MB
- k ≈ 7 hash functions

Compare to a hash set: ~40 MB for 1M strings.

```go
package bloom

import (
    "encoding/binary"
    "math"

    "github.com/spaolacci/murmur3"
)

type Filter struct {
    bits []uint64
    m    uint64 // bit array size
    k    uint   // number of hash functions
}

// New creates a Bloom filter for n expected elements with false positive rate p.
func New(n uint64, p float64) *Filter {
    m := optimalM(n, p)
    k := optimalK(m, n)
    return &Filter{
        bits: make([]uint64, (m+63)/64),
        m:    m,
        k:    k,
    }
}

func optimalM(n uint64, p float64) uint64 {
    return uint64(math.Ceil(-float64(n) * math.Log(p) / (math.Log(2) * math.Log(2))))
}

func optimalK(m, n uint64) uint {
    return uint(math.Round(float64(m) / float64(n) * math.Log(2)))
}

// hashes returns k hash positions for key.
func (f *Filter) hashes(key []byte) []uint64 {
    // Use double hashing: h_i(key) = h1(key) + i * h2(key)
    // This simulates k independent hash functions from two.
    h1 := murmur3.Sum64(key)
    // Second hash: seed with a different value
    buf := make([]byte, len(key)+4)
    binary.LittleEndian.PutUint32(buf, 0xDEADBEEF)
    copy(buf[4:], key)
    h2 := murmur3.Sum64(buf)

    positions := make([]uint64, f.k)
    for i := uint(0); i < f.k; i++ {
        positions[i] = (h1 + uint64(i)*h2) % f.m
    }
    return positions
}

func (f *Filter) set(pos uint64) {
    f.bits[pos/64] |= 1 << (pos % 64)
}

func (f *Filter) get(pos uint64) bool {
    return f.bits[pos/64]&(1<<(pos%64)) != 0
}

// Add inserts an element into the filter.
func (f *Filter) Add(key []byte) {
    for _, pos := range f.hashes(key) {
        f.set(pos)
    }
}

// AddString is a convenience wrapper.
func (f *Filter) AddString(s string) {
    f.Add([]byte(s))
}

// Contains returns false if the element is definitely not in the set,
// or true if it is probably in the set.
func (f *Filter) Contains(key []byte) bool {
    for _, pos := range f.hashes(key) {
        if !f.get(pos) {
            return false
        }
    }
    return true
}

// ContainsString is a convenience wrapper.
func (f *Filter) ContainsString(s string) bool {
    return f.Contains([]byte(s))
}

// EstimatedFPR returns the current estimated false positive rate
// given that n elements have been inserted.
func (f *Filter) EstimatedFPR(n uint64) float64 {
    // FPR = (1 - e^(-k*n/m))^k
    exponent := -float64(f.k) * float64(n) / float64(f.m)
    return math.Pow(1-math.Exp(exponent), float64(f.k))
}

// MemoryBytes returns the size of the bit array in bytes.
func (f *Filter) MemoryBytes() int {
    return len(f.bits) * 8
}
```

### Usage

```go
package main

import (
    "fmt"
    "bloom"
)

func main() {
    // Create filter for 1M elements, 1% FPR
    f := bloom.New(1_000_000, 0.01)

    // Add elements
    words := []string{"apple", "banana", "cherry", "date"}
    for _, w := range words {
        f.AddString(w)
    }

    // Query
    fmt.Println(f.ContainsString("apple"))    // true (it was inserted — true is "probably in set")
    fmt.Println(f.ContainsString("grape"))    // false (false always means definitely NOT in set)
    fmt.Println(f.ContainsString("mango"))    // false (usually — a rare false positive could return true)

    fmt.Printf("Memory: %d bytes\n", f.MemoryBytes()) // ~1.14 MB
}
```

### Real-world applications

- **Database query optimization**: check if a key might exist before a disk read (Cassandra, RocksDB, HBase all use this)
- **Web crawlers**: avoid re-crawling visited URLs (Chrome uses one for safe browsing)
- **Spam filtering**: check if an email address has been seen before
- **CDN cache**: check if a content hash is in the cache before fetching from origin

---

## 2. HyperLogLog — Counting Uniques

### Problem

How many unique visitors does your website have? With 100M daily visitors, storing a hash set takes gigabytes. HyperLogLog estimates the cardinality (count of distinct elements) using just a few kilobytes of memory with ~1-2% error (our implementation below uses 16 KB; Redis squeezes the same 16,384 registers into 12 KB).

### Concept

HyperLogLog hashes each element and looks at the leading zeros in the binary representation of the hash. An element with `r` leading zeros is rare (probability 1/2^r). The maximum number of leading zeros seen across all elements is a statistical estimator of log₂(cardinality).

To reduce variance, HyperLogLog uses `m` buckets (sub-streams) and takes a harmonic mean.

```go
package hll

import (
    "math"
    "math/bits"

    "github.com/spaolacci/murmur3"
)

const (
    precision = 14               // b=14: 2^14 = 16384 buckets
    m         = 1 << precision   // number of registers
    alphaMM   = 0.7213 / (1 + 1.079/m) * m * m
)

type HyperLogLog struct {
    registers [m]uint8
}

func New() *HyperLogLog { return &HyperLogLog{} }

func (h *HyperLogLog) Add(data []byte) {
    hash := murmur3.Sum64(data)
    // Use top `precision` bits as bucket index
    idx := hash >> (64 - precision)
    // Count leading zeros in remaining bits + 1
    remaining := hash<<precision | (1<<precision - 1)
    rho := uint8(bits.LeadingZeros64(remaining)) + 1
    if rho > h.registers[idx] {
        h.registers[idx] = rho
    }
}

func (h *HyperLogLog) AddString(s string) {
    h.Add([]byte(s))
}

func (h *HyperLogLog) Count() uint64 {
    sum := 0.0
    zeros := 0
    for _, r := range h.registers {
        sum += math.Pow(2, -float64(r))
        if r == 0 {
            zeros++
        }
    }
    estimate := alphaMM / sum

    // Small range correction
    if estimate <= 2.5*m && zeros > 0 {
        estimate = m * math.Log(float64(m)/float64(zeros))
    }
    // Large range correction
    const two32 = 1 << 32
    if estimate > two32/30.0 {
        estimate = -two32 * math.Log(1-estimate/two32)
    }
    return uint64(estimate)
}

func (h *HyperLogLog) Merge(other *HyperLogLog) *HyperLogLog {
    result := &HyperLogLog{}
    for i := range h.registers {
        if h.registers[i] > other.registers[i] {
            result.registers[i] = h.registers[i]
        } else {
            result.registers[i] = other.registers[i]
        }
    }
    return result
}
```

### Usage

```go
hll := hll.New()

// Simulate 1M unique user IDs
for i := 0; i < 1_000_000; i++ {
    hll.AddString(fmt.Sprintf("user_%d", i))
}

fmt.Printf("Estimated: %d (actual: 1,000,000)\n", hll.Count())
// Estimated: 999,843 (error < 0.02%)
fmt.Printf("Memory: %d bytes\n", m)  // 16,384 bytes ≈ 16 KB
```

---

## 3. Count-Min Sketch — Frequency Estimation

### Problem

Which URLs are being requested most? You have 100M requests/day — you can't keep a full count table. Count-Min Sketch estimates the frequency of any element using O(width × depth) space.

### Concept

A 2D array of counters (width × depth). To insert, hash the element with `depth` different hash functions, one per row, and increment that row's cell. To query, hash the same way and return the minimum of the counters (the minimum is the best estimate because false positives only over-count, never under-count).

```go
package cms

import (
    "encoding/binary"
    "math"

    "github.com/spaolacci/murmur3"
)

type Sketch struct {
    table  [][]uint64
    width  uint
    depth  uint
    seeds  []uint32
}

// New creates a Count-Min Sketch with the given error bounds:
// epsilon = error tolerance, delta = failure probability
// (e.g. delta=0.01 means the error bound holds with 99% confidence)
func New(epsilon, delta float64) *Sketch {
    width := uint(math.Ceil(math.E / epsilon))
    depth := uint(math.Ceil(math.Log(1 / delta)))

    table := make([][]uint64, depth)
    seeds := make([]uint32, depth)
    for i := range table {
        table[i] = make([]uint64, width)
        seeds[i] = uint32(i * 0x9747b28c)
    }
    return &Sketch{table: table, width: width, depth: depth, seeds: seeds}
}

func (s *Sketch) hash(key []byte, seed uint32) uint {
    buf := make([]byte, len(key)+4)
    binary.LittleEndian.PutUint32(buf, seed)
    copy(buf[4:], key)
    return uint(murmur3.Sum64(buf)) % s.width
}

func (s *Sketch) Increment(key []byte) {
    for i := uint(0); i < s.depth; i++ {
        col := s.hash(key, s.seeds[i])
        s.table[i][col]++
    }
}

func (s *Sketch) IncrementString(key string) {
    s.Increment([]byte(key))
}

func (s *Sketch) Estimate(key []byte) uint64 {
    min := ^uint64(0)
    for i := uint(0); i < s.depth; i++ {
        col := s.hash(key, s.seeds[i])
        if s.table[i][col] < min {
            min = s.table[i][col]
        }
    }
    return min
}

func (s *Sketch) EstimateString(key string) uint64 {
    return s.Estimate([]byte(key))
}
```

### Usage: Top-K URLs

```go
sketch := cms.New(0.001, 0.01) // 0.1% error, 1% failure probability (99% confidence)

// Count requests
for _, url := range requestLog {
    sketch.IncrementString(url)
}

// Query frequencies
for _, url := range []string{"/", "/api/users", "/api/orders"} {
    fmt.Printf("%s: ~%d requests\n", url, sketch.EstimateString(url))
}
```

---

## 4. When to Use Probabilistic DS

| Problem | Structure | Trade-off |
|---------|-----------|-----------|
| "Have I seen this before?" | Bloom Filter | False positives only, no delete |
| "How many unique X?" | HyperLogLog | ~1-2% error, fixed 16 KB |
| "How often does X appear?" | Count-Min Sketch | Over-estimates, never under |
| "What are the top-K items?" | CMS + min-heap | Approximate top-K |
| Exact membership, any size | Hash set | Exact, O(n) memory |

**Use probabilistic structures when:**
- The input stream is too large to store fully (logs, network packets, web requests)
- You can tolerate small error rates in exchange for huge memory savings
- The cost of a false positive is low (checking a cache, pre-filtering before a real lookup)

**Don't use them when:**
- You need exact counts (billing, compliance, deduplication)
- False positives are expensive (security checks, bank transactions)

---

## Summary

- **Bloom Filter**: O(k) insert/query, O(m) space. No false negatives; controls false positive rate via m and k. Use for fast membership pre-checking.
- **HyperLogLog**: O(1) add, O(1) count, ~16 KB fixed memory. ~1-2% error for cardinality estimation. Used in Redis, Google Analytics, PostgreSQL.
- **Count-Min Sketch**: O(d) update/query (d = depth, usually 4-8). Always over-estimates. Used for frequency estimation in streaming systems.
- All three use double hashing / multiple hash functions to distribute elements uniformly.

## Exercises

### Easy
1. Measure the empirical false positive rate of your Bloom Filter: insert 1M random strings, then query 10K strings that were NOT inserted. What fraction returns true? How close is it to the theoretical 1%?
2. Show that a Bloom filter with all bits set to 1 reports every element as present. Then implement `Reset()` that zeroes all bits.
3. Add a `Union(a, b *Filter) *Filter` function that combines two Bloom filters (bitwise OR of their bit arrays). When is the union valid? When is the intersection (bitwise AND) valid?

### Medium
4. Implement a **Counting Bloom Filter** that supports deletion: replace each bit with a small counter (4 bits). Increment on insert, decrement on delete, report false when any counter is 0. What new failure modes does this introduce?
5. Test HyperLogLog accuracy: estimate cardinality for inputs of size 10, 100, 1K, 10K, 100K, 1M. Calculate the error percentage at each scale. At what scale does the accuracy stabilize?
6. Build a **sliding window Bloom filter**: maintain two filters, rotating every N elements. Only the current and previous windows are kept. Query returns true if either filter contains the element. This approximates a membership set with an expiry window.

### Hard
7. Build a **Top-K tracker** using Count-Min Sketch: maintain a min-heap of size K alongside the sketch. On each increment, update the sketch and check if the new estimated count displaces the current K-th largest. Handle hash collisions gracefully by verifying candidates against the sketch before inserting them in the heap.
8. Implement **Cuckoo Filter** as an alternative to Bloom Filter: store fingerprints (short hashes) in a compact array with two candidate buckets per item. Cuckoo filters support deletion and have better lookup performance than Bloom filters. Compare memory efficiency and FPR for n=1M at 1% target FPR.
