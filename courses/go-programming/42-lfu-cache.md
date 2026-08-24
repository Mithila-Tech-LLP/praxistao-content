# Chapter 42: LRU and LFU Cache

Caches are bounded by memory. When a cache is full and a new item arrives, you must evict something. LRU (Least Recently Used) evicts the item you haven't touched in the longest time. LFU (Least Frequently Used) evicts the item accessed least often overall.

## Table of Contents

1. [LRU Cache](#1-lru-cache)
2. [LFU Cache](#2-lfu-cache)
3. [LRU vs LFU](#3-lru-vs-lfu)
4. [Concurrent Cache](#4-concurrent-cache)
5. [Summary](#summary)
6. [Exercises](#exercises)

---

## 1. LRU Cache

LRU is implemented with a hash map (for O(1) lookup) + a doubly-linked list (for O(1) move-to-front and remove-tail).

```go
package lru

type entry struct {
    key   string
    value any
    prev  *entry
    next  *entry
}

type Cache struct {
    cap   int
    items map[string]*entry
    head  *entry // most recently used
    tail  *entry // least recently used
}

func New(cap int) *Cache {
    if cap <= 0 { panic("cap must be positive") }
    head := &entry{}
    tail := &entry{}
    head.next = tail
    tail.prev = head
    return &Cache{
        cap:   cap,
        items: make(map[string]*entry, cap),
        head:  head,
        tail:  tail,
    }
}

func (c *Cache) remove(e *entry) {
    e.prev.next = e.next
    e.next.prev = e.prev
}

func (c *Cache) pushFront(e *entry) {
    e.next = c.head.next
    e.prev = c.head
    c.head.next.prev = e
    c.head.next = e
}

func (c *Cache) Get(key string) (any, bool) {
    e, ok := c.items[key]
    if !ok { return nil, false }
    // Move to front (most recently used)
    c.remove(e)
    c.pushFront(e)
    return e.value, true
}

func (c *Cache) Put(key string, value any) {
    if e, ok := c.items[key]; ok {
        e.value = value
        c.remove(e)
        c.pushFront(e)
        return
    }
    e := &entry{key: key, value: value}
    c.items[key] = e
    c.pushFront(e)

    if len(c.items) > c.cap {
        // Evict least recently used (tail.prev)
        lru := c.tail.prev
        c.remove(lru)
        delete(c.items, lru.key)
    }
}

func (c *Cache) Len() int { return len(c.items) }
```

**Usage:**
```go
cache := lru.New(3)
cache.Put("a", 1)
cache.Put("b", 2)
cache.Put("c", 3)

cache.Get("a")     // access "a" → it moves to front
cache.Put("d", 4)  // evicts "b" (LRU)

_, ok := cache.Get("b")
fmt.Println(ok)  // false — evicted
```

---

## 2. LFU Cache

LFU is trickier. We need to:
1. Track the access frequency of every key
2. Evict the key with the minimum frequency
3. When there's a tie (same frequency), evict the least recently used among them

**Data structures used:**
- `keyMap`: key → {value, frequency}
- `freqMap`: frequency → ordered list of keys with that frequency
- `minFreq`: track the current minimum frequency

When a key is accessed, we move it from `freqMap[freq]` to `freqMap[freq+1]`.

```go
package lfu

import "container/list"

type entry struct {
    key   string
    value any
    freq  int
}

type Cache struct {
    cap     int
    minFreq int
    keyMap  map[string]*list.Element // key → *list.Element containing *entry
    freqMap map[int]*list.List        // freq → doubly-linked list of *entry
}

func New(cap int) *Cache {
    if cap <= 0 { panic("cap must be positive") }
    return &Cache{
        cap:     cap,
        keyMap:  make(map[string]*list.Element),
        freqMap: make(map[int]*list.List),
    }
}

func (c *Cache) incrFreq(elem *list.Element) {
    e := elem.Value.(*entry)
    oldFreq := e.freq

    // Remove from current frequency list
    c.freqMap[oldFreq].Remove(elem)
    if c.freqMap[oldFreq].Len() == 0 {
        delete(c.freqMap, oldFreq)
        if c.minFreq == oldFreq {
            c.minFreq++
        }
    }

    // Add to new frequency list (at the back = most recently used for this freq)
    e.freq++
    if c.freqMap[e.freq] == nil {
        c.freqMap[e.freq] = list.New()
    }
    newElem := c.freqMap[e.freq].PushBack(e)
    c.keyMap[e.key] = newElem
}

func (c *Cache) Get(key string) (any, bool) {
    elem, ok := c.keyMap[key]
    if !ok { return nil, false }
    c.incrFreq(elem)
    return elem.Value.(*entry).value, true
}

func (c *Cache) Put(key string, value any) {
    if c.cap == 0 { return }

    // Update existing key
    if elem, ok := c.keyMap[key]; ok {
        elem.Value.(*entry).value = value
        c.incrFreq(elem)
        return
    }

    // Evict if at capacity
    if len(c.keyMap) == c.cap {
        // Evict the front of freqMap[minFreq] (least recently used at min freq)
        minList := c.freqMap[c.minFreq]
        evicted := minList.Front()
        if evicted != nil {
            minList.Remove(evicted)
            if minList.Len() == 0 {
                delete(c.freqMap, c.minFreq)
            }
            delete(c.keyMap, evicted.Value.(*entry).key)
        }
    }

    // Insert new key with freq=1
    e := &entry{key: key, value: value, freq: 1}
    if c.freqMap[1] == nil {
        c.freqMap[1] = list.New()
    }
    elem := c.freqMap[1].PushBack(e)
    c.keyMap[key] = elem
    c.minFreq = 1
}

func (c *Cache) Len() int { return len(c.keyMap) }
```

### LFU trace example

```go
cache := lfu.New(3)
cache.Put("a", 1)   // freq: a=1
cache.Put("b", 2)   // freq: a=1 b=1
cache.Put("c", 3)   // freq: a=1 b=1 c=1
cache.Get("a")      // freq: a=2 b=1 c=1
cache.Get("a")      // freq: a=3 b=1 c=1
cache.Get("b")      // freq: a=3 b=2 c=1

// Put "d" — evicts "c" (freq=1, LRU among freq-1 items)
cache.Put("d", 4)

_, ok := cache.Get("c")
fmt.Println(ok)  // false — evicted

// Frequencies now: a=3 b=2 d=1
cache.Put("e", 5)  // evicts "d" (freq=1)
```

---

## 3. LRU vs LFU

| Property | LRU | LFU |
|----------|-----|-----|
| Eviction policy | Least recently accessed | Least frequently accessed |
| Time complexity | O(1) get/put | O(1) get/put |
| Space complexity | O(n) | O(n) |
| Implementation complexity | Simple | Complex |
| Scan resistance | Poor | Good |
| Good for | Temporal locality (recent = likely needed) | Frequency locality (popular items) |
| Bad for | One-time large scans that flush hot items | Newly inserted items (freq=1 always at risk) |

**LRU scan problem:**
```
Cache size 3, hot items: A (freq=100), B (freq=100), C (freq=100)
Scan: D, E, F, G, H (each accessed once)
→ After scan, A B C all evicted despite being hot
```

**LFU cold-start problem:**
```
Cache size 3, items: A B C each accessed 100 times
New hot item D requested 1 time → D has freq=1
→ If the cache is full, D is immediately the eviction candidate
→ LFU never gives new items a chance
```

**TinyLFU** (used in Caffeine, Go's ristretto cache) solves the cold-start problem by using a Count-Min Sketch as an admission filter.

---

## 4. Concurrent Cache

Both LRU and LFU need synchronization for concurrent access:

```go
package cache

import "sync"

// ConcurrentLRU wraps an LRU cache with a mutex.
type ConcurrentLRU struct {
    mu sync.Mutex
    c  *lru.Cache
}

func NewConcurrentLRU(cap int) *ConcurrentLRU {
    return &ConcurrentLRU{c: lru.New(cap)}
}

func (cl *ConcurrentLRU) Get(key string) (any, bool) {
    cl.mu.Lock()
    defer cl.mu.Unlock()
    return cl.c.Get(key)
}

func (cl *ConcurrentLRU) Put(key string, value any) {
    cl.mu.Lock()
    defer cl.mu.Unlock()
    cl.c.Put(key, value)
}
```

For high-concurrency workloads, consider:
1. **Sharded caches**: split by key hash into N independent LRU caches, each with its own lock
2. **`github.com/dgraph-io/ristretto`**: production-grade concurrent cache using TinyLFU
3. **`github.com/allegro/bigcache`**: low-GC cache using pre-allocated byte pools

```go
// Sharded cache for high concurrency
type ShardedCache struct {
    shards [256]*ConcurrentLRU
}

func NewSharded(totalCap int) *ShardedCache {
    sc := &ShardedCache{}
    perShard := totalCap / 256
    if perShard < 1 { perShard = 1 }
    for i := range sc.shards {
        sc.shards[i] = NewConcurrentLRU(perShard)
    }
    return sc
}

func (sc *ShardedCache) shard(key string) *ConcurrentLRU {
    // FNV hash
    h := uint32(2166136261)
    for i := 0; i < len(key); i++ {
        h ^= uint32(key[i])
        h *= 16777619
    }
    return sc.shards[h%256]
}

func (sc *ShardedCache) Get(key string) (any, bool) { return sc.shard(key).Get(key) }
func (sc *ShardedCache) Put(key string, v any)      { sc.shard(key).Put(key, v) }
```

---

## Summary

- **LRU**: doubly-linked list (order by recency) + hash map (O(1) lookup). Evicts the node at the tail.
- **LFU**: per-frequency linked lists + hash maps. Tracks `minFreq` to find eviction candidate in O(1).
- Both achieve O(1) get and put.
- LRU is simpler and better for workloads with temporal locality.
- LFU is better for workloads where popular items should stay cached longer.
- For production: use `ristretto` (TinyLFU) or `bigcache` rather than rolling your own for concurrent use.

## Exercises

### Easy
1. Add a `Keys() []string` method to the LRU cache that returns keys from most-recently-used to least-recently-used. Verify the order with a test.
2. Implement `Delete(key string) bool` for the LRU cache. Return true if the key existed. Update the doubly-linked list and map correctly.
3. Trace through the LFU cache for this sequence and draw the state after each operation: `Put("x",1), Put("y",2), Put("z",3), Get("x"), Get("x"), Get("y"), Put("w",4)`.

### Medium
4. Implement a **TTL cache**: wrap the LRU cache so each entry has an expiry time. Items expire after a given duration. Lazy expiry: check on Get; periodic cleanup: a goroutine removes expired items every minute.
5. Implement a **write-through cache**: when `Put` is called, write to both the cache and a backing store (simulate with a `map[string]string`). On `Get`, check cache first; on miss, load from the backing store and insert into cache. Add a `LoadAll()` method that warms the cache from the backing store.
6. Benchmark sharded (N=256 shards) vs un-sharded LRU under high concurrency: 8 goroutines, 50% reads / 50% writes, 10M operations total. What is the throughput improvement from sharding?

### Hard
7. Implement **TinyLFU**: a Count-Min Sketch admission filter + a window LRU (for new items) + a main LFU cache. An item moves from the window LRU to the main LFU only if its estimated frequency (from the sketch) exceeds the frequency of the candidate to be evicted from the main cache. This gives new items a chance without polluting the main cache with one-off requests.
8. Compare LRU, LFU, and TinyLFU hit rates using a **Zipfian workload** (realistic for web traffic): generate 1M accesses where item k is requested with probability proportional to 1/k. Measure hit rate for each policy with cache sizes 1%, 5%, 10% of the working set. Plot the results.
