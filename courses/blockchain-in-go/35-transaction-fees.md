# Chapter 35: Transaction Fees

A block has limited space, but the mempool can hold far more transactions than fit in it. This chapter gives miners a rational way to choose which pending transactions win a place in the next block: fees, paid voluntarily by whoever wants their transaction prioritized.

## Table of Contents

1. [Why Miners Need an Incentive Beyond the Block Reward](#1-why-miners-need-an-incentive-beyond-the-block-reward)
2. [What a Fee Actually Is](#2-what-a-fee-actually-is)
3. [Adding Fee Calculation to core.Transaction](#3-adding-fee-calculation-to-coretransaction)
4. [Fee-Per-Byte: Why Raw Fee Isn't Enough](#4-fee-per-byte-why-raw-fee-isnt-enough)
5. [Selecting Transactions for a Block](#5-selecting-transactions-for-a-block)
6. [Wiring Selection in Front of MineBlock](#6-wiring-selection-in-front-of-mineblock)
7. [Worked Example: Five Transactions, One Small Block](#7-worked-example-five-transactions-one-small-block)
8. [Fee Markets Under Congestion](#8-fee-markets-under-congestion)
9. [Testing Fee Calculation and Selection](#9-testing-fee-calculation-and-selection)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Miners Need an Incentive Beyond the Block Reward

By Volume 4, mining already pays: every block a miner successfully mines includes a coinbase transaction (introduced properly in Chapter 37) that rewards them with newly created gochips. That reward alone is enough incentive to mine *some* block. But it says nothing about *which transactions* a miner should bother including in that block.

Imagine a busy post office with one truck leaving every hour, and a hundred packages waiting to go out, but the truck only has room for thirty. If every sender pays the exact same flat rate regardless of urgency, the post office has no rational way to decide whose thirty packages go out this hour instead of next. Real post offices solve this with priority/express shipping: pay more, and your package is more likely to make the next truck. A **transaction fee** is GoChain's version of that express-shipping surcharge — a little extra value a sender attaches to their transaction, which is not "sent" to anyone in particular, but instead becomes available to whichever miner successfully mines the block containing it.

This matters more and more as a chain matures. Early on, block rewards alone are generous enough that miners will happily include almost any pending transaction. But block rewards are typically designed to shrink over time (Bitcoin halves its reward periodically), and eventually most of a miner's income has to come from fees instead. Building fee-based prioritization now means GoChain behaves like a real, economically sustainable blockchain from the start, not one that quietly assumes free, unlimited space forever.

```
                         one truck, limited capacity
                                     |
   100 packages waiting  ---------> [ truck: room for 30 ] ---> delivered this hour
   (all paid the same                       ^
    flat postage)                            |
                                    which 30? no rational
                                    way to choose without
                                    an extra signal (a fee)
```

---

## 2. What a Fee Actually Is

Recall the UTXO model from Chapter 30: a transaction consumes some inputs (existing UTXOs) and produces some outputs (new UTXOs). Nothing requires that the total value going in exactly equals the total value coming out. A **fee** is defined as simply the difference:

```
fee = (sum of all input values) - (sum of all output values)
```

If Alice's transaction consumes a 50-gochip UTXO as its only input, and produces one output paying Bob 47 gochips, the fee is 3 gochips — value that existed, was legitimately claimed by the transaction, but was not assigned to any output. That "missing" 3 gochips doesn't vanish; it becomes available to whichever miner mines the block this transaction ends up in, folded into that miner's own coinbase reward (you'll see this combined in the worked example in Chapter 37).

Nobody explicitly writes "pay a fee of 3" anywhere in the transaction — a fee is never an output with its own `TxOutput`. It is entirely implicit: whatever value goes in and isn't accounted for in the outputs simply becomes the fee. This is why computing a transaction's fee requires knowing the *value* of every input, which is not stored directly on the `TxInput` itself.

```
Alice's transaction

  Input:  spends UTXO worth 50 gochips
                    |
                    v
  Output: pays Bob 47 gochips
                    |
                    v
  50 (in) - 47 (out) = 3 gochips fee
  (not written anywhere explicitly — it's the leftover)
```

---

## 3. Adding Fee Calculation to core.Transaction

A `TxInput` only stores a *reference* to a previous output — its `TxID` and `OutIndex` — not the value of that output. To compute a fee, we need to look up what each input is actually worth, by finding the original transaction it references. This is precisely the same `prevTXs` lookup pattern Chapter 33 built for `Sign` and `Verify`: a map from a transaction's hex ID to the transaction itself.

```go
package core

import "encoding/hex"

// IsCoinbase reports whether tx is a coinbase transaction — the special,
// no-input transaction that creates new supply, covered fully in Chapter
// 37. A coinbase transaction has exactly one input whose TxID is empty.
func (tx *Transaction) IsCoinbase() bool {
	return len(tx.Inputs) == 1 && len(tx.Inputs[0].TxID) == 0
}

// Fee returns tx's transaction fee: total input value minus total output
// value. prevTXs must map the hex-encoded ID of every transaction
// referenced by tx's inputs to that transaction itself, so we can look up
// how much each input is actually worth.
func (tx *Transaction) Fee(prevTXs map[string]Transaction) int64 {
	if tx.IsCoinbase() {
		return 0 // coinbase transactions create new supply; there's no "input value" to compare against
	}

	var inputTotal, outputTotal int64
	for _, in := range tx.Inputs {
		prevTX, ok := prevTXs[hex.EncodeToString(in.TxID)]
		if !ok {
			// A verified transaction should never hit this — Verify()
			// already required every referenced previous transaction to
			// be resolvable. We skip defensively rather than panic.
			continue
		}
		inputTotal += prevTX.Outputs[in.OutIndex].Value
	}

	for _, out := range tx.Outputs {
		outputTotal += out.Value
	}

	return inputTotal - outputTotal
}
```

`IsCoinbase` checks the one telltale sign of a coinbase transaction: an empty `TxID` on its single input, since a coinbase transaction has nothing real to point at (Chapter 37 explains exactly why). `Fee` walks every input, looks each one up in `prevTXs` to find out what it's actually worth, sums those values into `inputTotal`, sums every output's `Value` into `outputTotal`, and returns the difference. A coinbase transaction short-circuits to a fee of zero, since "fee" doesn't mean anything for a transaction that isn't spending anyone's existing value in the first place.

---

## 4. Fee-Per-Byte: Why Raw Fee Isn't Enough

Suppose Transaction A pays a 100-gochip fee, and Transaction B pays a 10-gochip fee. Transaction A looks like the obvious winner — until you learn Transaction A is enormous (say, it consolidates two thousand tiny inputs into one output) and takes up 500 times more space in the block than Transaction B. A miner filling a block of *limited byte size* doesn't just want the highest absolute fee — it wants the highest fee *per byte of space consumed*, exactly the way a moving company charges by weight and volume, not by "how much the customer offered to pay in total," because space on the truck is the actual scarce resource.

We need a transaction's serialized size in bytes to compute this. Every GoChain type has followed a `Serialize()` convention since Chapter 09 — `Transaction` is no exception:

```go
// Size returns the number of bytes tx occupies once serialized, using
// the same Serialize() method used for hashing and disk storage.
func (tx *Transaction) Size() int {
	return len(tx.Serialize())
}

// FeePerByte returns tx's fee divided by its serialized size — the
// metric miners actually care about when block space, not just fee
// amount, is the scarce resource.
func (tx *Transaction) FeePerByte(prevTXs map[string]Transaction) float64 {
	size := tx.Size()
	if size == 0 {
		return 0
	}
	return float64(tx.Fee(prevTXs)) / float64(size)
}
```

`Size` reuses `Serialize()` (already built for hashing transactions and for the disk format in Chapter 20) purely to measure how many bytes the transaction takes up. `FeePerByte` divides the absolute fee from Section 3 by that size, giving a rate — gochips per byte — that is directly comparable across transactions of wildly different sizes. This is the number GoChain's block-filling logic actually sorts by, not the raw fee.

If you've ever used a real Bitcoin wallet, this is precisely the idea behind the "sat/vByte" number many of them show you before you send — satoshis (Bitcoin's smallest unit) per virtual byte of transaction size. Wallets display that rate, not a flat total, for exactly the reason worked out above: byte-for-byte cost is what a miner actually compares across competing transactions, so it's the only fee number that behaves consistently regardless of how big or small your particular transaction happens to be. GoChain's `FeePerByte` is the same concept under a different unit name (gochips instead of satoshis).

---

## 5. Selecting Transactions for a Block

With `Fee` and `FeePerByte` in place, we can build the function that chooses which pending transactions make it into the next block. This lives alongside the other `core` package code, and it deliberately does **not** change the signature of `MineBlock` — it produces the `[]*Transaction` slice that gets *passed into* `MineBlock`, unchanged from Volume 4:

```go
package core

import "sort"

// MaxBlockBytes caps how much transaction data a single block may
// contain. A real chain would derive this from historical block size
// analysis and network bandwidth constraints; GoChain fixes it as a
// constant for clarity, the same way Bitcoin originally shipped with a
// simple fixed 1MB cap before later soft forks changed the accounting.
const MaxBlockBytes = 1_000_000

// SelectForBlock chooses which pending transactions to include in the
// next block. It sorts candidates by fee-per-byte, highest first, then
// greedily packs them into the block until no more fit within maxBytes.
// prevTXs must resolve every transaction referenced by any candidate's
// inputs, exactly like the map Sign/Verify and Fee use.
func SelectForBlock(pending []*Transaction, prevTXs map[string]Transaction, maxBytes int) []*Transaction {
	type candidate struct {
		tx         *Transaction
		size       int
		feePerByte float64
	}

	candidates := make([]candidate, 0, len(pending))
	for _, tx := range pending {
		candidates = append(candidates, candidate{
			tx:         tx,
			size:       tx.Size(),
			feePerByte: tx.FeePerByte(prevTXs),
		})
	}

	// Highest fee-per-byte first: this ordering IS the fee market. A
	// transaction offering more per byte of space consumed is more
	// likely to be picked when space is scarce.
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].feePerByte > candidates[j].feePerByte
	})

	selected := make([]*Transaction, 0, len(candidates))
	usedBytes := 0
	for _, c := range candidates {
		if usedBytes+c.size > maxBytes {
			// Doesn't fit right now — skip it, but keep checking the
			// rest, since a smaller, lower-fee transaction further down
			// the list might still fit in the remaining space.
			continue
		}
		selected = append(selected, c.tx)
		usedBytes += c.size
	}

	return selected
}
```

This is a **greedy** algorithm, not an optimal one — worth naming plainly, since it matters for Section 8. It sorts every candidate by `feePerByte` descending, then walks the sorted list once, adding each transaction to `selected` if it still fits in the remaining space, and skipping it (without giving up) if it doesn't. This means a big, no-longer-fitting transaction near the top of the list can be skipped in favor of smaller ones further down — but it does **not** guarantee the mathematically best possible combination of transactions for maximizing total fees collected, which is a much harder combinatorial problem (a variant of the classic "knapsack problem"). Real production miners, including Bitcoin's own reference implementation, use greedy heuristics very similar to this one for exactly this reason: it's fast, simple, and good enough in practice, even though it isn't provably optimal.

---

## 6. Wiring Selection in Front of MineBlock

`core.Blockchain.MineBlock(transactions []*Transaction) *Block` is unchanged from Volume 4 — what changes is what a miner passes into it. Here is the full pattern a node's mining loop follows, tying together the mempool from Chapter 34 and the selection logic from Section 5:

```go
func mineNextBlock(bc *core.Blockchain, mempool *core.Mempool, prevTXs map[string]core.Transaction, minerAddress string) *core.Block {
	pending := mempool.GetPending()

	// Choose which pending transactions actually make it into this block.
	chosen := core.SelectForBlock(pending, prevTXs, core.MaxBlockBytes)

	// Sum up the fees collected from every chosen transaction — this
	// gets folded into the miner's own coinbase reward (Chapter 37).
	var totalFees int64
	for _, tx := range chosen {
		totalFees += tx.Fee(prevTXs)
	}

	const blockReward = 50
	coinbase := core.NewCoinbaseTX(minerAddress, blockReward+totalFees)

	// The coinbase transaction always goes first; MineBlock's signature
	// and behavior are exactly as built in Volume 4.
	blockTxs := append([]*core.Transaction{coinbase}, chosen...)
	block := bc.MineBlock(blockTxs)

	// Only the transactions that actually made it into the block are
	// removed from the mempool — anything skipped for space stays
	// pending, waiting for the next block.
	for _, tx := range chosen {
		mempool.Remove(tx.ID)
	}

	return block
}
```

`mineNextBlock` is the glue: it pulls everything currently pending out of the mempool, asks `SelectForBlock` which of those transactions actually fit and are worth including, tallies up the fees they collectively pay, and folds that total into the miner's coinbase reward alongside the fixed `blockReward`. It then calls `bc.MineBlock` with the coinbase transaction first, followed by the chosen transactions — `MineBlock`'s own signature and mining logic haven't changed at all from Volume 4. Finally, and importantly, it only removes the transactions that were actually **chosen** from the mempool — anything `SelectForBlock` skipped because it didn't fit stays right where it was, waiting for a future block with more room.

Notice what happens at the two extremes. If the mempool is completely empty, `pending` is an empty slice, `chosen` stays empty, `totalFees` is zero, and the miner still mines a valid block containing only the coinbase transaction paying exactly `blockReward` — mining never grinds to a halt just because nobody has anything to send. At the other extreme, if the mempool has far more pending transactions (in bytes) than `MaxBlockBytes` can hold, `chosen` only ever contains whatever `SelectForBlock` decided fits, and the rest patiently wait for the next block — this is precisely the congested scenario Section 8 explores in detail.

---

## 7. Worked Example: Five Transactions, One Small Block

Let's make the greedy algorithm concrete with real numbers. Suppose five transactions are sitting in the mempool, and — to make the example easy to follow by hand — our block has a deliberately tiny `maxBytes` of 300:

| Transaction | Size (bytes) | Fee (gochips) | Fee-per-byte |
|---|---|---|---|
| Tx1 | 250 | 25 | 0.10 |
| Tx2 | 100 | 50 | 0.50 |
| Tx3 | 150 | 90 | 0.60 |
| Tx4 | 200 | 20 | 0.10 |
| Tx5 | 80  | 56 | 0.70 |

`SelectForBlock` first sorts these by fee-per-byte, descending:

```
Tx5 (0.70)  ->  Tx3 (0.60)  ->  Tx2 (0.50)  ->  Tx1 (0.10)  ->  Tx4 (0.10)
```

Then it walks the sorted list, greedily packing what fits into 300 bytes:

```
Step 1: Tx5, size 80.  usedBytes:   0 + 80 = 80   <= 300  -> SELECTED
Step 2: Tx3, size 150. usedBytes:  80 + 150 = 230 <= 300  -> SELECTED
Step 3: Tx2, size 100. usedBytes: 230 + 100 = 330 >  300  -> SKIPPED
Step 4: Tx1, size 250. usedBytes: 230 + 250 = 480 >  300  -> SKIPPED
Step 5: Tx4, size 200. usedBytes: 230 + 200 = 430 >  300  -> SKIPPED
```

The final block contains **Tx5 and Tx3**, using 230 of the available 300 bytes and collecting 56 + 90 = **146 gochips in fees**. Notice something a little unsatisfying: after selecting Tx5 and Tx3, there are 70 bytes of unused space left in the block, and Tx2 — which would collect a healthy 50-gochip fee and is only 100 bytes — doesn't fit in that remaining 70 bytes, so it's skipped even though it's the third-best transaction by fee-per-byte. This is exactly the greedy algorithm's known limitation from Section 5: a more exhaustive search might find that some *other* combination (say, dropping Tx3 in favor of Tx2 plus something smaller) collects more total fees or wastes less space. GoChain accepts this trade-off for simplicity — Exercise 8 asks you to measure how much fee revenue is actually left on the table by this greedy shortcut.

---

## 8. Fee Markets Under Congestion

A **fee market** is what emerges the moment more transactions want space in a block than can actually fit. When the mempool is nearly empty relative to block capacity, almost every transaction gets in regardless of its fee — there's no real competition for space, so a fee of zero or near-zero would still likely make the next block. But once the mempool consistently holds *more* transactions (measured in bytes) than a typical block's `MaxBlockBytes`, transactions start genuinely competing, and fee-per-byte becomes the deciding factor, exactly as in Section 7's worked example.

```
Quiet mempool (plenty of spare room)          Congested mempool (space is scarce)

+----------------------------------+          +----------------------------------+
| Block capacity: 300 bytes         |          | Block capacity: 300 bytes         |
|------------------------------------|         |------------------------------------|
| [Tx A, low fee ] [Tx B, low fee]   |         | [Tx5 ][Tx3 ][........unused........]|
| [........... unused ..............]|         | (Tx2, Tx1, Tx4 all missed the cut) |
+----------------------------------+          +----------------------------------+
Every waiting transaction fits;               Only the highest fee-per-byte
fee level barely matters.                     transactions get in; everyone
                                               else waits, or raises their fee.
```

This has real, observable consequences that show up in every production blockchain:

- **Fees rise during congestion.** If you want your transaction to beat out everyone else's during a busy period, you have to outbid them — offering a higher fee-per-byte than whatever the current "cutoff" happens to be for that block.
- **Fees fall during quiet periods.** With plenty of spare block space, even a low fee-per-byte transaction gets included, because there's no competition forcing miners to be selective.
- **Fee estimation becomes a real product problem.** Real wallets (and GoChain's own wallet, starting in Chapter 36) often look at recent blocks' fee-per-byte distributions to suggest a fee likely to get a new transaction included within some number of blocks — guessing wrong in either direction either overpays or risks the transaction languishing in the mempool indefinitely.
- **A transaction can get "stuck."** If a transaction's fee-per-byte never beats the going rate, it can sit in the mempool for a very long time, since `SelectForBlock` will simply keep skipping it in favor of anything better-paying that shows up. Real wallets sometimes offer "replace-by-fee" — rebroadcasting the same transaction with a higher fee to bump its priority — a topic outside this chapter's scope but worth knowing exists.
- **Block size limits are themselves a policy choice with trade-offs.** A larger `MaxBlockBytes` fits more transactions per block (lower fees, less congestion) but makes each block larger to download, verify, and store, which affects how easily ordinary computers can keep up as a full node — a genuine, contested trade-off in real blockchain communities (Bitcoin's block-size debates are the most famous real-world example).
- **Tiny fees create their own problem, called "dust."** A transaction offering a fee so small it isn't worth a miner's time to include is functionally the same as an unfunded package sitting in the warehouse forever — it clutters the mempool without any realistic path to being mined. Real nodes often reject transactions below some configurable "minimum relay fee" outright, rather than let them accumulate indefinitely; `SelectForBlock`'s optional minimum threshold (Exercise 5) implements exactly this idea for GoChain.

---

## 9. Testing Fee Calculation and Selection

A few focused tests cover the logic built in this chapter:

```go
func TestTransaction_Fee(t *testing.T) {
	prevTX := Transaction{
		ID:      []byte("prev-tx"),
		Outputs: []TxOutput{{Value: 50, PubKeyHash: []byte("alice-hash")}},
	}
	prevTXs := map[string]Transaction{
		hex.EncodeToString(prevTX.ID): prevTX,
	}

	tx := &Transaction{
		ID:      []byte("spending-tx"),
		Inputs:  []TxInput{{TxID: prevTX.ID, OutIndex: 0}},
		Outputs: []TxOutput{{Value: 47, PubKeyHash: []byte("bob-hash")}},
	}

	if got := tx.Fee(prevTXs); got != 3 {
		t.Fatalf("expected fee of 3, got %d", got)
	}
}

func TestTransaction_CoinbaseHasZeroFee(t *testing.T) {
	coinbase := NewCoinbaseTX("miner-address", 50)
	if got := coinbase.Fee(nil); got != 0 {
		t.Fatalf("expected coinbase fee of 0, got %d", got)
	}
}

// fakeTx builds a minimal transaction shaped only to hit a target
// serialized size and a target fee, so we can reproduce Section 7's
// worked example numbers directly, without needing a real prior
// transaction, wallet, or signature for this particular test.
func fakeTx(id string, targetSize int, fee int64) (*Transaction, Transaction) {
	prev := Transaction{
		ID:      []byte("prev-" + id),
		Outputs: []TxOutput{{Value: fee + 1000, PubKeyHash: []byte("owner")}},
	}
	tx := &Transaction{
		ID:     []byte(id),
		Inputs: []TxInput{{TxID: prev.ID, OutIndex: 0}},
		Outputs: []TxOutput{
			{Value: 1000, PubKeyHash: []byte("recipient")},
			// Padding bytes below pad the serialized size up toward
			// targetSize; a real transaction's size comes from its real
			// field contents, but for this test we only care about
			// reproducing specific byte counts deterministically.
			{Value: 0, PubKeyHash: make([]byte, targetSize)},
		},
	}
	return tx, prev
}

func TestSelectForBlock_PrefersHigherFeePerByte(t *testing.T) {
	sizes := map[string]int{"tx1": 250, "tx2": 100, "tx3": 150, "tx4": 200, "tx5": 80}
	fees := map[string]int64{"tx1": 25, "tx2": 50, "tx3": 90, "tx4": 20, "tx5": 56}

	var pending []*Transaction
	prevTXs := make(map[string]Transaction)
	for _, id := range []string{"tx1", "tx2", "tx3", "tx4", "tx5"} {
		tx, prev := fakeTx(id, sizes[id], fees[id])
		pending = append(pending, tx)
		prevTXs[hex.EncodeToString(prev.ID)] = prev
	}

	selected := SelectForBlock(pending, prevTXs, 300)

	if len(selected) != 2 {
		t.Fatalf("expected exactly 2 transactions selected, got %d", len(selected))
	}

	selectedIDs := map[string]bool{}
	for _, tx := range selected {
		selectedIDs[string(tx.ID)] = true
	}
	if !selectedIDs["tx5"] || !selectedIDs["tx3"] {
		t.Fatalf("expected tx5 and tx3 to be selected, got: %v", selectedIDs)
	}
}
```

`TestTransaction_Fee` confirms the basic arithmetic from Section 3: a 50-value input against a 47-value output yields a fee of 3. `TestTransaction_CoinbaseHasZeroFee` confirms the special case from `IsCoinbase` — a coinbase transaction never contributes a "fee" in the usual sense. `fakeTx` is a small test helper that builds a transaction shaped to hit a specific serialized size and fee, purely so the test can reproduce Section 7's exact worked-example numbers without needing real wallets or signatures. `TestSelectForBlock_PrefersHigherFeePerByte` uses it to confirm the greedy selection logic really does reproduce the worked example: tx5 and tx3 win, and the lower fee-per-byte transactions are correctly left out of the final selected set.

---

## Summary

- A **fee** is the leftover value in a transaction — total input value minus total output value — and is never written explicitly as its own output.
- `Transaction.Fee(prevTXs)` computes this by looking up each input's real value in a `prevTXs` map, the same lookup pattern `Sign`/`Verify` already use; coinbase transactions always report a fee of zero.
- Raw fee alone is misleading when transactions vary wildly in size — `FeePerByte` (fee divided by `Size()`, itself built on the existing `Serialize()` convention) is the metric that actually matters when block space is the scarce resource.
- `SelectForBlock` sorts pending transactions by fee-per-byte descending and greedily packs them into a block up to `MaxBlockBytes`, without changing `MineBlock`'s signature — it only decides what gets passed into it.
- The greedy packing approach is fast and simple but not provably optimal; it can leave a lower-fee transaction that would have fit skipped over unused block space, purely due to ordering.
- A **fee market** emerges once pending transaction volume regularly exceeds block capacity: fees rise under congestion, fall in quiet periods, and low-fee transactions can get stuck waiting indefinitely.
- Only the transactions actually chosen for a block are removed from the mempool — anything skipped for lack of space remains pending for a future block.

---

## Exercises

### Easy

1. Write a test where a transaction's inputs total exactly equal to its outputs total, and confirm `Fee()` returns 0.
2. Add a method `func (tx *Transaction) HasPositiveFee(prevTXs map[string]Transaction) bool` and write a table-driven test covering a positive fee, a zero fee, and (as an explicit edge case worth reasoning about) whether a *negative* fee should even be possible for a validly signed, `Verify()`-passing transaction. Explain your answer in a code comment.
3. Change `MaxBlockBytes` to a much smaller constant (say, 150) in a test and rerun the Section 7 worked example numbers through `SelectForBlock` by hand first, then verify your prediction against the actual function output.

### Medium

4. Complete `TestSelectForBlock_PrefersHigherFeePerByte` from Section 9 with the exact five transactions and byte/fee numbers from Section 7's table, asserting the selected set is precisely `{Tx5, Tx3}` and the total fee collected equals 146.
5. Modify `SelectForBlock` to accept an additional parameter, `minFeePerByte float64`, and skip any candidate below that threshold entirely, regardless of whether it would otherwise fit — simulating a miner that refuses to bother with transactions below a minimum acceptable rate. Write a test proving a transaction that would otherwise fit gets excluded once its fee-per-byte falls below the threshold.
6. Write a benchmark comparing `SelectForBlock` against a naive "select transactions in mempool iteration order until full, ignoring fees" version, using a synthetic mempool of 1,000 transactions with randomized fees and sizes. Report which one collects more total fees on average across several runs, and explain in a short comment why this demonstrates the fee market actually working as intended.

### Hard

7. Implement a second selection function, `SelectForBlockOptimal`, that uses dynamic programming (the classic 0/1 knapsack algorithm) to find the mathematically optimal set of transactions maximizing total fees within `maxBytes`, for a mempool of realistic size (a few hundred transactions with byte sizes small enough for DP to be tractable). Compare its result against `SelectForBlock`'s greedy result on the Section 7 worked example and on at least one adversarial case you construct where greedy provably leaves fee revenue on the table.
8. Using the adversarial case you constructed in Exercise 7, compute exactly how many gochips of fee revenue the greedy `SelectForBlock` leaves unclaimed compared to `SelectForBlockOptimal`, and write a short analysis (150-250 words) of whether this gap would matter in practice for a real chain, considering that `SelectForBlockOptimal`'s time complexity grows with `maxBytes` in a way `SelectForBlock`'s does not.
9. Design and implement a simple fee-estimation function, `EstimateFee(recentBlocks []*Block, targetConfirmationBlocks int) float64`, that looks at the fee-per-byte of transactions actually included in some number of recent blocks and returns a suggested fee-per-byte likely to get a new transaction included within `targetConfirmationBlocks` blocks. Write a test using a synthetic sequence of blocks with known fee-per-byte distributions and assert your function's suggestion falls within a reasonable range.
