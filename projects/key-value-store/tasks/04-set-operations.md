# Task 04: Set Operations

## What you will build

Unordered collections of unique strings. Redis sets power features like "users who liked this post", "online members", and tag systems. The defining property: each member appears at most once.

## Concepts

### Implementing a set in Go

Go has no built-in set type. The idiom is `map[string]struct{}`. The empty struct `struct{}{}` takes zero bytes, so only the keys cost memory.

```go
s := map[string]struct{}{}
s["alice"] = struct{}{}
s["alice"] = struct{}{}   // no-op, still one entry
_, ok := s["alice"]       // ok == true
```

### Idempotency

`SAdd` with a member already in the set should be a no-op for that member. The return value is how many *new* members were actually added (not the total passed in).

## Interface to implement

```go
// SAdd adds members to the set. Returns count of newly added members (duplicates ignored).
func (s *Store) SAdd(key string, members ...string) int

// SMembers returns all members of the set in arbitrary order.
func (s *Store) SMembers(key string) []string

// SIsMember returns true if member belongs to the set.
func (s *Store) SIsMember(key, member string) bool

// SRem removes members from the set. Returns count of actually removed members.
func (s *Store) SRem(key string, members ...string) int

// SCard returns the number of members in the set (cardinality).
func (s *Store) SCard(key string) int
```

## Hints

- `SMembers` on a non-existent key returns an empty (not nil) slice.
- `SCard` on a non-existent key returns 0.
- Iteration order over a Go map is not guaranteed — tests must not assume order for `SMembers`. Sort the slice before comparing if the test needs a stable result.

## Run the tests

```bash
cd starter/task-04-set-operations
go test ./...
```
