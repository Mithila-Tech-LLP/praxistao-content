---
title: Cryptography Foundations
---
Every guarantee a blockchain makes — that a block hasn't been tampered with, that a transaction really came from its claimed sender — comes down to cryptography. This is the math and code underneath everything else in this roadmap.

### Hashing & SHA-256
A hash function turns any amount of data into a fixed-size fingerprint. Change one byte of the input and the fingerprint changes completely — this one property is what makes a blockchain tamper-evident.

**Resources:**
- [Hashing and SHA-256](course:blockchain-in-go#08-hashing-and-sha256)
- [Implementing a Hasher in Go](course:blockchain-in-go#09-implementing-a-hasher-in-go)

### Merkle Trees
A Merkle tree lets you prove a single transaction is included in a block without downloading every transaction in it — the trick behind "light clients" in real blockchains.

**Resources:**
- [Merkle Trees Explained and Implemented](course:blockchain-in-go#10-merkle-trees-explained-and-implemented)

### Public-Key Cryptography & Digital Signatures
Wallets don't have passwords — they have key pairs. Understand how a private key signs a transaction and how anyone can verify that signature using only the public key.

**Resources:**
- [Public-Key Cryptography Basics](course:blockchain-in-go#11-public-key-cryptography-basics)
- [Digital Signatures with ECDSA](course:blockchain-in-go#12-digital-signatures-with-ecdsa)
- [Implementing Keys and Signatures in Go](course:blockchain-in-go#13-implementing-keys-and-signatures-in-go)

### Addresses & Encoding
Public keys are long and ugly. Addresses are a shorter, checksummed, human-shareable encoding of them — with a checksum specifically so a typo gets caught instead of silently sending funds into the void.

**Resources:**
- [Encoding: Hex, Base58, and Addresses with Checksums](course:blockchain-in-go#14-encoding-hex-base58-and-addresses-with-checksums)

### Practice: Build a File Integrity Verifier
> branches-from: Hashing & SHA-256

Apply hashing to a real, self-contained tool: detect whether a file has been modified since you last checked it.

**Resources:**
- [Mini Project: A File Integrity Verifier](course:blockchain-in-go#15-mini-project-a-file-integrity-verifier)
