# Chapter 111: Building a DNS Resolver

> **"Chapter 68 described a resolver walking from root to TLD to authoritative server as if narrating a trip. This chapter drives the car."**

---

## Table of Contents

1. [Recap: What Chapters 66-69 Described in Prose](#1-recap-what-chapters-66-69-described-in-prose)
2. [The Problem: Turning a Name Into an Address, One Server at a Time](#2-the-problem-turning-a-name-into-an-address-one-server-at-a-time)
3. [The Naive Shortcut We're Deliberately Not Taking](#3-the-naive-shortcut-were-deliberately-not-taking)
4. [The DNS Wire Format, Field by Field](#4-the-dns-wire-format-field-by-field)
5. [Code: Encoding a Query By Hand](#5-code-encoding-a-query-by-hand)
6. [Parsing a Response: Headers, Names, and the Compression Pointer Problem](#6-parsing-a-response-headers-names-and-the-compression-pointer-problem)
7. [Code: Decoding Names, Questions, and Resource Records](#7-code-decoding-names-questions-and-resource-records)
8. [The Recursive Walk: Root → TLD → Authoritative, in Code](#8-the-recursive-walk-root--tld--authoritative-in-code)
9. [Code: The Complete Resolver](#9-code-the-complete-resolver)
10. [Hands-On Experiment: Running It, and Checking Against `dig +trace`](#10-hands-on-experiment-running-it-and-checking-against-dig-trace)
11. [Common Pitfalls in Hand-Rolled DNS Parsing](#11-common-pitfalls-in-hand-rolled-dns-parsing)
12. [Production Notes: What Real Recursive Resolvers Do Differently](#12-production-notes-what-real-recursive-resolvers-do-differently)
13. [What's Simplified Here](#13-whats-simplified-here)
14. [Interview Questions & Model Answers](#interview-questions--model-answers)
15. [Exercises](#exercises)
16. [Summary](#summary)

---

## 1. Recap: What Chapters 66-69 Described in Prose

Chapter 67 described DNS's three tiers of servers — root, TLD, authoritative — and delegation between them. Chapter 68 described, in prose, what a recursive resolver does on a client's behalf: start at a root server, follow a referral to the right TLD server, follow another referral to the right authoritative server, and return the final answer, caching results along the way. Chapter 69 described the actual record types (A, NS, CNAME, and others) that flow through that process, plus the security and privacy add-ons (DNSSEC, DoH, DoT) layered on top of it.

None of those three chapters showed a single byte of what a DNS message actually looks like, because that was deliberately left for this chapter. Everything here is real: a real 12-byte binary header, real length-prefixed name encoding, a real UDP socket, and a real implementation of the root→TLD→authoritative walk — talking directly to root and TLD servers, never asking a public recursive resolver like `8.8.8.8` to do the work for us.

---

## 2. The Problem: Turning a Name Into an Address, One Server at a Time

Stated precisely: given a domain name like `example.com`, and starting with nothing but the well-known IP addresses of the 13 root server letters (Chapter 67), produce the correct IPv4 address for that name — by sending exactly the right bytes over UDP to exactly the right server at each step, correctly interpreting whatever comes back (an answer, or a referral to more specific servers), and repeating until an answer arrives.

Two sub-problems make this hard enough to deserve its own chapter. First, **encoding**: a DNS message is not text like HTTP — it's a binary format with a fixed-size header, bit-packed flags, and a name encoding scheme (length-prefixed labels) that doesn't look like anything in Chapters 109 or 110. Second, **decoding a referral correctly**: a real DNS response uses **compression pointers** to avoid repeating the same domain name (like `.com`) dozens of times in one packet, and a parser that doesn't understand those pointers will either crash or silently produce garbage names the moment it meets one — which is almost immediately, since even the root server's very first referral response uses them.

---

## 3. The Naive Shortcut We're Deliberately Not Taking

```go
ips, err := net.LookupHost("example.com")
```

This works, and is exactly what you should call in real Go programs. But it's a black box in two different ways at once. First, it hides the wire format entirely — no header, no QNAME encoding, no compression pointers, just a `[]string` of addresses. Second, and more specific to this chapter's ambition: `net.LookupHost` (in its default configuration) doesn't do the root→TLD→authoritative walk itself at all — it sends one query to whatever recursive resolver is configured in `/etc/resolv.conf` (often your router, or `8.8.8.8`, or `1.1.1.1`) and lets *that* server do the walking, entirely out of view. This chapter's resolver does the walking itself, in code you can read, acting as its own recursive resolver rather than asking someone else's.

---

## 4. The DNS Wire Format, Field by Field

Every DNS message — query or response — starts with the same fixed 12-byte header:

```
 0                   1                   2                   3
 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1 2 3 4 5 6 7 8 9 0 1
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                      ID (16 bits)                            |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|QR|   Opcode  |AA|TC|RD|RA|   Z    |    RCODE      |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    QDCOUNT (16 bits)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    ANCOUNT (16 bits)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    NSCOUNT (16 bits)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
|                    ARCOUNT (16 bits)                         |
+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+-+
```

| Field | Size | Meaning |
|---|---|---|
| ID | 16 bits | Transaction ID — the client picks it, the server echoes it back, so a UDP-based client can match a response to the right in-flight query |
| QR | 1 bit | 0 = this is a query, 1 = this is a response |
| Opcode | 4 bits | 0 = standard query (the only kind this chapter sends) |
| AA | 1 bit | Authoritative Answer — set when the responder owns the zone, not just relaying it |
| TC | 1 bit | Truncated — the answer didn't fit in this UDP response; retry over TCP |
| RD | 1 bit | Recursion Desired — set by a client asking a server to do the full walk itself |
| RA | 1 bit | Recursion Available — set by a server willing to do that walk |
| Z | 3 bits | Reserved, must be zero |
| RCODE | 4 bits | Response code: 0 = no error, 3 = NXDOMAIN (name doesn't exist), others for server failure etc. |
| QDCOUNT | 16 bits | Number of entries in the question section |
| ANCOUNT | 16 bits | Number of resource records in the answer section |
| NSCOUNT | 16 bits | Number of resource records in the authority section |
| ARCOUNT | 16 bits | Number of resource records in the additional section |

After the header comes the **question section**: one or more `(QNAME, QTYPE, QCLASS)` triples. QNAME is the part with no HTTP analog — a domain name encoded as a sequence of **length-prefixed labels**, terminated by a zero-length label. `www.example.com` becomes:

```
[3]www[7]example[3]com[0]
 ^   ^^^  ^   ^^^^^^^  ^   ^^^  ^
 |   |     |    |        |    |   |
 len "www" len "example"  len "com" terminator (length 0)
```

As raw bytes: `03 77 77 77 07 65 78 61 6D 70 6C 65 03 63 6F 6D 00` — each label's length byte tells the parser exactly how many bytes of text follow, so there's never any ambiguity about where one label ends and the next begins, and the final `00` byte marks the end of the name (the root zone itself, which has zero labels).

After the question section, a **response** (never a query) adds up to three more sections — answer, authority, and additional — each containing zero or more **resource records**, all sharing one format:

```
NAME (possibly compressed) | TYPE (2B) | CLASS (2B) | TTL (4B) | RDLENGTH (2B) | RDATA (RDLENGTH bytes)
```

`RDATA`'s meaning depends on `TYPE`: for an A record (`TYPE=1`), it's exactly 4 raw bytes — the IPv4 address itself, no text encoding at all. For an NS record (`TYPE=2`), `RDATA` is itself an encoded domain name — the nameserver's hostname — which can, and in real responses usually does, use the compression pointer trick covered next.

---

## 5. Code: Encoding a Query By Hand

```go
// encodeName turns "www.example.com" into DNS wire format: a sequence of
// length-prefixed labels, terminated by a zero-length label (Section 4).
func encodeName(name string) []byte {
	name = strings.TrimSuffix(name, ".")
	if name == "" {
		return []byte{0} // the root domain itself
	}
	var buf []byte
	for _, label := range strings.Split(name, ".") {
		buf = append(buf, byte(len(label)))
		buf = append(buf, []byte(label)...)
	}
	buf = append(buf, 0)
	return buf
}

// buildQuery constructs a complete 12-byte header plus one question section.
func buildQuery(id uint16, name string, qtype uint16, recursionDesired bool) []byte {
	var flags uint16
	if recursionDesired {
		flags = 0x0100 // RD is bit 8 of the flags word (Section 4's diagram)
	}
	buf := make([]byte, 12)
	binary.BigEndian.PutUint16(buf[0:2], id)
	binary.BigEndian.PutUint16(buf[2:4], flags)
	binary.BigEndian.PutUint16(buf[4:6], 1) // QDCOUNT = 1 — ANCOUNT/NSCOUNT/ARCOUNT stay 0 in a query

	buf = append(buf, encodeName(name)...)
	tail := make([]byte, 4)
	binary.BigEndian.PutUint16(tail[0:2], qtype)
	binary.BigEndian.PutUint16(tail[2:4], ClassIN)
	return append(buf, tail...)
}
```

`buildQuery(0x1234, "example.com", TypeA, false)` produces exactly:

```
12 34            ID = 0x1234
00 00            flags = 0 (QR=0 query, RD=0 — we walk the hierarchy ourselves)
00 01            QDCOUNT = 1
00 00 00 00 00 00   ANCOUNT, NSCOUNT, ARCOUNT = 0
07 65 78 61 6D 70 6C 65   "example" (length 7)
03 63 6F 6D              "com" (length 3)
00                        terminating zero label
00 01            QTYPE = 1 (A)
00 01            QCLASS = 1 (IN)
```

29 bytes total — the entire query, ready to write directly onto a UDP socket.

---

## 6. Parsing a Response: Headers, Names, and the Compression Pointer Problem

A response's header decodes with plain fixed-offset reads — no surprises. The surprise is in decoding names inside the answer/authority/additional sections. A root server's referral for `example.com` needs to mention `.com`'s TLD nameservers — `a.gtld-servers.net`, `b.gtld-servers.net`, and so on — a dozen or more times across the authority and additional sections. Repeating the full, uncompressed encoding of `gtld-servers.net` a dozen times would waste a meaningful fraction of a UDP packet's limited size. DNS's fix: a **compression pointer**.

Whenever a length byte's top two bits are both `1` (a real label length is 0–63, so this pattern — `0xC0` or higher — can never be a real length), the parser instead treats that byte plus the next byte as a 14-bit **offset from the start of the entire message**, and the name continues from there:

```
byte value 0xC0..0xFF at position P:
   top 2 bits = 11  -> this is a POINTER, not a length
   pointer = ((msg[P] & 0x3F) << 8) | msg[P+1]
   -> jump to byte `pointer` in the SAME message and keep reading labels from there
```

So a response might encode `ns1.example.com` as `[3]ns1` followed immediately by a pointer back to wherever `example.com` was already spelled out earlier in the same packet — no repetition, and the message stays inside a single UDP datagram. A parser that doesn't check for this pattern will read a compression-pointer byte (e.g. `0xC0`) as if it were a 192-byte label length, immediately overrun the buffer, and produce a corrupted or crashing parse — on the very first real-world response, since root and TLD responses use compression pervasively.

---

## 7. Code: Decoding Names, Questions, and Resource Records

```go
// decodeName reads a (possibly compressed) domain name starting at offset
// within msg (RFC 1035 Sec 4.1.4). It returns the decoded name and the
// offset of the first byte AFTER the name in the ORIGINAL stream — which,
// if the name used a pointer, is NOT the same as where the pointer jumped to.
func decodeName(msg []byte, offset int) (string, int, error) {
	var labels []string
	pos := offset
	resumeAt := -1
	jumps := 0

	for {
		if pos >= len(msg) {
			return "", 0, errors.New("dns: name runs past end of message")
		}
		length := int(msg[pos])

		if length == 0 { // terminating zero-length label
			pos++
			if resumeAt == -1 {
				resumeAt = pos
			}
			break
		}

		if length&0xC0 == 0xC0 { // compression pointer (Section 6)
			if pos+1 >= len(msg) {
				return "", 0, errors.New("dns: truncated compression pointer")
			}
			if resumeAt == -1 {
				resumeAt = pos + 2 // caller resumes right after this 2-byte pointer
			}
			pointer := (length&0x3F)<<8 | int(msg[pos+1])
			jumps++
			if jumps > 20 {
				return "", 0, errors.New("dns: too many compression jumps (possible loop)")
			}
			pos = pointer
			continue
		}

		pos++ // plain label: length byte, then that many bytes of text
		if pos+length > len(msg) {
			return "", 0, errors.New("dns: label runs past end of message")
		}
		labels = append(labels, string(msg[pos:pos+length]))
		pos += length
	}
	return strings.Join(labels, "."), resumeAt, nil
}
```

The `resumeAt` variable is the subtle part: once a name follows a pointer, the *caller's* next read must resume right after the 2-byte pointer in the original stream — not at wherever the pointer jumped to, which might be in the middle of an entirely different record. `resumeAt` is captured the first time a pointer is followed (or the name ends without one) and never overwritten again, exactly so the jump target's own internal position doesn't leak out as the wrong "next offset."

Resource records and questions build directly on `decodeName`:

```go
type resourceRecord struct {
	Name        string
	Type        uint16
	Class       uint16
	TTL         uint32
	RDLength    uint16
	RData       []byte
	RDataOffset int // absolute offset of RDATA's first byte — needed to decode a
	                // name (NS/CNAME target) living INSIDE rdata, since a pointer
	                // there is an offset into the whole message, not into RData alone
}

func decodeRR(msg []byte, offset int) (resourceRecord, int, error) {
	name, offset, err := decodeName(msg, offset)
	if err != nil {
		return resourceRecord{}, 0, err
	}
	if offset+10 > len(msg) {
		return resourceRecord{}, 0, errors.New("dns: resource record header truncated")
	}
	rr := resourceRecord{
		Name:     name,
		Type:     binary.BigEndian.Uint16(msg[offset : offset+2]),
		Class:    binary.BigEndian.Uint16(msg[offset+2 : offset+4]),
		TTL:      binary.BigEndian.Uint32(msg[offset+4 : offset+8]),
		RDLength: binary.BigEndian.Uint16(msg[offset+8 : offset+10]),
	}
	offset += 10
	if offset+int(rr.RDLength) > len(msg) {
		return resourceRecord{}, 0, errors.New("dns: RDATA truncated")
	}
	rr.RDataOffset = offset
	rr.RData = msg[offset : offset+int(rr.RDLength)]
	return rr, offset + int(rr.RDLength), nil
}
```

`RDataOffset` matters because an NS record's `RDATA` is a domain name that may itself contain a compression pointer — and that pointer, like every pointer in the message, is an offset from the *whole message's* start, not from the start of this one record's `RDATA` slice. Section 9's resolver decodes NS targets with `decodeName(msg.Raw, ns.RDataOffset)` for exactly this reason.

---

## 8. The Recursive Walk: Root → TLD → Authoritative, in Code

```mermaid
sequenceDiagram
    participant R as resolve("example.com")
    participant Root as Root server (198.41.0.4)
    participant TLD as .com TLD server
    participant Auth as example.com authoritative server

    R->>Root: query A example.com (RD=0)
    Root-->>R: no answer; referral: NS + glue A for .com servers
    R->>TLD: query A example.com (RD=0)
    TLD-->>R: no answer; referral: NS + glue A for example.com's own servers
    R->>Auth: query A example.com (RD=0)
    Auth-->>R: ANSWER: A record, e.g. 93.184.215.14
```

Each hop follows the same decision tree: does the response's **answer** section have what we asked for? If yes, done (or, for a CNAME, restart resolution for the alias target — Chapter 69's CNAME chasing, made concrete). If no, does the **authority** section contain NS records for a more specific zone? If yes, look in the **additional** section for their IP addresses ("glue" — provided specifically so the resolver doesn't have to separately look up the nameserver's own address, which would be a circular problem if the nameserver's name is itself inside the domain being delegated, e.g. `ns1.example.com` for `example.com`). Use those addresses as the next set of servers to query, and repeat.

---

## 9. Code: The Complete Resolver

```go
// dnsresolver.go
package main

import (
	"encoding/binary"
	"errors"
	"fmt"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	TypeA     uint16 = 1
	TypeNS    uint16 = 2
	TypeCNAME uint16 = 5
	ClassIN   uint16 = 1
)

// A small subset of the real 13 root server letters (Ch 67) — enough to
// demonstrate the walk; a production resolver ships all 13 for redundancy.
var rootServers = []string{
	"198.41.0.4",   // a.root-servers.net
	"199.9.14.201", // b.root-servers.net
	"192.33.4.12",  // c.root-servers.net
	"199.7.91.13",  // d.root-servers.net
}

type header struct {
	ID, Flags, QDCount, ANCount, NSCount, ARCount uint16
}

type question struct {
	Name  string
	Type  uint16
	Class uint16
}

type resourceRecord struct {
	Name        string
	Type        uint16
	Class       uint16
	TTL         uint32
	RDLength    uint16
	RData       []byte
	RDataOffset int
}

type message struct {
	Header      header
	Questions   []question
	Answers     []resourceRecord
	Authorities []resourceRecord
	Additionals []resourceRecord
	Raw         []byte
}

// -- encodeName, buildQuery: Section 5 --
// -- decodeName: Section 7 --

func decodeQuestion(msg []byte, offset int) (question, int, error) {
	name, offset, err := decodeName(msg, offset)
	if err != nil {
		return question{}, 0, err
	}
	if offset+4 > len(msg) {
		return question{}, 0, errors.New("dns: question section truncated")
	}
	q := question{
		Name:  name,
		Type:  binary.BigEndian.Uint16(msg[offset : offset+2]),
		Class: binary.BigEndian.Uint16(msg[offset+2 : offset+4]),
	}
	return q, offset + 4, nil
}

// -- decodeRR: Section 7 --

func decodeMessage(msg []byte) (*message, error) {
	if len(msg) < 12 {
		return nil, errors.New("dns: message shorter than a header")
	}
	h := header{
		ID:      binary.BigEndian.Uint16(msg[0:2]),
		Flags:   binary.BigEndian.Uint16(msg[2:4]),
		QDCount: binary.BigEndian.Uint16(msg[4:6]),
		ANCount: binary.BigEndian.Uint16(msg[6:8]),
		NSCount: binary.BigEndian.Uint16(msg[8:10]),
		ARCount: binary.BigEndian.Uint16(msg[10:12]),
	}
	m := &message{Header: h, Raw: msg}
	offset := 12

	for i := 0; i < int(h.QDCount); i++ {
		q, next, err := decodeQuestion(msg, offset)
		if err != nil {
			return nil, err
		}
		m.Questions = append(m.Questions, q)
		offset = next
	}
	readRRs := func(count int) ([]resourceRecord, error) {
		var rrs []resourceRecord
		for i := 0; i < count; i++ {
			rr, next, err := decodeRR(msg, offset)
			if err != nil {
				return nil, err
			}
			rrs = append(rrs, rr)
			offset = next
		}
		return rrs, nil
	}
	var err error
	if m.Answers, err = readRRs(int(h.ANCount)); err != nil {
		return nil, err
	}
	if m.Authorities, err = readRRs(int(h.NSCount)); err != nil {
		return nil, err
	}
	if m.Additionals, err = readRRs(int(h.ARCount)); err != nil {
		return nil, err
	}
	return m, nil
}

func queryServer(name string, qtype uint16, server string) (*message, error) {
	conn, err := net.Dial("udp", net.JoinHostPort(server, "53"))
	if err != nil {
		return nil, err
	}
	defer conn.Close()
	conn.SetDeadline(time.Now().Add(3 * time.Second))

	id := uint16(rand.Intn(65536))
	// RD=false: we walk the hierarchy ourselves (Section 8); root/TLD servers
	// would ignore RD=true anyway, since they never recurse for anyone.
	if _, err := conn.Write(buildQuery(id, name, qtype, false)); err != nil {
		return nil, err
	}

	buf := make([]byte, 4096) // comfortable margin over the classic 512-byte UDP limit
	n, err := conn.Read(buf)
	if err != nil {
		return nil, err
	}
	msg, err := decodeMessage(buf[:n])
	if err != nil {
		return nil, err
	}
	if msg.Header.ID != id {
		return nil, fmt.Errorf("dns: transaction ID mismatch (sent %d, got %d)", id, msg.Header.ID)
	}
	return msg, nil
}

// resolve walks root -> TLD -> authoritative, exactly as Section 8's diagram
// shows, asking one real server directly at every step.
func resolve(name string) (net.IP, error) {
	servers := rootServers
	for hop := 0; hop < 20; hop++ {
		if len(servers) == 0 {
			return nil, errors.New("dns: no servers left to query")
		}

		var msg *message
		var lastErr error
		var usedServer string
		for _, server := range servers {
			if msg, lastErr = queryServer(name, TypeA, server); lastErr == nil {
				usedServer = server
				break
			}
		}
		if msg == nil {
			return nil, fmt.Errorf("dns: no server at this level responded: %w", lastErr)
		}
		fmt.Printf("  queried %-15s -> %d answer(s), %d authority, %d additional\n",
			usedServer, len(msg.Answers), len(msg.Authorities), len(msg.Additionals))

		for _, ans := range msg.Answers {
			switch ans.Type {
			case TypeA:
				return net.IP(ans.RData), nil
			case TypeCNAME:
				target, _, err := decodeName(msg.Raw, ans.RDataOffset)
				if err != nil {
					return nil, err
				}
				fmt.Printf("  %s is a CNAME for %s -- restarting resolution\n", name, target)
				return resolve(target)
			}
		}

		if len(msg.Authorities) == 0 {
			return nil, fmt.Errorf("dns: %s gave neither an answer nor a referral", usedServer)
		}

		nsNames := make(map[string]bool)
		for _, ns := range msg.Authorities {
			if ns.Type != TypeNS {
				continue
			}
			target, _, err := decodeName(msg.Raw, ns.RDataOffset)
			if err != nil {
				return nil, err
			}
			nsNames[strings.ToLower(target)] = true
		}

		var nextServers []string
		for _, add := range msg.Additionals {
			if add.Type == TypeA && nsNames[strings.ToLower(add.Name)] {
				nextServers = append(nextServers, net.IP(add.RData).String())
			}
		}
		if len(nextServers) == 0 {
			// No glue supplied — resolve one nameserver's own address first.
			for nsName := range nsNames {
				if ip, err := resolve(nsName); err == nil {
					nextServers = append(nextServers, ip.String())
					break
				}
			}
		}
		if len(nextServers) == 0 {
			return nil, errors.New("dns: referral had no usable nameserver addresses")
		}
		servers = nextServers
	}
	return nil, errors.New("dns: too many referrals, giving up")
}

func main() {
	name := "example.com"
	fmt.Printf("resolving %s by walking the hierarchy ourselves:\n", name)
	ip, err := resolve(name)
	if err != nil {
		fmt.Println("resolution failed:", err)
		return
	}
	fmt.Printf("%s -> %s\n", name, ip)
}
```

(`encodeName`, `buildQuery`, and `decodeName` are exactly the code already shown in Sections 5 and 7 — included in the real file, elided here only to avoid repeating them a third time.)

---

## 10. Hands-On Experiment: Running It, and Checking Against `dig +trace`

**Step 1 — run the hand-rolled resolver** (this needs outbound UDP port 53 to be permitted, which is normal on most home and cloud networks but sometimes blocked on corporate ones):

```
$ go run dnsresolver.go
resolving example.com by walking the hierarchy ourselves:
  queried 198.41.0.4     -> 0 answer(s), 13 authority, 14 additional
  queried 192.5.6.30     -> 0 answer(s), 2 authority, 2 additional
  queried 199.43.135.53  -> 1 answer(s), 0 authority, 0 additional
example.com -> 93.184.215.14
```

Read this trace against Section 8's decision tree: the root server (`198.41.0.4`, one of the four hardcoded root IPs) had no answer for `example.com`, but its authority section listed 13 NS records — the `.com` TLD's own nameservers, `a.gtld-servers.net` through `m.gtld-servers.net` — with 14 additional (glue) A records supplying their addresses directly, no extra lookup needed. The resolver picked one, `192.5.6.30` (`a.gtld-servers.net`), which again had no answer but referred to `example.com`'s own two authoritative nameservers (`a.iana-servers.net` and `b.iana-servers.net`) with their glue addresses. Querying `199.43.135.53` (`a.iana-servers.net`) directly finally produced a real answer section with one A record.

**Step 2 — compare against the real tool built for exactly this trace:**

```
$ dig +trace example.com

; <<>> DiG 9.18.24 <<>> +trace example.com
;; global options: +cmd
.            518400  IN  NS  a.root-servers.net.
...
com.         172800  IN  NS  a.gtld-servers.net.
...
example.com. 172800  IN  NS  a.iana-servers.net.
example.com. 172800  IN  NS  b.iana-servers.net.
...
example.com.     86400   IN  A   93.184.215.14
;; Received 56 bytes from 199.43.135.53#53(a.iana-servers.net) in 22 ms
```

Same shape, same three hops, same final answer — `dig +trace` is doing precisely what Section 9's `resolve()` does, with far more polish (it also handles IPv6 glue, DNSSEC validation flags, and TCP fallback on truncation).

**Step 3 — deliberately corrupt a compression pointer to see the guard rail work.** Editing `decodeName` to skip the `jumps > 20` check and feeding it a two-byte pointer pair that points at itself (`C0 0C` at offset `0x0C`) would otherwise loop forever; with the check in place:

```
resolution failed: dns: too many compression jumps (possible loop)
```

---

## 11. Common Pitfalls in Hand-Rolled DNS Parsing

- **Not checking for compression pointers at all.** The single most common bug: treating every length byte as a real label length. Since a compression pointer's first byte is always `0xC0`–`0xFF` (192–255), and no real label is longer than 63 bytes, a parser that skips the `length&0xC0==0xC0` check will try to read a 192+-byte "label" and immediately run past the buffer on the very first real-world response.
- **Resuming from the wrong offset after following a pointer.** `decodeName`'s `resumeAt` variable exists specifically because the *next* record in the message continues right after the 2-byte pointer that was written in the stream — not at whatever earlier offset the pointer jumped to. Getting this backwards makes every record after the first compressed name parse from the wrong position.
- **Treating `RDataOffset` as relative to `RData` instead of the whole message.** A pointer inside an NS record's RDATA is still an absolute offset from byte 0 of the entire DNS message (Section 6) — decoding it with `decodeName(rr.RData, 0)` instead of `decodeName(msg.Raw, rr.RDataOffset)` looks plausible but silently produces wrong names the moment that RDATA happens to use a pointer.
- **Comparing domain names case-sensitively.** DNS names are case-insensitive; a server is free to echo `EXAMPLE.com` or `example.COM`. `resolve()`'s `strings.ToLower` on both sides of the `nsNames` lookup avoids silently failing to match glue records to their NS targets due to a casing mismatch that has no actual bearing on which name is meant.
- **Trusting the response without checking the transaction ID.** UDP has no built-in way to know a packet actually came from the server you sent to (Chapter 58) — the transaction ID is the only correlation DNS provides, and `queryServer`'s `msg.Header.ID != id` check is a minimal defense against accepting a stray or spoofed packet as the real answer.

---

## 12. Production Notes: What Real Recursive Resolvers Do Differently

- **Caching, aggressively, at every level (Chapter 68).** This resolver re-walks the entire hierarchy on every call. A real recursive resolver caches every record it sees — including the root and TLD referrals — respecting each record's TTL, so a second lookup for anything under `.com` skips the root query entirely for as long as that referral's TTL remains valid. This is the single biggest reason DNS is fast enough to be invisible in practice.
- **Cryptographically random transaction IDs and source ports.** `rand.Intn` here uses Go's default (not cryptographically secure) generator, and relies on `net.Dial` picking an arbitrary ephemeral source port. Chapter 83's discussion of DNS cache poisoning is directly about attacking exactly these two values when they're too predictable — production resolvers use `crypto/rand` and validate both.
- **All 13 root letters, plus IPv6 glue and fallback to TCP on truncation (`TC=1`).** This chapter's 4-server list and IPv4-only glue handling are deliberately trimmed for readability.
- **DNSSEC validation (Chapter 69).** A real security-conscious resolver checks signatures on every step of the chain, not just trusting whatever the network handed back — this resolver has no defense against a compromised or spoofed server returning a wrong answer.
- **Concurrency and timeouts tuned per real-world failure modes**, retrying different root/TLD servers in parallel rather than this chapter's simple sequential loop, and handling `SERVFAIL`/`NXDOMAIN` response codes explicitly rather than only branching on the presence of answers/referrals.

---

## 13. What's Simplified Here

This resolver looks up A records only (no AAAA, MX, TXT, or any other type from Chapter 69's list), ignores DNSSEC entirely, has no cache, uses only 4 of the 13 real root servers, does not fall back to TCP when a response is truncated (`TC=1`), and does not validate `RCODE` (a real `NXDOMAIN` response would currently be silently treated as "no answer, no referral" and reported as a generic failure rather than "this domain doesn't exist"). Each of these is a real, well-understood gap between a teaching implementation and a production one — not a shortcut that changes the core wire-format or walking logic this chapter set out to prove correct.

---

## Interview Questions & Model Answers

**Beginner: Why does DNS use length-prefixed labels (`[3]www[7]example[3]com[0]`) instead of just sending the domain name as a plain string with dots, the way a browser displays it?**

A length-prefixed encoding lets the parser know exactly how many bytes belong to each label without needing a special delimiter character (like a dot) that could theoretically appear inside label data, and without needing to scan ahead to find where a label ends. The parser just reads one length byte, then reads exactly that many bytes as the label, then reads the next length byte — a completely unambiguous, self-describing format that a zero-length byte can also cleanly terminate.

**Intermediate: Explain, using this chapter's code, exactly what a compression pointer is and why root and TLD server responses use them so heavily.**

A compression pointer is a 2-byte sequence, identified by its first byte's top two bits both being `1` (a pattern no real label length, capped at 63, can produce), that encodes a 14-bit offset back into the same DNS message where a domain name was already spelled out earlier. Root and TLD responses reference the same handful of parent-zone names (like `.com`, or `gtld-servers.net`) repeatedly across many NS and glue A records in one response; without compression, each repetition would need the full uncompressed encoding, potentially pushing the response past DNS's traditional 512-byte UDP limit. `decodeName`'s pointer-following logic — jumping to the referenced offset, reading labels from there, but reporting the *original* stream position as where the next field resumes — is exactly what lets a parser handle this without needing to special-case which fields might be compressed.

**Advanced: This resolver's `resolve()` function queries root and TLD servers with `RD=0` (Recursion Desired unset). Explain why, and what would (and wouldn't) go wrong if it set `RD=1` instead.**

Setting `RD=1` asks the server being queried to perform full recursive resolution itself and return only the final answer — behavior that root and TLD (authoritative-only) servers don't implement at all; they always answer with only what they directly know (a referral, in this case) regardless of the RD bit's value, since the RA (Recursion Available) bit in their response would correctly report "no" either way. So functionally, nothing would break — the responses would be identical. The reason to set `RD=0` deliberately is correctness of intent: this resolver is explicitly choosing to do the walk itself (Section 8), and a query correctly stating "I am not asking you to recurse" documents that intent honestly, matching what a real recursive resolver's own outbound queries to root and TLD servers look like — they too set `RD=0` on these upstream queries, reserving `RD=1` for the query a client sends *to* the recursive resolver itself.

---

## Exercises

### Easy
1. Change `main()` to resolve a different domain (e.g. `wikipedia.org`) and confirm the trace still shows the same three-level shape.
2. Add a case for `RCODE != 0` in `resolve()` that reports `NXDOMAIN` (RCODE 3) with a clear message instead of falling through to "gave neither an answer nor a referral."
3. Print the TTL of the final A record returned, and explain in a comment what it would mean for caching (Chapter 68) if this resolver kept results around.

### Medium
4. Add AAAA (IPv6, `TYPE=28`) support: a `resolveAAAA` variant of `resolve` that requests `TypeAAAA` and decodes a 16-byte `RData` into a `net.IP`.
5. Add a simple in-memory cache keyed by `(name, qtype)` that stores answers along with an expiration time computed from the record's TTL, and skip the network entirely on a cache hit within that window.
6. Add TCP fallback: if a response has the `TC` bit set (bit 9 of the flags word), reissue the same query over a TCP connection to the same server instead of trusting the truncated UDP response.

### Hard
7. Add all 13 real root server IPs (look up the current root hints file) and randomize which one is tried first on each top-level call, the way a real resolver load-balances across the root server set.
8. Implement basic negative caching: when a server returns `RCODE=3` (NXDOMAIN), cache that fact (with a TTL derived from the SOA record in the authority section) so repeated lookups for a nonexistent name don't repeat the full walk.
9. Harden `queryServer` against a spoofed response by also verifying the response's single question section echoes back the exact name and type that was queried, not just the transaction ID — and explain in a comment which specific attack this defends against (tie to Chapter 83's DNS cache poisoning discussion).

---

## Summary

| Term | Meaning |
|---|---|
| DNS header | Fixed 12-byte block: ID, bit-packed flags, and four section counts |
| QNAME encoding | Domain name as length-prefixed labels, terminated by a zero-length label |
| Compression pointer | A 2-byte back-reference (top two bits `11`) replacing a repeated name with an offset into the same message |
| Resource record (RR) | `NAME, TYPE, CLASS, TTL, RDLENGTH, RDATA` — the shared shape of every answer/authority/additional entry |
| Glue record | An A record supplied alongside an NS referral so the resolver doesn't need a separate lookup for the nameserver's own address |
| Referral | A response with no answer but authority/additional records pointing to more specific nameservers |
| Iterative walk | Querying root, then TLD, then authoritative servers directly yourself, instead of asking a recursive resolver to do it |

You've now built, by hand, the exact process Chapter 68 only described in prose — a resolver that speaks DNS's real binary wire format and walks the real global hierarchy to turn a name into an address. Chapter 112 moves to the next stage of a request's life once a name has resolved to an IP: building a reverse proxy that accepts a connection meant for one address and forwards it, correctly relabeled, to a backend server behind it.
