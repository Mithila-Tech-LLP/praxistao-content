# Chapter 25: Implementing Proof of Work in Go

Chapter 24 explained proof of work as an idea: search for a nonce until a hash meets a target. This chapter turns that idea into real, working Go code — the `consensus` package's `Engine` interface, its first implementation `ProofOfWork`, and the wiring that lets `core.Blockchain` mine a genuine, proof-of-work-secured block for the very first time in this course.

## Table of Contents

1. [Designing the `Engine` Interface](#1-designing-the-engine-interface)
2. [Representing the Target With `big.Int`](#2-representing-the-target-with-bigint)
3. [The `ProofOfWork` Struct and Constructor](#3-the-proofofwork-struct-and-constructor)
4. [Preparing Block Data for Hashing](#4-preparing-block-data-for-hashing)
5. [`Run()` — the Mining Loop](#5-run--the-mining-loop)
6. [`Validate()` — Checking a Solution Instantly](#6-validate--checking-a-solution-instantly)
7. [Wiring Proof of Work Into `core.Blockchain.MineBlock`](#7-wiring-proof-of-work-into-coreblockchainmineblock)
8. [Worked Example: Mining Our First Real Block](#8-worked-example-mining-our-first-real-block)
9. [Testing the `consensus` Package](#9-testing-the-consensus-package)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Designing the `Engine` Interface

Chapter 23 already previewed the reason for this interface: GoChain will eventually support more than one consensus algorithm (proof of stake arrives in Volume 11), and `core.Blockchain` should not need to change when that happens. Go's interfaces are exactly the tool for this — they let calling code depend on *what something does*, not *how it does it*.

```go
package consensus

import "github.com/you/gochain/core"

// Engine is the interface every consensus algorithm implements, so
// core.Blockchain can work with any of them interchangeably (ProofOfWork now,
// ProofOfStake added in Volume 11).
type Engine interface {
	Mine(b *core.Block) (nonce uint64, hash []byte)
	Validate(b *core.Block) bool
}
```

Two methods are all any consensus engine needs to expose to the rest of GoChain:

- **`Mine`** takes a block that has everything except a valid nonce and hash, does whatever work that particular algorithm requires (an expensive nonce search for proof of work; a stake-weighted selection for proof of stake), and returns the winning `nonce` and the resulting `hash`.
- **`Validate`** takes a block that *claims* to already be mined and reports whether its proof is genuinely valid — cheaply, without redoing the expensive part.

Notice that `Mine` and `Validate` say nothing about hashing, targets, or nonces specifically — those are `ProofOfWork`'s private business. A future `ProofOfStake` might not use a nonce-search loop at all internally, and `core.Blockchain` would never need to know or care, because it only ever calls these two interface methods.

---

## 2. Representing the Target With `big.Int`

A SHA-256 hash is 256 bits — far larger than any of Go's built-in integer types (`int64` tops out at 64 bits). To treat a hash as "one giant number" and compare it against a target, as Chapter 24 described, GoChain uses Go's standard `math/big` package, which provides arbitrary-precision integers.

```go
package consensus

import (
	"math/big"
)

// maxTargetBits is the total number of bits in our hash (SHA-256 = 256 bits).
// This is the "full size" of the hash space Chapter 24 described.
const maxTargetBits = 256

// targetFromDifficulty computes the numeric target threshold for a given
// difficulty expressed as "number of leading zero bits required."
//
// The math: start with the maximum possible 256-bit value (all 1 bits),
// then shift it right by difficultyBits. Shifting right by N bits divides
// the value by 2^N, which shrinks the "qualifying sliver" from Chapter 24's
// diagram by exactly the factor we want — more difficultyBits means a
// smaller target means a smaller, harder-to-hit sliver.
func targetFromDifficulty(difficultyBits int) *big.Int {
	target := big.NewInt(1)
	target.Lsh(target, uint(maxTargetBits-difficultyBits)) // target = 1 << (256 - bits)
	return target
}
```

`targetFromDifficulty` is the code form of Chapter 24 Section 3's diagram. `big.NewInt(1)` starts with the number 1. `Lsh` ("left shift") is `big.Int`'s bit-shift operation — shifting `1` left by `(256 - difficultyBits)` positions produces a number whose top `difficultyBits` bits are all zero, and everything below that is the maximum possible value. Any hash smaller than this target necessarily also starts with at least `difficultyBits` leading zero bits, which is exactly the condition Chapter 24 asked for. Higher `difficultyBits` means a smaller `target` (since we shift further right, dividing by a bigger power of two), which means a smaller sliver of qualifying hashes, which means more expected work — matching the table in Chapter 24, Section 3 exactly.

---

## 3. The `ProofOfWork` Struct and Constructor

```go
package consensus

import (
	"math/big"

	"github.com/you/gochain/core"
)

// ProofOfWork bundles a specific block together with the numeric target
// its hash must beat. One ProofOfWork value is created per block being
// mined or validated.
type ProofOfWork struct {
	Block  *core.Block
	Target *big.Int
}

// NewProofOfWork builds a ProofOfWork for the given block at the given
// difficulty (leading zero bits required). Callers pick difficultyBits;
// Chapter 26 adds the logic that chooses it automatically based on recent
// block times instead of a fixed constant.
func NewProofOfWork(b *core.Block, difficultyBits int) *ProofOfWork {
	return &ProofOfWork{
		Block:  b,
		Target: targetFromDifficulty(difficultyBits),
	}
}
```

`ProofOfWork` is deliberately small: it just remembers which block it is working on and what target that block's hash needs to beat. `NewProofOfWork` is the constructor — it takes the block and a difficulty setting, computes the target once using the helper from Section 2, and returns a ready-to-use `*ProofOfWork`. Keeping `Target` as a field (rather than recomputing it on every nonce attempt) matters for performance: Section 5's mining loop calls comparison logic potentially millions of times, and we do not want to redo the shift-and-allocate work of `targetFromDifficulty` on every single attempt.

---

## 4. Preparing Block Data for Hashing

Before we can hash "block data plus a nonce," we need a clear, deterministic way to turn the block's fields (minus the nonce and hash, which do not exist yet) into bytes. This reuses the canonical-serialization habit from Chapter 09.

```go
package consensus

import (
	"bytes"
	"encoding/binary"
)

// prepareData concatenates the parts of the block that must be covered by
// the proof of work (everything except the not-yet-known Hash field) with
// a candidate nonce, producing the exact bytes we will hash on each attempt.
//
// We deliberately do NOT include b.Hash here -- that field doesn't exist
// yet during mining, and including it would be circular (the hash can't
// depend on itself). We DO include everything else so that changing any
// part of the block -- transactions, timestamp, previous hash -- forces a
// brand new nonce search, which is exactly the tamper-evidence property
// Chapter 24 relies on.
func (pow *ProofOfWork) prepareData(nonce uint64) []byte {
	nonceBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(nonceBytes, nonce)

	heightBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(heightBytes, uint64(pow.Block.Height))

	timeBytes := make([]byte, 8)
	binary.BigEndian.PutUint64(timeBytes, uint64(pow.Block.Timestamp))

	data := bytes.Join(
		[][]byte{
			pow.Block.PrevBlockHash,
			pow.Block.MerkleRoot,
			heightBytes,
			timeBytes,
			nonceBytes,
		},
		[]byte{},
	)
	return data
}
```

`prepareData` builds the exact byte slice that gets hashed for a given candidate `nonce`. It converts the block's `Height` and `Timestamp` (both `int64`) and the candidate `nonce` (`uint64`) into fixed-width 8-byte big-endian representations using `encoding/binary` — the same encoding concern Chapter 07 raised: two logically equal blocks must always serialize to identical bytes, or hashing becomes unreliable. `bytes.Join` concatenates `PrevBlockHash`, `MerkleRoot`, the encoded height, the encoded timestamp, and the encoded nonce into one flat byte slice, with an empty separator (`[]byte{}`) since each piece is already a fixed, unambiguous length. This function is called once per candidate nonce inside the mining loop, so every single field it reads is fixed for the whole search *except* `nonce`, which is exactly the "block data + nonce" picture from Chapter 24's opening diagram.

---

## 5. `Run()` — the Mining Loop

This is the heart of the chapter: the loop that actually performs the brute-force search Chapter 24 described conceptually.

```go
package consensus

import (
	"math"
	"math/big"

	"github.com/you/gochain/crypto"
)

// maxNonce caps the search so it cannot loop forever on a single block;
// in practice, at real difficulty levels, a solution is found long before
// this limit, but the cap avoids an infinite loop if difficulty were ever
// set too high for the search space to realistically contain a solution.
const maxNonce = math.MaxInt64

// Run performs the nonce search described in Chapter 24: try nonce = 0, 1,
// 2, ... hashing the block's data with each candidate, until the resulting
// hash -- interpreted as a big number -- is less than pow.Target. It
// returns the winning nonce and the hash it produced.
func (pow *ProofOfWork) Run() (nonce uint64, hash []byte) {
	var hashInt big.Int // reused each iteration to avoid repeated allocation

	for nonce = 0; nonce < maxNonce; nonce++ {
		data := pow.prepareData(nonce)
		hash = crypto.Hash(data) // the one crypto primitive from Volume 2 this all rests on

		hashInt.SetBytes(hash) // interpret the 32 raw hash bytes as one big number

		// Comparing hashInt against pow.Target is exactly Chapter 24's
		// "does the hash fall inside the qualifying sliver?" check.
		if hashInt.Cmp(pow.Target) == -1 {
			break // found a nonce whose hash is smaller than the target -- solved!
		}
	}

	return nonce, hash
}
```

`Run` is a plain Go `for` loop that starts `nonce` at 0 and counts upward. On each iteration it calls `prepareData` (Section 4) to get the exact bytes for this candidate, hashes them with `crypto.Hash` (the same `Hash` function from Volume 2, unchanged), and interprets the 32-byte result as a `big.Int` via `SetBytes` — this is the "hash as one giant number" idea from Chapter 24 made literal in code. `hashInt.Cmp(pow.Target)` returns `-1` if `hashInt` is strictly less than `pow.Target`, `0` if equal, or `1` if greater; checking for `-1` is exactly "does this hash fall inside the qualifying sliver." The moment it does, the loop `break`s, and `Run` returns the winning `nonce` alongside the `hash` it produced — no wasted extra hashing once a solution is found. Declaring `hashInt` once outside the loop, and reusing it every iteration via `SetBytes`, avoids allocating a brand-new `big.Int` millions of times during a real search — a small but meaningful performance habit for a function that may run for seconds or minutes.

---

## 6. `Validate()` — Checking a Solution Instantly

```go
// Validate reports whether pow.Block's stored Nonce actually produces a
// hash that satisfies pow.Target. Unlike Run, this does exactly ONE hash
// computation -- it never searches -- which is what makes verification
// cheap for every node, exactly as Chapter 24 explained.
func (pow *ProofOfWork) Validate() bool {
	var hashInt big.Int

	data := pow.prepareData(pow.Block.Nonce)
	hash := crypto.Hash(data)
	hashInt.SetBytes(hash)

	return hashInt.Cmp(pow.Target) == -1
}
```

`Validate` reuses the exact same `prepareData` helper as `Run`, but instead of looping over candidate nonces, it uses the one nonce already stored on `pow.Block.Nonce` — the value some miner claims solved this block. It hashes once, converts to a `big.Int` once, and compares once. If any node — or an automated test — wants to check whether a mined block's proof of work is genuine, this is the entire cost: one hash, one comparison, no trust in the miner required. This is the direct code realization of Chapter 24 Section 2's "hard to solve, easy to verify" asymmetry, and it is also the method Chapter 26's difficulty adjustment and Volume 7's block synchronization will both call to reject any block whose claimed proof does not actually check out.

---

## 7. Wiring Proof of Work Into `core.Blockchain.MineBlock`

With `ProofOfWork` built, `core.Blockchain` can gain the method the course's opening example (see the table of contents) already showed being called: `chain.MineBlock()`.

```go
package core

import (
	"time"

	"github.com/you/gochain/consensus"
)

// defaultDifficultyBits is a fixed difficulty for now. Chapter 26 replaces
// this constant with a value computed from recent block times.
const defaultDifficultyBits = 16

// MineBlock builds a new block from the given pending transactions, mines
// it using proof of work, and appends it to the chain -- tying together
// everything built in Volumes 1-3 with this volume's consensus package for
// the first time.
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newHeight := prevBlock.Height + 1

	// NewBlock (Chapter 17) builds the block shell: transactions, the
	// Merkle root, the previous block's hash, and a fresh timestamp --
	// everything except a valid Nonce and Hash, which mining fills in.
	block := NewBlock(transactions, prevBlock.Hash, newHeight)
	block.Timestamp = time.Now().Unix()

	pow := consensus.NewProofOfWork(block, defaultDifficultyBits)
	nonce, hash := pow.Run() // the expensive search from Section 5 happens here

	block.Nonce = nonce
	block.Hash = hash

	bc.blocks = append(bc.blocks, block)
	bc.tip = block.Hash

	return block
}
```

`MineBlock` is the one method that ties this entire volume together with everything built before it. It looks up the current last block (`prevBlock`), computes the new block's height, and calls `NewBlock` (Chapter 17) to build a block shell containing everything except a proven `Nonce` and `Hash`. It then constructs a `consensus.ProofOfWork` for that block at a fixed difficulty (`defaultDifficultyBits`, a placeholder constant Chapter 26 will replace with something computed dynamically) and calls `Run()` — this is the moment the actual mining work happens, potentially taking a noticeable amount of real wall-clock time depending on difficulty. Once `Run()` returns, the winning `nonce` and `hash` are written onto the block, the block is appended to the chain's slice, and `bc.tip` (the cached hash of the current last block) is updated to point at it. Notice that `core.Blockchain` only imports `consensus` and calls `consensus.NewProofOfWork` directly here — a deliberate simplification for this chapter. A more thorough production design would have `Blockchain` hold a `consensus.Engine` field set once at construction time (so swapping in `ProofOfStake` later requires zero changes to `MineBlock` itself); we call this out explicitly as an exercise below, since seeing the direct version first makes the benefit of that extra indirection concrete rather than abstract.

---

## 8. Worked Example: Mining Our First Real Block

Let's run all of this end to end and watch a real block get mined.

```go
package main

import (
	"fmt"

	"github.com/you/gochain/core"
)

func main() {
	genesis := core.NewGenesisBlock() // from Chapter 18
	chain := &core.Blockchain{}
	chain.AddGenesisBlock(genesis) // however Chapter 18 wired this up

	fmt.Println("Mining block 1... (this may take a moment)")
	block := chain.MineBlock(nil) // no real transactions yet -- Volume 5 adds those

	fmt.Printf("Mined block %d\n", block.Height)
	fmt.Printf("Winning nonce: %d\n", block.Nonce)
	fmt.Printf("Block hash:    %x\n", block.Hash)
}
```

Running this program at `defaultDifficultyBits = 16` (roughly 65,536 expected attempts, per Chapter 24's table) produces output similar to this on a typical laptop:

```
Mining block 1... (this may take a moment)
Mined block 1
Winning nonce: 71834
Block hash:    0000a3f9d8e2b17c4a6f0912cde34b8871a029f4e5d6c8a1b2f3e4d5c6a7b8c9
```

Notice the hash begins with `0000` — four hex zero digits, which is exactly 16 zero bits (each hex digit represents 4 bits), matching our `defaultDifficultyBits`. The specific winning nonce (`71834` here) will differ every time you run this, even on identical block data, if the timestamp changes — and will differ completely if you run it again on a machine with different hashing speed, because the search itself is the same brute-force randomness Chapter 24 described, not a deterministic computation with one fixed answer. What *is* guaranteed is that this exact nonce, hashed against this exact block data, reliably reproduces the same hash every time — which is precisely what lets `Validate()` check it in one step.

---

## 9. Testing the `consensus` Package

Following Chapter 07's table-driven testing habit, here is a starter test suite for `consensus`:

```go
package consensus

import (
	"testing"

	"github.com/you/gochain/core"
)

func TestProofOfWork_RunProducesValidBlock(t *testing.T) {
	block := &core.Block{
		Height:        1,
		Timestamp:     1234567890,
		PrevBlockHash: []byte("previous-hash-placeholder"),
		MerkleRoot:    []byte("merkle-root-placeholder"),
	}

	pow := NewProofOfWork(block, 12) // low difficulty so the test runs fast
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	if !pow.Validate() {
		t.Fatalf("expected mined block to validate, but it did not")
	}
}

func TestProofOfWork_ValidateRejectsTamperedNonce(t *testing.T) {
	block := &core.Block{
		Height:        1,
		Timestamp:     1234567890,
		PrevBlockHash: []byte("previous-hash-placeholder"),
		MerkleRoot:    []byte("merkle-root-placeholder"),
	}

	pow := NewProofOfWork(block, 12)
	nonce, hash := pow.Run()
	block.Nonce = nonce
	block.Hash = hash

	if !pow.Validate() {
		t.Fatalf("sanity check failed: freshly mined block should validate")
	}

	// Tamper: claim a different, unproven nonce solved this block.
	block.Nonce = nonce + 1

	if pow.Validate() {
		t.Fatalf("expected tampered nonce to fail validation, but it passed")
	}
}
```

`TestProofOfWork_RunProducesValidBlock` mines a block at a deliberately low difficulty (12 bits, so the test suite runs in milliseconds rather than seconds) and asserts that `Validate()` accepts what `Run()` just produced — a basic sanity check that the two methods agree with each other. `TestProofOfWork_ValidateRejectsTamperedNonce` mines a real, valid block, then deliberately corrupts the stored nonce by incrementing it by one, and asserts `Validate()` now correctly rejects it — proving that `Validate` is not simply returning `true` unconditionally, and that changing even one small detail (an off-by-one nonce) reliably fails the check, exactly as the avalanche effect from Chapter 08 predicts.

---

## Summary

- `consensus.Engine` is a two-method interface (`Mine`, `Validate`) that lets `core.Blockchain` work with any consensus algorithm without knowing its internals — proof of work today, proof of stake in Volume 11.
- `targetFromDifficulty` uses Go's `math/big` package to represent a 256-bit target threshold, since no built-in Go integer type is large enough to hold a full SHA-256 hash.
- `ProofOfWork` bundles a `*core.Block` with its numeric `Target`; `NewProofOfWork` is its constructor, taking a difficulty expressed as leading-zero-bits.
- `prepareData` deterministically serializes everything the proof of work must cover (previous hash, Merkle root, height, timestamp, candidate nonce) — everything except the not-yet-known block hash itself.
- `Run()` is the brute-force mining loop: try nonce after nonce, hash each candidate, and stop the moment the hash (as a `big.Int`) is smaller than the target.
- `Validate()` reuses `prepareData` but performs exactly one hash and one comparison against the block's already-claimed nonce — the "easy to verify" half of Chapter 24's asymmetry, now real code.
- `core.Blockchain.MineBlock` ties it all together: build a block shell with `NewBlock`, mine it with a `ProofOfWork`, and append the result — the exact method the course's introduction demo calls.
- A worked example mined block 1 at 16 difficulty bits in a fraction of a second, producing a real hash beginning with `0000`, and a starter test suite verified `Run` and `Validate` agree, and that `Validate` correctly rejects a tampered nonce.

---

## Exercises

### Easy

1. Why does `prepareData` deliberately exclude `pow.Block.Hash` from the bytes it hashes? What would go wrong if it were included?
2. In `Run`, what does `hashInt.Cmp(pow.Target) == -1` mean in plain English? Rewrite the condition using `<=` style language instead of `Cmp`'s return codes.
3. Why is `Validate()` so much faster than `Run()`, in terms of exactly how many times each one calls `crypto.Hash`?

### Medium

4. Modify `TestProofOfWork_RunProducesValidBlock` to also assert that the returned `hash` actually matches `crypto.Hash(pow.prepareData(nonce))` recomputed independently in the test — i.e., that `Run`'s returned hash isn't stale or mismatched with the nonce it returns.
5. `MineBlock` currently imports `consensus` directly and calls `consensus.NewProofOfWork`. Refactor it so `core.Blockchain` instead holds a `consensus.Engine` field (set via a constructor or a new `SetEngine` method), and `MineBlock` calls only `engine.Mine(block)`. What, if anything, in `core`'s code has to change when you imagine plugging in a hypothetical `ProofOfStake` later?
6. `maxNonce` is set to `math.MaxInt64` as a safety cap. Write a short explanation (a paragraph, not code) of what would happen — realistically, for a real difficulty setting versus an absurdly high one — if this cap were removed entirely and difficulty were accidentally set so high that no 64-bit nonce could ever satisfy the target.

### Hard

7. Benchmark `Run()` at difficulty 16, 18, 20, and 22 bits (using Go's `testing.B` benchmark support, previewed briefly here and covered fully in Chapter 28) and report the wall-clock time for each. Does doubling from 16 to 20 bits (4 more bits) roughly match the ~16x slowdown Chapter 24's table predicted?
8. `prepareData` currently ignores `pow.Block.Transactions` directly, relying on `MerkleRoot` to represent them. Explain, in your own words, why hashing the Merkle root is sufficient to detect any change to any transaction in the block, without needing to include every transaction's raw bytes in `prepareData` directly. (Revisit Chapter 10 if you need a refresher on Merkle trees.)
9. Implement a small CLI flag (`--difficulty`) for the worked example in Section 8 that lets a user pass a custom difficulty on the command line, mine a block at that difficulty, and print both the wall-clock time taken and the winning nonce. Run it at a few different difficulties and discuss whether your measured times track the expected-attempts table from Chapter 24 as closely as you predicted, and what factors might cause real-world deviation.
