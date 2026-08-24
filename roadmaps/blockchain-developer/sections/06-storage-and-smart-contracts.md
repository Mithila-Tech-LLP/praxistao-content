---
title: Storage & Smart Contracts
---
A production node can't keep its entire chain in memory, and "smart contract" chains need a way to run untrusted code safely. This section covers real persistent storage and building a small virtual machine.

### Embedded Databases: BoltDB & Badger
Flat files don't scale past a toy project. Embedded key-value databases give you fast, durable, indexed storage without running a separate database server.

**Resources:**
- [Embedded Databases: BoltDB and Badger](course:blockchain-in-go#54-embedded-databases-boltdb-and-badger)

### Fast UTXO Lookups & State Tries
> optional

As the UTXO set grows, naive lookups get slow. See the indexing strategies and the Merkle-Patricia trie structure Ethereum uses to track account state efficiently.

**Resources:**
- [Fast UTXO Lookups and Indexing](course:blockchain-in-go#56-fast-utxo-lookups-and-indexing)
- [Merkle-Patricia Tries and State](course:blockchain-in-go#57-merkle-patricia-tries-and-state)

### What Are Smart Contracts
A smart contract is code that runs on the blockchain itself, with its execution and results agreed on by every node. Understand what problem that actually solves before building one.

**Resources:**
- [What Are Smart Contracts](course:blockchain-in-go#59-what-are-smart-contracts)

### Stack-Based Virtual Machines
Most smart-contract platforms run on a small, deterministic virtual machine. Build one — a stack-based VM is the simplest design that still runs real logic.

**Resources:**
- [Stack-Based Virtual Machines](course:blockchain-in-go#60-stack-based-virtual-machines)
- [Building Our VM in Go](course:blockchain-in-go#62-building-our-vm-in-go)

### Practice: A Tiny Contract VM
> branches-from: Stack-Based Virtual Machines

Build the VM task of the standalone project: a working stack machine that executes bytecode with gas metering.

**Resources:**
- [Build a Blockchain project](project:build-a-blockchain)

### Gas & Execution Limits
Without a cost model, a contract could loop forever and stall the whole network. Gas is how you make execution cost something and cap it.

**Resources:**
- [Gas and Execution Limits](course:blockchain-in-go#64-gas-and-execution-limits)

### Writing & Testing a Token Contract
Apply the VM to something concrete: a token contract with `mint`, `transfer`, and `balanceOf`, plus the tests that prove it behaves correctly.

**Resources:**
- [Writing a Simple Token Contract](course:blockchain-in-go#65-writing-a-simple-token-contract)
- [Testing Smart Contracts](course:blockchain-in-go#68-testing-smart-contracts)

### Major Project: Deploy a Token Contract
Deploy the token contract to a running GoChain node and interact with it end-to-end — mint, transfer, and query balances against a live chain.

**Resources:**
- [Major Project 3: Deploy a Token Contract](course:blockchain-in-go#69-major-project-3-deploy-a-token-contract)
