# Task 10: Snapshot

## What you will build

Persistence: save the entire store to an `io.Writer` and restore it from an `io.Reader`. This is the last piece that makes the store useful in production — it can survive a restart.

## Concepts

### Serialisation format

Use JSON. `encoding/json` is in the stdlib and handles `map[string]string` directly.

```json
{
  "keys": {
    "greeting": "hello",
    "counter":  "42"
  }
}
```

Only string keys need to be persisted in this task. Lists, sets, hashes, and sorted sets are bonus if you want to extend.

### io.Writer / io.Reader

Writing to a file, a network socket, or an in-memory buffer all look the same through these interfaces. The test uses `bytes.Buffer`, but your implementation works with any writer/reader.

```go
var buf bytes.Buffer
store.Save(&buf)        // writes JSON into buf
newStore := NewStore()
newStore.Load(&buf)     // restores from buf
```

### Snapshot consistency

Hold the read lock while serialising so the snapshot is consistent — no partial writes from a concurrent `Set`.

## Interface to implement

```go
// Save serialises all string key-value pairs to w as JSON.
// Returns an error if serialisation fails.
func (s *Store) Save(w io.Writer) error

// Load deserialises key-value pairs from r and merges them into the store.
// Existing keys are overwritten. Returns an error if deserialisation fails.
func (s *Store) Load(r io.Reader) error
```

## Hints

- Define a private struct for the JSON envelope: `type snapshot struct { Keys map[string]string }`.
- `json.NewEncoder(w).Encode(v)` and `json.NewDecoder(r).Decode(&v)` are cleaner than `Marshal`/`Unmarshal` for streaming.
- Hold `RLock` during `Save` (reading only). Hold `Lock` during `Load` (writing new keys).

## Run the tests

```bash
cd starter/task-10-snapshot
go test ./...
```

Congratulations — you have built a Redis-inspired key-value store from scratch.
