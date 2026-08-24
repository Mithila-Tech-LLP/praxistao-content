# Chapter 34: The Mempool and Preventing Double-Spending

A signed, verified transaction does not go straight into a block — it waits in a holding area until a miner picks it up. This chapter builds that holding area, `core.Mempool`, and gives it the one job that matters most: catching an attacker who tries to spend the same coins twice before anyone notices.

## Table of Contents

1. [What Problem Does the Mempool Solve?](#1-what-problem-does-the-mempool-solve)
2. [The Anatomy of a Pending Transaction](#2-the-anatomy-of-a-pending-transaction)
3. [Designing core.Mempool](#3-designing-coremempool)
4. [Detecting Double-Spends: The Spent-Outpoints Set](#4-detecting-double-spends-the-spent-outpoints-set)
5. [Implementing Add, Remove, and GetPending](#5-implementing-add-remove-and-getpending)
6. [A Deliberate Double-Spend — and Watching It Get Rejected](#6-a-deliberate-double-spend--and-watching-it-get-rejected)
7. [What the Mempool Does Not Check](#7-what-the-mempool-does-not-check)
8. [Removing Transactions Once They're Mined](#8-removing-transactions-once-theyre-mined)
9. [Concurrency: Why the Mempool Needs a Mutex](#9-concurrency-why-the-mempool-needs-a-mutex)
10. [Testing the Mempool End to End](#10-testing-the-mempool-end-to-end)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What Problem Does the Mempool Solve?

Think about an airport security line. Passengers don't teleport straight from the terminal entrance onto a plane — they queue, get checked, and wait until a gate agent is ready to actually board them. A blockchain has the exact same shape of problem: a transaction gets created and signed by a wallet at some arbitrary moment, but blocks only get mined every so often (every ten seconds, every ten minutes, whatever GoChain's target block time is from Volume 4). Something has to hold every transaction that has shown up *between* blocks, so a miner has something to choose from when it's time to build the next one.

That holding area is the **mempool** — short for "memory pool." It is exactly what it sounds like: a pool, kept in memory, of transactions that are valid and signed but have not yet been included in a mined block. Every node runs its own mempool. There is no single, shared, canonical mempool anywhere — each node's view is built from whatever transactions have reached it so far (a detail that becomes concrete once nodes talk to each other over the network, starting in Volume 7).

The mempool has three jobs:

1. **Hold pending transactions** until a miner is ready to build a block.
2. **Let a miner retrieve everything pending** so it can choose what to include (Chapter 35 adds fee-based prioritization on top of this).
3. **Reject a transaction that conflicts with another transaction already waiting** — this is the job this chapter is really about.

That third job deserves its own name: **double-spend prevention**. A double-spend is exactly what it sounds like — an attempt to spend the same coins twice. Since GoChain uses the UTXO model (Volume 5, Chapter 30), "the same coins" means, precisely, the same unspent transaction output. If two different transactions both try to consume the same UTXO as an input, at most one of them can ever be valid, because the moment either one is mined, that UTXO stops existing. The mempool is the first line of defense that catches this — before a miner ever has to think about it, before a block is ever built, before anything is written to the chain.

```
                    +-------------------+
Wallet signs a  --> |      Mempool      | --> Miner picks pending
transaction         |  (waiting room)   |     transactions for
                    +-------------------+     the next block
                            |
                            v
                 rejects conflicting /
                 already-claimed spends
```

---

## 2. The Anatomy of a Pending Transaction

Before writing `core.Mempool`, it's worth re-grounding what exactly is sitting in this waiting room. From Chapters 32 and 33, a `core.Transaction` looks like this:

```go
package core

type Transaction struct {
	ID        []byte
	Inputs    []TxInput
	Outputs   []TxOutput
	Timestamp int64
}

type TxInput struct {
	TxID      []byte
	OutIndex  int
	Signature []byte
	PublicKey []byte
}

type TxOutput struct {
	Value      int64
	PubKeyHash []byte
}
```

Every `TxInput` is a claim: "I am spending output number `OutIndex` of the transaction with ID `TxID`." That pair — `TxID` plus `OutIndex` — uniquely identifies one specific unspent output anywhere in GoChain's history. This pair has a name worth learning now, because it is the core unit the rest of this chapter revolves around: an **outpoint**. An outpoint is a pointer to one particular coin-sized chunk of value, the same way a claim ticket at a coat check points to one specific coat. Two transactions with inputs pointing at the *same outpoint* are, by definition, both trying to claim the same coat.

A transaction that reaches the mempool has already passed `Transaction.Verify()` from Chapter 33 — its signatures are legitimate, and it really was authorized by whoever controls the private keys behind its inputs. What `Verify()` does *not* check is whether some other transaction, sitting in the mempool right now, has already claimed the exact same outpoint. That is a mempool-level concern, not a single-transaction concern, because it requires comparing a transaction against everything else waiting alongside it.

---

## 3. Designing core.Mempool

Here is the type this chapter builds, exactly as specified by GoChain's shared contract:

```go
package core

type Mempool struct {
	pending map[string]*Transaction // keyed by hex transaction ID
	spent   map[string]bool         // keyed by "txid:outindex" to detect conflicting spends
}

func NewMempool() *Mempool
func (mp *Mempool) Add(tx *Transaction) error // rejects double-spends
func (mp *Mempool) Remove(txID []byte)
func (mp *Mempool) GetPending() []*Transaction
```

Two maps do all the work. `pending` is straightforward: it stores every waiting transaction, keyed by its ID (hex-encoded, since `[]byte` cannot be a Go map key but a string can). `spent` is the interesting one — it does not store transactions at all. It stores **outpoints**, encoded as a string like `"a1b2c3...:0"`, and the only thing it records about each one is that some pending transaction has already claimed it. Think of `spent` as a bouncer's clipboard: it doesn't care who you are, only whether your name (your outpoint) is already checked off the list.

```
Mempool
+-----------------------------+     +-----------------------------+
| pending (map)                |     | spent (map)                  |
|-------------------------------|     |-------------------------------|
| "a1b2..." -> *Transaction     |     | "f00d...:0" -> true           |
| "c3d4..." -> *Transaction     |     | "f00d...:1" -> true           |
|                               |     | "beef...:0" -> true           |
+-----------------------------+     +-----------------------------+
   full transactions, keyed          just outpoints, keyed by
   by their own ID                   "txid:index", no transaction
                                      data needed
```

---

## 4. Detecting Double-Spends: The Spent-Outpoints Set

Here's the everyday version of the attack this section defends against. Imagine you have exactly one $20 bill, and you try to buy something from two different online stores at nearly the same instant, using the same $20 bill's serial number as "proof of payment" in both. If neither store can see the other's pending order, both might tentatively accept it — until whichever store's warehouse ships first "spends" the bill for real, and the second store discovers, too late, that the serial number was already used. A single, shared point of bookkeeping — one ledger both stores check before shipping — would have caught this immediately.

The mempool is exactly that shared point of bookkeeping, but for one node's view of the network. The rule is simple: **the first transaction to claim an outpoint wins; any later transaction claiming the same outpoint is rejected outright**, without needing to inspect signatures, addresses, or amounts. It doesn't matter whether the second transaction is otherwise perfectly valid — the outpoint it wants no longer has a claim available.

```
UTXO: transaction f00d..., output #0, worth 50 gochips, owned by Alice

     Transaction A                      Transaction B
     spends f00d...:0                   spends f00d...:0
     pays 50 to Bob                     pays 50 to Carol
            \                                  /
             \                                /
              v                              v
        +---------------------------------------+
        |          Mempool.Add() checks          |
        |     "is f00d...:0 already claimed?"     |
        +---------------------------------------+
              |                              |
        A arrives first:               B arrives second:
        outpoint is free ->             outpoint already
        ACCEPTED, outpoint              claimed by A ->
        marked spent                    REJECTED
```

Whether Transaction A or Transaction B "arrives first" at any particular node depends on network timing (Volume 7 covers exactly how transactions propagate), but the important guarantee is this: **no single node will ever accept both**. That is enough to make the attack pointless in practice, because a miner cannot mine a transaction its own mempool never admitted, and — as you'll see in Chapter 37's block validation — a conflicting transaction that somehow reached a different miner would also fail the instant it tried to enter a block whose predecessor already spent that outpoint.

---

## 5. Implementing Add, Remove, and GetPending

Time to write it. Create `core/mempool.go`:

```go
package core

import (
	"encoding/hex"
	"fmt"
	"sync"
)

// Mempool holds transactions that are signed and verified but not yet
// included in a mined block. Every node keeps its own mempool; nodes do
// not share one global mempool directly — they only learn about each
// other's pending transactions through gossip (Volume 7).
type Mempool struct {
	mu      sync.Mutex               // protects both maps below from concurrent access
	pending map[string]*Transaction  // keyed by hex transaction ID
	spent   map[string]bool          // keyed by "txid:outindex" to detect conflicting spends
}

// NewMempool returns an empty mempool, ready to accept transactions.
func NewMempool() *Mempool {
	return &Mempool{
		pending: make(map[string]*Transaction),
		spent:   make(map[string]bool),
	}
}

// outpointKey turns a (txID, outIndex) pair into a single comparable
// string, since Go maps need a plain key type, not a struct with a
// []byte field ([]byte itself isn't a valid map key).
func outpointKey(txID []byte, outIndex int) string {
	return fmt.Sprintf("%s:%d", hex.EncodeToString(txID), outIndex)
}

// Add inserts tx into the mempool, rejecting it outright if any of its
// inputs claims an outpoint another pending transaction has already
// claimed — this is GoChain's mempool-level double-spend defense.
func (mp *Mempool) Add(tx *Transaction) error {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	idHex := hex.EncodeToString(tx.ID)
	if _, exists := mp.pending[idHex]; exists {
		return fmt.Errorf("transaction %s is already in the mempool", idHex)
	}

	// First pass: check every input BEFORE reserving any of them. This
	// matters — if we reserved input 1, then found input 2 conflicts, a
	// rejected transaction would still have left a stray reservation
	// behind for input 1.
	for _, in := range tx.Inputs {
		key := outpointKey(in.TxID, in.OutIndex)
		if mp.spent[key] {
			return fmt.Errorf("double-spend detected: outpoint %s is already claimed by a pending transaction", key)
		}
	}

	// Second pass: now that we know every input is free, claim all of
	// them at once.
	for _, in := range tx.Inputs {
		mp.spent[outpointKey(in.TxID, in.OutIndex)] = true
	}

	mp.pending[idHex] = tx
	return nil
}

// Remove takes a transaction out of the mempool and frees the outpoints
// it had claimed, so future transactions spending those same outpoints
// (which should no longer exist once this transaction is mined — see
// Section 8) don't get stuck permanently reserved.
func (mp *Mempool) Remove(txID []byte) {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	idHex := hex.EncodeToString(txID)
	tx, exists := mp.pending[idHex]
	if !exists {
		return // already gone; nothing to do
	}

	for _, in := range tx.Inputs {
		delete(mp.spent, outpointKey(in.TxID, in.OutIndex))
	}
	delete(mp.pending, idHex)
}

// GetPending returns every transaction currently waiting to be mined.
// Order is not guaranteed here — Chapter 35 adds fee-based ordering on
// top of this raw list.
func (mp *Mempool) GetPending() []*Transaction {
	mp.mu.Lock()
	defer mp.mu.Unlock()

	txs := make([]*Transaction, 0, len(mp.pending))
	for _, tx := range mp.pending {
		txs = append(txs, tx)
	}
	return txs
}
```

Walking through this by name: `outpointKey` is the small helper that turns a `TxInput`'s `TxID`/`OutIndex` pair into the string key both `Add` and `Remove` use to talk to the `spent` map. `Add` does its double-spend check in two clearly separated passes — check everything first, reserve everything second — specifically so a rejected transaction never leaves the mempool in a half-updated state. `Remove` is the mirror image of `Add`: it frees every outpoint the transaction had reserved before deleting the transaction itself, so the bookkeeping in `spent` never goes stale. `GetPending` just hands back a plain slice, copied out from under the lock, so callers can iterate it freely without holding the mempool's internal lock themselves.

---

## 6. A Deliberate Double-Spend — and Watching It Get Rejected

Talk is cheap; let's watch the rejection happen. This test builds two transactions that both claim the exact same fake outpoint and confirms only the first one is accepted:

```go
package core

import "testing"

func TestMempool_RejectsDoubleSpend(t *testing.T) {
	mp := NewMempool()

	// Both transactions below claim the SAME outpoint: output #0 of a
	// (fake, for this test) prior transaction. In a real chain this
	// would be someone's actual UTXO — say, 50 gochips sitting in
	// Alice's wallet, waiting to be spent.
	sharedOutpoint := TxInput{TxID: []byte("prior-tx-id"), OutIndex: 0}

	txToBob := &Transaction{
		ID:      []byte("tx-to-bob"),
		Inputs:  []TxInput{sharedOutpoint},
		Outputs: []TxOutput{{Value: 50, PubKeyHash: []byte("bob-hash")}},
	}
	txToCarol := &Transaction{
		ID:      []byte("tx-to-carol"),
		Inputs:  []TxInput{sharedOutpoint}, // <- claims the same coin as txToBob
		Outputs: []TxOutput{{Value: 50, PubKeyHash: []byte("carol-hash")}},
	}

	// The first transaction to arrive gets the outpoint.
	if err := mp.Add(txToBob); err != nil {
		t.Fatalf("expected txToBob to be accepted, got error: %v", err)
	}

	// The second transaction is a genuine double-spend attempt and must
	// be rejected, even though nothing is wrong with it in isolation.
	err := mp.Add(txToCarol)
	if err == nil {
		t.Fatal("expected txToCarol to be rejected as a double-spend, but Add() succeeded")
	}
	t.Logf("txToCarol correctly rejected: %v", err)

	pending := mp.GetPending()
	if len(pending) != 1 {
		t.Fatalf("expected exactly 1 pending transaction after the rejection, got %d", len(pending))
	}
	if hex.EncodeToString(pending[0].ID) != hex.EncodeToString(txToBob.ID) {
		t.Fatal("the surviving pending transaction should be txToBob")
	}
}
```

`sharedOutpoint` is deliberately the exact same `TxInput` value used in both transactions, standing in for "Alice's 50-gochip UTXO" from the airport-line and $20-bill analogies above. `txToBob` and `txToCarol` are two otherwise-valid-looking transactions that both try to spend it — one pays Bob, the other pays Carol, but they cannot both succeed, because there is only one coin. The test calls `mp.Add(txToBob)` first and expects success — this is the honest transaction winning the race. It then calls `mp.Add(txToCarol)` and expects an error — this is the double-spend attempt getting caught. Finally, it confirms the mempool still has exactly one pending transaction, and that it's the one that should have survived. Run it with `go test ./core/... -run TestMempool_RejectsDoubleSpend -v` and you'll see the rejection message logged, proving the defense works before a single block has been mined.

---

## 7. What the Mempool Does Not Check

It's worth being precise about the boundary of this chapter's responsibility, because the mempool is deliberately narrow in scope:

- **It does not verify signatures.** That already happened via `Transaction.Verify()` (Chapter 33) before a transaction is ever offered to `mp.Add()`. A well-behaved node calls `Verify()` first and only calls `Add()` on transactions that already passed.
- **It does not check whether an outpoint was already spent by a *mined* transaction.** This chapter's `spent` map only tracks conflicts *among pending transactions*. Checking against the blockchain's actual UTXO set (has this output already been consumed by a block that's already on the chain?) is the UTXO set's job, built in Chapter 32 and used again during block validation in Chapter 37. In a complete node, a transaction is checked against *both*: the UTXO set (is this outpoint real and currently unspent on-chain?) and the mempool (is anyone else already trying to spend it right now?).
- **It does not order transactions by anything meaningful yet.** `GetPending()` returns them in whatever order Go's map iteration happens to produce, which is intentionally randomized. Chapter 35 fixes this by sorting pending transactions by fee.
- **It does not persist to disk.** If a node restarts, its mempool is empty again. This is normal and matches real blockchains — the mempool is explicitly a transient, in-memory structure, not part of the permanent record. Anything genuinely pending will simply be rebroadcast by wallets or re-gossiped by peers once the node comes back up (Volume 7).

Real-world grounding: Bitcoin's mempool routinely holds tens of thousands of pending transactions during busy periods, and different nodes' mempools are never perfectly identical at any given instant — one node might have just heard about a transaction that hasn't reached another node yet. GoChain's single-node mempool in this chapter is the simplest possible case (one node, no network), which is exactly why it's the right place to build and test the double-spend logic in isolation, before Volume 7 adds the networked complexity of many mempools disagreeing about what they've each seen.

---

## 8. Removing Transactions Once They're Mined

A transaction's stay in the mempool should be temporary. The moment a miner successfully includes it in a mined block (`core.Blockchain.MineBlock`, built in Volume 4), it needs to come out of the mempool — otherwise it would sit there forever, and worse, its outpoint reservation in `spent` would permanently block any future (legitimate) attempt to reuse that entry once it's already settled on-chain.

The pattern every miner follows looks like this:

```go
pending := mempool.GetPending()
block := bc.MineBlock(pending) // proof-of-work + Merkle root, from Volume 4

// Now that these transactions are permanently recorded in a mined
// block, they no longer belong in the "waiting room."
for _, tx := range pending {
	mempool.Remove(tx.ID)
}
```

`bc.MineBlock(pending)` takes the transactions we retrieved and produces a real, proof-of-work-secured block containing them — this is the exact `Blockchain.MineBlock` method from Volume 4, unchanged. The loop afterward is the cleanup step: for every transaction that just got permanently written into the chain, we call `mempool.Remove` so it stops occupying space in `pending` and, just as importantly, so its claimed outpoints are freed from `spent`. Chapter 37's end-to-end demo exercises this exact sequence for real.

---

## 9. Concurrency: Why the Mempool Needs a Mutex

Recall from Volume 1 (Chapter 05) that a real GoChain node does many things "at once": it might be receiving a freshly gossiped transaction from a peer on one goroutine while a miner goroutine is simultaneously calling `GetPending()` to decide what to mine next. If two goroutines touched `pending` or `spent` at the same time without coordination, Go's race detector would (rightly) scream — maps are not safe for concurrent read/write access in Go, and a half-updated map can corrupt or crash the program.

That's what the `sync.Mutex` embedded in `Mempool` is for. Every exported method — `Add`, `Remove`, `GetPending` — begins with `mp.mu.Lock()` and immediately `defer mp.mu.Unlock()`, ensuring only one goroutine is ever inside the mempool's internal state at a time. This is a small, deliberate piece of defensive engineering: nothing about single-threaded testing would have caught its absence, but the moment gossip and mining run concurrently (which they will, starting in Volume 7), an unprotected mempool would produce rare, maddening, hard-to-reproduce bugs. Paying this cost now, while it's cheap and obvious, is far better than debugging it later under real network load.

You can verify this yourself with Go's built-in race detector:

```bash
go test ./core/... -race -run TestMempool
```

If two goroutines ever touched the mempool's maps without the mutex, `-race` would report a data race immediately. With the mutex in place, it reports nothing — exactly what we want.

---

## 10. Testing the Mempool End to End

Beyond the double-spend test in Section 6, a solid mempool test suite should cover the full lifecycle:

```go
func TestMempool_AddRemoveGetPending(t *testing.T) {
	mp := NewMempool()

	tx := &Transaction{
		ID:      []byte("tx-1"),
		Inputs:  []TxInput{{TxID: []byte("utxo-1"), OutIndex: 0}},
		Outputs: []TxOutput{{Value: 25, PubKeyHash: []byte("alice-hash")}},
	}

	// A brand-new transaction should be accepted.
	if err := mp.Add(tx); err != nil {
		t.Fatalf("unexpected error adding a fresh transaction: %v", err)
	}
	if len(mp.GetPending()) != 1 {
		t.Fatalf("expected 1 pending transaction, got %d", len(mp.GetPending()))
	}

	// Adding the exact same transaction twice should fail — it's
	// already there, not a fresh double-spend candidate.
	if err := mp.Add(tx); err == nil {
		t.Fatal("expected an error re-adding an already-pending transaction")
	}

	// Removing it should free the transaction AND its outpoint.
	mp.Remove(tx.ID)
	if len(mp.GetPending()) != 0 {
		t.Fatal("expected 0 pending transactions after Remove")
	}

	// Because Remove freed the outpoint, a NEW transaction spending the
	// same outpoint (simulating that it was legitimately mined and this
	// is an unrelated, later re-check) should now be accepted again.
	txAgain := &Transaction{
		ID:      []byte("tx-2"),
		Inputs:  []TxInput{{TxID: []byte("utxo-1"), OutIndex: 0}},
		Outputs: []TxOutput{{Value: 25, PubKeyHash: []byte("bob-hash")}},
	}
	if err := mp.Add(txAgain); err != nil {
		t.Fatalf("expected the freed outpoint to be reusable, got error: %v", err)
	}
}
```

This test exercises the whole loop: a transaction is accepted, a duplicate add is rejected, removal frees both the transaction and its reserved outpoint, and a fresh transaction can then claim that same freed outpoint. In production this last scenario matters because after `Remove` is called following a successful mine (Section 8), the *real* UTXO backing that outpoint is gone forever — any transaction that later tries to spend it would be rejected by the UTXO-set check instead, not the mempool. The mempool's job ends at "nobody else pending is claiming this right now."

---

## Summary

- The **mempool** is where signed, verified, not-yet-mined transactions wait; every node keeps its own, and it exists purely in memory.
- An **outpoint** — a `(TxID, OutIndex)` pair — uniquely identifies one specific coin-sized chunk of value; double-spend detection is fundamentally about outpoint conflicts.
- `core.Mempool` tracks two things: `pending` (full transactions, keyed by ID) and `spent` (just outpoints already claimed by something pending, keyed by `"txid:outindex"`).
- `Add` rejects a transaction outright if any of its inputs claims an outpoint another pending transaction already claimed — checking every input before reserving any of them, so rejections never leave partial state behind.
- `Remove` frees a transaction's claimed outpoints before deleting it, keeping the `spent` map accurate as transactions get mined.
- The mempool does **not** verify signatures, does **not** check against the on-chain UTXO set, does **not** order by fee, and does **not** persist to disk — each of those is a separate, deliberate responsibility that lives elsewhere.
- A `sync.Mutex` protects the mempool's maps, because a real node's mining and gossip logic will touch it from multiple goroutines concurrently.
- A deliberately constructed double-spend test — two transactions claiming the same outpoint — proves the defense works before it ever has a chance to matter on a real network.

---

## Exercises

### Easy

1. Write a test, `TestMempool_RemoveNonexistent`, that calls `mp.Remove()` with a transaction ID that was never added, and confirms it does not panic and leaves the mempool's pending count unchanged.
2. Add a method `func (mp *Mempool) Size() int` that returns the number of pending transactions without needing to call `GetPending()` and take `len()` of the result. Make sure it also locks the mutex correctly.
3. Modify `TestMempool_RejectsDoubleSpend` so that instead of two transactions each spending one outpoint, `txToBob` spends *two* outpoints and `txToCarol` conflicts on only the *second* one. Confirm `Add` still rejects `txToCarol` and that neither of `txToBob`'s outpoints leaked into an inconsistent state.

### Medium

4. Right now, `Add` silently overwrites nothing and just returns an error on a duplicate transaction ID. Change the error messages in `Add` to distinguish clearly, in the returned error text, between "duplicate transaction already pending" and "double-spend against a different pending transaction," and write a test asserting each specific error path via `errors.Is` or a string check.
5. Add a method `func (mp *Mempool) Conflicts(tx *Transaction) []*Transaction` that, given a candidate transaction, returns every *currently pending* transaction it would conflict with (without modifying the mempool). Use it to write a test where a single new transaction conflicts with two different existing pending transactions at once (spending outpoints claimed by each).
6. Benchmark `Add` and `GetPending` using Go's `testing.B` with 10,000 pending transactions already in the mempool. Report how `GetPending`'s allocation (via `make([]*Transaction, 0, len(mp.pending))`) compares to a naive version that starts with a nil slice and appends without pre-allocating capacity.

### Hard

7. The current design permanently loses a rejected double-spend transaction — the caller just gets an error and the transaction is discarded. Design (and implement) a small "orphan" tracking structure that remembers the losing transaction of a double-spend attempt for a short time, so that if the *winning* transaction's block later gets reorganized out of the chain (Volume 7, Chapter 50 — forks), the previously-rejected transaction could be reconsidered. Write a test simulating this exact sequence.
8. Extend `Mempool` with a maximum size limit (e.g., 5,000 transactions) and an eviction policy for when it's full: when a new, valid transaction arrives and the mempool is already at capacity, evict the single lowest-value pending transaction (sum of its outputs) to make room, unless the new transaction's own value is even lower, in which case reject the new one instead. Write tests covering both eviction and rejection paths.
9. Using Go's `-race` flag, write a stress test that spins up 50 goroutines simultaneously calling `Add` with transactions that deliberately have a 10% chance of conflicting with an outpoint another goroutine is also submitting, and 50 more goroutines simultaneously calling `GetPending` and `Remove` in a loop. Run it under `-race` and confirm it reports no data races, then explain in a short comment why the two-pass check-then-reserve pattern in `Add` is still safe under this concurrent load even though it is not a single atomic map operation.
