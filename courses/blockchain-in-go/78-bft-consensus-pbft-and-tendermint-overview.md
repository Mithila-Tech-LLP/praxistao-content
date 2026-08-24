# Chapter 78: BFT Consensus — PBFT and Tendermint, Overview

Both consensus engines GoChain has built so far — proof of work in Volume 4 and proof of stake in Chapter 77 — share a quiet assumption: agreement emerges gradually, through probability. A miner's block *probably* won't be reverted once enough blocks are stacked on top of it; a validator's proposal is *probably* legitimate because selection was fair. Byzantine Fault Tolerant (BFT) consensus takes a different, older approach: a fixed, known set of validators explicitly vote, in structured rounds, until a supermajority agrees — and once that supermajority is reached, the decision is final immediately, not "probably final after enough confirmations." This chapter is conceptual grounding, not a fourth engine for GoChain to implement: you'll walk through PBFT's three-phase vote and Tendermint's round-based model closely enough to recognize them (and reason about their trade-offs) the moment you meet a real system — Hyperledger Fabric, Cosmos, or any other — built on top of one.

## Table of Contents

1. [What "Byzantine Fault Tolerant" Actually Means](#1-what-byzantine-fault-tolerant-actually-means)
2. [Why PoW and PoS Aren't BFT in This Classical Sense](#2-why-pow-and-pos-arent-bft-in-this-classical-sense)
3. [PBFT's Three-Phase Vote](#3-pbfts-three-phase-vote)
4. [PBFT Sequence Diagram, Step by Step](#4-pbft-sequence-diagram-step-by-step)
5. [Tendermint's Round-Based Model](#5-tendermints-round-based-model)
6. [Tendermint Sequence Diagram, Step by Step](#6-tendermint-sequence-diagram-step-by-step)
7. [Comparing PBFT, Tendermint, PoW, and PoS](#7-comparing-pbft-tendermint-pow-and-pos)
8. [Recovering From a Faulty Primary: View Changes](#8-recovering-from-a-faulty-primary-view-changes)
9. [Where GoChain Could Plug In a BFT Engine](#9-where-gochain-could-plug-in-a-bft-engine)
10. [Beyond PBFT and Tendermint: A Wider Family](#10-beyond-pbft-and-tendermint-a-wider-family)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What "Byzantine Fault Tolerant" Actually Means

The name comes from the **Byzantine Generals Problem**, a thought experiment from a 1982 computer science paper: several generals, commanding different divisions surrounding an enemy city, must unanimously agree to attack or retreat, communicating only by messenger. The catch — some generals might be traitors, sending contradictory messages to different peers specifically to cause the loyal generals to act on conflicting information. A **Byzantine fault** is exactly this: a component that doesn't just crash or go silent (a "fail-stop" fault, the easy case), but actively lies, sends different answers to different peers, or behaves however is worst for the system — the digital equivalent of a traitorous general.

```
FAIL-STOP FAULT                          BYZANTINE FAULT

  Node crashes or goes offline             Node stays online but LIES:
  Other nodes: "node 3 is silent,           - tells node A "yes"
  let's proceed without it"                - tells node B "no"
                                            - sends a validly-signed
  Easy to detect, easy to work around        message for the WRONG value
                                            Much harder: from the outside,
                                            a lying node looks identical
                                            to an honest one, message by
                                            message
```

A consensus algorithm is **Byzantine Fault Tolerant** if it still reaches correct, consistent agreement even when some fraction of participants behave this way — not just crash, but actively try to sabotage the outcome. The classical result, proven in that same 1982 paper, is the "one-third bound": a BFT protocol can tolerate up to **f** faulty (including actively malicious) nodes only if the total number of participants is at least **3f + 1** — put differently, strictly more than two-thirds of participants must be honest for agreement to be guaranteed. This single number — 3f+1 — reappears, unchanged, in essentially every classical BFT protocol, including both systems this chapter covers.

Why exactly `3f + 1`, and not, say, `2f + 1` the way a simpler fail-stop-only system can get away with? The intuition is worth internalizing, because it explains every BFT design choice that follows. Suppose the network splits into three groups during a vote: some honest nodes, some Byzantine nodes actively lying, and a *second* set of honest nodes that simply haven't heard from the first group yet (perhaps their messages are just slow to arrive, not lost). A Byzantine node can exploit this by telling the slow-to-hear honest nodes one thing, and the already-informed honest nodes another. For the protocol to still guarantee that no two *honest* nodes ever disagree, the "definitely honest and already agreed" group needs to outnumber "the Byzantine nodes plus every honest node that could plausibly still be confused" — working through that arithmetic is exactly what produces the `3f + 1` bound (`f` Byzantine, up to `f` slow-but-honest, and still needing a strict majority of the *remaining* honest nodes to settle it, which arithmetically requires `f + 1` more). You do not need to reproduce this proof to use the result correctly — just remember that it is not an arbitrary safety margin someone picked by feel, but the tightest bound the problem's own logic allows.

---

## 2. Why PoW and PoS Aren't BFT in This Classical Sense

It's tempting to assume "BFT" is just a fancier synonym for "handles bad actors," which would make proof of work and proof of stake sound BFT too — after all, Volume 4 and Chapter 77 both discussed tolerating dishonest miners and validators. The distinction that matters is **how finality is reached**:

```
PoW / PoS                                  CLASSICAL BFT (PBFT, Tendermint)

  "Finality" is PROBABILISTIC:               "Finality" is IMMEDIATE and
  a block becomes safer the more              ABSOLUTE:
  blocks are mined on top of it,               once 2/3+ of validators have
  but in principle a long-enough               explicitly voted to commit a
  reorg could still undo it                    block, it is final -- full
                                                stop, no "wait a few more
  Open, permissionless membership:              blocks to be safer"
  anyone can join as a miner or
  validator at any time                       Requires a KNOWN, FIXED set of
                                                validators who explicitly
  No explicit voting round --                  vote in structured rounds
  agreement is implicit, via
  "longest chain" or "highest stake          Typically PERMISSIONED: the
  weight" rules                                validator set is known in
                                                advance and changes only
                                                through an explicit process
```

This is not "BFT is strictly better" — it's a different point in the design space, trading permissionless openness for immediate, absolute finality. Chapter 80 returns to this trade-off directly as part of a broader decision framework; for now, the important thing is recognizing that "tolerates a dishonest minority" is necessary but not sufficient for something to be classically BFT — the explicit, round-based voting protocol is the defining feature, and that's exactly what Sections 3 and 5 walk through.

---

## 3. PBFT's Three-Phase Vote

**PBFT** (Practical Byzantine Fault Tolerance), published in 1999, was the algorithm that made BFT consensus practical enough for real systems rather than a purely theoretical result — it underlies the general design philosophy of permissioned chains like Hyperledger Fabric. Every round has one designated **primary** (the validator proposing this round's value — conceptually similar to proof of stake's selected proposer from Chapter 77) and proceeds through exactly three phases of message exchange among all validators (often called **replicas** in PBFT's own vocabulary) before anything is considered committed.

An everyday analogy: imagine a jury of twelve, but one where a "guilty" verdict requires not just a majority, but explicit, repeated confirmation — first, one juror proposes a verdict; second, every juror announces to *every other juror* whether they'd accept that verdict; third, once a juror sees enough of their peers accepting it, they announce they're personally *committing* to it, and only once a juror has seen enough commit announcements do they treat the verdict as truly final. Three separate rounds of "everyone tells everyone" is what makes it safe against a lying juror — a single round of votes could be manipulated by a minority sending different answers to different peers, but three cross-checked rounds cannot be, as long as at least two-thirds of the jury is honest.

```
Phase 1: PRE-PREPARE                Phase 2: PREPARE                 Phase 3: COMMIT

  Primary proposes a value            Every replica broadcasts         Every replica broadcasts
  to all replicas                     "PREPARE" to every OTHER         "COMMIT" once it has seen
                                       replica, confirming it saw       enough matching PREPAREs
  Primary --> R1, R2, R3               the SAME pre-prepare              (2f+1 of them)

                                       This is the step that            Once a replica sees 2f+1
                                       CROSS-CHECKS that no two          matching COMMITs, the
                                       honest replicas were shown        value is final -- it will
                                       different proposed values          never be reverted
```

The reason three phases (not one, not two) are required comes directly from tolerating actively lying nodes, not just crashed ones: a single "primary proposes, everyone votes" round is exactly the failure mode the Byzantine Generals Problem describes — a lying primary (or a lying minority of replicas) could tell different honest replicas different things, and a single round gives no way to detect that discrepancy. The PREPARE phase's entire job is making every replica cross-check with every *other* replica ("did you also see this exact proposal?") before anyone commits to anything, which is precisely what makes a two-thirds-honest supermajority sufficient to guarantee agreement even when the remaining third is actively malicious.

---

## 4. PBFT Sequence Diagram, Step by Step

Here is a concrete round with one primary and three other replicas (four validators total, tolerating up to `f=1` Byzantine fault, since 4 = 3(1)+1):

```
        Client       Primary        R1          R2          R3
          |             |            |           |           |
          | request     |            |           |           |
          |------------>|            |           |           |
          |             |            |           |           |
          |             | PRE-PREPARE (proposed value V)      |
          |             |----------->|---------->|---------->|
          |             |            |           |           |
          |             |         PREPARE (V)  PREPARE (V) PREPARE (V)
          |             |<-----------|<----------|<----------|
          |             |----------->|---------->|---------->|
          |             |            |           |           |
          |         [ each replica now has 2f+1 = 3 matching  ]
          |         [ PREPAREs for V -- "prepared" state       ]
          |             |            |           |           |
          |             |         COMMIT (V)   COMMIT (V)   COMMIT (V)
          |             |<-----------|<----------|<----------|
          |             |----------->|---------->|---------->|
          |             |            |           |           |
          |         [ each replica now has 2f+1 = 3 matching  ]
          |         [ COMMITs for V -- V is FINAL, no reversal]
          |             |            |           |           |
          |          reply (V committed)         |           |
          |<------------|            |           |           |
```

Notice the client only ever hears back once, after COMMIT — everything in the PRE-PREPARE and PREPARE phases is internal cross-checking among validators, invisible to whoever originally submitted the request. Also notice this diagram assumes a healthy primary; PBFT's full specification includes a **view change** protocol (electing a new primary if the current one appears to be stalling or lying), which is exactly the kind of detail this conceptual overview intentionally leaves out — real PBFT deployments spend a great deal of their complexity budget on correctly handling exactly that case.

---

## 5. Tendermint's Round-Based Model

**Tendermint**, the consensus engine powering the Cosmos ecosystem, is a more modern BFT design built specifically for blockchains (PBFT predates blockchains entirely and was originally designed for fault-tolerant file systems). Tendermint restructures the same fundamental idea — propose, then two rounds of voting, then commit — into a repeating **round** structure explicitly designed to keep producing blocks continuously, one after another, rather than reaching a single one-off decision the way PBFT's original formulation did.

```
Each Tendermint ROUND has three steps, and produces (at most) one block:

  PROPOSE                    PREVOTE                    PRECOMMIT

  This round's proposer        Every validator            Every validator
  (chosen by round-robin       broadcasts a PREVOTE        broadcasts a PRECOMMIT
  or stake-weighted            for the proposed block       once it has seen 2/3+
  selection, similar in         (or a "nil" prevote if       matching PREVOTES
  spirit to Chapter 77's        it saw nothing valid)
  SelectValidator) proposes                                Once a validator sees
  a block                                                    2/3+ matching PRECOMMITs,
                                                              the block is COMMITTED
                                                              and the chain moves to
                                                              the next height
```

The naming echoes PBFT deliberately (PREVOTE parallels PREPARE, PRECOMMIT parallels COMMIT) because it is, structurally, the same three-phase idea — propose, cross-check, commit — adapted for a chain of blocks rather than a single isolated decision. The key addition Tendermint makes is an explicit **round number** within each block height: if a proposer fails to produce a valid block, or the network fails to gather enough matching votes within a timeout, validators simply move to the *next round* at the *same height* with a new proposer, rather than the whole protocol stalling indefinitely.

```
HEIGHT 100

  Round 0: Validator A proposes -- timeout, not enough PRECOMMITs gathered
  Round 1: Validator B proposes -- timeout again
  Round 2: Validator C proposes -- 2/3+ PRECOMMIT -- block 100 COMMITTED

HEIGHT 101

  Round 0: Validator D proposes -- 2/3+ PRECOMMIT -- block 101 COMMITTED
  (most rounds succeed on the first try in a healthy network; round
   failures are the exception, handled gracefully rather than fatally)
```

This round-and-retry structure is Tendermint's practical answer to the "what if the primary is stalling" problem PBFT's view-change protocol also had to solve — instead of a special-case recovery procedure, "try the next round with a new proposer" is simply the normal, everyday path every height goes through, whether the first proposer succeeds immediately (the common case) or several rounds are needed (the rare case).

---

## 6. Tendermint Sequence Diagram, Step by Step

A single successful round, height 100, four validators (again tolerating `f=1`):

```
      Proposer(A)      B            C            D
          |            |            |            |
          | PROPOSE (block for height 100)         |
          |----------->|----------->|----------->|
          |            |            |            |
          |         PREVOTE      PREVOTE      PREVOTE
          |<-----------|<-----------|<-----------|
          |----------->|----------->|----------->|
          |            |            |            |
      [ each validator now has 2/3+ matching PREVOTEs -- "polka" ]
          |            |            |            |
          |         PRECOMMIT    PRECOMMIT    PRECOMMIT
          |<-----------|<-----------|<-----------|
          |----------->|----------->|----------->|
          |            |            |            |
      [ each validator now has 2/3+ matching PRECOMMITs ]
      [ block 100 is COMMITTED -- move to height 101      ]
```

Tendermint's own terminology calls a set of 2/3+ matching PREVOTEs for the same block a **polka** — a term borrowed informally to describe "enough agreement to be worth acting on, but not yet final." A validator only casts its own PRECOMMIT once it has personally observed a polka, mirroring exactly the same "cross-check before committing" discipline PBFT's PREPARE phase enforces. Once 2/3+ matching PRECOMMITs exist, the block is final immediately — a Tendermint chain has no equivalent of "wait six confirmations to be safe" the way proof-of-work chains do, because finality was the entire point of the voting rounds that already happened.

---

## 7. Comparing PBFT, Tendermint, PoW, and PoS

| | Proof of Work (Ch. 25) | Proof of Stake (Ch. 77) | PBFT | Tendermint |
|---|---|---|---|---|
| **Membership** | permissionless, anyone with hardware | permissionless, anyone who stakes | permissioned, known validator set | typically permissioned (staked validator set, but structurally BFT) |
| **Finality** | probabilistic (more confirmations = safer) | probabilistic (similar to PoW) | immediate and absolute once committed | immediate and absolute once committed |
| **Fault tolerance** | tolerates a dishonest minority of hash power | tolerates a dishonest minority of stake | tolerates up to f Byzantine faults out of 3f+1 | tolerates up to f Byzantine faults out of 3f+1 |
| **Communication** | no explicit voting; implicit via longest chain | no explicit voting; implicit via stake-weighted selection | explicit multi-round voting, every validator to every other | explicit multi-round voting, every validator to every other |
| **Throughput/latency** | limited by block interval and confirmation depth | faster block times than PoW, still probabilistic | fast, but O(n²) message complexity limits validator count | fast, same O(n²) limitation, in practice tuned for tens to low hundreds of validators |
| **Real-world example** | Bitcoin | Ethereum (post-Merge, Chapter 82), many app chains | Hyperledger Fabric's design philosophy | Cosmos SDK chains |

The **O(n²) message complexity** row is worth sitting with: every validator broadcasting to every other validator, twice (PREPARE and COMMIT, or PREVOTE and PRECOMMIT), means the total number of messages grows with the *square* of the validator count. This is precisely why classical BFT protocols favor tens to a few hundred known validators rather than the tens of thousands of anonymous participants a permissionless proof-of-work or proof-of-stake network can support — it's a direct, structural trade-off, not an implementation shortcoming waiting to be optimized away.

A practical way to keep all four side by side in your head, the next time you're reading about a real chain's design: ask two questions. First, "can I, a stranger, join this network as a validator without anyone's permission?" — a "yes" points toward PoW or PoS; a "no" points toward PBFT or Tendermint. Second, "does this chain's documentation talk about 'confirmations' or 'finality depth,' or does it talk about blocks being final the instant they're produced?" — the former is PoW/PoS's probabilistic finality; the latter is classical BFT's immediate finality. Answering just those two questions correctly identifies which quadrant of this table you're looking at before you've read a single further detail of the protocol's internals.

---

## 8. Recovering From a Faulty Primary: View Changes

Sections 3 and 5 both mentioned, in passing, that a healthy round assumes the primary (or proposer) actually does its job — proposes a single, valid value promptly. Real networks cannot assume that: a primary might crash, might be slow, or might be actively malicious, deliberately withholding its proposal or sending conflicting ones to different replicas. PBFT's answer to this is called a **view change** — a **view** is PBFT's term for "the current assignment of which replica is primary," and a view change is the structured process of replacing a stalling or misbehaving primary with the next one in line, without ever losing agreement about what (if anything) was already committed.

Think of it like a company's succession plan for an out-of-office CEO who has stopped responding to anything time-sensitive: the board doesn't wait forever, and it doesn't just let chaos reign either — there's a specific, pre-agreed process (a vote among the board) for handing authority to the next person in the succession line, with clear rules for what counts as "already decided" versus "still open" at the moment of handover.

```
NORMAL OPERATION                          PRIMARY STALLS OR MISBEHAVES

  Primary proposes promptly                 Replicas start a timer the
  Replicas PREPARE, COMMIT                  moment they expect a proposal
  as in Section 4's diagram                 (e.g. after receiving a client
                                              request) -- if PRE-PREPARE
                                              never arrives before the
                                              timer expires...
                                                        |
                                                        v
                                            Replicas broadcast a
                                            VIEW-CHANGE message: "I no
                                            longer trust this primary,
                                            move to the next view"
                                                        |
                                                        v
                                            Once enough VIEW-CHANGE
                                            messages (2f+1) are collected,
                                            the NEXT replica in line
                                            becomes primary, broadcasts a
                                            NEW-VIEW message, and normal
                                            operation (PRE-PREPARE onward)
                                            resumes under the new primary
```

The genuinely tricky part of a real view change — deliberately left out of the simplified diagram above — is making sure nothing that was already "prepared" (had 2f+1 matching PREPAREs, per Section 3) under the old primary gets silently lost or contradicted under the new one. PBFT's actual view-change protocol requires each replica to include, in its VIEW-CHANGE message, proof of anything it had already prepared, so the new primary is forced to re-propose those same values first, before handling anything new — preserving the guarantee that once enough replicas agreed something was safe, that agreement can never quietly evaporate just because leadership changed hands.

Tendermint's answer, introduced back in Section 5, is structurally simpler precisely because it was designed after PBFT, with blockchains specifically in mind: rather than a special-case recovery procedure invoked only when something goes wrong, "move to the next round with a new proposer" is simply the *normal* path every single block height takes, whether the first proposer succeeds immediately or several rounds are needed. Tendermint traded PBFT's separate, more intricate view-change machinery for a design where "the leader might need replacing" is treated as routine, not exceptional — a good example of how a later protocol, built for a narrower and better-understood problem (continuously producing a chain of blocks, rather than PBFT's original, more general state-machine-replication goal), can simplify a genuinely hard problem by not trying to solve the fully general version of it.

---

## 9. Where GoChain Could Plug In a BFT Engine

Recall from Chapter 77 that `consensus.Engine` is deliberately narrow:

```go
type Engine interface {
	Mine(b *core.Block) (nonce uint64, hash []byte)
	Validate(b *core.Block) bool
}
```

A BFT engine would strain this interface's assumptions in a way proof of stake did not, and it's worth understanding exactly why, conceptually, even without building it. Both `ProofOfWork.Mine` and `ProofOfStake.Mine` are **single-node operations**: one node, on its own, can compute a nonce or select a proposer and return an answer immediately. A BFT engine's "mine," by contrast, is fundamentally a **multi-round, multi-party conversation** — a single node cannot, alone, produce a committed block; it must exchange PRE-PREPARE/PREPARE/COMMIT (or PROPOSE/PREVOTE/PRECOMMIT) messages with every other validator over the network first, using machinery that looks much more like Volume 7's P2P networking than like a local computation.

```
   ProofOfWork.Mine / ProofOfStake.Mine        A hypothetical BFT engine's "Mine"

     one function call,                          would require:
     one goroutine,                              - sending messages to every
     returns a value                                other validator (network I/O)
     synchronously                               - waiting for a threshold of
                                                    responses (blocking on OTHER
                                                    nodes' behavior, not just
                                                    local computation)
                                                  - handling timeouts and view/
                                                    round changes
                                                  ...none of which fits neatly
                                                  inside a synchronous function
                                                  that returns (nonce, hash)
```

This is exactly why real BFT implementations (Tendermint Core included) are not small pluggable "engines" behind a two-method interface — they are entire subsystems with their own networking, message queues, and state machines, deeply integrated with how blocks are proposed and gossiped in the first place. Recognizing *why* the shape doesn't fit is itself a valuable, transferable piece of understanding: it's the same reason production systems that want BFT guarantees (Cosmos SDK chains built on Tendermint, for instance) are architected around their consensus engine from the very beginning, rather than retrofitting it in behind an interface designed for something structurally simpler.

Concretely, if a future version of GoChain genuinely wanted BFT consensus, the change would ripple far past `consensus`. Volume 7's `network.Node` (Chapter 46) would need three brand-new message types alongside its existing `version`, `getblocks`, `inv`, `getdata`, `block`, and `tx` (Chapter 45) — something like `preprepare`/`propose`, `prepare`/`prevote`, and `commit`/`precommit` — each one addressed to the validator set specifically, not gossiped broadly to every peer the way a new transaction is. `core.Blockchain.MineBlock` (Chapter 25) would need to become asynchronous and stateful, tracking "which round are we in, for which height, and how many matching votes have I collected so far" across multiple incoming network messages, rather than returning a finished block from one direct function call. And the storage layer (Volume 8) would need a place to durably record "what did I personally vote for at height H, round R" — precisely so that if a validator restarts mid-round, it can never accidentally vote for two different values at the same height, which is exactly the kind of behavior that would look, to every other honest node, indistinguishable from a genuine Byzantine fault.

```
   WHAT PoS ADDED TO GoChain (Ch. 77)        WHAT BFT WOULD ADD (conceptual)

   - consensus.ProofOfStake, a new           - new P2P message types
     type implementing Engine                  (propose/prevote/precommit)
   - SelectValidator, Slash: pure             - MineBlock becomes stateful
     functions, no network I/O                  and asynchronous, spanning
   - Everything else in core, network,          many incoming messages
     storage: UNCHANGED                       - a durable "what did I vote
                                                 for" record per height/round
                                               - core.Blockchain, network.Node,
                                                 AND storage.Store would all
                                                 need real changes
```

---

## 10. Beyond PBFT and Tendermint: A Wider Family

PBFT (1999) and Tendermint (its blockchain-era successor) are the two most useful BFT protocols to recognize by name, but they are two points in a wider, still-active area of research and production systems, worth knowing exist even without walking through each one's mechanics in detail:

```
   PBFT (1999)              Tendermint (2014)         HotStuff (2018)

   The original             Blockchain-specific,       Simplifies BFT further:
   "practical" BFT           round-based restructuring  a LINEAR message pattern
   protocol -- O(n^2)         of the same 3-phase idea    (through a rotating
   messages, view-change      Powers Cosmos SDK chains    leader) instead of every
   for leader replacement                                 replica messaging every
                                                            other replica directly,
                                                            reducing communication
                                                            overhead -- adopted by
                                                            Diem/Libra and several
                                                            newer chains
```

**HotStuff**, published in 2018 and later adopted by Meta's (then Facebook's) Diem/Libra project and several other chains since, is worth naming specifically because it addresses the exact O(n²) bottleneck Section 7 flagged as PBFT and Tendermint's shared structural limitation: instead of every replica broadcasting directly to every other replica in each phase, HotStuff routes votes *through the current leader*, who aggregates them into a single compact certificate before relaying it onward — reducing the total message complexity from quadratic to linear in the number of validators, at the cost of a slightly deeper chain of phases per decision. This is a good illustration that "classical BFT" is not a single frozen design from 1999, but an active area where the core safety guarantee (the `3f+1` bound from Section 1) stays fixed while the *engineering* of how validators communicate to reach it keeps improving — precisely the kind of detail that separates recognizing a name from actually understanding what problem it solves.

```
   MESSAGE PATTERN: PBFT / TENDERMINT              MESSAGE PATTERN: HOTSTUFF

     every replica  ------->  every OTHER            every replica  ------->  LEADER
     replica, in EACH phase                          (votes aggregated into
     (PREPARE, then COMMIT --                         one certificate)
     each an all-to-all broadcast)
                                                       LEADER  ------->  every replica
     O(n^2) total messages per                        (broadcasts the aggregated
     decision                                          certificate onward)

                                                       O(n) total messages per
                                                       phase -- LINEAR, not
                                                       quadratic, in validator count
```

None of this changes anything about what GoChain would need to do to actually run a BFT engine (Section 9's conclusion holds regardless of which specific BFT protocol you'd pick) — the point of naming HotStuff here is purely so that, when you eventually read a real chain's whitepaper or documentation and see "we use a HotStuff-based consensus" or "we use Tendermint," you recognize both as members of the same classical-BFT family this chapter introduced, differing in communication pattern and engineering detail, not in their fundamental safety guarantee.

---

## Summary

- A Byzantine fault is a component that actively lies or behaves inconsistently toward different peers, not merely one that crashes — Byzantine Fault Tolerant consensus specifically guarantees correct agreement even when some participants behave this way.
- The classical BFT bound is `3f + 1`: a protocol can tolerate `f` Byzantine faults only when there are at least `3f + 1` total participants, meaning strictly more than two-thirds must be honest.
- PoW and PoS tolerate a dishonest minority too, but reach agreement *implicitly* and *probabilistically* (longest chain, stake-weighted selection); classical BFT protocols reach agreement *explicitly*, through multi-round voting, with immediate and absolute finality once a supermajority commits.
- PBFT's three phases — PRE-PREPARE (a primary proposes), PREPARE (replicas cross-check that they all saw the same proposal), and COMMIT (replicas finalize once enough matching PREPAREs exist) — are each necessary specifically because a single round of voting cannot be cross-checked against a lying minority.
- Tendermint restructures the same propose/cross-check/commit idea into repeating rounds (PROPOSE, PREVOTE, PRECOMMIT) per block height, so a stalling or dishonest proposer simply causes a move to the next round with a new proposer, rather than stalling the whole chain.
- Both PBFT and Tendermint have O(n²) message complexity — every validator communicates with every other validator, twice — which is why classical BFT systems favor small, known validator sets over permissionless, internet-scale participation.
- A BFT engine does not fit neatly behind GoChain's synchronous `consensus.Engine` interface the way `ProofOfStake` did, because producing a block under BFT consensus is fundamentally a multi-round, multi-party networked conversation, not a single-node computation.
- A view change is PBFT's structured recovery mechanism for replacing a stalling or misbehaving primary while preserving anything already agreed on; Tendermint achieves the same outcome more simply by treating "move to the next round with a new proposer" as the normal path rather than a special case.
- HotStuff and other post-Tendermint BFT protocols reduce the same O(n²) message bottleneck to a linear pattern by routing votes through a rotating leader, showing that classical BFT's core safety bound (`3f+1`) is fixed even as its engineering keeps evolving.
- Real permissioned and semi-permissioned systems — Hyperledger Fabric's design philosophy, and Cosmos SDK chains running on Tendermint — are the concrete, real-world destinations this conceptual grounding prepares you to read about with genuine understanding.

---

## Exercises

### Easy

1. In your own words, explain the difference between a "fail-stop" fault and a "Byzantine" fault, using an example other than the ones given in this chapter.
2. Why must at least two-thirds of validators be honest for classical BFT consensus to work, rather than a simple majority (just over one-half), as many people initially assume? (Hint: think about what happens if exactly half are honest and half are dishonest, split evenly.)
3. Name one concrete way PBFT's "finality" differs from proof of work's "finality," and explain why a payment processor might specifically prefer the former for high-value transactions.

### Medium

4. Draw (as an ASCII diagram, in your own words, not copied from this chapter) what happens in Tendermint when a proposer at height 100, round 0, crashes mid-round and produces no proposal at all. Which step's timeout detects this, and what happens next?
5. PBFT and Tendermint both have O(n²) message complexity. Calculate, roughly, how many total messages one full PBFT round (PRE-PREPARE + PREPARE + COMMIT) requires among 4 validators, then among 100 validators, and explain in a sentence or two why this makes 100 validators a meaningfully heavier network burden than 4, even though both are "small" by permissionless-blockchain standards.
6. Research Tendermint's actual validator set size in a real, currently running Cosmos SDK chain of your choice, and compare it to the number of active miners or validators on Bitcoin or Ethereum. What does the difference in scale tell you about the permissioned-vs-permissionless trade-off from Section 2?

### Hard

7. PBFT's original design includes a "view change" protocol for replacing a primary that appears to be stalling or lying, which this chapter deliberately did not diagram in detail. Research how PBFT's view change works, and explain, using the same style of sequence diagram this chapter used for the normal case, what messages are exchanged to elect a new primary.
8. Section 8 argued that a BFT engine doesn't fit `consensus.Engine`'s synchronous, single-node-call shape. Sketch (in Go, as a rough interface design rather than a full implementation) what a differently-shaped interface for a BFT engine might look like instead — one that acknowledges block production requires multiple network round trips and a timeout-driven round/view mechanism. What would `core.Blockchain` need to change to use it?
9. Research a real-world incident where a BFT-based blockchain (a Cosmos SDK chain, or another Tendermint-based network) experienced a halt or fork due to failing to reach the required two-thirds supermajority (for example, due to validator downtime or a network partition), and explain, using this chapter's vocabulary, exactly which phase of voting failed to gather enough matching votes and why the chain correctly refused to proceed rather than committing an unsafe block.
10. Section 10 claimed HotStuff reduces message complexity from O(n²) to O(n) by routing votes through a rotating leader instead of an all-to-all broadcast. Using the same style of calculation as Exercise 5, compute the total number of messages one decision requires under HotStuff's leader-routed pattern for 4 validators and for 100 validators, and compare both numbers directly against your PBFT answers from Exercise 5.
11. Research Meta's (Facebook's) Diem/Libra project, which adopted a HotStuff-based consensus protocol before the project was discontinued. Identify one specific design reason the project's documentation gave for choosing a linear-communication BFT protocol over a Tendermint- or classical-PBFT-style design, and relate it back to this chapter's O(n²)-versus-O(n) discussion.
