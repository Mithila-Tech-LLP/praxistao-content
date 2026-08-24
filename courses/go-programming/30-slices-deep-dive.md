# Chapter 30: Slices Deep Dive — Internals, Patterns, and Pitfalls

Slices are Go's most versatile built-in. Most Go bugs involving data corruption, unexpected sharing, or memory leaks trace back to a misunderstanding of how slices work internally. This chapter fixes that.

## Table of Contents

1. [Slice Header and Backing Array](#1-slice-header-and-backing-array)
2. [Append and Growth](#2-append-and-growth)
3. [Slice of Slice — Sharing](#3-slice-of-slice--sharing)
4. [copy and When You Need It](#4-copy-and-when-you-need-it)
5. [Common Patterns](#5-common-patterns)
6. [Common Pitfalls](#6-common-pitfalls)
7. [2D Slices](#7-2d-slices)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Slice Header and Backing Array

A slice is a three-field struct — **not** a pointer to an array:

```go
// The runtime representation:
type slice struct {
    array unsafe.Pointer  // pointer to the backing array
    len   int             // number of elements in use
    cap   int             // total capacity of the backing array
}
```

```go
// Visualizing a slice:
s := []int{10, 20, 30, 40, 50}

// Memory layout:
// s.array → [10][20][30][40][50]
// s.len   = 5
// s.cap   = 5
```

This is why slices are passed by value but still mutate the original — the copy of the header points to the same backing array.

```go
func double(nums []int) {
    for i := range nums {
        nums[i] *= 2       // modifies the backing array directly
    }
}

orig := []int{1, 2, 3}
double(orig)
fmt.Println(orig)  // [2 4 6] — modified!
```

---

## 2. Append and Growth

`append` adds elements to a slice, growing the backing array if necessary.

```go
s := make([]int, 3, 5)  // len=3, cap=5, backing array has 5 slots
s = append(s, 100)       // len=4, cap=5, same backing array
s = append(s, 200)       // len=5, cap=5, same backing array
s = append(s, 300)       // cap exceeded! new backing array allocated
                         // len=6, cap=10 (roughly doubled)
```

**Growth behavior in Go 1.18+:** starts at 2× for small slices, transitions to ~1.25× for large slices. The exact formula is in `runtime/slice.go`.

```go
// Verify growth with reflect:
import "fmt"

s := make([]int, 0)
prev := 0
for i := 0; i < 20; i++ {
    s = append(s, i)
    if cap(s) != prev {
        fmt.Printf("len=%d cap=%d\n", len(s), cap(s))
        prev = cap(s)
    }
}
// len=1 cap=1
// len=2 cap=2
// len=3 cap=4
// len=5 cap=8
// len=9 cap=16
// ...
```

### Pre-allocate when you know the size

```go
// Bad: append repeatedly reallocates
result := []int{}
for _, n := range input {
    result = append(result, n*2)
}

// Good: allocate once
result := make([]int, len(input))
for i, n := range input {
    result[i] = n * 2
}

// Also good: pre-allocate capacity, use append
result := make([]int, 0, len(input))
for _, n := range input {
    result = append(result, n*2)
}
```

---

## 3. Slice of Slice — Sharing

Slicing a slice does **not** copy the data. The sub-slice shares the same backing array.

```go
original := []int{1, 2, 3, 4, 5}
sub := original[1:4]  // [2, 3, 4]

// Both point to the same memory:
// original: ptr=X, len=5, cap=5
// sub:      ptr=X+1, len=3, cap=4 (4 remaining from position 1)

sub[0] = 99
fmt.Println(original)  // [1 99 3 4 5] — original modified!
```

This is intentional and efficient — slicing is O(1) with no allocation.

### Full slice expression — capping capacity

The three-index slice `a[low:high:max]` limits the capacity:

```go
a := []int{1, 2, 3, 4, 5}
b := a[1:3:3]  // len=2, cap=2 (limited to position 3)

// b = [2, 3], cap=2
// Appending to b won't affect a — it will allocate a new backing array
b = append(b, 99)  // new backing array, a is untouched
```

This is the right way to pass a sub-slice to a function that might append to it.

---

## 4. copy and When You Need It

`copy(dst, src)` copies elements between slices. It copies `min(len(dst), len(src))` elements.

```go
src := []int{1, 2, 3, 4, 5}

// Clone a slice (independent copy):
dst := make([]int, len(src))
copy(dst, src)
dst[0] = 99
fmt.Println(src) // [1 2 3 4 5] — unaffected

// Idiomatic one-liner clone:
dst = append([]int(nil), src...)  // or append([]int{}, src...)
```

### When you must copy

1. **Returning a sub-slice from a function** — the caller gets a reference to the whole backing array, preventing GC of large data.
2. **Growing a slice for concurrent access** — two goroutines cannot safely append to the same backing array.
3. **Isolating mutations** — when you need to modify a copy without affecting the original.

```go
// Memory leak: returning a small slice retains the big backing array
func parseFirst(bigData []byte) []byte {
    return bigData[:10]  // BAD: bigData never GC'd
}

func parseFirstSafe(bigData []byte) []byte {
    result := make([]byte, 10)
    copy(result, bigData[:10])
    return result  // bigData can now be GC'd
}
```

---

## 5. Common Patterns

### Filter in-place

```go
// Remove evens without allocation (reuse backing array)
func filterOdd(nums []int) []int {
    result := nums[:0]  // len=0, cap=len(nums), same backing array
    for _, n := range nums {
        if n%2 != 0 {
            result = append(result, n)  // writes into original backing array
        }
    }
    return result
}

nums := []int{1, 2, 3, 4, 5, 6, 7}
fmt.Println(filterOdd(nums))  // [1 3 5 7]
```

### Delete element at index i (order matters)

```go
// Maintain order: O(n) — shifts elements left
func deleteAt(s []int, i int) []int {
    return append(s[:i], s[i+1:]...)
}

// Don't maintain order: O(1) — swap with last
func deleteAtFast(s []int, i int) []int {
    s[i] = s[len(s)-1]
    return s[:len(s)-1]
}
```

### Insert element at index i

```go
func insertAt(s []int, i int, v int) []int {
    s = append(s, 0)          // grow by one (may reallocate)
    copy(s[i+1:], s[i:])      // shift right
    s[i] = v
    return s
}
```

### Flatten (2D → 1D)

```go
func flatten(matrix [][]int) []int {
    total := 0
    for _, row := range matrix { total += len(row) }
    
    result := make([]int, 0, total)  // pre-allocate
    for _, row := range matrix {
        result = append(result, row...)
    }
    return result
}
```

### Chunking (split into batches)

```go
func chunk[T any](s []T, size int) [][]T {
    if size <= 0 { return nil }
    chunks := make([][]T, 0, (len(s)+size-1)/size)
    for size < len(s) {
        s, chunks = s[size:], append(chunks, s[:size:size]) // cap limits each chunk
    }
    return append(chunks, s)
}
```

### Deduplication

```go
func dedupe[T comparable](s []T) []T {
    seen := make(map[T]struct{}, len(s))
    result := make([]T, 0, len(s))
    for _, v := range s {
        if _, ok := seen[v]; !ok {
            seen[v] = struct{}{}
            result = append(result, v)
        }
    }
    return result
}
```

---

## 6. Common Pitfalls

### Pitfall 1: Append sharing

```go
base := []int{1, 2, 3}
a := append(base, 4)  // may share backing array with base
b := append(base, 5)  // if cap(base) >= 4, writes to same position as a

// Safe pattern: copy before branching
a := append(append([]int(nil), base...), 4)
b := append(append([]int(nil), base...), 5)
```

### Pitfall 2: Loop variable capture

```go
// Bug: all elements capture the same i
nums := []int{1, 2, 3}
funcs := make([]func(), len(nums))
for i, n := range nums {
    i, n := i, n  // shadow: create new variables per iteration
    funcs[i] = func() { fmt.Println(n) }
}
// Without shadowing, all would print 3 (last value of n in Go < 1.22)
// In Go 1.22+, loop variables are scoped per iteration automatically
```

### Pitfall 3: nil vs empty slice

```go
var nilSlice []int      // nil — no backing array
emptySlice := []int{}   // not nil — has backing array (len=0, cap=0)

// Both have len=0 and work with range, append, copy
// The difference matters for JSON:
json.Marshal(nilSlice)    // → null
json.Marshal(emptySlice)  // → []

// Best practice: return nil for "no data", [] for "empty collection"
```

### Pitfall 4: Retaining large backing arrays

```go
data := make([]int, 1_000_000)
// ... fill data ...

// Bug: result retains the entire 1M-element backing array
result := data[:10]

// Fix: copy out the elements you need
result := make([]int, 10)
copy(result, data[:10])
data = nil  // now the big array can be GC'd
```

---

## 7. 2D Slices

```go
// Method 1: independent rows (safe for independent mutation)
rows, cols := 3, 4
matrix := make([][]int, rows)
for i := range matrix {
    matrix[i] = make([]int, cols)
}

// Method 2: contiguous backing array (cache-friendly)
backing := make([]int, rows*cols)
matrix2 := make([][]int, rows)
for i := range matrix2 {
    matrix2[i] = backing[i*cols : (i+1)*cols]
}
// All rows share one backing array — better memory locality
```

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Slice header | 3 fields: pointer, len, cap — passed by value |
| Backing array | Shared between the original and any sub-slice |
| append | O(1) amortized; may allocate new backing array |
| copy | Always makes an independent copy; use to isolate mutations |
| `s[:0]` | Reuse backing array for filter-in-place |
| `a[lo:hi:max]` | Limit cap so append can't corrupt the parent |
| nil slice | `json.Marshal → null`; empty slice → `[]` |
| Memory leak | Sub-slice pins the entire backing array in memory |

---

## Exercises

### Easy
1. Write a function `removeDuplicatesSorted(nums []int) []int` that removes duplicates from a sorted slice in-place (O(1) space, O(n) time). The result should be a prefix of the original backing array.
2. Predict and verify: `a := []int{1,2,3,4,5}; b := a[2:4]; b[0] = 99; fmt.Println(a)`. Then do `b = append(b, 100, 200, 300); b[0] = 77; fmt.Println(a)`. Explain why the second modification does not affect `a`.
3. Implement `Rotate(s []int, k int)` that rotates a slice left by k positions in-place. Example: `[1,2,3,4,5]` rotated by 2 → `[3,4,5,1,2]`. Use the three-reversal trick: O(n) time, O(1) space.

### Medium
4. Implement a generic `Stack[T]` backed by a slice. Add `Peek()`, `Push()`, `Pop()`, and `Len()`. Ensure that `Peek()` and `Pop()` return false instead of panicking on empty stacks. Benchmark push + pop of 1M elements and verify amortized O(1) per operation.
5. Implement `chunk[T any](s []T, size int) [][]T` using the three-index slice expression so each chunk has an independent capacity of exactly `size`. Write a test that verifies appending to one chunk does not affect adjacent chunks.
6. Given a 2D matrix of ints, implement `transpose(m [][]int) [][]int` that returns the transpose. Then implement an in-place version (only valid for square matrices). Compare memory usage.

### Hard
7. Implement a **sliding window** function `SlidingWindowMax(nums []int, k int) []int` that returns the max of each window of size k in O(n) time using a monotonic deque (backed by a slice). Verify that naive O(nk) and this O(n) solution produce identical results for all test cases.
8. Design a **safe sub-slice function** that prevents backing-array memory leaks: `SafeSlice(s []T, lo, hi int) []T` should return a copy if `hi-lo < len(s)/2` (i.e., if you're taking less than half). Write a benchmark showing that the leak-prevention copy becomes worthwhile when the original is large and short-lived.
