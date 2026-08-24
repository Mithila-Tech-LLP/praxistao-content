# Chapter 36: IPv4 Addresses — What They Are and Why They Exist

> **"A MAC address tells you who a device is. An IP address tells you where it is. The Internet could never have scaled past a few thousand machines without that second idea."**

---

## Table of Contents

1. [The Problem — Why MAC Addresses Don't Scale to a Planet](#1-the-problem--why-mac-addresses-dont-scale-to-a-planet)
2. [A Real-World Analogy: Phone Numbers and Postal Codes](#2-a-real-world-analogy-phone-numbers-and-postal-codes)
3. [The Naive Fix, and Why It Fails](#3-the-naive-fix-and-why-it-fails)
4. [The Real Solution: A 32-Bit Hierarchical Address](#4-the-real-solution-a-32-bit-hierarchical-address)
5. [Dotted-Decimal Notation](#5-dotted-decimal-notation)
6. [Binary to Decimal: Worked Conversions](#6-binary-to-decimal-worked-conversions)
7. [Where the Address Lives: A Peek at the IP Header](#7-where-the-address-lives-a-peek-at-the-ip-header)
8. [How Many Addresses Are There?](#8-how-many-addresses-are-there)
9. [MAC vs. IP Addressing, Contrasted](#9-mac-vs-ip-addressing-contrasted)
10. [A Hands-On Experiment](#10-a-hands-on-experiment)
11. [Common Misconceptions](#11-common-misconceptions)
12. [What's Simplified Here](#12-whats-simplified-here)
13. [Interview Questions & Model Answers](#13-interview-questions--model-answers)
14. [Exercises](#14-exercises)
15. [Summary](#15-summary)

---

## 1. The Problem — Why MAC Addresses Don't Scale to a Planet

Chapter 29 gave every network interface a 48-bit MAC address: a globally unique, burned-in identifier like `3C:22:FB:9A:11:04`. It works beautifully inside a single LAN — a switch learns which MAC address sits on which port (Chapter 31) and forwards frames accordingly.

Now stretch that idea to the entire planet. Suppose the Internet tried to route traffic using MAC addresses alone. A router somewhere in Mumbai receives a frame destined for `3C:22:FB:9A:11:04` and has to decide which of its several links to send it out on. What does it know about that address?

Nothing useful. The first 24 bits of a MAC address (the OUI, Chapter 29) identify the *manufacturer* of the network card — Intel, Realtek, Apple — not the device's location. A laptop with an Intel network card could be sitting in Mumbai, Lagos, or São Paulo; the address gives no clue. There is no way to look at `3C:22:FB:9A:11:04` and conclude "that's somewhere near Europe, send it west."

This is the core problem: **MAC addresses are flat**. Every address is equally "close" to every other address in the sense that none of them encode any structural information about where the device sits in the network. A flat address space can only be routed by one of two methods:

- **Global flooding** — every router keeps trying every link until the frame reaches its destination. This is roughly what a single switch does with an unknown MAC address (Chapter 31's "flooding" behavior) — and it barely scales to a building.
- **A complete directory** — every router keeps a table mapping every possible MAC address in existence to the correct next link. With roughly 2^48 (281 trillion) possible MAC addresses and billions of them active on the Internet at once, no router has anywhere near enough memory to hold that table, and every added or removed device would require updating that table everywhere.

Neither approach scales past a single LAN. The Internet needed an addressing scheme where a device's address itself encodes enough information about *where in the network* it lives that a router can make a forwarding decision without knowing about every device on Earth — just by looking at a prefix of the address.

That scheme is **IP addressing**, and this chapter covers its first (and still dominant) version: **IPv4**.

---

## 2. A Real-World Analogy: Phone Numbers and Postal Codes

Before computers, humans solved exactly this problem twice.

**Postal addresses** are hierarchical: country, then state/province, then city, then street, then house number. A sorting facility in Mumbai doesn't need to know which specific house `221B Baker Street, London, UK` refers to. It only needs to read "UK" and route the letter toward London. Once it reaches a London sorting office, *that* office reads "Baker Street" and routes it locally. Nobody, at any stage, needed a giant table of every house address on Earth — the address itself, read a piece at a time, tells you which direction to go next.

**Phone numbers** work the same way. `+91 22 6612 3456` breaks down as: `+91` (country code: India), `22` (city/area code: Mumbai), `6612 3456` (the actual subscriber line). An exchange in London routes a call to `+91...` toward India without knowing anything about the specific subscriber — the *country code* alone is enough for the first hop.

This is the essential idea: **a hierarchical address lets you make a correct routing decision using only part of the address**, and defer the rest of the decision to someone closer to the destination. It's exactly analogous to Chapter 24's argument for layering — decompose a hard global problem into smaller local ones.

MAC addresses have no such hierarchy tied to location. Two devices could have wildly different MAC addresses (different manufacturers) yet sit on the same desk, or nearly identical MAC addresses (same manufacturer, sequential serial numbers) yet sit on opposite sides of the planet. There is no "network area code" baked into a MAC address.

IP addresses fix exactly this.

---

## 3. The Naive Fix, and Why It Fails

Suppose you tried to patch MAC addressing instead of replacing it: keep 48-bit flat addresses, but give every router a complete table mapping every address to a next hop, and have devices announce themselves whenever they move.

This is worse than it sounds:

- **Table size.** Even restricting to "addresses actually in use" (a few billion, not 2^48), that is still billions of entries every core router on the Internet would need to hold and search on every single packet, in nanoseconds, with the table constantly changing as devices connect, disconnect, and move between networks (a laptop moving from a home Wi-Fi network to a coffee shop's).
- **Update storms.** Every time any device anywhere joined, left, or moved, that change would need to propagate to every router capable of reaching it — potentially every router on Earth. Chapter 46 will show that even *summarized* routing information takes real engineering to propagate efficiently; propagating billions of individual device locations is not feasible at any layer.
- **No way to summarize.** Because a flat address carries no structural information, there is no way to represent "all 3,000 devices behind this one office building's router" as a single compact routing table entry. You'd need one entry per device, forever.

This naive fix isn't a strawman — it is, in effect, what a single Ethernet switch already does at LAN scale (a MAC address table, Chapter 31), and it works there specifically *because* a LAN is small: one building, a few hundred to a few thousand devices, one administrative owner. The naive fix fails specifically at *global* scale, which is the scale the Internet actually operates at.

What's needed is an address that is assigned based on *where a device sits in the network topology* — not who manufactured its network card — so that the address itself can be summarized, the way `+91` summarizes "somewhere in India" without listing every phone number in the country.

---

## 4. The Real Solution: A 32-Bit Hierarchical Address

**Internet Protocol version 4 (IPv4)**, standardized in RFC 791 (1981), assigns every device on an IP network a **32-bit address**, split conceptually into two parts:

```
 +------------------------+------------------------+
 |     NETWORK portion    |      HOST portion       |
 +------------------------+------------------------+
  <-- identifies WHICH -->  <-- identifies WHICH  -->
  <--  network a host   -->  <--  specific host   -->
  <--    belongs to     -->  <--  on that network -->
```

The *network portion* plays the role of the country code or area code: a router far away only needs to look at this part to decide "send it in that general direction." The *host portion* plays the role of the subscriber's own line: it only matters once the packet has arrived at the correct network, where the local router or switch delivers it to the exact device.

Exactly where the line falls between "network" and "host" is not fixed — it's configurable per network, using a mechanism called a **subnet mask**, which is the entire subject of Chapter 37. For now, the important idea is simply that the split exists, and that it is what makes IP addresses hierarchical where MAC addresses are flat.

This is the single biggest structural difference from MAC addressing:

| | MAC address (Ch. 29) | IPv4 address |
|---|---|---|
| Length | 48 bits | 32 bits |
| Assigned by | Manufacturer (burned into hardware) | Network administrator / ISP / DHCP (Ch. 55) |
| Structure | Flat (OUI identifies vendor, not location) | Hierarchical (network portion + host portion) |
| Changes when device moves networks? | No — same MAC address anywhere | Yes — a device gets a *different* IP address on a different network |
| Enables routing by prefix? | No | Yes — this is the entire point |

That last row is the payoff. Because the network portion is a contiguous prefix of the address, a router doesn't need to know about every host in the world — it only needs to know, for each network prefix it has heard of, which direction to send traffic. Chapter 45 formalizes this as **longest prefix match**, the actual algorithm every IP router runs on every packet. None of it would be possible without the hierarchical structure introduced right here.

---

## 5. Dotted-Decimal Notation

A 32-bit number written out in raw binary is miserable for humans to read or type:

```
11000000101010000000000100001010
```

Is that a valid address? Where does one "field" end and the next begin? To make IPv4 addresses human-manageable, the 32 bits are split into four groups of 8 bits each — called **octets** — and each octet is written as its decimal value (0–255), separated by dots. This is **dotted-decimal notation**.

```
 32 bits total, split into four 8-bit octets:

 11000000 . 10101000 . 00000001 . 00001010
 <------>   <------>   <------>   <------>
  octet 1    octet 2    octet 3    octet 4

 converted to decimal:

    192    .   168    .    1     .    10
```

So `11000000101010000000000100001010` and `192.168.1.10` are the *exact same 32-bit number* — dotted-decimal is purely a human-readable shorthand. A computer never actually "sees" the dots; internally the address is just 32 bits (usually stored as a single 32-bit unsigned integer or a 4-byte array).

Because each octet is 8 bits, and 8 bits can represent 2^8 = 256 distinct values (0 through 255), **every octet in a valid IPv4 address must be between 0 and 255**. `192.168.1.999` is not a valid IPv4 address — 999 cannot be represented in 8 bits (the largest 8-bit value is 255). This single fact — "each part is 0–255 because each part is 8 bits" — is the most common thing people memorize without understanding *why*; now you know why.

---

## 6. Binary to Decimal: Worked Conversions

Every octet-to-decimal conversion uses the same idea: each of the 8 bit positions has a fixed **place value**, powers of two from left (most significant) to right (least significant):

```
 bit position:   7    6    5    4    3    2    1    0
 place value:   128   64   32   16   8    4    2    1
```

To convert binary to decimal: add up the place values of every position that has a `1`. To convert decimal to binary: repeatedly subtract the largest place value that still fits, marking that position `1`, and mark every other position `0`.

### Worked Example 1 — Binary to Decimal

Convert `10001110` to decimal.

```
 bit:          1    0    0    0    1    1    1    0
 place value: 128   64   32   16   8    4    2    1
 contributes: 128    -    -    -   8    4    2    -

 128 + 8 + 4 + 2 = 142
```

So `10001110` = `142`.

### Worked Example 2 — Decimal to Binary

Convert `250` to binary. Subtract the largest place value that fits, in order:

```
 250 - 128 = 122   → bit 7 = 1
 122 - 64  = 58    → bit 6 = 1
  58 - 32  = 26    → bit 5 = 1
  26 - 16  = 10    → bit 4 = 1
  10 -  8  =  2    → bit 3 = 1
   2 -  4        → doesn't fit → bit 2 = 0
   2 -  2  =  0    → bit 1 = 1
   0 -  1        → doesn't fit → bit 0 = 0

 result: 11111010
```

Check by adding back: 128+64+32+16+8+2 = 250. Correct.

### Worked Example 3 — A Full Address, Both Directions

Take the real, public IP address `142.250.80.46` (one of Google's) and convert every octet to binary:

```
 142 = 128 + 8 + 4 + 2                     = 10001110
 250 = 128 + 64 + 32 + 16 + 8 + 2          = 11111010
  80 = 64 + 16                             = 01010000
  46 = 32 + 8 + 4 + 2                      = 00101110

 142.250.80.46  =  10001110.11111010.01010000.00101110
```

And going the other direction, take the binary form `00001000.00001000.00001000.00001000` and convert each octet back to decimal:

```
 00001000 = 8   (only bit 3, value 8, is set)

 result: 8.8.8.8    (a real address — one of Google's public DNS resolvers,
                      which Chapter 66 will meet again)
```

Practice this until it's automatic — every worked example in Chapters 37 through 39 depends on being able to flip an octet between decimal and binary without hesitation, because that's *exactly* the operation a router or a network engineer performs to figure out which network an address belongs to.

### Worked Example 4 — A Table of Conversions

To build real fluency, here is a batch of conversions worked in full, mixing well-known public addresses with arbitrary practice values:

| Decimal | Binary | Check |
|---|---|---|
| 1.1.1.1 | `00000001.00000001.00000001.00000001` | Cloudflare's public DNS resolver |
| 208.67.222.222 | `11010000.01000011.11011110.11011110` | 208 = 128+64+16 = `11010000`; 67 = 64+2+1 = `01000011`; 222 = 128+64+16+8+4+2 = `11011110` |
| 93.184.216.34 | `01011101.10111000.11011000.00100010` | 93 = 64+16+8+4+1 = `01011101`; 184 = 128+32+16+8 = `10111000` |
| 255.255.255.255 | `11111111.11111111.11111111.11111111` | The IPv4 limited broadcast address (all bits set) — previewed here, covered fully in Chapter 40 |
| 0.0.0.0 | `00000000.00000000.00000000.00000000` | The "unspecified" address — also covered in Chapter 40 |

Notice the two extremes at the bottom of the table: `255.255.255.255` (every bit is 1) and `0.0.0.0` (every bit is 0). These aren't arbitrary — they're the two addresses you get when you push the binary-to-decimal conversion to its limits, and both turn out to have special reserved meanings, which is exactly why Chapter 38's subnetting math is so careful about never handing an all-zeros or all-ones host portion to an actual device (Chapter 38, Section 4).

### Deep Dive: How Software Actually Stores and Converts These Addresses

Everything above was done by hand to build intuition, but real software never manually subtracts place values — it uses a couple of standard, well-defined operations, and it's worth seeing what those look like, because they are the literal implementation of the same idea.

Internally, an IPv4 address is almost always stored as a 32-bit unsigned integer, with the four octets packed in a specific byte order called **network byte order** (big-endian: the most significant octet, the first one you'd type, comes first in memory) — this matters because different CPU architectures (x86 is little-endian internally) must all agree on one wire format, or two machines would disagree about which octet is which.

In Go, the standard library's `net` package handles this conversion for you, but you can also do it by hand to see the mechanism directly:

```go
package main

import (
	"fmt"
	"net"
)

func main() {
	ip := net.ParseIP("192.168.1.10").To4() // parses dotted-decimal into 4 raw bytes
	fmt.Printf("Bytes: %v\n", ip)           // [192 168 1 10]

	// Pack the 4 bytes into one 32-bit unsigned integer, most significant byte first —
	// this is exactly "network byte order," and exactly what a router or the kernel
	// works with internally, never the dotted-decimal string.
	var asUint32 uint32
	for _, octet := range ip {
		asUint32 = (asUint32 << 8) | uint32(octet)
	}
	fmt.Printf("As one 32-bit number: %d\n", asUint32) // 3232235786

	// And back the other way: extract each octet with a shift and a mask —
	// the exact inverse of the packing above.
	back := net.IPv4(
		byte(asUint32>>24), byte(asUint32>>16),
		byte(asUint32>>8), byte(asUint32),
	)
	fmt.Println(back.String()) // 192.168.1.10
}
```

Run this and `192.168.1.10` round-trips to `3232235786` and back. That number, `3232235786`, is the exact same 32-bit value this chapter has been writing as `192.168.1.10` and as `11000000.10101000.00000001.00001010` the whole time — three different notations for one underlying bit pattern. This is also why some databases and older logging systems store IP addresses as a single integer column: it's a smaller, faster-to-index representation of literally the same information, and converting between the two forms is nothing more than the shift-and-mask arithmetic above.

---

## 7. Where the Address Lives: A Peek at the IP Header

Chapter 27 introduced encapsulation: every layer wraps the layer above it in its own header. The IP layer's header is where these 32-bit addresses actually travel, once per packet, in two dedicated 32-bit fields:

```
  0                   1                   2                   3
  0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |Version|  IHL  |   ...other fields covered in later chapters... |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                       Source IP Address (32 bits)            |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 |                    Destination IP Address (32 bits)          |
 +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

Every single IPv4 packet, no matter what it carries above it (TCP, UDP, ICMP), has exactly one 32-bit **source address** field and one 32-bit **destination address** field in its header. Every router between source and destination reads the destination address field, and — as Chapter 45 details — matches it against its routing table to decide which link to forward the packet out of. The source address field is what lets the eventual reply find its way back.

Fields like TTL, protocol number, and header checksum sit alongside these two address fields; they matter for later chapters (TTL for Chapter 45's forwarding and Chapter 54's traceroute) but aren't the concern of this chapter. What matters here is simply: the address you've been reading in dotted-decimal is not decoration — it is a literal 32-bit value sitting at a fixed byte offset in every packet's header, read by every router along the path.

---

## 8. How Many Addresses Are There?

A 32-bit number has 2^32 possible values:

```
 2^32 = 4,294,967,296

 roughly 4.3 billion distinct IPv4 addresses, total, ever, worldwide.
```

That sounded like plenty in 1981, when the number of computers on Earth was in the thousands. It does not sound like plenty against ~8 billion humans, most carrying multiple IP-connected devices (phone, laptop, smart TV, watch), plus billions more servers, routers, IoT sensors, and cars. The practical, exhaustion-driven consequences of this ceiling — and the various strategies (NAT in Chapter 41, and ultimately IPv6 in Chapter 42) that have kept IPv4 usable decades past when it should have run out — are covered in depth later in this volume. For now, just note the number, because Chapters 38 and 39 are fundamentally about not wasting a scarce 4.3-billion-address budget.

The exhaustion wasn't a distant hypothetical — it's a dated historical fact, tracked by the five Regional Internet Registries (RIRs) that hand out address blocks to ISPs and large organizations:

| Registry | Region | Free pool exhausted |
|---|---|---|
| IANA (the global pool feeding all RIRs) | Global | February 2011 |
| APNIC | Asia-Pacific | April 2011 |
| RIPE NCC | Europe, Middle East | September 2012 |
| LACNIC | Latin America, Caribbean | June 2014 |
| ARIN | North America | September 2015 |
| AFRINIC | Africa | 2020 (last to exhaust) |

Every one of these exhaustion dates fell *decades* after 1981 — far longer than the original designers of a 32-bit address field ever expected the protocol to remain in service. That the Internet kept running smoothly through and after every one of these dates is a direct testament to how effectively NAT (Chapter 41) and disciplined, waste-minimizing allocation (Chapters 38–39) stretched a fixed 4.3-billion-address budget across a network that grew from thousands of hosts to tens of billions of connected devices.

---

## 9. MAC vs. IP Addressing, Contrasted

It's worth being explicit about how these two addressing systems, both introduced by this course, divide responsibility — because Chapter 53 (ARP) exists entirely to bridge the gap between them:

| Question | MAC address answers | IP address answers |
|---|---|---|
| "Who, physically, is this?" | Yes — a specific network interface | No |
| "Which network segment is this device on?" | No | Yes — via the network portion (Ch. 37) |
| "How do I deliver this within one LAN?" | Yes — switches forward by MAC (Ch. 31) | No — IP doesn't know about switch ports |
| "How do I deliver this across the entire Internet?" | No — no global structure to route by | Yes — via hierarchical prefixes (Ch. 45) |
| Assigned by | Hardware manufacturer, once, forever | Network administrator or DHCP (Ch. 55), can change |
| Layer (OSI, Ch. 25) | Layer 2 (Data Link) | Layer 3 (Network) |

Neither address replaces the other — every packet traveling across a LAN segment is simultaneously wrapped in both: an Ethernet frame (MAC addresses, for this one physical hop) carrying an IP packet (IP addresses, for the entire end-to-end journey). Chapter 35's full LAN trace showed this in action; Chapter 53 shows exactly how a device figures out *which* MAC address corresponds to a given IP address on its local network.

---

## 10. A Hands-On Experiment

You can see your own device's IPv4 address, and prove to yourself that it's a genuinely different concept from your MAC address, in under a minute.

**On macOS/Linux:**

```bash
$ ifconfig en0
en0: flags=8863<UP,BROADCAST,SMART,RUNNING,SIMPLEX,MULTICAST> mtu 1500
        ether 3c:22:fb:9a:11:04
        inet 192.168.1.42 netmask 0xffffff00 broadcast 192.168.1.255
```

**On Linux (modern):**

```bash
$ ip addr show eth0
2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500
    link/ether 3c:22:fb:9a:11:04 brd ff:ff:ff:ff:ff:ff
    inet 192.168.1.42/24 brd 192.168.1.255 scope global eth0
```

**On Windows:**

```
> ipconfig /all
Ethernet adapter Ethernet:
   Physical Address. . . . . . . . . : 3C-22-FB-9A-11-04
   IPv4 Address. . . . . . . . . . . : 192.168.1.42
   Subnet Mask . . . . . . . . . . . : 255.255.255.0
```

Notice both addresses printed side by side: `ether`/`Physical Address` (MAC, Chapter 29) and `inet`/`IPv4 Address` (IP, this chapter). Now do this experiment: connect the same laptop to a different Wi-Fi network — home, then a coffee shop, then a mobile hotspot — and run the command again each time.

**Prediction:** the MAC address (`3c:22:fb:9a:11:04`) will stay identical every time — it's burned into the network card. The IPv4 address will almost certainly be *different* on each network (perhaps `192.168.1.42` at home, `10.0.0.15` at the coffee shop, `172.20.10.3` on the hotspot) — because it was assigned based on which network you joined, not which piece of hardware you are. That single observation is the entire chapter, demonstrated on your own machine.

You can also do the binary-to-decimal conversion from Section 6 in one line, in Python, using the standard `ipaddress` module — useful for checking your own by-hand work quickly:

```python
>>> import ipaddress
>>> addr = ipaddress.IPv4Address("142.250.80.46")
>>> format(int(addr), "032b")             # the raw 32-bit integer, in binary
'10001110111110100101000000101110'
>>> int(addr)                              # the same value, as a plain integer
2398769198
```

Split that binary string into four 8-bit chunks (`10001110`, `11111010`, `01010000`, `00101110`) and it matches Section 6's Worked Example 3 exactly — confirming the by-hand conversion was correct, using a tool instead of a place-value table.

---

## 11. Common Misconceptions

- **"An IP address identifies a device."** Not quite — it identifies a device's *point of attachment to a particular network*. The same laptop gets a different IP address on every network it joins (as the experiment above shows), while its MAC address never changes. IP addresses answer "where," not "who."
- **"192.168.1.999 could theoretically be a valid address if we just made bigger networks."** No — each octet is fundamentally 8 bits, capable of representing exactly 0–255. There's no "bigger octet" option within IPv4; that constraint is why the whole system was eventually outgrown, prompting IPv6 (Chapter 42).
- **"The four numbers in an IP address are somehow arbitrary groupings, like a phone number's dashes."** They're not arbitrary — they correspond exactly to 8-bit boundaries in the underlying 32-bit binary number. Dotted-decimal is a direct, lossless transcription of the binary, not a separate format layered on top.
- **"IP addresses are randomly assigned."** They are anything but random — as Chapters 37–39 will show in detail, the specific bits of an address are chosen deliberately to reflect network topology, precisely so routers can make sense of them in bulk.

---

## 12. What's Simplified Here

This chapter deliberately did not cover: the historical Class A/B/C system (Chapter 39 covers it, and why it was replaced), private vs. public address ranges and reserved blocks like loopback and multicast (Chapter 40), how a device actually *obtains* its IP address in the first place (Chapter 55, DHCP), and IPv6's entirely different, much larger address format (Chapter 42). Each of those builds directly on the 32-bit, dotted-decimal foundation laid here.

---

## 13. Interview Questions & Model Answers

**Beginner: "Why can't we just route Internet traffic using MAC addresses, since every device already has one?"**

MAC addresses are flat — the address doesn't encode any information about where in the world the device is, only who manufactured its network card. A router would need a complete table mapping every device on Earth to a direction, with no way to summarize entries, and no way to keep that table current as devices connect, disconnect, and move. IP addresses fix this by being hierarchical: part of the address (the network portion) reflects the device's position in the network topology, so routers only need to know about network *prefixes*, not individual devices.

**Beginner: "What is dotted-decimal notation, and why do we use it?"**

It's a way of writing a 32-bit IPv4 address as four decimal numbers (0–255), separated by dots, where each number represents one 8-bit octet of the underlying binary address. It exists purely for human readability — computers work with the raw 32-bit binary value; dotted-decimal is a lossless shorthand.

**Intermediate: "Why must every octet in an IPv4 address be between 0 and 255?"**

Because each octet is exactly 8 bits, and 8 bits can represent 2^8 = 256 distinct values: 0 through 255. There's no way to represent, say, 300 in a single octet — it would require 9 bits. This isn't a rule imposed on top of the address format; it falls directly out of the binary structure.

**Intermediate: "What structurally distinguishes an IP address from a MAC address?"**

An IP address is hierarchical — split into a network portion and a host portion (Chapter 37) — while a MAC address is flat. This means IP addresses can be aggregated into prefixes that a router can reason about in bulk (Chapters 45 and 50), while MAC addresses cannot be meaningfully summarized at all; a switch has to learn every individual MAC address it has seen (Chapter 31).

**Advanced: "IPv4 has roughly 4.3 billion addresses. Explain, structurally, why that number is fixed and cannot be casually increased."**

The address is a fixed-width 32-bit field in every IPv4 packet's header (RFC 791). The field width is baked into the wire format that every router, host, and piece of network hardware on Earth parses; widening it isn't a software patch, it's a fundamentally different protocol, because every device that reads that header would need to agree on the new width simultaneously. That's precisely why the actual fix (IPv6, Chapter 42) is a new protocol version running in parallel with IPv4 rather than a revision to it, and why the transition (Chapter 43) has taken decades.

**Advanced: "A device's IP address changes when it joins a different network, but its MAC address does not. What does this imply about where each address is 'assigned,' and why does that matter for mobility?"**

A MAC address is baked into the network interface at manufacture time and is meaningless as location information. An IP address is assigned by whichever network the device currently attaches to (statically, or dynamically via DHCP, Chapter 55) precisely *because* the address must reflect the device's current position in the network's hierarchy for routing to work. This is also exactly why mobility is hard at the IP layer — moving to a new network fundamentally means getting a new address, which is why techniques like mobile IP, and more commonly, NAT and higher-layer session persistence, exist to paper over the disruption an address change causes to an in-progress connection.

---

## 14. Exercises

### Easy

1. Convert the following binary octets to decimal: `01100100`, `11110000`, `00000001`, `10000000`.
2. Convert the following decimal values to 8-bit binary: `17`, `200`, `255`, `1`.
3. Is `10.20.30.400` a valid IPv4 address? Explain why or why not, in terms of bits.

### Medium

4. Convert the IP address `216.58.211.14` fully to its 32-bit binary form, octet by octet, showing your subtraction work for at least two octets.
5. Convert the binary address `11011000.00111010.00001100.01100100` to dotted-decimal.
6. Explain, in your own words and without using the word "location," why a MAC address cannot be used the way an IP address is used for Internet-wide routing.

### Hard

7. Two laptops have MAC addresses that differ only in their last byte (same manufacturer, sequential units off the assembly line). Is it safe to assume they're on the same network? Justify your answer using what you know about OUIs (Chapter 29) versus IP's network portion (this chapter, expanded in Chapter 37).
8. A junior engineer claims: "IPv6 will eventually need a Chapter 39-style CIDR fix too, once it fills up, just like IPv4 did." Using only the numbers from Section 8 of this chapter and what you can infer about a fixed-width address field, argue for or against that claim (a rough comparison of 2^32 vs. what you'd expect from a 128-bit space, covered fully in Chapter 42, is enough to reason about it).
9. Using the shift-and-mask technique from Section 6's deep dive (not a library function), write out, step by step, how you would convert the 32-bit integer `167772161` back into dotted-decimal by hand. (Hint: divide by 256 repeatedly, keeping remainders, exactly the inverse of packing four bytes into one integer.)

---

## Production Usage Notes

A few places this chapter's ideas show up directly in day-to-day engineering work: IP addresses are frequently stored as a 32-bit unsigned integer column (`INET`/`INT UNSIGNED`) in databases that need to do fast range queries (e.g., "which rows fall within this CIDR block") — range comparisons on a single integer are far cheaper than string comparisons on a dotted-decimal `VARCHAR`, which is exactly why Section 6's deep dive matters beyond being an academic exercise. Log analysis and SIEM tools do the same conversion internally for the same reason. And when you see a raw hexadecimal netmask like `0xffffff00` in a config file or `ifconfig` output (Section 10), it's the same binary mask from Chapter 37 written in base 16 instead of base 10 — `ff` is `11111111`, exactly as derived in this chapter's binary tables.

---

## 15. Summary

| Term | Meaning |
|---|---|
| Flat address space | An address with no internal structure related to location (e.g., MAC addresses) |
| Hierarchical address space | An address split into parts, where part of it encodes topological location (e.g., IP addresses) |
| IPv4 | Internet Protocol version 4 — 32-bit hierarchical addressing, defined in RFC 791 |
| Octet | An 8-bit group; an IPv4 address has exactly four |
| Dotted-decimal notation | Writing a 32-bit IPv4 address as four decimal numbers (0–255) separated by dots |
| Network portion | The part of an IP address that identifies which network a host belongs to (Chapter 37) |
| Host portion | The part of an IP address that identifies a specific host within its network (Chapter 37) |
| Address space | The total set of possible addresses; IPv4's is 2^32 ≈ 4.3 billion |

You now know what an IPv4 address *is* and *why it's shaped the way it is* — but not yet exactly where the line falls between its network and host portions, or how a device uses that line to make a real decision. That's Chapter 37: network and host portions, subnet masks, and the binary AND operation a device performs on every single outgoing packet.
