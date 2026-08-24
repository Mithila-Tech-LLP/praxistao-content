---
title: Blocks & Consensus
---
With the cryptographic building blocks in place, assemble them into an actual block, chain blocks together, and solve the hardest problem in the whole system: how do independent nodes agree on which chain is the real one?

### Anatomy of a Block
A block is a struct: a header (previous hash, timestamp, nonce) plus a body (transactions). Simple in shape, but every field exists for a specific reason.

**Resources:**
- [Anatomy of a Block](course:blockchain-in-go#16-anatomy-of-a-block)
- [Building the Block Struct in Go](course:blockchain-in-go#17-building-the-block-struct-in-go)

### Linking Blocks & Validation
Each block references the hash of the one before it — that link is what makes tampering with an old block detectable, because it would invalidate every block after it.

**Resources:**
- [Linking Blocks and the Genesis Block](course:blockchain-in-go#18-linking-blocks-and-the-genesis-block)
- [Block Validation and Immutability](course:blockchain-in-go#19-block-validation-and-immutability)

### Practice: Tamper-Evident Log
> branches-from: Linking Blocks & Validation

Build a minimal chain of linked, hashed records and prove to yourself that modifying an old entry breaks the chain.

**Resources:**
- [Mini Project: A Tamper-Evident Log](course:blockchain-in-go#22-mini-project-a-tamper-evident-log)

### Consensus & Proof of Work
With no central authority, how do thousands of independent nodes agree on the same history? Proof of Work is Bitcoin's answer: make adding a block expensive enough that lying about the chain isn't worth it.

**Resources:**
- [What Is Consensus and Why Does It Matter](course:blockchain-in-go#23-what-is-consensus-and-why-does-it-matter)
- [Proof of Work Explained](course:blockchain-in-go#24-proof-of-work-explained)
- [Implementing Proof of Work in Go](course:blockchain-in-go#25-implementing-proof-of-work-in-go)

### Difficulty Adjustment & Concurrent Mining
> optional

As more miners join a network, blocks would be found faster and faster without a correction — difficulty adjustment is that correction. Concurrent mining is how you actually use multiple CPU cores to search for a valid nonce.

**Resources:**
- [Difficulty Adjustment](course:blockchain-in-go#26-difficulty-adjustment)
- [Concurrent Mining with Goroutines](course:blockchain-in-go#27-concurrent-mining-with-goroutines)

### Practice: Hash a Block & Mine It
> branches-from: Consensus & Proof of Work

Build the block and proof-of-work tasks of the standalone blockchain project — hash a block, chain it to the previous one, and mine it under a difficulty target.

**Resources:**
- [Build a Blockchain project](project:build-a-blockchain)
