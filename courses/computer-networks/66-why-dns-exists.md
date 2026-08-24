# Chapter 66: Why DNS Exists — Names vs. Numbers

> **"In 1982, the entire Internet's address book was one file, called HOSTS.TXT, maintained by one woman at one institution, downloaded by every machine on the network by hand. It worked. And then the Internet grew, and it stopped working — not gradually, but in a way that made it obvious a flat file could never scale. DNS is the answer to exactly that failure."**

---

## Table of Contents

1. [The Big Question](#1-the-big-question)
2. [Numbers Are for Routers, Names Are for Humans](#2-numbers-are-for-routers-names-are-for-humans)
3. [The Real History: HOSTS.TXT](#3-the-real-history-hoststxt)
4. [The Distribution Mechanism, Concretely](#4-the-distribution-mechanism-concretely)
5. [Why a Flat File Cannot Scale](#5-why-a-flat-file-cannot-scale)
6. [The Actual Growth Timeline](#6-the-actual-growth-timeline)
7. [What a Real Fix Needs](#7-what-a-real-fix-needs)
8. [Enter Paul Mockapetris and DNS](#8-enter-paul-mockapetris-and-dns)
9. [The Shift in Shape: From a List to a Tree](#9-the-shift-in-shape-from-a-list-to-a-tree)
10. [A Living Fossil: /etc/hosts Today](#10-a-living-fossil-etchosts-today)
11. [Naming Collisions: A Worked Example](#11-naming-collisions-a-worked-example)
12. [Common Misconceptions](#12-common-misconceptions)
13. [What's Simplified Here](#13-whats-simplified-here)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#summary)

---

## 1. The Big Question

Chapter 36 gave every device on the Internet a 32-bit (or, with IPv6, 128-bit) number: an IP address. That number is exactly what a router needs — a compact, hierarchical, comparable value it can use for longest-prefix-match forwarding (Chapter 45). It is exactly what a human being does not want to type, remember, or say out loud.

Try to imagine the Internet if it had never grown past IP addresses. To check your email, you'd type `142.250.80.46` into a browser. To tell a friend where your blog lives, you'd read them out four (or eight) groups of numbers. Businesses would print IP addresses on billboards. Nobody would remember them, typos would be common and silent (one wrong digit takes you to a different, real server, not an error page), and worse: if a company moved its website to new hardware with a new address, every billboard, business card, and bookmark in the world would instantly become wrong with no way to fix them.

So the big question this chapter answers is not "how do we look up a name" — that's Chapters 67 and 68. It's more basic: **why does looking up a name need an entire globally distributed system at all? Why not just... a list?**

That's not a rhetorical question. A list is exactly what the early Internet used. This chapter tells you what that list was, why extremely smart engineers were satisfied with it for over a decade, and exactly which moment turned "the list" from adequate into impossible.

---

## 2. Numbers Are for Routers, Names Are for Humans

**Intuitive level.** Think about how you navigate a city. Streets have names, but the postal system, GPS satellites, and delivery trucks ultimately operate on structured numeric codes — postal codes, GPS coordinates, building numbers within a range. You don't experience the city that way. You experience it as "the coffee shop on Elm Street near the library." Two representations of the same reality exist because they serve two different kinds of consumers: machines that need to compute a route, and humans who need to remember a place.

IP addresses and domain names are exactly this split. A router forwarding a packet (Chapter 44) does not care that the destination is "your bank's website" — it cares that the destination falls under a particular prefix reachable via a particular next hop. A human trying to reach their bank does not think in prefixes; they think in names, ideally names tied to the organization's identity (`chase.com`, not `23.15.128.9`).

**Where the analogy breaks.** A postal code rarely changes for a building. An IP address, in contrast, changes constantly and for mundane reasons: a company switches cloud providers, a load balancer adds or removes a server, a CDN reroutes you to a different edge node depending on where you are standing (Chapter 96). The naming system therefore has to do something a postal code never has to do: stay stable while the number underneath it changes freely, possibly many times a minute.

**Engineering terminology.** This is the classic **indirection** problem in computer science — introducing a layer between "the thing a consumer asks for" and "the thing that actually exists," so the underlying thing can change without breaking every consumer. DNS is indirection applied to network addressing: a **name** is a stable handle; the **address** behind it is free to move.

---

## 3. The Real History: HOSTS.TXT

This part is not a simplified retelling for teaching purposes — it is literally what happened.

In the ARPANET era (Chapter 10), before DNS existed, every host on the network needed a way to translate a human-friendly name (like `SRI-NIC` or `MIT-AI`) into a numeric host address. The solution the ARPANET community adopted was blunt and effective at small scale: **one single text file, called `HOSTS.TXT`**, that listed every hostname on the network next to its numeric address.

```
Excerpt from the spirit of HOSTS.TXT (illustrative format, ARPANET era):

NET   : 10.0.0.0
GATEWAY : MIT-GW, 10.0.0.77
HOST  : SRI-NIC, 10.0.0.51, DEC-2060, TOPS20, TCP/TELNET,TCP/FTP,TCP/SMTP
HOST  : MIT-AI, 10.0.0.41, PDP-10, ITS, TCP/TELNET,TCP/FTP
HOST  : UCLA-CCN, 10.2.0.1, IBM-360/91, OS/MVT, TCP/TELNET
...
```

This file was maintained centrally by the **Network Information Center (NIC)** at **SRI International** (Stanford Research Institute), a role famously associated with Elizabeth "Jake" Feinler and her team, who processed requests to add or change hostnames. The workflow was:

1. A site administrator wanted to add a new host to the network, or change an existing host's address.
2. They emailed (or, earlier, phoned) the request to the NIC.
3. The NIC manually updated the master `HOSTS.TXT` file.
4. Every other site on the network periodically **downloaded the entire updated file via FTP** from the NIC's host and replaced their local copy.
5. Each host on each site then used its local copy of `HOSTS.TXT` to resolve names to addresses for every program that needed it.

This was standardized informally and later documented in early RFCs (RFC 606, RFC 623, RFC 625, RFC 810, culminating in RFC 952's "DoD Internet Host Table Specification"). It worked, for years, because the ARPANET was genuinely small — a few dozen to a few hundred hosts, run by research institutions who mostly knew each other.

---

## 4. The Distribution Mechanism, Concretely

It's worth spelling out exactly *how* a site kept its copy of `HOSTS.TXT` synchronized, because the mechanism itself is what Section 5 shows breaking down. This wasn't an automated push system — it was a scheduled, manual pull, and the steps below are drawn directly from how RFC 952 and the operational practice around it describe the process:

```
Site administrator's actual weekly (or so) routine, ARPANET era:

  1. Connect via FTP to SRI-NIC's well-known host.
  2. Log in (often anonymously, or with a shared community login).
  3. Retrieve the current NETINFO:HOSTS.TXT file in its entirety.
  4. Overwrite the site's local copy with the freshly downloaded one.
  5. Restart or signal any local services that read the hosts table,
     so they'd pick up the new copy.
  6. Repeat next week (or whenever the administrator remembered to).
```

Two things about this process matter enormously for Section 5's argument. First, it is a **full-file replace**, not an incremental patch — there was no concept of "just send me the five lines that changed since last Tuesday." Second, it is **administrator-initiated and manual** — nothing automatically pushed an update to every site the moment the NIC changed the master file; propagation speed depended entirely on how often busy system administrators remembered to re-fetch it. Put those two facts together and you get exactly the consistency and bandwidth problems Section 5 walks through next.

---

## 5. Why a Flat File Cannot Scale

By the early 1980s, the number of hosts on the growing Internet was climbing past a few hundred and heading toward the thousands, doubling roughly every year. `HOSTS.TXT` did not fail gracefully — it failed along four distinct, independent axes at once. Understanding each axis is the entire point of this chapter, because DNS's design (Chapter 67) is a direct, deliberate answer to each one.

### 5.1 The Distribution Problem (Network Load)

Every site needed a synchronized copy of the whole file, so every site had to periodically re-download the **entire file**, even if only one line in it had changed. As the file grew from kilobytes to hundreds of kilobytes, and as the number of participating sites grew from dozens to hundreds, the total bandwidth spent just distributing the address book started to compete with the bandwidth spent doing actual work. This is a classic "N × total-size" scaling problem: total distribution cost grows as (number of sites) × (file size), and both factors were growing simultaneously.

```
Cost of keeping everyone in sync, roughly:

  Year      Hosts     File size (approx)     Sites needing full re-download
  1981      ~200      small                  ~200
  1984      ~1,000    much larger            ~1,000
  1987      ~10,000   far too large          ~10,000 (fails long before this)

Total bandwidth ≈ hosts × file_size — quadratic-ish growth against a
network that had not grown its capacity nearly as fast.
```

### 5.2 The Consistency Problem (Staleness)

Because the file was pulled periodically (not pushed instantly), different hosts on the network held different, aging snapshots of the truth at any given moment. A host that had just changed its address might be reachable under the old address at some sites and the new one at others, for however long it took every site to re-sync. There was no way to make an update "take effect everywhere at once."

### 5.3 The Namespace Problem (Collisions and Authority)

Every hostname had to be globally unique in one flat list, and only the NIC could grant a new name — there was no way to *delegate* a chunk of the namespace to someone else to manage independently. If Stanford and MIT both wanted to name a machine `AI`, someone had to adjudicate that centrally, by hand, for every single request, forever. This is fundamentally a single point of administrative bottleneck, not just a technical one.

### 5.4 The Single Point of Failure Problem

The master copy lived on one host at SRI. If that host was down, or the NIC's staff were unavailable, no new names could be registered and no updates could propagate, network-wide. A tree-shaped, delegated system (Chapter 67) would later spread this authority — and this risk — across thousands of independently operated servers.

**The moment this became undeniable.** By the mid-1980s, projections showed the Internet heading toward tens of thousands of hosts within a few years (a number that, in hindsight, undersold what was coming by many orders of magnitude — there are now well over a billion routable devices). A file-based, centrally-administered, periodically-redistributed address book was mathematically guaranteed to collapse under that trajectory. This wasn't a hypothetical engineering concern — RFC 882, which proposed DNS in 1983, opens by stating almost exactly this problem as its motivation.

---

## 6. The Actual Growth Timeline

It helps to see roughly how fast the underlying number — total hosts on the network — was actually moving, because "doubling every year" sounds abstract until you see what it does to Section 5.1's bandwidth-cost table over a realistic decade. The counts below are the widely cited estimates from historical retrospectives of ARPANET/early-Internet growth (see Section 13's honesty note on precision):

```
Year    Approximate host count     What this meant for HOSTS.TXT
1971    ~20                        Trivial; the file barely needs to exist
1977    ~100                       Still fully manageable by hand
1981    ~200                       RFC 882/883 not yet written; strain starting
1983    ~500                       Mockapetris publishes RFC 882/883 (DNS is born)
1984    ~1,000                     DNS specified but not yet widely deployed
1987    ~10,000                    RFC 1034/1035 (cleaned-up DNS spec); HOSTS.TXT
                                   long since impractical as the PRIMARY mechanism
1989    ~100,000                   DNS fully dominant; HOSTS.TXT relegated to
                                   local overrides and small private networks
```

Notice the timing: DNS was specified in 1983, *while the host count was still only in the hundreds* — Mockapetris and the community that reviewed his RFCs were not reacting to an already-collapsed system, they were extrapolating the trend line (the same doubling pattern visible from 1971 to 1983) and building the replacement ahead of the wall, not after hitting it. That's a meaningfully different, and more impressive, engineering story than "the old system broke, so someone scrambled to fix it" — the fix arrived because someone did the math on the growth curve early enough to matter.

---

## 7. What a Real Fix Needs

Before looking at what Paul Mockapetris actually built, it's worth pausing and deriving, from the four failures above, what *properties* any real replacement would need. This is the "naive attempt, then the real solution" arc this course follows everywhere else, and DNS is no exception.

| Failure of HOSTS.TXT | Property the fix must have |
|---|---|
| Whole file re-sent for one change | Only fetch the *specific* name you need, not the whole namespace |
| Central authority approves every name | **Delegated authority** — let organizations manage their own piece of the namespace |
| Central single point of failure | **Distribute** the responsibility of answering across many independent servers |
| Stale copies drift out of sync | **Timed expiration** of cached answers, tunable per-record (this becomes TTL — Chapter 68) |
| Flat namespace, one shared pool of unique names | A **hierarchical namespace**, where uniqueness is only required *within* a level, not globally |

That last row is the single most important idea in this entire volume, so it's worth restating plainly: in a flat namespace, `ai` must be unique across the *entire Internet*. In a hierarchical namespace, `ai.mit.edu` and `ai.stanford.edu` can coexist peacefully, because uniqueness is only enforced *within* `mit.edu` and *within* `stanford.edu` respectively — and MIT and Stanford each control their own slice without asking anyone's permission. This is precisely the same trick Chapter 36 used to make IP addressing scale (network portion + host portion) and the same trick postal addressing has always used (country, state, city, street, number).

---

## 8. Enter Paul Mockapetris and DNS

In 1983, Paul Mockapetris, working at the University of Southern California's Information Sciences Institute (USC-ISI), published RFC 882 ("Domain Names — Concepts and Facilities") and RFC 883 ("Domain Names — Implementation and Specification"). These two documents are the birth certificate of DNS. (They were later superseded and cleaned up by RFC 1034 and RFC 1035 in 1987, which remain the core specifications engineers cite today.)

The design directly answers every row in the table above:

- **Hierarchical namespace** (the tree structure Chapter 67 covers in full) replaces the flat list.
- **Delegation** lets any organization that owns a piece of the namespace run its own servers for it, with zero involvement from a central authority for day-to-day changes.
- **Distributed servers** replace the single NIC host — thousands of independent authoritative servers exist today, each responsible only for its own small slice.
- **Caching with per-record TTLs** replaces "redownload everything periodically" with "ask only when you need to, and only as often as the data owner says you should."

`HOSTS.TXT` did not disappear overnight. RFC 952's addressing table lived on for years as a fallback and for small private networks, and the transition period saw both systems running side by side. But by the late 1980s, DNS had fully taken over as the Internet's naming system, and it has not been meaningfully replaced since — an unusually long lifespan for any piece of 1983-era infrastructure, and testimony to how correctly Mockapetris identified the actual scaling problem.

---

## 9. The Shift in Shape: From a List to a Tree

It helps to see the shape of the fix, not just read about it, before Chapter 67 goes deep on the mechanics.

```
BEFORE — HOSTS.TXT: one flat, unordered list

  [ SRI-NIC -> 10.0.0.51 ]
  [ MIT-AI  -> 10.0.0.41 ]
  [ UCLA-CCN -> 10.2.0.1 ]
  [ ... one line per host on Earth, all in one file,
        all maintained by one organization ... ]


AFTER — DNS: a hierarchical, delegated tree

                          . (root)
                        /  |  \
                   .com  .org  .edu  ... (TLDs, separately operated)
                    |
              example.com  (delegated to its owner)
               /        \
        www.example.com   mail.example.com
        (owner's own authoritative server answers for this branch)

No single server holds the whole tree.
No single organization approves every name.
Each node only needs to know its own children and how to find its parent.
```

This is the same shift Chapter 24 made you expect in general: a monolithic, centralized design that worked at small scale gets replaced by a layered, delegated one that scales by *dividing responsibility*, not by making any single component bigger and more powerful. You will see this exact pattern again in Chapter 50 (BGP route aggregation replacing one giant routing table) and Chapter 99 (VXLAN replacing VLAN's flat 4096-network ceiling).

---

## 10. A Living Fossil: /etc/hosts Today

Remarkably, `HOSTS.TXT`'s direct descendant is still on your computer right now, under a different name: `/etc/hosts` on Linux and macOS, or `C:\Windows\System32\drivers\etc\hosts` on Windows. It is a tiny, local, manually-edited flat file mapping names to addresses — exactly the format of the original — and your operating system still checks it *before* asking DNS.

### Hands-On Experiment

```bash
# View your local hosts file (the same idea as 1983's HOSTS.TXT, just local-only now)
cat /etc/hosts

# Typical output:
# 127.0.0.1       localhost
# ::1             localhost
# 255.255.255.255 broadcasthost

# Add a private mapping (requires sudo) and watch it override DNS:
echo "93.184.216.34  myfakeexample.test" | sudo tee -a /etc/hosts
ping myfakeexample.test
# Your machine resolves this name WITHOUT ever asking a DNS server,
# because /etc/hosts is checked first (this order is itself
# configurable — see /etc/nsswitch.conf on Linux).
```

This is exactly why the flat-file approach still works for a *single machine's* handful of overrides — the same design that broke at a few hundred entries shared across a growing Internet works fine at ten entries on one laptop that never needs to synchronize with anyone else. Scale, not the idea itself, was always the actual enemy.

---

## 11. Naming Collisions: A Worked Example

Section 5.3 asserted that a flat namespace forces central adjudication of every naming conflict. It's worth seeing exactly what that adjudication looked like versus what DNS makes unnecessary.

```
UNDER HOSTS.TXT (flat namespace, ~1980):

  Stanford AI Lab wants to register a host named "AI".
  MIT AI Lab ALSO wants to register a host named "AI".

  Both requests reach the NIC. Only one can be granted "AI" outright.
  The NIC must intervene, by hand, and typically resolves this by
  assigning qualified variants instead — e.g., "SU-AI" (Stanford
  University AI) and "MIT-AI" — a workaround that is really just a
  human being manually inventing a crude, informal hierarchy inside
  a system that has no formal concept of one.

  This had to happen for EVERY name collision, forever, as the
  network grew — an ever-growing manual workload for one team.


UNDER DNS (hierarchical namespace, post-1983):

  Stanford registers under its own delegated zone: ai.stanford.edu
  MIT registers under its own delegated zone:      ai.mit.edu

  Neither organization needs to ask the other, or any central
  authority, anything. Both names are simultaneously valid because
  uniqueness is enforced only WITHIN each delegated zone
  (.stanford.edu and .mit.edu respectively), not across the whole
  Internet. The "SU-" and "MIT-" prefixing trick administrators were
  doing by hand under HOSTS.TXT is exactly what DNS's hierarchy now
  does automatically, structurally, for every organization at once.
```

This is worth sitting with, because it shows that DNS's hierarchy isn't a new idea invented from nothing — it's the formalization of a workaround administrators were *already* doing manually under the old system, turned into an actual structural guarantee instead of a case-by-case human judgment call.

---

## 11.5 Production Notes: The Same Split Still Exists Today

The three-way split DNS eventually formalized — *who owns the naming policy*, *who operates the registration process*, and *who actually answers queries for a given name* — did not disappear once the technology matured; it hardened into distinct commercial and organizational roles that still exist today, and that Chapter 67 names precisely:

- A **registry** owns and operates a TLD (Verisign for `.com`, NIXI for `.in`) — the modern, formalized descendant of "whoever gets to decide what names are valid at this level."
- A **registrar** is the customer-facing business you actually buy a domain from (Namecheap, GoDaddy, Google Domains' successor Squarespace Domains) — a role that didn't meaningfully exist in the HOSTS.TXT era, created specifically because DNS made it *possible* to have many independent businesses compete to sell registrations under a shared TLD, instead of one team at one institution processing every request by hand.
- A **DNS host** actually runs the authoritative servers answering for your domain's records (Cloudflare, AWS Route 53) — often a completely different company from your registrar, something that would have been meaningless under a single, monolithic HOSTS.TXT model.

Seeing this triangle helps make concrete just how much structural flexibility the shift from Section 5's centralized model bought the Internet: three independent, competitive markets exist today (registries, registrars, DNS hosting) where there used to be exactly one team, doing everything, by hand, for the whole network.

---

## 12. Common Misconceptions

- **"DNS is just a database."** It's a distributed system with a specific delegation model and caching semantics, not a single database anyone can query directly for the whole namespace. No server, anywhere, holds every DNS record on Earth.
- **"HOSTS.TXT failed because it was badly engineered."** It was *well* engineered for the scale it was built for. It failed because a specific input to the problem — the number of participating hosts — grew past the assumptions baked into the design. This is a lesson about scale, not competence.
- **"DNS replaced IP addresses."** No — DNS sits *on top of* IP addressing (Chapter 36) and translates names into it. IP addresses are still exactly what routers forward on; DNS never touches a router's forwarding table.
- **"/etc/hosts is deprecated."** It's still actively used today — for local development, container name resolution (Docker injects entries into it), and blocking known-malicious domains by pointing them at `0.0.0.0`.

---

## 12.5 Why This Story Keeps Repeating

It's worth noticing, before moving on, that "a flat namespace, centrally administered, works fine until participant count crosses some threshold, then collapses" is not a story unique to 1983's ARPANET. The identical shape of problem — and often, eventually, the identical hierarchical fix — shows up any time a naming or identity system that started small tries to keep growing:

- Early package registries for a programming language sometimes start with a flat namespace (first-come-first-served package names, no ownership hierarchy) and later have to bolt on scoped/namespaced packages (`@myorg/mypackage`) once collisions and squatting become a real problem — the exact same "flat pool, add hierarchy" move DNS made in 1983.
- Early internal company wikis or file shares sometimes start as one shared flat folder and inevitably get reorganized into per-team, per-project directory hierarchies once enough people are adding content that name collisions and "whose responsibility is this" become constant friction.
- Even human naming itself follows this shape historically: many cultures moved from single given names to given-name-plus-family-name systems specifically as populations within a single village or region grew large enough that a single name was no longer sufficiently unique.

The lesson worth carrying forward from this chapter, past DNS specifically, is a general one: **a flat namespace is not a design mistake — it's often the right choice at small scale, because it's simpler to build and reason about.** The mistake is failing to notice, or failing to plan for, the point at which growth in the number of participants makes a hierarchical, delegated structure not just nicer but mathematically necessary. Mockapetris got the timing right for DNS (Section 6). Not every system's designers do.

---

## 13. What's Simplified Here

This chapter tells the historically accurate shape of the HOSTS.TXT story, but a few honest caveats: the exact file format evolved across several RFCs (606, 623, 625, 810, 952) rather than being fixed once; the "one file, one maintainer" description is accurate for the *canonical* master copy, though some sites ran modified or delayed local copies; and the precise host counts by year in Section 6 are widely cited estimates from historical accounts (including Mockapetris's own RFC 882 and later retrospectives), not a precisely audited census — record-keeping from 1971–1984 was itself informal, appropriately for a research network of that era. The Stanford/MIT naming example in Section 11 is illustrative of the *kind* of adjudication the NIC had to perform, not a claim about a specific, documented historical dispute between those two institutions by name.

---

## 14. Interview Questions & Model Answers

**Beginner: Why do we need DNS if computers only understand IP addresses anyway?**//
Because IP addresses are not just hard for humans to remember — they also change over time (a server migrates, a load balancer changes, a CDN redirects you to a different edge node), while a name like `example.com` needs to stay stable. DNS provides a layer of indirection: humans and applications use names, and DNS translates those names to whatever the current address happens to be, without anyone needing to update every reference to that name.

**Intermediate: What specifically was wrong with using a single shared hosts file for the whole Internet, in engineering terms?**
Four independent failures: (1) bandwidth cost of redistributing the entire file to every site for even a single change, scaling roughly as sites × file size; (2) staleness — different sites held different snapshots at any given time since updates were pulled periodically, not pushed; (3) a single administrative bottleneck (the NIC) had to approve every new name, and every name had to be globally unique in one flat namespace with no way to delegate a subspace to another organization; (4) a single point of failure — if the master file's host was unreachable, no new names could be registered anywhere.

**Advanced: DNS solves the flat-namespace problem with hierarchy and delegation. Where else in this course have you seen the same "flat doesn't scale, hierarchy does" pattern, and why does it recur?**
At least three other places: IPv4/IPv6 addressing itself (Chapter 36–39) splits an address into network and host portions specifically so routers don't need a route for every individual host — only for aggregated prefixes; BGP route aggregation (Chapter 50) lets an ISP announce one summarized prefix instead of thousands of individual customer routes; and VXLAN (Chapter 99) replaces VLAN's flat 4096-ID space with a much larger, delegatable 24-bit identifier space. The pattern recurs because it's a structural fact about scaling any naming or addressing system: a flat namespace requires global coordination for every single entry, while a hierarchical one only requires coordination *within* each level, letting the system grow by adding independent branches instead of contesting a single shared pool.

**Advanced: Why was DNS specified in 1983 while HOSTS.TXT was still technically functioning, rather than after it had already failed?**
Because the community tracking host growth could see the trend line — roughly doubling year over year through the late 1970s and early 1980s (Section 6) — and could extrapolate that a flat, centrally-distributed file would become unmanageable within a few years even though it was still working at the time. This is a case of engineering ahead of a known, projectable scaling wall rather than reactively patching a system after user-visible failure, and it's part of why the transition to DNS, though it took several years, never involved a period where the Internet's naming system was actually broken for end users.

---

## 15. Exercises

### Easy
1. List the four distinct ways HOSTS.TXT failed to scale, in your own words, without looking back at the chapter.
2. Run `cat /etc/hosts` on your machine and identify which lines were put there by your operating system versus (if any) by you or an installed application (e.g., Docker Desktop often adds entries).
3. In your own words, explain what a site administrator had to manually do, and how often, to keep their copy of HOSTS.TXT current, and identify which specific failure mode (Section 5) this manual process directly causes.

### Medium
4. Suppose a network has 500 hosts and the master hosts file is 40 KB. Every host re-downloads the full file once a day. Roughly how much total bandwidth per day does this consume network-wide? Now suppose the network doubles to 1,000 hosts and the file doubles to 80 KB (since it lists twice as many hosts). By what factor did total daily bandwidth grow?
5. Explain, using the specific words "authority" and "delegation," why `ai.mit.edu` and `ai.stanford.edu` can both exist under DNS but `AI` and `AI` could never have both existed under the original ARPANET HOSTS.TXT scheme.
6. Using the growth timeline in Section 6, estimate in which year the total daily distribution bandwidth (hosts × file size, roughly) would have grown enough to plausibly make HOSTS.TXT impractical, and justify your reasoning.

### Hard
7. Design (on paper, no code needed) the minimum set of operations a distributed naming system needs to support to fix all four HOSTS.TXT failures from Section 5 — you should end up independently deriving something very close to the DNS zone/delegation/TTL model covered starting in Chapter 67. Write down each operation and which specific HOSTS.TXT failure it fixes.

---

## Summary

| Term | Meaning |
|---|---|
| IP address | The numeric, router-friendly address of a device (Chapter 36) — stable for machines, unfriendly for humans |
| Domain name | A human-friendly name that stands in for an IP address, free to be re-pointed without the name changing |
| HOSTS.TXT | The ARPANET-era flat text file mapping every hostname to its address, centrally maintained by SRI's NIC |
| Distribution problem | The cost of re-sending an entire shared file to every participant for even one small change |
| Consistency problem | Different sites holding different, aging snapshots of the truth between periodic re-syncs |
| Namespace problem | A flat list requires globally unique names and a central authority to grant them, with no delegation |
| Single point of failure | One master copy on one host means one outage blocks all naming updates, network-wide |
| Delegation | Handing authority over a piece of the namespace (e.g., `mit.edu`) to the organization that owns it |
| RFC 882 / 883 (later 1034 / 1035) | Paul Mockapetris's 1983 specifications that created DNS as the hierarchical, delegated replacement |
| /etc/hosts | The direct, still-living descendant of HOSTS.TXT — a small local flat file checked before DNS |

DNS's actual tree structure — root servers, TLD servers, authoritative servers, and how delegation physically works between them — is the subject of Chapter 67.
