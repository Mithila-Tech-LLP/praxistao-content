# Chapter 22: Sorting Algorithms — Ordering the World

> "Sorting is the fundamental act of imposing order on chaos. Almost every other algorithm assumes the world is already sorted."

---

## Overview

Sorting is one of the most studied problems in computer science, and for good reason. More than 25% of all computer time was historically spent sorting data. Every database you query, every search you perform, every autocomplete suggestion you see — they all depend on data being sorted. Binary search (which we cover in Chapter 23) requires sorted data. Database indexes are sorted. Your file system lists directories in sorted order.

This chapter takes you through the major sorting algorithms from the simplest (bubble sort) to the ones actually used in production (quicksort, Timsort). We will implement merge sort in Go step by step, understand why some sorts are stable and others are not, and see exactly where the Astra compiler needs sorting and why.

## What We Are Building

By the end of this chapter, you will understand the trade-offs between different sorting algorithms, be able to implement merge sort from scratch, know why the theoretical lower bound for comparison-based sorting is O(n log n), and see how Go's standard library sort works. We will also implement sorting utilities that the Astra compiler uses: sorting error messages by line number, struct fields by offset, and import paths alphabetically.

---

## Table of Contents

1. Why Sorting Matters
2. Stability in Sorting
3. Bubble Sort — Simple but Slow
4. Selection Sort
5. Insertion Sort — Small Arrays' Champion
6. Merge Sort — Divide and Conquer
7. Quicksort — The Practical King
8. Heap Sort
9. Counting Sort — Breaking the Lower Bound
10. Timsort — The Algorithm Python and Java Use
11. The Comparison-Based Sorting Lower Bound
12. Go's sort.Slice
13. Astra Build Milestone
14. Exercises
15. Summary

---

## 1. Why Sorting Matters

Let us be concrete about where sorting appears:

**Binary search requires sorted data.** If you want to find a name in a phone book in O(log n) time, the phone book must be sorted. Without sorting, you are stuck with O(n) linear search.

**Database indexes are B-trees — sorted structures.** When you write `SELECT * FROM users WHERE name = 'Alice'`, the database finds Alice in O(log n) time because the index is sorted.

**Sorting removes duplicates efficiently.** Sort first, then scan for adjacent duplicates — O(n log n) instead of O(n²).

**Compilers sort for deterministic output.** If your compiler generates different output each run (because map iteration order is random), debugging is impossible. The Astra compiler sorts struct fields, function names, and error messages for consistent output.

---

## 2. Stability in Sorting

A sorting algorithm is **stable** if it preserves the original relative order of equal elements.

```
Input: [(Alice, 30), (Bob, 25), (Carol, 30), (Dave, 25)]
Sorted by age:

Stable result:   [(Bob, 25), (Dave, 25), (Alice, 30), (Carol, 30)]
                   ──────Bob before Dave (original order)──────

Unstable result: [(Dave, 25), (Bob, 25), (Carol, 30), (Alice, 30)]
                   ──────Dave before Bob (order changed)──────
```

Stability matters when you sort by multiple keys. For example: sort by department first, then by name. If the name sort is stable, the department order is preserved within each name. If unstable, you need to do a single combined comparison.

| Algorithm      | Stable? | In-Place? | Best       | Average    | Worst      |
|----------------|---------|-----------|------------|------------|------------|
| Bubble Sort    | Yes     | Yes       | O(n)       | O(n²)      | O(n²)      |
| Selection Sort | No      | Yes       | O(n²)      | O(n²)      | O(n²)      |
| Insertion Sort | Yes     | Yes       | O(n)       | O(n²)      | O(n²)      |
| Merge Sort     | Yes     | No        | O(n log n) | O(n log n) | O(n log n) |
| Quicksort      | No      | Yes       | O(n log n) | O(n log n) | O(n²)      |
| Heap Sort      | No      | Yes       | O(n log n) | O(n log n) | O(n log n) |
| Counting Sort  | Yes     | No        | O(n+k)     | O(n+k)     | O(n+k)     |
| Timsort        | Yes     | No        | O(n)       | O(n log n) | O(n log n) |

---

## 3. Bubble Sort — Simple but Slow

Bubble sort compares adjacent elements and swaps them if they are in the wrong order. The largest element "bubbles up" to its correct position in each pass.

```
Pass 1: [5, 3, 8, 1, 9, 2]
Compare 5,3: swap → [3, 5, 8, 1, 9, 2]
Compare 5,8: ok   → [3, 5, 8, 1, 9, 2]
Compare 8,1: swap → [3, 5, 1, 8, 9, 2]
Compare 8,9: ok   → [3, 5, 1, 8, 9, 2]
Compare 9,2: swap → [3, 5, 1, 8, 2, 9]  ← 9 is now in place!

Pass 2: [3, 5, 1, 8, 2, 9]  (ignore last element, already sorted)
Compare 3,5: ok   → [3, 5, 1, 8, 2, 9]
Compare 5,1: swap → [3, 1, 5, 8, 2, 9]
Compare 5,8: ok   → [3, 1, 5, 8, 2, 9]
Compare 8,2: swap → [3, 1, 5, 2, 8, 9]  ← 8 is now in place!

...and so on until sorted
```

```go
// Go implementation of bubble sort
func bubbleSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        swapped := false
        for j := 0; j < n-i-1; j++ {
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
                swapped = true
            }
        }
        // Optimization: if no swaps in a pass, already sorted
        if !swapped {
            break  // This makes best case O(n) for already-sorted input
        }
    }
}
```

**Analysis:**
- Outer loop runs n-1 times
- Inner loop runs n-1, n-2, n-3, ... 1 times
- Total comparisons: (n-1) + (n-2) + ... + 1 = n(n-1)/2 = O(n²)
- Space: O(1) — in-place
- Stable: yes (we only swap when strictly greater, equal elements stay in order)

**When to use:** Never in production. Bubble sort exists to teach the concept of comparison-based sorting. Its O(n²) complexity makes it impractical for any serious use.

---

## 4. Selection Sort

Selection sort finds the minimum element in the unsorted portion and places it at the beginning.

```
[64, 25, 12, 22, 11]
Find min in [64,25,12,22,11] = 11, swap with position 0:
[11, 25, 12, 22, 64]
Find min in [25,12,22,64] = 12, swap with position 1:
[11, 12, 25, 22, 64]
Find min in [25,22,64] = 22, swap with position 2:
[11, 12, 22, 25, 64]
Find min in [25,64] = 25, already in place:
[11, 12, 22, 25, 64]  ← sorted!
```

```go
func selectionSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        minIdx := i
        for j := i + 1; j < n; j++ {
            if arr[j] < arr[minIdx] {
                minIdx = j
            }
        }
        arr[i], arr[minIdx] = arr[minIdx], arr[i]
    }
}
```

Selection sort is O(n²) in all cases — it always scans the entire unsorted portion even if already sorted. It makes at most O(n) swaps (fewer than bubble sort's O(n²) swaps), which can be useful when swapping is expensive. It is not stable.

---

## 5. Insertion Sort — Small Arrays' Champion

Insertion sort builds the sorted array one element at a time, inserting each new element into its correct position among the already-sorted elements.

Think of sorting a hand of playing cards. You pick up cards one at a time and insert each into the right position in your hand.

```
[5, 2, 4, 6, 1, 3]
Start: sorted portion = [5]

Take 2: find its position in [5]
  2 < 5, move 5 right → insert 2 → [2, 5]

Take 4: find its position in [2, 5]
  4 < 5, move 5 right
  4 > 2, insert 4 → [2, 4, 5]

Take 6: find its position in [2, 4, 5]
  6 > 5, insert at end → [2, 4, 5, 6]

Take 1: find its position in [2, 4, 5, 6]
  1 < 6, 1 < 5, 1 < 4, 1 < 2, insert at beginning → [1, 2, 4, 5, 6]

Take 3: find its position in [1, 2, 4, 5, 6]
  3 < 5, 3 < 4, 3 > 2, insert → [1, 2, 3, 4, 5, 6]  ← sorted!
```

```go
func insertionSort(arr []int) {
    for i := 1; i < len(arr); i++ {
        key := arr[i]
        j := i - 1
        // Shift elements that are greater than key one position to the right
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]
            j--
        }
        arr[j+1] = key
    }
}
```

**Why insertion sort is used in practice for small arrays:**

1. O(n) for nearly-sorted data (inner loop runs very few times)
2. Very small constant factor — the inner loop is extremely tight
3. Cache-friendly memory access patterns
4. Stable and in-place

Timsort (Python and Java's sort) uses insertion sort for small subarrays (32 elements or fewer). Go's standard library also falls back to insertion sort for small slices.

---

## 6. Merge Sort — Divide and Conquer

Merge sort is the classic divide-and-conquer sorting algorithm. It is elegant, provably O(n log n) in all cases, and stable.

**The idea:**
1. If the array has 0 or 1 elements, it is already sorted (base case)
2. Split the array in half
3. Recursively sort each half
4. Merge the two sorted halves into one sorted array

```
                [38, 27, 43, 3, 9, 82, 10]
                           │
               ┌───────────┴────────────┐
         [38, 27, 43, 3]          [9, 82, 10]
               │                        │
        ┌──────┴──────┐          ┌──────┴──────┐
    [38, 27]     [43, 3]      [9, 82]      [10]
        │             │          │
     ┌──┴──┐       ┌──┴──┐    ┌──┴──┐
   [38]  [27]   [43]   [3]  [9]  [82]
     │      │      │      │    │      │
     └──┬───┘      └──┬───┘   └──┬───┘
    [27, 38]     [3, 43]     [9, 82]
          └──────┬──────┘         │
           [3, 27, 38, 43]   [9, 10, 82]
                   └──────────┬──────────┘
                [3, 9, 10, 27, 38, 43, 82]
```

The merge step is the key. Given two sorted arrays, we can merge them in O(n) time:

```go
// Merge two sorted slices into one sorted slice
func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {   // <= for stability
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    
    // Append remaining elements (one of these will be empty)
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}

// Recursive merge sort
func mergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr  // base case
    }
    
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])   // sort left half
    right := mergeSort(arr[mid:])  // sort right half
    return merge(left, right)       // merge sorted halves
}
```

**Why is merge sort O(n log n)?**

```
Recursion depth: log₂(n)  (we halve the array each time)

Level 0: 1 array of size n     → O(n) merge work
Level 1: 2 arrays of size n/2  → O(n/2 + n/2) = O(n) merge work
Level 2: 4 arrays of size n/4  → O(n) merge work
...
Level log₂n: n arrays of size 1 → O(n) merge work

Total: O(n) work × log₂(n) levels = O(n log n)
```

**Weakness:** Merge sort requires O(n) extra space for the merged results. It is not in-place.

---

## 7. Quicksort — The Practical King

Quicksort is the most widely used sorting algorithm in practice. It is O(n log n) on average and sorts in-place (O(log n) extra space for the recursion stack).

**The idea:**
1. Choose a "pivot" element
2. Partition: move all elements smaller than pivot to the left, all larger to the right
3. The pivot is now in its final sorted position
4. Recursively sort the left and right partitions

```
[3, 6, 8, 10, 1, 2, 1]  — pick pivot = 1 (last element)

Partition around pivot 1:
Elements less than 1: (none)
Elements equal to 1: 1, 1
Elements greater than 1: 3, 6, 8, 10, 2

Result: [1, 1, 3, 6, 8, 10, 2]
              │  └──── recursively sort this
Pivot (1) is in its final position!

Recursively sort [3, 6, 8, 10, 2], pivot = 2:
Less than 2: (none)
Equal to 2: 2
Greater: 3, 6, 8, 10

Result: [2, 3, 6, 8, 10]

Final: [1, 1, 2, 3, 6, 8, 10]
```

```go
func quickSort(arr []int, low, high int) {
    if low < high {
        pivotIdx := partition(arr, low, high)
        quickSort(arr, low, pivotIdx-1)
        quickSort(arr, pivotIdx+1, high)
    }
}

// Lomuto partition scheme
func partition(arr []int, low, high int) int {
    pivot := arr[high]  // pivot is the last element
    i := low - 1        // i tracks the boundary of smaller elements
    
    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    
    arr[i+1], arr[high] = arr[high], arr[i+1]  // place pivot in final position
    return i + 1
}

// Usage: quickSort(arr, 0, len(arr)-1)
```

**The worst case problem:** If the pivot is always the smallest or largest element (e.g., sorted input with last-element pivot), quicksort degenerates to O(n²).

```
Sorted input: [1, 2, 3, 4, 5] with pivot = last element
Partition: pivot=5, [1,2,3,4] on left, nothing on right → 1 element reduction
Partition: pivot=4, [1,2,3] on left, nothing on right → 1 element reduction
...
n partitions of size n-1, n-2, ... → O(n²) total work!
```

**Solutions to worst-case quicksort:**
1. Random pivot selection: pick a random element as pivot
2. Median-of-three: pick the median of first, middle, and last elements
3. Use a different algorithm (introsort: quicksort + heapsort fallback)

Go's `sort.Slice` uses a combination of quicksort, heapsort, and insertion sort (pattern-defeating quicksort / pdqsort).

---

## 8. Heap Sort

Heap sort uses a binary heap data structure. It is O(n log n) in all cases (including worst case, unlike quicksort) and sorts in-place, but it is not stable and has poor cache performance.

```go
func heapSort(arr []int) {
    n := len(arr)
    
    // Build max-heap (rearrange array)
    for i := n/2 - 1; i >= 0; i-- {
        heapify(arr, n, i)
    }
    
    // Extract elements from heap one by one
    for i := n - 1; i > 0; i-- {
        arr[0], arr[i] = arr[i], arr[0]  // move current root to end
        heapify(arr, i, 0)                // heapify reduced heap
    }
}

func heapify(arr []int, n, i int) {
    largest := i
    left := 2*i + 1
    right := 2*i + 2
    
    if left < n && arr[left] > arr[largest] {
        largest = left
    }
    if right < n && arr[right] > arr[largest] {
        largest = right
    }
    if largest != i {
        arr[i], arr[largest] = arr[largest], arr[i]
        heapify(arr, n, largest)
    }
}
```

---

## 9. Counting Sort — Breaking the Lower Bound

All comparison-based sorts are at least O(n log n) in the worst case. But counting sort is O(n) — how?

The trick: counting sort does not compare elements. It counts how many of each value exist, then uses those counts to place elements directly.

**Constraint:** Only works when elements are integers in a known, bounded range [0, k].

```
Input: [4, 2, 2, 8, 3, 3, 1]  (all values in range [1, 8])

Step 1: Count occurrences
count[1]=1, count[2]=2, count[3]=2, count[4]=1, count[8]=1

Step 2: Output elements in order
1 appears 1 time → output: 1
2 appears 2 times → output: 1, 2, 2
3 appears 2 times → output: 1, 2, 2, 3, 3
4 appears 1 time  → output: 1, 2, 2, 3, 3, 4
8 appears 1 time  → output: 1, 2, 2, 3, 3, 4, 8
```

```go
func countingSort(arr []int, maxVal int) []int {
    count := make([]int, maxVal+1)
    
    // Count occurrences
    for _, v := range arr {
        count[v]++
    }
    
    // Build output
    result := make([]int, 0, len(arr))
    for val, cnt := range count {
        for i := 0; i < cnt; i++ {
            result = append(result, val)
        }
    }
    return result
}
// Time: O(n + k) where k = maxVal
// Space: O(n + k)
```

When k is O(n), counting sort is O(n). But when k is huge (e.g., sorting 64-bit integers), k dominates and it becomes impractical.

**Compiler use case:** Sorting error messages by line number. If a source file has n lines, error line numbers are in [1, n]. Counting sort gives O(n) sorting of error messages!

---

## 10. Timsort — The Algorithm Python and Java Use

Timsort is the algorithm used by Python (since 2002), Java (for objects since Java 7), Android, and many others. It is a hybrid of merge sort and insertion sort.

**Key insight:** Real-world data is rarely random. It often has "runs" — subsequences that are already sorted (ascending or descending).

Timsort exploits this:
1. Scan the array, identifying natural runs (already-sorted subsequences)
2. If a run is descending, reverse it to make it ascending
3. If a run is shorter than a minimum size (32 elements), extend it using insertion sort
4. Merge the runs using a merge sort-like process, but with special optimizations

```
Real-world data with natural runs:
[1, 3, 5, 7, 2, 4, 6, 9, 8, 10, 11, 12]
 └──run 1──┘  └──run 2──┘  └────run 3────┘

Timsort finds these runs and merges them:
Merge run1 + run2: [1, 2, 3, 4, 5, 6, 7, 9]
Merge result + run3: [1, 2, 3, 4, 5, 6, 7, 8, 9, 10, 11, 12]
```

Timsort's strengths:
- O(n) for already-sorted data
- O(n log n) worst case
- Stable
- Extremely cache-friendly due to run detection
- Handles partially sorted data much better than pure quicksort

---

## 11. The Comparison-Based Sorting Lower Bound

Here is a remarkable theorem: **Any comparison-based sorting algorithm must make at least O(n log n) comparisons in the worst case.**

This means merge sort, Timsort, and heapsort are optimal. You cannot do better with comparisons.

**The proof idea:** Think of sorting as navigating a decision tree.

```
              a[0] < a[1]?
             /             \
           Yes               No
    a[1] < a[2]?         a[0] < a[2]?
    /         \           /         \
  Yes         No        Yes          No
[1,2,3]    [1,3,2]   [2,1,3]     [2,3,1]
                                  ...
```

For n elements, there are n! possible orderings (permutations). Each is a leaf of this tree. A binary tree with n! leaves has depth at least log₂(n!).

By Stirling's approximation: log₂(n!) ≈ n log₂(n) - n log₂(e) = Ω(n log n)

So the depth (number of comparisons in worst case) is at least Ω(n log n). No comparison-based sort can be faster.

Non-comparison sorts (counting sort, radix sort, bucket sort) escape this bound by using element values directly instead of comparisons.

---

## 12. Go's sort.Slice

Go's standard library sorting is straightforward to use:

```go
import "sort"

// Sort integers
nums := []int{5, 2, 8, 1, 9}
sort.Ints(nums)  // [1, 2, 5, 8, 9]

// Sort strings
words := []string{"banana", "apple", "cherry"}
sort.Strings(words)  // ["apple", "banana", "cherry"]

// Sort custom types with sort.Slice
type Person struct {
    Name string
    Age  int
}
people := []Person{
    {"Alice", 30},
    {"Bob", 25},
    {"Carol", 35},
}

// Sort by age
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
// [{Bob 25} {Alice 30} {Carol 35}]

// Sort by name, stable
sort.SliceStable(people, func(i, j int) bool {
    return people[i].Name < people[j].Name
})
// [{Alice 30} {Bob 25} {Carol 35}]

// Check if sorted
isSorted := sort.IntsAreSorted(nums)  // true
```

Go's `sort.Slice` uses pdqsort (pattern-defeating quicksort), which combines:
- Quicksort for the general case (fast, in-place)
- Heapsort as a fallback when quicksort degrades (prevents O(n²) worst case)
- Insertion sort for small partitions (< 12 elements typically)

---

## Astra Build Milestone

The Astra compiler needs sorting in several concrete places. Let us implement them all.

```go
// File: compiler/sorting/compiler_sorts.go
package sorting

import (
    "fmt"
    "sort"
    "strings"
)

// ─── 1. SORT ERROR MESSAGES BY LINE NUMBER ───────────────────────────────────

// CompileError represents an error found during compilation
type CompileError struct {
    Line    int
    Column  int
    Message string
    File    string
}

func (e CompileError) String() string {
    return fmt.Sprintf("%s:%d:%d: error: %s", e.File, e.Line, e.Column, e.Message)
}

// SortErrors sorts compile errors by file, then line, then column.
// This ensures error messages are displayed in source order,
// making them easy to read and fix top-to-bottom.
func SortErrors(errors []CompileError) {
    sort.SliceStable(errors, func(i, j int) bool {
        ei, ej := errors[i], errors[j]
        // Primary sort: by file name
        if ei.File != ej.File {
            return ei.File < ej.File
        }
        // Secondary sort: by line number
        if ei.Line != ej.Line {
            return ei.Line < ej.Line
        }
        // Tertiary sort: by column
        return ei.Column < ej.Column
    })
}

// Example: unsorted errors from parallel compilation passes
func ExampleErrorSorting() {
    errors := []CompileError{
        {Line: 42, Column: 5, Message: "undefined variable 'x'", File: "main.astra"},
        {Line: 10, Column: 1, Message: "missing return statement", File: "utils.astra"},
        {Line: 42, Column: 1, Message: "type mismatch: expected int, got string", File: "main.astra"},
        {Line: 7, Column: 3, Message: "unknown type 'Foobar'", File: "main.astra"},
    }

    fmt.Println("Before sorting:")
    for _, e := range errors {
        fmt.Println(" ", e)
    }

    SortErrors(errors)

    fmt.Println("\nAfter sorting:")
    for _, e := range errors {
        fmt.Println(" ", e)
    }
}

// ─── 2. SORT STRUCT FIELDS FOR CONSISTENT MEMORY LAYOUT ──────────────────────

// StructField represents a field in an Astra struct
type StructField struct {
    Name   string
    Type   string
    Offset int // byte offset in memory layout
}

// SortStructFields sorts struct fields to minimize padding.
// In memory, fields must be aligned to their size.
// By sorting larger fields first, we reduce wasted space.
//
// Example:
//   struct Foo { a: bool; b: int64; c: bool }
//   Naive layout: bool(1) + 7 pad + int64(8) + bool(1) + 7 pad = 24 bytes
//   Sorted layout: int64(8) + bool(1) + bool(1) + 6 pad = 16 bytes
func SortStructFieldsByAlignment(fields []StructField) []StructField {
    typeSize := map[string]int{
        "int64":  8,
        "int32":  4,
        "int16":  2,
        "int8":   1,
        "bool":   1,
        "float64": 8,
        "float32": 4,
    }

    sorted := make([]StructField, len(fields))
    copy(sorted, fields)

    sort.SliceStable(sorted, func(i, j int) bool {
        si := typeSize[sorted[i].Type]
        sj := typeSize[sorted[j].Type]
        if si != sj {
            return si > sj  // larger types first (decreasing size)
        }
        return sorted[i].Name < sorted[j].Name  // alphabetical for equal sizes
    })

    // Assign offsets after sorting
    offset := 0
    for i := range sorted {
        size := typeSize[sorted[i].Type]
        // Align offset to field size
        if offset%size != 0 {
            offset += size - (offset % size)
        }
        sorted[i].Offset = offset
        offset += size
    }

    return sorted
}

// ─── 3. SORT IMPORT PATHS ALPHABETICALLY ─────────────────────────────────────

// ImportDecl represents an import statement in Astra
type ImportDecl struct {
    Path  string
    Alias string
}

// SortImports sorts import declarations alphabetically.
// Standard library imports come before third-party imports.
// This mirrors Go's import grouping convention and makes
// diffs cleaner (new imports are inserted in sorted position).
func SortImports(imports []ImportDecl) {
    sort.SliceStable(imports, func(i, j int) bool {
        pi, pj := imports[i].Path, imports[j].Path
        // Stdlib paths don't contain "/"
        iStdlib := !strings.Contains(pi, "/")
        jStdlib := !strings.Contains(pj, "/")
        if iStdlib != jStdlib {
            return iStdlib // stdlib comes first
        }
        return pi < pj
    })
}

// ─── 4. MERGE SORT FOR STABLE TOKEN SORT ─────────────────────────────────────

// Token represents a lexer token with position information
type Token struct {
    Type    string
    Value   string
    Line    int
    Column  int
}

// MergeSortTokens sorts tokens by position using merge sort.
// We use merge sort (not quicksort) because stability is important:
// tokens on the same line should retain their column order.
func MergeSortTokens(tokens []Token) []Token {
    if len(tokens) <= 1 {
        return tokens
    }
    mid := len(tokens) / 2
    left := MergeSortTokens(tokens[:mid])
    right := MergeSortTokens(tokens[mid:])
    return mergeTokens(left, right)
}

func mergeTokens(left, right []Token) []Token {
    result := make([]Token, 0, len(left)+len(right))
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        li, rj := left[i], right[j]
        // Compare by line first, then column
        if li.Line < rj.Line || (li.Line == rj.Line && li.Column <= rj.Column) {
            result = append(result, li)
            i++
        } else {
            result = append(result, rj)
            j++
        }
    }
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}

// ─── 5. DEMO: SHOW ALL SORTING IN ACTION ────────────────────────────────────

func RunCompilerSortingDemo() {
    fmt.Println("=== Astra Compiler Sorting Demo ===\n")

    // Error sorting demo
    ExampleErrorSorting()

    // Struct field sorting demo
    fmt.Println("\n--- Struct Field Sorting ---")
    fields := []StructField{
        {Name: "alive", Type: "bool"},
        {Name: "score", Type: "int64"},
        {Name: "level", Type: "int32"},
        {Name: "dead", Type: "bool"},
    }
    fmt.Println("Before (declaration order):")
    for _, f := range fields {
        fmt.Printf("  %s: %s\n", f.Name, f.Type)
    }
    sorted := SortStructFieldsByAlignment(fields)
    fmt.Println("After (optimized for alignment):")
    for _, f := range sorted {
        fmt.Printf("  [offset %2d] %s: %s\n", f.Offset, f.Name, f.Type)
    }

    // Import sorting demo
    fmt.Println("\n--- Import Sorting ---")
    imports := []ImportDecl{
        {Path: "github.com/user/utils"},
        {Path: "fmt"},
        {Path: "strings"},
        {Path: "github.com/astra/stdlib"},
    }
    SortImports(imports)
    fmt.Println("Sorted imports:")
    for _, imp := range imports {
        fmt.Printf("  import \"%s\"\n", imp.Path)
    }
}
```

Running this produces:

```
=== Astra Compiler Sorting Demo ===

Before sorting:
  main.astra:42:5: error: undefined variable 'x'
  utils.astra:10:1: error: missing return statement
  main.astra:42:1: error: type mismatch: expected int, got string
  main.astra:7:3: error: unknown type 'Foobar'

After sorting:
  main.astra:7:3: error: unknown type 'Foobar'
  main.astra:42:1: error: type mismatch: expected int, got string
  main.astra:42:5: error: undefined variable 'x'
  utils.astra:10:1: error: missing return statement

--- Struct Field Sorting ---
Before (declaration order):
  alive: bool
  score: int64
  level: int32
  dead: bool
After (optimized for alignment):
  [offset  0] score: int64
  [offset  8] level: int32
  [offset 12] alive: bool
  [offset 13] dead: bool

--- Import Sorting ---
Sorted imports:
  import "fmt"
  import "strings"
  import "github.com/astra/stdlib"
  import "github.com/user/utils"
```

---

## Exercises

1. **Implement insertion sort**: Write insertion sort in Go and trace through sorting `[7, 3, 9, 1, 5]` step by step, showing the array state after each insertion.

2. **Stability test**: Write a Go program that proves whether a sorting implementation is stable. Create a list of pairs where you can check if equal-keyed elements maintain their relative order after sorting.

3. **Merge sort for strings**: Implement merge sort that sorts `[]string` alphabetically. Handle the edge cases: empty slice, single element, two elements.

4. **Quicksort worst case**: Create a slice of 1,000 already-sorted integers and time how long quicksort takes using the last element as pivot vs. a random pivot. Measure and compare.

5. **Counting sort for line numbers**: Implement counting sort to sort a list of compiler error messages by line number. Assume line numbers are between 1 and 10,000.

6. **Compare all sorts**: Implement bubble sort, insertion sort, merge sort, and quicksort. Run all four on the same randomly generated array of 10,000 integers. Time them and produce a table showing the results. Does the real-world timing match the Big O predictions?

7. **The compiler error scenario**: A compiler finds 50 errors during type checking. The errors are discovered in the order the type checker processes nodes (not in source order). Write code to sort them for display. Use `sort.SliceStable` and explain why stability matters here.

8. **Prove the lower bound intuitively**: For n=3, draw the complete decision tree for comparison-based sorting. How many leaves does it have? How deep is the tree? Does the depth match O(n log n)?

---

## Summary Table

| Algorithm      | Best       | Average    | Worst      | Space    | Stable | When to Use                        |
|----------------|------------|------------|------------|----------|--------|------------------------------------|
| Bubble Sort    | O(n)       | O(n²)      | O(n²)      | O(1)     | Yes    | Teaching only                      |
| Selection Sort | O(n²)      | O(n²)      | O(n²)      | O(1)     | No     | Minimizing swaps                   |
| Insertion Sort | O(n)       | O(n²)      | O(n²)      | O(1)     | Yes    | Small arrays, nearly sorted data   |
| Merge Sort     | O(n log n) | O(n log n) | O(n log n) | O(n)     | Yes    | Stable sort needed, linked lists   |
| Quicksort      | O(n log n) | O(n log n) | O(n²)      | O(log n) | No     | General purpose, cache-friendly    |
| Heap Sort      | O(n log n) | O(n log n) | O(n log n) | O(1)     | No     | Guaranteed O(n log n), in-place    |
| Counting Sort  | O(n+k)     | O(n+k)     | O(n+k)     | O(k)     | Yes    | Integer keys in bounded range      |
| Timsort        | O(n)       | O(n log n) | O(n log n) | O(n)     | Yes    | Real-world data, Python/Java std   |

| Concept                        | Key Fact                                              |
|--------------------------------|-------------------------------------------------------|
| Stable sort                    | Preserves relative order of equal elements            |
| Lower bound for comparison sort | O(n log n) — provably optimal                        |
| Counting sort breaks the bound  | Uses values directly, not comparisons                |
| Go's sort.Slice                 | pdqsort: quicksort + heapsort + insertion sort        |
| Astra compiler uses sorting for | Error display, struct layout, import ordering         |

The key insight of this chapter: there is no single "best" sorting algorithm. The right choice depends on your data (is it nearly sorted? are there few unique values?), your constraints (do you need stability? is space limited?), and your use case (what are you sorting, and how often?). The Astra compiler uses different strategies in different places — and understanding those trade-offs is what separates a thoughtful engineer from someone who just picks the first algorithm they think of.
