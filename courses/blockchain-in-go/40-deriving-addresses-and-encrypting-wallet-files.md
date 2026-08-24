# Chapter 40: Deriving Addresses and Encrypting Wallet Files

Chapter 39 gave you a 64-byte seed you can regenerate from twelve words at will. This chapter turns that seed into an actual tree of usable addresses — `wallet.NewHDWallet()` and `wallet.DeriveWallet()` — and then closes the loop on a problem we've been quietly ignoring since Chapter 36: an unencrypted wallet file sitting on disk is a complete, standing invitation to anyone who gets access to that disk. We fix that with `wallet.EncryptWalletFile()` and `wallet.DecryptWalletFile()`.

## Table of Contents

1. [What We're Building](#1-what-were-building)
2. [From Seed to Master Key](#2-from-seed-to-master-key)
3. [The Child Derivation Step](#3-the-child-derivation-step)
4. [GoChain's Fixed BIP-44 Path](#4-gochains-fixed-bip-44-path)
5. [Implementing DeriveWallet()](#5-implementing-derivewallet)
6. [Trying It: One Seed, Many Addresses](#6-trying-it-one-seed-many-addresses)
7. [The Problem With Plaintext Wallet Files](#7-the-problem-with-plaintext-wallet-files)
8. [Password Stretching With scrypt](#8-password-stretching-with-scrypt)
9. [Implementing EncryptWalletFile() and DecryptWalletFile()](#9-implementing-encryptwalletfile-and-decryptwalletfile)
10. [Testing the Full Round Trip](#10-testing-the-full-round-trip)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What We're Building

Two small groups of functions, both specified in this course's shared contract:

```go
type HDWallet struct {
	Seed      []byte
	MasterKey []byte
}

func NewHDWallet(seed []byte) *HDWallet
func (hd *HDWallet) DeriveWallet(index uint32) (*Wallet, error)

func EncryptWalletFile(path string, w *Wallet, password string) error
func DecryptWalletFile(path string, password string) (*Wallet, error)
```

`NewHDWallet` takes the seed from Chapter 39 and computes the BIP-32 master key. `DeriveWallet` walks a fixed BIP-44 path down to one specific address, returning a plain `*Wallet` you already know how to use from Chapter 36. `EncryptWalletFile` and `DecryptWalletFile` then let you persist any `*Wallet` — HD-derived or not — to disk, protected by a password instead of sitting there in the clear.

---

## 2. From Seed to Master Key

BIP-32's master key derivation is itself refreshingly simple: run the seed through HMAC-SHA512 using a fixed, publicly known key string, and split the 64-byte output in half.

```
seed (64 bytes)
      |
      |  HMAC-SHA512(key = "GoChain seed", message = seed)
      v
64-byte output
+------------------------+------------------------+
|  bytes 0-31             |  bytes 32-63            |
|  master private key     |  master chain code      |
+------------------------+------------------------+
```

The left 32 bytes become the master **private key**; the right 32 bytes become the master **chain code** — the extra ingredient from Chapter 38, Section 4, that keeps child derivation unpredictable to anyone who only sees public keys. Real BIP-32 uses the literal key string `"Bitcoin seed"` for this step, regardless of which coin is being derived — it's simply a fixed domain-separation constant. GoChain uses its own string, `"GoChain seed"`, so that a GoChain seed and, say, a Bitcoin seed sharing the same 64 random bytes would never accidentally derive the same master key.

```go
package wallet

import (
	"crypto/hmac"
	"crypto/sha512"
)

// hdMasterKeyDomain is HMAC-SHA512's key input for master key
// derivation. Changing this string changes every wallet's derived
// addresses, so it must never vary between GoChain versions.
const hdMasterKeyDomain = "GoChain seed"

// NewHDWallet derives a 64-byte master key (32-byte private key
// followed by a 32-byte chain code) from a seed, storing both the
// seed and the derived master key on the returned HDWallet.
func NewHDWallet(seed []byte) *HDWallet {
	mac := hmac.New(sha512.New, []byte(hdMasterKeyDomain))
	mac.Write(seed)
	masterKey := mac.Sum(nil) // 64 bytes: 32-byte key || 32-byte chain code

	return &HDWallet{
		Seed:      seed,
		MasterKey: masterKey,
	}
}
```

`NewHDWallet` is short because BIP-32's master-key step *is* short: construct an HMAC keyed with the fixed domain string, feed it the seed, and keep the 64-byte result as-is. We store it whole, as a single `MasterKey []byte` field (matching this course's shared type contract), splitting it into its two 32-byte halves only where derivation actually needs each half separately — which is the next section.

---

## 3. The Child Derivation Step

BIP-32's child derivation step, introduced conceptually in Chapter 38 Section 4, takes a parent key, a parent chain code, and an index, and produces a child key and a new chain code. Here is the actual formula, in the **hardened** flavor (used throughout GoChain's fixed path, for the reasons Chapter 38 Section 5 explained):

```
data = 0x00 || parentPrivateKey (32 bytes) || index (4 bytes, big-endian)
       (the leading 0x00 pads the private key to line up with how a
        public key would be represented, per the BIP-32 spec)

I = HMAC-SHA512(key = parentChainCode, message = data)

IL = I[0:32]     -- combined with the parent key to form the child key
IR = I[32:64]    -- becomes the child chain code

childKey = (IL + parentPrivateKey) mod n
           (n is the order of the elliptic curve GoChain's keys use)
```

```
 parent key + parent chain code + index
                  |
                  |  HMAC-SHA512(chainCode, 0x00||parentKey||index)
                  v
        I (64 bytes) = IL || IR
                  |
        childKey = (IL + parentKey) mod n         childChainCode = IR
```

The `mod n` step is the one piece of real elliptic-curve math involved: `n` is the order of the curve (the size of the group the private keys live in), and adding two 32-byte numbers together can overflow past it, so the result wraps back around using modular arithmetic — exactly the same idea as clock arithmetic wrapping past 12. This guarantees the child key is always still a valid private key on the same curve.

```go
package wallet

import (
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/hmac"
	"crypto/sha512"
	"encoding/binary"
	"errors"
	"math/big"
)

// hardenedOffset marks an index as hardened, per BIP-32: any index at
// or above 2^31 is derived using the parent's private key rather than
// its public key (Chapter 38, Section 5).
const hardenedOffset = uint32(0x80000000)

// deriveChild implements one hardened BIP-32 derivation step: given a
// 32-byte parent private key, a 32-byte parent chain code, and an
// index, it returns the 32-byte child private key and 32-byte child
// chain code.
func deriveChild(parentKey, parentChainCode []byte, index uint32) (childKey, childChainCode []byte, err error) {
	if len(parentKey) != 32 || len(parentChainCode) != 32 {
		return nil, nil, errors.New("wallet: parent key and chain code must each be 32 bytes")
	}

	// Force this index into hardened range. GoChain's fixed path
	// (Section 4) always derives hardened, so every caller of this
	// unexported helper already expects that.
	hardenedIndex := index | hardenedOffset

	data := make([]byte, 0, 1+32+4)
	data = append(data, 0x00)         // padding byte, per BIP-32
	data = append(data, parentKey...) // the parent PRIVATE key, for hardened derivation
	indexBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(indexBytes, hardenedIndex)
	data = append(data, indexBytes...)

	mac := hmac.New(sha512.New, parentChainCode)
	mac.Write(data)
	i := mac.Sum(nil) // 64 bytes

	il, ir := i[:32], i[32:]

	curve := elliptic.P256()
	n := curve.Params().N

	// childKey = (IL + parentKey) mod n
	ilInt := new(big.Int).SetBytes(il)
	parentInt := new(big.Int).SetBytes(parentKey)
	childInt := new(big.Int).Add(ilInt, parentInt)
	childInt.Mod(childInt, n)

	if childInt.Sign() == 0 {
		// Astronomically unlikely, but per spec: a zero child key is
		// invalid and derivation should be retried with the next index.
		return nil, nil, errors.New("wallet: derived a zero child key, retry with a different index")
	}

	return childInt.Bytes(), ir, nil
}

// keyPairFromScalar reconstructs a full crypto.KeyPair-shaped ecdsa
// key pair from a raw 32-byte private scalar, by computing the public
// key as scalar * basePoint on the curve.
func keyPairFromScalar(scalar []byte) *ecdsa.PrivateKey {
	curve := elliptic.P256()
	d := new(big.Int).SetBytes(scalar)
	x, y := curve.ScalarBaseMult(d.Bytes())

	return &ecdsa.PrivateKey{
		PublicKey: ecdsa.PublicKey{Curve: curve, X: x, Y: y},
		D:         d,
	}
}
```

`deriveChild` is the direct Go translation of the formula above: it builds the `0x00 || parentKey || index` byte string, runs it through HMAC-SHA512 keyed with the parent chain code, splits the result into `IL`/`IR`, and computes the child key as `(IL + parentKey) mod n` using `math/big` for the modular arithmetic (Go's standard library has no built-in "mod-add" for byte slices, so we lift both operands into `big.Int`, add, and reduce modulo the curve's order `N`). `keyPairFromScalar` is a small but important helper: given only a raw 32-byte private scalar (which is all `deriveChild` produces), it reconstructs a complete `*ecdsa.PrivateKey` — public key included — by computing `scalar * basePoint` on the curve, exactly how any ECDSA private key's public half is derived from its private half.

We deliberately implement only the hardened derivation path here, since Section 4's fixed GoChain derivation path uses hardened derivation at every level it needs — a full BIP-32 implementation supporting normal (public-key-based) derivation too would add real value for advanced use cases like watch-only wallets, but is out of scope for the fixed, single-purpose path this volume builds.

---

## 4. GoChain's Fixed BIP-44 Path

Rather than exposing every level of BIP-44's path as a parameter, GoChain's `DeriveWallet(index uint32)` fixes everything except the final `address_index`, which is exactly what real wallet apps do for their default "just give me my next address" flow:

```
m / 44' / 9999' / 0' / 0 / index

44'    -- purpose: BIP-44 (fixed)
9999'  -- coin_type: GoChain's placeholder coin type (fixed)
0'     -- account: GoChain's single default account for this course (fixed)
0      -- change: 0 = receiving address (fixed; GoChain doesn't yet
          distinguish change addresses, unlike Bitcoin's UTXO change
          outputs, which Volume 5 handled differently — see note below)
index  -- address_index: the one thing DeriveWallet lets you vary
```

A quick, important caveat: real coin types are registered in [SLIP-44](https://github.com/satoshilabs/slips/blob/master/slip-0044.md) (Bitcoin is `0'`, Ethereum is `60'`) to avoid two different coins' wallets accidentally deriving colliding addresses from the same seed. `9999'` is not a registered coin type — it's a clearly-fake placeholder appropriate for a learning project. A production fork of GoChain intending real interoperability would need to register (or otherwise agree on) a real coin type.

```go
package wallet

const (
	bip44Purpose  = 44 | 0 // written as-is for clarity; combined with hardenedOffset below
	bip44CoinType = 9999   // placeholder; not a registered SLIP-44 coin type
	bip44Account  = 0
	bip44Change   = 0
)
```

---

## 5. Implementing DeriveWallet()

```go
package wallet

import (
	"crypto/elliptic"
	"errors"

	"github.com/you/gochain/crypto"
)

// DeriveWallet walks the fixed path m/44'/9999'/0'/0/index (Section 4)
// down from this HDWallet's master key, returning a plain *Wallet for
// the resulting address. Calling DeriveWallet(0), DeriveWallet(1), ...
// repeatedly produces the same sequence of addresses every time, for
// the same underlying seed.
func (hd *HDWallet) DeriveWallet(index uint32) (*Wallet, error) {
	if len(hd.MasterKey) != 64 {
		return nil, errors.New("wallet: HDWallet has no valid master key; call NewHDWallet first")
	}

	key := hd.MasterKey[:32]
	chainCode := hd.MasterKey[32:]

	// Walk each fixed level of the path in turn. Every level here is
	// hardened, including the final address_index, which keeps the
	// entire derivation safe even though DeriveWallet only exposes one
	// varying parameter to callers.
	path := []uint32{bip44Purpose, bip44CoinType, bip44Account, bip44Change, index}

	var err error
	for _, levelIndex := range path {
		key, chainCode, err = deriveChild(key, chainCode, levelIndex)
		if err != nil {
			return nil, err
		}
	}

	privKey := keyPairFromScalar(key)
	kp := &crypto.KeyPair{
		PrivateKey: privKey,
		PublicKey:  elliptic.Marshal(privKey.Curve, privKey.X, privKey.Y),
	}

	return &Wallet{KeyPair: kp}, nil
}
```

`DeriveWallet` walks the five-level path from Section 4 one level at a time, feeding each level's index through `deriveChild` and threading the resulting key and chain code into the next level. After the loop, `key` holds the final 32-byte private scalar for this specific address; `keyPairFromScalar` turns that into a full `*ecdsa.PrivateKey`, and we wrap it in a `crypto.KeyPair` (matching Volume 2's exact type) and then a `*Wallet` — the same type `wallet.New()` produces. Everything built on top of `*Wallet` since Chapter 36 — `Address()`, signing transactions, the CLI wallet — works unmodified on a `DeriveWallet`-produced wallet, exactly as promised at the end of Chapter 38.

---

## 6. Trying It: One Seed, Many Addresses

```go
package main

import (
	"fmt"
	"log"

	"github.com/you/gochain/wallet"
)

func main() {
	mnemonic, err := wallet.NewMnemonic()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Seed phrase (write this down!):", mnemonic)

	seed := wallet.SeedFromMnemonic(mnemonic, "")
	hd := wallet.NewHDWallet(seed)

	// Derive five addresses from the same seed, one at a time.
	for i := uint32(0); i < 5; i++ {
		w, err := hd.DeriveWallet(i)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("address %d: %s\n", i, w.Address())
	}
}
```

Running this program twice with the *same* mnemonic (say, by hardcoding one you saved from a first run and calling `wallet.SeedFromMnemonic` directly, skipping `NewMnemonic`) reproduces the identical five addresses both times — the whole HD wallet promise, now working end to end. Sample output looks like:

```
Seed phrase (write this down!): witch collapse practice feed shame open maple north rescue fee coyote crumble
address 0: 1Gz8kR9nQm...
address 1: 1PxT4vLwZs...
address 2: 1Kj2mNp7Ax...
address 3: 1FbW9sQdRt...
address 4: 1Ym3cVhU8e...
```

---

## 7. The Problem With Plaintext Wallet Files

Chapter 36's simple CLI wallet almost certainly saved its key pair to disk in some form so it would survive between runs — and if that file sits there unencrypted, anyone who gains access to the machine (a stolen laptop, a compromised backup, a misconfigured cloud storage bucket) can read the private key directly and drain every address it controls, instantly and irreversibly. There is no "revoke" button in a blockchain the way there is for a leaked credit card number.

The fix is not exotic: encrypt the wallet file with a key derived from a password the user chooses, so that even someone who obtains the raw file cannot do anything with it without also knowing (or guessing) that password.

---

## 8. Password Stretching With scrypt

A tempting-but-wrong approach is to hash the user's password once (say, with SHA-256) and use that hash directly as an AES key. The problem: SHA-256 is *fast* — modern hardware can compute billions of SHA-256 hashes per second, which means an attacker who steals an encrypted wallet file can brute-force plausible passwords at that same billions-per-second rate.

**scrypt** is a password-based key derivation function designed specifically to resist this kind of attack, by being deliberately expensive in both *time* and *memory* to compute. That second property — memory-hardness — matters because it makes scrypt resistant to attacks using specialized, cheap-to-mass-produce hardware (like the GPUs and ASICs used for mining in Volume 4): those devices have plenty of raw compute but comparatively little fast memory per unit, so a memory-hungry function narrows the attacker's speed advantage far more than a purely CPU-bound one would.

```go
package wallet

import "golang.org/x/crypto/scrypt"

// scrypt cost parameters. N controls CPU/memory cost (must be a power
// of two); r and p tune memory usage and parallelism. These specific
// values are a widely used, reasonable default for interactive
// (not high-throughput server) password hashing as of this writing.
const (
	scryptN      = 1 << 15 // 32768
	scryptR      = 8
	scryptP      = 1
	scryptKeyLen = 32 // AES-256 needs a 32-byte key
)

// deriveKeyFromPassword stretches a user-supplied password (plus a
// random salt unique to this wallet file) into a 32-byte AES key,
// using scrypt to make brute-forcing the password computationally
// expensive even if an attacker has the encrypted file.
func deriveKeyFromPassword(password string, salt []byte) ([]byte, error) {
	return scrypt.Key([]byte(password), salt, scryptN, scryptR, scryptP, scryptKeyLen)
}
```

`deriveKeyFromPassword` wraps `scrypt.Key` with fixed cost parameters: `scryptN` (32768) is the main CPU/memory cost knob, doubling it roughly doubles both the time and memory required per attempt; `scryptR` and `scryptP` are secondary tuning parameters for memory block size and parallelism that we leave at commonly recommended defaults. Crucially, the function takes a `salt` — Section 9 shows exactly how that salt is generated and stored — so that two users with the same password still end up with completely different derived keys, and so precomputed "rainbow table" style attacks against common passwords don't work.

---

## 9. Implementing EncryptWalletFile() and DecryptWalletFile()

With a password-derived AES key in hand, encrypting the wallet file itself is standard, well-trodden ground: AES-256 in GCM mode, which gives both confidentiality (nobody can read the plaintext without the key) and integrity (any tampering with the ciphertext is detected on decryption, rather than silently producing corrupted garbage).

```go
package wallet

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/json"
	"errors"
	"os"

	"github.com/you/gochain/crypto"
)

// walletFile is the on-disk (post-decryption) representation of a
// wallet: just enough to reconstruct a crypto.KeyPair. We store the
// raw private scalar and public key bytes rather than any Go-specific
// binary encoding, keeping the format simple and stable.
type walletFile struct {
	PrivateKeyD []byte `json:"private_key_d"`
	PublicKey   []byte `json:"public_key"`
}

// EncryptWalletFile serializes w, encrypts it with a key derived from
// password via scrypt, and writes salt || nonce || ciphertext to path.
// Anyone who obtains the file without the password sees only random-
// looking bytes.
func EncryptWalletFile(path string, w *Wallet, password string) error {
	plaintext, err := json.Marshal(walletFile{
		PrivateKeyD: w.KeyPair.PrivateKey.D.Bytes(),
		PublicKey:   w.KeyPair.PublicKey,
	})
	if err != nil {
		return err
	}

	// A fresh, random salt per file means two wallets encrypted with
	// the same password still derive different AES keys.
	salt := make([]byte, 16)
	if _, err := rand.Read(salt); err != nil {
		return err
	}

	key, err := deriveKeyFromPassword(password, salt)
	if err != nil {
		return err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return err
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err := rand.Read(nonce); err != nil {
		return err
	}

	// Seal appends the authentication tag to the ciphertext, so a
	// tampered file fails to decrypt rather than decrypting to
	// corrupted-but-unnoticed data.
	ciphertext := gcm.Seal(nil, nonce, plaintext, nil)

	// On-disk layout: [16-byte salt][12-byte nonce][ciphertext+tag]
	out := make([]byte, 0, len(salt)+len(nonce)+len(ciphertext))
	out = append(out, salt...)
	out = append(out, nonce...)
	out = append(out, ciphertext...)

	// 0600: only this file's owner can read or write it, on top of the
	// encryption itself — defense in depth against other local users.
	return os.WriteFile(path, out, 0600)
}

// DecryptWalletFile reads the file written by EncryptWalletFile,
// re-derives the AES key from password and the stored salt, and
// decrypts and reconstructs the original *Wallet. A wrong password
// (or a tampered file) causes decryption to fail with an error,
// rather than silently returning garbage key material.
func DecryptWalletFile(path string, password string) (*Wallet, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	if len(data) < 16+12 {
		return nil, errors.New("wallet: encrypted file is too short to be valid")
	}

	salt := data[:16]
	nonce := data[16:28]
	ciphertext := data[28:]

	key, err := deriveKeyFromPassword(password, salt)
	if err != nil {
		return nil, err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		// GCM's authentication check failed: either the password was
		// wrong (wrong derived key) or the file was tampered with.
		// We deliberately don't distinguish the two in the error
		// message, so an attacker can't use error text to confirm
		// they've guessed a correct password against a tampered file.
		return nil, errors.New("wallet: decryption failed (wrong password or corrupted file)")
	}

	var wf walletFile
	if err := json.Unmarshal(plaintext, &wf); err != nil {
		return nil, err
	}

	privKey := keyPairFromScalar(wf.PrivateKeyD)
	kp := &crypto.KeyPair{
		PrivateKey: privKey,
		PublicKey:  wf.PublicKey,
	}
	return &Wallet{KeyPair: kp}, nil
}
```

`EncryptWalletFile` serializes the wallet's raw key material to a small JSON structure, generates a random salt and nonce, derives an AES key from the password and salt via `deriveKeyFromPassword`, and seals the plaintext with AES-GCM, writing `salt || nonce || ciphertext` to disk with permissions restricted to the file's owner. `DecryptWalletFile` reverses every step: read the file, split out the stored salt and nonce, re-derive the same AES key from the supplied password, and call `gcm.Open` — which fails loudly (not silently) if either the password was wrong or the file's bytes were altered in any way, because GCM's built-in authentication tag check fails first. Note the deliberately vague error message on failure: distinguishing "wrong password" from "corrupted file" in the error text would hand an attacker a free oracle for password-guessing attempts, so both cases return the same message.

---

## 10. Testing the Full Round Trip

```go
package wallet

import (
	"os"
	"testing"
)

func TestEncryptDecryptWalletFileRoundTrip(t *testing.T) {
	w := New()

	dir := t.TempDir()
	path := dir + "/test-wallet.dat"

	if err := EncryptWalletFile(path, w, "correct horse battery staple"); err != nil {
		t.Fatalf("EncryptWalletFile failed: %v", err)
	}

	recovered, err := DecryptWalletFile(path, "correct horse battery staple")
	if err != nil {
		t.Fatalf("DecryptWalletFile failed: %v", err)
	}

	if recovered.Address() != w.Address() {
		t.Fatalf("recovered wallet address %s does not match original %s", recovered.Address(), w.Address())
	}
}

func TestDecryptWalletFileWrongPasswordFails(t *testing.T) {
	w := New()
	dir := t.TempDir()
	path := dir + "/test-wallet.dat"

	if err := EncryptWalletFile(path, w, "correct horse battery staple"); err != nil {
		t.Fatalf("EncryptWalletFile failed: %v", err)
	}

	if _, err := DecryptWalletFile(path, "wrong password"); err == nil {
		t.Fatal("expected DecryptWalletFile to fail with the wrong password, got nil error")
	}
}

func TestDecryptWalletFileDetectsTampering(t *testing.T) {
	w := New()
	dir := t.TempDir()
	path := dir + "/test-wallet.dat"

	if err := EncryptWalletFile(path, w, "correct horse battery staple"); err != nil {
		t.Fatalf("EncryptWalletFile failed: %v", err)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	data[len(data)-1] ^= 0xFF // flip the last byte, simulating tampering
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	if _, err := DecryptWalletFile(path, "correct horse battery staple"); err == nil {
		t.Fatal("expected DecryptWalletFile to detect tampering, got nil error")
	}
}
```

`TestEncryptDecryptWalletFileRoundTrip` checks the core promise: encrypt a wallet, decrypt it back, and confirm the resulting address is identical to the original. `TestDecryptWalletFileWrongPasswordFails` confirms a wrong password is rejected rather than silently producing a broken (or worse, subtly different but valid-looking) wallet. `TestDecryptWalletFileDetectsTampering` flips a single bit in the encrypted file on disk and confirms GCM's authentication check catches it — the same "detect any tampering" property Volume 3's block hashes gave you for the chain itself, now applied to a wallet file at rest.

---

## Summary

- BIP-32's master key comes from HMAC-SHA512 of the seed, split into a 32-byte private key and a 32-byte chain code; GoChain uses the domain string `"GoChain seed"` for this step.
- Hardened child derivation computes `childKey = (IL + parentKey) mod n`, where `I = HMAC-SHA512(chainCode, 0x00||parentKey||index)` and `n` is the elliptic curve's order.
- `DeriveWallet(index)` walks a fixed path, `m/44'/9999'/0'/0/index`, so only the final address index varies per call — `9999'` is a placeholder, not a registered SLIP-44 coin type.
- An unencrypted wallet file is a standing liability: anyone who reads it can drain every address it controls, with no way to revoke a leaked private key.
- **scrypt** stretches a user's password into an AES key slowly and memory-intensively on purpose, raising the cost of brute-force password guessing far above a single fast hash.
- `EncryptWalletFile`/`DecryptWalletFile` use scrypt for key derivation and AES-256-GCM for encryption, storing `salt || nonce || ciphertext` on disk, with GCM's authentication tag catching both wrong passwords and file tampering.
- Wrong-password and tampered-file failures are reported with the same generic error message, deliberately, to avoid handing an attacker a distinguishing oracle.
- Every derived `*Wallet` and every decrypted `*Wallet` is a completely normal `Wallet` — all downstream GoChain code (signing, `Address()`, the CLI) needs no changes to use either kind.

---

## Exercises

### Easy

1. **Run the Section 6 example program** twice, saving the printed mnemonic from the first run and hardcoding it (skipping `NewMnemonic()`) on the second run. Confirm all five printed addresses are identical between runs, and explain in your own words which specific function call is responsible for that determinism.

2. **Encrypt a wallet to a file, then open the file in a text editor** (or with `xxd`/`hexdump`). Confirm you cannot identify any recognizable structure or key material by eye, and explain what specifically in `EncryptWalletFile` is responsible for that (as opposed to, say, just base64-encoding the plaintext).

3. **Change `scryptN` to a much smaller value** (like `1 << 8`) in a local copy of the code and benchmark how much faster `deriveKeyFromPassword` runs. Explain, using Section 8's reasoning, why this "speedup" is actually a security downgrade.

---

### Medium

4. **Implement `TestDeriveWalletIsDeterministicAcrossCalls`**: derive `hd.DeriveWallet(3)` twice from the same `HDWallet` and assert the two resulting wallets' addresses are identical, then derive `hd.DeriveWallet(3)` and `hd.DeriveWallet(4)` and assert those two addresses are *different*.

5. **Add a `ChangePassword(path, oldPassword, newPassword string) error` function** that decrypts a wallet file with the old password and re-encrypts it (with a fresh salt and nonce) using the new password, without ever writing the wallet's private key to disk in plaintext at any intermediate step.

6. **Investigate what happens if two different HDWallets** (derived from two completely different, unrelated seeds) both call `DeriveWallet(0)`. Explain, using Section 3's derivation formula, why their resulting addresses are unrelated and unpredictable from one another, even though they used the identical fixed path.

---

### Hard

7. **Implement normal (non-hardened) derivation** as a second unexported helper, `deriveChildNormal`, using the parent's *public* key in the HMAC input instead of the private key (per BIP-32's normal-derivation formula), and write a test demonstrating that deriving a normal child's public key from only the parent's public key (no private key involved) produces the same public key as deriving the full child key pair and then extracting its public half.

8. **Design and implement an "extended public key" export** — a function that returns just the chain code and public key (never the private key) for a given HDWallet path level — and demonstrate, using your Exercise 7 implementation, that this exported data alone is sufficient to derive an entire subtree of receiving addresses without ever exposing any private key.

9. **Research real-world scrypt parameter recommendations** (as of your research date) for interactive password-based encryption versus server-side password hashing, and argue whether GoChain's chosen `scryptN = 32768` value in Section 8 is appropriately tuned for a wallet file a user decrypts occasionally on their own laptop, or whether it should be higher or lower, citing the specific trade-off (decryption latency vs. brute-force resistance) that drives your answer.
