# Chapter 51: Dynamic Programming — Advanced Patterns (2D DP, Memoization)

Chapter 50 gave you the DP toolkit: state, recurrence, base cases, and the first look at 2D tables with LCS, edit distance, and 0/1 knapsack. This chapter turns those single examples into *pattern families*. You'll learn to recognize when a problem is "grid DP in disguise" or "a knapsack variant", how to convert between top-down memoization and bottom-up tabulation mechanically, and how to cut a 2D table down to one row. The goal: given a new problem, you should be able to name the family it belongs to before writing a line of code.

## Table of Contents

1. [Recognizing Multi-Dimensional State](#1-recognizing-multi-dimensional-state)
2. [Grid DP Patterns](#2-grid-dp-patterns)
3. [Two-Sequence DP — The LCS Family](#3-two-sequence-dp--the-lcs-family)
4. [Knapsack Variants](#4-knapsack-variants)
5. [Top-Down vs Bottom-Up in Practice](#5-top-down-vs-bottom-up-in-practice)
6. [Space Optimization Techniques](#6-space-optimization-techniques)
7. [Choosing a Pattern](#7-choosing-a-pattern)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Recognizing Multi-Dimensional State

In 1D DP, one index describes a subproblem: `dp[i]` = "answer for the first `i` elements". A problem needs a second dimension when **one index can't uniquely identify a subproblem**. Ask: "what do I need to know to answer the smaller question?" Every independent piece of information becomes a dimension.

```
Problem                          State needed                       Dimensions
-------------------------------- ---------------------------------- ----------
Max sum of prefix                position i                          dp[i]
Path through a grid              row i AND column j                  dp[i][j]
Compare two strings              position in s AND position in t     dp[i][j]
Knapsack                         item index AND remaining capacity   dp[i][w]
Buy/sell stock with cooldown     day AND holding-or-not              dp[i][2]
Interval problems                left endpoint AND right endpoint    dp[i][j]
```

The recurrence then answers: "given the choices available at state `(i, j)`, which smaller states do I combine?" Everything else — memoization, tabulation, space optimization — is mechanical once the state is right.

**Note on `min`/`max`**: this chapter uses Go's built-in `min` and `max` functions (available since Go 1.21) — no helper functions needed.

---

## 2. Grid DP Patterns

Chapter 50 showed `UniquePaths` — count paths moving only right/down. The grid family is much bigger. The shared skeleton: `dp[i][j]` = best answer *at cell (i, j)*, computed from the cells you could have arrived from.

### Unique Paths with Obstacles

Same as `UniquePaths`, but some cells are blocked. The recurrence survives; only base cases need care:

```go
// UniquePathsWithObstacles: grid[i][j] == 1 means blocked.
func UniquePathsWithObstacles(grid [][]int) int {
    m, n := len(grid), len(grid[0])
    dp := make([][]int, m)
    for i := range dp { dp[i] = make([]int, n) }

    for i := 0; i < m; i++ {
        for j := 0; j < n; j++ {
            if grid[i][j] == 1 { continue }  // Blocked cell contributes 0 paths
            if i == 0 && j == 0 {
                dp[i][j] = 1
                continue
            }
            if i > 0 { dp[i][j] += dp[i-1][j] }
            if j > 0 { dp[i][j] += dp[i][j-1] }
        }
    }
    return dp[m-1][n-1]
}
// Lesson: don't hardcode "first row/column = 1" like plain UniquePaths.
// An obstacle in row 0 blocks every cell after it.
```

### Minimum Falling Path Sum

Move from any cell in the top row to the bottom row; each step goes down, down-left, or down-right. Now a cell has up to **three** predecessors:

```go
// MinFallingPathSum: minimize the sum of a top-to-bottom path.
// dp[i][j] = matrix[i][j] + min(dp[i-1][j-1], dp[i-1][j], dp[i-1][j+1])
func MinFallingPathSum(matrix [][]int) int {
    n := len(matrix)
    prev := append([]int(nil), matrix[0]...)  // Row 0 is its own base case

    for i := 1; i < n; i++ {
        curr := make([]int, n)
        for j := 0; j < n; j++ {
            best := prev[j]
            if j > 0 && prev[j-1] < best { best = prev[j-1] }
            if j < n-1 && prev[j+1] < best { best = prev[j+1] }
            curr[j] = matrix[i][j] + best
        }
        prev = curr
    }

    ans := prev[0]
    for _, v := range prev[1:] { ans = min(ans, v) }
    return ans
}
// Note: the answer is the min over the WHOLE last row — the path can end anywhere.
```

### Triangle

The grid doesn't have to be rectangular. For a triangle (row `i` has `i+1` cells), working **bottom-up** avoids edge cases entirely:

```go
// MinTriangleTotal: minimum path sum from top of triangle to bottom.
// Process from the bottom row upward: dp[j] = triangle[i][j] + min(dp[j], dp[j+1])
func MinTriangleTotal(triangle [][]int) int {
    dp := append([]int(nil), triangle[len(triangle)-1]...)

    for i := len(triangle) - 2; i >= 0; i-- {
        for j := 0; j <= i; j++ {
            dp[j] = triangle[i][j] + min(dp[j], dp[j+1])
        }
    }
    return dp[0]
}
// Going bottom-up: every cell has exactly 2 children — no boundary checks needed.
```

### Maximal Square

`dp[i][j]` doesn't have to mean "path". Here it means "side length of the largest all-ones square whose **bottom-right corner** is (i, j)":

```go
// MaximalSquare: area of the largest square of '1's in a binary matrix.
// dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])  if cell is '1'
func MaximalSquare(matrix [][]byte) int {
    if len(matrix) == 0 { return 0 }
    m, n := len(matrix), len(matrix[0])
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }  // Padded — row/col 0 stay zero

    best := 0
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if matrix[i-1][j-1] == '1' {
                dp[i][j] = 1 + min(dp[i-1][j], dp[i][j-1], dp[i-1][j-1])
                best = max(best, dp[i][j])
            }
        }
    }
    return best * best
}
// Why min of three neighbors? A k×k square at (i,j) requires (k-1)×(k-1) squares
// ending above, to the left, and diagonally — the weakest one limits the size.
```

### Dungeon Game — When Forward DP Fails

Sometimes the natural left-to-right direction is wrong. A knight walks from top-left to bottom-right; cells add or subtract health; health must stay ≥ 1 at all times. Find the minimum starting health.

Forward DP fails because two quantities matter at once (health gained so far *and* the worst dip along the way), and neither alone is a valid state. The fix: **compute backwards** — `dp[i][j]` = minimum health needed *when entering* `(i, j)`:

```go
// CalculateMinimumHP: minimum initial health to survive the dungeon.
// dp[i][j] = max(1, min(dp[i+1][j], dp[i][j+1]) - dungeon[i][j])
func CalculateMinimumHP(dungeon [][]int) int {
    m, n := len(dungeon), len(dungeon[0])
    const INF = 1 << 60

    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
        for j := range dp[i] { dp[i][j] = INF }  // Padding = "unreachable"
    }
    dp[m][n-1], dp[m-1][n] = 1, 1  // Exiting the dungeon requires health 1

    for i := m - 1; i >= 0; i-- {
        for j := n - 1; j >= 0; j-- {
            need := min(dp[i+1][j], dp[i][j+1]) - dungeon[i][j]
            dp[i][j] = max(1, need)  // Health can never drop below 1
        }
    }
    return dp[0][0]
}
// Rule of thumb: if the state depends on "what happens later", reverse the direction.
```

---

## 3. Two-Sequence DP — The LCS Family

Chapter 50 built the LCS and edit distance tables. Those two recurrences are the parents of a whole family — most "compare two strings" problems are small mutations of them.

### Longest Common Substring (vs Subsequence)

One word changes ("substring" = contiguous) and the recurrence changes in one place: a mismatch **resets to zero** instead of carrying the best-so-far.

```go
// LongestCommonSubstring: length of the longest contiguous common run.
// dp[i][j] = dp[i-1][j-1] + 1 if s[i-1] == t[j-1], else 0
func LongestCommonSubstring(s, t string) int {
    m, n := len(s), len(t)
    dp := make([][]int, m+1)
    for i := range dp { dp[i] = make([]int, n+1) }

    best := 0
    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if s[i-1] == t[j-1] {
                dp[i][j] = dp[i-1][j-1] + 1
                best = max(best, dp[i][j])
            }
            // Mismatch: dp[i][j] stays 0 — a substring can't skip characters.
        }
    }
    return best
}
// Also note: the answer is the max over the whole table, not dp[m][n] —
// the best run can end anywhere.
```

### Shortest Common Supersequence

The shortest string that contains both `s` and `t` as subsequences. No new table needed — share the overlap once:

```go
// SCSLength: every character of s and t appears once, and the LCS is shared.
func SCSLength(s, t string) int {
    return len(s) + len(t) - LCS(s, t)  // LCS from Chapter 50
}
```

To build the actual string, backtrack through the LCS table exactly like `LCSString` in Chapter 50, but emit *every* character you walk past instead of only the matches.

### Distinct Subsequences

Count how many times `t` appears in `s` as a subsequence. The state is the same `(i, j)`; the recurrence becomes a **sum of ways** instead of a max of lengths:

```go
// NumDistinct: number of distinct subsequences of s equal to t.
// dp[i][j] = ways to build t[:j] from s[:i]
// Always allowed: skip s[i-1]           → dp[i-1][j]
// If chars match: also use s[i-1]       → dp[i-1][j-1]
func NumDistinct(s, t string) int {
    m, n := len(s), len(t)
    dp := make([][]int, m+1)
    for i := range dp {
        dp[i] = make([]int, n+1)
        dp[i][0] = 1  // The empty t can be built exactly one way: take nothing
    }

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            dp[i][j] = dp[i-1][j]
            if s[i-1] == t[j-1] {
                dp[i][j] += dp[i-1][j-1]
            }
        }
    }
    return dp[m][n]
}
// Verify: NumDistinct("rabbbit", "rabbit") == 3
```

**The family resemblance**: state is always `(prefix of s, prefix of t)`. What changes is the *combiner* — `max` for longest, `min` for cheapest (edit distance), `+` for counting. When you meet a new two-string problem, start from this template and only redesign the combiner.

---

## 4. Knapsack Variants

Chapter 50 covered 0/1 knapsack and its backward-loop space optimization. Nearly every "pick items subject to a budget" problem is one of four variants — and they differ mainly in **loop direction and what dp stores**.

### Subset Sum

"Can some subset hit exactly `target`?" — 0/1 knapsack where value doesn't matter and dp stores a boolean:

```go
// SubsetSum: can any subset of nums sum to target?
func SubsetSum(nums []int, target int) bool {
    dp := make([]bool, target+1)
    dp[0] = true  // The empty subset sums to 0

    for _, num := range nums {
        for s := target; s >= num; s-- {  // Backwards: each num used at most once
            if dp[s-num] { dp[s] = true }
        }
    }
    return dp[target]
}
```

Partition Equal Subset Sum (Chapter 50, exercise 6) reduces directly to this: `SubsetSum(nums, total/2)`.

### Target Sum — Transform, Then Knapsack

Assign `+` or `-` to every number so the expression equals `target`. Looks nothing like knapsack — until you do algebra. Let `P` = sum of positives, `N` = sum of negatives. Then `P - N = target` and `P + N = total`, so `P = (total + target) / 2`. The problem becomes: *count* subsets summing to `P`:

```go
// TargetSumWays: number of ways to assign +/- to nums so the sum equals target.
func TargetSumWays(nums []int, target int) int {
    total := 0
    for _, n := range nums { total += n }
    if target > total || target < -total || (total+target)%2 != 0 {
        return 0  // P must be a non-negative integer
    }
    p := (total + target) / 2

    dp := make([]int, p+1)
    dp[0] = 1
    for _, num := range nums {
        for s := p; s >= num; s-- {  // Backwards again: 0/1 counting
            dp[s] += dp[s-num]
        }
    }
    return dp[p]
}
// Verify: TargetSumWays([]int{1,1,1,1,1}, 3) == 5
```

### Unbounded Knapsack — Rod Cutting

Each item may be used **any number of times**. The only code change from 0/1: iterate capacity **forwards**, so a state can reuse an item already counted in this round:

```go
// RodCutting: prices[i] = price of a piece of length i+1.
// Maximize revenue from cutting a rod of the given length.
func RodCutting(prices []int, length int) int {
    dp := make([]int, length+1)
    for l := 1; l <= length; l++ {
        for cut := 1; cut <= l && cut <= len(prices); cut++ {
            dp[l] = max(dp[l], prices[cut-1]+dp[l-cut])
        }
    }
    return dp[length]
}
// Verify: RodCutting([]int{1,5,8,9,10,17,17,20}, 8) == 22 (cut into 2 + 6)
```

`CoinChange` from Chapter 50 is unbounded knapsack too — same forward loop.

### The Loop-Direction Rule

This one table resolves most knapsack confusion:

| Variant | Inner loop over capacity | Why |
|---|---|---|
| 0/1 (each item once) | **Backwards** | `dp[s-num]` must be from the *previous* item row |
| Unbounded (unlimited reuse) | **Forwards** | `dp[s-num]` may already include the current item |
| Count combinations | Items **outer**, capacity inner | Fixes item order → `{2,3}` counted once |
| Count permutations | Capacity **outer**, items inner | Every ordering counted → `{2,3}` and `{3,2}` |

---

## 5. Top-Down vs Bottom-Up in Practice

Chapter 50 compared memoization and tabulation in a table. Here's the practical skill: **converting between them mechanically**. The mapping is exact:

```
Memoized recursion                    Tabulation
------------------------------------- -------------------------------------
memo key (i, j)                    →  table index dp[i][j]
recursion base case                →  table initialization
recursive calls solve(i-1, ...)    →  cells read: fill order must compute
                                       them first (reverse the call direction)
top-level call solve(n, m)         →  the cell you return at the end
```

### Worked Example: Both Directions

Top-down minimum falling path (compare with the bottom-up version in Section 2):

```go
import "math"

// MinFallingPathSumMemo: same problem, top-down.
func MinFallingPathSumMemo(matrix [][]int) int {
    n := len(matrix)
    memo := make([][]int, n)
    for i := range memo {
        memo[i] = make([]int, n)
        for j := range memo[i] { memo[i][j] = math.MinInt }  // Sentinel: "not computed"
    }

    var solve func(i, j int) int  // Min path sum from (i, j) down to the bottom row
    solve = func(i, j int) int {
        if j < 0 || j >= n { return math.MaxInt / 2 }  // Off the grid
        if i == n-1 { return matrix[i][j] }            // Base: bottom row
        if memo[i][j] != math.MinInt { return memo[i][j] }

        memo[i][j] = matrix[i][j] + min(solve(i+1, j-1), solve(i+1, j), solve(i+1, j+1))
        return memo[i][j]
    }

    ans := math.MaxInt
    for j := 0; j < n; j++ {
        ans = min(ans, solve(0, j))
    }
    return ans
}
```

Notice how the pieces line up: the recursion's base case (bottom row) became tabulation's initialization (`prev = matrix[0]` — or last row if you tabulate downward); the recursion goes *down* so tabulation fills *up* (or vice versa); the final loop over starting columns became the final loop over the last row.

### When to Prefer Which

- **Memoization wins** when the reachable state space is sparse. Example: `(i, remaining)` states where `remaining` takes only a few of its possible values — tabulation would fill millions of cells that are never asked about. Also wins when the fill order is hard to see (interval DP recurrences are easier to get right recursively).
- **Tabulation wins** when you'll touch most states anyway: it has no function-call overhead, no risk of stack overflow (Go goroutine stacks grow, but a 10⁶-deep recursion is still slow), and — critically — it **enables space optimization** (Section 6), which memoization can't do because old states must stay cached.
- **In Go specifically**: the closure-based memo pattern (`var solve func(...)`) is idiomatic and keeps the cache scoped to one call. Use a 2D slice with a sentinel, not a `map[[2]int]int` — the map is 5–10× slower.

---

## 6. Space Optimization Techniques

Every optimization here follows from one observation: **if row `i` only reads row `i-1`, you don't need rows `0..i-2`.**

### Two Rolling Rows (LCS)

```go
// LCSTwoRows: LCS length in O(min side) extra space.
func LCSTwoRows(s, t string) int {
    if len(t) > len(s) { s, t = t, s }  // Make t the shorter string
    m, n := len(s), len(t)

    prev := make([]int, n+1)
    curr := make([]int, n+1)

    for i := 1; i <= m; i++ {
        for j := 1; j <= n; j++ {
            if s[i-1] == t[j-1] {
                curr[j] = prev[j-1] + 1
            } else {
                curr[j] = max(prev[j], curr[j-1])
            }
        }
        prev, curr = curr, prev  // Swap slices — no allocation
    }
    return prev[n]  // After the final swap, prev holds the last computed row
}
```

### One Row Plus a Diagonal (Edit Distance)

LCS needed `prev[j-1]` (the diagonal), which the two-row version keeps naturally. You can go further — one row — by saving the diagonal in a variable *before* overwriting it:

```go
// EditDistanceOneRow: Levenshtein distance in O(n) space.
func EditDistanceOneRow(s, t string) int {
    m, n := len(s), len(t)
    dp := make([]int, n+1)
    for j := range dp { dp[j] = j }  // Row 0: insert j characters

    for i := 1; i <= m; i++ {
        prevDiag := dp[0]  // This is dp[i-1][0]
        dp[0] = i          // Column 0: delete i characters
        for j := 1; j <= n; j++ {
            temp := dp[j]  // Save dp[i-1][j] before overwriting
            if s[i-1] == t[j-1] {
                dp[j] = prevDiag
            } else {
                dp[j] = 1 + min(dp[j], dp[j-1], prevDiag)
                //              ↑ delete  ↑ insert  ↑ replace
            }
            prevDiag = temp
        }
    }
    return dp[n]
}
```

### In-Place on the Input

When you're allowed to destroy the input, the grid itself is the DP table — `MinFallingPathSum` could write into `matrix` directly for O(1) extra space. Production code usually shouldn't (surprising side effects), but it's worth knowing for memory-constrained cases.

### What You Give Up

Space optimization keeps the *answer* but destroys the *history*: reconstruction (like `LCSString` in Chapter 50) needs the full table to backtrack through. If you need the actual subsequence/path/item set, keep the 2D table — or use Hirschberg's algorithm (divide and conquer over the rows) to reconstruct in linear space, a classic hard exercise.

---

## 7. Choosing a Pattern

A field guide for mapping a fresh problem onto a family:

```
The problem mentions...                    Try...
------------------------------------------ ----------------------------------------
a grid, robot, paths, falling              Grid DP: dp[i][j] from neighbor cells
two strings/arrays being compared          LCS family: dp[i][j] over two prefixes
"can we reach exactly X" / budgets         Subset sum / knapsack (watch loop direction!)
items reusable                             Unbounded knapsack (forward loop)
count the ways vs best value               Same state — change combiner to +
must survive / constraint on the path      Consider REVERSE direction DP
subarray/substring (contiguous)            Reset-on-mismatch recurrences
"minimum starting X"                       Backwards DP from the goal
```

And the boundary with the next chapter: if at every step one choice is *provably* safe without looking ahead — no table needed — you have a greedy problem. Chapter 52 covers how to spot (and prove) that. When greedy's local choice can block a better global answer (0/1 knapsack, coin change with arbitrary denominations), you're back here in DP land.

---

## Summary

- A second DP dimension appears when one index can't identify a subproblem — each independent piece of information (second string position, remaining capacity, row *and* column) is a dimension
- **Grid DP**: `dp[i][j]` from the cells you could arrive from; watch for obstacles breaking base cases, answers spread across a whole row, and non-rectangular shapes (triangle → go bottom-up)
- **Reverse-direction DP** (dungeon game): when validity depends on the future, define the state as "what do I need *from here on*" and fill backwards
- **LCS family**: same `(i, j)` state, different combiner — `max` (LCS), reset-to-0 (common substring), `+` (distinct subsequences), arithmetic (shortest common supersequence = m + n − LCS)
- **Knapsack variants**: 0/1 → backward capacity loop; unbounded → forward; combinations → items outer; permutations → capacity outer
- **Memo ↔ tab conversion** is mechanical: key → index, base case → initialization, call direction → reversed fill order
- Prefer memoization for sparse state spaces and tricky fill orders; tabulation for full tables, speed, and space optimization
- **Space optimization**: two rolling rows, or one row + saved diagonal — but it destroys the history needed for reconstruction

---

## Exercises

### Easy
1. Implement `MinPathSum(grid [][]int) int` (top-left to bottom-right, moves right/down) **twice**: once top-down memoized, once bottom-up with a single rolling row. Verify both give the same answer on `[[1,3,1],[1,5,1],[4,2,1]]` → 7.
2. Implement `IsSubsequence(s, t string) bool` two ways: with the LCS table (`LCS(s,t) == len(s)`), and with a two-pointer scan (O(n), no table). When is the DP version actually needed? (Hint: many queries `s` against one long `t`.)
3. Convert Chapter 50's `LPS` (longest palindromic subsequence) from tabulation to top-down memoization. Which felt easier to write, and why?

### Medium
4. **Interleaving String**: given `s1`, `s2`, `s3`, determine if `s3` is formed by interleaving `s1` and `s2` (preserving each string's internal order). State: `dp[i][j]` = can `s1[:i]` and `s2[:j]` interleave into `s3[:i+j]`. Space-optimize to one row. Verify: `("aabcc", "dbbca", "aadbbcbcac")` → true.
5. **Last Stone Weight II**: smash stones together; each smash of stones `x ≤ y` leaves `y - x`. Minimize the final stone. Show this reduces to partitioning the array into two subsets with minimal difference — a subset-sum sweep over all reachable sums ≤ total/2. Verify: `[2,7,4,1,8,1]` → 1.
6. **Coin Change — Permutations**: Chapter 50's `CoinWays` counts combinations. Write `CoinPermutations(coins []int, amount int) int` that counts ordered sequences (so 1+2 and 2+1 differ) by swapping the loops. Verify: `coins=[1,2], amount=3` → 3 (`1+1+1`, `1+2`, `2+1`).

### Hard
7. **Longest Increasing Path in a Matrix**: from each cell you may move to a strictly larger neighbor (4 directions). Find the longest path. This is memoized DFS on a grid — there's no valid tabulation order until you realize the ordering is "by cell value". Implement both the memoized version and a tabulated version that sorts cells by value first. O(mn) states, O(mn) time.
8. **Hirschberg's LCS**: implement LCS *string reconstruction* in O(min(m,n)) space using divide and conquer: compute the LCS lengths of the first half of `s` against `t` (forward) and the second half against `t` (backward), find the split point of `t` that maximizes the sum, and recurse on the two halves. Verify against Chapter 50's `LCSString` on random inputs.
