# Chapter 35: Searching and Binary Search

Binary search is deceptively simple and endlessly subtle. Most programmers can describe the idea in ten seconds but write a buggy implementation in ten minutes — off-by-one errors lurk in every `mid` calculation and boundary update. This chapter builds binary search from first principles and shows the full generalization: searching for any monotonic predicate.

## Table of Contents

1. [Linear Search](#1-linear-search)
2. [Binary Search — The Classic Form](#2-binary-search--the-classic-form)
3. [The Generalized Pattern](#3-the-generalized-pattern)
4. [Binary Search on Answer](#4-binary-search-on-answer)
5. [Searching in Rotated / Special Arrays](#5-searching-in-rotated--special-arrays)
6. [Exponential and Interpolation Search](#6-exponential-and-interpolation-search)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Linear Search

```go
// LinearSearch scans every element. O(n). Works on unsorted data.
func LinearSearch[T comparable](arr []T, target T) int {
    for i, v := range arr {
        if v == target {
            return i
        }
    }
    return -1
}

// When linear search is the right choice:
// - Array is unsorted and only searched once
// - Array is tiny (< ~16 elements; cache effects dominate)
// - Searching for multiple matches (collect all occurrences)
```

---

## 2. Binary Search — The Classic Form

The invariant: `target`, if it exists, is always in `[lo, hi]`.

```go
// BinarySearch returns the index of target, or -1 if not found.
// Array must be sorted in ascending order.
func BinarySearch(arr []int, target int) int {
    lo, hi := 0, len(arr)-1

    for lo <= hi {
        // NEVER use (lo + hi) / 2 — integer overflow for large indices!
        mid := lo + (hi-lo)/2

        switch {
        case arr[mid] == target:
            return mid
        case arr[mid] < target:
            lo = mid + 1  // Target is in right half
        default:
            hi = mid - 1  // Target is in left half
        }
    }
    return -1
}
```

**Common bug: off-by-one in boundary update.** Writing `lo = mid` or `hi = mid` (without ±1) causes infinite loops when `hi - lo == 1`. Always update to `mid±1`.

### LowerBound — first position where arr[i] >= target
```go
// LowerBound returns the index of the first element >= target.
// Equivalent to C++ lower_bound.
// If all elements are < target, returns len(arr).
func LowerBound(arr []int, target int) int {
    lo, hi := 0, len(arr)  // Note: hi = len, not len-1

    for lo < hi {  // Note: strict < not <=
        mid := lo + (hi-lo)/2
        if arr[mid] < target {
            lo = mid + 1
        } else {
            hi = mid  // Keep mid as a candidate
        }
    }
    return lo
}
```

### UpperBound — first position where arr[i] > target
```go
// UpperBound returns the index of the first element > target.
// Equivalent to C++ upper_bound.
func UpperBound(arr []int, target int) int {
    lo, hi := 0, len(arr)

    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] <= target {
            lo = mid + 1
        } else {
            hi = mid
        }
    }
    return lo
}

// Count occurrences using lower/upper bound:
func CountOccurrences(arr []int, target int) int {
    return UpperBound(arr, target) - LowerBound(arr, target)
}
```

**Using Go's sort package:**
```go
import "sort"

arr := []int{1, 3, 3, 3, 7, 9}
// SearchInts returns leftmost index where arr[i] >= target (lower_bound):
i := sort.SearchInts(arr, 3)      // 1 — first occurrence of 3
j := sort.SearchInts(arr, 4)      // 4 — first position after all 3s

// Generic sort.Search(n, f) — find smallest i in [0, n) where f(i) is true:
i = sort.Search(len(arr), func(k int) bool { return arr[k] >= 3 })  // 1
j = sort.Search(len(arr), func(k int) bool { return arr[k] > 3 })   // 4
```

### Quick Check
> 1. Why is `mid = lo + (hi-lo)/2` better than `mid = (lo+hi)/2`?
> 2. What does `LowerBound` return when the target doesn't exist in the array?
> 3. What's the key invariant difference between `lo <= hi` and `lo < hi` variants?

---

## 3. The Generalized Pattern

Any binary search is searching for the **transition point** in a boolean predicate that is `false` for a prefix and `true` for the rest.

```
Index:     0  1  2  3  4  5  6  7
Array:     1  3  5  7  9  11 13 15
f(i) ≥ 7? F  F  F  T  T  T  T  T
                   ^--- LowerBound(7) = 3
```

**The template:**
```go
// findFirstTrue returns the smallest index where f(mid) is true,
// searching in [lo, hi] (inclusive).
// Precondition: f is false for all indices below the answer, true for all above.
func findFirstTrue(lo, hi int, f func(int) bool) int {
    for lo < hi {
        mid := lo + (hi-lo)/2
        if f(mid) {
            hi = mid        // mid could be the answer — keep it
        } else {
            lo = mid + 1    // mid definitely not the answer — skip it
        }
    }
    // lo == hi: single candidate remaining
    if f(lo) {
        return lo
    }
    return -1  // No valid index found
}

// Example: find first element >= target
func lowerBoundTemplate(arr []int, target int) int {
    return findFirstTrue(0, len(arr), func(i int) bool {
        if i == len(arr) { return false }
        return arr[i] >= target
    })
}
```

---

## 4. Binary Search on Answer

When the answer is a number in a range and you can write a function `canAchieve(x)` that returns true iff it's possible to achieve at most x (or at least x), you can binary search on the answer.

**Example: Koko Eating Bananas**
```go
// KokoEating: Given piles of bananas, find minimum speed k (bananas/hour)
// so Koko can finish all piles within h hours.
// Constraints: n = len(piles), 1 <= h, answer in [1, max(piles)]
func MinEatingSpeed(piles []int, h int) int {
    maxPile := 0
    for _, p := range piles {
        if p > maxPile { maxPile = p }
    }

    // Can Koko finish with speed k?
    canFinish := func(k int) bool {
        hours := 0
        for _, p := range piles {
            hours += (p + k - 1) / k  // ceiling division
        }
        return hours <= h
    }

    // Binary search: find smallest k where canFinish is true
    lo, hi := 1, maxPile
    for lo < hi {
        mid := lo + (hi-lo)/2
        if canFinish(mid) {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```

**Example: Allocate Minimum Pages**
```go
// MinPages: allocate n books to k students (each student reads contiguous books),
// minimize the maximum pages any student reads.
func AllocateBooks(pages []int, k int) int {
    if k > len(pages) { return -1 }

    sum := 0
    maxPage := 0
    for _, p := range pages {
        sum += p
        if p > maxPage { maxPage = p }
    }

    // Can we allocate with max limit = limit?
    canAllocate := func(limit int) bool {
        students, current := 1, 0
        for _, p := range pages {
            if p > limit { return false }
            if current+p > limit {
                students++
                current = p
                if students > k { return false }
            } else {
                current += p
            }
        }
        return true
    }

    lo, hi := maxPage, sum
    for lo < hi {
        mid := lo + (hi-lo)/2
        if canAllocate(mid) {
            hi = mid
        } else {
            lo = mid + 1
        }
    }
    return lo
}
```

**Example: Sqrt — classic integer binary search**
```go
func MySqrt(x int) int {
    if x < 2 { return x }

    lo, hi := 1, x/2
    for lo < hi {
        mid := lo + (hi-lo+1)/2  // Ceiling division to avoid infinite loop
        if mid*mid <= x {
            lo = mid
        } else {
            hi = mid - 1
        }
    }
    return lo
}
```

---

## 5. Searching in Rotated / Special Arrays

### Rotated Sorted Array
```go
// SearchRotated: binary search in a sorted array that was rotated at some pivot.
// [4, 5, 6, 7, 0, 1, 2] — rotated at index 4
func SearchRotated(arr []int, target int) int {
    lo, hi := 0, len(arr)-1

    for lo <= hi {
        mid := lo + (hi-lo)/2
        if arr[mid] == target { return mid }

        // Left half is sorted:
        if arr[lo] <= arr[mid] {
            if arr[lo] <= target && target < arr[mid] {
                hi = mid - 1
            } else {
                lo = mid + 1
            }
        } else { // Right half is sorted:
            if arr[mid] < target && target <= arr[hi] {
                lo = mid + 1
            } else {
                hi = mid - 1
            }
        }
    }
    return -1
}
```

### Find Minimum in Rotated Array
```go
func FindMin(arr []int) int {
    lo, hi := 0, len(arr)-1

    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] > arr[hi] {
            lo = mid + 1  // Min is in the right half
        } else {
            hi = mid      // Min is at mid or to the left
        }
    }
    return arr[lo]
}
```

### Search in 2D Matrix
```go
// SearchMatrix: each row is sorted, first element of each row > last element of previous row.
// Treat it as a flattened sorted array.
func SearchMatrix(matrix [][]int, target int) bool {
    if len(matrix) == 0 { return false }
    m, n := len(matrix), len(matrix[0])
    lo, hi := 0, m*n-1

    for lo <= hi {
        mid := lo + (hi-lo)/2
        val := matrix[mid/n][mid%n]
        if val == target { return true }
        if val < target { lo = mid + 1 } else { hi = mid - 1 }
    }
    return false
}
```

### Find Peak Element
```go
// FindPeak: find any peak element where arr[peak] > neighbors.
// Uses binary search: go toward higher neighbor.
func FindPeakElement(arr []int) int {
    lo, hi := 0, len(arr)-1

    for lo < hi {
        mid := lo + (hi-lo)/2
        if arr[mid] > arr[mid+1] {
            hi = mid    // Peak is on left side
        } else {
            lo = mid+1  // Peak is on right side
        }
    }
    return lo
}
```

---

## 6. Exponential and Interpolation Search

### Exponential Search — for unbounded arrays
```go
// ExponentialSearch: double the range until we overshoot, then binary search.
// O(log n). Useful when array is infinite/unbounded.
func ExponentialSearch(arr []int, target int) int {
    if arr[0] == target { return 0 }

    i := 1
    for i < len(arr) && arr[i] <= target {
        i *= 2
    }

    lo := i / 2
    hi := min(i, len(arr)-1)

    for lo <= hi {
        mid := lo + (hi-lo)/2
        if arr[mid] == target { return mid }
        if arr[mid] < target { lo = mid + 1 } else { hi = mid - 1 }
    }
    return -1
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
```

### Interpolation Search — for uniformly distributed data
```go
// InterpolationSearch: estimate position based on value distribution.
// O(log log n) for uniform distribution, O(n) worst case.
func InterpolationSearch(arr []int, target int) int {
    lo, hi := 0, len(arr)-1

    for lo <= hi && target >= arr[lo] && target <= arr[hi] {
        if lo == hi {
            if arr[lo] == target { return lo }
            return -1
        }
        // Estimate position proportionally:
        pos := lo + (target-arr[lo])*(hi-lo)/(arr[hi]-arr[lo])

        if arr[pos] == target { return pos }
        if arr[pos] < target { lo = pos + 1 } else { hi = pos - 1 }
    }
    return -1
}
```

---

## Summary

- **Binary search invariant**: maintain the search space — update `lo = mid+1` or `hi = mid-1` to always shrink
- **Use `lo + (hi-lo)/2`** not `(lo+hi)/2` to avoid integer overflow
- **LowerBound / UpperBound** pattern handles duplicate targets and "first/last occurrence" queries
- **`sort.Search(n, f)`** is Go's generic binary search — provide a monotonic predicate
- **Binary search on answer**: when you can express "can I achieve X?" as a monotonic function, binary search on the answer space
- **Rotated arrays**: identify which half is sorted, then decide which half the target is in
- **Exponential search** finds the range; **interpolation search** exploits uniform distribution

---

## Exercises

### Easy
1. Implement `FirstAndLastOccurrence(arr []int, target int) (int, int)` using lower and upper bound. Return `(-1, -1)` if not found. Write table-driven tests with: target not present, one occurrence, multiple occurrences, target at boundaries.
2. Implement `SearchInsertPosition(arr []int, target int) int` — return the index where target is or would be inserted to keep the array sorted. This is exactly `LowerBound`.
3. Implement `IsBadVersion(version int) bool` (mock it with a closure over a threshold). Find the first bad version using binary search. This is the classic "find first true" template.

### Medium
4. **Capacity to ship packages in D days**: Given array `weights` and integer `D`, find the minimum weight capacity of a ship to deliver all packages within D days. Packages must be loaded in order. Binary search on the capacity (range: `max(weights)` to `sum(weights)`). Verify: `weights = [1,2,3,4,5,6,7,8,9,10], D = 5` → 15.
5. **Find duplicate in [1..n+1] array without modifying it**: Given array of `n+1` integers in range `[1, n]`, binary search on the value space. For each mid, count elements `<= mid`. If count > mid, the duplicate is in `[1, mid]`, otherwise in `[mid+1, n]`. Time: O(n log n), Space: O(1).
6. **Median of two sorted arrays**: Implement `FindMedianSortedArrays(nums1, nums2 []int) float64` in O(log(min(m,n))). Binary search on the partition point of the smaller array. This is LeetCode hard — study the partition invariant carefully.

### Hard
7. **Aggressive cows**: Given farm positions and `k` cows, maximize the minimum distance between any two cows. Binary search on the minimum distance, and greedily check if `k` cows can be placed. Verify: `positions = [1,2,8,4,9], k = 3` → 3.
8. **K-th smallest element in sorted matrix**: Given n×n matrix where each row and column is sorted, find the k-th smallest element in O(n log(max-min)). Binary search on the value range; for each mid, count elements ≤ mid using a two-pointer walk from bottom-left corner.
