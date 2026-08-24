# Chapter 15: Analog vs. Digital Signals, and Modulation

*"A wire can only really do one thing: carry a voltage that changes over time. Everything else — text, video, voice, a stock trade — is a story we tell about that voltage."*

---

## Table of Contents

1. [The Problem This Chapter Solves](#1-the-problem-this-chapter-solves)
2. [Analog Signals — Continuous, Like the World Itself](#2-analog-signals--continuous-like-the-world-itself)
3. [Digital Signals — Discrete, By Agreement](#3-digital-signals--discrete-by-agreement)
4. [Why Digital Resists Noise (and Analog Doesn't)](#4-why-digital-resists-noise-and-analog-doesnt)
5. [The New Problem: Not Every Channel Wants a Square Wave](#5-the-new-problem-not-every-channel-wants-a-square-wave)
6. [The Carrier Wave — Borrowing a Channel's Favorite Frequency](#6-the-carrier-wave--borrowing-a-channels-favorite-frequency)
7. [Modulation, Defined](#7-modulation-defined)
8. [ASK — Amplitude Shift Keying](#8-ask--amplitude-shift-keying)
9. [FSK — Frequency Shift Keying](#9-fsk--frequency-shift-keying)
10. [PSK — Phase Shift Keying](#10-psk--phase-shift-keying)
11. [Comparing ASK, FSK, and PSK](#11-comparing-ask-fsk-and-psk)
12. [Real Systems That Use These](#12-real-systems-that-use-these)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Hands-On Experiment](#14-hands-on-experiment)
15. [Code: Simulating ASK/FSK/PSK in Go](#15-code-simulating-askfskpsk-in-go)
16. [What This Chapter Simplified](#16-what-this-chapter-simplified)
17. [Interview Questions & Model Answers](#interview-questions--model-answers)
18. [Exercises](#exercises)
19. [Summary](#summary)

---

## 1. The Problem This Chapter Solves

Chapter 14 ended with bits becoming voltage pulses on a wire: high voltage for a 1, low voltage for a 0, held steady for a fixed slice of time, then abruptly switching for the next bit. That square, on/off waveform is called a **digital baseband signal**, and it works beautifully — over a short, well-shielded copper wire.

But two enormous categories of real-world channels flatly refuse to carry that waveform:

1. **Channels that are physically bandlimited to a narrow range of frequencies.** The classic example: the telephone network was built, over a century, to carry human voice — frequencies roughly 300 Hz to 3,400 Hz — and nothing else. Send a sharp on/off digital pulse train down that line and the phone network's filters mangle it into mush, because a sharp square pulse actually requires a *very wide* range of frequencies to represent faithfully (a fact this chapter will explain and Chapter 16 will make precise).
2. **Channels with no wire at all.** Radio and Wi-Fi have to radiate energy through open air using an antenna. An antenna can only efficiently radiate a wave whose frequency matches its physical size (roughly, an antenna's length should be a meaningful fraction of the wave's *wavelength*). There is no way to "turn off" and "turn on" empty air the way you turn a wire's voltage up and down — you can only broadcast a wave and vary *properties* of that wave.

So here is the real question this chapter answers: **how do you send digital information (bits) through a channel that will only carry a specific, narrow band of continuous wave frequencies?** The answer, invented independently for telephone modems and radio decades before Wi-Fi existed, is **modulation** — and it's one of the most important ideas in all of Part 3.

---

## 2. Analog Signals — Continuous, Like the World Itself

**Intuitive version:** Turn a dimmer-switch light bulb slowly from off to full brightness. At every instant, as your hand moves, the bulb is at *some* brightness — there's no jump between "off" and "the next level." The brightness is a continuous quantity: it can be any value in a range, and it changes smoothly.

**Engineering definition:** An **analog signal** is one whose value (voltage, current, sound pressure, light intensity) can take on *any* value within a continuous range, and typically varies smoothly and continuously over time. Sound is analog: air pressure varies continuously as a sound wave passes. The original telephone system was entirely analog: your voice's air-pressure wave was converted directly into a voltage wave of the same shape, sent down a copper pair, and converted back into sound at the other end — no bits, ever, anywhere in that original design.

```
   Analog signal (e.g. voice, or the voltage
   from an old telephone handset's microphone):

   voltage
     |        .-.           .--.
     |       /   \         /    \      .
     |      /     \       /      \    / \
     |-----/-------\-----/--------\--/---\------> time
     |    /         \   /          \/     \
     |   /           `-'                   `-.
     |  /
```

Every value along that smooth curve is meaningful and distinguishable, in principle. That's both analog's strength (it can represent infinitely fine detail) and, as the next section shows, its fatal weakness for reliable long-distance communication.

---

## 3. Digital Signals — Discrete, By Agreement

**Intuitive version:** Now replace the dimmer switch with a normal light switch. It has exactly two positions: on or off. There is no "in-between" that means anything — the switch is always, unambiguously, in one of two states.

**Engineering definition:** A **digital signal** takes on only a finite number of discrete, pre-agreed values (most commonly two, as built up in Chapter 14). It doesn't matter, for interpretation purposes, whether the actual voltage is 3.29V or 3.31V — both are read as "1," because the receiver only checks which side of a threshold the voltage falls on (Chapter 14, Section 3).

```
   Digital signal (the NRZ waveform from Chapter 14):

   voltage
   3.3V |___     ___         ___     ___
        |   |   |   |       |   |   |   |
   0.0V |   |___|   |_______|   |___|   |___
        +---+---+---+---+---+---+---+---+---> time
             1   0   0   1   1   0   1   0
```

Only two vertical positions ever appear on this graph, by design, no matter what the underlying voltage actually measured at any given instant.

---

## 4. Why Digital Resists Noise (and Analog Doesn't)

This is the single most important comparison in the chapter, and it explains why the entire modern communications world — voice calls, video, radio, even AM/FM broadcast increasingly being replaced by digital equivalents — has moved to digital representation wherever possible.

**The core mechanism: regeneration.** A digital receiver doesn't need to reproduce the exact voltage that was sent — it only needs to correctly decide "was that closer to a 1 or a 0?" Any noise that's smaller than the gap between the two threshold zones (Chapter 14's noise margin) gets *completely erased* the instant the signal is read and re-emitted as a fresh, clean 1 or 0. This is why digital repeaters (used throughout the telephone backbone, undersea cables, and long fiber runs) can regenerate a signal perfectly, indefinitely, over any distance, as long as noise never gets large enough to cross a threshold.

An analog signal has no such luxury. Every detail of the waveform's shape *is* the information — there's no "close enough, round it to the nearest allowed value" step, because every value is allowed. An analog repeater can only *amplify* the signal (make it bigger), and amplification makes the noise bigger right along with the real signal. Every time you amplify an analog signal to compensate for cable loss, you also amplify all the noise that's already accumulated. Analog long-distance telephone calls, before all-digital trunk lines existed, would get progressively noisier with every amplifier hop — the crackle of an old long-distance call was this effect, audible.

```
  DIGITAL REGENERATION                    ANALOG AMPLIFICATION

  original:  1 0 1 1 0                    original: ~~/\~~
                 |                                        |
  + noise:  1.1V 0.3V 0.9V 1.2V -0.1V     + noise:  ~~/\~~ + wiggly noise
                 |  (still cleanly                        |
                 |   above/below threshold)                |
  regenerate: 1 0 1 1 0   <- IDENTICAL     amplify:  BIGGER ~~/\~~ + BIGGER
              to original                            wiggly noise (both scaled up)
```

**Worked numeric intuition:** Suppose a digital signal has a noise margin of ±0.6V (Chapter 14's example: "1" above 2.0V, "0" below 0.8V, out of a 0-3.3V swing). As long as accumulated electrical noise on the line stays under roughly 0.6V, every single regeneration recovers the *exact* original bit sequence — a thousand hops later, the signal is bit-for-bit identical to what left the sender. An analog signal has no such recovery: if noise adds even 0.1V of unwanted wiggle at every hop, after enough hops the accumulated noise can become comparable to the signal itself, degrading it permanently and unrecoverably.

This is why your voice, streamed as a digital phone call today, sounds crystal clear across a call routed through six countries and a dozen switching hops, in a way that would have been impossible on the pure-analog long-distance network of the 1970s.

---

## 5. The New Problem: Not Every Channel Wants a Square Wave

Section 1 named the problem in general terms; here it is precisely, using a fact you'll fully derive with Fourier reasoning in Chapter 16, but can accept intuitively for now:

**A perfectly sharp square wave (the digital signal from Section 3) is mathematically equivalent to an infinite sum of sine waves at many different frequencies, layered on top of each other.** The sharper and more sudden the voltage transitions, the *wider* the range of frequencies required to represent that shape faithfully. A truly instantaneous voltage jump would, in theory, require infinite bandwidth to reproduce exactly.

Real channels don't have infinite bandwidth — they have a limited range of frequencies they can carry (Chapter 16 defines this precisely as **bandwidth**). A telephone line physically cannot carry frequencies outside roughly 300-3,400 Hz — its wires, filters, and switching equipment were built and tuned specifically for that voice band. Send a digital square wave with sharp transitions down that line unmodified, and the phone system's filters strip away the high-frequency components that gave the signal its sharp edges, turning crisp square pulses into blurred, rounded, overlapping mush that a receiver cannot reliably tell apart.

Radio has the opposite but related problem: there's no wire at all, only open air, and an antenna can't radiate a DC voltage or a slow on/off toggle efficiently — physics (specifically, how efficiently an antenna converts electrical energy to radiated electromagnetic waves) demands a wave already oscillating at a frequency related to the antenna's size.

Both problems have the same solution.

---

## 6. The Carrier Wave — Borrowing a Channel's Favorite Frequency

**The idea:** instead of sending the digital bits directly, generate a smooth, continuous sine wave at a frequency the channel *is* happy to carry — this is called the **carrier wave** — and then vary some property of that wave in a pattern that encodes the bits. The channel sees only a well-behaved sine wave centered on a frequency it was built for; the receiver, knowing what to look for, can recover the original bits from how that wave changes.

A general sine wave is fully described by three, and only three, adjustable properties:

```
  s(t) = A · sin(2π f t + φ)
          |        |     |
          |        |     +-- phase (φ), in radians or degrees
          |        +-- frequency (f), in Hz
          +-- amplitude (A)
```

- **A (amplitude):** how "big" the wave is — its peak voltage/power.
- **f (frequency):** how many complete cycles the wave makes per second.
- **φ (phase):** where in its cycle the wave is at a given instant, relative to a reference.

**This is the single most important fact in this entire chapter:** there are exactly three properties of a wave you can vary to encode information, and every modulation scheme that has ever existed is built by varying one, two, or all three of them. Chapter 16 goes deep on each property individually and shows how modern systems (QAM) vary two simultaneously. This chapter introduces the three simplest schemes, each varying exactly one property — historically, exactly the order in which they were invented and adopted.

---

## 7. Modulation, Defined

**Modulation** is the general technique of encoding a digital (or analog) signal onto a carrier wave by systematically varying one or more of the carrier's amplitude, frequency, or phase, in a pattern the receiver can detect and reverse (**demodulation**) to recover the original data.

```
  BITS  --[modulator]-->  CARRIER WAVE, VARIED  --[channel]-->  [demodulator]-->  BITS
   101                     (e.g. loud-quiet-loud                                    101
                            = amplitude changes
                            encoding 1-0-1)
```

The three schemes below are named with the suffix "**-shift keying**" — "keying" being an old telegraph-era term (Chapter 7) for switching a signal between discrete states, carried forward into radio and modem terminology.

---

## 8. ASK — Amplitude Shift Keying

**The scheme:** vary only the amplitude (A). The simplest version, **On-Off Keying (OOK)**, sends the carrier at full strength for a 1 and turns it off (zero amplitude) for a 0.

```
  Bits to send:      1     0     1     1     0

  ASK/OOK waveform (carrier present = 1, absent = 0):

  /\/\        /\/\  /\/\
 /    \      /    \/    \
/      \____/            \____
```

**Worked example:** infrared TV remote controls historically used a close cousin of ASK/OOK — turning an infrared LED fully on and off in bursts to encode button presses, with the receiver simply detecting "light present" vs. "light absent" over successive time windows.

**Why it's simple, and why that's also its weakness:** ASK is trivial to generate and detect — you just need to measure signal strength. But real channels have varying noise and interference (Chapter 17) that also changes signal *amplitude* unpredictably; a burst of electrical noise can look exactly like a "1" was sent when it wasn't, or can wash out a genuine "1" so it looks like silence. ASK is therefore the *least* noise-resistant of the three basic schemes — a direct consequence of the fact that amplitude is the wave property most easily corrupted by real-world noise (which itself is essentially random amplitude fluctuation).

---

## 9. FSK — Frequency Shift Keying

**The scheme:** keep amplitude constant, vary only the frequency (f). A 1 is sent as a burst of the carrier at one frequency; a 0 is sent as a burst at a different frequency.

```
  Bits to send:      1        0        1        1        0

  FSK waveform (higher frequency = 1, lower frequency = 0):

  /\/\/\/\   /\  /\   /\/\/\/\  /\/\/\/\   /\  /\
 /        \ /  \/  \ /        \/        \ /  \/  \
```

**Real worked example — the Bell 103 modem (1962), one of the first widely deployed modems:** it used FSK over ordinary phone lines at 300 bits per second, with two separate pairs of frequencies depending on which side of the call was "originating" versus "answering" (so both directions could transmit simultaneously without interfering — full duplex):

| Role | Frequency for "0" (space) | Frequency for "1" (mark) |
|---|---|---|
| Originate (caller) | 1,070 Hz | 1,270 Hz |
| Answer (callee) | 2,025 Hz | 2,225 Hz |

Notice both frequency pairs sit comfortably inside the phone line's ~300-3,400 Hz voice passband (Section 1) — this is exactly the constraint that forced the invention of modulation in the first place. Notice also that the two directions use *entirely different frequency ranges*, so the originating modem's signal and the answering modem's signal never overlap and can be filtered apart cleanly at each end.

**Why it's better than ASK, but not free:** because amplitude is untouched, FSK is much more robust to the amplitude-distorting noise that devastates ASK — a receiver just has to determine *which* frequency is present, not *how strong* it is. The cost: FSK requires a wider range of frequencies (Chapter 16 will call this a wider bandwidth) than an equivalent ASK signal, since it needs room for two (or more) distinct frequency bands rather than one.

---

## 10. PSK — Phase Shift Keying

**The scheme:** keep amplitude and frequency constant, vary only the phase (φ) — the position of the wave within its repeating cycle.

The simplest version, **Binary Phase Shift Keying (BPSK)**, uses just two phase states, 180° apart:

```
  Bits to send:      1         0         1

  BPSK waveform (0° phase = 1, 180° phase = 0 -- i.e. the wave "flips"):

   /\      \/      /\
  /  \     /\     /  \
 /    \   /  \   /    \
     0°  180°   0°   <- phase reference points
```

A more advanced version, **Quadrature Phase Shift Keying (QPSK)**, uses four phase states 90° apart, letting each single wave "symbol" carry 2 bits instead of 1 (since log₂(4) = 2):

```
  Constellation diagram for QPSK (4 phase states, each encoding 2 bits):

              90°
               |
        01     |     11
               |
   180° -------+------- 0°
               |
        00     |     10
               |
              270°
```

(A **constellation diagram** plots each possible signal state as a point — for phase-only schemes, all points sit on a circle of fixed radius, at different angles. Chapter 16 extends this to QAM, where points also vary in *distance* from the center, representing amplitude.)

**Why phase resists noise even better than frequency, in practice:** modern receivers can measure phase very precisely relative to a reference, and phase itself (unlike amplitude) isn't directly corrupted by the amplitude fluctuations that dominate most real-world channel noise. PSK schemes, especially QPSK and its higher-order relatives, became the workhorse of satellite communication, Wi-Fi, and cellular modulation precisely because they pack multiple bits per symbol (like QPSK's 2 bits) while remaining highly noise-resistant — a theme Chapter 16 develops fully with QAM, which combines phase *and* amplitude variation to pack even more bits per symbol.

---

## 11. Comparing ASK, FSK, and PSK

| Scheme | Varies | Bits/symbol (basic version) | Noise resistance | Bandwidth needed | Historical/real use |
|---|---|---|---|---|---|
| ASK (OOK) | Amplitude | 1 | Low (amplitude noise directly corrupts it) | Narrow | IR remotes, early low-speed telegraph-derived links, some fiber (on/off keying) |
| FSK | Frequency | 1 (2 for 4-FSK, etc.) | Medium | Wider (needs room for multiple frequencies) | Bell 103/202 modems, older caller ID, some IoT radios (LoRa uses a chirp variant) |
| PSK (BPSK/QPSK) | Phase | 1 (BPSK) or 2 (QPSK) | High | Narrow (comparable to ASK) | Satellite links, Wi-Fi (lower data rates), cellular, dial-up's later, faster modem standards |

The historical trend visible in this table — ASK first (simplest to build with 1960s electronics), then FSK (better noise resistance for the modest cost of more bandwidth), then PSK (best noise resistance, efficient bandwidth use, but requires more precise electronics to measure phase accurately) — is not a coincidence. It mirrors exactly the order these schemes became practical to *build* as electronics improved, and it sets up Chapter 16's punchline directly: once you can measure both amplitude and phase precisely (which modern chips can), why not vary *both at once* and pack even more bits into a single symbol? That combined scheme is QAM, and it's how essentially every high-speed modern system (cable modems, DSL, Wi-Fi, LTE, 5G) actually works today.

---

## 12. Real Systems That Use These

- **Dial-up modems** (fully explored numerically in Chapter 18): early ones (300 bps) used FSK exactly as described in Section 9; later, faster standards (up to 33.6 kbps) moved to increasingly sophisticated PSK and QAM-family schemes to pack more bits into the same ~3 kHz phone-line bandwidth.
- **Wi-Fi** (Volume 13): uses BPSK and QPSK at its most robust (lowest, longest-range) data rates, and steps up to 16-QAM, 64-QAM, 256-QAM, even 1024-QAM at its fastest, shortest-range rates — a direct, real-world instance of the "more bits per symbol vs. more noise sensitivity" tradeoff from Chapter 14, Section 4.
- **Digital TV and radio broadcasting** use variants of PSK and QAM to pack a full HD television channel into an over-the-air frequency slot originally designed for one analog channel.
- **Fiber optic long-haul links** (Chapter 22) increasingly use phase and amplitude modulation of light itself (not just simple on/off keying), for exactly the same "pack more bits per symbol" reason.
- **Pagers and some low-power IoT radios** (e.g., the POCSAG paging protocol, and LoRa's "chirp spread spectrum," a close relative of FSK) still use frequency-based modulation today, prized specifically for FSK's simplicity and robustness over ASK, even at the cost of the extra bandwidth Section 11's table notes — a good real-world reminder that "older" and "simpler" doesn't mean "obsolete" when the application (long battery life, long range, low data rate) plays to that scheme's actual strengths.
- **GPS signals** are BPSK-modulated: each satellite broadcasts a very precisely timed BPSK signal, and a GPS receiver's entire ability to compute your location depends on measuring tiny differences in when each satellite's BPSK phase transitions arrive — a striking example of Section 10's "phase resists noise, and can be measured very precisely" property being pushed to its practical extreme (differences of nanoseconds in arrival time translate to positioning accuracy of meters).

### Production Note: The Historical Modem Standard Progression

Every step in the following table is a real ITU-T standard, and every jump in bits-per-second came from exactly the modulation-scheme evolution this chapter describes — moving from ASK-like/FSK schemes toward increasingly sophisticated PSK, and eventually (Chapter 16) QAM:

| Standard | Year | Speed | Modulation approach |
|---|---|---|---|
| Bell 103 | 1962 | 300 bps | FSK (Section 9) |
| Bell 212A / V.22 | 1980 | 1,200 bps | PSK (4-phase, a QPSK relative) |
| V.32 | 1984 | 9,600 bps | QAM (Chapter 16) with trellis coding |
| V.34 | 1994 | 33,600 bps | Adaptive, highly-optimized QAM, near the Shannon limit (Chapter 18) |
| V.90 | 1998 | 56,000 bps (download) | Exploits the digital phone network directly (fully explained in Chapter 18) |

Notice that between 1962 and 1998, the *underlying phone line itself barely changed* — its bandwidth (Chapter 16) and typical noise characteristics (Chapter 17) were roughly the same 3 kHz voice channel throughout. Nearly the entire 186x jump in speed (300 bps to 56,000 bps) came from modulation and coding cleverness, not from a better wire — a concrete, decades-long demonstration of exactly the idea this chapter and the next two exist to teach.

---

## 13. Common Misconceptions

- **"Modulation is only for wireless/radio."** False — it's equally essential for wired systems whose channel is bandlimited, like the telephone line in Section 9's Bell 103 example, or DSL over copper telephone wire, or cable modems over coax.
- **"Digital signals are noise-*proof*."** No — they are noise-*resistant*, not immune. If accumulated noise is large enough to push a voltage across a decision threshold (Chapter 14, Section 3), a digital receiver will misread the bit, just as certainly as an analog signal would be distorted. Chapter 17 quantifies exactly how much noise it takes, and Chapters 19-20 cover what happens when it does happen (error detection and correction).
- **"FSK/PSK/ASK are old, obsolete technology."** The *simple, single-bit-per-symbol* versions described here mostly are, for high-speed use. But they are the conceptual building blocks of QAM (Chapter 16), which is very much alive in every Wi-Fi router, cable modem, and cellular tower operating today.
- **"A carrier wave 'contains' the message, like a bus carries passengers."** More precise: the carrier wave's *shape over time* — specifically, how its amplitude/frequency/phase changes — *is* the encoded message. There's no separate "cargo" riding inside the wave; the pattern of change is the entire message.

---

## 14. Hands-On Experiment

You can hear an actual FSK-modulated signal, even today, with tools most computers already have.

```bash
# On many Linux systems, "minimodem" can both generate and decode audio FSK,
# almost exactly like the Bell 103 example in Section 9.
# (install: sudo apt install minimodem   /   brew install minimodem   on macOS)

echo "Hello from Chapter 15" | minimodem --tx 300

# This will play an audible tone through your speakers that alternates
# between two frequencies -- literally FSK, exactly as described in this
# chapter, audible to the human ear because 300 baud audio FSK uses
# frequencies inside the range humans can hear (unlike real telephone-line
# FSK, which used the same general idea over the wire rather than the air).
```

If you don't have `minimodem` available, search for "dial-up modem sound" — the screeching, buzzing sound of a real dial-up connection establishing (from the 1990s-2000s) is a mix of exactly the ASK/FSK/PSK-family tones this chapter describes, layered together as the two modems negotiate the fastest modulation scheme they can both reliably decode over that specific phone line's noise and bandwidth (a negotiation process directly connected to Chapter 18's Shannon limit discussion).

---

## 15. Code: Simulating ASK/FSK/PSK in Go

This program generates numeric sample values for each modulation scheme, so you can see, concretely, how "vary the amplitude / frequency / phase" translates into an actual signal, sample by sample — a direct extension of Chapter 14's waveform code.

```go
package main

import (
	"fmt"
	"math"
)

const (
	sampleRate   = 20     // samples per bit period, for a readable printout
	carrierCycle = 2      // carrier cycles per bit period (for readability)
)

// ask generates samples for Amplitude Shift Keying: full amplitude for
// bit=1, zero amplitude for bit=0.
func ask(bit int, samples []float64) {
	amplitude := 0.0
	if bit == 1 {
		amplitude = 1.0
	}
	for i := range samples {
		theta := 2 * math.Pi * carrierCycle * float64(i) / float64(sampleRate)
		samples[i] = amplitude * math.Sin(theta)
	}
}

// fsk generates samples for Frequency Shift Keying: a higher-frequency
// carrier for bit=1, lower-frequency for bit=0.
func fsk(bit int, samples []float64) {
	cycles := float64(carrierCycle)
	if bit == 1 {
		cycles = carrierCycle * 2 // double frequency for "1"
	}
	for i := range samples {
		theta := 2 * math.Pi * cycles * float64(i) / float64(sampleRate)
		samples[i] = math.Sin(theta)
	}
}

// bpsk generates samples for Binary Phase Shift Keying: 0 phase for bit=1,
// 180-degree (pi radian) phase shift for bit=0.
func bpsk(bit int, samples []float64) {
	phase := 0.0
	if bit == 0 {
		phase = math.Pi
	}
	for i := range samples {
		theta := 2*math.Pi*carrierCycle*float64(i)/float64(sampleRate) + phase
		samples[i] = math.Sin(theta)
	}
}

func printWave(name string, samples []float64) {
	fmt.Printf("%-6s ", name)
	for _, v := range samples {
		switch {
		case v > 0.5:
			fmt.Print("#")
		case v < -0.5:
			fmt.Print(".")
		default:
			fmt.Print("-")
		}
	}
	fmt.Println()
}

func main() {
	bits := []int{1, 0, 1, 1, 0}
	fmt.Printf("Bits to send: %v\n\n", bits)

	for _, bit := range bits {
		samples := make([]float64, sampleRate)
		fmt.Printf("--- bit = %d ---\n", bit)

		ask(bit, samples)
		printWave("ASK", samples)

		fsk(bit, samples)
		printWave("FSK", samples)

		bpsk(bit, samples)
		printWave("BPSK", samples)

		fmt.Println()
	}
}
```

Running this prints, bit by bit, a crude ASCII rendering of each scheme's waveform — for ASK you'll see the wave vanish entirely on a "0" bit; for FSK you'll see the wave oscillate visibly faster on a "1"; for BPSK the wave's shape looks identical in "loudness" but is offset (flipped) between the two bit values — exactly the three properties (amplitude, frequency, phase) named in Section 6, each isolated and made visible.

---

## 16. What This Chapter Simplified

- Real modems and radios rarely use *pure* single-property ASK, FSK, or BPSK at their fastest data rates — Chapter 16 shows how modern systems combine amplitude and phase (QAM) to go much further.
- This chapter modulated one bit per symbol throughout, for clarity. Real FSK and PSK systems often use more than two frequencies or phases per symbol (4-FSK, QPSK, 8-PSK) to pack multiple bits into each transmitted symbol — introduced briefly with QPSK in Section 10, developed fully in Chapter 16.
- Real demodulation (recovering bits from a received wave) requires careful synchronization between sender and receiver clocks, and filtering out noise — both simplified away here and picked up properly in Chapter 17 (noise) and Chapter 19 (error detection, for when demodulation still gets it wrong).
- Fiber optic "modulation" (Chapter 22) is conceptually the same idea applied to light instead of radio/electrical waves, but the physical mechanism (modulating a laser) differs in ways this chapter didn't cover.

---

## Interview Questions & Model Answers

**Beginner: "What is modulation and why is it necessary?"**

Modulation is the technique of encoding digital information onto a continuous carrier wave by varying its amplitude, frequency, or phase. It's necessary because many real channels — a phone line limited to voice frequencies, or open air that can only efficiently carry radio waves via an antenna — cannot directly carry a raw digital on/off signal. Modulation lets that digital information "ride" on a wave the channel is actually built to carry.

**Intermediate: "Compare ASK, FSK, and PSK in terms of noise resistance and why that difference exists."**

ASK varies amplitude, and since most real-world channel noise itself manifests as random amplitude fluctuation, ASK is the most vulnerable of the three — noise can be mistaken for a signal, or wash out a genuine signal. FSK varies frequency instead, so a receiver only needs to determine which frequency is present, which is comparatively unaffected by amplitude noise, making FSK more robust at the cost of needing a wider range of frequencies (more bandwidth). PSK varies phase, which can be measured very precisely by modern receivers and is not directly corrupted by amplitude noise, giving it the best noise resistance of the three while still using bandwidth efficiently — which is why PSK and its descendants (like QPSK, and eventually QAM) dominate modern high-speed systems.

**Advanced: "Why couldn't dial-up modems just send a raw digital signal down the phone line, and how does the choice of modulation scheme connect to the maximum speed a modem could achieve?"**

The phone network's copper wiring, filters, and switching equipment were engineered specifically for the human voice band, roughly 300-3,400 Hz — a sharp digital square wave requires far more bandwidth than that to preserve its shape, so the phone network's filters would destroy it. Modems instead used modulation to encode bits as variations of a carrier wave that fit inside that narrow voice band — early modems used FSK (Bell 103, two frequency pairs for full duplex), and as electronics improved, later, faster modem standards moved to PSK and eventually QAM-family schemes to pack more bits into each transmitted symbol without needing more bandwidth than the phone line could carry. This connects directly to Shannon's theorem (Chapter 18): given a fixed bandwidth (~3 kHz) and a fixed signal-to-noise ratio, there's a hard mathematical ceiling on how many bits per second can be reliably sent, no matter how clever the modulation scheme — which is exactly why dial-up modems converged on a ceiling near 33.6-56 kbps rather than continuing to improve indefinitely.

---

## Exercises

### Easy

1. In your own words, explain why a telephone line can't just carry a raw digital square-wave signal, referencing the specific frequency range a phone line supports.
2. Draw (on paper) a simple 6-bit sequence (e.g., 101100) as an ASK waveform, following the style of Section 8's diagram.
3. Name the three properties of a sine wave that can be varied to encode information, and give the modulation scheme name for each when varied alone.

### Medium

4. Explain, using Section 9's Bell 103 frequency table, why a single phone line could carry data in both directions (full duplex) at the same time without the two directions interfering with each other.
5. Draw a QPSK constellation diagram (following Section 10) and label which 2-bit value you'd assign to each of the 4 phase positions, if you were designing the scheme yourself. Explain your choice.
6. Modify the Go code in Section 15 to implement 4-FSK (four distinct frequencies, each representing a 2-bit value: 00, 01, 10, 11) instead of binary FSK.

### Hard

7. Explain why phase noise (small, random errors in a receiver's measurement of a wave's phase) becomes a *more* serious problem as you move from BPSK (2 phase states) to QPSK (4 states) to 8-PSK (8 states), even though the underlying phase-measurement noise itself hasn't changed. (Hint: think about how close together the valid phase states are as you add more of them within the same 360°.)
8. A engineering team wants to send data over a channel that is extremely bandwidth-limited (very narrow range of usable frequencies) but has an excellent signal-to-noise ratio (very little noise). Between ASK, FSK, and PSK, which would you recommend as the better starting point, and why? (You may reference Chapter 16's bandwidth concept even though it's introduced formally in the next chapter.)

---

## Summary

| Term | Meaning |
|---|---|
| Analog signal | Continuously variable value; every point on its range is meaningful |
| Digital signal | Discrete, pre-agreed set of values (commonly two); interpreted by threshold |
| Regeneration | A digital receiver's ability to recover an exact original bit despite noise, by re-emitting a clean signal after thresholding |
| Baseband signal | The raw digital on/off signal, before modulation |
| Carrier wave | A continuous sine wave at a frequency the channel can carry, used as the "vehicle" for modulation |
| Modulation | Encoding information by varying a carrier's amplitude, frequency, and/or phase |
| ASK | Amplitude Shift Keying — varies amplitude; simplest, least noise-resistant |
| FSK | Frequency Shift Keying — varies frequency; more noise-resistant, needs more bandwidth |
| PSK | Phase Shift Keying — varies phase; highly noise-resistant, bandwidth-efficient |
| Constellation diagram | A plot of each possible modulated signal state, used to visualize and design modulation schemes |

Chapter 15 established the three knobs — amplitude, frequency, phase — and showed the simplest scheme for varying each one alone. Chapter 16, **Frequency, Amplitude, and Phase, and Bandwidth**, goes deep on each of those three properties individually, shows what happens when you vary two of them *at once* (QAM), and finally gives bandwidth — used loosely so far — its precise, technical definition in Hz, carefully distinguished from the everyday (and technically wrong) use of "bandwidth" to mean internet speed in bits per second.

