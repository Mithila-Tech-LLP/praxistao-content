---
title: Sum of Digits
number: 9
difficulty: easy
duration: 10 minutes
concept: Number Decomposition, Modulo
---

## What You Need to Build

Write a function called `SumOfDigits` that takes an integer and returns the sum of all its digits.

## Function Signature

```go
func SumOfDigits(n int) int
```

## Examples

| Input | Output | Reason |
|-------|--------|--------|
| `123` | `6` | 1 + 2 + 3 = 6 |
| `456` | `15` | 4 + 5 + 6 = 15 |
| `0` | `0` | Single digit |
| `10` | `1` | 1 + 0 = 1 |
| `-123` | `6` | Treat as positive: 1 + 2 + 3 = 6 |

## Key Concept: Extracting Digits with Modulo

To extract digits from a number one at a time, use two operations repeatedly:
- `n % 10` gives the **last digit**
- `n / 10` removes the last digit (integer division)

```
n = 456
  456 % 10 = 6 → digit, sum = 6
  456 / 10 = 45
  45 % 10  = 5 → digit, sum = 11
  45 / 10  = 4
  4  % 10  = 4 → digit, sum = 15
  4  / 10  = 0 → stop (n == 0)
```

## Requirements

- `SumOfDigits(123)` → `6`
- `SumOfDigits(0)` → `0`
- Handle negative numbers by treating them as positive: `SumOfDigits(-42)` → `6`

## Hints

<details>
<summary>Hint 1 — Handle negatives</summary>

Use `if n < 0 { n = -n }` at the start to make the number positive before processing.

</details>

<details>
<summary>Hint 2 — Loop structure</summary>

```go
sum := 0
for n != 0 {
    sum += n % 10
    n /= 10
}
return sum
```

For `n = 0`, the loop doesn't execute and `sum` stays 0 — which is correct.

</details>

## How to Verify

```bash
lncli run
```
