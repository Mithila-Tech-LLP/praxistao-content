# Chapter 13: Implementing Keys and Signatures in Go

Chapter 11 built the concepts. Chapter 12 diagrammed the sign/verify pipeline in detail, step by step, without a single line of code. This chapter turns all of it into real, compiling Go: `gochain/crypto` gains a `KeyPair` type, a `NewKeyPair()` function, a `Sign()` function, and a `Verify()` function — the exact four pieces every wallet (Volume 6) and every transaction (Volume 5) will build on top of for the rest of this course. By the end of this chapter you will have watched Alice sign a real message, watched anyone verify it using only her public key, and watched that same verification fail the instant a single byte changes — the theory from Chapter 12, now running on your own machine.

## Table of Contents

1. [What This Chapter Builds](#1-what-this-chapter-builds)
2. [crypto/ecdsa and crypto/elliptic in Go](#2-cryptoecdsa-and-cryptoelliptic-in-go)
3. [Choosing and Fixing a Curve](#3-choosing-and-fixing-a-curve)
4. [The KeyPair Type](#4-the-keypair-type)
5. [Representing a Public Key as Bytes](#5-representing-a-public-key-as-bytes)
6. [Generating a Key Pair: NewKeyPair](#6-generating-a-key-pair-newkeypair)
7. [Signing Data: Sign](#7-signing-data-sign)
8. [Verifying a Signature: Verify](#8-verifying-a-signature-verify)
9. [Full Worked Example: Alice Pays Bob 5 Gochips](#9-full-worked-example-alice-pays-bob-5-gochips)
10. [Corrupting One Byte — Watching Verification Fail](#10-corrupting-one-byte--watching-verification-fail)
11. [Handling Errors Correctly](#11-handling-errors-correctly)
12. [Testing crypto.Sign and crypto.Verify](#12-testing-cryptosign-and-cryptoverify)
13. [Where This Fits in GoChain](#13-where-this-fits-in-gochain)
14. [Summary](#summary)
15. [Exercises](#exercises)

---

## 1. What This Chapter Builds

Everything in this chapter lives in the `crypto` package Chapter 09 started, alongside `Hash`, `HashHex`, and the Merkle tree code from Chapter 10. By the end of this chapter, that package will additionally export exactly this surface, matching the shape Chapter 12's diagrams promised:

```go
package crypto

type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte // uncompressed X||Y bytes
}

func NewKeyPair() (*KeyPair, error)
func Sign(priv *ecdsa.PrivateKey, data []byte) ([]byte, error)
func Verify(pubKey []byte, data, signature []byte) bool
func PublicKeyToBytes(pub *ecdsa.PublicKey) []byte
```

Notice the shape of each function mirrors a box from Chapter 12's diagrams exactly: `NewKeyPair` is Chapter 11, Section 2's "generate a private key, derive a public key" arrow. `Sign` is Chapter 12, Section 4's signing box: private key and data in, signature out. `Verify` is Chapter 12, Section 6's verifying box: public key, data, and signature in, a plain `bool` out — no private key anywhere in sight, exactly as the diagram required.

```
gochain/
├── go.mod
├── main.go
└── crypto/
    ├── hash.go        (Chapter 09)
    ├── merkle.go      (Chapter 10)
    └── keys.go        (this chapter)
```

---

## 2. crypto/ecdsa and crypto/elliptic in Go

Go's standard library splits ECDSA across two packages, and it is worth knowing what each one is responsible for before writing any code that uses them:

- **`crypto/elliptic`** defines the curves themselves — the mathematical objects Chapter 11, Section 6 described conceptually (a fixed shape, a generator point, a defined point-multiplication rule). It exposes ready-made, standard curves as simple function calls: `elliptic.P224()`, `elliptic.P256()`, `elliptic.P384()`, `elliptic.P521()`.
- **`crypto/ecdsa`** implements the actual signing and verifying algorithm from Chapter 12 *on top of* whatever curve you hand it. It provides `ecdsa.GenerateKey` (produce a new key pair on a given curve), `ecdsa.SignASN1` (produce a signature), and `ecdsa.VerifyASN1` (check one).

This split exists for the same reason Chapter 09 kept `Hash` a thin wrapper around `crypto/sha256` rather than reimplementing SHA-256's internals: the curve arithmetic and the signing algorithm built on top of it are both precise, easy-to-get-subtly-wrong mathematics that Go's standard library already implements correctly, and this course has no reason to duplicate that work by hand. GoChain's job is to wire these two packages together correctly and wrap them in the exact function signatures the rest of this course expects.

---

## 3. Choosing and Fixing a Curve

Chapter 12, Section 3 already made the decision: GoChain uses **P-256**, not Bitcoin's secp256k1, so that this entire course runs on Go's standard library alone. In code, that decision is a single line, declared once and reused everywhere a curve is needed:

```go
// crypto/keys.go
package crypto

import (
	"crypto/ecdsa"
	"crypto/elliptic"
)

// curve is the one elliptic curve every GoChain key pair, signature,
// and verification uses. Chapter 12, Section 3 explains why GoChain
// uses P-256 (built into Go's standard library) instead of Bitcoin's
// secp256k1 (which would require a third-party package).
var curve = elliptic.P256()
```

Declaring `curve` once as a package-level variable, rather than calling `elliptic.P256()` separately inside every function that needs it, matters for a reason beyond tidiness: it guarantees every key pair GoChain ever generates, and every signature GoChain ever verifies, agrees on the exact same curve, with no risk of one function accidentally using a different one. Chapter 11, Section 6 was explicit that everyone using a curve must agree on it in advance — this one line is where GoChain's code makes that agreement concrete and impossible to accidentally violate.

---

## 4. The KeyPair Type

`KeyPair` bundles the two halves of Chapter 11, Section 2's key pair into a single Go value:

```go
// KeyPair bundles a private key with its matching public key, in the
// exact shape the rest of gochain/crypto -- and, starting in Volume 6,
// gochain/wallet -- expects to work with.
type KeyPair struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte // uncompressed X||Y bytes
}
```

Two design choices here are worth pausing on. First, `PrivateKey` is stored as Go's own `*ecdsa.PrivateKey` type directly, rather than converted into some GoChain-specific representation — there is no reason to reinvent a type the standard library already gives us in exactly the shape `ecdsa.SignASN1` expects. Second, `PublicKey` is stored as a plain `[]byte`, not as an `*ecdsa.PublicKey` — this is a deliberate choice explained fully in Section 5, and it is the same shape `Verify()`'s `pubKey` parameter expects, so a `KeyPair`'s public key can be handed directly to `Verify` with zero conversion.

---

## 5. Representing a Public Key as Bytes

An `*ecdsa.PublicKey` in Go is a struct holding a curve and two large numbers, `X` and `Y` — the coordinates of the point Chapter 11, Section 6 called "the public key point," `d * G`. That struct is convenient inside a single running Go program, but it is not something you can send over a network, embed inside a transaction, or hash — those all need a plain sequence of bytes. GoChain's public key format takes the most direct possible approach: pad `X` and `Y` to a fixed width and concatenate them, with no extra framing:

```
  Public key point: (X, Y)

  X = 0x3f2a91...  (padded to 32 bytes for P-256)
  Y = 0xd817ec...  (padded to 32 bytes for P-256)

  PublicKey []byte = X bytes (32) followed by Y bytes (32) = 64 bytes total

  ┌────────────────────────────┬────────────────────────────┐
  │      X  (32 bytes)          │      Y  (32 bytes)          │
  └────────────────────────────┴────────────────────────────┘
```

The padding matters for exactly the reason Chapter 09, Section 5 cared about canonical serialization: `X` and `Y` are large numbers, and a number like `7` and a number like `1000007` do not naturally occupy the same number of bytes. Without padding every coordinate to the curve's fixed byte width (32 bytes for P-256's 256-bit numbers), two different public keys could accidentally produce ambiguous, misaligned byte layouts — precisely Chapter 09, Section 7's field-boundary bug, now in a new context. Go's `big.Int.FillBytes` method does exactly this padding for us, writing a number's bytes into a fixed-size buffer, left-padded with zeroes if the number is smaller than the buffer:

```go
// marshalPublicKey turns an *ecdsa.PublicKey into GoChain's public
// key byte format: X and Y, each padded to the curve's byte size,
// concatenated with no other framing.
func marshalPublicKey(pub *ecdsa.PublicKey) []byte {
	byteLen := (pub.Curve.Params().BitSize + 7) / 8

	buf := make([]byte, 2*byteLen)
	pub.X.FillBytes(buf[:byteLen])
	pub.Y.FillBytes(buf[byteLen:])

	return buf
}
```

`(pub.Curve.Params().BitSize + 7) / 8` is a standard rounding-up trick: P-256's `BitSize` is 256, so `(256+7)/8` is exactly `32` bytes. This formula works unchanged for any curve size, which matters if a future chapter ever swaps curves. Reversing this — turning 64 raw bytes back into an `*ecdsa.PublicKey` — is exactly what `Verify()` needs to do internally, covered in Section 8.

`marshalPublicKey` is unexported (lowercase) because, so far, only code inside this same package needs it. That changes the moment a *signer* — something outside `gochain/crypto` — needs to write its own public key onto a fresh, not-yet-verified value, rather than just checking a signature against a key it already has in byte form. Chapter 33's `core.Transaction.Sign()` is exactly that case: having just signed with a private key, it needs to stamp the *matching* public key, in this same byte format, onto the transaction input it just signed. A one-line exported wrapper covers this:

```go
// PublicKeyToBytes exposes marshalPublicKey to callers outside this
// package. core.Transaction.Sign (Chapter 33) uses it to convert the
// signer's *ecdsa.PublicKey into GoChain's raw public-key byte format
// immediately after producing a signature.
func PublicKeyToBytes(pub *ecdsa.PublicKey) []byte {
	return marshalPublicKey(pub)
}
```

There's no new logic here at all — `PublicKeyToBytes` is `marshalPublicKey` under an exported name, kept as two functions only so the package's own internal code (Section 6's `NewKeyPair`) can keep calling the terse, private name, while external callers get a name that reads clearly at the call site.

One more thing worth naming precisely, since the field comment says it directly: this is the *uncompressed* form of a public key — both coordinates in full — as opposed to a *compressed* form (storing `X` plus a single extra bit indicating which of the two possible `Y` values matches it, since a curve equation generally has two `Y` solutions for a given `X`). Compressed keys are roughly half the size and are what many production wallets use to save space; GoChain uses the simpler uncompressed form throughout this course, favoring clarity over the extra byte-shaving.

---

## 6. Generating a Key Pair: NewKeyPair

With `curve` declared and `marshalPublicKey` ready, `NewKeyPair` itself is short:

```go
// NewKeyPair generates a brand-new ECDSA key pair on GoChain's chosen
// curve, using Go's cryptographically secure random source. It never
// uses math/rand, which is fast but predictable and therefore
// completely unsuitable here (Chapter 11, Section 7).
func NewKeyPair() (*KeyPair, error) {
	priv, err := ecdsa.GenerateKey(curve, rand.Reader)
	if err != nil {
		return nil, err
	}

	return &KeyPair{
		PrivateKey: priv,
		PublicKey:  marshalPublicKey(&priv.PublicKey),
	}, nil
}
```

`ecdsa.GenerateKey` does the actual work Chapter 11, Section 2 described as a single arrow: it draws a large, properly random number for the private key from `rand.Reader` (Go's `crypto/rand` package, backed by the operating system's cryptographically secure randomness source), and derives the matching public key point by applying the curve's point-multiplication rule. The returned `*ecdsa.PrivateKey` struct actually embeds the public key already (as `priv.PublicKey`), which is why `marshalPublicKey` takes `&priv.PublicKey` directly rather than needing a separate derivation step.

The `error` return exists because `rand.Reader` can, in principle, fail — the underlying operating system's entropy source could be unavailable, though this is exceptionally rare on any real machine. Section 11 covers exactly why GoChain propagates this error rather than silently ignoring it, matching the honest error-handling standard Chapter 07 set for this whole course.

Try it immediately:

```go
kp, err := crypto.NewKeyPair()
if err != nil {
	log.Fatal(err)
}
fmt.Printf("private key D: %x\n", kp.PrivateKey.D)
fmt.Printf("public key:    %x\n", kp.PublicKey)
```

`kp.PrivateKey.D` is the raw private key number itself (`D` for the private scalar Chapter 11, Section 6 called `d`) — printed here only to make the abstraction concrete on your own screen. GoChain never prints, logs, or transmits this value anywhere in real code; Chapter 11, Section 4's rule applies from this line of code onward.

---

## 7. Signing Data: Sign

`Sign` takes a private key and any data, and returns a signature:

```go
// Sign produces an ECDSA signature over the SHA-256 hash of data,
// using priv. It signs the hash, not the raw data (Chapter 12,
// Section 5), so the exact same function works whether data is ten
// bytes or ten megabytes.
func Sign(priv *ecdsa.PrivateKey, data []byte) ([]byte, error) {
	hash := Hash(data) // crypto.Hash, from Chapter 09

	signature, err := ecdsa.SignASN1(rand.Reader, priv, hash)
	if err != nil {
		return nil, err
	}

	return signature, nil
}
```

Three things worth calling out explicitly. First, `Sign` calls `Hash(data)` itself, rather than expecting the caller to hash first — this makes it impossible to accidentally sign raw, unhashed data by forgetting a step, and it means every caller in this course, from this chapter's demo to Chapter 33's transaction signing, hashes the exact same way, with SHA-256, every time. Second, `ecdsa.SignASN1` takes `rand.Reader` as its first argument — this is Chapter 12, Section 10's random nonce, generated fresh, correctly, and internally by Go's standard library every single time `Sign` is called, so GoChain never has to manage that value by hand or risk the Sony-style reuse bug Chapter 12 described. Third, `SignASN1` returns the `(r, s)` pair (Chapter 12, Section 2) already encoded into a single self-describing byte slice using a standard format called ASN.1 DER — which is exactly why `Sign`'s return type is a plain `[]byte`, matching the contract, with no separate `r` and `s` fields to manage.

---

## 8. Verifying a Signature: Verify

`Verify` is `Sign`'s mirror image, and, matching Chapter 12, Section 6's diagram exactly, it never receives a private key at all:

```go
// Verify reports whether signature is a valid ECDSA signature over
// data's SHA-256 hash, produced by the private key matching pubKey.
func Verify(pubKey []byte, data, signature []byte) bool {
	pub, err := unmarshalPublicKey(pubKey)
	if err != nil {
		return false
	}

	hash := Hash(data)

	return ecdsa.VerifyASN1(pub, hash, signature)
}

// unmarshalPublicKey reverses marshalPublicKey, rebuilding an
// *ecdsa.PublicKey from GoChain's raw X||Y byte format.
func unmarshalPublicKey(pubKey []byte) (*ecdsa.PublicKey, error) {
	byteLen := (curve.Params().BitSize + 7) / 8

	if len(pubKey) != 2*byteLen {
		return nil, errors.New("crypto: invalid public key length")
	}

	x := new(big.Int).SetBytes(pubKey[:byteLen])
	y := new(big.Int).SetBytes(pubKey[byteLen:])

	return &ecdsa.PublicKey{Curve: curve, X: x, Y: y}, nil
}
```

`unmarshalPublicKey` reverses Section 5's byte layout exactly: split the 64-byte slice at its midpoint, and reconstruct `X` and `Y` from each half using `big.Int.SetBytes` (the inverse of `FillBytes`). If `pubKey` is not exactly the expected length, `unmarshalPublicKey` returns an error immediately, and `Verify` treats that as `false` rather than panicking — a malformed or corrupted public key should never crash a node, it should simply fail to verify, and Section 11 explains why that distinction matters.

`ecdsa.VerifyASN1` performs the actual mathematical check Chapter 12, Section 6 diagrammed conceptually, and returns a plain `bool` — exactly matching `Verify`'s own return type, with no error to check at all. This is a deliberate, important asymmetry worth noticing: `Sign` can fail (a broken random source, an invalid key), so it returns an `error`; `Verify` cannot meaningfully "fail" in the same sense — a signature is either valid or it is not, and an unparseable public key or a malformed signature is simply one more way to be invalid, not a separate failure mode a caller needs to branch on differently.

---

## 9. Full Worked Example: Alice Pays Bob 5 Gochips

Here is the complete pipeline from Chapter 12, Section 7, running as real Go code:

```go
// main.go
package main

import (
	"encoding/hex"
	"fmt"
	"log"

	"github.com/you/gochain/crypto"
)

func main() {
	// Alice generates her key pair -- once, ever, for this identity.
	alice, err := crypto.NewKeyPair()
	if err != nil {
		log.Fatal(err)
	}

	// The exact data Alice wants to authorize.
	message := []byte("pay Bob 5 gochips")

	// Alice signs it with her PRIVATE key.
	signature, err := crypto.Sign(alice.PrivateKey, message)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("message:         ", string(message))
	fmt.Println("Alice public key:", hex.EncodeToString(alice.PublicKey))
	fmt.Println("signature:       ", hex.EncodeToString(signature))

	// Anyone -- Bob, a stranger, every GoChain node -- can verify it
	// using ONLY Alice's PUBLIC key. Her private key never appears
	// again after the Sign call above.
	valid := crypto.Verify(alice.PublicKey, message, signature)
	fmt.Println("valid?           ", valid)
}
```

Running this prints something like:

```
message:          pay Bob 5 gochips
Alice public key: 3f2a91c4...d817ec4b   (64 bytes, 128 hex characters)
signature:        304402203a1f...        (roughly 70 bytes, ASN.1-encoded)
valid?            true
```

Every one of Chapter 12's promises is now something you watched happen on your own machine: Alice's private key (`alice.PrivateKey`) was used exactly once, inside `Sign`, and never appears again; `Verify` used only `alice.PublicKey`, `message`, and `signature`; and the result, `true`, is a claim any other computer running this same code, with the same three inputs, would reach independently — no server, no coordination, no trusted third party, exactly Chapter 11, Section 1's opening promise for this entire volume.

---

## 10. Corrupting One Byte — Watching Verification Fail

Now reproduce Chapter 12, Section 9's byte-corruption walkthrough for real. Extend the example with a tampered copy of the message, verified against Alice's original, genuine signature:

```go
	// An attacker intercepts the message and changes ONE detail --
	// the amount -- while keeping Alice's original signature attached.
	tampered := []byte("pay Bob 9 gochips")

	stillValid := crypto.Verify(alice.PublicKey, tampered, signature)
	fmt.Println("tampered valid?  ", stillValid)
```

This prints:

```
tampered valid?   false
```

Nothing about `Verify`'s code path changed between this call and the genuine one in Section 9 — the exact same function, the exact same public key, and the exact same signature bytes were used both times. The only thing that changed was one character inside `data`, and that was enough: `Hash(tampered)` produces a completely different 32-byte digest than `Hash(message)` did (the avalanche effect, Chapter 08, Section 3), so the mathematical relationship `ecdsa.VerifyASN1` checks between that digest, `alice.PublicKey`, and `signature` no longer holds. Try changing even a single character — a lowercase `b` to `B` in `"Bob"`, or an extra trailing space — and confirm for yourself that `stillValid` remains `false` every time. There is no "close enough" in this system, exactly as Chapter 12, Section 9 predicted.

```
  ┌──────────────────────┐     ┌──────────────────────┐
  │ "pay Bob 5 gochips"   │     │ "pay Bob 9 gochips"   │
  │ Alice's real message  │     │ attacker's tampered   │
  │                        │     │ version                │
  └──────────┬───────────┘     └──────────┬───────────┘
             │ Hash()                       │ Hash()
             ▼                              ▼
   hash_A (32 bytes)              hash_B (32 bytes, unrelated)
             │                              │
             └──────────────┬───────────────┘
                             │
              signature (r, s) was produced
              for hash_A specifically
                             │
             ┌───────────────┴───────────────┐
             ▼                                ▼
   Verify(pub, hash_A, sig)          Verify(pub, hash_B, sig)
        = true                            = false
```

---

## 11. Handling Errors Correctly

Two small but important Go habits are worth calling out explicitly, because getting them wrong would quietly undermine everything this chapter just built.

**Always check `NewKeyPair`'s and `Sign`'s errors.** Both can fail (in practice, almost never, since they only fail if the underlying secure random source is unavailable), but "almost never" is not "never," and silently ignoring an error here could mean generating a key pair from a broken or partially-failed random source — reintroducing exactly the weak-randomness risk Chapter 11, Section 7 and Chapter 12, Section 10 spent so much effort warning about. Go's `if err != nil { return nil, err }` pattern, used consistently throughout this section, propagates that failure to whoever called `NewKeyPair` or `Sign`, rather than pretending everything succeeded.

**Never let `Verify` panic.** A node running GoChain will call `Verify` on signatures and public keys it received over the network (starting in Volume 7) — data from strangers, some of it possibly malformed by accident or by a deliberately malicious peer. `unmarshalPublicKey`'s length check in Section 8 is exactly the guard that keeps a malformed public key from ever reaching `ecdsa.VerifyASN1` in a shape it cannot handle; returning `false` for anything that fails to parse, rather than letting a panic escape, means one bad piece of network data can never crash an entire node. This is a preview of a much larger theme Volume 7 returns to repeatedly: code that processes untrusted input must fail *safely*, never catastrophically.

---

## 12. Testing crypto.Sign and crypto.Verify

Following Chapter 09, Section 9's testing pattern, `crypto`'s signature tests should directly verify the properties this chapter and Chapter 12 both promised: a genuine signature verifies, a tampered message does not, the wrong public key does not verify a signature meant for someone else, and two signatures of the same data are never identical:

```go
// crypto/keys_test.go
package crypto

import (
	"bytes"
	"testing"
)

func TestSignAndVerify_ValidSignature(t *testing.T) {
	kp, err := NewKeyPair()
	if err != nil {
		t.Fatalf("NewKeyPair failed: %v", err)
	}

	data := []byte("pay Bob 5 gochips")
	sig, err := Sign(kp.PrivateKey, data)
	if err != nil {
		t.Fatalf("Sign failed: %v", err)
	}

	if !Verify(kp.PublicKey, data, sig) {
		t.Fatal("expected a genuine signature to verify")
	}
}

func TestVerify_TamperedDataFails(t *testing.T) {
	kp, _ := NewKeyPair()
	data := []byte("pay Bob 5 gochips")
	sig, _ := Sign(kp.PrivateKey, data)

	tampered := []byte("pay Bob 9 gochips")
	if Verify(kp.PublicKey, tampered, sig) {
		t.Fatal("expected verification to fail for tampered data")
	}
}

func TestVerify_WrongPublicKeyFails(t *testing.T) {
	alice, _ := NewKeyPair()
	mallory, _ := NewKeyPair() // an unrelated key pair
	data := []byte("pay Bob 5 gochips")
	sig, _ := Sign(alice.PrivateKey, data)

	if Verify(mallory.PublicKey, data, sig) {
		t.Fatal("expected verification to fail against an unrelated public key")
	}
}

func TestSign_DifferentSignaturesEachTime(t *testing.T) {
	kp, _ := NewKeyPair()
	data := []byte("pay Bob 5 gochips")

	sig1, _ := Sign(kp.PrivateKey, data)
	sig2, _ := Sign(kp.PrivateKey, data)

	if bytes.Equal(sig1, sig2) {
		t.Fatal("expected two signatures of the same data to differ (Chapter 12, Section 10)")
	}
	if !Verify(kp.PublicKey, data, sig1) || !Verify(kp.PublicKey, data, sig2) {
		t.Fatal("expected both signatures to verify successfully")
	}
}
```

`TestSign_DifferentSignaturesEachTime` directly exercises Chapter 12, Section 10's nonce discussion: it would fail immediately if `Sign` ever reused the same nonce (or, worse, used a fixed one, Sony-style) across calls, since two identical signatures would then be produced for identical input — this test is a small, permanent guardrail against ever accidentally regressing into that specific, historically real bug. Run `go test ./crypto/...` and all four should pass.

---

## 13. Where This Fits in GoChain

`crypto.NewKeyPair`, `crypto.Sign`, and `crypto.Verify` are the exact three functions Volume 5's `core.Transaction.Sign()` and `core.Transaction.Verify()` (Chapter 33) call underneath, and the exact three functions Volume 6's `wallet` package (starting Chapter 38) builds an entire user-facing wallet experience around. Nothing about signing or verifying changes from here forward — every later chapter is either calling these three functions directly, or building conveniences on top of them (an address in Chapter 14, an HD wallet's tree of key pairs in Chapter 38, a hardware-backed signer in Chapter 41). Chapter 14 takes the 64-byte `PublicKey` this chapter produces and turns it into something a human being can actually read, write down, and type correctly.

---

## Summary

- `gochain/crypto` now exports `KeyPair` (a private key plus its public key, stored as raw `X||Y` bytes), `NewKeyPair()`, `Sign()`, `Verify()`, and `PublicKeyToBytes()` — real, running code implementing Chapter 12's sign/verify pipeline.
- `crypto/elliptic` provides the curve (GoChain uses `elliptic.P256()`); `crypto/ecdsa` provides the signing and verifying algorithm built on top of whatever curve it is given.
- A public key is stored as a 64-byte slice: the `X` and `Y` coordinates, each padded to a fixed width with `big.Int.FillBytes`, concatenated with no other framing — the uncompressed form, chosen for clarity over the smaller compressed alternative.
- `NewKeyPair` calls `ecdsa.GenerateKey` with `crypto/rand`'s cryptographically secure random source, never `math/rand`, matching Chapter 11, Section 7's requirement for a properly random private key.
- `Sign` hashes its input with `crypto.Hash` before signing, so it always signs a fixed-size digest rather than raw data of arbitrary length, exactly as Chapter 12, Section 5 required.
- `Verify` never accepts or needs a private key, returns a plain `bool` (not an error), and fails safely — via an explicit length check, not a panic — when handed a malformed public key.
- The worked Alice/Bob example proved every claim from Chapter 12 concretely: a genuine signature verifies as `true`; changing even one byte of the signed message makes the exact same signature verify as `false`.
- Tests for this package should directly check the properties this course cares about: valid signatures verify, tampered data does not, the wrong public key does not verify, and two signatures of identical data are never byte-for-byte identical.

---

## Exercises

### Easy

1. Run this chapter's Alice/Bob example yourself, then print `len(alice.PublicKey)` and `len(signature)`. Confirm the public key is exactly 64 bytes (Section 5) and explain, in 2-3 sentences, why the signature's length is not perfectly fixed the way the public key's is (hint: consider what ASN.1 DER encoding, mentioned in Section 7, does with numbers of varying size).

2. Modify the worked example so that Bob, not Alice, generates a key pair and verifies Alice's signature. Confirm Bob's code never needs access to `alice.PrivateKey` at any point, and explain in your own words why that is exactly the property Chapter 11, Section 4 required.

3. `unmarshalPublicKey` (Section 8) returns an error if `pubKey` is not exactly 64 bytes. Write a small test that calls `Verify` with a deliberately truncated public key (say, `alice.PublicKey[:32]`) and confirm it returns `false` rather than panicking.

### Medium

4. Section 11 argued that `Verify` should never panic, even on malformed input. Deliberately pass a `nil` slice as `signature` to `Verify` and confirm the behavior. Then explain, in 150-200 words, why a node accepting data from untrusted network peers (previewed for Volume 7) especially cannot afford a signature-checking function that panics on bad input.

5. Write a small program that generates two independent key pairs (Alice's and Mallory's), has Alice sign a message, and then attempts to verify Alice's signature using Mallory's public key instead of Alice's. Confirm the result is `false`, and explain in your own words, referencing Chapter 12, Section 6, exactly why the verification math cannot succeed here even though `signature` and `data` are both entirely genuine.

6. Section 6 pads a public key's `X` and `Y` coordinates to a fixed byte width using `big.Int.FillBytes`. Research what `FillBytes` does if the destination buffer is *too small* for the number being written (check the Go standard library documentation), and explain why `marshalPublicKey`'s use of `(BitSize + 7) / 8` guarantees this never happens for a properly generated P-256 key.

### Hard

7. Benchmark `NewKeyPair`, `Sign`, and `Verify` using Go's `testing.B` benchmarking support (Chapter 07 introduced `go test`; research `func BenchmarkXxx(b *testing.B)` if you have not written a Go benchmark before). Report which of the three operations is slowest, and offer a hypothesis for why, grounded in what each function actually computes (Sections 6 through 8).

8. This chapter's `Sign` always hashes with SHA-256 via `crypto.Hash`. Design (as Go code or detailed pseudocode) a version of `Sign` and `Verify` that accepts an already-computed hash directly, instead of raw data, avoiding a redundant hash computation when a caller (like a future `core.Transaction`) has already hashed its own contents for other reasons. Discuss, in 200-300 words, one risk this redesign introduces that the current chapter's design avoids by hashing internally every time.

9. Research Go's `crypto/ecdh` package (added in Go 1.20) and how it relates to, but differs from, `crypto/ecdsa`. Write an explanation (250-350 words) of what problem Elliptic Curve Diffie-Hellman (ECDH) key exchange solves that is fundamentally different from what ECDSA solves, even though both are built on the exact same elliptic curve concepts from Chapter 11, and explain why GoChain's wallet and transaction signing needs ECDSA specifically, not ECDH.
