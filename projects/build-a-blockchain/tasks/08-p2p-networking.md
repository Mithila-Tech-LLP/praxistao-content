# Task 08: P2P Networking

## What you will build

A minimal peer-to-peer node that can send a message to another node over a real TCP connection and have it understood on the other end — the foundation every blockchain network is built on.

## Concepts

### Framing: knowing where one message ends

TCP gives you a stream of bytes with no built-in concept of "messages" — if you write two JSON objects back to back, the reader has no idea where the first one ends and the second begins unless you tell it. The standard fix is a length-prefixed frame: write the message's byte length first (as a fixed-size number), then the message itself. The reader always knows exactly how many bytes to read next.

```
  [4-byte length][message bytes....................]
  [4-byte length][message bytes..]
```

### A tiny message protocol

Real P2P protocols define several message types (handshake, inventory, get-data, and so on). For this task, two are enough: `MsgBlock` (here is a full block) and `MsgTx` (here is a full transaction) — everything a minimal node needs to tell another node about new data.

## Interface to implement

```go
type MessageType byte

const (
	MsgBlock MessageType = iota
	MsgTx
)

// SendMessage writes a length-prefixed frame to conn: one byte for
// msgType, then a 4-byte big-endian length, then payload.
func SendMessage(conn net.Conn, msgType MessageType, payload []byte) error

// ReadMessage reads one complete length-prefixed frame from conn and
// returns its type and payload.
func ReadMessage(conn net.Conn) (MessageType, []byte, error)

// Node listens on address and invokes onMessage for every message it
// receives on any connection, in a new goroutine per connection.
type Node struct {
	Address string
	// unexported fields
}

func NewNode(address string, onMessage func(MessageType, []byte)) *Node

// Listen starts accepting connections in the background. It returns
// once the listener is ready (or an error if binding the address failed).
func (n *Node) Listen() error

// Send connects to peerAddress and sends one message.
func (n *Node) Send(peerAddress string, msgType MessageType, payload []byte) error

// Close stops accepting new connections.
func (n *Node) Close() error
```

## Hints

- `encoding/binary`'s `binary.Write`/`binary.Read` (or `BigEndian.PutUint32`/`Uint32`) handle the 4-byte length prefix cleanly.
- `Listen` should start a goroutine running an `Accept()` loop; each accepted connection should be handled in its *own* goroutine so one slow peer never blocks another.
- For tests, listen on `"localhost:0"` and read back the actual assigned port from the listener, rather than hardcoding a port — this avoids test flakiness from port conflicts. Expose the resolved address (e.g. via a field or a method) so your test can connect to it.
- Test the full round trip: start a node with an `onMessage` callback that records what it received (e.g. into a channel), send it a message from a plain `net.Dial` connection (or a second `Node`), and assert the callback fires with the exact type and payload sent — using a channel with a short timeout rather than a raw `sync` primitive, since the message arrives on a different goroutine.
- Try wiring `Task 05`'s `Transaction` (gob-encoded) as a `MsgTx` payload, and confirm a node can decode it back into a working `Transaction` after receiving it.

## Run the tests

```bash
cd starter/task-08-p2p-networking
go test ./...
```

All tests must pass before moving to Task 09.
