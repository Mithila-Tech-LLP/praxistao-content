# Chapter 31: Linked Lists

A linked list is a chain of nodes where each node holds a value and a pointer to the next node. Unlike arrays (contiguous memory), linked list nodes can live anywhere in memory — connected by pointers. This makes insertions and deletions at the head O(1), while random access is O(n). Linked lists are the foundation for stacks, queues, and more complex data structures.

## Table of Contents

1. [Singly Linked List](#1-singly-linked-list)
2. [Common Operations](#2-common-operations)
3. [Doubly Linked List](#3-doubly-linked-list)
4. [Circular Linked List](#4-circular-linked-list)
5. [Classic Interview Problems](#5-classic-interview-problems)
6. [When to Use Linked Lists](#6-when-to-use-linked-lists)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Singly Linked List

Each node holds a value and a pointer to the next node:

```
head → [1|next] → [2|next] → [3|next] → nil
```

```go
package main

import "fmt"

// Node is a single element in the linked list.
type Node[T any] struct {
    Value T
    Next  *Node[T]
}

// LinkedList holds the head pointer and length.
type LinkedList[T any] struct {
    Head *Node[T]
    size int
}

func (l *LinkedList[T]) Len() int { return l.size }

func (l *LinkedList[T]) IsEmpty() bool { return l.Head == nil }
```

**Building a list manually:**
```go
// Manual construction:
n3 := &Node[int]{Value: 3}
n2 := &Node[int]{Value: 2, Next: n3}
n1 := &Node[int]{Value: 1, Next: n2}

list := &LinkedList[int]{Head: n1, size: 3}

// Traverse:
for cur := list.Head; cur != nil; cur = cur.Next {
    fmt.Printf("%d → ", cur.Value)
}
fmt.Println("nil")
// Output: 1 → 2 → 3 → nil
```

### Quick Check
> 1. What does a singly linked list node contain?
> 2. What is the time complexity of traversing a linked list?
> 3. What value does the last node's `Next` field hold?

---

## 2. Common Operations

### Prepend — O(1)
```go
// PushFront inserts a value at the head.
func (l *LinkedList[T]) PushFront(val T) {
    node := &Node[T]{Value: val, Next: l.Head}
    l.Head = node
    l.size++
}
```

### Append — O(n) without tail pointer
```go
// PushBack inserts a value at the tail.
func (l *LinkedList[T]) PushBack(val T) {
    node := &Node[T]{Value: val}
    if l.Head == nil {
        l.Head = node
        l.size++
        return
    }
    // Traverse to last node:
    cur := l.Head
    for cur.Next != nil {
        cur = cur.Next
    }
    cur.Next = node
    l.size++
}
```

### Pop Front — O(1)
```go
// PopFront removes and returns the head value.
func (l *LinkedList[T]) PopFront() (T, bool) {
    if l.Head == nil {
        var zero T
        return zero, false
    }
    val := l.Head.Value
    l.Head = l.Head.Next
    l.size--
    return val, true
}
```

### Delete by value — O(n)
```go
// Delete removes the first node with the given value.
func (l *LinkedList[T]) Delete(val T, eq func(T, T) bool) bool {
    if l.Head == nil {
        return false
    }
    // Special case: deleting head
    if eq(l.Head.Value, val) {
        l.Head = l.Head.Next
        l.size--
        return true
    }
    // General case: find the node before the target
    prev := l.Head
    for prev.Next != nil {
        if eq(prev.Next.Value, val) {
            prev.Next = prev.Next.Next  // Skip over the deleted node
            l.size--
            return true
        }
        prev = prev.Next
    }
    return false
}
```

### Insert at position — O(n)
```go
// Insert inserts val at position pos (0-indexed). O(n)
func (l *LinkedList[T]) Insert(pos int, val T) error {
    if pos < 0 || pos > l.size {
        return fmt.Errorf("position %d out of range [0, %d]", pos, l.size)
    }
    if pos == 0 {
        l.PushFront(val)
        return nil
    }
    // Find node at pos-1:
    cur := l.Head
    for i := 0; i < pos-1; i++ {
        cur = cur.Next
    }
    node := &Node[T]{Value: val, Next: cur.Next}
    cur.Next = node
    l.size++
    return nil
}
```

### Search — O(n)
```go
// Find returns the first node whose value satisfies the predicate.
func (l *LinkedList[T]) Find(pred func(T) bool) (*Node[T], int) {
    for i, cur := 0, l.Head; cur != nil; i, cur = i+1, cur.Next {
        if pred(cur.Value) {
            return cur, i
        }
    }
    return nil, -1
}
```

### Reverse — O(n)
```go
// Reverse reverses the list in-place.
func (l *LinkedList[T]) Reverse() {
    var prev *Node[T]
    cur := l.Head
    for cur != nil {
        next := cur.Next  // Save next
        cur.Next = prev   // Reverse pointer
        prev = cur        // Move prev forward
        cur = next        // Move cur forward
    }
    l.Head = prev
}

// Visualizing the reversal:
// Before: head → [1] → [2] → [3] → nil
// Step 1: prev=nil, cur=[1], next=[2] → [1].next=nil, prev=[1], cur=[2]
// Step 2: prev=[1], cur=[2], next=[3] → [2].next=[1], prev=[2], cur=[3]
// Step 3: prev=[2], cur=[3], next=nil → [3].next=[2], prev=[3], cur=nil
// After:  head → [3] → [2] → [1] → nil
```

### Convert to slice — O(n)
```go
func (l *LinkedList[T]) ToSlice() []T {
    result := make([]T, 0, l.size)
    for cur := l.Head; cur != nil; cur = cur.Next {
        result = append(result, cur.Value)
    }
    return result
}
```

**Complete working example:**
```go
func main() {
    list := &LinkedList[int]{}
    list.PushBack(1)
    list.PushBack(2)
    list.PushBack(3)
    list.PushFront(0)

    fmt.Println(list.ToSlice())  // [0 1 2 3]
    fmt.Println("len:", list.Len())  // 4

    list.Reverse()
    fmt.Println(list.ToSlice())  // [3 2 1 0]

    list.Delete(2, func(a, b int) bool { return a == b })
    fmt.Println(list.ToSlice())  // [3 1 0]
}
```

### Quick Check
> 1. Why is prepend O(1) but append O(n) for a singly linked list without a tail pointer?
> 2. In the reverse algorithm, what do the `prev`, `cur`, and `next` variables track?
> 3. Deleting the head is a special case — why?

---

## 3. Doubly Linked List

Each node has pointers to both the next AND previous nodes. This enables O(1) deletion given a pointer to any node:

```
nil ← [1|prev|next] ↔ [2|prev|next] ↔ [3|prev|next] → nil
       ↑ head                                 ↑ tail
```

```go
type DNode[T any] struct {
    Value T
    Prev  *DNode[T]
    Next  *DNode[T]
}

type DoublyLinkedList[T any] struct {
    Head *DNode[T]
    Tail *DNode[T]
    size int
}

func (l *DoublyLinkedList[T]) Len() int { return l.size }

// PushBack is now O(1) because we have a tail pointer:
func (l *DoublyLinkedList[T]) PushBack(val T) {
    node := &DNode[T]{Value: val, Prev: l.Tail}
    if l.Tail != nil {
        l.Tail.Next = node
    } else {
        l.Head = node  // Empty list
    }
    l.Tail = node
    l.size++
}

// PushFront:
func (l *DoublyLinkedList[T]) PushFront(val T) {
    node := &DNode[T]{Value: val, Next: l.Head}
    if l.Head != nil {
        l.Head.Prev = node
    } else {
        l.Tail = node  // Empty list
    }
    l.Head = node
    l.size++
}

// PopBack — O(1):
func (l *DoublyLinkedList[T]) PopBack() (T, bool) {
    if l.Tail == nil {
        var zero T
        return zero, false
    }
    val := l.Tail.Value
    l.Tail = l.Tail.Prev
    if l.Tail != nil {
        l.Tail.Next = nil
    } else {
        l.Head = nil  // List is now empty
    }
    l.size--
    return val, true
}

// DeleteNode — O(1) given a pointer to the node:
func (l *DoublyLinkedList[T]) DeleteNode(node *DNode[T]) {
    if node.Prev != nil {
        node.Prev.Next = node.Next
    } else {
        l.Head = node.Next  // Deleting head
    }
    if node.Next != nil {
        node.Next.Prev = node.Prev
    } else {
        l.Tail = node.Prev  // Deleting tail
    }
    node.Prev = nil  // Help GC
    node.Next = nil
    l.size--
}
```

### Quick Check
> 1. What is the advantage of a doubly linked list over singly?
> 2. Why does O(1) deletion require a pointer to the node itself?
> 3. What must you update when deleting the tail of a doubly linked list?

---

## 4. Circular Linked List

The tail's `Next` points back to the head — useful for round-robin scheduling, circular buffers:

```
head → [1] → [2] → [3] → back to head
```

```go
type CircularList[T any] struct {
    Head *Node[T]
    size int
}

func (l *CircularList[T]) PushBack(val T) {
    node := &Node[T]{Value: val}
    if l.Head == nil {
        node.Next = node  // Points to itself
        l.Head = node
        l.size++
        return
    }
    // Find the tail (node whose Next = Head):
    tail := l.Head
    for tail.Next != l.Head {
        tail = tail.Next
    }
    tail.Next = node
    node.Next = l.Head  // Complete the circle
    l.size++
}

// Traverse N steps (circular):
func (l *CircularList[T]) Traverse(steps int) {
    if l.Head == nil {
        return
    }
    cur := l.Head
    for i := 0; i < steps; i++ {
        fmt.Printf("%v ", cur.Value)
        cur = cur.Next  // Wraps around automatically
    }
    fmt.Println()
}
```

---

## 5. Classic Interview Problems

### Detect a cycle (Floyd's algorithm) — O(n) time, O(1) space
```go
// HasCycle detects if a linked list has a cycle using slow/fast pointers.
func HasCycle[T any](head *Node[T]) bool {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next        // Move one step
        fast = fast.Next.Next   // Move two steps
        if slow == fast {
            return true  // They met — cycle exists
        }
    }
    return false
}

// Visual: If cycle exists, fast laps slow. They'll eventually meet.
// No cycle:  fast reaches nil.
```

### Find the middle node — O(n) time, O(1) space
```go
// FindMiddle returns the middle node (second middle for even-length lists).
func FindMiddle[T any](head *Node[T]) *Node[T] {
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow  // Slow is at the middle when fast reaches the end
}

// [1→2→3→4→5]: slow=3 when fast=5 → middle is 3
// [1→2→3→4]:   slow=3 when fast=nil → middle is 3 (second)
```

### Merge two sorted lists — O(n+m)
```go
func MergeSorted[T any](a, b *Node[T], less func(T, T) bool) *Node[T] {
    dummy := &Node[T]{}  // Sentinel head — avoids special-casing the first node
    cur := dummy

    for a != nil && b != nil {
        if less(a.Value, b.Value) {
            cur.Next = a
            a = a.Next
        } else {
            cur.Next = b
            b = b.Next
        }
        cur = cur.Next
    }

    if a != nil {
        cur.Next = a
    } else {
        cur.Next = b
    }

    return dummy.Next
}
```

### Find Nth node from end — O(n) time, O(1) space
```go
// NthFromEnd returns the Nth node from the end (1-indexed).
func NthFromEnd[T any](head *Node[T], n int) *Node[T] {
    ahead, behind := head, head

    // Advance `ahead` n steps first:
    for i := 0; i < n; i++ {
        if ahead == nil {
            return nil  // n > list length
        }
        ahead = ahead.Next
    }

    // Now move both until `ahead` reaches the end:
    for ahead != nil {
        ahead = ahead.Next
        behind = behind.Next
    }

    return behind
}
// When ahead=nil, behind is n steps behind = Nth from end
```

### Check palindrome — O(n) time, O(1) space
```go
func IsPalindrome(head *Node[int]) bool {
    if head == nil || head.Next == nil {
        return true
    }

    // 1. Find middle:
    mid := FindMiddle(head)

    // 2. Reverse the second half in-place:
    var prev *Node[int]
    cur := mid
    for cur != nil {
        next := cur.Next
        cur.Next = prev
        prev = cur
        cur = next
    }
    secondHalf := prev

    // 3. Compare first and reversed second halves:
    p1, p2 := head, secondHalf
    for p2 != nil {
        if p1.Value != p2.Value {
            return false
        }
        p1 = p1.Next
        p2 = p2.Next
    }
    return true
}
```

### Quick Check
> 1. Why does Floyd's cycle detection use a slow (1-step) and fast (2-step) pointer?
> 2. What is the "dummy head" trick and why does it simplify merge?
> 3. In the "Nth from end" algorithm, why advance `ahead` by N steps first?

---

## 6. When to Use Linked Lists

| Operation | Array/Slice | Linked List |
|-----------|------------|-------------|
| Access by index | O(1) ✓ | O(n) |
| Search by value | O(n) | O(n) |
| Insert at head | O(n) (shift) | O(1) ✓ |
| Insert at tail | O(1) amortized | O(1) with tail ptr |
| Insert at middle | O(n) | O(n) traverse + O(1) insert |
| Delete at head | O(n) (shift) | O(1) ✓ |
| Memory | Contiguous, cache-friendly | Scattered, pointer overhead |

**Use linked lists when:**
- Frequent insertions/deletions at the head (stack, queue)
- You have a pointer to the node you want to delete (O(1) with doubly-linked)
- You need to merge two sorted sequences efficiently
- Implementing LRU caches (doubly-linked + hash map)

**Prefer slices when:**
- You need random access
- Cache performance matters
- The collection size is known upfront

In practice, **Go's slices beat linked lists for most use cases** due to cache locality. Use linked lists only when the access pattern specifically benefits from O(1) head operations or pointer-based deletion.

---

## Summary

- **Singly linked list**: each node has value + `Next`; prepend O(1), append O(n) without tail
- **Doubly linked list**: each node has `Prev` + `Next`; append and delete are O(1) with tail ptr
- **Circular**: tail's `Next` points to head — useful for round-robin
- **Reverse**: three-pointer technique (prev, cur, next)
- **Cycle detection**: Floyd's slow/fast pointer — O(n) time, O(1) space
- **Find middle**: slow/fast pointer — when fast reaches end, slow is at middle
- **Merge sorted**: dummy sentinel head, then zip the two lists
- **Nth from end**: two-pointer with N-step head start

---

## Exercises

### Easy
1. Implement `ToSlice` and `FromSlice[T any]([]T) *LinkedList[T]`. Verify round-trip: `FromSlice([1,2,3]).ToSlice() == [1,2,3]`.
2. Write `RemoveDuplicates[T comparable](head *Node[T]) *Node[T]` that removes duplicate values from a sorted list in-place. `[1→1→2→3→3→nil]` → `[1→2→3→nil]`.
3. Write `Print[T any](head *Node[T])` that prints `1 → 2 → 3 → nil`. Then write `PrintReverse[T any](head *Node[T])` using recursion.

### Medium
4. LRU eviction list: Implement a doubly linked list used as the internal data structure of an LRU cache. `MoveToFront(node *DNode[K])` moves a node to the front in O(1). `RemoveLast() *DNode[K]` removes and returns the tail (least recently used) in O(1). Write tests verifying order after multiple moves.
5. Flatten a nested list: Given a linked list where each node can have a `Child *Node[T]` pointer to another linked list, flatten it into a single list (depth-first). Example: `[1→2→[3→4]→5]` → `[1→2→3→4→5]`. Handle arbitrarily deep nesting.
6. Merge K sorted lists: Given K sorted linked lists (as a `[]*Node[int]`), merge them all into one sorted list. Naive: O(N×K). Optimal: use a min-heap to always pick the smallest — O(N log K). Implement both and benchmark for K=10, N=1000.

### Hard
7. Skip list: Implement a skip list — a probabilistic data structure offering O(log n) average search, insert, and delete. A skip list is a linked list with multiple levels, where higher levels are "express lanes" skipping over nodes. `Insert(val int)` with random level generation, `Search(val int) bool`, `Delete(val int) bool`, `Range(min, max int) []int`. Benchmark against a sorted slice for 100K elements.
8. Persistent linked list: Implement a persistent (immutable) linked list where each mutation creates a new list sharing unchanged nodes with the old list (structural sharing). `Cons(val T, list *PList[T]) *PList[T]` prepends. `Head(list) T` and `Tail(list) *PList[T]` access. This requires NO mutation of existing nodes. Verify: `list1 := Cons(1, empty)`, `list2 := Cons(2, list1)` — both `list1` and `list2` remain valid and list1 is unchanged.
