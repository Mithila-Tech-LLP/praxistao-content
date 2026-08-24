# Chapter 04: LAN, WAN, MAN, and PAN — Networks at Different Scales

> *"A network stretched across one desk and a network stretched across a continent obey the same rules of physics — but distance changes everything about what's practical, who owns what, and how fast is fast."*

---

## Table of Contents

1. [The Question: Does Distance Actually Change Anything?](#1-the-question-does-distance-actually-change-anything)
2. [A Thought Experiment: Stretching the Same Network](#2-a-thought-experiment-stretching-the-same-network)
3. [What Actually Changes With Distance](#3-what-actually-changes-with-distance)
4. [LAN — Local Area Network](#4-lan--local-area-network)
5. [WAN — Wide Area Network](#5-wan--wide-area-network)
6. [MAN — Metropolitan Area Network](#6-man--metropolitan-area-network)
7. [PAN — Personal Area Network](#7-pan--personal-area-network)
8. [Comparing All Four at a Glance](#8-comparing-all-four-at-a-glance)
9. [The Internet Is a WAN Built From Interconnected LANs](#9-the-internet-is-a-wan-built-from-interconnected-lans)
10. [The Physics of Distance, Quantified](#10-the-physics-of-distance-quantified)
11. [Real-World Bandwidth by Category](#11-real-world-bandwidth-by-category)
12. [Production Notes: How Real Organizations Mix These Categories](#12-production-notes-how-real-organizations-mix-these-categories)
13. [Hands-On Experiment: Measure Your Own Latency](#13-hands-on-experiment-measure-your-own-latency)
14. [Common Misconceptions](#14-common-misconceptions)
15. [Connections Backward and Forward](#15-connections-backward-and-forward)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)

---

## 1. The Question: Does Distance Actually Change Anything?

Chapter 3 defined a network as a set of computers connected by shared, addressable links. Nothing in that definition mentioned distance. By that definition, four computers on the same desk and four computers scattered across four different continents are equally valid networks, as long as they're joined by shared, addressable links.

And yet nobody who has ever used both a home Wi-Fi network and a video call to another country believes these two things "feel" the same. The home network is instant and nearly always available. The international call sometimes lags, occasionally drops, and clearly involves infrastructure your household doesn't own. Something is genuinely different — the question this chapter answers is: **what, precisely?**

---

## 2. A Thought Experiment: Stretching the Same Network

Imagine Chapter 3's bus-topology network — four computers, A, B, C, D, sharing one physical wire — built two different ways:

**Version 1:** All four computers sit on the same desk, connected by a 2-meter cable.

**Version 2:** The exact same four computers, but now A and B are in Mumbai, C is in London, and D is in São Paulo, connected by the "same" wire — just one that's now tens of thousands of kilometers long.

Logically, by Chapter 3's definition, these are the same network: a shared, addressable medium connecting four computers. Physically and practically, they are almost unrecognizably different systems. Version 1 can be built by one person in an afternoon with a spool of cable bought at a hardware store, and it costs almost nothing to operate. Version 2 would require crossing multiple countries' borders and, in reality, undersea cables (Chapter 23) owned by companies and consortia that took years to build and cost hundreds of millions of dollars — no individual person or even most companies could build Version 2 themselves.

This thought experiment is the entire point of this chapter: **the logical definition of a network doesn't change with distance, but nearly everything about the engineering, economics, and physics of building one does.** LAN, WAN, MAN, and PAN are simply names for the recurring points along that spectrum of distance where the practical differences become large enough to matter.

---

## 3. What Actually Changes With Distance

Four concrete things change as a network's geographic span grows, and every one of them is worth naming precisely, because each becomes the subject of later chapters.

### 1. Latency (an unavoidable consequence of Chapter 2's physics)

Chapter 2 established that no signal travels faster than light in its medium — about 200,000 km/s in typical optical fiber. This sets a hard floor on how quickly a signal can cross a given distance, no matter how good the equipment on either end is. A signal crossing a 2-meter desk cable takes about 10 nanoseconds to propagate — completely imperceptible. A signal crossing the roughly 14,000 km fiber-optic path between Mumbai and London (the real cable path is longer than the straight-line distance, since undersea cables follow actual geography) takes on the order of 70 milliseconds one-way, purely from the speed of light in fiber — before accounting for any processing delay at all. That's a difference of about seven million times, purely due to distance. This quantity — the delay caused purely by a signal's physical travel time — is called **propagation delay**, and Section 10 works through the exact math.

### 2. Ownership and administrative control

The 2-meter desk cable is owned, installed, and controlled entirely by whoever set up that desk. Nobody else's permission is needed to change it, extend it, or take it down. The undersea cable connecting continents is owned by a consortium of telecommunications companies and governments, is regulated by the laws of every country whose territory or waters it passes through, and cannot be modified by any single user of the network built on top of it. This is not a minor detail — it's the reason the word "Internet Service Provider" exists at all (previewed here, covered starting in Volume 2's history and formally throughout the rest of the course): almost nobody owns the long-distance links they use, they *lease access* to links someone else built and maintains.

### 3. Cost per unit of capacity

Running a cable across a room costs a few dollars. Running a cable across an ocean floor, that must survive water pressure, ship anchors, and marine life for decades, costs on a scale of many hundreds of millions of dollars for a single cable system. This isn't a difference in degree — it changes who can plausibly build the link at all, and it's why long-distance capacity is bought and sold ("leased") rather than simply built by whoever needs it, a business relationship covered properly in Chapter 51.

### 4. Reliability characteristics

A short local link fails rarely, and when it does, it's usually a simple, quickly-diagnosed physical problem (a bad cable, an unplugged connector). A network spanning a wide area is exposed to a much wider range of failure causes simultaneously — a backhoe cutting a buried cable in one city, a storm damaging an aerial line in another, a router misconfiguration at an intermediate point hundreds of kilometers from either endpoint. Wide-area networks are, as a direct consequence, engineered from the start assuming failures of individual links and paths are a routine, expected occurrence rather than a rare emergency — a design philosophy you'll see formalized throughout routing (Volume 7) and data center design (Volume 15).

These four consequences of distance — latency, ownership, cost, and reliability engineering — are the entire reason the categories in the rest of this chapter exist. They are not arbitrary size buckets; they're recognitions that a network's *design requirements* genuinely change as these four factors shift.

---

## 4. LAN — Local Area Network

### Definition

A **LAN (Local Area Network)** is a network confined to a single, limited physical location — typically one building or a small campus — that is owned, wired, and administered by a single organization or individual.

### Concrete examples

- Your home Wi-Fi network, connecting your laptop, phone, smart TV, and printer.
- The wired and wireless network inside a single office building, connecting employee desktops, printers, and internal servers.
- A university department's network, confined to one building.

### Characteristic properties

- **Latency:** typically well under 1 millisecond for the wired or wireless hop itself, since distances rarely exceed a few hundred meters.
- **Ownership:** the organization or individual using the network owns (or directly controls) essentially all of the physical infrastructure — cables, switches, access points.
- **Cost:** relatively low; consumer and small-business networking equipment is inexpensive and widely available.
- **Speed:** LANs are typically the fastest tier of network available to ordinary users — modern wired LANs commonly run at 1 Gbps or 10 Gbps, because the short distances and full ownership make high-bandwidth equipment practical and affordable.
- **Who builds it:** LANs are the only category in this chapter that an individual or small business can realistically build and fully control themselves, start to finish.

---

## 5. WAN — Wide Area Network

### Definition

A **WAN (Wide Area Network)** is a network that spans a large geographic area — often multiple cities, countries, or continents — typically connecting multiple separate LANs together, and typically relying on infrastructure leased from third-party carriers rather than owned outright by the organization using it.

### Concrete examples

- A company with offices in New York, London, and Singapore, with a private network connecting all three office LANs together.
- The network operated by an Internet Service Provider (ISP), connecting a whole city's or country's worth of homes and businesses back to the wider Internet.
- **The Internet itself** — which, as Section 9 explains precisely, is fundamentally a (very large) WAN.

### Characteristic properties

- **Latency:** ranges from single-digit milliseconds (within one country) to well over 100 milliseconds for the longest intercontinental paths, dominated by the propagation delay from Section 3.
- **Ownership:** almost never fully owned by the organization using it. Long-distance links are leased from telecommunications carriers who built and maintain the underlying fiber, satellite, or microwave infrastructure.
- **Cost:** significantly higher per unit of bandwidth than a LAN, because the underlying infrastructure (undersea cables, long-haul fiber routes, satellite links) is enormously expensive to build and maintain, and that cost is recovered by leasing capacity to many customers.
- **Speed:** historically much lower than LAN speeds for any individual connection, though core, shared backbone links (which carry many customers' traffic simultaneously, per Chapter 9's statistical multiplexing) run at extremely high aggregate capacities.
- **Who builds it:** almost always a specialized organization — a telecommunications carrier, an ISP, or a consortium of companies for shared infrastructure like undersea cables (Chapter 23) — rather than the end organizations that use the resulting network.

### WAN link technologies, briefly

It's worth naming a few real technologies specifically built for WAN-scale distances, because they look quite different from the Ethernet and Wi-Fi that dominate LANs, and the difference itself teaches something about this chapter's core theme (distance changes engineering requirements):

- **Leased lines and MPLS.** Historically, and still commonly today, companies pay a carrier for a dedicated, private circuit between two sites, using a technology called MPLS (Multiprotocol Label Switching) to route traffic predictably and privately across the carrier's shared backbone without exposing it to the public Internet. This trades a higher price for guaranteed, consistent performance — valuable for latency-sensitive applications like voice calls between offices.
- **SONET/SDH.** Much of the world's long-haul fiber backbone historically ran (and in many places still runs) on SONET/SDH, a family of standards specifically designed for extremely reliable, synchronized, long-distance optical transmission, with automatic failover to a backup path in milliseconds if a fiber is cut — a direct engineering response to Section 3's observation that wide-area links face a much broader range of failure causes than local ones.
- **Satellite links.** Where laying fiber isn't economical (remote regions, ships at sea, disaster recovery scenarios), a WAN link may instead bounce a signal off a satellite. Geostationary satellites orbit at about 35,786 km above the equator, meaning a single round trip through one adds roughly 240 milliseconds of unavoidable propagation delay (Chapter 4's own physics, applied vertically instead of horizontally) — which is why newer low-Earth-orbit satellite systems (covered as an emerging technology in Chapter 129) fly much closer to Earth specifically to cut this delay down.

None of these technologies exist in typical home or office LANs, precisely because a LAN's short distances and single-owner simplicity don't need them — Ethernet and Wi-Fi are already fast, cheap, and reliable enough at LAN scale. WAN technology exists specifically to solve problems (reliability over huge distances, guaranteed performance across shared carrier infrastructure, connectivity where no cable can reach) that simply don't arise at LAN scale.

---

## 6. MAN — Metropolitan Area Network

### Definition

A **MAN (Metropolitan Area Network)** sits between a LAN and a WAN: a network spanning a single city or metropolitan region — larger than any one building's LAN, but smaller and more geographically constrained than a WAN connecting distant cities or countries.

### Concrete examples

- A city government's network connecting multiple municipal buildings, traffic-signal control systems, and public libraries across the city.
- A university's network connecting several separate campuses located throughout the same city.
- A cable television/internet provider's local distribution network within one metropolitan area, before traffic leaves the city toward the wider Internet.

### Characteristic properties

- **Latency:** typically a few milliseconds — noticeably higher than a single-building LAN, but far lower than a cross-country or intercontinental WAN link.
- **Ownership:** often a mix — a municipal or institutional MAN might be fully owned by its operating organization (a city or university), whereas a commercial provider's metropolitan network is typically owned by that provider and used to serve many separate customers.
- **Why it's a distinct category at all:** MAN is the least sharply defined of the four terms in modern usage, and it's worth being honest about that: many of the practical distinctions that matter (ownership, latency, cost) shift gradually rather than at a clean boundary between "LAN," "MAN," and "WAN." The term remains useful mainly for describing infrastructure explicitly built and marketed at city scale — municipal fiber projects, city-wide public Wi-Fi initiatives, and metro-area carrier networks — more than as a precise engineering category with unique technical requirements of its own.

---

## 7. PAN — Personal Area Network

### Definition

A **PAN (Personal Area Network)** is the smallest category: a network confined to the immediate space around a single person, typically within a few meters, usually connecting personal devices to each other rather than to any wider infrastructure directly.

### Concrete examples

- A smartphone connected to a wireless earbud or smartwatch via Bluetooth.
- A laptop connected to a wireless mouse or keyboard.
- A fitness tracker syncing its data to a phone app over a short-range wireless link.

### Characteristic properties

- **Latency:** extremely low, since distances are at most a few meters.
- **Ownership:** entirely personal — a single individual typically owns every device involved.
- **Range:** deliberately short (often under 10 meters for Bluetooth), which is a feature, not a limitation — a PAN generally doesn't need or want to reach devices outside arm's reach, and a short range also reduces interference with other nearby PANs (your headphones shouldn't accidentally connect to a stranger's phone on the same train).
- **Relationship to LANs:** a PAN is often a "last few meters" extension that ultimately connects, through one of its devices, into a LAN or the wider Internet — your smartwatch talks to your phone over a PAN link (Bluetooth), and your phone separately talks to the Internet over a LAN link (Wi-Fi) or a cellular WAN link.

---

## 8. Comparing All Four at a Glance

| Property | PAN | LAN | MAN | WAN |
|---|---|---|---|---|
| Typical span | A few meters | One building/campus | One city/metro area | Countries/continents |
| Typical latency | Under 1 ms | Under 1–2 ms | A few ms | 10s–100s of ms |
| Typical ownership | One individual | One organization | Mixed (institution or carrier) | Leased from carriers/ISPs |
| Typical cost to build | Minimal (consumer devices) | Low–moderate | Moderate–high | Very high |
| Example technology | Bluetooth | Wired Ethernet, Wi-Fi | Metro fiber, city Wi-Fi | Undersea cable, long-haul fiber, satellite |
| Who can realistically build it | Anyone, instantly | Individual/small business | Institution or ISP | Telecom carrier/consortium |

---

## 9. The Internet Is a WAN Built From Interconnected LANs

Here is the idea this whole chapter has been building toward, and it's worth stating explicitly because it's easy to memorize the four category names without ever connecting them back to the network you actually use every day: **your home network is a LAN. Your Internet Service Provider connects that LAN, via WAN links, to millions of other LANs (homes, offices, data centers) all over the world. The Internet, as a whole structure, is the largest WAN in existence — and it is built almost entirely out of smaller LANs joined together by wide-area links.**

```
   HOME LAN                      ISP'S WAN INFRASTRUCTURE                DATA CENTER LAN
   (your devices,                (long-haul fiber, undersea               (a company's
    Wi-Fi router)                 cables, carrier equipment)               servers)

  [laptop]--+                                                          +--[server 1]
            |                                                          |
  [phone]---+---[router]=====long-distance links=====[router]---------+--[server 2]
            |    (LAN ends            (this is the                    |
  [TV]------+     here)                 "WAN" part)                   +--[server 3]
                                                                       (LAN ends here)
```

This is, precisely, what Chapter 6's title ("A Network of Networks") means, and this chapter has been the necessary setup for it: the Internet isn't one network at one scale — it's a nested structure of small, locally-owned LANs, connected through metropolitan and long-distance WAN infrastructure owned by many different organizations, none of which owns the whole thing. Chapter 6 draws this picture out fully; this chapter has given you the vocabulary (LAN, WAN, and the distance-driven reasons they differ) to understand why that structure exists in the first place rather than one giant flat network.

---

## 10. The Physics of Distance, Quantified

Section 3 asserted that propagation delay scales directly with distance. Let's compute real numbers, using the approximate speed of light in optical fiber (≈ 200,000 km/s, slower than the ≈300,000 km/s vacuum speed due to the fiber's refractive index, as introduced in Chapter 2 and covered fully in Chapter 22).

```go
package main

import "fmt"

const speedOfLightInFiberKmPerSec = 200000.0 // approximate, real fiber

func oneWayPropagationDelayMs(distanceKm float64) float64 {
	seconds := distanceKm / speedOfLightInFiberKmPerSec
	return seconds * 1000.0 // convert to milliseconds
}

func main() {
	routes := []struct {
		name       string
		distanceKm float64
	}{
		{"Across a desk (2 meters)", 0.002},
		{"Across an office building (200 meters)", 0.2},
		{"Across a city (30 km)", 30},
		{"Mumbai to Delhi (approx. 1,400 km)", 1400},
		{"New York to London (approx. 5,600 km, real cable route)", 5600},
		{"Mumbai to New York (approx. 13,000 km, real cable route)", 13000},
	}

	for _, r := range routes {
		delay := oneWayPropagationDelayMs(r.distanceKm)
		fmt.Printf("%-55s -> %10.6f ms one-way (propagation only)\n", r.name, delay)
	}
}
```

Output:

```
Across a desk (2 meters)                               ->   0.000010 ms one-way (propagation only)
Across an office building (200 meters)                 ->   0.001000 ms one-way (propagation only)
Across a city (30 km)                                  ->   0.150000 ms one-way (propagation only)
Mumbai to Delhi (approx. 1,400 km)                      ->   7.000000 ms one-way (propagation only)
New York to London (approx. 5,600 km, real cable route) ->  28.000000 ms one-way (propagation only)
Mumbai to New York (approx. 13,000 km, real cable route)-> 65.000000 ms one-way (propagation only)
```

Two things worth noticing. First, these are **one-way, propagation-only** numbers — a real round trip (like a web request and its response) roughly doubles this figure, and real measured latency is always somewhat higher still, because of processing delay at every router along the path (a subject this course returns to properly once routers are introduced in Chapter 44). Second, notice how the LAN-scale numbers (desk, office building) round to numbers so small they're irrelevant to human perception, while the WAN-scale numbers (New York–London, Mumbai–New York) are large enough to be clearly noticeable in a video call or a fast-paced online game. This single table is the quantitative heart of why LAN and WAN are treated as such different engineering problems, not just different sizes of the same problem.

---

## 11. Real-World Bandwidth by Category

Latency isn't the only number that shifts dramatically across these categories — so does typical bandwidth (how much data per second a link can carry), and it's worth seeing real figures side by side, because the direction of the difference often surprises beginners: unlike latency, which gets *worse* as you scale up from PAN to WAN, individual-link bandwidth generally shrinks as you scale up, even though the aggregate capacity of WAN backbones is enormous.

| Category | Typical individual-link bandwidth | Real example |
|---|---|---|
| PAN | 1–3 Mbps (Bluetooth Classic), up to ~50 Mbps (newer Bluetooth/UWB variants) | Wireless earbuds streaming audio |
| LAN | 1 Gbps–10 Gbps (wired), up to ~9.6 Gbps theoretical (Wi-Fi 6/6E) | An office Ethernet or Wi-Fi connection |
| MAN | 1–100 Gbps per metro fiber link | A city's municipal fiber ring |
| WAN (last-mile) | A few Mbps (older DSL) to ~1–2 Gbps (fiber to the home) | A home internet connection |
| WAN (core/backbone) | Individual long-haul fiber pairs now carry many terabits per second using advanced modulation (Chapter 16) | Undersea cables and carrier backbones |

The apparent contradiction — "WAN last-mile bandwidth is often lower than LAN bandwidth, yet WAN backbones carry far more data than any single LAN ever will" — resolves once you remember Chapter 9's statistical multiplexing: a single home's WAN connection only needs to carry that one household's traffic, but a carrier's backbone link is shared, simultaneously, across many thousands of customers' traffic at once, so it's engineered for enormously higher aggregate capacity even though no single conversation on it needs anywhere near that much. This is a first, concrete glimpse of an idea Volume 6 and 7 return to constantly: the capacity a network *needs* depends on how many independent things are sharing it, not just on the distance the signal travels.

---

## 12. Production Notes: How Real Organizations Mix These Categories

Real organizations rarely operate in just one of these four categories — they typically build a layered mix, and it's useful to see one concrete, realistic example end to end. Consider a mid-sized company with three offices:

- **Inside each office:** a LAN, owned and operated entirely by the company's own IT staff, connecting desks, meeting rooms, and printers via switches and Wi-Fi access points.
- **Between the three offices:** historically, companies leased dedicated WAN circuits (called MPLS lines) directly from a telecom carrier — expensive, but predictable and private. Increasingly, companies instead use **SD-WAN (Software-Defined WAN)**, which builds a private, centrally-managed overlay network on top of ordinary, cheaper Internet connections at each site, using encryption and smart routing to approximate the reliability of a dedicated leased line at a fraction of the cost — a real, modern example of software substituting for expensive dedicated infrastructure, made possible by the layering principle Chapter 5 previews and Chapter 24 formalizes.
- **Between the company and its cloud provider:** many organizations pay for a dedicated, private WAN-scale connection directly into a cloud provider's data center (e.g., AWS Direct Connect, Azure ExpressRoute) rather than routing that traffic over the shared public Internet, trading cost for more predictable latency and bandwidth — a real, production instance of choosing a private WAN link over a public one for specific traffic that needs it.
- **For individual mobile employees:** a personal PAN (their phone's Bluetooth connection to a headset) sits at one end, connected through their phone's cellular WAN link back to company systems, often through a VPN (Chapter 85) that makes their remote connection behave, logically, like part of the office LAN, regardless of physical distance.

The lesson worth taking from this: **LAN, MAN, WAN, and PAN are not competing choices an organization picks once — they are different tools an organization typically uses simultaneously, at the specific layer of its infrastructure where each one's trade-offs (Section 3's latency, ownership, cost, reliability) make sense.**

---

## 13. Hands-On Experiment: Measure Your Own Latency

Most operating systems include a `ping` command (Chapter 54 covers exactly what it does and how, using ICMP — for now, treat it as a stopwatch that measures round-trip time to another address). Run these three pings and compare the results:

```
# 1. Ping your own router (a LAN-scale hop)
ping 192.168.1.1          (use your router's actual local address, often
                            192.168.0.1 or 192.168.1.1 — check your device's
                            Wi-Fi/network settings if unsure)

# 2. Ping a server likely to be in your own country or continent (a shorter
#    WAN hop)
ping 1.1.1.1               (Cloudflare's public DNS resolver, widely
                             deployed with servers close to most users)

# 3. Ping a server you know or suspect is geographically distant from you
ping <a server address in a distant country, if you know one>
```

You should see round-trip times roughly in these bands: under 5 milliseconds for the router (LAN scale — matching Section 10's tiny numbers, doubled for a round trip and rounded up for real-world processing overhead), typically 10–40 milliseconds for a nearby, well-connected WAN destination, and potentially 150–300 milliseconds for a genuinely distant destination on the other side of the world. **This is Section 10's math, made real and measurable on your own machine, using a real diagnostic tool you'll learn properly in Chapter 54.**

---

## 14. Common Misconceptions

- **"WAN just means 'a bigger LAN.'"** As Section 3 details, the difference is not merely size — it's a qualitative shift in ownership (leased vs. owned infrastructure), cost structure, and reliability engineering, all driven by the unavoidable physics of distance. A LAN scaled up to cover more desks is still a LAN; a WAN exists specifically because certain distances cannot be crossed with infrastructure any single LAN-owning organization can build itself.
- **"MAN is a precise, universally agreed category with a fixed size range."** As Section 6 admits honestly, MAN is the fuzziest of the four terms in real-world usage, more a description of city-scale infrastructure projects than a category with sharp technical boundaries the way LAN and WAN's ownership and cost differences provide.
- **"The Internet is one big WAN, full stop — LANs aren't really part of it."** As Section 9 shows, this gets the relationship backward: the Internet's WAN links exist specifically to connect enormous numbers of separately-owned LANs together. Without the LANs at the edges (homes, offices, data centers), the WAN infrastructure would have nothing to connect — the Internet is best understood as LANs plus the WAN infrastructure joining them, not WAN instead of LANs.
- **"Higher latency always means something is wrong."** As Section 10's table shows, a large chunk of real-world latency to a distant server is simply the unavoidable, correctly-functioning result of the speed of light acting over a large distance — not a sign of a fault. Diagnosing "is this latency normal for this distance, or is something actually broken" is a real, recurring skill covered properly in Volume 18 (Observability & Debugging).
- **"A WAN link always has lower bandwidth than a LAN link."** As Section 11 shows, this is true for typical *individual, last-mile* WAN connections (like a home internet plan) compared to a LAN, but it inverts completely at the WAN's core: backbone fiber routes carry vastly more aggregate data than any single LAN, precisely because they're shared across enormous numbers of simultaneous users via statistical multiplexing.
- **"Satellite internet is just a slower version of fiber, otherwise the same."** As Section 5's WAN link technologies discussion shows, a geostationary satellite link adds roughly 240 milliseconds of unavoidable round-trip propagation delay purely from the physical distance to orbit and back — a fundamentally different latency profile from fiber, not merely a lower-bandwidth version of the same thing. This is precisely why latency-sensitive uses (video calls, competitive gaming) suffer disproportionately on geostationary satellite links even when the advertised bandwidth looks reasonable, and why low-Earth-orbit satellite systems (Chapter 129) exist specifically to address this.

---

## 15. Connections Backward and Forward

This chapter took Chapter 3's scale-agnostic definition of "network" and asked what changes as the same logical structure is stretched across greater and greater distances — arriving at latency, ownership, cost, and reliability as the four concrete forces that motivate distinguishing PAN, LAN, MAN, and WAN, and grounding all four in real bandwidth figures and a realistic production example of how organizations actually mix them. Chapter 5 uses this vocabulary immediately: it will show that "the Internet" (a WAN), "the Web" (an application running on top of it), and an "intranet" (a private network using the same technology, often built as a LAN or MAN) are three genuinely different things that get confused precisely because people haven't yet separated "the network" from "what runs on the network" — the distinction this chapter's LAN/WAN framing sets up perfectly to make.

---

## 16. Interview Questions & Model Answers

**Q1 (Beginner): What are the four main network categories introduced by geographic scale, and what does each stand for?**

*Model answer:* PAN (Personal Area Network) — devices within a few meters of one person, e.g., Bluetooth earbuds. LAN (Local Area Network) — a single building or campus, e.g., a home or office network. MAN (Metropolitan Area Network) — a single city or metro region, e.g., a municipal fiber network. WAN (Wide Area Network) — spanning countries or continents, e.g., an ISP's backbone or a multinational company's inter-office network. The Internet as a whole is the largest example of a WAN.

**Q2 (Beginner): Why does a network's geographic span affect its latency, beyond just "farther means slower"?**

*Model answer:* No signal can travel faster than the speed of light in its medium (about 200,000 km/s in typical optical fiber). This creates a hard, physics-based minimum time — propagation delay — proportional to distance, that no amount of better engineering can eliminate. A LAN's short distances make this delay negligible (fractions of a millisecond), while a WAN's long distances make it a significant, measurable component of total latency (tens to hundreds of milliseconds for intercontinental links).

**Q3 (Intermediate): Why do organizations typically own their LAN infrastructure outright but lease WAN infrastructure from carriers instead of building it themselves?**

*Model answer:* LAN infrastructure (cables, switches, access points within one building) is inexpensive and straightforward to install and maintain, making outright ownership practical for almost any organization. WAN infrastructure — long-haul fiber routes, undersea cables, satellite links — costs orders of magnitude more to build, often requires crossing multiple legal jurisdictions, and only makes economic sense when its cost is shared across many customers over a long operating lifetime. As a result, specialized telecommunications carriers build and maintain this infrastructure, and organizations lease capacity on it (through Internet Service Providers or dedicated WAN links) rather than building their own long-distance links.

**Q4 (Intermediate): A company has offices in three different cities, each with its own LAN. If those three LANs are connected together into one private company network, what category best describes the connecting links, and why?**

*Model answer:* The links connecting the three geographically separate office LANs together constitute a WAN — specifically, in this case, a private WAN, since the company controls how the LANs are connected together administratively, even though it almost certainly leases the actual long-distance transmission capacity from a telecommunications carrier rather than owning the physical long-haul links itself. Each individual office's internal network remains a LAN; the WAN is specifically the wide-area links joining the separate LANs into one company-wide network.

**Q5 (Advanced): Explain why MAN is described in this chapter as the least sharply defined of the four categories, and what practical purpose the term still serves despite that fuzziness.**

*Model answer:* The distinctions that clearly separate LAN from WAN — dramatically different latency, an ownership shift from full ownership to leased carrier infrastructure, and an order-of-magnitude jump in cost — shift gradually rather than sharply as scale increases from a single building to an entire city. A large university campus network and a small city's municipal network can have very similar technical characteristics, making "MAN" more a descriptive label for city-scale infrastructure projects (municipal fiber, metro-area carrier networks) than a category with unique engineering requirements the way LAN and WAN each have. It remains a useful term for describing infrastructure explicitly built and marketed at metropolitan scale, even without a crisp technical boundary separating it from its neighbors.

**Q6 (Advanced): Why might a modern company choose SD-WAN over a traditional dedicated leased line (MPLS) to connect its offices, and what trade-off is it making?**

*Model answer:* SD-WAN builds a private, centrally-managed, often encrypted overlay network on top of ordinary, comparatively inexpensive Internet connections at each site, rather than paying for dedicated point-to-point WAN circuits from a carrier. This can dramatically reduce cost and increase flexibility (new sites can be added quickly using whatever local Internet service is available), at the trade-off of relying on the shared public Internet's variable performance underneath the overlay, rather than a carrier's dedicated, more predictable circuit — a real, practical illustration of Chapter 3's latency/cost/reliability trade-offs playing out in a current industry decision, made possible by the layering principle that lets a private network be built logically on top of shared physical infrastructure.

**Q7 (Advanced): Why does a geostationary satellite link add a fixed, unavoidable amount of latency regardless of how good the satellite equipment is, and why doesn't this same penalty apply to fiber-based WAN links of similar geographic reach?**

*Model answer:* A geostationary satellite orbits at a fixed altitude of about 35,786 km above the equator, so any signal bounced off it must physically travel up to the satellite and back down — a round trip of over 70,000 km purely for one hop — at the speed of light, adding roughly 240 milliseconds that no equipment improvement can remove, since it's a direct consequence of distance and the physical speed limit established in Chapter 2. A fiber link covering a similar geographic reach (say, connecting two cities 3,000 km apart) follows a much shorter physical path along the Earth's surface rather than up to geostationary orbit and back, so its propagation delay, while still nonzero, is a small fraction of the satellite link's — the penalty is specifically about the physical path length involved, not about the general difficulty of long-distance communication.

---

## 17. Exercises

### Easy

1. Classify each of the following as PAN, LAN, MAN, or WAN: (a) a smartwatch talking to a phone via Bluetooth, (b) a company's network connecting its Tokyo and Berlin offices, (c) all the computers in a single school building, (d) a city's public library system's shared network across all its branches.
2. Using Section 10's table as a reference, explain in your own words why a video call between two people in the same city feels different from a video call between people on opposite sides of the world, even with identical software and equipment.
3. Run the three-step `ping` experiment from Section 13 on your own network and record the three round-trip times you observe.
4. Using Section 5's WAN link technologies, explain in one or two sentences why a company's video-conferencing traffic between two offices might be routed over a leased MPLS line instead of the ordinary public Internet.

### Medium

1. Modify the Go program in Section 10 to add two more example routes of your choosing (look up or estimate two real-world city-pair distances), and compute their one-way propagation delay.
2. Explain, using Section 3's four factors (latency, ownership, cost, reliability), why a company would choose to lease WAN connectivity from an ISP rather than attempt to lay its own long-distance fiber between two offices in different countries.
3. A friend claims "a MAN is just a LAN that covers a whole city instead of one building, so there's no real difference." Using Section 6's discussion, write a short rebuttal or agreement, and justify your position.
4. Using Section 11's bandwidth table, explain why a home's individual WAN connection can have lower bandwidth than its own LAN, while the Internet's WAN backbone as a whole carries far more traffic than any single LAN.

### Hard

1. Using real published distances (research these, or estimate using a map), compute the theoretical minimum one-way propagation delay for a signal traveling by fiber between two cities of your choosing on different continents. Compare your computed value to a real, publicly reported latency figure between those two cities (many websites publish real-world "ping time" tables between major cities) and explain the difference between your theoretical minimum and the real observed figure.
2. This chapter asserts that WAN links are "leased" rather than owned by the organizations using them. Research what an Internet Exchange Point (IXP) is, in one or two sentences (this will be covered fully in Chapter 51), and explain how it complicates the simple picture of "every organization leases WAN capacity from a carrier."
3. Propose a scenario where a network spans a large geographic distance (qualifying it as WAN-scale by Section 3's latency/distance reasoning) but is nonetheless fully owned by a single organization, without any leased third-party infrastructure. Is this plausible in the real world? What kind of organization might actually do this, and why?
4. Section 5 describes SONET/SDH's ability to automatically fail over to a backup path within milliseconds of a fiber cut. Using Section 3's four factors (latency, ownership, cost, reliability), explain why this kind of automatic failover is a much higher engineering priority for WAN infrastructure than it typically is for a single office's LAN cabling.

---

## 18. Summary

| Term | Meaning |
|---|---|
| Propagation delay | The unavoidable time a signal takes to physically travel a distance, bounded by the speed of light in the medium |
| PAN (Personal Area Network) | A network confined to a few meters around one person (e.g., Bluetooth) |
| LAN (Local Area Network) | A network confined to one building or campus, typically fully owned by its user |
| MAN (Metropolitan Area Network) | A network spanning a city or metro region; the least sharply defined of the four categories |
| WAN (Wide Area Network) | A network spanning countries or continents, typically built on leased carrier infrastructure |
| Internet (previewed) | The largest WAN in existence, built by interconnecting enormous numbers of separately-owned LANs |
| Ownership boundary | The dividing line between infrastructure an organization controls directly and infrastructure it leases from a carrier |
| Statistical multiplexing (preview) | Why WAN backbone links carry far more aggregate bandwidth than any single LAN, by sharing capacity across many simultaneous users |
| SD-WAN | A modern approach building a private, managed WAN overlay on top of ordinary Internet connections instead of dedicated leased circuits |

This chapter showed that distance changes latency, ownership, cost, and reliability engineering — not just the label attached to a network — and that the Internet is best understood as LANs joined by WAN infrastructure, not as a single flat structure. Chapter 5 uses this same vocabulary to untangle three terms people constantly confuse: the Internet, the Web, and an intranet.
