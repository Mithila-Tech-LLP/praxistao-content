---
title: Bloom Filter
step: 10
difficulty: hard
estimated: 45 min
---

## What You Are Building

A bloom filter answers the question "have I seen this item before?" in O(1) time using a tiny bit array — far smaller than storing the items themselves. The catch: it can have **false positives** (say "yes" when the answer is "no") but never false negatives (never says "no" when the answer is "yes").

Used in: databases (avoid disk lookups for nonexistent keys), web caches, spam filters, distributed systems.

## Key Concepts

**Bit array** — Instead of storing items, store a `[]bool` (or `[]byte`) of fixed size. A set item flips certain bits to `true`. A query checks whether all those bits are `true`.

**Multiple hash functions** — To reduce false positives, hash each item with `numHashes` different hash functions and set all corresponding bits. To check membership, verify all those bits are set.

**Simulating multiple hash functions with seeds** — We don't need `numHashes` totally different algorithms. A common technique is to run one base hash (FNV-1a) with different numeric seeds mixed in:

```go
func hashWithSeed(item string, seed uint32) uint32 {
    h := fnv.New32a()
    // write the seed bytes first, then the item
    h.Write([]byte{byte(seed), byte(seed >> 8), byte(seed >> 16), byte(seed >> 24)})
    h.Write([]byte(item))
    return h.Sum32()
}
```

Then: `bitIndex = hashWithSeed(item, uint32(i)) % uint32(len(bf.bits))`

**False positive rate** — With `m` bits and `k` hash functions and `n` inserted items, the false positive rate is approximately `(1 - e^(-kn/m))^k`. More bits and more hash functions = lower false positive rate, but higher memory and CPU cost.

**No delete** — You cannot delete from a basic bloom filter. Flipping a bit back might have been set by a different item.

## Struct Signatures

```go
type BloomFilter struct {
    bits      []bool
    numHashes int
}

func NewBloomFilter(size, numHashes int) *BloomFilter {
    return &BloomFilter{
        bits:      make([]bool, size),
        numHashes: numHashes,
    }
}
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Add(item string)` | Hash item with all seeds; set corresponding bits |
| `MightContain(item string) bool` | True if all bits for this item are set |

## Edge Cases to Handle

- `MightContain` on a freshly created filter: return `false` (no bits set)
- Items added are always found (`MightContain` must return `true` for added items — no false negatives)
- `MightContain` may return `true` for items not added (false positives are acceptable)

## Example

```go
bf := NewBloomFilter(1000, 3)

bf.Add("apple")
bf.Add("banana")
bf.Add("cherry")

fmt.Println(bf.MightContain("apple"))   // true  (definitely)
fmt.Println(bf.MightContain("banana"))  // true  (definitely)
fmt.Println(bf.MightContain("durian"))  // false (probably — might be true due to false positive)

// After adding, must always contain:
bf.Add("elderberry")
fmt.Println(bf.MightContain("elderberry")) // always true
```

## Hints

- Import `hash/fnv` from stdlib. `fnv.New32a()` returns a `hash.Hash32`.
- For each of the `numHashes` seeds (0 to numHashes-1), compute a bit index and set (or check) that bit.
- The bit index must wrap around: `index := int(hashWithSeed(item, uint32(i))) % len(bf.bits)`. Handle negative modulo by taking the absolute value or using unsigned arithmetic.
- `MightContain` returns `false` as soon as any bit is unset — early exit.
- A size of 1000 bits with 3 hash functions gives roughly 1% false positive rate after 100 insertions.
