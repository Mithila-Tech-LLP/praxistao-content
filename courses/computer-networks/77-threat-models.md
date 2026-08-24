# Chapter 77: Threat Models — Thinking Like an Attacker and a Defender

> **"Security is not a feature you add. It's a question you ask about every design decision: who might attack this, with what capability, and what would they gain if they succeeded?"**

---

## Table of Contents

1. [The Assumption This Course Has Been Making](#1-the-assumption-this-course-has-been-making)
2. [Why "Just Encrypt Everything" Isn't an Answer](#2-why-just-encrypt-everything-isnt-an-answer)
3. [What a Threat Model Actually Is](#3-what-a-threat-model-actually-is)
4. [STRIDE, Worked Example: Threat-Modeling a Login Form](#4-stride-worked-example-threat-modeling-a-login-form)
5. [The Cast of Attackers](#5-the-cast-of-attackers)
6. [Attacker Capabilities: Passive vs. Active, On-Path vs. Off-Path](#6-attacker-capabilities-passive-vs-active-on-path-vs-off-path)
7. [What an Attacker Wants: The CIA Triad and Beyond](#7-what-an-attacker-wants-the-cia-triad-and-beyond)
8. [Building the Threat Model for This Volume](#8-building-the-threat-model-for-this-volume)
9. [A Worked Example: Logging Into a Website Over Wi-Fi at a Cafe](#9-a-worked-example-logging-into-a-website-over-wi-fi-at-a-cafe)
10. [A Second Worked Example: An Internet-Connected Home Camera](#10-a-second-worked-example-an-internet-connected-home-camera)
11. [Real-World Incidents That Justify This Threat Model](#11-real-world-incidents-that-justify-this-threat-model)
12. [What Security Cannot Do](#12-what-security-cannot-do)
13. [Common Misconceptions](#13-common-misconceptions)
14. [Connecting Back and Looking Ahead](#14-connecting-back-and-looking-ahead)
15. [Interview Questions & Model Answers](#15-interview-questions--model-answers)
16. [Exercises](#16-exercises)
17. [Summary](#summary)

---

## 1. The Assumption This Course Has Been Making

Go back and reread, in your mind, everything from Chapter 1 through Chapter 76. Ethernet switches learn MAC addresses by trusting whatever source address shows up in a frame (Chapter 31). ARP resolves an IP to a MAC address by broadcasting a question and trusting whoever answers (Chapter 53). DNS resolvers cache whatever answer a nameserver gives them (Chapter 68). TCP establishes a connection by trusting that the packets carrying SYN, SYN-ACK, and ACK actually came from the IP addresses in their headers (Chapter 59). BGP routers announce "I can reach this network" and their neighbors mostly just believe them (Chapter 49).

Every one of those chapters had a quiet, unstated assumption baked into it: **the network is cooperative.** Every device follows the protocol correctly. Every participant is who its addresses claim it is. Nobody is lying, listening in for profit, or trying to redirect your traffic. That assumption is what let those chapters focus entirely on *how do we make this work*, without also asking *how do we stop someone from breaking it on purpose*.

That assumption was reasonable pedagogically — you cannot learn how TCP achieves reliability while simultaneously learning how an attacker forges TCP segments. But it was never true of the real Internet, and this volume exists because the assumption has to be dropped eventually.

The real Internet is not a room full of well-behaved engineers. It is a room full of billions of strangers, a meaningful fraction of whom have both the tools and the motive to lie, eavesdrop, impersonate, and redirect. Coffee shop Wi-Fi is shared with people you've never met. Your packets to a bank in another country cross dozens of routers owned by organizations you have no relationship with and no way to audit. The protocols from Chapters 1–76 will move your data across that hostile territory exactly as faithfully as they'd move it across a friendly one — including your password, your session cookie, and your credit card number, in plaintext, if nothing else is done about it.

This chapter does not teach a specific defense yet. It teaches the *discipline* that Chapters 78–85 are all answers to: threat modeling. Before you can evaluate whether AES-256 or TLS 1.3 or a firewall rule is "secure," you have to be able to answer a prior question — secure *against whom*, doing *what*?

---

## 2. Why "Just Encrypt Everything" Isn't an Answer

A tempting shortcut: skip the threat modeling and just say "encrypt everything, use HTTPS, done." This fails for a concrete, checkable reason: encryption is a tool that answers exactly one kind of threat (someone reading your data in transit) and does nothing at all for several other real threats.

Consider a login form. If you encrypt the connection with TLS but the server itself stores passwords in plaintext in its database, encryption did nothing — the threat was a database breach, not a network eavesdropper, and TLS has no opinion about databases. If you encrypt the connection but an attacker tricks the DNS resolver into pointing `yourbank.com` at their own server (Chapter 68's caching, abused), you will TLS-encrypt your password beautifully — straight to the attacker's own server, whose certificate you didn't check carefully enough. If you encrypt the connection but an attacker has already installed a keylogger on your laptop, no protocol on Earth can save that password, because the threat is at the endpoint, not on the wire.

Each of those is a different attacker, in a different position, with a different capability, defended against by a different mechanism (TLS, PKI/certificate validation, endpoint security). "Encrypt everything" isn't wrong, it's just an incomplete answer to a question you have to actually ask: **what, specifically, are we defending against here?**

That question is what a threat model formalizes.

---

## 3. What a Threat Model Actually Is

**Intuitive level.** A threat model is the security equivalent of "who's it for and what's it protecting against." You wouldn't design a lock for your front door without first deciding whether you're worried about a curious neighbor, a professional burglar, or a nation-state — because each of those calls for a wildly different (and differently expensive) door. Building the strongest possible defense against every conceivable attacker, all the time, is not free: it costs performance, usability, and money. A threat model is how you spend that budget on the risks that are actually real for your situation.

**Engineering terminology.** A threat model is a structured description of:

1. **Assets** — what are we protecting? (a password, a credit card number, the integrity of a firmware update, the availability of a website)
2. **Actors / adversaries** — who might attack it? (a curious hobbyist on the same Wi-Fi, a criminal group, a nation-state intelligence agency, a malicious insider)
3. **Capabilities** — what can each actor actually do? (read traffic passively, inject or modify traffic, control a router on the path, control an endpoint, break cryptography)
4. **Goals** — what does the actor get if they succeed? (steal credentials, redirect money, deny service, plant disinformation)
5. **Out of scope** — what are we explicitly *not* defending against, and why?

That last point matters as much as the first four. A threat model that claims to defend against everything defends well against nothing — it has no priorities. A good threat model states its boundaries out loud: "We defend against a network attacker who can read and modify packets. We do not defend against an attacker who has already compromised the user's device, because at that point the client-side application has already lost."

**Deep technical view.** In formal security engineering, this process is often done with structured frameworks — STRIDE (Section 4), attack trees, or DREAD scoring (Damage, Reproducibility, Exploitability, Affected users, Discoverability) for prioritizing which threats matter most once several have been identified. You do not need to memorize an acronym to threat-model well; what matters is the habit of explicitly naming attacker, capability, and goal before designing a defense. This course will not build a full formal threat-modeling document for every chapter, but every one of Chapters 78–85 is implicitly defending against the threat model this chapter sets up.

---

## 4. STRIDE, Worked Example: Threat-Modeling a Login Form

**STRIDE** is a mnemonic Microsoft developed in the late 1990s that's still the most widely taught systematic way to walk a design and ask "what could go wrong here, categorically?" Each letter names a category of threat, mapped to the property it violates from Section 7's CIA-triad-plus-authenticity vocabulary:

| Letter | Threat category | Property violated | Example against a login form |
|---|---|---|---|
| S | Spoofing | Authenticity | Attacker submits a login request claiming to be a different user |
| T | Tampering | Integrity | Attacker modifies the "amount" field in a POST request mid-transit |
| R | Repudiation | Non-repudiation | User denies having authorized a transaction that they actually did authorize |
| I | Information disclosure | Confidentiality | Password sent in plaintext, or leaked via a verbose error message |
| D | Denial of service | Availability | Attacker floods the login endpoint, locking out real users |
| E | Elevation of privilege | (all of the above, combined) | A regular user manipulates a hidden form field to log in as an administrator |

Working through a real login form with this checklist, rather than freeform brainstorming, catches things an unstructured review often misses. Applying it here: **S**poofing is possible if there's no rate limiting on login attempts (credential stuffing) or if session tokens are predictable. **T**ampering is possible if the form is submitted over plain HTTP, letting an active on-path attacker (Section 5) alter the submitted username. **R**epudiation matters if there's no server-side audit log tying a login to a timestamp and source. **I**nformation disclosure happens if the login page's error message says "no such username" versus "wrong password" — leaking which accounts exist even when authentication itself never succeeds. **D**enial of service happens if failed-login lockouts can be triggered by an attacker who doesn't know the real password, deliberately locking a target out. **E**levation of privilege happens if the role/permission check happens in JavaScript on the client rather than being re-verified on the server.

Notice that every single one of these threats already has a specific chapter in this volume aimed at it: TLS (Chapter 82) for tampering and disclosure over the wire, rate limiting and firewalls (Chapter 84) for denial of service, and PKI (Chapter 81) plus proper server-side session design for spoofing and elevation. STRIDE doesn't invent new defenses — it's a checklist that reliably surfaces which of the defenses this volume builds are actually needed for a given design, and it scales to systems far more complex than a login form (an entire microservice architecture, an IoT product, a payment pipeline).

---

## 5. The Cast of Attackers

Different chapters in this volume defend against different members of this cast. Naming them precisely now avoids vague language like "hackers" later.

**The passive eavesdropper.** Someone who can *see* your traffic but cannot change it. Classic examples: another user on the same open Wi-Fi network running a packet capture tool, an ISP logging which sites you visit, an intelligence agency tapping an undersea cable (Chapter 23). A passive eavesdropper's goal is almost always information disclosure — reading a password, a message, a document — not disruption.

**The active on-path attacker.** Someone positioned somewhere your traffic actually flows through — a router, a Wi-Fi access point they control, a compromised switch — who can not only read your packets but also drop, delay, duplicate, reorder, or *modify* them before they reach the destination. This is strictly more dangerous than a passive eavesdropper: it enables man-in-the-middle attacks (Chapter 83), where the attacker impersonates each side to the other, silently relaying (and altering) everything in between.

**The off-path attacker.** Someone who is *not* sitting on the actual path your packets travel, but who can still inject forged packets into the conversation — for example, by guessing sequence numbers (Chapter 60) or by racing a legitimate DNS response (Chapter 68). Off-path attacks are generally harder to pull off than on-path ones because the attacker is working blind, guessing values instead of reading them directly, but they matter because they don't require compromising any network infrastructure at all — just being on the same Internet as you.

**The malicious (or compromised) endpoint.** One of the two parties in the conversation — or software running on one of their machines — is itself the attacker, or has been taken over by one. A protocol that perfectly secures data in transit does nothing here, because the threat is inside the room, not on the road between rooms. This is why "secure the network" and "secure the endpoint" are separate, both-necessary disciplines.

**The malicious infrastructure operator.** A step beyond a compromised router: an entity that *legitimately* operates a piece of the path — an ISP, a Wi-Fi hotspot provider, a certificate authority (Chapter 81) — but chooses to abuse that trusted position. This is a genuinely different threat from an outside attacker breaking in, because the actor already has the access; the defense has to assume trusted parties can misbehave, not just that untrusted parties might intrude.

**The credentialed insider.** Someone who legitimately holds valid credentials — an employee, a contractor, a compromised legitimate account — and uses that access beyond its intended purpose. Network-layer defenses (encryption, authentication) are largely powerless here by design, since an insider already passes those checks; this threat is addressed by access control, least-privilege design, and audit logging, mostly outside this volume's scope, though Chapter 101's service mesh and mTLS touch on limiting *lateral* movement once an insider or compromised service is inside a network.

---

## 6. Attacker Capabilities: Passive vs. Active, On-Path vs. Off-Path

It's worth separating *position* (where is the attacker relative to your traffic?) from *capability* (what can they do from there?) because the two combine into a 2x2 grid that predicts exactly which defense is needed.

```
                    PASSIVE (read only)          ACTIVE (read + modify/inject)

ON-PATH             Traffic sniffing              Man-in-the-middle,
(sees your          (coffee-shop Wi-Fi            packet injection,
actual packets)     sniffer, ISP logging,         TLS-stripping,
                    tapped cable)                 DNS response forgery
                    -> defeated by encryption      -> defeated by encryption
                       (Ch 78-79)                     PLUS authentication
                                                       (Ch 80-82)

OFF-PATH             Traffic analysis               Blind spoofing,
(does not see        (timing/size inference,       sequence-number guessing,
your actual          rare, defended in Ch 83)      DNS cache poisoning races
packets)                                            -> defeated by
                                                        unpredictability +
                                                        authentication (Ch 80-83)
```

Notice the pattern already emerging, ahead of the chapters that will justify it in full: a **passive** attacker is defeated by making the data unreadable (encryption — Chapters 78–79). An **active** attacker requires something stronger: not just hiding the content, but *proving* it hasn't been tampered with and *proving* who really sent it (integrity and authentication — Chapters 80–82). Confidentiality alone is not enough once you assume an attacker who can inject and modify.

This is precisely why TLS, previewed here and built in full in Chapter 82, bundles three separate guarantees rather than just one:

| Guarantee | Question it answers | Defeats |
|---|---|---|
| Confidentiality | Can anyone else read this? | Passive eavesdropper |
| Integrity | Was this altered in transit? | Active on-path tamperer |
| Authenticity | Am I really talking to who I think I am? | Active impersonator / MITM |

A protocol that only provides confidentiality (say, plain AES encryption with no authentication, foreshadowing a warning that Chapter 78 will make explicit) is secure against the top-left box of the grid above and nothing else. That's a threat-modeling failure waiting to happen if your actual adversary is in the top-right box.

It's also worth being honest about *where*, physically and organizationally, an on-path position actually comes from in practice, since "on-path attacker" can otherwise sound abstract: a shared Wi-Fi access point (anyone else connected to it, or the AP operator itself); any router or switch between you and the destination (your home router, your ISP's edge router, a transit provider's backbone router, a peering exchange — Chapter 51); a compromised piece of network infrastructure (a hacked router firmware, a rogue employee at an ISP); or a nation-state with lawful or unlawful access to a backbone link or IXP (Chapter 51). The number of organizations and individuals who are, at some point, "on-path" for a typical Internet connection between two continents is easily in the dozens.

---

## 7. What an Attacker Wants: The CIA Triad and Beyond

Security engineering has a standard vocabulary for *what's at stake*, independent of who the attacker is or where they sit. It's usually taught as the **CIA triad**:

- **Confidentiality** — keeping data secret from those not authorized to see it. (Your password in transit; a company's trade secrets.)
- **Integrity** — ensuring data isn't modified without detection, whether in transit or at rest. (A firmware update that hasn't been tampered with; a bank transfer amount that wasn't silently changed from $10 to $10,000.)
- **Availability** — ensuring a system keeps working for legitimate users. (A website that stays up during a DDoS attack, previewed in Chapter 83.)

Two more properties matter constantly in the chapters ahead and are worth naming now because they're easy to conflate with the ones above:

- **Authenticity** — proving that data or a message really came from who it claims to (not the same as confidentiality: a message can be public and unencrypted, yet you still need to know it's genuinely from its claimed sender — this is exactly what a digital signature in Chapter 80 provides).
- **Non-repudiation** — ensuring the sender of an authenticated message cannot later credibly deny having sent it. This matters for legal and financial systems (a digitally signed contract) more than for everyday browsing, but it's a direct consequence of good digital-signature design.

Every attacker in Section 5 is ultimately going after one or more of these five properties. A passive eavesdropper wants your *confidentiality*. An active MITM wants to break your *integrity* and *authenticity* simultaneously (to inject a fake page while impersonating the real server). A DDoS attacker (Chapter 83) wants your *availability*. Naming the property under threat is often the fastest way to identify which chapter's tool actually applies.

These five properties also give a fast way to sanity-check a proposed defense: ask, out loud, exactly which of the five it improves. "We added a WAF (Chapter 84)" mostly buys availability and some integrity against application-layer attacks, but does nothing for confidentiality of data already inside the network. "We added TLS" buys confidentiality, integrity, and authenticity for data in transit, but nothing for data at rest or for availability under a volumetric DDoS flood. A security review that can't name which property a given control improves is usually a sign the control was added out of habit rather than in response to an actual identified threat.

---

## 8. Building the Threat Model for This Volume

Putting Sections 3–7 together, here is the explicit threat model that Chapters 78 through 85 are collectively answers to. Writing it out plainly, the way a real security design document would:

**Assets:** the confidentiality and integrity of application data crossing an untrusted network (credentials, personal data, financial transactions, page content); the authenticity of the server a client believes it's talking to.

**In-scope adversary:** a network-level attacker who may be passive (can read all traffic on some link the conversation crosses) or active (can additionally drop, delay, duplicate, reorder, inject, or modify any packet on that link — i.e., a full on-path Dolev–Yao-style attacker, the standard adversary model used in academic cryptographic protocol analysis). This adversary does **not** need physical access to either endpoint. This single adversary model is strong enough to already justify everything from Chapter 78 (why you need encryption at all) through Chapter 82 (why the TLS handshake is shaped the way it is).

**Also in-scope, addressed individually:** an off-path attacker attempting blind injection or spoofing (Chapter 83); a resource-exhaustion attacker attempting denial of service (Chapter 83, mitigated partly in Chapter 84); a malicious or negligent Certificate Authority abusing its trusted position (Chapter 81).

**Explicitly out of scope for this volume:** a compromised endpoint (malware already running on the client or server); physical theft of a device; social engineering of a human; supply-chain attacks on hardware or software before it reaches the network; a credentialed insider abusing legitimate access. These are all real, serious threats — but they are defended against by completely different mechanisms (endpoint security, physical security, secure software supply chains, access control) that live outside the scope of a *networking* course. Naming them as out of scope is not an oversight; it's the discipline from Section 3 in action.

This is exactly the kind of explicit statement every real security design should make, and it's why this chapter had to come before any cryptography: without it, Chapter 78's AES and Chapter 82's TLS handshake would look like arbitrary machinery instead of specific, motivated answers to a specific, stated adversary.

---

## 9. A Worked Example: Logging Into a Website Over Wi-Fi at a Cafe

Concrete scenario: you sit down at a cafe, join their open (no-password) Wi-Fi, and log into your email.

```
Your laptop                Cafe Wi-Fi AP           Cafe's router        ISP        Email server
     |                           |                        |               |             |
     |-- HTTP login form ------->|                        |               |             |
     |   (if unencrypted)        |------------------------|-------------->|------------>|
     |                           |                        |               |             |
Attacker sitting at the next table, same open Wi-Fi:
  - PASSIVE capability: runs a packet sniffer, sees every unencrypted frame the AP
    forwards, including yours (Chapter 30's "shared medium" problem resurfacing as
    a security problem, not just a performance one).
  - If you're using plain HTTP: username and password are sitting in the payload,
    in cleartext, decodable with a five-minute-old free tool.
  - ACTIVE capability (same attacker, slightly more effort): runs a rogue DHCP
    server (abusing Chapter 55's trust assumption) or ARP-spoofs the gateway
    (abusing Chapter 53's trust assumption) to become an on-path
    man-in-the-middle, and can now also *modify* the login page before it
    reaches you -- e.g., injecting JavaScript that copies your password as
    you type it, even if the site later loads over HTTPS.
```

Walking through the CIA-triad-plus-authenticity checklist from Section 7 for this one scenario: confidentiality is broken the instant the connection is plaintext HTTP. Integrity and authenticity are *both* still at risk even with HTTPS, if the attacker can trick you into trusting a fraudulent certificate, or if you dismiss a browser warning without reading it (the browser is trying to enforce Chapter 81's trust chain on your behalf — heed it). Availability isn't really this attacker's goal, so a firewall or DDoS mitigation (Chapter 84) would be effort spent on the wrong threat.

The fix that actually addresses the attacker described here — passive-or-active, on the local Wi-Fi segment — is exactly TLS: it makes the payload unreadable (defeats the passive sniffer), it makes tampering detectable (defeats the active MITM's modification), and it cryptographically authenticates the server via its certificate (defeats impersonation, provided you don't click through the warning). Every mechanism named in that sentence is unexplained machinery right now. Chapters 78 through 82 build each one, one piece at a time, specifically to answer this scenario.

---

## 10. A Second Worked Example: An Internet-Connected Home Camera

A second scenario, deliberately different in shape, to show the same discipline generalizes beyond "browser talks to website."

**Assets:** the video feed itself (confidentiality — nobody wants a stranger watching their living room); the camera's control channel (integrity/authenticity — nobody wants a stranger able to pan the camera or disable recording); the camera's availability (a family relying on it for security shouldn't have it knocked offline easily).

**In-scope adversaries and what STRIDE-style analysis (Section 4) surfaces:**

- **Spoofing:** if the camera authenticates to the cloud backend using a hardcoded, factory-default credential shared across every unit of that model, an attacker who extracts it from one device (or finds it in a leaked firmware image) can spoof *any* camera of that model to the backend.
- **Tampering:** if the video stream isn't encrypted end-to-end, an on-path attacker (a compromised home router, a malicious ISP, a rogue actor on the local Wi-Fi) could substitute a static, pre-recorded frame for the live feed — a well-documented real-world attack class against poorly designed IoT cameras.
- **Information disclosure:** if the video stream is sent unencrypted over the local network (common in cheaper devices that assume the LAN is "safe"), anyone else on that LAN, or anyone who compromises one other device on it, can passively watch the feed.
- **Denial of service:** many cheap IoT devices have no rate limiting or resource protection at all, making them trivial to knock offline with a small flood of malformed requests — and, notoriously (Section 11), trivial to recruit *into* a much larger denial-of-service attack against someone else entirely.

**Explicitly out of scope for the camera's own threat model** (though very much in scope for the vendor's overall security program): a thief physically stealing the camera and extracting keys from its flash storage; a malicious firmware update pushed by a compromised vendor build server. Naming these as out of scope for a *network-focused* threat model isn't ignoring them — it's correctly routing them to the mechanisms designed for them (physical tamper-resistance, secure firmware signing and supply-chain controls) rather than expecting TLS to solve a hardware problem.

The lesson this second example adds: the exact same five-property, position-and-capability framework from Sections 5–7 applies whether the "client" is a browser with a human at the keyboard or a $30 embedded device with no human present at all — and IoT devices, as Section 11 shows, are disproportionately where threat modeling gets skipped, with large-scale real-world consequences.

---

## 11. Real-World Incidents That Justify This Threat Model

Threat modeling can sound like an abstract exercise until it's grounded in things that actually happened, at scale, because these exact threats were left unaddressed.

**Firesheep (2010).** A browser extension released as a deliberate, public demonstration: it passively sniffed unencrypted Wi-Fi traffic at, say, a coffee shop, and automatically hijacked other users' logged-in sessions on sites like Facebook and Twitter — because those sites, at the time, only encrypted the login form itself, then dropped back to plain HTTP for the rest of the session, leaving the session cookie (Chapter 72) sitting in cleartext for any passive eavesdropper (Section 5) to steal and reuse. This single, widely publicized tool is a large part of why the industry shifted to encrypting entire sessions end-to-end ("HTTPS everywhere") rather than just the login form.

**The Mirai botnet (2016).** Hundreds of thousands of IoT devices — cameras, routers, DVRs, almost exactly the class of device in Section 10's worked example — were compromised at scale using nothing more sophisticated than a list of about sixty common factory-default username/password combinations that owners had never changed. The compromised devices were then used to launch some of the largest denial-of-service attacks recorded at the time, taking down major DNS provider Dyn (Chapter 66-69's DNS infrastructure) and, with it, a large portion of the Web for millions of users who had nothing to do with any camera or DVR. This incident is a direct, large-scale illustration of Section 10's spoofing threat (shared default credentials) cascading into an availability attack against unrelated third parties.

**DNS cache poisoning (Kaminsky attack, 2008).** Security researcher Dan Kaminsky demonstrated a practical, fast way for an off-path attacker (Section 5) to win the "race" between a forged DNS response and the real one, exploiting DNS's historically limited randomness in query IDs and source ports (Chapter 68's caching mechanism). A successful attack could redirect an entire domain's traffic to a server the attacker controlled, for as long as the poisoned, forged record sat in resolver caches. This incident drove a coordinated, industry-wide emergency patch (adding much stronger source-port randomization) and accelerated interest in DNSSEC (Chapter 69).

Each of these incidents maps precisely onto Sections 5–8's vocabulary: a passive eavesdropper exploiting a confidentiality gap (Firesheep), a spoofing threat cascading through weak default authentication into a global availability attack (Mirai), and an off-path attacker exploiting insufficient unpredictability in a trusted protocol (Kaminsky). None of them required breaking any cryptography — they exploited exactly the "cooperative network" trust assumptions Section 1 named, at real, Internet-wide scale.

---

## 12. What Security Cannot Do

An honest threat model is as much about limits as about coverage. Nothing in this volume will:

- Protect you if your own device is already compromised (malware can read your screen and keystrokes regardless of how well-encrypted the network is).
- Protect you from a service that stores your data insecurely at rest, even if it received that data over a perfectly secure channel.
- Make a phishing site "not dangerous" just because it has a valid TLS certificate — TLS proves you're talking to the domain in the address bar over an unreadable, untampered channel; it says nothing about whether that domain is operated by someone honest. `paypa1-secure-login.com` can have flawless TLS and still be a scam.
- Eliminate metadata leakage — even with perfect encryption, an eavesdropper on the path can usually still see *that* you connected to a given IP address, roughly *when*, and roughly *how much* data moved, unless additional tools (VPNs, Tor, traffic padding) are layered on top. Encryption hides content; it doesn't automatically hide the fact of communication.
- Substitute for good key management. Every mechanism in Chapters 78–82 assumes secret keys stay secret. A leaked private key (through a software bug, a careless developer, or a stolen laptop) breaks all of the cryptography built on top of it, no matter how strong the algorithm.
- Fix a weak default credential. As Mirai (Section 11) demonstrates at scale, the strongest cryptography in the world does nothing if the authentication step it protects can be bypassed with `admin`/`admin`.

---

## 13. Common Misconceptions

**"HTTPS means the site is legitimate."** No — it means the connection to that specific domain is encrypted and the certificate was validated. As Section 12 notes, attackers can and do get valid certificates for phishing domains.

**"My data is only at risk on public Wi-Fi."** A wired connection, your home Wi-Fi, and your ISP's network can all have passive or active attackers too (a compromised router, a malicious insider at the ISP, a nation-state tap on a backbone link). Public Wi-Fi is just the easiest and cheapest case to demonstrate, not the only real one.

**"Encryption solves security."** As Section 2 showed, encryption solves exactly one property (confidentiality). Integrity, authenticity, and availability each need their own mechanisms, and endpoint and human-layer threats need entirely different mechanisms outside this volume's scope.

**"If it's not encrypted, it's automatically insecure."** Not always true in context — some data genuinely doesn't need confidentiality (a public DNS query for a well-known domain isn't secret), but almost always still benefits from integrity and authenticity (you still want to know the DNS answer wasn't forged — see DNSSEC in Chapter 69 and DNS poisoning in Chapter 83). Threat modeling means matching the mechanism to the actual property at risk, not reflexively encrypting everything and calling it done.

**"Threat modeling is a one-time document you write and file away."** Real threat models are living documents, revisited whenever the system changes — a new feature, a new integration, a new class of user — because each change can introduce a new asset, a new attacker capability, or invalidate an "out of scope" assumption that used to be safe.

---

## 14. Connecting Back and Looking Ahead

Every trust assumption named in this chapter as "abused" by an attacker was a real mechanism you already learned, working exactly as designed, for a cooperative network: MAC learning (Chapter 31), ARP (Chapter 53), DHCP (Chapter 55), DNS caching (Chapter 68), TCP sequence numbers (Chapter 60), and BGP announcements (Chapter 49) are not *broken* protocols — they were correctly designed for the threat model of 1970s–1990s research and enterprise networks, where "everyone on this network is roughly trustworthy" was a reasonable working assumption. The gap between that old assumption and today's hostile, billions-of-strangers Internet is exactly what this volume closes.

The path forward: Chapter 78 starts with the single most direct tool for the confidentiality problem (symmetric cryptography), immediately runs into a hard problem it creates (how do two strangers share a secret key over a network an eavesdropper is watching?), and hands that problem to Chapter 79 (asymmetric cryptography) to solve. Chapter 80 adds integrity and authenticity on top (hashing and digital signatures). Chapter 81 closes the remaining gap — proving a key belongs to a *specific, real-world identity* like `google.com`, not just to *some* keyholder. Chapter 82 assembles all four into the TLS handshake that protects the vast majority of the Web today.

---

## 15. Interview Questions & Model Answers

**Beginner: "What is a threat model, and why do you need one before designing a security control?"**

A threat model is an explicit statement of what you're protecting (assets), who might attack it (adversaries), what those adversaries are capable of (passive reading vs. active tampering, network-level vs. endpoint-level access), and what's explicitly out of scope. You need one before choosing a defense because defenses are not universal — encryption defeats a passive eavesdropper but does nothing against a compromised endpoint, so picking a control without first naming the threat usually means solving the wrong problem, or over-engineering against a threat that isn't realistic for the situation.

**Beginner: "What does STRIDE stand for, and why is it useful?"**

STRIDE stands for Spoofing, Tampering, Repudiation, Information disclosure, Denial of service, and Elevation of privilege — six categories of threat, each mapped to a specific security property being violated. It's useful because it turns "think of everything that could go wrong" (open-ended and easy to miss things) into a systematic checklist applied to each part of a design, which reliably surfaces threats — like a missing rate limit, or a client-side-only permission check — that unstructured brainstorming tends to skip.

**Intermediate: "What's the practical difference between a passive and an active network attacker, and why does it matter for protocol design?"**

A passive attacker can only observe traffic — read packets, log metadata — without altering delivery. An active attacker can additionally drop, delay, duplicate, reorder, inject, or modify packets. The difference matters because confidentiality (encryption) alone fully defeats a passive attacker, but an active attacker can still cause damage even against encrypted traffic — for example, replaying old encrypted packets, or dropping packets selectively to trigger unwanted retransmission behavior — unless the protocol also provides integrity (detecting modification) and often freshness/replay protection. This is why TLS bundles confidentiality, integrity, and authenticity together instead of shipping encryption alone.

**Intermediate: "What lesson does the Mirai botnet teach about the limits of network security controls?"**

Mirai compromised hundreds of thousands of IoT devices using only their unchanged factory-default credentials — no cryptography was broken at all. It demonstrates that even a perfectly designed network protocol is defeated if the authentication step behind it is weak (shared, guessable, or never-rotated credentials), and that a threat model has to explicitly cover authentication strength, not just assume "we used TLS, so we're secure." It also shows that a weak link in one product (a $30 camera) can cascade into an availability attack against completely unrelated third parties (Dyn's DNS infrastructure and everything depending on it), which is why threat models should consider the blast radius of a compromised asset, not just the asset owner's own risk.

**Advanced: "How would you threat-model a new internal microservice API that currently has no authentication, assuming it will eventually be exposed to traffic from outside the data center?"**

Start by naming the assets (whatever data or actions the API exposes) and enumerate credible adversaries at each deployment stage: internally, a compromised neighboring service or a malicious insider on the same network segment; once exposed externally, add the full Dolev-Yao on-path/off-path network adversary from Section 8, plus unauthenticated internet-wide scanners and credential-stuffing bots. Map each adversary to which CIA-triad-plus-authenticity property they threaten: an internal passive listener threatens confidentiality if traffic is plaintext (mitigate with mTLS, foreshadowing Chapter 101's service mesh pattern); an external attacker with no authentication threatens every property simultaneously, since "no authentication" means anyone can claim to be an authorized caller. The concrete deliverable is a short document naming these adversaries and mapping each to a specific planned control (mTLS for transport, signed/short-lived tokens for caller identity, rate limiting for availability) — and explicitly stating what's still out of scope (e.g., "we are not defending against a fully compromised caller service with valid credentials; that's covered by least-privilege authorization, tracked separately").

---

## 16. Exercises

### Easy

1. For each of the following, name whether the attacker described is passive or active, and on-path or off-path: (a) an ISP logging which websites its customers visit; (b) someone at a cafe running a rogue Wi-Fi access point with the same name as the real one; (c) someone thousands of miles away trying to guess a TCP sequence number to inject a forged packet.
2. Explain in your own words why "the connection is encrypted" is not the same claim as "the connection is with who I think it's with."
3. List the CIA-triad-plus-authenticity property most directly threatened by: a keylogger, a DDoS attack, an ARP-spoofing MITM, and a stolen unencrypted laptop backup.
4. Apply STRIDE (Section 4) to a simple "forgot password" email-reset flow. Name at least one concrete threat under three different STRIDE letters.

### Medium

5. Write a two-paragraph threat model (assets, adversary, in-scope, out-of-scope) for a home smart lock that connects to the internet through a Wi-Fi router. Be explicit about what you are choosing not to defend against and why that's a reasonable engineering trade-off.
6. A team proposes: "We'll encrypt data in transit with TLS, so we don't need to worry about authentication." Identify the specific threat-model gap in this reasoning, using a concrete attacker scenario.
7. Explain why an attacker who can only observe encrypted traffic (not decrypt it) might still learn something useful, and name two pieces of information that leak even through strong encryption.
8. Using Section 10's home camera example, explain why "the video stream is encrypted end-to-end" and "the camera authenticates itself securely to the backend" are two separate security properties, and describe an attack that defeats one without defeating the other.

### Hard

9. Research (or recall from earlier chapters) one real, historical incident that exploited a "cooperative network" trust assumption named in Section 14 (an ARP spoofing incident, a BGP route hijack — Chapter 52 — or a DNS cache poisoning attack, potentially the Kaminsky attack from Section 11). Describe which threat-model box (Section 6's grid) the attacker occupied, and which specific chapter-78-through-85 mechanism would have prevented or mitigated it.
10. Design a threat model for a video call application, distinguishing what an on-path network attacker could achieve versus what a malicious participant already inside the call could achieve. Explain why these require fundamentally different defenses, and why no amount of TLS on the network connection addresses the second category.
11. Firesheep (Section 11) worked because sites encrypted only the login page, not the rest of the session. Explain, using Section 6's passive/active grid, exactly what capability the attacker needed to exploit this, and why simply making the login page use HTTPS (without also protecting the rest of the session) left the underlying threat completely unaddressed.

---

## Summary

| Term | Meaning |
|---|---|
| Threat model | Explicit statement of assets, adversaries, capabilities, goals, and out-of-scope threats |
| STRIDE | Spoofing, Tampering, Repudiation, Information disclosure, Denial of service, Elevation of privilege |
| Passive attacker | Can read/observe traffic but not alter its delivery |
| Active attacker | Can drop, delay, duplicate, reorder, inject, or modify traffic |
| On-path attacker | Positioned somewhere the real traffic physically flows through |
| Off-path attacker | Not on the traffic's path; must inject/guess to interfere |
| Man-in-the-middle (MITM) | Active on-path attacker impersonating each side to the other |
| Confidentiality | Data is unreadable to unauthorized parties |
| Integrity | Data has not been altered undetected |
| Availability | The system keeps working for legitimate users |
| Authenticity | Proof that data/identity is genuinely who/what it claims |
| Non-repudiation | Sender cannot credibly deny having sent an authenticated message |
| Dolev-Yao adversary | Standard academic model: attacker controls the network fully (read/write/inject) but not the endpoints |

Chapter 77 named the enemy and the battlefield: from here on, every mechanism in this volume assumes a network-level attacker who can read and tamper with anything in transit. Chapter 78 starts building the first defense — symmetric cryptography — and runs headfirst into the problem that makes Chapter 79 necessary.
