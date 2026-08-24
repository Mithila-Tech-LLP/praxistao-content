# Chapter 08: Hashing and SHA-256

Every blockchain, at its core, is built out of one deceptively simple idea: a function that turns any amount of data into a short, fixed-size fingerprint. This chapter builds a solid, plain-language understanding of what that function does and why it behaves the way it does, using SHA-256 — the exact hash function Bitcoin uses, and the one GoChain will use starting next chapter — as the running example.

## Table of Contents

1. [What Is a Hash Function?](#1-what-is-a-hash-function)
2. [Determinism — Same Input, Same Output, Always](#2-determinism--same-input-same-output-always)
3. [The Avalanche Effect](#3-the-avalanche-effect)
4. [One-Wayness — You Cannot Reverse a Hash](#4-one-wayness--you-cannot-reverse-a-hash)
5. [Collision Resistance](#5-collision-resistance)
6. [How SHA-256 Works, Conceptually](#6-how-sha-256-works-conceptually)
7. [Hashing in Everyday Software](#7-hashing-in-everyday-software)
8. [Why Blockchains Use SHA-256 Specifically](#8-why-blockchains-use-sha-256-specifically)
9. [Hashes as GoChain's Building Block](#9-hashes-as-gochains-building-block)
10. [Common Misconceptions, Addressed](#10-common-misconceptions-addressed)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What Is a Hash Function?

Imagine a machine at a shipping company. You feed it a box of any size — a single letter, a crate of books, an entire truckload of furniture — and no matter what you put in, it prints out a barcode label that is always exactly the same length: say, 12 digits. Two completely different boxes get two completely different barcodes. The *same* box, fed in twice, always gets the *same* barcode. That machine is a **hash function**.

A **hash function** is a piece of code that takes an input of any size — a single character, a short sentence, an entire book, a video file — and produces an output of a fixed size, called a **hash** (also called a **digest** or, in this book, a **fingerprint**). For SHA-256, that fixed size is always 256 bits, no matter how big or small the input was. 256 bits is 32 bytes, which we usually write as 64 hexadecimal characters (each hex character represents 4 bits, and 64 × 4 = 256).

```
                 ┌─────────────────────┐
   "hi"    ───▶  │                     │  ───▶  256-bit fingerprint
                 │   SHA-256 hash      │
 "War and        │      function       │  ───▶  256-bit fingerprint
  Peace"  ───▶   │                     │
                 └─────────────────────┘

   Any size in.                              Always the same size out.
```

Notice what the diagram is showing: a two-character string and an entire novel both go in on the left, and both come out on the right as a fingerprint of the exact same length. The hash function does not care how much data you gave it. It always compresses everything down to that same fixed size.

This single property is already useful for a blockchain. Instead of comparing two entire blocks byte-by-byte to check if they are identical (slow, and awkward to store), you can compare their 256-bit fingerprints instead. If the fingerprints match, the data matches. If even one byte differs anywhere in a multi-gigabyte input, the fingerprint comes out completely different, as you will see in Section 3.

A quick note on how we will *write down* these fingerprints throughout this book. A 256-bit number is enormous and unreadable in raw binary (a string of 256 ones and zeroes). Instead, we almost always write hashes in **hexadecimal** — base 16, using the digits `0-9` and the letters `a-f`, where each single hex character represents exactly 4 bits. This is why a 256-bit (32-byte) SHA-256 output is always written as exactly 64 hex characters: 32 bytes × 2 hex characters per byte = 64. You will see this "64 hex characters" shape constantly from here on — it is the visual signature of a SHA-256 hash. Chapter 09 covers exactly how Go converts raw hash bytes into this readable hex form.

Different hash functions produce different fixed sizes, which is worth knowing so you recognize them later. The table below compares a few you may encounter:

```
  Hash function     Output size      Hex characters     Used by
  ───────────────   ─────────────    ───────────────    ──────────────────
  MD5                128 bits         32 characters      legacy, broken
  SHA-1              160 bits         40 characters      legacy, weakened
  SHA-256            256 bits         64 characters      Bitcoin, GoChain
  SHA-512            512 bits        128 characters      some certificates
  Keccak-256         256 bits         64 characters      Ethereum
```

MD5 and SHA-1 are included here specifically as a warning: both were once considered secure and widely used, but researchers eventually found practical ways to produce collisions for each of them (Section 5 explains what a collision is). Neither should be used anywhere security matters today, including anywhere in GoChain. This is a useful, humbling reminder that "cryptographically secure" is not a permanent label — it is a claim that holds only until someone finds a break, which is exactly why SHA-256's decades-long, unbroken track record (covered in Section 8) matters so much when choosing a hash function for a system meant to protect real value.

---

## 2. Determinism — Same Input, Same Output, Always

**Determinism** means that a function, run on the same input any number of times, always produces the exact same output. There is no randomness, no dependency on the time of day, the computer you are running on, or how many times you have already called it before.

Think of a food scale. If you place the same apple on it ten times in a row, it should read the same weight every single time (ignoring tiny measurement noise a scale might have — a hash function has none of that noise at all; it is perfectly exact). If the scale gave you a different number every time for the same apple, it would be useless for anything serious — you could never trust a reading. A hash function is like a perfect scale: feed it the same bytes, get the same fingerprint, forever, on any computer, in any programming language.

```
  "hello"  ──▶  SHA-256  ──▶  2cf24dba5fb0a30e26e83b2ac5b9e29e
                              1e1b161e5c1fa7425e73043362938b9824

  "hello"  ──▶  SHA-256  ──▶  2cf24dba5fb0a30e26e83b2ac5b9e29e
    (again)                   1e1b161e5c1fa7425e73043362938b9824

  Same input, same output — every single time, on every machine.
```

Determinism is what makes hashing useful for a blockchain's tamper-evidence. If you hash a block today and store that fingerprint, then hash the same block again next year, you should get the identical fingerprint back — assuming nothing in the block changed. If the fingerprint you compute now does *not* match the one stored earlier, that is proof, not a guess, that something in the data changed between then and now. Every node on a real blockchain network can independently recompute the same hash of the same block and arrive at the exact same 256-bit answer, without needing to trust each other or compare notes — that shared, guaranteed-identical computation is a large part of what makes decentralized agreement possible at all.

Picture three friends — Alice, Bob, and Carol — each running their own independent copy of GoChain on their own laptop, in three different cities, having never spoken to each other today. All three receive the same block of transactions. Because SHA-256 is deterministic, all three laptops compute the exact same 256-bit hash for that block, independently, without any of the three needing to phone the other two and compare notes:

```
   Alice's laptop   ──▶  SHA-256(block)  ──▶  7dc91974...30948b58
   Bob's laptop     ──▶  SHA-256(block)  ──▶  7dc91974...30948b58
   Carol's laptop   ──▶  SHA-256(block)  ──▶  7dc91974...30948b58

   Three independent computers. Same input. Identical output —
   guaranteed by determinism, with zero coordination required.
```

This is a small but genuinely important preview of a much bigger idea covered starting in Chapter 23 (consensus): a blockchain network is made of computers that do not inherently trust one another, run by people who may never meet. Determinism is the quiet, unglamorous property that lets all of them agree on basic facts — "this is what block 7dc91974... looks like" — as a simple side effect of running the same well-defined math, rather than something they need a central authority to referee.

---

## 3. The Avalanche Effect

The **avalanche effect** is the property that changing even a single bit of the input completely and unpredictably changes the output — roughly half of the output bits flip, seemingly at random, even though nothing about the change was random at all.

Here is where the shipping-label analogy breaks down a little and a better one takes over: think of an avalanche starting on a mountainside. Kick a single small pebble loose near the top, and by the time the disturbance reaches the bottom of the slope, it is an entirely different, unrecognizable mass of snow, ice, and rock — not "the same avalanche, but slightly smaller." SHA-256 behaves the same way with input bits: nudge one bit, and the fingerprint at the bottom is a completely different-looking string, not "the same hash with one character changed."

Below are two *real* SHA-256 fingerprints, computed from two input strings that differ by exactly one letter — the final word ends in `dog` versus `dof`:

```
Input:  "The quick brown fox jumps over the lazy dog"
SHA-256: d7a8fbb307d7809469ca9abcb0082e4f8d5651e46d3cdb762d02d0bf37c9e59

Input:  "The quick brown fox jumps over the lazy dof"
SHA-256: a1cbac0e93075ab66ad59ff54c32c8abcaeb533f0568e109281ed57eb519685
```

Look closely: not one single hex character in the first fingerprint matches the character in the same position in the second one, aside from a handful of coincidental overlaps you would expect from pure chance. One letter changed in a 44-character sentence, and the entire 64-character fingerprint became unrecognizable. There is no visible "closeness" between the two hashes that would let you guess the inputs were almost identical — and that is exactly the point.

```
   dog  ──▶  SHA-256  ──▶  d7a8fbb3...c9e59
    │                            │
    │ one letter changes         │ every part of the
    ▼                            ▼ fingerprint changes
   dof  ──▶  SHA-256  ──▶  a1cbac0e...19685

   Small, deliberate change in.   Total, unpredictable change out.
```

This matters enormously for a blockchain. If the avalanche effect did not hold — if changing one transaction's amount inside a block only changed the block's hash "a little" — an attacker could search for a tampered version of a block whose hash *looked* plausible enough to slip past a lazy check. Because of the avalanche effect, there is no such thing as a "close enough" hash. A tampered block's hash is not a suspicious variant of the original; it is a totally different, unrelated-looking number, making tampering trivially and immediately obvious to anyone who recomputes it.

Here is a second real example, closer to home for this course, comparing the two words this book's currency is built on:

```
Input:  "gochain"
SHA-256: 94270d9fbb414dbcc537dc89ca5b5b4d7ca2d63732bad19d4d04b924ca6425e

Input:  "gochip"
SHA-256: b84c623ee4489528cd1dfd55c552cc3739df6bf3ade94251282f3c107373aa8
```

`"gochain"` and `"gochip"` share the first four letters (`goch`) and are conceptually related — one is the project, the other is its currency unit — but their fingerprints share no visible structure whatsoever. This is a useful mental check to keep running as you move through this course: no matter how "related" two pieces of data might seem to a human reader, their hashes give you no hint of that relationship at all. A hash of block 1,000 and a hash of block 1,001 look no more alike than a hash of block 1,000 and a hash of a random string typed by a cat walking across a keyboard.

---

## 4. One-Wayness — You Cannot Reverse a Hash

**One-wayness** means that, given only a hash's output, there is no practical way to work backward and recover the original input. You can go from input to output easily and quickly. Going from output back to input is, by design, not feasible — not "hard but doable with enough patience," but computationally infeasible even for the largest computers on Earth, for any input longer than a handful of characters.

An analogy: think about scrambling an egg. Turning a raw egg into scrambled egg is fast and easy — crack it, whisk it, cook it. But there is no process, no matter how clever, that turns the scrambled egg back into the original, intact raw egg with the yolk and white separated the way they started. The transformation only runs in one direction. SHA-256 is deliberately engineered to behave the same way with information: it is trivial to compute `Hash(data)`, and there is no shortcut to compute `data` from `Hash(data)` alone. Note this is different from *encryption*, which is designed to be reversible by anyone holding the right key — hashing has no key, and no reverse operation, at all.

```
  Forward direction (easy, fast, one line of code):

    "my secret password123"  ────▶  ef92b778bafe771e89245b89ecbc08a4

  Reverse direction (the hash function offers no way to do this):

    ef92b778bafe771e89245b89ecbc08a4  ────▶  ???
                                              (no feasible way back)
```

One-wayness is exactly why hashing raw passwords (with additional safeguards beyond plain SHA-256, which is a topic for later in this course) is standard practice: a service can store your password's hash, check a login attempt by hashing what you typed and comparing fingerprints, and never need to store — or be able to recover — your actual password at all. For a blockchain, one-wayness means a block's hash reveals nothing about its contents to someone who has not already seen those contents; the fingerprint proves you *know* the data (because you could reproduce its exact hash) without needing to publish the data alongside it in every context.

It is worth being precise about what one-wayness does *not* protect against, since this is a common point of confusion. One-wayness means there is no mathematical shortcut from output back to input. It says nothing about **guessing**: if the space of possible inputs is small — say, four-digit PIN codes, or a short dictionary word — an attacker does not need to "reverse" the hash at all. They can simply hash every possible PIN (there are only 10,000) or every word in a dictionary, and check which one matches the target hash. This is called a **brute-force** or **dictionary attack**, and it works not by breaking one-wayness, but by outrunning it with sheer exhaustive guessing over a small enough space of likely inputs. This is precisely why real password systems add a random, per-user **salt** (extra random bytes mixed into the password before hashing) and deliberately slow, memory-hard hash functions — techniques this course revisits when GoChain encrypts wallet files in Volume 6 — rather than relying on plain SHA-256 alone. For GoChain's own use of hashing (fingerprinting blocks and transactions, not guessable secrets), this distinction matters less: a block's contents are large and effectively unguessable, so plain SHA-256's one-wayness is exactly the right tool for the job.

---

## 5. Collision Resistance

A **collision** is a case where two *different* inputs happen to produce the exact same hash output. Because a hash function squeezes an unlimited amount of possible input down into a fixed 256-bit output, collisions must mathematically exist somewhere — there are infinitely many possible inputs but only 2²⁵⁶ possible outputs, so by simple counting, some two different inputs somewhere must share a fingerprint.

**Collision resistance** is the property that, even though collisions exist in theory, no one has a practical, faster-than-random-guessing way of *finding* one. Think of it like the number of grains of sand on every beach on Earth: there are certainly two grains somewhere that are, atom for atom, identical in shape and size — but finding that exact matching pair by searching is so wildly impractical that for every purpose that matters, you can treat it as if it never happens.

2²⁵⁶ is an almost incomprehensibly large number — vastly larger than the number of atoms in the observable universe. Even if every computer on Earth searched for a SHA-256 collision continuously, it would take far, far longer than the current age of the universe to have a meaningful chance of finding one. That is what "collision resistant" means in practice for SHA-256 today: not "impossible" in the strictest mathematical sense, but "so implausible that engineering a blockchain around the assumption that it will not happen is entirely reasonable" — exactly the assumption Bitcoin, Ethereum, and GoChain all make.

```
                 ┌─────────────────────────────┐
  Every possible │  input of every length,      │
  input, ever    │  from empty string to        │
                 │  entire libraries of books    │
                 └──────────────┬──────────────┘
                                 │  SHA-256
                                 ▼
                 ┌─────────────────────────────┐
                 │   only 2²⁵⁶ possible          │
                 │   fixed-size outputs          │
                 └─────────────────────────────┘

     Infinitely many inputs squeezed into a finite (but huge) output
     space guarantees collisions exist — but finding one is, for all
     practical purposes, computationally impossible.
```

This is why a blockchain can safely use a single hash as a stand-in for an entire block's contents: the chance that some *other*, different block would ever produce that same hash is close enough to zero to build a financial system on top of.

A common intuition trap worth naming here, which you will explore quantitatively in Exercise 8: collisions become findable *far* sooner than people expect once you are searching among a large group of inputs rather than checking one specific pair — this is known as the **birthday paradox**, named after the surprising fact that in a room of just 23 random people, there is already better than a 50% chance two of them share a birthday, even though there are 365 possible birthdays and only 23 people. Cryptographers account for this "search among many" effect when sizing hash outputs: SHA-256's 256-bit output is chosen to be large enough that even birthday-paradox-style searching across an astronomical number of inputs still leaves collision-finding computationally infeasible.

---

## 6. How SHA-256 Works, Conceptually

**SHA-256** stands for "Secure Hash Algorithm, 256-bit." It is one member of a family of hash function designs published by the U.S. National Institute of Standards and Technology (NIST), and it has been extensively studied by cryptographers worldwide since its publication in 2001 without anyone finding a practical way to break its one-wayness or collision resistance. This course does not implement SHA-256's internal math by hand — that math (bitwise rotations, modular addition, a scrambling schedule applied over 64 rounds) is intricate, easy to get subtly wrong, and already implemented correctly, and extremely fast, in Go's standard library. Instead, this section gives you the conceptual shape of what happens inside, so SHA-256 stops being a mysterious black box and becomes a specific, understandable (if intricate) process.

At a high level, SHA-256 works like this:

1. **Padding.** Your input is padded — extra bits are appended to its end — so that its total length becomes a clean multiple of 512 bits. This is bookkeeping, so the next step always has evenly sized chunks to work with, no matter what length you started with.
2. **Splitting into blocks.** The padded input is cut into 512-bit chunks, one after another.
3. **Compression, one chunk at a time.** SHA-256 keeps a running 256-bit internal state (think of it as a "scratchpad" that starts at a fixed, publicly known set of starting numbers). For each 512-bit chunk, it runs 64 rounds of mixing operations — shifting bits, rotating them, combining them with the scratchpad using addition and logical operations — that thoroughly stir the chunk's bits into the scratchpad.
4. **Chaining.** The scratchpad's state after processing one chunk becomes the starting state for processing the *next* chunk. This is what makes the avalanche effect happen even for a change near the very end of a huge input: every later chunk's processing depends on everything that came before it.
5. **Final digest.** Once every chunk has been processed, the scratchpad's final 256-bit state — after one last round of mixing — is the hash.

```
  Input (any length)
        │
        ▼
  ┌───────────┐
  │  Padding   │  add bits until length is a multiple of 512
  └─────┬─────┘
        ▼
  ┌───────────┬───────────┬───────────┐
  │  chunk 1   │  chunk 2   │  chunk 3   │   ... (512 bits each)
  └─────┬─────┴─────┬─────┴─────┬─────┘
        │            │            │
        ▼            ▼            ▼
   ┌────────┐   ┌────────┐   ┌────────┐
   │ 64      │──▶│ 64      │──▶│ 64      │──▶  final 256-bit digest
   │ rounds  │   │ rounds  │   │ rounds  │
   └────────┘   └────────┘   └────────┘
   scratchpad    scratchpad    scratchpad
   starts at     carried in    carried in
   fixed values  from chunk 1  from chunk 2
```

The important takeaway from this diagram is the chaining arrows between chunks: because each chunk's processing depends on the scratchpad state left behind by the previous chunk, a change to *any* chunk — first, middle, or last — ripples forward and changes the scratchpad state for every chunk after it, which is exactly the mechanism that produces the avalanche effect from Section 3. You do not need to memorize the 64 rounds or the specific bitwise operations used inside them to use SHA-256 correctly and safely — you need to trust (correctly, given decades of public scrutiny) that this process behaves deterministically, one-way, and with strong collision resistance, and let Go's standard library run it for you, which is exactly what Chapter 09 does.

To make the padding step in step 1 concrete with a tiny example: if you hash the three-character string `"abc"`, SHA-256 does not process just those 3 bytes (24 bits) directly. It appends a single `1` bit, then enough `0` bits to make room, then a 64-bit number recording the *original* message length (24, in this case), until the total length is exactly 512 bits — one full chunk. Only then does the 64-round compression process in step 3 run, once, on that single padded chunk. A multi-megabyte block of transactions works exactly the same way, just with many more 512-bit chunks chained together instead of one.

---

## 7. Hashing in Everyday Software

Before returning to blockchains specifically, it is worth seeing that hashing is not some exotic blockchain-only invention — it quietly runs underneath a huge amount of software you already use every day. Recognizing these everyday cases makes the blockchain use case in Section 9 feel like a natural extension of a well-worn idea, rather than something invented from scratch for cryptocurrency.

**Git commit IDs.** Every commit you have ever made in Git is identified by a hash (Git currently uses SHA-1, though it is migrating toward SHA-256 for exactly the collision-resistance reasons discussed in Section 5) computed over the commit's contents — the code changes, the author, the timestamp, and the hash of the *parent* commit. If you have ever seen a commit referred to as `a3f5c9d`, that short string is a prefix of a much longer hash fingerprint. This should look familiar: it is precisely the "each block stores the previous block's hash" chaining pattern this course builds starting in Volume 3, and it is why `git log` shows you an unbroken, tamper-evident history of a codebase.

**Download checksums.** When you download a large file — a Linux distribution image, a piece of software — the publisher often also publishes a SHA-256 checksum alongside it. After downloading, you can hash the file yourself and compare your result against the published one. If they match, you can be confident the download was not corrupted in transit and (assuming the publisher's checksum listing itself was not tampered with) that no attacker replaced the file with a malicious one along the way.

```
  Publisher's website:            Your computer, after downloading:

  ubuntu.iso                       ubuntu.iso  (downloaded copy)
  sha256sum:                       $ shasum -a 256 ubuntu.iso
  d94b5a...c8fa1                   d94b5a...c8fa1

                    Match?  ──▶  Yes  ──▶  download is intact
                             ──▶  No   ──▶  download is corrupted
                                            or was tampered with
```

**Deduplication and change detection.** Cloud storage and backup software (and, as you will build hands-on in Chapter 15, a simple file integrity tool) often hash every file rather than comparing full file contents, to quickly detect duplicates or figure out which files changed since the last backup. Two files with the same hash are (for all practical purposes, per Section 5) the same file; a changed hash means the file's contents changed.

**Hash tables (a different, related idea).** Programming languages' built-in map/dictionary types (Go's `map`, Python's `dict`) use a *different family* of hash functions internally — usually much faster and not designed to resist a determined attacker — to decide which "bucket" a key belongs in for quick lookup. It is worth knowing this distinction exists: not every hash function you will encounter in software is a *cryptographic* hash function with the one-wayness and collision-resistance properties this chapter covers. SHA-256 is deliberately slower and much more carefully designed than a hash table's internal hash function, precisely because it needs to resist a deliberate, well-funded adversary trying to break it — a hash table's hash function only needs to resist accidental, non-malicious collisions.

Every one of these everyday examples relies on the same three properties from Sections 2, 3, and 5: determinism (so the same file always produces the same checksum), the avalanche effect (so a corrupted or tampered file is instantly detectable), and collision resistance (so you can trust that a matching hash really does mean matching content). GoChain's use of hashing, starting in Chapter 09, is this same idea, applied to blocks and transactions instead of files and commits.

---

## 8. Why Blockchains Use SHA-256 Specifically

SHA-256 is not the only cryptographic hash function that exists — SHA-3, BLAKE2, and Keccak-256 (used by Ethereum, which is subtly different from standard SHA-3 despite the similar name) are all reasonable alternatives with similar security properties. GoChain, like Bitcoin, standardizes on SHA-256 for a few concrete, practical reasons:

- **It is battle-tested at enormous scale.** SHA-256 has secured essentially every Bitcoin transaction since 2009 — trillions of dollars of value, under constant attack from the best-funded adversaries on the planet — without a single practical break of its core properties ever being found.
- **It ships in every mainstream language's standard library**, including Go's `crypto/sha256`, with implementations that are audited, fast, and safe to use directly, with no third-party dependency needed.
- **Modern CPUs often have hardware instructions** (Intel's SHA extensions, ARM's cryptographic extensions) specifically accelerating SHA-256, making it fast even at blockchain scale, where a single node might hash thousands of blocks and transactions.
- **It has a large, well-understood security margin.** Cryptographers continue to study it and, decades in, still find no practical weakness — the kind of track record you want underpinning a system meant to hold real value indefinitely.
- **It holds up reasonably well against future quantum computers.** Unlike some public-key cryptography (a topic Chapter 11 covers, and where quantum computers pose a much more serious theoretical threat), a well-known result called Grover's algorithm only halves a hash function's effective security margin — for SHA-256, that means a very well-resourced future quantum computer would face a search roughly as hard as breaking a 128-bit hash today, which is still entirely impractical. This does not mean quantum computing is irrelevant to blockchains generally, only that SHA-256's specific role (fingerprinting data) is far less exposed than the digital-signature role covered later in this volume.

None of these reasons involve SHA-256 being mathematically "the best possible" hash function in some absolute sense — they are pragmatic, real-world reasons rooted in trust built over time and ecosystem support. That pragmatism is itself a lesson worth carrying into the rest of this course: production cryptographic choices usually favor "proven and well-supported" over "theoretically newest."

---

## 9. Hashes as GoChain's Building Block

Every major concept the rest of this course builds rests on this chapter's three properties working together:

- **Determinism** lets every independent node on the GoChain network compute the exact same hash for the exact same block and agree on it, without needing to compare raw data or trust each other's word.
- **The avalanche effect** guarantees that tampering with a block — even changing a single transaction amount by one gochip — is instantly, unmistakably detectable, because the resulting hash looks nothing like the original.
- **One-wayness and collision resistance** together mean a block's hash can safely stand in for the block's entire contents: nobody can forge a different block that happens to produce the same hash, and nobody can work backward from a hash to reconstruct data they were not supposed to have.

```
  Block N-1              Block N               Block N+1
┌───────────┐    hash  ┌───────────┐    hash  ┌───────────┐
│ data       │ ───────▶│ PrevHash   │ ───────▶│ PrevHash   │
│ PrevHash   │         │ data       │         │ data       │
│ Hash       │         │ Hash       │         │ Hash       │
└───────────┘         └───────────┘         └───────────┘

  Each block's Hash is a SHA-256 fingerprint of its own contents.
  Each next block stores that fingerprint as its PrevHash — so
  changing anything in Block N-1 changes its Hash, which breaks
  the PrevHash stored in Block N, which breaks Block N's own Hash,
  which breaks Block N+1's PrevHash... and so on down the chain.
```

This diagram is a preview of exactly what Volume 3 builds in Go. For now, the goal of this chapter was purely conceptual: you should be able to explain, in your own words, what a hash function is, why the same input always gives the same output, why a tiny input change gives a wildly different output, why you cannot reverse a hash, and why finding two different inputs with the same hash is not realistically possible. Chapter 09 turns all of this into real, running Go code — the first function in `gochain/crypto`.

---

## 10. Common Misconceptions, Addressed

Before moving on, it's worth directly naming a few misunderstandings about hashing that are common enough to head off explicitly.

- **"A hash is a kind of encryption."** It is not. Encryption is deliberately reversible by anyone holding the right key — that is the entire point of encrypting something you intend to read again later. A hash has no key and no reverse operation at all (Section 4). If you find yourself thinking "I'll hash this so I can decrypt it later," that is a sign the word you actually want is "encrypt," not "hash."
- **"If two hashes look similar, the inputs were probably similar."** The avalanche effect (Section 3) guarantees the opposite: a hash gives you zero partial credit for a "close" input. Two inputs differing by one character produce two hashes that share no more visible structure than two hashes of completely unrelated data. Never use "the hashes look kind of alike" as evidence of anything.
- **"SHA-256 will eventually be reversible with enough computing power."** One-wayness (Section 4) is not a claim about current computing power being insufficient — it is a claim that there is no known mathematical shortcut from hash back to input at all, for any amount of computing power realistically available now or in the foreseeable future (Section 8 addresses the specific case of quantum computers). More computing power lets you *guess* more inputs per second; it does not give you a way to skip guessing.
- **"A single collision would instantly break Bitcoin/GoChain."** A collision would be a serious, actively exploited problem if a practical way to *find* one were ever discovered (Section 5) — but SHA-256's 256-bit output makes even birthday-paradox-accelerated searching astronomically impractical today. A single accidental collision occurring in the wild, with nobody able to reproduce or exploit it, is a different (and vastly less likely, and less consequential) scenario than a discovered method for generating them on demand.
- **"Hashing a file and hashing a string are fundamentally different operations."** They are not — `Hash` (Chapter 09) takes a `[]byte` either way. A string is just bytes; a file's contents, once read into memory, are just bytes. The same function, the same three properties, apply identically regardless of what the bytes originally represented.

---

## Summary

- A **hash function** takes input of any size and produces a fixed-size fingerprint — SHA-256 always produces 256 bits, shown as 64 hexadecimal characters.
- **Determinism** means the same input always produces the same output, on any machine, forever — this lets every node independently compute and agree on the same fingerprint.
- The **avalanche effect** means a tiny change to the input (even one character) produces a completely different, unrecognizable output, making tampering impossible to hide.
- **One-wayness** means you cannot practically work backward from a hash to recover the original input — hashing is not encryption, and there is no key or reverse operation.
- **Collision resistance** means that, although two different inputs sharing a hash must mathematically exist, finding such a pair is computationally infeasible with current or foreseeable technology.
- SHA-256 works by padding the input, splitting it into 512-bit chunks, and running each chunk through 64 rounds of mixing that chain into the next chunk's starting state — this chaining is what produces the avalanche effect.
- Hashing already underpins everyday software you likely use — Git commit IDs, download checksums, backup deduplication — all leaning on the same determinism, avalanche effect, and collision resistance covered in this chapter.
- Blockchains, including GoChain, use SHA-256 specifically because of its long track record, standard-library availability, hardware acceleration, and large security margin — not because it is theoretically the newest or most exotic option.
- Every later GoChain concept — block linking, tamper-evidence, transaction IDs, Merkle roots, mining — rests directly on the three properties covered in this chapter.

---

## Exercises

### Easy

1. In your own words, write a short paragraph (4-6 sentences) explaining the difference between hashing and encryption to someone who has never heard of either. Your explanation must use the word "reversible" and make clear which one of the two is reversible and which is not, and why that difference matters for storing something like a password.

2. Using any hash generator you can find (a website, a command-line tool like `shasum -a 256` on macOS/Linux, or a small script in any language), compute the SHA-256 hash of your own first name, then compute the SHA-256 hash of your first name with the first letter capitalized differently (e.g., `alice` versus `Alice`). Write down both fingerprints and describe, in your own words, what you observe about how similar or different they look.

3. Explain why a hash function needs to produce a *fixed-size* output regardless of input size, using the shipping-label analogy or one of your own. What practical problem would arise for a blockchain if hashing a huge block produced a huge hash, while hashing a tiny transaction produced a tiny hash?

### Medium

4. The chapter states that 2²⁵⁶ is "vastly larger than the number of atoms in the observable universe" (estimated at roughly 10⁸² atoms). Look up or estimate the value of 2²⁵⁶ in scientific notation (base 10), and write a short comparison showing how many orders of magnitude larger it is than the atom estimate. Explain in your own words why this gap matters for collision resistance in practice.

5. Re-read Section 6's description of SHA-256's chaining step (each 512-bit chunk's processing depends on the scratchpad state left by the previous chunk). Draw your own version of the diagram in Section 6, but for a 3-chunk input where you deliberately change one bit *only* in chunk 2. Annotate your diagram to show which chunks' processing is affected and why chunk 1's processing is untouched while chunk 3's is not.

6. Research (via a web search or documentation) one other cryptographic hash function mentioned in Section 7 — SHA-3, BLAKE2, or Keccak-256 — and write a short comparison (150-250 words) against SHA-256: what is different about its internal design at a conceptual level, and what practical reason might a real project choose it over SHA-256?

### Hard

7. The avalanche effect implies that hashes "look random" even though the process is fully deterministic. Research the difference between "looking random" (statistical randomness, i.e., no detectable pattern) and "being random" (true unpredictability with no underlying determinism), and write an explanation (250-400 words) of why a cryptographic hash function needs the former property but is, underneath, entirely the latter's opposite — a fixed, repeatable algorithm.

8. Suppose SHA-256 had a fixed output size of only 32 bits instead of 256. Using the "birthday paradox" (research this term if it is new to you), estimate roughly how many random inputs you would need to hash before you had a 50% chance of finding two that collide, and explain why this number is so much smaller than you might intuitively expect. Then explain, using this estimate, why 256 bits was chosen instead of something smaller, in terms of the resulting collision-search cost.

9. Write a short design document (300-500 words) proposing what would go wrong for a blockchain if its chosen hash function had a practical collision — a way to find two different blocks with the same hash — that took only a few hours of computing time to execute. Walk through a concrete attack scenario an adversary could carry out against a chain using such a hash function, referencing the tamper-evidence diagram from Section 9.
