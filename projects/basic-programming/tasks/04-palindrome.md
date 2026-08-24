---
title: Palindrome
number: 4
difficulty: easy
duration: 10 minutes
concept: String Comparison, Two Pointers
---

## What You Need to Build

Write a function called `IsPalindrome` that returns `true` if the input string reads the same forwards and backwards, `false` otherwise.

## Function Signature

```go
func IsPalindrome(s string) bool
```

## Examples

| Input | Output | Reason |
|-------|--------|--------|
| `"racecar"` | `true` | Same forwards and back |
| `"hello"` | `false` | `"hello"` ≠ `"olleh"` |
| `""`  | `true` | Empty string is a palindrome |
| `"a"` | `true` | Single character always is |
| `"abba"` | `true` | Symmetric |
| `"abcd"` | `false` | Not symmetric |

## Key Concept: Two-Pointer Technique

Instead of reversing the whole string and comparing, you can use two pointers:

```
"racecar"
 ^     ^   r == r ✓ — move inward
  ^   ^    a == a ✓ — move inward
   ^ ^     c == c ✓ — move inward
    ^      middle — done, it's a palindrome
```

Start `left` at index 0 and `right` at the last index. While `left < right`:
- If `s[left] != s[right]`, return `false`
- Increment `left`, decrement `right`

If the loop finishes without finding a mismatch, return `true`.

## Requirements

- Case-sensitive: `"Racecar"` is NOT a palindrome (capital R ≠ lowercase r)
- Empty string → `true`
- Single character → `true`
- Even-length palindromes work: `"abba"` → `true`

## Hints

<details>
<summary>Hint 1 — Reuse your previous work</summary>

You can call `ReverseString(s)` from Task 03 and compare `s == ReverseString(s)`. But try implementing the two-pointer approach — it's more efficient.

</details>

<details>
<summary>Hint 2 — Rune indexing</summary>

Convert to `[]rune` before indexing so you correctly handle any Unicode characters.

</details>

## How to Verify

```bash
lncli run
```
