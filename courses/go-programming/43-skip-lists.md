# Chapter 43: Skip Lists and Concurrent Data Structures

A skip list is a probabilistic data structure that achieves O(log n) average search, insert, and delete — the same as a balanced BST — without complex rotations. Its layered structure also makes it easier to implement a concurrent version than a balanced tree.

## Table of Contents

1. [Skip List](#1-skip-list)
2. [Concurrent Skip List](#2-concurrent-skip-list)
3. [sync.Map — When to Use It](#3-syncmap--when-to-use-it)
4. [Concurrent Patterns](#4-concurrent-patterns)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. Skip List

### Concept

A skip list is multiple layers of linked lists:
- Level 0: a full linked list of all elements in sorted order
- Level 1: a subset of level 0 (each element promoted with probability p=0.5)
- Level 2: a subset of level 1
- ...

To search for key k, start at the top level. At each node, if the next node's key is ≤ k, move forward; otherwise, drop down a level.

```
Level 3: ─────────────────────────────── 50 ────────────────────
Level 2: ──── 10 ─────────────────────── 50 ────────────────────
Level 1: ──── 10 ────── 30 ──────────── 50 ──── 70 ─────────────
Level 0: 1 ── 10 ── 20 ─ 30 ── 40 ──── 50 ── 60 ─ 70 ── 80 ──
```

Search for 60:
1. Level 3: next=50 ≤ 60, move forward. next=∞ > 60, drop down.
2. Level 2: next=∞ > 60, drop down.
3. Level 1: next=70 > 60, drop down.
4. Level 0: next=60 — found!

Expected comparisons: O(log n).

```go
package skiplist

import "math/rand"

const maxLevel = 16
const probability = 0.5

type node struct {
    key     int
    value   any
    forward []*node // forward[i] = next node at level i
}

func newNode(key int, value any, level int) *node {
    return &node{
        key:     key,
        value:   value,
        forward: make([]*node, level+1),
    }
}

type SkipList struct {
    head    *node
    level   int // current maximum level in use
    length  int
}

func New() *SkipList {
    return &SkipList{
        head:  newNode(0, nil, maxLevel),
        level: 0,
    }
}

func randomLevel() int {
    level := 0
    for level < maxLevel && rand.Float64() < probability {
        level++
    }
    return level
}

// Search returns the value for key if it exists.
func (sl *SkipList) Search(key int) (any, bool) {
    curr := sl.head
    for i := sl.level; i >= 0; i-- {
        for curr.forward[i] != nil && curr.forward[i].key < key {
            curr = curr.forward[i]
        }
    }
    curr = curr.forward[0]
    if curr != nil && curr.key == key {
        return curr.value, true
    }
    return nil, false
}

// Insert adds or updates a key-value pair.
func (sl *SkipList) Insert(key int, value any) {
    // update[i] = the rightmost node at level i that is to the left of insertion point
    update := make([]*node, maxLevel+1)
    curr := sl.head
    for i := sl.level; i >= 0; i-- {
        for curr.forward[i] != nil && curr.forward[i].key < key {
            curr = curr.forward[i]
        }
        update[i] = curr
    }

    curr = curr.forward[0]
    if curr != nil && curr.key == key {
        curr.value = value // update existing
        return
    }

    lvl := randomLevel()
    if lvl > sl.level {
        for i := sl.level + 1; i <= lvl; i++ {
            update[i] = sl.head
        }
        sl.level = lvl
    }

    n := newNode(key, value, lvl)
    for i := 0; i <= lvl; i++ {
        n.forward[i] = update[i].forward[i]
        update[i].forward[i] = n
    }
    sl.length++
}

// Delete removes a key. Returns true if the key existed.
func (sl *SkipList) Delete(key int) bool {
    update := make([]*node, maxLevel+1)
    curr := sl.head
    for i := sl.level; i >= 0; i-- {
        for curr.forward[i] != nil && curr.forward[i].key < key {
            curr = curr.forward[i]
        }
        update[i] = curr
    }

    target := curr.forward[0]
    if target == nil || target.key != key {
        return false
    }

    for i := 0; i <= sl.level; i++ {
        if update[i].forward[i] != target { break }
        update[i].forward[i] = target.forward[i]
    }

    // Reduce level if top levels are now empty
    for sl.level > 0 && sl.head.forward[sl.level] == nil {
        sl.level--
    }
    sl.length--
    return true
}

// Range calls fn for all keys in [lo, hi].
func (sl *SkipList) Range(lo, hi int, fn func(key int, value any) bool) {
    curr := sl.head
    for i := sl.level; i >= 0; i-- {
        for curr.forward[i] != nil && curr.forward[i].key < lo {
            curr = curr.forward[i]
        }
    }
    curr = curr.forward[0]
    for curr != nil && curr.key <= hi {
        if !fn(curr.key, curr.value) { return }
        curr = curr.forward[0]
    }
}

func (sl *SkipList) Len() int { return sl.length }
```

**Usage:**
```go
sl := skiplist.New()
sl.Insert(30, "thirty")
sl.Insert(10, "ten")
sl.Insert(50, "fifty")
sl.Insert(20, "twenty")

v, ok := sl.Search(20)
fmt.Println(v, ok)  // twenty true

sl.Range(10, 30, func(k int, v any) bool {
    fmt.Println(k, v)  // 10 ten, 20 twenty, 30 thirty
    return true
})
```

---

## 2. Concurrent Skip List

The skip list's layered structure makes it amenable to lock-free concurrent access via CAS (Compare-And-Swap). A simpler approach uses fine-grained locking per node.

```go
package cskip

import (
    "math/rand"
    "sync"
    "sync/atomic"
    "unsafe"
)

const maxLevel = 16

type node struct {
    key     int
    value   atomic.Value   // safe for concurrent writes
    mu      sync.Mutex
    marked  atomic.Bool    // true if logically deleted
    fullyLinked atomic.Bool
    level   int
    forward [maxLevel]*node
}

type ConcurrentSkipList struct {
    head   *node
    level  int32 // current max level
}

func New() *ConcurrentSkipList {
    head := &node{key: -1<<62, level: maxLevel - 1}
    tail := &node{key: 1<<62, level: maxLevel - 1}
    for i := 0; i < maxLevel; i++ {
        head.forward[i] = tail
    }
    return &ConcurrentSkipList{head: head, level: 0}
}

// find locates where key would go at each level.
// Returns the predecessor and successor arrays.
func (sl *ConcurrentSkipList) find(key int, preds, succs []*node) int {
    foundLevel := -1
    curr := sl.head
    maxL := int(atomic.LoadInt32(&sl.level))
    for i := maxL; i >= 0; i-- {
        next := (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&curr.forward[i]))))
        for next != nil && next.key < key {
            curr = next
            next = (*node)(atomic.LoadPointer((*unsafe.Pointer)(unsafe.Pointer(&curr.forward[i]))))
        }
        if next != nil && next.key == key {
            foundLevel = i
        }
        preds[i] = curr
        succs[i] = next
    }
    return foundLevel
}

func (sl *ConcurrentSkipList) Contains(key int) bool {
    var preds, succs [maxLevel]*node
    foundLevel := sl.find(key, preds[:], succs[:])
    return foundLevel != -1 && succs[foundLevel].fullyLinked.Load() && !succs[foundLevel].marked.Load()
}

func (sl *ConcurrentSkipList) Add(key int, value any) bool {
    topLevel := randomLevel()
    var preds, succs [maxLevel]*node

    for {
        foundLevel := sl.find(key, preds[:], succs[:])
        if foundLevel != -1 {
            n := succs[foundLevel]
            if !n.marked.Load() {
                for !n.fullyLinked.Load() { /* spin */ }
                return false // already exists
            }
            continue
        }

        // Lock predecessors
        highest := topLevel
        var locked []*node
        valid := true

        for i := 0; i <= highest; i++ {
            pred := preds[i]
            if i == 0 || pred != preds[i-1] {
                pred.mu.Lock()
                locked = append(locked, pred)
            }
            succ := succs[i]
            if pred.marked.Load() || succ.marked.Load() || pred.forward[i] != succ {
                valid = false
                break
            }
        }

        if !valid {
            for _, l := range locked { l.mu.Unlock() }
            continue
        }

        n := &node{key: key, level: topLevel}
        n.value.Store(value)
        for i := 0; i <= topLevel; i++ {
            n.forward[i] = succs[i]
        }

        for i := 0; i <= topLevel; i++ {
            atomic.StorePointer(
                (*unsafe.Pointer)(unsafe.Pointer(&preds[i].forward[i])),
                unsafe.Pointer(n),
            )
        }

        n.fullyLinked.Store(true)
        for _, l := range locked { l.mu.Unlock() }

        // Update level
        for {
            curLevel := atomic.LoadInt32(&sl.level)
            if int32(topLevel) <= curLevel || atomic.CompareAndSwapInt32(&sl.level, curLevel, int32(topLevel)) {
                break
            }
        }
        return true
    }
}

func randomLevel() int {
    level := 0
    for level < maxLevel-1 && rand.Float64() < 0.5 {
        level++
    }
    return level
}
```

---

## 3. sync.Map — When to Use It

Go's `sync.Map` is a concurrent map optimized for specific access patterns.

```go
var m sync.Map

// Store
m.Store("key", 42)

// Load
v, ok := m.Load("key")
if ok {
    fmt.Println(v.(int)) // 42
}

// LoadOrStore — atomic: load if exists, else store
actual, loaded := m.LoadOrStore("key", 99)
fmt.Println(actual, loaded) // 42 true (was already there)

// Delete
m.Delete("key")

// Range — iterate (order not guaranteed)
m.Range(func(k, v any) bool {
    fmt.Println(k, v)
    return true // return false to stop
})
```

**`sync.Map` is optimized for:**
1. Keys written once, then read many times (read-heavy)
2. Multiple goroutines operating on disjoint sets of keys

**`sync.Map` is NOT good for:**
- Write-heavy workloads
- Uniform read/write patterns

For write-heavy or general-purpose concurrent maps, use a sharded `map[K]V` with per-shard mutexes (shown in Ch 42).

```go
// When sync.Map is fast: many goroutines reading, few writing
var cache sync.Map

func getConfig(key string) string {
    if v, ok := cache.Load(key); ok {
        return v.(string)  // fast path: no lock
    }
    val := fetchFromDB(key)
    cache.Store(key, val)
    return val
}
```

---

## 4. Concurrent Patterns

### Read-mostly: RWMutex

```go
type ReadHeavyMap struct {
    mu sync.RWMutex
    m  map[string]int
}

func (m *ReadHeavyMap) Get(key string) (int, bool) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    v, ok := m.m[key]
    return v, ok
}

func (m *ReadHeavyMap) Set(key string, val int) {
    m.mu.Lock()
    defer m.mu.Unlock()
    m.m[key] = val
}
```

### Write-once-read-many: `sync.Once` + atomic

```go
type Config struct {
    once   sync.Once
    values atomic.Pointer[map[string]string]
}

func (c *Config) Load() map[string]string {
    p := c.values.Load()
    if p != nil { return *p }
    c.once.Do(func() {
        m := loadFromFile()
        c.values.Store(&m)
    })
    return *c.values.Load()
}
```

### Copy-on-write

When reads greatly outnumber writes, keep an immutable snapshot. Writes create a new copy.

```go
type COWMap[K comparable, V any] struct {
    mu      sync.Mutex
    current atomic.Pointer[map[K]V]
}

func (m *COWMap[K, V]) Get(key K) (V, bool) {
    p := m.current.Load()
    if p == nil { var z V; return z, false }
    v, ok := (*p)[key]
    return v, ok
}

func (m *COWMap[K, V]) Set(key K, val V) {
    m.mu.Lock()
    defer m.mu.Unlock()

    // Copy the current map
    old := m.current.Load()
    var newMap map[K]V
    if old != nil {
        newMap = make(map[K]V, len(*old)+1)
        for k, v := range *old { newMap[k] = v }
    } else {
        newMap = make(map[K]V)
    }
    newMap[key] = val
    m.current.Store(&newMap)
}
```

---

## Summary

| Structure | Average ops | Use case |
|-----------|-------------|----------|
| Skip List | O(log n) | Sorted map, range queries |
| Concurrent Skip List | O(log n) | Concurrent sorted data |
| `sync.Map` | O(1) amortized | Read-heavy, disjoint keys |
| Sharded map | O(1) | General concurrent map |
| COW map | O(1) read, O(n) write | Very read-heavy |

- Skip list achieves balanced BST performance without rotations — uses random coin flips to maintain probabilistic balance
- Concurrent skip list uses fine-grained per-node locking — inserts at different positions can proceed in parallel
- `sync.Map` uses two internal maps: a "read" snapshot (lock-free) and a "dirty" map (with mutex). Dirty map is promoted to read when enough misses accumulate.
- For most concurrent map needs: sharded `sync.RWMutex` is fastest; `sync.Map` is simplest for read-heavy

## Exercises

### Easy
1. Add a `Min()` and `Max()` method to the skip list. These should be O(1) by maintaining head and tail sentinel nodes.
2. Implement `Rank(key int) int` for the skip list — how many keys are less than or equal to key. This requires augmenting each node with a span count per level (how many level-0 nodes does this pointer skip).
3. Trace through the concurrent skip list for `Add(10)`, `Add(20)`, `Add(15)`, `Contains(15)` and show what happens at each level during the `find` phase.

### Medium
4. Implement `DeleteRange(lo, hi int) int` for the skip list that removes all keys in [lo, hi] and returns the count of removed elements. This should be O(k + log n) where k is the number of deleted elements.
5. Build a **sorted set** (like Redis ZADD/ZRANGE) on top of the skip list: each element has a score, the skip list is ordered by score, and you can query `RangeByScore(min, max float64) []Element` efficiently.
6. Benchmark skip list vs `sort.Search` on a sorted slice for 1M elements: measure insert, delete, and range query at p10/p50/p99 latency. When does the skip list's log(n) pay off versus the slice's cache-friendly layout?

### Hard
7. Implement a **lock-free stack** using CAS (`sync/atomic`). The stack holds a pointer to the top node. Push: create new node, set `new.next = top`, CAS `top` from old → new. Pop: load `top`, CAS `top` from `top` → `top.next`. Handle the ABA problem using a version counter tagged to the pointer.
8. Build a **concurrent pipeline**: producer goroutines generate items, a pool of worker goroutines processes them, consumer goroutines collect results. Use channels and sync primitives to coordinate. Add back-pressure: workers should slow down producers if the queue depth exceeds a threshold.
