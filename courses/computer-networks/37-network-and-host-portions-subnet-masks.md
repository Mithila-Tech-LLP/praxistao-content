# Chapter 37: Network and Host Portions, Subnet Masks

> **"An IP address by itself is an incomplete sentence. `192.168.1.10` means nothing until you also know where the network ends and the host begins — and that second piece of information is the subnet mask."**

---

## Table of Contents

1. [The Problem: An Address Alone Doesn't Say Where "Home" Ends](#1-the-problem-an-address-alone-doesnt-say-where-home-ends)
2. [A First (Wrong) Assumption: A Fixed Split](#2-a-first-wrong-assumption-a-fixed-split)
3. [The Real Solution: The Subnet Mask](#3-the-real-solution-the-subnet-mask)
4. [Binary Representation of a Mask](#4-binary-representation-of-a-mask)
5. [The AND Operation, Worked Step by Step](#5-the-and-operation-worked-step-by-step)
6. [The Actual Decision: Is This Destination Local, or Do I Need a Router?](#6-the-actual-decision-is-this-destination-local-or-do-i-need-a-router)
7. [Slash Notation, Previewed](#7-slash-notation-previewed)
8. [Three More Worked Examples, Increasing in Difficulty](#8-three-more-worked-examples-increasing-in-difficulty)
9. [Why Masks Must Be Contiguous Ones](#9-why-masks-must-be-contiguous-ones)
10. [A Hands-On Experiment](#10-a-hands-on-experiment)
11. [Common Misconceptions](#11-common-misconceptions)
12. [What's Simplified Here](#12-whats-simplified-here)
13. [Interview Questions & Model Answers](#13-interview-questions--model-answers)
14. [Exercises](#14-exercises)
15. [Summary](#15-summary)

---

## 1. The Problem: An Address Alone Doesn't Say Where "Home" Ends

Chapter 36 established that an IPv4 address is split into a network portion and a host portion, and that this split is what makes IP addressing hierarchical instead of flat like a MAC address. But look again at a bare address:

```
192.168.1.10
```

Where does the network portion end and the host portion begin? Is it `192` (network) and `168.1.10` (host)? Is it `192.168` (network) and `1.10` (host)? Is it `192.168.1` (network) and `10` (host)? The 32 bits alone don't say. Nothing about the number `192.168.1.10` inherently marks a boundary — the split is a piece of *configuration*, not something derivable from the address bits themselves.

And this boundary isn't a minor detail — it's load-bearing for one of the most frequent decisions every networked device makes, many times per second: **when I want to send a packet to some destination IP address, can I hand it directly to that device over my local network link, or do I need to send it to my router instead, because the destination is somewhere else entirely?**

Get this wrong, and either you flood your local network trying to reach devices that are actually far away, or you send everything to your router even when the destination is sitting three feet away on the same cable — both wrong, both wasteful. The device needs a reliable, mechanical way to answer this question for every single packet, without human intervention, in microseconds. That mechanism is the **subnet mask**.

---

## 2. A First (Wrong) Assumption: A Fixed Split

The most natural first guess is: "just always split on an octet boundary — the first octet is always the network, the rest is host." That would make `192.168.1.10`'s network portion `192` and its host portion `168.1.10`, giving 2^24 (about 16.7 million) possible hosts on every network.

This was, in fact, close to how the earliest version of IP addressing worked — the historical "classful" system, which Chapter 39 covers in full, including exactly why it was abandoned. The short version of why a *fixed* split fails: a company with 10 employees and a company with 10,000 employees cannot both be well served by the same fixed-size network. Give everyone a 16.7-million-address network and you waste almost all of it on the small company; force everyone into a network sized for 10,000 hosts and small deployments still waste enormous amounts of address space (this exact waste is the worked example at the heart of Chapter 39).

The fix is to make the boundary **configurable per network**, rather than fixed. Some networks might split after 24 bits, others after 26, others after 16 — whatever fits how many hosts that particular network actually needs. But if the boundary can move, something has to specify, for any given address, exactly where it currently sits. That "something" is the subnet mask.

---

## 3. The Real Solution: The Subnet Mask

A **subnet mask** is a second 32-bit number, always paired with an IP address, whose only job is to mark where the network portion ends and the host portion begins. It does this in the simplest possible way: a mask bit is `1` for every position that belongs to the network portion, and `0` for every position that belongs to the host portion.

```
 IP address:    192  .  168  .   1   .   10
 Subnet mask:   255  .  255  .  255  .   0

 mask meaning:  <------ network ------>|<host->
                (first 24 bits are      (last 8 bits
                 the network portion)    are the host
                                         portion)
```

Read left to right: wherever the mask has `1`s, the corresponding bits of the IP address are the network portion; wherever the mask has `0`s, the corresponding bits are the host portion. The mask `255.255.255.0` says "the first three octets (24 bits) are network, the last octet (8 bits) is host" — which is exactly the traditional "Class C"-style split, but now it's a stated configuration value rather than an assumption baked into the address itself.

Crucially, a mask is *not* an address. It never appears as a source or destination in a packet (Chapter 36, Section 7). It's local configuration, held by each device (and visible in the `ifconfig`/`ip addr` output you already ran in Chapter 36's experiment) and used purely to interpret addresses.

---

## 4. Binary Representation of a Mask

Applying Chapter 36's binary-to-decimal skills to a mask makes the "1s mark network, 0s mark host" rule completely literal:

```
 255.255.255.0  in binary:

  255      255      255       0
 11111111.11111111.11111111.00000000
 <-------- 24 ones -------->|<8 zeros>
```

Every `255` octet is eight `1` bits (255 = 128+64+32+16+8+4+2+1, every place value present). Every `0` octet is eight `0` bits. So `255.255.255.0` really is nothing more than "24 ones, then 8 zeros" — a direct, literal statement of "the first 24 bits are network, the rest are host."

A different mask, `255.255.0.0`, is 16 ones followed by 16 zeros — network portion is the first two octets, host portion is the last two, giving room for far more hosts per network (2^16 ≈ 65,536) but far fewer distinct networks from the same address block. A mask like `255.255.255.192` is less obviously round in decimal, but in binary it's just as simple:

```
 255.255.255.192 in binary:

  255      255      255      192
 11111111.11111111.11111111.11000000
 <---------- 26 ones ------->|<6 zeros>
```

192 = 128+64 = two leading `1` bits followed by six `0` bits. Read across all four octets, that's 26 ones total, then 6 zeros — a network portion of 26 bits, a host portion of 6 bits (room for 2^6 = 64 addresses per network, including the network and broadcast addresses — Chapter 38 covers exactly what those cost you).

This is the pattern to memorize by *deriving*, not by rote: **count the leading 1 bits in the mask, in binary, across all 32 bits — that count is the size of the network portion.**

---

## 5. The AND Operation, Worked Step by Step

Here is the actual mechanism, the one every operating system's networking stack runs on every outgoing packet: a bitwise **AND** between the IP address and the subnet mask.

The AND operation is simple per bit: the result is `1` only if *both* input bits are `1`; otherwise the result is `0`.

```
 AND truth table:
   0 AND 0 = 0
   0 AND 1 = 0
   1 AND 0 = 0
   1 AND 1 = 1
```

Applied bit by bit across an entire 32-bit address and mask, AND-ing with the mask has exactly the effect you'd want: every bit in the network portion (where the mask is `1`) passes through unchanged, and every bit in the host portion (where the mask is `0`) gets forced to `0`. The result is called the **network address** — the IP address with its host bits zeroed out, representing "this network, generically," with no specific host identified.

### Worked Example: Finding the Network Address

Take the IP address `192.168.1.10` with mask `255.255.255.0`.

```
        IP address:  11000000.10101000.00000001.00001010   (192.168.1.10)
        Subnet mask: 11111111.11111111.11111111.00000000   (255.255.255.0)
                     ------------------------------------- AND, bit by bit
     Network address: 11000000.10101000.00000001.00000000   (192.168.1.0)
```

Go through it one octet at a time to see exactly why:

- Octet 1: `11000000` AND `11111111` = `11000000` → 192 (mask is all 1s, so the IP's bits pass through unchanged)
- Octet 2: `10101000` AND `11111111` = `10101000` → 168 (same — unchanged)
- Octet 3: `00000001` AND `11111111` = `00000001` → 1 (same — unchanged)
- Octet 4: `00001010` AND `00000000` = `00000000` → 0 (mask is all 0s, so every IP bit is forced to 0, regardless of what it was)

Result: `192.168.1.0` — the network address for this device. Notice the host portion of the original address (`10`, meaning "host number 10 on this network") is completely gone; the network address represents the network as a whole, not any specific device on it.

---

## 6. The Actual Decision: Is This Destination Local, or Do I Need a Router?

Now the payoff. Every device configured with an IP address and a subnet mask can determine, for any destination it wants to reach, whether that destination is on the same local network (reachable directly, over the LAN, using MAC addressing as Chapter 35 traced end to end) or on a different network entirely (requiring the packet be sent to a router — previewed here, formalized in Chapter 44).

The algorithm is exactly the AND operation from Section 5, applied twice and compared:

```
 1. Compute MY network address:      (my IP)          AND (my mask)
 2. Compute DESTINATION's network:   (destination IP) AND (my mask)
 3. If the two results are EQUAL     → same network  → send directly (local delivery)
    If the two results DIFFER        → different network → send to my router (default gateway)
```

### Worked Example: Same Network

A device is configured with IP `192.168.1.10`, mask `255.255.255.0`. It wants to reach `192.168.1.200`.

```
 My network:
   11000000.10101000.00000001.00001010   (192.168.1.10)
 AND
   11111111.11111111.11111111.00000000   (255.255.255.0)
 = 11000000.10101000.00000001.00000000   (192.168.1.0)

 Destination's network (using MY mask, since that's the only mask I have):
   11000000.10101000.00000001.11001000   (192.168.1.200)
 AND
   11111111.11111111.11111111.00000000   (255.255.255.0)
 = 11000000.10101000.00000001.00000000   (192.168.1.0)

 Compare: 192.168.1.0  ==  192.168.1.0   → SAME NETWORK
```

Both results are `192.168.1.0`. The device concludes `192.168.1.200` is local, and — after resolving its MAC address via ARP (Chapter 53) — sends the frame directly, no router involved.

### Worked Example: Different Network

Same device (`192.168.1.10`, mask `255.255.255.0`) wants to reach `192.168.2.50` instead.

```
 My network (already computed above):     192.168.1.0

 Destination's network:
   11000000.10101000.00000010.00110010   (192.168.2.50)
 AND
   11111111.11111111.11111111.00000000   (255.255.255.0)
 = 11000000.10101000.00000010.00000000   (192.168.2.0)

 Compare: 192.168.1.0  !=  192.168.2.0   → DIFFERENT NETWORK
```

The results differ — `192.168.1.0` versus `192.168.2.0`. The device concludes `192.168.2.50` is remote and hands the packet to its configured **default gateway** (its router) instead of trying to deliver it directly. The router, which sits on multiple networks and has a fuller view of the topology (Chapters 44–45), takes over from there.

This is, quite literally, the first decision your laptop makes every single time an application tries to open a connection — before DNS, before TCP, before anything else, the operating system's IP stack runs exactly this AND-and-compare check to decide how to hand the packet to the network interface.

### The Algorithm as Code

Everything in this section is small enough to write as a handful of lines of real code — and doing so removes any doubt that "AND, then compare" really is the entire mechanism, with nothing hidden:

```go
package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

// toUint32 packs a 4-byte IPv4 address into one 32-bit integer,
// exactly as Chapter 36's deep dive described.
func toUint32(ip net.IP) uint32 {
	return binary.BigEndian.Uint32(ip.To4())
}

// sameNetwork implements Section 6's entire decision algorithm:
// AND my IP with my mask, AND the destination with my mask, compare.
func sameNetwork(myIP, destIP, mask net.IP) bool {
	m := toUint32(mask)
	return (toUint32(myIP) & m) == (toUint32(destIP) & m)
}

func main() {
	me := net.ParseIP("192.168.1.10")
	mask := net.ParseIP("255.255.255.0")

	dest1 := net.ParseIP("192.168.1.200")
	dest2 := net.ParseIP("192.168.2.50")

	fmt.Println(sameNetwork(me, dest1, mask)) // true  — local, deliver directly
	fmt.Println(sameNetwork(me, dest2, mask)) // false — remote, send to default gateway
}
```

This is not a simplified toy version of what the kernel does — it is, at the level that matters, the actual check. Real IP stacks add layers of caching (a route cache, so this comparison isn't recomputed from scratch for every packet to the same destination) and generalize it to handle multiple configured routes rather than just "local network vs. default gateway" (Chapter 45's longest prefix match is the generalized form of exactly this same AND-and-compare idea, run against every entry in a routing table instead of just one mask). But the core arithmetic — an AND, then an equality check — is identical.

---

## 7. Slash Notation, Previewed

Writing out a full dotted-decimal mask every time is verbose, and — as Section 9 will explain — largely redundant, since a valid mask is fully described by a single number: how many leading `1` bits it has. That number is written after the address, separated by a slash:

```
 192.168.1.10/24

 means: IP address 192.168.1.10, with a 24-bit network portion —
        exactly equivalent to writing the mask 255.255.255.0 separately.
```

This is called **CIDR notation** (Classless Inter-Domain Routing notation), and it is used throughout the rest of this course from here on, because it's compact and because, as Chapter 39 will show, prefix lengths aren't restricted to the traditional 8/16/24-bit boundaries — `/26`, `/23`, `/30` are all completely legal. This chapter has already used a `/26` mask (`255.255.255.192`) in Section 4 without needing to invent new machinery for it — the AND operation works identically no matter what the prefix length is. Chapter 39 is dedicated to CIDR notation and its history; for now, just recognize `/N` as shorthand for "a mask with N leading 1 bits."

---

## 8. Three More Worked Examples, Increasing in Difficulty

### Example A: A Non-Octet-Aligned Mask

Device: `10.4.20.130`, mask `255.255.255.192` (i.e., `/26`). Is `10.4.20.190` local?

```
 My network:
   00001010.00000100.00010100.10000010   (10.4.20.130)
 AND
   11111111.11111111.11111111.11000000   (255.255.255.192)
 = 00001010.00000100.00010100.10000000   (10.4.20.128)

 Destination's network:
   00001010.00000100.00010100.10111110   (10.4.20.190)
 AND
   11111111.11111111.11111111.11000000   (255.255.255.192)
 = 00001010.00000100.00010100.10000000   (10.4.20.128)

 Compare: 10.4.20.128 == 10.4.20.128   → SAME NETWORK, deliver directly
```

Only the last octet needed real bit-level work here (`10000010` AND `11000000` = `10000000`, and `10111110` AND `11000000` = `10000000`) — the first three octets are entirely covered by 1-bits in the mask and pass through unchanged, exactly as in Section 5's worked example.

### Example B: Same Third Octet, Different Subnet

Device: `10.4.20.130`, mask `255.255.255.192` (`/26`). Is `10.4.20.70` local? (Same third octet as Example A — a common source of real-world confusion.)

```
 My network (from Example A):        10.4.20.128

 Destination's network:
   00001010.00000100.00010100.01000110   (10.4.20.70)
 AND
   11111111.11111111.11111111.11000000   (255.255.255.192)
 = 00001010.00000100.00010100.01000000   (10.4.20.64)

 Compare: 10.4.20.128 != 10.4.20.64   → DIFFERENT NETWORK, send to router
```

This is the important, non-obvious case: `10.4.20.70` and `10.4.20.130` share the exact same first three octets and would look "clearly local" to anyone reasoning only in dotted-decimal. But with a `/26` mask, the third octet's boundary doesn't matter — the split happens *inside* the fourth octet, at bit position 6. `.70` (binary `01000110`) and `.130` (binary `10000010`) fall into different 64-address blocks (`.64`–`.127` versus `.128`–`.191`), and only the binary AND operation reveals that reliably. This is exactly the kind of case Chapter 38's subnetting problems are built around.

### Example C: The All-Ones Mask (A Single Host)

Device: `172.16.5.1`, mask `255.255.255.255` (`/32` — thirty-two 1-bits, zero host bits). Is `172.16.5.2` local?

```
 My network:
   10101100.00010000.00000101.00000001   (172.16.5.1)
 AND
   11111111.11111111.11111111.11111111   (255.255.255.255)
 = 10101100.00010000.00000101.00000001   (172.16.5.1)   -- unchanged; every bit is "network"

 Destination's network:
   10101100.00010000.00000101.00000010   (172.16.5.2)
 AND
   11111111.11111111.11111111.11111111   (255.255.255.255)
 = 10101100.00010000.00000101.00000010   (172.16.5.2)   -- unchanged

 Compare: 172.16.5.1 != 172.16.5.2   → DIFFERENT NETWORK, even though only 1 apart!
```

A `/32` mask has zero host bits — the "network address" for any address under this mask is the address itself. With this mask, no other address, not even one numerically adjacent, is ever considered "local." This isn't a mistake or an edge case to work around — `/32` routes are genuinely used in real networks to refer to one single, specific host (a loopback address, or a server that needs an individually routable address), and Chapter 45's longest-prefix-match algorithm relies on `/32` being a completely ordinary, valid prefix length, just an extreme one.

---

## 9. Why Masks Must Be Contiguous Ones

Nothing about the AND mechanism in Section 5 technically *requires* the 1-bits in a mask to be contiguous from the left. A mask like `11111111.11111111.11110101.00000000` (some scattered 1s and 0s in the third octet) would still "work" as input to an AND operation. Early networking standards even technically permitted such masks.

In practice, discontiguous masks were never meaningfully used, and modern equipment generally rejects them outright, for a good structural reason: the entire value of the network/host split is that the network portion behaves as one contiguous, summarizable *block* of addresses (Chapter 39's aggregation, and Chapter 50's route summarization, depend on this completely). A contiguous run of `1` bits at the top of the mask carves out a clean, describable range — "every address that starts with these N bits." A scattered mask would carve out a scattered, non-contiguous, hard-to-reason-about set of addresses that couldn't be summarized as a single routing table entry at all, defeating the entire hierarchical-addressing point of Chapter 36.

So: by convention (and, in current standards, by requirement), a valid subnet mask is always some number of `1` bits, left-aligned, followed by `0` bits filling the rest — which is exactly why a single number (the CIDR prefix length from Section 7) is sufficient to fully describe any valid mask.

---

## 10. A Hands-On Experiment

Find your own device's mask and manually compute your network address, then confirm it against a real tool.

```bash
# macOS
$ ifconfig en0 | grep inet
        inet 192.168.1.42 netmask 0xffffff00 broadcast 192.168.1.255

# Linux
$ ip addr show eth0 | grep inet
    inet 192.168.1.42/24 brd 192.168.1.255 scope global eth0
```

`0xffffff00` is the mask written in hexadecimal — convert it: `ff` = `11111111` (three times), `00` = `00000000` — the same `255.255.255.0` from Section 3, just in a different base. `/24` on the Linux line says the same thing directly.

Now compute the network address by hand (AND your IP with your mask, as in Section 5) and check it against a calculator:

```bash
$ ipcalc 192.168.1.42/24
Address:   192.168.1.42
Netmask:   255.255.255.0 = 24
Network:   192.168.1.0/24
Broadcast: 192.168.1.255
HostMin:   192.168.1.1
HostMax:   192.168.1.254
```

(`ipcalc` may need installing — `sudo apt install ipcalc` on Debian/Ubuntu, `brew install ipcalc` on macOS.) Your hand-computed network address should match the tool's `Network:` line exactly — if it doesn't, redo the AND operation octet by octet; the bug is almost always a decimal-to-binary slip in one octet.

You can also watch the Linux kernel make Section 6's exact local-vs-remote decision live, using `ip route get`, which reports which interface and next hop it would actually use for a given destination:

```bash
$ ip route get 192.168.1.200
192.168.1.200 dev eth0 src 192.168.1.42 uid 1000
    cache

$ ip route get 8.8.8.8
8.8.8.8 via 192.168.1.1 dev eth0 src 192.168.1.42 uid 1000
    cache
```

The first command's output has no `via` — the kernel determined `192.168.1.200` is on the same network as `192.168.1.42/24` (both AND down to `192.168.1.0`, exactly as in Section 6's worked example) and would send the frame directly. The second command's output shows `via 192.168.1.1` — the kernel determined `8.8.8.8` is remote and would forward the packet to that address, the default gateway, instead. This is the AND-and-compare algorithm, running for real, on your own machine, printing its own conclusion.

---

## 11. Common Misconceptions

- **"The subnet mask is part of the IP address."** No — it's separate configuration, always paired with an address but never transmitted as part of a packet's source or destination address fields (Chapter 36, Section 7). Two devices on the same LAN must agree on the mask by configuration, not by reading it off each other's packets.
- **"You can tell the network/host split just by looking at the dotted-decimal address."** Only true for the three traditional class boundaries (/8, /16, /24) that classful addressing (Chapter 39) trained people to assume. With CIDR, the split can fall anywhere, even mid-octet, as Example B in Section 8 showed dramatically — you cannot know the split without being told the mask.
- **"If two addresses look close together (like `.70` and `.130`), they must be on the same network."** Disproven directly in Section 8, Example B. Numeric closeness in decimal has no reliable relationship to being on the same network; only the AND operation, using the actual configured mask, answers that question correctly.
- **"A /32 mask is a mistake — every network needs host bits."** Not a mistake — a /32 legitimately describes a route to one single specific address, and is common in real routing tables and firewall rules (Chapter 45).

---

## 12. What's Simplified Here

This chapter deliberately did not derive *how* to choose or design a mask when splitting an address block into smaller pieces — that is the real problem Chapter 38 tackles from first principles, worked by hand across several progressively harder scenarios. It also did not cover the historical classful system in depth (default masks tied to address ranges) — that's Chapter 39, alongside CIDR's fix for the waste that system caused. What this chapter did establish, thoroughly, is the *mechanism*: what a mask is, how AND works bit by bit, and how a device uses both to make the single most common decision in all of IP networking.

---

## 13. Interview Questions & Model Answers

**Beginner: "What is a subnet mask, and why is it needed?"**

A subnet mask is a 32-bit value, paired with an IP address, that marks which bits of the address are the network portion (mask bit = 1) and which are the host portion (mask bit = 0). It's needed because the network/host split in an IP address isn't fixed or self-evident from the address alone (Chapter 36) — the mask is the configuration that states where the boundary actually is for a given address.

**Beginner: "How do you determine a device's network address from its IP and subnet mask?"**

Perform a bitwise AND between the IP address and the subnet mask, bit by bit. Every bit where the mask is 1 passes through unchanged; every bit where the mask is 0 becomes 0. The result is the network address — the original address with all its host bits zeroed out.

**Intermediate: "How does a device decide whether to send a packet directly on the LAN or to its default gateway?"**

It computes its own network address (its IP AND its mask) and the destination's network address (the destination IP AND its own mask, since it has no other mask to use). If the two results match, the destination is on the same network and can be reached directly (after ARP resolves the MAC address, Chapter 53). If they differ, the destination is remote, and the packet is sent to the configured default gateway instead.

**Intermediate: "Why is a mask like `255.255.255.213` invalid, even though nothing stops you from typing it?"**

A valid mask must be a contiguous run of 1-bits followed by a contiguous run of 0-bits (Section 9). `213` in binary is `11010101` — the 1s and 0s are interleaved, not contiguous. Such a mask can't be described by a single CIDR prefix length and would carve out a scattered, non-summarizable set of addresses, defeating the purpose of hierarchical addressing established in Chapter 36.

**Advanced: "Two hosts, `10.0.0.70/26` and `10.0.0.130/26`, differ by only 60 in their last octet. Are they on the same subnet? Show your work."**

No. `/26` means 26 network bits, leaving 6 host bits — subnet blocks of size 2^6 = 64 in the last octet, at boundaries 0, 64, 128, 192. `.70` in binary is `01000110`; AND with the mask's last-octet byte `11000000` gives `01000000` = 64, so its network is `10.0.0.64/26`. `.130` in binary is `10000010`; AND with `11000000` gives `10000000` = 128, so its network is `10.0.0.128/26`. The networks differ, so the hosts are on different subnets despite looking numerically close — a device would send traffic between them via its default gateway, not directly on the LAN.

**Advanced: "Explain why the CIDR prefix length alone (e.g., `/26`) is sufficient to fully reconstruct the dotted-decimal subnet mask, with no ambiguity."**

Because a valid mask, by convention and modern requirement, is always a contiguous run of 1-bits from the most significant bit, followed by 0-bits for the remainder (Section 9). Given only the count of 1-bits (the prefix length N), you can reconstruct the mask uniquely: set the first N bit positions (out of 32) to 1, the rest to 0, then regroup into four octets and convert each to decimal. There is exactly one mask for any given valid N, which is why CIDR notation lost no information by dropping the full dotted-decimal form.

---

## 14. Exercises

### Easy

1. Convert the mask `255.255.240.0` to binary and state how many network bits and how many host bits it represents.
2. Given IP `172.16.10.5` and mask `255.255.255.0`, compute the network address by hand, showing the AND operation for each octet.
3. Write the CIDR (slash) notation equivalent for the masks `255.0.0.0`, `255.255.0.0`, and `255.255.255.0`.

### Medium

4. A host has IP `192.168.5.60` and mask `255.255.255.240` (`/28`). Determine whether `192.168.5.75` is on the same network, showing the full binary AND for the last octet of both addresses.
5. A host has IP `10.10.10.10/27`. List, in binary and decimal, the mask, the network address, and the range of usable host addresses you'd expect this network to contain (don't worry yet about excluding the network/broadcast address precisely — Chapter 38 formalizes that; just identify the block boundaries).
6. Explain why the mask `255.255.255.255` (`/32`) results in "every other address is remote," using the AND mechanism to justify it rather than just stating the rule.

### Hard

7. Two hosts are configured as `192.168.100.14/29` and `192.168.100.9/29`. Are they on the same network? Show every step of the binary AND for both, including deriving the mask's binary form from `/29` first.
8. A device is misconfigured with the *wrong* subnet mask: its real network uses `/25`, but it's configured with `/24`. Using the AND-and-compare decision process from Section 6, describe a concrete pair of addresses (one truly local, one truly remote) for which this misconfiguration would cause the device to reach the wrong conclusion, and explain in terms of the binary AND exactly why.

---

## Production Usage Notes

A mismatched subnet mask is one of the most common real-world causes of "I can ping some things on my network but not others" tickets — exercise 8 above is not a contrived scenario, it happens whenever a device is manually configured (or a DHCP scope is misconfigured, Chapter 55) with a mask that doesn't match the rest of its actual subnet. The symptom is asymmetric and confusing precisely because the AND-and-compare check (Section 6) is purely local to each device: a misconfigured host might wrongly decide a truly-local peer is remote (sending traffic needlessly through the gateway, which may or may not route it back correctly) or wrongly decide a truly-remote host is local (and then get no ARP reply, Chapter 53, because no such device exists on the LAN) — and the *other*, correctly-configured devices on the same network see none of this, because their own mask is fine. This is why "check the subnet mask" is one of the first steps in any real Layer 3 connectivity troubleshooting playbook (Chapter 122 covers this playbook systematically), well before considering routing protocols, firewalls, or DNS.

---

## 15. Summary

| Term | Meaning |
|---|---|
| Subnet mask | A 32-bit value marking which bits of an IP address are network (1) vs. host (0) |
| Network address | An IP address with all host bits set to 0; represents a network as a whole |
| Bitwise AND | An operation where the result bit is 1 only if both input bits are 1; used to apply a mask to an address |
| Default gateway | The router a device sends traffic to when the destination is not on its local network |
| CIDR / slash notation | Shorthand `/N` for a mask with N leading 1-bits (previewed here, fully covered in Chapter 39) |
| Contiguous mask | A valid mask: a run of 1-bits followed by a run of 0-bits, with no interleaving |
| /32 | A mask with zero host bits; describes exactly one specific address |

You now know how a device decides "local or remote" for any single destination — but not yet how an organization, handed one address block, deliberately *chooses* where to put that boundary in order to carve the block into several independently manageable networks. That real-world problem, solved by hand through several increasingly difficult worked examples, is Chapter 38: Subnetting From First Principles.
