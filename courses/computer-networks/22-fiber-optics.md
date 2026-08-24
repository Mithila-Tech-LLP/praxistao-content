# Chapter 22: Fiber Optics — Single-Mode, Multi-Mode, and Transceivers

> *Chapter 21 showed copper twisted pair fighting a constant electromagnetic battle against itself — crosstalk, attenuation, and a hard 100-meter ceiling, all consequences of pushing electrons through a conductor. Fiber optics sidesteps almost the entire fight by changing the medium entirely: instead of electrons in copper, it sends the signal as pulses of light trapped inside a strand of glass thinner than a human hair — a medium that is, by its physical nature, immune to the electromagnetic crosstalk that dominated the last chapter.*

---

## Table of Contents

1. [The Problem Copper Can't Solve: Distance and Bandwidth Together](#1-the-problem-copper-cant-solve-distance-and-bandwidth-together)
2. [Total Internal Reflection — How Light Gets Trapped in Glass](#2-total-internal-reflection--how-light-gets-trapped-in-glass)
3. [The Anatomy of an Optical Fiber](#3-the-anatomy-of-an-optical-fiber)
4. [Multi-Mode Fiber — Many Paths, Shorter Reach](#4-multi-mode-fiber--many-paths-shorter-reach)
5. [Single-Mode Fiber — One Path, Long Reach](#5-single-mode-fiber--one-path-long-reach)
6. [Why Fiber Beats Copper: The Physics, Made Concrete](#6-why-fiber-beats-copper-the-physics-made-concrete)
7. [The Speed of Light in Glass — A Real Number, Worked](#7-the-speed-of-light-in-glass--a-real-number-worked)
8. [Transceivers: SFP, SFP+, QSFP, and the Optics Ecosystem](#8-transceivers-sfp-sfp-qsfp-and-the-optics-ecosystem)
9. [Wavelengths, WDM, and How Fiber Scales to Terabits](#9-wavelengths-wdm-and-how-fiber-scales-to-terabits)
10. [Real-World Fiber Specs and Distance Limits](#10-real-world-fiber-specs-and-distance-limits)
11. [Fiber to the Home: PON and Why Your ISP Loves It](#11-fiber-to-the-home-pon-and-why-your-isp-loves-it)
12. [Hands-On: Reading a Transceiver and Testing a Fiber Link](#12-hands-on-reading-a-transceiver-and-testing-a-fiber-link)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#summary)

---

## 1. The Problem Copper Can't Solve: Distance and Bandwidth Together

Chapter 21 ended with a puzzle worth restating precisely: copper's electromagnetic signal actually travels at a speed (60-70% of light speed) not wildly different from what this chapter's medium achieves. So why does fiber utterly dominate every long-haul, high-bandwidth link on Earth — every intercontinental submarine cable (Chapter 23), every internet backbone, every data center's core network — while copper is relegated to the last 100 meters to a desk or an access point?

The answer isn't raw signal speed. It's two compounding limits copper cannot escape, both explained in Chapter 21: **attenuation** (electrical resistance and skin effect worsen with both distance and frequency) and **crosstalk** (any two nearby current-carrying conductors interfere with each other, worse at higher frequencies). Stack enough of Chapter 18's Shannon-Hartley math on top of those two limits, and copper's usable bandwidth-times-distance product hits a wall that no amount of clever twisting or shielding escapes. Fiber's entire value proposition is that it uses a physical carrier — light in glass — for which neither of those two limits meaningfully applies. This chapter explains exactly why.

---

## 2. Total Internal Reflection — How Light Gets Trapped in Glass

**The naive first idea:** just shine a light down a glass rod and read pulses at the other end (light on = 1, light off = 0 — an encoding scheme, note, that traces straight back to Chapter 14's original definition of a signal). The obvious problem: light, like anything traveling through a transparent medium, doesn't travel in a perfectly straight, lossless line forever — any curve in the rod, or any imperfection, would let light escape out the side, and a fiber cable that had to stay perfectly straight for kilometers would be useless.

**The real solution is a genuinely elegant piece of physics called total internal reflection (TIR).** Here's the intuitive version first: light travels at different speeds in different transparent materials (this is why a straw looks "bent" where it enters a glass of water — light bends, or *refracts*, when crossing between materials of different density). A fiber is built from two layers of glass with **deliberately different densities**: a **core** in the center (denser, so light travels slightly slower in it) surrounded by a **cladding** (less dense, light travels slightly faster in it). When light traveling inside the denser core hits the boundary with the less-dense cladding at a shallow enough angle, instead of passing through and escaping, **100% of it bounces back inward** — as if the boundary had become a perfect mirror.

**Where the "mirror" analogy is incomplete:** a real mirror reflects light because of a reflective metal coating; total internal reflection needs no coating or mirror at all — it's a pure consequence of the two materials' different **refractive indices** (a number describing how much a material slows down and bends light; core glass is engineered to have a refractive index `n1` slightly higher than the cladding's `n2`) and the angle at which light strikes the boundary. There's a precise **critical angle**, determined by the ratio of the two refractive indices, above which every single photon reflects perfectly back into the core, with essentially zero loss at that boundary. Fiber manufacturers engineer the core/cladding refractive index difference specifically so that light launched into the fiber within a certain acceptance cone always strikes the core-cladding boundary steeper than that critical angle, guaranteeing it stays trapped and bounces its way down the fiber, hugging the same core the entire distance — even around gentle bends — for kilometers.

```
                    cladding (lower refractive index, n2)
        ______________________________________________
       /                                                \
      |    core (higher refractive index, n1)             |
      |  light ---> \  /\  /\  /\  /\  /\  /---> light out |
      |               \/  \/  \/  \/  \/                   |
       \______________________________________________/
                    cladding (lower refractive index, n2)

  Each time the light ray hits the core-cladding boundary at a
  shallow enough angle, it reflects PERFECTLY back into the core
  (total internal reflection) instead of escaping. No metal, no
  mirror coating — purely a property of the two materials' densities.
```

---

## 3. The Anatomy of an Optical Fiber

A real fiber optic cable, from the inside out:

```
   ┌─────────────────────────────────────────┐
   │  Core       (glass, ~8-62.5 microns)     │  <- light actually travels here
   │  Cladding   (glass, ~125 microns total)  │  <- traps light via TIR (Section 2)
   │  Coating    (plastic buffer, protective)  │  <- protects the fragile glass
   │  Strengthening fibers (aramid/Kevlar)     │  <- tensile strength, so the glass
   │                                            │     itself never has to bear pulling force
   │  Outer jacket (PVC or plenum-rated)       │  <- physical/environmental protection
   └─────────────────────────────────────────┘
```

For scale: a human hair is roughly 50-100 microns wide. A single-mode fiber's light-carrying core (Section 5) is only about **9 microns** — meaning the actual "wire" the signal travels through is far thinner than a hair, and the glass carrying gigabits (or, with the techniques in Section 9, terabits) of data per second down an ocean floor is mechanically the most fragile part of the entire assembly — which is exactly why real fiber cables devote most of their bulk to protective layers, not the light-carrying core itself.

---

## 4. Multi-Mode Fiber — Many Paths, Shorter Reach

**The core insight for this section:** the word "mode" in "multi-mode" refers to the distinct paths, or angles, at which light can bounce down the fiber while still satisfying the total-internal-reflection condition from Section 2. A **wider core** (typically **50 or 62.5 microns**) is physically large enough for light to enter and bounce down it at many different valid angles simultaneously — hence "multi-mode."

**Why multiple modes is a problem, not a feature:** different modes (different bounce angles) travel genuinely different *physical path lengths* down the same length of fiber — a ray bouncing steeply zigzags back and forth far more than a ray traveling nearly straight down the core's center, even though both started at the same instant from the same LED or laser pulse. This means a single, sharp pulse of light launched into multi-mode fiber arrives at the far end **smeared out over time**, because its different modes traveled different effective distances and arrive at slightly different moments. This smearing is called **modal dispersion**, and it directly limits both the distance and the maximum data rate multi-mode fiber can reliably support — push data too fast or too far, and consecutive pulses smear into each other and become indistinguishable, an optical version of the intersymbol interference Chapter 16 touched on conceptually.

```
Sharp pulse launched:        |
                              |
Multiple modes bounce at
different angles, arrive
at different times:      ____/\____/\___  <- smeared, harder to
                                              distinguish from the
                                              next pulse if sent too soon

This modal dispersion is why multi-mode fiber's reach is limited
to a few hundred meters at high speed, not kilometers.
```

**Where it's used:** multi-mode fiber is cheaper to manufacture and terminate, and pairs well with cheaper light sources (LEDs at lower speeds, VCSEL lasers at higher speeds) — making it the standard choice for **short runs inside a single data center or building**: connecting switches within or between adjacent racks, where distances are measured in tens or a few hundred meters, not kilometers.

---

## 5. Single-Mode Fiber — One Path, Long Reach

**Single-mode fiber solves modal dispersion by removing the "multi" entirely:** its core is dramatically narrower — typically **about 9 microns**, close to the wavelength of the light itself — small enough that light can only travel down it along essentially **one path**, straight through the center, with no alternate bounce angles physically able to propagate. With no competing modes traveling different effective path lengths, there's no modal dispersion to smear pulses over distance, and single-mode fiber can carry a sharp, clean signal for tens or even hundreds of kilometers before it needs amplification or regeneration.

**The trade-off:** launching light precisely into a 9-micron core requires a tightly focused **laser** light source (rather than multi-mode's cheaper LEDs), and the transceivers, connectors, and alignment tolerances involved are correspondingly more expensive and less forgiving. Single-mode fiber isn't "strictly better" in an unconditional sense — it's a deliberate trade of higher per-connection cost for dramatically longer reach and higher achievable bandwidth, which is exactly the trade that makes sense once you're running a cable between buildings, cities, or continents rather than between two racks.

```
                 Multi-mode (50-62.5 micron core)     Single-mode (9 micron core)
Light source:    LED or VCSEL (cheaper)                Laser (precise, pricier)
Modes:           Many simultaneous paths               Effectively one path
Dispersion:      Modal dispersion limits reach          Minimal — mostly chromatic
                                                          dispersion at very long distance
Typical reach:   Tens of meters to ~550m (at 10G+)      Kilometers to 100+ km
Typical use:     Inside a data center/building          Between buildings, cities, oceans
Cost:            Lower                                   Higher (optics + laser)
```

---

## 6. Why Fiber Beats Copper: The Physics, Made Concrete

Bringing Chapter 21's copper limitations and this chapter's fiber mechanics side by side:

```
                              Copper (Ch. 21)              Fiber (this chapter)
Signal carrier:               Electric current              Light (photons)
Vulnerable to EMI/crosstalk?  Yes — needs twisting,          No — glass is a dielectric,
                               shielding, differential          carries no current, radiates
                               signaling to manage it            and picks up no electrical
                                                                  interference at all
Attenuation over distance:    Meaningful even over            Extremely low — ~0.2-0.35 dB/km
                               100m; worsens sharply             for single-mode at the
                               at higher frequencies             wavelengths used in Section 9
Practical distance limit:     ~100m (Ch. 21, Section 8-9)     Tens to 100+ km without
                                                                  amplification (single-mode)
Bandwidth ceiling:            Set by attenuation +             Set mostly by the electronics
                               crosstalk margins at              at each end, not the fiber
                               a given cable category            itself — one strand can carry
                                                                  many terabits with WDM (Sec. 9)
```

The crosstalk immunity deserves its own sentence, because it's the most fundamental physical difference: **glass carries no electrical current and generates no electromagnetic field**, so two fiber strands can run bundled together with essentially none of Chapter 21's crosstalk concern — there's simply no electromagnetic coupling mechanism between them the way there unavoidably is between two current-carrying copper conductors. Fiber is also immune to external electromagnetic interference for the same reason — a fiber cable running next to industrial machinery, high-voltage power lines, or a lightning-prone area is entirely unaffected by the electrical noise that would devastate a copper link in the same location.

---

## 7. The Speed of Light in Glass — A Real Number, Worked

Light in a vacuum travels at the famous constant `c ≈ 299,792 km/s`. Inside glass, light travels slower, because the glass's refractive index (Section 2's `n1`) describes exactly how much the medium slows light down:

```
speed of light in a medium = c / n

For typical single-mode fiber glass, n ≈ 1.47-1.48

speed in fiber ≈ 299,792 km/s / 1.47 ≈ 203,940 km/s
                                     ≈ ~200,000 km/s (the commonly quoted round figure)
```

This is the real, physically-derived origin of the "~200,000 km/s" figure you'll see quoted throughout this course (and in the wider networking literature) for how fast a signal travels down fiber — it isn't an approximation of "close enough to light speed," it's the actual, calculable consequence of light slowing down by a factor equal to the glass's refractive index. Compare this to Chapter 21's copper figure of 60-70% of `c` (~180,000-210,000 km/s) — the two media are, perhaps surprisingly, in a similar ballpark for raw propagation speed. The dramatic real-world difference in achievable *distance* and *bandwidth* comes entirely from Section 6's attenuation and interference story, not from light being "faster" than electricity in any simple sense.

This ~200,000 km/s figure is also the exact number Chapter 23 uses to calculate real, concrete latencies for transatlantic and trans-Pacific undersea cables — worth remembering precisely for that reason.

---

## 8. Transceivers: SFP, SFP+, QSFP, and the Optics Ecosystem

Somewhere, an electrical signal inside a switch or router (Chapter 44) has to become a pulse of light to enter a fiber, and vice versa on the receiving end. That conversion happens inside a small, pluggable module called a **transceiver**, which slots into a cage on the front of networking equipment and lets the same physical switch port be reconfigured for different fiber types, distances, or even copper, just by swapping the module — a genuinely useful piece of modularity that saves organizations from needing a different switch model for every possible fiber scenario.

```
Electrical signal (inside the switch)
        |
        v
[ Transceiver: electrical <-> optical conversion, laser driver, photodetector ]
        |
        v
Optical signal --------- fiber cable --------- Optical signal
                                                       |
                                                       v
                                          [ Transceiver at far end ]
                                                       |
                                                       v
                                          Electrical signal (far switch)
```

| Form factor | Typical speed | Common use |
|---|---|---|
| SFP | 1 Gbps | Gigabit Ethernet uplinks, older/lower-speed fiber links |
| SFP+ | 10 Gbps | Standard 10-gigabit data center and enterprise uplinks |
| SFP28 | 25 Gbps | Modern data center server/switch links |
| QSFP+ | 40 Gbps | Data center switch-to-switch aggregation |
| QSFP28 | 100 Gbps | Core/spine data center links (Chapter 94's leaf-spine architecture) |
| QSFP-DD / OSFP | 400 Gbps+ | Hyperscale data center backbones, newest core routers |

("SFP" = Small Form-factor Pluggable; "Q" prefixes generally mean "quad," bundling 4 lanes of the base rate together — a QSFP+ at 40G is, internally, 4 lanes of 10G each, and a QSFP28 at 100G is 4 lanes of 25G each — the same underlying multi-lane trick, at progressively higher per-lane speeds.)

Transceivers are also specified by which wavelength and fiber type they expect — a transceiver built for **850nm** light is meant for **multi-mode** fiber over short distances, while one built for **1310nm** or **1550nm** is meant for **single-mode** fiber over long distances (Section 9 explains exactly why those specific wavelengths were chosen). Plugging a single-mode transceiver into multi-mode fiber (or vice versa) is a classic, very real troubleshooting mistake — the link will typically either fail to come up at all or perform far below spec, because the transceiver's laser and the fiber's core diameter and dispersion characteristics simply don't match what each was engineered for.

---

## 9. Wavelengths, WDM, and How Fiber Scales to Terabits

**Why specific wavelengths (850nm, 1310nm, 1550nm) and not just "any light"?** Glass fiber has "windows" — narrow bands of wavelength where its attenuation (Section 6) is at a local minimum, a direct consequence of the specific physical and chemical properties of the silica glass used. Engineers simply chose to build lasers and transceivers that operate inside those low-loss windows:

```
~850nm  window: used by multi-mode transceivers, moderate loss, cheap VCSEL lasers
~1310nm window: used by single-mode transceivers, lower loss, moderate dispersion
~1550nm window: used by single-mode transceivers, LOWEST loss (~0.2 dB/km) —
                 the window of choice for the longest submarine and backbone links (Ch. 23)
```

**The single biggest reason fiber's bandwidth ceiling is so much higher than copper's:** a single strand of fiber isn't limited to carrying just one signal on just one wavelength. **Wavelength Division Multiplexing (WDM)** — and its denser cousin, **DWDM (Dense WDM)** — sends many independent signals down the *same physical strand of glass simultaneously*, each one riding on its own slightly different wavelength (color) of light, the way a prism can separate white light into many distinct colors traveling together but distinguishable at the far end.

```
Single fiber strand, DWDM:

  λ1 (1550.12nm) -----> [10G, 100G, or 400G of data on this wavelength]
  λ2 (1550.92nm) -----> [another independent 10G/100G/400G channel]
  λ3 (1551.72nm) -----> [another independent channel]
  ...
  λ80+           -----> [potentially 80+ independent channels on ONE strand]

  All traveling down the SAME physical glass strand at once,
  separated and recombined by optical filters (not electronics)
  at each end.
```

A modern DWDM system can pack 80 or more separate wavelength channels onto one fiber strand, each carrying 100 Gbps or more — meaning a **single physical strand of glass**, no thicker than a hair, can realistically carry **tens of terabits per second**. This is the direct, concrete answer to why fiber's bandwidth ceiling is "set by the electronics, not the fiber itself" (Section 6's table): as transceiver electronics get faster and DWDM systems pack in more/denser wavelengths, the same buried or submarine glass keeps carrying more data, with no need to dig up and replace the fiber itself — a huge economic advantage for long-haul and submarine infrastructure (directly relevant to Chapter 23's undersea cable capacity figures).

---

## 10. Real-World Fiber Specs and Distance Limits

| Standard | Fiber type | Wavelength | Typical max distance |
|---|---|---|---|
| 1000BASE-SX | Multi-mode | 850nm | ~550 m |
| 1000BASE-LX | Single-mode | 1310nm | ~10 km |
| 10GBASE-SR | Multi-mode | 850nm | ~300 m (OM3) / ~400 m (OM4) |
| 10GBASE-LR | Single-mode | 1310nm | ~10 km |
| 10GBASE-ER | Single-mode | 1550nm | ~40 km |
| 100GBASE-LR4 | Single-mode | ~1310nm (4 lanes) | ~10 km |
| 100GBASE-ZR | Single-mode | 1550nm (DWDM-compatible) | ~80-120 km |
| Submarine long-haul systems (Ch. 23) | Single-mode | 1550nm, DWDM | Hundreds of km between amplifier huts, thousands of km total with repeaters |

The pattern across this table directly reflects Sections 4-6: multi-mode (SX/SR) standards top out at a few hundred meters because of modal dispersion, while single-mode (LX/LR/ER/ZR) standards stretch to tens or over a hundred kilometers because they've eliminated that limit — and the true intercontinental distances in Chapter 23 are only reachable because submarine systems add periodic **optical amplifiers** (roughly every 50-100 km) that boost the light signal directly, without ever converting it back to electricity, keeping the entire multi-thousand-kilometer path in the optical domain.

---

## 11. Fiber to the Home: PON and Why Your ISP Loves It

Everything so far has assumed one dedicated fiber strand per link, exactly like copper Ethernet's dedicated pair per device (Chapter 21). Residential fiber internet (marketed as "FTTH," Fiber to the Home) usually works differently, for a straightforwardly economic reason: running one dedicated strand of fiber, plus one expensive laser transceiver, all the way from an ISP's facility to every single home individually would be enormously costly at neighborhood scale.

The real solution deployed by virtually every residential fiber ISP is a **Passive Optical Network (PON)**. "Passive" is the key word: a single fiber strand leaves the ISP's facility and runs to a **passive optical splitter** — a small device with no power supply and no electronics at all, purely glass and optical geometry — that splits the one incoming beam of light into many (commonly 32 or 64) outgoing directions, each continuing on to a different home over its own final short fiber run:

```
ISP facility (OLT)
        |
        | one fiber strand
        v
  [ Passive optical splitter — no power, pure optics ]
    /      |      |      \
   /       |      |       \
Home 1   Home 2  Home 3  ... Home 32/64
(ONT)    (ONT)   (ONT)        (ONT)
```

Downstream traffic (ISP to home) is simply broadcast down the shared trunk to every home simultaneously, with each home's equipment (the **ONT — Optical Network Terminal**, the box mounted at the customer's premises) picking out only the traffic addressed to it — conceptually similar to how a hub, back in Chapter 30, broadcasts everything to every port and relies on each device to just ignore what isn't addressed to it. Upstream traffic (home to ISP) is more delicate, since all the homes sharing that one trunk strand would otherwise collide if they transmitted light back up it at the same time — PON standards solve this with time-division scheduling, where the ISP's equipment (the **OLT — Optical Line Terminal**) grants each home a specific, brief time slot to transmit, so their upstream light pulses never overlap on the shared glass.

The economic payoff is exactly what motivated this design: the ISP needs only one expensive laser transceiver and one fiber strand per 32-64 homes, not per home, dramatically lowering the cost of reaching a whole neighborhood with fiber — which is a large part of why residential gigabit fiber internet has become widely affordable over the last decade.

---

## 12. Hands-On: Reading a Transceiver and Testing a Fiber Link

**Reading a transceiver's physical label** tells you most of what Sections 8-9 covered at a glance — a typical SFP+ module is printed with something like `10GBASE-LR 1310nm 10km`, directly telling you the standard, wavelength, and rated single-mode distance, exactly matching the table in Section 10.

**On Linux**, you can often query optical transceiver diagnostics directly from supported NICs/switches:

```
$ ethtool -m eth0
        Identifier                               : 0x03 (SFP)
        Transceiver type                         : 10G Ethernet: 10G Base-LR
        Laser wavelength                         : 1310nm
        Length (SMF,km)                          : 10km
        Optical diagnostics support               : Yes
        Laser output power                        : 1.23 mW / 0.90 dBm
        Receiver signal average optical power      : 0.45 mW / -3.47 dBm
```

That "receiver signal average optical power" reading is the fiber-optic equivalent of the SNR discussion from Chapter 17 — a receiver power far below the transceiver's rated sensitivity threshold is a direct, measurable early warning sign of a dirty connector, a fiber bend radius violation, or a run exceeding its rated distance, well before the link starts throwing CRC errors (Chapter 19) or dropping altogether.

**Experiment to try (conceptually, or with real gear if available):** a common real-world fiber troubleshooting step is inspecting a connector end-face under a fiber microscope — a speck of dust or a scratch on the polished glass end can scatter enough light to push received power below spec, exactly the kind of physical-layer problem this chapter's total-internal-reflection story predicts (light escaping instead of staying trapped, because the fiber's end-face, not its length, is now the source of loss).

---

## 13. Common Misconceptions

- **"Fiber is just 'faster' than copper because light is faster than electricity."** Not really, and Section 7 showed the actual numbers: light in fiber travels at ~200,000 km/s, copper's electromagnetic signal at ~180,000-210,000 km/s — genuinely comparable. Fiber's real advantage is much lower attenuation and total immunity to electromagnetic crosstalk (Section 6), which is what lets it run vastly longer distances and carry vastly more bandwidth, not raw propagation speed.
- **"Single-mode fiber is strictly better, so always use it."** No — Section 5 was explicit that single-mode trades higher-cost, less-forgiving laser optics for longer reach. For a 10-meter rack-to-rack data center link, multi-mode's cheaper transceivers are usually the more sensible engineering choice; single-mode's advantages only pay for themselves once distance requirements exceed what modal dispersion allows.
- **"You can splice single-mode and multi-mode fiber together with no issue, since it's all just glass."** No — the core diameters differ by nearly an order of magnitude (9 microns vs. 50-62.5 microns), and mismatched cores at a splice or connector cause severe signal loss or an outright non-functional link; likewise, mismatching a transceiver's expected fiber type (Section 8) is a very real and common cause of "the link just won't come up."
- **"WDM/DWDM means running more physical fiber strands."** No — the entire point (Section 9) is the opposite: many independent signals share one physical strand of glass simultaneously, distinguished purely by wavelength, which is precisely what makes upgrading a long-haul or submarine route's capacity (Chapter 23) a matter of installing new endpoint electronics rather than digging up or re-laying cable.
- **"Fiber has no distance limit at all."** No — attenuation in fiber is extremely low but not zero (Section 6, ~0.2-0.35 dB/km), and other effects like chromatic dispersion (pulse spreading due to different wavelengths of light traveling at very slightly different speeds even within single-mode fiber) still accumulate over very long distances, which is exactly why submarine cable systems (Chapter 23) need periodic optical amplification every 50-100 km rather than running one unbroken multi-thousand-kilometer strand.

---

## 14. Interview Questions & Model Answers

**Q1 (Beginner): Explain total internal reflection in your own words, and why it's the reason light stays inside a fiber optic cable.**

"A fiber has two layers of glass — a core in the center and a cladding around it — engineered so light travels very slightly slower in the core than it would in the cladding, because the core has a higher refractive index. When a beam of light inside the core hits the boundary with the cladding at a shallow enough angle, instead of passing through and escaping like light through a window, it reflects 100% back into the core, as if the boundary were a perfect mirror — with no actual mirror coating involved. This is called total internal reflection, and it happens because of the specific angle and the difference in refractive index between the two materials. As long as light enters the fiber within the angle range that guarantees this reflection, it keeps bouncing down the length of the core, even around gentle bends, for kilometers, with next to no light escaping through the sides."

**Q2 (Intermediate): What's the actual difference between single-mode and multi-mode fiber, and why does it matter for how far you can run each one?**

"The difference is the diameter of the light-carrying core: multi-mode fiber has a wide core, around 50 or 62.5 microns, wide enough that light can travel down it along many different valid bounce angles, or 'modes,' at once. Single-mode fiber has a much narrower core, around 9 microns, small enough that only one path is physically possible. The problem with multiple modes is that they travel different effective path lengths down the same length of cable — a steeply-bouncing ray covers more actual distance than one going nearly straight down the middle — so a sharp pulse of light launched into multi-mode fiber arrives at the far end smeared out in time, a phenomenon called modal dispersion. That smearing limits both the maximum distance and the maximum data rate multi-mode fiber can reliably carry — typically a few hundred meters at high speed. Single-mode fiber has no competing modes to cause that smearing, so it can carry a clean signal for tens or over a hundred kilometers, at the cost of needing more precise, more expensive laser-based transceivers to launch light accurately into such a narrow core."

**Q3 (Advanced): A data center wants to upgrade a long-haul link's capacity from 100 Gbps to several terabits per second without laying new fiber. How is this actually done, and why does it work?**

"The answer is DWDM — Dense Wavelength Division Multiplexing. A single strand of single-mode fiber can carry many independent optical signals simultaneously, each modulated onto its own distinct wavelength of light within the fiber's low-loss transmission windows, most commonly around 1550nm for long-haul systems. Optical filters at each end of the fiber combine (multiplex) the many wavelength channels onto the one strand for transmission and separate (demultiplex) them again at the receiving end, entirely in the optical domain — no electronics have to touch the combined signal in between. A modern DWDM system can pack 80 or more separate wavelength channels onto one strand, each carrying 100 Gbps or more using modern transceivers, which is how a single physical strand — the same glass that was laid years earlier — can be upgraded to carry tens of terabits per second just by installing newer, denser DWDM endpoint equipment. This is precisely why long-haul and submarine fiber routes (Chapter 23) are typically upgraded in capacity over their multi-decade physical lifespan without ever re-laying the underlying cable."

---

## 15. Exercises

### Easy

1. In your own words, explain why a fiber's cladding needs a *lower* refractive index than its core, and what would happen if the two layers had the same refractive index.
2. Using Section 10's table, identify which standard you'd choose to connect two switches 300 meters apart within the same data center, and which you'd choose to connect two buildings 8 km apart across a campus.
3. Why does a transceiver's physical label (e.g., "10GBASE-LR 1310nm 10km") tell you both the fiber type it expects and the maximum distance it supports?

### Medium

4. Using the formula from Section 7 (`speed = c / n`), calculate the speed of light in a fiber with a slightly different refractive index of `n = 1.5`, and compare it to the commonly-quoted ~200,000 km/s figure.
5. Explain, using Section 4's modal dispersion concept, why simply "sending data faster" (higher bit rate) makes multi-mode fiber's usable distance get *shorter*, not stay the same. (Hint: think about how close together consecutive pulses are in time at a higher bit rate, versus how much the pulses spread out over a given distance.)
6. A technician accidentally plugs a multi-mode 850nm transceiver into a single-mode fiber run. Using Sections 4, 5, and 8, predict what will likely happen to the link, and explain why.

### Hard

7. Using Section 9's DWDM concept and Section 10's real per-wavelength speeds (e.g., 100 Gbps per channel), calculate the total capacity of a fiber strand using 96 DWDM channels at 100 Gbps each, and compare that figure to a single strand's capacity without DWDM (one wavelength, one 100 Gbps channel).
8. Chapter 21 showed copper's speed of ~180,000-210,000 km/s is comparable to fiber's ~200,000 km/s (Section 7), yet real-world Ethernet distance limits differ by roughly 1000x (100m for copper vs. 100km+ for single-mode fiber). Write a short technical explanation (aimed at a junior engineer) of exactly which physical factors — not propagation speed — account for that 1000x gap, citing specific concepts from both this chapter and Chapter 21.
9. Research (or reason from Section 6's attenuation figures) how many optical amplifier sites a 6,600 km transatlantic submarine cable (a real, approximate distance covered in Chapter 23) would need, assuming amplifiers are spaced roughly every 70-100 km along the route, and explain why amplifying the signal optically (without converting to electricity) at each site is preferable to fully regenerating it electronically.

---

## Summary

| Term | Meaning |
|---|---|
| Total internal reflection (TIR) | Light reflecting perfectly at the core-cladding boundary due to differing refractive indices, no mirror coating needed |
| Core / cladding | The light-carrying center of a fiber and the surrounding layer that traps light via TIR |
| Multi-mode fiber | Wide core (50-62.5 microns); many light paths cause modal dispersion; short reach (hundreds of meters) |
| Single-mode fiber | Narrow core (~9 microns); effectively one light path; long reach (tens to 100+ km) |
| Modal dispersion | Pulse smearing caused by different light paths (modes) arriving at different times |
| Refractive index | A measure of how much a material slows down light; determines TIR and propagation speed |
| SFP / SFP+ / QSFP | Pluggable optical transceiver form factors for 1G / 10G / 40G+ links |
| WDM / DWDM | Sending many independent signals down one fiber strand simultaneously on different wavelengths |
| ~200,000 km/s | The real, calculated speed of light in typical single-mode fiber glass (c divided by ~1.47-1.48) |
| Optical amplifier | A device that boosts a fiber signal purely optically, without converting to electricity, used every 50-100 km on long-haul/submarine routes |

Fiber solves distance and bandwidth — but it's still a physical cable, and cables can't reach a phone in your pocket or a ship at sea, and someone still has to lay them across every ocean on Earth. Chapter 23 covers the two remaining pieces of the physical layer story: how wireless links (radio, microwave, satellite) carry data with no cable at all, and the physical reality — reinforced with real routes and real distances — that the overwhelming majority of the world's intercontinental internet traffic still travels through exactly the kind of single-mode fiber this chapter just explained, laid across the ocean floor.
