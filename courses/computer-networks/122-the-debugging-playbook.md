# Chapter 122: The Debugging Playbook — From Symptom to Root Cause

> **"'It doesn't work' is not a starting point for debugging — it's the reason a starting point is needed. This chapter turns the layer model from Chapters 24-26, taught back then as theory, into the literal checklist you run through under pressure."**

---

## Table of Contents

1. [Why This Chapter Exists](#1-why-this-chapter-exists)
2. [The Core Insight: Every Symptom Lives at a Layer](#2-the-core-insight-every-symptom-lives-at-a-layer)
3. [The Layer Checklist, Built From Chapters 24-45](#3-the-layer-checklist-built-from-chapters-24-45)
4. [The Decision Tree](#4-the-decision-tree)
5. [Layer 1-2: Can I Reach It At All Physically/Locally?](#5-layer-1-2-can-i-reach-it-at-all-physicallylocally)
6. [Layer 3: Does Routing Work?](#6-layer-3-does-routing-work)
7. [DNS: A Sideways Step Before Transport](#7-dns-a-sideways-step-before-transport)
8. [Layer 4: Does TCP (or UDP) Connect?](#8-layer-4-does-tcp-or-udp-connect)
9. [TLS: Does the Secure Channel Complete?](#9-tls-does-the-secure-channel-complete)
10. [Layer 7: Does the Application Respond Correctly?](#10-layer-7-does-the-application-respond-correctly)
11. [Binary Search Along the Path](#11-binary-search-along-the-path)
12. [Client vs. Server vs. Network: Isolating the Fault Side](#12-client-vs-server-vs-network-isolating-the-fault-side)
13. [Worked Mini-Example Using the Full Playbook](#13-worked-mini-example-using-the-full-playbook)
14. [Hands-On Experiment](#14-hands-on-experiment)
15. [Common Misconceptions](#15-common-misconceptions)
16. [Production Notes](#16-production-notes)
17. [What's Simplified Here](#17-whats-simplified-here)
18. [Interview Questions & Model Answers](#18-interview-questions--model-answers)
19. [Exercises](#19-exercises)
20. [Summary and the Bridge to Chapter 123](#20-summary-and-the-bridge-to-chapter-123)

---

## 1. Why This Chapter Exists

Chapters 119-121 gave you an arsenal: `tcpdump`, Wireshark, `ping`, `iperf3`, `mtr`, SNMP, flow logs, Prometheus, Grafana. An arsenal without a method is just a pile of tools — faced with a real, vague complaint ("it doesn't work," "it's slow," "it works sometimes"), the actual hard problem isn't knowing what `mtr` does; it's knowing *whether this is even a job for `mtr`* versus `dig`, `curl -v`, or a TLS-specific check. That decision — which tool, in which order, based on what evidence — is what this chapter builds.

The naive approach under pressure is to guess based on whatever you personally suspect first ("it's probably DNS," "it's probably the firewall") and start poking randomly. This sometimes gets lucky, and often wastes real time chasing a hunch while the actual fault sits one layer away, unexamined. The real solution, and the one this chapter formalizes, is much older than any tool in this course: **networking is built in layers (Chapters 24-27) specifically so that each layer can be tested independently** — and a methodical, layer-by-layer sweep from the bottom up (or, as Section 11 shows, via an even faster binary search) will always find the broken layer, because a working layer above a broken one is structurally impossible.

## 2. The Core Insight: Every Symptom Lives at a Layer

Every one of Chapter 123's nine scenarios — and virtually every real networking complaint you will ever encounter — can be restated as: **"communication succeeds up through layer N, and fails at layer N+1."** This reframing is the entire value of the layer model outside of its original teaching purpose in Volume 4: it converts an open-ended, anxiety-inducing question ("why doesn't this work?!") into a bounded, mechanical search problem ("which is the first layer, counting from the bottom, that fails?").

This works precisely *because* of encapsulation (Chapter 27): each layer's data is wrapped inside the layer below it, so a lower layer succeeding is a hard prerequisite for anything above it to even have a chance — an HTTP response (Layer 7) physically cannot arrive if the TCP segment carrying it (Layer 4) was never delivered, which itself cannot happen if the IP packet (Layer 3) never found a route, which itself requires the Ethernet frame (Layer 2) to have left the machine at all (Layer 1). The dependency runs strictly upward, and that strict ordering is what makes a bottom-up sweep guaranteed to terminate at the true fault, not just a plausible-looking one.

## 3. The Layer Checklist, Built From Chapters 24-45

This maps each debugging question directly onto the TCP/IP model (Chapter 26) and names the specific tool from Chapters 119-121 (or earlier, from Chapter 56) that answers it:

| Layer | Question | Primary Tool(s) | Chapters |
|---|---|---|---|
| 1 — Physical | Is the cable/radio link up at all? | `ip addr` (LOWER_UP flag), physical inspection | Ch 14-23 |
| 2 — Data Link | Can I reach this MAC address on my local segment? | `ip neigh`/`arp -a`, switch port status | Ch 28-35, 53 |
| 3 — Network | Do I have a route, and is the destination IP reachable? | `ip route get`, `ping`, `traceroute`/`mtr` | Ch 36-52 |
| — Name resolution | Does the name resolve to the address I expect? | `dig`, `nslookup` | Ch 66-69 |
| 4 — Transport | Does a TCP handshake complete (or UDP traffic flow) on the right port? | `curl -v`, `nc -zv`, `tcpdump`/Wireshark (Ch 119) | Ch 57-65 |
| — Security | Does the TLS handshake complete, and is the certificate valid? | `curl -v`, `openssl s_client` | Ch 77-82 |
| 7 — Application | Does the server return the expected status code and body? | `curl -v`, application logs | Ch 70-76 |

The row order is the order to check things in — not because lower layers fail more often in absolute terms (in real fleets, application-layer bugs are extremely common), but because **checking in dependency order is the only way to trust a negative result at a higher layer.** There is no point debugging why an HTTP response looks wrong if the TCP connection carrying it never actually completed in the first place.

## 4. The Decision Tree

```
                         START: "It doesn't work" (or "it's slow" / "it's flaky")
                                          |
                                          v
                    Can I reach the destination at Layer 2/3 at all?
                    (ip route get, ping the IP directly — Section 5-6)
                          |                                |
                         NO                               YES
                          |                                |
              L1/L2/L3 problem:                 Does the hostname resolve
              cable/Wi-Fi/ARP/routing/          to the IP I expect?
              firewall drop (Ch 122 §5-6,       (dig — Section 7)
              Ch 123 Scenarios 3, 4, 5)              |          |
                                                     NO         YES
                                                      |          |
                                          DNS problem: wrong record,   Does a TCP (or UDP)
                                          stale cache, resolver issue  connection complete
                                          (Ch 123 Scenario 1)          on the right port?
                                                                       (curl -v, nc -zv — Section 8)
                                                                            |         |
                                                                           NO        YES
                                                                            |         |
                                                              L3/L4 problem:    Is this HTTPS? Does
                                                              firewall, SG      the TLS handshake
                                                              rule, service     complete?
                                                              not listening     (Section 9)
                                                              (Ch 123 Scenario  |        |
                                                              4, 7, 9)         NO       YES
                                                                                |        |
                                                                    TLS problem:   Does the app
                                                                    cert, cipher,  respond correctly
                                                                    SNI, expiry    at Layer 7?
                                                                    (Ch 123        (curl -v body/
                                                                    Scenario 8)    status — Section 10)
                                                                                    |        |
                                                                                   NO        YES
                                                                                    |         |
                                                                        App-layer problem:  Working as
                                                                        logic bug, timeout,  designed —
                                                                        hang, wrong config   symptom is
                                                                        (Ch 123 Scenario 2)  elsewhere
                                                                                             (perf, jitter,
                                                                                              Ch 120)
```

Read top to bottom, each "NO" branch is a stopping point — a specific, named class of problem with its own tool and its own worked example in Chapter 123. Each "YES" branch means "this layer is healthy, keep climbing." The tree terminates either at a clear fault layer, or — the rightmost "YES" path all the way down — at "everything up through the application layer is technically working," which pushes the investigation toward Chapter 120's territory (is it *slow*, rather than *broken* — jitter, throughput, or intermittent loss, none of which this pass/fail tree by itself detects).

## 5. Layer 1-2: Can I Reach It At All Physically/Locally?

**The check:** is there a link at all, and if the destination is on your local segment, does ARP (Chapter 53) resolve to a real MAC address?

```
$ ip addr show eth0 | grep -o 'LOWER_UP'
LOWER_UP                             # good — physical link is up

$ ip neigh show 192.168.1.1
192.168.1.1 dev eth0 lladdr b8:27:eb:12:34:56 REACHABLE   # good
192.168.1.1 dev eth0 lladdr (incomplete)                   # bad — ARP never got a reply
```

An `(incomplete)` ARP entry for a same-subnet destination, exactly as Chapter 56, Section 4 flagged, means the local segment itself has a problem — the target device is off, its NIC is down, or something on the LAN (a misconfigured switch port, a VLAN mismatch from Chapter 32) is preventing the ARP reply from ever arriving. No amount of investigation at higher layers will resolve a problem sitting here — this is Chapter 123's Scenario 4 territory when the unreachable machine is on the same LAN.

## 6. Layer 3: Does Routing Work?

**The check:** does my machine have a route to the destination, and does the destination actually respond to that route?

```
$ ip route get 93.184.216.34
93.184.216.34 via 192.168.1.1 dev eth0 src 192.168.1.132   # a route exists

$ ping -c 3 93.184.216.34
3 packets transmitted, 3 received, 0% packet loss           # good — reachable

$ mtr -n -c 20 --report 93.184.216.34
                                                              # if ping fails, mtr shows *where*
                                                              # along the path it stops (Ch 120 §9)
```

A route existing (Section 6's first command) is necessary but not sufficient — Chapter 45's forwarding algorithm only guarantees your machine *knows which way to send the packet*, not that every subsequent hop along that path will actually deliver it. If `ping` to the IP address itself fails (note: to the *IP*, deliberately bypassing DNS for this specific check — see Section 7 for why that separation matters), `mtr` is the next tool, precisely because it localizes the failure to a specific hop rather than leaving "somewhere along the path" as the entire answer.

## 7. DNS: A Sideways Step Before Transport

DNS doesn't fit cleanly into a numbered OSI/TCP-IP layer — it's an application-layer protocol (it runs over UDP/TCP port 53, itself sitting on top of IP) that nonetheless functions as a *prerequisite* for almost everything above Layer 3 in practice, since nearly every real connection starts with a hostname, not a bare IP. This is precisely why it earns its own step in the checklist rather than being folded silently into either Layer 3 or Layer 7 — and why Section 6's `ping`/`mtr` checks deliberately targeted a raw IP address first: **that separation is the single most valuable diagnostic move in this entire chapter**, because it cleanly answers "is this a DNS problem or a network-path problem" with one substitution.

```
$ dig example.com +short
93.184.216.34                       # confirm this matches what you expect

$ dig example.com @8.8.8.8 +short
93.184.216.34                       # confirm a different resolver agrees (rules out local resolver/cache issue)
```

If `ping <hostname>` fails but `ping <the IP dig just returned>` succeeds, the fault is entirely in name resolution (Chapter 123's Scenario 1 territory), and every layer this chapter would otherwise check above Layer 3 is, at this point, provably irrelevant to the symptom — a fact worth confirming early precisely because it rules out several later checklist steps in one move.

## 8. Layer 4: Does TCP (or UDP) Connect?

**The check:** does a TCP three-way handshake (Chapter 59) actually complete on the specific port the application needs, as opposed to the network merely being reachable in general (Section 6 only proved ICMP works, which is answered by the kernel, not by the target application — Chapter 56, Section 5's misconception callout applies directly here).

```
$ nc -zv 93.184.216.34 443
Connection to 93.184.216.34 443 port [tcp/https] succeeded!

$ nc -zv 93.184.216.34 8080
nc: connect to 93.184.216.34 port 8080 (tcp) failed: Connection refused
```

`nc -zv` (`netcat`, "zero-I/O mode, verbose") attempts exactly a TCP handshake and nothing more — the single most surgical Layer 4 check, isolating "can I connect" from everything DNS, TLS, or HTTP might otherwise add to the picture. **"Connection refused" and "connection timed out" are meaningfully different failures, and this checklist treats them differently:**

- **Refused** means a RST came back (Chapter 64) — something *is* reachable at that IP and actively told you no port is listening there, or a stateful firewall (Chapter 84) explicitly rejected it. The network path itself is fine; the problem is at or very near the destination.
- **Timed out** (silence, no RST at all) usually means a packet filter dropped the SYN silently somewhere along the path (Chapter 84's stateless filtering, or a cloud security group — Chapter 97) rather than any host actively refusing it, since a genuinely unreachable or down host would typically be answered by an ICMP "destination unreachable" rather than dead silence in most real network configurations.

This single distinction is frequently the fastest way to tell "wrong port/service not running" (refused) apart from "something is silently blocking this traffic" (timeout) — Chapter 123's Scenarios 4, 7, and 9 each hinge on correctly reading exactly this signal.

## 9. TLS: Does the Secure Channel Complete?

**The check:** for any HTTPS/TLS-protected service, a successful TCP connection is necessary but says nothing about whether the encrypted session on top of it can actually be established — a distinct step with its own distinct failure modes (Chapter 82).

```
$ openssl s_client -connect example.com:443 -servername example.com </dev/null 2>&1 | grep -E 'Verify return|subject'
subject=CN=example.com
Verify return code: 0 (ok)
```

`Verify return code: 0 (ok)` confirms both that the handshake completed *and* that the certificate chain validated (Chapter 81) — a non-zero code here (common values: `10` expired certificate, `18`/`19` self-signed or unknown CA, `62` hostname mismatch) pinpoints exactly which trust check failed, without needing to guess. `curl -v`'s TLS section (Chapter 56, Section 10; Chapter 119, Section 9) gives largely the same information in a more familiar format if `openssl s_client`'s raw output feels unwieldy. This step is where Chapter 123's Scenario 8 lives entirely.

## 10. Layer 7: Does the Application Respond Correctly?

**The check:** with TCP (and TLS, if applicable) confirmed working, is the actual application producing the right response — the right status code, the right body, in a reasonable amount of time?

```
$ curl -v -o /dev/null -w 'HTTP %{http_code} in %{time_total}s\n' https://example.com/api/status
HTTP 200 in 0.083s          # healthy

$ curl -v -o /dev/null -w 'HTTP %{http_code} in %{time_total}s\n' https://example.com/api/status
HTTP 502 in 2.104s          # the connection/TLS worked; the app or something in front of it failed
```

A `502 Bad Gateway` here is a genuinely informative result *precisely because* it can only occur once TCP and TLS have already succeeded — it's a proxy or load balancer (Chapter 95) explicitly reporting that it couldn't get a good response from a backend, which immediately narrows the search to "backend health" rather than leaving "network or application" as an open question. A hang with **no** response at all, by contrast (TCP and TLS both completed but the connection simply sits open with nothing coming back), is Chapter 123's Scenario 2 exactly — and is diagnosed by looking at the *application's own* state (is it blocked waiting on a database, a lock, another downstream service) rather than anything this networking course's tools alone can resolve, since at that point the network layers have already done their job correctly.

## 11. Binary Search Along the Path

Sections 5-10 describe a strict bottom-up sweep, which is the safest default when you have no prior evidence pointing anywhere specific. But when the path between client and server is long — client, home router, ISP, several transit hops, a cloud load balancer, an application server, a database — checking every layer at every hop is slow. **The faster method, once you have at least one working reference point, is binary search: pick a point roughly in the middle of the suspected path and test from there, then repeatedly halve the remaining distance based on the result** — precisely the same algorithmic idea Chapter 48's link-state routing used for finding shortest paths, applied here to finding a *fault* instead of a route.

```
Client ──── Router ──── ISP ──── (thousands of km) ──── Cloud LB ──── App Server ──── Database
  |                                                          |
  can't reach the site                              SSH into the cloud LB directly
                                                      and curl the app server from there
                                                              |
                                                    it works fine from here
                                                              |
                                        conclusion: the fault is somewhere between the client
                                        and the cloud LB, not between the LB and the backend —
                                        half the path is now eliminated in one test
```

This is exactly why cloud providers and CDNs (Chapter 96) give customers the ability to SSH into intermediate infrastructure, or provide synthetic checks from multiple points along a path — each additional vantage point is another opportunity to bisect the remaining search space instead of walking it hop by hop.

## 12. Client vs. Server vs. Network: Isolating the Fault Side

A question that cuts across every layer in the checklist and deserves to be asked explicitly, early: **is this symptom specific to one client, one server, or does it affect the connection between many different pairs?**

| Observation | Points toward |
|---|---|
| Fails from this one machine, works from every other machine on the same network | client-local (Chapter 123 Scenario 3's territory — device-specific config, not the network path) |
| Fails from every machine on this network, works from elsewhere | network-path or local-network problem (Chapter 123 Scenario 3, 5) |
| Fails for everyone, everywhere | server-side or DNS-wide problem (Chapter 123 Scenario 1, 7) |
| Fails only for one specific destination, all others fine | that specific destination's server, or a route/policy specific to reaching it |

This single comparison — "does the same test succeed from a different vantage point" — is often faster than a full layer-by-layer sweep, because it immediately eliminates entire branches of Section 4's decision tree: if a second machine on the exact same LAN, hitting the exact same destination, succeeds where the first fails, then Layers 3 and above (routing, DNS, the destination server) are already proven fine, and the entire remaining search space collapses to "what's different about this one client" — Layer 1/2 or local configuration only.

## 13. Worked Mini-Example Using the Full Playbook

Symptom: **"I can't reach `internal-api.company.com` from my laptop."**

```
1. Ping the hostname directly:
   $ ping internal-api.company.com
   ping: cannot resolve internal-api.company.com: Unknown host
   -> Immediately stop the Layer 1-4 checklist: this is a DNS problem (Section 7),
      not a routing or transport problem at all.

2. Check DNS directly:
   $ dig internal-api.company.com
   ;; connection timed out; no servers could be reached
   -> Not a "record doesn't exist" answer — the DNS query itself never got a reply.
      This demotes the investigation one more level: is DNS resolution itself
      even reachable?

3. Check whether the configured resolver is reachable at all:
   $ cat /etc/resolv.conf
   nameserver 10.50.0.2
   $ ping -c 3 10.50.0.2
   100% packet loss
   -> The internal DNS resolver itself is unreachable — likely because this is
      a corporate-internal name requiring a VPN (Chapter 85) that isn't currently
      connected.

4. Confirm the VPN theory directly:
   $ ip route show | grep tun
   (no output — no VPN interface present)

Conclusion: not a DNS misconfiguration, not a broken destination server, not a
routing problem at the destination's end at all — the VPN simply isn't
connected, so the internal resolver (and the internal network generally) was
never reachable in the first place. Every step above ruled out one layer of
the checklist in order, arriving at the true cause in four short commands
instead of guessing.
```

## 14. Hands-On Experiment

Pick any real website and deliberately walk Sections 5-10's checklist against it, writing down the actual output at each step, even though you expect it to succeed:

```bash
ip route get $(dig +short example.com | head -1)   # Layer 3: is there a route
ping -c 3 $(dig +short example.com | head -1)       # Layer 3: is it reachable
dig example.com +short                              # DNS: does it resolve
nc -zv $(dig +short example.com | head -1) 443      # Layer 4: does TCP connect
openssl s_client -connect example.com:443 -servername example.com </dev/null 2>&1 | grep 'Verify return'  # TLS
curl -o /dev/null -s -w 'HTTP %{http_code}\n' https://example.com   # Layer 7
```

Then deliberately break one layer yourself (e.g., add a bogus entry to `/etc/hosts` to break DNS resolution for one hostname, or block port 443 outbound with a local firewall rule) and re-run the same sequence — watching exactly which step's output changes, and only that step, is the fastest way to internalize that the checklist genuinely isolates faults rather than just listing commands.

## 15. Common Misconceptions

- **"If `ping` works, the network is fine and the problem must be the application."** Section 8 directly contradicts this — `ping`'s ICMP is answered by the destination's kernel, not by whatever application or port you actually care about; a completely healthy ping alongside a completely broken TCP connection on a specific port is one of the most common patterns in this entire chapter.
- **"You should always start troubleshooting at the application layer, since that's usually where the bug is."** In terms of raw frequency across a large fleet, this is often true — but Section 2's dependency argument still holds: you can only *trust* an application-layer diagnosis once you've confirmed the layers beneath it aren't secretly the real cause, which is why the checklist orders by dependency, not by statistical likelihood.
- **"'Connection refused' and 'connection timed out' are basically the same failure, just different wording."** Section 8 draws a precise, actionable distinction between them — refused means something answered and said no; timed out usually means something silently dropped the attempt — and conflating them routinely sends people chasing the wrong fix (checking a firewall rule when the real issue is a crashed service, or vice versa).
- **"A decision tree like Section 4's guarantees you'll always find the right answer quickly."** It guarantees you'll never skip a layer and never trust an unproven assumption, which eliminates a huge class of wasted effort — but Chapter 123's scenarios include several (intermittent loss, sudden latency spikes) where the tree correctly routes you to "keep climbing, everything technically works," and the real answer requires Chapter 120's measurement tools, not a pass/fail check.

## 16. Production Notes

- **Runbooks in real organizations are, structurally, exactly this chapter's decision tree, specialized to specific services.** A well-run on-call rotation typically has a written runbook per critical service that is this playbook with the generic steps replaced by that service's specific commands, dashboards (Chapter 121), and known failure signatures — this chapter is the template, not a replacement for building one.
- **Automating the checklist is standard practice for synthetic monitoring.** Chapter 121's Prometheus/Grafana stack commonly runs scripted versions of Sections 5-10's exact checks on a fixed schedule from multiple locations, alerting on the first layer that fails rather than waiting for a human to run through the checklist manually during an incident.
- **The binary-search approach (Section 11) is the entire justification for maintaining "jump boxes" or bastion hosts at multiple points in a network topology** — without a reachable vantage point partway along a path, binary search degrades back into a bottom-up sweep from one end only.
- **Isolating client vs. server vs. network (Section 12) should be the very first question in any incident involving more than one affected user**, because it can eliminate half the checklist's later steps in a single comparison, and many incident-response processes formalize this as literally the first triage question asked.

## 17. What's Simplified Here

Real systems frequently have multiple layers of proxies, load balancers, and service meshes (Chapter 101) between a client and the "true" backend, each of which can independently fail at any of Sections 5-10's checks — a real investigation often has to run this entire checklist multiple times, once per hop of infrastructure, not just once for the whole path. This chapter also treats DNS as a single step, when in practice a resolver's own caching behavior (Chapter 68) can make "does it resolve" a moving target that changes between two consecutive checks. Finally, the decision tree in Section 4 presents clean binary branches for teaching clarity; real symptoms are sometimes genuinely ambiguous between two branches (partial packet loss looks like both "sometimes reachable" and "sometimes not"), requiring the judgment Chapter 120 built around interpreting noisy, statistical measurements rather than a single pass/fail check.

## 18. Interview Questions & Model Answers

**Beginner: Why does it make sense to check DNS and basic IP reachability before checking whether the application is responding correctly?**

*Model answer:* Because of encapsulation (Chapter 27) and dependency order — an HTTP response can't arrive if TCP never connected, which can't happen if the IP address wasn't reachable, which is often preceded by needing DNS to even know what IP to try. Checking a higher layer before confirming the layers beneath it means any result you get could be misleading — you can't tell whether the application is actually broken or whether the request never even reached it.

**Intermediate: A user reports a TCP connection to your service "times out." Another user reports theirs is "refused." Are these describing the same underlying problem? Explain.**

*Model answer:* No — these are meaningfully different signals (Section 8). "Refused" means a RST came back, meaning something at that address is reachable and actively rejecting the connection (commonly, nothing is listening on that port, or a stateful firewall explicitly rejected it). "Timed out" means no response came back at all, which usually indicates a packet filter or firewall silently dropping the SYN somewhere along the path, or a genuinely unreachable host. I'd investigate these as two different hypotheses, not treat them as the same complaint.

**Advanced: Describe how you would use the binary-search approach from Section 11, rather than a full bottom-up layer sweep, to diagnose a slow response affecting users of a service that sits behind a CDN, a load balancer, and an application server, when you have SSH access to the load balancer and the application server but not to any client machine.**

*Model answer:* I'd start by testing directly from the load balancer to the application server (`curl -v` with timing, Chapter 56 Section 12's format) — a point roughly in the middle of the full path. If that leg is fast and healthy, the problem is eliminated from that half of the path entirely, and I'd focus on the CDN-to-load-balancer leg, or the client-to-CDN leg, next — likely using the CDN provider's own edge diagnostics or synthetic monitoring from external vantage points (Chapter 121) since I lack direct client access. If the load-balancer-to-app-server leg is instead where the slowness appears, I've immediately eliminated the CDN and load balancer as likely causes and can focus investigation on the application server and whatever it depends on (database, downstream services) without wasting time re-checking the parts of the path already proven healthy.

## 19. Exercises

### Easy
1. Using Section 3's table, list the order in which you would check DNS, TCP connectivity, and TLS for an HTTPS connection failure, and explain why that order matters.
2. What is the practical difference between "connection refused" and "connection timed out" when using `nc -zv`?
3. Name one command from this chapter that tests Layer 3 reachability using a raw IP address specifically to rule DNS in or out of the investigation.

### Medium
4. A colleague says "the ping works, so the network must be fine — the bug has to be in the application." Explain precisely why this reasoning is incomplete, citing a specific section of this chapter.
5. Using Section 12's client/server/network isolation table, describe what conclusion you would draw if a service is unreachable from every machine in one specific office, but reachable from every other office and from the internet at large.
6. Walk through Section 4's decision tree for a symptom where DNS resolves correctly, a TCP connection completes, but no TLS ServerHello is ever received before the connection times out. At which decision point does this symptom diverge from the "everything healthy" path, and what would you check next?

### Hard
7. Design a specific, ordered checklist (using this chapter's sections and tools) for diagnosing "our service works when accessed via IP address directly, but not via its hostname" — identify exactly which layer this symptom pattern isolates the fault to, and explain your reasoning.
8. Using the binary-search idea from Section 11, describe a general strategy for finding the single faulty hop in an 8-hop network path when you have shell access to only the first and last hosts (not any of the 6 hops in between), and explain what information you would need from `mtr` (Chapter 120) to compensate for the lack of direct access to the middle of the path.
9. A service fails intermittently — roughly 1 request in 20 times out, and the other 19 succeed instantly with no observable pattern. Explain why Section 4's binary pass/fail decision tree is insufficient on its own for this symptom, and describe what change to the investigation approach (drawing on Chapter 120) would be necessary to make progress.

## 20. Summary and the Bridge to Chapter 123

| Step | Question | Chapter 123 Scenario(s) It Diagnoses |
|---|---|---|
| Layer 1-2 | Physical/local link and ARP | Scenario 3, 4 |
| Layer 3 | Routing and raw IP reachability | Scenario 4, 5, 6 |
| DNS | Does the name resolve correctly | Scenario 1 |
| Layer 4 | Does TCP/UDP connect on the right port | Scenario 4, 7, 9 |
| TLS | Does the secure handshake complete | Scenario 8 |
| Layer 7 | Does the application respond correctly | Scenario 2, 7 |
| Binary search / isolation | Which segment or side of the fault | all nine, as a speed-up |

You now have both the tools (Chapters 119-121) and the method (this chapter) — everything needed to turn a vague complaint into a confirmed root cause, methodically, without guessing. Chapter 123 is where that method gets exercised for real: nine fully worked, realistic scenarios, each one run through exactly this checklist from first symptom to final root cause, including the ninth — MTU and fragmentation failures — which requires one more piece of mechanism this course hasn't covered yet, introduced there from first principles before it's debugged.
