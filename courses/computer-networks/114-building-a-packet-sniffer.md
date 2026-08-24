# Chapter 114: Building a Packet Sniffer

> **"Every chapter so far asked you to trust a diagram of a header. This chapter asks you to catch the bytes yourself, straight off a real interface, and prove the diagram was telling the truth all along."**

---

## Table of Contents

1. [Recap: What Chapters 106–113 Quietly Assumed](#1-recap-what-chapters-106113-quietly-assumed)
2. [The Problem: Sockets Hide the Wire From You, On Purpose](#2-the-problem-sockets-hide-the-wire-from-you-on-purpose)
3. [Three Ways to See Raw Frames](#3-three-ways-to-see-raw-frames)
4. [Requirements, Permissions, and the Fallback Path](#4-requirements-permissions-and-the-fallback-path)
5. [The Shapes We're Decoding, Recapped](#5-the-shapes-were-decoding-recapped)
6. [Code: Parsing an Ethernet Frame By Hand](#6-code-parsing-an-ethernet-frame-by-hand)
7. [Code: Parsing an IPv4 Header By Hand](#7-code-parsing-an-ipv4-header-by-hand)
8. [Code: Parsing TCP and UDP By Hand](#8-code-parsing-tcp-and-udp-by-hand)
9. [Code: A Raw AF_PACKET Capture Loop (Linux, Needs Root)](#9-code-a-raw-af_packet-capture-loop-linux-needs-root)
10. [Code: A Fallback Simulation Mode (No Root, No Linux Required)](#10-code-a-fallback-simulation-mode-no-root-no-linux-required)
11. [Code: Wiring It Together — main()](#11-code-wiring-it-together--main)
12. [Hands-On Experiment: Capture Real Traffic](#12-hands-on-experiment-capture-real-traffic)
13. [Worked Example: One Captured SYN, Byte by Byte](#13-worked-example-one-captured-syn-byte-by-byte)
14. [Common Misconceptions](#14-common-misconceptions)
15. [Production Notes: tcpdump, Wireshark, gopacket, eBPF](#15-production-notes-tcpdump-wireshark-gopacket-ebpf)
16. [What's Simplified Here](#16-whats-simplified-here)
17. [Interview Questions & Model Answers](#17-interview-questions--model-answers)
18. [Exercises](#18-exercises)
19. [Summary](#19-summary)

---

## 1. Recap: What Chapters 106–113 Quietly Assumed

Every program built in Chapters 106 through 113 called `net.Dial` or `net.Listen` and then read and wrote clean, already-unwrapped bytes. `net.Conn.Read` on a TCP connection never handed back an Ethernet header, an IP header, or a TCP header — the kernel silently stripped all three away before your code ever saw a byte, and quietly demultiplexed the segment to your process based purely on port number (Chapter 57). That's precisely correct behavior for an application, and it's why chapters 106–113 could focus entirely on the protocols they were building without ever touching a header field directly.

This chapter breaks that abstraction deliberately. Instead of asking the kernel to hand you the payload of a connection you already know about, you ask it to hand you **every frame arriving on an interface**, headers and all, before any of that stripping happens — turning Chapter 28's Ethernet frame diagram, Chapter 36's IPv4 header, and Chapter 65's TCP header from pictures in a book into bytes you parse yourself.

---

## 2. The Problem: Sockets Hide the Wire From You, On Purpose

Try this thought experiment with the tools you already have: could Chapter 106's TCP server, using only `net.Listen("tcp", ...)`, tell you the source MAC address of a client that just connected? No — that information was consumed and discarded by the network stack before the `Accept()` call ever returned a `net.Conn`. Could it see a *different* connection's traffic, one destined for a different port entirely? Also no — the socket API is built around exactly one guarantee: "give me traffic for this one port/connection, and nothing else." That guarantee is a feature for almost every program ever written, and it is exactly what a sniffer needs to bypass.

A packet sniffer needs three things ordinary sockets refuse to provide:

- **Every frame on the wire**, not just frames matching one socket's port/connection tuple.
- **The headers themselves**, not just the payload above them.
- Frames **not addressed to this machine at all** (on a hub or a mirrored switch port — Section 14 covers exactly what a switched network does and doesn't allow here).

Getting all three requires a fundamentally different kind of socket — one that operates below the layer where the kernel normally starts making decisions on your behalf.

---

## 3. Three Ways to See Raw Frames

- **Raw sockets, `AF_PACKET` / `SOCK_RAW` (Linux only).** A privileged socket type that hands your process a copy of every frame the kernel sees on a given interface, starting from the Ethernet header, before any IP or TCP/UDP processing happens. This is the lowest-level, most direct option, and what Section 9 builds.
- **`libpcap` (and Go's `gopacket/pcap` wrapper around it).** A portable capture library — the same one `tcpdump` and Wireshark (Chapter 119) are built on — that abstracts over `AF_PACKET` on Linux, BPF devices (`/dev/bpf*`) on macOS/BSD, and `Npcap`/`WinPcap` on Windows, so the same capture code runs everywhere. It also supports kernel-level BPF filter expressions (`tcp port 443`) that discard unwanted frames before they're even copied into your process.
- **eBPF/XDP (Chapter 105).** The modern, highest-performance option: a program is loaded *into* the kernel (or even the NIC driver) and runs against every frame at line rate, only copying to userspace what's actually needed. This is what production traffic-analysis tools increasingly use instead of a userspace capture loop at all.

This chapter builds the `AF_PACKET` version directly, using nothing but Go's standard `syscall` package — no external library, no cgo, no libpcap installation required. It's Linux-only and needs elevated privileges, which is exactly why Section 4 also builds a fallback that needs neither.

---

## 4. Requirements, Permissions, and the Fallback Path

The raw-capture code in Section 9 needs:

- **Linux.** `AF_PACKET` is a Linux-specific address family; it does not exist on macOS or Windows.
- **Root, or `CAP_NET_RAW`.** Opening an `AF_PACKET`/`SOCK_RAW` socket is a privileged operation — run it with `sudo`, or grant the compiled binary the capability directly: `sudo setcap cap_net_raw+ep ./sniffer`.
- **A real interface with traffic on it.** `lo` (loopback) works for self-generated traffic; a real NIC (`eth0`, `wlan0`, etc.) sees real traffic, but on a modern switched network (Section 14) that traffic is overwhelmingly limited to frames addressed to this machine unless the switch is explicitly mirroring a port.

If any of that isn't available — no root, no Linux, no interesting traffic nearby — Section 10 builds a **simulation mode** that runs the exact same parsing code (Sections 6–8) against a small set of hand-built, byte-correct frames, so every reader can run the whole chapter's logic with nothing but `go run`.

---

## 5. The Shapes We're Decoding, Recapped

Three header formats, already fully specified in earlier chapters, condensed here as a byte-offset reference:

```
Ethernet (Chapter 28) — 14 bytes before the payload:
  bytes 0-5    Destination MAC
  bytes 6-11   Source MAC
  bytes 12-13  EtherType   (0x0800 = IPv4, 0x0806 = ARP, 0x86DD = IPv6)

IPv4 (Chapter 36) — 20 bytes minimum:
  byte  0      Version (4 bits) | IHL (4 bits, in 32-bit words)
  bytes 2-3    Total Length
  byte  8      TTL
  byte  9      Protocol   (1 = ICMP, 6 = TCP, 17 = UDP)
  bytes 10-11  Header Checksum
  bytes 12-15  Source IP
  bytes 16-19  Destination IP

TCP (Chapter 65) — 20 bytes minimum:
  bytes 0-1    Source Port
  bytes 2-3    Destination Port
  bytes 4-7    Sequence Number
  bytes 8-11   Acknowledgment Number
  byte  12     Data Offset (upper 4 bits, in 32-bit words)
  byte  13     Flags (CWR ECE URG ACK PSH RST SYN FIN)
  bytes 14-15  Window Size

UDP (Chapter 58) — 8 bytes, fixed:
  bytes 0-1    Source Port
  bytes 2-3    Destination Port
  bytes 4-5    Length
  bytes 6-7    Checksum
```

The code below does nothing but read these exact byte ranges out of a `[]byte` slice — there is no magic here, only arithmetic on offsets you've already seen drawn as diagrams.

---

## 6. Code: Parsing an Ethernet Frame By Hand

```go
package main

import (
	"encoding/binary"
	"fmt"
	"net"
)

const (
	EtherTypeIPv4 = 0x0800
	EtherTypeARP  = 0x0806
	EtherTypeIPv6 = 0x86DD
)

type EthernetFrame struct {
	DstMAC, SrcMAC net.HardwareAddr
	EtherType      uint16
	Payload        []byte
}

// parseEthernet reads the 14-byte Ethernet header (Chapter 28, Section 4) directly
// out of raw bytes handed back by the capture socket. Note there is no Preamble,
// SFD, or FCS here — the NIC hardware already stripped those before the kernel
// ever saw the frame (Chapter 28, Section 4's note on what's actually visible
// to software).
func parseEthernet(raw []byte) (EthernetFrame, error) {
	if len(raw) < 14 {
		return EthernetFrame{}, fmt.Errorf("frame too short for an Ethernet header: %d bytes", len(raw))
	}
	return EthernetFrame{
		DstMAC:    net.HardwareAddr(append([]byte{}, raw[0:6]...)),
		SrcMAC:    net.HardwareAddr(append([]byte{}, raw[6:12]...)),
		EtherType: binary.BigEndian.Uint16(raw[12:14]),
		Payload:   raw[14:],
	}, nil
}

func etherTypeName(t uint16) string {
	switch t {
	case EtherTypeIPv4:
		return "IPv4"
	case EtherTypeARP:
		return "ARP"
	case EtherTypeIPv6:
		return "IPv6"
	default:
		return fmt.Sprintf("0x%04x", t)
	}
}
```

---

## 7. Code: Parsing an IPv4 Header By Hand

```go
const (
	ProtoICMP = 1
	ProtoTCP  = 6
	ProtoUDP  = 17
)

type IPv4Header struct {
	Version  uint8
	IHL      uint8 // header length, in 32-bit words
	TotalLen uint16
	TTL      uint8
	Protocol uint8
	Checksum uint16
	Src, Dst net.IP
	Payload  []byte
}

// parseIPv4 mirrors Chapter 36's field-by-field layout exactly. IHL matters
// because the header can legally be longer than 20 bytes if Options are
// present (Chapter 65, Section 6 covers the identical Data-Offset idea for TCP) —
// we must read IHL before we know where the header actually ends.
func parseIPv4(raw []byte) (IPv4Header, error) {
	if len(raw) < 20 {
		return IPv4Header{}, fmt.Errorf("too short for an IPv4 header: %d bytes", len(raw))
	}
	verIHL := raw[0]
	version := verIHL >> 4
	ihl := verIHL & 0x0F
	headerLen := int(ihl) * 4
	if version != 4 {
		return IPv4Header{}, fmt.Errorf("not IPv4: version nibble was %d", version)
	}
	if len(raw) < headerLen {
		return IPv4Header{}, fmt.Errorf("IPv4 header claims %d bytes but only %d available", headerLen, len(raw))
	}
	return IPv4Header{
		Version:  version,
		IHL:      ihl,
		TotalLen: binary.BigEndian.Uint16(raw[2:4]),
		TTL:      raw[8],
		Protocol: raw[9],
		Checksum: binary.BigEndian.Uint16(raw[10:12]),
		Src:      net.IP(append([]byte{}, raw[12:16]...)),
		Dst:      net.IP(append([]byte{}, raw[16:20]...)),
		Payload:  raw[headerLen:],
	}, nil
}

func protocolName(p uint8) string {
	switch p {
	case ProtoICMP:
		return "ICMP"
	case ProtoTCP:
		return "TCP"
	case ProtoUDP:
		return "UDP"
	default:
		return fmt.Sprintf("proto-%d", p)
	}
}
```

---

## 8. Code: Parsing TCP and UDP By Hand

```go
import "strings"

type TCPHeader struct {
	SrcPort, DstPort uint16
	SeqNum, AckNum   uint32
	DataOffset       uint8 // in 32-bit words
	Flags            uint8
	Window           uint16
	Payload          []byte
}

// parseTCP mirrors Chapter 65, Section 2's byte-offset table exactly.
func parseTCP(raw []byte) (TCPHeader, error) {
	if len(raw) < 20 {
		return TCPHeader{}, fmt.Errorf("too short for a TCP header: %d bytes", len(raw))
	}
	dataOffset := raw[12] >> 4
	headerLen := int(dataOffset) * 4
	if len(raw) < headerLen {
		return TCPHeader{}, fmt.Errorf("TCP header claims %d bytes but only %d available", headerLen, len(raw))
	}
	return TCPHeader{
		SrcPort:    binary.BigEndian.Uint16(raw[0:2]),
		DstPort:    binary.BigEndian.Uint16(raw[2:4]),
		SeqNum:     binary.BigEndian.Uint32(raw[4:8]),
		AckNum:     binary.BigEndian.Uint32(raw[8:12]),
		DataOffset: dataOffset,
		Flags:      raw[13],
		Window:     binary.BigEndian.Uint16(raw[14:16]),
		Payload:    raw[headerLen:],
	}, nil
}

// tcpFlagsString decodes the flags byte (Chapter 65, Section 7) into a readable label.
func tcpFlagsString(flags uint8) string {
	var b strings.Builder
	add := func(bit uint8, name string) {
		if flags&bit != 0 {
			if b.Len() > 0 {
				b.WriteByte(' ')
			}
			b.WriteString(name)
		}
	}
	add(0x20, "URG")
	add(0x10, "ACK")
	add(0x08, "PSH")
	add(0x04, "RST")
	add(0x02, "SYN")
	add(0x01, "FIN")
	if b.Len() == 0 {
		return "-"
	}
	return b.String()
}

type UDPHeader struct {
	SrcPort, DstPort uint16
	Length           uint16
	Checksum         uint16
	Payload          []byte
}

// parseUDP mirrors Chapter 58's fixed 8-byte header — no Data Offset needed,
// because UDP's header never has options (that's the whole contrast Chapter 58
// draws against TCP).
func parseUDP(raw []byte) (UDPHeader, error) {
	if len(raw) < 8 {
		return UDPHeader{}, fmt.Errorf("too short for a UDP header: %d bytes", len(raw))
	}
	return UDPHeader{
		SrcPort:  binary.BigEndian.Uint16(raw[0:2]),
		DstPort:  binary.BigEndian.Uint16(raw[2:4]),
		Length:   binary.BigEndian.Uint16(raw[4:6]),
		Checksum: binary.BigEndian.Uint16(raw[6:8]),
		Payload:  raw[8:],
	}, nil
}
```

With Sections 6–8 in place, one function ties the layers together exactly the way Chapter 27's encapsulation model says they nest:

```go
func decodeAndPrint(raw []byte) {
	eth, err := parseEthernet(raw)
	if err != nil {
		fmt.Println("  [error]", err)
		return
	}
	fmt.Printf("ETH  %s -> %s  type=%s\n", eth.SrcMAC, eth.DstMAC, etherTypeName(eth.EtherType))

	if eth.EtherType != EtherTypeIPv4 {
		return // ARP/IPv6 decode omitted here — see Exercises
	}
	ip, err := parseIPv4(eth.Payload)
	if err != nil {
		fmt.Println("  [error]", err)
		return
	}
	fmt.Printf("  IP   %s -> %s  ttl=%d  proto=%s  len=%d\n",
		ip.Src, ip.Dst, ip.TTL, protocolName(ip.Protocol), ip.TotalLen)

	switch ip.Protocol {
	case ProtoTCP:
		tcp, err := parseTCP(ip.Payload)
		if err != nil {
			fmt.Println("    [error]", err)
			return
		}
		fmt.Printf("    TCP  %d -> %d  flags=[%s]  seq=%d ack=%d win=%d\n",
			tcp.SrcPort, tcp.DstPort, tcpFlagsString(tcp.Flags), tcp.SeqNum, tcp.AckNum, tcp.Window)
	case ProtoUDP:
		udp, err := parseUDP(ip.Payload)
		if err != nil {
			fmt.Println("    [error]", err)
			return
		}
		fmt.Printf("    UDP  %d -> %d  len=%d\n", udp.SrcPort, udp.DstPort, udp.Length)
	}
}
```

---

## 9. Code: A Raw AF_PACKET Capture Loop (Linux, Needs Root)

```go
//go:build linux

package main

import (
	"fmt"
	"net"
	"os"
	"syscall"
)

// htons converts a 16-bit value from host byte order to network byte order.
// This is the single most common bug in hand-written AF_PACKET code: forgetting
// this conversion means the socket ends up bound to the wrong (or a
// nonexistent) protocol number and silently captures nothing.
func htons(i uint16) uint16 {
	return (i << 8 & 0xff00) | (i >> 8)
}

func captureLive(ifaceName string) error {
	fd, err := syscall.Socket(syscall.AF_PACKET, syscall.SOCK_RAW, int(htons(syscall.ETH_P_ALL)))
	if err != nil {
		return fmt.Errorf("opening AF_PACKET socket (are you root? try sudo): %w", err)
	}
	defer syscall.Close(fd)

	iface, err := net.InterfaceByName(ifaceName)
	if err != nil {
		return fmt.Errorf("interface %q not found: %w", ifaceName, err)
	}

	addr := syscall.SockaddrLinklayer{
		Protocol: htons(syscall.ETH_P_ALL),
		Ifindex:  iface.Index,
	}
	if err := syscall.Bind(fd, &addr); err != nil {
		return fmt.Errorf("binding to %s: %w", ifaceName, err)
	}

	fmt.Fprintf(os.Stderr, "listening on %s (Ctrl-C to stop)...\n", ifaceName)
	buf := make([]byte, 65536)
	for {
		n, _, err := syscall.Recvfrom(fd, buf, 0)
		if err != nil {
			fmt.Fprintln(os.Stderr, "recvfrom:", err)
			continue
		}
		decodeAndPrint(buf[:n])
	}
}
```

This is the entire mechanism: `AF_PACKET` + `SOCK_RAW`, bound with `ETH_P_ALL` (accept every EtherType, not just IPv4), gets every frame the interface sees handed straight to `Recvfrom` — headers, payload, and all — with nothing stripped. Everything from Section 6 onward runs against exactly the bytes the NIC driver received.

---

## 10. Code: A Fallback Simulation Mode (No Root, No Linux Required)

For readers on macOS, Windows, or without root, this builds a handful of frames byte-by-byte — using the *inverse* of Sections 6–8's parsing logic — and feeds them through the exact same `decodeAndPrint`.

```go
func buildIPv4(ttl, proto uint8, src, dst net.IP, payload []byte) []byte {
	header := make([]byte, 20)
	header[0] = 0x45 // version 4, IHL 5 (20 bytes, no options)
	binary.BigEndian.PutUint16(header[2:4], uint16(20+len(payload)))
	header[8] = ttl
	header[9] = proto
	copy(header[12:16], src.To4())
	copy(header[16:20], dst.To4())
	// Header checksum left as 0 here for simplicity — a real NIC/stack computes
	// this over the header only; see Chapter 115, Section 7 for the algorithm.
	return append(header, payload...)
}

func buildTCPSyn(srcPort, dstPort uint16) []byte {
	header := make([]byte, 20)
	binary.BigEndian.PutUint16(header[0:2], srcPort)
	binary.BigEndian.PutUint16(header[2:4], dstPort)
	binary.BigEndian.PutUint32(header[4:8], 1000) // arbitrary ISN, Chapter 59
	header[12] = 5 << 4                           // data offset 5 words = 20 bytes, no options
	header[13] = 0x02                             // SYN flag only
	binary.BigEndian.PutUint16(header[14:16], 65535)
	return header
}

func buildUDP(srcPort, dstPort uint16, payload []byte) []byte {
	header := make([]byte, 8)
	binary.BigEndian.PutUint16(header[0:2], srcPort)
	binary.BigEndian.PutUint16(header[2:4], dstPort)
	binary.BigEndian.PutUint16(header[4:6], uint16(8+len(payload)))
	return append(header, payload...)
}

func buildEthernet(dst, src net.HardwareAddr, etherType uint16, payload []byte) []byte {
	frame := make([]byte, 14)
	copy(frame[0:6], dst)
	copy(frame[6:12], src)
	binary.BigEndian.PutUint16(frame[12:14], etherType)
	return append(frame, payload...)
}

func simulatedCapture() [][]byte {
	clientMAC := net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x01}
	serverMAC := net.HardwareAddr{0x02, 0x00, 0x00, 0x00, 0x00, 0x02}
	clientIP := net.ParseIP("192.0.2.10")
	serverIP := net.ParseIP("192.0.2.20")

	synSegment := buildTCPSyn(51234, 443) // a client opening a TLS connection
	frame1 := buildEthernet(serverMAC, clientMAC, EtherTypeIPv4,
		buildIPv4(64, ProtoTCP, clientIP, serverIP, synSegment))

	dnsQuery := buildUDP(53421, 53, []byte("(pretend DNS query bytes, Chapter 111)"))
	frame2 := buildEthernet(serverMAC, clientMAC, EtherTypeIPv4,
		buildIPv4(64, ProtoUDP, clientIP, serverIP, dnsQuery))

	return [][]byte{frame1, frame2}
}
```

---

## 11. Code: Wiring It Together — main()

```go
func main() {
	if len(os.Args) < 2 || os.Args[1] == "-simulate" {
		fmt.Println("=== simulation mode: no root, no live interface required ===")
		for _, frame := range simulatedCapture() {
			decodeAndPrint(frame)
		}
		return
	}
	if err := captureLive(os.Args[1]); err != nil {
		fmt.Fprintln(os.Stderr, "error:", err)
		os.Exit(1)
	}
}
```

Run it two ways:

```
$ go run .                       # simulation mode, works anywhere
$ sudo go run . eth0             # live capture, Linux + root required
```

---

## 12. Hands-On Experiment: Capture Real Traffic

On a Linux machine (a cloud VM is fine), run the live capture against loopback while generating your own traffic in a second terminal:

```
# terminal 1
$ sudo go run . lo

# terminal 2
$ curl http://localhost:80/    # or any local service; even a connection refused works
```

You should see, printed by your own code, the exact TCP three-way handshake (Chapter 59) your `curl` process just triggered: a SYN from an ephemeral port (Chapter 57) to port 80, decoded flag-by-flag as `[SYN]`, followed by whatever response comes back. If nothing is listening on port 80, you'll instead see a SYN followed by a TCP segment with the `RST` flag set — this is the byte-level view of "connection refused" that Chapter 59 described only in prose.

---

## 13. Worked Example: One Captured SYN, Byte by Byte

Take the `frame1` built in Section 10 and trace it exactly as `decodeAndPrint` would, field by field:

```
Ethernet header (14 bytes):
  Dst MAC: 02:00:00:00:00:02
  Src MAC: 02:00:00:00:00:01
  EtherType: 0x0800 -> IPv4

IPv4 header (20 bytes, IHL=5):
  Version/IHL byte: 0x45 -> version 4, IHL 5 words = 20 bytes
  TTL: 64
  Protocol: 6 -> TCP
  Src: 192.0.2.10
  Dst: 192.0.2.20

TCP header (20 bytes, Data Offset=5):
  Src Port: 51234 (an ephemeral port, Chapter 57)
  Dst Port: 443 (HTTPS, Chapter 82)
  Seq: 1000
  Flags byte: 0x02 -> SYN only
  Window: 65535
```

This is precisely Chapter 59's opening handshake segment — a client at an ephemeral source port initiating a connection to a well-known destination port with only the SYN flag set — except now every field's exact byte position and value is something you extracted yourself, not something you were told to trust.

---

## 14. Common Misconceptions

- **"A sniffer on my laptop sees everyone's traffic on the office WiFi."** On any modern switched or WPA2/WPA3-encrypted network, no — a switch (Chapter 30) forwards each frame only out the port the destination MAC lives on, and it takes a mirrored/SPAN port, a hub, or a successful ARP-spoofing attack (Chapter 83) for a sniffer to see traffic that isn't already addressed to or from its own machine. "Promiscuous mode" changes what the *NIC* will hand up to the OS (normally a NIC silently drops frames not addressed to it or to broadcast/multicast); it does nothing to make a switch deliver frames it wouldn't otherwise send to that port.
- **"Root is required because sniffing is inherently dangerous to the network."** It's required because `AF_PACKET`/`SOCK_RAW` sockets bypass the normal per-connection access control the kernel enforces for every other socket type (Chapter 57) — the privilege check is about what your *process* can see, not about any effect on the wire itself. Pure capture is passive and changes nothing about the traffic.
- **"If I can capture a TLS-protected connection's packets, I can read its contents."** You can read every header at every layer this chapter decodes — MACs, IPs, ports, sequence numbers, flags — but Chapter 82's TLS handshake means the TCP/UDP *payload* itself is ciphertext, unreadable without the session keys. A sniffer is a perfect tool for header/metadata analysis and a useless one for reading encrypted payloads.
- **"The FCS in a captured frame lets me verify checksums."** As Chapter 28, Section 8 noted, the NIC hardware validates and then strips the FCS before handing a frame to software in almost all cases — a userspace capture essentially never sees a real FCS value, correct or corrupted.

---

## 15. Production Notes: tcpdump, Wireshark, gopacket, eBPF

- **`tcpdump` and Wireshark** (Chapter 119 in full) are libpcap-based tools that do everything this chapter's code does, plus protocol dissection for hundreds of higher-level protocols, filter expressions, and a GUI — reach for them for real diagnostic work; this chapter exists so you understand what they're doing underneath.
- **`gopacket` and `gopacket/pcap`** wrap libpcap for Go specifically, giving cross-platform capture (Linux, macOS, Windows) plus built-in decoders for a huge range of protocols, at the cost of a cgo dependency on the system's libpcap/Npcap installation. Production Go tooling that needs to run on more than just Linux typically reaches for this instead of hand-rolled `AF_PACKET` code.
- **Kernel-level BPF filters** (the same "Berkeley Packet Filter" bytecode both `tcpdump -i eth0 port 443` and libpcap compile down to) let you discard uninteresting frames *inside the kernel*, before they're copied to userspace at all — critical on a busy interface where copying every single frame to userspace, only to discard 99% of them in your own code, would be wasted CPU and memory bandwidth.
- **eBPF/XDP** (Chapter 105) is the modern evolution of that same idea taken further: your filtering/aggregation logic runs *in the kernel*, sometimes even in the NIC driver before an skb is even allocated, letting tools like Cilium's Hubble or Cloudflare's DDoS mitigation inspect line-rate traffic with a fraction of this chapter's per-packet userspace overhead.
- **This is legally and ethically load-bearing.** Capturing traffic on a network you don't own or don't have explicit authorization to monitor is illegal in most jurisdictions (wiretapping and computer-crime statutes both commonly apply) — everything in this chapter should be run only against your own machine, your own lab network, or infrastructure you have explicit permission to inspect.

---

## 16. What's Simplified Here

- ARP (EtherType `0x0806`) and IPv6 (`0x86DD`) decoding were left as an exercise — the same offset-reading technique applies, just against different header layouts (Chapter 53 for ARP, Chapter 42 for IPv6).
- IPv4 Options (when IHL > 5) and TCP Options (when Data Offset > 5) are correctly *skipped* by this code (via the header-length arithmetic) but never *parsed* — a production dissector would decode SACK, timestamps, and MSS options (Chapter 63) from the TCP options area.
- No checksum verification is performed on captured frames; Chapter 115, Section 7 covers computing (and by extension, verifying) the IPv4 header checksum.
- The simulation mode's synthetic packets don't carry correct TCP/IP checksums, since nothing in this chapter's parsing logic reads them — a byte-perfect capture from a real NIC always has them computed correctly by hardware.
- Fragmented IPv4 packets (Chapter 36's Flags/Fragment Offset fields) are not reassembled here; each fragment would be decoded as if it were a complete, standalone packet, which is only correct for the first fragment.

---

## 17. Interview Questions & Model Answers

**Beginner: Why can't an ordinary TCP or UDP socket see another application's traffic, or the raw headers of its own traffic?**
Because the kernel's normal socket API is built specifically to hide that information: it demultiplexes incoming traffic by port/connection tuple and delivers only the payload of the one connection a socket is bound to, having already stripped the Ethernet, IP, and TCP/UDP headers during processing. Seeing raw frames requires a different, privileged socket type (like `AF_PACKET`) that operates below that demultiplexing step.

**Intermediate: What's the practical difference between hand-rolled `AF_PACKET` code and using `libpcap`/`gopacket`?**
`AF_PACKET` is Linux-specific, requires writing your own bind/receive loop, and gives you every frame with no kernel-side filtering unless you add BPF bytecode yourself. `libpcap` (and `gopacket`, its Go wrapper) is cross-platform — it abstracts over `AF_PACKET` on Linux, BPF devices on macOS/BSD, and Npcap on Windows — and provides a simple filter-expression syntax (`tcp port 443`) that gets compiled to kernel-level BPF automatically, discarding unwanted frames before they're copied to userspace.

**Advanced: A sniffer is running on a laptop connected to a switched office network. Explain precisely what it will and won't be able to capture, and why promiscuous mode doesn't change the answer.**
A switch (Chapter 30) learns which MAC addresses live behind which physical port and forwards each frame only out the port(s) needed to reach its destination — a laptop's port normally only ever receives frames addressed to that laptop (plus broadcast/multicast). Promiscuous mode changes what the laptop's *own NIC* passes up to the OS (normally frames not addressed to it are dropped in hardware) — but it has no effect on which frames the switch chose to deliver to that port in the first place. To see other hosts' unicast traffic, you need either a mirrored/SPAN port configured on the switch itself, a hub (which floods every frame to every port — Chapter 30), or an active attack like ARP spoofing (Chapter 83) that tricks the switch/hosts into sending traffic your way.

---

## 18. Exercises

### Easy
1. Given the raw byte `0x11` in an IPv4 header's Protocol field, what protocol is this packet carrying? Show your work in decimal.
2. Modify `decodeAndPrint` to also print the IPv4 Total Length field, and verify it matches `20 + len(payload)` for the simulated frames in Section 10.
3. What single byte change to `buildTCPSyn`'s flags byte would make it represent a SYN-ACK instead of a bare SYN (Chapter 59)?

### Medium
4. Add ARP decoding: EtherType `0x0806`, and parse the ARP packet's operation field (1 = request, 2 = reply) and sender/target IP and MAC (Chapter 53 has the exact layout).
5. Extend `simulatedCapture` to include a TCP segment with the `RST` flag set, and explain in one sentence what real-world event this represents (Chapter 59).
6. Add a simple BPF-style filter to `main()`: accept a `-port N` flag, and skip printing any TCP/UDP segment whose source and destination ports don't include `N`.

### Hard
7. Implement IPv4 header checksum verification: recompute the checksum from a captured header's bytes (with the checksum field itself zeroed, per Chapter 115 Section 7's algorithm) and flag any packet where it doesn't match.
8. Extend the raw capture loop to reassemble IPv4 fragments (Chapter 36's Flags/Fragment Offset fields) into a complete original packet before decoding the transport-layer header.
9. Port the raw-capture half of this chapter (Section 9) to use `gopacket/pcap` instead of hand-written `AF_PACKET` syscalls, and confirm it produces identical decoded output on macOS or Windows.

---

## 19. Summary

| Term | Meaning |
|---|---|
| Raw socket (`AF_PACKET`/`SOCK_RAW`) | A privileged Linux socket type delivering complete frames, headers included, bypassing normal per-connection demultiplexing |
| `ETH_P_ALL` | The EtherType value meaning "capture every protocol," not just one |
| `htons` | Host-to-network byte order conversion, required when binding an `AF_PACKET` socket |
| Promiscuous mode | A NIC setting that stops dropping frames not addressed to it — does not change what a switch delivers to that port |
| `libpcap` / `gopacket` | The portable, cross-platform capture library `tcpdump` and Wireshark are built on |
| BPF filter | Kernel-level bytecode that discards unwanted frames before they reach userspace |
| Simulation mode | Hand-built, byte-correct synthetic frames used to exercise the decoder without root or live traffic |

Chapter 114 turned three chapters' worth of header diagrams into working, byte-accurate parsing code. Chapter 115 goes one step further than reading headers: it decrements one (TTL), recomputes another (the checksum), and uses a third (the destination address) to actually decide where a packet goes next — building, in Go, the exact longest-prefix-match forwarding algorithm Chapter 45 proved on paper.
