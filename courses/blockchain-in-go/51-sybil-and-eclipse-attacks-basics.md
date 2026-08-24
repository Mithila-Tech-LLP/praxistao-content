# Chapter 51: Sybil and Eclipse Attacks — Basics

Every mechanism this volume has built so far — seed nodes and peer exchange (Chapter 47), gossip (Chapter 48), synchronization (Chapter 49), and fork resolution (Chapter 50) — quietly assumes that at least some of the peers a node talks to are honest. This chapter examines what happens when that assumption is attacked directly: not by sending an invalid block or a forged signature (both of which our existing validation already rejects), but by controlling *who a victim node thinks its peers are* in the first place. If an attacker can control that, none of the validation logic built so far gets a chance to matter.

## Table of Contents

1. [Identity Is Free in a Peer-to-Peer Network](#1-identity-is-free-in-a-peer-to-peer-network)
2. [Sybil Attacks: Many Masks, One Puppeteer](#2-sybil-attacks-many-masks-one-puppeteer)
3. [Eclipse Attacks: Surrounding a Single Victim](#3-eclipse-attacks-surrounding-a-single-victim)
4. [Diagram: A Victim Node, Fully Eclipsed](#4-diagram-a-victim-node-fully-eclipsed)
5. [What an Eclipsed Node Actually Experiences](#5-what-an-eclipsed-node-actually-experiences)
6. [Why Validation Alone Doesn't Save You](#6-why-validation-alone-doesnt-save-you)
7. [Mitigation 1: Diverse Peer Selection](#7-mitigation-1-diverse-peer-selection)
8. [Mitigation 2: Connection Limits, Applied Deliberately](#8-mitigation-2-connection-limits-applied-deliberately)
9. [Mitigation 3: Preferring Outbound Connections](#9-mitigation-3-preferring-outbound-connections)
10. [Putting the Mitigations Together](#10-putting-the-mitigations-together)
11. [What GoChain Still Doesn't Defend Against](#11-what-gochain-still-doesnt-defend-against)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Identity Is Free in a Peer-to-Peer Network

Think back to Chapter 14's wallet addresses: generating a new one costs nothing more than running `crypto/ecdsa`'s key generation function — no registration, no fee, no proof that a human being is actually behind it. That same freedom, which is a genuine feature for financial privacy, becomes a liability the moment it applies to *network* identity too. A GoChain node's "identity," from a peer's point of view, is nothing more than a listen address it announced in a `VersionPayload` (Chapter 45). Nothing stops one person, running one machine, from starting a thousand `network.Node` processes, each with a different listen address, each looking — from the outside — like a thousand different, independent participants.

This is the foundational fact both attacks in this chapter exploit: **identity in a permissionless P2P network is free to create, and free identities are indistinguishable from real, independent participants using nothing but the protocol itself.** Every mitigation this chapter discusses is really an attempt to work around this one fact, since the protocol has no built-in way to fix it directly.

---

## 2. Sybil Attacks: Many Masks, One Puppeteer

A **Sybil attack** (named after a famous case study of a person diagnosed with many distinct personalities) is exactly what Section 1 described: one attacker creates a large number of fake node identities to gain influence disproportionate to the resources they actually control. In GoChain's specific design, "influence" mostly means **being the peer other nodes happen to talk to** — if an attacker's fake identities make up a large fraction of the addresses circulating through Chapter 47's peer exchange, a large fraction of the *real* network's connection slots end up occupied by nodes that all ultimately answer to one puppeteer.

```
  BEFORE A SYBIL ATTACK:                AFTER A SYBIL ATTACK:
  a healthy, diverse address book       one attacker floods knownAddrs with
                                         hundreds of addresses they control

  knownAddrs:                           knownAddrs:
    73.12.4.9    (real, Alice)            73.12.4.9    (real, Alice)
    91.88.6.201  (real, Bob)              91.88.6.201  (real, Bob)
    12.4.90.5    (real, Carol)            12.4.90.5    (real, Carol)
                                           203.0.113.1  (attacker mask #1)
                                           203.0.113.2  (attacker mask #2)
                                           203.0.113.3  (attacker mask #3)
                                                 ... (hundreds more)
                                           203.0.113.N  (attacker mask #N)
```

It is important to notice what a Sybil attack, by itself, does *not* do: creating a thousand fake identities does not give the attacker any extra mining power, any extra coins, or any ability to forge a signature. Chapter 50's `ChainWork` comparison and every block/transaction validation this course has built are completely untouched by how many fake addresses exist. A Sybil attack's real payoff is entirely about *positioning* — and the specific, dangerous form that positioning takes is the eclipse attack.

---

## 3. Eclipse Attacks: Surrounding a Single Victim

An **eclipse attack** takes the raw material a Sybil attack produces (many attacker-controlled identities) and aims it at one specific victim: get every single one of that victim's peer connection slots filled with an attacker-controlled node, so the victim's entire view of the network is whatever the attacker chooses to show it.

The name is deliberate and visual: in an astronomical eclipse, the moon does not destroy the sun — it simply gets between the sun and an observer, blocking every bit of real light from reaching them, while the observer might not immediately realize anything is wrong at all. An eclipsed GoChain node is the same: its own `Chain`, its own `Mempool`, its own cryptography — none of it is broken. What is broken is *what information ever reaches it in the first place*.

An eclipse attack against a node with `maxPeers` set to 8 (Chapter 47's default) needs the attacker to win all 8 of that node's connection slots — either by being the first 8 addresses the victim happens to dial (for instance, immediately after the victim starts up and has an empty peer list), or by patiently out-competing legitimate peers over time as the victim's connections naturally churn (a peer disconnecting, as Chapter 46's `removePeerByConn` handles routinely, opens exactly the kind of slot an attacker is waiting to fill).

A patient version of this attack rarely happens all at once — it plays out over hours or days, one opportunistic slot at a time:

```
   TIME          EVENT                                    VICTIM'S PEERS
   ----          -----                                    --------------
   Day 1, 09:00  Node starts, Bootstraps from 2 honest     [Honest-1, Honest-2]
                 seeds, dialSome fills the rest with
                 whatever addresses circulate first
   Day 1, 09:05  6 more slots fill from peer exchange,     [Honest-1..4,
                 a mix of honest and (unluckily) 4          Attacker-1..4]
                 attacker-controlled addresses
   Day 1, 14:00  Honest-3's connection drops (its own      [Honest-1,2,4,
                 process restarted for an update) --        Attacker-1..4,
                 an attacker address, pre-positioned        Attacker-5]
                 in the victim's knownAddrs, wins the
                 now-empty slot
   Day 2, 03:00  Honest-1 and Honest-4 both churn during   [Attacker-1..6,
                 routine restarts, hours apart -- both      Honest-2]
                 replaced by attacker addresses
   Day 2, 03:15  The victim's last honest connection,      [Attacker-1..8]
                 Honest-2, churns. Every slot is now
                 attacker-controlled. The eclipse is
                 complete, and nothing about how it
                 happened looked unusual at any single
                 moment.
```

Each individual event in that timeline — a peer disconnecting, `dialSome` filling the resulting slot from `knownAddrs` — is completely ordinary, expected behavior, exactly as designed in Chapter 46 and Chapter 47. That is precisely what makes a patient eclipse attack hard to notice while it's happening: no single step looks like an attack, only the accumulated result does.

---

## 4. Diagram: A Victim Node, Fully Eclipsed

```
                    HONEST NETWORK (what actually exists)
                    --------------------------------------

         [Alice]------[Bob]------[Carol]------[Dave]
            |            |           |            |
            +----[Eve]---+----[Frank]+----[Grace]--+

                   VICTIM'S ACTUAL VIEW (all 8 peer slots)
                   ----------------------------------------

                          +------------------+
                          |   VICTIM NODE    |
                          |   maxPeers = 8   |
                          +------------------+
                           /   |   |   |    \
                          /    |   |   |     \
                    [Atk-1][Atk-2][Atk-3][Atk-4][Atk-5]
                          \    |   |   |     /
                           \   |   |   |    /
                          [Atk-6][Atk-7][Atk-8]

    Every single one of the victim's 8 connections is to an attacker-
    controlled process. Alice, Bob, Carol, Dave, Eve, Frank, and Grace
    still exist and are still honestly running the real protocol -- the
    victim simply has no path to any of them. Every message the victim
    sends goes only to the attacker. Every message the victim receives
    was chosen, or fabricated, by the attacker.
```

The critical detail this diagram is trying to make visible: from the victim's own perspective, nothing looks obviously wrong. `handleVersion` (Chapter 46) still logs sensible-looking height comparisons for each of the 8 connections. `Broadcast` (Chapter 46) still successfully writes bytes to 8 live TCP connections. There is no error message, no crash, no red flag — just eight peers, all of whom happen to be lying.

---

## 5. What an Eclipsed Node Actually Experiences

Once every peer slot is attacker-controlled, the attacker can show the victim any of the following, entirely undetected by anything this course has built so far:

- **A withheld chain.** The attacker's 8 fake peers simply never forward gossiped blocks or transactions from the real network (Chapter 48's gossip only works if at least one of a node's peers is honestly relaying what it hears — an eclipsed node has zero such peers).
- **A fabricated chain.** Worse, the attacker's fake peers can feed the victim an entirely different, internally self-consistent chain — recall from Chapter 50 that `ValidateBlock` only checks that a block's own hash, proof-of-work, and linkage are internally correct; it says nothing about whether that block is the one the *rest* of the honest network agrees on. An attacker willing to spend real computational effort can produce a real, validly-mined *alternate* history and feed it exclusively to the eclipsed victim.
- **A double-spend against the victim specifically.** This is the classic, concrete payoff: convince the eclipsed victim (perhaps a merchant's node) that a payment is confirmed, using a fabricated view of the chain, while the real network never saw that transaction at all — the merchant ships goods against a "confirmed" transaction that the honest network will never actually settle.
- **Delayed or blocked awareness entirely.** Even without any fabrication, simply refusing to forward *anything* new leaves the victim believing the network is quiet, when in reality blocks and transactions are flowing normally everywhere except to them.

A concrete walkthrough makes the double-spend payoff less abstract. Imagine a coffee shop running a GoChain node to accept payments, and an attacker who has quietly eclipsed it overnight, exactly as in Section 3's timeline. The attacker broadcasts a real transaction paying the shop 5 gochips, using a UTXO they actually own — but privately, off-network, they've also prepared a *second* transaction spending that same UTXO back to an address they control (a plain double-spend, the exact scenario Chapter 34's mempool logic is designed to catch — but only among the peers a node actually talks to). The eclipsed shop's node sees the payment, sees a fabricated "confirmation" fed to it exclusively by the attacker's fake peers, and the register prints a receipt. The coffee is served. Meanwhile, on the real, honest network — the one the shop's node has no path to — the second, conflicting transaction is what actually gets mined, and the first one the shop believed in never existed anywhere else. This is not a flaw in `Mempool.Add`'s double-spend detection; it is a flaw in the shop's node having no honest peer left to detect the conflict *with*.

---

## 6. Why Validation Alone Doesn't Save You

It is worth being precise about exactly why Chapters 19, 25, 34, and 50's validation logic — all of which is real, all of which works exactly as designed — does not protect against this. Every one of those checks answers the question "is this specific piece of data internally well-formed and self-consistent?" None of them can answer the very different question "is this the same data the rest of the honest network is also seeing?" — because answering that second question requires actually talking to a representative sample of the honest network, which is precisely the thing an eclipse attack denies you.

```
              VALIDATION ANSWERS:              ECLIPSE ATTACKS EXPLOIT:
   "Is this block's hash correct?"      "Is this the ONLY block I'm being
   "Is the proof-of-work sufficient?"    shown for this height, when in
   "Does it link to a real prior       reality other, different blocks
    block (in MY chain)?"              exist that I'll simply never see?"

          ANSWERED BY:                        ANSWERED BY:
      ValidateBlock (Ch. 19, 25)        having enough diverse, honest
                                         connections that a full picture
                                         of the network actually reaches
                                         you (this chapter)
```

This is exactly why this chapter's mitigations are not about cryptography or validation logic at all — they are entirely about **peer selection**: making it hard, expensive, or simply unlikely for one attacker to occupy every connection slot a node has.

---

## 7. Mitigation 1: Diverse Peer Selection

Chapter 47's `dialSome` already shuffles candidate addresses before dialing, which is a small defense against always preferring whichever address happens to arrive first in a peer-exchange message (an attacker's addresses might otherwise dominate an unshuffled list simply by being numerous). This chapter adds a sharper form of diversity: refusing to let too many connections come from addresses that look like they belong to the same operator.

A crude but genuinely useful signal is **IP subnet diversity** — an attacker running a thousand fake node processes on rented cloud infrastructure very often runs them from a small number of IP address ranges, since spinning up truly independent infrastructure across many unrelated networks is expensive and slow, exactly the friction Section 1 said the protocol itself cannot enforce, but which a peer-selection policy can lean on:

```go
import (
	"net"
	"strings"
)

// subnetKey returns a coarse "which /24 network is this address in"
// string, used only to estimate whether two addresses are likely to be
// controlled by the same operator -- a real attacker with resources
// across many unrelated networks can still defeat this, but it meaningfully
// raises the cost of a naive, single-provider Sybil attack.
func subnetKey(addr string) string {
	host, _, err := net.SplitHostPort(addr)
	if err != nil {
		host = addr
	}
	parts := strings.Split(host, ".")
	if len(parts) != 4 {
		return host // not a recognizable IPv4 address -- fall back to the whole thing
	}
	return strings.Join(parts[:3], ".") // e.g. "203.0.113" from "203.0.113.7"
}

// tooManyFromSameSubnet reports whether accepting addr as a new peer
// would push this node's connections from addr's /24 subnet above a
// small cap, regardless of how many total peer slots remain free.
func (n *Node) tooManyFromSameSubnet(addr string) bool {
	target := subnetKey(addr)
	const maxPerSubnet = 2 // no single /24 network gets more than 2 of our slots

	n.mu.RLock()
	defer n.mu.RUnlock()

	count := 0
	for existing := range n.peers {
		if subnetKey(existing) == target {
			count++
		}
	}
	return count >= maxPerSubnet
}
```

Wiring this into Chapter 47's `dialSome` is a one-line addition to its existing checks:

```go
// dialSome, extended from Chapter 47 with a subnet-diversity check.
func (n *Node) dialSome(candidates []string) {
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, addr := range candidates {
		n.mu.RLock()
		full := len(n.peers) >= n.maxPeers
		_, alreadyConnected := n.peers[addr]
		n.mu.RUnlock()

		if full {
			return
		}
		if alreadyConnected {
			continue
		}
		if n.tooManyFromSameSubnet(addr) { // new in this chapter
			log.Printf("[%s] skipping %s -- already have enough peers from that subnet", n.Address, addr)
			continue
		}
		if err := n.Dial(addr); err != nil {
			log.Printf("[%s] could not dial learned address %s: %v", n.Address, addr, err)
			continue
		}
	}
}
```

---

## 8. Mitigation 2: Connection Limits, Applied Deliberately

Chapter 47's `maxPeers` was already introduced mostly as a resource-management concern ("don't accumulate an unbounded number of open sockets"). It doubles as a security mitigation once you notice what it implies for an attacker: **a small `maxPeers` means an attacker needs to win a small, fixed number of specific connection slots**, not an unbounded number — which sounds at first like it makes the attack *easier*, but actually matters most in combination with Section 9's next mitigation. A large `maxPeers` gives an attacker with many fake identities a bigger, easier target to fill; a small one, defended by outbound preference, gives them far less room to work with, and forces them to win a fight for a handful of slots that are actively being protected rather than passively left open.

The connection limit's real security value is in making the *cost of the fight* concrete and boundable: with `maxPeers = 8`, an attacker knows exactly how many of a victim's slots exist to compete for, which is the necessary precondition for the next two mitigations (preferring outbound connections, and refusing to let inbound connections dominate) to have any teeth at all.

It helps to think of this the way a building's security guard thinks about doors, not headcount. A building with a hundred doors, all unlocked and unwatched, is not safer than one with eight doors that are each individually monitored — it is far less safe, because there are simply more places for an intruder to try. `maxPeers` is GoChain's way of deciding, on purpose, exactly how many doors this node has, so that the *other* mitigations in this chapter have a fixed, small, defensible surface to actually protect, rather than an ever-growing one that scales with however many fake identities an attacker is willing to manufacture.

---

## 9. Mitigation 3: Preferring Outbound Connections

Recall from Chapter 47 that `Peer.Outbound` records whether *we* dialed a peer (`true`) or they dialed *us* (`false`). This distinction matters enormously for eclipse resistance: an attacker attempting to eclipse a specific victim by dialing it repeatedly can only ever fill that victim's **inbound** slots — they cannot force the victim to dial *them*, since `dialSome` only ever dials addresses the victim itself learned about and chose to act on. This means a node that deliberately reserves a majority of its peer slots for outbound connections — ones it chose to make, from its own address book, built up over time through Chapter 47's peer exchange — is much harder to fill entirely with an attacker's inbound connection attempts alone.

```go
// tooManyInbound reports whether accepting one more inbound connection
// would push this node's ratio of inbound-to-total peers above a safe
// threshold. Reserving most slots for OUTBOUND connections -- ones this
// node itself chose to make -- means an attacker who can only dial IN
// to us can never occupy more than a minority of our peer table, no
// matter how many fake identities they throw at us.
func (n *Node) tooManyInbound() bool {
	const maxInboundFraction = 0.25 // at most 1 in 4 peers may be inbound

	n.mu.RLock()
	defer n.mu.RUnlock()

	if len(n.peers) == 0 {
		return false
	}
	inbound := 0
	for _, p := range n.peers {
		if !p.Outbound {
			inbound++
		}
	}
	return float64(inbound+1)/float64(len(n.peers)+1) > maxInboundFraction
}
```

This check belongs in `Listen`'s accept path (Chapter 46), right before a newly accepted connection is registered as a peer:

```go
// Listen, extended from Chapter 46 with an inbound-connection cap.
func (n *Node) Listen() error {
	listener, err := net.Listen("tcp", n.Address)
	if err != nil {
		return err
	}
	log.Printf("[%s] listening for peers", n.Address)

	go func() {
		for {
			conn, err := listener.Accept()
			if err != nil {
				log.Printf("[%s] accept error: %v", n.Address, err)
				continue
			}
			if n.tooManyInbound() { // new in this chapter
				log.Printf("[%s] rejecting inbound connection from %s -- inbound quota full", n.Address, conn.RemoteAddr())
				conn.Close()
				continue
			}
			go n.handleConnection(conn, true)
		}
	}()

	return nil
}
```

An attacker can still occupy that reserved 25% of inbound slots relatively cheaply — that much stays true even with this defense in place. What they *cannot* do anymore is fill the other 75%, because those slots only ever get filled by addresses this node itself chose to dial, drawn from an address book built up through Chapter 47's diverse, shuffled peer exchange, which the attacker can influence but not fully dictate unless they also dominate the addresses circulating through peer exchange generally — a meaningfully harder, more expensive attack than simply opening a lot of inbound connections.

---

## 10. Putting the Mitigations Together

None of Sections 7, 8, or 9 is individually sufficient, and that is by design — each closes off one specific angle an attacker could otherwise use cheaply, and together they raise the cost of a full eclipse substantially without claiming to make it theoretically impossible:

```
   ATTACKER'S GOAL: fill all 8 of the victim's peer slots

   Without mitigations:          With mitigations from this chapter:
   -------------------------     ------------------------------------
   Just open 8 connections        Inbound capped at ~2 of 8 slots
   to the victim. Done.           (Section 9) --

                                  The remaining ~6 outbound slots are
                                  only filled by addresses the victim
                                  itself dialed, chosen with subnet
                                  diversity (Section 7) from a peer-
                                  exchange address book the attacker
                                  would need to dominate broadly, not
                                  just locally, to control --

                                  and maxPeers (Section 8) keeps the
                                  whole fight bounded to a small,
                                  knowable number of slots rather than
                                  an ever-growing target.
```

---

## 11. What GoChain Still Doesn't Defend Against

Being honest about the limits of this chapter's defenses matters as much as the defenses themselves. None of Sections 7-9 stop an attacker with truly diverse infrastructure (real machines, spread across many real, unrelated networks and providers) from patiently building a large, subnet-diverse pool of addresses and dominating peer exchange legitimately, by the numbers, over a long period of time — the mitigations here raise the *cost* of an eclipse attack meaningfully, they do not make one *impossible*. Real production networks add further layers this course does not implement: hardcoded "anchor" connections to a small set of known-good peers that survive restarts, randomized connection eviction (occasionally dropping and replacing an existing peer even when nothing has gone wrong, so an attacker who has been patiently waiting for churn can't rely on it happening predictably), and out-of-band reputation or allow-lists maintained by node operators. Chapter 76, later in this course, revisits attacks on the network at a broader scale, including the 51% attack, which shares eclipse's core lesson: consensus safety depends on actually seeing what the honest majority of the network is doing, not merely on validating whatever you're shown.

It's also worth naming the trade-off every one of this chapter's mitigations makes, since none of them are free. Subnet diversity (Section 7) can wrongly treat two unrelated node operators who happen to share a hosting provider's address range as suspicious, and it can be defeated outright by an attacker willing to rent infrastructure across several providers. Reserving inbound slots (Section 9) means a smaller pool of connections is available to serve *other* nodes trying to reach us, which is a real cost to the network's overall connectivity, paid in exchange for our own node's safety. There is no configuration of these three mitigations that makes eclipse attacks free to defend against and costless to the honest network at the same time — this is the same kind of engineering trade-off Chapter 26's difficulty adjustment and Chapter 35's fee market both made peace with earlier in this course: security and efficiency pull in opposite directions, and the right answer is a deliberate, documented choice, not a free lunch.

---

## Summary

- Identity is free to create in a permissionless P2P network — nothing stops one attacker from running many fake node processes, each looking like an independent participant.
- A **Sybil attack** floods a network (or, more precisely, floods the addresses circulating through peer exchange) with many attacker-controlled identities to gain disproportionate influence over who real nodes end up connecting to.
- An **eclipse attack** aims that flood at one specific victim, trying to fill every one of its peer connection slots with attacker-controlled nodes, so the victim's entire view of the network's blocks and transactions is whatever the attacker chooses to show it.
- Existing validation logic (`ValidateBlock`, `Mempool.Add`, `ChainWork`) cannot detect an eclipse attack, because it only checks whether data is internally self-consistent, never whether it matches what the rest of the honest network actually sees.
- **Diverse peer selection** (rejecting too many connections from the same IP subnet, extending Chapter 47's `dialSome`) raises the cost of a single-operator Sybil flood dominating one node's connections.
- **Connection limits** (`maxPeers`, from Chapter 47) bound exactly how many slots an attacker needs to fight for, which is a precondition for the other mitigations to matter.
- **Preferring outbound connections** — reserving most peer slots for addresses this node itself chose to dial, and capping how many inbound (attacker-dialable) connections are allowed — means an attacker who can only connect *to* a victim, rather than influence who the victim dials *out* to, can never occupy more than a minority of that victim's peer table.
- None of these mitigations claim to make an eclipse attack impossible against a sufficiently resourced, patient attacker — they raise its cost and difficulty substantially, which is the same honest, probabilistic security story proof-of-work itself tells (Chapter 24).

---

## Exercises

### Easy

1. **Explain in your own words**, using the eclipse diagram in Section 4, why an eclipsed victim node's own logs (from `handleVersion` and `Broadcast`) would show nothing unusual, even while every single peer is lying to it.

2. **Compute the numbers**: with `maxPeers = 8` and `maxInboundFraction = 0.25` (Section 9), how many of a victim's slots can an attacker fill using only inbound connections? How many additional slots would the attacker need to influence via peer exchange (Section 7) to eclipse the victim completely?

3. **Trace `subnetKey`** (Section 7) on the addresses `"203.0.113.5:3000"`, `"203.0.113.199:3000"`, and `"198.51.100.5:3000"`. Which two are considered the "same subnet" by this function, and is that grouping actually guaranteed to mean the same operator controls both? Explain the false-positive and false-negative cases this heuristic can produce.

### Medium

4. **Implement and test `tooManyFromSameSubnet`**: write a unit test that registers three fake peers all in `203.0.113.0/24`, then asserts a fourth address from that same subnet is rejected by `dialSome` while an address from a different subnet is accepted (assuming `maxPeers` isn't yet full).

5. **Implement randomized peer eviction**: add a background goroutine that, every few minutes, closes one randomly chosen *outbound* peer connection (never an inbound one, to avoid making the inbound-preference mitigation pointless) even if nothing is wrong with it, forcing `dialSome` to eventually pick a replacement. Explain in a comment why deliberately disconnecting a perfectly healthy peer can improve security rather than just hurting uptime.

6. **Simulate a partial eclipse attempt** with an in-process test: create one victim `Node` with `maxPeers = 4`, then have four attacker-controlled `Node`s all try to connect inbound to it in quick succession. With Section 9's `tooManyInbound` check in place, confirm no more than one of the four inbound connections is accepted, and write a short comment explaining what would have happened without that check.

### Hard

7. **Design an anchor-connections feature**: extend `Node` with a small, separate list of "anchor" peer addresses (say, 2-3) that are exempt from the normal churn-and-eviction logic and are always the very first addresses reconnected to after a restart, specifically so a node that was eclipsed once cannot be trivially re-eclipsed the moment it restarts and rebuilds its peer list from scratch. Implement it and explain what new risk this introduces if an anchor address itself later becomes attacker-controlled.

8. **Research and summarize (200-300 words) a real eclipse attack finding** against Bitcoin or Ethereum (several academic papers exist describing practical eclipse attacks against Bitcoin's peer-to-peer layer). Compare the specific mitigations that paper recommends against the three implemented in this chapter (diverse peer selection, connection limits, outbound preference) — which of GoChain's defenses does the real-world research validate, and which additional defenses does it suggest that this chapter did not implement?

9. **Build a small attack-simulation harness**: spin up one victim `Node` and a configurable number of attacker `Node`s (all sharing a small number of simulated "subnets" via a test-only address-generation helper), have the attackers attempt to fill the victim's peer table using both inbound dials and by injecting `MsgAddr` messages advertising more attacker addresses, and measure what fraction of the victim's final peer table ends up attacker-controlled with all of Sections 7-9's mitigations enabled versus disabled. Report your numbers and discuss whether the mitigations behaved the way this chapter predicted.
