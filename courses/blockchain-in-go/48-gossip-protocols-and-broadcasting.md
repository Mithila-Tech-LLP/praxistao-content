# Chapter 48: Gossip Protocols and Broadcasting

Chapter 46 gave `network.Node` the ability to talk to a handful of peers, and Chapter 47 gave it a way to discover more of them. Neither chapter answered a more basic question: once a node hears about a new transaction, how does that news reach every other honest node on the network, without every node being directly connected to every other node? This chapter answers that with gossip — the same mechanism a rumor uses to spread through a crowded room — and gives `network.Node` its first real broadcasting logic.

## Table of Contents

1. [The Problem With Connecting to Everyone](#1-the-problem-with-connecting-to-everyone)
2. [Gossip Protocols: A Rumor in a Crowd](#2-gossip-protocols-a-rumor-in-a-crowd)
3. [The Missing Piece: Stopping a Rumor From Looping Forever](#3-the-missing-piece-stopping-a-rumor-from-looping-forever)
4. [Extending network.Node With Gossip State](#4-extending-networknode-with-gossip-state)
5. [Implementing HandleTx](#5-implementing-handletx)
6. [Wiring HandleTx Into the Message Dispatch Loop](#6-wiring-handletx-into-the-message-dispatch-loop)
7. [Implementing HandleBlock](#7-implementing-handleblock)
8. [Tracing a Rumor Across Six Nodes](#8-tracing-a-rumor-across-six-nodes)
9. [Limits of This Simple Gossip Scheme](#9-limits-of-this-simple-gossip-scheme)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Problem With Connecting to Everyone

The most obvious way to make sure every node hears about every new transaction is to have every node connect directly to every other node, and send each new transaction to all of them. This is called a **fully connected mesh**, and it falls apart almost immediately at any real scale.

```
 4 nodes, fully connected:            100 nodes, fully connected:
                                       each node needs ~99 connections
   A---B                              total connections ~= 100*99/2 = 4,950
   |\ /|
   | X |          6 connections       At 10,000 nodes:
   |/ \|                              ~50 MILLION connections
   C---D
```

Each node would need thousands of open TCP connections, and every single transaction would need to be sent out thousands of times by its origin node alone. Real peer-to-peer networks — Bitcoin, Ethereum, and GoChain — solve this differently: every node stays connected to only a small number of peers (Chapter 47's peer address book already limits this), and relies on those peers to pass information onward on its behalf.

---

## 2. Gossip Protocols: A Rumor in a Crowd

Think about how a rumor actually spreads through a crowded party. You don't stand on a chair and shout it to all two hundred guests. You tell the two or three people standing near you. Each of them tells the two or three people near *them*. Within a handful of these "hops," nearly everyone at the party has heard the rumor — even though no single person ever spoke to more than a few others.

A **gossip protocol** is exactly this idea, formalized for computer networks: when a node learns something new (a transaction, a block), it forwards that information to its own peers, who forward it to their own peers, and so on. No single node needs to know about — or connect to — the whole network. The information still reaches everyone, just a few hops later than if one node had shouted directly at everyone.

```
Round 0: A learns about tx X (someone submitted it directly to A)

    A       B       C       D       E       F
   [X]

Round 1: A forwards X to its peers, B and D

    A       B       C       D       E       F
   [X]     [X]             [X]

Round 2: B forwards to its peers C and E; D forwards to its peers E and F

    A       B       C       D       E       F
   [X]     [X]     [X]     [X]     [X]     [X]
                                     ^
                              (E hears it twice --
                               from both B and D)

Round 3: every node has X. Total: 2 hops, not 5 direct sends from A.
```

Two things about this diagram matter. First, **A never talked to C, E, or F directly** — it only had to know about B and D. Second, **E received the same transaction twice**, once relayed by B and once relayed by D. That is completely normal in a gossip network with more than one path between nodes — and it is exactly the problem the next section solves.

A new term worth pinning down here: a **hop** is one node-to-node forward. "The message reached everyone within 2 hops" means no honest node needed more than two forwards to hear it — a useful, very small number, considering a real GoChain network might have thousands of nodes, each connected to only a handful of peers.

---

## 3. The Missing Piece: Stopping a Rumor From Looping Forever

If E just blindly forwards every transaction it receives to *its* peers (including back to B and D), you get an infinite loop: B sends to E, E sends back to B, B sends to E again, forever, using more and more bandwidth and CPU for no benefit. Real gossip networks avoid this with a simple rule: **remember what you've already seen, and never forward the same thing twice.**

This is called **duplicate suppression**, and it needs a **seen-set**: a record of message identifiers (in GoChain's case, a transaction's or block's own hash) that a node has already processed. The rule becomes:

1. A new transaction or block arrives.
2. Compute its ID (its hash).
3. Is that ID already in our seen-set? If yes — we've handled this one before, whether we created it locally or heard it from a peer — drop it silently and do nothing further.
4. If no — record the ID in the seen-set, validate the message, and (if valid) forward it to our own peers.

This single check is what turns "shout everything you hear, forever" into "each node forwards each fact exactly once." It is also what makes E's duplicate copy in the diagram above harmless: E processes the *first* copy it receives (from either B or D, whichever arrives first) and silently drops the second.

---

## 4. Extending network.Node With Gossip State

`network.Node`, as built in Chapter 46, already tracks its `peers` and `knownAddrs`. This chapter adds the seen-sets described above as two more unexported fields:

```go
package network

import (
	"sync"

	"github.com/you/gochain/core"
)

// seenSet remembers message IDs (hex-encoded hashes) we've already
// processed, so HandleTx and HandleBlock never process -- or forward --
// the same message twice. A real network runs for months, so this cannot
// grow forever: we cap it and evict the oldest entry once it's full.
type seenSet struct {
	mu    sync.Mutex
	ids   map[string]bool
	order []string // insertion order, so we know what to evict first
	max   int
}

func newSeenSet(max int) *seenSet {
	return &seenSet{
		ids: make(map[string]bool),
		max: max,
	}
}

// seenOrAdd returns true if id was already recorded (meaning: drop this
// message, it's a duplicate). If id is new, it records it and returns false.
func (s *seenSet) seenOrAdd(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	if s.ids[id] {
		return true
	}

	s.ids[id] = true
	s.order = append(s.order, id)

	// Evict the oldest entry once we exceed our cap, so long-running nodes
	// don't accumulate an unbounded map over weeks of uptime.
	if len(s.order) > s.max {
		oldest := s.order[0]
		s.order = s.order[1:]
		delete(s.ids, oldest)
	}

	return false
}

// Node is the peer-to-peer node type from Chapter 46. This chapter adds two
// fields to the ones already there -- everything above the comment is
// unchanged from Chapter 46.
type Node struct {
	Address    string
	Chain      *core.Blockchain
	Mempool    *core.Mempool
	peers      map[string]*Peer
	knownAddrs map[string]bool

	// Added in this chapter: gossip duplicate suppression.
	seenTx    *seenSet
	seenBlock *seenSet
}
```

`seenSet` is a small helper type, not part of the public API — nothing outside `network` ever touches it directly. It wraps a plain `map[string]bool` (for fast "have we seen this?" lookups) with a `sync.Mutex`, because `HandleTx` and `HandleBlock` will be called concurrently from many goroutines, one per open peer connection. The `order` slice and the eviction logic in `seenOrAdd` exist purely so a node that runs for weeks doesn't slowly leak memory into an ever-growing map — once we've remembered `max` message IDs (a reasonable value is a few thousand), the oldest one is forgotten to make room for the newest.

`NewNode` (already defined in Chapter 46) should be updated to initialize these two new fields, the same way it already initializes `peers` and `knownAddrs`:

```go
func NewNode(address string, chain *core.Blockchain, mempool *core.Mempool) *Node {
	return &Node{
		Address:    address,
		Chain:      chain,
		Mempool:    mempool,
		peers:      make(map[string]*Peer),
		knownAddrs: make(map[string]bool),
		seenTx:     newSeenSet(5000),
		seenBlock:  newSeenSet(5000),
	}
}
```

---

## 5. Implementing HandleTx

`HandleTx` is the single entry point for "a transaction has entered my awareness, whether I created it or a peer sent it to me." Every transaction — whether it came from your own wallet CLI (Chapter 36) or from another node over the wire — should flow through this one method, so the dedup-and-forward logic only has to be written once.

```go
// HandleTx processes a transaction this node has just learned about --
// either created locally by a wallet, or received from a peer over the
// wire. It deduplicates using the seen-set, validates the transaction by
// trying to add it to the local mempool, and (only if it is new and valid)
// gossips it onward to this node's own peers.
func (n *Node) HandleTx(tx *core.Transaction) {
	id := tx.IDHex()

	// Step 1: dedup. If we've already handled this exact transaction --
	// whether we created it ourselves or heard it from a different peer --
	// there is nothing left to do. This single check is what stops the
	// rumor from looping forever between us and our peers.
	if n.seenTx.seenOrAdd(id) {
		return
	}

	// Step 2: validate. Mempool.Add (Volume 5) checks the signature, that
	// referenced UTXOs actually exist and aren't already spent by another
	// pending transaction, and that inputs cover outputs plus fee. A
	// transaction that fails any of this is simply dropped -- we never
	// forward garbage to our peers just because someone sent it to us.
	if err := n.Mempool.Add(tx); err != nil {
		return
	}

	// Step 3: gossip onward. Broadcast (Chapter 46) writes this message to
	// every currently connected peer. Some of those peers will already have
	// this transaction (they may be the very peer who sent it to us, or
	// another path may have beaten us to them) -- that's fine, their own
	// seen-set will catch and drop the duplicate on arrival.
	n.Broadcast(MsgTx, tx.Serialize())
}
```

Notice what `HandleTx` deliberately does *not* do: it does not try to track which peer sent it this transaction, so it can avoid sending it right back to them. That optimization is real and worth knowing about (production networks like Bitcoin's do it, to cut wasted bandwidth roughly in half), but it is not required for *correctness* — the seen-set alone guarantees the rumor stops spreading once every node has processed it once. We call this out explicitly in Section 9 and leave the optimization as an exercise, because the simpler version is easier to reason about and get right first.

---

## 6. Wiring HandleTx Into the Message Dispatch Loop

Chapter 46 built a per-connection read loop that decodes each incoming `Envelope` and routes it to a handler based on `Type`. `HandleTx` needs one new case added to that existing `switch`:

```go
import (
	"bytes"
	"encoding/gob"
	"log"
)

// handleEnvelope is the per-connection dispatch function from Chapter 46.
// This chapter adds the MsgTx and MsgBlock cases; everything else was
// already there.
func (n *Node) handleEnvelope(env Envelope, payload []byte) {
	switch env.Type {

	case MsgTx:
		var tx core.Transaction
		dec := gob.NewDecoder(bytes.NewReader(payload))
		if err := dec.Decode(&tx); err != nil {
			log.Printf("network: bad tx payload: %v", err)
			return
		}
		n.HandleTx(&tx)

	case MsgBlock:
		var b core.Block
		dec := gob.NewDecoder(bytes.NewReader(payload))
		if err := dec.Decode(&b); err != nil {
			log.Printf("network: bad block payload: %v", err)
			return
		}
		n.HandleBlock(&b)

	// MsgVersion, MsgGetBlocks, MsgInv, MsgGetData, MsgAddr: handled
	// elsewhere (Chapters 45-47, and MsgGetBlocks/MsgGetData get their real
	// implementations in Chapter 49).
	default:
		log.Printf("network: unhandled message type %d", env.Type)
	}
}
```

This function decodes each message's raw payload bytes back into a `core.Transaction` or `core.Block` using `encoding/gob` (the serialization format this course settled on back in Chapter 07), then simply calls `n.HandleTx` or `n.HandleBlock`. All of the interesting logic — dedup, validation, forwarding — lives inside those two methods, not in the dispatch loop itself.

---

## 7. Implementing HandleBlock

`HandleBlock` follows the identical shape as `HandleTx`: dedup, validate, forward. The one difference is what "validate" means for a block instead of a transaction.

```go
// HandleBlock processes a block this node has just learned about, either
// mined locally or received from a peer. Like HandleTx, it deduplicates,
// validates, and gossips onward -- but "validate" here means re-running the
// full block validation from Chapter 19 and Chapter 25 (hash correctness,
// proof-of-work, linkage to our current chain).
func (n *Node) HandleBlock(b *core.Block) {
	id := b.HashHex()

	if n.seenBlock.seenOrAdd(id) {
		return
	}

	// ValidateBlock (Volume 3 and Volume 4) checks the block's own hash,
	// its proof-of-work, and that it correctly links to the block before
	// it. A block that fails this is simply dropped -- accepting an
	// invalid block just because a peer sent it to us would let any single
	// dishonest peer corrupt our chain.
	if err := n.Chain.ValidateBlock(b); err != nil {
		log.Printf("network: rejecting invalid block %s: %v", id, err)
		return
	}

	if err := n.Chain.AddBlock(b); err != nil {
		// This is not necessarily an attack -- it commonly just means the
		// block doesn't connect to our current chain tip, because we're
		// missing some blocks in between (Chapter 49 fixes this with
		// SyncWithPeer) or because it's part of a competing fork (Chapter
		// 50 fixes this with ChainWork and ReplaceChain). For now we simply
		// decline to add it and move on.
		log.Printf("network: could not add block %s to chain: %v", id, err)
		return
	}

	log.Printf("network: accepted block %d (%s), forwarding to %d peers", b.Height, id, len(n.peers))
	n.Broadcast(MsgBlock, b.Serialize())
}
```

This version of `HandleBlock` is intentionally the simplest one that works when the network is healthy and every node's chain is already in sync. Chapters 49 and 50 will come back and make it smarter: Chapter 49 makes sure a node that's behind can catch up instead of just discarding blocks it can't yet connect, and Chapter 50 teaches it what to do when two different, individually valid blocks show up claiming the same spot in the chain.

---

## 8. Tracing a Rumor Across Six Nodes

Let's walk through the six-node network from Section 2 again, but now with the actual dedup logic from this chapter, and imagine E is connected to both B and D (a realistic shape, since Chapter 47's peer exchange tends to create a few redundant paths like this on purpose — more on why in Chapter 51).

```
Topology (lines are open peer connections):

    A --- B --- C
    |     |
    D --- E --- F

A submits a new transaction, tx X, straight from a wallet.
```

**Hop 0 — origin.** A wallet calls `chain.Send(...)`, producing `tx X`. A's own code calls `n.HandleTx(X)` directly (not over the network — it's local). A's seen-set does not contain X's ID yet, so it is recorded, added to A's own mempool, and broadcast to A's peers: B and D.

```
seenTx sets after hop 0:     A:{X}   B:{}    C:{}    D:{}    E:{}    F:{}
```

**Hop 1 — first relay.** B and D each receive `MsgTx` over the wire, decode it, and call `n.HandleTx(X)`. Neither has seen X before, so both add it to their own mempool and broadcast onward: B forwards to A and C; D forwards to A and E.

```
seenTx sets after hop 1:     A:{X}   B:{X}   C:{}    D:{X}   E:{}    F:{}
```

**Hop 2 — the interesting part.** A receives X back from both B and D — but A's seen-set already has X, so both arrivals are silently dropped; A does not re-broadcast. C receives X from B, adds it, and forwards to... B, its only peer (dropped there, already seen). E receives X from D, adds it, and forwards to its peers B and F.

```
seenTx sets after hop 2:     A:{X}   B:{X}   C:{X}   D:{X}   E:{X}   F:{}
```

**Hop 3 — the last node.** F receives X from E, adds it, forwards to E (dropped, already seen).

```
seenTx sets after hop 3:     A:{X}   B:{X}   C:{X}   D:{X}   E:{X}   F:{X}

Every node now has X. No further forwarding happens -- every peer that
receives X from this point on already has it in its seen-set and drops it
immediately.
```

Two things to notice: the loop A → B/D → A was harmless because of the seen-set, and the whole network converged in three hops despite no node connecting to more than three peers. This is the payoff of gossip: linear cost per node, network-wide convergence in a small, roughly logarithmic number of hops.

---

## 9. Limits of This Simple Gossip Scheme

This chapter's version is deliberately the simplest one that correctly and safely spreads information without looping forever. A few honest limitations worth naming, some of which later chapters and exercises address:

- **No origin-skip optimization.** As mentioned in Section 5, we always broadcast to *all* peers, including possibly the one who just sent it to us. This wastes some bandwidth (that peer's seen-set will just drop the duplicate) but costs nothing in correctness.
- **No hop-count limit.** Some gossip protocols attach a "time to live" counter that decrements each hop, so a message eventually stops propagating even in pathological topologies. Our seen-set already prevents infinite loops, so a hop counter is a refinement, not a requirement, here.
- **Eventual, not instant, consistency.** A node three hops away from the origin sees a new transaction slightly later than a node one hop away. GoChain (like Bitcoin and Ethereum) is an **eventually consistent** system: every honest, connected node converges on the same information given enough time, but "given enough time" might be a few seconds on a real network, not zero.
- **A malicious peer can still lie.** Gossip spreads information faithfully, but it does not by itself protect against a peer that refuses to forward things, or that floods you with junk. `Mempool.Add` and `ValidateBlock` protect you from *invalid* data; Chapter 51 covers protecting yourself from a peer that behaves *maliciously but validly* (for example, by surrounding you with only its own connections).

---

## Summary

- Fully connecting every node to every other node does not scale; gossip protocols spread information hop by hop through a small number of peer connections instead, the same way a rumor spreads through a crowd.
- A **seen-set** — a record of message IDs a node has already processed — is what stops a gossiped message from looping between peers forever; without it, "rumor spreading" degenerates into an infinite echo.
- `network.Node` gained two new unexported fields this chapter, `seenTx` and `seenBlock`, each a capped, mutex-protected `seenSet` that evicts its oldest entry once full so long-running nodes don't leak memory.
- `HandleTx(tx *core.Transaction)` is the single entry point for any transaction, local or remote: dedup, validate via `Mempool.Add`, then `Broadcast(MsgTx, ...)` onward.
- `HandleBlock(b *core.Block)` follows the identical shape, validating with `Chain.ValidateBlock` and `Chain.AddBlock` instead of `Mempool.Add`.
- A worked six-node trace showed a transaction reaching every node in 3 hops, with duplicate arrivals silently dropped by the seen-set rather than re-forwarded.
- This simple scheme has known, acceptable limitations (no origin-skip, no hop-count limit, eventual rather than instant consistency) that later chapters build on, particularly around what happens when a peer behaves maliciously (Chapter 51).

---

## Exercises

### Easy

1. **Trace by hand** a gossip broadcast on a different topology: a straight line of five nodes, `A - B - C - D - E` (A only connects to B, B connects to A and C, and so on). If A originates a transaction, how many hops does it take to reach E? Write out each hop's seen-set contents the way Section 8 did.

2. **Modify `seenSet`** to expose a `Len() int` method returning the number of currently remembered IDs, and write a short Go test that adds more than `max` entries and asserts the set never grows past `max`.

3. **Explain in your own words** why `HandleTx` checks the seen-set *before* validating the transaction with `Mempool.Add`, rather than after. What would go wrong (think about wasted CPU, not correctness) if the order were reversed on a busy network relaying the same popular transaction from five peers at once?

### Medium

4. **Implement the origin-skip optimization** mentioned in Section 9: change `Broadcast` (or add a new method, `BroadcastExcept(msgType MessageType, payload []byte, except string)`) so `HandleTx` and `HandleBlock` can skip re-sending a message back to the peer address it just arrived from. You'll need to know which peer a given incoming message came from — thread that information through from `handleEnvelope`.

5. **Add a hop-count field** to the wire format for `MsgTx` (you'll need a small wrapper struct around `core.Transaction` with an extra `Hops int` field). Have `HandleTx` refuse to forward a transaction once `Hops` exceeds some maximum (say, 20), and write a test proving a message with `Hops` already at the maximum is dropped without being re-broadcast.

6. **Simulate the six-node network from Section 8 in a single Go test**, without any real networking: represent each node as a `*Node` with an in-memory "peer" abstraction (a Go channel or a direct method call standing in for a TCP connection) instead of real sockets, originate a transaction at A, and assert every node's mempool contains it after gossip settles.

### Hard

7. **Measure gossip convergence time** empirically: build a small in-process simulation of 50 nodes arranged in a random topology (each node connected to a random 4-6 others), originate one transaction at a random node, and count how many hops it takes for the last node to receive it. Run this 100 times with different random topologies and report the average and worst-case hop count. What does this tell you about how gossip protocols scale as the network grows?

8. **Design (and if you like, implement) a "gossip storm" defense.** A malicious node could construct thousands of unique, individually-valid-looking transactions and originate all of them at once, forcing every node in the network to spend CPU deduplicating and forwarding all of them. Propose at least two concrete defenses GoChain's `HandleTx` could add (think about rate limiting, minimum fees, or per-peer message quotas), and explain the trade-off each one makes between blocking abuse and blocking legitimate heavy use.

9. **Compare GoChain's flood-based gossip (this chapter) with "epidemic" or "gossip-push-pull" protocols** used in some real distributed systems (for example, Amazon's Dynamo or Cassandra's anti-entropy gossip). Research how those systems periodically exchange summaries of what each node has seen, rather than eagerly forwarding every message the instant it arrives, and write a 250-350 word comparison of the two approaches' trade-offs for a blockchain specifically, where near-real-time propagation of new blocks matters much more than it does for, say, replicating a key-value store's data.
