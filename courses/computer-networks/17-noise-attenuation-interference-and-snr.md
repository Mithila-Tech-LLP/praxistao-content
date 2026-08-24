# Chapter 17: Noise, Attenuation, Interference, and SNR

*"Every signal that has ever traveled anywhere has arrived a little bit worse than it left. The entire discipline of communications engineering is the science of how much worse — and how to survive it."*

---

## Table of Contents

1. [The Problem: Real Channels Are Never Perfect](#1-the-problem-real-channels-are-never-perfect)
2. [Attenuation — Signals Get Weaker With Distance](#2-attenuation--signals-get-weaker-with-distance)
3. [Attenuation in Copper — Resistance and the Skin Effect](#3-attenuation-in-copper--resistance-and-the-skin-effect)
4. [Attenuation in Fiber — Absorption and Scattering](#4-attenuation-in-fiber--absorption-and-scattering)
5. [Attenuation in Free Space — The Inverse-Square Law](#5-attenuation-in-free-space--the-inverse-square-law)
6. [Interference — Unwanted Energy From Other Sources](#6-interference--unwanted-energy-from-other-sources)
7. [Crosstalk — Interference From the Next Wire Over](#7-crosstalk--interference-from-the-next-wire-over)
8. [Radio Interference — Why Microwave Ovens Disrupt Wi-Fi](#8-radio-interference--why-microwave-ovens-disrupt-wi-fi)
9. [Noise — The Unavoidable Background Hiss](#9-noise--the-unavoidable-background-hiss)
10. [Signal-to-Noise Ratio (SNR), Defined Precisely](#10-signal-to-noise-ratio-snr-defined-precisely)
11. [Worked Example: Computing Real SNR](#11-worked-example-computing-real-snr)
12. [Real Fixes: Repeaters, Shielding, Channel Selection, Coding](#12-real-fixes-repeaters-shielding-channel-selection-coding)
13. [Why Long Ethernet Runs Need Repeaters — the Full Explanation](#13-why-long-ethernet-runs-need-repeaters--the-full-explanation)
14. [Common Misconceptions](#14-common-misconceptions)
15. [Hands-On Experiment](#15-hands-on-experiment)
16. [Code: dB, Attenuation, and SNR in Go](#16-code-db-attenuation-and-snr-in-go)
17. [What This Chapter Simplified](#17-what-this-chapter-simplified)
18. [Interview Questions & Model Answers](#interview-questions--model-answers)
19. [Exercises](#exercises)
20. [Summary](#summary)

---

## 1. The Problem: Real Channels Are Never Perfect

Chapter 16 ended by admitting a quiet assumption: the Nyquist symbol-rate formula (2B) assumed a **noiseless** channel — one with no distortion, no unwanted energy, nothing standing between the transmitted signal and a perfect receiver. No such channel exists anywhere in the physical universe.

Every real signal, from the moment it leaves a transmitter, is under continuous physical attack from three distinct enemies, which this chapter treats as three separate, precisely defined phenomena because engineers fix each one differently:

1. **Attenuation** — the signal itself gets weaker simply from traveling, independent of anything else happening on the channel.
2. **Interference** — unwanted energy from *other, identifiable* signal sources leaks into the channel (another Wi-Fi network, a microwave oven, a neighboring wire).
3. **Noise** — unwanted, essentially random energy that exists in every channel, from fundamental physics, with no single identifiable "culprit" to blame or filter out.

Understanding these three separately, and then combining their combined effect into a single number — **Signal-to-Noise Ratio (SNR)** — is what finally lets Chapter 18 complete Shannon's formula and answer, with a real number, exactly how much data a real, imperfect channel can carry.

---

## 2. Attenuation — Signals Get Weaker With Distance

**Intuitive version:** shout across an empty room, and a friend hears you clearly. Shout the same way across a football field, and they hear something faint and hard to make out — the sound wave's energy has spread out and been absorbed by air along the way. Nothing "attacked" your shout; it simply weakened by traveling.

**Engineering definition: attenuation** is the reduction in a signal's amplitude (power) as it travels through a medium, caused by the medium itself absorbing, scattering, or spreading the signal's energy. It is expressed, using Chapter 16's dB scale, as a **loss** — a negative number of decibels, or a "dB per unit distance" figure that lets you calculate total loss over any length:

```
  Total attenuation (dB) = attenuation-per-unit-length x length

  Example: a cable rated at 0.5 dB/meter loss, run for 100 meters:
     Total loss = 0.5 x 100 = 50 dB of signal loss
```

Every transmission medium attenuates signals, but for entirely different physical reasons — worth examining separately, because the fixes (Section 12, and Chapters 21-23) differ completely between them.

---

## 3. Attenuation in Copper — Resistance and the Skin Effect

Copper wire has real electrical resistance, and pushing current through resistance always dissipates some energy as heat (this is simple Ohm's-law physics: power dissipated = I²R). Over distance, this steadily bleeds energy out of the signal.

There's a second, frequency-dependent effect that matters a great deal for the higher frequencies used by modern Ethernet: the **skin effect**. At higher frequencies, current in a conductor increasingly flows only near the conductor's outer surface ("skin") rather than through its full cross-section, effectively shrinking the wire's usable conducting area and *increasing* its effective resistance as frequency rises.

```
  Copper cable attenuation increases with BOTH:
    1. Length      (more resistance encountered = more loss)
    2. Frequency   (skin effect increases effective resistance at higher freq)

  This is precisely why Cat6/Cat6a cabling specifications (fully detailed
  in Chapter 21) publish attenuation figures as a function of BOTH cable
  length AND signaling frequency, e.g.:

     Cat5e @ 100 MHz, 100m:  ~22 dB attenuation (typical spec limit)
     Cat5e @ 100 MHz, 1m:    ~0.22 dB attenuation (roughly linear with length)
```

This dual dependency — worse at both greater length *and* higher frequency — is the direct physical reason Ethernet standards specify a hard maximum cable length (100 meters for common twisted-pair Ethernet) for a given signaling speed: run the cable longer, or push a higher-frequency (faster) signal down it, and attenuation eventually degrades the signal past the point where a receiver can reliably distinguish its logic levels (Chapter 14, Section 3's noise margin, being eaten away).

---

## 4. Attenuation in Fiber — Absorption and Scattering

Light traveling through glass fiber attenuates for different physical reasons than electricity in copper:

- **Absorption:** the glass itself isn't perfectly transparent — some fraction of light energy is absorbed by impurities in the glass and converted to heat, at every point along the fiber's length.
- **Scattering:** microscopic density variations in the glass (an unavoidable property of the material, called Rayleigh scattering) redirect a small fraction of light out of the fiber's intended path at every point along its length.

Both effects are remarkably small per unit length compared to copper's resistive loss — modern single-mode fiber (fully detailed in Chapter 22) attenuates light at roughly:

```
  Single-mode fiber attenuation: ~0.2 dB/km (at the standard 1550nm wavelength)

  Compare a 40 km fiber run:
     Total loss = 0.2 x 40 = 8 dB total

  A comparable-distance copper run would be utterly unusable --
  copper Ethernet is limited to about 100 METERS, not 40 KILOMETERS,
  precisely because of this enormous difference in attenuation-per-distance.
```

This roughly 500x-better attenuation-per-kilometer figure (0.2 dB/km for fiber vs. copper's much steeper loss) is the single biggest reason long-distance and undersea links (Chapter 23) are built exclusively from fiber, not copper.

---

## 5. Attenuation in Free Space — The Inverse-Square Law

A wireless signal radiating from an antenna doesn't travel through a solid, bounded medium at all — it spreads out in all directions through open space, and its power density falls off geometrically as it spreads over an ever-larger sphere. This is captured by the **Free-Space Path Loss (FSPL)** formula:

```
  FSPL (dB) = 20 log10(d) + 20 log10(f) + 20 log10(4π/c)

  A commonly used practical version, with d in kilometers and f in MHz:

  FSPL (dB) = 20 log10(d_km) + 20 log10(f_MHz) + 32.44
```

**Worked example — Wi-Fi at 2.4 GHz over 10 meters:**

```
  d = 0.01 km,  f = 2,400 MHz

  FSPL = 20 log10(0.01) + 20 log10(2400) + 32.44
       = 20 x (-2)      + 20 x 3.38      + 32.44
       = -40            + 67.6           + 32.44
       = 60.04 dB of path loss, just from the geometry of distance --
         before accounting for walls, furniture, or anything else
```

**Worked example — the same Wi-Fi frequency, but 100 meters instead of 10:**

```
  d = 0.1 km (10x farther)

  FSPL = 20 log10(0.1) + 20 log10(2400) + 32.44
       = 20 x (-1)     + 67.6           + 32.44
       = -20 + 67.6 + 32.44 = 80.04 dB
```

Notice: **10x the distance costs an extra 20 dB of loss** — this is the "inverse-square law" made concrete via the logarithm's math (doubling distance always costs the same fixed extra dB, regardless of starting distance, because 20·log10 of a ratio is constant for a constant ratio). This is exactly why moving twice as far from a Wi-Fi router doesn't just "somewhat" weaken the signal — it costs a fixed, substantial, and unavoidable dB penalty no engineering can eliminate, only compensate for (with more transmit power, better antennas, or — as Chapter 16 showed — falling back to a more robust, lower-order modulation).

**Notice also the second term:** higher frequency *also* costs more path loss for the same distance (20·log10(f) grows with f) — this is a real, physical reason (not the only reason, but a real contributing one) that 5 GHz Wi-Fi has shorter effective range than 2.4 GHz Wi-Fi at the same transmit power, a fact Chapter 16's exercises flagged and this section now fully explains.

---

## 6. Interference — Unwanted Energy From Other Sources

**Definition: interference** is unwanted energy in a channel that comes from an *identifiable, separate signal source* — another transmitter, another device, another system — as opposed to noise (Section 9), which has no single identifiable source and exists in a channel by fundamental physics regardless of any other equipment nearby.

The distinction matters practically: interference can sometimes be avoided (change a channel, move away from the source, add shielding), while noise is a fundamental floor you can reduce but never fully eliminate.

---

## 7. Crosstalk — Interference From the Next Wire Over

When two wires run close together and carry changing electrical signals, each wire's changing electromagnetic field induces a small, unwanted signal in its neighbor — this is **crosstalk**, and it's the single biggest reason twisted-pair cabling exists at all (fully developed in Chapter 21).

```
  Two parallel, untwisted wires, each carrying a signal:

  Wire A:  --~~~~~~--       (intended signal)
                 |
                 | (electromagnetic coupling -- unwanted)
                 v
  Wire B:  --~~~~~~--  + small unwanted copy of Wire A's signal
           (intended signal, now corrupted by crosstalk from A)
```

Two specific, named forms matter in real cabling specs:
- **NEXT (Near-End Crosstalk):** interference measured at the same end of the cable where the interfering signal originates — the worst case, since the interfering signal hasn't yet attenuated with distance.
- **FEXT (Far-End Crosstalk):** interference measured at the opposite end from where the interfering signal originates — generally less severe, since the interfering signal has attenuated over the cable's length by the time it reaches the far end.

Chapter 21 shows exactly how twisting each pair of wires around each other (rather than running them in straight parallel lines) makes each half-twist's induced interference cancel against the next half-twist's, dramatically reducing crosstalk — the entire reason "twisted pair" cabling is twisted at all, not just a manufacturing quirk.

---

## 8. Radio Interference — Why Microwave Ovens Disrupt Wi-Fi

This is one of the most concrete, checkable real-world examples in all of networking, and it follows directly and precisely from Chapter 16's frequency concepts.

**The 2.4 GHz Wi-Fi band (802.11b/g/n) occupies roughly 2.400-2.4835 GHz.** A microwave oven works by generating microwave radiation at approximately **2.45 GHz** — chosen specifically because that frequency is well-absorbed by water molecules (which is how it heats food) and happens to sit almost exactly in the middle of the very same unlicensed frequency band Wi-Fi was allocated to use.

```
  2.4 GHz Wi-Fi channels (US, non-overlapping: 1, 6, 11):

  2.400          2.412         2.437         2.462        2.4835 GHz
    |----ch1------|----ch6------|----ch11------|
                        ^
                        |
              Microwave oven leakage
              centered ~2.45 GHz, right in
              the middle of the band
```

Microwave ovens are shielded, but shielding is never perfect — a working oven leaks a small amount of 2.45 GHz radiation, especially through door seals as they age. That leaked energy is:
1. **High power** relative to a Wi-Fi signal — ovens use around 700-1100 watts internally, compared to Wi-Fi transmit power typically well under 1 watt.
2. **Directly overlapping** Wi-Fi's frequency band, unlike, say, a car engine's electrical noise, which occupies entirely different frequencies.

This combination — high power, exactly overlapping frequency — is why running a microwave can visibly degrade or drop a 2.4 GHz Wi-Fi connection for anyone standing near it, especially on channels near the middle of the band (channel 6 suffers more than channel 1 or 11, precisely because of the frequency alignment shown above). This has no equivalent effect on 5 GHz or 6 GHz Wi-Fi, because those bands sit at entirely different frequencies that ovens don't leak into — which is a real, practical reason modern routers push devices toward 5 GHz/6 GHz when possible, beyond just the extra bandwidth those bands offer.

---

## 9. Noise — The Unavoidable Background Hiss

**Definition: noise** is unwanted, essentially random signal energy present in a channel that has no single identifiable transmitting source — it arises from fundamental physical processes rather than from another device you could locate and remove.

The most fundamental and universal type is **thermal noise** (also called **Johnson-Nyquist noise**), caused by the random thermal motion of electrons in any conductor at any temperature above absolute zero. It exists in every wire, every antenna, every receiver circuit, all the time, with no way to eliminate it short of cooling the system to absolute zero (impossible in practice). Its power is given by a remarkably simple formula:

```
  N = k · T · B

  where:
    N = noise power, in watts
    k = Boltzmann's constant ≈ 1.38 x 10^-23 joules/kelvin
    T = absolute temperature, in kelvin
    B = bandwidth, in Hz  (Chapter 16's precise definition, doing real work here)
```

**Worked example:** thermal noise power in a receiver at room temperature (T ≈ 290 K, a standard reference temperature in RF engineering) with a 20 MHz (Wi-Fi channel) bandwidth:

```
  N = (1.38 x 10^-23) x 290 x (20 x 10^6)
    = 1.38 x 10^-23 x 290 x 20,000,000
    = 1.38 x 10^-23 x 5.8 x 10^9
    ≈ 8.0 x 10^-14 watts
    ≈ -131 dBm   (dBm = dB relative to 1 milliwatt, standard RF unit)
```

Notice the formula directly: **wider bandwidth admits more thermal noise power.** This is a genuinely important, non-obvious fact — it means Chapter 16's "just use more bandwidth for more throughput" advice has a real cost attached: a wider channel doesn't just admit more signal, it admits proportionally more noise too, which is one reason simply making channels wider and wider eventually yields diminishing returns (fully quantified in Chapter 18).

Beyond thermal noise, real systems also experience **quantization noise** (small errors introduced when an analog signal is digitized — relevant background for Chapter 18's dial-up example) and **impulse noise** (short, high-energy bursts from sources like electrical switching, lightning, or motors) — both mentioned here for completeness; thermal noise is the one with a clean formula and the one Shannon's theorem in Chapter 18 fundamentally assumes.

---

## 10. Signal-to-Noise Ratio (SNR), Defined Precisely

Everything above — attenuation, interference, and noise — ultimately needs to be compared against something: how strong is the wanted signal, relative to everything unwanted competing with it? That single comparison is **SNR**.

```
  SNR = P_signal / P_noise            (as a plain ratio, unitless)

  SNR (dB) = 10 · log10(P_signal / P_noise)
```

| SNR (dB) | Plain ratio | Practical meaning |
|---|---|---|
| 0 dB | 1:1 | Signal and noise equal power — essentially unusable |
| 10 dB | 10:1 | Signal 10x noise power — poor but potentially workable with robust modulation |
| 20 dB | 100:1 | Signal 100x noise — usable for moderate-order modulation |
| 30 dB | 1,000:1 | Good connection — supports higher-order QAM (Chapter 16) |
| 40+ dB | 10,000:1+ | Excellent connection — supports highest-order modulation (256-QAM, 1024-QAM) |

This table is precisely why Chapter 16's adaptive modulation exists: a Wi-Fi client continuously estimates its current SNR and selects the highest QAM order the current SNR can reliably support, per Chapter 16 Section 5's tradeoff.

---

## 11. Worked Example: Computing Real SNR

Combine Sections 5, 9, and 10 into one complete, realistic calculation: a Wi-Fi access point transmits at 20 dBm (100 milliwatts, a typical consumer router's transmit power) at 2.4 GHz, and a laptop sits 20 meters away.

**Step 1 — free-space path loss (Section 5):**

```
  d = 0.02 km, f = 2,400 MHz

  FSPL = 20 log10(0.02) + 20 log10(2400) + 32.44
       = 20 x (-1.70)   + 67.6           + 32.44
       = -34.0 + 67.6 + 32.44
       = 66.04 dB of loss
```

**Step 2 — received signal power:**

```
  Received power (dBm) = Transmit power (dBm) - Path loss (dB)
                        = 20 dBm - 66.04 dB
                        = -46.04 dBm
```

**Step 3 — noise floor (Section 9's result for a 20 MHz channel):**

```
  Noise floor ≈ -131 dBm + typical receiver noise figure (~6 dB, real hardware
                is never perfectly ideal) ≈ -125 dBm (a realistic practical figure)
```

**Step 4 — SNR:**

```
  SNR (dB) = Received signal (dBm) - Noise floor (dBm)
           = -46.04 - (-125)
           = 78.96 dB   ...

  (Real Wi-Fi SNR figures reported by devices are typically in the
  20-50 dB range in practice, because Step 1's idealized free-space
  model ignores walls, furniture, reflections, and other real losses
  Chapter 23 and Volume 13 will address -- but the CALCULATION METHOD
  above -- transmit power minus path loss minus noise floor -- is
  exactly the real method used, just with more realistic loss inputs.)
```

The number above is intentionally optimistic (a truly clear line of sight with no obstacles is rare indoors); the value of this worked example is the *method* — every real SNR estimate in real networking hardware is computed by exactly this three-step chain: transmit power, minus path/cable/fiber loss (Sections 2-5), compared against the receiver's noise floor (Section 9).

---

## 12. Real Fixes: Repeaters, Shielding, Channel Selection, Coding

Each problem in this chapter has a corresponding, standard real-world fix:

| Problem | Fix | Where covered |
|---|---|---|
| Attenuation (copper) | Repeaters/amplifiers regenerate the digital signal (Chapter 14 Section 4) before it degrades too far; length limits enforced by cabling standards | Chapter 21 |
| Attenuation (fiber) | Optical amplifiers (EDFAs) or full optical-electrical-optical regeneration for very long runs | Chapter 22 |
| Attenuation (wireless) | More transmit power, better antennas, additional access points to shorten distances | Volume 13 |
| Crosstalk | Twisting wire pairs, shielding (STP cable) | Chapter 21 |
| Radio interference (e.g. microwaves) | Moving to less-crowded channels/bands, using 5/6 GHz instead of 2.4 GHz, physical shielding | Volume 13 |
| Thermal/background noise | Cannot be eliminated — mitigated by lowering modulation order, adding error correction, increasing signal power | Chapters 16, 18, 19, 20 |
| All of the above, combined | Forward error correction and error detection to survive whatever noise gets through anyway | Chapters 19-20 |

---

## 13. Why Long Ethernet Runs Need Repeaters — the Full Explanation

This chapter can now fully answer a question the TOC posed at the start of this volume: why does a standard twisted-pair Ethernet cable have a hard 100-meter length limit?

Three of this chapter's mechanisms compound together over distance:

1. **Attenuation (Section 3)** reduces signal amplitude, shrinking the noise margin (Chapter 14) available at the receiver.
2. **Crosstalk (Section 7)**, while mitigated by twisting, is never perfectly cancelled, and its relative impact grows as the wanted signal weakens from attenuation.
3. **Timing/propagation delay** (not a noise phenomenon, but a real constraint) matters for collision-detection mechanisms in older half-duplex Ethernet (Chapter 30) — a cable that's too long lets timing windows for detecting collisions expire before a collision signal can propagate back.

The 100-meter figure in the Ethernet standards (Chapter 21) is the point at which, for the specified maximum signaling frequency (necessary for gigabit-class speeds), all of these effects combined would push a standard-quality cable's worst-case attenuation and crosstalk past the safety margin the standard's designers required — not an arbitrary round number, but a calculated, tested engineering limit. Beyond 100 meters, you don't get a "somewhat worse" connection — you install a **repeater** (or simply a switch, which digitally regenerates the signal exactly as Chapter 14 Section 4 described) and start a fresh 100-meter run.

---

## 14. Common Misconceptions

- **"Interference and noise are the same thing."** They aren't, by this chapter's precise definitions — interference comes from an identifiable other source (another Wi-Fi network, a microwave, a neighboring wire) and can often be avoided; noise is fundamental, unavoidable, and has no single source to blame or move away from.
- **"A weak Wi-Fi signal (low signal strength/RSSI) always means a bad connection."** Not necessarily — what matters is SNR, not raw signal strength. A weak signal in a very quiet (low-noise) environment can have excellent SNR and work fine; a strong signal in a very noisy environment can have poor SNR and perform badly. Signal strength alone, without knowing the noise floor, tells you only half the story.
- **"Doubling distance halves your Wi-Fi signal."** Section 5 showed the real relationship is logarithmic, not linear — 10x the distance costs a fixed 20 dB, and dB losses don't correspond to simple halving/doubling of perceived signal in a linear sense.
- **"Shielded cables eliminate crosstalk and interference completely."** Shielding reduces both substantially but never to zero — some coupling always survives, which is why even well-shielded, professionally-installed cabling still has published (nonzero) crosstalk specifications.

---

## 15. Hands-On Experiment

Measure real SNR and watch real interference happen, using your own devices.

```bash
# macOS: view live SNR-relevant figures for your Wi-Fi connection
/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I
# Look for "agrCtlRSSI" (signal, in dBm -- more negative = weaker) and
# "agrCtlNoise" (noise floor, in dBm). Subtract noise from signal
# (both are negative dBm numbers) to get your real, live SNR in dB --
# exactly Section 10's formula, applied to your own network right now.

# Linux:
iw dev wlan0 station dump
# Look for "signal:" (dBm)

# Experiment: run a microwave oven near your router/laptop while
# running a continuous ping to your router, and watch for increased
# latency or packet loss during the oven's operation:
ping -i 0.2 <your-router-ip>
# (Try this only if you have a 2.4 GHz-band device to test with --
# 5 GHz devices should show no effect, exactly as Section 8 predicts.)
```

---

## 16. Code: dB, Attenuation, and SNR in Go

```go
package main

import (
	"fmt"
	"math"
)

// dBFromPowerRatio implements Chapter 16 Section 3's dB formula.
func dBFromPowerRatio(pOut, pIn float64) float64 {
	return 10 * math.Log10(pOut/pIn)
}

// freeSpacePathLoss implements Section 5's FSPL formula (d in km, f in MHz).
func freeSpacePathLoss(dKm, fMHz float64) float64 {
	return 20*math.Log10(dKm) + 20*math.Log10(fMHz) + 32.44
}

// thermalNoisePowerWatts implements Section 9's N = kTB formula.
func thermalNoisePowerWatts(tempKelvin, bandwidthHz float64) float64 {
	const boltzmann = 1.38e-23
	return boltzmann * tempKelvin * bandwidthHz
}

// wattsToDBm converts a power in watts to dBm (dB relative to 1 milliwatt).
func wattsToDBm(watts float64) float64 {
	milliwatts := watts * 1000
	return 10 * math.Log10(milliwatts/1)
}

func main() {
	fmt.Println("--- Section 5: Free-space path loss, 2.4 GHz Wi-Fi ---")
	for _, dMeters := range []float64{10, 20, 100} {
		loss := freeSpacePathLoss(dMeters/1000, 2400)
		fmt.Printf("  distance=%5.0fm  FSPL=%.2f dB\n", dMeters, loss)
	}

	fmt.Println("\n--- Section 9: Thermal noise floor, 20 MHz channel, room temp ---")
	noiseWatts := thermalNoisePowerWatts(290, 20_000_000)
	fmt.Printf("  Noise power: %.3e W  (%.1f dBm)\n", noiseWatts, wattsToDBm(noiseWatts))

	fmt.Println("\n--- Section 11: Full SNR worked example ---")
	txPowerDBm := 20.0
	loss := freeSpacePathLoss(20.0/1000, 2400)
	rxPowerDBm := txPowerDBm - loss
	noiseFloorDBm := wattsToDBm(noiseWatts) + 6 // + realistic receiver noise figure
	snrDB := rxPowerDBm - noiseFloorDBm

	fmt.Printf("  TX power:      %.2f dBm\n", txPowerDBm)
	fmt.Printf("  Path loss:     %.2f dB\n", loss)
	fmt.Printf("  RX power:      %.2f dBm\n", rxPowerDBm)
	fmt.Printf("  Noise floor:   %.2f dBm\n", noiseFloorDBm)
	fmt.Printf("  SNR:           %.2f dB\n", snrDB)
}
```

```
Output (abridged):

--- Section 5: Free-space path loss, 2.4 GHz Wi-Fi ---
  distance=   10m  FSPL=60.04 dB
  distance=   20m  FSPL=66.04 dB
  distance=  100m  FSPL=80.04 dB

--- Section 9: Thermal noise floor, 20 MHz channel, room temp ---
  Noise power: 8.004e-14 W  (-130.97 dBm)

--- Section 11: Full SNR worked example ---
  TX power:      20.00 dBm
  Path loss:     66.04 dB
  RX power:      -46.04 dBm
  Noise floor:   -124.97 dBm
  SNR:           78.93 dB
```

Change the distance or transmit power constants and watch SNR fall in exactly the pattern Sections 5 and 10 describe.

---

## 17. What This Chapter Simplified

- The free-space path loss formula (Section 5) assumes a perfect, unobstructed line of sight — real indoor Wi-Fi suffers additional losses from walls, furniture, and multipath reflections that this simplified model omits entirely, which is why Section 11's optimistic result doesn't match typical real-world SNR readings.
- This chapter treated attenuation, interference, and noise as cleanly separable, but real systems often experience combinations that interact in more complex ways than a simple sum of dB losses.
- The thermal noise formula (Section 9) is exact physics, but real receiver hardware adds its own additional noise (captured loosely here as a "noise figure" fudge factor) from imperfect amplifiers and components — a full treatment belongs to RF engineering texts, not this course.
- Interference sources beyond microwave ovens and crosstalk (e.g., Bluetooth devices, cordless phones, neighboring Wi-Fi networks, radar in some 5 GHz sub-bands) exist in real deployments and are only briefly gestured at here.

---

## Interview Questions & Model Answers

**Beginner: "What is the difference between attenuation, interference, and noise?"**

Attenuation is a signal simply getting weaker as it travels, due to the medium itself (resistance in copper, absorption/scattering in fiber, spreading out in open air) — no other signal is involved. Interference is unwanted energy from an identifiable other source, like another Wi-Fi network or a microwave oven leaking energy into the same frequency band. Noise is unwanted, essentially random energy with no single identifiable source, arising from fundamental physics like the thermal motion of electrons — it exists in every channel all the time and can never be fully eliminated, only reduced or compensated for.

**Intermediate: "Why does a microwave oven specifically disrupt 2.4 GHz Wi-Fi but not 5 GHz Wi-Fi?"**

Microwave ovens generate radiation at approximately 2.45 GHz to heat food, because that frequency is well-absorbed by water molecules. That frequency happens to sit almost exactly in the middle of the 2.4 GHz Wi-Fi band (roughly 2.400-2.4835 GHz). Ovens are shielded but never perfectly, so some of that high-power radiation leaks out and directly overlaps Wi-Fi's operating frequencies, acting as strong interference. 5 GHz Wi-Fi operates at an entirely different frequency range that ovens don't emit into, so there's no frequency overlap and therefore no interference from this particular source.

**Advanced: "Explain, with the actual formula, why doubling the distance from a Wi-Fi access point doesn't simply halve the received signal strength."**

Free-space path loss follows FSPL(dB) = 20·log10(d) + 20·log10(f) + constant. Because the distance term uses a logarithm, doubling the distance always adds the same fixed amount of dB loss (about 6 dB for a doubling, or 20 dB for a 10x increase), regardless of the starting distance — it's a multiplicative, not additive, relationship in linear terms. In linear power terms, this logarithmic relationship means received power actually falls off with the square of distance (an "inverse-square law"), not linearly — so going from 10m to 20m away doesn't cost you "half your signal" in a simple sense, it costs you a specific, calculable dB penalty that compounds with every further doubling of distance.

---

## Exercises

### Easy

1. Explain, in one or two sentences each, the practical difference between interference and noise, using a real example of each.
2. A cable has an attenuation rating of 0.3 dB/meter. Calculate the total attenuation over a 50-meter run.
3. Why does a microwave oven interfere with 2.4 GHz Wi-Fi but not with a wired Ethernet connection?

### Medium

4. Using the FSPL formula from Section 5, calculate the path loss for a 5 GHz signal (5,000 MHz) at the same 20-meter distance used in Section 11's 2.4 GHz example. Compare the two results and explain the difference.
5. A receiver measures a signal at -50 dBm and a noise floor at -90 dBm. Calculate the SNR in dB, and using Chapter 16 Section 5's QAM table plus this chapter's Section 10 table, argue whether this connection could reliably support 256-QAM.
6. Explain why NEXT (near-end crosstalk) is generally worse than FEXT (far-end crosstalk), using Section 2's attenuation concept.

### Hard

7. Using the thermal noise formula N = kTB, explain quantitatively (with a calculation) why doubling a channel's bandwidth (e.g., from 20 MHz to 40 MHz Wi-Fi channels) increases the noise floor, and explain why this partially offsets the throughput gains Chapter 16 attributed to wider bandwidth.
8. A network engineer needs to run a cable 150 meters between two buildings and is deciding between adding an Ethernet repeater partway along a copper run, or switching to fiber optic cable for the entire run. Using Sections 3, 4, and 13's numbers, write a short technical justification for one approach over the other.

---

## Summary

| Term | Meaning |
|---|---|
| Attenuation | Signal weakening from traveling through a medium (resistance, absorption, scattering, spreading) |
| Skin effect | Higher-frequency current concentrating near a conductor's surface, increasing effective resistance |
| Free-space path loss (FSPL) | Signal loss from geometric spreading in open space; grows with log(distance) and log(frequency) |
| Interference | Unwanted energy from an identifiable other source (crosstalk, another radio, a microwave oven) |
| Crosstalk (NEXT/FEXT) | Interference between adjacent wires from electromagnetic coupling |
| Noise | Unwanted, essentially random energy with no identifiable source; e.g. thermal/Johnson-Nyquist noise |
| Thermal noise formula | N = kTB — noise power grows with temperature and bandwidth |
| SNR | Signal-to-Noise Ratio; SNR(dB) = 10·log10(signal power / noise power) |
| dBm | Decibels relative to 1 milliwatt; standard unit for absolute RF signal/noise power |

This chapter gave every real-world degradation a signal suffers a name, a cause, and a formula — and combined them all into one number, SNR, that finally captures "how good is this channel, really?" Chapter 18, **Shannon's Limit**, takes that single number and combines it with Chapter 16's bandwidth to answer the question this entire volume has been building toward: given a channel's bandwidth and its SNR, exactly how many bits per second can it possibly carry — no more, ever, regardless of how clever the engineering gets?

