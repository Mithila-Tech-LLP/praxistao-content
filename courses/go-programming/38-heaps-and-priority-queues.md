# Chapter 38: Heaps and Priority Queues

A heap is a **complete binary tree** stored in an array where the parent is always smaller (min-heap) or larger (max-heap) than its children. This structure gives O(1) access to the minimum/maximum and O(log n) insert and remove. Heaps power priority queues, heap sort, Dijkstra's shortest path, and scheduling algorithms.

## Table of Contents

1. [Heap Structure and the Array Trick](#1-heap-structure-and-the-array-trick)
2. [Min-Heap Implementation](#2-min-heap-implementation)
3. [Heap Operations](#3-heap-operations)
4. [Heap Sort](#4-heap-sort)
5. [Go's container/heap Package](#5-gos-containerheap-package)
6. [Classic Heap Problems](#6-classic-heap-problems)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Heap Structure and the Array Trick

A heap is a **complete binary tree** where every parent satisfies the heap property:
- **Min-heap**: parent ≤ both children (root = minimum)
- **Max-heap**: parent ≥ both children (root = maximum)

```
Min-heap:
        1           ← root = minimum
       / \
      3    2
     / \  / \
    7   4 5   6
```

**The array representation** stores the tree level by level:
```
index: 0  1  2  3  4  5  6
value: 1  3  2  7  4  5  6

Parent-child relationships (0-indexed):
  parent(i)      = (i - 1) / 2
  left_child(i)  = 2*i + 1
  right_child(i) = 2*i + 2
```

This is brilliant — no pointers needed! The tree structure is implicit in the array indices.

```go
func parent(i int) int      { return (i - 1) / 2 }
func leftChild(i int) int   { return 2*i + 1 }
func rightChild(i int) int  { return 2*i + 2 }
```

### Quick Check
> 1. Where is the minimum element in a min-heap?
> 2. What is the parent index of element at index 5?
> 3. Why doesn't a heap need pointers like a BST?

---

## 2. Min-Heap Implementation

```go
type MinHeap[T any] struct {
    data []T
    less func(a, b T) bool
}

func NewMinHeap[T any](less func(T, T) bool) *MinHeap[T] {
    return &MinHeap[T]{less: less}
}

func (h *MinHeap[T]) Len() int  { return len(h.data) }
func (h *MinHeap[T]) IsEmpty() bool { return len(h.data) == 0 }

func (h *MinHeap[T]) swap(i, j int) {
    h.data[i], h.data[j] = h.data[j], h.data[i]
}

// Peek returns the minimum without removing it — O(1):
func (h *MinHeap[T]) Peek() (T, bool) {
    if h.IsEmpty() {
        var zero T
        return zero, false
    }
    return h.data[0], true
}
```

---

## 3. Heap Operations

### Push — O(log n): add to end, bubble up
```go
// Push adds a new element and restores the heap property.
func (h *MinHeap[T]) Push(val T) {
    h.data = append(h.data, val)  // Add to end (complete tree)
    h.siftUp(len(h.data) - 1)    // Restore heap property
}

// siftUp moves element at index i up until heap property is satisfied.
func (h *MinHeap[T]) siftUp(i int) {
    for i > 0 {
        p := parent(i)
        if h.less(h.data[i], h.data[p]) {  // Child < parent: swap
            h.swap(i, p)
            i = p
        } else {
            break  // Heap property satisfied
        }
    }
}

// Visual: Push(1) to [2, 3, 5, 8, 4, 7]:
// Add 1 to end: [2, 3, 5, 8, 4, 7, 1]
// siftUp(6): parent(6)=2, 1<5 → swap: [2, 3, 1, 8, 4, 7, 5]
// siftUp(2): parent(2)=0, 1<2 → swap: [1, 3, 2, 8, 4, 7, 5]
// siftUp(0): i=0, done. Heap: [1, 3, 2, 8, 4, 7, 5]
```

### Pop — O(log n): swap root with last, remove last, sift down
```go
// Pop removes and returns the minimum element.
func (h *MinHeap[T]) Pop() (T, bool) {
    if h.IsEmpty() {
        var zero T
        return zero, false
    }
    min := h.data[0]
    last := len(h.data) - 1

    h.swap(0, last)              // Swap root with last element
    h.data = h.data[:last]       // Remove last (the old root)
    if len(h.data) > 0 {
        h.siftDown(0)            // Restore heap property
    }
    return min, true
}

// siftDown moves element at index i down until heap property is satisfied.
func (h *MinHeap[T]) siftDown(i int) {
    n := len(h.data)
    for {
        smallest := i
        l, r := leftChild(i), rightChild(i)

        if l < n && h.less(h.data[l], h.data[smallest]) {
            smallest = l
        }
        if r < n && h.less(h.data[r], h.data[smallest]) {
            smallest = r
        }

        if smallest == i {
            break  // Already in correct position
        }

        h.swap(i, smallest)
        i = smallest
    }
}
```

### Heapify — O(n): build heap from arbitrary array
```go
// Heapify builds a valid heap from an unordered slice in O(n).
// (Much faster than inserting n elements one by one: O(n log n))
func Heapify[T any](data []T, less func(T, T) bool) *MinHeap[T] {
    h := &MinHeap[T]{data: data, less: less}
    // Only internal nodes need sifting — start from last non-leaf
    for i := len(data)/2 - 1; i >= 0; i-- {
        h.siftDown(i)
    }
    return h
}

// Why O(n)? Intuition: nodes at deeper levels need less sifting.
// Half the nodes are leaves (0 work), quarter need 1 level of sifting, etc.
// Mathematical series: n/2*0 + n/4*1 + n/8*2 + ... = O(n)
```

### Quick Check
> 1. After pushing to a heap, which direction does the new element move?
> 2. After popping the root, why do we swap with the last element?
> 3. Why is Heapify O(n) instead of O(n log n)?

---

## 4. Heap Sort

Heap sort uses the heap to sort in O(n log n) time with O(1) extra space:

```go
// HeapSort sorts a slice in ascending order using a max-heap in-place.
func HeapSort(data []int) {
    n := len(data)
    less := func(a, b int) bool { return a > b }  // Max-heap: greater element on top

    // Phase 1: Build max-heap in O(n):
    for i := n/2 - 1; i >= 0; i-- {
        siftDownInPlace(data, i, n, less)
    }
    // data is now a max-heap: data[0] = maximum

    // Phase 2: Extract max one by one and place at end:
    for end := n - 1; end > 0; end-- {
        data[0], data[end] = data[end], data[0]  // Move max to end
        siftDownInPlace(data, 0, end, less)       // Restore heap for data[0:end]
    }
    // data is now sorted ascending
}

func siftDownInPlace(data []int, i, n int, less func(int, int) bool) {
    for {
        largest := i
        l, r := 2*i+1, 2*i+2
        if l < n && less(data[l], data[largest]) {
            largest = l
        }
        if r < n && less(data[r], data[largest]) {
            largest = r
        }
        if largest == i {
            break
        }
        data[i], data[largest] = data[largest], data[i]
        i = largest
    }
}
```

**Heap sort is not used much in practice** because quicksort's cache performance is better (heap sort jumps around memory), but it guarantees O(n log n) worst case with O(1) space.

### Quick Check
> 1. What heap type (min or max) do you use to sort in ascending order?
> 2. What is the space complexity of heap sort?
> 3. Why is heap sort rarely used in practice despite O(n log n) guarantee?

---

## 5. Go's container/heap Package

Go's standard library provides `container/heap` — it requires you to implement an interface:

```go
import "container/heap"

// You implement heap.Interface:
type Interface interface {
    sort.Interface      // Len, Less, Swap
    Push(x any)         // Add element
    Pop() any           // Remove and return last element (the implementation moves min to last)
}
```

**Task priority queue (common pattern):**
```go
type Task struct {
    Name     string
    Priority int  // Lower = higher priority
    Index    int  // Heap index (for efficient update)
}

type TaskQueue []*Task

func (pq TaskQueue) Len() int            { return len(pq) }
func (pq TaskQueue) Less(i, j int) bool  { return pq[i].Priority < pq[j].Priority }
func (pq TaskQueue) Swap(i, j int) {
    pq[i], pq[j] = pq[j], pq[i]
    pq[i].Index = i
    pq[j].Index = j
}

func (pq *TaskQueue) Push(x any) {
    n := len(*pq)
    task := x.(*Task)
    task.Index = n
    *pq = append(*pq, task)
}

func (pq *TaskQueue) Pop() any {
    old := *pq
    n := len(old)
    task := old[n-1]
    old[n-1] = nil   // Help GC
    task.Index = -1  // Mark as removed
    *pq = old[:n-1]
    return task
}

// Update a task's priority (requires knowing its heap index):
func (pq *TaskQueue) Update(task *Task, priority int) {
    task.Priority = priority
    heap.Fix(pq, task.Index)  // O(log n) rebalance
}

// Usage:
func main() {
    pq := &TaskQueue{}
    heap.Init(pq)

    heap.Push(pq, &Task{Name: "Low priority", Priority: 5})
    heap.Push(pq, &Task{Name: "Critical", Priority: 1})
    heap.Push(pq, &Task{Name: "Normal", Priority: 3})

    for pq.Len() > 0 {
        task := heap.Pop(pq).(*Task)
        fmt.Printf("Running: %s (priority %d)\n", task.Name, task.Priority)
    }
    // Output:
    // Running: Critical (priority 1)
    // Running: Normal (priority 3)
    // Running: Low priority (priority 5)
}
```

### Quick Check
> 1. How many methods does `heap.Interface` require?
> 2. What does `heap.Fix` do?
> 3. Why does the `Pop` method remove the LAST element, not the first?

---

## 6. Classic Heap Problems

### K Largest Elements — O(n log k)
```go
// KLargest returns the k largest elements from a stream.
// Uses a min-heap of size k — the smallest in the heap is the k-th largest.
func KLargest(nums []int, k int) []int {
    pq := &IntMinHeap{}
    heap.Init(pq)

    for _, n := range nums {
        heap.Push(pq, n)
        if pq.Len() > k {
            heap.Pop(pq)  // Remove smallest — keeps k largest
        }
    }

    result := make([]int, pq.Len())
    for i := pq.Len() - 1; i >= 0; i-- {
        result[i] = heap.Pop(pq).(int)
    }
    return result
}
// Min-heap of size k: heap contains the k largest seen so far.
// If new element > minimum in heap, it belongs in top k.
```

### Merge K Sorted Lists — O(n log k)
```go
type HeapItem struct {
    Val      int
    ListIdx  int  // Which list this came from
    NodePtr  *Node[int]
}

type ListHeap []HeapItem

func (h ListHeap) Len() int            { return len(h) }
func (h ListHeap) Less(i, j int) bool  { return h[i].Val < h[j].Val }
func (h ListHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *ListHeap) Push(x any)         { *h = append(*h, x.(HeapItem)) }
func (h *ListHeap) Pop() any {
    old := *h; n := len(old); x := old[n-1]; *h = old[:n-1]; return x
}

func MergeKSorted(lists []*Node[int]) *Node[int] {
    h := &ListHeap{}
    heap.Init(h)

    // Initialize with the head of each list:
    for i, list := range lists {
        if list != nil {
            heap.Push(h, HeapItem{list.Value, i, list.Next})
        }
    }

    dummy := &Node[int]{}
    cur := dummy

    for h.Len() > 0 {
        item := heap.Pop(h).(HeapItem)
        cur.Next = &Node[int]{Value: item.Val}
        cur = cur.Next

        if item.NodePtr != nil {
            heap.Push(h, HeapItem{item.NodePtr.Value, item.ListIdx, item.NodePtr.Next})
        }
    }

    return dummy.Next
}
```

### Kth Smallest in Matrix — O(k log min(k,n))
```go
// KthSmallest finds the kth smallest in a sorted n×n matrix.
// Each row and column is sorted in ascending order.
func KthSmallestMatrix(matrix [][]int, k int) int {
    n := len(matrix)
    type Cell struct{ val, row, col int }
    h := &struct {
        data []Cell
    }{}

    // Custom heap would go here — simplified for clarity:
    // Use a standard min-heap on Cell values
    // Initialize with first column elements, then advance cell by cell
    // Full implementation follows the same pattern as MergeKSorted
    _ = h
    _ = n
    return matrix[0][0]  // Placeholder — full impl in exercises
}
```

### Find Median from Data Stream — O(log n) insert, O(1) median
```go
// MedianFinder maintains a running median using two heaps:
// maxHeap (left half) | minHeap (right half)
// Invariant: maxHeap.Len() == minHeap.Len() ± 1
type MedianFinder struct {
    maxHeap *MaxIntHeap  // Lower half — top is the largest of lower half
    minHeap *IntMinHeap  // Upper half — top is the smallest of upper half
}

// IntMinHeap is the min-heap from Chapter 32. MaxIntHeap is the same
// with the comparison flipped: Less(i, j) returns h[i] > h[j].

func (mf *MedianFinder) AddNum(num int) {
    // Add to max-heap first:
    heap.Push(mf.maxHeap, num)

    // Balance: push maxHeap's top to minHeap:
    if mf.maxHeap.Len() > 0 {
        top := heap.Pop(mf.maxHeap).(int)
        heap.Push(mf.minHeap, top)
    }

    // Ensure maxHeap has same size or one more than minHeap:
    if mf.minHeap.Len() > mf.maxHeap.Len() {
        top := heap.Pop(mf.minHeap).(int)
        heap.Push(mf.maxHeap, top)
    }
}

func (mf *MedianFinder) FindMedian() float64 {
    if mf.maxHeap.Len() > mf.minHeap.Len() {
        return float64((*mf.maxHeap)[0])
    }
    return float64((*mf.maxHeap)[0]+(*mf.minHeap)[0]) / 2.0
}
```

---

## Summary

- **Heap**: complete binary tree; min-heap = root is minimum, max-heap = root is maximum
- **Array storage**: parent at `(i-1)/2`, children at `2i+1` and `2i+2`; no pointers needed
- **Push**: append to end + siftUp — O(log n)
- **Pop**: swap root with last + shrink + siftDown — O(log n)
- **Peek**: read `data[0]` — O(1)
- **Heapify**: O(n) — faster than N insertions; sift down from last internal node to root
- **Heap sort**: max-heap → extract max to end repeatedly — O(n log n), O(1) space
- **`container/heap`**: implement `Len`, `Less`, `Swap`, `Push`, `Pop`; then use `heap.Init/Push/Pop/Fix`
- **K largest**: min-heap of size k — O(n log k)
- **Merge K sorted**: min-heap seeded with list heads — O(n log k)
- **Median stream**: two heaps (max for lower half, min for upper half)

---

## Exercises

### Easy
1. Implement a max-heap by flipping the `Less` function. Verify: inserting `[5,1,9,3,7]` and popping returns `[9,7,5,3,1]`.
2. Write `KSmallest(nums []int, k int) []int` using a max-heap of size k. For each element, if it's smaller than the heap max, replace the max. Return the k smallest.
3. Write `SortColors(colors []string, order map[string]int)` using a priority queue that sorts arbitrary strings by a given priority map. Test with: `{red→1, green→2, blue→3}`.

### Medium
4. Running median: Implement `MedianFinder` completely with two heaps (max-heap for lower half, min-heap for upper half). `AddNum(int)` and `FindMedian() float64`. Test: add `[5,3,8,1,9,2]` one by one and verify median after each addition.
5. Sliding window median: Given `nums []int` and window size `k`, return the median of each window of size k. This is harder than the running median because you must efficiently REMOVE elements (not just add). Use two heaps with a "lazy deletion" map to defer removal.
6. Kth smallest in matrix: Implement `KthSmallest(matrix [][]int, k int) int` properly. The matrix has sorted rows and columns. Use a min-heap initialized with the first row, advancing column-by-column. Verify with k=1 (always returns matrix[0][0]) and k=n*n (returns the max).

### Hard
7. Task scheduler with deadlines: Implement a scheduler where tasks have a `Deadline` and `Duration`. The scheduler must pick tasks that maximize the number of tasks completed by their deadline. Use a greedy algorithm with a priority queue: sort by deadline, use a max-heap of durations to make room when behind schedule. Return the maximum number of completable tasks and the schedule.
8. Huffman encoding: Build a Huffman tree for text compression. `BuildHuffmanTree(text string) *HuffNode` counts character frequencies, then uses a min-heap to repeatedly merge the two least-frequent nodes into a new parent. `Encode(tree *HuffNode, text string) string` returns the binary encoding. `Decode(tree *HuffNode, encoded string) string` decodes it. Verify: encoding + decoding the original text produces the original text, and common characters get shorter codes.
