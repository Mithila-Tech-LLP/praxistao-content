# Chapter 48: OSPF and IS-IS — Link State Routing

> **"RIP asks your neighbors what they believe. OSPF asks your neighbors what they can see, gets the same raw facts everyone else gets, and then does its own math. Trusting facts instead of conclusions is the entire difference between a protocol that can loop for minutes and one that can't loop at all."**

---

## Table of Contents

1. [The Problem RIP Leaves Open](#1-the-problem-rip-leaves-open)
2. [The Link-State Idea](#2-the-link-state-idea)
3. [OSPF Mechanics: Hello, Neighbors, and the LSDB](#3-ospf-mechanics-hello-neighbors-and-the-lsdb)
4. [Link State Advertisements](#4-link-state-advertisements)
5. [The Designated Router](#5-the-designated-router)
6. [OSPF Cost — A Metric That Reflects Reality](#6-ospf-cost--a-metric-that-reflects-reality)
7. [Dijkstra's Algorithm, Explained From Scratch](#7-dijkstras-algorithm-explained-from-scratch)
8. [A Full Worked Dijkstra Walkthrough](#8-a-full-worked-dijkstra-walkthrough)
9. [Why Link-State Converges Faster and Doesn't Count to Infinity](#9-why-link-state-converges-faster-and-doesnt-count-to-infinity)
10. [OSPF Areas — Scaling the Link-State Idea](#10-ospf-areas--scaling-the-link-state-idea)
11. [IS-IS — The ISP-Favored Cousin](#11-is-is--the-isp-favored-cousin)
12. [OSPF vs. IS-IS vs. RIP — A Direct Comparison](#12-ospf-vs-is-is-vs-rip--a-direct-comparison)
13. [Code: A Minimal Dijkstra Implementation](#13-code-a-minimal-dijkstra-implementation)
14. [A Real OSPF Neighbor Table](#14-a-real-ospf-neighbor-table)
15. [What's Simplified Here](#15-whats-simplified-here)
16. [Common Misconceptions](#16-common-misconceptions)
17. [Hands-On Experiment](#17-hands-on-experiment)
18. [Interview Questions & Model Answers](#18-interview-questions--model-answers)
19. [Exercises](#19-exercises)
20. [Summary](#summary)

---

## 1. The Problem RIP Leaves Open

Chapter 47 traced RIP's count-to-infinity failure down to one root cause: **a RIP router never sees the network's actual topology — only ever-accumulating secondhand distance opinions.** When the fact behind those opinions changes, there's no way for a router to tell "genuinely new information" apart from "an old belief echoing back." Split horizon patches the most common case; the 15-hop ceiling bounds the worst case; but nothing about distance-vector routing, even fully mitigated, removes the underlying blindness. And separately, RIP's hop-count metric cannot express the difference between a saturated dial-up link and an idle 100 Gbps backbone trunk — both cost exactly 1.

Both problems trace back to the same design choice: **routers exchange conclusions, not facts.** The fix explored in this chapter inverts that choice entirely.

---

## 2. The Link-State Idea

**Intuitively:** instead of asking your neighbors "how far is it to Accounting?" and trusting their answer, imagine every person in the building instead writes down, on one small card, only the doors *directly* connected to their own room and how wide each one is — nothing else, no opinions about the rest of the building. Every card gets copied and handed to literally everyone in the building. Now everyone holds an identical, complete stack of cards — a complete floor plan, assembled entirely from small, first-hand, unopinionated facts nobody could get wrong about anyone else's territory. From that floor plan, *you* work out the shortest path to Accounting yourself; you don't have to trust anyone else's arithmetic, because you did the arithmetic.

**The real-world analogy, and where it breaks:** think of it as building a full subway map by collecting individual "this platform connects to these two tracks" placards from every single platform, rather than asking passengers "how many stops to downtown?" and trusting their memory. Where it breaks: in a real building, handing out and correctly reconciling thousands of little cards, then having everyone independently redo the shortest-path arithmetic, sounds like far more total work than just asking around — this is a legitimate cost, and Section 10 (OSPF areas) exists specifically to bound it as the network grows.

**Engineering definition:** in **link-state routing**, every router:

1. Discovers its directly connected neighbors and the cost of reaching each one,
2. Packages that information (only about itself, never any opinion about the wider network) into a small, signed message — a **Link State Advertisement (LSA)**,
3. **Floods** that LSA to every other router in the relevant area, unmodified, so every router ends up holding an identical copy of every other router's LSA,
4. Assembles all received LSAs into a **Link State Database (LSDB)** — a complete graph of the network's topology,
5. Runs **Dijkstra's shortest-path algorithm** (Section 7) independently, locally, on its own copy of that graph, to compute the best route to every destination.

No router ever asks another router "what's your route to X?" No router ever repeats anyone else's conclusion. Every router does its own, independent, verifiable computation from the same raw, first-hand data. This single change — facts instead of conclusions — is the entire reason link-state routing avoids count-to-infinity and converges dramatically faster, as Section 9 makes precise once Dijkstra's algorithm itself has been introduced.

---

## 3. OSPF Mechanics: Hello, Neighbors, and the LSDB

OSPF (Open Shortest Path First, RFC 2328 for OSPFv2/IPv4) is the dominant link-state IGP inside enterprise networks, data centers, and many ISPs' internal (not inter-AS) routing. Its mechanics, concretely:

- **Hello protocol.** OSPF routers send small Hello packets (every 10 seconds on most link types) to discover neighbors and continuously confirm they're still alive. Two routers become OSPF neighbors only after exchanging Hellos and agreeing on shared parameters (area ID, timers, authentication).
- **Adjacency and database exchange.** After becoming neighbors, routers on point-to-point links (and the Designated Router on broadcast links — Section 5) form a full **adjacency**, exchanging their complete LSDBs so both sides start from an identical baseline.
- **Flooding.** When a router's own link state changes (a link goes up/down, cost changes), it generates a new LSA and floods it to every neighbor, who immediately forward it onward (with loop-prevention via sequence numbers) until every router in the area has the update — typically within milliseconds to a few seconds of the actual change, not tied to any periodic timer the way RIP's 30-second cycle is.
- **Dead interval.** If no Hello is received from a neighbor for 40 seconds (4× the Hello interval, by default), the neighbor is declared down, triggering an immediate LSA flood announcing the change — much faster detection than RIP's 180-second route timeout, though still bounded by this timer rather than being instantaneous.
- **Recomputation.** Every router that receives a topology-changing LSA reruns Dijkstra's algorithm (Section 7) on its updated LSDB to recompute affected routes.

---

## 4. Link State Advertisements

An LSA is deliberately narrow in scope — a router only ever originates LSAs describing *itself*, never anyone else. The most important types:

| LSA Type | Originated by | Describes |
|---|---|---|
| Type 1 (Router LSA) | Every router | The router's own direct links, their costs, and neighbors — the raw building block of the topology graph |
| Type 2 (Network LSA) | The Designated Router (Section 5) on a broadcast segment | Which routers are attached to a shared LAN segment |
| Type 3 (Summary LSA) | Area Border Routers | Routes from one area, summarized, into another area (Section 10) |
| Type 5 (AS-External LSA) | Routers redistributing external routes | Routes learned from outside OSPF entirely (e.g., redistributed static or BGP routes) |

Every LSA carries a sequence number and an age, specifically so that a router receiving two versions of the same LSA (from different flooding paths, potentially arriving out of order) can always tell which is newer and safely discard the stale one — a small but important piece of engineering that prevents old, superseded topology facts from ever being mistaken for current ones, closing off exactly the kind of confusion Chapter 47's Section 6 showed distance-vector routing is vulnerable to.

---

## 5. The Designated Router

On a broadcast-capable segment (an Ethernet LAN with, say, five routers all attached to the same switch), the naive link-state approach would have every router form a full adjacency with every other router on that segment — for 5 routers, that's 10 separate adjacencies, each one a full LSDB exchange and an ongoing relationship to maintain, purely to describe one shared LAN segment that could be represented far more economically.

OSPF's fix: routers on a broadcast segment elect one **Designated Router (DR)** and one backup (**BDR**, for redundancy). Every other router on that segment forms a full adjacency only with the DR and BDR, not with each other — reducing 10 potential adjacencies down to roughly 2×(N−1). The DR is then responsible for originating the Type 2 Network LSA describing "these routers all share this segment," on behalf of everyone. This is purely a scaling optimization for broadcast media; on point-to-point links (exactly two routers, nobody else could possibly be attached), there's no need for a DR at all, since a single adjacency between the two routers already describes the whole segment completely.

---

## 6. OSPF Cost — A Metric That Reflects Reality

Where RIP's metric is a flat "1 per hop" regardless of what that hop actually is, OSPF's metric — **cost** — is explicitly derived from link bandwidth, directly fixing the fiber-vs-dial-up blindness flagged in Chapter 47, Section 9. The standard formula:

```
Cost = Reference Bandwidth / Interface Bandwidth
```

With the common default reference bandwidth of 100 Mbps:

| Link type | Bandwidth | OSPF Cost (default reference 100 Mbps) |
|---|---|---|
| 10 Mbps Ethernet | 10 Mbps | 10 |
| 100 Mbps Fast Ethernet | 100 Mbps | 1 |
| 1 Gbps Ethernet | 1,000 Mbps | 1 (rounds down — a known limitation) |
| 10 Gbps Ethernet | 10,000 Mbps | 1 (same rounding limitation) |

That rounding problem — gigabit and 10-gigabit links both computing to cost 1 under the historical 100 Mbps default — is a real, commonly-hit operational issue; production OSPF deployments on modern networks almost universally raise the reference bandwidth (e.g., to 100 Gbps or higher) via configuration specifically so that cost can still meaningfully distinguish between a 1G, 10G, and 100G link. The total cost of a path is the **sum of costs of every link along it**, and OSPF selects the path with the lowest total cost — this is exactly the quantity Dijkstra's algorithm (Sections 7–8) is built to minimize.

---

## 7. Dijkstra's Algorithm, Explained From Scratch

Edsger Dijkstra published this algorithm in 1959, for a completely general problem: given a graph with weighted edges (costs) and a starting node, find the shortest (lowest total cost) path from that starting node to every other node in the graph. It's a perfect fit for OSPF's needs because an LSDB *is* exactly this kind of graph — routers are nodes, links are edges, and OSPF cost is the edge weight.

**The algorithm, in plain steps:**

1. Assign the starting node (the router doing the computation — call it "self") a tentative distance of 0. Assign every other node a tentative distance of infinity.
2. Mark all nodes as unvisited. Maintain a set of "visited" (finalized) nodes, initially empty.
3. Repeat until every node is visited:
   a. Among unvisited nodes, pick the one with the smallest tentative distance — call it the **current node**.
   b. For each unvisited neighbor of the current node, calculate a candidate distance: (current node's tentative distance) + (cost of the edge to that neighbor). If this candidate is smaller than the neighbor's current tentative distance, update it, and record the current node as that neighbor's next hop toward `self`.
   c. Mark the current node as visited — its tentative distance is now final and will never change again.
4. When every node is visited, every node's tentative distance is the true shortest-path cost from `self`, and the recorded next-hop chain traces out the actual shortest path.

The key insight that makes step 3c valid — and makes this algorithm provably correct, not just a good heuristic — is that once a node has the smallest tentative distance among all remaining unvisited nodes, no future discovery through a not-yet-visited node could possibly produce a shorter path to it (since all edge weights are non-negative, any such alternate path would have to be at least as long as the direct comparison already made). This guarantee is exactly why link-state routing never needs to "wait and see" the way distance-vector routing's slow, uncertain convergence does — Dijkstra's algorithm reaches a mathematically final answer in one pass over the graph, not through iterative rounds of gossip.

---

## 8. A Full Worked Dijkstra Walkthrough

Take a small network of five routers with OSPF costs assigned to each link (using Section 6's cost values as realistic examples):

```
        (2)
    A ------- B
    |  \       |  \
   (1)  (4)   (1)  (3)
    |      \   |     \
    D ------ C -------- E
        (2)        (1)
```

Edges: A–B (2), A–D (1), A–C (4), B–C (1), B–E (3), D–C (2), C–E (1).

Compute the shortest paths **from A** to every other node.

**Initialization:**

| Node | Tentative distance | Visited? | Via |
|---|---|---|---|
| A | 0 | no | — |
| B | ∞ | no | — |
| C | ∞ | no | — |
| D | ∞ | no | — |
| E | ∞ | no | — |

**Step 1 — current node: A (distance 0, smallest unvisited).** Examine A's neighbors: B (edge cost 2), D (edge cost 1), C (edge cost 4).
- B: candidate = 0 + 2 = 2. Better than ∞ → update B to 2, via A.
- D: candidate = 0 + 1 = 1. Better than ∞ → update D to 1, via A.
- C: candidate = 0 + 4 = 4. Better than ∞ → update C to 4, via A.
Mark A visited.

| Node | Tentative distance | Visited? | Via |
|---|---|---|---|
| A | 0 | **yes** | — |
| B | 2 | no | A |
| C | 4 | no | A |
| D | 1 | no | A |
| E | ∞ | no | — |

**Step 2 — current node: D (distance 1, smallest unvisited).** Examine D's neighbors: A (visited, skip), C (edge cost 2).
- C: candidate = 1 + 2 = 3. Compare to C's current tentative distance, 4. **3 < 4**, so update C to 3, via D.
Mark D visited.

| Node | Tentative distance | Visited? | Via |
|---|---|---|---|
| A | 0 | yes | — |
| B | 2 | no | A |
| C | **3** | no | **D** |
| D | 1 | **yes** | A |
| E | ∞ | no | — |

Notice what just happened: C's shortest path was *not* the direct A–C edge (cost 4) that was discovered first — it's A→D→C (cost 1 + 2 = 3), discovered on the very next step. This is exactly why the algorithm can't just take the first path it finds and stop; it must keep exploring until a node is actually selected as the current (minimum-distance) node before that node's distance is considered final.

**Step 3 — current node: B (distance 2, smallest unvisited).** Examine B's neighbors: A (visited, skip), C (edge cost 1), E (edge cost 3).
- C: candidate = 2 + 1 = 3. Compare to C's current tentative distance, 3. **3 is not less than 3** — no update (a tie doesn't overwrite; either path is equally valid, and this walkthrough keeps the first-found tie, via D).
- E: candidate = 2 + 3 = 5. Better than ∞ → update E to 5, via B.
Mark B visited.

| Node | Tentative distance | Visited? | Via |
|---|---|---|---|
| A | 0 | yes | — |
| B | 2 | **yes** | A |
| C | 3 | no | D |
| D | 1 | yes | A |
| E | **5** | no | **B** |

**Step 4 — current node: C (distance 3, smallest unvisited).** Examine C's neighbors: A (visited), D (visited), B (visited), E (edge cost 1).
- E: candidate = 3 + 1 = 4. Compare to E's current tentative distance, 5. **4 < 5**, so update E to 4, via C.
Mark C visited.

| Node | Tentative distance | Visited? | Via |
|---|---|---|---|
| A | 0 | yes | — |
| B | 2 | yes | A |
| C | 3 | **yes** | D |
| D | 1 | yes | A |
| E | **4** | no | **C** |

**Step 5 — current node: E (distance 4, only unvisited node left).** All of E's neighbors (B, C) are already visited. Nothing to update. Mark E visited. **Every node is now visited — the algorithm terminates.**

**Final shortest-path tree, computed entirely from A's own local view of the LSDB:**

| Destination | Shortest cost from A | Path |
|---|---|---|
| B | 2 | A → B |
| C | 3 | A → D → C |
| D | 1 | A → D |
| E | 4 | A → D → C → E |

This table is precisely what Router A installs into its routing table's OSPF-derived entries: for each destination, the total cost and the **first hop** of the shortest path (for D, C, and E, that first hop is D — Router A doesn't need to remember the whole path, only the next hop, exactly as Chapter 45 defined it). Every other router in this network runs this exact same algorithm on its own copy of the identical LSDB, just starting from itself instead of A — and because every router computed from the same underlying facts using the same deterministic algorithm, their independently-computed routes are guaranteed to be mutually consistent, with no possibility of the kind of divergent, self-reinforcing error Chapter 47 documented for RIP.

---

## 9. Why Link-State Converges Faster and Doesn't Count to Infinity

With Dijkstra's algorithm now concrete, Chapter 46's claim can finally be made precise, on two separate fronts:

**Why it doesn't count to infinity:** count-to-infinity (Chapter 47, Section 6) happens because a distance-vector router's belief about a destination is built entirely from a neighbor's *conclusion*, with no way to independently verify it or trace it back to a first-hand fact. A link-state router never does this — every input to its Dijkstra computation is a directly-flooded, individually-attributable fact about one specific link, signed by the router that owns it. When A's connection to Net-behind-D fails, A floods a new Router LSA saying "I no longer have a link to D" (or D floods one saying "I no longer have a link to A" — whichever side detects it). Every router, upon receiving that single updated fact, reruns Dijkstra on the corrected graph and gets the mathematically correct answer in one pass — there is no possibility of two routers re-inflating a stale rumor between themselves, because neither one is trusting the other's *conclusion* about C's or D's distance; they're only trusting each other's *direct, falsifiable claims about their own links*, and a claim about your own link cannot be a stale echo of someone else's earlier mistake.

**Why it converges faster:** RIP's convergence time scales with (network diameter) × (update interval), because information has to propagate and be recomputed hop by hop, each hop waiting on the next periodic cycle (Chapter 47, Section 5). OSPF's flooding is not periodic and not hop-limited by a timer — a single new LSA propagates outward essentially at the speed the network can physically forward small packets, reaching every router in the area within a small, roughly constant number of milliseconds regardless of how many hops separate them, because flooding happens in parallel across all links at once rather than being paced by a 30-second clock. Once every router has the new LSA, recomputing Dijkstra (a well-understood, efficient algorithm — O((V+E) log V) with a proper priority queue, for V routers and E links) takes milliseconds even for networks with thousands of routers. The net result, well documented in production deployments: OSPF convergence after a link failure is typically sub-second, versus RIP's tens of seconds to several minutes for a comparable topology.

---

## 10. OSPF Areas — Scaling the Link-State Idea

Section 2 flagged the honest cost of link-state routing: flooding LSAs to *everyone* and having *everyone* run Dijkstra on the *entire* topology doesn't scale forever — a network with tens of thousands of routers would mean an enormous LSDB, expensive floods on every minor change anywhere in the network, and a large, frequently-rerun Dijkstra computation on every router.

OSPF's answer is **hierarchical areas**: the network is divided into areas, each identified by a number, with **Area 0** mandated as the **backbone area** that every other area must connect to (directly or via a virtual link). Routers sitting on the boundary between two areas — **Area Border Routers (ABRs)** — are the only routers that need a full LSDB for more than one area; they summarize routes from one area into the next using Type 3 Summary LSAs (Section 4) rather than flooding every individual link-state fact across area boundaries.

The payoff, concretely: a router inside Area 1 only needs the *full* topology detail (every link, every cost) for Area 1 itself — for everything outside Area 1, it only needs to know "these summarized destination ranges are reachable via this ABR, at this total cost," which is a vastly smaller amount of information than the full topology of the entire rest of the network. A link flapping deep inside Area 2 triggers a flood and a Dijkstra rerun only among Area 2's own routers; Area 1's routers are entirely insulated from that churn, seeing at most a small summary-cost update if the flap actually changes the best summarized route. This is precisely the "regional hub" idea from Chapter 44's postal analogy, reapplied at the level of a link-state protocol's own internal database.

```
                    Area 0 (backbone)
                 ___________________
                /    |         |    \
          [ABR1]  [ABR2]   [ABR3]  [ABR4]
             |        |        |       |
         Area 1   Area 2   Area 3  Area 4
      (own full   (own full  (own full (own full
       LSDB, own   LSDB...)  LSDB...)  LSDB...)
       Dijkstra)
```

---

## 11. IS-IS — The ISP-Favored Cousin

IS-IS (Intermediate System to Intermediate System, ISO/IEC 10589) shares OSPF's entire link-state DNA — Hello-based neighbor discovery, flooded link-state information, an independently-computed Dijkstra shortest-path tree, and a two-level hierarchy (Level 1 areas and a Level 2 backbone, directly analogous to OSPF's non-backbone areas and Area 0). The differences that matter in practice:

- **Origin and encapsulation.** IS-IS was originally designed for OSI's CLNS (Connectionless Network Service) protocol stack, not IP, and runs directly over the data-link layer (encapsulated straight in Ethernet frames, not inside IP packets the way OSPF is). This is precisely *why* ISPs have historically favored it: because IS-IS never depended on IP in the first place, extending it to route both IPv4 and IPv6 simultaneously, in a single unified protocol instance and a single topology database, was a comparatively natural extension. OSPF, by contrast, needed a genuinely separate protocol version (OSPFv3) to handle IPv6, since OSPFv2 is intrinsically tied to IPv4 addressing throughout its packet formats.
- **Addressing.** IS-IS routers are identified by NSAP (Network Service Access Point) addresses rather than IP addresses for the protocol's own internal identity — an artifact of its OSI heritage that mostly shows up as an operational quirk (router IDs looking unfamiliar to IP-only network engineers) rather than a meaningful functional difference.
- **Track record at scale.** IS-IS has an especially strong reputation for stability in very large, single-organization backbone networks — many of the world's largest ISPs and content networks run IS-IS as their internal IGP specifically because of its maturity, its clean separation from IP-layer addressing changes/renumbering, and long-accumulated operational trust, rather than any large theoretical advantage over a well-tuned OSPF deployment.

For nearly every conceptual purpose in this course — the link-state idea, LSAs (called LSPs, Link State Packets, in IS-IS), flooding, areas, and Dijkstra — IS-IS and OSPF should be understood as the same underlying idea with different packaging and different historical adoption patterns, not two fundamentally different approaches to the routing problem.

---

## 12. OSPF vs. IS-IS vs. RIP — A Direct Comparison

| | RIP | OSPF | IS-IS |
|---|---|---|---|
| Family | Distance-vector | Link-state | Link-state |
| Metric | Hop count (max 15) | Cost (derived from bandwidth) | Cost (configurable, similar idea) |
| What's exchanged | Summarized distances | Raw link facts (LSAs) | Raw link facts (LSPs) |
| Path computation | Accumulated via Bellman-Ford-style hop counting | Dijkstra's algorithm, run locally | Dijkstra's algorithm, run locally |
| Convergence speed | Slow (tens of seconds to minutes) | Fast (sub-second, typically) | Fast (sub-second, typically) |
| Count-to-infinity risk | Yes (mitigated, not eliminated) | No — not structurally possible | No — not structurally possible |
| Scaling mechanism | None (15-hop hard limit) | Areas (Area 0 backbone + numbered areas) | Levels (Level 1 areas + Level 2 backbone) |
| Typical deployment | Legacy, small labs | Enterprise, data center, ISP-internal | Large ISP/carrier backbones |
| Runs over | UDP/IP | Directly over IP (protocol 89) | Directly over Layer 2 (not IP-dependent) |
| IPv6 support | Separate protocol (RIPng) | Separate protocol version (OSPFv3) | Same protocol, extended address families |

---

## 13. Code: A Minimal Dijkstra Implementation

A direct, general-purpose implementation of Section 7–8's algorithm, applied to Section 8's exact graph so you can check its output against the worked table by hand:

```go
package main

import (
	"container/heap"
	"fmt"
	"math"
)

type Edge struct {
	To   string
	Cost int
}

type Graph map[string][]Edge

type Item struct {
	node string
	dist int
}
type PriorityQueue []Item

func (pq PriorityQueue) Len() int            { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool  { return pq[i].dist < pq[j].dist }
func (pq PriorityQueue) Swap(i, j int)       { pq[i], pq[j] = pq[j], pq[i] }
func (pq *PriorityQueue) Push(x interface{}) { *pq = append(*pq, x.(Item)) }
func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	item := old[n-1]
	*pq = old[:n-1]
	return item
}

func dijkstra(graph Graph, start string) (map[string]int, map[string]string) {
	dist := map[string]int{start: 0}
	prevHop := map[string]string{}
	visited := map[string]bool{}

	pq := &PriorityQueue{{node: start, dist: 0}}
	heap.Init(pq)

	for pq.Len() > 0 {
		current := heap.Pop(pq).(Item)
		if visited[current.node] {
			continue
		}
		visited[current.node] = true

		for _, edge := range graph[current.node] {
			if visited[edge.To] {
				continue
			}
			candidate := current.dist + edge.Cost
			existing, known := dist[edge.To]
			if !known || candidate < existing {
				dist[edge.To] = candidate
				prevHop[edge.To] = current.node
				heap.Push(pq, Item{node: edge.To, dist: candidate})
			}
		}
	}
	return dist, prevHop
}

func main() {
	graph := Graph{
		"A": {{"B", 2}, {"D", 1}, {"C", 4}},
		"B": {{"A", 2}, {"C", 1}, {"E", 3}},
		"C": {{"A", 4}, {"B", 1}, {"D", 2}, {"E", 1}},
		"D": {{"A", 1}, {"C", 2}},
		"E": {{"B", 3}, {"C", 1}},
	}

	dist, _ := dijkstra(graph, "A")
	for _, node := range []string{"A", "B", "C", "D", "E"} {
		d := dist[node]
		if d == math.MaxInt {
			fmt.Printf("%s: unreachable\n", node)
			continue
		}
		fmt.Printf("%s: cost %d from A\n", node, d)
	}
}

// Output (matches Section 8's worked table exactly):
// A: cost 0 from A
// B: cost 2 from A
// C: cost 3 from A
// D: cost 1 from A
// E: cost 4 from A
```

This is not a simplification of what OSPF does internally — production OSPF implementations run essentially this exact algorithm (with production-grade priority queues and incremental recomputation optimizations) over their LSDB every time a topology-affecting LSA arrives.

---

## 14. A Real OSPF Neighbor Table

From a router's perspective, OSPF's state is directly inspectable — here's representative output from a Cisco-style `show ip ospf neighbor`:

```
Router# show ip ospf neighbor

Neighbor ID     Pri   State           Dead Time   Address         Interface
10.0.0.2          1   FULL/DR         00:00:38    10.1.1.2        GigabitEthernet0/0
10.0.0.3          1   FULL/BDR        00:00:39    10.1.1.3        GigabitEthernet0/0
10.0.0.4          1   FULL/ -         00:00:35    10.1.2.2        GigabitEthernet0/1
```

`FULL` means this router has completed a full adjacency (LSDB exchange, Section 3) with that neighbor. `DR`/`BDR` identify which neighbor won the Designated Router election (Section 5) on the shared segment; the third neighbor, on a different, presumably point-to-point interface, shows no DR/BDR role at all because none is needed there. `Dead Time` is the countdown until the 40-second dead interval (Section 3) would expire without a fresh Hello — a live, continuously-refreshed number confirming the neighbor relationship is healthy right now.

---

## 15. What's Simplified Here

- Section 8's worked Dijkstra example used a small, hand-tractable 5-node graph; real OSPF areas can contain hundreds of routers, but the algorithm and its guarantees scale unchanged — only the practical computation is larger, which is exactly the scaling pressure Section 10's areas exist to relieve.
- Real OSPF implementations use incremental SPF (iSPF) optimizations that avoid fully recomputing the entire shortest-path tree from scratch on every minor change — this chapter described the conceptually complete, from-scratch Dijkstra computation for clarity.
- This chapter treated "OSPF cost" and "IS-IS metric" as directly comparable ideas; IS-IS historically used narrower default metric ranges (a legacy 6-bit metric per link in the original specification, later widened) — a real operational detail, but one that doesn't change the underlying link-state/Dijkstra story told here.
- Equal-cost paths (as previewed in Chapter 45's ECMP note) are common in real OSPF/IS-IS deployments; Section 8's walkthrough noted one tie (B→C vs. via D) and kept a single path for clarity, but production routers frequently install and load-balance across multiple equal-cost shortest paths simultaneously.

---

## 16. Common Misconceptions

- **"Link-state routers know the topology of the entire Internet."** No — an OSPF or IS-IS router knows the complete topology only of its own area(s)/level(s), which is itself typically one organization's internal network. This is an IGP (Chapter 46, Section 9); routing between organizations is BGP's job (Chapter 49), which deliberately does *not* use link-state flooding at global scale — for reasons that chapter explains in full.
- **"OSPF routers ask each other for the shortest path."** They never ask for or exchange path *conclusions* at all — only raw link facts (LSAs). Every router computes its own shortest paths independently; this is the entire point, and precisely what Section 9 credits for eliminating count-to-infinity.
- **"Dijkstra's algorithm finds a path, then moves on to find the next path separately."** A single run of Dijkstra's algorithm from one starting node computes the shortest path to *every* other reachable node simultaneously, as Section 8's single pass demonstrated — it's a one-to-all algorithm, not a one-to-one algorithm rerun repeatedly.
- **"More OSPF areas always means better performance."** Areas trade a real cost (less detailed visibility, some suboptimality possible for routes that could have used a shorter path a summarized route hides) for scalability; a small network gains nothing from carving itself into multiple areas and only adds complexity.

---

## 17. Hands-On Experiment

Using FRRouting on a small set of Linux VMs or containers (the same setup style as Chapter 47's experiment), configure OSPF instead of RIP:

```
router ospf
 network 10.1.1.0/24 area 0
 network 10.1.2.0/24 area 0
```

Then inspect the LSDB directly and confirm every router in the topology holds an identical copy:

```
$ vtysh -c "show ip ospf database"
$ vtysh -c "show ip route ospf"
```

Bring a link down on one router and watch (via repeated `show ip ospf database` and a packet capture on `tcpdump -i any proto ospf`) how quickly a new LSA appears and propagates, compared to RIP's 30-second periodic cadence from Chapter 47's experiment — the contrast in convergence speed between the two protocols is directly, visibly measurable on the exact same lab topology.

---

## 18. Interview Questions & Model Answers

**Beginner: What is the core difference between distance-vector and link-state routing?**
Distance-vector routers exchange summarized distance opinions with immediate neighbors and never see the actual topology. Link-state routers flood raw, first-hand facts about their own direct links to every router in the area, so every router ends up with an identical, complete topology map and computes shortest paths independently using Dijkstra's algorithm.

**Intermediate: Walk through, at a high level, what Dijkstra's algorithm does.**
Starting from a source node with distance 0 (and every other node at infinity), it repeatedly selects the unvisited node with the smallest tentative distance, updates its neighbors' tentative distances if a shorter path through the current node is found, and marks the current node as finalized/visited. It terminates once every node has been visited, at which point every node's distance and recorded next hop represent the true shortest path from the source.

**Intermediate: Why does OSPF need Designated Routers on broadcast segments but not on point-to-point links?**
On a broadcast segment with N routers, full mesh adjacencies between every pair would require roughly N×(N−1)/2 relationships purely to describe one shared segment. Electing a DR (and backup BDR) reduces this to roughly 2×(N−1) adjacencies, with the DR originating a single Network LSA describing the whole segment. A point-to-point link only ever has two routers, so there's no redundant adjacency to eliminate — a DR would add complexity with no benefit.

**Advanced: Explain precisely why link-state routing cannot suffer count-to-infinity, in terms of what information actually gets exchanged.**
Count-to-infinity requires a router to unknowingly trust another router's *conclusion* about a distance that conclusion itself was originally derived from — creating a loop of mutual, un-verifiable trust. Link-state routers never exchange conclusions about distances at all; they only exchange direct, individually-attributable, falsifiable facts about their own links ("I have a link to X, costing Y"), which every recipient independently verifies as consistent with the rest of the flooded LSDB and computes shortest paths from itself, using Dijkstra. There is no conclusion to echo back, so the specific mechanism behind count-to-infinity has no analog in the link-state model.

---

## 19. Exercises

### Easy
1. In your own words, explain what a Link State Advertisement contains, and — just as importantly — what it deliberately does not contain.
2. What does the OSPF Designated Router do, and why is it only needed on broadcast segments?
3. Name one concrete advantage IS-IS has historically offered ISPs over OSPF, and explain why.

### Medium
4. Using Section 8's graph but changing the A–C edge cost from 4 to 2, redo the Dijkstra walkthrough from scratch and determine whether the shortest path to C changes.
5. Explain, referencing Section 6, why a network still running OSPF with the default 100 Mbps reference bandwidth might make poor routing decisions between a 1 Gbps link and a 10 Gbps link, and how an operator would fix this.
6. A router's OSPF Dead Time counter for a neighbor drops to 0. Explain precisely what the router does next, and how quickly the rest of the area learns about it.

### Hard
7. Implement Section 13's Dijkstra algorithm from scratch (without copying the code) for a 6-node graph of your own design, including at least one node reachable by two paths of different costs, and verify your output finds the genuinely cheaper path.
8. Design a two-area OSPF topology (Area 0 backbone plus two non-backbone areas, each with an ABR) for a company with a headquarters and two branch offices. Explain exactly what information crosses each ABR, and what stays hidden inside each area.
9. Explain, precisely and without hand-waving, why a router running Dijkstra's algorithm never needs to reconsider a node's distance once that node has been marked visited — i.e., prove informally why the algorithm's greedy choice at each step is safe, given that all edge costs are non-negative.

---

## Summary

| Term | Meaning |
|---|---|
| Link-state routing | Routers flood raw facts about their own direct links to everyone, then independently compute shortest paths — used by OSPF and IS-IS |
| Link State Advertisement (LSA) | A small, signed message describing a router's own direct links, flooded unmodified throughout an area |
| Link State Database (LSDB) | The complete, identical topology graph every router in an area assembles from received LSAs |
| Dijkstra's algorithm | The shortest-path algorithm every link-state router runs locally on its own LSDB to compute routes |
| OSPF cost | OSPF's metric, derived from link bandwidth (Reference Bandwidth / Interface Bandwidth), unlike RIP's flat hop count |
| Designated Router (DR) | The router elected to represent a broadcast segment, reducing the number of adjacencies needed on that segment |
| OSPF Area / Area 0 | OSPF's hierarchical scaling mechanism; Area 0 is the mandatory backbone all other areas connect through |
| IS-IS | A link-state IGP originally built for OSI's CLNS, favored by many large ISPs for its native, protocol-independent support for both IPv4 and IPv6 |
| Level 1 / Level 2 (IS-IS) | IS-IS's equivalent of OSPF's non-backbone areas and Area 0 backbone |

OSPF and IS-IS solve routing *inside* one organization's network, where every router can be trusted to report its links honestly and everyone shares the same goal: the shortest path. Chapter 49 leaves that trusted world behind entirely — BGP is the protocol that routes traffic *between* organizations that don't share a goal, don't trust each other by default, and often care more about business relationships than about which path is actually shortest.
