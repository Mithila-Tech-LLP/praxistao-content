# Task 06: Sorted Set

## What you will build

A sorted set combines the uniqueness of a set with the ordering of a list. Each member has a `float64` score, and members are always iterated in ascending score order. Redis uses sorted sets for leaderboards, priority queues, and time-series indexes.

## Concepts

### Storage design

You need two structures per key:
1. `map[string]float64` — member → score, for O(1) score lookup.
2. `[]ZMember` — a slice kept in sorted order, for range queries.

When `ZAdd` is called, update the map and re-sort (or insert at the correct position in) the slice.

```go
type ZMember struct {
    Member string
    Score  float64
}
```

### Range and rank

`ZRange(key, 0, -1)` returns all members in ascending score order (same negative-index convention as `LRange`). `ZRank` returns the 0-based position of a member in that ordering.

## Interface to implement

```go
type ZMember struct {
    Member string
    Score  float64
}

// ZAdd adds or updates member with the given score.
// If member already exists, its score is updated and the sorted order is maintained.
func (s *Store) ZAdd(key string, score float64, member string)

// ZScore returns (score, true) if member exists, or (0, false) if not.
func (s *Store) ZScore(key, member string) (float64, bool)

// ZRange returns members from start to stop by rank (ascending). Supports negative indices.
func (s *Store) ZRange(key string, start, stop int) []string

// ZRangeWithScores returns ZMember structs instead of just strings.
func (s *Store) ZRangeWithScores(key string, start, stop int) []ZMember

// ZRank returns the 0-based rank of member (ascending order). Returns (0, false) if not found.
func (s *Store) ZRank(key, member string) (int, bool)
```

## Hints

- Use `sort.Slice` after any modification that changes order. For small sets this is fast enough; a production store would use a skip list.
- When `ZAdd` updates an existing member's score, remove the old entry from the slice before re-inserting.
- Ties in score should break alphabetically by member name — standard Redis behaviour. Add a secondary sort on `Member` to your comparator.

## Run the tests

```bash
cd starter/task-06-sorted-set
go test ./...
```
