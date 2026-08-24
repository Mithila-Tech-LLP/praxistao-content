# Chapter 80: Hashing and Digital Signatures

> **"A fingerprint doesn't tell you someone's name. But if you already know whose finger it came from, it tells you, beyond reasonable doubt, that it's really them — and that nothing about the print has been altered since it was taken."**

---

## Table of Contents

1. [The Gap Left Open by Chapter 79](#1-the-gap-left-open-by-chapter-79)
2. [The Problem: Proving a Message Wasn't Altered](#2-the-problem-proving-a-message-wasnt-altered)
3. [A Naive First Attempt: Just Compare the Whole Message](#3-a-naive-first-attempt-just-compare-the-whole-message)
4. [Cryptographic Hash Functions — The Real Solution](#4-cryptographic-hash-functions--the-real-solution)
5. [The Properties a Hash Function Must Have](#5-the-properties-a-hash-function-must-have)
6. [SHA-256: A Real Example](#6-sha-256-a-real-example)
7. [Hashing Alone Isn't Enough: The Missing Piece](#7-hashing-alone-isnt-enough-the-missing-piece)
8. [Digital Signatures — Combining Hashing and Asymmetric Crypto](#8-digital-signatures--combining-hashing-and-asymmetric-crypto)
9. [A Full Worked Example, Step by Step](#9-a-full-worked-example-step-by-step)
10. [What a Signature Proves — and What It Doesn't](#10-what-a-signature-proves--and-what-it-doesnt)
11. [Where This Shows Up in Networking](#11-where-this-shows-up-in-networking)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Production Usage Notes](#13-production-usage-notes)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#summary)

---

## 1. The Gap Left Open by Chapter 79

Chapter 79 ended on a specific, named weakness: Diffie-Hellman lets Alice compute a shared secret with *whoever answered her message* — but nothing in that math checks whether the answerer is really Bob or an active on-path attacker (Chapter 77's Eve) impersonating him. RSA has the same gap in a different form: anyone can encrypt a message with a public key claiming to belong to Bob, but nothing stops Eve from generating her own key pair, publishing her public key while claiming "this is Bob's," and reading everything Alice sends "to Bob."

Both gaps come down to the same missing capability: **a way to prove that a specific message genuinely came from a specific keyholder, and that it hasn't been altered since.** That capability is called a digital signature, and it's built from two ingredients — one entirely new (cryptographic hashing, this chapter's first half) and one you already have (asymmetric cryptography from Chapter 79).

---

## 2. The Problem: Proving a Message Wasn't Altered

State the problem independently of signatures first, because hashing solves a problem on its own, before it's ever combined with anything else: **given a large piece of data, how do you produce a small, fixed-size value that changes completely if even one bit of the data changes — so that comparing two small values tells you, with overwhelming confidence, whether two large pieces of data are identical, without ever comparing the large data directly?**

This shows up everywhere, independent of security: verifying a downloaded file wasn't corrupted in transit (Chapter 19's error-detection checksums were a primitive version of this idea), detecting duplicate files without comparing them byte-by-byte, or building the fast lookup structures (hash tables) used throughout computer science. What makes a hash function *cryptographic* — the subject of this chapter — is a much stronger set of guarantees needed specifically because an adversary, not just random noise, might be trying to produce a collision on purpose.

---

## 3. A Naive First Attempt: Just Compare the Whole Message

The obvious solution to "did this message change" is to keep a full copy of the original and compare byte-by-byte. This technically works but fails the moment the message is large (comparing gigabytes twice for every verification) or, more importantly for security, fails to give you anything to publish or transmit that lets *someone else* verify integrity — sending the "answer key" (a full copy of the original) alongside the message is redundant and doesn't prove anything an attacker who intercepted both couldn't also fake.

Chapter 19's error-detection tools (parity, checksums, CRC) are a real, useful, more compact version of this idea — but they were designed to catch *accidental* corruption from noise (Chapter 17), not to resist a deliberate, intelligent adversary. A CRC is a simple, fast, linear function; it is straightforward, given a target CRC value, to *construct* a different message that produces the exact same CRC on purpose. That's fine for catching random bit-flips from a noisy cable — and disastrous if an attacker wants to forge a document with the same CRC as a legitimate one. The next section builds a hash function strong enough to resist exactly that kind of deliberate attack.

---

## 4. Cryptographic Hash Functions — The Real Solution

**Intuitive level.** A cryptographic hash function is a mathematical blender: feed it any input — a single character or a 10-gigabyte video file — and it always outputs a fixed-size string of bits (256 bits, for the SHA-256 function this chapter focuses on) that looks completely random and bears no visible resemblance to the input. Run the exact same input through it a million times, and you get the exact same output every time (it's deterministic, not actually random) — but there is no known way to work backward from the output to reconstruct, or even partially guess, the input.

**Engineering terminology.** A cryptographic hash function `H` maps an arbitrary-length input to a fixed-length output (the **digest** or **hash**), and it is a *one-way function*: computing `H(input)` is fast and easy; given only `H(input)`, finding any `input'` such that `H(input') = H(input)` is believed computationally infeasible. It requires no secret key at all — unlike Chapter 78's AES or Chapter 79's RSA, hashing is a public, keyless operation that anyone can perform and anyone can verify.

**Deep technical view.** Modern cryptographic hash functions like SHA-256 (part of the SHA-2 family, designed by the NSA and standardized by NIST) work by processing input in fixed-size blocks (512 bits for SHA-256) through a compression function applied repeatedly — a **Merkle-Damgård construction** — where each block's processing depends on the running state accumulated from every previous block. This is exactly what produces the property from Section 1's opening quote: a single flipped bit anywhere in a large input completely changes the internal state from that point forward, cascading into an entirely different final digest (the same avalanche effect concept introduced for AES in Chapter 78, Section 6, applied here to hashing).

---

## 5. The Properties a Hash Function Must Have

A cryptographic hash function is judged against three specific properties, each defending against a specific kind of attack:

- **Preimage resistance (one-wayness).** Given a hash output `h`, it should be computationally infeasible to find *any* input `m` such that `H(m) = h`. This is what makes hashing safe for storing password verification data (Section 13) — even if the stored hash leaks, recovering the original password from it should be infeasible.
- **Second preimage resistance.** Given a specific input `m1`, it should be computationally infeasible to find a *different* input `m2` such that `H(m1) = H(m2)`. This defends against an attacker trying to substitute a fraudulent document for a specific known legitimate one while keeping the same hash.
- **Collision resistance.** It should be computationally infeasible to find *any* two different inputs `m1` and `m2` (neither specified in advance) such that `H(m1) = H(m2)`. This is a subtly *stronger* requirement than second preimage resistance — the attacker gets to freely choose both inputs, not just attack one fixed target — and it's the property that most often breaks first as hash functions age (see the MD5 and SHA-1 note in Section 12).

Why "one-wayness and collision-resistance are the whole point," as the chapter's brief puts it: if preimage resistance failed, an attacker could reconstruct secret inputs (like passwords) from leaked hashes. If collision resistance failed, an attacker could produce two different documents — say, a fair contract and a fraudulent one — that hash to the identical value, then get the fair one signed (Section 8) and later swap in the fraudulent one, since the signature only ever certifies the hash, not the document text directly (this exact mechanism is why collision resistance matters so much for digital signatures specifically).

**The birthday problem, briefly.** Because a hash's output space is finite (2^256 possible values for SHA-256), *some* collision must exist in principle — the question is only how hard it is to find one. The birthday paradox means finding *any* collision (not matching a specific target) takes roughly the square root of the total output space — about 2^128 attempts for a 256-bit hash, rather than the full 2^256 — which is still astronomically infeasible with any known hardware, but is the reason hash digests are sized at 256+ bits rather than, say, 128, for long-term collision resistance.

---

## 6. SHA-256: A Real Example

SHA-256 is the dominant cryptographic hash function in networking today — it's what TLS certificate signatures (Chapter 81), Bitcoin's proof-of-work, and Git's older object-addressing scheme are built on. Its output is always exactly 256 bits, conventionally written as 64 hexadecimal characters. Here are real, verified SHA-256 outputs, computed directly (not simulated) with the `shasum` command-line tool:

```
$ printf '' | shasum -a 256
e3b0c44298fc1c149afbf4c8996fb92427ae41e4649b934ca495991b7852b85

$ printf 'hello' | shasum -a 256
2cf24dba5fb0a30e26e83b2ac5b9e29e1b161e5c1fa7425e73043362938b9824

$ printf 'Hello' | shasum -a 256
185f8db32271fe25f561a6fc938b2e264306ec304eda518007d1764826381969

$ printf 'hello world' | shasum -a 256
b94d27b9934d3e08a52e52d7da7dabfac484efe37a5380ee9088f7ace2efcde9
```

Look closely at `hello` vs. `Hello` — a single-character case change (one bit's worth of difference in the ASCII encoding: `h` is `0x68`, `H` is `0x48`, differing only in one bit) produces two hashes that share no discernible pattern at all: `2cf24dba...` versus `185f8db3...`. This is the avalanche effect from Section 4, made concrete and directly checkable. Also notice: the empty string still hashes to a specific, well-defined 256-bit value (`e3b0c442...`) — a cryptographic hash function is defined for *every* input, including the trivial one, and this specific value is famous enough in security engineering to be immediately recognizable to most practitioners as "the hash of nothing."

**Hands-on experiment.** Run `printf 'hello' | shasum -a 256` yourself (or `sha256sum` on Linux, or `Get-FileHash` in PowerShell, or `hashlib.sha256(b"hello").hexdigest()` in Python). Then flip one character and rerun it. Confirm you get a completely different digest every time, and that running the same input twice always produces the identical digest — determinism and unpredictability existing together is the entire trick.

---

## 7. Hashing Alone Isn't Enough: The Missing Piece

A hash by itself proves only one thing: *this exact data, unmodified, hashes to this exact value.* It says nothing about *who* produced the data. If Alice publishes a document alongside its SHA-256 hash, and Eve intercepts both in transit, Eve can simply replace the document with her own fraudulent version *and* recompute a new matching hash — since hashing requires no secret at all, anyone can compute a valid hash for anything. A bare hash defends against *accidental* corruption (like Chapter 19's CRC, but far more robustly) and detects *any* tampering, but provides zero protection against a deliberate, active-on-path attacker who can just recompute a fresh, correct hash for their tampered version.

What's missing is exactly what Chapter 79's asymmetric cryptography provides: a way to bind an operation to a specific private key that only one party holds, so that anyone can *verify* the binding without being able to *forge* it.

---

## 8. Digital Signatures — Combining Hashing and Asymmetric Crypto

**The idea.** Instead of signing an entire (potentially huge) message directly with a slow asymmetric operation (Chapter 79, Section 8's speed problem), first hash the message down to a small, fixed-size digest, then apply the signer's *private* key to that digest. Anyone with the signer's *public* key can then verify the signature against the digest.

Step by step, using the vocabulary from Chapters 78–79:

1. **Signing (only the private-key holder can do this):**
   - Compute `h = H(message)` — the SHA-256 hash of the full message.
   - Compute `signature = Sign(h, private_key)` — a private-key operation over the digest (mechanically, for RSA, this looks structurally similar to RSA decryption from Chapter 79, Section 7, applied to the hash instead of a plaintext message; ECDSA and Ed25519 use different but analogous private-key operations over the curve from Chapter 79, Section 9).
   - Publish or send `(message, signature)` together.

2. **Verifying (anyone with the public key can do this):**
   - Recompute `h' = H(message)` independently from the received message.
   - Use the signer's public key to check whether `signature` is a valid signature over `h'` — mechanically, for RSA, this looks like the RSA encryption operation from Chapter 79, Section 7, applied to the signature, checking whether it recovers `h'`.
   - If the recovered/verified value matches `h'` exactly, the signature is valid: the message really was signed by the private-key holder, and it has not been altered since (because even a single-bit change to the message would produce a completely different `h'`, per Section 5's collision resistance, that would no longer match the signature).

This is the two ingredients combined precisely: hashing provides the compact, tamper-evident fingerprint (Sections 4–6); asymmetric cryptography provides the "only one party could have produced this, but everyone can check it" property (Chapter 79). Neither ingredient alone would work — a bare hash can be freely recomputed by anyone (Section 7); a raw asymmetric signature over an entire large message would be prohibitively slow (Chapter 79, Section 8) and, for some signature schemes, mathematically awkward to apply directly to arbitrary-length data.

---

## 9. A Full Worked Example, Step by Step

Walk through a concrete (still simplified, but mechanically accurate) example using RSA-style signing, reusing the toy RSA key pair from Chapter 79, Section 7 (`n = 55`, public exponent `e = 3`, private exponent `d = 27`) purely to keep the arithmetic checkable — real signatures use 2048-bit-plus keys and never sign a raw hash value this small without additional padding schemes, noted honestly in Section 13.

**Setup.** Bob wants to publish a message and prove he really wrote it.

```
Message: "PAY ALICE 100"
```

**Step 1 — Bob hashes the message.** In a real system this is the full 256-bit SHA-256 digest; for this toy walkthrough, imagine the hash has been compressed down to fit RSA's tiny `n = 55` modulus from the earlier example, giving a toy digest value:

```
h = H("PAY ALICE 100") = 2   (toy stand-in value, for arithmetic compatible with n=55)
```

**Step 2 — Bob signs the digest with his private key** (`d = 27`, `n = 55` — the same private key from Chapter 79, Section 7):

```
signature = h^d mod n = 2^27 mod 55

(same modular exponentiation structure as Chapter 79's examples,
 using the same private exponent d=27 and modulus n=55 -- signing
 is mechanically RSA "decryption" applied to the digest instead of a
 ciphertext)

2^1  = 2
2^2  = 4
2^4  = 16
2^8  = 256 mod 55 = 36        (55*4=220, 256-220=36)
2^16 = 36^2 mod 55 = 1296 mod 55 = 31   (55*23=1265, 1296-1265=31)

27 = 16 + 8 + 2 + 1, so:
2^27 = 2^16 * 2^8 * 2^2 * 2^1 mod 55
     = 31 * 36 * 4 * 2 mod 55
     = (31*36 = 1116 mod 55 = 16)      (55*20=1100, 1116-1100=16)
     -> (16*4 = 64 mod 55 = 9)
     -> (9*2 = 18)

signature = 18
```

Only Bob can produce this value, because only Bob knows the private exponent `d = 27`.

**Step 3 — Bob publishes the message and the signature together:**

```
("PAY ALICE 100", signature=18)
```

**Step 4 — Alice (or anyone) verifies, using only Bob's public key** (`e = 3`, `n = 55`):

```
1. Alice independently computes h' = H("PAY ALICE 100") = 2  (same hash function, same message)
2. Alice computes: signature^e mod n = 18^3 mod 55

   18^2 = 324 mod 55 = 49       (55*5=275, 324-275=49)
   18^3 = 18 * 49 mod 55 = 882 mod 55 = 2   (55*16=880, 882-880=2)

3. Result is 2, which equals h' = 2 -- the signature is valid.
```

**What this proves to Alice.** Only Bob's private key `d = 27` could have produced a value that, when raised to Bob's public exponent `e = 3` modulo `n = 55`, recovers exactly `h' = 2` — the hash of the exact message Alice received. If even one character of the message had been altered in transit (say, "PAY ALICE 100" became "PAY ALICE 900"), Alice's independently recomputed `h'` would be a completely different value (Section 6's avalanche effect), and it would not match what the signature verifies to — the signature check would fail loudly, and Alice would know the message was tampered with, without needing any separate integrity mechanism at all.

---

## 10. What a Signature Proves — and What It Doesn't

Precisely, using Chapter 77's vocabulary: a valid digital signature proves **integrity** (the message hasn't been altered since signing — guaranteed by the hash's collision resistance) and **authenticity** (the message was signed by whoever holds the corresponding private key — guaranteed by the asymmetric signing operation), together with **non-repudiation** (the signer cannot later credibly claim they didn't sign it, assuming their private key wasn't stolen — a property with real legal weight for signed contracts and code-signing).

What a signature does **not** prove, and this is the exact gap Chapter 81 exists to close: it does not prove *who the keyholder actually is* in the real world. Verifying Bob's signature with "Bob's public key" only means something if you already know, with confidence, that this specific public key genuinely belongs to Bob and not to Eve pretending to be Bob. Nothing in the mathematics of Sections 8–9 establishes that binding — it has to come from somewhere else entirely.

---

## 11. Where This Shows Up in Networking

Digital signatures and hashing appear throughout the stack, well beyond the TLS certificates this volume is building toward:

- **TLS certificates (Chapter 81, Chapter 82):** a Certificate Authority signs a hash of a certificate's contents, binding a public key to a domain name.
- **DNSSEC (Chapter 69):** DNS zones are signed so resolvers can verify records weren't forged in transit or by a compromised cache.
- **Software and firmware updates:** package managers and OS updaters verify a cryptographic signature before installing an update, so a compromised download mirror or an active on-path attacker can't silently substitute malware.
- **BGP route origin validation / RPKI (Chapter 52):** route announcements can be cryptographically signed to prove the announcing AS is actually authorized to originate that prefix — the same identity-binding problem as Chapter 81, applied to routing instead of TLS.
- **Git commit signing:** each commit's content is hashed (historically SHA-1, now migrating to SHA-256 for exactly the collision-resistance reasons in Section 5), and commits can optionally be GPG-signed to prove authorship.

---

## 12. Common Misconceptions

**"Hashing is a type of encryption."** No — encryption is reversible with the right key (Chapters 78–79); hashing is deliberately one-way and has no key or decryption operation at all. "Hashed" data cannot be "unhashed."

**"MD5 or SHA-1 are fine for security purposes because they're still 'hashing.'** Both have had practical, real-world collision attacks demonstrated (a SHA-1 collision was publicly demonstrated in 2017; MD5 collisions have been practical since the mid-2000s) — meaning their collision-resistance property (Section 5) is broken. They remain fine for *non-security* purposes like basic deduplication or checksumming against accidental corruption, but should never be used for signatures, certificates, or password storage. SHA-256 (and newer, SHA-3/BLAKE2/BLAKE3) are the current safe choices.

**"A digital signature is the same as a scanned or typed 'signature' on a document."** A cryptographic signature is a mathematical value computed over the specific bytes of a specific message; it's tied to the exact content and breaks completely (fails to verify) if a single byte changes. A scanned handwritten signature image can be copy-pasted onto any document with no such binding.

**"Verifying a signature tells you the signer is trustworthy."** It tells you the message came from whoever holds a specific private key, unaltered. It says nothing about whether that keyholder is honest, or — critically, per Section 10 — whether the public key you used for verification really belongs to who you think it does.

**"Hashing passwords with plain SHA-256 is secure password storage."** SHA-256 is *fast* by design (Section 6, and Chapter 78's speed emphasis applies to hash computation too) — which is exactly the wrong property for password storage, because it lets an attacker who steals a database of hashes try billions of guesses per second. Real password storage uses deliberately slow, memory-hard functions (bcrypt, scrypt, Argon2) instead, precisely to defeat that kind of brute-force guessing.

---

## 13. Production Usage Notes

Real digital signature schemes never sign a raw hash value directly the way Section 9's toy example did for simplicity — RSA signatures use a padding scheme (commonly PSS — Probabilistic Signature Scheme) that adds structured randomness before the modular exponentiation, both for provable security reasons and to prevent certain algebraic attacks that exist against naively "textbook" RSA signing. ECDSA and the newer Ed25519 (Chapter 79, Section 9's ECC) are increasingly preferred over RSA for signatures in modern protocols (TLS 1.3, SSH, and most new deployments), for the same reasons ECDHE is preferred for key exchange: smaller keys and signatures, comparable security, and faster computation. SHA-256 itself is standardized in FIPS 180-4 and is a required baseline for TLS 1.2 and 1.3 certificate signatures.

---

## 14. Interview Questions & Model Answers

**Beginner: "What is a cryptographic hash function, and what makes it different from a regular checksum like CRC?"**

A cryptographic hash function maps arbitrary-length input to a fixed-length digest such that it's computationally infeasible to reverse it (find an input producing a given output) or to find two different inputs producing the same output. A checksum like CRC (Chapter 19) is designed to catch accidental, random corruption efficiently, but is not designed to resist a deliberate adversary — it's often trivial to construct two different inputs with the same CRC on purpose, which makes CRC unsuitable for any application where an attacker might benefit from a collision.

**Intermediate: "Explain how a digital signature is created and verified, and why it needs both hashing and asymmetric cryptography."**

To sign, the signer hashes the message to a fixed-size digest, then applies their private key to that digest using an asymmetric algorithm (RSA, ECDSA, or Ed25519), producing the signature. To verify, anyone recomputes the hash of the received message independently and uses the signer's public key to check that the signature is valid for that hash. Hashing alone can't prove authorship — anyone can compute a hash, so an attacker could tamper with a message and just recompute a fresh, correct hash to match. Asymmetric cryptography alone would be too slow to apply directly to large messages and, for some signature schemes, isn't mathematically well-suited to arbitrary-length data. Combining them gets a fast, fixed-size representation of the message (the hash) bound to a specific private key (the signature).

**Advanced: "A digital signature verifies successfully. What exactly has been proven, and what critical assumption does that proof rest on that isn't part of the cryptographic math itself?"**

A successful verification proves the message matches exactly what was signed (integrity, via the hash's collision resistance) and that the signature was produced using the private key corresponding to the public key used for verification (authenticity, via the asymmetric algorithm's unforgeability). It does not, by itself, prove anything about who controls that private key in the real world — the verifier has to separately trust that the specific public key genuinely belongs to the entity they believe it does. That binding between a cryptographic public key and a real-world identity is not established by signature math at all; it requires an external trust mechanism, which is exactly the gap Public Key Infrastructure and Certificate Authorities are built to close.

---

## 15. Exercises

### Easy

1. Compute `printf 'networking' | shasum -a 256` (or the equivalent in your language of choice) and record the output. Then change one letter and recompute. Confirm the outputs share no visible similarity.
2. Explain in your own words why a hash function needs no secret key, while encryption and signing both do.
3. Name the three properties a cryptographic hash function must have (Section 5), and give one real-world consequence of each property failing.

### Medium

4. Explain why signing the hash of a message is preferable to signing the entire message directly, referencing Chapter 79 Section 8's discussion of asymmetric crypto's speed.
5. A colleague proposes storing user passwords as `SHA-256(password)` in a database. Explain, using Section 12's misconception correction, what's wrong with this approach and what should be used instead.
6. Using Section 9's worked example structure, explain what would happen to the verification step if an attacker intercepted Bob's message and changed "PAY ALICE 100" to "PAY ALICE 900" without also having access to Bob's private key. Walk through exactly where the verification fails.

### Hard

7. Research the 2017 "SHAttered" SHA-1 collision demonstration (or a similar documented hash collision attack). Explain, in terms of Section 5's three properties, exactly which property was broken, and why a successful *collision* (rather than a *preimage*) attack is specifically dangerous for digital signatures (tie your answer to Section 5's contract-substitution example).
8. Explain precisely why Section 10 says a valid signature does not prove the identity of the signer in the real world, using a concrete attack scenario: Eve generates her own RSA key pair and claims, out loud, "this public key belongs to Bob." Walk through what would happen if Alice believed her and verified a message using Eve's key, and identify exactly which chapter's mechanism is missing that would have prevented this.

---

## Summary

| Term | Meaning |
|---|---|
| Cryptographic hash function | Maps arbitrary-length input to a fixed-length digest; deterministic, one-way, keyless |
| Digest / hash | The fixed-size output of a hash function |
| Preimage resistance | Can't recover an input from its hash output |
| Collision resistance | Can't find two different inputs with the same hash output |
| SHA-256 | Widely deployed 256-bit cryptographic hash function (SHA-2 family) |
| Avalanche effect (hashing) | A tiny input change produces a completely different digest |
| Digital signature | Private-key operation over a message's hash, verifiable by anyone with the public key |
| Non-repudiation | Signer cannot credibly deny having signed a message |
| RSA-PSS / ECDSA / Ed25519 | Real-world signature schemes built on RSA or elliptic curves |

Digital signatures prove a message came from a specific private key, unaltered — but a verifier still has to trust that the public key genuinely belongs to the real-world identity it claims to. Chapter 81 closes that final gap: Public Key Infrastructure and Certificate Authorities, the system that lets your browser decide which signatures — and which keyholders' claims to be "google.com" — to actually believe.
