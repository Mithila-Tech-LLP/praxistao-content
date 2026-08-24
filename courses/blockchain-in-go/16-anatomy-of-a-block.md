# Chapter 16: Anatomy of a Block

A block is the unit of agreement in a blockchain: a bundle of everything that happened in one "round," stamped with proof of when it happened and what came before it. Before writing a single line of Go, this chapter dissects exactly which fields GoChain's block needs, why each one exists, and how the design compares to a real Bitcoin block.

## Table of Contents

1. [What Problem a Block Actually Solves](#1-what-problem-a-block-actually-solves)
2. [GoChain's Block Layout, Field by Field](#2-gochains-block-layout-field-by-field)
3. [Header vs. Body: Two Different Jobs](#3-header-vs-body-two-different-jobs)
4. [Side-by-Side: GoChain vs. a Real Bitcoin Block](#4-side-by-side-gochain-vs-a-real-bitcoin-block)
5. [Why Field Order and Structure Matter for Hashing](#5-why-field-order-and-structure-matter-for-hashing)
6. [Drawing the Chain: Blocks Linked by Hash Arrows](#6-drawing-the-chain-blocks-linked-by-hash-arrows)
7. [Common Points of Confusion](#7-common-points-of-confusion)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Problem a Block Actually Solves

Imagine three friends keeping a shared notebook of who owes whom money. If every single transfer got its own page, the notebook would be slow to check ("did we agree on page 47 yet?") and easy to mess with (rip out page 47, nobody notices for a while). Instead, the friends agree on a better rule: once a day, gather every transfer that happened since the last entry, write them all onto one page, and at the top of that page write a short summary of the page before it — so anyone flipping to page 12 can instantly tell it belongs after page 11, not floating on its own.

That daily page is a **block**: a batch of records (in GoChain's case, transactions) sealed together with a reference back to the block before it. A blockchain is nothing more than a sequence of these pages, each one anchored to the last.

```
Transactions arriving over time (in GoChain, they wait in a mempool,
built properly in Chapter 34; for now just think "a waiting room"):

  tx: Alice pays Bob 5
  tx: Carol pays Dave 2       ---->  bundled into  ---->  BLOCK N
  tx: Bob pays Erin 1                one round            (one Timestamp,
  tx: Frank pays Alice 3                                   one Hash,
                                                            one PrevBlockHash)
```

Notice each individual transaction does *not* get its own timestamp, its own hash-link to "the previous transaction," or its own place in a chain. Only the block does. This is a deliberate design choice: a transaction's trustworthiness comes from *which block it ends up inside*, not from any property of its own. Pull a transaction out of its block, and it loses its place in history entirely — which is exactly why tampering with even one transaction, as Chapter 19 demonstrates hands-on, breaks the hash of the entire block that contains it, and every block after that one.

Batching many transactions into one block instead of hashing and linking every transaction individually solves three problems at once:

- **Efficiency.** Hashing and verifying one block of a thousand transactions is far cheaper than doing it a thousand separate times.
- **Ordering.** Every node needs to agree not just on *which* transactions happened, but in *what order* — grouping them into height-numbered blocks gives everyone the same sequence to check.
- **A single point of validation.** Once a block is accepted, every transaction inside it is accepted together, as an atomic unit. There's no in-between state where three of a block's five transactions are "half-approved."

A **block**, then, is a data structure holding a list of transactions plus enough extra information (a timestamp, a link to the previous block, and — starting in Volume 4 — proof of work) to make that batch verifiable and tamper-evident.

This batching also has a very concrete, practical side: a block has to fit somewhere. Real Bitcoin blocks are capped around one to a few megabytes, which in practice holds somewhere between one and a few thousand ordinary transactions. GoChain does not impose a hard size limit in this volume (that becomes relevant once mining and network propagation costs matter, later in the course), but keep the picture in your head: a block is not an unbounded bucket, it's a single, boundable page in the ledger.

**New term — block:** one batch of transactions, plus metadata proving when it was created and what block came before it, treated by the network as a single, atomic unit of history.

---

## 2. GoChain's Block Layout, Field by Field

Here is GoChain's block, previewed as a Go struct. We won't add any methods until Chapter 17 — this chapter is purely about understanding what each field means and why it's there.

```go
package core

// Block represents one batch of transactions plus the proof that links it
// to the block before it. Implemented fully in Chapter 17.
type Block struct {
	Height        int64
	Timestamp     int64
	Transactions  []*Transaction
	PrevBlockHash []byte
	Hash          []byte
	Nonce         uint64
	MerkleRoot    []byte
}
```

Going through it field by field:

- **`Height`** — the block's position in the chain, counting from zero. The genesis block (Chapter 18) is height 0, the next block is height 1, and so on. Height gives every block an unambiguous place in history, independent of its hash.
- **`Timestamp`** — the Unix time (seconds since January 1, 1970) the block was created. This is a rough, self-reported clock — not a guarantee of exact real-world time — but it's still useful for ordering, difficulty adjustment (Chapter 26), and human-readable inspection (Chapter 21).
- **`Transactions`** — the actual payload of this round: a slice of pointers to `Transaction` values. GoChain's `Transaction` type is fully designed in Volume 5; for this volume, treat it as an opaque signed record — something that happened, with an ID we can hash, and nothing more.
- **`PrevBlockHash`** — the fingerprint (hash) of the block that came immediately before this one. This single field is what turns a pile of independent blocks into a *chain* — more on this in Chapter 18.
- **`Hash`** — this block's own fingerprint, computed from everything else in the block. Any node can recompute this from scratch and compare it to what's stored, which is exactly how tampering gets caught (Chapter 19).
- **`Nonce`** — a number that, starting in Volume 4, gets searched over during mining until the block's hash satisfies a difficulty target. For this entire volume, `Nonce` simply sits at its zero value — it exists in the struct now so we never have to change the block's shape later.
- **`MerkleRoot`** — a single hash that summarizes the entire `Transactions` list, built with the `MerkleTree` from Chapter 10. It lets us detect if even one transaction inside the block was altered, without needing to re-hash every transaction individually to notice.

```
+-------------------- Block --------------------+
| Height:        3                               |
| Timestamp:     1719792000                      |
| Transactions:  [tx0, tx1, tx2]                  |
| PrevBlockHash: 7f3a...c92e  (points backward)  |
| Hash:          d81b...44f0  (this block's ID)  |
| Nonce:         0            (Volume 4 fills in)|
| MerkleRoot:    9c2e...0a71  (summary of txs)   |
+-------------------------------------------------+
```

**New term — fingerprint (hash):** a short, fixed-size value produced by feeding data into a hash function (Chapter 08), such that any change to the input, however tiny, produces a completely different fingerprint. We use "hash" and "fingerprint" interchangeably in this course.

Two more design details are worth pausing on before moving to the next section, because both are questions almost every newcomer asks the first time they see this struct.

**Why is `Timestamp` an `int64` instead of Go's built-in `time.Time`?** `time.Time` is a rich type — it carries a monotonic clock reading, a location (time zone), and internal fields that can vary between Go versions and platforms in ways that are not guaranteed to serialize identically everywhere. An `int64` holding a plain Unix timestamp (seconds since 1970, always UTC, always just a number) has none of that baggage: it serializes to the same 8 bytes on every machine, in every Go version, forever. Section 5 of this chapter is all about why that kind of cross-machine determinism matters — `Timestamp` is a small, concrete example of the same principle showing up in a field type choice, not just in serialization code.

**Why is `PrevBlockHash` a `[]byte` (a hash) instead of a `*Block` (an actual pointer to the previous block)?** This one cuts to the heart of what makes a blockchain a *chain of custody* rather than just a linked list. A `*Block` pointer only means something inside one running program's memory — it can't be written to disk, sent across a network to another node, or checked by someone who doesn't already have the exact same block sitting in memory. A hash, by contrast, is a portable, self-contained proof: any node, anywhere, holding a copy of "the previous block" can recompute its hash and compare it to what's stored, with no shared memory and no trust required. Go's ordinary `*Block` pointers are perfect for chaining data structures *within* one process (a use you'll see constantly in this course), but the links *between* blocks in the blockchain itself have to survive being written to a file, sent over a socket, and checked by a stranger's computer — which is a job for a hash, not a pointer.

One design detail worth pausing on: `Transactions` is typed as `[]*Transaction` — a slice of *pointers* to transactions, not a slice of transaction values (`[]Transaction`). Chapter 04 introduced the general rule that Go prefers pointers when a value is large, needs to be shared, or needs to be modified through a reference rather than copied. Here, a `Transaction` (once fully designed in Volume 5) will carry a signature, a public key, and potentially several inputs and outputs — not tiny. Storing `*Transaction` avoids copying all of that data every time a block is passed around, constructed, or handed to a function, and it means every part of GoChain that holds a reference to "transaction #3 in block #7" is looking at the exact same underlying data, not an independent copy that could quietly drift out of sync.

It helps to see each field's Go type and its "empty" state side by side, since a few of them behave a little differently from what you might expect coming from other languages:

| Field           | Go type          | Size (typical)        | Zero value                         |
|-----------------|-------------------|-------------------------|-------------------------------------|
| `Height`        | `int64`           | 8 bytes                 | `0` (which is also the genesis block's real height — not a sentinel) |
| `Timestamp`     | `int64`           | 8 bytes                 | `0` (meaning midnight, Jan 1 1970 — never valid for a real block, so it's a good "did I forget to set this?" tell) |
| `Transactions`  | `[]*Transaction`  | varies                  | `nil` (a block with no transactions is still valid — see the empty-block note in Chapter 17) |
| `PrevBlockHash` | `[]byte`          | 32 bytes (SHA-256 output) | `nil` (only ever acceptable for a block that hasn't been built yet) |
| `Hash`          | `[]byte`          | 32 bytes                | `nil` until `ComputeHash` runs |
| `Nonce`         | `uint64`          | 8 bytes                 | `0` (which is also its starting value before mining begins in Volume 4) |
| `MerkleRoot`    | `[]byte`          | 32 bytes                | `nil` until computed from `Transactions` |

Two things worth noticing here. First, `Height`'s zero value (`0`) is not a "missing data" marker the way it often is in other contexts — it's the genesis block's *real, correct* height. Second, both `Hash` and `PrevBlockHash` are 32 bytes because GoChain's `crypto.Hash` function (Chapter 09) is built on SHA-256, which always produces exactly 32 bytes of output no matter how much data goes in — a fixed-size fingerprint for an unbounded amount of input. When you print a hash for a human to read (as Chapter 21's inspector does), those 32 raw bytes get shown in hexadecimal — 64 characters like `7f3ac9e1...`, two hex digits per byte — purely for readability; the actual value stored in memory and compared during validation is always the raw bytes.

To make all seven fields feel like real data rather than abstract names, here's what a single, fully-populated block looks like as a Go value, using placeholder byte slices where a real hash would normally go (we're only constructing the struct literal by hand here — no methods exist yet, that's Chapter 17):

```go
block := core.Block{
	Height:        2,
	Timestamp:     1719792000, // 2024-06-30 20:00:00 UTC
	Transactions:  []*core.Transaction{tx0, tx1, tx2},
	PrevBlockHash: []byte{0x7a, 0x3f, 0xc9, 0xe1 /* ...28 more bytes */},
	Hash:          []byte{0xc8, 0x1e, 0x44, 0xb0 /* ...28 more bytes */},
	Nonce:         0,           // untouched until Volume 4
	MerkleRoot:    []byte{0x9c, 0x2e, 0x0a, 0x71 /* ...28 more bytes */},
}
```

Printed for a human (again, jumping ahead to formatting Chapter 21 will do properly), this is the block from the diagram in Section 6 — height 2, pointing back at the block whose hash was `7a3fc9e1...`, itself producing the hash `c81e44b0...` that block 3 will need to point back at next.

---

## 3. Header vs. Body: Two Different Jobs

Real blockchains split a block conceptually — and sometimes literally, in code — into two parts with two different jobs:

- The **header** is a small, fixed-size summary: enough information to hash the block, link it to its predecessor, and verify proof of work, *without* needing the full list of transactions at all.
- The **body** is the full list of transactions — potentially large, potentially thousands of entries.

Why bother separating them? Because a **light client** (a program, like a mobile wallet, that doesn't want to store the entire multi-gigabyte blockchain) can download just the tiny headers for the whole chain, verify the proof-of-work chain is legitimate, and only fetch full transaction bodies for the specific blocks it actually cares about. Headers being small and self-contained is what makes this possible.

GoChain's `Block` struct doesn't split into two literal Go types — that's more machinery than we need this early. But the fields naturally fall into the same two groups:

```
Header-like fields (small, fixed-size):        Body (can be large):
  Height                                          Transactions
  Timestamp
  PrevBlockHash
  Nonce
  MerkleRoot   <-- this is the bridge: a small
                   fixed-size field that still
                   summarizes the (possibly huge)
                   Transactions list
```

Notice `MerkleRoot` is the hinge between the two halves: it's small and lives in the "header" group, but it's a fingerprint *of* the body. This is exactly why Chapter 10 built a Merkle tree in the first place — it lets the small header group make a verifiable claim about the (much larger) body without embedding the whole body.

Here's the payoff in concrete terms. Imagine a mobile wallet app that wants to prove to you, its user, that a specific payment really was included in block 800,000 of a real blockchain — without downloading that block's (possibly thousands of) other transactions.

```
Light client's problem:               Light client's solution:
"Was my payment included in           1. Download only the header of block
 block 800,000?"                         800,000 (tiny — ~80 bytes for
                                          Bitcoin).
Downloading the FULL block             2. Ask a full node for a Merkle proof
(megabytes, thousands of                  — a short list of sibling hashes
other people's transactions)              (Chapter 10) connecting your
would work, but it's slow and             transaction up to MerkleRoot.
wasteful for a phone that just        3. Recompute the path yourself and
wants to check ONE payment.               check it matches the header's
                                           MerkleRoot you already have.
```

This is precisely the "light client" use case Chapter 10 introduced Merkle proofs to solve, and it's the concrete reason real blockchains bother separating header from body at all: without that separation, "verify one payment" and "download the entire block" would be the same amount of work.

The storage difference this creates in practice is dramatic. A **full node** (one that validates and stores everything) needs the entire body of every block, forever. A **light client** only ever needs headers, plus the occasional small Merkle proof:

| | Full node | Light client |
|---|---|---|
| Stores headers | Yes | Yes |
| Stores full transaction bodies | Yes, every block, forever | No — only fetches specific ones on demand |
| Can independently verify proof of work | Yes | Yes (headers alone are enough) |
| Can independently verify a specific payment | Yes | Yes, via a Merkle proof from a full node |
| Storage for a chain with millions of blocks | Many gigabytes and growing | A few hundred megabytes of headers, roughly constant per block |

GoChain doesn't build a light-client mode of its own in this volume — that's a natural extension once Volume 7's networking exists — but the `MerkleRoot` field is what makes it *possible* later without redesigning the `Block` struct.

---

## 4. Side-by-Side: GoChain vs. a Real Bitcoin Block

Bitcoin's actual block header is a fixed 80 bytes, made of these fields (shown here as an illustrative struct — this is not literal Bitcoin source code, just a plain-language mapping of its header format):

```go
// Illustrative only — a simplified view of Bitcoin's real block header.
type BitcoinBlockHeader struct {
	Version       int32  // format version of the block
	PrevBlockHash [32]byte
	MerkleRoot    [32]byte
	Timestamp     uint32 // seconds since epoch
	Bits          uint32 // encoded difficulty target
	Nonce         uint32
}
```

| Concept                     | GoChain field           | Bitcoin header field     | Notes                                                        |
|-----------------------------|--------------------------|---------------------------|----------------------------------------------------------------|
| Position in chain            | `Height`                | *(not in header)*         | Bitcoin derives height from chain position; it's also encoded inside the coinbase transaction by convention (BIP34), not the header itself. |
| Creation time                 | `Timestamp`              | `Timestamp`               | Same idea; Bitcoin's is 32-bit, GoChain's is a 64-bit Unix timestamp. |
| Transactions                  | `Transactions`           | *(in the body, not header)* | Bitcoin's header holds only the `MerkleRoot`, not the transactions themselves. |
| Link to previous block        | `PrevBlockHash`          | `PrevBlockHash`            | Identical idea and role in both designs. |
| This block's own fingerprint   | `Hash`                  | *(computed, not stored)*  | Bitcoin doesn't store its own hash inside the header — it's always recomputed on demand. GoChain stores it for convenience; Chapter 17 explains why this doesn't cause the trap you might expect. |
| Proof-of-work search variable  | `Nonce`                 | `Nonce`                    | Same purpose; Bitcoin's is 32 bits, which is part of why real miners also vary the timestamp and coinbase transaction to find enough search space. |
| Difficulty target              | *(Volume 4)*             | `Bits`                     | GoChain adds an equivalent in `consensus.ProofOfWork`, starting in Chapter 24. |
| Summary of transaction list    | `MerkleRoot`             | `MerkleRoot`               | Identical role: a small fingerprint standing in for the whole transaction list. |
| Format/protocol version        | *(not yet needed)*       | `Version`                  | GoChain doesn't need this until later volumes add protocol evolution. |

The similarity is the point: GoChain's block isn't a simplified toy invented for this course, it's a close cousin of the exact structure securing billions of dollars of real value today. The differences that do exist (explicit `Height`, a stored `Hash`, 64-bit fields instead of 32-bit) are deliberate simplifications that make the code easier to read and test in Go, not corners cut on the underlying idea.

To make this concrete, here's (a simplified description of) Bitcoin's actual genesis block, mined by Satoshi Nakamoto on January 3rd, 2009:

```
Bitcoin Genesis Block (block 0)
  Timestamp:      1231006505  (2009-01-03 18:15:05 UTC)
  PrevBlockHash:  0000000000000000000000000000000000000000000000000000000000000000
                  (all zero — no predecessor, exactly like GoChain's genesis rule)
  MerkleRoot:     4a5e1e4baab89f3a32518a88c31bc87f618f76673e2cc77ab2127b7afdeda33b
                  (the hash of the single coinbase transaction it contains)
  Nonce:          2083236893  (the winning value found by mining)
  Hash:           000000000019d6689c085ae165831e934ff763ae46a2a6c172b3f1b60a8ce26f
                  (the famous "genesis hash," referenced by every Bitcoin block since)
```

Notice the same two conventions GoChain adopts: `PrevBlockHash` is all zeroes because there is truly nothing before the first block, and the block's `Hash` is not something Satoshi typed in — it's the *result* of hashing everything else, found only after searching many nonce values until the proof-of-work condition was met (a story Volume 4 tells in full).

As a final piece of real-world grounding: Bitcoin's block format hasn't stood still. The 2017 "SegWit" (Segregated Witness) upgrade changed how signature data is counted toward a block's size limit, precisely because signatures were taking up so much room that the effective *transaction* capacity of a block was shrinking. This is a good preview of a lesson GoChain returns to throughout the course: a block format that looks "final" on paper often has to evolve as a real network hits real limits — which is exactly why GoChain keeps its own `Block` struct intentionally minimal for now, adding fields only once a later volume gives a concrete reason to (proof of work in Volume 4, and so on), rather than guessing at every field a "real" blockchain might eventually want.

---

## 5. Why Field Order and Structure Matter for Hashing

A hash function (Chapter 08) doesn't understand "fields" — it only understands a sequence of bytes. Before we can hash a `Block`, its fields have to be turned into bytes in some specific, fixed way. This process is called **serialization**, and it was introduced properly in Chapter 09's discussion of *canonical* serialization: the same logical block must always produce the exact same byte sequence, on every machine, every time.

Here's why this matters so much for a blockchain specifically. Every node independently:

1. Receives (or builds) a block.
2. Serializes it into bytes.
3. Hashes those bytes.
4. Compares the result against the block's stored `Hash`.

If two nodes serialized the *same logical block* into two *different* byte sequences — say, because one iterated over a Go map in a different order, or padded a number differently — they would compute two different hashes for a block that is, in every meaningful sense, identical. One node would accept the block; the other would reject it as corrupted. The network would disagree about its own history. This is why GoChain avoids maps entirely in anything that gets hashed, and why Chapter 17 builds `ComputeHash` by hand from a fixed, deliberate field order rather than trusting a generic "serialize this struct" helper to do the right thing by default.

To make "different byte sequence, same logical data" concrete: suppose a block has `Height = 1` and `Nonce = 256`. One valid-looking approach might serialize `Height` before `Nonce`; another, equally reasonable-looking approach might do it the other way around:

```
Order A: [Height bytes][Nonce bytes]
         00 00 00 00 00 00 00 01   00 00 00 00 00 00 01 00
         \______ Height=1 ______/  \______ Nonce=256 _____/

Order B: [Nonce bytes][Height bytes]
         00 00 00 00 00 00 01 00   00 00 00 00 00 00 00 01
         \______ Nonce=256 ______/  \______ Height=1 _____/
```

Both byte strings encode exactly the same information — the same block, by any reasonable definition. But they are different sequences of bytes, so `crypto.Hash` (which only ever looks at raw bytes, never at "meaning") produces two completely different, unrelated-looking fingerprints for them. If even one node in the network picked "Order B" while everyone else used "Order A," that node would permanently disagree with everyone else about this block's hash — a silent, catastrophic bug that has nothing to do with any real tampering. GoChain avoids this by fixing the exact order once, in Chapter 17's `ComputeHash`, and never changing it.

There's a second, subtler trap worth naming here, even though we won't fix it until Chapter 17 shows the code: a block's `Hash` field cannot be part of what gets hashed to *produce* `Hash`. If it were, you'd have a circular definition — the hash would depend on itself. Keep this in the back of your mind; Chapter 17 walks through exactly why this trips up nearly everyone's first blockchain implementation, and exactly how GoChain avoids it.

```
Same logical block, serialized two different ways:

Node A serializes:                    Node B serializes:
[Height][Time][PrevHash][Root][Nonce]  [Nonce][Height][PrevHash][Time][Root]
        |                                      |
        v                                      v
   hash = 7a3f...                        hash = 91c2...

Two different hashes for "the same" block --> the network can no
longer agree on anything. Canonical, fixed field order prevents this.
```

---

## 6. Drawing the Chain: Blocks Linked by Hash Arrows

Put four blocks next to each other and the whole idea clicks into place. Each block's `PrevBlockHash` is literally a copy of the hash the previous block already computed:

```
 Block 0 (genesis)        Block 1                   Block 2                   Block 3
+-----------------+      +-----------------+       +-----------------+       +-----------------+
| Height: 0       |      | Height: 1       |       | Height: 2       |       | Height: 3       |
| PrevHash: 000..0|      | PrevHash: 7a3f..|<------| PrevHash: c81e..|<------| PrevHash: 44b0..|
| Hash:     7a3f..|----->| Hash:     c81e..|------>| Hash:     44b0..|------>| Hash:     9de2..|
+-----------------+      +-----------------+       +-----------------+       +-----------------+
```

Read the arrows as "points back to." Block 1's `PrevBlockHash` (`7a3f..`) is exactly Block 0's `Hash`. Block 2's `PrevBlockHash` (`c81e..`) is exactly Block 1's `Hash`. This is the entire linking mechanism — no separate index, no external table of "what comes after what," just each block quietly holding onto its predecessor's fingerprint. Chapter 18 turns this diagram into working Go code.

It's worth noticing what this design makes *hard* on purpose. Suppose someone wanted to slip a brand new block in between Block 1 and Block 2 — call it Block 1.5. They'd need Block 1.5's `PrevBlockHash` to equal Block 1's real `Hash` (fine, that's copyable), but they'd also need Block 2's `PrevBlockHash` to equal Block 1.5's `Hash` instead of Block 1's — which means editing Block 2, which changes Block 2's own `Hash`, which means editing Block 3's `PrevBlockHash` too, and so on for every block after the insertion point:

```
Inserting a block in the middle forces you to rewrite every link after it:

 Block 0        Block 1        [Block 1.5]        Block 2        Block 3
+--------+     +--------+      +----------+       +--------+     +--------+
| h: 7a3f|---->| h: c81e|--+   | h: NEW!  |    +-->| h: MUST|--+  | h: MUST|
+--------+     +--------+  |   +----------+    |   | CHANGE |  |  | CHANGE |
                           +-->| prev:c81e|----+   +--------+  +->+--------+
                               +----------+
```

Nothing stops you from doing this rewriting on your own private copy of the chain — the mechanism itself isn't magic. What stops you from getting *other people* to accept your rewritten history is the subject of Volume 4: making each block's hash expensive to (re-)produce, so redoing this cascade for a long chain becomes computationally infeasible rather than merely inconvenient. This chapter and the next few only build the tamper-*evident* half of the story — spotting a rewrite instantly. Making rewrites tamper-*proof* — actually expensive to pull off — comes next.

---

## 7. Common Points of Confusion

A few terminology mix-ups trip up almost everyone the first time they look closely at a block. Clearing them up now saves confusion for the rest of the course.

**"Block," "blockchain," and "chain" are not the same thing.** A *block* is one struct — one batch of transactions plus its metadata. The *blockchain* (`core.Blockchain`, built in Chapter 18) is the whole ordered sequence of blocks, from genesis to the current tip. "Chain" is casual shorthand for the same thing as "blockchain." You'll sometimes hear "the chain" used loosely to mean the whole system (as in "GoChain's chain is now three blocks long"), but it never refers to a single block on its own.

**`Height` is the same idea as "block number" or "block index."** Different projects and articles use different words for identical concepts — you will see "block height," "block number," and occasionally "block index" used completely interchangeably in blockchain writing (including, at times, in this very course). GoChain's field is called `Height` because that's the term Bitcoin's own documentation and community use most consistently.

**A block's `Hash` is not the same kind of thing as a transaction's `ID`, or an address.** All three are hashes — fixed-size fingerprints produced by the same underlying `crypto.Hash` function from Chapter 09 — but they're fingerprints *of different data*. A block's `Hash` fingerprints the block's own header-like fields. A transaction's `ID` (Volume 5) fingerprints that transaction's own contents. An address (Chapter 14) is derived by hashing a *public key*, not a block or a transaction at all. Seeing three hex strings that all "look the same" (64 hex characters) is normal; they answer three completely different questions.

**It's normal for `Nonce` to just sit at `0` throughout this entire volume.** Nothing in Chapters 16 through 22 ever changes it, and that's by design — this volume is deliberately building the *linking and validation* half of a blockchain before the *proof-of-work* half exists at all. Don't go looking for mining logic yet; Volume 4 is where `Nonce` starts actually doing something.

**A block does not "contain" the previous block.** `PrevBlockHash` is only ever a *fingerprint* of the previous block — 32 bytes — never a copy of the previous block itself. This is precisely what keeps blocks small and bounded: no block's size grows with the size of the chain behind it, no matter how many blocks came before it.

Putting the header/body split and the chain-linking mechanism into one combined picture, here's the complete anatomy this chapter has built up, piece by piece:

```
                         BLOCKCHAIN (Chapter 18)
   +----------------------------------------------------------------+
   |                                                                  |
   |   BLOCK 0            BLOCK 1            BLOCK 2                 |
   |  +---------+        +---------+        +---------+              |
   |  | HEADER  |        | HEADER  |        | HEADER  |              |
   |  |---------|        |---------|        |---------|              |
   |  | Height  |        | Height  |        | Height  |              |
   |  | Time    |        | Time    |        | Time    |              |
   |  | PrevHash|<-\     | PrevHash|<-\     | PrevHash|<-\           |
   |  | Nonce   |   \    | Nonce   |   \    | Nonce   |   \          |
   |  | MerkRoot|-\  \   | MerkRoot|-\  \   | MerkRoot|-\  \         |
   |  | Hash    |-)--)---+ Hash    |-)--)---+ Hash    |-)--)-- (next) |
   |  +---------+ | |      +---------+ | |      +---------+ | |      |
   |  | BODY    | | |      | BODY    | | |      | BODY    | | |      |
   |  |---------| | |      |---------| | |      |---------| | |      |
   |  | tx0     | | |      | tx0     | | |      | tx0     | | |      |
   |  | tx1     |-+ |      | tx1     |-+ |      | tx1     |-+ |      |
   |  +---------+   |      +---------+   |      +---------+   |      |
   |       ^        |           ^        |           ^        |      |
   |       |        |           |        |           |        |      |
   |       +-- MerkleRoot summarizes the BODY, and lives in the HEADER
   |                |                        |                        |
   |                +-- Hash summarizes the whole HEADER (never itself)|
   |                                                                  |
   +----------------------------------------------------------------+
```

Every idea from this chapter is in that one picture: the header/body split (Section 3), `MerkleRoot` bridging the two (also Section 3), `PrevBlockHash` linking backward (Section 6), and `Hash` summarizing everything except itself (the trap flagged in Section 5, and resolved in code next chapter).

---

## Summary

- A block batches many transactions into one round, giving the network efficiency, a shared ordering, and a single unit of validation.
- GoChain's `Block` has seven fields: `Height`, `Timestamp`, `Transactions`, `PrevBlockHash`, `Hash`, `Nonce`, and `MerkleRoot`.
- `Height` and `Timestamp` are simple bookkeeping; `Transactions` is the payload; `PrevBlockHash` and `Hash` form the chain's links; `Nonce` waits unused until Volume 4; `MerkleRoot` is a small fingerprint standing in for the (possibly large) transaction list.
- Real blockchains split blocks into a small, fixed-size **header** and a larger **body** — GoChain's fields fall into the same two conceptual groups even though they live in one Go struct.
- Bitcoin's real block header is strikingly similar to GoChain's design: `PrevBlockHash`, `MerkleRoot`, `Timestamp`, and `Nonce` all have direct equivalents.
- Deterministic, canonical field ordering during serialization is not a nitpick — without it, different nodes could compute different hashes for the same logical block and permanently disagree about the chain's history.
- A block's own `Hash` field must never be part of the data that gets hashed to compute it — a subtlety Chapter 17 resolves in code.
- The chain itself is nothing more than each block's `PrevBlockHash` matching the actual `Hash` of the block right before it.

---

## Exercises

### Easy

1. Draw (on paper or in a text file) five blocks by hand, at heights 0 through 4, and fill in fake but internally consistent `Hash` and `PrevBlockHash` values (e.g., short strings like `"h0"`, `"h1"`) so the chain links correctly. Then change block 2's `Hash` on paper and mark every block whose link is now broken.
2. List each of GoChain's seven `Block` fields and, in one sentence each, explain what would go wrong for the chain as a whole if that field were removed entirely.
3. Using the comparison table in Section 4, write down which Bitcoin header field has no GoChain equivalent yet, and which volume of this course you'd expect to add it in (a quick skim of the table of contents will tell you).

### Medium

4. Explain, in your own words, why `MerkleRoot` belongs in the "header-like" group of fields even though it's a summary of something that lives in the "body" group (`Transactions`). Specifically address:
   - What would a light client have to do to verify one payment if `MerkleRoot` did not exist at all?
   - How does the answer change if the block instead had 5,000 transactions rather than 3?
   - Why does the *size* of `MerkleRoot` (always 32 bytes, per Section 2's field-size table) staying constant regardless of transaction count matter here?
5. Suppose GoChain used Go's built-in `encoding/json` package to serialize a block for hashing, and two different machines happened to produce byte-for-byte identical JSON except for the order of fields (JSON produced from a map is not guaranteed to preserve field order across languages/versions). Using the Order A / Order B byte example from Section 5 as a model:
   - Write out (in plain text, not real JSON) two differently-ordered but logically identical serializations of a two-field block.
   - Explain concretely what symptom this would cause across a real network of nodes.
   - Propose one concrete rule GoChain's serialization code could follow to guarantee this never happens.
6. Bitcoin's header does not store `Height` directly. Research (or reason from what you already know about hashing and chains) how a node can still determine a given block's height without it being an explicit field, and explain one advantage and one disadvantage of GoChain's choice to store `Height` explicitly instead.

### Hard

7. Redesign the combined diagram from Section 6 to show what happens when a fifth block, Block 4, is added:
   - Extend the arrows and hash values so Block 4 correctly links to Block 3.
   - Explicitly show that Block 4 does *not* need to know anything about Block 0, 1, or 2 to be considered valid — only about Block 3.
   - Now suppose Block 1's `Timestamp` is discovered to be wrong. Mark, on your diagram, every field in every later block that would need to be recomputed to "fix" it without detection — using the insertion example from Section 6 as your guide.
8. Bitcoin's `Nonce` field is only 32 bits, which real miners exhaust quickly at modern hash rates — forcing them to also vary the timestamp and a piece of the coinbase transaction to keep searching. GoChain's `Nonce` is a 64-bit `uint64`. Estimate (rough order of magnitude is fine) how many more distinct values a 64-bit nonce provides compared to a 32-bit one, and explain whether you think GoChain will ever run into the same "nonce exhaustion" problem Bitcoin miners deal with.
9. Propose one additional field a production blockchain might need that GoChain's current design does not have (examples to get you thinking: a protocol version number, a difficulty/target value, a block size limit, a chain identifier to prevent replaying transactions from a different network). For your proposed field:
   - Justify why it would matter with a concrete scenario where its absence causes a real problem.
   - Sketch which existing field group (header-like or body-like, per Section 3) it belongs in.
   - Say whether you'd expect it to be a fixed size (like `Hash`) or variable size (like `Transactions`), and why that matters for the header/body split.
