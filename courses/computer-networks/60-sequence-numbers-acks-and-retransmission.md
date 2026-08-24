# Chapter 60: Sequence Numbers, Acknowledgments, and Retransmission

> **"TCP doesn't count packets. It counts bytes. That one design choice is why a lost segment can be replaced by a differently-sized one, and why 'how much data arrived' is always a precise, unambiguous number."**

---

## Table of Contents

1. [Where the Handshake Left Off](#1-where-the-handshake-left-off)
2. [Why Number Bytes, Not Packets](#2-why-number-bytes-not-packets)
3. [Cumulative Acknowledgments](#3-cumulative-acknowledgments)
4. [A Worked Example: Clean Delivery](#4-a-worked-example-clean-delivery)
5. [A Worked Example: A Lost Segment](#5-a-worked-example-a-lost-segment)
6. [Duplicate ACKs](#6-duplicate-acks)
7. [Retransmission Timeout (RTO)](#7-retransmission-timeout-rto)
8. [Computing RTO: The Jacobson/Karn Algorithm](#8-computing-rto-the-jacobsonkarn-algorithm)
9. [Selective Acknowledgment (SACK)](#9-selective-acknowledgment-sack)
10. [Seeing It Live](#10-seeing-it-live)
11. [Code: A Minimal Retransmission Timer](#11-code-a-minimal-retransmission-timer)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Production Notes](#13-production-notes)
14. [What's Simplified Here](#14-whats-simplified-here)
15. [Interview Questions & Model Answers](#interview-questions--model-answers)
16. [Exercises](#exercises)
17. [Summary](#summary)

---

## 1. Where the Handshake Left Off

Chapter 59 ended with both sides of a TCP connection in the `ESTABLISHED` state, each holding a confirmed, agreed-upon starting number for the other's byte stream — the client's ISN plus one, and the server's ISN plus one. That agreement was the entire point of the handshake, and this chapter is where it starts paying off.

Now the real problem from Chapter 59's opening returns in full force: data is about to flow over a network that can still lose, duplicate, and reorder anything, and both sides need a mechanism to detect exactly what got through and exactly what didn't — precisely enough that only the missing pieces get resent, never the whole stream.

---

## 2. Why Number Bytes, Not Packets

Here's a design choice worth pausing on, because it isn't the obvious one: TCP does not number the *segments* it sends (segment 1, segment 2, segment 3...). It numbers every individual **byte** of the data stream, and each segment's sequence number is just the position of its *first* byte in that stream.

Why not simply count segments? Walk through what breaks if you did. Suppose segment 1 carries bytes and gets lost, and the sender needs to retransmit it — but suppose the sender also wants to be free to retransmit it as a *differently sized* chunk (maybe combined with the next segment, or split in two, for efficiency). If receivers tracked "I'm missing segment 1," but the retransmission arrives packaged as a different segment altogether, the receiver has no way to recognize it as the same missing data. Packet-counting only works if every retransmission is byte-for-byte, size-for-size identical to the original — a needless, brittle constraint.

Byte-numbering sidesteps this completely. The receiver doesn't track "which segment numbers have I seen" — it tracks "which *byte positions* have I received," which is a much more flexible fact: it stays true no matter how the sender chooses to chop up (or re-chop up, on retransmission) the underlying stream into segments. This is also precisely why TCP is often described as delivering a **byte stream**, not a sequence of discrete messages — the application reads a continuous flow of bytes with no length markers of its own, and TCP internally is free to combine, split, or resize the segments that carry them however it finds convenient.

```
Application writes:  "GET / HTTP/1.1\r\n..."   (say, 517 bytes total)

TCP is free to send this as:
  one 517-byte segment, or
  two segments of 300 + 217 bytes, or
  517 individual 1-byte segments (wildly inefficient, but not incorrect)

Whatever the split, byte 1 of the data is sequence number ISN+1,
byte 2 is ISN+2, ... byte 517 is ISN+517 — regardless of which
segment each byte happens to travel in.
```

---

## 3. Cumulative Acknowledgments

TCP's core acknowledgment rule, and the single most important sentence in this chapter: **the acknowledgment number in a TCP segment means "I have correctly received every byte up to, but not including, this number — send me this one next."**

This is called a **cumulative acknowledgment** because it doesn't just say "I got segment X" — it says "I have an unbroken run of bytes starting from the beginning of this connection all the way up to this exact point, with no gaps." A single ACK, in other words, silently confirms *everything that came before it too*, not just the most recent segment.

The direct, useful consequence: if a sender has transmitted bytes 1001 through 3000 in three segments, and the receiver acknowledges `ack=3001`, the sender knows with certainty that all three segments arrived intact — one ACK confirmed all of them, because a cumulative ACK can only advance past a byte position if every byte before it, with no gaps, has actually arrived.

---

## 4. A Worked Example: Clean Delivery

Picking up directly from Chapter 59's worked handshake — client ISN 1000, server ISN 5000, connection `ESTABLISHED` with the client's next byte at 1001 and the server's at 5001 — suppose the client now sends 1,500 bytes of data, transmitted as three 500-byte segments.

```
Segment 1: seq=1001, 500 bytes (covers bytes 1001–1500)
Segment 2: seq=1501, 500 bytes (covers bytes 1501–2000)
Segment 3: seq=2001, 500 bytes (covers bytes 2001–2500)

All three arrive successfully, in order.

Server sends back:  ack=2501
                     ("I have everything through byte 2500;
                       send me byte 2501 next")
```

Note that the server doesn't have to send three separate ACKs, one per segment — a single ACK carrying `2501` is sufficient to confirm all three segments at once, precisely because of the cumulative property from Section 3. In practice, real TCP stacks often *do* send more frequent ACKs than strictly required (partly to help the sender's flow control and RTT estimation, covered in Chapter 61 and Section 8 below), but the *minimum* information needed to confirm receipt of all three segments is exactly this one number.

---

## 5. A Worked Example: A Lost Segment

Now the interesting case — segment 2 (bytes 1501–2000) is lost in transit, while segments 1 and 3 both arrive:

```
Segment 1: seq=1001, 500 bytes ──────────▶  arrives — server now has 1001–1500
Segment 2: seq=1501, 500 bytes ──── ✕ ────  LOST
Segment 3: seq=2001, 500 bytes ──────────▶  arrives, but out of order —
                                              server is missing 1501–2000,
                                              so it CANNOT advance its
                                              cumulative ACK past 1501

Server's ACKs:
  After segment 1: ack=1501   ("I have through 1500, send me 1501 next")
  After segment 3: ack=1501   ("still — I have through 1500 only;
                                 segment 3 is buffered, out of order,
                                 waiting on the gap at 1501 to close")
```

This is the crucial mechanic: **the cumulative ACK cannot skip over a gap**, even though the receiver has *already received* bytes 2001–2500 (segment 3) and is holding them in a buffer. It must keep re-announcing `ack=1501` — the position right at the start of the missing segment — until segment 2 actually arrives and closes the hole. Only then can the receiver's cumulative ACK jump all the way to 2501 in one step, since bytes 1501–2000 (finally received) plus the already-buffered 2001–2500 together form one unbroken run.

```
Segment 2 retransmitted: seq=1501, 500 bytes ──────────▶  arrives
Server's ACK jumps immediately: ack=2501
  ("Now I have an unbroken run through 2500 — the gap is closed,
    and I can credit the already-buffered segment 3 in the same step")
```

This single worked example is the reason TCP is called a **reliable, in-order** protocol: the receiver never hands the application bytes 2001–2500 until 1501–2000 has also arrived, preserving order even though the underlying network delivered the data out of order — and the sender is told, unambiguously, exactly which byte to resend, without needing to guess.

---

## 6. Duplicate ACKs

Look again at Section 5's sequence of ACKs: `1501`, then `1501` again after segment 3 arrives. That repeated identical ACK value is called a **duplicate ACK**, and it's a second, faster signal (beyond the timeout mechanism in Section 7) that something specific went wrong.

Here's the reasoning a sender can apply: a single ACK confirming `1501` is completely normal — it's just the ordinary acknowledgment of segment 1. But a *second* ACK arriving with the exact same value strongly suggests that something arrived at the receiver *after* the gap at 1501 (in this case, segment 3), and the receiver is essentially saying, over and over, "still waiting on 1501, but I'm getting later data around it." That pattern — the same cumulative ACK repeating — is a much stronger and faster signal of "something specific is missing" than simply not hearing back at all.

TCP's actual convention (elaborated fully in Chapter 63) is to treat **three duplicate ACKs** (the original ACK plus three more identical repeats) as confident enough evidence of loss to retransmit immediately, without waiting for a timeout — this is called **fast retransmit**, and it exists because timeouts (Section 7) are deliberately conservative and can be far slower than necessary when the evidence of loss is already sitting in the ACK stream. This chapter introduces the concept only far enough to explain *why* duplicate ACKs matter for retransmission at all; Chapter 63 covers the exact threshold, fast recovery, and the congestion-window interaction in full.

---

## 7. Retransmission Timeout (RTO)

Duplicate ACKs only help when *later* data manages to arrive and reveal the gap, as in Section 5. What about the case where a segment is lost and it was the **last** segment sent — there's no later data to trigger duplicate ACKs, because nothing came after it? The sender simply... hears nothing.

This is what the **Retransmission Timeout (RTO)** exists for: every time the sender transmits a segment, it starts a timer. If no acknowledgment covering that segment arrives before the timer expires, the sender assumes the segment (or its ACK) was lost, and retransmits — without needing any duplicate-ACK evidence at all.

```
Sender:  seq=2001, 500 bytes ──────────▶   sent at T=0.000s
                                             (timer started, RTO = 300ms)

         ... nothing arrives before T=0.300s ...

Timer expires at T=0.300s → sender assumes loss → retransmits seq=2001
```

The obvious hard question is: **how long should the timer be?** Too short, and the sender retransmits data that was simply still in flight, wasting bandwidth and potentially making congestion worse (Chapter 62). Too long, and a genuinely lost segment sits unretransmitted for an unnecessarily long time, stalling the whole connection (since, per Section 5, no data after a gap can be delivered to the application until the gap closes).

The right answer clearly depends on how long round trips on this specific connection actually take — a connection to a server on the same continent might have a 20ms RTT, while a connection to a server on the other side of the planet might have a 250ms RTT, and a fixed timeout can't be right for both. Section 8 covers exactly how TCP estimates this.

---

## 8. Computing RTO: The Jacobson/Karn Algorithm

TCP doesn't guess at a fixed timeout — it continuously *measures* round-trip time (RTT) on the live connection and adapts. The algorithm in wide use today descends directly from Van Jacobson's 1988 paper (with Phil Karn's earlier contribution about which samples are even valid to use), and it works in two parts.

**Part 1 — smoothing the RTT estimate.** Every time a segment's ACK arrives, the time between sending it and receiving that ACK is one RTT sample. Rather than reacting to each sample individually (which would make the timeout jump around wildly with normal network jitter), TCP maintains two running values that get updated with every new sample:

```
SRTT  = Smoothed Round-Trip Time     (a weighted moving average of RTT)
RTTVAR = Round-Trip Time Variation   (a weighted moving average of how
                                       much RTT is actually varying)

On each new RTT sample R:
  RTTVAR = (1 - β) × RTTVAR + β × |SRTT - R|      (β is typically 0.25)
  SRTT   = (1 - α) × SRTT   + α × R                (α is typically 0.125)

Then:
  RTO = SRTT + 4 × RTTVAR
```

Worked through with real numbers: suppose a connection has settled into `SRTT = 50ms`, `RTTVAR = 10ms` (giving `RTO = 50 + 40 = 90ms`), and a new RTT sample of `70ms` arrives:

```
RTTVAR = 0.75 × 10  + 0.25 × |50 - 70| = 7.5 + 5.0   = 12.5ms
SRTT   = 0.875 × 50 + 0.125 × 70       = 43.75 + 8.75 = 52.5ms

New RTO = 52.5 + 4 × 12.5 = 52.5 + 50 = 102.5ms
```

The `4 × RTTVAR` term is deliberately generous: it means the timeout sits well above the *typical* RTT, scaled by how *unstable* that RTT has recently been — a connection with wildly varying latency gets a proportionally larger safety margin, while a very stable connection can afford a tighter timeout.

**Part 2 — Karn's rule: which samples count?** There's a subtlety that makes this measurement genuinely tricky: if a segment is retransmitted, and an ACK then arrives, **which transmission does that ACK actually correspond to** — the original, or the retransmission? There's no way to tell for certain from the ACK alone (this is called retransmission ambiguity). Karn's algorithm's fix is simple and effective: **never use an RTT sample from a segment that was retransmitted** — only measure RTT from segments that were acknowledged on their very first attempt, sidestepping the ambiguity entirely rather than trying to resolve it.

**Exponential backoff.** If a retransmitted segment *also* times out (loss is persisting, or the network is badly congested), TCP doesn't just retry at the same RTO again — it doubles it on each successive failure (`RTO`, then `2×RTO`, then `4×RTO`, and so on, typically up to some maximum), specifically to avoid making an already-struggling network worse by hammering it with retransmissions at a fixed, possibly too-aggressive rate. This exponential backoff is one of the direct behavioral links between this chapter's retransmission logic and Chapter 62's congestion control — both exist because a network under stress needs senders to back off, not pile on.

---

## 9. Selective Acknowledgment (SACK)

Section 5's worked example showed a real inefficiency in plain cumulative ACKs: after segment 2 was lost, the receiver was silently holding segment 3 (bytes 2001–2500) in a buffer the whole time, but had no way to *tell the sender that* — the sender only knew "you're missing something starting at 1501," not "everything after that is already fine, just resend the one gap."

**SACK** (Selective Acknowledgment, negotiated as an option during the handshake — Chapter 59, Section 7 — via the SACK-permitted flag) fixes exactly this. A SACK-enabled receiver can report, in addition to the ordinary cumulative ACK, one or more specific *ranges* of out-of-order data it already has:

```
Without SACK:
  ack=1501   (sender only knows: "resend starting from 1501,
               I have no idea how much you need to resend")

With SACK:
  ack=1501, SACK block: [2001–2500]
  (sender now knows: "resend exactly 1501–2000 — I already have
    2001–2500, don't bother resending that part too")
```

This matters most when *multiple* segments are lost from the same window: without SACK, a sender might have to guess how much to retransmit (or wait for a full timeout and resend everything the receiver hasn't cumulatively confirmed, even parts it already has); with SACK, the sender gets a precise map of exactly which byte ranges are actually missing. SACK is now essentially universal in real-world TCP stacks and is the mechanism Chapter 63's fast recovery leans on for efficient loss recovery across multiple lost segments in one window.

---

## 10. Seeing It Live

Capturing a real connection under artificial packet loss makes Sections 5–7 directly observable. On Linux, `tc` (traffic control) can inject loss for an experiment:

```
$ sudo tc qdisc add dev eth0 root netem loss 20%   # drop ~20% of packets
$ curl -s http://example.com/largefile > /dev/null &
$ sudo tcpdump -i eth0 -n tcp -S
```

In the resulting capture, three patterns from this chapter become visible directly in the flag/seq/ack fields:

```
... seq 3001:3501, length 500                 # original segment
... ack 3001, ack 3001, ack 3001              # three duplicate ACKs —
                                                 fast-retransmit trigger
... seq 3001:3501, length 500  [retransmission]  # sender resends before
                                                    any timeout fired
```

or, for a timeout-driven case (no later data to trigger duplicate ACKs):

```
... seq 4501:5001, length 500                  # last segment sent
                                                  (long gap — no ACK, no
                                                   later segments to
                                                   trigger dup ACKs)
... seq 4501:5001, length 500  [retransmission]  # RTO expired, resent
```

`ss -ti` on a live connection also directly exposes the RTT and RTO values Section 8 computed by hand:

```
$ ss -ti
ESTAB  0  0  192.168.1.42:51342  142.250.80.46:443
    cubic rtt:52.5/12.5 rto:102 ...
```

That `rtt:52.5/12.5` is exactly `SRTT/RTTVAR` from the worked calculation in Section 8, and `rto:102` (rounded) matches the computed `102.5ms` — real, observable confirmation of the exact arithmetic this chapter walked through.

---

## 11. Code: A Minimal Retransmission Timer

A simplified but faithful sketch of Section 8's algorithm:

```go
type RTOEstimator struct {
    srtt, rttvar float64 // in milliseconds
    initialized  bool
}

const alpha, beta = 0.125, 0.25

func (e *RTOEstimator) Sample(rttMs float64) {
    if !e.initialized {
        e.srtt = rttMs
        e.rttvar = rttMs / 2
        e.initialized = true
        return
    }
    e.rttvar = (1-beta)*e.rttvar + beta*abs(e.srtt-rttMs)
    e.srtt = (1-alpha)*e.srtt + alpha*rttMs
}

func (e *RTOEstimator) RTO() float64 {
    rto := e.srtt + 4*e.rttvar
    if rto < 200 {
        rto = 200 // typical implementations enforce a minimum floor
    }
    return rto
}

func abs(x float64) float64 {
    if x < 0 {
        return -x
    }
    return x
}

// Karn's rule in practice: only feed Sample() with RTTs measured from
// segments that were NOT retransmitted, e.g.:
//
//   if !segment.WasRetransmitted {
//       estimator.Sample(timeAckReceived.Sub(segment.SentAt).Seconds() * 1000)
//   }
//
// A segment whose timer expires without an ACK doubles the next RTO
// (exponential backoff) rather than calling Sample() at all:
//
//   nextRTO := currentRTO * 2
```

---

## 12. Common Misconceptions

- **"TCP acknowledges every packet with its own ACK."** Cumulative ACKs mean one ACK can (and routinely does) confirm several segments at once, as Section 4 showed directly. Real stacks also frequently delay ACKs slightly (see "delayed ACK" in production TCP behavior) to piggyback them on outgoing data or batch several segments' worth of confirmation into fewer ACKs.
- **"A duplicate ACK means the ACK itself was duplicated on the wire."** It means the *same* cumulative position was acknowledged more than once, almost always because later, out-of-order data arrived while a gap remained open — not because a network fault literally copied one ACK segment.
- **"Sequence numbers count packets."** Section 2 covered this directly: they count bytes. Two connections sending the same *number of segments* can have wildly different sequence number progressions if their segments carry different amounts of data.
- **"RTO is a fixed timeout, like 3 seconds, that all TCP connections use."** It's computed per-connection, continuously, from Section 8's algorithm, and differs enormously between, say, a connection on the same LAN (RTO potentially near the algorithm's minimum floor) and a connection to another continent.
- **"Retransmission always means the original packet was lost."** It can also mean the *ACK* for that packet was lost — the sender has no way to distinguish "my data never arrived" from "my data arrived but the ACK confirming it didn't," and in both cases the correct, safe action is the same: retransmit (the receiver's cumulative ACK logic naturally handles a redundant retransmission of already-received data without any harm — it just re-acknowledges the same position).

---

## 13. Production Notes

- **Spurious retransmissions** happen when the RTO estimate is too aggressive for a connection whose latency just spiked temporarily (not due to loss) — the retransmission wasn't actually needed, and the extra traffic can make transient congestion worse. This is one motivation behind more sophisticated modern refinements layered on top of the classic algorithm (e.g., the TCP timestamp option feeding more precise, per-segment RTT samples rather than relying purely on ACK arrival timing).
- **Bufferbloat interacts badly with RTO.** If intermediate routers buffer excessively before dropping packets (rather than dropping promptly under congestion), measured RTT samples can balloon, inflating SRTT and RTO and slowing down loss detection — a real, documented issue that later congestion-control work (Chapter 62's BBR, in particular) was partly designed to address by not relying purely on loss as a congestion signal.
- **SACK is effectively mandatory in production today** — nearly every modern OS and server enables it by default, because the efficiency gain during multi-segment loss recovery (Section 9) is substantial, especially on high-bandwidth, high-latency connections where retransmitting more than strictly necessary is expensive.
- **Monitoring retransmission rate is a standard health signal** for production TCP-heavy services: a rising retransmission percentage on a server's connections is one of the most direct available signals of either network path degradation or the server itself becoming a bottleneck (unable to process and ACK data promptly).

---

## 14. What's Simplified Here

The RTO formula and constants (`α = 0.125`, `β = 0.25`, the `4×RTTVAR` term) follow RFC 6298's specification, which real operating systems implement closely but not always identically — Linux, for instance, layers additional refinements (like using timestamp-based RTT measurement per-segment rather than purely relying on cumulative ACK timing) on top of this baseline. Fast retransmit's exact interaction with the congestion window, and fast recovery's behavior after a fast retransmit fires, are deliberately deferred to Chapter 63, since they depend on congestion-control concepts (Chapter 62) not yet introduced. SACK's exact wire format (a TCP option carrying up to four discontiguous byte ranges) is likewise not detailed here — Section 9 covers the concept and its purpose, not its byte layout.

---

## Interview Questions & Model Answers

**Beginner: What does a TCP acknowledgment number actually mean?**

It means "I have correctly received every byte up through this position minus one, with no gaps — send me this byte next." It's cumulative: one ACK can confirm many previously-sent segments at once, as long as they form an unbroken run starting from the beginning of the connection.

**Intermediate: A receiver gets segments covering bytes 1–1000 and 2001–3000, but the segment covering 1001–2000 is lost. What ACK value does the receiver send, and why can't it acknowledge the data it already has past the gap?**

It sends `ack=1001` repeatedly — even after 2001–3000 arrives — because a cumulative ACK can only advance past a byte position once every byte up to that point has arrived with no gaps. The receiver is free to buffer 2001–3000 internally, but it cannot report it via the ACK number until 1001–2000 fills the hole; once it does, the ACK jumps straight to 3001, crediting both the retransmitted segment and the already-buffered one in a single step.

**Advanced: Why does TCP need both a timeout-based (RTO) and a duplicate-ACK-based mechanism for detecting loss, instead of relying on just one?**

They cover complementary failure shapes. A timeout is the only option when a lost segment is the last one sent on the connection (or followed by nothing for a while) — there's no later data to generate the duplicate ACKs a faster mechanism would need, so the sender has to wait out a timer calibrated to the connection's measured RTT and its variability. Duplicate ACKs handle the far more common case where loss happens mid-stream and later segments keep arriving and getting acknowledged with the same stuck cumulative position — that repetition is available immediately, often well before the RTO would have expired, so treating three duplicate ACKs as a loss signal (fast retransmit, Chapter 63) recovers much faster than waiting for the timer. Relying on RTO alone would leave a lot of unnecessary latency on the table whenever the faster signal was actually available; relying on duplicate ACKs alone would leave connections stalled indefinitely whenever no later data existed to generate them.

---

## Exercises

### Easy
1. A sender transmits three 400-byte segments starting at sequence number 5001. What are the sequence numbers of the first byte of each segment?
2. If a receiver has correctly received every byte through 8000 with no gaps, what ACK value does it send?
3. Explain, in one sentence, why TCP retransmits based on byte ranges rather than by resending "packet number 4."

### Medium
4. A connection's SRTT is 40ms and RTTVAR is 15ms. Compute the current RTO using the Jacobson/Karn formula.
5. Following on from question 4, a new RTT sample of 90ms arrives. Compute the updated RTTVAR, SRTT, and RTO.
6. Explain Karn's rule and construct a concrete scenario (with made-up sequence numbers and timestamps) showing why an RTT sample taken from a retransmitted segment's ACK would be ambiguous.

### Hard
7. A sender transmits five 200-byte segments covering bytes 10001–11000. Segments 2 and 4 (in transmission order) are both lost; segments 1, 3, and 5 arrive successfully. Write out every ACK the receiver sends as each segment arrives (in order 1, 3, 5), and explain precisely what the sender can and cannot conclude about which bytes are missing from cumulative ACKs alone — then explain what a SACK-enabled receiver would additionally report that a non-SACK receiver couldn't.
8. Design (in pseudocode or precise prose) a rule for when a sender applying exponential backoff after repeated RTO expirations should give up and declare the connection dead, rather than continue doubling the timeout forever. Justify your chosen stopping condition.

---

## Summary

| Term | Meaning |
|---|---|
| Sequence number | The position of the first byte of a segment within the connection's overall byte stream |
| Cumulative ACK | Acknowledgment meaning "everything up to this byte, with no gaps, has arrived" |
| Duplicate ACK | The same cumulative ACK value repeated, signaling later data arrived while a gap remains open |
| RTO (Retransmission Timeout) | How long a sender waits for an ACK before assuming loss and retransmitting |
| SRTT / RTTVAR | Smoothed running estimates of round-trip time and its variability, used to compute RTO |
| Karn's rule | Never use an RTT sample from a retransmitted segment, to avoid ambiguity about which transmission an ACK confirms |
| Exponential backoff | Doubling the RTO after each successive retransmission timeout, to avoid worsening congestion |
| SACK | Selective Acknowledgment — reports specific out-of-order byte ranges already received, beyond the cumulative ACK |

Sequence numbers and ACKs guarantee that data which *is* sent eventually arrives correctly and in order. But nothing so far has stopped a fast sender from simply sending more data than a slow receiver can hold in its buffer at once. Chapter 61 covers the fix: the sliding window.
