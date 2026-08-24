# Chapter 26: Difficulty Adjustment

Chapter 25 wired proof of work into `core.Blockchain.MineBlock` with one honest shortcut: `defaultDifficultyBits = 16`, a constant, never changing, no matter what. That was fine for mining a single demonstration block on one laptop. It falls apart the moment GoChain becomes a real, ongoing network, because the thing that determines how long mining a block actually takes — the total hashing power pointed at the chain — is never fixed in reality. This chapter replaces the constant with a real algorithm, modeled closely on Bitcoin's own, that watches how fast recent blocks have actually been arriving and nudges difficulty up or down to bring block times back toward a steady target.

## Table of Contents

1. [Why a Fixed Difficulty Breaks](#1-why-a-fixed-difficulty-breaks)
2. [Target Block Time and the Adjustment Window](#2-target-block-time-and-the-adjustment-window)
3. [The Core Idea: Compare Actual Time to Expected Time](#3-the-core-idea-compare-actual-time-to-expected-time)
4. [Recording Each Block's Own Difficulty](#4-recording-each-blocks-own-difficulty)
5. [Implementing NextDifficultyBits](#5-implementing-nextdifficultybits)
6. [Clamping: Preventing Wild Swings](#6-clamping-preventing-wild-swings)
7. [Wiring Adjustment Into MineBlock](#7-wiring-adjustment-into-mineblock)
8. [A Full Worked Numeric Example](#8-a-full-worked-numeric-example)
9. [Testing Difficulty Adjustment](#9-testing-difficulty-adjustment)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why a Fixed Difficulty Breaks

Picture a toll bridge with a single, fixed rule: "it takes exactly one hour to cross, no matter what." That rule might work perfectly on a quiet Tuesday morning with three cars. It falls apart completely the moment a thousand cars show up at once — either the bridge lets everyone through in five minutes (too fast, the "one hour" promise is meaningless) or it becomes so congested that crossing actually takes six hours (too slow, worse than promised). A fixed rule that ignores how much traffic is actually showing up cannot hold its promise under changing conditions.

GoChain's mining difficulty has exactly this problem. Recall from Chapter 24, Section 3 that the expected number of hash attempts to solve a block is `2^difficultyBits`, and expected time is `attempts ÷ hashRate`, where `hashRate` is the *combined* speed of every miner currently pointed at the chain. `difficultyBits` is a constant Chapter 25 fixed at 16. `hashRate` is not a constant at all — it changes every time a miner joins, upgrades their hardware, gets bored and quits, or a large mining operation switches to a different chain entirely. Hold `difficultyBits` fixed while `hashRate` moves, and expected block time necessarily moves too:

```
expected block time = 2^difficultyBits / hashRate

  hashRate doubles (more miners join)  -> expected block time HALVES
  hashRate quadruples                  -> expected block time QUARTERS
  hashRate drops by 10x (miners leave) -> expected block time grows 10x
```

Concretely: if GoChain's early network mines a block roughly every 10 seconds with 16 difficulty bits and a hundred hobbyist miners, and then a single large operation joins with a hundred times more combined hashing power, block times would collapse toward roughly one-tenth of a second — far too fast for transactions to meaningfully settle, for blocks to even propagate across the network before the next one arrives (a problem Volume 7 makes very concrete), and for anything depending on "roughly one block every N seconds" (an exchange counting confirmations, a game with an in-block clock) to behave predictably. The reverse is just as real: if enough miners abandon the chain, block times could stretch from seconds into minutes or hours, and the chain would grind to a crawl exactly when people most need it to keep working.

**Difficulty adjustment**, defined now precisely: an algorithm that periodically recalculates `difficultyBits` based on how fast blocks have *actually* been arriving recently, so that expected block time is pulled back toward a fixed target no matter how the network's total hash power changes. This is not a cosmetic feature — every production proof-of-work chain, Bitcoin included, has one, for exactly the reasons above.

---

## 2. Target Block Time and the Adjustment Window

Two numbers anchor the whole algorithm, and it's worth being precise about both before writing any Go:

- **Target block time** — the block interval GoChain is trying to hold steady. This course uses **10 seconds** (small enough to experiment with quickly on a laptop; Bitcoin's own target is 600 seconds, ten minutes).
- **Adjustment interval** — how many blocks pass between each difficulty recalculation. GoChain uses **10 blocks**. Bitcoin's real interval is 2016 blocks (roughly two weeks at a 10-minute target); GoChain's much smaller interval is a deliberate teaching choice, so a worked example (Section 8) and a test suite (Section 9) don't require mining thousands of blocks to see an adjustment happen.

```go
package consensus

const (
	// TargetBlockTimeSecs is the block interval GoChain tries to hold
	// steady, in seconds. Every difficulty adjustment is measured against
	// this one number.
	TargetBlockTimeSecs = 10

	// AdjustmentInterval is how many blocks pass between each difficulty
	// recalculation. Bitcoin's real interval is 2016 blocks; GoChain uses
	// a much shorter one so adjustments happen fast enough to observe
	// and test without mining a huge number of blocks first.
	AdjustmentInterval = 10
)
```

Multiplying these two together gives the **target timespan**: the total wall-clock time an entire adjustment window (10 blocks) is *supposed* to take if difficulty is well-calibrated — `10 blocks × 10 seconds = 100 seconds`. That number, compared against how long the last 10 blocks *actually* took, is the entire signal the algorithm reacts to.

An important design question worth answering explicitly: why not recompute difficulty after *every single block*, instead of waiting for a whole window of ten? Recall Chapter 24, Section 5's key insight: proof of work only guarantees an *expected* number of attempts, not a fixed one — any single block, mined at perfectly correct difficulty, might get lucky and solve in a tenth of the expected time, or unlucky and take five times longer, purely by chance. Reacting to one block's time would mean reacting mostly to random noise, not to any real change in network hash power — the mining equivalent of slamming the brakes because one car ahead of you tapped its brake lights. Averaging over a whole window of blocks smooths that noise out, so the signal the algorithm reacts to is actually "hash power seems to have changed," not "we got unlucky once."

---

## 3. The Core Idea: Compare Actual Time to Expected Time

Here is the whole algorithm in one sentence: **look at how long the last `AdjustmentInterval` blocks actually took, compare that to how long they were supposed to take, and scale difficulty by the ratio.**

```
actualTimespan  = timestamp of the LAST block in the window
                   - timestamp of the FIRST block in the window

targetTimespan  = AdjustmentInterval * TargetBlockTimeSecs
                   (what that same window SHOULD have taken)

  actualTimespan <  targetTimespan  -->  blocks arrived TOO FAST
                                          -->  network hash power went UP
                                          -->  INCREASE difficulty

  actualTimespan >  targetTimespan  -->  blocks arrived TOO SLOW
                                          -->  network hash power went DOWN
                                          -->  DECREASE difficulty

  actualTimespan == targetTimespan  -->  right on target, no change needed
```

Chapter 25's `targetFromDifficulty` already established the relationship between `difficultyBits` and expected mining time: `expected attempts = 2^difficultyBits`, and since attempts scale directly with time (at a roughly constant hash rate), **expected time is also proportional to `2^difficultyBits`.** That means a change in difficulty of exactly one bit doubles or halves expected time — which gives us a clean, exact way to convert a timing ratio into a bits adjustment, using `log2`:

```
targetTimespan / actualTimespan  =  the factor by which time was OFF

bitsAdjustment = round( log2( targetTimespan / actualTimespan ) )

newDifficultyBits = currentDifficultyBits + bitsAdjustment
```

Walk through why the direction is correct: if blocks came in *twice as fast* as intended, `actualTimespan` is half of `targetTimespan`, so `targetTimespan / actualTimespan = 2`, and `log2(2) = 1` — difficulty goes up by exactly 1 bit, which (by the doubling relationship above) is exactly enough to double expected time back toward the target. If blocks came in *twice as slow*, the ratio is `0.5`, `log2(0.5) = -1`, and difficulty drops by 1 bit, halving expected time. The formula is symmetric and self-correcting in both directions.

---

## 4. Recording Each Block's Own Difficulty

Before any of this math can run, there's a gap to close: Chapter 16's `Block` struct has no field remembering *which* difficulty a given block was actually mined at. That was fine when every block shared the same hardcoded constant — there was nothing block-specific to remember. Now that difficulty changes over time, a node re-validating an old block (or computing the next adjustment) needs to know exactly what difficulty applied to each block in the window, not just the current one. We close the gap with one new field:

```go
package core

// Block gains one new field this chapter: Bits. Every earlier field
// (Chapter 16) is unchanged.
type Block struct {
	Height        int64
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         uint64
	MerkleRoot    []byte

	// Bits is the difficulty -- leading zero bits required of this
	// block's hash -- that this specific block was mined and validated
	// at. Chapters 24-25 treated this as a fixed constant everywhere;
	// starting this chapter, it varies block by block, so each block
	// must carry its own record of what applied to it.
	Bits int
}
```

This mirrors real Bitcoin block headers directly — recall Chapter 16, Section 4's comparison table, which already noted Bitcoin's header carries a `Bits` field encoding difficulty, with a note that "GoChain adds an equivalent... starting in Chapter 24." This chapter is where that promise is finally kept. `Validate()` (Chapter 25) needs a small adjustment to match: instead of being told a difficulty from the outside, `consensus.NewProofOfWork` should read `block.Bits` directly when reconstructing the target for validation, so that validating an old block always uses *that block's own* recorded difficulty, not whatever difficulty is currently active on the tip of the chain.

```go
// NewProofOfWorkForValidation builds a ProofOfWork using a block's own
// recorded Bits field, for re-checking an existing block's proof --
// as opposed to NewProofOfWork (Chapter 25), which takes an explicit
// difficulty for a block about to be MINED, before Bits is even set.
func NewProofOfWorkForValidation(b *core.Block) *ProofOfWork {
	return NewProofOfWork(b, b.Bits)
}
```

---

## 5. Implementing NextDifficultyBits

With `Bits` recorded on every block, the timing-ratio formula from Section 3 becomes real, testable Go code:

```go
package consensus

import "math"

// maxAdjustmentFactor caps how much a single adjustment window can change
// difficulty, in either direction -- Section 6 explains exactly why this
// clamp exists. 4 mirrors Bitcoin's own real retarget clamp.
const maxAdjustmentFactor = 4.0

// NextDifficultyBits computes the difficulty bits the NEXT adjustment
// window's blocks should be mined at, given the difficulty the just-
// completed window was mined at (currentBits) and the real Unix
// timestamps of that window's first and last blocks.
func NextDifficultyBits(currentBits int, firstTimestamp, lastTimestamp int64) int {
	targetTimespan := float64(AdjustmentInterval * TargetBlockTimeSecs)
	actualTimespan := float64(lastTimestamp - firstTimestamp)

	// Guard against a zero or negative actual timespan (possible only
	// with corrupted or adversarial timestamps -- real block timestamps
	// always advance). Treating it as "instantaneous" and clamping below
	// is safer than dividing by zero or by a negative number.
	if actualTimespan <= 0 {
		actualTimespan = 1
	}

	ratio := targetTimespan / actualTimespan
	ratio = clampRatio(ratio)

	// log2(ratio) tells us exactly how many bits to add: doubling
	// expected time needs +1 bit, halving needs -1 bit (Section 3).
	// math.Round turns a fractional log into the single nearest whole
	// bit adjustment -- difficultyBits is always an integer.
	bitsAdjustment := int(math.Round(math.Log2(ratio)))

	newBits := currentBits + bitsAdjustment
	if newBits < 1 {
		newBits = 1 // never let difficulty collapse to "anything goes"
	}
	return newBits
}

// clampRatio restricts ratio to [1/maxAdjustmentFactor, maxAdjustmentFactor]
// -- see Section 6 for why this matters.
func clampRatio(ratio float64) float64 {
	if ratio > maxAdjustmentFactor {
		return maxAdjustmentFactor
	}
	if ratio < 1.0/maxAdjustmentFactor {
		return 1.0 / maxAdjustmentFactor
	}
	return ratio
}
```

`NextDifficultyBits` takes exactly the three numbers Section 3's diagram needs — the difficulty the window was mined at, and the two timestamps bracketing it — and returns a single new difficulty. `targetTimespan` is `AdjustmentInterval * TargetBlockTimeSecs` from Section 2 (`10 × 10 = 100` seconds with GoChain's constants). `ratio` is `targetTimespan / actualTimespan`, matching Section 3's formula exactly (note this is the *inverse* of "actual over target" — deliberately, so that `log2` of it gives the right sign of adjustment directly, without an extra negation step). `clampRatio` is Section 6's safety net, applied before the logarithm. The final `newBits < 1` guard is a floor: difficulty could never sensibly reach zero or negative bits (that would mean *any* hash at all qualifies, defeating the entire point of proof of work), so we refuse to go there no matter what the math suggests.

---

## 6. Clamping: Preventing Wild Swings

Without a clamp, one anomalous adjustment window could send difficulty flying to an extreme in a single jump — imagine a network split (Volume 7's territory) that briefly halved observed hash power, or a handful of blocks that happened to solve unusually fast purely by chance despite Section 2's window-averaging. An unclamped algorithm would swing difficulty by however large that one window's ratio happened to be, and an overcorrection in one direction sets up an overcorrection in the other direction next window — a feedback loop of instability that undermines the whole point of trying to hold block time *steady*.

Bitcoin's real difficulty adjustment guards against exactly this with a hard clamp: no single retarget is ever allowed to change the target by more than a factor of 4, in either direction. GoChain's `clampRatio` (Section 5) copies this exact number. In bits terms, this clamp has a clean, checkable consequence:

```
maxAdjustmentFactor = 4

log2(4)   =  2   -->  difficulty can increase by AT MOST 2 bits per window
log2(1/4) = -2   -->  difficulty can decrease by AT MOST 2 bits per window
```

No matter how extreme a single window's actual timing turns out to be — blocks arriving a thousand times too fast, or effectively never arriving at all — `NextDifficultyBits` will never move difficulty by more than 2 bits in one adjustment. If hash power genuinely has changed by more than a factor of 4, the *next* window's adjustment simply continues correcting further in the same direction, converging on the right difficulty over a few windows rather than lurching there (or past it) in one.

---

## 7. Wiring Adjustment Into MineBlock

`core.Blockchain.MineBlock` (Chapter 25) needs one change: instead of always using `defaultDifficultyBits`, it looks up the *current* difficulty by checking whether the block about to be mined lands on an adjustment boundary.

```go
package core

import (
	"time"

	"github.com/you/gochain/consensus"
)

// currentDifficultyBits decides what difficulty the NEXT block (the one
// about to be mined) should use. Most blocks simply inherit the last
// block's difficulty unchanged; only once every AdjustmentInterval
// blocks does a real recalculation happen.
func (bc *Blockchain) currentDifficultyBits() int {
	lastBlock := bc.blocks[len(bc.blocks)-1]
	nextHeight := lastBlock.Height + 1

	// Adjust only at each interval boundary -- e.g. heights 10, 20, 30...
	// with AdjustmentInterval = 10. Height 0 (genesis) is never adjusted;
	// it starts the chain at a fixed, hardcoded initial difficulty.
	if nextHeight == 0 || nextHeight%consensus.AdjustmentInterval != 0 {
		return lastBlock.Bits
	}

	// windowStart is the block AdjustmentInterval positions back --
	// the first block of the window that just completed.
	windowStart := bc.blocks[len(bc.blocks)-consensus.AdjustmentInterval]

	return consensus.NextDifficultyBits(lastBlock.Bits, windowStart.Timestamp, lastBlock.Timestamp)
}

// MineBlock now computes difficulty dynamically instead of using a fixed
// constant -- everything else is unchanged from Chapter 25.
func (bc *Blockchain) MineBlock(transactions []*Transaction) *Block {
	prevBlock := bc.blocks[len(bc.blocks)-1]
	newHeight := prevBlock.Height + 1

	block := NewBlock(transactions, prevBlock.Hash, newHeight)
	block.Timestamp = time.Now().Unix()
	block.Bits = bc.currentDifficultyBits() // <-- the only new line of substance

	pow := consensus.NewProofOfWork(block, block.Bits)
	nonce, hash := pow.Run()

	block.Nonce = nonce
	block.Hash = hash

	bc.blocks = append(bc.blocks, block)
	bc.tip = block.Hash

	return block
}
```

`currentDifficultyBits` is the decision point: for every height that isn't a multiple of `AdjustmentInterval`, the new block simply inherits whatever difficulty the previous block used — no recalculation, no change. Only when `nextHeight` lands exactly on a window boundary does it gather the two timestamps bracketing the just-completed window and hand them to `NextDifficultyBits`. `MineBlock` itself changes by exactly one meaningful line: `block.Bits` is now computed dynamically instead of read from a constant, and everything downstream (constructing the `ProofOfWork`, running the search, storing the result) is unchanged from Chapter 25 — proof that this chapter's entire contribution is contained in *deciding what difficulty to use*, not in *how mining or validation actually work*.

---

## 8. A Full Worked Numeric Example

Let's trace two concrete adjustment windows by hand, using GoChain's real constants: `TargetBlockTimeSecs = 10`, `AdjustmentInterval = 10`, so `targetTimespan = 100` seconds.

**Scenario A — a burst of new miners joins, blocks come in much faster than intended.**

Suppose the chain has been running at `currentBits = 20`. A large mining operation joins partway through, and the next adjustment window's 10 blocks — which should have taken about 100 seconds — actually complete in just 25 seconds (their combined timestamps: window start at `t = 1,000,000`, window end at `t = 1,000,025`).

```
Step 1 -- gather the numbers:
  targetTimespan = 10 blocks * 10 sec = 100 seconds
  actualTimespan = 1,000,025 - 1,000,000 = 25 seconds

Step 2 -- compute the raw ratio:
  ratio = targetTimespan / actualTimespan = 100 / 25 = 4.0

Step 3 -- clamp the ratio:
  maxAdjustmentFactor = 4, so 4.0 is exactly AT the clamp boundary --
  no clamping needed, ratio stays 4.0

Step 4 -- convert to a bits adjustment:
  bitsAdjustment = round(log2(4.0)) = round(2.0) = 2

Step 5 -- apply it:
  newBits = currentBits + bitsAdjustment = 20 + 2 = 22
```

Difficulty jumps from 20 to 22 bits. Sanity-check this against Chapter 24's table: going from 20 to 22 bits multiplies expected attempts (and therefore expected time) by `2^2 = 4` — exactly canceling out the 4x-too-fast timing this window measured. The very next window, mined at 22 bits with the same hash power, should land almost exactly back on the 100-second target.

**Scenario B — some miners leave, blocks slow down.**

Now suppose, a few windows later, the chain is at `currentBits = 22`, and enough miners have since dropped off that the next window's 10 blocks take 200 seconds instead of the target 100 (window start `t = 2,000,000`, window end `t = 2,000,200`).

```
Step 1 -- gather the numbers:
  targetTimespan = 100 seconds (unchanged, it's always AdjustmentInterval * TargetBlockTimeSecs)
  actualTimespan = 2,000,200 - 2,000,000 = 200 seconds

Step 2 -- compute the raw ratio:
  ratio = 100 / 200 = 0.5

Step 3 -- clamp the ratio:
  0.5 is well within [0.25, 4] -- no clamping needed

Step 4 -- convert to a bits adjustment:
  bitsAdjustment = round(log2(0.5)) = round(-1.0) = -1

Step 5 -- apply it:
  newBits = currentBits + bitsAdjustment = 22 + (-1) = 21
```

Difficulty eases from 22 down to 21 bits — half as many expected attempts, which should roughly halve expected time back down toward 100 seconds for the next window, matching the 2x-too-slow timing this window measured.

**Scenario C — an extreme spike that hits the clamp.**

One more pass to see the clamp actually engage: suppose a window somehow completes in just 5 seconds (an extreme, unrealistic spike, used here purely to exercise the clamp), at `currentBits = 21`.

```
Step 1: targetTimespan = 100, actualTimespan = 5
Step 2: raw ratio = 100 / 5 = 20.0
Step 3: clamp! 20.0 > maxAdjustmentFactor (4), so ratio is capped at 4.0
Step 4: bitsAdjustment = round(log2(4.0)) = 2
Step 5: newBits = 21 + 2 = 23
```

Even though the raw timing ratio was 20x too fast, the clamp restricts this single window's reaction to the same `+2` bits Scenario A saw — never more than a 4x change in one step, exactly as Section 6 promised, no matter how extreme the input.

---

## 9. Testing Difficulty Adjustment

A table-driven test suite covers the ratios worked by hand above, plus the clamp boundary itself:

```go
package consensus

import "testing"

func TestNextDifficultyBits(t *testing.T) {
	const target = AdjustmentInterval * TargetBlockTimeSecs // 100 seconds

	tests := []struct {
		name           string
		currentBits    int
		actualTimespan int64 // seconds the window actually took
		wantBits       int
	}{
		{"exactly on target: no change", 20, target, 20},
		{"4x too fast: hits clamp, +2 bits", 20, target / 4, 22},
		{"2x too fast: +1 bit", 20, target / 2, 21},
		{"2x too slow: -1 bit", 22, target * 2, 21},
		{"20x too fast: clamped to +2 bits, same as 4x", 21, target / 20, 23},
		{"20x too slow: clamped to -2 bits", 21, target * 20, 19},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NextDifficultyBits(tt.currentBits, 0, tt.actualTimespan)
			if got != tt.wantBits {
				t.Errorf("NextDifficultyBits(%d, 0, %d) = %d, want %d",
					tt.currentBits, tt.actualTimespan, got, tt.wantBits)
			}
		})
	}
}

func TestCurrentDifficultyBits_OnlyAdjustsAtWindowBoundaries(t *testing.T) {
	bc := NewBlockchainForTest(t) // test helper: fresh chain, genesis at height 0

	// Mine AdjustmentInterval - 1 blocks (heights 1 through 9). None of
	// these should trigger a recalculation -- each should simply inherit
	// the previous block's Bits unchanged.
	for i := 0; i < AdjustmentInterval-1; i++ {
		bc.MineBlock(nil)
	}

	last := bc.blocks[len(bc.blocks)-1]
	genesis := bc.blocks[0]
	if last.Bits != genesis.Bits {
		t.Fatalf("expected Bits to be unchanged before a window boundary, got %d, want %d",
			last.Bits, genesis.Bits)
	}

	// Mining one more block lands exactly on height 10 -- AdjustmentInterval
	// -- which SHOULD trigger a recalculation (the resulting Bits value
	// depends on how fast this test actually ran, so we only assert that
	// a recalculation was attempted, not its exact outcome).
	bc.MineBlock(nil)
	// (A real test would control timestamps precisely to assert an exact
	// resulting Bits value here -- left as Exercise 5.)
}
```

`TestNextDifficultyBits` exercises the pure function directly: on-target timing produces no change, 2x and 4x deviations produce the expected ±1/±2 bit shifts from Section 8, and deliberately extreme 20x deviations confirm the clamp caps the result at the same ±2 bits a plain 4x deviation would produce — proving the clamp is actually doing its job, not just theoretically present in the code. `TestCurrentDifficultyBits_OnlyAdjustsAtWindowBoundaries` checks the *timing* of adjustment rather than its exact math: mining nine blocks (one short of a full window) should never touch `Bits` at all, confirming Section 7's boundary check (`nextHeight % AdjustmentInterval != 0`) correctly skips every non-boundary height.

---

## Summary

- A fixed difficulty cannot hold block time steady once real miners join and leave, because expected block time is `2^difficultyBits / hashRate`, and `hashRate` is never actually constant on a real network.
- **Difficulty adjustment** periodically recalculates difficulty from how long recent blocks actually took, pulling expected block time back toward a fixed **target block time** (GoChain uses 10 seconds).
- Recalculation happens only once every **adjustment interval** (GoChain uses 10 blocks), not after every single block — averaging over a window filters out the natural per-block randomness Chapter 24 described, reacting to real hash-power trends instead of noise.
- Each `core.Block` now carries its own `Bits` field, recording the difficulty it was actually mined and validated at — necessary now that difficulty is no longer a single shared constant.
- `NextDifficultyBits` computes `ratio = targetTimespan / actualTimespan`, clamps it to `[1/4, 4]`, and applies `bitsAdjustment = round(log2(ratio))` — exploiting the fact that expected mining time is proportional to `2^difficultyBits`, so a timing ratio converts cleanly into a bits adjustment.
- The **clamp** (mirroring Bitcoin's real ±4x-per-window limit, or ±2 bits) prevents one anomalous window from swinging difficulty to an extreme in a single step, avoiding feedback-loop instability.
- `Blockchain.currentDifficultyBits` inherits the previous block's difficulty on every ordinary block, and only calls `NextDifficultyBits` when the next block's height lands exactly on a window boundary.
- A full worked example traced three scenarios by hand — 4x too fast (exactly at the clamp), 2x too slow, and an extreme 20x spike (clamped down to the same ±2-bit limit) — showing the same math a real GoChain node runs internally.

---

## Exercises

### Easy

1. Using the `expected time = 2^difficultyBits / hashRate` relationship from Section 1, explain in your own words why doubling `hashRate` should, left uncorrected, roughly halve block time.
2. Why does GoChain wait for a whole `AdjustmentInterval` of blocks before recalculating difficulty, instead of recalculating after every single block? Reference Chapter 24, Section 5 in your answer.
3. In Scenario B of Section 8, difficulty dropped from 22 to 21 bits. Using Chapter 24's expected-attempts table, explain roughly what fraction of the original expected attempts (and therefore expected time) that one-bit drop represents.

### Medium

4. `clampRatio` clamps the *timing ratio* to `[0.25, 4]` before taking its logarithm. Show, with the actual numbers, that this is equivalent to clamping the final `bitsAdjustment` to the range `[-2, 2]` — i.e., prove no input to `NextDifficultyBits` could ever produce a `bitsAdjustment` of 3 or more (in either direction).
5. Finish `TestCurrentDifficultyBits_OnlyAdjustsAtWindowBoundaries` from Section 9 properly: control the mined blocks' timestamps precisely (rather than letting them default to `time.Now()`) so that the tenth block's window has a known, exact `actualTimespan`, and assert the resulting `Bits` value exactly, using `NextDifficultyBits` to compute the expected value independently in the test.
6. `currentDifficultyBits` inherits `lastBlock.Bits` for every non-boundary height. Suppose a bug caused genesis's `Bits` field to be left at its zero value (`0`) instead of a real starting difficulty. Trace through what would happen to the first ten blocks mined on such a chain, and explain why `NewProofOfWork`'s target-computation math (Chapter 25, Section 2) would behave especially strangely at `difficultyBits = 0`.

### Hard

7. GoChain's `NextDifficultyBits` uses floating-point `math.Log2` and `math.Round`. Real Bitcoin's actual retarget algorithm avoids floating-point arithmetic entirely, instead multiplying and dividing the raw 256-bit `big.Int` target directly by the (clamped) timespan ratio using only integer arithmetic. Research why avoiding floating-point math matters specifically for consensus-critical code that every node must compute identically (hint: think about what could go wrong if two nodes' hardware or Go runtime handled floating-point rounding even slightly differently), and sketch (in comments or pseudocode) how you would rewrite `NextDifficultyBits` to adjust `pow.Target` directly with `big.Int` arithmetic instead of adjusting an integer bits count.
8. Extend `NextDifficultyBits` (or write a small standalone simulation) that models 50 consecutive adjustment windows, where the simulated network's hash power quadruples exactly once (at window 20) and then stays constant. Plot (in a text table is fine) how many windows it takes for difficulty to fully "catch up" to the new hash power within the ±2-bit clamp, and explain the shape of the convergence curve you observe.
9. This chapter's `AdjustmentInterval` of 10 blocks is much shorter than Bitcoin's 2016. Argue, using concrete trade-offs (responsiveness to real hash-power changes vs. vulnerability to a short-lived timestamp manipulation or lucky/unlucky streak dominating a small window), what interval you would actually recommend for a real, production GoChain-derived network, and whether it should differ from Bitcoin's 2016-block choice given GoChain's much shorter 10-second target block time.
