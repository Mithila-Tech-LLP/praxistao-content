---
title: Frequency Map
number: 10
difficulty: medium
duration: 15 minutes
concept: Maps, Sorting, Slices
---

## What You Need to Build

Write a function called `FrequencyMap` that takes a slice of strings and returns a map showing how many times each string appears.

Then write a second function `TopN` that takes the frequency map and an integer `n`, and returns the top `n` most frequent strings in descending order of frequency. Break ties alphabetically.

## Function Signatures

```go
func FrequencyMap(words []string) map[string]int
func TopN(freq map[string]int, n int) []string
```

## Examples

```
words = ["apple", "banana", "apple", "cherry", "banana", "apple"]

FrequencyMap(words) → {"apple": 3, "banana": 2, "cherry": 1}

TopN(freq, 2) → ["apple", "banana"]
TopN(freq, 3) → ["apple", "banana", "cherry"]
TopN(freq, 1) → ["apple"]
```

Tie-breaking example:
```
freq = {"cat": 2, "ant": 2, "bat": 2}
TopN(freq, 2) → ["ant", "bat"]  // alphabetical among ties
```

## Key Concept: Sorting with Custom Comparators

`FrequencyMap` is straightforward — same pattern as Task 08. The interesting part is `TopN`.

To sort by frequency descending, then alphabetically for ties, use `sort.Slice`:

```go
import "sort"

keys := []string{"apple", "banana", "cherry"}
sort.Slice(keys, func(i, j int) bool {
    if freq[keys[i]] != freq[keys[j]] {
        return freq[keys[i]] > freq[keys[j]]  // higher count first
    }
    return keys[i] < keys[j]  // alphabetical for ties
})
```

## Requirements

- `FrequencyMap` returns an empty map for empty input
- `TopN` returns at most `n` items (fewer if there are fewer unique words)
- Ties broken alphabetically (ascending)
- Results sorted by frequency descending

## Hints

<details>
<summary>Hint 1 — Building the top-N list</summary>

1. Extract all keys from the map into a `[]string`
2. Sort using `sort.Slice` with the comparator above
3. Return the first `n` elements (use `min(n, len(keys))` to avoid index out of bounds)

</details>

<details>
<summary>Hint 2 — Getting map keys</summary>

```go
keys := make([]string, 0, len(freq))
for k := range freq {
    keys = append(keys, k)
}
```

</details>

<details>
<summary>Hint 3 — Capping at n</summary>

```go
if n > len(keys) {
    n = len(keys)
}
return keys[:n]
```

</details>

## How to Verify

```bash
lncli run
```
