# Chapter 49: Divide and Conquer

Divide and conquer breaks a problem into smaller sub-problems of the same type, solves them independently (usually recursively), and combines their results. The key property: sub-problems must be independent (no shared state). When sub-problems overlap, use dynamic programming instead.

## Table of Contents

1. [The Pattern](#1-the-pattern)
2. [Master Theorem](#2-master-theorem)
3. [Classic Algorithms](#3-classic-algorithms)
4. [D&C in Practice](#4-dc-in-practice)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. The Pattern

```
divide_and_conquer(problem):
    if problem is small enough:
        solve directly (base case)
    
    sub-problems = divide(problem)
    solutions    = [divide_and_conquer(p) for p in sub-problems]
    return combine(solutions)
```

The three steps:
- **Divide**: split the problem into sub-problems (often at the midpoint)
- **Conquer**: solve each sub-problem recursively
- **Combine**: merge the sub-problem solutions into the answer

---

## 2. Master Theorem

For recurrences of the form `T(n) = a·T(n/b) + f(n)`:
- `a` = number of sub-problems
- `b` = factor by which input is reduced
- `f(n)` = cost of divide + combine steps

| Case | Condition | Solution |
|------|-----------|----------|
| 1 | f(n) = O(n^(log_b(a) - ε)) | T(n) = Θ(n^log_b(a)) |
| 2 | f(n) = Θ(n^log_b(a)) | T(n) = Θ(n^log_b(a) · log n) |
| 3 | f(n) = Ω(n^(log_b(a) + ε)) | T(n) = Θ(f(n)) |

**Common examples:**
- Merge sort: T(n) = 2T(n/2) + Θ(n) → Case 2 → O(n log n)
- Binary search: T(n) = T(n/2) + Θ(1) → Case 2 → O(log n)
- Strassen's matrix multiply: T(n) = 7T(n/2) + Θ(n²) → O(n^2.807)

---

## 3. Classic Algorithms

### Merge Sort

```go
func mergeSort(arr []int) []int {
    if len(arr) <= 1 { return arr }
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])
    right := mergeSort(arr[mid:])
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i]); i++
        } else {
            result = append(result, right[j]); j++
        }
    }
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}
```

### Count Inversions (merge sort variant)

An inversion is a pair (i, j) where i < j but arr[i] > arr[j]. Counting inversions requires O(n²) naively. With merge sort, count during the merge step:

```go
func countInversions(arr []int) (sorted []int, count int) {
    if len(arr) <= 1 { return arr, 0 }
    mid := len(arr) / 2
    left, lCount := countInversions(arr[:mid])
    right, rCount := countInversions(arr[mid:])
    merged, mCount := mergeCount(left, right)
    return merged, lCount + rCount + mCount
}

func mergeCount(left, right []int) ([]int, int) {
    result := make([]int, 0, len(left)+len(right))
    count := 0
    i, j := 0, 0
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i]); i++
        } else {
            // All remaining elements in left[i:] are > right[j]
            count += len(left) - i  // each is an inversion with right[j]
            result = append(result, right[j]); j++
        }
    }
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result, count
}
```

### Binary Search

```go
func binarySearch(nums []int, target int) int {
    lo, hi := 0, len(nums)-1
    for lo <= hi {
        mid := lo + (hi-lo)/2  // avoids integer overflow
        if nums[mid] == target  { return mid }
        if nums[mid] < target   { lo = mid + 1 } else { hi = mid - 1 }
    }
    return -1
}
```

### Quick Sort

```go
func quickSort(arr []int, lo, hi int) {
    if lo >= hi { return }
    p := partition(arr, lo, hi)
    quickSort(arr, lo, p-1)
    quickSort(arr, p+1, hi)
}

// Lomuto partition scheme
func partition(arr []int, lo, hi int) int {
    pivot := arr[hi]
    i := lo - 1
    for j := lo; j < hi; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[hi] = arr[hi], arr[i+1]
    return i + 1
}
```

### Quick Select (k-th smallest element)

Quick select finds the k-th smallest element in O(n) average time without sorting:

```go
func quickSelect(arr []int, k int) int {
    return qs(arr, 0, len(arr)-1, k-1) // k is 1-indexed
}

func qs(arr []int, lo, hi, k int) int {
    if lo == hi { return arr[lo] }
    p := partition(arr, lo, hi)
    if k == p     { return arr[p] }
    if k < p      { return qs(arr, lo, p-1, k) }
    return qs(arr, p+1, hi, k)
}
```

### Maximum Subarray (Divide & Conquer version)

```go
// O(n log n) divide and conquer — Kadane's is O(n) but this illustrates the pattern
func maxSubarray(nums []int, lo, hi int) int {
    if lo == hi { return nums[lo] }
    mid := (lo + hi) / 2
    
    // Max subarray is entirely in left half, right half, or crosses midpoint
    leftMax := maxSubarray(nums, lo, mid)
    rightMax := maxSubarray(nums, mid+1, hi)
    crossMax := maxCrossing(nums, lo, mid, hi)
    
    return max3(leftMax, rightMax, crossMax)
}

func maxCrossing(nums []int, lo, mid, hi int) int {
    leftSum, sum := math.MinInt64, 0
    for i := mid; i >= lo; i-- {
        sum += nums[i]
        if sum > leftSum { leftSum = sum }
    }
    rightSum, sum = math.MinInt64, 0
    for i := mid + 1; i <= hi; i++ {
        sum += nums[i]
        if sum > rightSum { rightSum = sum }
    }
    return leftSum + rightSum
}

func max3(a, b, c int) int {
    if a > b { a = b }  // wrong — should be max not min
    // Let me write this correctly:
    m := a
    if b > m { m = b }
    if c > m { m = c }
    return m
}
```

### Closest Pair of Points

Find the two closest points among n points. Naive: O(n²). Divide and conquer: O(n log n).

```go
import "math"

type Point struct{ X, Y float64 }

func dist(a, b Point) float64 {
    dx, dy := a.X-b.X, a.Y-b.Y
    return math.Sqrt(dx*dx + dy*dy)
}

func closestPair(points []Point) float64 {
    if len(points) <= 3 {
        return bruteForce(points)
    }
    
    // Sort by X (done once before recursion in practice)
    mid := len(points) / 2
    midX := points[mid].X
    
    dLeft := closestPair(points[:mid])
    dRight := closestPair(points[mid:])
    d := math.Min(dLeft, dRight)
    
    // Check strip: points within distance d of the midline
    var strip []Point
    for _, p := range points {
        if math.Abs(p.X-midX) < d {
            strip = append(strip, p)
        }
    }
    
    // Sort strip by Y
    sort.Slice(strip, func(i, j int) bool {
        return strip[i].Y < strip[j].Y
    })
    
    // Check strip — only need to check next 7 points per point
    for i, p := range strip {
        for j := i + 1; j < len(strip) && strip[j].Y-p.Y < d; j++ {
            if v := dist(p, strip[j]); v < d {
                d = v
            }
        }
    }
    return d
}

func bruteForce(points []Point) float64 {
    d := math.MaxFloat64
    for i := range points {
        for j := i + 1; j < len(points); j++ {
            if v := dist(points[i], points[j]); v < d {
                d = v
            }
        }
    }
    return d
}
```

---

## 4. D&C in Practice

### Parallel execution

D&C is naturally parallelizable — sub-problems are independent:

```go
func parallelMergeSort(arr []int) []int {
    if len(arr) <= 1000 { // threshold for sequential
        return mergeSort(arr)
    }
    mid := len(arr) / 2
    
    var left, right []int
    var wg sync.WaitGroup
    wg.Add(2)
    
    go func() {
        defer wg.Done()
        left = parallelMergeSort(arr[:mid])
    }()
    go func() {
        defer wg.Done()
        right = parallelMergeSort(arr[mid:])
    }()
    
    wg.Wait()
    return merge(left, right)
}
```

### D&C vs Dynamic Programming

| Aspect | D&C | DP |
|--------|-----|-----|
| Sub-problems | Independent | Overlapping |
| Re-use results | No (or memoize) | Yes (core idea) |
| Examples | Merge sort, binary search, quick sort | Fibonacci, LCS, knapsack |
| Direction | Top-down | Bottom-up (usually) |

---

## Summary

| Algorithm | T(n) recurrence | Complexity |
|-----------|----------------|------------|
| Binary search | T(n) = T(n/2) + O(1) | O(log n) |
| Merge sort | T(n) = 2T(n/2) + O(n) | O(n log n) |
| Quick sort avg | T(n) = 2T(n/2) + O(n) | O(n log n) |
| Quick select avg | T(n) = T(n/2) + O(n) | O(n) |
| Closest pair | T(n) = 2T(n/2) + O(n log n) | O(n log² n) |
| Count inversions | T(n) = 2T(n/2) + O(n) | O(n log n) |

- Use D&C when sub-problems are independent; use DP when they overlap
- Master Theorem gives complexity from the recurrence without working through the full math
- D&C naturally parallelizes — use goroutines for sub-problems above a size threshold

---

## Exercises

### Easy
1. Implement `power(base, exp float64) float64` using divide and conquer: `x^n = (x^(n/2))^2` for even n, `x * x^(n-1)` for odd n. This runs in O(log n) vs O(n) for the naive multiply loop.
2. Find the maximum element in an array using divide and conquer: split in half, recurse, return the max of both halves. This is O(n) but uses O(log n) stack space. Compare to the iterative O(1)-space version.
3. Use the Master Theorem to determine the time complexity of `T(n) = 4T(n/2) + O(n)`, `T(n) = 4T(n/2) + O(n²)`, and `T(n) = 4T(n/2) + O(n³)`.

### Medium
4. Implement **Karatsuba multiplication** for large integers: multiply two n-digit numbers in O(n^1.585) instead of O(n²). Split each number in half, perform 3 recursive multiplications instead of 4, combine.
5. Implement **majority element** (element appearing > n/2 times) using divide and conquer: if the majority element exists in both halves, it must be the majority element overall. O(n log n) — also solvable in O(n) with Boyer-Moore voting.
6. Find the **k-th largest** element using quick select with the "median of medians" pivot selection strategy, which guarantees O(n) worst case (not just average). This is more complex than random pivot but provides a worst-case guarantee.

### Hard
7. Implement **fast Fourier transform (FFT)** using divide and conquer to multiply two polynomials in O(n log n) instead of O(n²). This is the core algorithm behind fast integer multiplication and many signal processing applications.
8. Implement a **parallel merge sort** that uses a worker pool of goroutines. The parallelism should scale with the number of available CPU cores. Measure speedup on 10M elements with 1, 2, 4, and 8 goroutines.
