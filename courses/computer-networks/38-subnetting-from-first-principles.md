# Chapter 38: Subnetting From First Principles

> **"Subnetting is not a formula. It's a negotiation between how many independent networks you need and how many addresses you can afford to spend on each one — and every 'formula' you've ever seen is just a shortcut for redoing that negotiation by hand."**

---

## Table of Contents

1. [The Real Problem: One Block, Many Networks](#1-the-real-problem-one-block-many-networks)
2. [The Naive (Failed) First Attempt](#2-the-naive-failed-first-attempt)
3. [The Insight: Borrow Bits From the Host Portion](#3-the-insight-borrow-bits-from-the-host-portion)
4. [Every Subnet's Two Reserved Addresses](#4-every-subnets-two-reserved-addresses)
5. [Worked Problem 1: Splitting a /24 Into Four Equal Departments](#5-worked-problem-1-splitting-a-24-into-four-equal-departments)
6. [Worked Problem 2: Unequal Needs — Variable Length Subnet Masking](#6-worked-problem-2-unequal-needs--variable-length-subnet-masking)
7. [Worked Problem 3: Subnetting Hierarchically Across Multiple Sites](#7-worked-problem-3-subnetting-hierarchically-across-multiple-sites)
8. [Worked Problem 4: Reverse-Engineering a Subnet From a Given Address](#8-worked-problem-4-reverse-engineering-a-subnet-from-a-given-address)
9. [A Hands-On Experiment](#9-a-hands-on-experiment)
10. [Common Misconceptions and Mistakes](#10-common-misconceptions-and-mistakes)
11. [What's Simplified Here](#11-whats-simplified-here)
12. [Interview Questions & Model Answers](#12-interview-questions--model-answers)
13. [Exercises](#13-exercises)
14. [Summary](#14-summary)

---

## 1. The Real Problem: One Block, Many Networks

Here is the situation that subnetting exists to solve, stated the way it actually shows up in the real world, with no formulas yet:

> An organization has been handed one block of IP addresses — say, `192.168.10.0/24`, 256 addresses. That organization has four departments (Sales, Engineering, HR, Finance), each in its own physical area of the building, each wanting its own broadcast domain for isolation, security, and manageability (Chapter 30 established why you'd want separate broadcast domains at all — fewer devices hearing every broadcast, and a problem in one department's network doesn't spill into another's). How do you turn one 256-address block into four independently manageable, non-overlapping networks — without wasting addresses you don't have to spare?

Notice what's *not* given: nobody handed you four separate address blocks, one per department. You got exactly one block, and the problem is entirely about how to subdivide it — a problem of division, not allocation. This is the actual origin of the word **subnet**: a smaller network carved out of ("sub-") a larger one.

Chapter 37 gave you the mechanism (subnet masks, the AND operation) but always applied it to networks whose boundary was already decided for you. This chapter is about *deciding* that boundary yourself, deliberately, to solve a real organizational problem — and doing it by hand, in binary, so the reasoning is never a mystery.

---

## 2. The Naive (Failed) First Attempt

**Naive attempt #1: give every department the same network address.** Just configure every department's devices with an address in `192.168.10.0/24` and call it done. This fails immediately and obviously: there's no isolation at all. A broadcast from Sales reaches Engineering, HR, and Finance identically — Chapter 30's entire argument for separate broadcast domains is thrown away. This isn't subnetting; it's just one large flat network with rooms.

**Naive attempt #2: ask for four separate address blocks, one per department.** This "solves" the isolation problem but only by refusing to solve the actual problem posed — you were handed *one* block, and going back for three more (perhaps non-contiguous, perhaps from a different range entirely) creates a second problem: those four networks can no longer be summarized as a single routing table entry to the outside world (Chapter 39's aggregation, fully explored in Chapter 50, depends on a set of subnets forming one describable superblock). Every router upstream would need four separate entries instead of one, forever — that's a real cost, and it's exactly what CIDR-era Internet routing was designed to avoid at scale.

The only approach that gives you real isolation *and* keeps the whole allocation summarizable as one block is to carve the *existing* `192.168.10.0/24` into smaller pieces — which means moving the network/host boundary that Chapter 37 introduced, without leaving the original block at all.

---

## 3. The Insight: Borrow Bits From the Host Portion

A `/24` network has 8 host bits (32 total − 24 network = 8 host), giving 2^8 = 256 addresses. If you take one of those host bits and reclassify it as a network bit instead — moving the mask boundary one position to the right, from `/24` to `/25` — something useful happens: that one bit, which used to distinguish between 256 different *hosts* within one network, now distinguishes between 2 different *networks*, each with 7 remaining host bits (2^7 = 128 addresses).

```
 /24 (before borrowing):
   network bits: 24            host bits: 8   → 1 network of 256 addresses

 /25 (borrow 1 bit):
   network bits: 25            host bits: 7   → 2 networks of 128 addresses each

 /26 (borrow 2 bits):
   network bits: 26            host bits: 6   → 4 networks of 64 addresses each

 /27 (borrow 3 bits):
   network bits: 27            host bits: 5   → 8 networks of 32 addresses each
```

The pattern: **borrowing N bits from the host portion multiplies the number of networks by 2^N, and divides the size of each network by 2^N.** This isn't a rule to memorize — it falls directly out of what a bit *is*: one more bit of network address space doubles how many distinct network prefixes you can express, and correspondingly removes one bit's worth of addressing capacity (a halving) from what's left for hosts. Every subnetting problem in this chapter is just this trade-off, applied deliberately to fit a real requirement.

---

## 4. Every Subnet's Two Reserved Addresses

Before working the first real problem, one more piece of mechanism: within any subnet, exactly two addresses are never assignable to a host.

- The **network address** — all host bits set to 0 (Chapter 37, Section 5). This represents the subnet as a whole, not any device on it.
- The **broadcast address** — all host bits set to 1. Chapter 40 covers broadcast addressing in depth; for now, just know that a packet sent to this address is meant to reach every device on that subnet simultaneously, so no individual host can own it.

So a subnet with H host bits has 2^H total addresses, but only **2^H − 2** are usable by actual devices. This is the source of a classic real-world gotcha: a `/30` subnet (2 host bits, 4 total addresses) provides only **2** usable addresses — exactly enough for a point-to-point link between two routers, and a very common real allocation for exactly that reason.

---

## 5. Worked Problem 1: Splitting a /24 Into Four Equal Departments

Back to Section 1's scenario: `192.168.10.0/24`, four departments, each comfortably under 60 hosts, wanting equal-sized, isolated networks.

**Step 1 — how many bits to borrow?** Four departments means you need to be able to express 4 distinct network identifiers. 2^2 = 4, so borrowing 2 bits is exactly enough (borrowing 1 bit would only give 2 networks — not enough; borrowing 3 would give 8, more than needed, at the cost of smaller subnets than necessary).

**Step 2 — the new mask.** Original: `/24` (24 network bits). Borrow 2: `/26` (26 network bits, 6 host bits remaining). Each subnet holds 2^6 = 64 addresses, of which 62 are usable (Section 4) — comfortably above the ~60-host requirement per department.

**Step 3 — enumerate the subnets in binary.** The two borrowed bits are the *first two bits* of the fourth octet (the leftmost bits still available to borrow, since the first three octets are already fully "network" under the original `/24`). Those two bits can take four values: `00`, `01`, `10`, `11` — and each combination defines one subnet:

```
 Borrowed bits = 00:  fourth octet = 00 000000  through  00 111111  =   0 –  63
 Borrowed bits = 01:  fourth octet = 01 000000  through  01 111111  =  64 – 127
 Borrowed bits = 10:  fourth octet = 10 000000  through  10 111111  = 128 – 191
 Borrowed bits = 11:  fourth octet = 11 000000  through  11 111111  = 192 – 255
```

**Step 4 — assign and identify network/broadcast/usable range for each:**

| Subnet | Network address | Broadcast address | Usable range | Usable hosts |
|---|---|---|---|---|
| Sales | `192.168.10.0/26` | `192.168.10.63` | `.1` – `.62` | 62 |
| Engineering | `192.168.10.64/26` | `192.168.10.127` | `.65` – `.126` | 62 |
| HR | `192.168.10.128/26` | `192.168.10.191` | `.129` – `.190` | 62 |
| Finance | `192.168.10.192/26` | `192.168.10.255` | `.193` – `.254` | 62 |

**Verify one row in full binary**, to make sure the pattern isn't magic. Take Engineering, `192.168.10.64/26`, and find its network address the Chapter 37 way — AND the address with the mask:

```
   IP (last octet):     01000000   (64)
   Mask  (last octet):  11000000   (255.255.255.192 → /26)
   AND:                 01000000   (64)   → network address confirmed: 192.168.10.64
```

And its broadcast address — flip every host bit (the 6 rightmost bits) to 1:

```
   Network (last octet):  01 000000   (64)
   Flip 6 host bits to 1: 01 111111   (127)   → broadcast address: 192.168.10.127
```

Both match the table. All four subnets are now genuinely separate broadcast domains, carved from one original block, with zero addresses wasted beyond the unavoidable network/broadcast pair each subnet requires (4 subnets × 2 reserved addresses = 8 addresses "spent" on structure, out of 256 total — a 3% overhead for four fully isolated, equally-sized networks).

It helps to see the whole block as a single tree being split in half, twice, rather than four unrelated slices appearing from nowhere:

```
                         192.168.10.0/24  (256 addresses)
                                 |
                 borrow 1 bit  /   \  borrow 1 bit
                               /     \
              192.168.10.0/25         192.168.10.128/25
               (0-127, 128 addr)        (128-255, 128 addr)
                    |                          |
            borrow /  \ borrow          borrow /  \ borrow
                  /    \                      /    \
      .0/26            .64/26          .128/26      .192/26
    (0-63, 64)        (64-127, 64)   (128-191, 64) (192-255, 64)
     Sales           Engineering         HR           Finance
```

Every leaf in this tree is reachable by borrowing bits one at a time, and every leaf's address range is exactly half of its parent's — which is the direct, visual form of Section 3's "borrowing N bits multiplies subnet count by 2^N and divides subnet size by 2^N." This same tree shape is exactly what Worked Problem 2's VLSM approach exploits deliberately: instead of splitting every branch evenly, VLSM stops splitting a branch as soon as it's the right size for its intended purpose, and only keeps subdividing the branches that still need to be smaller.

---

## 6. Worked Problem 2: Unequal Needs — Variable Length Subnet Masking

Real organizations rarely need equal-sized departments. Suppose, from the same `192.168.20.0/24` block, the actual requirements are:

```
 Engineering:              100 hosts
 Sales:                     50 hosts
 HR:                        20 hosts
 Router-to-router link:      2 hosts
```

If you naively applied Problem 1's approach — pick one subnet size to fit the *largest* need (100 hosts needs a /25, since /26 only gives 62) and give every department a `/25` — you couldn't even fit two /25s and a /26 and a /30 into one /24 without careful placement, and you'd badly overallocate for HR and the router link. This is exactly the situation **Variable Length Subnet Masking (VLSM)** exists for: **use a different mask size for each subnet, sized to what that subnet actually needs**, rather than forcing every subnet in the block to be the same size.

The technique, derived rather than memorized: **sort requirements largest to smallest, and allocate the largest first**, from the start of the available space. Allocating largest-first avoids a specific failure mode — if you allocated small subnets first, you could end up with the remaining space fragmented in a way that no longer contains a big enough contiguous block for a large requirement that comes later.

**Step 1 — Engineering needs 100 hosts.** Find the smallest host-bit count H such that 2^H − 2 ≥ 100: H=6 gives 62 (too small), H=7 gives 126 (enough). So Engineering needs a `/25` (32−7=25 network bits, block size 2^7=128). Allocate the first available `/25`-sized block: `192.168.20.0/25`, covering `.0`–`.127`. Network `.0`, broadcast `.127`, usable `.1`–`.126` (126 hosts — comfortably covers 100).

Remaining space: `192.168.20.128` – `192.168.20.255` (128 addresses left, i.e., a `/25`-sized block still to subdivide).

**Step 2 — Sales needs 50 hosts.** Smallest H with 2^H−2 ≥ 50: H=6 gives 62 (enough). So Sales needs a `/26` (block size 64). Allocate the next available `/26`-sized block from the remaining space: `192.168.20.128/26`, covering `.128`–`.191`. Network `.128`, broadcast `.191`, usable `.129`–`.190` (62 hosts).

Remaining space: `192.168.20.192` – `192.168.20.255` (64 addresses left, a `/26`-sized block).

**Step 3 — HR needs 20 hosts.** Smallest H with 2^H−2 ≥ 20: H=5 gives 30 (enough). So HR needs a `/27` (block size 32). Allocate: `192.168.20.192/27`, covering `.192`–`.223`. Network `.192`, broadcast `.223`, usable `.193`–`.222` (30 hosts).

Remaining space: `192.168.20.224` – `192.168.20.255` (32 addresses left, a `/27`-sized block).

**Step 4 — the router-to-router link needs 2 hosts.** Smallest H with 2^H−2 ≥ 2: H=2 gives exactly 2 (enough — this is the classic `/30` point-to-point link). Allocate: `192.168.20.224/30`, covering `.224`–`.227`. Network `.224`, broadcast `.227`, usable `.225`–`.226` (exactly 2 hosts, one for each router).

**Remaining, unused:** `192.168.20.228` – `192.168.20.255` — 28 addresses left over, reserved for future growth.

| Subnet | Requirement | Mask | Block | Usable range | Usable hosts |
|---|---|---|---|---|---|
| Engineering | 100 | /25 | `192.168.20.0` – `.127` | `.1`–`.126` | 126 |
| Sales | 50 | /26 | `192.168.20.128` – `.191` | `.129`–`.190` | 62 |
| HR | 20 | /27 | `192.168.20.192` – `.223` | `.193`–`.222` | 30 |
| Router link | 2 | /30 | `192.168.20.224` – `.227` | `.225`–`.226` | 2 |
| *(spare)* | — | — | `192.168.20.228` – `.255` | — | 28 unused |

Compare the waste here to what Problem 1's *equal*-subnetting approach would have cost if forced onto these same requirements: giving the router link a `/26` (64 addresses) just because that's what "the standard subnet size" happened to be would waste 62 of those 64 addresses on a link that only ever needs 2. VLSM's entire value is visible in that one line of the table: a 2-host requirement costs exactly 4 addresses, not 64, because the mask was sized to the actual need rather than to a fixed convention.

### The Allocation Algorithm as Code

The largest-first VLSM procedure from this section is mechanical enough to implement directly — this is, in essence, what real IP address management (IPAM) tooling does when you hand it a block and a list of requirements:

```go
package main

import (
	"fmt"
	"math"
	"sort"
)

// hostsNeeded finds the smallest number of host bits H such that
// 2^H - 2 >= required usable hosts. This is Section 6's core
// per-subnet computation, generalized.
func hostBitsNeeded(required int) int {
	for h := 1; h <= 30; h++ {
		usable := int(math.Pow(2, float64(h))) - 2
		if usable >= required {
			return h
		}
	}
	panic("requirement too large for a single subnet")
}

type request struct {
	name  string
	hosts int
}

func allocateVLSM(baseNetwork uint32, baseHostBits int, reqs []request) {
	// Largest-first, exactly as derived in Section 6.
	sort.Slice(reqs, func(i, j int) bool { return reqs[i].hosts > reqs[j].hosts })

	next := baseNetwork
	totalAvailable := uint32(1) << uint(baseHostBits)
	used := uint32(0)

	for _, r := range reqs {
		hostBits := hostBitsNeeded(r.hosts)
		blockSize := uint32(1) << uint(hostBits)

		fmt.Printf("%-14s needs %3d hosts -> /%d block, %d addresses, network=%d, broadcast=%d\n",
			r.name, r.hosts, 32-hostBits, blockSize, next, next+blockSize-1)

		next += blockSize
		used += blockSize
	}
	fmt.Printf("Unallocated remaining: %d addresses\n", totalAvailable-used)
}

func main() {
	reqs := []request{
		{"Engineering", 100},
		{"Sales", 50},
		{"HR", 20},
		{"RouterLink", 2},
	}
	// 192.168.20.0/24 -> base host bits = 8 (2^8 = 256 addresses total)
	allocateVLSM(0, 8, reqs)
}
```

Running this prints the same allocation, in the same order, with the same leftover count (28 addresses) as the hand-worked table above — because it's the same algorithm, just executed by a computer instead of by hand. Seeing the two match is the real confirmation that "VLSM" isn't a separate topic from what you already know — it's Section 6's arithmetic, automated.

---

## 7. Worked Problem 3: Subnetting Hierarchically Across Multiple Sites

Now scale the problem up one level. A company is handed `10.1.0.0/22` (32−22=10 host bits → 2^10 = 1024 addresses) and has three branch offices — Site A, Site B, Site C — with a fourth site planned but not yet built. Each site independently needs to run its own subnetting (data VLAN + voice VLAN) internally.

This is a genuinely different kind of problem from Problems 1 and 2: it requires subnetting **in two layers** — first divide the big block among sites, *then* divide each site's piece among that site's own departments. Getting the order right matters, and it previews something Chapter 39 and Chapter 50 build on directly: dividing along clean binary boundaries at the *top* level is what makes it possible to later summarize all three (or four) sites back into one routing announcement.

**Layer 1 — divide the /22 among up to 4 sites.** A `/22` has 10 host bits. Borrowing 2 bits (to get 4 equal pieces, one per site, with room for the planned fourth site) gives `/24` blocks — 2^8 = 256 addresses each:

```
 10.1.0.0/22 borrowed by 2 bits → four /24 blocks:

   Site A:  10.1.0.0/24    (10.1.0.0   – 10.1.0.255)
   Site B:  10.1.1.0/24    (10.1.1.0   – 10.1.1.255)
   Site C:  10.1.2.0/24    (10.1.2.0   – 10.1.2.255)
   Site D:  10.1.3.0/24    (10.1.3.0   – 10.1.3.255)   ← reserved for the future site
```

Verify this in binary. A `/22` fixes the first 22 bits as network: the first 16 bits (both full octets) are network, plus 6 more bits reaching into the third octet (16+6=22). That leaves the third octet's *last* 2 bits, plus all 8 bits of the fourth octet, as the /22's 10 host bits (2+8=10, matching 32−22=10).

To split into 4 site blocks, borrow exactly those 2 leftover host bits of the third octet and reclassify them as network bits. Since `10.1.0.0`'s third octet is `00000000`, its fixed high-order 6 bits are all `0`; only the low-order 2 bits (the ones being borrowed) vary:

```
 third octet, high 6 bits fixed at 000000, low 2 bits borrowed:

    000000 00 = 0   → Site A: third octet 0
    000000 01 = 1   → Site B: third octet 1
    000000 10 = 2   → Site C: third octet 2
    000000 11 = 3   → Site D: third octet 3
```

Borrowing those 2 bits turns the third octet into fully-fixed network bits (6 fixed + 2 borrowed = all 8), which is exactly what a `/24` prefix means (24 network bits = 16 + 8). That leaves only the fourth octet's 8 bits as host bits per site (32−24=8, 2^8=256 addresses) — confirming each site correctly receives a full, clean 256-address `/24`, exactly as intended.

**Layer 2 — within Site A's `10.1.0.0/24`, subnet again for data + voice VLANs.** Site A needs a data VLAN for ~120 devices and a voice VLAN for ~50 phones — this is now just Problem 2's VLSM technique, applied one level deeper, on Site A's own /24:

```
 Data VLAN needs 120 hosts:  2^7-2=126 ≥ 120  → /25.  Allocate 10.1.0.0/25    (.0–.127)
 Voice VLAN needs 50 hosts:  2^6-2=62  ≥ 50   → /26.  Allocate 10.1.0.128/26  (.128–.191)
 Remaining spare:                                      10.1.0.192/26  (.192–.255) for growth
```

Site B and Site C would each run this exact same Layer-2 process independently, inside their own `/24`, with no coordination needed between sites — which is precisely the point of the Layer 1 split: once each site owns a clean, self-contained `/24`, that site's network administrator can subnet it however suits their local needs without touching anyone else's allocation. This independence — being able to subdivide your own piece without asking permission or coordinating bit-for-bit with another team — is the entire organizational payoff subnetting was introduced to provide back in Section 1, now demonstrated at two nested levels simultaneously. It is also, not coincidentally, exactly what makes it possible for someone upstream (an ISP, or this company's own core router) to announce just one route, `10.1.0.0/22`, and have it correctly cover all three (or four) sites — the aggregation idea Chapter 39 names and Chapter 50 fully formalizes.

---

## 8. Worked Problem 4: Reverse-Engineering a Subnet From a Given Address

The first three problems were all *design* problems: given requirements, choose subnets. The mirror-image problem — just as common in real troubleshooting — is: given one address and mask you found in a config file or a `ping` error message, figure out what subnet it belongs to, and what else lives in that subnet.

**Given:** `172.20.133.57/21`. What is this host's network address, broadcast address, and usable range?

**Step 1 — derive the mask from the prefix length.** `/21` = 21 network bits. Two full octets (16 bits) plus 5 more bits into the third octet: `11111111.11111111.11111000.00000000`. The third octet's mask byte, `11111000`, is `128+64+32+16+8 = 248`. So the mask is `255.255.248.0`.

**Step 2 — find the block size in the octet where the split falls.** The split is inside the third octet, with 5 bits fixed and 3 bits free (8−5=3 host bits in that octet, contributing to a larger total host portion that also includes all 8 bits of the fourth octet: 3+8 = 11 host bits total, and 32−21=11 — confirmed). The block size in the third octet is 2^3 = 8 — meaning valid subnet boundaries in the third octet fall every 8 units: `0, 8, 16, 24, ..., 128, 136, ..., 248`.

**Step 3 — locate where 133 (the given third octet) falls.** `133 ÷ 8 = 16.625` → the block starts at `16 × 8 = 128` (round down to the nearest multiple of 8). So the third octet of the network address is `128`.

**Step 4 — state the full network and broadcast addresses.** Since the block spans third-octet values `128` through `128+8-1=135`, and *every* value of the fourth octet within that span:

```
 Network address:    172.20.128.0
 Broadcast address:  172.20.135.255
 Usable range:        172.20.128.1  –  172.20.135.254
```

**Step 5 — sanity check.** Is `172.20.133.57` actually inside that range? Third octet `133` is between `128` and `135` — yes. Fourth octet `57` is a normal usable value, not `0` or `255` of the *overall* range boundary (the network/broadcast addresses here are `172.20.128.0` and `172.20.135.255` specifically — not every `.0` or `.255` fourth-octet value, since the block spans multiple third-octet values). Confirmed: `172.20.133.57/21` belongs to the network `172.20.128.0/21`, sharing that network with every other address from `172.20.128.1` through `172.20.135.254`.

This "find the nearest lower multiple of the block size" technique is the general-purpose reverse-engineering method — it works identically for a boundary that falls anywhere in any octet, and it's exactly the mental math a network engineer does when staring at a `show ip interface` output or an error log with no calculator handy.

---

## 9. A Hands-On Experiment

Pick any two of the subnets computed above and verify them independently using a subnet calculator, forcing yourself to explain any mismatch by finding the arithmetic slip:

```bash
$ ipcalc 192.168.20.128/26
Address:   192.168.20.128
Netmask:   255.255.255.192 = 26
Network:   192.168.20.128/26
HostMin:   192.168.20.129
HostMax:   192.168.20.190
Broadcast: 192.168.20.191

$ ipcalc 172.20.133.57/21
Address:   172.20.133.57
Netmask:   255.255.248.0 = 21
Network:   172.20.128.0/21
HostMin:   172.20.128.1
HostMax:   172.20.135.254
Broadcast: 172.20.135.255
```

Both match Sections 6 and 8 exactly. Then try one more, cold, with no worked example to check against first: what network does `10.5.9.200/23` belong to? Work it by hand (a `/23` has 9 host bits, block size 2^1=2 in the second-to-last relevant octet position — actually here the split is one bit into the third octet, giving block size 2 in the third octet: boundaries at every even number, 8 and 9 pair together), then confirm with `ipcalc`.

---

## 10. Common Misconceptions and Mistakes

- **"Subnetting means memorizing which mask goes with which host count."** That's a symptom, not the skill. Every mask-to-host-count mapping in this chapter was *derived* on the spot from "how many bits does it take to represent N values" (Section 3) — a table of memorized values breaks the moment a problem doesn't match a row you memorized, which is exactly why Problems 2–4 above were solved from the underlying bit arithmetic instead.
- **"All subnets from one allocation must be the same size."** Disproven directly by Worked Problem 2 (VLSM) — subnets from the same parent block can, and in real deployments usually do, have different mask lengths, chosen independently per subnet's actual need.
- **"You can just round the host count up to the nearest occupied address and ignore the network/broadcast overhead."** Every worked example above explicitly subtracted 2 from 2^H to get the *usable* count (Section 4) — forgetting this off-by-two is the single most common subnetting mistake in practice, and it matters most exactly at the small end (a /30 "for 4 hosts" only actually holds 2).
- **"A /21 or /22 splits cleanly on an octet boundary, like /24 does."** Worked Problems 3 and 4 both deliberately used non-octet-aligned prefixes (/22, /21) specifically to break this assumption — the split can fall in the middle of any octet, and only the binary block-size arithmetic (Section 8, Steps 2–3) locates it correctly.
- **Historical note — "subnet zero" and the "all-ones subnet":** older equipment and standards once forbade using the first subnet (all borrowed bits 0) and last subnet (all borrowed bits 1) of a block, out of concern they could be confused with the parent network's own network/broadcast address. Modern equipment (and RFC 1878) permits both; the worked examples in this chapter use them freely, matching current real-world practice, but you may still encounter this restriction mentioned in older documentation or legacy hardware.

---

## 11. What's Simplified Here

This chapter assumed you already had a clean address block to subdivide — it did not cover how that block was originally decided in size, or the historical rigid class system that once constrained what block sizes were even possible to request. That history, and the classless (CIDR) fix that makes flexible block sizes like `192.168.20.0/22` requestable in the first place, is Chapter 39. This chapter also didn't cover how subnetted networks get advertised to routers outside the organization, or how those routers combine many small subnets back into fewer, larger routing table entries — that's route aggregation, previewed at the end of Chapter 39 and covered fully in Chapter 50.

---

## 12. Interview Questions & Model Answers

**Beginner: "What problem does subnetting solve?"**

An organization is given one block of IP addresses but needs multiple separate, isolated networks (for different departments, buildings, or security zones) without requesting additional address blocks and without wasting the addresses it already has. Subnetting solves this by moving the network/host boundary within the existing block, trading host-address capacity for additional distinct network segments.

**Beginner: "If you borrow 3 bits from the host portion of a /24 network, how many subnets do you get, and how large is each one?"**

Borrowing 3 bits gives 2^3 = 8 subnets. The new prefix length is /27 (24+3), leaving 5 host bits, so each subnet has 2^5 = 32 addresses, of which 30 are usable (32 minus the network and broadcast addresses).

**Intermediate: "What is VLSM, and why is it usually better than giving every subnet the same size?"**

Variable Length Subnet Masking allows different subnets carved from the same parent block to use different prefix lengths, each sized to that subnet's actual host requirement. Giving every subnet the same (largest-needed) size wastes huge numbers of addresses on subnets that need far fewer hosts — for example, a router-to-router link only ever needs 2 addresses, and forcing it into a uniform /26 (62 usable) wastes 60 of them, whereas VLSM assigns it a /30 (2 usable) instead.

**Intermediate: "Why does the standard VLSM allocation strategy allocate the largest subnets first?"**

Because allocating from the front of the available space largest-first avoids fragmenting the remaining space in a way that could leave no contiguous block large enough for a still-unallocated large requirement. If small subnets were placed first, at arbitrary positions, the leftover space could end up split into pieces too small individually to satisfy a later, larger requirement even though enough total addresses remain.

**Advanced: "Given `10.20.75.0/22`, determine whether `10.20.76.14` is on the same network. Show your work."**

A /22 fixes 22 bits as network: the first 16 (both full octets) plus 6 more bits into the third octet, leaving that octet's last 2 bits (plus all of the fourth octet) as host bits. With 2 free bits in the third octet, the block boundaries there fall every 2^2=4 units: 0, 4, 8, ..., 72, 76, 80. `10.20.75.0`'s third octet, 75, falls between 72 and 75 (75÷4=18.75, floor to 18×4=72), so its true network is `10.20.72.0/22`, covering third-octet values 72–75. `10.20.76.14`'s third octet, 76, is itself an exact multiple of 4, so it starts the *next* block, `10.20.76.0/22`, covering 76–79. Since `10.20.72.0/22` and `10.20.76.0/22` are different networks, `10.20.76.14` is NOT on the same network as an address in `10.20.75.0/22`.

**Advanced: "Explain, from first principles (not a formula), why a /30 subnet provides exactly 2 usable host addresses, and why this makes it the conventional choice for router-to-router links."**

A /30 has 32-30=2 host bits, giving 2^2=4 total addresses. Of those 4, one is always the network address (all host bits 0) and one is always the broadcast address (all host bits 1) — mandatory reservations for any subnet, regardless of size (Section 4). That leaves exactly 4-2=2 usable addresses. A point-to-point link between two routers has exactly two endpoints and needs exactly two addresses, no more — using any larger subnet would waste addresses on a link that structurally can never have more than two hosts, so a /30 (or, in more modern practice, sometimes a /31, which repurposes both remaining addresses for point-to-point use per RFC 3021) is the natural, minimum-waste fit.

---

## 13. Exercises

### Easy

1. You are given `192.168.50.0/24` and need exactly 2 equal-sized subnets. What prefix length do you use, and what are the two resulting network addresses?
2. How many usable host addresses does a `/29` subnet provide? Show the 2^H − 2 computation.
3. Given `10.0.0.0/25`, state its network address, broadcast address, and usable host range.

### Medium

4. An office needs 3 subnets from `172.16.0.0/24`: one for 100 hosts, one for 40 hosts, one for 10 hosts. Using VLSM (largest first), determine each subnet's prefix length, network address, broadcast address, and usable range, plus how much address space (if any) is left unallocated.
5. A technician finds a device configured as `192.168.30.190/26`. Determine the network address, broadcast address, and usable range for this device's subnet, showing the binary AND for the last octet.
6. Explain why `10.1.4.0/22` and `10.1.5.0/22` describe the exact same network, using binary to justify it (hint: check whether 4 and 5 fall in the same block-size-4 boundary).

### Hard

7. A company has `10.50.0.0/20` and three regional offices, each needing a several-hundred-address block for further internal subnetting, plus room for 5 more future offices carved from that same original block (8 equal regions total). Design the top-level split (how many bits to borrow, resulting prefix length, and the network address of each of the 8 regional blocks), then, for just the first region, further subnet it via VLSM into a 60-host subnet and a 25-host subnet.
8. Given only the address `192.168.77.201` and the mask `255.255.255.224`, determine: the prefix length, the network address, the broadcast address, the usable host range, and whether `192.168.77.190` is on the same subnet. Show every step in binary — do not skip to the answer.

---

## Production Usage Notes

Real organizations rarely subnet by hand at the scale of Worked Problem 3 more than once — they use IP Address Management (IPAM) systems (like NetBox, phpIPAM, or a cloud provider's own VPC subnet planner) that implement exactly the largest-first VLSM algorithm from Section 6, tracked against a database so two teams never accidentally allocate overlapping blocks. But the underlying arithmetic these tools run is never anything more exotic than what this chapter worked by hand: this is precisely why understanding the derivation matters even though you'll rarely do Worked Problem 2's long division in a spreadsheet by hand at work — when the tool produces a surprising or seemingly-wrong answer (and it does happen — a stale IPAM record, a fat-fingered prefix length), the only way to know whether the tool or the record is wrong is to be able to redo the AND-and-block-boundary math yourself, exactly as Worked Problem 4 modeled. It's also worth planning deliberate spare capacity, as every worked problem in this chapter did (Problem 2 left 28 addresses unused, Problem 3 reserved a whole future site) — renumbering a live, in-production subnet later is disruptive in a way that reserving a little headroom up front is not.

---

## 14. Summary

| Term | Meaning |
|---|---|
| Subnetting | Dividing one address block into smaller, independently manageable networks by moving the network/host boundary |
| Borrowing bits | Reclassifying host bits as network bits, multiplying the number of subnets by 2^N and dividing each subnet's size by 2^N |
| Network address | All host bits = 0; represents a subnet as a whole (never assignable to a host) |
| Broadcast address | All host bits = 1; reaches every host on that subnet (never assignable to a host) |
| Usable hosts | 2^H − 2, where H is the number of host bits in that subnet's mask |
| VLSM (Variable Length Subnet Masking) | Using different prefix lengths for different subnets carved from the same parent block, sized to each one's actual need |
| Hierarchical subnetting | Subnetting in layers — dividing a block among sites, then dividing each site's piece again internally |
| /30 (or /31) | The conventional minimum-size subnet for a two-router point-to-point link |

Subnetting solved the problem of dividing a fixed block efficiently — but it assumed you already had a flexible, arbitrary-length prefix (a /26, a /21, a /22) available to request and to divide in the first place. That flexibility itself had to be invented; the original IP addressing system only offered three rigid sizes, and that rigidity wasted addresses on a massive scale. Chapter 39 tells that story and introduces CIDR — the classless addressing system that made every subnet mask in this chapter a legal, requestable thing in real-world routing, before previewing route aggregation, subnetting's mirror-image operation, which Chapter 50 covers in full.
