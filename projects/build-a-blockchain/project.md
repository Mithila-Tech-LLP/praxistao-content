---
title: Build a Blockchain
subtitle: Build a working proof-of-work blockchain in Go — blocks, mining, wallets, transactions, a P2P network, and a tiny smart-contract VM
category: Systems Programming
difficulty: intermediate
duration: 8-12 hours
accent: "#2dd4bf"
technologies: [Go]
skills: [Cryptographic Hashing, Digital Signatures, Proof of Work, UTXO Model, Concurrency, TCP Networking, Merkle Trees, Virtual Machines]
prerequisites: [basic-programming]
repo: build-a-blockchain
outcomes:
  - Build a tamper-evident chain of blocks linked by SHA-256 hashes
  - Implement a real proof-of-work mining loop with an adjustable difficulty target
  - Generate ECDSA key pairs and sign and verify data exactly like a real wallet
  - Build a UTXO-based transaction model with change outputs and a mempool
  - Implement a Merkle tree and a Merkle inclusion proof
  - Send blocks and transactions between two independent node processes over TCP
  - Build a minimal stack-based virtual machine with gas metering
  - Wire everything into a single command-line wallet that can mine, send, and check balances
---

## Overview

Bitcoin, Ethereum, and every blockchain since inherited the same handful of core ideas: a chain of blocks linked by hashes, proof of work to make rewriting history expensive, digital signatures to prove ownership, and a peer-to-peer network to keep everyone in sync. None of it is magic — it's cryptographic hashing, some careful data structures, and ordinary networking code, combined in a specific, clever order.

In this project you build a real, working blockchain from scratch in Go — not a toy that only makes sense on paper. By the end you'll have a chain that mines real proof-of-work blocks, wallets that sign real transactions, a mempool that rejects real double-spends, two node processes that gossip blocks to each other over a real TCP connection, and a tiny virtual machine that runs real, gas-metered bytecode.

Each task is self-contained — its own Go module, its own tests, its own focused piece of the system — so you can see exactly what you built and why at every step. This project is a condensed, hands-on companion to the full **[Blockchain in Go](/course/blockchain-in-go)** course, which goes far deeper into every one of these topics (and adds smart contracts, HD wallets, a block explorer, and real cloud deployment) if you want to keep going afterward.

## What You'll Build

By Task 10 you'll have a single CLI, `chain-wallet`, that can generate a keypair, check a balance, send coins, and mine a block — the same fundamental user experience as a real cryptocurrency wallet, backed entirely by code you wrote yourself:

```bash
$ chain-wallet new
Your new address: 1J9k...Qz7f

$ chain-wallet mine -miner 1J9k...Qz7f
Mined block 1 (reward: 50 coins)

$ chain-wallet balance -address 1J9k...Qz7f
Balance: 50

$ chain-wallet send -from 1J9k...Qz7f -to 1Bv1...k7c -amount 20
Transaction a1b2c3... submitted to the mempool

$ chain-wallet mine -miner 1J9k...Qz7f
Mined block 2 (reward: 50 coins, 1 transaction included)
```
