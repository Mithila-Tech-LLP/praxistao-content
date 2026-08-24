# Chapter 59: What Are Smart Contracts?

Every transaction GoChain has processed so far follows the same shape: Alice signs a message saying "move this value from me to Bob," and every node checks the signature and moves on. A smart contract asks a different question: what if the *rule* for moving value is not "whoever holds the right private key," but a small program — one that every node runs, and every node agrees on the result of? This chapter builds that idea from a concrete example, with no code yet, so the virtual machine we start building in Chapter 60 has an obvious job to do.

## Table of Contents

1. [A Motivating Problem: Escrow Without a Middleman](#1-a-motivating-problem-escrow-without-a-middleman)
2. [What a Smart Contract Actually Is](#2-what-a-smart-contract-actually-is)
3. [Why Deterministic Execution Is the Whole Point](#3-why-deterministic-execution-is-the-whole-point)
4. [Contracts vs. Plain Transactions](#4-contracts-vs-plain-transactions)
5. [Walking Through the Escrow Contract, Step by Step](#5-walking-through-the-escrow-contract-step-by-step)
6. [What Smart Contracts Are Not](#6-what-smart-contracts-are-not)
7. [Where This Fits Into GoChain](#7-where-this-fits-into-gochain)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. A Motivating Problem: Escrow Without a Middleman

Imagine Priya wants to buy a used laptop from a stranger, Devan, whom she has never met and has no reason to trust. Priya is nervous about sending money before receiving the laptop; Devan is nervous about shipping the laptop before receiving money. This is the oldest problem in commerce between strangers, and the traditional fix is **escrow**: both parties send their side of the deal to a trusted third party (an escrow agent), who holds the money until Priya confirms the laptop arrived, then releases it to Devan.

```
TRADITIONAL ESCROW

  Priya                Escrow Agent               Devan
    |                        |                       |
    |-- sends 200 gochips -->|                       |
    |                        |<---- ships laptop ----|
    |-- "I received it" ---->|                       |
    |                        |---- releases 200 ---->|
    |                        |     gochips to Devan   |

  Works, but: Priya and Devan must both trust the agent not to run
  off with the money, collude with one side, or simply make mistakes.
```

This works, but it reintroduces exactly the kind of trusted middleman that a decentralized system is supposed to make unnecessary. The escrow agent could be dishonest. The escrow agent could be slow, expensive, or simply unavailable in the country either party lives in. Every dollar of trust placed in that one agent is a dollar of risk.

A **smart contract** solves this without any person acting as the middleman at all: Priya sends her 200 gochips to a small program instead of to a person. That program's rules — written once, in code, and visible to both parties before they agree to use it — say exactly this: *hold the funds; release them to Devan only when Priya submits a message confirming delivery; refund Priya if a deadline passes with no confirmation.* Every node on the GoChain network runs that exact same program on the exact same inputs and reaches the exact same conclusion about whether the funds should move. No person is the referee. The code is the referee.

```
SMART CONTRACT ESCROW

  Priya                 Escrow Contract              Devan
    |                    (code, not a person)            |
    |-- sends 200 gochips -->|                           |
    |     (contract holds funds)                          |
    |                        |<------ ships laptop -------|
    |-- "confirm delivery" ->|                           |
    |               every node runs the SAME program      |
    |               on the SAME inputs, reaches the        |
    |               SAME conclusion: release funds         |
    |                        |----- releases 200 -------->|
    |                        |      gochips to Devan       |
```

The word **escrow** here means: a neutral holding arrangement where funds sit with a third party until agreed-upon conditions are met. In the smart-contract version, "the third party" is not a company or a person — it is code, running identically everywhere.

---

## 2. What a Smart Contract Actually Is

A **smart contract** is a program that is stored on the blockchain and executed by every node that processes the transaction invoking it, producing a result every honest node independently agrees on. That is a precise definition, so let's take it apart piece by piece:

- **"Stored on the blockchain"** — the contract's code itself lives in the chain's data, the same way a transaction does. Anyone can read it. Nobody (not even the person who wrote it) can secretly change it after it is deployed, because doing so would require rewriting history, which Volume 4's proof of work already made prohibitively expensive.
- **"Executed by every node"** — when a transaction calls a contract, every node that validates that transaction does not just check a signature; it actually *runs the contract's code*, starting from the same state, using the same inputs.
- **"A result every honest node independently agrees on"** — because every node ran the identical program on identical inputs, they all arrive at the identical answer, with no need to ask each other or trust any single node's report of what happened.

The escrow example from Section 1 is a smart contract in exactly this sense: its code is what everyone on the network agreed to trust *in advance*, not a person or company they have to trust *afterward*.

It helps to name the two roles in that example precisely, because we will reuse them for the rest of this volume:

- The **contract** is the stored program plus whatever data it is currently holding (its balance, and — starting in Chapter 66 — its own private storage).
- A **transaction that invokes the contract** is how someone triggers the contract's code to run, the same way calling a function triggers code to run — except here, "calling the function" means submitting a transaction that every node on the network will independently execute.

---

## 3. Why Deterministic Execution Is the Whole Point

**Deterministic** means: given the same input, a process always produces the same output, every single time, on every single machine. This property is not a nice-to-have for smart contracts — it is the entire reason the idea works at all. If ten thousand nodes each ran the escrow contract and reached ten thousand *different* conclusions about whether Devan should get paid, the network would have no way to agree on anything, and the whole point of removing the trusted middleman would collapse right back into needing one (whoever's answer you happen to trust).

This is why smart contracts cannot be arbitrary programs doing arbitrary things:

- **No real time.** A contract cannot ask "what time is it right now?" and get a consistent answer, because network delay means different nodes execute the same transaction at slightly different real-world moments. (Contracts that need a notion of time use the **block height** or a timestamp embedded in the block itself — a value every node agrees on, because it is part of the data being validated, not read live from a clock.)
- **No random numbers from the operating system.** `rand.Int()` on Node A and `rand.Int()` on Node B produce different numbers. Any "randomness" a contract uses has to come from something every node can compute identically (or, in advanced designs, from a verifiable randomness scheme — well outside this course's scope).
- **No network calls, no file reads, no anything that could return a different answer on a different machine.** A contract that asked "what does this external website say?" would immediately break determinism, because two nodes checking that website at different moments could see different content.

This is precisely why Chapter 60 does not let contracts run as ordinary compiled Go, Python, or JavaScript programs. Ordinary programs can do all of the non-deterministic things above without a second thought. GoChain instead runs contracts on a small, deliberately restricted virtual machine — one that is simply incapable of expressing "ask the operating system for the time" or "make an HTTP request," because those operations do not exist in its instruction set at all. Restricting *what a program can even say* is a much stronger guarantee than trusting a programmer to avoid non-determinism by discipline alone.

```
ORDINARY PROGRAM                        SMART CONTRACT (restricted VM)

can call: time.Now()                    CANNOT call anything like it —
can call: http.Get(url)                 those operations simply do not
can call: rand.Int()                    exist as instructions the VM
can read/write any file                 understands. The instruction set
                                          only contains safe, deterministic
  --> non-determinism is POSSIBLE        operations (push, add, compare, ...)

                                          --> non-determinism is IMPOSSIBLE
                                              by construction, not by promise
```

---

## 4. Contracts vs. Plain Transactions

By this point in the course, GoChain already has a working notion of a transaction (Volume 5): a signed instruction moving value from a sender's UTXOs to a recipient's new outputs. It is worth being precise about how a contract call relates to that, because Chapter 63 will show the two are more alike than they first appear.

| | Plain transaction | Smart contract call |
|---|---|---|
| What decides if it's valid | A signature matches the claimed owner | A small program runs and produces a pass/fail (or a computed result) |
| What it can express | "Move value from A to B" | Arbitrary logic: conditions, state, multi-step rules |
| Where the logic lives | Fixed, built into `core` itself | Stored on-chain, different per contract |
| Example | Alice pays Bob 10 gochips | "Release funds only if Priya confirms delivery" |

The key realization Chapter 63 will make concrete: a plain transaction's "check the signature" rule is *itself* just a very small, very common program — one that always does the same thing (verify one signature against one key). GoChain's virtual machine will be general enough to express both a plain payment and an escrow contract as programs on the same machine, which is what "unifying transactions and contracts under one execution model" (Chapter 63's subject) actually means. A smart contract, in other words, is not a separate universe from a transaction — it is what you get when you let the "is this spend allowed?" question be answered by an arbitrary small program instead of only ever "does the signature match."

---

## 5. Walking Through the Escrow Contract, Step by Step

Let's trace the escrow example all the way through, in plain language, before any GoChain-specific mechanics appear in later chapters. Assume the contract's rules, agreed on by both parties before any money moves, are:

1. Priya deposits 200 gochips into the contract.
2. If Priya later submits a "confirm delivery" message, signed with her private key, the contract releases the 200 gochips to Devan.
3. If block height passes a deadline with no confirmation from Priya, the contract instead refunds the 200 gochips to Priya.

```
STATE: contract holds 0 gochips, waiting

Priya's deposit transaction is mined
        |
        v
STATE: contract holds 200 gochips, waiting for confirmation or deadline
        |
        |------------------------+
        |                        |
   Priya sends signed        deadline block height
   "confirm delivery"        is reached with no
        |                    confirmation
        v                        v
  contract verifies          contract releases
  Priya's signature           200 gochips back
  matches the deposit         to Priya (refund)
  address, then releases
  200 gochips to Devan
```

Notice what every node does here: it does not "ask Priya" or "ask Devan" who is right. It looks at the contract's stored code and the contract's stored state (currently: 200 gochips held, no confirmation yet, deadline at block height N), applies the exact same rules, and reaches the exact same conclusion. If a dishonest node claimed the funds went to the wrong party, every honest node would simply disagree with it and reject its version of events — precisely the "agreement" property Chapter 23 introduced for consensus in general, now applied to contract execution specifically.

---

## 6. What Smart Contracts Are Not

A few honest myth-corrections, because "smart contract" is a term that invites overclaiming:

- **They are not "smart" in any AI sense.** A smart contract is a small, deterministic, usually quite simple program — often closer to a vending machine's internal logic than to anything resembling intelligence. "Smart" here really just means "self-executing, based on code rather than on a person's judgment call."
- **They are not legally binding contracts by themselves.** A smart contract enforces whatever its code says, exactly, even if the code has a bug that does something nobody intended. It does not know or care what the parties *meant* — only what the code *says*. This is precisely why Chapter 67 (reentrancy) matters so much: bugs in contract code are not "fixable after the fact" the way a bug in a normal web app is, because the contract's behavior on past invocations is already permanently recorded on the chain.
- **They do not remove all trust.** Both parties still have to trust that the contract's code actually does what it claims to do. In the escrow example, Priya and Devan should both read (or have someone trustworthy read) the escrow contract's code before agreeing to use it. What smart contracts remove is the need to trust a specific *person or company* to behave honestly after the fact — the trust shifts to "can I read and verify this code," which is a fundamentally more checkable kind of trust.

---

## 7. Where This Fits Into GoChain

This volume builds, in order, exactly what is needed to make the escrow example real:

```
gochain/
├── crypto/       (done)   Hash(), Sign(), Verify() — the escrow contract's
│                          signature check reuses Verify() directly.
├── core/         (done)   Transaction, TxOutput — a contract call is a
│                          transaction whose output carries a program
│                          instead of (or alongside) a simple address.
├── consensus/    (done)   Proof of work — mining a block that contains a
│                          contract call is no different from mining any
│                          other block; the VM runs during validation.
├── vm/           (THIS VOLUME) <- the restricted, deterministic machine
│   │                              that runs contract code.
│   ├── Ch. 60-62: the machine itself (stack, opcodes, execution loop)
│   ├── Ch. 63:    locking/unlocking scripts — how a TxOutput becomes
│   │              "spendable only by running a small program"
│   └── Ch. 64:    gas — what stops a buggy or malicious contract from
│                  running forever
├── wallet/       (done)
├── network/      (done)
└── storage/      (done)   Chapter 66 (next half of this volume) adds
                            contract-local storage on top of it.
```

Chapter 60 starts exactly where this chapter leaves off: if a contract must be a small, deterministic program, what does the *machine* that runs it actually look like? The answer — a stack-based virtual machine — is the subject of the rest of this volume.

---

## Summary

- A smart contract is a program stored on the blockchain, executed identically by every node, so its results require trusting code, not a person or company.
- Escrow motivates the idea concretely: a contract can hold funds and release them based on a rule (confirmation, or a deadline) without any human acting as the middleman.
- Deterministic execution — same input, same output, everywhere, always — is not optional; it is the entire reason every node can agree on a contract's result without asking each other.
- Contracts avoid non-determinism by construction: the machine that runs them simply has no instruction for reading the system clock, making network calls, or generating true randomness.
- A plain transaction and a smart contract call are more alike than they look: both are "a small program decides whether this spend is allowed," just with different programs (Chapter 63 makes this precise).
- Smart contracts are not artificially intelligent, are not automatically legally binding, and do not eliminate all trust — they shift trust from "trust this person" to "trust this code," which is independently verifiable.
- This volume builds, in order: the virtual machine itself (Ch. 60-62), locking/unlocking scripts that unify transactions and contracts (Ch. 63), and gas metering to stop runaway execution (Ch. 64) — with a token contract and contract storage following in Ch. 65-69.

---

## Exercises

### Easy

1. In your own words, explain why "every node runs the same program and gets the same result" is what lets a smart contract replace a trusted escrow agent.
2. List three operations an ordinary program can do that a deterministic smart contract cannot, and explain why each one would break agreement across nodes.
3. Define, in one or two sentences each: smart contract, deterministic, escrow.

### Medium

4. Rewrite the escrow example from Section 5 as a numbered list of rules for a *different* everyday situation (for example: a rental deposit, a bet between two friends, a crowdfunding pledge that only charges backers if a goal is met). Identify what "state" the contract needs to remember between steps.
5. A friend claims "smart contracts remove the need for trust entirely." Using the distinctions from Section 6, explain what is right and what is misleading about that claim.
6. Explain why a contract cannot safely use `time.Now()` but can safely use the block height a transaction was included in. What property does block height have that wall-clock time does not?

### Hard

7. Suppose the escrow contract's code has a subtle bug: the deadline check uses `>` instead of `>=`, so a refund one block later than intended goes to Devan instead of Priya. Once this transaction is mined into a block, can the mistake be corrected by "just fixing the code"? Explain why or why not, referencing how blocks and hashes work (Volumes 3-4).
8. Research one real, historical smart-contract exploit (for example, the 2016 DAO hack, previewed for you in Chapter 67, or any other well-documented incident). Summarize, in your own words, what the contract's code allowed that its authors did not intend, and connect it back to the "code is the referee, exactly as written" idea from Section 6.
9. Design, on paper, a simple smart contract for a scenario of your choosing that needs at least two participants and at least one condition based on time (block height) and one condition based on a signed confirmation message. Diagram its possible states and transitions the way Section 5 diagrammed the escrow contract.
