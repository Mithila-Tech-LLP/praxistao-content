# Chapter 84: Firewalls and Web Application Firewalls

> **"A firewall's hardest job was never blocking obviously bad traffic — attackers don't announce themselves. Its hardest job is looking at a packet that is, by every field in its header, indistinguishable from a legitimate one, and correctly deciding whether it belongs."**

---

## Table of Contents

1. [The Big Question](#1-the-big-question)
2. [The Naive Attempt: Block by Port or Address Alone](#2-the-naive-attempt-block-by-port-or-address-alone)
3. [Stateless Packet Filtering](#3-stateless-packet-filtering)
4. [Worked Example: The Same Packet, Evaluated Statelessly](#4-worked-example-the-same-packet-evaluated-statelessly)
5. [Stateful Firewalls — Remembering the Conversation](#5-stateful-firewalls--remembering-the-conversation)
6. [Worked Example: The Same Packet, Evaluated Statefully](#6-worked-example-the-same-packet-evaluated-statefully)
7. [Connection Tracking, NAT, and the State Table in Practice](#7-connection-tracking-nat-and-the-state-table-in-practice)
8. [The Firewall Family Tree](#8-the-firewall-family-tree)
9. [Web Application Firewalls — Application-Aware Defense](#9-web-application-firewalls--application-aware-defense)
10. [Worked Example: A WAF Evaluating an HTTP Request](#10-worked-example-a-waf-evaluating-an-http-request)
11. [WAF's Real Limitations](#11-wafs-real-limitations)
12. [Where Firewalls Actually Sit in Real Networks](#12-where-firewalls-actually-sit-in-real-networks)
13. [A Hands-On Experiment](#13-a-hands-on-experiment)
14. [Common Misconceptions](#14-common-misconceptions)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)

---

## 1. The Big Question

Chapter 83 walked through a tour of attacks that all rely on the same underlying fact: **a raw IP packet or TCP segment carries no field that says "I am malicious."** A SYN flood packet (Chapter 83, Section 8) has a perfectly valid TCP header. A spoofed packet (Chapter 83, Section 7) has a perfectly well-formed IP header — it's simply lying about its source. An HTTP request carrying a SQL injection attempt is, byte for byte, a syntactically valid HTTP request (Chapter 71).

This creates the central engineering problem of this chapter: **how do you build a device that sits on the path between a network and the outside world, and decides — packet by packet, or request by request — what to let through, when "good" and "bad" traffic can be structurally identical?**

The answer isn't one device — it's a spectrum, from very fast, very dumb rule-matching at the bottom, to slow, context-aware inspection at the top. This chapter walks that entire spectrum: stateless packet filters, stateful firewalls that remember what came before, and Web Application Firewalls that read all the way up into the meaning of an HTTP request. Each one trades speed for how much context it can bring to the decision — and each one closes gaps the one below it structurally cannot.

---

## 2. The Naive Attempt: Block by Port or Address Alone

The first instinct, for anyone encountering this problem fresh, is simple: **maintain a list of ports and addresses that are allowed, and drop everything else.** "Allow port 80 and 443 inbound (web traffic), allow port 22 from our office IP range only (SSH), block everything else."

This is a real and useful first layer — it's exactly what a home router's basic firewall does, and what a cloud security group (previewed for Chapter 97) starts from. But it immediately runs into a structural problem the moment you think about **return traffic**. A client outside your network that connects to your web server on port 443 needs the *response* to reach it — and that response, from the server's side, is a packet leaving your network, sourced from port 443, destined for whatever ephemeral port (Chapter 57) the client happened to pick, often a random high number.

If your rule is a static, direction-blind list ("allow anything on port 443, in or out"), you've just also allowed an attacker to send you traffic claiming to be *from* port 443, since the rule can't distinguish "this is a legitimate reply to a connection we started" from "this is a fresh, unsolicited packet that happens to use a common port number." **A rule with no memory of what came before has no way to make that distinction at all.** That gap is precisely what the rest of this chapter is about closing, one layer at a time.

---

## 3. Stateless Packet Filtering

**What it is.** A stateless packet filter — the original, and still the fastest, kind of firewall — evaluates **every single packet in complete isolation**, against a fixed list of rules (an Access Control List, or ACL), with zero memory of any packet that came before it. Each rule typically matches on some combination of: source/destination IP address, source/destination port, protocol (TCP/UDP/ICMP), and TCP flags.

```
Example stateless ACL (evaluated top to bottom, first match wins):

  1. ALLOW  tcp   dst port 443           (inbound HTTPS to our web server)
  2. ALLOW  tcp   dst port 80            (inbound HTTP)
  3. ALLOW  tcp   src port 443           (return traffic FROM our web server)
  4. ALLOW  tcp   src port 80            (return traffic FROM our web server)
  5. DENY   *     *                      (everything else)
```

**Why this is fast.** Because each packet is checked against the rule list independently, with no lookup into any per-connection state, this can be implemented directly in hardware (an ASIC in a router's forwarding path) or in the kernel's fast path, adding negligible latency even at very high packet rates. This is exactly why stateless filtering still exists today as the first, cheapest layer of defense in most real network designs — it's the traffic-light version of security: dumb, but nearly instantaneous.

**Why this is limited.** Rules 3 and 4 above exist specifically to solve the return-traffic problem from Section 2 — but notice what they actually say: **"allow any inbound TCP packet whose source port is 443, regardless of any other context."** That's not "allow the specific reply to the specific connection our web server initiated" — a stateless filter has no concept of "our specific connection" at all. It only knows "does this one packet's header match one of my static rules." This means an attacker can craft a packet with a forged source port of 443 and pass rules 3/4 even though no real connection ever solicited it — the filter genuinely cannot tell the difference, because telling the difference requires remembering something, and a stateless filter remembers nothing.

---

## 4. Worked Example: The Same Packet, Evaluated Statelessly

Consider one specific packet arriving at a firewall protecting an internal client machine:

```
Incoming packet:
  Source IP:    203.0.113.50   (unknown, external)
  Source port:  443
  Dest IP:      10.0.0.15      (internal client)
  Dest port:    51244          (an ephemeral port)
  Flags:        SYN, ACK
```

This packet claims to be a SYN-ACK — the second message of a TCP handshake (Chapter 59) — implying the internal client at `10.0.0.15` must have sent a SYN to `203.0.113.50:443` moments earlier and is now receiving the expected reply.

**A stateless filter's evaluation, using rule 3 from Section 3's list ("ALLOW tcp src port 443"):** the filter checks the packet's header fields against its rule list, finds source port 443 matches an ALLOW rule, and **lets it through — with no way to check whether the internal client ever actually sent a matching SYN in the first place.** If no such SYN was ever sent — if this packet is entirely unsolicited, perhaps part of a port scan or an attempt to smuggle traffic through a firewall by disguising it as "web reply traffic" — the stateless filter has no mechanism to notice. It answers only the question "does this packet's header look like something I generally allow," never "does this packet look like something I specifically expect right now."

---

## 5. Stateful Firewalls — Remembering the Conversation

**The real fix.** A stateful firewall maintains a **connection tracking table** — a running record of every connection currently considered "established" or "in progress," keyed by the same 4-tuple Chapter 57 introduced (source IP, source port, destination IP, destination port) plus protocol. Instead of asking "does this packet's header match a static rule," a stateful firewall asks the fundamentally different question: **"does this packet correspond to a connection I already know is actually in progress, or does it look like a genuine new connection attempt I'm configured to allow?"**

```
Connection tracking table (simplified):

  Src IP        Src port   Dst IP        Dst port   Protocol   State
  10.0.0.15     51244      203.0.113.50  443        TCP        SYN_SENT
```

The moment the internal client sends its SYN out through the firewall, the firewall creates exactly this entry, recording that `10.0.0.15:51244` is expecting a reply from `203.0.113.50:443`. When a packet later arrives claiming to be that reply, the firewall checks it not against a generic rule but against **this specific tracked entry** — matching the 4-tuple (reversed, since it's the return direction) and confirming the TCP flags are consistent with the connection's current tracked state (a SYN-ACK is a valid next step from `SYN_SENT`; an unexpected RST or a stray SYN-ACK for a connection not in the table is not).

This state machine mirrors TCP's own state machine directly: entries move through states that track the handshake (Chapter 59: `SYN_SENT` → `SYN_RECEIVED` → `ESTABLISHED`) and the connection close (Chapter 64: `FIN_WAIT` → `TIME_WAIT` → removed from the table). A stateful firewall is, in a very real sense, running a parallel copy of the TCP state machine for every connection crossing it, purely to answer the question "is this packet consistent with a connection that's actually happening?"

---

## 6. Worked Example: The Same Packet, Evaluated Statefully

Take the exact same packet from Section 4:

```
Incoming packet:
  Source IP:    203.0.113.50   Source port: 443
  Dest IP:      10.0.0.15      Dest port:   51244
  Flags:        SYN, ACK
```

**Case A — a real connection exists.** If the connection tracking table already has an entry for `10.0.0.15:51244 ↔ 203.0.113.50:443` in state `SYN_SENT` (because the internal client genuinely initiated this connection moments earlier), the incoming SYN-ACK matches that entry exactly, and a SYN-ACK is a valid transition from `SYN_SENT`. **The firewall allows it through and updates the table entry to `SYN_RECEIVED`.**

**Case B — no such connection exists (the attack case from Section 4).** If the connection tracking table has **no entry at all** for that 4-tuple — because `10.0.0.15` never actually sent a SYN to `203.0.113.50:443` — the stateful firewall **rejects the packet outright**, regardless of what its header claims and regardless of any static rule that might otherwise say "allow traffic from port 443." There is no tracked conversation this packet could plausibly belong to, so it's treated as unsolicited, unexpected, and dropped.

```
                    Same packet, two firewalls, two different verdicts:

  Stateless filter:  "Source port 443 matches my ALLOW rule."      → PASS
  Stateful firewall: "No tracked connection matches this 4-tuple." → DROP
```

This is the exact, concrete answer to the question Section 1 opened with: **the stateful firewall can tell a real reply from a spoofed one not by looking harder at the packet itself — the packet's header fields are identical in both cases — but by checking the packet against a memory of what actually happened before it arrived.** This is a direct, practical instance of Chapter 83's SYN flood defense (Chapter 83, Section 8) generalized: stateful inspection is what lets a firewall distinguish a genuine handshake in progress from noise designed to look like one.

---

## 7. Connection Tracking, NAT, and the State Table in Practice

Stateful connection tracking isn't just a security feature bolted onto a firewall — it's the same underlying mechanism that makes NAT (Chapter 41) work at all. A NAT gateway has to remember which internal (private-IP, ephemeral-port) tuple corresponds to which externally-visible (public-IP, translated-port) tuple, for exactly the same reason a stateful firewall remembers which connections are legitimate: **so a reply packet arriving from the outside can be correctly matched back to the internal host that solicited it, and nowhere else.** On Linux, this is literally the same subsystem — `conntrack`, part of the kernel's netfilter framework — that both NAT and the stateful firewall (`iptables`/`nftables` with `-m state` or `-m conntrack` rules) are built on top of. A single table entry serves both jobs simultaneously: it tells NAT how to rewrite addresses, and it tells the firewall whether an inbound packet corresponds to a real, tracked outbound request.

This has a real production consequence worth naming explicitly: the connection tracking table is **finite memory on a real device**. Under a sufficiently large flood of new connection attempts (echoing Chapter 83's SYN flood, but now aimed at exhausting the *firewall's* table instead of a single server's backlog), a stateful firewall's own conntrack table can fill up, at which point it may either refuse new connections outright or (misconfigured) fail open and stop tracking state altogether — which is precisely why high-throughput stateful firewalls need conntrack table sizing and connection-rate limiting as part of their own defense, not just a "set it and forget it" feature.

---

## 8. The Firewall Family Tree

Firewalls are usually described in generations, each adding context the previous one lacked:

| Generation | What it inspects | Can it tell a real reply from a spoofed one? | Example |
|---|---|---|---|
| Packet filter (stateless) | Individual packet headers only | No — no memory across packets | Simple router ACLs |
| Stateful inspection | Packet headers + connection tracking table | Yes, at the TCP/UDP/ICMP level | iptables/nftables (conntrack), most hardware firewalls |
| Circuit-level / proxy | Entire session, often relaying it through a proxy rather than just routing it | Yes, and can also enforce protocol correctness end-to-end | SOCKS proxies |
| Next-generation firewall (NGFW) | Application identity, user identity, deep packet inspection of payload | Yes, plus can distinguish applications sharing a port (e.g. all HTTPS on 443) | Palo Alto, Fortinet-class appliances |
| Web Application Firewall (WAF) | Full HTTP request/response semantics, application-layer payload | Not its job — operates one layer above connection-level state entirely | ModSecurity, AWS WAF, Cloudflare WAF |

Each generation up this table adds context (state, then payload content, then application semantics) at the cost of more processing per packet or request — which is exactly why a real production network typically layers several of these together rather than picking just one: a stateless ACL or cloud security group at the outermost edge for cheap, fast filtering, a stateful firewall behind it for connection-level correctness, and a WAF specifically in front of web application servers for content-aware defense. That last layer is where the rest of this chapter goes next.

---

## 9. Web Application Firewalls — Application-Aware Defense

**The gap a stateful firewall structurally cannot close.** A stateful firewall, no matter how sophisticated its connection tracking, answers questions about connections — is this a legitimate, in-progress TCP session on an allowed port? It has no concept of what an HTTP request's *body* or *query string* actually contains, because that's payload, not header, and reading it requires understanding the application protocol running on top of TCP (Chapter 71), not just TCP itself. A perfectly legitimate, fully established, correctly-stated TCP connection to port 443 can still carry an HTTP request designed to exploit the web application listening on the other end — a SQL injection attempt, a cross-site scripting payload, a path traversal attempt — and nothing about that connection's TCP state machine looks any different from a completely benign request.

**What a WAF actually is.** A Web Application Firewall sits logically in front of a web server (often literally as a reverse proxy — Chapter 76's pattern — or as a module inside the web server itself, or as a cloud/CDN-integrated service) and inspects the **full HTTP request**: method, path, query string, headers, cookies, and body — the same fields Chapter 71 and Chapter 72 taught you to read — looking for patterns known to be associated with common web application attacks.

**Two philosophies of matching:**

- **Negative security model (signature/pattern matching)**: maintain a list of known-bad patterns and block requests matching them — for example, flagging a request whose query string contains SQL syntax fragments that have no business appearing in ordinary input (`' OR '1'='1`, `UNION SELECT`), or HTML/JavaScript fragments in a field meant for plain text (`<script>`). This is analogous to antivirus signature matching — effective against known, common attack shapes, structurally blind to anything novel.
- **Positive security model (allowlisting)**: define exactly what a *legitimate* request to a given endpoint should look like (expected parameters, expected types, expected lengths) and reject anything that deviates, rather than trying to enumerate every possible bad pattern. This is more precise but requires detailed, maintained knowledge of the application's actual expected inputs — expensive to build and keep current as the application changes.

Most real WAF deployments (ModSecurity's OWASP Core Rule Set, cloud-provider WAFs) lean heavily on the negative model for broad, low-maintenance coverage, sometimes layering positive-model rules for specific, high-value endpoints.

---

## 10. Worked Example: A WAF Evaluating an HTTP Request

Consider a login form that submits a username to a backend query. Two requests arrive at a WAF sitting in front of the web server, both perfectly valid, fully-established HTTPS connections as far as any stateful firewall (Section 5) can tell — the WAF is the only layer that reads the actual field content:

```
Request A (legitimate):
  POST /login HTTP/1.1
  Host: example.com
  Content-Type: application/x-www-form-urlencoded

  username=alice&password=hunter2

Request B (SQL injection attempt, explained defensively — mechanism only):
  POST /login HTTP/1.1
  Host: example.com
  Content-Type: application/x-www-form-urlencoded

  username=admin' OR '1'='1&password=x
```

A stateful firewall sees two identical-looking established TCP connections on port 443, carrying HTTP, both correctly formed — nothing in the TCP layer distinguishes them at all. A WAF applying signature rules parses the request body, recognizes the `username` field's value contains an unescaped single quote followed by SQL boolean-logic syntax (`OR '1'='1`) — a classic pattern indicating an attempt to alter the structure of a backend SQL query rather than supply an ordinary username — and blocks Request B before it ever reaches the application server, typically logging the attempt and returning a generic 403 to the client. Request A, containing only expected alphanumeric field values, passes through untouched. **This is precisely the "application-aware" layer that neither a stateless filter nor a stateful firewall can provide, because both operate entirely below the layer where "username" and "SQL syntax" have any meaning at all.**

---

## 11. WAF's Real Limitations

A WAF is a genuinely valuable defense-in-depth layer, and it is not, and was never designed to be, a substitute for secure application code. Its real, well-documented limits:

- **Encoding and obfuscation evasion.** Signature matching looks for known bad patterns in the request as received — an attacker who encodes a payload differently than the signature expects (alternate character encodings, case variation, comment injection inside the malicious syntax itself) can sometimes slip past a signature that would have caught the unobfuscated version, in an ongoing arms race between rule authors and evasion techniques.
- **Zero-day and application-specific logic flaws.** A WAF's rules are built from known attack patterns; a genuinely novel technique, or a flaw specific to one application's own business logic (for example, a discount code that can be applied multiple times due to a race condition, with no "malicious syntax" anywhere in the request at all), produces requests that are syntactically completely normal and simply exploit a flaw in what the application *does* with valid-looking input — invisible to pattern matching by design.
- **TLS visibility.** Chapter 82 established that HTTPS traffic is encrypted end-to-end by default. A WAF can only inspect HTTP content if it terminates TLS itself (decrypting, inspecting, and typically re-encrypting before forwarding to the origin — the same TLS offloading pattern Chapter 96 discusses for CDNs) — meaning the WAF becomes another endpoint that must be trusted with the plaintext, and a WAF deployed purely as a network-level packet inspector with no TLS termination capability sees nothing but ciphertext.
- **False positives and false negatives.** Overly broad signatures block legitimate requests that happen to contain innocuous text resembling an attack pattern (a support ticket containing the literal string `SELECT * FROM` as an example in a bug report, for instance), which real teams have to spend ongoing effort tuning; overly narrow signatures miss real attacks that don't match a known shape. Tuning a WAF is a continuous operational task, not a one-time setup.
- **Doesn't replace secure coding.** Parameterized queries (preventing SQL injection at the source, regardless of what the WAF does or doesn't catch), proper output encoding (preventing XSS), and correct authorization checks are the actual fix for the underlying vulnerabilities a WAF's rules are trying to compensate for at the network edge. A WAF is best understood as **virtual patching** — a fast, network-level mitigation for a known vulnerability class while (or instead of, if the code can't be fixed promptly) the underlying application code is properly fixed — not a permanent substitute for fixing the code.

---

## 12. Where Firewalls Actually Sit in Real Networks

```
              Internet
                 │
        ┌────────▼────────┐
        │  Perimeter       │   stateless ACLs / border router filtering
        │  filtering       │   (Section 3 — cheap, first line of defense)
        └────────┬────────┘
                  │
        ┌─────────▼────────┐
        │  Stateful         │   connection tracking, NAT (Sections 5-7)
        │  firewall          │
        └─────────┬────────┘
                  │
        ┌─────────▼────────┐
        │  WAF / reverse     │   HTTP-aware inspection (Sections 9-11)
        │  proxy             │   in front of web servers specifically
        └─────────┬────────┘
                  │
        ┌─────────▼────────┐
        │  Web / app         │
        │  servers           │
        └───────────────────┘
```

In cloud environments (Chapter 97 covers this in depth), the same layered idea reappears under different names: **security groups** (stateful, instance-level firewalls attached to individual cloud resources), **Network ACLs** (stateless, subnet-level filtering), and a managed **WAF service** (AWS WAF, Cloudflare, Azure Front Door's WAF) sitting in front of a load balancer or CDN. The concepts in this chapter transfer directly — only the deployment model (dedicated hardware appliance vs. software-defined cloud construct) differs.

Host-based firewalls (Windows Firewall, `iptables`/`nftables` running directly on a Linux server) apply everything in this chapter at the level of a single machine rather than a network boundary — the same stateless-vs-stateful distinction, the same connection tracking table, just protecting one host instead of an entire network segment.

---

## 13. A Hands-On Experiment

1. On a Linux machine, run `sudo iptables -L -n -v` (or `sudo nft list ruleset` on a system using nftables) and look for rules referencing `ESTABLISHED,RELATED` — this is the literal stateful rule that says "if this packet matches an existing connection in the tracking table, allow it," the real-world equivalent of Section 5's mechanism.
2. Run `sudo conntrack -L` (installing `conntrack-tools` if needed) to view the live connection tracking table on a Linux machine with active network connections — you'll see entries with source/destination tuples and states (`ESTABLISHED`, `TIME_WAIT`) directly mirroring Section 5's diagram.
3. Look up the OWASP Core Rule Set (the open-source rule set used by ModSecurity, a widely deployed open-source WAF) and read a handful of its published rule descriptions — notice how many are explicitly pattern-based (Section 9's negative security model) and consider, for each, one plausible way legitimate traffic could accidentally trigger it.
4. If you run a personal website behind a CDN or hosting provider with WAF logging available, check its dashboard for blocked requests over the last 30 days — real production WAFs log far more blocked attempts than most people expect, mostly automated scanning traffic rather than targeted attacks.

---

## 14. Common Misconceptions

- **"A firewall protects against everything."** A firewall (of any generation) controls what traffic is allowed to reach a destination based on network/transport-layer information, or in a WAF's case, application-layer request content — it does nothing about a legitimate, allowed connection carrying an application vulnerability it wasn't specifically written to detect, or about an attacker who has already compromised an internal machine allowed to talk freely on the network.
- **"Stateful firewalls check payload content."** They check *connection state* — TCP/UDP/ICMP-level context — not the meaning of the data inside a request. That's the WAF's job, not the stateful firewall's.
- **"A WAF replaces the need for secure application code."** As Section 11 discussed, a WAF is a mitigation and a compensating control, most honestly described as virtual patching — not a substitute for parameterized queries, output encoding, and correct authorization logic in the application itself.
- **"NAT is a security feature."** NAT (Chapter 41) exists to conserve IPv4 addresses and happens to obscure internal addressing as a side effect, which provides a mild security benefit — but the actual security enforcement comes from the stateful connection tracking and firewall rules that typically run alongside it, not from address translation itself.
- **"More firewall layers always means more security with no downside."** Every layer in Section 12's stack adds latency and operational complexity (more rules to maintain, more places a misconfiguration can silently block legitimate traffic) — real designs balance defense-in-depth against these costs rather than stacking every possible layer everywhere.

---

## 15. What's Simplified Here

This chapter covers the conceptual core — stateless vs. stateful evaluation and application-aware WAF inspection — without going into the specific rule syntax of any one product (iptables/nftables, cloud security group JSON, ModSecurity's SecRule language), the internals of deep packet inspection engines used in next-generation firewalls, or the machine-learning-based anomaly detection increasingly layered on top of signature-based WAF rules in commercial products. Real production firewall and WAF deployments also involve logging, alerting, and incident response workflows that are out of scope here — this chapter's job was the reasoning behind *why* each layer exists and what specifically it can and cannot see.

---

## 16. Interview Questions & Model Answers

**Beginner: "What's the difference between a stateless and a stateful firewall?"**

*Model answer:* "A stateless firewall evaluates every packet independently against a fixed set of rules matching header fields like source/destination IP, port, and protocol — it has no memory of previous packets. A stateful firewall maintains a connection tracking table recording every connection currently in progress, and evaluates incoming packets against that table — asking not just 'does this header look allowed' but 'does this packet actually correspond to a connection I know is really happening.' This lets a stateful firewall reject a packet that has a header matching an allow rule but doesn't correspond to any real, tracked connection — something a stateless firewall structurally cannot do, because it never remembers what came before."

**Intermediate: "Give a concrete example of a packet that a stateless firewall would allow but a stateful firewall would correctly reject."**

*Model answer:* "Suppose a stateless firewall has a rule allowing any inbound TCP packet with source port 443, intended to let through replies from web servers a client connected to. An attacker can craft a SYN-ACK packet with a forged source port of 443, destined for an internal client that never actually initiated any connection to that source. The stateless firewall matches the header against its rule and lets it through. A stateful firewall checks its connection tracking table for an entry matching that 4-tuple in a state where a SYN-ACK would be a valid next step — finds none, because no SYN was ever sent from that internal client to that destination — and drops the packet, regardless of what its header claims."

**Advanced: "Why can't a stateful firewall detect a SQL injection attempt inside an HTTP request, and what does a WAF add that closes that gap?"**

*Model answer:* "A stateful firewall's decisions are based entirely on connection-level information — source/destination addresses and ports, protocol, and TCP/UDP state — it has no logic for parsing or understanding HTTP request bodies or query strings, because that requires understanding the application protocol running on top of the transport layer, not just the transport layer's own handshake and state machine. A fully established, correctly-behaving TCP connection carrying a malicious HTTP request looks identical, at the connection-tracking level, to one carrying a benign request — there's nothing in the TCP state machine that differs. A WAF closes this gap by operating one layer up: it parses the actual HTTP request — method, headers, query string, body — and matches its content against patterns known to indicate attacks like SQL injection or cross-site scripting, something that requires understanding HTTP and the application's expected input shapes, not just the underlying connection's state."

---

## 17. Exercises

### Easy

1. In your own words, explain why a stateless firewall needs a rule like "allow inbound traffic with source port 443" at all, and what problem that rule is trying to solve.
2. What specific piece of information does a stateful firewall have access to that a stateless firewall does not?
3. Name one type of attack a WAF is specifically designed to catch that neither a stateless nor a stateful firewall could ever detect, and explain why.

### Medium

4. Explain the relationship between NAT's translation table (Chapter 41) and a stateful firewall's connection tracking table — why are they so often implemented as literally the same underlying data structure (as in Linux's conntrack)?
5. A company deploys a WAF but continues to build SQL queries by directly concatenating user input into SQL strings instead of using parameterized queries. Explain, using Section 11's discussion of virtual patching, why this is a risky long-term security posture even if the WAF is well-tuned.
6. Explain why a WAF must either terminate TLS itself or be deployed behind whatever does terminate it, in order to inspect HTTP request content at all.

### Hard

7. A stateful firewall's connection tracking table has a fixed maximum size. Design (in prose) an attack scenario that targets this limit directly rather than targeting a web server's own resources, and explain how it differs mechanically from the SYN flood described in Chapter 83.
8. Compare the negative security model (signature/blocklist) and positive security model (allowlist) approaches to WAF rule design, and argue for a specific hybrid approach for a hypothetical company's public login endpoint versus its internal admin API — justify the trade-off in each case using the limitations discussed in Section 11.

---

## 18. Summary

| Term | Meaning |
|---|---|
| Stateless packet filter | Evaluates each packet independently against static header-matching rules; fast, no memory of prior packets |
| Access Control List (ACL) | The ordered list of allow/deny rules a stateless filter evaluates against |
| Stateful firewall | Maintains a connection tracking table and evaluates packets against tracked, in-progress connections |
| Connection tracking table (conntrack) | The state table recording each active connection's 4-tuple and current state, shared with NAT |
| Next-generation firewall (NGFW) | A firewall adding application/user identity awareness and deep packet inspection on top of stateful tracking |
| Web Application Firewall (WAF) | An application-layer-aware filter inspecting full HTTP request/response content, typically for SQLi/XSS-style patterns |
| Negative / positive security model | Blocking known-bad patterns vs. allowing only known-good request shapes |
| Virtual patching | Using a WAF rule as a fast, network-level mitigation for a known vulnerability while the underlying code is fixed |

Firewalls and WAFs decide what's allowed to reach a server at all. Chapter 85 turns to the complementary problem — once you're allowed through, how do you make the entire path between two networks private and safe over a public Internet you don't own, using tunneling and VPNs.
