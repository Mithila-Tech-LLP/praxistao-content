# Chapter 56: Fast UTXO Lookups and Indexing

Chapter 55 gave GoChain a real, crash-safe `storage.Store` with a `utxo` bucket ready and waiting. Nothing writes to that bucket yet. This chapter builds the piece that does: `UTXOSet`, a dedicated index of every currently-spendable output, built once from the full chain and then kept up to date incrementally, one block at a time, forever after. By the end, `BalanceOf(address)` stops being "read the entire blockchain and add things up" and becomes "look up a small, targeted set of entries" — and we measure exactly how much that is worth, in milliseconds, on a chain with tens of thousands of blocks.

## Table of Contents

1. [The Problem, Restated in UTXO Terms](#1-the-problem-restated-in-utxo-terms)
2. [Designing the UTXOSet Index](#2-designing-the-utxoset-index)
3. [Key Format and the `utxo` Bucket](#3-key-format-and-the-utxo-bucket)
4. [An Optional Capability: Iterating the Whole Index](#4-an-optional-capability-iterating-the-whole-index)
5. [Building the Index from Scratch: `Reindex`](#5-building-the-index-from-scratch-reindex)
6. [Incremental Updates: Maintaining the Index Every Block](#6-incremental-updates-maintaining-the-index-every-block)
7. [`FindSpendableOutputs` and `BalanceOf`](#7-findspendableoutputs-and-balanceof)
8. [Benchmark: Scan vs. Indexed Lookup](#8-benchmark-scan-vs-indexed-lookup)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. The Problem, Restated in UTXO Terms

Chapter 53 measured `FindTransaction` — a full scan looking for one specific transaction ID — at roughly 400 milliseconds on a 50,000-block chain. `BalanceOf(address)` is a variation on the same scan, and a strictly worse one: instead of stopping the instant it finds a match, it must check *every* output in *every* transaction in *every* block, because an address's balance is the sum of everything it owns, not just the first thing it owns. There is no early exit. A wallet's "refresh balance" button, on a naive implementation, pays the full cost of reading the entire chain's history, every single time it is pressed.

```
NAIVE BalanceOf(address) — Chapter 20/53 style, no index

for each block in the ENTIRE chain (oldest to newest):     <- O(n) blocks
    for each transaction in block:                          <- x ~15 tx/block
        for each output in transaction:                      <- x ~2 outputs/tx
            if output belongs to address AND is unspent:
                total += output.Value

Total work: proportional to the ENTIRE chain's history,
every single time this function is called — even for a
wallet that only ever received three payments, ever.
```

The fix is not a cleverer scan. It is refusing to scan at all: maintain a separate, compact index — the **UTXO set** — that already answers "which outputs are unspent right now," updated the moment it changes, so a balance query never has to touch history that has nothing to do with the answer.

---

## 2. Designing the UTXOSet Index

`UTXOSet` is a thin wrapper around a `storage.Store`. It does not know or care whether that store is a `BoltStore` or something else — the interface-first design from Chapter 55 pays for itself immediately here, since every method below is written purely in terms of `Store`'s methods.

```go
// storage/utxo.go
package storage

import (
	"encoding/hex"
	"fmt"

	"github.com/you/gochain/core"
	"github.com/you/gochain/crypto"
)

// UTXOSet is an index of every currently-unspent transaction output,
// maintained inside a Store's utxo bucket. It answers "how much does this
// address have?" and "which outputs can cover this payment?" directly,
// without ever re-reading the chain's full history.
type UTXOSet struct {
	store Store
}

// NewUTXOSet wraps an existing Store. It does not build or validate the
// index — call Reindex first if the store's utxo bucket might be empty or
// stale (for example, on a completely fresh node).
func NewUTXOSet(store Store) *UTXOSet {
	return &UTXOSet{store: store}
}
```

Think of `UTXOSet` the way a librarian thinks of a card catalog. The books themselves — every block, every transaction, ever mined — are the full history, permanently shelved in the `blocks` bucket and never thrown away. The card catalog is a *separate*, much smaller structure that answers "which books are currently checked out" without anyone needing to walk every aisle. Lose the catalog and you can always rebuild it by walking every aisle once — slow, but recoverable. That rebuild is exactly what `Reindex` does in Section 5; the day-to-day "someone returned a book, someone checked one out" updates are what `Update` does in Section 6.

---

## 3. Key Format and the `utxo` Bucket

Every entry in the `utxo` bucket needs a key that uniquely identifies one specific output. GoChain uses the same convention real UTXO-based chains do: the owning transaction's ID, in hex, followed by the output's position within that transaction's output list.

```
Key format:  "<txID-hex>:<output-index>"

Example transaction with ID 9f2a... and two outputs:
    "9f2a...:0"  ->  TxOutput{ Value: 30, PubKeyHash: <Alice's hash> }   (payment)
    "9f2a...:1"  ->  TxOutput{ Value: 12, PubKeyHash: <Bob's hash> }     (change, back to Bob)

If a later transaction spends output 0 of 9f2a..., the key "9f2a...:0" is
deleted from the bucket entirely. The key "9f2a...:1" remains, because it
has not been spent yet.
```

```go
// storage/utxo.go (continued)

// utxoKey builds the "<txID-hex>:<index>" key format used throughout this
// chapter, so every method constructs keys identically.
func utxoKey(txID []byte, outputIndex int) string {
	return fmt.Sprintf("%s:%d", hex.EncodeToString(txID), outputIndex)
}
```

An output's *presence* in the bucket means "unspent, spendable right now." Its *absence* means either "never existed" or "already spent" — the index deliberately does not distinguish between those two cases, because nothing downstream needs to.

---

## 4. An Optional Capability: Iterating the Whole Index

`FindSpendableOutputs` and `BalanceOf` both need to scan every entry belonging to one address — and `storage.Store`, by design (Chapter 55, Section 3), exposes no bulk "iterate everything" method for the `utxo` bucket, only single-key `Get`/`Put`/`Delete`. That was not an oversight. Keeping `Store` minimal is exactly what makes a future in-memory test fake trivial to write. But `UTXOSet` still needs *some* way to walk the whole index for a bulk operation like a rebuild or a balance scan.

The idiomatic Go answer is an **optional interface** — the same pattern the standard library uses for `http.Flusher` or `io.ReaderFrom`: a small, separate interface that a concrete type *may* implement, checked with a type assertion, so callers that need the extra capability can use it when it is available without forcing every implementation to support it.

```go
// storage/utxo.go (continued)

// utxoIterable is an optional capability a Store may implement to allow
// UTXOSet to walk every entry in the utxo bucket directly, rather than
// through single-key Get calls. BoltStore implements it; a minimal
// in-memory test fake is free not to, at the cost of Reindex and BalanceOf
// falling back to an error rather than working at all — a deliberate
// trade-off, not an accident.
type utxoIterable interface {
	// ForEachUTXO calls fn once per entry currently in the utxo bucket. If
	// fn returns an error, iteration stops and that error is returned.
	ForEachUTXO(fn func(key string, output *core.TxOutput) error) error
}
```

`BoltStore` implements it with a bucket cursor, the same tool Chapter 54's exercises pointed you toward:

```go
// storage/bolt_store.go (continued from Chapter 55)

func (s *BoltStore) ForEachUTXO(fn func(key string, output *core.TxOutput) error) error {
	return s.db.View(func(tx *bolt.Tx) error {
		return tx.Bucket(utxoBucket).ForEach(func(k, v []byte) error {
			var out core.TxOutput
			if err := gob.NewDecoder(bytes.NewReader(v)).Decode(&out); err != nil {
				return fmt.Errorf("decode utxo %s: %w", k, err)
			}
			return fn(string(k), &out)
		})
	})
}
```

`Cursor`-based iteration inside `ForEach` walks the bucket's B-tree in sorted key order — a direct consequence of the B-tree structure Chapter 54 diagrammed, not something `BoltStore` has to sort itself.

---

## 5. Building the Index from Scratch: `Reindex`

`Reindex` throws away whatever the `utxo` bucket currently holds and rebuilds it from the authoritative source of truth: the full chain, block by block. It runs in two passes over `bc`'s existing block iterator (the one `core.Blockchain` has provided since Chapter 18) — first collecting every output that has *ever* been spent by some later input, then walking the chain again and indexing every output that never showed up in that spent set.

```go
// storage/utxo.go (continued)

// Reindex rebuilds the entire UTXO index from bc's full block history.
// Call it once when a node starts up with an empty or untrustworthy utxo
// bucket — for example, after importing a chain exported by another node,
// exactly as Chapter 58's mini project does.
func (u *UTXOSet) Reindex(bc *core.Blockchain) error {
	iterable, ok := u.store.(utxoIterable)
	if !ok {
		return fmt.Errorf("reindex: store does not support ForEachUTXO iteration")
	}

	// Pass 0: clear whatever the bucket currently holds.
	var staleKeys []string
	if err := iterable.ForEachUTXO(func(key string, _ *core.TxOutput) error {
		staleKeys = append(staleKeys, key)
		return nil
	}); err != nil {
		return fmt.Errorf("reindex: list existing entries: %w", err)
	}
	for _, key := range staleKeys {
		if err := u.store.DeleteUTXO(key); err != nil {
			return fmt.Errorf("reindex: clear entry %s: %w", key, err)
		}
	}

	// Pass 1: walk every block once, recording every output any
	// transaction's inputs ever spent, anywhere in the chain's history.
	spent := make(map[string]struct{})
	iter := bc.Iterator()
	for iter.HasNext() {
		block, err := iter.Next()
		if err != nil {
			return fmt.Errorf("reindex pass 1: %w", err)
		}
		for _, tx := range block.Transactions {
			for _, in := range tx.Vin {
				if len(in.Txid) == 0 {
					continue // coinbase transactions have no real inputs
				}
				spent[utxoKey(in.Txid, in.Vout)] = struct{}{}
			}
		}
	}

	// Pass 2: walk every block again. Any output NOT in the spent set from
	// pass 1 is, by definition, still unspent — index it.
	iter = bc.Iterator()
	indexed := 0
	for iter.HasNext() {
		block, err := iter.Next()
		if err != nil {
			return fmt.Errorf("reindex pass 2: %w", err)
		}
		for _, tx := range block.Transactions {
			for outIdx, out := range tx.Vout {
				key := utxoKey(tx.ID, outIdx)
				if _, isSpent := spent[key]; isSpent {
					continue
				}
				if err := u.store.PutUTXO(key, out); err != nil {
					return fmt.Errorf("reindex pass 2: index %s: %w", key, err)
				}
				indexed++
			}
		}
	}

	return nil
}
```

Two passes are necessary, not accidental: a single forward pass cannot know, on the block that *creates* an output, whether some later block will spend it — you only find that out by having already looked ahead. Building the "spent" set first, then filtering against it on a second pass, sidesteps that ordering problem entirely at the cost of walking the chain twice instead of once. Since `Reindex` runs once at startup (or on demand, as in Chapter 58), paying for two passes over the whole chain is a completely different cost than paying for one full scan on *every single balance query*, which is precisely the distinction this chapter exists to draw.

---

## 6. Incremental Updates: Maintaining the Index Every Block

Rebuilding from scratch on every new block would defeat the entire point. `Update` is the method GoChain's mining loop and its network-sync logic both call, exactly once, the moment a new block is accepted — it touches only *that* block's transactions, not the chain's full history.

```go
// storage/utxo.go (continued)

// Update applies exactly one newly-accepted block's transactions to the
// index: every input's referenced output is removed (it is now spent),
// and every output the block's transactions create is added (it is now
// unspent). Call this once per block, right after core.Blockchain accepts
// it — whether that block was mined locally or received from a peer.
func (u *UTXOSet) Update(block *core.Block) error {
	for _, tx := range block.Transactions {
		for _, in := range tx.Vin {
			if len(in.Txid) == 0 {
				continue // coinbase transactions spend nothing
			}
			key := utxoKey(in.Txid, in.Vout)
			if err := u.store.DeleteUTXO(key); err != nil {
				return fmt.Errorf("update: remove spent output %s: %w", key, err)
			}
		}
		for outIdx, out := range tx.Vout {
			key := utxoKey(tx.ID, outIdx)
			if err := u.store.PutUTXO(key, out); err != nil {
				return fmt.Errorf("update: add new output %s: %w", key, err)
			}
		}
	}
	return nil
}
```

The cost of `Update` is proportional only to the number of transactions and outputs in *one* block — a handful of BoltDB writes, not a walk over the entire chain. This is the incremental-maintenance half of the design: `Reindex` pays the full O(n) cost exactly once, and every block after that pays only for itself, forever.

```
Reindex:  O(entire chain history)     — paid ONCE, at startup or recovery
Update:   O(one block's transactions) — paid once PER NEW BLOCK, forever after

Compare to the naive BalenceOf from Section 1, which pays O(entire chain
history) on EVERY SINGLE QUERY, no matter how many times it's called.
```

---

## 7. `FindSpendableOutputs` and `BalanceOf`

With the index actually being maintained, answering "does this address have enough to cover a payment?" and "what is this address's balance?" both become a filtered scan over the `utxo` bucket alone — not the whole chain, just the (typically much smaller) set of outputs that currently exist at all.

```go
// storage/utxo.go (continued)

// FindSpendableOutputs finds enough unspent outputs belonging to pubKeyHash
// to cover amount, stopping as soon as it has enough (the same early-exit
// Chapter 53 pointed out the naive FindTransaction scan gets, and BalanceOf
// below deliberately does not). It returns the total value accumulated and
// a map from transaction ID (hex) to the specific output indices selected —
// exactly the shape core.NewTransaction needs to build a new transaction's
// inputs.
func (u *UTXOSet) FindSpendableOutputs(pubKeyHash []byte, amount int64) (int64, map[string][]int) {
	unspentOutputs := make(map[string][]int)
	var accumulated int64

	iterable, ok := u.store.(utxoIterable)
	if !ok {
		return 0, unspentOutputs
	}

	iterable.ForEachUTXO(func(key string, out *core.TxOutput) error {
		if accumulated >= amount {
			return nil // already have enough; keep draining the callback cheaply
		}
		if !bytesEqual(out.PubKeyHash, pubKeyHash) {
			return nil
		}
		txIDHex, outIdx := splitUTXOKey(key)
		unspentOutputs[txIDHex] = append(unspentOutputs[txIDHex], outIdx)
		accumulated += out.Value
		return nil
	})

	return accumulated, unspentOutputs
}

// BalanceOf sums every unspent output belonging to address. Unlike
// FindSpendableOutputs, it cannot stop early — a balance is the sum of
// EVERYTHING owned, so every matching entry must be visited.
func (u *UTXOSet) BalanceOf(address string) int64 {
	pubKeyHash, err := addressToPubKeyHash(address)
	if err != nil {
		return 0
	}

	iterable, ok := u.store.(utxoIterable)
	if !ok {
		return 0
	}

	var total int64
	iterable.ForEachUTXO(func(_ string, out *core.TxOutput) error {
		if bytesEqual(out.PubKeyHash, pubKeyHash) {
			total += out.Value
		}
		return nil
	})
	return total
}

func bytesEqual(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func splitUTXOKey(key string) (txIDHex string, outputIndex int) {
	for i := len(key) - 1; i >= 0; i-- {
		if key[i] == ':' {
			txIDHex = key[:i]
			fmt.Sscanf(key[i+1:], "%d", &outputIndex)
			return
		}
	}
	return key, 0
}

// addressToPubKeyHash reverses the Base58Check-style address encoding from
// Chapter 14: a version byte, then the pubKeyHash, then a 4-byte checksum,
// all Base58-encoded. This mirrors gochain/crypto's own address decoding
// logic exactly, kept local here so storage.UTXOSet does not need to depend
// on the wallet package just to strip a checksum.
func addressToPubKeyHash(address string) ([]byte, error) {
	decoded := crypto.Base58Decode(address)
	if len(decoded) < 5 {
		return nil, fmt.Errorf("address %q decodes to fewer than 5 bytes", address)
	}
	return decoded[1 : len(decoded)-4], nil
}
```

Notice that `BalanceOf` is still, technically, a full scan — just a scan over the `utxo` bucket, not over the entire chain. This is the honest nuance worth sitting with: the index turns "scan every transaction that ever existed" into "scan every output that currently exists," and in a real chain those two numbers diverge enormously over time, because most historical outputs have already been spent and removed from the index. A chain with 750,000 total historical transactions might have only a few thousand outputs *currently* unspent at any moment — that difference is the entire performance win. Exercise 8 pushes this further, toward a true O(1) lookup keyed directly by address.

---

## 8. Benchmark: Scan vs. Indexed Lookup

Put real numbers next to the claim. The benchmark below builds a synthetic chain of 20,000 blocks (roughly 300,000 transactions, following the same ~15-transactions-per-block estimate Chapter 53 used), then times a `BalanceOf` call against a single, specific address once using the naive Chapter 53-style full-chain scan, and once using `UTXOSet.BalanceOf`.

```go
// storage/utxo_bench_test.go
package storage

import "testing"

func BenchmarkNaiveBalanceOf(b *testing.B) {
	bc := generateTestChain(20_000) // helper: 20,000 blocks, ~15 tx each
	address := knownTestAddress      // an address that received a handful of payments early on

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		naiveBalanceOf(bc, address) // full scan, Chapter 53 style
	}
}

func BenchmarkUTXOSet_BalanceOf(b *testing.B) {
	bc := generateTestChain(20_000)
	store, _ := OpenBoltStore(b.TempDir() + "/bench.db")
	defer store.Close()

	utxos := NewUTXOSet(store)
	if err := utxos.Reindex(bc); err != nil {
		b.Fatal(err)
	}

	address := knownTestAddress

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		utxos.BalanceOf(address)
	}
}
```

```
$ go test ./storage/ -bench BalanceOf -benchtime=20x
goos: darwin
goarch: arm64
pkg: gochain/storage
BenchmarkNaiveBalanceOf-8          20    162481903 ns/op    38412 B/op   2 allocs/op
BenchmarkUTXOSet_BalanceOf-8       20       118742 ns/op      896 B/op   9 allocs/op
PASS
```

That is roughly **162 milliseconds** for the naive scan against **0.12 milliseconds** for the indexed lookup — the index answers the same question about **1,370 times faster**, on a chain of only 20,000 blocks. The gap does not stay fixed as the chain grows; it widens, because the naive scan's cost is proportional to the *entire* chain while the indexed lookup's cost is proportional only to the *current* UTXO set size, which grows far more slowly (most outputs eventually get spent and leave the index):

```
Chain size    Naive BalanceOf (scan)   UTXOSet.BalanceOf (indexed)   Speedup
----------    ----------------------   ----------------------------   -------
    5,000              ~40 ms                    ~0.05 ms             ~800x
   20,000             ~162 ms   <- measured        ~0.12 ms   <- measured  ~1,370x
  100,000              ~810 ms                    ~0.3 ms              ~2,700x
  500,000               ~4.0 s                    ~0.8 ms              ~5,000x

The naive column grows in a straight line with chain size, exactly as
Chapter 53 showed. The indexed column grows far more slowly, because it
tracks the size of the LIVE utxo set, not the size of history.
```

This is the concrete, measured payoff of everything Chapters 54 and 55 set up: an embedded database with indexed lookups, a `Store` interface clean enough to build a real index on top of, and now a maintained index that turns an operation with unbounded, ever-growing cost into one with a cost that barely moves as the chain ages.

---

## Summary

- Computing a balance by scanning the entire chain is even worse than finding one transaction, because it cannot stop early — it must visit every output belonging to an address, every single call.
- `UTXOSet` wraps a `storage.Store` and maintains a dedicated index of every currently-unspent output, keyed by `"<txID-hex>:<output-index>"` in the `utxo` bucket.
- Since `storage.Store` deliberately has no bulk-iteration method, `UTXOSet` detects an optional `ForEachUTXO` capability via a type assertion — the same pattern the standard library uses for `http.Flusher` — so `BoltStore` can support fast bulk scans without forcing every future `Store` implementation to.
- `Reindex` rebuilds the entire index from a chain's full history in two passes (collect every spent output, then index everything not in that set), and is meant to run once, not on every query.
- `Update` applies exactly one new block's transactions to the index — removing spent outputs, adding new ones — and is what mining and network-sync both call after every accepted block.
- `FindSpendableOutputs` stops as soon as it has accumulated enough value; `BalanceOf` cannot stop early, but still only scans the current UTXO set, not the entire chain's history.
- A benchmark on a 20,000-block chain measured the indexed lookup at roughly 1,370 times faster than the naive full-chain scan, with the gap widening as the chain grows, because the index's cost tracks the live UTXO set's size, not the chain's entire history.

---

## Exercises

### Easy

1. Trace through `Update`'s behavior for a single coinbase transaction (no inputs, one output paying the miner). Confirm, in your own words, that the `for _, in := range tx.Vin` loop simply never executes, and only the output loop adds a new entry.
2. `splitUTXOKey` scans a key string backward looking for the last `:` character. Explain why scanning backward (rather than forward, stopping at the *first* `:`) is the correct choice, given that a transaction ID is a hex string that could theoretically be adjacent to other characters (it isn't here, but justify the choice anyway).
3. Run `BenchmarkUTXOSet_BalanceOf` (or a smaller version you build with `generateTestChain(2_000)`) on your own machine and record the actual `ns/op`. How does the *ratio* between your naive and indexed numbers compare to the ~1,370x measured in Section 8?

### Medium

4. Write a test, `TestUTXOSet_UpdateThenReindexAgree`, that builds a small chain, calls `Update` once per block as each is added, then separately calls `Reindex` on the finished chain into a *second*, fresh `UTXOSet`, and asserts both indexes contain exactly the same set of keys. This is exactly the property that should hold if incremental updates and full rebuilds are computing the same thing two different ways.
5. `FindSpendableOutputs`'s early-exit check (`if accumulated >= amount { return nil }`) still runs the callback for every remaining entry in the bucket — it just does nothing once satisfied. Rewrite `ForEachUTXO`'s contract (or add a new method) so that a callback can signal "stop iterating entirely," and use it to make `FindSpendableOutputs` genuinely stop early, the way `FindTransaction` in Chapter 53 did.
6. Estimate, using the growth pattern from Section 8's table, how many milliseconds `BalanceOf` (indexed) would take on a chain of 5,000,000 blocks, assuming the UTXO set's size grows roughly with the square root of chain size (a reasonable approximation, since older outputs are progressively more likely to have been spent). Compare this estimate to the naive scan's cost at the same chain size from Chapter 53's table.

### Hard

7. Design and implement a *secondary* index — a separate bucket keyed directly by `pubKeyHash`, mapping to a list of UTXO keys that address owns — maintained alongside the primary `utxo` bucket by both `Update` and `Reindex`. Show how this turns `BalanceOf` from "scan every entry in the utxo bucket" into "look up one key, then fetch only the entries that belong to this address" — a genuine O(1)-ish lookup rather than O(current UTXO set size).
8. The two-pass `Reindex` algorithm holds an entire `spent` map in memory for the whole first pass. Estimate, for a chain with 10 million total historical transactions (each producing on average 2 outputs, so roughly 20 million spent-output keys in the worst case), roughly how much memory that map would consume, and propose (in comments, no need to fully implement) an alternative approach that would not require holding the entire spent set in memory at once.
9. `UTXOSet.Update` and `Reindex` are not safe to call concurrently with each other on the same `Store` — nothing prevents a `Reindex` in progress from racing with an `Update` triggered by a newly mined block arriving mid-rebuild. Design (and, if you like, implement) a locking strategy inside `UTXOSet` itself that serializes these two operations correctly, and explain why putting this lock inside `UTXOSet` rather than inside `BoltStore` is the more correct layer for it.
