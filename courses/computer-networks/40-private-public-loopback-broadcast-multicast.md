# Chapter 40: Private vs. Public Addresses, Loopback, Broadcast, Multicast

> *"Not every address on your network is meant to reach the whole world — some are meant to reach nobody but yourself, some are meant to reach everybody nearby, and some are meant to reach a chosen few. IP had to invent all three."*

---

## Table of Contents

1. [The Problem: Every Address Cannot Mean "Reachable From Everywhere"](#1-the-problem-every-address-cannot-mean-reachable-from-everywhere)
2. [The Naive Assumption and Why It Breaks](#2-the-naive-assumption-and-why-it-breaks)
3. [RFC 1918 — The Private Address Ranges](#3-rfc-1918--the-private-address-ranges)
4. [What "Private" Actually Means on the Wire](#4-what-private-actually-means-on-the-wire)
5. [Public Addresses — The Ones the Whole Internet Agrees On](#5-public-addresses--the-ones-the-whole-internet-agrees-on)
6. [Loopback — Talking to Yourself](#6-loopback--talking-to-yourself)
7. [Subnet Broadcast — Reaching Everyone on This Wire](#7-subnet-broadcast--reaching-everyone-on-this-wire)
8. [Multicast — Reaching a Chosen Group](#8-multicast--reaching-a-chosen-group)
9. [Other Special-Use Ranges You Will Meet](#9-other-special-use-ranges-you-will-meet)
10. [Full Worked Example: Classifying a Batch of Addresses](#10-full-worked-example-classifying-a-batch-of-addresses)
11. [Hands-On Experiment](#11-hands-on-experiment)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Production Notes](#13-production-notes)
14. [What's Simplified Here](#14-whats-simplified-here)
15. [Interview Questions & Model Answers](#15-interview-questions--model-answers)
16. [Exercises](#16-exercises)
17. [Summary](#17-summary)

---

## 1. The Problem: Every Address Cannot Mean "Reachable From Everywhere"

Chapters 36 through 39 built up IPv4 addressing as if every address were interchangeable: any 32-bit number, sliced into network and host portions, subnetted and aggregated as needed. That picture is accurate as far as it goes, but it hides three practical problems that show up the moment you try to actually deploy IP on real networks.

**Problem one — there are not enough addresses to go around.** IPv4 has exactly 2^32 = 4,294,967,296 possible addresses. That sounds like a lot until you remember that a single household today might have a phone, a laptop, a smart TV, a thermostat, two tablets, and a game console all wanting an address simultaneously — and there are billions of households. If every device on Earth needed a unique, globally-routable address, IPv4 would have run out decades before it did (and, as Chapter 42 will show with real dates, it *has* run out).

**Problem two — most devices never need to be reached from the outside world.** Your laptop, sitting on your home Wi-Fi, needs to reach Google. Google does not need to reach your laptop unprompted. Neither does a random stranger's laptop in another country. The overwhelming majority of devices only ever *initiate* connections outward; they don't need a globally unique, globally reachable identity at all.

**Problem three — sometimes a device needs to talk to itself, or to everyone nearby, or to a specific group — not to "the internet."** A database server testing its own network stack needs to send a packet that never leaves the machine. A DHCP client (Chapter 55) that doesn't have an address *yet* still needs to ask "is there a DHCP server anywhere on this wire?" A video conferencing backend streaming the same feed to 10,000 subscribers shouldn't have to send 10,000 separate copies of every packet.

These three problems have three different answers, and this chapter covers all of them: **private addressing** (for problem one and two together), **loopback** (for talking to yourself), **broadcast** (for reaching everyone on a local wire), and **multicast** (for reaching a defined group). Chapter 41 then covers NAT, which is the mechanism that makes private addressing actually work with the public internet.

It's worth previewing, at a glance, how these fit together conceptually before diving into each one — they answer four genuinely different questions about "who should receive this packet":

```
"Nobody but me"                -> Loopback     (127.0.0.0/8 / ::1)
"Anyone, anywhere, uniquely"   -> Public address
"Reusable, not globally unique"-> Private address (RFC 1918)
"Literally everyone on my LAN" -> Broadcast     (subnet-directed or 255.255.255.255)
"A specific, opted-in group"   -> Multicast     (224.0.0.0/4)
```

Every address you will ever encounter on an IPv4 network falls into exactly one of these five categories — which is precisely what makes the classification exercise in Section 10 a complete, well-defined exercise rather than a matter of judgment.

---

## 2. The Naive Assumption and Why It Breaks

Imagine you're designing IP addressing for a new company network in the 1980s, before any of this had been standardized. The naive approach: give every device — every desktop, every printer, every internal test server — its own address drawn from the same global pool that public internet hosts use.

This naive approach works fine on paper for a company with 50 machines. It fails in three ways once it meets reality:

1. **Waste.** A company with 5,000 internal machines, most of which will never accept an inbound connection from the internet, has just consumed 5,000 globally unique addresses that could have gone to someone who actually needs global reachability.
2. **No renumbering pressure relief.** If the company ever needs to change ISPs, every single one of those 5,000 machines potentially needs a new address, because the address was tied to a specific place in the global routing hierarchy.
3. **No isolation.** Every one of those 5,000 machines is, in principle, directly reachable (and attackable) from anywhere on Earth, whether or not that was ever the intent.

The fix that emerged — formalized in **RFC 1918** (*Address Allocation for Private Internets*, February 1996, obsoleting RFC 1597) — is almost embarrassingly simple: carve out a few blocks of address space that **everyone is allowed to reuse**, and agree, by convention enforced at the routing level, that **no router on the public internet will ever forward a packet to or from those blocks**. Every home network, every office network, every cloud VPC (Chapter 97) on Earth can use `192.168.1.1` for its router, and none of them will ever collide with each other on the public internet — because packets addressed to `192.168.1.1` never leave the private network in the first place.

---

## 3. RFC 1918 — The Private Address Ranges

RFC 1918 reserves exactly three blocks of IPv4 address space for private use:

| Range | CIDR | Address Count | Typical Use |
|---|---|---|---|
| `10.0.0.0` – `10.255.255.255` | `10.0.0.0/8` | 16,777,216 | Large enterprises, cloud VPCs, ISPs' internal infrastructure |
| `172.16.0.0` – `172.31.255.255` | `172.16.0.0/12` | 1,048,576 | Medium enterprises, Docker's default bridge network |
| `192.168.0.0` – `192.168.255.255` | `192.168.0.0/16` | 65,536 | Home routers, small offices |

Notice the sizes are deliberately different — `10.0.0.0/8` is enormous (a full Class A equivalent), `172.16.0.0/12` is a middle ground carved out of what used to be Class B space, and `192.168.0.0/16` is a single Class B-sized block. This wasn't arbitrary: RFC 1918 was written to give organizations of every size a block that fits, without wasting more than necessary — the same philosophy Chapter 39 covered for CIDR aggregation, applied to address *reservation* rather than *allocation*.

**How the "no router will forward it" rule is enforced.** There's no cryptographic or protocol-level lock on these addresses — nothing stops a misconfigured router from trying to forward a packet destined for `10.5.5.5` out to the internet. What actually enforces the boundary is operational convention: every major ISP and every Tier-1 network (Chapter 51) is configured to silently drop, at its network edge, any packet whose source or destination falls inside these three ranges. This is sometimes called **bogon filtering** (a "bogon" being a packet that has no business appearing on the public internet — RFC 1918 space is one category, but so are unallocated IANA blocks). The system works because essentially the entire internet cooperates on this one rule.

**Why three separate blocks instead of one bigger one?** Politically and historically, `10.0.0.0/8` was chosen because it was the old (pre-CIDR) Class A network 10, already informally used for exactly this purpose inside large private networks. `192.168.0.0/16` mirrors old Class B usage. `172.16.0.0/12` fills the gap. All three predate RFC 1918's formal blessing — the RFC mostly standardized existing informal practice, which is a pattern you'll see again and again in this course (Chapter 39's CIDR did the same for classful addressing's mess).

**How this connects to the public allocation hierarchy.** Chapter 39 introduced the idea that public address blocks flow down a chain: IANA allocates large blocks to Regional Internet Registries (RIRs — ARIN for North America, RIPE NCC for Europe/Middle East, APNIC for Asia-Pacific, LACNIC for Latin America, AFRINIC for Africa), which allocate smaller blocks to ISPs, which assign individual addresses to customers. RFC 1918 space deliberately sits *outside* this entire hierarchy — no RIR, no ISP, and no IANA record ever needs to be consulted to use `192.168.1.1`, which is exactly what makes it free and immediately usable by anyone, anywhere, with zero paperwork. This is the direct structural reason private addressing solves the *scarcity* half of Chapter 41's opening problem, even though it creates the *reachability* half in the same stroke.

---

## 4. What "Private" Actually Means on the Wire

It's worth being precise about what "private" does and does not mean, because the word invites two wrong intuitions.

**It does not mean encrypted.** A packet from `192.168.1.50` to `192.168.1.1` travels in exactly the same plaintext Ethernet frames (Chapter 28) as any other packet. "Private" here is about *routability*, not confidentiality. Anyone on the same local network can still see that traffic with a packet sniffer (Chapter 119).

**It does not mean unique.** Quite the opposite — `192.168.1.1` is almost certainly the router address in *your* house, in your neighbor's house, and in ten million other houses simultaneously. This is only safe because none of those `192.168.1.1`s are ever visible to each other; each one lives inside its own isolated private network, walled off by NAT (Chapter 41) or a private WAN.

```
Your house                          Your neighbor's house
┌─────────────────────┐             ┌─────────────────────┐
│  laptop  192.168.1.50│             │  laptop  192.168.1.77│
│  phone   192.168.1.51│             │  TV      192.168.1.20│
│  router  192.168.1.1 │             │  router  192.168.1.1 │
└──────────┬───────────┘             └──────────┬───────────┘
           │ (public IP: 203.0.113.9)           │ (public IP: 198.51.100.44)
           └───────────────┬─────────────────────┘
                            │
                     Public Internet
```

Both routers are `192.168.1.1`. Both networks work perfectly. Neither network can see or address the other's internal machines directly — only their public-facing router addresses are visible to the outside world, and even those are reached indirectly through NAT, which Chapter 41 dissects mechanically.

**A wrinkle worth naming here even though Chapter 41 owns the full mechanism: double NAT.** It's increasingly common — especially with ISP-provided modem/router combos — for a household to unknowingly run *two* layers of private addressing stacked on top of each other: the ISP's own equipment hands out an address from `192.168.1.0/24` (or sometimes `100.64.0.0/10`, the CGNAT range from Section 9) to the customer's own router, which then runs its *own* separate `192.168.1.0/24` (or `192.168.0.0/24`) network for the customer's actual devices. From inside the house, this is invisible — your laptop still just sees `192.168.1.50` — but it means two full layers of address translation happen before a packet reaches the real internet, which is a common, very real cause of certain applications (game consoles wanting inbound connections, in particular) failing to work correctly even after the customer has "correctly" configured port forwarding on their own router, because that forwarding rule never reaches the ISP's outer layer at all.

---

## 5. Public Addresses — The Ones the Whole Internet Agrees On

A **public IP address** is, by contrast, any IPv4 address that is *not* reserved for private use, loopback, multicast, or one of the other special-use categories in this chapter — and that has been formally allocated by IANA (the Internet Assigned Numbers Authority) to a Regional Internet Registry (RIR), which in turn allocated it to an ISP or organization, which assigned it to a specific device or router.

The defining property of a public address is that it is meant to be **globally unique and globally routable** — every router on the internet either knows how to forward toward it, or knows to forward toward a router that does (this is the entire subject of Part 7: Routing, starting at Chapter 44). If your home router has public address `203.0.113.9`, that address means *exactly one thing* on the entire planet, and BGP (Chapter 49) is the machinery that makes sure a packet from Tokyo and a packet from São Paulo both find their way to it.

This is precisely why public addresses are scarce and private addresses are free: uniqueness at global scale is the expensive property. Anyone can hand out `192.168.1.x` addresses because nobody has to coordinate; only one entity on Earth can ever hold `203.0.113.9` at a time, and that requires the entire allocation hierarchy (IANA → RIR → ISP → you) that Chapter 39 introduced when discussing CIDR-based allocation.

| | Private Address | Public Address |
|---|---|---|
| Uniqueness | Only within its own private network | Globally unique |
| Routable on the internet | No — filtered at ISP edges | Yes |
| Who allocates it | Anyone, for free, from RFC 1918 ranges | IANA → RIR → ISP → you |
| Example | `192.168.1.50` | `142.250.80.46` (a Google address) |
| Reachable directly from outside? | No — needs NAT (Ch. 41) | Yes, if not firewalled |

**Seeing the allocation hierarchy directly.** Every public address's allocation history is a matter of public record, queryable with the `whois` tool:

```bash
$ whois 8.8.8.8 | head -20
NetRange:       8.8.8.0 - 8.8.8.255
CIDR:           8.8.8.0/24
NetName:        LVLT-GOGL-8-8-8
NetHandle:      NET-8-8-8-0-1
Organization:   Google LLC (GOGL)
RegDate:        2014-03-14
Updated:        2014-03-14
```

This is the same allocation hierarchy Chapter 39 introduced conceptually for CIDR — IANA to RIR to ISP or organization — made concrete: `8.8.8.8` isn't just "a public address," it's an address whose entire chain of custody, back to IANA, is on record and auditable by anyone. No such record exists, or could exist, for `192.168.1.1` — there's nothing to look up, because nobody had to be granted it in the first place.

---

## 6. Loopback — Talking to Yourself

**The problem.** Sometimes a program on a machine needs to talk to another program on the *same* machine, using the network stack — not because the data needs to physically travel anywhere, but because using sockets (Chapter 57) is a convenient, protocol-agnostic way for two processes to communicate, and because you want to test your networking code without needing an actual network. A developer running a local web server wants to hit it with a browser on the same laptop, with zero dependency on Wi-Fi, Ethernet, or any physical medium at all.

**The naive attempt.** You could just use your machine's real IP address (say, `192.168.1.50`) to talk to itself. This technically can work, but it's fragile: it depends on your network interface being up, it depends on your actual assigned address (which might change via DHCP, Chapter 55), and if your network cable is unplugged, "talking to yourself" now depends on hardware that has nothing to do with the conversation.

**The real solution.** IPv4 reserves an entire block, `127.0.0.0/8` (yes, a full /8 — 16.7 million addresses, though by overwhelming convention only `127.0.0.1` is ever actually used), as the **loopback range**. Any packet sent to an address in `127.0.0.0/8` is handled entirely inside the sending machine's own network stack. It never touches a network interface card, never becomes an electrical or optical signal, and never leaves the machine — the operating system's kernel loops the packet directly from the sending socket back to a receiving socket, hence the name.

```
Process A (client)                 Process B (server)
     |                                    |
     |  connect() to 127.0.0.1:8080       |
     |----------------------------------->|  (kernel loopback interface "lo")
     |                                    |
     |  no NIC, no cable, no switch,      |
     |  no electrical signal involved     |
```

`::1` is the IPv6 equivalent (Chapter 42 covers IPv6 addressing in full) — a single reserved address rather than a whole block, since IPv6 has no shortage to worry about.

**Why this matters in practice.** Every web server you've ever run locally (`localhost:3000`, `127.0.0.1:8080`) uses loopback. Databases often bind only to loopback by default specifically so they are *unreachable* from any other machine, private or public — a basic security posture. The word "localhost" is simply the conventional hostname that resolves to `127.0.0.1` (or `::1`) via your machine's hosts file or resolver — Chapter 66 covers hostname resolution in depth.

**A subtlety worth knowing:** because the entire `127.0.0.0/8` block is reserved, `127.0.0.2`, `127.1.2.3`, and even `127.255.255.255` are all also valid loopback addresses on most operating systems — some testing tools deliberately use different loopback addresses (e.g., `127.0.0.2`) to simulate "multiple hosts" on one machine without any real network involved.

**Why this matters for security, not just convenience.** Binding a service to `127.0.0.1` specifically, rather than to a machine's real network-facing address (or the catch-all `0.0.0.0`), is one of the single cheapest security decisions available to a developer: a service on loopback is, by construction, unreachable from any other machine on Earth, private or public, because — as this section has stressed — loopback traffic never reaches a NIC at all. This is why local development databases, internal metrics endpoints, and admin consoles are so often deliberately bound to `127.0.0.1:PORT` rather than left open on every interface: it's not an accident of convenience, it's a genuine access-control boundary enforced at the kernel level, with zero firewall rules required.

---

## 7. Subnet Broadcast — Reaching Everyone on This Wire

**The problem.** Chapter 55 will show that a device joining a network for the first time doesn't have an IP address yet — so it can't send a packet "to" the DHCP server, because it doesn't know the DHCP server's address, and in fact might not even know its *own* address well enough to have a valid source. It needs a way to say, in effect, "to whoever on this wire can help me — I don't know your name yet."

**The naive attempt.** You might imagine sending the request to every address in the subnet one at a time. On a `/24` subnet that's up to 254 individual packets for a single question, wasteful and slow, and the sender doesn't even know which of those addresses currently have a live device behind them.

**The real solution: broadcast.** IPv4 reserves the **all-ones host portion** of any subnet as its **broadcast address** — a single address that means "every host on this subnet, deliver this frame to yourselves." Recall from Chapter 37 that a subnet mask splits an address into a network portion and a host portion; the broadcast address is exactly that network portion followed by all 1-bits in the host portion.

**Worked example.** Take the subnet `192.168.1.0/24`:

```
Network:     192.168.1.0      (host bits = 00000000)
Broadcast:   192.168.1.255    (host bits = 11111111)
Usable hosts: 192.168.1.1 – 192.168.1.254
```

Any device on that subnet sending a packet to `192.168.1.255` expects every other device on the same physical (or switched, per Chapter 30) LAN segment to receive and process it. At the Ethernet layer, this maps directly onto the Ethernet broadcast MAC address `FF:FF:FF:FF:FF:FF` (Chapter 29) — a switch (Chapter 31) floods a frame addressed to that MAC out every port, because a broadcast frame has no "known destination" to look up.

There is also the **limited broadcast address**, `255.255.255.255`, which means "everyone on my local network, whatever subnet I happen to be on" — used precisely in situations like DHCP DISCOVER, where the sending device doesn't yet know its own subnet, so it cannot compute a subnet-specific broadcast address at all.

```
Directed (subnet) broadcast:  192.168.1.255   — "everyone on 192.168.1.0/24"
Limited broadcast:            255.255.255.255 — "everyone on my local wire, whatever subnet that is"
```

**Why broadcast doesn't scale past a LAN.** A broadcast is, by design, never forwarded by a router. If it were, a broadcast storm on any network anywhere on Earth would eventually reach every device on the internet — an obvious catastrophe. Routers are the boundary of a **broadcast domain** (Chapter 30 defined this precisely): broadcast traffic dies at the first router hop. This is one of the fundamental reasons IP needs *routing* at all rather than "broadcast to everyone and let them ignore what's not theirs" — it would never scale past a single LAN, a point Chapter 3's N² wiring discussion foreshadowed from the opposite direction.

---

## 8. Multicast — Reaching a Chosen Group

**The problem.** Broadcast is "everyone." Unicast (an ordinary IP address, everything you've seen so far) is "exactly one recipient." But plenty of real traffic is neither: a live video stream has many interested subscribers but is certainly not meant for every device on the LAN, let alone the internet. Routing protocols like OSPF (Chapter 48) need to talk to "all OSPF routers on this link" — a specific *group*, not literally everyone.

**The naive attempts and why they fail.** Sending the video stream once per subscriber (unicast to each) wastes enormous bandwidth — a busy stream with 10,000 subscribers sharing a link would need 10,000 copies of every packet crossing that link. Sending it as a broadcast forces every device on the network to receive and discard it even if only 3 of 500 devices care.

**The real solution: multicast.** IPv4 reserves the entire **Class D** range, `224.0.0.0/4` (that is, `224.0.0.0` through `239.255.255.255`), for multicast group addresses. A multicast address doesn't identify one specific device — it identifies a **group**, and any device can choose to "join" that group (using IGMP, the Internet Group Management Protocol, which tells the local router "please start forwarding traffic for group X to me"). The sender transmits exactly *one* copy of each packet; routers and switches along the path replicate it only toward branches of the network that actually have a joiner.

```
                          Sender (one video stream to 224.2.2.2)
                                    |
                              [ Router ]
                             /          \
                    (nobody joined)   (joined 224.2.2.2)
                        no copy sent         |
                                        [ Switch ]
                                       /    |    \
                                 (join) (no join) (join)
                                   PC-A         PC-C
                             receives one copy   receives one copy
                             PC-B never gets a copy at all
```

**Notable reserved multicast addresses:**

| Address | Meaning |
|---|---|
| `224.0.0.1` | "All hosts on this subnet" (multicast equivalent of local broadcast) |
| `224.0.0.2` | "All routers on this subnet" |
| `224.0.0.5` / `224.0.0.6` | OSPF "all OSPF routers" / "all designated routers" (Chapter 48) |
| `224.0.0.0/24` | Reserved for local network control traffic — never forwarded by routers |
| `239.0.0.0/8` | Administratively-scoped multicast (RFC 2365) — the multicast equivalent of RFC 1918 private space, meant to stay inside one organization |

**How a router actually learns who wants a group: IGMP.** Joining a multicast group isn't automatic — a device has to say so, using the **Internet Group Management Protocol (IGMP)**, which operates strictly between a host and its directly-attached router (it never crosses a router hop itself; multicast *routing*, handled by separate protocols like PIM, is what stitches groups together across multiple routers, and is out of scope for this addressing-focused chapter).

IGMP has gone through three versions, each fixing a specific weakness of the one before:

| Version | RFC | Key improvement |
|---|---|---|
| IGMPv1 | RFC 1112 (1989) | Basic join; leaving a group relied purely on a query timeout — a host that left had no way to say so explicitly, so routers kept forwarding traffic to an empty group for longer than necessary |
| IGMPv2 | RFC 2236 (1997) | Added an explicit "Leave Group" message, so a router can stop forwarding almost immediately once the last member leaves, instead of waiting for a timeout |
| IGMPv3 | RFC 3376 (2002) | Added **source filtering** — a host can join a group and specify "only from this specific sender," which is what makes modern **Source-Specific Multicast (SSM)**, using the `232.0.0.0/8` range, possible |

**A subtlety worth knowing: IGMP snooping.** IGMP is a Layer 3 (IP) protocol, but it travels inside ordinary Ethernet frames, which means a plain switch (Chapter 30), which only understands MAC addresses, would otherwise flood every multicast frame out every port — exactly the broadcast-like waste multicast was invented to avoid in the first place. **IGMP snooping** is a switch feature that inspects (or "snoops on") IGMP join/leave messages passing through it, and builds its own internal table of which switch ports actually have a group member behind them, so it can forward multicast frames only out the ports that need them — restoring multicast's efficiency benefit even at the switch level, not just the router level.

```
Switch with IGMP snooping enabled, group 224.2.2.2:

  Port 1 (PC-A): sent IGMP Join for 224.2.2.2  -> forward multicast here
  Port 2 (PC-B): no join seen                   -> do NOT forward here
  Port 3 (PC-C): sent IGMP Join for 224.2.2.2  -> forward multicast here
  Port 4 (uplink to router): always forward     -> multicast source may be upstream
```

Without snooping, the switch would forward the stream out every port, indistinguishable from a broadcast at the Ethernet level, even though the IP layer is technically doing multicast correctly.

**A second, older way multicast reach was controlled: TTL scoping.** Before RFC 2365 formalized administratively-scoped multicast (the `239.0.0.0/8` block from the table above), network operators controlled how far a multicast stream could travel using a much blunter instrument: the IP packet's Time-To-Live field itself (the same TTL field Chapter 45 covers in depth for ordinary unicast routing). A small TTL value on a multicast packet meant it could only survive a few router hops before being discarded, giving an approximate, easy-to-reason-about scope:

```
TTL 0   — never leaves the originating host at all
TTL 1   — stays on the local subnet only (never crosses a router)
TTL 15  — historically treated as "site-local" scope by convention
TTL 63  — historically treated as "region-local" by convention
TTL 127 — historically treated as "world-wide" scope
```

This convention (never formally standardized as strictly as the address-based scoping that replaced it) illustrates something worth internalizing: **the same TTL mechanism this course will use in Chapter 45 to explain ordinary packet forwarding, and again in Chapter 54 to explain how `traceroute` cleverly exploits TTL expiration, was also, historically, multicast's first tool for containing its own blast radius** — one field, three completely different-looking use cases, once you know where to look for it.

**Why multicast never fully took over the public internet.** In theory, multicast is the efficient answer to one-to-many delivery. In practice, inter-domain multicast (multicast that crosses from one ISP's network into another's) never saw wide deployment — it requires every router along the path to cooperate on group membership and forwarding state, across organizational and trust boundaries that, unlike unicast BGP routing, nobody had strong commercial incentive to build out universally. This is why large-scale "streaming to millions" on the modern internet (Chapter 96, CDNs) is solved with unicast delivery from many geographically distributed edge caches instead — an economically simpler, if less elegant, answer to the same problem multicast was invented to solve. Multicast remains very much alive, however, *within* single administrative domains: enterprise video distribution, financial market data feeds, and routing protocols like OSPF all depend on it daily.

---

## 9. Other Special-Use Ranges You Will Meet

RFC 1918, loopback, broadcast, and multicast are the ones you'll use constantly, but IANA's special-use address registry (RFC 6890 and successors) reserves several other ranges worth knowing by sight:

| Range | Name | Purpose |
|---|---|---|
| `169.254.0.0/16` | Link-local (APIPA) | A device auto-assigns one of these when DHCP (Ch. 55) fails, so it can at least talk to other devices on the same broken segment |
| `192.0.2.0/24`, `198.51.100.0/24`, `203.0.113.0/24` | Documentation (TEST-NET-1/2/3) | Reserved specifically for use in books, RFCs, and examples — like this very course — so nobody accidentally documents someone's real address |
| `100.64.0.0/10` | Shared Address Space (CGNAT) | Used by ISPs for Carrier-Grade NAT, previewed in Chapter 41 and revisited with IPv4 exhaustion in Chapter 42 |
| `0.0.0.0/8` | "This network" | `0.0.0.0` itself means "unspecified address," commonly used to mean "bind to all interfaces" |
| `240.0.0.0/4` | Reserved for future use | Never allocated; historically unusable on most stacks |

You'll notice this chapter has already used `203.0.113.9` and `198.51.100.44` as example public addresses — that's not a coincidence; they're drawn deliberately from the documentation ranges above.

---

## 10. Full Worked Example: Classifying a Batch of Addresses

Suppose a packet capture (Chapter 119 covers the real tool) shows the following destination addresses. Classify each one using everything in this chapter:

| Address | Classification | Why |
|---|---|---|
| `10.4.2.9` | Private | Falls in `10.0.0.0/8` |
| `172.20.5.5` | Private | Falls in `172.16.0.0/12` (16–31 in the second octet) |
| `172.32.5.5` | **Public** | 32 is outside `172.16.0.0/12`'s range (16–31) — common trap! |
| `192.168.50.1` | Private | Falls in `192.168.0.0/16` |
| `127.0.0.1` | Loopback | The classic loopback address |
| `127.5.5.5` | Loopback | Still inside `127.0.0.0/8` |
| `8.8.8.8` | Public | Google Public DNS — globally routable |
| `192.168.1.255` | Broadcast | All-ones host portion of `192.168.1.0/24` |
| `255.255.255.255` | Limited broadcast | Reserved special address |
| `224.0.0.5` | Multicast | OSPF "all routers," inside `224.0.0.0/4` |
| `169.254.1.5` | Link-local | APIPA self-assigned range |
| `169.254.1.5` | Link-local | (repeated to emphasize: this means DHCP likely failed on that host) |
| `203.0.113.7` | Public (in principle) — but reserved for documentation | Never actually routed on the real internet |
| `100.64.5.5` | Shared Address Space (CGNAT) | Inside `100.64.0.0/10` — not RFC 1918, not public, used only between a customer and their ISP's own NAT layer |
| `239.5.5.5` | Administratively-scoped multicast | Inside `239.0.0.0/8` — meant to stay within one organization, unlike globally-scoped multicast |
| `0.0.0.0` | Unspecified / "this network" | Commonly seen as a bind address meaning "all interfaces," never a valid destination on the wire |

The `172.32.5.5` row is the one worth staring at: RFC 1918's middle block is `172.16.0.0/12`, whose second octet range is 16–31 inclusive (because a /12 mask fixes only the top 4 bits of that octet). `172.15.x.x` and `172.32.x.x` look deceptively similar but are fully public — this is a genuinely common real-world misconfiguration, and Chapter 38's binary subnetting drills are exactly the skill that avoids it.

**Reading it off a real capture.** A Wireshark-style summary of a few seconds of traffic on a home LAN makes the same classification exercise concrete, straight from a live capture rather than a table you're handed:

```
No.  Time      Source           Destination        Protocol  Info
1    0.0001    192.168.1.50     239.255.255.250    SSDP      M-SEARCH * HTTP/1.1   (admin-scoped multicast, UPnP discovery)
2    0.0450    192.168.1.50     8.8.8.8            DNS       Standard query A example.com  (public — leaving the LAN via NAT)
3    0.0890    0.0.0.0          255.255.255.255    DHCP      DHCP Discover                 (unspecified source, limited broadcast dest)
4    0.1020    192.168.1.1      192.168.1.255      ARP-ish   (broadcast, subnet-directed — router announcing itself)
5    0.1500    127.0.0.1        127.0.0.1          TCP       [SYN] 54021 -> 5432           (loopback — local Postgres client)
```

Every single one of the five packets above is classifiable using nothing but the address ranges in this chapter — no payload inspection required.

---

## 11. Hands-On Experiment

Run these on any Linux or macOS machine (the network toolbox is covered in full in Chapter 56, but a preview is useful now):

```bash
# See your machine's actual assigned address(es)
$ ip addr show          # Linux
$ ifconfig               # macOS / older Linux

# Ping loopback — this never touches a physical NIC
$ ping -c 3 127.0.0.1
PING 127.0.0.1 (127.0.0.1): 56 data bytes
64 bytes from 127.0.0.1: icmp_seq=0 ttl=64 time=0.045 ms
64 bytes from 127.0.0.1: icmp_seq=1 ttl=64 time=0.038 ms

# Find your subnet's broadcast address
$ ip addr show eth0 | grep brd
    inet 192.168.1.42/24 brd 192.168.1.255 scope global eth0

# Ping the subnet broadcast (many OSes suppress the reply by default for safety)
$ ping -c 2 -b 192.168.1.255      # Linux requires -b to allow broadcast pings

# Confirm your public-facing address as seen from outside your NAT (Chapter 41 explains why
# this differs from `ip addr show`)
$ curl -s ifconfig.me
203.0.113.9
```

A minimal Python program that joins a multicast group and listens — a good way to *feel* the difference between unicast and multicast delivery:

```python
import socket
import struct

MCAST_GRP = "224.1.1.1"
MCAST_PORT = 5007

sock = socket.socket(socket.AF_INET, socket.SOCK_DGRAM, socket.IPPROTO_UDP)
sock.setsockopt(socket.SOL_SOCKET, socket.SO_REUSEADDR, 1)
sock.bind(("", MCAST_PORT))

# Tell the kernel: I want to receive traffic sent to this multicast group
mreq = struct.pack("4sl", socket.inet_aton(MCAST_GRP), socket.INADDR_ANY)
sock.setsockopt(socket.IPPROTO_IP, socket.IP_ADD_MEMBERSHIP, mreq)

print(f"Listening on multicast group {MCAST_GRP}:{MCAST_PORT} ...")
while True:
    data, addr = sock.recvfrom(1024)
    print(f"Received {data!r} from {addr}")
```

Run that on two machines on the same LAN, then send one UDP packet to `224.1.1.1:5007` from a third machine — both listeners receive it from a single transmission, which is the entire point of multicast.

---

## 12. Common Misconceptions

- **"Private IP addresses are more secure."** Private addressing provides *isolation from direct routing*, not encryption, authentication, or any cryptographic protection. A device on your LAN with a private address can still be attacked by another device on the same LAN, and once you cross a VPN (Chapter 85) or a compromised router, private addressing offers no protection at all.
- **"127.0.0.1 is the only loopback address."** It's the conventional one, but the entire `127.0.0.0/8` block loops back — `127.0.0.2` works identically on most systems.
- **"Broadcast traffic goes everywhere on the internet."** It's strictly confined to one broadcast domain (Chapter 30) — routers never forward it, by design, precisely to prevent this.
- **"172.16.x.x through 172.999.x.x is all private."** Only `172.16.0.0` through `172.31.255.255` is private. `172.32.0.0` and above is ordinary public space, as the worked example above showed.
- **"Multicast means the same thing as broadcast, just smaller."** Multicast is opt-in (via IGMP group membership); broadcast is unconditional delivery to everyone. Multicast can also cross router boundaries under multicast routing protocols in ways broadcast structurally cannot.
- **"A device's loopback address is reachable from other devices on the network if I just use its real hostname."** `127.0.0.1` (or `::1`) always means "myself, whoever is asking" — it can never be used by one machine to reach another, no matter what hostname or DNS trick is attempted; the kernel intercepts it before any packet is ever built for the wire.
- **"IGMP is what actually moves multicast traffic across the internet."** IGMP only handles the local host-to-router "I want to join this group" signaling (Section 8). Getting multicast traffic to actually flow correctly across multiple routers requires a separate multicast *routing* protocol (like PIM), which is a materially harder problem this chapter deliberately doesn't cover in depth.

---

## 13. Production Notes

- Cloud VPCs (Chapter 97) default to RFC 1918 ranges internally — AWS defaults to `172.31.0.0/16` for the default VPC, and it's common practice to carve `10.0.0.0/8` into per-team or per-environment `/16`s in larger organizations.
- Docker's default bridge network uses `172.17.0.0/16` by default, which is why you'll sometimes see conflicts if your corporate VPN also hands out addresses in `172.16.0.0/12` — both are drawing from the same RFC 1918 block.
- Databases, caches (Redis), and internal admin panels are frequently configured to bind only to `127.0.0.1` or a private-range address, specifically so a public-facing misconfiguration can't accidentally expose them — a real, common cause of data breaches when that convention is violated.
- Multicast is heavily used inside data centers for things you'd never suspect: some financial exchanges distribute market tick data over multicast because it must reach hundreds of subscribers simultaneously with minimal added latency from replication.
- IGMP snooping is enabled by default on essentially all modern managed switches, but on some older or budget unmanaged switches it may be absent entirely, which is a real, occasionally-encountered cause of unexplained "multicast video stream saturates the whole network" incidents in smaller office or home-lab deployments.
- Kubernetes and other container platforms lean heavily on the loopback-address convention: a pod's own containers commonly reach each other over `127.0.0.1` because they share a network namespace, a detail that becomes important once you reach the container-networking material in Chapter 103.

---

## 14. What's Simplified Here

- IPv6 has its own, structurally different set of private/special addresses (unique local `fc00::/7`, link-local `fe80::/10`) — Chapter 42 covers these on their own terms rather than as a direct IPv4 translation, because the design philosophy genuinely differs (notably, IPv6 has no broadcast at all — multicast fully replaces it).
- This chapter treats "the internet filters RFC 1918 traffic" as a simple universal rule; in reality it depends on correct configuration at every ISP's edge, and misconfigured routers occasionally do leak private-address traffic onto the public internet, which is one of the categories of "bogon" traffic that route-hijacking defenses (Chapter 52) also have to consider.
- Real-world multicast deployment (PIM-SM, PIM-DM, IGMP snooping on switches) is considerably more mechanically involved than "devices join a group" — this chapter gives you the addressing model, not the full multicast routing protocol stack, which is out of scope for this course.
- TTL-based multicast scoping (Section 8) is presented here as a simple historical convention; real deployments that still use it (a shrinking minority, given administrative scoping's cleaner semantics) often layer additional operator-specific TTL thresholds on top of the rough values shown, which this chapter does not attempt to catalog exhaustively.

---

## 15. Interview Questions & Model Answers

**Beginner: "What's the difference between a private and a public IP address?"**

*Model answer:* "A public IP address is globally unique and routable across the entire internet — routers everywhere know how to forward packets toward it. A private IP address, drawn from one of the RFC 1918 ranges (10.0.0.0/8, 172.16.0.0/12, or 192.168.0.0/16), is only meaningful inside one local network; it can be reused by millions of unrelated networks simultaneously because ISPs agree never to route that address space across the public internet. Devices with private addresses reach the internet via NAT, which rewrites their address to a shared public one."

**Intermediate: "Why can't `172.20.0.0` and `172.40.0.0` both be considered private?"**

*Model answer:* "RFC 1918's middle private block is `172.16.0.0/12`, not the whole `172.0.0.0/8`. A /12 mask fixes the first 12 bits, which constrains the second octet to values 16 through 31 in binary (00010000 through 00011111). `172.20.x.x` falls inside that range (20 is between 16 and 31), so it's private. `172.40.x.x` does not (40 is outside 16–31), so it's an ordinary public address, even though it superficially looks similar. This is a genuinely common misconfiguration source in real networks."

**Advanced: "Why did inter-domain IP multicast never see wide adoption on the public internet, even though it's clearly more efficient than unicast for one-to-many delivery?"**

*Model answer:* "Efficiency at the packet level doesn't automatically translate into deployability at the internet level. Inter-domain multicast requires every router along a potential path — across different ISPs, with different commercial relationships and no inherent trust — to maintain per-group forwarding state and cooperate on protocols like PIM and MSDP. Unlike unicast routing, where BGP (Chapter 49) gives ISPs a clean, incentive-aligned way to advertise reachability and get paid for transit, there was no comparable commercial mechanism that made ISPs want to carry someone else's multicast group traffic for free. The industry instead solved one-to-many delivery at internet scale with CDNs (Chapter 96) — many unicast copies served from distributed edge caches — which requires no cross-ISP cooperation at the routing layer at all, just money changing hands for hosting and bandwidth. Multicast remains standard practice *within* single administrative domains (enterprise networks, financial data feeds, data centers) where that trust and control problem doesn't exist."

**Advanced: "What is IGMP snooping, and why does a switch need it if multicast is already handled correctly at the IP layer?"**

*Model answer:* "IP-layer multicast correctly identifies which hosts *want* a group's traffic via IGMP, but a plain Ethernet switch operates at Layer 2 and has no native concept of an IP multicast group — a multicast IP address maps onto a multicast MAC address, but without any additional logic, a switch treats an unrecognized multicast MAC much like a broadcast, flooding it out every port. That defeats the entire efficiency purpose of using multicast in the first place. IGMP snooping fixes this by having the switch passively inspect the IGMP join and leave messages that pass through it — even though the switch isn't the intended recipient of those control messages — and build its own per-port membership table. It then forwards multicast frames only out the ports that have actually joined that group, restoring the efficiency benefit at Layer 2 that IGMP alone only guarantees at Layer 3."

---

## 16. Exercises

### Easy

1. Classify each address as private, public, loopback, broadcast, or multicast: `10.0.0.1`, `224.0.0.1`, `127.0.0.1`, `172.16.5.5`, `8.8.4.4`.
2. What is the broadcast address for the subnet `192.168.100.0/24`? For `10.0.0.0/16`?
3. Why does `ping 127.0.0.1` succeed even if your network cable is unplugged?

### Medium

4. A device shows itself configured with `169.254.32.10`. What almost certainly went wrong, and what protocol from a later chapter would normally have prevented it?
5. Explain, in your own words, why `172.16.0.0/12` covers addresses `172.16.x.x` through `172.31.x.x` and not further, using the binary breakdown of the second octet.
6. You are designing address space for a company with 40,000 internal devices spread across 12 sites. Which RFC 1918 block would you choose as your base, and why not the others?

### Hard

7. A misconfigured router is observed forwarding packets destined for `10.0.0.0/8` onto the public internet. What real-world category of problem does this fall into (see Chapter 52 for the vocabulary), and why is universal cooperation, rather than any single technical enforcement mechanism, the thing actually preventing this at scale?
8. Design (in pseudocode or plain English) the logic a router or switch would need to implement to correctly handle a packet addressed to `224.0.0.5` differently from a packet addressed to `224.5.5.5` — considering both are "multicast" but one is meant to never leave the local segment.
9. Explain why IPv6, covered in Chapter 42, was able to eliminate broadcast entirely and rely purely on multicast — what property of IPv6 deployment (hint: think about scale of addressing and design timing) made that possible where IPv4 was stuck with a broadcast mechanism from the start.
10. A switch without IGMP snooping is added to a network that already has a working multicast video stream. Predict, in concrete terms, what changes about the traffic pattern on that switch's ports, and why — using Section 8's discussion of snooping.

---

## 17. Summary

| Term | Meaning |
|---|---|
| RFC 1918 | The standard reserving `10.0.0.0/8`, `172.16.0.0/12`, `192.168.0.0/16` for private, freely-reusable addressing |
| Private address | Reusable, non-globally-unique address, unroutable on the public internet |
| Public address | Globally unique, IANA-allocated, routable address |
| Loopback | `127.0.0.0/8` (IPv4) / `::1` (IPv6) — a machine talking to itself, never touching a NIC |
| Subnet broadcast | All-host-bits-1 address of a subnet — reaches every device on that one LAN segment |
| Limited broadcast | `255.255.255.255` — "everyone on my local wire," regardless of subnet |
| Multicast | `224.0.0.0/4` (IPv4) — an opt-in group address; one transmission, many receivers, via IGMP membership |
| Link-local (APIPA) | `169.254.0.0/16` — a device's fallback self-assigned address when DHCP fails |
| Bogon filtering | ISPs' cooperative practice of dropping traffic that claims to be from/to private or unallocated space |

Private addressing solves the scarcity and isolation problems — but it creates a new one of its own: a device with a private address genuinely *cannot* be reached from, or reach, the public internet without help. Chapter 41 opens with exactly that problem, and shows the mechanism — Network Address Translation — that quietly makes billions of privately-addressed devices work on a public internet that never sees their real addresses at all.
