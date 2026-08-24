# Chapter 107: Building a UDP Server and Client

> **"Chapter 106's server spent half its code managing a connection's lifetime. This chapter's server has no connection to manage — and that missing half is exactly what you have to write yourself, by hand, the moment you actually need it."**

---

## Table of Contents

1. [Recap: What Chapter 106 Built, and What UDP Throws Away](#1-recap-what-chapter-106-built-and-what-udp-throws-away)
2. [The Problem: Same Goal, No Connection to Lean On](#2-the-problem-same-goal-no-connection-to-lean-on)
3. [A Naive Attempt: Porting the TCP Server's Shape Directly](#3-a-naive-attempt-porting-the-tcp-servers-shape-directly)
4. [The Real Solution: ListenUDP, ReadFromUDP, WriteToUDP](#4-the-real-solution-listenudp-readfromudp-writetoudp)
5. [Code: A Complete UDP Echo Server](#5-code-a-complete-udp-echo-server)
6. [Code: A Complete UDP Client — With Its Own Retry Logic](#6-code-a-complete-udp-client--with-its-own-retry-logic)
7. [TCP vs. UDP Code, Side by Side](#7-tcp-vs-udp-code-side-by-side)
8. [Hands-On Experiment: Run It, Break It, Watch It Not Fix Itself](#8-hands-on-experiment-run-it-break-it-watch-it-not-fix-itself)
9. [Deep Dive: Simulating Loss and Reordering On Purpose](#9-deep-dive-simulating-loss-and-reordering-on-purpose)
10. [Common Pitfalls in Concurrent Go UDP Code](#10-common-pitfalls-in-concurrent-go-udp-code)
11. [Production Notes: Buffers, MTU, Timeouts, Cancellation](#11-production-notes-buffers-mtu-timeouts-cancellation)
12. [What's Simplified Here](#12-whats-simplified-here)
13. [Interview Questions & Model Answers](#interview-questions--model-answers)
14. [Exercises](#exercises)
15. [Summary](#summary)

---

## 1. Recap: What Chapter 106 Built, and What UDP Throws Away

Chapter 106's TCP server spent real, deliberate code on things that had nothing to do with the application's actual logic: an accept loop, a goroutine per connection, read/write deadlines to protect against a silent peer, and careful handling of `Close()` and EOF. None of that was accidental complexity — it was the direct cost of the guarantees Chapter 59 through 64 described: a connection that exists, persists, and is reliably, cleanly torn down.

Chapter 58 already told you, in the abstract, that UDP has none of this: no handshake, no connection object, no ordering, no retransmission. This chapter makes that concrete by writing the code and watching what's missing. If Chapter 106 was about *managing* a connection's lifetime, this chapter is about discovering there is no lifetime to manage — and, where an application genuinely needs some of what TCP gave away for free, writing that logic by hand, deliberately, in the application layer where UDP always intended it to live (Chapter 58, Section 6).

---

## 2. The Problem: Same Goal, No Connection to Lean On

The task is superficially identical to Chapter 106's: a client sends something, a server receives it and sends something back. But every mechanism Chapter 106 used to solve sub-problems along the way is gone:

- **No `Accept()`.** There is no concept of "a new connection arriving" in UDP — Chapter 58, Section 6 established that the very first datagram *is* the entire interaction, with no setup beforehand. A UDP server doesn't accept connections; it just reads datagrams as they arrive on a single socket, from whoever happens to send one.
- **No per-client `net.Conn`.** TCP's `Accept()` handed back a brand-new socket, one per client, each with its own 4-tuple (Chapter 57, Section 6) and its own independent byte stream. A UDP server has exactly *one* socket for its entire lifetime, shared across every sender — the server has to track "who sent this" itself, datagram by datagram, using the sender's address that comes back alongside each read.
- **No ordering, no delivery guarantee, no automatic retry.** If a client sends a datagram and gets no reply, that could mean the request was lost, the reply was lost, or the server is simply slow — the client cannot tell which, and nothing beneath the application will resend anything on its behalf (Chapter 58, Section 6). If an application needs a retry, *the application* has to notice the absence of a reply and decide to resend.

This chapter's code has to confront every one of these honestly, rather than assuming them away.

---

## 3. A Naive Attempt: Porting the TCP Server's Shape Directly

Before reaching for UDP's real API, it's worth trying to reuse Chapter 106's server shape and watching exactly where it breaks.

```go
// This does not compile as written, and the comments explain exactly why.
func main() {
	listener, _ := net.Listen("udp", ":9000")  // net.Listen only supports
	                                            // stream-oriented networks (tcp,
	                                            // unix) — "udp" is rejected outright.
	for {
		conn, _ := listener.Accept()  // UDP has no concept of "a new
		                               // connection arriving" to accept.
		go handleConnection(conn)
	}
}
```

`net.Listen` is specifically Go's constructor for *connection-oriented, stream* listeners — trying `net.Listen("udp", ...)` returns an error immediately (`listen udp: unknown network udp`) precisely because "listening for new connections" is a concept that doesn't exist for a connectionless protocol. There is nothing to accept, because there is no separate, per-client object being created at all — every client, forever, is talking to the exact same single socket.

This isn't a minor API naming inconsistency; it's the same structural fact from Chapter 58, Section 6 showing up as a compiler error instead of a paragraph of prose: **UDP genuinely has no connection**, so an API modeled on "accept a new connection" cannot apply to it.

---

## 4. The Real Solution: ListenUDP, ReadFromUDP, WriteToUDP

Go's real UDP API reflects the honest shape of the protocol: one socket, read datagrams from anyone, and every read tells you exactly who sent it so you can reply to that specific sender.

```
                    ┌───────────────────────────┐
                    │   one UDP socket (:9000)  │
                    └───────────────────────────┘
                       ▲          ▲          ▲
             datagram  │          │          │  datagram
             from A    │          │          │  from C
                       │  datagram from B     │
                     ReadFromUDP() returns (n, addr, err)
                     for WHICHEVER datagram arrived next —
                     addr tells you it was A, B, or C.
```

```go
addr, _ := net.ResolveUDPAddr("udp", ":9000")
conn, _ := net.ListenUDP("udp", addr)   // ONE socket, for the server's entire life

buf := make([]byte, 1472)
for {
    n, senderAddr, err := conn.ReadFromUDP(buf)  // returns as soon as ANY
                                                    // datagram from ANYONE arrives
    if err != nil {
        continue
    }
    conn.WriteToUDP(buf[:n], senderAddr)  // reply to THIS SPECIFIC sender —
                                            // there is no other way to know
                                            // who to reply to
}
```

Notice there is no loop-inside-a-loop, no goroutine spawned per sender, and no separate object representing "the conversation with client A." `ReadFromUDP` is the entire receive-side API: it blocks until a datagram arrives (from anyone), and hands back the sender's address alongside the payload every single time, because that address is the *only* record of who sent it — there is no persistent socket object tying subsequent datagrams from the same sender together the way a TCP `net.Conn` would.

---

## 5. Code: A Complete UDP Echo Server

```go
// server.go
package main

import (
	"fmt"
	"log"
	"net"
)

func main() {
	addr, err := net.ResolveUDPAddr("udp", ":9000")
	if err != nil {
		log.Fatalf("resolve error: %v", err)
	}

	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		log.Fatalf("listen error: %v", err)
	}
	defer conn.Close()
	fmt.Println("UDP echo server listening on :9000")

	buf := make([]byte, 1472) // stays under typical Ethernet MTU minus IP/UDP headers
	for {
		n, senderAddr, err := conn.ReadFromUDP(buf)
		if err != nil {
			log.Printf("read error: %v", err)
			continue
		}

		received := string(buf[:n])
		fmt.Printf("received %d bytes from %s: %q\n", n, senderAddr, received)

		response := []byte(fmt.Sprintf("echo: %s", received))
		if _, err := conn.WriteToUDP(response, senderAddr); err != nil {
			log.Printf("write error to %s: %v", senderAddr, err)
			// deliberately no retry here — Section 6 explains why
			// that decision belongs to the CLIENT, not this server
		}
	}
}
```

Compare this to Chapter 106's server directly: there is no `Accept()`, no `go handleConnection(conn)`, and no per-client goroutine at all. **One goroutine, running one loop, serves every client this server will ever have**, because there is only ever one socket, and datagrams from different senders don't block each other the way a slow TCP read could (Chapter 106, Section 3) — `ReadFromUDP` returning one sender's datagram doesn't require waiting on any other sender in the first place.

---

## 6. Code: A Complete UDP Client — With Its Own Retry Logic

This is the section that makes Chapter 58's warnings concrete. A UDP client that wants any confidence its message was received has to build that confidence itself — nothing underneath it will retry on its behalf.

```go
// client.go
package main

import (
	"fmt"
	"log"
	"net"
	"time"
)

const (
	maxRetries     = 3
	timeoutPerTry  = 1 * time.Second
)

func main() {
	conn, err := net.Dial("udp", "127.0.0.1:9000")
	if err != nil {
		log.Fatalf("dial error: %v", err)
	}
	defer conn.Close()

	message := []byte("hello over UDP")
	reply, err := sendWithRetry(conn, message)
	if err != nil {
		fmt.Println("giving up:", err)
		return
	}
	fmt.Printf("server replied: %q\n", reply)
}

// sendWithRetry sends `message` and waits for a reply, resending up to
// maxRetries times if no reply arrives within timeoutPerTry. This entire
// function exists ONLY because UDP itself provides none of this — Chapter
// 106's TCP client needed nothing like it, because the kernel's TCP stack
// already retransmits lost segments automatically (Chapter 60).
func sendWithRetry(conn net.Conn, message []byte) ([]byte, error) {
	buf := make([]byte, 1472)

	for attempt := 1; attempt <= maxRetries; attempt++ {
		if _, err := conn.Write(message); err != nil {
			return nil, fmt.Errorf("write failed: %w", err)
		}

		conn.SetReadDeadline(time.Now().Add(timeoutPerTry))
		n, err := conn.Read(buf)
		if err == nil {
			return buf[:n], nil // success — got a reply in time
		}

		if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
			fmt.Printf("attempt %d: no reply within %v, retrying...\n", attempt, timeoutPerTry)
			continue // exactly the situation Chapter 58 warned about:
			         // we genuinely don't know if the REQUEST or the
			         // REPLY was lost — we just resend and hope
		}

		return nil, fmt.Errorf("read failed: %w", err) // a non-timeout error is fatal
	}

	return nil, fmt.Errorf("no reply after %d attempts", maxRetries)
}
```

Two details here are worth calling out explicitly, because they are exactly the honesty Chapter 58 promised this chapter would deliver on:

- **The client cannot distinguish "my request was lost" from "the reply was lost."** Both look identical from the client's point of view: silence. This is precisely why `sendWithRetry` simply resends the *original message* on every retry rather than trying to be clever about it — resending the request is the only thing that recovers from *either* failure mode.
- **A blind retry can produce a duplicate side effect.** If the server actually processed the first attempt and only its *reply* was lost, the retry causes the server to process the same message twice. For an idempotent operation (like this echo), that's harmless. For a non-idempotent one — "transfer $10," say — this is exactly the problem HTTP's idempotency keys (Chapter 71, Section 16) exist to solve, and any UDP-based protocol that needs this guarantee has to build something equivalent itself (DNS's query ID, discussed in Chapter 111, is one real example of exactly this kind of application-level bookkeeping).

---

## 7. TCP vs. UDP Code, Side by Side

| Concern | TCP (Chapter 106) | UDP (this chapter) |
|---|---|---|
| Server setup | `net.Listen` → accept loop → goroutine per connection | `net.ListenUDP` → single loop, one socket forever |
| Identifying a client | A dedicated `net.Conn` per client, created by `Accept()` | The sender's `*net.UDPAddr`, returned fresh on every `ReadFromUDP` call |
| Sending a reply | `conn.Write(data)` on that client's own connection | `conn.WriteToUDP(data, senderAddr)` — must specify the destination every time |
| Ordering | Guaranteed in-order delivery within one connection (Ch 60) | None — datagrams can arrive in a different order than sent |
| Loss handling | Automatic retransmission, invisible to the application (Ch 60) | None — the application must detect and retry, as Section 6 does |
| Connection setup cost | One handshake, one RTT, before any data flows (Ch 59) | None — the very first `Write` sends real bytes immediately |
| Concurrency model | One goroutine per connection, isolating slow clients from each other | One goroutine handling all senders (can be scaled with worker goroutines, Section 10) |
| What "close" means | A real four-way handshake tears the connection down (Ch 64) | Nothing — there's no connection to close; the socket itself is just stopped |

This table is the entire chapter's argument compressed into one place: every UDP row is either "the application must do this manually" or "not applicable, because there's nothing to track." Chapter 106's server code was long specifically because TCP's guarantees are doing real, non-trivial work under the hood; this chapter's server code is short for the mirror-image reason — there's genuinely less for it to manage.

---

## 8. Hands-On Experiment: Run It, Break It, Watch It Not Fix Itself

**Step 1 — run the server, and talk to it with `nc -u` (UDP mode):**

```
$ go run server.go
UDP echo server listening on :9000
```

```
$ nc -u 127.0.0.1 9000
hello there
echo: hello there
```

Server output:

```
received 11 bytes from 127.0.0.1:52210: "hello there"
```

**Step 2 — capture it with `tcpdump`, and note the absence of a handshake, directly contrasting Chapter 106, Section 9:**

```
$ sudo tcpdump -i lo -n udp port 9000 -X

12:00:01.000 IP 127.0.0.1.52210 > 127.0.0.1.9000: UDP, length 11
12:00:01.001 IP 127.0.0.1.9000 > 127.0.0.1.52210: UDP, length 16
```

Exactly two packets total — one datagram in, one datagram out — with **no SYN, no SYN-ACK, no ACK preceding them**. Compare this directly against Chapter 106, Section 9's capture, which needed three packets of pure handshake bookkeeping *before* a single byte of real data moved. This is Chapter 58's "connectionless" claim, made visible in your own packet capture.

**Step 3 — kill the server mid-conversation and watch the client not notice immediately:**

```
$ go run client.go
```
(In another terminal, kill the server with Ctrl-C right after it starts.)
```
attempt 1: no reply within 1s, retrying...
attempt 2: no reply within 1s, retrying...
attempt 3: no reply within 1s, retrying...
giving up: no reply after 3 attempts
```

Unlike a TCP client, which would get an immediate, explicit signal (a RST, or a `connection reset by peer` error, once it tried to write to a server that's gone), a UDP client gets nothing — it can only infer failure from silence, after waiting out its own self-imposed timeouts. This is the direct, observable cost of not having a connection object that the OS can definitively mark as dead.

---

## 9. Deep Dive: Simulating Loss and Reordering On Purpose

It's one thing to read that UDP doesn't guarantee ordering; it's another to force it to happen and watch application code cope (or fail to). This experiment deliberately introduces loss and reordering in a way that's easy to reproduce.

```go
// lossy_client.go — sends 5 numbered datagrams, deliberately dropping #3,
// and sends them in a scrambled order to simulate what the network is
// always, honestly, allowed to do.
package main

import (
	"fmt"
	"net"
	"time"
)

func main() {
	conn, _ := net.Dial("udp", "127.0.0.1:9000")
	defer conn.Close()

	order := []int{1, 2, 4, 5} // #3 is deliberately never sent — simulating loss
	for _, seq := range order {
		msg := fmt.Sprintf("packet #%d", seq)
		conn.Write([]byte(msg))
		fmt.Println("sent:", msg)
		time.Sleep(50 * time.Millisecond)
	}
}
```

Run this against Section 5's server and read its log:

```
received 9 bytes from 127.0.0.1:52310: "packet #1"
received 9 bytes from 127.0.0.1:52310: "packet #2"
received 9 bytes from 127.0.0.1:52310: "packet #4"
received 9 bytes from 127.0.0.1:52310: "packet #5"
```

Nothing in the server or the network complains about the missing `#3` — as far as UDP is concerned, four independent datagrams arrived, in the order they were sent, and that's a perfectly ordinary, unremarkable outcome. If `#3` mattered (an important game event, a critical status update), **only application-level code that put its own sequence number inside the payload** — exactly as shown here — would even be capable of detecting that a gap exists at all. This is the concrete mechanism behind Chapter 58, Section 8's claim that real-time protocols like RTP add their own sequence numbers on top of UDP: without them, a receiver has no way to know a gap occurred, only that the datagrams it did receive were internally in order relative to each other's arrival.

---

## 10. Common Pitfalls in Concurrent Go UDP Code

- **Assuming one socket needs one goroutine per sender, mirroring the TCP pattern.** A UDP server has no notion of a persistent per-client connection to hand off to a dedicated goroutine — `ReadFromUDP` already interleaves all senders through one call. If per-datagram processing is CPU-heavy enough to want concurrency, the correct pattern is a **worker pool** reading from a channel the main loop feeds, not a goroutine-per-sender model that has nothing stable to key off of:

```go
type packet struct {
    data []byte
    addr *net.UDPAddr
}

jobs := make(chan packet, 100)
for i := 0; i < 8; i++ {
    go worker(conn, jobs) // 8 workers processing datagrams concurrently
}
for {
    buf := make([]byte, 1472) // a FRESH buffer per read — see next bullet
    n, addr, _ := conn.ReadFromUDP(buf)
    jobs <- packet{data: buf[:n], addr: addr}
}
```

- **Reusing one shared buffer across `ReadFromUDP` calls when processing happens concurrently.** If a single `buf` is reused for every read (as Section 5's simple server does, safely, because it processes each datagram fully before reading the next), handing that same underlying array off to a worker goroutine while the main loop immediately overwrites it with the *next* datagram is a data race — the worker may read corrupted or half-overwritten data. The fix, shown above, is allocating a new buffer per datagram before dispatching to a worker.
- **Forgetting that `WriteToUDP` can partially fail per-destination without affecting anything else.** A write to one unreachable or firewalled client returning an error should not stop the server's main loop from continuing to serve every other client — Section 5's server deliberately logs and continues rather than returning or panicking.
- **Treating a UDP "connection" (from `net.Dial("udp", ...)`) as proof the destination exists.** `net.Dial("udp", ...)` succeeds immediately and unconditionally, because — unlike TCP's `Dial`, which performs a real handshake that can fail — it only records a destination address locally and never touches the network at all until the first `Write`. A successful `Dial` call for UDP proves nothing about reachability, which is precisely why Section 6's client needs its own retry logic rather than trusting `Dial`'s success as a signal.
- **Sizing the receive buffer smaller than the largest datagram a sender might send.** `ReadFromUDP` silently truncates a datagram larger than the provided buffer — there's no error, just quietly lost trailing bytes. Sizing the buffer at least to the practical MTU-safe ceiling (1472 bytes on typical Ethernet, per Chapter 58, Section 3) avoids this for well-behaved senders, though a server accepting arbitrary internet traffic should validate this defensively rather than assume it.

---

## 11. Production Notes: Buffers, MTU, Timeouts, Cancellation

- **Buffer sizing is a real, concrete production decision**, not an arbitrary constant. `1472` bytes (Chapter 58, Section 3) is the largest UDP payload that fits inside one standard 1500-byte Ethernet frame without IP fragmentation. A server expecting only small, controlled messages can use a smaller buffer; one accepting arbitrary UDP traffic from the internet should size for the protocol's actual maximum (UDP's theoretical ceiling is 65,507 bytes of payload) or reject oversized input explicitly.
- **Amplification risk is UDP-specific and directly relevant to any UDP server exposed to the internet.** Because there's no handshake proving the sender's address is real (Chapter 58, Section 13; Chapter 83 covers this fully), a server that replies with a response larger than the request it received can be abused: an attacker spoofs a victim's source address, sends a small request, and the server unwittingly blasts a much larger reply at the victim. Production UDP services (DNS resolvers are the classic example) deliberately cap response sizes and/or rate-limit per source address specifically to blunt this.
- **`context.Context` for graceful shutdown works differently than Chapter 106's version**, because there's no listener whose `Close()` naturally unblocks a pending `Accept()`. Closing the UDP `*net.UDPConn` itself is what unblocks a pending `ReadFromUDP` with an error, which the loop should treat as a deliberate shutdown signal, exactly mirroring Chapter 106, Section 13's pattern but applied to `conn.Close()` instead of `listener.Close()`.
- **Retry and backoff parameters (Section 6's `maxRetries` and `timeoutPerTry`) are real tuning knobs**, not implementation details to hardcode and forget. A client on a high-latency or lossy path (satellite links, Chapter 23; congested Wi-Fi, Volume 13) needs a longer per-attempt timeout and often exponential backoff between retries, to avoid making network congestion worse by retrying aggressively into an already-struggling path — precisely the kind of self-restraint TCP's congestion control (Chapter 62) provides automatically and UDP applications must approximate by hand if they care.
- **`SO_REUSEPORT` (Chapter 57, Section 8) applies to UDP sockets too**, letting multiple worker processes each bind the identical UDP port and have the kernel spread incoming datagrams across them by hash — a real way to scale a high-throughput UDP service (a game server or DNS resolver, say) across CPU cores without the worker-pool-behind-one-socket pattern from Section 10 being the only option.

---

## 12. What's Simplified Here

Section 6's retry logic is deliberately minimal — a real production client would typically use exponential backoff (increasing the wait between retries) rather than a fixed `timeoutPerTry`, and would often distinguish "give up entirely" from "the destination is unreachable right now, try later." This chapter also doesn't implement anything resembling the sequencing, jitter buffering, or forward error correction that real UDP-based protocols (RTP, QUIC) layer on top of raw datagrams — Section 9's numbered-packet experiment demonstrates the *problem* those protocols solve, not their actual solutions, which are substantial engineering efforts in their own right (QUIC's approach is previewed in Chapter 75).

---

## Interview Questions & Model Answers

**Beginner: Why doesn't Go's `net.Listen` function work for UDP servers?**

`net.Listen` is specifically for connection-oriented protocols where the runtime needs to accept a stream of *new connections* over time, each becoming its own independent socket — that's a TCP-shaped operation. UDP has no connections to accept at all; every datagram from every sender arrives on the same single socket. Go models this with a different, honest API — `net.ListenUDP` to create that one socket, and `ReadFromUDP`/`WriteToUDP` to receive from and reply to whichever sender's datagram is being handled, with the sender's address passed explicitly on every call instead of being implied by a dedicated connection object.

**Intermediate: A UDP client sends a request and receives no reply within its timeout. What can it actually conclude, and what can't it conclude?**

It can conclude only that no reply arrived in time — nothing more specific than that. It cannot distinguish whether its own request never reached the server, whether the server processed the request but its reply was lost on the way back, or whether the server is simply slow and a reply is still coming. All three look identical: silence. This ambiguity is exactly why a naive "just resend on timeout" strategy can cause a request to be processed twice if the original request actually succeeded and only the reply was lost — safe for an idempotent operation, unsafe for one that isn't, unless the application adds its own deduplication mechanism (an idea directly analogous to HTTP's idempotency keys, Chapter 71).

**Advanced: Explain, mechanically, how a naive UDP-based echo service could be abused for a DDoS amplification attack, and what design choice specifically defends against it.**

The attack relies on two properties: UDP has no handshake verifying that the sender's source IP address is genuine (unlike TCP, where a spoofed source address can't complete a real three-way handshake, Chapter 59), and the server replies to whatever source address the incoming datagram claims, without any prior interaction proving that address is reachable or willing to receive a reply. An attacker sends a request with the victim's IP forged as the source address; the server, having no way to detect the forgery, sends its reply directly to the victim instead of the real sender. If the reply is significantly larger than the request — which is common for protocols like DNS, where a small query can trigger a much larger response — the attacker achieves amplification: a small amount of attacker bandwidth becomes a much larger flood directed at the victim. Defenses include keeping response sizes close to request sizes where possible, rate-limiting replies per source address, and, at the network level, ingress filtering that rejects packets with obviously spoofed source addresses before they ever reach the application (BCP 38, mentioned again in Chapter 83's fuller treatment of spoofing).

---

## Exercises

### Easy
1. Run Section 5's server and Section 6's client, then kill the server before the client sends its message — describe exactly what the client prints and why.
2. Modify the server to print, alongside each received datagram, the total count of unique sender addresses it has seen so far in this run.
3. Using `nc -u`, send five separate one-line messages to the running server in a single session and confirm each gets its own independent echoed reply.

### Medium
4. Modify Section 6's client to use exponential backoff (e.g., 500ms, 1s, 2s) between retries instead of a fixed `timeoutPerTry`, and explain in a comment why this is friendlier to a congested network than fixed-interval retrying.
5. Extend Section 9's lossy-client experiment so the server, upon noticing a numbered packet arrived out of the expected sequence order (not just a gap), logs a warning — implement just enough application-level sequence tracking to detect this, without attempting to reorder or buffer anything.
6. Rewrite the worker-pool pattern from Section 10 so that each worker also writes a reply back to the correct sender, and verify with concurrent `nc -u` clients that replies never get sent to the wrong sender.

### Hard
7. Implement a simple request/response protocol on top of UDP where each request carries a client-generated random request ID (as part of the payload), and the server echoes that ID back in its reply. Modify the client so that if it receives a reply with an ID that doesn't match its most recent outstanding request (simulating a stale, delayed reply from an earlier retry arriving late), it discards it and keeps waiting — this is a simplified rehearsal of DNS's real query-ID matching (Chapter 111).
8. Add a basic rate limiter to the server (e.g., no more than 10 datagrams per second from any single source address, tracked with a map and timestamps) and explain, referencing Section 11, exactly what class of attack this defends against.
9. Benchmark the throughput difference between Chapter 106's TCP echo server and this chapter's UDP echo server, both under loopback conditions with no simulated packet loss, using a simple load-generating client for each. Explain the result in terms of what each protocol's connection setup and per-message overhead actually costs.

---

## Summary

| Term | Meaning |
|---|---|
| `net.ListenUDP` | Creates the single, permanent socket a UDP server uses for every sender it will ever talk to |
| `ReadFromUDP` / `WriteToUDP` | Receive from, and reply to, a specific sender identified by address on every call |
| No accept loop | UDP has no "new connection" event — there is nothing to accept |
| Application-level retry | The client's own responsibility, since UDP will not resend anything automatically |
| Ambiguous silence | A UDP client cannot tell a lost request from a lost reply from a slow server |
| Amplification attack | Abusing a spoofable sender address plus a large reply to flood an unwitting victim |
| Worker pool (for UDP) | The correct concurrency pattern for UDP, replacing TCP's goroutine-per-connection model |

You've now built both transport-layer options this course spent Volume 9 explaining, and seen exactly what code each one costs and saves you. Chapter 108 goes back to TCP and puts Chapter 106's accept loop to real, multi-user work: a chat server where dozens of clients are connected at once, and a message from one has to be broadcast to all the others — the first genuinely concurrent, stateful application this course builds.
