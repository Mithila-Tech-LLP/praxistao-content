# Chapter 19: Error Detection — Parity, Checksums, and CRC

> *Every wire lies sometimes. Chapter 17 already proved that noise, attenuation, and interference are physically unavoidable, and Chapter 18 showed that even a channel running safely under Shannon's limit still has some non-zero chance of flipping a bit. The question this chapter answers isn't "how do we stop errors" — that's often impossible. It's: "when a bit flips in transit, how does the receiver even find out?"*

---

## Table of Contents

1. [The Big Question: Can We Trust the Wire?](#1-the-big-question-can-we-trust-the-wire)
2. [Why Errors Happen At All](#2-why-errors-happen-at-all)
3. [The Naive Fix: Send It Twice](#3-the-naive-fix-send-it-twice)
4. [Parity Bits — The Simplest Real Check](#4-parity-bits--the-simplest-real-check)
5. [Two-Dimensional Parity — A Hint of Correction](#5-two-dimensional-parity--a-hint-of-correction)
6. [Checksums — Summing Your Way to Confidence](#6-checksums--summing-your-way-to-confidence)
7. [CRC — Cyclic Redundancy Check](#7-crc--cyclic-redundancy-check)
8. [Comparison — What Each Technique Catches and Misses](#8-comparison--what-each-technique-catches-and-misses)
9. [Where These Are Actually Used](#9-where-these-are-actually-used)
10. [Hands-On: Compute a Checksum and a CRC Yourself](#10-hands-on-compute-a-checksum-and-a-crc-yourself)
11. [Common Misconceptions](#11-common-misconceptions)
12. [Interview Questions & Model Answers](#12-interview-questions--model-answers)
13. [Exercises](#13-exercises)
14. [Summary](#summary)

---

## 1. The Big Question: Can We Trust the Wire?

Chapters 14–18 built up the physical layer from the ground up: bits become voltages or light pulses (Ch. 14), those get modulated onto carriers (Ch. 15), a wave's frequency/amplitude/phase encode more bits per symbol (Ch. 16), noise and attenuation degrade the signal as it travels (Ch. 17), and Shannon's limit (Ch. 18) tells you the theoretical *maximum* rate you can push through a channel of given bandwidth and SNR.

Here is the detail that trips people up about Shannon's limit: it tells you the maximum rate at which you *can*, in principle, drive the error rate arbitrarily close to zero using a sufficiently clever code — not that a real, practical system running below that limit has *zero* errors. In practice, every real link has a nonzero **bit error rate (BER)**:

```
Typical BER by medium (order of magnitude, well-engineered systems):
  Fiber optic backbone              ~1 in 10^12 bits  (1e-12)
  Copper Gigabit Ethernet           ~1 in 10^10 bits  (1e-10)
  Wi-Fi, good conditions            ~1 in 10^6  bits  (1e-6)
  Wi-Fi, marginal signal            ~1 in 10^4  bits  (1e-4) or worse
  Old dial-up modem                 ~1 in 10^5  bits
```

That sounds tiny. It isn't, at scale. At 1 Gbps with a BER of 1e-10, you get roughly **one flipped bit every 10 seconds** — million of bits corrupted per day on a single busy link, across a network with billions of links. Multiply that across a data center with a million Ethernet ports and errors are a constant, background fact of networking life, not a rare edge case.

So the real question a network has to answer, on every single frame, is: **did this arrive exactly as sent?** That's error *detection* — the subject of this chapter. Notice this is explicitly a smaller problem than error *correction* (Chapter 20). Detection only has to say "something is wrong here" — it doesn't have to say what, or where, or fix it. That weaker goal is why detection can be done far more cheaply than correction, which is why almost every protocol you'll meet in this course (Ethernet, IP, TCP, UDP, Wi-Fi) detects errors on every single frame, but very few of them correct errors in the data path itself.

---

## 2. Why Errors Happen At All

Briefly, recalling Chapter 17 in concrete terms — these are the actual physical events that flip a 1 into a 0 or vice versa:

- **Thermal noise**: the random jitter of electrons in any conductor at any temperature above absolute zero (Johnson-Nyquist noise). This is not a design flaw — it's physics, and it never goes away.
- **Attenuation**: a signal weakens with distance (Ch. 17). If it weakens enough that its voltage/light level for a "1" falls close to the voltage level for a "0," the receiver's sampling circuit can misread it.
- **Crosstalk and interference**: a neighboring wire's electromagnetic field induces a spurious voltage on your wire (fully explained in Chapter 21), or an external source — a motor, a microwave oven, a fluorescent ballast, another radio transmitter — injects energy into your channel.
- **Impulse noise / bursts**: a lightning strike, a relay switching, a bad solder joint arcing — brief but intense, and instead of flipping one random bit, it usually **wipes out a contiguous run of several bits at once**. This single fact — that real-world errors cluster into *bursts*, not lone isolated flips — is the single most important design constraint behind everything in Section 7 of this chapter.
- **Cosmic rays and radioactive decay**: a genuinely strange but real cause — a high-energy particle striking a memory cell or transistor can flip a stored bit. This matters more for Chapter 20 (ECC memory) than for wires, but it's the same underlying problem: a 1 becomes a 0 with no warning.

Given any channel with finite SNR, there's always *some* nonzero probability that a given symbol is decoded as its neighbor. Engineering never eliminates that probability — it only decides how much redundancy is worth spending to detect (this chapter) or correct (next chapter) the errors it can't prevent.

---

## 3. The Naive Fix: Send It Twice

The first idea anyone has, unprompted, is: **send the data twice, and compare the two copies.** If they match, assume it's correct. If they don't, ask for it again.

This "works," in the sense that it will catch some errors. But it fails as an engineering solution for three concrete reasons:

1. **Cost.** It doubles the bandwidth used for *every single transmission*, even though the overwhelming majority of frames (999,999,999 out of 1,000,000,000 at a 1e-9 BER) have no errors at all. You're paying a 100% tax to guard against an event that happens 0.0000001% of the time.
2. **Correlated failures.** A burst of noise (Section 2) that corrupts one copy is quite likely, depending on timing, to corrupt the *same bits* in a copy sent immediately after — especially if both copies are subject to the same environmental interference. Duplication is not automatically independent.
3. **It still can't tell you which copy is right.** If the two copies disagree, all you know is "one of these is wrong" — not which one, and not where. You're forced to discard both and ask the sender to retransmit anyway — which is exactly what a proper detection code lets you do, at a small fraction of the bandwidth cost.

The real engineering insight, which every scheme in this chapter builds on, is: **you don't need to duplicate the whole message. You need to add a small, mathematically-constructed summary of the message that is very unlikely to match by coincidence if even a few bits changed.** That's the entire idea behind parity, checksums, and CRC — each one is a cheaper, more clever "fingerprint" of the data than a second copy of the data itself.

---

## 4. Parity Bits — The Simplest Real Check

**Intuitive picture:** Imagine you write a shopping list, and at the bottom you jot down the total number of items. If someone accidentally drops or adds an item to your list before you check out, the count at the bottom won't match what you count at the register — you'll *know* something changed, even though the note at the bottom doesn't say which item. That's a parity bit: a one-bit summary of "how many 1s were in this data," attached so the receiver can recompute the same summary and compare.

**Mechanism:** For a block of bits, count how many bits are `1`. Append one extra bit chosen so that the *total* count of 1-bits (data + parity bit) is even (**even parity**) or odd (**odd parity** — less common, but used in old RS-232 serial links). The sender computes it once; the receiver recomputes it on arrival and compares.

**Worked example (even parity):**

```
Data (7 bits):      1 0 1 1 0 0 1
Count of 1-bits:    1+0+1+1+0+0+1 = 4   (already even)
Parity bit chosen:  0   (adding 0 keeps the count at 4, even)

Transmitted byte:   1 0 1 1 0 0 1 0
                     └─────data────┘ └parity┘
```

If a single bit flips in transit — say bit 3 flips from `1` to `0` — the received data becomes `1 0 0 1 0 0 1` (three 1-bits: odd), but the parity bit is still `0` (which claimed "even"). The receiver recomputes parity, gets a mismatch, and correctly flags an error.

**What parity catches:** any **odd number** of bit flips (1, 3, 5, 7, ...) in the block. Each single flip toggles the 1-count's parity, so an odd number of flips always leaves the count's parity different from what was declared.

**What parity misses — and this is the important part:** any **even number** of bit flips. Two flips cancel each other's effect on the parity, and the check silently passes over real corruption. Worked example:

```
Original data:      1 0 1 1 0 0 1     (four 1-bits, even)
Flip bit 1 AND bit 3:
New data:            0 0 0 1 0 0 1    (two 1-bits — still even!)

Parity bit (unchanged): 0
Receiver recomputes: even parity again → check PASSES.

Two bits of real corruption, zero bits of detection.
```

This is not a rare edge case — Section 2 explained that real interference tends to corrupt *bursts* of adjacent bits, and a burst that flips exactly two, four, or six bits is common. A single parity bit over a whole block is a genuinely weak defense on its own, which is exactly why real systems either apply parity per byte (limiting the "blast radius" of an undetected even-flip event) or move to the stronger schemes in Sections 6–7.

**Real use today:** parity was standard in early RS-232 serial links and pre-ECC computer memory. Modern systems have mostly replaced bare parity with the checksum and CRC techniques below for data in transit, and with Hamming-derived SECDED codes (Chapter 20) for memory.

---

## 5. Two-Dimensional Parity — A Hint of Correction

A clever extension pushes single-bit parity a step further: arrange the data as a rectangular grid, and compute a parity bit for **every row** and **every column**.

```
Data arranged as a 4x4 grid, with a parity column and parity row added:

           col0 col1 col2 col3 | row-parity
row0:        1    0    1    1  |    1   (1+0+1+1=3 → odd → parity=1 makes total even)
row1:        0    1    1    0  |    0
row2:        1    1    0    1  |    1
row3:        0    0    1    0  |    1
           ---------------------------
col-parity:  0    0    1    0  |   (corner bit, parity-of-parities)
```

If exactly **one** bit flips anywhere in the grid — say row1, col2 (originally `1`, flipped to `0`) — then:

- Row1's parity check now fails (its row no longer matches its stored row-parity bit).
- Col2's parity check now fails (same reason, column-wise).

The receiver doesn't just know *that* an error occurred — the **intersection of the failing row and failing column points at the exact bit** that flipped. This is the first real hint of *error correction*, not just detection: with enough structured redundancy, you can localize an error precisely enough to flip it back. Chapter 20 formalizes this idea into Hamming codes, which do the same trick with far less redundancy overhead by choosing parity groups more cleverly than "every row and column."

Two-dimensional parity still shares parity's core weakness: two errors in the *same* row and *same* column (or certain other symmetric patterns) can cancel out and go undetected. It's a stepping stone, not a production-grade scheme — but the idea of "structured overlapping parity groups localize an error" is exactly the idea Hamming codes formalize.

---

## 6. Checksums — Summing Your Way to Confidence

Parity is cheap (1 bit) but weak (misses any even number of flips). The next idea: instead of just counting 1-bits, treat the data as a sequence of *numbers* and add them up. A sum is far more sensitive to which bits changed than a simple 1-bit-count.

**Mechanism — the Internet checksum (RFC 1071), used by IPv4, TCP, UDP, and ICMP:**

1. Split the data into 16-bit words.
2. Add all the words together using **one's complement arithmetic**: whenever an addition overflows past 16 bits, take that overflow "carry" bit and add it back into the bottom of the sum (this is called an **end-around carry**).
3. The checksum transmitted is the **one's complement (bit-flip) of that final sum**.
4. The receiver adds all the received words *plus* the received checksum. If nothing was corrupted, the result is **all 1s** (16 bits of `1`) — any other result means an error occurred.

**Worked example (using 8-bit words to keep the arithmetic small — real Internet checksums use 16-bit words, the math is identical):**

```
Word A:  11111111   (255)
Word B:  00000001   (1)

Step 1 — ordinary binary addition:
    11111111
  + 00000001
  -----------
   100000000     (9 bits: overflowed past 8 bits!)

Step 2 — end-around carry: take the overflow bit (the 9th bit, a 1)
         and add it back into the 8-bit result:
    00000000
  +        1
  -----------
    00000001

Step 3 — checksum = one's complement (flip every bit) of that sum:
    00000001  →  11111110      ← this is the transmitted checksum

Transmitted: word A (11111111), word B (00000001), checksum (11111110)
```

**Receiver's check:** add A + B + checksum together (same end-around-carry rule):

```
    11111111
  + 00000001
  -----------
   100000000  → end-around carry → 00000001

    00000001
  + 11111110      (the checksum)
  -----------
    11111111      ← all 1s. Check PASSES — no error detected.
```

That "sum of everything, including the checksum itself, equals all 1s" trick is the elegant part of one's-complement checksums, and it's exactly what a real router or NIC does in hardware on every IPv4/TCP/UDP packet.

**What checksums catch that parity misses:** most patterns of 2+ bit corruption, because they usually change the *numeric value* of a word, and that changes the sum. Checksums are a real, meaningful improvement over a single parity bit.

**What checksums still miss — and this is the important nuance:** addition is *commutative* and *insensitive to reordering*. If two 16-bit words in the data are swapped with each other, the sum — and therefore the checksum — is **completely unchanged**, even though every byte of the payload is now in the wrong place:

```
Original: word1 = 0x1234, word2 = 0x5678  → sum includes 0x1234 + 0x5678
Swapped:  word1 = 0x5678, word2 = 0x1234  → sum includes 0x5678 + 0x1234

Same sum. Same checksum. Reordering is invisible to a checksum.
```

Checksums also miss "compensating errors" — if one word's value increases by exactly X while another word's value decreases by exactly X, the sum is unchanged. These are real, if less common, corruption patterns, and they motivate the stronger, structurally different technique in Section 7.

**One more real nuance worth internalizing now, because it surprises people later:** the **IPv4 header checksum only covers the IP header — never the payload.** It was designed that way in the 1980s partly for speed (routers had to recompute it on every hop, since the TTL field changes at every hop — see Chapter 45) and partly because it was assumed the layers above (TCP, UDP) would separately check the payload's integrity, which they do, via their own checksums that cover a "pseudo-header" (source/destination IP, protocol number, length) plus the entire TCP/UDP segment.

---

## 7. CRC — Cyclic Redundancy Check

Checksums are cheap and much better than parity, but Section 6 exposed a real gap: they're blind to reordering, and — more importantly for real physical links — they aren't mathematically guaranteed to catch **burst errors**, which Section 2 established as the dominant real-world error pattern (a single noise event corrupting several consecutive bits). CRC closes that gap with a completely different mathematical tool: **polynomial division over GF(2)** (the two-element field {0, 1}, where addition and subtraction are both just XOR).

**Intuitive picture:** treat the entire message as one enormous binary number. Both sender and receiver agree in advance on a fixed "magic number" (the **generator polynomial**, `G`). The sender appends a few extra bits to the message, chosen so that the *whole thing* — message plus extra bits — divides evenly (remainder zero) by `G`. The receiver does the same division on what it received; if the remainder isn't zero, something changed in transit.

**Real-world analogy — and where it breaks:** this resembles the old-school arithmetic trick of "casting out nines" to catch mistakes when adding by hand (if a number and the sum of its digits disagree mod 9 by too much, you made an arithmetic error). The analogy is useful for the shape of the idea — "check divisibility, not equality" — but it breaks down in strength: casting out nines misses the majority of real transcription errors. CRC, because it operates in the algebra of binary polynomials rather than base-10 digit sums, can be **proven** — not just observed empirically — to catch entire mathematically-defined classes of errors with certainty, as you'll see below.

**The mechanics: binary polynomial division = long division using XOR instead of subtraction.**

A bit string like `1101` represents the polynomial `x³ + x² + 1` (each bit is a coefficient, either 0 or 1, of a power of x). "Dividing" one bit string by another is done exactly like grade-school long division, except every subtraction step is replaced by XOR (which conveniently never needs to "borrow," because in GF(2), subtraction and addition are the same operation).

**Full worked example.** Let the message be `M = 1101011011` (10 bits) and the agreed generator be `G = 10011` (5 bits — this represents the polynomial `x⁴ + x + 1`, a real, simple 4th-degree generator used here for a hand-workable example; production systems use standardized generators like CRC-32, introduced below).

Since `G` has 5 bits (degree `r = 4`), the sender first appends `r = 4` zero bits to the message, then divides by `G`:

```
Message:            1101011011
Append r=4 zeros:   1101011011 0000
Dividend to divide: 11010110110000   (14 bits)
Divisor (G):        10011            (5 bits)
```

XOR long division, step by step (bring down one bit at a time, XOR whenever the leading bit of the current remainder is 1):

```
  1 1 0 1 0 1 1 0 1 1 0 0 0 0    <- dividend
  1 0 0 1 1                       <- XOR (leading bit was 1)
  ---------------------------
  0 1 0 0 1 1                     <- bring down next bit (1)
    1 0 0 1 1                     <- XOR (leading bit is 1)
    ---------------------------
    0 0 0 0 0 1                   <- bring down next bit (1): 00001
                                     leading bit 0 → no XOR, carry forward
    0 0 0 0 1 0                   <- bring down next bit (0): 00010, leading 0
    0 0 0 1 0 1                   <- bring down next bit (1): 00101, leading 0
    0 0 1 0 1 1                   <- bring down next bit (1): 01011, leading 0
    0 1 0 1 1 0                   <- bring down next bit (0): 10110, leading 1
      1 0 0 1 1                   <- XOR
      -----------------------
      0 0 1 0 1                   <- result 00101, bring down next bit (0): 01010, leading 0
      0 1 0 1 0 0                 <- bring down next bit (0): 10100, leading 1
        1 0 0 1 1                 <- XOR
        ---------------------
        0 0 1 1 1                 <- result 00111, bring down last bit (0): 01110

Final remainder (last 4 bits): 1 1 1 0
```

**CRC (remainder) = `1110`.**

The sender transmits the original message with this remainder appended in place of the zeros:

```
Transmitted codeword = 1101011011 (message) + 1110 (CRC) = 11010110111110
```

**Receiver's check:** divide the *entire received codeword* (message + CRC together) by the same generator `G = 10011`. If nothing was corrupted, the remainder comes out to exactly `0000`. If any covered bit changed, the remainder will (with the guarantees below) be nonzero, and the frame is discarded or flagged.

**What CRC mathematically guarantees, not just usually catches:**

- **Every single-bit error**, with certainty.
- **Every double-bit error**, as long as the generator polynomial has at least two nonzero terms and the message length stays under a bound related to the generator (true for all standard generators at realistic frame sizes).
- **Every burst error of length ≤ r** (the generator's degree), with certainty — this is the property that directly answers Section 2's observation that real-world noise corrupts *runs* of adjacent bits.
- **The overwhelming majority of longer bursts and other error patterns**: the probability that a longer, unlucky error pattern is *undetected* is bounded at `1/2^r`. For CRC-32 (`r = 32`, used by Ethernet), that's roughly a **1-in-4-billion** chance of a miss — astronomically better than a checksum's blind spots.

**What CRC still can't catch:** an error pattern that, by sheer coincidence, is itself exactly a multiple of the generator polynomial. By design, standard generators are chosen so that this class of undetectable errors is vanishingly rare and doesn't correspond to any physically common corruption pattern (like the reordering or compensating-sum errors that trip up plain checksums).

**Real generator polynomials you'll actually meet:**

| Name | Degree (r) | Used in |
|---|---|---|
| CRC-8 | 8 | Some sensor/automotive buses (e.g., 1-Wire) |
| CRC-16-CCITT | 16 | Bluetooth, PPP framing, XMODEM |
| CRC-32 | 32 | Ethernet frame trailer (FCS), Wi-Fi (802.11), ZIP, PNG, gzip |

CRC-32's actual polynomial (in hex, representing degree-32 coefficients) is `0x04C11DB7`. You never need to hand-divide 32-bit polynomials in practice — every NIC has this built into hardware, computing it at line rate as bits stream past — but the *mechanism* is exactly the XOR long division shown above, just with a bigger generator.

---

## 8. Comparison — What Each Technique Catches and Misses

| Technique | Redundancy cost | Reliably catches | Reliably misses | Relative strength |
|---|---|---|---|---|
| Send twice (naive) | 100% (doubles data) | Any difference between the two copies | Identical corruption of both copies (correlated bursts) | Expensive, still weak |
| Single parity bit | 1 bit per block | Any odd number of flips | Any even number of flips | Weakest real scheme |
| 2D parity | ~1 bit per row + column | Single-bit errors (and can *locate* them) | Certain symmetric multi-bit patterns | Weak, but hints at correction |
| Checksum (Internet checksum) | 16–32 bits per packet | Most single/multi-bit value changes | Reordered words, compensating errors | Moderate, very cheap to compute |
| CRC | 8–32 bits per frame | All single/double-bit errors, all bursts ≤ r bits, ~99.9999998% of everything else (CRC-32) | Errors that are exact multiples of the generator (astronomically rare) | Strongest of the three, still cheap in hardware |

---

## 9. Where These Are Actually Used

- **Ethernet frames** (Chapter 28) end with a 4-byte **Frame Check Sequence (FCS)** — a CRC-32 over the entire frame. A switch or NIC that gets a bad FCS silently drops the frame; there's no automatic "please resend" at this layer (that's TCP's job, much higher up — Chapter 60).
- **Wi-Fi (802.11)** frames also carry a CRC-32 FCS, for the same reason — radio is noisier than copper, so detection matters even more.
- **IPv4 header checksum** (Chapter 36): a 16-bit one's-complement Internet checksum, but — as emphasized in Section 6 — it covers **only the header**, not the payload, and it must be recomputed at every router hop because the TTL field (Chapter 45) changes at each hop.
- **TCP and UDP checksums** (Chapters 58, 65): 16-bit Internet checksums covering a "pseudo-header" (source/dest IP, protocol, length) plus the entire segment — this is what actually protects your payload data end-to-end at the transport layer.
- **ICMP checksum** (Chapter 54): same Internet checksum algorithm, covering the ICMP message.
- **ZIP, PNG, gzip, and most archive/compression formats** use CRC-32 to detect file corruption.
- **Parity in memory**: pure single-bit parity RAM is mostly historical; modern **ECC memory** uses Hamming-derived SECDED (Single Error Correction, Double Error Detection) codes, covered fully in Chapter 20.

**Where the FCS physically sits in a frame** (full field-by-field breakdown of the Ethernet frame is Chapter 28's job — this is just to make the CRC's position concrete):

```
| Preamble | Dest MAC | Src MAC | EtherType | Payload (46-1500B) | FCS (CRC-32) |
   7 bytes    6 bytes    6 bytes    2 bytes         variable          4 bytes
                                                                          ^
                                                              covers everything
                                                              from Dest MAC onward
```

The FCS is computed over the destination MAC through the end of the payload, appended as a trailer, and checked by the receiving NIC in hardware before the frame is even handed up to the operating system — a corrupted frame never reaches software at all.

**Production usage note:** hardware CRC-32 computation happens at line rate — a 100 Gbps NIC computes and checks a CRC-32 over every frame without adding measurable latency, because the XOR long division in Section 7 maps directly onto a small, fully pipelined shift-register circuit. This is precisely why CRC, despite looking like "more math" than a checksum, is actually cheaper to implement in silicon than you'd expect, and why link layers universally chose it over a plain checksum once transistor budgets allowed.

---

## 10. Hands-On: Compute a Checksum and a CRC Yourself

**Quick check using standard tools.** Every Linux/macOS system ships a CRC utility:

```
$ echo -n "hello" | cksum
```

And Python's standard library exposes the exact CRC-32 Ethernet/ZIP uses:

```python
>>> import zlib
>>> hex(zlib.crc32(b"hello"))
'0xb6cc4292'
>>> hex(zlib.crc32(b"hellp"))     # one character changed
'0x2fda4441'
```

Notice how completely different the two CRC values are for a one-character change — this "avalanche" behavior (small input change, large output change) is a useful, if informal, sign of a strong detection code.

**Go implementation — Internet checksum (RFC 1071), the exact algorithm from Section 6:**

```go
package main

import "fmt"

// internetChecksum implements the one's-complement, end-around-carry
// checksum used by IPv4, TCP, UDP, and ICMP (RFC 1071).
func internetChecksum(words []uint16) uint16 {
	var sum uint32
	for _, w := range words {
		sum += uint32(w)
		if sum > 0xFFFF {
			// end-around carry: fold the overflow back in
			sum = (sum & 0xFFFF) + (sum >> 16)
		}
	}
	return ^uint16(sum) // one's complement of the final sum
}

func main() {
	data := []uint16{0xFFFF, 0x0001} // matches the 8-bit worked example, widened to 16 bits
	cs := internetChecksum(data)
	fmt.Printf("checksum: %#04x\n", cs)

	// verification: sum everything including the checksum -> should be 0xFFFF
	verify := internetChecksum(append(data, cs))
	fmt.Printf("verification (want 0xffff): %#04x\n", ^verify) // note: sum, not complement, should be all 1s
}
```

**Go implementation — CRC via XOR long division, mirroring Section 7's hand-worked example exactly:**

```go
package main

import "fmt"

// crcRemainder performs binary polynomial division (mod 2) of `message`
// (with r zero bits already appended) by `generator`, returning the
// r-bit remainder — the CRC value.
func crcRemainder(dividend []byte, generator []byte) []byte {
	r := len(generator) - 1
	work := make([]byte, len(dividend))
	copy(work, dividend)

	for i := 0; i < len(work)-r; i++ {
		if work[i] == 1 {
			for j := 0; j < len(generator); j++ {
				work[i+j] ^= generator[j] // XOR = "subtraction" in GF(2)
			}
		}
	}
	return work[len(work)-r:]
}

func main() {
	message := []byte{1, 1, 0, 1, 0, 1, 1, 0, 1, 1}   // 1101011011
	generator := []byte{1, 0, 0, 1, 1}                // x^4 + x + 1
	r := len(generator) - 1

	dividend := append(append([]byte{}, message...), make([]byte, r)...)
	crc := crcRemainder(dividend, generator)
	fmt.Println("CRC remainder:", crc) // expect [1 1 1 0] — matches Section 7
}
```

**Experiment to run yourself:** take the transmitted codeword `11010110111110` from Section 7, flip a single bit anywhere in it, and re-run the division (by hand or by adapting the Go code above). Confirm the remainder is no longer `0000`. Then try flipping *two adjacent* bits (simulating a tiny burst) and confirm CRC still catches it — then compare against what a single parity bit over the same data would have caught (Section 4's worked example shows exactly this kind of even-count flip slipping past parity).

---

## 11. Common Misconceptions

- **"A passing checksum/CRC means the data is definitely correct."** No — it means the data is *probably* correct. Every scheme in this chapter has a nonzero (if tiny, for CRC-32) probability of missing an error. "Detected no error" is a strong probabilistic statement, not a proof.
- **"CRC-32 is a security mechanism / a way to detect tampering."** No. CRC (and even cryptographic-looking checksums) can be trivially recomputed by anyone, including an attacker who deliberately modifies data — there's no secret involved. CRC defends against *random* noise, not an adversary. Tamper-evidence needs cryptographic hashes and signatures (Chapter 80), which are a completely different tool built for a completely different threat model.
- **"More redundant bits automatically means the receiver can fix the error."** No — detection and correction are different capabilities. Every technique in this chapter tells you *that* something is wrong; none of them tell you *what the correct value was*. That's Chapter 20's job, and it costs meaningfully more redundancy to achieve.
- **"Parity catches most real errors."** No — it only catches odd counts of flips, and Section 2 established that real noise often flips multiple adjacent bits at once (bursts), which are just as likely to be an even count as odd. Parity's real-world hit rate on burst errors is close to a coin flip.
- **"The IP checksum protects my data end-to-end."** No — as Section 9 stressed, the IPv4 header checksum covers only the header, is recomputed at every hop, and says nothing about payload integrity. TCP/UDP checksums are what actually cover your data, and even those are 16-bit — far weaker than the CRC-32 protecting the same bits at the link layer.

---

## 12. Interview Questions & Model Answers

**Q1 (Beginner): What's the difference between error detection and error correction?**

"Error detection tells you *that* a frame or packet was corrupted in transit — it doesn't tell you which bits changed or what they should have been, only that a check (parity, checksum, or CRC) failed. Error correction goes further: it adds enough redundant information that the receiver can pinpoint the specific bit(s) that changed and flip them back, without asking the sender to retransmit. Detection is cheap (a few bits per frame) and is what Ethernet, Wi-Fi, IP, TCP, and UDP all do on every single packet. Correction is more expensive and is reserved for situations where retransmission is slow, costly, or impossible — satellite links, broadcast video, deep-space probes, ECC memory — which is exactly what Chapter 20 covers."

**Q2 (Intermediate): Why does CRC catch burst errors that a checksum might miss, when both add redundant bits to the message?**

"A checksum treats the message as a set of independent numeric words and adds them — this makes it blind to certain structural changes, like two words being swapped, or one word increasing while another decreases by the same amount, because addition doesn't care about position or order. CRC is fundamentally different: it treats the *entire* message as one long binary polynomial and checks divisibility by a fixed generator polynomial using XOR-based long division. Because of how polynomial division mod 2 works, the CRC value is a function of the *exact position and pattern* of every bit, not just a positionally-blind sum. This gives CRC a provable guarantee — not just an empirical tendency — to catch every burst error shorter than or equal to the generator's degree, which is precisely the error pattern real-world noise events (a lightning strike, a collision, a brief electromagnetic pulse) actually produce."

**Q3 (Advanced): The IPv4 header checksum only covers the IP header, not the payload. Why was it designed that way, and what are the practical consequences?**

"Two reasons, one historical/performance-driven and one architectural. Performance: the IPv4 TTL field decrements at every router hop (Chapter 45), which means the header checksum must be recomputed at every hop too — in 1980s router hardware, minimizing what had to be recomputed on the fast path mattered a lot, and the payload is usually much larger than the 20-byte header. Architecturally: the designers assumed — correctly — that end-to-end payload integrity was the transport layer's job, not the network layer's, since IP is explicitly a best-effort, unreliable protocol (Chapter 36) and shouldn't duplicate work TCP/UDP already do via their own checksums covering a pseudo-header plus the full segment. The practical consequence is that if you're ever debugging a scenario where IP-layer routing succeeds but payload data looks corrupted, you have to look at the TCP/UDP checksum. And because that checksum is only 16 bits — meaningfully weaker than the CRC-32 FCS the link layer already applied and stripped off before the packet ever reached IP — extremely rare silent data corruption that slips past every layer's check is a real, documented phenomenon in very large-scale systems (this is one motivation for application-level integrity checks, like content hashes, on top of everything the network stack already does)."

---

## 13. Exercises

### Easy

1. Compute the even-parity bit for the 8-bit data `11001101`. Show your work.
2. Given data `1010` with even parity bit `1` was sent as `10101`, and the receiver got `10001`, does the parity check catch this error? Why or why not?
3. Using the Internet checksum method from Section 6 (8-bit words for simplicity), compute the checksum for words `00000000` and `00000000`. What potential ambiguity does one's complement arithmetic avoid here compared to plain binary addition?

### Medium

4. Take the 2D parity grid from Section 5. Flip the bit at row2, col1 (originally `1`). Recompute the row and column parities and show which row/column checks fail, confirming the intersection correctly identifies the flipped bit.
5. Using generator `G = 1011` (degree 3) and message `M = 101100`, compute the CRC remainder by hand using XOR long division (append 3 zero bits first). Show every step.
6. Explain, with a concrete pair of 16-bit words, a specific corruption that an Internet checksum would fail to detect. Then explain whether a CRC-32 over the same words would catch it, and why.

### Hard

7. Prove informally (using the polynomial-division framing from Section 7) why a CRC with an r-bit generator is guaranteed to detect any burst error of length ≤ r. (Hint: think about what a burst error "looks like" as a polynomial added to the original message, and what has to be true of that error polynomial for it to be silently divisible by the generator.)
8. Write a program (Go or Python) that simulates a noisy channel: generate 10,000 random 64-bit messages, inject a random burst error of length 1–8 bits into each, and measure the empirical detection rate for (a) a single parity bit over the whole message, (b) a 16-bit Internet checksum, and (c) a CRC-16. Compare your empirical results to the theoretical guarantees from Section 8.
9. Ethernet's CRC-32 FCS is checked by the NIC hardware, and a frame that fails the check is silently dropped — there is no "corrupted frame, please resend" message sent back at the link layer. Explain, tracing forward to what you know is coming in Volume 9 (TCP), how a dropped Ethernet frame eventually still results in the correct data arriving at the application, and estimate (using round-trip time reasoning) how much slower that recovery path is compared to if link-layer retransmission existed.

---

## Summary

| Term | Meaning |
|---|---|
| Bit Error Rate (BER) | Fraction of bits corrupted in transit on a real channel; always nonzero even below Shannon's limit |
| Parity bit | One extra bit making the total 1-count even/odd; catches only odd numbers of flips |
| Two-dimensional parity | Row + column parity; can detect and locate a single-bit error |
| Checksum (Internet checksum) | One's-complement sum of 16-bit words with end-around carry; catches most value-changing errors, misses reordering |
| CRC (Cyclic Redundancy Check) | Remainder of polynomial division (mod 2, i.e. XOR) by a fixed generator polynomial; guarantees detection of all bursts ≤ generator degree |
| Generator polynomial | The fixed "magic divisor" both sender and receiver agree on for CRC (e.g., CRC-32's `0x04C11DB7`) |
| FCS (Frame Check Sequence) | The CRC-32 trailer on every Ethernet and Wi-Fi frame |
| Detection vs. correction | Detection says "something's wrong"; it does not say what or where — that's Chapter 20 |

Every technique in this chapter answers "is something wrong?" — none of them answer "what should it have been?" Chapter 20 picks up exactly there: Hamming codes show how a receiver can locate *and fix* a flipped bit without asking the sender for anything, and Forward Error Correction shows why that ability is not a curiosity but an operational necessity for satellite links, fiber backbones, and Wi-Fi.
