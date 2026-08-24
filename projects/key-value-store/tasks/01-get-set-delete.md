# Task 01: Get, Set, Delete

## What you will build

The foundation of every key-value store: storing a string value under a string key, retrieving it, and removing it. You will also make it safe for concurrent access using `sync.RWMutex`.

## Concepts

### Why a mutex?

A Go map is not goroutine-safe. If two goroutines write at the same time, the runtime panics with "concurrent map writes". A `sync.RWMutex` allows unlimited concurrent readers (`RLock`) but only one writer at a time (`Lock`). This is the classic *multiple-reader, single-writer* pattern.

```
Read:   RLock()  → read map → RUnlock()
Write:  Lock()   → write map → Unlock()
```

### Value type

Store values as `string`. That matches Redis's external wire format. Later tasks will layer richer types (lists, sets, hashes) on top.

## Interface to implement

```go
type Store struct { /* unexported fields */ }

func NewStore() *Store

// Set stores value under key, overwriting any previous value.
func (s *Store) Set(key, value string)

// Get returns the value and true if key exists, or ("", false) if not.
func (s *Store) Get(key string) (string, bool)

// Delete removes key. Returns true if the key existed.
func (s *Store) Delete(key string) bool

// Exists returns true if key is present in the store.
func (s *Store) Exists(key string) bool
```

## Hints

- Keep `Store` simple: one map and one mutex are enough.
- `Delete` on a missing key should not panic — just return `false`.
- For `Exists`, prefer `RLock` (read-only operation).

## Run the tests

```bash
cd starter/task-01-get-set-delete
go test ./...
```

All tests must pass before moving to Task 02.
