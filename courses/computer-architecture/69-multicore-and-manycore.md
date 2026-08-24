# Chapter 69: Multicore and Manycore — Parallel by Design

After Dennard scaling failed in 2004, the industry faced a wall: transistors were growing smaller but couldn't run faster without burning. The solution was to put multiple cores on the same die — instead of one fast core, make four (or forty, or forty-thousand) cores that collectively do more work. This chapter explains why multicore emerged, how it works at the hardware level, what challenges it creates for software, the difference between multicore and manycore, and how modern processors manage many parallel execution units.

## Table of Contents

1. [Why Multicore? The End of the Free Lunch](#1-why-multicore-the-end-of-the-free-lunch)
2. [Symmetric Multiprocessing (SMP)](#2-symmetric-multiprocessing-smp)
3. [Cache Coherence in Multicore Systems](#3-cache-coherence-in-multicore-systems)
4. [Memory Consistency Models](#4-memory-consistency-models)
5. [NUMA — Non-Uniform Memory Access](#5-numa--non-uniform-memory-access)
6. [Manycore Architectures — Hundreds of Cores](#6-manycore-architectures--hundreds-of-cores)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Why Multicore? The End of the Free Lunch

Herb Sutter's famous 2004 Dr. Dobb's article "The Free Lunch Is Over" announced that software engineers could no longer expect free performance gains from processor speed increases. The reason: Dennard scaling had broken.

**Before 2004 (the free lunch era):**
- Each CPU generation: same power, 2× more transistors, 43% faster clock
- Existing single-threaded software ran 43% faster without any changes
- Software engineers could write inefficient code and rely on hardware progress

**After 2004 (the multicore era):**
- CPU clock speed plateaued at ~3–4 GHz (thermal/power wall)
- Chip makers added more cores instead of faster cores
- Single-threaded software no longer got faster automatically
- Multi-threaded software could use additional cores

```
Clock frequency trend:
  
  1980-2004: doubles every 18 months (Dennard scaling)
  2004-2024: stagnant at 3-5 GHz (power wall)
  
  Core count trend:
  2003: 1 core (Intel Pentium 4)
  2005: 2 cores (Intel Core Duo)
  2007: 4 cores (Intel Core 2 Quad)
  2011: 6-8 cores (Intel Xeon E5)
  2014: 18 cores (Intel Xeon E5 v3)
  2019: 64 cores (AMD EPYC Rome)
  2023: 96 cores (AMD EPYC Genoa, 192 threads with SMT)
  2024: 128 cores (Intel Xeon 6 Granite Rapids)
```

**Amdahl's Law** (1967): The maximum speedup from parallelization is limited by the sequential fraction:

```
Speedup = 1 / (s + (1-s)/p)

Where:
  s = fraction of code that is sequential (cannot be parallelized)
  p = number of processors
  (1-s) = parallelizable fraction

Example: 80% parallelizable code (s=0.2), 8 cores (p=8):
  Speedup = 1 / (0.2 + 0.8/8) = 1 / (0.2 + 0.1) = 1 / 0.3 = 3.3×

Even with 8 cores: only 3.3× speedup due to 20% sequential code
With 1000 cores: 1 / (0.2 + 0/1000) = 1/0.2 = 5× maximum ever
```

Amdahl's Law explains why 128 cores doesn't give 128× speedup: most real programs have significant sequential portions.

### Quick Check
> 1. What is Herb Sutter's "free lunch" and when did it end?
> 2. State Amdahl's Law and explain its implication for multicore speedup.
> 3. If a program is 95% parallelizable, what is the maximum speedup with infinite cores?

---

## 2. Symmetric Multiprocessing (SMP)

**SMP (Symmetric Multiprocessing)** is the most common multicore organization: multiple identical cores on one chip, each with access to the same shared memory.

```
SMP architecture:
  
  Core 0       Core 1       Core 2       Core 3
  ┌──────┐    ┌──────┐    ┌──────┐    ┌──────┐
  │  L1  │    │  L1  │    │  L1  │    │  L1  │
  │  L2  │    │  L2  │    │  L2  │    │  L2  │
  └──┬───┘    └──┬───┘    └──┬───┘    └──┬───┘
     │           │           │           │
     └───────────┴───────────┴───────────┘
                         │
              ┌──────────┴──────────┐
              │      L3 Cache        │
              │  (shared by all)     │
              └──────────┬──────────┘
                         │
              ┌──────────┴──────────┐
              │   Memory Controller  │
              └──────────┬──────────┘
                         │
                    DRAM (shared)
```

**Symmetric**: Every core sees the same address space and can run any thread. The OS scheduler can assign any thread to any core.

**Shared memory**: All cores access the same physical DRAM. This is the programming convenience of multicore — threads share data through shared memory (though with synchronization requirements).

**SMT (Simultaneous Multithreading)**: Each physical core can appear as 2 (or more) logical cores to the OS. Intel calls this Hyper-Threading. Each logical core shares execution units, but has separate register state, PC, and reorder buffer. When one thread stalls (cache miss, branch misprediction), the other thread uses the idle execution units.

**Hybrid multicore** (Intel Alder Lake, ARM big.LITTLE): Not all cores are identical. Different core types for different workloads:
- P-cores (Performance): large, fast, power-hungry → for foreground tasks
- E-cores (Efficiency): small, lower power → for background tasks
- OS thread scheduler assigns threads to cores based on power/performance requirements

### Quick Check
> 1. What is SMP (Symmetric Multiprocessing)?
> 2. What is SMT (Simultaneous Multithreading) and what problem does it solve?
> 3. What is the difference between a P-core and an E-core in Intel's hybrid architecture?

---

## 3. Cache Coherence in Multicore Systems

The most complex hardware challenge of multicore: if each core has its own L1 cache, and both cores can read/write the same memory location, the caches can have inconsistent values.

```
Cache coherence problem:
  
  Step 1: Both cores read address 0x1000 → both caches have value = 5
  
  Core 0 cache:   [0x1000: 5]
  Core 1 cache:   [0x1000: 5]
  DRAM:           [0x1000: 5]
  
  Step 2: Core 0 writes 0x1000 = 10 (into its cache, write-back policy)
  
  Core 0 cache:   [0x1000: 10]   ← updated
  Core 1 cache:   [0x1000: 5]    ← STALE!
  DRAM:           [0x1000: 5]    ← not yet updated
  
  Step 3: Core 1 reads 0x1000 → gets 5 (WRONG!)
```

**MESI protocol** (Chapter 29 introduced this): The most common cache coherence protocol. Each cache line has one of four states:
- **M (Modified)**: Line is dirty, only in this cache, no other caches have it
- **E (Exclusive)**: Line is clean, only in this cache, matches DRAM
- **S (Shared)**: Line is clean, possibly in multiple caches, matches DRAM
- **I (Invalid)**: Line is not valid (miss, or invalidated by another core)

```
MESI state transitions for a write:
  
  Core 0 has line in S state (shared with Core 1)
  Core 0 writes to the line:
    1. Core 0 sends RFO (Request For Ownership) on coherence bus
    2. Core 1 receives RFO: invalidates its copy (S → I)
    3. Core 0 transitions: S → M (now the only copy, modified)
    4. Core 0 completes write to its own cache
  
  Now Core 1 reads the line:
    5. Core 1 requests: sends ReadReq on coherence bus
    6. Core 0's cache (state M) must supply data (cache-to-cache transfer)
    7. Core 0: M → S (after supplying data)
    8. Core 1: I → S
    9. Both caches now have the updated value
```

**Directory protocol**: In large NUMA systems (many nodes), a broadcast snooping bus becomes a bottleneck. Directory protocols maintain a directory that tracks which caches have each line — enables point-to-point coherence messages.

**False sharing**: Two different variables on the same cache line (64 bytes). Different threads write to different variables, but both writing to the same cache line causes constant coherence traffic even though they aren't actually sharing data.

```c
// False sharing example
struct { int a; int b; } data;  // a and b on same cache line

Thread 0: while(1) data.a++;   // constantly invalidates cache line
Thread 1: while(1) data.b++;   // both threads fight over the same cache line
// Performance drops to worse than single-threaded!

// Fix: pad to separate cache lines
struct { int a; char pad[60]; int b; } data;
```

### Quick Check
> 1. Describe the cache coherence problem in a multicore system.
> 2. What is a "Request For Ownership" (RFO) in the MESI protocol?
> 3. What is false sharing and why is it bad for performance?

---

## 4. Memory Consistency Models

Cache coherence ensures a single memory location is consistent. Memory consistency models define what ordering of memory operations is guaranteed across different locations.

**Sequential consistency (SC)**: The most intuitive model. All memory operations appear to execute in program order, and all cores see the same interleaving.
- Strong but slow: every memory operation must be visible to all cores before proceeding
- x86 uses TSO (Total Store Ordering) — a relaxation of SC that allows store buffers

**TSO (Total Store Order)** — x86's model:
- Loads cannot pass loads (reads are ordered)
- Stores cannot pass stores (writes are ordered)
- A load CAN pass a prior store to a different address → loads can observe own stores before other cores do
- This allows hardware store buffers → higher performance

**ARM's Weak Memory Model (WMO):**
- Nearly all reorderings are allowed by default
- Must use explicit memory barrier instructions (`dmb`, `dsb`, `isb`) to enforce ordering
- More aggressive performance optimization possible; harder to program correctly

**RISC-V RVWMO (Weak Memory Ordering):**
- Similar to ARM's model — relaxed, with explicit fences
- `fence r,rw` = read → read/write fence
- Required for multi-core Linux on RISC-V

**Why memory models matter for programmers:**
```c
// Thread 0:           // Thread 1:
data = 42;            while (!flag);
flag = 1;             print(data);  // Prints 42? Or garbage?

// On x86 (TSO): flag=1 cannot be seen by Thread 1 before data=42
//   (stores are ordered in TSO)
// On ARM/RISC-V: this is broken — need memory barriers:

data = 42;            while (!flag);
__sync_synchronize();  __sync_synchronize();
flag = 1;             print(data);  // Now correct on all architectures
```

### Quick Check
> 1. What is sequential consistency and why is it too slow for modern hardware?
> 2. What is TSO (Total Store Ordering) and how does it differ from SC?
> 3. Why do ARM and RISC-V programs need explicit memory barriers that x86 programs don't?

---

## 5. NUMA — Non-Uniform Memory Access

As core count grows beyond what one DRAM controller can serve, chips move to NUMA:

**NUMA architecture:**
Each processor has its own local DRAM. Access to local DRAM is fast; access to another processor's DRAM (remote access) is slower and must traverse the inter-processor interconnect (UPI, Infinity Fabric, QPI).

```
NUMA example (dual-socket server):
  
  Socket 0:                     Socket 1:
  ┌────────────────────┐         ┌────────────────────┐
  │ Core 0..31         │         │ Core 32..63         │
  │ L1/L2/L3           │◄────────►│ L1/L2/L3           │
  │ Memory Ctrl        │  UPI    │ Memory Ctrl        │
  └──────┬─────────────┘ link    └──────┬─────────────┘
         │                              │
      64GB DDR5                      64GB DDR5
      (local)                        (remote for socket 0)
  
  Socket 0 accesses its own 64GB: ~50ns latency
  Socket 0 accesses Socket 1's 64GB: ~100–150ns (2–3× slower)
```

**NUMA effects on performance:**
- OS tries to allocate memory on the same socket as the thread accessing it (first-touch allocation)
- NUMA-unaware programs may accidentally allocate remote memory → 2× slower accesses
- Database systems, JVMs, HPC codes are NUMA-aware: explicitly control memory placement

**HBM in CPU packages**: Intel Sapphire Rapids (Xeon 4th Gen) offers HBM2e memory on-package alongside DDR5:
- HBM bandwidth: ~1 TB/s, low latency
- DDR5 bandwidth: ~350 GB/s
- Programs can use HBM as a "fast tier" — hot data in HBM, cold data in DDR5
- This is another form of NUMA within a single socket

### Quick Check
> 1. What is NUMA and what causes the access latency difference?
> 2. How does the OS manage memory allocation in a NUMA system?
> 3. What is HBM-on-package and how does it relate to NUMA?

---

## 6. Manycore Architectures — Hundreds of Cores

Beyond "multi" core (2–128 cores) is "many" core: hundreds or thousands of simpler cores.

**GPU as manycore** (Chapter 38 in depth): 16,896 CUDA cores in H100. Each core is simpler than a CPU core but together they deliver orders of magnitude more throughput for parallel workloads.

**Intel Xeon Phi (Knights Landing, 2016):**
- 72 x86-compatible cores on one chip
- Each core runs 4 hardware threads (SMT4)
- 288 effective threads
- AVX-512 SIMD on each core
- Targeted at HPC (high-performance computing) — ran standard x86 code, unlike GPUs
- Discontinued 2020 — could not compete with NVIDIA GPUs for ML workloads

**Network-on-Chip (NoC) for manycore:**
As cores increase beyond 8–16, a shared bus becomes a bottleneck. Manycore chips use mesh or ring network-on-chip:
```
Mesh NoC (4×4 example):
  
  C - C - C - C
  |   |   |   |
  C - C - C - C
  |   |   |   |
  C - C - C - C
  |   |   |   |
  C - C - C - C
  
  Each C = core + router
  Messages route hop-by-hop
  Bandwidth scales with core count
  Latency increases with distance (core 0 to core 15: 6 hops)
```

**Intel Loihi** (neuromorphic, Chapter 51): 128+ cores connected by an asynchronous mesh NoC — a manycore spike-routing network.

**RISC-V manycore research**: OpenPiton (Princeton) and Celerity (Michigan) are open-source manycore RISC-V designs — 500+ cores. Used to study NoC design, cache coherence protocols, and programming models.

### Quick Check
> 1. What is the difference between "multicore" and "manycore"?
> 2. What is a Network-on-Chip (NoC) and why is it needed for manycore?
> 3. Why did Intel Xeon Phi fail to compete with NVIDIA GPUs for ML workloads?

---

## Summary

- **Multicore** emerged because Dennard scaling failed — additional transistors became cores, not clock speed.
- **Amdahl's Law**: maximum speedup = 1/s (s = sequential fraction). 20% sequential code limits speedup to 5× regardless of core count.
- **SMP**: identical cores, shared memory, each with private L1/L2, shared L3. SMT adds logical cores per physical core.
- **Cache coherence** (MESI): hardware protocol ensuring all cores see consistent memory. RFO for write ownership, invalidation/update for sharers.
- **Memory models**: x86 (TSO) is strong; ARM/RISC-V (WMO) require explicit fence instructions for ordering.
- **NUMA**: multiple memory controllers in large systems; local access is 2–3× faster than remote.
- **Manycore**: hundreds to thousands of simpler cores connected by NoC; GPUs are the primary commercial example.

---

## Exercises

### Easy
1. What is Amdahl's Law and what does it say about the benefit of adding more cores?
2. What is false sharing and how can you fix it in C code?
3. What is NUMA and why does it affect performance?

### Medium
4. Amdahl's Law calculation: A matrix multiply program is 98% parallelizable. (a) Maximum speedup with 16 cores? (b) With 64 cores? (c) With 1024 cores? (d) At what core count is the speedup within 5% of the theoretical maximum? (e) The program is restructured to be 99.9% parallelizable: new maximum speedup and optimal core count?
5. Cache coherence sequence: Three cores (0, 1, 2) all read address X (initial value = 10). Then: Core 0 writes X = 20. Core 1 reads X. Core 2 writes X = 30. Core 0 reads X. Trace the MESI state of the cache line in each core after each operation, and identify the RFO messages sent.
6. Memory consistency: On ARM (relaxed memory model), two threads: Thread A writes `data = 99` then `flag = 1`. Thread B polls `while (!flag);` then reads `data`. (a) Without barriers: can Thread B see `data != 99`? Why? (b) Add appropriate `__atomic_thread_fence(...)` or `std::atomic` C++ barriers to make this correct. (c) On x86: is the original code (without barriers) correct? Explain using TSO's guarantee.

### Hard
7. Cache coherence protocol design: Design a simplified MESI protocol for a 4-core system with a snooping bus. For each transition: (a) Core 0 reads X (initially invalid). (b) Core 1 reads X. (c) Core 2 writes X. (d) Core 3 reads X. Draw the state machine and specify all bus messages (ReadReq, ReadResp, RFO, InvalidateReq, InvalidateAck, WritebackReq).
8. NUMA-aware memory allocator: You are optimizing a key-value store for a 2-socket NUMA server (Socket 0 and Socket 1, each with 32 cores and 64GB RAM). Queries are balanced across all 64 cores. (a) Naive allocation: all data on Socket 0's memory. Predict the performance bottleneck. (b) NUMA-aware: each socket owns half the keyspace. When Socket 1 workers need Socket 0 data: how does the access latency compare to local access? (c) Design a data structure that minimizes cross-socket accesses for a hash table. (d) Write pseudocode for a NUMA-aware malloc that allocates on the calling thread's local socket. (e) The key-value store serves 1M requests/second. 20% require cross-socket access at 100ns vs 50ns local. What is the throughput impact?
