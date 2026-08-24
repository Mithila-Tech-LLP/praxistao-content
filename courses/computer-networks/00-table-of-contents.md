# Computer Networks
### From "What Is a Network?" to Designing the Internet's Backbone

---

## What This Course Is About

Right now, information is leaving this device, crossing rooms, cities, oceans, and satellites, and arriving at another device thousands of kilometers away — in milliseconds. Almost nobody who uses this every day actually knows how it works.

This course explains it. All of it. Starting from the single most basic question anyone can ask about networking — **"What is a network, and why do computers need one?"** — and ending with you being able to explain, at the level of physics, electronics, protocols, and global infrastructure, exactly what happens when a browser loads `https://www.google.com`.

The course is built around one teaching philosophy: **every protocol was invented to solve a specific, concrete problem.** Before you learn what TCP *is*, you'll be shown the problem TCP exists to solve — and you'll be invited to imagine solving it yourself, badly, before seeing how engineers actually solved it. Every major idea in this course is derived, not memorized.

By the end, you will understand:

- What actually travels between two computers, at the level of electrical signals and light pulses
- How Ethernet, switches, and VLANs move data across a local network
- How IP addressing, subnetting, and routing move data across the *entire* Internet
- How TCP turns an unreliable network into a reliable one, byte by byte
- How DNS turns names into addresses, and HTTP turns bytes into web pages
- How TLS makes it safe to send a password across a network full of strangers
- How Wi-Fi, 4G/5G, data centers, and cloud networks are built and operated
- How to build real network software yourself — TCP servers, an HTTP server, a DNS resolver, a load balancer, a router, a VPN
- How to debug real networking problems by reasoning from symptom to root cause
- Where networking is headed — satellites, quantum networks, 6G — and how to tell deployed technology from speculation

## Who This Course Is For

Anyone. This course assumes you can use a computer and have seen a little bit of programming, but it assumes **zero** prior knowledge of networking. It does not start with the OSI model or a list of protocols to memorize — it starts with the problem of getting one bit of information from one place to another, and builds upward from there, one honest step at a time.

If you are a beginner, read from Chapter 01. If you already know the basics and want the deep, production-grade internals (BGP route selection, TCP congestion control algorithms, eBPF, eBPF-based Kubernetes networking, eBPF/XDP), you can jump directly to the relevant volume — each chapter says what it assumes and links back to where that assumption was built.

## How This Course Teaches

For every major concept, the chapter follows the same arc:

1. **What problem exists?** — stated in plain language, before any jargon.
2. **Why is it hard?** — what makes the obvious solution insufficient.
3. **A naive attempt** — what you might try first, and why it breaks.
4. **The real solution** — how engineers actually solved it, and why.
5. **Its limits** — what the real solution still cannot do.
6. **What came next** — what replaced it or grew on top of it.

Every concept is explained at **three levels**: an intuitive explanation with a real-world analogy (and where that analogy breaks down), the engineering terminology and mechanism, and — where it matters for real engineering work — the deep technical detail: packet formats, algorithms, OS behavior, and production considerations.

## How Long Will This Take

```
PART  1 — ABSOLUTE BEGINNER                         (Vol 1)         ~  4 hours
PART  2 — HISTORY OF NETWORKING                      (Vol 2)         ~  5 hours
PART  3 — THE PHYSICAL LAYER                         (Vol 3)         ~ 10 hours
PART  4 — NETWORK ARCHITECTURE & LAYERING            (Vol 4)         ~  4 hours
PART  5 — ETHERNET AND LOCAL NETWORKS                (Vol 5)         ~  9 hours
PART  6 — THE INTERNET PROTOCOL (IP)                 (Vol 6)         ~  9 hours
PART  7 — ROUTING                                    (Vol 7)         ~ 10 hours
PART  8 — ARP, ICMP & NETWORK TOOLS                  (Vol 8)         ~  4 hours
PART  9 — THE TRANSPORT LAYER (TCP/UDP)              (Vol 9)         ~ 11 hours
PART 10 — DNS                                        (Vol 10)        ~  5 hours
PART 11 — HTTP AND THE WEB                           (Vol 11)        ~  8 hours
PART 12 — NETWORK SECURITY                           (Vol 12)        ~ 11 hours
PART 13 — WI-FI                                      (Vol 13)        ~  5 hours
PART 14 — CELLULAR NETWORKS                          (Vol 14)        ~  5 hours
PART 15 — DATA CENTERS & CLOUD NETWORKING            (Vol 15)        ~  6 hours
PART 16 — ADVANCED NETWORKING (SDN, eBPF, K8s)       (Vol 16)        ~  9 hours
PART 17 — NETWORK PROGRAMMING (build it yourself)    (Vol 17)        ~ 20 hours
PART 18 — OBSERVABILITY & DEBUGGING                  (Vol 18)        ~  7 hours
PART 19 — THE INTERNET AT GLOBAL SCALE               (Vol 19)        ~  5 hours
PART 20 — THE ULTIMATE CAPSTONE                      (Vol 20)        ~  2 hours
PART 21 — FUTURE NETWORKING                          (Vol 21)        ~  4 hours
                                                                     ----------
GRAND TOTAL:                                                        ~153 hours
```

At one hour a day, that's about five months. Read one chapter at a time — each one stands on the ones before it, and tells you exactly which chapters those were.

---

## Full Table of Contents

---

# PART 1 — ABSOLUTE BEGINNER

## Volume 1: Your First Mental Model of a Network

> Before switches, IP addresses, or the word "protocol," there is a much older question: how do you get information from one mind — or one machine — to another? This volume builds the mental model everything else in the course hangs on.

### [Chapter 01: What Is a Computer, Communication, and Information?](01-what-is-a-computer-communication-and-information.md)

What a computer actually does (stores and transforms symbols), what "communication" means stripped of all technology (a sender, a receiver, a shared code), and what "information" means in the engineering sense (something that resolves uncertainty). Sets up the question the entire course answers: how do two machines share information when they are not in the same room?

**Key topics:** computers as symbol processors, sender/receiver/channel/code, information vs. noise, why "meaning" requires agreement.

### [Chapter 02: What Is a Signal?](02-what-is-a-signal.md)

A signal is a physical thing (voltage, light, radio wave) that carries a symbol. Explains why any communication, no matter how abstract, must eventually become something physical that travels through space. Introduces the idea of encoding: agreeing in advance what a given voltage or pulse of light *means*.

**Key topics:** signals as physical carriers, encoding and decoding, why all communication is eventually physics.

### [Chapter 03: What Is a Network, and Why Do Computers Need One?](03-what-is-a-network.md)

The problem: one computer, alone, cannot do anything for anyone else. Naive solution: wire every computer directly to every other computer. Why that fails at scale (the math of full-mesh wiring). The real solution: shared, addressable, general-purpose links — a *network*. Defines "network" precisely for the first time.

**Key topics:** the N² wiring problem, addressing, shared medium vs. dedicated links, first definition of "network."

### [Chapter 04: LAN, WAN, MAN, and PAN — Networks at Different Scales](04-lan-wan-man-and-pan.md)

Why a network inside one building behaves completely differently from a network spanning a continent (distance, latency, ownership, cost). Defines LAN, WAN, MAN, and PAN with concrete examples of each, and explains why the Internet is fundamentally a WAN built out of interconnected LANs.

**Key topics:** LAN vs. WAN vs. MAN vs. PAN, latency and distance, ownership boundaries.

### [Chapter 05: The Internet, the Web, and Intranets — Untangling the Terms](05-the-internet-the-web-and-intranets.md)

The most commonly confused terms in computing, disentangled: the Internet (a physical and logical network of networks), the Web (one application that runs on top of it), and an intranet (a private network using the same technology). Explains why "the internet is down" and "the website is down" are almost always different problems.

**Key topics:** Internet vs. Web vs. intranet, why they get confused, layering as the reason they're separable.

### [Chapter 06: A Network of Networks — Your First Mental Model of the Internet](06-a-network-of-networks.md)

Builds the picture the rest of the course will refine for 130 more chapters: home networks connect to ISPs, ISPs connect to each other, and no single organization owns or controls the whole thing. Introduces, at an intuitive level only, the words "router," "protocol," and "address" that later volumes will define rigorously.

**Key topics:** network of networks, no central owner, first sketch of router/protocol/address, what this course will build toward.

---

# PART 2 — HISTORY OF NETWORKING

## Volume 2: How We Got Here

> Every major idea in modern networking was a response to a real limitation of what came before it. This volume tells that story in order, so that when you meet TCP/IP formally in Volume 9, you already understand *why* it looks the way it does.

### [Chapter 07: Before Computers Talked — The Telegraph and the Telephone](07-before-computers-talked.md)

The telegraph as the first system to turn information into an electrical signal sent over distance, and Morse code as the first shared "protocol." The telephone's leap to real-time analog voice. Why both systems' way of allocating a physical line per conversation set up the problem the next chapter solves.

**Key topics:** telegraph, Morse code as protocol, telephone, dedicated physical circuits.

### [Chapter 08: Circuit Switching — How the Phone Network Worked](08-circuit-switching.md)

The problem: connecting any caller to any other caller. The circuit-switched solution: reserve an unbroken physical (or logical) path for the whole call. Why this is simple and reliable but wastes capacity — a line sits idle during silence, and one call can't share a wire with another.

**Key topics:** circuit switching, dedicated end-to-end paths, capacity waste, telephone exchanges.

### [Chapter 09: Packet Switching — The Idea That Changed Everything](09-packet-switching.md)

The problem with circuit switching restated for computer data (bursty, not continuous). Paul Baran's and Donald Davies's insight: chop data into small, independently-routed packets that share links. Why this lets many conversations share the same wire, and why it means no single link failure takes down the whole network.

**Key topics:** packet switching, statistical multiplexing, resilience to failure, Baran and Davies.

### [Chapter 10: ARPANET and the IMPs — The Internet's First Ancestor](10-arpanet-and-the-imps.md)

The real 1969 network: four university computers, a funding agency (ARPA), and a purpose-built minicomputer (the IMP) that did the packet switching so host computers didn't have to. Why this division of labor — "let a dedicated box handle networking" — is still how routers work today.

**Key topics:** ARPANET, IMPs, host vs. network functions, first packet-switched WAN.

### [Chapter 11: TCP/IP and the Birth of "Inter-networking"](11-tcp-ip-and-internetworking.md)

The problem: ARPANET's own internal protocol only worked within ARPANET. What happens when you want to connect *networks of different kinds* together? Cerf and Kahn's answer: a protocol that doesn't care what's underneath it. Introduces the core idea of TCP/IP — a common language between different networks — without yet diving into packet formats (that's Volumes 6 and 9).

**Key topics:** internetworking problem, Cerf and Kahn, protocol independence from underlying network, "the network of networks" made literal.

### [Chapter 12: NSFNET, Privatization, and the Rise of the Commercial Internet](12-nsfnet-and-the-commercial-internet.md)

How a government-funded backbone (NSFNET) grew the Internet beyond research institutions, and why removing the "no commercial use" rule in the 1990s changed everything. The birth of commercial ISPs and Internet exchange points.

**Key topics:** NSFNET, acceptable use policy, privatization, birth of ISPs.

### [Chapter 13: The Web, Broadband, Smartphones, Cloud — The Modern Internet Takes Shape](13-the-web-broadband-smartphones-and-cloud.md)

Tim Berners-Lee's Web as one application on top of the existing Internet (not a new network). The shift from dial-up to broadband, from desktops to smartphones, and from owning servers to renting them (cloud computing). Sets up why the rest of this course spends so much time on HTTP, mobile networks, and data centers.

**Key topics:** World Wide Web, broadband, smartphones, cloud computing, why the Internet's *use* changed even as its core protocols didn't.

---

# PART 3 — THE PHYSICAL LAYER

## Volume 3: What Actually Travels Between Two Computers

> Strip away every protocol name and acronym, and a network is still just wires, light, and radio waves carrying voltage changes. This volume answers the question most courses skip: what is *physically* moving, and what are its hard physical limits?

### [Chapter 14: What Physically Travels Between Two Computers?](14-what-physically-travels-between-two-computers.md)

The starting question, made unavoidable: when your computer "sends a message," what physically leaves it? Bits, bytes, and binary as a human-readable description of physical states (voltage high/low, light on/off). Why everything — text, video, a bank transfer — eventually becomes a sequence of these physical states.

**Key topics:** bits and bytes, binary representation, signals as physical events, why all data reduces to this.

### [Chapter 15: Analog vs. Digital Signals, and Modulation](15-analog-vs-digital-signals-and-modulation.md)

The difference between a signal that varies continuously (analog) and one that takes discrete values (digital), and why digital signals resist noise so much better. Modulation as the technique of riding digital information on top of an analog carrier wave (needed for radio, DSL, cable, and fiber).

**Key topics:** analog vs. digital, why digital resists noise, modulation (ASK/FSK/PSK), carrier waves.

### [Chapter 16: Frequency, Amplitude, Phase, and Bandwidth](16-frequency-amplitude-phase-and-bandwidth.md)

The three properties of a wave you can vary to encode information (frequency, amplitude, phase), and how modulating combinations of them (QAM) packs more bits into the same wave. Bandwidth defined precisely as a range of frequencies, not "how fast your internet is" — and the crucial link between the two.

**Key topics:** frequency, amplitude, phase, QAM, bandwidth (Hz) vs. throughput (bps).

### [Chapter 17: Noise, Attenuation, Interference, and SNR](17-noise-attenuation-interference-and-snr.md)

Why a signal degrades over distance (attenuation), picks up unwanted energy from other sources (interference and noise), and how Signal-to-Noise Ratio (SNR) quantifies "how much can I trust this signal." Real examples: why long Ethernet runs need repeaters, why microwaves disrupt 2.4GHz Wi-Fi.

**Key topics:** attenuation, interference, noise floor, SNR, real-world causes and fixes.

### [Chapter 18: Shannon's Limit — The Physics of How Much Data Can Fit](18-shannons-limit.md)

The problem: is there a hard ceiling on how much data can be sent through a noisy channel of given bandwidth? Claude Shannon's channel capacity theorem, explained intuitively and then with the actual formula, and why it explains everything from dial-up's 56k ceiling to why fiber can carry terabits.

**Key topics:** Shannon-Hartley theorem, channel capacity, bandwidth vs. SNR trade-off, why physical limits (not engineering laziness) cap throughput.

### [Chapter 19: Error Detection — Parity, Checksums, and CRC](19-error-detection.md)

The problem: noise sometimes flips bits in transit. Naive solution: send everything twice and compare (wasteful). Real solutions built up in order of sophistication: parity bits, checksums, and Cyclic Redundancy Check (CRC), with worked binary examples of each catching (and missing) errors.

**Key topics:** parity bit, checksum, CRC, what each can and cannot detect, worked examples.

### [Chapter 20: Error Correction — Hamming Codes and Forward Error Correction](20-error-correction.md)

The next problem: detecting an error is not the same as fixing it without asking the sender to resend. Hamming codes as the classic example of adding just enough redundancy to locate and correct a flipped bit. Forward Error Correction (FEC) in real systems (satellite links, fiber, Wi-Fi) where asking for a retransmit is too slow or too expensive.

**Key topics:** Hamming codes (worked example), FEC, where correction is preferred over retransmission and why.

### [Chapter 21: Copper and Twisted Pair — Ethernet Cabling Explained](21-copper-and-twisted-pair.md)

How electrical signals travel down copper wire, why simply running parallel wires causes crosstalk, and why twisting pairs of wires together (twisted pair) cancels it out. Cat5e/Cat6/Cat6a cabling compared, RJ45 connectors, and real distance/speed limits.

**Key topics:** copper conduction, crosstalk, twisted pair, Cat5e/6/6a, RJ45, distance limits.

### [Chapter 22: Fiber Optics — Single-Mode, Multi-Mode, and Transceivers](22-fiber-optics.md)

Why light in glass fiber can travel much farther and faster than electricity in copper (total internal reflection, near-zero electrical interference). The difference between single-mode fiber (one path, long distance) and multi-mode fiber (many paths, shorter distance, cheaper), and what a transceiver (SFP/SFP+) actually does.

**Key topics:** total internal reflection, single-mode vs. multi-mode fiber, transceivers/SFPs, why fiber beats copper for long runs.

### [Chapter 23: Wireless Media, Satellites, and the Undersea Cables That Carry the Internet](23-wireless-media-satellites-and-undersea-cables.md)

How radio and microwave links carry data without wires, and why satellite links have inherently higher latency (the speed-of-light round trip to orbit and back). Then: the physical reality that ~95% of intercontinental Internet traffic travels through undersea fiber cables, with a real map of how continents are actually connected.

**Key topics:** radio and microwave links, satellite latency, undersea cables, physical map of the global Internet's backbone.

---

# PART 4 — NETWORK ARCHITECTURE & LAYERING

## Volume 4: Why Networking Is Built in Layers

> A single protocol trying to handle "send bits over copper" and "guarantee this file arrives intact" and "render a web page" at once would be an unmaintainable monster. This volume explains the problem layering solves, then teaches the two models (OSI and TCP/IP) that formalize it.

### [Chapter 24: Why Do We Need Layers? The Problem Layering Solves](24-why-do-we-need-layers.md)

The problem, made concrete: if Ethernet, IP, TCP, and HTTP were all one tangled protocol, a change to Wi-Fi would force a rewrite of web browsers. Naive alternative: one giant protocol. Why that fails to evolve. The real solution: separate concerns into independent layers that only need a defined interface with their neighbors.

**Key topics:** the tangled-protocol problem, separation of concerns, interfaces between layers, why layering enables independent evolution.

### [Chapter 25: The OSI Model, Layer by Layer](25-the-osi-model.md)

The seven-layer OSI reference model, introduced only now that you understand *why* a layered model is needed, and after you've already met concrete examples (signals, cabling) of the bottom layer. Each layer explained with a one-sentence job description and a real device or protocol that lives there.

**Key topics:** Physical, Data Link, Network, Transport, Session, Presentation, Application layers; what OSI is (a reference model) and isn't (not what the real Internet strictly follows).

### [Chapter 26: The TCP/IP Model, and OSI vs. TCP/IP](26-the-tcp-ip-model.md)

The four/five-layer TCP/IP model that the real Internet actually implements, and an honest mapping between it and OSI. Why textbooks teach OSI (rigorous conceptual separation) but engineers talk in TCP/IP terms (what's actually deployed).

**Key topics:** TCP/IP model layers, mapping to OSI, why both models coexist in practice.

### [Chapter 27: Encapsulation and Decapsulation — Frames, Packets, Segments, Datagrams](27-encapsulation-and-decapsulation.md)

The mechanics of layering made physical: how each layer wraps the layer above it's data in its own header, and unwraps it on the way back up. Precise definitions — finally — of frame, packet, segment, and datagram, with a byte-level diagram of one HTTP request wrapped in TCP, then IP, then Ethernet.

**Key topics:** encapsulation, decapsulation, headers vs. payload, frame/packet/segment/datagram terminology, full worked wrap-up diagram.

---

# PART 5 — ETHERNET AND LOCAL NETWORKS

## Volume 5: How Computers Talk on the Same Local Network

> Before a packet can cross the world, it has to leave the room. This volume answers, in complete mechanical detail: what exactly happens when Computer A sends data to Computer B sitting a few meters away?

### [Chapter 28: Ethernet and the Ethernet Frame](28-ethernet-and-the-ethernet-frame.md)

Ethernet as the dominant technology for wired LANs, and the structure of an Ethernet frame field by field (preamble, destination/source MAC, EtherType, payload, FCS). Minimum and maximum frame sizes and why they exist.

**Key topics:** Ethernet frame format, field-by-field breakdown, frame size limits.

### [Chapter 29: MAC Addresses — Every Device's Physical Name](29-mac-addresses.md)

The problem: Ethernet frames need a destination, but IP addresses (not yet introduced in depth) are logical and layered on top. The 48-bit MAC address burned into every network interface, its structure (OUI + device ID), and why it's called a "physical" or "hardware" address.

**Key topics:** MAC address structure, OUI, burned-in vs. configurable addresses, broadcast MAC.

### [Chapter 30: Hubs vs. Switches — Collision Domains and Broadcast Domains](30-hubs-vs-switches.md)

The naive first LAN device (the hub: repeat everything to everyone) and its fatal flaw at scale — collisions. The switch's solution: learn which device is on which port and forward frames only where they need to go. Precise definitions of collision domain and broadcast domain.

**Key topics:** hubs, collisions and CSMA/CD, switches, collision domain vs. broadcast domain.

### [Chapter 31: MAC Learning, Forwarding, and Flooding](31-mac-learning-forwarding-and-flooding.md)

The exact algorithm a switch runs: build a MAC address table by watching source addresses, forward known destinations out one port, and flood unknown destinations out all ports. Worked example with a table filling in frame by frame.

**Key topics:** MAC address table, learning algorithm, forwarding, flooding, aging.

### [Chapter 32: VLANs and 802.1Q Trunking](32-vlans-and-802-1q-trunking.md)

The problem: sometimes you want devices on the same physical switch to behave like they're on *different* logical networks (isolation, security, organization). VLANs as logical LAN segmentation, the 802.1Q tag that carries a VLAN ID between switches, and the difference between access ports and trunk ports.

**Key topics:** VLANs, 802.1Q tagging, access vs. trunk ports, use cases for segmentation.

### [Chapter 33: Spanning Tree Protocol — Preventing Loops](33-spanning-tree-protocol.md)

The problem: redundant links between switches (added for reliability) create loops, and a broadcast frame in a loop multiplies forever (a broadcast storm). STP's solution: elect a root bridge and mathematically block redundant paths until they're needed. RSTP as STP's faster-converging successor.

**Key topics:** broadcast storms, redundant links, STP algorithm (root bridge, blocking ports), RSTP improvements.

### [Chapter 34: Link Aggregation](34-link-aggregation.md)

The problem: one physical link has a fixed maximum speed, and losing that one link means losing connectivity entirely. Link aggregation (LACP/802.3ad) bonding multiple physical links into one logical, higher-throughput, fault-tolerant link.

**Key topics:** link aggregation, LACP, throughput and redundancy trade-offs.

### [Chapter 35: Full Trace — What Happens When Computer A Sends Data to Computer B on the Same LAN](35-full-trace-lan-communication.md)

The synthesis chapter for this volume: a complete, no-gaps trace of one ping between two machines on the same switch — ARP resolution (previewed, fully explained in Volume 8), frame construction, switch forwarding table lookup, and the reply — with a full sequence diagram.

**Key topics:** full LAN communication trace, sequence diagram, tying together MAC addresses, switching, and (a preview of) ARP.

---

# PART 6 — THE INTERNET PROTOCOL (IP)

## Volume 6: Addressing the Entire World

> MAC addresses work beautifully inside one LAN, but they don't scale to billions of devices spread across the planet — there's no meaningful way to "look up" a MAC address globally. This volume introduces the addressing scheme that does scale: IP.

### [Chapter 36: IPv4 Addresses — What They Are and Why They Exist](36-ipv4-addresses.md)

The scaling problem MAC addresses can't solve (no hierarchy, no way to route by "region"). IPv4's 32-bit hierarchical address, dotted-decimal notation, and why hierarchy (like phone numbers and postal codes) is what makes global routing possible at all.

**Key topics:** IPv4 address structure, dotted-decimal notation, why hierarchy enables scale, MAC vs. IP addressing contrasted.

### [Chapter 37: Network and Host Portions, Subnet Masks](37-network-and-host-portions-subnet-masks.md)

How one 32-bit address is split into a "network" part and a "host" part, and how a subnet mask tells you where that split is. Worked binary AND examples showing exactly how a device decides "is this destination on my local network, or do I need a router?"

**Key topics:** network vs. host bits, subnet mask, binary AND operation, local vs. remote destination decision.

### [Chapter 38: Subnetting From First Principles](38-subnetting-from-first-principles.md)

Not a formula to memorize — the actual problem: an organization is handed one block of addresses and needs to split it into smaller, independently manageable networks (for different buildings, departments, or security zones) without wasting addresses. Builds up subnetting by solving progressively harder versions of that problem by hand.

**Key topics:** why subnetting exists, borrowing host bits, worked subnetting problems, address/broadcast address per subnet.

### [Chapter 39: CIDR — Classless Addressing](39-cidr.md)

The problem with the original rigid class A/B/C system (massive waste — a company needing 300 addresses got handed 65,536). CIDR's slash notation as the fix: any prefix length, not just /8, /16, /24. Route aggregation as CIDR's other superpower, previewed for Volume 7.

**Key topics:** classful addressing and its waste, CIDR notation, prefix lengths, aggregation preview.

### [Chapter 40: Private vs. Public Addresses, Loopback, Broadcast, Multicast](40-private-public-loopback-broadcast-multicast.md)

Why some address ranges (10.0.0.0/8, 172.16.0.0/12, 192.168.0.0/16) are reserved for private use and never routed on the public Internet. Loopback (127.0.0.1) as "talk to myself," broadcast as "reach everyone on this subnet," and multicast as "reach a specific group."

**Key topics:** RFC 1918 private ranges, loopback, subnet broadcast address, multicast basics.

### [Chapter 41: NAT — Sharing One Address Among Many](41-nat.md)

The problem private addressing creates: a device with a private address can't be reached from, or reach, the public Internet directly. Network Address Translation's solution: a router rewrites addresses (and tracks a translation table) so many private devices can share one public IP. Worked example of a NAT table.

**Key topics:** the private-addressing-meets-public-Internet problem, NAT mechanics, NAT table, why NAT complicates certain applications (previewed for later security/VPN chapters).

### [Chapter 42: IPv6 — Solving Address Exhaustion](42-ipv6.md)

The real, unavoidable problem: IPv4's ~4.3 billion addresses were never going to be enough for billions of phones, sensors, and servers. IPv6's 128-bit address space, its hexadecimal notation and shorthand rules, and why NAT's role changes (mostly disappears) in an IPv6 world.

**Key topics:** address exhaustion, IPv6 address format, shorthand notation (::), IPv6 vs. IPv4 addressing philosophy.

### [Chapter 43: SLAAC, Neighbor Discovery, and IPv6 Transition Mechanisms](43-slaac-neighbor-discovery-and-transition.md)

How an IPv6 device can configure its own address without DHCP (SLAAC), how Neighbor Discovery Protocol replaces ARP and adds router discovery, and the practical reality that IPv4 and IPv6 must coexist for years — dual-stack, NAT64, and tunneling explained.

**Key topics:** SLAAC, Neighbor Discovery Protocol, dual-stack, NAT64, 6to4/tunneling transition mechanisms.

---

# PART 7 — ROUTING

## Volume 7: How a Packet Finds Its Way Across the World

> IP addressing tells you *where* something is. Routing is the much harder problem of *how to get there* — potentially through thousands of intermediate networks you've never heard of, owned by organizations that don't trust each other. This volume is the story of how the Internet solves that.

### [Chapter 44: What Is a Router? What Is a Routing Table?](44-what-is-a-router.md)

A router as a device that connects different networks and decides, packet by packet, which direction to send each one. The routing table as the router's decision-making data structure: a list of "to reach this network, go this way."

**Key topics:** router's job, routing table structure, router vs. switch contrasted.

### [Chapter 45: Forwarding, Next Hop, Default Route, Longest Prefix Match, and TTL](45-forwarding-next-hop-and-longest-prefix-match.md)

The exact algorithm a router runs on every packet: find the most specific matching route (longest prefix match), forward to that route's next hop, and decrement TTL. Why a "default route" (0.0.0.0/0) exists as the catch-all, and what happens when TTL hits zero (with a preview of how traceroute, in Volume 8, exploits this).

**Key topics:** forwarding algorithm, next hop, default route, longest prefix match, TTL and packet death.

### [Chapter 46: Static Routing vs. Dynamic Routing](46-static-vs-dynamic-routing.md)

The problem with manually configuring every route on every router: it doesn't scale and can't react to failures. Static routing's simplicity vs. dynamic routing protocols that let routers discover and adapt to the network's actual topology automatically.

**Key topics:** static routes, dynamic routing protocols, trade-offs of each, when each is actually used in practice.

### [Chapter 47: RIP — Distance Vector Routing](47-rip.md)

The simplest dynamic routing protocol: each router tells its neighbors "here's my distance to every network I know," and everyone updates accordingly. Why this is easy to understand and implement, but suffers from slow convergence and the "count to infinity" problem — worked example included.

**Key topics:** distance vector routing, RIP mechanics, hop count metric, count-to-infinity problem, split horizon.

### [Chapter 48: OSPF and IS-IS — Link State Routing](48-ospf-and-is-is.md)

The improvement: instead of trusting secondhand distance rumors, every router builds a complete map of the network (link state) and computes the best path itself using Dijkstra's algorithm. Why OSPF (and IS-IS, its ISP-favored cousin) converge faster and scale better than RIP.

**Key topics:** link state routing, OSPF areas, Dijkstra's algorithm applied to routing, IS-IS overview.

### [Chapter 49: BGP — The Protocol That Runs the Internet](49-bgp.md)

The hardest routing problem yet: routing between organizations (autonomous systems) that don't trust each other and have business relationships, not just topology, to consider. BGP as a path-vector protocol built for policy, not just shortest path — and why "best route" on the Internet often isn't the shortest one.

**Key topics:** BGP as path-vector protocol, policy-based routing, why BGP differs fundamentally from RIP/OSPF, BGP path attributes overview.

### [Chapter 50: Autonomous Systems and Route Aggregation](50-autonomous-systems-and-route-aggregation.md)

What an Autonomous System (AS) actually is — a network under one administrative control, with its own AS number — and why BGP is fundamentally the protocol that connects ASes, not routers. Route aggregation revisited: how ISPs advertise one summarized prefix instead of thousands of specific ones.

**Key topics:** Autonomous System definition, AS numbers, route aggregation at ISP scale, why this keeps global routing tables manageable.

### [Chapter 51: Peering, Transit, and Tier-1 Networks](51-peering-transit-and-tier-1-networks.md)

The business layer underneath the technical protocol: what it means for two networks to "peer" (exchange traffic for free) versus buy "transit" (pay to reach the rest of the Internet), and what makes a network "Tier-1" (needs to buy no transit at all). Internet Exchange Points (IXPs) as the physical places peering happens.

**Key topics:** peering vs. transit, Tier-1/2/3 networks, Internet Exchange Points, business relationships behind routing decisions.

### [Chapter 52: Route Leaks and Route Hijacking — When Routing Breaks](52-route-leaks-and-hijacking.md)

What happens when BGP's trust-based design goes wrong: a misconfigured or malicious AS announces routes it has no business announcing, redirecting traffic that was never meant for it. Real historical incidents, and why the Internet has spent decades building (partial) defenses like RPKI.

**Key topics:** route leaks, BGP hijacking, real-world incidents, RPKI and route origin validation as mitigations.

---

# PART 8 — ARP, ICMP & NETWORK TOOLS

## Volume 8: The Glue Protocols and the Engineer's Toolbox

> Some protocols don't move your data — they make the protocols that do possible in the first place. This volume covers those "glue" protocols, then hands you the command-line tools every network engineer actually uses day to day.

### [Chapter 53: ARP — Translating IP Addresses to MAC Addresses](53-arp.md)

The gap left open since Volume 5: a device knows the destination's IP address, but Ethernet needs a MAC address to actually send a frame. ARP's broadcast-then-cache solution, worked packet by packet, and the ARP cache that avoids repeating the process for every single packet.

**Key topics:** ARP request/reply, ARP cache, why IP and MAC addresses must be bridged, gratuitous ARP.

### [Chapter 54: ICMP, Ping, and Traceroute](54-icmp-ping-and-traceroute.md)

ICMP as the Internet's error-and-diagnostics protocol — not for carrying data, but for routers and hosts to report problems. How `ping` uses ICMP echo request/reply to test reachability, and how `traceroute` cleverly abuses TTL expiration (from Chapter 45) to map every hop along a path.

**Key topics:** ICMP message types, ping mechanics, traceroute mechanics (TTL exploitation), reading real output.

### [Chapter 55: DHCP — How Devices Get Their Address Automatically](55-dhcp.md)

The problem: manually configuring an IP address on every device that joins a network doesn't scale (and Volume 6 assumed devices "already have" an address — this chapter explains how they got it). DHCP's DORA process (Discover, Offer, Request, Acknowledge), leases, and renewal.

**Key topics:** DHCP DORA process, lease time and renewal, DHCP relay for routed networks.

### [Chapter 56: The Network Engineer's Toolbox — ip, ss, dig, arp, tcpdump, curl](56-the-network-engineers-toolbox.md)

A hands-on lab chapter: real commands for real diagnosis, each tied back to a protocol already taught. `ip`/`ifconfig` for interface configuration, `ss`/`netstat` for socket state, `ping`/`traceroute` revisited, `dig`/`nslookup` previewed for Volume 10, `arp -a`, `curl -v`, and `tcpdump` previewed for Volume 18.

**Key topics:** ip, ss, ping, traceroute, dig, nslookup, arp, route, curl, tcpdump — first hands-on pass, with exercises.

---

# PART 9 — THE TRANSPORT LAYER (TCP/UDP)

## Volume 9: Making an Unreliable Network Reliable

> IP promises nothing: packets can be lost, duplicated, corrupted, or arrive out of order. And yet file transfers complete perfectly and video calls (mostly) don't garble your voice. This volume is the story of how two very different protocols — UDP and TCP — sit on top of that chaos and give applications what they actually need.

### [Chapter 57: Ports and Sockets — How Many Programs Share One Network Connection](57-ports-and-sockets.md)

The problem: one machine, one IP address, but dozens of programs (browser, email client, game) needing network access simultaneously. Ports as the answer — a 16-bit number identifying *which program* on a machine a packet is for — and the socket as the (IP, port) pair that identifies one actual conversation.

**Key topics:** port numbers, well-known vs. ephemeral ports, sockets, the 4-tuple that identifies a connection.

### [Chapter 58: UDP — The Simplest Transport Protocol](58-udp.md)

Introduced first because it's the honest baseline: UDP adds almost nothing to IP except ports and a checksum — no ordering, no retransmission, no connection. Why that's a *feature*, not a limitation, for DNS, video, gaming, and VoIP, and the full UDP header explained field by field.

**Key topics:** UDP header format, connectionless communication, why "unreliable" is sometimes exactly right, real UDP use cases.

### [Chapter 59: TCP — The Three-Way Handshake](59-tcp-three-way-handshake.md)

The problem UDP leaves unsolved for most applications: file transfers and web pages need every byte, in order, with no gaps. Introduces TCP by first asking "how would you build reliability on top of a network that loses and reorders packets?" then walks through the actual solution — the three-way handshake (SYN, SYN-ACK, ACK) that establishes a connection.

**Key topics:** the reliability problem stated from first principles, three-way handshake, why two-way isn't enough, initial sequence numbers.

### [Chapter 60: Sequence Numbers, Acknowledgments, and Retransmission](60-sequence-numbers-acks-and-retransmission.md)

How TCP numbers every byte (not every packet), how acknowledgments tell the sender what's been received, and how timeouts and duplicate ACKs trigger retransmission of exactly what's missing — not the whole stream.

**Key topics:** byte sequence numbers, cumulative ACKs, retransmission timeout (RTO), duplicate ACKs.

### [Chapter 61: Flow Control and the Sliding Window](61-flow-control-and-sliding-window.md)

The problem: a fast sender can overwhelm a slow receiver's buffer. The sliding window mechanism: the receiver advertises how much it can currently accept, and the sender is only allowed that much data in flight at once. Worked example with a shrinking and growing window.

**Key topics:** receiver buffer problem, sliding window, window scaling, worked numeric example.

### [Chapter 62: Congestion Control — Slow Start, AIMD, CUBIC, and BBR](62-congestion-control.md)

A different problem from flow control: not "can the receiver keep up" but "can the *network* keep up." Slow start and AIMD (Additive Increase, Multiplicative Decrease) as TCP's original congestion-avoidance strategy, then the modern algorithms (CUBIC, the Linux default; BBR, Google's model-based alternative) and why they were invented.

**Key topics:** congestion vs. flow control distinguished, slow start, AIMD, congestion avoidance, CUBIC, BBR, why algorithms keep evolving.

### [Chapter 63: Fast Retransmit and Fast Recovery](63-fast-retransmit-and-fast-recovery.md)

The problem with waiting for a timeout to detect loss (too slow). Fast retransmit's trick: three duplicate ACKs are treated as a strong-enough loss signal to retransmit immediately, and fast recovery avoids dropping all the way back to slow start.

**Key topics:** duplicate-ACK-triggered fast retransmit, fast recovery, why this is faster than timeout-based recovery.

### [Chapter 64: Connection Termination — FIN, TIME_WAIT, and CLOSE_WAIT](64-connection-termination.md)

The four-way close handshake (FIN/ACK from each side), why closing is more delicate than opening (either side can still have data in flight), and the TIME_WAIT and CLOSE_WAIT states — including the real production problem of TIME_WAIT/port exhaustion on busy servers.

**Key topics:** four-way close, half-close, TIME_WAIT purpose and problems, CLOSE_WAIT and application-level leaks.

### [Chapter 65: The TCP Header, Field by Field](65-the-tcp-header.md)

The synthesis chapter for TCP: the complete 20+-byte TCP header laid out field by field (source/dest port, sequence/ack numbers, flags, window size, checksum, options), tying every field back to the mechanism (handshake, flow control, congestion control) that uses it, with a real captured header decoded byte by byte.

**Key topics:** full TCP header layout, every field mapped to a mechanism already taught, real packet decode.

---

# PART 10 — DNS

## Volume 10: Turning Names Into Numbers

> Nobody should have to remember 142.250.80.46 to check their email. This volume explains how a single hierarchical, cached, globally distributed system turns human-friendly names into the IP addresses the rest of the stack actually needs.

### [Chapter 66: Why DNS Exists — Names vs. Numbers](66-why-dns-exists.md)

The problem stated plainly: IP addresses (Volume 6) are perfect for routers, terrible for humans. Naive fix: a single giant text file mapping every name to every address (this literally existed — the ARPANET HOSTS.TXT file) — and why it couldn't scale past a few hundred hosts.

**Key topics:** the HOSTS.TXT precursor, why a flat file doesn't scale, the case for a hierarchical, distributed naming system.

### [Chapter 67: The DNS Hierarchy — Root, TLD, and Authoritative Servers](67-the-dns-hierarchy.md)

DNS's tree-shaped namespace and the three tiers of servers that mirror it: the 13 root server clusters, top-level domain (TLD) servers (.com, .org, .in), and authoritative servers for individual domains. Why no single server needs to know every name on Earth.

**Key topics:** DNS namespace as a tree, root servers, TLD servers, authoritative servers, delegation.

### [Chapter 68: Recursive Resolvers, Caching, and TTL](68-recursive-resolvers-caching-and-ttl.md)

The client's side of the lookup: what a recursive resolver (like your ISP's or 8.8.8.8) actually does on your behalf, walking the hierarchy so your device doesn't have to. Why caching at every level (and the TTL that controls how long) is what makes DNS fast enough to be invisible.

**Key topics:** recursive vs. iterative resolution, resolver's walk down the hierarchy, caching, TTL trade-offs.

### [Chapter 69: DNS Record Types, DNSSEC, DoH, DoT, and Anycast DNS](69-dns-record-types-and-dnssec.md)

The actual data DNS stores, record type by record type (A, AAAA, CNAME, MX, TXT, NS, SOA, SRV) with real examples. Then: DNS's security problems (plaintext, spoofable) and the fixes — DNSSEC (authenticity), DoH/DoT (privacy via encryption), and Anycast DNS (how root servers "exist" at hundreds of physical locations under one address).

**Key topics:** DNS record types, DNSSEC, DNS-over-HTTPS, DNS-over-TLS, Anycast addressing for DNS.

---

# PART 11 — HTTP AND THE WEB

## Volume 11: The Protocol Behind Every Web Page

> Everything in this volume assumes you already know how to get from a browser to a server's IP address (DNS) and how to open a reliable connection to it (TCP). This volume is about what happens once that connection exists: how a browser actually asks for a page, and how HTTP evolved as the Web's demands outgrew it.

### [Chapter 70: URLs and the Structure of a Web Request](70-urls-and-the-structure-of-a-web-request.md)

The URL decomposed piece by piece (scheme, host, port, path, query string, fragment) and how a browser turns a typed URL into the DNS lookup and TCP connection from earlier volumes before HTTP even begins.

**Key topics:** URL anatomy, scheme/host/port/path/query/fragment, what happens before the first HTTP byte is sent.

### [Chapter 71: The HTTP Request/Response Cycle — Methods, Headers, Status Codes](71-the-http-request-response-cycle.md)

The plain-text (in HTTP/1.1) structure of a request and response, HTTP methods (GET, POST, PUT, DELETE, and why each exists), headers as metadata, and status codes as a structured way to say what happened — grouped by class (2xx/3xx/4xx/5xx) with real examples.

**Key topics:** request/response structure, HTTP methods, headers, status code classes, real captured request/response.

### [Chapter 72: Cookies, Sessions, and Caching](72-cookies-sessions-and-caching.md)

The problem: HTTP is stateless, but logins and shopping carts need state. Cookies as the mechanism, sessions as the pattern built on top of them, and HTTP caching headers (`Cache-Control`, `ETag`) as the mechanism that makes the Web fast by avoiding unnecessary requests entirely.

**Key topics:** statelessness problem, cookies, sessions, cache-control/ETag, conditional requests.

### [Chapter 73: HTTP/1.0 and HTTP/1.1 — And Why They Struggled](73-http-1-0-and-http-1-1.md)

HTTP/1.0's one-request-per-connection design and its overhead, HTTP/1.1's improvements (keep-alive, pipelining), and the head-of-line blocking problem that pipelining couldn't actually fix — the browser workaround of opening 6 parallel connections per host, and why that workaround has its own costs.

**Key topics:** HTTP/1.0 vs 1.1, keep-alive, pipelining and its failure, head-of-line blocking, 6-connections-per-host workaround.

### [Chapter 74: HTTP/2 — Multiplexing Over One Connection](74-http-2.md)

The real fix for head-of-line blocking at the HTTP layer: binary framing and multiplexing multiple requests/responses over a single TCP connection, plus header compression (HPACK) and server push. Why HTTP/2 still can't fix head-of-line blocking *at the TCP layer* — setting up the motivation for HTTP/3.

**Key topics:** binary framing, multiplexing, HPACK header compression, server push, the TCP-layer limitation that remains.

### [Chapter 75: HTTP/3 and QUIC — Rebuilding Transport on UDP](75-http-3-and-quic.md)

The radical fix: stop using TCP at all. QUIC, built on UDP, reimplements reliability *per stream* so one lost packet no longer stalls every other stream, bakes in TLS 1.3 from the start, and supports 0-RTT reconnection. Why this required reinventing transport-layer plumbing that Volume 9 spent nine chapters on.

**Key topics:** QUIC over UDP, per-stream reliability, built-in encryption, 0-RTT, connection migration.

### [Chapter 76: WebSockets, Server-Sent Events, REST APIs, and Reverse Proxies](76-websockets-sse-rest-and-reverse-proxies.md)

Where plain request/response HTTP isn't enough: WebSockets for full-duplex, persistent connections (chat, live games) and Server-Sent Events for one-way server push (live feeds). REST as an architectural convention for designing HTTP APIs. Reverse proxies as the component sitting in front of real servers — previewed in depth for Volume 15 and built from scratch in Volume 17.

**Key topics:** WebSocket handshake and framing, SSE, REST conventions, reverse proxy role.

---

# PART 12 — NETWORK SECURITY

## Volume 12: Making an Untrusted Network Safe to Use

> Every protocol so far has assumed a cooperative network. This volume drops that assumption. It asks: what can go wrong when the network is full of strangers, some of them actively hostile — and what tools exist to defend against each threat?

### [Chapter 77: Threat Models — Thinking Like an Attacker and a Defender](77-threat-models.md)

Security isn't a checklist — it's a discipline of asking "who might attack this, with what capability, and what would they gain?" Introduces the threat model for a typical network conversation (an eavesdropper on the path, a malicious router, a spoofed endpoint) that every remaining chapter in this volume defends against.

**Key topics:** threat modeling, attacker capabilities (passive vs. active), what "secure" actually needs to mean here.

### [Chapter 78: Symmetric Cryptography — One Key, Shared Secret](78-symmetric-cryptography.md)

The core problem: making data unreadable to anyone without a secret. Symmetric encryption (AES) as the fast, practical solution when both sides already share a key — and the problem it does *not* solve (how do two strangers agree on a shared key over a network an eavesdropper is watching?), which the next chapter answers.

**Key topics:** symmetric encryption, AES basics, speed advantage, the key-distribution problem it creates.

### [Chapter 79: Asymmetric Cryptography — Diffie-Hellman, RSA, and ECC](79-asymmetric-cryptography.md)

The elegant answer to the key-distribution problem: mathematical trapdoor functions that let two strangers agree on a secret in public view (Diffie-Hellman) or let anyone encrypt a message only one person can decrypt (RSA). Elliptic Curve Cryptography (ECC) as the modern, faster alternative doing the same job.

**Key topics:** Diffie-Hellman key exchange (worked numeric example), RSA public/private keys, ECC, why asymmetric crypto is slow and mostly used to bootstrap symmetric keys.

### [Chapter 80: Hashing and Digital Signatures](80-hashing-and-digital-signatures.md)

Hashing as a one-way fingerprint function (and why that one-way property is the whole point), and digital signatures as the combination of hashing and asymmetric crypto that proves both authenticity ("this really came from X") and integrity ("this wasn't altered").

**Key topics:** cryptographic hash functions, collision resistance, digital signatures, authenticity vs. integrity vs. confidentiality.

### [Chapter 81: PKI and Certificate Authorities — Who Do You Trust?](81-pki-and-certificate-authorities.md)

The remaining gap: asymmetric crypto proves "this is signed by this key," but not "this key belongs to google.com." Public Key Infrastructure and Certificate Authorities as the trust hierarchy that closes that gap, how a certificate is structured, and how your browser decides which CAs to trust at all.

**Key topics:** the identity-binding problem, PKI, Certificate Authorities, certificate structure, root/intermediate/leaf trust chains.

### [Chapter 82: TLS and the TLS Handshake](82-tls-and-the-tls-handshake.md)

Putting Chapters 78–81 together into the protocol that actually protects HTTPS: the TLS handshake, walked step by step (ClientHello, ServerHello, certificate exchange, key exchange, Finished), TLS 1.2 vs. 1.3's simplifications, and how this connects back to HTTP/3's built-in TLS 1.3 from Chapter 75.

**Key topics:** full TLS handshake sequence diagram, TLS 1.2 vs. 1.3, session resumption, where HTTPS actually gets its security.

### [Chapter 83: Common Network Attacks — Sniffing, MITM, Spoofing, SYN Floods, and DDoS](83-common-network-attacks.md)

A tour of attacks, each mapped to the specific protocol weakness it exploits: passive packet sniffing (plaintext), man-in-the-middle attacks, ARP spoofing (Chapter 53's trust assumption abused), DNS cache poisoning (Chapter 68's caching abused), IP spoofing, SYN flood (Chapter 59's handshake abused), and DDoS at scale. Each attack is explained defensively — for recognition and mitigation, not exploitation.

**Key topics:** packet sniffing, MITM, ARP spoofing, DNS poisoning, IP spoofing, SYN flood, DDoS, session hijacking — each tied to the specific mechanism it abuses.

### [Chapter 84: Firewalls and Web Application Firewalls](84-firewalls-and-waf.md)

Stateless packet filtering vs. stateful firewalls that track connection state (and can therefore tell a real reply from a spoofed one), and Web Application Firewalls (WAFs) as an application-layer-aware defense sitting in front of web servers.

**Key topics:** stateless vs. stateful firewalls, connection tracking, WAF role and limitations.

### [Chapter 85: VPNs — IPsec, WireGuard, OpenVPN, and Tunneling](85-vpns.md)

The problem: you need a private, secure connection across a network you don't control (the public Internet). Tunneling as the general concept (wrapping one packet inside another), and a comparison of the three major implementations — IPsec (complex, ubiquitous), OpenVPN (flexible, TLS-based), and WireGuard (minimal, modern, fast).

**Key topics:** tunneling concept, IPsec, OpenVPN, WireGuard, trade-offs between them.

---

# PART 13 — WI-FI

## Volume 13: Networking Without Wires, in One Room

> Wi-Fi looks like magic until you realize it's solving a problem Ethernet never had to: many devices sharing one invisible, shared medium (the air) where anyone nearby can also transmit — and collide.

### [Chapter 86: Radio Communication and the Basics of 802.11](86-radio-communication-and-802-11-basics.md)

How data rides on a radio wave (building on modulation from Chapter 15), why Wi-Fi uses the 2.4GHz and 5GHz (and now 6GHz) bands specifically, and the 802.11 standard family as the umbrella all Wi-Fi generations belong to.

**Key topics:** radio basics recap, 2.4/5/6GHz bands and their trade-offs, 802.11 as a standard family.

### [Chapter 87: Access Points, SSID, BSSID, Channels, and CSMA/CA](87-access-points-ssid-bssid-and-csma-ca.md)

The role of an access point as a bridge between the wireless medium and the wired LAN, SSID (the network's name) vs. BSSID (the AP's actual MAC-based identity), channel selection to avoid interference, and CSMA/CA — Wi-Fi's answer to the shared-medium collision problem (contrasted with Ethernet's older CSMA/CD).

**Key topics:** access points, SSID vs BSSID, channels and overlap, CSMA/CA vs CSMA/CD, association and roaming.

### [Chapter 88: The Wi-Fi Generations — From 802.11a to Wi-Fi 7](88-the-wifi-generations.md)

The evolution of Wi-Fi standards in order (802.11a/b/g/n/ac/ax) with what problem each generation actually solved (speed, range, multi-device efficiency), landing on Wi-Fi 6/6E/7 and the technologies that got them there — MIMO and beamforming explained mechanically.

**Key topics:** 802.11a/b/g/n/ac/ax lineage, Wi-Fi 6/6E/7, MIMO, beamforming.

### [Chapter 89: Wi-Fi Security — WEP, WPA, WPA2, WPA3](89-wifi-security.md)

The security arms race on wireless networks: WEP's fatal flaws, WPA's stopgap fix, WPA2's AES-based solution (and the KRACK vulnerability it later faced), and WPA3's modern improvements (forward secrecy, protection against offline dictionary attacks).

**Key topics:** WEP weaknesses, WPA, WPA2 and AES, KRACK, WPA3 improvements.

---

# PART 14 — CELLULAR NETWORKS

## Volume 14: The Network in Your Pocket

> A phone's cellular connection solves problems Wi-Fi never has to: nationwide coverage, seamless handoff between moving towers, and serving voice and data to billions of devices simultaneously. This volume traces cellular technology from analog voice to 5G's software-defined core.

### [Chapter 90: 1G to 3G — GSM, CDMA, and the First Mobile Data](90-1g-to-3g.md)

1G as pure analog voice (and why it was neither secure nor efficient), 2G/GSM's digital leap (and SMS as an accidental byproduct of its signaling channel), CDMA as a competing approach, and 3G/UMTS bringing real mobile data for the first time.

**Key topics:** 1G analog, 2G/GSM digital voice, SMS origin, CDMA, 3G/UMTS mobile data.

### [Chapter 91: 4G and LTE — Mobile Broadband Arrives](91-4g-and-lte.md)

Why 4G/LTE was designed all-IP from the ground up (unlike 3G's voice-circuit legacy), the architecture (eNodeB, EPC) at a conceptual level, and VoLTE as the way voice calls finally moved onto the same IP data network as everything else.

**Key topics:** all-IP design, LTE architecture overview, VoLTE, typical real-world speeds.

### [Chapter 92: 5G — Architecture, mmWave, Massive MIMO, and Network Slicing](92-5g.md)

5G's three use-case pillars (fast mobile broadband, massive IoT, ultra-reliable low latency) and the technology behind each: Sub-6GHz vs. mmWave spectrum trade-offs, Massive MIMO and beamforming at cellular scale, the 5G Core's software-defined design, and network slicing as one physical network behaving like many virtual ones. Edge computing and private 5G as deployments this architecture enables.

**Key topics:** 5G use-case pillars, Sub-6 vs mmWave, Massive MIMO, 5G Core, network slicing, edge computing, private 5G.

### [Chapter 93: 6G and the Future of Mobile Networks](93-6g-and-the-future-of-mobile-networks.md)

What's actually true right now, stated carefully: 6G has no finalized standard yet. What research directions exist (terahertz spectrum, AI-native radio access, integrated sensing and communication) and — explicitly labeled — what is deployed technology versus standards work versus active research versus speculation.

**Key topics:** 6G research directions, terahertz communication, AI-native RAN, explicit deployed/research/speculative classification.

---

# PART 15 — DATA CENTERS & CLOUD NETWORKING

## Volume 15: How the Internet's Biggest Services Are Actually Built

> Behind every "google.com" response is not one server but a data center with tens of thousands of machines, engineered so that no single failure takes the service down and no single user waits in line behind another. This volume explains that machine.

### [Chapter 94: Inside a Data Center — Servers, NICs, and Leaf-Spine Architecture](94-inside-a-data-center.md)

The physical and logical structure of a modern data center network: servers and their NICs, top-of-rack (ToR) switches, and the leaf-spine architecture that replaced older tree topologies to give every server roughly equal bandwidth to every other server.

**Key topics:** ToR switches, leaf-spine topology, why it replaced 3-tier trees, border routers.

### [Chapter 95: Load Balancing — L4 vs. L7 and Reverse Proxies](95-load-balancing.md)

The problem: one server can't handle millions of users. Load balancing as the fix, split into Layer 4 (fast, connection-level, protocol-agnostic) and Layer 7 (slower, but can route by URL, header, or cookie). Reverse proxies revisited as the general pattern underneath both.

**Key topics:** why load balancing is necessary, L4 vs L7 load balancing, health checks, sticky sessions.

### [Chapter 96: CDNs, Caching, and Anycast at Scale](96-cdns-caching-and-anycast.md)

The problem of physical distance: a server in Virginia is slow for a user in Mumbai no matter how good the code is. Content Delivery Networks solving it by caching content at edge locations near users, and Anycast as the routing trick that lets the "same" IP address actually mean "the nearest of hundreds of servers."

**Key topics:** CDN caching, edge locations, Anycast routing, cache invalidation basics.

### [Chapter 97: Cloud Networking — VPCs, Subnets, Route Tables, and Security Groups](97-cloud-networking-vpcs-and-security-groups.md)

How cloud providers let customers build their own private, isolated network inside shared physical infrastructure: Virtual Private Clouds (VPCs), subnets within them, route tables controlling traffic flow, and security groups/NACLs as software-defined firewalls.

**Key topics:** VPCs, subnets in the cloud, route tables, security groups vs NACLs.

### [Chapter 98: NAT Gateways and Internet Gateways in the Cloud](98-nat-gateways-and-internet-gateways.md)

How cloud-hosted private resources reach the public Internet (NAT gateways, mirroring Chapter 41's NAT but managed as a service) and how public-facing resources are actually exposed (Internet gateways), with a full request path traced through a real cloud VPC.

**Key topics:** NAT gateways, Internet gateways, full traced request through a cloud VPC.

---

# PART 16 — ADVANCED NETWORKING

## Volume 16: Software Has Eaten the Network

> Everything so far assumed physical devices doing fixed jobs. This volume covers the shift to software-defined everything — networks built inside software on top of other networks, and the Linux kernel internals that make containers and Kubernetes possible.

### [Chapter 99: Network Virtualization — Overlays, Underlays, and VXLAN](99-network-virtualization-and-vxlan.md)

The problem: cloud and data-center operators need to run many isolated virtual networks on top of one physical network, at a scale (millions of VLANs' worth) that VLAN tagging (4096 max) can't reach. VXLAN's solution: encapsulate Ethernet frames inside UDP packets to build a virtual (overlay) network on top of a physical (underlay) one.

**Key topics:** overlay vs underlay networks, VXLAN encapsulation, why VLANs don't scale to cloud scale.

### [Chapter 100: Software-Defined Networking — Control Plane vs. Data Plane, OpenFlow, and NFV](100-software-defined-networking.md)

The insight behind SDN: separate the decision-making (control plane) from the packet-forwarding (data plane), and centralize the former in software. OpenFlow as the protocol that lets a central controller program switch behavior directly. Network Function Virtualization (NFV) as running what used to be dedicated hardware (firewalls, routers) as software instead.

**Key topics:** control plane vs data plane separation, SDN controllers, OpenFlow, NFV.

### [Chapter 101: Service Mesh — Envoy, mTLS, and Service Discovery](101-service-mesh.md)

The modern microservices version of "how does A find and securely talk to B": service discovery replacing static IP configuration, mutual TLS (mTLS) automatically encrypting and authenticating service-to-service traffic, and Envoy as the sidecar proxy pattern that implements both without changing application code.

**Key topics:** service discovery, mTLS, sidecar proxy pattern, Envoy's role in a service mesh.

### [Chapter 102: The Linux Networking Stack — Namespaces, veth, and Bridges](102-the-linux-networking-stack.md)

The kernel primitives everything in this volume (and Kubernetes, next) is actually built from: network namespaces as isolated network stacks within one machine, veth pairs as virtual "network cables" connecting them, and Linux bridges as software switches joining them together.

**Key topics:** network namespaces, veth pairs, Linux bridges, hands-on: building an isolated network by hand.

### [Chapter 103: Container Networking and the CNI](103-container-networking-and-cni.md)

How Docker gives each container what looks like its own network stack (using the namespaces and veth pairs from Chapter 102), and the Container Network Interface (CNI) as the standard plugin API that lets orchestrators like Kubernetes plug in different networking implementations.

**Key topics:** container network namespaces in practice, Docker bridge networking, CNI as a plugin standard.

### [Chapter 104: Kubernetes Networking](104-kubernetes-networking.md)

Kubernetes's networking model requirements (every pod gets its own IP, pods can reach each other without NAT), Services as stable virtual IPs in front of ephemeral pods, and kube-proxy's role — building directly on CNI (Chapter 103) and load balancing (Chapter 95).

**Key topics:** Kubernetes networking model, pod IPs, Services and kube-proxy, ClusterIP/NodePort/LoadBalancer.

### [Chapter 105: eBPF, XDP, TC, and Cilium](105-ebpf-xdp-tc-and-cilium.md)

The newest layer: eBPF as a way to run sandboxed programs inside the Linux kernel itself — safely and without kernel source changes — and how XDP (earliest possible packet processing) and TC (traffic control) use it for extreme-performance networking. Cilium as a real, widely-deployed Kubernetes CNI built entirely on eBPF.

**Key topics:** eBPF fundamentals, XDP, TC hooks, Cilium as an eBPF-based CNI, why eBPF is reshaping networking.

---

# PART 17 — NETWORK PROGRAMMING

## Volume 17: Build It Yourself

> Reading about protocols only gets you so far. This volume has you implement the concepts from every previous volume as working code — mostly in Go, chosen for how naturally it expresses concurrent network I/O, with Python where it makes a clearer teaching point.

### [Chapter 106: Building a TCP Server and Client](106-building-a-tcp-server-and-client.md)

The socket API made concrete: opening a listening socket, accepting connections, and building a minimal client — implementing, in code, the three-way handshake and connection lifecycle from Volume 9.

**Key topics:** Go net package, TCP listener/accept loop, TCP client dial, goroutine-per-connection pattern.

### [Chapter 107: Building a UDP Server and Client](107-building-a-udp-server-and-client.md)

The same exercise for UDP: no connection setup, no ordering guarantees — and code that has to handle that honestly, contrasting directly with Chapter 106's TCP version.

**Key topics:** Go UDP sockets, connectionless read/write, when to choose UDP in real code.

### [Chapter 108: Building a TCP Chat Application](108-building-a-tcp-chat-application.md)

The first multi-user project: a chat server handling many simultaneous client connections concurrently, broadcasting messages, and handling disconnects gracefully — built on Chapter 106's TCP server.

**Key topics:** concurrent connection handling, broadcast pattern, channels for coordination, graceful disconnect handling.

### [Chapter 109: Building an HTTP Server From Scratch](109-building-an-http-server-from-scratch.md)

Not using a framework's `http.ListenAndServe` as a black box — parsing raw HTTP requests off a TCP socket by hand (methods, headers, body) and writing a spec-correct response, tying directly back to Volume 11's request/response format.

**Key topics:** raw HTTP request parsing, building responses by hand, connecting HTTP structure to TCP bytes.

### [Chapter 110: Building an HTTP Client From Scratch](110-building-an-http-client-from-scratch.md)

The client-side mirror of Chapter 109: constructing a valid HTTP request by hand, sending it over a raw TCP connection, and parsing the response — including handling chunked transfer encoding.

**Key topics:** manual HTTP request construction, response parsing, chunked encoding.

### [Chapter 111: Building a DNS Resolver](111-building-a-dns-resolver.md)

Implementing a real (simplified) recursive DNS resolver: constructing DNS query packets by hand, parsing responses, and walking the root → TLD → authoritative hierarchy from Volume 10 in actual code.

**Key topics:** DNS wire format, UDP query construction, recursive resolution logic in code.

### [Chapter 112: Building a Reverse Proxy](112-building-a-reverse-proxy.md)

A working reverse proxy that accepts client connections and forwards them to a backend, rewriting headers as needed — the code-level realization of Chapter 76 and Chapter 95's concepts.

**Key topics:** proxying HTTP requests, header rewriting, connection pooling to backends.

### [Chapter 113: Building a Load Balancer](113-building-a-load-balancer.md)

Extending Chapter 112's reverse proxy into a load balancer implementing round-robin and least-connections strategies across multiple backends, with health checks — the code-level realization of Chapter 95.

**Key topics:** round-robin and least-connections algorithms, health checking, backend pool management.

### [Chapter 114: Building a Packet Sniffer](114-building-a-packet-sniffer.md)

Using raw sockets (or a packet-capture library) to capture and parse live traffic off a network interface, decoding Ethernet, IP, and TCP/UDP headers by hand — putting Volumes 5, 6, and 9's header formats directly into code.

**Key topics:** raw sockets / pcap, header parsing in code, live traffic capture and decode.

### [Chapter 115: Building a Simple Router](115-building-a-simple-router.md)

A software implementation of IP forwarding: maintaining a routing table, performing longest-prefix match, decrementing TTL, and forwarding packets between interfaces — the code-level realization of Volume 7's core algorithm.

**Key topics:** software IP forwarding, longest-prefix match implementation, TTL handling in code.

### [Chapter 116: Building a CDN-Like Cache](116-building-a-cdn-like-cache.md)

An HTTP caching proxy implementing `Cache-Control`/`ETag`-aware caching (Chapter 72) and cache invalidation — a simplified, single-node version of the CDN edge behavior from Chapter 96.

**Key topics:** HTTP cache semantics in code, cache key design, invalidation and expiry.

### [Chapter 117: Building a Simple VPN](117-building-a-simple-vpn.md)

A minimal tunneling implementation using a TUN/TAP virtual interface: encapsulating one IP packet inside another (UDP) tunnel and encrypting the payload — the code-level realization of Chapter 85's tunneling concept.

**Key topics:** TUN/TAP interfaces, packet-in-packet encapsulation, adding encryption to a tunnel.

### [Chapter 118: Building a Distributed Network Service](118-building-a-distributed-network-service.md)

The capstone project for this volume: a small distributed key-value service with multiple nodes discovering each other over the network and communicating via a custom wire protocol — bringing together sockets, serialization, and the concurrency patterns from every prior chapter in this volume.

**Key topics:** custom wire protocol design, node discovery, distributed service architecture, putting the whole volume together.

---

# PART 18 — OBSERVABILITY & DEBUGGING

## Volume 18: Reasoning From Symptom to Root Cause

> A working knowledge of protocols is only half of network engineering. The other half is diagnosis: given a vague symptom ("the site is slow," "it doesn't connect"), knowing which layer to look at, which tool to reach for, and what a normal-versus-broken result actually looks like.

### [Chapter 119: Wireshark and tcpdump — Seeing the Network](119-wireshark-and-tcpdump.md)

Hands-on mastery of the two essential capture tools: `tcpdump` for fast command-line capture and filtering, Wireshark for deep visual inspection, and how to read a capture of a TCP handshake, a DNS query, and a TLS handshake — all protocols already taught — as actual bytes on the wire.

**Key topics:** tcpdump filters, Wireshark UI and filters, reading captured handshakes from earlier volumes.

### [Chapter 120: Measuring the Network — Latency, Throughput, Jitter, and Packet Loss](120-measuring-the-network.md)

Precise definitions (and how to measure each): latency, throughput vs. bandwidth (revisited from Chapter 16), jitter (variance in latency, critical for VoIP/video), and packet loss — with the tools that measure each (`ping`, `iperf`, MTR).

**Key topics:** latency vs bandwidth vs throughput, jitter, packet loss measurement, iperf and MTR.

### [Chapter 121: SNMP, Flow Logs, Prometheus, and Grafana for Networks](121-snmp-flow-logs-and-grafana.md)

How production network infrastructure is monitored continuously rather than debugged reactively: SNMP as the classic device-monitoring protocol, flow logs (NetFlow/VPC flow logs) for traffic visibility, and Prometheus/Grafana as the modern metrics-and-dashboards stack applied to network data.

**Key topics:** SNMP basics, flow logs, Prometheus metrics for networking, Grafana dashboards.

### [Chapter 122: The Debugging Playbook — From Symptom to Root Cause](122-the-debugging-playbook.md)

The methodology chapter: a decision framework for any networking symptom — which layer is most likely responsible, what to check first, and how to use the OSI/TCP-IP model (Volume 4) as a literal debugging checklist rather than just theory.

**Key topics:** layer-by-layer debugging framework, symptom-to-layer mapping, when to reach for which tool from Chapter 119–121.

### [Chapter 123: Real Debugging Scenarios, Solved](123-real-debugging-scenarios.md)

Nine realistic, fully worked scenarios applying the Chapter 122 playbook: DNS resolves but HTTP fails; TCP connects but the app hangs; a site works on Wi-Fi but not mobile data; one machine can't reach another; intermittent packet loss; sudden latency spikes; a service works internally but not externally; TLS handshake failures; and MTU-caused mysterious failures.

**Key topics:** nine worked debugging case studies, each traced from symptom through the correct layer to root cause.

---

# PART 19 — THE INTERNET AT GLOBAL SCALE

## Volume 19: How One Network Became *The* Network

> Individually, every protocol in this course is well understood. This volume is about what happens when you compose all of them across tens of thousands of independently-operated networks and call the result "the Internet."

### [Chapter 124: ISP Tiers, IXPs, Peering, and Transit — Revisited at Global Scale](124-isp-tiers-ixps-peering-and-transit.md)

Volume 7's peering/transit concepts, now viewed as a whole system: how Tier-1, Tier-2, and Tier-3 ISPs fit together, what a real Internet Exchange Point looks like operationally, and why "the Internet" is really a marketplace of interconnection agreements.

**Key topics:** the global ISP tier system as a whole, real IXP operations, interconnection as an economic system.

### [Chapter 125: Global Routing, Anycast, and CDN Architecture at Scale](125-global-routing-anycast-and-cdn-architecture.md)

How BGP (Volume 7) and Anycast (Volume 15) combine to let a company like Cloudflare or Google announce the same IP address from hundreds of locations worldwide, and how a global CDN's architecture routes each user to a nearby, healthy edge node.

**Key topics:** BGP Anycast at global scale, global CDN architecture, edge node selection.

### [Chapter 126: Submarine Cables and the Physical Backbone of the Internet](126-submarine-cables.md)

Revisiting Chapter 23's undersea cables with full operational detail: who owns them (increasingly, the cloud providers themselves), how they're laid and repaired, landing stations, and why a single cable cut can measurably slow down a whole region's Internet.

**Key topics:** submarine cable ownership and consortia, laying and repair, landing stations, real outage case studies.

### [Chapter 127: How Google, Amazon, Microsoft, and Cloudflare Run Networks at Planet Scale](127-how-hyperscalers-run-networks-at-scale.md)

A synthesis chapter built entirely from publicly documented architecture (with assumptions about undisclosed internals clearly labeled as such): private global backbones, edge points of presence, and how hyperscalers reduce reliance on the public Internet by building their own.

**Key topics:** private backbone networks, points of presence, publicly documented hyperscaler architecture, explicit labeling of what's undisclosed.

---

# PART 20 — THE ULTIMATE CAPSTONE

## Volume 20: One Question, Every Layer

> This is the question this entire course has been building toward. Answering it completely — and being able to explain every step — is the actual finish line.

### [Chapter 128: What EXACTLY Happens When You Type https://www.google.com and Press Enter?](128-what-happens-when-you-type-google-com.md)

The complete, layer-by-layer trace: keyboard input → browser → OS → DNS cache → resolver → DNS hierarchy → IP address → routing table → Wi-Fi → Ethernet → home router → NAT → ISP → BGP → Internet backbone → Google's network → Anycast → load balancer → server → TLS handshake → HTTP request → application logic → cache/database → response, all the way back to pixels on your screen — with every step naming the exact chapter of this course that explained it, and a final answer to: what physically happened to the bits, from keypress to pixels?

**Key topics:** full end-to-end trace across every layer and volume of this course, application/TLS/TCP-or-QUIC/IP/Ethernet/physical views of the same request, the unified mental model.

---

# PART 21 — FUTURE NETWORKING

## Volume 21: What Comes Next

> The final volume looks forward — carefully. Every topic here is explicitly labeled as deployed, commercially emerging, standardized, active research, or speculative, because conflating these is how good engineers end up making bad predictions.

### [Chapter 129: Satellite Internet, LEO Constellations, and Edge Computing](129-satellite-internet-leo-and-edge-computing.md)

Low Earth Orbit satellite constellations (Starlink, Kuiper) as a deployed alternative to terrestrial last-mile Internet, inter-satellite laser links as an emerging capability, and edge computing/IoT as the pattern of pushing compute closer to where data is generated — each labeled by deployment status.

**Key topics:** LEO satellite Internet, inter-satellite links, edge computing, industrial IoT networking, deployment-status labeling.

### [Chapter 130: Quantum Networking and the Quantum Internet](130-quantum-networking.md)

Quantum Key Distribution as a real (if niche and expensive) deployed technology today, versus the "quantum Internet" (networked quantum computers sharing entangled states) as active, early-stage research — explained intuitively without requiring quantum mechanics background, and clearly separated from science fiction.

**Key topics:** Quantum Key Distribution (deployed), quantum entanglement networking (research), realistic timeline honesty.

### [Chapter 131: 6G, AI-Native Networks, and the Speculative Frontier](131-6g-ai-native-networks-and-the-speculative-frontier.md)

A closing survey of what's actively being standardized or researched beyond this course's timeline: AI-native network management, autonomous self-healing networks, integrated sensing and communication, reconfigurable intelligent surfaces, and digital twins of networks — each explicitly classified, and a final reflection on the one durable skill this whole course was actually teaching: thinking like a networking engineer, at any layer, with any future technology.

**Key topics:** AI-native networking, autonomous networks, integrated sensing and communication, reconfigurable intelligent surfaces, digital twins, closing synthesis of the whole course.
