# Chapter 47: RIP — Distance Vector Routing

> **"Tell your neighbors what you know. Believe what they tell you. Add one. Repeat forever." That's the entire algorithm behind the Internet's first automatic routing protocol — and, as this chapter's worked failure case shows, almost the entire reason it eventually needed to be replaced.**

---

## Table of Contents

1. [The Problem: Automating What Chapter 46 Did by Hand](#1-the-problem-automating-what-chapter-46-did-by-hand)
2. [The Distance-Vector Idea](#2-the-distance-vector-idea)
3. [RIP Mechanics](#3-rip-mechanics)
4. [RIPv1 vs. RIPv2](#4-ripv1-vs-ripv2)
5. [A Worked Convergence Example](#5-a-worked-convergence-example)
6. [The Count-to-Infinity Problem](#6-the-count-to-infinity-problem)
7. [Split Horizon, and Split Horizon With Poison Reverse](#7-split-horizon-and-split-horizon-with-poison-reverse)
8. [Other Mitigations: Hold-Down and Triggered Updates](#8-other-mitigations-hold-down-and-triggered-updates)
9. [RIP's Fundamental Limits](#9-rips-fundamental-limits)
10. [Code: A Minimal RIP-Style Router](#10-code-a-minimal-rip-style-router)
11. [A Real RIP Packet](#11-a-real-rip-packet)
12. [What's Simplified Here](#12-whats-simplified-here)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Hands-On Experiment](#14-hands-on-experiment)
15. [Interview Questions & Model Answers](#15-interview-questions--model-answers)
16. [Exercises](#16-exercises)
17. [Summary](#summary)

---

## 1. The Problem: Automating What Chapter 46 Did by Hand

Chapter 46 ended with static routing's core weakness fully exposed: someone has to know the correct path to every network, type it in, and manually fix it when the topology changes. The obvious next question, stated as plainly as possible: **could routers just tell each other what they know, and work out the routes themselves?**

RIP — the Routing Information Protocol — is the oldest serious answer to that question still in real (if now niche) use, standardized in RFC 1058 (1988), with roots going back to the Xerox PARC Universal Protocol's routing information protocol in the early 1980s. It is deliberately, almost aggressively simple — simple enough that this entire chapter can describe its complete algorithm in a few paragraphs, and simple enough that its failure modes are also easy to demonstrate by hand, which is exactly why it makes the best possible teaching vehicle for *why* dynamic routing protocols need to be more careful than "just believe your neighbors."

---

## 2. The Distance-Vector Idea

**Intuitively:** imagine you're new to a large office building and you don't know which floor any department is on. Instead of walking the whole building yourself, you ask everyone standing near you: "how many doors away is Accounting, from where you're standing?" Someone says "3 doors." You now know Accounting is 4 doors away from you (their 3, plus the 1 door between you and them) — without ever having seen Accounting yourself, and without needing the building's actual floor plan. You then tell *your* neighbors the same thing when they ask you. That's the entire idea: propagate secondhand distance estimates, one hop at a time, and let them accumulate into the right answer everywhere, without anyone needing a full picture of the whole building.

**The real-world analogy, and where it breaks:** this is exactly the childhood game of passing a rumor down a line of people, each one repeating (and slightly modifying — adding a step) what they heard. It works remarkably well *as long as everyone is telling the truth and up to date*. Where the analogy breaks — and where Section 6 hits hard — is what happens when the *original source* of the rumor becomes wrong (say, Accounting moved out entirely) but the rumor keeps circulating anyway, each person still trusting what they last heard from a neighbor who was, in turn, trusting someone else, with nobody able to tell that the root fact has changed. Chapter 46's Section 7 named this family "distance-vector": each router's update is a **vector** — a list of (destination, distance) pairs — describing only distances, never the actual path or topology behind those distances.

**Engineering definition:** in RIP, every router periodically sends its *entire* routing table (destination network + hop count) to each directly connected neighbor. A receiving router compares each advertised destination against its own table: if the neighbor's route, plus one hop for the link to that neighbor, is better than (or new compared to) what the router already has, it adopts that route, records the neighbor as the next hop, and will re-advertise this newly learned, incremented distance to its *own* neighbors on the next update cycle. No router ever learns the actual topology — only ever-accumulating hop counts.

---

## 3. RIP Mechanics

The concrete, standardized details:

- **Metric: hop count.** Every link, regardless of its actual bandwidth, latency, or reliability, costs exactly 1 hop. A route through three routers over gigabit fiber and a route through three routers over a saturated 1 Mbps satellite link look *identical* to RIP — both cost 3. This is RIP's single biggest technical limitation, flagged here and expanded in Section 9.
- **Maximum metric: 15.** Any destination 16 or more hops away is considered **unreachable** — 16 is defined as "infinity" in RIP. This deliberately tiny ceiling was a design choice to bound how long the count-to-infinity failure (Section 6) can possibly run before it's forced to stop.
- **Transport: UDP, port 520.** RIP messages ride directly on UDP (Chapter 58) — no TCP connection, no handshake, no per-neighbor session state, in keeping with the protocol's overall minimalism.
- **Update timing: every 30 seconds**, a router sends its complete routing table to every RIP-enabled neighbor, whether or not anything has changed. This periodic, full-table broadcast is itself a meaningful cost (constant background chatter) that link-state protocols (Chapter 48) specifically avoid.
- **Route timeout: 180 seconds.** If a router doesn't hear a refreshing advertisement for a given route within 180 seconds (six missed periodic updates), it marks that route as invalid (metric set to 16/infinity).
- **Garbage collection: 120 seconds** after being marked invalid, an unreachable route is finally removed from the table entirely (during this window, it's still advertised to neighbors as unreachable — this is what "poisoning" a route means, revisited in Section 7).

These timers alone explain why RIP is often described as "slow": a link failure can take up to 180 seconds just to be locally detected via timeout (faster if the router notices the physical interface itself going down, which most implementations do check directly), and then a further multiple of the 30-second update interval to fully propagate through the rest of the network, hop by hop.

---

## 4. RIPv1 vs. RIPv2

RIP exists in two versions still referenced in practice, and the difference matters directly to earlier chapters in this course:

| | RIPv1 (RFC 1058, 1988) | RIPv2 (RFC 2453, 1998) |
|---|---|---|
| Addressing | Classful only — no subnet mask carried in updates | Carries an explicit subnet mask per route — supports CIDR/VLSM (Chapter 39) |
| Delivery | Broadcast (`255.255.255.255`) | Multicast to `224.0.0.9` — more efficient, doesn't interrupt non-RIP hosts |
| Authentication | None | Simple plaintext or MD5 authentication supported |
| Compatibility with modern networks | Effectively obsolete — cannot represent a VLSM subnet correctly | Still occasionally seen in legacy/lab environments |

RIPv1's inability to carry a subnet mask is a direct, practical consequence of Chapter 39: without a mask in the advertisement, a receiving router has to *guess* the prefix length from the address's old classful A/B/C boundary, which breaks immediately in any network using CIDR-based subnetting more granular than a classful boundary. RIPv2 fixed exactly this gap, which is why it's the version actually referenced in any modern discussion of RIP, and the version this chapter's examples assume from here on.

---

## 5. A Worked Convergence Example

Take four routers in a line — the simplest topology that still shows real multi-hop propagation:

```
   Net A            Net B            Net C            Net D
 10.0.1.0/24     (link R1-R2)     (link R2-R3)     (link R3-R4)     10.0.4.0/24
       \               |                |                |               /
      [R1]============[R2]============[R3]============[R4]
```

Each router starts knowing only its own directly connected networks, with hop count 0 (a router considers its own attached network to be "0 hops away" — reachable without crossing any router at all). Assume all four routers power on simultaneously and begin sending RIP updates every 30 seconds. Track only the routers' knowledge of Net A (`10.0.1.0/24`) and Net D (`10.0.4.0/24`), the two endpoints, to see propagation across the full width of the network:

**Time 0 (before any updates exchanged):**

| Router | Route to Net A | Route to Net D |
|---|---|---|
| R1 | 0 hops, connected | unknown |
| R2 | unknown | unknown |
| R3 | unknown | unknown |
| R4 | unknown | 0 hops, connected |

**Time 30s (first round of updates exchanged with immediate neighbors):**

R1 advertises "Net A: 0 hops" to R2. R2 adopts it as "Net A: 1 hop, via R1." Simultaneously R4 advertises "Net D: 0 hops" to R3, who adopts "Net D: 1 hop, via R4." R2 and R3 also exchange (nothing new yet about A or D between them at this point in a synchronized round, but in practice updates aren't perfectly synchronized — this table assumes idealized lockstep rounds for clarity).

| Router | Route to Net A | Route to Net D |
|---|---|---|
| R1 | 0 hops, connected | unknown |
| R2 | 1 hop, via R1 | unknown |
| R3 | unknown | 1 hop, via R4 |
| R4 | unknown | 0 hops, connected |

**Time 60s (second round):**

R2 now advertises "Net A: 1 hop" to R3, who adopts "Net A: 2 hops, via R2." Symmetrically, R3 advertises "Net D: 1 hop" to R2, who adopts "Net D: 2 hops, via R3."

| Router | Route to Net A | Route to Net D |
|---|---|---|
| R1 | 0 hops, connected | unknown |
| R2 | 1 hop, via R1 | 2 hops, via R3 |
| R3 | 2 hops, via R2 | 1 hop, via R4 |
| R4 | unknown | 0 hops, connected |

**Time 90s (third round):**

R3 advertises "Net A: 2 hops" to R4, who adopts "Net A: 3 hops, via R3." Symmetrically, R2 advertises "Net D: 2 hops" to R1, who adopts "Net D: 3 hops, via R2."

| Router | Route to Net A | Route to Net D |
|---|---|---|
| R1 | 0 hops, connected | 3 hops, via R2 |
| R2 | 1 hop, via R1 | 2 hops, via R3 |
| R3 | 2 hops, via R2 | 1 hop, via R4 |
| R4 | 3 hops, via R3 | 0 hops, connected |

At **T=90 seconds**, every router now has a correct, fully converged route to every network in this small topology. Notice the pattern that generalizes: **a distance-vector network takes roughly (number of hops across the network) × (update interval) to fully converge from a cold start** — here, 3 hops across, 30-second intervals, 90 seconds total. For a larger network — say, 15 routers in a line, RIP's own maximum useful diameter given its 15-hop limit — a cold start could take up to 450 seconds (7.5 minutes) to fully converge, which is a very long time by the standards of an actual outage. This is the exact convergence-speed contrast Chapter 46, Section 7 previewed against link-state protocols.

---

## 6. The Count-to-Infinity Problem

Now take the fully converged network from Section 5 and break something: the link between R1 and Net A goes down (R1's own directly connected network fails — a NIC dies, a cable is cut on R1's LAN side). Trace what happens **without** any of Section 7's mitigations in place, to see the failure mode in full.

**Time 0 (failure occurs):** R1 detects Net A is gone (its own directly connected network — R1 notices this immediately, not via timeout, since it's R1's own interface). R1 removes its "0 hops, connected" route to Net A.

**Time 0, immediately after:** Here is the crux of the problem. R1 still remembers, from its last received update, that **R2 also claims a route to Net A** — R2 had advertised "Net A: 1 hop" back to R1 as part of R2's normal full-table updates (a distance-vector router advertises its entire table to *all* neighbors, including the one it originally learned a route from — this is the specific behavior Section 7's fix targets). R1, not knowing that R2's route to Net A was only ever valid *because of R1 in the first place*, believes this secondhand information: "R2 says it can reach Net A in 1 hop, so I can reach Net A in 2 hops, via R2."

| Router | Route to Net A (before break) | Route to Net A (right after break, naive RIP) |
|---|---|---|
| R1 | 0 hops, connected | **2 hops, via R2** (wrongly believing R2's stale advertisement) |
| R2 | 1 hop, via R1 | still 1 hop, via R1 (hasn't heard anything yet) |
| R3 | 2 hops, via R2 | still 2 hops, via R2 |
| R4 | 3 hops, via R3 | still 3 hops, via R3 |

**Time 30s:** R1 advertises its new (wrong) belief, "Net A: 2 hops," to R2. R2, which had been trusting R1 all along, now updates: "if R1 says 2 hops, and R1 is my neighbor, then I now believe Net A is 3 hops via R1."

| Router | Route to Net A |
|---|---|
| R1 | 2 hops, via R2 |
| R2 | **3 hops, via R1** |
| R3 | 2 hops, via R2 (unchanged so far) |
| R4 | 3 hops, via R3 (unchanged so far) |

**Time 60s:** R1 hears R2's new advertisement ("Net A: 3 hops") and updates again: "3 hops via R2, plus 1, equals 4 hops via R2." Meanwhile R3, still trusting R2's *old* advertisement from before R2 updated (timing offsets make this messier in reality, but the structural point holds regardless of exact timing), or upon receiving R2's new "3 hops" value, updates to "4 hops via R2."

| Router | Route to Net A |
|---|---|
| R1 | 4 hops, via R2 |
| R2 | 4 hops, via R1 (or R3, whichever offers a lower number as it also spirals) |
| R3 | 4 hops, via R2 |
| R4 | 5 hops, via R3 |

**And so on.** Every 30-second round, the believed distance to Net A climbs by roughly 1–2 hops, as R1 and R2 (and eventually R3 and R4) keep re-inflating each other's stale beliefs — each router is technically following the distance-vector algorithm correctly (trust your neighbor, add 1), but the *information itself* is a self-reinforcing echo with no connection to reality anymore, because nobody in this exchange has any way to know the original source of the route (R1's connection to Net A) is gone. This continues, round after round, until the metric finally reaches RIP's defined infinity, **16**, at which point (and only at which point) every router finally agrees the route is genuinely dead and removes it.

| Round | R1's belief | R2's belief | R3's belief | R4's belief |
|---|---|---|---|---|
| 0 (break) | dead (0→gone) | 1 | 2 | 3 |
| 1 | 2 | 1 (stale) | 2 | 3 |
| 2 | 2 | 3 | 2 | 3 |
| 3 | 4 | 3 | 4 | 3 |
| 4 | 4 | 5 | 4 | 5 |
| ... | ... | ... | ... | ... |
| ~13-15 | 16 (infinity) | 16 | 16 | 16 |

This is **count-to-infinity**: a routing loop, formed by mutual misinformation, that only terminates because RIP artificially caps its metric at a small number (15/16) specifically so that this exact failure mode can't run forever. It's a real, historically documented flaw — and it's why the number 15 was chosen as small as it practically could be for the networks RIP was originally designed for: the smaller the ceiling, the faster this worst-case scenario burns itself out, at the direct cost of limiting RIP's usable network diameter to 15 hops.

---

## 7. Split Horizon, and Split Horizon With Poison Reverse

The root cause of Section 6's failure, stated precisely: **R2 kept advertising a route to Net A back to R1 — the very neighbor R2 originally learned that route from.** This is obviously useless information at best (R1 doesn't need to be told about a route it taught R2 in the first place) and actively dangerous at worst (exactly the mechanism that let stale information bounce back and forth). **Split horizon** is the fix: a router must never advertise a route back out the same interface it learned that route from.

Reapplying split horizon to Section 6's scenario: R2 learned its route to Net A from R1 (via the R1–R2 link), so split horizon forbids R2 from ever advertising "Net A: 1 hop" (or anything about Net A at all) back to R1 across that same link. When R1's connection to Net A fails, R2 simply never receives any confusing feedback from R1 about a route to Net A that R2 itself is the sole source of — R1's route to Net A eventually just times out (Section 3's 180-second rule) and is correctly declared unreachable, with no artificial inflation in between.

**Split horizon with poison reverse** goes one step further: instead of *silently* omitting the route (plain split horizon), the router actively advertises that route back with a metric of **infinity (16)** — explicitly saying "I know about Net A, but I want you to know I consider it completely unreachable via this path." This is slightly more bandwidth (an update still gets sent) in exchange for slightly faster, more explicit confirmation that a loop-forming path shouldn't be trusted, rather than relying purely on the other end reasoning "I never heard about this destination from them, so I won't count them as a valid source" (which is what plain split horizon relies on implicitly).

With split horizon (either variant) applied to a simple linear topology like Section 5–6's example, the count-to-infinity scenario in that specific two-router echo (R1↔R2) is completely eliminated — the two routers directly facing each other across the broken link simply cannot re-inflate each other's belief, because the one who would provide the false confirmation is forbidden from ever sending it. It's worth being precise, though, about what split horizon does *not* fully solve, which Section 9 returns to: in topologies with a genuine loop involving three or more routers (not just two facing each other across one link), a longer indirect cycle can still form, because split horizon only prevents a route from bouncing straight back the way it came — it says nothing about a route that goes out, around a longer loop, and comes back from a *different* direction. This is why RIP still needed the hard 15-hop ceiling as a backstop even after split horizon was standard practice — split horizon is a major mitigation, not a complete, structural fix.

---

## 8. Other Mitigations: Hold-Down and Triggered Updates

Two further, complementary techniques are standard in real RIP implementations, stacking on top of split horizon:

- **Triggered updates**: rather than waiting for the next scheduled 30-second periodic update to announce a change, a router sends an update *immediately* the moment it detects a route has failed or changed metric. This alone dramatically shrinks the propagation delay worked out in Section 5–6's round-by-round tables — instead of waiting a full 30 seconds per hop, a genuine change can ripple outward at close to wire speed, hop by hop, limited only by processing and transmission time rather than the periodic timer.
- **Hold-down timers**: once a router learns that a route has become unreachable, it refuses to accept *any* new information about that same route (even from a different neighbor, with a seemingly better metric) for a fixed hold-down period (commonly 180 seconds in classic Cisco RIP implementations) — specifically to give the bad news time to propagate fully through the network before any router considers accepting what might well be more stale, out-of-date good news echoing back around a slower path.

The combination of split horizon (with poison reverse), triggered updates, and hold-down timers is what makes RIP a genuinely usable, if dated, protocol in practice rather than a theoretical curiosity that falls over on the first link failure — but as Section 9 makes explicit, all of these are mitigations layered onto a fundamentally limited core idea, not fixes that make distance-vector routing as fast or as robust as link-state routing.

---

## 9. RIP's Fundamental Limits

Even with every mitigation from Sections 7–8 in place, RIP carries structural limitations that no amount of tuning removes, because they follow directly from the distance-vector idea itself (Section 2):

- **Hop count ignores real cost.** A 1 Gbps fiber link and a 56 kbps dial-up link both cost exactly 1 hop. RIP will happily route traffic across three saturated slow links in preference to one single fast link crossing four hops, because it is structurally blind to anything except hop count. Chapter 48's OSPF, by contrast, uses a cost metric explicitly derived from link bandwidth.
- **The 15-hop diameter limit.** Any network requiring a path of 16 or more router-hops between any two points simply cannot be correctly routed by RIP at all — this rules out RIP for any reasonably large or geographically wide network outright.
- **Slow convergence, even with triggered updates.** Convergence in a distance-vector network fundamentally depends on information propagating hop-by-hop, each hop requiring at least one round of communication and local recomputation — Section 5 showed this scaling with the network's diameter, which link-state's flood-and-recompute approach (Chapter 48) avoids entirely.
- **No real topology awareness.** Because a RIP router never learns the actual shape of the network — only accumulated distances — it cannot make sophisticated routing decisions (like "avoid this specific link because I know it's the network's only path to three other regions") the way a router with a full topology map can.

These are exactly the limitations that motivated OSPF and IS-IS, covered fully in Chapter 48 — not as an arbitrary "better protocol came along," but as a direct, traceable engineering response to RIP's specific, well-understood weaknesses.

---

## 10. Code: A Minimal RIP-Style Router

A simplified simulation of one router's core RIP update-processing logic — receiving a neighbor's advertised table and updating its own, exactly following Section 2's rule (and, notably, *without* split horizon, so you can see the raw mechanism Section 6 exploited):

```go
package main

import "fmt"

const Infinity = 16

type Route struct {
	Metric  int
	NextHop string
}

type Router struct {
	Name  string
	Table map[string]Route // destination -> Route
}

// processUpdate models receiving a neighbor's full routing table
// advertisement and applying the distance-vector rule: candidate
// metric = neighbor's advertised metric + 1 (the cost of this link).
func (r *Router) processUpdate(neighborName string, neighborTable map[string]int) {
	for dest, neighborMetric := range neighborTable {
		candidate := neighborMetric + 1
		if candidate > Infinity {
			candidate = Infinity
		}
		existing, known := r.Table[dest]
		if !known || candidate < existing.Metric || existing.NextHop == neighborName {
			// Adopt if it's new, strictly better, or if this is simply a
			// refreshed update from the neighbor we already route through
			// (this second condition is exactly what lets Section 6's
			// stale-echo problem happen when split horizon isn't applied).
			r.Table[dest] = Route{Metric: candidate, NextHop: neighborName}
		}
	}
}

func main() {
	r2 := &Router{Name: "R2", Table: map[string]Route{
		"NetA": {Metric: 1, NextHop: "R1"},
	}}

	// R1's connection to NetA just failed -- R1 now (wrongly, without
	// split horizon) believes it can reach NetA via R2's stale advert.
	r1Advert := map[string]int{"NetA": 2} // R1 heard "1" from R2 last round, believes 2
	r2.processUpdate("R1", r1Advert)

	fmt.Println(r2.Table["NetA"])
	// Metric climbs to 3, NextHop becomes R1 -- exactly the count-to-infinity
	// echo from Section 6, reproduced by 20 lines of unprotected distance-vector logic.
}
```

Adding split horizon to this code is precisely Exercise 7 at the end of this chapter — it requires nothing more than tracking which neighbor a route was originally learned from and refusing to include that destination in the table sent back to that same neighbor.

---

## 11. A Real RIP Packet

RIPv2's on-the-wire format (RFC 2453), sent inside a UDP datagram to port 520:

| Field | Size | Meaning |
|---|---|---|
| Command | 1 byte | 1 = Request, 2 = Response |
| Version | 1 byte | 2 for RIPv2 |
| Unused | 2 bytes | Zero |
| Address Family Identifier | 2 bytes | 2 for IP |
| Route Tag | 2 bytes | Used to distinguish internal vs. external (redistributed) routes |
| IP Address | 4 bytes | The advertised destination network |
| Subnet Mask | 4 bytes | The mask for that network (RIPv2's key addition over RIPv1) |
| Next Hop | 4 bytes | Optional override of the advertising router as next hop |
| Metric | 4 bytes | Hop count, 1–16 |

A single RIP message can carry up to 25 such route entries (a limit imposed by the maximum practical UDP payload RIP was designed around); a router with a larger table simply sends multiple RIP messages back-to-back. Captured with a packet analyzer (Chapter 119 covers `tcpdump`/Wireshark properly), a RIPv2 response advertising two routes would show, in essence:

```
UDP  src=520 dst=520
RIP  Command: Response (2)
     Version: 2
     Entry 1: 10.0.1.0/24, metric 1
     Entry 2: 10.0.2.0/24, metric 2
```

---

## 12. What's Simplified Here

- Section 5–6's worked examples assumed perfectly synchronized 30-second update rounds across all routers simultaneously; real RIP routers' timers drift independently, so propagation in practice is messier and can be somewhat faster or slower in specific spots than the idealized round-by-round table shown, though the overall order of magnitude and structural failure mode are accurate.
- Split horizon with poison reverse is described as fully solving the two-router echo case, which is correct, but real-world RIP deployments still rely on the 15-hop limit as a hard backstop for indirect loops through three or more routers, as noted in Section 7 — split horizon reduces the frequency and severity of count-to-infinity scenarios, it does not eliminate every possible topology that could trigger one.
- Modern RIP implementations (RIPng for IPv6, and RIPv2 extensions) exist but are rarely deployed in any new network design today; RIP's practical footprint in the mid-2020s is almost entirely legacy equipment, small lab/teaching environments, and a small number of niche or cost-constrained deployments.

---

## 13. Common Misconceptions

- **"RIP picks the fastest path."** RIP picks the path with the *fewest hops*, which is frequently not the fastest path at all — see Section 9's fiber-vs-dial-up example.
- **"Count-to-infinity is a bug in specific RIP implementations."** It's an inherent property of the unprotected distance-vector algorithm itself, not an implementation defect — every distance-vector protocol without mitigations like split horizon faces this exact failure mode, which is exactly why those mitigations are considered mandatory, standard parts of any real deployment.
- **"Split horizon completely eliminates count-to-infinity in all topologies."** It eliminates the simplest, most common two-router case; more complex loops involving three or more routers can still, in principle, trigger a slower version of the same problem, which the 15-hop ceiling exists specifically to bound.
- **"RIP updates only happen when something changes."** By default, full-table updates happen every 30 seconds regardless of whether anything changed — "triggered updates" (Section 8) are an *additional* mechanism for reacting faster to genuine changes, not a replacement for the periodic schedule.

---

## 14. Hands-On Experiment

Most Linux distributions can run a real RIP daemon via FRRouting (`frr`) or the older Quagga suite. If you have access to two or three VMs (or containers) connected in a line:

```
# On each router (via FRR's vtysh CLI):
router rip
 network 10.0.1.0/24
 network 10.12.0.0/30
```

Once configured on each hop, watch real RIP traffic on the wire:

```
$ sudo tcpdump -i eth0 udp port 520 -v
```

You should see periodic broadcasts/multicasts roughly every 30 seconds, each carrying the router's current view of the network. Bring an interface down on the middle router and watch (via `show ip rip` in vtysh, or repeated `tcpdump` captures) the hop counts on the surviving routers climb over successive rounds — a live, real reproduction of Section 6's worked table, assuming split horizon is disabled for the experiment (FRR enables it by default, so intentionally disabling it with `no ip split-horizon` on an interface is the only way to actually observe count-to-infinity happening on real equipment).

---

## 15. Interview Questions & Model Answers

**Beginner: What metric does RIP use, and what is its maximum value?**
RIP uses hop count — the number of routers a packet must cross to reach a destination. The maximum usable value is 15; a metric of 16 is defined as infinity, meaning the destination is unreachable.

**Intermediate: Explain the count-to-infinity problem in your own words.**
When a route becomes unreachable, a distance-vector router may still receive stale advertisements about that same route from a neighbor who originally learned it from the now-failed router itself. Without protection, the two (or more) routers keep incrementing and re-believing each other's increasingly wrong information, causing the metric to slowly climb, round by round, until it reaches infinity (16) and is finally recognized as truly dead.

**Intermediate: What does split horizon do, and what specific scenario does it prevent?**
Split horizon prevents a router from advertising a route back out the same interface it learned that route from. It directly prevents the simplest two-router count-to-infinity echo, where a router would otherwise re-advertise a neighbor's own route right back to that neighbor after the underlying path fails.

**Advanced: Why does OSPF not suffer from count-to-infinity the way RIP does, at a fundamental level, not just "it uses a different algorithm"?**
RIP routers only ever exchange summarized distance opinions and never see the actual topology, so when a route's original source disappears, remaining routers have no way to distinguish "genuinely new, correct information" from "an old belief bouncing back around a longer path." OSPF routers instead flood raw, individually-signed topology facts (link-state advertisements) to every router in the area, and each router independently computes shortest paths from its own complete, first-hand copy of the topology (Chapter 48) — there is no secondhand rumor to become a stale echo, because no router ever depends on trusting another router's *conclusions*, only on trusting the raw facts each router reports about its own directly connected links.

---

## 16. Exercises

### Easy
1. In your own words, explain the difference between "hop count" and "actual link cost" as a routing metric, and give an example where they'd produce different route choices.
2. Why is RIP's maximum metric set to 15 (16 = infinity) rather than some much larger number?
3. What transport protocol and port number does RIP use?

### Medium
4. Redo Section 5's worked convergence example, but for a 5-router line instead of 4, and calculate the total time to full convergence from a cold start using RIP's 30-second update interval.
5. Explain precisely why split horizon requires per-interface (not per-router) bookkeeping — that is, why a router might correctly withhold a route advertisement on one interface while still advertising the same route out a different interface.
6. A network operator disables split horizon "to save a small amount of bandwidth." Explain the specific risk this introduces, referencing Section 6's worked example.

### Hard
7. Extend Section 10's Go code to implement split horizon: track which neighbor each route was learned from, and refuse to include that destination when building the advertisement sent back to that same neighbor. Verify it prevents the specific echo shown in Section 6.
8. Construct a three-router triangle topology (R1–R2, R2–R3, R3–R1, all interconnected) where, even with basic split horizon applied on every direct pair, a slower count-to-infinity-style loop can still form through the indirect third router. Walk through several rounds by hand.
9. Explain why triggered updates alone (without split horizon) would not have prevented Section 6's failure — i.e., why "send updates faster" and "don't create the false information in the first place" are solving two different problems.

---

## Summary

| Term | Meaning |
|---|---|
| Distance-vector routing | Routers share summarized distance opinions with neighbors and accumulate them hop by hop |
| RIP | The Routing Information Protocol — the classic distance-vector IGP, hop-count metric, UDP port 520 |
| Hop count | RIP's metric — number of routers crossed; maximum useful value 15, 16 = infinity/unreachable |
| Convergence | The time for all routers to agree on correct routes after a change; scales with network diameter × update interval for distance-vector protocols |
| Count-to-infinity | The failure mode where stale, mutually-reinforced route information causes metrics to climb slowly toward infinity instead of immediately reflecting a failure |
| Split horizon | Never advertise a route back out the interface it was learned from — the primary count-to-infinity mitigation |
| Poison reverse | Split horizon's stronger variant: actively advertise the withheld route as metric infinity, rather than omitting it silently |
| Triggered update | An immediate advertisement sent upon detecting a change, rather than waiting for the next periodic update |
| Hold-down timer | A period during which a router refuses new information about a route it just learned is unreachable, to let bad news propagate before accepting possibly-stale good news |

RIP proved dynamic routing works, but its blindness to real link cost, its 15-hop ceiling, and its slow, rumor-based convergence left an obvious opening. Chapter 48 covers the protocols built to close it — OSPF and IS-IS — by replacing secondhand rumors with a complete, shared map of the network and a hard mathematical guarantee (Dijkstra's algorithm) about finding the actual shortest path.
