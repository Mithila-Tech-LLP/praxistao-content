# Chapter 83: Common Network Attacks — Sniffing, MITM, Spoofing, SYN Floods, and DDoS

> **"Every attack in this chapter is not a bug. It's a protocol working exactly as designed — for a world that assumed cooperation instead of hostility. Understanding these attacks means understanding, precisely, which earlier chapter's design assumption got violated."**

---

## Table of Contents

1. [The Big Question, Restated From Chapter 77](#1-the-big-question-restated-from-chapter-77)
2. [A Map: Attack to Chapter to Defense](#2-a-map-attack-to-chapter-to-defense)
3. [Passive Packet Sniffing](#3-passive-packet-sniffing)
4. [Man-in-the-Middle Attacks, in General](#4-man-in-the-middle-attacks-in-general)
5. [ARP Spoofing — Abusing Chapter 53's Trust](#5-arp-spoofing--abusing-chapter-53s-trust)
6. [DNS Cache Poisoning — Abusing Chapter 68's Caching](#6-dns-cache-poisoning--abusing-chapter-68s-caching)
7. [IP Spoofing](#7-ip-spoofing)
8. [SYN Flood — Abusing Chapter 59's Handshake](#8-syn-flood--abusing-chapter-59s-handshake)
9. [DDoS at Scale: Volumetric vs. Application-Layer](#9-ddos-at-scale-volumetric-vs-application-layer)
10. [Session Hijacking](#10-session-hijacking)
11. [A Real, Safe Hands-On Experiment](#11-a-real-safe-hands-on-experiment)
12. [Common Misconceptions](#12-common-misconceptions)
13. [What's Simplified Here, and a Note on Framing](#13-whats-simplified-here-and-a-note-on-framing)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#16-summary)

---

## 1. The Big Question, Restated From Chapter 77

Chapter 77 introduced threat modeling: asking, for any protocol, "who might attack this, with what capability, and what would they gain?" Chapters 78–82 then spent five chapters building the defenses — encryption, authentication, PKI, and finally TLS as the protocol that assembles all of them.

This chapter asks the question from the other direction: **for every protocol taught in this course before Volume 12, what happens if you take away the assumption that everyone on the network is cooperating?**

Every single attack below is defensible in exactly this framing: find the protocol chapter, find the assumption it quietly made about a well-behaved network, and show what happens the moment that assumption is false. This is not a catalog to memorize — it's a lens. Once you can look at any protocol and ask "what did this assume, and who benefits from that assumption being false?", you can reason about attacks this chapter never explicitly lists.

**A framing note before we go further, stated once and meant seriously:** everything from here on is explained for *recognition and defense* — to help you identify an attack in progress, understand why a mitigation works, and reason about your own systems' exposure. Nothing here is a how-to for causing harm, and several operationally sensitive details (exact tool commands, precise payload construction) are deliberately omitted in favor of the underlying mechanism, which is the part that actually matters for defense.

---

## 2. A Map: Attack to Chapter to Defense

Keep this table as the chapter's spine — every section below expands exactly one row.

| Attack | Protocol assumption abused | Chapter it abuses | Primary defense |
|---|---|---|---|
| Passive packet sniffing | Data travels in plaintext | Chapter 77 (threat model baseline) | TLS (Chapter 82) |
| Man-in-the-middle (general) | Endpoints can't verify who they're really talking to | Chapters 78–81 | Certificate/identity verification (Chapter 81) |
| ARP spoofing | Any device's ARP reply is trusted, unauthenticated | Chapter 53 (ARP) | Dynamic ARP Inspection, static entries, port security |
| DNS cache poisoning | A resolver caches the first plausible-looking answer | Chapter 68 (DNS caching) | DNSSEC (Chapter 69), randomized query IDs/ports |
| IP spoofing | The source address field is never verified | Chapter 36 (IPv4 addressing) | Ingress/egress filtering (BCP 38) |
| SYN flood | Server allocates state on SYN, before the handshake completes | Chapter 59 (three-way handshake) | SYN cookies |
| Volumetric / application DDoS | Legitimate and malicious traffic can look identical at scale | Chapters 58, 71 (UDP, HTTP) | Scrubbing, anycast, rate limiting, CDN absorption |
| Session hijacking | A session is only as strong as the identifier that proves it | Chapter 72 (cookies/sessions) | Secure/HttpOnly cookies, TLS, session rotation |

---

## 3. Passive Packet Sniffing

**The assumption abused:** before TLS is in the picture, most protocols this course has taught (HTTP in Chapter 71, DNS queries in Chapter 66, plain FTP, plain SMTP) send their payload as plaintext bytes on the wire. The assumption, implicitly, is that nobody unauthorized is watching the wire at all.

**The mechanism.** Any device with access to the physical medium a frame travels over can, in principle, read every byte of every frame that reaches its network interface — this is precisely what a network interface card does when placed in **promiscuous mode**: instead of discarding frames not addressed to it (the normal behavior described in Chapter 28), it hands every frame it sees up to software. On a hub-based network (Chapter 30), this meant *any* device could see *all* traffic, because a hub is a dumb repeater with no forwarding intelligence at all. Modern switches limit this significantly — a switch (Chapter 31) only forwards a frame out the port where the destination MAC address lives, so a machine plugged into a switch normally sees only traffic addressed to it or broadcast to everyone. But switches don't eliminate the problem, only shrink it: an attacker on a shared Wi-Fi network, a compromised router along the path, or an attacker who can force traffic to flow through a device they control (Section 5's ARP spoofing is exactly how) can still capture everything that crosses that point.

**What it reveals, if the payload is plaintext:** login credentials, session cookies, private messages, the content of any unencrypted request or response — all readable as ordinary text.

```
tcp stream (plain HTTP, port 80), captured with tcpdump -A:

GET /login HTTP/1.1
Host: example.com
Cookie: session=a91f...

POST /login HTTP/1.1
Host: example.com
Content-Type: application/x-www-form-urlencoded

username=alice&password=hunter2
```

Notice: nothing above is obfuscated in any way — it's the literal bytes of the HTTP request, exactly as Chapter 71 described them, simply visible to anyone with a copy of the frames.

**The defense.** TLS (Chapter 82) is the complete answer to this specific attack: once payload bytes are encrypted with an AEAD cipher (Chapter 78) and a session key nobody outside the two endpoints holds, a passive observer captures only ciphertext — mathematically useless without the key. This is precisely why "is this connection using HTTPS" is the single most important habit for defending against sniffing, and why plaintext protocols (unencrypted FTP, Telnet, plain HTTP) are considered obsolete for anything sensitive.

---

## 4. Man-in-the-Middle Attacks, in General

**The assumption abused:** that the party you're establishing a secure channel with is actually who they claim to be, and that nobody sits between you and them able to intercept, read, or alter traffic in both directions.

**The mechanism.** A man-in-the-middle (MITM) attack is a category, not a single technique: the attacker positions themselves on the path between two endpoints — physically (a rogue Wi-Fi access point), logically (Section 5's ARP spoofing forces local traffic through the attacker), or at a higher layer (Section 6's DNS poisoning redirects a victim to an attacker-controlled server entirely) — and then either passively observes or actively rewrites traffic passing through.

```
Normal path:
  Client ───────────────────────────▶ Server

MITM path:
  Client ─────▶ Attacker ─────▶ Server
         (reads/modifies both directions, relays traffic so neither side notices)
```

**Why TLS mostly defeats a generic MITM.** If a client correctly performs certificate verification (Chapter 81), an attacker sitting in the middle cannot present the real server's certificate (they don't have its private key) and cannot forge a certificate a trusted CA would sign for a domain they don't control. The attacker's only remaining option is to present their own, different certificate — which a correctly implemented client will reject, precisely because that's what the whole PKI trust chain was built to catch. The practical failure mode is almost never "TLS was broken" — it's a client accepting a certificate it shouldn't (a user clicking through a browser warning, a mobile app with certificate validation disabled for debugging and shipped that way, or a corporate network that has installed its own trusted root CA to intentionally intercept and re-encrypt traffic, which is MITM by design, done with the endpoint's consent).

**Where MITM still bites even with TLS everywhere:** before a connection is established (an attacker on an open Wi-Fi network redirecting a plaintext HTTP request to a spoofed login page, exploiting the fact that most users type `example.com` without `https://`, an issue HSTS from Chapter 82's neighbor material addresses by forcing HTTPS after the first visit), and at any point where the "trusted" identity itself is compromised (a stolen or mis-issued certificate, or ARP/DNS attacks that redirect traffic to an attacker who has *also* somehow obtained a valid certificate, which is rare but not impossible given past CA compromises).

---

## 5. ARP Spoofing — Abusing Chapter 53's Trust

**The assumption abused, precisely.** Chapter 53 described ARP's entire design: a device broadcasts "who has IP address X? Tell me your MAC," and whichever device replies "that's me, here's my MAC" is simply believed, cached, and used for all future traffic to that IP — with **no authentication of any kind**. ARP was designed in an era (1982, RFC 826) where every device on a LAN was assumed to be a cooperating, known machine. Nothing checks that the reply came from the device that actually owns that IP address.

**The mechanism.** An attacker on the same local network sends unsolicited ARP replies (or "gratuitous ARP" announcements, which Chapter 53 noted exist for legitimate reasons like failover) claiming that the attacker's own MAC address corresponds to, say, the default gateway's IP address. Victims on the LAN update their ARP cache to point the gateway's IP at the attacker's MAC. From that point on, every packet a victim intends for the gateway — which usually means every packet leaving the LAN at all — is instead delivered to the attacker's machine.

```
Before spoofing:
  Victim's ARP cache:  192.168.1.1 (gateway)  ->  AA:AA:AA:AA:AA:AA (real gateway)

Attacker broadcasts forged ARP reply:
  "192.168.1.1 is at BB:BB:BB:BB:BB:BB"  (attacker's MAC)

After spoofing:
  Victim's ARP cache:  192.168.1.1 (gateway)  ->  BB:BB:BB:BB:BB:BB (attacker!)

Result: every frame the victim sends to the gateway now physically
arrives at the attacker's NIC first. The attacker can read it,
optionally forward it on to the real gateway (invisible MITM), or drop it.
```

This is the classic mechanism behind local-network MITM — it's specifically what lets an attacker on a coffee-shop Wi-Fi network intercept traffic from other devices on the same network, and it's why Section 4's MITM category so often starts here in practice.

**Defenses.** Static ARP entries for critical infrastructure (gateways, servers) that ignore unsolicited replies entirely; **Dynamic ARP Inspection (DAI)**, a switch feature that cross-checks incoming ARP replies against a known-good binding table (often built from DHCP snooping records) and drops replies that don't match; and, at the application layer, TLS (Chapter 82) means that even a successful ARP-spoof-based MITM only yields the attacker unreadable ciphertext for any properly configured HTTPS connection — ARP spoofing gets you *position* on the path, not automatically *content*.

---

## 6. DNS Cache Poisoning — Abusing Chapter 68's Caching

**The assumption abused, precisely.** Chapter 68 explained why DNS caching exists: without it, every single web request would require a full walk of the DNS hierarchy (Chapter 67), which would be far too slow. A recursive resolver caches an answer for the duration of its TTL and serves that cached answer to every subsequent client that asks, trusting that the answer it received and cached was genuine. Classic DNS, as taught up through Chapter 68, has essentially no way to verify that an answer actually came from the authoritative server it claims to be from — it's UDP-based (Chapter 58), and the only "proof" of authenticity is a 16-bit transaction ID matching the query.

**The mechanism.** An attacker races to send a resolver a forged DNS response — claiming to be the reply to a query the resolver just sent to a legitimate authoritative server — before the real answer arrives, guessing (or observing) the matching 16-bit query ID and matching source port. If the forged answer arrives first and the ID matches, the resolver caches it as if it were genuine, and **every client using that resolver for the remainder of the TTL is silently redirected** to whatever IP address the attacker specified — potentially a server designed to impersonate the real site. The 2008 Kaminsky attack made this practical at scale by showing an attacker could repeatedly query for random, non-existent subdomains of a target domain, giving unlimited attempts to win the race to poison the *parent* domain's nameserver record in the cache, rather than needing to guess the ID for one single, rarely-repeated query.

```
Normal resolution (Chapter 68):
  Resolver ──query (ID=4471)──▶ Authoritative NS
  Resolver ◀─answer (ID=4471)── Authoritative NS   [genuine, cached]

Poisoning attempt:
  Resolver ──query (ID=4471)──▶ Authoritative NS
  Resolver ◀─forged answer (ID guessed=4471)── Attacker   [race to arrive first]
  If attacker wins the race: resolver caches attacker's answer as if genuine
  Every client asking this resolver now gets the attacker's IP, until TTL expires
```

**Defenses.** Source port randomization (making the ID+port pair the attacker must guess far larger than 16 bits alone), 0x20 encoding (randomizing the case of letters in the query name as an extra unpredictable value the forged reply must match), and — the real structural fix — **DNSSEC** (Chapter 69), which adds cryptographic signatures to DNS records so a resolver can verify an answer actually came from the zone's real key holder, closing the gap the same way Chapter 81's PKI closes the "whose public key is this" gap for TLS.

---

## 7. IP Spoofing

**The assumption abused.** Chapter 36 introduced the IPv4 header's source address field as simply a field the sender fills in — nothing in the base IP protocol cryptographically ties a packet's claimed source address to the machine that actually sent it. Routers forward packets based on the *destination* address (Chapter 45's longest-prefix-match); almost nothing along a typical path checks whether the *source* address is plausible for the network the packet actually arrived from.

**The mechanism.** An attacker constructs packets with a forged source IP address — anything from a real, uninvolved third party's address to an address inside a network the attacker has no business sending from. Because return traffic goes to whatever address is in the source field, not to the attacker, "blind" IP spoofing (used, for instance, in some SYN flood variants in Section 8, and in DDoS amplification in Section 9) doesn't require the attacker to ever see the response — it's used precisely when the goal is to hide the true source, or to direct response traffic at an unwitting third party (the victim of a reflection attack) instead of back at the attacker.

**Defenses.** **BCP 38 (RFC 2827) — ingress/egress filtering**: an ISP or network operator configures routers at the network's edge to drop any outbound packet whose source address doesn't belong to that network's own assigned address block, and any inbound packet claiming a source address that belongs internally. Widely deployed at the network edge, this closes off spoofing at the source — but because it requires every network operator worldwide to configure it correctly, and many still don't, IP spoofing from poorly configured networks remains a real, ongoing enabler of the amplification attacks in Section 9.

---

## 8. SYN Flood — Abusing Chapter 59's Handshake

**The assumption abused, precisely.** Chapter 59 explained the three-way handshake: a client sends SYN, the server responds with SYN-ACK **and allocates connection state for that half-open connection** (a slot in its backlog queue, holding the initial sequence number and connection metadata), then waits for the final ACK to complete the connection. The handshake's design implicitly assumes clients that send a SYN will, in reasonable time, either complete the handshake or the half-open state will time out and get cleaned up without much cost. **The server commits real memory and a queue slot before it has any proof the client's claimed source address is genuine or that the client intends to finish.**

**The mechanism.** An attacker sends a flood of SYN packets — often with spoofed source addresses (Section 7), so the resulting SYN-ACKs go to unrelated, uninvolved third parties who simply drop them, never producing the final ACK — far faster than the server's backlog queue can be cleared by timeouts. Each SYN consumes a backlog slot; if the flood arrives faster than half-open connections expire, the backlog fills entirely, and **the server can no longer accept any new legitimate connection**, because it has no free slot to record one, even though the server's CPU and bandwidth may otherwise be nearly idle.

```
Normal handshake (Chapter 59):
  Client ──SYN──▶ Server (allocates backlog slot)
  Client ◀─SYN-ACK─ Server
  Client ──ACK──▶ Server (slot converted to established connection)

SYN flood:
  Attacker ──SYN (spoofed src A)──▶ Server (allocates slot 1)
  Attacker ──SYN (spoofed src B)──▶ Server (allocates slot 2)
  Attacker ──SYN (spoofed src C)──▶ Server (allocates slot 3)
  ... thousands more, arriving faster than half-open slots time out ...
  Server ◀─SYN-ACKs sent to A, B, C── (spoofed addresses never respond)
  Backlog queue full → legitimate client's SYN is dropped or refused
```

**Defenses.** **SYN cookies** are the elegant, standard mitigation: instead of allocating state on receipt of a SYN, the server encodes the connection's necessary state (a hash of source/destination address and port, a timestamp, and the negotiated options) directly into the sequence number it sends back in the SYN-ACK, and allocates **no backlog slot at all** until the final ACK arrives with a sequence number that decodes back into valid, consistent state. This means an attacker's spoofed SYNs, which never receive or return a valid ACK, cost the server essentially nothing beyond generating and sending the SYN-ACK itself — the expensive resource (backlog memory) is never committed to a connection that never completes. Rate limiting incoming SYNs per source, and reducing the backlog timeout so abandoned half-open connections are reclaimed faster, are complementary, coarser defenses.

---

## 9. DDoS at Scale: Volumetric vs. Application-Layer

A Distributed Denial-of-Service attack is, at its core, the SYN flood's underlying goal (exhaust some resource so legitimate users can't be served) generalized: instead of one attacker exhausting one specific resource (a backlog queue), many sources (often a botnet of compromised machines) simultaneously exhaust some resource on the target. It's useful to split DDoS into two categories that abuse completely different layers:

**Volumetric attacks — exhaust bandwidth itself.** The goal is simply to send more traffic than the target's network link (or an upstream link) can carry, so legitimate packets are dropped purely due to congestion, regardless of what those packets contain. A key technique here is **amplification/reflection**, which combines Section 7's IP spoofing with a protocol that returns a much larger response than the request that triggered it:

```
Reflection + amplification attack (defensive framing — mechanism, not a recipe):

  Attacker sends a small, spoofed-source query to an open, misconfigured
  service (historically: open DNS resolvers, NTP servers, memcached
  instances) — the spoofed source address is the victim's IP.

  Attacker's small query ──▶  Reflector service  ──▶  Large response sent to "victim"
       (e.g. 60 bytes)          (open DNS/NTP/etc.)      (e.g. 6,000+ bytes)

  Attacker never sends the large volume directly — the reflector's
  own bandwidth, multiplied by the response/request size ratio,
  is what actually floods the victim.
```

Amplification factors of 50-100x (and historically much higher for some misused protocols) mean a modest amount of attacker-controlled bandwidth can generate a genuinely massive flood at the victim — this is why the defenses in Section 7 (BCP 38 ingress filtering, which prevents the spoofed queries from ever leaving their network of origin in the first place) and properly configuring internet-facing UDP services to not respond to arbitrary requests from arbitrary sources are both taken seriously by network operators, not just by the eventual victims.

**Application-layer attacks — exhaust a specific server-side resource, at much lower bandwidth.** Instead of overwhelming the network link, these attacks send traffic that looks like ordinary, valid application traffic (Chapter 71's HTTP requests) but is specifically shaped to be expensive for the server to process relative to how cheap it is to send:

- **HTTP flood**: sending a large volume of ordinary-looking GET/POST requests, often targeting a specific expensive endpoint (a search query, a report generation page) rather than the homepage, so each request consumes disproportionate CPU or database time.
- **Slowloris-style attacks**: opening many connections and sending HTTP request headers *extremely slowly*, one byte at a time, keeping each connection open just long enough to avoid a timeout — since a naive web server allocates a worker thread/process per connection (Chapter 57's socket model), a modest number of slow, incomplete connections can exhaust the server's entire worker pool without ever sending much data or consuming much bandwidth at all.

| | Volumetric | Application-layer |
|---|---|---|
| Layer targeted | Network/transport (bandwidth, packet processing) | Application (Chapter 71's HTTP semantics) |
| Traffic volume needed | Very high (Gbps-Tbps) | Can be very low (a few hundred connections) |
| Looks like legitimate traffic? | No — often malformed or reflected | Often yes — valid HTTP requests |
| Typical defense | Upstream scrubbing centers, anycast (spreading load across many points of presence, Chapter 96), ISP-level filtering | Rate limiting, WAF pattern/behavior rules (Chapter 84), CAPTCHA challenges, connection timeouts |

**Defenses in general.** Because the attacking traffic often can't be perfectly distinguished from legitimate traffic packet-by-packet, real-world DDoS defense relies heavily on **scale and distribution**: scrubbing centers that can absorb and filter enormous volumes before forwarding only clean traffic on; Anycast (previewed in Chapter 69 for DNS, generalized in Chapter 96) which spreads a single logical destination across many physical locations so an attack's traffic is naturally divided among them instead of concentrated on one link; and rate limiting plus behavioral analysis (unusually high request rates from a single source, requests that don't follow a normal browsing pattern) at the CDN or WAF layer, which Chapter 84 covers in depth.

---

## 10. Session Hijacking

**The assumption abused.** Chapter 72 introduced cookies and sessions as the mechanism that gives stateless HTTP the illusion of a continuous, logged-in relationship — but a session, mechanically, is only as strong as the identifier (the session cookie) that proves "this request belongs to the same logged-in user as the last one." If that identifier can be observed, guessed, or stolen, an attacker can present it and be treated as the legitimate user, with no further authentication required.

**The mechanism, briefly, since it composes attacks already covered:** a session cookie sent over plaintext HTTP can simply be read via Section 3's packet sniffing; a cookie without the `Secure` flag might be sent over an accidental plaintext connection even on a site that mostly uses HTTPS; a cookie without `HttpOnly` is readable by any JavaScript running on the page, making it a target for cross-site scripting; and predictable session identifiers (sequential numbers, or values derived from guessable data) can sometimes be brute-forced or predicted outright without ever intercepting traffic. In every case, the attack yields the attacker a live, authenticated session — not a password, but something often just as useful.

**Defenses.** TLS everywhere (removing Section 3's plaintext exposure entirely), the `Secure` and `HttpOnly` cookie attributes (Chapter 72), cryptographically random, high-entropy session identifiers, and rotating session identifiers after login (so a session ID observed before authentication can't be reused after).

---

## 11. A Real, Safe Hands-On Experiment

Every step below is performed only against your own machine/network, for defensive learning:

1. On your own laptop, run `curl -v http://neverssl.com` (a site intentionally kept plaintext for testing) and, in a second terminal, run `sudo tcpdump -A -i any port 80 host neverssl.com` while the request runs — you'll see the entire HTTP request and response as plaintext, exactly as Section 3 describes.
2. Repeat against any HTTPS site (`curl -v https://example.com`) with the same `tcpdump` filter adjusted to port 443 — confirm you see only encrypted, unreadable bytes, demonstrating TLS's defense against exactly the same capture technique.
3. Run `arp -a` on your own machine and note the MAC address currently associated with your gateway's IP — this is the exact cache Section 5 described being poisoned; understanding what a legitimate entry looks like is the first step toward noticing an illegitimate one.
4. Look up your own ISP's or a public DNS resolver's behavior with `dig +dnssec example.com` and check whether a signed response (an `RRSIG` record) comes back — this shows DNSSEC (Section 6's defense) in action, or its absence, for a domain of your choosing.

---

## 12. Common Misconceptions

- **"A VPN makes you immune to all of this."** A VPN (Chapter 85) protects the segment of the path between you and the VPN endpoint — it does nothing about a compromised server you're legitimately connecting to, or about attacks (like DNS poisoning at a resolver you still use) that happen entirely on the other side of the tunnel.
- **"DDoS and DoS are the same thing, just bigger."** The "distributed" part matters mechanically, not just quantitatively — spreading an attack across many sources defeats simple defenses like "block this one IP," which is precisely why real-world DDoS mitigation depends on distributed absorption (scrubbing, anycast) rather than simple blocklisting.
- **"SYN cookies mean SYN floods are a solved problem."** They neutralize the *backlog exhaustion* mechanism specifically; a sufficiently large SYN flood can still exhaust raw bandwidth or CPU generating SYN-ACKs, which is why SYN cookies are one layer of defense, not a complete answer to every flood.
- **"ARP spoofing only matters on 'hacker' networks."** It's a routine risk on any shared, untrusted local network — coffee shop Wi-Fi, conference networks, and any LAN where you don't control or trust every other connected device.
- **"If a resolver supports DNSSEC, I'm automatically protected."** DNSSEC only protects domains that have signed their own records, and only if your resolver actually validates signatures rather than just passing them through — both ends of that chain have to hold.

---

## 13. What's Simplified Here, and a Note on Framing

This chapter deliberately covers the *mechanism* of each attack at the level needed to understand why a specific defense works — not exact tool syntax, packet construction detail, or step-by-step exploitation instructions, which is intentional and consistent throughout the chapter, not an oversight. Real attacks are often combined (an ARP spoof used to enable a session hijack; IP spoofing used inside a reflection DDoS) and real-world detection relies on tooling (intrusion detection systems, anomaly-based monitoring) this chapter didn't have room to cover. The goal here was recognition: given a symptom (a suddenly unreachable server, an unexpectedly redirected domain, a flood of half-open connections), you should now be able to reason backward to a plausible cause and the chapter it violates.

---

## 14. Interview Questions & Model Answers

**Beginner: "What makes a SYN flood effective, mechanically?"**

*Model answer:* "TCP's three-way handshake requires the server to allocate connection state — a backlog queue slot — as soon as it receives a SYN, before it has any confirmation the client is real or will ever complete the handshake. An attacker sends a flood of SYNs, often with spoofed source addresses so the resulting SYN-ACKs go nowhere useful, faster than half-open connections can time out and be cleared. Once the backlog queue is full, the server can't accept any new legitimate connection, even though it might have plenty of free CPU and bandwidth otherwise."

**Intermediate: "How do SYN cookies fix this without breaking the handshake?"**

*Model answer:* "Instead of storing connection state in a backlog queue on receipt of a SYN, the server encodes exactly what it would have stored — a hash of the source/destination address and port plus a timestamp — directly into the sequence number field of its own SYN-ACK reply, and allocates no memory for the connection at all. When (and only when) a legitimate client responds with the final ACK, the server recomputes the expected value from the ACK's sequence number and, if it matches, reconstructs the connection state on the spot. Spoofed SYNs that never receive a real ACK back cost the server nothing beyond generating one reply packet — there's no queue for them to exhaust."

**Advanced: "Why does DNS cache poisoning depend on winning a race, and what makes the Kaminsky variant more dangerous than a naive one?"**

*Model answer:* "A resolver only accepts a DNS reply as valid for a pending query if the reply's 16-bit transaction ID (and, with source port randomization, the source port too) matches what the resolver sent. An attacker has to guess that value and deliver a forged reply before the genuine authoritative server's real answer arrives — a race with fairly long odds against for a single attempt. The Kaminsky attack made this far more dangerous by not attacking a single, rarely-repeated record directly. Instead, it repeatedly queries the resolver for random, non-existent subdomains of the target domain — each one forces a brand-new race the attacker gets another shot at — while the forged replies don't just answer that fake subdomain, they include a forged NS record for the entire parent domain in the additional-records section. Win even one of those many, cheap races, and the attacker poisons the authoritative nameserver record for the whole domain, redirecting all future lookups for it, not just one subdomain. DNSSEC closes this by making the resolver cryptographically verify the reply came from the real zone owner, regardless of whether the transaction ID guess was correct."

---

## 15. Exercises

### Easy

1. For each of the eight attacks in Section 2's table, name the one earlier-chapter protocol assumption it violates, without looking at the table.
2. Explain, in plain language, why a switch (Chapter 31) reduces but does not eliminate the risk of packet sniffing on a LAN.
3. Why does a SYN flood attacker often use spoofed source addresses rather than their own real address?

### Medium

4. Explain why TLS defeats a generic passive-sniffing attacker but does not, by itself, defeat an ARP-spoofing-based MITM attacker who has already positioned themselves on the path — and what specifically TLS still protects even in that scenario.
5. A company's DNS resolver has DNSSEC validation enabled, but the target domain being queried has not signed its own DNS records. Explain whether this resolver is still vulnerable to cache poisoning of that particular domain, and why.
6. Explain the difference between a volumetric DDoS attack and an application-layer DDoS attack in terms of what resource is actually being exhausted, and why a defense effective against one (like adding more bandwidth) might do little against the other.

### Hard

7. Design (in prose, not code) a monitoring rule that could help a network operator detect an ARP spoofing attempt in progress on their LAN, using only information already available from routine ARP traffic — and explain one limitation of your rule (a case where it would miss a real attack, or flag a false positive).
8. A reflection/amplification DDoS attack depends on two separate failures happening at two different networks: an open, misconfigured reflector service, and a network that allows spoofed-source packets to leave it in the first place. Explain, referencing BCP 38 from Section 7, why fixing only the reflector-side misconfiguration (closing open resolvers) is not, by itself, a complete defense against this entire class of attack across the whole Internet.

---

## 16. Summary

| Term | Meaning |
|---|---|
| Passive sniffing | Reading plaintext traffic off a shared or intercepted link, defeated by encryption (TLS) |
| Man-in-the-middle (MITM) | An attacker positioned on the path between two endpoints, reading or altering traffic in both directions |
| ARP spoofing | Sending forged, unauthenticated ARP replies to redirect local traffic to the attacker |
| DNS cache poisoning | Racing a forged DNS reply into a resolver's cache before the genuine answer arrives |
| IP spoofing | Forging a packet's source address, enabling blind attacks and reflection/amplification |
| SYN flood | Flooding half-open TCP connections to exhaust a server's backlog queue before handshakes complete |
| SYN cookies | Encoding connection state into the SYN-ACK sequence number instead of allocating memory up front |
| Volumetric DDoS | Exhausting network bandwidth itself, often via reflection/amplification |
| Application-layer DDoS | Exhausting a specific server-side resource (CPU, worker threads, database time) with valid-looking requests |
| Session hijacking | Stealing or predicting a session identifier to impersonate an authenticated user |

Every attack in this chapter is a symptom of a network doing exactly what it was designed to do, aimed by someone the design never anticipated. Chapter 84 now turns to the devices built specifically to sit in the path and refuse to forward what shouldn't be forwarded — firewalls, and the web-aware layer built on top of them, the WAF.
