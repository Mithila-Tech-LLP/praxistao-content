# Chapter 49: Blockchain Synchronization

Gossip (Chapter 48) only reaches nodes that are already connected and already listening at the moment a message is sent. A brand-new node — or one that was offline for an hour — has missed everything that happened while it wasn't around, and no amount of future gossip will hand it that history. This chapter builds the mechanism that lets a node catch up on demand: `SyncWithPeer`.

## Table of Contents

1. [Why Gossip Alone Isn't Enough](#1-why-gossip-alone-isnt-enough)
2. [The Sync Sequence, Conceptually](#2-the-sync-sequence-conceptually)
3. [Designing the Sync Messages](#3-designing-the-sync-messages)
4. [The Server Side: Answering GetBlocks and GetData](#4-the-server-side-answering-getblocks-and-getdata)
5. [Implementing SyncWithPeer](#5-implementing-syncwithpeer)
6. [Never Trust, Always Verify](#6-never-trust-always-verify)
7. [A Worked Sync Walkthrough](#7-a-worked-sync-walkthrough)
8. [Handling Sync Failures Gracefully](#8-handling-sync-failures-gracefully)
9. [When to Trigger a Sync](#9-when-to-trigger-a-sync)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Gossip Alone Isn't Enough

Imagine you've been away from your friend group's chat for a week-long trip with no signal. The moment you're back online, nobody re-sends you the entire week's conversation just because you reappeared — you have to actively ask "what did I miss?" and catch up. A blockchain node has exactly the same problem the moment it starts up, or reconnects after being offline: gossip (Chapter 48) only relays *new* transactions and blocks going forward. It carries no memory of everything that happened before the node was listening.

**Blockchain synchronization** (often shortened to "sync," and the very first sync a brand-new node does is sometimes called **initial block download**, or IBD) is the process of a node actively requesting the blocks it's missing from a peer that already has them, in order, and validating each one before trusting it. This chapter builds exactly that as `SyncWithPeer`.

---

## 2. The Sync Sequence, Conceptually

Picture two nodes: **Fresh**, which just started up with nothing but a genesis block, and **Established**, an existing node that's been running for a while and has 100 blocks. Before any transactions or mining can safely happen on Fresh, it needs Established's history.

```
   Fresh (height 0)                          Established (height 100)
        |                                              |
        |  1. "What's your current tip? Mine is        |
        |      genesis, hash G"                        |
        |---------------------------------------------->|
        |                                              |
        |  2. "I'm at height 100. Here are the         |
        |      hashes of blocks 1 through 100"         |
        |<----------------------------------------------|
        |                                              |
        |  3. "Send me block 1"                        |
        |---------------------------------------------->|
        |  4. <full block 1 data>                       |
        |<----------------------------------------------|
        |     (Fresh validates block 1, appends it)     |
        |                                              |
        |  5. "Send me block 2"                         |
        |---------------------------------------------->|
        |  6. <full block 2 data>                       |
        |<----------------------------------------------|
        |     (Fresh validates block 2, appends it)     |
        |                                              |
        |             ... repeats for blocks 3-100 ...  |
        |                                              |
        |  Fresh is now at height 100, chain tips match |
```

Four ideas are packed into this diagram, and they map directly onto the wire-protocol message types from Chapter 45: Fresh asks "what do you have that I don't?" (`MsgGetBlocks`), Established answers with a manifest of hashes (`MsgInv`, short for *inventory*), Fresh requests the actual data for each one (`MsgGetData`), and Established sends the full block (`MsgBlock`) — the very same message type gossip already uses to broadcast newly mined blocks. Sync and gossip share a message vocabulary; they differ in who initiates the exchange and why.

---

## 3. Designing the Sync Messages

`MsgGetBlocks`, `MsgInv`, and `MsgGetData` already exist as `MessageType` values from Chapter 45's protocol design — this chapter gives them their first real payloads and handlers. Following the same `encoding/gob` convention this course has used for wire payloads since Chapter 07, we define three small payload types:

```go
package network

// GetBlocksPayload is sent by a node that wants to know which blocks a peer
// has that it does not. It describes the requester's current position so
// the peer can figure out exactly what's missing.
type GetBlocksPayload struct {
	Height uint64 // requester's current chain height
	TipHash []byte // hash of the requester's current tip block
}

// InvPayload ("inventory") is the peer's answer: a list of block hashes it
// has, that the requester is missing, in chain order. It intentionally does
// not include the full block data -- that would waste bandwidth if the
// requester already has some of them from another source.
type InvPayload struct {
	Hashes [][]byte
}

// GetDataPayload asks for the full contents of one specific block, named by
// hash, after the requester has seen it listed in an InvPayload.
type GetDataPayload struct {
	Hash []byte
}
```

`MsgBlock`'s payload is simply a serialized `core.Block`, exactly as gossip already uses it in Chapter 48 — no new type needed there. Splitting "here's what I have" (`Inv`) from "send me the actual data" (`GetData`) instead of just sending every block's full contents immediately might look like an unnecessary extra round trip, but it matters in practice: it lets a requester who is only missing a handful of recent blocks avoid re-downloading ones it might already have obtained from a different peer in the meantime, and it lets a real network throttle how much data flows to a given requester at once.

---

## 4. The Server Side: Answering GetBlocks and GetData

Before `SyncWithPeer` (the requesting side) makes sense, the *answering* side needs real logic. Chapter 46's dispatch loop already routes messages by type; here are the two new handlers it calls for sync requests:

```go
import (
	"bytes"
	"encoding/gob"
	"log"
)

// handleGetBlocks answers a peer's "what am I missing?" request with an
// inventory of the block hashes it doesn't have yet, in chain order. We cap
// how many hashes we send in one batch, mirroring how Bitcoin limits a
// single getblocks response -- sending 500 at a time keeps any one exchange
// from monopolizing the connection, and the requester simply asks again
// once it has caught up on the first batch.
const maxHashesPerBatch = 500

func (n *Node) handleGetBlocks(conn *Peer, payload []byte) {
	var req GetBlocksPayload
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(&req); err != nil {
		log.Printf("sync: bad GetBlocks payload: %v", err)
		return
	}

	ourHeight := n.Chain.Height()
	if req.Height >= ourHeight {
		// The requester is already at or ahead of us -- nothing to offer.
		return
	}

	var hashes [][]byte
	for h := req.Height + 1; h <= ourHeight && len(hashes) < maxHashesPerBatch; h++ {
		block, err := n.Chain.BlockAtHeight(h)
		if err != nil {
			log.Printf("sync: could not read block at height %d: %v", h, err)
			return
		}
		hashes = append(hashes, block.Hash)
	}

	inv := InvPayload{Hashes: hashes}
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(inv); err != nil {
		log.Printf("sync: failed to encode inventory: %v", err)
		return
	}
	n.sendToPeer(conn, MsgInv, buf.Bytes())
}

// handleGetData answers a specific "send me this block" request with the
// block's full serialized contents.
func (n *Node) handleGetData(conn *Peer, payload []byte) {
	var req GetDataPayload
	if err := gob.NewDecoder(bytes.NewReader(payload)).Decode(&req); err != nil {
		log.Printf("sync: bad GetData payload: %v", err)
		return
	}

	block, err := n.Chain.BlockByHash(req.Hash)
	if err != nil {
		log.Printf("sync: peer asked for a block we don't have: %v", err)
		return
	}

	n.sendToPeer(conn, MsgBlock, block.Serialize())
}
```

`handleGetBlocks` looks up how far ahead we are of the requester, gathers up to `maxHashesPerBatch` hashes of the blocks it's missing, and replies with an `Inv` message — just the hashes, not the full blocks. `handleGetData` handles the follow-up: given one specific hash from that inventory, it looks up the full block and sends it back. Both rely on two small additions to `core.Blockchain` this chapter assumes exist alongside `AddBlock` and `ValidateBlock`: `Height() uint64`, `BlockAtHeight(h uint64) (*core.Block, error)`, and `BlockByHash(hash []byte) (*core.Block, error)` — straightforward lookups against the chain's stored blocks. `sendToPeer` is a small helper, reused from Chapter 46, that writes one framed `Envelope` to a specific peer's connection rather than broadcasting to all of them.

---

## 5. Implementing SyncWithPeer

With the server side able to answer, `SyncWithPeer` is the client-side loop that drives the whole exchange: ask what's missing, then request and validate each block in order.

```go
import (
	"errors"
	"fmt"
)

// SyncWithPeer catches this node up to a specific peer's chain. It asks the
// peer what blocks it has beyond our own tip, then requests and validates
// them one at a time, appending each to our local chain only after it
// passes every check Chapter 19 and Chapter 25 already established.
func (n *Node) SyncWithPeer(peerAddr string) error {
	peer, ok := n.getPeer(peerAddr)
	if !ok {
		return fmt.Errorf("sync: not connected to peer %s", peerAddr)
	}

	tip, err := n.Chain.Tip()
	if err != nil {
		return fmt.Errorf("sync: could not read our own tip: %w", err)
	}

	// Step 1: ask the peer what we're missing.
	req := GetBlocksPayload{Height: n.Chain.Height(), TipHash: tip.Hash}
	var reqBuf bytes.Buffer
	if err := gob.NewEncoder(&reqBuf).Encode(req); err != nil {
		return fmt.Errorf("sync: encoding request: %w", err)
	}
	n.sendToPeer(peer, MsgGetBlocks, reqBuf.Bytes())

	// Step 2: wait for the peer's inventory reply. waitForMessage is a small
	// helper (shown below) that blocks until a specific message type arrives
	// on this peer's connection, or a timeout elapses.
	invPayload, err := n.waitForMessage(peer, MsgInv, 10*time.Second)
	if err != nil {
		return fmt.Errorf("sync: waiting for inventory: %w", err)
	}
	var inv InvPayload
	if err := gob.NewDecoder(bytes.NewReader(invPayload)).Decode(&inv); err != nil {
		return fmt.Errorf("sync: decoding inventory: %w", err)
	}

	if len(inv.Hashes) == 0 {
		return nil // we're already fully caught up with this peer
	}

	// Step 3: request and validate each block in order. Order matters here
	// -- we must apply block 1 before block 2 makes any sense, since each
	// block's validity depends on the chain state left behind by the one
	// before it.
	for _, hash := range inv.Hashes {
		dataReq := GetDataPayload{Hash: hash}
		var dataBuf bytes.Buffer
		if err := gob.NewEncoder(&dataBuf).Encode(dataReq); err != nil {
			return fmt.Errorf("sync: encoding data request: %w", err)
		}
		n.sendToPeer(peer, MsgGetData, dataBuf.Bytes())

		blockPayload, err := n.waitForMessage(peer, MsgBlock, 10*time.Second)
		if err != nil {
			return fmt.Errorf("sync: waiting for block %x: %w", hash, err)
		}
		var b core.Block
		if err := gob.NewDecoder(bytes.NewReader(blockPayload)).Decode(&b); err != nil {
			return fmt.Errorf("sync: decoding block %x: %w", hash, err)
		}

		// This is the whole point of the chapter: validate every single
		// block exactly as if we had mined it ourselves, no matter how
		// trustworthy the peer seems. See Section 6.
		if err := n.Chain.ValidateBlock(&b); err != nil {
			return fmt.Errorf("sync: peer %s sent an invalid block %x: %w", peerAddr, hash, err)
		}
		if err := n.Chain.AddBlock(&b); err != nil {
			return fmt.Errorf("sync: could not add validated block %x: %w", hash, err)
		}
	}

	// If the peer had more than maxHashesPerBatch blocks for us, keep going
	// until a round trip comes back with an empty inventory.
	if len(inv.Hashes) == maxHashesPerBatch {
		return n.SyncWithPeer(peerAddr)
	}
	return nil
}
```

Two supporting pieces deserve a closer look. `getPeer` is a small accessor over the existing `peers` map, guarded by the same mutex introduced for gossip's seen-sets. `waitForMessage` is new: because a peer connection in this design is a shared, long-lived pipe that the async gossip dispatch loop (Chapter 46) is *also* reading from concurrently, `SyncWithPeer` cannot simply call `conn.Read` itself without racing that loop. Instead, it registers a short-lived "waiter" that the dispatch loop checks first, before falling through to its normal gossip handling:

```go
// waiter is a one-shot request for the next message of a specific type on a
// specific peer connection, used to make SyncWithPeer's request/response
// exchange look synchronous even though the same connection is also being
// read by the ordinary async dispatch loop.
type waiter struct {
	msgType MessageType
	result  chan []byte
}

func (n *Node) waitForMessage(p *Peer, msgType MessageType, timeout time.Duration) ([]byte, error) {
	w := &waiter{msgType: msgType, result: make(chan []byte, 1)}

	n.mu.Lock()
	p.waiters = append(p.waiters, w)
	n.mu.Unlock()

	select {
	case payload := <-w.result:
		return payload, nil
	case <-time.After(timeout):
		return nil, errors.New("timed out waiting for response")
	}
}

// deliverToWaiters is called by the dispatch loop before its normal
// handling, giving any pending SyncWithPeer call first refusal on a
// matching message. If nothing is waiting for this type, it returns false
// and the caller falls through to ordinary gossip handling (HandleTx,
// HandleBlock, and so on).
func (n *Node) deliverToWaiters(p *Peer, msgType MessageType, payload []byte) bool {
	n.mu.Lock()
	defer n.mu.Unlock()

	for i, w := range p.waiters {
		if w.msgType == msgType {
			w.result <- payload
			p.waiters = append(p.waiters[:i], p.waiters[i+1:]...)
			return true
		}
	}
	return false
}
```

This "waiter" pattern is a common way to bolt a synchronous-feeling request/response exchange on top of an otherwise asynchronous, event-driven connection, without opening a second connection or blocking the dispatch loop. It is a small piece of extra bookkeeping, but it keeps `SyncWithPeer` readable as a straight-line sequence of steps rather than a tangle of callbacks.

---

## 6. Never Trust, Always Verify

The single most important line in `SyncWithPeer` is easy to skim past: `n.Chain.ValidateBlock(&b)`, called on *every single block*, before it is ever added to our chain — with no exceptions for peers we've been connected to for a long time, peers with a lot of other honest-looking history, or peers who happen to be the very seed node we bootstrapped from in Chapter 47.

This matters because a peer, however trustworthy it looks, could be malfunctioning, running modified software, or actively malicious. If `SyncWithPeer` trusted a peer's blocks just because they arrived in response to our own request, a single bad peer could feed a fresh node an entirely fabricated history — coins that were never really mined, transactions that were never really signed. `ValidateBlock` re-derives the truth independently: it recomputes the block's hash from its actual contents and checks it matches, checks the proof-of-work actually meets the required difficulty, and checks the block correctly links to the one before it. None of that depends on trusting the peer that happened to deliver the bytes — exactly the same "don't trust, verify" principle Chapter 19 first introduced for locally tampered blocks now protects a syncing node from a dishonest network peer.

```
                     Peer sends block N
                            |
                            v
              +----------------------------+
              | ValidateBlock(N):          |
              |  - recompute hash, compare |
              |  - check proof-of-work     |
              |  - check PrevBlockHash     |
              |    matches our chain tip   |
              +----------------------------+
                    |                |
                 PASS               FAIL
                    |                |
                    v                v
             AddBlock(N),      reject block N,
             continue sync     abort sync with error
                                (do NOT quietly skip it
                                 and keep going)
```

Notice the failure branch: `SyncWithPeer` does not skip an invalid block and try the next one. It stops and returns an error. Silently skipping a bad block and continuing would leave the chain with a gap or a subtly wrong state; stopping and surfacing the error lets the caller (Chapter 52's multi-node demo, or a future automatic peer-rotation policy) decide what to do next — for example, try syncing from a different peer instead.

---

## 7. A Worked Sync Walkthrough

Suppose Established has 3 blocks beyond genesis (heights 1, 2, 3) and Fresh has only genesis (height 0). Here is exactly what happens, in terms of the messages from Section 3:

```
1. Fresh -> Established:  GetBlocksPayload{Height: 0, TipHash: <genesis hash>}

2. Established computes: ourHeight=3, req.Height=0, so hashes for
   heights 1, 2, 3 are gathered.
   Established -> Fresh:  InvPayload{Hashes: [hash1, hash2, hash3]}

3. Fresh -> Established:  GetDataPayload{Hash: hash1}
   Established -> Fresh:  MsgBlock{<full block 1>}
   Fresh: ValidateBlock(block1) -> PASS -> AddBlock(block1). Fresh height: 1

4. Fresh -> Established:  GetDataPayload{Hash: hash2}
   Established -> Fresh:  MsgBlock{<full block 2>}
   Fresh: ValidateBlock(block2) -> PASS -> AddBlock(block2). Fresh height: 2

5. Fresh -> Established:  GetDataPayload{Hash: hash3}
   Established -> Fresh:  MsgBlock{<full block 3>}
   Fresh: ValidateBlock(block3) -> PASS -> AddBlock(block3). Fresh height: 3

6. len(inv.Hashes) == 3, which is less than maxHashesPerBatch (500), so
   SyncWithPeer returns nil -- sync complete. Fresh's tip hash now equals
   Established's tip hash.
```

If Established instead had 1,200 blocks, step 2's inventory would only contain the first 500 hashes (the `maxHashesPerBatch` cap from Section 4). After Fresh works through all 500, the final `if len(inv.Hashes) == maxHashesPerBatch` check in `SyncWithPeer` fires, and it recursively calls itself to request the next batch, starting from its new height of 500 — repeating until an inventory reply comes back with fewer than the maximum, signaling there's nothing left to fetch.

---

## 8. Handling Sync Failures Gracefully

Real networks are not perfectly reliable, and this chapter's design accounts for a few concrete failure modes:

- **Timeout.** `waitForMessage`'s `time.After(timeout)` branch fires if a peer goes silent mid-sync (its process crashed, its connection dropped). `SyncWithPeer` returns an error rather than hanging forever, so the caller can retry with a different peer.
- **Invalid block from a peer.** Covered in Section 6 — the sync aborts immediately rather than silently accepting or skipping bad data.
- **Peer has nothing new.** An empty `InvPayload.Hashes` is treated as success, not an error — "we're already caught up" is a perfectly normal outcome, not a failure.
- **Peer disconnects partway through a large sync.** Since `SyncWithPeer` applies blocks one at a time and calls `AddBlock` immediately after each one validates, a dropped connection partway through still leaves the node's chain correctly extended up to whatever it managed to receive — nothing is lost, and a subsequent `SyncWithPeer` call (to the same or a different peer) picks up exactly where it left off, since the new `GetBlocksPayload.Height` reflects the progress already made.

---

## 9. When to Trigger a Sync

`SyncWithPeer` itself doesn't decide *when* to run — that's a policy decision made by the code that owns a `Node`, typically in a few places: immediately after a successful handshake with a new peer (Chapter 46's `Dial`/accept path calling `SyncWithPeer` once a `MsgVersion` exchange confirms the peer is on a compatible protocol version), on a periodic timer (say, every 30 seconds, in case gossip missed something due to a dropped connection), and manually, from a CLI command like `gochain node sync <peer-address>` for operators who want to force it. Chapter 52's multi-node major project wires up the first of these — sync-on-connect — as the concrete example.

---

## Summary

- Gossip only relays *new* messages going forward; a node that starts fresh or reconnects after downtime needs to actively request the history it missed — this is **synchronization**, and a brand-new node's first sync is called **initial block download**.
- The sync sequence has four steps mapped onto existing message types: `MsgGetBlocks` ("what am I missing?"), `MsgInv` ("here are the hashes you're missing"), `MsgGetData` ("send me this one"), and `MsgBlock` (the full data) — the same `MsgBlock` type gossip already uses.
- `handleGetBlocks` and `handleGetData` are the server-side handlers that answer these requests, capping how many hashes go out in one inventory batch (`maxHashesPerBatch`) so no single exchange monopolizes a connection.
- `SyncWithPeer(peerAddr string) error` is the client-side driver: request inventory, then request and validate each block in order, recursing to fetch further batches if the peer had more than one batch's worth.
- A small "waiter" mechanism lets `SyncWithPeer`'s request/response exchange feel synchronous on top of the same connection the async gossip dispatch loop is also reading from.
- Every single block received during sync is passed through `Chain.ValidateBlock` before being added — no exception for trusted-looking peers — and any failure aborts the sync with an error rather than silently skipping the bad block.
- Sync handles timeouts, empty replies (already caught up), and mid-sync disconnects gracefully, always leaving the chain in a valid, resumable state rather than a partially-applied broken one.

---

## Exercises

### Easy

1. **Trace by hand** what `SyncWithPeer` does if Fresh is at height 2 and Established is also at height 2 with an identical tip hash. Which message is sent, what does `handleGetBlocks` compute, and what does `SyncWithPeer` return?

2. **Add a `SyncProgress` callback** parameter to `SyncWithPeer` (a `func(current, total int)` invoked after each block is validated and added) so a caller — like a CLI command — can print a progress bar (`"Syncing... 47/500 blocks"`) instead of sync running silently.

3. **Explain in your own words** why `handleGetBlocks` sends only hashes (`InvPayload`) instead of immediately sending the full blocks the requester is missing. What situation from Section 3 would be worse if the peer always sent full block data right away?

### Medium

4. **Implement `BlockAtHeight` and `BlockByHash`** on `core.Blockchain` if you haven't already (Chapter 20's flat-file storage or Chapter 55's future BoltDB-backed storage both need to support fast lookup by height and by hash), and write tests confirming both return the correct block for a chain of at least 10 blocks, and a clear error for a height or hash that doesn't exist.

5. **Simulate a mid-sync disconnect.** Using an in-memory test harness (two `*Node`s in the same process, without real TCP sockets), start a sync of 20 blocks, forcibly interrupt the connection after 8 blocks have been validated and added, and assert that calling `SyncWithPeer` again against a fresh connection resumes from block 9 rather than re-requesting or re-applying blocks 1-8.

6. **Add a `waitForMessage` timeout test**: spin up a `Node` with a peer that never responds to a `GetBlocks` request, and assert `SyncWithPeer` returns an error within roughly the configured timeout window rather than hanging indefinitely.

### Hard

7. **Extend `GetBlocksPayload` into a real "block locator"** like Bitcoin's, instead of a single height and tip hash. A locator is a list of hashes going backward from your tip with exponentially increasing gaps (tip, tip-1, tip-2, tip-4, tip-8, tip-16, ...) so a peer can efficiently find the most recent common ancestor even if your local chain has since branched off onto a fork the peer doesn't recognize at all. Implement `BuildLocator()` and update `handleGetBlocks` to walk backward through its own chain looking for the first hash in the locator it recognizes, then reply with everything after that point.

8. **Benchmark full sync from genesis** for a chain of 10,000 blocks (you can generate a synthetic test chain rather than actually mining 10,000 real ones). Measure wall-clock time for `SyncWithPeer` to fully catch up over an in-process connection, and identify the biggest bottleneck: is it `ValidateBlock`'s proof-of-work check, the batch round-trip overhead, or something else? Propose one concrete change that would meaningfully speed this up.

9. **Implement parallel sync from multiple peers.** Real Bitcoin nodes download different block ranges from different peers simultaneously to speed up initial block download. Design (and implement, if you're up for it) a version of sync that splits a large missing range across 2-3 connected peers, requests different sub-ranges from each concurrently, and reassembles and validates the results in the correct order once all sub-ranges arrive — being careful that a slow or failing peer for one sub-range doesn't stall the whole sync.
