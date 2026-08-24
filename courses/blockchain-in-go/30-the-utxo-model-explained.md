# Chapter 30: The UTXO Model, Explained

Bitcoin — and GoChain — does not store a running balance for each address anywhere. Instead, it tracks a pile of discrete, indivisible chunks of value called Unspent Transaction Outputs (UTXOs), and your "balance" is nothing more than the sum of the chunks that belong to you. This chapter builds that idea from a physical-cash analogy all the way to a full, worked numeric example tracing several transactions.

## Table of Contents

1. [Balances Are Not Stored — Bills and Coins Are](#1-balances-are-not-stored--bills-and-coins-are)
2. [What Is a UTXO?](#2-what-is-a-utxo)
3. [Spending Consumes, It Never Modifies](#3-spending-consumes-it-never-modifies)
4. [Change: Paying with a $20 for a $12 Item](#4-change-paying-with-a-20-for-a-12-item)
5. [Balance Is a Sum, Not a Number in a Cell](#5-balance-is-a-sum-not-a-number-in-a-cell)
6. [A Full Worked Example](#6-a-full-worked-example)
7. [Why Indivisible Chunks Instead of Just a Number?](#7-why-indivisible-chunks-instead-of-just-a-number)
8. [Where GoChain Fits](#8-where-gochain-fits)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Balances Are Not Stored — Bills and Coins Are

If you asked most people how Bitcoin tracks "how much money you have," they would probably guess there's a giant table somewhere: address on the left, balance on the right, updated every time money moves. That guess is wrong, and it's wrong in a way that matters.

Instead, think about the cash in your physical wallet right now. You don't have a little number stamped on your wallet that says "$47." You have a $20 bill, two $10 bills, a $5 bill, and two $1 bills — and if someone asks how much money you have, you have to actually add up the individual bills and coins to answer. Nobody updates a single "balance" field when you get paid or buy coffee; instead, bills and coins get handed over, pocketed, and occasionally broken into change.

This is *exactly* how Bitcoin — and GoChain — represents value. There is no field anywhere that says "Alice has 47 gochips." There is only a collection of discrete chunks of value, scattered across many past transactions, each one belonging to a specific address, waiting to be spent.

---

## 2. What Is a UTXO?

A **UTXO** (Unspent Transaction Output) is one of these chunks: an output created by some past transaction, that has not yet been used as the input to any later transaction. "Unspent" is the operative word — the moment it's used up as an input, it stops being a UTXO forever (more on this in Section 3).

Every UTXO carries exactly two pieces of information:

- **How much** value it holds (its amount).
- **Who** can spend it — in practice, whoever can produce a valid signature matching the address it's locked to.

The entire set of every UTXO that exists anywhere on the chain, across every address, is called the **UTXO set**. It is the single most important piece of "current state" a GoChain node needs to answer the question "can this transaction legally spend this?" — everything else (old, fully-spent transactions) is historical record-keeping that no longer affects what can be spent right now.

```
                     THE UTXO SET (a pile of unspent "bills")

  +-----------+   +-----------+   +-----------+   +-----------+
  | 20 gochips|   | 5 gochips |   | 10 gochips|   | 3 gochips |
  | -> Alice  |   | -> Bob    |   | -> Alice  |   | -> Carol  |
  +-----------+   +-----------+   +-----------+   +-----------+

  Alice's "balance" isn't stored anywhere -- it's the sum of every UTXO
  in this pile that's locked to her address: 20 + 10 = 30 gochips.
```

---

## 3. Spending Consumes, It Never Modifies

Here is the part that trips people up coming from a bank-account mental model: **you cannot partially spend a UTXO, and you cannot modify one in place.** A UTXO is all-or-nothing, exactly like a physical bill — you cannot hand a merchant "part of" a $20 bill. You can only hand over the whole bill and, if it's worth more than what you owe, receive change back.

When a transaction spends a UTXO, three things happen together, atomically:

1. The UTXO being spent is referenced as an **input** and is consumed *entirely* — it can never be referenced again by any other transaction, ever.
2. One or more brand-new UTXOs are created as **outputs** of the new transaction — these become the *only* new spendable chunks resulting from this transaction.
3. The old UTXO is removed from the UTXO set; the new output(s) are added to it.

Nothing about the old UTXO is ever edited, shrunk, or reused — it is destroyed completely and replaced by whatever new outputs the spending transaction creates. This is why double-spending (trying to use the same UTXO as an input in two different transactions) is something every node can detect just by checking whether that exact UTXO is still present in the UTXO set — Chapter 34 builds this check for real.

---

## 4. Change: Paying with a $20 for a $12 Item

If UTXOs can't be split, how does GoChain handle "I want to send exactly 12 gochips, but the only UTXO I have is worth 20"? The answer is the same one cash registers have used forever: **change**.

Imagine buying a $12 item with a $20 bill. The cashier doesn't tear the bill in half — they take the whole $20 bill (consuming it entirely) and hand you back an $8 bill (or equivalent coins) as change, along with your $12 item. Two new things exist after this transaction: the merchant now has your original $20 (in their till), and you now have a new $8 bill that didn't exist a moment ago.

A GoChain transaction that spends a 20-gochip UTXO to pay someone 12 gochips works identically:

```
   BEFORE:  Alice owns one UTXO: 20 gochips

   Transaction:
     Input:   Alice's 20-gochip UTXO           (consumed entirely)
     Output 1: 12 gochips -> Bob                 (the actual payment)
     Output 2: 8 gochips  -> Alice                (her own change)

   AFTER:   That 20-gochip UTXO no longer exists.
            Two new UTXOs exist: 12 -> Bob, 8 -> Alice
```

The "change" output isn't a special kind of transaction — it's just an ordinary output that happens to be locked back to the sender's own address. Wallets build this automatically (Chapter 32's `NewTransaction` does exactly this), and from the network's point of view there is nothing distinguishing a "change output" from a "payment output" — both are just UTXOs, waiting to be spent by whoever they're locked to.

---

## 5. Balance Is a Sum, Not a Number in a Cell

Once you accept that value lives entirely in a scattered pile of UTXOs, "checking a balance" has a precise, simple definition: **scan every UTXO currently in the UTXO set, add up the ones locked to a given address, and that sum is the balance.** There is no faster shortcut in the basic model — you must look at the actual pile of unspent chunks, not read a single stored number (Volume 8 builds an *index* that makes this scan fast without changing what it fundamentally means).

```
balance(address) = sum of Value for every UTXO in the UTXO set
                    whose PubKeyHash matches address
```

This has an important, sometimes-surprising consequence: an address's balance can be spread across many separate UTXOs of different sizes, just like your physical wallet might hold one $20, three $5s, and a $1. Sending a large payment might require a wallet to gather up *several* UTXOs as inputs to cover the amount — exactly the "selecting enough UTXOs to cover an amount" logic Chapter 32 implements.

---

## 6. A Full Worked Example

Let's trace the entire UTXO set through four transactions, step by step, so the mechanics from Sections 2-5 become completely concrete. Assume a coinbase-style reward transaction has already given Alice her starting funds (Chapter 37 covers coinbase transactions properly; for now, just accept that Alice starts with one 50-gochip UTXO from mining a block).

**Starting state — UTXO set after Alice mines a block:**

```
UTXO #1:  50 gochips -> Alice   (created by the mining reward)
```

Alice's balance: 50. Bob's balance: 0. Carol's balance: 0.

**Transaction 1 — Alice sends Bob 20 gochips:**

```
  Input:    UTXO #1 (50 gochips -> Alice)        -- consumed entirely
  Output A: 20 gochips -> Bob                     -- the payment
  Output B: 30 gochips -> Alice                   -- change
```

UTXO set after Transaction 1:

```
UTXO #1:  CONSUMED (no longer exists)
UTXO #2:  20 gochips -> Bob      (Transaction 1, Output A)
UTXO #3:  30 gochips -> Alice    (Transaction 1, Output B, her change)
```

Alice's balance: 30. Bob's balance: 20. Carol's balance: 0.

**Transaction 2 — Bob sends Carol 5 gochips:**

Bob only has one UTXO (#2, worth 20), so it must be spent in full even though he only wants to send 5:

```
  Input:    UTXO #2 (20 gochips -> Bob)           -- consumed entirely
  Output C: 5 gochips -> Carol                     -- the payment
  Output D: 15 gochips -> Bob                       -- Bob's change
```

UTXO set after Transaction 2:

```
UTXO #2:  CONSUMED
UTXO #3:  30 gochips -> Alice    (still unspent, untouched)
UTXO #4:  5 gochips  -> Carol    (Transaction 2, Output C)
UTXO #5:  15 gochips -> Bob      (Transaction 2, Output D, his change)
```

Alice's balance: 30. Bob's balance: 15. Carol's balance: 5.

**Transaction 3 — Alice sends Carol 30 gochips, exactly using up her one UTXO:**

```
  Input:    UTXO #3 (30 gochips -> Alice)         -- consumed entirely
  Output E: 30 gochips -> Carol                    -- the payment, exact amount, no change needed
```

Notice: when an input's value exactly matches the amount being sent, there is no change output at all — a transaction can have as few as one output.

UTXO set after Transaction 3:

```
UTXO #3:  CONSUMED
UTXO #4:  5 gochips  -> Carol    (still unspent)
UTXO #5:  15 gochips -> Bob      (still unspent)
UTXO #6:  30 gochips -> Carol    (Transaction 3, Output E)
```

Alice's balance: 0. Bob's balance: 15. Carol's balance: 35.

**Transaction 4 — Carol sends Alice 8 gochips, needing to combine two UTXOs:**

Carol wants to send 8 gochips, but neither of her individual UTXOs (5 and 30) alone determines the smallest possible payment — a real wallet would prefer spending just UTXO #6 (30) since it alone covers 8, leaving one input and simpler change. Suppose Carol's wallet selects UTXO #4 (5 gochips) *and* UTXO #6 (30 gochips), for a combined 35, to illustrate a transaction spending **multiple inputs** at once:

```
  Input 1:  UTXO #4 (5 gochips -> Carol)          -- consumed entirely
  Input 2:  UTXO #6 (30 gochips -> Carol)         -- consumed entirely
  Output F: 8 gochips -> Alice                     -- the payment
  Output G: 27 gochips -> Carol                     -- Carol's change (5+30-8)
```

UTXO set after Transaction 4:

```
UTXO #4:  CONSUMED
UTXO #6:  CONSUMED
UTXO #5:  15 gochips -> Bob      (still unspent, untouched this whole time)
UTXO #7:  8 gochips  -> Alice    (Transaction 4, Output F)
UTXO #8:  27 gochips -> Carol    (Transaction 4, Output G)
```

Final balances: Alice = 8, Bob = 15, Carol = 27. Check the arithmetic: 8 + 15 + 27 = 50, exactly the amount Alice originally mined — no value was created or destroyed anywhere along the way, only moved and occasionally split into change. This conservation of total value is a useful sanity check you'll rely on again when writing tests in Chapter 32.

---

## 7. Why Indivisible Chunks Instead of Just a Number?

It's worth pausing on *why* this seemingly roundabout system (indivisible chunks, change outputs, sometimes combining several inputs) is worth the complexity compared to simply storing one balance number per address, the way a bank does. Two properties fall out of it almost for free:

- **Natural parallelism.** Two transactions that spend *different* UTXOs don't touch each other's data at all — they can be validated and even processed concurrently, since there's no shared "balance" field they'd both need to update in a coordinated way. This becomes a real, exploitable property later (Volume 9's contract design and Volume 11's scaling discussions both return to it).
- **A cleaner privacy story.** Because value lives in many separate, unlinked chunks rather than one running total per address, it's harder for an outside observer to immediately see "this address has exactly this much" the way they could from a single public ledger entry — though, to be clear, Bitcoin-style UTXO chains are still fully public and traceable; this is a modest privacy improvement, not anonymity.

Chapter 31 puts these two properties side by side against the account model's own strengths, and explains exactly why GoChain adopts UTXO for its core ledger.

---

## 8. Where GoChain Fits

Every mechanic in this chapter maps directly onto Go types you'll write in Chapter 32: a UTXO is represented by a `TxOutput` sitting inside some transaction's `Outputs` slice that no later transaction's `Inputs` has referenced yet. The `UTXOSet` helper type you'll build scans the chain to answer exactly the two questions this chapter relied on informally: "what is this address's balance?" and "which specific UTXOs should a new transaction spend to cover an amount?" Chapter 34's mempool builds directly on the "a UTXO can only ever be consumed once" rule from Section 3 to detect double-spends before they're ever mined.

---

## Summary

- GoChain does not store per-address balances anywhere; it tracks a scattered set of **UTXOs** (Unspent Transaction Outputs), like bills and coins in a physical wallet.
- A UTXO is a discrete, indivisible chunk of value created as an output of some transaction, that no later transaction has yet consumed as an input.
- Spending a UTXO **consumes it entirely** — it is never partially spent or modified, only destroyed and replaced by new outputs.
- **Change** — an output paying the sender back their own leftover value — is how transactions handle spending a UTXO worth more than the amount being sent, just like a cashier giving change for a large bill.
- **Balance is computed, not stored**: it's the sum of every UTXO currently in the UTXO set that belongs to a given address.
- The worked example traced four transactions, showing UTXOs being created, consumed, combined (multiple inputs), and split (change outputs), while total value was conserved throughout.
- UTXO's key strengths — natural parallelism and a cleaner (though not anonymous) privacy story — set up Chapter 31's comparison against the account model.

---

## Exercises

### Easy

1. In your own words, explain why you cannot "spend half a UTXO" the same way you cannot hand a cashier half of a $10 bill.
2. Given a UTXO set of `{15 -> Alice, 4 -> Alice, 30 -> Bob}`, what is Alice's balance? What is the *minimum* number of UTXOs Alice's wallet would need to combine as inputs to pay someone exactly 19 gochips, and which ones?
3. In Transaction 3 of the worked example (Section 6), no change output was created. Explain precisely why, and describe the general rule for when a transaction needs a change output and when it doesn't.

### Medium

4. Continue the worked example in Section 6 with a fifth transaction: Bob (currently holding 15 gochips in UTXO #5) sends Alice 15 gochips exactly. Write out the transaction's inputs and outputs, the resulting UTXO set, and the final balances for all three parties. Verify total value is still conserved.
5. A friend claims "the UTXO model is basically the same as the account model, just with extra steps." Write a 150-200 word rebuttal (or defense, if you disagree) using specific mechanics from this chapter — spending, change, and the UTXO set — to argue your position.
6. Design (on paper) a UTXO-selection strategy for a wallet that needs to send a specific amount: given a list of available UTXOs and a target amount, describe step by step how you'd decide which UTXOs to use as inputs, and what tradeoffs your strategy makes (for example: fewest UTXOs used vs. smallest resulting change output vs. avoiding very small "dust" UTXOs).

### Hard

7. Suppose a malicious node tried to convince the network that a UTXO which was already spent in Transaction 2 (Section 6) was still spendable, by resubmitting a transaction referencing it. Walk through, mechanically, how an honest node checking the current UTXO set would catch and reject this, referencing the specific UTXO by number from the worked example.
8. Research how Bitcoin actually selects which UTXOs to spend when building a transaction (general knowledge is fine; this is a well-known problem called "coin selection"). Summarize, in 200-300 words, at least two real coin-selection strategies real wallets use and what problem each one is specifically trying to avoid (hint: transaction size/fees, privacy, or leaving many small "dust" UTXOs behind).
9. Prove, with a general argument (not just the one worked example), that the UTXO model can never create or destroy total value on its own — that is, argue why the sum of a transaction's output values can never legitimately exceed the sum of its input values, and explain what rule a node must enforce to guarantee this (this rule becomes fully precise once fees are introduced in Chapter 35).
