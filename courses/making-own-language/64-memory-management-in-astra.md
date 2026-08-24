# Chapter 64: Memory Management in Astra — Heap, Stack, and Garbage Collection

> "Memory management is like doing the dishes. Manual memory management means you wash them yourself — you might forget one and get mold. Garbage collection is like having a dishwasher — it runs automatically, but occasionally it runs when you are trying to eat dinner."
> — A systems programmer at a dinner party

---

## Overview

Every variable in your Astra program lives somewhere in memory. But where? For how long? Who decides when that memory can be reused? These questions sit at the heart of language design, and every language answers them differently.

C says: "You, the programmer, decide. Use `malloc` to allocate and `free` to release. If you forget, that is your problem." This leads to entire categories of bugs — memory leaks, double-frees, use-after-free — that have caused countless security vulnerabilities.

Astra says: "We will handle it. Just create objects; the garbage collector will clean them up when you are done with them."

This chapter builds Astra's garbage collector: a mark-and-sweep GC implemented in C, about 200 lines, that runs inside every Astra program and handles all memory management automatically.

**What you will understand after this chapter:**
- How stack and heap memory differ and when each is used
- Why C's manual memory management is dangerous
- How mark-and-sweep garbage collection works, step by step
- How to implement a complete GC in C
- How Astra's GC knows which pointers are "live"

---

## What We Are Building

```
ASTRA MEMORY LAYOUT
━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━━

High addresses
┌──────────────────────────────────────┐
│          KERNEL SPACE               │  ← OS, not accessible
├──────────────────────────────────────┤
│            STACK                    │  ← local variables, frames
│         (grows down ↓)              │     ~8 MB limit
│ fn foo() { let x = 5; ... }         │     automatic, LIFO
├──────────────────────────────────────┤
│     (unmapped gap for safety)        │
├──────────────────────────────────────┤
│            HEAP                     │  ← structs, strings, lists
│         (grows up ↑)                │     unlimited (until OOM)
│  gc_alloc() returns pointers here   │     GC manages lifetime
├──────────────────────────────────────┤
│        BSS SEGMENT                  │  ← global variables (zero-init)
├──────────────────────────────────────┤
│        DATA SEGMENT                 │  ← global variables (initialized)
├──────────────────────────────────────┤
│        TEXT SEGMENT                 │  ← compiled machine code
└──────────────────────────────────────┘
Low addresses (0x0)
```

---

## Table of Contents

1. Stack Memory — Fast, Automatic, Limited
2. Heap Memory — Dynamic, Flexible, Dangerous
3. malloc and free in C
4. Memory Safety Problems in C
5. Garbage Collection — The Four Main Strategies
6. Why Astra Chose Mark-and-Sweep
7. The Object Header — How the GC Tracks Objects
8. GC Roots — What the GC Must Never Collect
9. The Mark Phase — Tracing Live Objects
10. The Sweep Phase — Reclaiming Dead Objects
11. Write Barriers
12. The Complete GC Implementation
13. Integrating the GC with the Runtime
14. Astra Build Milestone
15. Exercises

---

## 1. Stack Memory — Fast, Automatic, Limited

When a function is called, the CPU pushes a **stack frame** onto the call stack. The stack frame holds:
- The function's local variables
- The function's arguments
- The return address (where to go when the function returns)
- Saved registers

Stack memory is **extremely fast** — allocating is just subtracting from the stack pointer register (one instruction). Deallocation is just adding back (also one instruction). The OS pre-maps the stack region, so there are no syscalls.

```
CALL STACK EXAMPLE

fn add(a: int, b: int) -> int {
    let result = a + b   // 'result' is on the stack
    return result
}

fn main() {
    let x = 10          // x is on the stack
    let y = 20          // y is on the stack
    let z = add(x, y)  // z is on the stack; add() gets a new frame
}

Stack (growing down):
┌─────────────────────┐  ← stack pointer (high)
│  main's frame:      │
│    x = 10           │
│    y = 20           │
│    z = (pending)    │
│    return addr      │
├─────────────────────┤  ← new frame when add() is called
│  add's frame:       │
│    a = 10           │
│    b = 20           │
│    result = 30      │
│    return addr      │
└─────────────────────┘  ← current stack pointer (low)

When add() returns: its frame is popped, stack pointer moves up.
```

**Limitations of stack memory:**
- **Size:** The default stack size is 8 MB. If you have too many nested function calls (deep recursion), you overflow the stack — "stack overflow."
- **Lifetime:** Stack variables live only as long as their function. You CANNOT return a pointer to a local variable — the frame disappears when the function returns.
- **Types:** Only fixed-size values work well. You cannot have a stack-allocated list that grows dynamically.

In Astra, the following go on the stack:
- `int`, `float`, `bool` local variables
- Small structs (if they fit in registers or a few words)
- Function arguments and return values

Everything else — structs with heap pointers, strings, lists, closures — goes on the heap.

---

## 2. Heap Memory — Dynamic, Flexible, Dangerous

The heap is a large region of memory (starting small, growing as needed by asking the OS for more via `mmap` or `brk`). You request chunks of any size, use them, and eventually mark them as free.

Unlike the stack, heap memory has **no automatic lifetime**. A heap-allocated object lives until you explicitly free it (in C) or until the GC collects it (in Astra).

**When does Astra allocate on the heap?**
- Every Astra string (`"hello"` as a runtime value with a mutable component)
- Every struct instance (`User { id: 1, name: "Alice" }`)
- Every list (`[1, 2, 3]`)
- Every closure (`fn(x) { x + 1 }` captured as a value)
- Every box/reference value

```astra
fn main() {
    let n = 42            // int — on the stack
    let s = "hello"       // string — data on heap, AstraString* on stack
    let u = User { ... }  // struct — on the heap, pointer on stack
    let lst = [1, 2, 3]   // list — on the heap, pointer on stack
}
```

---

## 3. malloc and free in C

Before we understand why garbage collection exists, we need to understand the manual approach. In C:

```c
// Allocate 100 bytes on the heap. Returns a pointer.
void* ptr = malloc(100);

// Use the memory
memset(ptr, 0, 100);

// Free the memory when done
free(ptr);
```

**How malloc works internally** (simplified):
```
malloc maintains a "free list" — a linked list of available memory blocks.

When you call malloc(100):
1. Find a free block with size >= 100 in the free list
2. Remove it from the free list
3. Return a pointer to it

When you call free(ptr):
1. Find the block that ptr points into
2. Mark it as available
3. Add it back to the free list (merging adjacent free blocks if possible)

If no block is big enough:
1. Ask the OS for more memory (mmap/sbrk)
2. Add the new region to the free list
3. Retry
```

malloc stores the block's size just before the pointer it returns:

```
What malloc returns and manages:

  ptr - 8 bytes            ptr
       ↓                   ↓
┌──────────────┬───────────────────────────────┐
│  size = 100  │    100 bytes of user memory    │
│ (hidden)     │                               │
└──────────────┴───────────────────────────────┘
               ↑ This is what you get back
```

---

## 4. Memory Safety Problems in C

C trusts programmers completely, which leads to a class of bugs that are responsible for a huge portion of security vulnerabilities. Understanding them explains WHY we build a GC.

**1. Memory leak: forgetting to call free()**
```c
void process() {
    char* buf = malloc(1024);
    if (error_condition) {
        return;  // BUG: forgot to free(buf)!
    }
    // ... use buf ...
    free(buf);
}
// Every call to process() with an error leaks 1024 bytes.
// Over time, the process uses more and more memory → crash.
```

**2. Dangling pointer: using memory after free()**
```c
char* ptr = malloc(100);
free(ptr);
ptr[0] = 'A';  // BUG: writing to freed memory!
// This can corrupt other data, cause crashes, or enable exploits.
```

**3. Double-free: calling free() twice**
```c
char* ptr = malloc(100);
free(ptr);
free(ptr);  // BUG: double free corrupts malloc's internal data!
// This can cause heap corruption or arbitrary code execution.
```

**4. Buffer overflow: writing past the end**
```c
char* buf = malloc(10);
strcpy(buf, "Hello, World!");  // BUG: 14 bytes into 10-byte buffer!
// Overwrites adjacent heap data. Classic security vulnerability.
```

Astra's GC eliminates problems 1, 2, and 3 completely. Problem 4 is prevented by bounds checking (Chapter 63).

---

## 5. Garbage Collection — The Four Main Strategies

There are four main approaches to automatic memory management. Every GC language uses some variant of these:

### Strategy 1: Reference Counting

Each object stores a count of how many pointers point to it. When the count reaches zero, the object is freed immediately.

```
Object A (count: 2) ◄── ptr1
                    ◄── ptr2

ptr1 = nil  → count drops to 1
ptr2 = nil  → count drops to 0 → FREE IMMEDIATELY
```

**Pros:** Simple to implement; memory is freed as soon as it becomes unreachable; no "stop-the-world" pauses.

**Cons:** Cannot handle cycles. If A points to B and B points to A, both have count >= 1 forever, even if nothing else points to them.

```
A (count: 1) ──► B (count: 1)
A (count: 1) ◄── B (count: 1)
     ↑                 ↑
     These will NEVER be freed!
```

Python uses reference counting + a separate cycle detector. Swift uses reference counting with `weak` references to break cycles. This is simple but requires programmer discipline.

### Strategy 2: Mark-and-Sweep

Periodically, the GC scans all memory in two phases:
1. **Mark:** Starting from "roots" (global variables, stack variables), follow all pointers and mark every reachable object.
2. **Sweep:** Scan all objects. Any unmarked object is unreachable (garbage) — free it.

```
BEFORE:     MARK:        SWEEP:
[A]●──►[B]  [A]✓──►[B]✓  [A]✓──►[B]✓
[C]●        [C]✓         [C]✓
[D]         [D]           [D → FREED]
[E]──►[F]   [E]──►[F]    [E → FREED]
                          [F → FREED]
```

**Pros:** Handles cycles naturally. B and F above would be collected even if they pointed to each other.

**Cons:** "Stop-the-world" pause while GC runs (the program cannot run during collection).

### Strategy 3: Copying / Compacting GC

Memory is divided into two halves: "from-space" and "to-space." The GC copies all live objects from from-space to to-space, compacting them together. Then the halves flip.

```
FROM-SPACE          TO-SPACE
┌───┬───┬───┬───┐   ┌───┬───┬───┬───┐
│ A │ B │ G │ C │   │   │   │   │   │
│(L)│(L)│(D)│(L)│   │   │   │   │   │
└───┴───┴───┴───┘   └───┴───┴───┴───┘
     GC copies live (L) objects:
┌───┬───┬───┬───┐   ┌───┬───┬───┬───┐
│   │   │   │   │   │ A │ B │ C │   │
│   │   │   │   │   │(L)│(L)│(L)│   │
└───┴───┴───┴───┘   └───┴───┴───┴───┘
  (all freed)          (compacted!)
```

**Pros:** Allocation is extremely fast (just bump a pointer). Eliminates heap fragmentation.

**Cons:** Requires 2x the memory. All pointers must be updated when objects move (complex).

### Strategy 4: Generational GC

Based on the empirical observation that "most objects die young" (they are created, used briefly, and become garbage):

```
                 GENERATION 0 (nursery) — collected every few ms
                 ┌─────────────────────────────────┐
                 │ new objects go here              │
                 │ [A][B][C][D][E][F][G][H]        │
                 │  ↑  ↑  ↑  ↑  ↑  ↑  ↑           │
                 │  survived → promoted             │
                 └──────────────┬──────────────────┘
                                ↓ survived several collections
                 GENERATION 1 (old space) — collected every few seconds
                 ┌─────────────────────────────────┐
                 │ [long-lived objects]             │
                 └─────────────────────────────────┘
```

**Pros:** Very efficient — most collections only process the small nursery, which is fast.

**Cons:** Complex to implement correctly. Go uses a variant of this.

---

## 6. Why Astra Chose Mark-and-Sweep

For Astra v1.0, we chose mark-and-sweep because:

1. **Correctness:** It correctly handles all cases, including cycles.
2. **Simplicity:** ~150 lines of C. We can implement it in one chapter.
3. **No 2x memory overhead:** Unlike copying GC.
4. **No cycle bugs:** Unlike reference counting.
5. **Educational:** The algorithm is clean and easy to understand.

The tradeoff: stop-the-world pauses. Our programs will pause briefly during collection. For v1.0 this is acceptable. A future volume will add a concurrent, generational GC.

---

## 7. The Object Header — How the GC Tracks Objects

Every heap-allocated object in Astra has a **header** prepended to its data. The header contains metadata that the GC needs:

```c
// runtime/gc.h

typedef struct AstraObject {
    uint32_t size;              // payload size in bytes
    uint32_t type_id;           // type identifier (ASTRA_TYPE_*)
    uint8_t  marked;            // GC mark bit: 0 = not marked, 1 = marked
    uint8_t  _padding[3];       // align to 4 bytes
    struct AstraObject* next;   // intrusive linked list of ALL objects
    // payload follows immediately after this struct
} AstraObject;

// The total size of an AstraObject with payload:
// sizeof(AstraObject) + payload_size
```

The **intrusive linked list** is crucial: every allocated object is linked into a global list. The GC can sweep all objects by walking this list — no separate bookkeeping needed.

```
GLOBAL OBJECT LIST

gc_head ──► [Header|"hello"] ──► [Header|User{..}] ──► [Header|[1,2,3]] ──► NULL
             marked=0             marked=0               marked=0

After marking (User{} and [1,2,3] are reachable, "hello" is not):

gc_head ──► [Header|"hello"] ──► [Header|User{..}] ──► [Header|[1,2,3]] ──► NULL
             marked=0             marked=1               marked=1

After sweeping ("hello" is freed):

gc_head ──► [Header|User{..}] ──► [Header|[1,2,3]] ──► NULL
             marked=0               marked=0
             (mark bits reset for next cycle)
```

To get from a user pointer to the header, we subtract `sizeof(AstraObject)`:

```c
// Given a pointer to the payload, get the header.
static AstraObject* obj_header(void* ptr) {
    return (AstraObject*)((char*)ptr - sizeof(AstraObject));
}

// Given a header, get the payload pointer.
static void* obj_payload(AstraObject* obj) {
    return (void*)((char*)obj + sizeof(AstraObject));
}
```

---

## 8. GC Roots — What the GC Must Never Collect

The GC starts tracing from **roots** — the set of all live pointers that exist at collection time. These are:

1. **Global variables:** Variables declared at the top level of Astra modules.
2. **Stack variables:** Local variables in all currently-active functions.
3. **CPU registers:** Some pointers may be in registers (the GC must scan these too, but we use a conservative approach for v1.0).

Finding stack roots is the hardest part. A production GC uses **stack maps** — metadata emitted by the compiler that says "at instruction address X, register R1 holds a pointer." This is complex.

**Astra v1.0's conservative approach:** We scan the entire stack as if every word might be a pointer. If a word looks like it falls within the heap's address range, we treat it as a pointer. This is called a **conservative GC** — it may accidentally keep some garbage alive (false positives), but it never collects live objects (no false negatives).

```c
// runtime/gc.c — conservative stack scanning (simplified)

extern char __data_start;  // start of data segment (linker symbol)
extern char _end;          // end of data segment (linker symbol)

// Scan the stack from current stack pointer to the top of the stack.
// Every word that looks like a heap pointer is treated as a root.
static void gc_scan_stack(void) {
    // Get current stack pointer
    void* sp;
    __asm__("mov %%rsp, %0" : "=r"(sp));

    // Walk up the stack, checking each word
    void** word = (void**)sp;
    while (word < (void**)stack_top) {
        void* potential_ptr = *word;
        if (gc_is_heap_ptr(potential_ptr)) {
            AstraObject* obj = obj_header(potential_ptr);
            gc_mark(obj);
        }
        word++;
    }
}

// Scan global variables (BSS + data segments).
static void gc_scan_globals(void) {
    void** word = (void**)&__data_start;
    void** end  = (void**)&_end;
    while (word < end) {
        void* potential_ptr = *word;
        if (gc_is_heap_ptr(potential_ptr)) {
            AstraObject* obj = obj_header(potential_ptr);
            gc_mark(obj);
        }
        word++;
    }
}
```

---

## 9. The Mark Phase — Tracing Live Objects

Once we have the roots, we trace all pointers recursively:

```
MARK PHASE ALGORITHM

1. Start with all objects unmarked (marked = 0).
2. For each root pointer:
   a. If the object it points to is already marked, skip.
   b. Mark the object (set marked = 1).
   c. For each pointer field inside the object, recursively mark.
3. When no more objects can be marked, phase is done.

Objects that are still unmarked after this phase are garbage.
```

The mark phase needs to know the pointer fields inside each object. This is where `type_id` comes in:

```c
// runtime/gc.c

// Mark a single object and all objects reachable from it.
void gc_mark(AstraObject* obj) {
    if (!obj || obj->marked) return;  // already marked or null
    obj->marked = 1;

    // Based on the object's type, we know where the pointer fields are.
    // This is a simplified dispatch — a real GC uses a type descriptor.
    void* payload = obj_payload(obj);

    switch (obj->type_id) {
        case ASTRA_TYPE_STRING: {
            // AstraString: data pointer (but data is raw bytes, not a GC obj)
            // No further tracing needed unless we store GC refs in strings.
            break;
        }
        case ASTRA_TYPE_LIST: {
            // AstraList: array of pointers to AstraObjects
            AstraList* list = (AstraList*)payload;
            for (int64_t i = 0; i < list->len; i++) {
                if (list->data[i]) {
                    gc_mark(obj_header(list->data[i]));
                }
            }
            break;
        }
        case ASTRA_TYPE_STRUCT: {
            // AstraStruct: pointer fields described by type descriptor
            AstraStructHeader* sh = (AstraStructHeader*)payload;
            for (uint32_t i = 0; i < sh->num_ptr_fields; i++) {
                void* field = sh->ptr_fields[i];
                if (field) gc_mark(obj_header(field));
            }
            break;
        }
        case ASTRA_TYPE_CLOSURE: {
            // Closure: captured variables (all treated as potential pointers)
            AstraClosure* cl = (AstraClosure*)payload;
            for (uint32_t i = 0; i < cl->num_captures; i++) {
                if (cl->captures[i]) {
                    gc_mark(obj_header(cl->captures[i]));
                }
            }
            break;
        }
        default:
            // RAW_BYTES and other non-pointer types: no further tracing
            break;
    }
}
```

---

## 10. The Sweep Phase — Reclaiming Dead Objects

After marking, the sweep phase walks the entire object list and frees anything that is not marked:

```c
// runtime/gc.c

void gc_sweep(void) {
    AstraObject** current = &gc_head;  // pointer to the 'next' field

    while (*current) {
        AstraObject* obj = *current;

        if (!obj->marked) {
            // This object was not marked: it is garbage. Free it.
            *current = obj->next;  // remove from list
            gc_bytes_freed += sizeof(AstraObject) + obj->size;
            free(obj);             // return to OS/malloc
        } else {
            // This object is live. Reset mark bit for next GC cycle.
            obj->marked = 0;
            current = &obj->next;  // advance pointer
        }
    }
}
```

Note the clever pointer-to-pointer trick (`AstraObject** current`). By maintaining a pointer to the `next` field of the previous object, we can remove nodes from the linked list without a separate `prev` pointer.

---

## 11. Write Barriers

There is a subtle problem: what if a program modifies a pointer *during* the mark phase? The GC might miss the new pointer. This is called the "write barrier" problem.

In Astra v1.0, we sidestep this by doing **stop-the-world** collection: the entire program pauses during GC. No pointers can be modified while the GC runs. This is safe but causes visible pauses.

A note for future improvement: concurrent GC (where collection runs alongside the program) requires write barriers — every pointer write must notify the GC. This is how modern Go works but adds ~5% overhead to pointer writes.

---

## 12. The Complete GC Implementation

```c
// runtime/gc.c
// Mark-and-sweep garbage collector for the Astra runtime.
// Build: cc -c -O2 gc.c -o gc.o

#include <stdlib.h>
#include <string.h>
#include <stdint.h>
#include <stdio.h>
#include "gc.h"
#include "runtime.h"

// ── Global state ──────────────────────────────────────────────

// Linked list head — all allocated objects
static AstraObject* gc_head = NULL;

// Heap statistics
static size_t gc_bytes_allocated = 0;
static size_t gc_bytes_freed     = 0;
static size_t gc_num_collections = 0;
static size_t gc_collection_threshold = 0;

// Stack top (set at startup)
static void* stack_top = NULL;

// Heap address bounds (for conservative root scanning)
static uintptr_t heap_lo = UINTPTR_MAX;
static uintptr_t heap_hi = 0;

// ── Initialization ────────────────────────────────────────────

void gc_init(size_t initial_heap_size) {
    // Record the approximate stack top (caller's frame)
    // Using __builtin_frame_address(0) is GCC/Clang-specific but portable enough
    stack_top = __builtin_frame_address(0);

    // Set initial collection threshold
    gc_collection_threshold = initial_heap_size;

    // heap_lo/hi will be updated as we allocate
    heap_lo = UINTPTR_MAX;
    heap_hi = 0;
}

void gc_shutdown(void) {
    // Free all remaining objects (run a full sweep)
    AstraObject* obj = gc_head;
    while (obj) {
        AstraObject* next = obj->next;
        free(obj);
        obj = next;
    }
    gc_head = NULL;
}

// ── Allocation ────────────────────────────────────────────────

// Allocate a new GC-managed object with the given payload size and type.
// Returns a pointer to the payload (NOT the header).
void* gc_alloc(size_t payload_size, uint32_t type_id) {
    // Trigger collection if we are above the threshold
    if (gc_bytes_allocated - gc_bytes_freed > gc_collection_threshold) {
        gc_collect();
        // Grow threshold after collection to avoid collecting too often
        gc_collection_threshold =
            (gc_bytes_allocated - gc_bytes_freed) * 2;
        if (gc_collection_threshold < 1024 * 1024) {
            gc_collection_threshold = 1024 * 1024;  // minimum 1MB
        }
    }

    // Allocate: header + payload
    size_t total_size = sizeof(AstraObject) + payload_size;
    AstraObject* obj  = (AstraObject*)malloc(total_size);

    if (!obj) {
        // Out of memory: run GC and try again
        gc_collect();
        obj = (AstraObject*)malloc(total_size);
        if (!obj) {
            // Still out of memory: panic
            // (use write directly to avoid re-entering the allocator)
            const char* oom = "[PANIC] out of memory\n";
            extern void astra_eprint(const char*);
            astra_eprint(oom);
            _exit(1);
        }
    }

    // Initialize the header
    memset(obj, 0, total_size);
    obj->size    = (uint32_t)payload_size;
    obj->type_id = type_id;
    obj->marked  = 0;

    // Prepend to the global object list
    obj->next = gc_head;
    gc_head   = obj;

    // Track heap bounds
    uintptr_t payload_addr = (uintptr_t)obj_payload(obj);
    if (payload_addr < heap_lo) heap_lo = payload_addr;
    if (payload_addr + payload_size > heap_hi) heap_hi = payload_addr + payload_size;

    // Update statistics
    gc_bytes_allocated += total_size;

    return obj_payload(obj);
}

// ── Is a word a heap pointer? ─────────────────────────────────

static int gc_is_heap_ptr(void* ptr) {
    uintptr_t addr = (uintptr_t)ptr;
    // Must be within known heap bounds and aligned to pointer size
    return addr >= heap_lo &&
           addr <  heap_hi &&
           (addr % sizeof(void*)) == 0;
}

// ── Pointer arithmetic helpers ────────────────────────────────

static AstraObject* obj_header(void* payload_ptr) {
    return (AstraObject*)((char*)payload_ptr - sizeof(AstraObject));
}

static void* obj_payload(AstraObject* obj) {
    return (void*)((char*)obj + sizeof(AstraObject));
}

// ── Mark phase ────────────────────────────────────────────────

void gc_mark(AstraObject* obj) {
    if (!obj || obj->marked) return;
    obj->marked = 1;

    void* payload = obj_payload(obj);

    switch (obj->type_id) {
        case ASTRA_TYPE_STRING: {
            // AstraString struct — its data field points to ASTRA_TYPE_RAW_BYTES
            AstraString* s = (AstraString*)payload;
            if (s->data) {
                AstraObject* data_obj = obj_header(s->data);
                if (gc_is_heap_ptr(s->data)) gc_mark(data_obj);
            }
            break;
        }
        case ASTRA_TYPE_LIST: {
            // Generic list — scan all element slots
            AstraList* list = (AstraList*)payload;
            for (int64_t i = 0; i < list->len; i++) {
                void* elem = list->data[i];
                if (elem && gc_is_heap_ptr(elem)) {
                    gc_mark(obj_header(elem));
                }
            }
            break;
        }
        case ASTRA_TYPE_STRUCT: {
            // Struct — uses embedded pointer map
            // For v1.0: conservatively scan the entire payload for pointers
            void** word = (void**)payload;
            void** end  = (void**)((char*)payload + obj->size);
            while (word < end) {
                if (*word && gc_is_heap_ptr(*word)) {
                    gc_mark(obj_header(*word));
                }
                word++;
            }
            break;
        }
        case ASTRA_TYPE_CLOSURE: {
            // Closure captures — scan conservatively
            void** word = (void**)payload;
            void** end  = (void**)((char*)payload + obj->size);
            while (word < end) {
                if (*word && gc_is_heap_ptr(*word)) {
                    gc_mark(obj_header(*word));
                }
                word++;
            }
            break;
        }
        case ASTRA_TYPE_RAW_BYTES:
        default:
            // No pointers inside raw bytes
            break;
    }
}

// ── Root scanning ─────────────────────────────────────────────

// Scan a memory region for potential heap pointers.
static void gc_scan_region(void* start, void* end) {
    if (start > end) {
        void* tmp = start; start = end; end = tmp;
    }
    void** word = (void**)((uintptr_t)start & ~(sizeof(void*) - 1));
    void** last = (void**)end;
    while (word < last) {
        void* potential_ptr = *word;
        if (potential_ptr && gc_is_heap_ptr(potential_ptr)) {
            gc_mark(obj_header(potential_ptr));
        }
        word++;
    }
}

void gc_mark_roots(void) {
    // 1. Scan the call stack (from current SP to the stack top we recorded at init)
    void* sp;
    __asm__("mov %%rsp, %0" : "=r"(sp));
    gc_scan_region(sp, stack_top);

    // 2. Scan global/static variables (BSS and DATA segments)
    // On Linux/macOS these linker symbols bracket the static data area.
    // For portability, we skip this and rely on stack scanning in v1.0.
    // A future version will use compiler-emitted root tables.
}

// ── Sweep phase ───────────────────────────────────────────────

void gc_sweep(void) {
    AstraObject** current = &gc_head;
    while (*current) {
        AstraObject* obj = *current;
        if (!obj->marked) {
            // Garbage: unlink and free
            *current = obj->next;
            size_t total = sizeof(AstraObject) + obj->size;
            gc_bytes_freed += total;
            free(obj);
        } else {
            // Live: clear mark, advance
            obj->marked = 0;
            current = &obj->next;
        }
    }
}

// ── Full collection cycle ─────────────────────────────────────

void gc_collect(void) {
    gc_num_collections++;
    gc_mark_roots();
    gc_sweep();
}

// ── Statistics ────────────────────────────────────────────────

void gc_print_stats(void) {
    printf("[GC] collections=%zu  allocated=%zu  freed=%zu  live=%zu bytes\n",
           gc_num_collections,
           gc_bytes_allocated,
           gc_bytes_freed,
           gc_bytes_allocated - gc_bytes_freed);
}

size_t gc_live_bytes(void) {
    return gc_bytes_allocated - gc_bytes_freed;
}
```

**Header file:**

```c
// runtime/gc.h

#ifndef ASTRA_GC_H
#define ASTRA_GC_H

#include <stdint.h>
#include <stddef.h>

// Object header — prepended before every GC-managed allocation
typedef struct AstraObject {
    uint32_t size;             // payload size in bytes
    uint32_t type_id;          // type identifier
    uint8_t  marked;           // GC mark bit
    uint8_t  _pad[3];          // alignment padding
    struct AstraObject* next;  // intrusive linked list
    // payload data follows immediately
} AstraObject;

// List type (pointed to by ASTRA_TYPE_LIST objects)
typedef struct {
    void**   data;    // array of pointers to AstraObjects
    int64_t  len;     // number of elements
    int64_t  cap;     // capacity
} AstraList;

// GC lifecycle
void  gc_init(size_t initial_heap_size);
void  gc_shutdown(void);

// Allocation
void* gc_alloc(size_t payload_size, uint32_t type_id);

// Collection
void  gc_mark_roots(void);
void  gc_mark(AstraObject* obj);
void  gc_sweep(void);
void  gc_collect(void);

// Statistics
void   gc_print_stats(void);
size_t gc_live_bytes(void);

#endif // ASTRA_GC_H
```

---

## 13. Integrating the GC with the Runtime

With the GC in place, all allocations in the runtime go through `gc_alloc` instead of `malloc`:

```c
// In runtime.c, replace:
//   void* buf = malloc(size);
// with:
//   void* buf = gc_alloc(size, ASTRA_TYPE_RAW_BYTES);

AstraString* astra_string_new(const char* data, size_t len) {
    // Allocate the string struct
    AstraString* s = (AstraString*)gc_alloc(sizeof(AstraString),
                                             ASTRA_TYPE_STRING);
    // Allocate the raw bytes
    s->data = (char*)gc_alloc(len + 1, ASTRA_TYPE_RAW_BYTES);
    memcpy(s->data, data, len);
    s->data[len] = '\0';
    s->len = len;
    s->cap = len;
    return s;
}
```

When you write in Astra:
```astra
fn process() {
    let s = "temporary string"
    // s goes out of scope at end of function
}
// After process() returns, 's' is no longer reachable.
// The next GC collection will free it automatically.
```

The Astra programmer never writes `free()`. The GC takes care of it.

---

## 14. Astra Build Milestone

The runtime directory is now:
```
runtime/
├── runtime.h        ← AstraString, type IDs, function declarations
├── runtime.c        ← startup, panic, I/O, strings, bounds check
├── gc.h             ← AstraObject, GC function declarations
└── gc.c             ← complete mark-and-sweep GC (~220 lines)
```

Let us verify the GC works:

```c
// test_gc.c

#include "gc.h"
#include "runtime.h"
#include <stdio.h>

void astra_user_main(void) {
    // Allocate many strings — most should be collected
    for (int i = 0; i < 10000; i++) {
        AstraString* s = astra_string_new("temporary", 9);
        (void)s;  // s immediately becomes garbage (no live reference)
    }

    printf("After loop, before GC: %zu bytes live\n", gc_live_bytes());
    gc_collect();
    printf("After GC:              %zu bytes live\n", gc_live_bytes());
    gc_print_stats();
}
```

Expected output:
```
After loop, before GC: ~560000 bytes live   (10000 * ~56 bytes each)
After GC:              ~0 bytes live         (all collected!)
[GC] collections=1  allocated=560000  freed=560000  live=0 bytes
```

The GC is working. Objects that are no longer reachable are collected automatically.

---

## 15. Exercises

**Exercise 64.1 — Count live objects**

Add a function `size_t gc_live_objects(void)` that returns the count of objects currently in the linked list. Print this alongside `gc_live_bytes()` in `gc_print_stats()`.

**Exercise 64.2 — GC triggered by threshold**

Currently the GC threshold doubles after each collection. Implement a smarter strategy: if the live bytes after collection are less than 25% of the old threshold, shrink the threshold by half. This prevents the threshold from growing without bound in long-running programs.

**Exercise 64.3 — Finalization**

Add support for "finalizers" — callbacks that run just before an object is freed. Add a function pointer `void (*finalizer)(void* payload)` to the `AstraObject` header. Call it in `gc_sweep()` before `free(obj)`. Write a test that uses a finalizer to print "freeing object" each time an object is collected.

**Exercise 64.4 — Heap fragmentation visualization**

After running many allocations and collections, the heap may be fragmented. Write a function `gc_visualize_heap()` that prints the sequence of all live objects as a grid, showing their sizes. For example: `[56][240][56][1024][56][56]`. This helps visualize fragmentation.

**Exercise 64.5 — Reference counting comparison**

Implement a simple reference-counting allocator (as a separate `rc.c` file, not replacing the GC). Each object has a `uint32_t refcount` field. Add `rc_retain(ptr)` (increment refcount) and `rc_release(ptr)` (decrement; free if 0). Write a test showing the cycle problem: two objects that point to each other and have no external references, but never get freed by reference counting.

**Exercise 64.6 — GC pause timing**

Add timing to `gc_collect()` using `clock_gettime(CLOCK_MONOTONIC)`. Print the GC pause duration (in milliseconds) to stderr when `ASTRA_DEBUG=1` is set. Run the test suite and measure how long GC pauses are.

**Exercise 64.7 — Incremental GC sketch**

The stop-the-world GC pauses the program for the entire collection. An "incremental" GC does a little bit of work each allocation, spreading the pause out. Sketch (in pseudocode, not full implementation) how you would implement incremental marking: instead of running the full mark phase at once, do a fixed number of marking steps per `gc_alloc` call.

---

## Chapter Summary

| Concept | Description | Implementation |
|---|---|---|
| Stack memory | Automatic, LIFO, ~8 MB limit | CPU stack pointer manipulation |
| Heap memory | Dynamic, any size, GC-managed | `gc_alloc()` |
| `malloc`/`free` | C's manual memory management | Replaced by `gc_alloc` in Astra |
| Memory leak | Forgetting to free | Eliminated by GC |
| Dangling pointer | Using freed memory | Eliminated by GC |
| Mark-and-sweep | Trace live objects, collect rest | `gc_mark_roots()` + `gc_sweep()` |
| Object header | `AstraObject` prepended to payload | 16 bytes of metadata |
| GC roots | Stack + globals — starting points | `gc_scan_region()` |
| Mark phase | Recursively mark reachable objects | `gc_mark()` |
| Sweep phase | Free all unmarked objects | `gc_sweep()` |
| Collection threshold | Trigger GC when live bytes exceed N | Doubles after each collection |
| Conservative GC | Treat any heap-range word as pointer | Safe but may keep some garbage alive |

In Chapter 65, we build the standard library — the first layer above the runtime that gives Astra programs useful, high-level capabilities like string manipulation, math, and time.
