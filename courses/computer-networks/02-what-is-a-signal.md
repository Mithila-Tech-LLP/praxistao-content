# Chapter 02: What Is a Signal?

> *"An idea cannot travel down a wire. Only physics can. Every act of communication, no matter how abstract it feels, ends the same way: as a change in some physical quantity that a receiver can detect."*

---

## Table of Contents

1. [The Gap Chapter 1 Left Open](#1-the-gap-chapter-1-left-open)
2. [What Is a Signal?](#2-what-is-a-signal)
3. [Why All Communication Is Eventually Physics](#3-why-all-communication-is-eventually-physics)
4. [Encoding: Agreeing in Advance What a Signal Means](#4-encoding-agreeing-in-advance-what-a-signal-means)
5. [A Worked Example: Turning Text Into a Signal](#5-a-worked-example-turning-text-into-a-signal)
6. [A First Glimpse of Analog vs. Digital](#6-a-first-glimpse-of-analog-vs-digital)
7. [The Electromagnetic Spectrum: Where Every Signal Type Actually Lives](#7-the-electromagnetic-spectrum-where-every-signal-type-actually-lives)
8. [Production Notes: Real Signal Levels and Units You'll Actually See](#8-production-notes-real-signal-levels-and-units-youll-actually-see)
9. [Hands-On Experiment: Send a Message With Light](#9-hands-on-experiment-send-a-message-with-light)
10. [Common Misconceptions](#10-common-misconceptions)
11. [Connections Backward and Forward](#11-connections-backward-and-forward)
12. [Interview Questions & Model Answers](#12-interview-questions--model-answers)
13. [Exercises](#13-exercises)
14. [Summary](#14-summary)

---

## 1. The Gap Chapter 1 Left Open

Chapter 1 built a five-step picture of communication: a message gets **encoded** into something, that something crosses a **channel**, and the receiver **decodes** it back into a message. It deliberately left one word undefined: *signal*. What, physically, is the "something" that actually leaves computer A, travels through a wire or through the air, and arrives at computer B?

This is not a rhetorical question with a philosophical answer. It has an extremely concrete, physical answer, and understanding it precisely is what separates "I've heard of networking" from "I understand what a network actually does." So let's ask it as plainly as possible: **when your phone sends a message, what physically leaves the phone?**

Not "data." Not "bits," exactly — not yet. Something even more basic has to leave first, and that something is a **signal**.

---

## 2. What Is a Signal?

### Intuitive explanation

A signal is any physical quantity that **changes over time or space in a way that carries information**, because sender and receiver have agreed in advance on what those changes mean. A drumbeat is a signal (the drummer varies loudness and timing; the listener has learned what patterns mean). A red traffic light is a signal (a wavelength of light; drivers have learned that this wavelength means "stop"). A raised hand is a signal (a change in the position of an arm; onlookers have learned this typically means "wait" or "I have a question").

Note what all three examples share: in each case, something *physical* — sound pressure, light wavelength, arm position — is doing the actual work of crossing the distance between sender and receiver. The "meaning" (stop, wait, this pattern means the letter A) is not physically present in the wire, the air, or the light itself. Meaning is something *both parties agree to attach* to a physical pattern. This is exactly the shared-code idea from Chapter 1, Section 5, applied to something concrete and physical rather than an abstract example.

### Engineering terms

A signal, in electronics and networking, is almost always one of these three physical carriers:

- **Electrical signal**: a voltage (or current) on a wire that varies over time. This is what travels down a copper Ethernet cable or a USB cable.
- **Optical signal**: pulses of light, typically infrared, sent down a glass fiber (or, less commonly, through open air with a laser). This is what travels through the fiber-optic backbone of the Internet.
- **Radio signal**: an electromagnetic wave broadcast through open space at a particular frequency. This is what travels between your phone and a cell tower, or between your laptop and a Wi-Fi router.

In every one of these three cases, the signal is **a physical quantity that varies over time** — voltage rising and falling, light turning on and off (or shifting between wavelengths), a radio wave's amplitude or frequency shifting. The variation itself is the alphabet. What the variation *means* is the code, agreed in advance.

### Deep technical view

Physically, all three signal types trace back to the same underlying phenomenon: electromagnetic fields. An electrical signal on a copper wire is a changing electric field driven by free electrons moving in the conductor. Light in a fiber is an electromagnetic wave at a much higher frequency (hundreds of terahertz) confined inside a glass core by total internal reflection (the mechanism behind this is covered fully in Chapter 22). A radio wave is the same phenomenon — an electromagnetic wave — at a lower frequency (megahertz to low gigahertz for most networking uses), radiating freely through space rather than being guided by a wire or fiber.

This is worth sitting with for a moment: **an Ethernet cable, a fiber-optic strand, and a Wi-Fi radio are three different physical implementations of the exact same underlying idea** — encode information as variation in an electromagnetic field, and let that field propagate from sender to receiver. Chapters 21–23 cover the engineering trade-offs (speed, distance, cost, interference resistance) between these three implementations in detail. For now, the important idea is that they're all *answering the same question* — "how do I turn a symbol into something physical that moves?" — with different physical media.

---

## 3. Why All Communication Is Eventually Physics

Here is a claim worth stating baldly, because it resolves a lot of confusion for beginners: **there is no such thing as sending "data" without sending a signal.** Every layer of software abstraction you will ever meet — a file, a web page, an encrypted message, a video call — is, at the moment it actually leaves one machine and enters a wire, a fiber, or the air, nothing but a physical signal: a voltage level, a pulse of light, a radio wave. Software never "skips" this step. It cannot. There is no wormhole for bits — only wires, glass, and open air, and physics governing what can move through each of them and how fast.

This single fact explains a huge amount of what the rest of this course will cover:

- **Why there's a speed limit.** No signal can travel faster than light in that medium — about 300,000 km/s in a vacuum, and roughly 200,000 km/s in typical optical fiber (light travels slower in a denser medium like glass). This hard physical limit is why a request from Mumbai to a server in Virginia cannot possibly complete in under about 80 milliseconds round-trip, no matter how good the software is — the light itself needs that long just to make the round trip at roughly half the speed of light in vacuum, ignoring routing and processing delays entirely. (Volume 19 works through real global latency numbers built on exactly this fact.)
- **Why signals degrade.** Every physical medium resists the signal somewhat — copper has electrical resistance, glass fiber has slight absorption, air scatters and absorbs radio waves. A signal that starts strong arrives weaker. This is called **attenuation**, and Chapter 17 covers it, along with the noise problem from Chapter 1, in full technical depth.
- **Why there's a hard ceiling on speed for a given medium.** As previewed in Chapter 1, Shannon's channel capacity theorem (Chapter 18) shows that for any real physical channel, there's a maximum rate at which bits can be reliably sent through it — a physics fact, not an engineering shortfall.

### Attenuation, made concrete with real numbers

It's worth putting real figures against "signals degrade," because the differences between media are large enough to shape entire engineering decisions later in this course. Modern optical fiber, at the 1550 nanometer wavelength commonly used for long-haul links, loses roughly **0.2 decibels of signal strength per kilometer** — remarkably little, which is exactly why a single unbroken fiber run can cross 60-80 km before needing an amplifier. A typical copper twisted-pair Ethernet cable, by contrast, loses signal strength far faster, especially at high frequencies, which is why Ethernet over copper (Chapter 21) is standardized with a strict 100-meter maximum segment length — beyond that, the signal has degraded enough that a receiver can no longer reliably distinguish the intended voltage levels from noise. Radio signals attenuate fastest of all in open space, following an inverse-square relationship with distance, which is why Wi-Fi range is measured in tens of meters while a fiber run is measured in tens of kilometers. None of these numbers need to be memorized yet — Chapters 17 and 21-23 build them up properly — but seeing them side by side here should make "signals degrade" feel like an engineering constraint with real teeth, not an abstract warning.

None of this should feel intimidating yet — it's simply the honest reason this course spends an entire volume (Volume 3, Chapters 14–23) on "boring" physical details before ever mentioning a router or an IP address. Every protocol later in this course inherits these physical constraints; understanding them first means every later chapter's design choices will make sense as *responses to physics*, not as arbitrary rules to memorize.

---

## 4. Encoding: Agreeing in Advance What a Signal Means

### The problem, stated precisely

Suppose two computers are connected by a plain copper wire, and one wants to send the other the letter "A". The sending computer can vary the voltage on that wire however it likes — say, between 0 volts and 5 volts. But a voltage level, by itself, means nothing. **0 volts is not "off" and 5 volts is not "on" until both machines agree that's what those levels represent** — and further agree on a scheme for representing something as rich as "the letter A" using nothing but a sequence of these two voltage levels over time.

This is encoding, applied to the concrete case of an electrical signal, and it happens in (at least) two layers that are worth separating clearly, because later chapters will build on exactly this distinction:

1. **Symbol-to-bits encoding**: agreeing how a meaningful symbol (a letter, a color, a musical note) maps to a sequence of bits. The most famous early example is **ASCII** (American Standard Code for Information Interchange), a table agreed upon in the 1960s that maps each English letter, digit, and punctuation mark to a 7-bit number. The letter "A" is, by this agreement, the number 65 — in binary, `01000001`.
2. **Bits-to-physical-signal encoding**: agreeing how each bit (0 or 1) maps to an actual physical state on the channel — e.g., 0 volts represents a `0` bit, and 5 volts represents a `1` bit. This is often called **line coding**, and real systems use schemes considerably cleverer than "low voltage = 0, high voltage = 1" (Chapter 15 covers why, and what schemes like Manchester encoding and NRZ actually do). For this chapter, the simple version is enough to build the right mental model.

### A small but important warning

Notice that ASCII's mapping ("A" → 65 → `01000001`) is itself just another instance of the shared-code problem from Chapter 1, Section 5: it only works because essentially every computer built since the 1960s has agreed to use it (or its modern superset, Unicode, for characters beyond plain English). There is nothing "natural" about the number 65 meaning "A" — it's a historical agreement, exactly as arbitrary and exactly as essential as "two torches means enemy sighted." If a system used a different table (and some historically did — EBCDIC, used on old IBM mainframes, assigns completely different numbers), the identical bit pattern would decode to a different letter, another instance of the silent-misinterpretation failure from Chapter 1.

---

## 5. A Worked Example: Turning Text Into a Signal

Let's trace one letter, "A", all the way from a keystroke to a (simplified, illustrative) physical signal, using Go to make each step explicit and inspectable.

```go
package main

import "fmt"

func main() {
	letter := 'A'

	// Step 1: symbol -> bits, using ASCII (a shared code agreed decades ago)
	asciiValue := int(letter)
	fmt.Printf("Letter:        %c\n", letter)
	fmt.Printf("ASCII value:   %d\n", asciiValue)
	fmt.Printf("As 8 bits:     %08b\n", asciiValue)

	// Step 2: bits -> a simplified physical signal
	// Convention (agreed in advance, exactly like the torches in Chapter 1):
	//   bit 1 -> HIGH voltage (we'll print 5V)
	//   bit 0 -> LOW voltage  (we'll print 0V)
	fmt.Println("\nSimplified electrical signal over time:")
	bits := fmt.Sprintf("%08b", asciiValue)
	for i, bit := range bits {
		voltage := "0V (LOW) "
		if bit == '1' {
			voltage = "5V (HIGH)"
		}
		fmt.Printf("  time slot %d: bit=%c -> %s\n", i, bit, voltage)
	}
}
```

Running it prints:

```
Letter:        A
ASCII value:   65
As 8 bits:     01000001

Simplified electrical signal over time:
  time slot 0: bit=0 -> 0V (LOW)
  time slot 1: bit=1 -> 5V (HIGH)
  time slot 2: bit=0 -> 0V (LOW)
  time slot 3: bit=0 -> 0V (LOW)
  time slot 4: bit=0 -> 0V (LOW)
  time slot 5: bit=0 -> 0V (LOW)
  time slot 6: bit=0 -> 0V (LOW)
  time slot 7: bit=1 -> 5V (HIGH)
```

We can even draw what that voltage would look like over time, as a crude square wave (this is a genuine, if simplified, oscilloscope-style view of the signal):

```
5V |     __                                                        __
   |    |  |                                                      |  |
0V |____|  |______________________________________________________|  |____
   +----+----+----+----+----+----+----+----+----+
     t0   t1   t2   t3   t4   t5   t6   t7

     0    1    0    0    0    0    0    1     <- bits
```

This picture — a voltage that jumps between two discrete levels over discrete time slots — is a genuinely accurate (if simplified) model of what happens, electrically, on a real digital link when the letter "A" is sent. Real systems complicate this picture in ways Chapter 15 covers in detail (the exact voltage levels, how the receiver knows exactly when one time slot ends and the next begins — called **clock recovery** — and cleverer encodings that avoid long runs of the same voltage, which cause practical problems). But the core idea will not change for the rest of this course: **symbol → agreed bit pattern → agreed physical signal**, and the receiver runs this whole process in reverse.

### Decoding: running it backward

```go
package main

import (
	"fmt"
	"strconv"
)

func main() {
	// The receiver observed this sequence of voltages...
	observedSignal := []string{"LOW", "HIGH", "LOW", "LOW", "LOW", "LOW", "LOW", "HIGH"}

	// Step 1: signal -> bits (using the same agreed convention as the sender)
	bits := ""
	for _, level := range observedSignal {
		if level == "HIGH" {
			bits += "1"
		} else {
			bits += "0"
		}
	}
	fmt.Println("Recovered bits:", bits)

	// Step 2: bits -> symbol (using the same shared code, ASCII)
	value, _ := strconv.ParseInt(bits, 2, 64)
	fmt.Printf("ASCII value: %d\n", value)
	fmt.Printf("Decoded letter: %c\n", rune(value))
}
```

Output:

```
Recovered bits: 01000001
ASCII value: 65
Decoded letter: A
```

This is the complete loop from Chapter 1's diagram, made concrete: a symbol became bits, bits became a physical signal, the signal crossed a channel, the receiver observed the signal, turned it back into bits, and turned those bits back into the original symbol — all of it resting on prior agreement (ASCII, and the voltage convention) that was never itself transmitted alongside the message.

---

## 6. A First Glimpse of Analog vs. Digital

The signal in Section 5 only ever took one of two voltage levels — 0V or 5V, nothing in between. A signal restricted to a small number of discrete levels like this is called a **digital signal**. A signal that can vary smoothly and continuously across a range — like the actual sound pressure wave of a human voice, or the smoothly varying voltage from an old analog telephone microphone — is called an **analog signal**.

This distinction deserves an entire chapter of its own (Chapter 15), because it explains, among other things, why computer networks are built almost entirely on digital signals even though the physical world (sound, light intensity, temperature) is fundamentally analog. The short preview: digital signals are dramatically more resistant to the noise problem from Chapter 1, because a receiver only has to distinguish between a small number of possible levels (here, just two), giving it a lot of tolerance before noise causes a wrong decision — whereas an analog signal's exact value *is* the information, so any noise directly corrupts the message. Chapter 15 will show this with a worked comparison; for now, it's enough to know that when this course says "signal," it will almost always mean a digital signal, unless it explicitly says analog.

It's also worth previewing one more idea by name, without yet explaining it: real physical channels (radio especially) often can't directly carry a simple on/off digital pattern efficiently over distance, so engineers instead vary a smooth, continuous "carrier" wave's properties (its amplitude, frequency, or phase) to represent digital bits. This technique is called **modulation**, and it's the entire subject of Chapter 15 and 16. If you've ever seen "QAM," "AM," or "FM" mentioned anywhere, that's what those terms are about — different strategies for riding digital information on top of an analog carrier wave.

---

## 7. The Electromagnetic Spectrum: Where Every Signal Type Actually Lives

Section 2 claimed that electrical, optical, and radio signals are "the same underlying phenomenon at different frequencies." It's worth making that claim visually concrete, because it explains a genuinely common point of confusion: why can't Wi-Fi just use the same frequencies as a cell tower, or fiber-optic light just use the same wavelength as a light bulb?

```
LOWER FREQUENCY                                                    HIGHER FREQUENCY
(LONGER WAVELENGTH)                                              (SHORTER WAVELENGTH)
   |                                                                        |
   v                                                                        v
+--------+--------+--------+--------+--------+--------+--------+--------+
|  AC    | RADIO  |  Wi-Fi |  4G/5G | Micro- | Infra- | VISIBLE| Ultra- |
| power  | (AM/FM)| 2.4/5/ | cell.  | wave   | red    | LIGHT  | violet |
| (60Hz) | (kHz-  |  6 GHz | (600MHz| ovens  | (fiber |(what   |  and   |
|        |  MHz)  |        | -mmWave)| (2.4GHz)| optics)|eyes see)| beyond |
+--------+--------+--------+--------+--------+--------+--------+--------+
                                                            ^
                                                    fiber optic signals
                                                    (~200 THz, just below
                                                     visible red light)
```

Every signal type this chapter has named occupies its own slice of this one continuous spectrum, governed by the same physics (Maxwell's equations, if you've encountered them) throughout. This immediately explains a few things that otherwise look arbitrary:

- **Why Wi-Fi, Bluetooth, and microwave ovens interfere with each other.** They all operate in the same 2.4 GHz slice of spectrum (a coincidence of 1980s-90s regulatory history, not a technical necessity), so signals meant for one can genuinely disturb another — this is a real, common cause of home Wi-Fi slowdowns, covered practically in Chapter 87.
- **Why fiber optic light is usually infrared, not visible light.** Engineers chose wavelengths (commonly 850nm, 1310nm, and 1550nm — just past the red end of what human eyes can see) specifically because glass fiber has particularly low attenuation at those wavelengths, not because infrared is inherently "better" for carrying data.
- **Why regulators assign specific frequency bands to specific uses.** Since radio spectrum is a shared, finite physical resource (much like Chapter 3's shared medium, but shared across an entire country or planet rather than one building), governments license specific frequency ranges to specific uses (a particular cellular carrier, aviation radio, emergency services) precisely to prevent the kind of interference described above at a societal scale.

None of this requires the underlying physics to be memorized — the useful takeaway is that "electrical signal," "radio signal," and "optical signal" are not three unrelated technologies invented separately; they're three engineering choices about *which frequency range of the identical underlying phenomenon* to exploit for a given medium and use case.

---

## 8. Production Notes: Real Signal Levels and Units You'll Actually See

Section 5's example used a clean, simplified 0V/5V convention. It's worth grounding that in numbers you'll actually encounter if you ever look at real hardware specifications, a Wi-Fi settings screen, or a fiber transceiver's datasheet — so these concepts don't stay purely theoretical until Volume 3 arrives.

- **Logic voltage levels.** Real digital electronics historically used **TTL** (Transistor-Transistor Logic) levels, roughly 0V for a "0" and 5V for a "1", matching Section 5's example closely. Most modern chips use **CMOS** logic at lower voltages (3.3V or even 1.8V) to save power, but the principle — two clearly separated voltage bands, with a "forbidden zone" in between that should never be sampled — remains identical. Older serial ports (RS-232) actually used the *opposite* convention from what you might expect, and a wider voltage swing (±3V to ±15V), a real historical example of the shared-code problem from Chapter 1: connecting incompatible-voltage devices without a level-shifting adapter can fail to communicate, or in bad cases, damage a device.
- **Wi-Fi signal strength, in dBm.** Every phone's Wi-Fi settings, and every laptop's network menu, report signal strength in **dBm (decibel-milliwatts)** — a logarithmic unit of power. A strong Wi-Fi signal is typically around **-30 dBm** (very close to the access point); a usable but weak signal is closer to **-70 to -80 dBm**; below about **-90 dBm**, a connection typically becomes unreliable or drops. Because the scale is logarithmic, every -10 dBm represents a 10x drop in actual signal power — meaning -80 dBm is one thousand times weaker than -50 dBm, not just "a bit weaker." This is precisely why moving a Wi-Fi router even a few extra meters, or adding one more wall, can degrade a connection so much more than the distance alone would suggest.
- **Fiber optic power budgets.** Fiber transceiver datasheets specify a transmit power (often around 0 dBm, i.e., 1 milliwatt) and a receiver sensitivity threshold (often around -20 to -28 dBm for common standards). The difference between the two is the "power budget" — how much attenuation (Section 3's 0.2 dB/km figure, plus losses at every connector and splice) the link can tolerate before the signal becomes too weak to read reliably. Network engineers use exactly this arithmetic when planning how far a fiber run can go without needing an amplifier or repeater.

A short Go snippet makes the logarithmic dBm scale concrete:

```go
package main

import (
	"fmt"
	"math"
)

// dBmToMilliwatts converts a power level in dBm to milliwatts, showing
// just how nonlinear this common real-world unit actually is.
func dBmToMilliwatts(dbm float64) float64 {
	return math.Pow(10, dbm/10)
}

func main() {
	levels := []float64{0, -10, -30, -50, -70, -90}
	for _, dbm := range levels {
		fmt.Printf("%6.0f dBm -> %10.6f mW\n", dbm, dBmToMilliwatts(dbm))
	}
}
```

```
     0 dBm ->   1.000000 mW
   -10 dBm ->   0.100000 mW
   -30 dBm ->   0.001000 mW
   -50 dBm ->   0.000010 mW
   -70 dBm ->   0.000000 mW  (0.0000001 mW, rounds to zero at this precision)
   -90 dBm ->   0.000000 mW  (0.000000001 mW)
```

Notice how a "weak but usable" -70 dBm Wi-Fi signal is already carrying roughly ten million times less power than the -30 dBm "excellent" signal a phone reports right next to a router — a striking, real illustration of just how sensitive a modern radio receiver has to be to recover a usable signal at all, and why Chapter 17 (Noise, Attenuation, Interference, and SNR) treats these numbers as central engineering quantities rather than trivia.

---

## 9. Hands-On Experiment: Send a Message With Light

This experiment makes "signal," "encoding," and "shared code" completely physical, using nothing but a phone flashlight (or any light source you can turn on and off) and a partner in another room or with their back turned.

**Setup — agree on a code first (crucially, before starting, exactly like Chapter 1's torches):**

```
Short flash (under 1 second) = the bit 0
Long flash (about 2 seconds) = the bit 1
A 3-second pause             = "end of letter, decode what you have using ASCII"
```

**Steps:**

1. Pick a single short word, like "hi". Look up (or compute, using the Go program from Section 5) the 8-bit ASCII pattern for each letter: 'h' = 104 = `01101000`, 'i' = 105 = `01101001`.
2. Using your flashlight, physically transmit each bit as a short or long flash, pausing 3 seconds between letters.
3. Have your partner write down "short" or "long" for every flash they see, without knowing what word you're sending.
4. After you finish, have your partner convert their observed short/long sequence back into 0s and 1s, group them into 8-bit chunks, and look up each chunk's ASCII value.

**What almost always happens:** at least one bit gets misread — a flash held slightly too long looks "long" instead of "short," or a pause is misjudged. This is a direct, physical demonstration of the **noise** and **attenuation** concepts from Chapter 1 and Section 3 — real channels (including the "channel" of one person watching another's flashlight) introduce timing and interpretation errors, and even a single flipped bit changes which letter gets decoded. If you want to see this exactly: flipping the last bit of 'h' (`01101000` → `01101001`) silently turns it into 'i' — no error, no warning, just a wrong-but-valid letter. This is exactly the silent-misinterpretation failure mode from Chapter 1, Section 5, and it's the entire motivation for Chapters 19 and 20 (error detection and correction), which exist precisely because real signals get corrupted exactly like this.

---

## 10. Common Misconceptions

- **"A signal is the same thing as data."** A signal is the *physical* carrier (a voltage, a light pulse, a radio wave). Data is the symbolic content (bits, bytes) that the signal is encoding. The same data can be carried by very different signals (electrical on copper, optical on fiber, radio through the air) — this is exactly why your home network can seamlessly mix a Wi-Fi laptop with a wired desktop: different signals, same underlying bits, translated at the boundary (a Wi-Fi router does exactly this translation).
- **"Digital means perfect, analog means imperfect."** Both are physical signals subject to the same noise and attenuation. Digital signals are simply *more tolerant* of noise, because the receiver only needs to distinguish between a small number of discrete levels rather than recover an exact continuous value. Digital signals can and do get corrupted — Chapters 19 and 20 exist because of this.
- **"Encoding is about making data secret."** In this chapter's sense — and in most of networking — encoding just means "converting a message into a form suitable for transmission," with no implication of secrecy at all. ASCII encoding doesn't hide anything; anyone who knows the ASCII table can read it immediately. Secrecy (encryption) is a completely different, much later topic (Volume 12), layered on top of ordinary encoding, not a synonym for it.
- **"A wire carries bits."** Strictly, a wire carries voltage that varies over time. "Bits" is the human-friendly interpretation of that voltage pattern, agreed upon in advance. Nothing labeled "0" or "1" is physically present in the copper — only voltage, which both ends have agreed to interpret as bits.
- **"Radio, Wi-Fi, and light are fundamentally different kinds of physics."** As Section 7 shows, they are the identical physical phenomenon (an electromagnetic wave) at different frequencies. Engineering differences between them (range, ability to pass through walls, how they're generated and detected) come from frequency, not from being different categories of physics.

---

## 11. Connections Backward and Forward

This chapter took the single word "signal" from Chapter 1's diagram (Section 6 of that chapter) and gave it a full, physical, three-level treatment: intuitive (a drumbeat, a red light), engineering (electrical/optical/radio signals, line coding), and deep technical (electromagnetic fields, the speed-of-light limit, ASCII as a real worked encoding, and real dBm power figures). It also quietly set up vocabulary — analog vs. digital, modulation, attenuation, noise, the electromagnetic spectrum — that Volume 3 (Chapters 14–23) will each get a dedicated, rigorous chapter of their own.

What this chapter has *not* yet done is explain what happens when there isn't just one channel between exactly two computers, but a whole building, city, or planet full of computers that all need to reach each other — potentially sharing the same physical wires and radio spectrum. That is the problem Chapter 3 opens with, and it's the moment this course starts using the word "network" precisely for the first time.

---

## 12. Interview Questions & Model Answers

**Q1 (Beginner): What is a signal, and how does it differ from "data" or "information"?**

*Model answer:* A signal is a physical quantity — voltage, light, or a radio wave — that varies over time or space in a way that carries a symbol, given a shared code agreed on in advance. Data and information are the symbolic/abstract content being carried; the signal is the physical carrier that data is encoded into so it can actually travel across a channel. The same data can be represented by different physical signals (electrical, optical, radio); the same signal type can carry completely different data depending on the encoding scheme in use.

**Q2 (Beginner): Name the three main physical carrier types used in networking and one real-world example of each.**

*Model answer:* Electrical signals (voltage on copper wire, e.g., a wired Ethernet cable or USB cable), optical signals (light pulses in glass fiber, e.g., long-haul internet backbone links), and radio signals (electromagnetic waves through open space, e.g., Wi-Fi between a laptop and a router, or a phone talking to a cell tower).

**Q3 (Intermediate): Why can no signal, in any medium, travel instantly, and why does this matter for real networks?**

*Model answer:* All three signal types are ultimately electromagnetic phenomena, and electromagnetic waves are bounded by the speed of light in whatever medium they're traveling through — about 300,000 km/s in vacuum, roughly 200,000 km/s in typical optical fiber due to the fiber's refractive index. This creates an unavoidable propagation delay proportional to distance, which is why, for example, a round trip between geographically distant data centers has a hard minimum latency no software optimization can eliminate — this becomes directly relevant when reasoning about global service latency (Volume 19) and satellite links, which must cross tens of thousands of kilometers to geostationary orbit and back.

**Q4 (Intermediate): Explain, with an example, why digital signals are more resistant to noise than analog signals.**

*Model answer:* An analog signal's exact continuous value is itself the information being conveyed, so any noise added to it directly corrupts the message — the receiver cannot easily tell "signal" from "signal plus noise." A digital signal only needs to be classified into one of a small number of discrete levels (commonly two, "high" and "low"); as long as noise doesn't push the received value past the threshold that separates one level from another, the receiver can still recover the exact original bit. This tolerance for imprecision, up to a threshold, is why nearly all modern networking uses digital signaling even though most underlying physical phenomena (sound, light intensity) are naturally analog.

**Q5 (Advanced): What problem does line coding (the "bits-to-physical-signal" mapping) solve beyond the simple "0 volts = 0, 5 volts = 1" scheme used in this chapter's example, and why can't real systems just use the simple scheme?**

*Model answer:* The simple scheme in this chapter assumed the receiver already knows exactly when each bit's time slot starts and ends — but in reality, sender and receiver clocks are never perfectly synchronized, and a long run of identical bits (e.g., many consecutive 0s) produces a flat, unchanging voltage with no transitions the receiver can use to resynchronize its timing, a problem called clock drift / loss of clock recovery. Real line coding schemes (e.g., Manchester encoding, covered with NRZ and others in Chapter 15) guarantee frequent voltage transitions regardless of the actual bit pattern, so the receiver can continuously recover timing information from the signal itself, and also often provide DC balance (equal average time spent high and low) to avoid practical electrical issues in real transmission hardware.

**Q6 (Advanced): A phone reports Wi-Fi signal strength of -85 dBm in one room and -55 dBm in another. Explain, quantitatively, how much difference in actual received power this represents, and why this matters more than the raw number difference (30) suggests.**

*Model answer:* dBm is a logarithmic unit, where every 10 dBm represents a 10x change in actual power (in milliwatts). A 30 dBm difference therefore represents a 10^3 = 1,000-fold difference in received signal power, not merely "30 units weaker." This matters because a receiver's ability to correctly distinguish the intended signal from background noise (its signal-to-noise ratio, formalized in Chapter 17) depends on this actual power ratio, not the dBm number's face value — a link at -85 dBm may already be operating close to the noise floor and prone to errors or drops, even though the numeric difference from -55 dBm looks moderate at first glance.

---

## 13. Exercises

### Easy

1. Using the ASCII value convention from Section 4 ('A' = 65), compute the ASCII value and 8-bit binary representation of the lowercase letter 'z'. (You may run the Go program in Section 5 with the letter changed.)
2. In your own words, explain the difference between a signal, encoding, and a shared code, using an example not from this chapter.
3. Give two real-world (non-computer) examples of an analog signal and two of a digital signal.
4. Using Section 7's spectrum diagram, explain in one or two sentences why a microwave oven can disrupt a Wi-Fi connection but not a Bluetooth-free FM radio broadcast.

### Medium

1. Take the 8-bit pattern for the letter 'h' (`01101000`) from Section 5 and flip a single bit of your choosing. Compute the new ASCII value and letter it now represents. Was the result still a valid, decodable letter, or did it fail to look up? What does this tell you about the risk of a single-bit error?
2. Modify the decoding Go program in Section 5 so that it also handles a signal reading of "UNKNOWN" (representing a corrupted or ambiguous voltage reading that is neither clearly HIGH nor LOW) by printing an explicit error rather than silently guessing. Explain why this is a meaningfully better design than always guessing HIGH or LOW.
3. Explain, using the speed-of-light figures from Section 3, roughly how long the minimum one-way propagation delay would be for a signal traveling 12,000 km through optical fiber (ignoring all router/switch processing delays). Show your arithmetic.
4. Run the dBm-to-milliwatts Go program in Section 8 with two additional dBm values of your choosing, and describe in words how much the actual power changes between them.

### Hard

1. Design a line-coding scheme (a rule for mapping bits to voltage levels over time) that guarantees at least one voltage transition every two time slots, regardless of the actual bit pattern being sent — even for the all-zero bit pattern `00000000`. You do not need to reproduce a real standard; invent your own rule and demonstrate it on a sample byte, showing the resulting voltage sequence.
2. Research the difference between ASCII (7-bit, English-centric) and Unicode/UTF-8 (variable-length, supports virtually all human scripts). Explain, using the shared-code framing from this chapter and Chapter 1, why the transition from ASCII to Unicode required backward-compatible design decisions rather than a clean break, and why backward compatibility matters for a globally shared code.
3. The hands-on experiment in Section 9 revealed timing errors as a real source of bit corruption. Propose (in plain language, without needing exact protocols) a way the sender could add extra information to their transmission that would let the receiver *detect*, even if not automatically fix, that at least one bit was probably misread. You are informally deriving the idea behind Chapter 19 (Error Detection) — a rough, correct intuition is the goal, not a complete algorithm.

---

## 14. Summary

| Term | Meaning |
|---|---|
| Signal | A physical quantity (voltage, light, radio wave) that varies over time to carry a symbol, given a shared, pre-agreed code |
| Electrical signal | Voltage/current variation on a conductor (e.g., copper Ethernet cable) |
| Optical signal | Light pulses guided through glass fiber |
| Radio signal | An electromagnetic wave broadcast through open space |
| ASCII | A historical shared-code table mapping English characters to 7-bit numbers |
| Line coding | The scheme mapping bits to physical signal states, addressing practical issues like clock recovery |
| Analog signal | A signal that can take any value continuously across a range |
| Digital signal | A signal restricted to a small number of discrete levels, giving it strong noise tolerance |
| Modulation (preview) | Varying a continuous carrier wave's amplitude, frequency, or phase to represent digital information |
| Attenuation | The weakening of a signal as it travels through a medium (e.g., ~0.2 dB/km in modern fiber) |
| Electromagnetic spectrum | The full range of frequencies electrical, radio, and optical signals occupy — one underlying phenomenon, many engineering uses |
| dBm | A logarithmic unit of signal power used to describe real Wi-Fi strength and fiber transceiver budgets |
| Propagation delay | The unavoidable time a signal takes to physically traverse a distance, bounded by the speed of light in that medium |

This chapter answered what a signal physically is and how encoding turns a symbol into one. Chapter 3 asks the next necessary question: once you know how to send one signal down one channel between two computers, what changes when there isn't just one pair of computers, but many computers that all need to reach each other?
