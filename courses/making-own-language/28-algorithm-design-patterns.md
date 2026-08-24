# Chapter 28: Algorithm Design Patterns — The Master Toolkit

> "Every problem that is solved has a pattern behind it. Learn the pattern, not the problem." — Donald Knuth (paraphrased)

---

## Overview

In the previous chapters we learned individual data structures and specific algorithms. But experienced programmers do not memorize hundreds of algorithms — they learn a smaller set of **design patterns** and recognize which pattern applies to each new problem. A pattern is a general blueprint, a way of thinking, that can be applied to an entire family of problems.

This chapter is your master toolkit. We will cover the eight most powerful algorithm design patterns: Divide and Conquer, Greedy Algorithms, Two Pointers, Sliding Window, Monotonic Stack, Prefix Sums, Bit Manipulation, and a systematic problem-solving approach you can apply to any new problem you encounter. Every pattern comes with a concrete Go implementation, and we end with a milestone showing how each phase of the Astra compiler uses one of these exact patterns.

By the end of this chapter, when you see a new problem, you will no longer stare at a blank page. You will think: "Does this look like a sliding window problem? A divide and conquer? A greedy?" Having that vocabulary changes everything.

---

## What We Are Building

By the end of this chapter you will have:

- Understood and implemented all eight major algorithm design patterns
- A mental framework for recognizing which pattern to apply
- Deep understanding of why each pattern works (not just how)
- Complete Go implementations of key examples for every pattern
- A mapping of the Astra compiler phases to their underlying algorithmic patterns
- A systematic 7-step approach to solving any algorithm problem

---

## Table of Contents

1. Divide and Conquer
2. Greedy Algorithms
3. Two-Pointer Technique
4. Sliding Window
5. Monotonic Stack
6. Prefix Sums
7. Bit Manipulation
8. How to Approach Any Algorithm Problem
9. Astra Build Milestone: The Compiler as Algorithmic Patterns

---

## 1. Divide and Conquer

### The Core Idea

Divide and conquer is one of the oldest and most powerful patterns in computer science. The idea is deceptively simple:

1. **Divide**: Split the problem into smaller subproblems of the same type.
2. **Conquer**: Solve each subproblem (recursively, or directly if small enough).
3. **Combine**: Merge the subproblem solutions into the overall solution.

The magic is that the subproblems are the same type as the original. This leads to recursive algorithms that often achieve logarithmic or log-linear time complexity.

```
solve(problem):
  if problem is small enough:
      return directSolve(problem)
  left = solve(leftHalf(problem))
  right = solve(rightHalf(problem))
  return combine(left, right)
```

### Time Complexity: The Master Theorem

When you divide a problem of size n into `a` subproblems of size `n/b`, and the combine step takes `O(n^d)` time, the total time complexity is:

```
T(n) = a * T(n/b) + O(n^d)

Result:
  If d > log_b(a):  T(n) = O(n^d)
  If d = log_b(a):  T(n) = O(n^d * log n)
  If d < log_b(a):  T(n) = O(n^(log_b(a)))
```

Let us verify this with merge sort:
- a = 2 (two subproblems)
- b = 2 (each is half size)
- d = 1 (merge step is O(n))
- log_2(2) = 1, and d = 1 → case 2 → O(n log n). Confirmed!

### Example 1: Merge Sort

Merge sort is the canonical divide and conquer algorithm. It:
1. Divides the array in half
2. Recursively sorts each half
3. Merges the two sorted halves

```go
package algorithms

// MergeSort sorts a slice of integers in O(n log n) time.
func MergeSort(arr []int) []int {
    // Base case: arrays of length 0 or 1 are already sorted
    if len(arr) <= 1 {
        return arr
    }

    // Divide: find midpoint and split
    mid := len(arr) / 2
    left := MergeSort(arr[:mid])
    right := MergeSort(arr[mid:])

    // Combine: merge the two sorted halves
    return merge(left, right)
}

// merge merges two sorted slices into one sorted slice.
func merge(left, right []int) []int {
    result := make([]int, 0, len(left)+len(right))
    i, j := 0, 0

    // Compare and pick the smaller element
    for i < len(left) && j < len(right) {
        if left[i] <= right[j] {
            result = append(result, left[i])
            i++
        } else {
            result = append(result, right[j])
            j++
        }
    }

    // Append any remaining elements
    result = append(result, left[i:]...)
    result = append(result, right[j:]...)
    return result
}
```

### Example 2: Binary Search

Binary search is divide and conquer at its most elegant: each step divides the search space in half.

```go
// BinarySearch finds target in a sorted slice. Returns index or -1.
func BinarySearch(arr []int, target int) int {
    low, high := 0, len(arr)-1

    for low <= high {
        mid := low + (high-low)/2 // avoid integer overflow

        if arr[mid] == target {
            return mid
        } else if arr[mid] < target {
            low = mid + 1 // target is in the right half
        } else {
            high = mid - 1 // target is in the left half
        }
    }
    return -1 // not found
}
```

**Why `mid = low + (high-low)/2` and not `(low+high)/2`?**

Because `low + high` can overflow a 32-bit integer if both are large. The subtraction form avoids this.

### Example 3: Fast Exponentiation

Computing `x^n` naively takes O(n) multiplications. Divide and conquer gives us O(log n):

```go
// FastPow computes x^n using divide and conquer (exponentiation by squaring).
func FastPow(x, n int) int {
    if n == 0 {
        return 1
    }
    if n%2 == 0 {
        // x^n = (x^(n/2))^2
        half := FastPow(x, n/2)
        return half * half
    }
    // x^n = x * x^(n-1)
    return x * FastPow(x, n-1)
}
```

This is O(log n) multiplications — computing `2^64` takes only 6 squarings instead of 64 multiplications.

---

## 2. Greedy Algorithms

### The Core Idea

A greedy algorithm makes the **locally optimal choice** at each step, hoping that this leads to a globally optimal solution. It never looks back — once a choice is made, it is final.

The key question with greedy algorithms is always: **does making the locally optimal choice lead to the globally optimal solution?** Sometimes it does (and greedy is both simple and efficient). Sometimes it does not (and you need dynamic programming instead).

### When Greedy Works

Two conditions make greedy algorithms provably correct:

1. **Greedy choice property**: A globally optimal solution can be built by making locally optimal (greedy) choices. You do not need to reconsider previous choices.

2. **Optimal substructure**: An optimal solution to the problem contains optimal solutions to its subproblems.

These are the same conditions as dynamic programming, but greedy adds the greedy choice property as an extra guarantee. When both hold, greedy is preferred (simpler, faster).

### Example 1: Activity Selection Problem

Given a set of activities with start and end times, select the maximum number of non-overlapping activities.

**Greedy choice**: Always pick the activity that finishes earliest.

**Why does this work?** By finishing early, we leave maximum room for future activities. Any solution that picks a later-finishing first activity can be improved by swapping it for the earlier-finishing one — this is called the **exchange argument**.

```go
// Activity represents a task with a start and end time.
type Activity struct {
    Start, End int
    Name       string
}

// ActivitySelection returns the maximum set of non-overlapping activities.
func ActivitySelection(activities []Activity) []Activity {
    if len(activities) == 0 {
        return nil
    }

    // Sort by end time — the greedy choice
    sort.Slice(activities, func(i, j int) bool {
        return activities[i].End < activities[j].End
    })

    selected := []Activity{activities[0]}
    lastEnd := activities[0].End

    for _, act := range activities[1:] {
        // Greedily pick if it starts after the last selected activity ends
        if act.Start >= lastEnd {
            selected = append(selected, act)
            lastEnd = act.End
        }
    }

    return selected
}
```

### Example 2: Huffman Coding (used in compression)

Huffman coding assigns shorter bit sequences to more frequent characters and longer sequences to less frequent ones, minimizing total encoded length. It uses a greedy algorithm: always merge the two least-frequent symbols.

```go
// HuffmanNode represents a node in the Huffman tree.
type HuffmanNode struct {
    Char      rune
    Frequency int
    Left      *HuffmanNode
    Right     *HuffmanNode
}

// HuffmanCodes builds optimal prefix-free codes for characters.
func HuffmanCodes(frequencies map[rune]int) map[rune]string {
    // Initialize priority queue with one node per character
    pq := &HuffmanPQ{}
    for char, freq := range frequencies {
        heap.Push(pq, &HuffmanNode{Char: char, Frequency: freq})
    }
    heap.Init(pq)

    // Greedy: repeatedly merge the two least frequent nodes
    for pq.Len() > 1 {
        left := heap.Pop(pq).(*HuffmanNode)
        right := heap.Pop(pq).(*HuffmanNode)
        merged := &HuffmanNode{
            Frequency: left.Frequency + right.Frequency,
            Left:      left,
            Right:     right,
        }
        heap.Push(pq, merged)
    }

    // The remaining node is the root of the Huffman tree
    root := heap.Pop(pq).(*HuffmanNode)

    // Walk the tree to assign codes
    codes := make(map[rune]string)
    var walkTree func(node *HuffmanNode, code string)
    walkTree = func(node *HuffmanNode, code string) {
        if node == nil {
            return
        }
        if node.Left == nil && node.Right == nil {
            // Leaf node: this is a character
            codes[node.Char] = code
            return
        }
        walkTree(node.Left, code+"0")
        walkTree(node.Right, code+"1")
    }
    walkTree(root, "")

    return codes
}

// HuffmanPQ implements heap.Interface for HuffmanNodes.
type HuffmanPQ []*HuffmanNode

func (pq HuffmanPQ) Len() int            { return len(pq) }
func (pq HuffmanPQ) Less(i, j int) bool { return pq[i].Frequency < pq[j].Frequency }
func (pq HuffmanPQ) Swap(i, j int)      { pq[i], pq[j] = pq[j], pq[i] }
func (pq *HuffmanPQ) Push(x interface{}) { *pq = append(*pq, x.(*HuffmanNode)) }
func (pq *HuffmanPQ) Pop() interface{} {
    old := *pq
    n := len(old)
    x := old[n-1]
    *pq = old[:n-1]
    return x
}
```

### When Greedy Fails

The classic counterexample is the **coin change problem** with arbitrary coin denominations. If coins are [1, 5, 6] and target is 10:

- Greedy (largest first): 6 + 1 + 1 + 1 + 1 = 5 coins
- Optimal: 5 + 5 = 2 coins

Greedy fails here. You need dynamic programming.

---

## 3. Two-Pointer Technique

### The Core Idea

The two-pointer technique uses two indices (pointers) that move through a data structure, usually an array or string. The trick is that by moving pointers intelligently based on what we see, we avoid redundant work and achieve O(n) time for problems that a brute force double loop would solve in O(n²).

### Three Major Variations

**Variation 1: Both pointers start at the left, one moves faster (slow+fast)**

Used for: cycle detection (Floyd's algorithm), finding the middle of a linked list, removing duplicates in-place.

**Variation 2: Pointers start from opposite ends, move toward each other**

Used for: finding pairs that sum to a target (sorted array), checking if a string is a palindrome, two-sum in sorted arrays.

**Variation 3: Both start at left, create a window between them**

Used for: longest substring, remove duplicates. (This transitions into the sliding window pattern.)

### Example 1: Two Sum in a Sorted Array

Find two numbers in a sorted array that add to target. O(n) time.

```go
// TwoSumSorted finds indices of two numbers in a sorted array that sum to target.
// Returns (-1, -1) if no pair exists.
func TwoSumSorted(arr []int, target int) (int, int) {
    left, right := 0, len(arr)-1

    for left < right {
        sum := arr[left] + arr[right]
        if sum == target {
            return left, right
        } else if sum < target {
            left++  // need a larger sum → move left pointer right
        } else {
            right-- // need a smaller sum → move right pointer left
        }
    }
    return -1, -1
}
```

Why does this work? In a sorted array, moving `left` right increases the sum; moving `right` left decreases it. We always make progress toward the target.

### Example 2: Remove Duplicates In-Place

Remove duplicates from a sorted array without using extra space. Return the new length.

```go
// RemoveDuplicates removes duplicates from a sorted array in-place.
// Returns the length of the deduplicated prefix.
func RemoveDuplicates(arr []int) int {
    if len(arr) == 0 {
        return 0
    }

    slow := 0 // slow points to the last unique element written

    for fast := 1; fast < len(arr); fast++ {
        if arr[fast] != arr[slow] {
            slow++
            arr[slow] = arr[fast]
        }
        // If arr[fast] == arr[slow], it's a duplicate → just advance fast
    }

    return slow + 1 // length of unique prefix
}
```

### Example 3: Container with Most Water

Given heights of vertical bars, find two bars that together with the x-axis form the container with the most water. Classic O(n) two-pointer solution.

```go
// MaxWater finds the maximum water container using two pointers.
func MaxWater(heights []int) int {
    left, right := 0, len(heights)-1
    maxWater := 0

    for left < right {
        // Width = right - left, height = min of the two bars
        width := right - left
        height := heights[left]
        if heights[right] < height {
            height = heights[right]
        }
        water := width * height
        if water > maxWater {
            maxWater = water
        }

        // Move the shorter bar inward — moving the taller bar can only hurt
        if heights[left] < heights[right] {
            left++
        } else {
            right--
        }
    }
    return maxWater
}
```

---

## 4. Sliding Window

### The Core Idea

The sliding window pattern maintains a contiguous subarray (or substring) called the "window" and slides it through the data. As the window slides right, we add the new element on the right and remove the old element on the left. This avoids recomputing the entire window from scratch at each position — the key to O(n) efficiency.

```
Array:   [2, 1, 5, 1, 3, 2]
Window (size 3):
  [2, 1, 5] → sum=8
  [1, 5, 1] → sum=7  (subtract 2, add 1)
  [5, 1, 3] → sum=9  (subtract 1, add 3)
  [1, 3, 2] → sum=6  (subtract 5, add 2)
```

### Fixed-Size Window vs Variable-Size Window

**Fixed window**: Window size k is given. Slide one step at a time, subtracting the leftmost element and adding the new rightmost element.

**Variable window**: Expand the window on the right; shrink from the left when a constraint is violated. The window grows and shrinks dynamically.

### Example 1: Maximum Sum Subarray of Size K (Fixed Window)

```go
// MaxSumSubarrayK finds the maximum sum of any subarray of size k.
func MaxSumSubarrayK(arr []int, k int) int {
    if len(arr) < k {
        return -1
    }

    // Compute sum of first window
    windowSum := 0
    for i := 0; i < k; i++ {
        windowSum += arr[i]
    }
    maxSum := windowSum

    // Slide the window: add arr[i], subtract arr[i-k]
    for i := k; i < len(arr); i++ {
        windowSum += arr[i] - arr[i-k]
        if windowSum > maxSum {
            maxSum = windowSum
        }
    }
    return maxSum
}
```

### Example 2: Longest Substring Without Repeating Characters (Variable Window)

This is the variable window pattern: expand right, shrink left when we see a repeat.

```go
// LongestUniqueSubstring finds the length of the longest substring
// that contains no repeating characters.
func LongestUniqueSubstring(s string) int {
    // charIndex[c] = last seen index of character c
    charIndex := make(map[byte]int)
    maxLen := 0
    left := 0 // left boundary of window

    for right := 0; right < len(s); right++ {
        c := s[right]

        // If c was seen and its last occurrence is inside the current window,
        // move left to exclude the duplicate.
        if idx, seen := charIndex[c]; seen && idx >= left {
            left = idx + 1
        }

        charIndex[c] = right // update last seen position

        windowLen := right - left + 1
        if windowLen > maxLen {
            maxLen = windowLen
        }
    }
    return maxLen
}
```

```
s = "abcabcbb"
right=0 (a): window=a,       left=0, len=1, max=1
right=1 (b): window=ab,      left=0, len=2, max=2
right=2 (c): window=abc,     left=0, len=3, max=3
right=3 (a): a seen at 0 ≥ left(0), left→1
             window=bca,     left=1, len=3, max=3
right=4 (b): b seen at 1 ≥ left(1), left→2
             window=cab,     left=2, len=3, max=3
right=5 (c): c seen at 2 ≥ left(2), left→3
             window=abc,     left=3, len=3, max=3
right=6 (b): b seen at 4 ≥ left(3), left→5
             window=cb,      left=5, len=2, max=3
right=7 (b): b seen at 6 ≥ left(5), left→7
             window=b,       left=7, len=1, max=3
Answer: 3
```

### The Expand/Shrink Pattern Template

```go
func slidingWindow(arr []int, condition func(window []int) bool) int {
    left := 0
    result := 0
    // Track window state (varies by problem)
    for right := 0; right < len(arr); right++ {
        // Expand: add arr[right] to window state
        for !condition(arr[left:right+1]) {
            // Shrink: remove arr[left] from window state
            left++
        }
        // Now window [left..right] satisfies condition
        windowLen := right - left + 1
        if windowLen > result {
            result = windowLen
        }
    }
    return result
}
```

---

## 5. Monotonic Stack

### The Core Idea

A **monotonic stack** is a stack where elements are always in monotonically increasing or decreasing order. When we push an element and it would violate the monotonic property, we pop elements until the property is restored.

The key insight: by maintaining this order, we can answer "what is the next greater element?" for every position in O(n) total time — even though it looks like it should take O(n²) to compare every element with every other element.

### ASCII Diagram: Next Greater Element

```
Array: [4, 6, 2, 5, 7, 1]
We want: nextGreater[i] = the first element to the right of i that is greater.

Process from right to left, using a stack:
i=5 (val=1): stack=[],    push 1.  nextGreater[5]=-1 (nothing to right)
i=4 (val=7): stack=[1],   pop 1 (1 < 7), stack=[]. push 7. nextGreater[4]=-1
i=3 (val=5): stack=[7],   7 > 5, nextGreater[3]=7. push 5. stack=[7,5]
i=2 (val=2): stack=[7,5], 5 > 2, nextGreater[2]=5. push 2. stack=[7,5,2]
i=1 (val=6): stack=[7,5,2], pop 2 (2<6), pop 5 (5<6). 7>6, nextGreater[1]=7.
             push 6. stack=[7,6]
i=0 (val=4): stack=[7,6], 6 > 4, nextGreater[0]=6. push 4. stack=[7,6,4]

Result: nextGreater = [6, 7, 5, 7, -1, -1]
```

### Complete Go Implementation

```go
// NextGreaterElement returns, for each element, the next greater element to its right.
// Returns -1 if no greater element exists.
func NextGreaterElement(arr []int) []int {
    n := len(arr)
    result := make([]int, n)
    for i := range result {
        result[i] = -1
    }

    // Stack stores indices (not values) so we can update result[i]
    stack := []int{} // monotonic decreasing stack (top has smallest)

    for i := n - 1; i >= 0; i-- {
        // Pop all elements from stack that are ≤ arr[i]
        for len(stack) > 0 && arr[stack[len(stack)-1]] <= arr[i] {
            stack = stack[:len(stack)-1]
        }
        // Top of stack is the next greater element
        if len(stack) > 0 {
            result[i] = arr[stack[len(stack)-1]]
        }
        stack = append(stack, i)
    }
    return result
}

// LargestRectangleInHistogram finds the area of the largest rectangle
// that can be formed in a histogram. O(n) using monotonic stack.
func LargestRectangleInHistogram(heights []int) int {
    stack := []int{} // indices of bars in increasing height order
    maxArea := 0

    for i := 0; i <= len(heights); i++ {
        currentHeight := 0
        if i < len(heights) {
            currentHeight = heights[i]
        }

        // Process all bars taller than current (they can't extend rightward)
        for len(stack) > 0 && heights[stack[len(stack)-1]] > currentHeight {
            topIdx := stack[len(stack)-1]
            stack = stack[:len(stack)-1]

            height := heights[topIdx]
            width := i
            if len(stack) > 0 {
                width = i - stack[len(stack)-1] - 1
            }
            area := height * width
            if area > maxArea {
                maxArea = area
            }
        }
        stack = append(stack, i)
    }
    return maxArea
}
```

**Why O(n)?** Each element is pushed onto the stack exactly once and popped at most once. So the total number of push + pop operations is O(n), even though the inner loop looks like it could be O(n) per outer iteration.

---

## 6. Prefix Sums

### The Core Idea

A prefix sum array (also called a cumulative sum array) precomputes the sum of all elements up to each index. Once computed, you can answer any range sum query in O(1) — even though computing the sum naively would take O(n) per query.

```
Array:  [3, 1, 4, 1, 5, 9, 2, 6]
Prefix: [0, 3, 4, 8, 9, 14, 23, 25, 31]

Sum of arr[2..5] = prefix[6] - prefix[2] = 23 - 4 = 19
Verify: 4 + 1 + 5 + 9 = 19 ✓
```

### Building and Querying

```go
// PrefixSumArray precomputes prefix sums for O(1) range queries.
type PrefixSumArray struct {
    prefix []int
}

// NewPrefixSumArray builds a prefix sum array from arr in O(n).
func NewPrefixSumArray(arr []int) *PrefixSumArray {
    prefix := make([]int, len(arr)+1)
    prefix[0] = 0
    for i, val := range arr {
        prefix[i+1] = prefix[i] + val
    }
    return &PrefixSumArray{prefix: prefix}
}

// RangeSum returns the sum of arr[left..right] (inclusive) in O(1).
func (p *PrefixSumArray) RangeSum(left, right int) int {
    return p.prefix[right+1] - p.prefix[left]
}

// SubarraySumEqualsK counts subarrays with sum equal to k.
// Uses prefix sums + hash map for O(n) time.
func SubarraySumEqualsK(arr []int, k int) int {
    // prefixCount[s] = number of times prefix sum s has appeared
    prefixCount := map[int]int{0: 1}
    currentSum := 0
    count := 0

    for _, val := range arr {
        currentSum += val
        // If currentSum - k has been seen before,
        // then there is a subarray ending here with sum k
        if c, found := prefixCount[currentSum-k]; found {
            count += c
        }
        prefixCount[currentSum]++
    }
    return count
}
```

### 2D Prefix Sums

For a 2D grid, we can extend the idea to answer rectangle sum queries in O(1):

```go
// PrefixSum2D precomputes 2D prefix sums for O(1) rectangle queries.
type PrefixSum2D struct {
    prefix [][]int
    rows   int
    cols   int
}

func NewPrefixSum2D(grid [][]int) *PrefixSum2D {
    rows, cols := len(grid), len(grid[0])
    prefix := make([][]int, rows+1)
    for i := range prefix {
        prefix[i] = make([]int, cols+1)
    }

    for r := 1; r <= rows; r++ {
        for c := 1; c <= cols; c++ {
            prefix[r][c] = grid[r-1][c-1] +
                prefix[r-1][c] + prefix[r][c-1] - prefix[r-1][c-1]
        }
    }

    return &PrefixSum2D{prefix: prefix, rows: rows, cols: cols}
}

// RectangleSum returns the sum of the rectangle from (r1,c1) to (r2,c2) inclusive.
func (p *PrefixSum2D) RectangleSum(r1, c1, r2, c2 int) int {
    return p.prefix[r2+1][c2+1] -
        p.prefix[r1][c2+1] -
        p.prefix[r2+1][c1] +
        p.prefix[r1][c1]
}
```

### Difference Arrays for Range Updates

A **difference array** is the inverse of prefix sums. It lets you perform range update operations (add a value to all elements in a range) in O(1), with O(n) at the end to reconstruct the final array.

```go
// DifferenceArray supports O(1) range updates, O(n) final reconstruction.
type DifferenceArray struct {
    diff []int
    n    int
}

func NewDifferenceArray(arr []int) *DifferenceArray {
    diff := make([]int, len(arr)+1)
    diff[0] = arr[0]
    for i := 1; i < len(arr); i++ {
        diff[i] = arr[i] - arr[i-1]
    }
    return &DifferenceArray{diff: diff, n: len(arr)}
}

// RangeAdd adds val to all elements from index l to r (inclusive) in O(1).
func (d *DifferenceArray) RangeAdd(l, r, val int) {
    d.diff[l] += val
    if r+1 <= d.n {
        d.diff[r+1] -= val
    }
}

// Reconstruct computes the final array after all range updates in O(n).
func (d *DifferenceArray) Reconstruct() []int {
    result := make([]int, d.n)
    result[0] = d.diff[0]
    for i := 1; i < d.n; i++ {
        result[i] = result[i-1] + d.diff[i]
    }
    return result
}
```

---

## 7. Bit Manipulation

### The Core Idea

At the lowest level, computers work on bits — 0s and 1s. Bit manipulation uses bitwise operators (&, |, ^, ~, <<, >>) to perform operations directly on the binary representation of numbers. This often achieves O(1) operations and avoids expensive division, modulo, or multiplication.

### Essential Bit Tricks

```go
// IsPowerOfTwo checks if n is a power of 2.
// A power of 2 in binary has exactly one '1' bit: 8 = 1000
// n - 1 flips all bits below the highest '1': 7 = 0111
// n & (n-1) clears the lowest set bit. For powers of 2, result is 0.
func IsPowerOfTwo(n int) bool {
    return n > 0 && (n&(n-1)) == 0
}

// CountSetBits counts the number of 1-bits in n (Hamming weight).
// Brian Kernighan's algorithm: n & (n-1) clears the lowest set bit.
func CountSetBits(n int) int {
    count := 0
    for n != 0 {
        n = n & (n - 1) // clear lowest set bit
        count++
    }
    return count
}

// GetBit returns the value of bit at position pos (0 = rightmost).
func GetBit(n, pos int) int {
    return (n >> pos) & 1
}

// SetBit sets bit at position pos to 1.
func SetBit(n, pos int) int {
    return n | (1 << pos)
}

// ClearBit clears bit at position pos (sets to 0).
func ClearBit(n, pos int) int {
    return n & ^(1 << pos)
}

// ToggleBit flips bit at position pos.
func ToggleBit(n, pos int) int {
    return n ^ (1 << pos)
}

// XORSwap swaps a and b without using a temporary variable.
// Note: only for educational purposes; real code should use: a, b = b, a
func XORSwap(a, b *int) {
    *a ^= *b
    *b ^= *a
    *a ^= *b
}

// FastModPowerOf2 computes n % m where m is a power of 2.
// n % m == n & (m - 1) when m is a power of 2.
func FastModPowerOf2(n, m int) int {
    return n & (m - 1)
}

// FindSingleNumber finds the element that appears once when all others appear twice.
// XOR of identical numbers is 0. XOR of 0 with any number is that number.
func FindSingleNumber(nums []int) int {
    result := 0
    for _, n := range nums {
        result ^= n
    }
    return result
}
```

### Application: Bitset for Set Membership

A bitset uses a single integer (or array of integers) as a compact set. Each bit represents membership of one element:

```go
// Bitset represents a set of integers from 0 to 63 using a single uint64.
type Bitset struct {
    bits uint64
}

func (b *Bitset) Add(n int)      { b.bits = SetBit(int(b.bits), n) }
func (b *Bitset) Remove(n int)   { b.bits = ClearBit(int(b.bits), n) }
func (b *Bitset) Contains(n int) bool { return GetBit(int(b.bits), n) == 1 }
func (b *Bitset) Union(other Bitset) Bitset { return Bitset{b.bits | other.bits} }
func (b *Bitset) Intersection(other Bitset) Bitset { return Bitset{b.bits & other.bits} }
```

A single integer can represent a set of 64 elements! This is 64x more memory-efficient than a boolean array and operations are O(1).

---

## 8. How to Approach Any Algorithm Problem

When faced with a new algorithm problem, panic is your enemy and process is your friend. Here is a systematic 7-step framework:

### Step 1: Understand the Problem Completely

Before writing a single line of code, restate the problem in your own words. What are the inputs? What are the outputs? What are the constraints? What does "optimal" mean here?

Ask: "Am I minimizing something? Maximizing? Counting? Searching? Sorting?"

### Step 2: Work Through Concrete Examples by Hand

Take the given examples and trace through them manually. Then create your own examples, including edge cases:
- Empty input (n=0)
- Single element (n=1)
- All elements the same
- Already sorted / reverse sorted
- Maximum possible input size

### Step 3: Write a Brute Force Solution First

Even if it is O(n³) or O(2^n), write a correct solution. This gives you:
- A reference implementation to check your optimized version against
- Clarity about what the problem is actually asking
- A working solution (partial credit in interviews, working code is better than no code)

### Step 4: Identify the Bottleneck and Pattern

Where is your brute force slow? What does it recompute repeatedly?
- Overlapping subproblems → Dynamic Programming or Memoization
- Array/string with contiguous subarray constraints → Sliding Window or Prefix Sum
- Need to find next greater/smaller → Monotonic Stack
- Sorted data, find pair → Two Pointers
- Make a local choice that doesn't need revision → Greedy
- Problem halves at each step → Divide and Conquer
- Graph connectivity, cycles, paths → BFS/DFS/Dijkstra

### Step 5: Implement the Optimized Solution

Now write the optimized version. Be methodical, not rushed. Name variables clearly. Add comments for non-obvious steps.

### Step 6: Test and Handle Edge Cases

Run your solution on:
1. The provided examples
2. Your edge cases from Step 2
3. The brute force output for random inputs

### Step 7: Analyze Complexity

State the time and space complexity and justify it. "It is O(n log n) because we sort once (O(n log n)) and then do a single linear scan (O(n))."

---

## 9. Astra Build Milestone: The Compiler as Algorithmic Patterns

Every phase of the Astra compiler is an application of the algorithm patterns we learned. This is not a coincidence — compiler writers have been using these patterns for 70 years because they match the structure of the problems compilers face.

```
Phase              | Algorithm Pattern          | Why This Pattern          | Complexity
───────────────────┼────────────────────────────┼───────────────────────────┼──────────────
Lexer              | Sliding window on text     | Scan chars, emit tokens   | O(n)
Parser             | Divide and conquer (RD)    | Rules call each other     | O(n)
Symbol resolution  | Two-pass + hash lookup     | Forward references        | O(n)
Type checking      | Constraint propagation     | Unify type constraints    | O(n log n)
Constant folding   | Tree walk (post-order)     | Evaluate leaves first     | O(n)
Dead code elim.    | DFS from entry point       | Reachability in call graph| O(V+E)
Register alloc.    | Greedy (linear scan)       | Assign regs greedily      | O(n)
Code generation    | Tree walk (pre-order)      | Emit parent before child  | O(n)
Linker             | Topological sort           | Link dependencies first   | O(V+E)
```

Let us look at each connection more carefully.

### The Lexer: Sliding Window on Source Text

The lexer scans the source text character by character. It maintains a window: the `start` marker points to the beginning of the current token, and `current` advances to scan characters.

```go
// lexer/lexer.go (simplified excerpt showing the sliding window)
type Lexer struct {
    source  string
    start   int // beginning of current token (left edge of window)
    current int // current character being examined (right edge of window)
    line    int
}

func (l *Lexer) NextToken() Token {
    l.start = l.current // slide the window: left = right

    if l.isAtEnd() {
        return l.makeToken(EOF)
    }

    c := l.advance()

    switch c {
    case '+': return l.makeToken(PLUS)
    case '-': return l.makeToken(MINUS)
    // ... more single-char tokens

    case '"': return l.scanString() // expand right until closing quote
    default:
        if isDigit(c) { return l.scanNumber() }  // expand right for number
        if isAlpha(c) { return l.scanIdent() }   // expand right for identifier
    }
    return l.errorToken("Unexpected character")
}

// Each scan function expands the right edge of the window:
func (l *Lexer) scanNumber() Token {
    for isDigit(l.peek()) {
        l.advance() // expand window right
    }
    if l.peek() == '.' && isDigit(l.peekNext()) {
        l.advance() // consume '.'
        for isDigit(l.peek()) { l.advance() }
    }
    return l.makeToken(NUMBER)
}
```

This is literally a sliding window: for each token, we slide `start` to `current`, then expand `current` rightward until we complete the token.

### The Parser: Divide and Conquer (Recursive Descent)

Recursive descent parsing is divide and conquer applied to grammar rules. Each grammar rule is a function that "divides" by calling sub-rules:

```go
// parser/parser.go (simplified excerpt)
// parseExpr → parseTerm → parseFactor → parsePrimary
// This is divide and conquer: each function handles one level of the grammar,
// calls child functions for sub-expressions, then combines results into an AST node.

func (p *Parser) parseExpr() ASTNode {
    return p.parseAddition() // entry point
}

func (p *Parser) parseAddition() ASTNode {
    left := p.parseMultiplication() // conquer left
    for p.match(PLUS, MINUS) {
        op := p.previous()
        right := p.parseMultiplication() // conquer right
        left = BinaryNode{Op: op, Left: left, Right: right} // combine
    }
    return left
}

func (p *Parser) parseMultiplication() ASTNode {
    left := p.parseUnary() // conquer left
    for p.match(STAR, SLASH) {
        op := p.previous()
        right := p.parseUnary() // conquer right
        left = BinaryNode{Op: op, Left: left, Right: right} // combine
    }
    return left
}

func (p *Parser) parseUnary() ASTNode {
    if p.match(MINUS, BANG) {
        op := p.previous()
        right := p.parseUnary() // recursive: D&C!
        return UnaryNode{Op: op, Right: right}
    }
    return p.parsePrimary()
}
```

Each `parse*` function is a divide and conquer step: it handles one grammar rule, calls sub-rule functions (the "divide"), and combines the returned AST nodes.

### Constant Folding: Post-Order Tree Walk

Constant folding (replacing `3 + 4` with `7` at compile time) uses a post-order walk — evaluate children before parents:

```go
// optimizer/constant_fold.go
func (o *Optimizer) foldConstants(node ASTNode) ASTNode {
    switch n := node.(type) {
    case *BinaryNode:
        // POST-ORDER: fold children first, then try to fold this node
        n.Left = o.foldConstants(n.Left)   // fold left child first
        n.Right = o.foldConstants(n.Right) // fold right child first

        // Now check if both children became constant literals
        leftLit, leftOk := n.Left.(*IntLiteral)
        rightLit, rightOk := n.Right.(*IntLiteral)

        if leftOk && rightOk {
            // Combine: evaluate at compile time
            switch n.Op.Type {
            case PLUS:  return &IntLiteral{Value: leftLit.Value + rightLit.Value}
            case MINUS: return &IntLiteral{Value: leftLit.Value - rightLit.Value}
            case STAR:  return &IntLiteral{Value: leftLit.Value * rightLit.Value}
            case SLASH:
                if rightLit.Value != 0 {
                    return &IntLiteral{Value: leftLit.Value / rightLit.Value}
                }
            }
        }
        return n // cannot fold, return as-is

    case *IntLiteral:
        return n // already a constant, nothing to fold

    default:
        return node
    }
}
```

### Register Allocation: Greedy Linear Scan

When generating code for physical registers, we use a greedy algorithm: scan through instructions in order, and greedily assign variables to any available register. When we run out of registers, we "spill" the variable least likely to be used soon to memory.

```go
// codegen/regalloc.go (simplified linear scan)
type LiveInterval struct {
    VarName string
    Start   int // first instruction that uses this variable
    End     int // last instruction that uses this variable
}

func LinearScanRegisterAlloc(intervals []LiveInterval, numRegisters int) map[string]int {
    // Sort by start point (greedy: process in order)
    sort.Slice(intervals, func(i, j int) bool {
        return intervals[i].Start < intervals[j].Start
    })

    assignment := make(map[string]int)    // var → register number
    freeRegs := []int{}
    for r := 0; r < numRegisters; r++ {
        freeRegs = append(freeRegs, r)
    }

    active := []LiveInterval{} // intervals currently in a register

    for _, interval := range intervals {
        // Expire intervals that ended before this one starts
        newActive := []LiveInterval{}
        for _, a := range active {
            if a.End < interval.Start {
                // Free the register
                freeRegs = append(freeRegs, assignment[a.VarName])
            } else {
                newActive = append(newActive, a)
            }
        }
        active = newActive

        if len(freeRegs) > 0 {
            // Greedy: assign any free register
            reg := freeRegs[len(freeRegs)-1]
            freeRegs = freeRegs[:len(freeRegs)-1]
            assignment[interval.VarName] = reg
            active = append(active, interval)
        } else {
            // Spill: no free registers — mark as spilled (-1)
            assignment[interval.VarName] = -1
        }
    }

    return assignment
}
```

This is the greedy pattern: make the locally best decision (assign any free register, spill the worst candidate) without backtracking.

---

## Exercises

1. **Merge sort + count inversions**: An inversion is a pair (i, j) where i < j but arr[i] > arr[j]. Modify merge sort to count inversions in O(n log n). (Hint: during the merge step, when you pick from the right half, all remaining left-half elements form inversions with it.)

2. **Minimum interval to include each query**: Given a list of intervals and a list of queries, for each query find the length of the smallest interval that contains it. (Hint: sort intervals and queries, use a greedy approach with a min-heap.)

3. **Sliding window maximum**: Given an array and window size k, find the maximum value in every window position in O(n). (Hint: use a monotonic deque — a combination of sliding window and monotonic structure.)

4. **Prefix sum for count of range sum**: Count the number of subarrays with sum in range [lower, upper]. (Hint: use prefix sums + a sorted structure to count pairs.)

5. **Bitset for type set**: The Astra type checker tracks sets of types that a variable could hold (for union types). Implement an efficient TypeSet using bit manipulation that supports union, intersection, and membership testing in O(1).

6. **Compiler phase analysis**: Pick one more phase of the Astra compiler (e.g., escape analysis, inlining, or loop unrolling) and analyze which algorithm design pattern it uses. Write out the reasoning clearly, as we did in the Milestone section.

---

## Summary Table

| Pattern | Key Idea | Time | Best Used When |
|---|---|---|---|
| Divide and Conquer | Split, conquer, combine | O(n log n) typical | Problem structure mirrors recursion |
| Greedy | Locally optimal choice at each step | O(n log n) typical | Greedy choice + optimal substructure proven |
| Two Pointers | Two indices scanning array | O(n) | Sorted array, find pair, remove duplicates |
| Sliding Window | Expand right, shrink left | O(n) | Contiguous subarray/substring problems |
| Monotonic Stack | Maintain sorted order in stack | O(n) | "Next greater/smaller element" queries |
| Prefix Sums | Precompute cumulative sums | O(n) build, O(1) query | Range sum queries, subarray sum counting |
| Bit Manipulation | Operate on binary representation | O(1) | Set membership, fast arithmetic, flags |
| Systematic approach | 7-step problem-solving framework | N/A | Any algorithm problem |
