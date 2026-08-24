# Chapter 39: Implementing Mnemonic Seed Phrases

Chapter 38 explained BIP-39 conceptually: random bytes become a checksummed, human-writable phrase, which stretches back into a seed. This chapter turns that into real, tested Go code — `wallet.NewMnemonic()` and `wallet.SeedFromMnemonic()` — starting from the one thing every wallet's security ultimately rests on: where the randomness comes from in the first place.

## Table of Contents

1. [What We're Building](#1-what-were-building)
2. [Secure Randomness With crypto/rand](#2-secure-randomness-with-cryptorand)
3. [The BIP-39 Wordlist: 2048 Words, 11 Bits Each](#3-the-bip-39-wordlist-2048-words-11-bits-each)
4. [From Entropy to Words: The Algorithm Step by Step](#4-from-entropy-to-words-the-algorithm-step-by-step)
5. [Implementing NewMnemonic()](#5-implementing-newmnemonic)
6. [From Mnemonic Back to a Seed: SeedFromMnemonic()](#6-from-mnemonic-back-to-a-seed-seedfrommnemonic)
7. [Round-Trip Tests: Generate, Then Recover](#7-round-trip-tests-generate-then-recover)
8. [How a Real Wordlist Gets Loaded](#8-how-a-real-wordlist-gets-loaded)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What We're Building

By the end of this chapter, `gochain/wallet` exposes exactly two new functions, matching this course's shared contract:

```go
func NewMnemonic() (string, error)
func SeedFromMnemonic(mnemonic, passphrase string) []byte
```

`NewMnemonic()` generates fresh, cryptographically secure randomness and returns a twelve-word phrase like `"witch collapse practice feed shame open maple north rescue fee coyote crumble"`. `SeedFromMnemonic()` takes that phrase (plus an optional passphrase) and deterministically reproduces the 64-byte seed that Chapter 40's HD derivation will build its entire key tree on top of. Run the same mnemonic and passphrase through `SeedFromMnemonic()` next year, on a different machine, and you get back the identical seed, byte for byte — that determinism is the entire point of Chapter 38's design.

---

## 2. Secure Randomness With crypto/rand

Everything downstream — the mnemonic, the seed, every derived private key — is only as unpredictable as the very first random bytes we generate. This is not a place to cut corners: Go's `math/rand` package is fine for shuffling a slice or picking a random test fixture, but it is **not safe for anything security-sensitive**, because its output is predictable if an attacker can guess or brute-force its internal seed. GoChain uses `crypto/rand` instead, which reads from the operating system's cryptographically secure randomness source (`/dev/urandom` on Linux/macOS, `CryptGenRandom` on Windows) — the same source real wallets, TLS libraries, and every serious cryptographic system rely on.

```go
package wallet

import "crypto/rand"

// secureEntropy returns n cryptographically secure random bytes, read
// directly from the OS's randomness source. This is the only place in
// the entire mnemonic pipeline where fresh randomness is introduced —
// everything after this function is a deterministic transformation.
func secureEntropy(n int) ([]byte, error) {
	b := make([]byte, n)
	// rand.Read blocks until it has filled b with real entropy; on any
	// error (extremely rare, but possible on some constrained systems)
	// we propagate it rather than silently falling back to something
	// weaker.
	if _, err := rand.Read(b); err != nil {
		return nil, err
	}
	return b, nil
}
```

`secureEntropy` is explained by name: it asks `crypto/rand.Read` to fill a byte slice of length `n` with secure random bytes, and returns an error rather than partial or predictable data if the read ever fails. Every mnemonic GoChain ever generates starts here.

---

## 3. The BIP-39 Wordlist: 2048 Words, 11 Bits Each

BIP-39's word list contains **exactly 2048 words** — no more, no fewer — because 2048 is 2¹¹. That is not a coincidence: it means every single word in the list can be represented by an **11-bit number** (0 through 2047), so converting a stream of random bits into words is just a matter of chopping the bit stream into 11-bit chunks and looking each chunk up as an index.

The real list (defined by BIP-39 for English) is designed with care beyond just being 2048 arbitrary words: every word is unique in its first four letters, so a wallet's autocomplete can identify a word from just a few keystrokes; visually or phonetically similar words (like "build" and "built") are deliberately avoided; and offensive or ambiguous words are excluded. Here is a representative excerpt — the real list starts alphabetically and continues to word index 2047:

```
index 0000: abandon
index 0001: ability
index 0002: able
index 0003: about
index 0004: above
index 0005: absent
index 0006: absorb
index 0007: abstract
     ...
index 0102: audit
     ...
index 1005: maple
     ...
index 1481: rescue
     ...
index 2044: zebra
index 2045: zero
index 2046: zone
index 2047: zoo
```

A real implementation does not retype all 2048 words by hand — it loads them from a data file bundled with the program, which Section 8 shows properly. For the algorithm walkthrough in this chapter, treat the wordlist as a fixed `[]string` of length 2048, indexed exactly like the excerpt above.

---

## 4. From Entropy to Words: The Algorithm Step by Step

Here is BIP-39's mnemonic-generation algorithm, stated precisely before any Go code, using 128 bits of entropy (the twelve-word case) as the running example:

1. **Generate ENT bits of entropy.** BIP-39 allows ENT to be 128, 160, 192, 224, or 256 bits; GoChain uses 128 (12 words) by default and supports 256 (24 words) for users who want extra margin.
2. **Compute a checksum.** Take SHA-256 of the entropy bytes, and keep only the first `ENT / 32` bits of the resulting hash. For 128 bits of entropy, that's `128 / 32 = 4` checksum bits.
3. **Append the checksum bits to the entropy bits.** `128 + 4 = 132` total bits.
4. **Split the combined bits into groups of 11.** `132 / 11 = 12` groups — exactly twelve words. (This is precisely why the checksum length is chosen as `ENT / 32`: it is the smallest number of extra bits that makes the total divide evenly by 11.)
5. **Look up each 11-bit group as an index into the 2048-word list**, producing the words in order.

```
128 bits entropy                              4 bits checksum
+----------------------------------------+    +----+
| 10110100 11001010 ... (128 bits total)  |    |1010|
+----------------------------------------+    +----+
                    |                             |
                    +--------------+--------------+
                                   v
                 132 bits, split into twelve 11-bit groups
        +-----------+-----------+-----------+-----+-----------+
        | 01101001100| 10101100 1| 011001... | ... | ...101010 |
        +-----------+-----------+-----------+-----+-----------+
              |            |           |               |
              v            v           v               v
           word[857]   word[1364]  word[...]  ...   word[...]
```

The checksum in step 2-3 is exactly what makes a typo detectable: if a user mistypes even one word when recovering their wallet, the last 4 (or 8, for 24 words) bits almost certainly no longer match the SHA-256 checksum recomputed from the other words, and recovery can fail loudly instead of silently producing a different, wrong wallet.

---

## 5. Implementing NewMnemonic()

```go
package wallet

import (
	"crypto/sha256"
	"errors"
	"strings"
)

// entropyBits is the amount of entropy GoChain generates by default.
// 128 bits produces a 12-word mnemonic; BIP-39 also permits 160, 192,
// 224, or 256 bits (24 words) for users who want a larger security
// margin at the cost of a longer phrase to write down and type back.
const entropyBits = 128

// NewMnemonic generates fresh secure entropy and encodes it as a
// BIP-39 mnemonic: a sequence of words from the standard 2048-word
// list, with a trailing checksum baked in so a single mistyped word
// is detectable on recovery.
func NewMnemonic() (string, error) {
	entropy, err := secureEntropy(entropyBits / 8)
	if err != nil {
		return "", err
	}
	return mnemonicFromEntropy(entropy)
}

// mnemonicFromEntropy implements BIP-39 section 4's algorithm exactly
// as described in Section 4: append a checksum, split into 11-bit
// groups, and map each group to a word.
func mnemonicFromEntropy(entropy []byte) (string, error) {
	entBits := len(entropy) * 8
	checksumBits := entBits / 32 // e.g. 128/32 = 4 bits for a 12-word phrase

	// SHA-256 of the raw entropy; we only need its first checksumBits.
	hash := sha256.Sum256(entropy)

	// Combine entropy bits and checksum bits into one bit string. We
	// use a simple string of '0'/'1' characters here for clarity —
	// real production code would pack this into a bitset for speed,
	// but the word count here is tiny (at most 264 bits), so clarity
	// wins over micro-optimization.
	bits := bytesToBits(entropy) + bytesToBits(hash[:])[:checksumBits]

	totalWords := len(bits) / 11
	words := make([]string, totalWords)
	for i := 0; i < totalWords; i++ {
		chunk := bits[i*11 : (i+1)*11]
		index := bitsToInt(chunk)
		if index >= len(wordlist) {
			return "", errors.New("wallet: word index out of range, wordlist corrupt")
		}
		words[i] = wordlist[index]
	}
	return strings.Join(words, " "), nil
}

// bytesToBits renders each byte as 8 '0'/'1' characters, most
// significant bit first, matching BIP-39's bit ordering.
func bytesToBits(b []byte) string {
	var sb strings.Builder
	for _, by := range b {
		for i := 7; i >= 0; i-- {
			if by&(1<<uint(i)) != 0 {
				sb.WriteByte('1')
			} else {
				sb.WriteByte('0')
			}
		}
	}
	return sb.String()
}

// bitsToInt parses an 11-character '0'/'1' string as an unsigned
// integer, used to turn each chunk into a wordlist index (0-2047).
func bitsToInt(bits string) int {
	n := 0
	for _, c := range bits {
		n <<= 1
		if c == '1' {
			n |= 1
		}
	}
	return n
}
```

Walking through this by name: `NewMnemonic()` is the public entry point — it pulls fresh secure entropy and hands it to `mnemonicFromEntropy`. `mnemonicFromEntropy` implements Section 4's five steps directly: it hashes the entropy, keeps just enough checksum bits (`entBits / 32`), concatenates entropy bits and checksum bits into one long bit string, chops that string into 11-character groups, and converts each group into a wordlist index. `bytesToBits` and `bitsToInt` are small helpers that convert between raw bytes and the bit-string representation the algorithm reasons about most naturally; a production-grade implementation would likely use bit-shifting on integers directly for speed, but the byte counts here are small enough that clarity is the right trade-off for a learning codebase.

---

## 6. From Mnemonic Back to a Seed: SeedFromMnemonic()

Generating a mnemonic is only half the story — Chapter 40's HD derivation needs a **seed**, not a sentence of words. BIP-39 defines a specific way to stretch a mnemonic (plus an optional passphrase) into a 64-byte seed: run it through **PBKDF2**, a password-based key derivation function, using HMAC-SHA512 as the underlying hash, with 2048 iterations.

The "2048 iterations" detail matters for a specific reason: it deliberately makes computing a seed from a candidate mnemonic *slow* (a few milliseconds, not micro­seconds), which raises the cost of an attacker brute-forcing possible mnemonics or passphrases. A legitimate user pays this cost once, when opening their wallet; an attacker trying millions of guesses pays it millions of times over.

```go
package wallet

import (
	"golang.org/x/crypto/pbkdf2"
	"crypto/sha512"
)

// seedFromMnemonicIterations follows BIP-39's specified value exactly;
// changing it would make GoChain wallets unable to recover seeds
// generated by any other BIP-39-compliant wallet, and vice versa.
const seedFromMnemonicIterations = 2048

// seedLength is 64 bytes, per BIP-39 — enough entropy to seed an
// entire BIP-32 key tree (Chapter 40) safely.
const seedLength = 64

// SeedFromMnemonic deterministically derives a 64-byte seed from a
// mnemonic phrase and an optional passphrase. The same mnemonic and
// passphrase always produce the same seed; a different passphrase
// (including the empty string vs. a non-empty one) produces a
// completely different, unrelated seed.
func SeedFromMnemonic(mnemonic, passphrase string) []byte {
	// BIP-39 salts the KDF with the literal string "mnemonic" followed
	// by the user's passphrase, so that even an empty passphrase still
	// produces a well-defined, standard salt.
	salt := "mnemonic" + passphrase

	return pbkdf2.Key(
		[]byte(mnemonic),
		[]byte(salt),
		seedFromMnemonicIterations,
		seedLength,
		sha512.New,
	)
}
```

`SeedFromMnemonic` is a thin, deliberate wrapper around `pbkdf2.Key`: the mnemonic sentence is the "password" being stretched, the salt is the fixed string `"mnemonic"` concatenated with the caller's passphrase (so two users with the same words but different passphrases get entirely different seeds — a real security feature, not an edge case to guard against), and the iteration count and output length are BIP-39's standard values, not tunable knobs. Getting any of these three constants wrong would silently produce seeds incompatible with every other BIP-39 wallet in existence, so GoChain treats them as fixed, not configurable.

---

## 7. Round-Trip Tests: Generate, Then Recover

The whole point of a seed phrase is that generating it once and recovering it later produce the *same* result. We test exactly that:

```go
package wallet

import (
	"bytes"
	"strings"
	"testing"
)

func TestNewMnemonicProducesTwelveWords(t *testing.T) {
	mnemonic, err := NewMnemonic()
	if err != nil {
		t.Fatalf("NewMnemonic failed: %v", err)
	}

	words := strings.Fields(mnemonic)
	if len(words) != 12 {
		t.Fatalf("expected 12 words, got %d: %q", len(words), mnemonic)
	}

	// Every word must actually be in our wordlist — a stray word here
	// would mean our bit-to-index math has a bug.
	valid := make(map[string]bool, len(wordlist))
	for _, w := range wordlist {
		valid[w] = true
	}
	for _, w := range words {
		if !valid[w] {
			t.Errorf("word %q is not in the BIP-39 wordlist", w)
		}
	}
}

func TestSeedFromMnemonicIsDeterministic(t *testing.T) {
	mnemonic, err := NewMnemonic()
	if err != nil {
		t.Fatalf("NewMnemonic failed: %v", err)
	}

	// Deriving the seed twice, from the same mnemonic and passphrase,
	// must produce byte-for-byte identical output every time — this is
	// the entire premise of "one backup, regenerate forever."
	seedA := SeedFromMnemonic(mnemonic, "")
	seedB := SeedFromMnemonic(mnemonic, "")

	if !bytes.Equal(seedA, seedB) {
		t.Fatal("SeedFromMnemonic is not deterministic for identical inputs")
	}
	if len(seedA) != 64 {
		t.Fatalf("expected a 64-byte seed, got %d bytes", len(seedA))
	}
}

func TestDifferentPassphrasesProduceDifferentSeeds(t *testing.T) {
	mnemonic, err := NewMnemonic()
	if err != nil {
		t.Fatalf("NewMnemonic failed: %v", err)
	}

	seedNoPass := SeedFromMnemonic(mnemonic, "")
	seedWithPass := SeedFromMnemonic(mnemonic, "extra-secret")

	if bytes.Equal(seedNoPass, seedWithPass) {
		t.Fatal("expected different passphrases to produce different seeds")
	}
}

func TestTwoMnemonicsAreVeryUnlikelyToCollide(t *testing.T) {
	a, err := NewMnemonic()
	if err != nil {
		t.Fatalf("NewMnemonic failed: %v", err)
	}
	b, err := NewMnemonic()
	if err != nil {
		t.Fatalf("NewMnemonic failed: %v", err)
	}
	if a == b {
		t.Fatal("two independently generated mnemonics were identical — entropy source may be broken")
	}
}
```

Each test explains itself by name: `TestNewMnemonicProducesTwelveWords` checks the mechanical shape of the output (right word count, every word a real wordlist entry). `TestSeedFromMnemonicIsDeterministic` is the most important test in this chapter — it directly verifies the "same input, same output, every time" property Chapter 38 promised. `TestDifferentPassphrasesProduceDifferentSeeds` guards the passphrase feature specifically, and `TestTwoMnemonicsAreVeryUnlikelyToCollide` is a lightweight sanity check that `crypto/rand` is actually being consulted (a hard-coded or broken entropy source would fail this test almost immediately).

---

## 8. How a Real Wordlist Gets Loaded

Section 3 showed a small excerpt of the wordlist for illustration, but production code needs the real, complete, official 2048-word list — retyping it by hand would be both tedious and dangerously error-prone (a single wrong word breaks compatibility with every other BIP-39 wallet). The standard approach is to ship the list as a plain text file, one word per line, and load it with Go's `embed` package so it becomes part of the compiled binary rather than a file that could go missing at runtime:

```go
package wallet

import (
	_ "embed"
	"strings"
)

// wordlistData embeds the official BIP-39 English wordlist (2048
// lines, one word each) directly into the compiled binary at build
// time. The file itself — wordlist.txt — is the standard list
// published alongside the BIP-39 specification, copied verbatim.
//go:embed wordlist.txt
var wordlistData string

// wordlist is the parsed, in-memory form used by mnemonicFromEntropy
// and its tests. Parsing happens once, at package init, rather than
// on every mnemonic generated.
var wordlist = strings.Fields(wordlistData)

func init() {
	if len(wordlist) != 2048 {
		// A corrupt or truncated wordlist file is a serious problem —
		// fail loudly at startup rather than generating subtly wrong
		// mnemonics later.
		panic("wallet: wordlist.txt does not contain exactly 2048 words")
	}
}
```

`//go:embed wordlist.txt` tells the Go compiler to bundle the contents of `wordlist.txt` — sitting alongside this source file — directly into the binary as the string `wordlistData`, so the program never depends on that file existing at runtime on whatever machine it's deployed to. The package-level `init()` function runs once when the package is first loaded and immediately verifies the list has exactly 2048 entries, failing fast and loudly if it doesn't, rather than letting a subtly corrupted wordlist silently generate mnemonics that would confuse other BIP-39 software. In a real GoChain checkout, `wordlist.txt` would be the official, unmodified BIP-39 English word list published alongside the specification — copying it verbatim (not retyping it) is the only safe way to guarantee interoperability with every other BIP-39-compliant wallet.

---

## Summary

- `crypto/rand`, not `math/rand`, is the only acceptable entropy source for anything security-sensitive, because it draws from the OS's cryptographically secure randomness pool.
- BIP-39's wordlist has exactly 2048 words because 2048 = 2¹¹, letting every word be addressed by an 11-bit index carved directly out of a bit stream.
- The mnemonic algorithm appends a `ENT / 32`-bit SHA-256 checksum to the raw entropy before splitting into 11-bit word indices — this checksum is what makes a single mistyped word detectable on recovery.
- `NewMnemonic()` generates fresh entropy and applies that algorithm; `SeedFromMnemonic()` stretches the resulting phrase (plus an optional passphrase) into a 64-byte seed via PBKDF2-HMAC-SHA512 with 2048 iterations, exactly per BIP-39.
- The PBKDF2 iteration count is deliberately expensive-but-tolerable, to raise the cost of brute-forcing candidate mnemonics or passphrases.
- Different passphrases combined with the same words produce entirely different, unrelated seeds — a real feature, not an edge case.
- A real wordlist is embedded from a verbatim copy of the official file via `//go:embed`, never retyped or reconstructed by hand.
- Round-trip tests — generate, then re-derive — are the correctness bar for this entire chapter: the same mnemonic and passphrase must always produce the identical seed.

---

## Exercises

### Easy

1. **Run `NewMnemonic()` twice in a small `main.go`** and print both phrases. Confirm by eye that they share no words in the same position, then explain in one or two sentences why this is expected given Section 2's entropy source.

2. **By hand, compute how many checksum bits and total words** BIP-39 would produce for 160 bits and for 256 bits of entropy, showing your division work the same way Section 4 did for 128 bits.

3. **Modify `TestDifferentPassphrasesProduceDifferentSeeds`** to also assert that the two resulting seeds are each still exactly 64 bytes long, and explain why that length check matters even though the test already checks the seeds differ.

### Medium

4. **Implement 24-word mnemonic support** by adding a second exported function, `NewMnemonic256() (string, error)`, that generates 256 bits of entropy instead of 128 and reuses `mnemonicFromEntropy` unchanged. Write a test confirming it produces exactly 24 words.

5. **Deliberately corrupt one word** in a mnemonic your test generates (swap it for an unrelated valid wordlist word) and write a `ValidateMnemonicChecksum(mnemonic string) bool` function that recomputes the checksum bits from the first 11 words and compares them against the last word's checksum bits, returning `false` for the corrupted phrase and `true` for the original.

6. **Benchmark `SeedFromMnemonic`** using Go's `testing.B` benchmark support, and report how many derivations per second your machine can perform. Explain, using Section 6's reasoning about brute-force cost, why a *slower* number here is arguably a *better* security property, up to a point.

### Hard

7. **Implement a minimal, from-scratch PBKDF2** (without using `golang.org/x/crypto/pbkdf2`) that reproduces BIP-39's exact seed derivation using only `crypto/hmac` and `crypto/sha512` from the standard library, and verify it produces identical output to the `pbkdf2.Key` version for the same mnemonic and passphrase.

8. **Investigate what happens if `wordlist.txt` contains a duplicate word** at two different indices. Explain concretely which property from this chapter would silently break, and design a stronger `init()` check that would catch a duplicate-word bug before it ever reached a user.

9. **Research BIP-39's non-English wordlists** (conceptually — no need to implement one) and explain what additional complexity a Japanese wordlist introduces for the mnemonic-to-seed algorithm specifically (hint: look into how BIP-39 normalizes text before hashing, and why that step matters more for some languages than others). Write 200-300 words summarizing your findings and what GoChain would need to change to support a second language.
