# Task 02: TTL Expiry

## What you will build

Keys that automatically disappear after a duration. This is how Redis implements cache expiry, session timeouts, and rate-limit windows.

## Concepts

### Two expiry strategies

**Lazy expiry** — check the timestamp on every `Get`. If the deadline has passed, treat the key as gone and delete it. Simple, but stale keys stay in memory until someone asks for them.

**Active expiry** — a background goroutine wakes up periodically, scans for expired keys, and deletes them. Keeps memory clean even for keys nobody reads.

You can implement either, or both. A production store does both: lazy for correctness, active for memory hygiene.

### Storing the deadline

Add a second map `expiry map[string]time.Time`. When `SetWithTTL` is called, store `time.Now().Add(ttl)`. When a key is accessed, check `time.Now().Before(deadline)`.

## Interface to add

```go
// SetWithTTL stores value under key and expires it after ttl.
// Calling Set on the same key removes any existing TTL.
func (s *Store) SetWithTTL(key, value string, ttl time.Duration)

// TTL returns the remaining lifetime of a key.
// Returns (remaining, true) if the key exists and has a TTL.
// Returns (0, false) if the key does not exist or has no TTL.
func (s *Store) TTL(key string) (time.Duration, bool)
```

**Modify `Get`** — an expired key must return `("", false)`, identical to a missing key.

**Modify `Set`** — setting a key without a TTL must clear any previous expiry on that key.

## Hints

- The background goroutine (if you use one) should stop when the store is garbage-collected. Use a `done` channel and a `Close()` method, or tie lifetime to a `context.Context`.
- Hold the write lock for the minimum time during background cleanup. You can collect expired keys under `RLock`, then delete them under `Lock`.
- `time.Duration` arithmetic: `deadline.Sub(time.Now())` gives remaining time.

## Run the tests

```bash
cd starter/task-02-ttl-expiry
go test ./...
```
