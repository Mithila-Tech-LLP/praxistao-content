# Chapter 09: Dynamic Programming — Patterns & Practice

Dynamic programming (DP) is the most feared topic in coding interviews. The fear is justified — DP problems require recognizing the right subproblem decomposition, which is not obvious. But there are only five core DP patterns. Learn to recognize each pattern and the code practically writes itself.

## Table of Contents

1. [What Makes a Problem a DP Problem](#1-what-makes-a-problem-a-dp-problem)
2. [Pattern 1: Linear DP](#2-pattern-1-linear-dp)
3. [Pattern 2: Grid DP](#3-pattern-2-grid-dp)
4. [Pattern 3: Interval DP](#4-pattern-3-interval-dp)
5. [Pattern 4: Knapsack DP](#5-pattern-4-knapsack-dp)
6. [Pattern 5: State Machine DP](#6-pattern-5-state-machine-dp)
7. [Top-Down vs Bottom-Up](#7-top-down-vs-bottom-up)
8. [Classic Hard Problems](#8-classic-hard-problems)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What Makes a Problem a DP Problem

A problem is a DP problem if it has:

1. **Optimal substructure:** The optimal solution to the problem contains optimal solutions to subproblems.
2. **Overlapping subproblems:** The same subproblems are solved multiple times in a naive recursive solution.

**How to identify DP in an interview:**
- Problem asks for maximum/minimum/count of something
- You are making choices at each step and need to optimize
- Future choices depend on past choices
- Recursive brute force would recompute the same subproblems

**The DP thought process:**
1. Define what `dp[i]` (or `dp[i][j]`) represents
2. Find the recurrence relation: how does `dp[i]` depend on smaller subproblems?
3. Identify base cases
4. Determine the computation order (bottom-up: smaller → larger)
5. Extract the answer

---

## 2. Pattern 1: Linear DP

The subproblem depends on one previous state. Classic problems: climbing stairs, Fibonacci, house robber.

### Climbing Stairs

```go
// dp[i] = number of ways to reach step i
// dp[i] = dp[i-1] (come from i-1 with one step) + dp[i-2] (come from i-2 with two steps)
func climbStairs(n int) int {
    if n <= 2 { return n }
    prev2, prev1 := 1, 2
    for i := 3; i <= n; i++ {
        curr := prev1 + prev2
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}
// Time: O(n), Space: O(1) — space-optimized (only need last 2 values)
```

### Longest Increasing Subsequence (LIS)

```go
// dp[i] = length of LIS ending at index i
func lengthOfLIS(nums []int) int {
    n := len(nums)
    dp := make([]int, n)
    for i := range dp { dp[i] = 1 } // every element is an LIS of length 1

    maxLen := 1
    for i := 1; i < n; i++ {
        for j := 0; j < i; j++ {
            if nums[j] < nums[i] { // nums[i] can extend the LIS ending at j
                if dp[j]+1 > dp[i] {
                    dp[i] = dp[j] + 1
                }
            }
        }
        if dp[i] > maxLen { maxLen = dp[i] }
    }
    return maxLen
}
// Time: O(n²), Space: O(n)

// Optimal O(n log n) approach using patience sorting (binary search):
func lengthOfLISOpt(nums []int) int {
    tails := []int{} // tails[i] = smallest tail element of LIS with length i+1
    for _, n := range nums {
        // Binary search for the first tail >= n
        lo, hi := 0, len(tails)
        for lo < hi {
            mid := (lo + hi) / 2
            if tails[mid] < n { lo = mid + 1 } else { hi = mid }
        }
        if lo == len(tails) {
            tails = append(tails, n) // n extends the longest LIS
        } else {
            tails[lo] = n // replace to keep tails as small as possible
        }
    }
    return len(tails)
}
```

### House Robber

```go
// dp[i] = max money robbing houses 0..i
// At house i: either rob it (dp[i-2] + nums[i]) or skip it (dp[i-1])
func rob(nums []int) int {
    if len(nums) == 1 { return nums[0] }
    prev2, prev1 := nums[0], max(nums[0], nums[1])
    for i := 2; i < len(nums); i++ {
        curr := max(prev1, prev2+nums[i])
        prev2 = prev1
        prev1 = curr
    }
    return prev1
}
```

### Word Break

```go
// dp[i] = true if s[0:i] can be segmented into words from wordDict
func wordBreak(s string, wordDict []string) bool {
    words := make(map[string]bool)
    for _, w := range wordDict { words[w] = true }

    n := len(s)
    dp := make([]bool, n+1)
    dp[0] = true // empty string is always valid

    for i := 1; i <= n; i++ {
        for j := 0; j < i; j++ {
            if dp[j] && words[s[j:i]] { // s[0:j] is valid and s[j:i] is a word
                dp[i] = true
                break
            }
        }
    }
    return dp[n]
}
// Time: O(n² * m) where m = average word length, Space: O(n)
```

---

## 3. Pattern 2: Grid DP

Subproblem is a cell in a 2D grid. Direction is usually top-left to bottom-right.

### Unique Paths

```go
// dp[r][c] = number of ways to reach cell (r, c)
// dp[r][c] = dp[r-1][c] + dp[r][c-1] (from above + from left)
func uniquePaths(m, n int) int {
    dp := make([][]int, m)
    for i := range dp {
        dp[i] = make([]int, n)
        dp[i][0] = 1 // leftmost column: only one way (go straight down)
    }
    for j := 0; j < n; j++ { dp[0][j] = 1 } // top row: only one way (go right)

    for r := 1; r < m; r++ {
        for c := 1; c < n; c++ {
            dp[r][c] = dp[r-1][c] + dp[r][c-1]
        }
    }
    return dp[m-1][n-1]
}
```

### Minimum Path Sum

```go
// dp[r][c] = minimum path sum from top-left to (r, c)
func minPathSum(grid [][]int) int {
    rows, cols := len(grid), len(grid[0])
    dp := make([][]int, rows)
    for i := range dp { dp[i] = make([]int, cols) }

    dp[0][0] = grid[0][0]
    for c := 1; c < cols; c++ { dp[0][c] = dp[0][c-1] + grid[0][c] }
    for r := 1; r < rows; r++ { dp[r][0] = dp[r-1][0] + grid[r][0] }

    for r := 1; r < rows; r++ {
        for c := 1; c < cols; c++ {
            dp[r][c] = grid[r][c] + min(dp[r-1][c], dp[r][c-1])
        }
    }
    return dp[rows-1][cols-1]
}
```

### Maximal Square

```go
// dp[r][c] = side length of largest square of 1s with bottom-right corner at (r,c)
// If matrix[r][c] == '1': dp[r][c] = min(dp[r-1][c], dp[r][c-1], dp[r-1][c-1]) + 1
func maximalSquare(matrix [][]byte) int {
    rows, cols := len(matrix), len(matrix[0])
    dp := make([][]int, rows)
    for i := range dp { dp[i] = make([]int, cols) }
    maxSide := 0

    for r := 0; r < rows; r++ {
        for c := 0; c < cols; c++ {
            if matrix[r][c] == '1' {
                if r == 0 || c == 0 {
                    dp[r][c] = 1
                } else {
                    dp[r][c] = min(dp[r-1][c], min(dp[r][c-1], dp[r-1][c-1])) + 1
                }
                maxSide = max(maxSide, dp[r][c])
            }
        }
    }
    return maxSide * maxSide
}
```

---

## 4. Pattern 3: Interval DP

Subproblem is a subarray or substring. Typically O(n²) space, O(n³) time.

### Palindrome Substrings

```go
// dp[i][j] = true if s[i:j+1] is a palindrome
func countSubstrings(s string) int {
    n := len(s)
    dp := make([][]bool, n)
    for i := range dp { dp[i] = make([]bool, n) }
    count := 0

    // Process by length: single chars first, then pairs, then longer
    for length := 1; length <= n; length++ {
        for i := 0; i+length-1 < n; i++ {
            j := i + length - 1
            if s[i] == s[j] && (length <= 2 || dp[i+1][j-1]) {
                dp[i][j] = true
                count++
            }
        }
    }
    return count
}
// Alternatively: expand around centers (O(n²) time, O(1) space)
```

### Burst Balloons

```go
// dp[i][j] = max coins from bursting all balloons between i and j (exclusive)
// Key trick: think of k as the LAST balloon to burst in range (i, j)
func maxCoins(nums []int) int {
    // Add boundary balloons of value 1
    balloons := []int{1}
    balloons = append(balloons, nums...)
    balloons = append(balloons, 1)
    n := len(balloons)

    dp := make([][]int, n)
    for i := range dp { dp[i] = make([]int, n) }

    // length = window size (at least 2, since boundaries are not burst)
    for length := 2; length < n; length++ {
        for left := 0; left < n-length; left++ {
            right := left + length
            for k := left + 1; k < right; k++ {
                coins := balloons[left]*balloons[k]*balloons[right]
                dp[left][right] = max(dp[left][right], dp[left][k]+coins+dp[k][right])
            }
        }
    }
    return dp[0][n-1]
}
// Time: O(n³), Space: O(n²)
```

---

## 5. Pattern 4: Knapsack DP

Choose items to maximize/minimize some objective within a constraint.

### 0/1 Knapsack

```go
// dp[i][w] = max value using first i items with capacity w
func knapsack(weights, values []int, capacity int) int {
    n := len(weights)
    dp := make([][]int, n+1)
    for i := range dp { dp[i] = make([]int, capacity+1) }

    for i := 1; i <= n; i++ {
        for w := 0; w <= capacity; w++ {
            dp[i][w] = dp[i-1][w] // don't take item i
            if weights[i-1] <= w {
                dp[i][w] = max(dp[i][w], dp[i-1][w-weights[i-1]]+values[i-1]) // take item i
            }
        }
    }
    return dp[n][capacity]
}
// Space optimization: use 1D dp array, iterate w from right to left
func knapsack1D(weights, values []int, capacity int) int {
    dp := make([]int, capacity+1)
    for i := 0; i < len(weights); i++ {
        for w := capacity; w >= weights[i]; w-- { // MUST iterate backwards for 0/1 knapsack
            dp[w] = max(dp[w], dp[w-weights[i]]+values[i])
        }
    }
    return dp[capacity]
}
```

### Coin Change (Unbounded Knapsack)

```go
// dp[i] = minimum coins to make amount i
func coinChange(coins []int, amount int) int {
    dp := make([]int, amount+1)
    for i := 1; i <= amount; i++ { dp[i] = amount + 1 } // infinity

    for i := 1; i <= amount; i++ {
        for _, c := range coins {
            if c <= i {
                dp[i] = min(dp[i], dp[i-c]+1)
            }
        }
    }
    if dp[amount] > amount { return -1 }
    return dp[amount]
}
// Note: iterate coins in outer loop or amount in outer loop — both work for unbounded
```

### Partition Equal Subset Sum

```go
// Can we partition array into two subsets with equal sum?
// Equivalent to: can we find a subset with sum = total/2?
func canPartition(nums []int) bool {
    total := 0
    for _, n := range nums { total += n }
    if total%2 != 0 { return false }
    target := total / 2

    dp := make([]bool, target+1)
    dp[0] = true

    for _, n := range nums {
        for j := target; j >= n; j-- { // backwards to avoid reusing same element
            dp[j] = dp[j] || dp[j-n]
        }
    }
    return dp[target]
}
```

---

## 6. Pattern 5: State Machine DP

The subproblem has multiple states at each position. Classic: stock buy/sell problems.

### Best Time to Buy and Sell Stocks with Cooldown

```go
// States: HOLDING (have stock), SOLD (just sold, in cooldown), IDLE (no stock, can buy)
// HOLDING: bought stock, waiting to sell
// SOLD: sold today, must wait one day (cooldown)
// IDLE: free to buy
func maxProfitCooldown(prices []int) int {
    holding := -prices[0] // profit if we hold a stock after day 0
    sold := 0             // profit if we just sold today
    idle := 0             // profit with no stock, not in cooldown

    for i := 1; i < len(prices); i++ {
        prevHolding, prevSold, prevIdle := holding, sold, idle
        holding = max(prevHolding, prevIdle-prices[i]) // keep holding OR buy (from idle)
        sold = prevHolding + prices[i]                  // sell today
        idle = max(prevIdle, prevSold)                  // stay idle OR come out of cooldown
    }
    return max(sold, idle)
}
```

---

## 7. Top-Down vs Bottom-Up

```go
// TOP-DOWN (memoization): natural recursion + cache
// Advantage: only computes what is needed, easier to think about
var memo map[[2]int]int

func editDistanceMemo(word1, word2 string, i, j int) int {
    if i == 0 { return j } // delete all of word2
    if j == 0 { return i } // delete all of word1
    if v, ok := memo[[2]int{i, j}]; ok { return v }

    result := 0
    if word1[i-1] == word2[j-1] {
        result = editDistanceMemo(word1, word2, i-1, j-1) // no operation needed
    } else {
        result = 1 + min(
            editDistanceMemo(word1, word2, i-1, j),   // delete from word1
            min(editDistanceMemo(word1, word2, i, j-1), // insert into word1
                editDistanceMemo(word1, word2, i-1, j-1)), // replace
        )
    }
    memo[[2]int{i, j}] = result
    return result
}

// BOTTOM-UP (tabulation): explicit table, fill in order
func editDistance(word1, word2 string) int {
    m, n := len(word1), len(word2)
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }

    for i := 0; i <= m; i++ { dp[i][0] = i }
    for j := 0; j <= n; j++ { dp[0][j] = j }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if word1[i-1] == word2[j-1] {
                dp[i][j] = dp[i-1][j-1]
            } else {
                dp[i][j] = 1 + min(dp[i-1][j], min(dp[i][j-1], dp[i-1][j-1]))
            }
        }
    }
    return dp[m][n]
}
// Time: O(m*n), Space: O(m*n) — can reduce to O(n) with row optimization
```

---

## 8. Classic Hard Problems

### Longest Common Subsequence

```go
func longestCommonSubsequence(text1, text2 string) int {
    m, n := len(text1), len(text2)
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if text1[i-1] == text2[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
            } else {
                dp[i][j] = max(dp[i-1][j], dp[i][j-1])
            }
        }
    }
    return dp[m][n]
}
```

### Regular Expression Matching

```go
// dp[i][j] = true if pattern p[0:j] matches string s[0:i]
func isMatch(s, p string) bool {
    m, n := len(s), len(p)
    dp := make([][]bool, m+1)
    for i := range dp { dp[i] = make([]bool, n+1) }
    dp[0][0] = true

    // Handle patterns like a*, a*b*, a*b*c* that match empty string
    for j := 2; j <= n; j++ {
        if p[j-1] == '*' { dp[0][j] = dp[0][j-2] }
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if p[j-1] == '*' {
                dp[i][j] = dp[i][j-2] // zero occurrences of preceding char
                if p[j-2] == '.' || p[j-2] == s[i-1] {
                    dp[i][j] = dp[i][j] || dp[i-1][j] // one+ occurrences
                }
            } else if p[j-1] == '.' || p[j-1] == s[i-1] {
                dp[i][j] = dp[i-1][j-1]
            }
        }
    }
    return dp[m][n]
}
```

---

## Summary

- DP = optimal substructure + overlapping subproblems.
- **5 patterns:** Linear (1D), Grid (2D), Interval (subarray), Knapsack (choices + constraint), State Machine (multiple states).
- **Process:** define dp[i], find recurrence, set base cases, fill in order, extract answer.
- **Memoization** (top-down): natural recursion + cache, only computes what is needed.
- **Tabulation** (bottom-up): explicit table, often more space-efficient.
- LIS key insight: O(n log n) with patience sorting. Coin change is unbounded knapsack. Edit distance is 2D DP.

---

## Exercises

### Easy
1. Fibonacci number using DP. Space-optimize to O(1).
2. Count number of ways to climb n stairs when you can take 1, 2, or 3 steps.
3. Find the minimum cost to reach the last index of an array (you can jump 1 or 2 steps, cost = value at current step).

### Medium
4. Decode ways: "12" can be decoded as "AB" (1,2) or "L" (12). Count total decodings.
5. Maximum product subarray (not sum — can be tricky with negatives).
6. Target sum: assign + or - to each number to reach target. Count the number of ways.

### Hard
7. Implement "Scramble String" — dp with 3D state or memoization.
8. Given a string, find the minimum number of cuts for palindrome partitioning.
9. Stone game: two players take from either end; return true if first player always wins with optimal play.
