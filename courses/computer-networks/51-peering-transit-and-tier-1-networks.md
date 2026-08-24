# Chapter 51: Peering, Transit, and Tier-1 Networks

> *"BGP decides which path a packet takes. Money — or the lack of it — decides which paths exist for BGP to choose from in the first place."*

---

## Table of Contents

1. [The Question BGP Never Answers](#1-the-question-bgp-never-answers)
2. [The Naive Assumption: Networks Just Connect](#2-the-naive-assumption-networks-just-connect)
3. [Transit: Paying to Reach the Whole Internet](#3-transit-paying-to-reach-the-whole-internet)
4. [Peering: Free, Mutual Traffic Exchange](#4-peering-free-mutual-traffic-exchange)
5. [Why Would Anyone Peer for Free?](#5-why-would-anyone-peer-for-free)
6. [Peering Policies: Open, Selective, Restrictive](#6-peering-policies-open-selective-restrictive)
7. [Tier-1, Tier-2, and Tier-3 Networks Defined](#7-tier-1-tier-2-and-tier-3-networks-defined)
8. [The Internet's Structure as a Hierarchy of Money](#8-the-internets-structure-as-a-hierarchy-of-money)
9. [Internet Exchange Points (IXPs)](#9-internet-exchange-points-ixps)
10. [How an IXP Actually Works, Mechanically](#10-how-an-ixp-actually-works-mechanically)
11. [Public Peering vs. Private Peering vs. PNI](#11-public-peering-vs-private-peering-vs-pni)
12. [Worked Example: Tracing a Packet Through the Business Layer](#12-worked-example-tracing-a-packet-through-the-business-layer)
13. [A Real Example: Seeing Peering and Transit in a Traceroute](#13-a-real-example-seeing-peering-and-transit-in-a-traceroute)
14. [Hands-On Experiment](#14-hands-on-experiment)
15. [Code: Modeling Peering Policy as a Route Filter in Go](#15-code-modeling-peering-policy-as-a-route-filter-in-go)
16. [Common Misconceptions](#16-common-misconceptions)
17. [Production Notes](#17-production-notes)
18. [What This Chapter Simplified](#18-what-this-chapter-simplified)
19. [Interview Questions & Model Answers](#19-interview-questions--model-answers)
20. [Exercises](#20-exercises)
21. [Summary](#21-summary)

---

## 1. The Question BGP Never Answers

Chapter 49 showed that BGP lets each Autonomous System encode its own policy — LOCAL-PREF, MED, AS-PATH filtering — instead of blindly picking the shortest path. Chapter 50 showed that ASes are the units this policy is expressed between. But neither chapter ever answered the more basic question sitting underneath all of it:

**Why does an AS agree to carry another AS's traffic at all? And why does it choose to prefer one neighbor's routes over another's?**

The honest answer, almost every single time, has nothing to do with technology. It's a **business relationship**. Two networks either have a contract where money changes hands, or they have a handshake agreement (sometimes literally just a handshake at a conference) that it's mutually beneficial to exchange traffic for free. BGP is the mechanism that *implements* whichever relationship exists — but the relationship itself is negotiated by people, in meetings, with lawyers, long before a single BGP session is configured.

This chapter is about that layer: the actual business structure of the Internet, which shapes almost every routing policy decision you saw examples of in Chapter 49.

---

## 2. The Naive Assumption: Networks Just Connect

If you've never thought about it, the natural assumption is something like: "the Internet is a bunch of networks, and they're all just... connected, the way roads connect between towns." This mental model, common among people who first meet Chapter 6's "network of networks" picture, quietly implies that connectivity is free and symmetric — if Network A can reach Network B, that connection cost nothing and both sides benefit equally.

**This is false, and the falseness matters enormously in practice.** Physical links between networks cost real money to build and maintain — routers, fiber, transceivers (Chapter 22), rack space in a facility, staff to manage it all. Someone has to pay for every link that exists. The question the Internet has spent decades answering, informally and through hard-nosed business negotiation, is: **who pays, for which links, and how much?**

Two answers emerged, and almost the entire commercial structure of the Internet is built from just these two:

---

## 3. Transit: Paying to Reach the Whole Internet

**Transit** is the simplest relationship to understand because it maps directly onto something familiar: it's a customer-provider relationship, just like buying electricity or water.

Network A pays Network B money (usually billed by committed bandwidth or actual usage, in dollars per megabit per second) in exchange for Network B carrying Network A's traffic to **anywhere else on the Internet** — not just to Network B itself, but through Network B to every other network Network B can reach, including networks Network B has never directly connected to.

**Intuitive analogy:** transit is exactly like a small retail shop paying a delivery company (FedEx, UPS) to ship its packages anywhere in the world. The shop doesn't need its own trucks reaching every city on Earth — it pays one company that already has that reach, and that company's job is precisely to get the package the rest of the way, wherever "the rest of the way" turns out to be. **Where the analogy breaks:** a delivery company charges per package with a fairly transparent price list; Internet transit pricing is opaque, individually negotiated, and varies wildly by volume, region, and the relative bargaining power of buyer and seller.

The key defining property of transit: **the provider is contractually obligated to carry your traffic to the entire rest of the Internet**, not just to itself. If AS 65001 buys transit from AS 65010, AS 65010 must accept AS 65001's routes and re-advertise them to *all* of AS 65010's own peers, customers, and upstreams — that's what "buying full transit" means, versus, say, a narrower agreement to reach only a subset of destinations.

In BGP terms (tying directly back to Chapter 49, Section 10's worked example): a transit relationship typically gets **lower LOCAL-PREF** on the paying customer's side, because sending traffic over a link you pay for by the byte is more expensive than sending it over a free peering link whenever one exists. This is the exact mechanic that made the "shortest path loses" example in Chapter 49 make business sense.

---

## 4. Peering: Free, Mutual Traffic Exchange

**Peering** is a fundamentally different relationship: two networks agree to exchange traffic **directly and for free**, but strictly limited to traffic **destined for each other** (or each other's own customers) — never traffic destined for some third party.

**Intuitive analogy:** peering is like two neighboring countries agreeing "we'll let each other's mail trucks cross our shared border directly, for free, as long as the mail is actually addressed to someone inside our own country." Country A's mail truck can't use that free crossing to relay mail through Country B on its way to Country C — that would be Country A getting free transit, disguised as peering, and Country B would (rightly) object. **Where the analogy breaks:** peering agreements have no physical border checkpoint enforcing this automatically; enforcement happens entirely through explicit BGP route filtering (only advertising your own and your customers' routes to your peer, never advertising routes you learned from your own transit providers or other peers) and through contract terms — violating this, whether by accident or on purpose, is precisely what a **route leak** is (Chapter 52).

The essential asymmetry that separates peering from transit: **a peer will only tell you about itself and its own downstream customers — never about the rest of the Internet.** If AS 65001 peers with AS 65777, AS 65001 can reach AS 65777 (and AS 65777's own customers) for free over that link, but if AS 65001 wants to reach some unrelated AS 65999, peering with AS 65777 does nothing for that — AS 65001 still needs its own transit or its own separate peering relationship to reach AS 65999.

```
              PEERING                              TRANSIT
     (free, limited scope)                  (paid, full Internet reach)

  AS A ◄──────────────► AS B              AS A ────pays $$$───► AS B
   │                     │                  │                    │
   ▼                     ▼                  ▼                    ▼
 A's own            B's own              A's own          Reaches B AND
 customers          customers            customers        everything B
 only                only                                 can reach --
                                                            the WHOLE
 A will NOT carry    B will NOT carry                      Internet
 traffic from B      traffic from A
 to some third       to some third
 network C           network D
```

---

## 5. Why Would Anyone Peer for Free?

If peering is free, why would a network bother, instead of just selling everyone transit? The answer is that peering is a mutual cost-saving move, not charity, and it makes obvious sense once you see the alternative:

Suppose AS A and AS B exchange a large, steady, roughly balanced amount of traffic with each other — say, both are large regional ISPs whose residential customers constantly stream video hosted by the other's customers, browse each other's customers' sites, and so on. Without peering, **every single byte of that traffic would have to go out through each network's paid transit provider and come back in through the other's paid transit provider** — both networks paying real money, twice over, to move traffic that both networks want moved anyway. By building one direct link (or meeting at a shared facility, Section 9) and agreeing to exchange that traffic for free, **both sides save money simultaneously** compared to paying transit for the same traffic. It also typically improves performance — a direct link is usually a shorter, more predictable, lower-latency path than going out to a transit provider and back.

This is why peering decisions are evaluated in terms of **traffic ratio and mutual benefit**: a network only wants to peer with another network if the relationship saves both parties roughly comparable amounts of money and neither side feels like it's providing an unfairly one-sided subsidy — a concern that Section 6 formalizes.

---

## 6. Peering Policies: Open, Selective, Restrictive

Not every network is equally willing to peer with just anyone. Networks generally publish (informally, or via databases like PeeringDB) one of three peering postures:

- **Open peering policy**: will peer with essentially any network that asks, with minimal requirements. Common for smaller networks and content-heavy networks that benefit from cheap, direct access to as many eyeball networks (residential ISPs) as possible.
- **Selective peering policy**: will peer, but only with networks that meet certain criteria — a minimum amount of traffic exchanged, presence at specific facilities, a roughly balanced traffic ratio, sufficient redundancy (multiple peering points in different regions), and so on. Most large regional ISPs and mid-size backbone networks fall here.
- **Restrictive peering policy** (sometimes called "peer of last resort" posture): peers only with a small, tightly controlled set of networks, often only other very large networks of comparable size and traffic volume — and generally refuses smaller or asymmetric-traffic networks, pushing them toward paid transit instead. This posture is closely associated with — and, per Section 7, is one informal way of recognizing — the very largest Tier-1 networks.

Peering negotiations sometimes genuinely break down, and can become public and contentious: large ISPs and large content/streaming providers have occasionally engaged in **"peering disputes"** where one side demands payment (arguing the traffic ratio is too imbalanced to justify free peering) and the other refuses, sometimes visibly degrading performance for end users caught in the middle until the dispute resolves — a real, recurring feature of Internet business history, not a hypothetical.

---

## 7. Tier-1, Tier-2, and Tier-3 Networks Defined

With peering and transit both defined, the "Tier" classification falls out almost immediately — it's not an official technical standard, but a widely-used industry convention describing a network's position in this business hierarchy:

- **Tier-1 network**: a network that can reach the **entire** global Internet routing table using **only settlement-free peering** — it purchases **zero** transit from anyone. By definition, every other Tier-1 network must be willing to peer with it (because if even one refused, the Tier-1 network couldn't reach that part of the Internet without paying someone for transit, disqualifying it). There is no official registry of "Tier-1" status — it's a claim networks make and the rest of the industry either accepts or disputes based on observed peering behavior — but a commonly cited, broadly agreed-upon list includes networks like **Lumen (formerly CenturyLink/Level 3), NTT Communications, Telia, Tata Communications, Zayo, and GTT**, among a small handful of others (the exact list shifts slowly over time through mergers and market changes).
- **Tier-2 network**: peers with some networks (usually many, often extensively) but still must buy transit from at least one Tier-1 (or large Tier-2) network to reach some portion of the Internet it can't reach via peering alone. The overwhelming majority of the world's mid-size and large regional ISPs are Tier-2.
- **Tier-3 network**: buys transit for essentially all of its Internet connectivity and does little or no peering — typically small, local, or regional ISPs and access networks whose traffic volumes don't justify the operational overhead of negotiating and maintaining peering relationships.

```
                         Tier-1 Networks
              (peer with each other ONLY, zero transit purchased)
        ┌─────────┐   peering   ┌─────────┐   peering   ┌─────────┐
        │ Tier-1 A│◄───────────►│ Tier-1 B│◄───────────►│ Tier-1 C│
        └────┬────┘             └────┬────┘             └────┬────┘
             │ sells transit         │ sells transit         │
             ▼                       ▼                       ▼
        ┌─────────┐             ┌─────────┐             ┌─────────┐
        │ Tier-2  │◄───peering─►│ Tier-2  │             │ Tier-2  │
        │ (also   │             │         │             │         │
        │  peers) │             └────┬────┘             └────┬────┘
        └────┬────┘                  │ sells transit         │
             │ sells transit         ▼                       ▼
             ▼                  ┌─────────┐             ┌─────────┐
        ┌─────────┐             │ Tier-3  │             │ Tier-3  │
        │ Tier-3  │             │(buys all│             │(buys all│
        │         │             │ transit)│             │ transit)│
        └─────────┘             └─────────┘             └─────────┘
```

**Important nuance:** it's entirely normal, and extremely common, for a network to *both* peer extensively *and* buy some transit — the tiers describe a spectrum of self-sufficiency, not a strict caste system. A large Tier-2 ISP might peer directly with hundreds of networks at IXPs while still buying transit from one or two Tier-1s purely as a safety net to guarantee reachability to the small remaining slice of the Internet it doesn't have direct peering relationships with.

---

## 8. The Internet's Structure as a Hierarchy of Money

Stepping back, this whole chapter is really describing one thing: **the "network of networks" picture from Chapter 6 is, underneath the technology, a hierarchy of financial relationships**, roughly (though very imperfectly) shaped like a pyramid — a small number of Tier-1 networks at the top peering only among themselves, a larger layer of Tier-2 networks buying some transit and peering some traffic, and a much larger base of Tier-3 access networks buying nearly all their connectivity.

This isn't a technical requirement of IP or BGP — nothing about the protocols *forces* this hierarchy to exist. It emerged because it's the economically efficient outcome: it would be absurd for every one of the roughly 70,000+ ASes on Earth to negotiate a direct settlement-free peering relationship with every other one (that's the same N² wiring problem from Chapter 3, recast at the business layer instead of the physical layer) — so instead, smaller networks pay larger, already-well-connected networks for aggregated reach, and only the very largest networks, whose mutual traffic volumes justify it, peer directly with each other.

---

## 9. Internet Exchange Points (IXPs)

Peering (Section 4) requires two networks' routers to actually be physically connected somewhere. Building a dedicated point-to-point link between every pair of networks that wants to peer would reintroduce exactly the N² scaling problem Chapter 3 warned about — if 500 networks all want to peer with each other, that's potentially 500×499/2 ≈ 124,750 individual physical links.

**Internet Exchange Points (IXPs)** solve this the same way a switch solves the equivalent problem inside a LAN (Chapter 30): instead of every network wiring directly to every other network, every participating network runs **one physical connection into a shared facility**, and that facility's switching fabric lets any two (or more) participants exchange traffic without needing a dedicated link between each pair.

Real, well-known IXPs include:

| IXP | Location | Notable scale (approximate, changes over time) |
|---|---|---|
| **DE-CIX Frankfurt** | Frankfurt, Germany | One of the world's largest by peak traffic, several hundred Tbps aggregate capacity, thousands of connected networks |
| **AMS-IX** | Amsterdam, Netherlands | One of the oldest and largest, founded 1994 |
| **LINX** | London, United Kingdom | The primary UK exchange, multiple interconnected sites |
| **Equinix IX** | Multiple cities worldwide | A commercial, multi-location exchange operated inside Equinix's own data centers |
| **NAPAfrica** | Johannesburg, South Africa | A major African exchange, notable for driving down local peering costs |
| **Any2 / IX Brasil / DE-CIX New York** | Various | Regional examples across the Americas |

**Intuitive analogy:** an IXP is like a wholesale farmers' market building — instead of every farmer driving a truck to every individual grocery store across the region (an N² problem again), every farmer and every grocery buyer shows up to one shared building, plugs into the same shared loading dock system, and trades directly with whichever other parties they've agreed to trade with. **Where the analogy breaks:** an IXP operator, unlike a market landlord, does not participate in or see the *content* of any of the trades happening across its fabric — it operates purely at Layer 2 (Chapter 28), switching Ethernet frames between participants' routers, and has no visibility into or involvement in the actual BGP relationships (who peers with whom, on what terms) that participants privately negotiate with each other while sitting on the same shared switch.

---

## 10. How an IXP Actually Works, Mechanically

Physically and logically, most IXPs are built around a **shared Layer 2 Ethernet fabric** — a large, often redundant switch infrastructure spanning one or more data centers in a region:

```
   Network A's router          Network B's router          Network C's router
        │                            │                            │
        │ (one physical port each)  │                            │
        └──────────┐        ┌───────┘        ┌───────────────────┘
                    ▼        ▼                ▼
              ┌───────────────────────────────────┐
              │   IXP Shared Layer-2 Switch Fabric  │
              │   (one shared broadcast/VLAN         │
              │    segment, one IP subnet)            │
              └───────────────────────────────────┘

  Each participant's router gets ONE IP address on the IXP's shared
  subnet (e.g. 198.51.100.0/24, purely illustrative).

  A peers with B: A and B configure a direct eBGP session between
    their two IXP-subnet IP addresses -- the underlying Ethernet
    frames physically cross the SAME shared switch fabric as every
    other participant's traffic, but the BGP session itself, and the
    resulting route exchange, is a private, bilateral agreement
    between A and B only.

  A does NOT automatically peer with C just by being on the same
  fabric -- BGP sessions must be explicitly, individually configured
  between each pair of networks that agree to peer.
```

This is worth restating because it surprises people: **joining an IXP does not automatically mean you're peering with everyone else there.** It means you have the *physical capability* to peer with any of them cheaply (one port, one cable, no dedicated circuit needed), but each individual peering relationship (Sections 4-6) is still a separate bilateral negotiation and a separately configured eBGP session.

Many IXPs also offer an optional **route server** — a piece of shared infrastructure that speaks BGP with every participant who opts in, and re-advertises everyone's routes to everyone else who's opted in, dramatically simplifying multilateral peering (instead of configuring N-1 separate BGP sessions to peer with everyone at the exchange, a network can configure one session to the route server and, via a single relationship, exchange routes with dozens or hundreds of other participants at once, as long as all sides agree to the resulting open policy).

---

## 11. Public Peering vs. Private Peering vs. PNI

Three distinct physical arrangements exist for actually exchanging peered traffic, worth distinguishing clearly:

- **Public peering**: happens over an IXP's shared switch fabric (Section 10) — cheap (one port serves potentially many peering relationships) but the available capacity is shared with, and ultimately limited by, the exchange's overall fabric.
- **Private peering / PNI (Private Network Interconnect)**: two networks run a dedicated, direct physical link (a fiber cross-connect, often within the same building or data center campus) solely between the two of them — no shared fabric involved. This is more expensive per relationship (a dedicated port and often a cross-connect fee) but gives guaranteed, dedicated capacity, and is typically what very large networks (especially content networks like Netflix or Google delivering enormous, latency-sensitive traffic volumes to residential ISPs) set up once their traffic with a specific partner grows large enough to justify a dedicated link instead of sharing exchange fabric capacity.
- **Remote peering**: a network reaches an IXP's fabric without physically owning equipment in that city at all, by leasing a Layer 2 circuit from a third party to extend its reach to the exchange — a way for smaller or more distant networks to participate in a major IXP's peering ecosystem without the capital cost of a full physical presence there.

---

## 12. Worked Example: Tracing a Packet Through the Business Layer

Put Chapters 49-51 together with a concrete scenario. A residential customer of "Regional ISP" (a Tier-2 network) requests a video from "StreamCo" (a large content provider):

```
1. Regional ISP checks its BGP table for StreamCo's announced prefix.
   It has TWO candidate routes:

     Route via Tier-1 Transit Provider:
       AS-PATH: [Tier1-AS, StreamCo-AS]     LOCAL-PREF: 100 (paid transit)

     Route via Direct Peering (negotiated at an IXP, Section 9):
       AS-PATH: [StreamCo-AS]                LOCAL-PREF: 200 (settlement-free peer)

2. Applying BGP's best-path algorithm (Chapter 49, Section 9):
   LOCAL-PREF is checked first. 200 beats 100. The peering route wins,
   even though in this case it ALSO happens to have the shorter AS-PATH
   (StreamCo peers directly with Regional ISP at a shared IXP).

3. Traffic flows: Regional ISP's router -> IXP shared fabric (Section 10)
   -> directly into StreamCo's network -> to StreamCo's server -> video
   streams back the same free, direct path.

4. Both sides win financially: Regional ISP avoids paying its Tier-1
   transit provider by the megabit for this popular, high-volume traffic;
   StreamCo avoids paying its own transit provider to reach Regional
   ISP's customers. This is EXACTLY the mutual-benefit case from Section 5
   -- which is why large content providers so aggressively pursue direct
   peering with residential ISPs worldwide.
```

---

## 13. A Real Example: Seeing Peering and Transit in a Traceroute

A `traceroute` (Chapter 54 covers the mechanism in depth) often visibly shows the business layer, if you know what to look for in the hostnames of intermediate hops — many networks name their routers to include facility and peer information:

```bash
$ traceroute www.example-cdn.com
 1  home-router.local (192.168.1.1)         1.2 ms
 2  isp-gateway.regionalisp.net (10.x.x.x)  8.1 ms
 3  core1.chi.regionalisp.net               12.4 ms
 4  ae3.edge-router.chi.regionalisp.net     13.0 ms
 5  ae1.chi-ix.regionalisp.net              13.5 ms     <- IXP-facing interface
 6  peer-gw.chicagoix.examplecdn.com        14.1 ms     <- CDN's IXP-facing router
 7  edge.chi.examplecdn.com                 14.6 ms
 8  cdn-server-142.chi.examplecdn.com       15.0 ms
```

Hop 5's interface name (`ae1.chi-ix...`) and the abrupt handoff directly to a differently-named organization at hop 6, with only a tiny latency increase, is a classic signature of a **direct peering handoff at an IXP** — versus a path that instead shows several more hops through a recognizably different, third-party transit network's naming convention (e.g., `...level3.net`, `...ntt.net`, `...cogentco.com`) sitting between the source and destination, which would indicate the traffic is being carried over paid transit rather than direct peering.

---

## 14. Hands-On Experiment

```bash
# 1. Look up a real network's peering posture and connected IXPs on
#    PeeringDB (the industry's shared, crowdsourced database of exactly
#    this kind of information):
#    https://www.peeringdb.com/net/433   (example: a real network's public record)

# 2. Run traceroute to a few different large content providers from
#    your own connection, and note how many hops appear BEFORE you
#    seem to leave your ISP's own naming convention -- a rough, informal
#    signal of how much peering vs. transit your own ISP relies on for
#    that particular destination:
traceroute www.google.com
traceroute www.netflix.com
traceroute www.wikipedia.org

# 3. Use Hurricane Electric's BGP towebsite (bgp.he.net) to look up any
#    ASN and see its "Peers" tab -- real large networks list dozens to
#    hundreds of settlement-free peers there, a direct look at Section 4
#    in the wild.
```

---

## 15. Code: Modeling Peering Policy as a Route Filter in Go

The single most important *mechanical* enforcement of the peering-vs-transit distinction (Section 4) is a route export filter: **never re-advertise a route learned from one peer/transit provider to another peer**, unless the route belongs to you or your own paying customer. Here's a minimal model of that filter — the same logic, vastly simplified, that a real router's route-map or policy-statement configuration implements:

```go
package main

import "fmt"

// PeerRelationship describes how a route was learned.
type PeerRelationship int

const (
	Customer PeerRelationship = iota // a route from someone who PAYS US (we owe them full reach)
	Peer                              // a route from a settlement-free peer (their own + their customers only)
	Transit                          // a route from someone WE pay (gives us full Internet reach)
)

// Route represents a BGP route annotated with where we learned it from.
type Route struct {
	Prefix     string
	LearnedVia PeerRelationship
}

// exportAllowed implements the fundamental peering-vs-transit export rule:
//   - Routes learned from a CUSTOMER: always safe to advertise to anyone
//     (transit providers, peers, other customers) -- that's the whole
//     point of selling transit to them.
//   - Routes learned from a PEER or a TRANSIT PROVIDER: only safe to
//     re-advertise to our OWN customers (who pay us for full reach) --
//     NEVER to another peer or transit provider, because doing so would
//     mean giving that other network free transit through us, which is
//     exactly what a route leak (Chapter 52) looks like mechanically.
func exportAllowed(learnedFrom PeerRelationship, exportingTo PeerRelationship) bool {
	if learnedFrom == Customer {
		return true // customer routes go everywhere -- that's the service we sell
	}
	// Peer or Transit-learned routes may ONLY go to our own customers.
	return exportingTo == Customer
}

func main() {
	scenarios := []struct {
		route       Route
		exportingTo PeerRelationship
		label       string
	}{
		{Route{"203.0.113.0/24", Customer}, Peer, "customer route -> to a peer"},
		{Route{"203.0.113.0/24", Customer}, Transit, "customer route -> to our transit provider"},
		{Route{"198.51.100.0/24", Peer}, Transit, "PEER-learned route -> to our transit provider"},
		{Route{"198.51.100.0/24", Peer}, Customer, "peer-learned route -> to our own customer"},
		{Route{"192.0.2.0/24", Transit}, Peer, "TRANSIT-learned route -> to another peer"},
	}

	for _, s := range scenarios {
		allowed := exportAllowed(s.route.LearnedVia, s.exportingTo)
		fmt.Printf("%-45s => export allowed: %v\n", s.label, allowed)
	}
	// The two cases marked in ALL CAPS above are exactly the dangerous
	// ones -- exporting a peer- or transit-learned route to another
	// peer or transit provider. exportAllowed correctly returns false
	// for both. A misconfigured router that skips this check is
	// precisely how a route leak happens in the real world.
}
```

---

## 16. Common Misconceptions

- **"Peering is inherently better than transit."** Peering is cheaper and often lower-latency for the specific traffic it covers, but it only reaches the peer's own network and customers — it can never substitute for transit's full-Internet reach. Real networks need both.
- **"Being at an IXP means you're peered with everyone there."** As Section 10 stressed, physical presence at an exchange only makes peering *possible* and cheap — every actual peering relationship is still a separate, explicit bilateral (or route-server-mediated) agreement.
- **"Tier-1 status is an official certification."** There is no governing body that certifies "Tier-1" — it's an industry-recognized claim, validated informally by observing that a network genuinely buys no transit from anyone and is willing to peer with all other claimed Tier-1s.
- **"Peering agreements are always symmetric and simple."** Many are negotiated with specific ratio requirements, minimum traffic commitments, and geographic redundancy requirements — informal "handshake" peering exists, but plenty of peering relationships involve real contracts.

---

## 17. Production Notes

- **PeeringDB** (peeringdb.com) is the industry-standard, crowdsourced, mostly-accurate public database of which networks are present at which IXPs, their peering policies, and traffic volumes — used daily by network engineers negotiating new peering relationships.
- Large content and cloud providers increasingly build **their own private global networks** (sometimes informally called "Tier-0" in casual conversation, though this isn't a standard term) that connect directly to residential ISPs worldwide via extensive private peering and dedicated edge caching infrastructure — this is a major reason a huge fraction of residential Internet traffic today never crosses a traditional Tier-1 backbone at all (previewed further in Chapter 96's CDN discussion).
- Peering disputes are a real operational risk that network engineers plan for: contracts and diversified transit relationships exist specifically so that a breakdown in one peering relationship doesn't cause a customer-visible outage, only a performance degradation as traffic reroutes over paid transit instead.
- The rise of large IXPs in emerging markets (e.g., NAPAfrica, various Latin American exchanges) has been explicitly, publicly credited with substantially lowering local Internet costs and latency, by letting local networks exchange local traffic locally instead of routing it through distant, expensive transit hubs (sometimes previously requiring a round trip through Europe or North America even for traffic between two networks in the same city or country) — a concrete real-world illustration of exactly why the concepts in this chapter matter economically, not just technically.
- **IP transit pricing has fallen dramatically over the Internet's history** — commonly cited industry figures put wholesale IP transit at roughly $1,200 per Mbps per month in the year 2000, falling to single-digit dollars per Mbps per month (and, in the most competitive major hubs, well under a dollar) by the early 2020s, a multiple-orders-of-magnitude decline driven by massive increases in fiber capacity, more competition among transit providers, and exactly the kind of aggressive settlement-free peering growth described in this chapter reducing how much traffic even needs to transit at all. This price collapse is a major reason smaller networks today can afford connectivity that would have been enterprise-only pricing a generation earlier.
- Many large ISPs now published **"peering requirements" documents** on their own websites (searchable, public, and genuinely used by network engineers) spelling out minimum traffic volume, presence at specific facilities, and traffic-ratio requirements a prospective peer must meet — turning what was once an informal, relationship-driven negotiation into a semi-standardized qualification process.

---

### A Worked Example: When Peering Breaks Down

Peering disputes (Section 6) aren't hypothetical — they've played out publicly enough times that the pattern is well understood. A simplified, generic version of how one unfolds:

```
Year 1: Network A (a large ISP) and Network B (a large content/streaming
        network) establish settlement-free peering at three IXPs. Traffic
        is roughly balanced -- A's customers download from B's servers
        at roughly the same volume B's customers (much smaller in this
        example) pull from A -- both sides are happy.

Year 3: B's video streaming product has grown enormously. Traffic across
        the peering link is now wildly asymmetric -- 50:1 or more in
        B's favor. A's network engineering team argues this no longer
        looks like "mutual" traffic exchange -- it looks like B is using
        A's free peering capacity as de facto free transit into A's
        entire residential customer base.

Year 3, continued: A asks B to either pay for the excess capacity, add
        substantially more peering ports at its own cost, or accept a
        formal paid-peering arrangement. B refuses, arguing peering
        should remain free regardless of ratio, since it's still
        traffic destined for A's own customers who explicitly requested
        it (not third-party transit in the Section 4 sense).

Outcome: The peering link saturates and A does not add capacity. Traffic
        that used to flow over the free peering link now competes for
        space on it, causing real, customer-visible degradation --
        buffering, slow load times -- for A's customers trying to reach
        B's service, until a new commercial arrangement (often some form
        of paid interconnection) is reached.
```

This scenario is deliberately generic because near-identical disputes, with the specific pairs of companies varying, have occurred publicly enough times in the industry's history that it represents a genuine, recurring pattern rather than a one-off — the underlying tension (Section 5's "mutual benefit" justification for free peering breaking down once traffic becomes heavily asymmetric) is structural, not particular to any one company's behavior.

---

## 18. What This Chapter Simplified

- Real peering and transit contracts involve pricing models, service-level agreements, and legal terms far beyond the technical/business sketch given here.
- The Tier-1/2/3 taxonomy is a useful mental model but is genuinely fuzzy at the edges in the real industry — some large networks blur the lines, and the informally-agreed Tier-1 list has changed over the decades through mergers and acquisitions.
- IXP fabric technology (redundant switch architectures, route server implementations like BIRD-based ones, DDoS scrubbing services some IXPs now offer) is far more elaborate in production than the single-shared-subnet model shown in Section 10.
- The Go code in Section 15 models only the simplest, single-rule version of export filtering; real router policy languages (route-maps, policy-statements) support far richer conditional logic layered on top of this base rule.

---

## 19. Interview Questions & Model Answers

**Beginner: "What's the difference between peering and transit?"**

*Model answer:* "Transit is a paid customer-provider relationship where the provider carries your traffic to the entire rest of the Internet, not just to itself. Peering is a free, mutual agreement between two networks to exchange traffic only between themselves and their own customers — never as a path to some third network. Peering saves both sides money on traffic they'd otherwise have to send through paid transit anyway, but it can never substitute for transit's full reach, because a peer will never carry your traffic onward to someone else."

**Intermediate: "What makes a network 'Tier-1', and why does that matter for how the Internet is structured?"**

*Model answer:* "A Tier-1 network is one that can reach the entire global Internet using only settlement-free peering — it buys zero transit from anyone. That requires every other Tier-1 network to be willing to peer with it, since if even one refused, it would need to pay someone for transit to reach that missing part of the Internet, disqualifying it. This matters because it creates a natural hierarchy: the Internet isn't flat, it's shaped like a pyramid, with a small number of Tier-1s peering only among themselves at the top, Tier-2 networks buying some transit and peering some traffic in the middle, and Tier-3 access networks buying nearly all their connectivity at the base. That hierarchy emerged for the same reason CIDR and route aggregation exist — full mesh peering between all ~70,000 ASes would be an unmanageable N-squared problem."

**Advanced: "A large streaming company's traffic to a residential ISP suddenly starts taking a longer, higher-latency path through a third-party transit network instead of a direct route. What business-layer explanations would you investigate before assuming it's a technical fault?"**

*Model answer:* "First I'd check whether there was previously a direct peering relationship — at an IXP or via a private interconnect — between the two networks, and whether that session is currently up. A sudden shift to a longer path through transit is the classic signature of a peering relationship going down, whether due to a technical failure of the peering session, a capacity upgrade being delayed, or an actual peering dispute — for instance, the ISP might be demanding payment for what it now considers an imbalanced traffic ratio, and the streaming company refusing, causing both sides to fail over to routing the traffic through paid transit instead while the dispute is negotiated. I'd check BGP session state on both sides' IXP-facing routers, check PeeringDB and public network operator mailing lists or status pages for any announced maintenance or disputes, and only after ruling out a business-layer cause would I dig into a purely technical explanation like a fiber cut or router failure at the peering point."

---

## 20. Exercises

### Easy

1. In one sentence each, define peering and transit, and state the one property that makes them fundamentally different (not just "one is free and one costs money" — what does that difference in cost actually mean for what traffic each carries?).
2. What does it mean for a network to have an "open" peering policy versus a "restrictive" one?
3. Name two real, well-known Internet Exchange Points and the cities they're in.

### Medium

4. Explain why an IXP's shared Layer-2 fabric does not, by itself, create any peering relationships between the networks connected to it.
5. A Tier-2 ISP peers with 40 networks at an IXP and also buys transit from a Tier-1 provider. Explain, using the best-path algorithm from Chapter 49, why traffic to a peer's network will almost always prefer the peering path over the transit path, even if the transit path's AS-PATH were shorter.
6. What's the difference between public peering, private peering (PNI), and remote peering? Give a plausible reason a very large content network might prefer PNI over public peering for its highest-volume traffic.

### Hard

7. Using the `exportAllowed` function from Section 15, extend it to also handle a fourth relationship type, `RouteServer` (a route learned via an IXP route server rather than a direct bilateral session), and reason about which export rules should apply to it.
8. Research (or reason from the economics described in Section 5) why a settlement-free peering relationship can break down into a paid one over time, even with no technical change — what change in *traffic ratio* between the two networks would make the previously-peering network start demanding payment?
9. Explain why "Tier-1" is described in this chapter as an industry claim rather than an official certification, and describe a real-world scenario (a merger, a network changing its transit-buying behavior) that could cause the informally-agreed list of Tier-1 networks to change.
10. Given the historical IP transit price decline described in Section 17 (roughly $1,200/Mbps/month in 2000 down to single-digit dollars/Mbps/month by the early 2020s), argue whether cheaper transit should be expected to make peering *more* or *less* attractive to a mid-size network over time, and justify your answer using Section 5's mutual-cost-saving logic.

---

## 21. Summary

| Term | Meaning |
|---|---|
| Transit | A paid relationship where a provider carries your traffic to the entire rest of the Internet |
| Peering | A free, mutual relationship to exchange traffic limited to each network's own traffic and customers |
| Open/selective/restrictive peering policy | How willing a network is to peer, and under what conditions |
| Tier-1 network | A network reaching the whole Internet via settlement-free peering only, buying zero transit |
| Tier-2 network | A network that both peers and buys some transit |
| Tier-3 network | A network that buys nearly all its connectivity as transit |
| IXP (Internet Exchange Point) | A shared physical facility where many networks connect once to enable cheap peering with each other |
| Public peering | Peering over an IXP's shared switch fabric |
| Private peering / PNI | A dedicated physical link between exactly two networks |
| Route server | Shared IXP infrastructure simplifying peering with many participants via one session |

Peering, transit, and IXPs explain why routes exist and why an AS prefers one over another — but they all assume everyone plays by the rules and only advertises routes they're actually entitled to. Chapter 52 is about what happens when that trust is broken, by accident or on purpose.
