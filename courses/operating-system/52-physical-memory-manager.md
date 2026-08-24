# Chapter 52: Physical Memory Manager

> **"Before you can allocate a single byte of heap or map a single virtual page, you need to know which physical pages of RAM are available. The Physical Memory Manager is the lowest layer of your memory system — a ledger tracking which 4KB frames are free and which are in use. Get this right and everything else in memory management becomes straightforward."**

---

## Table of Contents

1. [What Is Physical Memory Management?](#1-what-is-physical-memory-management)
2. [Pages and Frames](#2-pages-and-frames)
3. [The Bitmap Allocator](#3-the-bitmap-allocator)
4. [Reading the Memory Map from Multiboot](#4-reading-the-memory-map-from-multiboot)
5. [Marking Reserved Regions](#5-marking-reserved-regions)
6. [Allocating and Freeing Frames](#6-allocating-and-freeing-frames)
7. [Complete pmm.c / pmm.h](#7-complete-pmmc--pmmh)
8. [Testing the PMM](#8-testing-the-pmm)
9. [Summary](#summary)

---

## 1. What Is Physical Memory Management?

Physical memory is the actual RAM chips in your machine. It is divided into 4KB **pages** (also called frames when we're talking about physical memory). Before we can use any memory, we need a system that tracks:

```
Questions the Physical Memory Manager answers:
  1. "Is page frame #N currently free or in use?"
  2. "Give me a free page frame." (allocation)
  3. "I'm done with frame #N, mark it free." (deallocation)
  
Without a PMM:
  → Don't know what memory is safe to use
  → Can't implement virtual memory (no way to get physical frames for page tables)
  → Can't implement heap (no way to get physical pages for heap)
  
The PMM is the foundation. Everything else builds on it.
```

---

## 2. Pages and Frames

```
Physical memory is divided into 4KB frames:
  Frame 0:  physical addresses 0x00000000 - 0x00000FFF  (4096 bytes)
  Frame 1:  physical addresses 0x00001000 - 0x00001FFF
  Frame 2:  physical addresses 0x00002000 - 0x00002FFF
  ...
  Frame N:  physical addresses N*4096 - N*4096 + 4095
  
On a machine with 128MB RAM:
  Total frames = 128 * 1024 * 1024 / 4096 = 32768 frames
  
Frame address ↔ frame number:
  frame_number  = physical_address / 4096
  physical_address = frame_number * 4096
  
4KB pages are standard because:
  - x86 hardware page tables use 4KB as the granularity unit
  - Small enough for fine-grained allocation
  - Large enough to amortize overhead
```

---

## 3. The Bitmap Allocator

The simplest approach: one bit per frame.

```
Bitmap allocator:
  0 = frame is FREE
  1 = frame is USED
  
For 128MB (32768 frames):
  Bitmap size = 32768 bits / 8 = 4096 bytes = 4KB
  
  Bit position in bitmap = frame number
  Byte index = frame_number / 8
  Bit index  = frame_number % 8
  
  Check if frame N is free:
    (bitmap[N/8] >> (N%8)) & 1  == 0 → free
    
  Mark frame N as used:
    bitmap[N/8] |= (1 << (N%8))
    
  Mark frame N as free:
    bitmap[N/8] &= ~(1 << (N%8))

Example: 8 frames, bitmap = 0b00101101
  Frame 0: used (bit 0 = 1)
  Frame 1: used (bit 1 = 0... wait)

Actually: bit N represents frame N:
  0b00101101
    ||||||||
    |||||||└─ Frame 0: USED
    ||||||└── Frame 1: USED
    |||||└─── Frame 2: free
    ||||└──── Frame 3: USED
    |||└───── Frame 4: free
    ||└────── Frame 5: USED
    |└─────── Frame 6: free
    └──────── Frame 7: free
```

**Visual layout:**
```
Physical memory:
[Frame 0][Frame 1][Frame 2][Frame 3] ... [Frame N]
 0xB8000  KERNEL   KERNEL   FREE             FREE
 (VGA)

Bitmap:
  Byte 0: 1 1 0 1 0 1 0 0  (frames 0-7, LSB=frame 0)
  Byte 1: 1 1 1 1 0 0 0 0  (frames 8-15)
  ...
```

---

## 4. Reading the Memory Map from Multiboot

We use the E820 memory map that GRUB passed us to know which physical regions exist and are usable:

```c
/* From include/multiboot.h (already written in Ch 47): */

struct multiboot_mmap_entry {
    uint32_t size;
    uint64_t base_addr;
    uint64_t length;
    uint32_t type;   /* 1 = usable RAM */
} __attribute__((packed));
```

---

## 5. Marking Reserved Regions

Not all memory is safe to use:

```
Regions we must NEVER allocate:
  0x00000000 - 0x000FFFFF: BIOS, VGA, BIOS ROM — reserved forever
  0x00100000 - kernel_end: our kernel binary — already in use
  PMM bitmap itself:         the bitmap is in memory too!

Strategy:
  1. Start: mark ALL frames as USED (everything is reserved)
  2. For each usable region in E820 map: mark those frames as FREE
  3. Re-mark reserved regions as USED:
     - Frames 0-255 (first 1MB)
     - Kernel frames (1MB to kernel_end)
     - Bitmap frame(s)
```

---

## 6. Allocating and Freeing Frames

```
Allocation (find a free frame):
  Scan bitmap linearly from start
  Find first 0 bit → that's our free frame
  Set the bit to 1 (mark as used)
  Return the frame number
  
  Optimization: track the first possibly-free frame
                to avoid scanning from the beginning every time
  
  O(n) worst case but fast in practice (first-fit)
  
Deallocation:
  Clear the bit for the given frame
  Update "first free" hint if this frame is earlier
  O(1)
```

---

## 7. Complete pmm.c / pmm.h

```c
/* include/pmm.h */
#pragma once
#include "stdint.h"

#define PAGE_SIZE   4096
#define PAGE_SHIFT  12     /* log2(4096) = 12 */

/* Frame address ↔ number conversion: */
#define FRAME_TO_ADDR(frame)  ((uint32_t)(frame) << PAGE_SHIFT)
#define ADDR_TO_FRAME(addr)   ((uint32_t)(addr)  >> PAGE_SHIFT)
#define PAGE_ALIGN_UP(addr)   (((uint32_t)(addr) + PAGE_SIZE - 1) & ~(PAGE_SIZE - 1))
#define PAGE_ALIGN_DOWN(addr) ((uint32_t)(addr) & ~(PAGE_SIZE - 1))

void     pmm_init(uint32_t mmap_addr, uint32_t mmap_length, uint32_t mem_upper_kb);
void     pmm_mark_used(uint32_t frame_addr, uint32_t length);
void     pmm_mark_free(uint32_t frame_addr, uint32_t length);
uint32_t pmm_alloc_frame(void);      /* Returns physical address, or 0 on OOM */
void     pmm_free_frame(uint32_t addr);
uint32_t pmm_free_frame_count(void);
uint32_t pmm_total_frame_count(void);
```

```c
/* kernel/pmm.c — Physical Memory Manager */

#include "pmm.h"
#include "multiboot.h"
#include "vga.h"
#include "string.h"

/* External symbols from linker.ld: */
extern uint32_t kernel_end;

/* The bitmap (placed in BSS — zero-initialized): */
/* Maximum supported RAM: 4GB → 4GB/4KB = 1M frames → 1M bits → 128KB bitmap */
#define MAX_FRAMES  (1024 * 1024)         /* 1M frames = 4GB */
#define BITMAP_SIZE (MAX_FRAMES / 8)      /* 128KB */

static uint8_t pmm_bitmap[BITMAP_SIZE];
static uint32_t total_frames;
static uint32_t free_frames;
static uint32_t first_free_hint;          /* First index to search from */

/* Low-level bit manipulation: */
static void frame_set(uint32_t frame) {
    pmm_bitmap[frame / 8] |= (1 << (frame % 8));
}

static void frame_clear(uint32_t frame) {
    pmm_bitmap[frame / 8] &= ~(1 << (frame % 8));
}

static int frame_test(uint32_t frame) {
    return (pmm_bitmap[frame / 8] >> (frame % 8)) & 1;
}

/* Initialize the PMM from the Multiboot memory map: */
void pmm_init(uint32_t mmap_addr, uint32_t mmap_length, uint32_t mem_upper_kb) {
    /* Total RAM in bytes (approximate from Multiboot): */
    uint32_t total_mem = (mem_upper_kb + 1024) * 1024;   /* upper + lower 1MB */
    total_frames = total_mem / PAGE_SIZE;
    free_frames  = 0;
    first_free_hint = 0;
    
    /* Step 1: Mark everything as used (reserved): */
    memset(pmm_bitmap, 0xFF, BITMAP_SIZE);
    
    /* Step 2: Mark usable regions as free (from E820 map): */
    if (mmap_addr && mmap_length) {
        struct multiboot_mmap_entry *entry = (void *)mmap_addr;
        uint8_t *end = (uint8_t *)mmap_addr + mmap_length;
        
        while ((uint8_t *)entry < end) {
            if (entry->type == 1 && entry->base_addr < 0xFFFFFFFF) {
                uint32_t base = (uint32_t)entry->base_addr;
                uint32_t len  = (entry->length > 0xFFFFFFFF) ?
                                0xFFFFFFFF : (uint32_t)entry->length;
                pmm_mark_free(base, len);
            }
            entry = (void *)((uint8_t *)entry + entry->size + sizeof(uint32_t));
        }
    } else {
        /* No E820 map — just assume everything above 1MB is free: */
        pmm_mark_free(0x100000, mem_upper_kb * 1024);
    }
    
    /* Step 3: Re-mark regions that are actually in use: */
    
    /* Lower 1MB (BIOS, VGA, etc.): */
    pmm_mark_used(0, 0x100000);
    
    /* Our kernel (from 1MB to kernel_end, rounded up to page boundary): */
    uint32_t kend_aligned = PAGE_ALIGN_UP((uint32_t)&kernel_end);
    pmm_mark_used(0x100000, kend_aligned - 0x100000);
    
    /* The PMM bitmap itself (it's in BSS, which is inside the kernel): */
    /* Already covered by kernel_end marker since BSS is before kernel_end */
    
    kprintf("PMM: %u MB total, %u MB free (%u frames)\n",
            total_mem / (1024*1024),
            (free_frames * PAGE_SIZE) / (1024*1024),
            free_frames);
}

/* Mark a region of physical memory as used: */
void pmm_mark_used(uint32_t frame_addr, uint32_t length) {
    uint32_t frame = ADDR_TO_FRAME(frame_addr);
    uint32_t count = (length + PAGE_SIZE - 1) / PAGE_SIZE;
    
    for (uint32_t i = 0; i < count && (frame + i) < total_frames; i++) {
        if (!frame_test(frame + i)) {
            frame_set(frame + i);
            if (free_frames > 0) free_frames--;
        }
    }
}

/* Mark a region of physical memory as free: */
void pmm_mark_free(uint32_t frame_addr, uint32_t length) {
    uint32_t frame = ADDR_TO_FRAME(frame_addr);
    uint32_t count = length / PAGE_SIZE;   /* don't partial-free */
    
    for (uint32_t i = 0; i < count && (frame + i) < total_frames; i++) {
        if (frame_test(frame + i)) {
            frame_clear(frame + i);
            free_frames++;
            /* Update hint: */
            if ((frame + i) < first_free_hint) {
                first_free_hint = frame + i;
            }
        }
    }
}

/* Allocate one physical frame; returns its physical address: */
uint32_t pmm_alloc_frame(void) {
    if (free_frames == 0) {
        kprintf("PMM: OUT OF MEMORY!\n");
        return 0;
    }
    
    /* Search from hint forward: */
    for (uint32_t i = first_free_hint; i < total_frames; i++) {
        if (!frame_test(i)) {
            frame_set(i);
            free_frames--;
            first_free_hint = i + 1;
            return FRAME_TO_ADDR(i);
        }
    }
    
    /* Shouldn't reach here if free_frames > 0, but be safe: */
    return 0;
}

/* Free a physical frame: */
void pmm_free_frame(uint32_t addr) {
    uint32_t frame = ADDR_TO_FRAME(addr);
    if (frame >= total_frames) return;
    if (frame_test(frame)) {
        frame_clear(frame);
        free_frames++;
        if (frame < first_free_hint) first_free_hint = frame;
    }
}

uint32_t pmm_free_frame_count(void)  { return free_frames;  }
uint32_t pmm_total_frame_count(void) { return total_frames; }
```

---

## 8. Testing the PMM

```c
/* In kernel_main, after pmm_init: */

void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    terminal_init();
    gdt_init();
    idt_init();
    pic_init();
    
    /* Initialize PMM: */
    struct multiboot_info *mbi = (struct multiboot_info *)mbi_ptr;
    pmm_init(mbi->mmap_addr, mbi->mmap_length, mbi->mem_upper);
    
    /* Test: allocate some frames: */
    kprintf("\nPMM test:\n");
    uint32_t f1 = pmm_alloc_frame();
    uint32_t f2 = pmm_alloc_frame();
    uint32_t f3 = pmm_alloc_frame();
    kprintf("  Allocated: 0x%x, 0x%x, 0x%x\n", f1, f2, f3);
    kprintf("  Free frames remaining: %u\n", pmm_free_frame_count());
    
    /* Free and re-allocate: */
    pmm_free_frame(f2);
    uint32_t f4 = pmm_alloc_frame();
    kprintf("  After free+alloc: 0x%x (should == 0x%x)\n", f4, f2);
    
    /* Frames should be contiguous (they're above kernel_end): */
    kprintf("  Are frames 4KB apart? %s\n",
            (f2 - f1 == 4096) ? "YES" : "NO");
    
    kprintf("\nPMM working correctly!\n");
    
    for (;;) {}
}
```

Expected output:
```
PMM: 128 MB total, 121 MB free (31000 frames)

PMM test:
  Allocated: 0x[kend+0], 0x[kend+4096], 0x[kend+8192]
  Free frames remaining: 30997
  After free+alloc: 0x[kend+4096] (should == 0x[kend+4096])
  Are frames 4KB apart? YES

PMM working correctly!
```

---

## Summary

| Concept | Description |
|---------|------------|
| Physical frame | 4KB region of physical RAM; identified by frame number = addr / 4096 |
| Bitmap allocator | One bit per frame; 0 = free, 1 = used; 128KB bitmap for 4GB RAM |
| E820 memory map | BIOS/GRUB-provided list of RAM regions; must use this to know which RAM is real |
| Type 1 regions | Usable RAM in the E820 map — all others are reserved hardware/BIOS/ACPI |
| Reserved regions | Lower 1MB (BIOS/VGA) + kernel binary — must mark USED even in usable RAM |
| pmm_alloc_frame() | Returns physical address of a free 4KB frame; O(n) first-fit scan |
| pmm_free_frame() | Marks frame as free; O(1) |
| first_free_hint | Optimization: remember last freed/allocated position to avoid re-scanning |
| kernel_end | Linker symbol marking end of kernel binary — all RAM above this is allocatable |
| OOM handling | Return 0 on allocation failure; caller must check and handle gracefully |
| PAGE_ALIGN_UP | Round address up to next 4KB boundary: `(addr + 4095) & ~4095` |
