# Chapter 44: TCP Servers and Clients in Go

Chapter 43 gave you the vocabulary — addresses, ports, client-server vs. peer-to-peer. This chapter puts real Go code behind it: a working TCP server and client using the standard library's `net` package, and the length-prefixed framing trick that lets two programs agree on exactly where one message ends and the next begins. Everything GoChain sends over the network from here on — version handshakes, blocks, transactions — rides on top of the exact mechanics you build in this chapter.

## Table of Contents

1. [TCP, Briefly](#1-tcp-briefly)
2. [Your First TCP Server: net.Listen](#2-your-first-tcp-server-netlisten)
3. [Your First TCP Client: net.Dial](#3-your-first-tcp-client-netdial)
4. [The Problem: Streams Have No Message Boundaries](#4-the-problem-streams-have-no-message-boundaries)
5. [The Fix: Length-Prefixed Framing](#5-the-fix-length-prefixed-framing)
6. [Building a Length-Prefixed Echo Server](#6-building-a-length-prefixed-echo-server)
7. [Building the Matching Echo Client](#7-building-the-matching-echo-client)
8. [One Goroutine Per Connection](#8-one-goroutine-per-connection)
9. [Running It End to End](#9-running-it-end-to-end)
10. [Why This Underlies Everything in This Volume](#10-why-this-underlies-everything-in-this-volume)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. TCP, Briefly

**TCP** (Transmission Control Protocol) is the networking protocol that gives you a reliable, ordered, two-way stream of bytes between two programs — once a **connection** is established, whatever bytes one side writes arrive at the other side in the same order, with the network layer handling retransmission of lost data automatically. Think of it like a phone call: once connected, you can both talk and listen, in order, without worrying that a sentence might arrive scrambled or out of sequence. This reliability is exactly why almost every blockchain, including Bitcoin and Ethereum, builds its peer-to-peer layer on top of TCP rather than the faster-but-unreliable UDP — losing or reordering a block message silently would be far worse than the small extra overhead TCP adds.

Go's standard library gives you TCP through the `net` package, and the two functions you will use constantly are `net.Listen` (to become a server, waiting for connections) and `net.Dial` (to become a client, initiating a connection). Let's use both.

---

## 2. Your First TCP Server: net.Listen

A server needs to do three things: start listening on an address, accept incoming connections one at a time as they arrive, and handle each connection. Here is the smallest version that does all three:

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
)

func main() {
	// net.Listen opens a TCP socket bound to this address and starts
	// queuing up incoming connection attempts. It does not block yet —
	// it just returns a Listener we can Accept() from.
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()

	fmt.Println("listening on 127.0.0.1:3000")

	for {
		// Accept blocks until a client dials in, then returns a net.Conn —
		// a live, two-way connection to that specific client.
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue // don't let one bad accept kill the whole server
		}

		// Handle this one connection, then loop back to Accept the next.
		handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)
	line, err := reader.ReadString('\n')
	if err != nil {
		return
	}
	fmt.Print("received: ", line)
}
```

`net.Listen("tcp", "127.0.0.1:3000")` binds a socket to that address and port and starts queuing incoming connections — this is the Go equivalent of the `nc -l 3000` command from Chapter 43. `listener.Accept()` is a blocking call: the goroutine calling it pauses until some client actually connects, at which point it returns a `net.Conn`, which is Go's representation of one specific, live connection — you read from it and write to it exactly like a file. `handleConn` reads a single line (up to a newline character) using `bufio.NewReader`, a buffered reader that makes reading line-by-line convenient.

This server has an obvious flaw for our purposes: it handles exactly one connection at a time, and it only reads one line before closing. We will fix both, but first, let's write the matching client.

---

## 3. Your First TCP Client: net.Dial

A client is even simpler — it just needs to connect and then read and write:

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	// net.Dial actively reaches out and establishes a connection to a
	// listening server. This is the "client" side of client-server, and
	// the "outgoing" side of a P2P node's dual client/server role.
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Read one line from stdin (whatever the user types) and send it,
	// terminated with a newline so the server's ReadString('\n') stops
	// reading at the right place.
	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Print("message to send: ")
	text, _ := stdinReader.ReadString('\n')

	if _, err := conn.Write([]byte(text)); err != nil {
		log.Fatal(err)
	}
}
```

`net.Dial("tcp", "127.0.0.1:3000")` opens a connection to a server already listening at that address — if nothing is listening there, this call returns an error immediately (`connection refused`), which is worth testing on purpose so you recognize the error message later. Once connected, `conn.Write` sends bytes to the server exactly the way `conn.Read` on the server side receives them; a `net.Conn` implements the same `io.Reader`/`io.Writer` interfaces you already know from Chapter 07's serialization work, which is why the same `bufio` helpers work on both ends.

Run the server in one terminal, the client in another, type a message, and you will see it printed on the server's side. This is a real, working, if extremely limited, network conversation — the same two function calls, `net.Listen`/`Accept` and `net.Dial`, are the entire foundation `network.Node` builds on starting in Chapter 46.

---

## 4. The Problem: Streams Have No Message Boundaries

The line-based approach above works only because we got lucky: a simple text message with no newline in the middle, sent once, read once. Real protocols cannot rely on "read until a newline" — GoChain will be sending serialized blocks and transactions, which are arbitrary binary data that might legitimately contain any byte value, including bytes that look like a newline character.

The deeper issue is that TCP gives you a stream of bytes, not a sequence of messages. If one side calls `conn.Write` twice in a row, quickly, the other side's `conn.Read` might see both writes arrive as one combined chunk, or one write arrive split across two reads — TCP makes no promise about preserving your original "write boundaries." This is sometimes visualized as two rubber ducks dropped into a river one after another: downstream, all you see is duck-shaped stuff floating by; you cannot tell for certain, from the water alone, exactly where one duck ended and the next began if they happened to be pushed together.

```
  Sender writes:    [ MESSAGE ONE ][ MESSAGE TWO ]
                            |
                            v  (TCP just sees a byte stream)

  Receiver's buffer: [ MESSAGE ONETMESSAGE TWO ]
                                   ^
                       where does one message end and the next begin?
                       the raw bytes alone don't say.
```

If GoChain sent a serialized `core.Block` immediately followed by a serialized `core.Transaction` with no boundary information, the receiving node would have no reliable way to know where the block's bytes end and the transaction's bytes begin. We need to embed that boundary information into the data itself.

---

## 5. The Fix: Length-Prefixed Framing

The standard solution — used by Bitcoin, Ethereum, HTTP/2, and essentially every serious binary network protocol — is **length-prefixed framing**: before sending a message, first send a fixed-size number telling the receiver exactly how many bytes are coming, then send that many bytes. The receiver always reads the length first, then reads *exactly* that many bytes for the message body, no matter what boundaries the underlying TCP stream happened to chop things into.

```
   ONE FRAMED MESSAGE ON THE WIRE
   -------------------------------

   +----------------------+----------------------------------+
   |  length (4 bytes)    |  payload (length bytes)           |
   |  e.g. 0x0000002A      |  the actual message data           |
   +----------------------+----------------------------------+

   The receiver:
     1. reads exactly 4 bytes  -> decodes them as a uint32 -> N
     2. reads exactly N bytes  -> that is the complete message
     3. repeats from step 1 for the next message
```

A **frame** here just means "one complete, self-delimited unit of data." We will use a 4-byte, big-endian unsigned integer for the length prefix — `encoding/binary`'s `binary.BigEndian` handles converting between that 4-byte representation and a normal Go `uint32`. Four bytes gives us a maximum message size of about 4 billion bytes, more than enough for any block or transaction GoChain will realistically send; if a real message ever needed to be larger, we would simply split it, but that need does not arise in this course.

---

## 6. Building a Length-Prefixed Echo Server

Let's rebuild the echo server properly, using framing instead of newline-delimited text. First, two small helper functions that every later chapter in this volume will reuse directly:

```go
package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net"
)

// writeFrame sends data over conn prefixed with its length, so the
// receiver knows exactly how many bytes to read for this one message.
func writeFrame(conn net.Conn, data []byte) error {
	lengthPrefix := make([]byte, 4)
	binary.BigEndian.PutUint32(lengthPrefix, uint32(len(data)))

	// Write the length first, then the payload. Two separate Write calls
	// are fine here — TCP still delivers both in order to the same
	// connection, and the receiver reads them as two separate steps too.
	if _, err := conn.Write(lengthPrefix); err != nil {
		return err
	}
	_, err := conn.Write(data)
	return err
}

// readFrame reads exactly one length-prefixed message from conn,
// blocking until the full message has arrived.
func readFrame(conn net.Conn) ([]byte, error) {
	lengthPrefix := make([]byte, 4)
	// io.ReadFull keeps reading until the buffer is completely filled
	// (or an error occurs) — a plain conn.Read might return fewer bytes
	// than we asked for, since TCP doesn't guarantee one Read call gets
	// everything at once.
	if _, err := io.ReadFull(conn, lengthPrefix); err != nil {
		return nil, err
	}

	length := binary.BigEndian.Uint32(lengthPrefix)
	payload := make([]byte, length)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}
	return payload, nil
}
```

`writeFrame` computes the payload's length, encodes it into 4 bytes with `binary.BigEndian.PutUint32`, and writes the length followed by the payload. `readFrame` does the reverse: read exactly 4 bytes and decode them back into a length, then read exactly that many bytes for the payload. The use of `io.ReadFull` instead of a bare `conn.Read` is not optional — `conn.Read` is allowed to return fewer bytes than the buffer's size on any single call (this is normal, expected TCP behavior, not an error), so a naive single `Read` call could silently truncate a message. `io.ReadFull` loops internally until the buffer is completely full or a real error occurs, which is exactly the guarantee framing needs.

Now the server, rewritten to use these two helpers instead of newline-delimited reads:

```go
func main() {
	listener, err := net.Listen("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer listener.Close()
	fmt.Println("echo server listening on 127.0.0.1:3000")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Println("accept error:", err)
			continue
		}
		go handleConn(conn) // one goroutine per connection — Section 8
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()
	for {
		msg, err := readFrame(conn)
		if err != nil {
			if err != io.EOF {
				log.Println("read error:", err)
			}
			return // client disconnected, or a real error — either way, stop
		}
		fmt.Printf("received %d bytes: %q\n", len(msg), msg)

		// Echo the same bytes straight back, using the same framing.
		if err := writeFrame(conn, msg); err != nil {
			log.Println("write error:", err)
			return
		}
	}
}
```

Notice the server now loops on `readFrame` *inside* `handleConn`, handling as many messages as the client sends on that one connection, rather than reading exactly one line and stopping. This matches how GoChain nodes actually behave: once two nodes are connected, they keep exchanging messages for as long as the connection stays open, not just once.

---

## 7. Building the Matching Echo Client

The client uses the exact same `writeFrame`/`readFrame` helpers — this symmetry (both sides speaking the identical framing format) is the whole point of designing a shared protocol, which Chapter 45 does properly for GoChain's actual messages.

```go
package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

func main() {
	conn, err := net.Dial("tcp", "127.0.0.1:3000")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	stdinReader := bufio.NewReader(os.Stdin)
	fmt.Println("type a message and press enter (Ctrl+C to quit):")

	for {
		fmt.Print("> ")
		text, err := stdinReader.ReadString('\n')
		if err != nil {
			return
		}

		if err := writeFrame(conn, []byte(text)); err != nil {
			log.Fatal(err)
		}

		reply, err := readFrame(conn)
		if err != nil {
			log.Fatal(err)
		}
		fmt.Printf("echoed back: %q\n", reply)
	}
}
```

Each loop iteration reads one line of user input, frames and sends it with `writeFrame`, then blocks on `readFrame` waiting for the server's echoed reply. Because both sides agree on the exact same framing rules, there is never any ambiguity about where a message starts or stops, no matter how the underlying TCP packets happen to be split up in transit.

---

## 8. One Goroutine Per Connection

Look back at the server's `Accept` loop: every time a new connection arrives, it calls `go handleConn(conn)` instead of calling `handleConn(conn)` directly. This is the standard Go networking pattern, and it matters a great deal for GoChain: without the `go` keyword, the server could only ever talk to one client at a time — it would have to finish all of `handleConn`'s work (which, for a real P2P node, never truly finishes; it loops for as long as the peer stays connected) before it could `Accept` the next incoming connection.

```
  Accept loop (one goroutine, forever)
        |
        | new connection arrives
        v
  go handleConn(connA)  -----> runs in its own goroutine, loops reading frames
        |
        | Accept() returns immediately, loops back around
        v
  new connection arrives
        v
  go handleConn(connB)  -----> runs in its own goroutine too, independently

  Now connA and connB are both being served concurrently, and the
  Accept loop is free to keep accepting connC, connD, ...
```

Recall from Chapter 05 that goroutines are cheap — a Go program can comfortably run thousands of them. A real GoChain node might have dozens of peers connected simultaneously; each one gets its own goroutine reading and dispatching messages, completely independent of every other peer's connection. This is exactly the shape `network.Node`'s connection handling will take starting in Chapter 46: `Listen()` runs an `Accept` loop that spawns a goroutine per incoming connection, and `Dial()` similarly hands its outgoing connection off to its own goroutine once the initial handshake completes.

---

## 9. Running It End to End

Save the server code as `echoserver/main.go` and the client code as `echoclient/main.go` (each in its own directory, since Go does not allow two `main` functions in the same package), remembering to put `writeFrame` and `readFrame` in both packages (or, better, in a small shared internal package — we will do exactly this properly with `gochain/network` starting next chapter).

```bash
# terminal 1
go run ./echoserver

# terminal 2
go run ./echoclient
```

A sample session looks like this:

```
# terminal 1 (server)
$ go run ./echoserver
echo server listening on 127.0.0.1:3000
received 14 bytes: "hello, gochain\n"
received 6 bytes: "again\n"

# terminal 2 (client)
$ go run ./echoclient
type a message and press enter (Ctrl+C to quit):
> hello, gochain
echoed back: "hello, gochain\n"
> again
echoed back: "again\n"
```

Try opening a *third* terminal and running the client again while the first client is still connected — the server happily handles both, each on its own goroutine, proving the one-goroutine-per-connection pattern from Section 8 in practice.

---

## 10. Why This Underlies Everything in This Volume

Every message GoChain's network layer will ever send — a version handshake, a block, a transaction, a list of peer addresses — travels as exactly one length-prefixed frame, exchanged over exactly one `net.Conn`, handled by exactly one goroutine per connection. Chapter 45 replaces the raw `[]byte` payloads used in this chapter's echo example with a proper `Envelope` (message type plus length plus payload) and real Go structs for each message type, but the `writeFrame`/`readFrame` mechanics you just built do not change at all — they simply move into `gochain/network` as the low-level transport underneath the higher-level protocol.

If you understand why a length prefix is necessary, why `io.ReadFull` matters, and why one goroutine per connection lets a node serve many peers at once, you already understand the hardest part of this entire volume. Everything from here is building a richer message format and richer behavior on top of exactly this foundation.

---

## Summary

- TCP gives two connected programs a reliable, ordered, two-way byte stream, but it does not preserve "message boundaries" — that is the application's job.
- `net.Listen` + `Accept` makes a program a server; `net.Dial` makes a program a client; a `net.Conn` is the live, two-way connection either side uses to read and write.
- Without framing, two messages written back-to-back can arrive glued together, or a single write can arrive split — length-prefixed framing solves this by sending a fixed-size length before every message's payload.
- `writeFrame`/`readFrame`, built on `encoding/binary` and `io.ReadFull`, are the exact low-level primitives every GoChain network message rides on for the rest of this volume.
- `io.ReadFull` is required, not optional, because a single `conn.Read` call may legitimately return fewer bytes than requested.
- The standard Go networking pattern is one goroutine per connection, spawned from the `Accept` loop with `go handleConn(conn)`, so a server can handle many peers concurrently without one slow peer blocking the rest.
- A minimal, fully working, length-prefixed echo server and client together prove this pattern end to end before any GoChain-specific protocol exists.

---

## Exercises

### Easy

1. **Run the length-prefixed echo server and client from this chapter** on your own machine, exactly as shown in Section 9. Open a third terminal running another client instance and confirm the server handles all connected clients independently. Paste your terminal output for at least two simultaneous clients.

2. **Modify the echo client** so that instead of echoing back the exact same text, the server responds with the message converted to uppercase. Keep the same framing helpers — only the server's response content should change.

3. **Explain, in your own words, why `readFrame` cannot simply call `conn.Read(buffer)` once** and trust that it filled the whole buffer, using the rubber-duck-in-a-river analogy from Section 4 or one of your own.

### Medium

4. **Add a maximum message size check to `readFrame`** — if the decoded length prefix is larger than some sane limit (say, 10 MB), return an error immediately instead of attempting to allocate a buffer of that size. Explain in a comment why a real network server should never trust a length prefix from an untrusted remote peer without a limit like this (hint: think about what a malicious peer could send to try to exhaust your server's memory).

5. **Benchmark the difference between using one goroutine per connection and handling connections one at a time (no `go` keyword)** by writing a small test client that opens 20 connections simultaneously and measures total time to get all 20 replies, under both versions of the server. Report the two timings and explain the difference you observe.

6. **Extend the echo protocol with a second message "kind"**: alongside the plain echo, add a `"reverse"` command that, when a client sends a frame starting with the bytes `"REVERSE:"`, causes the server to reverse the remaining bytes before echoing them back. You are not using `MessageType`/`Envelope` yet (that's Chapter 45) — implement this with simple string prefix-checking on the raw payload, and note in a comment why this ad hoc approach will not scale to more message kinds.

### Hard

7. **Simulate a slow or malicious client** that connects but sends only 2 of the 4 length-prefix bytes and then goes silent forever (do not close the connection). Verify that this connection does not block the server from accepting and serving other, well-behaved clients, and explain precisely which part of your server's design (goroutines? timeouts? neither?) is responsible for that behavior. If nothing currently protects against this connection hanging forever from the server's own resource perspective, add a `conn.SetReadDeadline` and explain what changes.

8. **Implement a version of `writeFrame`/`readFrame` that uses a 2-byte length prefix instead of 4**, and explain, with a concrete numeric example, the maximum payload size this supports and why that would or would not be enough for a serialized `core.Block` containing many transactions (refer back to Volume 3's `Block` struct and Volume 5's `Transaction` struct if you need to estimate realistic sizes).

9. **Research Go's `bufio.Scanner` with a custom `SplitFunc`** as an alternative way to implement message framing (instead of manually calling `io.ReadFull` twice), implement `readFrame` using this approach, and write 150-200 words comparing the two implementations — is the custom `SplitFunc` version clearer, faster, or neither, and would you recommend it for `gochain/network`'s real implementation in the next chapter?
