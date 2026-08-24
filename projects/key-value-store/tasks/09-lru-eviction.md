# Task 09: LRU Eviction

## What you will build

A capacity-bounded store that evicts the least recently used key when full. This is how Redis `maxmemory-policy allkeys-lru` works and how most in-memory caches prevent unbounded growth.

## Concepts

### LRU tracking

Every time a key is read (`Get`) or written (`Set`/`SetWithTTL`), it becomes the most recently used. When the store is full and a new key arrives, the least recently used key is evicted to make room.

### Implementation approaches

**Simple (O(n) eviction):** keep a `map[string]time.Time` of last-access timestamps. On eviction, scan all keys and find the minimum. Fine for small stores.

**Classic (O(1) eviction):** doubly-linked list + map. Move accessed keys to the front; evict from the back. This is the textbook LRU cache structure.

Either approach passes the tests. Implement whichever feels right, then think about the trade-off.

## Interface to add

```go
// NewStoreWithCapacity creates a store that evicts LRU keys when maxKeys is exceeded.
// maxKeys <= 0 means no limit (behaves like NewStore).
func NewStoreWithCapacity(maxKeys int) *Store
```

The existing `Set` method must now check capacity and evict if needed. No changes to the method signature.

## Hints

- Eviction happens *before* inserting the new key, so the store never exceeds `maxKeys`.
- A key being overwritten (same key, new value) does not increase the count — update its access time and do not evict.
- Access order: `Get` on a non-existent key does NOT update order (there is nothing to touch).
- For the doubly-linked list approach, `container/list` from the stdlib is a ready-made implementation.

## Run the tests

```bash
cd starter/task-09-lru-eviction
go test ./...
```
