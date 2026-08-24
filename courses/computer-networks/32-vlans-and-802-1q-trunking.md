# Chapter 32: VLANs and 802.1Q Trunking

> *"A switch has ports, not opinions — plug two devices into the same switch and, by default, they're on the same network, whether or not that's what anyone actually wants. VLANs are how you tell the switch to have opinions."*

---

## Table of Contents

1. [The Problem: One Physical Switch, Many Logical Needs](#1-the-problem-one-physical-switch-many-logical-needs)
2. [A Naive Attempt: Separate Physical Switches](#2-a-naive-attempt-separate-physical-switches)
3. [The Real Solution: VLANs](#3-the-real-solution-vlans)
4. [Access Ports: One Device, One VLAN](#4-access-ports-one-device-one-vlan)
5. [The Problem Access Ports Alone Can't Solve](#5-the-problem-access-ports-alone-cant-solve)
6. [Trunk Ports and the 802.1Q Tag](#6-trunk-ports-and-the-8021q-tag)
7. [The 802.1Q Tag, Field by Field](#7-the-8021q-tag-field-by-field)
8. [Worked Example: Two Switches, Two VLANs, One Trunk](#8-worked-example-two-switches-two-vlans-one-trunk)
9. [Native VLAN](#9-native-vlan)
10. [How VLANs Change Collision and Broadcast Domains](#10-how-vlans-change-collision-and-broadcast-domains)
11. [Why VLANs Exist: Isolation, Security, Organization](#11-why-vlans-exist-isolation-security-organization)
12. [Deep Dive: Inter-VLAN Routing](#12-deep-dive-inter-vlan-routing)
13. [Deep Dive: VLAN Hopping](#13-deep-dive-vlan-hopping)
14. [A Real Example: Configuring VLANs](#14-a-real-example-configuring-vlans)
15. [Code: Parsing an 802.1Q Tag in Go](#15-code-parsing-an-8021q-tag-in-go)
16. [Hands-On Experiment](#16-hands-on-experiment)
17. [Common Misconceptions](#17-common-misconceptions)
18. [Production Notes](#18-production-notes)
19. [What's Simplified Here](#19-whats-simplified-here)
20. [Interview Questions & Model Answers](#20-interview-questions--model-answers)
21. [Exercises](#21-exercises)
22. [Summary, and the Next Problem](#22-summary-and-the-next-problem)

---

## 1. The Problem: One Physical Switch, Many Logical Needs

Chapters 30 and 31 built up a precise picture of what a switch does: it learns MAC addresses, forwards known unicast traffic precisely, and floods broadcasts and unknowns to every port. Section 7 of Chapter 30 was explicit about a consequence of this: every port on a switch belongs to **one single broadcast domain**, spanning the entire device (and any switches connected to it), with no exceptions built into the algorithm.

Now consider a real organization's actual needs. A company has an HR department handling sensitive employee data, an engineering team, and a guest Wi-Fi network for visitors — all in the same building, plausibly all cabled into the same physical switches in the same wiring closet, because running entirely separate cabling for each department would be absurd. But the company very much does *not* want a guest's laptop to be able to see broadcast traffic from HR's payroll system, or engineering's internal tools to be reachable from the same broadcast domain as an unmanaged guest device. Chapter 30 already told you this plainly: **plain switching does nothing to solve this.** Every device plugged into the switch is, as far as the algorithm from Chapter 31 is concerned, an equal citizen of one shared network.

The problem, precisely: how do you make devices that share the same physical switches, cabling, and wiring closets behave as though they're on **completely separate logical networks** — separate broadcast domains, unable to see each other's broadcast or unknown-unicast-flooded traffic at all — without buying and cabling a separate physical switch for every department?

## 2. A Naive Attempt: Separate Physical Switches

The most obvious answer: give HR its own switch, engineering its own switch, and guests their own switch, and never physically connect them (except through a router that can enforce policy between them, foreshadowing Chapter 44).

This actually *works*, in the narrow sense that it achieves genuine isolation — but it fails on nearly every practical dimension:

- **Cost.** Every department needs its own dedicated switching hardware, even if no single department is using anywhere near the capacity of a full switch.
- **Wasted capacity.** A 48-port switch dedicated to a 6-person HR team leaves 42 ports permanently unused, while another department's switch might be oversubscribed.
- **Inflexibility.** Reorganizing which desks belong to which department (something that happens constantly in any real office) means physically re-cabling devices to different switches — impossible to do quickly, and a nightmare to keep track of at scale.
- **It doesn't compose across floors or buildings.** If HR has people on three different floors, each with their own wiring closet, "HR's switch" would need to somehow span all three physical locations, which defeats the entire premise of a switch being a local device.

What's needed instead is a way to get the *isolation* benefit of separate physical switches while keeping the *flexibility and cost efficiency* of one shared physical infrastructure — logical separation on top of shared hardware, the same fundamental move that runs through this entire course starting with layering itself (Chapter 24).

## 3. The Real Solution: VLANs

A **VLAN (Virtual LAN)** lets a single physical switch (or a whole set of interconnected physical switches) be logically partitioned into multiple independent broadcast domains, each behaving — as far as any device attached to it can tell — exactly like its own separate, physically isolated LAN, even though the underlying wiring and hardware is entirely shared.

**Intuitive level:** think of a large open-plan office building with movable partition walls. The physical building (the switches and cabling) never changes, but by putting up partitions in different configurations, you can make it function as three separate suites, or one open floor, or five small offices — reconfigured in an afternoon, without moving a single brick.

**Engineering level:** every port on a VLAN-capable switch is assigned membership in one or more VLANs, identified by a **VLAN ID** (a number from 1 to 4094, per the IEEE 802.1Q standard). The switch's MAC learning and forwarding algorithm from Chapter 31 is now effectively run *per VLAN* — a broadcast frame from a device in VLAN 10 is only flooded to other ports that are also members of VLAN 10, never to a port in VLAN 20, even if that port sits on the exact same physical switch, one slot away.

**Deep technical level:** the mechanism that lets this VLAN membership information travel *between* switches (since a real network almost never has just one switch) is the 802.1Q tag, covered in Sections 6 and 7.

## 4. Access Ports: One Device, One VLAN

The simplest kind of VLAN-aware switch port is an **access port**: a port assigned to exactly one VLAN, used for connecting a single end device (a PC, a printer, a phone) that has no idea VLANs even exist.

Frames arriving on an access port are treated by the switch as belonging to that port's configured VLAN — and critically, **the frame itself is not modified in any way**. The device plugged into an access port sends and receives perfectly ordinary, untagged Ethernet frames exactly as described in Chapter 28; it has no VLAN-awareness at all, and needs none. The VLAN membership is a property the *switch* tracks internally, associated with the port the frame arrived on, not something carried in the frame.

```
Switch (single device, 4 access ports, 2 VLANs configured):

  Port 1 (VLAN 10) -- PC-A  (HR)
  Port 2 (VLAN 10) -- PC-B  (HR)
  Port 3 (VLAN 20) -- PC-C  (Engineering)
  Port 4 (VLAN 20) -- PC-D  (Engineering)

  A broadcast frame from PC-A (port 1, VLAN 10) is flooded ONLY to
  port 2 (also VLAN 10) -- NOT to ports 3 or 4 (VLAN 20).

  PC-A and PC-C are on the same physical switch, but as far as
  broadcast traffic and MAC learning are concerned, they might as
  well be on two entirely separate switches in two separate buildings.
```

This alone already solves the core problem from Section 1, for devices that all sit on one physical switch. But real organizations have more than one switch — and that's where a genuinely new mechanism is needed.

## 5. The Problem Access Ports Alone Can't Solve

Suppose HR has people on floor 1, connected to Switch A, and also on floor 2, connected to Switch B — both should be part of VLAN 10, and should be able to see each other's broadcasts as if they were on the same LAN (because logically, they are meant to be). Switch A and Switch B are connected to each other by a single cable, as switches typically are.

If that single link between Switch A and Switch B just carries ordinary, untagged frames, how would Switch B know that a given frame arriving from Switch A belongs to VLAN 10 and not VLAN 20? The frame itself, as defined all the way back in Chapter 28, carries no VLAN information whatsoever — it's just a destination MAC, source MAC, EtherType, payload, and FCS. Somehow, VLAN membership information has to travel across that inter-switch link too, or VLANs would only ever work within a single physical switch — useless for anything beyond a very small office.

## 6. Trunk Ports and the 802.1Q Tag

The solution is a second kind of port, a **trunk port**, used specifically for links *between* switches (or between a switch and another VLAN-aware device, like certain servers or virtualization hosts) that need to carry traffic for **multiple VLANs simultaneously** over a single physical link.

A trunk port doesn't send frames unmodified the way an access port does. Instead, before a frame leaves a trunk port, the switch inserts a small **802.1Q tag** into the frame, right after the source MAC address and before the EtherType field — explicitly stamping which VLAN that frame belongs to, so that whatever switch receives it on the other end of the trunk knows exactly how to handle it, without any ambiguity.

```
Untagged frame (as sent from an access port, or by any ordinary device):

  +--------------+--------------+----------+---------+-----+
  | Dest MAC (6) | Src MAC (6)  | EtherType| Payload | FCS |
  +--------------+--------------+----------+---------+-----+

Tagged frame (as carried across a trunk link):

  +--------------+--------------+----------+----------+---------+-----+
  | Dest MAC (6) | Src MAC (6)  | 802.1Q   | EtherType| Payload | FCS |
  +--------------+--------------+ Tag (4)  |          |         |     |
  +--------------+--------------+----------+----------+---------+-----+
                                 ^
                                 inserted here, between source MAC
                                 and the original EtherType field
```

This is exactly why Chapter 28 flagged the maximum frame size as "1518 bytes (1522 with a VLAN tag)" — the 4-byte 802.1Q tag adds directly onto the standard maximum frame size, and every switch and NIC on a trunked network has to be able to handle frames 4 bytes larger than the untagged maximum to accommodate it.

## 7. The 802.1Q Tag, Field by Field

The 4-byte 802.1Q tag itself breaks down into two 2-byte fields:

| Field | Size | Purpose |
|---|---|---|
| TPID (Tag Protocol Identifier) | 2 bytes | Fixed value `0x8100`, signaling "an 802.1Q tag follows" — this is exactly the EtherType value from Chapter 28's table, and it's how a receiving device distinguishes a tagged frame from an ordinary untagged one in the first place |
| TCI (Tag Control Information) | 2 bytes | Subdivided into three fields below |

The TCI's 16 bits are further divided:

```
TCI (2 bytes = 16 bits):

  bits: 15  14  13 | 12  11  10  9  8  7  6  5  4  3  2  1  0
        └─ PCP ──┘ D └──────────── VLAN ID (12 bits) ────────┘
        (3 bits)   (1 bit)

  PCP  (Priority Code Point, 3 bits): traffic priority/QoS class,
       0 (lowest) to 7 (highest) -- used by switches to decide which
       frames to service first under congestion.
  DEI  (Drop Eligible Indicator, 1 bit): marks a frame as a candidate
       to be dropped first if the network is congested (historically
       called CFI, Canonical Format Indicator, in older 802.1Q text).
  VLAN ID (12 bits): the actual VLAN number, 0-4095.
```

That 12-bit VLAN ID field is the number that matters most in everyday networking conversation — 12 bits gives a theoretical range of 0 to 4095, but two values are reserved (`0` means "no VLAN, priority tagging only" and `4095` is reserved for implementation use), leaving **4094 usable VLAN IDs** per network. This specific ceiling — 4094, imposed directly by the tag's 12-bit field width — is exactly the number Chapter 99 will return to when explaining why cloud and data-center operators eventually needed VXLAN: 4094 logical networks is nowhere near enough for a large multi-tenant cloud provider serving potentially millions of separate customers.

## 8. Worked Example: Two Switches, Two VLANs, One Trunk

```
   Switch A                                      Switch B
  +---------+                                   +---------+
  | Po1 (A) |-- PC-A (VLAN 10, HR)               | Po1 (A) |-- PC-C (VLAN 10, HR)
  | Po2 (A) |-- PC-B (VLAN 20, Eng)               | Po2 (A) |-- PC-D (VLAN 20, Eng)
  | Po3 (T) |======= TRUNK (carries VLAN 10 & 20) ========| Po3 (T) |
  +---------+                                   +---------+
   (A) = access port      (T) = trunk port
```

**PC-A (VLAN 10, on Switch A) sends a broadcast frame:**

1. PC-A sends an ordinary, untagged Ethernet frame — it has no idea VLANs exist. It arrives on Switch A's port 1, an access port configured for VLAN 10.
2. Switch A, per Chapter 31's algorithm, learns PC-A's MAC address on port 1 *within the context of VLAN 10* (the table is effectively per-VLAN now), and determines this broadcast needs to reach every other port that is a member of VLAN 10.
3. Port 2 is VLAN 20 — excluded entirely. It never sees this frame.
4. Port 3 is the trunk — a member of "every VLAN configured to cross it," which includes VLAN 10. Before sending the frame out the trunk, Switch A **inserts an 802.1Q tag with VLAN ID = 10** into the frame.
5. Switch B receives the tagged frame on its own trunk port. It reads the tag, sees VLAN ID 10, and **strips the tag off** before deciding where the (now-untagged-again) frame needs to go within its own VLAN 10 membership.
6. Switch B floods the frame out its own VLAN-10 access ports — port 1, where PC-C lives — after removing the tag, so PC-C receives an ordinary, untagged frame exactly as if PC-A and PC-C were on the very same physical switch. PC-C's NIC never has to know an 802.1Q tag was ever involved.

The essential insight to take from this trace: **the tag exists purely for the switches' benefit, on the wire between them.** It's added right before a frame crosses a trunk link and removed right after it arrives on the other side, before being delivered to any device on an access port. End devices, in the overwhelming majority of real deployments, never see or generate an 802.1Q tag themselves at all.

## 9. Native VLAN

There's one deliberate exception to "trunk ports always tag": every 802.1Q trunk port has a configurable **native VLAN** (commonly VLAN 1 by default, though best practice in production networks is to change this), and frames belonging to the native VLAN are sent across the trunk **untagged**. When an untagged frame arrives on a trunk port, the switch assumes it belongs to that trunk's configured native VLAN.

This exists mostly for backward compatibility with older equipment that doesn't understand 802.1Q tags at all, and for very specific management-traffic conventions — but it's also, as Section 13 explains, a genuine security liability if left at its careless default, because it means an attacker who can control what a switch perceives as "untagged" traffic on a trunk can potentially manipulate which VLAN their frames land in.

## 10. How VLANs Change Collision and Broadcast Domains

Chapter 30 defined collision domains and broadcast domains precisely, and pointed out that switching alone shrinks collision domains to one per port but leaves the broadcast domain as one single domain spanning the whole switch. VLANs are exactly the missing piece that finally changes that second half of the picture:

| Setup | Collision domains | Broadcast domains |
|---|---|---|
| One hub, N devices | 1 (Chapter 30) | 1 (Chapter 30) |
| One switch, N devices, no VLANs | N (one per port) | 1 (whole switch) |
| One switch, N devices, 2 VLANs configured | N (one per port — unaffected by VLANs) | 2 (one per VLAN) |
| Two trunked switches, 2 VLANs spanning both | N (unaffected) | 2 (one per VLAN, now spanning both physical switches) |

This table makes the relationship exact: **collision domains are a property of switching (Chapter 30); broadcast domains are a property of VLAN configuration (this chapter).** They are controlled by entirely different mechanisms, and it's a common — and easily corrected — beginner error to conflate the two, or to assume that adding a switch automatically shrinks broadcast domains the way it shrinks collision domains.

## 11. Why VLANs Exist: Isolation, Security, Organization

Tying back to the motivating problem in Section 1, VLANs earn their place in real networks for several concrete, overlapping reasons:

- **Security isolation.** A guest Wi-Fi VLAN cannot see broadcast traffic (or, with appropriate routing policy, any traffic at all) from an internal corporate VLAN, even though both may be served by the very same physical access points and switches.
- **Broadcast domain performance.** As a network grows to hundreds or thousands of devices, broadcast traffic (ARP requests, DHCP discovery — both covered fully in Volume 8) reaching every single device becomes a real, measurable performance cost. Splitting a large flat network into several VLANs keeps each broadcast domain a manageable size.
- **Organizational and compliance boundaries.** Regulatory frameworks (e.g., PCI DSS for payment card data) often require network segmentation between systems that handle sensitive data and systems that don't — VLANs are a standard, auditable way to demonstrate and enforce that boundary.
- **Flexibility without re-cabling.** Moving a desk from engineering to HR is now a configuration change (reassign one switch port's VLAN membership) rather than a physical re-cabling job — directly solving the inflexibility problem from Section 2's naive approach.
- **Traffic-type separation.** It's extremely common to run voice traffic (VoIP phones) on a separate VLAN from data traffic, both physically arriving over the same cable to a desk (many VoIP phones have a pass-through Ethernet port for a PC), specifically so voice traffic can be given prioritized QoS handling (the PCP field from Section 7) without competing with a user's bulk file downloads.

## 12. Deep Dive: Inter-VLAN Routing

One direct, important consequence of everything above: **because VLANs are separate broadcast domains, they are also — by the definitions this course has been building since Chapter 24's layering discussion — separate IP subnets in essentially every real deployment** (the relationship between broadcast domains and IP subnets is made fully precise in Chapter 37). Two devices in different VLANs cannot simply send Ethernet frames directly to each other's MAC address the way two devices in the same VLAN can; a Layer 2 switch will never bridge traffic between two different VLANs by design — that would defeat the entire point of configuring them.

For a device in VLAN 10 to reach a device in VLAN 20, the traffic must pass through something that operates at Layer 3 — a router (Chapter 44), which can be a dedicated router device, a **Layer 3 switch** (a switch with routing capability built in, extremely common in modern enterprise networks), or a single router interface configured with sub-interfaces for each VLAN, a setup informally nicknamed **"router on a stick."** This is intentionally previewed here rather than explained in depth, because it genuinely requires IP addressing and routing concepts (Volumes 6 and 7) this course hasn't built yet — but it's worth knowing now that VLAN isolation and IP routing between VLANs are two sides of the same coin in real network design.

## 13. Deep Dive: VLAN Hopping

Because trunk ports are the mechanism that lets tagged traffic cross between switches, they're also a meaningful attack surface, previewed here and covered alongside other Layer 2/network attacks in Chapter 83. Two classic attack techniques are worth naming precisely:

- **Switch spoofing:** some switches, by default, run a dynamic trunk-negotiation protocol on their ports (Cisco's DTP is the best-known example) that will automatically turn a port into a trunk if the device on the other end asks for it. An attacker's device pretending to be a switch can potentially negotiate its own access port into a trunk port, gaining visibility into every VLAN that trunk carries. The standard defense is disabling dynamic trunk negotiation entirely on any port connected to an end device.
- **Double tagging:** an attacker on the native VLAN (Section 9) of a trunk can craft a frame with *two* stacked 802.1Q tags — an outer tag matching the trunk's native VLAN (which gets stripped, unremarkably, by the first switch, since native VLAN traffic is expected to arrive untagged and this outer tag is discarded as redundant) and an inner tag naming the attacker's actual target VLAN, which the next switch in the path then honors, effectively smuggling a frame into a VLAN the attacker was never authorized to reach. The standard defense is never using the default native VLAN (VLAN 1) for any purpose, and explicitly tagging native VLAN traffic where switch hardware supports it.

Both attacks exist specifically because of design choices this chapter has already made concrete — dynamic trunk negotiation and the native VLAN's untagged special case — which is exactly why understanding the mechanism in depth (Sections 6–9) is what makes the attacks, and their fixes, make sense.

## 14. A Real Example: Configuring VLANs

Cisco IOS-style configuration, showing exactly the access-port and trunk-port distinction from this chapter:

```
! Create the VLANs
vlan 10
 name HR
vlan 20
 name Engineering

! Configure an access port for a single device
interface GigabitEthernet0/1
 switchport mode access
 switchport access vlan 10

! Configure a trunk port carrying multiple VLANs to another switch
interface GigabitEthernet0/3
 switchport mode trunk
 switchport trunk allowed vlan 10,20
 switchport trunk native vlan 999    ! deliberately not VLAN 1 (Section 13)
```

The equivalent on Linux, using an 802.1Q sub-interface to receive tagged VLAN 10 traffic on a trunk-connected interface `eth0`:

```bash
sudo ip link add link eth0 name eth0.10 type vlan id 10
sudo ip link set eth0.10 up
sudo ip addr add 192.168.10.5/24 dev eth0.10
```

That `type vlan id 10` is Linux directly implementing Section 7's tag insertion/stripping — `eth0.10` is a virtual interface that transparently adds an 802.1Q tag with VLAN ID 10 to outgoing frames and strips/filters incoming frames by that same tag, letting a single physical NIC on a trunk link participate in multiple VLANs simultaneously, exactly the way a switch's trunk port does.

## 15. Code: Parsing an 802.1Q Tag in Go

Extending Chapter 28's frame parser to detect and decode a VLAN tag when present:

```go
package main

import (
	"encoding/binary"
	"fmt"
)

const tpid8021Q = 0x8100

type ParsedFrame struct {
	DstMAC       [6]byte
	SrcMAC       [6]byte
	HasVLANTag   bool
	VLANID       uint16
	PCP          uint8
	DEI          bool
	EtherType    uint16
	PayloadStart int
}

func parseFrame(frame []byte) (ParsedFrame, error) {
	if len(frame) < 14 {
		return ParsedFrame{}, fmt.Errorf("frame too short")
	}
	var p ParsedFrame
	copy(p.DstMAC[:], frame[0:6])
	copy(p.SrcMAC[:], frame[6:12])

	next := binary.BigEndian.Uint16(frame[12:14])
	if next == tpid8021Q {
		// An 802.1Q tag is present: 2 bytes TPID (already read) + 2 bytes TCI.
		tci := binary.BigEndian.Uint16(frame[14:16])
		p.HasVLANTag = true
		p.PCP = uint8((tci >> 13) & 0x7)     // top 3 bits
		p.DEI = (tci>>12)&0x1 == 1           // next 1 bit
		p.VLANID = tci & 0x0FFF              // bottom 12 bits
		p.EtherType = binary.BigEndian.Uint16(frame[16:18])
		p.PayloadStart = 18
	} else {
		p.EtherType = next
		p.PayloadStart = 14
	}
	return p, nil
}

func main() {
	// dst, src, TPID(0x8100), TCI(PCP=3,DEI=0,VLAN=10), EtherType(0x0800 IPv4), payload...
	frame := []byte{
		0x3c, 0x22, 0xfb, 0x12, 0x34, 0x56,
		0xaa, 0xbb, 0xcc, 0x00, 0x11, 0x22,
		0x81, 0x00, // TPID: 802.1Q tag follows
		0x60, 0x0a, // TCI: PCP=3 (011), DEI=0, VLAN ID=0x00a=10
		0x08, 0x00, // EtherType: IPv4
		0x45, 0x00, 0x00, 0x3c,
	}

	p, err := parseFrame(frame)
	if err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Printf("Has 802.1Q tag: %v\n", p.HasVLANTag)
	if p.HasVLANTag {
		fmt.Printf("  VLAN ID: %d\n", p.VLANID)
		fmt.Printf("  PCP (priority): %d\n", p.PCP)
		fmt.Printf("  DEI: %v\n", p.DEI)
	}
	fmt.Printf("EtherType: 0x%04x\n", p.EtherType)
	fmt.Printf("Payload starts at byte offset: %d\n", p.PayloadStart)
}
```

Output:

```
Has 802.1Q tag: true
  VLAN ID: 10
  PCP (priority): 3
  DEI: false
EtherType: 0x0800
Payload starts at byte offset: 18
```

## 16. Hands-On Experiment

```bash
# 1. Create a tagged VLAN sub-interface on Linux (requires the 8021q
#    kernel module, usually already loaded):
sudo modprobe 8021q
sudo ip link add link eth0 name eth0.100 type vlan id 100
sudo ip link set eth0.100 up
ip -d link show eth0.100
# Look for "vlan protocol 802.1Q id 100" in the output -- this
# confirms the kernel is now tagging outgoing frames on this
# sub-interface with VLAN ID 100 exactly as Section 7 describes.

# 2. Capture traffic on the underlying physical interface while
#    generating traffic on the VLAN sub-interface, and look for the
#    802.1Q tag directly in the capture:
sudo tcpdump -i eth0 -e -n vlan 100
# tcpdump's own "vlan 100" filter is itself proof that VLAN tags are
# a real, independently filterable field in the captured frames.
```

## 17. Common Misconceptions

- **"VLANs and subnets are the same thing."** They're closely related in nearly every real deployment (one VLAN, one subnet is the overwhelmingly common convention) but conceptually distinct: a VLAN is a Layer 2 broadcast-domain boundary; a subnet (Chapter 37) is a Layer 3 addressing boundary. It's entirely possible, if unusual, to have a mismatch between the two, and networking gear enforces VLAN membership regardless of what IP addresses happen to be configured on a device.
- **"An end device needs special software to be 'on' a VLAN."** Not when connected to an access port — Section 4 is explicit that access-port devices send and receive completely ordinary, untagged frames and need zero VLAN awareness. Only devices connected to a trunk port, or explicitly configured with a tagged sub-interface (Section 14's Linux example), ever need to understand 802.1Q tags directly.
- **"More VLANs always means more security, with no downside."** VLANs enforce broadcast-domain isolation, but by themselves they don't encrypt traffic or authenticate devices, and misconfigured trunk ports (Section 13) can undermine the isolation they're meant to provide. VLANs are one layer of a defense-in-depth strategy, not a complete security solution on their own.
- **"The native VLAN is just VLAN 1 by convention, and that's fine to leave as-is."** Leaving the default native VLAN unchanged is a well-known, real security anti-pattern, precisely because of the double-tagging attack in Section 13 — production networks are generally configured to use a dedicated, unused VLAN ID as the native VLAN on every trunk.

## 18. Production Notes

- Nearly every enterprise Wi-Fi deployment maps SSIDs to VLANs — a "Guest" SSID and a "Corporate" SSID broadcast from the very same physical access point are typically bridged into two entirely separate VLANs the moment traffic reaches the wired network, giving Wi-Fi the same isolation wired VLANs provide.
- VoIP deployments commonly use a specific Cisco-originated but now widely supported feature called the "voice VLAN," letting a single access port serve both a VoIP phone (tagged, voice VLAN) and a PC plugged into the phone's pass-through port (untagged, data VLAN) simultaneously — a real, everyday combination of access-port and tagging behavior beyond the simple access-vs-trunk split presented in this chapter.
- The 4094-VLAN ceiling (Section 7) is a real, binding constraint for large multi-tenant service providers and cloud platforms, and is the direct motivation for VXLAN's 24-bit (16-million-ID) segment identifier, covered in Chapter 99.

## 19. What's Simplified Here

This chapter presents access and trunk ports as the two port types, which covers the overwhelming majority of real configurations; some vendors also support a "dynamic access" mode that assigns VLAN membership based on authentication (802.1X) rather than static port configuration, and voice-VLAN configurations (mentioned in Section 18) technically blend access and tagged behavior on a single port. This chapter also does not cover Q-in-Q (802.1ad, "VLAN stacking" for service-provider use, distinct from the malicious double-tagging in Section 13, though built on a similar nested-tag mechanism) or Private VLANs (a further sub-segmentation feature within a single VLAN) — both real, deployed features beyond this introductory treatment.

## 20. Interview Questions & Model Answers

**Beginner: What problem do VLANs solve, and what's the simplest alternative they replace?**
"VLANs let devices that share the same physical switches behave as though they're on separate, isolated networks — separate broadcast domains — without needing separate physical switches and cabling for each group. Before VLANs, the only way to isolate departments or traffic types on the same premises was to run genuinely separate physical network hardware for each one, which is expensive, wasteful of capacity, and inflexible whenever anyone needs to move between groups."

**Intermediate: What's the difference between an access port and a trunk port, and what problem does the 802.1Q tag solve?**
"An access port belongs to exactly one VLAN and carries ordinary, untagged frames to and from a single end device that has no VLAN awareness at all. A trunk port carries traffic for multiple VLANs over one physical link, typically between two switches, and needs some way to mark which VLAN each frame belongs to as it crosses that shared link — that's what the 802.1Q tag does: a 4-byte field inserted between the source MAC address and the EtherType, carrying a 12-bit VLAN ID (among a priority field), added when a frame goes out a trunk port and stripped again before the frame reaches an access port on the other side."

**Advanced: Why are VLANs, by themselves, not sufficient for two devices in different VLANs to communicate, and what has to happen for that communication to work?**
"VLANs are strictly a Layer 2 mechanism — a switch, even a VLAN-aware one, will never bridge a frame from one VLAN to another, because that would erase the entire isolation VLANs are configured to provide. For a device in VLAN 10 to reach a device in VLAN 20, traffic has to be routed at Layer 3, through something that has an interface (or sub-interface) in each VLAN and can make an IP-level forwarding decision between them — a dedicated router, a router-on-a-stick setup using 802.1Q sub-interfaces on a single physical interface, or more commonly today, a Layer 3 switch with routing built in. This is also why, in essentially every real network, one VLAN corresponds to exactly one IP subnet — the VLAN boundary and the routing boundary are designed to line up."

## 21. Exercises

### Easy
1. What is the maximum theoretical number of VLAN IDs the 802.1Q tag's VLAN ID field can represent, and how many are actually usable, after reserved values are excluded?
2. A PC is plugged into a switch port configured as `switchport mode access` `switchport access vlan 30`. Does this PC's NIC ever see an 802.1Q tag? Why or why not?
3. Name the two 802.1Q tag sub-fields other than the VLAN ID itself, and state what each is used for.

### Medium
4. Two switches are connected by a trunk carrying VLANs 10, 20, and 30. A device in VLAN 10 on Switch A sends a broadcast frame. Trace exactly what happens to the frame as it crosses to Switch B and is delivered (or not delivered) to devices in each VLAN on Switch B.
5. Explain, using the collision-domain/broadcast-domain table from Section 10, why adding VLANs to an already-switched network changes the broadcast domain count but leaves the collision domain count completely unchanged.
6. A network engineer wants to move one desk from the Engineering VLAN to the HR VLAN. Describe exactly what configuration change accomplishes this, and contrast it with what would have been required under the "separate physical switches" approach from Section 2.

### Hard
7. Extend the Go program in Section 15 to also handle Q-in-Q (double-tagged) frames — where a second 802.1Q tag (TPID 0x8100) immediately follows the first one's EtherType field position — and explain, referencing Section 13, why a switch that unconditionally trusts and strips exactly one tag per frame at a trunk boundary is vulnerable to the double-tagging attack.
8. A company's guest Wi-Fi VLAN and its internal corporate VLAN are both configured on the same access points and switches. List at least three distinct configuration or design choices (beyond simply "assign them different VLAN IDs") that would need to be correct for this to actually provide the security isolation the company expects, drawing on Sections 9, 12, and 13.

## 22. Summary, and the Next Problem

| Term | Meaning |
|---|---|
| VLAN | A logical partition of a switch (or switches) into a separate broadcast domain |
| VLAN ID | 12-bit number (1–4094 usable) identifying a VLAN, carried in the 802.1Q tag |
| Access port | Switch port belonging to exactly one VLAN; carries untagged frames to/from one end device |
| Trunk port | Switch port carrying tagged traffic for multiple VLANs, typically between switches |
| 802.1Q tag | 4-byte field (TPID + TCI) inserted into a frame to mark its VLAN, added/stripped at trunk boundaries |
| PCP | 3-bit priority field within the 802.1Q tag, used for QoS |
| Native VLAN | The one VLAN on a trunk port whose traffic is sent untagged, by convention |
| Inter-VLAN routing | Layer 3 forwarding required for devices in different VLANs to communicate |

VLANs solved the isolation problem — but notice what Section 8's worked example quietly assumed: exactly **one** trunk link connecting Switch A and Switch B. Real networks, for reliability, almost never rely on a single link between two switches — if that one cable or port fails, every device behind it loses connectivity to the rest of the network entirely. The obvious fix is to add a second, redundant link between the switches.

But here's the problem waiting at the door: a switch's flooding behavior (Chapter 31) sends broadcast and unknown-unicast frames out *every* port that's a member of the relevant VLAN — including, now, both of those redundant links. A frame flooded out link 1 arrives at the other switch, which floods it right back out link 2, which arrives back at the first switch, which floods it out link 1 again — a loop with no natural end, multiplying a single broadcast frame into an exponentially growing storm that can saturate an entire network within seconds. Redundant links, added specifically to make the network *more* reliable, can silently make it catastrophically less reliable the moment two switches are connected by more than one path. Solving that — without giving up the redundancy that motivated the second link in the first place — is exactly the subject of Chapter 33: the Spanning Tree Protocol.
