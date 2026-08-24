# Chapter 18: Shannon's Limit

*"You can build a more sensitive receiver. You can invent a cleverer modulation scheme. You can pay for a better cable. But you cannot argue with mathematics — and in 1948, Claude Shannon proved, once and for all, exactly how much data any channel can carry, and no engineering genius in the 80 years since has broken that ceiling. They've only gotten closer to it."*

---

## Table of Contents

1. [The Problem: Is There a Hard Ceiling?](#1-the-problem-is-there-a-hard-ceiling)
2. [Recap: Nyquist's Noiseless Ceiling](#2-recap-nyquists-noiseless-ceiling)
3. [Shannon's Insight, Intuitively](#3-shannons-insight-intuitively)
4. [The Shannon-Hartley Theorem — The Actual Formula](#4-the-shannon-hartley-theorem--the-actual-formula)
5. [Worked Example 1: Why Analog Dial-Up Modems Capped Near 33.6 kbps](#5-worked-example-1-why-analog-dial-up-modems-capped-near-336-kbps)
6. [Worked Example 2: Why 56K Modems Could (Almost) Beat Shannon](#6-worked-example-2-why-56k-modems-could-almost-beat-shannon)
7. [Worked Example 3: A Modern Wi-Fi Channel](#7-worked-example-3-a-modern-wi-fi-channel)
8. [Worked Example 4: Why Fiber Carries Terabits](#8-worked-example-4-why-fiber-carries-terabits)
9. [The Bandwidth/SNR Tradeoff, Visualized](#9-the-bandwidthsnr-tradeoff-visualized)
10. [Why This Is Physics, Not Engineering Laziness](#10-why-this-is-physics-not-engineering-laziness)
11. [What Shannon's Theorem Does NOT Tell You](#11-what-shannons-theorem-does-not-tell-you)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Hands-On Experiment](#13-hands-on-experiment)
14. [Code: Shannon Capacity Calculator in Go](#14-code-shannon-capacity-calculator-in-go)
15. [What This Chapter Simplified](#15-what-this-chapter-simplified)
16. [Interview Questions & Model Answers](#interview-questions--model-answers)
17. [Exercises](#exercises)
18. [Summary](#summary)
19. [Bridge to Chapter 19](#19-bridge-to-chapter-19)

---

## 1. The Problem: Is There a Hard Ceiling?

Every chapter in this volume has been building toward one question, asked plainly now: **given a channel with a certain bandwidth (Chapter 16) and a certain amount of noise (Chapter 17), is there an absolute, mathematically provable maximum number of bits per second it can carry — a ceiling no engineer, no matter how brilliant, can ever exceed?**

This isn't a rhetorical question with an obvious "no, just keep innovating" answer. In 1948, at Bell Labs — the same institution that had spent decades building the telephone network this entire volume keeps returning to — a mathematician named **Claude Shannon** answered it with a definitive, proven **yes**, in a paper widely regarded as founding the entire field of information theory. This chapter presents that answer, both intuitively and as the literal formula, and then uses it to explain two of the most concrete, checkable facts in networking history: why dial-up modems topped out where they did, and why a strand of glass can carry more data than every phone call in the 20th century combined.

---

## 2. Recap: Nyquist's Noiseless Ceiling

Chapter 16, Section 10 gave Harry Nyquist's 1928 result for a **noiseless** channel: the maximum symbol rate is 2B (twice the bandwidth in Hz), and combined with bits-per-symbol from a chosen modulation order M, the noiseless bit-rate ceiling is:

```
  C_noiseless = 2B · log2(M)
```

The gap Nyquist's formula leaves open is exactly what Chapter 17 spent an entire chapter quantifying: **how large can M actually be, given that the channel has noise?** Chapter 16, Section 5 showed that higher-order modulation (larger M) packs constellation points closer together, making them easier to confuse when noise is present. Shannon's genius was finding the precise mathematical relationship between noise (specifically, SNR) and the maximum *reliable* number of distinguishable levels — closing Nyquist's gap completely.

---

## 3. Shannon's Insight, Intuitively

**Intuitive version:** imagine trying to communicate with a friend across a room by holding up your hand at different heights, and your friend has to read off which of several height "levels" you're showing. If the room is perfectly still, you could use a hundred finely-spaced height levels and your friend could read each one precisely. But if there's a slight, constant tremor in your hand — noise — your friend can no longer reliably tell "level 47" from "level 48," because the tremor is bigger than the gap between adjacent levels. The tremor's size directly determines how many *distinguishable* levels you can safely use. A bigger tremor (more noise) forces you to use fewer, more widely-spaced levels (a coarser code) to stay reliably distinguishable — directly costing you information per "symbol" (each hand position), even though you're still holding your hand up just as often (the same symbol rate).

**The two levers Shannon combined:** Nyquist already answered "how many *symbols* per second can I send" (the hand's raising/lowering rate, capped by bandwidth). Shannon answered the missing half: "given the *noise level*, how many *distinguishable levels* can each symbol safely encode" — and multiplying those two answers together gives the true, complete capacity of the channel.

---

## 4. The Shannon-Hartley Theorem — The Actual Formula

```
  C = B · log2(1 + S/N)

  where:
    C   = channel capacity, in bits per second (bps) -- the absolute maximum
          reliable throughput, given infinite cleverness in coding
    B   = bandwidth, in Hz (Chapter 16's precise definition)
    S/N = the signal-to-noise ratio, as a plain power RATIO (not in dB!)
          -- to convert from dB: S/N (ratio) = 10^(SNR_dB / 10)
```

Every symbol in this formula has already been rigorously defined in this volume: B in Chapter 16, S/N (as SNR) in Chapter 17. Shannon's contribution was proving that this specific combination — bandwidth times the logarithm of (1 + SNR) — is the *exact* upper bound on reliable capacity, and, remarkably, that this bound can be *approached arbitrarily closely* with sufficiently clever coding (though never exceeded). This second half of Shannon's result — that the limit is not just an upper bound but an *achievable* one, in principle — is precisely why practical modulation and coding schemes (Chapters 15-16, and error correction in Chapter 20) have spent 80 years getting closer and closer to it.

**Why "1 + S/N" and not just "S/N"?** Adding 1 handles the edge case cleanly: even if there were zero noise at all (N → 0, so S/N → infinity), the formula still gives a sensible, finite-symbol-rate-based answer for any real, achievable SNR; and critically, if the signal is nonexistent (S = 0), capacity correctly comes out to B·log2(1) = 0 — no signal, no capacity, exactly as physically required.

### Spectral Efficiency — How Engineers Actually Talk About "Closeness to Shannon"

Rearranging Shannon's formula slightly gives a quantity real engineers use constantly when comparing modulation and coding schemes:

```
  Spectral efficiency = C / B = log2(1 + S/N)      (units: bits/second/Hz)
```

This single number — bits per second, per hertz of bandwidth used — is the standard metric for how "efficient" a real system's modulation and coding is, independent of how much raw bandwidth it happens to have available. A real modulation-and-coding scheme's spectral efficiency, measured in a lab or the field, can be directly compared to the Shannon-predicted maximum for the measured SNR, giving engineers a precise, quantitative answer to "how much headroom is left before we hit the theoretical wall?" Wi-Fi and cellular standards documents (5G, Wi-Fi 6/7) routinely publish exactly this figure for each supported modulation-and-coding scheme (MCS) index, alongside the SNR each one requires — a real, production artifact of the theorem this chapter is built around, and the mechanism behind Chapter 16 Section 5's adaptive modulation: a device's real-time SNR estimate selects the MCS index whose spectral efficiency the current channel can support.

---

## 5. Worked Example 1: Why Analog Dial-Up Modems Capped Near 33.6 kbps

A standard analog telephone line has approximately B = 3,000 Hz of usable bandwidth (the ~300-3,400 Hz voice band from Chapter 16, Section 7, rounded for this calculation). A good-quality analog line typically achieves an SNR around 30-35 dB.

**Step 1 — convert SNR from dB to a plain ratio (Section 4's note):**

```
  SNR = 30 dB  -->  S/N (ratio) = 10^(30/10) = 10^3 = 1,000
```

**Step 2 — apply Shannon-Hartley:**

```
  C = B · log2(1 + S/N)
    = 3,000 · log2(1 + 1,000)
    = 3,000 · log2(1,001)
    = 3,000 · 9.967
    ≈ 29,900 bps  ≈  ~30 kbps
```

**This lands almost exactly on the real-world ceiling of the V.34 modem standard (33.6 kbps), the fastest *purely analog* dial-up standard ever deployed** — a standard whose engineers, working directly against Shannon's formula with the real, measured SNR of the telephone network, pushed modulation (multi-level, adaptive QAM-like trellis-coded schemes, direct descendants of Chapter 16's ideas) as close to this hard ceiling as coding theory allowed. This is not a coincidence — it's the entire point of Shannon's theorem: it correctly predicts, decades in advance, where real engineering effort would eventually top out, because it's math, not a guess.

---

## 6. Worked Example 2: Why 56K Modems Could (Almost) Beat Shannon

Here is the genuinely surprising part, and the reason this worked example gets special attention: **56K modems (the V.90/V.92 standards) achieved download speeds up to 56 kbps — nearly double Section 5's ~30 kbps analog Shannon ceiling — and they did it without breaking Shannon's law, because they weren't operating over the same channel Section 5 analyzed.**

By the mid-1990s, the telephone network's *core* (between telephone company central offices) had become almost entirely **digital**, not analog. A voice call is digitized at the first central office it reaches, using **Pulse Code Modulation (PCM)**: it's sampled 8,000 times per second, and each sample is quantized into one of 256 levels (8 bits per sample) — precisely Chapter 14's bits-and-bytes idea, applied to voice:

```
  Digital PCM channel capacity = sample rate x bits per sample
                                = 8,000 samples/sec x 8 bits/sample
                                = 64,000 bps = 64 kbps
```

A 56K modem's ISP-side connection was **directly digital** — the ISP connected to the phone network over a digital trunk line (like a T1), skipping the analog conversion entirely on that end. Only the *user's* side still had one analog-to-digital conversion (from your phone line into the phone company's local switch). Because there was only **one** analog conversion in the whole path (versus **two**, one at each end, for a normal analog modem-to-modem call, which is what Section 5 assumed), the download direction avoided a full round-trip's worth of analog quantization noise — meaning it wasn't limited by Section 5's analog SNR calculation at all. It was limited instead by the **digital PCM channel's own raw ceiling: 64 kbps.**

That 64 kbps ceiling was then reduced by two further, very concrete real-world factors:

```
  64 kbps (raw PCM channel capacity)
  -  8 kbps  (lost to "robbed-bit signaling": in the US T1 standard, every
              6th frame steals its least-significant bit for signaling
              overhead, effectively reducing usable payload)
  = 56 kbps  (the theoretical V.90 ceiling, and where the standard's name
              comes from)

  In practice, further reduced to ~53 kbps by FCC Part 68 regulations
  limiting the maximum power telephone equipment may put on the shared
  phone network (to prevent crosstalk -- Chapter 17, Section 7 -- with
  other lines), which is why real-world 56K connections almost never
  actually reached the full 56.0 kbps figure printed on the box.
```

The *upload* direction of a 56K connection still had to pass through a real analog-to-digital conversion (your modem's analog signal, converted to digital at the phone company), so it remained subject to Section 5's ordinary Shannon-limited analog calculation — which is exactly why 56K modems were always asymmetric: fast downloads (~56 kbps, digital-PCM-limited), slower uploads (~33.6 kbps, analog-Shannon-limited). This asymmetry is a direct, checkable fingerprint of the exact mechanism described above, and it's a perfect illustration of Section 10's core lesson: Shannon's limit is never actually violated — you can only ever find a genuinely different, better channel to be limited by.

---

## 7. Worked Example 3: A Modern Wi-Fi Channel

Apply the same formula to Chapter 17, Section 11's worked SNR result: a 20 MHz Wi-Fi channel with an (optimistic, line-of-sight) SNR of roughly 79 dB.

```
  B = 20,000,000 Hz
  SNR = 79 dB --> S/N (ratio) = 10^(79/10) = 10^7.9 ≈ 79,432,823

  C = B · log2(1 + S/N)
    = 20,000,000 · log2(1 + 79,432,823)
    = 20,000,000 · log2(79,432,824)
    ≈ 20,000,000 · 26.24
    ≈ 524,800,000 bps  ≈  ~525 Mbps
```

Compare this to Chapter 16, Section 11's modulation-based estimate of ~1.28 Gbps *ceiling from 256-QAM alone* — Shannon's number here (525 Mbps) is the true information-theoretic ceiling given this specific (optimistic, idealized) SNR, while Chapter 16's number was the ceiling *if* you could use 256-QAM without any error at all. The two numbers being in the same broad neighborhood, despite being derived from completely different reasoning, is exactly the kind of consistency check that gives confidence both models are capturing real physics correctly — real, deployed 802.11ac hardware, with realistic (much lower than idealized free-space) indoor SNR values, achieves headline rates in the hundreds-of-Mbps range, consistent with both calculations once realistic SNR is used instead of the idealized figure from Chapter 17.

---

## 8. Worked Example 4: Why Fiber Carries Terabits

Chapter 16, Section 2 noted fiber operates around 190-196 THz — but the relevant number for Shannon's formula isn't the carrier frequency itself, it's the **usable bandwidth** around it, plus fiber's extraordinarily high achievable SNR (thanks to fiber's near-total immunity to electromagnetic interference, Chapter 17 Section 6, and its very low attenuation, Chapter 17 Section 4).

A single modern long-haul fiber system might use roughly 4-8 THz of usable bandwidth (spread across many individual wavelength channels via **Wavelength-Division Multiplexing**, previewed here and covered fully in Chapter 22) with an SNR that can exceed 40 dB thanks to fiber's clean transmission characteristics:

```
  B = 4 x 10^12 Hz (4 THz, illustrative combined bandwidth across many
                     WDM wavelength channels on one fiber)
  SNR = 40 dB --> S/N (ratio) = 10^(40/10) = 10,000

  C = B · log2(1 + S/N)
    = 4 x 10^12 · log2(10,001)
    ≈ 4 x 10^12 · 13.29
    ≈ 5.3 x 10^13 bps
    ≈ 53 Terabits per second
```

This is not a hypothetical exercise — real commercial long-haul fiber systems, combining WDM with dense modulation and advanced forward error correction (previewed in Chapter 20), have demonstrated aggregate capacities well into the tens of terabits per second on a single fiber strand. The reason isn't that fiber "cheats" Shannon's law — it's that fiber offers a genuinely enormous B (many terahertz of usable, low-noise spectrum available across the optical windows used for transmission) compared to any radio or copper system, and Shannon's formula scales capacity *linearly* with bandwidth — so a channel with a thousand times the bandwidth of Wi-Fi, and comparable or better SNR, can support proportionally vastly more throughput.

---

## 9. The Bandwidth/SNR Tradeoff, Visualized

Shannon's formula reveals that bandwidth and SNR are not independent contributors to capacity — they trade off against each other, but not symmetrically:

```
  C = B · log2(1 + S/N)

  Doubling B:          capacity roughly DOUBLES (linear relationship)
  Doubling S/N (ratio): capacity increases by log2(2x) worth of the
                        log term -- a much SMALLER effect once S/N
                        is already large, because log2 grows very slowly

  Concretely, from Section 7's Wi-Fi example (S/N ≈ 79.4 million):
    Going from S/N = 79.4M to S/N = 158.8M (doubling SNR, an enormous,
    often physically difficult power increase) only grows log2(1+S/N)
    from about 26.24 to about 27.24 -- a mere ~4% increase in capacity

    Going from B = 20 MHz to B = 40 MHz (doubling bandwidth, often as
    simple as combining two adjacent Wi-Fi channels) directly DOUBLES
    capacity, all else equal
```

**This single asymmetry is why real engineering effort overwhelmingly focuses on acquiring more bandwidth (wider channels, more spectrum, WDM in fiber) rather than chasing marginal SNR improvements once SNR is already reasonably good** — a transmitter would need to increase its power by a factor of 1,000 to gain the same capacity that simply doubling its bandwidth provides for free. This explains, among many other things, why cellular carriers spend billions of dollars bidding for additional spectrum (bandwidth) at government auctions rather than simply building more powerful towers, and why Wi-Fi 6E and Wi-Fi 7's headline feature is opening up entirely new frequency bands (6 GHz) rather than boosting transmit power.

### Sidebar: Why DSL Massively Outperformed Dial-Up, Using the Exact Same Copper Wire

This chapter's asymmetry lesson has one more, very concrete illustration, using a technology built on the very same physical copper telephone wire that dial-up modems used: **DSL (Digital Subscriber Line)**.

Dial-up modems were restricted to the ~3,000 Hz voice band (Section 5) because the phone company's *switching equipment* — not the copper wire itself — filtered everything outside that band, since only the voice band needed to pass through the voice-switched telephone network to reach another phone. DSL equipment, installed directly at the phone company's local switch and at the customer's premises, taps into the copper wire *before* it hits that voice-band-filtering equipment, and uses frequencies **above** the voice band — typically up to about 1.1 MHz for ADSL2+ — that a plain telephone was never designed to use, but that the same physical copper wire can still carry over short distances.

```
  Recomputing Shannon capacity for DSL vs. dial-up, same copper pair:

  Dial-up:  B = 3,000 Hz,     SNR ≈ 30 dB (S/N ratio = 1,000)
            C = 3,000 x log2(1,001) ≈ 29,900 bps

  ADSL2+:   B ≈ 1,100,000 Hz (up to 1.1 MHz, using frequencies a phone
            never needed), SNR varies a lot with distance from the
            exchange, but even a conservative 20 dB average across
            that much wider band gives:
            S/N ratio = 10^(20/10) = 100
            C = 1,100,000 x log2(101) ≈ 1,100,000 x 6.66 ≈ 7.3 Mbps
```

Even with a considerably *worse* average SNR assumption than dial-up's, DSL's dramatically larger bandwidth (1.1 MHz vs. 3 kHz — roughly 367x more) still produces a capacity roughly 240x higher. This is Section 9's bandwidth/SNR asymmetry, playing out as a real, historical technology transition: DSL didn't win by "trying harder" with cleverer modulation on the same narrow band — it won by physically accessing far more bandwidth on the exact same wire dial-up had always used, simply by not routing the signal through the equipment that had been filtering it down to voice-only frequencies for a century.

---

## 10. Why This Is Physics, Not Engineering Laziness

It's worth stating directly, because it's easy to assume any "limit" in technology is temporary and will eventually be engineered around: **Shannon's limit is a mathematically proven upper bound, not a current-technology bottleneck.** No future breakthrough in silicon fabrication, no cleverer algorithm, no quantum leap in materials science can ever push reliable throughput past B·log2(1+S/N) for a given, fixed bandwidth and SNR — because the proof isn't about the limits of 1948 (or 2026) engineering, it's a limit derived from information theory and probability itself, as fundamental within its domain as the speed of light is to physics.

What *does* keep improving, and will keep improving indefinitely, are the two inputs to the formula: engineers find ways to access **more bandwidth** (new frequency bands, WDM in fiber, wider Wi-Fi channels) and **better SNR** (better antennas, lower-noise receiver electronics, less interference through smarter channel management) — and coding schemes get closer to the theoretical ceiling that a given B and SNR already permit (the subject of Chapter 20's forward error correction). Every "internet just got faster" headline in the last 80 years has come from moving one or both of Shannon's two inputs, never from breaking his formula.

---

## 11. What Shannon's Theorem Does NOT Tell You

Shannon's theorem proves a capacity number *exists* and is *achievable in principle* — it does not say:

- **How** to achieve it. Shannon's original proof was a non-constructive existence proof — it showed that codes achieving near-capacity reliable communication must exist, without showing how to build one. It took engineers until the 1990s-2000s (with turbo codes and LDPC codes) to find *practical* codes that get close to the Shannon limit for real channels — over 50 years after the theorem was proven.
- **What happens when noise does corrupt data anyway.** Shannon's formula describes the maximum *reliable* rate — but "reliable" here means "with an error rate that can be made arbitrarily small via sufficiently long, clever coding," not "error will literally never happen." Real systems, running below Shannon's limit with real, finite-length codes, still experience individual bit errors from real noise events (Chapter 17). Detecting when that happens — and, ideally, fixing it — is a completely separate engineering problem that Shannon's theorem is silent on.
- **Latency, cost, or power consumption** of actually building a system that approaches the limit — a scheme that theoretically approaches Shannon's capacity might require impractically long codewords, enormous processing delay, or power consumption unsuitable for a battery-powered device. Real systems always make a practical tradeoff, operating somewhat below the theoretical ceiling in exchange for acceptable latency and cost.

That first gap — Shannon says a certain throughput is theoretically achievable with *vanishingly small* error, but doesn't guarantee *zero* error at any real, finite code length — is exactly why every real digital communication system, no matter how good its modulation or how close to Shannon's limit its coding gets, still needs a mechanism to catch and handle the errors that inevitably slip through. That mechanism is the entire subject of the next chapter.

---

## 12. Common Misconceptions

- **"56K modems proved Shannon's limit could be beaten."** No — Section 6 showed precisely why not: the 56K downstream direction was operating over a fundamentally different, digital PCM channel with its own separate 64 kbps ceiling, not the same analog voice-band channel Section 5 analyzed. Shannon's limit, applied correctly to whichever actual channel is in use, was never violated.
- **"More bandwidth always beats more SNR, in every situation."** Section 9 showed bandwidth's advantage assumes the *ratio* increase is comparable — in absolute terms, if SNR is currently very poor (near 0 dB), even modest SNR improvements can matter enormously, since log2(1+S/N) changes rapidly for small S/N. The "bandwidth wins" conclusion specifically applies once SNR is already reasonably good, which is the common case for most modern, well-engineered links.
- **"Shannon's formula tells you the actual speed you'll get."** It tells you the theoretical maximum, assuming ideal coding. Real, practical throughput is always somewhat below Shannon's number, due to protocol overhead, imperfect (though very good) real-world codes, and the practical latency/complexity tradeoffs from Section 11.
- **"Shannon's theorem is only relevant to old technology like dial-up."** It is exactly as relevant, and used just as directly by engineers, in the design of 5G, Wi-Fi 7, and the latest fiber optic systems — Sections 7 and 8 are not historical curiosities, they're live, current engineering practice.

---

## 13. Hands-On Experiment

Compute your own home connection's theoretical Shannon ceiling and compare it to what your ISP actually advertises.

```bash
# Estimate your Wi-Fi channel's bandwidth and current SNR (from Chapter 17's
# hands-on experiment):
/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I
# (macOS -- look at agrCtlRSSI and agrCtlNoise, subtract for SNR in dB)
# or: iw dev wlan0 station dump   (Linux -- look for "signal:")

# Then, using this chapter's Section 14 Go program (or by hand using
# Section 4's formula), plug in your measured SNR and your channel's
# known bandwidth (20/40/80/160 MHz, visible in your router's settings)
# to compute your theoretical Shannon ceiling.

# Compare that number to the actual negotiated PHY rate your OS reports
# (Chapter 16 Section 13's experiment) -- real PHY rate should sit
# somewhat BELOW your calculated Shannon ceiling, since real modulation
# and coding schemes approach, but never exceed, the theoretical limit.
```

---

## 14. Code: Shannon Capacity Calculator in Go

```go
package main

import (
	"fmt"
	"math"
)

// shannonCapacity implements the Shannon-Hartley theorem: C = B log2(1 + S/N)
// snrDB is the signal-to-noise ratio in decibels (converted internally to a ratio).
func shannonCapacity(bandwidthHz, snrDB float64) float64 {
	snrRatio := math.Pow(10, snrDB/10)
	return bandwidthHz * math.Log2(1+snrRatio)
}

func main() {
	fmt.Println("--- Section 5: Analog dial-up over a phone line ---")
	c := shannonCapacity(3000, 30)
	fmt.Printf("  B=3000 Hz, SNR=30 dB  ->  C = %.0f bps (%.2f kbps)\n", c, c/1000)

	fmt.Println("\n--- Section 6: Digital PCM ceiling (not a Shannon formula --")
	fmt.Println("    included for direct comparison to the analog figure above) ---")
	pcm := 8000.0 * 8 // sample rate x bits/sample
	fmt.Printf("  Raw PCM channel: %.0f bps (%.0f kbps)\n", pcm, pcm/1000)
	fmt.Printf("  Minus robbed-bit signaling (8 kbps): %.0f kbps (the V.90 ceiling)\n", pcm/1000-8)

	fmt.Println("\n--- Section 7: A modern 20 MHz Wi-Fi channel ---")
	c = shannonCapacity(20_000_000, 79)
	fmt.Printf("  B=20 MHz, SNR=79 dB  ->  C = %.0f bps (%.2f Mbps)\n", c, c/1_000_000)

	fmt.Println("\n--- Section 8: A long-haul fiber system (illustrative) ---")
	c = shannonCapacity(4e12, 40)
	fmt.Printf("  B=4 THz, SNR=40 dB  ->  C = %.3e bps (%.2f Tbps)\n", c, c/1e12)

	fmt.Println("\n--- Section 9: Bandwidth vs SNR asymmetry ---")
	base := shannonCapacity(20_000_000, 79)
	doubleBW := shannonCapacity(40_000_000, 79)
	doubleSNR := shannonCapacity(20_000_000, 82) // +3dB ~ doubling the ratio
	fmt.Printf("  Base:        %.2f Mbps\n", base/1e6)
	fmt.Printf("  Double BW:   %.2f Mbps (+%.0f%%)\n", doubleBW/1e6, (doubleBW/base-1)*100)
	fmt.Printf("  Double SNR:  %.2f Mbps (+%.1f%%)\n", doubleSNR/1e6, (doubleSNR/base-1)*100)

	fmt.Println("\n--- Sidebar: Spectral efficiency (bits/s/Hz) for a few SNR values ---")
	for _, snrDB := range []float64{0, 10, 20, 30, 40} {
		snrRatio := math.Pow(10, snrDB/10)
		eff := math.Log2(1 + snrRatio)
		fmt.Printf("  SNR=%4.0f dB  ->  spectral efficiency = %.2f bits/s/Hz\n", snrDB, eff)
	}

	fmt.Println("\n--- Sidebar: DSL vs dial-up, same copper wire ---")
	dialup := shannonCapacity(3000, 30)
	dsl := shannonCapacity(1_100_000, 20)
	fmt.Printf("  Dial-up (3 kHz,  30 dB): %.0f bps\n", dialup)
	fmt.Printf("  DSL     (1.1MHz, 20 dB): %.0f bps (%.0fx dial-up)\n", dsl, dsl/dialup)
}
```

```
Output:

--- Section 5: Analog dial-up over a phone line ---
  B=3000 Hz, SNR=30 dB  ->  C = 29901 bps (29.90 kbps)

--- Section 6: Digital PCM ceiling (not a Shannon formula --
    included for direct comparison to the analog figure above) ---
  Raw PCM channel: 64000 bps (64 kbps)
  Minus robbed-bit signaling (8 kbps): 56 kbps (the V.90 ceiling)

--- Section 7: A modern 20 MHz Wi-Fi channel ---
  B=20 MHz, SNR=79 dB  ->  C = 524717797 bps (524.72 Mbps)

--- Section 8: A long-haul fiber system (illustrative) ---
  B=4 THz, SNR=40 dB  ->  C = 5.310e+13 bps (53.10 Tbps)

--- Section 9: Bandwidth vs SNR asymmetry ---
  Base:        524.72 Mbps
  Double BW:   1049.44 Mbps (+100%)
  Double SNR:  545.18 Mbps (+3.9%)

--- Sidebar: Spectral efficiency (bits/s/Hz) for a few SNR values ---
  SNR=   0 dB  ->  spectral efficiency = 1.00 bits/s/Hz
  SNR=  10 dB  ->  spectral efficiency = 3.46 bits/s/Hz
  SNR=  20 dB  ->  spectral efficiency = 6.66 bits/s/Hz
  SNR=  30 dB  ->  spectral efficiency = 9.97 bits/s/Hz
  SNR=  40 dB  ->  spectral efficiency = 13.29 bits/s/Hz

--- Sidebar: DSL vs dial-up, same copper wire ---
  Dial-up (3 kHz,  30 dB): 29901 bps
  DSL     (1.1MHz, 20 dB): 7333066 bps (245x dial-up)
```

This single program reproduces every major numeric claim in this chapter — run it, change the inputs, and confirm the pattern Section 9 describes: bandwidth changes move capacity proportionally; SNR changes (once SNR is already large) barely move it at all.

---

## 15. What This Chapter Simplified

- Real telephone-line SNR varies significantly by line quality, distance from the central office, and country — the 30 dB figure in Section 5 is a reasonable, commonly-cited illustrative value, not a universal constant.
- Section 8's fiber figures are illustrative of the *scale* real WDM long-haul systems achieve, not a specific deployed system's exact spec — real systems vary widely in wavelength count, per-channel bandwidth, and modulation choice (detailed in Chapter 22).
- Shannon's theorem, as stated, assumes an "Additive White Gaussian Noise" (AWGN) channel model — a specific, idealized mathematical model of noise that real channels approximate well but never match perfectly; real channels also have interference (Chapter 17, Section 6) that isn't part of the pure AWGN model, though extended versions of the theory handle this.
- This chapter didn't derive the theorem mathematically (a task requiring probability theory and information entropy well beyond this course's scope) — it presented the intuition and the result, which is sufficient for every practical, engineering-facing purpose this course cares about.

---

## Interview Questions & Model Answers

**Beginner: "What does Shannon's theorem say, in plain language?"**

It says that for any communication channel with a given bandwidth and a given amount of noise, there is a hard mathematical maximum on how many bits per second can be sent reliably — and that this maximum can be approached, with sufficiently clever coding, but never exceeded, no matter how advanced the technology. The formula is C = B·log2(1+S/N), where B is bandwidth in Hz and S/N is the signal-to-noise power ratio.

**Intermediate: "Explain, using Shannon's theorem, why dial-up modems capped out where they did."**

A standard analog phone line has about 3,000 Hz of bandwidth and a typical SNR around 30 dB. Plugging those into C = B·log2(1+S/N) gives a theoretical ceiling of about 30 kbps, which closely matches the fastest purely analog modem standard, V.34 at 33.6 kbps. The later 56K standard appeared to beat this, but it wasn't actually beating Shannon's law for the same channel — it exploited the fact that the phone network's core had become digital, so the download direction only had one analog-to-digital conversion instead of two, putting it under a completely different ceiling: the digital PCM channel's own 64 kbps capacity (8,000 samples/second times 8 bits/sample), reduced to 56 kbps by signaling overhead and further to about 53 kbps by FCC power regulations.

**Advanced: "Given the Shannon-Hartley formula, explain why network engineers today invest more in acquiring additional spectrum/bandwidth than in improving signal-to-noise ratio, once SNR is already reasonably good."**

Because capacity scales linearly with bandwidth but only logarithmically with SNR, once SNR is already at a reasonably high value, further SNR improvements yield rapidly diminishing returns — doubling the signal-to-noise ratio (a substantial, often costly engineering effort, like increasing transmit power tenfold) might only increase capacity by a few percent, because you're adding a fixed amount to log2(1+S/N) which is already a large number. Doubling bandwidth, by contrast, directly and proportionally doubles capacity, with no diminishing-returns effect. This is precisely why cellular carriers spend enormous sums bidding for additional spectrum at auction, why Wi-Fi's biggest recent upgrades (6E, 7) center on opening new frequency bands rather than boosting radio power, and why long-haul fiber systems use wavelength-division multiplexing to add more usable bandwidth on a single strand rather than trying to further reduce an already-excellent noise floor.

---

## Exercises

### Easy

1. State the Shannon-Hartley formula from memory and explain what each symbol means.
2. Using the formula, calculate the capacity of a channel with B = 4,000 Hz and SNR = 20 dB.
3. In one or two sentences, explain why 56K modems' download and upload speeds were different (asymmetric).

### Medium

4. A channel's SNR improves from 20 dB to 40 dB (a 100x improvement in the linear ratio). Using B = 1 MHz, calculate the capacity at both SNR values and calculate the percentage increase in capacity. Compare this to the percentage increase you'd get from simply doubling the bandwidth instead, at the original 20 dB SNR.
5. Explain, using Section 6's numbers, why the "K" in "56K modem" was somewhat misleading in most real-world usage.
6. Using Section 14's Go program as a starting point, compute the Shannon capacity for a hypothetical 5G mmWave channel with B = 400 MHz and SNR = 25 dB. Compare it to the Section 7 Wi-Fi example.

### Hard

7. Shannon's proof is described in Section 11 as "non-constructive" — it proves capacity-achieving codes exist without showing how to build one. Research briefly (or reason from what you know of error-correcting codes) why this gap between "proven possible" and "practically achievable" might exist for a purely mathematical existence proof, and why it took engineers decades to close it with real codes (turbo codes, LDPC).
8. A satellite internet link has extremely limited bandwidth (a few MHz) but, due to its high transmit power and clear line of sight to space, an excellent SNR (45+ dB). A terrestrial fiber link has enormous bandwidth (several THz) but must share it across thousands of users. Using Shannon's formula and Section 9's asymmetry argument, explain why the satellite link, despite its excellent SNR, will likely never match the fiber link's aggregate capacity — and what the satellite operator's only real lever is for increasing its own capacity.

---

## Summary

| Term | Meaning |
|---|---|
| Nyquist rate | Noiseless maximum symbol rate = 2B (Chapter 16) |
| Shannon-Hartley theorem | C = B·log2(1+S/N); the absolute maximum reliable channel capacity |
| Channel capacity (C) | Maximum bits/second achievable with arbitrarily low error, given enough coding cleverness |
| S/N (ratio) | Signal-to-noise ratio as a plain power ratio, not in dB, for use in Shannon's formula |
| Bandwidth/SNR tradeoff | Capacity scales linearly with bandwidth, only logarithmically with SNR |
| Non-constructive proof | Shannon proved capacity-achieving codes exist without showing how to build them |
| PCM (Pulse Code Modulation) | Digitizing voice: 8,000 samples/sec x 8 bits/sample = 64 kbps digital channel |

---

## 19. Bridge to Chapter 19

Shannon's theorem draws a hard line around how much data a channel *can* carry — but it says nothing about what happens to the data that gets through anyway with an error, when real, finite-length codes fall short of the theoretical ideal, or when a sudden burst of noise (a lightning strike near a cable, a car door slamming near a Wi-Fi router, a scratch on a fiber connector) overwhelms even a well-engineered link for a fraction of a second. Every mechanism built across this entire volume — attenuation, interference, noise, SNR, and now Shannon's limit — describes *how likely* a bit is to be corrupted in transit. None of them describe what a receiver actually *does* when a bit arrives flipped: a 1 that was sent as a 0, or a 0 that arrived as a 1.

That is where Chapter 19, **Error Detection**, begins: given that noise sometimes corrupts a bit despite every precaution in this volume, how does a receiver even know it happened — and what is the simplest possible way to check?

