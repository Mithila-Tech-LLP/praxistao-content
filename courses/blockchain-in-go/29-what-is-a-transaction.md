# Chapter 29: What Is a Transaction?

A transaction is a signed instruction that moves value from one place to another — but unlike a bank transfer, it has to convince every stranger running a copy of the network, independently, using nothing but math. This chapter defines exactly what a transaction must contain to earn that trust, before Chapter 32 turns the definition into real Go code.

## Table of Contents

1. [Blocks Full of What, Exactly?](#1-blocks-full-of-what-exactly)
2. [A Transaction as a Signed Instruction](#2-a-transaction-as-a-signed-instruction)
3. [The Signed-Check Analogy](#3-the-signed-check-analogy)
4. [Where the Analogy Breaks: Math Instead of a Bank](#4-where-the-analogy-breaks-math-instead-of-a-bank)
5. [The Three Things Every Trustworthy Transaction Proves](#5-the-three-things-every-trustworthy-transaction-proves)
6. [Inputs and Outputs, Previewed](#6-inputs-and-outputs-previewed)
7. [From Wallet to Block: A Transaction's Journey](#7-from-wallet-to-block-a-transactions-journey)
8. [Where GoChain Fits](#8-where-gochain-fits)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Blocks Full of What, Exactly?

Since Volume 3, `core.Block` has carried a `Transactions []*Transaction` field, and since Volume 4 you have been mining blocks by calling `bc.MineBlock(transactions)`. Up to this point, though, we quietly glossed over what actually lives inside a `Transaction` — we just needed *something* to hash and put in a Merkle tree so we could focus on blocks and proof of work without getting distracted.

That deferral ends now. A blockchain that only knows how to chain together empty or placeholder blocks is not very useful. The entire point of GoChain — like Bitcoin before it — is to let people move value between each other without a bank, a payment processor, or any single company standing in the middle. That value-moving instruction is called a **transaction**, and this volume builds it for real: what it contains, how it proves it's legitimate, and how every node on the network can check it without trusting anyone.

---

## 2. A Transaction as a Signed Instruction

Strip away the cryptography for a moment and a transaction is a simple sentence: *"Move some amount of value from here to there."* What turns that plain sentence into something a stranger can trust, with no intermediary, is two additions:

- **A signature** — cryptographic proof that the sentence was actually written (authorized) by the person who controls the funds being moved, not by an impostor claiming to be them.
- **A reference to where the value already exists** — you cannot just say "send Bob 10 gochips" out of thin air; the instruction has to point at value that genuinely exists and genuinely belongs to the sender.

A **transaction**, then, is a signed instruction that consumes existing value the sender is entitled to and creates new value assigned to a recipient (and possibly back to the sender, as we'll see in Chapter 30). It's signed the moment it's created, and every future person who looks at it — a miner deciding whether to include it in a block, another node re-validating an already-mined block — checks that signature themselves. Nobody has to take the sender's word for it, and nobody has to take *any single node's* word for it either.

---

## 3. The Signed-Check Analogy

A paper check is a useful mental model, because it already has most of the shape we need. Look at a check and you'll find:

```
+--------------------------------------------------+
|  ALICE JOHNSON                     No. 1042       |
|  123 Main St                                       |
|                                                    |
|  Pay to the                                        |
|  order of      Bob Smith                $ 10.00   |
|                                                    |
|  Ten and 00/100 ---------------------------- DOLLARS |
|                                                    |
|  Memo: lunch                    Alice Johnson     |
|                                  (signature)       |
|                                                    |
|  :123456789: 0001042"" ACCOUNT #55512340""        |
+--------------------------------------------------+
```

Every check names an account the money is supposed to come from, a payee it should go to, an amount, and a signature that (in theory) only the account holder can produce. A bank teller cashing this check does three things: checks the signature looks like Alice's, checks Alice's account actually has at least $10.00 in it, and then debits Alice's account and credits Bob's.

A GoChain transaction plays exactly the same role. It says, in effect: "the owner of this specific chunk of value, proven by this signature, authorizes moving it to this recipient." Keep this picture in your head — source, destination, amount, signature — because Chapter 32 turns it directly into Go struct fields.

---

## 4. Where the Analogy Breaks: Math Instead of a Bank

Here is the crucial place the check analogy stops working, and it is the entire reason blockchains are interesting rather than just "digital checks."

When a real bank cashes Alice's check, exactly one institution — her bank — decides whether the signature is good and whether the funds exist. You trust that institution to be honest, competent, and available. If it makes a mistake, colludes with someone, or goes down, you have limited recourse.

A blockchain transaction has no single verifying institution. Instead:

```
        Alice's signed transaction
                   |
                   v
      +------------+------------+------------+
      |            |            |            |
      v            v            v            v
   Node A        Node B       Node C       Node D
  (checks         (checks      (checks      (checks
   signature      signature    signature    signature
   & source        & source     & source     & source
   independently)  independently) independently) independently)
```

Every single node — every independent copy of GoChain running anywhere in the world — receives the same transaction and runs the exact same two checks itself: *is the signature mathematically valid for this data and this public key?* and *does the referenced source of funds genuinely exist and genuinely belong to whoever signed this?* Nobody phones a central authority to ask. Nobody has to be trusted, because the verification is math (elliptic-curve signature verification, from Volume 2's `crypto.Verify`), not opinion. A node run by a stranger on the other side of the planet, who has never heard of you and has no reason to like you, will still accept your transaction — provided the math checks out — because the rules apply identically to everyone, including that node's own operator.

This is the payoff of everything Volumes 2 through 4 built: hashing (to fingerprint data), signatures (to prove authorization without a password ever being shared), and proof of work (to make the shared history expensive to rewrite). Transactions are where all three finally combine to do something people actually want: move value, trustlessly.

---

## 5. The Three Things Every Trustworthy Transaction Proves

Distilling the check analogy and the "math instead of a bank" difference into a precise checklist, a transaction that any node can independently trust must supply:

1. **Proof of where the funds come from.** Not a balance figure typed in by the sender (anyone could type a big number), but a concrete reference to specific, previously-existing value — something every node can look up and confirm is real and has not already been spent.
2. **Proof the sender is authorized to spend it.** A cryptographic signature, produced with the private key that controls the referenced funds, over the exact details of this transaction. Anyone holding only the corresponding *public* key (which is safe to share, unlike a private key or a bank password) can verify this signature without ever learning the private key itself — this is the same asymmetric-cryptography idea from Chapter 11.
3. **Where the funds should go.** One or more destinations, and how much value each one receives.

If any one of these three is missing, the transaction is not trustworthy on its own — a node would have to *trust* something outside the transaction itself (the sender's word, a central ledger, a middleman) to accept it. The entire design goal of this volume is making all three checkable by pure computation, using data any node can already see.

---

## 6. Inputs and Outputs, Previewed

Bitcoin — and GoChain, following it — has a specific, elegant way of satisfying point 1 and point 3 above: **inputs** and **outputs**.

- An **output** is a chunk of value created by some transaction, assigned to whoever controls a particular address. Think of it as a sealed envelope of money sitting out there, waiting to be claimed by its rightful owner.
- An **input** is a reference to one specific earlier output, together with the proof (the signature from point 2) that the person spending it is actually entitled to.

A new transaction's outputs become the *only* legitimate sources of funds that future transactions can reference as inputs. This is the seed of an idea called the **UTXO model** — Unspent Transaction Output — which Chapter 30 covers in full, worked-example depth. For this chapter, just hold onto the shape: **inputs point backward** at value that already exists, and **outputs point forward**, creating new value for someone to spend later.

```
   Earlier transaction              New transaction
  +------------------+          +----------------------+
  |  ... (inputs)    |          |  Input:              |
  |                  |          |   -> points at output |
  |  Output #0:      | <------- |      #0 of the        |
  |  5 gochips to    |          |      earlier tx       |
  |  Alice           |          |                        |
  +------------------+          |  Output #0:          |
                                 |   5 gochips to Bob    |
                                 +----------------------+
```

---

## 7. From Wallet to Block: A Transaction's Journey

It helps to see the full lifecycle a transaction goes through before it becomes a permanent part of GoChain's history, since the rest of this volume builds each stage in order:

```
 1. A wallet builds an unsigned transaction        (Chapter 32)
              |
              v
 2. The wallet signs it with the sender's          (Chapter 33)
    private key
              |
              v
 3. The transaction is broadcast into the           (Chapter 34)
    mempool -- a waiting room of valid, signed,
    not-yet-mined transactions
              |
              v
 4. A miner selects it (often based on fee,         (Chapters 34-35)
    covered soon) and includes it in a
    candidate block
              |
              v
 5. The block is mined (proof of work,              (Volume 4, already built)
    Chapter 25) and added to the chain
              |
              v
 6. Every node independently re-verifies the         (this chapter's
    transaction's signature and sources              Section 4, made real)
    before accepting the new block
```

Notice that step 6 is not optional and does not stop happening once a transaction is mined once. Any node that later downloads the chain from scratch (Volume 7 covers this) re-runs exactly the same signature and source checks. Trust is never "spent once and forgotten" — it is re-derived from scratch, from math, every single time anyone looks at the chain.

---

## 8. Where GoChain Fits

Starting in Chapter 32, `core.Transaction` becomes a real Go type with `Inputs []TxInput` and `Outputs []TxOutput` fields that implement exactly the input/output shape previewed in Section 6. Chapter 30 first builds the mental model those types encode (the UTXO model) with a full numeric example. Chapter 31 explains why we chose this model over Ethereum's alternative. Chapter 33 implements the signing and verification that makes Section 5's "proof of authorization" real, working Go code. By Chapter 37 — the end of this volume — you will have two wallets sending gochips back and forth, with every node-independent check from this chapter running for real.

---

## Summary

- Blocks have carried a `Transactions` field since Volume 3, but this volume is where `core.Transaction` finally becomes real.
- A transaction is a **signed instruction** that moves value, combining a reference to existing funds with cryptographic proof of authorization.
- The **signed-check analogy** captures the shape (source, destination, amount, signature) but breaks down at the verifier: a bank is one trusted institution, while a blockchain is verified independently by every node using math, not trust.
- A trustworthy transaction must supply three things: **proof of source**, **proof of authorization** (a signature), and **a destination**.
- GoChain represents source and destination using **inputs** (references to earlier value) and **outputs** (newly created value for someone to spend later) — the seed of the UTXO model covered fully in Chapter 30.
- Verification is never "done once and trusted forever" — every node re-checks every transaction's signature and sources independently, even when re-validating old, already-mined blocks.
- The rest of this volume builds, in order: the UTXO model (30), why we chose it over the account model (31), the real Go types and construction logic (32), and real signing/verification (33).

---

## Exercises

### Easy

1. Write out, in your own words and without using the word "blockchain," the three things a transaction must prove to be trustworthy without a central authority. Then match each one to the corresponding part of a paper check.
2. Draw (on paper or in a text file) the "signed-check" diagram from Section 3, but replace every field with a GoChain equivalent: instead of a bank account number, what would identify the source of funds? Instead of a payee name, what would identify the recipient? You do not need real syntax — just label the boxes.
3. Explain in 3-4 sentences why a transaction that only contained "send Alice 10 gochips" (no signature, no reference to existing funds) would be worthless on a real network, even if it were technically well-formatted.

### Medium

4. Section 4 argues that verification must happen independently on every node rather than through one central authority. Describe a concrete scenario where a single central verifier (like a bank) could behave dishonestly in a way that independent, math-based verification specifically prevents.
5. Re-read the transaction lifecycle diagram in Section 7. For each of the six steps, write one sentence describing what could go wrong if that step were skipped entirely (for example: what breaks if step 2, signing, is skipped?).
6. Using the input/output diagram from Section 6, sketch (on paper) a chain of three transactions: Transaction A creates an output paying Alice 10 gochips. Transaction B has one input pointing at Transaction A's output, and creates one output paying Bob 10 gochips. Transaction C has one input pointing at Transaction B's output, and creates one output paying Carol 10 gochips. Label every arrow.

### Hard

7. A malicious node could, in principle, tell a wallet "yes, your transaction was accepted and mined" without it actually having happened. Explain why this lie only fools a wallet that trusts a single node's word, and describe what a wallet should do instead to avoid being deceived this way (hint: think about how many independent nodes it could ask, and what Volume 7's networking will eventually enable).
8. Research (from general knowledge) how a real bank wire transfer is authorized and cleared internally, and write a 200-300 word comparison against the transaction lifecycle in Section 7, focusing specifically on *who* is trusted at each step and what happens if that party is unavailable or dishonest.
9. Suppose GoChain's design allowed a transaction to reference a source of funds without any signature at all, relying instead on the recipient simply promising, in good faith, that the source was legitimate. Write a short scenario describing exactly how an attacker would exploit this, and explain precisely which of the "three things a trustworthy transaction proves" (Section 5) was missing that made the exploit possible.
