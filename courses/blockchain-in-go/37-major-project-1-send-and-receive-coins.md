# Chapter 37: Major Project 1 — Send and Receive Coins

Every chapter since Volume 2 has built one specific piece: hashing, signatures, addresses, the block structure, proof of work, transactions, the mempool. Each piece was tested in isolation, on its own, with placeholder data standing in for whatever the surrounding pieces hadn't been built yet. This chapter builds nothing new at the cryptographic or algorithmic level — instead, it wires everything built across Volumes 2 through 5 into one continuous, runnable program, and watches real value move from one wallet to another, get permanently recorded, and get independently re-verified from scratch.

## Table of Contents

1. [Bringing Volumes 2-5 Together](#1-bringing-volumes-2-5-together)
2. [Why the Coinbase Transaction Has No Inputs](#2-why-the-coinbase-transaction-has-no-inputs)
3. [A New Helper: Blockchain.FindTransaction](#3-a-new-helper-blockchainfindtransaction)
4. [Walking Through the Demo Before Running It](#4-walking-through-the-demo-before-running-it)
5. [Major Project: Send-and-Receive Demo](#major-project-send-and-receive-demo)
6. [Checking the Arithmetic](#6-checking-the-arithmetic)
7. [What This Demo Does Not Yet Do](#7-what-this-demo-does-not-yet-do)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Bringing Volumes 2-5 Together

Here is the full inventory of what this one program is about to exercise, and where each piece came from:

```
crypto.KeyPair, Sign, Verify         <- Volume 2 (Chapters 11-13)
address encoding (Base58, checksum)  <- Volume 2 (Chapter 14)
core.Block, core.Blockchain          <- Volume 3 (Chapters 16-20)
consensus.ProofOfWork                <- Volume 4 (Chapters 24-25)
difficulty adjustment                <- Volume 4 (Chapter 26)
concurrent mining                     <- Volume 4 (Chapter 27)
core.Transaction, TxInput, TxOutput  <- Volume 5 (Chapter 32)
core.UTXOSet                         <- Volume 5 (Chapter 32)
Transaction.Sign, Transaction.Verify <- Volume 5 (Chapter 33)
core.Mempool                         <- Volume 5 (Chapter 34)
wallet.New, wallet.Wallet            <- Volume 5/6 (Chapter 36 preview)
```

Nothing on this list is new. What's new is watching all of it cooperate in one continuous run: two real key pairs, a real proof-of-work-secured chain, a real signed transaction, and a real balance recomputed independently at the end — not a unit test asserting one function's return value in isolation, but an actual, small, working economy.

---

## 2. Why the Coinbase Transaction Has No Inputs

Every transaction this course has built so far — starting with Chapter 32's `NewTransaction` — follows the same shape: consume one or more existing UTXOs as inputs, produce one or more new UTXOs as outputs. That shape has a hard requirement baked into it: an input can only exist by referencing a *specific, real, previously-created* output. You cannot spend a UTXO that was never created.

This creates an obvious chicken-and-egg problem the very first time any value needs to exist on a brand-new chain at all: **before anyone has ever received a payment, there is nothing anyone owns, and therefore nothing anyone could possibly reference as an input.** If every transaction required consuming an existing UTXO, GoChain's total supply of gochips would be permanently and inescapably zero — a perfectly consistent, perfectly useless currency that nobody could ever actually hold any of.

A **coinbase transaction** is GoChain's answer, and it is deliberately, structurally different from every other transaction in exactly one way: it has **no inputs**. Recall its exact shape from Chapter 32, Section 9:

```go
func NewCoinbaseTX(toAddress string, reward int64) *Transaction {
	txIn := TxInput{TxID: []byte{}, OutIndex: -1}
	txOut := NewTXOutput(reward, toAddress)

	tx := &Transaction{
		Inputs:    []TxInput{txIn},
		Outputs:   []TxOutput{*txOut},
		Timestamp: time.Now().Unix(),
	}
	tx.ID = tx.Hash()
	return tx
}
```

Here is the "why" in full, plain language, worth internalizing precisely because it's the single most commonly confused idea for newcomers to blockchain economics: **an ordinary transaction moves existing value from one place to another. A coinbase transaction creates brand-new value out of nothing, by protocol rule, as the reward for the specific, real-world work of mining a block.** These are not two variations on the same operation — they are opposites. An ordinary transaction's total output value can never exceed its total input value (Chapter 30's conservation-of-value principle — you cannot legally spend more than you have). A coinbase transaction has *no* input value to be bounded by, because it isn't *moving* anything that existed before it ran; it's minting.

```
ORDINARY TRANSACTION                    COINBASE TRANSACTION

  Input:  Alice's 50-gochip UTXO         Input:  NONE -- nothing to
          (already existed,                      reference, because
           already belonged to                   this transaction is
           someone)                               not spending anyone's
                                                    existing value
  Output: 20 -> Bob                      Output: 50 -> Miner
          30 -> Alice (change)                    (BRAND NEW value,
                                                    created by protocol
  Total value MOVED: 50                            rule, not moved from
  Total value CREATED: 0                            anywhere)

                                          Total value MOVED: 0
                                          Total value CREATED: 50
```

This is precisely how new gochips enter circulation at all, and it mirrors exactly how new bitcoins enter circulation in the real Bitcoin network: the reward for successfully mining a block is a brand-new coinbase transaction, paid to whoever mined it, created from nothing but the protocol's own rule that mining a valid block earns a fixed reward. `IsCoinbase()` (also from Chapter 32) is what lets the rest of GoChain's code — the UTXO scanner, `Sign`, `Verify` — recognize this special case and skip the "does this input have a valid signature over a real previous output" checks that would otherwise make no sense applied to a transaction with no inputs at all.

One detail worth flagging explicitly, since Section 5's demo relies on it: **every block a miner successfully mines earns them a fresh coinbase reward**, not just the very first block ever mined on a chain. Mining block 1 pays a coinbase reward; mining block 2 pays another, brand-new one; and so on for every block, forever (real Bitcoin's reward amount actually shrinks over time via periodic "halvings" — a detail outside this chapter's scope, but the "every block gets its own fresh coinbase reward" mechanic is identical).

---

## 3. A New Helper: Blockchain.FindTransaction

Chapter 33's `Transaction.Sign` needs a `prevTXs map[string]Transaction` — a lookup from a hex-encoded transaction ID to the full transaction that created whatever output is being spent. Nothing built so far actually *produces* that map for a real, caller-facing scenario; Chapter 33's own examples assumed it was simply "available." This chapter closes that gap with one small addition to `core.Blockchain`:

```go
package core

import (
	"encoding/hex"
	"errors"
)

// FindTransaction scans the chain for a transaction with the given ID,
// returning it if found. This is exactly the same "walk every block,
// every transaction" scan UTXOSet already does (Chapter 32) -- just
// searching by transaction ID instead of by which outputs are unspent.
func (bc *Blockchain) FindTransaction(id []byte) (*Transaction, error) {
	idHex := hex.EncodeToString(id)

	for _, block := range bc.blocks {
		for _, tx := range block.Transactions {
			if hex.EncodeToString(tx.ID) == idHex {
				return tx, nil
			}
		}
	}

	return nil, errors.New("core: transaction not found on this chain")
}
```

`FindTransaction` is a plain linear scan — the same "correctness first, speed later" philosophy Chapter 32's `UTXOSet` already adopted, and the same one Volume 8 eventually replaces with an indexed lookup. For a demo chain with a handful of blocks, scanning every transaction to find one by ID is instant; for a chain with millions of transactions, this exact function is exactly what Chapter 56's real index makes fast, without changing what `FindTransaction` *means*.

With this in place, building the `prevTXs` map Chapter 33's `Sign` needs becomes a small, mechanical helper any caller can write:

```go
package main

import (
	"encoding/hex"
	"log"

	"github.com/you/gochain/core"
)

// findPrevTXs looks up, for every input tx wants to spend, the full
// earlier transaction that created the output being spent -- exactly
// what Transaction.Sign (Chapter 33) needs to reconstruct each input's
// PubKeyHash before hashing and signing.
func findPrevTXs(bc *core.Blockchain, tx *core.Transaction) map[string]core.Transaction {
	prevTXs := make(map[string]core.Transaction)

	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.TxID)
		if err != nil {
			log.Fatalf("findPrevTXs: could not find transaction %x: %v", in.TxID, err)
		}
		prevTXs[hex.EncodeToString(in.TxID)] = *prevTX
	}

	return prevTXs
}
```

---

## 4. Walking Through the Demo Before Running It

Before looking at the full program, here is the shape of what it does, step by step, so the code in Section 5 reads as a confirmation of something you already expect rather than a surprise:

```
1. Create Wallet 1 and Wallet 2 -- two independent key pairs, two
   independent addresses. Neither knows or needs to know about the
   other's private key at any point.

2. Open a brand-new, empty blockchain -- just a genesis block, no
   transactions, no value anywhere yet (Chapter 18).

3. Mine block 1, containing ONE transaction: a coinbase transaction
   paying Wallet 1 a reward. This is the very first value to ever
   exist on this chain -- created, not moved (Section 2).

4. Check balances: Wallet 1 has the reward. Wallet 2 has zero.

5. Wallet 1 builds, signs, and verifies a transaction sending part of
   its balance to Wallet 2 (Chapters 32-33). This transaction is
   correct and fully authorized -- but it is NOT yet on the chain.

6. Mine block 2, containing TWO transactions: a fresh coinbase reward
   for whoever mined this block (Wallet 1, again, in this demo), and
   the transaction from step 5.

7. Check balances again -- by scanning the ENTIRE chain from scratch,
   both blocks, every transaction, with no shortcuts and no cached
   numbers left over from step 4.
```

Step 7 is the detail worth sitting with for a moment: this demo deliberately never trusts "whatever balance we computed last time, plus or minus this one transaction." Every single balance check in this program walks the whole chain from the beginning, exactly as Chapter 30's "balance is a sum, not a number in a cell" principle demands. If the coinbase transaction, the signature, the mining, or the UTXO scanning had any bug anywhere, the final balance check would be the thing to catch it — not because it's specially instrumented to catch bugs, but because it is running the exact same, ordinary `BalanceOf` any real GoChain node would run at any time.

---

## Major Project: Send-and-Receive Demo

```go
// cmd/senddemo/main.go
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/you/gochain/core"
	"github.com/you/gochain/wallet"
)

// miningReward is the fixed number of gochips a coinbase transaction
// pays out per block mined -- a protocol constant, not something a
// miner or wallet chooses. Real Bitcoin's own reward shrinks over time
// via halvings; GoChain keeps it fixed for this course's purposes.
const miningReward = 50

func main() {
	// --- Step 1: two independent wallets ---
	wallet1 := wallet.New()
	wallet2 := wallet.New()
	fmt.Println("Wallet 1 address:", wallet1.Address())
	fmt.Println("Wallet 2 address:", wallet2.Address())
	fmt.Println()

	// --- Step 2: a brand-new, empty chain ---
	bc := core.NewBlockchain()

	// --- Step 3: mine block 1, rewarding Wallet 1 via a coinbase tx ---
	// No inputs here -- see Section 2 for exactly why. This is new
	// supply entering circulation, not existing value moving around.
	coinbase1 := core.NewCoinbaseTX(wallet1.Address(), miningReward)
	block1 := bc.MineBlock([]*core.Transaction{coinbase1})
	fmt.Printf("Mined block %d (coinbase reward %d gochips to Wallet 1)\n",
		block1.Height, miningReward)

	// --- Step 4: check balances ---
	printBalances(bc, wallet1, wallet2)

	// --- Step 5: Wallet 1 sends Wallet 2 20 gochips ---
	utxoSet := core.NewUTXOSet(bc)
	tx, err := core.NewTransaction(wallet1, wallet2.Address(), 20, utxoSet)
	if err != nil {
		log.Fatalf("building transaction: %v", err)
	}

	prevTXs := findPrevTXs(bc, tx)
	tx.Sign(wallet1.KeyPair.PrivateKey, prevTXs)

	if !tx.Verify(prevTXs) {
		log.Fatal("transaction failed to verify immediately after signing -- this should never happen")
	}
	fmt.Printf("Wallet 1 signed a transaction sending 20 gochips to Wallet 2 (tx %x)\n", tx.ID)
	fmt.Println()

	// --- Step 6: mine block 2, containing a fresh coinbase reward AND
	//     the pending transaction from step 5 ---
	coinbase2 := core.NewCoinbaseTX(wallet1.Address(), miningReward)
	block2 := bc.MineBlock([]*core.Transaction{coinbase2, tx})
	fmt.Printf("Mined block %d (%d transactions: 1 coinbase reward + 1 transfer)\n",
		block2.Height, len(block2.Transactions))

	// --- Step 7: verify balances by scanning the WHOLE chain from
	//     scratch, from genesis through the block just mined ---
	printBalances(bc, wallet1, wallet2)
}

// printBalances recomputes both wallets' balances from scratch, by
// scanning every block on the chain -- never from a cached number left
// over from an earlier call. This is Chapter 30's "balance is a sum,
// not a number in a cell" principle, exercised for real.
func printBalances(bc *core.Blockchain, w1, w2 *wallet.Wallet) {
	fmt.Println("--- Balances (recomputed from scratch by scanning the chain) ---")
	fmt.Printf("Wallet 1: %d gochips\n", bc.BalanceOf(w1.Address()))
	fmt.Printf("Wallet 2: %d gochips\n", bc.BalanceOf(w2.Address()))
	fmt.Println()
}

// findPrevTXs looks up the full earlier transaction behind every input
// tx wants to spend -- see Section 3.
func findPrevTXs(bc *core.Blockchain, tx *core.Transaction) map[string]core.Transaction {
	prevTXs := make(map[string]core.Transaction)
	for _, in := range tx.Inputs {
		prevTX, err := bc.FindTransaction(in.TxID)
		if err != nil {
			log.Fatalf("findPrevTXs: could not find transaction %x: %v", in.TxID, err)
		}
		prevTXs[hex.EncodeToString(in.TxID)] = *prevTX
	}
	return prevTXs
}
```

Running `go run ./cmd/senddemo` produces output like this (addresses, transaction IDs, and mining time will differ every run, since key generation and proof-of-work search are both genuinely random):

```
Wallet 1 address: 1Fh2K9pXqLmN8vRtWzYcBaJd3Ee7GsUo1x
Wallet 2 address: 1Pz7VbNcRqLtXmK2WdEo9YhJf4Uu8AaGs3

Mined block 1 (coinbase reward 50 gochips to Wallet 1)
--- Balances (recomputed from scratch by scanning the chain) ---
Wallet 1: 50 gochips
Wallet 2: 0 gochips

Wallet 1 signed a transaction sending 20 gochips to Wallet 2 (tx 9f86d081884c7d659a2feaa0c55ad015...)

Mined block 2 (2 transactions: 1 coinbase reward + 1 transfer)
--- Balances (recomputed from scratch by scanning the chain) ---
Wallet 1: 80 gochips
Wallet 2: 20 gochips
```

---

## 6. Checking the Arithmetic

It's worth pausing on the final numbers, because they look almost too simple to be interesting — and that simplicity is exactly the point; it means everything underneath worked correctly. Trace where every gochip in the final balances actually came from:

```
Wallet 1's 80 gochips breaks down as:

  50  -- the ORIGINAL coinbase reward from block 1
- 20  -- sent to Wallet 2 in the transaction
+ 30  -- change from that same transaction, returned to Wallet 1
         (its one 50-gochip UTXO had to be spent in FULL -- Chapter 30,
          Section 4 -- so the leftover 30 comes back as change)
+ 50  -- the SECOND, fresh coinbase reward from mining block 2
------
  80

Wallet 2's 20 gochips is simply the one payment output it received.

Total gochips now in existence: 80 + 20 = 100
Total coinbase rewards minted so far: 50 (block 1) + 50 (block 2) = 100

100 == 100 -- every gochip in the final balances traces back to a real
coinbase reward. No value was created or destroyed by the transfer
itself (Chapter 30's conservation principle) -- ALL new value came
only from mining, never from an ordinary transaction.
```

This check is worth running yourself, by hand, on your own program's actual output every time you run the demo — it's a fast, reliable way to catch a subtle bug (a fee accidentally charged or not charged, a double-counted UTXO, a change output computed incorrectly) that a program simply printing "no errors" would never surface on its own.

---

## 7. What This Demo Does Not Yet Do

Being precise about scope, in the same spirit as Chapter 36, Section 9:

- **Everything runs in one process, in memory.** There is no persistence (Volume 8) and no networking (Volume 7) — Wallet 1 and Wallet 2 are not on two different computers, and nothing here has been broadcast anywhere. This demo proves the *transaction and mining pipeline* works end to end; it does not yet prove it works *across* independent nodes.
- **No mempool in this specific demo.** `bc.MineBlock` is called directly with an explicit transaction list, rather than pulling from a `core.Mempool` the way Chapter 34's design and Chapter 36's CLI both do. This is a deliberate simplification to keep this chapter's demo linear and easy to trace by hand; Exercise 5 asks you to rebuild it using a real mempool instead.
- **No fee.** The transaction in this demo pays no fee (Chapter 35) — Wallet 1 sends exactly 20 and gets exactly 30 back in change, with nothing left over for the miner beyond the coinbase reward itself. A real, congested network would need fee-based prioritization for a miner to have any reason to prefer one pending transaction over another.
- **Balances are correct, but slow to compute.** Every `BalanceOf` call here rescans the entire chain (Chapter 32's `UTXOSet`, unchanged). On a two-block demo chain this is instant; Chapter 56 is where this stops being fast enough to ignore.

---

## Summary

- This chapter's demo builds nothing new — it wires together `wallet.New`, `core.NewBlockchain`, `core.NewCoinbaseTX`, `core.MineBlock`, `core.NewTransaction`, `Transaction.Sign`/`Verify`, and `Blockchain.BalanceOf` into one continuous, runnable program.
- A **coinbase transaction** has no inputs because it does not *move* existing value — it *creates* brand-new value, by protocol rule, as the reward for successfully mining a block. Every ordinary transaction, by contrast, can only ever move value that already exists somewhere as a real UTXO.
- Every block mined earns a fresh coinbase reward, not just the very first one — GoChain's demo mines two blocks and pays out two separate 50-gochip rewards, both to Wallet 1 in this particular run.
- `Blockchain.FindTransaction` is a new helper this chapter adds: a linear scan over every block's transactions, used to build the `prevTXs` map `Transaction.Sign` needs to know which `PubKeyHash` each input is unlocking.
- The demo mines a reward block, checks balances, builds and signs a real transfer, mines a second block containing both a new reward and the transfer, and checks balances again — every single balance check recomputed from scratch by scanning the whole chain, never cached.
- The final numbers (Wallet 1: 80, Wallet 2: 20) can be fully reconstructed by hand from the two coinbase rewards and the one transfer's change output, and their total (100) exactly matches total coinbase rewards minted (100) — a concrete demonstration of Chapter 30's conservation-of-value principle holding up under real code, not just worked-example arithmetic.
- This demo intentionally has no persistence, no networking, no mempool, and no fees — each is a deliberate simplification this specific chapter leaves for later volumes or for the exercises below.

---

## Exercises

### Easy

1. Run the demo program yourself several times in a row. Which parts of the printed output change between runs, and which stay exactly the same (`50`, `80`, `20`)? Explain why each changing value is expected to differ and each fixed value is expected to stay the same.
2. Modify `main.go` to also mine a *third* block containing only a fresh coinbase reward (no transfer), and update `printBalances` calls so you can confirm Wallet 1's balance increases by exactly `miningReward` again, with Wallet 2 unaffected.
3. In your own words, using this chapter's exact numbers, explain why Wallet 1 ends up with 80 gochips rather than the more "obvious"-sounding 30 (its change) — what other source of value is easy to forget it also received in this same demo?

### Medium

4. The demo currently has Wallet 1 mine every block (and therefore collect every coinbase reward). Modify it so that Wallet 2 mines block 2 instead (i.e., the second `NewCoinbaseTX` call pays `wallet2.Address()`), and predict, before running it, what the new final balances for both wallets should be. Verify your prediction against the program's actual output.
5. Rebuild this demo using a real `core.Mempool` (Chapter 34) instead of passing transaction slices to `MineBlock` directly: create a mempool, add the signed transfer transaction to it with `mempool.Add`, retrieve pending transactions with `mempool.GetPending()` for the second `MineBlock` call, and call `mempool.Remove` afterward — matching the pattern Chapter 34, Section 8 described.
6. Add a second transfer to the demo: after block 2 is mined, have Wallet 2 send 5 of its newly-received 20 gochips back to Wallet 1, mine a third block containing that transaction, and verify the final balances by hand before checking them against the program's output.

### Hard

7. Deliberately introduce a bug: after building and signing `tx` in Step 5, but before mining block 2, mutate `tx.Outputs[0].Value` to a different number (simulating tampering, exactly like Chapter 33, Section 8's walkthrough). Add a check in `main.go` that calls `tx.Verify(prevTXs)` again right before mining and refuses to proceed (with a clear error) if verification now fails. Confirm your check actually catches the tampering.
8. This demo builds `prevTXs` by scanning the whole chain once per input, via `FindTransaction`, itself a linear scan over every block and every transaction. For a chain with many blocks and transactions, estimate (with rough big-O reasoning) how expensive building `prevTXs` for a transaction with several inputs referencing outputs scattered across many different earlier blocks would become, and describe what kind of index (without fully implementing it — that's Volume 8's job) would make this fast.
9. Extend the demo into a small three-wallet scenario: Wallet 1 receives a coinbase reward, sends part of it to Wallet 2, and Wallet 2 sends part of *that* to Wallet 3, all within a chain of three mined blocks. Verify, by hand and then by running the program, that total value across all three wallets' final balances still exactly equals the total of every coinbase reward minted across all three blocks — generalizing Section 6's arithmetic check to a longer chain of transfers.
