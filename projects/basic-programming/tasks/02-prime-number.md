---
title: Prime Number
number: 2
difficulty: easy
duration: 10-15 minutes
concept: Loops, Divisibility
---

## What You Need to Build

Write a function called `IsPrime` that takes an integer and returns `true` if it is a prime number, `false` otherwise.

A **prime number** is a number greater than 1 that has no divisors other than 1 and itself.

## Function Signature

```go
func IsPrime(n int) bool
```

## Examples

| Input | Output | Reason |
|-------|--------|--------|
| `2`   | `true`  | Smallest prime |
| `3`   | `true`  | Only divisible by 1 and 3 |
| `4`   | `false` | 4 = 2 × 2 |
| `17`  | `true`  | Prime |
| `1`   | `false` | By definition, 1 is not prime |
| `0`   | `false` | Not prime |
| `-5`  | `false` | Negative numbers are not prime |

## Key Concept: Trial Division

The straightforward approach: try dividing `n` by every number from `2` to `n-1`. If any of those divisions have zero remainder, `n` is not prime.

You can optimize this: you only need to check up to `√n`. If `n` has a divisor larger than `√n`, it must also have one smaller than `√n`.

```
Is 25 prime?
  Check 2: 25 % 2 = 1, not a divisor
  Check 3: 25 % 3 = 1, not a divisor
  Check 4: 25 % 4 = 1, not a divisor
  Check 5: 25 % 5 = 0, divisor found! → NOT prime
  √25 = 5, so we only needed to check up to 5.
```

## Requirements

- Numbers ≤ 1 are not prime
- `IsPrime(2)` → `true` (2 is the only even prime)
- All even numbers > 2 are not prime
- Correctly identifies primes up to at least 1,000

## Hints

<details>
<summary>Hint 1 — Base cases</summary>

Handle `n <= 1` first and return `false`. Then handle `n == 2` and return `true`.

</details>

<details>
<summary>Hint 2 — The loop</summary>

Loop `i` from `2` up to `i*i <= n`. Inside the loop, if `n % i == 0`, return `false`. After the loop, return `true`.

</details>

<details>
<summary>Hint 3 — Square root optimization</summary>

Instead of `i < n`, use `i*i <= n` as your loop condition. This avoids importing `math.Sqrt` and is equally correct.

</details>

## How to Verify

```bash
lncli run
```
