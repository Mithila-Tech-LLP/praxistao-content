# Chapter 32: Stacks and Queues

Stacks and queues are the two most fundamental abstract data types. A **stack** is LIFO (Last In, First Out) — like a pile of plates. A **queue** is FIFO (First In, First Out) — like a line at a coffee shop. Understanding them deeply unlocks dozens of algorithms: expression evaluation, tree/graph traversal, task scheduling, and undo/redo systems.

## Table of Contents

1. [Stack — LIFO](#1-stack--lifo)
2. [Stack Applications](#2-stack-applications)
3. [Queue — FIFO](#3-queue--fifo)
4. [Deque — Double-Ended Queue](#4-deque--double-ended-queue)
5. [Priority Queue](#5-priority-queue)
6. [Monotonic Stack/Queue](#6-monotonic-stackqueue)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Stack — LIFO

```
Push(1) → [1]
Push(2) → [1, 2]
Push(3) → [1, 2, 3]
Pop()   → returns 3, stack = [1, 2]
Peek()  → returns 2, stack = [1, 2]
```

**Slice-based implementation (most common in Go):**
```go
type Stack[T any] struct {
    items []T
}

func (s *Stack[T]) Push(v T) {
    s.items = append(s.items, v)
}

func (s *Stack[T]) Pop() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    top := s.items[len(s.items)-1]
    s.items = s.items[:len(s.items)-1]
    return top, true
}

func (s *Stack[T]) Peek() (T, bool) {
    if len(s.items) == 0 {
        var zero T
        return zero, false
    }
    return s.items[len(s.items)-1], true
}

func (s *Stack[T]) Len() int     { return len(s.items) }
func (s *Stack[T]) IsEmpty() bool { return len(s.items) == 0 }
```

**All operations are O(1) amortized** (append may occasionally reallocate, but amortized over many pushes it's O(1)).

### Quick Check
> 1. What does LIFO mean?
> 2. What is the time complexity of Push and Pop on a slice-based stack?
> 3. What is the difference between Pop and Peek?

---

## 2. Stack Applications

### Valid Parentheses
```go
// IsBalanced checks if brackets are properly balanced.
// "()[]{}" → true, "([)]" → false, "{[]}" → true
func IsBalanced(s string) bool {
    stack := &Stack[rune]{}
    matching := map[rune]rune{')': '(', ']': '[', '}': '{'}

    for _, ch := range s {
        switch ch {
        case '(', '[', '{':
            stack.Push(ch)
        case ')', ']', '}':
            top, ok := stack.Pop()
            if !ok || top != matching[ch] {
                return false
            }
        }
    }
    return stack.IsEmpty()
}
```

### Evaluate Postfix (Reverse Polish Notation)
```go
// EvalRPN evaluates a reverse Polish notation expression.
// ["2","1","+","3","*"] → (2+1)*3 = 9
func EvalRPN(tokens []string) int {
    stack := &Stack[int]{}

    for _, token := range tokens {
        switch token {
        case "+", "-", "*", "/":
            b, _ := stack.Pop()  // Right operand (popped first!)
            a, _ := stack.Pop()  // Left operand
            switch token {
            case "+": stack.Push(a + b)
            case "-": stack.Push(a - b)
            case "*": stack.Push(a * b)
            case "/": stack.Push(a / b)
            }
        default:
            n, _ := strconv.Atoi(token)
            stack.Push(n)
        }
    }

    result, _ := stack.Pop()
    return result
}
```

### Daily Temperatures (Next Greater Element)
```go
// DailyTemperatures returns, for each day, how many days until a warmer day.
// temps = [73,74,75,71,69,72,76,73] → [1,1,4,2,1,1,0,0]
func DailyTemperatures(temps []int) []int {
    result := make([]int, len(temps))
    stack := &Stack[int]{}  // Stack of indices

    for i, temp := range temps {
        // Pop all indices with temperatures cooler than today:
        for !stack.IsEmpty() {
            top, _ := stack.Peek()
            if temps[top] < temp {
                stack.Pop()
                result[top] = i - top  // Days until warmer
            } else {
                break
            }
        }
        stack.Push(i)
    }
    return result
}
```

### Undo/Redo System
```go
type Editor struct {
    content   string
    undoStack Stack[string]
    redoStack Stack[string]
}

func (e *Editor) Type(text string) {
    e.undoStack.Push(e.content)
    e.redoStack = Stack[string]{}  // Clear redo on new action
    e.content += text
}

func (e *Editor) Undo() bool {
    prev, ok := e.undoStack.Pop()
    if !ok {
        return false
    }
    e.redoStack.Push(e.content)
    e.content = prev
    return true
}

func (e *Editor) Redo() bool {
    next, ok := e.redoStack.Pop()
    if !ok {
        return false
    }
    e.undoStack.Push(e.content)
    e.content = next
    return true
}
```

### Quick Check
> 1. Why do you pop `b` (right operand) before `a` (left operand) in RPN evaluation?
> 2. In the "next greater element" pattern, what does the stack hold?
> 3. Why must the redo stack be cleared when the user types new text?

---

## 3. Queue — FIFO

```
Enqueue(1) → [1]
Enqueue(2) → [1, 2]
Enqueue(3) → [1, 2, 3]
Dequeue()  → returns 1, queue = [2, 3]
Front()    → returns 2
```

**Naive slice implementation (inefficient):**
```go
// BAD: Dequeue is O(n) because it shifts elements
type SlowQueue[T any] struct{ items []T }
func (q *SlowQueue[T]) Enqueue(v T) { q.items = append(q.items, v) }
func (q *SlowQueue[T]) Dequeue() (T, bool) {
    if len(q.items) == 0 { var z T; return z, false }
    v := q.items[0]
    q.items = q.items[1:]  // O(n) — copies all remaining elements
    return v, true
}
```

**Ring buffer implementation — O(1) for all operations:**
```go
type Queue[T any] struct {
    items []T
    head  int  // Index of front element
    tail  int  // Index of next empty slot
    size  int
}

func NewQueue[T any](capacity int) *Queue[T] {
    return &Queue[T]{items: make([]T, capacity)}
}

func (q *Queue[T]) Enqueue(v T) {
    if q.size == len(q.items) {
        q.grow()
    }
    q.items[q.tail] = v
    q.tail = (q.tail + 1) % len(q.items)
    q.size++
}

func (q *Queue[T]) Dequeue() (T, bool) {
    if q.size == 0 {
        var zero T
        return zero, false
    }
    v := q.items[q.head]
    var zero T
    q.items[q.head] = zero  // Help GC
    q.head = (q.head + 1) % len(q.items)
    q.size--
    return v, true
}

func (q *Queue[T]) Front() (T, bool) {
    if q.size == 0 {
        var zero T
        return zero, false
    }
    return q.items[q.head], true
}

func (q *Queue[T]) Len() int      { return q.size }
func (q *Queue[T]) IsEmpty() bool { return q.size == 0 }

func (q *Queue[T]) grow() {
    newCap := len(q.items) * 2
    if newCap == 0 {
        newCap = 4
    }
    newItems := make([]T, newCap)
    // Copy in order from head to tail:
    for i := 0; i < q.size; i++ {
        newItems[i] = q.items[(q.head+i)%len(q.items)]
    }
    q.items = newItems
    q.head = 0
    q.tail = q.size
}
```

**Channel-based queue (for goroutines):**
```go
// Buffered channel IS a queue for goroutines:
ch := make(chan int, 100)
ch <- 1      // Enqueue
ch <- 2
v := <-ch   // Dequeue: v=1
```

### Queue applications — BFS level order traversal:
```go
// BFS traversal of a binary tree level by level:
func LevelOrder(root *TreeNode) [][]int {
    if root == nil {
        return nil
    }

    var result [][]int
    queue := NewQueue[*TreeNode](16)
    queue.Enqueue(root)

    for !queue.IsEmpty() {
        levelSize := queue.Len()
        level := make([]int, 0, levelSize)

        for i := 0; i < levelSize; i++ {
            node, _ := queue.Dequeue()
            level = append(level, node.Val)
            if node.Left != nil {
                queue.Enqueue(node.Left)
            }
            if node.Right != nil {
                queue.Enqueue(node.Right)
            }
        }
        result = append(result, level)
    }
    return result
}
```

### Quick Check
> 1. Why is a ring buffer more efficient than a plain slice for a queue?
> 2. What does `(index + 1) % capacity` accomplish in a ring buffer?
> 3. When would you use a buffered channel as a queue?

---

## 4. Deque — Double-Ended Queue

A deque (pronounced "deck") supports O(1) push and pop from BOTH ends:

```go
// Built on doubly linked list:
type Deque[T any] struct {
    head *DNode[T]
    tail *DNode[T]
    size int
}

func (d *Deque[T]) PushFront(v T) {
    node := &DNode[T]{Value: v, Next: d.head}
    if d.head != nil {
        d.head.Prev = node
    } else {
        d.tail = node
    }
    d.head = node
    d.size++
}

func (d *Deque[T]) PushBack(v T) {
    node := &DNode[T]{Value: v, Prev: d.tail}
    if d.tail != nil {
        d.tail.Next = node
    } else {
        d.head = node
    }
    d.tail = node
    d.size++
}

func (d *Deque[T]) PopFront() (T, bool) {
    if d.head == nil { var z T; return z, false }
    val := d.head.Value
    d.head = d.head.Next
    if d.head != nil { d.head.Prev = nil } else { d.tail = nil }
    d.size--
    return val, true
}

func (d *Deque[T]) PopBack() (T, bool) {
    if d.tail == nil { var z T; return z, false }
    val := d.tail.Value
    d.tail = d.tail.Prev
    if d.tail != nil { d.tail.Next = nil } else { d.head = nil }
    d.size--
    return val, true
}

func (d *Deque[T]) PeekFront() (T, bool) {
    if d.head == nil { var z T; return z, false }
    return d.head.Value, true
}

func (d *Deque[T]) PeekBack() (T, bool) {
    if d.tail == nil { var z T; return z, false }
    return d.tail.Value, true
}

func (d *Deque[T]) Len() int { return d.size }
```

### Quick Check
> 1. What makes a deque different from a stack and queue?
> 2. What underlying data structure enables O(1) operations on both ends?

---

## 5. Priority Queue

A priority queue returns elements in priority order (highest/lowest first) regardless of insertion order. Implemented with a **heap** (Chapter 38). Here's the interface using Go's `container/heap`:

```go
import "container/heap"

// IntMinHeap implements heap.Interface for a min-heap of ints.
type IntMinHeap []int

func (h IntMinHeap) Len() int           { return len(h) }
func (h IntMinHeap) Less(i, j int) bool { return h[i] < h[j] }  // min-heap
func (h IntMinHeap) Swap(i, j int)      { h[i], h[j] = h[j], h[i] }
func (h *IntMinHeap) Push(x any)        { *h = append(*h, x.(int)) }
func (h *IntMinHeap) Pop() any {
    old := *h
    n := len(old)
    x := old[n-1]
    *h = old[:n-1]
    return x
}

func main() {
    h := &IntMinHeap{5, 2, 8, 1, 9}
    heap.Init(h)

    heap.Push(h, 3)
    fmt.Println(heap.Pop(h))  // 1 — smallest
    fmt.Println(heap.Pop(h))  // 2
    fmt.Println(heap.Pop(h))  // 3
}
```

**Generic priority queue:**
```go
type PriorityQueue[T any] struct {
    items    []T
    lessFunc func(a, b T) bool
}

func NewPQ[T any](less func(T, T) bool) *PriorityQueue[T] {
    return &PriorityQueue[T]{lessFunc: less}
}

func (pq *PriorityQueue[T]) Len() int           { return len(pq.items) }
func (pq *PriorityQueue[T]) Less(i, j int) bool { return pq.lessFunc(pq.items[i], pq.items[j]) }
func (pq *PriorityQueue[T]) Swap(i, j int)      { pq.items[i], pq.items[j] = pq.items[j], pq.items[i] }

func (pq *PriorityQueue[T]) Push(x any) {
    pq.items = append(pq.items, x.(T))
}
func (pq *PriorityQueue[T]) Pop() any {
    old := pq.items
    n := len(old)
    x := old[n-1]
    pq.items = old[:n-1]
    return x
}

func (pq *PriorityQueue[T]) Enqueue(v T) {
    heap.Push(pq, v)
}
func (pq *PriorityQueue[T]) Dequeue() (T, bool) {
    if pq.Len() == 0 { var z T; return z, false }
    return heap.Pop(pq).(T), true
}
func (pq *PriorityQueue[T]) Peek() (T, bool) {
    if pq.Len() == 0 { var z T; return z, false }
    return pq.items[0], true
}

// Usage:
pq := NewPQ[int](func(a, b int) bool { return a < b })  // min-heap
heap.Init(pq)
pq.Enqueue(5)
pq.Enqueue(1)
pq.Enqueue(3)
v, _ := pq.Dequeue()  // 1
```

### Quick Check
> 1. What ordering does a min-heap priority queue provide?
> 2. What interface must you implement to use Go's `container/heap`?

---

## 6. Monotonic Stack/Queue

A **monotonic stack** maintains elements in strictly increasing or decreasing order. Elements that violate the monotonic property are popped before a new element is pushed.

**Monotonic stack — largest rectangle in histogram:**
```go
// LargestRectangle finds the largest rectangle area in a histogram.
// heights = [2,1,5,6,2,3] → 10 (the 5×2 rectangle using bars 3 and 4)
func LargestRectangle(heights []int) int {
    stack := &Stack[int]{}  // Stack of indices (monotonically increasing by height)
    maxArea := 0
    heights = append(heights, 0)  // Sentinel: force all remaining bars to be processed

    for i, h := range heights {
        for !stack.IsEmpty() {
            top, _ := stack.Peek()
            if heights[top] > h {
                stack.Pop()
                width := i
                if !stack.IsEmpty() {
                    prev, _ := stack.Peek()
                    width = i - prev - 1
                }
                area := heights[top] * width
                if area > maxArea {
                    maxArea = area
                }
            } else {
                break
            }
        }
        stack.Push(i)
    }
    return maxArea
}
```

**Monotonic deque — sliding window maximum:**
```go
// SlidingWindowMax returns the maximum in each window of size k.
// nums = [1,3,-1,-3,5,3,6,7], k=3 → [3,3,5,5,6,7]
func SlidingWindowMax(nums []int, k int) []int {
    dq := &Deque[int]{}  // Stores indices, decreasing order of nums[i]
    result := []int{}

    for i, n := range nums {
        // Remove elements outside the window from front:
        for dq.Len() > 0 {
            front, _ := dq.PeekFront()
            if front <= i-k {
                dq.PopFront()
            } else {
                break
            }
        }
        // Remove elements smaller than n from back:
        for dq.Len() > 0 {
            back, _ := dq.PeekBack()
            if nums[back] <= n {
                dq.PopBack()
            } else {
                break
            }
        }
        dq.PushBack(i)

        // Window is full — record maximum (front of deque):
        if i >= k-1 {
            front, _ := dq.PeekFront()
            result = append(result, nums[front])
        }
    }
    return result
}
```

---

## Summary

- **Stack** (LIFO): Push/Pop/Peek — O(1); use for: bracket matching, RPN eval, undo, DFS
- **Queue** (FIFO): Enqueue/Dequeue — O(1) with ring buffer; use for: BFS, task queues
- **Ring buffer**: fixes the O(n) dequeue of naive slice queue; `(idx+1) % cap` wraps around
- **Deque**: push/pop from both ends — O(1) with doubly linked list; enables sliding window algorithms
- **Priority queue**: `container/heap`; implement `Len`, `Less`, `Swap`, `Push`, `Pop`
- **Monotonic stack**: elements in sorted order; pop elements violating order on push
- **Monotonic deque**: front = max/min of current window; use for O(n) sliding window problems

---

## Exercises

### Easy
1. Implement `Min Stack` — a stack that supports `Push`, `Pop`, `Top`, and `GetMin` all in O(1). Use two stacks internally: one for values, one to track the current minimum.
2. Write `QueueFromStacks[T any]` — implement a queue using two stacks. `Enqueue` uses Stack 1. `Dequeue` moves all of Stack 1 to Stack 2 (if Stack 2 is empty), then pops from Stack 2. Verify FIFO ordering.
3. Write `IsValidExpression(expr string) bool` that checks whether a mathematical expression has properly balanced parentheses AND operators (+, -, *, /) appear only between operands (basic structural check, not full parse).

### Medium
4. Task scheduler: Implement a `Scheduler` using a priority queue where tasks have a priority (1=highest, 5=lowest) and a name. `AddTask(name string, priority int)`, `RunNext() string` executes the highest-priority task (lowest number). `Reschedule(name string, newPriority int)` changes priority. Test with 20 tasks of mixed priorities added in random order — verify they execute in priority order.
5. Circular buffer with blocking: Implement a goroutine-safe bounded queue (ring buffer) that blocks on `Enqueue` when full and blocks on `Dequeue` when empty, using channels and a mutex. This is a classic producer-consumer queue. Test with 3 producers and 2 consumers, total 1000 items — verify no items are lost and no deadlock.
6. Stock span problem: For each day's stock price, find the span — the number of consecutive days before it (including itself) where the price was ≤ today's price. `Span([100,80,60,70,60,75,85])` → `[1,1,1,2,1,4,6]`. Implement in O(n) using a monotonic stack.

### Hard
7. LFU cache: Implement a Least Frequently Used cache with O(1) get and put. Use: a hash map of key→value, a hash map of key→frequency, a hash map of frequency→doubly linked list of keys, and track `minFreq`. When capacity is full, evict the key with lowest frequency (if tie: evict the least recently used among those). Test with the LeetCode 460 test cases.
8. Expression tree evaluator: Parse a mathematical expression string `"3 + 4 * 2 / ( 1 - 5 ) ^ 2"` into an expression tree using two stacks (operator stack and operand stack, shunting-yard algorithm). `Parse(expr string) *ExprNode`, `Eval(node *ExprNode) float64`, `String(node *ExprNode) string` (inorder with parens). Handle: `+`, `-`, `*`, `/`, `^` (power), unary minus, and parentheses with correct precedence.
