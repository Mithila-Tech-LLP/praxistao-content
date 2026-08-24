# Task 06: The Mempool

## What you will build

A mempool ("memory pool") that holds pending transactions before they're mined, and rejects a transaction that tries to spend a UTXO another pending transaction has already claimed — the mechanism that prevents double-spending.

## Concepts

### The waiting room

Transactions don't go straight into a block. They wait in the mempool until a miner picks them up. This is exactly where double-spend prevention has to happen: two transactions might both try to spend the same coin before either one is mined, and only one of them can be allowed to succeed.

```
   tx1: spend UTXO #7 -> pay Bob        }  both reference the SAME
   tx2: spend UTXO #7 -> pay Mallory    }  unspent output #7

   Mempool must accept tx1 (or tx2) -- whichever arrives first --
   and reject the other one outright.
```

### Detecting a conflict

Every transaction input references a specific `(TxID, OutIndex)` pair. If a new transaction's input matches a `(TxID, OutIndex)` pair already claimed by something sitting in the mempool, it's a conflicting spend and must be rejected, even though the transaction itself is otherwise perfectly well-formed.

## Interface to implement

```go
type Mempool struct {
	// unexported fields
}

func NewMempool() *Mempool

// Add validates tx against everything currently pending and, if nothing
// conflicts, adds it. Returns an error if tx tries to spend a
// (TxID, OutIndex) pair some other pending transaction already claims.
func (mp *Mempool) Add(tx *Transaction) error

// Remove takes a transaction out of the mempool (called after it's
// been mined into a block).
func (mp *Mempool) Remove(txID []byte)

// Pending returns every transaction currently waiting, in the order
// they were added.
func (mp *Mempool) Pending() []*Transaction

// Has reports whether a transaction with the given ID is currently
// pending.
func (mp *Mempool) Has(txID []byte) bool
```

## Hints

- Track two maps internally: pending transactions keyed by hex-encoded `TxID`, and claimed `(TxID, OutIndex)` pairs keyed by something like `fmt.Sprintf("%x:%d", txID, outIndex)`.
- `Add` should check *every* input of the incoming transaction against the claimed-outputs map before accepting any of it — reject the whole transaction if even one input conflicts.
- Write a test that adds a transaction, then attempts to add a second, different transaction spending the same `(TxID, OutIndex)`, and confirms the second `Add` returns an error while the first transaction is still present via `Has`.
- Also test that `Remove` followed by a fresh `Add` of a transaction spending the same output now succeeds — removing a transaction should free up whatever it had claimed.

## Run the tests

```bash
cd starter/task-06-the-mempool
go test ./...
```

All tests must pass before moving to Task 07.
