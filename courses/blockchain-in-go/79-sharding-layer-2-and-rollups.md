# Chapter 79: Sharding, Layer 2, and Rollups

Every design decision so far in this course has quietly assumed one thing: every honest node processes every single transaction, in the same order, forever. That assumption is exactly what makes GoChain's security model work — Chapter 19's tamper-evidence, Chapter 25's proof of work, Chapter 77's proof of stake, all of them lean on the idea that the whole network agrees on one linear history. It's also exactly what puts a hard ceiling on how much a chain can ever process: if every node must redo every transaction, the network's total throughput can never exceed what a *single* node can keep up with, no matter how many thousands of nodes join. This chapter is about the two broad families of techniques real chains use to push past that ceiling — sharding (splitting the *chain itself* into parallel pieces) and Layer 2 (moving computation *off* the chain while still relying on it for security) — explained conceptually, with a clear picture of exactly what each approach gives up to get its speed.

## Table of Contents

1. [The Bottleneck: One Chain, Every Node](#1-the-bottleneck-one-chain-every-node)
2. [Sharding, Conceptually](#2-sharding-conceptually)
3. [Three Kinds of Sharding](#3-three-kinds-of-sharding)
4. [Layer 2: Moving Computation Off-Chain](#4-layer-2-moving-computation-off-chain)
5. [State Channels](#5-state-channels)
6. [Optimistic Rollups](#6-optimistic-rollups)
7. [Zero-Knowledge Rollups](#7-zero-knowledge-rollups)
8. [The Arithmetic of a Rollup Batch](#8-the-arithmetic-of-a-rollup-batch)
9. [What Each Approach Still Verifies on the Main Chain](#9-what-each-approach-still-verifies-on-the-main-chain)
10. [What Scaling Techniques Actually Trade Away](#10-what-scaling-techniques-actually-trade-away)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Bottleneck: One Chain, Every Node

Picture a small town with exactly one notary public, and a town charter that requires every legal document — deeds, wills, contracts — to be personally read, stamped, and filed by that one notary, in the exact order they arrive, before it's considered official. As the town grows from a hundred people to a hundred thousand, the notary doesn't get faster; the line outside their office just gets longer. Hiring more notaries doesn't obviously help either, if the charter still insists every document must go through the *same* notary in the *same* order, for everyone to agree on what's official.

This is precisely GoChain's situation once the network has real load. Every honest node re-executes every transaction in every block to independently verify the chain, exactly as Chapter 19 designed it — which is the whole source of the chain's trustworthiness, but also means the network's total transaction throughput can never exceed roughly "one machine's worth" of processing, no matter how many thousands of nodes are running GoChain.

```
   ONE CHAIN, EVERY NODE PROCESSES EVERYTHING

   Node 1   Node 2   Node 3   Node 4   ...   Node 10,000
     |        |        |        |               |
     +--------+--------+--------+---------------+
                        |
              same blocks, same order,
              same transactions --
              re-executed independently
              by EVERY single node

   Adding more nodes improves DECENTRALIZATION and
   FAULT TOLERANCE, but does NOT raise throughput --
   the bottleneck is what ONE node can process, not
   how many nodes exist.
```

Both families of scaling techniques in this chapter attack this exact bottleneck, but from opposite directions: **sharding** (Sections 2-3) says "stop making every node process everything — split the work into parallel pieces." **Layer 2** (Sections 4-8) says "stop putting everything on the chain in the first place — do most of the work elsewhere, and only touch the chain when it truly matters."

---

## 2. Sharding, Conceptually

**Sharding** splits the network into multiple parallel sub-chains, called **shards**, each responsible for processing only a fraction of the network's total transactions — the way a busy restaurant might split its dining room into sections, each staffed by a different set of servers, rather than making every single server walk to every single table.

```
   WITHOUT SHARDING                    WITH SHARDING (3 shards)

   [ all transactions ]                [ transactions routed by
          |                              account/address ]
          v                                |      |      |
   +--------------+                        v      v      v
   |  every node   |                   Shard A  Shard B  Shard C
   |  processes     |                  (nodes    (nodes   (nodes
   |  EVERYTHING     |                  1-100)    101-200) 201-300)
   +--------------+
                                       Each shard has its OWN chain,
                                       its OWN validators, and processes
                                       its OWN transactions IN PARALLEL
                                       with the other shards.
```

The obvious question sharding raises immediately: what happens when a transaction needs to touch data on *two different shards* — say, an address on Shard A sending funds to an address on Shard B? This **cross-shard communication** problem is sharding's central engineering difficulty, and every real sharded system (Ethereum's own sharding research included) spends most of its complexity budget here, typically through some combination of:

- A **beacon chain** or coordinating layer that shards periodically report summaries to, so cross-shard references can eventually be verified against a shared source of truth.
- **Receipts and proofs**: Shard A produces a cryptographic receipt proving "this cross-shard transaction really happened here," which Shard B can verify without needing to replay Shard A's entire history.
- Accepting that cross-shard transactions are inherently slower and more complex than same-shard ones, and designing applications to minimize how often they're needed.

The trade-off sharding makes is direct: parallelism in exchange for a strictly more complicated security and communication model — every shard, individually, is now watched over by fewer validators than the whole network would provide, which is exactly the concern Section 10 returns to.

---

## 3. Three Kinds of Sharding

"Sharding" is often discussed as though it were one single technique, but real designs (including Ethereum's own multi-year sharding research program) actually distinguish three separate things that can each be split, independently of one another:

```
   NETWORK SHARDING              TRANSACTION SHARDING           STATE SHARDING

   Split which PEERS a           Split WHICH transactions       Split the STORAGE of
   node needs to connect          get assigned to which          account/UTXO state
   to and gossip with --          shard for processing --        itself across shards,
   a node only needs to           determined by, e.g., a         so no single node
   fully participate in           hash of the sender             needs to hold a copy
   the peer-to-peer mesh          address modulo the shard        of the ENTIRE network's
   for the shard(s) it            count                          state -- only its own
   actually cares about                                          shard's slice of it

   Reduces bandwidth and         Reduces COMPUTE load per        Reduces STORAGE load per
   connection overhead per       node -- each shard only         node -- the hardest of
   node                          re-executes ITS OWN               the three to get right,
                                 transactions                     because cross-shard
                                                                  reads (Section 2's
                                                                  problem) become
                                                                  fundamentally harder
                                                                  when even the DATA
                                                                  being read might live
                                                                  on a different shard
```

These three are not mutually exclusive — a mature sharding design typically layers all three together, since splitting transaction processing without also splitting storage would still leave every node holding a full copy of the entire network's state, just processing a subset of the transactions that touch it, which caps how much a design like this can ultimately scale. State sharding is consistently described, across real sharding research (Ethereum's included), as the hardest of the three to implement safely, precisely because Section 2's cross-shard communication problem gets sharper: a transaction on Shard A that needs to *read* account data living on Shard B cannot just wait for an eventual receipt the way a simple value-transfer can — it needs that data to actually be available, correctly, at execution time.

---

## 4. Layer 2: Moving Computation Off-Chain

Where sharding changes how the *base chain itself* processes transactions, **Layer 2** leaves the base chain (called **Layer 1**, or **L1** — GoChain itself, in our case) entirely unchanged, and instead builds a second system on top of it that does most of the actual work, only touching L1 when something needs L1's security guarantees specifically.

An analogy: think of Layer 1 as a country's central bank, and a Layer 2 system as a local credit union. You don't walk into the central bank for every single coffee purchase — the credit union handles day-to-day transactions quickly and cheaply, and only occasionally "settles" a summary of its activity with the central bank, which is where the ultimate, unforgeable record of who-owns-what is anchored. The credit union can only get away with this because, if it ever tried to cheat, the central bank's records would eventually expose the discrepancy — Layer 2's whole design challenge is making that same guarantee hold cryptographically, without a trusted credit union at all.

```
   LAYER 1 (GoChain itself)

     - the ultimate source of truth
     - slow-ish, but maximally secure and decentralized
     - every honest node verifies everything (Section 1's bottleneck)

              ^
              |  occasional, compressed settlement
              |
   LAYER 2 (a system built ON TOP of L1)

     - most activity happens HERE, fast and cheap
     - security is INHERITED from L1, not independently re-created
     - periodically posts a compact summary (or proof) back to L1
```

This chapter covers three concrete Layer 2 designs, in increasing order of cryptographic sophistication: state channels (Section 5), optimistic rollups (Section 6), and zero-knowledge rollups (Section 7). All three share the same core shape from the diagram above; they differ in exactly *what* gets posted back to L1 and *how* L1 is convinced it can trust that summary.

---

## 5. State Channels

A **state channel** lets two (or a small, fixed group of) participants transact with each other an unlimited number of times completely off-chain, touching the base chain only twice: once to open the channel, and once to close it.

```
   OPEN CHANNEL (on-chain)         OFF-CHAIN (unlimited, instant, free)        CLOSE CHANNEL (on-chain)

   Alice and Bob each lock          Alice pays Bob 5 gochips (signed,          Either party submits the
   10 gochips into a shared          exchanged directly, no chain)             LATEST signed balance to
   on-chain contract                                                            the chain, which pays out
                                    Bob pays Alice 2 gochips back               accordingly
   Chain now knows:                  (signed, exchanged directly)
   Alice: 10, Bob: 10                                                          Chain now knows:
                                    ... hundreds more, instantly,               Alice: 7, Bob: 13
                                    each just a new signed message               (only the OPEN and CLOSE
                                                                                  transactions ever touched
                                                                                  the chain)
```

The security property that makes this safe, not just convenient, is that every off-chain message is a **signed statement of the latest balance**, and either party can unilaterally submit the *most recent* one to the chain at any time to force a correct settlement — meaning neither party ever has to trust the other to "play fair," only to eventually publish the latest signed state if a dispute arises. The catch is right there in the name: a state channel only ever coordinates a **fixed set of participants** who opened the channel together (Alice and Bob need to know, in advance, that they'll transact with each other — this doesn't help a payment to a stranger you've never opened a channel with), which is why state channels have mostly proven useful for things like payment networks between frequent counterparties, rather than a general-purpose scaling solution for arbitrary strangers interacting.

---

## 6. Optimistic Rollups

A **rollup** generalizes the state-channel idea to *arbitrary* participants, not just a fixed pair who pre-opened a channel: a separate system (the rollup) executes a large batch of transactions off-chain, then posts a compressed summary back to L1 — the transaction data itself (compressed) plus the resulting new state root, a single hash committing to the entire updated state, exactly the kind of root Chapter 57's Merkle-Patricia trie was built to produce.

**Optimistic** rollups get their name from their core assumption: post the result *optimistically*, assume it's correct, and give anyone a window of time (typically about a week, in real deployments) to prove otherwise with a **fraud proof** — a piece of evidence showing the posted result doesn't match what honestly re-executing the batch would have produced.

```
   OPTIMISTIC ROLLUP

   Off-chain:  1,000 transactions executed, compressed,
               new state root computed

   Posted to L1:  [ compressed tx data | new state root ]
                            |
                            v
                  L1 accepts it OPTIMISTICALLY --
                  no immediate verification of correctness

                  CHALLENGE WINDOW (e.g. ~7 days)
                            |
              +-------------+-------------+
              |                           |
      nobody submits a            somebody submits a
      fraud proof --               FRAUD PROOF, showing
      the state root                the posted state root
      is now FINAL                  was WRONG
                                            |
                                            v
                                  L1 re-executes just the
                                  disputed step, reverts the
                                  bad state root, and
                                  PENALIZES whoever posted it
```

The trade-off is right there in "challenge window": funds withdrawn from an optimistic rollup back to L1 are not instantly final — a withdrawal must typically wait out the full challenge window (so that, if the batch it depends on turns out to be fraudulent, there's still time to catch and revert it) before L1 will honor it. This is optimistic rollups' central practical downside: fast, cheap execution off-chain, but a real delay before funds are provably, finally safe to withdraw back to L1.

---

## 7. Zero-Knowledge Rollups

A **zero-knowledge rollup** (**ZK-rollup**) solves the same "batch transactions off-chain, post a summary" problem with a fundamentally different guarantee: instead of *assuming* correctness and giving a window to dispute it, a ZK-rollup posts a **validity proof** — a small cryptographic proof (a zk-SNARK or zk-STARK, families of proof systems whose math is well beyond this course's scope, but whose *property* is simple to state) that mathematically demonstrates the new state root is the correct result of executing the batch, full stop, with no room for it to be wrong.

```
   ZERO-KNOWLEDGE ROLLUP

   Off-chain:  1,000 transactions executed, new state root
               computed, AND a validity proof generated
               proving "this state root is what you get from
               correctly executing these transactions"

   Posted to L1:  [ compressed tx data | new state root | validity proof ]
                            |
                            v
                  L1 VERIFIES THE PROOF (cheap and fast,
                  even though the proof attests to a
                  computation that would have been
                  expensive to redo directly)
                            |
                            v
                  Proof valid?  -->  state root accepted
                                     IMMEDIATELY -- no
                                     challenge window needed,
                                     because correctness was
                                     already mathematically
                                     proven, not merely assumed
```

The remarkable property of a validity proof — the part that makes "zero-knowledge" the right name — is that L1 can verify the proof is correct *without re-executing any of the underlying 1,000 transactions itself*, and without learning anything about their contents beyond "yes, this state transition is valid." This is what removes the need for a challenge window entirely: there is nothing to dispute, because the proof already constitutes airtight mathematical evidence. The trade-off moves elsewhere instead — generating a validity proof is computationally expensive (though verifying one is cheap), and the cryptography involved is considerably more sophisticated to implement correctly than optimistic rollups' comparatively simple "assume correct, allow fraud proofs" approach.

Both rollup styles from this section and the last are, as of this course's writing, in active production use, not purely theoretical designs — several optimistic rollups and several ZK-rollups process real transaction volume on top of Ethereum today, and it's worth recognizing why both flavors have found lasting adoption rather than one strictly displacing the other: optimistic rollups were simpler to build correctly earlier on, since fraud proofs only need to re-run a disputed computation using tooling that already existed (an ordinary virtual machine, much like Volume 9's own), while ZK-rollups required (and continue to require) genuinely novel cryptographic engineering to make proof generation fast enough to be practical for general-purpose computation, not just narrow, specialized calculations. The gap between them has been closing over time as ZK proving systems mature, but as of today both remain live, reasonable choices depending on a project's tolerance for withdrawal delay versus its tolerance for cryptographic implementation complexity — precisely the trade-off Section 10 makes explicit.

---

## 8. The Arithmetic of a Rollup Batch

It's worth working through, with made-up but realistic-shaped numbers, exactly *why* batching is where a rollup's cost savings actually come from — the mechanism is simpler than the cryptography surrounding it suggests. Imagine each individual GoChain transaction, submitted directly to L1 the ordinary way (Chapter 70's `POST /transactions`), costs a fixed amount of "on-chain space" — call it 200 bytes once serialized, matching roughly the shape of `core.Transaction` from Volume 5. Submitting 1,000 such transactions directly to L1, one at a time, costs L1 space for all 1,000 of them individually:

```
   1,000 transactions submitted DIRECTLY to L1

     1,000 x 200 bytes  =  200,000 bytes of L1 block space consumed
     1,000 separate transaction validations performed BY L1 directly
```

A rollup instead executes those same 1,000 transactions off-chain, and posts back to L1 only a *compressed* representation of what changed, plus the new state root:

```
   1,000 transactions batched through a ROLLUP

     Off-chain: full 200,000 bytes of transaction data,
                executed and verified off-chain

     Posted to L1:  compressed tx data (often a small
                     fraction of the original size, since much
                     of a batch of transactions is repetitive --
                     the same handful of addresses transacting
                     repeatedly compresses well)
                     + one new state root (a single, fixed-size
                       hash, regardless of batch size)
                     + (ZK-rollup only) one validity proof
                       (also fixed-size, regardless of batch size)

     L1 space consumed: a SMALL, ROUGHLY CONSTANT overhead
     PER BATCH, not per transaction
```

The key insight this arithmetic makes concrete: as the batch size grows from 1,000 transactions to 10,000, the *per-transaction* share of that roughly-constant per-batch overhead keeps shrinking, which is exactly why real rollups batch as many transactions together as they reasonably can before posting — the economics only get better with larger batches, up to whatever practical limits (compression ratios, proof-generation time for ZK-rollups, or acceptable latency before a batch is posted) apply. This is also precisely why a rollup's state root (Section 6 and 7's diagrams) is the one piece of information that must be posted no matter how large the batch grows — Chapter 57's Merkle-Patricia trie root is a single, fixed-size hash regardless of how much state it commits to, which is exactly the property that makes "one small posting settles an arbitrarily large batch" possible at all.

---

## 9. What Each Approach Still Verifies on the Main Chain

It's worth being explicit about exactly what work each technique leaves for L1 to do, since that's precisely what determines how much security each one actually inherits from the base chain versus how much it's trusting elsewhere.

```
                          WHAT L1 STILL VERIFIES

   SHARDING              Each shard's own validators verify their own
                         shard's blocks fully; L1 (or a coordinating
                         beacon chain) verifies cross-shard receipts
                         and shard-level commitments, NOT every
                         individual transaction directly.

   STATE CHANNELS        L1 verifies exactly two transactions per
                         channel (open and close) plus, RARELY, a
                         dispute if one party tries to submit a
                         stale, outdated signed balance.

   OPTIMISTIC ROLLUP     L1 verifies nothing about the batch's
                         correctness UP FRONT -- it only verifies a
                         fraud proof IF one is submitted during the
                         challenge window. Absence of a challenge is
                         treated as proof of correctness.

   ZK-ROLLUP             L1 verifies a validity proof for EVERY
                         single batch, immediately -- cheap to check,
                         but mathematically airtight, with nothing
                         optimistic or assumed about it.
```

This table is the single most useful thing to internalize from this whole chapter: "Layer 2" is not one technique with one security model — it's a spectrum from "L1 barely has to do anything, and trusts a challenge mechanism" (optimistic rollups) to "L1 does a small but non-negotiable amount of cryptographic verification on every single batch" (ZK-rollups), with state channels sitting in a different corner entirely (minimal on-chain footprint, but restricted to pre-agreed participants).

Sharding's own row in the table above deserves one further distinction, tying back to Section 3: a *transaction*-sharded design still has L1 (or a coordinating beacon layer) ultimately anchor every shard's block headers, even though it never re-executes each shard's individual transactions itself — meaning L1's actual verification burden per shard looks a great deal like "verify a compact commitment" (a header, a state root) rather than "re-run everything," the same underlying pattern every Layer 2 technique in this chapter also relies on. This is a useful thing to notice on its own: sharding and Layer 2 are usually presented as two entirely separate categories of scaling technique, but at the level of "what does the ultimate source of truth actually verify," both consistently boil down to the same move — replace "re-execute everything" with "verify a small, fixed-size commitment to a larger computation that happened elsewhere."

---

## 10. What Scaling Techniques Actually Trade Away

None of the techniques in this chapter are free lunches — each is a real, considered trade-off, not a strictly-better replacement for a single, simple chain:

- **Sharding** trades a simpler, whole-network security model for parallelism, but introduces cross-shard communication complexity and — critically — means any *single* shard is now watched over by a smaller subset of the network's total validators than the whole chain would provide, which is a real reduction in that shard's individual security margin unless carefully mitigated (frequent, randomized reassignment of validators to shards is one common mitigation, precisely to prevent an attacker from concentrating influence on one shard).
- **State channels** trade generality (only pre-committed participants can transact) for near-zero on-chain footprint and instant off-chain finality between those participants.
- **Optimistic rollups** trade withdrawal speed (the challenge window) for comparatively simple, well-understood cryptography and lower off-chain computational cost.
- **Zero-knowledge rollups** trade implementation complexity and proof-generation cost for immediate finality and no challenge window at all.

The throughline across every single one of these trade-offs is that **none of them make a chain unconditionally faster for free** — every technique in this chapter is a deliberate exchange of one property (simplicity, uniform security, generality, or instant finality) for another (parallelism, speed, or lower on-chain cost). Chapter 80 picks this exact thread back up as part of a broader decision framework: recognizing *when* a project's actual requirements justify reaching for any of these techniques at all, versus when a single, well-run chain (much like GoChain itself) is already more than sufficient.

---

## Summary

- GoChain's throughput ceiling comes directly from every honest node re-executing every transaction — the same property that makes the chain trustworthy also caps how much it can ever process, no matter how many nodes join.
- Sharding splits the network into parallel sub-chains (shards), each processing a fraction of total transactions independently, at the cost of a genuinely hard cross-shard communication problem and a smaller effective validator set securing any one shard.
- Real sharding designs actually layer three distinct techniques together — network sharding (who gossips with whom), transaction sharding (which shard processes which transactions), and state sharding (which shard stores which slice of account/UTXO state) — with state sharding consistently the hardest to get right.
- Layer 2 leaves the base chain (L1) unchanged and builds a separate system on top that does most of the work off-chain, touching L1 only when its security guarantees are specifically needed.
- State channels let a fixed, pre-agreed set of participants transact off-chain indefinitely, touching the chain only to open and close the channel, secured by the ability to unilaterally submit the latest signed balance at any time.
- Optimistic rollups post batched results to L1 assuming correctness, backed by a challenge window during which anyone can submit a fraud proof to revert an incorrect result — trading a withdrawal delay for comparatively simple cryptography.
- Zero-knowledge rollups post a validity proof alongside every batch, mathematically demonstrating correctness so L1 can verify it immediately and cheaply, with no challenge window, at the cost of expensive proof generation and more sophisticated cryptography.
- A rollup's on-chain cost per batch is roughly constant regardless of batch size (a compressed transaction blob plus one fixed-size state root, and for ZK-rollups one fixed-size proof), which is exactly why real rollups batch as many transactions together as practical — the per-transaction cost keeps shrinking as batches grow.
- Across sharding, state channels, and both rollup types, what varies is precisely how much verification work L1 still does directly versus how much it delegates to a challenge mechanism, a cryptographic proof, or a smaller subset of validators.
- Every scaling technique in this chapter is a deliberate trade-off — parallelism, speed, or lower cost in exchange for added complexity, a security-margin reduction, a withdrawal delay, or restricted participation — never an unconditional improvement.

---

## Exercises

### Easy

1. In your own words, explain why adding more nodes to GoChain's existing network (as built in Volume 7) improves decentralization and fault tolerance, but does not, by itself, raise the network's transaction throughput.
2. What is the key difference between how an optimistic rollup and a zero-knowledge rollup convince Layer 1 that a batch of off-chain transactions was executed correctly?
3. Why are state channels well suited to a payment channel between two frequent counterparties, but not to a one-time payment to a stranger you've never transacted with before?

### Medium

4. Draw your own ASCII sequence diagram (not copied from this chapter) showing Alice and Bob opening a state channel, exchanging three off-chain payments, and one of them submitting a stale (outdated) balance to try to cheat during channel close — and show how the chain's dispute mechanism catches it.
5. Research one real, currently-operating optimistic rollup and one real ZK-rollup (built on any base chain, not necessarily Ethereum), and report each one's actual challenge window length (for the optimistic rollup) and roughly how proof generation time compares to normal transaction execution time (for the ZK-rollup).
6. Section 10 mentions that sharding reduces the effective validator set securing any one shard, and that randomized validator reassignment is one mitigation. Explain, in your own words, how randomly and frequently reassigning validators across shards makes it harder for an attacker to concentrate malicious control over any single shard.
7. Using Section 8's arithmetic, suppose a rollup's per-batch overhead on L1 is a fixed 300 bytes (state root plus fixed overhead), regardless of batch size. Calculate the per-transaction L1 footprint at batch sizes of 100, 1,000, and 10,000 transactions, and explain in a sentence why this specific shape of cost curve is what makes larger batches economically attractive up to a point.

### Hard

8. Sketch, at a conceptual level (no Go code required), how you would add a simple two-shard design to GoChain: how would you decide which shard an address's transactions belong to, and what minimal cross-shard mechanism would you need for an address on Shard A to pay an address on Shard B? Identify, using Section 3's three-way distinction, which of network, transaction, or state sharding your design actually implements.
9. Research zk-SNARKs or zk-STARKs (pick one) at a conceptual level — you do not need to understand the underlying mathematics — and explain, in plain language, what "succinct" and "non-interactive" mean in this context, and why both properties matter for a rollup posting proofs to a base chain that must verify them cheaply.
10. Compare, in a short written analysis, the total time from "transaction submitted to a Layer 2" to "funds are irrevocably final on Layer 1" across a state channel, an optimistic rollup, and a ZK-rollup. Under what real-world circumstances would the optimistic rollup's challenge-window delay be an acceptable trade-off, and under what circumstances would it not be?
