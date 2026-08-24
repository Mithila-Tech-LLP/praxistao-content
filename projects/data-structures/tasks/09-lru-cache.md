---
title: LRU Cache
step: 9
difficulty: hard
estimated: 60 min
---

## What You Are Building

An LRU (Least Recently Used) Cache evicts the item that was used least recently when the cache is full. It provides O(1) `Get` and `Put` by combining two data structures: a hash map for O(1) lookup and a doubly-linked list to track usage order.

This is one of the most common system design and interview data structure questions.

## Key Concepts

**Why a doubly-linked list?** — We need to move a node to the front of the list in O(1). With a singly-linked list, removing a node requires finding its predecessor (O(n)). With a doubly-linked list, each node knows its previous and next neighbour, so removal is O(1).

**Why a hash map?** — To jump directly to a node in the linked list by key in O(1). Without it, finding the node to move to the front would be O(n).

**The combined structure:**
```
Map:  { 1 → nodeA, 2 → nodeB, 3 → nodeC }
List: head ↔ [3] ↔ [1] ↔ [2] ↔ tail
             MRU                    LRU
```

**Get**: find the node via the map, move it to the front of the list (it just became most-recently-used), return its value.

**Put**: if key exists, update value and move to front. If key is new and cache is full, remove the node at the tail (LRU), delete its key from the map, then add the new node at the front.

**Sentinel nodes** — Use a dummy `head` and `tail` node. This eliminates nil checks: the real nodes always sit between two non-nil sentinels.

```
head ↔ [real nodes...] ↔ tail
```

## Struct Signatures

```go
type lruNode struct {
    key, val   int
    prev, next *lruNode
}

type LRUCache struct {
    capacity int
    cache    map[int]*lruNode
    head     *lruNode // dummy head (MRU side)
    tail     *lruNode // dummy tail (LRU side)
}

func NewLRUCache(capacity int) *LRUCache
```

## Methods to Implement

| Method | Description |
|--------|-------------|
| `Get(key int) (int, bool)` | Return value and move to front; ok=false if missing |
| `Put(key, value int)` | Insert or update; evict LRU if over capacity |

## Edge Cases to Handle

- `Get` a key that doesn't exist: return `0, false`
- `Put` to update an existing key: update value, move to front, no eviction
- `Put` when at capacity: evict the least-recently-used item before inserting
- Cache with capacity 1: every new `Put` evicts the previous item

## Example

```go
cache := NewLRUCache(3)
cache.Put(1, 10)
cache.Put(2, 20)
cache.Put(3, 30)

val, ok := cache.Get(1)
fmt.Println(val, ok) // 10 true  (1 is now MRU)

cache.Put(4, 40) // evicts 2 (LRU)

_, ok = cache.Get(2)
fmt.Println(ok) // false (was evicted)

val, _ = cache.Get(3)
fmt.Println(val) // 30
```

## Hints

- Write two private helpers: `removeNode(node *lruNode)` and `insertFront(node *lruNode)`. Every `Get` and `Put` operation is just combinations of these two.
- `removeNode`: `node.prev.next = node.next; node.next.prev = node.prev`
- `insertFront` (after head sentinel): `node.next = h.head.next; node.prev = h.head; h.head.next.prev = node; h.head.next = node`
- In `NewLRUCache`, wire up sentinels: `head.next = tail; tail.prev = head`.
- The LRU node is always `h.tail.prev` (the real node just before the dummy tail).
