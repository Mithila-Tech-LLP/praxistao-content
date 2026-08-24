# Chapter 05: Stacks, Queues & Monotonic Structures

Stacks and queues are simple data structures that enable surprisingly powerful algorithms. The monotonic stack is the underrated pattern that many senior candidates miss — it solves "next greater element" style problems in O(n) that feel like they need O(n²).

## Table of Contents

1. [Stack Basics in Go](#1-stack-basics-in-go)
2. [Queue Basics in Go](#2-queue-basics-in-go)
3. [The Monotonic Stack](#3-the-monotonic-stack)
4. [The Deque (Double-Ended Queue)](#4-the-deque)
5. [Classic Stack Problems](#5-classic-stack-problems)
6. [Classic Queue Problems](#6-classic-queue-problems)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Stack Basics in Go

A stack is LIFO — Last In, First Out. In Go, implement it as a slice. No `Stack` type in the standard library.

```go
// Stack implementation using a slice
type Stack struct {
    items []int
}

func (s *Stack) Push(val int) {
    s.items = append(s.items, val)
}

func (s *Stack) Pop() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    n := len(s.items)
    val := s.items[n-1]
    s.items = s.items[:n-1]
    return val, true
}

func (s *Stack) Peek() (int, bool) {
    if len(s.items) == 0 {
        return 0, false
    }
    return s.items[len(s.items)-1], true
}

func (s *Stack) IsEmpty() bool { return len(s.items) == 0 }
func (s *Stack) Size() int     { return len(s.items) }
```

In interviews, most people just use a slice directly rather than wrapping it:

```go
// Idiomatic Go for stack operations in interview context
stack := []int{}
stack = append(stack, 5)           // push
top := stack[len(stack)-1]         // peek
stack = stack[:len(stack)-1]       // pop
isEmpty := len(stack) == 0         // check empty
```

---

## 2. Queue Basics in Go

A queue is FIFO — First In, First Out. Go has no standard queue, but `container/list` or a slice works. For BFS, a slice is most common in interviews.

```go
// Simple queue using slice (enqueue = append, dequeue = shift)
queue := []int{}
queue = append(queue, 1)     // enqueue
front := queue[0]            // peek at front
queue = queue[1:]            // dequeue (O(n) — shifts elements)

// For large queues, use a circular buffer or linked list to get O(1) dequeue.
// In interviews, the slice approach is acceptable — just mention the tradeoff.
```

### Thread-Safe Queue for Production (Go)

```go
// Buffered channel as a bounded queue — idiomatic Go
queue := make(chan int, 100)
queue <- 42        // enqueue (blocks if full)
val := <-queue     // dequeue (blocks if empty)
val, ok := <-queue // non-blocking dequeue (ok=false if empty and closed)
```

---

## 3. The Monotonic Stack

### The Core Idea

A monotonic stack is a stack that maintains a monotonically increasing or decreasing order. When you push a new element, you first pop all elements that violate the monotonic property. The popped elements are the ones for which the new element is the "answer" (e.g., the next greater element).

**When to use:** Any problem involving "next/previous greater/smaller element," "largest rectangle," "trapped water," or "stock span" patterns.

### Problem: Next Greater Element

**Input:** [2, 1, 2, 4, 3]
**Output:** [4, 2, 4, -1, -1] (for each element, the next element that is greater, or -1)

```go
func nextGreaterElement(nums []int) []int {
    n := len(nums)
    result := make([]int, n)
    for i := range result { result[i] = -1 } // default: no greater element

    // Stack stores indices of elements waiting for their next greater element.
    // The stack is monotonically decreasing (largest on bottom, smallest on top).
    stack := []int{}

    for i := 0; i < n; i++ {
        // While the current element is greater than the top of the stack,
        // the current element IS the "next greater element" for the top.
        for len(stack) > 0 && nums[i] > nums[stack[len(stack)-1]] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1] // pop
            result[idx] = nums[i]        // nums[i] is the answer for idx
        }
        stack = append(stack, i) // push current index
    }
    // Remaining elements in stack have no next greater element (result stays -1)
    return result
}
// Time: O(n) — each element is pushed and popped at most once
// Space: O(n) for the stack
```

### Problem: Largest Rectangle in Histogram

This is one of the hardest monotonic stack problems. It appears frequently at FAANG.

```go
// For each bar, find the largest rectangle it can be part of.
// Key insight: for bar i to be part of a rectangle of height h[i],
// it must be flanked by bars of at least the same height.
func largestRectangleArea(heights []int) int {
    // Monotonically increasing stack: stores indices of bars in increasing order.
    // When we find a bar shorter than the top, we calculate area for the top bar.
    stack := []int{}
    maxArea := 0

    // Process all bars + a sentinel bar of height 0 at the end
    // to flush remaining bars from the stack.
    for i := 0; i <= len(heights); i++ {
        h := 0
        if i < len(heights) { h = heights[i] }

        for len(stack) > 0 && h < heights[stack[len(stack)-1]] {
            top := stack[len(stack)-1]
            stack = stack[:len(stack)-1] // pop

            height := heights[top]
            width := i // width extends from left boundary to current i
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }
            maxArea = max(maxArea, height*width)
        }
        stack = append(stack, i)
    }
    return maxArea
}
// Time: O(n), Space: O(n)
```

### Problem: Trapping Rain Water

```go
// For each position, trapped water = min(maxLeft, maxRight) - height[i]
// Approach 1: Two arrays (O(n) space)
// Approach 2: Two pointers (O(1) space) — the elegant one

func trap(height []int) int {
    left, right := 0, len(height)-1
    leftMax, rightMax := 0, 0
    water := 0

    for left < right {
        if height[left] < height[right] {
            // Process left side: guaranteed that maxRight >= height[right] >= height[left]
            // So the limiting factor is leftMax
            if height[left] >= leftMax {
                leftMax = height[left]
            } else {
                water += leftMax - height[left]
            }
            left++
        } else {
            // Process right side
            if height[right] >= rightMax {
                rightMax = height[right]
            } else {
                water += rightMax - height[right]
            }
            right--
        }
    }
    return water
}
// Time: O(n), Space: O(1)
```

### Problem: Daily Temperatures

**Input:** [73, 74, 75, 71, 69, 72, 76, 73]
**Output:** [1, 1, 4, 2, 1, 1, 0, 0]
(For each day, how many days until a warmer temperature)

```go
func dailyTemperatures(temps []int) []int {
    n := len(temps)
    result := make([]int, n)
    stack := []int{} // indices of days waiting for a warmer day

    for i := 0; i < n; i++ {
        for len(stack) > 0 && temps[i] > temps[stack[len(stack)-1]] {
            idx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]
            result[idx] = i - idx // days waited = current index - waiting index
        }
        stack = append(stack, i)
    }
    return result
}
// Time: O(n), Space: O(n)
```

---

## 4. The Deque

A deque (double-ended queue) supports push/pop from both ends in O(1). Use it for sliding window maximum.

### Problem: Sliding Window Maximum

**Input:** nums = [1,3,-1,-3,5,3,6,7], k = 3
**Output:** [3,3,5,5,6,7]
(Maximum of each window of size k)

```go
func maxSlidingWindow(nums []int, k int) []int {
    // Deque stores indices. The front is always the maximum for the current window.
    // The deque is monotonically decreasing (largest at front).
    deque := []int{}
    result := []int{}

    for i := 0; i < len(nums); i++ {
        // Remove indices that are outside the current window
        for len(deque) > 0 && deque[0] < i-k+1 {
            deque = deque[1:] // pop front
        }

        // Remove indices from the back that are smaller than current element
        // (they can never be the maximum for any future window)
        for len(deque) > 0 && nums[i] > nums[deque[len(deque)-1]] {
            deque = deque[:len(deque)-1] // pop back
        }

        deque = append(deque, i) // push current index to back

        // Window is full — record the maximum (front of deque)
        if i >= k-1 {
            result = append(result, nums[deque[0]])
        }
    }
    return result
}
// Time: O(n) — each element is added and removed at most once
// Space: O(k) for the deque
```

---

## 5. Classic Stack Problems

### Valid Parentheses

```go
func isValid(s string) bool {
    stack := []rune{}
    pairs := map[rune]rune{')': '(', ']': '[', '}': '{'}

    for _, ch := range s {
        if ch == '(' || ch == '[' || ch == '{' {
            stack = append(stack, ch) // push opening bracket
        } else {
            // Closing bracket: stack must be non-empty and top must match
            if len(stack) == 0 || stack[len(stack)-1] != pairs[ch] {
                return false
            }
            stack = stack[:len(stack)-1] // pop
        }
    }
    return len(stack) == 0 // valid only if all brackets matched
}
```

### Min Stack — O(1) Minimum Retrieval

```go
// Design a stack that supports push, pop, top, and getMin in O(1).
type MinStack struct {
    stack    []int
    minStack []int // parallel stack tracking minimums
}

func (s *MinStack) Push(val int) {
    s.stack = append(s.stack, val)
    if len(s.minStack) == 0 || val <= s.minStack[len(s.minStack)-1] {
        s.minStack = append(s.minStack, val) // push new min
    } else {
        s.minStack = append(s.minStack, s.minStack[len(s.minStack)-1]) // carry over current min
    }
}

func (s *MinStack) Pop() {
    s.stack = s.stack[:len(s.stack)-1]
    s.minStack = s.minStack[:len(s.minStack)-1]
}

func (s *MinStack) Top() int      { return s.stack[len(s.stack)-1] }
func (s *MinStack) GetMin() int   { return s.minStack[len(s.minStack)-1] }
```

### Evaluate Reverse Polish Notation

```go
// Tokens: ["2","1","+","3","*"] → ((2+1)*3) = 9
func evalRPN(tokens []string) int {
    stack := []int{}
    ops := map[string]func(a, b int) int{
        "+": func(a, b int) int { return a + b },
        "-": func(a, b int) int { return a - b },
        "*": func(a, b int) int { return a * b },
        "/": func(a, b int) int { return a / b },
    }

    for _, token := range tokens {
        if op, ok := ops[token]; ok {
            b := stack[len(stack)-1]; stack = stack[:len(stack)-1]
            a := stack[len(stack)-1]; stack = stack[:len(stack)-1]
            stack = append(stack, op(a, b))
        } else {
            n, _ := strconv.Atoi(token)
            stack = append(stack, n)
        }
    }
    return stack[0]
}
```

---

## 6. Classic Queue Problems

### Implement Queue Using Two Stacks

```go
// push is O(1). pop/peek is amortized O(1).
type MyQueue struct {
    inbox  []int // for push
    outbox []int // for pop/peek
}

func (q *MyQueue) Push(val int) {
    q.inbox = append(q.inbox, val)
}

func (q *MyQueue) transfer() {
    // Transfer all from inbox to outbox (reverses order, making oldest element on top)
    if len(q.outbox) == 0 {
        for len(q.inbox) > 0 {
            n := len(q.inbox)
            q.outbox = append(q.outbox, q.inbox[n-1])
            q.inbox = q.inbox[:n-1]
        }
    }
}

func (q *MyQueue) Pop() int {
    q.transfer()
    val := q.outbox[len(q.outbox)-1]
    q.outbox = q.outbox[:len(q.outbox)-1]
    return val
}

func (q *MyQueue) Peek() int {
    q.transfer()
    return q.outbox[len(q.outbox)-1]
}

func (q *MyQueue) Empty() bool {
    return len(q.inbox) == 0 && len(q.outbox) == 0
}
```

---

## Summary

- In Go, implement stacks with a slice. Use `append` to push, `slice[:n-1]` to pop.
- **Monotonic stack** = pop elements that violate the monotonic property when pushing a new element. Each element is pushed/popped at most once → O(n) overall.
- Use a **monotonically decreasing** stack for "next greater element" style problems.
- Use a **monotonically increasing** stack for "largest rectangle" style problems.
- **Deque** = double-ended queue. Use for sliding window maximum — maintain a decreasing deque and the front is always the current window's maximum.
- Classic stack problems: valid parentheses, min stack, evaluate RPN, basic calculator.

---

## Exercises

### Easy
1. Implement a stack that returns its minimum in O(1) using a different strategy than the MinStack above — can you do it with a single stack that stores (value, currentMin) pairs?
2. Given a string of brackets, return the minimum number of parentheses to add to make it valid.
3. Implement a queue using a single stack plus recursion.

### Medium
4. Design a monotonic stack solution for "Previous Smaller Element" (mirror of Next Greater Element).
5. Given an array representing a histogram, find the total amount of water that can be trapped (using the stack approach, not the two-pointer approach).
6. Implement a browser history: `visit(url)`, `back(steps)`, `forward(steps)` — all using two stacks.

### Hard
7. Implement `maxSlidingWindow` without using a deque — use a monotonic decreasing stack-based approach instead. What are the tradeoffs?
8. The "Largest Rectangle in Histogram" solved with the monotonic stack approach — trace through an example of height = [2,1,5,6,2,3] step by step, drawing the stack state at each step.
9. Design a data structure that supports: `push(val)`, `pop()`, `popMax()` (remove and return the maximum element). Target O(log n) for all operations.
