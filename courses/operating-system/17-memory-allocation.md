# Chapter 17: Memory Allocation — malloc and free

> **"malloc() is one of the most important and most misunderstood functions in programming. It doesn't just 'give you memory.' It coordinates with the OS to acquire pages, maintains a complex free-list structure, and manages fragmentation — all in microseconds. Understanding how it works makes you a better programmer."**

---

## Table of Contents

1. [The Allocation Problem](#1-the-allocation-problem)
2. [How malloc() Works Internally](#2-how-malloc-works-internally)
3. [Allocation Strategies](#3-allocation-strategies)
4. [The Heap — Structure and Management](#4-the-heap--structure-and-management)
5. [brk() and mmap() — Getting Memory from the OS](#5-brk-and-mmap--getting-memory-from-the-os)
6. [Free — Returning Memory](#6-free--returning-memory)
7. [Memory Leaks](#7-memory-leaks)
8. [Common Memory Errors](#8-common-memory-errors)
9. [Kernel Memory Allocation — kmalloc and the Slab Allocator](#9-kernel-memory-allocation--kmalloc-and-the-slab-allocator)
10. [Custom Allocators](#10-custom-allocators)
11. [Summary](#summary)

---

## 1. The Allocation Problem

At any moment during program execution, you need to:
- Allocate a block of N bytes
- Use it for some time
- Free it when done

The allocator must satisfy requests of arbitrary sizes and in arbitrary order. It doesn't know in advance what sizes will be requested or when blocks will be freed.

**The fundamental challenge:** After many allocs and frees, memory looks like Swiss cheese — many small holes. Large requests may fail even though total free memory is sufficient.

```
After many alloc/free cycles:
[free: 10B][used][free: 5B][used][free: 20B][used][free: 8B][used]
Request: 30B → FAILS! Total free = 43B, but largest contiguous = 20B
```

A good allocator minimizes fragmentation while being fast.

---

## 2. How malloc() Works Internally

`malloc()` in glibc (the standard C library on Linux) uses a sophisticated allocator called **ptmalloc** (based on dlmalloc by Doug Lea).

**At a high level:**
1. Maintain a pool of free memory blocks (the heap)
2. When `malloc(n)` is called: find a free block ≥ n bytes
3. If no free block is large enough: ask OS for more memory (`brk()` or `mmap()`)
4. Carve out the block, return pointer to the user
5. When `free(ptr)` is called: put the block back in the free pool

**Block structure:**
Each heap block (both used and free) has a hidden header:
```
[  block header  ][ ... user data ... ][optional footer]
  prev_size (4B)  ←── user gets pointer to HERE
  size+flags (4B)     
```

The block header stores:
- **size:** total size of the block (including header) — in multiples of 8 bytes
- **prev_size:** size of the previous block (used for backward merging)
- **flags:** P (previous block in use), M (mmap'd block), A (arena)

When you call `malloc(n)` and get pointer `p`, the actual block is at `p - 8` (or `p - 16`). `free(p)` uses `p - 8` to find the header.

---

## 3. Allocation Strategies

Different ways to find a free block:

**First Fit:**
Scan the free list from the beginning; return the first block that's large enough.
```
Free list: [10B][50B][25B][8B][100B]
Request: 20B → returns [50B] (split: [20B used][30B free])
```
- Fast to implement
- Creates fragmentation near the beginning of the heap

**Best Fit:**
Find the smallest free block that's still large enough.
```
Free list: [10B][50B][25B][8B][100B]
Request: 20B → returns [25B] (best fit: [20B used][5B free])
```
- Leaves larger blocks free (better for future large requests)
- Leaves many tiny fragments (5B block rarely useful)

**Worst Fit:**
Find the LARGEST free block.
- Theory: large blocks split evenly, less tiny unusable fragments
- In practice: worse than best fit

**Next Fit:**
Like First Fit, but start scanning from where the last allocation was made.
- More evenly distributed fragmentation
- Slightly faster than First Fit (less scanning from start)

**Modern allocators use segregated free lists:**
Instead of one free list, maintain separate lists for different size classes:
```
Size class 16B:  [block]─[block]─[block]─...
Size class 32B:  [block]─[block]─...
Size class 64B:  [block]─...
Size class 128B: [block]─...
...
Large (>512B):   best-fit within a sorted tree
```
`malloc(20)` → look in the 32B list (smallest class ≥ 20). O(1) if list is non-empty!

---

## 4. The Heap — Structure and Management

The **heap** is a contiguous region of memory starting just after the program's BSS segment. It grows upward.

```
Low address
┌───────────────────────┐ ← heap start (static, set at load time)
│                       │
│  Heap region          │
│  (grows upward ↑)     │
│                       │
│  [block1][block2]...  │
│                       │
├───────────────────────┤ ← "program break" (brk pointer)
│  (unmapped space)     │  OS doesn't allocate physical pages here yet
│                       │
High address
```

**The program break (brk):**
The program break is the top of the heap. The allocator extends the heap by calling `brk()` (to move the break up) to get more memory from the OS.

**Chunk coalescing:**
When two adjacent free blocks exist, merge them into one larger free block:
```
[free: 32B][free: 64B]  →  [free: 96B]
```
glibc does this eagerly (on every `free()` call, checks neighbors).

**Thread caching (tcmalloc, jemalloc):**
In multi-threaded programs, every thread allocating and freeing from the same heap requires locking. Modern allocators (TCMalloc by Google, jemalloc by Facebook) give each thread its own cache:
- Each thread has a thread-local cache of small blocks (no lock needed)
- Falls back to central heap (with lock) only when thread cache is empty/full
- Dramatically reduces lock contention in multi-threaded programs

---

## 5. brk() and mmap() — Getting Memory from the OS

`malloc()` gets pages from the OS via two system calls:

**brk(addr):**
Extends the heap by moving the program break.
```c
// Current heap ends at 0x10000
brk(0x20000);  // extend heap to 0x20000 (add 64KB)
// Now malloc can use addresses 0x10000 – 0x20000
```

Used for small allocations (glibc uses this for the main heap).

**mmap(NULL, size, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0):**
Maps a new anonymous region of memory anywhere in the address space.
- Not necessarily contiguous with the existing heap
- Can be unmapped individually with `munmap()`

glibc uses `mmap()` for:
- Large allocations (>128KB by default): each gets its own mmap region
- Thread-specific heaps (each thread gets its own mmap'd arena)

```c
// What malloc actually does internally for a large request:
void *chunk = mmap(NULL, 200000,
                   PROT_READ | PROT_WRITE,
                   MAP_PRIVATE | MAP_ANONYMOUS,
                   -1, 0);
// free() of this block calls munmap() — returns pages to OS
```

**Why does this matter?**
For small allocations, `free()` returns the memory to glibc's free list — it does NOT immediately give pages back to the OS. The process's RSS (Resident Set Size) stays high.

For large (`mmap`'d) allocations, `free()` calls `munmap()` — pages go back to OS immediately.

You can force glibc to return pages: `malloc_trim(0)` trims the heap, returning unused pages to OS.

---

## 6. Free — Returning Memory

```c
void free(void *ptr) {
    if (ptr == NULL) return;  // free(NULL) is always safe
    
    // Get the block header (8 bytes before ptr)
    Block *blk = (Block *)((char *)ptr - HEADER_SIZE);
    
    // Mark as free
    blk->flags &= ~IN_USE;
    
    // Coalesce with adjacent free blocks
    if (!prev_block_in_use(blk))
        blk = merge_with_prev(blk);
    if (!next_block_in_use(blk))
        blk = merge_with_next(blk);
    
    // Return to appropriate free list
    add_to_free_list(blk);
}
```

**What `free()` does NOT do:**
- Does NOT zero the memory (a common misconception)
- Does NOT immediately return pages to OS (for small blocks)
- Does NOT prevent use-after-free (the memory might look valid for a while)

---

## 7. Memory Leaks

A **memory leak** occurs when allocated memory is never freed. The program's memory usage grows indefinitely.

```c
// Leak example:
while (1) {
    char *buf = malloc(1024);
    // ... use buf ...
    // forgot free(buf)!
    // Every iteration: 1KB lost forever
}
// After 1 hour: program consumed 3.6GB
```

**Finding leaks:**

**Valgrind:**
```bash
valgrind --leak-check=full --show-leak-kinds=all ./program

# Output:
# LEAK SUMMARY:
# definitely lost: 4,096 bytes in 1 blocks
# indirectly lost: 0 bytes in 0 blocks
#   still reachable: 1,024 bytes in 1 blocks
# ==1234== at 0x483577F: malloc
# ==1234== by 0x401234: create_buffer (program.c:15)
# ==1234== by 0x401289: main (program.c:42)
```

**AddressSanitizer (ASan):**
```bash
gcc -fsanitize=address -o program program.c
./program
# Reports: heap-use-after-free, heap-buffer-overflow, memory leaks
```

**LeakSanitizer:**
```bash
gcc -fsanitize=leak -o program program.c
./program
# Reports memory leaks at program exit
```

**RAII (C++) — automatic memory management:**
```cpp
// With raw pointers: leak-prone
{
    int *arr = new int[100];
    if (error) return;  // BUG: arr leaked!
    delete[] arr;
}

// With RAII (smart pointers):
{
    auto arr = std::make_unique<int[]>(100);
    if (error) return;  // OK: arr automatically deleted when scope exits
}
```

---

## 8. Common Memory Errors

**1. Buffer overflow:**
Writing beyond the end of an allocated buffer:
```c
char buf[10];
strcpy(buf, "Hello, World!");  // 14 chars > 10 → overflow!
// Overwrites whatever is adjacent in memory
// Classic security vulnerability (stack smashing)
```

**2. Use-after-free:**
Using memory after it's been freed:
```c
int *p = malloc(4);
*p = 42;
free(p);
printf("%d\n", *p);  // undefined behavior! p is now a dangling pointer
*p = 100;            // might corrupt malloc's free list!
```

**3. Double free:**
Freeing the same pointer twice:
```c
int *p = malloc(4);
free(p);
free(p);  // double free → heap corruption → security vulnerability
```

**4. Null pointer dereference:**
```c
int *p = malloc(4);
if (p == NULL) { /* handle error */ }
// If you skip the check:
*p = 42;  // crash if malloc returned NULL (out of memory)
```

**5. Stack overflow:**
```c
void recursive(int n) {
    char big_array[1000000];  // 1MB per call
    recursive(n + 1);         // infinite recursion
    // → stack grows until SIGSEGV (stack overflow)
}
```

**6. Heap corruption:**
Any of the above bugs can corrupt the heap's metadata (block headers, free list pointers), causing subsequent malloc/free to crash or behave incorrectly.

---

## 9. Kernel Memory Allocation — kmalloc and the Slab Allocator

The kernel needs its own memory allocator (it can't use glibc's malloc — that's user space).

**`kmalloc(size, flags)`:**
The kernel's equivalent of `malloc()`:
```c
// Allocate 4KB of kernel memory
struct my_struct *s = kmalloc(sizeof(struct my_struct), GFP_KERNEL);
if (!s) return -ENOMEM;

// Use it
s->value = 42;

// Free it
kfree(s);
```

**GFP flags:** Control behavior:
- `GFP_KERNEL`: Normal allocation, may sleep (can't use in interrupt context)
- `GFP_ATOMIC`: Can't sleep (for interrupt handlers, spinlock-protected code)
- `GFP_USER`: User-space allocation
- `GFP_DMA`: Must be in DMA-accessible memory (<16MB on legacy ISA DMA)

**The Slab Allocator:**
The kernel frequently allocates and frees many objects of the same type (task_struct, inode, file, socket, etc.). The slab allocator creates dedicated caches for each:

```c
// Register a slab cache for task_struct:
struct kmem_cache *task_cache = 
    kmem_cache_create("task_struct",           // name
                       sizeof(struct task_struct),
                       __alignof__(struct task_struct),
                       SLAB_HWCACHE_ALIGN,     // flags
                       NULL);                  // constructor

// Allocate from the cache:
struct task_struct *t = kmem_cache_alloc(task_cache, GFP_KERNEL);

// Free back to the cache:
kmem_cache_free(task_cache, t);
```

**Why slabs are faster:**
- Objects in the cache are pre-initialized (constructor called once on first allocation)
- The freed object stays in the cache — next allocation is instant (just take from cache, no initialization needed for common fields)
- Objects in the same cache are on the same memory pages (cache locality)

**Slab allocator variants in Linux:**
- SLAB: original allocator
- SLUB: simplified slab (default since 2.6.23) — fewer per-CPU data structures, better debugging
- SLOB: ultra-compact allocator for embedded systems

---

## 10. Custom Allocators

For specific use cases, a custom allocator can beat glibc:

**Arena/pool allocator:**
All objects in a pool are the same size. Allocation = take from pool. Free = return to pool. No fragmentation possible.

```c
#define POOL_SIZE 1000

struct Object pool[POOL_SIZE];  // pre-allocated array
int free_list[POOL_SIZE];       // indices of free objects
int free_count = POOL_SIZE;

struct Object *alloc_object() {
    if (free_count == 0) return NULL;
    return &pool[free_list[--free_count]];
}

void free_object(struct Object *obj) {
    free_list[free_count++] = obj - pool;
}
```

Used in: database connection pools, network packet buffers, game object managers.

**Stack/linear/bump allocator:**
For objects that are all freed together:
```c
char arena[1024*1024];
size_t offset = 0;

void *alloc(size_t size) {
    void *p = &arena[offset];
    offset += size;  // just bump a pointer!
    return p;
}

void free_all() {
    offset = 0;  // free everything at once!
}
```

Allocation is O(1): just increment a counter. Free is O(1): reset the counter.
Used in: HTTP request handling (allocate per-request, free all at end), parser memory.

---

## Summary

| Topic | Key Points |
|-------|-----------|
| malloc() | Finds a free block ≥ requested size; asks OS for more if needed |
| Block header | Hidden metadata before each allocation: size, flags, prev_size |
| Free list | Linked list of free blocks; malloc searches this |
| Best fit | Finds smallest sufficient block; reduces fragmentation |
| First fit | Finds first sufficient block; faster but more fragmentation |
| Segregated lists | Separate free lists per size class; O(1) allocation |
| brk() | Extends heap for small allocations |
| mmap() | Gets large anonymous regions; freed immediately on free() |
| Memory leak | Allocated memory never freed; process grows indefinitely |
| Valgrind/ASan | Tools to detect memory errors and leaks |
| kmalloc | Kernel's malloc; uses GFP flags for different contexts |
| Slab allocator | Pre-allocate same-type objects; fast, cache-friendly |
