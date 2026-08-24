# Chapter 19: Block Validation and Immutability

Chapter 18 left a deliberate gap: `AddBlock` would happily accept a block with a stale hash, a broken link, or transactions altered after the fact. This chapter closes that gap for real. `ValidateBlock` gives GoChain a precise, callable answer to "is this block actually telling the truth about itself and about what came before it?" — and then this chapter does something most tutorials only describe in prose: it actually tampers with a real block in a running chain and watches, hash by hash, exactly how far the damage spreads.

## Table of Contents

1. [Two Questions Every Valid Block Must Answer](#1-two-questions-every-valid-block-must-answer)
2. [ValidateBlock — Checking One Block's Own Integrity](#2-validateblock--checking-one-blocks-own-integrity)
3. [Checking the Link — PrevBlockHash Against Reality](#3-checking-the-link--prevblockhash-against-reality)
4. [Rewriting AddBlock to Validate Before Accepting](#4-rewriting-addblock-to-validate-before-accepting)
5. [Validating the Whole Chain From Genesis](#5-validating-the-whole-chain-from-genesis)
6. [Worked Example — Tampering With an Old Block](#6-worked-example--tampering-with-an-old-block)
7. [Tamper-Evident vs. Tamper-Proof](#7-tamper-evident-vs-tamper-proof)
8. [Testing Validation](#8-testing-validation)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Two Questions Every Valid Block Must Answer

Every block that claims a place in GoChain's chain must answer two independent questions truthfully, and "valid" means nothing more or less than both answers checking out. The first question actually needs two separate checks to answer honestly, for a reason worth being precise about before any code.

- **Is this block honest about itself?** This splits into two checks, because `Block.Hash` (Chapter 17) is computed over the block's *header* fields — including the stored `MerkleRoot` — never over the raw transaction list directly (Chapter 17, Section 4 explained why: `MerkleRoot` already stands in for the transactions). That means recomputing `ComputeHash()` alone only re-checks the header against itself; it says nothing about whether the header's `MerkleRoot` still honestly summarizes the transactions actually sitting in `b.Transactions` right now. So this question needs both:
  - **Header self-consistency:** does a freshly computed `ComputeHash()` match the stored `Hash`?
  - **Body-matches-header:** does a freshly built Merkle root over `b.Transactions` (Chapter 10's tooling, reused directly) match the stored `MerkleRoot`? If someone altered a transaction after the block was built, without recomputing `MerkleRoot` to match, this is the check that catches it — the header-hash check alone would not.
- **Is this block honest about its place in the chain?** Does its `PrevBlockHash` actually equal the *real*, current `Hash` of the block that is supposed to come before it? If it points to a hash that does not match reality — whether because of an honest mistake or a deliberate forgery — this check catches that too.

Think of a bank teller checking a paper check. First, the teller checks the check itself is not altered — no smudged, changed amount, no forged signature (self-consistency: is this document internally consistent?). Then, separately, the teller checks that the amount written in numerals actually matches the amount spelled out in words on the same check — a document can be internally tidy and still lie about what it claims to represent (body-matches-header: does the summary match the details?). Finally, the teller checks the check against the account it claims to draw from (the link: does this document's claim about the world outside itself hold up?). A block needs all three kinds of scrutiny, and `ValidateBlock` performs all three, in that order.

```
  Self-consistency        Body matches header       Link to reality

  freshly computed hash    freshly built Merkle       b.PrevBlockHash
        ==?                  root over b.Transactions        ==?
  b.Hash (stored)                  ==?                previous block's
                             b.MerkleRoot (stored)      real Hash

  All three must hold for a block to be considered valid.
```

## 2. ValidateBlock — Checking One Block's Own Integrity

The first question from Section 1 needs nothing more than the tools Chapters 10 and 17 already built: call `ComputeHash()` fresh and compare it against the stored `Hash`, then call `MerkleRootOf()` fresh over the actual transactions and compare it against the stored `MerkleRoot`.

```go
// core/validate.go
package core

import (
	"bytes"
	"fmt"
)

// ValidateBlock checks that b is internally consistent, that its
// header honestly summarizes its own body, and that it is correctly
// linked to the block that should precede it in bc. It never trusts
// any of b's own stored fields at face value -- it always recomputes
// and compares.
func (bc *Blockchain) ValidateBlock(b *Block) error {
	if b == nil {
		return fmt.Errorf("cannot validate a nil block")
	}

	// Check 1 (Section 1): is the header itself self-consistent?
	if !bytes.Equal(b.Hash, b.ComputeHash()) {
		return fmt.Errorf(
			"block %d: stored hash %x does not match recomputed hash %x",
			b.Height, b.Hash, b.ComputeHash(),
		)
	}

	// Check 2 (Section 1): does the header's MerkleRoot actually
	// summarize the transactions the block body claims to contain?
	// This is the check that catches a transaction altered after the
	// block was built, without its MerkleRoot being recomputed to match.
	freshRoot := MerkleRootOf(b.Transactions)
	if !bytes.Equal(b.MerkleRoot, freshRoot) {
		return fmt.Errorf(
			"block %d: stored MerkleRoot %x does not match a freshly computed Merkle root %x over its transactions",
			b.Height, b.MerkleRoot, freshRoot,
		)
	}

	// Check 3 (Section 1): is this block honest about its place in
	// the chain? Section 3 implements this check.
	if err := bc.validateLink(b); err != nil {
		return err
	}

	return nil
}
```

The header-hash check is the exact comparison Chapter 17, Section 8 asked you to try by hand — `block.Hash` versus `block.ComputeHash()` — now wrapped in a real method with a real, specific error message instead of a silently-ignored mismatch. The Merkle-root check reuses `MerkleRootOf`, the exact function Chapter 17's `NewBlock` calls to build a block's `MerkleRoot` in the first place (Chapter 17 exports it for exactly this reason), so "does the body match the header" is never a separate, parallel implementation to keep in sync — it is the same Merkle-tree logic, called again, on demand. Every error message here includes both the stored and the freshly recomputed value: when a check fails in a real debugging session (Chapter 21's inspector CLI will surface exactly this), seeing both values side by side is far more useful than a bare "invalid block" message.

## 3. Checking the Link — PrevBlockHash Against Reality

The second question needs one more piece: a way to find "the block that should precede `b`" inside the chain. Since `bc.blocks` is ordered by height starting at 0 (Chapter 18, Section 4), the block at height `b.Height - 1` is exactly that predecessor — as long as it actually exists in the chain yet.

```go
// core/validate.go

// blockAtHeight returns the block stored at height, or nil if no
// block at that height exists yet in bc.
func (bc *Blockchain) blockAtHeight(height int64) *Block {
	if height < 0 || height >= int64(len(bc.blocks)) {
		return nil
	}
	return bc.blocks[height]
}

// validateLink checks that b's PrevBlockHash correctly points back
// to the real Hash of the block before it -- or, for a genesis
// block, that PrevBlockHash is the all-zero placeholder from Chapter
// 18, Section 2.
func (bc *Blockchain) validateLink(b *Block) error {
	if b.Height == 0 {
		zeroHash := make([]byte, 32)
		if !bytes.Equal(b.PrevBlockHash, zeroHash) {
			return fmt.Errorf("genesis block must have an all-zero PrevBlockHash, got %x", b.PrevBlockHash)
		}
		return nil
	}

	prev := bc.blockAtHeight(b.Height - 1)
	if prev == nil {
		return fmt.Errorf("block %d: no block found at height %d to link against", b.Height, b.Height-1)
	}

	if !bytes.Equal(b.PrevBlockHash, prev.Hash) {
		return fmt.Errorf(
			"block %d: PrevBlockHash %x does not match block %d's real hash %x",
			b.Height, b.PrevBlockHash, prev.Height, prev.Hash,
		)
	}

	return nil
}
```

Notice `validateLink` reads `prev.Hash` directly — the *stored* hash of the previous block — rather than recomputing it again from scratch here. That is not a shortcut taken for convenience; it is correct by construction, because `ValidateBlock` (Section 2) is exactly what guarantees, for any block already accepted into the chain, that its stored `Hash` genuinely matches its own contents. Validation composes: as long as every block already in `bc.blocks` passed `ValidateBlock` on the way in (Section 4 makes sure of this), trusting `prev.Hash` here is not a leap of faith, it is relying on a check that already happened.

## 4. Rewriting AddBlock to Validate Before Accepting

With `ValidateBlock` in place, Chapter 18's `AddBlock` gets the one-line change that turns it from "accepts anything" into an actual gatekeeper:

```go
// core/blockchain.go

// AddBlock validates b against the chain's current state and, only
// if it passes, appends it and advances the tip. This replaces
// Chapter 18's version, which performed no real validation at all.
func (bc *Blockchain) AddBlock(b *Block) error {
	if b == nil {
		return fmt.Errorf("cannot add a nil block")
	}

	if err := bc.ValidateBlock(b); err != nil {
		return fmt.Errorf("rejecting block %d: %w", b.Height, err)
	}

	bc.blocks = append(bc.blocks, b)
	bc.tip = b.Hash

	return nil
}
```

Try the exact experiment Chapter 18, Section 8 suggested and left unresolved — construct a block whose `PrevBlockHash` is made-up bytes that do not match `chain.Tip()`, and call `chain.AddBlock(fakeBlock)`:

```go
chain := core.NewBlockchain()

fake := core.NewBlock([]*core.Transaction{testTx("tx1")}, []byte("not the real previous hash!!"), 1)
err := chain.AddBlock(fake)
fmt.Println(err)
// rejecting block 1: block 1: PrevBlockHash 6e6f742074686520... does not match block 0's real hash 7a3fc9e1...
```

This is the gap Chapter 18 named and left open, now closed. `%w` in the `fmt.Errorf` call wraps the underlying error rather than just formatting it into a string, which is idiomatic Go error handling from Chapter 04 — a caller further up the stack can use `errors.Is` or `errors.Unwrap` to inspect the original cause if it needs to, rather than only ever seeing a flattened message.

## 5. Validating the Whole Chain From Genesis

`ValidateBlock` checks one block against the chain it is joining. A separate, equally useful question is: given a chain that already exists, is *every single block in it*, from genesis all the way to the tip, still valid right now? This is exactly what Chapter 21's `--verify` flag needs, so it earns its own method:

```go
// core/validate.go

// ValidateChain walks every block in bc, from genesis to tip, and
// returns an error describing the first invalid block it finds, or
// nil if the entire chain checks out.
func (bc *Blockchain) ValidateChain() error {
	for _, b := range bc.blocks {
		if err := bc.ValidateBlock(b); err != nil {
			return err
		}
	}
	return nil
}
```

`ValidateChain` deliberately reuses `ValidateBlock` for every single block, rather than writing a separate, parallel "check the whole chain" implementation. This is worth pausing on: `ValidateBlock` already checks both self-consistency and the link backward, so calling it once per block, walking forward from genesis, automatically checks every link in the entire chain — there is no additional logic `ValidateChain` needs of its own. Reusing a single, well-tested building block like this instead of re-deriving the same logic twice is exactly the kind of small design discipline that keeps a growing codebase trustworthy.

## 6. Worked Example — Tampering With an Old Block

Now for the hands-on part this chapter exists to deliver: build a real, four-block chain, reach in and alter a transaction inside an old block directly, and watch precisely what breaks — and what it takes to hide the tampering, block by block.

```go
chain := core.NewBlockchain()

for i := 1; i <= 3; i++ {
	tx := &core.Transaction{ID: []byte(fmt.Sprintf("tx-%d", i)), Timestamp: time.Now().Unix()}
	next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)
	if err := chain.AddBlock(next); err != nil {
		log.Fatal(err)
	}
}

fmt.Println("Before tampering:", chain.ValidateChain()) // <nil>
```

**Step 1 — tamper with Block 1's transaction, change nothing else.** Reach directly into `chain.Blocks()[1]` and alter its transaction's `ID`, simulating an attacker editing stored data directly on disk without going through any of GoChain's normal code paths:

```go
tampered := chain.Blocks()[1]
tampered.Transactions[0].ID = []byte("tx-1-FORGED")

fmt.Println("After tampering Block 1, before any fix-up at all:")
fmt.Println(chain.ValidateBlock(chain.Blocks()[1]))
// block 1: stored MerkleRoot 9c2e0a71... does not match a freshly computed Merkle root 3d8f1a29... over its transactions
```

Block 1 fails immediately — but notice exactly *which* check catches it. Its stored `Hash` is untouched, and `ComputeHash()` only re-hashes the header fields (Height, Timestamp, PrevBlockHash, the stored `MerkleRoot`, Nonce, transaction count), none of which changed, so the header-hash check alone would have missed this entirely. What actually catches the tampering is Section 2's second check: a freshly built Merkle root over Block 1's (now-forged) transaction no longer matches the `MerkleRoot` still sitting in its header, unchanged since the block was honestly built. This is Chapter 10's Merkle tree machinery, and the avalanche effect (Chapter 08, Section 3) underneath it, doing exactly the job they were built for — one changed byte inside one transaction, and the freshly recomputed root looks nothing like the one still stored in `MerkleRoot`.

**Step 2 — the attacker tries to cover their tracks by recomputing Block 1's MerkleRoot and hash.** Suppose the attacker is not careless — they know both checks from Section 2 must pass, so they fix both, in the correct order (`MerkleRoot` first, since `Hash` is computed over it):

```go
tampered.MerkleRoot = core.MerkleRootOf(tampered.Transactions) // recomputed over the forged transaction
tampered.Hash = tampered.ComputeHash()                          // recomputed over the new MerkleRoot

fmt.Println("Block 1 alone, after the attacker 'fixes' both fields:")
fmt.Println(chain.ValidateBlock(chain.Blocks()[1]))
// <nil> -- Block 1 now passes its OWN self-checks!
```

Block 1, checked entirely on its own, now looks perfectly self-consistent — its `MerkleRoot` really does summarize its (forged) transaction, and its `Hash` really does match its (also forged) header. This is an important, slightly unsettling thing to see with your own eyes: self-consistency alone is not enough to catch tampering, which is exactly why Section 1 insisted on checking the *link* too, not just the block in isolation.

**Step 3 — but Block 2 still points at the old hash.** Block 2's `PrevBlockHash` was set, back when it was built, to Block 1's *original*, honest hash — and that field has not changed, because nobody touched Block 2 at all:

```go
fmt.Println("Block 2, unchanged, checked against the new Block 1:")
fmt.Println(chain.ValidateBlock(chain.Blocks()[2]))
// block 2: PrevBlockHash c81e44b0... does not match block 1's real hash 5f3ab291...

fmt.Println("Whole-chain check:")
fmt.Println(chain.ValidateChain())
// block 2: PrevBlockHash c81e44b0... does not match block 1's real hash 5f3ab291...
```

The forgery is caught anyway — one hop later, at Block 2's link check instead of Block 1's self-check, but caught nonetheless. To make Block 2 pass, the attacker would have to update Block 2's `PrevBlockHash` to the new value and recompute *Block 2's* hash too — which then breaks Block 3's link, for the exact same reason. Continue this "fix one, break the next" pattern all the way to the current tip, and you have reproduced, by hand, precisely the cascading rewrite Chapter 16, Section 6 described on paper: tampering with block `N` forces recomputing every single block from `N` to the tip, in order, with no way to stop partway through and still have a chain that validates.

```
   Block 0        Block 1 (forged)     Block 2 (unchanged)    Block 3 (unchanged)
  +--------+      +--------+           +--------+             +--------+
  | h: 7a3f|----->| h: 5f3a|    X      | prev: c81e (STALE) |  | prev: 44b0 |
  +--------+      +--------+   /|\     +--------+             +--------+
                                 |
                       link check fails here --
                       Block 2 still expects the ORIGINAL Block 1 hash
```

## 7. Tamper-Evident vs. Tamper-Proof

Section 6 demonstrated something worth naming precisely, because the vocabulary matters for the rest of this course. Nothing in `ValidateBlock` or `ValidateChain` made it *impossible* to tamper with Block 1, or to fix up Block 2 and Block 3 afterward to hide it. An attacker with direct access to the data — exactly the scenario Section 6 simulated — genuinely can recompute every hash from the tampered block forward, and the resulting chain will pass `ValidateChain()` with no complaint at all, because every link really would hold, honestly, in the forged version.

What GoChain has built through Chapter 19 is **tamper-evident**, not **tamper-proof**:

- **Tamper-evident** means any tampering *leaves a detectable trace* — Section 6 showed that trace is unavoidable the moment you check any single block in isolation right after tampering, and it forces a cascading, mechanical rewrite of every later block to erase.
- **Tamper-proof** would mean tampering is not just detectable but *infeasible to pull off at all* — even the cascading rewrite from Section 6 would be too expensive or too slow to actually carry out.

GoChain, as built through this chapter, gives you the first property but not yet the second. Rewriting three blocks by hand, as Section 6 did, took a few lines of Go and ran in microseconds. Nothing stopped it. What *would* stop it — making each block's hash expensive enough to compute that redoing that cascade for a long chain becomes genuinely, practically infeasible, not merely inconvenient — is proof of work, and it is the entire subject of Volume 4. This chapter's honest, unglamorous conclusion is exactly the one Chapter 16, Section 6 previewed: tamper-evidence is real, valuable, and now working code — but it is only half the story.

```
                    TAMPER-EVIDENT              TAMPER-PROOF
                  (this volume builds)         (Volume 4 builds)

  Tampering        Immediately detectable       Immediately detectable
  detection?        by ValidateBlock              by ValidateBlock

  Tampering         Cheap -- rewrite every       Expensive -- redoing
  cost?             later block's hash by         proof-of-work for every
                    hand, in microseconds         later block, competing
                                                    against the entire
                                                    honest network
```

## 8. Testing Validation

The two failure modes Section 6 walked through by hand belong in the permanent test suite, alongside the "everything is fine" happy path.

```go
// core/validate_test.go
package core

import (
	"strings"
	"testing"
)

func buildTestChain(t *testing.T, n int) *Blockchain {
	t.Helper()
	chain := NewBlockchain()
	for i := 1; i <= n; i++ {
		next := NewBlock([]*Transaction{testTx("tx")}, chain.Tip(), chain.Height()+1)
		if err := chain.AddBlock(next); err != nil {
			t.Fatalf("unexpected error building test chain: %v", err)
		}
	}
	return chain
}

func TestValidateChain_HealthyChainPasses(t *testing.T) {
	chain := buildTestChain(t, 4)
	if err := chain.ValidateChain(); err != nil {
		t.Fatalf("expected a healthy chain to validate, got: %v", err)
	}
}

func TestValidateBlock_CatchesHashTamper(t *testing.T) {
	chain := buildTestChain(t, 2)

	victim := chain.Blocks()[1]
	victim.Hash = []byte("not the real hash at all")

	err := chain.ValidateBlock(victim)
	if err == nil {
		t.Fatal("expected an error for a block with a corrupted Hash field, got nil")
	}
	if !strings.Contains(err.Error(), "does not match recomputed hash") {
		t.Fatalf("expected a hash-mismatch error, got: %v", err)
	}
}

func TestValidateBlock_CatchesMerkleRootTamper(t *testing.T) {
	chain := buildTestChain(t, 3)

	tampered := chain.Blocks()[1]
	tampered.Transactions[0].ID = []byte("forged")
	// MerkleRoot and Hash are deliberately NOT recomputed here.

	err := chain.ValidateBlock(tampered)
	if err == nil {
		t.Fatal("expected an error for a block whose MerkleRoot no longer matches its transactions")
	}
	if !strings.Contains(err.Error(), "does not match a freshly computed Merkle root") {
		t.Fatalf("expected a Merkle-root-mismatch error, got: %v", err)
	}
}

func TestValidateChain_CatchesBrokenLinkAfterFullFixUp(t *testing.T) {
	chain := buildTestChain(t, 3)

	tampered := chain.Blocks()[1]
	tampered.Transactions[0].ID = []byte("forged")
	tampered.MerkleRoot = MerkleRootOf(tampered.Transactions) // covers the Merkle-root check
	tampered.Hash = tampered.ComputeHash()                    // covers the header-hash check

	err := chain.ValidateChain()
	if err == nil {
		t.Fatal("expected ValidateChain to catch the broken link at block 2, got nil")
	}
	if !strings.Contains(err.Error(), "does not match block") {
		t.Fatalf("expected a broken-link error, got: %v", err)
	}
}

func TestAddBlock_RejectsBadPrevBlockHash(t *testing.T) {
	chain := NewBlockchain()

	fake := NewBlock([]*Transaction{testTx("tx1")}, []byte("not the real hash"), 1)
	if err := chain.AddBlock(fake); err == nil {
		t.Fatal("expected AddBlock to reject a block with a bad PrevBlockHash")
	}
}
```

`TestValidateChain_CatchesBrokenLinkAfterFullFixUp` is the automated version of Section 6's entire worked example: tamper, fully "fix" the tampered block's own `MerkleRoot` and `Hash`, and confirm the *chain-wide* check still catches it — because the block one hop later still expects the original hash. Run `go test ./core/...`; all five should pass, and they now encode, permanently, the exact tamper-evidence guarantee this chapter spent its worked example demonstrating by hand.

---

## Summary

- A valid block must pass three checks: is its header self-consistent (`Hash == ComputeHash()`), does its header's `MerkleRoot` actually summarize its current transactions (`MerkleRoot == MerkleRootOf(Transactions)`), and is it honest about its place in the chain (`PrevBlockHash` matches the real previous block's `Hash`)?
- The header-hash check alone cannot catch a tampered transaction, because `ComputeHash()` only hashes header fields (including the *stored* `MerkleRoot`), never the raw transaction list directly — the Merkle-root check is what actually re-examines the transactions themselves.
- `ValidateBlock` checks all three, in order, returning a specific error describing exactly which check failed and with what values.
- `validateLink` treats a genesis block (`Height == 0`) as a special case requiring an all-zero `PrevBlockHash`, and every other block as requiring a real match against the block at `Height - 1`.
- `AddBlock`, rewritten this chapter, now calls `ValidateBlock` before ever appending a block, closing the gap Chapter 18 deliberately left open.
- `ValidateChain` walks the entire chain from genesis, reusing `ValidateBlock` on every block, to confirm every single link still holds right now.
- The worked example proved, hands-on, that tampering with one block's transaction is detected instantly by that block's Merkle-root check, and that "fixing" the tampered block's own `MerkleRoot` and `Hash` only pushes the detectable failure one hop forward, to the next block's link check — the tampering can be hidden only by cascading the fix-up forward, block by block, all the way to the tip.
- GoChain, as built through this chapter, is **tamper-evident** — tampering always leaves a detectable trace — but not yet **tamper-proof** — nothing yet makes recomputing that cascade of hashes expensive. That is Volume 4's proof-of-work chapters.

---

## Exercises

### Easy

1. Build a chain of 5 blocks, then deliberately corrupt the *stored* `Hash` field of block 3 (set it to some arbitrary bytes) without touching its actual contents. Call `chain.ValidateBlock` on block 3 and report the exact error message. Which of Section 1's three checks does this failure correspond to, and why does corrupting `Hash` alone (leaving `MerkleRoot` untouched) not also trigger the Merkle-root check?
2. Explain, in your own words, why `validateLink` treats `Height == 0` as a special case instead of always looking up `bc.blockAtHeight(b.Height - 1)` unconditionally.
3. `ValidateChain` returns only the *first* invalid block it finds, then stops. Propose a small change that would instead collect and return *every* invalid block found, and explain one situation where that would be more useful to a caller than stopping at the first failure.

### Medium

4. Reproduce Section 6's full worked example yourself in Go: build a 4-block chain, tamper with block 2's transaction, observe the Merkle-root check failure, "fix" block 2's `MerkleRoot` and `Hash`, observe the link-check failure move to block 3, and then fully cascade the fix-up all the way to the tip so `ValidateChain()` finally passes again. Report how many fields across how many blocks you had to change in total.
5. Add a new validation rule to `ValidateBlock`: reject any block whose `Timestamp` is earlier than the previous block's `Timestamp` (a block cannot claim to have been created before its own predecessor). Write a test that constructs such a block and confirms it is now rejected.
6. `blockAtHeight` does a direct slice index (`bc.blocks[height]`), which relies on `Height` values in `bc.blocks` never having gaps or duplicates. Write a test that manually constructs a `Blockchain` (bypassing `AddBlock`) with a gap in heights (blocks at height 0, 1, and 3, skipping 2) and describe what `blockAtHeight(2)` and `blockAtHeight(3)` actually return in that broken state. Propose one guard `AddBlock` could add to make this state impossible to reach through normal use.

### Hard

7. Section 7 draws a sharp line between tamper-evident and tamper-proof. Write a design note (300-450 words) estimating, in concrete terms, how long it would take a single attacker with access to GoChain's current (non-proof-of-work) code to tamper with a block 10,000 blocks behind the current tip and fully cascade the fix-up to the tip, assuming each `NewBlock` call (and its Merkle root computation) takes on the order of microseconds. Contrast this with your intuition for how expensive the same cascade would need to become for GoChain to be considered tamper-proof, previewing Volume 4's proof-of-work chapters.
8. Implement a `RepairChain(bc *core.Blockchain) int` function that walks a chain from genesis, and for any block whose `Hash` no longer matches `ComputeHash()`, recomputes and overwrites `Hash`, then continues forward fixing every subsequent broken link the same way, returning the number of blocks it had to fix. Discuss, in a short paragraph, why building this function is educational but why no real blockchain node should ever call anything like it in production.
9. Design (in prose, 300-500 words, with at least one diagram) an alternative to height-indexed lookup in `blockAtHeight` that instead looks up a block by its own `Hash` value using a Go map (`map[string]*Block`, keyed by the hex-encoded hash). Explain what new capability this would add (hint: consider Chapter 50's later discussion of forks, where two blocks might briefly claim the same height), and what `Blockchain` would need to track differently to support it correctly.
