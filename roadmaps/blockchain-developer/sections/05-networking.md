---
title: Peer-to-Peer Networking
---
A blockchain isn't one program — it's many independent nodes that need to find each other, share blocks and transactions, and agree on which chain is the real one when they briefly disagree.

### Networking Fundamentals for Blockchain
Before a peer-to-peer protocol, understand the plain TCP connections it's built on top of.

**Resources:**
- [Intro to Networking for Blockchain](course:blockchain-in-go#43-intro-to-networking-for-blockchain)
- [TCP Servers and Clients in Go](course:blockchain-in-go#44-tcp-servers-and-clients-in-go)

### Designing a P2P Protocol
Design the actual message format nodes use to talk to each other — what a "here's a new block" message looks like on the wire.

**Resources:**
- [Designing a P2P Protocol](course:blockchain-in-go#45-designing-a-p2p-protocol)
- [Building a P2P Node in Go](course:blockchain-in-go#46-building-a-p2p-node-in-go)

### Practice: P2P Networking
> branches-from: Designing a P2P Protocol

Build the networking task of the standalone project: nodes that connect to each other and exchange blocks over TCP.

**Resources:**
- [Build a Blockchain project](project:build-a-blockchain)

### Peer Discovery & Gossip
New nodes need a way to find existing ones, and the network needs a way to spread a new block to everyone without every node connecting to every other node directly.

**Resources:**
- [Peer Discovery and Handshakes](course:blockchain-in-go#47-peer-discovery-and-handshakes)
- [Gossip Protocols and Broadcasting](course:blockchain-in-go#48-gossip-protocols-and-broadcasting)

### Handling Forks: Longest Chain Rule
Two miners can find a valid block at almost the same time, temporarily splitting the network's view of the chain. The longest-chain rule is how the network converges back to one history.

**Resources:**
- [Blockchain Synchronization](course:blockchain-in-go#49-blockchain-synchronization)
- [Handling Forks: The Longest Chain Rule](course:blockchain-in-go#50-handling-forks-the-longest-chain-rule)

### Attack Basics: Sybil & Eclipse
> optional

A peer-to-peer network has its own attack surface — understand the basics of what a Sybil attack and an eclipse attack are before you deploy anything real.

**Resources:**
- [Sybil and Eclipse Attacks: Basics](course:blockchain-in-go#51-sybil-and-eclipse-attacks-basics)

### Major Project: A Real Multi-Node Network
Run several independent GoChain nodes, have them discover each other, and watch a block mined on one propagate to all of them.

**Resources:**
- [Major Project 2: A Real Multi-Node Network](course:blockchain-in-go#52-major-project-2-a-real-multi-node-network)
