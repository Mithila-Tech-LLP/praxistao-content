---
title: Hash Map
step: 5
difficulty: medium
estimated: 45 min
---

## What You Are Building

A hash map (also called hash table or dictionary) provides O(1) average-case get, set, and delete using a hash function to map keys to bucket indices. It is the most widely used data structure in all of software.

```
Set("name", 42)
Set("age",  25)
Get("name") → 42, true
Get("city") → 0, false
```

## Key Concepts

**Hash functions** — A hash function deterministically maps an arbitrary key to an integer. A good hash function distributes keys uniformly across buckets. We'll use a simple version of the **djb2** algorithm:

```go
func hash(key string) int {
    h := 5381
    for _, c := range key {
        h = h*33 + int(c)
    }
    return h
}
```

To get a bucket index: `bucketIndex := hash(key) % len(buckets)`. Always take the absolute value: `if bucketIndex < 0 { bucketIndex = -bucketIndex }`.

**Separate chaining** — Multiple keys can hash to the same bucket (a *collision*). We handle this by storing a slice (or linked list) of entries at each bucket. When getting a key, we scan the bucket's entries for an exact key match.

```
Bucket 0: []
Bucket 1: [("name", 42), ("city", 7)]  ← collision
Bucket 2: [("age", 25)]
...
```

**Load factor** — Performance degrades as the map fills up. A real implementation resizes (rehashes) when `len(entries) / len(buckets) > 0.75`. For this task, a fixed 16 buckets is fine.

## Struct Signatures

```go
type entry struct {
    key string
    val int
}

type HashMap struct {
    buckets [][]entry
    size    int
}

func NewHashMap() *HashMap {
    return &HashMap{
        buckets: make([][]entry, 16),
    }
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Set(key string, val int)` | Insert or update the key |
| `Get(key string) (int, bool)` | Return value; ok=false if missing |
| `Delete(key string)` | Remove the key if present |
| `Keys() []string` | Return all current keys (any order) |

## Edge Cases to Handle

- `Get` a key that doesn't exist: return `0, false`
- `Set` an existing key: update the value, don't create a duplicate entry
- `Delete` a key that doesn't exist: no-op, no panic
- `Keys` on empty map: return empty slice

## Example

```go
m := NewHashMap()
m.Set("go", 1)
m.Set("python", 2)
m.Set("rust", 3)

val, ok := m.Get("go")
fmt.Println(val, ok) // 1 true

m.Set("go", 99)
val, _ = m.Get("go")
fmt.Println(val) // 99

m.Delete("python")
_, ok = m.Get("python")
fmt.Println(ok) // false

fmt.Println(len(m.Keys())) // 2
```

## Hints

- Write a private `bucketIndex(key string) int` helper so all methods agree on which bucket to use.
- In `Set`, scan the bucket first. If you find the key, update its value and return. Otherwise append a new entry.
- In `Delete`, rebuild the bucket slice by copying entries that don't match the key.
- `Keys`: loop all buckets, loop all entries in each bucket, collect `e.key`.
