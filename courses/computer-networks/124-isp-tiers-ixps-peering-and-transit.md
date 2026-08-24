# Chapter 124: ISP Tiers, IXPs, Peering, and Transit — Revisited at Global Scale

> *"Chapters 49-51 taught you the pieces: BGP, autonomous systems, peering, transit. This chapter puts them all in one room, at planet scale, and asks the question those chapters left for later: who actually runs this thing? The honest answer is nobody — and that turns out to be the point."*

---

## Table of Contents

1. [What Chapters 49-51 Taught, and What's Missing](#1-what-chapters-49-51-taught-and-whats-missing)
2. [The Naive Picture: The Internet as a Designed Hierarchy](#2-the-naive-picture-the-internet-as-a-designed-hierarchy)
3. [The Real Picture: A Marketplace, Not a Blueprint](#3-the-real-picture-a-marketplace-not-a-blueprint)
4. [Tier-1 Networks, Revisited as a Club With No Membership Office](#4-tier-1-networks-revisited-as-a-club-with-no-membership-office)
5. [Tier-2: The Networks That Actually Do Most of the Work](#5-tier-2-the-networks-that-actually-do-most-of-the-work)
6. [Tier-3: Where Actual Humans Connect](#6-tier-3-where-actual-humans-connect)
7. [The Networks That Don't Fit the Tiers At All](#7-the-networks-that-dont-fit-the-tiers-at-all)
8. [Inside a Real IXP: DE-CIX Frankfurt and AMS-IX](#8-inside-a-real-ixp-de-cix-frankfurt-and-ams-ix)
9. [The Peering LAN, Route Servers, and Onboarding a New Member](#9-the-peering-lan-route-servers-and-onboarding-a-new-member)
10. [A Full Packet Trace Across the Whole Tier System](#10-a-full-packet-trace-across-the-whole-tier-system)
11. [Who Governs the Internet, If Nobody Runs It?](#11-who-governs-the-internet-if-nobody-runs-it)
12. [Real Numbers: IXP Traffic and Membership at Scale](#12-real-numbers-ixp-traffic-and-membership-at-scale)
13. [Hands-On: Exploring the Marketplace Yourself](#13-hands-on-exploring-the-marketplace-yourself)
14. [Code: Simulating the Tier System in Go](#14-code-simulating-the-tier-system-in-go)
15. [Common Misconceptions](#15-common-misconceptions)
16. [Production Notes](#16-production-notes)
17. [What This Chapter Simplified](#17-what-this-chapter-simplified)
18. [Interview Questions & Model Answers](#18-interview-questions--model-answers)
19. [Exercises](#19-exercises)
20. [Summary, and the Bridge to Chapter 125](#20-summary-and-the-bridge-to-chapter-125)

---

## 1. What Chapters 49-51 Taught, and What's Missing

Chapter 49 explained BGP as a path-vector protocol built for policy. Chapter 50 explained the Autonomous System as the unit that policy is expressed between. Chapter 51 explained the two relationships — peering and transit — and sketched the resulting Tier-1/2/3 pyramid, plus a first look at Internet Exchange Points.

All of that was necessarily taught **one relationship at a time**: one AS, comparing two candidate routes, choosing based on LOCAL-PREF. That's correct, but it's a single frame from a much bigger movie. This chapter zooms out to the whole reel: tens of thousands of these individual decisions, made independently, with no central coordinator, somehow adding up to a single connected Internet that a user in Lagos can use to load a page hosted in Oregon. It also goes somewhere Chapter 51 only sketched: what an IXP building actually looks like, who works there, and what a new member has to do to get connected.

---

## 2. The Naive Picture: The Internet as a Designed Hierarchy

Most people's mental model of "the Internet," if pressed, is something like a clean org chart: a handful of Tier-1 "backbone" companies at the top, ISPs in the middle buying from them, and everyone else at the bottom buying from the ISPs — a single, static pyramid, engineered top-down like a road network with a central highway authority.

**This naive picture is wrong in a specific, important way: nobody designed the pyramid, and it isn't static.** There is no Internet Routing Authority that assigns a network's tier, approves its peering relationships, or enforces the shape of the hierarchy. The pyramid in Chapter 51's diagram is an *emergent* pattern — the accumulated result of tens of thousands of independent, self-interested, bilateral business decisions, each made by one network's own engineers and executives, optimizing for that network's own costs and reliability, with zero visibility into the global picture they're collectively producing.

---

## 3. The Real Picture: A Marketplace, Not a Blueprint

The accurate mental model is an economic one: **the Internet's physical topology is the visible residue of tens of thousands of separately negotiated contracts**, the same way a country's supply chains are the residue of separately negotiated purchase orders, not a single planned distribution network. Nobody sat down in 1995 and decided "here is the final shape interconnection should take." Instead:

- A small ISP in Nairobi decides buying transit from a single European provider is too slow and too expensive for local traffic, so it joins the local IXP (Section 8) instead.
- A streaming company decides its peering costs with a particular regional ISP have become large enough to justify a dedicated private link (Chapter 51, Section 11) instead of continuing over shared exchange fabric.
- A mid-size ISP that used to buy transit exclusively from one Tier-1 diversifies to two providers after one outage cost it a very public, very expensive few hours of downtime.

Every one of these is a local, self-interested decision. None of them requires (or waits for) permission from any global authority. And yet the sum of millions of such decisions, repeated and adjusted continuously over three decades, is a network that — almost all of the time — successfully routes a packet between any two of its ~70,000+ participating networks. This is the single most important idea in this chapter: **the Internet's connectivity is an emergent property of an incentive structure, not a specification anyone implemented.**

```
     Chapter 3's original insight, recast at the business layer:

     Full-mesh direct connections between N networks   ~  N² contracts
     (nobody actually does this -- same math, same reason it fails,
      as full-mesh WIRING failed for N computers in Chapter 3)

     What actually happens instead:
       - A minority of networks (Tier-1s) settle for peering only
         among themselves -- a small, dense core.
       - A much larger set of networks (Tier-2s) buy aggregated
         reach from a few providers instead of negotiating with
         everyone individually.
       - The largest set (Tier-3s, content networks) mostly buy,
         rarely sell, and increasingly build their own private
         reach directly to the edge (Chapter 127).

     This is EXACTLY the same hierarchical-aggregation trick that
     made CIDR (Ch. 39) and route aggregation (Ch. 50) work at the
     addressing layer -- applied here to the business layer instead.
```

---

## 4. Tier-1 Networks, Revisited as a Club With No Membership Office

Chapter 51 defined a Tier-1 network as one reaching the entire Internet via settlement-free peering alone, buying zero transit. Worth restating at global scale: **there is no certifying body, no application form, no annual audit.** "Tier-1" is a claim a network makes about its own peering relationships, which the rest of the industry either believes (because that network is, observably, willing to peer with every other claimed Tier-1, and visibly buys no transit from anyone) or disputes.

A commonly cited (and slowly shifting) list, as of the mid-2020s, includes networks such as **Lumen (formerly Level 3/CenturyLink), NTT Communications, Telia Carrier (Arelion), Tata Communications, Zayo, GTT, and Cogent** — though Cogent's precise status is a long-running point of industry debate, since it has historically had contentious depeering disputes with other large networks precisely over whether its traffic ratios justified settlement-free treatment. This dispute-and-negotiate texture is itself the point: **Tier-1 status is continuously re-litigated in practice, not permanently awarded.**

The Tier-1 "club" behaves less like a hierarchy and more like a small, tightly-knit **cartel of mutual necessity**: each member needs every other member to keep peering with it, because losing even one Tier-1 peer would force buying transit somewhere — which would, by definition, disqualify it from Tier-1 status. This mutual dependency is what keeps the small club stable even with zero formal enforcement: self-interest alone holds it together.

---

## 5. Tier-2: The Networks That Actually Do Most of the Work

If Tier-1 networks are a small, famous club, **Tier-2 networks are the vast, unglamorous majority of the Internet's actual carrying capacity** — the large national and regional ISPs, cable operators, and mobile carriers that most people have actually heard of: Comcast, Deutsche Telekom, Orange, Telstra, Airtel, Vodafone, and thousands of smaller counterparts worldwide.

A Tier-2 network typically:

- Peers extensively — often with hundreds of other networks at multiple IXPs (Section 8) — to keep as much of its traffic off paid transit as possible.
- Still buys transit from one or more Tier-1s (or large Tier-2s) as a safety net, guaranteeing full reachability to the long tail of networks it hasn't negotiated direct peering with.
- Often *sells* transit itself, to smaller Tier-3 networks in its region — meaning the same network is simultaneously a transit customer (upstream) and a transit provider (downstream), a completely normal, extremely common position in the hierarchy.

This dual role is why Chapter 51's pyramid diagram should be read as describing a *relationship*, not a fixed caste: the same AS number can appear as a customer in one BGP session and a provider in another, at the same moment, on different links.

---

## 6. Tier-3: Where Actual Humans Connect

Tier-3 networks — often called **eyeball networks**, because their defining trait is having actual human end users (eyeballs) rather than servers — buy nearly all their connectivity as transit and do little or no settlement-free peering. Small regional ISPs, municipal broadband providers, and many enterprise networks fall here.

Two things are worth being explicit about, because they're easy to miss from Chapter 51's zoomed-in view:

1. **"Eyeball network" versus "content network" is a useful, real industry distinction, not a formal tier.** An eyeball network's traffic is overwhelmingly *inbound* (users downloading video, web pages, app data). A content network's traffic is overwhelmingly *outbound* (servers sending video, web pages, app data). This asymmetry is exactly what drives Section 5's peering-ratio disputes (Chapter 51, Section 6) — a content network sending a large, one-directional flow into an eyeball network's users looks, from a pure traffic-ratio standpoint, uncomfortably close to free transit, which is precisely why eyeball networks scrutinize these relationships so closely.
2. **Tier-3 status is about self-sufficiency, not size.** A Tier-3 network can still be very large in absolute subscriber count (a national mobile carrier serving tens of millions of phones) while still buying essentially all its upstream Internet reachability as paid transit, because peering only makes financial sense when both sides' traffic volumes and ratios justify the operational overhead (Chapter 51, Section 5) — many large eyeball networks simply never accumulate enough symmetric, bilateral traffic with any one partner to bother.

---

## 7. The Networks That Don't Fit the Tiers At All

The Tier-1/2/3 model was coined to describe an Internet that looked a certain way in the 1990s and 2000s: networks whose job was fundamentally to **carry** other people's traffic. A newer category of network doesn't fit that model cleanly at all: **large content and cloud providers** — Google, Meta, Netflix, Amazon, Microsoft, Cloudflare — that generate or serve an enormous fraction of all Internet traffic themselves, rather than carrying third parties' traffic between them.

These networks often:

- Peer as aggressively and broadly as any Tier-1, sometimes more so (open peering policies at nearly every major IXP worldwide), specifically to get their traffic as close to eyeball networks as possible.
- Buy little or no traditional transit for their core traffic, not because they qualify as Tier-1 in the classic sense, but because they increasingly build **their own private global backbone** connecting their own data centers and edge locations directly (a full architectural treatment is Chapter 127's entire subject).
- Are sometimes informally called "Tier-0" or "hyperscale" networks in casual industry conversation — not a standardized term, but a real, widely-understood acknowledgment that the classic three-tier model doesn't capture what these networks actually are.

This is worth flagging clearly here because it previews the course's own trajectory: Chapter 125 revisits how these networks route users to themselves at scale, Chapter 126 revisits who now physically owns the undersea cables carrying their traffic, and Chapter 127 is the full synthesis of how they run networks that increasingly bypass the traditional Tier-1/2/3 marketplace described in this chapter entirely.

---

## 8. Inside a Real IXP: DE-CIX Frankfurt and AMS-IX

Chapter 51, Section 9 introduced IXPs conceptually. It's worth actually describing what one physically is, using two of the world's largest as concrete examples.

**DE-CIX Frankfurt**, operated by DE-CIX Management GmbH, is not one building but a **switching fabric distributed across more than a dozen carrier-neutral data centers** throughout the Frankfurt metro area (facilities operated by companies like Equinix, Interxion/Digital Realty, and others), all interconnected by DE-CIX's own fiber links so that a router plugged in at any one site can reach every other member regardless of which building they're physically in. As of the mid-2020s, DE-CIX Frankfurt reports **peak traffic well into the several-Tbps range (historically among the highest of any exchange globally) and several thousand connected ports from well over 900 member networks** spanning ISPs, content networks, cloud providers, and enterprises. Frankfurt's role as a exchange hub is itself a historical accident worth knowing: Germany's position as a geographic and political crossroads of Europe, plus early, deliberate investment by German network operators in shared neutral interconnection, made it the natural place for many European and international networks to meet.

**AMS-IX (Amsterdam Internet Exchange)**, founded in 1994 by a small group of Dutch ISPs, is one of the oldest exchanges still operating and one of the largest by connected member count, also spanning multiple data centers around Amsterdam rather than a single physical building, following the identical multi-site model DE-CIX uses.

**What actually sits in these buildings:** row after row of standard data-center racks, but instead of servers, they're filled with **very large, high-port-density Ethernet switches** (the "switching fabric" from Chapter 51, Section 10) — enterprise/carrier-grade gear from vendors like Nokia, Juniper, or Arista, configured essentially as one enormous, redundant Layer-2 broadcast domain spanning the whole metro area. A member network doesn't ship a whole router to the IXP; it runs a **fiber cross-connect** from its own router (physically located in the same building, or reached via a leased circuit — Chapter 51's "remote peering") to a port on the IXP's switch fabric. That's the entire physical footprint of "joining an IXP": one port, one cross-connect, one router interface configured with an IP address on the exchange's shared peering LAN subnet.

```
                DE-CIX / AMS-IX metro-area fabric (simplified)

  Data Center A            Data Center B            Data Center C
 ┌─────────────┐          ┌─────────────┐          ┌─────────────┐
 │ Member A's   │          │ Member B's   │          │ Member C's   │
 │ router  ──┐  │          │ router  ──┐  │          │ router  ──┐  │
 │           │  │          │           │  │          │           │  │
 │      IXP switch         │      IXP switch         │      IXP switch
 │      (this site)  ══════╪══════ (this site) ══════╪══════ (this site)
 └─────────────┘   fiber   └─────────────┘   fiber   └─────────────┘
                  between IXP's OWN switches at each site

  From any member's point of view: it's ONE flat peering LAN,
  regardless of which physical building they're plugged into.
```

---

## 9. The Peering LAN, Route Servers, and Onboarding a New Member

Every IXP publishes a single shared IP subnet — the **peering LAN** — that every connected member's router gets one address on (e.g., a /22 or /23 IPv4 block reserved specifically for this purpose, plus an IPv6 equivalent). Chapter 51, Section 10 already showed the logical picture; here is the actual member-onboarding sequence a real network engineer follows to join one:

1. **Apply for membership** with the IXP operator, providing the applicant's own ASN, expected traffic volume, and the data center site(s) where they can present a cross-connect.
2. **Order a physical port** (common increments: 1G, 10G, 100G, and increasingly 400G at the largest exchanges) and a cross-connect from the IXP's operator or the data center's own cross-connect service.
3. **Configure the router interface** with the assigned peering-LAN IP address and bring up the physical link — at this point the member can *see* the shared Ethernet segment but has exchanged zero routes with anyone.
4. **Decide on a peering strategy:** configure individual bilateral eBGP sessions with specific other members whose peering policies (Chapter 51, Section 6) and traffic profiles make sense, and/or configure a single eBGP session to the exchange's **route server** — a piece of IXP-operated infrastructure (commonly running open-source BGP daemons like BIRD) that re-advertises every opted-in member's routes to every other opted-in member, turning what would otherwise be dozens or hundreds of individually-configured sessions into one.
5. **Publish the resulting presence on PeeringDB** (Section 12) so other networks can find and request peering, the same industry-standard directory Chapter 51, Section 17 introduced.

The route server is worth dwelling on because it changes the *shape* of peering at scale: without it, joining an exchange with 900 members and wanting to peer with, say, 200 of them would mean configuring, monitoring, and troubleshooting 200 separate BGP sessions. With a route server, a network configures **one** session to the route server, and — as long as both sides have opted their routes into it — effectively multilateral-peers with every other opted-in member through that single relationship. Most large exchanges report a meaningful fraction of their members using route servers specifically for this reason, alongside a smaller set of very high-volume bilateral sessions for their most important, highest-traffic relationships (which still get dedicated bilateral BGP sessions for finer policy control and dedicated capacity planning).

---

## 10. A Full Packet Trace Across the Whole Tier System

Put Sections 4-9 together into one continuous, worked trace. A residential customer of "RegionalISP" (Tier-2) in Warsaw requests a page from "GlobalContent" (a large content network, Section 7), whose nearest presence is a data center in Frankfurt reachable via DE-CIX.

```
1. RegionalISP's edge router looks up GlobalContent's announced prefix.
   Two candidates exist:

     a) Via Transit Provider "BackboneCo" (a Tier-1):
        AS-PATH: [BackboneCo-AS, GlobalContent-AS]
        LOCAL-PREF: 100 (paid transit -- Ch. 51 Section 3)

     b) Via direct peering at DE-CIX Frankfurt (Section 8):
        AS-PATH: [GlobalContent-AS]
        LOCAL-PREF: 200 (settlement-free peer -- Ch. 51 Section 4)

2. Best-path algorithm (Ch. 49 Section 9): LOCAL-PREF checked first.
   200 beats 100. RegionalISP's traffic to GlobalContent flows over
   the free DE-CIX peering link, never touching BackboneCo at all.

3. Physically: RegionalISP's router in Warsaw sends the packet over
   its own backbone to its router physically present at a DE-CIX
   Frankfurt site (Section 8), across the IXP's shared fabric to
   GlobalContent's own router at the same exchange, then onto
   GlobalContent's private network to wherever the actual server
   or cache lives (previewed fully in Ch. 127).

4. RegionalISP never paid BackboneCo a cent for this specific
   request. BackboneCo's entire involvement in this transaction:
   none -- it's simply the safety-net path RegionalISP keeps paying
   for to reach the smaller, long-tail networks it hasn't bothered
   to negotiate direct peering with.
```

This single trace touches every layer of the marketplace at once: a Tier-2 eyeball network, a Tier-1 safety-net transit provider that ends up unused for this particular flow, an IXP's physical switching fabric, and a large content network that doesn't cleanly fit any tier at all — exactly Sections 4-7's taxonomy, operating simultaneously, with zero central coordination deciding any of it in real time.

---

## 10b. A Real Historical Peering Dispute, With Numbers

Section 6's generic dispute pattern happened, in a well-documented, publicly reported real form, between Netflix and several large US residential ISPs (including Comcast and Verizon) around 2013-2014. It is worth walking through because it shows Sections 3-7's abstractions playing out with real commercial stakes.

By the early 2010s, Netflix's video streaming traffic had grown to represent a very large share — publicly reported estimates at the time put Netflix at roughly a third of all North American fixed-line Internet traffic during peak evening hours — flowing overwhelmingly in one direction: from Netflix's content delivery infrastructure into residential ISPs' networks. Netflix had been relying partly on third-party transit and CDN providers to reach these ISPs, and those interconnection links became publicly, measurably congested, with independently published streaming-quality data showing degraded video performance for Netflix subscribers on the affected ISPs during this period.

The affected ISPs argued the traffic ratio was now so overwhelmingly one-directional that continuing to accept it without any payment (whether structured as paid interconnection, paid peering, or a similar arrangement) was not equivalent to the mutually-beneficial relationship peering is supposed to represent (Section 5's mutual-benefit logic). Netflix, and some public commentary at the time, framed the same situation differently — as ISPs deliberately allowing congestion to build in order to extract payment for what should, in this view, be treated as ordinary, reciprocal interconnection. Netflix ultimately entered into paid direct-interconnection agreements with the major ISPs involved, after which the publicly measured streaming performance issues were reported to improve.

**Why this real case matters for this chapter's argument:** it is a genuine, publicly documented instance of Section 6's generic pattern (an initially free or third-party-mediated relationship becoming untenable as traffic asymmetry grows), it shows the dispute playing out as *public* controversy rather than a purely private, invisible negotiation (regulatory bodies and the press covered it extensively at the time), and it is a large part of *why* large content networks (Chapter 124, Section 7) increasingly build their own direct, extensive private interconnection and caching infrastructure (previewed here, covered fully in Chapter 127) rather than relying on third-party transit or standard peering to reach residential ISPs at all.

---

## 11. Who Governs the Internet, If Nobody Runs It?

It's worth being precise about what *does* exist centrally, because "nobody runs the Internet" is easy to overstate into "there's no coordination at all," which isn't true either. A small number of real, functioning bodies coordinate specific, narrow technical resources — but critically, **none of them controls routing, peering, or transit decisions**:

| Body | What it actually coordinates | What it does NOT control |
|---|---|---|
| **IETF (Internet Engineering Task Force)** | Publishes the RFCs that define protocols (BGP itself, TCP, IP, DNS) | Whether any network chooses to use a protocol correctly, or at all |
| **ICANN / IANA** | Coordinates the DNS root zone and allocates large address blocks to Regional Internet Registries | Which networks peer, or how any AS routes its own traffic |
| **RIRs (ARIN, RIPE NCC, APNIC, LACNIC, AFRINIC)** | Allocate IP address blocks and AS numbers (Ch. 50) to networks in their region | Any network's business or peering decisions once it has its numbers |
| **IXP operators (DE-CIX, AMS-IX, LINX, etc.)** | Operate the shared physical switching fabric (Section 8) | Which of their members choose to peer with which others, or on what terms |

Every one of these bodies coordinates a specific, narrow resource — names, numbers, and shared physical fabric — and every one of them deliberately stays out of the actual routing and business decisions covered by this chapter. **The "who decides what path my packet takes" question has no institutional answer**, because the honest answer is: thousands of independent network engineers, at thousands of independent companies, each optimizing their own AS's BGP policy (Chapter 49) according to their own peering and transit contracts (Section 3) — with the emergent result being the Internet.

---

## 11b. A Field Guide: Reading a Network's Business Position From the Outside

Because none of this is centrally published (Section 11), engineers routinely have to infer a network's business position from observable signals. A quick reference for what to actually look at:

| Signal | What it suggests |
|---|---|
| Zero upstream providers listed on bgp.he.net | Likely claims Tier-1 status (Section 4) |
| Very high peer count (hundreds) on PeeringDB, open peering policy | Likely a large Tier-2, or a content/hyperscale network (Section 7) chasing broad reach |
| Small peer count, one or two upstream transit providers listed | Likely a Tier-3 / regional eyeball network (Section 6) |
| AS-PATH in a traceroute repeatedly passes through the same well-known transit AS name across many different destinations | That network is likely this ISP's primary paid upstream (Section 3) |
| Router hostnames referencing a city + "ix" pattern, with an abrupt organization change one hop later | A direct peering handoff at an IXP (Chapter 51, Section 13) |
| A network's PeeringDB entry lists dozens of facility presences worldwide | Likely a hyperscaler or major CDN building broad physical reach (Section 7, and Chapter 127 in full) |

None of these signals is individually conclusive — PeeringDB entries can be stale, and a network's self-reported peering policy is aspirational, not a guarantee — but combined, they let an outside engineer reconstruct a rough, usually accurate picture of where a given AS sits in the marketplace described in this chapter, purely from public data.

---

## 12. Real Numbers: IXP Traffic and Membership at Scale

Concrete figures (approximate, and changing continuously, but structurally accurate as of the mid-2020s) to anchor Sections 4-9's descriptions in scale:

| Metric | Approximate figure |
|---|---|
| Number of IXPs worldwide | Roughly 600-700+ across essentially every country with meaningful Internet infrastructure |
| DE-CIX Frankfurt peak traffic | Regularly cited in the several-Tbps range at peak, among the highest of any single exchange globally |
| AMS-IX connected networks | Several hundred to close to a thousand member networks, depending on how affiliated regional sites are counted |
| Largest exchanges by member count | Often exceed 900-1,000 individually connected ASes at a single exchange |
| Global count of Autonomous Systems | Roughly 70,000-75,000+ actively announced ASNs in the global routing table |
| Global IPv4 routing table size | Roughly 950,000+ distinct announced prefixes |

The scale gap between "70,000+ ASes exist" and "under 1,000 of them show up at even the largest single exchange" is itself informative: the overwhelming majority of the world's Autonomous Systems are small Tier-3 networks that never need, or can't justify, direct IXP presence at all — they simply buy transit from a Tier-2, which in turn peers on their behalf at the exchanges that matter, exactly the aggregation pattern Section 3 described.

---

## 13. Hands-On: Exploring the Marketplace Yourself

```bash
# 1. Look up a real large IXP's live statistics page (most publish these
#    publicly, updated in near-real-time):
#    https://www.de-cix.net/en/locations/germany/frankfurt/statistics
#    https://www.ams-ix.net/ams/statistics

# 2. Search PeeringDB for a real network you use daily and see every
#    IXP it's connected to, and its declared peering policy:
#    https://www.peeringdb.com  (search e.g. "Netflix", "Cloudflare",
#    or your own ISP's name)

# 3. Use bgp.he.net to look up a claimed Tier-1's ASN and inspect its
#    "Peers" tab -- count roughly how many peers it lists, and note
#    that it should show ZERO upstream transit providers, which is
#    the practical, observable signature of Tier-1 status (Section 4).

# 4. Traceroute to a destination you know is served by a large content
#    network, and look for a hostname segment referencing a real IXP
#    (common patterns include city+ix names in router hostnames):
traceroute www.some-large-site.example
```

---

## 14. Code: Simulating the Tier System in Go

A small simulation showing how independent, self-interested peering/transit decisions by individual ASes converge into the pyramid shape from Section 3 — with no central coordinator making the shape happen:

```go
package main

import "fmt"

type Tier int

const (
	Tier1 Tier = iota
	Tier2
	Tier3
)

type AS struct {
	Name        string
	Tier        Tier
	PeersWith   []string // settlement-free peers
	BuysTransit []string // paid upstream providers
}

// classify infers an AS's effective tier purely from its OWN local
// relationships -- exactly how the real Internet "decides" tiers:
// no central authority, just an observable pattern.
func classify(a AS) Tier {
	if len(a.BuysTransit) == 0 {
		return Tier1 // reaches everyone via peering alone
	}
	if len(a.PeersWith) > 0 {
		return Tier2 // mixes peering and transit
	}
	return Tier3 // transit only
}

func main() {
	networks := []AS{
		{Name: "BackboneCo", PeersWith: []string{"OtherTier1A", "OtherTier1B"}},
		{Name: "RegionalISP", PeersWith: []string{"GlobalContent"}, BuysTransit: []string{"BackboneCo"}},
		{Name: "SmallTownISP", BuysTransit: []string{"RegionalISP"}},
	}

	for _, n := range networks {
		t := classify(n)
		label := map[Tier]string{Tier1: "Tier-1", Tier2: "Tier-2", Tier3: "Tier-3"}[t]
		fmt.Printf("%-14s peers=%-2d transit-providers=%-2d  => classified as %s\n",
			n.Name, len(n.PeersWith), len(n.BuysTransit), label)
	}
	// Nobody assigned these labels -- they fall directly out of each
	// AS's own, independently-made peering and transit decisions,
	// exactly as Section 3 argues happens on the real Internet.
}
```

---

## 15. Common Misconceptions

- **"There's a governing body that decides the Internet's routing topology."** As Section 11 detailed, real bodies coordinate names, numbers, and shared physical fabric — none of them touches routing or business decisions, which remain fully decentralized.
- **"Tier-1 status is permanent."** It's continuously re-litigated through ongoing peering relationships; a merger, an acquisition, or a large enough dispute can change the informally-agreed list (Section 4).
- **"Joining an IXP automatically gives you transit-free reachability to the whole Internet."** No — Section 9 was explicit that presence at an exchange only creates the *opportunity* for cheap peering; a network still typically needs paid transit as a safety net unless its peering coverage is genuinely comprehensive (which is essentially the definition of Tier-1).
- **"Content networks like Google or Netflix are Tier-1 ISPs."** They don't fit the classic model at all (Section 7) — they generate/serve traffic rather than carrying third parties' traffic, and increasingly bypass the traditional marketplace with their own private backbones (Chapter 127).

---

## 16. Production Notes

- Real network engineers treat **IXP diversity** (presence at multiple, geographically separate exchanges) as a basic resilience requirement, mirroring the same "don't rely on one path" principle from packet switching itself (Chapter 9) and submarine cable redundancy (Chapter 126).
- **PeeringDB accuracy is community-maintained and imperfect** — production peering coordinators cross-check it against direct outreach and NOC-to-NOC relationships rather than trusting it as a sole source of truth.
- Large IXPs increasingly offer value-added services beyond plain Ethernet switching — DDoS scrubbing, blackholing services (Chapter 83), and even direct cloud on-ramps letting members reach major cloud providers' networks without separate cross-connects.
- Emerging-market IXP growth (Section 3's Nairobi example is realistic, not hypothetical) is an active, ongoing area of Internet infrastructure development, explicitly credited by network operators and NGOs with measurably lowering local access costs and latency in regions that previously depended on distant transit hubs.

---

## 17. What This Chapter Simplified

- Real IXP fabric architectures involve redundant switch pairs, multiple routing-engine failover mechanisms, and often several independent route server instances for fault tolerance — simplified here to one logical shared segment.
- The Tier-1 list given in Section 4 is illustrative and genuinely contested at the margins; treat it as directionally accurate, not an authoritative registry (because, per Section 4 itself, no such registry exists).
- Real onboarding (Section 9) involves considerably more paperwork, SLAs, and cross-connect logistics than the five-step summary given here.

---

## 17b. What's on Every IXP Member's Checklist Before Going Live

A short, concrete summary of Section 9's onboarding steps, worth having as a standalone checklist:

- [ ] ASN and prefix registration confirmed with the relevant RIR (Chapter 50).
- [ ] Membership application approved by the IXP operator.
- [ ] Physical port ordered at the correct speed tier for expected traffic.
- [ ] Cross-connect ordered and installed to the member's own router.
- [ ] Peering-LAN IP address (v4 and v6) configured on the router interface.
- [ ] Route filters and prefix limits configured (Chapter 49, Section 16) before accepting any routes at all.
- [ ] Decision made: route server, bilateral sessions, or both.
- [ ] PeeringDB entry published and kept current (Section 12).
- [ ] Monitoring in place for the new BGP sessions' state (Established vs. flapping).

Skipping the route-filter step in particular is a direct, real cause of exactly the route-leak failure mode Chapter 52 covers — a brand-new IXP member accepting a peer's entire unfiltered table, with no limits, is a classic, avoidable misconfiguration.

---

## 18. Interview Questions & Model Answers

**Beginner: "Is there a company or organization that runs the Internet?"**

*Model answer:* "No single entity runs it. A handful of bodies coordinate narrow technical resources — ICANN and the RIRs hand out IP addresses and AS numbers, the IETF publishes the protocol specifications, and IXP operators run shared physical switching fabric — but none of them decides which networks peer with which others, or which paths traffic takes. Those are independent, bilateral business decisions made by tens of thousands of separately-owned networks, and the Internet's actual topology is just the accumulated result of those decisions."

**Intermediate: "Why does an IXP's route server matter for how peering scales?"**

*Model answer:* "Without it, peering with N other networks at an exchange means configuring, monitoring, and troubleshooting N separate bilateral BGP sessions — the same N-squared-style scaling problem the course keeps running into, from full-mesh wiring in Chapter 3 to full-mesh peering itself. A route server lets a network configure one BGP session and, as long as both sides opt in, exchange routes with every other opted-in member through that single relationship, turning an operationally expensive many-to-many problem into a one-to-one one."

**Advanced: "A large content network's traffic to a specific regional eyeball ISP is growing rapidly. Walk through how that relationship is likely to evolve, referencing the concepts in this chapter."**

*Model answer:* "Initially the two networks probably meet at a shared IXP and set up ordinary public peering (Chapter 51, Section 11) since traffic volumes are modest. As the content network's traffic grows, it starts consuming a meaningful share of the exchange's shared fabric capacity on that link, and eventually the pair typically moves the highest-volume portion of the relationship to a dedicated private peering interconnect (PNI) to get guaranteed capacity outside the shared exchange fabric. If growth continues and the traffic ratio becomes heavily asymmetric, the eyeball ISP may start questioning whether continued free peering is fair, potentially triggering the kind of dispute described in Chapter 51 — or, increasingly in practice, the content network may instead build out its own private caching infrastructure physically inside the ISP's own network (previewed in Chapter 127), sidestepping the exchange relationship almost entirely for that traffic."

---

## 19. Exercises

### Easy

1. In your own words, explain why "the Internet's topology is an emergent economic pattern, not a designed hierarchy" is a more accurate statement than "the Internet is organized as Tier-1, Tier-2, and Tier-3 networks by design."
2. Name two things ICANN/IANA coordinate, and one thing they explicitly do not control.
3. What is a "route server," and what specific scaling problem does it solve for a network joining a large IXP?

### Medium

4. A network wants to know whether it should describe itself as Tier-1. Using Section 4's criteria, write the two conditions it would need to satisfy, and explain why "peers with every other Tier-1" is a self-reinforcing, not externally enforced, requirement.
5. Explain why "eyeball network" and "content network" are useful terms even though they aren't part of the formal Tier-1/2/3 vocabulary, tying your answer to the peering-ratio dispute pattern from Chapter 51.
6. Using Section 10's full trace, explain what would change about the packet's path if RegionalISP's peering session with GlobalContent at DE-CIX went down entirely, referencing the BGP best-path algorithm from Chapter 49.

### Hard

7. Using the Go code in Section 14, extend the `classify` function to detect a fourth, informal category some engineers call "Tier-1.5" — a network that peers extensively with Tier-1s and Tier-2s but still buys a small amount of transit purely as a redundancy safety net. Decide on a reasonable rule and justify it.
8. Research one real, historical case of a network's Tier-1 status being publicly disputed (or a real depeering incident) and summarize what triggered it and how it was resolved, tying it back to Section 4's "continuously re-litigated" framing.
9. Argue, using Section 3's marketplace framing and Section 7's content-network category, whether the classic Tier-1/2/3 model will still be a useful description of Internet topology twenty years from now, or whether it is likely to be replaced by a different taxonomy centered on hyperscaler private backbones.

---

## 20. Summary, and the Bridge to Chapter 125

| Term | Meaning |
|---|---|
| Emergent hierarchy | The Tier-1/2/3 pyramid, produced by independent decisions with no central designer |
| Tier-1 club | Networks buying zero transit, sustained by mutual peering necessity, not certification |
| Eyeball network | A network whose defining traffic is inbound, to human end users |
| Content network | A network whose defining traffic is outbound, from its own servers |
| Peering LAN | The shared IP subnet every IXP member gets one address on |
| Route server | Shared IXP infrastructure that multilateral-peers many members through one BGP session |
| ICANN / IANA / RIRs / IETF | Bodies coordinating names, numbers, and protocols — never routing or business decisions |

This chapter showed *who* the Internet's networks are and *how* they physically and commercially connect. It didn't yet explain how one single company can make its servers appear to be "everywhere at once" from any of those networks' point of view — which is exactly the trick Anycast and modern CDN architecture perform, and which Chapter 125 now explains at the same full global scale this chapter just established.
