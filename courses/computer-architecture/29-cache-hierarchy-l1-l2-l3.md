# Chapter 29: Cache Hierarchy — L1, L2, and L3

One cache level is not enough. The L1 cache must be fast enough to deliver data every cycle without stalling the CPU — but that constraint limits its size. A larger L2 cache catches the misses that L1 can't, but is slower. An even larger L3 cache acts as a final buffer before reaching the very slow DRAM. This multi-level hierarchy is not a compromise — it is a carefully engineered system where each level is optimized for a different point on the size/speed trade-off curve.

## Table of Contents

1. [Why One Cache Level Isn't Enough](#1-why-one-cache-level-isnt-enough)
2. [L1 Cache — The Speed King](#2-l1-cache--the-speed-king)
3. [L2 Cache — The Middle Manager](#3-l2-cache--the-middle-manager)
4. [L3 Cache — The Last Resort Before DRAM](#4-l3-cache--the-last-resort-before-dram)
5. [Inclusive, Exclusive, and Non-Inclusive Policies](#5-inclusion-policies)
6. [Cache Coherence Basics — The MESI Protocol](#6-cache-coherence-basics--the-mesi-protocol)
7. [Cache Prefetching](#7-cache-prefetching)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why One Cache Level Isn't Enough

The fundamental tension: the CPU wants a cache that is simultaneously **very fast** (answer in 1 cycle) and **very large** (hold the entire working set). These goals conflict.

A larger cache has more lines to search and longer wire delays — slower. L1 caches are close to the core (physically small on chip), using the fastest SRAM cells, designed for single-cycle access. Making L1 larger would slow it down below the 1-cycle target.

```
Access latency vs cache size trade-off (modern 3 GHz CPU):
  Size    |  Latency (cycles) | What it's used for
  --------|-------------------|--------------------
  32 KB   |  4 cycles         | L1 (per-core, hot data)
  256 KB  |  12 cycles        | L2 (per-core, warm data)
  8-64 MB |  30-50 cycles     | L3 (shared, cold data)
  16+ GB  |  200 cycles       | DRAM
```

Each level handles misses from the level above. L1 misses go to L2; L2 misses go to L3; L3 misses go to DRAM.

### Miss Rate Cascade

With typical workloads:
- L1 hit rate: ~95% → 5% of requests go to L2
- L2 hit rate: ~85% (of misses from L1) → 0.75% of original requests go to L3
- L3 hit rate: ~90% (of misses from L2) → 0.075% of original requests go to DRAM

Only 1 in 1,333 accesses reaches DRAM. The hierarchy is remarkably effective.

### Quick Check
> 1. Why does increasing the L1 cache size make it slower?
> 2. If L1 hit rate is 95%, L2 hit rate is 85%, and L3 hit rate is 90%, what fraction of requests reach DRAM?
> 3. The L1 must deliver data in 1 cycle. At 5 GHz, what is the maximum round-trip time in nanoseconds?

---

## 2. L1 Cache — The Speed King

The L1 (Level 1) cache is the closest cache to the CPU execution units. It is split into two separate caches: **L1-I** (instructions) and **L1-D** (data). This split eliminates structural hazards — the fetch unit and load/store unit can access the cache simultaneously.

**Typical modern L1 specs:**

| CPU | L1-I | L1-D | Associativity | Latency |
|-----|------|------|--------------|---------|
| Intel Core (Raptor Lake) | 32KB | 48KB | 8-way | 4 cycles |
| AMD Zen 4 | 32KB | 32KB | 8-way | 4 cycles |
| Apple M2 Firestorm | 192KB | 128KB | 12-way | 3-4 cycles |
| ARM Cortex-A77 | 64KB | 64KB | 4-way | 4 cycles |

Apple's absurdly large L1 (up to 192KB for instructions) is a major contributor to M-series CPUs' exceptional IPC — more hot code fits in the ultra-fast L1, reducing pressure on L2.

**L1 must support:**
- One load and one store per cycle (to avoid structural hazards with out-of-order execution)
- ECC or parity bits for data integrity
- Snooping for cache coherence in multicore systems

### Critical Path and Banking

L1 is on the critical timing path — if it takes longer than the processor's clock period, the whole CPU must slow down. Modern L1 is implemented as a banked array: multiple sub-arrays that can be accessed in parallel. If an address maps to bank 3, only bank 3 is read while banks 0, 1, 2, 4... remain idle. This allows fast access (fewer cells to search per bank) while keeping the total size large.

### Quick Check
> 1. Why is the L1 split into separate instruction and data caches?
> 2. Apple's M2 has 192KB L1-I vs Intel's 32KB. Why might this matter for performance?
> 3. What does "cache banking" accomplish in terms of speed?

---

## 3. L2 Cache — The Middle Manager

The L2 cache is the first line of defense against L1 misses. It is larger and slower than L1 but still fast enough to avoid going to L3 or DRAM for most misses.

**Typical modern L2 specs:**

| CPU | L2 Size | Associativity | Latency |
|-----|---------|--------------|---------|
| Intel Core i9-13900K | 2MB/core | 16-way | 12 cycles |
| AMD Zen 4 (Ryzen 9) | 1MB/core | 8-way | 12 cycles |
| Apple M2 Firestorm | 16MB (shared) | — | ~15 cycles |
| ARM Cortex-A77 | 256KB–512KB | 8-way | 10-12 cycles |

L2 is usually **unified** (not split into I and D). By the time a request reaches L2, the distinction between instruction and data is less important since L1 has already resolved most of the traffic.

### L2 and the Load-to-Use Penalty

When the CPU issues a load instruction and L1 misses:
1. The memory access hardware generates an L2 lookup request
2. L2 responds in ~12 cycles
3. The loaded value is forwarded directly to the dependent instruction (if it's ready in the OOO window)

For an out-of-order processor with a large instruction window (e.g., Apple M1's 630-entry ROB), a 12-cycle L2 miss can often be hidden by executing other independent instructions while waiting. This is why deep OOO machines tolerate L2 latency better than simple in-order machines.

### Quick Check
> 1. Why is L2 unified (not split into I and D) while L1 is split?
> 2. How does a deep OOO instruction window help hide L2 miss latency?
> 3. AMD Zen 4 has 1MB of L2 per core; Intel Raptor Lake has 2MB. For a workload with a 1.5MB working set, which processor avoids more L3 accesses?

---

## 4. L3 Cache — The Last Resort Before DRAM

The L3 cache is the last level of on-chip cache (LLC — Last Level Cache). It is large, shared between all cores, and much slower than L1 or L2, but still enormously faster than DRAM.

**Typical modern L3 specs:**

| CPU | L3 Size | Configuration | Latency |
|-----|---------|---------------|---------|
| Intel Core i9-13900K | 36MB | Shared, all cores | 40-45 cycles |
| AMD Ryzen 9 7950X | 64MB + 96MB 3D V-Cache | Shared | 40-50 cycles |
| Apple M2 Ultra | 192MB | Shared | ~60 cycles |
| AMD Threadripper 7980X | 256MB | Shared | 50-60 cycles |

**AMD's 3D V-Cache**: AMD stacks an extra SRAM die directly on top of the CPU die using 3D packaging. The Ryzen 9 7950X with 3D V-Cache has 96MB of extra L3 — particularly beneficial for games where the working set (textures, AI data, geometry) fits in the huge L3 but not in a normal-sized L3.

### L3 and Multicore Sharing

When Core 0 modifies a cache line that Core 1 has in its L1, the coherence protocol must ensure Core 1 gets the updated value. The L3 acts as the coherence hub — it holds the "true" version of shared data and coordinates between cores' private L1/L2 caches.

```
Core 0 L1 ──→ Core 0 L2 ──→
                              Shared L3 ──→ DRAM
Core 1 L1 ──→ Core 1 L2 ──→
```

### Quick Check
> 1. What does "LLC" stand for and which cache level is it?
> 2. How does a huge L3 (like AMD's 3D V-Cache) benefit gaming workloads specifically?
> 3. Why is the L3 shared between cores while L1 and L2 are private per-core?

---

## 5. Inclusion Policies

When a line is in L1, is it also in L2 and L3? The answer depends on the **inclusion policy**.

### Inclusive Cache

Every line in L1 is also present in L2 and L3. The set of data in L1 is a **subset** of L2, which is a subset of L3.

**Advantages**:
- Cache coherence is simpler: when a line must be invalidated (e.g., for coherence), you only need to check L3 — if it's not there, it's not in any private cache.
- On a cache-to-cache transfer, L3 is the authoritative source.

**Disadvantages**:
- Wastes L2/L3 space — lines in the private L1 also occupy L3 slots. For a core with 32KB L1 and a 36MB L3, this costs 32KB of L3 capacity per core.

Intel has historically used inclusive L3.

### Exclusive Cache

A line is in at most one level of the hierarchy. If it's in L1, it's not in L2. If evicted from L1, it moves to L2.

**Advantages**: Maximum total cache capacity (no duplication).  
**Disadvantages**: Cache coherence is harder (must search all levels to find a line).

Early AMD processors used exclusive L2.

### Non-Inclusive (Neither Inclusive Nor Exclusive) — NINE

A line in L1 may or may not be in L2. L2 does not maintain strict inclusion or strict exclusion. This is a practical compromise used by modern processors.

AMD Zen: the L2 and L3 relationship is NINE. Apple M-series uses a victim cache approach (evictions from L1 go to L2, but L2 doesn't necessarily hold the same lines as L1). Intel shifted toward NINE in recent designs too.

### Quick Check
> 1. In an inclusive L3 cache, if a line is in L2 of Core 1, is it also in L3? Why?
> 2. What is the main disadvantage of inclusive caches?
> 3. In an exclusive L1/L2 hierarchy, what happens to a cache line when it is evicted from L1?

---

## 6. Cache Coherence Basics — The MESI Protocol

With multiple cores, each with their own private L1/L2 caches, a problem arises: two cores might have copies of the same cache line, and one modifies it. Without coherence, the other core reads a stale value.

**Cache coherence** guarantees that all cores agree on the current value of every memory location.

### MESI Protocol

MESI is the most widely used cache coherence protocol. Each cache line has one of four states:

| State | Meaning |
|-------|---------|
| **M** odified | This cache has the only copy, and it has been modified. Memory is stale. Must write back on eviction. |
| **E** xclusive | This cache has the only copy, and it matches memory. Can write without broadcasting. |
| **S** hared | Multiple caches may have this line. Matches memory. Must invalidate others before writing. |
| **I** nvalid | This line is not valid. Must fetch before use. |

**State transitions (simplified):**
```
Read a line not in cache:
  → Check other caches (snoop bus/directory)
  → If nobody has it: fetch from memory → state = E
  → If another core has it in S: fetch → state = S (both now S)
  → If another core has it in M: that core writes back → fetch → state = S/E

Write to a line:
  → If state = M: just write (already exclusive and modified)
  → If state = E: write and transition to M
  → If state = S or I: broadcast INVALIDATE to other cores → they set their copy to I
                       fetch line exclusively → write → state = M
```

The MESI protocol ensures coherence: when you write, all other cached copies are invalidated. When you read, you get the most recent version.

**False sharing**: Two threads write to different variables that happen to occupy the same 64-byte cache line. The line ping-pongs between cores (both wanting the line in M state), even though there's no true data sharing. A classic performance bug in multithreaded code.

### Quick Check
> 1. What does the MESI acronym stand for?
> 2. Core A has a line in state M. Core B tries to read that same line. What happens?
> 3. What is "false sharing" and why is it a performance problem?

---

## 7. Cache Prefetching

Even with a perfectly organized cache hierarchy, cold and capacity misses are unavoidable. **Prefetching** is the technique of loading data into cache before it is explicitly requested, hiding the miss latency.

### Hardware Prefetchers

Modern CPUs contain multiple hardware prefetchers that observe the stream of cache misses and try to predict future accesses:

**Stride prefetcher**: Detects a regular stride in access addresses. If it sees accesses at addresses A, A+64, A+128, it predicts A+192, A+256... and prefetches them. Excellent for sequential array accesses.

**Stream prefetcher**: Detects sequential access patterns and prefetches an entire "stream" of consecutive lines. Intel's processors have multiple stream prefetchers, each tracking a different memory region simultaneously.

**Spatial prefetcher (Next-Line)**: Always prefetch the next cache line after any miss. Simple, effective for sequential code.

**Indirect/complex pattern prefetchers**: Newer AMD Zen 4 and Apple M-series have advanced prefetchers that can follow pointer chains (linked lists, trees) — very hard to do correctly.

### Software Prefetching

The ISA provides explicit prefetch instructions:
- x86: `PREFETCHNTA` (prefetch, no temporal hint), `PREFETCHT0` (prefetch to all cache levels)
- ARM: `PLD` (preload data)

Compilers can insert these automatically (auto-vectorization/unrolling often includes prefetch). Hand-tuned code for HPC/database workloads manually inserts prefetch instructions to hide latency on known access patterns.

**Prefetch too early**: The line might be evicted before it's needed.  
**Prefetch too late**: The data isn't in cache when needed.  
**Wrong prefetch**: Wastes bandwidth, evicts useful data.

The hardware prefetcher is "free" in terms of programmer effort, but getting software prefetch right requires careful tuning.

### Quick Check
> 1. What is a "stride prefetcher" and what access pattern does it handle well?
> 2. Why might prefetching too early be harmful?
> 3. For a linked-list traversal (following `next` pointers), does a stride prefetcher help? Why or why not?

---

## Summary

- The cache hierarchy (L1 → L2 → L3 → DRAM) balances speed and capacity: each level is slower but larger than the previous.
- **L1** (32–192KB, 4 cycles) is split into instruction and data caches, optimized for single-cycle access.
- **L2** (256KB–16MB, 12 cycles) is unified, per-core, and catches most L1 misses.
- **L3** (8MB–256MB, 40–60 cycles) is shared between all cores and acts as the last on-chip barrier before DRAM.
- **Inclusion policies** determine whether a line in L1 also exists in L2/L3. Modern designs use NINE (Non-Inclusive/Non-Exclusive) for better capacity utilization.
- **MESI protocol** ensures cache coherence between cores. Lines can be Modified, Exclusive, Shared, or Invalid. False sharing is a common multicore performance pitfall.
- **Prefetching** (hardware stride prefetchers, software `PREFETCH` instructions) hides miss latency by loading cache lines before they're needed.

---

## Exercises

### Easy
1. A CPU has L1 hit latency = 4 cycles (miss rate 5%), L2 hit latency = 12 cycles (miss rate 20% of L1 misses), L3 hit latency = 40 cycles (miss rate 10% of L2 misses), DRAM latency = 200 cycles. Calculate the average memory access time.
2. What does "false sharing" mean? Give a code example in C where two threads write to adjacent array elements and explain why this causes false sharing.
3. What are the four MESI states? For each, write one sentence describing when a line is in that state.

### Medium
4. Core A has cache line X in state M (modified). Core B tries to read line X. Trace the MESI protocol steps: what messages are exchanged, what state does each cache end up in, and where does Core B get the data from?
5. AMD's 3D V-Cache adds 64MB of extra L3 on top of the Zen 4 die. For a game with 80MB of texture/AI working set: (a) Without 3D V-Cache (16MB L3), what fraction of the working set fits? (b) With 3D V-Cache (80MB L3), what fraction fits? (c) How does this affect frame rate for that game?
6. A software engineer writes a multithreaded counter array: `int counters[16]`, with thread i incrementing `counters[i]`. With 64-byte cache lines and 4-byte ints, how many threads share the same cache line? How would you fix this to avoid false sharing?

### Hard
7. The inclusive L3 policy wastes L3 space by duplicating lines already in private L1/L2. Suppose a 16-core CPU has: L1 = 48KB/core, L2 = 2MB/core, L3 = 32MB (shared). With inclusive L3: (a) How much L3 space is "wasted" by holding copies of all L1 and L2 data? (b) What fraction of L3 is effectively wasted? (c) Argue whether NINE policy would significantly improve L3 effective capacity for a workload where each core has a 3MB working set.
8. Describe the hardware needed to implement a stride prefetcher. It should: detect a stride pattern after 3 accesses, prefetch the next 4 cache lines, and handle multiple simultaneous streams. Specify the data structures (tables, counters) and the update logic. How many simultaneous streams can you support with 512 bytes of hardware state?
