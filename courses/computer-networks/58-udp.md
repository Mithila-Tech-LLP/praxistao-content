# Chapter 58: UDP — The Simplest Transport Protocol

> **"UDP doesn't lie to you. It never promised delivery, order, or anything else — and that honesty is exactly what a surprising number of applications need."**

---

## Table of Contents

1. [What's Left to Solve After Ports](#1-whats-left-to-solve-after-ports)
2. [The Design Question: What's the Minimum Useful Addition to IP?](#2-the-design-question-whats-the-minimum-useful-addition-to-ip)
3. [The UDP Header, Field by Field](#3-the-udp-header-field-by-field)
4. [A Full Worked Packet, Byte by Byte](#4-a-full-worked-packet-byte-by-byte)
5. [The Checksum, Computed by Hand](#5-the-checksum-computed-by-hand)
6. [Connectionless: What It Actually Means](#6-connectionless-what-it-actually-means)
7. [Why "Unreliable, Unordered" Is a Feature](#7-why-unreliable-unordered-is-a-feature)
8. [Real Use Cases, One by One](#8-real-use-cases-one-by-one)
9. [Broadcast and Multicast: What TCP Cannot Do](#9-broadcast-and-multicast-what-tcp-cannot-do)
10. [UDP Code: A Minimal Client and Server](#10-udp-code-a-minimal-client-and-server)
11. [Seeing It Live](#11-seeing-it-live)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Production Notes](#13-production-notes)
14. [What's Simplified Here](#14-whats-simplified-here)
15. [Interview Questions & Model Answers](#15-interview-questions--model-answers)
16. [Exercises](#16-exercises)
17. [Summary](#17-summary)

---

## 1. What's Left to Solve After Ports

Chapter 57 fixed one gap left by IP: identifying *which program* on a machine a packet is for, via ports. But IP (Chapter 36) leaves several other gaps completely open, on purpose — it was designed to do one job only (best-effort delivery of a packet from one address to another) and to push everything else up to higher layers:

- IP can silently **drop** a packet if a router's queue is full.
- IP can **duplicate** a packet if a link retransmits at a lower layer and both copies survive.
- IP can deliver packets **out of order**, since different packets from the same conversation can take different routes (Chapter 45) and arrive at different times.
- IP has **no concept of a "connection"** at all — every packet is routed independently, with no memory of previous packets.

None of that is a bug. It's the entire point of packet switching (Chapter 9): keep the network's core simple and stateless, and let the endpoints decide how much extra reliability they need, since not every application needs the same amount.

This raises the actual design question this chapter is about: **what is the smallest possible addition to IP that still makes it usable by real applications?**

---

## 2. The Design Question: What's the Minimum Useful Addition to IP?

Suppose you were designing a transport protocol and had committed to the smallest possible header. IP already handles addressing, routing, and fragmentation. What's genuinely missing that an application can't get from IP alone?

Two things, and only two:

1. **Ports.** As Chapter 57 established, IP has no notion of "which program." Whatever sits above IP has to carry source and destination port numbers, or applications can never share a machine.
2. **A way to detect corruption end-to-end.** IP does have its own header checksum (protecting just the IP header, not the payload), and lower layers like Ethernet have their own frame check sequence (Chapter 28). But no layer below the application actually validates that the *payload bytes* survived the trip untouched, end to end, across every hop, link technology, and NAT box in between. A transport-layer checksum closes that last gap.

That's it. No handshake, no sequence numbers, no acknowledgments, no retransmission, no flow control. **UDP (User Datagram Protocol, RFC 768, 1980) is exactly this minimum: ports plus a checksum, wrapped around whatever payload the application hands it.** Its entire specification fits on a single page — it is, quite deliberately, the least you can add to IP and still call the result a "transport protocol."

```
        Application data
               │
               ▼
        ┌─────────────┐
        │  UDP header │  8 bytes: src port, dst port, length, checksum
        └─────────────┘
               │
               ▼
        ┌─────────────┐
        │  IP header  │  handles addressing + routing (Ch 36, 45)
        └─────────────┘
               │
               ▼
         Ethernet frame  (Ch 28)
```

The resulting unit is called a **UDP datagram** — "datagram" precisely because, like an IP packet, each one is self-contained and independently routed, with no relationship enforced between one datagram and the next.

---

## 3. The UDP Header, Field by Field

The entire UDP header is 8 bytes — four fields, each exactly 16 bits (2 bytes):

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|          Source Port          |       Destination Port       |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|            Length             |           Checksum           |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                                                               |
|                       Application Data                       |
|                                                               |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Size | Meaning |
|---|---|---|
| Source Port | 16 bits | The sending program's port (Chapter 57). Can be 0 if the sender expects no reply. |
| Destination Port | 16 bits | The receiving program's port. Never 0 for a meaningful datagram. |
| Length | 16 bits | Total length of the UDP datagram **including this 8-byte header**, in bytes. Minimum legal value is 8 (header, zero-length payload). |
| Checksum | 16 bits | Error-detection value covering the UDP header, the payload, and a "pseudo-header" borrowed from the IP layer (Section 5). Optional in IPv4 (a value of all-zero means "not computed"); mandatory in IPv6. |

Compare that against the TCP header you'll meet starting in Chapter 59, and fully in Chapter 65: TCP's minimum header is 20 bytes and carries sequence numbers, acknowledgment numbers, flags, and a window size, none of which UDP has. That difference in size and complexity is the entire story of this volume — it's the cost of the guarantees TCP adds on top of exactly this same foundation.

The 16-bit Length field caps a single UDP datagram at 65,535 bytes total (65,527 bytes of actual payload), though in practice IP fragmentation and network MTUs (typically 1,500 bytes on Ethernet) mean most real UDP payloads are kept well under 1,472 bytes to avoid fragmentation altogether — a detail DNS resolvers (Chapter 69) and VoIP applications both engineer around carefully.

---

## 4. A Full Worked Packet, Byte by Byte

Let's build one real UDP datagram completely, the way a packet capture tool like Wireshark or `tcpdump` (Chapter 119) would show it to you.

**Scenario:** a client at 203.0.113.5, using ephemeral port 51342, sends the 5-byte ASCII payload `hello` to a server at 198.51.100.7, port 53.

(Real DNS messages have their own binary format, covered starting in Chapter 66 — we're using a plain 5-byte payload here so the UDP framing itself, and the checksum arithmetic in Section 5, stay easy to follow by hand.)

```
Source port      51342  = 0xC08E
Destination port    53  = 0x0035
Length     8 (header) + 5 (payload) = 13 = 0x000D
Checksum         (computed in Section 5) = 0x94FD
Payload          "hello" = 0x68 0x65 0x6C 0x6C 0x6F
```

The datagram, as bytes on the wire (13 bytes total):

```
Offset  Bytes           Field
0x00    C0 8E           Source Port      = 51342
0x02    00 35           Destination Port = 53
0x04    00 0D           Length           = 13
0x06    94 FD           Checksum         = 0x94FD
0x08    68 65 6C 6C 6F  Payload          = "hello"
```

As a single hex dump, the way `tcpdump -xx` or Wireshark's hex pane would print it:

```
0000   c0 8e 00 35 00 0d 94 fd 68 65 6c 6c 6f
```

That's the entire transport-layer contribution to this packet — 8 bytes of header wrapped tightly around whatever the application wanted to send, with no further ceremony.

---

## 5. The Checksum, Computed by Hand

The checksum field is the one part of the header that isn't a plain copy of a number an application handed the OS — it has to be *computed*, and it's worth doing once by hand so it stops being a black box.

UDP's checksum doesn't just cover the UDP header and payload — it also covers a **pseudo-header**: a handful of fields borrowed from the IP header (source IP, destination IP, protocol number, and UDP length) that aren't actually transmitted as part of the UDP datagram, but are included in the checksum calculation anyway. This is deliberate: it means a UDP checksum also protects against a datagram somehow arriving with a corrupted destination or source IP address — otherwise a UDP payload could pass its checksum perfectly while quietly having been delivered to (or claimed to be from) the wrong machine.

**Pseudo-header for our example:**

```
Source IP:        203.0.113.5   → CB 00 71 05
Destination IP:    198.51.100.7  → C6 33 64 07
Zero byte:                          00
Protocol (UDP=17):                  11
UDP length:        13            → 00 0D
```

**The algorithm (RFC 768 + RFC 1071):** treat the pseudo-header, the UDP header (with the checksum field itself temporarily set to zero), and the payload (padded with one zero byte if its length is odd) as a sequence of 16-bit words. Add them all together using one's-complement arithmetic (any carry out of the top bit wraps around and gets added back in). Then take the one's complement (bitwise NOT) of the final sum. That's the checksum.

Laid out as 16-bit words:

```
Pseudo-header:  CB00  7105  C633  6407  0011  000D
UDP header:     C08E  0035  000D  0000   (checksum field zeroed for calculation)
Payload+pad:    6865  6C6C  6F00          ("he" "ll" "o"+pad)
```

Summing all thirteen words with end-around carry gives a folded 16-bit sum of `0x6B02`. Taking the one's complement:

```
  0xFFFF
- 0x6B02
--------
  0x94FD
```

That `0x94FD` is exactly the checksum value used in Section 4's hex dump. The receiver redoes the identical calculation — pseudo-header plus the received header (checksum field included this time, unmodified) plus payload — and if the datagram wasn't corrupted in transit, the sum of *everything including the checksum field itself* comes out to `0xFFFF` (all one bits), which is the standard way one's-complement checksums self-verify.

**When checksums are skipped:** in IPv4, a sender is allowed to set the UDP checksum field to all zeros to mean "I didn't compute one" — a leftover from an era when CPU cycles for checksum computation were expensive and some link layers already guaranteed integrity. IPv6 removed this option: since IPv6 has no IP-layer checksum of its own, UDP's checksum became mandatory there, because otherwise nothing at all would protect payload integrity end to end.

---

## 6. Connectionless: What It Actually Means

UDP is called **connectionless**, and it's worth being precise about what that means mechanically, not just as a slogan.

There is no handshake before the first UDP datagram is sent (contrast this with TCP's three-way handshake, Chapter 59) — a program can call `sendto()` on a fresh socket and a real datagram leaves the machine immediately, with no prior negotiation. There is also no shared state kept about "this conversation" on either side purely because of UDP itself: the operating system doesn't track "datagram 1 of a session," doesn't buffer datagrams waiting for missing ones to arrive, and doesn't remember what was sent five seconds ago when deciding what to do with what arrives now. Each datagram is handled completely independently, exactly like the IP packet carrying it.

```
Client                                   Server
  │                                         │
  │  UDP datagram (payload A) ────────────▶│   no handshake, no setup —
  │                                         │   the very first datagram
  │  UDP datagram (payload B) ────────────▶│   IS the "connection"
  │                                         │
  │◀──────────── UDP datagram (reply) ─────│
```

A UDP "server" is really just a program that calls `bind()` on a port and reads whatever datagrams show up — it has no way to know, from UDP alone, whether a given sender's earlier datagram ever arrived, or whether the sender is still "there" at all. (Some programs build a *session concept* on top of UDP anyway, using their own sequence numbers and keep-alive messages in the payload — that's exactly what QUIC, Chapter 75, and many game engines do. UDP itself, though, provides none of it.)

---

## 7. Why "Unreliable, Unordered" Is a Feature

It's tempting to read "no guaranteed delivery, no ordering, no retransmission" as a list of missing features — as if UDP were an unfinished TCP. That's backwards. For a specific, common category of application, every one of TCP's guarantees is actively the wrong trade-off, and here's the reasoning, worked through rather than asserted:

**The core trade-off:** guaranteed, in-order delivery requires *waiting*. If TCP detects a lost segment, it must hold back everything that arrived after it — even data the application already has in hand — until the missing piece is retransmitted and arrives (this is "head-of-line blocking," which you'll see again from a different angle in Chapters 73–75). For a file download, that wait is obviously the right choice: a file with a missing chunk in the middle is useless. But for other applications, the wait itself is the failure.

**DNS (Chapter 66):** a DNS query is one small request, one small response. If it's lost, the resolver just asks again — cheaply, in milliseconds, with no state to clean up on either side, because there was never a connection to begin with. Building a full TCP handshake for a single 60-byte question would roughly triple the number of round trips needed for the common case, for a guarantee (ordering) that's meaningless when there's only ever one packet in flight at a time. DNS does fall back to TCP for large responses (Chapter 57 mentioned TCP/53 exists for exactly this), but the default, common-case path is UDP precisely because the guarantees TCP adds aren't worth their cost here.

**Live video and voice calls (VoIP):** if frame 500 of a video stream is lost, waiting for it to be retransmitted means the whole stream stalls while a frame that's now half a second stale finally arrives — at which point the *viewer* has already moved past that moment. The right behavior is almost always "skip it, keep playing," which is only possible if nothing lower down is silently blocking newer data behind older, lost data. Real-time video/voice protocols (RTP, itself usually carried over UDP) build their own much lighter-weight handling on top — sequence numbers for the *application* to detect loss and conceal it (interpolate audio, hold the last video frame briefly), never automatic retransmission.

**Online gaming:** a player's position update from two ticks ago is worse than useless once a newer one has already arrived — retransmitting the old one and blocking the new one behind it would make the game feel laggy in exchange for delivering information the game no longer needs. Fast-paced multiplayer games almost universally use UDP and layer their own "only the newest position matters" logic on top, exactly because TCP's insistence on strict order is a mismatch for data that expires within milliseconds.

The general shape of the argument, stated once so it doesn't need repeating per example: **TCP is correct when every byte matters and staleness is unacceptable; UDP (plus lightweight application-level handling) is correct when timeliness matters more than completeness, and stale data is actively worse than missing data.**

---

## 8. Real Use Cases, One by One

| Application | Why UDP | What's built on top, if anything |
|---|---|---|
| DNS (Ch 66–69) | Single small request/response; retry is cheap | TCP fallback for large responses |
| VoIP (e.g. SIP/RTP) | Stale audio worse than missing audio | RTP sequence numbers + jitter buffers for the app to manage ordering/loss itself |
| Video conferencing | Same as VoIP, plus adaptive bitrate reacts faster without TCP's own retransmission logic in the way | Application-level forward error correction, frame dropping |
| Online multiplayer games | Newest state matters, not complete history | Custom lightweight reliability only for events that truly must arrive (e.g. "player fired a shot") |
| DHCP (Ch 55) | No established connection exists yet — device doesn't even have an IP address | DORA process handles retries itself |
| NTP (time sync) | Single tiny query/response, timing precision matters more than reliability | Client-side retry and averaging |
| QUIC / HTTP/3 (Ch 75) | Wants TCP-like reliability but *per stream*, not globally — building it fresh on UDP avoids head-of-line blocking across unrelated streams | Full custom reliability, congestion control, and encryption layered inside QUIC itself |
| SNMP (network monitoring, Ch 121) | Frequent small polls; occasional loss is tolerable | Polling naturally retries |

That QUIC row is worth pausing on: it's proof that "connectionless and unreliable" is a *starting point*, not a ceiling. QUIC takes UDP's bare 8-byte header and builds essentially all of TCP's reliability back on top of it — in user space, per application — specifically because doing so per-*stream* (rather than per-*connection*, as TCP does) solves a problem TCP's own design can't.

---

## 9. Broadcast and Multicast: What TCP Cannot Do

There's a structural reason UDP, not TCP, is the transport underneath DHCP (Chapter 55), certain routing protocol announcements, and IPTV-style video distribution, and it isn't just about avoiding handshake overhead — it's that **TCP's entire design assumes exactly two endpoints**, and broadcast/multicast delivery is fundamentally a one-to-many operation that a strictly two-party, connection-oriented protocol cannot express at all.

**Broadcast** means "deliver to every host on this local subnet" — the destination address `255.255.255.255` (or a subnet's directed broadcast address, Chapter 40) is not a specific machine at all, but an instruction to every receiver on the link. DHCP's DISCOVER message (Chapter 55) is broadcast for an unavoidable reason: a device that doesn't yet have an IP address has no way to know which server to address a unicast request to, so it has to shout to everyone on the local network and let any listening DHCP server reply.

**Multicast** is the more general, more efficient version: "deliver to every host that has explicitly joined this particular group," using a reserved range of addresses (224.0.0.0–239.255.255.255 for IPv4, Chapter 40) rather than every host on the subnet indiscriminately. A multicast video stream, for example, can be sent once by the source, and switches/routers along the way replicate it only toward the specific downstream links that actually have a subscribed listener — dramatically more efficient than sending an individual unicast copy to every viewer.

Why can TCP never do either? A TCP socket's entire identity — the 4-tuple from Chapter 57 — presupposes exactly one source and exactly one destination, and its reliability machinery (Chapters 59–61) requires a private, bidirectional stream of sequence numbers and acknowledgments negotiated between precisely two parties during a handshake. There is no coherent way to "acknowledge" a segment on behalf of an unbounded, dynamically-changing group of receivers, some of whom may not even exist yet when the data is sent, and no sensible way to retransmit a lost segment to "the one receiver who missed it" without knowing who that is in the first place. UDP's datagrams, by contrast, have no such assumption baked in — a UDP socket can simply be told "send this datagram to this multicast address" and the network layer (with IGMP managing group membership, briefly touched on again in Chapter 96's discussion of multicast-adjacent CDN techniques) handles fan-out, with no acknowledgment expected from anyone.

```
Unicast (ordinary UDP or TCP):    one sender ───▶ one receiver

Broadcast (UDP only):             one sender ───▶ every host on the subnet
                                              ─▶
                                              ─▶
                                              ─▶

Multicast (UDP only):             one sender ───▶ only the subscribed hosts
                                              ─▶ (not this one — didn't join)
                                              ─▶
```

---

## 10. UDP Code: A Minimal Client and Server

A complete, minimal UDP echo server and client in Go — deliberately short, because that's the point: there is no handshake to negotiate, no connection object to maintain state for, just send and receive datagrams.

```go
// server.go — listens on UDP port 9000 and echoes back whatever it receives
package main

import (
    "fmt"
    "net"
)

func main() {
    addr, _ := net.ResolveUDPAddr("udp", ":9000")
    conn, _ := net.ListenUDP("udp", addr)
    defer conn.Close()

    buf := make([]byte, 1472) // stay under typical Ethernet MTU minus headers
    for {
        n, clientAddr, err := conn.ReadFromUDP(buf)
        if err != nil {
            continue
        }
        fmt.Printf("received %d bytes from %s: %q\n", n, clientAddr, buf[:n])
        conn.WriteToUDP(buf[:n], clientAddr) // echo it straight back
    }
}
```

```go
// client.go — sends one datagram and waits (briefly) for the echo
package main

import (
    "fmt"
    "net"
    "time"
)

func main() {
    conn, _ := net.Dial("udp", "127.0.0.1:9000")
    defer conn.Close()

    conn.Write([]byte("hello"))

    conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    buf := make([]byte, 1472)
    n, err := conn.Read(buf)
    if err != nil {
        fmt.Println("no reply — datagram may have been lost, and nothing will retry it for us")
        return
    }
    fmt.Printf("server echoed: %q\n", buf[:n])
}
```

Notice the client's `SetReadDeadline` and the comment beside it: because UDP guarantees nothing, an application that cares about a reply has to build its own timeout and retry logic — exactly the kind of thing TCP does for you automatically, at the cost of the connection setup this code doesn't have to do at all. `net.Dial("udp", ...)` doesn't send any packets by itself; it just records the destination so `Write` doesn't need it repeated every call. No bytes cross the network until the first `Write`.

---

## 11. Seeing It Live

Capture the exchange above with `tcpdump` (a real hands-on experiment you can run):

```
$ sudo tcpdump -i lo -n udp port 9000 -X

12:00:01.001 IP 127.0.0.1.51342 > 127.0.0.1.9000: UDP, length 5
    0x0000:  4500 0021 0000 4000 4011 3ca0 7f00 0001
    0x0010:  7f00 0001 c8ae 2328 000d 0f2c 6865 6c6c
    0x0020:  6f
12:00:01.001 IP 127.0.0.1.9000 > 127.0.0.1.51342: UDP, length 5
    0x0000:  4500 0021 0000 4000 4011 3ca0 7f00 0001
    0x0010:  7f00 0001 2328 c8ae 000d 6237 6865 6c6c
    0x0020:  6f
```

Everything from Sections 3–5 is right there: at offset `0x0c` in each UDP payload you can find the 8-byte UDP header (`c8ae 2328` = source port 51342, destination port 9000 in the first packet; `2328 c8ae` reversed on the reply), followed immediately by `000d` (length 13) and a checksum, then `6865 6c6c 6f` — `hello` in ASCII — with no additional framing of any kind. There is no visible handshake preceding it, and no acknowledgment following it, because at the UDP layer, none exists.

---

## 12. Common Misconceptions

- **"UDP is unreliable, so it's low-quality or a fallback."** UDP is exactly as reliable as its designers intended: not reliable at all, by design, because reliability has a cost (waiting, buffering, complexity) that many applications shouldn't have to pay. "Unreliable" here is a precise technical term, not a value judgment.
- **"UDP doesn't have error detection at all."** It does — the checksum, computed in Section 5, detects (though cannot correct) corrupted data, covering the payload, the UDP header, and key IP-layer fields via the pseudo-header. What UDP lacks is a *response* to detected errors: a corrupted UDP datagram is simply dropped silently by the receiving kernel, with no automatic request for retransmission.
- **"UDP is always faster than TCP."** Usually lower-latency per packet (no handshake, no ordering delay), but "faster" depends entirely on what the application needs. If an application needs delivery guarantees, building them on top of UDP badly can end up slower and buggier than just using TCP, which already implements them correctly (Chapter 59 onward).
- **"A UDP 'connection' timing out means the server crashed."** UDP has no connection to time out in the first place — what actually happened is that a request datagram, or its reply, was lost somewhere along the path, and it's purely the *application's own* retry logic (if it has any) that decides what a timeout even means.
- **"UDP can't be used for anything that needs to be reliable."** QUIC (Chapter 75) is the clearest counterexample: it rebuilds reliable, ordered, per-stream delivery entirely on top of UDP, in user-space code, precisely because doing so gives it flexibility the kernel's fixed TCP implementation doesn't allow.

---

## 13. Production Notes

- **UDP flood attacks and amplification.** Because UDP requires no handshake, an attacker can send a stream of datagrams with a spoofed source IP address, and the recipient (or worse, a *third-party* server that replies to the spoofed address) will send responses to a victim who never asked for them. DNS and NTP have both been abused this way for amplification DDoS attacks (previewed here, covered in depth in Chapter 83) — a small forged UDP query can trigger a much larger response directed at the victim.
- **MTU and fragmentation.** Since UDP itself does nothing about datagrams larger than a single link's MTU, IP has to fragment them (or the OS avoids the problem by capping payloads under ~1,472 bytes on typical Ethernet). Fragmented UDP datagrams are more fragile — losing *one* IP fragment loses the entire UDP datagram, since there's no partial delivery. Applications sensitive to this (DNS with large responses, for instance) actively try to keep payloads under the fragmentation threshold.
- **Firewalls and NAT are harder on UDP.** Because there's no connection state (Chapter 41's NAT tables, and Chapter 84's stateful firewalls, both lean on visible handshake events to know when a "connection" starts and ends), UDP "connections" through NAT devices are tracked using timeouts and heuristics rather than a definitive close signal — this is part of why VoIP and gaming traffic sometimes struggles to traverse NAT cleanly, a problem protocols like STUN/TURN exist specifically to work around.
- **Checksum offload.** Modern NICs frequently compute UDP (and TCP) checksums in hardware rather than the kernel's CPU-based implementation, which matters for high-throughput applications like video streaming servers pushing many gigabits per second of UDP traffic.

---

## 14. What's Simplified Here

The worked checksum example in Section 5 uses a plain ASCII payload rather than a real protocol message purely for arithmetic clarity — real DNS, RTP, and QUIC payloads have their own internal binary structure that this chapter doesn't attempt to teach (DNS's real wire format is Chapters 66–69's job). This chapter also doesn't cover UDP-Lite (a rarely-used variant that allows partial checksum coverage for streaming media) or the details of how specific applications (RTP, QUIC) build reliability, sequencing, or congestion control on top of UDP — those are covered where relevant (Chapter 75 for QUIC) rather than here, since UDP itself contributes nothing to them.

---

## 15. Interview Questions & Model Answers

**Beginner: What fields does a UDP header contain, and how big is it?**

Eight bytes total, four 16-bit fields: source port, destination port, a length field covering the header plus payload, and a checksum. That's the entire header — no sequence numbers, no flags, no window size.

**Intermediate: Why does DNS use UDP instead of TCP for most queries?**

A typical DNS query and its response are each a single small packet, and if one is lost, the resolver can just resend the query — cheap, since there's no connection state to rebuild. Using TCP would mean paying for a three-way handshake (Chapter 59) for every single question, roughly tripling the round trips needed for the overwhelmingly common case of a small query and small response. DNS does fall back to TCP when a response is too large to fit in one UDP datagram, since large datagrams are more likely to fragment and lose the whole message if any one fragment is lost.

**Advanced: QUIC is built on top of UDP but implements reliable, ordered delivery. Doesn't that make UDP pointless in that case — why not just use TCP?**

TCP's reliability and ordering guarantees are enforced across the *entire connection*, not per logical stream within it — so if a browser opens one HTTP/2 connection carrying five independent resource downloads and one TCP segment is lost, all five streams stall until it's retransmitted, even though four of them have no actual dependency on the lost data (the head-of-line blocking problem, Chapters 74–75). QUIC needs reliability and ordering *scoped to each individual stream*, which the kernel's fixed TCP implementation can't provide — TCP's state machine has no concept of "stream" at all. Building on UDP gives QUIC a blank slate: no baked-in ordering or reliability semantics to fight against, so it can implement its own per-stream logic in user-space code that ships and evolves with the browser, rather than waiting years for operating system kernels worldwide to add a new TCP-layer feature.

---

## 16. Exercises

### Easy
1. List the four fields of the UDP header and their sizes.
2. Explain, without using the word "unreliable," what it means for UDP to be connectionless.
3. Name three applications that use UDP and, for each, state in one sentence why waiting for guaranteed in-order delivery would hurt more than it helps.

### Medium
4. A UDP datagram's Length field says 20. How many bytes of application payload does it contain?
5. Using the pseudo-header method from Section 5, explain what would happen to a UDP checksum if a NAT device rewrote the packet's source IP address but forgot to recompute the UDP checksum. Would the receiver detect a problem?
6. Why is UDP's checksum optional in IPv4 but mandatory in IPv6? (Hint: think about what other checksums exist at the IP layer in each version.)

### Hard
7. Compute, by hand or with a calculator, the UDP checksum for a datagram with source IP 10.0.0.1, destination IP 10.0.0.2, source port 5000, destination port 6000, and a 4-byte payload of `0x01 0x02 0x03 0x04`. Show the pseudo-header, the word-by-word sum, and the final one's complement.
8. A UDP-based multiplayer game sends 60 position updates per second, and design a lightweight application-level scheme (using only the game's own payload, not any new transport protocol) that lets a client detect when an update was lost, without ever requesting retransmission of stale data. What single extra field would you add to each payload, and how would the receiver use it?

---

## 17. Summary

| Term | Meaning |
|---|---|
| UDP | User Datagram Protocol — the minimal transport protocol: ports plus a checksum, nothing else |
| UDP header | 8 bytes: source port, destination port, length, checksum |
| Datagram | One self-contained, independently routed UDP message, with no relationship to others |
| Pseudo-header | IP-layer fields (source/dest IP, protocol, length) folded into the UDP checksum calculation |
| Connectionless | No handshake, no persistent per-conversation state kept by UDP itself |
| Head-of-line blocking | The wait a reliable, ordered protocol imposes when earlier data is missing — exactly what UDP avoids |
| Amplification attack | Abusing UDP's lack of a handshake to spoof requests that trigger large replies toward a victim |

UDP is the honest, minimal baseline: it tells applications exactly what the underlying network already provides, and nothing more. Most applications, though — file transfers, web pages, databases — need guarantees UDP was never designed to give. Chapter 59 starts from that gap and builds, from first principles, the protocol that closes it: TCP.
