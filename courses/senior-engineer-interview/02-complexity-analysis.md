# Chapter 02: Complexity Analysis — Big O, Time & Space

Complexity analysis is the language of coding interviews. Every time you write code in an interview, the interviewer expects you to characterize how its performance scales. If you cannot do this fluently and automatically, it signals a gap in fundamentals. This chapter makes complexity analysis second nature.

## Table of Contents

1. [Why Complexity Analysis Matters](#1-why-complexity-analysis-matters)
2. [Big O Notation — The Formal Definition Made Simple](#2-big-o-notation--the-formal-definition-made-simple)
3. [The Complexity Classes You Must Know Cold](#3-the-complexity-classes-you-must-know-cold)
4. [Analyzing Code in Go](#4-analyzing-code-in-go)
5. [Space Complexity](#5-space-complexity)
6. [Amortized Analysis](#6-amortized-analysis)
7. [Common Traps and Mistakes](#7-common-traps-and-mistakes)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Complexity Analysis Matters

Imagine you write a function that finds duplicates in a slice. You test it on 100 elements — it runs in 1 millisecond. Works great. Your interviewer then says: "What if the input has 10 million elements?"

If you cannot answer that question precisely, you cannot reason about production systems. Real systems handle millions of requests, process billions of records, and run with tight latency budgets. The engineer who understands complexity is the one who can design systems that do not fall over under load.

Big O notation is just a precise way of answering the question: **as input grows, how does performance grow?**

---

## 2. Big O Notation — The Formal Definition Made Simple

**Formal definition:** A function f(n) is O(g(n)) if there exist positive constants c and n₀ such that for all n > n₀: f(n) ≤ c · g(n).

**In plain English:** Big O describes the upper bound of how an algorithm's run time or memory usage grows relative to the input size n. We drop constants and lower-order terms because we care about growth rate, not exact values.

### Why We Drop Constants

f(n) = 5n and f(n) = n are both O(n). Why? Because if you double the input, both run twice as long. The constant 5 does not change this relationship. At large enough n, constants become irrelevant — the growth rate dominates.

f(n) = n² + 100n + 50 is O(n²). The n² term grows so much faster than 100n or 50 that for large inputs, only the n² term matters.

```go
// Example: both of these are O(n), even though one does more work
func sumSliceA(nums []int) int {
    total := 0
    for _, n := range nums { // one pass through nums
        total += n
    }
    return total
}

func sumSliceB(nums []int) int {
    total := 0
    for _, n := range nums { // first pass
        total += n
    }
    for _, n := range nums { // second pass (same size)
        total += n * 2
    }
    return total
}
// sumSliceB runs 2n iterations. But O(2n) = O(n) — we drop the constant.
```

### Big O Describes the Worst Case

Unless stated otherwise, Big O is worst-case analysis. For a linear search in a slice of n elements, the worst case is when the element is last — you scan all n elements. So linear search is O(n).

Some algorithms have different best/average/worst cases, and interviewers will sometimes ask about all three.

---

## 3. The Complexity Classes You Must Know Cold

From fastest to slowest, here are the complexity classes you will encounter in interviews:

```
O(1)        Constant    — does not grow with input
O(log n)    Logarithmic — doubles input = one more step
O(n)        Linear      — doubles input = double the work
O(n log n)  Linearithmic — sorting, merge sort, heap sort
O(n²)       Quadratic   — nested loops over same data
O(n³)       Cubic       — triple nested loops
O(2^n)      Exponential — every element makes a binary choice
O(n!)       Factorial   — every possible ordering of n elements
```

### Visual Growth at Different Input Sizes

```
n =        10       100      1,000     1,000,000

O(1)        1         1          1             1
O(log n)    3         7         10            20
O(n)       10       100      1,000     1,000,000
O(n log n) 33       664     10,000    20,000,000
O(n²)     100    10,000  1,000,000  10^12 (1 trillion!)
O(2^n)   1,024  10^30    10^301    ASTRONOMICAL
```

This table explains why O(n²) algorithms are fine for n=100 but catastrophic for n=1,000,000.

### O(1) — Constant Time

```go
// Hash map lookup: O(1) average
func getUser(users map[int]string, id int) string {
    return users[id] // one operation regardless of map size
}

// Array index access: O(1)
func firstElement(nums []int) int {
    return nums[0] // directly access memory address
}
```

### O(log n) — Logarithmic

The hallmark of logarithmic complexity: the input is halved each step.

```go
// Binary search: O(log n)
// Each iteration cuts the search space in half.
func binarySearch(nums []int, target int) int {
    left, right := 0, len(nums)-1

    for left <= right {
        mid := left + (right-left)/2 // avoids integer overflow vs (left+right)/2

        if nums[mid] == target {
            return mid
        } else if nums[mid] < target {
            left = mid + 1  // target must be in right half
        } else {
            right = mid - 1 // target must be in left half
        }
    }
    return -1
}
// With 1,000,000 elements: at most log₂(1,000,000) ≈ 20 iterations.
// With 1,000,000,000 elements: at most 30 iterations. That is logarithmic growth.
```

### O(n) — Linear

```go
// Find maximum: O(n) — must look at every element once
func findMax(nums []int) int {
    max := nums[0]
    for _, n := range nums {
        if n > max {
            max = n
        }
    }
    return max
}
```

### O(n log n) — Linearithmic

Most efficient comparison-based sorting algorithms are O(n log n). This is the minimum complexity for any comparison sort (proven by information theory).

```go
import "sort"

// Go's sort.Slice uses a hybrid of quicksort/heapsort/insertion sort
// Average case: O(n log n)
sort.Slice(nums, func(i, j int) bool {
    return nums[i] < nums[j]
})
```

### O(n²) — Quadratic

The sign of nested loops over the same data.

```go
// Bubble sort: O(n²) — avoid this in interviews unless asked specifically
func bubbleSort(nums []int) {
    n := len(nums)
    for i := 0; i < n; i++ {       // outer loop: n iterations
        for j := 0; j < n-i-1; j++ { // inner loop: up to n iterations
            if nums[j] > nums[j+1] {
                nums[j], nums[j+1] = nums[j+1], nums[j]
            }
        }
    }
}
// Total operations: n * n = n². For n=1000, that's 1,000,000 operations.
```

### O(2^n) — Exponential

Typical of problems that explore all subsets.

```go
// Generate all subsets of a slice: O(2^n)
// For every element, it is either in the subset or not: 2^n possibilities
func subsets(nums []int) [][]int {
    result := [][]int{{}}
    for _, num := range nums {
        newSubsets := make([][]int, len(result))
        for i, s := range result {
            subset := make([]int, len(s)+1)
            copy(subset, s)
            subset[len(s)] = num
            newSubsets[i] = subset
        }
        result = append(result, newSubsets...)
    }
    return result
}
// For n=20: 2^20 = 1,048,576 subsets. Manageable.
// For n=40: 2^40 = 1,099,511,627,776. Not manageable.
```

---

## 4. Analyzing Code in Go

### Rules for Analysis

**Rule 1: Sequential code — add complexities**

```go
func example(nums []int) {
    sort.Slice(nums, ...)   // O(n log n)
    findMax(nums)           // O(n)
    // Total: O(n log n + n) = O(n log n) — keep the dominant term
}
```

**Rule 2: Nested loops — multiply complexities**

```go
func allPairs(nums []int) {
    for i := 0; i < len(nums); i++ {        // O(n)
        for j := i+1; j < len(nums); j++ {  // O(n)
            fmt.Println(nums[i], nums[j])
        }
    }
    // Total: O(n * n) = O(n²)
}
```

**Rule 3: Recursion — depends on call depth and branching**

```go
// Fibonacci (naive): O(2^n) — each call branches into 2 more calls, depth n
func fib(n int) int {
    if n <= 1 { return n }
    return fib(n-1) + fib(n-2) // T(n) = T(n-1) + T(n-2) → O(2^n)
}

// Fibonacci (memoized): O(n) — each subproblem solved once
func fibMemo(n int, memo map[int]int) int {
    if n <= 1 { return n }
    if v, ok := memo[n]; ok { return v }
    memo[n] = fibMemo(n-1, memo) + fibMemo(n-2, memo)
    return memo[n]
}
```

### Analyzing Map Operations

In Go, map operations (get, set, delete) are O(1) average but O(n) worst case (hash collisions). In interviews, always say "O(1) average" and mention the worst case if asked.

```go
// Building a frequency map: O(n)
func frequency(words []string) map[string]int {
    freq := make(map[string]int)
    for _, w := range words { // O(n) iterations
        freq[w]++              // O(1) per operation
    }
    return freq
    // Total: O(n * 1) = O(n)
}
```

### The Two-Pointer Pattern: Better Than O(n²)

A classic interview scenario: naive O(n²) vs. optimal O(n).

```go
// NAIVE: Find two numbers that sum to target — O(n²)
func twoSumNaive(nums []int, target int) (int, int) {
    for i := 0; i < len(nums); i++ {
        for j := i + 1; j < len(nums); j++ { // nested loop = O(n²)
            if nums[i]+nums[j] == target {
                return i, j
            }
        }
    }
    return -1, -1
}

// OPTIMAL: Use a hash map — O(n) time, O(n) space
func twoSumOptimal(nums []int, target int) (int, int) {
    seen := make(map[int]int) // value -> index
    for i, n := range nums {
        complement := target - n
        if j, ok := seen[complement]; ok {
            return j, i // found it!
        }
        seen[n] = i // remember this value
    }
    return -1, -1
}
// One pass, O(n) time. Trading space for time — a classic pattern.
```

---

## 5. Space Complexity

Space complexity measures how much extra memory your algorithm uses relative to input size. Input space itself is usually not counted — we count auxiliary (extra) space.

```go
// O(1) space — only a few variables, no matter how large nums is
func sumArray(nums []int) int {
    sum := 0                   // one integer variable
    for _, n := range nums {
        sum += n
    }
    return sum
}

// O(n) space — creates a new slice of same size as input
func doubled(nums []int) []int {
    result := make([]int, len(nums)) // allocates n integers
    for i, n := range nums {
        result[i] = n * 2
    }
    return result
}

// O(n) space — recursion stack depth proportional to n
func factorial(n int) int {
    if n <= 1 { return 1 }
    return n * factorial(n-1) // n recursive calls on the call stack
}

// O(log n) space — binary search uses recursion depth of log n
func binarySearchRecursive(nums []int, target, left, right int) int {
    if left > right { return -1 }
    mid := left + (right-left)/2
    if nums[mid] == target { return mid }
    if nums[mid] < target {
        return binarySearchRecursive(nums, target, mid+1, right)
    }
    return binarySearchRecursive(nums, target, left, mid-1)
}
// Call stack depth is log n — O(log n) space
```

### Space vs. Time Tradeoff

One of the most important engineering judgments is deciding when to trade space for time.

```go
// Example: caching vs recomputing
// Option A: Recompute every time — O(1) space, O(n) time per call
// Option B: Cache results — O(n) space, O(1) time per call after warming up
```

In interviews, always discuss this tradeoff explicitly when you use a hash map or cache.

---

## 6. Amortized Analysis

Amortized analysis tells you the average cost of an operation over a sequence of operations, even if individual operations occasionally cost more.

The canonical example is a dynamic array (Go's slice with `append`).

```go
// When you append to a Go slice:
nums := []int{}
for i := 0; i < 100; i++ {
    nums = append(nums, i) // sometimes O(1), sometimes O(n)
}
```

When the slice is full and you append, Go allocates a new, larger array and copies everything — that's O(n). But this happens rarely. Most appends are O(1). The amortized cost is O(1) because:

- Suppose the slice doubles each time it resizes.
- Resizes happen at sizes 1, 2, 4, 8, 16, ...
- Total copy work for n elements = 1 + 2 + 4 + 8 + ... + n = 2n operations.
- Spread over n appends: 2n / n = 2 operations per append on average = O(1) amortized.

**In an interview:** If asked about `append` or a stack's `push` operation, say "O(1) amortized." If asked why, explain the doubling strategy above.

---

## 7. Common Traps and Mistakes

### Trap 1: Forgetting about the inner loop's range

```go
// What is the complexity?
for i := 0; i < n; i++ {
    for j := i; j < n; j++ { // j starts at i, not 0!
        // ...
    }
}
// At first glance you might say O(n²).
// Correct: it IS O(n²) — (n + (n-1) + ... + 1) = n(n+1)/2 ≈ n²/2 → O(n²)
// The constant 1/2 is dropped.
```

### Trap 2: Recursion with memoization

Without memoization, recursive Fibonacci is O(2^n). With memoization, it is O(n) time and O(n) space. This difference is enormous and interviewers check whether you know it.

### Trap 3: String concatenation in a loop

```go
// This looks O(n) but is actually O(n²) in many languages
result := ""
for _, word := range words {
    result += word // creates a new string each iteration!
}
// Each concatenation copies the existing string, making it O(1+2+3+...+n) = O(n²)

// Correct O(n) approach:
import "strings"
result = strings.Join(words, "")
// Or use strings.Builder
```

In Go, string concatenation with `+` creates a new string each time. `strings.Builder` avoids this.

### Trap 4: Multiple inputs require multiple variables

```go
// What is the complexity of this function?
func intersect(a, b []int) []int {
    set := make(map[int]bool)
    for _, v := range a { set[v] = true }  // O(a)
    result := []int{}
    for _, v := range b {                   // O(b)
        if set[v] { result = append(result, v) }
    }
    return result
}
// Answer: O(a + b) — not O(n)!
// Always use separate variables when your algorithm has multiple input sizes.
```

---

## 8. Interview Questions & Model Answers

### Q: What is the time complexity of your solution?

After coding, always state it proactively: "This solution is O(n) time and O(n) space. The single loop over the input gives us O(n) time, and the hash map can hold up to n entries in the worst case."

### Q: Can you do better?

This is a probe. Have a prepared response: "The current approach is O(n) time. For this type of problem, we can sometimes reduce space to O(1) if we are allowed to modify the input, or we could achieve O(n log n) time with O(1) space using a two-pointer approach on a sorted array. Would you like me to explore that tradeoff?"

### Q: What is the worst-case complexity of a hash map lookup?

"O(n) worst case due to hash collisions, but O(1) average case with a good hash function. In practice for interviews, we assume O(1). In production, we would choose hash functions designed to avoid collision attacks."

### Q: What is amortized O(1)?

"Amortized O(1) means the operation is O(1) on average across a sequence of operations, even if individual operations are occasionally more expensive. The classic example is Go's `append` — most appends are O(1) but a resize copies all elements. Because resizes happen exponentially less often, the amortized cost per append is O(1)."

### Q: Compare O(n log n) and O(n²) for n = 1,000,000

"O(n log n) for n=1,000,000 is about 20,000,000 operations. O(n²) is 10^12 operations — one trillion. At 10^9 operations per second, O(n log n) takes 20 milliseconds. O(n²) takes over 15 minutes. This is why O(n log n) sorting algorithms are essential."

---

## Summary

- Big O describes the upper bound of how an algorithm's performance grows with input size.
- Drop constants and lower-order terms — only the dominant term matters.
- Key classes: O(1) → O(log n) → O(n) → O(n log n) → O(n²) → O(2^n).
- Nested loops multiply complexities; sequential code adds them.
- Space complexity counts auxiliary memory, not the input itself.
- Amortized O(1) means O(1) on average across many operations.
- Common traps: string concatenation, multiple input sizes, memoization changes exponential to polynomial.

---

## Exercises

### Easy

1. Analyze the time and space complexity of this function:
   ```go
   func containsDuplicate(nums []int) bool {
       seen := make(map[int]bool)
       for _, n := range nums {
           if seen[n] { return true }
           seen[n] = true
       }
       return false
   }
   ```

2. Why is `O(n/2)` the same as `O(n)`? Write a one-paragraph explanation in your own words.

3. A function processes all pairs in a slice. Write the function and state its complexity.

### Medium

4. You have an algorithm with the following recurrence: T(n) = 2T(n/2) + O(n). Use the Master Theorem to determine its complexity. (Answer: O(n log n) — this is merge sort's complexity.)

5. Write a Go function that finds the intersection of three sorted arrays in O(n) time where n is the total number of elements. What is the space complexity?

6. A hash map has O(1) average lookup. Describe a scenario where a real production system could experience O(n) hash map lookups and how you would prevent it.

### Hard

7. Analyze the time and space complexity of BFS on a graph with V vertices and E edges. Why is it O(V + E) and not O(V²)?

8. A Go slice `s` starts with capacity 1 and doubles on each resize. After 2^k appends, what is the total number of element copy operations? Prove that amortized cost per append is O(1).

9. Design an O(n) time, O(1) space algorithm to find the majority element (appears more than n/2 times) in a slice. Why does the Boyer-Moore Voting Algorithm work?
