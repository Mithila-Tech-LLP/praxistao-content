# Chapter 06: Trees — DFS, BFS, BST & LCA

Trees are the most tested data structure in senior interviews after arrays. They combine recursion, iteration, and problem decomposition into a single topic. Once you internalize the recursive tree mindset — "solve for the root, trust the recursion handles the rest" — tree problems become significantly more approachable.

## Table of Contents

1. [Tree Fundamentals in Go](#1-tree-fundamentals-in-go)
2. [Tree Traversals — DFS](#2-tree-traversals--dfs)
3. [Level-Order Traversal — BFS](#3-level-order-traversal--bfs)
4. [Binary Search Trees](#4-binary-search-trees)
5. [Lowest Common Ancestor](#5-lowest-common-ancestor)
6. [Tree Construction Problems](#6-tree-construction-problems)
7. [Classic Tree Problems](#7-classic-tree-problems)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Tree Fundamentals in Go

```go
type TreeNode struct {
    Val   int
    Left  *TreeNode
    Right *TreeNode
}

// Build a tree from level-order array (nil = missing node)
// [3, 9, 20, nil, nil, 15, 7] builds:
//       3
//      / \
//     9  20
//        / \
//       15   7
func buildTree(vals []interface{}) *TreeNode {
    if len(vals) == 0 || vals[0] == nil {
        return nil
    }
    root := &TreeNode{Val: vals[0].(int)}
    queue := []*TreeNode{root}
    i := 1
    for len(queue) > 0 && i < len(vals) {
        node := queue[0]
        queue = queue[1:]
        if i < len(vals) && vals[i] != nil {
            node.Left = &TreeNode{Val: vals[i].(int)}
            queue = append(queue, node.Left)
        }
        i++
        if i < len(vals) && vals[i] != nil {
            node.Right = &TreeNode{Val: vals[i].(int)}
            queue = append(queue, node.Right)
        }
        i++
    }
    return root
}
```

### The Recursive Tree Mindset

Before writing any tree solution, ask:
1. What does this function return for a single node?
2. What does it return for a null node (base case)?
3. How do I combine results from left and right subtrees?

If you can answer these three questions, you can write the recursive solution.

---

## 2. Tree Traversals — DFS

Three DFS traversal orders. Know all three without thinking.

```go
// PREORDER: root → left → right (good for copying/serializing a tree)
func preorder(root *TreeNode, result *[]int) {
    if root == nil { return }
    *result = append(*result, root.Val)  // process root first
    preorder(root.Left, result)
    preorder(root.Right, result)
}

// INORDER: left → root → right (gives sorted order for BST)
func inorder(root *TreeNode, result *[]int) {
    if root == nil { return }
    inorder(root.Left, result)
    *result = append(*result, root.Val)  // process root between children
    inorder(root.Right, result)
}

// POSTORDER: left → right → root (good for deleting a tree, computing subtree sizes)
func postorder(root *TreeNode, result *[]int) {
    if root == nil { return }
    postorder(root.Left, result)
    postorder(root.Right, result)
    *result = append(*result, root.Val)  // process root last
}
```

### Iterative DFS (Interviewers Often Ask For This)

```go
// Iterative inorder using an explicit stack
func inorderIterative(root *TreeNode) []int {
    result := []int{}
    stack := []*TreeNode{}
    curr := root

    for curr != nil || len(stack) > 0 {
        // Go as far left as possible, pushing nodes onto the stack
        for curr != nil {
            stack = append(stack, curr)
            curr = curr.Left
        }
        // Pop the leftmost node, process it, then explore its right subtree
        curr = stack[len(stack)-1]
        stack = stack[:len(stack)-1]
        result = append(result, curr.Val)
        curr = curr.Right
    }
    return result
}
```

### Tree Height / Maximum Depth

```go
func maxDepth(root *TreeNode) int {
    if root == nil { return 0 }
    // The depth of this node = 1 + max depth of its subtrees
    leftDepth := maxDepth(root.Left)
    rightDepth := maxDepth(root.Right)
    return 1 + max(leftDepth, rightDepth)
}
// Time: O(n) — visit every node once
// Space: O(h) where h is tree height (O(log n) balanced, O(n) skewed)
```

### Check if Balanced

```go
// A tree is balanced if every node's left and right subtrees differ in height by at most 1.
func isBalanced(root *TreeNode) bool {
    return checkHeight(root) != -1
}

// Returns height if subtree is balanced, -1 if not
func checkHeight(root *TreeNode) int {
    if root == nil { return 0 }

    left := checkHeight(root.Left)
    if left == -1 { return -1 } // left subtree is unbalanced

    right := checkHeight(root.Right)
    if right == -1 { return -1 } // right subtree is unbalanced

    if abs(left-right) > 1 { return -1 } // current node is unbalanced

    return 1 + max(left, right) // return height
}
// This is O(n) — it computes height and balance check in a single pass.
// The naive O(n²) approach calls maxDepth from every node separately.
```

### Diameter of a Binary Tree

```go
// Diameter = longest path between any two nodes (may not pass through root)
// The diameter through a node = leftHeight + rightHeight
func diameterOfBinaryTree(root *TreeNode) int {
    diameter := 0
    
    var height func(*TreeNode) int
    height = func(node *TreeNode) int {
        if node == nil { return 0 }
        left := height(node.Left)
        right := height(node.Right)
        diameter = max(diameter, left+right) // update global diameter
        return 1 + max(left, right)
    }
    height(root)
    return diameter
}
```

---

## 3. Level-Order Traversal — BFS

BFS processes nodes level by level. Essential for problems involving levels, level-wise operations, or shortest paths in trees.

```go
func levelOrder(root *TreeNode) [][]int {
    if root == nil { return nil }
    
    result := [][]int{}
    queue := []*TreeNode{root}

    for len(queue) > 0 {
        levelSize := len(queue) // number of nodes at this level
        level := make([]int, levelSize)

        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            level[i] = node.Val

            if node.Left != nil  { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
        result = append(result, level)
    }
    return result
}
// Time: O(n), Space: O(w) where w is maximum width of tree
```

### Right Side View

```go
// Return the value visible from the right side at each level
func rightSideView(root *TreeNode) []int {
    if root == nil { return nil }
    result := []int{}
    queue := []*TreeNode{root}

    for len(queue) > 0 {
        levelSize := len(queue)
        for i := 0; i < levelSize; i++ {
            node := queue[0]
            queue = queue[1:]
            if i == levelSize-1 { // last node at this level = rightmost
                result = append(result, node.Val)
            }
            if node.Left != nil  { queue = append(queue, node.Left) }
            if node.Right != nil { queue = append(queue, node.Right) }
        }
    }
    return result
}
```

---

## 4. Binary Search Trees

In a BST: left subtree values < root < right subtree values. Inorder traversal gives sorted order.

### Validate a BST

```go
// Common mistake: only checking parent-child relationship.
// A node's value must be within a valid range [min, max].
func isValidBST(root *TreeNode) bool {
    return validate(root, nil, nil)
}

func validate(node *TreeNode, min, max *int) bool {
    if node == nil { return true }
    if min != nil && node.Val <= *min { return false }
    if max != nil && node.Val >= *max { return false }
    // Left subtree: all values must be < node.Val
    // Right subtree: all values must be > node.Val
    return validate(node.Left, min, &node.Val) &&
           validate(node.Right, &node.Val, max)
}
```

### Insert into a BST

```go
func insertIntoBST(root *TreeNode, val int) *TreeNode {
    if root == nil {
        return &TreeNode{Val: val}
    }
    if val < root.Val {
        root.Left = insertIntoBST(root.Left, val)
    } else {
        root.Right = insertIntoBST(root.Right, val)
    }
    return root
}
```

### Delete from a BST

```go
func deleteNode(root *TreeNode, key int) *TreeNode {
    if root == nil { return nil }

    if key < root.Val {
        root.Left = deleteNode(root.Left, key)
    } else if key > root.Val {
        root.Right = deleteNode(root.Right, key)
    } else {
        // Found the node to delete
        if root.Left == nil  { return root.Right } // no left child
        if root.Right == nil { return root.Left }  // no right child

        // Both children exist: replace with inorder successor (min of right subtree)
        successor := root.Right
        for successor.Left != nil { successor = successor.Left }
        root.Val = successor.Val
        root.Right = deleteNode(root.Right, successor.Val)
    }
    return root
}
```

### Kth Smallest in BST

```go
// Inorder traversal gives sorted order — the kth visited node is the kth smallest
func kthSmallest(root *TreeNode, k int) int {
    count := 0
    result := 0

    var inorder func(*TreeNode)
    inorder = func(node *TreeNode) {
        if node == nil || count >= k { return }
        inorder(node.Left)
        count++
        if count == k {
            result = node.Val
            return
        }
        inorder(node.Right)
    }
    inorder(root)
    return result
}
```

---

## 5. Lowest Common Ancestor

LCA is one of the most commonly asked tree problems at senior level. Know two approaches.

### LCA of a Binary Tree (Not BST)

```go
// The elegant recursive approach:
// If root equals p or q, return root.
// Recurse both sides. If both return non-nil, root IS the LCA.
// If only one side returns non-nil, that side contains both p and q.
func lowestCommonAncestor(root, p, q *TreeNode) *TreeNode {
    if root == nil || root == p || root == q {
        return root
    }

    left := lowestCommonAncestor(root.Left, p, q)
    right := lowestCommonAncestor(root.Right, p, q)

    if left != nil && right != nil {
        return root // p and q are in different subtrees — root is LCA
    }
    if left != nil {
        return left // both p and q are in the left subtree
    }
    return right // both p and q are in the right subtree
}
// Time: O(n), Space: O(h)
```

### LCA of a BST (More Efficient)

```go
// In a BST, we can use the BST property to guide our search
func lcaBST(root, p, q *TreeNode) *TreeNode {
    for root != nil {
        if p.Val < root.Val && q.Val < root.Val {
            root = root.Left  // both in left subtree
        } else if p.Val > root.Val && q.Val > root.Val {
            root = root.Right // both in right subtree
        } else {
            return root // they split here — root is LCA
        }
    }
    return nil
}
// Time: O(h) — O(log n) for balanced BST, O(n) worst case
// Space: O(1) iterative
```

---

## 6. Tree Construction Problems

### Build Tree from Preorder and Inorder Traversal

```go
// preorder[0] is always the root.
// Find root's position in inorder — everything left is left subtree, right is right subtree.
func buildTreeFromPreIn(preorder []int, inorder []int) *TreeNode {
    if len(preorder) == 0 { return nil }

    rootVal := preorder[0]
    root := &TreeNode{Val: rootVal}

    // Find root in inorder
    mid := 0
    for mid < len(inorder) && inorder[mid] != rootVal { mid++ }

    // Left subtree has 'mid' nodes
    root.Left = buildTreeFromPreIn(preorder[1:mid+1], inorder[:mid])
    root.Right = buildTreeFromPreIn(preorder[mid+1:], inorder[mid+1:])
    return root
}
// Time: O(n²) naive, O(n) with a hashmap for inorder index lookup
```

### Serialize and Deserialize a Binary Tree

```go
// Serialize: preorder traversal, use "#" for nil nodes
// "3,9,#,#,20,15,#,#,7,#,#"
func serialize(root *TreeNode) string {
    if root == nil { return "#" }
    return fmt.Sprintf("%d,%s,%s", root.Val, serialize(root.Left), serialize(root.Right))
}

func deserialize(data string) *TreeNode {
    parts := strings.Split(data, ",")
    idx := 0
    var build func() *TreeNode
    build = func() *TreeNode {
        if idx >= len(parts) || parts[idx] == "#" {
            idx++
            return nil
        }
        val, _ := strconv.Atoi(parts[idx])
        idx++
        node := &TreeNode{Val: val}
        node.Left = build()
        node.Right = build()
        return node
    }
    return build()
}
```

---

## 7. Classic Tree Problems

### Maximum Path Sum

```go
// Path can go from any node to any node (not necessarily through root).
// Through each node: leftGain + node.Val + rightGain
// Each subtree returns the max "arm" it can contribute to its parent.
func maxPathSum(root *TreeNode) int {
    maxSum := root.Val

    var dfs func(*TreeNode) int
    dfs = func(node *TreeNode) int {
        if node == nil { return 0 }
        leftGain := max(0, dfs(node.Left))   // ignore negative contributions
        rightGain := max(0, dfs(node.Right))
        maxSum = max(maxSum, node.Val+leftGain+rightGain) // path through this node
        return node.Val + max(leftGain, rightGain)        // best arm for parent
    }
    dfs(root)
    return maxSum
}
```

### Count Good Nodes in a Binary Tree

```go
// A node is "good" if no node on its path from root is greater than it.
func goodNodes(root *TreeNode) int {
    var dfs func(*TreeNode, int) int
    dfs = func(node *TreeNode, maxSoFar int) int {
        if node == nil { return 0 }
        good := 0
        if node.Val >= maxSoFar { good = 1 }
        maxSoFar = max(maxSoFar, node.Val)
        return good + dfs(node.Left, maxSoFar) + dfs(node.Right, maxSoFar)
    }
    return dfs(root, root.Val)
}
```

---

## Summary

- **Preorder:** root first (copying trees). **Inorder:** sorted for BST. **Postorder:** children first (deletion, size).
- For balanced check, diameter, max path sum: compute in a single DFS pass returning the metric you need.
- **BFS:** Use a queue. Process `levelSize` nodes per iteration to separate levels.
- **BST:** use the ordering property to guide search. Validate with [min, max] range, not just parent comparison.
- **LCA:** elegant recursion — if both subtrees return non-nil, root is the LCA.
- **Recursive mindset:** "What does this return for nil? For a leaf? How do I combine left and right results?"

---

## Exercises

### Easy
1. Check if two binary trees are identical (same structure and values).
2. Find all root-to-leaf paths and return them as strings like "1->2->5".
3. Invert/mirror a binary tree.

### Medium
4. Find the minimum depth of a binary tree (shortest path from root to any leaf).
5. Given a binary tree and a target sum, find all root-to-leaf paths that sum to target.
6. Check if a binary tree is symmetric (mirror image of itself).

### Hard
7. Given a binary tree, return the vertical order traversal. Nodes at the same column are sorted by row, then by value.
8. Binary Tree Cameras: place the minimum number of cameras to monitor every node. A camera at a node monitors the node, its parent, and its children.
9. Recover a BST where exactly two nodes were swapped. Fix the BST in O(1) space (no recursion — use Morris traversal).
