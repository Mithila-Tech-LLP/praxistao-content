# Chapter 77: Proof of Stake, Explained and Implemented

Proof of work makes rewriting history expensive by burning real electricity. Proof of stake makes it expensive a different way: by requiring a participant to put up real financial value as collateral, which they lose if they misbehave. This chapter explains that trade-off conceptually, then builds `consensus.ProofOfStake` as a second, fully working implementation of the same `consensus.Engine` interface `consensus.ProofOfWork` already satisfies — so GoChain can run on either algorithm without a single line of `core` changing.

## Table of Contents

1. [The Core Idea: Skin in the Game Instead of Burned Electricity](#1-the-core-idea-skin-in-the-game-instead-of-burned-electricity)
2. [Validator Selection: Weighted by Stake](#2-validator-selection-weighted-by-stake)
3. [Slashing: What Happens When a Validator Misbehaves](#3-slashing-what-happens-when-a-validator-misbehaves)
4. [Why Proof of Stake Fits Behind the Same `Engine` Interface](#4-why-proof-of-stake-fits-behind-the-same-engine-interface)
5. [Implementing `Validator` and `ProofOfStake`](#5-implementing-validator-and-proofofstake)
6. [Implementing Weighted-Random Selection](#6-implementing-weighted-random-selection)
7. [Implementing `Mine`: Propose and Sign, Not Grind](#7-implementing-mine-propose-and-sign-not-grind)
8. [Implementing `Validate`](#8-implementing-validate)
9. [Implementing `Slash`](#9-implementing-slash)
10. [Testing Selection Fairness and Slashing](#10-testing-selection-fairness-and-slashing)
11. [A Small GoChain Testnet Under PoS](#11-a-small-gochain-testnet-under-pos)
12. [Comparing Block-Time Behavior: PoW vs. PoS](#12-comparing-block-time-behavior-pow-vs-pos)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. The Core Idea: Skin in the Game Instead of Burned Electricity

Recall proof of work's whole security argument from Volume 4: an attacker who wants to rewrite history has to out-race the honest network's *computational effort*, which costs real electricity and hardware, so honesty is cheaper than dishonesty as long as honest participants control most of the hash power. **Proof of stake** reaches for the same outcome — make dishonesty expensive — through a different mechanism entirely: instead of spending electricity to *earn* the right to propose the next block, a participant locks up (**stakes**) a meaningful amount of the chain's own currency as collateral, and that collateral is at risk of being partially or fully destroyed (**slashed**) if they are caught behaving dishonestly.

A useful everyday analogy: think of a security deposit on a rented apartment. The landlord does not need to constantly watch every tenant to trust they will not trash the place — the deposit itself creates the right incentive, because a tenant who damages the apartment loses money they already handed over, whether or not the landlord catches every act of damage in real time. Proof of stake works the same way: a validator does not need to be watched every second, because they have already put money on the line that they stand to lose the moment provable misbehavior is detected.

```
PROOF OF WORK                              PROOF OF STAKE

  earn the right to propose a block          earn the right to propose a block
  by SPENDING COMPUTATION                    by LOCKING UP MONEY (stake)

  cheating costs: the electricity            cheating costs: your staked
  you would have spent mining honestly,       gochips, which get destroyed
  wasted, plus a lost block reward            (slashed) if you are caught

  "security" = honest hash power              "security" = honest staked value
   > dishonest hash power                      > dishonest staked value
```

This is why proof of stake is often described as more energy-efficient than proof of work: nobody is racing to solve a puzzle nobody else can shortcut, so there is no reason to run enormous amounts of specialized hardware continuously. The trade-off is a different kind of trust: proof of stake's security rests on "a majority of staked value is honest," in the same way proof of work's rests on "a majority of hash power is honest" — different currencies of trust, not an unqualified improvement of one over the other. Chapter 80 returns to this trade-off directly when comparing when each is the right engineering choice.

---

## 2. Validator Selection: Weighted by Stake

In proof of work, "who gets to propose the next block" is decided by a race: whoever finds a valid nonce first. Proof of stake instead needs an explicit selection rule, because there is no puzzle to race — anyone could technically propose a block at any time, so the protocol must pick exactly one **validator** (proof of stake's term for a proof-of-work miner: a participant eligible to propose and sign blocks) for each round, in a way every node can independently verify was fair.

The natural rule, and the one this chapter implements, is **weighted-random selection**: each validator's probability of being chosen is proportional to how much they have staked. A validator who has staked twice as much as another is twice as likely to be picked in any given round — mirroring how, in proof of work, a miner with twice the hash power finds valid blocks roughly twice as often.

```
Three validators, weighted by stake:

  Alice   staked 500 gochips   |=====================|  50%
  Ben     staked 300 gochips   |=============|          30%
  Chidi   staked 200 gochips   |========|                20%
                                0                       1000 (total stake)

  A random point in [0, 1000) lands in Alice's range 50% of the time,
  Ben's range 30% of the time, Chidi's range 20% of the time --
  proportional to stake, exactly like a raffle where buying more
  tickets (staking more) improves your odds without guaranteeing a win.
```

The randomness itself must be unpredictable in advance (so validators cannot manipulate who gets picked) but verifiable after the fact (so every honest node agrees the selection was fair). Real systems typically derive this randomness from something already agreed upon and hard to bias, such as a hash of the previous block combined with the current round number — GoChain's implementation below takes exactly this approach, using a `seed []byte` the caller supplies.

---

## 3. Slashing: What Happens When a Validator Misbehaves

**Slashing** is the mechanism that gives staking teeth: if a validator is caught breaking the protocol's rules — proposing two conflicting blocks for the same height (called **equivocation**), signing an invalid block, or being provably offline when it was their turn — a portion (or all) of their staked gochips is destroyed, and they may be removed from the validator set entirely. This is the direct proof-of-stake analogue of a proof-of-work miner wasting real electricity on a block that gets rejected: in both systems, dishonesty (or even just unreliability) has a real, unavoidable cost, not merely a missed opportunity.

```
NORMAL OPERATION                        SLASHING EVENT

  Validator proposes valid block          Validator caught proposing TWO
  Validator's stake: unaffected           different blocks at the same height
  Validator eligible for future rounds    (equivocation)
                                                |
                                                v
                                          Validators.Stake -= penalty
                                          (partially or fully destroyed)
                                                |
                                                v
                                          If stake drops to 0 (or below a
                                          minimum), removed from the
                                          validator set entirely
```

Slashing only works as a deterrent if the *expected* penalty for cheating exceeds the *expected* gain from cheating — the same cost/benefit logic Chapter 76 applied to 51% attacks under proof of work. This is why real proof-of-stake systems set slashing penalties deliberately harsh relative to any plausible one-time gain from double-signing or censoring transactions.

---

## 4. Why Proof of Stake Fits Behind the Same `Engine` Interface

Recall the interface `consensus.ProofOfWork` already implements, unchanged since Chapter 25:

```go
package consensus

type Engine interface {
	Mine(b *core.Block) (nonce uint64, hash []byte)
	Validate(b *core.Block) bool
}
```

`core.Blockchain.MineBlock` and `core.ValidateBlock` (Volume 3 and 4) were both written to depend only on this interface, never on `ProofOfWork` directly — exactly the design decision Chapter 23 flagged as paying off later. Proof of stake's `Mine` does something conceptually different internally (select a validator and sign, rather than grind through nonces), but it returns the exact same shape of data, so every other package in GoChain — `core`, `network`, `api` — can use a `ProofOfStake` engine without being aware anything changed. This is the practical payoff of programming against interfaces: **swapping the entire consensus algorithm is a one-line change** at the point where a `Blockchain` is constructed, not a rewrite of the blockchain itself.

---

## 5. Implementing `Validator` and `ProofOfStake`

Here are the exact types this course's shared contract requires, plus the constructor:

```go
// consensus/pos.go
package consensus

import (
	"crypto/sha256"
	"encoding/binary"
	"errors"

	"github.com/you/gochain/core"
)

// Validator represents one participant eligible to propose and sign blocks
// under proof of stake. Stake is denominated in gochips -- the same currency
// unit transactions move -- so staking is a real, felt financial commitment,
// not a separate reputation score.
type Validator struct {
	Address string
	Stake   int64 // amount of gochips staked
}

// ProofOfStake is a second, complete implementation of consensus.Engine. It
// holds the current validator set and uses it to select a proposer for each
// block instead of racing to solve a computational puzzle.
type ProofOfStake struct {
	Validators []Validator
}

// NewProofOfStake builds a ProofOfStake engine from an initial validator set.
// In a real network this set would be seeded from a genesis configuration
// (Section 11) and updated over time as validators join, leave, or get
// slashed -- this constructor just captures the starting point.
func NewProofOfStake(validators []Validator) *ProofOfStake {
	// Copy the slice so callers cannot mutate our internal validator set by
	// holding onto and modifying the slice they passed in -- the same
	// defensive-copy habit used for core.Block's transaction list.
	vs := make([]Validator, len(validators))
	copy(vs, validators)
	return &ProofOfStake{Validators: vs}
}

// totalStake sums every validator's stake -- the denominator used throughout
// weighted-random selection in Section 6.
func (pos *ProofOfStake) totalStake() int64 {
	var total int64
	for _, v := range pos.Validators {
		total += v.Stake
	}
	return total
}
```

`Validator` and `ProofOfStake` are defined exactly as the shared contract specifies: a validator is just an address and a stake amount, and the engine is just a slice of validators. `NewProofOfStake` copies the incoming slice defensively — a habit already familiar from `core.Block.Transactions` in Volume 3 — so nothing outside this package can silently mutate the validator set out from under the engine. `totalStake` is a small private helper the rest of this file relies on repeatedly.

---

## 6. Implementing Weighted-Random Selection

```go
// SelectValidator picks one validator, weighted by stake, using seed as the
// source of randomness. Passing the same seed always returns the same
// validator -- determinism is essential here, because every honest node
// must independently compute the SAME selection given the same seed, or
// they will disagree about whose block is legitimate.
func (pos *ProofOfStake) SelectValidator(seed []byte) *Validator {
	if len(pos.Validators) == 0 {
		return nil // no validators registered; caller must handle this
	}

	total := pos.totalStake()
	if total == 0 {
		return nil // nobody has staked anything; nobody is eligible
	}

	// Turn the seed into a deterministic, uniformly-distributed number in
	// [0, total) -- this is the "raffle ticket" from Section 2's diagram.
	// Hashing the seed first (rather than using it directly) smooths out
	// any patterns in how seeds are chosen, so small differences in seed
	// input still land uniformly across the full range.
	digest := sha256.Sum256(seed)
	ticket := int64(binary.BigEndian.Uint64(digest[:8])) % total
	if ticket < 0 {
		ticket += total // Uint64 truncation to int64 can go negative; correct it
	}

	// Walk the validator list, accumulating stake ranges, until the ticket
	// falls inside one validator's range -- exactly the raffle diagram from
	// Section 2, walked in code instead of drawn as a bar.
	var cumulative int64
	for i := range pos.Validators {
		cumulative += pos.Validators[i].Stake
		if ticket < cumulative {
			return &pos.Validators[i]
		}
	}

	// Defensive fallback: floating-point-free integer math should make this
	// unreachable, but returning the last validator instead of nil keeps
	// callers safe against a subtle rounding bug rather than crashing them.
	return &pos.Validators[len(pos.Validators)-1]
}
```

`SelectValidator` is the weighted-random raffle from Section 2's diagram, written as code. It first hashes the caller-supplied `seed` with SHA-256 (from `gochain/crypto`'s underlying primitive, Volume 2) to get a well-distributed number, reduces it modulo the total stake to land somewhere in `[0, total)`, and then walks the validator list accumulating stake ranges until it finds which validator's range contains that number — precisely the "which colored segment of the bar does the random point fall into" logic from the diagram. Determinism is the crucial property: every honest node computes `SelectValidator` with the identical seed (typically derived from the previous block's hash, so nobody can predict it far in advance, but everyone can verify it after the fact) and must therefore always agree on the same answer.

---

## 7. Implementing `Mine`: Propose and Sign, Not Grind

```go
// Mine, for ProofOfStake, does not search for a nonce at all -- "mining"
// here is a holdover name from the shared Engine interface, kept so
// core.Blockchain can call it identically regardless of which engine is
// active. What actually happens is: select this round's validator (using
// the previous block's hash as the seed), confirm the caller IS that
// validator, and produce the block's hash by hashing its finalized contents
// -- there is no puzzle to solve, only a signature to produce.
func (pos *ProofOfStake) Mine(b *core.Block) (nonce uint64, hash []byte) {
	seed := b.PrevBlockHash // deterministic, already-agreed-upon randomness source
	proposer := pos.SelectValidator(seed)
	if proposer == nil {
		return 0, nil // no eligible validator; caller should not accept this block
	}

	// Unlike proof of work, the "nonce" carries no puzzle-solving meaning.
	// We repurpose the field to record which validator produced this block,
	// encoded as a simple index, so Validate (Section 8) can cheaply check
	// that the block's declared proposer matches who SelectValidator would
	// have actually chosen for this seed.
	nonce = uint64(pos.indexOf(proposer.Address))

	// The block's hash is computed the same way it always has been --
	// core.Block.ComputeHash() from Volume 3 -- because every downstream
	// consumer (storage, networking, the explorer) already expects a
	// block's Hash field to be its own content hash, regardless of which
	// consensus engine produced it.
	hash = b.ComputeHash()
	return nonce, hash
}

// indexOf returns the position of the validator with the given address in
// pos.Validators, or -1 if not found. A small helper used by Mine and
// Validate to keep the "nonce" field's PoS-specific meaning consistent.
func (pos *ProofOfStake) indexOf(address string) int {
	for i, v := range pos.Validators {
		if v.Address == address {
			return i
		}
	}
	return -1
}
```

`Mine` fulfills the `Engine` interface's signature exactly, but its meaning is different from proof of work's grinding loop, as the contract's comment ("mine here just means propose+sign") promises: it selects this round's validator deterministically via `SelectValidator`, using the previous block's hash as an unpredictable-in-advance-but-verifiable-after-the-fact seed, and repurposes the returned `nonce` to record *which* validator produced the block rather than a solved puzzle number. The block's `hash` is still computed the ordinary way, via `ComputeHash`, so nothing downstream — storage, networking, the explorer — needs to know or care which consensus engine is active.

---

## 8. Implementing `Validate`

```go
// Validate checks that a block was legitimately proposed under this round's
// rules: the validator recorded in the block's Nonce field must be exactly
// who SelectValidator would independently choose given the block's own
// PrevBlockHash as the seed, and the block's stored Hash must match a fresh
// recomputation of its contents -- the same tamper-evidence check every
// engine performs, regardless of how the block was produced.
func (pos *ProofOfStake) Validate(b *core.Block) bool {
	seed := b.PrevBlockHash
	expected := pos.SelectValidator(seed)
	if expected == nil {
		return false // nobody was eligible; this block cannot be legitimate
	}

	claimedIndex := int(b.Nonce)
	if claimedIndex < 0 || claimedIndex >= len(pos.Validators) {
		return false // out-of-range index; clearly not a real validator
	}

	claimed := pos.Validators[claimedIndex]
	if claimed.Address != expected.Address {
		return false // someone other than the rightfully-selected validator proposed this
	}

	// Recompute the hash fresh and compare -- this is the same tamper check
	// core.ValidateBlock already performs; ProofOfStake repeats it here so
	// that Validate alone, called in isolation, is a complete check.
	return string(b.Hash) == string(b.ComputeHash())
}
```

`Validate` mirrors `SelectValidator`'s logic exactly, because that is precisely what makes selection *verifiable*: any node, given only the block itself, can recompute who *should* have proposed it and compare that against who the block claims did. This closes the loop opened in Section 2 — validator selection is unpredictable enough that nobody can game it in advance, yet fully checkable by anyone after the fact, exactly like proof of work's puzzle is hard to solve but trivial to verify.

---

## 9. Implementing `Slash`

```go
// Slash reduces the named validator's stake by penalty gochips, as
// punishment for provable misbehavior (equivocation, invalid block
// proposals, or extended unavailability -- detecting WHICH of these
// occurred is a networking-and-evidence concern outside this function's
// job; Slash's job is purely to apply the punishment once misbehavior has
// already been established). If the resulting stake would drop to zero or
// below, the validator is removed from the set entirely -- they lose their
// eligibility to be selected in any future round.
func (pos *ProofOfStake) Slash(address string, penalty int64) error {
	if penalty <= 0 {
		return errors.New("consensus: slash penalty must be positive")
	}

	idx := pos.indexOf(address)
	if idx == -1 {
		return errors.New("consensus: no such validator: " + address)
	}

	pos.Validators[idx].Stake -= penalty

	if pos.Validators[idx].Stake <= 0 {
		// Remove the validator entirely: swap with the last element and
		// truncate. Order among validators does not matter for
		// SelectValidator's correctness, so this O(1) removal is safe.
		last := len(pos.Validators) - 1
		pos.Validators[idx] = pos.Validators[last]
		pos.Validators = pos.Validators[:last]
	}

	return nil
}
```

`Slash` is the punishment half of the incentive structure described in Section 3. It looks up the offending validator by address, subtracts the penalty from their staked amount, and — if that brings their stake to zero or below — removes them from the validator set entirely using Go's standard "swap with the last element, then truncate" trick for O(1) removal from a slice, since validator order has no effect on `SelectValidator`'s correctness. Returning an error for an unknown address or a non-positive penalty follows the same defensive style as every other GoChain package from earlier volumes: fail loudly and immediately on nonsensical input, rather than silently doing something unintended.

---

## 10. Testing Selection Fairness and Slashing

```go
// consensus/pos_test.go
package consensus

import (
	"fmt"
	"testing"
)

// TestSelectValidator_WeightedFairness checks that, across many rounds with
// varying seeds, each validator is chosen roughly in proportion to their
// stake -- not exactly (this is randomness, not arithmetic), but close
// enough that a validator with 50% of total stake is picked close to half
// the time, not, say, an equal third of the time regardless of stake.
func TestSelectValidator_WeightedFairness(t *testing.T) {
	pos := NewProofOfStake([]Validator{
		{Address: "alice", Stake: 500},
		{Address: "ben", Stake: 300},
		{Address: "chidi", Stake: 200},
	})

	counts := map[string]int{}
	const rounds = 10000
	for i := 0; i < rounds; i++ {
		seed := []byte(fmt.Sprintf("round-%d", i))
		v := pos.SelectValidator(seed)
		counts[v.Address]++
	}

	// Allow generous tolerance (+/- 5 percentage points) since this is a
	// statistical property, not an exact one -- a flaky test that fails on
	// ordinary random variation is worse than no test at all.
	assertRoughly(t, counts["alice"], rounds, 0.50, 0.05)
	assertRoughly(t, counts["ben"], rounds, 0.30, 0.05)
	assertRoughly(t, counts["chidi"], rounds, 0.20, 0.05)
}

func assertRoughly(t *testing.T, count, total int, expectedFraction, tolerance float64) {
	t.Helper()
	actual := float64(count) / float64(total)
	if actual < expectedFraction-tolerance || actual > expectedFraction+tolerance {
		t.Errorf("expected fraction near %.2f, got %.2f (%d/%d)", expectedFraction, actual, count, total)
	}
}

// TestSlash_RemovesValidatorAtZeroStake confirms that slashing a validator's
// entire stake removes them from future eligibility -- the concrete,
// testable consequence of Section 3's "cheating has a real cost" claim.
func TestSlash_RemovesValidatorAtZeroStake(t *testing.T) {
	pos := NewProofOfStake([]Validator{
		{Address: "alice", Stake: 100},
		{Address: "ben", Stake: 300},
	})

	if err := pos.Slash("alice", 100); err != nil {
		t.Fatalf("Slash returned unexpected error: %v", err)
	}

	if len(pos.Validators) != 1 {
		t.Fatalf("expected 1 validator remaining, got %d", len(pos.Validators))
	}
	if pos.Validators[0].Address != "ben" {
		t.Fatalf("expected ben to remain, got %s", pos.Validators[0].Address)
	}

	// A fully-slashed validator must never be selectable again.
	for i := 0; i < 100; i++ {
		seed := []byte(fmt.Sprintf("post-slash-%d", i))
		if v := pos.SelectValidator(seed); v.Address == "alice" {
			t.Fatalf("slashed validator alice was still selected")
		}
	}
}

// TestSlash_PartialPenaltyKeepsValidatorEligible confirms a partial slash
// reduces stake without removing eligibility entirely -- a lighter penalty
// for a lighter offense, distinct from the full removal case above.
func TestSlash_PartialPenaltyKeepsValidatorEligible(t *testing.T) {
	pos := NewProofOfStake([]Validator{
		{Address: "alice", Stake: 500},
		{Address: "ben", Stake: 300},
	})

	if err := pos.Slash("alice", 100); err != nil {
		t.Fatalf("Slash returned unexpected error: %v", err)
	}

	if len(pos.Validators) != 2 {
		t.Fatalf("expected both validators still present after partial slash")
	}
	if pos.Validators[pos.indexOf("alice")].Stake != 400 {
		t.Fatalf("expected alice's stake reduced to 400, got %d", pos.Validators[pos.indexOf("alice")].Stake)
	}
}
```

`TestSelectValidator_WeightedFairness` runs ten thousand selection rounds with varying seeds and checks that each validator's observed selection frequency lands within five percentage points of their stake's expected share — a statistical test with generous tolerance, appropriate for verifying a randomized process rather than an exact calculation. `TestSlash_RemovesValidatorAtZeroStake` slashes a validator's entire stake and confirms two things: the validator set shrinks to exclude them, and — just as important — running selection a hundred more times never picks them again, proving `Slash` has a real, lasting consequence rather than a cosmetic one. `TestSlash_PartialPenaltyKeepsValidatorEligible` checks the gentler case: a partial penalty reduces stake without full removal, matching the "partial slash for a partial offense" behavior real proof-of-stake systems implement (a validator briefly offline is typically penalized more lightly than one caught equivocating).

---

## 11. A Small GoChain Testnet Under PoS

Wiring `ProofOfStake` into a running node requires no changes to `core.Blockchain` at all — only a different engine passed in at construction time, exactly as Section 4 promised:

```go
// cmd/postestnet/main.go
package main

import (
	"log"

	"github.com/you/gochain/consensus"
	"github.com/you/gochain/core"
)

func main() {
	// Seed the validator set: three participants who have staked gochips,
	// standing in for a real genesis configuration where stake would come
	// from actual on-chain deposits.
	engine := consensus.NewProofOfStake([]consensus.Validator{
		{Address: "gochain1alice...", Stake: 5000},
		{Address: "gochain1ben...", Stake: 3000},
		{Address: "gochain1chidi...", Stake: 2000},
	})

	chain, err := core.OpenBlockchain("./data/pos-testnet", "gochain1alice...")
	if err != nil {
		log.Fatalf("open chain: %v", err)
	}
	defer chain.Close()

	// core.Blockchain.MineBlock takes any consensus.Engine -- it has no idea
	// (and no need to know) that this one selects validators by stake
	// instead of grinding through nonces.
	for i := 0; i < 10; i++ {
		block := chain.MineBlock(engine)
		log.Printf("PoS block %d proposed, proposer index (nonce) = %d, hash = %x",
			block.Height, block.Nonce, block.Hash[:4])
	}
}
```

This program is almost identical in shape to `startHonestNode` from Chapter 76, with one line changed: `consensus.NewProofOfStake(...)` in place of `consensus.NewProofOfWork(...)`. Everything else — `core.OpenBlockchain`, `chain.MineBlock`, the resulting `core.Block` values — is exactly the code you already built, exercised through a completely different consensus mechanism underneath it, which is the entire point of designing to the `Engine` interface in the first place.

---

## 12. Comparing Block-Time Behavior: PoW vs. PoS

Run this testnet alongside the proof-of-work honest network from Chapter 76 (or re-run Chapter 25's original miner) and compare the timing behavior directly:

```
PROOF OF WORK (low difficulty, Chapter 76's lab settings)

  block 14  mined at t=0.00s   (had to search for a valid nonce)
  block 15  mined at t=0.61s   (variance: proof-of-work solve time is random)
  block 16  mined at t=1.02s
  block 17  mined at t=2.31s   <- got unlucky, took much longer this round
  block 18  mined at t=2.55s

PROOF OF STAKE (this chapter's testnet)

  block 14  proposed at t=0.00s   (just select + sign, no search)
  block 15  proposed at t=0.05s   (consistent: nearly the same each round)
  block 16  proposed at t=0.11s
  block 17  proposed at t=0.16s
  block 18  proposed at t=0.21s
```

The qualitative difference is the point: proof of work's block times are inherently variable, because they depend on a random search whose duration follows an exponential distribution — sometimes you get lucky and solve it almost instantly, sometimes you do not. Proof of stake's block times are far more consistent, because there is no search at all — selecting a validator and producing a signature takes roughly the same, small amount of time every round, bounded mainly by network latency (getting the proposal to other validators) rather than by luck. This consistency is one of proof of stake's genuinely attractive practical properties, independent of the energy-efficiency argument from Section 1 — it is a large part of why chains that value fast, predictable block times (many application-specific and gaming chains, discussed further in Chapter 80) favor proof-of-stake or BFT-style consensus over proof of work.

---

## Summary

- Proof of stake replaces "spend electricity to earn the right to propose a block" with "lock up gochips as collateral, and lose them (slashing) if caught misbehaving" — a different currency of security, not an unqualified improvement over proof of work.
- Validator selection is weighted-random by stake: a validator with twice the stake of another is twice as likely to be selected in any given round, computed deterministically from an agreed-upon seed so every honest node reaches the same answer.
- Slashing gives staking real teeth: a caught validator loses part or all of their stake, and a fully-slashed validator is removed from future eligibility entirely.
- `consensus.ProofOfStake` implements the exact same `consensus.Engine` interface as `consensus.ProofOfWork` from Volume 4, so `core.Blockchain` requires zero changes to run under either algorithm.
- `SelectValidator` hashes a seed (typically the previous block's hash) into a deterministic "raffle ticket" and walks accumulated stake ranges to pick a winner — unpredictable in advance, verifiable after the fact.
- `Mine` repurposes the shared interface's `nonce` return value to record which validator proposed the block, and `Validate` independently recomputes who *should* have been selected to check it.
- `Slash` reduces a validator's stake and removes them from the set entirely once their stake reaches zero, tested directly against both full and partial-penalty scenarios.
- Proof of stake's block times are far more consistent than proof of work's inherently random solve times, because proposing a block is a fixed-cost signature rather than an open-ended search.

---

## Exercises

### Easy

1. In your own words, explain the security-deposit analogy from Section 1. What does a validator "lose" by cheating, and how is that similar to (and different from) what a dishonest proof-of-work miner loses?
2. Why must `SelectValidator` be deterministic given the same seed, rather than using ordinary, non-reproducible randomness?
3. What is the difference between a full slash and a partial slash, and why might a real system use different penalty sizes for different offenses?

### Medium

4. Run `TestSelectValidator_WeightedFairness` with the tolerance in `assertRoughly` reduced from 0.05 to 0.005. Does the test still pass reliably across several runs? What does this tell you about how many rounds are needed to get a statistically tight estimate of a validator's true selection frequency?
5. `Mine`'s implementation uses `b.PrevBlockHash` as the selection seed. Explain, using Section 2's "unpredictable in advance, verifiable after the fact" requirement, why using something like the current wall-clock time as the seed instead would be a security flaw.
6. Extend `ProofOfStake` with a method `TotalStake() int64` that exposes `totalStake()` publicly, and use it to write a test asserting that after slashing a validator to zero, `TotalStake()` correctly reflects the reduced total.
7. Wire `ProofOfStake` into a version of Chapter 76's three-node honest network (replacing `consensus.NewProofOfWork`), run it for the same duration, and compare the variance in block-arrival times against the PoW version's output from Section 12.

### Hard

8. Real proof-of-stake systems typically require a minimum stake to become a validator at all (to prevent an attacker from registering thousands of tiny-stake "validators" purely to increase their odds of being selected across many identities rather than through legitimate stake). Add a `MinStake int64` field to `ProofOfStake` and reject validators below it in `NewProofOfStake`. Explain why this defense is conceptually similar to, but not identical to, Chapter 51's discussion of Sybil attacks.
9. `SelectValidator`'s use of `b.PrevBlockHash` as a seed means a validator who is *about to be selected* for the next block already knows the seed for the block after that, once their own block is finalized. Research "proposer front-running" or "stake grinding" attacks in real proof-of-stake systems, and explain, in your own words, one concrete way a validator could try to manipulate future selections in their own favor, and one mitigation real systems (such as Ethereum's beacon chain RANDAO) use against it.
10. Implement a simple slashing-evidence mechanism: a function `DetectEquivocation(a, b *core.Block) bool` that returns true if two blocks have the same `Height` and the same claimed proposer (via the `Nonce` field) but different `Hash` values — direct proof the same validator signed two conflicting blocks for the same round. Wire it so that when a node observes two such blocks, it automatically calls `Slash` against the offending validator.
