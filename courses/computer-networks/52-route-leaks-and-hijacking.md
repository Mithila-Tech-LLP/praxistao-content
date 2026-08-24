# Chapter 52: Route Leaks and Route Hijacking — When Routing Breaks

> *"BGP was built for a world where every network told the truth about what it owned. The Internet outgrew that world decades ago, and the protocol never fully caught up."*

---

## Table of Contents

1. [The Trust Assumption at the Center of Everything](#1-the-trust-assumption-at-the-center-of-everything)
2. [Two Different Failures: Leaks vs. Hijacks](#2-two-different-failures-leaks-vs-hijacks)
3. [Why Longest-Prefix Match Makes This So Dangerous](#3-why-longest-prefix-match-makes-this-so-dangerous)
4. [Route Leaks, Mechanically](#4-route-leaks-mechanically)
5. [Route Hijacking, Mechanically](#5-route-hijacking-mechanically)
6. [Real Incident #1: Pakistan Telecom vs. YouTube (2008)](#6-real-incident-1-pakistan-telecom-vs-youtube-2008)
7. [Real Incident #2: Google's Route Leak Breaks Japan (2017)](#7-real-incident-2-googles-route-leak-breaks-japan-2017)
8. [Real Incident #3: The MyEtherWallet BGP Hijack (2018)](#8-real-incident-3-the-myetherwallet-bgp-hijack-2018)
9. [Why These Incidents Took So Long to Detect and Fix](#9-why-these-incidents-took-so-long-to-detect-and-fix)
10. [The Real, Partial Defense: RPKI](#10-the-real-partial-defense-rpki)
11. [How Route Origin Validation Actually Works](#11-how-route-origin-validation-actually-works)
12. [What RPKI Does NOT Fix](#12-what-rpki-does-not-fix)
13. [Beyond RPKI: ASPA, BGPsec, and MANRS](#13-beyond-rpki-aspa-bgpsec-and-manrs)
14. [Packet/Object View: A ROA](#14-packetobject-view-a-roa)
15. [A Real Example: Checking a Prefix's RPKI Status](#15-a-real-example-checking-a-prefixs-rpki-status)
16. [Hands-On Experiment](#16-hands-on-experiment)
17. [Code: A Minimal Route Origin Validator in Go](#17-code-a-minimal-route-origin-validator-in-go)
18. [Common Misconceptions](#18-common-misconceptions)
19. [Production Notes](#19-production-notes)
20. [What This Chapter Simplified](#20-what-this-chapter-simplified)
21. [Interview Questions & Model Answers](#21-interview-questions--model-answers)
22. [Exercises](#22-exercises)
23. [Summary, and the Bridge to Part 8](#23-summary-and-the-bridge-to-part-8)

---

## 1. The Trust Assumption at the Center of Everything

Every chapter in this volume has quietly relied on one assumption: **that when an Autonomous System announces "I can reach this prefix," it's telling the truth.** BGP's path-vector design (Chapter 49) elegantly solves loop detection and lets networks express rich policy — but notice what it was never designed to solve: **verifying that an AS actually has the right to originate a given prefix in the first place.**

When AS 65001 sends a BGP UPDATE announcing `203.0.113.0/24`, nothing in the base BGP protocol checks whether AS 65001 is the legitimate holder of that address block. Every neighboring router simply believes it, subject only to whatever manual filtering (Chapter 49, Section 16; Chapter 51, Section 15) that neighbor happens to have configured. This is precisely analogous to a postal system where anyone can walk up and declare "I am now the official delivery address for 123 Main Street" and, absent anyone checking, the mail starts arriving at their door instead.

For most of the Internet's history, this worked well enough because the community was small, mostly cooperative, and mistakes were rare and usually caught fast by attentive operators. As the Internet grew to tens of thousands of independently-run networks (Chapter 50) with wildly varying levels of operational maturity, that trust model became a real, repeatedly-exploited liability — sometimes by accident, sometimes on purpose.

---

## 2. Two Different Failures: Leaks vs. Hijacks

It's worth being precise about the difference between these two terms, because they get conflated constantly in casual conversation, but they describe genuinely different failures:

- **A route leak** happens when a network **propagates a route beyond its intended scope**, typically by violating the export rules from Chapter 51, Section 15 — for example, re-advertising a route learned from one peer or transit provider to another peer or transit provider, which effectively (and usually accidentally) offers free transit service the leaking network never intended to provide and isn't provisioned to carry. The origin AS in a route leak is often perfectly legitimate — the *path* the route takes is what's wrong.
- **A route hijack** happens when a network **originates a route for a prefix it does not own or have any right to announce** — impersonating the legitimate holder of that address space, whether by mistake (a fat-fingered configuration) or with malicious intent (redirecting traffic for surveillance, censorship, or theft).

The two failures often *combine* in practice: a hijacked route, once originated by the wrong AS, still has to *propagate* through the Internet to actually cause damage — and it propagates through exactly the same unfiltered acceptance and re-advertisement behavior that makes leaks possible. Section 6's incident, in fact, is usually described as beginning as an intentional (if narrowly-scoped) hijack that turned into a global incident specifically *because* it also leaked far beyond its intended audience.

---

## 3. Why Longest-Prefix Match Makes This So Dangerous

Chapter 50, Section 10 already planted the exact mechanism this chapter is built on, so it's worth restating with full weight: **forwarding uses longest-prefix match (Chapter 45), and that lookup happens independently of, and prior to, BGP's own best-path comparison between competing routes to the identical prefix.**

This means if the legitimate holder of `198.51.100.0/24` (Chapter 50's whole point of that block being aggregated) is announcing it as part of a larger `/16`, and *anyone else, anywhere on the Internet*, announces a more-specific `/24` or smaller sub-block matching part of it, every router that hears both routes will send traffic for addresses inside that more-specific block to the impostor — **completely regardless of AS-PATH length, LOCAL-PREF, or any other BGP policy attribute**, because those attributes only matter when comparing multiple routes to the *exact same* prefix; a more-specific prefix simply wins the forwarding decision outright, every time, by design. This is not a bug in longest-prefix match — it's the correct, intended behavior that lets legitimate traffic engineering work (Chapter 50, Section 10) — but it is precisely the mechanism every incident in this chapter exploits.

---

## 4. Route Leaks, Mechanically

The IETF's RFC 7908 formally catalogs several distinct route leak patterns, but the core mechanical failure in nearly all of them is the same violation modeled in Chapter 51, Section 15's `exportAllowed` function: **a route learned from a peer or a transit provider gets re-advertised to another peer or transit provider**, instead of being restricted to the leaking network's own customers.

```
                    LEGITIMATE (no leak)

  Transit Provider T  ──sells transit──►  AS X  ──sells transit──►  Customer C
       (T's routes reach the whole Internet through X, as paid for)


                    ROUTE LEAK

  Peer P  ◄──peering──►  AS X  ◄──peering──►  Peer Q

  AS X incorrectly re-advertises routes learned from Peer P to Peer Q
  (and vice versa) -- AS X is now accidentally offering FREE TRANSIT
  between P and Q, something neither P nor Q agreed to, and something
  AS X's own network was never provisioned (in bandwidth or peering
  contract) to actually carry.
```

The damage a leak causes comes from two directions at once: traffic that should have taken an efficient, provisioned path (through Peer P and Peer Q's *actual* transit arrangements) instead gets pulled toward AS X, because — remember Chapter 49, Section 10 — a shorter or more-preferred-looking AS-PATH through the leaking network can genuinely look attractive to routers elsewhere on the Internet, even though the leaking network has no real capacity or business relationship to carry that traffic. The result is typically congestion, packet loss, or a de facto blackhole at the leaking network, for traffic that had nothing to do with it.

---

## 5. Route Hijacking, Mechanically

A hijack is simpler to describe and, mechanically, is just an ordinary BGP UPDATE with a false claim in it: an AS originates a route for a prefix using its own ASN as the origin, for an address block it has no legitimate allocation for. There's no special "hijack packet" — it's a completely ordinary-looking UPDATE message (Chapter 49, Section 12), indistinguishable at the protocol level from a legitimate announcement, unless something *outside* BGP itself (Section 10's RPKI, a routing registry, or a human noticing) checks whether the origin is actually authorized.

Hijacks fall into a rough spectrum of intent and precision:

- **Accidental / fat-finger hijacks**: an engineer misconfigures a router — copies a customer's prefix into the wrong `network` statement, or mistypes a subnet mask — and originates a route they never meant to announce at all. These are, by far, the most common category historically.
- **Narrowly-targeted intentional hijacks**: a network deliberately originates a specific block for a specific, limited purpose (Section 6's incident began this way — Pakistan Telecom fully intended to announce a false route, just not for the whole world to see).
- **Malicious, traffic-interception hijacks**: an attacker deliberately hijacks a prefix specifically to intercept, inspect, or redirect the traffic of others — for surveillance, censorship, or, in Section 8's case, financial theft.

---

## 6. Real Incident #1: Pakistan Telecom vs. YouTube (2008)

This is the single most-cited BGP incident in networking education, because it demonstrates the leak/hijack combination (Section 2) with unusual clarity, and its cause and effect are both well-documented.

**Background:** On February 24, 2008, the Pakistani government ordered domestic ISPs to block access to YouTube, in response to a video hosted there that the government considered blasphemous.

**What Pakistan Telecom actually did:** To implement the block *within Pakistan*, Pakistan Telecom (AS 17557) configured one of its routers to originate a bogus, highly-specific route for a `/24` block that fell inside YouTube's real address space (YouTube, at the time under AS 36561, was legitimately announcing that space as a `/22` — a larger, less-specific block). Pakistan Telecom's router pointed this `/24` at a null interface — a deliberate blackhole, a standard technique for a network to drop traffic to a destination *for its own users* without actually contacting that destination.

**What went wrong:** This configuration was only ever meant to affect routing *inside Pakistan Telecom's own network*, for its own customers. But Pakistan Telecom then advertised this same bogus `/24` route to its upstream international transit provider, **PCCW Global** (AS 3491) in Hong Kong — and PCCW's router **did not filter or reject the announcement**, instead accepting it and re-advertising it onward to PCCW's own peers and customers across the global Internet, exactly the unfiltered-export failure described in Section 4.

**The result:** Because Pakistan Telecom's announced block was a `/24` — more specific than YouTube's legitimate `/22` — longest-prefix match (Section 3) meant that **routers all over the world**, having now heard both YouTube's real `/22` and this bogus, more-specific `/24`, preferred the bogus route for any address falling inside it. Traffic from users worldwide who were trying to reach YouTube's servers within that `/24` instead got routed to Pakistan Telecom's network, where it was silently dropped by the null route Pakistan Telecom had configured purely for its own domestic censorship purposes.

**YouTube went effectively unreachable for a significant fraction of the entire global Internet for roughly two hours.**

**How it was resolved:** YouTube's own network operations team noticed the anomaly and responded by announcing even more specific routes of their own (`/24`s matching and beating Pakistan Telecom's announcement, and in some accounts even more specific blocks still, to reliably win the longest-prefix-match contest back). Simultaneously, once the incident became visible and was reported through operator mailing lists and monitoring services, PCCW withdrew the erroneous announcement and put filtering in place. Global reachability was restored within a few hours of the initial event.

**Why this incident remains the canonical teaching example:** it cleanly demonstrates all three moving parts at once — an AS announcing a block it doesn't own (a hijack, even if narrowly intended), that announcement propagating far beyond its intended scope because an upstream failed to filter it (a leak), and the entire failure mode being enabled purely by longest-prefix match having no concept of authorization (Section 3).

---

## 7. Real Incident #2: Google's Route Leak Breaks Japan (2017)

**Background:** On August 25, 2017, Google's network (AS 15169) — one of the largest, most sophisticated network operators on Earth — mistakenly propagated a large batch of BGP routes it had received from one of its own peering or customer BGP sessions, re-advertising them onward as though Google itself were a legitimate transit path for them. This is the leak pattern from Section 4, occurring inside a network most people would assume was far too well-defended for it to happen.

**What went wrong:** Among the routes affected were prefixes belonging to major Japanese network operators, including **NTT Communications (part of OCN, one of Japan's largest ISPs) and KDDI**. Google's network briefly appeared, to large parts of the global routing system, to be a valid path to reach these Japanese networks' address space — despite Google having no actual transit relationship or provisioned capacity to carry that traffic on their behalf.

**The result:** Substantial volumes of Japan-bound (and Japan-internal) traffic that should have traveled directly between Japanese networks, or through their actual, provisioned transit and peering arrangements, was instead pulled toward Google's network by the leaked announcements — and once it arrived there, Google's network was not built to actually forward it onward as genuine transit traffic at that path. The practical effect was severe congestion and packet loss for a significant portion of Japan's domestic Internet traffic for roughly **an hour**, with real-world, publicly reported impact including disruptions to some banking services and railway-related systems that depended on the affected connectivity, alongside general slowdowns and outages reported broadly across Japanese networks and services during the incident window.

**How it was resolved:** Once detected — both by automated BGP-monitoring services run by third parties that track global routing table anomalies in near-real-time, and by the affected networks' own operations teams — Google corrected the erroneous route advertisements and withdrew them, restoring normal routing within roughly an hour of the leak's onset.

**Why this incident matters for this chapter:** it's a sobering demonstration that route leaks are not a problem exclusive to small, under-resourced networks — even a network with Google's engineering resources produced, in this case, an internal error (widely attributed to an automated configuration or route-redistribution process behaving unexpectedly during an internal network change) that manifested as exactly the same fundamental failure mode Pakistan Telecom's incident showed nine years earlier: a route reaching an audience, and a scope of impact, its origin network never intended and had no way to actually serve.

---

## 8. Real Incident #3: The MyEtherWallet BGP Hijack (2018)

**Background:** On April 24, 2018, attackers carried out a BGP hijack with a specific, financially-motivated target: the DNS infrastructure behind **MyEtherWallet**, a popular web-based cryptocurrency wallet interface.

**What the attackers did:** Rather than hijacking MyEtherWallet's own servers directly, the attackers hijacked a block of IP addresses (`205.251.192.0/24`, part of the range used by **Amazon's Route 53** DNS service) by announcing a bogus, more-specific route for it via a rogue BGP session. This is the pure hijacking pattern from Section 5, applied with unusually clear malicious precision: rather than targeting web servers or application infrastructure, the attackers targeted the **DNS resolution layer itself**, redirecting DNS queries for `myetherwallet.com` toward malicious DNS servers under their control.

**The result:** Those malicious DNS servers returned forged DNS records pointing `myetherwallet.com` at a server the attackers controlled, hosting a **near-identical phishing clone** of the real MyEtherWallet site. Users who visited the real domain name, with no visible warning sign in their browser's address bar (this attack operated below the level any normal user could detect, since the domain name itself was correct — only the underlying routing and DNS resolution had been subverted), were served the fake site instead, which prompted them to enter their wallet credentials and private keys. Reports at the time estimated roughly **$150,000 in cryptocurrency** was stolen from users who fell for the phishing page before the hijack was identified and shut down.

**How it was resolved:** The anomalous BGP announcement was detected by third-party BGP monitoring services within about two hours, publicized rapidly across the network operator community, and the affected upstream providers stopped propagating the rogue route, restoring correct routing to Amazon's legitimate Route 53 infrastructure.

**Why this incident matters for this chapter:** it demonstrates that BGP hijacking isn't only a source of accidental outages — it is a viable, real-world **attack technique**, capable of undermining security assumptions (like "the domain name in my browser's address bar is enough to trust this site") that have nothing to do with routing on their face, because DNS resolution (Chapter 66-69) itself depends entirely on IP-layer routing delivering DNS queries to the *correct* server in the first place. A hijack of the IP address a DNS server lives at is, functionally, a hijack of every domain name that server is authoritative for.

---

## 9. Why These Incidents Took So Long to Detect and Fix

A recurring, uncomfortable pattern across all three incidents: detection took anywhere from roughly one to a few hours, not seconds. This is a direct consequence of how BGP itself is designed and operated:

- **No protocol-level authorization check exists by default** (Section 1) — a hijacked or leaked route looks, to any router receiving it, exactly like a legitimate one, so nothing in the protocol itself raises an alarm.
- **BGP convergence and propagation, while fast by human standards, still takes real time** (Chapter 49, Section 16) — by the time a bad route has propagated globally, reversing it requires the correction to propagate just as widely, which isn't instantaneous either.
- **Detection has historically relied on humans and third-party monitoring services** noticing anomalies — a route appearing for a prefix from an AS that has never announced it before, or a sudden, suspicious change in a prefix's usual AS-PATH — rather than on any built-in protocol mechanism. Services like the historical **BGPmon** (now part of other commercial offerings) and academic/nonprofit efforts pioneered exactly this kind of "does this new route make sense given routing history" anomaly detection, which is genuinely useful but is fundamentally a *monitoring* layer bolted on top of BGP, not a fix to BGP's trust model itself.

This detection gap — real, human-hours-long windows where a hijack or leak actively misdirects global traffic — is exactly the gap Section 10's defense is designed to close.

---

## 10. The Real, Partial Defense: RPKI

**RPKI (Resource Public Key Infrastructure)** is the Internet community's actual, deployed answer to Section 1's core problem: *how can a router cryptographically verify that an AS is actually authorized to originate a given prefix, without trusting the announcement itself?*

RPKI extends the existing IP address allocation hierarchy — IANA at the top, delegating to the five Regional Internet Registries (Chapter 50, Section 4), which delegate further to ISPs and end organizations — into a proper **public key infrastructure**, conceptually parallel to the certificate authority system that secures TLS (Chapter 81). Each RIR operates as a certificate authority for the address space it has allocated, and legitimate address holders can use this infrastructure to create a specific, narrowly-scoped, cryptographically-signed statement:

> **"AS 65001 is authorized to originate the prefix `203.0.113.0/24` (or any more-specific prefix up to a stated maximum length)."**

This signed statement is called a **ROA — Route Origin Authorization**. It's published into a globally-replicated repository system that any network operator's routers (or a nearby validation server acting on their behalf) can fetch and check incoming BGP announcements against.

---

## 11. How Route Origin Validation Actually Works

**Route Origin Validation (ROV)** is the operational process of actually using published ROAs to check live BGP announcements:

```
1. A dedicated RPKI validator (software like Routinator, RPKI-client,
   OctoRPKI, or FORT Validator -- typically run on a small server near
   the router, NOT on the router's own limited CPU) periodically fetches
   the full, current set of published ROAs from the global RPKI
   repository system.

2. The router itself uses a lightweight protocol (RPKI-to-Router, RFC
   8210) to pull a simplified table from the local validator: a list of
   (prefix, max-length, authorized-origin-ASN) tuples.

3. For every incoming BGP UPDATE, the router checks the announced
   prefix and origin ASN against this table, and marks the route:

     VALID    -- a matching ROA exists, and the origin ASN and prefix
                 length are both authorized by it.
     INVALID  -- a ROA exists covering this prefix, but the origin ASN
                 or prefix length does NOT match what's authorized.
     NOT FOUND (sometimes "Unknown") -- no ROA exists covering this
                 prefix at all (very common, since RPKI/ROA coverage
                 is still partial across the Internet, per Section 12).

4. Router policy then decides what to DO with each verdict -- RPKI
   itself doesn't force any particular action:

     Common real-world policy: accept VALID and NOT FOUND routes
     normally, but REJECT (or, more cautiously, heavily de-preference
     via LOCAL-PREF, Chapter 49 Section 8) any route marked INVALID.
```

Applied retroactively to Section 6's incident: Pakistan Telecom's bogus `/24` announcement, had ROV been widely deployed and enforced at the time, would have been marked **INVALID** the moment it reached any router checking it — the legitimate ROA for that block would have named YouTube's actual AS as the only authorized origin — and a router enforcing a "reject invalid" policy would have discarded the bogus route immediately, never propagating it any further, regardless of whether PCCW's own filtering happened to catch it or not. This is exactly why RPKI/ROV is described as the real, deployed defense against *this specific class* of incident.

---

## 12. What RPKI Does NOT Fix

It's essential to be honest about RPKI's limits, both because overselling security mechanisms is a real, recurring failure mode in the industry, and because the house style of this course demands it:

- **RPKI only validates the *origin* AS, not the full path.** A route with a perfectly legitimate, ROA-authorized origin ASN can still have a manipulated or leaked *path* to get there — Section 7's Google/Japan incident, for instance, involved routes whose ultimate origin ASNs were legitimate (the real Japanese networks); the problem was an unauthorized *AS in the middle of the path* (Google) inserting itself as a transit hop it had no business providing. Standard RPKI/ROV, as deployed today, does **not** catch this class of leak at all.
- **Coverage is still partial.** As of the mid-2020s, roughly 40-50%+ of announced IPv4 address space (and a growing but still incomplete share of IPv6 space) has a published ROA — meaning a large fraction of the Internet's routes are still "NOT FOUND" under ROV, which by common policy are treated the same as VALID (accepted normally) rather than rejected, precisely because rejecting all unvalidated routes today would break reachability to a very large, still-legitimate portion of the Internet.
- **A misconfigured ROA can itself cause an outage** — for instance, publishing a ROA with too restrictive a max-length, or for the wrong ASN, can cause a network's own *legitimate* announcements to be marked INVALID and rejected by networks enforcing strict ROV policy, which is a real operational risk that makes some network operators cautious about aggressive enforcement.
- **RPKI doesn't stop a legitimate ROA holder from misusing their own authorization**, or from being compromised — if an attacker gains control of an AS that genuinely holds a valid ROA for a prefix, RPKI provides no defense at all against that AS misusing its legitimate authorization.

---

## 13. Beyond RPKI: ASPA, BGPsec, and MANRS

Because Section 12's gap — RPKI validates origin but not path — is well understood, the community has pursued further, less mature mechanisms:

- **BGPsec** (RFC 8205) extends the cryptographic signing idea to the *entire AS-PATH*, having each AS in a path cryptographically sign that it genuinely received and is forwarding the route from the previous AS in the chain, making a fabricated or leaked path detectable, not just a fabricated origin. In practice, BGPsec has seen very little real-world deployment, largely due to the computational cost of signing and verifying on every router for every route, and the requirement that essentially the whole path support it to get meaningful protection — a classic "network effect" adoption problem. **Status: standardized, but not meaningfully deployed.**
- **ASPA (Autonomous System Provider Authorization)** is a newer, lighter-weight IETF proposal specifically targeting route *leaks* (rather than hijacks): each AS publishes a signed statement of which ASes are its legitimate upstream providers, letting a validating router detect a path that implies an impossible or unauthorized provider relationship — closer to directly addressing the Section 7-style leak pattern than RPKI's origin-only validation does. **Status: standardized (RFC 9582, 2024) and actively emerging in early deployment, not yet ubiquitous.**
- **MANRS (Mutually Agreed Norms for Routing Security)** is not a protocol at all, but an industry initiative — a set of concrete, actionable commitments (filtering routes at the network edge, anti-spoofing, coordinating routing information, global validation via RPKI) that participating network operators publicly commit to following, with the explicit goal of reducing exactly the kinds of incidents this chapter describes through better collective operational hygiene, rather than through new cryptography alone. **Status: deployed and actively growing, with several thousand participating networks as of the mid-2020s.**

---

## 14. Packet/Object View: A ROA

A Route Origin Authorization is not a BGP message at all — it's a separate, offline, cryptographically-signed object (technically, a CMS-signed object per RFC 6482) published into the RPKI repository system, entirely independent of the live BGP session. Simplified to its essential logical content:

```
ROA {
  ASID:            65001                  // the authorized origin AS
  ipAddrBlocks: [
    {
      addressFamily: IPv4
      addresses: [
        { address: 203.0.113.0/24, maxLength: 24 }
      ]
    }
  ]
  signature:       <cryptographic signature, chaining up through the
                     issuing RIR's certificate, ultimately to IANA's
                     trust anchor -- structurally parallel to a TLS
                     certificate chain, Chapter 81>
}
```

The `maxLength` field matters precisely because of Section 3's longest-prefix-match danger: a ROA authorizing `203.0.113.0/24` with `maxLength: 24` means AS 65001 may announce *exactly* `/24`, but a more-specific `/25` or smaller announcement — even from the legitimate AS 65001 — would be marked INVALID, specifically to prevent the exact "more-specific sub-block" hijack technique from Section 3 from succeeding even if it were somehow launched *from* the legitimate origin AS's own infrastructure by a compromised or careless actor.

---

## 15. A Real Example: Checking a Prefix's RPKI Status

Public tools let anyone check real, current RPKI validation status for any prefix and ASN pair, no special access required:

```bash
# RIPEstat's RPKI validation API -- check whether a given (ASN, prefix)
# pair would be considered VALID, INVALID, or NOT FOUND:
curl -s "https://stat.ripe.net/data/rpki-validation/data.json?resource=15169&prefix=8.8.8.0/24" \
  | python3 -m json.tool

# Expected structure (abbreviated):
# {
#   "data": {
#     "validating_roas": [
#       {"origin": "15169", "prefix": "8.8.8.0/24", "max_length": 24,
#        "validity": "valid"}
#     ],
#     "status": "valid"
#   }
# }
```

You can also browse **rpki-validator.ripe.net** or **rpki.cloudflare.com** interactively to check any prefix, and see, in real time, exactly which real-world prefixes currently have ROA coverage and which don't — a direct, hands-on look at Section 12's "coverage is still partial" claim.

---

## 16. Hands-On Experiment

```bash
# 1. Check RPKI/ROA coverage and validity for a few real prefixes you
#    care about, comparing a large, security-conscious network against
#    a smaller or older one:
curl -s "https://stat.ripe.net/data/rpki-validation/data.json?resource=15169&prefix=8.8.8.0/24"
curl -s "https://stat.ripe.net/data/rpki-validation/data.json?resource=13335&prefix=1.1.1.0/24"

# 2. Look at RIPEstat's routing history for a prefix to see whether it
#    has ever had an unexpected origin AS appear in its history --
#    a real, hands-on hijack-detection exercise:
curl -s "https://stat.ripe.net/data/routing-history/data.json?resource=8.8.8.0/24" | python3 -m json.tool

# 3. Read a real, historical BGP incident report from a project that
#    monitors global routing anomalies in near-real-time (search for
#    "bgpstream.com" or "Oracle Internet Intelligence Map" incident
#    archives) and identify: was it classified as a leak, a hijack, or
#    both? How long did detection and resolution take?
```

---

## 17. Code: A Minimal Route Origin Validator in Go

This program implements the essential logic of Section 11's Route Origin Validation — given a small in-memory set of ROAs (standing in for what a real validator would fetch from the RPKI repository system) and an incoming route announcement, it returns Valid, Invalid, or NotFound, exactly the three-way verdict a real router's RPKI-to-Router session would receive:

```go
package main

import "fmt"

// ROA mirrors the essential fields of a real Route Origin Authorization
// (Section 14): who may originate this prefix, and up to what length.
type ROA struct {
	Prefix       string
	MaxLength    int
	AuthorizedAS int
}

// Verdict mirrors the three outcomes a real RPKI-to-Router session
// reports to a BGP router (Section 11).
type Verdict int

const (
	NotFound Verdict = iota
	Valid
	Invalid
)

func (v Verdict) String() string {
	return [...]string{"NOT FOUND", "VALID", "INVALID"}[v]
}

// Announcement mirrors the two facts a router actually needs from a
// live BGP UPDATE to run origin validation: the announced prefix
// (with its length) and the origin AS at the end of the AS-PATH.
type Announcement struct {
	Prefix       string
	PrefixLength int
	OriginAS     int
}

// validateOrigin implements Section 11's core check: does any ROA
// cover this exact prefix, and if so, does the announcement's origin
// AS and prefix length fall within what that ROA actually authorizes?
func validateOrigin(roas []ROA, ann Announcement) Verdict {
	found := false
	for _, roa := range roas {
		if roa.Prefix != ann.Prefix {
			continue
		}
		found = true
		if roa.AuthorizedAS == ann.OriginAS && ann.PrefixLength <= roa.MaxLength {
			return Valid
		}
	}
	if found {
		return Invalid // a ROA exists for this prefix, but this announcement doesn't match it
	}
	return NotFound // no ROA at all covers this prefix
}

func main() {
	// The legitimate holder of 203.0.113.0/24 has published a ROA
	// authorizing exactly AS 65001 to announce it, up to a /24.
	roas := []ROA{
		{Prefix: "203.0.113.0/24", MaxLength: 24, AuthorizedAS: 65001},
	}

	scenarios := []Announcement{
		{Prefix: "203.0.113.0/24", PrefixLength: 24, OriginAS: 65001}, // legitimate
		{Prefix: "203.0.113.0/24", PrefixLength: 24, OriginAS: 66666}, // Section 6-style hijack: wrong origin AS
		{Prefix: "203.0.113.0/24", PrefixLength: 26, OriginAS: 65001}, // Section 3-style more-specific sub-hijack
		{Prefix: "198.51.100.0/24", PrefixLength: 24, OriginAS: 65002}, // no ROA published at all for this prefix
	}

	for _, s := range scenarios {
		v := validateOrigin(roas, s)
		fmt.Printf("Announcement %s/%d from AS%d => %s\n",
			s.Prefix, s.PrefixLength, s.OriginAS, v)
	}
	// Output:
	// Announcement 203.0.113.0/24 from AS65001 => VALID
	// Announcement 203.0.113.0/24 from AS66666 => INVALID
	// Announcement 203.0.113.0/26 from AS65001 => INVALID   (exceeds maxLength)
	// Announcement 198.51.100.0/24 from AS65002 => NOT FOUND
}
```

---

## 18. Common Misconceptions

- **"RPKI stops all BGP hijacks."** It stops origin-based hijacks (the Pakistan Telecom and MyEtherWallet pattern) when both properly deployed *and* enforced by the networks in the affected path — it does essentially nothing against path-based leaks like the 2017 Google/Japan incident (Section 12).
- **"These incidents only happen to small, careless networks."** Section 7 is a direct counterexample — one of the most technically sophisticated network operators on Earth produced an incident with real-world, nationwide impact.
- **"A hijack requires special hacking tools."** Mechanically, originating a bogus BGP route requires nothing more exotic than access to a BGP-speaking router and an upstream willing to accept the announcement without filtering — the barrier is operational access and a lack of filtering by others, not sophisticated exploitation.
- **"Once RPKI is 100% deployed, BGP will be fully secure."** Even complete RPKI/ROV deployment leaves path manipulation (Section 12) as an open problem, which is exactly why ASPA and BGPsec (Section 13) exist as separate, additional efforts.

---

## 19. Production Notes

- Major cloud and content networks (Cloudflare, Google, AWS, and others) have publicly committed to, and largely deployed, RPKI ROA publication for their own address space and Route Origin Validation on their own routers — a real, visible shift in industry practice since roughly the late 2010s.
- Some large transit providers now publicly enforce a **"drop invalid" ROV policy** on customer-facing BGP sessions, meaning a customer whose own ROAs are misconfigured can find their legitimate routes silently rejected — a strong, real operational incentive for correct ROA hygiene, and a frequently-cited reason operators are cautious rolling out strict enforcement without first auditing their own ROAs carefully.
- **BGPStream** and similar public monitoring tools/APIs, along with dashboards like **Cloudflare Radar's routing/BGP anomaly views** and the **Oracle (formerly Dyn/Renesys) Internet Intelligence** incident archives, let any operator (or student) subscribe to or browse near-real-time alerts for possible hijacks or leaks affecting prefixes they care about — a practical, deployable piece of Section 9's detection gap being actively narrowed today.
- Regulatory and industry pressure has grown for large transit providers to publicly commit to MANRS-style practices, partly in direct response to the public visibility of incidents like the three covered in this chapter.
- RFC 7908 formally distinguishes several sub-types of route leaks by exactly which relationship boundary was crossed — for example, a route learned from one transit provider re-advertised to a different transit provider, versus a route learned from a customer re-advertised to a peer in a way that violates an agreed scope. Operationally, most networks don't need to memorize the full taxonomy; the practical takeaway is the same one from Chapter 51, Section 15: **only ever re-advertise externally-learned routes to your own customers, never to another peer or transit provider**, and nearly every leak sub-type in the RFC is some variation of that one rule being violated.

---

## 20. What This Chapter Simplified

- RFC 7908's formal route-leak taxonomy has six distinct sub-types; this chapter covers the core mechanical pattern common to most of them rather than enumerating all six.
- The Google/Japan incident's precise internal root cause (the specific automated process or configuration error inside Google's network) was never disclosed by Google in full technical detail publicly; this chapter describes the well-corroborated externally observable facts (what leaked, which networks were affected, roughly how long it lasted) rather than Google's undisclosed internal diagnosis.
- RPKI's underlying certificate hierarchy, repository synchronization protocols (rsync vs. RRDP), and the RPKI-to-Router protocol's full message set are considerably more detailed in production than this chapter's simplified validator model.
- ASPA is presented at a conceptual level; its full validation algorithm (handling both customer and provider directions, "up-ramp/down-ramp" path shape checks) is more intricate than described here.

---

## 21. Interview Questions & Model Answers

**Beginner: "What's the difference between a route leak and a route hijack?"**

*Model answer:* "A hijack is when a network originates a route for a prefix it doesn't legitimately own, impersonating the real holder — like Pakistan Telecom announcing part of YouTube's address space in 2008. A leak is when a network propagates a route beyond where it was supposed to go, usually by re-advertising a route learned from one peer or transit provider onward to another peer or transit provider, effectively offering unintended free transit — like Google accidentally leaking Japanese networks' routes in 2017. The origin AS in a leak is often completely legitimate; it's the path the route travels that's wrong. In practice they often combine, because a hijacked route still needs to leak beyond its intended scope to cause global damage."

**Intermediate: "Explain how RPKI and Route Origin Validation would have prevented the 2008 Pakistan Telecom/YouTube incident, and why it would NOT have prevented the 2017 Google/Japan incident."**

*Model answer:* "In the 2008 incident, Pakistan Telecom announced a /24 for address space it didn't own, with its own ASN as the origin. If YouTube's legitimate ASN had a published ROA for that block, any router checking Route Origin Validation would have seen Pakistan Telecom's announcement fail — wrong origin AS for that prefix — and marked it INVALID, and a router enforcing a drop-invalid policy would never have accepted or propagated it. The 2017 Google incident is different: the routes involved had legitimate origin ASNs — they belonged to real Japanese networks — the problem was that Google's network inserted itself as an unauthorized transit hop in the middle of the path. Standard RPKI/ROV only validates the origin AS at the end of the path, not the AS-PATH as a whole, so it wouldn't flag anything wrong with a route whose origin is correct but whose path includes an AS that shouldn't be there. That gap is exactly what ASPA and BGPsec are trying to close."

**Advanced: "Given RPKI's partial deployment and its origin-only scope, what would a defense-in-depth strategy look like for a network trying to protect itself from both hijacks and leaks today?"**

*Model answer:* "I'd layer several things. First, publish accurate, correctly-scoped ROAs for all of my own address space, including sensible max-length values, so anyone doing Route Origin Validation against my prefixes gets a correct answer. Second, deploy Route Origin Validation myself on inbound sessions and drop or heavily de-preference INVALID routes, to protect my own network from accepting hijacked routes elsewhere on the Internet. Third, implement strict prefix filtering and max-prefix limits on every BGP session with customers and peers, so I don't become the next accidental leak — this is basic MANRS-style hygiene and doesn't require any cryptography at all. Fourth, where available, adopt ASPA to catch leak-style path anomalies that origin validation alone can't see. And finally, subscribe to third-party BGP monitoring for my own prefixes, since even with all of the above, detection speed still matters — all three historical incidents in this chapter took real, damaging hours to detect and fix, and monitoring is what shrinks that window."

---

## 22. Exercises

### Easy

1. In one or two sentences, explain why longest-prefix match (not any BGP policy attribute) is the specific mechanism that made the Pakistan Telecom/YouTube incident possible.
2. What is a ROA, and what single fact does it cryptographically assert?
3. Name the three possible verdicts Route Origin Validation can produce for a given announcement, and explain the difference between the last two.

### Medium

4. Explain, referencing Chapter 51 Section 15's `exportAllowed` function, exactly which rule Google's network violated in the 2017 incident.
5. A network publishes a ROA for `192.0.2.0/24` with `maxLength: 24` and `AuthorizedAS: 65001`. AS 65001 later announces `192.0.2.0/25`. What verdict does Route Origin Validation produce, and why — even though the origin AS is correct?
6. Why did the MyEtherWallet attackers target Amazon's Route 53 DNS infrastructure specifically, rather than attacking MyEtherWallet's own web servers directly?

### Hard

7. Using the `validateOrigin` function from Section 17, extend it to also check a plausible AS-PATH-based rule: reject an announcement if a given "known-bad" ASN appears anywhere in a provided AS-PATH slice, even if the origin AS itself is VALID under RPKI. Explain, in a comment, why this simple extension is not the same thing as real BGPsec or ASPA validation.
8. Research the real, current global percentage of IPv4 address space covered by published ROAs (a specific current or recent figure, not the approximate "40-50%" range given in this chapter), and discuss what practical operational reason might explain why coverage isn't closer to 100% even years after RPKI became widely available.
9. Design, in words, an ASPA-style check that would have caught the 2017 Google/Japan leak, given that standard origin-only RPKI/ROV would not have. What additional piece of published, signed information would a validating router need, beyond a ROA, to detect that Google was an unauthorized transit hop for those Japanese prefixes?

---

## 23. Summary, and the Bridge to Part 8

| Term | Meaning |
|---|---|
| Route leak | A route propagated beyond its intended scope, typically by violating peer/transit export rules |
| Route hijack | An AS originating a route for a prefix it does not legitimately hold |
| Longest-prefix match danger | A more-specific route always wins forwarding, regardless of BGP policy attributes, legitimate or not |
| RPKI | Resource Public Key Infrastructure — a certificate hierarchy proving legitimate address-block ownership |
| ROA | Route Origin Authorization — a signed statement of which AS may originate which prefix, up to what length |
| ROV | Route Origin Validation — checking live BGP announcements against published ROAs (Valid/Invalid/Not Found) |
| BGPsec | Full AS-PATH cryptographic validation; standardized but not meaningfully deployed |
| ASPA | Autonomous System Provider Authorization; targets route leaks specifically; standardized, early deployment |
| MANRS | An industry initiative of operational best practices to reduce leaks and hijacks, not a protocol |

This closes Part 7. Across Chapters 44 through 52, you've followed routing from "what is a router and a routing table" all the way to the business relationships, trust failures, and cryptographic defenses that shape how packets actually cross the real, adversarial, multi-organization Internet — a world Chapters 44-48's single-organization RIP and OSPF were never designed to survive in.

But this entire volume has quietly assumed something else, since as far back as Chapter 35's full LAN trace: that once a router decides "send this packet to next-hop IP address X," some lower-level mechanism *actually knows how to get a frame to that address on the physical wire.* That mechanism has a name, and it's the very first chapter of Part 8: **ARP — the protocol that translates an IP address into the MAC address Ethernet actually needs to deliver a frame (Chapter 53).**
