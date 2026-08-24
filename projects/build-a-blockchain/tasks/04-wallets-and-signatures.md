# Task 04: Wallets and Signatures

## What you will build

A `Wallet` type backed by a real ECDSA key pair, plus `Sign` and `Verify` functions — the exact mechanism that lets someone prove they authorized a piece of data without ever revealing their private key.

## Concepts

### Two keys, one relationship

A wallet is a private/public key pair. The private key signs; the public key verifies. Anyone holding only the public key can confirm a signature is genuine, but cannot forge a new one — that asymmetry is the entire point.

```
   Alice's private key  --sign(message)-->  signature
   Alice's public key + message + signature --verify()--> true / false
```

### Sign the hash, not the raw data

Real signing schemes sign a fixed-size hash of the data, not the raw data itself — that way the same signing function works whether the message is 10 bytes or 10 megabytes, and it composes cleanly with the hashing you already built in Task 01.

### Never store or print a private key carelessly

Once generated, a private key should only ever be used to sign — never logged, never sent anywhere. Every real key-management disaster in cryptocurrency history traces back to a private key ending up somewhere it shouldn't.

## Interface to implement

```go
type Wallet struct {
	PrivateKey *ecdsa.PrivateKey
	PublicKey  []byte // uncompressed X||Y bytes
}

// NewWallet generates a fresh ECDSA key pair on the P-256 curve.
func NewWallet() *Wallet

// Address derives a short, printable identifier from the wallet's public
// key -- hash the public key and hex-encode the result. (Real chains
// Base58-encode with a checksum; a hex string is enough for this task.)
func (w *Wallet) Address() string

// Sign signs the SHA-256 hash of data with priv, returning the signature.
func Sign(priv *ecdsa.PrivateKey, data []byte) []byte

// Verify reports whether signature is a valid signature over data's hash,
// produced by the private key matching pubKey.
func Verify(pubKey []byte, data, signature []byte) bool
```

## Hints

- Use `crypto/ecdsa`, `crypto/elliptic` (P-256, `elliptic.P256()`), and `crypto/rand` for key generation — never `math/rand` for anything cryptographic.
- Represent a public key as `X` and `Y`, each padded to 32 bytes with `big.Int.FillBytes`, concatenated — 64 bytes total. You'll need the reverse (bytes back into `X`/`Y`) inside `Verify`.
- `ecdsa.SignASN1` and `ecdsa.VerifyASN1` handle the actual signing math — you're wiring them up, not reimplementing ECDSA.
- Write a test that signs a message, verifies it succeeds, then flips one byte of the message and confirms verification now fails.
- Write a second test that verifies Alice's signature against Mallory's public key and confirms it fails, even though both keys and the signature are individually "valid."

## Run the tests

```bash
cd starter/task-04-wallets-and-signatures
go test ./...
```

All tests must pass before moving to Task 05.
