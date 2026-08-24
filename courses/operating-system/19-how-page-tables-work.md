# Chapter 19: How Page Tables Work

> **"The page table is the kernel's contract with hardware. The CPU will enforce exactly what the OS writes into the page table — no more, no less. Understanding every bit of a PTE is understanding exactly how the CPU enforces process isolation and memory permissions."**

---

## Table of Contents

1. [Review: Why We Need Page Tables](#1-review-why-we-need-page-tables)
2. [The Linear Page Table (Simple)](#2-the-linear-page-table-simple)
3. [The Problem with Flat Page Tables](#3-the-problem-with-flat-page-tables)
4. [Hierarchical Page Tables](#4-hierarchical-page-tables)
5. [x86-32 Two-Level Paging](#5-x86-32-two-level-paging)
6. [x86-64 Four-Level Paging in Detail](#6-x86-64-four-level-paging-in-detail)
7. [The CR3 Register](#7-the-cr3-register)
8. [Walking the Page Table — Hardware and Software](#8-walking-the-page-table--hardware-and-software)
9. [Kernel Mapping in Every Process](#9-kernel-mapping-in-every-process)
10. [Inverted Page Tables](#10-inverted-page-tables)
11. [Hashed Page Tables](#11-hashed-page-tables)
12. [How Linux Manages Page Tables](#12-how-linux-manages-page-tables)
13. [Summary](#summary)

---

## 1. Review: Why We Need Page Tables

Chapter 18 introduced virtual memory. To recap:
- Programs use virtual addresses
- CPU hardware (MMU) must translate virtual → physical
- The page table is the data structure the MMU uses to do this translation
- Each process has its own page table (its own view of memory)
- The OS maintains and updates page tables

Now we go deeper: how are page tables actually structured, and what are all the details?

---

## 2. The Linear Page Table (Simple)

The simplest page table is a flat array, indexed by virtual page number:

```c
// 32-bit address space, 4KB pages
// Virtual pages: 2^32 / 2^12 = 2^20 = 1,048,576 pages
// Each PTE: 4 bytes
// Total size: 4MB per process!

uint32_t page_table[1048576];  // 4MB of page table per process

// To translate virtual address V:
uint32_t vpn = V >> 12;         // page number = upper 20 bits
uint32_t offset = V & 0xFFF;    // offset = lower 12 bits
uint32_t pte = page_table[vpn]; // look up entry
uint32_t pfn = pte >> 12;       // physical frame number
uint32_t phys = (pfn << 12) | offset; // physical address
```

**This works!** But 4MB per process is expensive:
- 100 processes × 4MB = 400MB just for page tables
- And most of the 4MB is never used (most virtual pages are unmapped)

---

## 3. The Problem with Flat Page Tables

A typical process only uses a small fraction of its virtual address space:
- Text segment: first few MB
- Libraries: somewhere in the middle
- Stack: top few MB
- Heap: growing area

```
Virtual address space (32-bit = 4GB):
0x00000000 ──► [code + data: first 10MB] — USED
               [gap: several GB]          — UNUSED
0xC0000000 ──► [kernel: last 1GB]        — USED
```

With a flat page table: we must store PTEs for ALL 1M virtual pages, even the 999,990 that are never used. Wasteful!

**64-bit is catastrophically worse:**
- 64-bit: 2^48 = 256TB of virtual address space
- Pages: 256TB / 4KB = 64 billion pages
- Flat PTE array: 64 billion × 8 bytes = 512 GB per process — impossible!

We need a smarter structure.

---

## 4. Hierarchical Page Tables

The solution: build the page table as a TREE. Only allocate nodes for parts of the address space that are actually used.

```
2-level page table (32-bit):

Virtual Address:  [ 10-bit PD index ] [ 10-bit PT index ] [ 12-bit offset ]

Page Directory (PD): 1024 entries (4KB)
    Each entry → physical address of a Page Table (or "not present")
    
Page Table (PT): 1024 entries (4KB)
    Each entry → physical frame number + flags

If a program only uses addresses 0x00000000 – 0x003FFFFF (first 4MB):
  Page Directory: only entry 0 is present (points to one PT)
  Pages 1–1023: NOT PRESENT (no page table allocated, 4KB each not needed)
  
  Memory for page tables: 1 PD (4KB) + 1 PT (4KB) = 8KB total
  vs. flat: 4MB total
```

**Key insight:** Non-present page directory entries mean NO page table is allocated for that region. Those 1023 page tables (× 4KB each = 4MB) are simply not allocated.

A program using only 4MB of its 4GB space needs only 8KB of page tables instead of 4MB.

---

## 5. x86-32 Two-Level Paging

Classic 32-bit x86 paging (PAE disabled):

**Address split:**
```
32-bit virtual address:
┌─────────────────┬─────────────────┬─────────────────────┐
│  PD Index       │   PT Index      │      Offset         │
│    (10 bits)    │   (10 bits)     │     (12 bits)       │
└─────────────────┴─────────────────┴─────────────────────┘
     Bits 31-22         Bits 21-12          Bits 11-0
```

**Translation:**
```
CR3 register → physical address of Page Directory

1. PD[bits 31-22] → PDE (Page Directory Entry)
   If PDE.Present == 0: page fault!
   PDE.address → physical address of Page Table

2. PT[bits 21-12] → PTE (Page Table Entry)
   If PTE.Present == 0: page fault!
   PTE.address → physical frame number

3. Physical address = (PTE.frame_number << 12) | (VA bits 11-0)
```

**x86-32 PDE format (4 bytes):**
```
Bit 31-12: Physical address of page table (right-shifted 12)
Bit 8:     Ignored
Bit 7:     PS (Page Size): 0 = 4KB, 1 = 4MB large page
Bit 6:     Ignored
Bit 5:     Accessed bit (set by CPU on access)
Bit 4:     PCD (Page Cache Disable)
Bit 3:     PWT (Page Write-Through)
Bit 2:     U/S (User/Supervisor)
Bit 1:     R/W (Read/Write)
Bit 0:     P (Present)
```

**x86-32 PTE format (4 bytes):**
```
Bit 31-12: Physical Frame Number
Bit 11-9:  Available (OS can use for its own purposes)
Bit 8:     Global
Bit 7:     PAT
Bit 6:     Dirty (set by CPU on write)
Bit 5:     Accessed (set by CPU on access)
Bit 4:     PCD
Bit 3:     PWT
Bit 2:     U/S
Bit 1:     R/W
Bit 0:     P (Present)
```

---

## 6. x86-64 Four-Level Paging in Detail

Modern 64-bit x86 (most common today):

**Address split (48-bit virtual addresses, 4KB pages):**
```
64-bit virtual address (only 48 bits used, bits 63-48 must be sign-extended from bit 47):
┌────────────┬────────────┬────────────┬────────────┬─────────────────┐
│  PML4 idx  │  PDPT idx  │  PD idx    │  PT idx    │    Offset       │
│  (9 bits)  │  (9 bits)  │  (9 bits)  │  (9 bits)  │   (12 bits)     │
└────────────┴────────────┴────────────┴────────────┴─────────────────┘
  Bits 47-39    Bits 38-30   Bits 29-21   Bits 20-12    Bits 11-0
```

Each level: 512 entries × 8 bytes = 4096 bytes = exactly one page!

**Four-level walk:**
```
CR3 → PML4 table (512 entries)
       │
       └─[PML4 index]─► PML4 Entry → PDPT table
                                       │
                                       └─[PDPT index]─► PDPT Entry → PD table
                                                                       │
                                                                       └─[PD index]─► PD Entry → PT table
                                                                                                  │
                                                                                                  └─[PT index]─► PTE → Physical Frame
                                                                                                                        +offset = Physical Address
```

**x86-64 PTE format (8 bytes):**
```
Bit 63:    XD/NX (Execute-Disable)
Bits 62-52: Ignored (12 bits available for OS use)
Bits 51-12: Physical Frame Number (40 bits → up to 4 petabytes physical address space)
Bit 11:    Available
Bit 10:    Available
Bit 9:     Available
Bit 8:     Global (1 = don't flush on CR3 change — for kernel pages)
Bit 7:     PAT (Page Attribute Table)
Bit 6:     Dirty
Bit 5:     Accessed
Bit 4:     PCD (Page Cache Disable)
Bit 3:     PWT (Page Write-Through)
Bit 2:     U/S (0 = kernel only, 1 = user accessible)
Bit 1:     R/W (0 = read-only, 1 = read-write)
Bit 0:     P (Present)
```

**Huge pages (skipping levels):**
- **2MB huge page:** Set PS bit in PD entry. Skip the PT level entirely.
  - PD entry directly maps a 2MB physical region
  - Physical address: PD_entry.base_addr + (VA bits 20-0)
- **1GB huge page:** Set PS bit in PDPT entry. Skip PD and PT.
  - Physical address: PDPT_entry.base_addr + (VA bits 29-0)

Benefits: fewer TLB entries needed for large allocations (one 2MB TLB entry vs. 512 4KB entries).

---

## 7. The CR3 Register

CR3 (Control Register 3) is the pivot of virtual memory:

```
x86-64 CR3 layout:
Bits 63-52: PCID (Process Context Identifier) — optional, 12 bits
Bits 51-12: Physical address of PML4 table (40 bits, shifted right 12)
Bits 11-5:  Reserved
Bit 4:      PCD (cache disable for PML4 page itself)
Bit 3:      PWT (write-through for PML4 page itself)
Bits 2-0:   Reserved
```

**Loading CR3:**
```c
// Load a new page table (change the process's address space):
asm volatile("mov %0, %%cr3" : : "r"(new_pml4_phys_addr) : "memory");
// This:
// 1. Flushes the TLB (all non-global entries invalidated)
// 2. From now on, all virtual translations use the new PML4
```

**When does the OS change CR3?**
- On every context switch (switching to a different process)
- When creating a new process (setting up its page tables)
- When exec()-ing (replacing the process's address space)

**PCID (Process Context ID):**
With PCID support (Intel INVPCID instruction, Linux uses this):
- TLB entries are tagged with PCID (12-bit identifier per process)
- CR3 load with PCID doesn't flush TLB for OTHER PCIDs
- Only invalidates the TLB when necessary

This dramatically reduces TLB flush overhead on context switches.

---

## 8. Walking the Page Table — Hardware and Software

**Hardware page table walk:**
When a TLB miss occurs, hardware automatically walks the page table:

```
1. Read CR3 → physical address of PML4 (PA_pml4)
2. Read memory[PA_pml4 + PML4_index * 8] → PML4 entry
   Check: present bit set? else: page fault
3. Extract PDPT physical address from PML4 entry
4. Read memory[PA_pdpt + PDPT_index * 8] → PDPT entry
   Check: present bit set? PS bit? else: page fault
5. Extract PD physical address from PDPT entry
6. Read memory[PA_pd + PD_index * 8] → PD entry
   Check: present bit set? PS bit? else: page fault
7. Extract PT physical address from PD entry
8. Read memory[PA_pt + PT_index * 8] → PTE
   Check: present bit, R/W, U/S permissions vs. access type
9. Extract physical frame number from PTE
10. Physical address = (frame_number << 12) | offset
11. Cache in TLB
```

This walk reads 4 physical memory locations (one per level). With a good TLB hit rate (~99%), most translations don't need to walk.

**Software page table walk (in the kernel):**
The OS needs to read/modify page tables (for mmap, fork, exec, swap, etc.):

```c
// Linux: walk a user process's page table to find or create a PTE
pmd_t *pmd = pmd_alloc(mm, pud, address);
pte_t *pte = pte_alloc_map(mm, pmd, address);
set_pte(pte, mk_pte(page, prot));
```

The kernel uses a set of functions (`pgd_offset`, `p4d_offset`, `pud_offset`, `pmd_offset`, `pte_offset`) to navigate page table levels, hiding hardware differences.

---

## 9. Kernel Mapping in Every Process

**Every process's page table maps the kernel at high addresses.**

On 32-bit Linux (classic layout):
- Lower 3GB (0x00000000–0xBFFFFFFF): user space
- Upper 1GB (0xC0000000–0xFFFFFFFF): kernel (same for ALL processes)

On 64-bit Linux:
- Lower half (0x0000000000000000–0x00007FFFFFFFFFFF): user space
- Upper half (0xFFFF800000000000–0xFFFFFFFFFFFFFFFF): kernel

**Why map kernel in every process?**
When a user process makes a system call, the CPU stays in the SAME process context — it just switches to kernel mode and starts executing kernel code. The kernel code needs its own data to be accessible. If the kernel wasn't mapped in the process's page table, the kernel code couldn't access kernel data structures.

**KPTI (Kernel Page Table Isolation):**
Spectre/Meltdown attacks (2018) could read kernel memory from user space by exploiting CPU speculative execution. The fix: KPTI.

With KPTI: the user-space page table has ALMOST NO kernel mappings. When a syscall occurs, the CPU switches to a separate kernel page table that HAS full kernel mappings. This requires a TLB flush on every syscall/return → significant overhead (~5–30% on some workloads).

AMD is not vulnerable to Meltdown, so KPTI can be disabled on AMD CPUs.

---

## 10. Inverted Page Tables

Alternative to hierarchical page tables: instead of indexing by virtual address, index by physical frame.

**Structure:**
One entry per physical frame (not one per virtual page):
```
Inverted page table[physical_frame_number] = {
    process_id,
    virtual_page_number
}
```

**Advantage:** Size proportional to physical RAM, not virtual address space.
For 8GB RAM: 8GB / 4KB = 2M entries × 16 bytes ≈ 32MB — fixed, regardless of address space size.

**Disadvantage:** Looking up a virtual address requires searching for a matching (PID, VPN) entry — O(n) search! In practice, use a hash table → O(1) average.

Used by: IBM PowerPC, HP PA-RISC, early Alpha processors.

---

## 11. Hashed Page Tables

Hash the virtual page number → find the physical frame.

```
hash(VPN) → bucket → search chain for matching VPN → physical frame
```

Used when virtual address space >> physical RAM (virtual address space can be much larger without wasted page table space).

Used by: many RISC architectures, some HPAs.

---

## 12. How Linux Manages Page Tables

Linux abstracts hardware differences with a 5-level generic page table:
```
pgd_t (Page Global Directory)    ← top level (PML4 on x86-64)
p4d_t (Page 4th Directory)       ← added for 5-level paging
pud_t (Page Upper Directory)     ← PDPT on x86-64
pmd_t (Page Middle Directory)    ← PD on x86-64
pte_t (Page Table Entry)         ← PT on x86-64
```

On 4-level hardware, p4d = pgd (the level is folded away). On 2-level hardware, pud and pmd are folded too.

**Linux memory zones:**
```
ZONE_DMA    (0–16MB)      — for legacy ISA DMA devices
ZONE_DMA32  (0–4GB)       — for 32-bit DMA devices
ZONE_NORMAL (4GB+)        — standard kernel memory
ZONE_HIGHMEM (32-bit only) — above kernel virtual mapping limit
```

**Memory regions (vm_area_struct):**
Linux tracks each process's virtual memory as a linked list of VMAs (Virtual Memory Areas):
```c
struct vm_area_struct {
    unsigned long vm_start;  // start virtual address
    unsigned long vm_end;    // end virtual address
    struct vm_area_struct *vm_next;
    pgprot_t vm_page_prot;   // page protection (read/write/exec)
    unsigned long vm_flags;  // VM_READ, VM_WRITE, VM_EXEC, VM_SHARED, ...
    struct file *vm_file;    // mapped file (if mmap'd), or NULL
};
```

VMAs are coarser than page tables. A VMA describes a region's policy; the page table handles individual page status.

---

## Summary

| Concept | Definition |
|---------|-----------|
| Flat page table | One entry per virtual page; too large for 64-bit |
| Hierarchical page table | Tree of tables; only allocated for used regions |
| x86-32 2-level | PD (10 bits) → PT (10 bits) → offset (12 bits) |
| x86-64 4-level | PML4 (9) → PDPT (9) → PD (9) → PT (9) → offset (12) |
| CR3 | Points to top-level page table; loaded on context switch |
| PTE present bit | 0 = page fault; 1 = in RAM, translation proceeds |
| PTE NX bit | 1 = cannot execute code from this page |
| PTE dirty bit | Set by CPU on write; tells OS if page must be written to swap |
| TLB miss | Triggers hardware page table walk; result cached in TLB |
| KPTI | Security: separate page tables for user/kernel to prevent Spectre/Meltdown |
| VMA | Linux: higher-level description of a virtual memory region |
| Huge pages | 2MB or 1GB pages; skip PT/PD levels; reduce TLB pressure |
