# Chapter 82: Ethereum Architecture, Deep Dive

Chapter 81 looked at Bitcoin, the chain whose UTXO model GoChain's core ledger directly matches. This chapter looks at the other chain that shaped this course's design at every turn, but from the opposite direction: Ethereum kept the "chain of blocks secured by proof of work" idea from Bitcoin, then rebuilt almost everything else around a different core decision — track running balances per account instead of a set of spendable coins — specifically so that general-purpose programs (smart contracts) would have somewhere sensible to keep their own state. Volume 9's `vm.VM` and Volume 8's simplified trie are both, in a real sense, small, teachable descendants of Ethereum's design. This chapter opens up the real thing: the account model, the real EVM's opcode categories, gas pricing under actual network congestion, the real state trie, and "The Merge" — the day in 2022 Ethereum swapped its entire consensus engine out from under a live, running, trillion-dollar-adjacent system without stopping it.

## Table of Contents

1. [Why Ethereum, After Bitcoin](#1-why-ethereum-after-bitcoin)
2. [The Account Model, For Real](#2-the-account-model-for-real)
3. [The Real Ethereum State Trie](#3-the-real-ethereum-state-trie)
4. [The EVM: Real Opcode Categories](#4-the-evm-real-opcode-categories)
5. [Gas Pricing in a Live, Congested Network](#5-gas-pricing-in-a-live-congested-network)
6. [Contract Accounts and External Accounts](#6-contract-accounts-and-external-accounts)
7. [The Merge: Swapping Consensus in a Live System](#7-the-merge-swapping-consensus-in-a-live-system)
8. [Ethereum vs. GoChain: Side-by-Side](#8-ethereum-vs-gochain-side-by-side)
9. [What GoChain Simplified, and Why That Was the Right Call](#9-what-gochain-simplified-and-why-that-was-the-right-call)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Ethereum, After Bitcoin

Think of Bitcoin as a single, extremely reliable vending machine: it does one job — move coins from one owner to another under strict rules — and it does that job with an almost obsessive commitment to never being surprised. Ethereum asked a different question: what if the machine could run *any* small program you handed it, as long as every other machine in the network agreed to run the exact same program and got the exact same answer? That single change — from "verify a fixed set of spend rules" to "execute arbitrary, deterministic code" — is the entire reason Ethereum's architecture looks different from Bitcoin's at nearly every layer, even though both are, underneath, a chain of cryptographically linked blocks secured by a consensus engine.

You already own a small version of exactly that idea. Volume 9's `vm.VM` runs arbitrary (bounded, gas-metered) programs. Chapter 65's token contract keeps its own persistent state via `OpSLoad`/`OpSStore`. Chapter 57's simplified trie commits to a whole set of key-value data with one root hash. None of this is a coincidence — Volume 9 and Chapter 57 were built with Ethereum's architecture specifically in mind, at a scale one person can implement and fully understand. This chapter shows you the full-scale versions.

As with Chapter 81, the goal here is architectural durability, not a snapshot of any specific client version. Ethereum's account model, its state trie's overall shape, and its EVM's core opcode categories have been stable, load-bearing design decisions since the network's earliest years — the details this chapter describes are safe to treat as durable, even as specific gas costs, opcode additions, and client implementations continue to evolve.

---

## 2. The Account Model, For Real

Chapter 31 already introduced you to the account model conceptually and explained why GoChain chose UTXO for its base ledger anyway. Here is the real thing Chapter 31 was describing.

Picture the difference this way: Bitcoin's UTXO model is a wallet stuffed with individual bills and coins — your "balance" is just whatever you'd get by counting them up. Ethereum's account model is a bank's ledger book — one line per account, showing a single running number that goes up and down as transactions are applied, with the bank (every node, independently) keeping that ledger perfectly in sync.

Concretely, Ethereum's **world state** is a mapping from every address to an **account**, and every account carries exactly four fields:

```
  AN ETHEREUM ACCOUNT
  +----------------------------------------------------------+
  |  nonce        -- for an externally-owned account, a       |
  |                  count of transactions sent (prevents      |
  |                  replaying the same signed transaction      |
  |                  twice); for a contract account, a count    |
  |                  of contracts it has created                |
  |  balance       -- the account's current balance, in wei     |
  |                  (the smallest indivisible unit of ether,    |
  |                  playing the exact role GoChain's "gochip"   |
  |                  plays for gochips)                          |
  |  storageRoot   -- root hash of this account's OWN storage    |
  |                  trie (empty for a plain account, populated  |
  |                  for a contract — see Section 3)             |
  |  codeHash      -- hash of this account's code (empty for a   |
  |                  plain account; the contract's bytecode      |
  |                  for a contract account)                     |
  +----------------------------------------------------------+
```

A **nonce**, in this context, has nothing to do with the mining nonce from Chapter 24 — it is simply a strictly increasing counter, one per account, whose entire job is making sure a signed transaction can only ever be applied once. This is a direct, real-world illustration of a problem GoChain's UTXO model sidesteps by construction: because a UTXO is consumed entirely and can never be spent twice (Chapter 34's double-spend check enforces this directly against the set of outputs), there is no equivalent "replay" risk to guard against with a counter. The account model reintroduces that risk the moment it stops tracking discrete, single-use outputs — and the nonce is the account model's answer to it.

```
  UTXO MODEL (GoChain, Bitcoin)          ACCOUNT MODEL (Ethereum)
  --------------------------------       --------------------------------
  "balance" = sum of UTXOs owned         balance = one stored number
  spend = consume UTXOs, create new      spend = decrement sender's
          ones                                   number, increment
                                                  receiver's
  double-spend prevented by:             double-spend prevented by:
    an output can only be spent            a transaction can only be
    once, ever (Ch. 34)                    applied once, tracked by
                                            the sender's nonce
  naturally parallel: unrelated          naturally sequential per
  spends never conflict                   account: nonce N+1 must
                                           follow nonce N
```

That last row matters in practice: because every transaction from a given account must apply its nonce in strict order, two transactions from the *same* account can't be processed out of order or in parallel — a real, structural cost of the account model that Chapter 31 gestured at when it called UTXO "naturally parallel."

Sending value under the account model is not "select UTXOs, build outputs, sign" (Chapter 32) — it is simply "look up the sender's balance, check it covers the amount plus gas costs, subtract from sender, add to recipient, increment sender's nonce." GoChain's transactions do not need this exact mechanism, because Chapter 32-34's UTXO machinery already solves the same underlying problem (don't let value be spent twice) a different way — but the account model is precisely why Ethereum's smart contracts got to have persistent state as naturally as a contract account's own row in the ledger, which brings us to the trie that actually stores all of this.

Here is that same balance transfer, worked through concretely, next to GoChain's UTXO equivalent from Chapter 32, so the "same job, different bookkeeping" framing is visible at the level of actual numbers rather than just prose:

```
  ACCOUNT MODEL: Alice sends 5 (ether-equivalent units) to Bob

  BEFORE:  Alice { balance: 20, nonce: 7 }   Bob { balance: 3, nonce: 2 }

  Alice signs a transaction: { to: Bob, value: 5, nonce: 7 }
  Every node applies it identically:
    - check nonce 7 matches Alice's stored nonce exactly (else reject)
    - check Alice's balance (20) >= 5
    - Alice.balance -= 5   -->  15
    - Bob.balance   += 5   -->  8
    - Alice.nonce   += 1   -->  8

  AFTER:   Alice { balance: 15, nonce: 8 }   Bob { balance: 8, nonce: 2 }


  UTXO MODEL: Alice sends 5 gochips to Bob (Ch. 32-34)

  BEFORE:  Alice owns one UTXO worth 20 gochips

  Alice builds a core.Transaction:
    - input:  spend the 20-gochip UTXO entirely
    - output: 5 gochips locked to Bob's address
    - output: 15 gochips locked back to Alice's OWN address (change)
  Every node applies it identically:
    - check the referenced UTXO exists and is unspent (Ch. 34)
    - mark the 20-gochip UTXO as spent, forever
    - add two brand-new UTXOs to the UTXO set: 5 (Bob), 15 (Alice)

  AFTER:   Alice owns a new 15-gochip UTXO. Bob owns a new 5-gochip UTXO.
           The original 20-gochip UTXO no longer exists at all.
```

Notice what each model had to invent to make this safe: the account model needed the nonce, specifically because "Alice.balance -= 5" is a mutation that could otherwise be replayed; the UTXO model needed a change output, specifically because a UTXO cannot be partially spent — it is entirely consumed and replaced, never edited in place. Neither invention is optional under its own model; each is the direct consequence of the core data-model choice made at the very top of this section.

---

## 3. The Real Ethereum State Trie

Chapter 57 built a simplified Merkle-Patricia trie and used it to commit to GoChain's UTXO set with one root hash per block. That chapter was, quite deliberately, teaching you a smaller version of exactly the data structure this section now shows you in full.

The analogy from Chapter 57 still holds without modification: a **trie** (short for re*trie*val tree, though pronounced "try" to avoid confusion with "tree") is a tree where the *path you walk down* to reach a value is determined by the key itself, letter by letter (or nibble by nibble, for Ethereum's hex-keyed trie) — like a filing cabinet where a document's folder is determined entirely by spelling out its name one letter at a time, so two documents with similar names share most of their folder path. A **Merkle-Patricia trie** layers Chapter 10's Merkle-tree tamper-evidence (every node's hash commits to everything beneath it) on top of a **Patricia trie**'s path-compression trick (a chain of single-child nodes collapses into one combined "extension" node, so keys that share a long common prefix don't waste space on redundant single-branch levels).

Ethereum does not use one trie — it uses several, nested:

```
  ETHEREUM'S NESTED TRIES (one block)

  BLOCK HEADER
  +------------------------------------------------------+
  |  stateRoot        --> STATE TRIE                       |
  |                        key: address                     |
  |                        value: account (nonce, balance,  |
  |                                storageRoot, codeHash)    |
  |                                                            |
  |                        for each account with storage:     |
  |                        storageRoot --> STORAGE TRIE         |
  |                                        key: storage slot      |
  |                                        value: stored word     |
  |                                        (this is the real,      |
  |                                        full version of your    |
  |                                        Ch. 66 contract storage)|
  |                                                            |
  |  transactionsRoot --> TRANSACTIONS TRIE                 |
  |                        (one Merkle-Patricia trie of      |
  |                        this block's transactions)         |
  |                                                            |
  |  receiptsRoot     --> RECEIPTS TRIE                      |
  |                        (one entry per transaction:        |
  |                        did it succeed, how much gas did   |
  |                        it use, what events did it log)    |
  +------------------------------------------------------+
```

The **state trie** is the big one: it is, in effect, Ethereum's entire world state — every account that has ever held a nonzero balance, sent a transaction, or been the target of contract code, keyed by address — committed to with a single 32-byte **state root** stored directly in the block header. Every contract account additionally owns its *own* storage trie (reachable via its `storageRoot` field), which is the real, full version of the per-contract storage slots Chapter 66's `ContractStore` gave GoChain's contracts. A **transactions trie** and a **receipts trie** exist per block for the same Chapter 10 reason your `core.Block` stores a `MerkleRoot` — so any single transaction (or its result) can be proven included, or excluded, without needing the whole block.

The single most important practical consequence of this design: because the state root sits in the block header, any two nodes that computed the same state root after processing the same blocks are *provably* in complete agreement about every account's balance, every contract's storage, everything — without either node needing to send the other its entire multi-gigabyte state to check. This is the exact same guarantee Chapter 57's simplified trie demonstrated at UTXO-set scale; Ethereum simply runs the identical idea across every account and every contract's storage, continuously, for a state that has grown into the hundreds of gigabytes.

---

## 4. The EVM: Real Opcode Categories

Chapter 60 already told you, honestly, that "a stack-based VM" is the conceptual ancestor of both Bitcoin Script and the EVM. Chapter 81 showed you Bitcoin Script's real, deliberately narrow version. Here is the EVM's real, deliberately *broad* version — the same execution model (a stack, a program counter, an instruction dispatch loop — exactly what `vm.VM.Execute()` implements), stretched to support arbitrary computation.

```
  EVM OPCODE CATEGORIES                          GOCHAIN vm.VM EQUIVALENT
  ------------------------------------------      ------------------------------
  Arithmetic & comparison                          OpAdd, OpSub, OpEqual,
    ADD, SUB, MUL, DIV, MOD, LT, GT, EQ            OpGreaterThan (Ch. 61)

  Bitwise & stack manipulation                     OpDup, OpPop (Ch. 61);
    AND, OR, XOR, NOT, DUP1-16, SWAP1-16,           no direct GoChain
    POP, PUSH0-32                                  equivalent for bitwise ops

  Memory                                           no direct equivalent --
    MLOAD, MSTORE, MSTORE8, MSIZE                  GoChain's VM has no
                                                    separate scratch memory
                                                    region, only the stack
                                                    and contract storage

  Storage                                          OpSLoad, OpSStore (Ch. 66)
    SLOAD, SSTORE

  Control flow                                     OpJump, OpJumpIfFalse
    JUMP, JUMPI, JUMPDEST, PC, STOP                (Ch. 61)

  Environment / context info                       no direct equivalent --
    CALLER, CALLVALUE, ADDRESS, ORIGIN,             GoChain's opcode table
    GASPRICE, BLOCKHASH, TIMESTAMP, NUMBER          has no opcodes exposing
                                                    caller/block context

  Logging                                          no direct equivalent --
    LOG0, LOG1, LOG2, LOG3, LOG4                   GoChain has no on-chain
                                                    event-log mechanism

  System / calls                                   no direct equivalent --
    CALL, DELEGATECALL, STATICCALL, CREATE,         Ch. 67 Section 8's
    CREATE2, RETURN, REVERT, SELFDESTRUCT            exercise sketches what
                                                     an OpCall might look like

  Cryptographic                                    OpHash160, OpCheckSig
    SHA3 (Keccak-256), and signature recovery        (Ch. 61-63)
    exposed via a precompiled contract
```

Four categories are worth a sentence of context, because they explain *why* the EVM needed gas at all, in a way Bitcoin Script never did:

**Memory** gives a contract a scratch pad — temporary, wiped clean at the end of each call — for holding data too large or too transient to be worth persisting in storage (building up a return value, say). GoChain's VM has no equivalent: every value either sits on the stack or is written straight to persistent contract storage, a deliberate simplification that keeps `vm.VM` small at the cost of forcing every intermediate value through the stack.

**Environment/context opcodes** let a contract read facts about the world it's executing in — who called it, what block it's running in, how much value was sent with the call. GoChain's opcode table has no equivalent family, which is one reason a real token contract in production Ethereum can do things (like emit an event other software watches for) that Chapter 65's simplified token contract does not attempt.

**Logging** (`LOG0`-`LOG4`) lets a contract emit an **event** — a structured record attached to the transaction receipt, not stored in the state trie at all, that off-chain software (a block explorer, a wallet, an exchange) can watch for without needing to re-execute the contract. This is architecturally distinct from Chapter 66's `SSTORE`-backed contract storage: an event is *cheaper* to write precisely because it is never meant to be read back *by* a contract — only by outside observers.

**Control flow** (`JUMP`/`JUMPI`) is the crucial difference from Bitcoin Script named in Chapter 81 Section 4: the EVM permits genuine backward jumps, meaning genuine loops, meaning a contract *can* run forever if nothing stops it. Chapter 64 already told you why gas exists; here is the sharper version of that story — Bitcoin Script avoided the whole problem by never allowing loops in the first place, while the EVM chose expressiveness and paid for it with a mandatory metering system for every single opcode, no exceptions. GoChain's VM, with its `OpJump`/`OpJumpIfFalse` (Chapter 61) and gas metering (Chapter 64) from day one, sits on the EVM's side of this trade-off, not Bitcoin Script's — Chapter 81 Section 4 said as much in passing; this is the fuller picture.

This course does not enumerate an exact, current EVM opcode count for the same reason Chapter 81 didn't enumerate Bitcoin Script's: opcodes have been added over the network's history (and, rarely, repriced or restricted) as the EVM has evolved, so an exact count is a snapshot, not a durable fact. The *categories* above, however, are the structural, long-standing shape of the machine.

Here is an illustrative (not copy-pasted from any specific client) trace showing the EVM's broader scope in action — a tiny "return double the input" contract, stepped through instruction by instruction, next to the closest GoChain VM program can manage with Chapter 61's smaller instruction set:

```
Illustrative EVM bytecode, representative form (input value already
pushed onto the stack by the calling convention):

  PUSH1 0x02      [2]
  MUL             [input * 2]
  PUSH1 0x00      [0, input*2]
  MSTORE          []          -- writes input*2 into memory at offset 0
  PUSH1 0x20      [32]
  PUSH1 0x00      [0, 32]
  RETURN          []          -- returns the 32 bytes starting at
                                  memory offset 0 to the caller

  Notice MSTORE and RETURN: the EVM stages its answer in its scratch
  MEMORY region before handing it back to whatever called this
  contract -- a mechanism Section 6 will show you contract accounts
  rely on constantly when contracts call each other.
```

```go
// The closest GoChain vm equivalent (Chapter 61-62), illustrative --
// not literal bytecode from any real client. GoChain's VM has no
// memory region or RETURN opcode, so the "answer" is simply whatever
// is left on top of the stack when Execute() finishes.
program := []vm.Op{
    vm.OpPushData, // 2
    vm.OpMul,      // stack: [input * 2]
}
```

The GoChain program is shorter not because it is more efficient — it is shorter because it has nowhere else to put an answer except the stack itself, and no way to hand control back to a caller that expects a return value formatted a particular way. That is a genuine, deliberate capability gap, not a coincidence of style.

---

## 5. Gas Pricing in a Live, Congested Network

Chapter 64 gave every opcode a small, fixed gas cost and had the caller supply a gas limit up front — a clean, deterministic story that is entirely correct as far as it goes. What Chapter 64 didn't need to cover, because GoChain is a teaching chain with no real, competing users, is what happens to gas *pricing* the moment thousands of people are trying to get their transactions into the very next block at the very same time.

Think of gas as the metered "how much work am I asking the network to do" number Chapter 64 already built — it does not change from block to block. Gas *price* is different: it is "how much am I willing to pay, per unit of gas, in ether" — and unlike gas itself, this number is a decision the sender makes fresh for every single transaction, based on how badly they want to get included soon.

```
  TOTAL FEE PAID = gas USED * gas PRICE

  gas USED   -- how much computational work the transaction actually
               did, metered opcode-by-opcode exactly like Ch. 64
  gas PRICE  -- how much the sender bids to pay per unit of gas,
               which the sender chooses fresh each time
```

Because block space is limited (every block has a gas limit — a cap on the total gas all its transactions may consume together) and demand for that space fluctuates with real-world activity, gas price behaves like a live auction: when many people want in at once, the going rate to be included soon rises; when the network is quiet, it falls. This is the same fee-market dynamic Chapter 35 introduced conceptually ("what happens when the mempool has more transactions waiting than fit in a block") — Ethereum's mainnet is simply where that dynamic plays out continuously, at real scale, with real financial stakes attached to how long a transaction waits.

Real Ethereum's current fee mechanism (introduced by a network upgrade commonly referred to by its improvement-proposal number, EIP-1559) splits the fee a sender pays into two parts, specifically to make that auction less erratic than a pure highest-bidder system:

```
  TOTAL FEE PER GAS = BASE FEE + PRIORITY FEE (a "tip")

  BASE FEE       -- algorithmically set by the protocol itself, rising
                    or falling automatically block-by-block based on
                    how full the PREVIOUS block was (fuller than
                    target --> base fee rises; emptier --> it falls).
                    Every sender pays this regardless of whom they're
                    "competing" with. This portion is destroyed
                    (burned) rather than paid to the miner/validator.

  PRIORITY FEE   -- a tip the sender adds on top, paid directly to
                    whoever proposes the block, as the actual
                    incentive to include this transaction over another
                    one competing for the same limited block space.
```

This is a genuinely useful, durable engineering idea worth naming on its own: instead of leaving 100% of the price-discovery work to senders guessing at an opaque market (the older, pure "highest bidder wins" model), the protocol *itself* algorithmically tracks recent demand and adjusts a large, predictable portion of the fee automatically — leaving senders to haggle only over the much smaller "how much do I want to jump the queue right now" portion. It is the same instinct as Chapter 26's difficulty adjustment, applied to fee pricing instead of mining difficulty: let the system measure its own recent load and correct itself, rather than relying on any single participant's guess.

```
   BASE FEE SELF-ADJUSTMENT, ACROSS THREE BLOCKS

   Block N   : fuller than target   --> base fee for block N+1 rises
   Block N+1 : still fuller than target --> base fee for block N+2 rises again
   Block N+2 : emptier than target  --> base fee for block N+3 falls
```

Here is a small worked example of the arithmetic involved, so "self-adjusting" doesn't stay an abstract phrase. Suppose a block's target size is 15 million gas, its actual size was 20 million gas (a third over target), and the current base fee is 40 gwei per gas:

```
  fullness ratio   = 20,000,000 / 15,000,000 = 1.333 (33% over target)

  next base fee (illustrative, simplified formula)
                   = current base fee + (current base fee *
                     fullness excess * an adjustment constant)
                   = 40 + (40 * 0.333 * ~0.125)
                   ~= 41.7 gwei per gas

  If instead the block had been HALF full (7.5M of 15M target):
  fullness ratio   = 0.5 (50% under target)
  next base fee    ~= 40 - (40 * 0.5 * ~0.125) ~= 37.5 gwei per gas
```

The exact constant real Ethereum uses is an implementation detail that has been tuned and could be revisited in a future protocol upgrade; the *shape* of the formula — measure how full the last block was relative to a target, nudge the price up or down proportionally, apply automatically with no human or market operator deciding the number — is the durable, structural idea this section is teaching, and it is exactly the same idea Chapter 26's difficulty adjustment applies to mining difficulty instead of fee pricing.

GoChain's Chapter 64 gas model deliberately stops short of any of this: every opcode has one fixed cost, and Chapter 35's fee-sorting greedily fills a block by fee-per-byte, with no algorithmic base fee, no burning, and no live auction dynamics, because reproducing a real fee market requires many competing, independently-motivated senders and miners — exactly the kind of large-scale, adversarial condition a single-learner teaching chain cannot simulate honestly. The *concept* — gas meters work, price is a separate, market-driven number layered on top — is fully present in GoChain; the *market*, with its real congestion and real financial incentives, is what a live network the size of Ethereum's actually adds.

---

## 6. Contract Accounts and External Accounts

Section 2 already showed you that Ethereum has exactly one account structure with four fields. What distinguishes a plain wallet from a smart contract is not a different type — it's whether `codeHash` is empty.

An **externally-owned account (EOA)** is controlled by a private key, exactly like every GoChain `wallet.New()` address from Chapter 36 onward — its `codeHash` is empty, and it can only ever *initiate* a transaction by being signed for, the same way Chapter 33's `Transaction.Sign()` authorizes a spend.

A **contract account** has no private key controlling it at all — it is controlled entirely by the code stored at its `codeHash`, running inside the EVM whenever a transaction (or another contract) calls it. This is the real, full-scale version of Chapter 65's token contract and Chapter 66's `ContractStore`: the token contract's balance table *is* a contract account's storage trie, and its `transfer`/`mint`/`balanceOf` logic *is* the code the EVM runs when that account is called.

```
  EXTERNALLY-OWNED ACCOUNT (EOA)          CONTRACT ACCOUNT
  --------------------------------        --------------------------------
  codeHash: empty                          codeHash: hash of deployed code
  controlled by: a private key             controlled by: its own code,
                                            run by the EVM on every call
  can initiate a transaction: yes          can initiate a transaction: no
                                            (only ever reacts to being
                                            called by an EOA or another
                                            contract)
  GoChain equivalent: wallet.New()         GoChain equivalent: a deployed
  address (Ch. 36)                         token contract (Ch. 65-66)
```

Only an EOA can ever be the one to originate activity on Ethereum — a contract account can call other contracts (this is exactly the mechanism Chapter 67 Section 9's exercise sketches as a hypothetical `OpCall` for GoChain, and exactly the mechanism the real reentrancy bug from Chapter 67 depends on), but every chain of calls has to start from some EOA's signed transaction, tracing all the way back to a real private key, the same root of trust every GoChain transaction has always had since Chapter 33.

---

## 7. The Merge: Swapping Consensus in a Live System

Chapter 77 built `consensus.ProofOfStake` as an alternative, swappable engine behind the exact same interface `consensus.ProofOfWork` implements — deliberately designed so a GoChain node could run under either algorithm without any other package needing to change. This section is why that design choice matters far beyond a teaching exercise: Ethereum actually did this, live, in production, in 2022.

For its first several years, Ethereum secured its chain with proof of work — the same mining-a-nonce-until-the-hash-qualifies mechanism Chapter 24 explained and Chapter 25 implemented, just running Ethereum's own block format instead of Bitcoin's. In parallel, starting in late 2020, Ethereum ran a separate, independent proof-of-stake network — the **Beacon Chain** — that did not yet control the "real" Ethereum (the one holding actual account balances and running real contracts) at all. It existed purely so validators could stake, be tested, and be slashed for misbehavior (Chapter 77's slashing, at real scale) on a live network with real economic stakes, without risking the live, value-bearing chain while the new consensus mechanism proved itself.

**The Merge**, completed in September 2022, was the moment the two were joined: Ethereum's existing chain — with its full transaction history, every account balance, every contract's storage untouched — stopped being secured by proof-of-work mining and started being secured by the Beacon Chain's proof-of-stake validators instead. This is the single most important fact to take away from this section: The Merge changed *only* the consensus layer. Account balances did not change. Contract code and storage did not change. The EVM did not change. Every transaction anyone had ever sent remained exactly as valid as it had always been. Only the answer to "how does the network agree on which block comes next" changed.

```
  BEFORE THE MERGE                       AFTER THE MERGE
  +---------------------------+          +---------------------------+
  | EXECUTION LAYER             |          | EXECUTION LAYER             |
  |  accounts, EVM, contracts,  |          |  accounts, EVM, contracts,  |
  |  transactions -- UNCHANGED   | ------> |  transactions -- UNCHANGED   |
  +---------------------------+  MERGE   +---------------------------+
  | CONSENSUS LAYER              |  EVENT  | CONSENSUS LAYER              |
  |  Proof of Work               |          |  Proof of Stake              |
  |  (miners, mining difficulty)  |          |  (validators, staking,       |
  |                                |          |   the Beacon Chain)          |
  +---------------------------+          +---------------------------+
```

Why is this hard, and why is it worth a whole section? Because a live blockchain cannot be taken offline to be upgraded the way an ordinary piece of software can — every node, run by independent, uncoordinated operators around the world, has to switch to agreeing on the new rule at the exact same point in the chain's history, or the network splits. The Beacon Chain's multi-year, separate-network "trial run" was specifically how Ethereum de-risked this: proof-of-stake consensus proved itself, under real economic incentive to misbehave, on a chain that mattered enough to take seriously but not so much that a bug there could touch a single real account balance — before it was ever allowed to become the mechanism securing the live chain.

It's worth being concrete about what "validators" actually do under proof of stake, since Chapter 77 introduced the term conceptually but a live network gives it real texture. A **validator** is, at heart, the proof-of-stake counterpart to a miner: instead of racing to solve a computational puzzle (Chapter 24), a validator locks up ("stakes") a meaningful amount of the network's own currency as collateral, and is then periodically selected — pseudorandomly, weighted roughly by how much is staked — to propose the next block or to **attest** (cast a vote confirming) a block someone else proposed. If a validator proposes conflicting blocks, or attests to something dishonest, a portion of their staked collateral can be destroyed — Chapter 77's **slashing** — which is what replaces "wasted electricity" as proof-of-stake's cost for attempting to cheat. None of this changed what a block or a transaction *is* on Ethereum; it only changed *who gets to decide* which block comes next and *what they risk* by deciding dishonestly.

This maps directly onto why Chapter 77 designed `consensus.ProofOfStake` behind the exact same interface as `consensus.ProofOfWork`, rather than as a hard fork of GoChain's entire codebase: swapping *how consensus is reached* should be a change that is, as much as possible, isolated from everything consensus secures. GoChain's teaching version of that lesson is "you can run a small testnet under either engine and compare them" (Chapter 77's own exercise). Ethereum's real version of that lesson is "you can migrate a live network worth an enormous amount of real economic value from one consensus engine to a completely different one, with a multi-year rehearsal first, and have every single balance and every single contract survive completely untouched." Same architectural principle — swappable `consensus.Engine`-style boundary — at two wildly different scales of consequence.

---

## 8. Ethereum vs. GoChain: Side-by-Side

```
  +----------------------+---------------------------+---------------------------+
  | AXIS                 | ETHEREUM (real)             | GOCHAIN (this course)      |
  +----------------------+---------------------------+---------------------------+
  | Data model            | Accounts (balance, nonce,   | UTXO (core.Transaction,   |
  |                        | storageRoot, codeHash)       | core.TxInput/TxOutput)    |
  | Smallest unit          | wei                          | gochip                     |
  | State commitment       | state trie (nested Merkle-   | simplified trie (Ch. 57), |
  |                        | Patricia tries per block)    | UTXO set only              |
  | VM                     | EVM: broad opcode set,       | gochain/vm: smaller,       |
  |                        | memory, logging, CALL family | bounded opcode set (Ch.61) |
  | Gas pricing            | base fee (burned) + priority | fixed per-opcode cost,     |
  |                        | fee, live auction under load | Ch. 35 fee-sort, no market |
  | Consensus              | Proof of stake (post-Merge,  | Proof of work (Ch. 25),    |
  |                        | formerly proof of work)      | optional proof of stake   |
  |                        |                              | (Ch. 77), same interface  |
  | Contract state         | contract account's own       | ContractStore, keyed by    |
  |                        | storage trie                 | contract address (Ch. 66) |
  | Events/logging          | LOG0-4 opcodes, event logs   | no on-chain event log      |
  |                        | in transaction receipts       | mechanism                  |
  | Replay protection       | per-account nonce             | not needed -- UTXOs are    |
  |                        |                                | single-use by construction |
  +----------------------+---------------------------+---------------------------+
```

Read this table next to Chapter 81's Bitcoin table and a pattern emerges: GoChain did not simply "pick Bitcoin" or "pick Ethereum" — it borrowed Bitcoin's data model (UTXO) and locking-script pattern, and Ethereum's VM philosophy (a general, gas-metered, loop-capable instruction set with persistent contract storage) and married them into one system. That is not an accident or a compromise; it is precisely the kind of deliberate, informed architectural choice this course has been training you to make since Chapter 31 first laid the two models side by side.

---

## 9. What GoChain Simplified, and Why That Was the Right Call

- **UTXO instead of the account model.** Chapter 31 already gave the reasoning: UTXO's natural parallelism and its lack of a replay-protection nonce made it the simpler ledger to build and fully understand first, while Volume 9's contract storage still gave you an account-*like* persistence story for contracts specifically, without requiring the whole base ledger to be rebuilt around accounts.
- **One trie, not several nested tries.** Chapter 57's simplified Merkle-Patricia trie commits to the UTXO set alone. Building separate state, storage, transaction, and receipt tries — each contributing its own root to the block header — is real, valuable engineering, but it is additional structural complexity layered on the exact same core idea (a trie root committing to a large key-value set) Chapter 57 already teaches completely.
- **No memory region, no logging opcodes, no CALL family.** GoChain's VM keeps every intermediate value on the stack or in persistent storage, has no event mechanism, and (per Chapter 67 Section 8's exercise) does not yet implement inter-contract calls. Each of these is a real EVM capability GoChain intentionally left out to keep `vm.VM`'s opcode table small enough to implement, test, and reason about completely in one volume.
- **No live fee market.** Chapter 64's fixed gas costs and Chapter 35's greedy fee-sort capture the *concept* of metering and prioritization without needing a real, adversarial, many-participant auction to demonstrate it — because a teaching chain, by construction, does not have thousands of independent, competing senders driving genuine congestion.
- **A course-scale consensus swap, not a live one.** Chapter 77 lets you run a small testnet under proof of work or proof of stake and compare them side by side — a genuinely useful hands-on lesson. The Merge did the same swap under a live network holding real value, with a multi-year rehearsal network (the Beacon Chain) de-risking it first — a scale of operational care no course exercise needs to reproduce to teach the underlying architectural lesson.

Every one of these, again, is a deliberate trade-off in service of a course you can complete and fully understand, not a gap that makes GoChain's design "wrong." Chapter 83 now takes both this chapter's comparisons and Chapter 81's and lines them up in one single, systematic table.

---

## Summary

- Ethereum's account model tracks one running balance per address (plus a nonce, storage root, and code hash) rather than a set of spendable UTXOs — trading UTXO's natural parallelism for a natural home for smart-contract state.
- The real Ethereum state trie is several nested Merkle-Patricia tries per block (state, per-contract storage, transactions, receipts), each committed to by a root hash in the block header — the full-scale version of the single simplified trie Chapter 57 built for GoChain's UTXO set.
- The EVM's opcode categories — arithmetic, bitwise, stack, memory, storage, control flow (including genuine loops via `JUMP`/`JUMPI`), environment context, logging, and system/call opcodes — are a strictly larger, more general instruction set than GoChain's `gochain/vm`, and the price of that generality is mandatory gas metering everywhere.
- Real gas pricing splits into an algorithmically self-adjusting base fee (burned) and a sender-chosen priority fee (a tip to the block proposer), letting the network's own recent congestion set most of the price automatically — a live fee market GoChain's fixed per-opcode costs and greedy fee-sort don't attempt to simulate.
- An externally-owned account (controlled by a private key, like every GoChain wallet address) and a contract account (controlled entirely by its own code) are the same account structure, distinguished only by whether `codeHash` is empty.
- The Merge (September 2022) swapped Ethereum's entire consensus layer from proof of work to proof of stake on a live network, while leaving every account balance, every contract's storage, and the EVM itself completely untouched — a real-world validation of exactly the swappable-consensus-engine design Chapter 77's `consensus.ProofOfStake` demonstrates at course scale.
- GoChain deliberately combines Bitcoin's UTXO data model with an Ethereum-style general, gas-metered VM rather than copying either chain wholesale — a considered architectural choice, not an accidental middle ground.

---

## Exercises

### Easy

1. **In your own words (100-150 words), explain why the account model needs a nonce but GoChain's UTXO model does not.** Reference Chapter 34's double-spend check directly in your explanation.
2. **Draw the four fields of an Ethereum account** (Section 2) next to the closest GoChain equivalent for each — some fields (like `storageRoot`) will map to a specific chapter's mechanism (Chapter 66), and at least one will have no direct GoChain equivalent. Name which, and why.
3. **List three EVM opcode categories from Section 4 that GoChain's `vm.VM` has no equivalent for**, and for each one, write one sentence on what capability that gives real Ethereum contracts that GoChain's token contract (Chapter 65) does not have.

### Medium

4. **Work through a base-fee/priority-fee calculation** (Section 5): a transaction uses 21,000 gas, the current base fee is 30 gwei per gas, and the sender sets a priority fee (tip) of 2 gwei per gas. Compute the total fee paid, and state which portion is burned and which portion goes to the block proposer.
5. **Explain, in 150-200 words, why The Merge (Section 7) is described as changing "only the consensus layer."** Specifically address what happened to account balances and contract storage during the transition, and why that separation was possible at all — tying your answer back to the interface design behind `consensus.ProofOfWork` and `consensus.ProofOfStake` from Chapter 77.
6. **Compare the state trie's role (Section 3) to Chapter 57's simplified trie** by describing, in your own words, what would have to change about Chapter 57's implementation to support a *separate* per-contract storage trie (as real Ethereum does) instead of one flat UTXO-set trie.

### Hard

7. **Implement a `nonce` field on GoChain's wallet type** (extending Chapter 36's CLI wallet) purely as a design exercise: even though GoChain's UTXO model does not need one for double-spend protection, add a monotonically increasing counter to each wallet, require every new transaction to include the sender's current nonce plus one, and write a test proving a transaction reusing an old nonce is rejected. Discuss, in a comment, whether this actually adds any security GoChain didn't already have, and why.
8. **Design (on paper) a minimal "event log" mechanism for GoChain's VM**, modeled loosely on the EVM's `LOG0`-`LOG4` opcodes from Section 4: what would an `OpLog` opcode need to do, where would logged data be stored (hint: consider Chapter 66's separation of contract storage from block data), and how would a block explorer (Chapter 72) query for all events emitted by a specific contract address across many blocks?
9. **Research the real EIP-1559 fee mechanism** (the network upgrade that introduced the base fee/priority fee split described in Section 5) and write a 200-300 word summary of one problem it was specifically designed to solve with the fee market that existed before it, citing your source. Connect your answer back to Chapter 35's fee-market discussion.
