---
title: Even or Odd
number: 1
difficulty: easy
duration: 5-10 minutes
concept: Conditionals, Modulo Operator
---

## What You Need to Build

Write a function called `EvenOrOdd` that takes a single integer and returns the string `"even"` if the number is divisible by 2, or `"odd"` if it is not.

## Function Signature

```go
func EvenOrOdd(n int) string
```

## Examples

| Input | Output |
|-------|--------|
| `2`   | `"even"` |
| `3`   | `"odd"`  |
| `0`   | `"even"` |
| `-4`  | `"even"` |
| `-7`  | `"odd"`  |

## Key Concept: The Modulo Operator

The modulo operator `%` gives you the **remainder** after division.

```
10 % 3 = 1   (10 divided by 3 is 3, with remainder 1)
10 % 2 = 0   (10 divided by 2 is 5, with remainder 0)
 7 % 2 = 1   (7 divided by 2 is 3, with remainder 1)
```

A number is even if dividing it by 2 leaves **no remainder** — i.e., `n % 2 == 0`.

## Requirements

- Handle positive numbers: `EvenOrOdd(4)` → `"even"`
- Handle odd numbers: `EvenOrOdd(7)` → `"odd"`
- Handle zero: `EvenOrOdd(0)` → `"even"` (zero is even)
- Handle negative numbers: `EvenOrOdd(-6)` → `"even"`

## Hints

<details>
<summary>Hint 1 — Structure</summary>

Use an `if`/`else` statement. If `n % 2 == 0`, return `"even"`. Otherwise return `"odd"`.

</details>

<details>
<summary>Hint 2 — Negative numbers</summary>

In Go, `%` on negative numbers follows the sign of the dividend: `-6 % 2 == 0`, `-7 % 2 == -1`. Check for `== 0` rather than `== 1` for odd, and it works for all integers.

</details>

## How to Verify

After editing `task-01-even-odd/main.go`, run:

```bash
lncli run
```

The CLI will execute the test suite and show you which tests pass and which fail. Fix any failures, run again, and once all tests pass your progress is saved automatically.
