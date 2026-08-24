# Task 03: List Operations

## What you will build

Doubly-ended lists keyed by name. Redis lists are used for task queues (`LPUSH` / `BRPOP`), activity feeds, and message logs. You will implement push, pop, and range on both ends.

## Concepts

### Type coexistence

The store now holds two distinct value types: strings (from Task 01/02) and lists. Keep them separate — a key that holds a list is a different thing from a key that holds the string `"[a,b,c]"`.

One clean approach: a second map `lists map[string][]string`. An attempt to `LPush` a key that already holds a string value should return an error (or panic — your choice; tests will not mix types).

### LRange indexing

Redis `LRANGE key 0 -1` returns all elements. Negative indices count from the end: `-1` is the last element, `-2` is second-to-last. Implement this convention.

```
list: ["a", "b", "c", "d"]
LRange(0, -1)  → ["a","b","c","d"]
LRange(0, 1)   → ["a","b"]
LRange(-2, -1) → ["c","d"]
```

Out-of-bounds indices should be clamped, not panicked.

## Interface to implement

```go
// LPush prepends vals to the list. Returns new list length.
func (s *Store) LPush(key string, vals ...string) int

// RPush appends vals to the list. Returns new list length.
func (s *Store) RPush(key string, vals ...string) int

// LPop removes and returns the first element. Returns ("", false) if empty.
func (s *Store) LPop(key string) (string, bool)

// RPop removes and returns the last element. Returns ("", false) if empty.
func (s *Store) RPop(key string) (string, bool)

// LRange returns elements from start to stop (inclusive), supporting negative indices.
func (s *Store) LRange(key string, start, stop int) []string

// LLen returns the length of the list, or 0 if key does not exist.
func (s *Store) LLen(key string) int
```

## Hints

- `LPush("q", "a", "b", "c")` should result in `["c","b","a"]` — each element is prepended in order, so the last argument ends up at the front. This matches Redis behaviour.
- A slice satisfies the list requirement. `append` for RPush, prepend (`append([]string{v}, existing...)`) for LPush.
- Normalize negative indices before slicing: `if i < 0 { i = len(list) + i }`.

## Run the tests

```bash
cd starter/task-03-list-operations
go test ./...
```
