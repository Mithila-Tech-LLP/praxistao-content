---
title: Count Characters
number: 8
difficulty: easy
duration: 10 minutes
concept: String Iteration, Counting
---

## What You Need to Build

Write a function called `CountCharacters` that counts how many times each character appears in a string and returns the result as a map.

## Function Signature

```go
func CountCharacters(s string) map[rune]int
```

## Examples

| Input | Output |
|-------|--------|
| `"hello"` | `{'h':1, 'e':1, 'l':2, 'o':1}` |
| `"aab"` | `{'a':2, 'b':1}` |
| `""` | `{}` (empty map) |
| `"aaa"` | `{'a':3}` |

## Key Concept: Maps in Go

A Go map stores key-value pairs. For character counting, the key is a `rune` (character) and the value is an `int` (count).

```go
counts := make(map[rune]int)
```

When you read a key that doesn't exist yet in a Go map, it returns the **zero value** for its type — for `int`, that's `0`. This makes counting very clean:

```go
counts[ch]++   // first time: 0 + 1 = 1, second time: 1 + 1 = 2, etc.
```

No need to check if the key exists first.

**Iterating over a string with `range`** gives you each character as a `rune`:
```go
for _, ch := range "hello" {
    // ch is a rune: 'h', 'e', 'l', 'l', 'o'
}
```

## Requirements

- Count every character including spaces and symbols
- Return an empty map (not `nil`) for an empty string
- Use `rune` as the map key type (handles Unicode correctly)

## Hints

<details>
<summary>Hint 1 — Full solution structure</summary>

```go
func CountCharacters(s string) map[rune]int {
    counts := make(map[rune]int)
    for _, ch := range s {
        counts[ch]++
    }
    return counts
}
```

</details>

<details>
<summary>Hint 2 — Checking the result</summary>

Maps in Go don't have a guaranteed order, so when your tests verify the result, they check specific keys: `result['l'] == 2`, not the full map string.

</details>

## How to Verify

```bash
lncli run
```
