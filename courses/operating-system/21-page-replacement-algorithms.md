# Chapter 21: Page Replacement Algorithms

> **"When RAM is full and a new page must be loaded, something must be evicted. Which page? The wrong choice cascades into thrashing. The right choice is invisible. This decision, made millions of times per day by the OS, is governed by page replacement algorithms — some elegant, some heuristic, all important."**

---

## Table of Contents

1. [The Page Replacement Problem](#1-the-page-replacement-problem)
2. [OPT — The Optimal Algorithm](#2-opt--the-optimal-algorithm)
3. [FIFO — First-In, First-Out](#3-fifo--first-in-first-out)
4. [LRU — Least Recently Used](#4-lru--least-recently-used)
5. [Clock Algorithm (Approximating LRU)](#5-clock-algorithm-approximating-lru)
6. [LFU — Least Frequently Used](#6-lfu--least-frequently-used)
7. [Working Set Model](#7-working-set-model)
8. [Linux's Page Replacement (LRU Approximation)](#8-linuxs-page-replacement-lru-approximation)
9. [Belady's Anomaly](#9-beladys-anomaly)
10. [Summary](#summary)

---

## 1. The Page Replacement Problem

When RAM is full and a page fault occurs, the OS must:
1. Find a victim page to evict
2. If dirty: write it to disk/swap
3. Load the new page into the freed frame
4. Update both page tables (old → evicted, new → loaded)

The question: **which page to evict?**

**The perfect choice:** Evict the page that won't be needed for the longest time in the future. But we don't know the future.

**The goal:** Minimize the number of page faults (each fault = potential disk I/O = performance hit).

**Reference string:**
To analyze algorithms, we use a sequence of page accesses:
```
Example reference string: 1, 2, 3, 4, 1, 2, 5, 1, 2, 3, 4, 5
(each number is a page number being accessed)
```

**Frame count:**
All algorithms perform better with more physical frames (more RAM).

---

## 2. OPT — The Optimal Algorithm

**Rule:** Evict the page that will not be used for the longest time in the future.

**Example:**
```
Frames = 3, Reference string: 1, 2, 3, 4, 1, 2, 5, 1, 2, 3, 4, 5

Access  Frames          Action
1       [1, -, -]       miss (compulsory)
2       [1, 2, -]       miss (compulsory)
3       [1, 2, 3]       miss (compulsory)
4       [1, 2, 4]       miss: evict 3 (next use of 3 at step 10, furthest away)
1       [1, 2, 4]       hit
2       [1, 2, 4]       hit
5       [1, 2, 5]       miss: evict 4 (4 used at step 11, further than 5's next use)
1       [1, 2, 5]       hit
2       [1, 2, 5]       hit
3       [3, 2, 5]       miss: evict 1 (1 not used after this; 5 used at step 12)
4       [3, 4, 5]       miss: evict 2 (2 not used after; 3 next at ∞, wait 5 next)
5       [3, 4, 5]       hit

Total page faults: 7 ← minimum possible with 3 frames
```

**Why OPT matters:**
OPT is the theoretical minimum (Oracle algorithm). We can't implement it (requires knowing the future), but we can benchmark other algorithms against it. How close does a real algorithm get to OPT's fault rate?

---

## 3. FIFO — First-In, First-Out

**Rule:** Evict the page that has been in memory the longest (first loaded = first evicted).

**Implementation:** Queue. New pages added to tail; victim is the head.

```
Frames = 3, Reference string: 1, 2, 3, 4, 1, 2, 5, 1, 2, 3, 4, 5

Access  Frames (queue front→back)  Action
1       [1, -, -]                  miss: add 1
2       [1, 2, -]                  miss: add 2
3       [1, 2, 3]                  miss: add 3
4       [2, 3, 4]                  miss: evict 1 (oldest)
1       [3, 4, 1]                  miss: evict 2 (oldest) — oops, 1 was needed!
2       [4, 1, 2]                  miss: evict 3 (oldest)
5       [1, 2, 5]                  miss: evict 4 (oldest)
1       [1, 2, 5]                  hit
2       [1, 2, 5]                  hit
3       [2, 5, 3]                  miss: evict 1 (oldest)
4       [5, 3, 4]                  miss: evict 2
5       [5, 3, 4]                  hit

Total page faults: 9  (vs OPT: 7)
```

**Problems with FIFO:**
- Doesn't consider usage frequency — evicts a heavily-used page if it's "old"
- Suffers from Belady's Anomaly (see section 9)

**FIFO is rarely used alone** in production systems.

---

## 4. LRU — Least Recently Used

**Rule:** Evict the page that was least recently used. Based on temporal locality: if a page hasn't been used recently, it probably won't be used soon.

**Example:**
```
Frames = 3, Reference string: 1, 2, 3, 4, 1, 2, 5, 1, 2, 3, 4, 5

Access  Frames (ordered by recency, most recent last)  Action
1       [1]             miss
2       [1, 2]          miss
3       [1, 2, 3]       miss
4       [2, 3, 4]       miss: evict 1 (LRU)
1       [3, 4, 1]       miss: evict 2 (LRU) — 2 is older than 3 and 4
2       [4, 1, 2]       miss: evict 3 (LRU)
5       [1, 2, 5]       miss: evict 4 (LRU)
1       [2, 5, 1]       hit: 1 becomes MRU
2       [5, 1, 2]       hit: 2 becomes MRU
3       [1, 2, 3]       miss: evict 5 (LRU)
4       [2, 3, 4]       miss: evict 1 (LRU)
5       [3, 4, 5]       miss: evict 2 (LRU)

Total page faults: 8  (vs OPT: 7, FIFO: 9)
```

**LRU is theoretically strong** because it's based on temporal locality — a proven characteristic of real programs.

**Implementation challenges:**
True LRU requires tracking the EXACT last-use time of every page. Two approaches:

**1. Timestamp counter:**
On every memory access, record the current time in the page's PTE (use extra bits). On eviction, scan all pages for the smallest timestamp.
- Problem: needs hardware support to update timestamps on EVERY memory access → enormous overhead

**2. Ordered stack:**
Keep all pages in a stack. On each access, move the accessed page to the top.
- Problem: requires updating a linked list on every memory access → too slow

**Neither is practical for hardware implementation.** Real systems use LRU APPROXIMATIONS.

---

## 5. Clock Algorithm (Approximating LRU)

The **Clock Algorithm** (also called the "Second Chance" or NRU — Not Recently Used algorithm) approximates LRU using only the **Accessed bit** in each PTE (which hardware sets automatically on any access).

**Algorithm:**
```
Maintain a circular list of all pages (like a clock hand pointing at each)

When a page fault occurs:
  Look at the page the clock hand points to:
  
  If Accessed bit == 0:
    This page hasn't been used since our last sweep
    → Evict this page
    → Advance clock hand
    
  If Accessed bit == 1:
    Give this page a "second chance" — clear the bit
    → Accessed bit = 0
    → Advance clock hand
    → Continue to next page

Keep going until we find a page with Accessed bit == 0
```

**Visual example:**
```
Clock hands point to pages in a circle: [A][B][C][D][E][F]
Each page has Accessed bit (0 or 1)

                A(1)
            F(1)    B(0)  ← clock hand here
          E(0)        C(1)
                D(0)

Page fault! Find victim starting at B:
B: Accessed=0 → EVICT B, load new page here, advance hand to C
```

**Why this works:**
- Pages accessed recently have Accessed=1 (set by hardware) → get a second chance
- Pages not accessed recently have Accessed=0 → candidates for eviction
- Exactly one bit per page → very low overhead
- Hardware sets the bit on every access → no software intervention per access

**Enhanced clock (two bits: Accessed + Dirty):**
```
(A=Accessed, D=Dirty)

(0,0): Best victim — not recently used, not dirty (no disk write needed)
(0,1): Not recently used but dirty (need to write to disk — worse than 0,0)
(1,0): Recently used, not dirty (OK to evict if forced)
(1,1): Worst victim — recently used AND dirty (need disk write)

On each pass, demote (1,x) to (0,x) — give second chance
Prefer (0,0) → (0,1) → second pass (1,0) → (1,1)
```

**The enhanced clock is what Linux uses as the basis for its page replacement.**

---

## 6. LFU — Least Frequently Used

**Rule:** Evict the page with the lowest access count.

**Implementation:** Each page has a counter, incremented on each access. Evict the page with count = 0 (or minimum).

**Problem:** A page heavily used during startup has a high count. Even if never used again, it stays in memory forever because of its historical high count.

**Fix:** Aging — periodically right-shift all counters. Recent accesses count more than old ones.

LFU is rarely used for OS page replacement. It works better for application-level caching (CDN edge caches, web caches).

---

## 7. Working Set Model

Peter Denning's **Working Set Model** (1968) defines the set of pages a process NEEDS at any moment.

**Working set W(t, Δ):**
The set of pages referenced by a process in the last Δ time units, at time t.

```
Page reference history: ...1, 2, 1, 3, 4, 3, 2, 1...
                                        ↑ current time t
Working set W(t, 5) = pages referenced in last 5 accesses = {3, 4, 2, 1}
```

**Key insight:**
A process should be in memory only if ALL its working set pages fit in RAM.

**Working set policy:**
- If a process's working set fits in RAM: let it run
- If not: swap the entire process out (don't let it thrash)

**The Working Set Model explains why systems thrash:**
Thrashing happens when: Σ(working set sizes) > physical RAM

**Practical implementation:**
Tracking exact working sets is expensive. Linux approximates with:
- LRU list aging
- Page reference bits
- Memory pressure watermarks

---

## 8. Linux's Page Replacement (LRU Approximation)

Linux maintains two LRU lists per memory zone:

```
Active LRU list:   [recently accessed pages]    ← hot pages
Inactive LRU list: [less recently accessed pages] ← cold pages (candidates)

Pages flow:
New page → Inactive list (tail)
Accessed again → promoted to Active list
Aging → Active page → Inactive list (after time)
Eviction → takes from Inactive list head
```

**kswapd — the kernel swap daemon:**
Linux runs a kernel thread `kswapd` that wakes up when free pages drop below a threshold:
```
free_pages < low watermark:
  kswapd wakes up
  Scans inactive LRU list
  Evicts pages from tail of inactive list
  Writes dirty pages to swap
  Runs until free_pages > high watermark
  Goes back to sleep
```

**Direct reclaim:**
If memory pressure is severe enough that kswapd can't keep up, the process that triggered the fault directly reclaims memory before proceeding.

**Memory watermarks:**
```
/proc/sys/vm/min_free_kbytes  — minimum free pages (hard floor)
high watermark = min + 2×zone_pages/256
low watermark  = min + zone_pages/256
```

**Transparent Huge Pages (THP):**
Linux can automatically promote groups of 512 adjacent 4KB pages into one 2MB huge page — reduces TLB pressure and page table size. Configurable:
```bash
cat /sys/kernel/mm/transparent_hugepage/enabled
# [always] madvise never
echo madvise > /sys/kernel/mm/transparent_hugepage/enabled
```

---

## 9. Belady's Anomaly

A counterintuitive result: adding MORE frames can cause MORE page faults with FIFO.

**Example with FIFO:**
```
Reference string: 1, 2, 3, 4, 1, 2, 5, 1, 2, 3, 4, 5

With 3 frames: 9 page faults
With 4 frames: 10 page faults  ← MORE FRAMES, MORE FAULTS!
```

**Why?**
With 3 frames: page 5 is loaded at step 7, evicting page 4.
With 4 frames: page 5 is loaded at step 7, evicting a different page — causing faults on 3 and 4 later.

**FIFO is susceptible to Belady's Anomaly. LRU is not.**

This is a key reason LRU (and its approximations) are preferred over FIFO.

**Stack algorithms (like LRU, OPT):**
A page replacement algorithm is a "stack algorithm" if the set of pages in memory with n frames is always a subset of the set with n+1 frames. Stack algorithms never exhibit Belady's Anomaly.

LRU is a stack algorithm. FIFO is not.

---

## Summary

| Algorithm | Eviction policy | Fault rate | Belady? | Real use |
|-----------|----------------|------------|---------|---------|
| OPT | Furthest future use | Best possible | No | Benchmark only |
| FIFO | Oldest in memory | Poor | Yes | Rarely used alone |
| LRU | Least recently used | Near-OPT | No | Approximations used |
| Clock/NRU | Accessed bit clear | Good approximation | No | Linux, macOS |
| LFU | Least access count | Poor (aging helps) | No | Application caching |
| Working Set | Keep active working set | Good | N/A | Theory; approximated |

**Linux's real approach:**
- Active/Inactive LRU lists
- Clock scan of inactive list for victims
- kswapd daemon manages reclaim in background
- Direct reclaim under pressure
- Per-zone management with watermarks

**The bottom line:** LRU (or its approximations) is the right choice for general-purpose page replacement. It leverages temporal locality — the principle that recently used data is likely to be used again — which holds for virtually all real programs.
