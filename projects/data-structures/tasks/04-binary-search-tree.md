---
title: Binary Search Tree
step: 4
difficulty: medium
estimated: 45 min
---

## What You Are Building

A Binary Search Tree (BST) organises values so that every node's left subtree contains only smaller values and every node's right subtree contains only larger values. This ordering makes search, insert, and traversal efficient — O(log n) on average.

```
      5
     / \
    3   8
   / \   \
  1   4   9
```

## Key Concepts

**The BST property** — For any node N:
- All values in N's left subtree are less than N.Val
- All values in N's right subtree are greater than N.Val
- Duplicates: we will skip them (if the value already exists, insert is a no-op)

**Recursion is natural here** — Every BST operation can be expressed as: "handle the base case (nil node), otherwise recurse left or right."

```go
func insert(node *BSTNode, val int) *BSTNode {
    if node == nil {
        return &BSTNode{Val: val}
    }
    if val < node.Val {
        node.Left = insert(node.Left, val)
    } else if val > node.Val {
        node.Right = insert(node.Right, val)
    }
    return node
}
```

**In-order traversal** — Visit left, then current, then right. For a BST this visits all nodes in sorted ascending order. This is the key insight: an in-order traversal of a BST is always sorted.

**Finding Min/Max** — Min is the leftmost node (keep going left until `node.Left == nil`). Max is the rightmost node (keep going right).

## Struct Signatures

```go
type BSTNode struct {
    Val   int
    Left  *BSTNode
    Right *BSTNode
}

type BST struct {
    Root *BSTNode
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Insert(val int)` | Add val; ignore duplicates |
| `Search(val int) bool` | True if val exists |
| `InOrder() []int` | Return all values sorted ascending |
| `Min() (int, bool)` | Smallest value; ok=false if empty |
| `Max() (int, bool)` | Largest value; ok=false if empty |

## Edge Cases to Handle

- `Search` on empty tree: return `false`
- `Min`/`Max` on empty tree: return `0, false`
- Inserting a duplicate: do nothing (no second node)
- `InOrder` on empty tree: return `[]int{}`

## Example

```go
bst := &BST{}
for _, v := range []int{5, 3, 8, 1, 4, 9} {
    bst.Insert(v)
}

fmt.Println(bst.InOrder())  // [1 3 4 5 8 9]
fmt.Println(bst.Search(4))  // true
fmt.Println(bst.Search(7))  // false

min, _ := bst.Min()
max, _ := bst.Max()
fmt.Println(min, max)       // 1 9
```

## Hints

- Write private helper functions that accept a `*BSTNode` parameter. The public methods just call those helpers starting from `bst.Root`.
- For `InOrder`, use a closure that appends to a result slice: `var result []int; var traverse func(*BSTNode); traverse = func(n *BSTNode) { ... }`.
- For `Min`, walk left until `node.Left == nil`, then return `node.Val`.
