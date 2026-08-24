# Chapter 119: Wireshark and tcpdump — Seeing the Network

> **"Every protocol in this course — the handshake, the DNS query, the TLS exchange — has so far been a diagram on a page. This chapter hands you the tools that turn each one into bytes you can point at, on your own machine, right now."**

---

## Table of Contents

1. [Why This Chapter Exists](#1-why-this-chapter-exists)
2. [tcpdump: Fast, Text, Everywhere](#2-tcpdump-fast-text-everywhere)
3. [tcpdump Filter Syntax, From First Principles](#3-tcpdump-filter-syntax-from-first-principles)
4. [Wireshark: The Visual Instrument](#4-wireshark-the-visual-instrument)
5. [Capture Filters vs. Display Filters](#5-capture-filters-vs-display-filters)
6. [Wireshark Display Filter Syntax](#6-wireshark-display-filter-syntax)
7. [Reading a Captured TCP Handshake, Byte by Byte](#7-reading-a-captured-tcp-handshake-byte-by-byte)
8. [Reading a Captured DNS Query](#8-reading-a-captured-dns-query)
9. [Reading a Captured TLS 1.3 Handshake](#9-reading-a-captured-tls-13-handshake)
10. [Saving and Replaying Captures: pcap Files](#10-saving-and-replaying-captures-pcap-files)
11. [A Real Diagnostic Capture Walkthrough](#11-a-real-diagnostic-capture-walkthrough)
12. [Hands-On Experiment](#12-hands-on-experiment)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Production Notes](#14-production-notes)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary and the Bridge to Chapter 120](#18-summary-and-the-bridge-to-chapter-120)

---

## 1. Why This Chapter Exists

Chapter 56 ended by handing you `tcpdump` for exactly one line of output — enough to recognize `Flags [S]`, `[S.]`, and `[.]` as the three-way handshake in flight, with a promise that a full chapter was coming. Chapter 59 built the handshake from first principles, ISN by ISN. Chapter 65 decoded a single captured TCP header, byte offset by byte offset. Chapter 66–69 walked through the DNS hierarchy and record types. Chapter 82 walked the TLS 1.3 handshake message by message. Every one of those chapters described what *should* be on the wire. This chapter is about actually looking at the wire and confirming it.

The problem this chapter solves is a specific, practical one: **an application-level symptom ("the page won't load," "the connection resets") gives you almost no information about which layer is actually broken.** `curl -v` (Chapter 56, Section 10) gives you a first read, but it only shows *its own* view of one connection, summarized and after the fact. To see what genuinely crossed the wire — including packets involving *other* processes, other machines' conversations you're allowed to observe, retransmissions, resets, and malformed traffic that no well-behaved client would ever construct — you need a tool that captures raw frames off a network interface before any application gets to interpret them.

Two tools dominate this job, and they are not competitors so much as two views of the same underlying capability:

- **`tcpdump`** — a command-line packet capture and filter tool, present by default on nearly every Linux and macOS system, built for speed, scriptability, and use over SSH on a remote server with no graphical interface at all.
- **Wireshark** — a graphical packet analyzer built on the same underlying capture library (`libpcap`/`WinPcap`/`Npcap`) as `tcpdump`, adding a visual protocol tree, color-coded packet lists, "Follow TCP Stream" reconstruction, and a much richer filter and statistics engine.

Both read the exact same bytes off the exact same interface. The choice between them is almost never "which one is more correct" — it's "am I on a headless server over SSH" (reach for `tcpdump`) or "do I have a GUI and want to visually dig through hundreds of packets" (reach for Wireshark, often by capturing with `tcpdump -w` on the remote box and opening the file locally in Wireshark — Section 10 shows exactly how).

## 2. tcpdump: Fast, Text, Everywhere

`tcpdump` is a thin, fast wrapper around `libpcap`, the packet capture library that talks to the operating system's raw packet interface (on Linux, `AF_PACKET` sockets; conceptually the same idea Chapter 114 has you build a simplified version of in Go). Every `tcpdump` invocation has the same basic shape:

```
tcpdump [options] [filter expression]
```

The two options you'll type constantly:

| Flag | Meaning |
|---|---|
| `-i <iface>` | which interface to capture on (`eth0`, `en0`, `any` for all interfaces) |
| `-n` | don't resolve IP addresses to hostnames or ports to service names — faster, and shows the literal numbers |
| `-v`, `-vv`, `-vvv` | increasing verbosity of protocol decode |
| `-X` | show packet payload as hex **and** ASCII |
| `-c N` | capture exactly N packets, then stop |
| `-w file.pcap` | write raw captured packets to a file instead of printing them |
| `-r file.pcap` | read and print packets from a previously saved file |
| `-A` | show payload as ASCII only (handy for spotting plaintext HTTP) |

A minimal capture, requiring root or `CAP_NET_RAW` because reading raw frames off an interface is a privileged operation by design (any unprivileged process being able to see everyone else's traffic would be a serious security hole):

```
$ sudo tcpdump -i eth0 -n
tcpdump: verbose output suppressed, use -v[v]... for full protocol decode
listening on eth0, link-type EN10MB (Ethernet), snapshot length 262144 bytes
10:41:02.001122 IP 192.168.1.132.51234 > 93.184.216.34.443: Flags [S], seq 1865791255, win 64240, length 0
10:41:02.014988 IP 93.184.216.34.443 > 192.168.1.132.51234: Flags [S.], seq 2984651007, ack 1865791256, win 65160, length 0
10:41:02.015103 IP 192.168.1.132.51234 > 93.184.216.34.443: Flags [.], ack 2984651008, win 502, length 0
```

With no filter at all, this is a firehose — every packet touching the interface, from every process, every protocol, every remote host — which is exactly why Section 3 exists.

## 3. tcpdump Filter Syntax, From First Principles

`tcpdump`'s filter language is **BPF (Berkeley Packet Filter)** expression syntax — the same underlying filter mechanism eBPF (Chapter 105) generalized decades later into a full in-kernel virtual machine. A BPF filter is built from **primitives** combined with boolean operators (`and`, `or`, `not`, or their symbolic forms `&&`, `||`, `!`).

**Primitives by what they match:**

| Primitive | Matches | Example |
|---|---|---|
| `host <addr>` | either source or destination is this IP | `host 93.184.216.34` |
| `src host <addr>` / `dst host <addr>` | one direction only | `src host 192.168.1.132` |
| `net <cidr>` | source or destination is in this network | `net 192.168.1.0/24` |
| `port <n>` | either source or destination port is n | `port 443` |
| `src port <n>` / `dst port <n>` | one direction only | `dst port 53` |
| `tcp` / `udp` / `icmp` / `arp` | protocol type | `tcp` |
| `ether host <mac>` | Ethernet frame's MAC address (Chapter 29) | `ether host aa:bb:cc:dd:ee:10` |

**Combining primitives** with boolean logic is the same as any programming condition:

```
$ sudo tcpdump -i eth0 -n host 93.184.216.34 and port 443
$ sudo tcpdump -i eth0 -n 'src 192.168.1.132 and (dst port 80 or dst port 443)'
$ sudo tcpdump -i eth0 -n 'not port 22'
```

The parentheses need quoting in the shell because `(` and `|` are shell metacharacters — this is a shell-quoting detail, not a BPF one, but it trips up almost everyone the first time.

**The powerful, less obvious feature: reaching directly into packet bytes.** BPF lets you index into a packet's raw bytes and test bits directly, which is exactly how you filter on TCP flags — a field Chapter 65 placed at byte offset 13 of the TCP header:

```
tcp[13]
```

means "byte 13 of the TCP header" — the flags byte. `tcp[tcpflags]` is a friendlier named alias for the same byte, and `tcp-syn`, `tcp-ack`, `tcp-fin`, `tcp-rst`, `tcp-push`, `tcp-urg` are named bit-mask constants for each flag bit from Chapter 65, Section 7. So the filter the task's own motivating example uses:

```
$ sudo tcpdump -i eth0 'tcp[tcpflags] & tcp-syn != 0'
```

reads literally as: "take the flags byte, bitwise-AND it with the SYN bit's mask, and keep the packet if the result is non-zero" — i.e., **any packet with the SYN flag set**, which includes both the initial SYN and the SYN-ACK of every handshake on the interface (Chapter 59). To see *only* the very first SYN of new connection attempts (excluding the SYN-ACK, which also has ACK set), add a second condition that excludes ACK:

```
$ sudo tcpdump -i eth0 'tcp[tcpflags] & (tcp-syn|tcp-ack) == tcp-syn'
```

This is genuinely just Boolean algebra over the same flags byte Chapter 65 decoded by hand — `tcpdump` gives you a way to write that same bit test as a filter instead of reading it visually off a hex dump. A few more concrete, commonly-used filters:

```
# Only RST packets — useful for spotting refused connections or firewall drops
$ sudo tcpdump -i eth0 'tcp[tcpflags] & tcp-rst != 0'

# Only packets carrying actual TCP payload (not pure ACKs), by checking IP total length
# minus header lengths is nonzero — commonly approximated instead with:
$ sudo tcpdump -i eth0 'tcp[tcpflags] & tcp-push != 0'

# Only DNS traffic (UDP or TCP port 53 — Chapter 66)
$ sudo tcpdump -i eth0 'port 53'

# Only ARP (Chapter 53) — no host/port makes sense here, ARP isn't IP at all
$ sudo tcpdump -i eth0 arp
```

## 4. Wireshark: The Visual Instrument

Wireshark opens the same raw capture (live from an interface, or a `.pcap`/`.pcapng` file saved by `tcpdump` or by Wireshark itself) into three coordinated panes:

```
+----------------------------------------------------------------+
|  Packet List  (one row per packet: No., Time, Source,          |
|                Destination, Protocol, Length, Info)             |
+----------------------------------------------------------------+
|  Packet Details  (the protocol tree: Frame > Ethernet >         |
|                    IP > TCP > ... — expandable, field by field) |
+----------------------------------------------------------------+
|  Packet Bytes  (raw hex + ASCII, with the field selected above  |
|                 in Packet Details highlighted here too)         |
+----------------------------------------------------------------+
```

Click any packet in the top pane, and the middle pane shows the full encapsulation stack from Chapter 27 as a literal, expandable tree — click "Transmission Control Protocol" and every field Chapter 65 named (Source Port, Sequence Number, Flags, Window Size, ...) appears with its actual decoded value, and clicking any individual field highlights the exact bytes that encode it in the bottom hex pane. This is, functionally, an interactive version of Chapter 65 Section 13's manual byte decode — Wireshark is doing the same offset arithmetic you did by hand, for every packet, instantly.

A feature with no `tcpdump` equivalent worth calling out immediately: **right-click any TCP packet → Follow → TCP Stream**. This reassembles every segment of one connection, strips the TCP framing, and shows you the reconstructed application-layer conversation (the literal HTTP request and response text, if unencrypted) as a continuous, readable transcript — genuinely one of the fastest ways to see "what did this connection actually say" without manually piecing segments back together.

## 5. Capture Filters vs. Display Filters

This is the single most common point of confusion for anyone moving from `tcpdump` to Wireshark, so it's worth stating precisely: **Wireshark has two entirely different filter languages, applied at two different times.**

- **Capture filters** are applied *before* a packet is even written to memory or disk — set in Wireshark's "Capture Options" dialog. They use **the exact same BPF syntax as `tcpdump`** from Section 3 (`host 93.184.216.34 and port 443`), because they're implemented by the same underlying `libpcap` engine.
- **Display filters** are applied *after* capture, to decide what's shown in the Packet List pane out of everything already captured. They use **a completely different, Wireshark-specific syntax** based on protocol field names (`tcp.flags.syn == 1`), described in full in Section 6.

The practical guidance: capture broadly (or with a light BPF capture filter, if the interface is extremely busy and you already know roughly what you're after) and refine with a display filter afterward — display filters can be typed and retyped instantly without re-capturing, while a capture filter mistake means recapturing from scratch. Most working Wireshark sessions use a wide-open (or no) capture filter and do essentially all of their narrowing with display filters.

## 6. Wireshark Display Filter Syntax

Display filter expressions are built from **dotted field names**, mirroring the protocol tree structure from Section 4 — `tcp.flags.syn` is literally "the `flags.syn` field inside the `tcp` protocol layer." Comparison operators are the familiar ones (`==`, `!=`, `<`, `>`, `<=`, `>=`), plus `and`/`or`/`not`.

| Display Filter | Matches | Chapter It Reads |
|---|---|---|
| `ip.addr == 93.184.216.34` | source or destination IP | Ch 36 |
| `tcp.port == 443` | source or destination TCP port | Ch 57 |
| `tcp.flags.syn == 1` | SYN flag set (both SYN and SYN-ACK) | Ch 59, 65 |
| `tcp.flags.syn == 1 and tcp.flags.ack == 0` | SYN only, not SYN-ACK | Ch 59 |
| `tcp.flags.reset == 1` | RST — refused or aborted connections | Ch 64, 83 |
| `tcp.analysis.retransmission` | Wireshark's own inferred retransmissions | Ch 60 |
| `dns` | any DNS packet | Ch 66-69 |
| `dns.qry.name == "example.com"` | DNS queries for this exact name | Ch 66-69 |
| `tls.handshake.type == 1` | TLS ClientHello specifically | Ch 82 |
| `http.request` | HTTP request messages only | Ch 71 |
| `frame contains "GET"` | raw byte search anywhere in the frame | (any protocol) |

Two things make display filters meaningfully more powerful than BPF: **field-level access to fully-parsed protocol data** (`dns.qry.name`, `tls.handshake.type` — BPF has no concept of "the DNS question name," only raw byte offsets), and **derived analysis fields that Wireshark itself computes**, like `tcp.analysis.retransmission`, `tcp.analysis.duplicate_ack`, and `tcp.analysis.zero_window` — Wireshark tracks every TCP stream's state across the whole capture and flags segments that don't fit the expected sequence, directly surfacing the retransmission and flow-control mechanics Chapters 60–61 described in the abstract. There is no `tcpdump` BPF equivalent to "flag every retransmission" — that requires stateful tracking across the whole stream, which BPF's stateless, per-packet filter model cannot do.

## 7. Reading a Captured TCP Handshake, Byte by Byte

Chapter 59 built the three-way handshake conceptually; Chapter 65, Section 13 decoded one raw SYN's hex bytes. Here is the same handshake, captured end to end, read the way you'd actually read it in Wireshark's packet list, then confirmed field by field:

```
No.  Time      Source              Destination         Protocol  Info
1    0.000000  192.168.1.132       93.184.216.34       TCP       51234 → 443 [SYN] Seq=0 Win=64240 Len=0 MSS=1460 SACK_PERM TSval=... WS=128
2    0.013866  93.184.216.34       192.168.1.132       TCP       443 → 51234 [SYN, ACK] Seq=0 Ack=1 Win=65160 Len=0 MSS=1460 SACK_PERM TSval=... WS=256
3    0.013981  192.168.1.132       93.184.216.34       TCP       51234 → 443 [ACK] Seq=1 Ack=1 Win=64256 Len=0
```

A detail worth explaining immediately: Wireshark shows `Seq=0` and `Ack=1`, not the huge real numbers (`1865791255`, etc.) — by default it displays **relative sequence numbers**, counting from 0 at the start of each captured stream, because that's dramatically easier to read across a whole conversation than tracking arbitrary 32-bit ISNs in your head. (Right-click → Protocol Preferences → "Relative sequence numbers" turns this off if you specifically need the raw wire values, e.g. to cross-check against a `tcpdump` capture of the same traffic.)

Expanding packet 1's TCP layer in the Packet Details pane shows exactly the fields Chapter 65 named:

```
Transmission Control Protocol, Src Port: 51234, Dst Port: 443, Seq: 0, Len: 0
    Source Port: 51234
    Destination Port: 443
    [Stream index: 0]
    Sequence Number: 0    (relative sequence number)
    Sequence Number (raw): 1865791255
    Acknowledgment Number: 0
    Header Length: 40 bytes (10)          <- Data Offset field, Ch 65 §6
    Flags: 0x002 (SYN)
        .... ..1. = SYN: Set
        .... ...0 = FIN: Not set
    Window: 64240
    Checksum: 0x8e4a [correct]
    Options: (20 bytes), Maximum segment size, No-Operation, Window scale, ...
```

Every one of these lines is Chapter 65's header table, populated with real values instead of a generic diagram: `Header Length: 40 bytes` confirms 20 bytes of fixed header plus 20 bytes of TCP options (MSS, SACK-permitted, timestamps, window scale — the exact negotiation Chapter 61's window scaling depends on). The `tcpdump` equivalent of this same packet, with `-v`, shows the same information more tersely:

```
$ sudo tcpdump -i eth0 -n -v 'tcp port 443 and tcp[tcpflags] & tcp-syn != 0' -c 1
10:41:02.001122 IP (tos 0x0, ttl 64, id 0, offset 0, flags [DF], proto TCP (6), length 60)
    192.168.1.132.51234 > 93.184.216.34.443: Flags [S], cksum 0x8e4a (correct), seq 1865791255, win 64240, options [mss 1460,sackOK,TS val 3847201 ecr 0,nop,wscale 7], length 0
```

Same bytes, same fields — `tcpdump -v` prints the raw (non-relative) sequence number by default, which is why it shows `1865791255` where Wireshark's default view showed `0`.

## 8. Reading a Captured DNS Query

Chapter 56, Section 8 introduced `dig`'s human-formatted output; here is the *actual UDP packet* `dig` (or any resolver stub) sends and receives, as Wireshark parses it — the wire format Chapter 111 later has you construct by hand in code:

```
No.  Time      Source            Destination        Protocol  Info
1    0.000000  192.168.1.132     127.0.0.53         DNS       Standard query 0x11a9 A example.com
2    0.011932  127.0.0.53        192.168.1.132       DNS       Standard query response 0x11a9 A example.com A 93.184.216.34
```

Expanding packet 1's DNS layer:

```
Domain Name System (query)
    Transaction ID: 0x11a9
    Flags: 0x0100 Standard query
        0... .... .... .... = Response: Message is a query
        .000 0... .... .... = Opcode: Standard query (0)
        .... ..0. .... .... = Recursion desired: Do query recursively   <- Ch 68 "rd" flag
    Questions: 1
    Queries
        example.com: type A, class IN
            Name: example.com
            Type: A (Host Address)                                     <- Ch 69 record type
            Class: IN (0x0001)
```

And the response:

```
Domain Name System (response)
    Transaction ID: 0x11a9              <- matches the query; this is how a resolver pairs replies (Ch 68)
    Flags: 0x8180 Standard query response, No error
        1... .... .... .... = Response: Message is a response
        .... .1.. .... .... = Recursion available
    Answers
        example.com: type A, class IN, addr 93.184.216.34
            Name: example.com
            Type: A
            Class: IN
            Time to live: 3577                                         <- the TTL, Ch 68's caching lifetime
            Address: 93.184.216.34
```

The `Transaction ID` matching between query and response is the mechanism Chapter 68 mentioned only in passing: a stub resolver may have several outstanding queries at once over the same UDP socket (Chapter 58's connectionless model means there's no per-conversation state to rely on), so this 16-bit ID — plus matching source port — is literally how a reply gets matched to its request. `tcpdump -v` shows the same query far more tersely, since DNS record parsing wasn't `tcpdump`'s original design focus:

```
$ sudo tcpdump -i eth0 -n port 53 -c 2
10:41:05.002211 IP 192.168.1.132.51900 > 127.0.0.53.53: 4521+ A? example.com. (29)
10:41:05.014033 IP 127.0.0.53.53 > 192.168.1.132.51900: 4521 1/0/1 A 93.184.216.34 (56)
```

`4521+` is the transaction ID (with `+` meaning recursion desired), `A?` is the query type and the `?` marking it as a question, and `1/0/1` on the response line means "1 answer, 0 authority records, 1 additional record" — a compact summary of exactly the message-section counts Chapter 67's hierarchy discussion assumed you'd eventually see on the wire.

## 9. Reading a Captured TLS 1.3 Handshake

Chapter 82, Section 7 walked TLS 1.3's 1-RTT handshake as a sequence diagram of named messages. Captured live, those exact message names appear as Wireshark's parsed `tls.handshake.type` values:

```
No.  Time      Source           Destination       Protocol  Info
4    0.001200  192.168.1.132    93.184.216.34     TLSv1.3   Client Hello
5    0.015310  93.184.216.34    192.168.1.132     TLSv1.3   Server Hello
6    0.015400  93.184.216.34    192.168.1.132     TLSv1.3   Change Cipher Spec, Application Data  (encrypted: Certificate, CertificateVerify, Finished)
7    0.015900  192.168.1.132    93.184.216.34     TLSv1.3   Change Cipher Spec, Application Data  (encrypted: Finished)
8    0.016400  192.168.1.132    93.184.216.34     TLSv1.3   Application Data                       (the actual HTTP request, encrypted)
```

Notice packet 6's label: after the ServerHello, TLS 1.3 encrypts everything (Chapter 82, Section 7's "Everything after ServerHello is encrypted" point) — so Wireshark can no longer parse the Certificate or CertificateVerify messages as distinct rows; it can only tell you *that* encrypted application data was exchanged, not what it contains, unless it has been given the session's decryption keys. Expanding the Client Hello (packet 4), the very first plaintext message, shows the negotiation Chapter 82 described in prose as actual captured fields:

```
Handshake Protocol: Client Hello
    Version: TLS 1.2 (0x0303)          <- legacy compatibility value; real version is in an extension
    Random: 4f2a91b3...  (32 bytes)     <- random_C from Ch 82 §5
    Cipher Suites (16 suites)
        Cipher Suite: TLS_AES_128_GCM_SHA256 (0x1301)
        Cipher Suite: TLS_CHACHA20_POLY1305_SHA256 (0x1303)
        ...
    Extension: server_name (len=16)
        Server Name: example.com                              <- the SNI field, Ch 82 §12's plaintext leak
    Extension: supported_versions (len=5)
        TLS 1.3
    Extension: key_share (len=38)
        Key Share Entry: x25519, ...                           <- the speculative key share, Ch 82 §7-8
    Extension: application_layer_protocol_negotiation (len=14)
        ALPN Protocol: h2                                       <- Ch 74/82 §13's ALPN
```

Two things worth pointing at directly, both already flagged in Chapter 82: the `Version: TLS 1.2 (0x0303)` field is not a mistake or a downgrade — TLS 1.3's actual version lives in the `supported_versions` extension for backward-compatibility reasons, so a naive reading of just that top-level field would wrongly conclude this is a TLS 1.2 handshake. And `Server Name: example.com`, sitting in plain, unencrypted view in this very first packet, is exactly the metadata leak Chapter 82, Section 12 described — anyone capturing this traffic on the path (a passive eavesdropper, Chapter 77's threat model) can see which hostname you requested even though the request and response bodies that follow are fully encrypted.

## 10. Saving and Replaying Captures: pcap Files

A capture is only useful in the moment unless it can be saved, shared, and re-examined — the standard format for that is `.pcap` (or the newer `.pcapng`, which adds richer per-packet metadata). `tcpdump -w` and Wireshark's "Save As" both write this same format, and either tool can read a file the other wrote:

```
# Capture 200 packets of one connection's traffic to a file, on a remote server over SSH
$ sudo tcpdump -i eth0 -w handshake.pcap -c 200 host 93.184.216.34

# Copy it to a local machine
$ scp user@remote-server:handshake.pcap .

# Open it in Wireshark's GUI locally, or re-read it from the command line
$ tcpdump -r handshake.pcap -n -v
```

This "capture headless, analyze visually" pattern is the standard production workflow: production servers rarely run a GUI, so the capture happens over SSH with `tcpdump -w`, and the actual investigation — clicking through the protocol tree, following a TCP stream, applying display filters — happens afterward, locally, in Wireshark. `tshark`, Wireshark's command-line sibling (same parsing engine, no GUI), can also apply full display-filter syntax directly against a saved file without ever opening a window — `tshark -r handshake.pcap -Y 'tcp.flags.syn==1'` — useful for scripting analysis of a capture on a machine with no display at all.

## 11. A Real Diagnostic Capture Walkthrough

Putting Sections 3, 6, 7, and 9 together against a real complaint: **"HTTPS requests to this one backend intermittently hang for a few seconds."**

```
Step 1 — Capture broadly around the suspect connection, on the client side
  $ sudo tcpdump -i eth0 -w hang.pcap host backend.internal and port 443

Step 2 — Reproduce the hang, then stop the capture (Ctrl-C) and open in Wireshark

Step 3 — Filter to just this stream's handshake
  Display filter: tcp.port == 443 and ip.addr == 10.20.4.11

Step 4 — Look at the timestamps between the SYN and the SYN-ACK
  No.  Time       Info
  1    0.000000   [SYN] Seq=0
  2    3.041002   [SYN, ACK] Seq=0 Ack=1     <- 3+ seconds before the reply!
  3    3.041200   [ACK]

Step 5 — Check for a retransmission first, using Wireshark's own analysis:
  Display filter: tcp.analysis.retransmission
  -> No results. The 3-second gap is real network delay, not a lost-and-resent SYN.

Step 6 — Confirm this isn't a TLS-layer stall by checking what follows the ACK
  Display filter: tls.handshake.type == 1 (ClientHello) through Finished
  -> The TLS handshake itself completes in 14ms once the TCP handshake finishes.

Conclusion: the hang is entirely inside the TCP handshake's SYN-ACK round trip,
not TLS and not the application — pointing at network path delay, a loaded
backend not accepting connections promptly, or a firewall/security-group
rule doing slow first-packet inspection, rather than anything wrong with
the application code that Chapter 122's playbook would otherwise send you
chasing at the wrong layer.
```

This is the payoff of the whole chapter: without a capture, "intermittently hangs" could be blamed on the application, DNS, TLS, or the network, almost at random. With one, the exact 3-second gap is visible, located precisely between two named packets, at one specific layer.

## 12. Hands-On Experiment

Run this on your own machine (Linux or macOS; `sudo` required for raw capture):

```bash
# Terminal 1 — start a capture, filtering to just one target and writing to a file
sudo tcpdump -i any -w /tmp/lab.pcap host example.com or host 8.8.8.8

# Terminal 2 — generate traffic across all three protocols this chapter covered
dig @8.8.8.8 example.com          # DNS, Section 8
curl -v https://example.com >/dev/null   # TCP handshake + TLS, Sections 7 and 9

# Back in Terminal 1, Ctrl-C to stop, then inspect
tcpdump -r /tmp/lab.pcap -n
```

Then open the same `/tmp/lab.pcap` in Wireshark and apply, in order: `dns`, then `tcp.flags.syn==1`, then `tls.handshake.type==1`. Confirm you can find, by eye, the exact packets this chapter's Sections 7–9 described — this is the single fastest way to convert this chapter from something read to something known.

## 13. Common Misconceptions

- **"Wireshark can decrypt any HTTPS traffic I capture."** Only if it's also given the session's symmetric keys (via a `SSLKEYLOGFILE` environment variable set on the client before capture, or an RSA private key for older non-forward-secret cipher suites) — capturing alone, as Section 9 showed, gives you only "encrypted Application Data was exchanged," never its contents. Modern TLS 1.3's exclusive use of forward-secret key exchange (Chapter 79, Chapter 82 §5) means even a server's long-term private key cannot decrypt a captured session after the fact.
- **"A capture filter and a display filter are the same thing, just typed in different boxes."** They use genuinely different languages (Section 5) — a BPF expression like `tcp[tcpflags] & tcp-syn != 0` is not valid Wireshark display-filter syntax, and `tcp.flags.syn==1` is not valid as a `tcpdump`/capture-filter expression. Mixing them up produces a filter-syntax error, not a wrong-but-working filter.
- **"If I don't see a SYN in my capture, the connection never started."** More often it means the capture started *after* the SYN, or was taken on an interface the SYN didn't cross (e.g., capturing on the server while the client's SYN was dropped somewhere upstream) — always check the capture's actual start time and location against the symptom's timing before concluding a packet never existed.
- **"tcpdump's summary line and Wireshark's parsed fields might disagree."** They read the same bytes through the same underlying decode logic (both link against similar dissection code); a genuine disagreement almost always means one tool is using a stale/older protocol dissector, not that the wire data itself is ambiguous.

## 14. Production Notes

- **Capture as narrowly as your BPF filter allows, especially in production.** An unfiltered capture on a busy production interface both risks dropping packets (if disk/CPU can't keep up) and captures far more sensitive traffic than intended — `host <target> and port <n>` should be treated as close to a minimum, not an optional refinement.
- **Rotate and cap capture file sizes.** `tcpdump -w capture.pcap -C 100 -W 10` writes 100MB files, keeping only the most recent 10 — an unbounded `-w` on a live production interface is a real, if boring, way to fill a disk.
- **`SSLKEYLOGFILE` is the sanctioned way to decrypt your own TLS traffic for debugging.** Set it before launching a browser or `curl`, point Wireshark's TLS protocol preferences at the resulting key log file, and previously-opaque `Application Data` packets become fully readable — invaluable when debugging your own service's HTTPS behavior, never usable against traffic you don't control the client for.
- **Elevated capture privileges are a real operational friction point.** `CAP_NET_RAW` (a narrower grant than full root) is the standard way to let a monitoring agent or a specific engineer run `tcpdump` without full `sudo` — worth knowing the capability exists rather than reaching for blanket root access by default.

## 15. What's Simplified Here

This chapter shows default output formats; both tools have vastly more depth than covered here — Wireshark's statistics menus (Conversations, I/O graphs, Expert Info's automatic anomaly flagging), `tcpdump`'s `-X`/`-XX` hex+ASCII dumps, and dozens more BPF primitives (`vlan`, `greater <n>`, `ip[8] == 64` for arbitrary byte tests) are all real and commonly used but out of scope here. The captures shown are illustrative reconstructions of realistic traffic, not literal byte-for-byte dumps from one specific real session — real TLS random values, sequence numbers, and timing will differ every time, which is the point: the *shape* and *field names* shown are what to expect, not literal numbers to memorize.

## 16. Interview Questions & Model Answers

**Beginner: What is the difference between `tcpdump` and Wireshark, and when would you choose one over the other?**

*Model answer:* Both capture the same raw packets via `libpcap`. `tcpdump` is command-line, text-based, and available almost everywhere by default, making it the right choice on a headless remote server over SSH or in a script. Wireshark is graphical, with a much richer protocol tree, statistics, and stream-reconstruction UI, making it better for deep visual investigation — commonly used by capturing with `tcpdump -w` remotely and opening the saved file in Wireshark locally.

**Intermediate: Write a `tcpdump` filter that captures only TCP RST packets to or from host 10.0.0.5, and explain what seeing a burst of RSTs from that host would suggest.**

*Model answer:* `sudo tcpdump -i eth0 -n 'host 10.0.0.5 and tcp[tcpflags] & tcp-rst != 0'`. A burst of RSTs suggests connections being actively refused or aborted rather than timing out — commonly a service not listening on the expected port (Chapter 64's RST-on-unmatched-segment behavior), a firewall configured to reject rather than silently drop, or an application-level abort after detecting an error, and it's worth checking whether the RSTs come from the host itself or from something on the path pretending to be it.

**Advanced: A capture shows a completed TCP three-way handshake and a completed TLS 1.3 handshake, but the HTTP response never arrives, and eventually the connection sends a RST. Using only this chapter and Chapters 59–82, what layer would you focus on next, and why?**

*Model answer:* Since both TCP (Chapter 59) and TLS (Chapter 82) completed cleanly, the failure is happening above them — at the HTTP/application layer (Chapter 71) or in the application logic behind it, not in the network path or the encryption negotiation. I'd apply a display filter on `http` or use "Follow TLS Stream" (if `SSLKEYLOGFILE` decryption is available) to see whether a request was actually sent and what, if anything, the server sent back before the RST — a RST after a completed handshake with no response often points to the server-side application crashing, timing out, or explicitly rejecting the request in a way that closes the socket abruptly rather than sending a proper HTTP error response.

## 17. Exercises

### Easy
1. Write a `tcpdump` filter (with `-i eth0 -n`) that captures only traffic to or from `10.1.2.3` on port `22`.
2. In Wireshark's display filter bar, write a filter that shows only TCP segments with the RST flag set.
3. Explain, in one sentence each, the difference between a Wireshark capture filter and a display filter.

### Medium
4. Using Section 3's BPF syntax, write a `tcpdump` filter that matches TCP segments with **both** the ACK and FIN flags set (a normal connection-close acknowledgment, Chapter 64) but excludes plain ACK-only segments.
5. You capture a DNS query and its response (Section 8) but the response's Transaction ID does not match the query's. Propose two different explanations for this mismatch, tying each to a specific mechanism from Chapter 68.
6. Explain why `tls.handshake.type == 1` (ClientHello) is visible and filterable in a TLS 1.3 capture, but a filter like `tls.handshake.type == 11` (Certificate) will typically show zero results in a TLS 1.3 session captured without a key log file — tie your answer to Section 9 and Chapter 82, Section 7.

### Hard
7. Design a `tcpdump` capture command that would run safely, unattended, on a production server for an hour, capturing only SYN packets destined for port 443 (to build a rough log of new connection attempts) without risking filling the disk. State every flag you'd add and why.
8. Given the walkthrough in Section 11, propose a *different* possible root cause than the one concluded there for the same observed 3-second SYN-to-SYN-ACK gap, and describe exactly what additional capture or tool (from this chapter or Chapter 120) you would use to distinguish your alternative from the walkthrough's conclusion.
9. A colleague claims that because Wireshark shows `Server Name: example.com` in cleartext inside a TLS 1.3 ClientHello (Section 9), TLS 1.3 offers "no real privacy benefit" over plaintext HTTP for hiding which site a user is visiting. Using Chapter 82, Section 12 and this chapter's Section 9, write a precise rebuttal that states exactly what TLS 1.3 does and does not hide, and name the specific extension (mentioned in Chapter 82) that addresses the SNI leak specifically.

## 18. Summary and the Bridge to Chapter 120

| Concept | Meaning | Chapter Connection |
|---|---|---|
| `tcpdump` | fast, text-based, scriptable packet capture | built on the raw capture Ch 114 has you implement |
| BPF filter syntax | `host`, `port`, `tcp[tcpflags] & tcp-syn` | reads Ch 65's header byte offsets directly |
| Wireshark | graphical capture/analysis with a protocol tree | parses every field from Ch 27, 65, 69, 82 |
| Capture filter | BPF syntax, applied before capture | same language as `tcpdump` |
| Display filter | `tcp.flags.syn==1`, dotted field names | Wireshark-specific, applied after capture |
| pcap / pcapng | saved capture file format | shared between `tcpdump -w`/`-r` and Wireshark |
| Follow TCP Stream | reassembled application-layer conversation | undoes Ch 60's segmentation for reading |

You can now see, with your own eyes, every protocol this course has built up through Chapter 118 — but seeing one packet, or even one connection, is not the same as knowing whether a network is healthy *right now*, under real load, over time. That's a different kind of question — not "what does this packet say" but "how fast, how consistent, and how lossy is this path" — and it needs different tools and precise definitions of terms like latency, throughput, jitter, and packet loss that get thrown around loosely in everyday conversation. Chapter 120 defines each one precisely and hands you the tools — `ping`, `iperf3`, and `mtr` — built specifically to measure them.
