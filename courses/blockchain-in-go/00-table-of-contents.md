# Blockchain in Go

### From Zero to Building, Networking, and Deploying Your Own Live Blockchain in Go

---

## What You Will Build

By the end of this course you will have built **GoChain** — a real, working blockchain, written entirely by you in Go. Not a toy simulation that runs on one computer and pretends to be a network. A real system: it hashes and mines blocks, verifies signed transactions, talks to other nodes over the network, stores its data on disk, runs simple smart contracts, and can be deployed to the cloud as a live, multi-node testnet that other people can connect to.

Here is what using GoChain looks like once you are done:

```go
package main

import (
	"fmt"
	"log"

	"github.com/you/gochain/core"
	"github.com/you/gochain/wallet"
)

func main() {
	// Create a wallet. This generates a private/public key pair and an
	// address — think of the address as an account number anyone can send
	// money to, and the private key as the only thing that can spend it.
	w := wallet.New()
	fmt.Println("Your address:", w.Address())

	// Open (or create) the local blockchain, stored on disk in ./data.
	chain, err := core.OpenBlockchain("./data/chain.dat")
	if err != nil {
		log.Fatal(err)
	}

	// Send 10 coins to another address. Under the hood this builds a
	// transaction, signs it with your private key, and adds it to the
	// mempool — the waiting room for transactions that haven't been mined
	// into a block yet.
	tx, err := chain.Send(w, "1Bv11k7c...recipient", 10)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Transaction submitted:", tx.IDHex())

	// Mine a new block. This picks up pending transactions from the
	// mempool, solves the proof-of-work puzzle, and appends the new
	// block to the chain — permanently, as long as no one out-competes it.
	block := chain.MineBlock(nil)
	fmt.Printf("Mined block %d with hash %x\n", block.Height, block.Hash)

	// Check the balance the same way a real wallet app would.
	balance := chain.BalanceOf(w.Address())
	fmt.Println("New balance:", balance)
}
```

A note on reading this preview: every call above is either exactly what a later chapter builds (`wallet.New`, `OpenBlockchain`, `MineBlock`, `BalanceOf` all appear with these exact signatures starting in Volumes 3–5) or a thin, obvious convenience a real caller would write in a few lines on top of primitives this course builds explicitly and separately (`Send` composes `NewTransaction` + `Sign` + `Mempool.Add`, exactly as Chapter 36's CLI wallet does; `IDHex` is a one-line `hex.EncodeToString(tx.ID)`). Nothing here is magic, and nothing here is skipped — it is simply the destination, shown once before the two-hundred-some pages of chapters that actually build it.

By Volume 7 this node is talking to other GoChain nodes over the network, exchanging blocks and transactions automatically. By Volume 9 it can run small programs — smart contracts — that move coins based on code, not just simple transfers. By the final volume, you deploy a small network of GoChain nodes to real cloud servers, wire up a block explorer website, and open a public testnet that friends can join from their own laptops.

Every line of GoChain — the cryptography, the block format, the mining loop, the networking protocol, the storage engine, the virtual machine, the API, the deployment scripts — is code you will have written yourself, understood completely, and be able to explain to anyone.

---

## Who Is This Course For

This course assumes you are a complete beginner to blockchain. You do not need to have read the Bitcoin whitepaper. You do not need a math or cryptography background. You do not need prior Go experience either — Volume 1 teaches you everything about Go you need before a single line of blockchain code appears, and every Go concept used later is explained again, in plain language, the first time it shows up.

This course is for you if:

- You keep hearing about "blockchain," "crypto," "smart contracts," and "Web3" and want to actually understand what is really going on underneath the buzzwords.
- You know some programming but have never touched cryptography, peer-to-peer networking, or consensus algorithms.
- You have used Go for small projects and want a large, real, satisfying system to build with it.
- You want to work in blockchain/Web3 engineering and need to prove — to yourself and to employers — that you can build one from scratch, not just call an SDK.
- You want to *stop being confused* about how Bitcoin and Ethereum actually work, by building a smaller version of the same ideas with your own hands.

Every new term — hash, nonce, Merkle root, UTXO, mempool, consensus, gas, gossip protocol — is defined in plain words, with an everyday analogy, the moment it first appears. Nothing is assumed. Nothing is hand-waved.

---

## How Long Will This Take

```
PART 1 — FOUNDATIONS (Volumes 1-3)
├── Volume 1: Orientation & Go Essentials         ~  8 hours
├── Volume 2: Cryptography Foundations            ~ 10 hours
└── Volume 3: The Blockchain Data Structure       ~  8 hours
                                                  ----------
                                        Total:    ~ 26 hours

PART 2 — CONSENSUS & VALUE TRANSFER (Volumes 4-6)
├── Volume 4: Proof of Work & Mining              ~  8 hours
├── Volume 5: Transactions & the UTXO Model       ~ 12 hours
└── Volume 6: Wallets & Key Management            ~  6 hours
                                                  ----------
                                        Total:    ~ 26 hours

PART 3 — NETWORKING & DECENTRALIZATION (Volume 7)
└── Volume 7: Peer-to-Peer Networking             ~ 14 hours

PART 4 — STORAGE, CONTRACTS & TOOLING (Volumes 8-10)
├── Volume 8:  Storage & State                    ~  8 hours
├── Volume 9:  Smart Contracts & a Virtual Machine ~ 16 hours
└── Volume 10: APIs, Explorer & Developer Tools    ~  8 hours
                                                  ----------
                                        Total:    ~ 32 hours

PART 5 — SECURITY & THE WIDER WORLD (Volumes 11-12)
├── Volume 11: Security, Alt Consensus & Scaling  ~ 10 hours
└── Volume 12: Real-World Case Studies             ~  6 hours
                                                  ----------
                                        Total:    ~ 16 hours

PART 6 — SHIPPING IT (Volume 13)
└── Volume 13: Deployment & Going Live            ~ 16 hours

GRAND TOTAL:                                      ~130 hours
```

At an hour a day, that is roughly four and a half months. At two hours a day, about two months. Most people find the second half — networking, contracts, and deployment — the most addictive part, because that is where GoChain stops being a private exercise and starts being a real, reachable system.

---

## The Learning Progression

```
Week  1-2   |  Part 1 starts: what a blockchain actually is, stripped of hype.
             |  You set up Go, learn what you need of it, and write your first
             |  hash functions and Merkle trees.

Week  3-4   |  You define the Block and Blockchain types and link blocks
             |  together into a real, tamper-evident chain, saved to disk.

Week  5-6   |  Proof of work: you make your chain expensive to rewrite by
             |  requiring real computational work to add a block, and you
             |  build a genuine multi-threaded miner.

Week  7-9   |  Transactions: the UTXO model, signing, verifying, the mempool,
             |  double-spend prevention, and a working CLI wallet that can
             |  send and receive coins.

Week 10-11  |  Wallets get serious: HD wallets, seed phrases, encrypted key
             |  storage — the same techniques real wallets like MetaMask use.

Week 12-14  |  Networking: your node talks to other nodes. Peer discovery,
             |  gossiping transactions and blocks, resolving forks, running a
             |  real multi-node network on your own machine.

Week 15-16  |  Storage: you replace flat files with a real embedded database,
             |  add fast balance lookups, and understand state tries.

Week 17-19  |  Smart contracts: you build a small virtual machine, a scripting
             |  language, gas metering, and deploy a working token contract.

Week 20-21  |  APIs and tooling: a JSON-RPC/REST API, a block explorer, and a
             |  polished CLI — the tools that make a blockchain usable.

Week 22-23  |  Security and the wider world: real attacks, alternative
             |  consensus (proof of stake, BFT), scaling, and how Bitcoin and
             |  Ethereum really work under the hood.

Week 24-26  |  You deploy. Docker, multi-node testnets, cloud servers,
             |  monitoring, CI/CD, and a public testnet with a faucet that
             |  anyone can connect to. GoChain goes live.
```

---

## Full Table of Contents

---

# PART 1 — FOUNDATIONS

---

## Volume 1: Orientation & Go Essentials for Blockchain

> Before you write a single line of blockchain code, you need two things: a clear, hype-free picture of what a blockchain actually is, and just enough Go to build one. This volume gives you both.

### [Chapter 01: What Is a Blockchain, Really?](01-what-is-a-blockchain-really.md)

Strip away the hype. A blockchain is a list of records (blocks) where each record contains a fingerprint (hash) of the one before it — so changing any old record breaks every fingerprint after it, making tampering obvious. We build this intuition with a paper-and-pencil exercise before any code: a shared notebook, a chain of wax-sealed envelopes, a shared ledger three friends keep in sync. Then we map each analogy onto the real vocabulary: block, chain, hash, node, ledger, decentralization.

**Key topics:** what a blockchain is, why it matters, hash chains as tamper-evidence, decentralization vs. a normal database, honest myth-busting (blockchain is not magic, not always necessary, not the same as "crypto prices").

---

### [Chapter 02: Why Go Is a Great Fit for Blockchain](02-why-go-is-a-great-fit-for-blockchain.md)

Real blockchains (Ethereum's `go-ethereum`, Hyperledger Fabric, Cosmos SDK chains) are written in Go, and for good reasons this chapter makes concrete: goroutines make handling thousands of network peers simple, Go's standard library already ships strong cryptography, static typing catches entire classes of bugs before they ever touch real money, and compiled binaries make deployment painless. A short tour of `go-ethereum`'s public source layout to see these ideas in a production system.

**Key topics:** why Go for blockchain, goroutines and concurrency at scale, Go's `crypto` standard library, static typing and safety, real-world Go blockchains.

---

### [Chapter 03: Environment Setup and the GoChain Project Scaffold](03-environment-setup-and-project-scaffold.md)

Installing Go, a code editor, and Git. Creating the `gochain` Go module that will grow across this entire course. Laying out the package structure we will fill in volume by volume: `core` (blockchain data structures), `crypto` (hashing and signatures), `consensus` (proof of work and beyond), `network` (P2P), `storage` (persistence), `wallet` (keys and addresses), `vm` (smart contracts), `api` (RPC/REST), and `cmd/gochain` (the CLI that ties it all together). Your first commit: a `go.mod` and an empty `main.go` that prints a welcome banner.

**Key topics:** installing Go, `go mod init`, project layout, Git basics for this course, the GoChain package map.

---

### [Chapter 04: A Go Crash Course for Blockchain Developers](04-a-go-crash-course-for-blockchain-developers.md)

Everything you need from Go before block one: variables and types, structs (the shape we will use for blocks and transactions), methods, interfaces (how we will swap storage engines and consensus algorithms later without rewriting everything), error handling the Go way, slices and maps, and pointers (why a block holds a *pointer* to the previous block's hash, not a copy). Every concept is demonstrated with a tiny, blockchain-flavored example, not a generic tutorial snippet.

**Key topics:** structs, methods, interfaces, error handling, slices, maps, pointers vs. values.

---

### [Chapter 05: Concurrency — Goroutines and Channels](05-concurrency-goroutines-and-channels.md)

A blockchain node does many things "at once": mining, listening for new peers, receiving transactions, answering API requests. Goroutines are Go's lightweight threads — you can run thousands of them cheaply. Channels are typed pipes goroutines use to hand data to each other safely, avoiding the classic bugs of shared memory. We build a miniature simulation: one goroutine "discovers" fake transactions every second and sends them over a channel to a goroutine that collects them into a mempool — the exact pattern GoChain's real mempool will use in Volume 5.

**Key topics:** goroutines, channels, `select`, the producer-consumer pattern, race conditions and how channels avoid them.

---

### [Chapter 06: Project Layout and Go Modules](06-project-layout-and-go-modules.md)

How Go modules version and share code, and how GoChain's own packages will import each other (`core` imports `crypto`, `network` imports `core`, and so on) without circular dependencies. Internal packages, why we keep `cmd/` separate from library code, and writing a `Makefile` with `build`, `test`, and `run` targets you will use for the rest of the course.

**Key topics:** Go modules, package dependencies, avoiding import cycles, `internal/` packages, project tooling with Make.

---

### [Chapter 07: Binary Encoding and Testing in Go](07-binary-encoding-and-testing-in-go.md)

Blocks and transactions must be turned into bytes to be hashed, signed, sent over the network, and saved to disk. We compare three encodings — `encoding/json` (readable, larger), `encoding/gob` (Go-native, fast), and a hand-rolled binary format with `encoding/binary` (smallest, used by real blockchains) — and pick gob for GoChain's early chapters, upgrading later. We also set up Go's built-in testing framework (`go test`, table-driven tests, `testing.T`) since every GoChain package from here on ships with tests.

**Key topics:** `encoding/json`, `encoding/gob`, `encoding/binary`, choosing a serialization format, Go testing basics, table-driven tests.

---

## Volume 2: Cryptography Foundations

> Cryptography is the load-bearing wall of every blockchain. Hashing makes tampering detectable. Digital signatures make forging transactions impossible without a private key. This volume builds both from first principles and puts real, working implementations into the `gochain/crypto` package.

### [Chapter 08: Hashing and SHA-256](08-hashing-and-sha256.md)

A hash function takes any input — a word, a book, a block of transactions — and produces a fixed-size fingerprint. Change one letter and the fingerprint changes completely (the avalanche effect). Hashes are one-way (you cannot reverse a fingerprint back into the original) and deterministic (the same input always produces the same fingerprint). We walk through SHA-256 conceptually — the same hash function Bitcoin uses — without needing to implement the math ourselves, since Go's standard library already provides it correctly and securely.

**Key topics:** what a hash function is, determinism, the avalanche effect, one-wayness, collision resistance, why SHA-256 specifically.

---

### [Chapter 09: Implementing a Hasher in Go](09-implementing-a-hasher-in-go.md)

Using Go's `crypto/sha256` package to hash strings, structs, and byte slices. Building `gochain/crypto`'s first function: `Hash(data []byte) []byte`. Hashing a struct correctly requires a *canonical* byte representation — if two equal-looking blocks serialize to different bytes, they will hash differently, which breaks everything downstream. We write and test a `Serialize()` method pattern every GoChain type will follow from here on.

**Key topics:** `crypto/sha256` in Go, hex encoding for readability, canonical serialization, writing your first `crypto` package tests.

---

### [Chapter 10: Merkle Trees — Explained and Implemented](10-merkle-trees-explained-and-implemented.md)

Instead of hashing a list of transactions as one long blob, a Merkle tree hashes pairs of transactions together, then hashes pairs of those hashes, repeatedly, up to a single Merkle root. This lets you prove one transaction is included in a block without needing every other transaction (a Merkle proof) — the trick "light clients" like mobile wallets rely on. We build `MerkleTree` and `MerkleRoot()` in Go step by step, with an ASCII diagram of the tree at every stage, then write a Merkle proof verifier.

**Key topics:** Merkle tree construction, Merkle root, Merkle proofs, light clients, implementing a Merkle tree in Go.

---

### [Chapter 11: Public-Key Cryptography Basics](11-public-key-cryptography-basics.md)

A key pair is two mathematically linked numbers: a private key (secret, proves ownership) and a public key (shareable, lets others verify you). Unlike a password, you never send your private key to prove who you are — you use it to produce a signature that anyone can check against your public key without learning the key itself. An everyday analogy: a wax seal only your unique stamp can produce, that anyone can recognize as genuinely yours by sight.

**Key topics:** public vs. private keys, asymmetric cryptography, why you never share a private key, the seal analogy, elliptic curve cryptography at a conceptual level.

---

### [Chapter 12: Digital Signatures with ECDSA](12-digital-signatures-with-ecdsa.md)

ECDSA (Elliptic Curve Digital Signature Algorithm) is what Bitcoin and Ethereum both use to sign transactions. A signature proves two things at once: that the owner of a specific private key approved this exact data, and that the data has not been altered since (because the signature depends on the data's hash). We walk through signing and verifying conceptually with diagrams before touching code, so the Go implementation in the next chapter feels obvious rather than magical.

**Key topics:** ECDSA explained conceptually, the secp256k1 curve, sign vs. verify, why signatures bind to exact data.

---

### [Chapter 13: Implementing Keys and Signatures in Go](13-implementing-keys-and-signatures-in-go.md)

Using Go's `crypto/ecdsa` and `crypto/elliptic` packages to generate a key pair, sign a message, and verify a signature. Building `gochain/crypto`'s `KeyPair`, `Sign()`, and `Verify()` functions — the exact primitives `gochain/wallet` will build on top of starting in Volume 6. A full worked example: Alice signs "pay Bob 5 coins," Bob (or anyone) verifies it using only Alice's public key, and we deliberately corrupt one byte to see verification fail.

**Key topics:** `crypto/ecdsa` in Go, key generation, signing bytes, verifying signatures, handling signature failures correctly.

---

### [Chapter 14: Encoding — Hex, Base58, and Addresses with Checksums](14-encoding-hex-base58-and-addresses-with-checksums.md)

Raw public keys are ugly, long, and error-prone to type. Base58 encoding (used by Bitcoin) removes visually confusing characters (0, O, I, l) so addresses are easier for humans to read and transcribe correctly. A checksum — a few extra bytes derived from a hash of the address — lets software detect a typo before broadcasting a transaction to the wrong address. We build GoChain's address format: hash the public key, add a version byte, append a checksum, and Base58-encode the result.

**Key topics:** hex vs. Base58 vs. Base64, why Base58 for addresses, checksums for typo detection, building GoChain's address format.

---

### [Chapter 15: Mini Project — A File Integrity Verifier](15-mini-project-a-file-integrity-verifier.md)

Before touching blocks, put hashing to work on something concrete and immediately useful: a CLI tool that fingerprints every file in a folder, saves the fingerprints, and later tells you exactly which files changed, were added, or were deleted — the same core idea `git status` and antivirus tools use, built entirely from the `gochain/crypto.Hash` function you wrote in Chapter 09.

**Key topics:** applying hashing practically, recursive file walking in Go, detecting changes via stored hashes.

**Mini Project:** `fileguard` — a command-line integrity checker that snapshots a directory's file hashes and reports drift on subsequent runs.

---

## Volume 3: The Blockchain Data Structure

> This is where GoChain becomes a real blockchain. You define the Block, link blocks into a chain, make tampering detectable, and persist everything to disk.

### [Chapter 16: Anatomy of a Block](16-anatomy-of-a-block.md)

A block bundles together everything that happened in one "round": a list of transactions, a timestamp, a reference to the previous block, and the proof that real work went into creating it. We diagram GoChain's exact block layout field by field, and compare it side-by-side with a simplified view of a real Bitcoin block header, so you can see how closely our design mirrors production systems.

**Key topics:** block header vs. body, the fields every block needs, comparison with a real Bitcoin block, why order and structure matter for hashing.

---

### [Chapter 17: Building the Block Struct in Go](17-building-the-block-struct-in-go.md)

Defining `core.Block` in Go with the fields decided in Chapter 16: `Height`, `Timestamp`, `Transactions`, `PrevBlockHash`, `Hash`, `Nonce`, and `MerkleRoot`. Writing `NewBlock()`, `Serialize()`, and `ComputeHash()` methods. This is the exact struct every later chapter builds on — we get its shape right once, here, and never redefine it.

**Key topics:** the `core.Block` struct, constructor functions in Go, computing a block's own hash from its contents.

---

### [Chapter 18: Linking Blocks and the Genesis Block](18-linking-blocks-and-the-genesis-block.md)

A chain forms the moment each new block stores the *hash* of the block before it. The very first block — the genesis block — has no predecessor, so its `PrevBlockHash` is conventionally all zeroes. We build `core.Blockchain` with an in-memory slice of blocks, a `NewGenesisBlock()` function, and an `AddBlock()` method, then draw the full chain as an ASCII diagram of boxes and arrows to make the linking mechanism completely concrete.

**Key topics:** the genesis block, linking via `PrevBlockHash`, the `core.Blockchain` type, `AddBlock`.

---

### [Chapter 19: Block Validation and Immutability](19-block-validation-and-immutability.md)

Adding a block is not enough — GoChain must refuse invalid ones. We implement `ValidateBlock()`: does the block's stored hash match a freshly recomputed hash of its contents? Does `PrevBlockHash` actually match the previous block's real hash? Then we intentionally tamper with a transaction inside an old block and watch every hash from that point forward stop matching — proving, hands-on, exactly why blockchains are called "tamper-evident" rather than "tamper-proof" (tampering is *detectable*, not impossible without more work — that "more work" part is Volume 4).

**Key topics:** block validation rules, chain-wide validation, hands-on tamper detection, tamper-evident vs. tamper-proof.

---

### [Chapter 20: Persisting Blocks to Disk](20-persisting-blocks-to-disk.md)

An in-memory chain disappears the moment your program exits — useless for anything real. We build GoChain's first storage layer: a simple append-only file that saves each serialized block and can rebuild the in-memory chain on startup by reading it back in order. This intentionally simple version sets up the motivation for Volume 8, where we replace it with a real embedded database.

**Key topics:** append-only file storage, saving and loading a chain from disk, why this simple approach will need to change later.

---

### [Chapter 21: Building a Chain Inspector CLI](21-building-a-chain-inspector-cli.md)

A command-line tool, `gochain inspect`, that prints every block in the chain in a readable format: height, timestamp, transaction count, hash, and previous hash — plus a `--verify` flag that walks the whole chain checking every link. This becomes your everyday debugging tool for the rest of the course, and your first taste of building GoChain's CLI, which grows steadily until Volume 10.

**Key topics:** building CLI tools in Go with `flag`, formatting blockchain data for humans, a reusable debugging workflow.

---

### [Chapter 22: Mini Project — A Tamper-Evident Log](22-mini-project-a-tamper-evident-log.md)

Apply everything from this volume outside the blockchain context: build a tamper-evident audit log for a fictional file-sharing app — every action (upload, delete, rename) becomes a "block" linked by hash to the action before it, so if anyone edits the log's history file directly, the chain of hashes breaks and your `--verify` command catches it immediately.

**Key topics:** applying block-chaining to general audit logging, a real, non-cryptocurrency use of blockchain structure.

**Mini Project:** `auditlog` — a tamper-evident event log built directly on GoChain's `core.Block` linking logic.

---

# PART 2 — CONSENSUS & VALUE TRANSFER

---

## Volume 4: Proof of Work & Mining

> Being tamper-evident is not enough — anyone could still rewrite history and just recompute all the hashes. Proof of work makes rewriting history expensive, by requiring real computational effort to add each block.

### [Chapter 23: What Is Consensus, and Why Does It Matter?](23-what-is-consensus-and-why-does-it-matter.md)

In a single-computer chain, there is no disagreement to resolve. The moment multiple independent nodes each keep their own copy of the chain, you need a rule for what happens when two nodes disagree about what the "real" chain is. Consensus is that rule. We introduce the problem with a simple thought experiment — three friends keeping the same shared notebook, two of whom write different entries at the same time — before naming any algorithm.

**Key topics:** the consensus problem, why a single source of truth is hard in a distributed system, a preview of the algorithms this course covers.

---

### [Chapter 24: Proof of Work, Explained](24-proof-of-work-explained.md)

Proof of work asks: "find a number (a nonce) that, combined with the block's data, produces a hash starting with a certain number of zeroes." There is no shortcut — you must try nonces one by one (or in parallel) until you get lucky. This deliberate difficulty is the "work" in proof of work, and it is what makes it expensive to rewrite old blocks: you would have to redo the work for that block *and every block after it*, faster than the rest of the honest network combined.

**Key topics:** the proof-of-work puzzle, why it must be hard to solve but easy to verify, the "51% attack" idea previewed conceptually.

---

### [Chapter 25: Implementing Proof of Work in Go](25-implementing-proof-of-work-in-go.md)

Building `consensus.ProofOfWork`: a `Run()` method that loops over candidate nonces, hashing the block's data plus the nonce each time, until the resulting hash meets the target (starts with enough zero bits). A `Validate()` method anyone can use to instantly check a solved block without redoing the search. We wire this into `core.Blockchain.MineBlock()` and mine our first real, proof-of-work-secured block.

**Key topics:** the mining loop in Go, target and difficulty as a numeric threshold, `Validate` vs. `Run`, wiring PoW into `MineBlock`.

---

### [Chapter 26: Difficulty Adjustment](26-difficulty-adjustment.md)

A fixed difficulty is wrong for a real network: as more miners join (or leave), block times drift away from your target (say, one block every 10 seconds). We implement a difficulty-adjustment algorithm, modeled on Bitcoin's, that looks at how long the last N blocks actually took and raises or lowers the difficulty to bring future block times back toward the target — with a worked numeric example showing the adjustment math step by step.

**Key topics:** target block time, difficulty adjustment algorithms, why fixed difficulty breaks as network hash power changes.

---

### [Chapter 27: Concurrent Mining with Goroutines](27-concurrent-mining-with-goroutines.md)

A single CPU core searching nonces one at a time is slow. We rewrite the mining loop to split the nonce search space across multiple goroutines running in parallel, using a channel to report the first valid solution found and a `context.Context` to cancel the remaining goroutines the instant a winner is found — a direct, practical payoff of the concurrency skills from Volume 1.

**Key topics:** parallelizing the nonce search, `context.Context` for cancellation, coordinating goroutines with channels, benchmarking the speedup.

---

### [Chapter 28: Mini Project — Multi-Threaded Miner Benchmark](28-mini-project-multi-threaded-miner-benchmark.md)

Build a small benchmarking harness that mines a batch of blocks at several difficulty levels, once single-threaded and once with your Chapter 27 concurrent miner, and prints a table comparing hashes-per-second and total time — turning "concurrency makes mining faster" into a number you measured yourself.

**Key topics:** writing Go benchmarks, comparing single- vs. multi-threaded performance, presenting benchmark results.

**Mini Project:** `minebench` — a CLI that benchmarks GoChain's miner across difficulty levels and thread counts.

---

## Volume 5: Transactions & the UTXO Model

> A blockchain that only tracks blocks is not very useful — it needs to move value between people. This volume builds real transactions, the accounting model behind them, and the mempool that holds pending ones.

### [Chapter 29: What Is a Transaction?](29-what-is-a-transaction.md)

A transaction is a signed instruction to move value from somewhere to somewhere else. We define what a transaction must contain to be trustworthy: proof of where the funds come from, proof the sender is authorized to spend them (a signature), and where the funds should go. An analogy: a transaction is like a signed check, except the "bank" verifying it is every node on the network independently, using math instead of trust.

**Key topics:** what a transaction represents, the signed-check analogy, why every node must be able to verify a transaction independently.

---

### [Chapter 30: The UTXO Model, Explained](30-the-utxo-model-explained.md)

Bitcoin does not track "account balances" directly. Instead, it tracks Unspent Transaction Outputs (UTXOs) — discrete, indivisible chunks of value, like physical bills and coins in a wallet. Spending means consuming one or more UTXOs entirely as inputs and creating new UTXOs as outputs (including "change" back to yourself, just like paying with a $20 bill for a $12 item and getting $8 in change). Your balance is simply the sum of every UTXO that belongs to your address.

**Key topics:** UTXOs as discrete "bills," inputs and outputs, change outputs, computing balance as a UTXO sum.

---

### [Chapter 31: The Account Model, and Choosing One for GoChain](31-the-account-model-and-choosing-one-for-gochain.md)

Ethereum instead uses an account model: a running balance per address, updated in place, like a bank account ledger. We compare both models side by side — UTXO's natural parallelism and privacy properties versus the account model's simplicity for smart contracts — and explain, with reasons, why GoChain adopts the UTXO model for its core ledger (matching Bitcoin) while previewing that our smart-contract volume will introduce account-like state for contracts specifically.

**Key topics:** account model vs. UTXO model, trade-offs of each, why GoChain uses UTXO for its base ledger.

---

### [Chapter 32: Building Transactions in Go](32-building-transactions-in-go.md)

Defining `core.Transaction`, `core.TxInput`, and `core.TxOutput` exactly as specified in this course's shared type contract. Writing `NewTransaction()`, which selects enough UTXOs to cover an amount, creates a change output if needed, and computes the transaction's own ID by hashing its contents (reusing `gochain/crypto.Hash` from Volume 2).

**Key topics:** the `Transaction`/`TxInput`/`TxOutput` types, selecting UTXOs to spend, generating change outputs, transaction IDs.

---

### [Chapter 33: Signing and Verifying Transactions](33-signing-and-verifying-transactions.md)

Every input in a transaction must be signed with the private key that controls the UTXO it spends — otherwise anyone could "spend" your coins by just naming your address as the source. We implement `Transaction.Sign()` and `Transaction.Verify()` using the ECDSA primitives from Volume 2, and walk through exactly what bytes get signed (a "trimmed" copy of the transaction, to avoid signing the signature itself — a subtle, important detail explained with a diagram).

**Key topics:** signing a transaction's inputs, verifying signatures against public keys, why you sign a trimmed copy, avoiding malleability bugs.

---

### [Chapter 34: The Mempool and Preventing Double-Spending](34-the-mempool-and-preventing-double-spending.md)

The mempool ("memory pool") is where valid, signed, but not-yet-mined transactions wait. We implement `core.Mempool` with add, remove, and "get all pending" operations, and the critical validation check: does this transaction try to spend a UTXO that another pending (or already-mined) transaction has already claimed? This is exactly what stops a double-spend — spending the same coin twice — and we build a test that deliberately attempts one and watches it get rejected.

**Key topics:** the mempool data structure, double-spend detection, transaction validation before mining, mempool eviction.

---

### [Chapter 35: Transaction Fees](35-transaction-fees.md)

Miners choose which pending transactions to include in the next block — a fee (the difference between a transaction's total inputs and total outputs) is their incentive to prioritize yours. We add fee calculation to `core.Transaction`, update `MineBlock()` to sort the mempool by fee-per-byte and greedily fill a block, and discuss fee markets: what happens when the mempool has more transactions waiting than fit in a block.

**Key topics:** how fees are calculated from inputs minus outputs, fee-based transaction selection, fee markets under congestion.

---

### [Chapter 36: Building a CLI Wallet](36-building-a-cli-wallet.md)

A `gochain wallet` command-line tool: generate a new key pair and address, check a balance by scanning the UTXO set, and send coins by building, signing, and broadcasting (for now, just submitting locally) a transaction. This is the first time in the course you interact with GoChain the way a real user would — through a wallet, not through internal Go function calls.

**Key topics:** building a CLI wallet, wiring keys + addresses + transactions together, a real end-user workflow.

---

### [Chapter 37: Major Project 1 — Send and Receive Coins](37-major-project-1-send-and-receive-coins.md)

Bring Volumes 2 through 5 together into one working, end-to-end demo: create two wallets, mine a block that rewards the first wallet with new coins (the "coinbase" transaction — a special transaction with no inputs, which we introduce here), send coins from the first wallet to the second, mine that transaction into a block, and verify both wallets' balances update correctly by scanning the chain from scratch.

**Key topics:** the coinbase transaction, full send/receive lifecycle, recomputing balances from the entire chain, your first true end-to-end GoChain demo.

**Major Project:** A two-wallet send-and-receive demo that exercises the entire transaction, mining, and balance pipeline built so far.

---

## Volume 6: Wallets & Key Management

> A wallet is not a place coins are "stored" — it is a keychain. This volume makes GoChain's wallets production-grade: one seed phrase, many addresses, and encryption at rest.

### [Chapter 38: HD Wallets — BIP-32, BIP-39, and BIP-44 Explained](38-hd-wallets-explained.md)

Real wallets do not use one key pair — they derive an entire tree of key pairs from a single master seed, so you only ever need to back up one thing. BIP-39 turns randomness into a human-writable list of words (a seed phrase). BIP-32 defines how to deterministically derive a whole tree of child keys from that seed. BIP-44 defines a standard path structure so different wallets and coins can derive keys the same predictable way. We explain each standard's job with a family-tree diagram before any code.

**Key topics:** hierarchical deterministic (HD) wallets, BIP-39 seed phrases, BIP-32 key derivation, BIP-44 derivation paths.

---

### [Chapter 39: Implementing Mnemonic Seed Phrases](39-implementing-mnemonic-seed-phrases.md)

Generating cryptographically secure randomness, mapping it onto the standard BIP-39 wordlist, and adding a checksum word so a typo in your seed phrase is detected immediately rather than silently losing your funds. Implementing `wallet.NewMnemonic()` and `wallet.SeedFromMnemonic()` in Go, and testing round-trip generation and recovery.

**Key topics:** secure randomness in Go (`crypto/rand`), the BIP-39 wordlist and checksum, mnemonic generation and recovery.

---

### [Chapter 40: Deriving Addresses and Encrypting Wallet Files](40-deriving-addresses-and-encrypting-wallet-files.md)

Deriving multiple child key pairs (and therefore multiple addresses) from one seed, following BIP-44's path convention. Then, since a wallet file sitting unencrypted on disk is a huge liability, encrypting it with a password-derived key (using `scrypt` or `argon2` to slow down password-guessing attacks) before writing it to disk, and decrypting it on load.

**Key topics:** deriving multiple addresses from one seed, password-based key derivation, encrypting wallet files at rest.

---

### [Chapter 41: Hardware Wallet Concepts](41-hardware-wallet-concepts.md)

You will not build physical hardware in this course, but understanding how hardware wallets like Ledger and Trezor work clarifies *why* wallet software is designed the way it is: the private key never leaves the device, transactions are signed on-device after a physical confirmation, and the computer only ever sees the finished signature — never the key. We map this model onto GoChain's own `wallet.Signer` interface, so a future hardware integration would not require rewriting the rest of the system.

**Key topics:** the hardware wallet security model, signing without exposing keys, designing software around a `Signer` interface.

---

### [Chapter 42: Mini Project — HD Wallet CLI](42-mini-project-hd-wallet-cli.md)

A polished `gochain-wallet` CLI: `new` (generate a seed phrase), `recover` (restore from an existing phrase), `addresses` (list derived addresses), `balance`, and `send` — all backed by the encrypted, HD wallet machinery built in this volume.

**Key topics:** wiring HD wallets into a real CLI tool, a realistic wallet UX.

**Mini Project:** `gochain-wallet` — a full HD, password-encrypted command-line wallet for GoChain.

---

# PART 3 — NETWORKING & DECENTRALIZATION

---

## Volume 7: Peer-to-Peer Networking

> Everything so far has run on one computer. A blockchain is not decentralized until independent nodes, run by different people, can find each other, share data, and agree on a single history. This is the volume where GoChain becomes a real network.

### [Chapter 43: Introduction to Networking for Blockchain](43-intro-to-networking-for-blockchain.md)

A quick, practical grounding in IP addresses, ports, and the client-server vs. peer-to-peer distinction — in a P2P network, every node can act as both a client and a server to every other node. We diagram what GoChain's network needs to do: find peers, exchange blocks and transactions, and keep every honest node's chain in sync, before writing any networking code.

**Key topics:** IP addresses and ports, client-server vs. peer-to-peer, what GoChain's network layer must accomplish.

---

### [Chapter 44: TCP Servers and Clients in Go](44-tcp-servers-and-clients-in-go.md)

Building a bare TCP server and client using Go's `net` package: `net.Listen`, `net.Dial`, reading and writing length-prefixed messages so the receiver always knows exactly where one message ends and the next begins. This low-level skill underlies every P2P message GoChain will send from here on.

**Key topics:** `net.Listen`, `net.Dial`, TCP framing, length-prefixed messages, handling one connection per goroutine.

---

### [Chapter 45: Designing a P2P Protocol](45-designing-a-p2p-protocol.md)

Before more Go code, we design GoChain's wire protocol on paper: message types (`version` for handshakes, `getblocks`, `inv` for inventory announcements, `getdata`, `block`, `tx`), and a simple binary envelope (message type + length + payload) every message follows — closely modeled on Bitcoin's own original protocol, so the design choices carry over directly to real systems.

**Key topics:** designing a message protocol, message types, envelope format, why protocol design happens before implementation.

---

### [Chapter 46: Building a P2P Node in Go](46-building-a-p2p-node-in-go.md)

Implementing `network.Node`: it listens for incoming connections, dials outgoing ones, and routes each incoming message to a handler based on its type from Chapter 45. We connect two GoChain nodes running on the same machine (different ports) and watch them exchange their first handshake message.

**Key topics:** the `network.Node` type, message routing, your first working two-node handshake.

---

### [Chapter 47: Peer Discovery and Handshakes](47-peer-discovery-and-handshakes.md)

New nodes need a way to find others without a central directory. We implement a simple seed-node approach (a small hardcoded or configurable list of known-good addresses to connect to first) and a peer-exchange message that lets connected nodes share the addresses of *other* peers they know about, so the network's connectivity grows organically from just a few seed nodes.

**Key topics:** seed nodes, peer exchange (addr messages), building a peer address book, growing network connectivity.

---

### [Chapter 48: Gossip Protocols and Broadcasting](48-gossip-protocols-and-broadcasting.md)

Instead of every node connecting to every other node (which does not scale), gossip protocols spread information by having each node forward what it hears to a handful of its own peers, who forward it further — like a rumor spreading through a crowd, reaching everyone in a few "hops" without anyone needing to shout to the whole room at once. We implement gossip-based broadcasting for new transactions, including duplicate-message suppression so the same transaction is not forwarded forever in a loop.

**Key topics:** gossip protocols, broadcast without a central hub, duplicate suppression, hop-based propagation.

---

### [Chapter 49: Blockchain Synchronization](49-blockchain-synchronization.md)

When a node joins the network (or reconnects after being offline), it needs to catch up. We implement chain synchronization: a new node asks its peers for their current chain height, requests missing blocks in order, and validates each one exactly as in Volume 3 before accepting it — refusing anything that fails validation, even from a trusted-looking peer.

**Key topics:** initial block download, requesting missing blocks, validating blocks received from the network, sync progress.

---

### [Chapter 50: Handling Forks — The Longest Chain Rule](50-handling-forks-the-longest-chain-rule.md)

Two miners can solve a valid block at nearly the same moment, splitting the network's view into two competing chains (a fork). We implement the longest-chain rule (more precisely, the chain with the most accumulated proof of work): nodes track competing chains and switch to whichever one becomes longest, "reorganizing" by rolling back and replaying transactions as needed — demonstrated with a hands-on simulation of two nodes mining a fork and one chain eventually winning.

**Key topics:** forks, the longest/heaviest chain rule, chain reorganization, simulating and resolving a fork.

---

### [Chapter 51: Sybil and Eclipse Attacks — Basics](51-sybil-and-eclipse-attacks-basics.md)

A Sybil attack floods the network with fake identities to gain disproportionate influence. An eclipse attack isolates a specific victim node by surrounding it entirely with attacker-controlled peers, feeding it a false view of the network. We explain both conceptually with diagrams and discuss GoChain's (and real blockchains') practical mitigations: diverse peer selection, connection limits, and outbound-connection preferences.

**Key topics:** Sybil attacks, eclipse attacks, peer diversity as a defense, why identity is "free" in P2P networks and what that implies.

---

### [Chapter 52: Major Project 2 — A Real Multi-Node Network](52-major-project-2-a-real-multi-node-network.md)

Run three or more independent GoChain node processes on your machine (or across machines on your home network), each with its own wallet. Submit a transaction on Node A, watch it gossip to Nodes B and C, mine it on Node B, and watch the resulting block propagate and get validated everywhere — with a short script that automates starting the whole network and asserting every node converges on an identical chain.

**Key topics:** running a real multi-node network, observing gossip and sync end-to-end, automated convergence testing.

**Major Project:** A scripted, multi-node GoChain testnet running entirely on your own machine, with automated convergence checks.

---

# PART 4 — STORAGE, CONTRACTS & TOOLING

---

## Volume 8: Storage & State

> The flat-file storage from Volume 3 will not survive a real network's load or a crash mid-write. This volume replaces it with a real embedded database and adds the indexes a usable blockchain needs.

### [Chapter 53: Why Flat Files Are Not Enough](53-why-flat-files-are-not-enough.md)

Concretely reproducing the problems: what happens if GoChain crashes while appending a block (partial write, corrupted file)? How slow is "scan every block to find one transaction" once you have tens of thousands of blocks? What happens when two goroutines try to write at once? Each problem is demonstrated hands-on with the existing Volume 3 storage code before we fix it.

**Key topics:** crash-safety problems, linear scan performance, concurrent write hazards, motivating a real storage engine.

---

### [Chapter 54: Embedded Databases — BoltDB and Badger](54-embedded-databases-boltdb-and-badger.md)

An embedded database runs inside your own process (no separate server), stores data in a single file, and gives you crash-safe, indexed key-value storage for free. We compare BoltDB (simple, B-tree based, single-writer) and Badger (LSM-tree based, built for higher write throughput) and choose BoltDB for GoChain's clarity-first storage layer, with Badger discussed as the production-scale alternative.

**Key topics:** embedded key-value databases, BoltDB vs. Badger, B-trees vs. LSM-trees at a conceptual level, choosing a storage engine.

---

### [Chapter 55: Designing the Storage Layer](55-designing-the-storage-layer.md)

Defining `storage.Store`, an interface with methods like `PutBlock`, `GetBlock`, `PutUTXO`, `GetUTXO`, and `Iterator` — so the rest of GoChain depends on this interface, not on BoltDB directly, meaning we could swap in Badger (or a test-only in-memory store) later without touching any other package. Designing BoltDB's "buckets" (its equivalent of tables) for blocks, transactions, and the UTXO set.

**Key topics:** the storage interface pattern, designing BoltDB buckets, decoupling business logic from a specific database.

---

### [Chapter 56: Fast UTXO Lookups and Indexing](56-fast-utxo-lookups-and-indexing.md)

Scanning every transaction in every block to compute a balance does not scale. We build and maintain a dedicated UTXO index inside BoltDB, updated incrementally every time a block is mined or received (add new outputs, remove spent ones), so `BalanceOf(address)` becomes a fast, direct lookup instead of a full-chain scan — with a benchmark comparing both approaches on a chain with thousands of blocks.

**Key topics:** UTXO set indexing, incremental index maintenance, benchmarking scan vs. indexed lookup.

---

### [Chapter 57: Merkle-Patricia Tries and State](57-merkle-patricia-tries-and-state.md)

A deeper look at the data structure Ethereum uses to commit to its entire account state in one root hash: the Merkle-Patricia Trie, a hybrid of a Merkle tree (for tamper-evidence) and a Patricia trie (a prefix tree optimized for keys that share common prefixes, like addresses). We build a simplified version in Go and use it to commit to GoChain's UTXO set with a single root hash stored in each block header, previewing exactly the kind of state root smart contracts will need in Volume 9.

**Key topics:** Merkle-Patricia tries, prefix trees, committing to state with a single root hash, building a simplified trie in Go.

---

### [Chapter 58: Mini Project — Fast Balance Index](58-mini-project-fast-balance-index.md)

Build and benchmark a standalone balance-index service: given a folder of exported blocks, it builds the UTXO index from Volume 8 from scratch, then serves balance lookups over a tiny local API, with before/after benchmarks proving the index turns an O(n) scan into an O(1) lookup.

**Key topics:** building an index from historical data, exposing it via a simple API, proving performance improvements with real numbers.

**Mini Project:** `balanceindex` — a standalone service that builds and serves a fast UTXO-backed balance index.

---

## Volume 9: Smart Contracts & a Virtual Machine

> Plain transfers can only do so much. Smart contracts let a blockchain run small programs that move value based on logic, not just a signature. This volume builds a real, working virtual machine and a token contract that runs on it.

### [Chapter 59: What Are Smart Contracts?](59-what-are-smart-contracts.md)

A smart contract is a program stored on the blockchain that runs identically on every node, so its results are trustworthy without trusting any single party. We ground this with a concrete example before any code: an escrow contract that only releases funds to a seller once a buyer confirms delivery — no intermediary needed, because the *rules* are enforced by every node re-executing the same program and agreeing on the result.

**Key topics:** what a smart contract is, why deterministic execution matters, an escrow walkthrough, contracts vs. plain transactions.

---

### [Chapter 60: Stack-Based Virtual Machines](60-stack-based-virtual-machines.md)

Rather than running arbitrary native code (dangerous and non-deterministic across machines), blockchains run a small, restricted virtual machine. A stack-based VM keeps a single stack of values and executes simple instructions (push a number, add the top two values, compare them) that each read from and write back to that stack — the same core idea behind Bitcoin Script and, at a conceptual level, the Ethereum Virtual Machine (EVM). We trace a tiny program by hand, instruction by instruction, watching the stack change.

**Key topics:** stack machines, why not run native code on-chain, tracing instruction execution by hand.

---

### [Chapter 61: Designing an Instruction Set](61-designing-an-instruction-set.md)

Designing GoChain VM's opcode table: arithmetic (`OpAdd`, `OpSub`), stack manipulation (`OpDup`, `OpPop`), comparison (`OpEqual`, `OpGreaterThan`), control flow (`OpJump`, `OpJumpIfFalse`), and blockchain-specific opcodes (`OpCheckSig` to verify a signature from within a contract, `OpBalance` to read an account's balance). Each opcode is specified with its exact stack effect before implementation.

**Key topics:** opcode design, stack effect notation, arithmetic and control-flow instructions, blockchain-specific opcodes.

---

### [Chapter 62: Building Our VM in Go](62-building-our-vm-in-go.md)

Implementing `vm.VM`: a `Stack`, a program counter, and an `Execute()` loop that fetches the next opcode, dispatches to its handler, and updates the stack and program counter — with a `switch` statement over opcodes as the interpreter's core, and unit tests for every single opcode from Chapter 61.

**Key topics:** implementing an interpreter loop in Go, the `switch`-based dispatch pattern, testing a VM opcode by opcode.

---

### [Chapter 63: A Scripting Language for Locking Coins](63-a-scripting-language-for-locking-coins.md)

Real UTXOs are not just "owned by an address" — they carry a small locking script that must evaluate to true for the output to be spendable, and the spender provides an unlocking script (typically a signature) that gets run together with it. We implement this exact pattern for GoChain: a standard "pay to public key hash" locking script that our VM evaluates using `OpCheckSig`, unifying "plain" transactions and "smart" ones under one execution model.

**Key topics:** locking and unlocking scripts, pay-to-public-key-hash, unifying transactions and contracts through the VM.

---

### [Chapter 64: Gas and Execution Limits](64-gas-and-execution-limits.md)

An infinite loop in a smart contract would hang every node in the network forever. Gas solves this: every opcode costs a small, fixed amount of gas, the caller supplies a gas limit up front, and execution halts the instant that limit is exceeded — turning "runs forever" into "fails safely and predictably, costing the caller only what they agreed to spend." We add gas accounting to every opcode handler from Chapter 62 and test a deliberately infinite loop running out of gas exactly as expected.

**Key topics:** why gas exists, per-opcode gas costs, gas limits and out-of-gas failures, testing a runaway contract.

---

### [Chapter 65: Writing a Simple Token Contract](65-writing-a-simple-token-contract.md)

Using the opcodes and scripting model built so far, we design and deploy GoChain's first real application-level contract: a fungible token, with `transfer`, `balanceOf`, and `mint` operations, each backed by a small VM program and its own contract-local storage (introduced properly in the next chapter).

**Key topics:** designing a token contract's operations, mapping high-level operations onto VM opcodes, a first deployed contract.

---

### [Chapter 66: Contract Storage and State](66-contract-storage-and-state.md)

Unlike a plain transaction, a contract needs persistent storage of its own — the token contract's balance table must survive between calls. We add a contract-storage layer keyed by contract address plus a storage slot, backed by the `storage.Store` interface from Volume 8, and wire `OpSLoad`/`OpSStore` opcodes so contract code can read and write it.

**Key topics:** persistent contract storage, storage slots, `SLOAD`/`SSTORE` opcodes, isolating one contract's storage from another's.

---

### [Chapter 67: Reentrancy and Contract Security Pitfalls](67-reentrancy-and-contract-security-pitfalls.md)

The most infamous smart-contract bug class: a contract calls out to another (potentially malicious) contract mid-execution, and that other contract calls back into the first one before its state has finished updating — draining funds through repeated withdrawals. We reproduce a simplified version of the real 2016 DAO hack's bug pattern inside GoChain's own VM, then fix it using the checks-effects-interactions pattern (update your own state *before* making external calls).

**Key topics:** reentrancy attacks, the DAO hack pattern explained, checks-effects-interactions, general smart contract security hygiene.

---

### [Chapter 68: Testing Smart Contracts](68-testing-smart-contracts.md)

Writing a proper test suite for the token contract: unit tests for each operation, an adversarial test that attempts the reentrancy bug from Chapter 67 and asserts it now fails safely, and gas-consumption assertions that catch accidental performance regressions in contract code — treating contract code with the same (or greater) rigor as any other financial software, because contract bugs cannot be quietly patched after the fact on a real chain.

**Key topics:** contract test suites, adversarial testing, gas regression tests, why contract bugs are especially costly.

---

### [Chapter 69: Major Project 3 — Deploy a Token Contract](69-major-project-3-deploy-a-token-contract.md)

End to end: deploy the token contract from Chapter 65 onto a running GoChain node, mint an initial supply to your wallet, transfer tokens to a second wallet through a real signed transaction that invokes the contract, and query balances directly from contract storage — the same fundamental flow behind every ERC-20 token on Ethereum, built and running on your own chain.

**Key topics:** full contract deployment and invocation lifecycle, minting and transferring tokens, querying contract state.

**Major Project:** A deployed, working GoChain token contract with mint and transfer, exercised end to end.

---

## Volume 10: APIs, Explorer & Developer Tools

> A blockchain nobody can query or inspect is not very useful to anyone but you. This volume builds the API, the block explorer, and the polished CLI that make GoChain usable by other programs and other people.

### [Chapter 70: Building a JSON-RPC and REST API](70-building-a-json-rpc-and-rest-api.md)

Most real blockchains expose a JSON-RPC API (a simple "call a named method with parameters, get a JSON result back" convention) alongside — or instead of — a REST API. We implement both for GoChain: JSON-RPC methods like `getBlock`, `getBalance`, and `sendTransaction`, and equivalent REST endpoints, using Go's standard `net/http`, so other programs (wallets, explorers, exchanges) can talk to a GoChain node without needing to speak our internal P2P protocol.

**Key topics:** JSON-RPC design, REST API design, `net/http` servers in Go, exposing node functionality safely.

---

### [Chapter 71: Real-Time Updates with WebSockets](71-real-time-updates-with-websockets.md)

Polling an API repeatedly to check "is there a new block yet?" is wasteful. We add a WebSocket endpoint that pushes new-block and new-transaction events to subscribed clients the instant they happen, using Go's `gorilla/websocket` (or the standard library's newer WebSocket support), and build a tiny terminal client that prints each event live as it streams in.

**Key topics:** WebSockets vs. polling, pushing live blockchain events, a real-time terminal client.

---

### [Chapter 72: Building a Block Explorer Backend](72-building-a-block-explorer-backend.md)

A block explorer is a website that lets anyone search and browse a blockchain — think of a site like Etherscan. We build its backend: endpoints for browsing recent blocks, viewing a single block's full transaction list, viewing a single transaction's inputs and outputs, and searching by address to see its full transaction history, all backed by the indexes from Volume 8.

**Key topics:** block explorer data requirements, address history queries, pagination for large result sets.

---

### [Chapter 73: Block Explorer Frontend — Overview](73-block-explorer-frontend-overview.md)

A conceptual, practically-oriented tour of building a minimal explorer frontend that consumes the Chapter 72 API: a homepage listing recent blocks, a block detail page, a transaction detail page, and an address page — enough HTML/JS (or a minimal framework of your choice) to have something real to click through and demo, with the emphasis kept on how the frontend maps onto the backend API rather than frontend framework mastery.

**Key topics:** mapping API endpoints to explorer pages, a minimal working frontend, what a production explorer adds beyond this.

---

### [Chapter 74: Building a Polished CLI with Cobra](74-building-a-polished-cli-with-cobra.md)

Replacing the ad-hoc `flag`-based CLI commands built across earlier volumes with a single, well-organized `gochain` command built on the popular Cobra library: `gochain node start`, `gochain wallet new`, `gochain tx send`, `gochain chain inspect`, each with proper help text, flags, and subcommands — the same pattern used by tools like `kubectl`, `docker`, and `go` itself.

**Key topics:** the Cobra CLI framework, subcommands and flags, designing a coherent CLI surface, help text and discoverability.

---

### [Chapter 75: Mini Project — Block Explorer API](75-mini-project-block-explorer-api.md)

A focused, deployable version of the Chapter 72-73 explorer: a single Go binary that serves both the JSON API and a minimal bundled frontend, configurable to point at any running GoChain node, ready to be containerized in Part 6.

**Key topics:** bundling a frontend with a Go binary, configuration for pointing at different nodes, preparing for deployment.

**Mini Project:** `gochain-explorer` — a self-contained block explorer service for any GoChain node.

---

# PART 5 — SECURITY & THE WIDER WORLD

---

## Volume 11: Security, Alternative Consensus & Scaling

> Understanding attacks, alternative consensus algorithms, and scaling techniques turns you from someone who *built* a blockchain into someone who *understands* blockchains broadly — including the ones you did not build.

### [Chapter 76: Common Attacks and an Attack Simulation Lab](76-common-attacks-and-an-attack-simulation-lab.md)

A tour of the attacks that matter most in practice: the 51% attack (an attacker controlling a majority of mining power can rewrite recent history), double-spending in more detail than Volume 5's basic case, replay attacks (rebroadcasting a valid transaction from one context in another), and race attacks on unconfirmed transactions. We build a hands-on lab: a modified GoChain node that deliberately attempts a 51% attack against a small local test network, so you watch the attack — and the honest network's defense — happen in real time.

**Key topics:** the 51% attack in depth, double-spend and replay attacks, a hands-on attack simulation against your own network.

---

### [Chapter 77: Proof of Stake, Explained and Implemented](77-proof-of-stake-explained-and-implemented.md)

Proof of stake replaces "spend electricity to earn the right to propose a block" with "put up a financial stake, and risk losing it (slashing) if you misbehave." We explain validator selection, staking, and slashing conceptually, then implement `consensus.ProofOfStake` as an alternative, swappable consensus engine for GoChain (behind the same interface used by `consensus.ProofOfWork`), so you can run a small GoChain testnet under either algorithm and directly compare their behavior.

**Key topics:** proof of stake fundamentals, validator selection, slashing, implementing PoS in Go behind a shared consensus interface.

---

### [Chapter 78: BFT Consensus — PBFT and Tendermint, Overview](78-bft-consensus-pbft-and-tendermint-overview.md)

Byzantine Fault Tolerant (BFT) consensus algorithms tolerate not just crashed nodes but actively malicious ones, using explicit multi-round voting among known validators to reach agreement — the approach used by permissioned chains like Hyperledger Fabric and by Tendermint (which powers the Cosmos ecosystem). We walk through PBFT's three-phase voting process and Tendermint's round-based model conceptually, with sequence diagrams, without a full implementation, so you can recognize and reason about these systems when you meet them in the real world.

**Key topics:** Byzantine fault tolerance, PBFT's phases, Tendermint's round-based BFT, permissioned vs. permissionless consensus.

---

### [Chapter 79: Sharding, Layer 2, and Rollups](79-sharding-layer-2-and-rollups.md)

As a chain grows, a single sequence of blocks processed by every node becomes a bottleneck. Sharding splits the network into parallel sub-chains that each handle a portion of the load. Layer 2 systems (state channels, rollups) move most computation off the main chain while still inheriting its security, periodically settling a compressed summary back to it. We explain optimistic rollups and zero-knowledge rollups conceptually, and diagram exactly what "the main chain" still has to verify in each approach.

**Key topics:** sharding, state channels, optimistic rollups, zero-knowledge rollups, what scaling techniques actually trade away.

---

### [Chapter 80: Choosing the Right Consensus and Architecture](80-choosing-the-right-consensus-and-architecture.md)

A decision framework, in the same spirit as choosing a database earlier in your learning: permissionless vs. permissioned, proof of work vs. proof of stake vs. BFT, when sharding or layer 2 is actually necessary versus premature optimization. Real projects and which choices they made, and why, as concrete anchors for the framework.

**Key topics:** a consensus/architecture decision framework, matching design to actual requirements, avoiding premature scaling.

---

## Volume 12: Real-World Case Studies

> You have built enough of your own blockchain to read real blockchain source code and design docs with genuine understanding. This volume turns that understanding outward.

### [Chapter 81: Bitcoin Architecture, Deep Dive](81-bitcoin-architecture-deep-dive.md)

A structural tour of Bitcoin itself, mapped explicitly back onto GoChain's own design: its UTXO model (identical in spirit to GoChain's), its actual block and transaction formats, Bitcoin Script (the real ancestor of GoChain's VM), its real difficulty adjustment algorithm, and the mempool policies real Bitcoin nodes use — reading excerpts of Bitcoin Core's actual behavior descriptions with the vocabulary you now already know cold.

**Key topics:** Bitcoin's real architecture, comparing it point-by-point with GoChain, Bitcoin Script, real difficulty adjustment.

---

### [Chapter 82: Ethereum Architecture, Deep Dive](82-ethereum-architecture-deep-dive.md)

Ethereum's account-based state model, the real Ethereum Virtual Machine and its actual opcode set (compared directly with GoChain's smaller VM from Volume 9), gas pricing in a live, congested network, and the Ethereum state trie (the real, full version of the simplified trie you built in Volume 8). The 2022 transition from proof of work to proof of stake ("The Merge") as a real-world case study in swapping consensus algorithms in a live, running system.

**Key topics:** Ethereum's account model and state trie, the real EVM, gas pricing in production, The Merge as a case study.

---

### [Chapter 83: Comparing GoChain to Real-World Chains](83-comparing-gochain-to-real-world-chains.md)

A structured, side-by-side comparison table across every major design axis covered in this course — data model, consensus, virtual machine, networking, storage — between GoChain, Bitcoin, and Ethereum, making explicit exactly which simplifications GoChain made for teachability and what a production chain adds on top of each one.

**Key topics:** structured architecture comparison, identifying deliberate simplifications, the gap between a teaching chain and a production one.

---

### [Chapter 84: Permissioned Chains — Hyperledger Fabric, Overview](84-permissioned-chains-hyperledger-fabric-overview.md)

Not every blockchain is public. Permissioned chains like Hyperledger Fabric restrict who can even join the network, which is a natural fit for consortiums of companies (banks, supply chain partners) who want shared, tamper-evident record-keeping without full public transparency. We explain Fabric's channel-based privacy model and its "endorsement" transaction flow conceptually, and discuss when a permissioned design is the right engineering choice rather than a public chain.

**Key topics:** permissioned vs. permissionless chains, Hyperledger Fabric's channels and endorsement flow, when permissioned chains make sense.

---

### [Chapter 85: Lessons from Real-World Incidents and Hacks](85-lessons-from-real-world-incidents-and-hacks.md)

Case studies of real, consequential blockchain incidents — the 2016 DAO reentrancy hack (which you already reproduced in miniature in Volume 9), major exchange hacks caused by key-management failures rather than protocol flaws, and at least one real 51%-attack incident on a smaller proof-of-work chain — each broken down into what actually went wrong, and which specific chapter of this course covers the concept that would have caught it.

**Key topics:** the DAO hack in full, key-management failures at exchanges, real 51% attack incidents, connecting incidents back to course concepts.

---

# PART 6 — SHIPPING IT

---

## Volume 13: Deployment & Going Live

> Everything you have built has run on your own laptop. This final volume takes GoChain to real servers, wires up monitoring and CI/CD, and ends with your own live, public testnet that other people can actually connect to.

### [Chapter 86: Containerizing a Node with Docker](86-containerizing-a-node-with-docker.md)

Writing a `Dockerfile` that builds the GoChain binary in one stage and produces a small, final runtime image in a second stage (multi-stage builds) — so your deployed image contains just the compiled binary and nothing else, keeping it small, fast to pull, and secure. Running your first containerized GoChain node and connecting to it from your host machine.

**Key topics:** Docker fundamentals, multi-stage builds, minimal runtime images, running and inspecting a containerized GoChain node.

---

### [Chapter 87: Docker Compose — A Multi-Node Testnet](87-docker-compose-a-multi-node-testnet.md)

A `docker-compose.yml` that starts several GoChain nodes, each in its own container, pre-configured to discover and connect to each other automatically, plus the block explorer from Volume 10 pointed at one of them — turning the manual multi-terminal setup from Volume 7's major project into a single `docker-compose up` command.

**Key topics:** `docker-compose.yml`, multi-container networking, service discovery between containers, one-command testnets.

---

### [Chapter 88: Deploying to a Cloud VM](88-deploying-to-a-cloud-vm.md)

Provisioning a real virtual machine on a cloud provider (DigitalOcean, AWS EC2, or similar — the concepts transfer directly regardless of which you pick), installing Docker on it, and running your `docker-compose` testnet on a real, internet-reachable server for the first time, with basic firewall rules opening only the ports GoChain actually needs.

**Key topics:** provisioning a cloud VM, remote Docker deployment, firewall and port configuration, your first internet-reachable node.

---

### [Chapter 89: Kubernetes for Blockchain Nodes, Overview](89-kubernetes-for-blockchain-nodes-overview.md)

For a network of nodes larger than `docker-compose` comfortably manages, Kubernetes automates deployment, scaling, and recovery. We explain Kubernetes' core concepts (pods, deployments, services) conceptually, and walk through a sample Kubernetes manifest for running a small GoChain node deployment, framed clearly as an optional next step rather than a requirement for the rest of this volume.

**Key topics:** Kubernetes core concepts, pods/deployments/services, a sample GoChain Kubernetes manifest, when Kubernetes is (and is not) worth the complexity.

---

### [Chapter 90: Setting Up a Public Testnet and a Faucet](90-setting-up-a-public-testnet-and-a-faucet.md)

Turning your deployed nodes into an actual public testnet: publishing seed-node addresses so anyone's GoChain node can join and sync, and building a faucet — a small web service that gives out free test coins to any address that requests them, exactly like real blockchain testnets do, so new users (and your own future demos) can get started without needing to mine first.

**Key topics:** public seed nodes, onboarding external participants, building a coin faucet service, testnet vs. mainnet framing.

---

### [Chapter 91: Monitoring with Prometheus and Grafana](91-monitoring-with-prometheus-and-grafana.md)

Instrumenting GoChain with Prometheus metrics — block height, mempool size, peer count, mining hash rate, API request latency — using Go's `prometheus/client_golang` library, and building a Grafana dashboard that visualizes all of it live, so you (and anyone running a node) can actually see the network's health rather than guessing at it from logs.

**Key topics:** Prometheus metrics in Go, key blockchain metrics to track, building a Grafana dashboard, observability as an operational habit.

---

### [Chapter 92: CI/CD for a Blockchain Project](92-ci-cd-for-a-blockchain-project.md)

Setting up GitHub Actions to automatically run GoChain's full test suite (every package's tests from every volume) on every push, build the Docker image, and — on tagged releases — push it to a container registry, so deploying a new version becomes a matter of pushing a tag rather than manually rebuilding and redeploying by hand.

**Key topics:** GitHub Actions basics, automated test runs, automated Docker builds, tagged releases and continuous deployment.

---

### [Chapter 93: TLS, Domains, and Hosting the Explorer](93-tls-domains-and-hosting-the-explorer.md)

Putting a real domain name and HTTPS (via Let's Encrypt/Certbot, or a reverse proxy like Caddy that handles it automatically) in front of the block explorer from Volume 10, so it is reachable at a real, secure URL instead of a bare IP address and port — the last piece that makes your testnet feel like a genuine, presentable, shareable product rather than a personal experiment.

**Key topics:** domain names and DNS basics, TLS/HTTPS, reverse proxies, automatic certificate renewal.

---

### [Chapter 94: Backups and Disaster Recovery](94-backups-and-disaster-recovery.md)

What happens if a node's disk fails? We implement scheduled backups of each node's BoltDB data file and wallet files, a documented restore procedure, and a chaos-test: deliberately kill a node's container, restore it from backup, and verify it rejoins the network and catches up correctly — because a backup you have never tested restoring is not a real backup.

**Key topics:** backup strategy for blockchain data, wallet backup, restore procedures, chaos-testing your own recovery process.

---

### [Chapter 95: Final Capstone — Launch Your Own Testnet Blockchain](95-final-capstone-launch-your-own-testnet-blockchain.md)

The capstone. Using everything from all thirteen volumes, you stand up a complete, live, monitored, multi-node GoChain testnet: several nodes deployed across real cloud servers, discovering each other over the public internet, mining and gossiping blocks and transactions, a running faucet, a publicly reachable block explorer behind HTTPS, Prometheus/Grafana monitoring, and CI/CD wired up so future changes deploy automatically. You invite someone else to run a GoChain node, connect it to your testnet, and watch it sync from scratch — proof that the entire system you built genuinely, independently works.

**Key topics:** full-system integration, a real live deployment checklist, onboarding an external node operator, what "done" looks like for this course.

**Capstone Project:** A live, public, monitored, multi-node GoChain testnet, deployable and joinable by someone other than you.

---

## Summary — What You Have Built

By the time you reach the end of Chapter 95, here is what you have done:

- **Understood blockchain from first principles** — not as a buzzword, but as a specific, well-understood combination of hashing, digital signatures, consensus, and peer-to-peer networking. You can read Bitcoin's or Ethereum's design documents and actual behavior and understand exactly what is happening and why.

- **Built GoChain** — a real blockchain in Go with a hand-built cryptography layer, a proof-of-work consensus engine (and an alternative proof-of-stake engine), a UTXO-based transaction model with a working mempool, HD wallets with encrypted key storage, a peer-to-peer networking layer with gossip and fork resolution, a real embedded-database storage layer with fast indexes, a stack-based virtual machine with gas metering, and a deployed token smart contract.

- **Built the tools around it** — a CLI wallet, a JSON-RPC/REST API, a WebSocket live-event feed, and a block explorer, mirroring the tooling every real blockchain ecosystem needs to be usable by anyone other than its own developers.

- **Understood security in depth** — common attacks (51%, double-spend, replay, Sybil, eclipse), smart contract vulnerabilities (reentrancy and the DAO hack), and how real incidents map back onto specific concepts you now know cold.

- **Deployed a live system** — containerized with Docker, orchestrated with Docker Compose (and, optionally, Kubernetes), monitored with Prometheus and Grafana, automated with CI/CD, and reachable by other people over the real internet through a proper domain, HTTPS, and a working faucet.

- **Developed engineering judgment** — the ability to look at any blockchain project, new or old, permissioned or public, and ask the right questions: how does it reach consensus? How does it prevent double-spending? How does it scale? What did it choose to trade away, and why? These questions will serve you for the rest of your career, in blockchain engineering or anywhere else.

---

## A Note on the Journey

This course does not end with a slide deck about "what blockchain is." It ends with a folder of Go code, written by you, that mines blocks, moves value, runs contracts, talks to other computers over the internet, and is currently running on a server somewhere, reachable by anyone who wants to connect to it.

That is the entire point. Blockchain stops being a mysterious word the moment you have built one — even a small one — with your own hands, and watched it actually work. Every "how does Bitcoin really prevent double-spending" or "what does gas actually mean" question you will ever be asked again, you will be able to answer not from memory, but from having built the exact mechanism yourself and watched it succeed and fail under your own tests.

That knowledge does not go away. It is yours.

Now turn to Chapter 1, and let's build GoChain.

---

*Blockchain in Go — From Zero to Building, Networking, and Deploying Your Own Live Blockchain in Go*
