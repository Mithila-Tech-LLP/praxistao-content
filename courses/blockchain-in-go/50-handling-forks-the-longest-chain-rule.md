# Chapter 50: Handling Forks — The Longest Chain Rule

Chapter 48's `HandleBlock` gave up the moment a new block didn't cleanly attach to the current chain tip, and Chapter 49's `SyncWithPeer` assumed there was one single, agreed-upon history to catch up on. Neither assumption survives contact with a real network: two miners can solve a valid block at nearly the same instant, and for a little while, the network genuinely disagrees about which one is "next." This chapter teaches GoChain to resolve that disagreement automatically, in favor of whichever chain represents the most total work.

## Table of Contents

1. [How a Fork Happens](#1-how-a-fork-happens)
2. [Why "Longest" Really Means "Heaviest"](#2-why-longest-really-means-heaviest)
3. [Implementing ChainWork](#3-implementing-chainwork)
4. [Detecting a Fork in HandleBlock](#4-detecting-a-fork-in-handleblock)
5. [Implementing ReplaceChain](#5-implementing-replacechain)
6. [What Happens to the Losing Chain's Transactions](#6-what-happens-to-the-losing-chains-transactions)
7. [Hands-On Simulation: Two Nodes, One Fork, One Winner](#7-hands-on-simulation-two-nodes-one-fork-one-winner)
8. [Deep Reorgs and Why They're Dangerous](#8-deep-reorgs-and-why-theyre-dangerous)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. How a Fork Happens

Proof of work (Volume 4) is a race with no starting gun and no referee — every miner is independently searching for a nonce that produces a hash meeting the difficulty target, and nothing stops two different miners, on opposite sides of a large network, from both finding a valid solution within moments of each other, each building on the same previous block.

```
                         Block 10 (agreed by everyone)
                                |
                +---------------+---------------+
                |                               |
          Block 11-A                      Block 11-B
       (mined by Node X)               (mined by Node Y)
       broadcast to nearby              broadcast to nearby
       peers first                      peers first
```

For a short time, roughly half the network might see Block 11-A first and consider it the tip, while the other half sees Block 11-B first. Both blocks are individually, completely valid — correct proof of work, correct linkage to Block 10, correctly signed transactions inside. There is nothing *wrong* with either one. The network simply, temporarily, disagrees about which one comes next. This situation — two or more valid blocks competing to extend the same chain — is called a **fork**.

A fork is not a bug, and not evidence of an attack (though Chapter 76, later in the course, covers how an attacker can deliberately try to *cause* one for their own benefit). It's an expected, ordinary consequence of many independent miners racing without coordination. What matters is that the network needs an automatic, deterministic rule for resolving it — one every honest node applies the same way, so they all converge back onto a single agreed history without anyone needing to intervene by hand.

---

## 2. Why "Longest" Really Means "Heaviest"

The classic name for this rule is the **longest chain rule**: when faced with two competing chains, adopt whichever one is longer. That name is a useful first mental model, but it is slightly imprecise, and this chapter implements the more precise version real systems actually use.

Chain "length" by itself is easy for an attacker to fake if it just means block *count* — mining 50 quick, low-difficulty blocks would out-count 40 blocks that each required enormous work, even though the 40-block chain represents far more real computational effort and is far harder to have produced dishonestly. What actually matters is **total accumulated proof of work**: sum up how much computational effort went into producing every block in a candidate chain, and prefer whichever chain required the most. In the overwhelmingly common case, more work also just means more blocks (since difficulty doesn't change wildly between neighboring blocks), so "longest chain" is usually a fine shorthand — but "heaviest chain" is the actual, precise rule, and the one we implement as `ChainWork()`.

Recall from Chapter 24 and Chapter 25 that GoChain's proof-of-work target is expressed as a required number of leading zero bits in a block's hash, tracked on each block as a `Difficulty` field (the number of leading zero bits that block's proof of work had to satisfy). The expected number of hash attempts needed to find such a solution grows exponentially with `Difficulty` — roughly `2^Difficulty` attempts on average — so we use exactly that as each block's individual contribution to the chain's total work.

```
Two competing chains from Block 10 (Difficulty stored per block):

  Chain A:  ... -> Block 10 -> Block 11-A (Difficulty 20) -> Block 12-A (Difficulty 20)
            work(11-A) + work(12-A) = 2^20 + 2^20 = 2,097,152

  Chain B:  ... -> Block 10 -> Block 11-B (Difficulty 20) -> Block 12-B (Difficulty 21)
            work(11-B) + work(12-B) = 2^20 + 2^21 = 3,145,728

  Both chains have exactly 2 blocks (the same "length"), but Chain B
  required roughly 50% more actual computational effort, because Block 12-B
  was mined at a slightly higher difficulty. ChainWork correctly prefers B.
```

---

## 3. Implementing ChainWork

```go
package core

import (
	"math/big"
)

// ChainWork returns the total accumulated proof-of-work across every block
// in this chain, used to compare competing chains during fork resolution.
// Each block's own contribution is approximately 2^Difficulty, since that is
// the expected number of hash attempts required to find a nonce producing a
// hash with that many leading zero bits (Chapter 24). Summing this across
// every block, rather than just counting blocks, is what makes this the
// "heaviest chain" rule rather than a naive "longest chain" rule -- a chain
// with fewer, harder-won blocks can rightfully outweigh one with more, easier
// blocks.
func (bc *Blockchain) ChainWork() *big.Int {
	total := big.NewInt(0)

	for h := uint64(0); h <= bc.Height(); h++ {
		block, err := bc.BlockAtHeight(h)
		if err != nil {
			continue // should not happen on a well-formed chain, but don't panic on it
		}
		// work = 2^Difficulty, computed as a left shift so we never overflow
		// a plain int64 even at high difficulties.
		work := new(big.Int).Lsh(big.NewInt(1), uint(block.Difficulty))
		total.Add(total, work)
	}

	return total
}

// chainWorkOf is a small unexported helper used internally by ReplaceChain
// to compute the work contributed by an arbitrary slice of blocks, without
// requiring them to already be part of this Blockchain.
func chainWorkOf(blocks []*Block) *big.Int {
	total := big.NewInt(0)
	for _, b := range blocks {
		work := new(big.Int).Lsh(big.NewInt(1), uint(b.Difficulty))
		total.Add(total, work)
	}
	return total
}
```

`ChainWork` walks every block currently in the chain, computes `2^Difficulty` for each one using `math/big` (since these numbers get astronomically large at real difficulties — a plain `int64` would overflow long before a genuine network's total work would), and sums them. `chainWorkOf` is the same computation applied to an arbitrary candidate slice of blocks that *isn't* part of our chain yet — exactly what `ReplaceChain` needs when comparing our current chain against a competitor's.

---

## 4. Detecting a Fork in HandleBlock

Chapter 48's `HandleBlock` simply gave up when `AddBlock` failed. Now it can do something smarter: recognize *why* it failed, and if the reason is "this block doesn't extend our tip, but it does look like it might be extending a shorter competing chain," gather that competing chain and let `ReplaceChain` decide whether to switch.

```go
// HandleBlock, updated from Chapter 48 to attempt fork resolution instead
// of simply discarding a block that doesn't extend our current tip.
func (n *Node) HandleBlock(b *core.Block) {
	id := b.HashHex()

	if n.seenBlock.seenOrAdd(id) {
		return
	}

	if err := n.Chain.ValidateBlock(b); err != nil {
		log.Printf("network: rejecting invalid block %s: %v", id, err)
		return
	}

	if err := n.Chain.AddBlock(b); err == nil {
		// Common case: it extended our current tip cleanly. Gossip onward
		// exactly as in Chapter 48.
		n.Broadcast(MsgBlock, b.Serialize())
		return
	}

	// AddBlock failed -- most likely this block's PrevBlockHash points at a
	// block that isn't our current tip. That could mean we're missing
	// blocks entirely (Chapter 49's job), or it could mean this block is
	// the tip of a competing fork we haven't fully seen yet. Try to fetch
	// the rest of that competing chain from whoever sent it to us, then let
	// ReplaceChain decide on the merits (accumulated work), not on whoever
	// happened to reach us first.
	peerAddr := n.peerAddrFor(b)
	competing, err := n.fetchCompetingChain(peerAddr, b)
	if err != nil {
		log.Printf("network: could not investigate possible fork for block %s: %v", id, err)
		return
	}

	if err := n.Chain.ReplaceChain(competing); err != nil {
		// Not heavier than what we already have, or otherwise not a valid
		// replacement -- we simply keep our current chain. This is the
		// normal, expected outcome most of the time a fork briefly appears.
		log.Printf("network: kept current chain over competing fork: %v", err)
		return
	}

	log.Printf("network: reorganized to a heavier competing chain, new tip %s", id)
	n.Broadcast(MsgBlock, b.Serialize())
}
```

`fetchCompetingChain` is a small helper built on the same `SyncWithPeer` machinery from Chapter 49: instead of requesting blocks by height, it walks backward from the received block's `PrevBlockHash`, requesting each ancestor by hash via `MsgGetData` until it reaches a block our own chain already recognizes — that shared ancestor is the **fork point**. Everything from the fork point forward, on the competing side, becomes the `competing []*core.Block` slice handed to `ReplaceChain`.

---

## 5. Implementing ReplaceChain

`ReplaceChain` is where the actual decision and, if warranted, the actual **reorganization** (or "reorg") happens: rolling back to the fork point and replaying the heavier chain's blocks in its place.

```go
import (
	"bytes"
	"errors"
	"fmt"
)

// ReplaceChain considers a competing sequence of blocks -- newBlocks, which
// must connect to some block already in our chain -- and, only if that
// competing chain represents strictly more accumulated work than our
// current one, switches our chain to it. This is the heaviest-chain rule in
// action: we never switch just because a competing chain is merely
// different, or even merely longer in block count, only because it took
// more real proof-of-work effort to produce.
func (bc *Blockchain) ReplaceChain(newBlocks []*core.Block) error {
	bc.mu.Lock()
	defer bc.mu.Unlock()

	if len(newBlocks) == 0 {
		return errors.New("replacechain: empty candidate chain")
	}

	// Step 1: locate the fork point -- the block in our own chain that
	// newBlocks[0] builds directly on top of.
	forkHeight, err := bc.heightOfHash(newBlocks[0].PrevBlockHash)
	if err != nil {
		return fmt.Errorf("replacechain: candidate chain doesn't connect to our history: %w", err)
	}

	// Step 2: validate every block in the candidate chain, in order, as if
	// we were receiving each one fresh. A heavier chain full of garbage is
	// still garbage -- accumulated work alone is not sufficient, every
	// block must also independently be well-formed.
	prevHash := newBlocks[0].PrevBlockHash
	for _, b := range newBlocks {
		if !bytes.Equal(b.PrevBlockHash, prevHash) {
			return errors.New("replacechain: candidate chain is not properly linked")
		}
		if err := bc.ValidateBlock(b); err != nil {
			return fmt.Errorf("replacechain: candidate block %x failed validation: %w", b.Hash, err)
		}
		prevHash = b.Hash
	}

	// Step 3: compare total accumulated work. This is the actual
	// heaviest-chain decision.
	ourTailBlocks, err := bc.blocksFrom(forkHeight + 1)
	if err != nil {
		return fmt.Errorf("replacechain: reading our own chain: %w", err)
	}
	ourTailWork := chainWorkOf(ourTailBlocks)
	candidateWork := chainWorkOf(newBlocks)

	if candidateWork.Cmp(ourTailWork) <= 0 {
		return errors.New("replacechain: competing chain is not heavier than our current one")
	}

	// Step 4: reorg. Roll back every block after the fork point, returning
	// their transactions to the mempool (Section 6), then replay the
	// heavier chain's blocks on top of the fork point, applying their
	// transactions to the UTXO set exactly as AddBlock normally would.
	orphaned := ourTailBlocks
	if err := bc.rollbackTo(forkHeight); err != nil {
		return fmt.Errorf("replacechain: rolling back to fork point: %w", err)
	}
	for _, b := range newBlocks {
		if err := bc.applyBlock(b); err != nil {
			// This should be rare, since we already validated every block
			// above, but if it does happen we are now in a bad spot: we've
			// rolled back and can't cleanly move forward. Surface the error
			// loudly rather than leaving the chain half-updated.
			return fmt.Errorf("replacechain: failed applying candidate block after rollback: %w", err)
		}
	}

	bc.returnOrphanedTransactions(orphaned, newBlocks)
	return nil
}
```

Four steps, each doing one job: **locate** the fork point, **validate** the candidate chain block by block (heaviness alone is not enough — a heavier pile of invalid blocks is still rejected), **compare** total work above the fork point, and only then **reorganize** — rolling the chain back and replaying the new blocks. `rollbackTo` and `applyBlock` are the same underlying primitives `AddBlock` already uses internally to update the UTXO set one block at a time; `ReplaceChain` just runs that machinery backward, then forward again along the new path.

---

## 6. What Happens to the Losing Chain's Transactions

When Chain A's blocks are rolled back in favor of Chain B, every transaction that was inside those now-orphaned blocks doesn't just vanish — it goes back to being unconfirmed. Some of those transactions might also appear in Chain B (if both miners happened to include the same pending transaction), in which case they're now confirmed there instead. The rest need to return to the mempool so they have a chance to be mined again.

```go
// returnOrphanedTransactions puts every transaction from the rolled-back
// blocks back into the mempool as pending, unless that exact transaction
// also appears in the newly adopted chain (in which case it's already
// confirmed there and doesn't need to wait again).
func (bc *Blockchain) returnOrphanedTransactions(orphaned []*core.Block, adopted []*core.Block) {
	confirmedInNewChain := make(map[string]bool)
	for _, b := range adopted {
		for _, tx := range b.Transactions {
			confirmedInNewChain[tx.IDHex()] = true
		}
	}

	for _, b := range orphaned {
		for _, tx := range b.Transactions {
			if tx.IsCoinbase() {
				// A coinbase reward only exists because that specific block
				// was mined -- if the block is gone, the reward is gone
				// with it. There's nothing to "return" to the mempool.
				continue
			}
			if confirmedInNewChain[tx.IDHex()] {
				continue // already confirmed on the winning chain
			}
			// Best effort: the transaction goes back into the waiting room.
			// It might now conflict with something the new chain already
			// confirmed (a double-spend of the same UTXO), in which case
			// Mempool.Add will correctly reject it.
			_ = bc.Mempool.Add(tx)
		}
	}
}
```

Two subtleties are worth calling out by name. The **coinbase transaction** (the special, input-less transaction that pays a miner their block reward, introduced in Chapter 37) that was only "real" because a now-orphaned block was mined simply disappears along with that block — there is no version of "give the miner their reward anyway" once their block has lost the race. And any orphaned transaction that would now double-spend a UTXO the winning chain already consumed is correctly rejected by `Mempool.Add`'s existing double-spend check from Chapter 34 — the fork-resolution code doesn't need any special-case logic for that; it falls naturally out of validation machinery that already exists.

---

## 7. Hands-On Simulation: Two Nodes, One Fork, One Winner

Here is a concrete, step-by-step simulation you can run (or trace by hand) using two in-process `Node`s sharing a fake, controllable network connection instead of real sockets — the same style of test harness suggested in Chapter 48's exercises.

```
Setup: Node X and Node Y are both caught up at height 10, identical tips.
Both have a pending transaction, tx T, in their mempools (received via
gossip earlier).

Step 1: Simulate a network partition -- temporarily disconnect X and Y so
neither can gossip to the other.

Step 2: Node X mines Block 11-X (Difficulty 20), including tx T.
         Node Y, at nearly the same moment, mines Block 11-Y (Difficulty 20),
         also including tx T (both miners had it in their mempool).

    X's chain: ... -> Block 10 -> Block 11-X   (ChainWork += 2^20)
    Y's chain: ... -> Block 10 -> Block 11-Y   (ChainWork += 2^20)

Step 3: Node Y additionally mines Block 12-Y (Difficulty 21) before the
        partition is healed, extending its own lead.

    Y's chain: ... -> Block 10 -> Block 11-Y -> Block 12-Y

Step 4: Restore the network connection between X and Y. Gossip resumes:
        X receives Block 11-Y (via HandleBlock), Y receives Block 11-X.

    - X's HandleBlock(Block 11-Y): AddBlock fails (doesn't extend X's tip,
      which is Block 11-X). fetchCompetingChain retrieves [11-Y, 12-Y] from
      Y. ReplaceChain compares: X's tail work = work(11-X) = 2^20.
      Candidate work = work(11-Y) + work(12-Y) = 2^20 + 2^21 = 3*2^20.
      Candidate is heavier -> X reorganizes onto Y's chain.

    - Y's HandleBlock(Block 11-X): AddBlock fails (Y is already past this
      point, at height 12). Y's own chain is unaffected -- Y already has
      more work than this single block could ever represent alone.

Step 5: Final state -- both X and Y agree: chain tip is Block 12-Y, height
        12. Block 11-X is orphaned. Its coinbase reward is gone. tx T is
        already confirmed in Block 11-Y, so X's returnOrphanedTransactions
        correctly does NOT re-add it to the mempool as pending.

Expected printed output on Node X during Step 4:

    network: kept current chain over competing fork: ...          (not printed here --
      this line only appears when we DON'T switch; contrast with:)
    network: reorganized to a heavier competing chain, new tip <hash of Block 12-Y>
```

This simulation deliberately mirrors the earlier ASCII diagram in Section 1: two nodes, briefly disagreeing, converging automatically the moment they can compare notes — with no human intervention, no coordination beyond the gossip and sync machinery already built in Chapters 48 and 49, and the `ChainWork` comparison from this chapter as the tiebreaker.

---

## 8. Deep Reorgs and Why They're Dangerous

The simulation above resolved a **shallow fork** — just one or two blocks deep, resolved within moments. A **deep reorg**, where a competing chain many blocks longer suddenly displaces a chain you'd been building on for a while, is a much bigger deal, and worth understanding even though this chapter's code handles it mechanically the same way.

The danger is this: if you treated a transaction as "confirmed" the moment it appeared in a block, and then a deep reorg orphans that block, the transaction you thought was final might vanish (or, worse, a conflicting version of it might now be confirmed instead) — this is exactly the mechanism behind a **51% attack**, covered in depth in Chapter 76, where an attacker with enough mining power deliberately builds a longer, secret chain to reverse a transaction they already got real-world value for (for example, they paid for goods with a transaction, waited for the merchant to ship the goods once it looked "confirmed," and only then revealed a heavier, secretly-mined chain that doesn't include their payment at all).

The practical mitigation, used by every real proof-of-work blockchain and worth adopting in GoChain's own wallet and merchant-facing code, is **waiting for confirmations**: instead of treating a transaction as final the instant it's in the latest block, treat it as final once it's buried under several additional blocks (Bitcoin's traditional rule of thumb is 6 confirmations for high-value transactions). The deeper a block is buried, the more accumulated work an attacker would need to out-produce to reorg past it — which is precisely why `ChainWork`, not simple block count, is the right measure of how "safe" a given block really is.

```
      Block 10          Block 11         Block 12         Block 13 (tip)
     [confirmed]       [1 confirm]      [2 confirms]     [3 confirms,
                                                            just mined]

  Reorging past Block 13 requires out-producing very little extra work.
  Reorging past Block 10 requires out-producing FOUR blocks' worth of
  accumulated work -- increasingly impractical the deeper it's buried.
```

---

## Summary

- Two miners solving a valid block at nearly the same time causes a **fork** — a temporary disagreement about which block comes next — which is normal, expected behavior in proof-of-work systems, not a bug or automatically an attack.
- The precise resolution rule is the **heaviest chain rule**: prefer whichever competing chain represents the most total accumulated proof-of-work, not simply whichever has more blocks — "longest chain" is a common but slightly imprecise nickname for the same idea.
- `ChainWork() *big.Int` sums `2^Difficulty` across every block in a chain, using `math/big` to avoid overflow at realistic difficulties.
- `HandleBlock` was extended from Chapter 48 to recognize when a received block doesn't extend the current tip, fetch the rest of that possible competing chain, and hand it to `ReplaceChain` rather than simply discarding it.
- `ReplaceChain(newBlocks []*core.Block) error` locates the fork point, validates the entire candidate chain, compares total work, and only if the candidate is strictly heavier, rolls back to the fork point and replays the new blocks — a **reorganization**, or reorg.
- Rolled-back transactions return to the mempool as pending, unless they're already confirmed in the newly adopted chain or conflict with something that is; orphaned coinbase rewards simply disappear along with their block.
- Deep reorgs are far more dangerous than shallow ones — they can undo transactions everyone believed were final — which is exactly why real systems wait for multiple confirmations before treating high-value transactions as settled, and exactly what a 51% attack (Chapter 76) tries to exploit.

---

## Exercises

### Easy

1. **Trace the Section 7 simulation by hand** but change Step 3 so Y mines only one additional block (Block 12-Y at Difficulty 19, slightly *lower* than usual) instead of Difficulty 21. Recompute both sides' total work above the fork point and determine whether X should still reorganize onto Y's chain.

2. **Write a unit test for `ChainWork`** on a small, manually constructed chain of blocks with known `Difficulty` values (for example, 3 blocks at difficulties 10, 12, and 11), and assert it returns exactly `2^10 + 2^12 + 2^11`.

3. **Explain in your own words** why `ReplaceChain` still calls `ValidateBlock` on every block in the candidate chain, even though the candidate chain is, by definition, "heavier" than our current one. What could go wrong if it skipped validation and only compared `ChainWork`?

### Medium

4. **Implement `heightOfHash`, `blocksFrom`, `rollbackTo`, and `applyBlock`** on `core.Blockchain` if you haven't already — the four internal helpers `ReplaceChain` relies on — and write a test that reorgs a 5-block chain onto a competing 6-block chain that forks off after block 2, asserting the final chain matches the competitor exactly and that `ChainWork` afterward equals the competitor's own precomputed total.

5. **Write a test for `returnOrphanedTransactions`** that constructs an orphaned block containing two transactions — one that also appears in the newly adopted chain, and one that doesn't — and asserts only the second one ends up back in the mempool after a reorg.

6. **Simulate a 3-block-deep fork** (extend Section 7's simulation so Node Y is 3 blocks ahead by the time the partition heals, not just 1) and confirm, either by hand or with a real test, that `ReplaceChain` correctly rolls back and replays all 3 blocks, and that a coinbase transaction inside one of the orphaned blocks is correctly dropped rather than resurrected.

### Hard

7. **Implement confirmation-count tracking.** Add a method `Confirmations(blockHash []byte) (int, error)` to `core.Blockchain` that returns how many blocks currently sit on top of the given block (0 if it's not on the current chain at all, or if it's the current tip). Wire this into the CLI wallet from Chapter 36 so `gochain wallet balance` can optionally report "confirmed" vs. "pending" balance based on a configurable confirmation threshold.

8. **Simulate a shallow 51%-style attack** in a controlled test: two nodes, X (honest) and M (secretly holding back mined blocks instead of broadcasting them immediately). Have M privately mine 2 blocks while X mines and broadcasts 1 publicly (with a transaction X believes is now "confirmed"). Then have M reveal its 2-block chain. Confirm X reorganizes onto M's chain via `ReplaceChain`, and that the transaction X thought was confirmed is either gone or conflicts with something in M's chain. Write up, in a short comment block, what specifically about GoChain's current design allowed this to work, and what confirmation count would have prevented X from being fooled.

9. **Investigate and implement a "checkpoint" defense.** Real blockchains sometimes hard-code a handful of known-good block hashes at specific heights into the client software itself, so a node will refuse to reorg past a checkpoint no matter how much work a competing chain claims, protecting against implausibly deep reorgs from a chain of forged history. Add an optional `[]Checkpoint{Height uint64; Hash []byte}` list to `core.Blockchain`, and modify `ReplaceChain` to reject any candidate chain that would roll back past a configured checkpoint, even if its `ChainWork` is technically higher. Discuss, in a short comment, the trade-off this introduces between security against deep reorgs and reintroducing a small amount of centralized trust (whoever chooses the checkpoints).
