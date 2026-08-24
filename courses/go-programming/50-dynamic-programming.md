# Chapter 36: Dynamic Programming

Dynamic programming (DP) is the art of remembering answers to avoid recomputation. It converts exponential brute-force into polynomial solutions by identifying overlapping subproblems and optimal substructure. The hard part isn't coding — it's formulating the recurrence. Once you see the pattern, the code writes itself.

## Table of Contents

1. [What is DP?](#1-what-is-dp)
2. [Memoization (Top-Down)](#2-memoization-top-down)
3. [Tabulation (Bottom-Up)](#3-tabulation-bottom-up)
4. [Classic 1D Problems](#4-classic-1d-problems)
5. [Classic 2D Problems](#5-classic-2d-problems)
6. [Interval and Partition DP](#6-interval-and-partition-dp)
7. [DP Patterns Reference](#7-dp-patterns-reference)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What is DP?

**Two required properties:**

1. **Optimal substructure**: the optimal solution contains optimal solutions to subproblems
2. **Overlapping subproblems**: the same subproblems are solved multiple times

**DP vs Greedy**: Greedy makes locally optimal choices and never looks back. DP tries all possibilities and picks the best.

**DP vs Divide & Conquer**: D&C subproblems don't overlap (merge sort partitions non-overlapping halves). DP subproblems do overlap (fibonacci(5) and fibonacci(4) both need fibonacci(3)).

**The process:**
1. Define the state: what does `dp[i]` (or `dp[i][j]`) represent?
2. Write the recurrence: how does `dp[i]` relate to smaller states?
3. Identify base cases
4. Determine computation order (ensure dependencies are computed first)
5. Extract the answer from the dp table

---

## 2. Memoization (Top-Down)

Memoization = recursion + caching. Write the natural recursive solution and add a cache.

```go
// Fibonacci — the classic DP introduction
// Naive recursion: O(2^n)
func fibNaive(n int) int {
    if n <= 1 { return n }
    return fibNaive(n-1) + fibNaive(n-2)
}

// Memoized: O(n) time, O(n) space
func FibMemo(n int) int {
    memo := make(map[int]int)
    var fib func(int) int
    fib = func(n int) int {
        if n <= 1 { return n }
        if v, ok := memo[n]; ok { return v }
        memo[n] = fib(n-1) + fib(n-2)
        return memo[n]
    }
    return fib(n)
}

// Slice memo is faster than map for integer indices:
func FibMemoSlice(n int) int {
    memo := make([]int, n+1)
    for i := range memo { memo[i] = -1 }
    var fib func(int) int
    fib = func(n int) int {
        if n <= 1 { return n }
        if memo[n] != -1 { return memo[n] }
        memo[n] = fib(n-1) + fib(n-2)
        return memo[n]
    }
    return fib(n)
}
```

---

## 3. Tabulation (Bottom-Up)

Tabulation fills in the DP table iteratively, starting from base cases.

```go
// Fibonacci tabulation: O(n) time, O(n) space
func FibTab(n int) int {
    if n <= 1 { return n }
    dp := make([]int, n+1)
    dp[0], dp[1] = 0, 1
    for i := 2; i <= n; i++ {
        dp[i] = dp[i-1] + dp[i-2]
    }
    return dp[n]
}

// Fibonacci space-optimized: O(1) space
func FibOptimal(n int) int {
    if n <= 1 { return n }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

**Memoization vs Tabulation:**
| | Memoization | Tabulation |
|--|-------------|------------|
| Order | Top-down (natural) | Bottom-up (explicit) |
| Base cases | Handled by recursion | Explicit initialization |
| Stack overflow | Risk for deep recursion | None |
| Only needed states | Computed lazily | All states computed |
| Cache | Map or slice with sentinel | DP array |

---

## 4. Classic 1D Problems

### Climbing Stairs
```go
// ClimbingStairs: n stairs, can climb 1 or 2 at a time.
// dp[i] = number of ways to reach step i
func ClimbStairs(n int) int {
    if n <= 2 { return n }
    prev2, prev1 := 1, 2
    for i := 3; i <= n; i++ {
        prev2, prev1 = prev1, prev2+prev1
    }
    return prev1
}
```

### House Robber
```go
// HouseRobber: Rob max money from houses; can't rob adjacent houses.
// dp[i] = max money robbing houses 0..i
// dp[i] = max(dp[i-1], dp[i-2] + nums[i])
func Rob(nums []int) int {
    if len(nums) == 0 { return 0 }
    if len(nums) == 1 { return nums[0] }

    prev2, prev1 := nums[0], max(nums[0], nums[1])
    for i := 2; i < len(nums); i++ {
        prev2, prev1 = prev1, max(prev1, prev2+nums[i])
    }
    return prev1
}

func max(a, b int) int {
    if a > b { return a }
    return b
}
```

### Coin Change — Minimum Coins
```go
// CoinChange: minimum number of coins summing to amount.
// dp[i] = min coins to make sum i
// dp[i] = min over all coins c: 1 + dp[i-c]  if i >= c
func CoinChange(coins []int, amount int) int {
    const INF = amount + 1
    dp := make([]int, amount+1)
    for i := range dp { dp[i] = INF }
    dp[0] = 0

    for i := 1; i <= amount; i++ {
        for _, c := range coins {
            if c <= i && dp[i-c]+1 < dp[i] {
                dp[i] = dp[i-c] + 1
            }
        }
    }
    if dp[amount] == INF { return -1 }
    return dp[amount]
}
```

### Coin Change — Number of Ways
```go
// CoinWays: count distinct combinations that sum to amount.
// Outer loop over coins, inner over amounts — avoids counting permutations.
func CoinWays(coins []int, amount int) int {
    dp := make([]int, amount+1)
    dp[0] = 1  // One way to make 0: use no coins

    for _, c := range coins {
        for i := c; i <= amount; i++ {
            dp[i] += dp[i-c]
        }
    }
    return dp[amount]
}
// Key insight: outer=items, inner=capacity → combinations (not permutations)
// Swap the loops: outer=capacity, inner=items → permutations
```

### Longest Increasing Subsequence (LIS)
```go
// LIS O(n²): dp[i] = length of LIS ending at index i
func LIS_N2(nums []int) int {
    n := len(nums)
    if n == 0 { return 0 }
    dp := make([]int, n)
    for i := range dp { dp[i] = 1 }

    best := 1
    for i := 1; i < n; i++ {
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] && dp[j]+1 > dp[i] {
                dp[i] = dp[j] + 1
            }
        }
        if dp[i] > best { best = dp[i] }
    }
    return best
}

// LIS O(n log n): patience sorting
// tails[i] = smallest tail element for increasing subseq of length i+1
func LIS(nums []int) int {
    tails := []int{}

    for _, n := range nums {
        lo, hi := 0, len(tails)
        for lo < hi {
            mid := lo + (hi-lo)/2
            if tails[mid] < n { lo = mid + 1 } else { hi = mid }
        }
        if lo == len(tails) {
            tails = append(tails, n)
        } else {
            tails[lo] = n
        }
    }
    return len(tails)
}
```

---

## 5. Classic 2D Problems

### Longest Common Subsequence (LCS)
```go
// LCS: length of longest common subsequence of two strings.
// dp[i][j] = LCS of s[:i] and t[:j]
// Recurrence:
//   if s[i-1] == t[j-1]: dp[i][j] = 1 + dp[i-1][j-1]
//   else:                 dp[i][j] = max(dp[i-1][j], dp[i][j-1])
func LCS(s, t string) int {
    m, n := len(s), len(t)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if s[i-1] == t[j-1] {
                dp[i][j] = 1 + dp[i-1][j-1]
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}

// Reconstruct the actual LCS string:
func LCSString(s, t string) string {
    m, n := len(s), len(t)
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if s[i-1] == t[j-1] {
                dp[i][j] = 1 + dp[i-1][j-1]
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }

    // Backtrack:
    result := []byte{}
    i, j := m, n
    for i > 0 && j > 0 {
        if s[i-1] == t[j-1] {
            result = append(result, s[i-1])
            i--; j--
        } else if dp[i-1][j] > dp[i][j-1] {
            i--
        } else {
            j--
        }
    }
    // Reverse:
    for lo, hi := 0, len(result)-1; lo < hi; lo, hi = lo+1, hi-1 {
        result[lo], result[hi] = result[hi], result[lo]
    }
    return string(result)
}
```

### Edit Distance (Levenshtein)
```go
// EditDistance: min operations (insert, delete, replace) to convert s to t.
// dp[i][j] = edit distance between s[:i] and t[:j]
func EditDistance(s, t string) int {
    m, n := len(s), len(t)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
        dp[i][0] = i  // Delete all chars in s[:i]
    }
    for j := 0; j <= n; j++ {
        dp[0][j] = j  // Insert all chars of t[:j]
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if s[i-1] == t[j-1] {
                dp[i][j] = dp[i-1][j-1]  // No op
            } else {
                dp[i][j] = 1 + min3(
                    dp[i-1][j],   // Delete from s
                    dp[i][j-1],   // Insert into s
                    dp[i-1][j-1], // Replace
                )
            }
        }
    }
    return dp[m][n]
}

func min3(a, b, c int) int {
    if a < b { if a < c { return a }; return c }
    if b < c { return b }
    return c
}
```

### 0/1 Knapsack
```go
// Knapsack01: max value fitting in capacity W.
// dp[i][w] = max value using items 0..i-1 with weight capacity w
// Either skip item i: dp[i-1][w]
// Or take item i (if fits): values[i-1] + dp[i-1][w-weights[i-1]]
func Knapsack01(weights, values []int, capacity int) int {
    n := len(weights)
    dp := make([][]int, n+1)
    for i := range dp { dp[i] = make([]int, capacity+1) }

    for i := 1; i <= n; i++ {
        for w := 0; w <= capacity; w++ {
            dp[i][w] = dp[i-1][w]  // Skip
            if weights[i-1] <= w {
                if v := values[i-1] + dp[i-1][w-weights[i-1]]; v > dp[i][w] {
                    dp[i][w] = v   // Take
                }
            }
        }
    }
    return dp[n][capacity]
}

// Space-optimized to O(capacity) — iterate weight BACKWARDS:
func Knapsack01Optimized(weights, values []int, capacity int) int {
    dp := make([]int, capacity+1)

    for i := range weights {
        for w := capacity; w >= weights[i]; w-- {
            if v := dp[w-weights[i]] + values[i]; v > dp[w] {
                dp[w] = v
            }
        }
    }
    return dp[capacity]
}
// Key: iterate backwards to avoid using same item twice (0/1 constraint)
```

### Unique Paths (Grid DP)
```go
// UniquePaths: robot moves from top-left to bottom-right (only right/down).
// dp[i][j] = paths to reach (i,j) = dp[i-1][j] + dp[i][j-1]
func UniquePaths(m, n int) int {
    dp := make([][]int, m)
    for i := range dp {
        dp[i] = make([]int, n)
        dp[i][0] = 1  // Only one way to reach any cell in column 0
    }
    for j := 0; j < n; j++ { dp[0][j] = 1 }  // Only one way in row 0

    for i := 1; i < m; i++ {
        for j := 1; j < n; j++ {
            dp[i][j] = dp[i-1][j] + dp[i][j-1]
        }
    }
    return dp[m-1][n-1]
}
```

---

## 6. Interval and Partition DP

Interval DP processes subarrays of increasing length.

### Matrix Chain Multiplication
```go
// MatrixChain: find the minimum cost to multiply a chain of matrices.
// dp[i][j] = min operations to multiply matrices i through j
// Split at k: dp[i][k] + dp[k+1][j] + dims[i]*dims[k+1]*dims[j+1]
func MatrixChain(dims []int) int {
    n := len(dims) - 1  // Number of matrices
    dp := make([][]int, n)
    for i := range dp { dp[i] = make([]int, n) }

    // Length of chain from 2 to n:
    for length := 2; length <= n; length++ {
        for i := 0; i <= n-length; i++ {
            j := i + length - 1
            dp[i][j] = 1<<62
            for k := i; k < j; k++ {
                cost := dp[i][k] + dp[k+1][j] + dims[i]*dims[k+1]*dims[j+1]
                if cost < dp[i][j] { dp[i][j] = cost }
            }
        }
    }
    return dp[0][n-1]
}
```

### Longest Palindromic Subsequence
```go
// LPS: length of longest palindromic subsequence.
// dp[i][j] = LPS in s[i..j]
func LPS(s string) int {
    n := len(s)
    dp := make([][]int, n)
    for i := range dp {
        dp[i] = make([]int, n)
        dp[i][i] = 1  // Single char is a palindrome of length 1
    }

    for length := 2; length <= n; length++ {
        for i := 0; i <= n-length; i++ {
            j := i + length - 1
            if s[i] == s[j] {
                if length == 2 {
                    dp[i][j] = 2
                } else {
                    dp[i][j] = 2 + dp[i+1][j-1]
                }
            } else {
                dp[i][j] = max(dp[i+1][j], dp[i][j-1])
            }
        }
    }
    return dp[0][n-1]
}
```

---

## 7. DP Patterns Reference

| Problem Type | State | Recurrence |
|---|---|---|
| 1D sequence | `dp[i]` = answer for prefix `[0..i]` | `dp[i] = f(dp[i-1], dp[i-2], ...)` |
| 2D grid | `dp[i][j]` = answer at cell `(i,j)` | `dp[i][j] = f(dp[i-1][j], dp[i][j-1])` |
| Two sequences | `dp[i][j]` = answer for `s[:i]` and `t[:j]` | `dp[i][j] = f(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])` |
| Knapsack | `dp[i][w]` = max value with i items, capacity w | `dp[i][w] = max(skip, take)` |
| Interval | `dp[i][j]` = answer for `arr[i..j]` | `dp[i][j] = min over k: dp[i][k] + dp[k+1][j] + cost` |
| Bitmask | `dp[mask]` = answer for subset `mask` | `dp[mask] = f(dp[mask ^ (1<<bit)])` |

**Space optimization**: if `dp[i]` only depends on `dp[i-1]`, use two rows (or one row traversed carefully).

---

## Summary

- DP requires **optimal substructure** (optimal solution uses optimal sub-solutions) + **overlapping subproblems**
- **Top-down (memoization)**: natural recursive solution + cache. Use a slice with `-1` sentinel over a map for performance.
- **Bottom-up (tabulation)**: fill table in dependency order from base cases
- **1D patterns**: Fibonacci / climbing stairs (previous 2 states), house robber (skip or take), coin change (min or count), LIS
- **2D patterns**: LCS (diagonal when equal, max of skip otherwise), edit distance (3 operations), knapsack (skip or take), grid paths
- **Interval DP**: process lengths 2..n, split at k, `dp[i][j] = best over all k`
- **Space optimization**: when `dp[i]` only uses `dp[i-1]`, reduce to O(1) or O(n) space

---

## Exercises

### Easy
1. Write `MaxSubarraySum(nums []int) int` using Kadane's algorithm. This is DP: `dp[i] = max(nums[i], dp[i-1]+nums[i])`. Return the maximum sum.
2. Write `MinPathSum(grid [][]int) int` — find minimum sum path from top-left to bottom-right, moving only right or down. `dp[i][j] = grid[i][j] + min(dp[i-1][j], dp[i][j-1])`.
3. Implement the climbing stairs problem but with step sizes {1, 2, 3}. Generalize to an arbitrary set of step sizes passed as a parameter.

### Medium
4. **Longest Palindromic Substring** (not subsequence): `dp[i][j] = true` if `s[i..j]` is a palindrome. `dp[i][j] = s[i]==s[j] && dp[i+1][j-1]`. Find the longest. Time: O(n²). Verify: `"babad"` → `"bab"`, `"racecar"` → `"racecar"`.
5. **Word Break**: Given a string and a word dictionary, can the string be segmented into valid words? `dp[i] = true` if `s[:i]` can be segmented. `dp[i] = any dp[j] && s[j:i] in dict`. Verify: `"leetcode"`, `["leet","code"]` → true.
6. **Partition Equal Subset Sum**: Given an array, can it be partitioned into two subsets with equal sum? This is a 0/1 knapsack variant: `dp[i][s] = can we pick subset from first i items summing to s`. Space-optimize to O(sum). Verify: `[1,5,11,5]` → true (because [1,5,5] and [11]).

### Hard
7. **Burst Balloons**: Given `n` balloons with values, burst all balloons to maximize coins. If you burst balloon `i`, you earn `nums[i-1]*nums[i]*nums[i+1]`. Interval DP: `dp[i][j]` = max coins bursting all balloons in `(i,j)`. Key insight: think of `k` as the LAST balloon to burst. Time: O(n³).
8. **Minimum cost to cut a stick**: Given a stick of length `n` and an array of cut positions, find the minimum total cost to make all cuts. The cost of cutting is the length of the piece being cut. Interval DP: add 0 and n as boundaries, `dp[i][j]` = min cost to make all cuts between position `i` and `j`. Verify: `n=7, cuts=[1,3,4,5]` → 16.
