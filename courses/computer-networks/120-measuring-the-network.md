# Chapter 120: Measuring the Network — Latency, Throughput, Jitter, and Packet Loss

> **"'The network is slow' is not a measurement. It's a symptom looking for a number. This chapter supplies the numbers — and the tools that produce them honestly."**

---

## Table of Contents

1. [Why This Chapter Exists](#1-why-this-chapter-exists)
2. [Latency: The Time for One Bit to Get There](#2-latency-the-time-for-one-bit-to-get-there)
3. [Measuring Latency with ping](#3-measuring-latency-with-ping)
4. [Bandwidth vs. Throughput vs. Goodput](#4-bandwidth-vs-throughput-vs-goodput)
5. [Measuring Throughput with iperf3](#5-measuring-throughput-with-iperf3)
6. [Jitter: Why Consistency Matters More Than Averages](#6-jitter-why-consistency-matters-more-than-averages)
7. [Why Jitter Specifically Breaks VoIP and Video](#7-why-jitter-specifically-breaks-voip-and-video)
8. [Packet Loss: Definition and Causes](#8-packet-loss-definition-and-causes)
9. [Measuring Loss Per Hop with mtr](#9-measuring-loss-per-hop-with-mtr)
10. [Bandwidth-Delay Product: Where Latency and Throughput Meet](#10-bandwidth-delay-product-where-latency-and-throughput-meet)
11. [A Real Combined Measurement Session](#11-a-real-combined-measurement-session)
12. [Hands-On Experiment](#12-hands-on-experiment)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes](#14-production-notes)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary and the Bridge to Chapter 121](#18-summary-and-the-bridge-to-chapter-121)

---

## 1. Why This Chapter Exists

Chapter 119 gave you the tools to see individual packets. But "the site is slow" and "the video call keeps breaking up" are rarely diagnosed by staring at one packet — they're diagnosed by *measuring*, over time, four specific, precisely defined quantities: how long a round trip takes (latency), how much data actually moves per second (throughput, as distinct from the theoretical ceiling, bandwidth), how much that latency varies from packet to packet (jitter), and what fraction of packets never arrive at all (packet loss).

The problem this chapter solves is that these four words get used loosely and interchangeably in ordinary conversation — "slow internet" gets blamed on "bandwidth" almost reflexively — when in reality each one has a distinct root cause, a distinct fix, and a distinct tool built specifically to measure it. Confusing them wastes time chasing the wrong fix: buying a faster internet plan (more bandwidth) does nothing for a video call that stutters because of jitter, and a perfectly fast, low-latency connection can still make a large file transfer crawl if packet loss is silently forcing constant retransmission (Chapter 60).

## 2. Latency: The Time for One Bit to Get There

**Intuitive definition:** latency is the delay between sending something and it arriving — the time a letter spends in transit, not how many letters the postal service can carry per day.

**Precise definition:** **latency** is the time for a single unit of data (conventionally, one packet) to travel from sender to receiver. **Round-trip time (RTT)** — what `ping` actually measures — is latency there *and back*: send, arrive, receiver responds, response arrives. One-way latency is genuinely harder to measure directly (it requires tightly synchronized clocks on both ends, since you're comparing a timestamp from one machine against a timestamp from another), which is exactly why RTT, measured entirely by one clock on one machine, is the practical, ubiquitous metric.

**What latency is made of**, in order of typical contribution over a real long-distance path:

| Component | Cause | Typical order of magnitude |
|---|---|---|
| Propagation delay | speed of light in the medium (~200,000 km/s in fiber, Chapter 22) × distance | ~4-5ms per 1,000km one-way in fiber |
| Transmission delay | time to push all the bits of a frame onto the wire = frame size / link bandwidth | microseconds on modern links, more on slow ones |
| Queuing delay | time spent waiting in a router's buffer behind other traffic | 0 to tens/hundreds of ms under congestion |
| Processing delay | time a router or switch spends deciding where to forward a packet (Chapter 45) | microseconds on modern hardware |

Propagation delay is a hard physical floor — Chapter 23 already established that a round trip to a geostationary satellite (~36,000km up) costs roughly 240ms one-way from physics alone, no amount of engineering budget removes that. Queuing delay, by contrast, is the *variable* component that engineering (buffer sizing, congestion control from Chapter 62, quality-of-service prioritization) can and does influence — and it's the component most responsible for latency that changes minute to minute on an otherwise-fixed physical path.

## 3. Measuring Latency with ping

Chapter 54 built `ping`'s mechanism (ICMP Echo Request/Reply); Chapter 56, Section 5 used it for a first diagnostic read. Here, the same tool is used specifically as a *measurement instrument*, with attention to the statistics it reports:

```
$ ping -c 20 93.184.216.34
PING 93.184.216.34 (93.184.216.34) 56(84) bytes of data.
64 bytes from 93.184.216.34: icmp_seq=1 ttl=56 time=14.021 ms
64 bytes from 93.184.216.34: icmp_seq=2 ttl=56 time=13.887 ms
64 bytes from 93.184.216.34: icmp_seq=3 ttl=56 time=14.502 ms
...
64 bytes from 93.184.216.34: icmp_seq=20 ttl=56 time=13.951 ms

--- 93.184.216.34 ping statistics ---
20 packets transmitted, 20 received, 0% packet loss, time 19023ms
rtt min/avg/max/mdev = 13.887/14.203/15.877/0.482 ms
```

Reading the final line precisely — this is the single most information-dense line `ping` produces, and it's worth being able to decode all four numbers on sight:

- **min (13.887ms)** — the best-case RTT observed; the closest this path gets to its physical floor (propagation + transmission + processing, with essentially zero queuing delay at that instant).
- **avg (14.203ms)** — the mean RTT across all replies; the number most people mean by "the latency," though by itself it hides everything about consistency.
- **max (15.877ms)** — the worst-case RTT observed; a max much larger than the avg is an early, cheap signal of jitter (Section 6) even before you've computed anything formally.
- **mdev (0.482ms)** — mean deviation, a rough (not the textbook-strict standard deviation, though used the same way in practice) measure of how much RTT varies sample to sample — this number *is*, informally, `ping`'s built-in jitter estimate, revisited precisely in Section 6.

**A single ping tells you almost nothing reliable; a run of pings tells you the actual distribution.** One measurement can catch a momentary fluke in either direction; `-c 20` (or more, for a noisy path) is the minimum for `avg`/`mdev` to mean anything.

## 4. Bandwidth vs. Throughput vs. Goodput

Chapter 16 defined bandwidth in the physical-layer sense: a *range of frequencies* (measured in Hz), and drew the crucial distinction that "bandwidth" in everyday tech conversation almost always actually means **throughput** — a data rate in bits per second. This chapter revisits that distinction specifically as a measurement problem, and adds a third term that matters just as much in practice.

| Term | Precise meaning | Analogy |
|---|---|---|
| **Bandwidth** | the theoretical maximum data rate a link *could* carry, given its physical properties (Chapter 16, 18) | the number of lanes on a highway |
| **Throughput** | the data rate *actually achieved* over a real link, at a given moment, including all protocol overhead | how many cars per minute actually cross a specific point |
| **Goodput** | the rate of *useful application data* actually delivered, after subtracting protocol headers, retransmissions, and encryption overhead | how many cars actually reach their destination, excluding ones that got lost and had to loop back |

Throughput is always less than or equal to bandwidth (you cannot exceed the physical ceiling), and goodput is always less than or equal to throughput (headers and retransmitted bytes consume capacity without delivering new application data). A concrete illustration: a "1 Gbps" advertised link (bandwidth) might, under real conditions with TCP/IP/Ethernet overhead (roughly 5-8% typical), sustain closer to 940 Mbps of raw throughput — and if packet loss (Section 8) is forcing 10% of segments to be retransmitted (Chapter 60), the actual goodput delivered to the application could be meaningfully lower still, even though the "line rate" never changed.

This is precisely why "I pay for 500 Mbps internet but my download is only showing 420 Mbps" is not automatically a fault — some of that gap is expected protocol overhead, and only a much larger, more persistent gap (or one that correlates with measurable packet loss) is a real problem worth chasing.

## 5. Measuring Throughput with iperf3

`ping`'s ICMP packets are tiny and infrequent — useless for measuring how much *sustained* data a path can carry. `iperf3` exists specifically to answer that question: it runs a server and client that deliberately try to push as much data as possible (or a controlled, specified rate) across a connection and report the achieved throughput.

```
# On the receiving/server machine
$ iperf3 -s
-----------------------------------------------------------
Server listening on 5201
-----------------------------------------------------------

# On the sending/client machine
$ iperf3 -c 10.0.0.5
Connecting to host 10.0.0.5, port 5201
[  5] local 10.0.0.20 port 51022 connected to 10.0.0.5 port 5201
[ ID] Interval           Transfer     Bitrate         Retr  Cwnd
[  5]   0.00-1.00   sec  112 MBytes   940 Mbits/sec    0   1.41 MBytes
[  5]   1.00-2.00   sec  111 MBytes   933 Mbits/sec    0   1.41 MBytes
[  5]   2.00-3.00   sec  110 MBytes   925 Mbits/sec    2   1.20 MBytes
[  5]   3.00-4.00   sec  111 MBytes   931 Mbits/sec    0   1.35 MBytes
...
[  5]   0.00-10.00  sec  1.09 GBytes   937 Mbits/sec  2             sender
[  5]   0.00-10.00  sec  1.09 GBytes   935 Mbits/sec               receiver
```

Reading this against earlier chapters: the **Bitrate** column is throughput measured live, second by second — this is exactly the number Chapter 62's congestion control algorithm (CUBIC or BBR, depending on OS defaults) is actively producing in real time by adjusting the sender's window. The **Retr** column counts TCP segment retransmissions during that second (Chapter 60) — the `2` retransmissions in the third interval, correlating with a slightly lower bitrate that interval, is a live, visible instance of packet loss forcing exactly the throughput cost Section 4 described. **Cwnd** is TCP's congestion window (Chapter 62) at that instant — watching it shrink after a loss event and grow again afterward is literally AIMD, happening on your screen instead of in a diagram.

Useful flags worth knowing:

```
$ iperf3 -c 10.0.0.5 -u -b 10M     # test UDP instead of TCP, at a fixed 10 Mbps target rate
$ iperf3 -c 10.0.0.5 -R            # reverse direction: measure server-to-client instead of client-to-server
$ iperf3 -c 10.0.0.5 -P 4          # 4 parallel streams — useful for saturating a link one TCP stream alone can't fill
```

`-u` matters specifically because UDP mode reports **packet loss directly and explicitly** (since UDP, per Chapter 58, has no retransmission to mask it) — TCP mode instead shows loss indirectly, as retransmissions and reduced throughput, since TCP is actively working to hide loss from the application:

```
$ iperf3 -c 10.0.0.5 -u -b 50M
[ ID] Interval           Transfer     Bitrate         Jitter    Lost/Total Datagrams
[  5]   0.00-10.00  sec  59.6 MBytes  50.0 Mbits/sec  0.412 ms  3/42800 (0.0070%)
```

This single line packs three of this chapter's four core quantities into one output: **Bitrate** (throughput, Section 4), **Jitter** (Section 6, computed by `iperf3` itself using the RFC 3550 algorithm originally defined for RTP/VoIP traffic — no coincidence, given Section 7), and **Lost/Total Datagrams** (packet loss, Section 8) — measured directly, in real conditions, rather than inferred.

## 6. Jitter: Why Consistency Matters More Than Averages

**Naive assumption:** if average latency is low, the connection is "fast enough" for anything.

**Why that's wrong:** average latency says nothing about *consistency*. A connection where every packet takes exactly 100ms is, for many applications, in a completely different category of "good" than a connection averaging 60ms but swinging wildly between 20ms and 150ms packet to packet — even though the second connection has the *lower* average.

**Precise definition: jitter is the variation in latency between consecutive packets** — formally, in the RFC 3550 (RTP) definition `iperf3` implements, a smoothed average of the absolute difference in one-way transit time between successive packets. Intuitively: not "how long does a packet take," but "how much does that time change from one packet to the next."

```
Low jitter (consistent, even if not blazing fast):
  Packet 1: 45ms   Packet 2: 47ms   Packet 3: 44ms   Packet 4: 46ms

High jitter (same rough average, wildly inconsistent):
  Packet 1: 20ms   Packet 2: 90ms   Packet 3: 15ms   Packet 4: 105ms
```

Both examples above average to roughly the same latency — but only the first is usable for real-time interactive traffic, for the specific reason Section 7 explains.

## 7. Why Jitter Specifically Breaks VoIP and Video

Here is the mechanism, not just the claim. A VoIP call or video stream is fundamentally a sequence of small chunks of audio/video generated at a fixed rate (say, one 20-millisecond audio frame every 20 milliseconds, continuously) that must be *played back* at that same fixed rate for the result to sound and look natural — a human voice sped up or stretched even slightly is immediately, unpleasantly noticeable.

The receiving application handles this by buffering incoming packets briefly in a **jitter buffer**: instead of playing each packet the instant it arrives, it holds a small cushion of packets and plays them out at the steady, correct rate, smoothing over small arrival-time variations. This is the direct fix for jitter — but it has a hard limit, and that limit is exactly where jitter becomes a real problem:

- **If jitter is low and predictable**, a small jitter buffer (say, 40-60ms) comfortably absorbs the variation, and playback stays smooth even though individual packets arrive at slightly different times.
- **If jitter is high**, the buffer must either grow large enough to absorb the worst-case variation — directly adding to perceived latency, which for a live conversation becomes its own separate problem past roughly 150ms one-way (Chapter 23's satellite-latency discussion already established when a round trip starts feeling laggy in conversation) — or, if the buffer is kept small to preserve low latency, packets that arrive *later* than the buffer can wait for are simply **dropped and treated as lost**, even though they weren't lost on the network at all — they just arrived too late to be useful.

This is the precise, mechanical answer to "why does jitter specifically break VoIP/video even when average latency looks fine": **a real-time stream doesn't get to ask the network to wait for a slow packet — it has a hard deadline to play something, every 20 milliseconds, forever**, and average latency describes nothing about whether any individual packet meets that deadline. A bulk file transfer (Chapter 60's TCP model) has no such deadline — it happily waits arbitrarily long for a retransmission and reassembles everything perfectly in order, which is exactly why the same underlying network jitter that ruins a video call is often completely invisible in a file download's total transfer time.

## 8. Packet Loss: Definition and Causes

**Precise definition:** packet loss is the fraction of packets sent that never arrive at their destination at all, expressed as a percentage (`3/42800 (0.0070%)` in Section 5's `iperf3` UDP example).

**Real causes, roughly in order of how often each shows up in practice:**

| Cause | Mechanism |
|---|---|
| Buffer overflow / congestion | a router's queue fills faster than it can drain, and it drops packets rather than delay them indefinitely (Chapter 62's congestion signal) |
| Physical-layer errors | noise/attenuation (Chapter 17) corrupts a frame badly enough that a CRC check (Chapter 19) fails and it's discarded |
| Wireless interference/range | a weak or contested Wi-Fi/cellular signal drops frames the physical layer never delivers cleanly at all |
| Misconfiguration | a firewall rule, ACL, or routing loop silently discards traffic that was never meant to be dropped |
| Hardware failure | a failing NIC, cable, or switch port produces above-baseline loss on one specific path or interface |

TCP (Chapter 60) hides loss from the application by silently retransmitting — which is why loss shows up *indirectly*, as reduced throughput and increased latency (Section 4-5), rather than as an outright application error, unless loss becomes severe enough that retransmissions themselves start timing out. UDP (Chapter 58) has no such safety net at all — an application using UDP (VoIP, live video, gaming, DNS) either handles loss itself at the application layer or simply experiences it directly as a glitch, a dropped frame, or a missing DNS reply that has to be retried from scratch.

## 9. Measuring Loss Per Hop with mtr

`ping` measures loss to one specific final destination. `traceroute` (Chapter 54, Chapter 56 Section 6) shows the hops along a path, but only as a handful of individual probes — it doesn't sustain measurement long enough to reveal *which specific hop* is dropping packets over time. **`mtr` ("My TraceRoute") combines both ideas**: it continuously sends probes like `traceroute`, but keeps a running, live statistical summary — including loss percentage — for *every hop along the path*, not just the final destination.

```
$ mtr -n -c 100 --report 93.184.216.34
Start: 2026-08-09T10:55:02+0000
HOST: client                       Loss%   Snt   Last   Avg  Best  Wrst StDev
  1.|-- 192.168.1.1                 0.0%   100    0.5   0.5   0.4   1.1   0.1
  2.|-- 10.20.0.1                   0.0%   100    5.1   5.3   4.9   9.8   0.6
  3.|-- 172.16.4.9                 12.0%   100    5.4  15.2   5.1  84.3  18.4
  4.|-- 152.195.64.1                0.0%   100   13.4  13.6  13.2  16.1   0.4
  5.|-- 93.184.216.34               0.0%   100   14.0  14.2  13.7  17.9   0.5
```

This single report is one of the most decisive diagnostic outputs in networking, precisely because of what it isolates: **hop 3 shows 12% loss and a high, unstable Avg/Wrst/StDev, but hops 4 and 5 — further along the *same path*, necessarily passing through hop 3 to get there — show 0% loss and low, stable latency.** The correct reading, and the one that trips up almost everyone the first time, is: **the loss is a hop-3 *reporting* artifact, not necessarily a real hop-3 problem**, because many routers deprioritize or rate-limit the ICMP TTL-exceeded replies `mtr`/`traceroute` rely on (Chapter 54, Section 12's "*" discussion) far more aggressively than they prioritize actual forwarded traffic. **The trustworthy signal is loss that persists at a hop and at every hop after it** — since a packet genuinely dropped at hop 3 could never be measured arriving at hop 4 or 5 at all. Loss isolated to one hop, with clean hops beyond it, usually means "that router is slow to reply to probes," not "that router is dropping your actual traffic."

A pattern that *does* indicate a real problem: loss appearing at hop N and persisting, at a similar or growing percentage, at every subsequent hop — that's consistent with real packet loss actually occurring at or before that hop, since every later hop can only report packets that actually survived to reach it.

## 10. Bandwidth-Delay Product: Where Latency and Throughput Meet

One more measurement concept worth naming precisely, because it explains a real, common confusion ("why is my high-bandwidth, high-latency satellite connection still slow for a single download?"): the **bandwidth-delay product (BDP)** is bandwidth × RTT, and it represents *how much data can be "in flight" on the path at any given instant* — the amount of data that has left the sender but hasn't yet been acknowledged.

```
BDP = bandwidth × RTT

Example: a 100 Mbps link with 200ms RTT (a realistic satellite or
long-haul intercontinental path):

  BDP = 100,000,000 bits/sec × 0.2 sec = 20,000,000 bits ≈ 2.5 MB
```

This directly connects to Chapter 61's sliding window: **if a single TCP connection's window (the maximum unacknowledged data it's allowed in flight) is smaller than the path's BDP, that one connection can never use the full available bandwidth**, no matter how fast the link physically is — it will always be sitting idle, waiting for acknowledgments, for part of every round trip. This is precisely why window scaling (Chapter 61) matters enormously on long, fast paths, and why a single-threaded download over a high-BDP path often underperforms compared to several parallel connections (each with its own window) even on the exact same physical link — the same reason `iperf3 -P 4` in Section 5 exists as an option at all.

## 11. A Real Combined Measurement Session

Putting all four quantities together for the complaint: **"our video calls to the branch office are choppy, but file transfers to the same office are fine."**

```
Step 1 — Baseline latency and its own consistency
  $ ping -c 50 branch-office.internal
  rtt min/avg/max/mdev = 38.112/41.847/198.223/22.940 ms
  -> avg looks fine (42ms), but max (198ms) and mdev (23ms) are both
     enormous relative to avg — a strong early jitter signal (Section 6).

Step 2 — Confirm and quantify jitter and loss directly with a UDP test,
  which is what the video call itself actually uses on the wire
  $ iperf3 -c branch-office.internal -u -b 2M -t 30
  Jitter: 24.318 ms    Lost/Total Datagrams: 41/15000 (0.27%)
  -> Jitter and loss both measured directly, not inferred — confirms
     Step 1's suspicion with the actual RFC 3550 jitter calculation.

Step 3 — Confirm raw TCP throughput is fine, explaining why file transfers
  aren't affected
  $ iperf3 -c branch-office.internal -t 10
  937 Mbits/sec average, 0 retransmissions
  -> Plenty of bandwidth, essentially no loss over TCP's timescale —
     confirms file transfers have no reason to be affected (Section 4/8).

Step 4 — Localize the jitter/loss to a specific hop
  $ mtr -n -c 200 --report branch-office.internal
  hop 4: Loss 0.0%, StDev 1.1     (stable)
  hop 5: Loss 0.3%, StDev 34.7    (unstable, and loss persists to the end)
  hop 6: Loss 0.3%, StDev 33.9    (unstable, loss persists here too)
  -> Loss persisting at and beyond hop 5, with a large StDev at exactly
     the point it appears, points at hop 5 as the real source — likely
     a congested or oversubscribed link between the two sites.

Conclusion: the path has plenty of raw bandwidth and low average latency
(both fine for file transfers, which don't care about per-packet timing
consistency), but real, hop-5-localized jitter and loss severe enough to
disrupt a real-time UDP stream's fixed playback schedule, even though
that same loss and jitter is nearly invisible in a TCP transfer's total
completion time.
```

## 12. Hands-On Experiment

If you have access to two machines (a home server and a cloud VM, or two devices on your LAN), run:

```bash
# On machine B (the "server")
iperf3 -s

# On machine A
ping -c 30 <machine B's IP>                    # latency + informal jitter (mdev)
iperf3 -c <machine B's IP>                     # TCP throughput
iperf3 -c <machine B's IP> -u -b 20M -t 20     # UDP throughput, real jitter, real loss
mtr -n -c 100 --report <machine B's IP>        # per-hop loss (mostly interesting over a WAN path)
```

Compare the `mdev` from `ping` against the `Jitter` figure `iperf3 -u` reports for the same path — they measure conceptually the same thing (variation in delay) through two different algorithms, and seeing both agree in magnitude on a real link is a strong way to internalize what "jitter" actually is, beyond the definition in Section 6.

## 13. Common Misconceptions

- **"Higher bandwidth always means a faster-feeling connection."** As Section 7 and the BDP discussion in Section 10 both show, a fast, high-bandwidth connection with bad jitter or a large BDP mismatch can feel worse for real-time or single-stream use than a slower but more consistent one — bandwidth is one of four independent quantities, not a stand-in for all of them.
- **"0% packet loss to the final destination means the whole path is loss-free."** `mtr`'s own Section 9 example directly contradicts this — intermediate-hop loss readings can be entirely an artifact of ICMP deprioritization, and the only fully trustworthy full-path loss number is the one measured end to end (as `ping`'s or `iperf3`'s final-destination statistics do), which is exactly why `mtr`'s value is in showing loss *patterns* across hops, not treating any single hop's number as gospel in isolation.
- **"Jitter is just another word for high latency."** They are independent: Section 6 showed two examples with nearly identical average latency and wildly different jitter — a connection can have very low latency and very high jitter (rare, but possible with certain shared/contended wireless mediums) or high latency and near-zero jitter (a stable, if slow, long-haul satellite link).
- **"`iperf3` measures what a user actually experiences on a normal web browsing session."** `iperf3` measures the maximum sustained rate a path can carry under a deliberately generated, saturating load — a genuinely different measurement from typical real-world usage, which is bursty and rarely saturates a link the way `iperf3` intentionally does; treat `iperf3`'s number as "the ceiling," not "the typical experience."

## 14. Production Notes

- **Synthetic monitoring runs these exact tools continuously, not just reactively.** Production observability systems commonly schedule `ping`/`mtr`-style probes between key points in the infrastructure on a fixed interval, feeding the resulting latency/loss numbers into the metrics pipelines Chapter 121 covers — turning this chapter's manual, one-off commands into an always-on early-warning system.
- **`iperf3` tests consume real bandwidth and should be run deliberately, not casually, on production links** — a saturating throughput test on a live production path can itself cause the congestion and loss it's trying to measure for other traffic sharing that link; scheduling tests during low-traffic windows, or against a dedicated test circuit, is standard practice.
- **Jitter buffers are tunable, and real VoIP/video systems expose that tuning.** Applications like WebRTC-based conferencing tools dynamically resize their jitter buffer based on recently observed jitter — a larger buffer trades latency for smoothness, and this adaptive behavior is precisely why the same network can produce a smooth call on one app and a choppy one on another with a less adaptive buffering strategy.
- **BDP-aware tuning matters on long-haul or satellite links specifically.** Cloud providers and CDNs (Chapter 96) explicitly tune TCP window scaling and sometimes deploy multiple parallel streams or QUIC's per-stream model (Chapter 75) specifically to avoid the single-connection BDP ceiling Section 10 described.

## 15. What's Simplified Here

Real jitter calculation (RFC 3550) uses a specific exponentially-weighted smoothing formula, simplified here to "variation in delay between consecutive packets" for intuition — the exact formula matters for interoperability between RTP implementations but not for understanding why it matters operationally. This chapter also treats latency, throughput, jitter, and loss as if cleanly separable; in reality they interact (severe loss increases effective latency through retransmission, e.g.), and a real diagnostic session, as Section 11 showed, usually has to measure several together rather than one in isolation. `mtr`'s per-hop numbers, as emphasized in Section 9 and 13, require real judgment to interpret correctly — the tool doesn't automatically tell you "hop 5 is broken," it gives you data that requires the same reasoning shown in this chapter to interpret correctly.

## 16. Interview Questions & Model Answers

**Beginner: What is the difference between bandwidth and throughput?**

*Model answer:* Bandwidth (Chapter 16) is the theoretical maximum data rate a link's physical properties allow — an upper ceiling. Throughput is the data rate actually achieved in practice, at a given time, including all protocol overhead — always less than or equal to bandwidth. "My internet plan is 500 Mbps but I'm only seeing 420 Mbps" is describing the gap between bandwidth (the plan's ceiling) and throughput (what's actually measured), and a moderate gap from protocol overhead is normal and expected.

**Intermediate: A VoIP call has an average one-way latency of only 60ms, well within acceptable limits, but users report it's unusable. What single additional measurement would you take, and why might it explain the complaint despite the good average latency?**

*Model answer:* I'd measure jitter directly (via `iperf3 -u` or examining `ping`'s mdev) rather than trusting the average latency alone. Since real-time audio must be played back at a fixed rate from a jitter buffer, high variation in per-packet delay — even with a fine average — causes packets to arrive too late for their playback deadline and be treated as lost, producing exactly the choppy, unusable experience described, despite an average latency number that looks completely healthy.

**Advanced: Using `mtr`, you observe 15% packet loss at hop 6 of a 10-hop path, with 0% loss at every hop before and after it. Is this hop the cause of real packet loss on the path? Justify your answer precisely.**

*Model answer:* Not necessarily, and the pattern actually argues against it: a hop with real packet loss would, by definition, cause every hop *after* it to also show elevated loss for the packets that continue past it — a packet cannot be measured as successfully reaching hop 7, 8, 9, and 10 if it was truly dropped at hop 6. Loss isolated to exactly one hop, with clean hops immediately afterward, is the classic signature of that specific router deprioritizing or rate-limiting its own ICMP TTL-exceeded replies (used by `mtr`/`traceroute` to build the report) rather than dropping the actual forwarded traffic — the trustworthy signal for real loss is loss that persists from a given hop through to the final destination.

## 17. Exercises

### Easy
1. Explain, precisely, the difference between latency and jitter.
2. A `ping` run reports `rtt min/avg/max/mdev = 20.1/20.4/20.9/0.2 ms`. Would you describe this path as having high or low jitter, and why?
3. Which single tool from this chapter would you use to measure packet loss at every hop along a path, not just the final destination?

### Medium
4. Explain why a UDP-mode `iperf3` test reports packet loss directly and explicitly, while a TCP-mode `iperf3` test does not — tie your answer to specific chapters covering each protocol.
5. Compute the bandwidth-delay product for a 1 Gbps link with a 150ms RTT, and explain in one sentence what a TCP connection needs (from Chapter 61) to actually use the full bandwidth on such a path.
6. An `mtr` report shows rising, persistent packet loss starting at hop 4 and continuing at a similar percentage through the final destination at hop 9. Explain why this pattern, unlike the isolated-hop pattern from Section 9's example, is trustworthy evidence of real loss occurring around hop 4.

### Hard
7. Design a measurement plan (naming specific tools and flags from this chapter) to determine whether a reported "slow website" complaint is caused by insufficient bandwidth, high latency, packet loss, or none of the above (i.e., an application-layer problem outside this chapter's scope). Justify the order you'd run your tests in.
8. A file transfer over a single TCP connection to a satellite-linked office achieves only 8 Mbps despite the link's advertised 50 Mbps bandwidth and near-zero measured packet loss. Using Section 10's bandwidth-delay product concept and Chapter 61, propose a specific, concrete fix and explain the mechanism by which it would help.
9. Explain why a jitter buffer that's tuned too large for the current network conditions can itself become the cause of a poor real-time call experience, even on a network with genuinely low jitter — tie your answer precisely to Section 7's mechanism.

## 18. Summary and the Bridge to Chapter 121

| Term | Precise Meaning | Primary Tool |
|---|---|---|
| Latency | one-way (or RTT for round trip) delay for a packet | `ping` |
| Bandwidth | theoretical maximum data rate of a link | link specification, Ch 16 |
| Throughput | actually achieved data rate, including overhead | `iperf3` |
| Goodput | useful application data rate, excluding overhead/retransmits | derived from throughput minus overhead |
| Jitter | variation in latency between consecutive packets | `iperf3 -u`, `ping`'s mdev |
| Packet loss | fraction of packets that never arrive | `iperf3 -u`, `mtr` |
| Bandwidth-delay product | bandwidth × RTT; data "in flight" capacity | derived; ties to Ch 61's window |
| `mtr` | continuous per-hop loss and latency, combining ping + traceroute | Section 9 |

You can now measure a network path precisely, on demand, whenever you run a command. But most real network problems don't announce themselves while you happen to be watching — they happen at 3 a.m., intermittently, on a device you don't have terminal access to. Chapter 121 covers how production infrastructure watches these same quantities continuously and automatically: SNMP polling individual devices, flow logs recording traffic patterns without full packet capture, and Prometheus and Grafana turning all of it into dashboards and alerts.
