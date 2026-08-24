# Chapter 03: What Is a Network, and Why Do Computers Need One?

> *"Two computers that can talk to each other are a curiosity. A hundred computers that can all talk to each other, without each one needing a private wire to every other one, is a network — and getting there is a harder problem than it first looks."*

---

## Table of Contents

1. [The Problem: One Computer, Alone, Can't Do Much For Anyone](#1-the-problem-one-computer-alone-cant-do-much-for-anyone)
2. [The Naive Solution: Wire Everyone to Everyone](#2-the-naive-solution-wire-everyone-to-everyone)
3. [Why Full Mesh Fails — The Math of N²](#3-why-full-mesh-fails--the-math-of-n)
4. [The Real Solution, Part 1: A Shared Medium](#4-the-real-solution-part-1-a-shared-medium)
5. [The New Problem Sharing Creates: Who Is This For?](#5-the-new-problem-sharing-creates-who-is-this-for)
6. [The Real Solution, Part 2: Addressing](#6-the-real-solution-part-2-addressing)
7. [The Problem Addressing Alone Doesn't Solve: Taking Turns](#7-the-problem-addressing-alone-doesnt-solve-taking-turns)
8. [Defining "Network," Precisely, for the First Time](#8-defining-network-precisely-for-the-first-time)
9. [Topologies: Different Shapes for the Same Idea](#9-topologies-different-shapes-for-the-same-idea)
10. [Production Notes: Shared-Medium-Plus-Addressing Beyond Ethernet](#10-production-notes-shared-medium-plus-addressing-beyond-ethernet)
11. [Hands-On Experiment: Count the Cables](#11-hands-on-experiment-count-the-cables)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Connections Backward and Forward](#13-connections-backward-and-forward)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#16-summary)

---

## 1. The Problem: One Computer, Alone, Can't Do Much For Anyone

Chapters 1 and 2 solved communication between exactly two computers: a message, encoded into bits, encoded again into a physical signal, sent across a channel, and decoded back on the other end. That's a complete, working answer — for two computers, joined by one dedicated wire, that already know exactly where that wire goes.

Real computing was never that simple. From nearly the moment computers existed, people wanted more than one pair talking on one private wire. A university wanted its dozens of research computers to share data with each other. A company wanted its employees' desktops to share one expensive printer. A researcher on one campus wanted to use processing time on a powerful computer at a different campus, hundreds of kilometers away. In every one of these cases, the requirement is the same: **more than two computers, all of which might need to reach any of the others, without a human manually running a new wire every time two new computers need to talk.**

This chapter asks the obvious next question, and works through it exactly the way real engineers had to: state the problem, try the obvious solution, watch it break, and then find the real one.

---

## 2. The Naive Solution: Wire Everyone to Everyone

If Chapter 1 and 2 taught us that two computers need a dedicated channel between them, the most direct extension to more computers is: **give every pair of computers their own dedicated wire.** This is called a **full mesh** topology, and for a small number of computers, it works exactly as well as you'd hope.

Here is a full mesh of 4 computers — every one of the 6 possible pairs gets its own direct link:

```
        A-------B
        |\     /|
        | \   / |
        |  \ /  |
        |   X   |
        |  / \  |
        | /   \ |
        |/     \|
        D-------C
```

Every computer can reach every other computer directly, with no intermediary, no ambiguity, and — crucially — no need to figure out "who is this signal for," because each wire only ever connects exactly two computers, just like Chapter 1 and 2's model. This is genuinely the simplest possible design, and for very small, fixed groups of machines (a handful of servers in the same rack that never change), full mesh wiring is still used today for exactly this reason: simplicity and no shared-medium complications.

---

## 3. Why Full Mesh Fails — The Math of N²

### The problem, made concrete

Now suppose the university in Section 1 doesn't have 4 computers — it has 50. Or the company doesn't have 4 employees — it has 5,000. How many wires does full mesh require?

The number of unique pairs among `N` computers is given by the combination formula:

```
links required = N × (N - 1) / 2
```

Let's compute this for a range of realistic values, in Go, to see exactly how bad this gets:

```go
package main

import "fmt"

func fullMeshLinks(n int) int {
	return n * (n - 1) / 2
}

func main() {
	sizes := []int{2, 4, 10, 50, 100, 1000, 10000}
	for _, n := range sizes {
		fmt.Printf("%-6d computers -> %d dedicated links needed\n", n, fullMeshLinks(n))
	}
}
```

Output:

```
2      computers -> 1 dedicated links needed
4      computers -> 6 dedicated links needed
10     computers -> 45 dedicated links needed
50     computers -> 1225 dedicated links needed
100    computers -> 4950 dedicated links needed
1000   computers -> 499500 dedicated links needed
10000  computers -> 49995000 dedicated links needed
```

Look at the jump: doubling the number of computers from 50 to 100 doesn't double the number of links — it nearly quadruples it (1,225 → 4,950). This is because the link count grows **quadratically** (proportional to N²), not linearly, with the number of computers. This is the single most important piece of math in this chapter, and it's worth internalizing precisely, because it's the mathematical fact underneath why every real network in existence — including the Internet — is *not* built as a full mesh.

### Why quadratic growth is fatal in practice, not just annoying

It's not merely that 49,995,000 wires for 10,000 computers is "a lot." It's that this approach has several kinds of failure baked into its structure, each worth naming precisely because each one motivates a specific real design later in this course:

1. **Physical impossibility.** Every computer in a full mesh needs one physical network connector *per other computer*. A computer with 10,000 peers would need 9,999 separate ports. Real network interfaces have a small, fixed number of ports (often just one or two). Full mesh doesn't just get expensive at scale — past a certain size, it becomes physically unbuildable with any real hardware.
2. **No room to grow.** Adding a single new computer to an N-computer full mesh requires running N new wires — one to every existing computer. A network that gets harder to extend every time you extend it does not scale, in the engineering sense of that word (a system "scales" if the cost of growing it grows proportionally to, or more slowly than, the growth itself — full mesh's N² cost is the textbook opposite of scaling).
3. **Massive idle capacity.** Most of those wires would sit unused most of the time. Computer A might talk to computer Z once a week; a dedicated, always-available wire between them is enormous overkill for that usage pattern. This is the same wasted-capacity problem Chapter 8 will show plagued the old telephone network, decades before computer networking existed — full mesh is really a special case of a much older, much more general problem: dedicating a whole channel to a conversation that doesn't need it constantly.

### Intuitive explanation, and where it breaks

This is exactly like a small town where everyone builds a private road directly to every other house. It works fine for 5 houses. It becomes obviously unworkable for a town of 50,000 houses — not because private roads are a bad idea in principle, but because the *number* of private roads required grows far faster than the number of houses. The analogy is useful because it points immediately at the real solution real towns actually use: shared roads, with houses connecting to a common road network, and an addressing system (a street address) so that mail and visitors can find the right house without a dedicated road to every one. That is precisely the shape of the solution the rest of this chapter builds.

---

## 4. The Real Solution, Part 1: A Shared Medium

### The insight

If dedicating a private wire to every pair is the problem, the fix is to stop doing that: instead of many private wires, use **one shared wire (or shared radio spectrum) that every computer connects to.**

```
    A       B       C       D
    |       |       |       |
    +-------+-------+-------+
              SHARED WIRE
      (every computer's signal
       travels down this one
       common link)
```

This single change turns an N² wiring problem into a linear one: each new computer needs exactly one connection — to the shared medium — regardless of how many other computers are already present. Adding a 10,001st computer to this design costs one new connection, not 10,000. This is the single most important structural idea in all of networking, and it will resurface, in more sophisticated forms, for the rest of this course: Ethernet's original design (Chapter 28) was literally a shared coaxial cable that every computer tapped into; Wi-Fi (Volume 13) is a shared medium too, except the "shared wire" is empty air; even the modern Internet backbone, at its core, is built from links that many conversations share simultaneously rather than each getting a dedicated line (a concept called statistical multiplexing, formalized properly in Chapter 9).

### Deep technical view

A shared medium in early real networks was often literally a single length of coaxial cable, with computers connected via taps along its length (the original 10BASE5/10BASE2 Ethernet standards worked exactly this way, decades before switches — Chapter 30 covers this history). In a wireless shared medium, "the wire" is simply the region of radio spectrum and physical space where transmitted signals can be heard by other devices — physically different from a cable, but functionally identical for the purposes of this chapter's reasoning.

---

## 5. The New Problem Sharing Creates: Who Is This For?

Solving the wiring explosion immediately creates a brand-new problem that full mesh never had. On the full mesh network in Section 2, if a signal arrived on the wire connecting A and B, there was zero ambiguity about who it was for — that wire only ever connected those two computers. On a shared medium, **every signal put onto the wire is physically received by every computer connected to it**, whether or not it was meant for them.

If computer A wants to send a message only to computer C, but B and D are also listening on the same shared wire, how does the shared medium know to deliver the message only to C? The honest answer, at the physical level, is: **it doesn't — it can't.** The wire has no intelligence. Every computer on it receives every signal. What has to happen instead is that **the message itself must say who it's for**, and every computer must inspect every signal it physically receives and decide, based on that inspection, whether to actually pay attention to it or ignore it.

This is a genuinely new kind of problem — not a wiring problem, but an *identification* problem — and it's the direct motivation for the next idea.

---

## 6. The Real Solution, Part 2: Addressing

### The idea

Give every computer on the shared medium a unique identifier — an **address** — and require every message sent onto the shared medium to explicitly state, as part of the message itself, the address of the intended recipient. Every computer on the wire still physically receives every signal (that part of the physics from Section 5 doesn't change), but now each computer can look at the stated destination address, compare it to its own address, and simply discard anything not addressed to it.

```
     A sends: "TO: C -- meeting at 3pm"
                    |
    +---------------+---------------+---------------+
    |               |               |               |
    v               v               v               v
    A               B               C               D
 (sender,      (receives      (receives       (receives
  ignores       signal,        signal, sees    signal, sees
  own           sees dest      dest = C, its   dest = C, NOT
  broadcast)    = C, not       own address --  its own
                itself --      ACCEPTS and     address --
                DISCARDS)      processes it)   DISCARDS)
```

This is a direct, mechanical implementation of the shared-code idea from Chapter 1 — B and D aren't confused or malfunctioning when they "receive" a message not meant for them; they are behaving exactly as agreed: inspect the address field, and only act on messages addressed to you. Every real networking address you will meet later in this course — a MAC address (Chapter 29), an IP address (Chapter 36) — exists to solve exactly this problem, at a different scale and with different additional requirements layered on top. The core mechanism, though, was fully present the moment this chapter introduced it: **a unique identifier per device, stated explicitly in every message, checked by every recipient.**

### A worked simulation

Here's a small Go program simulating exactly this: one shared medium, several computers, and address-based filtering.

```go
package main

import "fmt"

type Computer struct {
	Address string
}

type Message struct {
	From string
	To   string
	Body string
}

// broadcastToSharedMedium simulates the physical reality: every computer
// on the shared wire receives every signal, regardless of whom it's for.
func broadcastToSharedMedium(computers []Computer, msg Message) {
	fmt.Printf("Computer %s puts a signal on the shared wire, addressed TO: %s\n", msg.From, msg.To)
	for _, c := range computers {
		if c.Address == msg.From {
			continue // the sender doesn't process its own transmission
		}
		if c.Address == msg.To {
			fmt.Printf("  Computer %s: address matches -> ACCEPTS: %q\n", c.Address, msg.Body)
		} else {
			fmt.Printf("  Computer %s: address does not match -> discards silently\n", c.Address)
		}
	}
}

func main() {
	network := []Computer{{"A"}, {"B"}, {"C"}, {"D"}}
	broadcastToSharedMedium(network, Message{From: "A", To: "C", Body: "meeting at 3pm"})
}
```

Output:

```
Computer A puts a signal on the shared wire, addressed TO: C
  Computer B: address does not match -> discards silently
  Computer C: address matches -> ACCEPTS: "meeting at 3pm"
  Computer D: address does not match -> discards silently
```

Every computer physically received the identical signal (this program models that faithfully — `broadcastToSharedMedium` really does hand the message to every computer). Only one of them acted on it. That is the entire concept of addressing on a shared medium, in eleven lines of logic.

### Extending the simulation: measuring wasted work

It's worth quantifying exactly how much "wasted" inspection work this approach costs every computer on the medium, because it's directly relevant to Chapter 30's later discussion of why switches were eventually invented to reduce it. If N computers share one medium and only one message is sent, N-1 computers do the work of receiving and inspecting a signal that isn't for them, and only 1 computer does useful work:

```go
func wastedInspections(n int) int {
	return n - 2 // everyone except the sender and the one intended recipient
}
```

For a 4-computer network, that's 2 wasted inspections per message — trivial. For a 1,000-computer shared medium, that's 998 computers doing pointless work for every single message sent anywhere on the network. This inefficiency, small at Section 2's scale but serious at real organizational scale, is precisely the reason Chapter 30 exists: a switch is, at its core, a device that lets a shared-looking network stop actually broadcasting every signal to every computer, while keeping all the addressing benefits this chapter derived.

### Comparing setup cost, not just link count

Section 3 measured full mesh's cost in *links*. It's worth also measuring the cost in *work required to set the network up in the first place*, because this is what an installer or network engineer actually experiences, and it makes the N² problem even more visceral:

```go
package main

import "fmt"

// fullMeshSetupSteps assumes each of the N(N-1)/2 links requires physically
// running a cable and configuring both endpoints to recognize each other.
func fullMeshSetupSteps(n int) int {
	return n * (n - 1) / 2 * 2 // one config step per endpoint, per link
}

// sharedMediumSetupSteps assumes each new computer requires exactly one
// connection to the shared medium, plus one address assignment.
func sharedMediumSetupSteps(n int) int {
	return n * 2 // one cable/radio connection + one address, per computer
}

func main() {
	for _, n := range []int{4, 10, 50, 200} {
		fmt.Printf("N=%-4d  full mesh: %6d setup steps   shared medium: %4d setup steps\n",
			n, fullMeshSetupSteps(n), sharedMediumSetupSteps(n))
	}
}
```

```
N=4     full mesh:     12 setup steps   shared medium:    8 setup steps
N=10    full mesh:     90 setup steps   shared medium:   20 setup steps
N=50    full mesh:   2450 setup steps   shared medium:  100 setup steps
N=200   full mesh:  39800 setup steps   shared medium:  400 setup steps
```

By N=200, full mesh requires roughly 100 times more setup work than the shared-medium design — not because anyone made a mistake, but because full mesh's cost is fundamentally quadratic (O(N²), in the notation this course will use again when discussing algorithm and protocol efficiency) while the shared-medium design's cost is linear (O(N)). This distinction — linear versus quadratic growth — reappears constantly later in this course, including when comparing routing algorithms in Volume 7 and evaluating why certain naive approaches to problems like MAC address learning (Chapter 31) or spanning tree computation (Chapter 33) had to be replaced with more efficient ones as networks grew.

---

## 7. The Problem Addressing Alone Doesn't Solve: Taking Turns

Addressing solves "who should process this," but it does not solve a different, purely physical problem: **what happens if two computers transmit onto the shared wire at exactly the same moment?** On a copper wire, two computers driving different voltages onto the same conductor at once don't produce two separate signals — they produce electrical interference, and the resulting garbled voltage is unreadable to everyone. On shared radio spectrum, two overlapping transmissions similarly interfere and can become undecodable to receivers. This is called a **collision**, and it's a purely physical consequence of sharing one medium among more than one potential simultaneous transmitter — Chapter 1's "noise" idea, but self-inflicted by the network's own participants rather than an external disturbance.

This chapter will not solve this problem — that's a deliberate choice, because the real solutions (carrier sensing, collision detection, and the switch-based redesign that made the problem mostly disappear on wired LANs) require vocabulary this course hasn't built yet. What matters right now is simply naming the problem honestly and marking exactly where it gets solved: Chapter 30 covers how early Ethernet handled this with a scheme called CSMA/CD, Chapter 31 covers how modern switches sidestep the problem almost entirely by not sharing the medium the way this chapter describes, and Chapter 87 covers Wi-Fi's version of the same problem (CSMA/CA), since air is a shared medium every device in range must still take turns using.

---

## 8. Defining "Network," Precisely, for the First Time

We now have everything needed for a precise definition, rather than a vague gesture at "computers that are connected":

> **A network is a set of computers connected by shared, addressable communication links, such that any computer can send a message to any other computer by specifying its address, without requiring a dedicated physical link to every possible recipient.**

Every word in that sentence was earned in this chapter, not asserted: "shared" because Section 4 showed why dedicated links to everyone don't scale; "addressable" because Section 6 showed why a shared medium needs a way to identify intended recipients; "any computer to any other" because that's the actual requirement Section 1 opened with; and "without requiring a dedicated physical link to every possible recipient" because that's precisely the property full mesh (Section 2) lacked and a shared, addressed medium (Sections 4–6) provides.

This definition will get refined, not replaced, as the course progresses. Chapter 6 will extend it to cover networks made of other networks (a genuinely new idea, not just "more of the same"). Chapter 44 will introduce the router as a device that connects separate networks together — something this chapter's single shared medium doesn't yet need. But the core requirement — shared links, addresses, any-to-any reachability without full-mesh wiring — never goes away. Everything in this 131-chapter course is, in some sense, an increasingly sophisticated way of satisfying this one definition at increasing scale.

---

## 9. Topologies: Different Shapes for the Same Idea

The "one shared wire" picture in Section 4 is called a **bus topology**, and it was how the earliest real Ethernet networks were physically built. It's not the only physical shape that achieves the same logical goal (shared, addressable connectivity without full mesh), and it's worth seeing two more, because you'll meet all three again in Volume 5:

```
BUS TOPOLOGY                    STAR TOPOLOGY                RING TOPOLOGY

A   B   C   D                    A     B                      A---B
|   |   |   |                     \   /                       |    \
+---+---+---+                      \ /                        |     C
  (one shared                    [HUB/SWITCH]                 |    /
   cable)                          / \                        D---+
                                   /   \
                                  C     D
```

- **Bus**: every device shares one physical cable directly (Section 4's picture exactly). A single break in the cable can take down the entire network — a real, practical downside of true physical sharing.
- **Star**: every device has its own cable to a central device (a hub or switch). This looks superficially like full mesh's "everyone gets their own wire," but the crucial difference is that each device only needs **one** wire — to the central point — not one wire per peer. The central device is then responsible for getting a signal to the right place, which, as Chapter 30 will show, a switch does far more intelligently than a hub does. This is the topology essentially every modern wired office and home network actually uses.
- **Ring**: each device connects to exactly two neighbors, forming a loop, and a message may pass through several devices before reaching its destination. Less common in modern networks but historically significant (e.g., Token Ring, an early competitor to Ethernet).

The important thing to notice: **all three topologies satisfy this chapter's Section 8 definition of a network** — none of them require full mesh's N² wiring, and all of them need some form of addressing to work (even ring topology, where a device must recognize whether a passing message is addressed to it or should be forwarded onward). Topology is about the physical or logical *shape* of the connections; it doesn't change the fundamental requirements a network has to satisfy.

---

## 10. Production Notes: Shared-Medium-Plus-Addressing Beyond Ethernet

It's worth seeing that this chapter's core pattern — shared medium, plus addresses, plus a way to take turns — is not a networking-specific trick invented once for Ethernet. It's a general engineering solution to a general problem, and it reappears in systems you may not think of as "networks" at all:

- **USB.** A single USB host controller can address up to 127 devices over what is, at the protocol level, a shared bus — every USB device is assigned a unique address during enumeration (the process that happens the instant you plug in a new device), and the host explicitly addresses which device it wants to talk to for each transaction, exactly matching Section 6's addressing pattern.
- **I2C and CAN bus.** Inside a car, a laptop, or almost any modern electronic device with multiple chips, small internal "networks" connect sensors and controllers using a shared pair of wires (I2C) or a shared differential pair (CAN bus, used throughout automotive systems for engine, braking, and infotainment communication). Every device on these buses has an address, and collisions or bus contention are handled by rules strikingly similar in spirit to Section 7's problem, adapted for the specific electrical characteristics of each bus.
- **Bluetooth piconets.** A small Bluetooth network (a "piconet") of up to eight active devices shares radio spectrum in a small area, with one device acting as a coordinator that assigns short-lived addresses to the others — a miniature, single-room version of exactly the shared-medium-plus-addressing pattern this chapter derived for entire buildings and cities.

The pattern generalizes because the underlying problem generalizes: **any time more than two things need to communicate over a physical resource that can't scale as a private link to everyone, you will find shared media, addresses, and turn-taking rules, in some form.** Recognizing this pattern is often more useful for a working engineer than memorizing any one specific protocol's details, because it tells you what questions to ask about an unfamiliar system: what's the shared resource, how are participants addressed, and how do they take turns?

---

## 11. Hands-On Experiment: Count the Cables

You can verify Section 3's N² math with your own environment.

**Steps:**

1. Count the number of network-capable devices in your home or workplace that could plausibly want to reach each other directly — laptops, phones, smart TVs, printers, smart speakers. Call this number N.
2. Compute `N × (N - 1) / 2` — this is how many dedicated cables a full-mesh design would require to let every device talk directly to every other device.
3. Now count how many actual network cables (or a single Wi-Fi connection, which counts as "one connection to the shared medium" per device) each device actually has. It should be close to 1 per device — because your home network is a star or a shared-medium design (Wi-Fi), not a full mesh.
4. If you have access to a router's admin page (usually reachable via a browser at an address like `192.168.1.1` or `192.168.0.1` — don't worry yet about what that address format means, Chapter 36 covers it fully), look for a page listing "connected devices." Each one has a unique address assigned to it — a live, real-world instance of Section 6's addressing solution.

**What this demonstrates:** your real home network, almost certainly without you ever thinking about it in these terms, made exactly the engineering trade-off this chapter derived from first principles — shared medium plus addressing, instead of full mesh — because a household with even 10 devices would need 45 dedicated cables under full mesh, and nobody builds home networks that way.

---

## 12. Common Misconceptions

- **"More connections between computers is always better for reliability."** Full mesh does have a genuine reliability advantage (no single shared link can take down the whole network), but this chapter shows that advantage is bought at an unsustainable cost past a small scale. Real networks recover reliability differently — through redundant paths in a shared design (Chapter 33's Spanning Tree Protocol exists specifically to safely add *some* redundancy back into a shared-medium-style network without reintroducing full-mesh's cost).
- **"A network is just 'a bunch of connected computers.'"** As Section 8 makes precise, a network specifically requires shared, addressable links that let any device reach any other without needing a dedicated physical path to each one. Two computers joined only by a single dedicated wire (Chapter 1 and 2's whole setup) technically let them communicate, but by this chapter's definition, calling that pair-of-two a "network" undersells what the word is doing — the interesting engineering starts exactly at the point where more than two devices need to share infrastructure.
- **"Addressing is only an IP thing."** IP addressing (Chapter 36) is one specific, very important instance of the general addressing idea introduced in Section 6 — but it's neither the first nor the only one. MAC addresses (Chapter 29) solve the identical shared-medium addressing problem at a different layer, and Section 10's USB, I2C, CAN bus, and Bluetooth examples show the same pattern reused far outside traditional "networking" contexts entirely.
- **"Shared medium means slow."** A shared medium does mean devices must coordinate to avoid collisions (Section 7), which does add overhead compared to a private dedicated link — but "slow" is not automatic. Modern shared media (switched Ethernet, modern Wi-Fi) use much cleverer coordination than early bus-topology Ethernet did, and can approach the throughput of dedicated links for most realistic traffic patterns. Chapters 30–31 show exactly how switches largely eliminated the classic shared-medium slowdown for wired LANs.
- **"O(N²) growth is just a slower version of O(N) growth — it doesn't really matter which one you have."** As Section 10's setup-cost comparison shows, the two curves diverge dramatically, not gradually — at N=200, full mesh already needs roughly 100 times more setup work than the shared-medium design, and the ratio keeps growing without bound as N increases. This distinction between linear and quadratic cost is one of the most practically important ideas in all of computer science and networking, not a minor technicality.

---

## 13. Connections Backward and Forward

This chapter took the two-computer, one-channel model from Chapters 1 and 2 and asked what breaks when you need many computers to reach each other — arriving, through the N² wiring failure and the shared-medium-plus-addressing solution, at a precise definition of "network." Every remaining part of this course is a refinement of this chapter's two central ideas — sharing and addressing — at greater scale and sophistication: Chapter 4 asks what changes about this picture when the computers are spread across a building versus across a continent; Chapters 28–31 give the fully engineered version of shared-medium addressing (Ethernet and MAC addresses); Chapter 36 does the same for addressing across the entire globe (IPv4); and Volume 7 (Routing) solves an even harder version of "how does a message find its way to the right address" when there's no longer a single shared medium at all, but thousands of independently-owned networks joined together — which is exactly the subject Chapter 6 starts previewing.

---

## 14. Interview Questions & Model Answers

**Q1 (Beginner): Why doesn't a full-mesh network (a dedicated link between every pair of computers) scale to large numbers of devices?**

*Model answer:* The number of links required grows quadratically with the number of computers, following N × (N − 1) / 2. For example, 100 computers require 4,950 dedicated links, and 10,000 computers require nearly 50 million. This is not just expensive — each computer would need one physical network port per peer, which quickly becomes physically impossible with real hardware, and adding a single new computer requires wiring it to every existing computer, meaning the network gets harder to extend the larger it gets.

**Q2 (Beginner): What problem does a shared medium introduce that a full-mesh network doesn't have, and how is it solved?**

*Model answer:* On a full-mesh network, a wire only ever connects two specific computers, so there's never any doubt who a signal is for. On a shared medium, every device physically receives every signal sent on it, so the network needs a way to identify the intended recipient. This is solved with addressing: every device is given a unique identifier, every message states the address of its intended recipient, and every device compares the stated address against its own, discarding messages not addressed to it.

**Q3 (Intermediate): What is a "collision" in the context of a shared network medium, and why can addressing alone not prevent it?**

*Model answer:* A collision occurs when two or more devices transmit onto the same shared medium (a wire or shared radio spectrum) at the same time, causing their signals to physically interfere and become unreadable to receivers. Addressing only determines who should process a message that has been successfully received — it doesn't control when devices are allowed to transmit in the first place. Preventing or handling collisions requires a separate mechanism for coordinating access to the shared medium over time, such as CSMA/CD (used by early Ethernet, Chapter 30) or CSMA/CA (used by Wi-Fi, Chapter 87).

**Q4 (Intermediate): Compare bus, star, and ring topologies in terms of how they satisfy the requirement of "any device can reach any other device without full-mesh wiring."**

*Model answer:* All three avoid the N² wiring cost of full mesh, but achieve it differently. Bus topology has every device tap into one shared physical cable, so each device needs only one connection, at the cost of the whole network being vulnerable to a single cable fault. Star topology gives each device a single dedicated link to a central device (hub or switch), which is then responsible for delivering signals appropriately — this isolates individual link failures better than bus topology. Ring topology connects each device to exactly two neighbors in a loop, and a message may need to pass through intermediate devices to reach its destination, each of which must check whether the message is addressed to it or should be forwarded onward.

**Q5 (Advanced): The chapter defines a network as requiring "shared, addressable communication links." Explain why both properties — shared and addressable — are individually necessary, using a scenario where only one of the two is present.**

*Model answer:* A set of computers connected by dedicated, one-to-one links (like full mesh) is addressable trivially (a wire's two endpoints are unambiguous) but not shared — and as shown, this fails to scale past a small number of devices due to N² link growth. Conversely, a set of computers connected to one shared medium with no addressing scheme at all would let every device physically receive every transmission, but with no way to determine which transmissions matter to which device — every device would need to process, or attempt to interpret, every message ever sent on the medium, indistinguishable from noise for messages not meant for it, and no way to build point-to-point conversations at all. Both properties are necessary together: sharing solves the wiring scalability problem, and addressing solves the resulting problem of determining who a given message is meant for.

**Q6 (Advanced): This chapter shows the shared-medium-plus-addressing pattern in USB, I2C, CAN bus, and Bluetooth, well outside traditional computer networking. Why does the same pattern reappear in such different systems?**

*Model answer:* The underlying problem — multiple devices needing to communicate over a shared physical resource without paying the cost of a dedicated link to every possible peer — is not specific to Ethernet or to computer networks generally; it's a general consequence of physical and economic constraints whenever more than two devices must communicate. Any system facing that constraint (a car's internal electronics, a computer's internal peripherals, a room full of Bluetooth devices) will tend to converge on the same structural solution: a shared channel, unique addresses per participant, and a coordination rule for taking turns, because these are the minimum ingredients that solve the problem, not an arbitrary Ethernet-specific design choice.

---

## 15. Exercises

### Easy

1. Compute, by hand, the number of dedicated links a full-mesh network would need for 6 computers, and for 20 computers. Show your work using the formula from Section 3.
2. In your own words, explain why a shared medium alone (without addressing) is not sufficient to build a usable network.
3. Identify which topology (bus, star, or ring) your home or workplace network most closely resembles, and explain how you determined this.
4. From Section 10, pick one non-Ethernet example (USB, I2C/CAN bus, or Bluetooth) and identify, in your own words, what plays the role of "shared medium" and what plays the role of "address" in that system.

### Medium

1. Modify the Go simulation in Section 6 to add a fifth computer, "E", and simulate two messages being sent in sequence: one from A to E, and one from B to D. Print the full output and verify every computer behaves according to the addressing rule.
2. A colleague argues, "We should wire our 30-person office as a full mesh for maximum speed and reliability." Using this chapter's math and reasoning, write a short (3–5 sentence) explanation of why this is a poor engineering choice, and what you would suggest instead.
3. Explain why adding a single new computer to a full-mesh network of N computers requires N new links, while adding a single new computer to a shared-medium (bus) network requires only 1 new link. Use the formula from Section 3 in your explanation.
4. Using Section 6's "wasted inspections" idea, compute how many computers on a 500-device shared medium do "wasted" work for a single message sent from one device to another, and explain why this number motivates the invention of switches (previewed here, covered fully in Chapter 30).

### Hard

1. Suppose a shared bus network of 6 computers experiences a physical break in the cable exactly in the middle, separating computers {A, B, C} from {D, E, F}. Using this chapter's definition of a network from Section 8, is the result still one network? Justify your answer using the definition's exact wording, and explain what would need to be true for the two halves to be considered separate networks.
2. Design (in plain language, no code required) a rule a device on a shared medium could use to decide when it's "safe" to transmit, in order to reduce (not necessarily eliminate) the chance of a collision with another device also wanting to transmit. You are informally deriving part of the idea behind CSMA (Chapter 30) — a reasonable, well-justified guess is the goal.
3. This chapter's addressing solution (Section 6) assumes every device already knows its own unique address and every other device's address when it needs to send a message. In a real, newly-assembled network, how would a device likely first learn its own address, and how might it learn another device's address before sending it a message? You are not expected to know the real protocols yet (DHCP and ARP, covered in Chapters 55 and 53) — reason from first principles about what information would need to be exchanged, and when.

---

## 16. Summary

| Term | Meaning |
|---|---|
| Full mesh | A topology with a dedicated link between every pair of devices; link count grows as N × (N − 1) / 2 |
| N² wiring problem | The unsustainable quadratic growth in required links as a full-mesh network's device count increases |
| Shared medium | A single physical (or logical) link that multiple devices connect to and transmit signals on |
| Addressing | Assigning each device a unique identifier so messages can state an intended recipient |
| Collision | Interference caused by two or more devices transmitting on a shared medium at the same instant |
| Network (first precise definition) | A set of computers connected by shared, addressable links, letting any computer reach any other without a dedicated physical link to each one |
| Bus topology | All devices share one physical cable directly |
| Star topology | Each device connects individually to a central device (hub/switch), which handles delivery |
| Ring topology | Devices connect in a loop, each to exactly two neighbors |
| Scaling (engineering sense) | A system's cost of growth increasing proportionally to, or more slowly than, the growth itself |
| Shared-medium-plus-addressing pattern | The general engineering pattern (shared channel, unique addresses, turn-taking) that reappears in USB, I2C, CAN bus, and Bluetooth, not just Ethernet |

This chapter derived why networks are built on shared, addressable links rather than full-mesh wiring, and gave "network" its first precise definition. Chapter 4 asks what changes about this picture — physically, practically, and in terms of ownership — when the computers involved are a few meters apart versus scattered across an entire continent.
