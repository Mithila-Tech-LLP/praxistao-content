# Chapter 45: Forwarding, Next Hop, Default Route, Longest Prefix Match, and TTL

> **"A routing table is a list of maybes. Longest prefix match is the rule that turns a pile of maybes into exactly one answer, every single time, in the same number of steps whether the table has ten rows or a million."**

---

## Table of Contents

1. [The Problem Chapter 44 Left Open](#1-the-problem-chapter-44-left-open)
2. [A Naive Algorithm, and Why It's Wrong](#2-a-naive-algorithm-and-why-its-wrong)
3. [Longest Prefix Match, Precisely](#3-longest-prefix-match-precisely)
4. [Worked Example: One Packet, Four Candidate Routes](#4-worked-example-one-packet-four-candidate-routes)
5. [Doing It in Binary, Bit by Bit](#5-doing-it-in-binary-bit-by-bit)
6. [The Default Route as the Ultimate Fallback](#6-the-default-route-as-the-ultimate-fallback)
7. [Next Hop — What It Actually Means](#7-next-hop--what-it-actually-means)
8. [Point-to-Point vs. Broadcast Links](#8-point-to-point-vs-broadcast-links)
9. [TTL — Preventing Packets From Living Forever](#9-ttl--preventing-packets-from-living-forever)
10. [The Full Per-Packet Forwarding Algorithm](#10-the-full-per-packet-forwarding-algorithm)
11. [Seeing It for Real: `ip route get`](#11-seeing-it-for-real-ip-route-get)
12. [Code: Implementing Longest Prefix Match](#12-code-implementing-longest-prefix-match)
13. [Where Metric Fits In (and Where It Doesn't)](#13-where-metric-fits-in-and-where-it-doesnt)
14. [What's Simplified Here](#14-whats-simplified-here)
15. [Common Misconceptions](#15-common-misconceptions)
16. [Hands-On Experiment](#16-hands-on-experiment)
17. [Interview Questions & Model Answers](#17-interview-questions--model-answers)
18. [Exercises](#18-exercises)
19. [Summary](#summary)

---

## 1. The Problem Chapter 44 Left Open

Chapter 44 established that a router's routing table can hold multiple entries, and that a single packet's destination address can plausibly fall inside more than one of them at once. The closing example there was `10.0.0.0/8` and `10.1.0.0/16` both technically covering `10.1.2.5`. That wasn't a contrived edge case — it's the *normal* condition of a real routing table. A router deliberately keeps a broad, aggregated route (`10.0.0.0/8`, "everything in this whole region, roughly this direction") alongside a narrower, more specific route (`10.1.0.0/16`, "this smaller sub-region, more precisely this direction") because both pieces of information are useful and correct at the same time — Chapter 39's CIDR aggregation exists precisely so that one router can summarize what a hundred more specific routers know, while still allowing an intermediate router to carry a more precise exception for part of that space.

So the unavoidable question is: **when more than one row in the table matches a packet's destination, which row wins?** This chapter answers that question completely, then finishes the job by covering the two other pieces of the puzzle every forwarded packet touches: what "next hop" concretely means, and why every IP packet carries a self-destruct counter (TTL) that Chapter 54 will later turn into a diagnostic tool.

---

## 2. A Naive Algorithm, and Why It's Wrong

The Go code at the end of Chapter 44 modeled the naive-but-tempting approach: scan the table in whatever order it happens to be stored, and return the *first* row that matches. Concretely, imagine this routing table:

```
Row 1:  10.0.0.0/8      via 192.0.2.1   (a broad, aggregated route)
Row 2:  10.1.0.0/16     via 192.0.2.2   (a more specific exception)
Row 3:  10.1.2.0/24     via 192.0.2.3   (an even more specific exception)
```

A packet destined for `10.1.2.5` matches **all three rows** — `10.1.2.5` is inside `10.0.0.0/8`, inside `10.1.0.0/16`, and inside `10.1.2.0/24` simultaneously, because CIDR ranges nest inside each other by design. "First match wins" would send this packet via `192.0.2.1`, obeying Row 1 purely because of table ordering — an accident of how the table happens to be sorted, not a reflection of which route is actually correct. If Row 3 exists at all, it's because someone (an administrator, or a routing protocol) deliberately carved out a more precise, more current answer for that specific /24 — maybe that subnet moved to a new provider, maybe it needs special treatment. Ignoring it in favor of a broader, older route silently breaks exactly the case the more specific route was created to handle.

"Last match wins" (scan the whole table, keep overwriting your answer) fixes this specific example by accident only if the table happens to be sorted from least to most specific — which is not a property any real routing table is required to have, and isn't how any real table is stored or updated. The naive approaches all share the same flaw: they let *table order* decide the answer, when the answer should be decided by the *actual specificity of the match*.

---

## 3. Longest Prefix Match, Precisely

**Intuitively:** if someone asks you for directions and you know both "it's somewhere in Asia" and "it's at 221B Baker Street," you give the more specific answer, always, no matter which piece of information you happened to think of first. Longest Prefix Match (LPM) is that instinct, formalized: among every table entry whose prefix matches the destination address, pick the one with the **longest prefix length** (the largest number after the `/`).

**The rule, stated completely:**

> Among all routing table entries whose network prefix contains the destination address, select the entry with the greatest prefix length (the most specific match). If there is still a tie (two entries with the same prefix length both matching — which can only happen if they describe genuinely identical networks from different sources), break it using metric and administrative distance (Chapter 46).

Why "longest" and not "shortest," concretely: a longer prefix means fewer host bits are left unconstrained, which means the entry describes a *smaller, more precisely targeted* range of addresses. `/24` fixes 24 bits and leaves only 8 free (256 addresses); `/8` fixes only 8 bits and leaves 24 free (16.7 million addresses). A route that had to be specifically carved out to cover a narrow slice of address space is, by construction, describing something more deliberate and more current than a route that happens to sweep in that slice as a side effect of covering a much larger region.

---

## 4. Worked Example: One Packet, Four Candidate Routes

Take a router with this routing table — deliberately built with overlapping entries to make the mechanism concrete:

| # | Destination Prefix | Prefix Length | Next Hop | Does `10.1.2.5` fall inside? |
|---|---|---|---|---|
| 1 | `0.0.0.0/0` | 0 | `203.0.113.1` (default route) | Yes — matches everything |
| 2 | `10.0.0.0/8` | 8 | `192.0.2.1` | Yes — `10.1.2.5` is in `10.0.0.0`–`10.255.255.255` |
| 3 | `10.1.0.0/16` | 16 | `192.0.2.2` | Yes — `10.1.2.5` is in `10.1.0.0`–`10.1.255.255` |
| 4 | `10.1.2.0/24` | 24 | `192.0.2.3` | Yes — `10.1.2.5` is in `10.1.2.0`–`10.1.2.255` |

All four rows match. The router doesn't pick arbitrarily and doesn't scan for "the first one that works" — it filters for the **longest prefix length among the matches**, which is unambiguous: Row 4, `/24`, is longer than Rows 3, 2, and 1 (`/16`, `/8`, `/0`, respectively). **The packet is forwarded via `192.0.2.3`.**

Now change the destination address to `10.1.9.200` and redo the exercise:

| # | Destination Prefix | Prefix Length | Does `10.1.9.200` fall inside? |
|---|---|---|---|
| 1 | `0.0.0.0/0` | 0 | Yes |
| 2 | `10.0.0.0/8` | 8 | Yes |
| 3 | `10.1.0.0/16` | 16 | Yes — `10.1.9.200` is inside `10.1.0.0/16` |
| 4 | `10.1.2.0/24` | 24 | **No** — `10.1.9.200` is outside `10.1.2.0`–`10.1.2.255` |

This time Row 4 doesn't match at all — `10.1.9.200` isn't in that narrow /24 range. Among the rows that *do* match (1, 2, 3), Row 3 (`/16`) has the longest prefix, so **the packet goes via `192.0.2.2`**. One more: destination `10.9.9.9`.

| # | Destination Prefix | Does `10.9.9.9` fall inside? |
|---|---|---|
| 1 | `0.0.0.0/0` | Yes |
| 2 | `10.0.0.0/8` | Yes — `10.9.9.9` is inside `10.0.0.0/8` |
| 3 | `10.1.0.0/16` | No |
| 4 | `10.1.2.0/24` | No |

Only Rows 1 and 2 match; Row 2's `/8` beats Row 1's `/0`, so this packet goes via `192.0.2.1`. And finally, destination `172.16.5.1`, which matches nothing except the default:

| # | Destination Prefix | Does `172.16.5.1` fall inside? |
|---|---|---|
| 1 | `0.0.0.0/0` | Yes — always |
| 2, 3, 4 | (all `10.x`) | No |

Only the default route matches, so the packet falls all the way through to `203.0.113.1` — Section 6 covers exactly why this "catches everything" behavior is by design, not a bug.

---

## 5. Doing It in Binary, Bit by Bit

Chapter 37 showed that subnet masking is really a binary AND operation. LPM is the same operation, just applied to every candidate row and then compared. Take the `10.1.2.5` example from Section 4, and check Row 3 (`10.1.0.0/16`) by hand:

```
Destination:      10.1.2.5      = 00001010.00000001.00000010.00000101
Route prefix:     10.1.0.0/16   = 00001010.00000001.00000000.00000000
Mask (/16):                       11111111.11111111.00000000.00000000

Destination AND Mask:            00001010.00000001.00000000.00000000
                                  = 10.1.0.0

Compare to route prefix:         10.1.0.0   -> EQUAL -> MATCH
```

And Row 4 (`10.1.2.0/24`):

```
Destination:      10.1.2.5      = 00001010.00000001.00000010.00000101
Route prefix:     10.1.2.0/24   = 00001010.00000001.00000010.00000000
Mask (/24):                       11111111.11111111.11111111.00000000

Destination AND Mask:            00001010.00000001.00000010.00000000
                                  = 10.1.2.0

Compare to route prefix:         10.1.2.0   -> EQUAL -> MATCH
```

Both are genuine matches at the bit level — this is exactly why the earlier tables showed both rows matching. The router performs this AND-and-compare operation against **every candidate prefix length**, collects every row that comes back equal, and then — and only then — picks the one whose mask had the most 1-bits (the longest prefix). Section 9 of Chapter 44 explained that in practice this isn't done as a literal one-row-at-a-time scan; a trie walk or a TCAM lookup produces the same answer as this brute-force bit comparison, just far faster.

---

## 6. The Default Route as the Ultimate Fallback

`0.0.0.0/0` is not a special-cased keyword the router treats magically — notice in Section 5's binary logic that a `/0` mask is **all zero bits** (`00000000.00000000.00000000.00000000`). ANDing *any* address with an all-zero mask always produces `0.0.0.0`, which always equals the route's prefix (also `0.0.0.0`). That's the entire mechanism: a `/0` prefix matches literally every possible destination address, because the mask constrains zero bits, leaving the entire 32-bit space "don't care." It falls out of the exact same rule as every other route in this chapter — it isn't a separate special case, it's the trivial, shortest-possible-prefix end of the same algorithm.

Because it has the shortest possible prefix length (0), it will *always* lose to any more specific route that also matches — which is exactly the desired behavior. The default route is only ever selected when nothing more specific matches at all, making it the network's designated **catch-all**: "if you don't have better information, send this toward whoever probably does." On a home router, the default route points at the ISP. On an ISP's edge router, the default route (if it has one at all — many core routers deliberately do not, relying only on explicit routes learned via BGP, Chapter 49) points toward a larger transit provider. The phrase Cisco IOS uses for this — "Gateway of last resort," seen in Chapter 44's example output — is a precise description: it's where a packet goes only after every more specific option has been exhausted.

---

## 7. Next Hop — What It Actually Means

Once LPM has chosen a route, that route's **next hop** field says where the packet goes physically next. It's worth being exact about what this value is and is not:

- The next hop is the **IP address of the next router** on the path — specifically, a router that is *directly, physically reachable* from this router (on the same connected network as one of this router's own interfaces). A next hop is never a distant, multi-hop-away address; if it were, the router receiving the packet next wouldn't know how to reach it at Layer 2 at all.
- For a **connected route** (Chapter 44, Section 7), there is no next hop in the usual sense — the destination network *is* directly attached, so the router just needs to find the specific destination host on that wire (via ARP, Chapter 53) and hand the frame over directly.
- For every other route, the next hop tells the router: "don't try to resolve the final destination's MAC address directly — you can't, it's not on any network you're attached to. Instead, resolve *this next-hop router's* MAC address, and hand the packet to it; it's this router's job to figure out the rest."

This is worth connecting explicitly back to Chapter 44 Section 4's key contrast: the packet's **destination IP address never changes** across this entire journey, but the **destination MAC address changes at every single hop**, recomputed fresh via ARP against whatever the current next hop is. The routing table's next-hop field is the input to that per-hop MAC resolution — it's the answer to "who do I physically hand this frame to right now," which is a completely different question from "where is this packet ultimately going."

---

## 8. Point-to-Point vs. Broadcast Links

There's a subtlety worth flagging here because it resurfaces directly in Chapter 48's discussion of OSPF's Designated Router: not every outgoing interface needs a next-hop *address* at all.

- On a **broadcast-capable link** (an Ethernet segment, e.g., the router's LAN-facing interface, or a link shared with multiple other routers), a routing table entry needs a real next-hop IP address, because there could be more than one device reachable on that wire — the router genuinely needs to know *which* neighbor to send the frame to, and then resolve that neighbor's MAC via ARP.
- On a genuine **point-to-point link** (a dedicated serial line, a `/30` or `/31` link directly between exactly two routers, many WAN/leased-line connections, or the routing table entry `via 203.0.113.1 dev eth0` from Chapter 44's example, where `eth0` might in some configurations connect to exactly one other device), there is, by definition, only one possible device on the other end. Some router configurations allow omitting the next-hop address entirely for such links and specifying only the outgoing interface — "just send it out this wire, there's nobody else it could go to."

Most home and enterprise setups you'll encounter specify both interface and next-hop address regardless, which is always correct and unambiguous; the interface-only shorthand is a convenience some platforms allow specifically for point-to-point topologies where it can't possibly be ambiguous.

---

## 9. TTL — Preventing Packets From Living Forever

**The problem TTL solves:** routing tables can, through misconfiguration or during the brief moments a dynamic routing protocol (Chapters 47–48) is still converging after a change, contain a **routing loop** — Router A thinks the best path to some destination is via Router B, and Router B thinks the best path to that same destination is via Router A. A packet caught in that loop would circulate between A and B forever, consuming bandwidth and CPU on both routers, and never actually being delivered or discarded. Without a hard limit, a single misrouted packet could effectively live forever, and a network under a transient loop condition could saturate its links purely with packets going in circles.

**The mechanism:** every IPv4 packet's header carries an 8-bit **Time To Live (TTL)** field — despite the name, it's not measured in seconds; it's a hop counter. A sending host sets an initial value (commonly 64 on Linux, 128 on Windows, 255 on some network equipment). **Every router that forwards the packet decrements TTL by exactly 1** before sending it onward. If a router decrements TTL and finds the result is 0, it does not forward the packet at all — it drops it immediately, and (this is the important part Chapter 54 builds an entire diagnostic tool on top of) sends an **ICMP Time Exceeded** message back to the original source, reporting exactly where the packet died.

```
Packet leaves host:         TTL = 64
Arrives at Router 1:        TTL = 64 -> decrement -> 63, forward
Arrives at Router 2:        TTL = 63 -> decrement -> 62, forward
Arrives at Router 3:        TTL = 62 -> decrement -> 61, forward
...
Arrives at Router 64:       TTL = 1  -> decrement -> 0, DROP
                             Router 64 sends ICMP Time Exceeded back to original sender
```

If a routing loop existed between, say, Routers 20 and 21, a packet caught in it would still only survive a bounded number of trips around that loop before its TTL hit zero — TTL guarantees that even a badly misrouted packet dies in a bounded number of hops rather than living forever. This is a blunt, simple, effective safety net: it does not *prevent* loops (that's what Chapters 33's Spanning Tree Protocol does at Layer 2, and what correct routing protocol design does at Layer 3), it only guarantees loops can't cause packets to circulate infinitely.

**The preview worth sitting with:** Chapter 54 shows that this "flaw" — a router being forced to announce, via ICMP, exactly where a packet died — is not a flaw at all from a diagnostic point of view. `traceroute` deliberately sends packets with TTL starting at 1 and increasing by 1 each time, forcing each successive router along the path to be the one that hits zero and replies with "I killed this one" — turning a safety mechanism into a complete, hop-by-hop map of the path a packet takes. Nothing about TTL was designed with `traceroute` in mind; `traceroute` is simply the cleverest possible use of a side effect that was already there.

---

## 10. The Full Per-Packet Forwarding Algorithm

Putting everything in this chapter together, here is the complete algorithm a router runs for every single IP packet it forwards:

```
function forward(packet):
    if packet.destination_ip is one of this router's own interface addresses:
        deliver locally, stop here          # it's for the router itself, not a hop-through packet

    packet.ttl = packet.ttl - 1
    if packet.ttl == 0:
        drop packet
        send ICMP Time Exceeded to packet.source_ip
        stop here

    candidates = [ route for route in routing_table
                   if route.network.contains(packet.destination_ip) ]

    if candidates is empty:
        drop packet
        send ICMP Destination Unreachable to packet.source_ip
        stop here

    best_route = candidate in candidates with the longest prefix length
                 (break ties by lowest administrative distance, then lowest metric — Chapter 46)

    if best_route.next_hop is directly connected (a "connected" route):
        destination_mac = ARP_resolve(packet.destination_ip)     # Chapter 53
    else:
        destination_mac = ARP_resolve(best_route.next_hop)       # Chapter 53

    recompute IP header checksum (TTL changed, so the checksum must too)
    build new Layer 2 frame with destination_mac, send out best_route.interface
```

Every clause in this pseudocode maps directly onto a section above: the empty-candidates branch is what happens when even the default route is missing (rare, but possible on a router deliberately configured without one); the TTL branch is Section 9; the "longest prefix length" selection is Sections 3–5; and the ARP resolution step is the next-hop mechanic from Section 7, previewing Chapter 53 in full. Note also that the IP header's checksum (a simple error-detection field, in the spirit of Chapter 19) covers only the header, and TTL is part of that header — so decrementing TTL forces every single router along the path to recompute that checksum, one more small but mandatory piece of per-hop work.

---

## 11. Seeing It for Real: `ip route get`

Rather than just trusting the pseudocode, you can ask a real Linux routing table to run its own LPM lookup and show you the winning route directly:

```
$ ip route show
10.1.0.0/16 dev eth1 proto kernel scope link src 10.1.0.1
10.1.2.0/24 dev eth2 proto kernel scope link src 10.1.2.1
default via 203.0.113.1 dev eth0

$ ip route get 10.1.2.5
10.1.2.5 dev eth2 src 10.1.2.1
    cache

$ ip route get 10.1.9.200
10.1.9.200 dev eth1 src 10.1.0.1
    cache

$ ip route get 172.16.5.1
172.16.5.1 via 203.0.113.1 dev eth0 src 203.0.113.2
    cache
```

This mirrors Section 4's worked table exactly: `10.1.2.5` matches both `10.1.0.0/16` and `10.1.2.0/24`, and the kernel picks the longer prefix (`eth2`, the `/24`). `10.1.9.200` only matches the `/16`, so it goes out `eth1`. `172.16.5.1` matches nothing but the default, so it falls through to `eth0` via the ISP. `ip route get` is not a simulation — it invokes the exact same kernel lookup logic the machine uses for real traffic, which makes it one of the most useful commands for debugging "why is this packet going the way it's going" on any Linux box acting as a router.

---

## 12. Code: Implementing Longest Prefix Match

Chapter 44 ended with a deliberately broken `bestMatch` function that returned the first match in table order. Here is the fix — still a linear scan for clarity (a production implementation would use the trie from Chapter 44 Section 9 for speed at scale), but now correctly implementing LPM:

```go
package main

import (
	"fmt"
	"net"
)

type Route struct {
	Network *net.IPNet
	NextHop net.IP
	Iface   string
}

type RoutingTable []Route

// bestMatch now implements real longest-prefix-match: among every
// route whose network contains dst, return the one with the longest
// prefix (the smallest, most specific network).
func (t RoutingTable) bestMatch(dst net.IP) (Route, bool) {
	var best Route
	found := false
	bestLen := -1

	for _, r := range t {
		if !r.Network.Contains(dst) {
			continue
		}
		ones, _ := r.Network.Mask.Size()
		if ones > bestLen {
			best = r
			bestLen = ones
			found = true
		}
	}
	return best, found
}

func main() {
	_, netDefault, _ := net.ParseCIDR("0.0.0.0/0")
	_, net8, _ := net.ParseCIDR("10.0.0.0/8")
	_, net16, _ := net.ParseCIDR("10.1.0.0/16")
	_, net24, _ := net.ParseCIDR("10.1.2.0/24")

	table := RoutingTable{
		{Network: netDefault, NextHop: net.ParseIP("203.0.113.1"), Iface: "eth0"},
		{Network: net8, NextHop: net.ParseIP("192.0.2.1"), Iface: "eth1"},
		{Network: net16, NextHop: net.ParseIP("192.0.2.2"), Iface: "eth2"},
		{Network: net24, NextHop: net.ParseIP("192.0.2.3"), Iface: "eth3"},
	}

	for _, dst := range []string{"10.1.2.5", "10.1.9.200", "10.9.9.9", "172.16.5.1"} {
		ip := net.ParseIP(dst)
		route, _ := table.bestMatch(ip)
		fmt.Printf("%-12s -> next hop %-14s via %s\n", dst, route.NextHop, route.Iface)
	}
}

// Output:
// 10.1.2.5     -> next hop 192.0.2.3      via eth3
// 10.1.9.200   -> next hop 192.0.2.2      via eth2
// 10.9.9.9     -> next hop 192.0.2.1      via eth1
// 172.16.5.1   -> next hop 203.0.113.1    via eth0
```

This output reproduces every result from Section 4's worked-by-hand tables — which is exactly the point: LPM is a completely mechanical, deterministic algorithm, not a judgment call, and the same rule applied by hand, by `ip route get`'s kernel code, and by 30 lines of Go all agree.

---

## 13. Where Metric Fits In (and Where It Doesn't)

A common and understandable confusion: "doesn't the router pick the route with the best/lowest metric?" **Only among routes that already tied on prefix length.** Metric (hop count in RIP, cost in OSPF — Chapters 47–48) and administrative distance (Chapter 46) are tie-breakers used only when longest-prefix-match alone doesn't produce a unique winner — which happens only when two or more routes describe the *exact same* prefix (same network, same mask) but were learned from different neighbors or different protocols.

To make this concrete: if a router has both `10.1.0.0/16 via OSPF, cost 10` and `10.1.0.0/16 via RIP, hop count 3` — same prefix, same length, two different sources — administrative distance decides between them (Chapter 46 covers this ranking; OSPF's default administrative distance of 110 beats RIP's 120, so the OSPF route wins regardless of the individual protocol metrics). But that router will *never* compare a `/16` route's metric against a `/24` route's metric to decide which wins — the `/24`, if it matches, wins purely on specificity, full stop, no matter how "expensive" its metric says it is compared to a broader alternative. Prefix length is checked first and is decisive by itself whenever it differs; metric is a tie-breaker of last resort, not a competing criterion.

---

## 14. What's Simplified Here

- Real hardware forwarding paths (Chapter 44 Section 9's TCAM/trie discussion) don't literally iterate every table row the way the pseudocode and the Go example do — they use structures built specifically so the *result* is identical to this brute-force description, but the *mechanism* is a single fast lookup, not a loop.
- This chapter assumed IPv4 throughout for continuity with earlier volumes; IPv6 uses an identical longest-prefix-match principle over 128-bit addresses, with an equivalent Hop Limit field standing in for TTL (also decremented per hop, also triggering an ICMPv6 Time Exceeded at zero).
- Equal-Cost Multi-Path (ECMP) routing, common on real ISP and data-center routers, allows a router to install *multiple* routes of equal prefix length, equal administrative distance, and equal metric simultaneously, load-balancing traffic across them (often using a hash of the packet's source/destination to keep one flow consistently on one path) — this chapter's algorithm assumed a single best route is always selected for simplicity.
- The pseudocode's ARP resolution step assumes a cache miss requires a fresh lookup; in practice, this is almost always served from a warm ARP cache (Chapter 53) and adds negligible latency.

---

## 15. Common Misconceptions

- **"The router picks the shortest matching route because it's the 'simplest' path."** Backwards — it picks the *longest* matching prefix, because longer prefixes are more specific and therefore more likely to reflect a deliberate, current, correct exception.
- **"TTL is a countdown in seconds, like a packet's expiration date."** No — it's a hop counter, decremented once per router traversed, with no relationship to wall-clock time at all. A packet with TTL 64 that crosses 64 routers dies on the 64th hop regardless of whether that trip took 2 milliseconds or 2 seconds.
- **"A route with a lower metric always wins, even over a route with a shorter prefix."** No — prefix length is checked first and decides the outcome whenever the prefixes differ. Metric only matters when two routes tie on prefix length (Section 13).
- **"The default route is a special keyword the router hardcodes logic around."** It's an ordinary route with prefix length 0, correctly handled by the exact same LPM algorithm as every other row — it simply always loses when anything more specific also matches, and always wins when nothing else does.
- **"A dropped packet due to TTL expiry just silently disappears."** The router that drops it is required to notify the original sender via an ICMP Time Exceeded message — silence is not the intended behavior, though on the real Internet some routers or firewalls do suppress or rate-limit these messages, which is precisely why some `traceroute` output shows gaps (Chapter 54).

---

## 16. Hands-On Experiment

If you have access to a Linux machine (or a VM), try this directly:

```
# Add two deliberately overlapping static routes and watch LPM pick between them:
$ sudo ip route add 10.1.0.0/16 via 192.0.2.2 dev eth1
$ sudo ip route add 10.1.2.0/24 via 192.0.2.3 dev eth1

$ ip route get 10.1.2.5
# Watch it report the /24 route, not the /16 — even though the /16 was added first.

$ ip route get 10.1.9.1
# Watch it fall back to the /16, since /24 doesn't cover this address.
```

You can also directly observe TTL decrementing using a packet capture (previewed properly in Chapter 119):

```
$ ping -c 1 8.8.8.8    # send one ICMP echo with default TTL 64
$ sudo tcpdump -c 1 -v icmp
# Look for "ttl 64" (or however many hops remain by the time it's captured on your own interface)
```

---

## 17. Interview Questions & Model Answers

**Beginner: What does "longest prefix match" mean, in one sentence?**
When a packet's destination address matches more than one entry in the routing table, the router selects the entry with the most specific (longest) prefix — the one that describes the smallest, most precisely targeted range of addresses.

**Intermediate: Why does the default route (`0.0.0.0/0`) always lose when any other route also matches?**
Because it has a prefix length of 0 — the shortest possible — which mathematically means every other prefix length (1 through 32) is equal to or longer, so any route that also matches is by definition at least as specific and will be chosen over it under LPM.

**Intermediate: What actually happens when a packet's TTL reaches 0?**
The router about to forward it decrements TTL, finds the result is 0, and instead of forwarding, drops the packet and sends an ICMP Time Exceeded message back to the original source IP address, identifying itself as the point of failure.

**Advanced: A network engineer says "I added a more specific static route, but traffic is still going the old way." Given everything in this chapter, list two possible causes unrelated to a typo in the prefix itself.**
(1) The new route might have a worse administrative distance than an existing route to the *exact same* prefix from another source, so it isn't the one actually installed as active — but this only matters if the prefixes are identical in length; a genuinely more specific prefix always wins on length alone, so if the new route is truly more specific, distance/metric shouldn't be the cause. (2) The route might not have actually been applied/committed on the router in question, or might exist in the RIB but failed to be pushed into the FIB the data plane actually consults (Chapter 44's RIB/FIB distinction) — this is a common real-world gap between "the control plane thinks it has a route" and "the forwarding hardware is actually using it."

---

## 18. Exercises

### Easy
1. Given routes `192.168.0.0/16` and `192.168.5.0/24`, which one is selected for a packet destined to `192.168.5.10`? Why?
2. Explain in your own words why TTL is measured in hops, not seconds.
3. What ICMP message is sent, and by whom, when a packet's TTL reaches zero?

### Medium
4. Build a table like Section 4's for a router with routes `172.16.0.0/12`, `172.20.0.0/16`, `172.20.4.0/22`, and a default route. Work out, showing your reasoning, which route wins for destinations `172.20.4.5`, `172.20.9.9`, and `172.30.1.1`.
5. Explain why a router needs to recompute the IP header checksum every time it forwards a packet, tying your answer back to which field changes on every hop.
6. A route's next hop is `192.0.2.2`. Explain exactly what the router must do before it can actually transmit the frame — name the protocol involved and which earlier chapter covers it.

### Hard
7. Implement (in Go or Python) a longest-prefix-match function from scratch, without looking at Section 12's code, that takes a list of (prefix, next-hop) pairs and a destination address, and returns the correct next hop. Test it against all four destinations from Section 4.
8. A routing loop briefly exists between two misconfigured routers, each pointing at the other for a particular /24. A packet with TTL 10 enters this loop. Trace exactly what happens to this packet, hop by hop, and explain precisely which router sends the ICMP message and why the loop, while wasteful, cannot run forever.
9. Explain, precisely, why administrative distance and metric can never override a genuinely longer-prefix match — i.e., why a router will never prefer a `/16` OSPF route over a `/24` static route to a destination the `/24` also covers, no matter how favorable the OSPF metric looks.

---

## Summary

| Term | Meaning |
|---|---|
| Longest Prefix Match (LPM) | The rule that, among all matching routing table entries, selects the one with the longest (most specific) prefix |
| Next hop | The IP address of the next directly-reachable router a packet should be handed to |
| Default route (`0.0.0.0/0`) | A route with prefix length 0 that matches every address; used only when nothing more specific matches |
| TTL (Time To Live) | An 8-bit hop counter in the IP header, decremented by 1 at every router; packet is dropped when it reaches 0 |
| ICMP Time Exceeded | The message a router sends back to a packet's source when it drops that packet due to TTL reaching 0 |
| Administrative distance / metric | Tie-breakers used only when two or more routes share the exact same prefix length |
| Point-to-point vs. broadcast link | Link types affecting whether a next-hop address is strictly required for an outgoing route |
| ECMP | Equal-Cost Multi-Path — installing and load-balancing across multiple routes that tie on every criterion |

Chapter 45 completed the exact mechanics of a single forwarding decision. Chapter 46 zooms out to ask how the routing table itself gets built and kept correct in the first place — starting with the simplest possible answer, a human typing in every route by hand, and why that stops working almost immediately at any real scale.
