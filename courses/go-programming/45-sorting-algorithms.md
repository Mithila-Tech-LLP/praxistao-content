# Chapter 34: Sorting Algorithms

Sorting is the most studied problem in computer science. Understanding sorting algorithms means understanding fundamental trade-offs: time vs space, best-case vs worst-case, stability, and when to use comparison-based vs distribution-based sorting. Go's standard library uses a hybrid sort — knowing why helps you choose the right algorithm for your data.

## Table of Contents

1. [Comparison Sorting Basics](#1-comparison-sorting-basics)
2. [Simple Sorts — O(n²)](#2-simple-sorts--on)
3. [Efficient Sorts — O(n log n)](#3-efficient-sorts--on-log-n)
4. [Linear Sorts — O(n)](#4-linear-sorts--on)
5. [Go's sort Package](#5-gos-sort-package)
6. [Choosing the Right Sort](#6-choosing-the-right-sort)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Comparison Sorting Basics

**Key properties:**
- **Stable**: equal elements maintain their original relative order
- **In-place**: O(1) extra space (or O(log n) for recursive call stack)
- **Adaptive**: runs faster on nearly-sorted input

**The lower bound**: any comparison-based sort must make at least Ω(n log n) comparisons. This is proven by decision tree theory — the algorithm must distinguish between n! possible orderings, requiring log₂(n!) ≈ n log n comparisons.

**Therefore**: merge sort, heap sort, and quicksort are all asymptotically optimal for comparison sorting.

```go
// Helper: swap two elements
func swap(arr []int, i, j int) {
    arr[i], arr[j] = arr[j], arr[i]
}

// Helper: verify sorted
func isSorted(arr []int) bool {
    for i := 1; i < len(arr); i++ {
        if arr[i] < arr[i-1] {
            return false
        }
    }
    return true
}
```

---

## 2. Simple Sorts — O(n²)

### Bubble Sort
```go
// BubbleSort repeatedly swaps adjacent elements that are out of order.
func BubbleSort(arr []int) {
    n := len(arr)
    for i := 0; i < n-1; i++ {
        swapped := false
        for j := 0; j < n-1-i; j++ {  // Largest elements "bubble" to end
            if arr[j] > arr[j+1] {
                arr[j], arr[j+1] = arr[j+1], arr[j]
                swapped = true
            }
        }
        if !swapped {  // Optimization: stop early if already sorted
            break
        }
    }
}
// Stable: YES | In-place: YES | Best: O(n) | Worst: O(n²)
// Use: never in production — only for learning
```

### Selection Sort
```go
// SelectionSort finds the minimum element and places it at the front.
func SelectionSort(arr []int) {
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
// Stable: NO | In-place: YES | Always: O(n²) — no adaptive behavior
// Advantage: minimum number of writes (useful if writes are expensive)
```

### Insertion Sort
```go
// InsertionSort inserts each element into its correct position in the sorted prefix.
func InsertionSort(arr []int) {
    for i := 1; i < len(arr); i++ {
        key := arr[i]
        j := i - 1
        for j >= 0 && arr[j] > key {
            arr[j+1] = arr[j]  // Shift right
            j--
        }
        arr[j+1] = key
    }
}
// Stable: YES | In-place: YES | Best: O(n) for sorted input | Worst: O(n²)
// Use: excellent for small arrays (< 20 elements) and nearly sorted arrays
// Go's sort.Slice uses insertion sort for n ≤ 12
```

### Quick Check
> 1. Which of the three simple sorts is stable?
> 2. Which simple sort has the best performance on nearly-sorted data?
> 3. Why might selection sort be preferred despite O(n²) complexity?

---

## 3. Efficient Sorts — O(n log n)

### Merge Sort
```go
// MergeSort splits array in half, recursively sorts, then merges.
func MergeSort(arr []int) []int {
    if len(arr) <= 1 {
        return arr
    }

    mid := len(arr) / 2
    left := MergeSort(arr[:mid])
    right := MergeSort(arr[mid:])
    return merge(left, right)
}

func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0

    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}

// In-place merge sort (avoids O(n) allocation per call):
func MergeSortInPlace(arr []int) {
    if len(arr) <= 1 {
        return
    }
    mid := len(arr) / 2
    MergeSortInPlace(arr[:mid])
    MergeSortInPlace(arr[mid:])
    mergeInPlace(arr, mid)
}

func mergeInPlace(arr []int, mid int) {
    tmp := make([]int, len(arr))
    copy(tmp, arr)
    i, j, k := 0, mid, 0

    for i < mid && j < len(arr) {
        if tmp[i] <= tmp[j] {
            arr[k] = tmp[i]; i++
        } else {
            arr[k] = tmp[j]; j++
        }
        k++
    }
    for i < mid { arr[k] = tmp[i]; i++; k++ }
    for j < len(arr) { arr[k] = tmp[j]; j++; k++ }
}

// Stable: YES | Space: O(n) | Best/Worst/Average: O(n log n)
// Use: sorting linked lists, external sorting, when stability is required
```

### Quick Sort
```go
// QuickSort picks a pivot, partitions, and recursively sorts each half.
func QuickSort(arr []int) {
    quickSort(arr, 0, len(arr)-1)
}

func quickSort(arr []int, low, high int) {
    if low < high {
        pivotIdx := partition(arr, low, high)
        quickSort(arr, low, pivotIdx-1)
        quickSort(arr, pivotIdx+1, high)
    }
}

// Lomuto partition scheme — pivot is last element:
func partition(arr []int, low, high int) int {
    pivot := arr[high]
    i := low - 1  // Index of smaller element

    for j := low; j < high; j++ {
        if arr[j] <= pivot {
            i++
            arr[i], arr[j] = arr[j], arr[i]
        }
    }
    arr[i+1], arr[high] = arr[high], arr[i+1]
    return i + 1
}

// Three-way partition (Dutch National Flag) — handles duplicates better:
func quickSort3Way(arr []int, low, high int) {
    if low >= high {
        return
    }
    pivot := arr[low]
    lt, gt := low, high  // arr[lt..gt] are all equal to pivot
    i := low + 1

    for i <= gt {
        if arr[i] < pivot {
            arr[lt], arr[i] = arr[i], arr[lt]
            lt++; i++
        } else if arr[i] > pivot {
            arr[i], arr[gt] = arr[gt], arr[i]
            gt--
        } else {
            i++
        }
    }
    quickSort3Way(arr, low, lt-1)
    quickSort3Way(arr, gt+1, high)
}

// Stable: NO | Space: O(log n) avg stack | Average: O(n log n) | Worst: O(n²)
// Use: default general-purpose sort (best constant factors in practice)
// Pitfall: always choosing first/last pivot → O(n²) on sorted arrays → use random pivot
```

**Randomized quicksort:**
```go
import "math/rand"

func randomPartition(arr []int, low, high int) int {
    randIdx := low + rand.Intn(high-low+1)
    arr[randIdx], arr[high] = arr[high], arr[randIdx]
    return partition(arr, low, high)
}
```

### Heap Sort
```go
// HeapSort sorts using a max-heap (see Chapter 30 for details).
func HeapSort(arr []int) {
    n := len(arr)
    // Build max-heap:
    for i := n/2 - 1; i >= 0; i-- {
        heapify(arr, n, i)
    }
    // Extract max one by one:
    for end := n - 1; end > 0; end-- {
        arr[0], arr[end] = arr[end], arr[0]
        heapify(arr, end, 0)
    }
}

func heapify(arr []int, n, i int) {
    largest := i
    l, r := 2*i+1, 2*i+2

    if l < n && arr[l] > arr[largest] { largest = l }
    if r < n && arr[r] > arr[largest] { largest = r }

    if largest != i {
        arr[i], arr[largest] = arr[largest], arr[i]
        heapify(arr, n, largest)
    }
}

// Stable: NO | In-place: YES | Always: O(n log n) | Space: O(1)
// Use: when guaranteed O(n log n) + O(1) space is required
```

### Quick Check
> 1. Why is merge sort preferred over quicksort for linked lists?
> 2. What is the worst-case input for naive quicksort?
> 3. Why does heap sort have poor cache performance?

---

## 4. Linear Sorts — O(n)

These work only on specific data types (integers, bounded values).

### Counting Sort — O(n + k) where k = range of values
```go
// CountingSort sorts non-negative integers with values in [0, maxVal].
func CountingSort(arr []int, maxVal int) []int {
    count := make([]int, maxVal+1)

    // Count occurrences:
    for _, v := range arr {
        count[v]++
    }
    // Compute prefix sums (position of each element in output):
    for i := 1; i <= maxVal; i++ {
        count[i] += count[i-1]
    }
    // Build output (iterate backward for stability):
    output := make([]int, len(arr))
    for i := len(arr) - 1; i >= 0; i-- {
        output[count[arr[i]]-1] = arr[i]
        count[arr[i]]--
    }
    return output
}
// Stable: YES | Space: O(n+k) | Time: O(n+k)
// Use: when k (value range) is small, e.g., sorting ages, grades
```

### Radix Sort — O(d * (n + k)) where d = digits, k = base
```go
// RadixSort sorts non-negative integers by processing one digit at a time.
func RadixSort(arr []int) {
    if len(arr) == 0 {
        return
    }
    max := arr[0]
    for _, v := range arr {
        if v > max { max = v }
    }
    // Sort by each digit (1s, 10s, 100s, ...):
    for exp := 1; max/exp > 0; exp *= 10 {
        countingSortByDigit(arr, exp)
    }
}

func countingSortByDigit(arr []int, exp int) {
    n := len(arr)
    output := make([]int, n)
    count := make([]int, 10)  // Digits 0-9

    for _, v := range arr {
        digit := (v / exp) % 10
        count[digit]++
    }
    for i := 1; i < 10; i++ {
        count[i] += count[i-1]
    }
    for i := n - 1; i >= 0; i-- {
        digit := (arr[i] / exp) % 10
        output[count[digit]-1] = arr[i]
        count[digit]--
    }
    copy(arr, output)
}
// Stable: YES | Time: O(d(n+k)) ≈ O(n) for fixed-width integers | Space: O(n+k)
// Use: sorting large arrays of integers or fixed-length strings (phone numbers, IPs)
```

### Bucket Sort — O(n) average for uniform distribution
```go
// BucketSort sorts floats in [0, 1) into buckets.
func BucketSort(arr []float64) []float64 {
    n := len(arr)
    buckets := make([][]float64, n)

    for _, v := range arr {
        idx := int(v * float64(n))
        if idx == n { idx = n - 1 }
        buckets[idx] = append(buckets[idx], v)
    }

    result := arr[:0]
    for _, bucket := range buckets {
        sort.Float64s(bucket)  // Sort each bucket (typically insertion sort)
        result = append(result, bucket...)
    }
    return result
}
// Average: O(n) for uniform distribution | Worst: O(n²) if all in one bucket
// Use: uniformly distributed data (e.g., sampling, hashing)
```

---

## 5. Go's sort Package

```go
import "sort"

// Sort a slice of ints:
nums := []int{5, 2, 8, 1, 9, 3}
sort.Ints(nums)  // [1, 2, 3, 5, 8, 9]

// Sort a slice of strings:
words := []string{"banana", "apple", "cherry"}
sort.Strings(words)  // [apple banana cherry]

// Sort any slice with a custom comparator:
people := []struct{ Name string; Age int }{
    {"Alice", 30}, {"Bob", 25}, {"Carol", 35},
}
sort.Slice(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})
// [{Bob 25} {Alice 30} {Carol 35}]

// Stable sort (preserves order of equal elements):
sort.SliceStable(people, func(i, j int) bool {
    return people[i].Age < people[j].Age
})

// Binary search in sorted slice:
nums = []int{1, 3, 5, 7, 9}
idx := sort.SearchInts(nums, 5)   // 2 — index of 5
idx = sort.SearchInts(nums, 4)    // 2 — index where 4 would be inserted

// Check if sorted:
fmt.Println(sort.IntsAreSorted(nums))  // true
```

**Go 1.21 slices package:**
```go
import "slices"

nums := []int{5, 2, 8, 1}
slices.Sort(nums)                           // In-place sort
slices.SortFunc(nums, func(a, b int) int { return a - b })  // Custom comparator
idx, _ := slices.BinarySearch(nums, 5)     // Returns index + found bool
```

**What algorithm does Go use?**
Go's `sort.Slice` uses **pdqsort** (Pattern-Defeating Quicksort) — a hybrid of:
- Insertion sort for small arrays (≤12 elements)
- Heap sort as fallback (prevents O(n²) worst case)
- Quicksort with median-of-3 pivot + three-way partitioning

---

## 6. Choosing the Right Sort

```
General purpose, unknown data      → sort.Slice (pdqsort) — O(n log n) worst case
Stability required                 → sort.SliceStable
Small array (< 15 elements)        → Insertion sort — low overhead
Integers in small range [0, k]     → Counting sort — O(n+k)
Fixed-width integers/strings       → Radix sort — O(n) for large datasets
Nearly sorted input                → Insertion sort or timsort
Memory constrained                 → Heap sort or quicksort (in-place)
Linked list                        → Merge sort (no random access needed)
```

**Sorting complexity summary:**
| Algorithm | Best | Average | Worst | Space | Stable |
|-----------|------|---------|-------|-------|--------|
| Bubble | O(n) | O(n²) | O(n²) | O(1) | Yes |
| Selection | O(n²) | O(n²) | O(n²) | O(1) | No |
| Insertion | O(n) | O(n²) | O(n²) | O(1) | Yes |
| Merge | O(n log n) | O(n log n) | O(n log n) | O(n) | Yes |
| Quick | O(n log n) | O(n log n) | O(n²) | O(log n) | No |
| Heap | O(n log n) | O(n log n) | O(n log n) | O(1) | No |
| Counting | O(n+k) | O(n+k) | O(n+k) | O(n+k) | Yes |
| Radix | O(dn) | O(dn) | O(dn) | O(n+k) | Yes |

---

## Summary

- **O(n²) sorts**: bubble (learning only), selection (min writes), insertion (small/nearly-sorted)
- **O(n log n) sorts**: merge (stable, linked lists), quicksort (fastest in practice), heapsort (O(1) space, guaranteed)
- **Linear sorts**: counting (small integer range), radix (fixed-width), bucket (uniform distribution)
- **Go's sort.Slice**: pdqsort — insertion sort for small n, quicksort with median pivot, heapsort fallback
- **Stability matters** when: sorting by one key then another (multi-key sort), preserving original order of equals
- **Adaptive behavior**: insertion sort is O(n) on sorted data; pdqsort detects and exploits patterns

---

## Exercises

### Easy
1. Implement all three O(n²) sorts (bubble, selection, insertion). Write a benchmark comparing them on: sorted input, reverse-sorted input, random input of size 1000.
2. Write `SortByFrequency(nums []int) []int` — sort elements by their frequency (most frequent first). Elements with the same frequency appear in any order. Use counting + sort.SliceStable.
3. Write `MergeSortedArrays(arrays [][]int) []int` — merge N sorted arrays (no duplicates) into one sorted array without using sort.Slice. Use a priority queue (min-heap).

### Medium
4. Custom sort stability test: Create a test that verifies `sort.Slice` is NOT stable and `sort.SliceStable` IS stable. Use records with equal keys but distinct identifiers. Confirm that stable sort preserves original order while unstable sort doesn't.
5. Dutch National Flag: Implement the three-way partitioning to sort an array with only three distinct values (e.g., 0, 1, 2 representing red/white/blue) in O(n) time and O(1) space. Verify it's one-pass — no second scan of the array.
6. External sort simulation: Implement an external sort that handles data larger than memory. Given 1GB of integers but only 100MB of "RAM" (simulate with a fixed buffer size): (1) split data into sorted chunks that fit in RAM, (2) merge chunks using a k-way merge with a priority queue. Test with 10M integers split into 10 chunks of 1M.

### Hard
7. Tim sort: Implement a simplified Timsort — the algorithm used by Python and Java. Key features: (1) detect natural runs (already sorted sequences), (2) if run is too short, extend with insertion sort, (3) merge runs using a merge stack (galloping mode). Benchmark against `sort.Slice` for: random data, nearly-sorted data (90% sorted), runs of ascending then descending.
8. Parallel merge sort: Implement a parallel merge sort using goroutines. Below a threshold (e.g., 10,000 elements), fall back to sequential sort. Above the threshold, sort each half in a goroutine concurrently. Use `sync.WaitGroup` to synchronize. Benchmark for 1M, 10M, 100M integers. Does parallelism help? At what size does it start to win? What is the optimal threshold?
