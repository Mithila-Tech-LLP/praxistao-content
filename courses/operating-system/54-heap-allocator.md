# Chapter 54: Heap Allocator

> **"The PMM gives us 4KB blocks of physical RAM. The VMM gives us virtual address space. But neither gives us malloc-style allocations — the ability to ask for 37 bytes or 10,000 bytes and get exactly that. The heap allocator bridges the gap: it manages a large virtual region, sub-allocates from it in arbitrary sizes, and handles fragmentation."**

---

## Table of Contents

1. [Why a Heap Allocator?](#1-why-a-heap-allocator)
2. [Heap Layout in Memory](#2-heap-layout-in-memory)
3. [Free List Allocator — Block Headers](#3-free-list-allocator--block-headers)
4. [Splitting and Coalescing](#4-splitting-and-coalescing)
5. [kmalloc and kfree](#5-kmalloc-and-kfree)
6. [Expanding the Heap (sbrk-style)](#6-expanding-the-heap-sbrk-style)
7. [Complete heap.c / heap.h](#7-complete-heaph)
8. [Testing the Heap](#8-testing-the-heap)
9. [Summary](#summary)

---

## 1. Why a Heap Allocator?

```
Without a heap:
  pmm_alloc_frame()  → always gives you exactly 4096 bytes
  
  Want 100 bytes? You get 4096 (3996 wasted per allocation)
  Want 10,000 bytes? You need to manage multiple frames manually
  
With a heap (kmalloc/kfree):
  kmalloc(100)  → 100 bytes (small overhead for bookkeeping)
  kmalloc(10000) → 10,000 bytes (internal fragmentation: ~10,000 bytes used)
  kfree(ptr)    → returns memory for reuse
  
Use cases:
  Kernel data structures of variable size (process descriptors, file handles)
  String buffers for path names, command lines
  Any allocation where size isn't known at compile time
```

---

## 2. Heap Layout in Memory

```
Kernel heap: a contiguous virtual address region the heap allocator manages.

Virtual memory map:
  0x00000000 - 0x000FFFFF: identity-mapped low memory (VGA, BIOS)
  0x00100000 - kernel_end: kernel binary (text, data, bss, stack)
  0x00400000 - 0x00FFFFFF: kernel heap (we'll put it here — 12MB)
  
The heap starts as a small region (one or a few 4KB pages).
As kmalloc needs more memory, we extend the heap by mapping more pages.
The current "top" of the heap is tracked by heap_top.

Initial state:
  heap_start = 0x00400000
  heap_top   = 0x00400000 (empty — no pages mapped yet)
  heap_end   = 0x00FFFFFF (maximum we'll ever grow to)
```

---

## 3. Free List Allocator — Block Headers

Every allocated (and free) block has a small header describing it:

```
Block layout in memory:
  ┌──────────────┬──────────────────────────────────┐
  │   Header     │   User data (or free list links) │
  └──────────────┴──────────────────────────────────┘
  8 bytes          'size' bytes

Header structure:
  magic:  0xDEADBEEF if used, 0xFEEDFACE if free (for debugging)
  size:   size of the USER DATA (not including header)
  used:   1 = allocated, 0 = free
  next:   pointer to next block (for free blocks: next free block)
  prev:   pointer to previous block (for coalescing)
  
The heap is one big list of blocks, both free and used:
  [hdr][data][hdr][data][hdr][data]...
  
Visualized:
  +─────────────────────────────────────────────+
  │ HDR(32B,used) │ 32 bytes data               │ ← allocated
  │ HDR(64B,free) │ 64 bytes free               │ ← free
  │ HDR(16B,used) │ 16 bytes data               │ ← allocated
  │ HDR(100B,free)│ 100 bytes free              │ ← free
  +─────────────────────────────────────────────+
```

---

## 4. Splitting and Coalescing

**Splitting** (when allocating):
If a free block is larger than needed, split it into two:
```
Before: [HDR(100B,free) | 100 bytes]
Request: kmalloc(30)
After:  [HDR(30B,used)  | 30 bytes][HDR(62B,free) | 62 bytes]
        (8 bytes for the new header)
        Only split if remaining free portion is large enough (e.g., > 16 bytes)
```

**Coalescing** (when freeing):
When freeing a block, merge adjacent free blocks to avoid fragmentation:
```
Before: [HDR(30B,free)][HDR(62B,free)][HDR(16B,used)]
After coalesce: [HDR(100B,free)][HDR(16B,used)]

Why important:
  Without coalescing: lots of small free blocks
  Can't satisfy a large allocation even though total free > request
  This is "external fragmentation"
```

---

## 5. kmalloc and kfree

```c
/* kernel/heap.c */

#include "heap.h"
#include "vmm.h"
#include "pmm.h"
#include "string.h"
#include "vga.h"

/* Heap boundaries: */
#define HEAP_START    0x00800000   /* 8MB virtual address */
#define HEAP_MAX      0x01000000   /* 16MB virtual address (8MB heap) */
#define HEAP_MAGIC_USED 0xDEADBEEF
#define HEAP_MAGIC_FREE 0xFEEDFACE
#define HEAP_MIN_SPLIT  16          /* Minimum split remainder */

typedef struct heap_block {
    uint32_t magic;
    uint32_t size;           /* Size of user data */
    uint8_t  used;
    struct heap_block *next; /* Next block in the list */
    struct heap_block *prev; /* Previous block */
} heap_block_t;

#define HEADER_SIZE sizeof(heap_block_t)

static heap_block_t *heap_start_block = NULL;
static uint32_t      heap_current_top;

/* Extern page directory (from vmm.c): */
extern uint32_t kernel_page_dir[];

/* Expand the heap by mapping more physical pages: */
static void *heap_sbrk(uint32_t bytes) {
    uint32_t old_top = heap_current_top;
    uint32_t new_top = old_top + bytes;
    
    if (new_top > HEAP_MAX) {
        kprintf("heap: out of virtual address space!\n");
        return NULL;
    }
    
    /* Map each new page: */
    for (uint32_t addr = old_top; addr < new_top; addr += PAGE_SIZE) {
        uint32_t phys = pmm_alloc_frame();
        if (!phys) {
            kprintf("heap: out of physical memory!\n");
            return NULL;
        }
        vmm_map_page(kernel_page_dir, addr, phys, PTE_WRITABLE);
    }
    
    heap_current_top = new_top;
    return (void *)old_top;
}

/* Initialize the heap: */
void heap_init(void) {
    heap_current_top = HEAP_START;
    
    /* Map one initial page: */
    void *p = heap_sbrk(PAGE_SIZE);
    if (!p) {
        kprintf("heap_init: failed to map initial page!\n");
        return;
    }
    
    /* Create one large free block covering the initial page: */
    heap_start_block = (heap_block_t *)HEAP_START;
    heap_start_block->magic = HEAP_MAGIC_FREE;
    heap_start_block->size  = PAGE_SIZE - HEADER_SIZE;
    heap_start_block->used  = 0;
    heap_start_block->next  = NULL;
    heap_start_block->prev  = NULL;
    
    kprintf("Heap initialized at 0x%x, initial size: %u bytes\n",
            HEAP_START, PAGE_SIZE - (uint32_t)HEADER_SIZE);
}

/* Allocate 'size' bytes from the heap: */
void *kmalloc(uint32_t size) {
    if (size == 0) return NULL;
    
    /* Align to 8 bytes for performance: */
    size = (size + 7) & ~7;
    
    /* First-fit search: find first free block large enough: */
    heap_block_t *block = heap_start_block;
    while (block) {
        if (!block->used && block->size >= size) {
            /* Found a suitable free block: */
            
            /* Should we split it? */
            if (block->size >= size + HEADER_SIZE + HEAP_MIN_SPLIT) {
                /* Split: create new free block for remainder: */
                heap_block_t *new_block = (heap_block_t *)((uint8_t *)block
                                          + HEADER_SIZE + size);
                new_block->magic = HEAP_MAGIC_FREE;
                new_block->size  = block->size - size - HEADER_SIZE;
                new_block->used  = 0;
                new_block->next  = block->next;
                new_block->prev  = block;
                
                if (block->next) block->next->prev = new_block;
                block->next = new_block;
                block->size = size;
            }
            
            block->magic = HEAP_MAGIC_USED;
            block->used  = 1;
            return (void *)((uint8_t *)block + HEADER_SIZE);
        }
        block = block->next;
    }
    
    /* No free block found — expand the heap: */
    uint32_t expand_size = (size + HEADER_SIZE + PAGE_SIZE - 1) & ~(PAGE_SIZE - 1);
    heap_block_t *new_block = heap_sbrk(expand_size);
    if (!new_block) return NULL;  /* OOM */
    
    new_block->magic = HEAP_MAGIC_FREE;
    new_block->size  = expand_size - HEADER_SIZE;
    new_block->used  = 0;
    new_block->next  = NULL;
    
    /* Append to end of list: */
    block = heap_start_block;
    while (block->next) block = block->next;
    new_block->prev = block;
    block->next = new_block;
    
    /* Now allocate from this new block: */
    return kmalloc(size);   /* Recursion — guaranteed to succeed now */
}

/* Free memory: */
void kfree(void *ptr) {
    if (!ptr) return;
    
    heap_block_t *block = (heap_block_t *)((uint8_t *)ptr - HEADER_SIZE);
    
    /* Sanity check: */
    if (block->magic != HEAP_MAGIC_USED) {
        kprintf("kfree: double free or corruption at 0x%x!\n", (uint32_t)ptr);
        return;
    }
    
    block->magic = HEAP_MAGIC_FREE;
    block->used  = 0;
    
    /* Coalesce with NEXT free block: */
    if (block->next && !block->next->used) {
        heap_block_t *next = block->next;
        block->size += HEADER_SIZE + next->size;
        block->next  = next->next;
        if (next->next) next->next->prev = block;
    }
    
    /* Coalesce with PREVIOUS free block: */
    if (block->prev && !block->prev->used) {
        heap_block_t *prev = block->prev;
        prev->size += HEADER_SIZE + block->size;
        prev->next  = block->next;
        if (block->next) block->next->prev = prev;
    }
}

/* Allocate and zero-initialize: */
void *kcalloc(uint32_t count, uint32_t size) {
    void *p = kmalloc(count * size);
    if (p) memset(p, 0, count * size);
    return p;
}

/* Reallocate with a new size: */
void *krealloc(void *ptr, uint32_t new_size) {
    if (!ptr) return kmalloc(new_size);
    if (!new_size) { kfree(ptr); return NULL; }
    
    heap_block_t *block = (heap_block_t *)((uint8_t *)ptr - HEADER_SIZE);
    if (block->size >= new_size) return ptr;  /* already big enough */
    
    void *new_ptr = kmalloc(new_size);
    if (!new_ptr) return NULL;
    memcpy(new_ptr, ptr, block->size);
    kfree(ptr);
    return new_ptr;
}
```

---

## 6. Expanding the Heap (sbrk-style)

The `heap_sbrk` function (already shown above) expands the heap by:
1. Calculating how many new pages are needed
2. Calling `pmm_alloc_frame()` for physical frames
3. Calling `vmm_map_page()` to create virtual mappings
4. Updating `heap_current_top`

This is analogous to Unix's `sbrk()` system call, which moves the program break upward.

---

## 7. Complete heap.h

```c
/* include/heap.h */
#pragma once
#include "stdint.h"

void  heap_init(void);
void *kmalloc(uint32_t size);
void  kfree(void *ptr);
void *kcalloc(uint32_t count, uint32_t size);
void *krealloc(void *ptr, uint32_t new_size);
```

---

## 8. Testing the Heap

```c
/* kernel_main after vmm_init(): */

heap_init();

/* Basic allocation: */
char *str = (char *)kmalloc(32);
const char *hello = "Hello, heap!";
for (int i = 0; hello[i]; i++) str[i] = hello[i];
str[12] = '\0';
kprintf("kmalloc test: '%s'\n", str);

/* Multiple allocations: */
uint32_t *arr = (uint32_t *)kmalloc(10 * sizeof(uint32_t));
for (int i = 0; i < 10; i++) arr[i] = i * i;
kprintf("Array: %u %u %u %u %u\n", arr[0], arr[1], arr[2], arr[3], arr[4]);

/* Free and reallocate: */
kfree(str);
char *str2 = (char *)kmalloc(32);  /* Should reuse the freed block */
kprintf("Reused block: 0x%x (was 0x%x)\n", (uint32_t)str2, (uint32_t)str);

/* Large allocation: */
void *big = kmalloc(8000);
kprintf("Large alloc at: 0x%x\n", (uint32_t)big);
kfree(big);

/* kcalloc: */
uint8_t *buf = (uint8_t *)kcalloc(256, 1);
int all_zero = 1;
for (int i = 0; i < 256; i++) if (buf[i] != 0) all_zero = 0;
kprintf("kcalloc zeroed: %s\n", all_zero ? "YES" : "NO");
kfree(buf);
kfree(arr);

kprintf("Heap tests passed!\n");
```

Expected output:
```
Heap initialized at 0x800000, initial size: 4088 bytes
kmalloc test: 'Hello, heap!'
Array: 0 1 4 9 16
Reused block: 0x800008 (was 0x800008)
Large alloc at: 0x800058
kcalloc zeroed: YES
Heap tests passed!
```

---

## Summary

| Concept | Description |
|---------|------------|
| Heap allocator | Sub-allocates within a large virtual region; handles arbitrary sizes |
| Block header | Metadata before each allocation: magic, size, used flag, next/prev pointers |
| First-fit | Scan from start, use first free block that's large enough; simple but can fragment |
| Splitting | When free block is larger than requested: split into [used][free] to avoid waste |
| Coalescing | When freeing: merge adjacent free blocks to reduce fragmentation |
| Magic numbers | 0xDEADBEEF (used) / 0xFEEDFACE (free): help detect double-free/corruption |
| heap_sbrk | Expand heap by mapping new physical pages into virtual heap region |
| 8-byte alignment | Round up allocation sizes to 8 bytes for proper struct alignment |
| kmalloc(0) | Returns NULL (standard behavior) |
| kfree(NULL) | Does nothing (standard behavior) |
| kcalloc | malloc + memset(0) — useful for zero-initialized arrays |
| krealloc | Grow an allocation: malloc new, memcpy, free old |
| External fragmentation | Lots of small free blocks; can't satisfy large request; fixed by coalescing |
| Internal fragmentation | Allocated block is larger than requested; unavoidable overhead |
