# Chapter 11: Arrays and Dynamic Arrays — Storing Collections of Data

> "Arrays are the most fundamental data structure in computing — a contiguous block of memory holding values of the same type. Everything else builds on them."
> — Anonymous

---

## Chapter Overview

Almost every useful program needs to work with **collections** of data — a list of users, a sequence of scores, a batch of commands, a stream of pixels. The most fundamental collection in computing is the **array**: a fixed-size, contiguous block of memory that holds multiple values of the same type. Arrays are everywhere in programming, from the pixels on your screen (an array of color values) to the characters in a string (an array of bytes) to the DNS records that route internet traffic (an array of IP addresses).

But real programs often don't know in advance how many items they will need to store. How many users will sign up? How many lines will a file have? How many search results will come back? For these cases, we need **dynamic arrays** — arrays that can grow as needed. In Go, dynamic arrays are called **slices**. In Astra, they will be called `List<T>`. Both are built on top of a plain array with some clever bookkeeping.

In this chapter, we will understand arrays and slices deeply — not just how to use them, but how they work internally. We will look at memory layouts, cache behavior, amortized complexity, and the three-word slice header. We will also design the C runtime data structure that will power Astra's `List<T>`, giving you a preview of how language runtimes are implemented.

---

## What We're Building

Our Astra compiler and runtime work together: the compiler translates Astra source code into instructions, and the runtime provides the data structures those instructions operate on. In this chapter, we design the `AstraList` C struct that the runtime will use to implement `List<T>`. This is a real piece of the Astra runtime — the first runtime code we will write in this guide.

---

## Table of Contents

1. What is an Array? — A Contiguous Block of Memory
2. Fixed Arrays in Go: `[N]T`
3. Memory Layout — Why Arrays Are Fast
4. Array Indexing: 0-Based, Bounds Checking, and Panics
5. Why Fixed Arrays Are Limiting
6. Dynamic Arrays (Slices) in Go — The Three-Word Header
7. Creating and Using Slices
8. `append()` — How Growth Works and Amortized O(1)
9. Slice Operations: Slicing, `copy()`, Clearing
10. 2D Arrays and Slices
11. Arrays in Astra: Fixed `[N]T` and Dynamic `List<T>`
12. Cache Locality — Why Arrays Beat Linked Lists in Practice
13. Choosing the Right Data Structure
14. Astra Build Milestone: The `AstraList` Runtime Structure
15. Exercises
16. Summary

---

## 1. What is an Array? — A Contiguous Block of Memory

An **array** is a collection of values of the same type, stored **contiguously** in memory — one right after another, with no gaps. The word "contiguous" is key: it means the elements are stored in adjacent memory locations, not scattered around.

```
Array: [10, 20, 30, 40, 50]  (array of 5 integers)

Memory (each int = 8 bytes on 64-bit):
Address: 1000  1008  1016  1024  1032
         ┌─────┬─────┬─────┬─────┬─────┐
Value:   │ 10  │ 20  │ 30  │ 40  │ 50  │
         └─────┴─────┴─────┴─────┴─────┘
Index:     [0]   [1]   [2]   [3]   [4]
```

**To find element at index `i`:**
```
address = base_address + (i * element_size)
element[2] = 1000 + (2 * 8) = 1016  → holds value 30
```

This is called **random access** and it is O(1) — constant time. It doesn't matter if the array has 5 elements or 5 billion elements; finding element `i` takes the same amount of time because it's just a single arithmetic operation.

**Real-world analogy:** An array is like a row of mailboxes in an apartment building. If mailbox 1 is at the entrance, mailbox 2 is 30 cm to the right, mailbox 3 is another 30 cm, and so on — you can go directly to any mailbox without looking through all the others. Just count from the entrance.

---

## 2. Fixed Arrays in Go: `[N]T`

In Go, a fixed-size array is declared with the size as part of the type:

```go
// Declaration and initialization
var a [5]int              // [0 0 0 0 0] — zero-initialized
b := [3]string{"a", "b", "c"}   // ["a" "b" "c"]
c := [...]int{1, 2, 3, 4}       // [...] lets compiler count: [4]int

// Accessing elements
fmt.Println(b[0])   // "a"
fmt.Println(b[2])   // "c"

// Modifying elements
a[0] = 10
a[1] = 20
fmt.Println(a)  // [10 20 0 0 0]

// Length
fmt.Println(len(a))  // 5

// Iterating
for i, v := range b {
    fmt.Printf("b[%d] = %s\n", i, v)
}
```

**Important: in Go, the size is part of the type.**

```go
var x [5]int
var y [10]int
// x and y have DIFFERENT types: [5]int vs [10]int
// You cannot assign x = y — type mismatch!
```

This is unusual compared to most languages. The practical consequence is that `[5]int` and `[10]int` are as different to Go as `int` and `string`. Functions that accept `[5]int` cannot accept `[10]int`.

**Arrays are values in Go** — when you assign one array to another, or pass one to a function, the entire array is **copied**:

```go
a := [3]int{1, 2, 3}
b := a          // b is a COPY of a
b[0] = 999
fmt.Println(a)  // [1 2 3] — a is unchanged
fmt.Println(b)  // [999 2 3] — only b was modified
```

This copy semantics is why arrays are rarely used directly in Go for large data — you'd be copying megabytes of data on every function call.

---

## 3. Memory Layout — Why Arrays Are Fast

The contiguous memory layout of arrays has a profound effect on performance, thanks to the CPU's **cache system**.

Modern CPUs are thousands of times faster than RAM. To compensate, CPUs have multiple levels of **cache** — small, fast memory banks close to the processor:

```
Access speed (approximate):
┌───────────────────────────────────────────────────────┐
│  CPU Registers  │  0.3 ns    (basically instant)      │
│  L1 Cache       │  1 ns      (on chip)                │
│  L2 Cache       │  4 ns      (on chip)                │
│  L3 Cache       │  10 ns     (on chip)                │
│  RAM            │  100 ns    (100x slower than L1)    │
│  SSD            │  100,000 ns (100x slower than RAM)  │
│  HDD            │  10,000,000 ns (1000x slower)       │
└───────────────────────────────────────────────────────┘
```

When the CPU reads one byte from RAM, it actually fetches a whole **cache line** (typically 64 bytes). If the next byte you need is also in that cache line (as it would be in an array), it is already in cache — a **cache hit**, which is essentially free.

```
Array access pattern (cache-friendly):
Read arr[0] → cache loads: arr[0], arr[1], ..., arr[7] (64 bytes)
Read arr[1] → CACHE HIT! Already loaded.
Read arr[2] → CACHE HIT! Already loaded.
... all 8 elements in one cache line = very fast

Linked list access pattern (cache-unfriendly):
Read node[0] → cache loads 64 bytes around node[0]'s address
               node[1] is at a RANDOM address elsewhere in RAM
Read node[1] → CACHE MISS! Must go to RAM again. (100ns penalty)
Read node[2] → CACHE MISS! Another 100ns penalty.
```

This is why arrays are often faster than linked lists even when the algorithmic complexity (Big-O) is the same. The hardware behavior matters enormously.

---

## 4. Array Indexing: 0-Based, Bounds Checking, and Panics

Arrays in Go (and Astra) are **0-indexed**: the first element is at index 0, the last at index `len-1`.

```go
arr := [5]int{10, 20, 30, 40, 50}
//      [0]   [1]  [2]  [3]  [4]
fmt.Println(arr[0])  // 10 — first element
fmt.Println(arr[4])  // 50 — last element
fmt.Println(arr[5])  // PANIC: index out of range [5] with length 5
```

**Why 0-indexed?** It comes from C, where array indexing is defined as pointer arithmetic: `arr[i]` means "start at `arr`, add `i * element_size` bytes." With that definition, the first element is naturally at offset 0. Python, Java, Go, Rust, and most modern languages follow this convention.

**Bounds checking** is performed automatically in Go and Astra. If you try to access `arr[n]` where `n >= len(arr)` or `n < 0`, the runtime immediately panics (throws a fatal error) rather than reading garbage memory (which is what C does — causing security vulnerabilities).

```go
// In Go:
arr := [5]int{1, 2, 3, 4, 5}
i := 10
arr[i]   // panic: runtime error: index out of range [10] with length 5
         // The program stops immediately with a helpful error message

// In C:
int arr[] = {1, 2, 3, 4, 5};
int i = 10;
arr[i];   // undefined behavior — reads from memory you don't own
          // could crash, could return garbage, could be a security hole
```

Bounds checking is one of the most important safety features of memory-safe languages. The Astra compiler will emit bounds-check code whenever a list is indexed.

---

## 5. Why Fixed Arrays Are Limiting

Fixed arrays have a serious problem: **you must know the size at compile time.**

```go
// This works:
var arr [5]int   // size known at compile time

// This does NOT work in Go:
n := getUserInput()
var arr [n]int   // ERROR: n is not a constant — size must be known at compile time
```

For most real programs, you don't know in advance how much data you will have:
- How many users are in the database? Unknown.
- How many lines in a file? Unknown.
- How many results from a search? Unknown.

This is why **dynamic arrays** exist. A dynamic array is an array that can grow (and sometimes shrink) at runtime as you add or remove elements.

---

## 6. Dynamic Arrays (Slices) in Go — The Three-Word Header

Go's dynamic arrays are called **slices**. A slice is not just an array — it is a small **header** structure that describes a view into an underlying array.

```
A slice is internally a three-word struct:
┌─────────────────────────────────────────────────────────┐
│  ptr │ 0xC0001A0000   pointer to the first element      │
│  len │ 5              number of elements in the slice   │
│  cap │ 8              number of elements the array can  │
│                       hold before needing to grow       │
└─────────────────────────────────────────────────────────┘
          │
          ↓
Backing array (on heap):
┌────┬────┬────┬────┬────┬────┬────┬────┐
│ 10 │ 20 │ 30 │ 40 │ 50 │    │    │    │
└────┴────┴────┴────┴────┴────┴────┴────┘
  [0]  [1]  [2]  [3]  [4]   cap  cap  cap
                              ↑
                     capacity extends here
                     (unused, reserved space)
```

- **`ptr`**: a pointer to the start of the backing array in heap memory
- **`len`** (length): how many elements are currently in the slice
- **`cap`** (capacity): how many elements the backing array can hold

When you append elements and `len < cap`, Go just extends `len` (fast). When `len == cap` and you append, Go must allocate a *larger* array and copy everything over (slow, but rare).

---

## 7. Creating and Using Slices

```go
// Method 1: slice literal
s := []int{1, 2, 3, 4, 5}
fmt.Println(s)      // [1 2 3 4 5]
fmt.Println(len(s)) // 5
fmt.Println(cap(s)) // 5 (or more, depends on implementation)

// Method 2: make — creates a slice with specified length and capacity
s2 := make([]int, 5)       // length=5, capacity=5, all zeros
s3 := make([]int, 3, 10)   // length=3, capacity=10 (pre-allocate 10 slots)

// Method 3: from an array
arr := [5]int{10, 20, 30, 40, 50}
s4 := arr[:]    // slice of the whole array
s5 := arr[1:3]  // slice of arr[1] and arr[2]: [20, 30]

// Accessing and modifying
s[0] = 100
fmt.Println(s)  // [100 2 3 4 5]

// Length and capacity
fmt.Printf("len=%d cap=%d\n", len(s3), cap(s3))   // len=3 cap=10

// Nil slice (zero value of a slice type)
var nilSlice []int
fmt.Println(nilSlice == nil)  // true
fmt.Println(len(nilSlice))    // 0 — safe to call len/cap on nil slice
```

---

## 8. `append()` — How Growth Works and Amortized O(1)

The `append` function adds elements to a slice:

```go
s := []int{1, 2, 3}   // len=3, cap=3

// Appending when there is room (cap > len):
// (hypothetically, if cap were larger)

// Appending when cap is full:
s = append(s, 4)   // len was 3, cap was 3 → must grow!
// Go allocates a new array (often 2× the old capacity)
// copies all elements over
// sets s.ptr to the new array
// returns the new slice header
fmt.Println(s)     // [1 2 3 4]
```

**Growth strategy (approximately):**

```
Initial:       cap=3
After append:  cap=6   (doubled)
After more appends (when full): cap=12 (doubled again)
And so on...
```

```
Amortized O(1) explained:

Imagine appending n elements one by one, starting from an empty slice.

Most appends: just increment len, set the element. O(1).
Occasional appends: must copy all elements to a new array. O(n).

But the total work for n appends is:
  n (fast appends) + (1 + 2 + 4 + 8 + ... + n) (copy costs)
= n + 2n (geometric series)
= 3n total work for n appends

Average work per append = 3n / n = 3 = O(1)

This is called "amortized O(1)" — the expensive copies are rare enough
that the average cost is still constant.
```

```
Visual growth trace:
append(s, 1):     [1]           len=1  cap=1
append(s, 2):     [1,2]         len=2  cap=2  ← grew! copied 1 element
append(s, 3):     [1,2,3]       len=3  cap=4  ← grew! copied 2 elements
append(s, 4):     [1,2,3,4]     len=4  cap=4
append(s, 5):     [1,2,3,4,5]   len=5  cap=8  ← grew! copied 4 elements
append(s, 6):     [1,2,3,4,5,6] len=6  cap=8
...
```

**Best practice:** if you know approximately how many elements you will have, pre-allocate:

```go
// Instead of:
var s []int
for i := 0; i < 1000000; i++ {
    s = append(s, i)   // grows ~20 times, lots of copying
}

// Do this:
s := make([]int, 0, 1000000)   // pre-allocate capacity for 1M elements
for i := 0; i < 1000000; i++ {
    s = append(s, i)   // never needs to grow — always fast
}
```

---

## 9. Slice Operations: Slicing, `copy()`, Clearing

**Sub-slicing (creating a slice from a slice):**

```go
s := []int{10, 20, 30, 40, 50}

// s[low:high] — elements from index low to high-1
a := s[1:4]    // [20, 30, 40]  len=3, cap=4
fmt.Println(a)

// s[:high] — from start to high-1
b := s[:3]     // [10, 20, 30]
fmt.Println(b)

// s[low:] — from low to end
c := s[2:]     // [30, 40, 50]
fmt.Println(c)

// WARNING: a, b, c all share the SAME backing array as s!
a[0] = 999
fmt.Println(s[1])  // 999 — s was modified through a!
```

**The `copy()` function — make an independent copy:**

```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, len(src))
n := copy(dst, src)     // copies min(len(dst), len(src)) elements
fmt.Println(n)          // 5 — number of elements copied
fmt.Println(dst)        // [1 2 3 4 5]

dst[0] = 999
fmt.Println(src[0])     // 1 — src NOT modified (dst is independent)
```

**Clearing a slice:**

```go
s := []int{1, 2, 3, 4, 5}
s = s[:0]            // len=0, cap=5 (backing array still there)
// OR
s = nil              // releases the backing array (eligible for GC)
// OR (Go 1.21+)
clear(s)             // zeroes all elements, keeps length
```

**Deleting an element at index `i` (order-preserving, O(n)):**

```go
s := []int{10, 20, 30, 40, 50}
i := 2  // delete index 2 (value 30)
s = append(s[:i], s[i+1:]...)
fmt.Println(s)  // [10 20 40 50]
```

---

## 10. 2D Arrays and Slices

A 2D array is an array of arrays — useful for matrices, grids, game boards, and images.

```go
// Fixed 2D array: 3 rows, 4 columns
var grid [3][4]int
grid[0][0] = 1
grid[2][3] = 99

// 2D slice (more flexible)
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

fmt.Println(matrix[1][2])  // 6

// Creating a 2D slice dynamically
rows, cols := 4, 5
m := make([][]int, rows)
for i := range m {
    m[i] = make([]int, cols)
}
m[2][3] = 42
```

**Memory layout of a 2D slice:**

```
matrix is a slice of slices:
matrix → [ ptr_row0, ptr_row1, ptr_row2 ]
              │           │          │
              ↓           ↓          ↓
         [1, 2, 3]   [4, 5, 6]  [7, 8, 9]
         (heap)      (heap)     (heap)
```

Note that the rows are NOT contiguous — each row is a separate heap allocation. For performance-critical matrix operations (like image processing), you want a truly contiguous 2D array stored as a 1D slice:

```go
// Flat representation — better cache performance
rows, cols := 4, 5
flat := make([]int, rows*cols)

// Access element at (row, col):
flat[row*cols + col] = value
```

---

## 11. Arrays in Astra: Fixed `[N]T` and Dynamic `List<T>`

Astra provides both fixed and dynamic array types:

```astra
// Fixed array: [N]T — size is part of the type
let rgb: [3]int = [255, 128, 0]
let matrix: [3][3]float = [
    [1.0, 0.0, 0.0],
    [0.0, 1.0, 0.0],
    [0.0, 0.0, 1.0],
]

// Dynamic array: List<T> — can grow at runtime
let numbers: List<int> = [1, 2, 3, 4, 5]

// Adding elements
numbers.push(6)         // adds to the end: [1,2,3,4,5,6]
numbers.push_front(0)   // adds to the front: [0,1,2,3,4,5,6]

// Removing elements
let last = numbers.pop()     // removes and returns 6
let first = numbers.pop_front() // removes and returns 0

// Accessing
print(numbers[0])        // 1
print(numbers.length())  // 5
print(numbers.is_empty()) // false

// Slicing
let slice = numbers[1..4]   // [2, 3, 4] (a new list)

// Iterating
for n in numbers {
    print(n)
}

// Map, filter, reduce (functional style)
let doubled = numbers.map(fn(n: int) -> int { return n * 2 })
let evens   = numbers.filter(fn(n: int) -> bool { return n % 2 == 0 })
let sum     = numbers.reduce(0, fn(acc: int, n: int) -> int { return acc + n })
```

**Type aliases for common cases:**

```astra
// Astra provides type aliases for clarity
type IntList    = List<int>
type StringList = List<string>
type Matrix     = List<List<float>>
```

---

## 12. Cache Locality — Why Arrays Beat Linked Lists in Practice

We touched on cache locality earlier. Here is a concrete demonstration:

**Benchmark scenario:** sum all elements in a collection of 1 million integers.

```
Array (slice) implementation:
  Data: [1, 2, 3, 4, ..., 1000000]  — all in one contiguous block
  Access pattern: sequential (arr[0], arr[1], arr[2], ...)
  Cache behavior: entire cache line loaded per miss (64 bytes = 8 ints)
  Cache misses: ~125,000 (1M / 8)
  Time: ~1 millisecond

Linked list implementation:
  Data: Node{1, next→} → Node{2, next→} → ... → Node{1000000, nil}
  Each node is at a random address in the heap
  Access pattern: pointer-chasing (follow ptr to find next element)
  Cache behavior: every access is likely a cache miss
  Cache misses: ~1,000,000
  Time: ~100 milliseconds (100x slower!)
```

This is not a theoretical difference — it is a real performance gap that shows up in production systems. Unless you need efficient insertions in the middle (linked list O(1) vs array O(n)), use arrays for collections.

**When linked lists win:**
- Frequent insertion/deletion at the beginning or middle of large lists
- Very large elements where copying is expensive (though you'd typically use pointers in an array)
- Implementing queues with frequent enqueue/dequeue from both ends

For most use cases in everyday programming, `List<T>` (backed by a dynamic array) is the right choice.

---

## 13. Choosing the Right Data Structure

```
I need to...                         Use...
──────────────────────────────────────────────────────────────────────
Store N items, N known at compile time → [N]T (fixed array)
Grow a collection dynamically          → []T (Go slice) / List<T> (Astra)
Look up items by a unique key          → map[K]V (hash map) — Chapter 14
Access items in insertion order        → []T (slice/list)
Always access the most recently added  → stack (slice with push/pop from end)
Always access the oldest item first    → queue (channel in Go, or slice)
Find if an item exists, fast           → map (hash set)
Need sorted order                      → sort a slice; or use a tree (Chapter 30)
Fixed-size 2D grid                     → [R][C]T or flat []T with manual indexing
```

---

## 14. Astra Build Milestone: The `AstraList` Runtime Structure

The Astra runtime is written in C and provides the built-in data structures that Astra programs use. Here is the complete implementation of `AstraList` — the runtime backing of Astra's `List<T>`:

```c
// runtime/list.h
// The AstraList is the runtime implementation of Astra's List<T> type.
// It is a generic dynamic array, similar to Go's slice.
// Written in C for maximum portability and control.

#ifndef ASTRA_LIST_H
#define ASTRA_LIST_H

#include <stddef.h>   // size_t
#include <stdlib.h>   // malloc, realloc, free
#include <string.h>   // memcpy, memmove
#include <assert.h>   // assert

// AstraList: the header structure (analogous to Go's slice header)
// All fields are part of the public interface.
typedef struct {
    void*  data;       // pointer to the backing array on the heap
    size_t len;        // current number of elements (len <= cap)
    size_t cap;        // current allocated capacity
    size_t elem_size;  // size of each element in bytes (e.g., 8 for int64)
} AstraList;

// INITIAL_CAP: starting capacity when a list is first created.
// Chosen small enough to avoid wasting memory for empty/tiny lists.
#define ASTRA_LIST_INITIAL_CAP 8

// GROWTH_FACTOR: when the list must grow, multiply capacity by this factor.
// 2.0 gives amortized O(1) appends (see Chapter 11).
#define ASTRA_LIST_GROWTH_FACTOR 2

// ─────────────────────────────────────────────────────────────
// list_new: create a new, empty AstraList
// elem_size: the size in bytes of each element (e.g., sizeof(int64_t) = 8)
// ─────────────────────────────────────────────────────────────
AstraList* list_new(size_t elem_size) {
    AstraList* list = (AstraList*)malloc(sizeof(AstraList));
    assert(list != NULL && "failed to allocate AstraList header");

    list->elem_size = elem_size;
    list->len       = 0;
    list->cap       = ASTRA_LIST_INITIAL_CAP;

    // Allocate the backing array
    list->data = malloc(elem_size * ASTRA_LIST_INITIAL_CAP);
    assert(list->data != NULL && "failed to allocate AstraList backing array");

    return list;
}

// ─────────────────────────────────────────────────────────────
// list_free: free an AstraList and its backing array
// ─────────────────────────────────────────────────────────────
void list_free(AstraList* list) {
    if (list == NULL) return;
    free(list->data);   // free the backing array first
    free(list);         // then free the header
}

// ─────────────────────────────────────────────────────────────
// list_grow: internal helper — grow the capacity to at least min_cap.
// Implements the doubling growth strategy.
// ─────────────────────────────────────────────────────────────
static void list_grow(AstraList* list, size_t min_cap) {
    size_t new_cap = list->cap;

    // Keep doubling until we have enough capacity
    while (new_cap < min_cap) {
        new_cap *= ASTRA_LIST_GROWTH_FACTOR;
    }

    // Reallocate the backing array with the new capacity
    // realloc copies existing elements and frees the old allocation
    void* new_data = realloc(list->data, new_cap * list->elem_size);
    assert(new_data != NULL && "failed to grow AstraList");

    list->data = new_data;
    list->cap  = new_cap;
}

// ─────────────────────────────────────────────────────────────
// list_push: append an element to the end of the list.
// elem: pointer to the element to copy into the list.
//       The list makes a copy — the caller retains ownership of elem.
// ─────────────────────────────────────────────────────────────
void list_push(AstraList* list, const void* elem) {
    // If full, grow first
    if (list->len == list->cap) {
        list_grow(list, list->cap + 1);
    }

    // Calculate the destination address: base + len * elem_size
    void* dest = (char*)list->data + (list->len * list->elem_size);

    // Copy the element into the backing array
    memcpy(dest, elem, list->elem_size);

    // Increment length
    list->len++;
}

// ─────────────────────────────────────────────────────────────
// list_get: get a pointer to the element at index.
// Returns NULL if index is out of bounds (Astra runtime will panic on NULL).
// The caller must not free the returned pointer.
// ─────────────────────────────────────────────────────────────
void* list_get(AstraList* list, size_t index) {
    if (index >= list->len) {
        return NULL;  // Astra runtime will convert this to a panic
    }
    return (char*)list->data + (index * list->elem_size);
}

// ─────────────────────────────────────────────────────────────
// list_set: set the element at index to a new value.
// Returns 0 on success, -1 if index is out of bounds.
// ─────────────────────────────────────────────────────────────
int list_set(AstraList* list, size_t index, const void* elem) {
    if (index >= list->len) {
        return -1;   // out of bounds
    }
    void* dest = (char*)list->data + (index * list->elem_size);
    memcpy(dest, elem, list->elem_size);
    return 0;
}

// ─────────────────────────────────────────────────────────────
// list_pop: remove and return a copy of the last element.
// Writes the popped value to 'out_elem' (caller must provide buffer).
// Returns 0 on success, -1 if the list is empty.
// ─────────────────────────────────────────────────────────────
int list_pop(AstraList* list, void* out_elem) {
    if (list->len == 0) {
        return -1;  // empty list
    }
    list->len--;
    void* src = (char*)list->data + (list->len * list->elem_size);
    if (out_elem != NULL) {
        memcpy(out_elem, src, list->elem_size);
    }
    return 0;
}

// ─────────────────────────────────────────────────────────────
// list_insert: insert an element at 'index', shifting subsequent elements right.
// O(n) time because we must shift elements.
// ─────────────────────────────────────────────────────────────
void list_insert(AstraList* list, size_t index, const void* elem) {
    assert(index <= list->len && "insert index out of bounds");

    if (list->len == list->cap) {
        list_grow(list, list->cap + 1);
    }

    // Shift elements from index onward to the right by one position
    // memmove handles overlapping regions correctly (unlike memcpy)
    void* src  = (char*)list->data + (index * list->elem_size);
    void* dest = (char*)list->data + ((index + 1) * list->elem_size);
    size_t bytes_to_move = (list->len - index) * list->elem_size;
    memmove(dest, src, bytes_to_move);

    // Write the new element at index
    memcpy(src, elem, list->elem_size);
    list->len++;
}

// ─────────────────────────────────────────────────────────────
// list_delete: remove the element at 'index', shifting subsequent elements left.
// O(n) time.
// ─────────────────────────────────────────────────────────────
void list_delete(AstraList* list, size_t index) {
    assert(index < list->len && "delete index out of bounds");

    void* dest = (char*)list->data + (index * list->elem_size);
    void* src  = (char*)list->data + ((index + 1) * list->elem_size);
    size_t bytes_to_move = (list->len - index - 1) * list->elem_size;
    memmove(dest, src, bytes_to_move);

    list->len--;
}

// ─────────────────────────────────────────────────────────────
// list_len: return the current number of elements.
// ─────────────────────────────────────────────────────────────
size_t list_len(const AstraList* list) {
    return list->len;
}

// ─────────────────────────────────────────────────────────────
// list_cap: return the current capacity.
// ─────────────────────────────────────────────────────────────
size_t list_cap(const AstraList* list) {
    return list->cap;
}

// ─────────────────────────────────────────────────────────────
// list_is_empty: returns 1 if empty, 0 if not.
// ─────────────────────────────────────────────────────────────
int list_is_empty(const AstraList* list) {
    return list->len == 0;
}

// ─────────────────────────────────────────────────────────────
// list_clear: reset length to 0 (does not free memory).
// Useful for reusing the list's capacity.
// ─────────────────────────────────────────────────────────────
void list_clear(AstraList* list) {
    list->len = 0;
}

#endif  // ASTRA_LIST_H
```

**Test file to verify the implementation:**

```c
// runtime/list_test.c
// Compile with: gcc -o list_test runtime/list_test.c && ./list_test
// Expected output: all assertions pass (no crashes)

#include <stdio.h>
#include <stdint.h>
#include <assert.h>
#include "list.h"

int main() {
    printf("Testing AstraList...\n");

    // Create a list of int64_t (Astra's int type)
    AstraList* list = list_new(sizeof(int64_t));
    assert(list_is_empty(list));
    assert(list_len(list) == 0);

    // Push elements
    int64_t values[] = {10, 20, 30, 40, 50};
    for (int i = 0; i < 5; i++) {
        list_push(list, &values[i]);
    }
    assert(list_len(list) == 5);

    // Get elements
    int64_t* elem = (int64_t*)list_get(list, 0);
    assert(*elem == 10);
    elem = (int64_t*)list_get(list, 4);
    assert(*elem == 50);

    // Out-of-bounds get returns NULL
    assert(list_get(list, 5) == NULL);
    assert(list_get(list, 100) == NULL);

    // Pop
    int64_t popped;
    int ok = list_pop(list, &popped);
    assert(ok == 0);
    assert(popped == 50);
    assert(list_len(list) == 4);

    // Set
    int64_t new_val = 999;
    list_set(list, 0, &new_val);
    assert(*(int64_t*)list_get(list, 0) == 999);

    // Insert at index 1
    int64_t inserted = 555;
    list_insert(list, 1, &inserted);
    assert(list_len(list) == 5);
    assert(*(int64_t*)list_get(list, 1) == 555);
    assert(*(int64_t*)list_get(list, 2) == 20);  // shifted right

    // Delete at index 1
    list_delete(list, 1);
    assert(list_len(list) == 4);
    assert(*(int64_t*)list_get(list, 1) == 20);  // shifted back

    // Test growth: push 100 elements (forces several reallocations)
    list_clear(list);
    for (int64_t i = 0; i < 100; i++) {
        list_push(list, &i);
    }
    assert(list_len(list) == 100);
    assert(*(int64_t*)list_get(list, 99) == 99);

    // Clean up
    list_free(list);

    printf("All tests passed!\n");
    return 0;
}
```

Build and run:

```bash
gcc -Wall -Wextra -o list_test runtime/list_test.c && ./list_test
# Output: All tests passed!
```

---

## 15. Exercises

1. **Memory Math** — An array of 1000 `int64` values (8 bytes each) is stored starting at memory address 4096. What is the address of:
   - Element at index 0?
   - Element at index 10?
   - Element at index 999?
   Show your calculations.

2. **Slice Header Inspection** — In Go, write a program that creates a slice, appends to it several times, and after each append prints `len`, `cap`, and the pointer address. Observe when the capacity doubles and when the pointer changes (indicating a reallocation).
   *Hint: use `reflect.SliceHeader` or just print `len(s)` and `cap(s)`*

3. **Off-by-One in Slicing** — What does each of the following Go expressions produce? Which ones panic?
   ```go
   s := []int{10, 20, 30, 40, 50}
   a := s[0:3]
   b := s[1:4]
   c := s[:5]
   d := s[5:]
   e := s[5:5]
   f := s[6:6]  // does this panic?
   ```

4. **2D Grid** — In Astra, write a function that creates a multiplication table as a `List<List<int>>` for values 1 through 5. The element at `table[i][j]` should be `(i+1) * (j+1)`.

5. **Implement `list_contains`** — Add a `list_contains` function to `list.h` that returns 1 if a given element is in the list. The function should accept a comparator function pointer:
   ```c
   int list_contains(AstraList* list, const void* elem,
                     int (*equals)(const void*, const void*));
   ```
   Write the implementation and a test for it.

6. **Stack with AstraList** — Using `AstraList`, implement a stack data structure with `push`, `pop`, and `peek` operations in C. A stack is LIFO (Last In, First Out). *Hint: `list_push` and `list_pop` already work from the end.*

7. **Cache Experiment** — Write two Go programs that sum 10 million integers. One uses a `[]int64` slice (contiguous). The other simulates a linked list by using random index access (`s[rand.Intn(len(s))]` in a loop). Time both programs. What do you observe?
   *Hint: use `time.Now()` and `time.Since()` for timing*

8. **Design `AstraList` for Strings** — Strings in Astra are `(ptr, len)` pairs (like Go strings). Design how you would store a `List<string>` in the `AstraList` struct. What would `elem_size` be? What does each element look like in memory?

---

## 16. Summary

| Concept | Go | Astra | Notes |
|---|---|---|---|
| Fixed array type | `[N]T` | `[N]T` | Size is part of the type |
| Dynamic array type | `[]T` (slice) | `List<T>` | Pointer + len + cap |
| Create fixed | `[5]int{1,2,3,4,5}` | `[5]int` | Zero-initialized by default |
| Create dynamic | `make([]T, len, cap)` or literal | `List<int>` literal or new | Pre-allocate with capacity |
| Append | `append(s, x)` | `list.push(x)` | Amortized O(1) |
| Get element | `s[i]` | `list[i]` | O(1), bounds-checked |
| Length | `len(s)` | `list.length()` | Returns current element count |
| Capacity | `cap(s)` | `list.capacity()` | Returns backing array size |
| Sub-slice | `s[low:high]` | `list[low..high]` | Shares backing array in Go |
| Independent copy | `copy(dst, src)` | `list.clone()` | Makes a new independent list |
| 0-indexed | Yes | Yes | First element at index 0 |
| Bounds checking | Runtime panic | Runtime panic | Prevents buffer overflows |
| Growth strategy | ~2x doubling | 2x doubling | Amortized O(1) append |
| Cache friendly | Yes | Yes | Contiguous memory layout |

**Key takeaways:**
- Arrays store same-type values contiguously in memory, enabling O(1) random access
- Fixed arrays have their size in the type; dynamic arrays (slices/List) grow at runtime
- Go slices have a three-word header: pointer, length, capacity
- `append()` has amortized O(1) cost — expensive reallocation is rare
- Bounds checking prevents buffer overflow security vulnerabilities
- Contiguous memory layout makes arrays cache-friendly and much faster than linked lists in practice
- The `AstraList` C struct is the runtime foundation of Astra's `List<T>` type

---

*Next chapter: Chapter 12 — Strings and Character Encoding*
