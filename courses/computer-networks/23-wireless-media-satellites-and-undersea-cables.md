# Chapter 23: Wireless Media, Satellites, and the Undersea Cables That Carry the Internet

> *Ask most people how a message gets from New York to London, and "satellite" is a common guess. It's almost always wrong. This chapter closes Part 3 with the physical layer's biggest surprise: the internet that feels weightless and wireless is, at intercontinental scale, running through actual glass cables lying on the ocean floor — and understanding exactly why reveals the same physics (light in fiber, the speed of light itself) that every chapter in this volume has been building toward.*

---

## Table of Contents

1. [Two Remaining Problems: No Wires, and No Oceans](#1-two-remaining-problems-no-wires-and-no-oceans)
2. [Radio and Microwave Point-to-Point Links](#2-radio-and-microwave-point-to-point-links)
3. [Satellites — Trading a Cable for an Orbit](#3-satellites--trading-a-cable-for-an-orbit)
4. [GEO Satellite Latency — A Real Round-Trip Calculation](#4-geo-satellite-latency--a-real-round-trip-calculation)
5. [LEO Satellites — Starlink and the Low-Latency Alternative](#5-leo-satellites--starlink-and-the-low-latency-alternative)
6. [The Physical Reality: Undersea Cables Carry the Internet](#6-the-physical-reality-undersea-cables-carry-the-internet)
7. [A Real Map of Major Submarine Cable Routes](#7-a-real-map-of-major-submarine-cable-routes)
8. [How Undersea Cables Are Actually Built and Laid](#8-how-undersea-cables-are-actually-built-and-laid)
9. [What Threatens Undersea Cables](#9-what-threatens-undersea-cables)
10. [Comparing Every Physical Medium Side by Side](#10-comparing-every-physical-medium-side-by-side)
11. [Hands-On: Measuring Real Latency to Estimate Distance](#11-hands-on-measuring-real-latency-to-estimate-distance)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Interview Questions & Model Answers](#13-interview-questions--model-answers)
14. [Exercises](#14-exercises)
15. [Summary, and the Bridge to Part 4](#summary-and-the-bridge-to-part-4)

---

## 1. Two Remaining Problems: No Wires, and No Oceans

Chapters 21 and 22 solved the physical layer's core engineering problem — moving a fast, reliable signal from point A to point B — for two specific cases: electricity in copper over short distances, and light in glass over long ones. Both share one assumption this chapter finally drops: **a continuous physical cable exists between A and B.**

Two entire classes of real, everyday networking break that assumption:

1. **The device is moving, or a cable is simply impractical** — a phone in your pocket, a ship at sea, a research station in Antarctica, a rural home 40 km from the nearest fiber run. You can't lay a cable to a moving car, and you often can't justify the cost of laying one to a handful of homes across difficult terrain.
2. **The distance is an ocean.** Fiber (Chapter 22) is the best terrestrial medium there is, but continents are separated by thousands of kilometers of open water, and someone still has to physically get a cable from one shore to the other.

This chapter treats both problems in turn: wireless links (radio, microwave, satellite) for problem 1, and — because it surprises almost everyone the first time they learn it — the actual physical answer to problem 2, which turns out to still be fiber, just underwater.

---

## 2. Radio and Microwave Point-to-Point Links

The most direct wireless analog to a physical cable is a **point-to-point microwave or radio link**: two directional dish antennas, aimed precisely at each other, exchanging a modulated radio signal (Chapter 15's modulation concepts apply exactly the same way here as they do to any other carrier wave) with no cable in between at all.

**Why this works, and its inherent limit:** unlike Wi-Fi (Chapter 86) or cellular (Chapter 90+), which deliberately broadcast in most directions to serve many devices, a point-to-point link focuses nearly all of its transmitted energy into a narrow beam aimed at one specific receiver — dramatically increasing effective range and resistance to interference for a given transmit power, at the cost of needing precise physical alignment and, critically, a clear **line of sight** between the two dishes. Microwave frequencies (commonly 6-80 GHz for terrestrial backhaul links) don't bend around hills or buildings the way lower-frequency AM radio does, and the curvature of the Earth itself imposes a hard geometric limit — a typical tower-to-tower microwave hop tops out around 40-50 km before the Earth's curvature drops the horizon below line of sight, which is exactly why long microwave backhaul routes are built as a **chain of towers**, each one relaying the signal to the next, rather than one impossibly long single hop.

```
Tower A ))))  >>>>>>>>>>>>>>>>>>  (((( Tower B ))))  >>>>>>>>>>>>>>>>>>  (((( Tower C

   Each hop: ~20-50 km, line-of-sight required, limited by Earth's
   curvature and terrain — not by signal attenuation alone.
```

**Real use:** microwave backhaul links are the unglamorous workhorse connecting many cellular towers (Chapter 91-92) back to the core network in areas where trenching fiber is too expensive or too slow to deploy, and they're a standard tool for quickly connecting two buildings across a city without digging up streets for a fiber run. They're also the technology behind large financial firms' famous ultra-low-latency microwave links between stock exchanges (e.g., Chicago to New York), chosen specifically because microwave through air is very slightly *faster* than light in fiber — air's refractive index is far closer to 1 than glass's ~1.47-1.48 (Chapter 22, Section 7), so a microwave signal traveling in a straight line through air genuinely beats an equivalent fiber route over the same distance, a real, if narrow and expensive, latency edge worth knowing exists.

---

## 3. Satellites — Trading a Cable for an Orbit

When a tower chain isn't practical either — the middle of an ocean, a mountain range, a genuinely remote area with no infrastructure at all — the remaining option is to bounce the signal off a relay parked far above the atmosphere: a **satellite**. The tradeoff a satellite link makes is stated most honestly as: *unlimited reach, in exchange for distance-driven latency that no engineering can avoid, because that distance is fixed by orbital mechanics, not signal design.*

Three orbital regimes matter for internet connectivity, and the altitude difference between them is the single most important fact in this section, because — as Sections 4 and 5 show with real numbers — altitude directly determines latency via nothing more complicated than the speed of light:

```
GEO (Geostationary Earth Orbit):   ~35,786 km altitude
                                     Satellite orbits once per 24 hours,
                                     matching Earth's rotation, so it appears
                                     FIXED in the sky from the ground —
                                     one satellite can cover roughly a third
                                     of the globe continuously.

MEO (Medium Earth Orbit):          ~2,000-35,786 km altitude
                                     Used mostly for GPS/GNSS satellites,
                                     less common for internet connectivity.

LEO (Low Earth Orbit):             ~340-2,000 km altitude
                                     (Starlink operates mostly around ~550 km)
                                     Satellite moves rapidly relative to the
                                     ground — a constellation of MANY satellites
                                     is needed for continuous coverage of one spot,
                                     since any single satellite passes overhead
                                     and out of view again within minutes.
```

---

## 4. GEO Satellite Latency — A Real Round-Trip Calculation

This is the calculation that explains, more concretely than any qualitative description could, why GEO satellite internet has a reputation for feeling "laggy" even with plenty of bandwidth.

**Step 1 — one-way travel time from the ground to a GEO satellite:**

```
GEO altitude:             35,786 km
Speed of light (vacuum):  ~300,000 km/s   (the signal travels through
                                            atmosphere and space, both very
                                            close to vacuum speed — unlike
                                            fiber's ~200,000 km/s from Ch. 22)

one-way time = distance / speed
             = 35,786 km / 300,000 km/s
             = 0.11929 seconds
             ≈ 119.3 ms
```

**Step 2 — a full, realistic round trip.** A satellite internet user's traffic doesn't just bounce once — it typically goes user → satellite → ground gateway (the ISP's connection to the wider internet) on the way there, and gateway → satellite → user on the way back:

```
User uplink to satellite:            ~119.3 ms
Satellite downlink to ground gateway: ~119.3 ms
                                       ---------
One-way, user-to-gateway:            ~238.6 ms

Response takes the same path back:
Gateway uplink to satellite:          ~119.3 ms
Satellite downlink to user:           ~119.3 ms
                                       ---------
One-way, gateway-to-user:            ~238.6 ms

THEORETICAL MINIMUM ROUND TRIP:      ~477.2 ms
```

**Real-world GEO satellite internet services (e.g., HughesNet, Viasat) commonly report round-trip latencies around 500-600+ ms** — noticeably higher than the ~477 ms speed-of-light minimum above, because real systems add processing delay at the satellite itself, switching and queuing delay at the ground gateway, and the terrestrial network path from the gateway to whatever server the user is actually reaching. The speed-of-light figure is a hard floor no engineering can improve on; everything above that floor is where real-world equipment and network design add their own (reducible, but never zero) overhead.

**Why this matters practically:** ~500-600 ms of round-trip latency is well within the range that makes real-time applications (video calls, competitive online gaming) feel noticeably laggy, even though the *bandwidth* GEO satellite links offer can be perfectly reasonable for downloading files or streaming video (which tolerate latency far better than they tolerate low throughput). This is a direct, practical illustration of a distinction Chapter 16 introduced abstractly — bandwidth and latency are different properties of a link, and a link can be excellent on one axis while being fundamentally limited on the other by nothing more than geometry.

---

## 5. LEO Satellites — Starlink and the Low-Latency Alternative

Run the exact same calculation from Section 4, but for a LEO constellation like Starlink orbiting at roughly **550 km** instead of GEO's 35,786 km:

```
one-way time = 550 km / 300,000 km/s
             = 0.001833 seconds
             ≈ 1.83 ms

Round trip (user -> satellite -> gateway -> satellite -> user,
same 4-hop structure as Section 4):
             ≈ 4 x 1.83 ms ≈ 7.3 ms   (theoretical minimum)
```

That's roughly **65x lower** than GEO's theoretical minimum, purely because the satellite is roughly 65x closer — a direct, clean illustration of why altitude is the single dominant factor in satellite latency.

**Why real-world LEO latency is higher than 7.3 ms, and what that reveals.** Measured Starlink latency in practice is typically reported in the range of **20-60 ms**, not the ~7 ms theoretical floor. The gap comes from real, non-speed-of-light factors that matter far more here than they did for GEO's already-large latency budget: signal processing time in the satellite and user terminal, the routing path once traffic reaches a ground gateway (which might be hundreds of kilometers from both the user and the eventual destination server), and — for early Starlink specifically — the fact that a single satellite's coverage window over any one spot lasts only minutes before handoff to the next satellite is needed (a wireless-specific problem conceptually related to cellular handoff, previewed properly in Chapter 91-92). Newer LEO constellations increasingly add **inter-satellite laser links** — satellites relaying data directly to each other in space rather than always routing back down to a ground station — which can route long-distance traffic through vacuum (where light travels at the full ~300,000 km/s, faster than fiber's ~200,000 km/s) for part of the journey, an active area of real deployment worth watching.

```
                    Theoretical minimum RTT     Typical real-world RTT
GEO (35,786 km):    ~477 ms                     ~500-600+ ms
LEO (~550 km):      ~7.3 ms                     ~20-60 ms
Terrestrial fiber
(same-continent):   varies with distance         ~1-40 ms typical
```

---

## 6. The Physical Reality: Undersea Cables Carry the Internet

Given Section 3-5's honest accounting of satellite latency, and given that satellites also have real bandwidth ceilings per beam that are far below what a single fiber strand can carry (Chapter 22, Section 9's DWDM figures — tens of terabits on one strand — dwarf any current satellite's per-beam capacity), the obvious engineering question is: **how does intercontinental internet traffic actually cross oceans, if not mostly by satellite?**

The answer, and this genuinely surprises most people encountering it for the first time: **well over 95% (commonly cited figures range from 95% to 99%) of intercontinental internet and telephone traffic travels through undersea fiber optic cables laid directly on the ocean floor** — physically the same single-mode fiber technology explained in full in Chapter 22, just insulated, armored, and run for thousands of kilometers underwater instead of tens of kilometers between buildings. Satellites remain important for ships, aircraft, remote areas with no cable access, and as a backup path — but they are not, and were never, the primary way continents talk to each other.

This is precisely why the "somewhere in space" mental model of the internet — while forgivable, since Wi-Fi and cellular genuinely are wireless for the "last mile" — is physically wrong for the vast majority of a message's actual journey. Every video call to another continent, every international bank transfer, nearly every byte of transoceanic web traffic is, for the overwhelming majority of its physical trip, light bouncing down a glass fiber lying on the seabed, exactly as described in Chapter 22 — just for a few thousand kilometers instead of a few hundred.

---

## 7. A Real Map of Major Submarine Cable Routes

As of the early-to-mid 2020s, there are **over 500 active submarine cables** spanning a combined length of more than **1.4 million kilometers** — enough to circle the Earth roughly 35 times. A simplified schematic of the busiest intercontinental routes:

```
                     NORTH AMERICA
                          |
          (transatlantic) |  (transpacific)
                 ________/ \________
                /                    \
            EUROPE                  ASIA
               |                      |
        (Europe-Asia/Africa)   (intra-Asia, Asia-Oceania)
               |                      |
             AFRICA               OCEANIA
```

**Approximate real routes, distances, and latencies** (theoretical latency computed using Chapter 22's ~200,000 km/s speed of light in fiber; real-world figures reflect actual measured round-trip times, which are always somewhat higher due to equipment, routing, and the cable not running in a perfectly straight line):

| Route | Approx. cable route distance | A real cable on this route | Theoretical one-way (light-in-fiber) | Typical real-world RTT |
|---|---|---|---|---|
| New York ↔ London | ~5,900 km | TAT-14; Grace Hopper (Google, 2022) | ~29.5 ms | ~70-76 ms |
| Virginia Beach ↔ Bilbao, Spain | ~6,600 km | MAREA (Meta/Microsoft/Telxius, 2017; ~160 Tbps design capacity) | ~33 ms | ~65-75 ms |
| Virginia Beach ↔ France | ~6,400 km | Dunant (Google, 2021; ~250 Tbps design capacity) | ~32 ms | ~65-75 ms |
| Los Angeles ↔ Tokyo/Chiba | ~9,000 km | FASTER (Google + consortium, 2016; ~60 Tbps) | ~45 ms | ~100-120 ms |
| Los Angeles ↔ Chikura, Japan | ~9,000 km | Unity (Google + consortium, 2010) | ~45 ms | ~100-120 ms |
| Sydney ↔ Los Angeles | ~12,000 km | Southern Cross Cable Network | ~60 ms | ~130-160 ms |
| Singapore ↔ Marseille (Asia-Europe) | ~15,000-20,000 km depending on route | SEA-ME-WE cable family | ~75-100 ms | ~160-220 ms |
| Circling Africa, connecting Europe/Middle East/Africa | ~45,000 km total length | 2Africa (Meta + consortium; largest subsea cable project by length as of its construction) | n/a (multi-segment) | varies by segment |

Notice the consistent pattern: real-world round-trip times run roughly **1.5-2.5x higher** than the pure speed-of-light theoretical minimum for every route. That gap is the same story as Section 4's GEO satellite figures — the underlying cable route is never a perfectly straight line (it follows safe seabed contours, avoiding fault lines, fishing zones, and existing infrastructure), and real traffic passes through routers, switches, and regeneration/amplification equipment (Chapter 22, Section 10) at multiple points along the way, each adding a small, real processing delay on top of the physical travel time.

---

## 8. How Undersea Cables Are Actually Built and Laid

A submarine cable is not simply Chapter 22's fiber cable dropped in the ocean unmodified — it's built in armored layers specifically engineered for an environment no terrestrial fiber run ever faces:

```
   ┌───────────────────────────────────────────────────┐
   │  Polyethylene outer sheath (or steel wire armor     │
   │    in shallow, high-risk zones near shore)          │
   │  Steel wire strands (tensile strength for laying    │
   │    and to resist anchors/fishing trawls)            │
   │  Copper or aluminum power conductor (powers the      │
   │    optical amplifiers spaced along the route,         │
   │    Chapter 22 Section 10)                            │
   │  Petroleum jelly / water-blocking layer               │
   │  Optical fiber core(s) — typically 8-24 individual    │
   │    fiber pairs, each one exactly the single-mode      │
   │    fiber technology explained in Chapter 22           │
   └───────────────────────────────────────────────────┘
```

The cable is manufactured in one continuous run at a specialized factory, coiled onto a purpose-built **cable-laying ship**, and paid out along a carefully surveyed seabed route as the ship moves — in shallow coastal waters, a **plow** buries the cable a meter or more beneath the seabed for protection from anchors and trawling nets (a real, common cause of cable damage, covered in Section 9); in deep ocean, the cable is simply laid on the seabed surface, since the water itself is generally protection enough at those depths. Just as Chapter 22 explained for terrestrial long-haul fiber, submarine routes need periodic **optical amplifiers** — housed in pressure-resistant repeater units placed roughly every 50-100 km along the cable, powered by a DC current run down the cable's own copper conductor from stations on shore — to keep boosting the light signal purely optically across a multi-thousand-kilometer journey without ever converting it back to electricity mid-ocean.

---

## 9. What Threatens Undersea Cables

Given how much of the internet's intercontinental traffic depends on this relatively small number of physical cables, their vulnerabilities are a genuine, actively-managed operational concern, not a hypothetical:

- **Fishing trawls and ship anchors** are, by a wide margin, the most common real cause of submarine cable damage — a dragged anchor or trawl net in shallow coastal waters can sever a cable that isn't buried deeply enough, which is exactly why the plowing/burial practice mentioned in Section 8 exists.
- **Underwater landslides and seabed movement** near geologically active areas (e.g., a well-documented 2006 incident near Taiwan following an earthquake, which damaged multiple cables and disrupted regional internet connectivity for days) can sever multiple cables at once if they're routed too close together through the same geologically risky corridor.
- **Deliberate sabotage** is a real, publicly-discussed geopolitical and military concern in recent years, precisely because so much critical infrastructure and communication depends on a physically identifiable, largely undefended set of cables on the seabed.
- **Redundancy is the standard mitigation, not physical hardening.** Rather than trying to make any single cable indestructible, operators and countries deliberately build **multiple independent cables along different routes** between the same regions (Section 7's table already shows several separate cables serving the same New York-Europe and US-Asia corridors), so that damage to any one cable reroutes traffic across the others with a latency and capacity cost, rather than causing an outright regional outage. This is the same core resilience principle packet switching itself was built on back in Chapter 9 — don't rely on one path, rely on many, and let the failure of any single part degrade rather than collapse the whole system.

---

## 10. Comparing Every Physical Medium Side by Side

Bringing together everything from Chapters 21-23:

| Medium | Typical reach | Typical latency driver | Typical bandwidth ceiling | Best use case |
|---|---|---|---|---|
| Copper twisted pair (Ch. 21) | ~100 m | Negligible over short runs | 1-40 Gbps (category-dependent) | Last few meters: desk, access point |
| Terrestrial fiber (Ch. 22) | Tens to 100+ km per span | Distance / speed-of-light-in-glass | Tens of Tbps per strand (DWDM) | Backbone, data center, long-haul |
| Point-to-point microwave (this chapter) | ~20-50 km per hop | Distance + processing at each relay | Up to a few Gbps per link | Cellular backhaul, quick building-to-building links |
| GEO satellite (this chapter) | Near-global (fixed footprint) | Orbital altitude (~477 ms RTT floor) | Moderate, shared across a wide footprint | Remote areas, maritime, broadcast |
| LEO satellite (this chapter) | Near-global (via constellation) | Lower orbital altitude (~7 ms RTT floor) + processing | Improving rapidly, still below fiber | Remote/rural broadband, mobility, backup paths |
| Submarine fiber cable (this chapter) | Thousands to tens of thousands of km | Distance / speed-of-light-in-glass | Tens of Tbps per cable (many fiber pairs, DWDM) | The actual backbone of intercontinental internet traffic |

---

## 11. Hands-On: Measuring Real Latency to Estimate Distance

You can observe Sections 4, 5, and 7's math directly, using nothing but `ping` (fully explained in Chapter 54, but usable right now):

```
$ ping -c 4 london-based-server.example
64 bytes from ...: icmp_seq=1 ttl=54 time=72.3 ms
64 bytes from ...: icmp_seq=2 ttl=54 time=71.9 ms
64 bytes from ...: icmp_seq=3 ttl=54 time=73.1 ms
64 bytes from ...: icmp_seq=4 ttl=54 time=72.5 ms
```

If you're pinging from New York, ~72 ms round trip lines up closely with Section 7's real-world NYC-London figure — direct, personal, measured evidence of a signal's round trip through an actual transatlantic fiber cable, not a hypothetical number from a table.

**Experiment to run yourself:** pick a handful of well-known servers or websites hosted in different, known regions (a CDN edge or cloud region in Europe, one in East Asia, one in Australia — Chapter 96 covers how to reason about where a given service is actually served from), `ping` each one, and use Section 7's ~200,000 km/s figure to back-calculate an approximate one-way distance from your measured round-trip time (`distance ≈ (RTT_ms / 2 / 1000) × 200,000 km`). Compare your estimate to the real geographic distance — the gap between your estimate and reality is a direct, hands-on measurement of how much "extra" latency comes from routing, equipment, and non-straight-line cable paths, exactly as Section 7 described in the abstract.

---

## 12. Common Misconceptions

- **"Most international internet traffic goes through satellites."** No — Section 6 was explicit that well over 95% of intercontinental traffic runs through undersea fiber cables. Satellites matter enormously for coverage in places cables can't reach, but they are not the backbone of intercontinental connectivity.
- **"Satellite internet is slow because of weak signals or bad weather."** Not fundamentally — Section 4's calculation showed GEO satellite latency is dominated by the unavoidable speed-of-light round trip to a fixed, very high orbital altitude, not by signal quality. Weather and equipment add some real overhead, but the ~477 ms theoretical floor exists purely from geometry and would remain even with a perfect signal.
- **"LEO satellites like Starlink will make undersea cables obsolete."** Unlikely with current and near-future technology — Section 10's bandwidth comparison shows a single fiber strand's DWDM capacity (tens of terabits) still dwarfs any current satellite constellation's aggregate throughput per region, and Section 6's economics (once a cable is laid, its capacity can be upgraded for decades via new DWDM endpoint electronics, per Chapter 22 Section 9) remain very favorable for high-volume, fixed-location intercontinental traffic.
- **"A single cable cut takes down the whole internet between two continents."** No — Section 9 explained that redundancy across multiple independently-routed cables is the standard design response precisely to avoid this; a real cable cut typically degrades capacity and increases latency (as traffic reroutes) rather than causing a total outage, though a sufficiently severe multi-cable event (like the 2006 Taiwan earthquake mentioned in Section 9) can still cause serious, real regional disruption.
- **"Microwave/radio links are obsolete now that fiber exists everywhere."** No — Section 2 covered real, current use cases (cellular backhaul in hard-to-trench areas, rapid building-to-building links, and even genuine latency-arbitrage use in financial trading) where microwave remains the better engineering or economic choice specifically because it avoids the cost/time of trenching fiber, or because air's near-vacuum speed of light gives it a real (if narrow) latency edge over fiber for a given distance.

---

## 13. Interview Questions & Model Answers

**Q1 (Beginner): Why does a GEO satellite internet connection have noticeably higher latency than a typical wired or fiber connection, even with a strong signal and modern equipment?**

"The dominant factor is pure geometry, not signal quality. A geostationary satellite orbits at about 35,786 km above the Earth, and a signal traveling at the speed of light takes about 119 milliseconds to cover that distance one way. A typical request has to go up to the satellite and down to a ground gateway, then the response has to go back up and down again — four one-way hops of roughly 119 ms each, which adds up to a theoretical minimum round trip of about 477 milliseconds, before accounting for any real-world processing or routing delay. That's a hard floor set by physics and the satellite's fixed altitude; no improvement in the modem, antenna, or signal strength can get below it. That's why real GEO satellite services typically report round-trip times of 500-600+ milliseconds, noticeably worse than the tens of milliseconds typical of wired or terrestrial fiber connections, even though GEO satellite bandwidth can be perfectly reasonable for latency-tolerant tasks like downloads or video streaming."

**Q2 (Intermediate): Why do LEO satellite constellations like Starlink have such dramatically lower latency than GEO satellites, and why is their real-world measured latency still higher than the pure speed-of-light calculation suggests?**

"LEO satellites orbit at roughly 550 km, compared to GEO's approximately 35,786 km — around 65 times closer to Earth. Since latency here is fundamentally a function of distance divided by the speed of light, that altitude difference translates directly into a proportionally much lower theoretical round-trip floor, on the order of 7 milliseconds for LEO versus about 477 milliseconds for GEO. In practice, though, measured Starlink latency is typically in the 20-60 millisecond range, not anywhere near that 7 millisecond floor. The gap comes from factors the pure speed-of-light calculation doesn't capture: processing delay in the satellite and user terminal, the terrestrial routing path from the ground gateway to whatever server the user is actually trying to reach (which might be hundreds of kilometers away and add its own latency), and the practical realities of a fast-moving satellite constellation, including handoffs between satellites as any single one passes overhead and out of range within minutes. Even accounting for all of that overhead, LEO's real-world latency is still typically an order of magnitude better than GEO's, which is precisely why LEO constellations are viable for latency-sensitive applications that GEO satellite internet never realistically supported."

**Q3 (Advanced): If satellites can, in principle, offer global coverage without the enormous cost and years-long construction of laying a submarine cable, why does undersea fiber remain the dominant carrier of intercontinental internet traffic?**

"It comes down to the combination of bandwidth economics and the cost curve over time. A single strand of modern single-mode fiber, using dense wavelength division multiplexing, can carry tens of terabits per second, and a submarine cable typically bundles many such strands together — giving a single cable an aggregate capacity that dwarfs what any current satellite constellation can deliver to a comparable region, because satellite bandwidth is fundamentally limited by available radio spectrum and per-satellite beam capacity, a much more constrained resource than adding more wavelengths to an already-laid fiber. There's also a crucial cost-over-time argument: once a submarine cable is laid — admittedly a very expensive, multi-year undertaking — its capacity can be upgraded for decades afterward simply by installing newer, denser DWDM transceiver equipment at each end, without ever touching the cable itself again. Satellites, by contrast, have a fixed operational lifespan and need to be physically replaced or supplemented with new launches to grow capacity or maintain coverage. For the extremely high, sustained traffic volumes between major population centers on different continents, cables remain the more capacity-efficient and, over their operational lifetime, more cost-efficient choice. Satellites earn their place precisely where cables can't reach economically at all — ships, aircraft, remote regions, and as a resilient backup path — rather than as a wholesale replacement for the undersea cable backbone."

---

## 14. Exercises

### Easy

1. Using Section 4's method, calculate the one-way and round-trip theoretical latency for a satellite at MEO altitude (assume ~20,000 km, roughly a GPS satellite's altitude) and compare it to the GEO and LEO figures already calculated.
2. List two reasons a point-to-point microwave link might be chosen over trenching a new fiber run, from Section 2.
3. Using Section 7's table, identify which listed route has the largest gap between its theoretical one-way latency and its typical real-world round-trip time, and suggest one possible reason (referencing Section 7's explanation of that gap).

### Medium

4. A financial firm wants the lowest possible latency between two data centers 1,200 km apart and is choosing between a straight-line microwave relay chain and a fiber route that (due to available rights-of-way) must run 1,600 km. Using Chapter 22's fiber speed (~200,000 km/s) and this chapter's note about air's near-vacuum speed of light (~300,000 km/s), calculate the theoretical one-way latency for each option and explain which one wins and by how much.
5. Using Section 6 and Section 10's comparison table, write a two-to-three sentence explanation (as if to a non-technical stakeholder) of why a company should not plan to replace its transatlantic fiber-based connectivity with a LEO satellite service for a high-bandwidth, latency-tolerant application like nightly data backups.
6. Explain, using Section 9's redundancy principle and Chapter 9's original packet-switching resilience argument, why building three separate transatlantic cables along three different routes is a better resilience strategy than building one cable with three times the fiber pairs bundled inside it.

### Hard

7. Using Section 7's real cable table, pick the Los Angeles-Tokyo route and independently verify the ~45 ms theoretical one-way latency figure using the great-circle distance between the two cities (research or estimate it) and Chapter 22's ~200,000 km/s fiber speed figure. Discuss any discrepancy between your calculation and the table's ~9,000 km route-distance assumption, and explain why a submarine cable's actual laid length is typically longer than the straight-line great-circle distance between its endpoints.
8. Research one real historical submarine cable disruption (for example, the 2006 Taiwan earthquake cable damage mentioned in Section 9, or a more recent publicly-reported incident) and write a short technical incident summary: what was damaged, what regions were affected, how traffic was rerouted, and how long full restoration took.
9. Design, at a conceptual level, a redundancy strategy for a hypothetical new data center campus in Singapore that needs guaranteed connectivity to both Europe and North America even if any single submarine cable route is damaged. Reference Section 7's real routes and Section 9's redundancy principle in your answer, and explicitly state what latency and capacity trade-offs your backup path(s) would introduce compared to the primary route.

---

## Summary, and the Bridge to Part 4

| Term | Meaning |
|---|---|
| Point-to-point microwave link | Directional radio link between two dishes, needs line of sight, limited by Earth's curvature (~20-50 km/hop) |
| GEO (Geostationary orbit) | ~35,786 km altitude; fixed position in the sky; ~477 ms theoretical round-trip latency floor |
| LEO (Low Earth orbit) | ~340-2,000 km altitude (Starlink ~550 km); ~7 ms theoretical round-trip floor; needs a satellite constellation |
| Submarine cable | Armored single-mode fiber cable laid on the ocean floor; carries 95%+ of intercontinental internet traffic |
| Optical amplifier (repeater) | Powered unit every 50-100 km on a long submarine route, boosts the light signal purely optically |
| Cable redundancy | Multiple independent cables on different routes between the same regions, to survive a single cable's damage |
| Inter-satellite laser link | Newer LEO technology relaying data satellite-to-satellite through vacuum, bypassing some ground routing |

Parts 1-3 of this course have now covered, in order: what a network fundamentally is (Part 1), how networking got here historically (Part 2), and — across Chapters 14 through this one — exactly what physically moves between two computers, down to real speed-of-light numbers in copper, glass, and vacuum, and exactly how noise, attenuation, and interference are detected and corrected along the way. Every protocol in the rest of this course — Ethernet frames, IP packets, TCP segments, HTTP requests, TLS handshakes — ultimately reduces to voltage changes on copper, light pulses in fiber, or radio waves through air, exactly as this volume described.

But a single physical medium can't be the whole story. If Ethernet, Wi-Fi, IP addressing, and TCP reliability all had to be handled by one tangled, do-everything protocol, a change to how Wi-Fi encodes bits would force a rewrite of every web browser on Earth. Chapter 24 opens Part 4 by asking the next unavoidable question directly: **why does networking need to be built in layers at all** — and shows, with a concrete before-and-after, exactly what a networking stack without layering would actually look like, and why it collapses under its own weight.
