# Chapter 34: Hash Tables

Hash tables give you O(1) average-case lookup, insert, and delete — the fastest of any data structure for exact-match queries. Go's built-in `map` is a hash table, but understanding how it works under the hood makes you a better engineer: you'll know why maps aren't ordered, why they can't be used in comparisons, and how to handle collisions in your own implementations.

## Table of Contents

1. [How Hash Tables Work](#1-how-hash-tables-work)
2. [Hash Functions](#2-hash-functions)
3. [Collision Resolution](#3-collision-resolution)
4. [Implementing a Hash Map](#4-implementing-a-hash-map)
5. [Go's Built-in Map Internals](#5-gos-built-in-map-internals)
6. [Advanced Patterns](#6-advanced-patterns)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. How Hash Tables Work

A hash table stores key-value pairs. To find a value by key:
1. **Hash the key** — run it through a hash function to get an integer index
2. **Use index** — store/retrieve the value at that array index

```
Key "alice" → hash("alice") % 8 → index 3
Key "bob"   → hash("bob")   % 8 → index 7
Key "carol" → hash("carol") % 8 → index 1

Array: [_, carol, _, alice, _, _, _, bob]
       idx: 0     1   2   3   4   5   6   7

Lookup "alice":
  1. hash("alice") % 8 = 3
  2. return array[3] → "alice's value"
```

The core challenge: **two different keys may hash to the same index** (a collision). Every hash table implementation must handle this.

### Performance:
```
Average case (good hash, low load):
  Insert: O(1)
  Lookup: O(1)
  Delete: O(1)

Worst case (all keys collide, or extremely high load factor):
  All operations: O(n)
```

### Quick Check
> 1. What two steps does a hash table use to find a value by key?
> 2. What is a "collision"?
> 3. What is the worst-case time complexity of a hash table lookup?

---

## 2. Hash Functions

A good hash function:
- Is **deterministic**: same input always gives same output
- Is **fast** to compute
- **Distributes keys uniformly**: minimizes collisions
- Is **avalanche**: small input change → large output change

**Simple integer hash:**
```go
func hashInt(n int, buckets int) int {
    return n % buckets
}
// Problem: poor distribution for clustered keys
// Better: multiply by a prime:
func hashIntBetter(n, buckets int) int {
    h := (n * 2654435761) % buckets  // Knuth multiplicative hashing
    if h < 0 {
        h = -h  // Multiplication can overflow to negative — keep index valid
    }
    return h
}
```

**String hash (FNV — Fowler-Noll-Vo):**
```go
const fnvPrime  = 16777619
const fnvOffset = 2166136261

func hashString(s string, buckets int) int {
    h := fnvOffset
    for i := 0; i < len(s); i++ {
        h ^= int(s[i])  // XOR with byte
        h *= fnvPrime   // Multiply by prime
    }
    if h < 0 {
        h = -h
    }
    return h % buckets
}
```

**Go's built-in hash**: Go uses a randomized hash function seeded at startup (different each run) — this prevents hash-flooding attacks where an adversary crafts inputs that all hash to the same bucket.

```go
import "hash/fnv"

func hashKey(key string) uint64 {
    h := fnv.New64a()
    h.Write([]byte(key))
    return h.Sum64()
}
```

### Quick Check
> 1. What properties make a hash function "good"?
> 2. Why does Go randomize its hash seed on each program run?
> 3. What is a "hash-flooding attack"?

---

## 3. Collision Resolution

### Separate Chaining — each bucket is a linked list
```
bucket 3: → [alice, val1] → [dave, val2] → nil
bucket 7: → [bob, val3] → nil
```

```go
type entry[K comparable, V any] struct {
    key   K
    value V
    next  *entry[K, V]
}

type ChainHashMap[K comparable, V any] struct {
    buckets []*entry[K, V]
    size    int
    hash    func(K) int
}

func (m *ChainHashMap[K, V]) Get(key K) (V, bool) {
    idx := m.hash(key) % len(m.buckets)
    for e := m.buckets[idx]; e != nil; e = e.next {
        if e.key == key {
            return e.value, true
        }
    }
    var zero V
    return zero, false
}

func (m *ChainHashMap[K, V]) Set(key K, val V) {
    idx := m.hash(key) % len(m.buckets)
    for e := m.buckets[idx]; e != nil; e = e.next {
        if e.key == key {
            e.value = val  // Update existing
            return
        }
    }
    // Prepend new entry (O(1)):
    m.buckets[idx] = &entry[K, V]{key: key, value: val, next: m.buckets[idx]}
    m.size++
}
```

**Load factor**: `size / num_buckets`. When > 0.75, rehash (double buckets and re-insert all entries).

### Open Addressing — linear probing
When a bucket is occupied, probe to the next one:

```
Insert "alice" (hashes to 3): bucket[3] is empty → store there
Insert "dave"  (hashes to 3): bucket[3] occupied → try 4 → empty → store at 4

Lookup "dave": hash=3, bucket[3] ≠ dave → probe 4 → bucket[4] = dave ✓
```

```go
type probeEntry[K comparable, V any] struct {
    key   K
    value V
    used  bool
}

type OpenHashMap[K comparable, V any] struct {
    slots []probeEntry[K, V]
    size  int
    hash  func(K) int
}

func (m *OpenHashMap[K, V]) Get(key K) (V, bool) {
    n := len(m.slots)
    idx := m.hash(key) % n

    for i := 0; i < n; i++ {
        probe := (idx + i) % n
        if !m.slots[probe].used {
            break  // Empty slot = key not present
        }
        if m.slots[probe].key == key {
            return m.slots[probe].value, true
        }
    }
    var zero V
    return zero, false
}
```

**Comparison:**

| | Chaining | Open Addressing |
|--|----------|----------------|
| Memory | Extra pointer per entry | No extra pointers |
| Cache | Less cache-friendly | More cache-friendly (data is contiguous) |
| Load | Works well past 0.75 | Degrades rapidly past 0.7 |
| Delete | Simple (unlink) | Complex (need tombstones) |
| Go's map uses | Hybrid (bucket arrays) | — |

### Quick Check
> 1. In separate chaining, what happens when two keys hash to the same bucket?
> 2. In linear probing, how do you find a key that collided?
> 3. What is the load factor?

---

## 4. Implementing a Hash Map

A complete, resizable hash map with separate chaining:

```go
const defaultCapacity = 16
const loadFactorThreshold = 0.75

type HashMap[K comparable, V any] struct {
    buckets  [][]*kv[K, V]
    size     int
    hashFunc func(K, int) int
}

type kv[K comparable, V any] struct {
    key K
    val V
}

func NewHashMap[K comparable, V any](hashFunc func(K, int) int) *HashMap[K, V] {
    return &HashMap[K, V]{
        buckets:  make([][]*kv[K, V], defaultCapacity),
        hashFunc: hashFunc,
    }
}

func (m *HashMap[K, V]) idx(key K) int {
    h := m.hashFunc(key, len(m.buckets))
    if h < 0 {
        h = -h
    }
    return h % len(m.buckets)
}

func (m *HashMap[K, V]) Get(key K) (V, bool) {
    i := m.idx(key)
    for _, entry := range m.buckets[i] {
        if entry.key == key {
            return entry.val, true
        }
    }
    var zero V
    return zero, false
}

func (m *HashMap[K, V]) Set(key K, val V) {
    i := m.idx(key)
    for _, entry := range m.buckets[i] {
        if entry.key == key {
            entry.val = val
            return
        }
    }
    m.buckets[i] = append(m.buckets[i], &kv[K, V]{key, val})
    m.size++

    if float64(m.size)/float64(len(m.buckets)) > loadFactorThreshold {
        m.resize()
    }
}

func (m *HashMap[K, V]) Delete(key K) bool {
    i := m.idx(key)
    for j, entry := range m.buckets[i] {
        if entry.key == key {
            // Remove by swapping with last and truncating:
            last := len(m.buckets[i]) - 1
            m.buckets[i][j] = m.buckets[i][last]
            m.buckets[i] = m.buckets[i][:last]
            m.size--
            return true
        }
    }
    return false
}

func (m *HashMap[K, V]) resize() {
    newCap := len(m.buckets) * 2
    newBuckets := make([][]*kv[K, V], newCap)
    
    oldHashFunc := m.hashFunc
    for _, bucket := range m.buckets {
        for _, entry := range bucket {
            h := oldHashFunc(entry.key, newCap)
            if h < 0 { h = -h }
            idx := h % newCap
            newBuckets[idx] = append(newBuckets[idx], entry)
        }
    }
    m.buckets = newBuckets
}

func (m *HashMap[K, V]) Len() int { return m.size }

func (m *HashMap[K, V]) Keys() []K {
    keys := make([]K, 0, m.size)
    for _, bucket := range m.buckets {
        for _, entry := range bucket {
            keys = append(keys, entry.key)
        }
    }
    return keys
}
```

### Quick Check
> 1. When should a hash map resize?
> 2. What happens to all existing entries when you resize?
> 3. Why multiply the capacity by 2 (not by 1.5 or 3)?

---

## 5. Go's Built-in Map Internals

Go's `map` is not a simple chained hash table — it uses a more sophisticated bucket structure:

```
map[string]int internally:
  
  Header:
    - bucket array pointer
    - count (number of key-value pairs)
    - B (log2 of number of buckets)
    - random hash seed
  
  Each bucket holds 8 key-value pairs:
    - 8 top-hash bytes (quick pre-filter)
    - 8 keys
    - 8 values
    - overflow pointer (for more than 8 in this bucket)
```

> Note: this describes the classic implementation (Go ≤ 1.23). Go 1.24 switched maps to a "Swiss Table" design — different layout, same ideas (buckets of slots + a small hash byte as a pre-filter) and the same O(1) behavior. Everything below still applies.

**Key behaviors to remember:**
```go
// 1. Order is randomized — intentionally non-deterministic:
m := map[int]string{1: "a", 2: "b", 3: "c"}
for k, v := range m {
    fmt.Println(k, v)  // Order varies each run
}

// 2. Concurrent access requires external synchronization:
// Reading while writing: undefined behavior, can crash
var m = make(map[string]int)
go func() { m["key"] = 1 }()  // RACE with line below
go func() { _ = m["key"] }()  // Data race!
// Use sync.RWMutex or sync.Map (Chapter 21)

// 3. Map key must be comparable (==):
// OK:   map[int], map[string], map[struct{a,b int}]
// NOT:  map[[]int], map[map[string]int]  — slices/maps aren't comparable

// 4. Nil map panics on write, returns zero on read:
var m map[string]int
_ = m["key"]   // Returns 0 — safe
m["key"] = 1   // PANIC: assignment to entry in nil map

// 5. Growing a map during iteration is safe (new keys may or may not appear):
for k := range m {
    m[k+"_copy"] = m[k]  // Safe but new keys may not be iterated
}
```

**Map performance:**
```go
// Frequent optimization: check existence ONCE:
// Bad (two lookups):
if _, ok := m[key]; ok {
    use(m[key])
}

// Good (one lookup):
if v, ok := m[key]; ok {
    use(v)
}

// Pre-size maps when you know the count:
m := make(map[string]int, 1000)  // Hint: ~1000 entries
```

---

## 6. Advanced Patterns

**Frequency counter:**
```go
func WordFrequency(words []string) map[string]int {
    freq := make(map[string]int, len(words))
    for _, w := range words {
        freq[w]++  // Zero value of int is 0 — works on first access
    }
    return freq
}
```

**Grouping:**
```go
func GroupByLength(words []string) map[int][]string {
    groups := make(map[int][]string)
    for _, w := range words {
        groups[len(w)] = append(groups[len(w)], w)
    }
    return groups
}
```

**Set (using `map[T]struct{}`):**
```go
type Set[T comparable] map[T]struct{}

func (s Set[T]) Add(v T)          { s[v] = struct{}{} }
func (s Set[T]) Contains(v T) bool { _, ok := s[v]; return ok }
func (s Set[T]) Remove(v T)        { delete(s, v) }

func Intersection[T comparable](a, b Set[T]) Set[T] {
    result := make(Set[T])
    for k := range a {
        if b.Contains(k) {
            result.Add(k)
        }
    }
    return result
}
```

**Memoization:**
```go
func MemoizeFib() func(int) int {
    cache := make(map[int]int)
    var fib func(int) int
    fib = func(n int) int {
        if n <= 1 {
            return n
        }
        if v, ok := cache[n]; ok {
            return v
        }
        result := fib(n-1) + fib(n-2)
        cache[n] = result
        return result
    }
    return fib
}
```

**Two-sum (classic interview):**
```go
// TwoSum finds indices of two numbers that add up to target.
func TwoSum(nums []int, target int) (int, int) {
    seen := make(map[int]int)  // value → index
    for i, n := range nums {
        complement := target - n
        if j, ok := seen[complement]; ok {
            return j, i
        }
        seen[n] = i
    }
    return -1, -1
}
```

---

## Summary

- **Hash table**: key → hash function → index → array slot; O(1) average for all ops
- **Hash function**: deterministic, fast, uniform distribution, avalanche effect
- **Separate chaining**: each bucket = linked list; handles any load factor
- **Open addressing**: probe linearly on collision; cache-friendly, degrades >0.7 load
- **Load factor**: `size/capacity`; resize when > 0.75 (double the array, rehash everything)
- **Go map**: 8-entry buckets, randomized hash seed, non-deterministic iteration, concurrent reads safe but not concurrent reads+writes
- **Nil map**: reads return zero; writes panic — always `make(map[K]V)` before writing
- **Patterns**: frequency counter (m[k]++), grouping, sets (map[T]struct{}), memoization, two-sum

---

## Exercises

### Easy
1. Write `Anagram(a, b string) bool` using a frequency map. Two strings are anagrams if they have the same character frequencies. O(n) time.
2. Write `FirstNonRepeating(s string) rune` using a map to count characters, then iterate the original string to find the first with count 1.
3. Write `LongestConsecutive(nums []int) int` using a set. Insert all numbers, then for each number n, if n-1 is NOT in the set, count up from n. Return the longest sequence found.

### Medium
4. LRU Cache: Implement an LRU (Least Recently Used) cache with O(1) `Get` and `Put`. Use: a `map[int]*DNode[int]` for O(1) lookup, and a doubly linked list to maintain usage order (most recent at front). `Put(key, value)` — if full, evict the least recently used. `Get(key) int` — return -1 if not found, otherwise move to front. Test with capacity=2: put(1,1), put(2,2), get(1)→1, put(3,3)→evicts 2, get(2)→-1.
5. Consistent hashing: Implement consistent hashing for distributing keys across N nodes. A hash ring maps nodes to positions. `AddNode(name string)`, `RemoveNode(name string)`, `GetNode(key string) string` returns which node owns the key. Use virtual nodes (100 per physical node) to ensure even distribution. Test: add 3 nodes, distribute 1000 keys, remove 1 node — verify ≈ 1/3 of keys migrate (not all of them).
6. HashMap with O(1) random access: Implement `RandomizedSet` supporting `Insert(val int) bool`, `Remove(val int) bool`, `GetRandom() int` — all in O(1). Use a hash map + dynamic array. The tricky part: O(1) remove without preserving order (swap with last element in array).

### Hard
7. External hash table: Implement a hash table that stores values on "disk" (simulate with a []byte buffer). `Set(key string, value []byte)`, `Get(key string) ([]byte, bool)`. Use a fixed number of pages (4KB each), separate chaining with a page-based linked list, and an in-memory directory (hash → page number). Handle page overflow by allocating new overflow pages. Test with 10,000 entries, verify all are retrievable.
8. Cuckoo hashing: Implement a hash table using cuckoo hashing — each key has two possible positions (from two different hash functions). On collision, the existing key is "kicked out" to its other position (displacing that key to its other position, and so on). This gives worst-case O(1) lookup. Implement with a displacement limit (cycle detection) and rehash on failure. Benchmark against chaining: cuckoo should have better worst-case lookup at the cost of more complex insert.
