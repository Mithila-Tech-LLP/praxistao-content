---
title: Factorial
number: 5
difficulty: easy
duration: 10 minutes
concept: Recursion
---

## What You Need to Build

Write a function called `Factorial` that computes the factorial of a non-negative integer.

The factorial of `n` (written `n!`) is the product of all positive integers from 1 to `n`.

```
5! = 5 × 4 × 3 × 2 × 1 = 120
0! = 1  (by definition)
1! = 1
```

## Function Signature

```go
func Factorial(n int) int
```

## Examples

| Input | Output |
|-------|--------|
| `0` | `1` |
| `1` | `1` |
| `5` | `120` |
| `10` | `3628800` |

## Key Concept: Recursion

A recursive function calls **itself** to solve a smaller version of the same problem.

Factorial has a perfect recursive structure:
- `factorial(0) = 1` ← **base case** (stops the recursion)
- `factorial(n) = n × factorial(n - 1)` ← **recursive case**

Tracing `factorial(4)`:
```
factorial(4)
  = 4 × factorial(3)
  = 4 × 3 × factorial(2)
  = 4 × 3 × 2 × factorial(1)
  = 4 × 3 × 2 × 1 × factorial(0)
  = 4 × 3 × 2 × 1 × 1
  = 24
```

Every recursive solution needs:
1. A **base case** that returns a fixed value without recursing
2. A **recursive case** that gets closer to the base case with each call

## Requirements

- `Factorial(0)` → `1`
- `Factorial(1)` → `1`
- `Factorial(5)` → `120`
- You may use either recursion or an iterative loop

## Hints

<details>
<summary>Hint 1 — Recursive structure in Go</summary>

```go
func Factorial(n int) int {
    if n == 0 {
        return 1  // base case
    }
    return n * Factorial(n-1)  // recursive case
}
```

</details>

<details>
<summary>Hint 2 — Iterative alternative</summary>

Use a loop: `result := 1; for i := 2; i <= n; i++ { result *= i }; return result`.

</details>

## How to Verify

```bash
lncli run
```
