---
title: Queue
step: 3
difficulty: easy
estimated: 25 min
---

## What You Are Building

A queue is a First-In, First-Out (FIFO) collection. Items enter at the back and leave from the front — like a line at a counter. Queues are essential for BFS, job schedulers, message brokers, and rate limiters.

```
Enqueue(1) → front:[1]:back
Enqueue(2) → front:[1, 2]:back
Enqueue(3) → front:[1, 2, 3]:back
Dequeue()  → 1, true   (queue: [2, 3])
Peek()     → 2, true   (queue unchanged)
```

## Key Concepts

**The naive slice approach** — You can back a queue with a slice where `Enqueue` appends to the end and `Dequeue` returns `slice[0]` and re-slices. This is fine for this task. The downside is that `Dequeue` does not free the memory at index 0 (the backing array grows forever), but correctness is what we care about here.

```go
// Enqueue
q.data = append(q.data, val)

// Dequeue
val := q.data[0]
q.data = q.data[1:]
```

**Two-pointer / two-stack approach (bonus)** — A more memory-efficient design uses two stacks: one for enqueueing (push-stack) and one for dequeueing (pop-stack). When the pop-stack empties, reverse the push-stack into it. Amortized O(1) per operation. Try this after the basic version works.

**Same comma-ok pattern as Stack** — `Dequeue` and `Peek` return `(int, bool)` because there is no sensible value when the queue is empty.

## Struct Signature

```go
type Queue struct {
    data []int
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Enqueue(val int)` | Add val to the back |
| `Dequeue() (int, bool)` | Remove and return front; ok=false if empty |
| `Peek() (int, bool)` | Return front without removing; ok=false if empty |
| `IsEmpty() bool` | True if no elements |
| `Size() int` | Number of elements |

## Edge Cases to Handle

- `Dequeue` on empty queue: return `0, false`
- `Peek` on empty queue: return `0, false`
- Multiple enqueue/dequeue cycles should work correctly

## Example

```go
q := &Queue{}
fmt.Println(q.IsEmpty()) // true

q.Enqueue(10)
q.Enqueue(20)
q.Enqueue(30)
fmt.Println(q.Size()) // 3

front, ok := q.Peek()
fmt.Println(front, ok) // 10 true

val, ok := q.Dequeue()
fmt.Println(val, ok) // 10 true

val, ok = q.Dequeue()
fmt.Println(val, ok) // 20 true

fmt.Println(q.Size()) // 1
```

## Hints

- A queue is the mirror image of a stack: push goes on the right (append), pop comes from the left (index 0).
- `Dequeue`: guard with `len(q.data) == 0` before accessing `q.data[0]`.
- `Peek`: same guard, but return `q.data[0]` without modifying the slice.
