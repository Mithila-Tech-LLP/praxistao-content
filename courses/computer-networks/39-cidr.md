# Chapter 39: CIDR — Classless Addressing

> **"A company needing 300 addresses was once handed 65,536 of them, because the rulebook only offered three sizes: too small, too big, and way too big. CIDR threw out the rulebook."**

---

## Table of Contents

1. [The Problem: A Rulebook With Only Three Sizes](#1-the-problem-a-rulebook-with-only-three-sizes)
2. [Classful Addressing, As It Actually Worked](#2-classful-addressing-as-it-actually-worked)
3. [Worked Example: The Waste, With Real Numbers](#3-worked-example-the-waste-with-real-numbers)
4. [The Second, Independent Crisis: Routing Table Growth](#4-the-second-independent-crisis-routing-table-growth)
5. [The Real Solution: Classless Inter-Domain Routing](#5-the-real-solution-classless-inter-domain-routing)
6. [Worked Example: Right-Sizing the 300-Address Company](#6-worked-example-right-sizing-the-300-address-company)
7. [CIDR Prefix Lengths, In Full](#7-cidr-prefix-lengths-in-full)
8. [Route Aggregation — CIDR's Other Superpower](#8-route-aggregation--cidrs-other-superpower)
9. [Worked Example: Aggregating Four Networks Into One Route](#9-worked-example-aggregating-four-networks-into-one-route)
10. [A Hands-On Experiment](#10-a-hands-on-experiment)
11. [Common Misconceptions](#11-common-misconceptions)
12. [What's Simplified Here](#12-whats-simplified-here)
13. [Interview Questions & Model Answers](#13-interview-questions--model-answers)
14. [Exercises](#14-exercises)
15. [Summary](#15-summary)

---

## 1. The Problem: A Rulebook With Only Three Sizes

Chapters 37 and 38 have already been using CIDR notation freely — `/24`, `/26`, `/21` — without much comment, because by the time IP addressing needed to be taught in this course, that flexible notation had already won. This chapter goes back and tells the story properly: what addressing looked like *before* CIDR, why it was replaced, and exactly how much it wasted.

For the first decade of IP's existence, an organization requesting an address block wasn't offered a custom-sized allocation. It was offered exactly one of three fixed sizes — full stop. If your organization needed 300 host addresses, there was no "just give me a block for 300" option. You got the smallest available size that could technically contain 300 addresses, however enormous that size was, or you were rejected and told to apply for a bigger one. This chapter's central worked example (Section 3) makes that concrete with real numbers, but the shape of the problem should already sound familiar: it's the exact same rigid, non-flexible sizing problem Chapter 37, Section 2 described as the "naive fixed split" that classful addressing represented — and now it's time to see exactly how bad that rigidity actually was in practice, and how it was fixed.

---

## 2. Classful Addressing, As It Actually Worked

The original IPv4 addressing scheme (RFC 791, 1981) divided the entire 32-bit address space into fixed classes, distinguished by the leading bits of the first octet:

```
 Class A:  0xxxxxxx . host . host . host        first octet 0–127    (mask /8,  16,777,214 usable hosts)
 Class B:  10xxxxxx . xxxxxxxx . host . host     first octet 128–191  (mask /16,     65,534 usable hosts)
 Class C:  110xxxxx . xxxxxxxx . xxxxxxxx . host  first octet 192–223  (mask /24,        254 usable hosts)
 Class D:  1110xxxx . ...                        first octet 224–239  (multicast, Chapter 40)
 Class E:  1111xxxx . ...                        first octet 240–255  (reserved/experimental)
```

The class was determined entirely by the leading bits of the address — a router (or a human) could tell a Class A address from a Class C address just by looking at the first octet's range, with no separate mask needing to be communicated at all. That was, in fact, the original appeal: the network/host split (Chapter 37's whole subject) was implied by the address itself, so early IP didn't even need to transmit a mask alongside every address.

But notice the enormous jumps between class sizes: a Class C network holds 254 usable hosts; the next size up, Class B, holds 65,534 — a jump of 258×. There is no size in between. An organization's actual host count almost never lines up neatly with one of these three numbers, and that mismatch is where the waste comes from.

---

## 3. Worked Example: The Waste, With Real Numbers

Take the scenario in this chapter's opening quote literally. A company needs IP addresses for 300 hosts — a mid-sized office, comfortably real-world sized.

**Step 1 — check if a Class C block fits.** A Class C network provides 2^8 − 2 = 254 usable addresses (8 host bits, minus the network and broadcast addresses, exactly as Chapter 38, Section 4 established). 254 < 300. **It does not fit.** There is no way to get a bigger Class C — the class system has no "large Class C" option.

**Step 2 — move up to the next class.** The only larger option is Class B: 2^16 − 2 = 65,534 usable addresses (16 host bits). 65,534 ≥ 300 — it fits, easily.

**Step 3 — compute the waste.** The company is now allocated a block sized for 65,534 hosts and uses 300 of them:

```
 Addresses allocated:   65,536   (the full Class B block, /16)
 Addresses actually used:  300
 Addresses wasted:      65,236

 Utilization:  300 / 65,536  =  0.46%
 Waste:        99.54%
```

Ninety-nine and a half percent of that block sits completely unused, forever, by design — not through mismanagement, but because the classful system offered no size between "not quite enough" (Class C) and "wildly excessive" (Class B). Multiply this by every mid-sized organization on the growing 1980s–1990s Internet — universities, companies, government agencies, each independently needing a few hundred to a few thousand addresses, each forced into a Class B allocation sized for tens of thousands — and the 4.3-billion-address space (Chapter 36, Section 8) was on track to run out from *waste alone*, decades before anywhere near 4.3 billion devices actually existed.

---

## 4. The Second, Independent Crisis: Routing Table Growth

Address waste was only half of the problem CIDR was designed to solve — and understanding the second half explains why the fix (Section 5) does more than just "add more size options."

Every core Internet router keeps a table of known network prefixes and which direction to send traffic for each (Chapters 44–45 cover this table's mechanics in depth). In the strict classful world, every allocated Class A, B, or C network was, by default, its own separate, unrelated routing table entry — because nothing about classful addressing gave routers a reliable way to treat several smaller allocations as one bigger, summarizable chunk (this is the exact capability Chapter 37, Section 9 flagged as the entire point of requiring contiguous masks, previewed there for this chapter).

As the number of individually-allocated Class C networks exploded through the early 1990s — precisely *because* so many organizations needed something Class-C-sized or a little larger — the number of distinct entries backbone routers needed to hold, search, and update on every topology change grew at an alarming rate, threatening to outpace the memory and CPU capacity of the router hardware of that era. This was reported at the time as the **routing table growth crisis**, running in parallel with the address exhaustion crisis, and it needed a fix that solved *both* problems with one mechanism — which is exactly what CIDR turned out to be.

---

## 5. The Real Solution: Classless Inter-Domain Routing

**CIDR** (RFC 1518 and RFC 1519, 1993) threw out the fixed A/B/C class boundaries entirely. Instead of the network/host split being implied by which class an address's leading bits placed it in, CIDR made the split fully explicit and fully arbitrary — precisely the "subnet mask, communicated separately, any length allowed" system Chapters 37 and 38 have already been using throughout this volume.

```
 Classful (implied split):    an address's leading bits determine A/B/C, which determines
                               the mask — only three possible network-portion lengths exist: 8, 16, 24

 Classless / CIDR (explicit split):  the mask (or prefix length) is stated alongside the address,
                                      completely independent of the address's own bit pattern —
                                      any length from /0 to /32 is legal
```

This is why the notation is called "classless": the class-based rule for inferring the split from the address's leading bits is gone. `/13`, `/21`, `/27` — all previously nonsensical, since no class boundary landed there — are now first-class, fully legal prefix lengths, communicated explicitly via the slash notation Chapter 37, Section 7 already introduced. Every worked subnetting problem in Chapter 38 — a `/26`, a `/27`, a `/30`, a `/21` — was only a legal, requestable allocation *because* CIDR removed the classful restriction; before 1993, none of those prefix lengths could have been requested from an upstream provider as an actual allocation, only carved out internally via subnetting of a classful block you already owned.

---

## 6. Worked Example: Right-Sizing the 300-Address Company

Return to Section 3's company, now requesting an address block under CIDR instead of the classful system.

**Step 1 — find the smallest power-of-two block that fits 300 usable hosts.** Using Chapter 38, Section 6's exact technique: find the smallest H such that 2^H − 2 ≥ 300.

```
 H=8:  2^8 - 2 = 254   (too small — same as a Class C, doesn't fit)
 H=9:  2^9 - 2 = 510   (fits, with room to grow)
```

H=9 host bits means a prefix length of 32−9 = **/23**.

**Step 2 — compare the waste.**

```
                    Classful (Class B)     CIDR (/23)
 Allocated:              65,536                512
 Used:                      300                300
 Wasted:                 65,236                212
 Utilization:              0.46%              58.6%
```

The CIDR allocation still has some slack (212 spare addresses, room to grow toward roughly 510 hosts before needing more space) — but 58.6% utilization versus 0.46% is a difference of more than two orders of magnitude in efficiency, for identically meeting the same 300-host requirement. This is not a hypothetical improvement; it is the literal, mechanical result of being allowed to pick H=9 instead of being forced to jump straight to H=16 because no size in between was ever offered.

---

## 7. CIDR Prefix Lengths, In Full

With the class restriction gone, every prefix length from `/0` (the entire IPv4 address space, used only as the "match everything" default route, Chapter 45) to `/32` (a single host, Chapter 37, Section 8, Example C) is legal. A reference table for the ranges most commonly seen in real allocations and subnetting problems:

| Prefix | Mask | Host bits | Total addresses | Usable hosts | Old classful equivalent |
|---|---|---|---|---|---|
| /8 | 255.0.0.0 | 24 | 16,777,216 | 16,777,214 | Class A |
| /16 | 255.255.0.0 | 16 | 65,536 | 65,534 | Class B |
| /20 | 255.255.240.0 | 12 | 4,096 | 4,094 | *(no classful equivalent)* |
| /22 | 255.255.252.0 | 10 | 1,024 | 1,022 | *(no classful equivalent)* |
| /23 | 255.255.254.0 | 9 | 512 | 510 | *(no classful equivalent)* |
| /24 | 255.255.255.0 | 8 | 256 | 254 | Class C |
| /27 | 255.255.255.224 | 5 | 32 | 30 | *(no classful equivalent)* |
| /30 | 255.255.255.252 | 2 | 4 | 2 | *(no classful equivalent)* |

Every row without a classful equivalent is a size that simply could not be requested as an original allocation before 1993 — it could only ever appear as the result of subnetting a classful block you already owned internally (exactly Chapter 38's entire subject). CIDR made every one of these rows a legitimate, independently requestable, independently routable allocation.

---

## 8. Route Aggregation — CIDR's Other Superpower

Section 4 described a second crisis CIDR needed to solve: routing table growth. Removing the classful size restriction, on its own, only fixes address waste — it doesn't obviously fix routing table size, since now organizations get more finely-sized (and therefore potentially *more numerous*) individual allocations. The second half of CIDR's design is what actually addresses this: because prefixes are now arbitrary-length and explicit rather than implied, a set of smaller, adjacent, correctly-aligned prefixes can be represented by a *single, larger* prefix — exactly the reverse operation of Chapter 38's subnetting.

This is called **route aggregation** (or **supernetting**): an ISP that has been allocated (or has itself allocated to customers) several contiguous, power-of-two-aligned blocks can announce just *one* covering prefix to the rest of the Internet, instead of one entry per individual block. Chapter 50 covers this in full production depth — how ISPs plan their allocations specifically to stay aggregatable, and how Autonomous Systems use this to keep the global routing table from re-exploding the way Section 4 described. This chapter previews the mechanism itself, worked by hand, because it is — pleasingly — nothing more than Chapter 38's binary techniques run in reverse.

---

## 9. Worked Example: Aggregating Four Networks Into One Route

An ISP has allocated four separate `/24` blocks to four different customers:

```
 198.51.100.0/24
 198.51.101.0/24
 198.51.102.0/24
 198.51.103.0/24
```

Rather than announcing all four as separate routes to its upstream providers, can it summarize them as one? Convert the varying third octet of each block to binary:

```
 100 = 01100100
 101 = 01100101
 102 = 01100110
 103 = 01100111
       ^^^^^^--   the first 6 bits (011001) are IDENTICAL across all four
             ^^   only the last 2 bits vary, taking every possible combination: 00, 01, 10, 11
```

Because the first 6 bits of the third octet are identical across all four blocks, and the remaining 2 bits of the third octet (combined with the /24's original network bits) sweep through every possible value of those 2 bits exactly once, these four `/24`s combine perfectly into one larger block — a `/22` (two bits shorter than /24, since 2 bits were "given back"):

```
 198.51.100.0/22

 covers exactly:  198.51.100.0 - 198.51.103.255
                  (all four original /24 blocks, and nothing else)
```

Verify no extra addresses leak in: the aggregate block's host portion is now 10 bits (32−22), covering the low 2 bits of the third octet plus all 8 bits of the fourth octet — precisely the same address range as the four original `/24`s stacked together, no more and no less. The ISP's upstream routers now need to hold **one** routing table entry, `198.51.100.0/22`, instead of four — a 4× reduction for this one case, and in real ISP deployments, aggregation routinely collapses thousands of customer-facing `/24`s into a handful of summarized announcements.

This only worked because the four blocks were **contiguous and aligned** — starting at a third-octet value that is itself a multiple of 4 (100 ÷ 4 = 25 exactly), and covering exactly the next 4 consecutive values. If the ISP had instead allocated `198.51.100.0/24`, `198.51.101.0/24`, `198.51.105.0/24`, and `198.51.110.0/24` — technically the same *number* of addresses, but scattered rather than contiguous — no single CIDR prefix could summarize them, and all four would need separate routing table entries regardless of CIDR's existence. This is exactly why real ISPs plan customer address allocations deliberately around aggregation boundaries from the start, rather than handing out blocks in whatever order requests happen to arrive — the full operational story of this planning, plus what happens when it goes wrong, is Chapter 50.

---

## 10. A Hands-On Experiment

You can see real, live CIDR allocations and real route aggregation using public tools that query actual Internet routing data.

**Look up who owns a real CIDR block, and its actual allocated size:**

```bash
$ whois 8.8.8.0/24 | grep -i cidr
CIDR:           8.8.8.0/24

$ whois 1.1.1.0/24 | grep -i -E "netname|cidr"
NetName:        APNIC-1-1-1-0
CIDR:           1.1.1.0/24
```

**See real aggregation in your own router's table, or a public route-server's:**

```bash
$ ip route show
default via 192.168.1.1 dev eth0
192.168.1.0/24 dev eth0 proto kernel scope link src 192.168.1.42
```

Public route servers (operated by real ISPs and IXPs, for exactly this kind of inspection) let you see aggregation at Internet scale — for example, `route-views.routeviews.org` (a public BGP route collector, reachable via Telnet or a web-based looking glass) shows real backbone routing tables where you can search for a specific IP and see the exact aggregated prefix — often something like a `/19` or `/20` — that an ISP announces to cover thousands of individual customer allocations underneath it, mirroring Section 9's worked example at a much larger scale.

**Confirm a CIDR waste calculation with a script:**

```python
>>> def usable_hosts(prefix_len):
...     host_bits = 32 - prefix_len
...     return 2**host_bits - 2
...
>>> usable_hosts(24)   # old Class C size
254
>>> usable_hosts(16)   # old Class B size
65534
>>> usable_hosts(23)   # right-sized CIDR block for 300 hosts
510
>>> 300 / usable_hosts(16) * 100   # classful utilization
0.4577624053764861
>>> 300 / usable_hosts(23) * 100   # CIDR utilization
58.8235294117647
```

This matches Section 6's table exactly, confirming the waste calculation independent of the by-hand arithmetic.

---

## 11. Common Misconceptions

- **"Classful addressing determined a network's size by looking at its first octet, and that's basically the same idea as a subnet mask."** Related, but not the same: classful addressing had exactly three possible network sizes, rigidly tied to which numeric range the first octet fell into, with no separate mask needing to be communicated. A subnet mask (Chapter 37) — and CIDR's prefix length, its modern equivalent — can describe *any* length, is *always* communicated explicitly, and has no relationship to the address's own numeric value. Chapter 37, Section 2 already flagged classful addressing as the "naive fixed split" this whole volume moved past.
- **"CIDR is just about writing masks with a slash instead of dotted-decimal — a notation change, not a real technical fix."** The notation (`/24` instead of `255.255.255.0`) is genuinely just shorthand — but the underlying technical change (arbitrary-length prefixes replacing fixed classes) is what actually solved both the address-waste crisis (Section 3) and the routing-table-growth crisis (Section 4). The slash is a side effect of the real fix, not the fix itself.
- **"Route aggregation works on any group of same-sized blocks."** Worked Example 9 showed the requirement precisely: the blocks must be *contiguous* and *aligned* to a power-of-two boundary matching their combined size. Four /24s chosen at random, even if there are exactly four of them, will not aggregate unless their third-octet values happen to form one properly-aligned group of 4 (Section 9's closing paragraph gives a concrete counterexample).
- **"CIDR replaced subnetting."** The opposite relationship holds: CIDR made subnetting's arbitrary-length prefixes (Chapter 38's /26s, /27s, /21s) legally *requestable* from an upstream provider in the first place. Subnetting (dividing a block you own into smaller pieces) and CIDR (the classless system that lets any prefix length exist at all, and lets prefixes recombine via aggregation) are complementary, not competing, ideas — subnetting is division, aggregation is the same arithmetic run in reverse.

---

## 12. What's Simplified Here

This chapter covered the addressing-efficiency side of CIDR in full, and previewed route aggregation's mechanics, but deliberately did not cover: how routers actually store and search CIDR prefixes at wire speed (longest prefix match, Chapter 45), how Autonomous Systems and ISPs coordinate real-world address allocation and aggregation policy at Internet scale (Chapter 50), or the reserved special-purpose blocks (private ranges, loopback, multicast) that sit alongside ordinary CIDR-allocated space — that's Chapter 40, immediately next.

---

## 13. Interview Questions & Model Answers

**Beginner: "What problem did CIDR solve, and why did classful addressing (Class A/B/C) fail to solve it?"**

Classful addressing only offered three fixed network sizes (roughly 16.7 million, 65,534, and 254 usable hosts). An organization needing an in-between amount — say, 300 hosts — had to be allocated the next size up (Class B, 65,534 addresses), wasting over 99% of that allocation. CIDR replaced the fixed classes with arbitrary-length prefixes (any /N from 0 to 32), letting allocations be sized to the actual requirement instead of rounded up to one of three fixed options.

**Beginner: "What does the 'classless' in CIDR actually mean?"**

It means the network/host split is no longer inferred from the leading bits of the address itself (which is what determined whether an address was Class A, B, or C). Instead, the split (the prefix length) is stated explicitly, independent of the address's own value — any prefix length is legal for any address.

**Intermediate: "A company needs 300 host addresses. Compute how many addresses they'd waste under the old classful system versus under CIDR, and show your work."**

Under classful addressing, a Class C (/24, 254 usable) is too small for 300, forcing an upgrade to Class B (/16, 65,534 usable) — wasting 65,534−300 = 65,234 addresses (over 99.5%). Under CIDR, the smallest block that fits is found by checking 2^H−2 ≥ 300: H=9 gives 510, which fits, corresponding to a /23. That wastes only 510−300 = 210 addresses (about 41%) — dramatically better, because CIDR permits sizes classful addressing simply didn't offer.

**Intermediate: "What two distinct problems was CIDR designed to solve, and how does one mechanism solve both?"**

Address exhaustion (classful over-allocation wasting huge fractions of the 4.3-billion-address space) and routing table growth (each classful allocation being a separate, non-summarizable routing table entry, growing unsustainably as more organizations got their own Class C networks). Arbitrary-length, explicit prefixes solve the first problem directly (right-sized allocations) and enable the second problem's fix, route aggregation — because prefixes are now explicit and arbitrary-length, several small, contiguous, aligned allocations can be represented as one larger prefix, shrinking the routing table even as the number of individual customer networks grows.

**Advanced: "Explain, using a worked binary example, exactly what conditions must hold for four /24 networks to be aggregatable into a single /22."**

The four /24 blocks' varying octet (typically the third octet, for a "/24 to /22" aggregation) must, in binary, share the same leading 6 bits and, across the four blocks, take on every combination of the remaining 2 bits exactly once (00, 01, 10, 11) — equivalently, the block's starting value must be an exact multiple of 4, and the four blocks must be consecutive. For example, third-octet values 100, 101, 102, 103 (binary 01100100 through 01100111, sharing the prefix 011001) aggregate cleanly into `.100.0/22`; values 100, 101, 105, 110 do not, because they aren't the four consecutive values starting at a multiple of 4, and no single prefix boundary covers exactly that set without also covering unallocated addresses in between.

**Advanced: "Why couldn't the routing table growth crisis (Section 4) have been solved just by adding more router memory, rather than redesigning the addressing scheme?"**

Because the growth was structurally unbounded under classful addressing — every new classful allocation was, by design, a separate, non-aggregatable routing table entry, so the table's size scaled linearly with the number of organizations on the Internet with no ceiling and no way to summarize. Adding memory buys time against a linear trend but doesn't change the trend itself; CIDR's aggregation, by contrast, changes the *growth rate* itself, letting the routing table scale closer to the number of large network operators and their deliberately-planned aggregate blocks, rather than the much larger and faster-growing number of individual end-customer allocations.

---

## 14. Exercises

### Easy

1. Under the old classful system, which class (A, B, or C) would an organization needing 1,000 host addresses have been forced into, and how many addresses would be wasted?
2. Compute the CIDR prefix length that would right-size an allocation for 1,000 hosts, and compute the resulting waste. Compare to your answer for question 1.
3. Is `/9` a legal CIDR prefix length? Was it a legal classful boundary? Explain the difference.

### Medium

4. Four /24 blocks — `10.20.4.0/24`, `10.20.5.0/24`, `10.20.6.0/24`, `10.20.7.0/24` — are candidates for aggregation. Show, in binary, whether they aggregate cleanly into a single /22, and state that /22's network address if so.
5. A different set of four /24 blocks — `10.20.4.0/24`, `10.20.5.0/24`, `10.20.6.0/24`, `10.20.9.0/24` — is proposed for the same aggregation. Explain, using binary, why this set cannot be represented as a single CIDR prefix.
6. A company is allocated `172.30.64.0/19`. Using the technique from Section 6, determine how many usable host addresses this provides, and identify which old classful block size (if any) is closest in size without exceeding it.

### Hard

7. An ISP has allocated eight consecutive /24 blocks to eight customers, `203.0.104.0/24` through `203.0.111.0/24`. Determine the single CIDR prefix that aggregates all eight (show the binary justification), and explain what would go wrong for aggregation if the ninth customer were instead given `203.0.120.0/24` and the ISP tried to fold it into the same announcement.
8. Revisit Chapter 38's Worked Problem 3 (the `10.1.0.0/22` split into four /24 sites). Explain, referencing this chapter's Section 9, exactly why that specific split (borrowing the low 2 bits of the third octet) was chosen instead of some other 2-bit split — and what property of the resulting four /24s makes it possible for the company's border router to still advertise just `10.1.0.0/22` to its own ISP, even after each site has been further subnetted internally.

---

## 15. Summary

| Term | Meaning |
|---|---|
| Classful addressing | The original IPv4 system with three fixed allocation sizes (Class A /8, B /16, C /24) |
| CIDR (Classless Inter-Domain Routing) | The system replacing fixed classes with arbitrary-length, explicitly-stated prefixes |
| Prefix length | The number of network bits in a CIDR address, written as `/N` |
| Address waste | Allocated-but-unusable addresses, caused by being forced into a size larger than needed |
| Routing table growth crisis | The 1990s problem of non-aggregatable classful allocations overwhelming router memory |
| Route aggregation (supernetting) | Representing several contiguous, aligned smaller prefixes as one larger prefix |
| Aggregatable | A set of blocks that are contiguous and power-of-two aligned, and can therefore be summarized |

CIDR fixed how much of the address space gets wasted and how large routing tables need to grow — but it assumed, throughout, that every address in a block is equally reachable on the public Internet. That assumption is about to be questioned: some address ranges are deliberately kept off the public Internet entirely, some addresses mean "talk to yourself," and some mean "talk to everyone" or "talk to a specific group." Chapter 40 covers private and public addresses, loopback, broadcast, and multicast — the special-purpose exceptions to everything this volume has built so far.
