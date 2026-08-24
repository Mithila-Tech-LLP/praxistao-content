# Chapter 16: Physical Memory Basics

> **"Memory is the most scarce and most fought-over resource in any computer. Every running program wants more of it. The OS is the arbiter — it decides who gets what, when, and how much. To do this fairly and efficiently, it must deeply understand what physical memory actually is."**

---

## Table of Contents

1. [What Is Physical Memory?](#1-what-is-physical-memory)
2. [Memory Addresses — The Address Space](#2-memory-addresses--the-address-space)
3. [How Programs Use Memory — The Memory Model](#3-how-programs-use-memory--the-memory-model)
4. [Memory Fragmentation](#4-memory-fragmentation)
5. [Early Memory Allocation — Fixed and Variable Partitions](#5-early-memory-allocation--fixed-and-variable-partitions)
6. [Physical Memory Layout at Boot](#6-physical-memory-layout-at-boot)
7. [Memory Pages — The Unit of Management](#7-memory-pages--the-unit-of-management)
8. [The Physical Memory Manager](#8-the-physical-memory-manager)
9. [Memory Types: Stack vs Heap vs Static](#9-memory-types-stack-vs-heap-vs-static)
10. [Memory Alignment](#10-memory-alignment)
11. [Summary](#summary)

---

## 1. What Is Physical Memory?

**Physical memory** is the actual RAM chips installed in your computer. When you buy "16GB RAM," you're buying 16 gigabytes of physical memory.

**Physical properties:**
- Divided into 8-byte (64-bit) or 4-byte (32-bit) addressable locations
- Each location has a unique physical address
- Volatile: contents are lost when power is cut
- Random access: any address can be read in ~100 nanoseconds (vs. milliseconds for SSD)
- Limited: 8GB, 16GB, 32GB — whatever is installed

**How it's organized:**
```
Physical Memory (e.g., 8GB = 8,589,934,592 bytes)
┌─────────────────────────────────────────────────┐
│ Address 0x00000000 │ byte value                  │
│ Address 0x00000001 │ byte value                  │
│ Address 0x00000002 │ byte value                  │
│ ...                │                             │
│ Address 0x1FFFFFFF │ byte value  ← 8GB = 0x200000000
└─────────────────────────────────────────────────┘
```

The OS must manage this entire range — track which bytes are in use, by whom, and allocate/free regions on demand.

---

## 2. Memory Addresses — The Address Space

A **memory address** is a number that identifies a specific location in memory.

**32-bit address space:**
- 32-bit pointer → 2³² possible addresses = 4,294,967,296 bytes = 4GB maximum
- This was the limit for 32-bit operating systems (Windows XP, 32-bit Linux)

**64-bit address space:**
- 64-bit pointer → 2⁶⁴ possible addresses = 18 exabytes (theoretical)
- In practice, modern CPUs only implement 48 bits: 2⁴⁸ = 256TB
- This is far more than any current physical RAM

**Physical vs. Virtual addresses:**

| Type | What it is | Who sees it |
|------|-----------|------------|
| Physical address | Actual byte location in RAM chips | CPU hardware, memory controller |
| Virtual address | Address a program thinks it has | Programs, programmer |

A program running today never sees physical addresses. It sees virtual addresses. The MMU (Memory Management Unit, a hardware component of the CPU) translates virtual → physical.

```
Program writes to virtual address 0x7fff_c000_1234
    ↓ (MMU translation)
CPU actually writes to physical address 0x0000_0001_4b28_0000
```

This translation is what enables virtual memory (Chapter 18).

---

## 3. How Programs Use Memory — The Memory Model

A typical C program uses memory in these regions:

```
High virtual address
┌──────────────────────────────┐
│  Kernel space                │  ← OS uses this; user can't access
├──────────────────────────────┤  0xFFFF8000_00000000 (64-bit Linux)
│  Stack                       │  ← grows ↓
│  (local variables,           │    8MB default size
│   function call frames)      │    stack overflow → SIGSEGV
│                              │
│  [unused space between       │
│   stack and heap]            │
│                              │
│  Memory-mapped region        │  ← mmap(), shared libraries (.so)
│                              │
│  Heap                        │  ← grows ↑
│  (malloc/free allocations)   │    managed by C library
│                              │
│  BSS segment                 │  ← uninitialized globals (zeroed at start)
│  Data segment                │  ← initialized globals
│  Text segment                │  ← executable code (read-only)
│  NULL page                   │  ← address 0 is unmapped (catches null pointer deref)
└──────────────────────────────┘
Low virtual address (0x0)
```

**Sizes:**
- Text: size of your compiled code (typically KBs to MBs)
- Data + BSS: global variables (typically small)
- Heap: grows as needed via malloc()
- Stack: fixed limit (typically 8MB per thread on Linux)
- Memory-mapped: shared libraries (libc ~2MB, plus other libs)

---

## 4. Memory Fragmentation

**Fragmentation** is when memory has enough total free space but not enough CONTIGUOUS free space to satisfy a request.

**External fragmentation:**
Free memory exists in small scattered pieces that can't be combined:
```
Memory:  [USED][FREE 10B][USED][FREE 20B][USED][FREE 5B][USED]
Request for 30 bytes → FAILS (no 30-byte contiguous block)
Total free: 35 bytes (10 + 20 + 5), but not contiguous!
```

External fragmentation is a problem for the physical memory allocator.

**Internal fragmentation:**
Allocated block is larger than requested:
```
Request: 13 bytes
Allocator gives: 16 bytes (padded to alignment)
Wasted: 3 bytes (inside the allocated block, unusable)
```

**Why fragmentation matters:**
Over time, malloc() and free() create a swiss cheese of small free blocks. A large allocation may fail even when total free memory is sufficient.

**Solutions:**
- **Compaction:** Physically move data to merge free spaces. Very expensive (copies lots of data). Only practical for managed runtimes (Java GC).
- **Segregated lists:** Keep separate free lists for different sizes. Small blocks go to the small-block list. Reduces fragmentation for common patterns.
- **Page-based allocation:** Allocate in units of pages (4KB). Reduces external fragmentation to at most one page per allocation.
- **Buddy system:** Allocate in powers of 2. Enables fast splitting and merging.

---

## 5. Early Memory Allocation — Fixed and Variable Partitions

Early operating systems used simple approaches:

**Fixed partitions (IBM OS/360, ~1960s):**
RAM is divided into fixed-size regions (partitions). Each partition runs one program.
```
RAM: [OS][Partition1: 2MB][Partition2: 2MB][Partition3: 4MB]
     Program A (1.5MB) fits in Partition1
     Program B (3MB) fits in Partition3
     1MB wasted in Partition1 (internal fragmentation)
     2MB wasted — can't fit Program C (2.5MB) in remaining 2MB partition
```
Problems: internal fragmentation, inflexible sizes.

**Variable partitions:**
Partitions sized to fit each program exactly.
```
RAM: [OS][Prog A: 1.5MB][Prog B: 3MB][FREE: 1.5MB]
Prog A exits → [OS][FREE: 1.5MB][Prog B: 3MB][FREE: 1.5MB]
New 2MB program → can't fit in either free block!
```
External fragmentation. Compaction needed.

Modern OSes don't use these approaches. They use **paging** (Chapter 18), which eliminates both problems.

---

## 6. Physical Memory Layout at Boot

When the computer boots, the BIOS/UEFI reports to the OS which physical memory regions are usable.

**BIOS memory map (E820 map):**
```
BIOS E820 memory map:
  Address     Length      Type
  0x00000000  0x0009fc00  Usable (RAM)
  0x0009fc00  0x00000400  Reserved (BIOS data area)
  0x000f0000  0x00010000  Reserved (ROM BIOS)
  0x00100000  0x7fef0000  Usable (1MB to ~2GB: main RAM)
  0x7fef0000  0x00110000  Reserved (ACPI tables)
  0x7ff00000  0x00100000  Reserved (ACPI non-volatile)
  0xfec00000  0x00001000  Reserved (APIC)
  0xfee00000  0x00001000  Reserved (APIC)
  0xfffc0000  0x00040000  Reserved (BIOS ROM)
```

The OS must:
1. Read this map at boot (from BIOS/UEFI)
2. Mark all usable pages as FREE
3. Mark reserved regions as IN USE (never allocate from them)
4. Mark where the kernel code itself was loaded (mark as IN USE)
5. Start serving page allocations from the free list

---

## 7. Memory Pages — The Unit of Management

Modern OS memory management works in units of **pages** — not individual bytes.

**Page size:**
- x86/x86-64: 4096 bytes (4KB) standard page
- Also supports huge pages: 2MB (x86 Large Page) and 1GB (x86 Huge Page)
- ARM: typically 4KB, also 16KB/64KB on some processors

**Why pages?**
- Hardware (MMU) translates addresses in page-sized units
- Page tables describe mappings for entire pages, not individual bytes
- Disk swap: swap in/out in page-sized chunks
- Memory allocator: manages physical "page frames"

**Page frame number:**
Physical address is divided into:
```
Physical Address: 0x0012_3456
                  │    │
                  │    └── offset within page (12 bits = 0-4095 for 4KB pages)
                  └──────── page frame number (the rest)

0x0012_3456 → page frame = 0x123, offset = 0x456
```

The OS physical memory manager works in terms of page frames — 4KB blocks. Allocating memory = allocating a page frame.

---

## 8. The Physical Memory Manager

The **Physical Memory Manager (PMM)** is the kernel subsystem that tracks which physical page frames are free vs. in use.

**The simplest PMM: bitmap allocator**

One bit per physical page frame:
- 0 = FREE
- 1 = IN USE

```c
// For 4GB RAM with 4KB pages: 4GB / 4KB = 1,048,576 pages
// 1,048,576 bits = 128 KB for the bitmap

#define PAGE_SIZE 4096
#define TOTAL_PAGES (RAM_SIZE / PAGE_SIZE)

uint8_t bitmap[TOTAL_PAGES / 8];  // 1 bit per page

// Mark page as used:
bitmap[page_num / 8] |= (1 << (page_num % 8));

// Mark page as free:
bitmap[page_num / 8] &= ~(1 << (page_num % 8));

// Check if free:
return !(bitmap[page_num / 8] & (1 << (page_num % 8)));

// Allocate: find first free page
for (int i = 0; i < TOTAL_PAGES; i++) {
    if (is_free(i)) {
        mark_used(i);
        return i * PAGE_SIZE;  // return physical address
    }
}
```

**Linux's buddy allocator (more sophisticated):**
- Manages pages in groups of 2^n (buddies)
- When a large block is freed, it merges with its "buddy" (adjacent block of same size)
- Fast allocation of contiguous pages
- Reduces fragmentation for kernel allocations

```
Order 0: individual pages (4KB)
Order 1: 2-page blocks (8KB)
Order 2: 4-page blocks (16KB)
...
Order 10: 1024-page blocks (4MB)

Allocate 16KB:
  Check order-2 free list → found! Mark as used. Done.
  
Allocate 4KB after freeing a 4KB block next to another free 4KB block:
  Mark 4KB free → check if buddy is also free → yes → merge into 8KB → 
  check if buddy is also free → no → done. Free 8KB block on order-1 list.
```

---

## 9. Memory Types: Stack vs Heap vs Static

**Stack memory:**
```c
void function() {
    int x = 5;         // stack memory — allocated automatically
    int arr[100];       // 400 bytes on stack
    // ...
}   // ← all stack memory freed automatically here
```
- Allocated/freed automatically by function call/return
- Very fast (just move the stack pointer: `SUB RSP, 400`)
- Limited size (typically 8MB per thread)
- No fragmentation (LIFO = perfect stacking)
- Cannot return pointer to local variable (stack is freed on return)

**Heap memory:**
```c
int *arr = malloc(100 * sizeof(int));  // 400 bytes on heap
arr[0] = 42;
free(arr);  // must manually free!
```
- Explicitly allocated with malloc/new; explicitly freed with free/delete
- Can persist beyond function scope
- Larger size (limited by virtual address space and RAM)
- Fragmentation possible
- Memory leaks if you forget to free

**Static memory:**
```c
static int counter = 0;    // data segment (initialized)
int lookup_table[256];     // BSS segment (uninitialized, zeroed)
```
- Fixed size, determined at compile time
- Exists for lifetime of the program
- Not allocated/freed dynamically

---

## 10. Memory Alignment

CPUs access memory most efficiently when data is **aligned** to its natural size.

**Natural alignment:**
- `int` (4 bytes): best at address divisible by 4
- `long` (8 bytes): best at address divisible by 8
- `__m256` (32 bytes): best at address divisible by 32 (SIMD)

**Misaligned access:**
On modern x86: works but potentially slower (spans cache lines).
On ARM/MIPS: may trigger hardware exception (bus error) or be split into two accesses.

**Example:**
```c
struct Misaligned {
    char a;      // offset 0
    int b;       // offset 1? NO — compiler pads to offset 4
    char c;      // offset 8
    // trailing pad to multiple of 4
};
// sizeof(struct Misaligned) = 12, not 6!

// Memory layout:
// [a][pad][pad][pad][b...b...b...b][c][pad][pad][pad]
//  0   1    2    3   4   5   6   7  8   9   10   11
```

Understanding alignment is crucial for:
- Writing OS code that interfaces with hardware (memory-mapped registers have strict alignment)
- Writing high-performance code (cache line alignment)
- Parsing binary file formats or network packets
- Writing kernel memory allocators (must return properly aligned blocks)

---

## Summary

| Concept | Definition |
|---------|-----------|
| Physical memory | Actual RAM chips; finite, volatile, byte-addressable |
| Physical address | The actual hardware address of a byte in RAM |
| Virtual address | The address a program sees; translated to physical by MMU |
| Page | Fixed-size unit of memory management (typically 4KB) |
| Page frame | A 4KB block of physical memory, numbered from 0 |
| Physical memory manager | Kernel tracks which physical pages are free vs used |
| Bitmap allocator | Simplest PMM: 1 bit per page, 0=free, 1=used |
| Buddy allocator | Linux's PMM: manages power-of-2 blocks, merges buddies |
| External fragmentation | Free memory scattered, can't satisfy large contiguous request |
| Internal fragmentation | Allocated block larger than requested; wasted within block |
| Memory alignment | Data at addresses evenly divisible by its size — required by CPU |
