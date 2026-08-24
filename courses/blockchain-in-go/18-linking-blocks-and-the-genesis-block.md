# Chapter 18: Linking Blocks and the Genesis Block

Chapter 17 gave GoChain a `Block` that can hash itself. A single hashed block, sitting alone, is not yet a blockchain — it is just one well-formed struct. A chain forms the moment a *second* block stores the *first* block's hash inside its own `PrevBlockHash` field, and a third stores the second's, and so on. This chapter builds that link, handles the one block that has no predecessor at all — the genesis block — and wraps the whole thing in a `core.Blockchain` type you can grow one block at a time.

## Table of Contents

1. [What Makes a Chain a Chain](#1-what-makes-a-chain-a-chain)
2. [The Genesis Block — A Block With No Predecessor](#2-the-genesis-block--a-block-with-no-predecessor)
3. [NewGenesisBlock in Go](#3-newgenesisblock-in-go)
4. [The Blockchain Type](#4-the-blockchain-type)
5. [NewBlockchain — Starting a Chain From Genesis](#5-newblockchain--starting-a-chain-from-genesis)
6. [AddBlock — Extending the Chain](#6-addblock--extending-the-chain)
7. [Drawing the Full Chain — Four Blocks, Linked by Hash Arrows](#7-drawing-the-full-chain--four-blocks-linked-by-hash-arrows)
8. [What AddBlock Does Not Yet Check](#8-what-addblock-does-not-yet-check)
9. [Testing core.Blockchain](#9-testing-coreblockchain)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Makes a Chain a Chain

Picture a paper trail used to prove a document's history: each new page is stapled directly on top of the one before it, and at the top of every new page someone writes down a short description of the page underneath — "this continues the page ending in the sentence about the harbor." Anyone holding just the newest page can flip backward, checking each description against the page it actually sits on top of, all the way down to the very first page. If any page in the middle were secretly swapped out, the description written on the page above it would stop matching — a mismatch anyone could notice just by comparing two adjacent pages.

`core.Block.PrevBlockHash` is that written description, except instead of a human-written sentence, it is an exact SHA-256 fingerprint (Chapter 08) of the entire previous block. A **chain** exists the moment this one condition holds true for every block in a sequence, without exception:

```
  block[i].PrevBlockHash  ==  block[i-1].Hash        (for every i > 0)
```

That single equation, checked repeatedly, is the entire linking mechanism. There is no separate index of "what comes after what," no external table of contents — each block simply carries a copy of its predecessor's fingerprint, and the chain's integrity is nothing more than every one of those copies matching reality.

```
   Block 0            Block 1            Block 2
  +--------+         +--------+         +--------+
  | Hash A |  <-----  |PrevHash: A|      |          |
  +--------+         +--------+         +--------+
                      | Hash B |  <-----  |PrevHash: B|
                      +--------+         +--------+
                                          | Hash C |
                                          +--------+

  Read the arrows as "points back to." No block needs to know
  anything about blocks before the one it points to directly.
```

Notice something worth sitting with: Block 2 does not need to know anything about Block 0 at all — only about Block 1. Trust does not have to travel all the way back to the beginning in one giant leap; it travels one hop at a time, and if every single hop holds, the whole chain holds, all the way back. Chapter 19 turns "checking that every hop holds" into a real, callable function; this chapter's job is making sure that link exists to check in the first place.

## 2. The Genesis Block — A Block With No Predecessor

Every chain has to start somewhere, and the very first block has a problem none of the others do: there is no previous block to point back to. Bitcoin's solution, which GoChain adopts unchanged, is a simple convention rather than a special case in the hashing logic: the **genesis block** — height 0, the first block in any GoChain blockchain — sets `PrevBlockHash` to 32 zero bytes, a value no real block will ever actually produce as a hash (SHA-256's output space is so vast, per Chapter 08, Section 5, that a genuine all-zero hash is not a realistic occurrence).

```
   Genesis Block (Height 0)
  +---------------------------+
  | Height:        0           |
  | PrevBlockHash: 00000...000 |  <-- 32 zero bytes; "there is
  |                             |      nothing before me"
  | Hash:          7a3fc9e1... |
  +---------------------------+
```

**New term — genesis block:** the first block in a blockchain, at height 0, whose `PrevBlockHash` is conventionally all zeroes because no real predecessor exists. Everything else about the genesis block — it still gets hashed, still gets a `MerkleRoot`, still holds real transactions — works exactly the same as any other block. The *only* thing special about it is what it points backward to: nothing, represented honestly as an all-zero placeholder rather than left as `nil` or some other ambiguous value.

Real Bitcoin's genesis block, mined by Satoshi Nakamoto on January 3rd, 2009, follows this exact convention — Chapter 16, Section 4 showed its real `PrevBlockHash`: 32 zero bytes, exactly what GoChain uses. This is one of the small, concrete places where GoChain's design is not a simplified approximation of a real blockchain but the identical convention, just implemented by you instead of Satoshi.

## 3. NewGenesisBlock in Go

Building a genesis block turns out to need no new machinery at all — `NewBlock` from Chapter 17 already does everything required, given the right arguments:

```go
// core/blockchain.go
package core

// NewGenesisBlock builds the first block in a chain: height 0, with
// PrevBlockHash set to 32 zero bytes (Section 2) because no real
// predecessor exists. coinbase is typically a reward transaction --
// Chapter 37 explains coinbase transactions properly, once real
// mining rewards exist. For this volume, any *Transaction is enough.
func NewGenesisBlock(coinbase *Transaction) *Block {
	prevBlockHash := make([]byte, 32) // 32 zero bytes, per Section 2
	return NewBlock([]*Transaction{coinbase}, prevBlockHash, 0)
}
```

`make([]byte, 32)` allocates a 32-byte slice, and Go zero-initializes every element of a newly allocated slice by default — so this line produces exactly the all-zero value Section 2 described, with no explicit loop needed. Everything downstream of this call — `NewBlock` computing a `MerkleRoot` over `[coinbase]`, then computing `Hash` over the resulting header fields — runs completely unchanged from Chapter 17. The genesis block is not a different code path; it is `NewBlock`, called once, with a deliberately chosen `prevBlockHash` and `height`.

## 4. The Blockchain Type

A single `*Block` is not yet a usable chain — nothing tracks the sequence, and nothing remembers which block is currently the newest. `core.Blockchain` fixes that with the simplest structure that works for this volume: an in-memory slice of every block, plus a cached pointer to the newest one.

```go
// core/blockchain.go

// Blockchain is GoChain's in-memory (and, starting in Chapter 20,
// flat-file-backed) ledger: every block ever added, in order, plus a
// cached "tip" -- the current newest block's hash. A real database
// replaces this simple slice entirely in Volume 8; until then, this
// is intentionally the simplest structure that can hold a real,
// linked, tamper-evident chain.
type Blockchain struct {
	tip    []byte
	blocks []*Block
}
```

Both fields are unexported (lowercase), which is a deliberate Go design choice this course has followed since Chapter 04: nothing outside `core` should be able to reach in and directly append to `blocks` or overwrite `tip`, bypassing the linking and (starting next chapter) validation rules `AddBlock` enforces. Every legitimate way to change a `Blockchain`'s contents goes through a method. A handful of small, read-only accessor methods round this out:

```go
// core/blockchain.go

// Tip returns the hash of the current newest block -- exactly the
// PrevBlockHash the next block built with NewBlock should use.
func (bc *Blockchain) Tip() []byte {
	return bc.tip
}

// Height returns the height of the current newest block. The
// genesis block is height 0, so a chain holding only genesis
// reports height 0, not 1 -- consistent with Chapter 16's "Height
// counts from zero" rule.
func (bc *Blockchain) Height() int64 {
	return bc.blocks[len(bc.blocks)-1].Height
}

// Blocks returns every block in the chain, from genesis to tip, in
// order. Chapter 21's inspector CLI walks this slice directly.
func (bc *Blockchain) Blocks() []*Block {
	return bc.blocks
}

// LastBlock returns the current newest block.
func (bc *Blockchain) LastBlock() *Block {
	return bc.blocks[len(bc.blocks)-1]
}
```

## 5. NewBlockchain — Starting a Chain From Genesis

A `Blockchain` should never exist without a genesis block already inside it — an empty chain with no block zero at all is not a smaller version of a real chain, it is simply not a chain yet. `NewBlockchain` enforces this by building the genesis block itself, every time:

```go
// core/blockchain.go

import (
	"time"

	"github.com/you/gochain/crypto"
)

// NewBlockchain creates a brand new chain containing only its
// genesis block. The genesis block's single transaction is a
// placeholder coinbase-style transaction -- Volume 5 gives real
// coinbase transactions actual mining-reward meaning; here, its only
// job is to give the genesis block something to compute a
// MerkleRoot over, matching every other block's shape.
func NewBlockchain() *Blockchain {
	coinbase := &Transaction{
		ID:        crypto.Hash([]byte("gochain-genesis-coinbase")),
		Timestamp: time.Now().Unix(),
	}

	genesis := NewGenesisBlock(coinbase)

	return &Blockchain{
		tip:    genesis.Hash,
		blocks: []*Block{genesis},
	}
}
```

Every `Blockchain` your code ever creates by calling `NewBlockchain()` starts life with exactly one block, at height 0, with an all-zero `PrevBlockHash` and a real, computed `Hash` — never an empty slice waiting to be filled in later. This mirrors a design decision Chapter 16, Section 2 already flagged: `Height`'s zero value is not a "nothing here yet" placeholder, it is the genesis block's true, correct height, and `NewBlockchain` makes sure that block actually exists from the very first call.

## 6. AddBlock — Extending the Chain

Growing the chain by one block is, at this stage of the course, a small and honest method — it appends a block and moves the tip forward. It intentionally does *not* yet check whether the block is actually valid; that is Chapter 19's entire subject, and giving it its own chapter is not padding, it is treating "does this block link in correctly" as a serious question deserving a serious answer, rather than a one-line afterthought bolted onto `AddBlock` today.

```go
// core/blockchain.go

import "errors"

// AddBlock appends b to the chain and advances the tip to b's hash.
// This version performs only the most basic sanity check -- Chapter
// 19 rewrites this method's body to call ValidateBlock first,
// rejecting any block whose hash or link to the previous block does
// not check out.
func (bc *Blockchain) AddBlock(b *Block) error {
	if b == nil {
		return errors.New("cannot add a nil block")
	}

	bc.blocks = append(bc.blocks, b)
	bc.tip = b.Hash

	return nil
}
```

Using it looks like this — build a block with `NewBlock`, using the chain's current `Tip()` and `Height()+1`, then hand it to `AddBlock`:

```go
chain := core.NewBlockchain()

tx := &core.Transaction{ID: []byte("tx-a"), Timestamp: time.Now().Unix()}
next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)

if err := chain.AddBlock(next); err != nil {
	log.Fatal(err)
}
```

This two-step dance — build with `NewBlock`, extend with `AddBlock` — is worth noticing as a pattern, because it stays exactly the same shape all the way through Volume 4's mining loop: `NewBlock` decides *what* the next block contains and computes its hash; `AddBlock` decides *whether* that block is allowed to join the chain. Keeping those two responsibilities in two different functions, rather than one giant "build and add" function, is what makes Chapter 19's validation logic a clean addition rather than a rewrite.

## 7. Drawing the Full Chain — Four Blocks, Linked by Hash Arrows

Building on genesis three more times makes the linking mechanism fully concrete. Here is the loop, and the chain it produces:

```go
chain := core.NewBlockchain()

for i := 1; i <= 3; i++ {
	tx := &core.Transaction{
		ID:        []byte(fmt.Sprintf("tx-%d", i)),
		Timestamp: time.Now().Unix(),
	}
	next := core.NewBlock([]*core.Transaction{tx}, chain.Tip(), chain.Height()+1)
	if err := chain.AddBlock(next); err != nil {
		log.Fatal(err)
	}
}
```

```
 Block 0 (genesis)         Block 1                    Block 2                    Block 3
+------------------+      +------------------+       +------------------+       +------------------+
| Height: 0        |      | Height: 1        |       | Height: 2        |       | Height: 3        |
| PrevHash: 000..0 |      | PrevHash: 7a3f.. |<------| PrevHash: c81e.. |<------| PrevHash: 44b0.. |
| Hash:     7a3f.. |----->| Hash:     c81e.. |------>| Hash:     44b0.. |------>| Hash:     9de2.. |
+------------------+      +------------------+       +------------------+       +------------------+
     tx: coinbase              tx: tx-1                   tx: tx-2                   tx: tx-3
```

Read this diagram exactly the way Section 1 described: Block 1's `PrevHash` (`7a3f..`) is a literal copy of Block 0's `Hash`. Block 2's `PrevHash` (`c81e..`) is a literal copy of Block 1's `Hash`. Block 3's `PrevHash` (`44b0..`) is a literal copy of Block 2's `Hash`. `chain.Tip()` after this loop returns `9de2..` — Block 3's hash — because `AddBlock` moved `bc.tip` forward on every single call.

You can verify this yourself with nothing more than the accessors from Section 4:

```go
blocks := chain.Blocks()
for i := 1; i < len(blocks); i++ {
	match := bytes.Equal(blocks[i].PrevBlockHash, blocks[i-1].Hash)
	fmt.Printf("Block %d links correctly to Block %d: %v\n", blocks[i].Height, blocks[i-1].Height, match)
}
// Block 1 links correctly to Block 0: true
// Block 2 links correctly to Block 1: true
// Block 3 links correctly to Block 2: true
```

This little loop is, in essence, a hand-rolled preview of `ValidateBlock`, which Chapter 19 turns into a proper method with real error messages instead of a printed boolean.

## 8. What AddBlock Does Not Yet Check

It is worth being explicit about what Section 6's `AddBlock` will happily accept right now, because naming the gap clearly is what makes Chapter 19 feel necessary rather than optional. As written, `AddBlock` will accept:

- A block whose `PrevBlockHash` does **not** actually match `bc.Tip()` — silently creating a broken link that Section 7's verification loop would immediately catch, but that `AddBlock` itself does nothing to prevent.
- A block whose stored `Hash` does **not** match what `ComputeHash()` would actually produce from its own contents — meaning a caller could hand `AddBlock` a block with tampered transactions and a stale, no-longer-accurate `Hash`, and it would be added without complaint.
- A block at the wrong `Height` — skipping a number, or repeating one already used.

None of this is a bug in the code written so far; it is an honest, temporary gap, deliberately left open so Chapter 19 has something real and consequential to close. Try it yourself right now, before moving on: construct a block with a `PrevBlockHash` set to some made-up bytes that do not match `chain.Tip()` at all, call `chain.AddBlock(fakeBlock)`, and confirm it succeeds with no error. That successful call is exactly the gap Chapter 19's `ValidateBlock` exists to close.

## 9. Testing core.Blockchain

Three properties are worth pinning down now, so a future refactor cannot quietly break them: a fresh chain always starts with a valid genesis block, `AddBlock` correctly advances the tip, and the full link chain holds across several additions.

```go
// core/blockchain_test.go
package core

import (
	"bytes"
	"testing"
)

func TestNewBlockchain_StartsAtGenesis(t *testing.T) {
	chain := NewBlockchain()

	if chain.Height() != 0 {
		t.Fatalf("expected height 0 for a fresh chain, got %d", chain.Height())
	}

	genesis := chain.LastBlock()
	zeroHash := make([]byte, 32)
	if !bytes.Equal(genesis.PrevBlockHash, zeroHash) {
		t.Fatal("genesis block's PrevBlockHash must be 32 zero bytes")
	}
}

func TestAddBlock_AdvancesTip(t *testing.T) {
	chain := NewBlockchain()

	next := NewBlock([]*Transaction{testTx("tx1")}, chain.Tip(), chain.Height()+1)
	if err := chain.AddBlock(next); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if !bytes.Equal(chain.Tip(), next.Hash) {
		t.Fatal("Tip() must equal the most recently added block's Hash")
	}
	if chain.Height() != 1 {
		t.Fatalf("expected height 1 after one AddBlock call, got %d", chain.Height())
	}
}

func TestChain_EveryBlockLinksToItsPredecessor(t *testing.T) {
	chain := NewBlockchain()

	for i := 1; i <= 5; i++ {
		next := NewBlock([]*Transaction{testTx("tx")}, chain.Tip(), chain.Height()+1)
		if err := chain.AddBlock(next); err != nil {
			t.Fatalf("unexpected error adding block %d: %v", i, err)
		}
	}

	blocks := chain.Blocks()
	for i := 1; i < len(blocks); i++ {
		if !bytes.Equal(blocks[i].PrevBlockHash, blocks[i-1].Hash) {
			t.Fatalf("block %d's PrevBlockHash does not match block %d's Hash", i, i-1)
		}
	}
}
```

Run `go test ./core/...` and all three should pass. `TestChain_EveryBlockLinksToItsPredecessor` is the automated version of Section 7's printed loop — worth keeping permanently in the test suite, since it is exactly the property that "makes a chain a chain" from Section 1, expressed as code instead of prose.

---

## Summary

- A blockchain is nothing more than a sequence of blocks where each one's `PrevBlockHash` equals the actual `Hash` of the block immediately before it.
- The **genesis block** is the first block, at height 0, with `PrevBlockHash` set to 32 zero bytes by convention, because no real predecessor exists — exactly the convention real Bitcoin uses.
- `NewGenesisBlock` needs no new hashing machinery at all; it is `NewBlock`, called with a zeroed `prevBlockHash` and `height` 0.
- `core.Blockchain` holds an unexported `tip` (the newest block's hash) and an unexported `blocks` slice, with read-only accessors (`Tip`, `Height`, `Blocks`, `LastBlock`) as the only way to inspect it from outside `core`.
- `NewBlockchain()` always returns a chain that already contains a real genesis block — a `Blockchain` with zero blocks is never a state this course allows to exist.
- This chapter's `AddBlock` performs only a basic nil check; it does not yet verify a block's hash or its link to the previous block — that gap is deliberate, and Chapter 19 closes it with `ValidateBlock`.
- Drawing (or printing) four or more linked blocks side by side, with arrows from each `PrevBlockHash` back to the previous block's `Hash`, makes the entire chaining mechanism visually obvious in a way the underlying equation alone does not.

---

## Exercises

### Easy

1. Using `NewBlockchain` and a loop of `NewBlock`/`AddBlock` calls, build a chain of 6 blocks and print each block's height alongside the first 8 hex characters of its `Hash` and `PrevBlockHash`. Confirm by eye that each block's `PrevBlockHash` prefix matches the previous block's `Hash` prefix.
2. Explain, in 2-3 sentences, why `NewGenesisBlock` takes a `coinbase *Transaction` parameter instead of building a hardcoded, empty transaction list itself.
3. `Blockchain.Height()` reads `bc.blocks[len(bc.blocks)-1].Height` rather than simply returning `len(bc.blocks) - 1`. Both would give the same answer today. Explain one reason the first version is a safer choice to keep, even though it looks more roundabout.

### Medium

4. Write a function `VerifyLinks(bc *core.Blockchain) []int64` that returns the heights of every block whose `PrevBlockHash` does *not* match its predecessor's `Hash`. Test it against a normal, correctly-built chain (expect an empty result) and against a chain where you manually corrupt one block's `PrevBlockHash` after building it (expect that block's height in the result).
5. Section 8 lists three checks `AddBlock` does not yet perform. For each of the three, write a small Go test that constructs a deliberately broken block (wrong `PrevBlockHash`, stale `Hash`, or wrong `Height`) and calls `AddBlock` with it, asserting that it currently succeeds without error. Label these tests clearly as documenting a temporary gap, to be revisited in Chapter 19.
6. `Blockchain.blocks` is a plain Go slice. Explain what would happen, concretely, if two goroutines called `AddBlock` on the same `*Blockchain` at the same time (consider Chapter 05's discussion of race conditions). Propose one concrete fix using tools already introduced in this course.

### Hard

7. Redraw Section 7's four-block diagram to show what happens when someone attempts to insert a fabricated "Block 1.5" between Block 1 and Block 2, using real (fabricated but internally consistent) hash values of your own choosing. Show explicitly which existing blocks' fields would need to change to make the fabricated insertion link in without `VerifyLinks` (Exercise 4) noticing anything wrong.
8. Implement a `Fork() *core.Blockchain` method that returns a deep copy of a `Blockchain` (a genuinely separate `blocks` slice and blocks, not sharing underlying memory with the original) so that adding blocks to the copy never affects the original. Write a test proving the two chains diverge independently after the fork.
9. Research (or reason from Chapter 16's Bitcoin comparison) what value Bitcoin's genesis block actually uses for fields that GoChain's genesis block also has to fill in somehow — specifically, its `Timestamp` and its single coinbase transaction's contents. Write a short comparison (250-350 words) of how "real" or arbitrary you believe GoChain's current placeholder genesis coinbase transaction is, and propose one concrete, realistic value you would give it once Volume 5's real transactions exist.
