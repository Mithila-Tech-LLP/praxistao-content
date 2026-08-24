# Chapter 35: Trees and Binary Search Trees

Trees are hierarchical data structures — unlike linked lists (linear) or arrays (flat), trees branch. A **binary search tree (BST)** maintains the invariant that every left child is smaller and every right child is larger than the parent. This gives O(log n) search, insert, and delete on balanced trees — the foundation for ordered collections, range queries, and sorted iteration.

## Table of Contents

1. [Tree Terminology](#1-tree-terminology)
2. [Binary Tree Implementation](#2-binary-tree-implementation)
3. [Tree Traversals](#3-tree-traversals)
4. [Binary Search Tree (BST)](#4-binary-search-tree-bst)
5. [BST Operations](#5-bst-operations)
6. [Balanced Trees — Concepts](#6-balanced-trees--concepts)
7. [Classic Tree Problems](#7-classic-tree-problems)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Tree Terminology

```
            10          ← root (no parent)
           /  \
          5    15       ← internal nodes
         / \     \
        3   7    20     ← leaves (no children)
       /
      1
```

- **Root**: the top node (no parent)
- **Leaf**: a node with no children
- **Height**: longest path from root to a leaf (root alone = height 0)
- **Depth**: distance from root to a node (root depth = 0)
- **Level**: all nodes at the same depth
- **Subtree**: a node and all its descendants
- **Balanced**: height is O(log n) — left and right subtrees have similar heights
- **Full binary tree**: every node has 0 or 2 children
- **Complete binary tree**: all levels full except possibly the last (filled left-to-right)
- **Perfect binary tree**: all internal nodes have 2 children and all leaves at the same level

### Quick Check
> 1. What makes a node a "leaf"?
> 2. What is the height of a tree with just one node?
> 3. What is the maximum number of nodes at level k?

---

## 2. Binary Tree Implementation

```go
type TreeNode[T any] struct {
    Val   T
    Left  *TreeNode[T]
    Right *TreeNode[T]
}

func NewNode[T any](val T) *TreeNode[T] {
    return &TreeNode[T]{Val: val}
}

// Build the example tree manually:
//        10
//       /  \
//      5    15
//     / \     \
//    3   7    20

func buildExampleTree() *TreeNode[int] {
    root := NewNode(10)
    root.Left = NewNode(5)
    root.Right = NewNode(15)
    root.Left.Left = NewNode(3)
    root.Left.Right = NewNode(7)
    root.Right.Right = NewNode(20)
    return root
}

// Tree height (number of edges on longest root-to-leaf path):
func Height[T any](node *TreeNode[T]) int {
    if node == nil {
        return -1  // Height of empty tree = -1; height of single node = 0
    }
    left := Height(node.Left)
    right := Height(node.Right)
    if left > right {
        return left + 1
    }
    return right + 1
}

// Count nodes:
func Count[T any](node *TreeNode[T]) int {
    if node == nil {
        return 0
    }
    return 1 + Count(node.Left) + Count(node.Right)
}
```

### Quick Check
> 1. What is the height of a single-node tree?
> 2. Why does a nil node have height -1 in our convention?

---

## 3. Tree Traversals

Tree traversal means visiting every node exactly once. The three DFS traversals differ only in WHEN you visit the current node relative to its children:

```
Tree:      1
          / \
         2   3
        / \
       4   5

Inorder   (Left, Root, Right): 4 2 5 1 3  ← sorted for BST!
Preorder  (Root, Left, Right): 1 2 4 5 3  ← clone, serialize
Postorder (Left, Right, Root): 4 5 2 3 1  ← delete, evaluate
```

### Recursive Traversals
```go
// Inorder: Left → Root → Right
func Inorder[T any](node *TreeNode[T], visit func(T)) {
    if node == nil {
        return
    }
    Inorder(node.Left, visit)
    visit(node.Val)
    Inorder(node.Right, visit)
}

// Preorder: Root → Left → Right
func Preorder[T any](node *TreeNode[T], visit func(T)) {
    if node == nil {
        return
    }
    visit(node.Val)
    Preorder(node.Left, visit)
    Preorder(node.Right, visit)
}

// Postorder: Left → Right → Root
func Postorder[T any](node *TreeNode[T], visit func(T)) {
    if node == nil {
        return
    }
    Postorder(node.Left, visit)
    Postorder(node.Right, visit)
    visit(node.Val)
}
```

### Iterative Inorder (using explicit stack)
```go
// Iterative inorder — avoids potential stack overflow on deep trees:
func InorderIterative(root *TreeNode[int]) []int {
    var result []int
    stack := &Stack[*TreeNode[int]]{}
    cur := root

    for cur != nil || !stack.IsEmpty() {
        // Go as far left as possible:
        for cur != nil {
            stack.Push(cur)
            cur = cur.Left
        }
        // Process the node:
        cur, _ = stack.Pop()
        result = append(result, cur.Val)
        // Move to right subtree:
        cur = cur.Right
    }
    return result
}
```

### Level-Order (BFS) Traversal
```go
// LevelOrder visits nodes level by level (uses a queue):
func LevelOrder[T any](root *TreeNode[T]) [][]T {
    if root == nil {
        return nil
    }

    var result [][]T
    queue := []*TreeNode[T]{root}

    for len(queue) > 0 {
        levelSize := len(queue)
        level := make([]T, 0, levelSize)

        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            level = append(level, node.Val)
            if node.Left != nil {
                queue = append(queue, node.Left)
            }
            if node.Right != nil {
                queue = append(queue, node.Right)
            }
        }
        result = append(result, level)
    }
    return result
}
```

### Quick Check
> 1. For a BST, which traversal gives nodes in sorted order?
> 2. What data structure does level-order traversal use?
> 3. What is the difference between DFS and BFS tree traversal?

---

## 4. Binary Search Tree (BST)

**The BST invariant**: for every node N:
- All values in N's LEFT subtree < N.Val
- All values in N's RIGHT subtree > N.Val

```go
type BST[T any] struct {
    root *TreeNode[T]
    less func(a, b T) bool  // Comparison function
    size int
}

func NewBST[T any](less func(T, T) bool) *BST[T] {
    return &BST[T]{less: less}
}

// Search — O(log n) average, O(n) worst case (degenerate tree):
func (b *BST[T]) Search(val T) (*TreeNode[T], bool) {
    return search(b.root, val, b.less)
}

func search[T any](node *TreeNode[T], val T, less func(T, T) bool) (*TreeNode[T], bool) {
    if node == nil {
        return nil, false
    }
    if less(val, node.Val) {
        return search(node.Left, val, less)
    }
    if less(node.Val, val) {
        return search(node.Right, val, less)
    }
    return node, true  // Neither less: equal
}
```

### Quick Check
> 1. State the BST invariant.
> 2. What is the time complexity of search in a balanced BST?
> 3. What is the worst-case time complexity of BST operations?

---

## 5. BST Operations

### Insert — O(log n) average
```go
func (b *BST[T]) Insert(val T) {
    b.root = insert(b.root, val, b.less)
    b.size++
}

func insert[T any](node *TreeNode[T], val T, less func(T, T) bool) *TreeNode[T] {
    if node == nil {
        return &TreeNode[T]{Val: val}
    }
    if less(val, node.Val) {
        node.Left = insert(node.Left, val, less)
    } else if less(node.Val, val) {
        node.Right = insert(node.Right, val, less)
    }
    // Equal: duplicate — don't insert (or handle per requirements)
    return node
}
```

### Delete — O(log n) average
Delete has three cases:
1. **Node is a leaf**: just remove it
2. **Node has one child**: replace node with its child
3. **Node has two children**: replace with in-order successor (smallest in right subtree), then delete that successor

```go
func (b *BST[T]) Delete(val T) bool {
    var deleted bool
    b.root, deleted = deleteNode(b.root, val, b.less)
    if deleted {
        b.size--
    }
    return deleted
}

func deleteNode[T any](node *TreeNode[T], val T, less func(T, T) bool) (*TreeNode[T], bool) {
    if node == nil {
        return nil, false
    }

    var deleted bool
    if less(val, node.Val) {
        node.Left, deleted = deleteNode(node.Left, val, less)
    } else if less(node.Val, val) {
        node.Right, deleted = deleteNode(node.Right, val, less)
    } else {
        // Found the node to delete:
        deleted = true
        if node.Left == nil {
            return node.Right, deleted  // Case 1 & 2: no left child
        }
        if node.Right == nil {
            return node.Left, deleted   // Case 2: no right child
        }
        // Case 3: two children — find in-order successor (min of right subtree)
        successor := minNode(node.Right)
        node.Val = successor.Val
        node.Right, _ = deleteNode(node.Right, successor.Val, less)
    }
    return node, deleted
}

func minNode[T any](node *TreeNode[T]) *TreeNode[T] {
    for node.Left != nil {
        node = node.Left
    }
    return node
}
```

### Min, Max, Floor, Ceiling — O(log n)
```go
// Min returns the smallest value in the BST:
func (b *BST[T]) Min() (T, bool) {
    if b.root == nil {
        var zero T
        return zero, false
    }
    return minNode(b.root).Val, true
}

// Max returns the largest value:
func (b *BST[T]) Max() (T, bool) {
    if b.root == nil {
        var zero T
        return zero, false
    }
    node := b.root
    for node.Right != nil {
        node = node.Right
    }
    return node.Val, true
}

// Floor: largest value ≤ val
func (b *BST[T]) Floor(val T) (T, bool) {
    return floor(b.root, val, b.less)
}

func floor[T any](node *TreeNode[T], val T, less func(T, T) bool) (T, bool) {
    if node == nil {
        var zero T
        return zero, false
    }
    if !less(node.Val, val) && !less(val, node.Val) {
        return node.Val, true  // Exact match
    }
    if less(val, node.Val) {
        return floor(node.Left, val, less)  // val < node: floor must be in left
    }
    // val > node: floor might be node or in right subtree
    if f, ok := floor(node.Right, val, less); ok {
        return f, true
    }
    return node.Val, true
}
```

### In-order iteration (sorted output)
```go
func (b *BST[T]) InOrder(visit func(T)) {
    Inorder(b.root, visit)
}

// Collect all values in sorted order:
func (b *BST[T]) Sorted() []T {
    result := make([]T, 0, b.size)
    b.InOrder(func(v T) { result = append(result, v) })
    return result
}
```

### Quick Check
> 1. What are the three cases for BST delete?
> 2. What is the "in-order successor" of a node?
> 3. Why does inorder traversal of a BST give sorted output?

---

## 6. Balanced Trees — Concepts

A BST can degenerate to a linked list if elements are inserted in sorted order:
```
Insert 1,2,3,4,5 into a BST:
1 → 1
     \
      2
       \
        3  ← O(n) search, not O(log n)!
         \
          4
           \
            5
```

**Self-balancing BSTs** maintain O(log n) height automatically (covered in depth in Chapter 37):

**AVL Tree**: maintains the invariant that every node's left and right subtrees differ in height by at most 1. Rotations restore balance after insert/delete. Most strictly balanced — best for read-heavy workloads.

**Red-Black Tree**: nodes are colored red or black with rules that limit imbalance. Fewer rotations than AVL — used in Java's `TreeMap`, C++'s `std::map`, Linux kernel.

**B-Tree / B+ Tree**: N-way tree (not binary) for disk storage. Pages hold many keys. Minimizes disk I/O — used in databases and file systems.

**Go's standard library**: doesn't include a built-in balanced BST. For production use, the `btree` package by Google (`github.com/google/btree`) or a treap/skiplist are common choices.

**Checking BST balance:**
```go
func IsBalanced[T any](node *TreeNode[T]) bool {
    return checkHeight(node) != -1
}

// Returns height if balanced, -1 if unbalanced:
func checkHeight[T any](node *TreeNode[T]) int {
    if node == nil {
        return 0
    }
    left := checkHeight(node.Left)
    if left == -1 {
        return -1
    }
    right := checkHeight(node.Right)
    if right == -1 {
        return -1
    }
    diff := left - right
    if diff < -1 || diff > 1 {
        return -1  // Unbalanced at this node
    }
    if left > right {
        return left + 1
    }
    return right + 1
}
```

---

## 7. Classic Tree Problems

### Lowest Common Ancestor (LCA)
```go
// LCA finds the lowest node that is an ancestor of both p and q.
func LCA(root *TreeNode[int], p, q int) *TreeNode[int] {
    if root == nil {
        return nil
    }
    // For BST: use value comparison
    if p < root.Val && q < root.Val {
        return LCA(root.Left, p, q)
    }
    if p > root.Val && q > root.Val {
        return LCA(root.Right, p, q)
    }
    return root  // root is between p and q (or equal to one)
}
```

### Validate BST
```go
// IsValidBST checks if a tree satisfies the BST invariant.
func IsValidBST(root *TreeNode[int]) bool {
    return validate(root, nil, nil)
}

func validate(node *TreeNode[int], min, max *int) bool {
    if node == nil {
        return true
    }
    if min != nil && node.Val <= *min {
        return false
    }
    if max != nil && node.Val >= *max {
        return false
    }
    return validate(node.Left, min, &node.Val) &&
           validate(node.Right, &node.Val, max)
}
```

### Diameter of Binary Tree
```go
// Diameter finds the longest path between any two nodes (in edges).
func Diameter(root *TreeNode[int]) int {
    maxDiameter := 0
    var depth func(*TreeNode[int]) int
    depth = func(node *TreeNode[int]) int {
        if node == nil {
            return 0
        }
        left := depth(node.Left)
        right := depth(node.Right)
        if left+right > maxDiameter {
            maxDiameter = left + right  // Path through this node
        }
        if left > right {
            return left + 1
        }
        return right + 1
    }
    depth(root)
    return maxDiameter
}
```

### Serialize and Deserialize
```go
// Serialize converts tree to string: "10,5,3,nil,nil,7,nil,nil,15,nil,20,nil,nil"
func Serialize(root *TreeNode[int]) string {
    var sb strings.Builder
    var dfs func(*TreeNode[int])
    dfs = func(node *TreeNode[int]) {
        if node == nil {
            sb.WriteString("nil,")
            return
        }
        sb.WriteString(strconv.Itoa(node.Val))
        sb.WriteByte(',')
        dfs(node.Left)
        dfs(node.Right)
    }
    dfs(root)
    return sb.String()
}

func Deserialize(data string) *TreeNode[int] {
    tokens := strings.Split(strings.TrimRight(data, ","), ",")
    idx := 0
    var build func() *TreeNode[int]
    build = func() *TreeNode[int] {
        if idx >= len(tokens) || tokens[idx] == "nil" {
            idx++
            return nil
        }
        val, _ := strconv.Atoi(tokens[idx])
        idx++
        node := &TreeNode[int]{Val: val}
        node.Left = build()
        node.Right = build()
        return node
    }
    return build()
}
```

---

## Summary

- **BST invariant**: left < node < right at every node
- **Traversals**: Inorder (sorted), Preorder (clone/serialize), Postorder (delete/evaluate), Level-order (BFS)
- **Search/Insert**: O(log n) average, O(n) worst (degenerate); O(log n) guaranteed with balanced BST
- **Delete**: three cases — leaf, one child, two children (in-order successor)
- **Min/Max**: traverse left-/right-most path
- **Floor/Ceiling**: binary search using BST structure
- **Balanced BSTs**: AVL, Red-Black, B-Tree — all guarantee O(log n) height
- **Classic problems**: LCA (for BST: compare values), validate BST (pass min/max bounds), diameter (DFS returning height)

---

## Exercises

### Easy
1. Build a BST from a sorted array such that the resulting tree is balanced (height = O(log n)). `SortedArrayToBST(nums []int) *TreeNode[int]` — pick the middle element as root recursively.
2. Write `CountNodes[T](root *TreeNode[T]) int` iteratively (no recursion) using a queue.
3. Write `MirrorTree[T](root *TreeNode[T])` that swaps left and right children at every node (in-place mirror/flip).

### Medium
4. Kth smallest in BST: Write `KthSmallest(root *TreeNode[int], k int) int` using in-order traversal. First: recursive solution. Then: iterative solution using an explicit stack. Which is more memory efficient for large k? Why?
5. Range sum BST: Write `RangeSum(root *TreeNode[int], low, high int) int` that returns the sum of all values in the BST between `low` and `high` (inclusive). Prune subtrees using BST property to avoid unnecessary traversal. Verify it's faster than a full traversal on a tree with 10,000 nodes.
6. Flatten BST to sorted doubly linked list: Convert the BST to a sorted doubly linked list **in-place** (reuse the nodes, `Left` = prev, `Right` = next). Return the head of the doubly linked list. Do it with O(1) extra space (no stack/queue, only recursion counts as O(h) space where h is tree height).

### Hard
7. AVL tree: Implement an AVL self-balancing BST. Each node stores a height. After each insert/delete, check the balance factor (left height - right height). If it's ±2, perform rotations (single left, single right, double left-right, double right-left) to restore balance. Verify: insert 1 through 1000 in order — the tree should have height ≈ log₂(1000) ≈ 10, not 1000.
8. BST iterator: Implement `BSTIterator` that does in-order traversal lazily (on demand, not computing the whole traversal upfront). `Next() int` returns the next smallest value. `HasNext() bool`. Use O(h) memory (h = tree height). This models how a database might iterate over an index without loading all rows. Verify it works on trees of 10,000 nodes without precomputing the sorted array.
