# Chapter 18: Virtual Memory and Paging

> **"Virtual memory is the great illusion. Every program thinks it owns the entire machine. Each process believes it has gigabytes of private memory all to itself. The reality is more complex — and the elegant mechanism that maintains this illusion, while safely sharing physical RAM between all processes, is one of the greatest achievements in computer engineering."**

---

## Table of Contents

1. [The Problem Virtual Memory Solves](#1-the-problem-virtual-memory-solves)
2. [The Big Idea — Virtual Address Space](#2-the-big-idea--virtual-address-space)
3. [Pages and Frames — The Mapping Unit](#3-pages-and-frames--the-mapping-unit)
4. [The Page Table](#4-the-page-table)
5. [How Address Translation Works](#5-how-address-translation-works)
6. [The MMU — Hardware Address Translator](#6-the-mmu--hardware-address-translator)
7. [Page Table Entry — Bits and Flags](#7-page-table-entry--bits-and-flags)
8. [The TLB — Caching Translations](#8-the-tlb--caching-translations)
9. [What Virtual Memory Enables](#9-what-virtual-memory-enables)
10. [Virtual Memory on Real Hardware (x86-64)](#10-virtual-memory-on-real-hardware-x86-64)
11. [Setting Up Paging — A Conceptual Walk-Through](#11-setting-up-paging--a-conceptual-walk-through)
12. [Summary](#summary)

---

## 1. The Problem Virtual Memory Solves

Without virtual memory, there are three major problems:

**Problem 1: Address conflicts**
If Program A runs at address 0x1000 and Program B also wants to run at address 0x1000, they collide.

**Problem 2: Not enough physical RAM**
If you have 8GB RAM and want to run programs that together need 12GB, you're stuck.

**Problem 3: Security**
Without isolation, Program A can read and modify Program B's memory — steal passwords, corrupt data.

**Virtual memory solves all three simultaneously.**

---

## 2. The Big Idea — Virtual Address Space

**Every process gets its own "virtual" address space** — a private view of memory that looks like the entire address space is dedicated to it.

```
Process A's view:          Process B's view:
0x0000 → [A's code]        0x0000 → [B's code]
0x1000 → [A's data]        0x1000 → [B's data]
0x2000 → [A's stack]       0x2000 → [B's stack]

Physical RAM:
0x5000 → [A's code page]   0x9000 → [B's code page]
0x6000 → [A's data page]   0xA000 → [B's data page]
0x7000 → [A's stack page]  0xB000 → [B's stack page]
```

Both A and B "think" they're using addresses 0x0000–0x3000. But they're actually using completely different physical memory. The OS maintains a separate translation table for each process.

When Process A writes to 0x1000, the hardware translates it to physical 0x6000. When Process B writes to 0x1000, the hardware translates it to physical 0xA000. They NEVER interfere.

---

## 3. Pages and Frames — The Mapping Unit

Virtual memory doesn't work byte-by-byte. The address space is divided into fixed-size chunks:

**Page (virtual):** A fixed-size chunk of the virtual address space. Typically 4096 bytes (4KB).

**Frame (physical):** A fixed-size chunk of physical RAM. Same size as a page.

The OS maps pages to frames. Each virtual page maps to one physical frame (or to nothing — unmapped).

```
Virtual address space:      Physical RAM:
┌─────────────┐             ┌─────────────┐
│  Page 0     │─────────►  │  Frame 5    │
├─────────────┤             ├─────────────┤
│  Page 1     │──────────►  │  Frame 2    │
├─────────────┤             ├─────────────┤
│  Page 2     │ (unmapped)  │  Frame 8    │
├─────────────┤             ├─────────────┤
│  Page 3     │────────────►│  Frame 0    │
└─────────────┘             └─────────────┘
```

Key facts:
- Virtual pages can map to ANY physical frame (not necessarily in order)
- Not all virtual pages need to be mapped (gaps are valid — accessing them causes a page fault)
- Multiple processes can't map the same frame (unless explicitly sharing)
- The OS + hardware together enforce this

---

## 4. The Page Table

A **page table** is a data structure maintained by the OS for each process. It maps virtual page numbers to physical frame numbers.

**Simple (conceptual) page table:**
```
Process A's page table:
Virtual Page 0 → Physical Frame 5  (present, read-write)
Virtual Page 1 → Physical Frame 2  (present, read-only)
Virtual Page 2 → NOT PRESENT       (not mapped)
Virtual Page 3 → Physical Frame 0  (present, read-write)
...
```

Each entry in the page table is a **Page Table Entry (PTE)**.

**Page table location:**
The page table itself is stored in physical memory. The CPU has a special register (CR3 on x86) that holds the physical address of the current process's page table.

When the OS switches processes (context switch), it:
1. Saves the old process's register state
2. Loads the new process's CR3 register (points to new process's page table)
3. Loads the new process's registers
→ Now ALL virtual address translations use the new process's page table. Instant switch between address spaces!

---

## 5. How Address Translation Works

When a program accesses virtual address V:

```
Virtual Address: 0x00403456

Step 1: Split into page number and offset
  Page size = 4096 = 2^12
  Offset = lower 12 bits = 0x456
  Virtual page number (VPN) = upper bits = 0x403

Step 2: Look up page table
  Index page_table[0x403]
  → PTE says: physical frame number = 0x789, flags: present, readable, writable

Step 3: Construct physical address
  Physical address = (frame number << 12) | offset
                   = (0x789 << 12) | 0x456
                   = 0x789456

Step 4: Access physical RAM at 0x789456
```

This translation happens for EVERY memory access — every instruction fetch, every data read, every data write. The hardware does it automatically (the MMU).

---

## 6. The MMU — Hardware Address Translator

The **MMU (Memory Management Unit)** is hardware inside the CPU that automatically translates virtual to physical addresses.

**Without MMU:** Software would need to translate every address → impossibly slow.

**With MMU:** Hardware translation is transparent, adds ~1ns overhead (especially with TLB caching).

```
CPU core
┌─────────────────────────────────────────────┐
│                                              │
│  ┌──────────────┐    Virtual   ┌──────────┐ │
│  │  Execution   │── address ──►│   MMU    │ │
│  │  Units       │              │          │ │
│  └──────────────┘              │ Translates│ │
│                                │ V → P    │ │
│                                └────┬─────┘ │
│                                     │       │
└─────────────────────────────────────┼───────┘
                                      │ Physical address
                                      ↓
                                   RAM / Cache
```

**What the MMU checks for each access:**
1. Is the page mapped? (present bit in PTE)
   - If not: page fault exception → OS handles it
2. Does the access type match page permissions?
   - Write to read-only page? → protection fault → SIGSEGV to process
   - User code accessing kernel page? → protection fault → SIGSEGV
3. Translate virtual → physical
4. Issue the physical memory access

---

## 7. Page Table Entry — Bits and Flags

A page table entry (PTE) isn't just a physical address. It contains flags:

**x86-64 PTE flags:**
```
Bit 63    : XD (Execute-Disable) — 1 = can't execute code from this page
Bits 51-12: Physical Frame Number (12 bits for offset → 40 bits for frame number)
Bit 11-9  : Available (OS can use for its own purposes)
Bit 8     : Global — don't flush TLB on CR3 change (for kernel pages)
Bit 7     : Page Size — for huge pages (1 = 2MB page at this level)
Bit 6     : Dirty — CPU sets this when page is written to
Bit 5     : Accessed — CPU sets this when page is read or written
Bit 4     : PCD (Page Cache Disable) — don't cache this page
Bit 3     : PWT (Page Write-Through) — use write-through caching
Bit 2     : User/Supervisor — 0 = kernel only, 1 = user can access
Bit 1     : Read/Write — 0 = read-only, 1 = read-write
Bit 0     : Present — 1 = page is in RAM, 0 = not present (unmapped or swapped)
```

**Key flags:**

**Present (P) bit:**
- 1: page is in physical RAM — translation proceeds
- 0: page is NOT in RAM — page fault! OS must handle it

**Read/Write (RW) bit:**
- 0: read-only (code pages, shared data, copy-on-write pages)
- 1: writable

**User/Supervisor (US) bit:**
- 0: only kernel can access (ring 0)
- 1: user can access (ring 3)

**Dirty (D) bit:**
Hardware sets this when the page is written. OS uses it for swap: if dirty, page must be written to disk before evicting; if clean, can just discard (can be re-read from file).

**Accessed (A) bit:**
Hardware sets this when the page is read or written. OS uses it for page replacement: frequently accessed pages should stay in RAM.

**NX (Execute-Disable) bit:**
1 = code cannot be executed from this page. Used for security:
- Stack is NX: programs can't execute code on the stack
- Heap is NX: programs can't execute code they malloc'd
- Prevents many exploits that inject code into data

---

## 8. The TLB — Caching Translations

For every memory access, the CPU needs to walk the page table — which is itself in memory! That would mean 2 memory accesses per user memory access (one to fetch PTE, one for actual data).

The **TLB (Translation Lookaside Buffer)** is a small, ultra-fast cache of recent virtual → physical translations.

```
Virtual Address → ┌──────────────┐
                  │  TLB lookup  │ → TLB HIT: physical address instantly (O(1), ~1 cycle)
                  └──────────────┘ → TLB MISS: walk page table, cache result (~20-50 cycles)
                  (typically 32-1024 entries, fully associative)
```

**TLB hit rate is critical:**
Modern programs achieve ~99% TLB hit rate because they have good locality (accessing the same pages repeatedly). This makes virtual memory practical.

**TLB flush:**
When the OS switches to a different process (different page table), the TLB entries from the old process are invalid for the new process. The OS must flush (invalidate) the TLB.

On x86: `mov cr3, rax` flushes the TLB. This is expensive! (All future accesses are TLB misses until the TLB warms up.)

**ASID (Address Space Identifier):**
Modern CPUs tag TLB entries with an ASID, allowing TLB entries from multiple processes to coexist. This avoids full flushes on context switch — only entries for the new ASID are used.

Linux uses ASIDs on x86-64 (process context identifiers/PCIDs).

---

## 9. What Virtual Memory Enables

Virtual memory is the foundation for many OS features:

**1. Process isolation:**
Process A and B have separate page tables → separate virtual address spaces → can't access each other's memory.

**2. Larger address space than RAM:**
Virtual address space can be much larger than physical RAM. Not all pages need to be in RAM at once.

**3. Demand paging:**
Pages are only loaded into RAM when first accessed. Program starts instantly, loads code on demand.

**4. Swap space:**
When RAM is full, OS can move pages to disk (swap), making room for new pages. Process continues working (slowly) despite RAM pressure.

**5. Copy-on-write (COW) after fork():**
After `fork()`, parent and child share the same physical pages (mapped read-only in both page tables). When either writes to a page, a page fault fires, the OS copies the page, and both have their own copy. This makes `fork()` + `exec()` extremely fast (no actual copying until necessary).

**6. Shared libraries:**
The C library (libc) code is mapped (read-only) into every process's address space. Only ONE physical copy exists, but all processes see it.

**7. Memory-mapped files:**
Files can be mapped into the address space. Accessing the memory reads from the file. Modifying the memory writes to the file.

**8. Guard pages:**
Placing an unmapped page below the stack catches stack overflows: the overflowing code accesses the guard page, triggers a fault, OS sends SIGSEGV. Without this, a stack overflow would silently corrupt heap memory.

**9. No address conflict:**
Programs compiled to run at address 0x400000 (the default on x86-64 Linux) can all run simultaneously without conflict — each has its own virtual address space.

---

## 10. Virtual Memory on Real Hardware (x86-64)

x86-64 uses a 4-level (or 5-level) page table hierarchy.

**Why hierarchy?**
A flat page table for a 64-bit address space would be massive:
```
48-bit virtual address space → 2^48 / 2^12 = 2^36 = 68 billion pages
Each PTE is 8 bytes → 68 billion × 8 = 512 GB just for one process's page table!
```

A hierarchical page table only allocates entries for pages that actually exist:

**x86-64 four-level paging (48-bit addresses):**
```
Virtual Address (48 bits):
┌──────────┬──────────┬──────────┬──────────┬────────────┐
│  PML4    │   PDP    │    PD    │    PT    │   Offset   │
│  (9bits) │  (9bits) │  (9bits) │  (9bits) │  (12 bits) │
└──────────┴──────────┴──────────┴──────────┴────────────┘
    47-39      38-30      29-21      20-12       11-0

Each level: 512 entries × 8 bytes = 4KB (fits exactly in one page!)
```

**Page table walk:**
```
CR3 → Physical address of PML4 table
      ↓ index with bits[47:39]
      PML4 Entry → Physical address of PDPT table
                   ↓ index with bits[38:30]
                   PDPT Entry → Physical address of PD table
                                ↓ index with bits[29:21]
                                PD Entry → Physical address of PT table
                                           ↓ index with bits[20:12]
                                           PT Entry → Physical Frame Number
                                                      ↓ + bits[11:0] (offset)
                                                      Physical Address!
```

**4 levels = 4 memory accesses** to translate one virtual address (without TLB).

**Huge pages:**
Setting the Page Size bit in a PD entry means it points directly to a 2MB physical frame (skipping the PT level). Reduces TLB pressure for large allocations (databases, JVM).

Setting the Page Size bit in a PDPT entry gives a 1GB page (skipping both PD and PT levels).

---

## 11. Setting Up Paging — A Conceptual Walk-Through

When writing an OS, enabling paging is one of the first things you do. Here's the sequence (for 32-bit x86 protected mode):

**1. Prepare the page directory:**
```c
// 32-bit: one page directory, 1024 PDEs, each pointing to a page table
uint32_t page_directory[1024];  // aligned to 4KB boundary!

// Identity map the first 4MB (kernel startup code needs this)
uint32_t page_table_0[1024];

// Map virtual 0x00000000 – 0x003FFFFF → physical 0x00000000 – 0x003FFFFF
for (int i = 0; i < 1024; i++) {
    page_table_0[i] = (i * 4096) | 0x3;  // present | read-write
}
page_directory[0] = (uint32_t)page_table_0 | 0x3;  // present | read-write
```

**2. Load CR3:**
```c
// Write physical address of page directory into CR3
asm volatile("mov %0, %%cr3" : : "r"((uint32_t)page_directory));
```

**3. Enable paging in CR0:**
```c
uint32_t cr0;
asm volatile("mov %%cr0, %0" : "=r"(cr0));
cr0 |= 0x80000000;  // set PG (paging) bit
asm volatile("mov %0, %%cr0" : : "r"(cr0));
// From this instruction forward, all addresses are virtual!
```

After step 3, the CPU translates every address through the page table. If any address used before step 3 isn't mapped in the new page table, you'll immediately get a page fault.

---

## Summary

| Concept | Definition |
|---------|-----------|
| Virtual address | Address a program sees; may not correspond to physical location |
| Physical address | Actual location in RAM chips |
| Page | 4KB chunk of virtual address space |
| Frame | 4KB chunk of physical RAM |
| Page table | Per-process mapping from virtual pages to physical frames |
| PTE | Page Table Entry: physical frame number + flags |
| Present bit | 0 = page not in RAM (triggers page fault); 1 = in RAM |
| NX bit | 1 = code cannot execute from this page |
| MMU | CPU hardware that translates virtual → physical addresses |
| TLB | Cache of recent virtual → physical translations |
| CR3 | CPU register pointing to current process's page directory |
| Page fault | Exception when MMU can't translate (page not present, or protection violation) |
| COW | Copy-on-write: share physical pages after fork(); copy only on write |

**Virtual memory is non-optional in modern OSes.** It enables process isolation, demand paging, swap, shared libraries, and memory-mapped files. Without it, we'd be back to the problems of the 1960s: address conflicts, no protection, programs limited to physical RAM.
