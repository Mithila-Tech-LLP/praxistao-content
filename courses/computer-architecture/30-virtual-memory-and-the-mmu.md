# Chapter 30: Virtual Memory and the MMU

Every program thinks it has the entire memory address space to itself. A browser can be at address 0x400000, and a music player can also be at address 0x400000 — at the same time, on the same machine. How is this possible? Virtual memory: an abstraction that gives each process its own private address space, hides the physical layout of RAM, and enables features like memory-mapped files, demand paging, and copy-on-write that modern operating systems depend on.

## Table of Contents

1. [Why Virtual Memory Exists](#1-why-virtual-memory-exists)
2. [The Address Translation Problem](#2-the-address-translation-problem)
3. [Page Tables — Mapping Virtual to Physical](#3-page-tables--mapping-virtual-to-physical)
4. [The TLB — Speeding Up Translation](#4-the-tlb--speeding-up-translation)
5. [Page Faults and Demand Paging](#5-page-faults-and-demand-paging)
6. [The MMU and OS Cooperation](#6-the-mmu-and-os-cooperation)
7. [Large Pages and HUGE TLB](#7-large-pages-and-huge-tlb-coverage)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Virtual Memory Exists

Consider a world without virtual memory, where programs use physical addresses directly:

**Problem 1: Each process needs to know where in RAM it will be loaded.** You'd have to compile a program for a specific memory address and it would break if another program was already there.

**Problem 2: Programs could read each other's memory.** A buggy program could scribble over the OS kernel. A malicious program could steal another process's data.

**Problem 3: Programs couldn't use more memory than physically exists.** No overcommit, no demand paging, no swap.

Virtual memory solves all three:
- Every process has a **virtual address space** (e.g., 0–2^48 on 64-bit Linux) independent of physical location
- The OS maps each process's virtual pages to physical pages, maintaining isolation
- Physical pages can be swapped to disk, enabling more total memory use than DRAM capacity

```
Process A (browser):    Virtual 0x401000 → Physical 0x2003000
Process B (music):      Virtual 0x401000 → Physical 0x7F10000
                        Same virtual address → different physical addresses
```

### Quick Check
> 1. Why can't two programs both be loaded at virtual address 0x400000 on a system without virtual memory?
> 2. List three things virtual memory provides beyond just "translation."
> 3. A process has a 256GB virtual address space but the machine only has 16GB of RAM. How is this possible?

---

## 2. The Address Translation Problem

Every memory access in a program uses a **virtual address**. Before accessing DRAM, this virtual address must be **translated** to a **physical address**.

The translation must be:
- **Fast**: Translation happens on every memory access (billions per second)
- **Protected**: Only the OS should be able to change mappings
- **Flexible**: Different virtual regions can map to different physical locations, with different permissions (read/write/execute)

The unit of mapping is a **page** — typically 4KB. Translation works at page granularity:
```
Virtual address:  [Virtual Page Number (VPN)][Page Offset]
                   upper bits                  lower 12 bits

Physical address: [Physical Frame Number (PFN)][Page Offset]
                   translated by MMU            unchanged
```

Only the VPN is translated. The page offset (lower 12 bits for 4KB pages) is passed through unchanged — within a page, the relative offset stays the same.

```
4KB page example:
  Virtual address:  0x00401A34
  Page size = 4KB = 0x1000
  VPN = 0x00401A34 >> 12 = 0x401      (virtual page 0x401)
  Offset = 0x00401A34 & 0xFFF = 0xA34

  MMU looks up VPN 0x401 → PFN = 0x2003 (say)

  Physical address: 0x2003 << 12 | 0xA34 = 0x2003A34
```

### Quick Check
> 1. What are the two parts of a virtual address?
> 2. Why is the page offset not translated?
> 3. With 4KB pages and a 48-bit virtual address, how many bits are the VPN? The offset?

---

## 3. Page Tables — Mapping Virtual to Physical

The page table is the OS-maintained data structure that maps VPNs to PFNs. There is one page table per process.

### Single-Level Page Table

The simplest design: an array indexed by VPN.

```
Page Table for Process A:
Index (VPN) | PFN    | Valid | R | W | X
0x000       | 0x100  |  1    | 1 | 1 | 0   (data page)
0x001       | 0x201  |  1    | 1 | 0 | 1   (code page, read+execute)
...
0x401       | 0x2003 |  1    | 1 | 1 | 0
...
0xFFF...    | -      |  0    | -   (not mapped — accessing causes page fault)
```

**Problem**: With 48-bit virtual addresses and 4KB pages, there are 2^36 ≈ 64 billion possible VPNs. A full array would require 64 billion × 8 bytes = 512GB of memory just for the page table — larger than the physical memory! This is impractical.

### Multi-Level Page Tables

The solution: make the page table **hierarchical**. Only allocate page table memory for virtual regions that are actually mapped.

**x86-64 uses a 4-level page table** (48-bit virtual addresses):

```
Virtual address: [L4 index (9b)][L3 index (9b)][L2 index (9b)][L1 index (9b)][Offset (12b)]

CR3 register → L4 table (512 entries × 8 bytes = 4KB)
                  ↓ (indexed by L4 bits)
              L3 table (only allocated for mapped regions)
                  ↓ (indexed by L3 bits)
              L2 table
                  ↓ (indexed by L2 bits)
              L1 table (page table entries = PTE)
                  ↓ (indexed by L1 bits)
              Physical frame
```

For a process using 10MB of code/data: instead of a 512GB flat page table, you allocate only the L4 table (always present) + a few L3/L2/L1 tables for the mapped region. Total: maybe 20KB of page table memory.

**Page Table Entry (PTE)** bits:
- `P` (Present): is this page in DRAM? (0 = page fault)
- `R/W` (Read/Write): can this page be written?
- `U/S` (User/Supervisor): accessible in user mode?
- `X` (Execute Disable, XD/NX): can code be executed from this page?
- `Dirty`: has this page been written since last time?
- `Accessed`: has this page been accessed recently? (used by OS page replacement)
- `PFN`: the physical frame number (upper bits)

### Quick Check
> 1. Why is a single-level page table impractical for 48-bit virtual addresses?
> 2. x86-64 uses 4 levels of page tables with 9 bits per level. How many table entries does each level have? How large is each table?
> 3. What does the "Present" bit do in a page table entry?

---

## 4. The TLB — Speeding Up Translation

A multi-level page table walk requires 4 memory accesses before we even get to the data! That would mean every load/store takes 5× as long — catastrophic.

The **TLB (Translation Lookaside Buffer)** is a small, fully-associative cache that stores recently used VPN→PFN mappings. Since programs exhibit **temporal and spatial locality in their page accesses**, most translations repeat.

```
TLB lookup (every memory access):
  1. Extract VPN from virtual address
  2. Check TLB (16-1024 entries, fully-associative or set-associative)
  3. TLB HIT (>99% of the time): get PFN directly — access L1 cache in same cycle
  4. TLB MISS (~0.1% of accesses): walk page tables (hardware or software)
                                    → 4 memory accesses for x86-64
                                    → store result in TLB
                                    → retry access
```

**TLB miss cost**: 4 extra memory accesses = potentially 4 × 200 cycles = 800 cycles if each page table level misses in L1 cache! In practice, page table entries are frequently accessed and cached in L1/L2, so TLB miss cost is typically 10–50 cycles.

**TLB size**: L1 TLB typically has 64–128 entries; L2 TLB has 1024–2048 entries. With 4KB pages, 64 TLB entries covers only 64 × 4KB = 256KB of address space — fine for hot code/data but insufficient for large working sets. This is why **huge pages** matter (Section 7).

**Context switch and TLB flush**: When the OS switches between processes, all TLB entries from the old process are invalid (they map to that process's virtual addresses). The OS can flush the entire TLB (simple, expensive) or use **ASID (Address Space ID)** tags on each TLB entry to avoid flushing (modern approach: x86 PCID, ARM ASID).

```
TLB structure (simplified, with ASID):
  ASID | VPN    | PFN    | Dirty | Valid | Permissions
  0x05 | 0x401  | 0x2003 |  0    |  1    | R/W
  0x07 | 0x401  | 0x9F00 |  1    |  1    | R/W/X
  ...
```

Multiple processes can have entries in the TLB simultaneously because the ASID distinguishes them.

### Quick Check
> 1. What problem does the TLB solve?
> 2. With 64 TLB entries and 4KB pages, what is the maximum address range covered by a full TLB?
> 3. Why must the TLB be flushed (or use ASID tags) on a context switch?

---

## 5. Page Faults and Demand Paging

When the CPU accesses a virtual address and the TLB misses, it walks the page table. If it finds a PTE with Present = 0, it generates a **page fault** — a trap to the OS kernel.

The OS page fault handler:
1. Determines why the page is not present
2. Takes appropriate action
3. Sets Present = 1 in the PTE
4. Returns to user code, which retries the faulting instruction

### Types of Page Faults

**Demand paging**: When a program starts, its pages are not loaded into RAM. The OS only allocates physical pages when they are first accessed. The first access to any page triggers a page fault; the OS reads the page from the executable file and maps it. This makes program startup faster — only actually-used pages consume RAM.

**Stack growth**: When the stack grows to a new page (by calling deeply nested functions), the first write to the new stack page triggers a page fault. The OS allocates a new physical page and maps it.

**Copy-on-write (COW)**: When `fork()` creates a child process, instead of copying all the parent's memory, the OS maps the same physical pages into both processes as read-only. The first write to any shared page triggers a page fault; the OS then copies just that one page and maps it as writable in the faulting process. Elegant and efficient.

**Swap (page-out/page-in)**: When DRAM is full, the OS evicts cold pages to disk (swap space), freeing physical memory. Accessing an evicted page triggers a page fault; the OS reads it back from disk. This allows total virtual memory usage across all processes to exceed physical RAM — at the cost of disk access latency (milliseconds vs nanoseconds).

**Segfault**: If a process accesses an address that is not mapped at all (e.g., dereferencing a NULL pointer), the OS sends `SIGSEGV` — segmentation fault. The process terminates.

```
Page fault handling time:
  OS page fault handler overhead: ~1 µs
  Reading page from NVMe SSD: ~100 µs
  Reading page from HDD: ~10 ms
  Total for swap page-in from SSD: ~101 µs = 303,000 CPU cycles!
```

This is why thrashing (excessive swapping) is catastrophically bad for performance.

### Quick Check
> 1. What is "demand paging" and why does it speed up program startup?
> 2. Explain copy-on-write (COW) fork in one paragraph. When exactly does a physical copy happen?
> 3. A process dereferences a NULL pointer (address 0). What happens, step by step?

---

## 6. The MMU and OS Cooperation

The **MMU (Memory Management Unit)** is the hardware that performs address translation. In modern CPUs, the MMU is integrated into the processor die itself. The OS and MMU work together:

**OS responsibilities:**
- Maintain page tables for each process
- Handle page faults
- Map/unmap memory regions (via `mmap()`, `brk()`, etc.)
- Manage swap

**MMU responsibilities:**
- On every memory access: look up TLB, walk page tables on TLB miss
- Enforce permission bits (no write to read-only pages, no execute from data pages, no kernel access from user mode)
- Generate page fault exception when access is illegal

**Privilege levels**: The MMU enforces the distinction between **kernel mode (ring 0)** and **user mode (ring 3)**. User-mode code cannot access pages marked supervisor-only (kernel memory). Attempting to do so triggers a general protection fault and the process is terminated (or a kernel panic if it's the kernel itself misbehaving).

**Physical Address Extension (PAE) and 5-level paging**: x86 was originally 32-bit, limiting physical memory to 4GB. PAE extended physical addresses to 36 bits (64GB). 64-bit x86 uses 48-bit virtual addresses (256TB) and 52-bit physical addresses. Intel Tiger Lake added 5-level paging for 57-bit virtual addresses (128PB) — relevant for huge in-memory databases.

### Quick Check
> 1. What is the hardware component that performs virtual-to-physical address translation on every access?
> 2. What happens when user-mode code tries to access a kernel memory page?
> 3. Why was PAE (Physical Address Extension) needed for 32-bit x86 systems?

---

## 7. Large Pages and Huge TLB Coverage

With 4KB pages, 64 TLB entries cover only 256KB. A process with 1GB of active data needs 262,144 TLB entries to map it all — far more than any real TLB. This causes **TLB thrashing**: the working set's page mappings don't fit in the TLB, causing frequent TLB misses (page table walks) that dominate performance.

**Solution: Large pages** (also called huge pages):
- Standard: 4KB pages
- Large (x86): **2MB pages** — single TLB entry covers 2MB instead of 4KB (512× improvement)
- Huge (x86): **1GB pages** — single TLB entry covers 1GB

With 2MB pages, 64 TLB entries cover 64 × 2MB = 128MB. For a 1GB working set, only 512 TLB entries are needed.

**Tradeoffs of large pages:**
- **Pro**: Massive reduction in TLB misses for large working sets (databases, JVMs, HPC)
- **Con**: Internal fragmentation — if a process allocates a 2MB page but only uses 1 byte, 2MB-1B is wasted
- **Con**: Allocation pressure — contiguous 2MB aligned physical memory blocks can be scarce

**Transparent Huge Pages (THP)**: Linux can automatically promote 4KB pages to 2MB pages in the background for anonymously mapped memory. Applications using `mmap()` or `malloc()` can benefit without code changes.

**Explicit huge pages**: Databases (PostgreSQL, Oracle) and JVMs explicitly allocate huge pages for their buffer pools. PostgreSQL can be configured with `huge_pages = on` to map the shared buffer using 2MB pages.

```
TLB coverage comparison:
  4KB pages, 64 TLB entries:   256KB covered
  2MB pages, 64 TLB entries:   128MB covered   (512× better)
  1GB pages, 64 TLB entries:   64GB covered    (262,144× better)
```

### Quick Check
> 1. Calculate how many 4KB-page TLB entries you'd need to cover 1GB of working set data.
> 2. What is "internal fragmentation" with huge pages?
> 3. Why do databases benefit so much from huge pages?

---

## Summary

- **Virtual memory** gives each process an isolated, private address space. The OS maps virtual pages to physical frames.
- Addresses are divided into **Virtual Page Number (VPN)** (translated) and **page offset** (unchanged). Standard page size is 4KB.
- **Multi-level page tables** (4 levels on x86-64) avoid allocating huge flat arrays for sparse virtual address spaces.
- The **TLB** caches recently used VPN→PFN mappings. A TLB hit takes 0 extra cycles; a miss triggers a page table walk (10–50+ cycles). **ASID** tags avoid TLB flushes on context switches.
- **Page faults** occur when a page is not present (demand paging, swap, COW) or when access is illegal (segfault). The OS handles them.
- The **MMU** enforces permissions (read/write/execute, user/kernel) on every access.
- **Large/huge pages** (2MB or 1GB) dramatically improve TLB coverage, critical for databases and large in-memory workloads.

---

## Exercises

### Easy
1. A process accesses virtual address 0x00402B10. The page size is 4KB (0x1000). What is the VPN? What is the page offset?
2. What is the TLB and why is it needed alongside page tables?
3. List four different situations that can cause a page fault and describe what the OS does for each.

### Medium
4. A 4-level x86-64 page table uses 9 bits per level. A process maps virtual addresses 0x0000000000400000 through 0x0000000000600000 (2MB). (a) How many L1 page tables are needed? (b) How many L2, L3, L4 entries are needed? (c) What is the total memory used for this process's page tables?
5. A process accesses addresses that form a 2GB working set, uniformly distributed. The TLB has 128 entries. (a) With 4KB pages, what fraction of working set pages can the TLB cover? (b) With 2MB huge pages, what fraction? (c) Assuming each TLB miss costs 40 cycles and memory accesses are every 10 instructions, what is the CPI contribution from TLB misses in each case?
6. Explain copy-on-write fork in terms of page table entries. After fork(), Process A and Process B share the same physical pages. Process B writes to virtual address 0x401000. Trace the exact steps from the write instruction to the point where B has its own private copy.

### Hard
7. A database server runs on a machine with 64GB RAM. It allocates a 32GB buffer pool. With 4KB pages and a 128-entry TLB, estimate the TLB miss rate for random buffer pool accesses (assume uniform access pattern). Then estimate the improvement if the buffer pool uses 2MB huge pages. Calculate the performance impact assuming memory accesses constitute 40% of all instructions and each TLB miss costs 50 cycles.
8. Describe the "Meltdown" vulnerability (a different side-channel attack than Spectre): the kernel maps all physical memory into kernel virtual address space. Before Meltdown mitigations, the kernel's page table entries were present in user-mode process page tables (but marked supervisor-only). How did Meltdown exploit speculative execution to read kernel memory from user mode? What is the Kernel Page Table Isolation (KPTI) mitigation, and why does it have a performance cost?
