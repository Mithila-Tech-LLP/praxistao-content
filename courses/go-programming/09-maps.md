# Chapter 09: Maps

A **map** is Go's built-in key-value store — also called a hash map or dictionary in other languages. Maps let you associate any key with any value and retrieve values by key in constant time O(1). They are used everywhere: caching results, counting occurrences, grouping data, building indexes, and much more.

## Table of Contents

1. [Creating and Using Maps](#1-creating-and-using-maps)
2. [Reading from Maps Safely](#2-reading-from-maps-safely)
3. [Deleting and Iterating](#3-deleting-and-iterating)
4. [Maps Internals](#4-maps-internals)
5. [Common Map Patterns](#5-common-map-patterns)
6. [Maps as Sets](#6-maps-as-sets)
7. [Nested Maps](#7-nested-maps)
8. [Concurrency and Maps](#8-concurrency-and-maps)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Creating and Using Maps

**Declaration and initialization:**
```go
// Empty map (NOT nil — ready to use)
ages := map[string]int{}

// Map literal:
capitals := map[string]string{
    "India":   "New Delhi",
    "France":  "Paris",
    "Japan":   "Tokyo",
    "Germany": "Berlin",  // trailing comma required on last line
}

// make (preferred when size is unknown but expected to grow):
scores := make(map[string]int)
scores["Alice"] = 95
scores["Bob"] = 87
scores["Carol"] = 91
```

**Key types — anything comparable:**
```go
// Valid key types (must be comparable with ==):
map[string]int
map[int]string
map[float64]bool
map[[3]int]string  // Arrays are comparable — valid key
// map[[]int]string  // ERROR: slices are not comparable

// Struct keys (if all fields are comparable):
type Point struct{ X, Y int }
distances := map[Point]float64{
    {0, 0}: 0,
    {3, 4}: 5.0,
}
```

**Adding and updating:**
```go
m := map[string]int{}

m["Alice"] = 95   // Add new key
m["Alice"] = 97   // Update existing key (no error, just overwrites)
m["Bob"] = 87

fmt.Println(len(m))  // 2
```

### Quick Check
> 1. What is the difference between `var m map[string]int` and `m := map[string]int{}`?
> 2. Can you use a slice as a map key? Can you use a struct?
> 3. What happens when you assign to a key that already exists?

---

## 2. Reading from Maps Safely

**The two-value form** — always use this to check if a key exists:
```go
ages := map[string]int{"Alice": 30, "Bob": 25}

// One-value form: returns zero value if key missing
age := ages["Alice"]    // 30
missing := ages["Dave"] // 0 — zero value, but key doesn't exist!

// Two-value form: tells you if the key existed
age, ok := ages["Alice"]
if ok {
    fmt.Printf("Alice is %d years old\n", age)
} else {
    fmt.Println("Alice not found")
}

// Compact form (init statement in if):
if age, ok := ages["Dave"]; ok {
    fmt.Println("Dave is", age)
} else {
    fmt.Println("Dave not found")  // This branch runs
}
```

**Nil map: reads are safe, writes panic:**
```go
var m map[string]int  // nil map

// SAFE: reading from nil map returns zero value
fmt.Println(m["Alice"])  // 0 (no panic)
v, ok := m["Alice"]      // 0, false (no panic)

// PANIC: writing to nil map panics!
m["Alice"] = 30  // panic: assignment to entry in nil map

// Always initialize before writing:
m = make(map[string]int)
m["Alice"] = 30  // Safe now
```

**Default values pattern:**
```go
// Count word frequencies (zero value of int is 0, so this works!)
words := []string{"go", "is", "great", "go", "is", "go"}
freq := make(map[string]int)

for _, w := range words {
    freq[w]++  // If key doesn't exist, zero value (0) is returned, then incremented
}
// freq = {"go": 3, "is": 2, "great": 1}
```

### Quick Check
> 1. What does `m["missing"]` return if "missing" is not in the map?
> 2. Is reading from a nil map safe? What about writing?
> 3. Why does `freq[word]++` work even for words not yet in the map?

---

## 3. Deleting and Iterating

**Deleting:**
```go
ages := map[string]int{"Alice": 30, "Bob": 25, "Carol": 35}

delete(ages, "Bob")            // Remove key "Bob"
delete(ages, "nonexistent")    // Safe: deleting missing key does nothing

fmt.Println(len(ages))  // 2

// Clear all entries (Go 1.21+ built-in):
clear(ages)

// Or in a loop (works on any Go version):
for k := range ages {
    delete(ages, k)
}
```

**Iterating:**
```go
ages := map[string]int{"Alice": 30, "Bob": 25, "Carol": 35}

// Iteration order is RANDOM — Go deliberately randomizes this
for name, age := range ages {
    fmt.Printf("%s: %d\n", name, age)
}

// Keys only:
for name := range ages {
    fmt.Println(name)
}

// Sorted iteration (sort keys first):
import "slices"

names := make([]string, 0, len(ages))
for k := range ages {
    names = append(names, k)
}
slices.Sort(names)

for _, name := range names {
    fmt.Printf("%s: %d\n", name, ages[name])
}
```

**`maps` package (Go 1.21+):**
```go
import "maps"

a := map[string]int{"x": 1, "y": 2}
b := map[string]int{"y": 2, "z": 3}

// Clone:
c := maps.Clone(a)  // c is a shallow copy of a (keys/values copied by assignment)

// Equal:
fmt.Println(maps.Equal(a, b))  // false

// Copy (merge b into a):
maps.Copy(a, b)
// a = {"x": 1, "y": 2, "z": 3}

// Collect keys and values (Go 1.23+ — maps.Keys/Values return iterators):
keys := slices.Collect(maps.Keys(a))
vals := slices.Collect(maps.Values(a))
```

### Quick Check
> 1. What happens if you `delete` a key that doesn't exist?
> 2. Is map iteration order guaranteed in Go?
> 3. How do you iterate over a map in sorted key order?

---

## 4. Maps Internals

Understanding how maps work helps you use them correctly:

```
Hash Map Structure:

Key "Alice" ──hash()──→ bucket index 3
Key "Bob"   ──hash()──→ bucket index 7
Key "Carol" ──hash()──→ bucket index 3 (collision! — stored in same bucket)

Buckets array:
  [0]: empty
  [1]: empty
  [2]: empty
  [3]: "Alice"→30, "Carol"→35  (both hash to 3)
  [4]: empty
  ...
  [7]: "Bob"→25

Lookup "Alice":
  1. hash("Alice") → 3
  2. Check bucket 3
  3. Compare keys: "Alice" == "Alice" → found! return 30
```

**Key points about Go's map implementation:**
- Go's map is a **hash map** with buckets (each bucket holds 8 key-value pairs)
- When a bucket is full, it **overflows** to a linked bucket
- When load factor exceeds ~6.5 entries/bucket, the map **grows** (rehashes all entries)
- Map growth is O(n) amortized but individual operations are O(1) on average

**Maps are reference types** — passing a map to a function shares the underlying data:
```go
func addEntry(m map[string]int, key string, value int) {
    m[key] = value  // Modifies the original map
}

scores := map[string]int{"Alice": 95}
addEntry(scores, "Bob", 87)
fmt.Println(scores)  // map[Alice:95 Bob:87] — Bob was added
```

**Maps are NOT safe for concurrent access** (covered in Chapter 21):
```go
// RACE CONDITION: two goroutines reading/writing the same map
go func() { m["key"] = 1 }()  // goroutine 1 writes
go func() { _ = m["key"] }()  // goroutine 2 reads
// This will panic at runtime: concurrent map read and map write
```

### Quick Check
> 1. What is a hash collision and how does Go handle it?
> 2. If you pass a map to a function and the function modifies it, is the original changed?
> 3. Is Go's built-in map safe for concurrent use by multiple goroutines?

---

## 5. Common Map Patterns

**Grouping (group items by a key):**
```go
type Product struct {
    Name     string
    Category string
    Price    float64
}

products := []Product{
    {"iPhone", "Electronics", 999},
    {"MacBook", "Electronics", 1299},
    {"Chair", "Furniture", 299},
    {"Desk", "Furniture", 499},
    {"Keyboard", "Electronics", 129},
}

// Group by category:
byCategory := make(map[string][]Product)
for _, p := range products {
    byCategory[p.Category] = append(byCategory[p.Category], p)
}

for category, items := range byCategory {
    fmt.Printf("%s: %d items\n", category, len(items))
}
// Electronics: 3 items
// Furniture: 2 items
```

**Caching / Memoization:**
```go
var cache = make(map[int]int)

func fibonacci(n int) int {
    if n <= 1 {
        return n
    }
    if v, ok := cache[n]; ok {
        return v
    }
    result := fibonacci(n-1) + fibonacci(n-2)
    cache[n] = result
    return result
}
```

**Frequency counting:**
```go
func charFrequency(s string) map[rune]int {
    freq := make(map[rune]int)
    for _, r := range s {
        freq[r]++
    }
    return freq
}

freq := charFrequency("hello world")
// 'h':1, 'e':1, 'l':3, 'o':2, ' ':1, 'w':1, 'r':1, 'd':1
```

**Existence check (not caring about value):**
```go
// Use struct{} as value — zero memory, just need the key
visited := make(map[string]struct{})

visited["pageA"] = struct{}{}
visited["pageB"] = struct{}{}

if _, ok := visited["pageA"]; ok {
    fmt.Println("Already visited")
}
```

**Index map:**
```go
// Fast lookup of struct by ID:
type User struct {
    ID   int
    Name string
    Email string
}

users := []User{
    {1, "Alice", "alice@example.com"},
    {2, "Bob", "bob@example.com"},
}

// Build index:
userByID := make(map[int]*User, len(users))
for i := range users {
    userByID[users[i].ID] = &users[i]
}

// Now O(1) lookup:
if u, ok := userByID[2]; ok {
    fmt.Println(u.Name)  // Bob
}
```

### Quick Check
> 1. How do you group a slice of items by some field using a map?
> 2. What is `map[string]struct{}` used for and why `struct{}`?
> 3. Why store `*User` instead of `User` as map values when users is a slice?

---

## 6. Maps as Sets

Go doesn't have a built-in set type. Maps are used to simulate sets:

```go
// Set using map[T]struct{} (most memory-efficient):
type Set[T comparable] struct {
    m map[T]struct{}
}

func NewSet[T comparable]() *Set[T] {
    return &Set[T]{m: make(map[T]struct{})}
}

func (s *Set[T]) Add(v T) {
    s.m[v] = struct{}{}
}

func (s *Set[T]) Contains(v T) bool {
    _, ok := s.m[v]
    return ok
}

func (s *Set[T]) Remove(v T) {
    delete(s.m, v)
}

func (s *Set[T]) Size() int {
    return len(s.m)
}

// Usage:
words := NewSet[string]()
words.Add("go")
words.Add("python")
words.Add("go")  // Duplicate — ignored
fmt.Println(words.Size())         // 2
fmt.Println(words.Contains("go")) // true
```

**Set operations:**
```go
func Intersection[T comparable](a, b *Set[T]) *Set[T] {
    result := NewSet[T]()
    for k := range a.m {
        if b.Contains(k) {
            result.Add(k)
        }
    }
    return result
}

func Union[T comparable](a, b *Set[T]) *Set[T] {
    result := NewSet[T]()
    for k := range a.m { result.Add(k) }
    for k := range b.m { result.Add(k) }
    return result
}

func Difference[T comparable](a, b *Set[T]) *Set[T] {
    result := NewSet[T]()
    for k := range a.m {
        if !b.Contains(k) {
            result.Add(k)
        }
    }
    return result
}
```

### Quick Check
> 1. Why does Go use `map[T]struct{}` for sets instead of `map[T]bool`?
> 2. How do you implement set intersection using maps?
> 3. Write the type constraint for a map key (one word hint: comparable).

---

## 7. Nested Maps

```go
// Map of maps (like a 2D map):
// permissions[userID][resource] = allowed
permissions := map[string]map[string]bool{
    "alice": {
        "read":  true,
        "write": true,
        "admin": false,
    },
    "bob": {
        "read":  true,
        "write": false,
        "admin": false,
    },
}

// Safe nested access (check outer key first!):
func canAccess(permissions map[string]map[string]bool, user, action string) bool {
    if perms, ok := permissions[user]; ok {
        return perms[action]  // returns false if action not found (zero value)
    }
    return false
}

fmt.Println(canAccess(permissions, "alice", "write"))  // true
fmt.Println(canAccess(permissions, "bob", "admin"))    // false
fmt.Println(canAccess(permissions, "charlie", "read")) // false (user missing)
```

**Building nested maps dynamically:**
```go
// Count word occurrences per document:
wordsByDoc := make(map[string]map[string]int)

docs := map[string]string{
    "doc1": "go is great go",
    "doc2": "python is fun go is",
}

for docName, content := range docs {
    wordsByDoc[docName] = make(map[string]int)  // Initialize inner map!
    for _, word := range strings.Fields(content) {
        wordsByDoc[docName][word]++
    }
}
// wordsByDoc["doc1"]["go"] = 2
// wordsByDoc["doc2"]["is"] = 2
```

**Important: initialize inner maps before writing:**
```go
var m map[string]map[string]int  // nil outer map

// DON'T do this — outer map is nil:
// m["key1"]["key2"] = 1  // panic: assignment to entry in nil map

// DO:
if m == nil {
    m = make(map[string]map[string]int)
}
if m["key1"] == nil {
    m["key1"] = make(map[string]int)
}
m["key1"]["key2"] = 1
```

### Quick Check
> 1. What must you do before writing to the inner map of a nested map?
> 2. What is `map[string]map[string]bool` useful for?
> 3. Why is `m["missing"]["key"]` safe (no panic) when accessing a missing outer key?

---

## 8. Concurrency and Maps

Go's built-in map is **not safe for concurrent use**. If multiple goroutines read and write simultaneously, the program panics or produces corrupted data.

**`sync.Map` — thread-safe map:**
```go
import "sync"

var sm sync.Map

// Store:
sm.Store("key1", "value1")
sm.Store("key2", 42)

// Load:
if v, ok := sm.Load("key1"); ok {
    fmt.Println(v.(string))  // "value1" (requires type assertion)
}

// LoadOrStore — atomic load-or-insert:
actual, loaded := sm.LoadOrStore("key3", "default")
fmt.Println(actual, loaded)  // "default", false (was not there before)

// Delete:
sm.Delete("key1")

// Range:
sm.Range(func(key, value interface{}) bool {
    fmt.Printf("%v: %v\n", key, value)
    return true  // return false to stop iteration
})
```

**When to use `sync.Map`:**
- When keys are written once but read many times (cache-heavy workloads)
- When goroutines don't share keys (each goroutine writes distinct keys)
- For everything else: use a regular map protected by `sync.RWMutex` (covered in Chapter 21)

**Regular map + RWMutex (often better than sync.Map):**
```go
type SafeMap struct {
    mu   sync.RWMutex
    data map[string]int
}

func (m *SafeMap) Set(key string, value int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.data[key] = value
}

func (m *SafeMap) Get(key string) (int, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    v, ok := m.data[key]
    return v, ok
}
```

### Quick Check
> 1. What happens if two goroutines write to the same map concurrently?
> 2. What is `sync.Map` and when should you use it?
> 3. What is the advantage of `sync.RWMutex` over a regular `sync.Mutex` for a read-heavy map?

---

## Summary

- **Creation**: `map[K]V{}` or `make(map[K]V)` — always initialize before writing; nil map panics on write
- **Two-value read**: `v, ok := m[key]` — always use `ok` when unsure if key exists
- **Zero value on miss**: `m["missing"]` returns zero value, never panics (for initialized map)
- **Delete**: `delete(m, key)` — safe even for missing keys
- **Iteration**: `for k, v := range m` — order is random; sort keys for deterministic output
- **Internals**: hash map with buckets; reference type (shared on function pass); grows automatically
- **Not concurrent**: use `sync.Map` or `sync.RWMutex` for concurrent access
- **Patterns**: grouping, frequency counting, memoization, sets (`map[T]struct{}`), indexes

---

## Exercises

### Easy
1. Write a function `wordCount(s string) map[string]int` that counts how many times each word appears in a string (case-insensitive, strip punctuation).
2. Write a function `invertMap(m map[string]int) map[int]string` that swaps keys and values. What happens if two keys have the same value?
3. Write a function `intersection(a, b []string) []string` that returns elements present in both slices, using a map for O(n) time.

### Medium
4. Anagram detection: Write `groupAnagrams(words []string) [][]string` that groups anagrams together. "eat", "tea", "tan", "ate", "nat", "bat" → [["eat","tea","ate"], ["tan","nat"], ["bat"]]. Key insight: anagrams have the same sorted letters.
5. LRU Cache: Implement an LRU (Least Recently Used) cache with capacity N using a map + doubly linked list. It must support: `Get(key int) int` (returns -1 if not found) and `Put(key, value int)` (adds or updates; evicts least recently used when at capacity). Both operations O(1).
6. Two Sum: Classic algorithm problem using maps: given a slice of integers and a target sum, return the indices of the two numbers that add up to target. Example: `[2,7,11,15]`, target=9 → `[0,1]` (because 2+7=9). Solve in O(n) time using a map. Then extend to return ALL pairs, and then to Three Sum (three numbers summing to target).

### Hard
7. Concurrent word index: Build a concurrent document indexer. Given a directory path: (a) Scan all `.txt` files concurrently (one goroutine per file). (b) Each goroutine builds a local word frequency map. (c) Merge all local maps into a global inverted index: `map[word]map[filename]count`. (d) Use `sync.Mutex` to protect the merge step. (e) Support queries: `Search(word string) []FileMatch` where `FileMatch = {Filename, Count}`, sorted by count descending. Benchmark sequential vs concurrent with 100 files of 10,000 words each.
8. Graph as adjacency map: Represent a weighted directed graph as `map[string]map[string]float64` where `graph[src][dst] = weight`. Implement: `AddEdge(src, dst string, weight float64)`, `RemoveEdge(src, dst string)`, `Neighbors(node string) []string`, `BFS(start string) []string` (nodes in BFS order), `DFS(start string) []string`, `HasCycle() bool` (using DFS with visited/in-stack tracking), `TopologicalSort() ([]string, error)` (returns error if cycle detected). Test with a dependency graph (packages that depend on each other).
