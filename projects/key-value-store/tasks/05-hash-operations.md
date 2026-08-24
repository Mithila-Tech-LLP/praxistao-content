# Task 05: Hash Operations

## What you will build

Hashes are maps-within-a-map: a single key holds a collection of field→value pairs. Redis uses them to store structured records — a user profile, a product listing, or a session object — without serialising to JSON.

## Concepts

### Nested maps

A hash is `map[string]map[string]string` at the storage level. The outer key is the store key; the inner map holds fields.

```
store["user:42"] → {"name": "Ada", "email": "ada@example.com", "role": "admin"}
```

### Partial updates

`HSet` sets a single field without affecting others. `HDel` removes specific fields. `HGetAll` returns a snapshot of all fields at once.

## Interface to implement

```go
// HSet sets field to value in the hash stored at key.
// Creates the hash if it does not exist.
func (s *Store) HSet(key, field, value string)

// HGet returns (value, true) for an existing field, or ("", false) if the hash
// or field does not exist.
func (s *Store) HGet(key, field string) (string, bool)

// HGetAll returns a copy of all field-value pairs in the hash.
// Returns an empty map (not nil) if key does not exist.
func (s *Store) HGetAll(key string) map[string]string

// HDel removes one or more fields. Returns count of actually deleted fields.
func (s *Store) HDel(key string, fields ...string) int

// HExists returns true if field exists in the hash.
func (s *Store) HExists(key, field string) bool
```

## Hints

- `HGetAll` should return a *copy*, not a reference to the internal map. Callers mutating the returned map must not affect the store.
- After `HDel` removes the last field from a hash, you may optionally delete the top-level key. Tests will not rely on this either way.
- `HSet` on a key that holds a string type would be a type conflict in Redis (`WRONGTYPE` error). For simplicity you can ignore this or panic — the tests will not mix types on a single key.

## Run the tests

```bash
cd starter/task-05-hash-operations
go test ./...
```
