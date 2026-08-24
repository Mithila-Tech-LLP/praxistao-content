# Chapter 47: Peer Discovery and Handshakes

Chapter 46 gave `network.Node` the ability to `Listen` and `Dial`, but only if you already knew the exact address of another node to connect to. That is fine for two nodes you started yourself in two terminals — it falls apart the moment a stranger wants to join the GoChain network without you personally handing them an address list. This chapter fixes that with two complementary mechanisms: a small list of **seed nodes** to bootstrap from, and a **peer-exchange** protocol (`MsgAddr`, defined back in Chapter 45 but never implemented) that lets connectivity spread from those seeds the same way gossip spreads a transaction.

## Table of Contents

1. [The Problem: No Central Directory](#1-the-problem-no-central-directory)
2. [Seed Nodes: A Small, Trusted Starting Point](#2-seed-nodes-a-small-trusted-starting-point)
3. [Bootstrapping From Seeds](#3-bootstrapping-from-seeds)
4. [MsgAddr, Revisited](#4-msgaddr-revisited)
5. [Implementing handleAddr for Real](#5-implementing-handleaddr-for-real)
6. [Choosing Which New Addresses to Dial](#6-choosing-which-new-addresses-to-dial)
7. [Sharing What We Know: Periodic Peer Exchange](#7-sharing-what-we-know-periodic-peer-exchange)
8. [Guarding Against an Unbounded Peer List](#8-guarding-against-an-unbounded-peer-list)
9. [Watching a Network Grow: One Seed to Five Nodes](#9-watching-a-network-grow-one-seed-to-five-nodes)
10. [Running the Five-Node Demo](#10-running-the-five-node-demo)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Problem: No Central Directory

Imagine moving to a brand-new city with no phone book, no internet search, and no friend to introduce you around. How would you ever meet anyone? In practice, you'd probably start with one or two well-known public places — a town square, a community board — where you know *some* people tend to gather, and from there, ask the people you meet there who else they know.

A peer-to-peer network has exactly this problem, deliberately. There is no server anywhere that maintains "the official list of every GoChain node" — that would just recreate the central point of control and failure that decentralization is trying to avoid in the first place. A brand-new node, running for the very first time, knows about exactly zero other nodes. Chapter 46's `Dial` needs an address to connect to; where does that first address come from, if not from you typing it in by hand?

---

## 2. Seed Nodes: A Small, Trusted Starting Point

The practical answer real networks use — and the one GoChain uses too — is a small, fixed list of **seed nodes**: addresses of a handful of long-running, generally reliable nodes that a new node tries first, purely to get its very first connection or two. A seed node is not special in any technical sense — it is an ordinary `network.Node`, running the exact same code as everyone else. The only thing that makes it a "seed" is that its address is well-known in advance, the way a town square's location is well-known even though the square itself has no special powers.

```
   BEFORE PEER DISCOVERY                AFTER BOOTSTRAPPING FROM SEEDS
   ------------------------             --------------------------------

   New Node                             New Node
   (knows nobody)                       (connected to Seed A)
        ?                                     |
        ?                                     |
        ?          <-- no way to        Seed A ------ (Seed A's own
        ?              start                            peers, learned
                                                          via Section 5)
```

This chapter's seed list is deliberately simple: a hardcoded Go slice, `DefaultSeedNodes`, with the option to override it at startup (for instance, when running a private testnet whose seeds aren't the public ones). A real, large-scale deployment might instead resolve seed addresses from DNS (a technique called a "DNS seed," which Bitcoin uses) or from a config file shipped alongside the node's binary — Chapter 87's Docker Compose testnet, later in this course, uses a config-file variant of exactly this idea. The mechanism for finding the *first* address differs; what every one of these approaches shares is the core idea: start from a small, trusted list, and grow from there.

```go
package network

// DefaultSeedNodes is a small, hardcoded list of addresses a brand-new
// GoChain node can try first, before it has learned about any other node
// on the network from any other source. Nothing about these addresses is
// technically special -- they are ordinary nodes that happen to be
// well-known and (ideally) reliably running, the way a town square is
// well-known without having any special authority itself.
var DefaultSeedNodes = []string{
	"127.0.0.1:3000",
	"127.0.0.1:3001",
	"127.0.0.1:3002",
}
```

---

## 3. Bootstrapping From Seeds

With a seed list defined, `Bootstrap` is the method a node calls once, at startup, to try connecting to each one:

```go
import "log"

// Bootstrap dials every address in seeds that isn't this node's own
// address, logging (but not failing on) any seed that doesn't answer --
// a single retired or temporarily offline seed should never stop a fresh
// node from joining through whichever OTHER seeds do respond. Reaching
// even one live seed is enough: Section 5's peer exchange takes over from
// there and grows the connection from that single point.
func (n *Node) Bootstrap(seeds []string) {
	for _, addr := range seeds {
		if addr == n.Address {
			continue // never dial ourselves, even if we happen to be a seed
		}
		if err := n.Dial(addr); err != nil {
			log.Printf("[%s] could not reach seed %s: %v", n.Address, addr, err)
			continue
		}
		log.Printf("[%s] connected to seed %s", n.Address, addr)
	}
}
```

`Bootstrap` deliberately keeps going after a failed `Dial` rather than returning early — Go's `continue` here is doing real work: one unreachable seed (maybe it's down for maintenance, maybe its operator retired it) is a routine, expected event on a real network, not a reason to give up on joining entirely. Every successful `Dial` already triggers Chapter 46's version handshake automatically, so by the time `Bootstrap` returns, this node has exchanged a `MsgVersion` with every seed that was actually reachable.

---

## 4. MsgAddr, Revisited

Chapter 45 already defined the wire format this chapter needs — `MsgAddr` as a `MessageType` value, and `AddrPayload` as its payload struct:

```go
// AddrPayload is the body of a MsgAddr message -- peer addresses being
// shared. Defined back in Chapter 45; this chapter is the first to give
// it a real handler and a real reason to be sent.
type AddrPayload struct {
	Addresses []string // e.g. []string{"127.0.0.1:3002", "127.0.0.1:3003"}
}
```

Chapter 46's `handleAddr` was an honest stub: it decoded the payload correctly, proving the plumbing worked, then logged a note that real peer-exchange logic would arrive here. That note was a promise this chapter now keeps.

---

## 5. Implementing handleAddr for Real

When a `MsgAddr` message arrives, it means a peer is telling us: "here are some addresses I know about that you might not." The real `handleAddr` needs to record any addresses we didn't already know, and decide whether to act on them:

```go
// handleAddr replaces Chapter 46's stub. It records any newly learned
// addresses in knownAddrs (the set Chapter 46 already reserved space for
// on Node, specifically anticipating this chapter), then hands anything
// genuinely new off to dialSome (Section 6) to decide whether to connect.
func (n *Node) handleAddr(payload []byte) {
	var p AddrPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad addr payload: %v", n.Address, err)
		return
	}

	var fresh []string
	n.mu.Lock()
	for _, addr := range p.Addresses {
		if addr == n.Address || n.knownAddrs[addr] {
			continue // it's us, or we've already heard of it -- nothing new
		}
		n.knownAddrs[addr] = true
		fresh = append(fresh, addr)
	}
	n.mu.Unlock()

	if len(fresh) == 0 {
		return // every address in this message was already familiar
	}

	log.Printf("[%s] learned %d new address(es) via peer exchange: %v", n.Address, len(fresh), fresh)
	n.dialSome(fresh)
}
```

Notice the shape here is deliberately similar to the seen-set idea Chapter 48 will formalize for gossiped transactions: check whether we already know something, and only act on what's actually new. `knownAddrs` (already part of the `Node` struct's contract, and initialized empty by `NewNode`) is playing the same role for addresses that a seen-set will later play for message IDs — a simple, mutex-guarded record of "have I already dealt with this?"

---

## 6. Choosing Which New Addresses to Dial

Learning about an address is not the same as connecting to it — a node that immediately dialed every single address it ever heard about would quickly accumulate an unbounded, unmanageable number of open connections. `dialSome` decides which of the freshly learned addresses (if any) are actually worth connecting to right now:

```go
import "math/rand"

// dialSome tries to connect to a handful of newly learned addresses,
// stopping the moment we hit maxPeers. It shuffles the candidates first,
// so a node doesn't always prefer whichever address happened to arrive
// first in somebody else's list -- a small first step toward the more
// deliberate diverse peer selection Chapter 51 builds directly on top of
// this function.
func (n *Node) dialSome(candidates []string) {
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, addr := range candidates {
		n.mu.RLock()
		full := len(n.peers) >= n.maxPeers
		_, alreadyConnected := n.peers[addr]
		n.mu.RUnlock()

		if full {
			log.Printf("[%s] at max peers (%d) -- not dialing further addresses this round", n.Address, n.maxPeers)
			return
		}
		if alreadyConnected {
			continue // we know the address, but we're already talking to them
		}
		if err := n.Dial(addr); err != nil {
			log.Printf("[%s] could not dial learned address %s: %v", n.Address, addr, err)
			continue
		}
		log.Printf("[%s] dialed newly learned peer %s", n.Address, addr)
	}
}
```

`maxPeers` is a new unexported field on `Node`, set once inside `NewNode` (whose signature stays exactly as the contract specifies — we only add to the struct literal it already builds, the same pattern Chapter 48 used to add `seenTx`/`seenBlock`):

```go
const defaultMaxPeers = 8 // deliberately small -- Chapter 51 explains why "small" matters for security, not only for resource limits

func NewNode(address string, chain *core.Blockchain, mempool *core.Mempool) *Node {
	return &Node{
		Address:    address,
		Chain:      chain,
		Mempool:    mempool,
		peers:      make(map[string]*Peer),
		knownAddrs: make(map[string]bool),
		maxPeers:   defaultMaxPeers,
	}
}
```

Eight might sound low for a network that could eventually have thousands of participants, but it does not need to be high: as Chapter 48's gossip protocol will show directly, a node does not need to be connected to everyone to eventually hear from everyone — it only needs a small number of honest, reasonably diverse connections for information to reach it within a few hops.

---

## 7. Sharing What We Know: Periodic Peer Exchange

`handleAddr` and `dialSome` cover receiving addresses. The other half of peer exchange is *sending* them — proactively telling our own peers about every address we know, on a regular schedule, so the whole network's address books converge over time:

```go
import (
	"bytes"
	"encoding/gob"
	"time"
)

// StartPeerExchange begins periodically sharing this node's full address
// book with every currently connected peer. Call it once, right after
// Listen, so a node keeps spreading what it knows for as long as it runs
// -- not just once at startup.
func (n *Node) StartPeerExchange(interval time.Duration) {
	go func() {
		ticker := time.NewTicker(interval)
		defer ticker.Stop()
		for range ticker.C {
			n.shareKnownAddresses()
		}
	}()
}

// shareKnownAddresses gob-encodes this node's entire knownAddrs set into
// one AddrPayload and broadcasts it as a single MsgAddr to every
// connected peer, using Broadcast exactly as Chapter 46 defined it.
func (n *Node) shareKnownAddresses() {
	n.mu.RLock()
	addrs := make([]string, 0, len(n.knownAddrs))
	for a := range n.knownAddrs {
		addrs = append(addrs, a)
	}
	n.mu.RUnlock()

	if len(addrs) == 0 {
		return // nothing worth sharing yet
	}

	payload := AddrPayload{Addresses: addrs}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(payload); err != nil {
		log.Printf("[%s] failed encoding address list: %v", n.Address, err)
		return
	}
	n.Broadcast(MsgAddr, buf.Bytes())
}
```

`Broadcast`'s contract (Chapter 46, Section 9) expects an already gob-encoded payload, not a raw struct, which is exactly why `shareKnownAddresses` encodes the `AddrPayload` by hand with `gob.NewEncoder` before calling it — the same pattern later chapters use for `MsgTx` and `MsgBlock` payloads. Running this on a `time.Ticker` in its own goroutine means peer exchange keeps happening for the entire lifetime of the node, not just once at startup — which matters because `knownAddrs` keeps growing as more `MsgAddr` messages arrive, and each round of sharing might carry addresses that didn't exist in the network yet the last time this node shared its list.

---

## 8. Guarding Against an Unbounded Peer List

Two small, easy-to-miss details keep this scheme well-behaved as a network grows into the hundreds or thousands of nodes. First, `knownAddrs` is allowed to grow much larger than `maxPeers` — a node can *know about* far more addresses than it is actively *connected to*, which is exactly the distinction between an address book and a phone call. Second, `dialSome`'s `full` check inside its loop, not just before it, means a node that starts a dialing round with room for three more peers stops the instant it fills those three slots, even partway through a longer candidate list — it never overshoots `maxPeers` by racing ahead of its own connection count.

```
   knownAddrs (address book):   50 addresses, learned over time
   peers (active connections):  8 addresses, currently connected

   Every connected peer's address is also in knownAddrs, but most
   addresses in knownAddrs have no live connection right now -- they're
   just remembered for later, the way you might remember a phone number
   without being on a call with that person at this exact moment.
```

This distinction is not just tidy bookkeeping — it is the seed for Chapter 51's entire discussion of eclipse attacks. A node that dialed literally every address it ever heard about would have no room left to make deliberate, careful choices about *which* peers to trust with its view of the network; keeping `maxPeers` small and `knownAddrs` merely informative is what leaves room for the smarter selection logic Chapter 51 builds directly on top of `dialSome`.

---

## 9. Watching a Network Grow: One Seed to Five Nodes

Let's trace an entire small network coming into existence from a single seed, over a few rounds of bootstrapping and peer exchange.

```
ROUND 0 -- Seed S starts. It knows nobody. It just Listens.

    [S]

ROUND 1 -- Node 2 starts, Bootstraps from S. Node 3 starts moments later,
also Bootstraps from S. Both dial S directly; S now has two peers.

    [S]---[2]
     |
    [3]

    knownAddrs:  S:{2,3}   2:{S}   3:{S}

ROUND 2 -- S's periodic peer exchange (Section 7) fires. S broadcasts its
address book -- which now includes 2 and 3 -- to every one of its peers.
Node 2 receives S's list, learns about 3 (new!), and dialSome connects to
it. Node 3 receives S's list, learns about 2, but 2 has already dialed it
first -- dialSome sees "already connected" and does nothing further.

    [S]---[2]
     |     |
    [3]----+

    knownAddrs:  S:{2,3}   2:{S,3}   3:{S,2}

ROUND 3 -- Node 4 starts, Bootstraps from S (dials S directly, exactly as
2 and 3 did). Node 5 starts, and its Bootstrap call happens to reach S
too, dialing it directly.

    [S]---[2]
   / | \   |
  4  3  5  |
     |     |
     +-----+

    knownAddrs now includes 4 and 5 at S, and (after the next peer
    exchange round) will spread to 2 and 3 as well, exactly as in Round 2.

ROUND 4 -- S's next peer exchange tick fires. Every one of S's five peers
receives S's full address book. Node 4 and Node 5 each learn about 2 and
3 (and each other) and dialSome connects a few of those new addresses,
bringing the network from a hub-and-spoke shape (everyone connected only
to S) toward a more connected mesh -- with S no longer the only node
anyone can reach.

    Final shape after a few more exchange rounds (illustrative, exact
    edges depend on maxPeers and dialSome's random shuffle):

      [2]---[S]---[4]
       |     |      \
      [3]---[5]------+
```

Two things about this trace matter beyond the specific edges drawn. First, notice that no single node ever had to know about the *whole* network in advance — every connection after Round 1 was discovered, not configured. Second, notice that the network did not become fully connected (Chapter 48's opening section explains directly why that would be undesirable at scale even if it happened) — it became a small, connected mesh with some redundancy, which is exactly the shape gossip protocols are designed to work well on.

---

## 10. Running the Five-Node Demo

We can extend Chapter 46's `cmd/gochain-node/main.go` with a `-seeds` flag and a call to `Bootstrap` and `StartPeerExchange`, keeping everything else from that chapter unchanged:

```go
package main

import (
	"flag"
	"log"
	"strings"
	"time"

	"github.com/you/gochain/core"
	"github.com/you/gochain/network"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:3000", "this node's listen address")
	seeds := flag.String("seeds", "", "comma-separated seed addresses to bootstrap from (optional)")
	flag.Parse()

	chain := core.NewBlockchain()
	mempool := core.NewMempool()

	node := network.NewNode(*addr, chain, mempool)
	if err := node.Listen(); err != nil {
		log.Fatal(err)
	}

	// Peer exchange every 10 seconds -- frequent enough for this demo to
	// converge in seconds rather than minutes, generous enough not to
	// spam a real network with redundant address lists.
	node.StartPeerExchange(10 * time.Second)

	seedList := network.DefaultSeedNodes
	if *seeds != "" {
		seedList = strings.Split(*seeds, ",")
	}
	node.Bootstrap(seedList)

	select {} // block forever
}
```

Run five terminals: the first starts with no `-seeds` flag at all (it is itself the only seed, using the default list, none of which it can reach yet, which is fine), and each subsequent one points at the first:

```bash
# terminal 1 -- the first node in the network
go run ./cmd/gochain-node -addr 127.0.0.1:3000

# terminal 2
go run ./cmd/gochain-node -addr 127.0.0.1:3001 -seeds 127.0.0.1:3000

# terminal 3
go run ./cmd/gochain-node -addr 127.0.0.1:3002 -seeds 127.0.0.1:3000

# terminal 4
go run ./cmd/gochain-node -addr 127.0.0.1:3003 -seeds 127.0.0.1:3000

# terminal 5
go run ./cmd/gochain-node -addr 127.0.0.1:3004 -seeds 127.0.0.1:3000
```

Within a few peer-exchange intervals (about 10-20 seconds with the settings above), terminal 1's logs will show something like:

```
2026/07/31 11:00:00 [127.0.0.1:3000] listening for peers
2026/07/31 11:00:04 [127.0.0.1:3000] peer 127.0.0.1:3001 is caught up (height 0)
2026/07/31 11:00:07 [127.0.0.1:3000] peer 127.0.0.1:3002 is caught up (height 0)
2026/07/31 11:00:11 [127.0.0.1:3000] peer 127.0.0.1:3003 is caught up (height 0)
2026/07/31 11:00:14 [127.0.0.1:3000] peer 127.0.0.1:3004 is caught up (height 0)
2026/07/31 11:00:14 [127.0.0.1:3000] learned 0 new address(es) via peer exchange: []
```

...and a few seconds later, terminal 2's logs will show it discovering and connecting to nodes it never dialed directly:

```
2026/07/31 11:00:24 [127.0.0.1:3001] learned 3 new address(es) via peer exchange: [127.0.0.1:3002 127.0.0.1:3003 127.0.0.1:3004]
2026/07/31 11:00:24 [127.0.0.1:3001] dialed newly learned peer 127.0.0.1:3002
2026/07/31 11:00:24 [127.0.0.1:3001] dialed newly learned peer 127.0.0.1:3003
2026/07/31 11:00:24 [127.0.0.1:3001] dialed newly learned peer 127.0.0.1:3004
```

Node 2 never typed in node 3's, 4's, or 5's addresses anywhere — it discovered every one of them purely through node 1's peer exchange, exactly as the diagram in Section 9 predicted.

---

## Summary

- Peer-to-peer networks have no central directory; a brand-new node needs some other way to make its very first connection, which this chapter solves with **seed nodes** — a small, well-known list of addresses to try first.
- `Bootstrap(seeds []string)` dials every seed address in turn, logging (not failing on) any that don't answer, since one unreachable seed should never block joining through the others.
- `MsgAddr`/`AddrPayload`, defined back in Chapter 45 and stubbed in Chapter 46, get their first real implementation here: `handleAddr` records genuinely new addresses into `knownAddrs` and hands them to `dialSome`.
- `dialSome` shuffles candidate addresses and stops the moment `maxPeers` (a small, deliberately conservative default of 8) is reached, so a node's connection count never grows unbounded just because its address book does.
- `StartPeerExchange` runs a background goroutine that periodically broadcasts this node's entire `knownAddrs` set as a `MsgAddr`, which is exactly what lets connectivity spread outward from a handful of seeds without any node needing to know the whole network's shape in advance.
- `knownAddrs` (an address book) and `peers` (active connections) are deliberately different sizes — knowing about an address is not the same as being connected to it, and that gap is exactly what leaves room for the more careful peer-selection logic Chapter 51 builds next.
- A five-node demo showed exactly this: nodes 3, 4, and 5 were never dialed directly by node 2, yet node 2 connected to all of them purely through peer exchange relayed by node 1.

---

## Exercises

### Easy

1. **Run the five-node demo from Section 10 on your own machine.** Capture logs from at least three of the five terminals and annotate, in your own words, which log lines correspond to `Bootstrap`, which correspond to `handleAddr` learning something new, and which correspond to `dialSome` acting on that new knowledge.

2. **Reduce `defaultMaxPeers` to 2** and rerun the five-node demo. Confirm (from the logs) that at least one node's `dialSome` logs the "at max peers" message and stops partway through a batch of candidates, and explain in one or two sentences why a smaller `maxPeers` makes the network take longer to become well-connected.

3. **Add a `KnownAddresses() []string` method to `Node`** that returns a safe, read-locked snapshot of `knownAddrs`, and use it to print each node's full address book (not just its active peers) every 30 seconds, so you can watch `knownAddrs` grow independently of `peers`.

### Medium

4. **Change `Bootstrap` to try seeds in parallel** (using a goroutine per seed and a `sync.WaitGroup`) instead of one at a time, and measure whether this meaningfully speeds up joining a network with several seeds, some of which are slow or unreachable. Explain what could go wrong if two goroutines both call `registerPeer` for the same address at nearly the same moment, and whether the existing `n.mu` protects against it.

5. **Add a `RemoveKnownAddress(addr string)` method** and wire it into `handleConnection`'s cleanup path (alongside `removePeerByConn` from Chapter 46) so that an address whose connection has failed repeatedly (track a small failure count per address) is eventually forgotten rather than retried forever by future `dialSome` calls.

6. **Write an integration test** with four in-process `Node`s: start node A with no seeds, start B, C, and D each bootstrapping only from A, call `shareKnownAddresses` manually (without waiting for the ticker) on A, B, C, and D in some order, and assert that every node's `knownAddrs` eventually contains all four addresses.

### Hard

7. **Implement address expiry.** Real address books (Bitcoin's included) eventually forget addresses that have not answered in a long time, the same way old, stale phone numbers become useless. Add a `lastSeen` timestamp per known address, updated whenever that address answers a `Dial` or sends any message, and a background goroutine that removes addresses whose `lastSeen` is older than some threshold (say, one hour) from `knownAddrs`.

8. **Simulate a seed node going offline mid-bootstrap.** Start a seed node, have three other nodes begin `Bootstrap`-ing from it, then kill the seed process partway through. Confirm the three nodes still discover each other eventually (hint: think about what needs to happen for at least one of them to have already learned about the others before the seed died) and explain what would happen instead if the seed died before any peer exchange had occurred at all.

9. **Design and implement a basic reputation score per known address**: increment a counter each time a `Dial` to that address succeeds and the resulting peer stays connected for at least some minimum duration, decrement it on failed or very short-lived connections, and change `dialSome`'s shuffle into a weighted selection that prefers higher-reputation addresses when `maxPeers` forces a choice. Discuss, in a short comment, why this is a genuinely double-edged mitigation: it improves connection quality for an honest network, but could an attacker try to game a reputation system like this one, and how?
