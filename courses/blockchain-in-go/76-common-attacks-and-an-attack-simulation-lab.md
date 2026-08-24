# Chapter 76: Common Attacks and an Attack Simulation Lab

Every mechanism you have built so far — hashing, signatures, proof of work, the mempool, gossip — exists because something goes wrong without it. This chapter names the attacks those mechanisms defend against precisely, then puts you on the other side of the table: you will build a small, deliberately dishonest GoChain node whose entire job is to rewrite history, and watch it try (and fail, and eventually succeed against a network too small to resist it) against a local test network you also control.

## Table of Contents

1. [Why Attack Your Own Chain](#1-why-attack-your-own-chain)
2. [The 51% Attack, In Depth](#2-the-51-attack-in-depth)
3. [Double-Spending Beyond the Basics](#3-double-spending-beyond-the-basics)
4. [Replay Attacks](#4-replay-attacks)
5. [Race Attacks on Unconfirmed Transactions](#5-race-attacks-on-unconfirmed-transactions)
6. [Lab Setup: A Small Honest Test Network](#6-lab-setup-a-small-honest-test-network)
7. [Building the Attacker: A Secret Mining Node](#7-building-the-attacker-a-secret-mining-node)
8. [Running the Attack and Reading the Output](#8-running-the-attack-and-reading-the-output)
9. [Measuring Cost and Detectability](#9-measuring-cost-and-detectability)
10. [Real-World Defenses](#10-real-world-defenses)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Attack Your Own Chain

Understanding a lock by staring at it teaches you less than trying, honestly and methodically, to pick it. Every attack in this chapter targets a specific mechanism you built in earlier volumes — proof of work's cost (Volume 4), the mempool's double-spend check (Volume 5), transaction signing (Volume 5), and the longest-chain fork rule (Volume 7) — so understanding the attack is really a second, adversarial pass over material you already know, this time from the attacker's chair.

A quick vocabulary note before the details: an **attacker**, in this chapter, is any participant who deviates from the protocol's rules to gain an advantage — spending the same coins twice, rewriting history, or tricking another party into accepting a payment that later disappears. None of this requires the attacker to be some shadowy criminal genius. It just requires ordinary hardware, patience, and a willingness to lie.

---

## 2. The 51% Attack, In Depth

Chapter 24 previewed this idea conceptually: proof of work makes rewriting old blocks expensive because you would have to redo the puzzle for that block *and every block mined after it*, faster than the rest of the network combined. A **51% attack** is what happens when someone actually has that much mining power — more than half of the network's total hash rate — and chooses to use it dishonestly.

Here is precisely what a majority of hash power buys an attacker, and — just as important — what it does *not* buy them:

- **It can rewrite recent history.** The attacker can secretly mine an alternative chain, starting from some block in the recent past, and because they out-pace the honest network, their secret chain eventually has more accumulated proof of work than the public one. Broadcasting it triggers the longest-chain (heaviest-chain) rule from Chapter 50, and every honest node switches over, discarding the blocks and transactions that were only ever on the old chain.
- **It can double-spend.** By rewriting history, the attacker can erase a payment they made (say, to an exchange) after the exchange already credited them and let them withdraw something else in return, while keeping the coins they "spent" on their own, private version of history.
- **It cannot forge transactions from other people.** Signatures (Volume 2) are still cryptographically sound. The attacker cannot spend Alice's coins without Alice's private key, no matter how much hash power they control. A 51% attack rewrites *whose block wins*, not *who is allowed to spend what*.
- **It cannot create coins out of thin air.** Every block, including the attacker's, still has to satisfy every validation rule from `core.ValidateBlock` and `core.ValidateTransaction` — the attacker is not exempt from consensus rules, only powerful enough to out-race everyone else at proposing blocks that follow them.
- **It gets more expensive, not cheaper, the deeper the target block is buried.** Rewriting a block from ten minutes ago is far easier than rewriting one from ten thousand blocks ago, because the honest network has kept mining the whole time, and the attacker must out-mine *all* of that accumulated work, not just replace one block.

```
HONEST CHAIN (public, growing normally)

   B0 --- B1 --- B2 --- B3 --- B4 --- B5      <- tip everyone sees
                          ^
                          |
              attacker secretly forks here,
              mines B3' B4' B5' B6' in private,
              telling nobody

ATTACKER'S SECRET CHAIN (mined off to the side, never broadcast until ready)

   B0 --- B1 --- B2 --- B3' --- B4' --- B5' --- B6'   <- one block longer

THE MOMENT THE ATTACKER BROADCASTS:

   every honest node compares total accumulated work,
   finds the attacker's chain heavier (6 blocks vs 5 past the fork point),
   and REORGANIZES onto it -- B3, B4, B5 and anything only they contained
   are discarded as if they never happened.

   B0 --- B1 --- B2 --- B3' --- B4' --- B5' --- B6'   <- new accepted tip
                  \
                   `-- B3 -- B4 -- B5   (orphaned; transactions here that
                                          don't also appear in the new chain
                                          are back in the mempool, unconfirmed,
                                          or gone if they conflict)
```

The key number to internalize is: an attacker with hash power `p` (as a fraction of the total network) and a target `k` blocks behind the tip has a success probability that shrinks rapidly as `k` grows, and shrinks to essentially zero once `p` is comfortably under 0.5 and `k` is more than a handful of blocks. This is the precise mathematical reason "wait for six confirmations" is common advice for high-value Bitcoin transactions: at six blocks deep, an attacker with anything less than a large fraction of total hash power has a vanishingly small chance of ever catching up, even if they get lucky early on.

A rough intuition for the numbers, without deriving the full random-walk math behind it (that derivation is a good research exercise at the end of this chapter):

```
Approximate probability the attacker EVER catches up, by hash share and depth

  attacker's share (p)   1 block deep    3 blocks deep   6 blocks deep
  10%                    ~10%            ~1%             ~0.02%
  30%                    ~30%            ~13%            ~2.6%
  45%                    ~45%            ~40%            ~32%
  51%                    100% eventually  100% eventually  100% eventually
```

Two things stand out. First, below 50% hash share, the attacker's odds fall off fast as the required depth grows — this is exactly why confirmations work as a defense even against attackers with a meaningful, non-trivial share of hash power. Second, once an attacker's share crosses 50%, depth stops helping at all: given unlimited time, an attacker with a genuine hash-power majority will *eventually* out-mine any fixed number of confirmations, which is why this whole class of attack is named for the 51% threshold rather than for any specific confirmation count.

---

## 3. Double-Spending Beyond the Basics

Volume 5 introduced double-spending as simply "the mempool rejects a second transaction that tries to spend an already-claimed UTXO." That is correct as far as it goes, but it only stops the *simplest* version of the attack — one where both conflicting transactions reach the same mempool. Real double-spend attacks are more devious about **when** and **where** the two conflicting transactions are shown to whom.

- **The basic double-spend (Volume 5).** Attacker broadcasts transaction A (pay merchant), merchant sees it in their own mempool and ships goods immediately, without waiting for it to be mined. Attacker then gets a conflicting transaction B (pay themselves, spending the same UTXO) mined instead — accomplished by, for example, having pre-arranged for a friendly miner to include B instead of A, or by getting B relayed with a much higher fee so miners prefer it. The merchant is left with a mempool entry that will never confirm.
- **The Finney attack.** Named after Hal Finney, this variant requires the attacker to *be* a miner. They pre-mine a block containing a payment to themselves (transaction B) but do not broadcast it. They then make a real-world purchase, paying the merchant with transaction A, which spends the *same* coins. If the merchant accepts the payment before it is confirmed, the attacker immediately broadcasts their pre-mined block containing B — which the network now sees as the earlier, valid spend, and A is rejected as a double-spend of an already-consumed UTXO.
- **Vector76 (a "race plus one confirmation" attack).** A hybrid: the attacker mines a block privately containing transaction B (pay self), then simultaneously broadcasts transaction A (pay merchant) to the network and their own pre-mined block to a single, specifically targeted node (often the merchant's own node, or a mining pool likely to build on it) at the same instant. If the merchant's node happens to accept the attacker's private block, the merchant may see "one confirmation" for a payment that the wider network will actually reorganize away moments later, once it receives a competing block from elsewhere.

```
BASIC DOUBLE-SPEND                    FINNEY ATTACK
                                       (attacker must also be a miner)

Attacker -> Tx A (pay merchant)        Attacker privately mines a block
              \                        containing Tx B (pay self),
               \-> merchant's mempool  keeps it SECRET
                   (unconfirmed)
                                       Attacker -> Tx A (pay merchant),
Attacker also arranges for                        merchant ships goods
Tx B (pay self, same UTXO)                          on 0-confirmation trust
to get mined instead
                                       Attacker broadcasts the pre-mined
Merchant's Tx A never confirms;        block (containing Tx B) immediately
goods already shipped, no payment      after receiving goods

                                       Network sees Tx B was "already"
                                       mined; Tx A is rejected as a
                                       double-spend
```

The common thread across every variant: **the merchant acted on a payment before it was buried under enough confirmations to make rewriting it prohibitively expensive.** Every one of these attacks becomes dramatically harder, and eventually practically impossible, the more confirmations a merchant requires before treating a payment as final — which is exactly the trade-off between convenience (accept 0-conf payments instantly) and safety (wait for confirmations) every real payment processor has to make explicitly.

---

## 4. Replay Attacks

A **replay attack** takes a transaction that was validly signed and broadcast in one context and rebroadcasts it in a *different* context where it is unexpectedly still valid — even though nobody involved intended it to apply there. The classic real-world case is a blockchain **fork** that splits into two independent, ongoing chains (as opposed to Chapter 50's temporary forks that quickly resolve) — this happened when Ethereum split into Ethereum and Ethereum Classic in 2016. Because both resulting chains initially shared the exact same transaction and signature format, a transaction signed and broadcast on one chain was often *also* a perfectly validly signed transaction on the other — so if Alice sent coins to Bob on Chain A, an attacker (or Bob himself) could take that exact same signed transaction and rebroadcast it, unmodified, on Chain B, moving Alice's Chain-B coins too, without ever needing her private key again.

```
BEFORE THE SPLIT                 AFTER THE SPLIT

     shared chain                Chain A (original)      Chain B (new fork)
  ... -- B100 -- B101            ... -- B101 -- B102a    ... -- B101 -- B102b
                                          |                       |
                              Alice signs Tx: "pay Bob 5"  same signed bytes
                              broadcasts on Chain A        are ALSO a valid
                                                            transaction here --
                                                            replayed without
                                                            Alice's involvement
```

GoChain is not immune to this in principle: because `Transaction.Sign` (Volume 5) signs a trimmed copy of the transaction's own fields, a signed transaction says nothing about *which chain* it was intended for unless we explicitly make it say so. Real systems close this gap by including **chain-specific replay protection** in every signature — Ethereum's `EIP-155`, for example, folds a numeric chain ID directly into the data that gets signed, so a transaction signed for chain ID 1 produces a completely different signature than the "same" transaction signed for chain ID 61, and is rejected by the other chain's nodes as an invalid signature rather than silently replayed. If GoChain ever forked into two independently maintained networks, the fix would be to add a `ChainID` field to `core.Transaction` and include it in the bytes `Sign` hashes and signs — a small change with an outsized security benefit.

---

## 5. Race Attacks on Unconfirmed Transactions

A **race attack** is the narrowest, fastest version of double-spending: it does not require any mining power at all, only the ability to get two conflicting transactions in front of two different victims before either one confirms. The attacker broadcasts transaction A (pay merchant) directly to the merchant's node, while simultaneously broadcasting a conflicting transaction B (pay self, same UTXO) widely across the rest of the network. If B propagates faster and gets picked up by the miner who solves the next block, A never confirms — but the merchant, watching only their own node's mempool, may have already treated A's mere presence as good enough to release goods.

This is precisely why GoChain's mempool (Volume 5) enforcing "no double-spend among pending transactions" only protects *that node's own view* — it does nothing to protect a merchant who acts on an unconfirmed transaction before the network as a whole has settled which of two conflicting transactions wins. The practical lesson mirrors Section 3's: **zero-confirmation trust is a business decision to accept a specific, quantifiable risk**, appropriate for a cup of coffee, and never appropriate for a house.

---

## 6. Lab Setup: A Small Honest Test Network

Now the hands-on part. We build a tiny, three-node honest GoChain network, reusing the `network.Node` type from Volume 7 and `consensus.ProofOfWork` from Volume 4, deliberately running with **low difficulty** so a single attacker's laptop can plausibly out-mine it — the same reason real-world 51% attacks have historically targeted smaller, lower-hash-rate chains rather than Bitcoin itself.

```go
// cmd/attacksim/honestnet.go
package main

import (
	"log"
	"time"

	"github.com/you/gochain/consensus"
	"github.com/you/gochain/core"
	"github.com/you/gochain/network"
)

// startHonestNode boots one ordinary, rule-following GoChain node: it mines
// on top of whatever it believes is the current tip, gossips new blocks to
// its peers, and accepts a longer valid chain from anyone who broadcasts one
// -- exactly the behavior built in Volumes 4 and 7, nothing attacker-specific.
func startHonestNode(addr string, peers []string, dataDir string) (*network.Node, *core.Blockchain) {
	chain, err := core.OpenBlockchain(dataDir, "genesis-reward-address")
	if err != nil {
		log.Fatalf("open chain: %v", err)
	}

	// Deliberately low difficulty: this lab is about the DYNAMICS of a 51%
	// attack, not about burning real electricity to demonstrate it.
	engine := consensus.NewProofOfWork(12) // 12 leading zero bits

	node := network.NewNode(addr, chain, engine)
	for _, p := range peers {
		node.Dial(p) // connect out to the other honest nodes
	}

	go node.Listen() // accept inbound connections and gossip messages

	// A background goroutine that mines continuously, one block at a time,
	// picking up whatever transactions are waiting in the mempool -- this is
	// what makes the "honest network" keep growing while the attacker works.
	go func() {
		for {
			block := chain.MineBlock(engine)
			node.BroadcastBlock(block) // gossip it to peers immediately
			time.Sleep(500 * time.Millisecond)
		}
	}()

	return node, chain
}

func main() {
	// Three honest nodes, all peered with each other, all mining and
	// gossiping normally -- the "victim" network the attacker will target.
	_, _ = startHonestNode(":9001", []string{":9002", ":9003"}, "./data/honest1")
	_, _ = startHonestNode(":9002", []string{":9001", ":9003"}, "./data/honest2")
	_, _ = startHonestNode(":9003", []string{":9001", ":9002"}, "./data/honest3")

	select {} // keep the process alive; nodes run in their own goroutines
}
```

`startHonestNode` boots one ordinary GoChain node using components you already built: `core.OpenBlockchain` from Volume 8, `consensus.NewProofOfWork` from Volume 4, and `network.NewNode` from Volume 7. The only thing new in this function is the deliberately low difficulty (12 leading zero bits instead of a realistic production value), chosen so this lab finishes in seconds on a laptop rather than requiring real mining hardware. `main` starts three of these, fully peered with each other, and lets them run — this is the honest network the rest of the lab attacks.

---

## 7. Building the Attacker: A Secret Mining Node

The attacker is structurally almost identical to an honest node — it uses the exact same `consensus.ProofOfWork` engine and the exact same block validation rules, because **a 51% attack does not require breaking any cryptography or any validation rule**. It requires exactly one behavioral difference: mine in secret, on a chain fork of your own choosing, and only reveal it once it is ready to win.

```go
// cmd/attacksim/attacker.go
package main

import (
	"fmt"
	"log"
	"time"

	"github.com/you/gochain/consensus"
	"github.com/you/gochain/core"
	"github.com/you/gochain/network"
)

// AttackerNode holds a private, disconnected copy of the chain that the
// attacker mines on secretly. It only talks to the honest network at the
// very end, to broadcast the finished, longer chain all at once.
type AttackerNode struct {
	privateChain *core.Blockchain
	engine       consensus.Engine
	targetHeight int // the honest block height we are trying to rewrite from
}

// NewAttackerNode forks the private chain starting from forkHeight -- i.e. it
// keeps every honest block up to and including forkHeight, then discards
// anything after it, so the attacker can start mining an alternative history
// from that point without needing to redo the entire chain from genesis.
func NewAttackerNode(honestChain *core.Blockchain, forkHeight int) *AttackerNode {
	private := honestChain.CopyUpTo(forkHeight) // duplicate blocks [0..forkHeight]
	return &AttackerNode{
		privateChain: private,
		engine:       consensus.NewProofOfWork(12), // same difficulty as the honest network
		targetHeight: forkHeight,
	}
}

// mineSecretly mines `count` blocks on the private fork, entirely offline --
// no gossip, no broadcast, nobody else on the network learns this is
// happening. This is the "secretly mine a longer competing chain" step.
func (a *AttackerNode) mineSecretly(count int) {
	for i := 0; i < count; i++ {
		block := a.privateChain.MineBlock(a.engine)
		fmt.Printf("[attacker] secretly mined block %d (hash %x), total private height %d\n",
			block.Height, block.Hash[:4], a.privateChain.Height())
	}
}

// launch connects to the honest network for the first time and broadcasts
// every block of the private chain, oldest first. If the private chain has
// more accumulated work than what the honest nodes currently hold, their own
// longest-chain rule (Volume 7, Chapter 50) does the rest: they reorganize.
func (a *AttackerNode) launch(honestPeers []string) {
	node := network.NewNode(":9099", a.privateChain, a.engine)
	for _, p := range honestPeers {
		node.Dial(p)
	}

	fmt.Println("[attacker] broadcasting private chain to honest network NOW")
	for _, block := range a.privateChain.BlocksFrom(a.targetHeight + 1) {
		node.BroadcastBlock(block)
		time.Sleep(50 * time.Millisecond) // small delay so peers can process in order
	}
}

func main() {
	honestChain, err := core.OpenBlockchain("./data/honest1", "")
	if err != nil {
		log.Fatalf("could not read honest chain: %v", err)
	}

	forkPoint := honestChain.Height() - 3 // fork three blocks back from the current tip
	attacker := NewAttackerNode(honestChain, forkPoint)

	// The attacker must out-mine however many blocks the honest network adds
	// WHILE they are working, not just match the current height -- so they
	// mine more blocks than the gap that currently exists, to build a margin.
	attacker.mineSecretly(6)

	attacker.launch([]string{":9001", ":9002", ":9003"})
}
```

`AttackerNode` is deliberately built from the same pieces as `startHonestNode` — the difference is entirely in how it is *used*. `NewAttackerNode` calls `CopyUpTo`, a new small helper on `core.Blockchain` that duplicates blocks up to a chosen fork point and discards the rest, giving the attacker a private starting point identical to the honest chain up to that block. `mineSecretly` then mines additional blocks on this private copy without ever calling `node.BroadcastBlock` — nothing here is visible to the honest network, which keeps mining its own, public continuation the whole time. Finally, `launch` connects to the honest peers for the very first time and broadcasts every block from the fork point forward, in order — at which point each honest node's existing fork-resolution logic (Chapter 50) takes over automatically, with no attacker-specific code needed on the honest side at all.

---

## 8. Running the Attack and Reading the Output

Run the honest network in one terminal, let it mine for thirty seconds or so, then run the attacker in a second terminal.

```bash
# terminal 1
$ go run ./cmd/attacksim/honestnet.go
[node :9001] mined block 14 (hash 0a3f...)
[node :9002] mined block 15 (hash 71c2...)
[node :9003] mined block 16 (hash 9e01...)
...

# terminal 2, started after the honest network has reached height ~16
$ go run ./cmd/attacksim/attacker.go
[attacker] secretly mined block 14 (hash 55b0...), total private height 14
[attacker] secretly mined block 15 (hash 2a9d...), total private height 15
[attacker] secretly mined block 16 (hash c810...), total private height 16
[attacker] secretly mined block 17 (hash 004f...), total private height 17
[attacker] secretly mined block 18 (hash 88e2...), total private height 18
[attacker] secretly mined block 19 (hash f130...), total private height 19
[attacker] broadcasting private chain to honest network NOW
```

Back in terminal 1, the honest nodes react the instant the attacker's blocks arrive:

```
[node :9001] received block 14' (hash 55b0...) -- competing fork detected
[node :9001] comparing accumulated work: local chain height 17, incoming chain height 19
[node :9001] incoming chain is heavier -- REORGANIZING
[node :9001] rolled back blocks 15, 16, 17 (local); replaying mempool for any
             transactions that were only in the discarded blocks
[node :9001] adopted attacker's chain up to block 19'
[node :9002] received block 14' ... REORGANIZING ... adopted up to block 19'
[node :9003] received block 14' ... REORGANIZING ... adopted up to block 19'
```

If you re-run the lab with `mineSecretly(2)` instead of `mineSecretly(6)`, and the honest network has kept mining the whole time, you will usually see the opposite outcome: the honest chain remains heavier when the attacker broadcasts, and every honest node correctly rejects the shorter private fork, logging something like `incoming chain is NOT heavier -- ignoring`. This is the exact mechanism from Section 2 made concrete: whether the attack succeeds depends entirely on whether the attacker's private hash power, sustained over the attack window, exceeded the honest network's combined hash power over that same window.

---

## 9. Measuring Cost and Detectability

A 51% attack is not free, and this lab lets you measure exactly how not-free it is:

- **Compute cost.** Count how many nonces the attacker's `consensus.ProofOfWork` engine tried across `mineSecretly`'s six blocks versus how many the three honest nodes tried, combined, over the same wall-clock window. At realistic mainnet difficulty, this ratio translates directly into real electricity and hardware rental cost — this is precisely why 51% attacks are calculated, expensive business decisions for real attackers, not casual pranks.

```go
// costEstimate turns a raw nonce count into a rough, illustrative dollar
// figure, so the lab's abstract "hash power" becomes a number you can reason
// about the same way a real attacker (or a real exchange's risk team) would.
func costEstimate(noncesTried uint64, hashesPerSecondPerDollar float64) float64 {
	// hashesPerSecondPerDollar approximates how much compute one dollar of
	// rented cloud CPU or GPU time buys per second -- swap in a real,
	// current figure for your own experiments.
	seconds := float64(noncesTried) / hashesPerSecondPerDollar
	return seconds // interpret as "dollar-seconds" of rented compute
}
```

`costEstimate` is intentionally simple: it exists to turn "the attacker tried six billion nonces" into a number with real-world units attached, which is the same translation any serious cost/benefit analysis of a real attack has to make before deciding whether the attack is worth attempting at all.
- **Detectability.** A deep chain reorganization is not silent. Real nodes (and GoChain's `chain inspector` CLI from Chapter 21) can log and alert on any reorg deeper than some small threshold — a reorg of 1-2 blocks is normal background noise from honest forks (Chapter 50), but a reorg of 6+ blocks, arriving all at once from a single previously-unseen peer, is a strong signal worth paging a human about. Add a simple check to your lab's honest nodes: `if reorgDepth > 2 { log.Printf("ALERT: deep reorg of %d blocks from peer %s", reorgDepth, peerAddr) }`.
- **Economic exposure.** Whatever value the attacker double-spent is bounded by how much a victim was willing to accept at how few confirmations — tying this lab directly back to Sections 3 and 5's core lesson: the real defense against double-spending was never purely technical, it is the confirmation policy a merchant chooses to enforce.

```
COST VS. DEPTH -- why attacking a deeply-buried block is much harder

  target depth (blocks to rewrite)     work required (relative)
  1 block back                          1x
  3 blocks back                         3x  (must also out-race ongoing mining)
  6 blocks back                         6x+ (common "safe" confirmation count)
  1000 blocks back                      impractical for all but the largest
                                         possible attackers, against any chain
                                         with meaningful honest hash power
```

---

## 10. Real-World Defenses

None of this lab's attacks have a single silver-bullet fix — they are managed with layered, practical defenses:

- **Wait for confirmations proportional to value.** A coffee shop can reasonably accept 0 confirmations; an exchange crediting a large deposit should not, and most real exchanges enforce chain-specific minimum confirmation counts for exactly this reason.
- **Checkpointing.** Some chains hard-code (or socially agree on) specific historical block hashes as unquestionably final, so no amount of hash power can ever be used to rewrite history before that point — a pragmatic, if slightly centralized, patch on top of pure proof-of-work security.
- **Monitoring hash rate concentration.** A sudden, large increase in one mining pool's share of total network hash rate is a leading indicator worth watching, well before any attack is attempted.
- **Chain-specific replay protection**, as discussed in Section 4, closes the replay-attack class entirely by making a transaction signed for one chain simply invalid on another.
- **Choosing a chain with enough honest hash power in the first place.** This is the single biggest real-world factor: 51% attacks have repeatedly succeeded against smaller proof-of-work coins with modest total hash rate (several real incidents are covered in Volume 12's Chapter 85), precisely because renting enough hash power to exceed theirs is commercially affordable, whereas doing the same against Bitcoin is not.

---

## Summary

- A 51% attack lets an attacker rewrite recent history and double-spend, but cannot forge other people's signatures or create coins from nothing — it only changes *whose blocks win*, not the underlying validation rules.
- Double-spending has several real variants beyond the basic mempool case: the Finney attack (attacker is also a miner), Vector76 (a targeted hybrid), and the plain race attack (no mining power needed at all, just fast, targeted broadcasting).
- Replay attacks rebroadcast a validly signed transaction in a different context (typically after a chain split) where it is unexpectedly still valid; chain-specific replay protection (like Ethereum's `EIP-155`) fixes this by folding a chain ID into what gets signed.
- This chapter's lab builds a real, working attacker using the same `consensus.ProofOfWork` and `network.Node` types from earlier volumes — the attack requires no new cryptography, only mining secretly on a private fork and broadcasting it once it is longer.
- Whether the attack succeeds depends entirely on whether the attacker's sustained hash power exceeded the honest network's, over the attack window — exactly the mathematical relationship from Chapter 24, now observed directly in your own terminal output.
- Deep chain reorganizations are detectable and alertable; this is a practical, code-level defense worth building into any real node, not just a theoretical mitigation.
- The strongest real-world defense against every attack in this chapter is a confirmation policy matched to the value at risk — there is no purely technical fix that replaces this business decision.

---

## Exercises

### Easy

1. In your own words, list the two things a 51% attack lets an attacker do, and the two things it explicitly does not let them do.
2. Explain the difference between a race attack and a Finney attack. Which one requires the attacker to also be a miner?
3. Why does `Section 8`'s lab use a deliberately low proof-of-work difficulty (12 bits) instead of a realistic production value?

### Medium

4. Run the lab with `mineSecretly(2)` instead of `mineSecretly(6)` while the honest network keeps mining for at least 15 seconds before you launch the attack. Report whether the attack succeeded or failed, and explain why, referencing the accumulated-work comparison from Section 2.
5. Add the reorg-depth alert described in Section 9 to the honest node code from Section 6. Trigger it with an attack and paste the resulting log line.
6. Explain, using Section 4's diagram, exactly what field you would add to `core.Transaction` to prevent replay attacks across two independently forked GoChain networks, and what that field would need to be included in.

### Hard

7. Extend `AttackerNode` so that, instead of mining entirely offline before broadcasting, it interleaves secret mining with periodically checking the honest network's current height (without revealing its own chain), and stops mining and launches the moment its private chain is exactly 2 blocks ahead. Measure whether this "just enough" strategy uses meaningfully less compute than `mineSecretly(6)`'s fixed approach.
8. Modify the lab so the "honest" nodes require a configurable number of confirmations before considering a transaction spendable for a simulated real-world payment (e.g., before printing "goods shipped"). Re-run the double-spend scenario from Section 3 at 0, 1, and 6 required confirmations, and report at which confirmation count the simulated merchant stops losing money.
9. Research a real, historical 51% attack against a smaller proof-of-work cryptocurrency (Volume 12's Chapter 85 covers one in depth, but find your own if you prefer). Report the attacker's approximate rented hash power, the chain's total hash rate at the time, the estimated cost of the attack, and the value the attacker is believed to have double-spent. Was the attack profitable, based on your research?
