# Chapter 84: Permissioned Chains — Hyperledger Fabric, Overview

Every chain this course has looked at so far — GoChain itself, Bitcoin, Ethereum — shares one deep assumption: anyone, anywhere, can download the software, generate a key pair, and join the network as a full participant, with no one's permission required. That assumption is called **permissionless**, and it is so central to how this course has taught blockchain that it is easy to mistake it for a defining property of the word "blockchain" itself. It isn't. This chapter introduces the other branch of the family tree: **permissioned chains**, where a known, vetted set of organizations decides who is allowed to run a node at all — using Hyperledger Fabric, the most widely deployed permissioned blockchain framework, as the concrete example. You will not build any Fabric code in this chapter; the goal is to understand its architecture well enough to recognize, in a real engineering conversation, exactly when a permissioned design is the right tool and when it is a needless restriction on a problem GoChain's own public design already solves better.

## Table of Contents

1. [Not Every Chain Wants to Be Public](#1-not-every-chain-wants-to-be-public)
2. [Permissioned vs. Permissionless: The Core Distinction](#2-permissioned-vs-permissionless-the-core-distinction)
3. [Hyperledger Fabric's Architecture, Conceptually](#3-hyperledger-fabrics-architecture-conceptually)
4. [Channels: Fabric's Privacy Model](#4-channels-fabrics-privacy-model)
5. [The Endorsement Transaction Flow](#5-the-endorsement-transaction-flow)
6. [Chaincode: Fabric's Smart Contracts](#6-chaincode-fabrics-smart-contracts)
7. [Ordering and Consensus in Fabric](#7-ordering-and-consensus-in-fabric)
8. [When Permissioned Is the Right Call](#8-when-permissioned-is-the-right-call)
9. [Fabric vs. GoChain: Side-by-Side](#9-fabric-vs-gochain-side-by-side)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Not Every Chain Wants to Be Public

Picture two very different kinds of shared notebook. The first is Chapter 1's original analogy: a shared ledger that literally anyone in the world can pick up, write a page into (following the rules), and read every page of, forever. That's the public, permissionless model GoChain, Bitcoin, and Ethereum all share.

The second notebook belongs to a specific group of banks that have agreed to jointly track interbank settlements, or a specific group of shipping companies, customs authorities, and warehouses that have agreed to jointly track a shared supply chain. These parties already know and have a legal relationship with each other — they are not anonymous strangers relying on cryptography and economic incentives to keep each other honest, the way Bitcoin's miners are. What they still want, though, is exactly the property Chapter 1 introduced as blockchain's core value: a shared, tamper-evident record that no single participant can quietly rewrite after the fact, and that every participant can independently verify rather than trusting one central database operator's word for it. They just don't want that record open to the entire internet, and they don't want to hand block-production rights to anonymous strangers competing via proof of work.

That combination — tamper-evidence and shared, independently-verifiable history, without public membership or public readability — is exactly the gap **permissioned blockchains** fill, and Hyperledger Fabric (an open-source project hosted by the Linux Foundation, originally contributed by IBM) is the most widely deployed framework built specifically for it.

---

## 2. Permissioned vs. Permissionless: The Core Distinction

The distinction is simpler than it sounds once you separate two questions that are easy to conflate: *who can read the ledger*, and *who can help produce new blocks on it*.

```
  PERMISSIONLESS (GoChain, Bitcoin, Ethereum)
  +--------------------------------------------------------------+
  | Anyone can join the network, run a node, and read every       |
  | block, with no application or approval process.                |
  |                                                                  |
  | Anyone can attempt to produce the next block (mine, or be       |
  | selected as a validator), subject only to the consensus          |
  | rule (proof of work, proof of stake) -- not anyone's approval.   |
  |                                                                  |
  | Identity is cheap and anonymous: a new keypair is a new,         |
  | unvetted "identity" (Chapter 51's Sybil-attack discussion is     |
  | a direct consequence of this).                                   |
  +--------------------------------------------------------------+

  PERMISSIONED (Hyperledger Fabric)
  +--------------------------------------------------------------+
  | Only organizations that have been explicitly admitted to the    |
  | network (a consortium) can run a node or submit transactions.    |
  |                                                                    |
  | Which organizations are allowed to help order/validate            |
  | transactions is an explicit, governed decision -- not an open     |
  | competition.                                                        |
  |                                                                      |
  | Every participant has a real-world, vetted identity, issued and     |
  | managed through a certificate authority -- there is no anonymous    |
  | "just generate a keypair and join" path at all.                     |
  +--------------------------------------------------------------+
```

This is not a spectrum of "more secure" versus "less secure" — it is a difference in *threat model*. GoChain, Bitcoin, and Ethereum are built to stay honest even when a majority of participants might be complete strangers, some possibly malicious, coordinated only by economic incentives and math (Chapter 19's tamper-evidence, Chapter 24's proof of work, Chapter 33's signatures — every one of these exists specifically because no participant's honesty can be assumed). Fabric is built for a different, narrower problem: participants already have a real-world legal relationship (a consortium agreement, a regulatory requirement, a business contract), and the blockchain's job is to make *that specific group's* shared record-keeping tamper-evident and auditable to each other — not to coordinate value transfer among total strangers at internet scale.

---

## 3. Hyperledger Fabric's Architecture, Conceptually

Fabric names its moving pieces differently from anything you've built in GoChain, and getting the vocabulary straight up front makes the rest of this chapter much easier to follow.

```
  A HYPERLEDGER FABRIC NETWORK, AT A GLANCE

  +----------------+   +----------------+   +----------------+
  | ORGANIZATION A  |   | ORGANIZATION B  |   | ORGANIZATION C  |
  |  (e.g. Bank 1)   |   |  (e.g. Bank 2)   |   |  (e.g. Regulator)|
  |                  |   |                  |   |                  |
  |  PEER NODES      |   |  PEER NODES      |   |  PEER NODES      |
  |  (endorse,       |   |  (endorse,       |   |  (endorse,       |
  |   commit, hold    |   |   commit, hold    |   |   commit, hold    |
  |   the ledger)      |   |   the ledger)      |   |   the ledger)      |
  |                  |   |                  |   |                  |
  |  MSP (issues       |   |  MSP (issues       |   |  MSP (issues       |
  |  identities via     |   |  identities via     |   |  identities via     |
  |  its own CA)         |   |  its own CA)         |   |  its own CA)         |
  +----------------+   +----------------+   +----------------+
             \                   |                    /
              \                  |                   /
               \                 |                  /
                +-----------------------------------+
                |         ORDERING SERVICE            |
                |   (a shared, consortium-run          |
                |    component -- sequences and         |
                |    batches transactions into           |
                |    blocks, distributes them to          |
                |    every peer -- Section 7)              |
                +-----------------------------------+
```

A **peer** is the closest thing Fabric has to GoChain's `core.Blockchain` node — it holds a copy of the ledger, runs chaincode (Section 6), and validates and commits blocks. Unlike GoChain, where every node does everything, Fabric peers can specialize: an **endorsing peer** runs a proposed transaction and signs off ("endorses") that it produced a particular result, while every peer that holds the ledger acts as a **committing peer**, applying already-ordered blocks to its own copy of the ledger.

An **MSP (Membership Service Provider)** is Fabric's identity system: each organization runs its own MSP, backed by a **certificate authority (CA)**, which issues cryptographic identities (X.509 certificates, if you want the precise term) to that organization's own peers, applications, and administrators. Every single action in Fabric — endorsing a transaction, submitting one, administering a channel — is authenticated against one of these issued identities. This is the concrete mechanism behind Section 2's claim that "every participant has a real-world, vetted identity": nobody transacts on a Fabric network without first being issued a certificate by an organization's MSP, and that issuance is a deliberate, governed administrative act — not a free, anonymous key generation the way `wallet.New()` (Chapter 36) is on GoChain.

The **ordering service** is the piece with no direct GoChain equivalent at all, and Section 7 is dedicated to it — for now, think of it as a shared, consortium-operated component whose only job is deciding the order transactions get batched into blocks, entirely separately from validating what's inside them.

Here is the same architecture one more time, but mapped explicitly against GoChain's own package names, so the vocabulary shift from "GoChain node" to "Fabric peer/orderer/MSP" doesn't feel like starting from zero:

```
  ROUGH CONCEPTUAL MAPPING (not a 1:1 equivalence -- see Section 9
  for the precise, axis-by-axis comparison)

  GOCHAIN                              HYPERLEDGER FABRIC
  --------------------------------     --------------------------------
  core.Blockchain (Ch. 18)              the ledger a peer maintains
  network.Node (Ch. 46)                 a peer (endorsing/committing)
  wallet.New() (Ch. 36) --              MSP-issued X.509 certificate --
    free, anonymous keypair               vetted, administratively issued
  consensus.Engine (Ch. 25, 77)         the ordering service (Section 7)
  gochain/vm + ContractStore            chaincode + world state
  (Ch. 61-66)                            (Section 6)
  one shared chain, every full node     one ledger PER CHANNEL, visible
  sees everything (Ch. 48's gossip)      only to that channel's members
                                          (Section 4)
```

Treat this table as a rough orientation aid, not a claim that each row is functionally identical — the point of the rest of this chapter is precisely to show *where* the mapping breaks down, and why those breaks are deliberate, threat-model-driven design choices rather than gaps.

---

## 4. Channels: Fabric's Privacy Model

Here is the feature that most sharply distinguishes Fabric from every chain covered earlier in this course. GoChain, Bitcoin, and Ethereum each maintain exactly one ledger, and every full node holds a complete copy of it — there is no way for two nodes on the same network to see different, private subsets of history. Fabric's answer to "our consortium members don't all need to see each other's every transaction" is the **channel**.

Think of a channel the way you'd think of a private meeting room inside a larger building shared by several companies: everyone in the building (the overall consortium network) knows the room exists and who's allowed in it, but only the specific people invited into that particular room can see what's discussed inside — and a different private room, down the hall, can have an entirely different guest list and an entirely separate conversation, invisible to the first room's participants.

```
  ONE FABRIC NETWORK, MULTIPLE CHANNELS

  +---------------------------------------------------------------+
  |                     CONSORTIUM NETWORK                          |
  |   (Bank 1, Bank 2, Bank 3, Regulator all belong to it)           |
  |                                                                   |
  |   CHANNEL "settlement-1-2"        CHANNEL "settlement-1-3"        |
  |   +----------------------+        +----------------------+        |
  |   | Members: Bank 1,      |        | Members: Bank 1,      |        |
  |   |          Bank 2         |        |          Bank 3         |        |
  |   | Own separate ledger      |        | Own separate ledger      |        |
  |   | Bank 3 cannot see this    |        | Bank 2 cannot see this    |        |
  |   | channel's transactions     |        | channel's transactions     |        |
  |   +----------------------+        +----------------------+        |
  +---------------------------------------------------------------+
```

Each channel maintains its own completely separate ledger and its own instance of chaincode — a transaction on one channel is invisible to organizations that are not members of that channel, even though those organizations belong to the same overall network and might even run peers connected to the same ordering service. This gives a consortium the ability to run, say, a network-wide "member directory" or shared reference data on one broadly-joined channel, while keeping bilateral settlement details between just two of the members on a separate, narrowly-scoped channel — a genuinely different privacy model from anything GoChain, Bitcoin, or Ethereum offer, all of which broadcast every transaction to every full node by design (Chapter 48's gossip protocol has no concept of "don't tell this particular peer").

---

## 5. The Endorsement Transaction Flow

GoChain's transaction flow, since Chapter 29, has followed what's usually called an **order-execute** model, and so do Bitcoin and Ethereum: a transaction is first ordered into a block (via mining, Chapter 25), and *then* every node executes it (applies its effect to the UTXO set or account state) and validates the result independently, afterward. Fabric flips the middle step, following a model usually called **execute-order-validate**, specifically because it lets Fabric parallelize execution across multiple peers *before* anything is finalized in a single, agreed order.

```
  GOCHAIN / BITCOIN / ETHEREUM:  ORDER --> EXECUTE
  1. Transaction goes into a block (mining/ordering)
  2. Every node executes it and checks the result afterward

  HYPERLEDGER FABRIC:  EXECUTE --> ORDER --> VALIDATE
  1. Transaction is proposed and EXECUTED (simulated) by a
     designated set of endorsing peers, BEFORE any ordering happens
  2. Endorsing peers sign off on the result (the "endorsement")
  3. THEN the transaction (with its endorsements attached) is
     handed to the ordering service to be sequenced into a block
  4. Every peer VALIDATES the already-ordered, already-endorsed
     transaction before committing it to its own ledger copy
```

Walking through it step by step, using the Bank 1 / Bank 2 settlement example from Section 4:

1. **Proposal.** A client application (acting on behalf of, say, Bank 1) sends a transaction proposal to a set of endorsing peers required by the channel's **endorsement policy** — a rule like "this transaction is valid only if both Bank 1's peer and Bank 2's peer endorse it," the Fabric analogue of Chapter 81 Section 4's `OP_CHECKMULTISIG`, but applied to *which organizations* must agree a transaction's result is correct, rather than to which keys must sign a spend.
2. **Simulation and endorsement.** Each required endorsing peer executes (simulates) the proposed chaincode against its own current view of the ledger, without committing anything yet, and — if it agrees with the result — signs a response containing that result. Crucially, this simulation happens *before* any global ordering exists, meaning different transactions proposed at the same moment can be endorsed in parallel, entirely independently of each other.
3. **Submission to ordering.** The client collects the required endorsements and submits the transaction, endorsements attached, to the ordering service.
4. **Ordering.** The ordering service's only job is to take a stream of already-endorsed transactions from potentially many clients and put them into a strict, agreed sequence, batched into blocks — deliberately not re-executing any chaincode logic itself (Section 7 explains why that separation exists).
5. **Distribution and validation.** The ordering service distributes the resulting block to every peer on the relevant channel. Each peer then **validates** it: does this transaction actually carry the endorsements its policy requires, and — critically — does it still make sense given everything else that got ordered before it (a read/write conflict check, since another transaction touching the same data might have been ordered first, in the gap between this transaction's simulation and its final ordering)?
6. **Commit.** Only transactions that pass validation get applied to the peer's ledger and its current world state; failed ones are still recorded on the ledger (for auditability) but flagged as invalid and never applied.

```
  ENDORSEMENT FLOW, AS A SEQUENCE

  CLIENT --propose--> ENDORSING PEERS --execute & sign--> CLIENT
    |
    +--submit (tx + endorsements)--> ORDERING SERVICE
                                            |
                                            +--sequence into a block-->
                                                      ALL PEERS
                                                          |
                                                   validate + commit
```

The single biggest structural payoff of execute-order-validate, worth naming explicitly: because simulation happens before ordering and can run on multiple endorsing peers in parallel, Fabric can execute many unrelated transactions' chaincode simultaneously, and only needs strict, single-threaded agreement for the comparatively cheaper step of *sequencing* already-computed results — a genuinely different scalability trade-off from GoChain's (and Bitcoin's and Ethereum's) approach of executing everything, by every node, strictly in the one agreed order.

---

## 6. Chaincode: Fabric's Smart Contracts

**Chaincode** is Fabric's name for smart contract code — conceptually the same idea Chapter 59 introduced (a program stored and run identically by every relevant node, so its results are trustworthy without trusting any single party), but with two concrete differences from Chapter 61-66's `gochain/vm` worth naming precisely.

First, chaincode is not bytecode for a small, bespoke, gas-metered VM the way GoChain's token contract (Chapter 65) is — it is ordinary code written in a general-purpose language (commonly Go, Java, or Node.js) and run inside its own isolated container alongside a peer, rather than interpreted by a narrow custom instruction set. This is only possible because of Section 2's threat model difference: because Fabric's participants are known, vetted organizations rather than anonymous strangers, there is less need to bound execution with something as strict as Chapter 64's per-opcode gas metering — a misbehaving chaincode container can be identified, traced back to a specific accountable organization, and addressed through the consortium's governance process, an option simply unavailable on a public, permissionless network.

Second, chaincode reads and writes a **world state** — a key-value store representing the current state of the ledger, conceptually similar to the `storage.Store`-backed state GoChain's `ContractStore` (Chapter 66) gives each contract, but shared across an entire channel's chaincode rather than partitioned strictly per contract address. The full ledger (every block, ever) is Fabric's equivalent of `core.Blockchain`'s append-only history; the world state is a fast, current-value index over it, playing exactly the role Chapter 56's UTXO index plays for GoChain's balance lookups — a derived, queryable snapshot, rebuildable from the full ledger if it were ever lost.

Here is an illustrative (not copy-pasted from any real Fabric SDK) sketch of what a tiny piece of chaincode looks like, next to the closest equivalent operation in GoChain's own token contract, so the "general-purpose language, isolated container" description in this section reads as something concrete rather than only a description:

```go
// Illustrative chaincode, representative form -- a transfer function
// for a supply-chain "asset" chaincode, written in ordinary Go and
// run inside its own container by the Fabric peer runtime, NOT
// interpreted by a bespoke opcode-based VM the way GoChain's
// contracts are.
func (c *AssetContract) Transfer(ctx TransactionContext, assetID string, newOwner string) error {
    asset, err := ctx.GetWorldState(assetID)
    if err != nil {
        return err
    }
    asset.Owner = newOwner
    return ctx.PutWorldState(assetID, asset)
}
```

```go
// The closest GoChain equivalent: Chapter 65's token contract logic,
// expressed as VM opcodes operating through Chapter 66's SLOAD/SSTORE
// against ContractStore -- illustrative, not literal bytecode.
program := []vm.Op{
    vm.OpSLoad,     // read the current owner from contract storage
    vm.OpPushData,  // push the new owner
    vm.OpSStore,    // write the new owner back to contract storage
}
```

Both accomplish the same conceptual job — read some persisted state, change it, write it back — but chaincode does so as ordinary, uncontained-by-gas-metering Go code running in its own container, while GoChain's contract does so as a small, explicitly gas-metered sequence of opcodes interpreted by `vm.VM`. That difference is not accidental; it is the direct consequence of Section 2's threat-model distinction, restated here in code rather than prose.

---

## 7. Ordering and Consensus in Fabric

Section 5 deliberately deferred one question: what, exactly, does the ordering service use to agree on a single sequence of transactions, and how does that compare to Chapter 25's proof of work or Chapter 77's proof of stake?

Because Fabric's participants are known, vetted organizations rather than anonymous strangers, the ordering service does not need a consensus mechanism designed to resist a majority of completely unknown, potentially adversarial participants (the exact problem proof of work and proof of stake are built to solve). Instead, Fabric's ordering service commonly runs a **crash fault-tolerant (CFT)** consensus protocol — one designed to keep working correctly as long as a majority of ordering nodes are up and behaving as configured, but not necessarily designed to tolerate a node that is actively, maliciously lying (a **Byzantine fault**, in the vocabulary Chapter 78 already introduced). Raft is the commonly used protocol for Fabric's ordering service today; Fabric's design also supports Byzantine-fault-tolerant ordering for deployments that need to tolerate actively malicious ordering nodes rather than merely crashed ones, since not every consortium's trust assumptions are identical.

```
  CONSENSUS TRUST ASSUMPTIONS, COMPARED

  GOCHAIN / BITCOIN / ETHEREUM (PoW or PoS)
    assumes: some fraction of participants may be actively
             malicious, anonymous, and economically motivated
             to cheat -- consensus must resist that directly

  FABRIC'S ORDERING SERVICE (commonly Raft, a CFT protocol)
    assumes: ordering nodes are run by known, accountable
             organizations that might crash or go offline, but
             are not the network's primary line of defense
             against malicious behavior -- endorsement policies
             and peer-side validation (Section 5) carry more of
             that weight instead

  FABRIC (BFT ordering option, and Ch. 78's PBFT/Tendermint)
    assumes: some ordering/validator nodes might actively lie,
             requiring explicit multi-round voting to tolerate it
```

This is a genuinely important, durable point about consensus in general, worth connecting back to Chapter 78 directly: "which consensus mechanism is correct" is not a universal ranking — it is a question that can only be answered once you know what kind of dishonesty the mechanism actually has to survive. Public, permissionless GoChain-style networks assume the worst (anonymous strangers, some fraction actively adversarial) and pay a real cost (mining energy, or staked capital at risk) for consensus that survives that assumption. A permissioned consortium, having already solved "who is allowed to participate" through real-world vetting and legal agreements, can reasonably choose a lighter-weight consensus mechanism for ordering specifically, and lean on endorsement policies and per-peer validation to catch the kinds of dishonesty vetting alone doesn't rule out.

---

## 8. When Permissioned Is the Right Call

Chapter 80 already gave you a general decision framework for consensus and architecture; here is the same kind of thinking, applied specifically to "public or permissioned."

A permissioned design like Fabric tends to be the right engineering choice when most or all of the following are true: the participants are a known, bounded, already-vetted set of organizations (not an open public); there is a real business or regulatory reason some transaction data must stay private to a subset of participants (Section 4's channels exist for exactly this); the participants already have enforceable, real-world legal recourse against each other if someone misbehaves, reducing how much the *consensus mechanism itself* needs to defend against malice; and there is no need for anyone outside the consortium to independently verify the ledger's contents.

A public, permissionless design like GoChain's — or Bitcoin's, or Ethereum's — is the right call when the opposite conditions hold: the whole point is letting arbitrary strangers transact with each other with no prior relationship and no legal recourse, verifiability by any outside observer (a regulator, a member of the public, a competing business) is itself a requirement rather than something to avoid, or there is no natural "consortium" to define membership around in the first place (a public payments network, by definition, cannot have a bounded, pre-vetted membership list without ceasing to be public).

```
  A ROUGH DECISION GUIDE

  Choose PERMISSIONED (Fabric-style) when:
    - participants are known, bounded, and already have legal
      relationships/recourse with each other
    - some data genuinely must stay private to a subset of
      participants (Section 4's channels)
    - throughput and low latency matter more than resisting
      anonymous, adversarial strangers

  Choose PERMISSIONLESS (GoChain/Bitcoin/Ethereum-style) when:
    - anyone, including total strangers with no legal recourse
      against each other, must be able to participate
    - public verifiability by anyone, not just consortium
      members, is itself a requirement
    - there is no natural, defensible way to bound membership
      without defeating the point of the system
```

Neither answer is more "advanced" than the other — they are different tools engineered for genuinely different threat models and business requirements, the same lesson Chapter 80's broader framework already established for consensus mechanisms specifically.

A short, worked example makes the guide concrete. Imagine three companies — a manufacturer, a shipper, and a customs authority — want a shared, tamper-evident record of a shipping container's custody as it moves from factory to port to destination. They already have contracts with each other; none of them wants the container's contents or price visible to the general public; and they want fast, cheap confirmation, not a multi-minute proof-of-work wait. Every condition in the "choose permissioned" list from this section is satisfied, which is exactly why real-world supply-chain tracking is one of Hyperledger Fabric's most common deployment patterns. Now imagine instead a payments network meant to let any two strangers, anywhere in the world, exchange value without a bank or intermediary vetting either of them first — there is no natural consortium to define membership around, verifiability by total strangers is the entire point, and no participant has any prior legal relationship with any other. That is exactly GoChain's (and Bitcoin's and Ethereum's) problem, not Fabric's.

---

## 9. Fabric vs. GoChain: Side-by-Side

```
  +----------------------+---------------------------+---------------------------+
  | AXIS                  | HYPERLEDGER FABRIC          | GOCHAIN (this course)      |
  +----------------------+---------------------------+---------------------------+
  | Network membership     | Permissioned, vetted         | Permissionless, anyone     |
  |                        | consortium (MSP/CA-issued    | with a keypair (Ch. 36)     |
  |                        | identities)                    |                              |
  | Ledger visibility       | Per-channel; participants     | One shared ledger, visible |
  |                        | outside a channel cannot       | to every full node          |
  |                        | see its transactions           |                              |
  | Transaction flow         | Execute-order-validate         | Order-execute (mine, then  |
  |                        | (Section 5)                     | every node validates)        |
  | Consensus (ordering)      | Crash fault-tolerant (e.g.       | Proof of work (Ch. 25) or  |
  |                        | Raft), or BFT ordering option     | proof of stake (Ch. 77)     |
  | Smart contracts            | Chaincode: general-purpose        | gochain/vm: small, gas-      |
  |                        | languages, isolated containers    | metered, bespoke bytecode   |
  | Identity model              | Certificate-issued via each        | Free key generation, no      |
  |                        | org's MSP -- no anonymous path      | vetting (Ch. 13, Ch. 36)     |
  | Native currency             | None required -- Fabric is not     | gochip (native to the       |
  |                        | inherently a cryptocurrency          | chain's own ledger)          |
  +----------------------+---------------------------+---------------------------+
```

One last, easily-missed point this table makes visible: Hyperledger Fabric does not need a native token or currency at all — a consortium tracking shared supply-chain state or interbank settlement records has no inherent need to invent its own coin, whereas GoChain, Bitcoin, and Ethereum's entire economic security model (mining incentives, gas fees) depends on one existing. That is not a limitation of permissioned chains; it's a sign that "does this system need its own currency" and "is this system permissioned or permissionless" are two genuinely independent design questions, easy to conflate if the only chains you've studied so far all happened to need one.

---

## Summary

- Permissioned chains like Hyperledger Fabric restrict network membership to a known, vetted consortium of organizations, in contrast to GoChain's, Bitcoin's, and Ethereum's permissionless model where anyone can join with no approval required.
- Fabric's architecture names its components differently: peers (endorse and commit transactions), organizations (each running its own Membership Service Provider and certificate authority for issuing identities), and a separate ordering service.
- Channels give Fabric a genuine privacy model no chain covered earlier in this course has: separate sub-networks with their own ledgers, visible only to their specific member organizations, even within one larger consortium network.
- Fabric's transaction flow is execute-order-validate — proposed transactions are simulated and endorsed by required peers *before* being ordered into a block, then validated again after ordering — a deliberate inversion of the order-execute model GoChain, Bitcoin, and Ethereum all share.
- Chaincode is Fabric's smart-contract mechanism: general-purpose code in isolated containers rather than a small, custom, gas-metered VM, made possible precisely because Fabric's participants are known and accountable rather than anonymous.
- Fabric's ordering service commonly uses a crash fault-tolerant protocol like Raft rather than proof of work or proof of stake, because known-organization participants change what a consensus mechanism actually needs to defend against — with a BFT ordering option available for consortiums that need to tolerate actively malicious ordering nodes.
- Choosing permissioned versus permissionless is a threat-model and requirements question, not a maturity ranking: permissioned fits bounded, vetted consortiums with real privacy needs and legal recourse; permissionless fits systems that must serve anonymous strangers with public verifiability as a requirement.
- A native currency and a permissioned/permissionless design are independent decisions — Fabric needs no native token at all, while GoChain's, Bitcoin's, and Ethereum's economic security depends on one.

---

## Exercises

### Easy

1. **In 100-150 words, explain the difference between "who can read the ledger" and "who can help produce new blocks,"** and state Fabric's answer to each question versus GoChain's answer to each.
2. **Draw Section 3's diagram from memory**, labeling each Fabric component (organization, peer, MSP, ordering service) with the closest GoChain concept it resembles, or writing "no direct equivalent" where none exists.
3. **In your own words, explain why Fabric's channels (Section 4) have no equivalent in GoChain's gossip-based broadcasting (Chapter 48).** What would have to change about Chapter 48's design for GoChain to support something channel-like?

### Medium

4. **Walk through Section 5's six-step endorsement flow using a concrete example of your own devising** (not the Bank 1/Bank 2 settlement example used in this chapter) — name your own consortium members, describe your own endorsement policy, and trace a single transaction from proposal through commit.
5. **Compare Fabric's execute-order-validate model to GoChain's order-execute model (Section 5) in 200-250 words**, specifically addressing what happens in each model if two transactions submitted around the same time try to spend/modify the exact same piece of data — where does the conflict get caught in each model, and by which component?
6. **Research Raft consensus** (the crash fault-tolerant protocol commonly used by Fabric's ordering service, mentioned in Section 7) at a conceptual level, and write a 200-word explanation of what "crash fault-tolerant but not Byzantine fault-tolerant" means in practice, using an example of a failure Raft is designed to survive and one it is not.

### Hard

7. **Design (on paper, no implementation required) a "permissioned mode" for GoChain**: sketch what would need to change in `network.Node` (Chapter 46) to require an admitted, certificate-like identity before a peer connection is accepted, and describe how you would represent something channel-like (Section 4) using GoChain's existing `core.Blockchain` and `network` packages without a full Fabric-style rewrite.
8. **Implement a minimal endorsement-policy checker in Go**, independent of GoChain's actual transaction pipeline: given a list of required endorser identities (as strings) and a set of signatures actually collected for a transaction, write a function that returns whether the endorsement policy (e.g., "at least 2 of {OrgA, OrgB, OrgC}" or "OrgA AND OrgB") is satisfied, and write tests covering both AND-style and threshold-style policies.
9. **Research a real, publicly documented Hyperledger Fabric deployment** (for example, a supply-chain or trade-finance consortium use case) and write a 250-350 word case study: who are the consortium members, what problem does the shared ledger solve for them that a normal shared database would not, and which of Section 8's "choose permissioned" conditions does this real deployment actually satisfy? Cite your source.
