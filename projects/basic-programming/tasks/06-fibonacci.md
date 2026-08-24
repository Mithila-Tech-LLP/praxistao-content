---
title: Fibonacci
number: 6
difficulty: easy
duration: 10-15 minutes
concept: Recursion, Iteration
---

## What You Need to Build

Write a function called `Fibonacci` that returns the nth Fibonacci number.

The Fibonacci sequence: `0, 1, 1, 2, 3, 5, 8, 13, 21, 34, ...`

Each number is the sum of the two before it.

## Function Signature

```go
func Fibonacci(n int) int
```

## Examples

| Input (n) | Output | Sequence position |
|-----------|--------|-------------------|
| `0` | `0` | First: 0 |
| `1` | `1` | Second: 1 |
| `2` | `1` | Third: 1 |
| `5` | `5` | Sixth: 5 |
| `10` | `55` | Eleventh: 55 |

## Key Concept: Recursion vs Iteration

**Recursive approach** — elegant but slow for large n:
```
fib(5) = fib(4) + fib(3)
       = (fib(3) + fib(2)) + (fib(2) + fib(1))
       = ...  (many repeated calculations)
```

The recursive definition:
- `fib(0) = 0`
- `fib(1) = 1`
- `fib(n) = fib(n-1) + fib(n-2)`

**Iterative approach** — efficient:
Keep track of the previous two values and compute forward:
```
n=0: a=0
n=1: a=0, b=1
n=2: a=1, b=1  (b = a+b, a = old b)
n=3: a=1, b=2
n=4: a=2, b=3
n=5: a=3, b=5  ← answer
```

## Requirements

- `Fibonacci(0)` → `0`
- `Fibonacci(1)` → `1`
- `Fibonacci(10)` → `55`
- You may use either recursion or iteration

## Hints

<details>
<summary>Hint 1 — Recursive solution</summary>

```go
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    return Fibonacci(n-1) + Fibonacci(n-2)
}
```

</details>

<details>
<summary>Hint 2 — Iterative solution</summary>

```go
func Fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    a, b := 0, 1
    for i := 2; i <= n; i++ {
        a, b = b, a+b
    }
    return b
}
```

</details>

## How to Verify

```bash
lncli run
```
