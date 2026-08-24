# Chapter 28: Cache Memory — The Speed Gap Solution

There is a brutal gap in modern computers: the CPU can execute an instruction in 0.3 nanoseconds (at 3 GHz), but accessing RAM takes 60–100 nanoseconds — 200–300 times longer. If the CPU had to wait for memory on every instruction, it would be idle more than 99% of the time. Cache memory is the solution: a small, extremely fast memory placed between the CPU and RAM that stores copies of recently accessed data, so most memory accesses find what they need in 1–4 cycles rather than 200.

## Table of Contents

1. [The Memory Speed Gap](#1-the-memory-speed-gap)
2. [Locality — Why Caches Work](#2-locality--why-caches-work)
3. [Cache Basics: Hit, Miss, Line, Size](#3-cache-basics-hit-miss-line-size)
4. [Cache Organization: Direct-Mapped, Set-Associative, Fully-Associative](#4-cache-organization)
5. [Replacement Policies](#5-replacement-policies)
6. [Write Policies — Write-Through vs Write-Back](#6-write-policies--write-through-vs-write-back)
7. [The Three Cs of Cache Misses](#7-the-three-cs-of-cache-misses)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Memory Speed Gap

```
Memory hierarchy access times (approximate, modern hardware):
  CPU Register:      0.3 ns  (1 cycle at 3 GHz)
  L1 Cache:          1 ns    (4 cycles)
  L2 Cache:          3 ns    (12 cycles)
  L3 Cache:          10 ns   (30-40 cycles)
  DRAM (RAM):        60 ns   (200 cycles)
  NVMe SSD:          100 µs  (300,000 cycles!)
  HDD:               10 ms   (30,000,000 cycles)
```

If every instruction needed to access DRAM, a 3 GHz processor would actually run at the equivalent of ~15 MHz — 200× slower than its nominal speed. This is the **memory wall** — and cache is the dam that holds it back.

Why is DRAM so slow? DRAM stores bits as tiny electrical charges in capacitors. These capacitors leak and must be **refreshed** thousands of times per second. The access requires charging/discharging row lines, selecting column lines, sensing tiny voltage differences — all of which takes tens of nanoseconds. DRAM is cheap per bit and dense, but inherently slow.

SRAM (used for caches) stores bits in a 6-transistor latch — no charge to refresh, stable, fast — but uses 6× more transistors per bit and is much more expensive. A 1MB SRAM cache uses more silicon than 1GB of DRAM.

### Quick Check
> 1. How many times slower is DRAM than an L1 cache?
> 2. Why is DRAM slower than SRAM even though both use transistors?
> 3. If a CPU executes 1 billion instructions/second and 10% require DRAM access (200 cycles each), how many effective billion "useful" instructions/second does it achieve?

---

## 2. Locality — Why Caches Work

Cache works because programs exhibit **locality**: they tend to access the same data repeatedly, and nearby data together. Without locality, a small cache holding recently used data would have a very low hit rate and be useless.

### Temporal Locality
"If you used it recently, you'll probably use it again soon."

```
for (int i = 0; i < N; i++) {
    sum += array[i];    // 'sum' accessed every iteration — classic temporal locality
}
```

The variable `sum` is accessed N times. After the first access loads it into cache, all subsequent N-1 accesses find it there instantly.

### Spatial Locality
"If you used it, you'll probably use something nearby soon."

```
for (int i = 0; i < N; i++) {
    sum += array[i];    // array[0], array[1], array[2]... are adjacent in memory
}
```

When `array[0]` is loaded, the cache brings a full **cache line** (typically 64 bytes = 16 ints) into cache. So `array[1]` through `array[15]` are already in cache when the loop reaches them — no additional memory access needed.

Spatial locality is why caches load data in **cache lines** (chunks) rather than individual bytes. The gamble: if you used byte 100, you'll probably soon use bytes 101-163. This gamble wins almost always for arrays, structs, code (sequential instructions), and stack variables.

### The Compiler and Cache

Programmers and compilers can improve cache performance by organizing data to maximize spatial locality:
- Store frequently accessed data in arrays (sequential) rather than linked lists (pointers scattered in memory)
- Keep "hot" data together in structs (avoid padding that wastes cache space)
- Tile matrix operations to fit working set in cache

### Quick Check
> 1. Define temporal locality and give an example other than the one above.
> 2. Why does loading data in 64-byte cache lines work well for arrays?
> 3. Why does a linked list have poor spatial locality compared to an array?

---

## 3. Cache Basics: Hit, Miss, Line, Size

### Cache Hit and Miss

When the CPU needs a memory address:
1. It checks the cache first
2. **Cache hit**: the data is in cache → return it in 1–4 cycles
3. **Cache miss**: the data is NOT in cache → fetch from the next level (L2, L3, RAM) and store a copy in cache → takes 100–300 cycles

**Hit rate** = fraction of accesses that hit in cache. A hit rate of 95% means 1 in 20 accesses goes to the next level. This sounds good, but those 5% misses at 200 cycles each dominate if the hit is only 1 cycle.

```
Effective access time = hit_rate × hit_time + (1 - hit_rate) × miss_time
= 0.95 × 1 + 0.05 × 200 = 0.95 + 10 = 10.95 cycles average
```

Going from 95% to 99% hit rate: `0.99 × 1 + 0.01 × 200 = 2.99 cycles` — a 3.7× improvement from that 4% hit rate increase.

### Cache Line

The cache does not store individual bytes — it stores **cache lines** (also called cache blocks). A cache line is typically 64 bytes on modern x86/ARM CPUs (some embedded: 32 bytes; some server: 128 bytes).

When there's a miss, the entire 64-byte line containing the requested address is fetched from the next level. This exploits spatial locality.

```
Cache line containing address 0x1040 (assuming 64-byte lines):
  Line start: 0x1040 & ~0x3F = 0x1040  (lower 6 bits zeroed)
  Line end:   0x107F
  All 64 bytes (0x1040–0x107F) loaded together
```

### Cache Size and Hit Rate

Bigger caches have higher hit rates — more data fits, so more accesses find what they need. But bigger caches are slower (longer to search) and more expensive.

```
Typical L1 cache sizes and hit rates for integer code:
  8KB:   ~90% hit rate
  32KB:  ~95% hit rate
  64KB:  ~97% hit rate
  256KB: ~99% hit rate
```

### Quick Check
> 1. A cache has a 90% hit rate. Hit time is 2 cycles; miss penalty (time to go to RAM) is 200 cycles. What is the effective memory access time?
> 2. Why are cache lines 64 bytes rather than 1 byte or 4096 bytes?
> 3. If you double the cache size, does the hit rate double too? Explain.

---

## 4. Cache Organization

How does the hardware know if a given address is in the cache, and where to find it?

### Direct-Mapped Cache

Each memory address maps to exactly one cache location. The address is split into:
- **Tag** bits: identify which memory block occupies this cache slot
- **Index** bits: select which cache slot to check
- **Offset** bits: select which byte within the cache line

```
Address bits [31 ... 12 | 11 ... 6 | 5 ... 0]
              Tag (20b)   Index (6b)  Offset (6b)
              
For a 4KB direct-mapped cache with 64-byte lines:
  4KB / 64B = 64 lines → 6 index bits
  64-byte line → 6 offset bits
  Remaining 20 bits = tag
```

**Lookup**: Use index bits to select a cache slot. Compare the stored tag with the address's tag. If they match (and the valid bit is set) → HIT; else → MISS.

**Problem**: Two addresses that differ only in the tag but share the same index will evict each other every time both are accessed — **conflict miss**. Example: in a 4KB direct-mapped cache, addresses 0x0000 and 0x1000 map to the same slot. Accessing them alternately causes 100% miss rate.

### Set-Associative Cache

An **N-way set-associative** cache groups N cache lines into a **set**. Any address can map to any of the N lines within its set.

```
4-way set-associative cache:
Set 0: [line 0a][line 0b][line 0c][line 0d]   ← 4 possible homes for addresses mapping to set 0
Set 1: [line 1a][line 1b][line 1c][line 1d]
...

Address maps to a set (via index bits), and is compared against all 4 tags in parallel.
```

4-way or 8-way set-associative is the most common configuration. Higher associativity → fewer conflict misses but more complex/slower lookup hardware (must compare N tags in parallel).

### Fully-Associative Cache

Any address can go in any cache slot — maximum flexibility, zero conflict misses. But searching the entire cache for a tag is only feasible for very small caches (TLB, small victim cache). All modern L1/L2/L3 caches use set-associative organization.

### Quick Check
> 1. In a direct-mapped cache, what is a "conflict miss"?
> 2. How does set-associativity reduce conflict misses?
> 3. Why can't you make the L1 cache fully-associative?

---

## 5. Replacement Policies

When a miss occurs and the target set is full, which line gets evicted to make room?

**LRU (Least Recently Used)**: Evict the line that was accessed least recently. Intuition: if you haven't used it in a while, you're less likely to need it soon. Optimal for many access patterns but requires tracking access recency for all lines in each set — complex for high associativity.

**Pseudo-LRU**: An approximation of LRU using a binary tree of "last used" bits — much cheaper hardware while capturing most of LRU's benefit. Used in most modern CPUs.

**Random**: Evict a randomly chosen line. Surprisingly effective — avoids pathological cases that defeat LRU. Simple hardware.

**FIFO**: Evict the line that has been in the cache longest. Simple but worse than LRU for most workloads.

**Optimal (Bélády's Algorithm)**: Evict the line that will be used furthest in the future. Theoretically optimal but requires knowing the future — impossible in hardware, used only as a theoretical benchmark.

For most workloads, Pseudo-LRU performs nearly as well as true LRU while being much simpler. Intel uses a variant called PLRU (bit-PLRU or tree-PLRU) in its L1 and L2 caches.

### Quick Check
> 1. Why does LRU work well for temporal locality workloads?
> 2. In a 4-way set-associative cache with LRU replacement, trace the hits and misses for access pattern: A B C D A B C D A (all map to same set). How many misses?
> 3. For a workload that scans a large array sequentially (each element accessed exactly once), which replacement policy is best and why?

---

## 6. Write Policies — Write-Through vs Write-Back

When the CPU writes to a cached address, what happens?

### Write-Through

Every write to cache is **immediately** also written to the next level (L2 or RAM). The cache and memory always agree on the current value.

**Advantages**: Simple, memory is always up-to-date, no data loss on cache eviction.  
**Disadvantages**: Every store requires a memory write — high memory bandwidth pressure. Usually combined with a **write buffer** (a small FIFO queue that absorbs write traffic so the CPU doesn't stall).

### Write-Back (Copy-Back)

Writes update the cache only. The value in main memory may be stale. A **dirty bit** in each cache line indicates whether the cached value differs from memory. When a dirty line is evicted, it must be written back to the next memory level before the slot is reused.

**Advantages**: Dramatically reduces memory bandwidth (multiple writes to the same line are absorbed in cache, only one writeback needed).  
**Disadvantages**: More complex; on eviction, must write dirty data back; memory is inconsistent with cache.

Modern CPUs universally use **write-back** for L1/L2/L3. The bandwidth savings are enormous — a tight loop modifying a local variable would otherwise flood the memory bus with writes.

### Write Allocation vs No-Write Allocation

On a **write miss** (writing to an address not in cache):
- **Write-allocate** (most common): Fetch the cache line containing the address into cache, then perform the write. Future reads of nearby addresses benefit from spatial locality.
- **No-write-allocate**: Write directly to the next level without fetching the line into cache. Used with write-through for streaming write workloads.

### Quick Check
> 1. What is the "dirty bit" and when is it set?
> 2. Why does write-back reduce memory bandwidth compared to write-through?
> 3. A processor uses write-back with write-allocate. Describe exactly what happens when the CPU writes to address 0x2000 which is not currently in cache.

---

## 7. The Three Cs of Cache Misses

Cache misses are classified into three categories, each requiring a different solution:

### Compulsory Misses (Cold Misses)
The very first access to any cache line — the data has never been in cache. **Unavoidable** on the first access. Prefetching can help by loading lines before they are explicitly needed.

*Solution*: **Hardware prefetcher** — detects stride patterns and prefetches upcoming lines. **Software prefetch** instructions (PREFETCH in x86) let compilers explicitly request data in advance.

### Capacity Misses
The cache is too small to hold the working set. Even with perfect placement, lines keep getting evicted before being reused.

*Solution*: Bigger cache. Or — restructure the algorithm to reduce the working set (e.g., **cache blocking/tiling** for matrix operations).

### Conflict Misses
Two frequently used addresses map to the same cache set and evict each other. Only possible in direct-mapped and low-associativity caches.

*Solution*: Higher associativity. Or — **cache-aware data placement** by the compiler/OS to avoid conflicts.

```
Miss rate breakdown for typical L1 miss (32KB, 8-way):
  Compulsory: ~5%  (first-time accesses)
  Capacity:   ~40% (working set exceeds cache)
  Conflict:   ~55% (replaced before reuse — even with 8-way, some still conflict)
```

### Quick Check
> 1. Which type of miss can never be eliminated regardless of cache size or associativity?
> 2. A matrix multiply algorithm accesses elements in a pattern that causes many capacity misses. What algorithmic technique can reduce these?
> 3. Increasing associativity from 4-way to 8-way reduces which type of miss? Which two types are unaffected?

---

## Summary

- The **memory speed gap** (DRAM is 200× slower than L1 cache) would make fast processors nearly useless without caches.
- Caches exploit **temporal locality** (reuse) and **spatial locality** (nearby data) to achieve high hit rates with small, fast storage.
- A **cache line** (64 bytes) is loaded on a miss — exploiting spatial locality.
- Cache organization: **direct-mapped** (simple, conflict-prone), **set-associative** (balance of simplicity and flexibility), **fully-associative** (flexible, complex, used only for tiny structures like TLBs).
- **LRU/Pseudo-LRU** replacement evicts the least recently used line — effective for most workloads.
- **Write-back** policy reduces memory bandwidth by absorbing multiple writes to the same line; **dirty bit** tracks lines that need writeback.
- The **three Cs of misses**: Compulsory (unavoidable first access), Capacity (working set too large), Conflict (mapping collision in low-associativity caches).

---

## Exercises

### Easy
1. An L1 cache has 32KB, 64-byte lines, and 8-way associativity. How many sets are there? How many index bits? How many offset bits?
2. Calculate effective memory access time for: 99% hit rate, 4-cycle hit time, 200-cycle miss penalty (to DRAM).
3. What is the dirty bit for? When is it set, and when is it cleared?

### Medium
4. Consider a direct-mapped 8KB cache with 64-byte lines. Addresses 0x0000 and 0x2000 — do they conflict? Show your work using tag/index/offset decomposition.
5. A matrix multiply: `C[i][j] += A[i][k] * B[k][j]`. The inner loop over k accesses B column-major (poor spatial locality). Explain why this causes capacity/conflict misses and how loop tiling (blocking) fixes the problem.
6. A hardware prefetcher detects a stride-2 access pattern (accessing array[0], array[2], array[4]...). It prefetches array[6] when array[4] is accessed. If the miss penalty is 200 cycles and the prefetch accurately predicts the next 3 accesses, how many cycles of stall are eliminated?

### Hard
7. A 4KB direct-mapped cache (64-byte lines) is used with this code:
   ```c
   float a[256], b[256], c[256];
   for (int i=0; i < 256; i++) c[i] = a[i] + b[i];
   ```
   Each float is 4 bytes. Arrays a, b, c are placed at addresses 0x0000, 0x0400, 0x0800.
   (a) Do a, b, c conflict in the cache? (b) What is the miss rate on each array? (c) How would you reorder or pad the arrays to eliminate conflicts?

8. Design a fully-associative 4-entry cache (for simplicity, ignore offset bits — each entry holds exactly one word). Simulate the access sequence A B C D A B C D A B C D using LRU replacement. Count hits and misses. Then repeat with FIFO replacement. Which performs better for this cyclic pattern and why?
