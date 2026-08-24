# Chapter 46: Building a P2P Node in Go

With the protocol fully designed on paper in Chapter 45, this chapter builds the real thing: `network.Node`, a type that listens for incoming connections, dials outgoing ones, and routes every message it receives to the right handler based on its type. By the end of this chapter you will run two GoChain nodes on your own machine, watch them dial each other, and read their first real handshake exchange straight from the terminal.

## Table of Contents

1. [From Protocol to a Working Node](#1-from-protocol-to-a-working-node)
2. [The Peer Type](#2-the-peer-type)
3. [The Node Type](#3-the-node-type)
4. [Listening for Incoming Connections](#4-listening-for-incoming-connections)
5. [Dialing Out to a Peer](#5-dialing-out-to-a-peer)
6. [Routing Messages to Handlers](#6-routing-messages-to-handlers)
7. [Handling the Version Handshake](#7-handling-the-version-handshake)
8. [The Remaining Handlers, For Now](#8-the-remaining-handlers-for-now)
9. [Broadcasting to Every Peer](#9-broadcasting-to-every-peer)
10. [Running Two Nodes and Watching Them Meet](#10-running-two-nodes-and-watching-them-meet)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. From Protocol to a Working Node

Everything in this chapter assembles pieces you already have: `net.Listen`/`net.Dial` and one-goroutine-per-connection from Chapter 44, and `MessageType`, `Envelope`, the payload structs, and `EncodeMessage`/`ReadMessage`/`DecodePayload` from Chapter 45. `network.Node` is the type that owns a node's listening socket, its set of active peer connections, and the logic that decides what to do with each kind of message once it arrives.

Think of `Node` as the receptionist and switchboard of a small office: it answers the door (`Listen`) for anyone who shows up, it can also go knock on other offices' doors itself (`Dial`), and once a conversation is underway, it reads what each visitor says and routes them to the right department (`handleVersion`, `handleBlock`, and so on) rather than trying to handle every kind of request itself in one giant function.

---

## 2. The Peer Type

Before defining `Node` itself, we need a small type to represent one connected peer — enough information to write messages back to it later, and to display its address for logging:

```go
package network

import "net"

// Peer represents one other node this Node currently has an open
// connection to, either because we dialed them or they dialed us.
type Peer struct {
	Address string   // the peer's real listen address, e.g. "127.0.0.1:3001"
	Conn    net.Conn  // the live TCP connection to that peer
}
```

`Address` deliberately stores the peer's *listen* address — the one it advertises in its `VersionPayload.Address` field (Chapter 45, Section 5) — rather than whatever ephemeral address the underlying `net.Conn` happens to report for the other end. As Chapter 45's Section 7 explained, those two addresses are not the same for the dialing side, and only the advertised listen address is useful for dialing this peer again later (which Chapter 47's peer exchange depends on directly).

---

## 3. The Node Type

Here is `network.Node`, matching this course's established shape exactly, plus one small addition — a mutex, since a node's peer map will be read and written from multiple goroutines concurrently:

```go
package network

import (
	"log"
	"net"
	"sync"

	"github.com/you/gochain/core"
)

// Node is one GoChain participant: it can accept incoming connections
// and dial outgoing ones, and it holds a reference to the same chain
// and mempool every other package in this course has been building
// toward since Volumes 3 and 5.
type Node struct {
	Address    string
	Chain      *core.Blockchain
	Mempool    *core.Mempool
	peers      map[string]*Peer
	knownAddrs map[string]bool

	mu sync.RWMutex // guards peers and knownAddrs against concurrent access
}

// NewNode creates a Node ready to Listen and Dial, wired to an existing
// chain and mempool (the same ones core.OpenBlockchain and core.Mempool
// already give you from earlier volumes).
func NewNode(address string, chain *core.Blockchain, mempool *core.Mempool) *Node {
	return &Node{
		Address:    address,
		Chain:      chain,
		Mempool:    mempool,
		peers:      make(map[string]*Peer),
		knownAddrs: make(map[string]bool),
	}
}
```

`Address` is this node's own `host:port` string — the one it will `Listen` on, and the one it announces to peers inside its own `VersionPayload`. `Chain` and `Mempool` are pointers to the exact same `core.Blockchain` and `core.Mempool` types built in Volumes 3 and 5 — `network.Node` does not reimplement any blockchain logic; it is purely the messenger that lets a chain and mempool built entirely by earlier chapters talk to another node's chain and mempool. `peers` maps a peer's advertised address to its `*Peer` (its live connection); `knownAddrs` is a simple set (a `map[string]bool` used as a set, a pattern from Chapter 04) of every address this node has ever heard about, whether or not it currently has an open connection there — Chapter 47 uses this field to track peer-exchange candidates. The unexported `mu sync.RWMutex` field is not part of the shared contract's minimal shape, but is a necessary addition: every one of `peers` and `knownAddrs` will be touched from the `Accept` loop's goroutines, from `Dial`, and from message handlers running on separate connection goroutines, all potentially at once.

---

## 4. Listening for Incoming Connections

`Listen` starts the accept loop from Chapter 44, but wraps each accepted connection in the routing logic this chapter builds, and — unlike Chapter 44's blocking example — returns once the listener is up, running the accept loop in the background so a node's `main` function can go on to dial seed peers immediately afterward:

```go
// Listen starts accepting incoming peer connections in the background
// and returns immediately once the listener is bound. Each accepted
// connection is handled on its own goroutine (Chapter 44, Section 8).
func (n *Node) Listen() error {
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		return err
	}
	log.Printf("[%s] listening for peers", n.Address)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[%s] accept error: %v", n.Address, err)
				continue
			}
			// inbound=true: this peer connected to us, so they will
			// speak first (send MsgVersion) — we reply once we see it.
			go n.handleConnection(conn, true)
		}
	}()

	return nil
}
```

The `go func() { for { ... } }()` wrapping the whole accept loop is what makes `Listen` non-blocking: the outer goroutine keeps calling `Accept` forever, in the background, while `Listen` itself has already returned `nil` to its caller. Each individual accepted connection then gets its own goroutine via `go n.handleConnection(conn, true)`, exactly the one-goroutine-per-connection pattern from Chapter 44 — the boolean `true` records that this is an *inbound* connection, meaning the other side dialed us, which `handleVersion` (Section 7) needs to know in order to decide whether it must reply with its own version or not.

---

## 5. Dialing Out to a Peer

`Dial` is the client half: it connects to another node's listen address, registers it as a peer, and — because the dialing side always speaks first in our handshake design from Chapter 45 — immediately sends its own version:

```go
// Dial connects to another node's listen address, registers it as a
// peer, and sends the first MsgVersion (the dialer always speaks
// first, per Chapter 45's handshake design).
func (n *Node) Dial(address string) error {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		return err
	}

	n.registerPeer(address, conn)

	if err := n.sendVersion(conn); err != nil {
		conn.Close()
		return err
	}

	// inbound=false: we dialed them, so we already spoke first — their
	// reply version, when it arrives, should not trigger another reply.
	go n.handleConnection(conn, false)
	return nil
}

// registerPeer records a live connection under the peer's address,
// safe for concurrent use from Listen's accept loop and from Dial.
func (n *Node) registerPeer(address string, conn net.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	n.peers[address] = &Peer{Address: address, Conn: conn}
	n.knownAddrs[address] = true
}

// sendVersion encodes and writes this node's own MsgVersion to conn.
func (n *Node) sendVersion(conn net.Conn) error {
	payload := VersionPayload{
		Version:    1,
		BestHeight: n.Chain.Height(), // core.Blockchain.Height(), from Volume 3
		Address:    n.Address,
	}
	msg, err := EncodeMessage(MsgVersion, payload)
	if err != nil {
		return err
	}
	_, err = conn.Write(msg)
	return err
}
```

`Dial` calls `net.Dial` exactly as in Chapter 44, then immediately calls `registerPeer` — recording this connection in `n.peers` right away, before the handshake even completes, so that if this node needs to send this peer something while the handshake is in flight, the connection is already trackable. `sendVersion` builds a `VersionPayload` using `n.Chain.Height()` (assume `core.Blockchain` exposes a `Height()` method returning the height of its most recently added block, as established in Volume 3) and this node's own `Address`, then uses `EncodeMessage` from Chapter 45 to produce the complete wire bytes — envelope and gob-encoded payload together — and writes them directly to the connection. Finally, `Dial` spawns `handleConnection` on its own goroutine with `inbound=false`, so this node can keep reading whatever the peer sends back (starting with their reply `MsgVersion`) without blocking `Dial`'s caller any longer than the initial connect-and-greet.

---

## 6. Routing Messages to Handlers

`handleConnection` is the loop every accepted or dialed connection runs on its own goroutine: read one message, figure out its type, hand it to the matching handler, repeat until the connection breaks.

```go
import "io"

// handleConnection reads messages from conn in a loop for as long as
// the connection stays open, dispatching each one by MessageType.
func (n *Node) handleConnection(conn net.Conn, inbound bool) {
	defer func() {
		conn.Close()
		n.removePeerByConn(conn)
	}()

	for {
		msgType, payload, err := ReadMessage(conn)
		if err != nil {
			if err != io.EOF {
				log.Printf("[%s] connection error: %v", n.Address, err)
			}
			return // peer disconnected, or the stream is unrecoverable
		}

		switch msgType {
		case MsgVersion:
			n.handleVersion(conn, inbound, payload)
		case MsgGetBlocks:
			n.handleGetBlocks(payload)
		case MsgInv:
			n.handleInv(conn, payload)
		case MsgGetData:
			n.handleGetData(conn, payload)
		case MsgBlock:
			n.handleBlock(payload)
		case MsgTx:
			n.handleTx(payload)
		case MsgAddr:
			n.handleAddr(payload)
		default:
			log.Printf("[%s] unknown message type: %d", n.Address, msgType)
		}
	}
}

// removePeerByConn drops any peer entry whose connection matches conn,
// used once a connection has closed or errored out.
func (n *Node) removePeerByConn(conn net.Conn) {
	n.mu.Lock()
	defer n.mu.Unlock()
	for addr, peer := range n.peers {
		if peer.Conn == conn {
			delete(n.peers, addr)
			log.Printf("[%s] peer %s disconnected", n.Address, addr)
			return
		}
	}
}
```

The `switch msgType` block is the entire routing table this chapter promised: every one of the seven `MessageType` constants from Chapter 45 maps to exactly one handler method. This is a deliberately simple, explicit dispatch — the same `switch`-based dispatch pattern Chapter 62 will later use for the virtual machine's opcode interpreter — rather than something cleverer like a `map[MessageType]func(...)`, because explicitness here makes it trivial to see, at a glance, everything a GoChain node knows how to respond to. `removePeerByConn` is deferred so it runs no matter *why* the loop exits — a clean disconnect, a network error, or a malformed message that made `ReadMessage` fail — keeping `n.peers` accurate without needing to duplicate cleanup logic at every exit point.

---

## 7. Handling the Version Handshake

`handleVersion` is the one handler this chapter implements fully, since it is the only message type Chapters 43-45 have fully specified end to end:

```go
// handleVersion processes an incoming MsgVersion: it logs the peer's
// reported height relative to ours, records their real listen address,
// and — if they connected to us — replies with our own version, since
// the dialing side always speaks first (Chapter 45, Section 7).
func (n *Node) handleVersion(conn net.Conn, inbound bool, payload []byte) {
	var v VersionPayload
	if err := DecodePayload(payload, &v); err != nil {
		log.Printf("[%s] bad version payload: %v", n.Address, err)
		return
	}

	n.registerPeer(v.Address, conn)

	myHeight := n.Chain.Height()
	switch {
	case v.BestHeight > myHeight:
		log.Printf("[%s] peer %s is ahead (height %d vs our %d) — sync needed (Chapter 49)",
			n.Address, v.Address, v.BestHeight, myHeight)
	case v.BestHeight < myHeight:
		log.Printf("[%s] peer %s is behind (height %d vs our %d)",
			n.Address, v.Address, v.BestHeight, myHeight)
	default:
		log.Printf("[%s] peer %s is caught up (height %d)", n.Address, v.Address, myHeight)
	}

	if inbound {
		// They dialed us and spoke first; we haven't sent our own
		// version yet on this connection, so send it now.
		if err := n.sendVersion(conn); err != nil {
			log.Printf("[%s] failed replying to %s: %v", n.Address, v.Address, err)
		}
	}
}
```

Decoding uses `DecodePayload` from Chapter 45 to unpack the raw payload bytes into a `VersionPayload` struct. `registerPeer` is called again here — this time keyed by the peer's *real* advertised address from the payload, which matters specifically for inbound connections: `Listen`'s accept loop never learns a peer's real listen address until this exact moment, since `conn.RemoteAddr()` would only report the peer's ephemeral outgoing port, not the address they are listening on themselves. The three-way `switch` is exactly the height comparison the Chapter 45 handshake diagram walked through — full synchronization logic (actually requesting and applying the missing blocks) is Chapter 49's job; for now, a node simply notices and logs where it stands relative to each peer. The final `if inbound` block is what makes the handshake symmetric without infinite-looping: only the side that *received* a version without having sent one first replies with its own — the dialing side already sent its version in `Dial` before this handler ever runs on its connection.

---

## 8. The Remaining Handlers, For Now

The other six message types get minimal, honest stub handlers in this chapter — they decode their payload and log what arrived, but their real logic is built in later chapters (Chapter 47 for `MsgAddr`, Chapter 48 for gossip-driven `MsgBlock`/`MsgTx` handling, and Chapter 49 for `MsgGetBlocks`/`MsgInv`/`MsgGetData`-driven synchronization):

```go
func (n *Node) handleGetBlocks(payload []byte) {
	var p GetBlocksPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad getblocks payload: %v", n.Address, err)
		return
	}
	log.Printf("[%s] received getblocks from %s (sync logic arrives in Chapter 49)", n.Address, p.Address)
}

func (n *Node) handleInv(conn net.Conn, payload []byte) {
	var p InvPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad inv payload: %v", n.Address, err)
		return
	}
	log.Printf("[%s] received inv: %d %s hash(es) (fetch logic arrives in Chapter 49)", n.Address, len(p.Hashes), p.Kind)
}

func (n *Node) handleGetData(conn net.Conn, payload []byte) {
	var p GetDataPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad getdata payload: %v", n.Address, err)
		return
	}
	log.Printf("[%s] received getdata for a %s (fulfilling this arrives in Chapter 49)", n.Address, p.Kind)
}

func (n *Node) handleBlock(payload []byte) {
	var p BlockPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad block payload: %v", n.Address, err)
		return
	}
	log.Printf("[%s] received a block (%d bytes) — validation/insertion arrives in Chapter 48", n.Address, len(p.Block))
}

func (n *Node) handleTx(payload []byte) {
	var p TxPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad tx payload: %v", n.Address, err)
		return
	}
	log.Printf("[%s] received a transaction (%d bytes) — mempool insertion arrives in Chapter 48", n.Address, len(p.Transaction))
}

func (n *Node) handleAddr(payload []byte) {
	var p AddrPayload
	if err := DecodePayload(payload, &p); err != nil {
		log.Printf("[%s] bad addr payload: %v", n.Address, err)
		return
	}
	log.Printf("[%s] received %d address(es) — peer-exchange logic arrives in Chapter 47", n.Address, len(p.Addresses))
}
```

Every stub follows the same shape: decode the payload into its proper struct (proving the routing and decoding machinery works end to end for all seven message types, not just `version`), then log what arrived with an honest note about where the real behavior lives. This matters more than it might look: it means the *entire* message pipeline — framing, envelope, routing, decoding — is fully exercised and testable right now, and every later chapter only needs to replace a log statement's body with real logic, never touch the routing or decoding again.

---

## 9. Broadcasting to Every Peer

The last piece of this chapter's contract is `Broadcast`, which sends the same message to every currently connected peer — the mechanism Chapter 48's gossip protocol builds directly on top of:

```go
import "encoding/binary"

// Broadcast sends msgType with the given (already payload-encoded)
// bytes to every currently connected peer. Callers typically produce
// payload by gob-encoding one of the payload structs from Chapter 45,
// e.g. a TxPayload, before calling this.
func (n *Node) Broadcast(msgType MessageType, payload []byte) {
	header := make([]byte, 5)
	header[0] = byte(msgType)
	binary.BigEndian.PutUint32(header[1:], uint32(len(payload)))

	n.mu.RLock()
	defer n.mu.RUnlock()

	for addr, peer := range n.peers {
		if _, err := peer.Conn.Write(header); err != nil {
			log.Printf("[%s] broadcast header to %s failed: %v", n.Address, addr, err)
			continue
		}
		if _, err := peer.Conn.Write(payload); err != nil {
			log.Printf("[%s] broadcast payload to %s failed: %v", n.Address, addr, err)
		}
	}
}
```

`Broadcast` builds the same 5-byte envelope header (1 byte type, 4-byte length) that `EncodeMessage` builds internally, but takes the payload as already-encoded bytes rather than a struct, since a broadcaster typically already has the encoded bytes on hand (for instance, gossip logic in Chapter 48 will hold a newly-mined block's gob-encoded `BlockPayload` once, and needs to write those same bytes to every peer without re-encoding per peer). It reads under `n.mu.RLock()` — a read lock, not a write lock — since iterating `n.peers` here only reads the map; multiple goroutines are allowed to hold read locks simultaneously, which matters once several parts of a busy node want to broadcast at the same time without unnecessarily blocking each other.

---

## 10. Running Two Nodes and Watching Them Meet

Let's wire this into a tiny runnable program, `cmd/gochain-node/main.go`, that starts one node, optionally dials a peer, and then blocks forever so we can watch its logs:

```go
package main

import (
	"flag"
	"log"
	"select"

	"github.com/you/gochain/core"
	"github.com/you/gochain/network"
)

func main() {
	addr := flag.String("addr", "127.0.0.1:3000", "this node's listen address")
	dial := flag.String("dial", "", "peer address to dial on startup (optional)")
	flag.Parse()

	// A fresh in-memory chain and mempool, same types from Volumes 3 and 5.
	chain := core.NewBlockchain()
	mempool := core.NewMempool()

	node := network.NewNode(*addr, chain, mempool)
	if err := node.Listen(); err != nil {
		log.Fatal(err)
	}

	if *dial != "" {
		if err := node.Dial(*dial); err != nil {
			log.Fatal(err)
		}
	}

	select {} // block forever — Listen's accept loop runs in the background
}
```

(Note: replace the placeholder `"select"` import with nothing — `select{}` is a built-in statement, not a package; it needs no import. This is included here only to make the blocking-forever idiom explicit in the listing above.)

Run two instances in two terminals — Node A starts first and just listens; Node B starts second and dials Node A immediately:

```bash
# terminal 1 — Node A, height 10, just listens
go run ./cmd/gochain-node -addr 127.0.0.1:3000

# terminal 2 — Node B, height 4, dials Node A on startup
go run ./cmd/gochain-node -addr 127.0.0.1:3001 -dial 127.0.0.1:3000
```

A real run produces logs like this, matching Chapter 45's handshake diagram almost line for line:

```
# terminal 1 (Node A, 127.0.0.1:3000)
2026/07/31 10:02:01 [127.0.0.1:3000] listening for peers
2026/07/31 10:02:04 [127.0.0.1:3000] peer 127.0.0.1:3001 is behind (height 4 vs our 10)

# terminal 2 (Node B, 127.0.0.1:3001)
2026/07/31 10:02:04 [127.0.0.1:3001] listening for peers
2026/07/31 10:02:04 [127.0.0.1:3001] peer 127.0.0.1:3000 is ahead (height 10 vs our 4) — sync needed (Chapter 49)
```

Walking through what just happened: Node B's `Dial` connected to Node A and immediately called `sendVersion`, so Node A's `handleConnection` goroutine read a `MsgVersion` and dispatched it to `handleVersion` with `inbound=true` — which is why Node A's very first log line is about the peer, not about sending its own version (that happens silently, inside the `if inbound` branch). Node B's own `handleConnection` goroutine, meanwhile, was already running (spawned right after `sendVersion` in `Dial`), waiting for Node A's reply — which arrives moments later and produces Node B's log line, with `inbound=false` correctly preventing Node B from replying a second time. Both logs land within the same instant because both messages cross the wire almost simultaneously, exactly as the sequence diagram in Chapter 45 predicted.

---

## Summary

- `network.Peer` bundles a peer's advertised listen address with its live `net.Conn`, matching the shared contract's `map[string]*Peer` field on `Node`.
- `network.Node` holds `Address`, `Chain`, `Mempool`, an unexported `peers` map, and an unexported `knownAddrs` set, plus a `sync.RWMutex` needed for safe concurrent access from multiple connection goroutines.
- `Listen` runs `net.Listen`/`Accept` in a background goroutine and returns immediately, spawning one `handleConnection` goroutine per accepted connection, tagged `inbound=true`.
- `Dial` connects out, registers the peer, sends this node's own `MsgVersion` first (since the dialer always speaks first), then spawns its own `handleConnection` goroutine tagged `inbound=false`.
- `handleConnection` loops on `ReadMessage`, dispatching every message to exactly one handler via a `switch` on `MessageType` — the routing table promised by this volume's design.
- `handleVersion` is fully implemented: it logs each peer's height relative to this node's own, and replies with its own version only if it did not already speak first.
- The remaining six handlers are honest, fully-decoding stubs that prove the pipeline works end to end, with real logic arriving in Chapters 47-49.
- `Broadcast` fans a message out to every connected peer using the same envelope format as `EncodeMessage`, forming the foundation Chapter 48's gossip protocol builds on directly.

---

## Exercises

### Easy

1. **Run the two-node demo from Section 10 on your own machine**, with Node A started first and Node B dialing it second. Capture and paste both terminals' logs, and annotate, in your own words, which log line corresponds to which step of Chapter 45's handshake sequence diagram.

2. **Start three nodes** (`127.0.0.1:3000`, `3001`, `3002`), with node 2 and node 3 both dialing node 1 (but not each other). Confirm node 1's logs show two separate peers registering, and explain, referencing `n.peers`, why node 1 can tell them apart.

3. **Add a log line inside `registerPeer`** printing the current size of `n.peers` every time a peer is added, and explain, in one or two sentences, why reading `len(n.peers)` safely requires holding `n.mu` even though you're "just reading."

### Medium

4. **Add a `Peers() []string` method to `Node`** that returns a snapshot list of every currently connected peer's address, safely under a read lock, and use it to print each node's peer list every 5 seconds in a background goroutine started from `Listen`.

5. **Simulate a peer disconnecting mid-handshake**: start Node A and Node B, let them handshake successfully, then kill Node B's process (Ctrl+C) and confirm Node A's logs show the "peer disconnected" message from `removePeerByConn`. Explain exactly which error `ReadMessage` returns in this situation and why it triggers cleanup rather than a crash.

6. **Change `handleVersion` to reject connections from a peer claiming a nonsensical `BestHeight`** (for example, a negative number), closing the connection and logging a warning instead of registering the peer. Explain briefly why validating *any* data received from the network — even something as simple as a height number — matters for a program that will eventually handle real financial data.

### Hard

7. **Add a version-mismatch rejection**: if a peer's `VersionPayload.Version` doesn't equal this node's own hardcoded protocol version, close the connection immediately after logging why, before ever registering the peer or exchanging any other message. Test it by hand by temporarily hardcoding two different version numbers into two node instances and confirming they refuse each other.

8. **Diagnose and fix a potential deadlock**: `Broadcast` holds `n.mu.RLock()` while calling `peer.Conn.Write` for every peer. Research (or reason through) what happens if one peer's TCP connection is stalled (its receive buffer is full and it isn't reading) — could `Broadcast` end up blocked indefinitely on that one slow peer, and if so, holding the read lock the whole time, what does that do to every other goroutine that needs `n.mu`? Propose and implement a fix, such as using per-write timeouts (`conn.SetWriteDeadline`) or moving the actual writes off the locked section.

9. **Write an integration test** (a `_test.go` file, no `main` function, using Go's `testing` package from Chapter 07) that starts two `Node`s in-process (both `Listen`, then one `Dial`s the other), waits briefly, and asserts — using the `Peers()` method from Exercise 4, or by adding test-only introspection — that each node ends up with exactly one registered peer, whose reported height matches what you configured. This is the first automated test in this course for anything involving real goroutines and real network sockets.
