# Chapter 123: Real Debugging Scenarios, Solved

> **"A method proves itself only against real problems. Here are nine — each one a complaint a real engineer has actually received, each one solved with nothing but Chapter 122's checklist and Chapters 119-121's tools."**

---

## Table of Contents

1. [Why This Chapter Exists](#1-why-this-chapter-exists)
2. [How to Read Each Scenario](#2-how-to-read-each-scenario)
3. [Scenario 1: DNS Resolves but HTTP Fails](#3-scenario-1-dns-resolves-but-http-fails)
4. [Scenario 2: TCP Connects but the App Hangs](#4-scenario-2-tcp-connects-but-the-app-hangs)
5. [Scenario 3: Works on Wi-Fi but Not Mobile Data](#5-scenario-3-works-on-wi-fi-but-not-mobile-data)
6. [Scenario 4: One Machine Can't Reach Another](#6-scenario-4-one-machine-cant-reach-another)
7. [Scenario 5: Intermittent Packet Loss](#7-scenario-5-intermittent-packet-loss)
8. [Scenario 6: Sudden Latency Spikes](#8-scenario-6-sudden-latency-spikes)
9. [Scenario 7: Works Internally but Not Externally](#9-scenario-7-works-internally-but-not-externally)
10. [Scenario 8: TLS Handshake Failures](#10-scenario-8-tls-handshake-failures)
11. [MTU, Fragmentation, and PMTUD From First Principles](#11-mtu-fragmentation-and-pmtud-from-first-principles)
12. [Scenario 9: MTU-Caused Mysterious Failures](#12-scenario-9-mtu-caused-mysterious-failures)
13. [Cross-Scenario Patterns Worth Memorizing](#13-cross-scenario-patterns-worth-memorizing)
14. [Hands-On Experiment](#14-hands-on-experiment)
15. [Common Misconceptions](#15-common-misconceptions)
16. [Production Notes](#16-production-notes)
17. [What's Simplified Here](#17-whats-simplified-here)
18. [Interview Questions & Model Answers](#18-interview-questions--model-answers)
19. [Exercises](#19-exercises)
20. [Summary and the Bridge to Chapter 124](#20-summary-and-the-bridge-to-chapter-124)

---

## 1. Why This Chapter Exists

Chapter 122 built a method. A method that's never been exercised is just a diagram. This chapter runs that exact method — the layer checklist, the decision tree, the client/server/network isolation — against nine realistic, independently common scenarios, each one taken from symptom to confirmed root cause with real commands and real (representative) output. By the end, you will have seen the playbook actually work nine separate times, against nine genuinely different failure classes, which is what turns a checklist you've read into a checklist you reach for by instinct.

## 2. How to Read Each Scenario

Every scenario below follows the same four-part structure, deliberately mirroring Chapter 122's own arc: **Symptom** (exactly what was reported, in the vague, unhelpful form a real user or colleague would actually use), **Investigation** (the actual commands run, in order, with realistic output), **Root Cause** (the specific, named mechanism responsible), and **Fix and Lesson** (what resolves it, and which general principle from Chapter 122 or earlier volumes it reinforces).

## 3. Scenario 1: DNS Resolves but HTTP Fails

**Symptom:** "I can `ping` and `dig` the site fine, but the page won't load in my browser — it just spins."

**Investigation:**

```
$ dig api.example.com +short
203.0.113.44

$ ping -c 3 203.0.113.44
3 packets transmitted, 3 received, 0% packet loss

$ nc -zv 203.0.113.44 443
nc: connect to 203.0.113.44 port 443 (tcp) failed: Connection refused

$ dig api.example.com @1.1.1.1 +short
198.51.100.9
```

**Root Cause:** the local DNS resolver has a **stale cached A record** (Chapter 68) pointing to `203.0.113.44` — an old IP the service was migrated away from — while every other resolver (`1.1.1.1`, queried directly in the last command) already has the current, correct answer, `198.51.100.9`. The old IP still exists (it responds to `ping`, since something else now owns that address or it's still partially routed) but nothing is listening on port 443 there anymore, hence the immediate `Connection refused` rather than a timeout.

**Fix and Lesson:** flush the local resolver cache (or wait out the TTL, Chapter 68), and confirm the fix by re-running `dig` before re-testing HTTP. The lesson is Chapter 122, Section 7's entire point made concrete: **"DNS resolves" and "DNS resolves *correctly*" are different claims** — a successful `dig` only proves an answer came back, not that the answer is the current, intended one. Cross-checking against a second, independent resolver is the fastest way to tell "my DNS is stale" apart from "the site itself is actually broken at that IP."

## 4. Scenario 2: TCP Connects but the App Hangs

**Symptom:** "The connection opens fine — I can see it in `netstat` — but the request just sits there forever. No error, no response."

**Investigation:**

```
$ curl -v --max-time 15 https://api.example.com/v1/report
* Connected to api.example.com (203.0.113.44) port 443
* TLS handshake ... SSL certificate verify ok.
> GET /v1/report HTTP/1.1
> Host: api.example.com
>
* (nothing further — curl times out after 15s with no response)

$ ss -t dst 203.0.113.44
State  Recv-Q  Send-Q  Local Address:Port     Peer Address:Port
ESTAB  0       0       10.0.1.20:51900        203.0.113.44:443
```

TCP and TLS both completed cleanly (per Chapter 122, Section 8-9's checklist) — `ss` confirms the connection is genuinely `ESTAB` on both sides, not stuck in a handshake state at all. The request was sent (visible in `curl -v`'s `>` lines), and nothing at all came back — not even a slow response, an error, or a RST. On the server, checking what the request is actually doing:

```
# On the server, inspecting the application's own state
$ ps aux | grep report-worker
report-worker   PID 8821   99.8% CPU   00:14:32 elapsed

$ (application logs)
[10:41:02] Received request GET /v1/report
[10:41:02] Acquiring lock on report_cache...
(no further log lines for this request)
```

**Root Cause:** the request handler is **blocked waiting on a lock** that never gets released — a classic application-level deadlock or a downstream dependency (a database query, another internal service) that itself hung, with the request thread simply parked indefinitely rather than timing out and returning an error. This is explicitly *not* a networking problem in the sense this course has built up: every layer through TCP and TLS did its job correctly and continues to do so (the connection stays open and healthy the entire time).

**Fix and Lesson:** the fix lives entirely in application code — add a timeout to the lock acquisition (or the downstream call actually stuck behind it) so the request fails fast with a clear error instead of hanging forever. The lesson is Chapter 122, Section 10's point stated plainly: **once TCP and TLS are confirmed healthy, a hang is not a networking problem anymore, and no tool from Chapters 119-121 can diagnose it further** — the investigation must move to the application's own logs, thread dumps, and dependency health, which is precisely why the checklist's job is to *rule networking in or out*, not to solve every problem itself.

## 5. Scenario 3: Works on Wi-Fi but Not Mobile Data

**Symptom:** "The app works perfectly on my phone over Wi-Fi at home, but the moment I switch to mobile data, it fails to connect."

**Investigation:**

```
# On Wi-Fi
$ curl -v https://api.example.com/health
< HTTP/1.1 200 OK

# On mobile data (via a laptop tethered to the phone's hotspot, for easier tooling)
$ curl -v https://api.example.com/health
* Connected to api.example.com (203.0.113.44) port 443
* TLS handshake, Client hello (1)
* (connection stalls, eventually times out)

$ curl -v -6 https://api.example.com/health
* Trying [2606:2800:220:1:248:1893:25c8:1946]:443...
* connect to 2606:2800:220:1:248:1893:25c8:1946 port 443 failed: Network unreachable
```

**Root Cause:** the mobile carrier's network provides **IPv6-only or IPv6-preferred connectivity** (a common, deliberate carrier configuration, since Chapter 42 established IPv4 exhaustion pushed mobile networks toward IPv6 earlier and more completely than most home ISPs), and `api.example.com` either has no working AAAA record, or its IPv6 path is broken/unreachable while its IPv4 path (used successfully over Wi-Fi/the home ISP) works fine. The client's dual-stack happy-eyeballs behavior (attempting IPv6 first per RFC 8305-style logic) stalls on the broken IPv6 path before falling back, producing exactly the "hangs, then eventually fails" symptom on mobile data specifically.

**Fix and Lesson:** either fix the broken IPv6 path server-side (Chapter 42-43), or, as an immediate mitigation, remove the faulty AAAA record until IPv6 is properly working, forcing all clients back onto the known-good IPv4 path. The lesson connects directly to Chapter 122, Section 12's client-isolation logic taken one step further: **"works on one network, not another" isolates the fault to something specific about that network path** — and mobile carriers' IPv6-forward posture is a specific, common enough cause that it deserves to be checked directly (`curl -6` / `curl -4` forcing one family or the other) rather than treated as a mysterious, unexplainable carrier quirk.

## 6. Scenario 4: One Machine Can't Reach Another

**Symptom:** "Server A can't reach Server B on port 5432 (Postgres), but both are in the same cloud VPC and I can SSH into both."

**Investigation:**

```
# From Server A
$ nc -zv 10.0.2.15 5432
nc: connect to 10.0.2.15 port 5432 (tcp) timed out

$ ping -c 3 10.0.2.15
3 packets transmitted, 3 received, 0% packet loss   # basic L3 reachability is fine

# Checking the cloud VPC Flow Logs (Chapter 121, Section 9) for this exact flow
$ (flow log query filtered to srcaddr=10.0.1.20 dstaddr=10.0.2.15 dstport=5432)
srcaddr=10.0.1.20 dstaddr=10.0.2.15 dstport=5432 protocol=6 action=REJECT
```

**Root Cause:** `ping` succeeds because ICMP is a different protocol from TCP and can be — and here, is — governed by a completely different security-group rule; the actual TCP SYN to port 5432 is being **explicitly rejected by Server B's security group** (Chapter 97), which allows ICMP from the VPC's CIDR but was never updated to allow TCP/5432 from Server A's specific subnet after a recent network re-segmentation.

**Fix and Lesson:** add an explicit inbound rule on Server B's security group allowing TCP/5432 from Server A's subnet. The lesson is a direct, sharp instance of Chapter 122, Section 8's central warning: **`ping` succeeding proves nothing about a specific port's TCP reachability**, since ICMP and TCP are filtered independently by design — and the VPC Flow Logs' explicit `action=REJECT` field (Chapter 121, Section 9) turned what would otherwise be a guessing exercise between "wrong port," "service not running," and "firewall rule" into a one-line, directly confirmed answer.

## 7. Scenario 5: Intermittent Packet Loss

**Symptom:** "Our video calls to the Mumbai office are fine most of the time, but every few minutes there's a burst of garbled audio for a few seconds."

**Investigation:**

```
$ ping -i 0.2 -c 300 mumbai-office.internal | grep -c "icmp_seq"
300
$ ping -i 0.2 -c 300 mumbai-office.internal | grep -v "time=" 
Request timeout for icmp_seq 812
Request timeout for icmp_seq 813
Request timeout for icmp_seq 814
(a burst of ~8-10 consecutive drops, then resumes cleanly)

$ mtr -n -c 500 --report mumbai-office.internal
  1.|-- 10.0.0.1                  0.0%   500    0.4   0.5   0.4   2.1   0.2
  2.|-- 203.0.113.1                0.0%   500    2.1   2.4   2.0   5.8   0.4
  3.|-- 198.51.100.9               1.8%   500   45.2  46.9  44.8 190.3  22.1
  4.|-- 172.16.9.4                 1.8%   500   46.1  47.5  45.9 191.0  22.4
  5.|-- 10.10.5.20 (mumbai-office) 1.8%   500   46.4  47.8  46.2 191.5  22.6
```

**Root Cause:** loss (1.8%) appears starting at hop 3 and **persists at the same rate through every subsequent hop, including the final destination** — the trustworthy pattern for real loss established in Chapter 120, Section 9, ruling out an ICMP-deprioritization artifact. Cross-referencing hop 3's IP against known infrastructure identifies it as a specific transit provider's link that, per that provider's own status page, is running near capacity during that time window — the bursty timing (every few minutes, in short clusters) is consistent with periodic congestion-driven queue overflow (Chapter 120, Section 8) at that link rather than a constant, low-level loss rate.

**Fix and Lesson:** since the loss is on a third-party transit link outside direct control, the practical fix is either an alternate network path (a different ISP/transit arrangement, or a VPN/SD-WAN route around the congested provider) or escalating directly to that provider with the `mtr` evidence in hand. The lesson: Chapter 120, Section 9's "loss persisting past the hop it first appears at" rule is what separates an actionable finding from a red herring here — without it, hop 3's number alone would be ambiguous.

## 8. Scenario 6: Sudden Latency Spikes

**Symptom:** "Everything was fine all week, and starting this morning at roughly 9am, response times to our database tripled and haven't recovered."

**Investigation:**

```
# Grafana dashboard (Chapter 121, Section 12) showing a Prometheus-recorded
# latency metric for the DB connection, with a sharp step change at 09:02

$ (Grafana panel: rate of change and absolute value of db_query_duration_seconds,
   annotated against other metrics on the same dashboard)

  09:02:14  network_interface_errors_total{device=eth0} starts climbing
  09:02:14  db_query_duration_seconds (p50) steps from 4ms to 14ms and stays there

$ ethtool -S eth0 | grep -i error
     rx_crc_errors: 84213
     rx_errors: 84213
```

**Root Cause:** correlating two metrics on the same Grafana dashboard (exactly the value of storing history, Chapter 121, Section 2's original motivation) shows the latency step change lines up, to the second, with a spike in **physical-layer CRC errors** (Chapter 19) on the database server's network interface — consistent with a failing NIC, a degrading cable, or a bad switch port (Chapter 17's attenuation/interference discussion) rather than anything at the application or database layer. Every affected query has to retransmit lost/corrupted segments (Chapter 60), inflating latency uniformly and persistently rather than intermittently, matching the "steps up and stays" pattern rather than Scenario 5's bursty one.

**Fix and Lesson:** physically reseat or replace the cable, and if errors persist, the NIC or switch port itself — a hardware, Layer 1 fix for a symptom that first appeared as a database performance problem. The lesson: **a persistent, immediate, and unrecovering latency step change is a different signature from Scenario 5's intermittent bursts**, and reaching for continuous, correlated metrics (Chapter 121) rather than a one-off `ping` was what made the true, physical-layer cause visible in seconds instead of requiring a slow elimination of every layer above it first.

## 9. Scenario 7: Works Internally but Not Externally

**Symptom:** "Everyone in the office can reach the new staging site, but nobody outside the office network can load it at all."

**Investigation:**

```
# From inside the office network
$ curl -I http://staging.example.com
HTTP/1.1 200 OK

# From an external host (a cloud VM outside the office network)
$ curl -v http://staging.example.com
* Trying 203.0.113.77:80...
* connect to 203.0.113.77 port 80 failed: Connection timed out

$ dig staging.example.com +short
203.0.113.77                          # same IP resolved from both locations

$ (checking the office firewall/router's outbound NAT and port-forwarding rules)
No inbound port-forwarding rule exists for 203.0.113.77:80 -> internal staging server
```

**Root Cause:** `staging.example.com`'s public DNS record correctly points to the office's public IP, but the office's edge router/firewall (Chapter 41, Chapter 84) has **no inbound NAT/port-forwarding rule** allowing external traffic on port 80 through to the internal staging server — internal clients succeed because their traffic never has to cross that NAT boundary at all (their requests resolve to the same public IP, but the router's internal DNS/hairpin-NAT handling, or a split-horizon internal DNS entry, routes them directly), while genuinely external clients hit the router's public interface and get silently dropped, since by default an edge firewall denies unsolicited inbound connections (Chapter 84's stateful default-deny posture).

**Fix and Lesson:** add the missing inbound port-forwarding/NAT rule on the office router for port 80 (and 443) to the internal server's private address. The lesson is Chapter 122, Section 12's isolation table put to direct use: **"works internally, fails externally" is one of the cleanest possible isolating signals in this entire chapter** — it immediately points at the network boundary between the two populations (the NAT/firewall edge) rather than at DNS, the application, or the server itself, all of which are proven fine by the fact that internal access works perfectly.

## 10. Scenario 8: TLS Handshake Failures

**Symptom:** "Some users get a browser security warning when visiting our site; others say it loads fine."

**Investigation:**

```
$ openssl s_client -connect www.example.com:443 -servername www.example.com </dev/null 2>&1 | grep -E 'Verify return|subject|notAfter'
subject=CN=www.example.com
notAfter=Aug 8 23:59:59 2026 GMT
Verify return code: 0 (ok)

$ date -u
Sat Aug  9 10:20:00 UTC 2026
```

**Root Cause:** the certificate's `notAfter` date, `Aug 8 23:59:59 2026`, is **one day in the past** relative to the current date — the TLS certificate (Chapter 81) expired less than 24 hours ago. Users "for whom it loads fine" are almost certainly hitting a different edge server or CDN node (Chapter 96) that already received a renewed certificate through an automated rotation process, while others are hitting a node that hasn't yet picked up the renewal — a partial, in-progress rollout rather than a single, universal failure, which explains the inconsistent reports perfectly.

**Fix and Lesson:** force or verify completion of the certificate renewal and rollout across every edge node (most production setups use automated renewal via ACME/Let's Encrypt-style protocols specifically to prevent this), and add active certificate-expiry monitoring (a Grafana alert, Chapter 121 Section 12, on days-until-expiry) so this becomes a scheduled non-event rather than a surprise incident. The lesson: **`Verify return code` and the certificate's exact validity dates are the single fastest, most direct check for any TLS failure** — Chapter 122, Section 9's checklist step exists precisely to catch this class of problem (expiry, hostname mismatch, untrusted CA) before wasting time suspecting the network path or the application.

## 11. MTU, Fragmentation, and PMTUD From First Principles

Scenario 9 requires one mechanism this course has referenced but never built from first principles: what happens when a packet is simply **too big** for a link along its path.

**The problem, stated plainly:** every physical and data-link medium has a maximum frame size it can carry in one piece — Chapter 28 already established Ethernet's standard maximum frame size caps the payload at 1500 bytes (the **Maximum Transmission Unit**, or MTU, for that link). But a path from your laptop to a server crosses many links — Wi-Fi, your home router, your ISP, possibly a VPN tunnel, several ISP backbone hops, a data center network — and **there is no guarantee every link along that entire path shares the same MTU.**

**Naive assumption:** since IP packets can be up to 65,535 bytes, and Ethernet says 1500, applications should just... send whatever size makes sense, and let the network sort it out.

**Why that breaks:** if a router receives an IP packet larger than the MTU of the link it needs to forward it out on, it has exactly two choices, and IPv4 explicitly supports both: **fragment** the packet into smaller pieces that each fit the outgoing link's MTU (with the destination reassembling them later), or, if the sender marked the packet "Don't Fragment" (the **DF bit**, a real flag in the IP header), **drop it and report back with an ICMP error**.

**Why fragmentation itself is a problem worth avoiding, not just a mechanism worth using:** fragmenting is expensive for routers (a slow, CPU-costly path historically deprioritized in high-speed router hardware precisely because it's rare and awkward), it multiplies the chance of loss (Chapter 120, Section 8 — losing *any one* fragment means the *entire original packet* must be discarded and retransmitted, since partial reassembly isn't possible), and it plays badly with modern stateful firewalls and security devices that struggle to inspect fragmented traffic correctly. For this combination of reasons, **virtually all modern IP traffic sets the DF bit and relies on discovering the right size in advance instead of fragmenting in flight.**

**The real solution: Path MTU Discovery (PMTUD).** A sender starts by assuming its own local interface's MTU (commonly 1500 bytes for Ethernet) and sends packets with **DF set**. If some router along the path has a smaller-MTU outgoing link (a common real case: a VPN tunnel, Chapter 85, whose encryption overhead eats into the usable payload size, or certain older DSL/PPPoE links with a famously reduced 1492-byte MTU), that router cannot forward the oversized, DF-marked packet — so it drops it and sends back an **ICMP "Fragmentation Needed" message (Type 3, Code 4)**, which critically includes the actual MTU of the link that couldn't take it. The sender receives this, lowers the size it uses for that specific destination, and retries — converging, packet by packet if necessary, on the largest size the *entire path* can actually carry without fragmentation, which the running example's diagram below makes concrete:

```
Client (MTU 1500)                                                    Server
      |                                                                 |
      |--- IP packet, 1500 bytes, DF=1 ------->  Router A (out-link MTU 1400)
      |                                                |
      |<---- ICMP Type 3 Code 4: "Fragmentation Needed, MTU=1400" ------|
      |
      |--- IP packet, 1400 bytes, DF=1 ------->  Router A --> ... --> Server
      |                                          (fits every remaining
      |                                           link's MTU; delivered)
```

**Why this mechanism can silently break entirely — the exact setup Scenario 9 walks through:** PMTUD's entire feedback loop depends on that ICMP "Fragmentation Needed" message actually getting back to the original sender. **If a firewall anywhere along the path (commonly, overzealously, blocking all ICMP as a blanket "security" measure — Chapter 84's tooling misused) drops that specific ICMP message, the sender never learns its packets are too big, keeps sending at the original size with DF set, and those packets keep being silently discarded** at the same router, forever, with no error ever reported back to either side — a failure mode with the specific, notorious name **"black-holed PMTUD."**

## 12. Scenario 9: MTU-Caused Mysterious Failures

**Symptom:** "SSH sessions to our new VPN-connected office work fine for typing commands, but the moment I `cat` a large file or run anything with big output, the session just freezes."

**Investigation:**

```
$ ssh office-server "cat small_file.txt"       # works fine, few bytes
$ ssh office-server "cat large_report.log"     # hangs indefinitely, no output at all

$ ping -M do -s 1472 office-server              # -M do sets DF; 1472 + 28-byte
                                                  # ICMP/IP header = 1500 total
PING office-server: 1472 data bytes
--- office-server ping statistics ---
1 packets transmitted, 0 received, 100% packet loss

$ ping -M do -s 1372 office-server              # try a smaller size, guessing
                                                  # for the VPN's typical overhead
64 bytes from office-server: icmp_seq=1 ttl=63 time=41.2 ms
--- office-server ping statistics ---
1 packets transmitted, 1 received, 0% packet loss
```

Small SSH commands succeed because their packets are tiny — well under any plausible MTU on the path. The moment real data volume flows (a large file's output), the connection tries to send full-size, 1500-byte, DF-set segments, and those are exactly the ones going silently missing.

```
# Confirming a firewall is eating the PMTUD feedback, by checking for the
# ICMP error a healthy path would show in a packet capture (Chapter 119)
$ sudo tcpdump -i eth0 -n icmp
(nothing — no ICMP "Fragmentation Needed" message ever arrives, despite
 the 1472-byte ping above being dropped somewhere along the path)
```

**Root Cause:** the site-to-site VPN tunnel (Chapter 85) adds encryption/encapsulation overhead that reduces the effective usable payload below the standard 1500-byte Ethernet MTU (a real, extremely common VPN side effect — IPsec or WireGuard overhead routinely costs 50-100+ bytes) — so any packet built assuming the full 1500 bytes doesn't fit once it reaches the tunnel's entry point. PMTUD *should* automatically discover and adapt to this (Section 11's mechanism exists for exactly this situation) — but a firewall somewhere on the path is **blocking the ICMP "Fragmentation Needed" replies** that would tell the sender to shrink its packets, so the sender never finds out, and oversized packets are silently, permanently black-holed instead of being gracefully resized.

**Fix and Lesson:** the durable fix is to explicitly permit ICMP Type 3 (Destination Unreachable, which includes Code 4 Fragmentation Needed) through every firewall on the path — never blanket-block all ICMP. A common, more targeted mitigation specific to VPN/tunnel scenarios is **TCP MSS clamping** — the tunnel endpoint itself rewrites the Maximum Segment Size option (Chapter 65, Section 11) in the TCP handshake to a value that already accounts for the tunnel's overhead, so TCP senders never even attempt to build a segment that would need fragmenting or PMTUD in the first place, sidestepping the ICMP-blocking problem entirely rather than depending on it being fixed. The lesson, and the reason this scenario earns the label "mysterious": **an MTU/PMTUD failure produces a symptom (works for small data, silently hangs for large data, no error message anywhere) that looks nothing like a classic network failure** — no refused connection, no DNS error, no TLS alert — which is exactly why it belongs at the end of Chapter 122's checklist as a specific, named pattern to check for once everything else has passed, rather than something you'd stumble onto by generic layer-by-layer testing alone.

## 13. Cross-Scenario Patterns Worth Memorizing

Stepping back across all nine scenarios, a small number of recurring diagnostic moves did almost all of the real work:

| Pattern | Scenarios It Solved | Why It Works |
|---|---|---|
| Query a second, independent DNS resolver | 1 | isolates local cache/resolver staleness from a true server-side problem |
| Test the raw IP, not just the hostname | 1, 3, 7 | separates DNS from everything above it (Ch 122 §7) |
| Distinguish "refused" from "timed out" | 4, 7, 9 | tells "actively rejected" apart from "silently dropped" |
| Check flow logs / firewall rules directly, don't guess | 4 | turns a guess into a directly confirmed answer |
| Loss persisting past the hop it starts at | 5 | separates real path loss from `mtr`'s ICMP-artifact false positives |
| Correlate two time-series on one dashboard | 6 | finds causation by timing alignment, not by guesswork |
| "Works internally, fails externally" | 7 | immediately isolates the fault to the network boundary between the two |
| Check certificate dates and verify code directly | 8 | the fastest, most direct TLS-layer check available |
| Test with progressively smaller packet sizes and DF set | 9 | directly probes for a path MTU problem instead of assuming one |

## 14. Hands-On Experiment

Reproduce Scenario 9's core technique yourself, safely, against any real host:

```bash
# Find the largest packet your own path to a real server can carry without
# fragmentation, by binary-searching the -s size argument
ping -M do -s 1472 example.com   # 1500 total; likely fails on many real paths
ping -M do -s 1400 example.com   # try smaller
ping -M do -s 1300 example.com   # keep narrowing until one succeeds

# Once you find a passing size, confirm your system's actual PMTUD-discovered
# value for that destination (Linux exposes this directly):
ip route get example.com
# -> look for an "mtu" field in the output if PMTUD has cached a discovered value
```

Running this over a VPN connection, if you have access to one, is especially instructive — VPN overhead is one of the most common real-world triggers for exactly the reduced-MTU condition Scenario 9 walked through.

## 15. Common Misconceptions

- **"MTU problems would show up as an obvious error message."** Scenario 9 is built specifically to correct this — a black-holed PMTUD failure produces total silence, not an error, which is precisely what makes it "mysterious" and worth learning to recognize by its distinctive symptom pattern (small data fine, large data hangs) rather than by an explicit diagnostic message.
- **"Blocking all ICMP is always a safe, purely beneficial security hardening step."** Section 11 shows directly why blanket ICMP blocking is a common, real cause of production outages — PMTUD specifically depends on one ICMP message type reaching the sender, and blocking it trades a marginal, often overstated security benefit for a serious, hard-to-diagnose functional break.
- **"If `ping` (with default, small packets) succeeds, the path is fully healthy."** Every one of Scenario 4's and Scenario 9's investigations depended on going beyond a default-size `ping` — testing a specific port (Scenario 4) or a specific packet size with DF set (Scenario 9) — because a small default ICMP packet can succeed on a path that fails for exactly the traffic that actually matters.
- **"Nine scenarios means there are only nine kinds of networking bugs."** These nine are common, representative categories, chosen because each isolates a genuinely distinct mechanism — real production incidents frequently combine two or more of these patterns at once (an MTU problem masked by an unrelated DNS staleness issue, say), which is exactly why Chapter 122's general method, not a lookup table of nine fixed answers, is the actual durable skill.

## 16. Production Notes

- **TCP MSS clamping (Section 12) is standard, default configuration on most VPN and tunnel products precisely because relying on PMTUD's ICMP round trip working end to end across the public internet is considered unreliable in practice** — production VPN deployments should treat MSS clamping as a required setting, not an optional tweak.
- **Certificate expiry monitoring (Scenario 8) should alert well before expiry — commonly 30, 14, and 3 days out — specifically because a rolling multi-node deployment (CDN edges, load-balanced fleets) can take real time to fully propagate a renewal**, and Scenario 8's "some users affected, others not" symptom is the direct, predictable consequence of catching a renewal too close to the deadline.
- **Flow logs (Chapter 121) and the layer-by-layer checklist (Chapter 122) together turn most of Scenario 4 and Scenario 7's class of problems into minutes-long investigations rather than hours-long ones** — the investment in having those systems already running, before an incident, is what made this chapter's investigations look this fast; without them, several scenarios would have required much slower manual packet-capture work instead.
- **A dashboard that shows multiple correlated metrics on one timeline (Scenario 6) is worth building deliberately, in advance, for every critical service** — the causal link in Scenario 6 was only fast to find because interface error counters and application latency happened to already be on the same Grafana panel; if they'd lived in two separate, unrelated tools, the same investigation could easily have taken hours longer.

## 17. What's Simplified Here

Every scenario above presents a clean, single root cause, reached in a handful of steps — real incidents are frequently messier, involving multiple contributing factors, red herrings, and dead ends before the true cause is found, and the investigations shown here are compressed, idealized versions of what a real multi-hour incident often looks like. IPv6 has its own path MTU discovery mechanism (RFC 8201) that differs in detail from IPv4's (notably, IPv6 routers never fragment packets in flight at all — fragmentation, when it happens, occurs only at the original source), simplified here to the IPv4 case for a single, consistent worked example. Real MSS clamping, PMTUD caching behavior, and DNS resolver failover logic all have considerably more internal detail (timeout values, retry counts, per-OS differences) than shown in these worked examples.

## 18. Interview Questions & Model Answers

**Beginner: A user says a website's DNS resolves fine and `ping` succeeds, but the site itself won't load. What single additional check would you run first, and why?**

*Model answer:* I'd check whether a TCP connection actually completes on the specific port the site uses (`nc -zv <ip> 443`), because `ping`'s ICMP success only confirms the host's kernel is reachable, not that the actual web service is listening or reachable on that specific port — exactly the distinction Scenario 1 and Scenario 4 both hinge on.

**Intermediate: Explain, in your own words, why a network engineer would deliberately want TCP connections to avoid IP fragmentation rather than simply allowing routers to fragment oversized packets as needed.**

*Model answer:* Fragmentation is expensive for routers to perform, increases the chance that losing just one fragment forces the entire original packet to be discarded and retransmitted (since partial reassembly isn't possible), and is often poorly or inconsistently handled by stateful firewalls and security devices. For these reasons, virtually all modern TCP traffic sets the IP header's "Don't Fragment" bit and instead relies on Path MTU Discovery to find the largest size the whole path supports in advance, avoiding fragmentation entirely rather than depending on it working correctly.

**Advanced: A newly established site-to-site VPN causes large file transfers to hang while small requests work fine, and no error appears anywhere in application or system logs. Walk through how you would confirm this is a Path MTU Discovery black-hole, and explain the mechanism precisely.**

*Model answer:* I'd use `ping` with the Don't Fragment flag set (`-M do` on Linux) and progressively smaller packet sizes to find the largest size that can actually cross the path without being silently dropped — if a large size fails silently while a smaller size succeeds, that strongly suggests an MTU mismatch somewhere on the path, most likely introduced by the VPN's encapsulation overhead reducing the effective MTU below the standard 1500 bytes. I'd then check, via a packet capture, whether the sender ever receives an ICMP "Fragmentation Needed" (Type 3, Code 4) message in response to an oversized packet — if it doesn't, some firewall on the path is blocking that specific ICMP message, meaning Path MTU Discovery's feedback loop is broken: the sender keeps sending oversized, "Don't Fragment"-marked packets that keep getting silently discarded, with neither side ever being told why, producing exactly the silent-hang-on-large-data symptom described. The fix is either to permit that ICMP message through every firewall on the path, or to configure TCP MSS clamping at the tunnel endpoints so senders never build oversized segments in the first place.

## 19. Exercises

### Easy
1. In Scenario 1, what single piece of evidence proved the problem was DNS-related rather than a genuine server outage?
2. What is the difference between what `ping` tests and what `nc -zv <host> <port>` tests, and which scenario in this chapter depends most directly on that difference?
3. In your own words, define MTU and explain why different links along one network path can have different MTUs.

### Medium
4. In Scenario 6, explain specifically why having interface error counters and database query latency on the same time-correlated dashboard was faster than investigating each metric in a separate tool.
5. Explain why Scenario 7's symptom ("works internally, fails externally") is one of the most immediately diagnostic patterns in this entire chapter, referencing Chapter 122's client/server/network isolation table.
6. A `ping -M do -s 1472` to a destination fails, but `ping -M do -s 1400` to the same destination succeeds. Using Section 11, explain what this tells you about the path's MTU, and what you would check next to find out why PMTUD isn't automatically handling this.

### Hard
7. Design an investigation plan for a new symptom not covered by this chapter's nine scenarios: "uploads to our file storage service are fast for the first 10 seconds, then slow to a crawl for the rest of the transfer, every time." Using tools and reasoning patterns from this chapter and Chapters 119-121, propose at least two plausible root causes and describe how you would distinguish between them.
8. Scenario 9 concluded that blocking ICMP Type 3 Code 4 is a common cause of PMTUD black-holing. Propose a network configuration change that would prevent this specific failure mode without requiring every firewall operator on the public internet to change their ICMP policy, and explain why it works.
9. Combine two of this chapter's scenarios into one incident: a certificate renewal (Scenario 8) is deployed at the same time a new VPN tunnel (Scenario 9) is added to a path, and both a TLS warning and a large-file-hang symptom appear simultaneously. Using Chapter 122's checklist, explain the order in which you would investigate these two symptoms and why treating them as one root cause would be a mistake.

## 20. Summary and the Bridge to Chapter 124

| Scenario | Root Cause | Key Diagnostic Move |
|---|---|---|
| 1. DNS resolves, HTTP fails | stale cached DNS record | query a second resolver directly |
| 2. TCP connects, app hangs | application-level lock/deadlock | confirm TCP/TLS healthy, then move to app logs |
| 3. Wi-Fi works, mobile fails | broken IPv6 path on carrier network | force `-4`/`-6` to isolate the address family |
| 4. One machine can't reach another | security group silently rejecting TCP | check flow logs' explicit allow/reject field |
| 5. Intermittent packet loss | congested transit link | loss persisting past the hop it starts at |
| 6. Sudden latency spikes | physical-layer CRC errors on a NIC | correlate two metrics on one dashboard |
| 7. Internal works, external fails | missing NAT/port-forward rule | internal-vs-external isolation |
| 8. TLS handshake failures | expired certificate, partial rollout | check `Verify return code` and expiry date directly |
| 9. MTU-caused mysterious failures | black-holed Path MTU Discovery | test with DF set at progressively smaller sizes |

Every scenario in this chapter, and every tool in Chapters 119-122, has operated within one machine, one office, or one data center's boundary — a real but bounded slice of the Internet. The Internet these nine scenarios' packets actually crossed is not one network at all: it's tens of thousands of independently owned and operated networks, held together by business agreements and a single, trust-based routing protocol, spanning oceans via physical cables you've only seen mentioned in passing since Chapter 23. Chapter 124 zooms out to that full scale — revisiting peering, transit, and Tier-1 ISPs not as isolated concepts but as one interconnected global system, the first of Part 19's chapters on how one network became *the* network.
