---
title: Maximum of Array
number: 7
difficulty: easy
duration: 10 minutes
concept: Array Traversal, Comparison
---

## What You Need to Build

Write a function called `MaxOfArray` that takes a slice of integers and returns the largest value in the slice.

## Function Signature

```go
func MaxOfArray(nums []int) int
```

## Examples

| Input | Output |
|-------|--------|
| `[3, 1, 4, 1, 5, 9, 2, 6]` | `9` |
| `[1]` | `1` |
| `[-3, -1, -4]` | `-1` |
| `[0, 0, 0]` | `0` |
| `[100]` | `100` |

## Key Concept: Linear Scan

The classic approach: assume the first element is the max, then walk through the rest. Whenever you find a larger value, update your assumption.

```
nums = [3, 1, 4, 1, 5, 9, 2, 6]
         ^
max = 3  (start here)

Walk forward:
  1 < 3 → no change
  4 > 3 → max = 4
  1 < 4 → no change
  5 > 4 → max = 5
  9 > 5 → max = 9
  2 < 9 → no change
  6 < 9 → no change

Result: 9
```

## Requirements

- Returns the maximum value, not its index
- Works for negative numbers: `[-5, -2, -9]` → `-2`
- Works for a single-element slice
- You can assume the slice is non-empty (no need to handle `nil` or empty input)

## Hints

<details>
<summary>Hint 1 — Initialize max correctly</summary>

Set `max := nums[0]` (the first element), then loop from index 1. This works for all-negative arrays because you're comparing against an actual value, not zero.

</details>

<details>
<summary>Hint 2 — Loop pattern</summary>

```go
max := nums[0]
for _, v := range nums[1:] {
    if v > max {
        max = v
    }
}
return max
```

</details>

## How to Verify

```bash
lncli run
```
