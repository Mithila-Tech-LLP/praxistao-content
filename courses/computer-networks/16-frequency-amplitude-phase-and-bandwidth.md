# Chapter 16: Frequency, Amplitude, Phase, and Bandwidth

*"Two modems can share the exact same wire, the exact same voltage range, and the exact same amount of time — and still send wildly different amounts of data, because one of them figured out how to say more with each wave."*

---

## Table of Contents

1. [The Problem: One Bit Per Symbol Is Leaving Performance on the Table](#1-the-problem-one-bit-per-symbol-is-leaving-performance-on-the-table)
2. [Frequency, Precisely](#2-frequency-precisely)
3. [Amplitude, Precisely](#3-amplitude-precisely)
4. [Phase, Precisely](#4-phase-precisely)
5. [Combining Amplitude and Phase: QAM](#5-combining-amplitude-and-phase-qam)
6. [Symbols vs. Bits — Baud Rate vs. Bit Rate](#6-symbols-vs-bits--baud-rate-vs-bit-rate)
7. [Bandwidth, Defined Precisely (in Hz)](#7-bandwidth-defined-precisely-in-hz)
8. [Why a Sharp Digital Pulse Needs So Much Bandwidth](#8-why-a-sharp-digital-pulse-needs-so-much-bandwidth)
9. [Bandwidth (Hz) vs. Throughput (bps) — The Crucial Distinction](#9-bandwidth-hz-vs-throughput-bps--the-crucial-distinction)
10. [The Nyquist Rate — A Preview of the Noiseless Limit](#10-the-nyquist-rate--a-preview-of-the-noiseless-limit)
11. [Worked Example: A Real Wi-Fi Channel](#11-worked-example-a-real-wi-fi-channel)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Hands-On Experiment](#13-hands-on-experiment)
14. [Code: Bits Per Symbol and QAM Constellations in Go](#14-code-bits-per-symbol-and-qam-constellations-in-go)
15. [What This Chapter Simplified](#15-what-this-chapter-simplified)
16. [Interview Questions & Model Answers](#interview-questions--model-answers)
17. [Exercises](#exercises)
18. [Summary](#summary)

---

## 1. The Problem: One Bit Per Symbol Is Leaving Performance on the Table

Chapter 15 built three modulation schemes, each varying exactly one property of a carrier wave — amplitude (ASK), frequency (FSK), or phase (PSK) — and each, in its simplest binary form, sending exactly **one bit per transmitted symbol** (one wave "event"). QPSK, introduced briefly at the end of that chapter, showed a hint of something better: by using four distinct phase states instead of two, it squeezed **two bits into every symbol** instead of one, without needing any more bandwidth.

That single fact should make you suspicious that one bit per symbol was never a hard limit — it was just the simplest place to start. The real question this chapter answers: **if a wave has three independent, simultaneously-adjustable properties (amplitude, frequency, phase), why stop at varying just one, and why stop at just two states per property?**

The answer is a modulation scheme called **Quadrature Amplitude Modulation (QAM)**, which varies amplitude *and* phase together, and which underlies nearly every high-speed digital communication system in use today — Wi-Fi, cable modems, DSL, LTE, 5G. To understand it precisely, this chapter first nails down exact, rigorous definitions of frequency, amplitude, and phase (used loosely so far), then builds QAM from them, and finally gives **bandwidth** — a word used constantly and incorrectly in everyday speech — its real, technical meaning.

---

## 2. Frequency, Precisely

**Definition:** frequency (f) is the number of complete wave cycles that occur per second, measured in **hertz (Hz)**. One hertz means one complete cycle per second. Frequency and **period** (T, the time for one complete cycle) are reciprocals of each other:

```
  f = 1 / T          T = 1 / f

  Example: a wave with period T = 0.0005 seconds (0.5 milliseconds)
  has frequency f = 1 / 0.0005 = 2,000 Hz = 2 kHz
```

**Wavelength**, the physical distance a wave travels during one complete cycle, connects frequency to the speed the wave travels through its medium:

```
  λ = v / f

  where λ (lambda) = wavelength, v = wave speed in the medium, f = frequency
```

For electromagnetic waves in a vacuum (or very close to it, in air), v = c ≈ 300,000,000 m/s (the speed of light). This is why higher-frequency radio waves have shorter wavelengths — and why, as Chapter 15 mentioned, an antenna's practical size is tied to the wavelength (and therefore frequency) it's designed to radiate efficiently. A 2.4 GHz Wi-Fi wave has wavelength λ = 300,000,000 / 2,400,000,000 ≈ 0.125 m (12.5 cm) — which is why Wi-Fi antennas are a few centimeters long, not a few meters.

**The electromagnetic spectrum, and where networking technologies sit on it:**

| Frequency range | Common name | Used for |
|---|---|---|
| 300 Hz - 3,400 Hz | Voice band | Telephone lines (analog voice, dial-up modems) |
| 0 - 1.1 MHz | — | DSL over copper phone lines |
| 3 kHz - 30 MHz | HF (High Frequency) radio | Shortwave radio, some long-range comms |
| 2.4 GHz, 5 GHz, 6 GHz | ISM/UNII bands | Wi-Fi, Bluetooth, some cordless phones, microwave ovens (2.45 GHz) |
| 600 MHz - 6 GHz (Sub-6) | Cellular Sub-6 | 4G LTE, 5G Sub-6 |
| 24-100 GHz | mmWave | 5G mmWave (very high speed, very short range) |
| ~190-196 THz | Infrared/near-infrared light | Fiber optic transmission (Chapter 22) |

Notice fiber's frequency is measured in **terahertz** (10¹² Hz) — roughly 100,000 times higher than Wi-Fi's frequency — a fact that will matter enormously in Chapter 18's discussion of why fiber can carry vastly more data than radio or copper.

---

## 3. Amplitude, Precisely

**Definition:** amplitude (A) is the maximum displacement (or strength) of a wave from its resting/zero value — informally, how "big" or "loud" or "strong" the signal is. For an electrical signal, amplitude is typically measured as peak voltage; for a radio wave, it corresponds to signal power/strength.

Amplitude is most often discussed in networking not as a raw number, but as a *ratio*, expressed in **decibels (dB)** — a logarithmic scale that will become essential in Chapter 17 for discussing noise and attenuation:

```
  dB = 10 · log10(P2 / P1)          (for power ratios)

  Example: if a signal's power drops from 1 mW to 0.1 mW (a 10x drop),
  dB = 10 · log10(0.1/1) = 10 · log10(0.1) = 10 · (-1) = -10 dB

  A drop of 10 dB always means "1/10th the power," regardless of the
  starting power level -- this is the entire point of a logarithmic scale:
  ratios, not absolute differences, stay constant.
```

Chapter 17 develops this fully; for now, just fix the idea that amplitude is a precisely measurable quantity (voltage or power), and that ratios of amplitude are conventionally expressed in dB because real-world signal strengths vary over such enormous ranges (factors of millions or more) that a logarithmic scale is far more practical than a linear one.

---

## 4. Phase, Precisely

**Definition:** phase (φ) describes the position of a point in time within a wave's repeating cycle, usually measured in degrees (0°-360°) or radians (0-2π), relative to some reference point in time.

```
  Two waves of the same frequency and amplitude, but different phase:

  Wave A (phase = 0°):        Wave B (phase = 90°, "quarter cycle ahead"):

    /\      /\      /\          _    _    _
   /  \    /  \    /  \        / \  / \  / \
  /    \  /    \  /    \      |   \/   \/   |
 /      \/      \/      \    _|             |_
                                (shifted left/earlier by a quarter cycle)
```

Phase only makes sense as a *relative* concept — "90° ahead of what?" A receiver needs a reference (usually recovered from the incoming signal itself, or from a known pilot/reference signal) against which to measure the phase of each received symbol. This synchronization requirement is a real engineering challenge in PSK and QAM receivers, briefly noted here and expanded in the "what this chapter simplified" section.

---

## 5. Combining Amplitude and Phase: QAM

**The idea:** Chapter 15 varied phase alone (PSK) to get multiple bits per symbol. QAM goes further: vary **both amplitude and phase simultaneously**, so each symbol is defined by a *point in a two-dimensional plane* — its distance from the center (amplitude) and its angle (phase) — rather than a single value on a circle or a single loudness level.

A QAM constellation diagram plots every valid symbol as a point on this plane:

```
  16-QAM constellation (16 possible symbol positions,
  arranged in a 4x4 grid -- each point differs from its
  neighbors in BOTH amplitude (distance from center)
  AND phase (angle from center)):

         .   .   .   .
         .   .   .   .
         .   .   .   .        <- each "." is one of 16 valid
         .   .   .   .           (amplitude, phase) combinations,
                                  each representing a unique 4-bit value
```

Because each point is a unique, distinguishable combination of amplitude and phase, and there are 16 such points, each single transmitted symbol can represent **log₂(16) = 4 bits** — four times the information of a single ASK/FSK/BPSK symbol, using the same amount of time and (with care) a comparable amount of bandwidth.

**The general formula connecting the number of constellation points to bits per symbol:**

```
  bits per symbol = log2(M)

  where M = the number of distinct symbol states ("QAM order")
```

| QAM order (M) | Bits per symbol | Used in |
|---|---|---|
| 4-QAM (= QPSK) | log2(4) = 2 | Robust/long-range Wi-Fi, satellite |
| 16-QAM | log2(16) = 4 | Wi-Fi, cable modems (DOCSIS), LTE |
| 64-QAM | log2(64) = 6 | Wi-Fi, DOCSIS, digital cable TV |
| 256-QAM | log2(256) = 8 | Wi-Fi 5/6 (short range, high SNR), DOCSIS 3.1 |
| 1024-QAM | log2(1024) = 10 | Wi-Fi 6 (very short range, excellent SNR) |
| 4096-QAM | log2(4096) = 12 | Wi-Fi 6E/7, high-end DOCSIS 3.1 |

**The tradeoff this table hides, and that Chapter 14 flagged as the central tension of this entire subject:** as M grows, the constellation points must pack more tightly into the same amplitude/phase plane, so the "noise margin" between adjacent valid points shrinks. A router or modem doesn't get to pick 4096-QAM for free — it needs a very clean, high-SNR (Chapter 17) connection to reliably tell such closely-spaced points apart. This is exactly why your Wi-Fi router automatically drops to a lower-order, more-robust modulation (16-QAM, or even QPSK) as you walk farther from it and your signal gets noisier — it is dynamically trading bits-per-symbol for reliability, symbol by symbol, in real time. This behavior is called **adaptive modulation and coding**, and it's one of the most important practical consequences of everything built in this chapter.

---

## 6. Symbols vs. Bits — Baud Rate vs. Bit Rate

This distinction is subtle enough that it's worth its own section, because it is one of the most common sources of confusion in networking, and it directly sets up bandwidth's precise definition.

- A **symbol** is one discrete transmitted "event" — one specific combination of amplitude/frequency/phase, held for a fixed slice of time.
- The **symbol rate** (or **baud rate**, named after telegraph pioneer Émile Baudot) is how many symbols are transmitted per second, measured in **baud** (symbols/second).
- The **bit rate** is how many bits are transmitted per second, measured in **bits per second (bps)**.

They are related by:

```
  bit rate = symbol rate  x  bits per symbol
           = baud          x  log2(M)
```

**Worked example:** a modem transmitting at 2,400 baud (2,400 symbols per second) using 16-QAM (4 bits per symbol, from Section 5's table):

```
  bit rate = 2,400 baud  x  4 bits/symbol  =  9,600 bps
```

This is precisely how real dial-up modems achieved bit rates far higher than their baud rate ever changed — the actual telephone-line symbol rate for most dial-up standards stayed around 2,400-3,200 baud (limited by the phone line's bandwidth, as Section 10 will show), while the *bit rate* climbed from 2,400 bps in early modems to 33,600 bps in late-1990s modems purely by using progressively higher-order modulation (more bits per symbol) — not by sending symbols any faster. Chapter 18 revisits this exact fact with the full numbers.

---

## 7. Bandwidth, Defined Precisely (in Hz)

Here is the definition this entire course has been building toward, stated with no hedging:

**Bandwidth is the range of frequencies a channel can carry, or that a signal occupies — measured in hertz (Hz).** It is a property of *frequency*, full stop. It has nothing to do, in its technical definition, with how many bits per second you can send — that connection exists, and it's crucial (Section 9 and, fully, Chapter 18), but it is a *consequence*, not the definition.

```
  Bandwidth = highest frequency the channel/signal uses
              MINUS
              lowest frequency the channel/signal uses

  B = f_high - f_low
```

**Worked examples:**

```
  Telephone voice band:
    f_low = 300 Hz, f_high = 3,400 Hz
    Bandwidth = 3,400 - 300 = 3,100 Hz  (~3 kHz, commonly rounded to "3 kHz")

  A single Wi-Fi 20 MHz channel:
    (the "20 MHz" in "20 MHz channel" IS the bandwidth, by definition)
    Bandwidth = 20,000,000 Hz = 20 MHz

  Standard cable TV / DOCSIS channel:
    Bandwidth = 6 MHz (US) or 8 MHz (Europe/DVB-C)
```

Notice: none of these numbers say anything, by themselves, about bits per second. A 20 MHz Wi-Fi channel using QPSK carries far fewer bits per second than the same 20 MHz channel using 256-QAM — same bandwidth, wildly different throughput, because bandwidth measures *frequency range*, and throughput depends on *how many bits you pack into each symbol within that frequency range* (Section 6), and on how fast you can transmit symbols within that range (Section 10).

---

## 8. Why a Sharp Digital Pulse Needs So Much Bandwidth

Chapter 15, Section 5 asserted (without proof) that a sharp square wave is equivalent to summing many sine waves of different frequencies. This section makes that concrete enough to trust, without requiring calculus.

**Intuitive version:** imagine trying to build the shape of a square wave using only smooth, rounded sine-wave "bricks." A single sine wave is much too round to make a flat top and a sharp corner. But if you add a second, higher-frequency sine wave on top of the first, the corners get a little sharper and the top gets a little flatter. Keep adding higher and higher frequency sine waves, each with a smaller contribution, and the sum gets closer and closer to a true flat-topped, sharp-cornered square wave. This is a real mathematical fact (a **Fourier series**), not just an analogy — any repeating waveform, including a perfect square wave, can be expressed as a sum of sine waves at specific frequencies and amplitudes.

```
  1 sine wave:           +2nd (3x freq):        +3rd (5x freq):     ...approaching
     __                     __                     __                a square wave:
    /  \                   /''\                   /''\                ___
   /    \        -->      |    |        -->      |    |      -->     |   |
  /      \                |    |                 |    |              |   |
                           \__/                    \__/               |___|
```

**The practical consequence:** the sharper you want a digital pulse's edges to be (the faster it needs to transition from 0 to 1), the *more* high-frequency sine-wave components are needed to build that sharp edge — which means the signal occupies more bandwidth. A channel with limited bandwidth (like the phone line's 3 kHz) can only faithfully carry the *low-frequency* components of a signal; it acts as a filter that strips away everything above roughly its upper cutoff frequency, which is exactly why a raw digital square wave sent unmodified down a bandlimited phone line arrives blurred and unreadable, as Chapter 15 asserted from the start.

This fact — that faster transitions require more bandwidth — is also *why* symbol rate (Section 6) is fundamentally capped by a channel's bandwidth: transmitting more symbols per second means each symbol's voltage/amplitude/phase has less time to settle before the next one begins, which (by the reasoning above) requires more bandwidth to represent without the symbols blurring into each other.

---

## 9. Bandwidth (Hz) vs. Throughput (bps) — The Crucial Distinction

Now the everyday-language collision can be named directly and corrected permanently.

**Everyday, technically incorrect usage:** "I have 500 Mbps internet bandwidth" — using "bandwidth" to mean data throughput, in bits per second.

**Technically correct usage:** bandwidth is a frequency range, in Hz (Section 7). What an ISP sells you, measured in Mbps, is **throughput** (or **data rate**, or **bit rate**) — the actual number of bits successfully delivered per second, which depends on:

```
  throughput  ≈  bandwidth (Hz)  x  bits per symbol (modulation, Section 6)
                 (further reduced by noise, protocol overhead, etc. --
                  Chapter 17 and 18 make this precise, including the
                  hard upper bound Shannon's theorem places on it)
```

| Correct term | Unit | What it measures |
|---|---|---|
| Bandwidth | Hz (hertz) | Range of frequencies a channel/signal occupies |
| Throughput / bit rate / data rate | bps (bits/second), or Mbps, Gbps | Actual bits delivered per second |
| Symbol rate / baud rate | baud (symbols/second) | How many discrete symbols are transmitted per second |

**Worked comparison showing why conflating them is actively misleading:** two channels can have *identical* bandwidth and wildly different throughput:

```
  Channel A: 20 MHz bandwidth, QPSK (2 bits/symbol), reasonable symbol rate
             --> lower throughput

  Channel B: 20 MHz bandwidth (SAME), 256-QAM (8 bits/symbol), same symbol rate
             --> 4x the throughput of Channel A

  Same bandwidth. Same frequency range. Completely different bps.
```

This is precisely why "bandwidth" in the strict engineering sense is not what changed when your ISP upgraded your connection's advertised speed — usually your bandwidth (frequency allocation) stayed the same or grew only modestly, while improvements in modulation efficiency (higher-order QAM), signal processing, and reduced noise did most of the work of raising your actual bps.

---

## 10. The Nyquist Rate — A Preview of the Noiseless Limit

Before Claude Shannon's full noisy-channel theorem (Chapter 18), telegraph engineer Harry Nyquist derived, in 1928, the maximum symbol rate a channel of a given bandwidth can support **in the total absence of noise**:

```
  Nyquist symbol rate:   R_max = 2 · B

  where B = the channel's bandwidth in Hz
```

**Worked example:** a channel with 3,000 Hz of bandwidth (roughly the phone line's voice band from Section 7) can, at most, carry:

```
  R_max = 2 x 3,000 = 6,000 symbols per second (baud), in a noiseless channel
```

Combined with Section 6's formula, this gives a noiseless upper bound on bit rate:

```
  bit rate (noiseless) = R_max x log2(M) = 2B x log2(M)
```

**Why this matters, and why it's only half the story:** Nyquist's formula says how many symbols per second bandwidth alone permits — it says nothing about noise. In the real world, noise limits how many distinguishable amplitude/phase levels (M) you can safely use before adjacent constellation points become indistinguishable (Section 5's tradeoff). Chapter 18 combines Nyquist's bandwidth-only limit with a noise-aware limit on M to arrive at Shannon's full formula — the true, complete answer to "how much data can possibly fit through this channel." Everything in this chapter (frequency, amplitude, phase, QAM, symbol rate, bandwidth) is the exact vocabulary Chapter 18's formula is built from.

---

## 11. Worked Example: A Real Wi-Fi Channel

Put the whole chapter together on a single, realistic modern example.

**Given:**
- A Wi-Fi 5 (802.11ac) channel with 80 MHz bandwidth (a real, common Wi-Fi 5 channel width)
- Using 256-QAM (8 bits/symbol, from Section 5's table) — realistic for a device close to the router with an excellent SNR
- A real 802.11ac system doesn't run at the raw Nyquist rate — it uses a specific, standardized OFDM symbol structure — but this simplified calculation illustrates the scaling relationship correctly

```
  Bandwidth (B)         = 80 MHz = 80,000,000 Hz
  Bits per symbol (QAM) = 8  (256-QAM)

  Rough noiseless symbol-rate ceiling (Nyquist, Section 10):
     R_max = 2 x 80,000,000 = 160,000,000 symbols/sec (theoretical upper bound)

  Rough resulting bit-rate ceiling:
     bit rate = R_max x 8 = 1,280,000,000 bps = 1.28 Gbps
```

Real 802.11ac (with its actual OFDM subcarrier structure, guard intervals, and multiple spatial streams via MIMO) achieves headline rates in this same general multi-hundred-Mbps-to-several-Gbps range for wide channels and high-order QAM — the simplified Nyquist-based estimate above lands in the right neighborhood precisely because it captures the two real levers (bandwidth and bits-per-symbol) that actually determine throughput, even though it skips the full OFDM engineering detail.

**Now watch what happens as the device walks away from the router and SNR drops** (Chapter 17 explains why SNR drops with distance): the router's adaptive modulation (Section 5) steps down — 256-QAM → 64-QAM → 16-QAM → QPSK — cutting bits-per-symbol from 8 down to 2, a 4x drop in throughput, using the *exact same 80 MHz of bandwidth* the whole time. This single worked example is the entire chapter's content, applied to hardware you likely own right now.

---

## 12. Common Misconceptions

- **"More bandwidth always means more speed."** More bandwidth (Hz) *permits* more speed, by giving room for a higher symbol rate (Section 10) and more frequency-domain room for higher-order modulation to work reliably — but the actual speed also depends heavily on modulation order and SNR (Section 5, Chapter 17). Doubling bandwidth without improving anything else roughly doubles throughput; doubling modulation order (with sufficient SNR) can do the same without touching bandwidth at all.
- **"Bandwidth is measured in Mbps."** No — that's throughput. Bandwidth is measured in Hz. This chapter exists specifically to correct this extremely common conflation (Section 9).
- **"A higher frequency signal automatically carries more data."** Frequency alone says nothing about data rate — a very high carrier frequency (like fiber's ~190 THz) with simple on/off modulation carries less data than a much lower-frequency channel with wide bandwidth and 256-QAM. What matters for capacity is bandwidth (the *range* around that carrier) and modulation order, not the carrier frequency itself, though in practice, very high carrier frequencies do make very wide absolute bandwidths easier to allocate (this is part of why fiber's frequency range gives it access to enormous absolute bandwidth — more in Chapter 18 and Chapter 22).
- **"256-QAM is strictly 'better' than QPSK."** It's better only when the channel's SNR supports it. On a noisy channel, using 256-QAM would produce a flood of misread symbols and a worse effective throughput than a robust, lower-order scheme that at least gets its (fewer) bits through correctly. This is exactly why adaptive modulation exists.

---

## 13. Hands-On Experiment

See adaptive modulation happening in real time on your own Wi-Fi connection.

```bash
# On macOS:
/System/Library/PrivateFrameworks/Apple80211.framework/Versions/Current/Resources/airport -I
# Look for "lastTxRate" and "rssi" -- move around your house and re-run this,
# and watch the negotiated PHY rate change as your distance (and therefore
# SNR, Chapter 17) from the router changes.

# On Linux:
iw dev wlan0 link
# Look for "tx bitrate" -- note the modulation scheme shown, e.g.
# "tx bitrate: 866.7 MBit/s VHT-MCS 9 80MHz short GI" -- the "MCS" (Modulation
# and Coding Scheme) index directly maps to a specific QAM order and coding
# rate; a lower MCS index means a more robust, lower-order modulation.

# On Windows (PowerShell):
netsh wlan show interfaces
# Look at "Receive rate (Mbps)" and "Signal" -- walk further from your
# router and re-run; watch both drop together.
```

Standing right next to your router versus three rooms away should show a dramatic difference in the reported transmit rate — that drop is Section 5's adaptive modulation and coding in action, stepping down through the QAM-order table as SNR degrades.

---

## 14. Code: Bits Per Symbol and QAM Constellations in Go

```go
package main

import (
	"fmt"
	"math"
)

// bitsPerSymbol implements Section 5's core formula: bits/symbol = log2(M)
func bitsPerSymbol(qamOrder int) float64 {
	return math.Log2(float64(qamOrder))
}

// throughput implements Section 9's approximate relationship:
// throughput ≈ symbol rate x bits per symbol
func throughputBps(symbolRateBaud float64, qamOrder int) float64 {
	return symbolRateBaud * bitsPerSymbol(qamOrder)
}

// nyquistSymbolRate implements Section 10's noiseless upper bound: 2B
func nyquistSymbolRate(bandwidthHz float64) float64 {
	return 2 * bandwidthHz
}

func main() {
	qamOrders := []int{4, 16, 64, 256, 1024, 4096}

	fmt.Println("QAM order -> bits per symbol (Section 5 table, reproduced):")
	for _, m := range qamOrders {
		fmt.Printf("  %5d-QAM -> %.0f bits/symbol\n", m, bitsPerSymbol(m))
	}

	fmt.Println("\nSection 11's Wi-Fi worked example, reproduced:")
	bandwidth := 80_000_000.0 // 80 MHz
	symbolRate := nyquistSymbolRate(bandwidth)
	fmt.Printf("  Bandwidth: %.0f Hz\n", bandwidth)
	fmt.Printf("  Nyquist symbol-rate ceiling: %.0f symbols/sec\n", symbolRate)

	for _, m := range []int{4, 16, 64, 256} {
		bps := throughputBps(symbolRate, m)
		fmt.Printf("  Using %4d-QAM (%.0f bits/symbol): %.2f Gbps ceiling\n",
			m, bitsPerSymbol(m), bps/1_000_000_000)
	}
}
```

```
Output:

QAM order -> bits per symbol (Section 5 table, reproduced):
      4-QAM -> 2 bits/symbol
     16-QAM -> 4 bits/symbol
     64-QAM -> 6 bits/symbol
    256-QAM -> 8 bits/symbol
   1024-QAM -> 10 bits/symbol
   4096-QAM -> 12 bits/symbol

Section 11's Wi-Fi worked example, reproduced:
  Bandwidth: 80000000 Hz
  Nyquist symbol-rate ceiling: 160000000 symbols/sec
     4-QAM (2 bits/symbol): 0.32 Gbps ceiling
    16-QAM (4 bits/symbol): 0.64 Gbps ceiling
    64-QAM (6 bits/symbol): 0.96 Gbps ceiling
   256-QAM (8 bits/symbol): 1.28 Gbps ceiling
```

This directly reproduces Section 11's numbers and Section 5's adaptive-modulation tradeoff: same bandwidth throughout, throughput scaling linearly with bits-per-symbol.

---

## 15. What This Chapter Simplified

- Real Wi-Fi and cellular systems don't use one giant single-carrier symbol stream — they use **OFDM (Orthogonal Frequency-Division Multiplexing)**, splitting the channel's bandwidth into many parallel narrow sub-carriers, each independently QAM-modulated. This chapter's single-carrier Nyquist calculation captures the right scaling relationship but skips OFDM's real structure (a topic for Volume 13's Wi-Fi chapters).
- MIMO (multiple antennas sending independent spatial streams simultaneously) can multiply throughput further, independent of everything in this chapter — also deferred to Volume 13.
- Real symbol timing recovery, phase synchronization between transmitter and receiver, and the specific "constellation shaping" used to reduce peak power are all real engineering challenges glossed over here for clarity.
- The dB formula in Section 3 was given for power ratios; a slightly different formula (20·log10 instead of 10·log10) applies when working directly with voltage or amplitude ratios rather than power — Chapter 17 uses the power-ratio form throughout for consistency.

---

## Interview Questions & Model Answers

**Beginner: "What is the difference between bandwidth and throughput?"**

Bandwidth is a range of frequencies, measured in hertz (Hz) — it describes how much of the frequency spectrum a channel or signal occupies. Throughput (or bit rate, or data rate) is the actual number of bits successfully delivered per second, measured in bits per second (bps). They're related — more bandwidth generally permits higher throughput — but they're not the same thing and don't have the same units. A 20 MHz Wi-Fi channel (bandwidth) might deliver anywhere from tens of Mbps to over a gigabit per second (throughput) depending entirely on the modulation scheme and signal quality in use.

**Intermediate: "Explain how QAM allows a system to send more bits per symbol than simple ASK, FSK, or PSK."**

QAM varies both the amplitude and the phase of a carrier wave simultaneously, so each transmitted symbol corresponds to a specific point in a two-dimensional amplitude/phase plane rather than just a position on a single amplitude scale, frequency choice, or phase circle. With M distinct constellation points, each symbol can represent log2(M) bits — for example, 16-QAM's 16 constellation points each encode 4 bits, four times what a simple binary ASK, FSK, or BPSK symbol carries. The tradeoff is that as M grows, the constellation points are packed more closely together, so the scheme becomes more sensitive to noise, which is why real systems dynamically choose a lower QAM order when the connection is noisy (adaptive modulation).

**Advanced: "A Wi-Fi client's reported link rate drops from 866 Mbps to 200 Mbps as it moves away from the router, even though it's still using the exact same 80 MHz channel. Explain what's actually changing, using precise terminology."**

The channel's bandwidth (80 MHz, a fixed range of frequencies) hasn't changed at all. What's changed is the router's adaptive modulation and coding scheme selection: as the client moves farther away, its signal-to-noise ratio drops, and the router steps down from a high-order QAM scheme (like 256-QAM, 8 bits per symbol) to a lower-order, more noise-resistant scheme (like 16-QAM or QPSK, 4 or 2 bits per symbol), because the closely-packed constellation points of high-order QAM would otherwise be misread as noise corrupts the signal. Since throughput is approximately symbol rate multiplied by bits per symbol, and the symbol rate (bounded by the fixed 80 MHz bandwidth) hasn't changed, the entire throughput drop comes from the reduction in bits per symbol — a direct, real-world instance of the bandwidth-vs-modulation-order distinction this chapter establishes.

---

## Exercises

### Easy

1. Convert a wave with a period of 0.001 seconds into its frequency in Hz.
2. A signal's power drops from 4 mW to 0.4 mW. Calculate the loss in dB using Section 3's formula.
3. State, in one sentence each, the difference between bandwidth, throughput, and symbol rate.

### Medium

4. A channel has 6 MHz of bandwidth (a real DOCSIS cable channel, from Section 7). Using the Nyquist formula, calculate the noiseless maximum symbol rate. Then calculate the resulting bit-rate ceiling if the system uses 64-QAM (6 bits/symbol).
5. Explain, using Section 8's Fourier reasoning, why increasing symbol rate (sending symbols faster) on a fixed-bandwidth channel eventually causes symbols to blur into each other.
6. Using the table in Section 5, calculate how much the maximum theoretical bit rate would drop (as a percentage) if a system had to step down from 1024-QAM to 64-QAM due to a noisy channel, assuming symbol rate stays constant.

### Hard

7. Research (or reason from Section 2's wavelength formula) why 5 GHz Wi-Fi has a shorter usable range than 2.4 GHz Wi-Fi, all else being equal, even though a higher frequency doesn't inherently mean "weaker." (Hint: think about free-space signal spreading and obstacle penetration, both of which Chapter 17 will formalize.)
8. A designer proposes a communication system using 65,536-QAM (16 bits per symbol) to maximize throughput on a fixed-bandwidth channel. Using the concepts from Section 5 and Section 12, argue for or against this design choice for a mobile phone that will frequently be used far from a cell tower, in areas with significant electrical interference.

---

## Summary

| Term | Meaning |
|---|---|
| Frequency (f) | Cycles per second, in Hz; f = 1/T |
| Wavelength (λ) | Physical distance per cycle; λ = v/f |
| Amplitude (A) | Signal strength/magnitude; ratios expressed in dB |
| Phase (φ) | Position within a wave's cycle, relative to a reference |
| QAM | Modulation varying amplitude and phase together; bits/symbol = log2(M) |
| Symbol / baud rate | Discrete transmitted events per second |
| Bit rate / throughput | Actual bits delivered per second (bps) |
| Bandwidth | Range of frequencies a channel/signal occupies, in Hz — NOT bps |
| Nyquist rate | Noiseless maximum symbol rate = 2 x bandwidth |
| Adaptive modulation | Dynamically lowering QAM order as SNR degrades, to maintain reliability |

This chapter gave frequency, amplitude, and phase their precise definitions, showed how combining amplitude and phase (QAM) multiplies the information carried per symbol, and — critically — nailed down bandwidth as a frequency-domain quantity in Hz, distinct from throughput in bps. But every calculation so far has been "noiseless" — the Nyquist rate assumed a perfect channel. Real channels are never perfect: they attenuate signals, absorb interference from other sources, and carry an inescapable noise floor. Chapter 17, **Noise, Attenuation, Interference, and SNR**, quantifies exactly how much a real channel degrades a signal — the missing ingredient Chapter 18 needs to complete Nyquist's noiseless formula into Shannon's full, real-world channel capacity theorem.

