# Chapter 126: Submarine Cables — Ownership, Laying, Repair, and Real Outages

> *"Chapter 23 told you the internet runs through glass on the ocean floor, not satellites. This chapter tells you who actually owns that glass, who sails out to fix it when a fishing trawler drags it in half, and why that repair can take three weeks — not three hours."*

---

## Table of Contents

1. [What Chapter 23 Established, and What's Missing](#1-what-chapter-23-established-and-whats-missing)
2. [The Traditional Ownership Model: Telecom Consortia](#2-the-traditional-ownership-model-telecom-consortia)
3. [The Shift: Hyperscalers Building Their Own Cables](#3-the-shift-hyperscalers-building-their-own-cables)
4. [Why a Cloud Company Wants to Own a Cable](#4-why-a-cloud-company-wants-to-own-a-cable)
5. [A Real Table of Major Cables and Who Owns Them](#5-a-real-table-of-major-cables-and-who-owns-them)
6. [How a Cable Is Actually Manufactured and Laid](#6-how-a-cable-is-actually-manufactured-and-laid)
7. [Landing Stations: Where Ocean Meets Country](#7-landing-stations-where-ocean-meets-country)
8. [Cable Repair: The Ships, the Process, and Why It Takes Weeks](#8-cable-repair-the-ships-the-process-and-why-it-takes-weeks)
9. [Case Study: The 2006 Hengchun Earthquake](#9-case-study-the-2006-hengchun-earthquake)
10. [Case Study: The 2008 Mediterranean Cable Cuts](#10-case-study-the-2008-mediterranean-cable-cuts)
11. [Case Study: Tonga, 2022 — One Cable, One Nation](#11-case-study-tonga-2022--one-cable-one-nation)
12. [Case Study: The Red Sea Cuts and Baltic Sea Incidents](#12-case-study-the-red-sea-cuts-and-baltic-sea-incidents)
13. [Geopolitics: Protection Zones, Landing Rights, and Sabotage Concerns](#13-geopolitics-protection-zones-landing-rights-and-sabotage-concerns)
14. [Hands-On: Exploring the Real Cable Map](#14-hands-on-exploring-the-real-cable-map)
15. [Code: Modeling Redundancy Risk in Go](#15-code-modeling-redundancy-risk-in-go)
16. [Common Misconceptions](#16-common-misconceptions)
17. [Production Notes](#17-production-notes)
18. [What This Chapter Simplified](#18-what-this-chapter-simplified)
19. [Interview Questions & Model Answers](#19-interview-questions--model-answers)
20. [Exercises](#20-exercises)
21. [Summary, and the Bridge to Chapter 127](#21-summary-and-the-bridge-to-chapter-127)

---

## 1. What Chapter 23 Established, and What's Missing

Chapter 23 made one physically surprising claim and backed it with real numbers: well over 95% of intercontinental Internet traffic crosses oceans through fiber-optic cables lying on the seabed, not satellites. It explained the physics (single-mode fiber, Chapter 22, insulated and armored for underwater use), gave real route distances and latencies, and briefly touched on threats (fishing trawls, earthquakes, redundancy as the standard mitigation).

What it didn't cover, because it belonged in a physical-layer volume, not a global-infrastructure one: **who actually owns these cables today, what a cable-laying ship's job actually looks like day to day, what happens operationally in the weeks after a cable is cut, and what real, dated incidents reveal about how much a single point of physical failure can still matter** even in a supposedly redundant system. This chapter answers all four, in full operational detail.

---

## 2. The Traditional Ownership Model: Telecom Consortia

For most of submarine cable history — from the first transatlantic telegraph cable in 1858 through the fiber-optic era beginning in the late 1980s — the dominant ownership model has been the **consortium**: a group of telecommunications companies (often one from each country the cable will land in, plus international carriers) jointly funding construction, in exchange for each member receiving a guaranteed share of the cable's total capacity proportional to their investment.

**Why this model made sense for most of the industry's history:** a single transoceanic cable project can cost hundreds of millions of dollars and take years to plan, survey, and build. No single national telecom typically had a business case to fund an entire intercontinental cable alone, when its own domestic traffic needs were a fraction of the cable's total capacity. Pooling investment across many telecoms — each buying exactly the capacity share it needed — spread both the cost and the risk, and gave every participating carrier assured, contractually-guaranteed rights to a slice of the finished cable.

Real, long-running examples of this consortium model include the **SEA-ME-WE** cable family (SouthEast Asia-Middle East-Western Europe, now on its fifth generation, SEA-ME-WE 5/6, built and operated by a rotating consortium of dozens of national and international carriers) and **TAT-14** (Trans-Atlantic Telephone cable system, a consortium of major North American and European carriers). This model remains active today — it hasn't disappeared — but Section 3 covers the structural shift that's happened alongside it.

---

## 3. The Shift: Hyperscalers Building Their Own Cables

Starting roughly in the mid-2010s and accelerating sharply since, a fundamentally different class of owner has entered submarine cable construction: **the large cloud and content companies themselves — Google, Meta, Microsoft, and Amazon** — funding cables either entirely alone or as the dominant partner in a much smaller consortium than the traditional telecom model, sometimes retaining exclusive or majority use of the resulting capacity rather than sharing it broadly.

This is a genuine structural shift in the industry, not a minor variation on the consortium model, and it's directly connected to Chapter 124's observation that large content networks don't fit the classic Tier-1/2/3 taxonomy at all: exactly the same companies that build their own private terrestrial backbones (previewed there, covered fully in Chapter 127) have extended that same "build it ourselves" logic across the ocean floor.

---

## 4. Why a Cloud Company Wants to Own a Cable

The economic and operational logic, stated plainly:

- **Traffic volume justifies it.** Google, Meta, Microsoft, and Amazon each move an enormous, continuously growing volume of traffic between their own data centers on different continents (video, cloud storage replication, search index synchronization, and inter-region cloud backbone traffic previewed in Chapter 127) — traffic volumes large enough that owning dedicated capacity outright can be more cost-effective over the cable's decades-long operational life than perpetually leasing capacity from a traditional carrier consortium.
- **Control over capacity and upgrade timing.** A consortium member sharing a cable with dozens of other carriers has to negotiate collectively for capacity upgrades and route decisions. A company that owns (or owns the controlling majority of) a cable outright controls exactly when and how it upgrades the terminal equipment (Chapter 22, Section 9's DWDM electronics, the actual lever that grows a cable's usable capacity over its lifetime without touching the fiber itself).
- **Latency and reliability control for their own services.** Owning the physical path between two of your own data centers means not depending on a third party's business decisions, maintenance schedules, or pricing changes for a link your service's user-facing latency directly depends on.
- **It reduces reliance on the public Internet and public transit markets entirely** — directly extending Chapter 124's Section 7 observation that hyperscalers increasingly sidestep the traditional peering/transit marketplace, now literally at the level of owning the physical medium itself.

---

## 5. A Real Table of Major Cables and Who Owns Them

| Cable | Route | Owner(s) | Notable fact |
|---|---|---|---|
| **MAREA** | Virginia Beach, USA ↔ Bilbao, Spain | Microsoft, Meta, and Telxius (Telefónica's infrastructure arm) — a small joint venture, not a broad traditional consortium | ~6,600 km; designed for roughly 160 Tbps of capacity; completed 2017 |
| **Dunant** | Virginia Beach, USA ↔ Saint-Hilaire-de-Riez, France | Google (sole owner) | ~6,400 km; one of the first major cables built and owned entirely by a single cloud company; designed for roughly 250 Tbps; completed 2021 |
| **Grace Hopper** | New York, USA ↔ Bude, UK, and Bilbao, Spain | Google (sole owner) | Named after the pioneering computer scientist and US Navy rear admiral; completed 2022 |
| **Curie** | Los Angeles, USA ↔ Valparaíso, Chile | Google (sole owner) | One of the first Google-built cables to serve South America's west coast directly |
| **Equiano** | Portugal ↔ multiple landing points along the West African coast ↔ South Africa | Google (sole owner) | Explicitly built to bring higher-capacity, lower-cost connectivity to underserved West/Southern African markets |
| **2Africa** | Circling almost the entire African continent, plus landings in Europe and the Middle East | Meta, together with a consortium including China Mobile International, MTN GlobalConnect, Orange, Vodafone, and others | At roughly 45,000 km, one of the longest submarine cable systems ever built; a hybrid model — Meta-led but still a genuine multi-carrier consortium |
| **SEA-ME-WE 5 / 6** | Southeast Asia ↔ Middle East ↔ Western Europe | Traditional multi-carrier consortium (dozens of national and international telecoms) | Represents the long-running traditional consortium model (Section 2), still very much active alongside the newer hyperscaler-owned cables |
| **Southern Cross Cable Network** | Australia ↔ New Zealand ↔ USA (Pacific) | A dedicated cable operating company (Southern Cross Cables Limited), historically carrier-neutral rather than either a broad legacy consortium or a single hyperscaler | An example of a third ownership pattern: a specialized, independent cable operator selling capacity commercially to many customers |

The pattern across this table is worth naming directly: **the newest, most cutting-edge, highest-capacity cables increasingly skew toward single-hyperscaler or hyperscaler-led ownership**, while a large base of existing global capacity remains under the traditional multi-carrier consortium model from Section 2 — both models coexist, and neither has replaced the other, but the direction of new investment has shifted meaningfully since the mid-2010s.

---

## 6. How a Cable Is Actually Manufactured and Laid

Chapter 23, Section 8 gave the cross-sectional structure (armor layers, power conductor, fiber core). Here's the operational process, start to finish, that gets a finished cable from a factory onto the seabed:

1. **Route survey.** Before any cable is manufactured, a survey ship spends months mapping the exact intended seabed route in detail — depth, terrain, known fault lines, existing cables and pipelines to avoid, fishing zones, and shipping lanes — using sonar and remotely operated underwater vehicles. This survey directly determines the cable's final length, which is why real cable lengths (Chapter 23, Section 7's table) always exceed the straight-line great-circle distance between endpoints: the route deliberately winds around seabed hazards rather than taking the shortest possible path.
2. **Manufacturing.** The cable is manufactured in one continuous run — sometimes thousands of kilometers in a single unbroken production run — at a specialized cable factory operated by one of a small handful of companies with the capability to do this at all: **SubCom** (USA), **Alcatel Submarine Networks/ASN** (France, part of Nokia), and **NEC Corporation** (Japan) are the three companies that manufacture the overwhelming majority of the world's submarine cable systems.
3. **Loading onto a cable-laying ship.** The finished cable is coiled into massive circular tanks aboard a purpose-built **cable-laying ship** — vessels operated by the same small set of manufacturers (SubCom, ASN, NEC) or by specialized marine contractors like **Orange Marine**. Real, named vessels in this fleet include ships like **CS Durable**, **CS Cable Innovator**, and **Léon Thévenin** — a genuinely small global fleet, worth returning to in Section 8, because that same small fleet also handles repairs.
4. **Laying.** The ship follows the surveyed route (step 1) precisely, paying the cable out over the stern at a controlled rate matched to the ship's speed and the water depth. In shallow coastal waters (roughly the first several dozen kilometers from shore, where fishing trawls and ship anchors pose the greatest risk, per Chapter 23, Section 9), a **submersible plow**, towed behind the ship or operated by a remotely operated vehicle, simultaneously cuts a trench and buries the cable a meter or more below the seabed. In deep ocean, well below where trawling or anchors reach, the cable is simply laid directly on the seabed surface — burial isn't needed once you're deep enough that fishing gear and anchors physically can't reach it.
5. **Splicing multiple cable segments together.** A single cable system is often too long to load onto one ship in one continuous piece, or is built and shipped in segments; joining two segments at sea (or joining a repaired section, Section 8) requires precisely fusion-splicing every individual optical fiber (Chapter 22's fiber-splicing techniques, at ocean scale) inside a specialized housing that will then be sealed and lowered back to the seabed.
6. **Landing and testing.** Once the ship reaches shore, the cable's final stretch is pulled ashore (sometimes by a smaller specialized shore-landing vessel or even manually in shallow water) and connected at the **landing station** (Section 7), and the entire system is tested end-to-end for signal quality before being declared operational — a process that itself can take weeks after physical laying is complete.

```
   Survey ship (months)  ->  Factory manufacture (custom per project)
        │
        ▼
   Cable-laying ship departs, following the surveyed route
        │
        ├── Shallow coastal water: plow BURIES cable ~1m+ deep
        │
        ├── Deep ocean: cable simply LAID on the seabed surface
        │
        ▼
   Cable reaches the far shore -> landing station connection (Sec. 7)
        │
        ▼
   End-to-end testing -> system declared operational
```

---

## 7. Landing Stations: Where Ocean Meets Country

A **cable landing station** is the physical building, typically located within a few kilometers of the actual beach where the cable comes ashore, where the submarine cable's undersea electronics interface with the terrestrial network. This is where the cable's copper power conductor (Chapter 23, Section 8) connects to a shore-based power feeding equipment station that supplies the DC current powering every optical amplifier along the entire undersea route, and where the optical signal is handed off to terrestrial fiber networks connecting onward to inland data centers and Internet exchange points (Chapter 124).

Landing stations are, deliberately, unglamorous-looking buildings — often resembling a nondescript industrial facility rather than anything visibly identifiable as critical global infrastructure — but their locations are a matter of public record and a genuine subject of geopolitical attention (Section 13), because a country's cable landing stations are, quite literally, the small number of physical points where its entire international Internet connectivity crosses from sea to land. A single country can have anywhere from a handful to several dozen landing stations, and countries actively pursuing digital infrastructure investment (again, the West African Equiano example from Section 5) explicitly treat winning new cable landings as a strategic economic development goal, since a new landing station typically brings a substantial, immediate capacity and cost improvement to a region's Internet access.

**Landing rights are also a real, sometimes contentious, regulatory and geopolitical matter** — a cable operator needs formal government permission to bring a cable ashore in any given country, and that permission process has, in some well-documented recent cases, become entangled with broader national security reviews of which companies (and which countries' companies) are allowed to build and operate landing infrastructure at all.

---

## 8. Cable Repair: The Ships, the Process, and Why It Takes Weeks

This is the section that most directly explains the "measurably degrade a region's internet for weeks" half of this chapter's mandate. When a cable is damaged — by a trawler, an anchor, an earthquake, or deliberate action — restoring it is a genuinely slow, physically demanding process, and understanding why requires walking through it step by step.

**1. Detection and localization.** Cable operators continuously monitor their systems' optical signal characteristics from the landing stations at both ends; a break shows up immediately as a loss of signal, and the equipment can localize the fault to within a few kilometers using the same time-domain reflectometry principle used for terrestrial fiber fault-finding (an optical pulse's reflection timing reveals distance to the break) — but a few kilometers of *ocean*, at depths that can exceed several kilometers, is still an enormous area to search physically.

**2. Ship availability — the real bottleneck.** There is a genuinely small global fleet of cable repair ships — commonly cited industry figures put the worldwide dedicated repair fleet at roughly **50-60 vessels total**, covering every ocean on Earth, operated by a mix of the cable manufacturers themselves (SubCom, ASN, NEC) and dedicated maintenance consortiums that pool repair-ship access across many cable owners specifically because no single cable owner can justify keeping a repair ship permanently on standby near every cable they operate. If the nearest available repair ship is already committed to another repair job elsewhere in the ocean, or is in port for its own maintenance, a newly damaged cable simply has to wait its turn — this single fact is the single biggest, most common reason repairs stretch into weeks rather than days.

**3. Transit to the fault location.** Once a ship is available, it has to physically sail to the fault's location — which, for a fault in the middle of an ocean thousands of kilometers from the nearest port, can itself take several days.

**4. Grappling and retrieval.** In deep water, the ship uses a specialized grappling hook lowered on a long line to physically snag the cable on the seabed and haul the damaged section up to the surface — a slow, weather-dependent process, since the ship needs sufficiently calm seas to safely operate this equipment, and open-ocean weather windows aren't guaranteed on demand.

**5. Cutting, splicing, and testing.** Once the damaged section is aboard, the crew cuts out the broken portion, splices in a new length of cable (fusion-splicing every individual fiber strand, as in Section 6), tests the repaired splice's optical quality, and then carefully lowers the repaired cable back to the seabed — a delicate operation in its own right, since the repaired section needs to be laid without introducing new stress or damage.

**6. Permits and territorial waters.** If the fault lies within another country's territorial waters or exclusive economic zone, the repair ship typically needs formal government permission to operate there at all — a real, sometimes slow, diplomatic and regulatory step that has, in some documented cases, added meaningful delay to repairs entirely independent of the physical difficulty of the repair itself.

**Putting it together — why "weeks" is normal, not exceptional:** ship availability delay (days to weeks) + transit time (days) + weather-dependent grappling and splicing (days, sometimes delayed further by weather) + permitting (potentially days to weeks) routinely adds up to **two to four weeks being a completely normal total repair timeline** for a single deep-ocean cable fault, and multi-cable faults (Section 9) can take considerably longer if they exhaust the locally available repair-ship capacity.

---

## 9. Case Study: The 2006 Hengchun Earthquake

On December 26, 2006, a magnitude 7.1 earthquake struck near Hengchun, Taiwan, triggering underwater landslides that severed **multiple submarine cables simultaneously** in the Luzon Strait — a corridor carrying a large share of Asia's international cable traffic, including segments of major systems connecting Taiwan, mainland China, Hong Kong, and other regional hubs to the wider Internet.

The regional impact was severe and well-documented: significantly degraded or disrupted Internet and international telephone connectivity across large parts of East and Southeast Asia for **days to weeks**, as remaining, undamaged cables and satellite backup capacity absorbed traffic they weren't provisioned to carry alone, and the small global repair fleet (Section 8) had to prioritize and sequence multiple simultaneous faults in the same difficult, seismically active region. Full restoration of all affected systems took **on the order of weeks**, directly illustrating this chapter's core claim: even with Chapter 23's redundancy principle in effect (multiple cables *did* exist on this corridor), damaging enough of them simultaneously, in a geologically hazardous shared corridor, still produced a measurable, multi-day-to-multi-week regional degradation — redundancy reduces the *probability* and *severity* of a total outage, but a severe enough shared-risk event can still overwhelm it.

---

## 10. Case Study: The 2008 Mediterranean Cable Cuts

In late January 2008, multiple major cables — including segments of the **SEA-ME-WE 4** and **FLAG (FALCON/FLAG Europe-Asia)** systems — were cut in the Mediterranean Sea near Alexandria, Egypt, within days of each other. The disruption measurably degraded Internet connectivity across **Egypt, India, Pakistan, Saudi Arabia, and other countries** along this corridor, with some regions reportedly losing a large fraction of their international bandwidth for the duration of the outage.

The exact cause of the initial cuts was debated publicly at the time (early reporting suggested ship anchors, though this was never conclusively confirmed for every cut in the cluster), but the operational lesson is the same as Section 9's: this corridor, like the Luzon Strait, carries a disproportionate share of traffic for the region it serves, so even a handful of cable faults concentrated in one geographically narrow chokepoint can produce outsized, multi-country impact — a direct, real-world illustration of why cable route *diversity* (not just cable *count*) matters for genuine redundancy.

---

## 11. Case Study: Tonga, 2022 — One Cable, One Nation

In January 2022, the underwater eruption of the Hunga Tonga–Hunga Haʻapai volcano severed the **single submarine cable** connecting the Pacific island nation of Tonga to the rest of the world's Internet. Because Tonga depended on just one cable with no redundant path, the country was left almost entirely without international Internet connectivity for **roughly five weeks**, until repairs could be completed — repair complicated in this case by the same access and weather difficulties described generically in Section 8, compounded by the remoteness of the location and the aftermath of the volcanic event itself.

This is the cleanest, starkest real-world illustration in this chapter of Chapter 23's redundancy principle stated in the negative: **a single cable, with no alternate path, means a single point of failure for an entire nation's international connectivity** — exactly the scenario the redundancy principle exists to prevent, made concrete by a real country that, at the time, simply didn't have a second cable to fail over to.

---

## 12. Case Study: The Red Sea Cuts and Baltic Sea Incidents

More recently, cable faults in the **Red Sea** (affecting cables carrying significant Europe-Asia traffic, in a period of heightened regional maritime tension) and a series of **Baltic Sea cable and pipeline incidents** (involving both power and data cables, occurring amid heightened geopolitical tension in the region) have drawn substantial public and government attention specifically because of the *suspected* — though, in several cases, publicly unconfirmed or disputed — possibility of deliberate interference rather than routine accidental damage.

These recent incidents are included here deliberately as **current, actively-developing events** rather than settled historical case studies: at the time of writing, causation in several of these specific incidents remains publicly debated, under investigation, or attributed differently by different governments and outlets. The durable, non-speculative lesson that both older confirmed cases (Sections 9-11) and these newer, more geopolitically charged cases share is the same one this chapter keeps returning to: **the world's Internet still depends, physically, on a countable, mappable, and — as Section 13 discusses — increasingly securitized set of cables**, and that dependency is treated as a serious matter of national infrastructure policy by a growing number of governments.

---

## 12b. A Worked Repair Timeline, Day by Day

Section 8 described the repair process as a sequence of steps. To make "weeks, not days" concrete, here is a realistic, worked day-by-day timeline for a single mid-ocean cable fault, synthesizing publicly reported patterns from real incidents like those in Sections 9-12 (specific durations are illustrative composites, not one single documented case):

```
Day 0:       Fault detected instantly via loss-of-signal monitoring
             at both landing stations. Optical time-domain
             reflectometry narrows the fault to within a few km.

Day 0-3:     Cable owner notifies the maintenance consortium /
             repair-ship operator. Nearest available repair ship
             is identified -- if it's mid-repair on another fault
             elsewhere in the ocean, this step alone can take much
             longer than 3 days.

Day 3-4:     Permits requested for any territorial-water transit
             or repair-zone access needed along the ship's route
             and at the fault site itself.

Day 4-9:     Repair ship transits from its home port or previous
             job site to the fault location -- for a genuinely
             mid-ocean fault, several days of open-sea travel.

Day 9-12:    Grappling operations: the ship lowers a grapnel,
             locates and hooks the cable on the seabed, and hauls
             the damaged section to the surface. Delayed further
             by any bad weather window during this period.

Day 12-14:   Damaged section cut out; new cable segment spliced in,
             fiber by fiber, and optically tested.

Day 14-15:   Repaired cable carefully lowered back to the seabed
             along the original route; final end-to-end signal
             quality verification across the whole system.

TOTAL:       ~2 weeks in this composite scenario -- and this
             assumes a repair ship was available almost
             immediately and weather cooperated. Real incidents
             with ship queueing delays or difficult weather
             windows routinely stretch into 3-4+ weeks, consistent
             with Tonga's roughly five-week restoration in 2022
             (Section 11).
```

---

## 13. Geopolitics: Protection Zones, Landing Rights, and Sabotage Concerns

Given Sections 9-12's case studies, it's unsurprising that submarine cable security has become an explicit policy topic, not just a private industry operational concern:

- **Cable Protection Zones**: several countries have designated legally protected zones around cable landing approaches where anchoring, trawling, and certain other seabed activities are restricted or banned specifically to reduce accidental damage — a direct regulatory response to Section 8's "fishing trawls and anchors are the most common cause of damage" fact from Chapter 23.
- **Landing rights as a national security review matter**: as Section 7 noted, several governments now explicitly review, and have in some documented cases blocked or renegotiated, proposed cable landings on national security grounds tied to which companies (and which countries) would control the resulting infrastructure.
- **Naval and coast guard monitoring**: several governments have publicly increased naval patrol and monitoring attention specifically around known cable corridors and landing approaches in response to the kind of incidents described in Section 12.
- **The core vulnerability is structural, not merely a security lapse**: cable routes and landing station locations are, necessarily, public information (ships need to know where cables are to avoid damaging them accidentally, and cable operators need public coordination to prevent exactly that) — meaning the same information that supports safe, cooperative marine operations is also, unavoidably, information that would matter to anyone seeking to cause deliberate damage. This tension has no clean resolution and is an active area of ongoing policy discussion, not a solved problem.

---

## 13b. Real Numbers: Global Cable Infrastructure at a Glance

| Metric | Approximate figure (mid-2020s) |
|---|---|
| Active submarine cables worldwide | Over 600 systems (this figure grows year over year as new cables enter service) |
| Total combined length | More than 1.4-1.5 million km — enough to circle Earth roughly 35+ times |
| Share of intercontinental traffic carried | 95-99%, per Chapter 23 |
| Global dedicated repair-ship fleet | Roughly 50-60 vessels covering every ocean |
| Cable manufacturers capable of building new systems | Effectively three: SubCom, Alcatel Submarine Networks (ASN), NEC |
| Typical new major transoceanic cable cost | Commonly several hundred million US dollars, sometimes exceeding a billion for the largest, longest systems (e.g., 2Africa) |
| Typical time from planning to service | Multiple years — route survey, permitting, manufacturing, and laying together routinely span 2-4+ years for a large system |

These figures are worth holding next to Chapter 124's ISP/IXP numbers: both chapters describe a global-scale system with a surprisingly *small* number of physical chokepoints and specialized operators relative to the sheer size of "the Internet" as an abstraction — a small number of cables, a small number of manufacturers, a small number of repair ships, and a small number of IXP operators, all quietly underpinning something that feels, from a browser, like an infinite, borderless network.

---

## 14. Hands-On: Exploring the Real Cable Map

```bash
# 1. Explore the real, continuously updated global submarine cable
#    map (a widely used public industry reference):
#    https://www.submarinecablemap.com
#    Zoom into any region and click individual cables to see their
#    real owners, landing points, and length -- directly cross-check
#    a few entries against Section 5's table.

# 2. Look up a specific country's cable landing stations and count
#    how many independent cables serve it -- a small number (as in
#    Tonga's 2022 case, Section 11) is a direct, visible redundancy
#    risk signal.

# 3. Run traceroute/ping to a server in a distant region and compare
#    the measured RTT to Chapter 23 Section 7's theoretical
#    speed-of-light figures for the real cable route between your
#    location and that region -- the same experiment Chapter 23
#    introduced, now informed by knowing which real, named cable
#    likely carries your traffic.
traceroute www.some-distant-region-site.example
```

---

## 15. Code: Modeling Redundancy Risk in Go

A small graph-connectivity model showing why a country with only one cable (Section 11's Tonga scenario) is fundamentally more fragile than one with several, independent of any of this chapter's real geopolitics — pure graph theory doing the explanatory work:

```go
package main

import "fmt"

// Country represents a node connected to the rest of the world only
// through its list of submarine cables (edges to a shared "rest of
// world" node, simplified).
type Country struct {
	Name   string
	Cables []string // names of independent cables serving this country
}

// survivesLoss simulates losing one specific cable and reports
// whether the country retains ANY international connectivity.
func survivesLoss(c Country, lostCable string) bool {
	remaining := 0
	for _, cable := range c.Cables {
		if cable != lostCable {
			remaining++
		}
	}
	return remaining > 0
}

func worstCaseSurvival(c Country) bool {
	// The worst case is losing the country's single most critical
	// cable; if the country has only ONE cable, losing it is fatal
	// to connectivity (Section 11's Tonga case, generalized).
	if len(c.Cables) <= 1 {
		return false
	}
	return true
}

func main() {
	countries := []Country{
		{Name: "SingleCableNation (2022-style)", Cables: []string{"OnlyCable"}},
		{Name: "DualCableNation", Cables: []string{"CableA", "CableB"}},
		{Name: "MajorHub (e.g. Singapore-style)", Cables: []string{"CableA", "CableB", "CableC", "CableD", "CableE"}},
	}

	for _, c := range countries {
		fmt.Printf("%-32s cables=%d  survives worst-case single loss? %v\n",
			c.Name, len(c.Cables), worstCaseSurvival(c))
	}
	// Output shows the single-cable nation as the only one that
	// FAILS worst-case survival -- exactly Tonga's 2022 real
	// situation (Section 11), modeled as plain graph connectivity.
}
```

---

## 15b. A Field Guide: Comparing the Three Ownership Models Side by Side

Bringing Sections 2, 3, and 5 together into one comparison, since the three ownership patterns this chapter described are easy to blur together:

| Ownership model | Who funds it | Who can use the capacity | Real example | Primary motivation |
|---|---|---|---|---|
| Traditional multi-carrier consortium | Dozens of national/international telecoms, proportionally | Consortium members, per their funding share | SEA-ME-WE 5/6 | Spread cost and risk across many carriers with modest individual needs |
| Hyperscaler-led consortium | One or two dominant tech companies plus a smaller set of telecom partners | Mostly the lead company, with partners retaining agreed shares | 2Africa (Meta-led), MAREA (Microsoft/Meta/Telxius) | Guarantee capacity and control while sharing some cost and local regulatory relationships |
| Single-hyperscaler ownership | One company alone | That company, primarily for its own inter-datacenter and edge traffic | Dunant, Curie, Equiano (all Google) | Full control over capacity, upgrade timing, and route, independent of any partner's decisions |
| Independent carrier-neutral operator | A dedicated cable operating company, not itself a telecom or hyperscaler | Sold commercially to many customers | Southern Cross Cable Network | Operate cable infrastructure as a standalone commercial business, neutral among customers |

No single model has "won" — new construction includes examples of all four patterns being actively built in the same recent years, which is itself evidence that Chapter 124's marketplace framing (multiple viable economic structures coexisting, none centrally mandated) applies just as much to cable ownership as it does to peering and transit.

---

## 16. Common Misconceptions

- **"Submarine cables are mostly owned by governments."** Ownership today is split between traditional telecom consortia (Section 2) and, increasingly, private hyperscale companies (Section 3) — government ownership of the cables themselves is the exception, not the rule, though governments regulate landing rights (Section 7, Section 13) heavily.
- **"A cable repair is basically like fixing a break in a land cable — quick once you find it."** Section 8 showed the real bottleneck is almost always ship availability and weather-dependent deep-ocean operations, not the technical difficulty of the splice itself, which is comparatively fast once a ship is actually on-site.
- **"Redundancy always prevents regional outages."** Sections 9-11 all show cases where either concentrated, correlated damage (Hengchun, Mediterranean) or genuinely insufficient redundancy (Tonga) still produced real, measurable, multi-day-to-multi-week regional impact.
- **"Cable cuts are always accidental."** Most documented historical cases are (fishing, anchors, natural disasters), but Section 12's recent incidents show deliberate interference is now a live, publicly-discussed possibility taken seriously at a policy level, even where individual cases remain disputed.

---

## 16b. A Resilience Checklist for a Region Depending on Submarine Cables

Synthesizing Sections 5-13 into a practical checklist a country or large operator can hold itself to, directly informed by the contrast between Tonga (Section 11) and a major, well-connected hub:

- [ ] At least two independent cables, on genuinely different physical routes, not just different cable systems sharing the same narrow corridor (Sections 9-10's shared-corridor lesson).
- [ ] Landing stations in more than one physical location, so a single onshore incident (fire, flood, targeted attack) can't take out every cable's terrestrial handoff at once.
- [ ] A contractual or consortium relationship with repair capacity (Section 8) that doesn't depend entirely on being at the front of a shared repair-ship queue during a regional multi-fault event.
- [ ] Awareness of which named, real cables (Section 5) actually carry the region's traffic, and their ownership structure, since a single owner's business decisions can affect multiple "independent" cables at once.
- [ ] A tested fallback plan (satellite backup, Chapter 23, or rerouted capacity through a neighboring region) for the period between a cable cut and full repair, given that "weeks" is the realistic repair timeline established in Section 8.

---

## 17. Production Notes

- Large content and cloud companies increasingly disclose (via their own public blogs and infrastructure announcements) not just that they own specific cables, but their reasoning for route diversity across their whole cable portfolio — deliberately avoiding shared geological corridors (echoing Section 9-10's lessons) across their different cable investments.
- Network engineers at companies dependent on international connectivity for critical services maintain **cable-diversity awareness** as part of their disaster-recovery planning — knowing which physical cables carry which of their traffic paths, so that a single cable fault's blast radius on their own service can be estimated in advance rather than discovered during an incident.
- The **International Cable Protection Committee (ICPC)**, a real, long-running industry body, coordinates cable-protection best practices, incident data sharing, and engagement with governments and other seabed users (fishing industry, shipping) globally — a genuine, if low-profile, piece of the governance landscape Chapter 124, Section 11 discussed for the Internet more broadly.
- Repair-fleet capacity (Section 8's ~50-60 ships worldwide) is itself a subject of ongoing industry and government concern, given the growing number of cables in service relative to a repair fleet that hasn't grown proportionally — a real, structural risk factor independent of any single incident.

---

## 17b. Quick Reference: Who to Contact for What

A short, practical reference tying Sections 2-8's cast of characters together, useful for holding the whole operational picture in one place:

| Role | Who fills it, in practice |
|---|---|
| Route survey | Specialized marine survey contractors, often subcontracted by the cable owner |
| Manufacturing | SubCom, Alcatel Submarine Networks (ASN), or NEC — effectively the only three at global scale |
| Laying | The manufacturer's own cable-laying ships, or contracted marine operators like Orange Marine |
| Ownership/financing | A telecom consortium (Section 2), a hyperscaler-led group (Section 3), a single hyperscaler (Section 5), or an independent carrier-neutral operator |
| Repair | A shared maintenance consortium or the manufacturer's own repair-ship fleet (Section 8), coordinated across many cable owners |
| Landing rights and protection zones | National telecom regulators and, increasingly, national security review bodies (Section 13) |
| Industry coordination | The International Cable Protection Committee (ICPC), Section 17 |

---

## 18. What This Chapter Simplified

- Real cable ownership agreements (even the "single owner" cables in Section 5) often still involve minor partners, financing arrangements, or capacity-lease agreements with other carriers not fully captured by the simplified "owner" label used here.
- The repair process (Section 8) is described at a level that omits considerable engineering detail in fault localization, cable recovery techniques for different depths, and the specific splicing and testing standards used in practice.
- The geopolitical material (Section 12-13) reflects the state of public understanding and reporting as of this writing; several of the specific incidents referenced remain under investigation or subject to ongoing dispute, and should be treated as illustrative of a real, active policy concern rather than settled historical fact in every specific causal detail.

---

## 19. Interview Questions & Model Answers

**Beginner: "Who owns the submarine cables that carry the internet between continents?"**

*Model answer:* "Historically, most were owned by consortia of national and international telecom companies, each contributing funding in exchange for a guaranteed share of capacity — a model that's still active today. Since roughly the mid-2010s, large cloud and content companies like Google, Meta, Microsoft, and Amazon have increasingly built and owned cables themselves, sometimes entirely alone, because their own traffic volumes between their own global data centers are large enough to justify it, and it gives them direct control over capacity and upgrade timing rather than depending on a shared consortium's collective decisions."

**Intermediate: "Why does repairing a damaged submarine cable typically take weeks rather than days?"**

*Model answer:* "The technical splicing work itself is comparatively fast once a ship is actually on site. The real bottleneck is almost always ship availability — there's a genuinely small global repair fleet, on the order of fifty to sixty vessels covering every ocean, shared across many cable owners, so a newly damaged cable often has to wait its turn if the nearest available ship is already committed elsewhere. On top of that, the ship has to physically transit to the fault location, which can take days for a mid-ocean fault; deep-water grappling and retrieval is weather-dependent and can be delayed by rough seas; and if the fault is in another country's territorial waters, the ship may need government permits to operate there, which is a real, sometimes slow, regulatory step. All of those delays stack, which is why two to four weeks is a completely normal total repair timeline."

**Advanced: "Tonga lost nearly all international internet connectivity for about five weeks after a single cable was cut in 2022, while a similarly severe event affecting a well-connected hub like Singapore would likely cause only a partial slowdown. What structural difference explains this, and what would you recommend to reduce this kind of risk?"**

*Model answer:* "The difference is redundancy, not resilience of any individual cable. Tonga depended on a single submarine cable with no alternate physical path to the rest of the world, so cutting that one cable was equivalent to cutting Tonga's entire international connectivity — there was nothing for BGP or any other mechanism to reroute onto. A hub like Singapore, by contrast, is served by many independent cables along different routes, so losing any single one degrades available capacity and may increase latency as traffic reroutes across the survivors, but doesn't cause a total outage, mirroring the same packet-switching resilience principle from Chapter 9 — many independent paths, not one perfect path. The concrete recommendation for a country in Tonga's position is exactly what several Pacific island nations have pursued since: fund at least a second, geographically independent cable route, even though the economics are much harder to justify for a small population than they are for a major traffic hub, precisely because the cost of NOT doing so is a total national connectivity outage during any single cable's downtime."

---

## 20. Exercises

### Easy

1. Name the three companies that manufacture the overwhelming majority of the world's submarine cables.
2. What is a cable landing station, and what two things happen there?
3. List the two broad ownership models for submarine cables described in this chapter and one real example cable for each.

### Medium

4. Using Section 8's repair process, identify which single step is described as the "real bottleneck," and explain why owning more repair ships wouldn't fully solve the problem even if it were affordable.
5. Compare the 2006 Hengchun earthquake case (Section 9) and the 2022 Tonga case (Section 11): both involved a natural event severing cables, but the "redundancy failed" explanation differs between them. Explain the difference.
6. Explain, using Section 4's reasoning, why a cloud company might prefer building a cable entirely alone (like Google's Dunant) over joining a small joint venture (like Microsoft/Meta's MAREA) or a broad traditional consortium (like SEA-ME-WE).

### Hard

7. Using the Go code in Section 15, extend the model to account for "correlated risk" — two cables that are technically independent but run through the same narrow geographic corridor (like the Luzon Strait in Section 9) and are therefore both at risk from a single event. Sketch how you'd represent this in the `Country` struct and what `worstCaseSurvival` would need to check differently.
8. Research one submarine cable incident not covered in this chapter (any country, any decade) and write a short technical incident summary in the same style as Sections 9-12: what was damaged, what region was affected, and how long restoration took.
9. Argue, using Section 3's hyperscaler-ownership trend and Section 13's geopolitics, whether increasing hyperscaler ownership of submarine cables makes global Internet infrastructure more or less resilient overall, considering both the redundancy hyperscalers can afford to build and the national-security concerns their ownership raises in some countries.

---

## 21. Summary, and the Bridge to Chapter 127

| Term | Meaning |
|---|---|
| Consortium ownership | Multiple telecoms jointly funding a cable for proportional capacity shares |
| Hyperscaler-owned cable | A cable built and owned mostly or entirely by one large cloud/content company |
| Cable-laying ship | Purpose-built vessel that manufactures-to-order and lays cable along a surveyed route |
| Landing station | The shore facility where a submarine cable connects to terrestrial networks and power |
| Cable repair ship | One of a small global fleet (~50-60 vessels) that locates, retrieves, and splices damaged cables |
| Repair bottleneck | Ship availability, weather, and permitting — not splicing difficulty — as the real driver of multi-week repairs |
| Correlated risk corridor | A geographic chokepoint (Luzon Strait, Mediterranean near Alexandria) where multiple "independent" cables share the same physical risk |

This chapter has now shown who owns the ocean floor's fiber, and it deliberately kept surfacing the same handful of company names — Google, Meta, Microsoft, Amazon — building their own cables for their own traffic between their own data centers. That's not incidental. Chapter 127 is the full synthesis of exactly what those companies (plus Cloudflare) do with that private global infrastructure once it's built: how they run planet-scale networks that increasingly bypass the traditional peering, transit, and even cable-consortium marketplace this whole volume has been describing.
