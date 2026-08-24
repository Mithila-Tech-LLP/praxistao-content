# Chapter 29: Complexity Analysis — Big O, Time and Space

Before you can choose the right data structure or algorithm, you need a way to compare them objectively. Big O notation is that language. It describes how performance scales with input size, independent of hardware or language.

## Table of Contents

1. [Why Complexity Matters](#1-why-complexity-matters)
2. [Big O Notation](#2-big-o-notation)
3. [Common Complexity Classes](#3-common-complexity-classes)
4. [Analyzing Code](#4-analyzing-code)
5. [Space Complexity](#5-space-complexity)
6. [Amortized Analysis](#6-amortized-analysis)
7. [Go Benchmark Verification](#7-go-benchmark-verification)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Complexity Matters

```go
// Find if a slice contains a target — two approaches:
func linearSearch(nums []int, target int) bool {
    for _, n := range nums { if n == target { return true } }
    return false
}

func binarySearch(nums []int, target int) bool { // requires sorted input
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := (lo + hi) / 2
        if nums[mid] == target { return true }
        if nums[mid] < target  { lo = mid + 1 } else { hi = mid - 1 }
    }
    return false
}
```

For n = 1,000,000:
- Linear search: up to **1,000,000** comparisons
- Binary search: at most **20** comparisons (`log₂(1,000,000) ≈ 20`)

This difference doesn't show up at n = 10. It's catastrophic at n = 10,000,000. Complexity analysis predicts this before you write a benchmark.

---

## 2. Big O Notation

Big O describes the **worst-case upper bound** on how operations grow relative to input size `n`.

**Formal definition:** f(n) = O(g(n)) if there exist constants c > 0 and n₀ such that f(n) ≤ c·g(n) for all n ≥ n₀.

**In practice:** drop constants and lower-order terms.

```
f(n) = 3n² + 100n + 5000
→ O(n²)     (3 is a constant, 100n and 5000 are lower-order)

f(n) = n/2 + log(n)
→ O(n)      (n dominates log(n))
```

### Related notations

| Notation | Meaning | Usage |
|----------|---------|-------|
| O(f(n)) | Upper bound (worst case) | Most common |
| Ω(f(n)) | Lower bound (best case) | Theoretical |
| Θ(f(n)) | Tight bound (both) | Exact characterization |

In interviews and daily work, "O(n)" usually means "grows linearly in the worst case."

---

## 3. Common Complexity Classes

Ordered from fastest to slowest as n grows:

```
O(1) < O(log n) < O(n) < O(n log n) < O(n²) < O(2ⁿ) < O(n!)
```

### O(1) — Constant

The operation takes the same time regardless of input size.

```go
// Hash map lookup, array index, push/pop from stack
m := map[int]string{1: "a", 2: "b"}
v := m[42]          // O(1) average
arr := [3]int{1, 2, 3}
x := arr[1]         // O(1) always
```

### O(log n) — Logarithmic

Input is halved at each step. Grows very slowly.

```go
// Binary search, balanced BST lookup, heap insert
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {                    // loops log₂(n) times
        mid := lo + (hi-lo)/2
        if nums[mid] == target { return mid }
        if nums[mid] < target  { lo = mid + 1 } else { hi = mid - 1 }
    }
    return -1
}
// n=1,000,000 → 20 iterations
// n=1,000,000,000 → 30 iterations
```

### O(n) — Linear

One pass through the input.

```go
// Sum of a slice, linear search, single-pass scan
func sum(nums []int) int {
    total := 0
    for _, n := range nums { total += n }  // n iterations
    return total
}
```

### O(n log n) — Linearithmic

Sorting a comparison-based data structure. This is the optimal lower bound for comparison sort.

```go
// Merge sort, heap sort, Go's sort.Slice (pdqsort)
sort.Slice(data, func(i, j int) bool { return data[i] < data[j] })
// n=1,000,000 → ~20,000,000 comparisons
```

### O(n²) — Quadratic

Nested loops over the same input.

```go
// Bubble sort, checking all pairs
func allPairs(nums []int) {
    for i := 0; i < len(nums); i++ {
        for j := i + 1; j < len(nums); j++ {  // n*(n-1)/2 pairs
            fmt.Println(nums[i], nums[j])
        }
    }
}
// n=1,000 → 499,500 pairs
// n=10,000 → 49,995,000 pairs
```

### O(2ⁿ) — Exponential

Each element can be included or excluded. Brute-force subsets.

```go
// Generating all subsets: 2^n subsets for n elements
// n=10 → 1,024    n=20 → 1,048,576    n=40 → 1 trillion
func subsets(nums []int) [][]int {
    result := [][]int{{}}
    for _, n := range nums {
        current := make([][]int, len(result))
        for i, s := range result {
            newSubset := append([]int{}, s...)
            current[i] = append(newSubset, n)
        }
        result = append(result, current...)
    }
    return result
}
```

### O(n!) — Factorial

All permutations. Only feasible for n ≤ ~12.

```
n=10 → 3,628,800
n=15 → 1,307,674,368,000
```

---

## 4. Analyzing Code

### Rule 1: Count dominant operations

```go
func findMax(nums []int) int {
    max := nums[0]                    // O(1)
    for _, n := range nums[1:] {      // O(n)
        if n > max { max = n }        // O(1) per iteration
    }
    return max                        // O(1)
}
// Total: O(1) + O(n) * O(1) + O(1) = O(n)
```

### Rule 2: Add sequential blocks

```go
func twoPass(nums []int) {
    // First pass: O(n)
    for _, n := range nums { fmt.Println(n) }
    // Second pass: O(n)
    for _, n := range nums { fmt.Println(n * 2) }
    // Total: O(n) + O(n) = O(2n) = O(n)
}
```

### Rule 3: Multiply nested loops

```go
func nestedLoops(nums []int) {
    for i := range nums {           // O(n)
        for j := range nums {       // O(n)
            fmt.Println(nums[i] + nums[j])
        }
    }
    // Total: O(n) * O(n) = O(n²)
}
```

### Rule 4: Nested loops with independent ranges

```go
func matrix(rows, cols []int) {
    for _, r := range rows {        // O(m)
        for _, c := range cols {    // O(k)
            fmt.Println(r, c)
        }
    }
    // Total: O(m * k), not O(n²) unless m == k == n
}
```

### Rule 5: Recognize halving

```go
func logLoop(n int) {
    for i := n; i > 0; i /= 2 {    // i = n, n/2, n/4, ... 1
        fmt.Println(i)              // log₂(n) iterations
    }
    // O(log n)
}
```

### Tricky example: two nested loops, inner depends on outer

```go
func triangleLoop(n int) {
    for i := 0; i < n; i++ {
        for j := 0; j < i; j++ {   // 0+1+2+...+(n-1) = n*(n-1)/2
            fmt.Println(i, j)
        }
    }
    // Total: n*(n-1)/2 = O(n²)
}
```

---

## 5. Space Complexity

Space complexity counts the **extra memory** an algorithm uses, not counting the input itself (unless the problem requires it).

```go
// O(1) space — in-place
func reverseInPlace(nums []int) {
    for i, j := 0, len(nums)-1; i < j; i, j = i+1, j-1 {
        nums[i], nums[j] = nums[j], nums[i]
    }
}

// O(n) space — creates a new slice
func reverseCopy(nums []int) []int {
    result := make([]int, len(nums))
    for i, j := 0, len(nums)-1; i < len(nums); i, j = i+1, j-1 {
        result[i] = nums[j]
    }
    return result
}
```

### Recursion stack space

Each function call consumes stack space. Deep recursion = O(depth) space.

```go
// O(n) space — n recursive calls on the stack
func factorial(n int) int {
    if n <= 1 { return 1 }
    return n * factorial(n-1)  // stack depth n
}

// O(log n) space — binary search recursion
func bsRecursive(nums []int, lo, hi, target int) int {
    if lo > hi { return -1 }
    mid := lo + (hi-lo)/2
    if nums[mid] == target { return mid }
    if nums[mid] < target { return bsRecursive(nums, mid+1, hi, target) }
    return bsRecursive(nums, lo, mid-1, target)
    // maximum log₂(n) stack frames at once
}
```

---

## 6. Amortized Analysis

Some operations are occasionally expensive but cheap on average. Amortized analysis asks: what is the average cost **per operation** over a sequence of operations?

### Go slice append

```go
s := []int{}
for i := 0; i < n; i++ {
    s = append(s, i)  // "usually" O(1), "rarely" O(n) for resize
}
// Total work: n + n/2 + n/4 + ... = 2n = O(n)
// Amortized per operation: O(n)/n = O(1)
```

Each resize doubles capacity, so resizes happen at sizes 1, 2, 4, 8, ..., n. Total copy work = 1 + 2 + 4 + ... + n = 2n. Amortized O(1) per append.

### Stack with dynamic backing array

```go
type Stack[T any] struct {
    data []T
}

func (s *Stack[T]) Push(v T) {
    s.data = append(s.data, v) // O(1) amortized
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.data) == 0 { var z T; return z, false }
    n := len(s.data) - 1
    v := s.data[n]
    s.data = s.data[:n]        // O(1)
    return v, true
}
```

---

## 7. Go Benchmark Verification

You can verify complexity empirically with benchmarks. If an O(n) function truly is linear, doubling n should roughly double the time.

```go
package main_test

import (
    "testing"
)

func linearSearch(nums []int, target int) bool {
    for _, n := range nums { if n == target { return true } }
    return false
}

func BenchmarkLinear1k(b *testing.B) {
    nums := make([]int, 1_000)
    for i := range nums { nums[i] = i }
    target := 999
    b.ResetTimer()
    for i := 0; i < b.N; i++ { linearSearch(nums, target) }
}

func BenchmarkLinear10k(b *testing.B) {
    nums := make([]int, 10_000)
    for i := range nums { nums[i] = i }
    target := 9999
    b.ResetTimer()
    for i := 0; i < b.N; i++ { linearSearch(nums, target) }
}

// If linear: time(10k) ≈ 10 × time(1k)
// If quadratic: time(10k) ≈ 100 × time(1k)
```

```bash
go test -bench=. -benchmem
# BenchmarkLinear1k     2000000     742 ns/op
# BenchmarkLinear10k     200000    7423 ns/op
# Ratio ≈ 10.0 → confirms O(n)
```

---

## Summary

| Complexity | Example | n=1k | n=1M |
|------------|---------|------|------|
| O(1) | Hash map get | 1 op | 1 op |
| O(log n) | Binary search | 10 ops | 20 ops |
| O(n) | Linear scan | 1,000 | 1,000,000 |
| O(n log n) | Merge sort | 10,000 | 20,000,000 |
| O(n²) | Bubble sort | 1,000,000 | 10¹² |
| O(2ⁿ) | Subsets | Too slow | ∞ |

**Analysis rules:**
1. Drop constants: O(2n) = O(n)
2. Drop lower-order terms: O(n² + n) = O(n²)
3. Sequential blocks: O(n) + O(n) = O(n)
4. Nested loops: O(n) * O(n) = O(n²)
5. Halving = log: `i /= 2` → O(log n)
6. Recursion stack depth contributes to space complexity

---

## Exercises

### Easy
1. Classify each snippet: `for i := 0; i < n; i++ { for j := 0; j < 100; j++ { } }` — is this O(n) or O(n²)? Explain.
2. Write a function that finds all duplicates in a slice in O(n) time and O(n) space. Then write another version in O(n²) time and O(1) space. Benchmark both at n = 10,000.
3. Analyze the space complexity of a recursive Fibonacci function. Then implement an iterative version with O(1) space and compare.

### Medium
4. Implement and benchmark three versions of "find two numbers that sum to target": brute force O(n²), sort + two pointers O(n log n), hash map O(n). Plot the results — at what n does the O(n) version win?
5. Analyze `sort.Slice` on increasingly large random slices (n = 1k, 10k, 100k, 1M). If it's truly O(n log n), the ratio of `time(10n)/time(n)` should converge to 10 × log(10)/log(1) ≈ 3.3 per 10× input growth. Verify empirically.
6. Analyze this code: `func f(n int) { if n <= 1 { return } f(n/2); f(n/2); for i := 0; i < n; i++ { } }`. Time complexity is O(n log n). Prove it by writing the recurrence T(n) = 2T(n/2) + n and applying the Master Theorem.

### Hard
7. Build a **complexity verifier**: write a function `DetectComplexity(bench func(n int) time.Duration, ns []int) string` that runs `bench` at each n value, computes the ratios between consecutive timings, and classifies the function as O(1), O(log n), O(n), O(n log n), or O(n²) by comparing observed growth rates to expected ratios.
8. Analyze the amortized complexity of Go's `map`: implement a hash table from scratch with open addressing and load factor resize at 75%. Measure the amortized cost per insert over 1M insertions and verify it stays O(1). Measure the worst-case cost for a single insert that triggers a resize.
