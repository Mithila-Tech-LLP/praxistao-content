# Chapter 09: Implementing a Hasher in Go

Chapter 08 built a conceptual understanding of what SHA-256 does and why it behaves the way it does. This chapter turns that understanding into real, running Go code: the first function in `gochain/crypto`, a single call that will sit underneath every block, every transaction, and every wallet address for the rest of this course. Along the way, you will run into a problem that has nothing to do with SHA-256 itself and everything to do with a much more ordinary source of bugs — turning a Go struct into bytes in a way that is actually trustworthy to hash. Getting this one detail right, once, here, is what makes every later chapter's hashing "just work."

## Table of Contents

1. [Setting Up the crypto Package](#1-setting-up-the-crypto-package)
2. [crypto/sha256 in Go — Your First Hash Function](#2-cryptosha256-in-go--your-first-hash-function)
3. [Hex-Encoding for Human-Readable Fingerprints](#3-hex-encoding-for-human-readable-fingerprints)
4. [Hashing Strings, Numbers, and Anything Else](#4-hashing-strings-numbers-and-anything-else)
5. [The Canonical Serialization Problem](#5-the-canonical-serialization-problem)
6. [Bug One — Randomized Map Order Breaks Determinism](#6-bug-one--randomized-map-order-breaks-determinism)
7. [Bug Two — Ambiguous Field Boundaries Create Accidental Collisions](#7-bug-two--ambiguous-field-boundaries-create-accidental-collisions)
8. [The Serialize Method Pattern](#8-the-serialize-method-pattern)
9. [Testing the crypto Package](#9-testing-the-crypto-package)
10. [Where This Fits Going Forward](#10-where-this-fits-going-forward)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Setting Up the crypto Package

Chapter 03 laid out `gochain`'s package map, and `crypto` is the one this course fills in first, because almost everything else — blocks in Volume 3, signatures in this same volume, addresses in Chapter 14 — depends on it. Inside your `gochain` module (module path `github.com/you/gochain`, created back in Chapter 03), create a `crypto` directory with a single file to start:

```
gochain/
├── go.mod
├── main.go
└── crypto/
    └── hash.go
```

Every file in this directory declares `package crypto`, and every other GoChain package that needs to hash something will import it as `github.com/you/gochain/crypto`. Nothing in this package will ever import from `core`, `wallet`, or any other GoChain package — hashing is a foundation, and foundations do not depend on the things built on top of them.

## 2. crypto/sha256 in Go — Your First Hash Function

Go's standard library ships a correct, audited, fast implementation of SHA-256 in `crypto/sha256`, exactly as Chapter 08 promised: nobody in this course hand-writes the padding, chunking, or 64-round compression from Chapter 08, Section 6. We simply call it.

```go
// crypto/hash.go
package crypto

import "crypto/sha256"

// Hash computes the SHA-256 fingerprint of data and returns the raw
// 32-byte digest. It is deterministic (Chapter 08, Section 2): the
// same bytes in always produce the same 32 bytes out, on any machine,
// forever.
func Hash(data []byte) []byte {
	sum := sha256.Sum256(data)
	return sum[:]
}
```

Two small Go details worth pausing on if you have not used `crypto/sha256` before. First, `sha256.Sum256` returns a `[32]byte` — a fixed-size *array*, not a slice — because Go can guarantee the output is always exactly 32 bytes at compile time. Second, `sum[:]` converts that fixed array into a `[]byte` slice, which is the shape every other GoChain function will expect, since slices (unlike arrays) can be passed around, appended to, and stored in structs without their length being baked into the type itself.

That is the entire function. It looks almost too small to matter — and that smallness is the point. `Hash` is a thin, honest wrapper around a primitive the standard library already gets right; the interesting engineering in this chapter is not this function, it is making sure everything you *feed* to this function is trustworthy, which Sections 5 through 8 cover in depth.

## 3. Hex-Encoding for Human-Readable Fingerprints

`Hash` returns raw bytes — 32 of them, each a number from 0 to 255. Printed directly, they look like garbage on a terminal (non-printable bytes render as boxes, question marks, or control characters). Chapter 08, Section 1 already established the fix: encode the bytes as hexadecimal, the same 64-character format every hash shown in that chapter used. Go's standard library provides this too, in `encoding/hex`:

```go
package main

import (
	"encoding/hex"
	"fmt"

	"github.com/you/gochain/crypto"
)

func main() {
	fingerprint := crypto.Hash([]byte("gochain"))
	fmt.Println(hex.EncodeToString(fingerprint))
	// 94270d9fbb414dbcc537dc89ca5b5b4d7ca2d63732bad19d4d04b924ca6425e8
}
```

`hex.EncodeToString` takes a `[]byte` and returns a `string` where every byte becomes exactly two hex characters (32 bytes → 64 characters, matching Chapter 08's table exactly). There is a matching `hex.DecodeString` for the reverse direction, which you will use starting in Chapter 21 to read a hash back in from a human-typed or stored hex string. Because this conversion is so common, it is worth wrapping in a tiny helper right away:

```go
// crypto/hash.go

// HashHex computes the SHA-256 fingerprint of data and returns it as
// a lowercase, 64-character hex string -- the readable form used
// anywhere a hash needs to be printed, logged, or stored as text.
func HashHex(data []byte) string {
	return hex.EncodeToString(Hash(data))
}
```

From here on, this book uses `Hash` when a function needs to *work* with a fingerprint (compare it, store it, embed it in another struct) and `HashHex` (or a bare `hex.EncodeToString(Hash(...))`) whenever a human needs to *read* one — printed to a terminal, shown in a block explorer in Volume 10, or written into this book's own diagrams.

## 4. Hashing Strings, Numbers, and Anything Else

`Hash` takes a `[]byte`, but most of the data you actually care about hashing does not start out as bytes. A string converts trivially — `[]byte("some string")` — because Go strings are already just an immutable sequence of bytes under the hood, and this conversion copies them into a mutable slice `Hash` can accept:

```go
crypto.Hash([]byte("10 gochips"))
```

A number needs a decision about *which* bytes represent it, because "the number 10" is not a self-evidently unique byte sequence — is it the string `"10"` (two bytes, `0x31 0x30`), or the 8-byte big-endian encoding of the integer `10` (`encoding/binary`'s job), or something else? Both are valid choices, but they produce *different* hashes for what a person would call "the same value" — which is a preview of this chapter's real subject. Whichever encoding you pick, you must pick exactly one and use it everywhere that number gets hashed, or two representations of "10" will silently stop matching each other. Section 5 generalizes this exact problem to entire structs, where it matters far more.

## 5. The Canonical Serialization Problem

Here is the problem this chapter exists to solve. `Hash` only ever sees bytes — it has no idea a `Block` or a `Note` or a `Transaction` exists as a concept. Before you can hash a Go struct, something has to first flatten it into a `[]byte`. That flattening step is called **serialization**, and the function that does it, by convention throughout this course, is named `Serialize()`.

The trouble is that "flatten a struct into bytes" sounds like a solved, mechanical problem, but it hides a sharp requirement: serialization must be **canonical**. A canonical representation means that two values a person would call "the same data" always produce the *exact same bytes* — not "usually the same," not "the same on this run of the program," but always, unconditionally the same. If it is not, everything Chapter 08 promised about hashing quietly breaks, without SHA-256 itself doing anything wrong at all.

Think of two notaries asked to produce a certified copy of the same contract. If one notary happens to staple the pages in a different order, or writes the date in `07/31/2026` while the other writes `31 July 2026`, the two "copies" are legally the same contract but are two different *documents* — and if you fingerprinted each document rather than the contract's meaning, you would get two different fingerprints for something that should count as identical. A hash function has no concept of "meaning" at all; it only ever sees the exact bytes you hand it. If your `Serialize()` method can produce different bytes for data a person considers equal — or, just as dangerously, the *same* bytes for data that is actually different — the hash built on top of it inherits that flaw, silently.

```
        Note A                    Note B
   (logically identical    (logically identical
    to Note B)              to Note A)
        │                          │
        ▼                          ▼
   Serialize()                Serialize()
        │                          │
        ▼                          ▼
   bytes: 7a3f...            bytes: 91c0...   ◀── different!
        │                          │
        ▼                          ▼
     Hash()                     Hash()
        │                          │
        ▼                          ▼
   fingerprint A             fingerprint B   ◀── different, wrongly
```

This is not a hypothetical. The next two sections walk through two real, specific ways an ordinary-looking `Serialize()` method breaks this guarantee in Go — one that makes identical data hash *differently*, and one that makes different data hash *the same* — using a small example type built specifically to demonstrate the pattern, before Section 8 fixes both.

## 6. Bug One — Randomized Map Order Breaks Determinism

Define a small type to experiment with — not a type GoChain uses later, just a stand-in built for this chapter to demonstrate the `Serialize()` pattern on:

```go
// crypto/note.go
package crypto

// Note is a small, self-contained example type used only to
// demonstrate the Serialize() pattern. GoChain's real types --
// starting with core.Block in Chapter 17 -- follow this exact pattern.
type Note struct {
	Title string
	Body  string
	Tags  map[string]bool
}
```

A natural first instinct is to reach for a generic encoder rather than write serialization by hand — Chapter 07 introduced `encoding/gob` as GoChain's early-chapters choice for exactly this reason. So try it:

```go
var buf bytes.Buffer
gob.NewEncoder(&buf).Encode(note)
fingerprint := crypto.Hash(buf.Bytes())
```

This works — right up until a `Note` contains a map. Go deliberately randomizes the iteration order of a `map` every time one is created, specifically so programmers cannot accidentally rely on a stable order that was never guaranteed. `encoding/gob` does not sort map keys before writing them out; it writes them in whatever order the map happens to iterate in. The result: two `Note` values with the exact same title, body, and tags — built by inserting those same tags in a different order — can serialize to two different byte sequences, because their two `Tags` maps iterate differently, even within the very same run of your program:

```go
n1 := Note{Title: "genesis", Tags: map[string]bool{
	"go": true, "blockchain": true, "hash": true,
}}
n2 := Note{Title: "genesis", Tags: map[string]bool{
	"hash": true, "go": true, "blockchain": true, // same tags, different insert order
}}

// gob-encoding n1 and n2 can produce DIFFERENT bytes here, even
// though every human reading this code would call n1 and n2 the
// same Note.
```

This is exactly the failure diagrammed in Section 5: logically identical data, different serialized bytes, different hashes. It is not a SHA-256 bug — `Hash` is doing exactly what determinism (Chapter 08, Section 2) promises: identical bytes in, identical fingerprint out. The bug is upstream, in handing a generic encoder something whose byte layout it does not fully control. This matters well beyond this toy example: Volume 5 represents the entire UTXO set as a Go map, and Volume 8 builds indexes on top of maps too — the exact same trap, at much higher stakes, if it goes unnoticed.

## 7. Bug Two — Ambiguous Field Boundaries Create Accidental Collisions

The opposite failure is just as easy to write by accident, and just as dangerous. Suppose you "fix" serialization by writing your own encoder instead of trusting a generic one, and simply concatenate every field's bytes directly:

```go
// Do not do this -- see below.
func (n Note) BrokenSerialize() []byte {
	return []byte(n.Title + n.Body)
}
```

This looks perfectly reasonable, and it is deterministic (the same `Note` always concatenates to the same string). The problem is that concatenation without marking where one field ends and the next begins is *ambiguous*: it throws away information about the boundary between fields. Consider these two clearly different `Note` values:

```go
a := Note{Title: "ab", Body: "c"}
b := Note{Title: "a", Body: "bc"}
```

`a` and `b` are obviously not the same note — different titles, different bodies. But `a.Title + a.Body` is `"ab" + "c"` = `"abc"`, and `b.Title + b.Body` is `"a" + "bc"` = `"abc"` too. Both serialize to the identical byte string `"abc"`, so both hash to the identical fingerprint:

```
   Note{Title: "ab", Body: "c"}  ──▶  "ab"+"c" ──▶ "abc" ──▶ Hash ──▶ fingerprint X
   Note{Title: "a",  Body: "bc"} ──▶  "a"+"bc" ──▶ "abc" ──▶ Hash ──▶ fingerprint X

   Two different Notes.  Same bytes.  Same hash.  A self-inflicted
   collision -- nothing to do with SHA-256's real collision
   resistance from Chapter 08, Section 5, and entirely our own bug.
```

This is worth sitting with, because it is easy to conflate with the birthday-paradox-style collisions Chapter 08 discussed — but it is a completely different, much more mundane kind of failure. SHA-256's collision resistance is about how hard it is to find two *different byte strings* that hash the same; it offers zero protection here, because `"abc"` and `"abc"` are not two different byte strings at all — our own serialization code handed `Hash` the identical input twice, for data a person would insist is different. No amount of cryptographic strength downstream can fix a mistake made upstream of it.

## 8. The Serialize Method Pattern

Both bugs share one root cause: handing byte-layout decisions to something that was not designed to make them correctly for our purposes — a generic encoder that does not know maps need sorting, or naive concatenation that does not know fields need marking. The fix, and the pattern every GoChain type follows from here on, is to give each type an explicit `Serialize()` method that fully controls its own byte layout: fixed field order, every collection sorted into a deterministic sequence first, and every variable-length field marked so its boundary can never be confused with the next field's.

```go
// crypto/note.go
package crypto

import (
	"bytes"
	"encoding/binary"
	"sort"
)

// writeField appends a length-prefixed copy of s to buf. Prefixing
// the length means two fields written back to back can never be
// misread as one longer field, or split at the wrong boundary --
// exactly the bug from Section 7.
func writeField(buf *bytes.Buffer, s string) {
	var lenBytes [4]byte
	binary.BigEndian.PutUint32(lenBytes[:], uint32(len(s)))
	buf.Write(lenBytes[:])
	buf.WriteString(s)
}

// Serialize converts a Note into a single, canonical byte slice: the
// same logical Note always produces the exact same bytes, regardless
// of what order its fields were set in or what order its Tags map
// happens to iterate in on a given run -- fixing both Section 6 and
// Section 7's bugs.
func (n Note) Serialize() []byte {
	var buf bytes.Buffer

	writeField(&buf, n.Title)
	writeField(&buf, n.Body)

	// Map iteration order is randomized (Section 6), so we cannot
	// loop over n.Tags directly. Sorting the keys first makes the
	// output canonical regardless of insertion order.
	keys := make([]string, 0, len(n.Tags))
	for k := range n.Tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	binary.Write(&buf, binary.BigEndian, uint32(len(keys)))
	for _, k := range keys {
		writeField(&buf, k)
		if n.Tags[k] {
			buf.WriteByte(1)
		} else {
			buf.WriteByte(0)
		}
	}

	return buf.Bytes()
}

// Hash returns the SHA-256 fingerprint of the Note's canonical byte
// representation.
func (n Note) Hash() []byte {
	return Hash(n.Serialize())
}
```

Three small decisions here matter far more than the four lines of code they take. First, **fixed field order**: `Serialize()` always writes `Title`, then `Body`, then `Tags`, regardless of the order the struct literal listed them in — struct field order in Go source does not affect memory layout guarantees the way this method's *own* ordering choice does. Second, **length-prefixing every string** with `writeField`, so `"ab"` followed by `"c"` can never be confused with `"a"` followed by `"bc"` — Section 7's bug, closed. Third, **sorting map keys** before writing them, so two `Tags` maps with identical contents always serialize identically no matter how they were built — Section 6's bug, closed.

This is the pattern — a hand-written `Serialize()` method, one per type, that treats byte layout as a decision worth making explicitly rather than delegating — that Chapter 17 applies to `core.Block`, Chapter 32 applies to `core.Transaction`, and every serializable type after them applies in turn. `Note` will not appear again after this chapter; it exists purely so this pattern could be taught on something small before you meet it again on something that carries real value.

## 9. Testing the crypto Package

Chapter 07 set up `go test` and table-driven tests; this is the first package that puts them to real use. Three kinds of test matter here: that `Hash` is deterministic, that a *known* input produces a specific, previously verified digest (a "known-answer test," catching any accidental change to how you call `crypto/sha256`), and that `Serialize()` actually is canonical — directly testing the property Sections 6 through 8 were built around.

```go
// crypto/hash_test.go
package crypto

import (
	"bytes"
	"encoding/hex"
	"testing"
)

func TestHash_Deterministic(t *testing.T) {
	data := []byte("gochain")
	h1 := Hash(data)
	h2 := Hash(data)
	if !bytes.Equal(h1, h2) {
		t.Fatalf("expected identical hashes, got %x and %x", h1, h2)
	}
}

func TestHash_KnownVector(t *testing.T) {
	got := hex.EncodeToString(Hash([]byte("gochip")))
	want := "b84c623ee4489528cd1dfd55c552cc3739df6bf3ade94251282f3c107373aa83"
	if got != want {
		t.Fatalf("got %s, want %s", got, want)
	}
}

func TestNoteSerialize_SameContentSameBytes(t *testing.T) {
	n1 := Note{
		Title: "genesis",
		Body:  "first note",
		Tags:  map[string]bool{"go": true, "blockchain": true, "hash": true},
	}
	n2 := Note{
		Title: "genesis",
		Body:  "first note",
		Tags:  map[string]bool{"hash": true, "go": true, "blockchain": true},
	}

	if !bytes.Equal(n1.Serialize(), n2.Serialize()) {
		t.Fatal("expected identical Serialize() output for logically identical Notes")
	}
	if !bytes.Equal(n1.Hash(), n2.Hash()) {
		t.Fatal("expected identical Hash() output for logically identical Notes")
	}
}

func TestNoteSerialize_NoFieldAmbiguity(t *testing.T) {
	a := Note{Title: "ab", Body: "c"}
	b := Note{Title: "a", Body: "bc"}

	if bytes.Equal(a.Serialize(), b.Serialize()) {
		t.Fatal("Title/Body boundary is ambiguous: two different Notes serialized identically")
	}
}
```

`TestNoteSerialize_SameContentSameBytes` is the test that would have caught Section 6's bug immediately, had you written it against the broken `gob`-based version: build the same logical `Note` twice with `Tags` populated in a different order, and assert the two `Serialize()` outputs are byte-for-byte equal. `TestNoteSerialize_NoFieldAmbiguity` is the test that catches Section 7's bug: build two obviously different `Notes` whose naive concatenation would collide, and assert their serialized bytes differ. Run `go test ./crypto/...` and all four should pass — and if you temporarily swap `Serialize()` back to the broken `gob` or naive-concatenation versions, watch the corresponding test fail, which is a worthwhile five minutes spent confirming these tests actually test what they claim to.

## 10. Where This Fits Going Forward

Two functions and one method leave this chapter and go on to underpin the entire rest of GoChain: `crypto.Hash([]byte) []byte`, `crypto.HashHex([]byte) string`, and the `Serialize() []byte` pattern demonstrated on `Note`. Chapter 10 immediately puts `Hash` to work at a larger scale, combining many transaction hashes into a single Merkle root instead of hashing one flat blob. Chapter 17 applies the exact `Serialize()` pattern from Section 8 to `core.Block` — fixed field order, sorted collections, length-prefixed variable data — so that a block's hash can be trusted the same way every hash in this book has been trusted since Chapter 08: the same block, hashed anywhere, by anyone, always produces the same fingerprint, and no two different blocks ever accidentally look the same to `Hash`.

---

## Summary

- `crypto/sha256` in Go's standard library provides a correct, audited SHA-256 implementation; `gochain/crypto.Hash(data []byte) []byte` wraps `sha256.Sum256` into the shape the rest of GoChain will call.
- `encoding/hex.EncodeToString` turns a raw hash into the readable, 64-character form used everywhere a fingerprint needs to be printed, logged, or compared by a human — wrapped here as `crypto.HashHex`.
- Hashing a struct requires first turning it into bytes (**serialization**), and that serialization must be **canonical**: logically identical data must always produce identical bytes, or the hash built on top inherits the mismatch without SHA-256 doing anything wrong.
- A generic encoder like `encoding/gob` can silently break canonical serialization when a struct contains a `map`, because Go randomizes map iteration order and `gob` does not sort keys before encoding — identical data can produce different bytes on the same run.
- Naive field concatenation without marking field boundaries can create the opposite bug: two genuinely different values (e.g., `Title:"ab",Body:"c"` vs. `Title:"a",Body:"bc"`) serializing to the identical byte string, and therefore the identical hash — a self-inflicted collision unrelated to SHA-256's real collision resistance.
- The fix and the pattern GoChain uses from here on: every serializable type gets its own hand-written `Serialize()` method with a fixed field order, sorted collections, and length-prefixed variable-length fields.
- Tests for a `Serialize()` method should directly verify canonicality: build the same logical value two different ways and assert identical output; build two different values that a naive encoder would collide and assert distinct output.

---

## Exercises

### Easy

1. Add a `Serialize()`-based `String() string` method to `Note` that returns `hex.EncodeToString(n.Hash())`. Write a short paragraph explaining why this method should call `Hash()` rather than just returning `n.Title + n.Body` directly.

2. Using the `Note` type from this chapter, construct two `Note` values with identical `Title` and `Body` but `Tags` built by inserting the same three tag names in two different orders. Print `hex.EncodeToString(n.Hash())` for both and confirm, by running the code yourself, that they match.

3. The chapter's `writeField` function prefixes each string with a 4-byte big-endian length. Explain in your own words what would go wrong if you used a 1-byte length prefix instead, for a `Body` field that might contain more than 255 characters.

### Medium

4. Temporarily replace `Note.Serialize()` with a version that uses `encoding/gob` directly on the whole struct (as shown broken in Section 6). Run `TestNoteSerialize_SameContentSameBytes` repeatedly (`go test -run TestNoteSerialize_SameContentSameBytes -count=20 ./crypto/...`) and report whether it passes every time, sometimes, or never. Explain your result in terms of how map iteration order is actually randomized in Go (research `runtime` map internals if you want more depth than the chapter gives).

5. Add an `Attachments []string` field to `Note` (an ordered list, not a map) and extend `Serialize()` to include it. Since a Go slice, unlike a map, already has a well-defined order, explain in 3-5 sentences why you still need a length prefix for the *number* of attachments, even though you do not need to sort them the way `Tags` needed sorting.

6. Write a new test, `TestNoteSerialize_FieldOrderMatters`, that builds two `Note` values with `Title` and `Body` *swapped* (e.g., `Note{Title: "x", Body: "y"}` vs. `Note{Title: "y", Body: "x"}`) and asserts their `Serialize()` output differs. Explain why this test would still pass even without the length-prefixing fix from Section 8 — i.e., why field order and field-boundary ambiguity are two separate bugs, not the same one.

### Hard

7. Research how `encoding/json` handles map keys in Go (specifically: since which Go version, and what rule it applies) and explain, in a short comparison (200-300 words), why `encoding/json` would not have suffered Section 6's bug the way `encoding/gob` did, if `Note.Serialize()` had used `json.Marshal` instead. Then explain why this course still recommends hand-written `Serialize()` methods over `encoding/json` for hashing purposes, even knowing this.

8. Design (in Go pseudocode or real code) a `Serialize()` method for a hypothetical type with a field of type `float64`. Research IEEE 754's handling of `NaN` and signed zero (`-0.0` vs `0.0`), and explain a scenario in which two `float64` values that a person might consider "the same number" could still produce different raw bytes when serialized naively — and what a canonical serializer would need to do about it.

9. The chapter's fix for Section 6 sorts `Tags`' keys alphabetically before serializing. Suppose a future GoChain type needs to serialize a collection of *structs* (not strings) that have no natural total ordering a person would agree is "correct" (say, a set of `TxInput` values). Write a short design note (300-400 words) proposing a rule for canonically ordering such a collection before serialization, and explain why the rule itself does not need to be meaningful to a human — only deterministic and total (any two distinct elements must have a definite order) — to solve the canonicalization problem.
