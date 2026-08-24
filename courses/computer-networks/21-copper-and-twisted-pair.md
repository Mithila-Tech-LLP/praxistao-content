# Chapter 21: Copper and Twisted Pair — Ethernet Cabling Explained

> *Every Ethernet cable running behind your desk, under your office floor, or through your walls is, underneath the plastic jacket, four pairs of ordinary copper wire twisted around each other. That twisting is not decoration and not accidental — it is the single mechanical trick that makes it possible to run gigabit signals down cheap copper without every pair jamming every other pair with noise.*

---

## Table of Contents

1. [The Problem: Copper Is a Shared, Leaky Medium](#1-the-problem-copper-is-a-shared-leaky-medium)
2. [How an Electrical Signal Actually Travels Down a Wire](#2-how-an-electrical-signal-actually-travels-down-a-wire)
3. [Crosstalk — The Naive Cable's Fatal Flaw](#3-crosstalk--the-naive-cables-fatal-flaw)
4. [Why Twisting the Pair Cancels Crosstalk](#4-why-twisting-the-pair-cancels-crosstalk)
5. [Differential Signaling — The Other Half of the Trick](#5-differential-signaling--the-other-half-of-the-trick)
6. [Inside a Cat Cable: Cat5e, Cat6, Cat6a, and Beyond](#6-inside-a-cat-cable-cat5e-cat6-cat6a-and-beyond)
7. [The RJ45 Connector and Pinout](#7-the-rj45-connector-and-pinout)
8. [Real Distance and Speed Limits](#8-real-distance-and-speed-limits)
9. [Attenuation, Length, and the 100-Meter Rule](#9-attenuation-length-and-the-100-meter-rule)
10. [Power over Ethernet — Copper Carries More Than Data](#10-power-over-ethernet--copper-carries-more-than-data)
11. [Auto-Negotiation — How Two Devices Agree on Speed](#11-auto-negotiation--how-two-devices-agree-on-speed)
12. [Hands-On: Reading a Cable's Real Specs](#12-hands-on-reading-a-cables-real-specs)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#summary)

---

## 1. The Problem: Copper Is a Shared, Leaky Medium

Chapters 14–20 built up the *logical* view of a physical link: bits become signals (Ch. 14), signals get modulated (Ch. 15), bandwidth and noise set hard limits (Ch. 16–18), and error detection/correction (Ch. 19–20) cope with whatever corruption still slips through. This chapter and the next two ask a more concrete question: **what is the cable itself, physically, and why is it built the way it is?**

The obvious, naive design for a network cable is: one wire per signal, running in parallel, like a ribbon cable inside an old PC. It's the simplest possible thing that could work, and — this is the point of this chapter — **it doesn't work well at any real distance or speed**, for a reason that has nothing to do with the wire's resistance and everything to do with electromagnetism: **any two current-carrying wires running near each other for any distance will induce a signal on each other.** A cable is never really an isolated pipe for one signal; it's always, unavoidably, a little bit of a shared, leaky medium with its electrical neighbors. Understanding *why* that happens — and the one mechanical fix that solves it — is this chapter's core.

---

## 2. How an Electrical Signal Actually Travels Down a Wire

**Intuitive picture first, and where it's wrong.** Most people imagine a signal "traveling down a wire" the way water flows through a pipe — electrons at one end pushing electrons at the other end, all the way down, at the speed the electrons themselves move. That's the wrong mental model, and the actual speed gives it away: individual electrons in a copper wire drift at only millimeters per second under normal signaling currents, yet a signal change at one end of a cable reaches the other end in nanoseconds. Something else is doing the traveling.

**What's actually happening:** a change in voltage at the sending end creates a changing electric and magnetic field that propagates *along and around* the conductor as an electromagnetic wave — guided by the copper and the insulating dielectric around it, much like a wave guided along the surface of a pond, not "pushed" through it like water in a pipe. The electrons barely move; the *field* is what races down the cable, at a meaningful fraction of the speed of light.

**How fast, concretely:** the propagation speed in copper Ethernet cable is typically **60–70% of the speed of light in a vacuum** (roughly 180,000–210,000 km/s), a figure called the cable's **velocity of propagation** or **nominal velocity of propagation (NVP)**, determined by the dielectric (insulating) material surrounding the conductor. This is why, in Chapter 28's discussion of Ethernet frame timing and in real cable-length testing (Section 12), engineers use a specific NVP figure (often around 0.64–0.72 depending on the exact cable) rather than assuming the vacuum speed of light.

```
Speed of light in vacuum:              ~300,000 km/s
Typical copper Ethernet cable (NVP):   ~180,000-210,000 km/s (60-70% of c)
Compare: light in glass fiber (Ch. 22): ~200,000 km/s (about 67% of c, for different physical reasons)
```

**The engineering signal itself:** modern Ethernet doesn't send a simple "high voltage = 1, low voltage = 0" pattern; it uses encoding schemes (like the multi-level PAM encoding introduced conceptually back in Chapter 15's discussion of modulation) that pack more information into each voltage transition, and it applies that voltage as a small, carefully controlled difference between two wires — which brings us to the actual physical problem this chapter is about.

---

## 3. Crosstalk — The Naive Cable's Fatal Flaw

**The naive attempt:** bundle several single wires together in a cable — one for each of the multiple simultaneous signals modern Ethernet uses (a real Cat cable carries 4 separate wire pairs, all active at once, to reach gigabit and multi-gigabit speeds). Run them close together, in parallel, for convenience and cost.

**Why this fails:** every wire carrying a changing current generates a changing magnetic field around it (basic electromagnetism — this is not an engineering flaw, it's physics). A changing magnetic field, in turn, **induces** a voltage in any nearby conductor — including the wire right next to it in the bundle. The signal on wire A "leaks" a faint, distorted copy of itself onto wire B, and vice versa. This leakage is called **crosstalk**, and it comes in two flavors that matter for real cable testing and troubleshooting:

- **NEXT (Near-End Crosstalk):** the leaked signal measured at the *same end* of the cable where the interfering signal originates — the interference and the desired signal are both strong here, so NEXT is often the harder problem to control.
- **FEXT (Far-End Crosstalk):** the leaked signal measured at the *opposite end* of the cable from where the interference originates — by the time it arrives, the interfering signal has also been attenuated by the cable's length, which changes its practical impact.

**Why this matters concretely:** in a naive parallel-wire cable carrying four simultaneous pairs (as real Ethernet does), each pair's data is a source of noise for the other three. Left unaddressed, this crosstalk would set a hard ceiling on how fast, or how far, you could reliably run multiple pairs through one cable jacket — you'd effectively be back to Chapter 17's noise/SNR problem, self-inflicted by your own cable design.

```
       Wire A (signal): ----->---->---->---->
                          |  induced   |
                          v  crosstalk v
       Wire B (signal): ----->---->---->---->
                        (Wire B's receiver now sees its own signal
                         PLUS a faint, distorted echo of Wire A's)
```

---

## 4. Why Twisting the Pair Cancels Crosstalk

**The real solution, and it's mechanical, not electronic:** instead of running two wires of a pair straight and parallel, **twist them around each other** at a tight, regular interval (real Cat cables use a different twist rate for each of the 4 pairs inside one jacket — more on that in Section 6). This is the single most important physical fact about every Ethernet cable you've ever plugged in.

**Intuitive picture:** imagine two people shouting insults at each other while walking side by side in a straight line next to another pair of people also walking in a straight line — everyone on the left side of the group is closer to the same neighbor the entire way, so the interference is consistently one-directional and adds up. Now imagine both pairs constantly crisscrossing positions, swapping "left" and "right" every few steps, in a synchronized weave — over any reasonably long stretch, each person spends equal time near the interfering pair's noisiest and quietest points, and the *net* interference picked up by "the pair" as a whole tends to cancel out.

**The actual electromagnetic mechanism:** because the two wires of a twisted pair constantly swap physical position relative to any external noise source (whether that's another pair in the same cable, or outside interference like a fluorescent light ballast), each wire spends roughly equal time closer to and farther from that noise source. **A given noise source therefore induces very nearly the same amount of interference on both wires of the pair — and Section 5 shows exactly why "the same interference on both wires" is precisely the failure mode differential signaling is immune to.** The interference doesn't stop being generated (physics doesn't care about the twist), but by the time the receiver looks at the actual signal it needs, the twist has converted "unequal interference that corrupts the message" into "equal, common interference that cancels."

**Where this analogy is incomplete:** twisting doesn't eliminate every source of crosstalk — the four pairs inside one cable jacket are still close enough together, and their twist rates are still finite, that some residual coupling remains (which is exactly why Section 6's cable categories are rated with hard NEXT/FEXT limits in decibels, not "zero crosstalk"). Twisting reduces crosstalk to a manageable, specified level; it doesn't make it physically zero.

---

## 5. Differential Signaling — The Other Half of the Trick

Twisting alone would do nothing if a receiver only looked at one wire's voltage relative to a fixed ground reference — a noise spike hitting that one wire would still corrupt the reading. The mechanism that makes twisting actually pay off is **differential signaling**: instead of sending a signal as "this wire's voltage relative to ground," Ethernet sends it as **the voltage difference between the two wires of a pair**, and the receiver reads *only that difference*.

```
Sender puts:   Wire A = +0.5V,  Wire B = -0.5V     → difference = +1.0V → "1"
Sender puts:   Wire A = -0.5V,  Wire B = +0.5V     → difference = -1.0V → "0"
```

Now combine this with Section 4's result: a nearby noise source induces very nearly the **same** extra voltage on both wire A and wire B (this is what twisting guarantees). If a noise spike adds, say, +0.2V to both wires:

```
Without twisting (noise hits unevenly, e.g. +0.2V on A only):
   Wire A = 0.5 + 0.2 = 0.7V,  Wire B = -0.5V   → difference = 1.2V  (signal corrupted!)

With twisting (noise hits both wires nearly equally, +0.2V on both):
   Wire A = 0.5 + 0.2 = 0.7V,  Wire B = -0.5 + 0.2 = -0.3V
   difference = 0.7 - (-0.3) = 1.0V   → unchanged! Noise cancels out of the DIFFERENCE.
```

This is called **common-mode noise rejection**, and it's the actual payoff of twisting: the receiver's differential circuit subtracts the two wires, and any noise that hit both wires equally — which is exactly what twisting engineers the noise to do — vanishes in that subtraction, leaving the real signal untouched. Twisting and differential signaling are a matched pair of techniques; neither one alone solves the crosstalk problem as effectively as the two combined.

---

## 6. Inside a Cat Cable: Cat5e, Cat6, Cat6a, and Beyond

A standard Ethernet cable — technically an **"8P8C" connector cable**, universally (if technically incorrectly) called "RJ45" — contains **4 twisted pairs** (8 wires total), each pair twisted at its own distinct rate, precisely to reduce crosstalk *between* the 4 pairs as well as from outside sources (if every pair used the same twist rate, they could still line up periodically and re-create the very interference twisting is meant to prevent).

```
                    ___________________________________
                   /   Pair 1 (twisted, own rate)       \
                  |    Pair 2 (twisted, different rate)   |
   Cable jacket → |    Pair 3 (twisted, different rate)   |
                  |    Pair 4 (twisted, different rate)   |
                   \___________________________________/
```

| Category | Max bandwidth | Max speed | Max distance (at rated speed) | Notes |
|---|---|---|---|---|
| Cat5e | 100 MHz | 1 Gbps | 100 m | Enhanced Cat5; improved crosstalk specs; still extremely common |
| Cat6 | 250 MHz | 1 Gbps (100m) / 10 Gbps (up to ~55m) | 100 m (1G) / ~37-55 m (10G) | Tighter twists, often a physical spline separating pairs |
| Cat6a | 500 MHz | 10 Gbps | 100 m | Thicker, better-shielded; full 10G at full 100m distance |
| Cat7 | 600 MHz | 10 Gbps+ | 100 m | Individually shielded pairs (S/FTP); requires non-standard connectors for full spec |
| Cat8 | 2000 MHz | 25/40 Gbps | 30 m | Designed for short data-center runs, not general office cabling |

The pattern across this table is exactly what Sections 4–5 predicted: **higher bandwidth (Chapter 16) demands tighter crosstalk control**, because higher-frequency signals both radiate more energy to their neighbors and are more sensitive to what they pick up. Cat6's tighter twist rates and often-added internal spline, and Cat6a/7/8's added metallic shielding around individual pairs, are direct, physical engineering responses to needing more bandwidth (and therefore, per Shannon's limit from Chapter 18, more throughput) out of the same basic twisted-copper idea.

**A crucial, commonly-missed nuance about Cat6 in the table above:** unlike Cat5e and Cat6a, Cat6's 10-gigabit capability is **distance-limited** to well under the standard 100-meter run (typically 37–55 meters depending on cable quality and external interference) — this is a genuinely common source of real-world confusion when someone installs "Cat6" cabling expecting guaranteed 10Gbps at full building-wiring distances and doesn't get it.

**Shielding variants — the letters before and after "TP".** Beyond the Cat number, cables are also classified by how much metallic shielding they add on top of twisting, using a `XX/YTP` naming scheme (`U` = unshielded, `F` = foil, `S` = braided screen, `TP` = twisted pair):

```
UTP   (Unshielded Twisted Pair)         - relies purely on twisting; cheapest, most common,
                                           fine for Cat5e/Cat6 in typical office environments
F/UTP (overall Foil-shielded, UTP inside) - one foil wrap around all 4 pairs together;
                                             common on Cat6a for extra crosstalk/EMI margin
S/FTP (overall braided Screen,
       Foil-shielded pairs individually)  - each pair individually foil-wrapped, PLUS an
                                             overall braided shield; typical of Cat7/Cat8,
                                             needed for their much higher frequency ratings
```

Shielded cable variants require the shield to be properly and continuously **grounded** at both ends to be effective — an ungrounded or partially-grounded shield can, counterintuitively, perform *worse* than no shield at all, by acting as an antenna that picks up and re-radiates interference rather than blocking it. This is a real, common installation mistake worth remembering: shielding is only as good as its grounding.

---

## 7. The RJ45 Connector and Pinout

The connector on the end of every Ethernet cable is technically an **8P8C (8 Position, 8 Contact)** modular connector, universally called "RJ45" even though that name, strictly, refers to a similar-looking but electrically different telephone connector standard. Two standardized wiring orders exist for arranging the 8 wires into the 8 pins, **T568A** and **T568B** — functionally equivalent, but a cable must be wired consistently at both ends (a "straight-through" cable) or deliberately swapped (a "crossover" cable, mostly obsolete now that Auto-MDI/MDI-X, described below, handles this automatically):

```
Pin:        1     2     3     4     5     6     7     8
T568B:    W-Or   Or   W-Grn  Blu  W-Blu  Grn  W-Brn  Brn
T568A:    W-Grn  Grn  W-Or   Blu  W-Blu  Or   W-Brn  Brn

("W-Or" = white/orange striped wire, etc. — the standard color-coded
 wires inside every Cat5e/6/6a cable.)
```

For 1000BASE-T (Gigabit Ethernet) and faster, **all 4 pairs are used simultaneously, bidirectionally** — a major shift from 10/100 Mbps Ethernet, which only used 2 of the 4 pairs (pins 1,2 and 3,6) and left the other two pairs unused. This is why a cable that works fine at 100 Mbps can sometimes fail at 1000 Mbps: a fault in a previously-unused pair only becomes visible once gigabit signaling starts using all four.

**Auto-MDI/MDI-X**, standard on essentially all modern Ethernet hardware, lets a device automatically detect whether it's talking to another end device or to a switch/hub and electrically swap its transmit/receive pairs in software — this is why "crossover cables" (once required to directly connect two PCs, or two switches, without an intervening device) are essentially unnecessary today; any standard straight-through cable works in any port.

---

## 8. Real Distance and Speed Limits

The **100-meter maximum run** quoted throughout Section 6 is not an arbitrary round number — it comes from two independent, compounding physical constraints:

1. **Attenuation and crosstalk margins (Section 9):** the cable's signal degrades and picks up noise as it gets longer; the various Cat standards are engineered and tested to guarantee their rated bandwidth and speed up to, but not meaningfully beyond, 100 meters.
2. **Historical timing constraints (half-duplex CSMA/CD, Chapter 30):** original Ethernet needed every device on a shared segment to be able to detect a collision before it finished sending the smallest legal frame — which put a hard cap on how large (in propagation delay) a shared collision domain could be. Modern full-duplex switched Ethernet (Chapter 30) has removed collision domains from the picture almost entirely, but the 100m figure that came out of these original timing budgets stuck as the standard because it also happens to line up well with Section 9's attenuation limits.

```
Real-world numbers to anchor intuition:

  100BASE-TX (Fast Ethernet):     100 m max, well within Cat5e's specs
  1000BASE-T (Gigabit):           100 m max, Cat5e/6/6a all support this
  10GBASE-T:                      100 m on Cat6a/7/8; only ~37-55 m on Cat5e/Cat6
  25GBASE-T / 40GBASE-T:          30 m on Cat8 — short-reach, data-center-only
```

Beyond ~100 meters on copper, the standard answer is not "a longer/thicker copper cable" but either an **Ethernet repeater/extender**, a **media converter to fiber** (Chapter 22 explains exactly why fiber doesn't share copper's distance ceiling), or restructuring the network topology so no single copper run needs to be that long (e.g., placing a switch closer to the far end).

---

## 9. Attenuation, Length, and the 100-Meter Rule

Recall Chapter 17's definition of attenuation: signal strength decreasing with distance, generally worse at higher frequencies. Copper twisted pair is a textbook case: the cable's own electrical resistance and the "skin effect" (high-frequency current tends to travel closer to the surface of the conductor, effectively increasing resistance at higher frequencies) both cause more signal loss per meter as frequency — and therefore data rate — goes up. This is precisely why Cat6a's higher-frequency (500 MHz) signal needs thicker conductors and better shielding than Cat5e's 100 MHz signal to reach the same 100-meter distance without the signal decaying below a usable SNR (Chapter 17) at the receiver.

Cable testing standards specify this precisely as **insertion loss** (how much the signal weakens over the cable's length, in dB) and pair it with the **NEXT/FEXT** crosstalk figures from Section 3 into a combined metric, **ACR (Attenuation-to-Crosstalk Ratio)**, which is really just Chapter 17's SNR concept applied specifically to a twisted-pair cable: how much stronger is my actual signal, after weakening over distance, than the crosstalk noise it's competing against? A cable that passes its Cat6a certification is one where, at every point along a 100-meter run, that ratio stays above the threshold Chapter 18's Shannon-Hartley math says is needed to sustain 10 Gbps reliably.

**Real approximate insertion loss figures**, to make "attenuation gets worse at higher frequency" concrete rather than abstract:

```
Cat5e, measured at 100 MHz, 100 m run:   ~22-24 dB of signal loss
Cat6,  measured at 250 MHz, 100 m run:   ~19.8-21 dB of signal loss (better cable, but tested at higher frequency)
Cat6a, measured at 500 MHz, 100 m run:   ~20-21 dB of signal loss (again, better cable offsetting a much higher test frequency)

(dB loss is logarithmic: -20 dB means the signal power arriving at
the far end is only about 1% of what was transmitted — and the
receiver still has to pull a clean signal out of that, in the
presence of whatever crosstalk and external noise remains.)
```

The pattern is not "higher categories have less attenuation" in some absolute sense — it's that each category is engineered to hold attenuation roughly *flat* even as its rated test frequency (and therefore usable bandwidth, per Chapter 16) climbs dramatically, which is the real engineering achievement each Cat generation represents.

---

## 10. Power over Ethernet — Copper Carries More Than Data

A practical bonus of copper (something fiber, per Chapter 22, structurally cannot do): the same twisted-pair wires carrying a data signal can simultaneously carry **DC electrical power**, a standard called **Power over Ethernet (PoE)**. This works because the differential data signal (Section 5) rides on top of a DC voltage applied across the same wire pairs — the receiving device's circuitry separates the DC power from the AC-like data signal using simple filtering, since they occupy very different, easily-separable frequency ranges.

```
IEEE 802.3af (PoE):     up to ~15.4W delivered   — IP phones, basic access points
IEEE 802.3at (PoE+):    up to ~30W delivered     — pan-tilt-zoom cameras, better APs
IEEE 802.3bt (PoE++):   up to ~60-100W delivered — laptops, high-power access points
```

This is why a single Cat5e/6 cable can power and network a Wi-Fi access point mounted on a ceiling, or a security camera, with no separate electrical wiring run needed at all — a real, practical advantage copper retains even as fiber (Chapter 22) wins on raw distance and bandwidth.

---

## 11. Auto-Negotiation — How Two Devices Agree on Speed

A question this chapter has left implicit: when you plug a gigabit-capable laptop into a switch, how do the two ends agree to actually run at 1000 Mbps rather than falling back to some slower, safer speed? The answer is **auto-negotiation (IEEE 802.3 Clause 28)**, and it happens entirely at the physical layer, before any Ethernet frame (Chapter 28) is ever exchanged.

**The mechanism:** each end sends a burst of precisely-timed voltage pulses called **Fast Link Pulses (FLPs)** — a direct descendant of the simple **Normal Link Pulses (NLPs)** older 10BASE-T equipment used just to confirm "someone's alive on the other end of this cable." An FLP burst encodes a bitmap of every speed and duplex mode the sending device supports (10/100/1000 Mbps, half/full duplex, flow control, and so on). Both ends exchange these bitmaps and independently apply the same tie-breaking rule — **the highest commonly-supported speed and duplex setting wins** — so both sides land on the same conclusion without needing a back-and-forth negotiation protocol layered on top.

```
Device A advertises:  10H, 10F, 100H, 100F, 1000F
Device B advertises:  10H, 10F, 100H, 100F
                          (B has no gigabit hardware)

Highest mode BOTH support: 100F (100 Mbps, full duplex)
Both ends independently select 100F. Link comes up at 100 Mbps.
```

**Why this matters for real troubleshooting:** a link stuck at a surprisingly low speed or running in half-duplex is very often an auto-negotiation mismatch — historically caused by one end being manually forced to a fixed speed/duplex while the other is left on auto-negotiate, a classic and still-common cause of the "duplex mismatch" symptom (very slow, error-riddled throughput despite an apparently "up" link) that shows up in real-world network debugging (previewed further in Chapter 122's debugging playbook). The fix is almost always: set both ends to auto-negotiate, or manually fix both ends to *identical* settings — never one auto, one fixed.

---

## 12. Hands-On: Reading a Cable's Real Specs

**On Linux**, check what speed and duplex your Ethernet link actually negotiated (this reflects both the cable's real category and the electronics on both ends):

```
$ ethtool eth0
Settings for eth0:
        Supported link modes:   10baseT/Half 10baseT/Full
                                 100baseT/Half 100baseT/Full
                                 1000baseT/Full
        Speed: 1000Mb/s
        Duplex: Full
        Port: Twisted Pair
        Link detected: yes
```

**On macOS**, similar information is available via:

```
$ networksetup -getMedia "Ethernet"
```

**Physically identifying a cable's category** is usually printed directly on the cable jacket itself — look for text like `CAT6 UTP` (Unshielded Twisted Pair) or `CAT6A F/UTP` (foil-shielded overall, unshielded individual pairs) printed repeatedly along the cable's length. A cable tester (inexpensive ones are widely available) can verify continuity on all 8 conductors and, on better models, measure real length using the NVP figure from Section 2 (send a pulse down the cable, time how long it takes to reflect back off the far end or an open circuit, and multiply by NVP/2 — the same "time-of-flight" principle behind radar and, as Chapter 54 will show, `traceroute`).

**Experiment to try:** run `ethtool -S eth0` (Linux) or check a managed switch's port statistics, and look for counters like CRC errors, or (on gigabit links) counters distinguishing which of the 4 pairs reports errors — modern switch chips can often report per-pair signal quality, which is a direct, observable trace of exactly the crosstalk and attenuation physics this chapter described.

---

## 13. Common Misconceptions

- **"A higher category cable is always strictly better, so just buy the most expensive one."** Not quite — a Cat6a cable run for a 100 Mbps device gains nothing functionally (the bottleneck is the device, not the cable), while adding cost and, for shielded variants, a stiffer cable that's harder to route. Match the cable category to the actual speed you need, with reasonable headroom for future upgrades.
- **"Twisting eliminates crosstalk completely."** No — Section 4 was explicit that twisting reduces crosstalk to a specified, tested level (the NEXT/FEXT limits each Cat standard is certified against), not to physical zero. Untwisted "silver satin" style flat cables (rare today, but historically used for telephone wiring repurposed for data) perform dramatically worse precisely because they skip this mechanism.
- **"RJ45 is the technically correct name for the connector."** Strictly, "RJ45" refers to a telephone registered-jack standard with a different pin/key arrangement; the connector on Ethernet cables is properly called 8P8C. The "RJ45" name stuck through decades of common usage and is universally understood, but it's worth knowing the distinction exists.
- **"Gigabit Ethernet only needs 2 of the 4 pairs, like 100 Mbps did."** No — Section 7 was explicit that 1000BASE-T uses all 4 pairs simultaneously and bidirectionally on each pair; a cable or wiring job with only 2 good pairs (common in old telephone-repurposed wiring) will not reliably reach gigabit speeds even if it worked fine at 100 Mbps.
- **"The 100-meter limit is just a conservative safety margin and can usually be stretched."** No — it's derived from real, specified attenuation and crosstalk limits (Section 9); exceeding it doesn't fail gracefully by "going a bit slower," it typically produces intermittent CRC errors (Chapter 19) and retransmissions that are far more frustrating to diagnose than an outright failure would be.

---

## 14. Interview Questions & Model Answers

**Q1 (Beginner): Why are the wires inside an Ethernet cable twisted instead of running straight and parallel?**

"Any two current-carrying wires running near each other induce a voltage on each other through their changing electromagnetic fields — this is called crosstalk, and it's unavoidable physics, not a manufacturing defect. If the wires ran straight and parallel, one wire would consistently sit closer to a given noise source (whether that's another pair in the same cable or outside interference) than the other, so the interference would land unevenly and corrupt the signal. Twisting the two wires of a pair around each other makes them constantly swap relative position to any nearby noise source, so both wires end up picking up very nearly the same amount of interference. Combined with differential signaling — where the receiver only cares about the *difference* in voltage between the two wires, not their absolute levels — that equal, 'common-mode' interference cancels out in the subtraction, leaving the real signal intact."

**Q2 (Intermediate): A user reports their Cat6 cable 'only gets 1 Gbps, not 10 Gbps, even though it's rated for 10G.' What's actually going on?**

"Cat6's 10-gigabit rating is distance-limited, unlike Cat5e (which never supports 10G) or Cat6a (which supports 10G at the full standard 100-meter run). Cat6 can typically only sustain 10GBASE-T reliably up to somewhere around 37 to 55 meters, depending on cable quality, bundling with other cables, and ambient electrical interference — well short of the 100-meter figure most people associate with Ethernet cabling. If the run in question is longer than that, the link will either fail to establish 10G at all and fall back to 1G (which is almost certainly what's happening here), or in some marginal cases run at 10G with a much higher error rate. The fix is either shortening the run, upgrading to Cat6a for that segment, or accepting the 1G fallback if the distance can't be changed."

**Q3 (Advanced): Explain, at the level of the underlying physics, why higher Cat standards (Cat6a, Cat7, Cat8) increasingly rely on shielding rather than just tighter twisting.**

"Twisting exploits the fact that a nearby noise source induces roughly equal interference on both wires of a pair, which differential signaling then cancels out (Sections 4-5). But that cancellation is never perfect — some residual coupling always remains, and critically, the *amount* of energy a signal radiates to its neighbors, and the amount of external interference a wire picks up, both scale up as frequency increases (this connects directly to Chapter 16's frequency/bandwidth material — higher-frequency signals are inherently 'louder' electromagnetically for a given size of conductor). As data rates climb from Cat5e's 100 MHz to Cat6a's 500 MHz and beyond, tighter twisting alone increasingly fails to keep NEXT/FEXT within the levels needed to sustain Chapter 18's Shannon-limit-derived SNR requirements for the target bit rate. Adding a metallic shield — around each individual pair (as in Cat7's S/FTP construction) or around the whole cable — provides a physical barrier that blocks electromagnetic coupling directly, rather than relying purely on geometric cancellation, which is why the highest-speed copper standards (Cat7, Cat8) essentially always specify shielding as well as tight, precisely engineered twist rates."

---

## 15. Exercises

### Easy

1. Explain in your own words why a copper Ethernet signal reaches the far end of a 100-meter cable in roughly 500 nanoseconds, even though the individual electrons inside the wire are moving at only millimeters per second.
2. List the four wire pairs' color codes in the T568B standard (Section 7) and identify which pins they occupy.
3. Why can a single Cat5e/6 cable both power a device (PoE) and carry its network data simultaneously, without the power and data signals interfering with each other?

### Medium

4. A network installer runs a Cat6 cable exactly 70 meters between a switch and an access point and expects a reliable 10 Gbps link. Using Section 6 and 8's specifications, explain why this expectation is likely to fail, and propose two different fixes.
5. Using the common-mode noise cancellation math worked out in Section 5, show what happens to the received signal if a noise spike adds +0.3V to wire A but only +0.1V to wire B (i.e., the twisting/cancellation is imperfect). How much residual noise leaks into the final differential reading?
6. Explain why 100 Mbps Ethernet (100BASE-TX) only needs 2 of the 4 wire pairs, but Gigabit Ethernet (1000BASE-T) needs all 4, connecting your answer to Chapter 16's bandwidth/throughput relationship.

### Hard

7. Research (or reason from first principles using Chapter 17's attenuation concepts) why "skin effect" causes copper's effective resistance to increase at higher frequencies, and explain why this specifically works against the higher data rates that Cat6a/7/8 are trying to achieve, beyond just crosstalk considerations.
8. A data center wants to run 25GBASE-T over copper between racks 20 meters apart. Using Section 6's Cat8 specifications and Section 9's ACR (Attenuation-to-Crosstalk Ratio) concept, explain what physical cable properties had to improve relative to Cat6a to make this possible at all, and why Cat8 is explicitly not marketed for general office/building cabling despite its higher speed rating.
9. Compare the propagation speed figures for copper twisted pair (Section 2, ~60-70% of c) and fiber optic cable (Chapter 22, ~67% of c, ~200,000 km/s) — they're surprisingly close. Given that, explain in your own words why fiber still dramatically outperforms copper for long-haul links, if raw signal speed isn't the deciding factor. (Hint: think about what limits maximum cable length in each medium, from Sections 8-9 here versus what Chapter 22 will describe for fiber.)

---

## Summary

| Term | Meaning |
|---|---|
| Velocity of propagation (NVP) | Speed an electromagnetic signal actually travels down a cable; ~60-70% of light speed in copper |
| Crosstalk (NEXT/FEXT) | Unwanted signal induced from one wire pair onto another, measured at the near/far end |
| Twisted pair | Two wires twisted together so external noise induces nearly equal interference on both |
| Differential signaling | Encoding data as the voltage *difference* between two wires, immune to common-mode noise |
| Common-mode noise rejection | The cancellation of equal interference on both wires when only their difference is read |
| Cat5e / Cat6 / Cat6a | Twisted-pair cable standards with increasing bandwidth (100/250/500 MHz) and speed capability |
| 8P8C ("RJ45") | The 8-position, 8-contact modular connector used on virtually all Ethernet cables |
| T568A / T568B | The two standardized wire-to-pin orderings for wiring an 8P8C connector |
| ACR (Attenuation-to-Crosstalk Ratio) | Chapter 17's SNR concept applied specifically to twisted-pair cable certification |
| Power over Ethernet (PoE) | Delivering DC power over the same twisted pairs that carry data |

Twisting and differential signaling solve crosstalk well enough for copper to reach 10, even 40 Gbps — but only over short distances, and only by fighting attenuation and interference every step of the way. Chapter 22 shows the alternative that sidesteps almost all of this: instead of electrons and electromagnetic fields in a conductor, send the signal as light, guided down a strand of glass that's immune to electromagnetic crosstalk by its very nature.
