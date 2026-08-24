# Chapter 31: The Account Model, and Choosing One for GoChain

Bitcoin's UTXO model is not the only way to track who owns what. Ethereum instead gives every address a single running balance, updated in place, exactly like a bank account ledger. This chapter puts both models side by side, weighs their real tradeoffs, and explains — with concrete reasons, not just "because Bitcoin does it" — why GoChain's core ledger uses UTXO.

## Table of Contents

1. [The Account Model, Explained](#1-the-account-model-explained)
2. [Side by Side: The Same Four Transactions, Two Ways](#2-side-by-side-the-same-four-transactions-two-ways)
3. [UTXO's Strengths: Parallelism and Privacy](#3-utxos-strengths-parallelism-and-privacy)
4. [The Account Model's Strength: Simplicity for Smart Contracts](#4-the-account-models-strength-simplicity-for-smart-contracts)
5. [Other Practical Differences](#5-other-practical-differences)
6. [Why GoChain Chooses UTXO](#6-why-gochain-chooses-utxo)
7. [A Preview: Account-Like State for Contracts in Volume 9](#7-a-preview-account-like-state-for-contracts-in-volume-9)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Account Model, Explained

Where Bitcoin tracks a scattered pile of indivisible UTXOs, Ethereum tracks something much closer to what you'd expect from an ordinary bank: a single table mapping each address directly to a current balance, updated in place every time value moves.

```
   THE ACCOUNT MODEL (a running ledger)

   +---------+-----------+
   | Address | Balance   |
   +---------+-----------+
   | Alice   | 50        |
   | Bob     | 0         |
   | Carol   | 0         |
   +---------+-----------+
```

A transaction in this model is much simpler to describe: "subtract 20 from Alice's balance, add 20 to Bob's balance." There is no concept of a UTXO, no inputs referencing specific earlier outputs, and no change — because there is nothing to give change *from*. You just have a number, and the number goes down when you spend and up when you receive, the exact way most people already think about money before they ever encounter Bitcoin.

An **account**, in this model, is simply an address plus its current balance (and, for contract accounts, its code and storage — more on this in Section 7). Updating an account's balance in place is sometimes called a **state transition**: the network agrees on a new "state" (the full table of every account's balance) after each transaction, rather than agreeing on a growing pile of unspent chunks.

---

## 2. Side by Side: The Same Four Transactions, Two Ways

Let's replay Chapter 30's worked example — Alice mines 50 gochips, sends Bob 20, Bob sends Carol 5, Alice sends Carol 30, Carol sends Alice 8 — but this time under the account model, so the contrast is concrete rather than abstract.

```
                    UTXO MODEL                    |              ACCOUNT MODEL
--------------------------------------------------|--------------------------------------------
Start:   UTXO#1: 50 -> Alice                      | Alice=50, Bob=0, Carol=0
                                                   |
Tx1: Alice sends Bob 20                           | Tx1: Alice sends Bob 20
  in:  UTXO#1 (50)                                 |   Alice.balance -= 20
  out: 20->Bob, 30->Alice (change)                |   Bob.balance   += 20
  Set: {20->Bob, 30->Alice}                        | Alice=30, Bob=20, Carol=0
                                                   |
Tx2: Bob sends Carol 5                            | Tx2: Bob sends Carol 5
  in:  Bob's 20 UTXO                               |   Bob.balance   -= 5
  out: 5->Carol, 15->Bob (change)                  |   Carol.balance += 5
  Set: {30->Alice, 5->Carol, 15->Bob}               | Alice=30, Bob=15, Carol=5
                                                   |
Tx3: Alice sends Carol 30                         | Tx3: Alice sends Carol 30
  in:  Alice's 30 UTXO (exact, no change)          |   Alice.balance -= 30
  out: 30->Carol                                   |   Carol.balance += 30
  Set: {5->Carol, 15->Bob, 30->Carol}               | Alice=0, Bob=15, Carol=35
                                                   |
Tx4: Carol sends Alice 8                          | Tx4: Carol sends Alice 8
  in:  Carol's 5 + 30 UTXOs (combined)              |   Carol.balance -= 8
  out: 8->Alice, 27->Carol (change)                |   Alice.balance += 8
  Final Set: {15->Bob, 8->Alice, 27->Carol}          | Alice=8, Bob=15, Carol=27
```

Both models arrive at the identical final answer — Alice = 8, Bob = 15, Carol = 27 — because they're describing the same real-world movement of value. What differs entirely is *how the system gets there and what it has to track along the way*. The account model never mentions "which specific chunk" of Carol's money is being spent; it just decrements a single number. The UTXO model never has a single number at all; balance is always a derived sum.

---

## 3. UTXO's Strengths: Parallelism and Privacy

**Natural parallelism.** Look again at Tx1 and Tx2 above under the UTXO model: Tx1 consumes UTXO#1 and creates two new UTXOs. Tx2 consumes one of *those* new UTXOs. But if instead two *unrelated* transactions each spent a UTXO that neither one touches, a node validating them has no shared, mutable state to coordinate over — each transaction only reads and destroys UTXOs it specifically references. This means many independent transactions can, in principle, be validated concurrently by separate goroutines (or separate CPU cores, or even separate machines) with no risk of one transaction's update clobbering another's, because there is no single shared balance field anyone is racing to update. Under the account model, two transactions that both touch Alice's account balance *do* have to be applied one after another, in some agreed order, or her final balance could come out wrong — a coordination problem that grows with how "hot" (frequently used) a given account is.

**A cleaner privacy story.** Under the account model, every transaction touching Alice's account is trivially linkable to the same, single, long-lived "Alice" balance field — anyone watching the chain sees a running total tied to one address across her entire history. Under the UTXO model, GoChain (like Bitcoin) encourages generating a *fresh* address for change and for receiving payments, so Alice's funds are scattered across many UTXOs on many addresses that are not obviously all "the same person" without extra off-chain analysis. This is a modest improvement, not true anonymity — sophisticated chain analysis can often still link UTXOs together — but it is a genuine structural difference favoring UTXO.

---

## 4. The Account Model's Strength: Simplicity for Smart Contracts

The account model's biggest advantage shows up the moment you want to run more than simple transfers: **smart contracts**, which Volume 9 covers in depth.

A smart contract needs its own persistent storage — a token contract, for example, needs to remember *every holder's* token balance between calls, and update many balances arbitrarily as a single contract call executes complex logic (imagine a decentralized exchange trade touching four different users' balances in one atomic operation). This maps onto the account model almost for free: a contract is just another account, with a balance field *and* a storage area, and a contract call is simply a sequence of ordinary balance/storage updates, exactly like the account-model transactions in Section 2 — just orchestrated by contract code instead of a plain transfer instruction.

Representing the same thing in a pure UTXO model is markedly more awkward. There's no natural place to keep "the running list of every token holder's balance" as a single mutable table, because UTXO deliberately has no mutable, shared state at all — only discrete chunks. Real UTXO-based smart contract systems exist, but they generally have to work *around* this limitation rather than *with* it, often by encoding state transitions cleverly across chains of UTXOs. This is a large part of why Ethereum, designed with rich smart contracts as a first-class goal from day one, chose the account model over Bitcoin's UTXO approach.

---

## 5. Other Practical Differences

A few more differences are worth naming plainly, since they come up constantly once you start reading about either Bitcoin or Ethereum:

- **Transaction size and complexity.** A UTXO transaction that needs to combine many small UTXOs (imagine an address that received hundreds of tiny payments) can become large, since every input must be listed and signed individually. An account-model transaction is a constant, small size regardless of an account's history — it's just "from, to, amount."
- **Replay and ordering.** Because UTXOs are consumed exactly once, ever, the UTXO model has a built-in, structural defense against accidentally processing the same transaction twice — a spent UTXO simply cannot be spent again. The account model needs an explicit extra field (a **nonce**, a per-account transaction counter) to achieve the same guarantee, since "subtract 20 from Alice" is a generic instruction that says nothing about whether it's already been applied.
- **State size.** The UTXO set only grows by the *unspent* remainder of history — fully-spent chains of transactions stop contributing new entries. Ethereum's account state must retain a live entry for every account that has ever held a nonzero balance or any contract storage, which has historically made Ethereum's total "state" larger and more expensive to store and synchronize than Bitcoin's UTXO set at comparable levels of usage.

---

## 6. Why GoChain Chooses UTXO

Putting Sections 2 through 5 together, GoChain adopts the **UTXO model for its core ledger**, matching Bitcoin, for three concrete reasons:

1. **This course's core teaching goal is understanding Bitcoin's design deeply, end to end.** UTXO is the model that best illustrates the ideas this course has built toward since Chapter 1 — hash chains, signatures binding to specific, unforgeable references rather than generic instructions, and a genuinely decentralized, parallel-friendly ledger. Building it yourself, by hand, in Go, is the single best way to understand *why* Bitcoin looks the way it does.
2. **UTXO's natural parallelism and cleaner privacy story are real, structural wins for a base value-transfer layer**, and they cost nothing extra once you've built the machinery (Chapter 32's `UTXOSet` and Chapter 34's mempool checks) — a plain "send money" system doesn't need mutable shared state, so there's no reason to pay for the coordination complexity the account model requires.
3. **The account model's main advantage — simplicity for smart contracts — doesn't need to apply to the entire chain.** GoChain does not need every plain coin transfer to carry the coordination overhead a general-purpose contract platform requires, just to make the eventual contract layer easier to build.

---

## 7. A Preview: Account-Like State for Contracts in Volume 9

That third point deserves an important clarification, because it is not "GoChain avoids the account model forever." When Volume 9 introduces smart contracts and a virtual machine, contracts will get their own **account-like state** — persistent storage keyed by contract address, exactly as described in Section 4 — layered *on top of* the UTXO-based core ledger, not replacing it.

Concretely: plain coin transfers between ordinary addresses will keep working exactly as this volume builds them, as UTXOs being consumed and created. Smart contracts, introduced starting in Chapter 59, will additionally get persistent storage slots (Chapter 66's `SLOAD`/`SSTORE` opcodes) that behave like a mutable account balance table, scoped to each contract, because that's genuinely the right tool for the job contracts need to do. This mirrors a real, deliberate design pattern: even some real UTXO-based systems (like Cardano's "extended UTXO" model) add account-like elements specifically to support smart contracts, rather than picking one model and forcing every use case through it. GoChain does the same thing on a smaller, more understandable scale — UTXO where UTXO is the better fit, a lightweight account-like layer where contracts need it.

---

## Summary

- The **account model** (used by Ethereum) tracks a single running balance per address, updated in place — like a bank ledger — rather than a scattered pile of UTXOs.
- The same sequence of transactions produces identical final balances under either model; what differs is how the system tracks and validates value along the way.
- **UTXO's advantages**: natural parallelism (no shared mutable state to coordinate over) and a somewhat cleaner privacy story (funds scattered across many addresses rather than one running total).
- **The account model's advantage**: dramatically simpler support for smart contracts, which need persistent, arbitrarily-updatable storage that maps naturally onto account balances and storage slots.
- Other practical differences: UTXO transactions can grow large with many inputs but need no explicit nonce for replay protection; account-model transactions stay small but need one.
- **GoChain uses UTXO for its core ledger**, matching Bitcoin, because it best teaches this course's core ideas, its parallelism/privacy wins are essentially free for a base transfer layer, and the account model's main advantage (contract simplicity) doesn't need to apply chain-wide.
- Volume 9 will layer **account-like persistent storage** onto individual smart contracts, without replacing the UTXO-based core ledger — the right tool used only where it's actually needed.

---

## Exercises

### Easy

1. In your own words, explain the single biggest structural difference between how a UTXO-model transaction and an account-model transaction each describe "send 20 gochips to Bob."
2. List two advantages of the UTXO model and two advantages of the account model, using this chapter's terms (parallelism, privacy, contract simplicity, nonce/replay protection, state size).
3. Why can't a UTXO model transaction subtract "half a UTXO" the way an account model transaction can subtract "half a balance"? Refer back to Chapter 30 if needed.

### Medium

4. Redo the Section 2 side-by-side table with a *new* scenario of your own choosing: three people, a starting balance, and three transactions between them. Write out both columns fully (UTXO inputs/outputs and resulting set vs. account balances after each step) and confirm both models reach the same final balances.
5. A blockchain needs to process transactions from many different, unrelated senders as fast as possible. Explain, referencing Section 3's parallelism argument specifically, why a UTXO-based chain has an easier time validating a large batch of unrelated transactions concurrently than an account-based chain does.
6. Explain why the account model needs an explicit per-account "nonce" (transaction counter) to prevent the same transaction from being processed twice, while the UTXO model gets this protection "for free" from how spending works. What would go wrong, concretely, if an account-model chain forgot to include nonces?

### Hard

7. Section 7 previews that GoChain's smart contracts (Volume 9) will get account-like persistent storage layered on top of the UTXO core ledger. Sketch, on paper, how you imagine a contract call could be represented as a special kind of transaction in a UTXO-based chain — what would its inputs and outputs need to look like, and where would the contract's *other* state (its storage table) actually live, if not directly as a UTXO?
8. Research (general knowledge, no formal citation needed) one real practical consequence of Ethereum's account-model state growing large over time (for example, "state bloat" and its effect on running a full node), and one real practical consequence of Bitcoin's UTXO set growing large over time. Compare the two in 200-300 words: are they really the same problem in different clothes, or genuinely different challenges?
9. Argue, as persuasively as you can, for the *opposite* of this chapter's conclusion: that GoChain should have used the account model for its entire core ledger, contracts included, the way Ethereum does. Use at least two specific points from Sections 3-5, and be honest about which of GoChain's stated design goals (Section 6) your argument would sacrifice.
