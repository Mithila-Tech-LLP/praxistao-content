# Chapter 81: Bitcoin Architecture, Deep Dive

Every volume of this course has quietly been reading Bitcoin's design over your shoulder: your `core.Block`, your UTXO model, your proof-of-work miner, your stack-based VM are all deliberately simplified descendants of ideas Bitcoin introduced. Now that you have built and tested working versions of all of them yourself, this chapter turns around and looks directly at the real thing — Bitcoin's actual block format, its real scripting language, its real difficulty adjustment, its handling of transaction malleability, and its real mempool and node-type landscape — pointing out exactly where GoChain matches it, and exactly where GoChain simplified.

## Table of Contents

1. [Why Look at Bitcoin Now](#1-why-look-at-bitcoin-now)
2. [Bitcoin's UTXO Model, For Real](#2-bitcoins-utxo-model-for-real)
3. [The Real Bitcoin Block and Transaction Format](#3-the-real-bitcoin-block-and-transaction-format)
4. [Bitcoin Script: The Real Ancestor of GoChain's VM](#4-bitcoin-script-the-real-ancestor-of-gochains-vm)
5. [Transaction Malleability and the SegWit Fix](#5-transaction-malleability-and-the-segwit-fix)
6. [Real Difficulty Adjustment: The 2016-Block Retarget](#6-real-difficulty-adjustment-the-2016-block-retarget)
7. [Node Types: Full, Pruned, and Lightweight](#7-node-types-full-pruned-and-lightweight)
8. [Mempool Policy on Real Bitcoin Nodes](#8-mempool-policy-on-real-bitcoin-nodes)
9. [Bitcoin vs. GoChain: Side-by-Side](#9-bitcoin-vs-gochain-side-by-side)
10. [What GoChain Simplified, and Why That Was the Right Call](#10-what-gochain-simplified-and-why-that-was-the-right-call)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Look at Bitcoin Now

Think of the last twenty-three volumes as an apprenticeship. A furniture apprentice does not start by copying a master craftsman's finished cabinet joint by joint — they learn to cut a joint themselves, on scrap wood, until their hands understand *why* the joint is shaped the way it is. Only then does looking at the master's cabinet become genuinely useful, because now every design choice reads as a choice, not a mystery.

That is exactly the position you are in. You have built a UTXO ledger (Volume 5), a proof-of-work miner with difficulty adjustment (Volume 4), a stack-based virtual machine with gas (Volume 9), a gossip-based P2P network (Volume 7), and BoltDB-backed storage (Volume 8). Bitcoin is the system that inspired every one of those designs, running continuously, in production, since 2009. Reading its real architecture now is not learning new concepts — it is discovering how far a codebase can push concepts you already own.

One important framing before we go further: this chapter describes Bitcoin's durable, long-standing architectural properties — the shape of its data structures, the logic of its algorithms — rather than specific current version numbers, exact current opcode counts, or precise current default configuration values, all of which change over time as the software evolves. The architecture underneath those details has been remarkably stable for over a decade, and that stability is itself worth noticing: it is a direct consequence of the same principle you learned in Chapter 19 — a network of independent validators makes sweeping changes to core rules extremely hard to coordinate, on purpose.

---

## 2. Bitcoin's UTXO Model, For Real

You already know the mental model from Volume 5: instead of tracking a running balance per account, Bitcoin tracks a giant set of **Unspent Transaction Outputs (UTXOs)** — individual, indivisible chunks of value, like the specific bills and coins sitting in a physical wallet. Spending consumes some UTXOs completely as inputs and creates new UTXOs as outputs, including a change output back to the spender when the input value exceeds what's being sent.

GoChain's `core.Transaction`, with its `TxInput` and `TxOutput` slices, is not "inspired by" this model in some loose sense — it is structurally the same model, field for field, with the names Volume 5 chose for it. The one real difference worth naming precisely is what actually *locks* each output.

In GoChain, an output is conceptually "owned by an address," and Volume 9 taught you the real mechanism underneath that phrase: every output carries a small locking script, and every input carries an unlocking script, and your VM evaluates unlocking-script-then-locking-script together to decide if the spend is authorized. Bitcoin works identically — because GoChain's design in Chapter 63 was deliberately modeled on it. Every Bitcoin output has a **scriptPubKey** (the locking script — GoChain's "lock"), and every input that spends it supplies a **scriptSig** (the unlocking script — GoChain's "unlock"). The most common scriptPubKey template, **Pay-to-Public-Key-Hash (P2PKH)**, is functionally identical to the standard locking script you built in Chapter 63.

```
  BITCOIN OUTPUT                         GOCHAIN OUTPUT (core.TxOutput)
  --------------------------------       --------------------------------
  Value: 50,000 satoshis                 Value: 50,000 gochips
  scriptPubKey:                          LockingScript:
    OP_DUP OP_HASH160                      OP_DUP OP_HASH160
    <pubKeyHash> OP_EQUALVERIFY             <pubKeyHash> OP_EQUALVERIFY
    OP_CHECKSIG                             OP_CHECKSIG

  Spending input supplies a scriptSig:   Spending input supplies unlock:
    <signature> <publicKey>                <signature> <publicKey>
```

The only genuinely new vocabulary here is Bitcoin's unit: a **satoshi** is the smallest indivisible unit of bitcoin (1 BTC = 100,000,000 satoshis), the exact same role GoChain's smallest indivisible "gochip" unit plays for gochips. Everything else — inputs referencing prior outputs by transaction ID and output index, the UTXO set as the source of truth for balances, double-spend prevention by checking whether a referenced output is already spent — is the same mechanism you implemented and tested in Chapters 29 through 34.

One real-world wrinkle worth naming: the **UTXO set** itself — the complete collection of every currently-unspent output across the entire chain's history — has to be kept in fast, indexed storage for a node to check spends efficiently, exactly the problem your Chapter 56 UTXO index solves for GoChain. On a chain that has been running for over a decade with a large, continuously-used monetary base, this set is a substantial, ever-shifting piece of state, which is precisely why Chapter 56's benchmark (scan vs. indexed lookup) is not a toy exercise — it demonstrates, at small scale, the exact engineering problem real Bitcoin (and every UTXO-model chain) has to solve at large scale.

A second real detail: newly created coins from mining (the **coinbase output** — Chapter 37's coinbase transaction) cannot be spent immediately. Bitcoin enforces a **coinbase maturity** rule requiring a fixed number of confirmations (further blocks mined on top) before a coinbase output becomes spendable, specifically to reduce the damage a chain reorganization (Chapter 50) could do — if a block gets orphaned by a longer competing chain, any coinbase reward it minted simply never existed on the winning chain, and maturity rules keep that risk from spreading into transactions that spent those still-immature rewards.

---

## 3. The Real Bitcoin Block and Transaction Format

Your `core.Block` groups a header's worth of metadata with a list of transactions. Bitcoin's block splits these two concerns even more sharply: a compact, fixed-size **block header** and a separate, variable-length list of transactions.

```
                     A REAL BITCOIN BLOCK
  +----------------------------------------------------------+
  |                     BLOCK HEADER (fixed-size)             |
  |  version           -- which consensus rules this block    |
  |                        follows                             |
  |  previous block hash -- links to the prior block           |
  |                        (exactly like PrevBlockHash)         |
  |  merkle root       -- one hash committing to every          |
  |                        transaction below (Chapter 10)       |
  |  timestamp         -- roughly when the block was created    |
  |  bits              -- the current difficulty target,        |
  |                        compactly encoded                    |
  |  nonce             -- the number miners search over          |
  |                        (Chapter 24-25)                       |
  +----------------------------------------------------------+
  |                  TRANSACTION LIST (variable-size)          |
  |  tx count                                                   |
  |  tx[0]  (the coinbase transaction — Chapter 37)              |
  |  tx[1]                                                       |
  |  tx[2]                                                       |
  |  ...                                                         |
  +----------------------------------------------------------+
```

Every one of those header fields already exists on your `core.Block`, under names you chose in Chapter 17: `Height` and `Timestamp`, `PrevBlockHash`, `MerkleRoot`, `Nonce`, and `Hash`. Bitcoin's header does not store height directly (it is inferred from the header's position in the chain rather than stored as a field), and it encodes its difficulty target in a compact form called **bits** rather than storing the target as a full-width number — a byte-saving trick that matters when a header must be small enough to be handed around and checked by lightweight software.

The header alone is deliberately tiny and fixed-size — this is exactly what lets a **light client** (a wallet on a phone that doesn't store the whole blockchain) download and verify just the headers of every block quickly, and then use Merkle proofs (Chapter 10) to check that one specific transaction it cares about really is included in a block, without ever downloading that block's full transaction list. This is precisely the light-client trick your Chapter 10 Merkle proof implementation was built to demonstrate — Bitcoin's real design is the reason that trick matters at all.

Each individual transaction, in turn, matches the shape you built in Chapter 32:

```
  A REAL BITCOIN TRANSACTION
  +------------------------------------------+
  | version                                    |
  | input count                                |
  | inputs[]                                   |
  |   - previous tx ID + output index          |
  |   - scriptSig (unlocking script)           |
  |   - sequence number                        |
  | output count                               |
  | outputs[]                                  |
  |   - value (in satoshis)                    |
  |   - scriptPubKey (locking script)          |
  | locktime                                   |
  +------------------------------------------+
```

The one field with no direct GoChain counterpart is **locktime**: a value that says "do not consider this transaction valid for inclusion in a block until a certain block height (or time) is reached." It enables a class of contracts built entirely out of ordinary transactions — for example, a payment that cannot be broadcast until 30 days from now — without needing a smart contract VM at all. It's a useful reminder that not every piece of blockchain functionality needs Volume 9's machinery; sometimes a well-placed field on a plain transaction is enough.

The **sequence number** on each input is a related, older mechanism originally intended for a different purpose than it ended up serving — today it's most commonly used to signal opt-in behaviors like replace-by-fee (covered in Section 8) and relative timelocks. It's a good example of a real system accumulating a field whose actual, dominant use case shifted over time — a pattern worth watching for in any codebase you inherit rather than design from scratch yourself.

---

## 4. Bitcoin Script: The Real Ancestor of GoChain's VM

Chapter 60 told you, in a single sentence, that GoChain's stack-based VM was "the same core idea behind Bitcoin Script." This section makes that concrete.

**Bitcoin Script** is a small, stack-based scripting language: each instruction (opcode) either pushes data onto a stack or pops values off the stack, operates on them, and pushes a result back — exactly the execution model your `vm.VM` and its `Execute()` loop implement. The plain-language analogy from Chapter 60 still applies without modification: think of it as a calculator that only ever operates on "whatever's currently on top of the pile," never on named variables sitting elsewhere in memory.

Here is an illustrative (not copy-pasted from any specific client) trace of a standard P2PKH unlock-then-lock script executing, side by side with the equivalent GoChain program from Chapter 63:

```
Illustrative Bitcoin Script (P2PKH spend), representative form:

  scriptSig:     <sig> <pubKey>
  scriptPubKey:  OP_DUP OP_HASH160 <pubKeyHash> OP_EQUALVERIFY OP_CHECKSIG

  Combined execution, stack shown top-first:

  push <sig>            [sig]
  push <pubKey>          [pubKey, sig]
  OP_DUP                 [pubKey, pubKey, sig]
  OP_HASH160             [hash(pubKey), pubKey, sig]
  push <pubKeyHash>      [pubKeyHash, hash(pubKey), pubKey, sig]
  OP_EQUALVERIFY         [pubKey, sig]        (fails script if not equal)
  OP_CHECKSIG            [true/false]         (verifies sig against pubKey)
```

```go
// The same shape, expressed as GoChain vm opcodes (Chapter 61-63),
// illustrative — not literal bytecode from any real client:
program := []vm.Op{
    vm.OpDup,
    vm.OpHash160,
    vm.OpPushData, // pubKeyHash
    vm.OpEqualVerify,
    vm.OpCheckSig,
}
```

The single most important architectural fact about Bitcoin Script is what it deliberately *cannot* do: it has **no loop instructions** (no jump-backward opcodes) and executes a bounded, finite sequence of opcodes with no way to run indefinitely. This is not a limitation someone forgot to lift — it is a design decision to keep script evaluation trivially predictable and cheap for every node to verify, at the cost of expressiveness. Ethereum's EVM (Chapter 82) chose the opposite trade-off — genuine loops and arbitrary computation — and had to invent gas (Chapter 64) specifically to cope with the risk that choice introduces. GoChain's VM sits closer to Bitcoin Script's spirit for its locking scripts, but, like Ethereum, added gas metering (Chapter 64) because its instruction set is closer to general-purpose than Bitcoin Script's intentionally narrow one.

Bitcoin Script opcodes fall into a small number of durable categories, the same categories you designed for GoChain's own instruction set in Chapter 61:

```
  CATEGORY               EXAMPLE REAL OPCODES        GOCHAIN EQUIVALENT
  ---------------------   --------------------------   -------------------
  Stack manipulation      OP_DUP, OP_DROP, OP_SWAP      OpDup, OpPop
  Arithmetic              OP_ADD, OP_SUB                OpAdd, OpSub
  Comparison              OP_EQUAL, OP_GREATERTHAN       OpEqual, OpGreaterThan
  Cryptographic           OP_HASH160, OP_CHECKSIG        OpHash160, OpCheckSig
  Multi-signature          OP_CHECKMULTISIG               (composed from OpCheckSig)
  Timelocks                OP_CHECKLOCKTIMEVERIFY          (no direct equivalent)
  Flow control (limited)  OP_IF, OP_ELSE, OP_ENDIF        OpJumpIfFalse
  Constants               OP_0 .. OP_16, OP_PUSHDATA      OpPushData
```

Two of those categories are worth a sentence each because they show Bitcoin Script solving real problems entirely within its narrow, loop-free design. **OP_CHECKMULTISIG** verifies that some number of valid signatures out of a larger set of authorized public keys are present (an "M-of-N" multisig — for example, 2-of-3 company signers must all agree to release funds), without needing anything beyond the existing stack-and-opcode model. **OP_CHECKLOCKTIMEVERIFY** enforces the same kind of "not spendable until block height/time X" logic as the transaction-level locktime field from Section 3, but as an opcode a script can check directly — letting a locking script itself refuse to unlock early, rather than relying on the spender's transaction to set locktime honestly.

This course intentionally does not enumerate every real opcode or give an exact current opcode count, because that count is an implementation detail of specific Bitcoin software that has changed over time (some opcodes were disabled early in Bitcoin's history for security reasons and never re-enabled). The categories above, however, are structural and durable — they are the same categories any stack-based scripting language for authorizing spends needs, which is exactly why your Chapter 61 design converged on nearly the same list independently.

---

## 5. Transaction Malleability and the SegWit Fix

Chapter 33 spent real effort on one specific, subtle detail: you sign a *trimmed copy* of a transaction, specifically to avoid signing the signature data itself, and the chapter called this "an important detail" without fully spelling out what breaks if you get it wrong. Here is the real-world story that detail is drawn from.

**Transaction malleability** is the property that, for certain transaction formats, someone can take a valid, already-signed transaction and produce a *different-looking* transaction — with a different transaction ID, because the ID is a hash of the whole transaction including its scriptSig — that is still valid and still does exactly the same thing (spends the same inputs, creates the same outputs). This happens because parts of a scriptSig (padding within signature encoding, for instance) can sometimes be altered without invalidating the signature itself, and since the transaction ID hashes the scriptSig along with everything else, changing those bytes changes the ID without changing what the transaction *does*.

```
  Original transaction:      TXID = hash(version, inputs incl. scriptSig, outputs, locktime)
  Malleated transaction:     TXID' = hash(version, inputs incl. TWEAKED scriptSig, outputs, locktime)

  Same spend. Same outputs. Same effect on balances.
  Different TXID.
```

Why this matters in practice: any software that tracked an unconfirmed payment *by its transaction ID* — for example, a service watching for "transaction abc123 to confirm" before shipping a product — could be confused if a malleated version of that same transaction (with a different ID) is the one that actually gets mined, even though the payment genuinely went through. This was never a way to steal funds or double-spend by itself, but it was a real, exploitable annoyance for any system built on the assumption that a transaction's ID is a stable, unchanging reference to it before confirmation.

The real, adopted fix is called **Segregated Witness (SegWit)**: it restructures transactions so that the signature data (the "witness" — the scriptSig-equivalent proving authorization) is stored separately from the data that determines the transaction ID, rather than being hashed together with it. With witness data segregated out, altering it can no longer change the transaction ID, closing the malleability gap entirely for transactions using the new format.

```
  BEFORE (malleable):                    AFTER SEGWIT (witness segregated):
  +---------------------------+          +---------------------------+
  | version, inputs (incl.    |          | version, inputs (refs     |
  |   scriptSig), outputs,    |  hash    |   only, no sig data),     |  hash --> TXID
  |   locktime                | ------>  |   outputs, locktime       |
  +---------------------------+  TXID    +---------------------------+
                                          | witness data (signatures) |  stored
                                          | -- kept separate           |  separately,
                                          +---------------------------+  doesn't affect TXID
```

This is a genuinely useful case study in a durable engineering lesson, independent of Bitcoin specifically: **a hash-based identifier is only as stable as the exact set of bytes it's computed over.** Chapter 33's insistence on signing a trimmed copy rather than the whole (eventually-signed) transaction was solving a closely related problem one layer earlier — avoiding a scenario where your own signature process could accidentally make a transaction's hash depend on data that hadn't stabilized yet. GoChain's simpler transaction format sidesteps the historical malleability problem by design rather than needing a SegWit-style retrofit, precisely because Volume 5 had the benefit of hindsight this real incident provided.

---

## 6. Real Difficulty Adjustment: The 2016-Block Retarget

Chapter 26 built a difficulty-adjustment algorithm "modeled on Bitcoin's." Here is the real one, in full.

Recall the plain-language idea from Chapter 26: if blocks are coming faster than the target block time, raise the difficulty; if slower, lower it. Bitcoin's real version applies that idea over a fixed window rather than continuously:

- Bitcoin targets an average of **10 minutes** per block.
- Every **2016 blocks**, every full node independently recalculates the difficulty, using only information already in the chain — no central authority announces a new difficulty.
- 2016 blocks at 10 minutes each is exactly **two weeks** of blocks, if the network were hashing at exactly the right rate — which is precisely why 2016 was chosen as the window size.

The recalculation compares how long that last 2016-block window *actually* took against the two-week target:

```
  actual_time_taken = timestamp(block 2016) - timestamp(block 0)
  target_time       = 2016 blocks * 10 minutes  (exactly two weeks)

  new_difficulty = old_difficulty * (target_time / actual_time_taken)

  If actual_time_taken < target_time  -->  network was faster than expected
                                       -->  new_difficulty > old_difficulty
                                            (raise it, to slow blocks back down)

  If actual_time_taken > target_time  -->  network was slower than expected
                                       -->  new_difficulty < old_difficulty
                                            (lower it, to speed blocks back up)
```

One durable safety detail your Chapter 26 implementation is worth revisiting in light of: Bitcoin clamps the adjustment so difficulty can change by at most a bounded factor (traditionally at most 4x up or down) in any single retarget, no matter how extreme the measured timing was. Without that clamp, a short burst of unusual timestamps — honest or manipulated — could swing difficulty to an extreme value in one step. This is the same class of defensive thinking behind Chapter 19's tamper-evidence checks: never trust a single data point more than the system's overall design can tolerate being wrong.

```
                 THE 2016-BLOCK RETARGET WINDOW

  Block 0 -----------------------------------> Block 2016
     ^ window starts                                ^ window ends, recompute here
     |                                               |
     |<---------- ~2 weeks, if on schedule --------->|

  If the window took only 10 days (network too fast):
     difficulty goes UP  -->  next window's blocks slow back toward 10 min

  If the window took 20 days (network too slow, e.g. after miners left):
     difficulty goes DOWN -->  next window's blocks speed back toward 10 min
```

Your GoChain implementation from Chapter 26 almost certainly uses a shorter, simpler window (a smaller retarget period, useful for a teaching chain where you don't want to wait two weeks to see an adjustment happen) — that's the correct simplification for a course, not a flaw. The mechanism — measure actual time over a fixed window of blocks, compare to target, scale difficulty proportionally, clamp the adjustment — is identical in spirit either way.

---

## 7. Node Types: Full, Pruned, and Lightweight

Volume 8 built one storage story for GoChain: every node keeps its own BoltDB-backed copy of the chain and a UTXO index. Real Bitcoin deployments actually run several different tiers of node, trading off storage cost against independence and trust, and the vocabulary is worth having precisely.

A **full node** stores and independently validates every block since genesis, checking every rule (Chapter 19-style block validation, script validation, difficulty checks) itself, and trusts no other node's word for anything. This is the model GoChain assumes throughout the course — every `core.Blockchain` you've run has been acting as a full node.

A **pruned node** is also a full node in terms of validation — it downloads and checks every block just like a full node does — but after validating a block, it discards the old block data it no longer needs for future validation, keeping only a smaller recent window plus the UTXO set (the current spendable balances) needed to keep validating new blocks. It trusts nothing it hasn't checked itself; it just doesn't keep the full historical record sitting on disk. This maps directly onto the distinction your `storage.Store` interface (Chapter 55) was designed to make swappable — a pruned mode is, architecturally, "the same validation logic, a smaller retention policy in the storage layer."

A **lightweight (SPV) client** — "Simplified Payment Verification" — does not validate full blocks at all. It downloads only block headers (Section 3's fixed-size header, the reason it can afford to be tiny), and relies on Merkle proofs (Chapter 10) requested from full nodes to confirm that a specific transaction it cares about is included in a block whose header it has independently checked (proof of work satisfied, chain correctly linked). It trusts the *network's aggregate work* implicitly, without re-validating every rule of every transaction ever mined — a meaningfully weaker trust model, appropriate for a mobile wallet that cannot store hundreds of gigabytes of history, but not appropriate for, say, an exchange's back-end ledger.

```
   TRUST AND STORAGE TRADE-OFF ACROSS NODE TYPES

   FULL NODE            PRUNED NODE           SPV / LIGHT CLIENT
   validates every       validates every        validates only
   block ever, keeps      block, keeps only       headers + Merkle
   all of it on disk       a recent window +       proofs for txs
                            UTXO set                 it cares about
   ---------------        ---------------          ---------------
   most storage            moderate storage          minimal storage
   most trust-free          still trust-free          trusts aggregate
                                                        proof-of-work only
```

GoChain's course design has, until now, only ever asked you to build the leftmost column — every `core.Blockchain` instance you've run validates everything and stores everything. That was the right call for learning: pruning and SPV are optimizations layered on top of full validation, and you cannot meaningfully understand what to safely discard until you've built the thing that keeps everything and checks everything itself.

---

## 8. Mempool Policy on Real Bitcoin Nodes

Chapter 34 built a mempool that tracks pending transactions and rejects double-spends. Chapter 35 added fee-based prioritization. Real Bitcoin nodes layer several additional, practical policies on top of that same foundation — none of which are part of Bitcoin's *consensus* rules (a block isn't invalid for violating them), but which shape which transactions actually make it into the network's mempools and eventually into blocks.

**Minimum relay fee.** A node will typically refuse to relay (or hold in its own mempool) a transaction whose fee-per-byte falls below a locally configured minimum — protecting the node's own resources from being filled with transactions nobody will ever mine. This is a direct, real-world extension of the fee-sorting logic you built in Chapter 35: fee-per-byte doesn't just decide mining *order*, it can decide relay *eligibility* in the first place.

**Standardness rules.** Not every script that would be valid under Bitcoin's consensus rules is relayed by default nodes — many nodes only forward transactions using a small set of common, well-understood script templates (like the P2PKH pattern from Section 2), rejecting unusual-but-technically-valid scripts from their mempool and relay behavior as an extra layer of caution against unknown attack surface. A miner could still choose to mine an unusual, "non-standard" but valid transaction directly, but it won't casually propagate node to node the way a standard transaction does.

**Replace-by-fee (RBF).** A sender who broadcast a transaction with too low a fee (and sees it stuck, unconfirmed) can, if they signaled intent to allow it (via the sequence number from Section 3), broadcast a new transaction spending the same inputs with a higher fee, replacing the original in mempools across the network. This is a controlled, explicit exception to the "first one wins" mempool rule you'd otherwise expect from Chapter 34's double-spend rejection — the key difference from an actual double-spend attack is that RBF only replaces a transaction *before* it's mined, and the network-wide policy for exactly when replacement is allowed is itself a deliberate design choice, not an accident.

**Ancestor/descendant limits.** A mempool tracks not just individual transactions but chains of dependent, still-unconfirmed transactions (a transaction spending an output from another transaction that hasn't been mined yet). Nodes cap how deep and how large these unconfirmed chains are allowed to grow, to keep worst-case mempool re-validation work bounded — the same instinct behind Chapter 64's gas limits, applied to mempool bookkeeping instead of VM execution.

```
   MEMPOOL POLICY LAYERS (none of these are consensus rules)
   +--------------------------------------------------------+
   | 4. Ancestor/descendant limits — bound unconfirmed chains |
   | 3. Replace-by-fee — controlled, signaled replacement      |
   | 2. Standardness rules — only relay known-safe templates   |
   | 1. Minimum relay fee — floor on fee-per-byte to propagate  |
   +--------------------------------------------------------+
   | core.Mempool: double-spend rejection (Ch.34), fee sort   |
   |               and greedy block-filling (Ch.35)            |
   +--------------------------------------------------------+
```

None of this changes what makes a block *valid* — that's still governed entirely by consensus rules any node can check independently (Chapter 19, Chapter 25). Mempool policy only governs what an individual node chooses to *hold onto and forward* before mining happens, which is exactly why different nodes' mempools can differ from each other at any given moment without the network disagreeing about history.

---

## 9. Bitcoin vs. GoChain: Side-by-Side

```
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | BITCOIN (real)              | GOCHAIN (this course)     |
  +----------------------+---------------------------+---------------------------+
  | Data model            | UTXO                        | UTXO (same model)          |
  | Smallest unit          | satoshi                     | gochip                      |
  | Locking mechanism      | scriptPubKey/scriptSig       | LockingScript/unlock (Ch63) |
  | Scripting language     | Bitcoin Script (no loops)    | gochain/vm (bounded, gas)   |
  | Consensus              | Proof of work                | Proof of work (+ opt. PoS)  |
  | Retarget window        | 2016 blocks (~2 weeks)        | shorter, teaching window    |
  | Header size             | fixed, small (headers-only    | fixed-size header fields    |
  |                         | sync for light clients)       | (Chapter 17)                 |
  | Node types              | full / pruned / SPV            | full only (course scope)     |
  | Malleability             | historical issue, fixed by     | avoided by design (no        |
  |                          | SegWit's witness segregation    | witness/scriptSig split)     |
  | Mempool                | fee floor, standardness,      | fee sort + double-spend      |
  |                         | RBF, ancestor limits           | rejection (Ch. 34-35)        |
  | Curve / signatures      | ECDSA over secp256k1           | ECDSA (Go's crypto/ecdsa)    |
  +----------------------+---------------------------+---------------------------+
```

Reading this table, the honest takeaway is not "Bitcoin is complicated and GoChain is simple" as a criticism — it's that GoChain built the same *shape* of every one of these mechanisms, at a scale and level of polish appropriate for one person to build and fully understand in a few months, rather than the scale a globally-distributed, adversarial, trillion-dollar network needs after more than a decade of hardening.

---

## 10. What GoChain Simplified, and Why That Was the Right Call

It's worth naming the simplifications explicitly, because each one was a deliberate teaching decision, not an oversight:

- **A shorter difficulty retarget window.** Waiting two real weeks to see a difficulty adjustment would make Chapter 26 impossible to demonstrate hands-on. A short window teaches the identical algorithm at a pace a single learner mining on a laptop can actually observe.
- **No mandatory secp256k1-specific curve library.** Bitcoin standardized on one specific elliptic curve, secp256k1, and real implementations depend on highly optimized, audited libraries for it. GoChain uses Go's standard `crypto/ecdsa` package, which supports standard curves out of the box, prioritizing "runs correctly with zero extra dependencies" over "matches Bitcoin's exact curve choice" — the cryptographic *concepts* (Chapters 11-13) transfer regardless of which specific curve sits underneath.
- **No mandatory standardness/relay policy layer.** GoChain's mempool (Chapters 34-35) implements the consensus-relevant core — double-spend rejection and fee-based ordering — without also implementing a full policy layer of minimum relay fees, standardness templates, and replace-by-fee, because those are operational refinements on top of the core idea, not new concepts.
- **A smaller, more general instruction set with gas from the start**, rather than Bitcoin Script's intentionally narrow, loop-free design. This mirrors Ethereum's choice more than Bitcoin's (more on that trade-off in Chapter 82) because a slightly more expressive VM makes richer example contracts (Chapter 65's token) easier to teach.
- **Only one node type — the equivalent of "full node."** GoChain doesn't implement pruning or SPV-style light clients, because those are storage and bandwidth optimizations on top of full validation, and Volume 8's whole point was building and understanding the full-validation storage layer those optimizations sit on top of.
- **No historical malleability problem to retrofit a fix for.** Because GoChain's transaction ID and signing scheme were designed from Chapter 33 onward with Bitcoin's malleability history already known, GoChain never needed a SegWit-style redesign — a benefit of hindsight, not a claim that GoChain's design is more sophisticated than Bitcoin's was in 2009.

Every one of these is exactly the kind of decision a real engineering team makes too — optimize for the problem you actually have. Bitcoin's specific choices were optimized for being the first, most conservative, most battle-tested public monetary network in the world, iterating carefully under the constraint that changes must stay compatible with a live network worth enormous real value. GoChain's choices were optimized for you understanding, completely, everything happening under the hood. Chapter 83 makes this comparison exercise even more systematic, running it against both Bitcoin and Ethereum side by side.

---

## Summary

- Bitcoin's UTXO model, locking/unlocking scripts (scriptPubKey/scriptSig), and P2PKH spending pattern are structurally identical to GoChain's UTXO model and Chapter 63's locking scripts — GoChain's design was modeled directly on Bitcoin's, down to coinbase maturity mirroring the reorg-safety concerns behind Chapter 50.
- Bitcoin's real block splits into a small, fixed-size header (version, previous hash, Merkle root, timestamp, bits, nonce) and a separate transaction list — the same fields your `core.Block` already has, under the same names.
- Bitcoin Script is the real, production ancestor of GoChain's VM: a stack-based language, deliberately without loops, bounded and cheap to verify, with special-purpose opcodes like `OP_CHECKMULTISIG` and `OP_CHECKLOCKTIMEVERIFY` solving real problems within that narrow design.
- Transaction malleability — the same TXID could describe subtly different byte sequences — was a real historical problem, fixed by Segregated Witness moving signature data out of the transaction-ID hash; it's the production-scale version of the "sign a trimmed copy" caution from Chapter 33.
- Bitcoin's real difficulty adjustment recalculates every 2016 blocks (about two weeks at a 10-minute target), comparing actual elapsed time to the target and scaling difficulty proportionally, with a clamp on how much any single adjustment can move.
- Real Bitcoin deployments run full nodes, pruned nodes, and lightweight SPV clients, trading storage for trust — GoChain's course scope only builds the full-node equivalent, which is the correct foundation to build before layering on those optimizations.
- Real Bitcoin mempool policy — minimum relay fees, standardness rules, replace-by-fee, ancestor/descendant limits — sits on top of the same core double-spend rejection and fee-sorting your Chapter 34-35 mempool implements, and is not part of Bitcoin's consensus rules.
- GoChain's simplifications were each deliberate teaching trade-offs, not missing features — the same shape of mechanism, built at a scale one learner can fully own and understand.

---

## Exercises

### Easy

1. **Draw GoChain's `core.Block` fields next to Bitcoin's real block header fields** (Section 3), matching each GoChain field to its Bitcoin counterpart by name. Note the one Bitcoin header field that has no direct stored equivalent on your `core.Block` (hint: it's inferable from position in the chain instead).

2. **Trace the P2PKH script execution in Section 4 by hand** on paper, writing out the stack's contents after every single opcode, for a made-up signature and public key of your choosing (they don't need to be real cryptographic values — just labels like `sig1` and `pub1`).

3. **In 100-150 words, explain why Bitcoin Script deliberately has no loop instructions**, and why that omission makes script verification easier for every node on the network to reason about.

### Medium

4. **Read the "Block chain" and "Transactions" sections of the original Bitcoin whitepaper** (Satoshi Nakamoto, 2008) and write a short summary (150-250 words) connecting at least three specific ideas in it to specific chapters of this course (for example: "Section 2 describes chaining hashed blocks, which is exactly Chapter 18's `PrevBlockHash` linking").

5. **Work through the difficulty-adjustment formula in Section 6 with concrete numbers**: assume `old_difficulty = 1000`, the 2016-block window's target time is 20,160 minutes (2 weeks), and the window actually took 15,000 minutes. Compute the new difficulty (before any clamp), then explain in your own words whether this result makes sense given that the network apparently mined blocks faster than the 10-minute target.

6. **Explain transaction malleability in your own words (150-200 words)** using a non-cryptocurrency analogy of your own devising (something other than the ones used in this chapter), and describe specifically why segregating witness data out of the transaction-ID hash closes the problem.

### Hard

7. **Implement a minimal "standardness check"** for GoChain's mempool: reject any transaction whose locking script isn't the standard P2PKH pattern from Chapter 63 before it's admitted to `core.Mempool`, and write a test proving a transaction with an unusual, hand-crafted locking script is rejected by the mempool but would still pass `core.Block` validation if it somehow made it into a mined block (demonstrating that this really is a relay-policy choice, not a consensus rule).

8. **Implement the real 2016-block-style retarget algorithm exactly** (not a shortened teaching window) as an alternate, configurable difficulty-adjustment mode in `consensus.ProofOfWork`, parameterized by window size and target block time, and write a test that exercises it with a fast window, a slow window, and a window extreme enough to hit your implementation's clamp.

9. **Design (on paper, no implementation required) a "pruned mode" for GoChain's `storage.Store`** (Chapter 55): describe exactly what data it would keep versus discard after validating a block, what new failure mode it would introduce (for example, being unable to serve old blocks to a syncing peer per Chapter 49), and how a pruned GoChain node would need to signal its limitation to peers requesting old data it no longer has.
