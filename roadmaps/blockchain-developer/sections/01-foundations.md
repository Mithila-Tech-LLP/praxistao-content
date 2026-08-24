---
title: Foundations
---
Start with the concept before the code: what a blockchain actually is, why Go fits the problem well, and getting a project scaffold in place.

### What Is a Blockchain, Really
Strip away the hype and a blockchain is a specific, fairly simple data structure: an append-only chain of records, each one cryptographically linked to the one before it.

**Resources:**
- [What Is a Blockchain, Really](course:blockchain-in-go#01-what-is-a-blockchain-really)

### Why Go for Blockchain
Go's static typing, straightforward concurrency, and small binaries are exactly why most real blockchain node software (Ethereum's go-ethereum, Hyperledger Fabric) is written in it.

**Resources:**
- [Why Go Is a Great Fit for Blockchain](course:blockchain-in-go#02-why-go-is-a-great-fit-for-blockchain)

### Environment Setup & Project Scaffold
Get your Go module, folder layout, and tooling in place before writing any blockchain-specific code.

**Resources:**
- [Environment Setup and Project Scaffold](course:blockchain-in-go#03-environment-setup-and-project-scaffold)

### Go Crash Course
> optional

If you're new to Go, this is a fast on-ramp covering exactly what the rest of the roadmap assumes you know.

**Resources:**
- [A Go Crash Course for Blockchain Developers](course:blockchain-in-go#04-a-go-crash-course-for-blockchain-developers)

### Concurrency Basics
A blockchain node does many things simultaneously — mining, networking, serving API requests. Goroutines and channels are how Go handles that.

**Resources:**
- [Concurrency: Goroutines and Channels](course:blockchain-in-go#05-concurrency-goroutines-and-channels)
