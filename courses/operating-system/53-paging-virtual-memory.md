# Chapter 53: Paging and Virtual Memory

> **"Paging is the magic that lets every process believe it owns the entire address space. The kernel lives at 0xC0000000 in every process simultaneously. Processes can't see each other's memory. The heap can be larger than physical RAM. All of this is accomplished by one mechanism: a hardware-enforced mapping from virtual addresses to physical addresses via page tables."**

---

## Table of Contents

1. [Virtual vs Physical Addresses](#1-virtual-vs-physical-addresses)
2. [How Paging Works — Page Directory and Page Tables](#2-how-paging-works--page-directory-and-page-tables)
3. [Two-Level Paging in x86](#3-two-level-paging-in-x86)
4. [Page Table Entry Format](#4-page-table-entry-format)
5. [Enabling Paging](#5-enabling-paging)
6. [Mapping Physical Memory 1:1](#6-mapping-physical-memory-11)
7. [Page Fault Handler](#7-page-fault-handler)
8. [vmm.c — Virtual Memory Manager](#8-vmmh)
9. [Testing Paging](#9-testing-paging)
10. [Summary](#summary)

---

## 1. Virtual vs Physical Addresses

```
Before paging (what we have after Ch 47-52):
  Every address in your code IS a physical address
  mov eax, [0x1000]  → reads from physical RAM at 0x1000
  Two processes running would see the same physical memory
  Process A could read/write Process B's data
  No protection, no isolation

After paging:
  Every address in your code is a VIRTUAL address
  mov eax, [0x1000]  → CPU translates via page tables → reads physical addr X
  Process A's virtual 0x1000 maps to physical X
  Process B's virtual 0x1000 maps to physical Y (completely different)
  
  Advantages:
    ✓ Isolation: each process has its own virtual address space
    ✓ Protection: pages marked read-only, user/kernel, present/absent
    ✓ Virtual kernel: kernel can be at same virtual address in every process
    ✓ Lazy allocation: map a virtual page but don't allocate physical until access
    ✓ Demand paging / swapping: physical frame can be on disk, mapped when needed
```

---

## 2. How Paging Works — Page Directory and Page Tables

```
Virtual address translation (x86 32-bit, two-level paging):

32-bit virtual address:
  ┌──────────────┬──────────────┬────────────────┐
  │  DIR [31:22] │  TBL [21:12] │ OFFSET [11:0]  │
  └──────────────┴──────────────┴────────────────┘
     10 bits         10 bits         12 bits
  = 1024 entries   = 1024 entries   = 4096 bytes
  
Step 1: CPU reads CR3 → physical address of Page Directory
Step 2: Use DIR[31:22] as index into Page Directory
        → Get Page Directory Entry (PDE) → physical address of a Page Table
Step 3: Use TBL[21:12] as index into Page Table
        → Get Page Table Entry (PTE) → physical address of the page
Step 4: Add OFFSET → final physical address

Example: virtual 0xC0001234
  DIR    = 0xC0001234 >> 22 = 768        → PD[768]
  TBL    = (0xC0001234 >> 12) & 0x3FF = 1 → PT[1]
  OFFSET = 0xC0001234 & 0xFFF = 0x234
  
  If PD[768] → page table at physical 0x400000
  And PT[1]  → physical page at 0x200000
  Then virtual 0xC0001234 → physical 0x200234
```

---

## 3. Two-Level Paging in x86

```
Page Directory (PD):
  4KB structure (1024 × 4-byte entries)
  One per address space (per process)
  Physical address stored in CR3
  
  PD[0]   → page table covering virtual 0x00000000 - 0x003FFFFF
  PD[1]   → page table covering virtual 0x00400000 - 0x007FFFFF
  ...
  PD[768] → page table covering virtual 0xC0000000 - 0xC03FFFFF
  ...
  PD[1023]→ page table covering virtual 0xFFC00000 - 0xFFFFFFFF

Page Table (PT):
  4KB structure (1024 × 4-byte entries)
  One per 4MB virtual address region
  Physical address stored in PDE
  
  PT[0]   → physical page at (physical address in bits[31:12])
  ...
  PT[1023]→ physical page
  
Total addressable: 1024 PDs × 1024 PTs × 4KB = 4GB (full 32-bit space)
```

---

## 4. Page Table Entry Format

Both PDE and PTE use the same 32-bit format:

```
31                  12 11   9  8   7   6   5   4   3   2   1   0
┌─────────────────────┬──────┬───┬───┬───┬───┬───┬───┬───┬───┬───┐
│  Frame addr [31:12] │ AVL  │ G │ PS│ D │ A │ C │ W │ U │ R │ P │
└─────────────────────┴──────┴───┴───┴───┴───┴───┴───┴───┴───┴───┘

Bit 0: P  = Present (1 = page in memory; 0 = not present → page fault)
Bit 1: R  = Read/Write (0 = read only; 1 = read+write)
Bit 2: U  = User/Supervisor (0 = kernel only; 1 = user can access)
Bit 3: W  = Write-Through caching
Bit 4: C  = Cache Disable
Bit 5: A  = Accessed (CPU sets when page is read or written)
Bit 6: D  = Dirty (CPU sets when page is written; PTE only)
Bit 7: PS = Page Size (PDE only): 0=4KB, 1=4MB (PSE) 
Bit 8: G  = Global (TLB not flushed on CR3 reload; used for kernel pages)
Bits 9-11: AVL = Available for OS use
Bits 31:12: Physical frame address (top 20 bits; bottom 12 always 0 → 4KB aligned)
```

**Common flag combinations:**
```c
#define PTE_PRESENT    (1 << 0)
#define PTE_WRITABLE   (1 << 1)
#define PTE_USER       (1 << 2)
#define PTE_ACCESSED   (1 << 5)
#define PTE_DIRTY      (1 << 6)

/* Kernel read-write page: */
#define PTE_KERNEL_RW  (PTE_PRESENT | PTE_WRITABLE)

/* User read-write page: */
#define PTE_USER_RW    (PTE_PRESENT | PTE_WRITABLE | PTE_USER)

/* Kernel read-only page: */
#define PTE_KERNEL_RO  (PTE_PRESENT)
```

---

## 5. Enabling Paging

```c
/* Enable paging:
   1. Fill in page directory and page tables
   2. Load CR3 with physical address of page directory
   3. Set CR0.PG bit */

static inline void enable_paging(uint32_t page_dir_phys) {
    /* Load page directory: */
    __asm__ volatile("mov %0, %%cr3" : : "r"(page_dir_phys));
    
    /* Enable paging (set bit 31 in CR0): */
    uint32_t cr0;
    __asm__ volatile("mov %%cr0, %0" : "=r"(cr0));
    cr0 |= 0x80000000;
    __asm__ volatile("mov %0, %%cr0" : : "r"(cr0));
}

/* Flush TLB for one page (after changing a mapping): */
static inline void flush_tlb_page(uint32_t vaddr) {
    __asm__ volatile("invlpg (%0)" : : "r"(vaddr) : "memory");
}

/* Flush entire TLB (by reloading CR3 — all non-global entries invalidated): */
static inline void flush_tlb_all(void) {
    uint32_t cr3;
    __asm__ volatile("mov %%cr3, %0" : "=r"(cr3));
    __asm__ volatile("mov %0, %%cr3" : : "r"(cr3));
}
```

---

## 6. Mapping Physical Memory 1:1

For our initial kernel setup, we do an **identity map**: virtual address = physical address. This makes the transition into paging seamless — our kernel code is still reachable at the same addresses.

```c
/* kernel/vmm.c — Virtual Memory Manager */

#include "vmm.h"
#include "pmm.h"
#include "string.h"

/* Page directory for the kernel (allocated statically): */
static uint32_t kernel_page_dir[1024] __attribute__((aligned(4096)));

/* Map a virtual page to a physical frame: */
void vmm_map_page(uint32_t *page_dir, uint32_t vaddr, uint32_t paddr, uint32_t flags) {
    uint32_t dir_idx = vaddr >> 22;          /* top 10 bits */
    uint32_t tbl_idx = (vaddr >> 12) & 0x3FF; /* middle 10 bits */
    
    /* Does the page table exist? */
    if (!(page_dir[dir_idx] & PTE_PRESENT)) {
        /* Allocate a new page table: */
        uint32_t pt_phys = pmm_alloc_frame();
        if (!pt_phys) return; /* OOM */
        memset((void *)pt_phys, 0, PAGE_SIZE);
        
        /* Install the page table in the directory: */
        page_dir[dir_idx] = pt_phys | PTE_PRESENT | PTE_WRITABLE | PTE_USER;
    }
    
    /* Get pointer to the page table: */
    uint32_t *page_table = (uint32_t *)(page_dir[dir_idx] & ~0xFFF);
    
    /* Map the page: */
    page_table[tbl_idx] = (paddr & ~0xFFF) | (flags & 0xFFF) | PTE_PRESENT;
}

/* Unmap a virtual page: */
void vmm_unmap_page(uint32_t *page_dir, uint32_t vaddr) {
    uint32_t dir_idx = vaddr >> 22;
    uint32_t tbl_idx = (vaddr >> 12) & 0x3FF;
    
    if (!(page_dir[dir_idx] & PTE_PRESENT)) return;
    
    uint32_t *page_table = (uint32_t *)(page_dir[dir_idx] & ~0xFFF);
    page_table[tbl_idx] = 0;
    flush_tlb_page(vaddr);
}

/* Set up kernel paging: identity-map the first 4MB + kernel memory: */
void vmm_init(void) {
    /* Zero the page directory: */
    memset(kernel_page_dir, 0, PAGE_SIZE);
    
    /* Identity-map the first 4MB (covers kernel, VGA, BIOS areas):
       Use a PSE (4MB) entry for simplicity, or map page by page. */
    
    /* We'll map page by page for correctness: */
    /* Map 0x0 - 0x400000 (first 4MB): */
    for (uint32_t addr = 0; addr < 0x400000; addr += PAGE_SIZE) {
        vmm_map_page(kernel_page_dir, addr, addr, PTE_WRITABLE);
    }
    
    /* Enable paging: */
    enable_paging((uint32_t)kernel_page_dir);
    
    kprintf("VMM: Paging enabled. Kernel at virtual = physical addresses.\n");
}

/* Allocate and map a virtual page, returning its virtual address: */
uint32_t vmm_alloc_page(uint32_t *page_dir, uint32_t vaddr, uint32_t flags) {
    uint32_t phys = pmm_alloc_frame();
    if (!phys) return 0;
    vmm_map_page(page_dir, vaddr, phys, flags);
    return vaddr;
}

/* Get the physical address for a virtual address: */
uint32_t vmm_get_phys(uint32_t *page_dir, uint32_t vaddr) {
    uint32_t dir_idx = vaddr >> 22;
    uint32_t tbl_idx = (vaddr >> 12) & 0x3FF;
    uint32_t offset  = vaddr & 0xFFF;
    
    if (!(page_dir[dir_idx] & PTE_PRESENT)) return 0;
    
    uint32_t *pt = (uint32_t *)(page_dir[dir_idx] & ~0xFFF);
    if (!(pt[tbl_idx] & PTE_PRESENT)) return 0;
    
    return (pt[tbl_idx] & ~0xFFF) | offset;
}

/* Create a new page directory (copy kernel mappings): */
uint32_t *vmm_create_address_space(void) {
    uint32_t *new_dir = (uint32_t *)pmm_alloc_frame();
    if (!new_dir) return NULL;
    memset(new_dir, 0, PAGE_SIZE);
    
    /* Copy kernel mappings (upper half — dir entries 768-1023): */
    for (int i = 768; i < 1024; i++) {
        new_dir[i] = kernel_page_dir[i];
    }
    
    return new_dir;
}

/* Switch to a different address space: */
void vmm_switch(uint32_t *page_dir) {
    __asm__ volatile("mov %0, %%cr3" : : "r"((uint32_t)page_dir));
}
```

---

## 7. Page Fault Handler

We already handle page faults in `isr_handler`. Let's make it more informative:

```c
/* In isr.c, make the page fault handler more specific: */
if (regs->int_no == 14) {
    uint32_t faulting_addr;
    __asm__ volatile("mov %%cr2, %0" : "=r"(faulting_addr));
    
    int present   = !(regs->err_code & 0x1);  /* 1 = protection violation */
    int write     = regs->err_code & 0x2;      /* 1 = write, 0 = read */
    int user      = regs->err_code & 0x4;      /* 1 = user mode, 0 = kernel */
    int reserved  = regs->err_code & 0x8;      /* 1 = reserved bit set in PTE */
    int inst_fetch = regs->err_code & 0x10;    /* 1 = instruction fetch */
    
    kprintf("PAGE FAULT at 0x%x\n", faulting_addr);
    kprintf("Cause: %s %s from %s\n",
            present ? "protection violation" : "not-present page",
            write ? "on write" : "on read",
            user ? "user mode" : "kernel mode");
    if (reserved) kprintf("  (reserved bit in PTE was set)\n");
    if (inst_fetch) kprintf("  (instruction fetch)\n");
    kprintf("EIP: 0x%x\n", regs->eip);
    
    /* For now, halt. Later we could handle demand paging here. */
    for (;;) __asm__ volatile("cli; hlt");
}
```

---

## 8. vmm.h

```c
/* include/vmm.h */
#pragma once
#include "stdint.h"
#include "pmm.h"

/* Page flags: */
#define PTE_PRESENT   (1 << 0)
#define PTE_WRITABLE  (1 << 1)
#define PTE_USER      (1 << 2)
#define PTE_ACCESSED  (1 << 5)
#define PTE_DIRTY     (1 << 6)

void      vmm_init(void);
void      vmm_map_page(uint32_t *pd, uint32_t vaddr, uint32_t paddr, uint32_t flags);
void      vmm_unmap_page(uint32_t *pd, uint32_t vaddr);
uint32_t  vmm_alloc_page(uint32_t *pd, uint32_t vaddr, uint32_t flags);
uint32_t  vmm_get_phys(uint32_t *pd, uint32_t vaddr);
uint32_t *vmm_create_address_space(void);
void      vmm_switch(uint32_t *pd);
```

---

## 9. Testing Paging

```c
/* kernel_main after pmm_init: */

vmm_init();   /* Enable paging */

/* Test: map a new virtual page and write to it: */
uint32_t test_vaddr = 0x400000;  /* 4MB — just beyond our identity map */
vmm_alloc_page(kernel_page_dir, test_vaddr, PTE_WRITABLE);

volatile uint32_t *ptr = (volatile uint32_t *)test_vaddr;
*ptr = 0xDEADBEEF;
kprintf("Wrote 0x%x to virtual 0x%x\n", *ptr, test_vaddr);

/* Check physical address: */
uint32_t phys = vmm_get_phys(kernel_page_dir, test_vaddr);
kprintf("Virtual 0x%x → Physical 0x%x\n", test_vaddr, phys);

/* Read it back from physical address: */
volatile uint32_t *phys_ptr = (volatile uint32_t *)phys;
kprintf("Physical read: 0x%x (should be 0xDEADBEEF)\n", *phys_ptr);
```

If paging is working:
1. Writing to virtual 0x400000 succeeds (page is mapped)
2. Physical address returned is a valid RAM address
3. Reading from physical address gives back 0xDEADBEEF

---

## Summary

| Concept | Description |
|---------|------------|
| Virtual address | Address in code; CPU translates to physical via page tables |
| Physical address | Real location in RAM; what the memory bus sees |
| Page Directory | 4KB array of 1024 PDEs; one per process; physical address in CR3 |
| Page Table | 4KB array of 1024 PTEs; one per 4MB virtual region |
| PDE/PTE | 32-bit entry: top 20 bits = frame address, bottom 12 = flags |
| PTE_PRESENT | Bit 0: page is in RAM. If 0 and accessed → page fault (exception 14) |
| PTE_WRITABLE | Bit 1: writes allowed. If 0 and write attempted → page fault |
| PTE_USER | Bit 2: user mode can access. If 0 and user accesses → page fault |
| CR3 | Points to current page directory; loading new value switches address space |
| invlpg | Invalidate TLB entry for one page (must call after changing a PTE) |
| Identity map | Virtual address = physical address; easy transition when enabling paging |
| Page fault handler | Exception 14; CR2 has faulting virtual address; error code has cause bits |
| TLB | Translation Lookaside Buffer: hardware cache of recent virtual→physical mappings |
| flush_tlb_all | Reload CR3 to invalidate all non-global TLB entries |
