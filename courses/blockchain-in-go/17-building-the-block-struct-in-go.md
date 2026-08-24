# Chapter 17: Building the Block Struct in Go

Chapter 16 designed every field of GoChain's block on paper and explained why each one exists. This chapter turns that design into real, compiling, tested Go code: the `core.Block` struct, a constructor that builds one from a list of transactions, and the two methods — `Serialize()` and `ComputeHash()` — that let any node, anywhere, independently compute the exact same fingerprint for the exact same block. Along the way, this chapter walks straight into the single trap that catches almost every first-time blockchain implementer: what happens when a struct's hash field accidentally becomes part of the data being hashed.

## Table of Contents

1. [From Diagram to Code — Setting Up core](#1-from-diagram-to-code--setting-up-core)
2. [A Stand-In for Transaction](#2-a-stand-in-for-transaction)
3. [The Block Struct, Field for Field](#3-the-block-struct-field-for-field)
4. [Serialize — A Canonical Byte Layout for the Block](#4-serialize--a-canonical-byte-layout-for-the-block)
5. [The Self-Reference Trap — Why Hash Must Exclude Itself](#5-the-self-reference-trap--why-hash-must-exclude-itself)
6. [ComputeHash — Fingerprinting a Block](#6-computehash--fingerprinting-a-block)
7. [NewBlock — Constructing a Block From Its Transactions](#7-newblock--constructing-a-block-from-its-transactions)
8. [Putting It Together — Building and Printing a Real Block](#8-putting-it-together--building-and-printing-a-real-block)
9. [Testing core.Block](#9-testing-coreblock)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. From Diagram to Code — Setting Up core

Chapter 03 sketched out `gochain`'s package map, and `core` is the package that has waited the longest for real code: it is where `Block` and, starting next chapter, `Blockchain` live. Add it to your module now:

```
gochain/
├── go.mod
├── main.go
├── crypto/
│   ├── hash.go
│   └── merkle.go
└── core/
    ├── transaction.go
    └── block.go
```

`core` will import `crypto` — every hash `core.Block` ever computes is `crypto.Hash` under the hood, exactly as Chapter 09 built it, and every Merkle root it stores comes from `crypto.NewMerkleTree`, exactly as Chapter 10 built it. `crypto` will never import `core` back; the dependency only flows one direction, from foundation to building, which is precisely the layering Chapter 06 asked you to keep in mind.

Think of this chapter as pouring concrete into the formwork Chapter 16 built out of plywood and string. The shape was decided already — seven fields, in a specific order, with specific meanings. Nothing about *what* the block contains changes here. What changes is that it becomes a real Go type you can construct, hash, print, and — starting in Chapter 19 — actually validate.

## 2. A Stand-In for Transaction

`Block.Transactions` is a `[]*Transaction`, and `Transaction` is a type Volume 5 designs properly — signing, UTXOs, transaction IDs derived from real cryptographic rules. Building `Block` today still needs *something* named `Transaction` to compile against, so this chapter adds exactly the minimal, honest stand-in the contract for this course promises: the fields Volume 5 will use, plus a `Serialize()` method good enough for this volume's hashing and tamper-evidence work, and nothing more.

```go
// core/transaction.go
package core

import (
	"bytes"
	"encoding/binary"
)

// Transaction is fully designed starting in Chapter 32 (Volume 5) --
// signing, unlocking scripts, real transaction IDs. For this volume,
// treat it as an opaque signed record: something that happened, with
// an ID and a way to turn it into bytes, and nothing more. Every
// field below is exactly what Volume 5 will build on top of; nothing
// here needs to change later, only grow.
type Transaction struct {
	ID        []byte
	Inputs    []TxInput
	Outputs   []TxOutput
	Timestamp int64
}

// TxInput references one output of an earlier transaction that this
// transaction is spending, plus the proof (Signature, PublicKey) that
// whoever built this transaction was allowed to spend it. Chapter 33
// explains exactly what gets signed and why.
type TxInput struct {
	TxID      []byte
	OutIndex  int
	Signature []byte
	PublicKey []byte
}

// TxOutput is a discrete chunk of value assigned to whoever can prove
// they control PubKeyHash. Chapter 30 explains this "UTXO" model in
// full; for now it is just a value and a destination.
type TxOutput struct {
	Value      int64
	PubKeyHash []byte
}

// writeBytes appends a length-prefixed copy of data to buf -- the
// exact same helper Chapter 09 introduced as writeField, so two
// variable-length fields written back to back can never be misread
// as one field, or split at the wrong boundary.
func writeBytes(buf *bytes.Buffer, data []byte) {
	var lenBytes [4]byte
	binary.BigEndian.PutUint32(lenBytes[:], uint32(len(data)))
	buf.Write(lenBytes[:])
	buf.Write(data)
}

// Serialize returns a canonical byte representation of the
// transaction -- fixed field order, every variable-length field
// length-prefixed, exactly the pattern Chapter 09 built on Note.
// This is deliberately minimal; Chapter 32 replaces it with the real
// serialization used for signing and transaction IDs. It exists here
// purely so core.Block has something real to hash and Merkle-root
// against in this volume.
func (tx *Transaction) Serialize() []byte {
	var buf bytes.Buffer

	writeBytes(&buf, tx.ID)
	binary.Write(&buf, binary.BigEndian, tx.Timestamp)

	binary.Write(&buf, binary.BigEndian, uint32(len(tx.Inputs)))
	for _, in := range tx.Inputs {
		writeBytes(&buf, in.TxID)
		binary.Write(&buf, binary.BigEndian, int64(in.OutIndex))
		writeBytes(&buf, in.Signature)
		writeBytes(&buf, in.PublicKey)
	}

	binary.Write(&buf, binary.BigEndian, uint32(len(tx.Outputs)))
	for _, out := range tx.Outputs {
		binary.Write(&buf, binary.BigEndian, out.Value)
		writeBytes(&buf, out.PubKeyHash)
	}

	return buf.Bytes()
}
```

Nothing about this file will need to be un-learned later. Volume 5 adds signing, verification, and a proper way to compute `ID`, but the shape — fields plus a canonical `Serialize()` — is exactly the pattern this entire course uses for every hashable type, and it is the same pattern Section 4 now applies to `Block` itself.

## 3. The Block Struct, Field for Field

Here, at last in code, is the exact struct Chapter 16 designed:

```go
// core/block.go
package core

import (
	"bytes"
	"encoding/binary"
	"time"

	"github.com/you/gochain/crypto"
)

// Block represents one batch of transactions plus the proof that
// links it to the block before it. See Chapter 16 for the full
// rationale behind each field.
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

This is not a new design decision — it is the same seven fields, same types, same order Chapter 16, Section 2 laid out on paper. What Section 4 below adds is the part Chapter 16 explicitly deferred: turning this struct into bytes in a way every node can trust.

## 4. Serialize — A Canonical Byte Layout for the Block

Chapter 16, Section 5 already explained *why* this matters: a hash function does not understand fields, only bytes, and if two nodes turn the "same" logical block into two different byte sequences, they compute two different, permanently disagreeing hashes for data everyone would call identical. `Serialize()` is where GoChain fixes that byte layout, once, following the exact pattern Chapter 09 built on `Note`: a fixed field order, and every variable-length field wrapped with `writeBytes` so field boundaries can never blur into each other.

```go
// core/block.go

// Serialize returns a canonical byte representation of the block, in
// a fixed field order, suitable for hashing. Notice what is
// deliberately absent: the Hash field itself. Section 5 explains
// exactly why that omission is not an oversight.
func (b *Block) Serialize() []byte {
	var buf bytes.Buffer

	binary.Write(&buf, binary.BigEndian, b.Height)
	binary.Write(&buf, binary.BigEndian, b.Timestamp)
	writeBytes(&buf, b.PrevBlockHash)
	writeBytes(&buf, b.MerkleRoot)
	binary.Write(&buf, binary.BigEndian, b.Nonce)
	binary.Write(&buf, binary.BigEndian, uint32(len(b.Transactions)))

	return buf.Bytes()
}
```

Walk through why each line is there. `Height`, `Timestamp`, and `Nonce` are all fixed-size integers, so `binary.Write` (from Go's `encoding/binary` package) writes them as a predictable, constant number of bytes every time — no length prefix needed, because their size never varies. `PrevBlockHash` and `MerkleRoot` are `[]byte` — technically fixed at 32 bytes each in practice, but `writeBytes` prefixes them with their length anyway, for the same defense-in-depth reason Chapter 09 prefixed every string: it costs four bytes and permanently closes off the field-boundary-ambiguity bug from Chapter 09, Section 7, even if a future refactor ever changes a hash's size.

Notice, too, what `Serialize()` does *not* do: it never loops over `b.Transactions` and serializes each transaction's full bytes into this buffer. It only writes `len(b.Transactions)` — a count. This is deliberate, and it is the header/body split from Chapter 16, Section 3 showing up directly in code: `MerkleRoot` already is a fingerprint of the entire transaction list (Section 7 shows exactly how it gets computed), so re-serializing every transaction's bytes into the block's own hash input would be redundant work that adds nothing — if a single transaction changes, `MerkleRoot` already changes completely, which is all `ComputeHash` needs to notice. Real Bitcoin block headers work the same way: the 80-byte header that gets hashed contains a Merkle root, never the raw transaction list itself.

```
Serialize() writes, in this exact order, forever:

  [Height: 8 bytes][Timestamp: 8 bytes]
  [len+PrevBlockHash][len+MerkleRoot]
  [Nonce: 8 bytes][TxCount: 4 bytes]

  Every node that serializes "the same" block produces
  byte-for-byte identical output -- canonical, by construction.
```

## 5. The Self-Reference Trap — Why Hash Must Exclude Itself

Here is the trap Chapter 16, Section 5 flagged and promised to resolve: `Block.Hash` is a field that holds the block's own fingerprint. If `Serialize()` included `Hash` in its output, you would be asking a hash function to answer an impossible question — "what is the hash of a message that includes its own hash?" — which has no well-defined answer at all.

An everyday analogy makes the circularity obvious. Imagine trying to write your own signature on a form, in a box labeled "attach a photocopy of your completed signature here," *before* you have signed the form. There is no way to satisfy that instruction — the photocopy needs a finished signature to copy, but the signature is not finished until the photocopy is already attached. The instruction is not just hard, it is logically impossible to satisfy in order.

```
The impossible version (never do this):

   Hash = SHA-256( Height, Timestamp, ..., Hash, Nonce )
                                        ^^^^
                            Hash depends on Hash --
                            no well-defined answer exists.

The correct version (what Serialize() actually does):

   Hash = SHA-256( Height, Timestamp, PrevBlockHash,
                    MerkleRoot, Nonce, TxCount )
                                        ^^^^^^^^^^^^
                     everything EXCEPT Hash -- well-defined,
                     computable, and the same on every machine.
```

Look back at Section 4's `Serialize()` — it never once mentions `b.Hash`. That is not a missing line; it is the entire fix. `Hash` is the *output* of hashing everything else, so it can never also be part of the *input*. This is precisely why the field-size table in Chapter 16, Section 2 listed `Hash`'s zero value as `nil` "until `ComputeHash` runs" — before that call, there is nothing to be circular about, because the field simply has not been filled in yet.

A common beginner mistake is writing a generic "serialize the whole struct" helper — reflection-based, or a blind `encoding/gob` call on the entire `Block` value — and only noticing the circularity once two supposedly identical blocks stubbornly hash differently, or once a block's hash keeps changing every time you print it. GoChain avoids the trap entirely by writing `Serialize()` by hand, field by field, and simply never writing a line that touches `b.Hash`. There is nothing to remember to exclude at hash time, because it was never included in the first place.

## 6. ComputeHash — Fingerprinting a Block

With `Serialize()` doing the hard work of producing canonical, self-reference-free bytes, `ComputeHash` becomes almost embarrassingly small — which is exactly the sign of a well-factored design, the same lesson Chapter 09, Section 2 made about `crypto.Hash` itself.

```go
// core/block.go

// ComputeHash returns the SHA-256 fingerprint of the block's
// canonical serialized bytes -- everything the block is, except its
// own Hash field. Any node, anywhere, holding an identical block can
// call this and get the identical 32-byte answer back, per Chapter
// 08's determinism guarantee.
func (b *Block) ComputeHash() []byte {
	return crypto.Hash(b.Serialize())
}
```

One line. `crypto.Hash` (Chapter 09) takes the `[]byte` that `Serialize()` (Section 4) hands it and runs SHA-256 over it, exactly as it has for every other GoChain type since `Note`. There is no special-casing for `Block` inside `crypto.Hash` at all — it has no idea a block exists — which is exactly the point of building `crypto` as a foundation package with no knowledge of what sits on top of it.

`ComputeHash` is named deliberately differently from the `Hash` *field*. Calling it recomputes a fresh fingerprint from the block's current contents, every single time you call it — it never reads or trusts `b.Hash` at all. This distinction is what makes Chapter 19's `ValidateBlock` possible next chapter: comparing a freshly recomputed `ComputeHash()` against the *stored* `b.Hash` is exactly how tampering gets caught, and that comparison only means something because the two are computed independently.

## 7. NewBlock — Constructing a Block From Its Transactions

A constructor ties `Serialize()` and `ComputeHash()` together into the one function most of GoChain will actually call: hand it a list of transactions, the previous block's hash, and a height, and get back a fully-formed, already-hashed `Block`.

```go
// core/block.go

// NewBlock builds a new block from transactions, linking it to
// prevBlockHash and stamping it at height. It computes the block's
// MerkleRoot over transactions and its Hash over everything else,
// leaving Nonce at its zero value -- untouched until Volume 4's
// mining loop starts searching over it.
func NewBlock(transactions []*Transaction, prevBlockHash []byte, height int64) *Block {
	block := &Block{
		Height:        height,
		Timestamp:     time.Now().Unix(),
		Transactions:  transactions,
		PrevBlockHash: prevBlockHash,
		Nonce:         0,
	}

	block.MerkleRoot = MerkleRootOf(transactions)
	block.Hash = block.ComputeHash()

	return block
}

// MerkleRootOf builds a Merkle tree (Chapter 10) over transactions'
// serialized bytes and returns its root -- a single 32-byte
// fingerprint standing in for the entire list, exactly as Chapter 16,
// Section 3 described. An empty transaction list still gets a valid,
// well-defined root: the hash of an empty input.
//
// This is exported, rather than kept package-private, for two
// reasons that only become concrete in later chapters: Chapter 19's
// ValidateBlock needs to recompute it independently, from outside a
// single NewBlock call, to check that a block's body still matches
// what its header claims; and Chapter 22 reuses it directly, outside
// any cryptocurrency context at all, to build a general-purpose
// tamper-evident log.
func MerkleRootOf(transactions []*Transaction) []byte {
	if len(transactions) == 0 {
		return crypto.Hash([]byte{})
	}

	data := make([][]byte, len(transactions))
	for i, tx := range transactions {
		data[i] = tx.Serialize()
	}

	tree := crypto.NewMerkleTree(data)
	return tree.MerkleRoot()
}
```

Notice the order of operations inside `NewBlock`, because it matters: `MerkleRoot` is computed *before* `Hash`, since `Serialize()` (Section 4) reads `b.MerkleRoot` as one of its fields — if `Hash` were computed first, using a not-yet-set `MerkleRoot`, the fingerprint would be wrong the instant `MerkleRoot` was filled in afterward. `Hash` genuinely must be the last field set, precisely because it is the fingerprint of everything that came before it.

An empty block — zero transactions — is intentionally still valid here, matching Chapter 16, Section 2's note that `Transactions` being `nil` is a legitimate, not broken, state. `crypto.Hash([]byte{})` is a perfectly well-defined SHA-256 computation (an empty input is still an input), so an empty block gets a real, stable `MerkleRoot` rather than a crash or a `nil` fingerprint.

## 8. Putting It Together — Building and Printing a Real Block

Here is `NewBlock` in action, building a block over three placeholder transactions:

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/you/gochain/core"
)

func main() {
	txs := []*core.Transaction{
		{ID: []byte("tx1"), Timestamp: 1719792000},
		{ID: []byte("tx2"), Timestamp: 1719792001},
		{ID: []byte("tx3"), Timestamp: 1719792002},
	}

	genesisHash := make([]byte, 32) // Chapter 18 explains this all-zero convention

	block := core.NewBlock(txs, genesisHash, 1)

	fmt.Println("Height:       ", block.Height)
	fmt.Println("Timestamp:    ", block.Timestamp)
	fmt.Println("Tx count:     ", len(block.Transactions))
	fmt.Println("PrevBlockHash:", hex.EncodeToString(block.PrevBlockHash))
	fmt.Println("MerkleRoot:   ", hex.EncodeToString(block.MerkleRoot))
	fmt.Println("Hash:         ", hex.EncodeToString(block.Hash))
	fmt.Println("Nonce:        ", block.Nonce)
}
```

Running this prints something like:

```
Height:        1
Timestamp:     1719858213
Tx count:      3
PrevBlockHash: 0000000000000000000000000000000000000000000000000000000000000000
MerkleRoot:    9c2e0a71d84f3b6e1a5c9d02f7e8b1a3c4d5e6f708192a3b4c5d6e7f8091a2b3c4
Hash:          c81e44b0f0a291d7e3b5c8a9f2013d4e5f6a7b8c9d0e1f203141516171819202
Nonce:         0
```

Run it a second time (without changing the transactions) and every value stays identical except `Timestamp` — that field alone depends on wall-clock time, so `Hash` will differ between runs purely because `Timestamp` did. Call `block.ComputeHash()` yourself, right after building `block`, and compare it to `block.Hash` — they must be byte-for-byte equal, since `NewBlock` sets one directly from the other. That equality is not a coincidence to double-check once and forget; Chapter 19 makes it the very first rule `ValidateBlock` checks on every block, forever.

## 9. Testing core.Block

Following Chapter 09's testing pattern, three properties are worth locking down with real tests: determinism, the self-reference exclusion from Section 5, and the avalanche effect showing up correctly at the block level.

```go
// core/block_test.go
package core

import (
	"bytes"
	"testing"
)

func testTx(id string) *Transaction {
	return &Transaction{ID: []byte(id), Timestamp: 1719792000}
}

func TestNewBlock_HashMatchesComputeHash(t *testing.T) {
	block := NewBlock([]*Transaction{testTx("tx1")}, make([]byte, 32), 1)

	if !bytes.Equal(block.Hash, block.ComputeHash()) {
		t.Fatal("block.Hash must always equal a freshly computed ComputeHash()")
	}
}

func TestSerialize_ExcludesHashField(t *testing.T) {
	block := NewBlock([]*Transaction{testTx("tx1")}, make([]byte, 32), 1)

	before := block.Serialize()
	block.Hash = []byte("a completely different, fake hash value")
	after := block.Serialize()

	if !bytes.Equal(before, after) {
		t.Fatal("Serialize() output changed after mutating Hash -- Hash must never affect its own input")
	}
}

func TestNewBlock_Deterministic(t *testing.T) {
	txs := []*Transaction{testTx("tx1"), testTx("tx2")}
	prevHash := make([]byte, 32)

	// Building two blocks from identical inputs, at the same height,
	// must serialize identically in every field EXCEPT Timestamp --
	// so we compare MerkleRoot directly rather than Hash, since Hash
	// legitimately differs when Timestamp differs between calls.
	b1 := NewBlock(txs, prevHash, 1)
	b2 := NewBlock(txs, prevHash, 1)

	if !bytes.Equal(b1.MerkleRoot, b2.MerkleRoot) {
		t.Fatal("identical transactions must always produce an identical MerkleRoot")
	}
}

func TestNewBlock_AvalancheEffect(t *testing.T) {
	prevHash := make([]byte, 32)

	b1 := NewBlock([]*Transaction{testTx("tx1")}, prevHash, 1)
	b2 := NewBlock([]*Transaction{testTx("tx1-tampered")}, prevHash, 1)

	if bytes.Equal(b1.Hash, b2.Hash) {
		t.Fatal("changing a single transaction must change the block's Hash completely")
	}
}
```

`TestSerialize_ExcludesHashField` is the test that directly proves Section 5's fix works: mutate `block.Hash` after the fact, re-serialize, and assert the bytes did not change at all. If a future refactor ever accidentally added `writeBytes(&buf, b.Hash)` to `Serialize()`, this test fails immediately, loudly, and specifically — exactly the kind of test worth keeping in a codebase for years after the bug it guards against feels obvious. Run `go test ./core/...` and all four should pass.

---

## Summary

- `core.Block` is now a real, compiling Go struct with all seven fields Chapter 16 designed: `Height`, `Timestamp`, `Transactions`, `PrevBlockHash`, `Hash`, `Nonce`, `MerkleRoot`.
- A minimal `core.Transaction` stand-in (with `ID`, `Inputs`, `Outputs`, `Timestamp`, and a `Serialize()` method) lets `Block` compile and be tested now, without designing Volume 5's real transaction logic early.
- `Serialize()` turns a block into canonical bytes in a fixed field order, using length-prefixed variable fields — the exact pattern Chapter 09 introduced on `Note`.
- `Serialize()` deliberately omits the `Hash` field: including it would make the hash depend on itself, a circular, ill-defined computation with no correct answer.
- `ComputeHash()` is `crypto.Hash(b.Serialize())` — one line, made possible entirely by `Serialize()` already having done the hard work of producing trustworthy, self-reference-free bytes.
- `NewBlock()` computes `MerkleRoot` before `Hash`, because `Serialize()` reads `MerkleRoot` as one of its fields — `Hash` must always be the very last field set.
- A block's `MerkleRoot` is computed over its transactions' serialized bytes via `crypto.NewMerkleTree`, so `Serialize()` itself only needs a transaction *count*, not the full transaction list, mirroring the header/body split from Chapter 16.
- Tests should directly verify the self-reference exclusion — mutate `Hash`, re-serialize, and assert nothing changed — not just check that hashing "seems to work."

---

## Exercises

### Easy

1. Add a `String() string` method to `Block` that returns a short, human-readable one-line summary (height, transaction count, and the first 8 hex characters of its hash). Explain why it should call `hex.EncodeToString` on `b.Hash` rather than printing the raw byte slice with `%v`.

2. Build two blocks with `NewBlock`, using the exact same transactions and `prevBlockHash` but different `height` values (say, 1 and 2). Print both blocks' `Hash` values and confirm they differ, then explain in one sentence which field inside `Serialize()`'s output is responsible for that difference.

3. `MerkleRootOf` returns `crypto.Hash([]byte{})` for an empty transaction list rather than, say, `nil` or a slice of 32 zero bytes. Write a short paragraph explaining why a real, computed hash is a better choice here than either of those alternatives, referencing Chapter 16's field-size table.

### Medium

4. Temporarily "break" `Serialize()` by adding `writeBytes(&buf, b.Hash)` as its very first line, then run `TestSerialize_ExcludesHashField`. Report exactly what happens: does the test fail, panic, or hang? Explain, using Section 5's diagram, whether the resulting `ComputeHash()` value is well-defined at all, or just accidentally consistent within a single run.

5. `NewBlock` sets `Nonce` to `0` and never touches it again. Write a short test that constructs two otherwise-identical blocks, manually sets one's `Nonce` to `12345` after construction (without calling `ComputeHash` again), and checks whether their `Hash` fields still match. Explain what this demonstrates about the difference between a block's *stored* `Hash` and what `ComputeHash()` would return if called fresh.

6. Add a `TxCount() int` convenience method to `Block` that returns `len(b.Transactions)`. Then explain, in 3-5 sentences, why `Serialize()` already writes a transaction count as raw bytes (Section 4) even though `TxCount()` itself is never called during hashing — what would happen to two "same number of transactions" blocks with a different actual transaction count if that line were removed from `Serialize()`?

### Hard

7. Design and implement a `BenchmarkNewBlock` Go benchmark (using `testing.B`) that measures how long `NewBlock` takes to build a block containing 1, 100, and 10,000 placeholder transactions. Report your results and explain, referencing Chapter 10's discussion of Merkle tree height, why you would expect the *Merkle root* computation, not the final `ComputeHash` call, to dominate the time as transaction count grows.

8. Rewrite `Serialize()` to use `encoding/gob` on the whole `Block` struct directly, instead of the hand-written version in Section 4, but first set `b.Hash = nil` before encoding (to sidestep the exact circularity described in Section 5). Run the existing test suite against this version and report which tests fail, if any, and why — paying particular attention to whether `gob` encoding of a `nil` vs. an explicitly zero-length `[]byte` for `PrevBlockHash` produces identical output.

9. Chapter 16, Section 4 compared GoChain's design to Bitcoin's real header, noting Bitcoin never stores its own hash — it is always recomputed on demand. Write a design proposal (300-450 words) for an alternative `core.Block` that removes the `Hash` field entirely, relying only on `ComputeHash()` being called whenever a hash is needed. Identify at least two places elsewhere in this course's planned chapters (skim the table of contents) where *not* having a stored `Hash` field would add repeated computation, and argue for which design — stored `Hash` or always-recomputed — you would choose for GoChain, and why.
