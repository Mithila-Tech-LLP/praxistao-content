# Chapter 50: Autonomous Systems and Route Aggregation

> *"An IP address tells you where a machine is. An AS number tells you who is legally, financially, and operationally responsible for the network it's on."*

---

## Table of Contents

1. [The Missing Piece From Chapter 49](#1-the-missing-piece-from-chapter-49)
2. [What an Autonomous System Actually Is](#2-what-an-autonomous-system-actually-is)
3. [The Naive Alternative: Route Every Individual Network](#3-the-naive-alternative-route-every-individual-network)
4. [AS Numbers — Structure, Ranges, and Who Assigns Them](#4-as-numbers--structure-ranges-and-who-assigns-them)
5. [Types of Autonomous Systems](#5-types-of-autonomous-systems)
6. [Why BGP Is Really "The Protocol That Connects ASes"](#6-why-bgp-is-really-the-protocol-that-connects-ases)
7. [Route Aggregation, Revisited at ISP Scale](#7-route-aggregation-revisited-at-isp-scale)
8. [How Aggregation Actually Happens in BGP](#8-how-aggregation-actually-happens-in-bgp)
9. [Worked Example: One ISP, Thousands of Customers, One Announcement](#9-worked-example-one-isp-thousands-of-customers-one-announcement)
10. [The Cost of Aggregation: What You Lose](#10-the-cost-of-aggregation-what-you-lose)
11. [The Growth of the Global Routing Table](#11-the-growth-of-the-global-routing-table)
12. [Packet/Message View: AGGREGATOR and ATOMIC_AGGREGATE](#12-packetmessage-view-aggregator-and-atomic_aggregate)
13. [A Real Example: Looking Up an AS](#13-a-real-example-looking-up-an-as)
14. [Hands-On Experiment](#14-hands-on-experiment)
15. [Code: Aggregating a List of Prefixes in Go](#15-code-aggregating-a-list-of-prefixes-in-go)
16. [Common Misconceptions](#16-common-misconceptions)
17. [Production Notes](#17-production-notes)
18. [What This Chapter Simplified](#18-what-this-chapter-simplified)
19. [Interview Questions & Model Answers](#19-interview-questions--model-answers)
20. [Exercises](#20-exercises)
21. [Summary](#21-summary)

---

## 1. The Missing Piece From Chapter 49

Chapter 49 used the term "Autonomous System" dozens of times without ever formally defining it, because you needed to see BGP's *behavior* first to appreciate why the definition matters. Now we go back and nail it down, because "Autonomous System" is not a fuzzy, hand-wavy term — it's a precise, globally-registered, numbered entity, and understanding exactly what it is (and isn't) resolves a lot of confusion about how the Internet's routing actually scales.

We'll also finish a thread left open all the way back in Chapter 39: CIDR let a network be described by *any* prefix length instead of rigid class A/B/C blocks, and Chapter 39 mentioned, in passing, that this also enables **route aggregation** — combining many small announcements into one big one. This chapter shows exactly how that works at the scale that keeps the *entire Internet's* routing table from collapsing under its own weight.

---

## 2. What an Autonomous System Actually Is

**The problem, stated plainly:** BGP (Chapter 49) needs some way to answer the question "who is on the other side of this routing decision?" It can't reason in terms of individual routers — a network might have thousands of them — and it can't reason in terms of individual IP prefixes, because a single organization might own hundreds of them, all governed by the exact same routing policy. BGP needs a unit that represents **one coherent zone of routing policy and administrative control.**

That unit is the **Autonomous System (AS)**: a collection of IP networks and routers, under the control of a single administrative entity (a company, university, government, or ISP), that presents a common, clearly-defined routing policy to the rest of the Internet.

Three things make something a genuine "Autonomous System" rather than just "a network":

1. **Single administrative control.** One organization decides how routing works inside it — which internal routing protocol to run (OSPF, IS-IS, or even static routes), how to connect to the outside world, and what policy to apply.
2. **A coherent, external-facing routing policy.** From the outside, the rest of the Internet doesn't need to know or care that AS 65001 is internally running OSPF with three areas and three redundant core routers. It only needs to know: "AS 65001 can reach these prefixes, and here's how to reach AS 65001 itself."
3. **A globally unique number — the ASN — identifying it in BGP.** This is what actually appears in every AS-PATH (Chapter 49, Section 8).

**Intuitive analogy:** think of an AS as a *country* in the context of international shipping. A shipping company doesn't need to know India's internal highway network to route a package to Mumbai — it only needs to know "this package is going to India" and hand it to India's own postal/customs system at the border, trusting that country to handle internal delivery its own way. BGP is the treaty between countries about which shipping lanes exist and who agreed to carry whose packages; OSPF/IS-IS is the internal highway system each country builds and manages entirely on its own. **Where the analogy breaks:** countries have fixed geographic borders enforced by international law and often physical force; an AS's "border" is just wherever its administrators chose to run BGP instead of an IGP, and that boundary can be redrawn (a company splitting into two ASes, or two companies merging their ASes) far more easily than a national border.

---

## 3. The Naive Alternative: Route Every Individual Network

Before accepting "Autonomous System" as the right abstraction, it's worth asking: why not just have BGP-like routers exchange routes for individual *networks* or even individual *routers*, without grouping them into ASes at all?

The immediate problem is the same scaling problem from Chapter 49, Section 2, but sharper: if BGP treated every individual router as its own routing entity, the global routing table wouldn't have ~70,000 entries worth of "who owns what" — it would have to represent every router on Earth, likely in the tens of millions, with no natural grouping to summarize them. Loop detection (which relies on "have I seen this AS number before" in the AS-PATH) would also become useless, since the path would be a list of individual routers rather than administrative domains, defeating the entire elegance of path-vector routing.

The AS abstraction is what makes BGP's approach *hierarchical*: from the outside, one AS number stands in for however many routers, subnets, and internal protocols that organization has chosen to run. This mirrors exactly the reasoning behind CIDR and subnetting (Chapters 37-39): hierarchy is what lets a system scale, because it lets you summarize "everything under here" without needing to enumerate every leaf.

---

## 4. AS Numbers — Structure, Ranges, and Who Assigns Them

An **ASN (Autonomous System Number)** is assigned by one of the five **Regional Internet Registries (RIRs)** — the same organizations responsible for allocating IP address blocks (Chapter 36):

| RIR | Region |
|---|---|
| ARIN | United States, Canada, parts of the Caribbean |
| RIPE NCC | Europe, Middle East, parts of Central Asia |
| APNIC | Asia-Pacific |
| LACNIC | Latin America and the Caribbean |
| AFRINIC | Africa |

**Original 16-bit ASNs** (RFC 1930) ranged from 0 to 65,535 — about 65,000 possible numbers, which felt limitless in the early 1990s and was exhausted within a couple of decades as the number of independently-routed networks exploded. **32-bit ASNs** (RFC 6793, "Four-octet AS Number Space," standardized in 2007 and now near-universally deployed) extended the range up to roughly 4.29 billion, following exactly the same "we ran out of address space, so widen the field" pattern you already saw with IPv4 → IPv6 (Chapter 42).

Certain ranges are reserved and must never appear as an origin AS on the public Internet:

```
AS 0                         : reserved, means "invalid/unspecified" in RPKI (Chapter 52)
AS 23456                     : reserved as an "AS_TRANS" placeholder for backward
                                compatibility during the 16-bit -> 32-bit transition
AS 64512 - AS 65534          : private use (16-bit range) -- exactly analogous to
                                RFC 1918 private IP space (Chapter 40)
AS 65535                     : reserved
AS 4200000000 - 4294967294   : private use (32-bit range)
AS 4294967295                : reserved
```

Private ASNs are legitimately used inside enterprise networks that run BGP internally (for instance, a company multihomed to two different ISPs might use a private ASN on its own edge router) but, exactly like private IP addresses, must be stripped or translated before reaching the public Internet — an ISP would never accept a route with a private ASN anywhere in its AS-PATH from a public peer.

You can find any organization's real ASN with a simple `whois` query:

```bash
$ whois -h whois.radb.net AS15169
aut-num:    AS15169
as-name:    GOOGLE
descr:      Google LLC
country:    US
```

---

## 5. Types of Autonomous Systems

Not every AS plays the same role. Three broad categories, defined by how they connect to the rest of the Internet, matter for understanding routing policy (and set up Chapter 51's peering/transit discussion directly):

- **Stub AS.** Connects to exactly one other AS, and only carries traffic to and from itself — it never carries transit traffic for anyone else. A small business or university with a single ISP connection is a classic stub AS. Its BGP policy is trivial: "send everything not destined for me to my one upstream."
- **Multihomed AS.** Connects to two or more other ASes for redundancy or performance, but — like a stub AS — does not provide transit for third parties; it only originates and receives its own traffic. Many mid-size companies and content providers are multihomed for reliability (if one upstream link fails, traffic shifts to the other).
- **Transit AS.** Carries traffic *between* other Autonomous Systems as its core business — this is what an ISP fundamentally is. A transit AS must run real BGP policy logic, because it's constantly deciding, for traffic that isn't even its own, which of its several upstream and peer connections to forward it through.

```
   Stub AS               Multihomed AS              Transit AS
  (university)          (mid-size company)          (an ISP)

      │                    │        │              ╲   │   ╱
   ┌──┴──┐              ┌──┴──┐  ┌──┴──┐         ┌───┴───┴───┴───┐
   │ AS X│              │ AS Y│  │     │         │     AS Z       │
   └─────┘              └─────┘  (ISP2)          └───────────────┘
      │                    │                       │    │    │
   (ISP1)               (ISP1)                 (customer)(peer)(customer)

   Only sends/           Sends/receives         Actively forwards
   receives own          own traffic via         traffic BETWEEN
   traffic               either upstream          other ASes
```

---

## 6. Why BGP Is Really "The Protocol That Connects ASes"

It's worth stating plainly what falls out of Sections 2-5: **BGP is not, at its core, a protocol for routing IP prefixes. It's a protocol for expressing relationships between Autonomous Systems, which happens to carry IP prefixes as the payload of that relationship.**

Every AS-PATH attribute (Chapter 49, Section 8) is literally a list of AS numbers — not router IDs, not IP addresses, not link costs. Every routing policy an ISP writes ("prefer routes from my peers over my transit provider," "never accept a route whose AS-PATH contains AS 64512") is expressed in terms of AS numbers. The entire "who trusts whom, who pays whom" business layer that Chapter 51 unpacks is fundamentally a graph of **AS-to-AS relationships** — and BGP is the mechanism by which that graph of relationships gets turned into an actual, working, packet-forwarding decision at every router on Earth.

This is why, when network engineers say "BGP connects the Internet," they don't mean it connects individual computers, or even individual routers — they mean it connects **Autonomous Systems**, tens of thousands of independently-run networks, into one routable whole.

---

## 7. Route Aggregation, Revisited at ISP Scale

Chapter 39 introduced route aggregation conceptually: CIDR's variable-length prefixes mean that if you own a contiguous block of addresses split into many smaller subnets, you can advertise just the one larger, containing prefix instead of every smaller piece — *as long as every one of those smaller pieces is actually reachable the same way.*

Chapter 39's examples were modest — an organization aggregating a handful of /24s into a /22. At **ISP scale**, this same idea is what makes the entire global routing system survive at all.

**The problem, concretely:** A large regional ISP might have been allocated a `/16` (65,536 addresses) by its RIR, and it hands out pieces of that block — a `/24` here, a `/26` there — to thousands of individual business and residential customers over the years. If that ISP announced every single one of those customer allocations as its own separate BGP route, a single ISP's presence in the global routing table could balloon into tens of thousands of entries — and there are tens of thousands of ISPs on Earth doing the same thing. Multiply that out and the global BGP table would need to hold **billions** of routes, something no router's memory or CPU (running the best-path algorithm on every table change) could keep up with.

**The naive alternative — announce everything individually — genuinely was closer to how things worked in the early, pre-CIDR Internet**, and it's a big part of *why* CIDR was urgently adopted in 1993: the classful system (Chapter 39) plus a lack of aggregation discipline had the global routing table on a trajectory to outgrow router hardware of the era within a few years.

**The real solution:** an ISP advertises **one summarized prefix** — say, its entire `/16` — to the rest of the Internet, instead of the thousands of more specific customer routes making it up. Every router outside that ISP sees a single entry: "to reach any address in this /16, send it to this ISP." What the ISP does *internally* with that /16 — how it's split into a /24 here and a /26 there among its own customers — is the ISP's own business, invisible to the rest of the world, exactly the way an AS's internal OSPF topology (Section 2) is invisible to the rest of the world.

---

## 8. How Aggregation Actually Happens in BGP

Concretely, a router (usually the ISP's border router, or a purpose-built route reflector) is configured to originate an **aggregate route** — a summarized prefix that the router itself creates and injects into BGP, distinct from any prefix it directly learned from a customer or peer.

Two BGP attributes exist specifically to mark this (introduced briefly in Chapter 49, Section 8):

- **AGGREGATOR**: identifies which AS number (and, often, which specific router) performed the aggregation — useful for debugging and for downstream networks to understand where a summarized route actually originated.
- **ATOMIC_AGGREGATE**: a flag meaning "this route represents an aggregation of multiple more-specific routes, and as a result, some path attribute information from the individual constituent routes (like a fully accurate AS-PATH for every sub-piece) may have been lost in the summarization." It's an honesty flag — a warning that the aggregate is a simplification, not a perfectly faithful summary of every constituent route's exact path.

Most real-world configuration uses a simple declarative statement — for example, in Cisco IOS style:

```
router bgp 65001
 aggregate-address 203.0.113.0 255.255.0.0 summary-only
```

The `summary-only` keyword is important: without it, the router would advertise **both** the aggregate *and* all the individual more-specific routes it's built from — which defeats the purpose. `summary-only` suppresses the individual routes from being advertised externally, sending only the one summarized announcement.

---

## 9. Worked Example: One ISP, Thousands of Customers, One Announcement

Imagine "Example ISP" (AS 65010) has been allocated `203.0.113.0/16`... (using this only as a numerically illustrative example — in reality 203.0.113.0/24 is the actual reserved documentation block per RFC 5737, but we'll treat it here purely as a teaching-scale stand-in for "a /16 an ISP owns," as is conventional in networking pedagogy).

```
Example ISP owns:            203.0.113.0/16   (65,536 addresses)

Internally allocated as:
  203.0.113.0/24    -> Customer A (small business)
  203.0.114.0/24    -> Customer B
  203.0.115.0/26    -> Customer C (residential fiber block)
  ... (thousands more individual customer allocations) ...
  203.0.200.0/23    -> Customer Z (a larger enterprise customer)

Without aggregation:
  Example ISP would need to announce THOUSANDS of individual
  BGP routes to the rest of the Internet -- one per customer
  allocation, each one taking up a slot in every other AS's
  routing table on Earth.

With aggregation:
  Example ISP announces ONE route:

      203.0.113.0/16   AS-PATH: 65010

  Every other AS on the Internet stores exactly this ONE entry
  to reach ALL of Example ISP's customers, regardless of how many
  thousand individual customer subnets exist behind it.
```

The savings compound globally: if every mid-size-or-larger ISP on Earth aggregates its customer allocations this way (and, per Section 17, this is heavily enforced in practice), the global routing table stays in the hundreds of thousands of entries instead of ballooning into the hundreds of millions.

---

## 10. The Cost of Aggregation: What You Lose

Aggregation isn't free — it trades away information, and that trade-off matters in two specific ways:

1. **Loss of path granularity.** If Customer A and Customer Z (from Section 9) are reached via genuinely different physical paths inside Example ISP's network, the rest of the Internet has no way to know that from the outside — it only sees one route to the whole /16 and picks one entry/exit point for all of it. This is invisible and harmless almost all the time, because internal routing (OSPF/IS-IS) handles the actual delivery once traffic arrives inside the ISP; but it does mean the rest of the Internet cannot make finer-grained routing decisions about individual customers within an aggregated block.

2. **The "more-specific-wins" escape hatch — and its abuse.** Longest-prefix match (Chapter 45) means that if *any* AS anywhere announces a **more specific** route than the aggregate — say, a stray `/24` inside Example ISP's aggregated `/16` — every router that sees both routes will prefer the more specific one for addresses inside that /24, *regardless of AS-PATH length, LOCAL-PREF, or any other BGP attribute*, because longest-prefix match happens at the forwarding table lookup stage, entirely separate from and prior to BGP's best-path comparison of routes to the *same* prefix. This is normal and useful — an ISP intentionally announces a more-specific route to steer traffic to a particular customer's backup link, for instance. But it is *also exactly the mechanism that route hijacking abuses* (Chapter 52 covers this at length): if someone else announces a more-specific sub-block of your address space, the Internet will, by design, prefer their announcement over your legitimate aggregate for that sub-block, with no protocol-level check on whether they had any right to do so.

---

## 11. The Growth of the Global Routing Table

The numbers make the stakes concrete. The IPv4 global routing table (as tracked by public BGP monitoring projects like CIDR Report and RIPE RIS) has grown roughly as follows:

```
Year        Approx. IPv4 routing table size
1994          ~20,000 routes    (shortly after CIDR's introduction)
2001          ~100,000 routes
2010          ~340,000 routes
2019          ~780,000 routes
2024          ~950,000+ routes   (plus several hundred thousand IPv6 routes)
```

This growth happened *despite* aggregation discipline, driven mostly by continued Internet growth (more networks, more multihoming, more IPv6 deployment) — but without aggregation, and without CIDR replacing the old classful system in the first place, industry consensus at the time (the early-to-mid 1990s "routing table growth crisis") was that router memory and route-computation capacity of that era would have been overwhelmed well before the year 2000. Aggregation is a big part of why that crisis didn't materialize the way it was feared to.

---

## 12. Packet/Message View: AGGREGATOR and ATOMIC_AGGREGATE

Both attributes introduced in Section 8 travel as ordinary optional path attributes inside a BGP UPDATE message (recall the general UPDATE structure from Chapter 49, Section 12):

```
Path Attribute: AGGREGATOR
  Attribute Type: Optional Transitive, code 7
  Length: 6 or 8 bytes (4-byte ASN + 4-byte IP address of aggregating router)
  Value: <ASN of aggregating router> <Router ID / IP of aggregating router>

Path Attribute: ATOMIC_AGGREGATE
  Attribute Type: Well-known Discretionary, code 6
  Length: 0 bytes  (it's a pure flag -- its mere presence is the signal)
```

Note that `ATOMIC_AGGREGATE` carries zero bytes of value — its entire meaning is conveyed just by whether the attribute is present at all in the UPDATE, a nice minimal example of a flag attribute versus a value-carrying one.

---

## 13. A Real Example: Looking Up an AS

Public tools make it trivial to explore real Autonomous Systems and see aggregation in action:

```bash
# Look up which AS originates a given prefix, and its full route:
$ whois -h whois.radb.net 8.8.8.0/24
route:      8.8.8.0/24
origin:     AS15169
descr:      Google Public DNS block

# See ALL prefixes a given ASN announces (a good way to spot aggregation:
# does the ISP announce one big block, or hundreds of tiny ones?):
$ whois -h whois.radb.net -i origin AS15169

# Use RIPEstat's API to see a visual/history breakdown of announcements
# for a given resource over time, including aggregation and deaggregation
# events:
curl -s "https://stat.ripe.net/data/routing-history/data.json?resource=AS15169" | head -80
```

Browsing `bgp.he.net/AS15169` (Hurricane Electric's BGP toolkit, mentioned already in Chapter 49) in a browser shows Google's actual announced prefixes and immediately makes the aggregation pattern visible: large content/cloud providers typically announce a manageable number of sizable, aggregated blocks rather than a sprawl of tiny ones — good aggregation hygiene is treated as an operational best practice and is watched by tools like the CIDR Report, which specifically calls out ASes announcing excessive numbers of unaggregatable, overly-specific routes.

---

## 14. Hands-On Experiment

```bash
# Compare the number of routes an ISP announces to a rough estimate of
# how many distinct /24s it "should" need, given its allocated space,
# to spot how aggressively (or poorly) it aggregates:

# 1. Get an ASN's announced prefix count:
curl -s "https://stat.ripe.net/data/announced-prefixes/data.json?resource=AS3356" \
  | python3 -c "import json,sys; d=json.load(sys.stdin); print(len(d['data']['prefixes']))"

# 2. Compare the SIZE of the largest few prefixes it announces (a well-
#    aggregating large ISP will show a mix of large /13-/16 blocks;
#    a poorly-aggregating one will show mostly /24s):
curl -s "https://stat.ripe.net/data/announced-prefixes/data.json?resource=AS3356" \
  | python3 -c "
import json, sys
d = json.load(sys.stdin)
prefixes = [p['prefix'] for p in d['data']['prefixes']]
sizes = sorted(set(p.split('/')[1] for p in prefixes if '.' in p))
print('Prefix lengths seen:', sizes)
"
```

Try this against a large, old, well-run Tier-1 network's ASN versus a small regional ISP's ASN, and compare the spread of prefix lengths you see — it's a direct, hands-on view of aggregation discipline in the wild.

---

## 15. Code: Aggregating a List of Prefixes in Go

Here's a minimal but genuinely useful Go program that takes a list of contiguous /24 allocations and finds the smallest set of CIDR blocks that summarizes them — the same core problem a real router's `aggregate-address` configuration solves, simplified to operate on a small in-memory list rather than a live routing table:

```go
package main

import (
	"fmt"
	"net"
	"sort"
)

// prefixToUint32 converts an IPv4 network's base address to a uint32
// for numeric comparison and bit manipulation.
func prefixToUint32(n *net.IPNet) uint32 {
	ip := n.IP.To4()
	return uint32(ip[0])<<24 | uint32(ip[1])<<16 | uint32(ip[2])<<8 | uint32(ip[3])
}

// canMerge checks whether two equally-sized CIDR blocks are adjacent and
// aligned such that they can be combined into one block with a prefix
// length one bit shorter -- the fundamental aggregation test a router
// applies when deciding whether two routes can become one.
func canMerge(a, b *net.IPNet) (*net.IPNet, bool) {
	onesA, bitsA := a.Mask.Size()
	onesB, bitsB := b.Mask.Size()
	if onesA != onesB || bitsA != bitsB || onesA == 0 {
		return nil, false
	}
	addrA := prefixToUint32(a)
	addrB := prefixToUint32(b)
	blockSize := uint32(1) << (32 - onesA)

	// They must be adjacent AND the pair must be aligned on a boundary
	// of the larger (merged) block -- e.g. 203.0.112.0/24 and
	// 203.0.113.0/24 merge into 203.0.112.0/23, but 203.0.113.0/24 and
	// 203.0.114.0/24 do NOT (not aligned to a /23 boundary).
	lower := addrA
	if addrB < lower {
		lower = addrB
	}
	if lower%(blockSize*2) != 0 {
		return nil, false
	}
	if addrB != addrA+blockSize && addrA != addrB+blockSize {
		return nil, false
	}

	mergedIP := net.IPv4(byte(lower>>24), byte(lower>>16), byte(lower>>8), byte(lower))
	mergedMask := net.CIDRMask(onesA-1, bitsA)
	return &net.IPNet{IP: mergedIP.To4(), Mask: mergedMask}, true
}

// aggregate repeatedly merges adjacent, aligned prefixes until no more
// merges are possible -- a simplified version of what a router's
// aggregation logic (Section 8) does continuously as routes come and go.
func aggregate(prefixes []*net.IPNet) []*net.IPNet {
	changed := true
	for changed {
		changed = false
		sort.Slice(prefixes, func(i, j int) bool {
			return prefixToUint32(prefixes[i]) < prefixToUint32(prefixes[j])
		})
		for i := 0; i < len(prefixes)-1; i++ {
			if merged, ok := canMerge(prefixes[i], prefixes[i+1]); ok {
				newList := append([]*net.IPNet{}, prefixes[:i]...)
				newList = append(newList, merged)
				newList = append(newList, prefixes[i+2:]...)
				prefixes = newList
				changed = true
				break
			}
		}
	}
	return prefixes
}

func main() {
	_, a, _ := net.ParseCIDR("203.0.112.0/24")
	_, b, _ := net.ParseCIDR("203.0.113.0/24")
	_, c, _ := net.ParseCIDR("203.0.114.0/24")
	_, d, _ := net.ParseCIDR("203.0.115.0/24")

	result := aggregate([]*net.IPNet{a, b, c, d})
	fmt.Println("Aggregated announcement(s):")
	for _, p := range result {
		fmt.Println(" ", p.String())
	}
	// Output: a single "203.0.112.0/22" -- four /24 customer allocations
	// collapsed into one route, exactly as Section 9 describes at scale.
}
```

---

## 16. Common Misconceptions

- **"An AS is the same thing as an ISP."** Most large ISPs are ASes, but plenty of ASes are not ISPs at all — universities, large enterprises, cloud providers, and content networks all hold their own ASNs without ever selling transit to anyone.
- **"A bigger ASN means a bigger or more important network."** ASN values (especially in the 32-bit range) are just assigned sequentially by RIRs as requested — a large legacy network might hold a low, early 16-bit ASN, while a newer but equally important network holds a high 32-bit one. Size/importance is unrelated to the number's magnitude.
- **"Route aggregation always makes routing worse because it's less specific."** It only becomes a real problem when a more-specific route exists somewhere else in a way you didn't intend (Section 10) — used correctly, aggregation is a pure scalability win with essentially no cost, because internal routing still handles final delivery accurately.
- **"An ISP's aggregate announcement means all its customers are on one physical link."** No — aggregation is purely a *routing table* summarization. It says nothing about the ISP's actual internal topology, which could span an entire continent's worth of physical infrastructure.

---

## 17. Production Notes

- Real-world route aggregation is actively monitored by the community: the **CIDR Report** (cidr-report.org) has, for over two decades, published rankings of which ASes contribute the most "unaggregatable" or excessively specific routes to the global table, applying informal peer pressure toward good aggregation hygiene.
- Most RIRs enforce a **minimum allocation size** they will assign (commonly /24 for IPv4 in most regions today) partly because routes longer than /24 are widely filtered and dropped by many networks' BGP configurations — a /25 or smaller generally won't even propagate globally, which is itself an aggregation-forcing mechanism baked into common operational practice.
- Multihomed customers (Section 5) sometimes *must* announce a more-specific route alongside an upstream's aggregate, specifically to ensure return traffic prefers their intended link — a legitimate, everyday use of the "more-specific wins" rule from Section 10, done deliberately and with the knowledge of both upstreams involved.
- **AS-SETs**, a related BGP concept, let a network's aggregate route also carry a *set* of the ASNs that contributed to it (rather than one clean AS-PATH), used in specific route-registry and filtering contexts — an advanced detail beyond this chapter's scope but worth knowing the term for.
- **IPv6 aggregation follows the identical logic** but starts from a much more generous position: RIRs typically hand out IPv6 allocations in large, deliberately hierarchical blocks (often a /32 or larger to an ISP, who can then sub-allocate /48s to individual customers, Chapter 42) specifically so that clean, aligned aggregation is the natural default rather than something operators have to work to achieve — one of the quieter, less-discussed benefits of IPv6's vastly larger address space is that it makes the "the boundaries just happen to align" problem from Section 10's `canMerge` logic far less common in the first place.
- Real ASN population statistics illustrate Section 5's taxonomy at scale: of the roughly 70,000+ ASNs actively announcing routes on the Internet today, the overwhelming majority are stub or multihomed ASes (enterprises, universities, and smaller networks that only originate their own traffic); only a few thousand ASes ever appear as a **transit** hop — meaning most of the global routing table's actual path diversity is carried by a comparatively small set of networks, a concrete numeric preview of the Tier-1/Tier-2/Tier-3 hierarchy Chapter 51 formalizes next.

---

## 18. What This Chapter Simplified

- The worked examples use `203.0.113.0/16`-style addressing purely for illustrative round numbers; in reality this exact block is RFC 5737 documentation space reserved at /24, and real ISP allocations come from real, RIR-assigned blocks with less "clean" boundaries.
- AS-SETs, communities used for aggregation-related signaling, and confederations (a way to make one AS look like several for internal iBGP-scaling reasons) are mentioned but not built out — they're operational depth beyond a first pass at the concept.
- The Go aggregation code handles only same-size, IPv4 /24-style merges for clarity; production aggregation logic must handle arbitrary prefix lengths and IPv6 as well.

---

## 19. Interview Questions & Model Answers

**Beginner: "What is an Autonomous System, and why does it need its own number?"**

*Model answer:* "An Autonomous System is a network, or group of networks, under one administrative authority with a single, coherent external routing policy — a company, university, or ISP. It needs a unique number, an ASN, because BGP's entire job is to describe reachability in terms of which Autonomous Systems a route passes through — the AS-PATH attribute. Without a stable, unique identifier per organization, BGP couldn't detect routing loops or let operators write policy like 'prefer routes from this AS.'"

**Intermediate: "Why does route aggregation matter at the scale of the global Internet, and what's the trade-off?"**

*Model answer:* "Without aggregation, every individual customer subnet an ISP hands out would need its own entry in the global BGP table, which every other router on the Internet has to store and process on every best-path computation. With hundreds of thousands of ISPs, that would make the table unmanageably large — this was a real, feared crisis in the early-to-mid 1990s that CIDR was specifically designed to prevent. The trade-off is that aggregation hides internal detail: the outside world can't see or route around problems within an aggregated block, and — more seriously — if anyone else announces a more-specific route inside your aggregated space, routers will prefer that more-specific route regardless of your legitimate ownership, which is the exact mechanism BGP hijacking abuses."

**Advanced: "An ISP owns a /16 and aggregates it into one announcement, but you notice via a looking glass that a completely different, unrelated AS is also announcing a /24 that falls inside that /16. What's happening, and how would you confirm whether it's legitimate?"**

*Model answer:* "This is a textbook signature of either a legitimate multihoming arrangement — the /24 might belong to a customer that's intentionally announcing it themselves for a backup link, with the ISP's knowledge — or a route hijack, where an unrelated AS has, accidentally or maliciously, originated a prefix it has no right to. Because BGP applies no default check on origin authorization, longest-prefix match means routers seeing that /24 announcement will prefer it over the ISP's /16 for addresses inside it, regardless of intent. To confirm legitimacy, I'd check whether an RPKI ROA (Chapter 52) exists authorizing that AS to originate that specific prefix, check the prefix's registration in a routing registry like RADB, and if neither confirms it, treat it as a likely hijack and escalate — potentially contacting the legitimate AS's NOC and the hijacking AS's upstream providers to get the false announcement withdrawn."

**Advanced: "Why do most Autonomous Systems on the Internet never appear as a transit hop in anyone's AS-PATH, and what does that imply about how routing policy is actually distributed across the Internet?"**

*Model answer:* "Most ASes are stub or multihomed networks — they originate their own traffic and consume connectivity from one or more upstreams, but they never carry traffic on behalf of a third party, so they simply never show up in the middle of anyone else's AS-PATH. Only ASes that sell transit, by definition, appear as an intermediate hop. That means the actual shape of the global routing graph is far more concentrated than the raw count of ~70,000+ ASNs suggests: a comparatively small number of transit-providing networks carry the bulk of inter-AS path diversity, while the vast majority of ASes are, topologically, leaves hanging off that smaller transit core. This is exactly the structure that makes the Tier-1/Tier-2/Tier-3 hierarchy in Chapter 51 a useful and accurate simplification rather than an arbitrary one — it reflects a real asymmetry in how routing responsibility is distributed, not just a naming convention."

---

## 20. Exercises

### Easy

1. In your own words, explain the difference between a stub AS and a transit AS, and give a real-world example of each.
2. Why is `AS 64512–65534` reserved, and what earlier chapter's addressing scheme is this range directly analogous to?
3. A router announces `10.0.0.0/16` with `summary-only` configured. What does `summary-only` actually change about what gets advertised?

### Medium

4. Given the four /24 blocks `198.51.100.0/24`, `198.51.101.0/24`, `198.51.102.0/24`, `198.51.103.0/24`, what is the single smallest CIDR block that aggregates all four? Show the binary reasoning (which bits are shared).
5. Explain, step by step, why `203.0.112.0/24` and `203.0.113.0/24` can be merged into a `/23`, but `203.0.113.0/24` and `203.0.114.0/24` cannot — even though both pairs are numerically adjacent.
6. What does the `ATOMIC_AGGREGATE` attribute communicate, and why does it carry zero bytes of value?

### Hard

7. Extend the Go `aggregate` function from Section 15 to handle prefixes of *different* input sizes (e.g., a mix of /24s and /25s) by first normalizing appropriately, or explain in detail why that's a materially harder problem than same-size merging.
8. Using the historical routing-table-growth numbers in Section 11, and given that the table grew roughly 5x between 2001 and 2019 despite CIDR and aggregation discipline, argue for or against the claim that aggregation alone is a sufficient long-term solution to routing table growth — and mention at least one recent trend (e.g., growth in multihoming, IPv6 adoption, more specific announcements for anti-hijack/traffic-engineering reasons) that works against aggregation's effectiveness.
9. Design (in words, not code) a policy an ISP could configure to intentionally announce a more-specific /24 route alongside its aggregate /16, specifically to influence which of two upstream transit providers inbound traffic to that /24 arrives through — referencing the "more-specific wins" rule from Section 10.

---

## 21. Summary

| Term | Meaning |
|---|---|
| Autonomous System (AS) | A network under one administrative authority with a coherent, external-facing routing policy |
| ASN | The globally unique number identifying an AS, assigned by an RIR; 16-bit (legacy) or 32-bit (modern) |
| Stub AS | Connects to one upstream, carries only its own traffic |
| Multihomed AS | Connects to multiple upstreams for redundancy, still carries only its own traffic |
| Transit AS | Carries traffic between other ASes as its core function (an ISP) |
| Route aggregation | Advertising one summarized CIDR block instead of many smaller constituent routes |
| AGGREGATOR | BGP attribute identifying which AS/router performed an aggregation |
| ATOMIC_AGGREGATE | Zero-length flag attribute warning that a route is a summarized aggregate |
| Longest-prefix match risk | A more-specific route inside an aggregated block will always be preferred, whether legitimate or hijacked |

Autonomous Systems and aggregation explain the *technical* shape of Internet routing — one number per organization, one summarized route per block. But they say nothing about *why* one AS agrees to carry another's traffic at all, or why some connections are free and others cost real money. That business layer — peering, transit, and the tiered structure of the Internet's backbone — is Chapter 51.
