# Task 08: Transactions

## What you will build

Atomic batches of commands. A transaction queues up multiple operations and executes them all-at-once under a single lock. No other goroutine can interleave reads or writes while the transaction is executing. This is the foundation of Redis `MULTI`/`EXEC`.

## Concepts

### Command queue

A `Tx` is a list of closures (or command structs). Each method call on `Tx` records what to do without doing it.

```go
tx := store.Begin()
tx.Set("a", "1")
tx.Delete("b")
tx.Set("c", "3")
tx.Exec() // all three run atomically
```

### Atomicity via locking

`Exec` acquires the store's write lock once, runs every queued command against the store's internal map directly (bypassing the per-method locking), then releases the lock. This ensures no other operation observes the store in a half-applied state.

### Discard

`Discard` simply clears the queue. Nothing is applied to the store.

## Interface to implement

```go
type Tx struct { /* unexported */ }

// Begin creates a new transaction bound to this store.
func (s *Store) Begin() *Tx

// Set queues a Set command in the transaction.
func (t *Tx) Set(key, value string)

// Delete queues a Delete command in the transaction.
func (t *Tx) Delete(key string)

// Exec executes all queued commands atomically.
// Returns an error if the transaction has already been executed or discarded.
func (t *Tx) Exec() error

// Discard abandons the transaction without executing any commands.
func (t *Tx) Discard()
```

## Hints

- Store commands as `[]func()` — each closure captures what it needs. `Exec` loops over the slice and calls each function while holding the lock.
- Add a `done bool` field to `Tx` to prevent double-exec.
- The queued closures should call unexported helper methods that assume the lock is already held, avoiding dead-lock from re-entrant locking.

## Run the tests

```bash
cd starter/task-08-transactions
go test ./...
```
