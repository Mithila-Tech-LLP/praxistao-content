# Chapter 47: Recursion and Backtracking

Recursion solves problems by breaking them into smaller instances of the same problem. Backtracking is recursion with the ability to undo choices — it builds a solution incrementally and abandons partial solutions that can't possibly work.

## Table of Contents

1. [Recursion Fundamentals](#1-recursion-fundamentals)
2. [Memoization vs Tabulation](#2-memoization-vs-tabulation)
3. [Backtracking Framework](#3-backtracking-framework)
4. [Classic Problems](#4-classic-problems)
5. [Pruning Strategies](#5-pruning-strategies)
6. [Summary](#summary)
7. [Exercises](#exercises)

---

## 1. Recursion Fundamentals

Every recursive function has two parts:
1. **Base case**: the simplest input, handled directly
2. **Recursive case**: break the problem into a smaller sub-problem, solve it, combine

```go
// Classic: factorial
// n! = n * (n-1)! with base case 0! = 1
func factorial(n int) int {
    if n <= 1 { return 1 }           // base case
    return n * factorial(n-1)         // recursive case
}

// Call tree for factorial(4):
// factorial(4) = 4 * factorial(3)
//                    = 3 * factorial(2)
//                         = 2 * factorial(1)
//                              = 1
```

### Thinking recursively: trust the recursion

The hardest part of recursion is trusting that the recursive call works correctly. To write `mergeSort(arr)`:

1. **Define what the function does**: returns a sorted copy of arr
2. **Trust the recursive call**: `mergeSort(left)` returns a sorted left half
3. **Combine**: merge two sorted halves

```go
func mergeSort(arr []int) []int {
    if len(arr) <= 1 { return arr }  // base case
    mid := len(arr) / 2
    left := mergeSort(arr[:mid])      // trust: returns sorted
    right := mergeSort(arr[mid:])     // trust: returns sorted
    return merge(left, right)         // combine sorted halves
}
```

### Stack depth and tail recursion

Go doesn't optimize tail recursion, so deep recursion can overflow the stack. Convert to iteration when depth is proportional to n.

```go
// Recursive: O(n) stack space
func sum(nums []int, i int) int {
    if i == len(nums) { return 0 }
    return nums[i] + sum(nums, i+1)
}

// Iterative: O(1) stack space
func sumIter(nums []int) int {
    total := 0
    for _, n := range nums { total += n }
    return total
}
```

---

## 2. Memoization vs Tabulation

When a recursion recomputes the same sub-problems, cache the results.

```go
// Fibonacci without memoization: O(2^n) time — catastrophically slow
func fib(n int) int {
    if n <= 1 { return n }
    return fib(n-1) + fib(n-2)
}

// With memoization: O(n) time, O(n) space
func fibMemo(n int, memo map[int]int) int {
    if n <= 1 { return n }
    if v, ok := memo[n]; ok { return v }
    result := fibMemo(n-1, memo) + fibMemo(n-2, memo)
    memo[n] = result
    return result
}
```

The difference: without memo, `fib(40)` makes ~2 billion calls. With memo, it makes 40.

---

## 3. Backtracking Framework

Backtracking explores a decision tree. At each node, you:
1. Make a choice
2. Recurse into the next decision
3. Undo the choice (backtrack)

```go
// Generic backtracking template:
func backtrack(state State, choices []Choice, results *[]Solution) {
    if isSolution(state) {
        *results = append(*results, copySolution(state))
        return
    }
    for _, choice := range choices {
        if isValid(state, choice) {
            applyChoice(state, choice)          // make choice
            backtrack(state, nextChoices, results) // recurse
            undoChoice(state, choice)            // undo choice
        }
    }
}
```

**Key insight**: backtracking is a DFS over the decision tree, where branches that can't lead to a solution are pruned early.

---

## 4. Classic Problems

### Subsets (all 2^n subsets)

```go
func subsets(nums []int) [][]int {
    var result [][]int
    var backtrack func(start int, current []int)
    backtrack = func(start int, current []int) {
        // Every state is a valid subset — add a copy
        tmp := make([]int, len(current))
        copy(tmp, current)
        result = append(result, tmp)

        for i := start; i < len(nums); i++ {
            current = append(current, nums[i])      // choose
            backtrack(i+1, current)                  // explore
            current = current[:len(current)-1]       // unchoose
        }
    }
    backtrack(0, []int{})
    return result
}
// subsets([1,2,3]) = [[], [1], [1,2], [1,2,3], [1,3], [2], [2,3], [3]]
```

### Permutations

```go
func permutations(nums []int) [][]int {
    var result [][]int
    used := make([]bool, len(nums))
    
    var backtrack func(current []int)
    backtrack = func(current []int) {
        if len(current) == len(nums) {
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for i, n := range nums {
            if used[i] { continue }
            used[i] = true
            current = append(current, n)
            backtrack(current)
            current = current[:len(current)-1]
            used[i] = false
        }
    }
    backtrack([]int{})
    return result
}
```

### Combinations (choose k from n)

```go
func combinations(n, k int) [][]int {
    var result [][]int
    var backtrack func(start int, current []int)
    backtrack = func(start int, current []int) {
        if len(current) == k {
            tmp := make([]int, k)
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        // Pruning: remaining = n - start + 1, need = k - len(current)
        // If remaining < need, can't complete — prune
        for i := start; i <= n-(k-len(current))+1; i++ {
            current = append(current, i)
            backtrack(i+1, current)
            current = current[:len(current)-1]
        }
    }
    backtrack(1, []int{})
    return result
}
```

### N-Queens

Place N queens on an N×N board so no two queens attack each other.

```go
func solveNQueens(n int) [][]string {
    var result [][]string
    board := make([]int, n) // board[row] = col of queen in that row

    var backtrack func(row int)
    backtrack = func(row int) {
        if row == n {
            result = append(result, buildBoard(board, n))
            return
        }
        for col := 0; col < n; col++ {
            if isValid(board, row, col) {
                board[row] = col
                backtrack(row + 1)
            }
        }
    }
    backtrack(0)
    return result
}

func isValid(board []int, row, col int) bool {
    for r := 0; r < row; r++ {
        c := board[r]
        if c == col || abs(c-col) == abs(r-row) {
            return false // same column or diagonal
        }
    }
    return true
}

func buildBoard(board []int, n int) []string {
    rows := make([]string, n)
    for r, c := range board {
        row := make([]byte, n)
        for i := range row { row[i] = '.' }
        row[c] = 'Q'
        rows[r] = string(row)
    }
    return rows
}

func abs(x int) int { if x < 0 { return -x }; return x }
```

### Sudoku Solver

```go
func solveSudoku(board [][]byte) {
    solve(board)
}

func solve(board [][]byte) bool {
    for r := 0; r < 9; r++ {
        for c := 0; c < 9; c++ {
            if board[r][c] != '.' { continue }
            for digit := byte('1'); digit <= '9'; digit++ {
                if canPlace(board, r, c, digit) {
                    board[r][c] = digit
                    if solve(board) { return true }
                    board[r][c] = '.'  // backtrack
                }
            }
            return false // no digit works → dead end
        }
    }
    return true // all cells filled
}

func canPlace(board [][]byte, row, col int, digit byte) bool {
    boxRow, boxCol := (row/3)*3, (col/3)*3
    for i := 0; i < 9; i++ {
        if board[row][i] == digit { return false }  // same row
        if board[i][col] == digit { return false }  // same col
        if board[boxRow+i/3][boxCol+i%3] == digit { return false }  // same box
    }
    return true
}
```

### Word Search

Find if a word exists in a grid by moving in 4 directions.

```go
func exist(board [][]byte, word string) bool {
    rows, cols := len(board), len(board[0])
    
    var dfs func(r, c, idx int) bool
    dfs = func(r, c, idx int) bool {
        if idx == len(word) { return true }
        if r < 0 || r >= rows || c < 0 || c >= cols { return false }
        if board[r][c] != word[idx] { return false }

        temp := board[r][c]
        board[r][c] = '#'  // mark visited
        found := dfs(r+1, c, idx+1) || dfs(r-1, c, idx+1) ||
                 dfs(r, c+1, idx+1) || dfs(r, c-1, idx+1)
        board[r][c] = temp  // restore
        return found
    }

    for r := range board {
        for c := range board[r] {
            if dfs(r, c, 0) { return true }
        }
    }
    return false
}
```

---

## 5. Pruning Strategies

Pruning makes backtracking practical by cutting branches early.

```go
// Sum to target: find all combinations summing to target
// Without pruning: tries all 2^n subsets
// With pruning: sort first, skip once sum exceeds target
func combinationSum(candidates []int, target int) [][]int {
    sort.Ints(candidates)  // sort enables early termination
    var result [][]int

    var backtrack func(start, remaining int, current []int)
    backtrack = func(start, remaining int, current []int) {
        if remaining == 0 {
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remaining { break }  // PRUNE: sorted, rest are all larger
            current = append(current, candidates[i])
            backtrack(i, remaining-candidates[i], current) // i, not i+1 (can reuse)
            current = current[:len(current)-1]
        }
    }
    backtrack(0, target, []int{})
    return result
}

// Skip duplicates (avoid identical results)
func combinationSumUnique(candidates []int, target int) [][]int {
    sort.Ints(candidates)
    var result [][]int

    var backtrack func(start, remaining int, current []int)
    backtrack = func(start, remaining int, current []int) {
        if remaining == 0 {
            tmp := make([]int, len(current))
            copy(tmp, current)
            result = append(result, tmp)
            return
        }
        for i := start; i < len(candidates); i++ {
            if candidates[i] > remaining { break }
            if i > start && candidates[i] == candidates[i-1] { continue } // SKIP DUPLICATE
            current = append(current, candidates[i])
            backtrack(i+1, remaining-candidates[i], current)
            current = current[:len(current)-1]
        }
    }
    backtrack(0, target, []int{})
    return result
}
```

---

## Summary

| Problem | Pattern |
|---------|---------|
| All subsets | Recurse with `start` index, add copy at each node |
| All permutations | `used[]` array, iterate all positions |
| Combinations | Recurse with `start` + size check |
| Grid search | Mark visited, restore on return |
| Constraint satisfaction (N-Queens, Sudoku) | Check constraint before placing |
| Sum/count problems | Sort + early break when exceeded |

**Backtracking cost**: exponential in the worst case. Pruning is essential to make it practical. Time complexity is typically O(k × 2^n) or O(k × n!) where k is the work done per leaf.

---

## Exercises

### Easy
1. Generate all valid parentheses combinations for `n` pairs (e.g., n=3 → `((()))`, `(()())`, `(())()`, `()(())`, `()()()`). Use backtracking with pruning: only add `(` if open count < n, only add `)` if close count < open count.
2. Implement `letterCombinations(digits string)` that returns all letter combinations a phone number could represent (2="abc", 3="def", etc.).
3. Find all paths from top-left to bottom-right of a grid that sum to a target value. You can only move right or down.

### Medium
4. Solve the **Hamiltonian Path** problem: given a graph, find a path that visits every vertex exactly once. Use backtracking with a `visited` set.
5. Implement **expression generator**: given digits `[1, 2, 3]` and target `6`, find all ways to insert `+`, `-`, `*` between digits to reach the target. Handle multi-digit numbers and operator precedence.
6. **Palindrome partitioning**: given a string, find all ways to partition it such that every substring is a palindrome. Prune by checking palindrome before recursing.

### Hard
7. Implement a **SAT solver** using backtracking with DPLL (Davis–Putnam–Logemann–Loveland): given a boolean formula in CNF (conjunctive normal form), find an assignment of variables that satisfies all clauses. Implement unit propagation and pure literal elimination as pruning.
8. Implement a **constraint propagation + backtracking** Sudoku solver: before backtracking, propagate constraints (if a cell has only one possibility, fill it immediately). Measure how many cells can be solved by constraint propagation alone vs requiring backtracking.
