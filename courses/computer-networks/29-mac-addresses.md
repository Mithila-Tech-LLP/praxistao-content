# Chapter 29: MAC Addresses — Every Device's Physical Name

> *"An IP address is more like a mailing address — it can change depending on where you're standing. A MAC address is more like a fingerprint — it's supposed to be etched into the device itself, before it's ever plugged into any network at all."*

---

## Table of Contents

1. [The Problem: Ethernet Needs an Address, But Which Kind?](#1-the-problem-ethernet-needs-an-address-but-which-kind)
2. [A Naive Attempt: Just Use the IP Address](#2-a-naive-attempt-just-use-the-ip-address)
3. [The Real Solution: A 48-Bit Physical Address](#3-the-real-solution-a-48-bit-physical-address)
4. [Anatomy of a MAC Address](#4-anatomy-of-a-mac-address)
5. [Special Bits: Unicast/Multicast and Universal/Local](#5-special-bits-unicastmulticast-and-universallocal)
6. [Special Addresses: Broadcast, Multicast, All-Zero](#6-special-addresses-broadcast-multicast-all-zero)
7. [MAC vs. IP: A Direct Comparison](#7-mac-vs-ip-a-direct-comparison)
8. [Deep Dive: Are MAC Addresses Really Globally Unique?](#8-deep-dive-are-mac-addresses-really-globally-unique)
9. [MAC Randomization and Privacy](#9-mac-randomization-and-privacy)
10. [A Real Example: Reading Your Own MAC Address](#10-a-real-example-reading-your-own-mac-address)
11. [Code: Validating and Decoding a MAC Address](#11-code-validating-and-decoding-a-mac-address)
12. [Hands-On Experiment](#12-hands-on-experiment)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes](#14-production-notes)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)

---

## 1. The Problem: Ethernet Needs an Address, But Which Kind?

Chapter 28 showed you the exact byte layout of an Ethernet frame: a 6-byte destination address field and a 6-byte source address field, sitting right after the preamble. But it deliberately left one question open: what actually *goes* in those fields? What kind of address is small enough to fit in 6 bytes, unique enough to identify one specific device out of billions, and — critically — available the instant a device is powered on, before any network configuration has happened?

That last requirement is the crux of the problem, and it's worth sitting with. Think about the very first frame a laptop ever sends after you plug in an Ethernet cable, before it has requested an IP address from anything (that process, DHCP, is Chapter 55 — and DHCP itself relies on Ethernet frames to even get its own request out). At that moment, the device has no IP address yet. So whatever address identifies it at the Ethernet layer cannot depend on IP addressing existing first. It needs to be something that comes baked into the hardware itself.

## 2. A Naive Attempt: Just Use the IP Address

If you already know that IP addresses (fully covered starting Chapter 36) identify computers on a network, the obvious question is: why not just put the IP address in the Ethernet destination field and skip having two different addressing schemes entirely?

This fails for several concrete reasons:

- **Chicken-and-egg problem.** As noted above, a device frequently doesn't have an IP address yet when it needs to send its very first frame — including the DHCP request that will eventually *give* it an IP address. You cannot bootstrap IP configuration using an addressing scheme that depends on IP configuration already being done.
- **Ethernet predates and outlives any one network-layer protocol.** Ethernet frames don't just carry IPv4 — they carry IPv6, ARP messages, and historically other network-layer protocols entirely (IPX, AppleTalk). If the Ethernet address field were hard-wired to one network-layer addressing format, Ethernet would have to be redesigned every time a new network-layer protocol appeared. Layering (Chapter 24) exists precisely to prevent this kind of coupling.
- **IP addresses are logical and reassignable; hardware addressing needs to be intrinsic.** An IP address describes *where* something currently sits in a network's topology — Chapter 36 will show that IP addresses are structured specifically so that routers can make decisions based on network location. A NIC, in contrast, needs an identity that has nothing to do with which network it's currently plugged into, precisely so it can be manufactured, boxed, shipped, and sold without any network in mind at all.
- **Address-family independence.** Not every device speaking Ethernet even runs IP — some run only legacy or specialized protocols. The Ethernet frame format has to work regardless.

So the real solution needs an addressing scheme that lives entirely at the hardware/link layer, independent of whatever logical addressing scheme sits above it.

## 3. The Real Solution: A 48-Bit Physical Address

The answer is the **MAC address** — Media Access Control address, a 48-bit (6-byte) number assigned to a network interface at the time it's manufactured, stored in read-only or flash memory on the NIC itself. Because it's burned in at the factory rather than configured by an administrator or protocol, it's commonly called a **burned-in address (BIA)**, **hardware address**, or **physical address** — all three terms mean the same thing and you'll see all three in real documentation and tools.

**Intuitive level:** every network interface card ships from the factory with a unique serial number etched into it, the way every car engine has a stamped serial number regardless of what license plate (its "logical," reassignable address) gets bolted onto the car later. The engine number doesn't change when you sell the car to someone in a different city; the license plate does.

**Engineering level:** the MAC address is a flat (non-hierarchical) 48-bit identifier, structured into two 24-bit halves — a manufacturer identifier and a device-specific serial — that together are intended to be globally unique across every Ethernet (and Wi-Fi — 802.11 devices use the same MAC address format, previewed for Volume 13) device ever made.

**Deep technical level:** Section 4 breaks down exactly how those 48 bits are structured, and Section 8 examines just how strong the "globally unique" guarantee really is in practice.

## 4. Anatomy of a MAC Address

A MAC address is written as six groups of two hexadecimal digits, most commonly separated by colons (Linux/Unix convention) or hyphens (Windows convention):

```
3c:22:fb:12:34:56          (colon notation)
3C-22-FB-12-34-56          (hyphen notation — same address)

Byte:     1        2        3        4        5        6
Bits:  00000000 00000000 00000000 00000000 00000000 00000000
       └──────── OUI (24 bits) ───────────┘└── NIC-specific (24 bits) ──┘
       assigned by IEEE to a manufacturer   assigned by the manufacturer
```

The 48 bits split into two structurally distinct 24-bit halves:

| Bytes | Name | Assigned by | Purpose |
|---|---|---|---|
| First 3 bytes (bytes 1–3) | **OUI** (Organizationally Unique Identifier) | IEEE Registration Authority | Identifies the manufacturer (e.g., Apple, Intel, Cisco) |
| Last 3 bytes (bytes 4–6) | NIC-specific identifier (sometimes called the "extension identifier") | The manufacturer itself | A serial number the manufacturer must not reuse within that OUI |

A manufacturer applies to IEEE and pays a fee for a block of OUIs. Each OUI grants that manufacturer 2^24 (about 16.7 million) possible addresses to assign as it manufactures NICs. For example, `3c:22:fb` is a real OUI registered to Apple, Inc. — meaning every device with a MAC address starting `3c:22:fb` has an Apple-made network interface. You can look up any OUI in IEEE's public registry to identify a device's manufacturer just from its MAC address — a real forensic and network-inventory technique used constantly in practice.

## 5. Special Bits: Unicast/Multicast and Universal/Local

Two individual bits within the very first byte of a MAC address carry special meaning, and they're worth knowing precisely because they explain behavior you'll see in real tools.

```
First byte of a MAC address, bit numbering (bit 0 = least significant bit):

  bit:  7  6  5  4  3  2  1  0
        .  .  .  .  .  .  .  I/G   <- bit 0: Individual/Group bit
        .  .  .  .  .  .  U/L .    <- bit 1: Universal/Local bit
```

- **I/G bit (bit 0 of the first byte):** 0 means this address identifies a single device (**unicast**); 1 means it identifies a group of devices (**multicast**) that should all receive the frame. The broadcast address (Section 6) is a special case of multicast where every bit is 1.
- **U/L bit (bit 1 of the first byte):** 0 means this is a **universally administered address** — assigned by the manufacturer from its IEEE-issued OUI, intended to be globally unique. 1 means this is a **locally administered address** — one that has been manually or programmatically overridden by software or an administrator, breaking the link to the original OUI. This bit is exactly how you can tell, just by looking at a MAC address, whether it's the "real" factory address or one that's been changed (Section 9 covers a major real-world reason this happens constantly today).

A quick way to check: if the second hex digit of a MAC address is `2`, `6`, `a`, or `e` (binary `x010`, `x110`, `x010`, `x110` — i.e., bit 1 set, bit 0 clear), it's a locally-administered unicast address. Apple's real OUI `3c:22:fb` has second digit `c` (binary `1100`) — bit 1 is 0, confirming it's a universally administered address, consistent with being a genuine factory-assigned OUI.

## 6. Special Addresses: Broadcast, Multicast, All-Zero

A handful of MAC address values are reserved for specific meanings rather than identifying one physical NIC:

| Address | Meaning |
|---|---|
| `ff:ff:ff:ff:ff:ff` | **Broadcast** — every bit set to 1; delivered to every device on the local network segment. Used by ARP requests (Chapter 53) and DHCP discovery (Chapter 55), among others. |
| `01:00:5e:xx:xx:xx` | **IPv4 multicast** range — the I/G bit is set, and this specific OUI-like prefix is reserved by IEEE for mapping IPv4 multicast group addresses onto Ethernet, previewed here and covered properly alongside IP multicast in Chapter 40. |
| `33:33:xx:xx:xx:xx` | **IPv6 multicast** mapping, used by Neighbor Discovery Protocol (Chapter 43) among other IPv6 mechanisms. |
| `00:00:00:00:00:00` | All-zero — never a valid assigned address; sometimes seen as an uninitialized or placeholder value in software, and treated as invalid by most stacks. |

The broadcast address matters enormously for the switch-forwarding behavior you'll see in Chapter 31: a switch must forward any frame addressed to `ff:ff:ff:ff:ff:ff` out every port (except the one it arrived on), because the destination field itself says "everyone on this segment gets a copy."

## 7. MAC vs. IP: A Direct Comparison

Because the two addressing systems are so often mentioned in the same breath, it's worth being explicit and precise about how they differ — this table previews concepts (like hierarchy) that Chapter 36 develops in full, and is meant purely as a contrast, not as a substitute for that later chapter:

| Property | MAC Address | IP Address |
|---|---|---|
| Layer | Data Link (Layer 2) | Network (Layer 3) |
| Size | 48 bits | 32 bits (IPv4) or 128 bits (IPv6) |
| Structure | Flat — OUI + serial, no topological meaning | Hierarchical — network portion + host portion (Chapter 37) |
| Assigned by | Manufacturer (burned in), or overridden locally | Network administrator, DHCP server, or ISP |
| Changes when? | Rarely — tied to the physical NIC | Frequently — changes when a device moves to a different network |
| Scope | Only meaningful within one local network segment | Globally routable (public) or scoped to a private network |
| Used for | Delivering a frame to the next physical hop | Identifying a device's location across the entire Internet |
| Analogy | A device's serial number / fingerprint | A mailing address, which changes when you move |

The single most important idea in this table, worth restating plainly: a MAC address answers "which physical device, right here on this wire, should receive this frame?" An IP address answers a completely different, much harder question — "how do I find this device anywhere on a planet-spanning network of networks?" — which is exactly why the Internet needs both a flat, hardware-level addressing scheme for the last physical hop and a hierarchical, routable addressing scheme for everything before that. You'll see this division of labor made completely concrete in Chapter 53, when ARP is introduced specifically to translate between the two.

## 8. Deep Dive: Are MAC Addresses Really Globally Unique?

The IEEE registry design (a manufacturer gets an OUI, and is trusted to never reuse a NIC-specific suffix within it) is *intended* to guarantee global uniqueness, and for the vast majority of practical purposes, it does. But it's worth being honest about the ways this guarantee weakens in practice:

- **Manufacturing errors and firmware bugs** have, on real occasions, produced duplicate MAC addresses shipped in consumer hardware — rare, but documented.
- **Virtual machines** are commonly assigned MAC addresses from a range reserved for virtualization (e.g., VMware's `00:50:56` OUI, or ranges hypervisors generate pseudo-randomly), and cloning a VM image carelessly can duplicate a MAC address across multiple running instances on the same network — a genuinely common real-world networking bug.
- **Local administration** (Section 5's U/L bit) explicitly allows software to assign any MAC address it wants, overriding the factory value entirely — used deliberately for things like consistent MACs across a failover cluster (so a backup server can "become" the primary by adopting its MAC), and, as Section 9 covers, for privacy.
- **Duplicate MAC addresses on the same local segment cause real, hard-to-diagnose problems** — because switches (Chapter 31) learn one port per MAC address, two devices claiming the same MAC will cause the switch's forwarding table to flap between ports, and traffic intended for one device may intermittently arrive at the other.

The practical takeaway: MAC addresses only need to be unique *within a given broadcast domain* (Chapter 30) to function correctly — global uniqueness is a strong convention that makes this reliably true almost everywhere, not a cryptographically or physically enforced guarantee.

## 9. MAC Randomization and Privacy

Here's a very real, current-day reason locally administered addresses matter to you personally: modern phones and laptops (iOS, Android, Windows, macOS) now generate a **random, locally administered MAC address** for every new Wi-Fi network they join, instead of broadcasting their real factory MAC address.

Why? Because a stable, unique MAC address is a powerful tracking identifier. Before MAC randomization became standard (roughly 2014 onward across major platforms), retail stores, airports, and advertising networks could — and did — track a specific device (and by extension, a specific person) moving through physical space over weeks or months, just by passively observing which MAC address kept showing up on their Wi-Fi access points' probe requests, with no cooperation from the user needed at all. Randomizing the MAC address per network (and in more aggressive implementations, periodically even on the same network) breaks that long-term tracking capability, at the cost of losing the convenience of a permanently identifiable device.

You can directly observe the U/L bit doing its job here: a randomized MAC address will always have bit 1 of its first byte set to 1, distinguishing it at a glance from a genuine factory-assigned address.

## 10. A Real Example: Reading Your Own MAC Address

```bash
# Linux
ip link show

# Example output (trimmed):
# 2: eth0: <BROADCAST,MULTICAST,UP,LOWER_UP> mtu 1500 ...
#     link/ether 3c:22:fb:12:34:56 brd ff:ff:ff:ff:ff:ff

# macOS
ifconfig en0 | grep ether
# ether 3c:22:fb:12:34:56

# Windows
ipconfig /all
# Physical Address. . . . . . . . . : 3C-22-FB-12-34-56
```

Every one of these outputs is showing you exactly the structure from Section 4: three bytes of OUI (here, `3c:22:fb`, a real Apple OUI, used purely as an example), three bytes of manufacturer-assigned serial. Note also that `ip link show` explicitly prints `brd ff:ff:ff:ff:ff:ff` — the broadcast address from Section 6, listed as a property of the interface because every interface must recognize and accept frames sent to it.

## 11. Code: Validating and Decoding a MAC Address

```go
package main

import (
	"fmt"
	"net"
)

// A tiny, illustrative OUI table — real tools use IEEE's full public registry.
var ouiTable = map[string]string{
	"3c:22:fb": "Apple, Inc.",
	"00:50:56": "VMware, Inc.",
	"b8:27:eb": "Raspberry Pi Foundation",
}

func describeMAC(macStr string) {
	mac, err := net.ParseMAC(macStr)
	if err != nil {
		fmt.Printf("%s: invalid MAC address (%v)\n", macStr, err)
		return
	}

	firstByte := mac[0]
	isMulticast := firstByte&0x01 != 0  // I/G bit
	isLocal := firstByte&0x02 != 0      // U/L bit
	oui := fmt.Sprintf("%02x:%02x:%02x", mac[0], mac[1], mac[2])

	fmt.Printf("MAC:              %s\n", mac)
	fmt.Printf("  Unicast/Multicast: %s\n", ternary(isMulticast, "multicast", "unicast"))
	fmt.Printf("  Universal/Local:   %s\n", ternary(isLocal, "locally administered", "universally administered (factory)"))
	if vendor, ok := ouiTable[oui]; ok && !isLocal {
		fmt.Printf("  Vendor (by OUI):   %s\n", vendor)
	}
	fmt.Println()
}

func ternary(cond bool, a, b string) string {
	if cond {
		return a
	}
	return b
}

func main() {
	describeMAC("3c:22:fb:12:34:56") // real-looking Apple OUI, unicast
	describeMAC("ff:ff:ff:ff:ff:ff") // broadcast
	describeMAC("02:00:00:00:00:01") // locally administered, unicast (randomized-style)
	describeMAC("01:00:5e:00:00:01") // IPv4 multicast mapping
}
```

Output:

```
MAC:              3c:22:fb:12:34:56
  Unicast/Multicast: unicast
  Universal/Local:   universally administered (factory)
  Vendor (by OUI):   Apple, Inc.

MAC:              ff:ff:ff:ff:ff:ff
  Unicast/Multicast: multicast
  Universal/Local:   locally administered

MAC:              02:00:00:00:00:01
  Unicast/Multicast: unicast
  Universal/Local:   locally administered

MAC:              01:00:5e:00:00:01
  Unicast/Multicast: multicast
  Universal/Local:   universally administered (factory)
```

## 12. Hands-On Experiment

```bash
# 1. Find your machine's real MAC address:
ip link show     # Linux
ifconfig         # macOS

# 2. Look up its OUI (first 3 bytes) against IEEE's public registry:
#    https://standards-oui.ieee.org/oui/oui.txt
#    Search for the first 6 hex digits (no separators) — you'll find
#    the exact company that manufactured your NIC.

# 3. If you're on Wi-Fi, join a network, then check your MAC address
#    again with the same command. On a modern phone or laptop, compare
#    it to the MAC shown on a *different* Wi-Fi network you've joined —
#    on iOS/Android/modern Windows/macOS, these will very likely differ,
#    demonstrating per-network MAC randomization from Section 9.

# 4. Check the U/L bit yourself: take the second hex digit of the MAC.
#    If it's 2, 6, a, or e, bit 1 is set -> locally administered
#    (i.e., randomized, not the factory address).
```

## 13. Common Misconceptions

- **"MAC addresses are globally routable, like IP addresses."** They're not routable at all beyond the local network segment — a MAC address has no hierarchy a router could use to decide "which direction is this device in." That's precisely why IP addressing (Chapter 36) exists as a separate, hierarchical scheme.
- **"A MAC address can never be changed."** Software absolutely can override it (the U/L bit exists for exactly this), and operating systems now do this constantly and deliberately for privacy (Section 9). The *factory-assigned* value is fixed; what the OS actually presents to the network is not.
- **"Two devices can never accidentally have the same MAC address."** They can, and it happens — cloned VM images and rare manufacturing defects are the most common real-world causes (Section 8).
- **"MAC addresses identify a computer."** More precisely, they identify a *network interface*. A laptop with both Wi-Fi and Ethernet has two different MAC addresses, one per interface.

## 14. Production Notes

- Network inventory and security tools routinely use OUI lookups to fingerprint unknown devices on a network (e.g., "this MAC's OUI belongs to a known IoT camera vendor").
- MAC address filtering as a Wi-Fi security measure is weak and largely obsolete for stopping determined attackers, since MAC addresses are trivially spoofable in software — real Wi-Fi security relies on WPA2/WPA3 (Chapter 89), not MAC filtering.
- Enterprise network access control (802.1X, NAC systems) often uses MAC addresses as one signal among several, precisely because they're easy to spoof and shouldn't be trusted alone.
- MAC randomization has real operational side effects: enterprise Wi-Fi networks that historically whitelisted devices by MAC address had to adapt when clients started rotating their MACs per network.

## 15. What's Simplified Here

This chapter treats the IEEE OUI registry as authoritative and stable; in reality IEEE has also introduced smaller allocation blocks (CID, MA-M, MA-S) alongside the traditional 24-bit OUI to stretch a shrinking address-block supply, which changes the exact bit-split in some newer allocations but doesn't change anything about how MAC addresses function on the wire. The mapping from IPv4/IPv6 multicast addresses to their reserved Ethernet multicast ranges (`01:00:5e:...` and `33:33:...`) is shown here only as a preview; the actual mapping algorithm is covered together with IP multicast in Chapter 40.

## 16. Interview Questions & Model Answers

**Beginner: What is a MAC address, and why is it called a "physical" address?**
"It's a 48-bit address burned into a network interface at manufacture time, used to identify a device at the Ethernet (data link) layer. It's called a physical or hardware address because, unlike an IP address, it's tied to the actual piece of hardware rather than assigned by network configuration — it's meant to work the same no matter what network the device is plugged into."

**Intermediate: What's the structural difference between a MAC address and an IPv4 address, and why does that difference matter?**
"A MAC address is flat — 48 bits with no internal hierarchy beyond identifying a manufacturer versus a specific device. An IPv4 address is hierarchical — split into a network portion and a host portion, which lets routers make forwarding decisions based on which network a destination belongs to without knowing about every individual host. A flat address space like MAC doesn't scale to global routing because there's no way to summarize 'these million addresses are all roughly in this direction' — you'd need a lookup entry per device. That's exactly why the Internet layers a hierarchical address (IP) on top of a flat one (MAC)."

**Advanced: How does your phone's MAC address randomization interact with how switches and access points build their forwarding/association state, and why doesn't it break normal connectivity?**
"MAC randomization typically generates a new locally-administered address per SSID (and sometimes per connection), but for the duration of any single association to a given network, the device uses one consistent MAC address. Switches (Chapter 31) and access points don't require a MAC address to be the factory-assigned one — they just need it to be internally consistent for the duration of that session so their MAC learning and association tables stay coherent. It only becomes a problem for systems that specifically rely on a device's MAC staying the same across different networks or over long periods, like MAC-based access control lists, which is a large part of why enterprise Wi-Fi has had to move toward 802.1X-based authentication instead."

## 17. Exercises

### Easy
1. Split the MAC address `b8:27:eb:a1:02:c3` into its OUI and device-specific portions.
2. Is `ff:ff:ff:ff:ff:ff` a unicast or multicast address, and how can you tell just from the bit pattern?
3. Name two reasons a device cannot simply use its IP address as its Ethernet-layer address.

### Medium
4. Given the MAC address `02:1a:2b:3c:4d:5e`, determine whether it is universally or locally administered, and explain what real-world scenario would likely produce a MAC address like this.
5. Explain why duplicate MAC addresses on the same LAN cause problems, but duplicate MAC addresses on two completely separate, unconnected LANs cause none.
6. A security researcher wants to identify what kind of device (e.g., "a Raspberry Pi") is connected to a network, using only its MAC address and no other information. Explain exactly what they can and cannot reliably learn this way.

### Hard
7. Extend the Go program from Section 11 to accept a MAC address from the command line and print whether it looks like a MAC-randomization-style address (locally administered, unicast) versus a plausible factory address, and justify what confidence level that conclusion deserves given Section 8's discussion of the uniqueness guarantee.
8. Design (in prose, no code required) a scheme a switch vendor could use to detect and alert on a duplicate-MAC-address situation on a live network, using only information available to individual switches from the MAC learning process described conceptually in this chapter (full algorithm comes in Chapter 31).

## 18. Summary

| Term | Meaning |
|---|---|
| MAC address | 48-bit hardware address identifying a network interface |
| OUI | First 24 bits; IEEE-assigned identifier for the manufacturer |
| Burned-in address (BIA) | The factory-assigned MAC address stored in NIC hardware |
| I/G bit | Bit 0 of the first byte; 0 = unicast, 1 = multicast |
| U/L bit | Bit 1 of the first byte; 0 = universally administered (factory), 1 = locally administered (overridden) |
| Broadcast address | `ff:ff:ff:ff:ff:ff` — delivered to every device on the segment |
| MAC randomization | OS-level privacy feature generating a locally administered MAC per network |

You now know what goes in an Ethernet frame's address fields and why. The next question is what a device on the other end of the wire actually *does* with those addresses — and that's where the story splits sharply in two: the naive, collision-prone hub, and the switch that replaced it. That's Chapter 30.
