---
title: Min Heap
step: 7
difficulty: medium
estimated: 40 min
---

## What You Are Building

A min-heap is a complete binary tree where every parent is smaller than or equal to its children. The minimum element is always at the root. Heaps power priority queues, heap sort, Dijkstra's algorithm, and scheduling systems.

```
       1
      / \
     3   5
    / \
   7   8
```

## Key Concepts

**Slice-backed tree** — A heap is stored as a flat slice. The tree structure is encoded with index arithmetic:
- Root is at index 0
- Parent of node at index i: `(i - 1) / 2`
- Left child of node at index i: `2*i + 1`
- Right child of node at index i: `2*i + 2`

No `left` or `right` pointers needed — it is all arithmetic.

**Sift-up (after Push)** — After appending to the end of the slice, the new element may be smaller than its parent, violating the heap property. Repeatedly swap it with its parent until it is in the right place:

```go
func (h *MinHeap) siftUp(i int) {
    for i > 0 {
        parent := (i - 1) / 2
        if h.data[i] < h.data[parent] {
            h.data[i], h.data[parent] = h.data[parent], h.data[i]
            i = parent
        } else {
            break
        }
    }
}
```

**Sift-down (after Pop)** — Pop replaces the root with the last element (maintaining a complete tree), then sifts it down: repeatedly swap with the smaller child until the heap property holds.

**Heap property** — After every push and pop, the smallest element must be at `data[0]`.

## Struct Signature

```go
type MinHeap struct {
    data []int
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Push(val int)` | Insert val maintaining heap property |
| `Pop() (int, bool)` | Remove and return min; ok=false if empty |
| `Peek() (int, bool)` | Return min without removing; ok=false if empty |
| `Size() int` | Number of elements |

## Edge Cases to Handle

- `Pop` on empty heap: return `0, false`
- `Peek` on empty heap: return `0, false`
- `Pop` with one element: remove it, return it
- Pushing duplicate values is fine

## Example

```go
h := &MinHeap{}
h.Push(5)
h.Push(3)
h.Push(8)
h.Push(1)
h.Push(4)

min, _ := h.Peek()
fmt.Println(min) // 1

val, _ := h.Pop()
fmt.Println(val) // 1

val, _ = h.Pop()
fmt.Println(val) // 3

fmt.Println(h.Size()) // 3
```

## Hints

- `Push`: `h.data = append(h.data, val)`, then call `siftUp(len(h.data) - 1)`.
- `Pop`: swap `h.data[0]` with the last element, shrink the slice, then call `siftDown(0)`.
- In `siftDown`, at each step find the smaller of the two children (if they exist), and swap with the current node if that child is smaller.
- Always check child index bounds (`2*i+1 < len(h.data)`) before accessing children.
