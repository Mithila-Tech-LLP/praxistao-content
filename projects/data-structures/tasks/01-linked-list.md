---
title: Linked List
step: 1
difficulty: easy
estimated: 30 min
---

## What You Are Building

A singly linked list — the most fundamental pointer-based data structure. Each element (a *node*) holds a value and a pointer to the next node. The list itself only tracks the head.

```
Head → [3] → [7] → [12] → nil
```

## Key Concepts

**Pointers in Go** — A `*Node` is either a valid memory address or `nil`. Before you dereference a pointer (write `n.Next`) you must check that `n != nil`, or your program will panic.

**Traversal pattern** — Almost every linked list operation starts the same way:

```go
current := l.Head
for current != nil {
    // do something with current.Val
    current = current.Next
}
```

**Append vs Prepend** — Prepend is O(1): create a node, point it at `Head`, update `Head`. Append is O(n): you must walk the whole list to find the last node.

**Deleting a node** — You cannot delete a node directly. You find the node *before* it and update that node's `Next` pointer to skip over the target. The special case is deleting the head, where you update `l.Head` instead.

## Struct Signatures

```go
type Node struct {
    Val  int
    Next *Node
}

type LinkedList struct {
    Head *Node
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Append(val int)` | Add val at the tail |
| `Prepend(val int)` | Add val at the head |
| `Delete(val int)` | Remove first occurrence of val |
| `Contains(val int) bool` | Return true if val exists |
| `ToSlice() []int` | Return all values in order |

## Edge Cases to Handle

- `Delete` on an empty list should do nothing (no panic)
- `Delete` a value that does not exist should do nothing
- `Delete` the head node: update `l.Head`
- `ToSlice` on an empty list returns `[]int{}`

## Example

```go
l := &LinkedList{}
l.Append(1)
l.Append(2)
l.Append(3)
l.Prepend(0)
fmt.Println(l.ToSlice())   // [0 1 2 3]

l.Delete(2)
fmt.Println(l.ToSlice())   // [0 1 3]
fmt.Println(l.Contains(3)) // true
fmt.Println(l.Contains(2)) // false
```

## Hints

- For `Append`, walk until you find the node where `current.Next == nil`, then set `current.Next = &Node{Val: val}`.
- For `Delete`, keep a `prev` pointer one step behind `current` so you can update `prev.Next = current.Next`.
- You need a nil-guard at the start of `Delete`: if `l.Head == nil`, return immediately.
