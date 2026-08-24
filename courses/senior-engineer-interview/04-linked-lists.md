# Chapter 04: Linked Lists

Linked lists feel simple on paper but trip up candidates under pressure because the pointer manipulation requires careful thinking. The good news: there are exactly four techniques that solve almost every linked list problem. Master these and you can handle any linked list question confidently.

## Table of Contents

1. [Linked List Basics in Go](#1-linked-list-basics-in-go)
2. [The Four Techniques](#2-the-four-techniques)
3. [Fast & Slow Pointers](#3-fast--slow-pointers)
4. [List Reversal](#4-list-reversal)
5. [Merge Operations](#5-merge-operations)
6. [Two-List Techniques](#6-two-list-techniques)
7. [Classic Interview Problems](#7-classic-interview-problems)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Linked List Basics in Go

```go
// The standard linked list node definition used in most interview problems
type ListNode struct {
    Val  int
    Next *ListNode
}

// Helper to build a list from a slice (useful for testing)
func buildList(vals []int) *ListNode {
    dummy := &ListNode{}
    curr := dummy
    for _, v := range vals {
        curr.Next = &ListNode{Val: v}
        curr = curr.Next
    }
    return dummy.Next
}

// Helper to print a list (useful for debugging)
func printList(head *ListNode) {
    for head != nil {
        fmt.Printf("%d -> ", head.Val)
        head = head.Next
    }
    fmt.Println("nil")
}
```

### The Dummy Node Trick

Always use a dummy (sentinel) node when building or modifying a linked list. It eliminates the special case of handling the head node separately.

```go
// WITHOUT dummy: awkward edge cases for empty list or modifying head
func withoutDummy(head *ListNode, val int) *ListNode {
    if head == nil || head.Val == val {
        // special case... messy
    }
    // ...
}

// WITH dummy: uniform handling
func withDummy(head *ListNode, val int) *ListNode {
    dummy := &ListNode{Next: head} // dummy.Next always points to real head
    curr := dummy
    for curr.Next != nil {
        if curr.Next.Val == val {
            curr.Next = curr.Next.Next // remove node: clean, no special case
        } else {
            curr = curr.Next
        }
    }
    return dummy.Next // real head (might have changed)
}
```

---

## 2. The Four Techniques

| Technique | When to Use |
|---|---|
| Fast & Slow Pointers | Cycle detection, finding middle, nth from end |
| Reversal | Reverse a list or part of it |
| Merge | Combine two sorted lists, merge k lists |
| Two-List Traversal | Find intersection, compare lists |

---

## 3. Fast & Slow Pointers

### The Core Idea

Two pointers traverse the list at different speeds. The fast pointer moves 2 steps per iteration; the slow pointer moves 1. When fast reaches the end, slow is at the middle. If there is a cycle, fast will eventually lap slow.

### Finding the Middle of a Linked List

```go
func middleNode(head *ListNode) *ListNode {
    slow, fast := head, head

    // When fast reaches the end, slow is at the middle.
    // For even-length lists [1,2,3,4]: slow ends at 3 (second middle).
    for fast != nil && fast.Next != nil {
        slow = slow.Next       // move 1 step
        fast = fast.Next.Next  // move 2 steps
    }
    return slow
}
// Time: O(n), Space: O(1)
```

### Detecting a Cycle (Floyd's Algorithm)

```go
func hasCycle(head *ListNode) bool {
    slow, fast := head, head

    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next

        // If there is a cycle, fast will eventually lap slow and they meet.
        if slow == fast {
            return true
        }
    }
    // fast reached nil: no cycle
    return false
}
// Time: O(n), Space: O(1)
```

### Finding the Cycle Entry Point

```go
// If there is a cycle, find the node where the cycle begins.
func detectCycle(head *ListNode) *ListNode {
    slow, fast := head, head

    // Phase 1: detect if there is a cycle
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
        if slow == fast {
            break
        }
    }

    // No cycle detected
    if fast == nil || fast.Next == nil {
        return nil
    }

    // Phase 2: find the entry point of the cycle.
    // Mathematical proof: after meeting point, moving one pointer to head
    // and advancing both at speed 1 makes them meet at the cycle entry.
    slow = head
    for slow != fast {
        slow = slow.Next
        fast = fast.Next
    }
    return slow
}
// Time: O(n), Space: O(1)
// This is Floyd's Tortoise and Hare algorithm — interviewers love asking about it.
```

### Find Nth Node From End

```go
// Remove the nth node from the end of the list.
// Trick: advance fast by n steps first, then move both until fast.Next is nil.
func removeNthFromEnd(head *ListNode, n int) *ListNode {
    dummy := &ListNode{Next: head}
    slow, fast := dummy, dummy

    // Advance fast n+1 steps (using dummy, so one extra step to handle head removal)
    for i := 0; i <= n; i++ {
        fast = fast.Next
    }

    // Move both until fast is nil
    for fast != nil {
        slow = slow.Next
        fast = fast.Next
    }

    // slow.Next is the node to remove
    slow.Next = slow.Next.Next
    return dummy.Next
}
// Time: O(n), Space: O(1), single pass
```

---

## 4. List Reversal

### Reverse a Linked List (Iterative)

```go
// Standard iterative reversal — must know this cold
func reverseList(head *ListNode) *ListNode {
    var prev *ListNode // prev starts as nil (new tail's next)
    curr := head

    for curr != nil {
        next := curr.Next  // save next before we overwrite curr.Next
        curr.Next = prev   // reverse the link
        prev = curr        // advance prev
        curr = next        // advance curr
    }
    return prev // prev is now the new head
}
// Time: O(n), Space: O(1)

// Trace for [1->2->3->nil]:
// Step 1: next=2, 1->nil,  prev=1, curr=2
// Step 2: next=3, 2->1,    prev=2, curr=3
// Step 3: next=nil, 3->2,  prev=3, curr=nil
// Return: 3->2->1->nil ✓
```

### Reverse a Linked List (Recursive)

```go
// Recursive version — elegant but uses O(n) stack space
func reverseListRecursive(head *ListNode) *ListNode {
    // Base case: empty list or single node
    if head == nil || head.Next == nil {
        return head
    }

    // Recursively reverse the rest of the list
    newHead := reverseListRecursive(head.Next)

    // Make the next node point back to current node
    head.Next.Next = head
    head.Next = nil // current node becomes the new tail

    return newHead
}
```

### Reverse a Sublist (Between Positions m and n)

```go
func reverseBetween(head *ListNode, left, right int) *ListNode {
    if head == nil || left == right {
        return head
    }

    dummy := &ListNode{Next: head}
    prev := dummy

    // Move prev to the node just before position 'left'
    for i := 1; i < left; i++ {
        prev = prev.Next
    }

    curr := prev.Next
    // Reverse 'right - left' times using the "insert at front" technique
    for i := 0; i < right-left; i++ {
        next := curr.Next
        curr.Next = next.Next
        next.Next = prev.Next
        prev.Next = next
    }
    return dummy.Next
}
// Time: O(n), Space: O(1)
```

### Check if Linked List is Palindrome

```go
// Strategy: find middle, reverse second half, compare both halves
func isPalindrome(head *ListNode) bool {
    // Find middle using fast/slow
    slow, fast := head, head
    for fast != nil && fast.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    // Reverse the second half
    second := reverseList(slow)
    first := head

    // Compare both halves
    check := second // save to restore list later (optional but good practice)
    for check != nil {
        if first.Val != check.Val {
            return false
        }
        first = first.Next
        check = check.Next
    }
    return true
}
```

---

## 5. Merge Operations

### Merge Two Sorted Lists

```go
func mergeTwoLists(l1, l2 *ListNode) *ListNode {
    dummy := &ListNode{}
    curr := dummy

    for l1 != nil && l2 != nil {
        if l1.Val <= l2.Val {
            curr.Next = l1
            l1 = l1.Next
        } else {
            curr.Next = l2
            l2 = l2.Next
        }
        curr = curr.Next
    }

    // Attach remaining nodes (one of them is nil)
    if l1 != nil {
        curr.Next = l1
    } else {
        curr.Next = l2
    }
    return dummy.Next
}
// Time: O(m+n), Space: O(1)
```

### Merge K Sorted Lists — Min-Heap Approach

```go
import "container/heap"

// A min-heap of ListNode pointers
type NodeHeap []*ListNode
func (h NodeHeap) Len() int            { return len(h) }
func (h NodeHeap) Less(i, j int) bool  { return h[i].Val < h[j].Val }
func (h NodeHeap) Swap(i, j int)       { h[i], h[j] = h[j], h[i] }
func (h *NodeHeap) Push(x interface{}) { *h = append(*h, x.(*ListNode)) }
func (h *NodeHeap) Pop() interface{} {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func mergeKLists(lists []*ListNode) *ListNode {
    h := &NodeHeap{}
    heap.Init(h)

    // Push the head of each list into the heap
    for _, list := range lists {
        if list != nil {
            heap.Push(h, list)
        }
    }

    dummy := &ListNode{}
    curr := dummy

    for h.Len() > 0 {
        // Pop the smallest node
        node := heap.Pop(h).(*ListNode)
        curr.Next = node
        curr = curr.Next

        // Push the next node from the same list
        if node.Next != nil {
            heap.Push(h, node.Next)
        }
    }
    return dummy.Next
}
// Time: O(N log k) where N = total nodes, k = number of lists
// Space: O(k) for the heap
```

---

## 6. Two-List Techniques

### Find Intersection of Two Linked Lists

```go
// Two lists intersect at a node (same pointer, not just same value).
// Key insight: if you traverse A then B, and B then A, both pointers
// cover the same total length and meet at the intersection.
func getIntersectionNode(headA, headB *ListNode) *ListNode {
    a, b := headA, headB

    // If they don't intersect, both will reach nil at the same time
    // (after traversing each other's list), so we return nil correctly.
    for a != b {
        if a != nil {
            a = a.Next
        } else {
            a = headB // switch to the other list
        }
        if b != nil {
            b = b.Next
        } else {
            b = headA // switch to the other list
        }
    }
    return a // either the intersection node or nil
}
// Time: O(m+n), Space: O(1)
```

---

## 7. Classic Interview Problems

### Sort a Linked List in O(n log n) — Merge Sort

```go
// Merge sort is ideal for linked lists — no random access needed
func sortList(head *ListNode) *ListNode {
    if head == nil || head.Next == nil {
        return head
    }

    // Find middle and split into two halves
    mid := getMid(head)
    right := mid.Next
    mid.Next = nil // cut the list

    left := sortList(head)
    right = sortList(right)
    return mergeTwoLists(left, right)
}

func getMid(head *ListNode) *ListNode {
    slow, fast := head, head
    for fast.Next != nil && fast.Next.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }
    return slow // slow is at the middle
}
// Time: O(n log n), Space: O(log n) for recursion stack
```

### Copy List with Random Pointer

```go
type Node struct {
    Val    int
    Next   *Node
    Random *Node
}

func copyRandomList(head *Node) *Node {
    if head == nil { return nil }

    // Map original nodes to their copies
    nodeMap := make(map[*Node]*Node)

    // First pass: create all copy nodes
    curr := head
    for curr != nil {
        nodeMap[curr] = &Node{Val: curr.Val}
        curr = curr.Next
    }

    // Second pass: connect Next and Random pointers
    curr = head
    for curr != nil {
        if curr.Next != nil {
            nodeMap[curr].Next = nodeMap[curr.Next]
        }
        if curr.Random != nil {
            nodeMap[curr].Random = nodeMap[curr.Random]
        }
        curr = curr.Next
    }
    return nodeMap[head]
}
// Time: O(n), Space: O(n)
```

### Reorder List

**Problem:** Given [1, 2, 3, 4, 5], reorder to [1, 5, 2, 4, 3]

```go
func reorderList(head *ListNode) {
    if head == nil || head.Next == nil { return }

    // Step 1: Find the middle
    slow, fast := head, head
    for fast.Next != nil && fast.Next.Next != nil {
        slow = slow.Next
        fast = fast.Next.Next
    }

    // Step 2: Reverse the second half
    second := reverseList(slow.Next)
    slow.Next = nil // split the list

    // Step 3: Merge two halves by interleaving
    first := head
    for second != nil {
        tmp1 := first.Next
        tmp2 := second.Next
        first.Next = second
        second.Next = tmp1
        first = tmp1
        second = tmp2
    }
}
// Time: O(n), Space: O(1)
```

---

## Summary

- Use a **dummy node** to simplify edge cases when building or modifying lists.
- **Fast/slow pointers:** find middle (fast 2x), detect cycle (they meet), nth from end (offset by n first).
- **Reversal:** save `next`, flip `curr.Next = prev`, advance both. Three variables, clean iteration.
- **Merge:** use a dummy node + current pointer. Always prefer iterative over recursive for O(1) space.
- **Intersection:** traverse A→B and B→A simultaneously — they equalize at the intersection point.
- For O(n log n) sort on a linked list: use merge sort (no random access needed).

---

## Interview Questions & Model Answers

**Q: Why is reversing a linked list O(1) space?**
"We only need three pointers: prev, curr, and next. No extra data structure is needed regardless of list length. Compare this to recursion which uses O(n) stack space."

**Q: What happens if you forget to set `mid.Next = nil` in merge sort?**
"The list is not split — you have infinite recursion. The second half still points back into the first half, so you never reach the base case."

**Q: How does Floyd's cycle detection algorithm work mathematically?**
"When slow has traveled distance d, fast has traveled 2d. If there is a cycle of length C, fast laps slow every C steps. So they meet after at most n steps where n is the list length. For the cycle entry: from the meeting point, the distance to the entry equals the distance from head to the entry — this is the mathematical invariant that makes phase 2 work."

---

## Exercises

### Easy
1. Reverse a linked list iteratively and recursively. Verify both solutions match.
2. Merge two sorted linked lists. Test with one empty list, two empty lists.
3. Find the middle node of a linked list. Test with even and odd lengths.

### Medium
4. Given a linked list with duplicates (sorted), remove all elements that appear more than once. Example: 1→1→2→3→3 → 2.
5. Rotate a linked list to the right by k places. Handle k > list length.
6. Partition a linked list around value x: all nodes with val < x come before nodes with val >= x.

### Hard
7. Reverse nodes in groups of k: given [1,2,3,4,5] and k=2, return [2,1,4,3,5]. If remaining nodes are fewer than k, leave them as-is.
8. Given a linked list where each node has a `child` pointer (like a doubly linked list but with levels), flatten the list.
9. Implement an LRU cache using a doubly linked list and a hash map. Support `get` and `put` in O(1).
