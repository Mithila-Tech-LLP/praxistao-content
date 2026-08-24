---
title: Reverse String
number: 3
difficulty: easy
duration: 10 minutes
concept: String Iteration, Runes
---

## What You Need to Build

Write a function called `ReverseString` that takes a string and returns it with the characters in reversed order.

## Function Signature

```go
func ReverseString(s string) string
```

## Examples

| Input | Output |
|-------|--------|
| `"hello"` | `"olleh"` |
| `"Go"` | `"oG"` |
| `""` | `""` |
| `"a"` | `"a"` |
| `"racecar"` | `"racecar"` |

## Key Concept: Strings and Runes in Go

In Go, strings are sequences of **bytes**, not characters. For ASCII text (letters a–z, A–Z, digits, symbols), one character = one byte, so you can index directly.

However, the idiomatic way to reverse a string safely is to convert it to a **rune slice** first. A `rune` represents a single Unicode code point (character):

```go
s := "hello"
runes := []rune(s)  // convert to rune slice
// now runes[0] = 'h', runes[1] = 'e', etc.
```

To build the reversed string:
1. Convert `s` to `[]rune`
2. Walk from the last index down to 0
3. Append each rune to a result slice
4. Convert back to string with `string(result)`

## Requirements

- `ReverseString("hello")` → `"olleh"`
- Empty string returns empty string
- Single character returns the same character
- Works correctly with ASCII input

## Hints

<details>
<summary>Hint 1 — Two-pointer swap</summary>

An elegant approach: convert to `[]rune`, then use two pointers (`left = 0`, `right = len-1`) and swap characters, moving inward until they meet.

</details>

<details>
<summary>Hint 2 — Build by iteration</summary>

Alternatively, range over the runes in reverse and append each to a `[]rune` result. Then `return string(result)`.

</details>

## How to Verify

```bash
lncli run
```
