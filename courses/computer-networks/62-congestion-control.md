# Chapter 62: Congestion Control — Slow Start, AIMD, CUBIC, and BBR

> **"Flow control protects the receiver. Congestion control protects the network. TCP had to learn the second lesson the hard way — by watching the early Internet nearly destroy itself."**

---

## Table of Contents

1. [The Question Flow Control Doesn't Answer](#1-the-question-flow-control-doesnt-answer)
2. [October 1986: Congestive Collapse](#2-october-1986-congestive-collapse)
3. [Two Windows, Not One: cwnd vs. rwnd](#3-two-windows-not-one-cwnd-vs-rwnd)
4. [Slow Start — Exponential Growth on Purpose](#4-slow-start--exponential-growth-on-purpose)
5. [Congestion Avoidance — AIMD](#5-congestion-avoidance--aimd)
6. [A Worked Example: The Sawtooth](#6-a-worked-example-the-sawtooth)
7. [The Limits of Reno-Style AIMD](#7-the-limits-of-reno-style-aimd)
8. [CUBIC — Built for Long, Fat Pipes](#8-cubic--built-for-long-fat-pipes)
9. [BBR — Modeling the Network Instead of Guessing](#9-bbr--modeling-the-network-instead-of-guessing)
10. [CUBIC vs. BBR, Side by Side](#10-cubic-vs-bbr-side-by-side)
11. [ECN: A Congestion Signal That Isn't a Loss](#11-ecn-a-congestion-signal-that-isnt-a-loss)
12. [Packet-Level View and Reading `ss -i`](#12-packet-level-view-and-reading-ss--i)
13. [Hands-On Experiment](#13-hands-on-experiment)
14. [Code: Simulating AIMD in Go](#14-code-simulating-aimd-in-go)
15. [Common Misconceptions](#15-common-misconceptions)
16. [Production Usage Notes](#16-production-usage-notes)
17. [Interview Questions & Model Answers](#17-interview-questions--model-answers)
18. [Exercises](#18-exercises)
19. [Summary](#19-summary)

---

## 1. The Question Flow Control Doesn't Answer

Chapter 61 solved a specific problem: a fast sender can overrun a slow *receiver*. The fix was the sliding window — the receiver advertises `rwnd`, the amount of buffer space it currently has, and the sender never has more than `rwnd` bytes unacknowledged in flight. That is **flow control**, and it is entirely a conversation between two endpoints. Nothing about `rwnd` knows or cares what is happening to the packets *between* those endpoints.

But there is a second, completely different way a sender can overwhelm something: it can overwhelm the **network itself**.

Picture two hosts connected not directly, but through a chain of routers — the normal case for anything beyond a LAN. Each router has a finite amount of buffer memory for packets waiting to go out an interface. If a sender pushes data faster than the slowest link along the path (the "bottleneck link") can drain it, packets pile up in that router's queue. If the sender keeps pushing, the queue fills completely, and the router starts dropping packets — not because the receiver's application buffer is full, but because a piece of physical hardware in the middle of the path ran out of memory.

The receiver in this scenario might have gigabytes of free buffer space and be advertising a huge `rwnd`. Flow control has nothing to say about this. A second, independent mechanism is needed to answer a different question: **not "can the receiver keep up?" but "can the network in between keep up?"** That mechanism is **congestion control**, and every modern TCP sender runs it, even when talking to a receiver with an infinite appetite for data.

This is worth stating as plainly as possible, because it is the single most common point of confusion in this part of TCP:

```
FLOW CONTROL                          CONGESTION CONTROL
-------------------------------       -------------------------------
Protects: the RECEIVER                Protects: the NETWORK (routers,
                                       links, everyone sharing them)

Signal used: rwnd, advertised by      Signal used: packet loss, delay,
the receiver in every ACK             or explicit marks (ECN) — inferred
                                       by the SENDER, not told by anyone

Who enforces it: the receiver         Who enforces it: the sender,
tells the sender its limit            policing itself against a limit
                                       it has to estimate

Where the limit lives: the            Where the limit lives: the
receiver's socket buffer              sender's estimate of the
                                       bottleneck link's capacity

Chapter that covers it: 61            Chapter that covers it: this one
```

The two limits combine. As you will see precisely in a moment, the sender's actual send window on any given round trip is:

```
effective window = min(cwnd, rwnd)
```

where `rwnd` is the receiver-advertised flow-control window from Chapter 61, and `cwnd` — the **congestion window** — is a value the sender maintains *itself*, based on what it has inferred about the network's capacity. Nobody tells the sender what `cwnd` should be. There is no "network window" field in any packet. The sender has to estimate it, adjust it continuously, and get it wrong gracefully. That estimation problem is the subject of this chapter.

---

## 2. October 1986: Congestive Collapse

This isn't a hypothetical problem invented to justify a chapter. It happened, publicly, to the early Internet.

In October 1986, a link between Lawrence Berkeley Laboratory and UC Berkeley — physically capable of 32 kbit/s — saw its effective throughput drop to about 40 bit/s. Not a typo: a thousand-fold collapse. The link wasn't broken. It was drowning in its own retransmissions.

Here is the mechanism, and it's worth understanding precisely because it explains why congestion control cannot be optional:

1. Early TCP implementations had no congestion control at all. A sender with data to send and a receiver with buffer space simply sent as fast as the flow-control window allowed.
2. When several senders shared a bottleneck link, the combined offered load exceeded the link's capacity. The router's queue filled and packets were dropped.
3. TCP, doing exactly what Chapter 60 says it should do, noticed the missing ACKs and **retransmitted** the lost data.
4. But the senders didn't reduce their sending rate before retransmitting — they just added the retransmissions on top of the same aggressive rate. This *increased* the offered load on an already-overloaded link.
5. More packets dropped. More retransmissions. The link spent almost all of its capacity re-sending data that had already been sent at least once, while almost no new, useful data got through. Throughput on the wire could stay high while **goodput** — data that actually made it through exactly once — cratered.

This is **congestive collapse**: a positive feedback loop where the network's own repair mechanism (retransmission) becomes the thing that keeps it congested. Van Jacobson, watching this happen, wrote the fix — the algorithms in this chapter — and shipped it in 4.3BSD Tahoe in 1988. It is not an exaggeration to say that TCP's congestion control is one of the reasons the public Internet survived being a shared, uncoordinated, best-effort network at all. Nobody polices how fast you're allowed to send; every well-behaved TCP stack polices itself, all the time, cooperatively, because Jacobson's algorithms made self-restraint the only way to get good throughput.

---

## 3. Two Windows, Not One: cwnd vs. rwnd

Before the algorithms, the bookkeeping. A sending TCP stack maintains (at minimum):

- **`rwnd`** — the receiver's advertised window (Chapter 61). Told to the sender explicitly, in the Window field of every segment (Chapter 65 covers that field precisely).
- **`cwnd`** — the congestion window. A number of bytes (conceptually; some descriptions use segments), computed and updated entirely by the sender's own logic, based on what it infers about the network.
- **`ssthresh`** — "slow start threshold." A byte value marking the sender's current best guess at where the network's capacity roughly is, used to decide which of the two algorithms below should currently be running.

At every point, the sender is allowed to have this many unacknowledged bytes in flight:

```
effective window = min(cwnd, rwnd)
```

If the receiver has a huge buffer (`rwnd` is large) but the network path is thin, `cwnd` is the binding constraint — the sender is throttled by *its own estimate of the network*, not by anything the receiver said. This is by far the common case on the modern Internet, where receive buffers are usually generous. If you've ever wondered why a fast machine downloading from a fast server over a slow, congested link still trickles data in slowly even though both endpoints have plenty of RAM — this is why. `cwnd` is doing its job.

`cwnd` is not sent in any packet. It exists only in the sender's memory. You cannot observe another host's `cwnd` from a packet capture; you can only infer it by watching how many unacknowledged bytes it's willing to have outstanding. On Linux, you *can* read your own kernel's live `cwnd` for a connection with `ss -i` (Section 12).

---

## 4. Slow Start — Exponential Growth on Purpose

**The problem restated:** a brand-new TCP connection has zero information about the path's capacity. It might be a data-center link capable of 100 Gbit/s, or a congested mobile link capable of 200 kbit/s. Guess too high, and the very first burst of data can overflow a router queue before a single ACK has come back to say "slow down." Guess too low, and a fast, empty link sits idle for no reason.

**A naive attempt:** start by sending as much as `rwnd` allows, since that's the only limit that existed before congestion control. Section 2 showed exactly why this fails — it is precisely the behavior that caused congestive collapse.

**Another naive attempt:** start with a fixed, conservative rate (say, one segment per RTT, linear increase from there) and creep up slowly. This is *safe* but wastes enormous amounts of capacity on any link with real bandwidth — on a modern data-center or fiber link, linear ramp-up would take so long that most short connections (a web page fetch, an API call) would finish before ever reaching a useful sending rate.

**The real solution — slow start:** begin conservatively, but grow *exponentially*, doubling the congestion window roughly every round trip, until either the network signals trouble (a loss) or the sender reaches `ssthresh`, its current estimate of safe capacity.

The mechanics:

- `cwnd` starts at an **initial window (IW)**. Historically 1 MSS; RFC 6928 (2013) raised the modern recommended default to **10 MSS** (about 14,600 bytes at a 1460-byte MSS), reflecting the fact that most modern paths have far more capacity than 1980s links did, and most Internet transfers (like fetching a small web page) are short enough that a bigger initial burst finishes faster without risking much.
- For every ACK received that acknowledges new data, `cwnd` increases by 1 MSS: `cwnd += MSS`.
- Because the whole current window's worth of segments produces roughly that many ACKs within one round trip, and each ACK adds 1 MSS to `cwnd`, the practical effect is that `cwnd` **doubles every RTT**: 10 → 20 → 40 → 80 MSS, and so on. This is why the phase is called "slow start" — a slightly misleading name, since after the first RTT the growth is anything but slow. Jacobson's own note was that it starts more slowly than "blast the whole window immediately," which is the comparison that mattered in 1988.
- Slow start ends the moment one of two things happens:
  - `cwnd` reaches `ssthresh` — the sender transitions to congestion avoidance (Section 5).
  - A loss is detected — the sender treats this as its first real signal about the network's actual capacity, sets `ssthresh` to (roughly) half the current `cwnd`, and either restarts slow start from a small window (classic timeout recovery) or enters fast recovery (Chapter 63 covers this precisely — it's a large enough topic to deserve its own chapter).

```
cwnd growth during slow start (IW = 10 MSS, MSS = 1460 bytes)

RTT 0:  cwnd = 10 MSS  (14,600 bytes)
RTT 1:  cwnd = 20 MSS  (29,200 bytes)
RTT 2:  cwnd = 40 MSS  (58,400 bytes)
RTT 3:  cwnd = 80 MSS  (116,800 bytes)
RTT 4:  cwnd = 160 MSS (233,600 bytes)
   ...doubling each RTT until ssthresh or a loss...
```

**Three levels:**

- **Intuitive:** you don't floor the accelerator the first time you drive an unfamiliar road in the dark — you creep forward, and as nothing goes wrong, you speed up quickly, doubling your confidence (and speed) each time you don't crash. The moment something goes wrong, you slam the brakes and recalibrate.
- **Engineering:** exponential window growth, bounded by ACK clocking — the arrival of ACKs (which can only arrive as fast as the network delivers data) is what paces the growth, so the algorithm is naturally self-limiting to what the path can currently sustain.
- **Deep technical:** slow start is a **probing algorithm**. It has no direct measurement of link capacity, so it uses binary-search-like exponential growth to find the *order of magnitude* of available capacity quickly, then hands off to a much more careful, linear algorithm (congestion avoidance) to fine-tune around that estimate. The `ssthresh` value is exactly the sender's memory of "the last place a loss happened" — an estimate that gets revised every time the network provides new evidence.

---

## 5. Congestion Avoidance — AIMD

Once `cwnd` reaches `ssthresh` (or after a loss has already forced a reset), doubling every RTT is too aggressive — the sender is now operating near its best estimate of the network's actual capacity, and continuing to double risks overshooting it violently. The algorithm switches to **congestion avoidance**, governed by **AIMD — Additive Increase, Multiplicative Decrease**:

- **Additive increase:** for every full RTT during which all data was acknowledged without loss, `cwnd += MSS` (roughly — the canonical per-ACK formula is `cwnd += MSS * MSS / cwnd`, which works out to approximately one MSS of growth per round trip regardless of how many ACKs arrive within it). This is *linear* growth — cautious, steady probing for a little more capacity than last time.
- **Multiplicative decrease:** the instant a loss is detected, cut `cwnd` in half: `cwnd = cwnd / 2` (and set `ssthresh` to that same halved value). This is an aggressive, immediate retreat — the sender assumes the loss means it just found the network's ceiling, and backs off hard rather than testing the theory further.

Why this exact asymmetric shape — slow climb, fast retreat — and not something more symmetric? This isn't an arbitrary choice; it comes out of a real piece of control theory, the **Chiu-Jain analysis** (1989): among the simple linear control rules a sender could use, **only** additive-increase/multiplicative-decrease provably converges every competing flow sharing a bottleneck toward an equal, fair share of that link's capacity over time, regardless of where each flow started. Additive-increase-additive-decrease doesn't converge to fairness as reliably; multiplicative-increase-multiplicative-decrease amplifies any initial unfairness instead of correcting it. AIMD is the shape that self-corrects.

The intuition for *why* AIMD is fair: imagine two flows sharing one bottleneck link, one currently using more of it than the other. Every RTT, both flows add the *same fixed amount* (1 MSS) — so the flow with the smaller window grows by a *larger percentage* of itself than the flow with the bigger window. That nudges them toward equality. When a loss eventually happens (as it must, since both keep growing until the link's queue overflows), both flows are cut by the *same fraction* (half) — which preserves whatever ratio additive increase had already been closing. Repeated over many cycles, this converges the two flows toward roughly equal shares. Chiu and Jain proved this converges from *any* starting point, which is exactly the property you need in a network with no central coordinator.

---

## 6. A Worked Example: The Sawtooth

Put slow start and AIMD together and plot `cwnd` over time, and you get TCP's signature shape — the **congestion window sawtooth**. Here is a worked example with concrete numbers: `MSS = 1460 bytes`, initial `ssthresh = 64` (an arbitrary but plausible estimate), starting `cwnd = 10 MSS`.

```
RTT   Phase                cwnd (MSS)   Event
---   -----                ----------   -----
 0    slow start           10           —
 1    slow start           20           —
 2    slow start           40           —
 3    slow start (capped)  64           cwnd hits ssthresh → switch to AIMD
 4    congestion avoidance 65           +1 MSS
 5    congestion avoidance 66           +1 MSS
 6    congestion avoidance 67           +1 MSS
  ...(linear climb continues)...
40    congestion avoidance 101          LOSS DETECTED
41    congestion avoidance 50           cwnd halved; ssthresh = 50
42    congestion avoidance 51           +1 MSS, climbing again
  ...(linear climb resumes from the new, lower baseline)...
76    congestion avoidance 85           LOSS DETECTED
77    congestion avoidance 42           cwnd halved; ssthresh = 42
```

Plotted, `cwnd` traces a shape that looks like the teeth of a saw: a slow initial ramp (the exponential slow-start curve), then a long, gentle, straight-line climb, then a sudden vertical drop to half its peak, then climbing again from there.

```
cwnd
(MSS)
 101 |                                   /|
     |                                 /  |
  85 |                               /    |    /|
     |                             /      |  /  |
  64 |          /                /        |/    |
     |        /                /         (loss)(loss)
  50 |      /                 |
     |    /                  (loss, halve)
  20 |  /
     |/
  10 +--------------------------------------------------> time (RTTs)
     0    1    2    3    4 ...           40  41       76 77
     |--- slow start ---|--- congestion avoidance (AIMD) ---|
```

This sawtooth is the visible fingerprint of a *loss-based* congestion controller — the sender keeps deliberately pushing a little past the network's comfortable capacity, waiting for the network to say "too far" via a dropped packet, backing off, and trying again. It is, by design, always either under-using the link slightly or briefly over-using it. That design choice is exactly what Section 7 and Section 9 come back to challenge.

---

## 7. The Limits of Reno-Style AIMD

The algorithm in Sections 4–5 (with the exact recovery mechanics filled in by Chapter 63) is known as **TCP Reno**, and it served the Internet well for over a decade. But it has real, well-documented limits that became increasingly painful as networks changed:

**1. It's slow to fill high-bandwidth, high-latency links (the "long fat network" problem).** Additive increase adds a fixed 1 MSS per RTT, no matter how big `cwnd` already is or how much spare capacity the link has. On a link with a large **bandwidth-delay product** (BDP — the amount of data that can be "in flight" at once, equal to bandwidth × RTT) — say, a 10 Gbit/s transatlantic link with a 150ms RTT, whose BDP is on the order of 187 MB — climbing from a modest `cwnd` back up to a `cwnd` that uses the full link after every loss, at 1 MSS per RTT, can take **thousands of round trips** — tens of minutes of leaving most of an expensive, fast link idle. Reno was tuned for the modest links of the late 1980s, not modern long-haul fiber.

**2. It punishes short RTTs and rewards long ones, unevenly.** Because additive increase happens once per RTT, a flow with a 10ms RTT gets to add 1 MSS ten times as often as a flow with a 100ms RTT sharing the same bottleneck — the short-RTT flow grows its window, and therefore its share of the link, much faster. This is called **RTT unfairness**, and it means Reno doesn't actually deliver the equal-share fairness Chiu-Jain promised once real, unequal RTTs enter the picture.

**3. It treats all loss as a congestion signal, and all loss is not congestion.** On a wireless link, a packet can be corrupted by radio interference and dropped with zero relation to how full any router's queue is. A pure loss-based algorithm halves `cwnd` anyway, throttling itself for no real reason — a persistent problem for cellular and Wi-Fi paths.

**4. It relies on queues filling up before it learns anything — a phenomenon later named "bufferbloat."** Loss-based algorithms need a router buffer to actually overflow before they get a signal to slow down. Modern routers and cheap consumer equipment often have very large buffers (to smooth out bursts), which means a loss-based sender keeps filling that buffer for a long time before it overflows — adding huge amounts of queueing delay to every flow sharing that link, well before any packet is actually lost. Interactive traffic (a VoIP call, an SSH session) sharing a link with a large bulk transfer can experience latency spikes of hundreds of milliseconds or more, purely because a Reno-style flow is deliberately keeping the queue as full as it can get away with.

These limits are exactly what came next — Section 8 and Section 9 — were built to fix.

---

## 8. CUBIC — Built for Long, Fat Pipes

**The problem CUBIC targets specifically:** Reno's fixed "+1 MSS per RTT" climb is far too timid on the high-bandwidth, high-latency links that became increasingly common as data centers, backbone links, and international fiber routes got faster through the 2000s (this is literally called the "long fat network," or LFN, problem in the RFCs). CUBIC (Ha, Rhee, Xu, 2008) is designed to recover *aggressively* toward the last known-good window size after a loss, then probe carefully and increasingly aggressively for more, **independent of RTT** — directly attacking limit #1 and limit #2 from Section 7.

**The core idea:** instead of growing `cwnd` linearly with each ACK (which ties growth rate to how many ACKs arrive, which ties it to RTT), CUBIC grows `cwnd` as a **cubic function of real, wall-clock time** since the last congestion event. The formula:

```
W(t) = C(t - K)^3 + W_max

where:
  W_max = the cwnd value at the moment of the last loss
          (the window size the network last complained about)
  t     = real elapsed time since the last window reduction
  C     = a scaling constant (standardized default: 0.4)
  K     = the time it takes the cubic function to reach W_max again,
          computed as  K = cube_root(W_max * beta / C)
          where beta (the multiplicative decrease factor) is 0.7 for
          CUBIC — a gentler cut than Reno's 0.5
```

What this formula produces, walked through in plain language:

1. **Immediately after a loss**, `cwnd` is cut to `W_max * 0.7` (not halved as aggressively as Reno's 0.5 — CUBIC backs off less on the theory that the loss doesn't necessarily mean the network is drastically over capacity).
2. **Concave region:** right after the cut, the cubic curve rises *quickly* — because `t` is close to `0` and `(t-K)` is a large negative number, but cubed and added to `W_max`, this produces fast initial regrowth back toward the window size that was working before. This is the opposite of Reno's uniform, slow linear climb — CUBIC assumes "we were just fine at `W_max` a moment ago, let's get back there quickly."
3. **Plateau near `W_max`:** as `t` approaches `K`, growth flattens almost to zero right around the previous ceiling — the algorithm is cautious exactly where it has the most reason to believe a loss might happen again.
4. **Convex region:** past `t = K`, the curve starts climbing again, slowly at first and then increasingly steeply — probing for *new* capacity beyond what worked last time, on the theory that if it's been a while since the last loss, the network may now have more room (other flows may have finished, or link capacity may simply be greater than last measured).

The critical property this delivers: **CUBIC's growth curve depends only on elapsed real time since the last loss, not on RTT.** Two CUBIC flows sharing a bottleneck, one with a 10ms RTT and one with a 100ms RTT, converge toward the *same* window trajectory over wall-clock time — solving Reno's RTT-unfairness problem directly. And because the convex region's growth accelerates the longer it's been since a loss, CUBIC can climb to fill a large BDP link far faster than Reno's flat, RTT-gated +1-MSS-per-RTT ever could.

**CUBIC is the default congestion control algorithm in the Linux kernel** since 2.6.19 (2006), which — because Linux runs the overwhelming majority of the world's servers — makes it, in practical terms, the most-used TCP congestion control algorithm on the Internet today.

---

## 9. BBR — Modeling the Network Instead of Guessing

CUBIC is a better *loss-based* algorithm, but it is still fundamentally loss-based — it still needs the network to drop a packet before it learns anything, and it will still, by design, keep growing until that happens (limit #4 from Section 7, bufferbloat, is untouched by CUBIC). **BBR** (Bottleneck Bandwidth and Round-trip propagation time; Cardwell, Cheng, Yeganeh, Jacobson at Google, published 2016) starts from a completely different premise:

**The core idea:** instead of inferring network capacity indirectly from *when it breaks* (loss), directly **measure** the two physical quantities that actually define a path's capacity, and compute a sending rate from them instead of a window size reactively adjusted after the fact:

- **`BtlBw` (bottleneck bandwidth):** the maximum delivery rate observed recently — the fastest rate at which acknowledgments have confirmed data actually got through. This is measured continuously by comparing bytes delivered against the time it took to deliver them.
- **`RTprop` (round-trip propagation time):** the minimum RTT observed recently — specifically the RTT *with the queue empty*, i.e., not including any queueing delay. Because RTT = propagation delay + queueing delay, and queueing delay can only ever be added on top of the true minimum, the smallest RTT seen over a sufficiently long window is a good estimate of the path's true propagation delay.

Multiplying an accurate `BtlBw` by an accurate `RTprop` gives BBR a direct estimate of the ideal BDP — the amount of data that should be in flight to exactly fill the pipe without building a queue. BBR then **paces** packets out at a rate close to `BtlBw` (rather than sending a whole window in a burst and waiting) and keeps just enough data in flight to match the BDP — no more, which means, unlike Reno or CUBIC, **BBR does not deliberately try to fill router queues at all.** It tries to stay right at the edge of the pipe's true capacity without spilling into the queue, which is precisely the behavior that avoids bufferbloat.

BBR operates as a state machine cycling through four phases:

```
STARTUP    Doubles the sending rate each round trip (similar in spirit
           to slow start) until BtlBw stops increasing — i.e., until
           it detects the pipe is full — rather than waiting for a loss.

DRAIN      Deliberately sends at a reduced rate for one RTT to drain
           any queue that STARTUP's aggressive doubling built up,
           bringing queueing delay back down to (near) zero.

PROBE_BW   The steady-state phase. Cycles the pacing rate through a
           repeating pattern — briefly above the estimated BtlBw (to
           test whether more bandwidth has become available), then
           briefly below it (to drain any queue that probe created),
           then at the estimated rate for the remaining cycles. This
           is how BBR keeps re-discovering capacity without ever
           needing a loss to tell it "you went too far."

PROBE_RTT  Periodically (roughly every 10 seconds), BBR deliberately
           cuts the in-flight data down to a very small amount for a
           short interval, specifically to get a fresh, uninflated
           RTT sample — guarding against RTprop's estimate slowly
           drifting upward because the queue never fully empties
           during ordinary operation.
```

**Three levels:**

- **Intuitive:** Reno and CUBIC are like driving a fogged-in road by tapping the brakes every time you feel a bump (loss) and then speeding back up — reactive, and you always find the pothole. BBR is like using working headlights and a speedometer: it measures the road's actual speed limit and its actual distance, and drives right at the limit without needing to hit anything to find out where it was.
- **Engineering:** a model-based rate controller. It maintains running estimates of two physical quantities (`BtlBw`, `RTprop`) instead of a single reactive window, and paces packets at a computed rate instead of releasing a burst equal to a window size every RTT.
- **Deep technical:** BBR fundamentally changes what "full utilization" and "low latency" mean simultaneously achievable — Kleinrock's earlier work on network delay showed that the *optimal* operating point for a network path is exactly at `BtlBw × RTprop`: any less underutilizes the link, any more only adds queueing delay without adding throughput. Loss-based algorithms, by design, overshoot this point on every cycle (that's what builds the sawtooth's peak) and then undershoot it after every cutback. BBR tries to sit at the optimal point directly, which is why it both keeps throughput high *and* keeps latency low — a combination Reno/CUBIC structurally cannot offer, because their only tool for learning "am I at capacity yet" is overshooting it.

---

## 10. CUBIC vs. BBR, Side by Side

```
                     CUBIC (Reno family, loss-based)   BBR (model-based)
------------------  ---------------------------------  --------------------------------
Primary signal      Packet loss (a router queue        Measured delivery rate (BtlBw)
                    actually overflowed)                and measured min RTT (RTprop)

Behavior toward     Deliberately grows cwnd until       Deliberately avoids filling
router queues       a queue overflows — treats a        queues; targets the BDP, the
                    full queue as "found the limit"     point of full utilization with
                                                         an empty queue

RTT fairness        Fixed by design (growth is a        Also RTT-independent — driven
                    function of wall-clock time,         by directly measured bandwidth,
                    not RTT)                             not by RTT-gated window growth

Reaction to         Backs off (wrongly assumes           Largely unaffected — a lone
non-congestion      congestion) — a real problem on       corrupted packet doesn't
loss (e.g. Wi-Fi    lossy wireless/cellular links         change the BtlBw/RTprop
corruption)                                               estimates much

Where it shines     Bulk transfers on well-provisioned   Lossy/variable links (mobile,
                    wired/fiber paths — data centers,     satellite), and latency-
                    backbone transfers, most servers      sensitive + throughput-hungry
                                                           mixes (video streaming)

Default on          Linux kernel (since 2.6.19)          Used at large scale by Google
                                                           (YouTube, Google Cloud), available
                                                           as a selectable Linux module
                                                           (tcp_bbr) since kernel 4.9

Known criticism      Deliberately induces queueing       Early BBR (v1) could be
                     delay under load (bufferbloat);      unfair to loss-based flows
                     slower to react to real congestion   sharing the same bottleneck,
                     that is genuinely link capacity      and could underestimate BtlBw
                     shrinking (e.g. a new competing       on very lossy links; BBRv2/v3
                     flow joining)                         address this with more
                                                           conservative pacing gains
```

Why do algorithms keep evolving instead of one "final" answer being declared? Because the network keeps changing underneath every assumption: link speeds grew by six orders of magnitude since Reno was written, wireless became a dominant access medium with loss patterns Reno never anticipated, and router buffer sizes grew in ways that made loss-based signaling arrive later and later. CUBIC and BBR are each the right tool for different, real, still-coexisting conditions — which is exactly why Linux lets you choose (`net.ipv4.tcp_congestion_control`), and why large operators (Google notably) run different algorithms for different traffic classes.

---

## 11. ECN: A Congestion Signal That Isn't a Loss

One more piece completes the picture, and it connects directly to the header fields Chapter 65 will lay out precisely: **Explicit Congestion Notification (ECN)**, defined in RFC 3168. ECN lets a router that is starting to get congested — but hasn't yet had to drop a packet — **mark** a packet instead of dropping it, by setting two bits in the IP header. When the marked packet arrives, the receiving TCP stack echoes that mark back to the sender using two flag bits in the TCP header itself: **ECE** (ECN-Echo) and **CWR** (Congestion Window Reduced). A sender that sees an ECE-marked ACK reduces `cwnd` exactly as if a loss had occurred — but without ever actually having to drop and retransmit a packet, avoiding the wasted round trip and the retransmission overhead that a real loss would have cost. This is the single cleanest fix available for the loss-based approach's core weakness: it lets the network say "please slow down" *before* things get bad enough to require a drop, if both the router and both endpoints support it (ECN has to be negotiated at connection setup and is not universally enabled, though support has grown substantially since the 2010s).

---

## 12. Packet-Level View and Reading `ss -i`

Congestion control leaves **no dedicated packet field** for `cwnd` itself — it is pure sender-side state, invisible on the wire except through its effects (how much data the sender actually transmits before waiting for ACKs) and through the ECE/CWR flag bits when ECN is in play. On Linux, you can read your own kernel's live view of a connection's congestion state directly:

```
$ ss -i dst example.com

tcp   ESTAB  0  0  10.0.0.5:51322  93.184.216.34:443
     cubic wscale:7,7 rto:212 rtt:10.4/1.2 mss:1418 pmtu:1500
     rcvmss:1418 advmss:1418 cwnd:36 ssthresh:22 bytes_sent:187342
     bytes_acked:187342 bytes_received:942811 segs_out:412 segs_in:801
     data_segs_out:132 data_segs_in:0 send 490.3Mbps lastsnd:4
     lastrcv:4 lastack:4 pacing_rate 980.5Mbps delivery_rate 465.2Mbps
     busy:1284ms rwnd_limited:0ms(0.0%) unacked:1 rcv_space:65536
     rcv_ssthresh:65536 minrtt:9.887
```

Reading the fields that matter for this chapter: `cubic` names the active congestion control algorithm for this socket; `cwnd:36` is the live congestion window in segments (36 × 1418-byte MSS ≈ 51 KB currently allowed in flight); `ssthresh:22` shows the last recorded slow-start threshold; `pacing_rate` and `delivery_rate` are BBR-style measured rates that modern Linux tracks and reports even for non-BBR connections; `minrtt` is exactly the `RTprop`-style minimum RTT estimate described in Section 9. If the algorithm were BBR, `ss -i` would show `bbr` in place of `cubic` along with additional BBR-specific fields.

---

## 13. Hands-On Experiment

You can watch congestion control behave in real time on Linux without any special tools beyond `ss` and `tc` (traffic control, the kernel's network emulation facility):

```bash
# 1. See which congestion control algorithm your kernel uses by default
sysctl net.ipv4.tcp_congestion_control
# cubic   (on most stock Linux distributions)

# 2. See which algorithms are available to switch between
sysctl net.ipv4.tcp_available_congestion_control
# reno cubic bbr

# 3. Start a long transfer to a real server and watch cwnd grow live
curl -o /dev/null http://speedtest-server.example/100MB.bin &
watch -n 0.2 'ss -i dst speedtest-server.example'
# Watch cwnd climb, notice ssthresh, watch for a drop after any loss

# 4. Simulate a lossy link with netem and watch the same transfer's
#    cwnd get halved every time a simulated loss occurs
sudo tc qdisc add dev eth0 root netem loss 2%
curl -o /dev/null http://speedtest-server.example/100MB.bin &
watch -n 0.2 'ss -i dst speedtest-server.example'
sudo tc qdisc del dev eth0 root netem   # clean up afterward

# 5. Switch this machine's default algorithm and repeat the test
sudo sysctl -w net.ipv4.tcp_congestion_control=bbr
```

With `netem loss 2%` active, watch `cwnd` in the `ss -i` output repeatedly climb and then drop sharply — the sawtooth from Section 6, made visible on a real, if artificially degraded, connection.

---

## 14. Code: Simulating AIMD in Go

The core AIMD loop is simple enough to simulate directly, which is a good way to build intuition for the sawtooth shape and how quickly loss rate determines average throughput:

```go
package main

import "fmt"

// simulateAIMD models a simplified Reno-style AIMD congestion window
// over a sequence of simulated RTTs, given a fixed loss pattern.
func simulateAIMD(rtts int, lossEveryNRTTs int, initCwnd, ssthresh float64) {
	cwnd := initCwnd
	inSlowStart := true

	for rtt := 1; rtt <= rtts; rtt++ {
		lost := lossEveryNRTTs > 0 && rtt%lossEveryNRTTs == 0

		if lost {
			ssthresh = cwnd / 2
			cwnd = ssthresh // Reno-style multiplicative decrease
			inSlowStart = false
			fmt.Printf("RTT %3d: LOSS -> cwnd=%.1f ssthresh=%.1f\n", rtt, cwnd, ssthresh)
			continue
		}

		if inSlowStart {
			cwnd *= 2 // exponential growth
			if cwnd >= ssthresh {
				cwnd = ssthresh
				inSlowStart = false
			}
		} else {
			cwnd += 1 // additive increase, ~1 MSS per RTT
		}

		fmt.Printf("RTT %3d: cwnd=%.1f (%s)\n", rtt, cwnd,
			map[bool]string{true: "slow start", false: "congestion avoidance"}[inSlowStart])
	}
}

func main() {
	// Initial window 10 MSS, ssthresh 64 MSS, a simulated loss every 30 RTTs
	simulateAIMD(90, 30, 10, 64)
}
```

Running this with different values of `lossEveryNRTTs` demonstrates something worth internalizing: average throughput under AIMD is roughly proportional to `MSS / (RTT * sqrt(loss_rate))` — a relationship known as the **Mathis equation** (1997). It formalizes exactly what the simulation shows visually: halving the loss rate roughly doubles average throughput (a square-root relationship, not linear), and this is also *why* Reno performs so poorly on the long-fat-network links from Section 7 — with a large RTT in that denominator, even a tiny loss rate caps throughput far below the link's real capacity, motivating both CUBIC's RTT-independent growth and BBR's abandonment of loss as the primary signal entirely.

---

## 15. Common Misconceptions

- **"The Window field in a TCP header is the congestion window."** No — the Window field (Chapter 65) is `rwnd`, the *receiver's* flow-control window. `cwnd` is never transmitted; it's private sender-side state, observable only indirectly (via `ss -i` on your own machine).
- **"Faster round trips always mean faster growth for everyone equally."** True for Reno-style AIMD (which is exactly the RTT-unfairness problem in Section 7), but explicitly *not* true for CUBIC or BBR, both engineered specifically to be RTT-independent.
- **"Packet loss always means the network is congested."** Not on wireless or satellite links, where corruption-based loss is common and unrelated to queue occupancy — precisely the scenario BBR was partly designed to handle better than loss-based algorithms.
- **"Slow start is slow."** It's exponential growth — the fastest a sender can ramp up without any prior information. It's slow only compared to "send everything at once," which is what it's replacing.
- **"BBR ignores loss completely."** Modern BBR (v2/v3) does incorporate some loss and ECN signals as guardrails; the defining difference from CUBIC is that loss is not BBR's *primary* signal, not that BBR is blind to it entirely.
- **"CUBIC and Reno are unrelated algorithms."** CUBIC is a direct descendant of Reno's AIMD family — same fundamental idea (loss-triggered multiplicative decrease, probe for more afterward) with the linear increase function replaced by a cubic one.

---

## 16. Production Usage Notes

- Check and set the algorithm on Linux with `sysctl net.ipv4.tcp_congestion_control` / `net.ipv4.tcp_available_congestion_control`; per-socket overrides are possible via `setsockopt(fd, IPPROTO_TCP, TCP_CONGESTION, "bbr", ...)` for applications that want a non-default algorithm for specific connections.
- Google has published production experience running BBR at scale for Google.com, YouTube, and Google Cloud, reporting meaningfully higher throughput and lower latency on lossy and long-distance paths compared to CUBIC — one of the reasons BBR moved from a research proposal to a widely deployed algorithm within a few years.
- Content-heavy platforms (video streaming particularly) often prefer BBR-family algorithms specifically because they avoid the latency spikes loss-based algorithms can cause under sustained bulk transfer — important when a video stream on the same network as an interactive session shouldn't make that session laggy.
- Mixing algorithms on a shared bottleneck matters: an early, valid criticism of BBRv1 was that it could take a larger, unfair share of a bottleneck when directly competing with CUBIC flows, because BBR doesn't back off in response to loss the way CUBIC does. This motivated the BBRv2/BBRv3 revisions, which add loss- and ECN-aware guardrails specifically to coexist better with the still-dominant CUBIC traffic on the Internet.
- Initial window size (`IW`) is itself a tunable and a real point of production engineering — RFC 6928's 10-MSS default reflects a real, measured tradeoff for the modern web (finish short flows fast) vs. the original RFC 2581's conservative 1–4-MSS default (safer for a 1999 Internet).

---

## 17. Interview Questions & Model Answers

**Q (Beginner): What is the difference between flow control and congestion control?**

*Model answer:* "Flow control protects the receiver from being overwhelmed — the receiver advertises `rwnd`, how much buffer space it has, and the sender never exceeds that. Congestion control protects the network itself — routers and links between the two endpoints — from being overwhelmed. The sender maintains its own `cwnd` estimate of safe network capacity, and the actual amount of data allowed in flight is `min(cwnd, rwnd)`. Flow control is a message from the receiver; congestion control is the sender policing itself based on inferred signals like loss."

**Q (Intermediate): Walk through what happens to `cwnd` from the start of a connection through its first packet loss.**

*Model answer:* "The connection starts in slow start with `cwnd` at the initial window — commonly 10 MSS today. Each ACK for new data adds 1 MSS to `cwnd`, which because roughly the whole window's worth of ACKs arrive per RTT, means `cwnd` doubles each round trip. This continues until either `cwnd` reaches `ssthresh` (switching to linear congestion avoidance / AIMD) or a loss occurs. On loss, `ssthresh` is set to about half the current `cwnd`, and depending on the recovery mechanism (classic timeout vs. fast retransmit/fast recovery, covered in Chapter 63), `cwnd` either resets to a small value and slow-starts again, or drops to the new `ssthresh` and resumes linear growth directly."

**Q (Advanced): Why was CUBIC created given that Reno's AIMD already existed, and how does BBR differ fundamentally from both?**

*Model answer:* "Reno's additive increase adds a fixed 1 MSS per RTT regardless of how large `cwnd` already is, which on high-bandwidth, high-latency 'long fat network' links takes an impractically long time to refill the window after a loss — the bandwidth-delay product can be enormous, and Reno's growth rate doesn't scale with it. Reno's growth is also RTT-dependent, so flows with shorter RTTs unfairly out-grow flows with longer RTTs sharing the same bottleneck. CUBIC fixes both by making window growth a cubic function of wall-clock time since the last loss rather than a function tied to RTT and per-ACK increments — it recovers quickly toward the last known-good window and then probes increasingly aggressively for more, independent of RTT. BBR is a more fundamental departure from both: rather than reactively growing a window until a loss occurs, it actively measures the bottleneck bandwidth and minimum RTT and computes a target sending rate and in-flight data amount directly from those measurements — deliberately trying to operate at the point of full utilization without ever intentionally filling a router's queue, which is what gives it both high throughput and low latency simultaneously, a combination loss-based algorithms structurally cannot achieve since they require overshooting capacity to detect it."

---

## 18. Exercises

### Easy

1. Explain, in one or two sentences each, why a receiver with a very large `rwnd` can still limit a TCP connection's throughput to well below what the receiver's link could otherwise handle.
2. Given `MSS = 1000 bytes` and an initial window of 10 MSS, write out `cwnd` for the first 5 RTTs of slow start assuming no loss occurs and `ssthresh` is never reached.
3. In your own words, explain why AIMD is asymmetric — slow additive growth but fast multiplicative cutback — rather than symmetric in both directions.

### Medium

4. A connection has `cwnd = 80 MSS` when a loss occurs. Using classic Reno-style multiplicative decrease, what are the new values of `cwnd` and `ssthresh`? If congestion avoidance then proceeds for 20 more RTTs without further loss, what is `cwnd` at the end?
5. Explain the RTT-fairness problem with Reno's AIMD using a concrete numeric example: two flows sharing a bottleneck, one with a 20ms RTT and one with a 200ms RTT, both in congestion avoidance. Over 100 real-time milliseconds, roughly how many additive-increase steps does each flow get, and what does that imply about their relative window sizes over time?
6. Using the CUBIC formula in Section 8, explain qualitatively (no need to compute exact numbers) what happens to the growth curve's shape if a second loss occurs *before* `t` reaches `K` from the first loss (i.e., before the window has climbed back to `W_max`).

### Hard

7. Explain precisely why BBR does not need packet loss as its primary signal, referencing both `BtlBw` and `RTprop` and how they're each measured. Then explain a realistic scenario where a BBR flow and a CUBIC flow share one bottleneck link and describe, mechanistically, why BBR might claim more than an equal share of that link's capacity.
8. Using the Mathis equation intuition from Section 14 (`throughput ∝ MSS / (RTT × sqrt(loss_rate))`), explain why a satellite link with a 600ms RTT and a 0.5% packet loss rate performs so much worse under classic Reno-style AIMD than a data-center link with a 1ms RTT and the same 0.5% loss rate — and explain concretely which of Reno's design assumptions this scenario violates.
9. Design (in pseudocode or prose) a minimal PROBE_BW-style cycling scheme: describe how you would periodically send slightly faster than your current bandwidth estimate to test for more capacity, and slightly slower afterward to drain any queue you built while probing, without ever deliberately causing a router queue to overflow.

---

## 19. Summary

| Term | Meaning |
|---|---|
| Congestion control | Sender-side self-throttling to avoid overwhelming the *network*, distinct from receiver-protecting flow control (Ch. 61) |
| `cwnd` | Congestion window — sender's own estimate of safe in-flight data; never transmitted on the wire |
| `rwnd` | Receiver's advertised flow-control window (Ch. 61); effective send limit is `min(cwnd, rwnd)` |
| `ssthresh` | Slow-start threshold — the boundary between slow start and congestion avoidance |
| Congestive collapse | The 1986 real-world failure mode where uncontrolled retransmission amplified congestion into near-total throughput loss |
| Slow start | Exponential `cwnd` growth (roughly doubling per RTT) used to probe capacity quickly at connection start |
| AIMD | Additive Increase, Multiplicative Decrease — Reno's congestion-avoidance rule, provably fair per Chiu-Jain analysis |
| Sawtooth | The characteristic `cwnd`-over-time shape produced by loss-based AIMD: linear climb, sudden halving |
| Long fat network problem | High-BDP links on which Reno's fixed +1-MSS-per-RTT growth is far too slow to recover after loss |
| CUBIC | Linux's default congestion control; grows `cwnd` as a cubic function of time since last loss, RTT-independent |
| BBR | Google's model-based algorithm; paces sending using measured `BtlBw` and `RTprop` instead of reacting to loss |
| ECN | Explicit Congestion Notification — router marks (rather than drops) a packet; echoed via TCP's ECE/CWR flags |
| Bufferbloat | Excess queueing delay caused by loss-based senders deliberately filling large router buffers before backing off |

Congestion control decides *how much* data a sender is willing to risk putting on the wire. It says nothing about *what to do the instant something goes wrong* faster than a full timeout would allow — that faster, more surgical reaction is the subject of Chapter 63: fast retransmit and fast recovery.
