# Chapter 12: Digital Signatures with ECDSA

Chapter 11 ended with a promise: a private key can produce something called a **signature**, and anyone holding only the matching public key can check that signature without ever touching the private key itself. That promise was left deliberately abstract, because "produce a signature" is not one universal operation — it is a specific mathematical recipe, and different recipes exist. This chapter names GoChain's recipe (ECDSA), and, just as Chapter 11 explained elliptic curves conceptually before Chapter 13 wrote a single line of curve-math code, this chapter diagrams exactly what happens during signing and verifying, in plain language, before any Go appears. By the end, Chapter 13's `Sign()` and `Verify()` functions should feel like a direct, obvious translation of what you already understand, not a leap into new territory.

## Table of Contents

1. [What a Signature Actually Has to Prove](#1-what-a-signature-actually-has-to-prove)
2. [Meet ECDSA](#2-meet-ecdsa)
3. [The Curve Underneath: secp256k1](#3-the-curve-underneath-secp256k1)
4. [Signing, Conceptually](#4-signing-conceptually)
5. [Why You Sign the Hash, Not the Raw Message](#5-why-you-sign-the-hash-not-the-raw-message)
6. [Verifying, Conceptually](#6-verifying-conceptually)
7. [The Full Sign/Verify Pipeline, End to End](#7-the-full-signverify-pipeline-end-to-end)
8. [Signature Size and What Travels With a Transaction](#8-signature-size-and-what-travels-with-a-transaction)
9. [What Happens When Even One Byte Changes](#9-what-happens-when-even-one-byte-changes)
10. [The Nonce: ECDSA's Sharpest Edge](#10-the-nonce-ecdsas-sharpest-edge)
11. [Common Misconceptions About Signatures](#11-common-misconceptions-about-signatures)
12. [Where This Fits in GoChain](#12-where-this-fits-in-gochain)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. What a Signature Actually Has to Prove

Before naming an algorithm, it is worth being precise about the job that algorithm has to do, because "prove Alice approved this" is actually two separate claims bundled together, and a real signature scheme has to deliver both at once:

1. **Authenticity** — this exact data was approved by whoever controls this specific private key, and nobody else could have produced this signature.
2. **Integrity** — this exact data has not been altered since it was signed. Change even one byte, and the signature that was valid a moment ago becomes invalid.

Notice that a signature is never "attached" to a private key in the abstract — it is always a signature *over a specific piece of data*. Alice's private key does not have one single, reusable signature the way a rubber stamp has one fixed imprint. Every new piece of data she signs produces a brand-new signature, tied to that data and that data alone. This is the single most important difference between a cryptographic signature and Chapter 11's wax seal analogy, and it is worth stating plainly now so it does not cause confusion later: a wax seal looks the same on every envelope; a digital signature looks completely different for every message, even messages signed by the same person with the same private key.

```
  Wax seal (Chapter 11):                Digital signature (this chapter):

  Same ring, same imprint,               Same private key, but a DIFFERENT
  on every envelope.                     signature for every different
                                          piece of data signed with it.

  ┌────────┐   ┌────────┐               "pay Bob 5"  ──▶  signature X
  │ envelope│  │ envelope│               "pay Bob 6"  ──▶  signature Y (totally
  │  [seal] │  │  [seal] │                                 different from X)
  └────────┘   └────────┘
   identical      identical
```

This "one signature per exact message" property is not a limitation — it is exactly what makes integrity checking possible. If Alice's signature were reusable across any message, Bob could take a signature Alice produced for "pay Bob 5 gochips" and reattach it to "pay Bob 500 gochips," and it would still "look" like Alice's mark. Section 4 makes concrete why that attack does not work.

---

## 2. Meet ECDSA

**ECDSA** stands for **Elliptic Curve Digital Signature Algorithm**. It is not a new idea competing with Chapter 11's elliptic curve concepts — it is the specific, standardized *procedure* for using an elliptic curve key pair to produce and check signatures. Chapter 11 explained that a private key is a number and a public key is a point on a curve, derived one-way from that number. ECDSA is the recipe that takes those two pieces (plus the data to be signed) and turns them into a signature, and the matching recipe that takes a signature, a public key, and the data, and turns them into a yes-or-no answer.

ECDSA is not exotic or GoChain-specific — it is the exact signature algorithm Bitcoin has used since its very first block in 2009, and the one Ethereum still uses for every transaction today. Learning it here means you are learning the actual algorithm securing trillions of dollars of real value, not a simplified stand-in invented for this course.

The name is also a useful landmark if you ever go looking for more depth than this course provides. ECDSA is the elliptic-curve variant of an older algorithm, **DSA** (the Digital Signature Algorithm), standardized by the U.S. National Institute of Standards and Technology (NIST) in the early 1990s using ordinary modular arithmetic rather than elliptic curves. ECDSA takes the exact same overall shape as DSA — the same "sign produces a pair of numbers, verify checks a relationship between them" structure this section is about to describe — and simply swaps in elliptic curve point-multiplication (Chapter 11, Section 6) wherever DSA used plain modular exponentiation. Both are published, formally standardized algorithms (ECDSA appears in NIST's FIPS 186 standard, among others), not a home-grown invention by any single company or, for that matter, by Bitcoin's own creator — one more reason a stranger's laptop, running GoChain code it has never seen before, can trust the math without trusting the messenger.

A finished ECDSA signature is not one number — it is a *pair* of numbers, conventionally called **r** and **s**. You do not need to memorize what r and s individually represent to use ECDSA correctly (Go's standard library computes both internally), but it is worth knowing the shape exists, because Section 10 and Chapter 33 both return to exactly this pair:

```
  ECDSA signature = (r, s)

  r  ── derived from a random value chosen fresh for this one signature
  s  ── derived from r, the private key, and the hash of the signed data

  Both numbers travel together as "the signature." Verifying checks
  a mathematical relationship between r, s, the public key, and the
  hash — Section 6 diagrams this without the underlying algebra.
```

---

## 3. The Curve Underneath: secp256k1

Chapter 11, Section 6 explained that elliptic curve cryptography needs everyone to agree, in advance, on one specific curve and one specific generator point. Bitcoin and Ethereum both standardized on a particular curve named **secp256k1** — a name that, decoded, just describes its technical parameters (a 256-bit curve from a standardized family, the "k1" identifying which specific variant). There is nothing mysterious about secp256k1 itself; it is one of several curves cryptographers consider well-studied and secure, and Bitcoin's original design simply chose it.

Here is a detail worth being honest about before Chapter 13 writes any code: Go's standard library (`crypto/elliptic`) ships built-in support for a *different* family of curves — P-224, **P-256**, P-384, and P-521 — but does not include secp256k1 out of the box. Using secp256k1 in Go requires a third-party package. GoChain's goal in this volume is to teach ECDSA correctly using only Go's standard library, so **GoChain uses P-256** (also called secp256r1) instead of Bitcoin's secp256k1.

This is a real, deliberate trade-off, not an oversight, and it is worth being precise about what it costs and what it does not:

```
  secp256k1 (Bitcoin, Ethereum)          P-256 (GoChain, this course)

  - Not in Go's standard library         - Built into crypto/elliptic
  - Requires a third-party package       - Zero extra dependencies
  - Same "easy forward, infeasible       - Same "easy forward, infeasible
    backward" security shape from          backward" security shape from
    Chapter 11                             Chapter 11
  - Widely used in production            - Widely used in production too
    cryptocurrency                         (TLS certificates, many other
                                            real systems, just not Bitcoin)
```

Nothing about *how* ECDSA works changes based on which curve sits underneath it — Sections 4 through 8 apply identically to secp256k1 or P-256. The only thing that changes is which specific curve's point-multiplication rule gets used, and GoChain picks the one that lets this entire course run on the Go standard library alone, with no external dependencies to install, version, or trust. If you later want GoChain to be wire-compatible with real Bitcoin signatures, swapping in a secp256k1 package is a drop-in change to `NewKeyPair()` alone — everything else in this chapter and the next stays exactly the same.

---

## 4. Signing, Conceptually

Signing takes three inputs — a private key, the data to be signed, and a source of fresh randomness — and produces one output: the signature pair `(r, s)` from Section 2.

```
                    ┌───────────────────┐
   private key  ───▶│                    │
   data         ───▶│   ECDSA signing    │───▶  signature (r, s)
   fresh randomness ▶│                    │
                    └───────────────────┘
```

Conceptually, here is the shape of what happens inside that box, without the underlying algebra (Chapter 13's `crypto/ecdsa` package handles the real arithmetic correctly, exactly as `crypto/sha256` handled SHA-256's internals in Chapter 09):

1. A fresh random number (sometimes called a **nonce**, short for "number used once") is generated, just for this one signature.
2. That random number is combined with the curve's generator point (Chapter 11, Section 6) to produce the `r` component.
3. The data being signed is hashed (Section 5 explains why this step matters enormously), and that hash is combined mathematically with `r`, the private key, and the random nonce to produce the `s` component.
4. The pair `(r, s)` together *is* the signature — it gets attached to the data and sent along with it.

The private key is used *during* this process, but it never leaves the computer that ran it, and it cannot be recovered from the resulting `(r, s)` pair — the same one-way guarantee Chapter 11 built the entire chapter around, now applied to signature production specifically. Anyone who intercepts a signature learns nothing about the private key that made it.

---

## 5. Why You Sign the Hash, Not the Raw Message

Step 3 above quietly did something worth pulling out and examining on its own: it hashed the data before signing it, rather than feeding the raw message directly into the signing math. This is not an implementation detail — it is essential to how ECDSA works at all, for two concrete reasons.

**First, a practical one.** ECDSA's underlying math expects a fixed-size number as input, not an arbitrarily long message. A transaction might be a few hundred bytes; a large file might be gigabytes. Chapter 08 already built exactly the tool for turning "any amount of data" into "one fixed-size fingerprint": `crypto.Hash`. Signing `crypto.Hash(data)` instead of `data` itself means ECDSA always receives the same fixed-size input (32 bytes, for SHA-256) no matter how large or small the original message was.

**Second, and more important, a security one.** Recall the **avalanche effect** from Chapter 08, Section 3: changing a single bit of input scrambles a SHA-256 hash completely and unpredictably. Signing the hash rather than the raw message means that Section 1's integrity guarantee — "change one byte, and the old signature stops validating" — is inherited directly from a property you already trust, rather than needing ECDSA to somehow detect tiny data changes on its own:

```
  "pay Bob 5 gochips"   ──Hash──▶  3f2a91c...   ──sign──▶  signature A

  "pay Bob 6 gochips"   ──Hash──▶  d817ec4...   ──sign──▶  signature B

  One character changed the message. The avalanche effect (Ch. 08)
  made the hash completely different. A completely different hash
  produces a completely different, unrelated signature. Signature A
  will NOT verify against "pay Bob 6 gochips" -- Section 9 shows why.
```

This is exactly the kind of layered design this course keeps returning to: ECDSA does not need to reinvent tamper-detection, because it is built directly on top of a hash function that already provides it.

---

## 6. Verifying, Conceptually

Verifying is the mirror image of signing, but notice carefully what goes in and what does not: the private key is completely absent.

```
                    ┌───────────────────┐
   public key   ───▶│                    │
   data         ───▶│   ECDSA verifying  │───▶  true  (signature is valid)
   signature (r,s)──▶│                    │───▶  false (signature is invalid)
                    └───────────────────┘
```

Conceptually: the verifier hashes the data themselves (using the exact same hash function used during signing — SHA-256, in GoChain's case), then performs a mathematical check involving that hash, the public key, and the `(r, s)` pair from the signature. If a specific equation holds, the signature is valid; if it does not, the signature is invalid. There is no ambiguous "probably valid" outcome — the check is a clean, deterministic yes or no, and any two honest computers running this check on the same inputs always agree, exactly the "independently verifiable, no coordination needed" property Chapter 11, Section 1 demanded from the very start.

Crucially, this check can only ever succeed for data that was actually signed (in Section 4's sense) by the private key matching the public key being used to verify. There is no way to construct a valid `(r, s)` pair for arbitrary data without access to the private key — that is precisely the elliptic curve discrete logarithm problem (Chapter 11, Section 6) working in the verifier's favor. Anyone can run this check; nobody except the private key's owner can pass it for new data.

---

## 7. The Full Sign/Verify Pipeline, End to End

Putting Sections 4 and 6 together gives the complete picture Chapter 13 will implement in Go, function for function:

```
   ALICE (has the private key)              ANYONE (has only the public key)
   ───────────────────────────              ────────────────────────────────

   "pay Bob 5 gochips"
          │
          ▼
   Hash("pay Bob 5 gochips")
          │
          ▼
   Sign(privateKey, hash)  ──▶ (r, s)
          │
          ▼
   send data + (r, s)  ──────────────────▶  receive data + (r, s)
                                                       │
                                                       ▼
                                            Hash(received data)
                                                       │
                                                       ▼
                                            Verify(publicKey, hash, (r, s))
                                                       │
                                              ┌────────┴────────┐
                                              ▼                 ▼
                                             true              false
                                        "genuinely           "reject --
                                         Alice's, data        forged or
                                         unaltered"           corrupted"
```

Every arrow in the right-hand column can be run by anyone: Bob, a stranger, every single node on the GoChain network independently, all reaching the identical `true` or `false` answer with zero coordination between them. That is the entire point Chapter 11, Section 1 opened this volume with, now made concrete as an actual, runnable pipeline.

---

## 8. Signature Size and What Travels With a Transaction

It is worth pausing on a very concrete, practical question before moving on: how big is one of these `(r, s)` signatures, actually, in bytes? For P-256 (Chapter 13's curve, Section 3), each of `r` and `s` is a number up to 256 bits — 32 bytes — long. A raw signature is therefore about 64 bytes, and Go's standard encoding of the pair (using a compact wrapper format called ASN.1 DER, which Chapter 13 uses directly via `ecdsa.SignASN1`) adds a small amount of overhead on top, typically landing a P-256 signature somewhere around 70-72 bytes total. secp256k1, being the same 256-bit size class, produces essentially the same size signature.

This number matters more than it might seem, and it connects directly back to Chapter 11, Section 8's argument for elliptic curves over RSA in the first place. Every transaction GoChain builds starting in Volume 5 carries at least one signature alongside it, and every signed transaction gets stored, forever, in a block, and every block gets copied to every node on the network (Volume 7). A roughly 70-byte signature per transaction is small enough that this cost stays manageable even as GoChain's chain grows to thousands or millions of transactions — exactly the "smaller signatures compound into real savings at scale" argument Chapter 11 made about ECC broadly, now attached to one specific, concrete number instead of a general claim:

```
  One transaction, roughly:

  ┌─────────────────────────────────────────────┐
  │  transaction data (sender, recipient,        │
  │  amount, etc.)              ~100-300 bytes    │
  ├─────────────────────────────────────────────┤
  │  ECDSA signature (r, s)      ~70 bytes        │
  ├─────────────────────────────────────────────┤
  │  public key (Section 3,      ~64 bytes        │
  │  Chapter 11)                                  │
  └─────────────────────────────────────────────┘

  Multiply by thousands of transactions per block,
  and thousands of blocks over the chain's lifetime --
  a few extra bytes per signature is not academic.
```

This is also a preview of a genuine, practical question Chapter 14 answers in full: a 64-byte public key is unwieldy to write down, read aloud, or type into a wallet by hand, which is exactly the problem addresses are built to solve.

---

## 9. What Happens When Even One Byte Changes

It is worth tracing through, step by step, exactly why corrupting even a single byte of signed data causes verification to fail — this is precisely what Chapter 13's worked example demonstrates in real Go code, and seeing the mechanism here first makes that demonstration feel inevitable rather than surprising.

Suppose Alice signs `"pay Bob 5 gochips"` and produces a valid signature `(r, s)`. Now suppose an attacker intercepts the message and changes it to `"pay Bob 9 gochips"`, hoping to reuse Alice's genuine signature on the modified text:

```
   Original:   "pay Bob 5 gochips"  ──Hash──▶  hash_A
   Tampered:   "pay Bob 9 gochips"  ──Hash──▶  hash_B   (completely different
                                                          from hash_A, by the
                                                          avalanche effect)

   Verify(Alice's public key, hash_B, (r, s))

   The signature (r, s) was mathematically produced FOR hash_A.
   The verification equation checks a relationship that only holds
   between (r, s) and the exact hash it was generated for. Checking
   it against hash_B instead fails the equation outright.

   Result: false. Verification rejects the tampered message.
```

The attacker cannot "patch" the signature to match the new hash either, because doing so would require solving for a new valid `(r, s)` pair without knowing Alice's private key — exactly the elliptic curve discrete logarithm problem, again, standing in the way. This is the concrete mechanism behind both of Section 1's promises at once: authenticity (only Alice's private key produces signatures Alice's public key accepts) and integrity (any change to the data breaks the match between hash and signature), delivered by the same single mathematical check.

It is worth being precise about what "changes" actually means here, because it is easy to picture only obvious edits like swapping "5" for "9." The avalanche effect (Chapter 08, Section 3) makes no distinction between a dramatic edit and the smallest possible one: flipping a single bit anywhere in the message — capitalizing one letter, adding one invisible trailing space, changing one digit of an amount by one gochip — scrambles the hash just as completely as rewriting the whole sentence. There is no such thing as a "minor" tamper that a signature might overlook. Either the bytes verified are bit-for-bit identical to the bytes originally signed, or verification fails; there is no partial credit, and no notion of an edit being "close enough."

---

## 10. The Nonce: ECDSA's Sharpest Edge

Section 4 mentioned, almost in passing, that signing needs "a source of fresh randomness" — a random nonce, generated new for every single signature. This detail deserves more attention than a passing mention, because it is the one place ECDSA's security depends on something *outside* the elegant math: getting the randomness right, every single time, without exception.

If the same nonce is ever reused across two different signatures made with the same private key, an attacker who obtains both signatures can combine them mathematically to recover the private key itself — completely bypassing the elliptic curve discrete logarithm problem entirely, because a nonce reuse leaks exactly the information that problem is supposed to protect. This is not a theoretical worry invented for this course: in 2010, a security research group discovered that Sony had used a *fixed*, non-random value for this nonce when signing software for the PlayStation 3, across every single signature Sony ever produced. Reusing the same "random" value every time let researchers directly compute Sony's private signing key from just two of its signatures, permanently and completely compromising a key that could never be safely reused.

```
  Correct: a brand-new random nonce for every signature.

    Sig 1 uses nonce k1  ─┐
    Sig 2 uses nonce k2  ─┤  all different, all secret
    Sig 3 uses nonce k3  ─┘

  Sony's bug: the SAME nonce reused for every signature.

    Sig 1 uses nonce k  ─┐
    Sig 2 uses nonce k  ─┤  identical -- catastrophic
    Sig 3 uses nonce k  ─┘

    Two signatures + the same k = attacker recovers the private
    key directly, with simple algebra, no brute force needed.
```

This is precisely why Chapter 13 leans entirely on Go's `crypto/ecdsa` package to generate this nonce internally, using `crypto/rand` (Go's cryptographically secure random source) under the hood, rather than GoChain ever managing this value by hand. The lesson generalizes past this one nonce: Section 7's beautiful, symmetric math is only as trustworthy as the randomness it depends on, and "properly random, every single time, with no exceptions" is a genuinely hard operational requirement to get right by hand — exactly the kind of detail worth delegating to a well-audited standard library rather than reimplementing.

---

## 11. Common Misconceptions About Signatures

**"A signature encrypts the message."** No. ECDSA signatures do not hide or encrypt anything — the message travels in plain view, alongside its signature, and anyone can read it. A signature only proves who approved it and that it has not changed; if you also wanted to keep the message secret from onlookers, that is a completely different problem (encryption) solved by different tools, not something ECDSA does as a side effect.

**"The same private key produces the same signature every time."** No, and Section 10 explains exactly why this must be true: a fresh random nonce goes into every signature, so signing the exact same message twice with the exact same private key produces two different, equally valid `(r, s)` pairs. Both verify successfully against the same public key. A signature scheme where the output never changed for the same input would leak information an attacker could use, so this variability is a deliberate security feature, not a bug.

**"Verifying a signature proves the data is true."** No — verification only proves that whoever holds a specific private key approved these exact bytes. If Alice signs a message claiming something false, a valid signature just proves *Alice said it*, not that it is accurate. This distinction matters enormously once GoChain reaches transactions in Volume 5: a valid signature proves Alice authorized spending from her own address; it says nothing about whether Alice actually has the funds she is claiming to spend, which is a completely separate check the blockchain itself must perform.

**"If I see a signature, I can tell whose private key made it just by looking."** Not from the signature alone. A signature is only meaningful paired with a specific public key to check it against — `(r, s)` by itself is just two numbers. Verification always requires all three pieces together: the data, the signature, and the public key you are checking it against.

**"A bigger private key always makes for a stronger signature."** Not in isolation. Once a private key is properly sized for its curve (256 bits, for both P-256 and secp256k1) and properly random (Chapter 11, Section 7), making it bigger buys nothing further, because the security of the whole scheme is capped by the curve's own strength, not by how large a number you feed into it. What actually undermines a signature scheme in practice is almost never "the key was a few bits too short" — it is one of the operational failures this chapter and Chapter 11 both catalogued: a leaked private key (Chapter 11, Section 4), a reused nonce (Section 10), or a bug in how the surrounding software calls the underlying, correctly-implemented math.

---

## 12. Where This Fits in GoChain

Every transaction GoChain produces starting in Volume 5 will be signed exactly this way: the transaction's contents get hashed, that hash gets signed with the spender's private key, and every node that receives the transaction independently re-hashes it and verifies the attached signature against the sender's public key before accepting it. Chapter 33 returns to one subtlety this chapter only foreshadowed — signing a carefully *trimmed* copy of a transaction rather than the whole thing, to avoid a signature-malleability bug that ECDSA's `(r, s)` pair can otherwise introduce — but the core mechanism, sign-the-hash and verify-with-the-public-key, is exactly what this chapter just diagrammed. The Chapter 61 virtual machine's `OpCheckSig` instruction, much later in this course, is ultimately just this same verification check, wired up so a smart contract can demand a valid signature before releasing funds — proof that a concept learned this early keeps paying off across the entire rest of the book. Chapter 13 turns every arrow in Section 7's diagram into real, compiling Go code: `crypto.NewKeyPair()`, `crypto.Sign()`, and `crypto.Verify()`.

---

## Summary

- A signature must prove two things at once: **authenticity** (this specific private key's owner approved this) and **integrity** (the data has not changed since). Unlike a wax seal, a digital signature is different for every message, never a single reusable stamp.
- **ECDSA** (Elliptic Curve Digital Signature Algorithm) is the specific, standardized procedure GoChain uses to turn a private key and some data into a signature, and a public key, data, and signature into a true/false verification result.
- An ECDSA signature is a pair of numbers, **r** and **s**, produced using the private key, the data's hash, and a fresh random nonce.
- Bitcoin and Ethereum use the **secp256k1** curve; Go's standard library does not ship it built in, so GoChain deliberately uses **P-256** instead, trading exact Bitcoin compatibility for a zero-dependency, standard-library-only implementation — the same underlying "easy forward, infeasible backward" security shape either way.
- Signing always signs a **hash** of the data, not the raw data itself — this inherits the avalanche effect from Chapter 08, so changing even one byte of the original data produces a completely unrelated hash, which the old signature will no longer validate against.
- The random **nonce** used during signing must be fresh and unpredictable for every signature; reusing it (as Sony did, famously, for the PlayStation 3) lets an attacker recover the private key directly from just two signatures, with no brute force required.
- Verifying never requires the private key — only the public key, the data, and the signature — which is exactly what lets every node on the network check a signature independently, with no coordination.

---

## Exercises

### Easy

1. In your own words, explain why a digital signature is fundamentally different from a wax seal in one specific way: does a private key produce the *same* signature every time it signs, or a *different* one? Reference Section 1 and Section 10 in your answer.

2. Section 5 gave two separate reasons ECDSA signs a hash of the data rather than the raw data itself. Name both reasons in your own words, and say which one is a practical/engineering reason and which one is a security reason.

3. A friend tells you: "ECDSA signatures keep your message secret from anyone who intercepts it." Correct this misconception in 3-4 sentences, using Section 11.

### Medium

4. Section 3 explained that GoChain uses P-256 instead of Bitcoin's secp256k1. Write a short explanation (150-250 words) of what would have to change in `crypto.NewKeyPair()` if a future version of GoChain wanted to switch to secp256k1 for real Bitcoin compatibility, and what would *not* need to change anywhere else in this chapter's sign/verify pipeline.

5. Walk through Section 9's byte-corruption example again, but instead of changing the amount ("5" to "9"), imagine an attacker tries to change the recipient's name instead ("Bob" to "Eve") while keeping Alice's original signature attached. Explain, step by step, using the diagram style from Section 9, exactly where verification fails.

6. Research (via documentation or a web search) what "signature malleability" means for ECDSA specifically — the property that a valid `(r, s)` pair can sometimes be transformed into a second, different, but still-valid `(r, s)` pair for the exact same message and key. Write a 150-250 word explanation of why this could be dangerous for a system that identifies transactions by hashing the entire signed message (including its signature), previewing Chapter 33.

### Hard

7. Research the real 2010 Sony PlayStation 3 ECDSA nonce-reuse incident referenced in Section 10. Write an explanation (250-400 words) of, at a conceptual level (no need for the underlying algebra), why having two signatures produced with the *same* nonce and the *same* private key is enough information to solve for that private key directly, and why this is a fundamentally different kind of failure than a weak or short private key (Chapter 11, Section 7).

8. Design a short test plan (as a numbered list, 200-350 words) for a hypothetical `Verify()` function, covering at minimum: a genuinely valid signature, a signature checked against the wrong public key, a signature checked against tampered data, and a signature with a deliberately corrupted `r` or `s` value. For each case, state what `Verify()` should return and why, referencing Sections 6, 7, and 9.

9. Write a design note (300-450 words) considering this question: since Section 10 established that ECDSA's security depends heavily on generating a correct, unpredictable nonce every single time, would it be safer for a signature scheme to derive the nonce *deterministically* from the private key and the message (so the same message always produces the same nonce, with no randomness involved at all), rather than generating a fresh random one each time? Research RFC 6979 (deterministic ECDSA) at a high level to inform your answer, and explain what problem it solves relative to the failure mode described in Section 9.
