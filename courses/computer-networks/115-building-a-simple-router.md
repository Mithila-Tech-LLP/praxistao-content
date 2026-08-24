# Chapter 115: Building a Simple Router

> **"Chapter 45 proved longest-prefix-match on paper, with a routing table and a pencil. This chapter puts a real packet through it, decrements a real TTL, recomputes a real checksum, and sends the bytes somewhere else — the only difference between this and a rack-mounted router is scale, not logic."**

---

## Table of Contents

1. [Recap: Chapter 45's Algorithm, Chapter 114's Header Code](#1-recap-chapter-45s-algorithm-chapter-114s-header-code)
2. [The Problem: A Router Needs Interfaces, and We Don't Have Any](#2-the-problem-a-router-needs-interfaces-and-we-dont-have-any)
3. [The Simplification: UDP Sockets as Virtual Wires](#3-the-simplification-udp-sockets-as-virtual-wires)
4. [Designing the Toy Network](#4-designing-the-toy-network)
5. [Code: Building and Parsing an IPv4 Packet](#5-code-building-and-parsing-an-ipv4-packet)
6. [Code: The Routing Table and Longest Prefix Match](#6-code-the-routing-table-and-longest-prefix-match)
7. [Code: TTL Decrement and Checksum Recomputation](#7-code-ttl-decrement-and-checksum-recomputation)
8. [Code: The Forwarding Function](#8-code-the-forwarding-function)
9. [Code: Two Hosts and a Router, Wired Together](#9-code-two-hosts-and-a-router-wired-together)
10. [Code: main() — Running the Whole Network](#10-code-main--running-the-whole-network)
11. [Hands-On Experiment: Watch LPM Choose, Watch TTL Expire](#11-hands-on-experiment-watch-lpm-choose-watch-ttl-expire)
12. [Worked Example: Tracing One Packet Through the Code](#12-worked-example-tracing-one-packet-through-the-code)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes: What Real Routers Add](#14-production-notes-what-real-routers-add)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)

---

## 1. Recap: Chapter 45's Algorithm, Chapter 114's Header Code

Chapter 45 ended with a `bestMatch` function that correctly implements longest-prefix-match over a `RoutingTable` of `net.IPNet` entries, and a complete pseudocode algorithm (Section 10 of that chapter) describing exactly what a router does to every packet it forwards: check if the packet is locally destined, decrement TTL, drop and notify on expiry, run longest-prefix-match, resolve the next hop, recompute the checksum, and send the frame onward. Chapter 114 built byte-accurate parsing for the IPv4 header those algorithms operate on. This chapter combines both: a real routing table, a real forwarding loop, and real bytes moving between processes.

---

## 2. The Problem: A Router Needs Interfaces, and We Don't Have Any

A real router forwards a packet by receiving it on one physical network interface and transmitting it out a different one. Reproducing that exactly would mean either multiple physical NICs and real Ethernet segments, or — as Chapter 117 does properly — a TUN virtual interface that the OS treats as a real network device, both of which need root privileges and non-trivial setup before you've written a single line of forwarding logic.

That setup cost would bury the actual point of this chapter, which is Chapter 45's algorithm, not interface plumbing. So this chapter makes a deliberate trade: it fakes the "interface" concept with something anyone can run instantly, and saves the real TUN-based version for Chapter 117.

---

## 3. The Simplification: UDP Sockets as Virtual Wires

Each "interface" in this chapter is a **UDP socket bound to a fixed `127.0.0.1:port`**, and a "wire" connecting two interfaces is nothing more than both ends knowing each other's UDP address in advance. This is a closer analogy to real Ethernet than it might first appear: an Ethernet segment (Chapter 28) is just a shared medium where devices know how to address frames to each other by MAC; here, a loopback UDP socket pair plays exactly that role, minus the MAC address, ARP resolution, and physical signaling — all covered in earlier chapters and not the point of this one.

Every byte that crosses one of these UDP "wires" in this chapter is a **hand-built, complete IPv4 packet** — the same header format Chapter 114 parsed, now also being *constructed*. Nothing about the routing table, longest-prefix-match, TTL handling, or checksum recomputation is faked; only the physical transport between interfaces is a stand-in.

---

## 4. Designing the Toy Network

```
   Host A                      Router                       Host B
 10.0.1.10                 ifaceA      ifaceB               10.0.2.10
     |                    10.0.1.1    10.0.2.1                  |
     |                       |            |                     |
     +---- UDP :20001 <---> :20002   :20003 <---> UDP :20004 ---+
              "wire 1" (10.0.1.0/24)   "wire 2" (10.0.2.0/24)
```

| Device | Role | UDP address |
|---|---|---|
| Host A | sender, IP `10.0.1.10` | `127.0.0.1:20001` |
| Router, `ifaceA` | directly connected to `10.0.1.0/24` | `127.0.0.1:20002` |
| Router, `ifaceB` | directly connected to `10.0.2.0/24` | `127.0.0.1:20003` |
| Host B | receiver, IP `10.0.2.10` | `127.0.0.1:20004` |

The router's table also carries a **deliberately overlapping aggregate route**, `10.0.0.0/8`, pointing at a fictional, unwired "upstream" — mirroring Chapter 45, Section 4's worked example exactly, so a packet to `10.0.2.10` genuinely has two matching candidates and longest-prefix-match has real work to do, not just a single trivial match.

---

## 5. Code: Building and Parsing an IPv4 Packet

```go
package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

// buildIPv4 constructs a complete IPv4 header + payload, matching Chapter 36's
// field layout and Chapter 114 Section 7's parsing code exactly, byte for byte.
func buildIPv4(ttl, protocol uint8, src, dst net.IP, payload []byte) []byte {
	pkt := make([]byte, 20+len(payload))
	pkt[0] = 0x45 // version 4, IHL 5 (20-byte header, no options)
	binary.BigEndian.PutUint16(pkt[2:4], uint16(20+len(payload)))
	pkt[8] = ttl
	pkt[9] = protocol
	copy(pkt[12:16], src.To4())
	copy(pkt[16:20], dst.To4())
	copy(pkt[20:], payload)
	setChecksum(pkt) // Section 7 — must run AFTER every other field is final
	return pkt
}

type IPv4Header struct {
	IHL      uint8
	TTL      uint8
	Protocol uint8
	Checksum uint16
	Src, Dst net.IP
}

// parseIPv4 is Chapter 114 Section 7's parser, unchanged — a router reads
// packets with exactly the same code a sniffer does; the only difference is
// what happens next.
func parseIPv4(raw []byte) (IPv4Header, []byte, error) {
	if len(raw) < 20 {
		return IPv4Header{}, nil, fmt.Errorf("too short for IPv4: %d bytes", len(raw))
	}
	ihl := raw[0] & 0x0F
	headerLen := int(ihl) * 4
	if len(raw) < headerLen {
		return IPv4Header{}, nil, fmt.Errorf("truncated IPv4 header")
	}
	h := IPv4Header{
		IHL:      ihl,
		TTL:      raw[8],
		Protocol: raw[9],
		Checksum: binary.BigEndian.Uint16(raw[10:12]),
		Src:      net.IP(append([]byte{}, raw[12:16]...)),
		Dst:      net.IP(append([]byte{}, raw[16:20]...)),
	}
	return h, raw[headerLen:], nil
}
```

---

## 6. Code: The Routing Table and Longest Prefix Match

This is Chapter 45, Section 12's `bestMatch` function, unchanged in its core logic, with one addition: a `NextHopAddr` field standing in for what ARP (Chapter 53) would normally resolve — here, a static map from "next hop" to a UDP address, since there's no real Layer 2 to resolve against.

```go
type Route struct {
	Network     *net.IPNet
	Iface       string       // which "interface" (UDP socket) sends this out
	Connected   bool         // true = destination is directly reachable on this wire
	NextHopAddr *net.UDPAddr // where to send the bytes on this wire
}

type RoutingTable []Route

// bestMatch is Chapter 45 Section 12's algorithm, verbatim: among every route
// whose network contains dst, return the one with the longest prefix.
func (t RoutingTable) bestMatch(dst net.IP) (Route, bool) {
	var best Route
	found := false
	bestLen := -1
	for _, r := range t {
		if !r.Network.Contains(dst) {
			continue
		}
		ones, _ := r.Network.Mask.Size()
		if ones > bestLen {
			best = r
			bestLen = ones
			found = true
		}
	}
	return best, found
}
```

The table this chapter's router actually uses, matching Section 4's topology plus the deliberate `/8` overlap from Chapter 45:

```go
func buildRoutingTable(ifaceAAddr, ifaceBAddr *net.UDPAddr, hostBAddr *net.UDPAddr) RoutingTable {
	_, netA, _ := net.ParseCIDR("10.0.1.0/24")
	_, netB, _ := net.ParseCIDR("10.0.2.0/24")
	_, netAggregate, _ := net.ParseCIDR("10.0.0.0/8")

	return RoutingTable{
		// Directly connected: no next hop needed, destination is on this wire.
		{Network: netA, Iface: "ifaceA", Connected: true},
		// Connected route to host B's subnet — but the specific UDP address of
		// each host still needs a static "neighbor" lookup (Section 8), the toy
		// stand-in for ARP (Chapter 53).
		{Network: netB, Iface: "ifaceB", Connected: true, NextHopAddr: hostBAddr},
		// A broader, aggregate route that ALSO matches 10.0.2.10 — exactly
		// Chapter 45 Section 4's setup — pointed at a fictional upstream this
		// toy network never actually wires up.
		{Network: netAggregate, Iface: "upstream", Connected: false,
			NextHopAddr: &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 29999}},
	}
}
```

---

## 7. Code: TTL Decrement and Checksum Recomputation

```go
// ipChecksum computes the IPv4 header checksum: the one's-complement sum of
// every 16-bit word in the header, then one's-complemented. The checksum
// field itself must be treated as zero while computing this, which is why
// setChecksum below writes it as zero first (Chapter 45 Section 10's note
// that a changed TTL forces a checksum recompute every hop).
func ipChecksum(header []byte) uint16 {
	var sum uint32
	for i := 0; i+1 < len(header); i += 2 {
		sum += uint32(header[i])<<8 | uint32(header[i+1])
	}
	if len(header)%2 == 1 {
		sum += uint32(header[len(header)-1]) << 8
	}
	for sum>>16 != 0 {
		sum = (sum & 0xFFFF) + (sum >> 16)
	}
	return ^uint16(sum)
}

func setChecksum(pkt []byte) {
	ihl := int(pkt[0]&0x0F) * 4
	pkt[10], pkt[11] = 0, 0 // zero the checksum field before computing over it
	sum := ipChecksum(pkt[:ihl])
	binary.BigEndian.PutUint16(pkt[10:12], sum)
}

// decrementTTL implements Chapter 45 Section 9's rule exactly: subtract one,
// and report whether the packet just expired.
func decrementTTL(pkt []byte) (expired bool) {
	pkt[8]--
	return pkt[8] == 0
}
```

---

## 8. Code: The Forwarding Function

This is Chapter 45, Section 10's pseudocode, translated line for line:

```go
import "net"

type Interface struct {
	Name string
	IP   net.IP
	Conn *net.UDPConn
}

type Router struct {
	interfaces map[string]*Interface
	table      RoutingTable
}

func (r *Router) forward(pkt []byte) {
	hdr, _, err := parseIPv4(pkt)
	if err != nil {
		fmt.Println("[router] drop: malformed packet:", err)
		return
	}

	// Step 1: is this packet addressed to the router itself?
	for _, iface := range r.interfaces {
		if iface.IP.Equal(hdr.Dst) {
			fmt.Printf("[router] packet for %s is addressed to me (%s), delivering locally\n", hdr.Dst, iface.Name)
			return
		}
	}

	// Step 2: TTL (Chapter 45 Section 9).
	if expired := decrementTTL(pkt); expired {
		fmt.Printf("[router] TTL expired for packet to %s — dropping, would send ICMP Time Exceeded to %s\n", hdr.Dst, hdr.Src)
		return
	}

	// Step 3: longest-prefix-match (Chapter 45 Sections 3-5, Section 6 above).
	route, ok := r.table.bestMatch(hdr.Dst)
	if !ok {
		fmt.Printf("[router] no route to %s — dropping, would send ICMP Destination Unreachable\n", hdr.Dst)
		return
	}
	ones, _ := route.Network.Mask.Size()
	fmt.Printf("[router] %s matched %s (/%d) via %s\n", hdr.Dst, route.Network, ones, route.Iface)

	iface, ok := r.interfaces[route.Iface]
	if !ok {
		fmt.Printf("[router] route points at unreachable interface %q (not wired in this toy network) — dropping\n", route.Iface)
		return
	}

	// Step 4: recompute the checksum, since TTL just changed (Chapter 45 Section 10).
	setChecksum(pkt)

	// Step 5: resolve the next hop's transport address — our stand-in for ARP
	// (Chapter 53) — and send the bytes out the chosen "wire."
	dest := route.NextHopAddr
	if route.Connected && dest == nil {
		fmt.Printf("[router] no neighbor address known for %s on %s — dropping (real ARP would resolve this)\n", hdr.Dst, iface.Name)
		return
	}
	if _, err := iface.Conn.WriteToUDP(pkt, dest); err != nil {
		fmt.Println("[router] forward error:", err)
		return
	}
	fmt.Printf("[router] forwarded to %s via %s, new ttl=%d\n", hdr.Dst, iface.Name, pkt[8])
}
```

---

## 9. Code: Two Hosts and a Router, Wired Together

```go
func startHostB(addr *net.UDPAddr) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		panic(err)
	}
	buf := make([]byte, 65536)
	for {
		n, _, err := conn.ReadFromUDP(buf)
		if err != nil {
			continue
		}
		hdr, payload, _ := parseIPv4(buf[:n])
		fmt.Printf("[host B] received from %s, ttl=%d: %q\n", hdr.Src, hdr.TTL, string(payload))
	}
}

func sendFromHostA(routerIfaceA *net.UDPAddr, dst net.IP, ttl uint8, payload []byte) {
	conn, err := net.DialUDP("udp", nil, routerIfaceA)
	if err != nil {
		panic(err)
	}
	defer conn.Close()
	pkt := buildIPv4(ttl, 17 /* UDP, Chapter 58 */, net.ParseIP("10.0.1.10"), dst, payload)
	conn.Write(pkt)
	fmt.Printf("[host A] sent to %s, ttl=%d: %q\n", dst, ttl, string(payload))
}
```

---

## 10. Code: main() — Running the Whole Network

```go
func main() {
	ifaceAAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 20002}
	ifaceBAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 20003}
	hostBAddr := &net.UDPAddr{IP: net.ParseIP("127.0.0.1"), Port: 20004}

	connA, _ := net.ListenUDP("udp", ifaceAAddr)
	connB, _ := net.ListenUDP("udp", ifaceBAddr)

	router := &Router{
		interfaces: map[string]*Interface{
			"ifaceA": {Name: "ifaceA", IP: net.ParseIP("10.0.1.1"), Conn: connA},
			"ifaceB": {Name: "ifaceB", IP: net.ParseIP("10.0.2.1"), Conn: connB},
		},
		table: buildRoutingTable(ifaceAAddr, ifaceBAddr, hostBAddr),
	}

	// Router's receive loops — one goroutine per interface, exactly the
	// goroutine-per-connection pattern from Chapter 106.
	for _, iface := range router.interfaces {
		go func(iface *Interface) {
			buf := make([]byte, 65536)
			for {
				n, _, err := iface.Conn.ReadFromUDP(buf)
				if err != nil {
					continue
				}
				pkt := append([]byte{}, buf[:n]...)
				router.forward(pkt)
			}
		}(iface)
	}

	go startHostB(hostBAddr)

	time.Sleep(200 * time.Millisecond) // let listeners bind before traffic starts

	// Packet 1: a normal packet, ttl=4, plenty of hops left.
	sendFromHostA(ifaceAAddr, net.ParseIP("10.0.2.10"), 4, []byte("hello from host A"))
	time.Sleep(200 * time.Millisecond)

	// Packet 2: ttl=1 — will expire AT the router, demonstrating Chapter 45
	// Section 9's drop-and-notify behavior instead of a silent black hole.
	sendFromHostA(ifaceAAddr, net.ParseIP("10.0.2.10"), 1, []byte("this one should die at the router"))
	time.Sleep(200 * time.Millisecond)
}
```

---

## 11. Hands-On Experiment: Watch LPM Choose, Watch TTL Expire

```
$ go run .
[host A] sent to 10.0.2.10, ttl=4: "hello from host A"
[router] 10.0.2.10 matched 10.0.2.0/24 (/24) via ifaceB
[router] forwarded to 10.0.2.10 via ifaceB, new ttl=3
[host B] received from 10.0.1.10, ttl=3: "hello from host A"
[host A] sent to 10.0.2.10, ttl=1: "this one should die at the router"
[router] TTL expired for packet to 10.0.2.10 — dropping, would send ICMP Time Exceeded to 10.0.1.10
```

Notice the router's log line explicitly names both candidate routes' relationship: `10.0.2.10` matches *both* `10.0.2.0/24` and the `10.0.0.0/8` aggregate route from Section 4, and `bestMatch` picks the `/24` — the longer prefix — without ever considering table order, exactly Chapter 45 Section 3's rule. Try deleting the `10.0.2.0/24` route from `buildRoutingTable` and re-running: the packet now matches only the `/8`, and the forwarding log should show it being handed to the unwired `"upstream"` interface and dropped, since nothing is listening on port `29999` in this toy network.

---

## 12. Worked Example: Tracing One Packet Through the Code

Reproducing Chapter 45 Section 4's exact table format for the packet 10.0.1.10 → 10.0.2.10:

| # | Route | Prefix Length | Does `10.0.2.10` match? |
|---|---|---|---|
| 1 | `10.0.1.0/24` via ifaceA (connected) | 24 | No — wrong subnet |
| 2 | `10.0.2.0/24` via ifaceB (connected) | 24 | **Yes** |
| 3 | `10.0.0.0/8` via upstream (unwired) | 8 | Yes |

Both Route 2 and Route 3 match. `bestMatch` compares prefix lengths — 24 beats 8 — and Route 2 wins, exactly as Section 11's captured output showed. If Route 2 didn't exist, Route 3 would win by being the only match, and the packet would be handed to an interface this toy network never actually wires up, printing the "not wired" drop message instead of silently vanishing — a direct, deliberate echo of Chapter 45 Section 10's "no route" and "route to nowhere" branches.

---

## 13. Common Misconceptions

- **"A router looks at the payload or port numbers to decide where to forward a packet."** No — everything in Section 8's `forward` function reads only the IP header (destination address, TTL). This is the sharp contrast with Chapter 113's Layer-7 load balancer, which deliberately does inspect application data to make its routing decision; a plain IP router (Layer 3) never needs to.
- **"The next-hop field in a route always points at the final destination."** As Chapter 45 Section 7 explained, it points at the next *directly reachable* router or host — Route 2 in Section 12's table stores Host B's own UDP address only because Host B happens to be directly attached to that "wire"; for a multi-hop path, it would instead be the next router's address, which would then run its own separate `bestMatch` lookup.
- **"TTL expiring means the packet is corrupted."** It means only that it has traveled its allotted number of hops (Chapter 45 Section 9) — nothing about the packet's *contents* is wrong; the field exists purely to bound how long a misrouted packet can circulate.
- **"Checksum recomputation is optional if the payload didn't change."** The IPv4 header checksum covers the *header*, not the payload, and TTL — a header field — changes on every single hop, which is precisely why Section 7's `setChecksum` must run after every `decrementTTL` call, unconditionally.

---

## 14. Production Notes: What Real Routers Add

- **Real routers don't use ARP-style resolution as a static map** — Chapter 53's actual ARP protocol resolves next-hop IPs to MAC addresses dynamically, with a cache and timeout, not a hardcoded lookup table baked in at startup like this chapter's `NextHopAddr` fields.
- **Real forwarding tables aren't scanned linearly** — as Chapter 44 Section 9 and Chapter 45 Section 5 both noted, production routers use tries or TCAM hardware lookups that produce the same longest-prefix-match answer as `bestMatch`'s loop, just in a handful of clock cycles regardless of table size, which matters enormously at millions of routes and millions of packets per second.
- **Linux itself is a software router.** Setting `sysctl net.ipv4.ip_forward=1` turns any Linux box with multiple real interfaces into exactly this chapter's router, using the kernel's own routing table (`ip route`, Chapter 45 Section 11) and forwarding path instead of hand-written Go — Chapter 102 covers this stack in depth.
- **Real ICMP replies are real packets, not log lines.** Section 8's TTL-expiry and no-route branches print a message where a production router would construct and send an actual ICMP Time Exceeded or Destination Unreachable packet (Chapter 54) back toward the original source.
- **A production router forwards millions of packets per second across many interfaces simultaneously**; this chapter's single-process, two-interface, artificially delayed demo is built for clarity, not throughput.

---

## 15. What's Simplified Here

- UDP sockets on loopback stand in for real Ethernet interfaces and physical wires — Chapter 117 replaces this with a real TUN interface carrying genuine OS traffic.
- Next-hop resolution is a static map, not real ARP (Chapter 53).
- ICMP Time Exceeded and Destination Unreachable are logged, not actually constructed and sent (Chapter 54 covers their real format).
- IPv4 Options (IHL > 5) and fragmentation are not handled — every packet here is assumed to be a single, unfragmented, option-free datagram.
- The routing table is static and hardcoded at startup; Chapter 46 through 49 cover how real tables are built and kept current dynamically.
- There is no administrative-distance or metric tie-breaking (Chapter 45 Section 13) because this toy topology never produces two routes with an identical prefix length.

---

## 16. Interview Questions & Model Answers

**Beginner: What two things must a router do to a packet's IP header on every single hop it forwards, regardless of anything else?**
Decrement the TTL field by one, and recompute the header checksum — the checksum must change because it covers the header, and TTL (part of the header) just changed.

**Intermediate: Why does this chapter's router check whether an incoming packet's destination matches one of its own interface IPs before running longest-prefix-match at all?**
Because a packet addressed directly to the router itself should be delivered locally, not forwarded onward — forwarding it would be both wrong (the router is the actual destination, not a hop-through point) and potentially create a loop if a route also happened to match that address.

**Advanced: This chapter's router needs no root privileges and runs on any OS, unlike Chapter 114's raw-socket sniffer. Explain precisely why, referencing what a UDP socket does and doesn't require, and what would have to change to make this a "real" router.**
A UDP socket is an entirely ordinary, unprivileged socket — `net.ListenUDP` doesn't ask the kernel for anything beyond a normal port binding, and every "packet" this chapter forwards is just UDP payload data the application itself constructed and interpreted; the OS never treats these bytes as real network-layer traffic at all. Making this a genuine router would require capturing and injecting real IP packets on real interfaces — either raw sockets (Chapter 114's approach, privileged and OS-specific) or a TUN device (Chapter 117), both of which require elevated privileges precisely because they let a process see or inject traffic outside its own connections, which Section 2 of Chapter 114 explained is exactly what the ordinary socket API is designed to prevent.

---

## 17. Exercises

### Easy
1. Add a third host, Host C, on a new `10.0.3.0/24` subnet behind a new router interface, and confirm packets route correctly to all three subnets.
2. Change Host A's initial TTL to 2 instead of 4 and trace by hand how many hops it can survive before expiring, given this chapter's one-router topology.
3. Explain in your own words why the checksum field must be zeroed before `ipChecksum` runs over the header.

### Medium
4. Implement an actual ICMP-Time-Exceeded-style reply: when `forward` detects TTL expiry, construct and send a small UDP packet back to the original source's address describing the drop, instead of only printing a log line.
5. Add a second route to `10.0.2.0/25` (a more specific sub-split of Host B's subnet) and verify `bestMatch` prefers it over the existing `/24` for addresses that fall within it.
6. Modify the router to reject and drop any packet whose header checksum (as received) doesn't match a freshly computed one — simulating corruption detection.

### Hard
7. Chain two routers together (Router 1 connects Host A's subnet to a middle "backbone" subnet; Router 2 connects that backbone to Host B's subnet) and confirm TTL decrements exactly twice end to end.
8. Implement dynamic next-hop resolution: replace the static `NextHopAddr` map with a simple "ARP-like" broadcast-and-reply exchange over an additional UDP socket per wire, resolving an IP to a UDP address on first use and caching the result.
9. Introduce a deliberate routing loop between two chained routers (each with a route pointing back at the other for some destination) and confirm a packet's TTL still forces it to be dropped within a bounded number of hops, exactly as Chapter 45 Section 9 describes.

---

## 18. Summary

| Term | Meaning |
|---|---|
| Virtual interface (this chapter) | A UDP socket bound to a fixed loopback port, standing in for a physical NIC |
| `bestMatch` | Chapter 45's longest-prefix-match function, reused verbatim against `net.IPNet` routes |
| Connected route | A route whose destination network is directly reachable on one of the router's own interfaces |
| `setChecksum` / `ipChecksum` | Functions recomputing the IPv4 header checksum after any header field changes |
| `decrementTTL` | Implements Chapter 45 Section 9's hop-counter rule and expiry detection |
| Static neighbor map | This chapter's stand-in for ARP (Chapter 53) — a hardcoded IP-to-transport-address table |
| Forwarding loop | The per-interface goroutine reading incoming bytes and handing them to `forward` |

Chapter 115 turned Chapter 45's paper algorithm into a running program that genuinely forwards, decrements, and recomputes. Chapter 116 leaves routing behind and builds the other end of the spectrum this volume has been building toward: an HTTP cache that decides, correctly, when it's allowed to reuse a response at all — turning Chapter 72's Cache-Control and ETag rules into code a real client can hit and get a real `304 Not Modified` from. Chapters 117 (a real TUN-based VPN tunnel) and 118 (a distributed key-value service) close out this volume after that.
