# Chapter 79: Asymmetric Cryptography — Diffie-Hellman, RSA, and ECC

> **"Two strangers, standing in a crowded room where everyone can hear every word they say, can still agree on a secret that nobody else in the room can figure out — not because anyone whispers, but because of a trick of arithmetic that's easy to do forward and, for all practical purposes, impossible to undo."**

---

## Table of Contents

1. [The Problem, Restated Precisely](#1-the-problem-restated-precisely)
2. [The Key Insight: Trapdoor Functions](#2-the-key-insight-trapdoor-functions)
3. [Diffie-Hellman Key Exchange — The Idea](#3-diffie-hellman-key-exchange--the-idea)
4. [Diffie-Hellman — A Full Worked Numeric Example](#4-diffie-hellman--a-full-worked-numeric-example)
5. [Why Eve Really Can't Compute It](#5-why-eve-really-cant-compute-it)
6. [RSA — Encrypting With a Public Key](#6-rsa--encrypting-with-a-public-key)
7. [RSA — A Full Worked Numeric Example](#7-rsa--a-full-worked-numeric-example)
8. [Why Asymmetric Crypto Is Slow — and What That Means in Practice](#8-why-asymmetric-crypto-is-slow--and-what-that-means-in-practice)
9. [Elliptic Curve Cryptography (ECC) — The Modern Alternative](#9-elliptic-curve-cryptography-ecc--the-modern-alternative)
10. [Public Key, Private Key: Getting the Vocabulary Precise](#10-public-key-private-key-getting-the-vocabulary-precise)
11. [The Remaining Gap](#11-the-remaining-gap)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Production Usage Notes](#13-production-usage-notes)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#summary)

---

## 1. The Problem, Restated Precisely

Chapter 78 ended at what looked like a dead end. Alice and Bob need a shared AES key. They have never met. Every bit they exchange crosses a network Eve can read in full. Symmetric cryptography — same key locks and unlocks — cannot help, because getting the key to the other side *is* the problem, and encrypting the key just moves the same problem up one level.

The way out requires giving up an assumption that felt obvious: that "shared secret" has to mean "identical piece of information that only the two of you know, transmitted secretly." What if, instead, Alice and Bob could each do some *private* computation, exchange only the *results* of those computations in full public view, and then each combine the other's public result with their own private input to arrive at the *same* final value — a value that turns out to be computationally impractical for anyone who only saw the public exchange (Eve) to reconstruct?

That sentence describes something that sounds like it shouldn't be mathematically possible. Section 4 proves, with real numbers you can check by hand, that it is.

---

## 2. The Key Insight: Trapdoor Functions

The mathematical object that makes this work is called a **trapdoor function**: a function that's easy to compute in one direction, and — without one specific extra piece of information (the "trapdoor") — prohibitively expensive to reverse, even though reversing it is not *impossible* in principle, just impractically slow.

**Intuitive analogy.** Mixing two paint colors together is trivial: pour yellow and blue into the same can, stir, get green. Un-mixing that green paint back into pure separate yellow and pure separate blue is, for all practical purposes, something you cannot do — not because there's a law of physics against it, but because there is no efficient procedure; you'd have to somehow separate individual pigment molecules. Where the paint analogy breaks: paint mixing doesn't have a "trapdoor" — no extra piece of information un-mixes it easily. Real trapdoor functions do have one, which is exactly the private key.

**Engineering terminology.** Two specific mathematical trapdoor problems underpin essentially all classical asymmetric cryptography in use today:

- **The discrete logarithm problem**, which powers Diffie-Hellman (Sections 3–4): given a prime `p`, a generator `g`, and the value `g^x mod p`, it's easy to compute `g^x mod p` if you know `x`, but recovering `x` from `g^x mod p` alone is believed to require time exponential in the size of `p`, for large enough `p`.
- **The integer factorization problem**, which powers RSA (Sections 6–7): given two large prime numbers, multiplying them together is trivial; given only their product, finding the two original primes back out is believed to require impractically large amounts of computation, for large enough primes.

Both problems share the same shape: easy forward, hard backward, with no known efficient general shortcut — and decades of intense, well-funded public and academic cryptanalysis attention have failed to find one for well-chosen parameter sizes. That track record, not a mathematical proof of impossibility (neither problem has one), is why these are trusted as the foundation of modern cryptography.

---

## 3. Diffie-Hellman Key Exchange — The Idea

**The naive approach that doesn't work.** Alice could pick a secret key and text it to Bob. Eve reads it. Dead on arrival — this is exactly Chapter 78's unsolved problem restated.

**The real approach.** Alice and Bob first publicly agree on two numbers that don't need to be secret at all: a large prime `p` and a "generator" `g` (a number with a specific mathematical property described in Section 5). Anyone, including Eve, can know `p` and `g` — they're often even fixed, standardized values reused across millions of connections.

Then:

1. Alice picks a **private** number `a`, known only to her, and computes `A = g^a mod p`. She sends `A` — and only `A`, never `a` — to Bob, in public view.
2. Bob picks a **private** number `b`, known only to him, and computes `B = g^b mod p`. He sends `B` — and only `B`, never `b` — to Alice, in public view.
3. Alice takes Bob's public `B` and raises it to her own private `a`: computes `B^a mod p`.
4. Bob takes Alice's public `A` and raises it to his own private `b`: computes `A^b mod p`.

The punchline, which is pure algebra and not magic: `B^a mod p = (g^b)^a mod p = g^(ba) mod p`, and `A^b mod p = (g^a)^b mod p = g^(ab) mod p`. Since `ab = ba`, **these are the same number.** Alice and Bob have both computed `g^(ab) mod p` — a value neither of them could have produced alone, using two private numbers that never once crossed the network — and that value becomes their shared AES key (or, in practice, is fed into a key-derivation function to produce one, as Section 13 notes).

Eve, meanwhile, saw `p`, `g`, `A`, and `B` — everything that was ever transmitted. She did not see `a` or `b`. To compute the shared secret herself, she would need to recover `a` from `A = g^a mod p` (or `b` from `B`) — exactly the discrete logarithm problem from Section 2, believed computationally infeasible for large enough `p`.

---

## 4. Diffie-Hellman — A Full Worked Numeric Example

Real cryptographic Diffie-Hellman uses primes hundreds of digits long. For a worked example you can check by hand, use a small prime — the math is identical in structure, just breakable by brute force at this size (which is exactly why real deployments use enormous primes instead).

**Publicly agreed values (Eve sees these):** prime `p = 23`, generator `g = 5`.

**Alice's side (private `a` never leaves her machine):**

```
Alice picks private key: a = 6
Alice computes:  A = g^a mod p = 5^6 mod 23

5^1 = 5
5^2 = 25 mod 23     = 2
5^3 = 5 * 2 mod 23  = 10
5^4 = 5 * 10 mod 23 = 50 mod 23 = 4
5^5 = 5 * 4 mod 23  = 20
5^6 = 5 * 20 mod 23 = 100 mod 23 = 8

A = 8   <- Alice sends this to Bob, in the open
```

**Bob's side (private `b` never leaves his machine):**

```
Bob picks private key: b = 15
Bob computes:  B = g^b mod p = 5^15 mod 23

(continuing the pattern above)
5^7  = 5 * 8  mod 23 = 40  mod 23 = 17
5^8  = 5 * 17 mod 23 = 85  mod 23 = 16
5^9  = 5 * 16 mod 23 = 80  mod 23 = 11
5^10 = 5 * 11 mod 23 = 55  mod 23 = 9
5^11 = 5 * 9  mod 23 = 45  mod 23 = 22
5^12 = 5 * 22 mod 23 = 110 mod 23 = 18
5^13 = 5 * 18 mod 23 = 90  mod 23 = 21
5^14 = 5 * 21 mod 23 = 105 mod 23 = 13
5^15 = 5 * 13 mod 23 = 65  mod 23 = 19

B = 19  <- Bob sends this to Alice, in the open
```

**Eve, watching the whole exchange, now knows:** `p = 23`, `g = 5`, `A = 8`, `B = 19`. That's it. She does not know `a = 6` or `b = 15`.

**Alice computes the shared secret** using Bob's public value and her own private key:

```
shared_secret = B^a mod p = 19^6 mod 23

19^2 = 361 mod 23 = 16      (23 * 15 = 345, 361 - 345 = 16)
19^3 = 19 * 16 mod 23 = 304 mod 23 = 5   (23 * 13 = 299, 304 - 299 = 5)
19^6 = (19^3)^2 mod 23 = 5^2 mod 23 = 25 mod 23 = 2

shared_secret = 2
```

**Bob computes the shared secret** using Alice's public value and his own private key:

```
shared_secret = A^b mod p = 8^15 mod 23

8^2  = 64 mod 23 = 18
8^4  = 18^2 mod 23 = 324 mod 23 = 2     (23 * 14 = 322, 324 - 322 = 2)
8^8  = 2^2 mod 23 = 4
8^15 = 8^8 * 8^4 * 8^2 * 8^1 mod 23
     = 4 * 2 * 18 * 8 mod 23
     = (4*2=8) -> 8*18=144 mod 23 = 6  (23*6=138, 144-138=6)
     -> 6*8 = 48 mod 23 = 2

shared_secret = 2
```

**Both arrive at 2.** Alice and Bob now share the secret value `2` (in a real deployment, a much larger number, subsequently fed through a key-derivation function to produce, say, a 256-bit AES key) — and they arrived at it having only ever exchanged `p`, `g`, `A`, and `B` in full public view. Eve saw every single one of those four numbers and still cannot produce `2` without solving a discrete logarithm — recovering `a` from `8 = 5^a mod 23` or `b` from `19 = 5^b mod 23` — which, at this tiny scale, could actually be brute-forced by trying all 22 possible exponents in a fraction of a second. That brute-forceability is entirely a function of the example's small size, not a flaw in the method: real Diffie-Hellman uses primes with 2048+ bits (over 600 decimal digits), where the identical algebraic trick holds, but the discrete logarithm becomes intractable for any known algorithm running on any conceivable hardware.

---

## 5. Why Eve Really Can't Compute It

It's worth being precise about exactly what Eve is stuck on, because it's a genuinely different obstacle from "the numbers are just really big."

Given `p` and `g` and the public value `A = g^a mod p`, computing `A` from `a` (Alice's actual computation) takes a number of multiplications proportional to the number of *bits* in `a` — a fast, efficient operation called **modular exponentiation**, doable even for enormous numbers via the repeated-squaring technique used implicitly in Section 4's calculations. But going the other way — given `g`, `p`, and `A`, finding `a` — has no known algorithm that's similarly efficient. The best known classical algorithms (index calculus and its refinements) still take time that grows *sub-exponentially but still explosively* with the size of `p`, which is precisely why choosing `p` large enough (again, 2048 bits or more in real deployments) pushes the attacker's required computation time past any practical limit — millennia on all of Earth's combined computing power, for well-chosen parameters.

This asymmetry — modular exponentiation is fast forward, the discrete logarithm is (believed) slow backward — is not proven mathematically to be permanently unbreakable the way the one-time pad from Chapter 78 was proven. It rests on the empirical, decades-long failure of the world's cryptography research community to find a fast general algorithm for the discrete logarithm problem on classical computers. (A notable, honest caveat: Shor's algorithm, running on a sufficiently large fault-tolerant quantum computer — which does not exist at meaningful scale today — would solve the discrete logarithm problem, and RSA's factorization problem, efficiently. This is why "post-quantum cryptography" is an active area of standardization, out of scope for this chapter but worth knowing exists.)

---

## 6. RSA — Encrypting With a Public Key

Diffie-Hellman solves "how do two parties agree on a shared secret in public view." RSA (named for its inventors Rivest, Shamir, and Adleman, 1977) solves a related but distinct problem: **let anyone encrypt a message that only one specific party can decrypt** — without that anyone needing to have exchanged any prior secret with the recipient at all.

**The idea.** Bob generates a **key pair**: a **public key** he can hand out to literally anyone, and a **private key** he never shares with anyone. Anyone who has Bob's public key can encrypt a message such that only Bob's private key can decrypt it. Crucially, having the public key does not help you decrypt — nor does it let you derive the private key, for the same "hard in the reverse direction" reason as Section 2's trapdoor functions.

RSA's trapdoor rests on the integer factorization problem: Bob's public key is built from the *product* of two large secret prime numbers; multiplying two primes together is trivial, but factoring that product back into the original two primes, without already knowing them, is believed computationally infeasible for large enough primes. Bob's private key is derived using knowledge of the two original primes — knowledge nobody else has.

---

## 7. RSA — A Full Worked Numeric Example

Real RSA uses primes hundreds of digits long; this toy example uses tiny primes so every step can be verified by hand (and, like Section 4's example, is trivially breakable at this size — that's expected and fine for illustration).

**Key generation (Bob does this once, privately):**

```
1. Pick two prime numbers:           p = 5,  q = 11
2. Compute their product (public):   n = p * q = 55
3. Compute Euler's totient (secret): phi(n) = (p-1)(q-1) = 4 * 10 = 40
4. Pick a public exponent e, coprime with phi(n):  e = 3
   (gcd(3, 40) = 1, so this is a valid choice)
5. Compute the private exponent d, such that e * d = 1 mod phi(n):
   3 * d = 1 mod 40  ->  d = 27
   (check: 3 * 27 = 81 = 2*40 + 1, so 81 mod 40 = 1  correct)

Bob's PUBLIC key:  (e=3, n=55)   <- Bob publishes this to anyone
Bob's PRIVATE key: (d=27, n=55) <- Bob keeps this secret forever
```

**Alice encrypts a message using only Bob's public key.** Say the message is the number `m = 2` (real RSA encrypts padded blocks of bytes, not single small integers, but the arithmetic is identical in structure):

```
ciphertext c = m^e mod n = 2^3 mod 55 = 8 mod 55 = 8

Alice sends c = 8 to Bob, in the open. Eve sees c = 8, e = 3, and n = 55.
```

**Bob decrypts using his private key, which nobody else has:**

```
recovered_m = c^d mod n = 8^27 mod 55

8^1  = 8
8^2  = 64 mod 55 = 9
8^4  = 9^2 mod 55 = 81 mod 55 = 26
8^8  = 26^2 mod 55 = 676 mod 55 = 16    (55*12=660, 676-660=16)
8^16 = 16^2 mod 55 = 256 mod 55 = 36    (55*4=220, 256-220=36)

27 = 16 + 8 + 2 + 1, so:
8^27 = 8^16 * 8^8 * 8^2 * 8^1 mod 55
     = 36 * 16 * 9 * 8 mod 55
     = (36*16 = 576 mod 55 = 26)         (55*10=550, 576-550=26)
     -> (26*9 = 234 mod 55 = 14)          (55*4=220, 234-220=14)
     -> (14*8 = 112 mod 55 = 2)           (55*2=110, 112-110=2)

recovered_m = 2
```

**Bob recovers `m = 2`**, exactly Alice's original message, using only his private key `d = 27`. Eve, who saw `c = 8`, `e = 3`, and `n = 55` — everything that ever crossed the network — cannot decrypt without either recovering `d` (which requires factoring `n = 55` back into `p = 5` and `q = 11`, trivial at this toy size, believed infeasible for the 2048-bit-plus values `n` takes in real deployments) or otherwise breaking the underlying mathematical assumption.

---

## 8. Why Asymmetric Crypto Is Slow — and What That Means in Practice

Chapter 78, Section 9 already flagged the punchline; this section explains the mechanism behind it precisely enough to make the trade-off concrete.

Both Diffie-Hellman and RSA fundamentally rely on **modular exponentiation of very large numbers** — computing something like `g^a mod p` where `g`, `a`, and `p` might each be thousands of bits long. Even using the efficient repeated-squaring technique demonstrated (at tiny scale) in Sections 4 and 7, this requires many multiplications of huge numbers, each of which is itself a nontrivial operation for a CPU — nothing like AES's single-cycle hardware-accelerated table lookups and XORs from Chapter 78, Section 9. Concretely, a single RSA-2048 private-key operation (like a decryption or a signature — foreshadowing Chapter 80) typically takes on the order of a millisecond on modern hardware, versus gigabytes-per-second throughput for AES-GCM — a gap of many orders of magnitude per byte processed.

This gap dictates the architecture of essentially every real secure protocol, and is worth stating as the chapter's central practical conclusion: **asymmetric cryptography is never used to encrypt bulk application data.** It is used briefly, once per connection (or once per some longer-lived session), purely to solve the key-distribution problem from Chapter 78 — either via Diffie-Hellman (agree on a shared secret directly, as in Section 4) or via RSA (one side generates a random symmetric key and encrypts it with the other's public key so only they can read it). Once that brief, expensive asymmetric handshake produces a shared symmetric key, the connection switches entirely to fast symmetric AES (or ChaCha20) for every subsequent byte of actual data. This exact two-phase structure — expensive asymmetric handshake, then fast symmetric bulk transfer — is precisely what Chapter 82's TLS handshake implements, and now you know why it's shaped that way before you've even seen the handshake's message names.

---

## 9. Elliptic Curve Cryptography (ECC) — The Modern Alternative

RSA and classic (finite-field) Diffie-Hellman both work, but they require increasingly large keys to stay ahead of improving factorization and discrete-logarithm algorithms — 2048-bit or 3072-bit RSA keys are standard today, and those large keys make the already-slow modular exponentiation from Section 8 even slower, and make the keys themselves bulky to store and transmit.

**Intuitive level.** Elliptic Curve Cryptography replaces the "numbers modulo a large prime" trapdoor with a different mathematical structure: points on a specific kind of curve, combined using a geometrically-defined "addition" operation. Just as with Diffie-Hellman, there's an operation that's easy to do forward (repeatedly "adding" a point to itself — called **scalar multiplication**) and believed hard to reverse (given the starting point and the result, find out how many times it was added — the **elliptic curve discrete logarithm problem**). Same shape of trapdoor as Section 2, different underlying mathematics.

**Why it matters practically.** The elliptic curve discrete logarithm problem is believed to be significantly harder, key-size for key-size, than the classic discrete logarithm or factorization problems. A 256-bit elliptic curve key is estimated to provide roughly the same real-world security as a 3072-bit RSA key. That means ECC-based key exchange (commonly **ECDHE — Elliptic Curve Diffie-Hellman Ephemeral**, the dominant key exchange in modern TLS, previewed for Chapter 82) and ECC-based signatures (**ECDSA** or the newer **Ed25519**, previewed for Chapter 80) run faster, need smaller keys to transmit, and consume less bandwidth and CPU than RSA at equivalent security — which matters enormously at the scale of a server handling millions of TLS handshakes per second, or a low-power IoT device with limited CPU and battery.

**Deep technical note, honestly scoped.** This chapter deliberately does not derive elliptic curve point addition geometrically or algebraically — that requires more abstract algebra than fits this course's arc, and understanding *that* ECC exists, *why* it's smaller/faster than RSA for equivalent security, and *where* it shows up (ECDHE for key exchange, ECDSA/Ed25519 for signatures) is what matters for the networking chapters ahead. The core conceptual lesson carries over unchanged from Diffie-Hellman: a different trapdoor function, same overall shape — cheap forward, believed prohibitively expensive backward.

---

## 10. Public Key, Private Key: Getting the Vocabulary Precise

Worth pinning down precisely, since "asymmetric cryptography" quietly covers two related-but-distinct use cases that are easy to blur together:

| Use case | Who has which key | What it achieves |
|---|---|---|
| Key exchange (Diffie-Hellman / ECDHE) | Both parties generate a private/public pair; both public values are exchanged | Two parties compute the same shared secret in public view, with neither party's private value ever transmitted |
| Public-key encryption (RSA) | Recipient publishes a public key; keeps a private key secret | Anyone can encrypt a message only the private-key holder can decrypt |
| Digital signatures (RSA/ECDSA — Chapter 80) | Signer keeps a private key secret; publishes a public key | Anyone with the public key can verify a message was signed by the private-key holder, without being able to forge new signatures |

Notice signatures use the key pair in the *opposite* direction from encryption: encryption uses the recipient's public key to lock a message only their private key opens; signatures use the signer's private key to produce something anyone can verify with their public key. Chapter 80 builds this out in full.

---

## 11. The Remaining Gap

Diffie-Hellman and RSA both solve real, hard problems — but notice what neither of them has actually established yet. In Section 4's worked example, Alice computed a shared secret with *whoever sent her the value `B = 19`*. Nothing in the math verifies that the entity on the other end really is Bob, and not Eve, actively impersonating Bob by intercepting Alice's message and substituting her own public value (exactly the active on-path attacker from Chapter 77, Section 4). Diffie-Hellman, run by itself, is vulnerable to exactly this **man-in-the-middle attack**: Eve can run the entire protocol twice, once with Alice (pretending to be Bob) and once with Bob (pretending to be Alice), ending up with two separate shared secrets and silently relaying (and reading, and modifying) everything in between.

The math is not broken — the math correctly does exactly what it claims: produce a shared secret between whoever ran the two halves of the exchange. What's missing is a way to *bind* a public key to a real-world identity — proof that "this specific public key belongs to Bob," not just "here is a public key from whoever's on the other end of this connection right now." That identity-binding problem is exactly what Chapter 80's digital signatures start to solve mechanically, and what Chapter 81's Certificate Authorities solve at Internet scale.

---

## 12. Common Misconceptions

**"Public key" means "the encrypted data," and "private key" means "the decryption password."** No — public and private key are both cryptographic key material, mathematically linked; the public key is meant to be shared openly, the private key never is. Neither is the ciphertext itself.

**"Asymmetric crypto is 'more secure' than symmetric crypto because it uses two keys."** Not a meaningful comparison — they solve different problems (Chapter 78's fast bulk confidentiality vs. this chapter's key-agreement/identity problem) and are used together, not as competing alternatives, in essentially every real protocol.

**"Diffie-Hellman gives you authentication."** Corrected explicitly in Section 11 — plain Diffie-Hellman gives you a shared secret with *whoever answered*, with zero guarantee about who that is. Authentication is a separate property, added on top (Chapters 80–82).

**"RSA and Diffie-Hellman are basically the same thing."** They solve related but distinct problems (Section 6 vs. Section 3) and rest on different hard problems (factorization vs. discrete logarithm), even though both are "asymmetric cryptography" and both are vulnerable to the exact same class of MITM attack described in Section 11 if used without authentication.

**"Bigger keys are always better."** Bigger keys cost more CPU and bandwidth for the same operation (Section 8); the goal is choosing a key size that matches the threat model's realistic attacker capability and desired lifetime of security, not maximizing size unconditionally — which is part of why ECC's smaller-but-equally-strong keys (Section 9) have displaced RSA in many modern deployments.

---

## 13. Production Usage Notes

Real TLS deployments almost always use **ECDHE** (ECC-based, *ephemeral* Diffie-Hellman — a fresh key pair generated per connection, not reused) rather than RSA for key exchange, specifically because ephemeral keys provide **forward secrecy**: even if a server's long-term private key is later stolen, past recorded traffic (whose per-session keys were never stored anywhere, only computed transiently) cannot be decrypted retroactively. RSA key exchange (where the client encrypts a symmetric key directly with the server's long-lived public key) lacks this property — a compromised server key can, in principle, decrypt every past session ever recorded — which is a major reason TLS 1.3 (Chapter 82) removed RSA key exchange entirely and mandates ephemeral Diffie-Hellman-family exchanges. The Diffie-Hellman shared secret from Section 4, in real deployments, is never used as the AES key directly — it's fed through a **key derivation function (KDF)**, along with other handshake data, to produce the actual symmetric session keys, partly to avoid subtle mathematical structure in the raw DH output leaking into the key material.

---

## 14. Interview Questions & Model Answers

**Beginner: "What problem does Diffie-Hellman solve that symmetric cryptography cannot?"**

Symmetric cryptography requires both parties to already possess the same secret key, but provides no mechanism for two parties with no prior relationship to agree on that key over a network an attacker is watching. Diffie-Hellman lets two parties each combine a private number with public values exchanged in the open, arriving at an identical shared secret that an eavesdropper who saw every exchanged value still cannot compute, because doing so requires solving the discrete logarithm problem, believed computationally infeasible for well-chosen parameter sizes.

**Intermediate: "Walk through, at a high level, what happens in RSA encryption and why the private key can decrypt but the public key cannot."**

RSA key generation multiplies two large secret primes together to form the public modulus `n`; the public exponent `e` and modulus `n` together form the public key, while the private exponent `d` (derivable only from the original two primes) forms the private key. Encrypting raises the message to the power `e` modulo `n`; decrypting raises the ciphertext to the power `d` modulo `n`, and the mathematical relationship between `e` and `d` (built from Euler's totient of `n`) guarantees this recovers the original message. Only someone who knows the original two primes can compute `d`; recovering `d` from the public `(e, n)` alone requires factoring `n`, which is computationally infeasible for large enough primes — this is why the public key can encrypt but not decrypt.

**Advanced: "Why does modern TLS use ephemeral Diffie-Hellman (ECDHE) rather than RSA for key exchange, even though RSA can also be used to establish a shared symmetric key?"**

RSA key exchange has the client encrypt a randomly generated symmetric key directly with the server's long-lived public key; if that private key is ever compromised — even years later — an attacker who recorded the encrypted traffic at the time can retroactively decrypt every past session that used it. Ephemeral Diffie-Hellman generates a brand-new key pair per session that's discarded after use; the session's shared secret depends on private values that never existed anywhere durable, so compromising the server's long-term identity key later provides no way to recover past session keys. This property, forward secrecy, is why TLS 1.3 dropped RSA key exchange entirely and mandates (EC)DHE-family exchanges, at the cost of the extra CPU work of generating a fresh ephemeral key pair (or performing a fresh scalar multiplication for ECDHE) on every handshake.

---

## 15. Exercises

### Easy

1. In your own words, explain the difference between what Diffie-Hellman achieves and what RSA achieves.
2. Redo the Diffie-Hellman worked example from Section 4 with `p = 23`, `g = 5`, but change Alice's private key to `a = 4`. Compute Alice's public value `A`, and verify by hand that it's different from the `A = 8` computed in the chapter.
3. Why is a trapdoor function described as "easy forward, hard backward," and why is that exact property necessary for both Diffie-Hellman and RSA?

### Medium

4. Using the RSA key pair generated in Section 7 (`n = 55`, `e = 3`, `d = 27`), encrypt the message `m = 4` and then decrypt your resulting ciphertext to confirm you recover `4`. Show your modular arithmetic.
5. Explain precisely how a man-in-the-middle attacker could defeat a plain (unauthenticated) Diffie-Hellman exchange between Alice and Bob, referencing Chapter 77's definition of an active on-path attacker. What specific extra information would Alice need to detect the attack?
6. Explain, using Section 8's throughput numbers, why a server handling 50,000 new TLS connections per second could not perform an RSA-2048 private-key operation for every byte of every response, and why the hybrid asymmetric-then-symmetric handshake design avoids this entirely.

### Hard

7. Section 4's example is breakable by brute force (only 22 possible exponents to try for `p = 23`). Write out, in plain terms (no need to actually compute it), the brute-force procedure Eve would follow to recover Alice's private key `a` from the public values `p = 23`, `g = 5`, `A = 8`, and explain concretely why this same brute-force procedure becomes impossible in practice once `p` has 2048 bits instead of 5.
8. Research and explain, in a few sentences, what "forward secrecy" means in the context of TLS, why RSA key exchange lacks it, and why ephemeral Diffie-Hellman (ECDHE) provides it — connect your answer to Section 13's production usage notes and to Chapter 78 Section 6's caveat about quantum computers and Grover's/Shor's algorithms.

---

## Summary

| Term | Meaning |
|---|---|
| Trapdoor function | Easy to compute forward, believed infeasible to reverse without extra secret information |
| Discrete logarithm problem | Given `g^x mod p`, finding `x` is believed computationally infeasible for large `p` |
| Diffie-Hellman key exchange | Two parties compute an identical shared secret from public values, without ever transmitting their private inputs |
| Integer factorization problem | Given the product of two large primes, finding the primes is believed computationally infeasible |
| RSA | Public-key cryptosystem: encrypt with a public key, decrypt only with the matching private key |
| Public key / private key | A mathematically linked pair; public is shared openly, private is never shared |
| ECC | Elliptic Curve Cryptography; same trapdoor shape as DH/RSA, smaller keys, faster, based on elliptic curve point arithmetic |
| ECDHE | Ephemeral, ECC-based Diffie-Hellman; dominant TLS key exchange, provides forward secrecy |
| Forward secrecy | Compromising a long-term key later cannot decrypt past recorded sessions |
| Hybrid encryption | Use slow asymmetric crypto once to exchange a key, then fast symmetric crypto for bulk data |

Asymmetric cryptography solves the key-distribution cliffhanger from Chapter 78, but it opens a new gap: it proves you share a secret with *whoever answered*, not that the answerer is really who they claim to be. Chapter 80 closes half of that gap with hashing and digital signatures — proving a message truly came from a specific keyholder, unaltered.
