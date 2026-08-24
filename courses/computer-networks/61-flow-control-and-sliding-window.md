# Chapter 61: Flow Control and the Sliding Window

> **"Being able to detect a lost byte, from Chapter 60, doesn't help you if the real problem is that the receiver's buffer overflowed before the byte was ever lost to the network at all."**

---

## Table of Contents

1. [A Different Problem: The Receiver, Not the Network](#1-a-different-problem-the-receiver-not-the-network)
2. [The Naive Fix: Stop-and-Wait](#2-the-naive-fix-stop-and-wait)
3. [Why Stop-and-Wait Wastes Bandwidth](#3-why-stop-and-wait-wastes-bandwidth)
4. [The Real Solution: The Sliding Window](#4-the-real-solution-the-sliding-window)
5. [The Window Field in the TCP Header](#5-the-window-field-in-the-tcp-header)
6. [A Full Worked Example Across Several Round Trips](#6-a-full-worked-example-across-several-round-trips)
7. [Zero Window and the Persist Timer](#7-zero-window-and-the-persist-timer)
8. [Window Scaling](#8-window-scaling)
9. [Bandwidth-Delay Product](#9-bandwidth-delay-product)
10. [Silly Window Syndrome and Nagle's Algorithm](#10-silly-window-syndrome-and-nagles-algorithm)
11. [Seeing It Live](#11-seeing-it-live)
12. [Code: A Minimal Sliding Window Simulation](#12-code-a-minimal-sliding-window-simulation)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes](#14-production-notes)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#interview-questions--model-answers)
17. [Exercises](#exercises)
18. [Summary](#summary)
19. [Bridge to Chapter 62](#19-bridge-to-chapter-62)

---

## 1. A Different Problem: The Receiver, Not the Network

Chapter 60 solved a problem caused by the *network*: packets get lost or reordered somewhere between sender and receiver, and sequence numbers plus cumulative ACKs let the receiver detect exactly what's missing and the sender know exactly what to resend.

This chapter is about a completely different failure, caused not by the network at all, but by the **receiver itself**. Every TCP receiver holds incoming data in a finite buffer in memory until the application actually reads it out. If the sender transmits data faster than the receiving application consumes it, that buffer fills up — and once it's full, any further arriving data has nowhere to go and must simply be discarded by the receiver's own operating system, with no network failure involved at all.

Concretely: imagine a fast server on a data-center-grade network link sending data to a phone on a slow, congested cellular connection, where the app reading the data is itself busy doing something else (parsing, writing to disk, waiting on a slow database) and only reads a trickle of bytes out of its socket buffer every few hundred milliseconds. If the server just sent as fast as its own network link allowed, the phone's buffer would fill up in a fraction of a second, and every byte after that point would be dropped the instant it arrived — not because a router lost it, but because there was simply nowhere on the receiving end to put it.

This is exactly the problem **flow control** exists to solve, and it's worth being precise about the distinction this chapter is setting up: flow control asks "can the *receiver* keep up," a question about one machine's buffer capacity. Chapter 62's congestion control asks an entirely different question — "can the *network* keep up" — and conflating the two, as Section 19 will stress again at the end of this chapter, is one of the most common sources of confusion in how people talk about TCP.

---

## 2. The Naive Fix: Stop-and-Wait

The simplest possible fix, and worth trying before reaching for anything cleverer: what if the sender simply sent **one segment**, then waited for an acknowledgment before sending the next one? If the receiver's buffer can hold at least one segment's worth of data, it can never overflow, because the sender never has more than one segment "in flight" — sent but not yet acknowledged — at any given moment.

```
Sender                                   Receiver
  │  segment 1  ──────────────────────▶  │
  │                                       │  (buffer: 1 segment, fine)
  │◀────────────────────────  ack 1      │
  │  segment 2  ──────────────────────▶  │
  │                                       │
  │◀────────────────────────  ack 2      │
  │  segment 3  ──────────────────────▶  │
  │                                       │
```

This is called **stop-and-wait**, and it does correctly solve the overflow problem: the receiver's buffer requirement never exceeds one segment's worth of data, no matter how slow the application reading it is. It's a real, historically used approach (early, simple data-link protocols used exactly this scheme). But applied to TCP over a real network, it has a devastating cost, worked out precisely in Section 3.

---

## 3. Why Stop-and-Wait Wastes Bandwidth

Consider a connection with a round-trip time of 100ms (not unreasonable for a cross-country or transoceanic connection) and a link capable of carrying 100 Mbps, carrying 1,460-byte segments (a typical MSS, Chapter 59, Section 7).

Under strict stop-and-wait, the sender transmits one segment, then must wait a full RTT before the ACK returns and it's allowed to send the next one:

```
Time to send one 1460-byte segment:  ~0.1ms   (negligible on a 100 Mbps link)
Time to wait for its ACK:             100ms   (the full RTT)

Effective throughput = 1460 bytes / 0.1001 seconds ≈ 14,585 bytes/sec
                      ≈ 0.117 Mbps

Compare to the link's actual capacity: 100 Mbps

Utilization: 0.117 / 100 ≈ 0.117%
```

The link is sitting almost entirely idle — 99.88% of its capacity is wasted, not because the network can't carry more, but because the sender has voluntarily throttled itself to "one segment, then wait an entire round trip." The higher the RTT, and the higher the link's actual bandwidth, the *worse* this gets — a phenomenon directly related to the bandwidth-delay product covered properly in Section 9. Stop-and-wait solves the overflow problem completely and is completely unacceptable in practice. The real solution needs to let *multiple* segments be in flight simultaneously — but only as many as the receiver has actually promised it can hold.

---

## 4. The Real Solution: The Sliding Window

TCP's actual mechanism, called the **sliding window**, keeps stop-and-wait's core safety property — never send more than the receiver can currently hold — while dropping its fatal restriction to exactly one segment at a time.

The idea: the receiver continuously tells the sender, in every ACK it sends, exactly how much *additional* buffer space it currently has available. This value is called the **receive window (rwnd)**. The sender is then free to have up to that many bytes of data outstanding — sent, but not yet acknowledged — at any given moment, sending new segments continuously as long as the total unacknowledged data stays within that limit, without waiting for each individual ACK first.

```
                    ┌─────────────── receive window (rwnd) ───────────────┐
Sent & ACKed        │  Sent, NOT yet ACKed (in flight)     │  Not sent yet │
────────────────────┼───────────────────────────────────────┼─────────────
1000              2000                                    5000          ...
                    ▲                                        ▲
              window's left edge                      window's right edge
           (advances as new ACKs arrive)         (left edge + current rwnd)
```

As ACKs arrive confirming more data, the window's left edge slides forward — freeing up room for more data to be sent — which is exactly where the mechanism gets its name: the window of "currently permitted to be in flight" data slides along the byte stream as the connection progresses. This lets a sender keep the network link continuously full of data (bounded only by the receiver's actual buffer capacity), rather than idling for a full RTT after every single segment.

> **Intuitive analogy:** think of a conveyor belt loading boxes onto a truck, with a worker at the far end unloading them. If the loader is only allowed to put one box on the belt and then must wait until the worker confirms it's been unloaded before adding another, the belt spends nearly all its time empty. If instead the worker simply shouts back "I've got room for 10 more boxes" and the loader keeps the belt fed up to that limit, updating as boxes get unloaded and the announced capacity changes, the belt stays productively full — but never overloaded past what the worker actually said they could handle.
>
> **Where the analogy breaks:** a real conveyor belt has one lane; TCP's sliding window is really tracking a *range of byte positions*, not physical objects moving past a point — the "boxes" already delivered but not yet read by the application are still sitting in the receiver's buffer, counted against the receiver's stated free capacity, until the application actually consumes them and the window can grow again.

---

## 5. The Window Field in the TCP Header

The receive window isn't a separate control message — it rides in every single TCP segment sent, including pure ACKs with no data of their own, in the 16-bit **Window** field of the TCP header (the same header sketched in Chapter 59, Section 7, and covered field-by-field in full in Chapter 65):

```
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|Offset | Rsvd  |U A P R S F|                                  |
|       |       |R C S S Y I|          Window Size             |
|       |       |G K H T N N|                                  |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

A 16-bit field caps the raw advertised window at 65,535 bytes (2^16 − 1) — a real limitation, and the exact problem Section 8's window scaling exists to lift. Every time the receiving application reads data out of its socket buffer, freeing up space, the *next* ACK the receiver sends reflects the new, larger available window — this is how the receiver continuously communicates its current capacity, in real time, using the exact same segments it was already sending for ordinary acknowledgment purposes.

---

## 6. A Full Worked Example Across Several Round Trips

Let's trace a complete, numeric scenario: a sender transmitting data to a receiver whose application periodically gets busy and stops reading, then catches up — showing the window genuinely shrink and then grow again, exactly as the chapter plan requires.

**Setup:** MSS = 1000 bytes (round numbers, for readability). Sender's data starts at sequence number 10001. Receiver's buffer capacity is 5000 bytes total.

**Round trip 1 — receiver's application is keeping up:**

```
Receiver's buffer: empty (5000 bytes free) → advertises rwnd = 5000

Sender transmits, without waiting for individual ACKs, up to 5000 bytes:
  seg A: seq=10001, 1000 bytes (10001–11000)
  seg B: seq=11001, 1000 bytes (11001–12000)
  seg C: seq=12001, 1000 bytes (12001–13000)
  seg D: seq=13001, 1000 bytes (13001–14000)
  seg E: seq=14001, 1000 bytes (14001–15000)
  ── sender now has 5000 bytes in flight, exactly at the rwnd limit,
     and must stop sending new data until the window advances ──

Receiver's application reads all 5000 bytes promptly.
Receiver ACKs: ack=15001, rwnd=5000
  ("I have everything through 15000, and my buffer is now
    completely empty again — you may send up to 5000 more")
```

**Round trip 2 — receiver's application gets busy (stops reading):**

```
Sender, seeing rwnd=5000 still available, sends another full window:
  seg F: seq=15001, 1000 bytes
  seg G: seq=16001, 1000 bytes
  seg H: seq=17001, 1000 bytes
  seg I: seq=18001, 1000 bytes
  seg J: seq=19001, 1000 bytes

This time, the receiver's APPLICATION does not read any of it —
it's busy doing something else (e.g. writing a previous batch to
a slow disk). The data sits in the kernel's receive buffer.

Receiver still ACKs each segment as it arrives (the DATA is received
fine — flow control is about the buffer, not about loss), but the
buffer is now full: rwnd = 5000 - 5000 = 0

Receiver's last ACK: ack=20001, rwnd=0
  ("I've received everything through 20000 — but my buffer is now
    completely full. Do not send anything more until I say otherwise.")
```

**Round trip 3 — the window has shrunk to zero; sender must stop:**

```
Sender has NO permission to send any further data, even though it
may have more queued and ready to go. This is the "shrinking window"
the chapter plan calls for — and it happened purely because the
receiving application fell behind, with no packet loss anywhere.

Sender enters the zero-window state described fully in Section 7,
periodically probing rather than simply giving up.
```

**Round trip 4 — the application catches up; the window grows again:**

```
Receiver's application finally reads all 5000 buffered bytes.
Buffer is empty again → rwnd = 5000

Receiver sends a fresh ACK (prompted by a window probe, Section 7):
  ack=20001, rwnd=5000
  ("Same data position as before — I haven't received anything new —
    but I now have room again. You may resume sending.")

Sender resumes:
  seg K: seq=20001, 1000 bytes
  seg L: seq=21001, 1000 bytes
  ... up to 5000 bytes again, and the cycle repeats.
```

This four-round-trip trace is the entire mechanism, end to end: the window isn't a fixed connection property — it's a live, continuously-updated number reflecting one very specific, very local fact (how much free buffer space the receiving application has actually created by reading data out), and the sender is contractually bound to respect whatever value it was most recently told, even down to zero.

---

## 7. Zero Window and the Persist Timer

Section 6 left the sender stalled with `rwnd=0`. This creates a subtle problem worth working through: once the window hits zero, the sender stops transmitting data — but the *next* ACK that would tell it the window has grown again can only be sent by the receiver, and the receiver has no new data to acknowledge, so it has no natural reason to send anything at all. If that next informative ACK (announcing a restored window) is ever lost, sender and receiver could **deadlock indefinitely** — the sender waiting forever for a "window open" signal that was sent but never arrived, and the receiver having no reason to know its update didn't get through.

TCP's fix is the **persist timer**: once a sender's window has been zero for a while, it periodically sends a small **window probe** segment (typically containing just one byte of data) purely to provoke a fresh ACK from the receiver, carrying the receiver's current window value. If the window is still zero, the receiver just re-confirms zero and the sender waits and probes again, backing off the probe interval over time; the moment the window has actually opened up, the probe's ACK carries the good news and normal transmission resumes — closing the exact deadlock risk described above.

```
Sender (window has been 0 for a while):
   ──── 1-byte probe ────▶
                            Receiver: still full → ack=20001, rwnd=0
   ... wait, backoff ...
   ──── 1-byte probe ────▶
                            Receiver: freed up → ack=20001, rwnd=5000
Sender resumes normal transmission.
```

---

## 8. Window Scaling

Section 5 flagged the ceiling this section now confronts directly: the Window field is only 16 bits, capping the advertised window at 65,535 bytes. On a modern, high-bandwidth, high-latency connection, that ceiling becomes a serious real limitation — Section 9 works out precisely how serious with the bandwidth-delay product.

**Window scaling** (RFC 7323) fixes this without changing the 16-bit field itself — which can't be resized without breaking every existing TCP implementation on Earth. Instead, both sides negotiate a **scale factor** during the handshake (as a TCP option, mentioned in Chapter 59, Section 7), and from that point on, the *actual* window is the value in the 16-bit field, left-shifted by that factor:

```
Actual window = (16-bit Window field value) × 2^(scale factor)

Example: Window field = 4000, scale factor = 7
Actual window = 4000 × 2^7 = 4000 × 128 = 512,000 bytes
```

The scale factor can be up to 14, which extends the maximum representable window from 65,535 bytes all the way to `65,535 × 2^14 ≈ 1,073,725,440 bytes` — roughly one gigabyte, comfortably enough for even extremely high-bandwidth, high-latency ("long fat network") paths. Since the scale factor is fixed for the lifetime of the connection (negotiated only during the handshake, when both sides' capabilities are exchanged) and both endpoints must support the option for it to take effect, an older or non-participating stack simply sees the plain, unscaled 16-bit value and behaves exactly as it always has — a clean backward-compatible extension.

---

## 9. Bandwidth-Delay Product

Here is the precise reasoning for why window scaling was necessary at all, not just a nice-to-have. The **bandwidth-delay product (BDP)** answers the question "how much data needs to be in flight, unacknowledged, at any given moment to keep a link of a given speed and round-trip time continuously full?"

```
BDP = bandwidth × round-trip time
```

Worked example, a genuinely realistic modern case: a 1 Gbps link with a 100ms RTT (plausible for a transcontinental connection):

```
BDP = 1,000,000,000 bits/second × 0.1 seconds
    = 100,000,000 bits
    = 12,500,000 bytes
    ≈ 11.92 MiB
```

To keep that link fully utilized — the entire point of the sliding window mechanism from Section 4 — the sender needs roughly 12.5 million bytes in flight at once. But the unscaled 16-bit window field caps that at 65,535 bytes — **about 190 times too small**. Without window scaling, a connection over exactly this kind of high-bandwidth, high-latency path could never use more than a tiny fraction of the available bandwidth, no matter how fast the underlying network actually was, because the receiver is structurally incapable of advertising a large enough window to keep the pipe full. This is precisely the scenario window scaling (Section 8) was standardized to fix, and it's why virtually every modern OS negotiates it by default on every connection.

---

## 10. Silly Window Syndrome and Nagle's Algorithm

A related, smaller pathology worth knowing by name: **Silly Window Syndrome** occurs when a receiver's application reads data out of its buffer in tiny increments, and the receiver dutifully advertises each small freed sliver of space as soon as it appears — say, announcing a window of just a few bytes at a time. A sender honoring that literally would send a stream of tiny, mostly-header segments (a 20+-byte TCP header wrapped around a handful of payload bytes), wasting bandwidth on overhead vastly out of proportion to actual data carried.

The standard fixes are conventions, not new header fields: a well-behaved receiver should wait until a *meaningfully large* amount of buffer space has freed up before advertising a bigger window (rather than announcing every single byte the instant it's free), and a well-behaved sender applies **Nagle's algorithm** — holding small outgoing segments briefly to see if more application data is about to be written, so it can send one reasonably-sized segment instead of several tiny ones. Nagle's algorithm is occasionally disabled deliberately (the `TCP_NODELAY` socket option) for latency-sensitive applications, like real-time trading systems or interactive game state, where waiting even a few milliseconds to batch small writes is worse than the minor efficiency loss from sending them immediately.

---

## 11. Seeing It Live

The window field is visible in every packet capture. Deliberately slow down a receiving application's reads (for instance, a simple script that reads one byte and sleeps) while a fast sender pushes data, and capture with `tcpdump`:

```
$ sudo tcpdump -i lo -n tcp port 9000

12:00:01.001  IP receiver.9000 > sender.51342: Flags [.], ack 20001, win 5760, length 0
12:00:01.050  IP receiver.9000 > sender.51342: Flags [.], ack 20001, win 2048, length 0
12:00:01.090  IP receiver.9000 > sender.51342: Flags [.], ack 20001, win 0,    length 0
12:00:03.200  IP receiver.9000 > sender.51342: Flags [.], ack 20001, win 5760, length 0
```

Note the `ack 20001` staying fixed while `win` shrinks across the first three lines — exactly Section 6's trace: no new data is being acknowledged (the receiver isn't losing anything or receiving anything new), but its available buffer space is visibly draining to zero as the application falls behind, then recovering once it catches up. `ss -tin` on Linux exposes the currently negotiated window scale factor directly:

```
$ ss -tin
ESTAB 0 0 192.168.1.42:51342 142.250.80.46:443
    cubic wscale:7,7 rto:204 rtt:52.5/12.5 ...
```

`wscale:7,7` means both sides negotiated a scale factor of 7 — directly matching Section 8's worked example.

---

## 12. Code: A Minimal Sliding Window Simulation

A simplified simulation capturing the shrink-and-grow behavior from Section 6:

```go
type Receiver struct {
    bufferCap  int
    buffered   int // bytes currently held, unread by the application
}

func (r *Receiver) AdvertisedWindow() int {
    return r.bufferCap - r.buffered
}

func (r *Receiver) OnDataArrives(n int) {
    r.buffered += n // fills the buffer; ACK carries the new (smaller) window
}

func (r *Receiver) OnApplicationReads(n int) {
    r.buffered -= n // frees space; NEXT ACK carries the new (larger) window
}

type Sender struct {
    inFlight int
}

// CanSend enforces the core sliding-window safety rule: never exceed
// what the receiver most recently advertised.
func (s *Sender) CanSend(amount, lastKnownRwnd int) bool {
    return s.inFlight+amount <= lastKnownRwnd
}

func (s *Sender) OnSegmentSent(n int)     { s.inFlight += n }
func (s *Sender) OnCumulativeAck(n int)   { s.inFlight -= n } // window "slides"
```

Running this with the exact numbers from Section 6 — `bufferCap = 5000`, five 1000-byte arrivals with no reads in between, then a full 5000-byte read — reproduces the shrink-to-zero and grow-back-to-5000 sequence exactly.

---

## 13. Common Misconceptions

- **"Flow control and congestion control are the same thing."** They solve different problems with different information: flow control is bounded by the *receiver's* buffer, using a value the receiver itself reports (rwnd); congestion control (Chapter 62) is bounded by the *network's* capacity, which nothing directly reports and has to be inferred. Section 19 restates this distinction as the explicit bridge to the next chapter, because conflating them is genuinely one of the most common errors in casual descriptions of TCP.
- **"A zero window means the connection is broken."** It means the receiver's buffer is temporarily full — a completely normal, self-correcting condition, resolved automatically once the application reads more data and the persist timer's next probe reveals the recovered window (Section 7).
- **"The sliding window is about which segments have been acknowledged."** That's Chapter 60's cumulative-ACK mechanism. The sliding window is about *how much unacknowledged data is allowed to exist at once* — a distinct, complementary concept layered on top of sequence numbers and ACKs, not a replacement for them.
- **"A bigger advertised window always means faster transfer."** Only up to the bandwidth-delay product (Section 9) — advertising more window than the BDP requires provides no additional throughput benefit, since the link is already kept continuously full; it mostly just risks larger receiver memory consumption per connection.
- **"Window scaling changes the TCP header format."** It doesn't touch the 16-bit Window field's size at all — it's purely an interpretation agreement (a multiplier) negotiated as a header *option* during the handshake, layered on top of the unchanged field.

---

## 14. Production Notes

- **Socket buffer sizing is a real tuning knob.** Operating systems expose settings (e.g., Linux's `net.core.rmem_max`, `net.ipv4.tcp_rmem`) controlling the maximum receive buffer size a connection can grow to — directly capping the largest window a receiver can ever advertise, and therefore the maximum achievable throughput on high-BDP paths (Section 9), independent of how fast the network itself actually is.
- **Autotuning** is standard in modern kernels: rather than a fixed buffer size, the OS dynamically grows a connection's receive buffer (and therefore its advertised window) based on observed throughput and available system memory, removing the need for most applications to manually tune these values.
- **Long-haul, high-bandwidth links (satellite, transcontinental fiber) are the textbook case where window scaling and buffer tuning genuinely matter** — a misconfigured or scaling-disabled connection across such a path can silently cap throughput far below the link's real capacity, a common real-world performance debugging scenario (revisited in Chapter 120's measurement techniques).
- **Zero-window events are a monitorable signal** of an application-side bottleneck (slow consumer) rather than a network problem — distinguishing "the network is congested" (Chapter 62's territory) from "my own application isn't reading fast enough" (this chapter's territory) is a genuinely useful diagnostic split when debugging a slow service.

---

## 15. What's Simplified Here

This chapter uses small, round numbers (MSS = 1000, buffer = 5000) for the worked example in Section 6 purely for readability — real MSS values are typically 1460 bytes and real buffer sizes are often tens of megabytes on modern systems with autotuning enabled. The persist timer's exact backoff schedule and probe segment format are described only at the concept level here, not the precise byte layout. This chapter also does not cover the receiver-side congestion window interaction (the sender's actual send limit is the *smaller* of the receive window from this chapter and the congestion window from Chapter 62 — a detail that only makes sense once both concepts exist, and is picked up properly at the start of the next chapter).

---

## Interview Questions & Model Answers

**Beginner: What problem does the TCP sliding window solve, and how is it different from the problem sequence numbers and ACKs solve?**

It prevents a fast sender from overwhelming a slow receiver's buffer — a purely local, receiver-side capacity issue, not a network reliability issue. Sequence numbers and cumulative ACKs (Chapter 60) detect and recover from data that the *network* lost, duplicated, or reordered; the sliding window instead limits how much unacknowledged data may exist at once, based on how much buffer space the receiver has actually reported it can currently hold, so data is never dropped simply for lack of somewhere to put it.

**Intermediate: Why is stop-and-wait an inadequate flow-control strategy for a real network connection, even though it technically prevents receiver buffer overflow?**

Because it forces the sender to wait a full round-trip time after every single segment before sending the next one, regardless of how much spare capacity the network link and the receiver's buffer actually have. On a connection with meaningful latency, this collapses effective throughput to a tiny fraction of the link's real capacity — a 100 Mbps link with 100ms RTT sending 1460-byte segments under strict stop-and-wait achieves well under 1% utilization. The sliding window fixes this by allowing many segments to be in flight simultaneously, up to the receiver's advertised capacity, rather than one at a time.

**Advanced: A high-bandwidth, high-latency connection is achieving far less throughput than the link should allow, even though there's no packet loss. What's the likely cause, and how would you confirm it?**

The likely cause is the connection's usable window being capped below its bandwidth-delay product — either because window scaling wasn't negotiated (an old stack, a middlebox stripping the option, or a misconfiguration disabling it), or because a socket buffer size limit on one side is preventing the OS from advertising a window large enough to keep the pipe full, even with scaling available. Confirm it by computing the BDP (bandwidth × RTT) and comparing it against the connection's actual observed window size (visible via `ss -tin`'s window scale factor and effective window, or from a packet capture's advertised `win` values multiplied by the negotiated scale); if the usable window is meaningfully smaller than the BDP, the receiver is structurally prevented from ever fully utilizing the link, independent of anything happening at the congestion-control layer.

---

## Exercises

### Easy
1. In one sentence, explain the difference between what "flow control" limits and what "congestion control" (previewed for Chapter 62) will limit.
2. A receiver advertises a window of 0. What must happen before the sender can transmit new data again?
3. Explain why stop-and-wait guarantees no receiver buffer overflow but is unacceptable for real network throughput.

### Medium
4. A link has a bandwidth of 500 Mbps and an RTT of 60ms. Compute the bandwidth-delay product in bytes.
5. A receiver's Window field shows 8000 and the negotiated scale factor is 5. What is the actual advertised window in bytes?
6. Explain, using Section 6's four-round-trip trace as a template, a scenario where the window would need to shrink and grow *twice* within one connection.

### Hard
7. A connection has a measured BDP of 25 MB but the receiving host's maximum socket buffer size is capped by OS configuration at 4 MB. Explain precisely what throughput ceiling this creates, independent of the actual link speed or window scaling being correctly negotiated, and describe the one configuration change that would remove this ceiling.
8. Design a rule (in prose or pseudocode) for how often a sender in the zero-window state should send persist-timer probes, and justify a trade-off between probing too often (wasting bandwidth, adding load to an already-struggling receiver) and too rarely (leaving the connection stalled longer than necessary once the window actually reopens).

---

## Summary

| Term | Meaning |
|---|---|
| Flow control | Preventing a sender from overwhelming the receiver's buffer — a receiver-capacity problem |
| Receive window (rwnd) | The amount of additional data the receiver currently has buffer space for; advertised in every segment |
| Sliding window | Allowing multiple segments in flight, up to rwnd, sliding forward as ACKs arrive |
| Stop-and-wait | The naive one-segment-at-a-time approach; safe but wastes almost all available bandwidth on any real-latency link |
| Zero window / persist timer | The state when rwnd hits 0, and the periodic probing mechanism that safely detects when it reopens |
| Window scaling | A negotiated multiplier extending the 16-bit window field's effective range up to roughly 1 GB |
| Bandwidth-delay product (BDP) | The data that must be in flight to keep a link of given speed and RTT continuously full |
| Silly Window Syndrome / Nagle's algorithm | The pathology of advertising/sending tiny increments, and the convention-based fixes for it |

---

## 19. Bridge to Chapter 62

Everything in this chapter has been about one very specific constraint: how much data the *receiving application* can currently absorb. Notice what was never once mentioned in any of Section 6's worked example, or Section 9's bandwidth-delay product calculation: nothing about routers in the middle of the path running low on buffer space, links becoming saturated, or the network itself struggling under load. Flow control assumes the network can deliver whatever the receiver says it can accept — its only question is "can the receiver keep up?"

But the network is not infinite, and it is shared by every other connection on Earth competing for the same routers and links. A sender could have a receiver advertising a generous window, a receiver reading data instantly, buffers as large as anyone likes — and still be sending far more data than the actual network path between them can carry, causing routers along the way to queue up and eventually drop packets under their own load, entirely independent of anything the receiver reported. That is a different problem, with a different answer, from a different piece of information nobody explicitly reports at all: **can the network itself keep up?** Chapter 62 is where TCP answers it.
