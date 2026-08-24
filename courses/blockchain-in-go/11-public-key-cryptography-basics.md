# Chapter 11: Public-Key Cryptography Basics

Chapters 08 through 10 solved one problem: proving data has not changed, using hashes and Merkle trees. This chapter starts solving a completely different problem: proving *who* did something. A blockchain has no bank teller checking your ID, no central server holding a list of usernames and passwords, and no notary standing over your shoulder. Yet somehow, when a GoChain transaction says "send 5 gochips from this address," every node on the network needs to be convinced that the actual owner of that address approved it, not an imposter. Public-key cryptography is the tool that makes this possible, and this chapter builds a solid, plain-language understanding of how, before Chapter 12 introduces the specific algorithm (ECDSA) GoChain uses, and Chapter 13 turns it into real Go code.

## Table of Contents

1. [The Problem: Proving Identity Without a Middleman](#1-the-problem-proving-identity-without-a-middleman)
2. [What Is a Key Pair?](#2-what-is-a-key-pair)
3. [The Wax Seal Analogy](#3-the-wax-seal-analogy)
4. [Why You Never Share Your Private Key](#4-why-you-never-share-your-private-key)
5. [Symmetric vs. Asymmetric Cryptography](#5-symmetric-vs-asymmetric-cryptography)
6. [Elliptic Curve Cryptography, Conceptually](#6-elliptic-curve-cryptography-conceptually)
7. [Why "Easy Forward, Infeasible Backward" Matters](#7-why-easy-forward-infeasible-backward-matters)
8. [Why Elliptic Curves Instead of RSA](#8-why-elliptic-curves-instead-of-rsa)
9. [From Key Pair to Signature to Address](#9-from-key-pair-to-signature-to-address)
10. [Common Misconceptions, Addressed](#10-common-misconceptions-addressed)
11. [Where This Fits in GoChain](#11-where-this-fits-in-gochain)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. The Problem: Proving Identity Without a Middleman

Think about how you prove your identity in everyday life. At a bank, you show a driver's license, and the teller trusts the government agency that issued it. On a website, you type a password, and the website's server checks it against a stored copy and trusts its own database. In both cases, proving who you are depends on trusting some third party — the government, the bank, the website's server — to vouch for you or to keep a secret safely on your behalf.

A blockchain has no such third party by design. Chapter 01 already established that a blockchain is meant to work without a central authority everyone must trust. So when Alice wants to send gochips to Bob, and every node on the network needs to agree that this really is Alice's request and not a forgery, there is no bank teller and no login server to ask. The network needs a way to verify identity using nothing but math — math that any node, anywhere, can check independently, with zero coordination, the same way Chapter 08 showed nodes independently computing identical hashes without comparing notes.

It is worth knowing this problem is not new, and was not invented for blockchains. In 1976, Whitfield Diffie and Martin Hellman (building on ideas also independently developed around the same time by Ralph Merkle — no relation to this course's Merkle trees beyond sharing an inventor's name) published the first practical scheme for two parties to establish a shared secret over a channel someone else might be listening to, without ever transmitting the secret itself. That paper is widely regarded as the founding moment of public-key cryptography, and nearly fifty years later, the exact family of ideas it introduced is what makes it possible for a stranger's laptop, running GoChain code it downloaded from the internet, to verify a transaction from someone it has never met and will never meet, with no institution vouching for either party.

Passwords will not work for this. A password is a **shared secret** — you and the party checking it both need to know the same string, and the moment you use it to "prove" who you are, you have to send it somewhere for comparison, which means it can be intercepted, or the party checking it could be dishonest, or a data breach could leak every password in the checking party's database at once. Public-key cryptography solves the entire problem differently: it lets you prove you know a secret without ever revealing the secret itself, and without needing to send it anywhere for anyone to check against a stored copy.

```
  Password-based proof:                 Public-key-based proof:

  You  ──"my password is X"──▶  Server  You  ──"here is proof I know X,
         (X travels, and must         but X itself never leaves my
          be stored somewhere)         computer"──▶  Anyone
```

---

## 2. What Is a Key Pair?

A **key pair** is two large numbers, mathematically linked to each other, generated together as a matched set:

- The **private key** is a secret, randomly generated number that only its owner ever knows. It is never shared, never transmitted, never shown to anyone — not even to prove who you are.
- The **public key** is derived from the private key through a specific mathematical process (Section 6 explains this process conceptually), and it is meant to be shared freely. You can hand your public key to anyone, post it publicly, or attach it to every transaction you send, with no loss of security.

The critical property that makes this whole scheme work is a one-way relationship between the two: computing the public key *from* the private key is fast and straightforward. Computing the private key back *from* the public key is, for all practical purposes, impossible — not "hard, but a determined attacker with enough patience could eventually do it," but the same flavor of "impossible" Chapter 08, Section 5 used for finding a SHA-256 collision: technically not ruled out by pure logic, but so far beyond what any computer on Earth could accomplish that treating it as impossible is the only sane engineering assumption.

```
              generation
  random number  ────────▶  private key (secret, kept forever)
                                    │
                                    │  one-way derivation (Section 6)
                                    ▼
                            public key (shared freely, everywhere)

  Forward (private → public): fast, easy, one function call.
  Backward (public → private): computationally infeasible.
```

Notice the shape of this diagram: it is the exact same "easy forward, infeasible backward" shape as Chapter 08's one-wayness property for hash functions. That is not a coincidence — it is the same underlying engineering principle (build a function that is trivial to compute in one direction and infeasible to invert) applied to a different mathematical problem. Chapter 08 hashed *data* one-way. This chapter derives a *public key* from a *private key* one-way. Keep this parallel in mind; it will make Section 6 feel far less like a brand-new idea than it might otherwise.

Every participant on GoChain generates their own key pair completely independently, with no registration step and no need to inform anyone else first — the same permissionless spirit Chapter 01 introduced for the network as a whole. Alice, Bob, and Carol, from Chapter 08's determinism example, each just run a key-generation function on their own laptop, in their own city, and each ends up with their own private/public key pair that nobody else had any hand in creating:

```
   Alice's laptop:  generate key pair  ──▶  private key A, public key A
   Bob's laptop:    generate key pair  ──▶  private key B, public key B
   Carol's laptop:  generate key pair  ──▶  private key C, public key C

   Three independent computers. No coordination. No registration
   authority. Each key pair is unique because each private key was
   drawn from an astronomically large space of random numbers,
   overwhelmingly unlikely to ever collide with anyone else's.
```

This is worth sitting with for a moment, because it is easy to assume some central system must be handing out key pairs to keep everyone unique, the way a government issues unique passport numbers. Nothing of the sort happens. Chapter 13's `NewKeyPair()` function will simply reach into Go's cryptographically secure random number source, pull enough random bits to make a private key from an astronomically large space of possibilities, and derive the matching public key — and the odds of two people, anywhere, ever generating the same private key by chance are close enough to zero that (exactly as Chapter 08, Section 5 argued for hash collisions) treating it as impossible is the only sane engineering assumption.

---

## 3. The Wax Seal Analogy

Picture a medieval merchant sending an important letter. Before sealing it, the merchant presses a custom-carved signet ring into a blob of hot wax on the envelope's flap. That ring is unique — carved specifically for this one merchant, and no one else owns an identical one. Anyone receiving the letter can look at the wax seal's imprint and instantly recognize, just by sight, that it came from this specific merchant, because they (or people they trust) have seen this merchant's seal before and know its exact pattern.

Here is the part that matters most for this chapter: the recipient never needs to *borrow the ring itself* to confirm the letter is genuine. They only need to see the impression the ring left in the wax. The ring — the thing that actually produces the mark — never leaves the merchant's possession. If the ring ever fell into someone else's hands, that person could now forge seals that look exactly like the merchant's, which is exactly why the merchant guards it so carefully.

```
   Merchant's ring (kept secret,        Wax seal on envelope
   never leaves the merchant)                (shared publicly,
        │                                   anyone can inspect it)
        │  presses into wax
        ▼
   ┌─────────┐         produces         ┌─────────┐
   │  ring    │  ───────────────────▶  │   seal   │  ──▶  Recipient checks
   │ (private)│                          │ (public) │       the seal's pattern
   └─────────┘                          └─────────┘       against what they
                                                             know of this
                                                             merchant's mark
```

Map this directly onto public-key cryptography:

| Wax seal analogy | Public-key cryptography |
|---|---|
| The signet ring | The **private key** |
| The wax impression on the envelope | A **signature** (Chapter 12) |
| Recognizing the merchant's mark | **Verifying** the signature using a **public key** |
| The ring never leaves the merchant | The private key never leaves its owner |
| Anyone can inspect a sealed envelope | Anyone can verify a signature |

One place the analogy needs sharpening: a real wax seal can be physically copied by someone skilled enough, given enough time with the ring or a good enough mold. A cryptographic signature has no such weakness — Section 6 and Chapter 12 explain exactly why forging one is not merely difficult in the way carving a duplicate ring is difficult, but infeasible in the much stronger sense Chapter 08 established for hash collisions: the search space involved is so astronomically large that no realistic amount of computing power gets you there.

---

## 4. Why You Never Share Your Private Key

It is worth stating this as plainly as possible, because it is the single most important operational rule in this entire volume: **you never send, show, transmit, or reveal your private key to prove your identity — to anyone, for any reason, ever.** This is a sharp break from how passwords work, and it is easy to instinctively reach for the wrong mental model here, so it is worth walking through exactly why.

With a password, proving who you are *requires* revealing the secret (typing it into a login box) so the other party can check it against their stored copy. Every time you do this, the secret travels somewhere, and every place it travels is a place it could be intercepted, logged, or leaked.

With a key pair, proving who you are means producing a **signature** — a separate, derived value, tied to one specific piece of data, that anyone can check against your *public* key without ever needing your private key at all. Chapter 12 covers exactly how a signature is produced and checked; for now, the important shape to internalize is this:

```
  Password model:
    prove identity  ──▶  reveal the actual secret (the password itself)

  Key pair model:
    prove identity  ──▶  produce a signature derived from the secret
                          (the secret itself never travels anywhere)
```

If your private key were ever exposed — stolen from your computer, leaked through a bug, guessed because it was generated carelessly — an attacker could produce signatures indistinguishable from your genuine ones, and every node on the network would accept them as authentically yours, because the entire system's trust rests on the assumption that only you could have produced a valid signature under your public key. This is why Volume 6 dedicates an entire chapter (Chapter 40) to encrypting private key files at rest, and why hardware wallets (Chapter 41) exist specifically to ensure a private key never touches an internet-connected computer's memory at all. None of that machinery would matter if the underlying rule from this section — never transmit or expose the private key itself — were not true in the first place.

This is not a hypothetical worry invented for this course. A large share of the real-world cryptocurrency theft you may have read about in the news traces back to exactly this rule being broken somewhere: an exchange storing thousands of customers' private keys on an internet-connected server that got breached, a developer accidentally committing a private key into a public code repository, a phishing site tricking someone into typing their private key (or the seed phrase that generates it, previewed in Exercise 9) into a fake wallet interface. None of these are failures of the underlying mathematics from Sections 6 and 7 — the elliptic curve discrete logarithm problem was not broken in any of these incidents. Every one of them is a failure of *operational* security: the private key touched somewhere it should never have touched, and once it did, the cryptography protecting it stopped mattering, the same way a bank vault's steel walls stop mattering the moment someone leaves the combination taped to the front door.

---

## 5. Symmetric vs. Asymmetric Cryptography

It helps to place public-key cryptography next to the kind of cryptography most people encounter first: encrypting a file with a single password so only someone who also knows that password can unlock it (a ZIP file with a password, for example). That style is called **symmetric cryptography**: the *same* key both locks and unlocks the data, which means anyone who can decrypt something could equally well have been the one who encrypted it, and the key must somehow be shared between both parties in advance, safely, without anyone else intercepting it — a real, practical headache called the **key distribution problem**.

**Asymmetric cryptography** (also called **public-key cryptography** — the two terms describe the same idea) uses two *different* mathematically linked keys instead of one shared key, exactly as Section 2 described: one for locking or signing, a different one for unlocking or verifying, and only one of the two ever needs to be secret. This sidesteps the key distribution problem entirely for the use case this course cares about most: your public key can be handed out to literally anyone, published on a website, or embedded directly in a transaction, and none of that exposure weakens your security at all, because the public key alone cannot be used to forge your signature.

```
  Symmetric (one shared key):              Asymmetric (two linked keys):

  Alice's key: 🔑                          Alice's private key: 🔑 (secret)
  Bob's key:   🔑  (must be the            Alice's public key:  🔓 (shared)
               same key as Alice's,
               shared secretly first)      Bob only ever needs Alice's
                                           public key — nothing secret
                                           has to be exchanged at all.
```

Symmetric cryptography is not obsolete — it is generally much faster and is exactly what Chapter 40's wallet-file encryption will use, because encrypting your *own* file with your *own* password has no key-distribution problem to solve (there is only one party involved: you). Asymmetric cryptography earns its computational cost specifically for the problem this volume is solving: proving identity and authorizing actions between parties who have never met and share no prior secret, which is precisely GoChain's situation with every transaction it processes.

---

## 6. Elliptic Curve Cryptography, Conceptually

Public-key cryptography needs a specific mathematical process to actually derive a public key from a private key — Section 2's "one-way derivation" arrow has to mean something concrete. **Elliptic curve cryptography (ECC)** is the family of techniques GoChain (and Bitcoin, and Ethereum) uses to do this. This section explains what ECC accomplishes and why it behaves the way it does, deliberately without the underlying algebra — the actual curve equations and point-arithmetic rules are precise, well-studied mathematics that Go's standard library already implements correctly (exactly as Chapter 08 treated SHA-256's internals: a process worth understanding at a conceptual level, not worth hand-rolling yourself).

Here is the conceptual shape. Everyone using a given elliptic curve agrees, publicly and in advance, on one fixed, publicly known starting point on that curve, conventionally called the **generator point**. There is also a defined operation — think of it as a specific, well-defined kind of "step" you can take from any point on the curve to reach another point on the same curve — that behaves, for the purposes of this chapter, like a strange, one-directional form of multiplication: you can take the generator point and "multiply" it by a number to land on a new point, and you can take that new point and multiply it again, and so on.

A **private key** in ECC is simply a single, enormous, randomly chosen number. A **public key** is the point on the curve you land on after "multiplying" the generator point by that private number, using the curve's defined operation:

```
  Fixed, publicly known generator point:  G

  Alice's private key (a huge secret number):  d

  Alice's public key = "d applied to G" via the curve's
  point-multiplication operation, written  d * G

       G  ──── apply d times via curve's rule ────▶  Alice's public key point
      (public,                                        (public, shareable,
       known to                                        derived from a
       everyone)                                       secret nobody else
                                                        knows)
```

The one-way property Section 2 promised comes directly from this operation's behavior. Going forward — given the number `d` and the starting point `G`, computing the resulting point `d * G` — is fast; modern computers do this in microseconds even for enormous private keys. Going backward — given only the starting point `G` and the resulting public point, figuring out *what number `d`* was used to get there — has no known efficient shortcut at all. This backward problem has a name, the **elliptic curve discrete logarithm problem**, and cryptographers have studied it intensely for decades without finding a practical way to solve it for a well-chosen curve and a sufficiently large private key. It is, in spirit, exactly the same kind of asymmetry as Chapter 08, Section 4's one-wayness: not a coincidence, but the same design principle solving a related problem.

An analogy that captures the *feel* of this, without claiming to be literally accurate to the curve math: imagine a locked wheel that only ever spins in one direction, starting at a fixed, publicly marked position. Give someone the number of "clicks" you spun it, and they can spin an identical wheel the same number of clicks and land on your exact final position in seconds. But hand someone *only* the final resting position of your wheel, with no record of how many clicks it took to get there, and their only option is to start spinning their own wheel from the beginning and count, one click at a time, hoping to land on the same spot — and for a real elliptic curve with a properly sized private key, the number of possible starting click-counts is so large that this brute-force approach would take dramatically longer than the current age of the universe, using every computer on Earth at once, echoing exactly the scale argument Chapter 08, Section 5 made about SHA-256 collisions.

It is worth being precise about what "no known efficient shortcut" actually claims, because it is easy to overstate. Nobody has *proven*, in the strict mathematical sense that a proof in geometry class is proven, that the elliptic curve discrete logarithm problem is impossible to solve efficiently — it remains, technically, an open question whether some future breakthrough could crack it. What can be said, with decades of evidence behind it, is that an enormous, well-funded, globally distributed community of cryptographers and mathematicians has been actively trying to find a faster way to solve this exact problem for a very long time, on curves used to secure trillions of dollars of value, and nobody has succeeded. That is the same "battle-tested, not mathematically absolute" style of confidence Chapter 08, Section 8 built around SHA-256 — a claim earned through sustained, adversarial scrutiny rather than a claim proven once and settled forever.

---

## 7. Why "Easy Forward, Infeasible Backward" Matters

It is worth being explicit about why this specific mathematical shape — trivial to compute one way, infeasible the other way — is exactly what a signature scheme needs, rather than some other arrangement of numbers.

If deriving a public key from a private key were *also* easy to reverse, the entire scheme would collapse instantly: anyone who ever saw your public key (which, remember, is meant to be shared with literally everyone) could immediately work backward to recover your private key, and from there forge signatures indistinguishable from your genuine ones. The whole point of splitting a key pair into a shareable half and a secret half only works because sharing the public half provably does not leak the private half.

```
  If public → private were reversible (it is NOT, but imagine it were):

    Alice's public key  ────▶  anyone can recover  ────▶  Alice's private key
    (shared with everyone,                                (should be secret,
     posted in every transaction)                          now trivially exposed)

    Result: total collapse. Every public key doubles as a leaked
    private key. Nothing in this volume would work.

  Reality — the actual, infeasible-to-reverse relationship:

    Alice's public key  ────▶  no practical way back  ──▶  Alice's private key
                                                              stays secret
```

This is also why Section 6 emphasized that the private key needs to be *enormous* and *properly random* — a private key chosen from too small a range of possibilities, or generated from a predictable source (a weak pseudo-random number generator, a private key derived from something guessable like a birthday or a short passphrase), reintroduces exactly the brute-force weakness Chapter 08, Section 4 warned about for passwords: an attacker does not need to break the underlying math at all if they can simply guess or narrow down the private key some other way. Chapter 13's Go implementation leans entirely on `crypto/rand`, Go's cryptographically secure random number source, specifically to avoid ever being the weak link in an otherwise sound system.

---

## 8. Why Elliptic Curves Instead of RSA

Elliptic curve cryptography is not the only way to build a one-way, asymmetric key relationship. The older, more famous alternative is **RSA**, named after its inventors (Rivest, Shamir, and Adleman), which bases its one-wayness on a different hard problem: factoring the product of two very large prime numbers back into those two primes. Multiplying two huge primes together is fast; given only the product, finding the two original primes is believed to be infeasible for a sufficiently large product — a different mathematical problem than the elliptic curve discrete logarithm problem from Section 6, but the same "easy forward, infeasible backward" shape running underneath it.

If RSA solves the same basic problem, why does GoChain (like Bitcoin and Ethereum before it) use elliptic curves instead? The practical answer is key size, and it matters more than it might sound like it should. Security researchers estimate how much brute-force effort it would take to break a key of a given size, for a given underlying hard problem, and the two problems do not scale the same way at all — RSA needs a much larger number of bits to reach the same practical security level as ECC:

```
  Roughly equivalent security levels (illustrative, not exact):

  ECC key size        RSA key size        Rough security level
  ─────────────        ─────────────        ─────────────────────
  256 bits              ~3072 bits           strong, widely used today
  384 bits              ~7680 bits           very strong

  A 256-bit ECC private key is a number small enough to write on
  a postcard. A comparably secure RSA key is more than ten times
  as many bits -- and RSA's public keys, ciphertexts, and
  signatures grow correspondingly larger too.
```

For a blockchain, this size difference is not academic. Every signature GoChain produces gets embedded directly inside a transaction, and every transaction gets stored, forever, inside a block, and every block gets copied to every full node on the network (Volume 7) and re-downloaded by every new node that ever joins (Chapter 49). A signature scheme that produces smaller signatures for the same security level means smaller blocks, faster network transmission, less disk space per node, and faster signature verification during mining — all real, compounding costs that RSA's larger keys and signatures would make meaningfully worse at blockchain scale. This is the same style of pragmatic, real-world reasoning Chapter 08, Section 8 used to justify SHA-256 over exotic alternatives: not "elliptic curves are mathematically superior in some absolute sense," but "smaller keys and signatures at an equivalent security level matter enormously once you multiply that saving across every transaction in a system meant to run forever."

---

## 9. From Key Pair to Signature to Address

This chapter has stayed deliberately abstract — "a private key," "a public key," "a one-way derivation" — because the next three chapters each turn one piece of this into something concrete and runnable. It is worth previewing the whole pipeline now, so each upcoming chapter's job is clear before you are in the middle of it:

```
  Chapter 11 (this chapter)     Chapter 12              Chapter 13
  ┌─────────────────────┐      ┌─────────────────┐      ┌──────────────────┐
  │ private key (secret) │      │ ECDSA: the      │      │ crypto/ecdsa in  │
  │        │              │ ──▶  │ specific        │ ──▶  │ Go. NewKeyPair,  │
  │        ▼              │      │ algorithm that  │      │ Sign, Verify --  │
  │ public key (shared)   │      │ signs & verifies│      │ real, running    │
  └─────────────────────┘      └─────────────────┘      │ code             │
                                                          └──────────────────┘
                                                                    │
                                                                    ▼
                                                          Chapter 14
                                                          ┌──────────────────┐
                                                          │ Turn the public  │
                                                          │ key into a short,│
                                                          │ human-friendly,  │
                                                          │ typo-resistant   │
                                                          │ address          │
                                                          └──────────────────┘
```

Chapter 12 names the specific curve-and-signature algorithm GoChain uses (ECDSA, running on top of the elliptic curve concepts from Section 6) and diagrams exactly how signing and verifying work, before any Go code appears. Chapter 13 makes it real: `gochain/crypto.NewKeyPair()`, `Sign()`, and `Verify()`, with a full worked example of Alice signing a message and anyone verifying it. Chapter 14 solves a problem this chapter has not touched at all: a raw public key is a long, ugly string of bytes nobody wants to type by hand, so we compress it down into a short, readable **address**, with a built-in typo-catching checksum.

---

## 10. Common Misconceptions, Addressed

**"My public key is like my username, and my private key is like my password."** Close, but the mental model breaks in an important way. A username and password are two independent secrets you choose (or are assigned); a public and private key are one mathematically *linked* pair, and the entire security model depends on that link being one-way. Also, unlike a password, you genuinely never type or transmit your private key to use it — Section 4 covered exactly why.

**"If someone has my public key, they can access my funds."** No. A public key alone lets someone *verify* signatures you produce and (once addresses exist, Chapter 14) send funds *to* you. It grants zero ability to *spend* anything or produce a valid signature on your behalf. This is the entire point of the one-way relationship from Section 6.

**"Elliptic curve cryptography is unbreakable."** No cryptographic system earns the word "unbreakable" in an absolute sense — Chapter 08, Section 5 made the same point about SHA-256's collision resistance. What is true, and is the actually defensible claim, is that for a well-chosen curve and a properly generated, sufficiently large private key, no known algorithm solves the elliptic curve discrete logarithm problem faster than an infeasible brute-force search — a claim that has held up under decades of public scrutiny by cryptographers actively trying to break it, which is exactly the kind of pragmatic, track-record-based trust Chapter 08, Section 8 built for SHA-256.

**"Losing my private key means I just need to reset it, like a forgotten password."** No, and this misconception is worth dwelling on because it has caused real, permanent losses. A forgotten website password can be reset because the website controls a separate recovery mechanism (an email link, a support ticket) — a central authority exists that can override the check. A blockchain has no such authority by design (Section 1). If a private key is lost, with no seed phrase or backup, whatever funds or identity that key controlled are gone, permanently and unrecoverably, because no one — not GoChain's developers, not any node on the network, not anyone — has any mechanism to override the math and generate a replacement. This is the flip side of the same design that makes the system trustworthy without a central authority: nobody can arbitrarily reassign your funds to a new key, but that also means nobody can rescue you if the original key disappears.

**"Quantum computers already break this."** Not yet, but this deserves a precise answer rather than a shrug. A sufficiently powerful quantum computer running an algorithm called Shor's algorithm *could*, in theory, solve the elliptic curve discrete logarithm problem efficiently — a fundamentally more serious threat than the mild, well-understood erosion Grover's algorithm poses to hash functions (Chapter 08, Section 8). No quantum computer today comes anywhere close to the scale needed to threaten real-world ECC keys, and this remains an active area of research (into "post-quantum" signature schemes) rather than an immediate practical concern for a learning project like GoChain — but it is a genuinely different, larger risk than the one hashing faces, and it is worth knowing the two are not the same size of problem.

---

## 11. Where This Fits in GoChain

Every transaction GoChain will build starting in Volume 5 needs exactly the properties this chapter introduced: proof that a specific private key's owner authorized this specific transaction, verifiable by anyone holding only the corresponding public key, with the private key itself never appearing anywhere on the network, in any block, or in any message sent between nodes. Wallets (Volume 6) are, at their core, just a friendlier interface around generating, storing, and using exactly the key pairs this chapter described. Addresses (Chapter 14) are just a public key, reshaped into something a human can read and type without a doctorate in cryptography. None of it works without the one-way relationship this chapter spent most of its length explaining — everything else in this volume is building machinery on top of that single mathematical fact.

---

## Summary

- A **key pair** consists of a secret **private key** and a shareable **public key**, mathematically linked so the public key is easy to derive from the private key, but the reverse is computationally infeasible.
- The **wax seal analogy**: the private key is like a signet ring that never leaves its owner; the public key is like the recognizable seal pattern anyone can check an impression against, without ever needing the ring itself.
- You **never transmit your private key to prove identity** — instead, you produce a **signature** (Chapter 12) that anyone can verify using your public key alone, unlike a password, which must be revealed to be checked.
- **Symmetric cryptography** uses one shared key for both locking and unlocking; **asymmetric (public-key) cryptography** uses two linked keys, sidestepping the problem of safely sharing a secret key in advance.
- **Elliptic curve cryptography** derives a public key by applying a defined, one-way "point-multiplication" operation to a fixed, publicly known starting point, using the private key as the (secret) number of times to apply it.
- The hard direction of that operation is called the **elliptic curve discrete logarithm problem** — going from a public key back to the private key that produced it has no known practical shortcut, echoing the same "easy forward, infeasible backward" shape as hash functions' one-wayness from Chapter 08.
- A private key must be an enormous, properly random number — a weak or guessable private key breaks the whole scheme through brute-force guessing, with no need to break the underlying math at all.
- Chapters 12 through 14 turn this chapter's concepts into a specific algorithm (ECDSA), real Go code, and a human-readable address format, in that order.

---

## Exercises

### Easy

1. In your own words, extend the wax seal analogy from Section 3 to explain why a recipient checking a seal never needs to borrow the merchant's actual signet ring. Then explain, in the same paragraph, one way the wax seal analogy is *weaker* than real cryptography (hint: see the end of Section 3).

2. Explain, in 4-6 sentences, why "my public key is like my username, my private key is like my password" is a misleading way to think about key pairs, using at least one specific difference from Section 10.

3. List three real-world places you have already trusted a third party (a bank, a website, a government ID office) to vouch for your identity, and for each one, briefly explain what public-key cryptography removes the need for in that scenario.

### Medium

4. Section 5 introduced symmetric cryptography's "key distribution problem." Write a short scenario (150-250 words) describing how two people who have never met before and cannot securely exchange a secret in advance could still each get the other's *public* key safely, and explain why receiving a public key over an insecure channel (like a public website) does not create the same danger that receiving a symmetric key over an insecure channel would.

5. Section 6 used a "one-way spinning wheel" analogy for elliptic curve point multiplication. Design your own analogy (not the wheel, and not the wax seal) for the same "easy forward, infeasible backward" relationship, and write 150-200 words explaining where your analogy holds up and where it breaks down, the way Section 3 did for the wax seal.

6. Research (via a web search or documentation) what "curve parameters" actually are for an elliptic curve used in cryptography (you do not need to understand the underlying equation, just what a curve's public parameters describe at a high level: the curve's shape, the generator point, the size of the number space). Write a short explanation (150-250 words) of why every participant needs to agree on the exact same curve parameters for signatures produced under one set of parameters to be verifiable by someone using a different set.

### Hard

7. Section 10 distinguished the quantum computing threat to hashing (mild, from Grover's algorithm) from the quantum computing threat to elliptic curve cryptography (more serious, from Shor's algorithm). Research Shor's algorithm at a conceptual level (no need to understand the quantum mechanics) and write an explanation (250-400 words) of why an algorithm that can efficiently solve certain math problems on a quantum computer poses a fundamentally different kind of risk to ECC than to SHA-256, referencing the elliptic curve discrete logarithm problem from Section 6.

8. Write a design note (300-450 words) arguing for or against this claim: "Since public keys are safe to share, GoChain should let users choose their own public key directly, instead of deriving it from a randomly generated private key." Your note should explain what would go wrong (or why it would actually be fine) if a user could pick any point on the curve as their "public key" without it being derived from a private key they control, referencing the one-way relationship from Sections 6 and 7.

9. Research the difference between a **private key** and a **seed phrase** (previewed briefly for Chapter 39's BIP-39 implementation). Write a short explanation (250-350 words) of how a seed phrase relates to a private key (hint: it is not the same thing directly, but a way of encoding something a private key can be deterministically derived from) and why losing a seed phrase is often described as equivalent to losing every private key it can generate, even future ones not yet derived.
