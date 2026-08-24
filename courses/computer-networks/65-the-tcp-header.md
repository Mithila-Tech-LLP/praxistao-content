# Chapter 65: The TCP Header, Field by Field

> **"Every mechanism in the last seven chapters — the handshake, the retransmission logic, the flow control, the congestion control, the graceful close — has to fit into 20 bytes that travel in front of every single segment. This chapter is where you finally see all of it laid out at once, byte by byte."**

---

## Table of Contents

1. [Why This Chapter Comes Last](#1-why-this-chapter-comes-last)
2. [The Full Header Layout](#2-the-full-header-layout)
3. [Source Port and Destination Port](#3-source-port-and-destination-port)
4. [Sequence Number](#4-sequence-number)
5. [Acknowledgment Number](#5-acknowledgment-number)
6. [Data Offset, Reserved Bits, and the Flags Byte](#6-data-offset-reserved-bits-and-the-flags-byte)
7. [The Flags, One by One](#7-the-flags-one-by-one)
8. [Window Size](#8-window-size)
9. [Checksum](#9-checksum)
10. [Urgent Pointer](#10-urgent-pointer)
11. [Options](#11-options)
12. [Every Field Mapped to Its Mechanism](#12-every-field-mapped-to-its-mechanism)
13. [Decoding a Real Captured SYN Packet, Byte by Byte](#13-decoding-a-real-captured-syn-packet-byte-by-byte)
14. [Building a TCP Header Parser in Go](#14-building-a-tcp-header-parser-in-go)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Common Misconceptions](#16-common-misconceptions)
17. [Production Usage Notes](#17-production-usage-notes)
18. [Interview Questions & Model Answers](#18-interview-questions--model-answers)
19. [Exercises](#19-exercises)
20. [Summary and Bridge to Part 10](#20-summary-and-bridge-to-part-10)

---

## 1. Why This Chapter Comes Last

Every chapter in this volume introduced a mechanism first and a field second, deliberately, because a header field with no context is just a number in a diagram — memorizable, but meaningless. You now know *why* TCP needs sequence numbers (Chapter 60), *why* it needs a window field (Chapter 61), *why* it needs SYN and FIN flags (Chapters 59 and 64), and *why* an option field for SACK exists at all (Chapter 63). This chapter does the opposite of every chapter before it: it starts from the wire format and asks you to point backward. If you can look at a raw 20 (or more) bytes of TCP header and say what mechanism each one belongs to and why it's shaped the way it is, you understand TCP as an engineer, not as a diagram-memorizer.

---

## 2. The Full Header Layout

The classic RFC-793-style header diagram, one row per 32 bits (4 bytes), with byte offsets on the left:

```
 Byte  0                   1                   2                   3
       0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
  0   |          Source Port          |       Destination Port       |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
  4   |                        Sequence Number                       |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
  8   |                    Acknowledgment Number                     |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 12   |  Data |Rsv|N|C|E|U|A|P|R|S|F|                                |
      | Offset|   |S|W|C|R|C|S|S|Y|I|         Window Size            |
      |       |   | |R|E|G|K|H|T|N|N|                                |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 16   |            Checksum           |        Urgent Pointer         |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
 20   |               Options (if Data Offset > 5)  ...  |  Padding   |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
      |                             Data ...                          |
      +-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

The minimum TCP header — no options — is exactly **20 bytes**, the same fixed size UDP's 8-byte header is compared against in Chapter 58. Every field before Options has a fixed position and fixed width; that fixed layout is precisely what lets a router or a stack parse a TCP header in a handful of instructions without any variable-length scanning until it reaches Options. Here is the same information as a flat table, since a table is easier to reference while reading a hex dump:

```
Byte offset   Field                              Size
-----------   -----                              ----
 0 - 1        Source Port                        2 bytes  (16 bits)
 2 - 3        Destination Port                   2 bytes  (16 bits)
 4 - 7        Sequence Number                    4 bytes  (32 bits)
 8 - 11       Acknowledgment Number              4 bytes  (32 bits)
12            Data Offset (4 bits, upper nibble) +
              Reserved (3 bits) + NS (1 bit)     1 byte
13            CWR ECE URG ACK PSH RST SYN FIN    1 byte   (the flags)
14 - 15       Window Size                        2 bytes  (16 bits)
16 - 17       Checksum                           2 bytes  (16 bits)
18 - 19       Urgent Pointer                     2 bytes  (16 bits)
20 - (4*DO-1) Options (present only if
              Data Offset > 5)                   0-40 bytes
```

---

## 3. Source Port and Destination Port

**Bytes 0-1 (Source Port) and 2-3 (Destination Port), 16 bits each.**

Chapter 57 established the problem these fields solve: one IP address, many programs, and a way to say which program a segment belongs to. Each port is a plain 16-bit unsigned integer, giving a range of 0-65535. Destination Port is usually the more "meaningful" of the two in casual reading — port 443 for HTTPS, port 22 for SSH — but the *pair* of source and destination port, combined with source and destination IP address (which live in the IP header, one layer down, per Chapter 27's encapsulation model), is what actually identifies one specific connection. This is the 4-tuple Chapter 57 introduced and Chapter 64 relied on directly to explain why a `TIME_WAIT` socket blocks reuse of one specific 4-tuple, not an entire port.

---

## 4. Sequence Number

**Bytes 4-7, 32 bits.**

The field Chapter 60 built its entire reliability story on: not a packet counter, but a **byte counter**. The Sequence Number field holds the sequence number of the *first* byte of data in this particular segment (or, for a SYN, the Initial Sequence Number itself, since a SYN consumes one sequence number even though it carries no application data — Chapter 59). Being 32 bits means sequence numbers wrap around after about 4.3 billion bytes (2^32), a wraparound TCP is explicitly designed to handle correctly via comparison arithmetic that accounts for it (this is part of why an unpredictable Initial Sequence Number, chosen at connection setup, matters for security as well as correctness — Chapter 83 returns to ISN prediction as an attack vector).

---

## 5. Acknowledgment Number

**Bytes 8-11, 32 bits.**

Only meaningful when the ACK flag (Section 7) is set — which, in practice, is true for almost every segment after the initial SYN. Holds the next sequence number the receiver **expects**, which is exactly the cumulative-ACK rule Chapter 60 and Chapter 63 both depend on: "everything up through (Ack Number - 1) has been received." Chapter 63's entire fast-retransmit mechanism is built purely on watching this one field repeat across multiple incoming segments.

---

## 6. Data Offset, Reserved Bits, and the Flags Byte

**Byte 12.**

The upper 4 bits are the **Data Offset** (also called "header length") — the number of 32-bit *words* in the entire TCP header, including options. A value of `5` means 5 × 4 = 20 bytes (no options, the minimum); a value of `10` means 40 bytes (20 bytes of options). This field is exactly why a receiver can tell where the header ends and the data begins even when Options are present — it's the same fixed-width solution to a variable-length-field problem used in several other headers across this course (compare to the IHL field in the IPv4 header, Chapter 36).

The next 3 bits are **Reserved** — set to zero, kept for future use, a pattern this course has already seen elsewhere in protocol design.

The last bit of byte 12 is **NS** (ECN-Nonce Sum), added by RFC 3540 as an experimental extension to ECN (Chapter 62, Section 11) meant to help a sender detect if a network path was deliberately hiding an ECN mark; RFC 3540 was later reclassified historic and NS is not in meaningful production use today, but the bit position is real and still shows up in header diagrams and packet captures.

---

## 7. The Flags, One by One

**Byte 13 — 8 single-bit flags, in this bit order (MSB to LSB): CWR, ECE, URG, ACK, PSH, RST, SYN, FIN.** The 6 classic flags (URG, ACK, PSH, RST, SYN, FIN) come from the original 1981 TCP spec; CWR and ECE were added in 1999 by RFC 3168 to support ECN.

```
Flag   Meaning                                    Mechanism / Chapter
----   -------                                    -------------------
CWR    Congestion Window Reduced — sender is       ECN (Ch. 62 §11):
       telling the peer "I saw your ECE mark        sender echoes back
       and already reduced my cwnd"                 that it reacted

ECE    ECN-Echo — receiver is telling the sender    ECN (Ch. 62 §11):
       "a router marked a packet on this            router's congestion
       connection as experiencing congestion"        mark, relayed back

URG    Urgent — the Urgent Pointer field (Sec. 10)  Legacy out-of-band
       is valid and points to urgent data within     signaling (rarely
       this segment                                  used in practice)

ACK    Acknowledgment — the Acknowledgment          Ch. 60 (reliability);
       Number field is valid; set on almost every    set on every segment
       segment after the initial SYN                 after the handshake

PSH    Push — sender is asking the receiving TCP    Tells the OS to hand
       stack to deliver buffered data to the         data to the app now,
       receiving application immediately, rather      not wait to batch
       than waiting to accumulate more                more into one read

RST    Reset — abort the connection immediately,    Ch. 64 §10: abrupt
       no graceful close, no TIME_WAIT               termination on error
                                                      or deliberate abort

SYN    Synchronize — set on the first two           Ch. 59: three-way
       segments of the handshake, used to            handshake; consumes
       exchange Initial Sequence Numbers             one sequence number

FIN    Finish — sender has no more data to send      Ch. 64: four-way
       in this direction; consumes one sequence      close; each FIN gets
       number just like SYN                          its own ACK
```

A useful pattern to notice: **SYN and FIN are structurally identical in how they're handled** — both consume exactly one byte of sequence-number space, both get individually acknowledged, both can be retransmitted using the same loss-recovery machinery as ordinary data (Chapter 60). This is a deliberate design choice, not a coincidence: it means TCP didn't need a *separate* reliability mechanism for "reliably deliver a connection-state change" — it just reused the one it already had for reliably delivering data.

Common flag combinations you'll see in a real capture:

```
[SYN]              first handshake segment (Ch. 59)
[SYN, ACK]         second handshake segment (Ch. 59)
[ACK]              third handshake segment, and most ordinary data segments
[PSH, ACK]         a data segment the sender wants delivered to the app now
[FIN, ACK]         a close segment (Ch. 64) — FIN combined with acking
                   whatever data the sender had already received
[RST]              abrupt termination (Ch. 64 §10)
[RST, ACK]         abrupt termination that also acknowledges received data
```

---

## 8. Window Size

**Bytes 14-15, 16 bits.**

This is `rwnd` — the **receiver's** flow-control window from Chapter 61, not the sender's congestion window from Chapter 62 (a distinction worth restating precisely here, since it's the single most common point of confusion this entire volume addresses). A raw 16-bit field can only express window sizes up to 65,535 bytes directly — far too small for the high-bandwidth, high-latency links Chapter 62 discussed as CUBIC's motivating case. The **Window Scale option** (Section 11) fixes this by having both sides agree, during the handshake, on a scaling factor to multiply this field by, letting the *effective* window reach into the gigabytes while the field on the wire stays 16 bits.

---

## 9. Checksum

**Bytes 16-17, 16 bits.**

TCP's error-detection field, computed the same conceptual way as UDP's checksum (Chapter 58) and connecting back to the error-detection theory of Chapter 19: a 16-bit one's-complement sum computed over the TCP header, the TCP payload, and a **pseudo-header** built from fields borrowed from the IP layer (source IP, destination IP, protocol number, and TCP segment length) — specifically so that a segment accidentally delivered to, or addressed from, the wrong IP address gets caught even though the IP addresses themselves aren't literally part of the TCP header. Like UDP's checksum, this is a weak, non-cryptographic check — it catches accidental bit-flips reliably but is trivially forgeable, which is exactly the gap TLS (Chapter 82) exists to close at a much stronger, cryptographic level.

---

## 10. Urgent Pointer

**Bytes 18-19, 16 bits.**

Only meaningful when the URG flag (Section 7) is set. Points to the sequence-number offset (relative to this segment's Sequence Number) marking the end of "urgent data" — historically intended to let an application signal something like a Telnet interrupt (Ctrl-C) that should jump ahead of ordinary buffered data. In practice, this mechanism is rarely used by modern application protocols; most designs that need an equivalent of "urgent, out-of-band signaling" build it at the application layer instead (a special message type inside the ordinary data stream) rather than relying on URG, partly because different TCP/IP stacks historically implemented the exact semantics of the urgent pointer slightly inconsistently. It remains a real, standardized field, just a mostly dormant one in modern production traffic.

---

## 11. Options

**Bytes 20 onward, present only when Data Offset (Section 6) is greater than 5.** Options are encoded as a sequence of `(Kind, Length, Value)` triples (with a couple of single-byte exceptions), padded with `NOP` (Kind 1) bytes as needed so the whole header ends on a 32-bit boundary — which is exactly why Data Offset is measured in 32-bit words rather than raw bytes.

```
Kind  Name                Length   Introduced by / used for
----  ----                ------   -------------------------
  0   End of Option List   1       Marks the end of the options list
  1   No-Operation (NOP)   1       Padding / alignment filler
  2   Maximum Segment      4       Negotiated at handshake time (Ch. 59):
      Size (MSS)                    "don't send me a segment bigger than
                                     this" — sized to avoid IP fragmentation
  3   Window Scale         3       Ch. 61: multiplies the 16-bit Window
                                     field by 2^shift, negotiated once at
                                     handshake, only on SYN segments
  4   SACK-Permitted       2       Ch. 63: negotiates whether both sides
                                     support Selective Acknowledgment
  5   SACK                variable  Ch. 63: the actual selective-ack blocks,
                                     sent on ordinary segments once
                                     SACK-Permitted was negotiated
  8   Timestamps           10      Ch. 60/62: TSval/TSecr pair, used for
                                     more accurate RTT measurement (feeding
                                     RTO calculation) and safe TIME_WAIT
                                     reuse (Ch. 64's tcp_tw_reuse)
```

Every one of these options connects directly to a mechanism already built in an earlier chapter — Options exist specifically because the fixed 20-byte header couldn't grow to accommodate new mechanisms (window scaling, SACK, better RTT measurement) invented well after 1981 without breaking every existing implementation that only understood the fixed fields; Options let TCP evolve without changing its core wire format.

---

## 12. Every Field Mapped to Its Mechanism

A single reference table tying every field in this chapter back to the chapter that explains why it exists:

```
Field                     Mechanism                          Chapter
-----                     ---------                          -------
Source/Destination Port   Multiplexing many programs          57
Sequence Number           Byte-level reliable ordering        60
Acknowledgment Number     Cumulative acknowledgment            60
SYN flag                  Three-way handshake, ISN exchange   59
FIN flag                  Four-way graceful close              64
RST flag                  Abrupt termination                   64
ACK flag                  Marks Ack Number as valid             60
PSH flag                  Immediate delivery to application    (app-layer)
URG flag / Urgent Pointer Legacy out-of-band signaling          (legacy)
Window Size               Receiver flow control (rwnd)          61
Window Scale option       Scaling rwnd past 65,535 bytes        61
MSS option                Segment sizing at handshake time      59
CWR / ECE flags           Explicit Congestion Notification      62
SACK-Permitted / SACK     Precise multi-loss recovery            63
Timestamps option         RTT measurement, safe TIME_WAIT reuse 60 / 64
Checksum                  Error detection                       19 / 58
Data Offset               Locating the start of Options/data    (this ch.)
```

If you can reproduce this table from memory, you have genuinely learned Volume 9, not just read it.

---

## 13. Decoding a Real Captured SYN Packet, Byte by Byte

Below is a representative TCP SYN segment, formatted exactly the way `tcpdump -xx` or a hex-view pane in Wireshark would display it, for a client (port 51234) opening a connection to an HTTPS server (port 443) with a typical modern Linux option set: MSS, SACK-Permitted, Timestamps, and Window Scale.

```
0000   c8 22 01 bb 6f 3a 9d 17 00 00 00 00 a0 02 fa f0
0010   8f 3c 00 00 02 04 05 b4 04 02 08 0a 1a 2b 3c 4d
0020   00 00 00 00 01 03 03 07
```

Decoded field by field, using the byte offsets from Section 2:

```
Offset  Bytes              Field                    Value
------  -----              -----                    -----
 0- 1   c8 22              Source Port              0xC822 = 51234
 2- 3   01 bb              Destination Port         0x01BB = 443 (HTTPS)
 4- 7   6f 3a 9d 17        Sequence Number          0x6F3A9D17 = 1,865,791,255
                                                     (the Initial Sequence
                                                     Number, chosen pseudo-
                                                     randomly per Ch. 59)
 8-11   00 00 00 00        Acknowledgment Number    0 (irrelevant - ACK
                                                     flag not yet set)
12      a0                 Data Offset / Rsv / NS   0xA0 = 1010 0000
                                                     Data Offset = 1010 = 10
                                                     -> header is 10*4 = 40
                                                     bytes (20 fixed + 20
                                                     of options, matches
                                                     the 0x0028=40 total
                                                     bytes shown below)
13      02                 Flags                    0x02 = 0000 0010
                                                     -> only SYN is set
14-15   fa f0              Window Size              0xFAF0 = 64,240 bytes
                                                     (typical Linux default
                                                     initial advertised
                                                     window before scaling)
16-17   8f 3c              Checksum                 0x8F3C (computed over
                                                     header+pseudo-header;
                                                     verify, don't hand-
                                                     compute, in practice)
18-19   00 00              Urgent Pointer           0 (URG not set, unused)
20-23   02 04 05 b4        Option: MSS              Kind=2, Len=4,
                                                     Value=0x05B4=1460 bytes
24-25   04 02              Option: SACK-Permitted   Kind=4, Len=2
26-35   08 0a 1a 2b 3c 4d  Option: Timestamps       Kind=8, Len=10,
        00 00 00 00        (10 bytes)               TSval=0x1A2B3C4D,
                                                     TSecr=0 (no prior
                                                     timestamp to echo yet)
36      01                 Option: NOP              Kind=1 (pure padding,
                                                     aligns next option)
37-39   03 03 07           Option: Window Scale     Kind=3, Len=3,
                                                     Shift count = 7
                                                     (effective window =
                                                     Window Size << 7)
```

Reading this the way a real network engineer would, top to bottom: a client at some ephemeral port is opening a connection to port 443 (so, almost certainly, an HTTPS/TLS session — the TLS handshake itself, Chapter 82, will start only after this TCP handshake completes); the flags byte confirms this is the very first segment of a brand-new connection (SYN only, no ACK yet); the options reveal a modern, well-configured stack — MSS capped at 1460 bytes (the standard value for a 1500-byte Ethernet MTU minus 20 bytes of IPv4 header and 20 bytes of TCP header, tying directly back to Chapter 28's frame-size discussion), SACK support offered, a timestamp for RTT measurement, and window scaling by a shift of 7 (multiplying the eventual advertised window by 128, letting this connection's effective `rwnd` reach well past the raw 16-bit field's 65,535-byte ceiling — directly enabling the large in-flight windows Chapter 61 and Chapter 62's long-fat-network discussion assume are possible).

*(A note on honesty, in keeping with this course's standard: the exact byte values above are constructed to be internally consistent and field-accurate — they decode exactly as shown, and every value is realistic for a real modern client — rather than lifted from one specific packet capture file. A genuine capture from your own machine, taken with `tcpdump -xx` as shown in Section 12's earlier chapters' hands-on exercises, will show the identical structure with different specific numbers.)*

---

## 14. Building a TCP Header Parser in Go

A minimal, from-scratch parser for the fixed 20-byte portion of a TCP header plus its options, operating directly on a byte slice the way `tcpdump`/`gopacket`-style tools do internally:

```go
package main

import (
	"encoding/binary"
	"fmt"
)

type TCPHeader struct {
	SrcPort    uint16
	DstPort    uint16
	SeqNum     uint32
	AckNum     uint32
	DataOffset uint8 // in 32-bit words
	Flags      uint8
	Window     uint16
	Checksum   uint16
	UrgentPtr  uint16
}

const (
	FlagFIN = 1 << 0
	FlagSYN = 1 << 1
	FlagRST = 1 << 2
	FlagPSH = 1 << 3
	FlagACK = 1 << 4
	FlagURG = 1 << 5
	FlagECE = 1 << 6
	FlagCWR = 1 << 7
)

func ParseTCPHeader(b []byte) (*TCPHeader, error) {
	if len(b) < 20 {
		return nil, fmt.Errorf("packet too short for a TCP header: %d bytes", len(b))
	}
	h := &TCPHeader{
		SrcPort:   binary.BigEndian.Uint16(b[0:2]),
		DstPort:   binary.BigEndian.Uint16(b[2:4]),
		SeqNum:    binary.BigEndian.Uint32(b[4:8]),
		AckNum:    binary.BigEndian.Uint32(b[8:12]),
		Window:    binary.BigEndian.Uint16(b[14:16]),
		Checksum:  binary.BigEndian.Uint16(b[16:18]),
		UrgentPtr: binary.BigEndian.Uint16(b[18:20]),
	}
	h.DataOffset = b[12] >> 4  // top 4 bits of byte 12
	h.Flags = b[13]            // the full flags byte
	return h, nil
}

func (h *TCPHeader) HeaderLen() int { return int(h.DataOffset) * 4 }

func (h *TCPHeader) FlagString() string {
	names := []struct {
		bit  uint8
		name string
	}{
		{FlagCWR, "CWR"}, {FlagECE, "ECE"}, {FlagURG, "URG"}, {FlagACK, "ACK"},
		{FlagPSH, "PSH"}, {FlagRST, "RST"}, {FlagSYN, "SYN"}, {FlagFIN, "FIN"},
	}
	out := ""
	for _, f := range names {
		if h.Flags&f.bit != 0 {
			if out != "" {
				out += ","
			}
			out += f.name
		}
	}
	if out == "" {
		return "(none)"
	}
	return out
}

func main() {
	// The exact SYN segment decoded by hand in Section 13
	raw := []byte{
		0xc8, 0x22, 0x01, 0xbb, 0x6f, 0x3a, 0x9d, 0x17, 0x00, 0x00, 0x00, 0x00,
		0xa0, 0x02, 0xfa, 0xf0, 0x8f, 0x3c, 0x00, 0x00,
		0x02, 0x04, 0x05, 0xb4, 0x04, 0x02, 0x08, 0x0a,
		0x1a, 0x2b, 0x3c, 0x4d, 0x00, 0x00, 0x00, 0x00,
		0x01, 0x03, 0x03, 0x07,
	}

	h, err := ParseTCPHeader(raw)
	if err != nil {
		panic(err)
	}

	fmt.Printf("Src Port:     %d\n", h.SrcPort)
	fmt.Printf("Dst Port:     %d\n", h.DstPort)
	fmt.Printf("Seq Number:   %d\n", h.SeqNum)
	fmt.Printf("Ack Number:   %d\n", h.AckNum)
	fmt.Printf("Header Len:   %d bytes\n", h.HeaderLen())
	fmt.Printf("Flags:        %s\n", h.FlagString())
	fmt.Printf("Window Size:  %d\n", h.Window)
	fmt.Printf("Checksum:     0x%04x\n", h.Checksum)
}
```

Running this against the exact bytes decoded by hand in Section 13 should print `Src Port: 51234`, `Dst Port: 443`, `Header Len: 40 bytes`, and `Flags: SYN` — a good sanity check that the manual decode and the programmatic decode agree, and a genuinely useful skeleton to extend (parsing Options is a natural next step, left as Exercise 9).

---

## 15. What's Simplified Here

In the interest of an honest account of complexity, a few things this chapter simplified:

- Real Options parsing has to handle Kind values this chapter didn't list (e.g., experimental options, Kind values used by newer RFCs like TCP Fast Open's Kind 34), and must correctly handle malformed or truncated option lists defensively — production-grade parsers (and the Linux kernel itself) are far more defensive than the teaching parser in Section 14.
- The Checksum field's exact computation (the pseudo-header construction, one's-complement arithmetic, handling of the payload's odd/even byte length) was described conceptually, not implemented — Chapter 19's error-detection chapter is the right place to see the actual bit-level algorithm.
- This chapter treated IPv4's pseudo-header contribution to the checksum as a given; IPv6's pseudo-header (used when TCP runs over IPv6, per Chapter 42) has a different, larger structure to accommodate 128-bit addresses, though the checksum's role is identical.
- Modern high-performance network stacks (and NICs with checksum/segmentation offload) frequently compute or verify the checksum in hardware, and the kernel may hand a NIC an unsegmented "superpacket" for it to slice into MSS-sized real segments on transmit (TCP Segmentation Offload) — the clean, tidy one-segment-per-packet picture this chapter draws is a correct logical model but not always literally how the bytes move through a modern high-speed NIC.

---

## 16. Common Misconceptions

- **"The TCP header is always 20 bytes."** Only when Data Offset is 5 (no options). Modern connections using MSS, SACK-Permitted, Timestamps, and Window Scale together — as in Section 13's example — routinely have a 32- or 40-byte header.
- **"Sequence numbers count packets."** They count bytes (Section 4) — a point worth restating one more time here because it's easy to slip back into "packet number" thinking when staring at a hex dump instead of a data stream.
- **"The Window field is the congestion window."** Restated one final time because Chapter 62 built an entire chapter around this exact confusion: the Window field on the wire is always `rwnd`, the receiver's flow-control advertisement — `cwnd` is never transmitted, ever, in any TCP header.
- **"PSH means 'this is urgent.'** PSH and URG are unrelated flags with unrelated purposes — PSH is about delivery timing to the local application, URG is about a specific offset of "urgent" data inside the stream. Confusing the two is a common but incorrect reading of a packet capture.
- **"A SYN packet can't carry options because it's the first packet."** The opposite is true — the SYN and SYN-ACK segments are specifically where the most important options (MSS, Window Scale, SACK-Permitted) are negotiated, because those options only make sense if agreed once, at the start of the connection, by both sides.

---

## 17. Production Usage Notes

- `tcpdump -xx` or `tcpdump -X` on a real interface will show you a live version of exactly the hex dump decoded in Section 13 — this is a genuinely worthwhile five-minute exercise: capture one real outbound SYN from your own machine and decode it by hand against Section 2's table before trusting Wireshark's automatic decode.
- Wireshark's packet-detail pane effectively performs the Section 12 mapping automatically and interactively — expanding the "Transmission Control Protocol" section of any captured segment shows every field from this chapter labeled and, for flags, individually clickable.
- Middlebox interference with TCP Options is a real, still-relevant operational hazard: some older or misconfigured firewalls and NAT devices strip unrecognized options (or, historically, the Window Scale option specifically), silently disabling a feature both real endpoints support — a known, documented cause of unexpectedly poor throughput that is diagnosed exactly by comparing what options were sent versus what the far end's replies indicate were understood.
- Load balancers and proxies that terminate TCP (Chapter 95) construct entirely new TCP headers for the connection to the backend — meaning fields like the Initial Sequence Number, Window Scale shift, and Timestamps are independently negotiated on each leg, a detail that matters when tracing a Sequence Number mismatch across a proxy hop during debugging (Chapter 122's debugging playbook returns to this).

---

## 18. Interview Questions & Model Answers

**Q (Beginner): What is the minimum size of a TCP header, and what makes it larger?**

*Model answer:* "The minimum TCP header is 20 bytes — that's everything up through the Urgent Pointer field, with no Options. It gets larger when Options are present, such as Maximum Segment Size, Window Scale, SACK-Permitted, or Timestamps, which are common on the SYN and SYN-ACK segments of a modern connection. The Data Offset field, in the upper 4 bits of byte 12, records the total header length in 32-bit words specifically so a receiver knows exactly where the header ends and the payload begins, even with variable-length options present."

**Q (Intermediate): Explain the difference between the SYN flag and the ACK flag, and describe a segment where both would be set.**

*Model answer:* "SYN is set on segments that are establishing a connection and exchanging an Initial Sequence Number — it's set on the first segment of the handshake and, alongside ACK, on the second. ACK simply indicates the Acknowledgment Number field is valid and being used to acknowledge previously received data; it's set on virtually every segment after the initial SYN. A segment with both SYN and ACK set is exactly the second message of the three-way handshake — the server is simultaneously proposing its own Initial Sequence Number (SYN) and acknowledging the client's (ACK)."

**Q (Advanced): Walk through how the Window field and the Window Scale option interact, and why this design (rather than simply making the Window field bigger) was chosen.**

*Model answer:* "The Window field on the wire is a fixed 16-bit value, capping the raw advertised window at 65,535 bytes — too small for high-bandwidth, high-latency paths where the bandwidth-delay product can be many megabytes. Rather than redefining the fixed header (which would break every existing TCP implementation that only understands a 16-bit Window field), the Window Scale option is negotiated once, only on the SYN and SYN-ACK segments, and specifies a shift count both sides agree to apply: the actual, effective window is the value in the Window field left-shifted by that count. This keeps the wire format backward compatible — an old stack that doesn't understand Window Scale just ignores the unknown option and behaves as if scaling were 0 (shift by zero, i.e., no change) — while letting modern, cooperating stacks negotiate effective windows into the gigabytes. It's the same evolutionary pattern SACK and Timestamps follow: extend capability through Options rather than changing the fixed fields everyone already depends on."

---

## 19. Exercises

### Easy

1. Using Section 2's table, state the byte offset and size of the Acknowledgment Number field.
2. A captured segment has flags byte `0x11`. Using Section 7's bit assignments, which two flags are set, and what real-world event does that combination most likely represent?
3. Why is Data Offset measured in 32-bit words instead of raw bytes?

### Medium

4. A TCP header has Data Offset = 8. How many total bytes is the header, and how many bytes of Options does it contain?
5. Explain, referencing both this chapter and Chapter 61, exactly what it means for the Window field to hold `0xFFFF` on a connection that negotiated a Window Scale shift of 7 — what is the effective window in bytes?
6. Take the flags byte `0x12` and the flags byte `0x18`; decode both using Section 7's table, and explain what kind of segment each most plausibly represents in a real connection's lifecycle.

### Hard

7. Given the raw hex bytes `45 00 00 00 01 bb c8 22 00 00 00 00 6f 3a 9d 18 50 10 fa f0 00 00 00 00` for a *response* segment's TCP header (assume no options, Data Offset=5), decode every field using Section 2's layout and explain what event in a connection's lifecycle this segment most likely represents. (Hint: check the flags byte and compare the Ack Number to Section 13's original SYN's Sequence Number.)
8. Extend the Go parser in Section 14 to also parse the Options list (Kind/Length/Value triples) for at least the MSS, SACK-Permitted, Timestamps, and Window Scale options, and print each one's decoded value.
9. Explain why a middlebox silently stripping the Window Scale option from a SYN segment (but passing everything else through) can cause a measurable, real throughput problem on a long-distance, high-bandwidth connection, even though the connection still completes its handshake successfully and works. Reference Chapter 61's bandwidth-delay-product discussion directly.

---

## 20. Summary and Bridge to Part 10

| Field | Size | Mechanism it belongs to |
|---|---|---|
| Source Port | 16 bits | Multiplexing (Ch. 57) |
| Destination Port | 16 bits | Multiplexing (Ch. 57) |
| Sequence Number | 32 bits | Byte-ordered reliability (Ch. 60) |
| Acknowledgment Number | 32 bits | Cumulative ACK (Ch. 60) |
| Data Offset | 4 bits | Locating end of header/options |
| Reserved + NS | 4 bits | Reserved; NS is a dormant ECN extension |
| CWR, ECE | 2 bits | Explicit Congestion Notification (Ch. 62) |
| URG, ACK, PSH, RST, SYN, FIN | 6 bits | Urgent data, ack validity, push, reset, handshake, close (Ch. 59, 60, 64) |
| Window Size | 16 bits | Receiver flow control, `rwnd` (Ch. 61) |
| Checksum | 16 bits | Error detection (Ch. 19, 58) |
| Urgent Pointer | 16 bits | Legacy out-of-band signaling |
| Options (MSS, WScale, SACK, Timestamps, etc.) | variable | Handshake sizing, large windows, precise loss recovery, RTT measurement (Ch. 59, 61, 63) |

That table is the entire Transport Layer volume compressed into one page: every mechanism you learned from Chapter 57 through Chapter 64 exists because some field in this header needed to carry the information that mechanism depends on, and every field in this header exists because some real, concrete problem — multiplexing, reliability, flow control, congestion, fast recovery, or graceful shutdown — needed solving.

You now have the complete Transport Layer as it actually exists on the wire. But a header full of correctly-numbered bytes is useless if the two hosts don't know *where to send it in the first place*. Every example in this chapter quietly assumed you already had an IP address to put in the IP header underneath this TCP header — and no human ever actually knows one of those from memory. Chapter 66 opens Part 10 by asking the question that makes DNS necessary at all: nobody should have to remember `93.184.216.34` to reach a server, so how does a name like `example.com` become that address in the first place?
