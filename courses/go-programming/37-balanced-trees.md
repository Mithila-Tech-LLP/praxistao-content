# Chapter 37: Balanced Trees — AVL and Red-Black Trees

A binary search tree with randomly inserted data has O(log n) operations on average. But in the worst case (sorted input), it degenerates to a linked list: O(n). Balanced trees guarantee O(log n) in all cases by automatically rebalancing after insertions and deletions.

## Table of Contents

1. [Why Balance Matters](#1-why-balance-matters)
2. [AVL Trees](#2-avl-trees)
3. [Red-Black Trees](#3-red-black-trees)
4. [AVL vs Red-Black](#4-avl-vs-red-black)
5. [Go's Balanced Tree (btree package)](#5-gos-balanced-tree-btree-package)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Why Balance Matters

```
Sorted insertion into a BST:
Insert 1 → 2 → 3 → 4 → 5

    1
     \
      2
       \
        3
         \
          4
           \
            5

Height = n-1 → O(n) search, insert, delete
```

A balanced tree enforces that the height stays O(log n):

```
After balance:
        3
       / \
      2   4
     /     \
    1       5

Height = 2 → O(log n) operations
```

---

## 2. AVL Trees

An AVL tree is a BST where every node's **balance factor** (height of left subtree − height of right subtree) is in {-1, 0, 1}. After every insertion or deletion, we check and fix the balance factor up the tree.

### Balance factor and rotations

When a node becomes unbalanced (factor ±2), we restore it with one of four rotations:

```
Right rotation (left-heavy):          Left rotation (right-heavy):
    z                                      z
   / \                                    / \
  y   T4           →                    T1   y
 / \                                        / \
x   T3                                    T2   x
/ \                                           / \
T1  T2                                       T3  T4

After right rotate at z:              After left rotate at z:
    y                                      y
   / \                                    / \
  x   z                                  z   x
 / \ / \                                / \ / \
T1 T2 T3 T4                           T1 T2 T3 T4
```

```go
package avl

type Node struct {
    Key    int
    Left   *Node
    Right  *Node
    Height int
}

func height(n *Node) int {
    if n == nil { return 0 }
    return n.Height
}

func updateHeight(n *Node) {
    lh := height(n.Left)
    rh := height(n.Right)
    if lh > rh {
        n.Height = lh + 1
    } else {
        n.Height = rh + 1
    }
}

func balanceFactor(n *Node) int {
    if n == nil { return 0 }
    return height(n.Left) - height(n.Right)
}

// rotateRight performs a right rotation at y.
func rotateRight(y *Node) *Node {
    x := y.Left
    T2 := x.Right

    x.Right = y
    y.Left = T2

    updateHeight(y)
    updateHeight(x)
    return x
}

// rotateLeft performs a left rotation at x.
func rotateLeft(x *Node) *Node {
    y := x.Right
    T2 := y.Left

    y.Left = x
    x.Right = T2

    updateHeight(x)
    updateHeight(y)
    return y
}

// balance restores the AVL property at node n.
func balance(n *Node) *Node {
    updateHeight(n)
    bf := balanceFactor(n)

    // Left-heavy
    if bf > 1 {
        if balanceFactor(n.Left) < 0 {
            n.Left = rotateLeft(n.Left) // Left-Right case
        }
        return rotateRight(n)
    }

    // Right-heavy
    if bf < -1 {
        if balanceFactor(n.Right) > 0 {
            n.Right = rotateRight(n.Right) // Right-Left case
        }
        return rotateLeft(n)
    }

    return n
}

// Insert inserts a key and returns the new root.
func Insert(root *Node, key int) *Node {
    if root == nil {
        return &Node{Key: key, Height: 1}
    }
    if key < root.Key {
        root.Left = Insert(root.Left, key)
    } else if key > root.Key {
        root.Right = Insert(root.Right, key)
    } else {
        return root // duplicate key — no insert
    }
    return balance(root)
}

// Search returns true if key is in the tree.
func Search(root *Node, key int) bool {
    if root == nil { return false }
    if key == root.Key { return true }
    if key < root.Key { return Search(root.Left, key) }
    return Search(root.Right, key)
}

func minNode(n *Node) *Node {
    for n.Left != nil { n = n.Left }
    return n
}

// Delete removes a key and returns the new root.
func Delete(root *Node, key int) *Node {
    if root == nil { return nil }

    if key < root.Key {
        root.Left = Delete(root.Left, key)
    } else if key > root.Key {
        root.Right = Delete(root.Right, key)
    } else {
        // Found the node to delete
        if root.Left == nil { return root.Right }
        if root.Right == nil { return root.Left }
        // Two children: replace with in-order successor (minimum of right subtree)
        successor := minNode(root.Right)
        root.Key = successor.Key
        root.Right = Delete(root.Right, successor.Key)
    }
    return balance(root)
}

// InOrder returns sorted keys (left, root, right).
func InOrder(root *Node) []int {
    if root == nil { return nil }
    left := InOrder(root.Left)
    right := InOrder(root.Right)
    return append(append(left, root.Key), right...)
}
```

### Usage

```go
var root *avl.Node
for _, key := range []int{5, 3, 7, 1, 4, 6, 8, 2} {
    root = avl.Insert(root, key)
}

fmt.Println(avl.InOrder(root)) // [1 2 3 4 5 6 7 8]
fmt.Println(avl.Search(root, 4)) // true
root = avl.Delete(root, 3)
fmt.Println(avl.InOrder(root)) // [1 2 4 5 6 7 8]
```

### AVL height guarantee

After n insertions, height ≤ 1.44 × log₂(n+2). This is the tightest possible guarantee for a balanced BST.

---

## 3. Red-Black Trees

A red-black tree (RBT) trades AVL's tighter balance for simpler rebalancing, making insertions and deletions faster in practice.

**Five invariants every node must satisfy:**
1. Each node is RED or BLACK
2. The root is BLACK
3. Every leaf (nil node) is BLACK
4. Red nodes have only black children (no two consecutive reds)
5. All paths from any node to its leaf descendants have the same number of black nodes (**black-height**)

These invariants guarantee height ≤ 2 × log₂(n+1).

```go
package rbt

type color bool

const (
    RED   color = true
    BLACK color = false
)

type Node struct {
    Key         int
    Color       color
    Left, Right *Node
    Parent      *Node
}

type Tree struct {
    root *Node
    size int
}

func New() *Tree { return &Tree{} }

func (t *Tree) isRed(n *Node) bool {
    return n != nil && n.Color == RED
}

func (t *Tree) rotateLeft(x *Node) {
    y := x.Right
    x.Right = y.Left
    if y.Left != nil { y.Left.Parent = x }
    y.Parent = x.Parent
    if x.Parent == nil {
        t.root = y
    } else if x == x.Parent.Left {
        x.Parent.Left = y
    } else {
        x.Parent.Right = y
    }
    y.Left = x
    x.Parent = y
}

func (t *Tree) rotateRight(y *Node) {
    x := y.Left
    y.Left = x.Right
    if x.Right != nil { x.Right.Parent = y }
    x.Parent = y.Parent
    if y.Parent == nil {
        t.root = x
    } else if y == y.Parent.Right {
        y.Parent.Right = x
    } else {
        y.Parent.Left = x
    }
    x.Right = y
    y.Parent = x
}

func (t *Tree) Insert(key int) {
    n := &Node{Key: key, Color: RED}
    t.bstInsert(n)
    t.fixInsert(n)
    t.size++
}

func (t *Tree) bstInsert(n *Node) {
    var parent *Node
    curr := t.root
    for curr != nil {
        parent = curr
        if n.Key < curr.Key {
            curr = curr.Left
        } else if n.Key > curr.Key {
            curr = curr.Right
        } else {
            return // duplicate
        }
    }
    n.Parent = parent
    if parent == nil {
        t.root = n
    } else if n.Key < parent.Key {
        parent.Left = n
    } else {
        parent.Right = n
    }
}

// fixInsert restores red-black properties after insertion.
func (t *Tree) fixInsert(n *Node) {
    for t.isRed(n.Parent) {
        if n.Parent == n.Parent.Parent.Left {
            uncle := n.Parent.Parent.Right
            if t.isRed(uncle) {
                // Case 1: Uncle is red — recolor
                n.Parent.Color = BLACK
                uncle.Color = BLACK
                n.Parent.Parent.Color = RED
                n = n.Parent.Parent
            } else {
                if n == n.Parent.Right {
                    // Case 2: Node is right child — rotate left
                    n = n.Parent
                    t.rotateLeft(n)
                }
                // Case 3: Node is left child — rotate right
                n.Parent.Color = BLACK
                n.Parent.Parent.Color = RED
                t.rotateRight(n.Parent.Parent)
            }
        } else {
            // Mirror cases
            uncle := n.Parent.Parent.Left
            if t.isRed(uncle) {
                n.Parent.Color = BLACK
                uncle.Color = BLACK
                n.Parent.Parent.Color = RED
                n = n.Parent.Parent
            } else {
                if n == n.Parent.Left {
                    n = n.Parent
                    t.rotateRight(n)
                }
                n.Parent.Color = BLACK
                n.Parent.Parent.Color = RED
                t.rotateLeft(n.Parent.Parent)
            }
        }
        if n == t.root { break }
    }
    t.root.Color = BLACK
}

func (t *Tree) Search(key int) bool {
    curr := t.root
    for curr != nil {
        if key == curr.Key { return true }
        if key < curr.Key { curr = curr.Left } else { curr = curr.Right }
    }
    return false
}

func (t *Tree) InOrder() []int {
    var result []int
    var traverse func(*Node)
    traverse = func(n *Node) {
        if n == nil { return }
        traverse(n.Left)
        result = append(result, n.Key)
        traverse(n.Right)
    }
    traverse(t.root)
    return result
}

func (t *Tree) Size() int { return t.size }
```

---

## 4. AVL vs Red-Black

| Property | AVL | Red-Black |
|----------|-----|-----------|
| Height bound | ≤ 1.44 log₂(n) | ≤ 2 log₂(n+1) |
| Balance after insert | Always strict | Looser (fewer rotations) |
| Rotations per insert | ≤ 2 | ≤ 2 |
| Rotations per delete | O(log n) | ≤ 3 |
| Lookup speed | Faster (smaller height) | Slightly slower |
| Insert/delete speed | Slower (more rebalancing) | Faster |
| Memory | 1 height field per node | 1 color bit per node |
| Use case | Read-heavy (more lookups) | Write-heavy (more inserts/deletes) |

**In practice:**
- `std::map` (C++), Java's `TreeMap`/`TreeSet`, and the Linux kernel (process scheduling, memory management) all use Red-Black trees
- Databases with read-heavy workloads sometimes use AVL

---

## 5. Go's Balanced Tree (btree package)

Go doesn't ship a balanced BST in the standard library. The most common choice is `github.com/google/btree` — a B-tree (N-way balanced tree, not binary), which offers the same sorted-map operations with better cache behavior.

```bash
go get github.com/google/btree
```

```go
package main

import (
    "fmt"
    "github.com/google/btree"
)

type IntItem int

func (i IntItem) Less(other btree.Item) bool {
    return i < other.(IntItem)
}

func main() {
    t := btree.New(32) // 32 = branching factor

    for _, v := range []int{5, 3, 7, 1, 4, 6, 8} {
        t.ReplaceOrInsert(IntItem(v))
    }

    t.Ascend(func(item btree.Item) bool {
        fmt.Print(item.(IntItem), " ")
        return true // return false to stop
    })
    // 1 3 4 5 6 7 8

    fmt.Println(t.Has(IntItem(4)))  // true
    fmt.Println(t.Min(), t.Max())   // 1 8
    t.Delete(IntItem(3))
    fmt.Println(t.Len())            // 6
}
```

### When to use a sorted map vs built-in map

| Requirement | Use |
|-------------|-----|
| O(1) average lookup | `map[K]V` |
| Sorted iteration | `btree` or sorted slice |
| Range queries (all keys between a and b) | `btree` |
| Predecessor/successor queries | `btree` |
| Concurrent read-heavy | `sync.Map` |

---

## Summary

- **BST worst case** O(n) for sorted input → balanced trees fix this with guaranteed O(log n)
- **AVL**: strict balance (height ≤ 1.44 log n); best for read-heavy workloads
- **Red-Black**: looser balance (height ≤ 2 log n); fewer rotations on writes; used in most standard library implementations
- **Rotations**: O(1) operations that preserve BST ordering while restoring balance
- Both trees use the same four rotation types: Left, Right, Left-Right (double), Right-Left (double)
- In Go, use `github.com/google/btree` for practical sorted-map needs

## Exercises

### Easy
1. Implement `Height(root *avl.Node) int` that returns the actual height of an AVL tree by traversal (not using the cached field). Verify it matches the cached height for all test trees.
2. Show that inserting keys `1, 2, 3, 4, 5` into an AVL tree produces a balanced tree by tracing each rotation. Draw the tree state after each insertion.
3. Add a `Min()` and `Max()` method to the AVL tree. Both should run in O(log n) time.

### Medium
4. Add a `RangeQuery(root *avl.Node, lo, hi int) []int` function that returns all keys in [lo, hi] in sorted order. This should run in O(k + log n) where k is the number of results, not O(n).
5. Implement `Rank(root *avl.Node, key int) int` — the 1-based rank of a key (how many keys are ≤ key). To do this efficiently, augment each AVL node with a `Size int` field (count of nodes in its subtree). Update it during insertions and rotations.
6. Implement `Select(root *avl.Node, k int) (int, bool)` — the k-th smallest key (1-indexed). This is the inverse of Rank. Both should be O(log n) using the augmented Size field.

### Hard
7. Implement the full red-black deletion algorithm. It's more complex than insertion: you need to handle 6 cases for "double-black" nodes after removing a black node. Test with a fuzz test that inserts and deletes random keys and verifies all 5 RBT invariants after each operation.
8. Build an **interval tree** on top of AVL: each node stores an interval [lo, hi]. Augment each node with the maximum `hi` value in its subtree. Implement `OverlappingQuery(lo, hi int) []Interval` that returns all stored intervals overlapping [lo, hi] in O(k log n) time.
