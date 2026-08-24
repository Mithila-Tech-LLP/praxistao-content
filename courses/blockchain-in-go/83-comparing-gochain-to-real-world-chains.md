# Chapter 83: Comparing GoChain to Real-World Chains

Chapter 81 put Bitcoin next to GoChain, axis by axis. Chapter 82 did the same for Ethereum. Each of those chapters built its comparison table at the end, after walking through the details that justified every row. This chapter does the opposite: it starts from the table. Think of the last two chapters as two separate site visits to two different, real factories, each one focused entirely on the factory you were standing in. This chapter is the moment you go back to your own workshop, lay all three blueprints side by side on one table, and trace, axis by axis, exactly where your machine matches each of theirs, and exactly where — and why — it deliberately doesn't.

## Table of Contents

1. [Why a Combined Table Matters](#1-why-a-combined-table-matters)
2. [Data Model Axis](#2-data-model-axis)
3. [Consensus Axis](#3-consensus-axis)
4. [Virtual Machine / Scripting Axis](#4-virtual-machine--scripting-axis)
5. [Networking Axis](#5-networking-axis)
6. [Storage and State Commitment Axis](#6-storage-and-state-commitment-axis)
7. [Cryptography and Addressing Axis](#7-cryptography-and-addressing-axis)
8. [The Master Comparison Table](#8-the-master-comparison-table)
9. [Reading the Gaps: What "Production-Ready" Actually Adds](#9-reading-the-gaps-what-production-ready-actually-adds)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why a Combined Table Matters

A single comparison table, built once and read carefully, teaches something neither Chapter 81 nor Chapter 82 could teach on its own: *which* of GoChain's design decisions are "the UTXO/Bitcoin way," *which* are "the account/Ethereum way," and *which* are neither — genuine, independent simplifications made purely for teachability, that neither production chain shares. Without laying all three next to each other, it is easy to mistake "GoChain simplified this" for "GoChain copied Bitcoin" or "GoChain copied Ethereum" when the real answer, on several axes, is "GoChain built its own smaller version of an idea both real chains also had to solve."

This chapter is organized by *axis* — data model, consensus, VM, networking, storage, cryptography — rather than by chain, precisely so you can read straight down a single row and see all three systems' answers to the exact same design question at once. Each section below builds the reasoning for one row (or a small cluster of rows) of Section 8's master table; if you want the full table first and the reasoning after, skip ahead and come back.

---

## 2. Data Model Axis

The plain-language version, worth restating once more because it is the single most consequential fork in this entire comparison: a **UTXO model** tracks discrete, spendable chunks of value, like bills and coins in a wallet; an **account model** tracks one running balance per address, like a line in a bank's ledger.

```
  DATA MODEL
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN                     | ETHEREUM                    |
  +----------------------+---------------------------+---------------------------+
  | Model                 | UTXO                        | Account                     |
  | "Balance" is           | sum of owned, unspent        | one stored number per       |
  |                        | outputs                      | address                     |
  | Replay protection      | not needed (each output       | per-account nonce            |
  |                        | spends exactly once)          |                              |
  | Parallelism            | naturally high (unrelated    | naturally lower (same-      |
  |                        | spends never conflict)        | account txs must order by   |
  |                        |                                | nonce)                       |
  +----------------------+---------------------------+---------------------------+

                                | GOCHAIN
                                +---------------------------+
                                | UTXO (matches Bitcoin,      |
                                | Ch. 30, core.Transaction)    |
                                | sum of owned UTXOs           |
                                | not needed, by construction  |
                                | naturally high                |
                                +---------------------------+
```

GoChain's choice here is not a simplification of either real chain — it is a straightforward adoption of Bitcoin's model, for the reasons Chapter 31 already gave in full: UTXO is the more teachable base ledger, and Chapter 65-66's contract storage gives contracts an account-*like* persistence story without requiring the whole ledger to be rebuilt around it.

It is worth being precise about exactly where GoChain's design stops matching Bitcoin's on this axis, though, since Chapter 82 already surfaced the answer: Chapter 65-66's contract storage. A plain GoChain transfer (Chapter 32-34) is pure UTXO, indistinguishable in shape from Bitcoin's. But the moment a GoChain contract needs its own persistent state — Chapter 65's token contract's balance table — Chapter 66 reaches for exactly the account-model idea Section 2 just described: a key-value store, keyed by contract address plus storage slot, that gets read and mutated in place, not consumed-and-replaced the way a UTXO is. This is precisely why Chapter 31 previewed "our smart-contract volume will introduce account-like state for contracts specifically" rather than promising a purely UTXO experience end to end — the two models solve different problems, and a well-designed system can use each one exactly where it fits, rather than forcing one model to do both jobs badly.

---

## 3. Consensus Axis

Consensus answers one question — "how does a network of independent, mutually-distrusting computers agree on what happened, and in what order?" — and each chain answers it slightly differently.

```
  CONSENSUS
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN                     | ETHEREUM                    |
  +----------------------+---------------------------+---------------------------+
  | Mechanism              | Proof of work                | Proof of stake (since The  |
  |                        |                                | Merge, Ch. 82 Sec. 7;      |
  |                        |                                | formerly proof of work)     |
  | Sybil resistance via    | spent computation (hashing)  | staked capital (at risk of |
  |                        |                                | slashing)                   |
  | Block time (target)    | ~10 minutes                   | ~12 seconds                 |
  | Retarget/adjustment     | every 2016 blocks (Ch. 81)    | dynamic, protocol-level    |
  +----------------------+---------------------------+---------------------------+

                                | GOCHAIN
                                +---------------------------+
                                | Proof of work (Ch. 25),      |
                                | proof of stake also           |
                                | available (Ch. 77), same      |
                                | consensus.Engine interface    |
                                | spent computation OR staked   |
                                | value, learner's choice        |
                                | short, teaching-scale window  |
                                | shorter, teaching window       |
                                | (Ch. 26)                       |
                                +---------------------------+
```

This is the one axis where GoChain doesn't pick a side at all — it implements *both* mechanisms, behind one shared `consensus.Engine`-style interface (Chapter 77), specifically so you can run a small testnet under either and compare them directly. That design decision is itself modeled on the architectural boundary Chapter 82 Section 7 showed you Ethereum actually relied on during The Merge: a consensus mechanism that can be swapped without touching the rest of the system.

It's worth naming what stays constant on either side of that swap, because it's the same list of things The Merge left untouched (Chapter 82 Section 7): block validation rules for transaction contents (Chapter 19), the UTXO set and its double-spend checks (Chapter 34), the VM and every deployed contract's storage (Volume 9) — none of it cares whether the block it's validating was produced by a winning proof-of-work nonce or a selected proof-of-stake validator. Only the *rule for choosing who gets to propose the next block, and what they risk by lying about it* changes. That narrow, well-defined seam is exactly what makes "swap the consensus engine underneath a running system" a tractable engineering problem instead of a full rewrite, at both GoChain's teaching scale and Ethereum's production scale.

---

## 4. Virtual Machine / Scripting Axis

Here the two real chains diverge sharply from each other, and GoChain sits in between them by design.

```
  VIRTUAL MACHINE / SCRIPTING
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN                     | ETHEREUM                    |
  +----------------------+---------------------------+---------------------------+
  | Language               | Bitcoin Script                | EVM bytecode                |
  | Loops allowed?          | No (Ch. 81 Sec. 4)            | Yes (JUMP/JUMPI)            |
  | Gas metering?           | No -- bounded by construction | Yes, mandatory (Ch. 82 Sec.5)|
  | Persistent contract     | No -- scripts only lock/       | Yes -- storage trie per     |
  | storage?                | unlock outputs                | contract account            |
  | Memory region?          | No                             | Yes (MLOAD/MSTORE)          |
  | Event/log mechanism?    | No                             | Yes (LOG0-4)                |
  | Inter-contract calls?   | No                             | Yes (CALL family)           |
  +----------------------+---------------------------+---------------------------+

                                | GOCHAIN
                                +---------------------------+
                                | gochain/vm                   |
                                | Yes (OpJump/OpJumpIfFalse,    |
                                | Ch. 61)                        |
                                | Yes, mandatory (Ch. 64)        |
                                | Yes (ContractStore, Ch. 66)    |
                                |                                 |
                                | No                              |
                                | No                              |
                                | No (sketched only as an        |
                                | exercise, Ch. 67 Sec. 8/9)      |
                                +---------------------------+
```

GoChain's VM is architecturally closer to Ethereum's philosophy (loops, gas, persistent storage) than to Bitcoin Script's (bounded, loop-free, no persistence beyond the UTXO it's locking) — Chapter 81 Section 4 and Chapter 82 Section 9 both said this in passing; seeing it as one table row makes the shape of that choice unambiguous. GoChain then stops well short of the EVM's full scope: no memory region, no logging, no inter-contract calls — each a deliberate line drawn to keep `vm.VM`'s opcode table small enough for one learner to implement completely.

There's a second, quieter consequence of this table worth drawing out: the *cost* of expressiveness scales with exactly what a VM allows. Bitcoin Script's row is short and simple because it deliberately can't do much — there's very little for a verifier to reason about, and correspondingly very little attack surface (Chapter 81 Section 4 already made this trade-off explicit). Ethereum's row is long specifically because the EVM can do almost anything a general-purpose language can, and Chapter 67's entire reentrancy lesson exists *because* one specific capability on that longer list — inter-contract calls — makes an entirely new bug class possible that a loop-free, call-free scripting language like Bitcoin Script structurally cannot have. GoChain's middle position is not indecision; it's a specific, bounded slice of that expressiveness (enough to teach gas, storage, and a real token contract) chosen precisely because it lets Chapter 67 teach reentrancy's real lesson without yet needing the full CALL family that would make the vulnerability even easier to construct by accident.

---

## 5. Networking Axis

Every peer-to-peer chain has to solve the same three problems — find peers, propagate data efficiently, and agree on which competing chain wins — and Volume 7 built GoChain's own working answer to each.

```
  NETWORKING
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN                     | ETHEREUM                    |
  +----------------------+---------------------------+---------------------------+
  | Peer discovery         | seed nodes + address        | seed nodes + a discovery   |
  |                        | exchange (Ch. 47)            | protocol (Kademlia-based)  |
  | Propagation             | gossip (Ch. 48)              | gossip, over a publish/     |
  |                        |                                | subscribe network layer    |
  | Fork resolution          | most-accumulated-work chain  | fork-choice rule built for |
  |                        | wins (Ch. 50)                 | proof-of-stake finality    |
  | Node types               | full / pruned / SPV (Ch. 81)  | full / light client;        |
  |                        |                                | separate execution and      |
  |                        |                                | consensus-layer clients     |
  |                        |                                | since The Merge              |
  +----------------------+---------------------------+---------------------------+

                                | GOCHAIN
                                +---------------------------+
                                | seed nodes + address        |
                                | exchange (Ch. 47)            |
                                | gossip broadcast (Ch. 48)     |
                                | longest-accumulated-work      |
                                | chain wins (Ch. 50)            |
                                | full node only (Ch. 81 Sec.7)  |
                                +---------------------------+
```

This is the axis where GoChain's design is closest to Bitcoin's *specifically*, because Chapter 45 said so outright: GoChain's wire protocol was "closely modeled on Bitcoin's own original protocol." Ethereum's real networking stack has since grown its own separate peer-discovery protocol and, after The Merge, a genuinely more complex split between execution-layer and consensus-layer client software communicating with each other locally — complexity that exists because Ethereum's post-Merge architecture actually runs two cooperating pieces of software per node, a wrinkle neither Bitcoin nor GoChain has any reason to introduce.

That two-client wrinkle is worth a moment of extra attention, because it is a direct, structural consequence of Chapter 82 Section 7's Merge case study: once consensus (the Beacon Chain's proof-of-stake validators) and execution (the EVM, accounts, transactions) became genuinely separate concerns handled by genuinely separate software components communicating over a local API, a single Ethereum "node" stopped being one program and became two cooperating ones. GoChain's `network.Node` (Chapter 46) and `consensus.Engine` (Chapter 77) remain a single process specifically because GoChain never needed to de-risk a live consensus migration the way Ethereum did — the swappable-interface boundary is there in GoChain's code, but nothing forces the two sides of that boundary to run as separate operating-system processes the way Ethereum's post-Merge design does.

---

## 6. Storage and State Commitment Axis

```
  STORAGE AND STATE COMMITMENT
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN                     | ETHEREUM                    |
  +----------------------+---------------------------+---------------------------+
  | Persistent storage      | full node's own database    | full node's own database    |
  |                        | (implementation-specific)     | (implementation-specific)   |
  | Fast balance lookups     | UTXO set kept indexed        | state trie IS the account/  |
  |                        |                                | balance index                |
  | State commitment         | no single "state root" --   | state trie root, per block   |
  |                        | UTXO set isn't hashed into    | header (Ch. 82 Sec. 3)       |
  |                        | the block header itself        |                              |
  | Per-transaction proof    | Merkle root over txs         | transactionsRoot AND         |
  |                        | (Ch. 10)                       | receiptsRoot (Ch. 82 Sec. 3) |
  +----------------------+---------------------------+---------------------------+

                                | GOCHAIN
                                +---------------------------+
                                | BoltDB (Ch. 54-55)            |
                                | dedicated UTXO index           |
                                | (Ch. 56)                        |
                                | simplified trie root over the  |
                                | UTXO set (Ch. 57)               |
                                | Merkle root over txs (Ch. 10)   |
                                | only -- no receipts trie         |
                                +---------------------------+
```

Worth pausing on one subtlety this row surfaces cleanly: Bitcoin's block header commits to its transaction list (via the Merkle root) but not to the UTXO set itself as a single root hash — a node computes and indexes its own UTXO set locally, but no single number in the header lets two nodes instantly confirm they agree on the *entire current state* the way Ethereum's state root does. GoChain's Chapter 57 trie actually goes a step further than real Bitcoin here, adding a state-root-like commitment over the UTXO set that real Bitcoin's header format never included — a small but genuine example of GoChain's design borrowing an Ethereum-style idea and applying it inside an otherwise Bitcoin-shaped ledger.

---

## 7. Cryptography and Addressing Axis

```
  CRYPTOGRAPHY AND ADDRESSING
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN                     | ETHEREUM                    |
  +----------------------+---------------------------+---------------------------+
  | Signature scheme        | ECDSA over secp256k1         | ECDSA over secp256k1        |
  | Hash function(s)        | SHA-256 (and RIPEMD-160       | Keccak-256                   |
  |                        | for address hashing)           |                              |
  | Address encoding         | Base58Check (Ch. 14)          | hex-encoded, with a          |
  |                        |                                | checksum-by-capitalization   |
  |                        |                                | convention                    |
  | HD wallets               | BIP-32/39/44 (Ch. 38-40)       | BIP-32/39/44 (same standards)|
  +----------------------+---------------------------+---------------------------+

                                | GOCHAIN
                                +---------------------------+
                                | ECDSA via Go's crypto/ecdsa   |
                                | (standard curves, Ch. 13)      |
                                | SHA-256 (Ch. 8-9)               |
                                | Base58 with checksum (Ch. 14)   |
                                | BIP-32/39/44 (Ch. 38-40)         |
                                +---------------------------+
```

This is a pleasant surprise worth calling out: on HD wallet standards specifically, GoChain doesn't just resemble Bitcoin and Ethereum's approach — it implements the *same* standards (BIP-32, BIP-39, BIP-44) both real ecosystems actually use, because Chapter 38 built GoChain's wallet hierarchy directly from those specifications rather than a simplified stand-in. Signature scheme and hashing are the more typical story: GoChain uses Go's standard-library ECDSA (any standard curve) and SHA-256 rather than committing to Bitcoin's specific secp256k1 library dependency or Ethereum's Keccak-256 variant, exactly as Chapter 81 Section 10 already explained.

It's worth being explicit about why this particular axis produced GoChain's *closest* match to both real chains, when almost every other axis in this chapter shows at least a partial simplification. Cryptographic primitives and address-derivation standards are, by nature, more portable and more standardized across the whole industry than a chain's data model or VM design — an ECDSA signature is an ECDSA signature regardless of which curve or hash function backs it, and BIP-32/39/44 were explicitly written as cross-chain standards, not Bitcoin-only or Ethereum-only ones, precisely so that wallet software (and courses like this one) could implement them once and have them work correctly against real production expectations. Data models and virtual machines, by contrast, are exactly the axes where a chain's *specific* design philosophy shows up — which is exactly why those are the axes where GoChain's deliberate teaching simplifications are most visible.

---

## 8. The Master Comparison Table

```
  +------------------------+------------------------+------------------------+------------------------+
  | AXIS                    | GOCHAIN (this course)   | BITCOIN (real)          | ETHEREUM (real)         |
  +------------------------+------------------------+------------------------+------------------------+
  | Data model               | UTXO                     | UTXO                     | Account                  |
  | Smallest unit             | gochip                    | satoshi                  | wei                       |
  | Consensus                 | PoW (default) + PoS       | Proof of work             | Proof of stake (post-     |
  |                           | option, same interface     |                          | Merge; formerly PoW)      |
  | Block time (target)       | short, teaching window     | ~10 minutes              | ~12 seconds               |
  | Difficulty/validator       | shorter retarget window,   | 2016-block retarget       | dynamic, protocol-level   |
  | selection                 | or staking (Ch. 77)         | window (Ch. 81)          | validator duties          |
  | Locking mechanism          | LockingScript / unlock      | scriptPubKey / scriptSig | EVM contract code +       |
  |                           | (Ch. 63)                    |                          | account-level auth        |
  | VM loops allowed?          | Yes                          | No                        | Yes                       |
  | VM gas metering            | Yes, fixed costs (Ch. 64)   | N/A (bounded by design)  | Yes, dynamic market       |
  |                           |                               |                          | (Ch. 82 Sec. 5)           |
  | Persistent contract        | Yes (ContractStore,          | No                        | Yes (per-account          |
  | storage                    | Ch. 66)                       |                          | storage trie)             |
  | Event/logging mechanism    | No                            | No                        | Yes (LOG0-4)              |
  | Inter-contract calls       | No (sketched as exercise)     | No                        | Yes (CALL family)         |
  | Networking model            | Gossip P2P (Ch. 48)           | Gossip P2P                | Gossip P2P (separate       |
  |                            |                                |                          | discovery protocol)        |
  | Fork resolution              | Most accumulated work         | Most accumulated work    | Fork-choice for PoS         |
  |                            | (Ch. 50)                       |                          | finality                    |
  | Node types                  | Full only                       | Full / pruned / SPV      | Full / light client         |
  | Storage engine               | BoltDB (Ch. 54-55)              | implementation-specific  | implementation-specific    |
  | State commitment              | Simplified trie over UTXO set   | none in header (Ch. 83   | Full state trie in header  |
  |                              | (Ch. 57)                          | Sec. 6)                  | (Ch. 82 Sec. 3)             |
  | Signature scheme               | ECDSA, standard curves           | ECDSA / secp256k1        | ECDSA / secp256k1           |
  | Hash function(s)                | SHA-256                          | SHA-256 + RIPEMD-160     | Keccak-256                  |
  | Address encoding                 | Base58 + checksum (Ch. 14)      | Base58Check               | hex + capitalization        |
  |                                  |                                  |                          | checksum                    |
  | HD wallet standard                | BIP-32/39/44 (Ch. 38-40)          | BIP-32/39/44             | BIP-32/39/44                |
  +------------------------+------------------------+------------------------+------------------------+
```

Read down any single column and you get that chain's philosophy in one glance. Read across any single row and you get every chain's answer to one specific engineering question, side by side. Both readings are useful, and neither one alone tells the whole story of *why* GoChain's answers look the way they do — which is exactly what Section 9 addresses directly.

It's worth grouping this table's rows one more way before moving on, because it makes "which real chain did GoChain borrow from, on this axis" impossible to miss:

```
  WHICH ROWS MATCH WHICH REAL CHAIN

  MATCHES BITCOIN                        MATCHES ETHEREUM'S PHILOSOPHY
  --------------------------------       --------------------------------
  Data model (UTXO)                       VM loops allowed
  Locking mechanism (script-style)         VM gas metering (mandatory)
  Networking model (gossip, explicitly     Persistent contract storage
    modeled per Ch. 45)                     (Ch. 66's ContractStore)
  Fork resolution (most accumulated work)

  MATCHES NEITHER -- GOCHAIN'S OWN CHOICE
  --------------------------------------------------------
  Consensus: BOTH PoW and PoS behind one swappable interface (Ch.77)
  State commitment: a trie root in the header, layered onto an
    otherwise Bitcoin-shaped UTXO ledger (Section 6 of this chapter)
  Node types: full-node-only, simpler than either real chain's
    full/pruned/SPV or full/light-client spectrum
  Event/logging mechanism, inter-contract calls: present in neither
    real chain's simplest form and absent from GoChain entirely
```

Seeing the rows grouped this way is the clearest evidence in this whole chapter that GoChain's architecture was assembled with intent, axis by axis, rather than being a diluted copy of either real system.

---

## 9. Reading the Gaps: What "Production-Ready" Actually Adds

Every empty-looking gap in GoChain's column above is not a missing feature by accident — it's a specific category of real-world hardening that a production, adversarial, globally-distributed network needs and a single-learner teaching chain does not. Grouping the gaps this way is more useful than listing them chapter by chapter:

**Scale-driven hardening.** A dedicated UTXO index (Chapter 56) matters more the larger the UTXO set gets; a fast state trie matters more the larger the account set gets. Real chains' storage and indexing engineering exists because their state has grown to a size a simple approach genuinely cannot serve fast enough — not because GoChain's simpler approach is wrong at the scale GoChain actually runs at.

**Adversary-driven hardening.** Mempool policy layers (Chapter 81 Section 8's minimum relay fee, standardness rules, replace-by-fee), a live fee market (Chapter 82 Section 5), and the entire node-type spectrum (full/pruned/SPV) all exist because real networks have to assume some participants are actively trying to exploit the system for profit, at scale, continuously — a threat model a course exercise can *demonstrate* (Chapter 76's attack lab) without needing to defend against in earnest the way live software must.

**Expressiveness-driven complexity.** The EVM's memory region, logging opcodes, and inter-contract call family (Chapter 82 Section 4) exist because Ethereum committed to letting contracts do essentially anything, including calling each other — a decision that multiplies both what's possible and what can go wrong (Chapter 67's entire reentrancy lesson exists *because* inter-contract calls are possible at all). GoChain's narrower VM sidesteps an entire category of bug by simply not yet offering the capability that enables it.

**Operational maturity.** Real difficulty/fee/validator parameters were tuned, sometimes painfully, against years of live network behavior no course can simulate in a few weeks. GoChain's shorter retarget windows and simpler fee model aren't wrong answers to the same question — they're the right answer to a different, more tractable question: "can one learner watch this mechanism work, end to end, on their own laptop, in an afternoon?"

None of this makes the comparison a scorecard where real chains "win." It makes it a map: everywhere GoChain's column looks thinner, that's a pointer to a specific, nameable kind of engineering a production system had to do — most of which you are now equipped to actually understand, because you've built the smaller version of the exact mechanism it's built on top of.

As one last, concrete illustration of everything this chapter has been building toward, here is a single real-world action — "send 5 units of value to someone else" — traced through all three systems side by side, using nothing but vocabulary this course has already defined:

```
  ONE ACTION, THREE SYSTEMS: "SEND 5 UNITS OF VALUE"

  GOCHAIN (and Bitcoin, matching on this axis)
  +------------------------------------------------------------+
  | 1. Wallet selects UTXO(s) covering >= 5 gochips (Ch. 32)      |
  | 2. Builds core.Transaction: input(s), a 5-gochip output        |
  |    locked to the recipient, a change output locked back        |
  |    to the sender (Ch. 32)                                       |
  | 3. Signs every input with the sender's private key (Ch. 33)      |
  | 4. Broadcast via gossip to peers (Ch. 48)                          |
  | 5. Enters mempool; rejected if it double-spends (Ch. 34)            |
  | 6. Mined into a block once proof of work is solved (Ch. 25)          |
  | 7. Every node validates and updates its UTXO index (Ch. 56)           |
  +------------------------------------------------------------+

  ETHEREUM
  +------------------------------------------------------------+
  | 1. Sender signs a transaction: { to, value: 5, nonce } (Ch.82) |
  | 2. Sender picks a gas limit and a priority fee/tip (Ch. 82 S.5)  |
  | 3. Broadcast via gossip to peers                                  |
  | 4. Enters mempool; rejected if nonce is wrong or reused             |
  | 5. Ordered into a block by a selected validator (Ch. 82 Sec. 7)      |
  | 6. Every node applies it: sender.balance -= 5, recipient += 5,        |
  |    sender.nonce += 1, updating the state trie (Ch. 82 Sec. 3)          |
  | 7. New state root committed into the block header                      |
  +------------------------------------------------------------+
```

Reading the two side by side one final time makes the whole chapter's argument concrete in a single glance: seven structurally similar steps, in a similar overall order (sign, broadcast, order, apply, commit) — but every single step's *details* trace directly back to the data-model choice Section 2 opened with. That one decision — UTXOs versus accounts — is not one design detail among many; it is the fork in the road every other axis in this chapter's table ends up reflecting in some way.

---

## Summary

- A single, axis-organized comparison table across GoChain, Bitcoin, and Ethereum makes visible which of GoChain's choices came from Bitcoin (UTXO, gossip networking, longest-chain rule), which came from Ethereum (gas-metered, loop-capable VM with persistent contract storage), and which are GoChain's own independent simplifications (fixed gas costs, no fee market, full-node-only).
- On the data model axis, GoChain adopted Bitcoin's UTXO model outright rather than simplifying Ethereum's account model.
- On the consensus axis, GoChain uniquely implements both proof of work and proof of stake behind one swappable interface, mirroring the architectural boundary Ethereum's real Merge relied on.
- On the VM axis, GoChain's design sits closer to the EVM's philosophy (loops, mandatory gas, persistent storage) than to Bitcoin Script's, while deliberately omitting memory regions, event logging, and inter-contract calls.
- On networking, GoChain is closest to Bitcoin by explicit design (Chapter 45), while real Ethereum's post-Merge networking has grown a genuinely more complex, two-client-per-node architecture.
- On storage and state commitment, GoChain's Chapter 57 trie actually adds an Ethereum-style state-root commitment that real Bitcoin's header format never included — proof the comparison isn't a strict subset relationship in either direction.
- The gaps between GoChain and production chains cluster into four recognizable categories: scale-driven hardening, adversary-driven hardening, expressiveness-driven complexity, and operational maturity earned over years of live operation.

---

## Exercises

### Easy

1. **Pick any three rows from the master table in Section 8** and, for each, write one sentence stating whether GoChain's choice on that axis matches Bitcoin, matches Ethereum, or is a genuinely independent choice shared by neither.
2. **Identify the one row in Section 6** where GoChain's design actually goes further than real Bitcoin's does, and explain in 2-3 sentences what capability that extra piece of design gives GoChain that plain Bitcoin headers don't have.
3. **Using Section 9's four gap categories** (scale-driven, adversary-driven, expressiveness-driven, operational maturity), classify GoChain's lack of a live fee market (Section 5's networking-adjacent gas discussion, carried over from Chapter 82) into the category you think fits best, and justify your choice in 2-3 sentences.

### Medium

4. **Write a 200-300 word memo, as if to a new team member joining a real blockchain engineering team**, explaining which three architectural decisions from GoChain's design you would recommend keeping as-is even in a production system, and which three you would recommend hardening first, using this chapter's tables as your evidence.
5. **Take the "Cryptography and Addressing" row for hash functions** (Section 7): Bitcoin uses SHA-256 plus RIPEMD-160, Ethereum uses Keccak-256, GoChain uses SHA-256 only. Research the practical difference between Keccak-256 and the standardized SHA-3 (a common point of confusion, since Ethereum predates SHA-3's finalization), and summarize it in 150 words, citing your source.
6. **Design a fourth column for this chapter's master table** representing Hyperledger Fabric (previewed in Chapter 84), filling in your best guess for at least six axes based on what you already know about permissioned chains from Chapter 78's BFT discussion, before reading Chapter 84 — then revisit your answers after finishing Chapter 84 and note what you got right or wrong.

### Hard

7. **Implement a small Go program that walks GoChain's actual `core`, `consensus`, `vm`, `network`, and `storage` packages** and generates a machine-readable (JSON or CSV) version of a comparison row set — for example, listing every opcode currently implemented in `gochain/vm`, to make Section 4's category comparison automatically re-verifiable against the real codebase rather than hand-maintained prose.
8. **Pick one axis from this chapter (data model, consensus, VM, networking, or storage) and prototype the "next increment of production-hardening" for GoChain on that axis** — for example, a real fee market with a self-adjusting base fee (Chapter 82 Section 5) layered on top of Chapter 35's mempool, or a pruned-node mode (Chapter 81 Section 9's exercise) for the storage axis. Implement it, and write a short design note explaining which "gap category" from Section 9 it addresses.
9. **Research one real production blockchain not covered in this course** (for example, Solana, Cosmos, or Cardano) and produce your own single-column addition to the master table in Section 8, citing sources for at least four axes. Note any axis where your research reveals a genuinely different third answer this chapter's Bitcoin/Ethereum comparison didn't anticipate.
