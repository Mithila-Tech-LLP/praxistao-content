# Chapter 31: MAC Learning, Forwarding, and Flooding

> *"A switch is not born knowing your network's layout. It learns it, one frame at a time, purely by paying attention to who's talking — and it forgets what it learned if it goes quiet for too long."*

---

## Table of Contents

1. [The Problem: How Does a Switch Know Where Anything Is?](#1-the-problem-how-does-a-switch-know-where-anything-is)
2. [A Naive Attempt: Static Configuration](#2-a-naive-attempt-static-configuration)
3. [The Real Solution: Self-Learning From Source Addresses](#3-the-real-solution-self-learning-from-source-addresses)
4. [The Exact Algorithm](#4-the-exact-algorithm)
5. [Worked Example: Filling In a MAC Table Frame by Frame](#5-worked-example-filling-in-a-mac-table-frame-by-frame)
6. [Aging: Why Entries Expire](#6-aging-why-entries-expire)
7. [Flooding, Precisely](#7-flooding-precisely)
8. [Deep Dive: The CAM Table](#8-deep-dive-the-cam-table)
9. [Deep Dive: MAC Flooding Attacks](#9-deep-dive-mac-flooding-attacks)
10. [A Real Example: Reading a Switch's MAC Table](#10-a-real-example-reading-a-switchs-mac-table)
11. [Code: Simulating MAC Learning in Go](#11-code-simulating-mac-learning-in-go)
12. [Hands-On Experiment](#12-hands-on-experiment)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes](#14-production-notes)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)

---

## 1. The Problem: How Does a Switch Know Where Anything Is?

Chapter 30 established what a switch *does* differently from a hub: instead of blindly repeating every frame out every port, it reads a frame's destination MAC address and sends it out only the one port where that device actually lives. But it deliberately left open exactly *how* a switch knows which port a given MAC address is behind. A brand-new switch, fresh out of the box, has never seen any of the devices connected to it. It doesn't come pre-loaded with a wiring diagram of your office. So where does that knowledge come from?

This is a real, concrete problem with a genuinely elegant solution, and it's worth appreciating why it's non-trivial: a switch has no out-of-band way to ask "who's plugged into port 4?" It only ever sees electrical signals arriving on its ports, structured as Ethernet frames. Whatever mechanism builds its address-to-port table has to be derived entirely from the frames the switch happens to see pass through it — nothing else.

## 2. A Naive Attempt: Static Configuration

The obvious first idea: an administrator manually enters a table mapping each known MAC address to the port it's plugged into, the same way you might write "Alice's laptop → port 4" on a sticky note taped to the switch.

This fails for reasons that will feel familiar from earlier chapters' derivation arcs:

- **It doesn't scale.** A switch with 48 ports in an office with hundreds of devices moving between desks, joining, and leaving would require constant manual updates.
- **It breaks the moment anything changes.** Unplug a laptop from port 4 and plug it into port 9 (someone moved desks, or a cable got swapped during maintenance), and the static table is now simply wrong until a human notices and fixes it — during which time, frames for that device are misdirected.
- **It requires the administrator to already know every device's MAC address in advance**, which for guest devices, temporary equipment, or a network of any real size is impractical.

A dynamic, self-updating alternative is clearly needed — one that requires zero administrator input and correctly reacts when devices move.

## 3. The Real Solution: Self-Learning From Source Addresses

Here's the genuinely clever insight that makes MAC learning possible with no external configuration at all: **every frame that arrives on a switch port already tells the switch exactly what it needs to know** — not from the destination address (which is what the switch is trying to look up), but from the **source** address.

If a frame arrives on port 4 with source MAC address `aa:aa:aa:aa:aa:aa`, the switch can conclude, with total confidence, that device `aa:aa:aa:aa:aa:aa` is reachable via port 4 — because that's physically the only way that frame could have arrived. The switch doesn't need to be told anything; it just needs to watch, and remember.

**Intuitive level:** imagine you work at the front desk of a building with many rooms, and every time someone walks up to your desk from a particular hallway to hand you a letter, you jot down "this person's office is down that hallway." You never had to ask anyone directly — you just paid attention to where people came from, every single time, and built up a directory purely from that.

**Engineering level:** the switch maintains a **MAC address table** (also called a **forwarding table**, a **switching table**, or historically a **CAM table** — Section 8 explains why), mapping each learned MAC address to the port it was last seen arriving on, refreshed continuously as new frames arrive.

**Deep technical level:** this table lookup happens in dedicated hardware, in constant time, on every single frame the switch handles — Section 8 covers exactly how.

## 4. The Exact Algorithm

Every switch, on every single frame it receives on any port, performs exactly two independent steps — a **learning** step (always) and a **forwarding decision** step (always) — in this order:

```
ON EVERY FRAME received on port P, with source MAC S and
destination MAC D:

STEP 1 — LEARN:
    Look up S in the MAC address table.
    If S is not present, OR is present but mapped to a different port:
        Record (S -> P) in the table, with a fresh timestamp.
    (If S is already correctly mapped to P, just refresh the timestamp.)

STEP 2 — FORWARD:
    If D is the broadcast address (ff:ff:ff:ff:ff:ff),
       or a multicast address the switch treats as flood-worthy:
        FLOOD: send the frame out every port except P.
    Else, look up D in the MAC address table:
        If D is found, mapped to some port Q:
            If Q == P:
                DROP the frame. (The destination is on the same
                port/segment the frame arrived from — sending it
                back out would be pointless and, on a hub-connected
                segment behind that port, could even create a loop
                of unnecessary traffic.)
            Else:
                FORWARD the frame out port Q only.
        If D is not found in the table (an "unknown unicast"):
            FLOOD: send the frame out every port except P.
```

Three things about this algorithm are worth calling out explicitly, because they're the source of almost every subtlety and every interview question about it:

1. **Learning and forwarding are two separate, independent decisions**, both made on every single frame. A switch learns from the *source* address regardless of what it decides to do about the *destination* address, and it makes a forwarding decision about the destination regardless of what it just learned about the source.
2. **Flooding is not a special case bolted on — it's simply what happens when the forwarding lookup finds nothing** (or when the destination explicitly requests it, as with broadcast). It's the switch's honest fallback: "I don't know where this needs to go yet, so the only safe thing to do is send it everywhere and let the right device recognize itself."
3. **Every device on a switch's network sees every frame at least once** — the first frame from a source that hasn't spoken to a given destination yet — even on a fully switched network. This is a direct, unavoidable consequence of the algorithm, not a flaw in it: there's no way to avoid flooding an unknown destination without introducing some other mechanism (like a centralized directory) that Ethernet deliberately doesn't have.

## 5. Worked Example: Filling In a MAC Table Frame by Frame

Consider a single switch with four devices attached, one per port:

```
Port 1 -- PC-A (MAC: AA:AA:AA:AA:AA:AA)
Port 2 -- PC-B (MAC: BB:BB:BB:BB:BB:BB)
Port 3 -- PC-C (MAC: CC:CC:CC:CC:CC:CC)
Port 4 -- PC-D (MAC: DD:DD:DD:DD:DD:DD)
```

The switch powers on with a completely empty MAC address table:

```
MAC Address Table (initial state)
+-----+------+
| MAC | Port |
+-----+------+
| (empty)    |
+-----+------+
```

**Frame 1: PC-A sends a frame to PC-C.** (Source AA, destination CC, arrives on port 1.)

```
Step 1 (learn):  AA is not in the table -> learn (AA -> Port 1)
Step 2 (forward): CC is not in the table -> UNKNOWN -> FLOOD
                   Frame sent out ports 2, 3, and 4.
                   (PC-B and PC-D also receive this frame and
                    silently discard it, since it's not addressed
                    to them — only PC-C's NIC accepts it.)

MAC Address Table (after frame 1)
+-------------------+------+
| MAC               | Port |
+-------------------+------+
| AA:AA:AA:AA:AA:AA | 1    |
+-------------------+------+
```

**Frame 2: PC-C replies to PC-A.** (Source CC, destination AA, arrives on port 3.)

```
Step 1 (learn):  CC is not in the table -> learn (CC -> Port 3)
Step 2 (forward): AA IS in the table, mapped to port 1 -> KNOWN
                   Frame sent out port 1 ONLY.
                   (PC-B and PC-D never see this frame at all.)

MAC Address Table (after frame 2)
+-------------------+------+
| MAC               | Port |
+-------------------+------+
| AA:AA:AA:AA:AA:AA | 1    |
| CC:CC:CC:CC:CC:CC | 3    |
+-------------------+------+
```

Notice already: the very first frame of this conversation had to be flooded (nobody knew where CC was yet), but every frame *after* that, in both directions, is delivered with a single, precise, unicast forward — this is the whole benefit of switching materializing in exactly two frames.

**Frame 3: PC-B sends a frame to PC-A.** (Source BB, destination AA, arrives on port 2.)

```
Step 1 (learn):  BB is not in the table -> learn (BB -> Port 2)
Step 2 (forward): AA IS in the table, mapped to port 1 -> KNOWN
                   Frame sent out port 1 ONLY.

MAC Address Table (after frame 3)
+-------------------+------+
| MAC               | Port |
+-------------------+------+
| AA:AA:AA:AA:AA:AA | 1    |
| CC:CC:CC:CC:CC:CC | 3    |
| BB:BB:BB:BB:BB:BB | 2    |
+-------------------+------+
```

**Frame 4: PC-D sends a broadcast frame (e.g., an ARP request, Chapter 53).** (Source DD, destination FF:FF:FF:FF:FF:FF, arrives on port 4.)

```
Step 1 (learn):  DD is not in the table -> learn (DD -> Port 4)
Step 2 (forward): destination is the broadcast address -> FLOOD
                   Frame sent out ports 1, 2, and 3
                   (every device on the switch receives it and
                    processes it, per the broadcast address's
                    entire purpose).

MAC Address Table (after frame 4 — now fully populated)
+-------------------+------+
| MAC               | Port |
+-------------------+------+
| AA:AA:AA:AA:AA:AA | 1    |
| CC:CC:CC:CC:CC:CC | 3    |
| BB:BB:BB:BB:BB:BB | 2    |
| DD:DD:DD:DD:DD:DD | 4    |
+-------------------+------+
```

After just four frames, the switch has learned every device on the network purely by observation, with zero configuration. From this point on, as long as no device moves to a different port and every entry stays fresh (Section 6), every unicast conversation between any two of these four devices will be a precise, single-port forward — no more flooding needed for traffic between them.

**Frame 5 (bonus — what happens on a move): PC-A is unplugged from port 1 and plugged into port 5.** The very next frame PC-A sends will arrive on port 5 with source AA. Step 1 of the algorithm sees that AA is already in the table but mapped to a *different* port (1, not 5) — so it overwrites the entry to `AA -> Port 5`. The table self-corrects on the very next frame from the moved device, with no administrator involvement, exactly the property the naive static-table approach (Section 2) couldn't offer.

## 6. Aging: Why Entries Expire

A MAC address table entry isn't kept forever — each one has an **aging timer** that resets every time a fresh frame is seen with that source MAC, and the entry is deleted if the timer expires without being refreshed. Common real-world default aging times are around 300 seconds (5 minutes), though this is configurable on managed switches.

Why age entries out at all? Two concrete reasons:

- **Devices genuinely disconnect or move**, and an entry pointing traffic at a port where the device no longer exists would otherwise cause silent, permanent misdelivery (or more precisely, silent drops, since the algorithm forwards to that specific port and nothing arrives).
- **Table capacity is finite** (Section 8) — a switch handling a network with high device churn (public Wi-Fi guest networks, hot-desking offices) needs to reclaim space from addresses no longer actively communicating.

The trade-off aging time represents: too short, and the switch re-floods traffic to devices that are still present but simply quiet for a few minutes (wasteful, though not incorrect); too long, and stale entries linger uselessly (or, in a rare edge case with unlucky timing around a device moving to a new port and the old entry not yet expired, can cause a brief window of misdelivery until the source-learning correction in Section 5's Frame 5 example kicks in on the very next frame from the moved device).

## 7. Flooding, Precisely

It's worth being fully explicit about the three distinct situations that cause a switch to flood a frame, since "flooding" is used loosely in casual conversation but means something exact here:

| Situation | Why it floods |
|---|---|
| Destination MAC is the broadcast address (`ff:ff:ff:ff:ff:ff`) | By definition — a broadcast is meant for everyone on the segment |
| Destination MAC is a multicast address not otherwise being managed | The switch (absent smarter multicast-specific handling like IGMP snooping, a real optimization beyond this chapter's scope) treats it as "could be needed by more than one device" |
| Destination MAC is a valid unicast address, but simply isn't yet in the table ("unknown unicast flooding") | The switch has no idea where it lives yet, and flooding is the only way to guarantee delivery until it learns |

Flooding always excludes the port the frame arrived on — there's never a reason to send a frame back out the same wire it just came in on, since whatever generated that segment's traffic already received it directly.

## 8. Deep Dive: The CAM Table

The MAC address table in real switch hardware isn't a simple linear list scanned entry by entry — at any meaningful scale, that would be far too slow to keep up with line-rate traffic on many ports simultaneously. Instead, it's implemented using **Content-Addressable Memory (CAM)**, a specialized type of memory that can search its *entire contents* for a match in a single clock cycle, rather than the conventional approach of stepping through addresses one by one. This is precisely why the MAC address table is very commonly just called the **CAM table** in real networking documentation, vendor CLI output, and job interviews alike — the terms are used interchangeably in practice, even though "CAM table" technically refers to the hardware implementation and "MAC address table" refers to the logical structure it implements.

CAM hardware is fast but expensive and power-hungry per bit compared to ordinary memory, which is why switches have finite, sometimes surprisingly small, MAC table capacities — anywhere from a few thousand entries on inexpensive access switches to hundreds of thousands on large data-center-class switches. What happens when a table fills up completely is implementation-specific, but the safe, standards-compliant behavior is to simply stop learning new entries (or evict the oldest) and fall back to flooding for addresses that don't fit — which leads directly to a real, exploitable security weakness.

## 9. Deep Dive: MAC Flooding Attacks

If an attacker connected to a switch port deliberately sends a very high volume of frames with rapidly varying, fabricated source MAC addresses, they can exhaust the switch's finite CAM table capacity (Section 8) far faster than legitimate traffic would. Once the table is full, the switch — per the honest fallback logic in Section 4 — has no choice but to start **flooding legitimate unicast traffic** for any MAC address it can no longer hold in its table, out every port.

This is called a **MAC flooding attack** (or CAM table overflow attack), and its effect is genuinely striking: it forces a switch, which exists specifically to avoid the hub's shared-everything behavior (Chapter 30), to temporarily behave like a hub — broadcasting traffic that was meant for one specific device out to every port instead, including the attacker's. This is a real, historically significant Layer 2 attack, because it lets an attacker on a switched network passively capture traffic that switching was specifically supposed to prevent them from seeing.

Real defenses exist and are covered as a preview here (the full attack taxonomy is Chapter 83): **port security** features on managed switches let an administrator cap the number of MAC addresses learnable on a given port, or lock a port to a specific known set of addresses, directly closing off this attack.

## 10. A Real Example: Reading a Switch's MAC Table

On a real managed switch (Cisco IOS syntax shown, but every vendor has an equivalent command), the table from Section 5 would look like this once queried:

```
Switch# show mac address-table

          Mac Address Table
-------------------------------------------

Vlan    Mac Address       Type        Ports
----    -----------       --------    -----
1       aaaa.aaaa.aaaa    DYNAMIC     Gi0/1
1       bbbb.bbbb.bbbb    DYNAMIC     Gi0/2
1       cccc.cccc.cccc    DYNAMIC     Gi0/3
1       dddd.dddd.dddd    DYNAMIC     Gi0/4
```

Note the `Type` column showing `DYNAMIC` — that's exactly the self-learned entries from Section 4's algorithm, as opposed to `STATIC`, which would indicate an administrator manually pinned a MAC address to a port (occasionally still done for critical infrastructure like a server that must never be subject to a MAC flooding attack's flooding fallback). On Linux, the equivalent view of a software bridge's forwarding database uses:

```bash
bridge fdb show
# aa:aa:aa:aa:aa:aa dev eth1 master br0
# bb:bb:bb:bb:bb:bb dev eth2 master br0
```

## 11. Code: Simulating MAC Learning in Go

This program directly implements the algorithm from Section 4 and reproduces the worked example from Section 5:

```go
package main

import "fmt"

const broadcastMAC = "FF:FF:FF:FF:FF:FF"

type Switch struct {
	macTable map[string]int // MAC address -> port number
}

func NewSwitch() *Switch {
	return &Switch{macTable: make(map[string]int)}
}

type Frame struct {
	SrcMAC string
	DstMAC string
}

// Handle implements Section 4's two-step algorithm: learn, then forward.
func (s *Switch) Handle(f Frame, arrivalPort, totalPorts int) {
	// Step 1: LEARN from the source address, unconditionally.
	if existing, ok := s.macTable[f.SrcMAC]; !ok || existing != arrivalPort {
		s.macTable[f.SrcMAC] = arrivalPort
		fmt.Printf("  [learn] %s -> port %d\n", f.SrcMAC, arrivalPort)
	}

	// Step 2: FORWARD based on the destination address.
	if f.DstMAC == broadcastMAC {
		fmt.Printf("  [flood] broadcast frame from %s -> all ports except %d\n",
			f.SrcMAC, arrivalPort)
		return
	}

	if destPort, ok := s.macTable[f.DstMAC]; ok {
		if destPort == arrivalPort {
			fmt.Printf("  [drop]  %s is on the same port it arrived from\n", f.DstMAC)
			return
		}
		fmt.Printf("  [forward] %s -> port %d ONLY\n", f.DstMAC, destPort)
		return
	}

	fmt.Printf("  [flood] unknown destination %s -> all ports except %d\n",
		f.DstMAC, arrivalPort)
}

func main() {
	sw := NewSwitch()

	frames := []struct {
		f    Frame
		port int
	}{
		{Frame{"AA:AA:AA:AA:AA:AA", "CC:CC:CC:CC:CC:CC"}, 1},
		{Frame{"CC:CC:CC:CC:CC:CC", "AA:AA:AA:AA:AA:AA"}, 3},
		{Frame{"BB:BB:BB:BB:BB:BB", "AA:AA:AA:AA:AA:AA"}, 2},
		{Frame{"DD:DD:DD:DD:DD:DD", broadcastMAC}, 4},
	}

	for i, item := range frames {
		fmt.Printf("Frame %d: %s -> %s (arrived on port %d)\n",
			i+1, item.f.SrcMAC, item.f.DstMAC, item.port)
		sw.Handle(item.f, item.port, 4)
	}

	fmt.Println("\nFinal MAC address table:")
	for mac, port := range sw.macTable {
		fmt.Printf("  %s -> port %d\n", mac, port)
	}
}
```

Running this reproduces exactly the frame-by-frame table evolution from Section 5, including the flood-then-learn pattern on the first two frames and the precise single-port forward from frame 2 onward.

## 12. Hands-On Experiment

```bash
# On a Linux box acting as a software bridge (or any managed switch
# you have CLI access to), watch the forwarding database populate
# in real time:

# 1. Set up a Linux bridge with two or more interfaces (or use an
#    existing one, e.g. a container/VM host's default bridge):
bridge fdb show

# 2. Generate traffic from a couple of different devices/VMs attached
#    to the bridge (even simple pings between them), then re-run:
bridge fdb show
#    Watch new dynamic entries appear, each tied to the interface the
#    corresponding source traffic arrived on.

# 3. Leave one device idle for longer than the aging timeout (check
#    it with `bridge -d link show` or a switch's `show mac
#    address-table aging-time` equivalent) and re-check the table —
#    the idle device's entry should disappear once its aging timer
#    expires, exactly as Section 6 describes.
```

## 13. Common Misconceptions

- **"A switch looks up the destination address first, and only learns if that lookup fails."** No — Section 4 is explicit that learning from the source address happens unconditionally, on every frame, regardless of what the forwarding decision about the destination turns out to be.
- **"Flooding means something has gone wrong."** Flooding an unknown unicast destination is normal, expected, correct behavior — it's how a switch safely handles addresses it hasn't learned yet. Only *excessive, persistent* flooding (e.g., caused by a MAC flooding attack, Section 9, or a genuine network loop, foreshadowed for Chapter 33) indicates a problem.
- **"The MAC table is permanent once an entry is learned."** Entries expire via aging (Section 6) specifically so the table stays accurate as devices move or disconnect.
- **"Switches can't be tricked, since they only learn from real traffic."** They absolutely can — a MAC flooding attack (Section 9) exploits exactly this trusting, unconditional learning behavior; the switch has no way to distinguish a legitimate source address from a forged one.

## 14. Production Notes

- Port security (limiting learned MAC addresses per port, or locking a port to specific addresses) is a standard, widely deployed defense against MAC flooding attacks on any network exposed to untrusted physical access (e.g., publicly accessible ports, conference rooms, retail floors).
- Seeing the *same* MAC address rapidly and repeatedly relearn on *different* ports ("MAC flapping") in switch logs is a classic diagnostic sign of either a physical network loop (Chapter 33) or a duplicate MAC address (Chapter 29) somewhere on the network — worth knowing as a real troubleshooting heuristic.
- Static MAC table entries are occasionally configured deliberately for critical servers, specifically so that even a full CAM table overflow attack cannot force their traffic into a flooded, more easily sniffed state.

## 15. What's Simplified Here

This chapter presents the algorithm as if a switch only ever has one MAC address table shared across the whole device; in reality, once VLANs are introduced (Chapter 32), the table is more precisely keyed by (VLAN, MAC address) pairs, since the same MAC address could theoretically appear in different VLANs. IGMP snooping (an optimization that lets a switch forward IP multicast traffic only to ports that have explicitly joined a multicast group, rather than flooding it like ordinary unknown unicast or broadcast traffic) is mentioned but not detailed, since it depends on IP concepts not introduced until Volume 6.

## 16. Interview Questions & Model Answers

**Beginner: How does a switch learn which port a device is on?**
"By examining the source MAC address of every frame it receives and recording which port that frame arrived on. It doesn't need to be told anything in advance — the very act of a device sending any frame teaches the switch where that device lives."

**Intermediate: What does a switch do when it receives a frame for a destination MAC address it has never seen before?**
"It floods the frame — sends it out every port except the one it arrived on — because it has no way to know which specific port the destination is behind yet. This is a normal, expected part of the algorithm, not an error condition. Once the actual destination device responds (or otherwise sends any frame), the switch learns its address-to-port mapping from that frame's source address, and all subsequent traffic to it becomes a precise, single-port forward instead of a flood."

**Advanced: Explain a MAC flooding attack and why it's effective against switches specifically because of how MAC learning works.**
"A switch's MAC address table has finite hardware capacity, implemented as content-addressable memory (CAM). A MAC flooding attack sends a very high volume of frames with rapidly varying, often fabricated source MAC addresses from one port, filling the table with junk entries faster than it can be reasonably sized to hold. Because switches by design learn unconditionally from any source address they see with no authentication, they can't distinguish this from legitimate traffic. Once the table is full, the switch's standards-compliant fallback for any address it can't fit is to flood — meaning legitimate unicast traffic for MAC addresses the switch could no longer accommodate gets broadcast out every port instead of delivered precisely, letting the attacker passively capture traffic that switching, by design, was supposed to keep them from seeing. Port security features that cap or lock the number of MAC addresses learnable per port are the standard defense."

## 17. Exercises

### Easy
1. A switch's MAC table is completely empty. PC-X sends a frame to PC-Y for the first time. What does the switch do, step by step?
2. What is the difference between a `DYNAMIC` and a `STATIC` entry in a switch's MAC address table?
3. Why does a switch never forward a frame back out the same port it arrived on?

### Medium
4. Trace through Section 4's algorithm for this sequence of frames arriving at an empty 3-port switch: (1) src=X, dst=Y, port=1; (2) src=Y, dst=X, port=2; (3) src=Z, dst=X, port=3. Show the MAC table's state after each frame and state explicitly which frames were flooded versus forwarded.
5. Explain why a MAC address table entry aging out after a period of inactivity is a reasonable design choice rather than a flaw, using a concrete scenario where NOT aging out entries would cause a real problem.
6. A switch's `show mac address-table` output shows the same MAC address relearning on rapidly alternating ports over just a few seconds ("flapping"). List two distinct root causes that could produce this symptom.

### Hard
7. Extend the Go simulation in Section 11 to implement aging: add a simulated "clock tick" concept, record a last-seen timestamp per table entry, and expire entries that haven't been refreshed within a configurable number of ticks. Demonstrate a device's entry being learned, going idle, expiring, and then needing to be re-flooded to on its next communication.
8. Design, in prose, a simple heuristic a switch could use to detect a likely MAC flooding attack in progress (Section 9), using only information available from the ordinary MAC learning process (rate of new, never-before-seen source MAC addresses per port, for example), and discuss one legitimate scenario that could trigger a false positive under your heuristic.

## 18. Summary

| Term | Meaning |
|---|---|
| MAC address table (CAM table) | Switch's dynamic mapping of learned MAC addresses to ports |
| Learning | Recording (source MAC -> arrival port) from every frame received |
| Forwarding | Sending a frame out only the port its destination is known to be on |
| Flooding | Sending a frame out every port except the one it arrived on, for unknown or broadcast destinations |
| Aging | Expiring table entries after a period of inactivity so the table stays accurate |
| Unknown unicast | A valid unicast destination address not yet present in the table |
| MAC flooding attack | Deliberately overflowing the CAM table to force a switch into hub-like flooding behavior |

You now know exactly how one switch learns and forwards within its own network. But real networks aren't one switch and a handful of devices sitting equally — organizations want to segment devices on the *same* physical switch into logically separate networks entirely, for isolation, security, and manageability. That's Chapter 32: VLANs and 802.1Q trunking.
