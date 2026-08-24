# Chapter 41: Networking — TCP/UDP, HTTP/1.1 to HTTP/3, TLS & WebSockets

Senior engineers must understand the full networking stack: how TCP works, why HTTP/2 was invented, what TLS actually does, and how WebSockets differ from HTTP. These topics appear in system design and deep-dive interviews.

## Table of Contents

1. [TCP vs UDP](#1-tcp-vs-udp)
2. [How TCP Works (Deeply)](#2-how-tcp-works-deeply)
3. [HTTP/1.1, HTTP/2, and HTTP/3](#3-http11-http2-and-http3)
4. [TLS & HTTPS — What Actually Happens](#4-tls--https--what-actually-happens)
5. [DNS Resolution](#5-dns-resolution)
6. [CDN Internals](#6-cdn-internals)
7. [Interview Questions & Model Answers](#7-interview-questions--model-answers)
8. [Summary](#summary)

---

## 1. TCP vs UDP

| | TCP | UDP |
|---|---|---|
| Connection | Required (3-way handshake) | None |
| Reliability | Guaranteed delivery, ordered | No guarantees, no ordering |
| Error detection | Retransmission on loss | No retransmission |
| Flow control | Yes (prevents sender from overwhelming receiver) | No |
| Congestion control | Yes (backs off under network congestion) | No |
| Latency | Higher (handshake + retransmission) | Lower |
| Use cases | HTTP, SSH, databases, file transfer | Video streaming, DNS, gaming, VoIP |

**When UDP is better:** when losing some packets is acceptable but latency is critical. Video calls: a lost frame is better than stuttering. Online gaming: stale position is better than lag. DNS: a single packet request, can retry easily.

---

## 2. How TCP Works (Deeply)

### Three-Way Handshake

```
Client                          Server

SYN (seq=100) ──────────────▶
              ◀────────────── SYN-ACK (seq=500, ack=101)
ACK (ack=501) ──────────────▶

Connection established! Data can flow in both directions.

Each step takes 1 RTT (round trip time)
Total handshake: 1.5 RTT before first data byte
```

### Sequence Numbers and Reliability

```
TCP assigns a sequence number to every byte:
  Sender: "I'm sending bytes 1000-1999 in this packet"
  Receiver: "I received bytes 1000-1999. Send me 2000 next."
  
If packet is lost:
  Receiver stops acknowledging (or sends duplicate ACKs)
  Sender's timeout fires → retransmits
  
Sliding window:
  Sender can have up to window_size bytes in-flight (sent but not yet acknowledged)
  Receiver advertises window size (based on buffer space)
  This is flow control — prevents sender from overwhelming receiver
```

### TCP Congestion Control

```
Slow start:
  Connection starts with small congestion window
  Window doubles each RTT (exponential growth)
  When packet loss detected: window halved, linear growth resumes
  
This is why a new connection starts slow: TCP is being conservative.
Large file transfers benefit from connection reuse — window has grown.
HTTP/1.1 connection keep-alive exists for this reason.
```

### TIME_WAIT State

```
After closing: connection enters TIME_WAIT for 2×MSL (Maximum Segment Lifetime, ~2 min)
Purpose: ensures any late packets from the old connection are discarded
Problem: too many connections in TIME_WAIT can exhaust ports
Fix: SO_REUSEADDR socket option, or increase ephemeral port range
```

---

## 3. HTTP/1.1, HTTP/2, and HTTP/3

### HTTP/1.1 Problems

```
Head-of-line blocking:
  One TCP connection can serve one request at a time
  Response B can't be sent until response A is done
  
  Workaround: browsers open 6 connections per host
  But: 6 TCP handshakes, 6 congestion windows, more overhead
  
Request headers not compressed:
  Each request sends full headers (User-Agent, Cookie, etc.) — 800 bytes average
  For 100 requests/page: 80KB of headers overhead
```

### HTTP/2 — Multiplexing Over One TCP Connection

```
Key features:
  1. Multiplexing: multiple requests over ONE TCP connection simultaneously
     Request A and B are sent, responses interleave as they complete
     No head-of-line blocking at HTTP level
  
  2. Header compression (HPACK):
     Headers are compressed using static and dynamic tables
     Common headers sent as small integers, not full strings
  
  3. Server push:
     Server can proactively send resources (CSS, JS) before browser requests them
     
  4. Stream prioritization:
     Client signals priority of each stream; server can reorder
  
Binary protocol (not human-readable text):
  More efficient parsing, but harder to debug (need Wireshark/special tools)

HTTP/2 head-of-line blocking:
  Multiplexing fixes HTTP-level blocking but NOT TCP-level blocking
  If one TCP packet is lost, ALL streams in that connection stall
  This is why HTTP/3 was created
```

### HTTP/3 — QUIC Protocol

```
HTTP/3 runs over QUIC instead of TCP:
  QUIC is built on UDP but implements:
    - Reliable delivery (like TCP but per-stream)
    - Encryption (TLS 1.3 built-in)
    - 0-RTT connection setup (no separate TLS handshake)
    - Independent streams: packet loss on stream A doesn't affect stream B

Connection setup comparison:
  HTTP/1.1: TCP handshake (1 RTT) + TLS handshake (1-2 RTT) = 2-3 RTT
  HTTP/2:   TCP handshake (1 RTT) + TLS handshake (1-2 RTT) = 2-3 RTT
  HTTP/3:   QUIC+TLS 1.3 combined (1 RTT, or 0-RTT for resumption)
  
Connection migration:
  QUIC connections identified by Connection ID (not IP:port tuple)
  Phone switches from WiFi to 4G → connection continues without interruption
```

---

## 4. TLS & HTTPS — What Actually Happens

### TLS 1.3 Handshake (Simplified)

```
Client                                    Server

ClientHello ──────────────────────────▶
  (supported cipher suites, key_share)
                 ◀────────────────────── ServerHello
                                          (chosen cipher, key_share)
                                          Certificate
                                          CertificateVerify
                                          Finished
Client verifies certificate
Finished ──────────────────────────────▶
              
Application data flows (TLS 1.3: 1-RTT total)
```

### Certificate Verification

```
Certificate contains:
  - Server's public key
  - Domain name (CN/SAN)
  - Validity period
  - Digital signature from a Certificate Authority (CA)

Client's browser has a list of trusted CAs (built into OS/browser)

Verification:
  1. Is the certificate signed by a trusted CA? (chain of trust)
  2. Is the domain name correct? (matches the URL I'm accessing)
  3. Is the certificate still valid? (not expired, not revoked)

Certificate Transparency (CT):
  CAs must log all issued certificates to public CT logs
  Browsers check CT logs to detect misissued certificates
```

### What TLS Protects

```
Confidentiality: data encrypted with symmetric key (AES-256-GCM)
                 Only client and server can decrypt
                 
Integrity: HMAC/AEAD ensures data wasn't modified in transit
           Any tampering invalidates the message

Authentication: server proves it owns the private key corresponding to its certificate
                Prevents man-in-the-middle attacks
                
Note: TLS does NOT guarantee the server is legitimate — only that it owns the certificate
      A phishing site can have a valid TLS certificate for "paypa1.com"
```

### HSTS (HTTP Strict Transport Security)

```
Response header: Strict-Transport-Security: max-age=31536000; includeSubDomains; preload

Tells browsers: always use HTTPS for this domain, never HTTP, for the next N seconds

Prevents SSL stripping attacks (where attacker downgrades HTTPS to HTTP)
HSTS preload: domain is hardcoded into browsers, HTTPS enforced even on first visit
```

---

## 5. DNS Resolution

```
User types "api.stripe.com" in browser:

1. Browser checks its DNS cache (TTL-based)
2. OS checks its DNS cache + /etc/hosts
3. OS sends query to configured recursive resolver (usually ISP's or 8.8.8.8)
4. Recursive resolver checks its cache
5. If cache miss:
   a. Resolver asks root nameserver: "Who handles .com?"
   b. Root returns: "com nameservers are a.gtld-servers.net, ..."
   c. Resolver asks .com nameserver: "Who handles stripe.com?"
   d. Stripe.com NS returns IP address for api.stripe.com
6. Resolver caches the result (per TTL), returns to browser
7. Browser connects to that IP address

Total time (cold): 50-100ms (multiple network round trips)
Warm (cached): <1ms

DNS record types:
  A:     hostname → IPv4 address
  AAAA:  hostname → IPv6 address
  CNAME: hostname → another hostname (alias)
  MX:    mail server for domain
  TXT:   arbitrary text (SPF, DKIM, domain verification)
  NS:    authoritative nameservers for domain
  SOA:   start of authority (primary NS, admin email, TTL info)
```

---

## 6. CDN Internals

```
CDN (Content Delivery Network): geographically distributed servers that cache content close to users

How it works:
  1. Your DNS record points to CDN (CNAME to CDN hostname)
  2. User requests asset → CDN routes to nearest edge node (via anycast or DNS geolocation)
  3. Edge node checks cache: HIT → return cached content
  4. Cache MISS → edge fetches from origin, caches it, serves user
  5. Future users in same region: HIT → served from edge (no trip to origin)

CDN caching:
  Static assets: long TTL (JS/CSS with hashed filenames: 1 year)
  Dynamic content: short TTL or no cache (user-specific data)
  
  Cache-Control: public, max-age=31536000, immutable
  — tells CDN and browser to cache for 1 year, content will never change

CDN and TLS:
  CDN terminates TLS at the edge (decrypts client's request)
  Re-encrypts when fetching from origin (or uses HTTP internally on secure private network)
  This is "TLS offloading" — reduces CPU load on origin servers
```

---

## 7. Interview Questions & Model Answers

**Q: Explain what happens when you type "https://google.com" in a browser and press Enter.**

"The browser checks its DNS cache for google.com. If not cached, the OS queries a recursive resolver, which walks the DNS tree: root → .com nameservers → google.com nameservers → gets an A record with Google's IP. The browser then opens a TCP connection (3-way handshake) to that IP on port 443. After TCP is established, the TLS handshake begins — client sends supported cipher suites, server responds with its certificate. The browser verifies the certificate (trusted CA, correct domain, valid dates). A shared symmetric key is derived via Diffie-Hellman key exchange. Now encrypted HTTP/2 (or HTTP/3) begins. The browser sends a GET request, Google responds with HTML. The browser parses HTML, finds CSS/JS/images, fetches those (often from CDN), renders the page."

**Q: What is HTTP/2 multiplexing and what problem does it solve?**

"HTTP/1.1 can only serve one request per TCP connection at a time. Loading a page with 100 resources requires browsers to open 6 connections per host (browser limit) and queue remaining requests. HTTP/2 introduces streams — multiple request-response pairs can be in-flight simultaneously over a single TCP connection. Response B can return before response A finishes if it's ready sooner. This eliminates request queuing and reduces the need for many parallel connections. It also adds header compression (HPACK) and server push. The remaining limitation: a single lost TCP packet stalls all streams on that connection. HTTP/3 (QUIC) solves this by having independent streams at the transport layer."

---

## Summary

- **TCP:** connection-oriented, reliable, ordered. 3-way handshake. Flow control (window size). Congestion control (slow start).
- **UDP:** no connection, no guarantees, low overhead. For: DNS, video, gaming.
- **HTTP/1.1:** one request per connection. 6 parallel connections per host. Head-of-line blocking.
- **HTTP/2:** multiplexed streams over one TCP connection. Header compression (HPACK). Binary protocol.
- **HTTP/3:** QUIC (UDP-based). Independent streams — one lost packet doesn't stall others. 0-RTT resumption.
- **TLS 1.3:** 1-RTT handshake. Certificates prove server identity. AES-256-GCM for confidentiality + integrity.
- **DNS:** hierarchical resolution. Root → TLD → authoritative NS. Cache everything with TTL.
- **CDN:** edge nodes cache content near users. TLS termination at edge. Cache-Control headers drive caching behavior.
