# Chapter 46: Static Routing vs. Dynamic Routing

> **"A static route is a promise a human made about the network, frozen at the moment they typed it. The network keeps changing after that. Dynamic routing exists because promises don't survive contact with a cut fiber cable."**

---

## Table of Contents

1. [The Problem: Someone Has to Fill In the Table](#1-the-problem-someone-has-to-fill-in-the-table)
2. [Static Routing — The Simplest Possible Answer](#2-static-routing--the-simplest-possible-answer)
3. [A Worked Static Topology](#3-a-worked-static-topology)
4. [Where Static Routing Breaks: A Link Failure](#4-where-static-routing-breaks-a-link-failure)
5. [Why This Gets Worse, Not Better, With Scale](#5-why-this-gets-worse-not-better-with-scale)
6. [Dynamic Routing — Let Routers Talk to Each Other](#6-dynamic-routing--let-routers-talk-to-each-other)
7. [The Two Families of Dynamic Routing Protocols](#7-the-two-families-of-dynamic-routing-protocols)
8. [Administrative Distance — Choosing Between Sources](#8-administrative-distance--choosing-between-sources)
9. [IGP vs. EGP](#9-igp-vs-egp)
10. [Where Static Routing Still Wins in Production](#10-where-static-routing-still-wins-in-production)
11. [Hybrid Designs: Redistribution and Floating Static Routes](#11-hybrid-designs-redistribution-and-floating-static-routes)
12. [Code: A Tiny Static Routing Table Loader](#12-code-a-tiny-static-routing-table-loader)
13. [What's Simplified Here](#13-whats-simplified-here)
14. [Common Misconceptions](#14-common-misconceptions)
15. [Hands-On Experiment](#15-hands-on-experiment)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#summary)

---

## 1. The Problem: Someone Has to Fill In the Table

Chapters 44 and 45 treated the routing table as a given — a list of rows that already exists, ready for the longest-prefix-match algorithm to consult. But every one of those rows had to get there somehow. Chapter 44, Section 7 named the three sources: connected routes (free, automatic, but only cover networks a router is directly attached to), static routes, and dynamic routes. Connected routes alone are nowhere near enough — they only tell a router about its own immediate neighbors, and the entire point of a router is to reach networks it is *not* directly attached to. Something has to supply every other row.

The question this chapter answers: **who — or what — decides that "to reach network X, go via router Y," and writes that decision into every affected routing table?**

---

## 2. Static Routing — The Simplest Possible Answer

The most direct possible answer: a human works out the network's topology on paper (or in their head), decides the correct path to every network that matters, and types the corresponding route into every router's configuration by hand. This is **static routing** — Chapter 44's example (`198.51.100.0/24 via 203.0.113.1`) was exactly this: someone, at some point, sat down and configured that line.

A static route is defined by three properties, all following directly from how it's created:

- **It never changes on its own.** A static route stays in the routing table exactly as configured, permanently, until a human edits or removes it — even if the path it describes stops working.
- **It requires exact, correct topology knowledge at configuration time.** Whoever wrote it needed to already know which next hop is correct.
- **It consumes no protocol overhead.** No router sends any messages to any other router to maintain a static route — it's just a line in a config file, evaluated the same way every connected or dynamic route is evaluated (Chapter 45's LPM algorithm doesn't care where a route came from).

That last property is genuinely valuable, not just a consolation prize: static routes are completely predictable, consume no bandwidth or CPU keeping themselves updated, and are trivial to audit — reading the configuration file *is* reading the exact behavior, with no protocol state, timers, or convergence process to reason about. This is precisely why, as Section 10 covers, static routing never actually disappeared even after dynamic protocols existed — it remains the right tool for genuinely simple or security-sensitive cases.

---

## 3. A Worked Static Topology

Consider three small offices connected in a line, each with its own router and LAN:

```
   Office A LAN          Office B LAN          Office C LAN
  10.1.0.0/24           10.2.0.0/24           10.3.0.0/24
       |                     |                     |
    [Router A]===========[Router B]===========[Router C]
       |    10.12.0.0/30      |    10.23.0.0/30     |
   .1                     .1  .2                .1  .2
```

For every router to reach every LAN, someone must configure:

**Router A's table:**
```
10.1.0.0/24    connected, via LAN interface
10.12.0.0/30   connected, via link to B
10.2.0.0/24    via 10.12.0.2   (Router B)      <- static, hand-typed
10.3.0.0/24    via 10.12.0.2   (Router B)      <- static, hand-typed
```

**Router B's table:**
```
10.2.0.0/24    connected, via LAN interface
10.12.0.0/30   connected, via link to A
10.23.0.0/30   connected, via link to C
10.1.0.0/24    via 10.12.0.1   (Router A)      <- static, hand-typed
10.3.0.0/24    via 10.23.0.2   (Router C)      <- static, hand-typed
```

**Router C's table:**
```
10.3.0.0/24    connected, via LAN interface
10.23.0.0/30   connected, via link to B
10.1.0.0/24    via 10.23.0.1   (Router B)      <- static, hand-typed
10.2.0.0/24    via 10.23.0.1   (Router B)      <- static, hand-typed
```

Notice Router A's route to `10.3.0.0/24` (Office C's LAN): Router A isn't directly connected to Router C at all, but its static route correctly points at Router B, trusting that Router B knows how to get the rest of the way. This is exactly Chapter 44's delegated-knowledge principle in action — it works perfectly, and for exactly six routes across three small offices, it took only a few minutes to design and type. Six routes, three routers: entirely manageable.

---

## 4. Where Static Routing Breaks: A Link Failure

Now suppose the physical link between Router A and Router B fails — a cut cable, a dead port, an ISP outage on that leased line. What happens to a packet from Office A destined for Office C?

Router A's routing table still contains the exact same line it always did:
```
10.3.0.0/24    via 10.12.0.2   (Router B)
```

Nothing about a static route causes it to notice the link is down and remove or update itself — that is not what "static" not changing means literally, but it is the practical failure mode: unless Router A's operating system independently detects that the *interface* used by that route's path is down (many implementations do withdraw a static route automatically if its next hop becomes unreachable via a directly connected interface — a detail Section 13 revisits), the route can remain listed as usable long after it no longer is. Even in the best case where the OS does correctly notice the interface is down and stops using that specific route, there is **no other route in Router A's table that could replace it** — nobody configured a backup path, because the topology only offered one path in the first place, and even where alternate physical paths exist, nobody pre-typed a fallback route for this specific failure. The result: packets for Office C are dropped (Router A replies with an ICMP Destination Unreachable, if it detects the failure at all) until a human notices the outage, diagnoses it, and manually reconfigures a working route — which could be minutes, hours, or, for a route nobody's actively monitoring, days.

This is the structural problem, stated plainly: **a static route is a snapshot of what was true when it was configured, with no mechanism to notice or react when that stops being true.** The network's actual condition and the router's belief about the network's condition can silently diverge, for an unbounded length of time, until a human intervenes.

---

## 5. Why This Gets Worse, Not Better, With Scale

Three offices needed six hand-typed routes. It's worth working out how this scales, because the answer is the entire motivation for everything in Chapters 47–48.

For **N routers arranged so every router needs a route to every other router's networks** (the worst case, but a realistic shape for anything beyond a simple line or star), the number of static route entries needed grows roughly with N² — each of N routers potentially needs on the order of N routes (one per other network), so the total configuration burden across the whole network is on the order of N × N. Ten routers: on the order of 90 route entries to configure and keep synchronized by hand, scattered across ten separate configuration files that a human has to keep mentally consistent. A hundred routers: on the order of 9,900. And every one of those entries has the fragility from Section 4 baked in — a single topology change (a new subnet, a decommissioned link, a network renumbered) potentially requires editing dozens of routers' configurations simultaneously, correctly, with no tooling flagging an inconsistency until traffic actually breaks and someone notices.

This is, not coincidentally, structurally the exact same N² wall Chapter 03 raised about full-mesh wiring at the very start of this entire course: manually managing every pairwise relationship in a growing system is a shape of problem that always eventually breaks, whichever layer of the stack it shows up at. The Internet's actual scale — hundreds of thousands of networks — makes pure static routing not just inconvenient but categorically impossible to operate by hand. Something has to let routers discover and adapt to the network's real topology themselves.

---

## 6. Dynamic Routing — Let Routers Talk to Each Other

**The core idea:** instead of a human working out every path and typing it in, routers run a **routing protocol** — they exchange messages directly with their neighbors, describing what they know about the network, and each router uses that information to compute its own routing table entries automatically. When the topology changes — a link fails, a new network is added, a router reboots — the affected routers detect it and propagate updated information to their neighbors, who update their own tables in turn, without any human editing a configuration file.

This directly answers both of Section 4 and 5's problems:

- **Failure recovery becomes automatic.** If the A–B link in Section 3's topology fails and an alternate physical path exists anywhere in the network, a dynamic routing protocol will detect the failure (typically within seconds, via missed periodic messages or an interface-down notification) and recompute a working route through the surviving path — no human required.
- **Configuration burden stops scaling with N².** Each router only needs to be told about its own directly connected networks and which neighbors to exchange routing information with — a burden that scales with the router's own number of interfaces, not with the size of the entire network. The protocol itself handles propagating that information everywhere it needs to go.

The trade-off, previewed here and explored fully in Chapters 47–48: dynamic routing protocols consume real bandwidth and CPU maintaining themselves (routers are constantly talking to each other, even when nothing has changed), take a nonzero amount of time to detect and fully react to a change (**convergence time** — the gap during which some routers may have stale or even temporarily incorrect information), and are a meaningfully larger, more complex thing to configure correctly and secure than a static route. Nothing about dynamic routing is "free" — it trades one set of problems (manual toil, slow failure recovery) for a different, more manageable set (protocol overhead, convergence delay, configuration complexity of the protocol itself).

---

## 7. The Two Families of Dynamic Routing Protocols

Every interior dynamic routing protocol you'll meet in this volume belongs to one of two families, previewed here and covered in full in the next two chapters:

| | Distance-vector | Link-state |
|---|---|---|
| What a router shares with neighbors | "Here's my distance to every destination I know" — a summarized opinion, not raw topology | "Here are my direct links and their costs" — raw, unprocessed facts about itself only |
| What each router builds | A table of distances, trusting neighbors' summaries | A complete map (link-state database) of the *entire* network's topology |
| How the best path is computed | Accumulated hop-by-hop from neighbor gossip (Bellman-Ford-style) | Independently, by each router running Dijkstra's shortest-path algorithm on its own full map |
| Convergence speed | Slow — often tied to periodic update timers | Fast — triggered immediately by any topology change |
| Example protocol | **RIP** (Chapter 47) | **OSPF, IS-IS** (Chapter 48) |

Distance-vector protocols are conceptually simpler and were developed first — you'll see in Chapter 47 exactly why "trust your neighbor's summarized opinion" is both easy to implement and prone to a specific, serious failure mode (count-to-infinity) that link-state protocols were designed specifically to avoid.

---

## 8. Administrative Distance — Choosing Between Sources

Chapter 44 flagged this and Chapter 45 relied on it as a tie-breaker; now it can be defined properly, because it only makes sense once you know static and dynamic routes coexist in the same table.

A single router can simultaneously have a static route, a RIP-learned route, and an OSPF-learned route all describing the exact same destination prefix — perhaps an administrator configured a static backup route "just in case" a dynamic protocol's route disappears. When two or more routes tie on prefix length (Chapter 45, Section 13), the router needs a rule for which *source* to trust more. That rule is **administrative distance (AD)** — a per-source trust ranking, independent of the routing protocol's own internal metric, configured (with sensible vendor defaults) on the router itself:

| Route source | Typical default administrative distance (lower = more trusted) |
|---|---|
| Directly connected | 0 |
| Static route | 1 |
| eBGP (external BGP, Chapter 49) | 20 |
| OSPF | 110 |
| IS-IS | 115 |
| RIP | 120 |
| iBGP (internal BGP) | 200 |
| Unreachable / unusable | 255 |

Reading this table: a static route (AD 1) will always be preferred over an OSPF route (AD 110) to the exact same prefix, on the same router, even though OSPF is the more sophisticated protocol — administrative distance says nothing about which protocol is "smarter," only which source the operator (via vendor defaults, which are themselves just sensible starting assumptions) trusts more when there's a genuine tie. This is the exact mechanism Section 11 relies on for "floating static routes" — a deliberate design pattern that only works because static beats OSPF by default.

---

## 9. IGP vs. EGP

One more piece of vocabulary worth placing correctly before Chapters 47–49: dynamic routing protocols split into two categories based on *where* they operate, not how they compute routes.

- An **Interior Gateway Protocol (IGP)** runs *within* a single organization's network — one administrative domain, one entity making all the decisions. RIP, OSPF, and IS-IS are all IGPs. Their job is "help all the routers inside my company/campus/ISP agree on the best path to every network I own."
- An **Exterior Gateway Protocol (EGP)** runs *between* separate organizations — separate autonomous systems (Chapter 50) that don't share administrative control and, critically, don't necessarily trust each other or agree on what "best path" even means. BGP (Chapter 49) is the only EGP in practical use today.

This distinction matters because it explains why BGP, covered starting next chapter-volume-section, looks so different from RIP and OSPF: an IGP can assume cooperation (every router inside one company wants the same thing — the fastest, most reliable path) and optimize purely for topology. An EGP cannot assume that at all — two ISPs might have completely different business incentives about which paths they prefer, a subject Chapter 49 opens with directly.

---

## 10. Where Static Routing Still Wins in Production

Dynamic routing did not make static routing obsolete — it's worth being concrete about exactly when static routing remains the right engineering choice, even in modern, well-run networks:

- **Stub networks** — a network with exactly one way in or out (a small office with a single ISP link, a branch with one WAN connection) gains nothing from a dynamic protocol's ability to find alternate paths, because no alternate path exists. A single static default route is simpler, has zero protocol overhead, and cannot possibly compute the "wrong" answer because there's only one answer.
- **The last-mile default route to an ISP.** Even large enterprise routers commonly use a static default route pointing at their upstream ISP rather than learning it dynamically, because that relationship is stable, contractual, and doesn't benefit from automated discovery.
- **Security-sensitive or deliberately isolated segments**, where an administrator wants to guarantee, by inspection of a config file, that traffic to a sensitive network takes one specific, audited path and can never silently shift onto a different path a routing protocol might select on its own.
- **Explicit traffic engineering** — occasionally an operator wants to force traffic down a specific link for cost, contractual, or performance reasons that a routing protocol's automatic metric wouldn't naturally choose; a well-placed static route (or an even more surgical policy-based route) achieves that deterministically.
- **Small, genuinely static topologies** — Section 3's three-office example, if the topology truly never changes, isn't wrong to run statically; the case *against* static routing is a case about scale and volatility, not a blanket rule.

The honest framing: static routing is the right default when a network is small, its topology is stable, and simplicity/predictability outweigh the cost of manual maintenance. Dynamic routing is the right default the moment any of those three stop being true.

---

## 11. Hybrid Designs: Redistribution and Floating Static Routes

Real production networks are rarely purely static or purely dynamic — they blend both deliberately:

- **Redistribution** is configuring a router to inject routes learned from one source (say, static routes, or one routing protocol) into a different dynamic protocol's advertisements, so the rest of the network learns about them automatically. A common pattern: statically configure a route to a small partner network on one edge router, then redistribute it into OSPF so every other router in the company learns about it without needing its own static entry.
- **A floating static route** is a static route deliberately configured with a *worse* (higher) administrative distance than the dynamic protocol already in use — say, AD 130, worse than OSPF's 110. Under normal conditions, the router always prefers the OSPF-learned route (lower AD wins ties, Section 8). But if the OSPF route to that destination ever disappears entirely (the dynamic protocol has no path at all — not just a worse one, an absent one), the floating static route becomes the *only* remaining option and takes over automatically as a backup, with no human intervention required. This is a deliberately engineered combination of static routing's predictability (you know exactly what the backup path is) and dynamic routing's automatic failover (no one has to notice the failure and switch it in by hand).

---

## 12. Code: A Tiny Static Routing Table Loader

A simplified but realistic illustration of what "static routing" looks like as actual machine behavior — a program that reads a small static configuration file and installs the routes, with no protocol exchange of any kind:

```go
package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strings"
)

type StaticRoute struct {
	Network *net.IPNet
	NextHop net.IP
}

// loadStaticRoutes reads lines like:
//   10.2.0.0/24 10.12.0.2
// and returns them as routes. No neighbor communication, no timers,
// no reaction to failure -- exactly what makes static routing both
// simple and, per Section 4, fragile.
func loadStaticRoutes(path string) ([]StaticRoute, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	var routes []StaticRoute
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		fields := strings.Fields(scanner.Text())
		if len(fields) != 2 {
			continue
		}
		_, network, err := net.ParseCIDR(fields[0])
		if err != nil {
			continue
		}
		routes = append(routes, StaticRoute{
			Network: network,
			NextHop: net.ParseIP(fields[1]),
		})
	}
	return routes, scanner.Err()
}

func main() {
	routes, err := loadStaticRoutes("routes.conf")
	if err != nil {
		fmt.Println("error:", err)
		return
	}
	for _, r := range routes {
		fmt.Printf("static route: %s via %s\n", r.Network, r.NextHop)
	}
	// Nothing in this program ever re-reads routes.conf on its own,
	// pings the next hop, or notices a failed link. That's the whole
	// static-routing bargain: total predictability, zero self-healing.
}
```

Contrast this mentally with Chapter 47's RIP implementation sketch, which will instead *listen* for neighbor messages and *rebuild* its table continuously — the structural difference between static and dynamic routing is really the difference between "read a file once" and "run an ongoing conversation."

---

## 13. What's Simplified Here

- Many real operating systems and router platforms *do* automatically withdraw a static route if its directly-attached next-hop interface goes down (this is sometimes called route tracking or interface-dependent static routes) — Section 4 slightly overstated static routing's rigidity for clarity; the deeper, unavoidable problem is that even when a static route is correctly withdrawn, nothing automatically supplies a *working alternative* unless a human pre-configured one.
- Administrative distance values in Section 8's table are Cisco's well-known defaults; other vendors use different numeric scales for the same underlying idea, and AD itself is not part of any standard protocol specification — it's a per-implementation convention for resolving multi-source ties.
- Real convergence time for dynamic protocols depends heavily on specific timer tuning, network size, and protocol (RIP's default timers make it converge in tens of seconds to minutes; OSPF can converge in under a second with tuned timers) — Chapters 47 and 48 give protocol-specific numbers rather than the generic "dynamic is faster" claim made here.

---

## 14. Common Misconceptions

- **"Dynamic routing is strictly better, so static routing is legacy/deprecated."** False — Section 10 lists concrete, current production cases where static is the correct engineering choice, not a compromise.
- **"A static route, once configured, is guaranteed to keep working."** No — it's guaranteed to keep *existing* in the table (or be withdrawn if the platform tracks interface state); it says nothing about whether the path it describes is still physically functional unless something else is actively verifying that.
- **"Dynamic routing protocols automatically find the best possible path with no downsides."** They find the best path *according to their own metric* (hop count for RIP, cost for OSPF) once fully converged — but convergence takes time, and during that window, routes can be temporarily wrong, looped, or simply missing.
- **"Administrative distance and a routing protocol's metric are the same number."** They are two entirely separate values serving different purposes — AD ranks *sources*, metric ranks *paths within one source's own view*. See Section 8.

---

## 15. Hands-On Experiment

On a Linux machine, add and inspect a static route directly, then simulate a "failure" by removing the interface it depends on:

```
$ sudo ip route add 198.51.100.0/24 via 203.0.113.1 dev eth0
$ ip route show | grep 198.51.100.0
198.51.100.0/24 via 203.0.113.1 dev eth0

# Simulate the link going down:
$ sudo ip link set eth0 down
$ ip route show | grep 198.51.100.0
# On most systems, the route either disappears or is marked unusable --
# but nothing replaces it, because you (the human) are the only
# mechanism that knows an alternative even exists.

$ sudo ip link set eth0 up   # restore it
```

This small experiment is Section 4's failure scenario, reproduced on real hardware you can hold in your hands (or a VM you're running right now): the route's disappearance is automatic, but recovery is entirely on you.

---

## 16. Interview Questions & Model Answers

**Beginner: What is the main disadvantage of static routing?**
It requires a human to manually configure every route, and it does not automatically adapt when the network's topology changes (a link failure, a new subnet) — traffic can be silently dropped until someone notices and manually reconfigures the affected routers.

**Intermediate: Why does static routing's configuration burden grow roughly with the square of the number of routers, in a fully-meshed requirement scenario?**
Because in the worst case, each of N routers may need a route to every other router's network — roughly N destinations per router, across N routers, giving on the order of N × N = N² total route entries to configure and keep consistent by hand.

**Intermediate: What is administrative distance, and why is it necessary?**
It's a per-source trust ranking (connected, static, OSPF, RIP, BGP, etc.) that a router uses to decide which route to actually install when two or more sources offer routes to the exact same prefix (same length) — it resolves a tie that longest-prefix-match itself cannot resolve, since LPM only distinguishes routes by specificity, not by source.

**Advanced: Design a floating static route as a backup for an OSPF-learned default route, and explain precisely under what condition it activates.**
Configure a static default route (`0.0.0.0/0` via some backup next hop) with an administrative distance deliberately set higher than OSPF's default of 110 — e.g., 200. Under normal operation, the OSPF-learned default route (AD 110) is preferred and installed. The static route only becomes active if the OSPF route disappears from the table entirely — not merely "gets worse," but is completely withdrawn (e.g., OSPF neighbor relationship lost) — at which point the static route is the only remaining candidate and is installed automatically, with no human action.

---

## 17. Exercises

### Easy
1. List two production scenarios where static routing is still the correct choice today, and explain why in each case.
2. Explain, in one or two sentences, the core difference between how a static route and a dynamic route enter a routing table.
3. What does "administrative distance" resolve that longest-prefix-match cannot?

### Medium
4. For a network of 5 routers arranged so every router needs a route to every other router's LAN, estimate (using Section 5's reasoning) roughly how many static route entries would need to be configured in total.
5. Explain the difference between an IGP and an EGP, and name one example of each.
6. A static route and an OSPF route both exist for the exact same `/24` prefix on one router. Explain which one is installed, and why the answer doesn't depend on either route's individual metric.

### Hard
7. Design (on paper) a floating static route setup for the three-office topology in Section 3, assuming Router A and Router B are now connected by *two* physical links (the original, plus a new backup link with worse bandwidth). Explain how you'd use administrative distance so that OSPF would need to be introduced for full automatic failover — and what would still be missing if you tried to solve this with static routes alone.
8. Explain, using Section 4's failure scenario, precisely what "convergence" would need to mean if this network of three offices were instead running a dynamic routing protocol, and why "convergence time" is a genuinely new concept that has no equivalent at all under pure static routing.
9. A network engineer redistributes a static route into OSPF (Section 11). Explain what could go wrong if that static route were actually incorrect (pointing at a next hop that doesn't lead where the prefix claims), and why redistribution can propagate a single human mistake much further and much faster than a pure static-routing mistake would.

---

## Summary

| Term | Meaning |
|---|---|
| Static route | A manually configured, unchanging routing table entry |
| Dynamic route | A routing table entry learned automatically via a routing protocol exchanging information with neighbors |
| Convergence | The process (and time taken) for all routers in a network to agree on correct routes after a topology change |
| Distance-vector | A dynamic routing family where routers share summarized distance opinions with neighbors (e.g., RIP) |
| Link-state | A dynamic routing family where routers share raw topology facts and independently compute shortest paths (e.g., OSPF, IS-IS) |
| Administrative distance (AD) | A per-source trust ranking used to break ties when multiple sources offer a route to the exact same prefix |
| IGP | Interior Gateway Protocol — routes within one administrative domain (RIP, OSPF, IS-IS) |
| EGP | Exterior Gateway Protocol — routes between separate administrative domains (BGP) |
| Floating static route | A static route deliberately given a worse AD than a dynamic protocol, acting as an automatic backup if the dynamic route disappears |
| Redistribution | Injecting routes from one source (e.g., static) into a dynamic protocol so they propagate automatically |

Static routing answers "who fills in the table" with "a human, once." Chapter 47 introduces the first automatic answer — RIP, the simplest dynamic routing protocol — and shows exactly how its simplicity comes with a serious, well-documented failure mode of its own.
