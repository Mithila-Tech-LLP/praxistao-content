# Chapter 08: Arrays and Slices

When you need to work with multiple values of the same type — a list of user IDs, a buffer of bytes, a sequence of prices — you need a collection type. Go has two: **arrays** (fixed-size, rarely used directly) and **slices** (dynamic, the workhorse of Go). Understanding slices deeply is critical because they are used everywhere in Go — from reading files to handling HTTP requests to building data pipelines.

## Table of Contents

1. [Arrays — Fixed-Size Collections](#1-arrays--fixed-size-collections)
2. [Slices — Dynamic Collections](#2-slices--dynamic-collections)
3. [Slice Internals — The Three-Field Header](#3-slice-internals--the-three-field-header)
4. [Creating and Growing Slices](#4-creating-and-growing-slices)
5. [Slicing — Cutting Sub-Slices](#5-slicing--cutting-sub-slices)
6. [Common Slice Operations](#6-common-slice-operations)
7. [2D Slices](#7-2d-slices)
8. [Common Pitfalls](#8-common-pitfalls)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Arrays — Fixed-Size Collections

An array has a **fixed size** that is part of its type. `[5]int` and `[10]int` are different types.

```go
// Declaration
var a [5]int              // Zero-initialized: [0 0 0 0 0]
b := [5]int{1, 2, 3, 4, 5}
c := [3]string{"Go", "is", "great"}

// Let the compiler count the elements:
d := [...]int{10, 20, 30}  // Type is [3]int (compiler counts)

// Access by index (0-based):
fmt.Println(b[0])   // 1
fmt.Println(b[4])   // 5
b[2] = 99

// Length:
fmt.Println(len(b)) // 5

// Iterate:
for i, v := range b {
    fmt.Printf("b[%d] = %d\n", i, v)
}
```

**Arrays are values** — when you assign or pass them, you get a copy:
```go
a := [3]int{1, 2, 3}
b := a              // b is a COPY of a
b[0] = 99
fmt.Println(a[0])   // 1 (unchanged — a and b are independent)
fmt.Println(b[0])   // 99
```

**Arrays are comparable:**
```go
a := [3]int{1, 2, 3}
b := [3]int{1, 2, 3}
c := [3]int{4, 5, 6}
fmt.Println(a == b)  // true
fmt.Println(a == c)  // false
```

**When to use arrays:** Rarely. Use them when you need a fixed-size collection that is part of a type definition (like SHA-256 hash: `[32]byte`), or for small temporary buffers. For almost everything else, use slices.

### Quick Check
> 1. What is the difference between `[5]int` and `[10]int` in Go?
> 2. Are arrays passed by value or by reference in Go?
> 3. When would you use an array instead of a slice?

---

## 2. Slices — Dynamic Collections

A slice is a **view into an array** that can grow. Slices are used for almost all list/sequence operations in Go:

```go
// Slice declaration (no size in brackets — that's the key difference)
var s []int               // nil slice (not the same as empty!)
t := []int{1, 2, 3, 4, 5}  // Slice literal

// Length and capacity:
fmt.Println(len(t))  // 5 — number of elements
fmt.Println(cap(t))  // 5 — total capacity of underlying array

// Access same as array:
t[0] = 10
fmt.Println(t[2])    // 3

// Append: add elements (may allocate new backing array)
t = append(t, 6, 7, 8)
fmt.Println(t)       // [10 2 3 4 5 6 7 8]
fmt.Println(len(t))  // 8
fmt.Println(cap(t))  // 10 (Go doubled capacity from 5 to 10)
```

**Nil slice vs empty slice:**
```go
var nilSlice []int    // nil — len=0, cap=0, == nil
emptySlice := []int{} // empty — len=0, cap=0, != nil

fmt.Println(nilSlice == nil)    // true
fmt.Println(emptySlice == nil)  // false
fmt.Println(len(nilSlice))      // 0 (safe!)
fmt.Println(len(emptySlice))    // 0

// Both are safe to range over:
for _, v := range nilSlice { _ = v }  // No iterations, no panic

// append works on nil slices:
nilSlice = append(nilSlice, 1, 2, 3)
fmt.Println(nilSlice)  // [1 2 3]
```

**Why nil vs empty matters:**
```go
// In JSON encoding:
var nilSlice []int
emptySlice := []int{}

json.Marshal(nilSlice)    // "null"
json.Marshal(emptySlice)  // "[]"

// API difference: "no items" vs "empty list" — semantically different!
```

### Quick Check
> 1. What is the difference between `var s []int` and `s := []int{}`?
> 2. What is the difference between `len` and `cap` for a slice?
> 3. What does `json.Marshal` produce for a nil slice vs an empty slice?

---

## 3. Slice Internals — The Three-Field Header

Understanding slice internals prevents subtle bugs:

```go
// A slice is a struct with three fields:
type slice struct {
    ptr *T   // pointer to first element in backing array
    len int  // number of elements visible
    cap int  // total size of backing array from ptr
}
```

**Visualized:**
```
Backing array:
index:  0    1    2    3    4    5    6    7
value: [10]  [20] [30] [40] [50] [60] [70] [80]
                   ↑                        ↑
                   |____ slice s___________|
                   ptr=&array[2]
                   len=4 (sees [30, 40, 50, 60])
                   cap=6 (from index 2 to end)
```

**Creating a slice from array:**
```go
array := [8]int{10, 20, 30, 40, 50, 60, 70, 80}

s := array[2:6]  // Slice: ptr=&array[2], len=4, cap=6
fmt.Println(s)          // [30 40 50 60]
fmt.Println(len(s), cap(s))  // 4 6

// Modifying slice modifies the underlying array!
s[0] = 99
fmt.Println(array[2])   // 99 — array was modified!
```

**When append allocates a new array:**
```go
s := make([]int, 3, 5)  // len=3, cap=5

// Append within capacity — same backing array:
s = append(s, 4)         // len=4, cap=5 — no new allocation
s = append(s, 5)         // len=5, cap=5 — no new allocation

// Append exceeds capacity — NEW backing array:
s = append(s, 6)         // len=6, cap=10 — new array! (cap doubled)
// s no longer shares memory with the original array
```

**Go's capacity growth strategy:**
- If cap < 256: double it
- If cap >= 256: grow by ~25%
- (Exact algorithm varies by Go version)

**Slice header is copied, not the array:**
```go
a := []int{1, 2, 3}
b := a  // b gets a COPY of the header (same ptr, len, cap)

b[0] = 99
fmt.Println(a[0])  // 99 — both point to same array!

// But after append that reallocates:
b = append(b, 4, 5, 6, 7, 8, 9, 10)  // Exceeds capacity, new array
b[0] = 42
fmt.Println(a[0])  // 99 — a still points to old array
```

### Quick Check
> 1. What are the three fields in a slice header?
> 2. What happens to the backing array when append exceeds capacity?
> 3. If two slices share the same backing array, what happens when you modify one?

---

## 4. Creating and Growing Slices

**`make([]T, len, cap)`** — create with specific length and capacity:
```go
// make(type, length, capacity)
s1 := make([]int, 5)        // len=5, cap=5, all zeros
s2 := make([]int, 0, 10)    // len=0, cap=10 (pre-allocated)
s3 := make([]int, 3, 10)    // len=3, cap=10

fmt.Println(s1)     // [0 0 0 0 0]
fmt.Println(len(s2), cap(s2))  // 0 10
```

**When to pre-allocate with `make`:**
```go
// BAD: Append in a loop without pre-allocation
// Each append may allocate a new backing array
result := []int{}
for i := 0; i < 1000; i++ {
    result = append(result, i)
}

// GOOD: Pre-allocate when you know the size
result := make([]int, 0, 1000)  // Pre-allocate 1000-element capacity
for i := 0; i < 1000; i++ {
    result = append(result, i)
}
// No reallocations! Much faster.
```

**`append` in detail:**
```go
// Append one element
s = append(s, 1)

// Append multiple elements
s = append(s, 1, 2, 3)

// Append another slice (use ... to spread)
other := []int{4, 5, 6}
s = append(s, other...)

// Append to nil slice (works fine!)
var s []int
s = append(s, 1, 2, 3)  // [1 2 3]
```

**`copy(dst, src)` — copy elements between slices:**
```go
src := []int{1, 2, 3, 4, 5}
dst := make([]int, 3)  // Only copies min(len(dst), len(src)) elements

n := copy(dst, src)
fmt.Println(n, dst)  // 3 [1 2 3]

// Full copy:
fullCopy := make([]int, len(src))
copy(fullCopy, src)
// Now fullCopy is independent from src
```

**Growing a slice manually:**
```go
// Double capacity when limit is reached:
func growSlice(s []int) []int {
    newCap := cap(s) * 2
    if newCap == 0 { newCap = 1 }
    grown := make([]int, len(s), newCap)
    copy(grown, s)
    return grown
}
```

### Quick Check
> 1. What is the difference between `make([]int, 5)` and `make([]int, 0, 5)`?
> 2. Why is pre-allocating with `make([]T, 0, n)` faster than starting with `[]T{}`?
> 3. What does `copy` return and what does it copy?

---

## 5. Slicing — Cutting Sub-Slices

You can create a sub-slice using the slice expression `s[low:high]`:

```go
s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}
//         0  1  2  3  4  5  6  7  8  9

// s[low:high] — from index low (inclusive) to high (exclusive)
a := s[2:5]   // [2, 3, 4]  len=3, cap=8
b := s[:3]    // [0, 1, 2]  (low defaults to 0)
c := s[7:]    // [7, 8, 9]  (high defaults to len)
d := s[:]     // [0..9]     (full copy of header, same array)

fmt.Println(a, len(a), cap(a))  // [2 3 4] 3 8
```

**Three-index slicing** — control the capacity:
```go
s := []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}

// s[low:high:max] — cap = max - low
a := s[2:5:5]  // [2, 3, 4]  len=3, cap=3
// Appending to a now allocates a NEW backing array,
// so it can no longer overwrite s[5], s[6]
```

**Deleting elements:**
```go
s := []int{1, 2, 3, 4, 5}

// Delete element at index 2 (value 3):
i := 2
s = append(s[:i], s[i+1:]...)
fmt.Println(s)  // [1 2 4 5]
```

**Inserting elements:**
```go
s := []int{1, 2, 4, 5}

// Insert 3 at index 2 (the safe manual way):
i := 2
s = append(s, 0)      // Grow the slice by one (value doesn't matter)
copy(s[i+1:], s[i:])  // Shift elements right to make room
s[i] = 3
fmt.Println(s)  // [1 2 3 4 5]

// Simpler with slices package (Go 1.21+):
import "slices"
s = slices.Insert(s, 2, 3)
```

### Quick Check
> 1. What does `s[2:5]` return for `s = []int{0,1,2,3,4,5}`?
> 2. What is three-index slicing `s[2:5:7]` and what does it control?
> 3. How do you delete the element at index `i` from a slice?

---

## 6. Common Slice Operations

**The `slices` package (Go 1.21+)** — standard library for slice operations:
```go
import "slices"

s := []int{3, 1, 4, 1, 5, 9, 2, 6}

// Sort:
slices.Sort(s)
fmt.Println(s)  // [1 1 2 3 4 5 6 9]

// Binary search (slice must be sorted):
idx, found := slices.BinarySearch(s, 5)
fmt.Println(idx, found)  // 5 true

// Contains:
fmt.Println(slices.Contains(s, 4))  // true

// Index of element:
fmt.Println(slices.Index(s, 4))  // 4 (-1 if not found)

// Reverse:
slices.Reverse(s)
fmt.Println(s)  // [9 6 5 4 3 2 1 1]

// Remove duplicates (sorted slice):
slices.Sort(s)
s = slices.Compact(s)
fmt.Println(s)  // [1 2 3 4 5 6 9]

// Max and Min:
fmt.Println(slices.Max(s))  // 9
fmt.Println(slices.Min(s))  // 1
```

**Before Go 1.21 — `sort` package:**
```go
import "sort"

// Sort ints:
sort.Ints(s)

// Sort strings:
words := []string{"banana", "apple", "cherry"}
sort.Strings(words)

// Custom sort:
sort.Slice(s, func(i, j int) bool {
    return s[i] > s[j]  // Descending order
})

// Stable sort (preserves order of equal elements):
sort.SliceStable(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
```

**Useful patterns:**

*Stack (using slice):*
```go
stack := []int{}
stack = append(stack, 1)         // push
stack = append(stack, 2)         // push
top := stack[len(stack)-1]       // peek
stack = stack[:len(stack)-1]     // pop
```

*Queue (circular buffer is better for production):*
```go
queue := []int{}
queue = append(queue, 1)         // enqueue
queue = append(queue, 2)         // enqueue
front := queue[0]                // peek front
queue = queue[1:]                // dequeue
```

*Deduplication:*
```go
func dedupe(s []int) []int {
    seen := make(map[int]bool)
    result := make([]int, 0, len(s))
    for _, v := range s {
        if !seen[v] {
            seen[v] = true
            result = append(result, v)
        }
    }
    return result
}
```

*Flatten 2D slice:*
```go
func flatten(s [][]int) []int {
    total := 0
    for _, inner := range s {
        total += len(inner)
    }
    result := make([]int, 0, total)
    for _, inner := range s {
        result = append(result, inner...)
    }
    return result
}
```

### Quick Check
> 1. What package in Go 1.21+ provides modern slice operations?
> 2. How do you sort a slice of integers?
> 3. How do you implement a stack using a slice?

---

## 7. 2D Slices

```go
// 2D slice (slice of slices):
matrix := [][]int{
    {1, 2, 3},
    {4, 5, 6},
    {7, 8, 9},
}

fmt.Println(matrix[1][2])  // 6

// Access row:
row := matrix[0]  // [1, 2, 3]

// Creating a 2D slice programmatically:
rows, cols := 3, 4
grid := make([][]int, rows)
for i := range grid {
    grid[i] = make([]int, cols)
}
grid[0][0] = 1
grid[2][3] = 99
```

**Jagged (non-rectangular) 2D slices:**
```go
// Each row can have different length:
triangle := [][]int{
    {1},
    {1, 2},
    {1, 2, 3},
}
```

### Quick Check
> 1. How do you declare a 2D slice in Go?
> 2. How do you create an m×n grid programmatically?
> 3. Can rows in a 2D slice have different lengths?

---

## 8. Common Pitfalls

**Pitfall 1: Modifying a slice you got from another function:**
```go
// You don't know if the caller will be surprised by modifications
func process(data []byte) {
    data[0] = 0xFF  // Modifies the original!
}

// Safe: make a copy if you need to modify
func processSafely(data []byte) {
    local := make([]byte, len(data))
    copy(local, data)
    local[0] = 0xFF  // Only local copy modified
}
```

**Pitfall 2: Memory leak with sub-slicing:**
```go
// Large slice backing array stays alive as long as ANY sub-slice references it
large := make([]byte, 1_000_000)  // 1MB
small := large[:10]  // Header has small len but large cap -- 1MB still in memory!

// Fix: copy only what you need
small = append([]byte(nil), large[:10]...)  // 10 bytes, 1MB freed
```

**Pitfall 3: Append to a sub-slice:**
```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3]  // b = [2, 3], cap=4

b = append(b, 99)  // Writes 99 to a[3]!
fmt.Println(a)     // [1 2 3 99 5] — a was silently modified!

// Fix: use three-index slice to limit capacity
b := a[1:3:3]  // cap=2, append will allocate new array
b = append(b, 99)
fmt.Println(a)     // [1 2 3 4 5] — unchanged
```

**Pitfall 4: Range variable captures:**
```go
// s[i] is a copy in range
s := []int{1, 2, 3}
for _, v := range s {
    v *= 2  // Modifies the copy, not the original
}
fmt.Println(s)  // [1 2 3] unchanged

// Fix: modify via index
for i := range s {
    s[i] *= 2
}
fmt.Println(s)  // [2 4 6]
```

### Quick Check
> 1. Why can a small sub-slice cause a memory leak?
> 2. What unexpected thing happens when you `append` to a sub-slice within capacity?
> 3. Does modifying `v` in `for _, v := range s` modify the original slice?

---

## Summary

- **Arrays**: fixed-size `[N]T`; value type (copied on assign); comparable with `==`; rarely used directly
- **Slices**: dynamic views `[]T`; three-field header (ptr, len, cap); used everywhere
- **nil vs empty**: `var s []int` is nil; `[]int{}` is empty, not nil; both have len 0; JSON differs
- **Internals**: slice shares backing array; modifications visible through all slices sharing it
- **make**: `make([]T, len, cap)` — pre-allocate for performance
- **append**: may allocate new backing array when cap exceeded; always reassign result
- **copy**: explicit copy between slices; returns number of elements copied
- **Slicing**: `s[low:high]`; three-index `s[low:high:max]` to control capacity
- **Operations**: `slices` package (1.21+) for sort, search, dedup; stack/queue patterns
- **Pitfalls**: sub-slice memory leaks, append to sub-slice clobbers, range copies

---

## Exercises

### Easy
1. Write a function `reverse(s []int) []int` that returns a new slice with elements in reverse order (do not modify the original).
2. Write a function `removeDuplicates(s []int) []int` that removes duplicate values while preserving order.
3. Write a function `flatten(s [][]int) []int` that flattens a 2D slice into a 1D slice.

### Medium
4. Sliding window maximum: Given a slice of integers and a window size k, find the maximum value in each window. Return a slice of maximums. Example: `[1,3,-1,-3,5,3,6,7]`, k=3 → `[3,3,5,5,6,7]`. Implement in O(n) time using a deque (slice as deque).
5. Rotate slice: Implement `rotate(s []int, k int)` that rotates the slice k positions to the right in-place. `[1,2,3,4,5]` rotated 2 → `[4,5,1,2,3]`. Implement in O(n) time, O(1) extra space. Hint: three-step reverse.
6. Dutch national flag: Given a slice of 0s, 1s, and 2s (unsorted), sort them in-place in O(n) time and O(1) space so all 0s come first, then 1s, then 2s. Use only swaps. This is a classic Dijkstra partition algorithm.

### Hard
7. Dynamic array implementation: Implement a generic dynamic array (ArrayList) in Go using generics. It must: store any type `T`, support `Append(v T)`, `Get(i int) T`, `Set(i int, v T)`, `Delete(i int)`, `Insert(i int, v T)`, `Len() int`, `Cap() int`, internal resizing (double when full, halve when 25% full), panic with a meaningful message on out-of-bounds access. Write benchmarks comparing it to using Go's built-in slice directly.
8. Merge k sorted slices: Given k sorted slices of integers, merge them into a single sorted slice. Naive O(n×k) approach: merge pairwise. Optimal approach using a min-heap: O(n log k) where n is total elements. Implement both and benchmark them with k=100 slices each containing 10,000 elements. Bonus: implement a streaming version that returns an iterator (channel) rather than materializing the whole result.
