# Chapter 45: Designing a P2P Protocol

Before writing another line of networking code, we design GoChain's wire protocol on paper: what kinds of messages nodes send each other, and the exact envelope every one of them travels inside. This is deliberate — getting the message shapes right once, here, means every later chapter (and the parallel volume work happening in Chapters 48-52) can build directly on top of them without redesigning anything.

## Table of Contents

1. [Why Protocol Design Comes Before Code](#1-why-protocol-design-comes-before-code)
2. [A Tour of GoChain's Message Types](#2-a-tour-of-gochains-message-types)
3. [The Envelope: One Header for Every Message](#3-the-envelope-one-header-for-every-message)
4. [Choosing a Payload Encoding](#4-choosing-a-payload-encoding)
5. [Defining the Payload Structs](#5-defining-the-payload-structs)
6. [Encoding and Decoding Messages in Go](#6-encoding-and-decoding-messages-in-go)
7. [The Handshake Sequence, End to End](#7-the-handshake-sequence-end-to-end)
8. [How This Compares to Bitcoin's Protocol](#8-how-this-compares-to-bitcoins-protocol)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Protocol Design Comes Before Code

A **protocol**, in networking, is simply an agreed-upon set of rules for how two programs communicate — what messages exist, what shape each one has, and what each side is supposed to do when it receives one. Two nodes running completely different codebases (imagine a hypothetical Rust or Python reimplementation of GoChain) could still interoperate perfectly, as long as both sides follow the same protocol. The protocol is the contract; the code in any one language is just one implementation of it.

Designing the protocol before writing the Go implementation is exactly like Chapter 16's approach to the `Block` struct: we decided the fields a block needs on paper, with diagrams, before defining `core.Block` in Go. Getting the shape right up front avoids a much more painful problem later — two independently-written nodes that disagree about what a message *means*, discovered only after both are already running and refusing to talk to each other.

This chapter's design is deliberately modeled closely on Bitcoin's original wire protocol (Section 8 compares them directly), both because it is a well-proven design and because recognizing these exact ideas will make real blockchain source code far less intimidating the next time you encounter it.

---

## 2. A Tour of GoChain's Message Types

Every message GoChain nodes exchange falls into one of seven kinds. Here they are, in the order a typical conversation between two nodes tends to use them:

- **`version`** — sent immediately after connecting, by both sides, before anything else. It says, in effect, "hello, here's my protocol version and my current chain height." This is the handshake message; nothing else should be trusted from a peer until a version message has been exchanged.
- **`getblocks`** — "send me the hashes of blocks you have that come after this specific hash." Used when a node suspects it is behind and wants to know what it is missing, without yet requesting the (potentially large) full block data.
- **`inv`** (short for **inventory**) — "here are hashes of things I have" — could be block hashes or transaction hashes. This is an *announcement*, not the data itself; it lets a peer say "here's what I've got" cheaply, letting the receiver decide what it actually still needs.
- **`getdata`** — "I don't have these specific hashes yet — please send me the full data for them." The natural response to receiving an `inv` message that mentions something you don't already have.
- **`block`** — a full, serialized `core.Block`, sent in response to a `getdata` request (or gossiped proactively once a block is mined — Chapter 48).
- **`tx`** — a full, serialized `core.Transaction`, the transaction equivalent of `block`.
- **`addr`** — "here are peer addresses I know about." This is how the peer-to-peer network's connectivity grows organically from a handful of seed nodes, which Chapter 47 implements in full.

```
  MESSAGE TYPES, GROUPED BY PURPOSE
  ----------------------------------

  HANDSHAKE           version

  "WHAT DO YOU HAVE?"  getblocks --> inv --> getdata
  (ask, then get told   (ask for      (peer      (ask for the
   what's available,     hashes       announces   actual full
   then ask for the      after mine)  what it     data for
   actual data)                       has)        specific hashes)

  ACTUAL DATA          block, tx

  GROWING THE NETWORK  addr
```

Every one of these becomes exactly one value of a Go type called `MessageType` — a single byte is more than enough to distinguish seven message kinds, and using the smallest type that fits keeps every message's header as compact as possible, which matters once a busy node is exchanging thousands of messages a minute.

---

## 3. The Envelope: One Header for Every Message

No matter which of the seven message types is being sent, every message on the wire starts with the exact same fixed-size header, which we call the **envelope** — the same word Chapter 01 used for the wax-sealed physical envelopes, and not a coincidence: an envelope here, too, is a wrapper that tells you what kind of thing is inside and how much of it there is, before you open it and read the contents.

```
   ONE MESSAGE ON THE WIRE, FULL BYTE LAYOUT
   -------------------------------------------

   +------------+------------------+---------------------------+
   |  Type      |  Length          |  Payload                   |
   |  1 byte    |  4 bytes (u32)   |  Length bytes               |
   +------------+------------------+---------------------------+
        ^               ^                       ^
        |               |                       |
   which of the    how many bytes         the actual message
   7 MessageType    the payload is         data — its exact shape
   values this is                          depends on Type
```

This should look immediately familiar: it is exactly the length-prefixed framing from Chapter 44, with one addition — a single byte at the front identifying the message's *type*, before the length. In Go:

```go
package network

// MessageType identifies what kind of P2P message this is.
type MessageType byte

const (
	MsgVersion   MessageType = iota // handshake / "hello, here's my chain height"
	MsgGetBlocks                    // "send me block hashes after this one"
	MsgInv                          // "here are hashes I have" (inventory announcement)
	MsgGetData                      // "send me the full data for these hashes"
	MsgBlock                        // a full serialized block
	MsgTx                           // a full serialized transaction
	MsgAddr                         // "here are peer addresses I know about"
)

// Envelope is the fixed header every message starts with.
type Envelope struct {
	Type   MessageType
	Length uint32 // length of Payload in bytes
}
```

`MessageType` is defined as a `byte` (an alias for `uint8`), so it occupies exactly one byte on the wire — plenty of room for seven kinds of message, with headroom to add more later without changing the envelope's shape. The `iota` keyword (recall Chapter 04) auto-numbers the constants starting from 0, so `MsgVersion` is `0`, `MsgGetBlocks` is `1`, and so on up to `MsgAddr` at `6`. `Envelope` bundles the `Type` and `Length` together as the logical header, even though on the wire they are written as two separate fixed-size fields, exactly like `writeFrame`/`readFrame` wrote the 4-byte length and the payload as two separate `conn.Write` calls in Chapter 44.

---

## 4. Choosing a Payload Encoding

The envelope's `Type` and `Length` fields have an obvious, fixed binary representation — a byte and a `uint32`. The *payload* — the actual `version` info, or the actual list of transaction hashes — is more complex data, and we need to pick how to turn it into bytes.

Chapter 07 compared `encoding/json`, `encoding/gob`, and hand-rolled `encoding/binary` for GoChain's block and transaction serialization, and settled on `gob` for its combination of simplicity and reasonable efficiency, with the explicit note that we might upgrade later. We make the same choice here, for the same reasons: `encoding/gob` is Go-native (no external schema needed), handles nested structs and slices automatically, and is meaningfully more compact than JSON — good enough for every message payload in this course, while keeping the code in this chapter approachable.

One detail matters: only the *payload* is gob-encoded. The envelope's `Type` and `Length` fields are written and read using the same fixed-width `encoding/binary` approach as Chapter 44's framing, not gob — this keeps the header a predictable, constant size (5 bytes: 1 for `Type`, 4 for `Length`) that can always be read first, before we even know how large the payload's gob-encoded bytes will be.

---

## 5. Defining the Payload Structs

Each message type gets its own Go struct describing exactly what data it carries. These are new types in `gochain/network`, alongside `MessageType` and `Envelope`:

```go
package network

// VersionPayload is the body of a MsgVersion message — the handshake.
type VersionPayload struct {
	Version    int    // protocol version, so future incompatible changes can be detected
	BestHeight int     // the sender's current chain height (core.Blockchain height)
	Address    string  // the sender's own listen address, e.g. "127.0.0.1:3001"
}

// GetBlocksPayload is the body of a MsgGetBlocks message.
type GetBlocksPayload struct {
	Address string // who is asking, so the response can be sent back correctly
}

// InvPayload is the body of a MsgInv message — an inventory announcement.
type InvPayload struct {
	Kind   string   // "block" or "tx" — what kind of hash this inventory lists
	Hashes [][]byte // the hashes being announced
}

// GetDataPayload is the body of a MsgGetData message — a request for one
// specific piece of full data, identified by its hash.
type GetDataPayload struct {
	Kind string // "block" or "tx"
	Hash []byte
}

// BlockPayload is the body of a MsgBlock message: a full serialized block.
type BlockPayload struct {
	Block []byte // core.Block, serialized with Block.Serialize() from Chapter 17
}

// TxPayload is the body of a MsgTx message: a full serialized transaction.
type TxPayload struct {
	Transaction []byte // core.Transaction, serialized the same way
}

// AddrPayload is the body of a MsgAddr message — peer addresses being shared.
type AddrPayload struct {
	Addresses []string // e.g. []string{"127.0.0.1:3002", "127.0.0.1:3003"}
}
```

`VersionPayload` carries exactly what a handshake needs: a protocol version number (so two nodes speaking incompatible future versions of this protocol could, in principle, detect the mismatch and refuse to proceed, rather than misinterpreting each other's bytes), the sender's current chain height (so the receiver can immediately tell whether it is ahead, behind, or caught up relative to this peer), and the sender's own listen address (so the receiving side knows where to dial back, which matters once Chapter 47 wires up peer exchange). `InvPayload` and `GetDataPayload` both use a plain `Kind string` field (`"block"` or `"tx"`) rather than a separate `MessageType` for this purpose, since the inventory itself is not a message type — it is a small piece of metadata describing what category of hash is being listed inside one `MsgInv` or `MsgGetData` message. `BlockPayload` and `TxPayload` deliberately hold already-serialized `[]byte` (the exact output of `core.Block.Serialize()` from Chapter 17 and the equivalent transaction serialization from Volume 5) rather than the live Go structs directly — this keeps `gochain/network` decoupled from needing to know every internal detail of how `core` encodes its own types; it just forwards opaque bytes and lets `gochain/core` handle deserializing them on the receiving end.

---

## 6. Encoding and Decoding Messages in Go

With the envelope and payload types defined, we can write the two functions every part of `gochain/network` will use to actually put a message on the wire and take one off:

```go
package network

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"fmt"
	"io"
)

// EncodeMessage turns a MessageType and a payload struct into the full
// wire format: [1 byte type][4 byte length][gob-encoded payload].
func EncodeMessage(msgType MessageType, payload interface{}) ([]byte, error) {
	var payloadBuf bytes.Buffer
	// gob.NewEncoder writes Go values as a compact, self-describing binary
	// format — we use it here purely for the payload, not the header.
	if err := gob.NewEncoder(&payloadBuf).Encode(payload); err != nil {
		return nil, fmt.Errorf("encoding payload: %w", err)
	}
	payloadBytes := payloadBuf.Bytes()

	var out bytes.Buffer
	out.WriteByte(byte(msgType)) // 1-byte type comes first

	lengthPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthPrefix, uint32(len(payloadBytes)))
	out.Write(lengthPrefix) // then the 4-byte length

	out.Write(payloadBytes) // then the payload itself
	return out.Bytes(), nil
}

// ReadMessage reads exactly one full message from r: its type, and its
// already gob-decoded payload bytes (the caller decodes those into the
// specific payload struct matching the returned MessageType).
func ReadMessage(r io.Reader) (MessageType, []byte, error) {
	header := make([]byte, 5) // 1 byte type + 4 bytes length
	if _, err := io.ReadFull(r, header); err != nil {
		return 0, nil, fmt.Errorf("reading envelope: %w", err)
	}

	msgType := MessageType(header[0])
	length := binary.BigEndian.Uint32(header[1:5])

	payload := make([]byte, length)
	if _, err := io.ReadFull(r, payload); err != nil {
		return 0, nil, fmt.Errorf("reading payload: %w", err)
	}

	return msgType, payload, nil
}

// DecodePayload gob-decodes raw payload bytes into out, which must be a
// pointer to the payload struct matching the message's MessageType
// (for example, &VersionPayload{} for a MsgVersion message).
func DecodePayload(payload []byte, out interface{}) error {
	return gob.NewDecoder(bytes.NewReader(payload)).Decode(out)
}
```

`EncodeMessage` builds the complete wire format for one message: it gob-encodes the payload struct first (so we know its exact byte length), then writes the 1-byte type, the 4-byte big-endian length, and finally the payload bytes, in that order — precisely the layout diagrammed in Section 3. `ReadMessage` is the mirror image, and it is deliberately generic: it reads the 5-byte header, decodes the `Type` and `Length`, then reads exactly `Length` bytes for the payload using `io.ReadFull` (recall from Chapter 44 why this must be `ReadFull` and not a bare `Read`) — but it returns the payload as raw, still-gob-encoded bytes rather than decoding it, because at this point in the code we don't yet know which payload struct to decode into; that decision depends on the `MessageType` we just read. `DecodePayload` is the small final step a caller performs once it knows, from the returned `MessageType`, which concrete struct (`VersionPayload`, `InvPayload`, and so on) the bytes should be unpacked into.

This three-function split — encode a message, read a message's header and raw payload, then decode that payload into the right struct — is exactly what `network.Node`'s connection-handling loop uses in Chapter 46 to route each incoming message to the correct handler based on its type.

---

## 7. The Handshake Sequence, End to End

With the wire format fully specified, we can trace through the very first exchange any two GoChain nodes will have, byte-by-byte in spirit if not literally. Suppose Node A, at `127.0.0.1:3000` with a chain height of 10, dials Node B, at `127.0.0.1:3001` with a chain height of 4.

```
   Node A (127.0.0.1:3000, height 10)         Node B (127.0.0.1:3001, height 4)
              |                                            |
              |------ TCP connect (net.Dial) ------------->|
              |                                            |  (Listen's Accept()
              |                                            |   returns a new conn)
              |                                            |
              |-- MsgVersion{Version:1, BestHeight:10, -->|
              |             Address:"127.0.0.1:3000"}      |
              |                                            |  B sees A is ahead
              |                                            |  (10 > 4) -- it will
              |                                            |  need to sync (Ch. 49)
              |                                            |
              |<-- MsgVersion{Version:1, BestHeight:4,  ---|
              |             Address:"127.0.0.1:3001"}      |
              |                                            |
   A sees B is behind (4 < 10) --                          |
   A also now knows B's listen address                     |
   from the payload, not just this                         |
   ephemeral connection                                    |
```

Both sides send a `MsgVersion` — the connection initiator does not wait to be asked first, and the receiver replies with its own version immediately, rather than staying silent. This symmetry matters: it means either side can independently learn the other's height and listen address from a single round of messages, with neither node needing to be treated as more "senior" than the other. Once both versions have been exchanged, both nodes now know enough to decide, on their own, whether they need to request blocks from the other (Chapter 49 implements that decision fully; for this volume's first half, we stop at "both sides know where they stand").

Note the `Address` field's importance: the TCP connection itself has an *ephemeral* source port on the dialing side (the operating system picks some arbitrary high port number for A's outgoing connection, not `3000`), so without explicitly including `"127.0.0.1:3000"` inside the `VersionPayload` itself, Node B would have no reliable way to dial *back* to Node A later — it would only see whatever ephemeral port the connection happened to arrive from. This is exactly why `VersionPayload` carries the sender's own listen address as data, rather than relying on the connection's metadata.

---

## 8. How This Compares to Bitcoin's Protocol

If you go on to read Bitcoin Core's source code or its protocol documentation, you will recognize almost everything in this chapter immediately, because GoChain's design intentionally mirrors it closely:

- Bitcoin's messages also start with a fixed-size header containing a message type identifier and a payload length (Bitcoin's version also adds a network "magic number" to distinguish mainnet from testnet traffic, and a checksum of the payload for extra corruption-detection — both reasonable additions you could add to `Envelope` as an exercise, though we omit them here for simplicity).
- Bitcoin has a `version` message exchanged first, exactly as the handshake, followed by a `verack` (version-acknowledge) confirmation — we simplify this slightly by treating the exchange of two `version` messages as sufficient, without a separate acknowledgment message.
- Bitcoin's `getblocks`, `inv`, `getdata`, `block`, and `tx` messages are essentially identical in *purpose* to GoChain's, though Bitcoin's `inv` and `getdata` messages can mix multiple types of inventory (blocks and transactions together) in one message using explicit per-item type tags, where our `InvPayload.Kind` field applies to the whole message at once for simplicity.
- Bitcoin's `addr` message serves exactly the peer-exchange purpose GoChain's `MsgAddr` will serve starting in Chapter 47, right down to the idea of nodes proactively sharing addresses they know about so the network's connectivity graph grows without any central directory.

Understanding *why* GoChain's protocol is shaped this way means you are not just copying Bitcoin's design out of tradition — you have independently arrived at, and can explain, the same reasoning that shaped one of the most successful peer-to-peer protocols ever deployed.

---

## Summary

- A protocol is an agreed-upon message format two independently-written programs can both implement and interoperate through — designing it precedes writing code, the same way `core.Block`'s fields were designed before being implemented in Go.
- GoChain's seven message types are `MsgVersion` (handshake), `MsgGetBlocks`, `MsgInv`, `MsgGetData` (the "what do you have / send it" cycle), `MsgBlock` and `MsgTx` (actual data), and `MsgAddr` (peer exchange).
- Every message shares one fixed-size `Envelope`: a 1-byte `MessageType` followed by a 4-byte payload length, exactly mirroring Chapter 44's length-prefixed framing with a type byte added.
- Payloads are gob-encoded, matching Chapter 07's earlier choice of `encoding/gob` for GoChain's serialization; only the envelope's header fields use fixed-width `encoding/binary`.
- `VersionPayload`, `GetBlocksPayload`, `InvPayload`, `GetDataPayload`, `BlockPayload`, `TxPayload`, and `AddrPayload` are the seven payload structs, one per message type.
- `EncodeMessage`, `ReadMessage`, and `DecodePayload` are the three functions that turn a typed message into wire bytes and back, and they become the shared foundation `network.Node`'s message routing uses in Chapter 46.
- The full handshake is symmetric: both sides send `MsgVersion` independently, each learning the other's chain height and listen address from the payload itself, not from the ephemeral TCP connection metadata.
- This design deliberately mirrors Bitcoin's original wire protocol closely, so the reasoning transfers directly to reading real blockchain source code later.

---

## Exercises

### Easy

1. **Draw the envelope byte layout from Section 3 from memory**, labeling each field's exact byte size, then write out, byte by byte, what the envelope (header only, not payload) would look like for a `MsgTx` message with a payload of exactly 200 bytes. (Hint: what is `MsgTx`'s numeric value given the `iota` ordering in Section 3?)

2. **For each of the seven message types**, write one sentence in your own words describing what it means and give one concrete example scenario in which a GoChain node would send it.

3. **Trace through `EncodeMessage` by hand** for a `VersionPayload{Version: 1, BestHeight: 7, Address: "127.0.0.1:3005"}` sent as a `MsgVersion`. You don't need the exact gob bytes — just describe, step by step, what order the header and payload pieces get assembled in, referencing the code in Section 6.

### Medium

4. **Add a `Checksum` field to `Envelope`** (a hash of the payload, using `gochain/crypto.Hash` from Volume 2) and update `EncodeMessage`/`ReadMessage` to compute and verify it. Explain in a comment what kind of problem this checksum protects against that the length prefix alone does not.

5. **Design (on paper, then in Go) a `MsgPing`/`MsgPong` pair of message types** that are not in the original seven, meant to let a node check whether a peer connection is still alive. Add the necessary `MessageType` constants and a `PingPayload`/`PongPayload` (they can be nearly empty structs), and explain where in the `MessageType` enum you'd add them without disrupting the existing seven values used elsewhere in this course.

6. **Compare gob and JSON encodings size-wise for a realistic `InvPayload`** containing 50 hashes (each a 32-byte `[]byte`, as real SHA-256 hashes would be). Write a small Go program that encodes the same `InvPayload` both ways and prints the resulting byte counts, then explain the size difference you observe in 100-150 words.

### Hard

7. **Implement a version-negotiation rule**: when two nodes' `VersionPayload.Version` fields don't match, decide (and implement) what should happen — refuse the connection outright, proceed anyway with a warning, or something else — and justify your choice, considering what real-world protocol version mismatches (like HTTP/1.1 vs. HTTP/2) typically do.

8. **Research Bitcoin's actual `inv`/`getdata` message format** (general knowledge or public documentation, no need for formal citation) where each inventory item carries its own type tag rather than one `Kind` field for the whole message. Redesign `InvPayload` and `GetDataPayload` to support mixed block-and-transaction hashes in a single message, implement the change, and explain in 150-200 words what new scenario this makes possible that GoChain's simpler, single-kind version cannot handle.

9. **Write a fuzz-style test** that generates random byte slices of varying lengths and feeds them directly into `ReadMessage` (simulating a malformed or malicious peer sending garbage instead of a well-formed message), and make sure your implementation never panics — only returns an error — even for inputs shorter than 5 bytes, a length prefix claiming a payload larger than the actual bytes provided, or a `Type` byte value outside the seven known `MessageType` constants. Fix any panics you find.
