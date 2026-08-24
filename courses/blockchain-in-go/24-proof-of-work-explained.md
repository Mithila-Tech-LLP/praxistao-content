# Chapter 24: Proof of Work, Explained

Chapter 23 established that GoChain's nodes need a rule for agreeing on which chain is real, without trusting any single participant. This chapter explains the specific rule: proof of work. The idea is almost embarrassingly simple to state — find a number that makes a hash start with enough zeros — and the entire security of a proof-of-work blockchain rests on that one simple-sounding task being deliberately, provably hard to do quickly, yet trivially easy for anyone else to check once it is done.

## Table of Contents

1. [The Puzzle: Find a Nonce](#1-the-puzzle-find-a-nonce)
2. [Hard to Solve, Easy to Verify](#2-hard-to-solve-easy-to-verify)
3. [The Target: What "Enough Zeros" Actually Means](#3-the-target-what-enough-zeros-actually-means)
4. [Why This Makes Rewriting History Expensive](#4-why-this-makes-rewriting-history-expensive)
5. [A Worked Example, By Hand](#5-a-worked-example-by-hand)
6. [The 51% Attack, Previewed](#6-the-51-attack-previewed)
7. [Where This Idea Came From: Hashcash](#7-where-this-idea-came-from-hashcash)
8. [Energy Use and the Criticisms of Proof of Work](#8-energy-use-and-the-criticisms-of-proof-of-work)
9. [Common Misconceptions, Corrected](#9-common-misconceptions-corrected)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Puzzle: Find a Nonce

Recall from Chapter 09 that `crypto.Hash` takes any bytes and produces a fixed-size, unpredictable-looking fingerprint (SHA-256, 32 bytes). Proof of work asks miners to solve a puzzle built entirely out of this one function you already have:

> Find a number — called a **nonce** ("number used once") — such that when you hash the block's data *together with* that number, the resulting hash starts with at least N zero bits.

That's the whole puzzle. There is no cleverness available to shortcut it, because of a property you already learned in Chapter 08: the **avalanche effect** — changing even one bit of the input scrambles the entire output hash unpredictably. This means there is no pattern connecting "nonce = 41" and "nonce = 42" in the resulting hashes; knowing one tells you nothing whatsoever about the next. The only strategy available is **brute force**: try nonce = 0, hash it, check the result; if it doesn't qualify, try nonce = 1, hash it, check again; keep going until you get lucky.

```
Block data (transactions, prev hash, timestamp, ...)
         |
         |  concatenate with a candidate nonce
         v
+------------------------+
| block data + nonce = 0 |  --> Hash() --> 7f3a91c2...  (doesn't start with enough zeros, reject)
+------------------------+
| block data + nonce = 1 |  --> Hash() --> a08e44b1...  (doesn't start with enough zeros, reject)
+------------------------+
| block data + nonce = 2 |  --> Hash() --> 91cd0053...  (doesn't start with enough zeros, reject)
+------------------------+
             ...
+------------------------+
| block data + nonce = N |  --> Hash() --> 00000a7f...  (starts with enough zeros -- SOLVED!)
+------------------------+
```

**Mining** is exactly this: the process of repeatedly trying candidate nonces until one produces a qualifying hash. A **miner** is whoever is running this search, typically hoping to be rewarded with newly created gochips for being the first to solve a given block (a mechanism GoChain formalizes as the coinbase transaction in Volume 5). The word "work" in "proof of work" refers to the real CPU cycles burned trying nonce after nonce — this is deliberate, wasted-looking effort, and the entire point of this chapter is explaining why that waste is not actually wasteful at all.

---

## 2. Hard to Solve, Easy to Verify

Proof of work's genius is not that it is hard — lots of things are hard. It is that solving the puzzle is hard, but *checking someone else's solution* is nearly instant. This asymmetry is what makes the whole system work, so it is worth making completely concrete:

- **Solving** requires trying, on average, roughly one nonce per unit of "difficulty" (Section 3 makes this precise) before finding a qualifying one — potentially billions or trillions of hash attempts for a real network's difficulty level.
- **Verifying** requires exactly one hash computation: take the nonce the miner claims worked, hash the block data with it once, and check whether the result actually starts with enough zero bits. One computation, microseconds, done.

```
MINER (solving)                          EVERYONE ELSE (verifying)

try nonce=0  -> hash -> no             receive: "nonce = 8,742,193 solves this block"
try nonce=1  -> hash -> no
try nonce=2  -> hash -> no
      ...                              hash(block data + 8,742,193) once
try nonce=8,742,193 -> hash -> YES!          |
                                              v
     (billions of attempts,               00000a7f... starts with enough
      real CPU time spent)                 zeros? YES -> valid, done.
                                            (one hash, microseconds)
```

This "hard to solve, easy to verify" shape is not unique to hashing — it shows up throughout computer science and everyday life. Solving a jigsaw puzzle takes real effort; checking that a *completed* jigsaw puzzle is correctly assembled takes seconds, just by looking at it. Finding the correct combination to a padlock by trying every number takes time; checking whether a *given* combination opens the lock takes one attempt. Proof of work engineers this exact shape into a computational problem: the SHA-256 hash function has no known shortcut for "solving" (finding an input that produces a specific kind of output) faster than trying inputs one by one, while "verifying" (checking one specific input's output) is always a single, cheap computation, by the very definition of what a hash function does.

Because verification is cheap, **every node in the network can independently check that a block's proof of work is genuine**, without trusting the miner who produced it and without redoing the expensive search themselves. Chapter 25's `Validate()` method implements exactly this one-hash check.

---

## 3. The Target: What "Enough Zeros" Actually Means

"Starts with enough zero bits" needs a precise definition, because "enough" is a dial GoChain (and every real proof-of-work chain) needs to be able to turn up or down — this dial is called **difficulty**, and Chapter 26 builds the logic that adjusts it automatically. For now, understand what the dial controls.

A SHA-256 hash is 256 bits long. If you require the hash (interpreted as one giant number) to be *less than* some threshold value — called the **target** — you are effectively requiring some number of leading zero bits, because a smaller maximum value means fewer possible hashes qualify, and the hashes that do qualify are the smallest ones, which are exactly the ones starting with the most zeros.

```
256-bit hash space (every possible SHA-256 output), smallest to largest:

0x000...000                                                    0xFFF...FFF
|<---- target ---->|<--------------- everything else --------------->|
   ^ hashes in here                    ^ hashes in here do NOT qualify
     DO qualify
     (a tiny sliver of
      the whole space)

Lower target  = smaller sliver = fewer qualifying hashes = HARDER
Higher target = bigger sliver  = more qualifying hashes  = EASIER
```

Requiring, say, 20 leading zero *bits* means the qualifying sliver is 1 out of 2^20 (about 1 in a million) of the whole hash space — a miner trying random nonces will need, on average, about a million attempts before landing in that sliver. Requiring 24 leading zero bits shrinks the sliver by another factor of 16, to roughly 1 in 16 million, meaning roughly 16 times more work on average. This is exactly how **difficulty** is expressed numerically in Chapter 25's code: as a target threshold (represented with Go's `math/big.Int`, since 256-bit numbers do not fit in any native integer type), where a smaller target means a smaller sliver means more expected work per block.

Because each attempt is essentially an independent coin flip weighted by the sliver's size, the number of attempts needed follows what statisticians call a **geometric distribution** — the same shape as "how many times do I flip a coin before it lands heads." The expected (average) number of attempts is always `2^(leading zero bits)`, but any single mining attempt can take far fewer or far more tries than that average purely by chance, which is exactly what Section 5's worked example demonstrates numerically. Here is that relationship at a few difficulty levels, assuming a modest single machine hashing at one million attempts per second (1 MH/s):

| Leading zero bits | Expected attempts | Expected time @ 1 MH/s |
|---|---|---|
| 8 | 256 | 0.0003 seconds |
| 16 | 65,536 | 0.07 seconds |
| 20 | ~1,048,576 | ~1 second |
| 24 | ~16,777,216 | ~17 seconds |
| 28 | ~268,435,456 | ~4.5 minutes |
| 32 | ~4.3 billion | ~72 minutes |

Notice each additional 4 bits multiplies the expected work by 16 — this is the exact lever Chapter 26's difficulty-adjustment algorithm turns up and down to keep GoChain's average block time steady as the total mining power on the network changes.

**Difficulty**, defined precisely now for the first time in this course: a number representing how small the target sliver is, and therefore how many hash attempts are expected, on average, before a miner finds a qualifying nonce. Higher difficulty means more expected attempts means more real time and real electricity spent per block, assuming a roughly constant hashing speed.

---

## 4. Why This Makes Rewriting History Expensive

Chapter 19 showed that tampering with an old block is *detectable*, because every hash from that point forward stops matching. But detectable is not the same as *impossible* — nothing in Chapter 19's design stopped an attacker from simply recomputing every hash downstream of their tampered block, one after another, in the time it takes to run a `for` loop. That is precisely the gap proof of work closes.

With proof of work in place, every block's hash is not just "recomputed" — it has to satisfy the target, which means re-*mining* it: running the same expensive nonce search all over again. And it's worse than that for an attacker: because each block's hash feeds into the next block's data (via `PrevBlockHash`), changing one old block invalidates every block's proof of work after it too. To rewrite a transaction buried five blocks deep, an attacker must re-mine that block *and all five blocks after it*, from scratch, while the honest network is simultaneously mining new blocks on top of the real chain.

```
Honest chain:    B0 -> B1 -> B2 -> B3 -> B4 -> B5 -> B6 (growing every ~10s)

Attacker wants to change a transaction inside B2. They must:

  1. Re-mine B2 with the altered transaction (redo the nonce search)
  2. Re-mine B3 (its PrevBlockHash must match the NEW B2's hash)
  3. Re-mine B4 (must match the NEW B3's hash)
  4. Re-mine B5 (must match the NEW B4's hash)
  5. Re-mine B6 (must match the NEW B5's hash)
  ... and keep re-mining new blocks faster than the honest network mines
      real ones, or the honest chain will always be longer and win
      (the longest / most-work chain rule from Chapter 23, built in Volume 7)
```

This is the core insight worth internalizing: **proof of work does not make tampering impossible — it makes tampering cost real, ongoing, compounding computational work, growing with every block that gets mined on top.** An attacker who is slower than the honest network combined can never catch up, because the honest network is adding new proof-of-work-secured blocks the entire time the attacker is trying to redo old ones. This is exactly why more confirmations (Chapter 23, Section 7) mean more safety: each additional block on top is one more re-mining task an attacker would have to redo, and one more block of a head start the honest network gets to extend its own lead.

---

## 5. A Worked Example, By Hand

Let's make this fully concrete with small, hand-traceable numbers (real GoChain difficulty will require far more leading zero bits than this toy example, but the mechanics are identical).

Suppose the puzzle requires a hash — again, treating it as one gigantic number — to be less than a target that, in this simplified example, corresponds to "the first 2 hex characters of the hash must both be `0`" (this is a stand-in for "8 leading zero bits," simplified to hex digits for readability by hand). Take the (fictional, shortened) block data string `"block-7|prevhash=9f2a...|txcount=3"` and imagine hashing it with different nonces appended:

```
nonce =  0  ->  Hash("block-7|...:0")  ->  a3 f9 91 22 ...   (starts "a3", NOT a match)
nonce =  1  ->  Hash("block-7|...:1")  ->  71 0c d4 88 ...   (starts "71", NOT a match)
nonce =  2  ->  Hash("block-7|...:2")  ->  5e 22 aa 01 ...   (starts "5e", NOT a match)
nonce =  3  ->  Hash("block-7|...:3")  ->  d0 4b 6f 3a ...   (starts "d0", NOT a match)
nonce =  4  ->  Hash("block-7|...:4")  ->  00 91 2c 7f ...   (starts "00", MATCH!)
```

It took 5 attempts (nonce 0 through 4) to find a hash starting with `00`. With this toy 2-hex-digit target, roughly 1 out of every 256 hashes (2^8 possibilities for the first byte) will qualify by chance, so on average you would expect to try about 256 nonces before finding one — this particular run got lucky and found it on attempt 5, purely because hash outputs look random and any individual run can be faster or slower than the average. This randomness is fundamental and important: proof of work does not guarantee a fixed number of attempts, only an *expected* (average) number, which is exactly why real block-mining times vary block to block even at constant difficulty, and exactly why Chapter 26's difficulty adjustment has to look at recent *averages* over many blocks, not judge any single block's mining time in isolation.

Once nonce 4 is found, verification is one step: anyone can compute `Hash("block-7|...:4")` themselves, see it starts with `00`, and immediately accept the block — without needing to try nonces 0 through 3 themselves, and without trusting the miner's word for it.

---

## 6. The 51% Attack, Previewed

Section 4 established that rewriting history requires out-pacing the honest network's ongoing mining. This leads to the natural question: what if an attacker actually *could* out-pace the honest network — not by being clever, but simply by controlling more raw hashing power than everyone else combined?

This scenario has a name: the **51% attack** (sometimes called a **majority attack**). If a single attacker (or a colluding group) controls more than half of the network's total mining power, they can, in principle, out-mine the honest network on average, allowing them to eventually build a longer competing chain in secret and then reveal it, causing the network to switch over to the attacker's version under the longest-chain rule — potentially reversing transactions the attacker themselves made (a **double-spend**: paying someone, waiting for them to accept the payment as confirmed, then rewriting history to undo it while keeping the goods or funds).

```
Honest network's combined hash power:  ~49%  ---->  builds Chain H
Attacker's hash power:                  51%  ---->  builds Chain A (in secret)

Given enough time, Chain A grows faster on average than Chain H,
because the attacker alone out-paces everyone else combined.
Eventually Chain A is longer -- the network adopts it, discarding
whatever was only confirmed on Chain H.
```

Two things are worth previewing now, with full technical depth arriving in Volume 11, Chapter 76's hands-on attack lab: first, a 51% attack requires an enormous, genuinely difficult-to-acquire amount of real computing hardware and electricity for any network with meaningful total hash power — this cost *is* the security, not a coincidental side effect. Second, even a successful 51% attack cannot forge transactions from other people's accounts or steal funds directly (it still cannot produce a valid signature it does not hold the private key for, a topic Volume 2 already covered) — its damage is limited to reordering or reversing *its own* recent transactions and briefly disrupting the network's progress, not arbitrary theft. This distinction — "can rewrite recent history" versus "can forge anything" — is one of the most commonly misunderstood aspects of blockchain security, and Volume 11 will let you attempt this exact attack against your own small GoChain test network to see precisely where its power begins and ends.

---

## 7. Where This Idea Came From: Hashcash

Proof of work is not an invention original to Bitcoin or to blockchains at all — it predates them by more than a decade, and knowing its original purpose makes the whole idea feel less like cryptocurrency-specific magic and more like a general-purpose tool blockchains happened to borrow. In 1997, computer scientist Adam Back proposed **Hashcash**, a system designed to fight email spam. The problem Hashcash targeted was economic: sending a million spam emails costs a spammer almost nothing, so even a tiny success rate is profitable. Hashcash's fix was to require every outgoing email to include a small proof-of-work stamp — the sender's mail program had to find a hash meeting a difficulty target (using SHA-1 in the original design) before the email would be accepted by a receiving server, costing a fraction of a second of computer time per message.

```
Sending ONE email (legitimate use):        Sending ONE MILLION emails (spam):
  compute one PoW stamp                      compute one million PoW stamps
  cost: a fraction of a second               cost: hours or days of CPU time,
                                              suddenly no longer "basically free"
```

For a normal person sending a handful of emails a day, this cost is imperceptible. For a spammer trying to send millions of messages, the same small per-email cost, multiplied a million times, becomes prohibitively expensive — the exact same "cheap for one honest actor, expensive at attacker scale" shape that Section 4 described for rewriting blockchain history. Hashcash never achieved widespread adoption for email (spam filtering evolved in other directions instead), but Satoshi Nakamoto's 2008 Bitcoin whitepaper explicitly cites Hashcash as the direct inspiration for using proof of work to secure a public, trustless ledger — repurposing "make cheating expensive by requiring provable computational effort" from a spam-fighting tool into the entire security foundation of a currency. GoChain's `consensus.ProofOfWork`, built in the next chapter, is a direct descendant of this fifteen-years-earlier idea.

---

## 8. Energy Use and the Criticisms of Proof of Work

It would be incomplete to explain proof of work without naming its most famous criticism honestly: it is, by design, energy-intensive. The security argument from Section 4 — "attacking the chain costs real, ongoing computational work" — only functions as a security guarantee *because* that work is genuinely expensive to produce. There is no way to have the security property without the expense; they are the same thing viewed from two different angles. Real proof-of-work networks like Bitcoin consume electricity comparable to a mid-sized country, and this has drawn substantial environmental criticism, along with hardware arms races (specialized ASIC — Application-Specific Integrated Circuit — mining hardware built for nothing but computing SHA-256 as fast as physically possible) that push mining into the hands of well-capitalized operations rather than hobbyists with a home computer.

This is precisely the motivation behind **proof of stake** (previewed in Chapter 23, built in Volume 11, Chapter 77), which achieves a broadly similar security goal — making it expensive to attack the network — by requiring participants to lock up *financial value* they stand to lose (be "slashed") if they misbehave, rather than requiring them to burn *electricity* continuously whether they misbehave or not. Neither approach is a strictly better engineering choice than the other in every dimension; they trade off differently on decentralization, energy use, hardware requirements, and the exact shape of the guarantees they provide, which is why this course builds both, behind the same `consensus.Engine` interface, so you can compare them directly with working code rather than only in the abstract.

---

## 9. Common Misconceptions, Corrected

A few misunderstandings come up often enough with new blockchain learners that they are worth naming and correcting explicitly before moving on to code:

- **"Miners are solving a useful math problem."** They are not. The nonce search has no meaning outside of qualifying for the target — it is not factoring numbers, folding proteins, or doing anything with value beyond making cheating expensive. Some alternative schemes (outside this course's scope) do attempt "useful work" proof of work, but standard SHA-256 proof of work is deliberately meaningless busywork by design.
- **"A more powerful computer can find a mathematical shortcut."** No amount of algorithmic cleverness changes the number of *expected* attempts — only raw hashing speed (more attempts per second) helps, which is why mining hardware competes on hashes-per-second, not on smarter search strategies.
- **"Once you solve a block, you keep trying the same nonce for the next one."** Each new block has entirely different data (different transactions, a different previous hash, a new timestamp), so the search starts over completely from nonce zero — nothing about the previous block's winning nonce carries over or helps.
- **"Proof of work prevents all forms of cheating."** It specifically prevents cheaply rewriting history and lets nodes agree on one chain without trust. It does not, by itself, prevent someone from spending coins they legitimately own in ways you might not like, and it cannot forge a signature for coins the attacker never controlled — that protection comes from the ECDSA signatures built in Volume 2.

---

## Summary

- Proof of work's puzzle: find a **nonce** such that `Hash(block data + nonce)` produces a result below a numeric **target** (equivalently, with enough leading zero bits).
- There is no shortcut to solving it — the avalanche effect means each nonce's hash is unpredictable from the last, so **mining** means brute-force trying nonces one after another.
- The puzzle is deliberately **hard to solve** (many attempts on average) but **easy to verify** (exactly one hash computation), and this asymmetry is what lets every node check a block without trusting the miner or redoing the search.
- **Difficulty** is the dial controlling how small the target sliver is, and therefore how many attempts are expected on average — Chapter 26 makes this adjust automatically to keep block times steady.
- Tampering with an old block requires re-mining it *and every block after it*, faster than the honest network is simultaneously mining new ones — this is what makes deep history practically irreversible, and why more confirmations mean more safety.
- A worked hand example showed candidate nonces 0-4 hashed against a toy 2-hex-digit target, landing on a match at nonce 4 — five attempts is faster than the theoretical average of 256, illustrating that proof of work guarantees an *expected* number of attempts, not a fixed one.
- The **51% attack** is what happens if an attacker controls a majority of network hash power: they can eventually out-mine the honest network and reverse their own recent transactions, but they still cannot forge signatures or steal others' funds directly.
- Proof of work's security is inseparable from its energy cost — this is exactly the trade-off proof of stake (Volume 11) is designed to address differently.

---

## Exercises

### Easy

1. In your own words, explain the "hard to solve, easy to verify" property of proof of work, using the padlock-combination analogy or one of your own.
2. Why can't a miner predict, from having just tried nonce = 100 and failed, anything useful about whether nonce = 101 will succeed?
3. Define, in your own words: nonce, target, difficulty, 51% attack.

### Medium

4. If a target requires 16 leading zero bits instead of 8, roughly how many times more hash attempts would you expect a miner to need on average? Show your reasoning using powers of two.
5. Explain why re-mining a single old block is not enough for an attacker to successfully rewrite history — what else must they do, and why does the honest network's continued mining make this harder over time?
6. A friend argues: "Since verifying a proof-of-work solution only takes one hash, and computers can do billions of hashes per second, verification isn't really that much faster than solving." Explain what is wrong with this reasoning, being precise about what "verifying" and "solving" each actually require.

### Hard

7. Using the toy example from Section 5, extend the hand-traced search: continue trying nonces 5 through 10 with hashes of your own invention (you don't need real SHA-256 output — invent plausible-looking hex strings), and identify the first one (if any) that would satisfy a *stricter* toy target requiring the first THREE hex characters to be `0`. Explain how much rarer this qualifying condition is compared to the two-character version from the chapter.
8. Research the actual mechanism ASIC mining hardware uses to be faster than a general-purpose CPU or GPU at SHA-256 hashing specifically. Why can't a similar speedup be achieved for, say, general-purpose Go programs, and what does this imply about the "meritocracy" of proof-of-work mining in practice?
9. Suppose GoChain's target network block time is 10 seconds (a number this course will use starting in Chapter 26), and the total network hash rate suddenly quadruples overnight because a large new mining operation joins. Without yet knowing Chapter 26's exact algorithm, reason qualitatively: what will happen to real block times if difficulty is not adjusted? What problems could that cause for applications built on top of GoChain that assume roughly steady block times (for example, an exchange crediting deposits after a fixed number of confirmations)?
