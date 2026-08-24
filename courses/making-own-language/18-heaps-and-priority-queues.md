# Chapter 18: Heaps and Priority Queues

> "Not all things that wait are equal. Some things are more urgent than others." — The reason priority queues exist

---

## Overview

In Chapter 15, we teased the **priority queue**: a queue where instead of serving the oldest item (FIFO), you serve the highest-priority item first. But we needed the right underlying data structure to make this efficient. That data structure is the **heap**.

A heap is a complete binary tree with one special rule: every parent is "better than" (higher or lower priority than) all of its children. This single constraint, maintained through elegant insertion and extraction algorithms, gives us O(log n) priority queue operations — far better than the O(n) you'd get from a sorted array.

Heaps are used in:
- **Dijkstra's shortest path algorithm** (essential in mapping software)
- **A* search** (pathfinding in games and robotics)
- **Huffman coding** (the compression algorithm inside ZIP, PNG, JPEG)
- **Operating system schedulers** (running the most important process first)
- **Heap sort** (an in-place O(n log n) sorting algorithm)
- **Event-driven simulation** (processing events in order of their scheduled time)

This chapter covers:
- What a heap is: the complete binary tree invariant
- Max-heap vs min-heap
- Representing a heap as an array (the beautiful index math)
- Heapify-up: efficient insert
- Heapify-down: efficient extraction
- Building a heap from an unsorted array in O(n) — better than O(n log n)!
- Heap sort
- Go's `container/heap` interface
- A generic heap in Go
- D-ary heaps
- **Astra Build Milestone**: Compiler optimization pass ordering with a priority queue

---

## What We're Building

By the end of this chapter you will have implemented a generic heap in Go and understood how the Astra compiler uses a priority queue to order its optimization passes, ensuring that cheap, high-impact optimizations run before expensive ones.

---

## Table of Contents

1. What Is a Heap?
2. Max-Heap vs Min-Heap
3. Representing a Heap as an Array
4. Heapify-Up: Insertion
5. Heapify-Down: Extraction
6. Building a Heap in O(n)
7. Heap Sort
8. Applications of Heaps
9. Go's container/heap Interface
10. Generic Heap[T] in Go
11. D-ary Heaps
12. Astra Build Milestone: Optimization Pass Ordering

---

## 1. What Is a Heap?

A **heap** is a **complete binary tree** that satisfies the **heap property**.

A **complete binary tree** is filled level by level, with the last level filled from left to right:

```
Complete (valid):           Not complete (invalid):
      1                          1
    /   \                       / \
   3     2                     3   2
  / \   /                     / \   \
 6   5 4                     6   5   4

All levels full except       Last level has a gap
the last, filled L→R          (node added on right before left)
```

The **heap property** (for a max-heap):
- Every parent node is **greater than or equal to** both of its children
- This applies to every node in the tree (not just the root)

```
Max-heap example:
          9
        /   \
       7     8
      / \   / \
     5   6 3   4
    / \
   1   2

Check: 9 >= {7, 8} ✓, 7 >= {5, 6} ✓, 8 >= {3, 4} ✓, etc.

The ROOT always contains the MAXIMUM value.
```

The heap property does **not** say anything about the relationship between siblings or nodes at the same level. A node on the left may be smaller than a node on the right — that's fine. The only constraint is parent ≥ children.

---

## 2. Max-Heap vs Min-Heap

A **max-heap** places the maximum element at the root. Extraction gives you the maximum.

A **min-heap** places the minimum element at the root. Extraction gives you the minimum.

```
Max-heap (root = 9):    Min-heap (root = 1):
        9                       1
      /   \                   /   \
     7     8                 3     2
    / \   / \               / \   / \
   5   6 3   4             5   4 7   6

extract max → 9           extract min → 1
next max → 8              next min → 2
```

To convert between them: negate all values (a max-heap of negatives behaves like a min-heap of positives). This trick lets you use one implementation for both.

In practice, **min-heaps are more commonly used** in algorithms (Dijkstra, Huffman, task scheduling) because you usually want the "smallest cost" or "highest priority (low number)" first.

---

## 3. Representing a Heap as an Array

Here is the brilliant insight that makes heaps so efficient: a complete binary tree can be stored **perfectly** in a flat array, with no pointers needed.

For a node at index `i` (0-indexed):
- **Left child** is at index `2*i + 1`
- **Right child** is at index `2*i + 2`
- **Parent** is at index `(i-1) / 2`

```
Tree:                     Array:
        9                 [9, 7, 8, 5, 6, 3, 4, 1, 2]
      /   \                0  1  2  3  4  5  6  7  8
     7     8
    / \   / \
   5   6 3   4
  / \
 1   2

Node at index 3 (value 5):
  Left child: 2*3+1 = 7 → value 1 ✓
  Right child: 2*3+2 = 8 → value 2 ✓
  Parent: (3-1)/2 = 1 → value 7 ✓
```

This array representation has extraordinary cache performance: all elements are contiguous in memory, and parent-child relationships are pure arithmetic (no pointer dereferences). This is why heaps are faster in practice than pointer-based trees even though both have O(log n) operations.

```go
// Index arithmetic for a 0-indexed heap array.
func parent(i int) int      { return (i - 1) / 2 }
func leftChild(i int) int   { return 2*i + 1 }
func rightChild(i int) int  { return 2*i + 2 }
func hasParent(i int) bool  { return i > 0 }
```

---

## 4. Heapify-Up: Insertion

To insert a new element into the heap:
1. Append it at the end of the array (maintaining the complete binary tree property)
2. "Bubble up" (or "sift up") — repeatedly swap with the parent if the heap property is violated, until it's satisfied

```
Insert 10 into max-heap [9, 7, 8, 5, 6, 3, 4]:

Step 1: Append 10 at the end.
Array: [9, 7, 8, 5, 6, 3, 4, 10]
Tree:
        9
      /   \
     7     8
    / \   / \
   5   6 3   4
  /
 10

10 > parent(5 at index 1)? Yes. Swap.
Array: [9, 10, 8, 5, 7, 3, 4, 6]  ← wait, let me redo with correct index

Append 10 at index 7 (parent = (7-1)/2 = 3, value = 5):
10 > 5? Yes → swap indices 7 and 3.

Array: [9, 7, 8, 10, 6, 3, 4, 5]
10 is now at index 3. Parent = (3-1)/2 = 1, value = 7.
10 > 7? Yes → swap indices 3 and 1.

Array: [9, 10, 8, 7, 6, 3, 4, 5]
10 is now at index 1. Parent = (1-1)/2 = 0, value = 9.
10 > 9? Yes → swap indices 1 and 0.

Array: [10, 9, 8, 7, 6, 3, 4, 5]
10 is now at index 0 (the root). No parent. Done!

Final tree:
        10
      /   \
     9     8
    / \   / \
   7   6 3   4
  /
 5
```

Heapify-up is O(log n) because the tree height is O(log n) and we do at most one swap per level.

---

## 5. Heapify-Down: Extraction

To extract the maximum (the root):
1. Swap the root with the last element
2. Remove the last element (the old root is now extracted)
3. "Bubble down" (or "sift down") the new root — repeatedly swap with the larger child until the heap property is restored

```
Extract max from [10, 9, 8, 7, 6, 3, 4, 5]:

Step 1: Swap root (10) with last element (5).
Array: [5, 9, 8, 7, 6, 3, 4, 10]

Step 2: Remove last element. Extracted value = 10.
Array: [5, 9, 8, 7, 6, 3, 4]

Step 3: Heapify-down starting at index 0.
5 at index 0. Children: left=9 (idx 1), right=8 (idx 2).
Larger child = 9. 5 < 9 → swap.

Array: [9, 5, 8, 7, 6, 3, 4]
5 at index 1. Children: left=7 (idx 3), right=6 (idx 4).
Larger child = 7. 5 < 7 → swap.

Array: [9, 7, 8, 5, 6, 3, 4]
5 at index 3. Children: left=? (idx 7, out of bounds), right=? (out of bounds).
No children. Done!

Result: extracted 10, heap is [9, 7, 8, 5, 6, 3, 4].
```

Heapify-down is also O(log n) — at most one swap per level of the tree.

---

## 6. Building a Heap in O(n)

Given an unsorted array, you can build a heap in **O(n)** — better than the O(n log n) you might expect from inserting elements one by one.

The trick: start from the last internal node (parent of the last leaf) and heapify-down each node, going from right to left:

```go
// BuildHeap converts an arbitrary slice into a valid max-heap in O(n).
func BuildHeap[T any](items []T, greater func(a, b T) bool) {
    n := len(items)
    // Start from the last non-leaf node: parent((n-1)) = (n-2)/2
    for i := (n - 2) / 2; i >= 0; i-- {
        heapifyDown(items, i, n, greater)
    }
}

func heapifyDown[T any](items []T, i, n int, greater func(a, b T) bool) {
    for {
        largest := i
        l, r := leftChild(i), rightChild(i)
        if l < n && greater(items[l], items[largest]) {
            largest = l
        }
        if r < n && greater(items[r], items[largest]) {
            largest = r
        }
        if largest == i {
            break  // heap property satisfied
        }
        items[i], items[largest] = items[largest], items[i]
        i = largest
    }
}
```

Why O(n)? Although heapify-down is O(log n), most nodes are near the leaves where the subtree height is small. A careful mathematical analysis (summing over all levels) shows the total work is O(n), not O(n log n).

---

## 7. Heap Sort

Heap sort is a brilliant in-place O(n log n) sorting algorithm:

1. Build a max-heap from the array: O(n)
2. Repeatedly extract the max (root) and place it at the end of the array: O(n log n)

```go
// HeapSort sorts a slice in ascending order using heap sort.
// Time: O(n log n), Space: O(1) in-place.
func HeapSort[T any](items []T, less func(a, b T) bool) {
    n := len(items)
    greater := func(a, b T) bool { return less(b, a) }

    // Phase 1: Build max-heap
    BuildHeap(items, greater)

    // Phase 2: Extract elements one by one
    for size := n - 1; size > 0; size-- {
        // The max is at index 0. Move it to the end.
        items[0], items[size] = items[size], items[0]
        // Restore heap property for the reduced heap (size elements).
        heapifyDown(items, 0, size, greater)
    }
    // Result: items is now sorted in ascending order.
}
```

Trace on `[4, 10, 3, 5, 1]`:

```
Initial:     [4, 10, 3, 5, 1]
Build heap:  [10, 5, 3, 4, 1]   (max-heap)

Extract 10: swap with last → [1, 5, 3, 4, | 10]
Heapify:    [5, 4, 3, 1, | 10]

Extract 5:  swap with last → [1, 4, 3, | 5, 10]
Heapify:    [4, 1, 3, | 5, 10]

Extract 4:  swap → [3, 1, | 4, 5, 10]
Heapify:    [3, 1, | 4, 5, 10]

Extract 3:  swap → [1, | 3, 4, 5, 10]
Heapify:    [1, | 3, 4, 5, 10]

Result: [1, 3, 4, 5, 10] ✓ sorted!
```

Heap sort has the rare combination of O(n log n) worst-case time and O(1) space, making it theoretically superior to quicksort (which is O(n²) worst case). However, quicksort has better cache behavior in practice, so it's typically faster on real hardware despite the worse theoretical guarantee.

---

## 8. Applications of Heaps

### Dijkstra's Shortest Path Algorithm

Find the shortest path from a source node to all other nodes in a weighted graph. Uses a min-heap to always process the nearest unvisited node:

```
Min-heap: [(distance, node)]
Start: push (0, source)

Loop:
  1. Pop the node with smallest distance → (d, u)
  2. For each neighbor v of u:
     if d + weight(u,v) < dist[v]:
        dist[v] = d + weight(u,v)
        push (dist[v], v) into heap
```

Without a heap: O(V²). With a heap: O((V + E) log V). For sparse graphs, this is dramatically faster.

### Huffman Coding (Data Compression)

Huffman coding assigns shorter bit strings to more frequent characters. Uses a min-heap:

```
Characters: a(freq=5), b(freq=2), c(freq=1), d(freq=3)

Min-heap: [(1,c), (2,b), (3,d), (5,a)]

Step 1: Extract two smallest: c(1) and b(2). Merge into node(3).
        Push (3, cb) back. Heap: [(3,cb), (3,d), (5,a)]

Step 2: Extract (3,cb) and (3,d). Merge into node(6).
        Push (6, cbd). Heap: [(5,a), (6,cbd)]

Step 3: Extract (5,a) and (6,cbd). Merge into root(11).

Final Huffman tree:
        (11)
       /    \
      a(5)  (6)
            / \
           (3) d(3)
           / \
          c(1) b(2)

Codes: a=0, d=10, c=110, b=111
Most frequent (a) gets shortest code.
```

### OS Process Scheduling

The Linux kernel's Completely Fair Scheduler uses a red-black tree (functionally similar to a heap) to track runnable processes by their virtual runtime — the process that has run the least gets to run next.

---

## 9. Go's container/heap Interface

Go's standard library provides a heap in `container/heap`. It requires you to implement a specific interface:

```go
// container/heap requires your type to implement:
// - sort.Interface (Len, Less, Swap)
// - heap.Interface (Push, Pop)

package main

import (
    "container/heap"
    "fmt"
)

// IntMinHeap implements heap.Interface for a min-heap of ints.
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] }  // min-heap
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }

func (h *IntMinHeap) Push(x interface{}) {
    *h = append(*h, x.(int))
}

func (h *IntMinHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func main() {
    h := &IntMinHeap{5, 3, 8, 1, 4}
    heap.Init(h)  // build min-heap from unsorted slice: O(n)

    heap.Push(h, 2)
    heap.Push(h, 7)

    fmt.Println("Extracting in sorted order:")
    for h.Len() > 0 {
        fmt.Printf("%d ", heap.Pop(h).(int))
    }
    // Output: 1 2 3 4 5 7 8
    fmt.Println()
}
```

The `container/heap` package is low-level. For a nicer API, we build our own generic wrapper.

---

## 10. Generic Heap[T] in Go

Here is a complete, user-friendly generic heap:

```go
package heap

import "fmt"

// Heap[T] is a generic heap (min or max depending on the comparator).
// The comparator `priority(a, b)` returns true if a should be extracted before b.
// For a min-heap: priority = func(a, b T) bool { return a < b }
// For a max-heap: priority = func(a, b T) bool { return a > b }
type Heap[T any] struct {
    items    []T
    priority func(a, b T) bool  // true if a has higher priority than b
}

// New creates a new heap.
func New[T any](priority func(a, b T) bool) *Heap[T] {
    return &Heap[T]{
        items:    make([]T, 0, 16),
        priority: priority,
    }
}

// NewMinHeap creates a min-heap for ordered types.
// For ints: NewMinHeap(func(a, b int) bool { return a < b })
func NewMinHeap[T any](less func(a, b T) bool) *Heap[T] {
    return New(less)
}

// NewMaxHeap creates a max-heap for ordered types.
// For ints: NewMaxHeap(func(a, b int) bool { return a > b })
func NewMaxHeap[T any](greater func(a, b T) bool) *Heap[T] {
    return New(greater)
}

// Push inserts an item into the heap — O(log n).
func (h *Heap[T]) Push(item T) {
    h.items = append(h.items, item)
    h.bubbleUp(len(h.items) - 1)
}

// Pop removes and returns the highest-priority item — O(log n).
func (h *Heap[T]) Pop() (T, bool) {
    if h.IsEmpty() {
        var zero T
        return zero, false
    }
    top := h.items[0]
    last := len(h.items) - 1
    h.items[0] = h.items[last]
    h.items = h.items[:last]
    if !h.IsEmpty() {
        h.bubbleDown(0)
    }
    return top, true
}

// Peek returns the highest-priority item without removing it — O(1).
func (h *Heap[T]) Peek() (T, bool) {
    if h.IsEmpty() {
        var zero T
        return zero, false
    }
    return h.items[0], true
}

// Size returns the number of items.
func (h *Heap[T]) Size() int { return len(h.items) }

// IsEmpty returns true if the heap has no items.
func (h *Heap[T]) IsEmpty() bool { return len(h.items) == 0 }

// BuildFrom creates a heap from an existing slice — O(n).
func BuildFrom[T any](items []T, priority func(a, b T) bool) *Heap[T] {
    h := &Heap[T]{
        items:    make([]T, len(items)),
        priority: priority,
    }
    copy(h.items, items)
    // Heapify: sift down from last internal node to root
    for i := (len(h.items) - 2) / 2; i >= 0; i-- {
        h.bubbleDown(i)
    }
    return h
}

func (h *Heap[T]) parent(i int) int      { return (i - 1) / 2 }
func (h *Heap[T]) leftChild(i int) int   { return 2*i + 1 }
func (h *Heap[T]) rightChild(i int) int  { return 2*i + 2 }

func (h *Heap[T]) bubbleUp(i int) {
    for i > 0 {
        p := h.parent(i)
        if h.priority(h.items[i], h.items[p]) {
            h.items[i], h.items[p] = h.items[p], h.items[i]
            i = p
        } else {
            break
        }
    }
}

func (h *Heap[T]) bubbleDown(i int) {
    n := len(h.items)
    for {
        best := i
        l, r := h.leftChild(i), h.rightChild(i)
        if l < n && h.priority(h.items[l], h.items[best]) {
            best = l
        }
        if r < n && h.priority(h.items[r], h.items[best]) {
            best = r
        }
        if best == i {
            break
        }
        h.items[i], h.items[best] = h.items[best], h.items[i]
        i = best
    }
}

// ToSortedSlice extracts all elements in priority order (destructive).
func (h *Heap[T]) ToSortedSlice() []T {
    result := make([]T, 0, h.Size())
    for !h.IsEmpty() {
        item, _ := h.Pop()
        result = append(result, item)
    }
    return result
}

func (h *Heap[T]) String() string {
    return fmt.Sprintf("Heap%v", h.items)
}
```

Usage example:

```go
package main

import (
    "fmt"
    "your-module/heap"
)

func main() {
    // Min-heap of integers
    minH := heap.NewMinHeap(func(a, b int) bool { return a < b })
    for _, v := range []int{5, 3, 8, 1, 4, 7, 2} {
        minH.Push(v)
    }
    fmt.Println("Min:", minH.ToSortedSlice())  // [1 2 3 4 5 7 8]

    // Max-heap of strings by length
    maxH := heap.New(func(a, b string) bool { return len(a) > len(b) })
    for _, s := range []string{"go", "python", "c", "javascript", "rust"} {
        maxH.Push(s)
    }
    fmt.Print("Longest first: ")
    for !maxH.IsEmpty() {
        s, _ := maxH.Pop()
        fmt.Printf("%q ", s)
    }
    // Output: "javascript" "python" "rust" "go" "c"
    fmt.Println()

    // Priority queue of tasks
    type Task struct {
        Name     string
        Priority int  // higher = more urgent
    }
    tasks := heap.New(func(a, b Task) bool { return a.Priority > b.Priority })
    tasks.Push(Task{"Write tests", 2})
    tasks.Push(Task{"Fix prod bug", 10})
    tasks.Push(Task{"Update docs", 1})
    tasks.Push(Task{"Code review", 5})

    fmt.Println("\nTask execution order:")
    for !tasks.IsEmpty() {
        t, _ := tasks.Pop()
        fmt.Printf("  [priority %2d] %s\n", t.Priority, t.Name)
    }
    // [priority 10] Fix prod bug
    // [priority  5] Code review
    // [priority  2] Write tests
    // [priority  1] Update docs
}
```

---

## 11. D-ary Heaps

A **d-ary heap** (or d-heap) is a generalization of the binary heap where each node has at most **d** children (instead of just 2).

Index arithmetic for a d-ary heap (0-indexed):
- **Parent of i**: `(i - 1) / d`
- **Children of i**: `d*i + 1` through `d*i + d`

```
4-ary heap (each node has up to 4 children):

            1
     / |  |  \
    2  3  4   5
  / | | \
 6  7 8  9

Index: 0  1  2  3  4  5  6  7  8  9
Value: 1  2  3  4  5  6  7  8  9 10

Children of index 0: 1, 2, 3, 4
Children of index 1: 5, 6, 7, 8
Children of index 2: 9, 10, 11, 12 (if they exist)
```

**Trade-offs of d-ary heaps:**
- **Shallower tree**: height is log_d(n) instead of log_2(n) — fewer levels
- **Faster push (bubbleUp)**: fewer levels to traverse
- **Slower pop (bubbleDown)**: must compare d children instead of 2 at each level
- **Better cache performance**: larger d = more siblings in the same cache line

D-ary heaps (typically d=4 or d=8) are used in high-performance priority queues because they hit cache lines more efficiently than binary heaps.

The optimal d depends on the cache line size (typically 64 bytes) and element size. For 8-byte elements, d=8 means all children of a node fit in a single 64-byte cache line.

---

## 12. Astra Build Milestone: Optimization Pass Ordering

After parsing and type checking, the Astra compiler runs a series of **optimization passes** on the program's intermediate representation (IR). Each pass transforms the IR to make the final code faster or smaller.

The key insight: optimization passes have **dependencies and priorities**. You should run:
1. **Constant folding** first (simplify `2 + 3` to `5`) — it reveals opportunities for other passes
2. **Dead code elimination** second (remove code that can never run) — enabled by constant folding
3. **Inlining** later (replace function calls with the function's body) — more expensive, needs folding/DCE first
4. **Register allocation** last (very expensive, only useful after all other optimizations)

A priority queue ensures passes run in the right order automatically:

```go
// compiler/optimizer.go

package compiler

import (
    "fmt"
    "your-module/heap"
)

// IR represents the compiler's intermediate representation of a program.
// (We'll build this fully in later chapters.)
type IR struct {
    Name         string
    Instructions []string  // simplified for this example
}

func (ir *IR) Clone() *IR {
    clone := &IR{Name: ir.Name, Instructions: make([]string, len(ir.Instructions))}
    copy(clone.Instructions, ir.Instructions)
    return clone
}

// OptimizationPass represents a single transformation on the IR.
type OptimizationPass struct {
    Name     string
    Priority int              // higher number = runs first
    Run      func(*IR) *IR   // the actual transformation
}

// PassQueue is a max-priority-heap of OptimizationPass (higher priority first).
type PassQueue = heap.Heap[*OptimizationPass]

// NewPassQueue creates a priority queue for optimization passes.
func NewPassQueue() *PassQueue {
    return heap.New(func(a, b *OptimizationPass) bool {
        return a.Priority > b.Priority  // higher priority runs first
    })
}

// ConstantFoldingPass evaluates constant expressions at compile time.
// e.g., 2 + 3 → 5, true && false → false
func ConstantFoldingPass() *OptimizationPass {
    return &OptimizationPass{
        Name:     "constant-folding",
        Priority: 100,  // highest priority: run first
        Run: func(ir *IR) *IR {
            result := ir.Clone()
            // In a real compiler, this walks the AST/IR and folds constants.
            fmt.Printf("  [constant-folding] Processing %d instructions\n",
                len(result.Instructions))
            // Simplified: mark folded instructions
            for i, inst := range result.Instructions {
                if inst == "CONST_ADD" {
                    result.Instructions[i] = "CONST_FOLDED"
                }
            }
            return result
        },
    }
}

// DeadCodeElimPass removes code that can never be executed.
// e.g., code after unconditional return, branches where condition is constant false
func DeadCodeElimPass() *OptimizationPass {
    return &OptimizationPass{
        Name:     "dead-code-elimination",
        Priority: 80,
        Run: func(ir *IR) *IR {
            result := ir.Clone()
            fmt.Printf("  [dead-code-elimination] Processing %d instructions\n",
                len(result.Instructions))
            // Simplified: remove marked-dead instructions
            live := result.Instructions[:0]
            for _, inst := range result.Instructions {
                if inst != "DEAD" {
                    live = append(live, inst)
                }
            }
            result.Instructions = live
            return result
        },
    }
}

// InliningPass replaces function call sites with the function body.
// Can increase code size, so runs after DCE removes unnecessary calls.
func InliningPass(maxInlineSize int) *OptimizationPass {
    return &OptimizationPass{
        Name:     "inlining",
        Priority: 50,
        Run: func(ir *IR) *IR {
            result := ir.Clone()
            fmt.Printf("  [inlining] Inlining functions smaller than %d instructions\n",
                maxInlineSize)
            return result
        },
    }
}

// StrengthReductionPass replaces expensive operations with cheaper equivalents.
// e.g., x * 2 → x + x, or x * 4 → x << 2
func StrengthReductionPass() *OptimizationPass {
    return &OptimizationPass{
        Name:     "strength-reduction",
        Priority: 60,
        Run: func(ir *IR) *IR {
            result := ir.Clone()
            fmt.Printf("  [strength-reduction] Replacing expensive ops\n")
            return result
        },
    }
}

// CommonSubexprElimPass avoids recomputing the same expression twice.
// e.g., a+b+c where a+b appears twice → compute once, reuse
func CommonSubexprElimPass() *OptimizationPass {
    return &OptimizationPass{
        Name:     "common-subexpr-elimination",
        Priority: 70,
        Run: func(ir *IR) *IR {
            result := ir.Clone()
            fmt.Printf("  [common-subexpr-elimination] Eliminating redundant computations\n")
            return result
        },
    }
}

// LoopInvariantCodeMotionPass moves loop-invariant code outside the loop.
// e.g., for i in 0..n { let x = a + b; ... } → let x = a + b; for i in 0..n { ... }
func LoopInvariantCodeMotionPass() *OptimizationPass {
    return &OptimizationPass{
        Name:     "loop-invariant-code-motion",
        Priority: 40,  // runs later: needs inlining results first
        Run: func(ir *IR) *IR {
            result := ir.Clone()
            fmt.Printf("  [loop-invariant-code-motion] Hoisting loop-invariant code\n")
            return result
        },
    }
}

// Optimizer runs all registered passes in priority order.
type Optimizer struct {
    passes *PassQueue
    debug  bool
}

// NewOptimizer creates a new optimizer.
func NewOptimizer(debug bool) *Optimizer {
    return &Optimizer{
        passes: NewPassQueue(),
        debug:  debug,
    }
}

// Register adds an optimization pass to the queue.
func (o *Optimizer) Register(pass *OptimizationPass) {
    o.passes.Push(pass)
}

// RegisterAll adds all standard optimization passes.
func (o *Optimizer) RegisterAll() {
    o.Register(ConstantFoldingPass())
    o.Register(DeadCodeElimPass())
    o.Register(CommonSubexprElimPass())
    o.Register(StrengthReductionPass())
    o.Register(InliningPass(50))
    o.Register(LoopInvariantCodeMotionPass())
}

// Optimize runs all passes in priority order, returning the optimized IR.
func (o *Optimizer) Optimize(ir *IR) *IR {
    if o.debug {
        fmt.Printf("Starting optimization of %q (%d passes queued)\n",
            ir.Name, o.passes.Size())
    }

    current := ir
    passCount := 0

    for !o.passes.IsEmpty() {
        pass, _ := o.passes.Pop()
        if o.debug {
            fmt.Printf("Pass %d: %s (priority %d)\n",
                passCount+1, pass.Name, pass.Priority)
        }
        current = pass.Run(current)
        passCount++
    }

    if o.debug {
        fmt.Printf("Optimization complete: ran %d passes\n", passCount)
    }
    return current
}
```

Usage in the Astra compiler's main compilation pipeline:

```go
package main

import "fmt"

func main() {
    // Create a sample IR (simplified)
    ir := &IR{
        Name: "main.as",
        Instructions: []string{
            "CONST_ADD",  // will be folded
            "LOAD x",
            "CONST_ADD",  // another constant add
            "DEAD",       // dead code
            "CALL add",   // will be inlined
            "STORE result",
        },
    }

    fmt.Println("=== Before Optimization ===")
    fmt.Println("Instructions:", ir.Instructions)
    fmt.Println()

    // Set up optimizer with all standard passes
    optimizer := NewOptimizer(true)
    optimizer.RegisterAll()

    fmt.Println("=== Running Optimization Passes ===")
    optimized := optimizer.Optimize(ir)

    fmt.Println()
    fmt.Println("=== After Optimization ===")
    fmt.Println("Instructions:", optimized.Instructions)
}
```

Expected output:

```
=== Before Optimization ===
Instructions: [CONST_ADD LOAD x CONST_ADD DEAD CALL add STORE result]

=== Running Optimization Passes ===
Starting optimization of "main.as" (6 passes queued)
Pass 1: constant-folding (priority 100)
  [constant-folding] Processing 6 instructions
Pass 2: dead-code-elimination (priority 80)
  [dead-code-elimination] Processing 6 instructions
Pass 3: common-subexpr-elimination (priority 70)
  [common-subexpr-elimination] Eliminating redundant computations
Pass 4: strength-reduction (priority 60)
  [strength-reduction] Replacing expensive ops
Pass 5: inlining (priority 50)
  [inlining] Inlining functions smaller than 50 instructions
Pass 6: loop-invariant-code-motion (priority 40)
  [loop-invariant-code-motion] Hoisting loop-invariant code
Optimization complete: ran 6 passes

=== After Optimization ===
Instructions: [CONST_FOLDED LOAD x CONST_FOLDED CALL add STORE result]
```

The priority queue guarantees that constant folding always runs before dead code elimination (which needs folding results to identify dead branches), which runs before inlining (which should inline only live code), which runs before loop optimization (which assumes inlining has already expanded loops).

Adding a new optimization pass is as simple as:
```go
optimizer.Register(&OptimizationPass{
    Name:     "my-new-pass",
    Priority: 75,  // runs between dead-code-elim (80) and strength-reduction (60)
    Run:      myNewPassFunction,
})
```

The priority queue automatically inserts it in the right position. No manual ordering code needed.

---

## Exercises

1. **K largest elements**: Given an array of n integers, find the k largest elements in O(n log k) time using a min-heap of size k.

2. **Merge k sorted lists**: Given k sorted arrays, merge them into one sorted array in O(n log k) where n is the total number of elements. Use a min-heap.

3. **Running median**: Given a stream of integers, find the median after each insertion. Maintain two heaps: a max-heap for the lower half and a min-heap for the upper half.

4. **Lazy deletion heap**: Implement a heap that supports `decrease-key` (changing a value to a smaller one) without O(n) search. Use a "lazy" approach: mark old entries as invalid and ignore them when they're popped.

5. **Event-driven simulation**: Implement a discrete event simulator. Events have a timestamp and a handler function. Use a min-heap ordered by timestamp to process events in chronological order.

6. **Implement heap sort**: Implement heap sort without using your `Heap[T]` class — directly manipulate the array with the heapify algorithms.

7. **Top k frequent words**: Given a long text, find the k most frequent words in O(n + k log n) time using a heap.

8. **Astra extension**: Add a `Cancel(passName string)` method to the `Optimizer` that removes a pass from the queue before it runs. How do you efficiently find and remove an element from the middle of a heap? (Hint: swap with the last element, then heapify.)

---

## Summary

| Concept                | Key Point                                                       |
|------------------------|-----------------------------------------------------------------|
| Heap                   | Complete binary tree with heap property (parent dominates children) |
| Max-heap               | Root = maximum element; extract gives maximum                  |
| Min-heap               | Root = minimum element; extract gives minimum                  |
| Array representation   | Parent at i → children at 2i+1, 2i+2; parent at (i-1)/2      |
| Push (heapify-up)      | Append + bubble up to restore heap property — O(log n)         |
| Pop (heapify-down)     | Swap root with last + bubble down — O(log n)                   |
| Build from array       | Heapify-down from last internal node to root — O(n)            |
| Heap sort              | Build heap + repeated extract → O(n log n), O(1) space         |
| Priority queue         | Heap-backed queue serving highest-priority item first           |
| Dijkstra's algorithm   | Uses min-heap for O((V+E) log V) shortest paths               |
| Huffman coding         | Uses min-heap to build optimal prefix codes                     |
| D-ary heap             | d children per node; shallower tree, better cache, slower pop  |
| Astra usage            | Priority queue orders optimization passes by dependency/cost    |
