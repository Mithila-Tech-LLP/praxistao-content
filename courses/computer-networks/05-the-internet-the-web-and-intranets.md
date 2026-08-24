# Chapter 05: The Internet, the Web, and Intranets — Untangling the Terms

> *"'The internet is down' and 'the website is down' are, almost always, two completely different problems — and you cannot tell which one you have until you understand that the internet and the Web were never the same thing."*

---

## Table of Contents

1. [The Confusion, Stated Plainly](#1-the-confusion-stated-plainly)
2. [Why People Conflate Them — A Reasonable Mistake](#2-why-people-conflate-them--a-reasonable-mistake)
3. [What the Internet Actually Is](#3-what-the-internet-actually-is)
4. [What the Web Actually Is](#4-what-the-web-actually-is)
5. [Other Applications That Use the Internet But Aren't the Web](#5-other-applications-that-use-the-internet-but-arent-the-web)
6. [What an Intranet Actually Is](#6-what-an-intranet-actually-is)
7. [Why Layering Is What Makes These Separable](#7-why-layering-is-what-makes-these-separable)
8. [Deep Dive: The Same URL, Different Reachability](#8-deep-dive-the-same-url-different-reachability)
9. [Production Notes: VPNs, Extranets, and the Blurring Boundary](#9-production-notes-vpns-extranets-and-the-blurring-boundary)
10. [A Worked Scenario: Diagnosing "Internet Down" vs. "Website Down"](#10-a-worked-scenario-diagnosing-internet-down-vs-website-down)
11. [Hands-On Experiment: Prove the Separation Yourself](#11-hands-on-experiment-prove-the-separation-yourself)
12. [Common Misconceptions](#12-common-misconceptions)
13. [Connections Backward and Forward](#13-connections-backward-and-forward)
14. [Interview Questions & Model Answers](#14-interview-questions--model-answers)
15. [Exercises](#15-exercises)
16. [Summary](#16-summary)

---

## 1. The Confusion, Stated Plainly

Say the sentence "the internet is down" to almost anyone, and picture what they imagine: usually, a browser that won't load any web page. Say "check if google.com is down," and most people will open a browser and try to load it — using, again, the Web. For the overwhelming majority of people, in practice, "the internet" and "the Web" refer to the identical daily experience: typing something into a browser and getting a page back.

They are not the same thing. They were never the same thing. And the difference is not pedantic trivia — it's the exact reason a network engineer's first diagnostic question, when someone reports "the internet is down," is almost always "down for everything, or just for one website?" — because the answer determines whether the problem is in the Internet's infrastructure (Chapters 4, 6, and Volumes 5–8) or in one specific application running on top of it (Volume 11, on HTTP and the Web). This chapter exists to make that distinction as sharp and as useful as a working network engineer actually needs it to be.

---

## 2. Why People Conflate Them — A Reasonable Mistake

The confusion isn't foolish — it has a real historical cause, and understanding that cause is itself useful (Chapter 13 tells this story properly). For most people who came online any time after the mid-1990s, the Web was the very first, and for years the *only*, thing they ever did with an Internet connection. Email existed earlier and separately, but for a huge number of users, "going online" meant, specifically, opening a browser. As smartphones made apps (which also use the Internet, but not the Web, as Section 5 shows) the dominant way people spend time online, the same conflation just shifted shape: "the internet" now often means "whatever app or website I'm currently trying and failing to use," regardless of whether that thing is technically part of the Web at all.

This chapter's job is to pull these two ideas apart cleanly, using the vocabulary Chapters 3 and 4 already built.

---

## 3. What the Internet Actually Is

Using Chapter 3's definition directly: the Internet is a network — a (very large) set of computers connected by shared, addressable communication links, letting any computer reach any other by specifying an address. Using Chapter 4's vocabulary: it is, structurally, an enormous WAN, built by interconnecting an enormous number of separately-owned LANs (homes, offices, university campuses, data centers) through leased long-distance infrastructure (undersea cables, long-haul fiber, satellite links).

Crucially, notice what this definition does **not** mention: web pages, browsers, HTML, or `https://`. Nothing about "network of interconnected, addressable computers" requires the Web to exist at all. And historically, it didn't for decades — the Internet's direct ancestor, ARPANET (Chapter 10), was running in 1969, moving real data between real computers, more than twenty years before the Web existed. The Internet is the plumbing: physical wires and radio links, the addressing scheme that lets one machine find another (Chapter 36 and beyond), and the rules (protocols) for moving data reliably across that plumbing (Volumes 6–9). It is infrastructure, in the same sense that a city's road network or electrical grid is infrastructure — it exists to be used *for* things, and it doesn't particularly care what those things are.

```
                          THE INTERNET
     (addressing, routing, physical links — the "plumbing")

  +------------------------------------------------------------+
  |                                                              |
  |   Any two connected computers can reach each other, using   |
  |   addresses (Ch.36+) and routes (Volume 7) — regardless of  |
  |   what application is using that connection.                |
  |                                                              |
  +------------------------------------------------------------+
```

---

## 4. What the Web Actually Is

The **World Wide Web** is one specific application that runs *on top of* the Internet — not a separate network, not an alternative to the Internet, but a particular way of using the connectivity the Internet already provides. It was invented by Tim Berners-Lee in 1989–1991 at CERN, more than two decades after the Internet's predecessor network began operating, specifically to let researchers share and link documents to each other easily (Chapter 13 covers this history properly).

The Web, concretely, is the combination of three ideas working together:

1. **Hypertext documents** — pages of content (originally text, now much richer) that can contain links to other pages, letting a reader jump between related documents instantly.
2. **A common protocol for requesting and delivering those documents** — HTTP (HyperText Transfer Protocol), which defines the rules a browser and a server follow to ask for a page and receive it (the full subject of Volume 11, starting at Chapter 70).
3. **A universal addressing scheme for documents** — the URL (Uniform Resource Locator), which identifies exactly which document, on which server, a browser is asking for.

None of this is what makes two computers able to reach each other in the first place — that's the Internet's job, already done, one layer down, by the time HTTP ever gets involved. The Web simply assumes that connectivity already exists and adds a specific, agreed-upon *application* — documents, links, and a request/response protocol — on top of it. This is precisely the layering idea Chapter 1 first hinted at (the shared code can be swapped without touching the physical channel underneath) and that Chapter 7 formalizes for the whole course.

### Deep dive: URL schemes as an explicit protocol selector

It's worth noticing something that hides in plain sight every time you type a web address: the `https://` (or `http://`) at the start of a URL is not decoration — it's an explicit instance of Chapter 1's shared-code idea, telling your browser exactly which protocol to speak before it sends a single byte. This is called the **URL scheme**, and the Web's `http`/`https` schemes are only two of many:

| Scheme | What it tells the computer to do |
|---|---|
| `http://`, `https://` | Speak HTTP (Chapter 71) to fetch a web page — `https` additionally wraps the connection in TLS encryption (Chapter 82) |
| `mailto:` | Hand off to the user's configured email application, not a browser rendering a page at all |
| `ftp://` | Speak FTP, an older file-transfer protocol, entirely unrelated to HTTP |
| `ssh://` | Open a secure remote terminal session (Volume 17 covers building tools like this) |
| `ws://`, `wss://` | Open a WebSocket (Chapter 76) — a persistent, two-way connection, a different traffic pattern from an ordinary web page request |

A small Go program makes the dispatch logic concrete — this is, in miniature, exactly what a browser or operating system does the instant you click a link:

```go
package main

import (
	"fmt"
	"strings"
)

// dispatchByScheme models how a browser decides what to do with a URL
// before ever making a network connection — a direct, real instance of
// the shared-code idea from Chapter 1 applied to protocol selection.
func dispatchByScheme(url string) string {
	switch {
	case strings.HasPrefix(url, "https://"):
		return "Open a TLS-encrypted HTTP connection (Ch.82, Ch.71)"
	case strings.HasPrefix(url, "http://"):
		return "Open a plain HTTP connection (Ch.71)"
	case strings.HasPrefix(url, "mailto:"):
		return "Hand off to the default email application -- no HTTP involved"
	case strings.HasPrefix(url, "ftp://"):
		return "Speak FTP, an entirely different protocol from HTTP"
	default:
		return "Unrecognized scheme -- browser cannot decide what to do"
	}
}

func main() {
	examples := []string{
		"https://www.example.com/page",
		"http://internal-tool.local/dashboard",
		"mailto:someone@example.com",
		"ftp://files.example.com/report.pdf",
		"gopher://old.example.com",
	}
	for _, u := range examples {
		fmt.Printf("%-40s -> %s\n", u, dispatchByScheme(u))
	}
}
```

```
https://www.example.com/page            -> Open a TLS-encrypted HTTP connection (Ch.82, Ch.71)
http://internal-tool.local/dashboard    -> Open a plain HTTP connection (Ch.71)
mailto:someone@example.com              -> Hand off to the default email application -- no HTTP involved
ftp://files.example.com/report.pdf      -> Speak FTP, an entirely different protocol from HTTP
gopher://old.example.com                -> Unrecognized scheme -- browser cannot decide what to do
```

This is worth connecting directly back to Section 5: `mailto:` and `ftp://` links prove, in a single browser address bar, that "the Web" (HTTP-based) and "things a browser can open" are not even the same category — a browser is a general-purpose *client* for several different protocols, of which HTTP/the Web is only the most common. The last row is a real, deliberate demonstration of Chapter 1's shared-code failure mode: a scheme neither the sender nor the browser has agreed upon in advance simply cannot be acted on, no matter how well-formed the rest of the URL is.

### Intuitive explanation, and where it breaks

The Internet is like the national postal road system, and the Web is like one specific delivery company (say, a magazine subscription service) that uses those roads. The roads don't know or care what's being delivered — magazines, furniture, letters — they just move vehicles from address to address. If the magazine company goes out of business, the roads are completely unaffected, and a different company (furniture delivery) keeps using them just fine. The analogy holds up well for this chapter's core point, but it breaks in one place worth flagging: unlike a delivery company, the Web isn't merely *one of many* things that happened to use existing roads — it became so dominant so quickly that, as Section 2 discussed, it started getting confused with the road system itself in casual language. No delivery company has ever caused people to say "the roads are down" when they meant "my magazines haven't arrived."

---

## 5. Other Applications That Use the Internet But Aren't the Web

If the Web is only one application among many that can run on top of the Internet, what are some others? This list is worth sitting with, because every item on it is proof that "the internet" and "the Web" cannot be the same thing — these all keep working when a browser can't load a single page, and all of them stop working if the underlying Internet connectivity genuinely fails, regardless of the Web's status.

- **Email** (using protocols like SMTP and IMAP, not HTTP) — predates the Web, and remains architecturally separate from it, even though many people now check email *through* a web page, which is itself just the Web being used as a convenient front-end to a non-Web protocol underneath.
- **Instant messaging apps** (WhatsApp, Signal, iMessage) — these use the Internet's connectivity but generally define their own application-level protocols, not HTTP/the Web, for the actual message delivery (though some use HTTP for parts of their infrastructure, like media uploads).
- **Video/voice calls** (Zoom, FaceTime, WhatsApp calls) — typically use specialized real-time protocols (built largely on UDP, Chapter 58) tuned for low latency, distinct from how a browser loads a page.
- **Online multiplayer games** — often use custom, highly optimized protocols over UDP, prioritizing speed over the guarantees HTTP and the Web rely on.
- **File transfer and remote access tools** (SSH, FTP, and their modern equivalents) — move data across the Internet using their own protocols, older in some cases than the Web itself.

Every one of these is a real, everyday demonstration of Section 3's claim that the Internet doesn't care what runs on top of it. When your video call is crystal clear but a specific website won't load, you have just directly experienced the Internet working and the Web (at least, one particular part of it) failing — proof, in your own daily life, that they are not the same system.

---

## 6. What an Intranet Actually Is

An **intranet** is a private network that uses the exact same underlying technologies as the public Internet — the same addressing schemes, the same protocols, often literally the same Web technology (HTTP, browsers, web pages) — but is restricted to a specific organization and not reachable by the general public.

The key insight, and the reason this chapter groups intranets with "Internet vs. Web" rather than treating it as an unrelated third topic, is this: **an intranet is not a different kind of network technology — it's the same technology, deployed with restricted access.** A company's internal HR portal, reachable only from inside the company's office network (or via a VPN, Chapter 85), is, technically, a website — built with HTTP, HTML, and a web server, exactly like any public website. What makes it an "intranet" resource rather than part of "the Web" in the public sense is purely a matter of **access control**: who is allowed to reach it, not what technology it's built from.

```
   PUBLIC INTERNET                          COMPANY INTRANET
   (anyone, anywhere,                       (only devices inside the
    can reach these                          company's network, or
    servers)                                 connected via VPN, can
                                              reach these servers)

  [google.com]                              [internal-wiki.company.local]
  [wikipedia.org]                           [hr-portal.company.local]
  [your bank's public site]                 [internal-only build server]

  Same protocols (HTTP, TCP/IP) used on both sides of this line.
  The difference is who is allowed to reach what — a matter of
  network access and addressing scope (Chapter 40's private
  address ranges are one real mechanism for enforcing exactly
  this kind of boundary), not a different kind of technology.
```

This is why the term **extranet** also exists in professional settings, worth mentioning briefly for completeness: an extranet is a limited, controlled extension of an intranet to specific external parties (e.g., letting a trusted supplier access a specific internal ordering system), again using identical underlying technology with different access rules layered on top. Section 9 returns to this with a fuller, more production-grounded picture.

---

## 7. Why Layering Is What Makes These Separable

Everything in this chapter rests on one architectural idea this course has been quietly using since Chapter 1 and will finally name and formalize starting in Chapter 24: **networking is built in layers, where each layer only needs to know how to talk to the layers directly above and below it, not how every other layer works.**

This is *why* the Web could be invented in 1989–1991 and deployed onto an Internet that had already existed, unchanged in its core addressing and routing design, for two decades. The people who invented HTTP and HTML didn't need to redesign how computers find each other or route data across the world — that problem was already solved, one layer down, and the Web simply had to speak the existing lower layers' language to make use of it. It's also why an intranet can reuse literally the same protocol stack as the public Internet, just deployed with different access boundaries — layering doesn't care whether the network beneath an application is public or private, only that the interfaces between layers are respected.

You don't need the full seven-layer or four-layer models yet (Chapters 25 and 26 build those properly) — the idea to hold onto for now is simply: **the reason "Internet," "Web," and "intranet" can be cleanly separated as concepts is that networking was deliberately built so that the layer providing connectivity (the Internet) and the layer providing an application experience (the Web, email, a game, an intranet's internal tools) don't have to be redesigned together.** Chapter 24 explains, from first principles, why that separation was such a deliberate and important engineering choice, rather than an accident.

---

## 8. Deep Dive: The Same URL, Different Reachability

Layering (Section 7) explains *why* separation is possible; this section shows what that separation looks like from a single machine's point of view, since it's easy to nod along with the theory without seeing it play out concretely.

Imagine an employee's laptop with two network states: connected only to public Wi-Fi at a coffee shop, versus connected to that same Wi-Fi *plus* a company VPN (Chapter 85 covers exactly how a VPN accomplishes this; for now, treat it as a tunnel that makes the laptop behave as if it's plugged directly into the office network, regardless of physical location).

```
STATE 1: Coffee shop Wi-Fi only              STATE 2: Coffee shop Wi-Fi + VPN active

[laptop] --> public Internet --> reaches:    [laptop] --> public Internet --> VPN tunnel --> reaches:
   - www.wikipedia.org        (YES)             - www.wikipedia.org           (YES, same as before)
   - www.google.com           (YES)             - www.google.com              (YES, same as before)
   - hr-portal.company.local  (NO — not          - hr-portal.company.local     (YES — now reachable,
                                routable            because the VPN makes the
                                outside the         laptop appear to be inside
                                company                the company's network)
                                network at all)
```

Notice precisely what changed between the two states, and what didn't: the laptop's hardware, its browser, and the public websites' technology are all identical in both columns. What changed is purely the laptop's *network reachability* — in State 1, `hr-portal.company.local` isn't just "blocked," it's **not addressable at all** from a coffee shop's network, because that hostname resolves to a private address (Chapter 40) that only makes sense inside the company's own network. In State 2, the VPN tunnel effectively extends the company's private network to include the laptop, wherever it physically is, making the same address reachable. This is Section 6's "same technology, different access" claim, made completely concrete: nothing about HTTP, HTML, or the web server changed between the two states — only the laptop's position, logically, within which network it currently belongs to.

---

## 9. Production Notes: VPNs, Extranets, and the Blurring Boundary

The clean "public Internet vs. private intranet" picture from Section 6 is the right first mental model, but real organizations complicate it in ways worth knowing about honestly, since you'll encounter all three of these in real jobs:

- **VPN-based remote access** (Chapter 85) is now the dominant way employees reach intranet resources from outside the office, exactly as shown in Section 8. This means a modern "intranet" is often not confined to one physical LAN at all — it's a *logical* boundary (who has a valid VPN credential) rather than a purely physical one (who is plugged into the office network).
- **Extranets** extend limited intranet access to specific external organizations — for example, a manufacturer might give a parts supplier narrow access to one internal inventory system, without giving that supplier access to the rest of the company's intranet. This is neither purely "public Internet" nor purely "private intranet" — it's a deliberately controlled middle ground, using the same underlying access-control mechanisms (Chapter 84's firewalls, Chapter 85's VPNs) to grant exactly the access needed and nothing more.
- **Zero Trust architecture** is a more recent, increasingly common approach that partially dissolves the classic "inside the network = trusted, outside = untrusted" boundary this chapter has been describing. Instead of assuming anything connected to the office LAN or VPN is automatically trustworthy, a Zero Trust design checks every single request's identity and authorization individually, regardless of whether it originated "inside" or "outside" the traditional intranet perimeter. This doesn't eliminate the intranet/Internet distinction this chapter teaches — a company's internal-only resources still exist — but it changes *how* access to them is enforced, moving the enforcement point closer to each individual resource rather than relying solely on network-level boundaries. It's a real, current industry trend directly built on top of this chapter's core distinction between technology (which stays the same) and access control (which is where all the real engineering complexity actually lives).

The honest takeaway: this chapter's clean three-way split (Internet / Web / intranet) remains the correct starting mental model, but real production networks increasingly implement the "who can reach what" boundary through software-defined, identity-based rules rather than simple physical network location — a theme that resurfaces throughout Volume 12 (Network Security) and Volume 16 (Advanced Networking).

---

## 10. A Worked Scenario: Diagnosing "Internet Down" vs. "Website Down"

Let's make Section 1's opening claim concrete with a realistic troubleshooting scenario, using tools this course will properly teach in Chapters 54 and 56, but that are worth seeing in action now.

**Symptom reported:** "The internet is down! I can't load `example-shop.com`."

**Step 1 — is connectivity to the wider Internet actually working at all?**

```
$ ping 1.1.1.1
PING 1.1.1.1 (1.1.1.1): 56 data bytes
64 bytes from 1.1.1.1: icmp_seq=0 ttl=57 time=14.221 ms
64 bytes from 1.1.1.1: icmp_seq=1 ttl=57 time=13.998 ms
64 bytes from 1.1.1.1: icmp_seq=2 ttl=57 time=14.502 ms
```

This succeeds — a well-known, generally very reliable public server responded normally. **This already tells us the Internet, in this chapter's sense, is working**: your device has a working address, a working route to a distant server, and the underlying infrastructure is functioning.

**Step 2 — is DNS (name lookup, covered fully in Volume 10) resolving the specific site's name?**

```
$ ping example-shop.com
ping: cannot resolve example-shop.com: Unknown host
```

This fails. Combined with Step 1's success, this strongly suggests the problem is specific to `example-shop.com` — perhaps its DNS records are misconfigured, or its hosting provider is having an outage — and has nothing to do with your own Internet connectivity, which Step 1 already proved was fine.

**Step 3 — confirm by trying a different, known-good website:**

```
$ curl -I https://www.wikipedia.org
HTTP/2 200
```

Success. This confirms the diagnosis: **the Internet is fine. The Web, broadly, is fine. One specific website has a problem** — very likely on the server side, or in that specific domain's DNS configuration, not anywhere in your own network or in "the internet" as a whole.

This three-step process — check raw connectivity, check name resolution, check a known-good alternative — is a genuinely real, professionally used diagnostic pattern, and it only makes sense once you've internalized this chapter's core distinction: connectivity (the Internet) and a specific application's availability (one website, part of the Web) are different things that fail independently and need to be tested separately.

### A second scenario: "The intranet is down," diagnosed differently

Compare the above to a superficially similar complaint with a genuinely different cause: "I can't reach `hr-portal.company.local`, but Google works fine."

```
$ ping 1.1.1.1
64 bytes from 1.1.1.1: icmp_seq=0 ttl=57 time=14.221 ms      <- public Internet: fine

$ ping hr-portal.company.local
ping: cannot resolve hr-portal.company.local: Unknown host   <- fails, but for a
                                                                  completely different
                                                                  reason than Step 2
                                                                  in the scenario above
```

The symptom looks identical to Step 2 of the earlier scenario (a name fails to resolve), but the correct next question is entirely different, because of Section 6 and Section 8's access-control distinction: `hr-portal.company.local` is very likely a **private, non-publicly-routable hostname** that only resolves at all when the requester is inside the company's network or connected via VPN (Section 8) — it was never going to resolve from a coffee shop's public Wi-Fi, no matter how healthy that public Internet connection is. The fix here isn't "check if the destination server is down" (Volume 10's DNS troubleshooting) — it's "check whether you're currently inside the intranet's access boundary at all" (confirm VPN status, per Chapter 85). Diagnosing this correctly requires recognizing which of the two very different failure categories — a public Internet/Web problem, or an intranet access-control problem — you're actually looking at, which is exactly the distinction this entire chapter has been building toward.

---

## 11. Hands-On Experiment: Prove the Separation Yourself

You can reproduce Section 10's reasoning on your own device right now, without needing anything to actually be broken.

1. Open a terminal (Terminal on macOS/Linux, Command Prompt or PowerShell on Windows) and run `ping 1.1.1.1` (or `ping 8.8.8.8`). Confirm you get replies — this demonstrates raw Internet connectivity, independent of any website.
2. Open a video call app, or make a regular phone call over Wi-Fi calling if your phone supports it, and confirm it works. This is a non-Web application using the same underlying Internet connectivity you just tested.
3. Now open a browser and load any website. Notice: you've now tested three completely different things (raw connectivity, a non-Web application, and the Web) that all depend on the same Internet infrastructure but are otherwise entirely independent of each other.
4. If your workplace or university has an internal-only site (a wiki, an HR portal, a VPN-gated tool), try loading it while connected to that network, then try loading the exact same address after disconnecting (e.g., turning off Wi-Fi and using cellular data, or disconnecting a VPN). If it stops loading, you've just directly observed Section 8's exact experiment for yourself — the technology didn't change, only your permission and reachability did.

---

## 12. Common Misconceptions

- **"If a website won't load, the internet is down."** As Section 10 demonstrates step by step, this is very often false, and confirming or ruling it out takes under a minute using tools this course fully explains in Chapters 54 and 56.
- **"The Web and HTTP are the same thing as 'going online.'"** As Section 5 lists concretely, an enormous amount of everyday "online" activity — messaging, calling, gaming, email — does not use the Web or HTTP at all, even though all of it depends on the same underlying Internet.
- **"An intranet is a completely different, more primitive kind of network."** As Section 6 makes explicit, an intranet typically uses identical technology to the public Internet and Web — the difference is exclusively about who is allowed to access it, not what protocols or software are involved.
- **"The Web was part of the Internet's original design."** As Section 4 states plainly, the Web was invented over twenty years after the Internet's predecessor network (ARPANET) began operating, specifically as an application layered on top of an Internet that already worked without it. Volume 2 of this course (Chapters 7–13) tells this history in the correct chronological order.
- **"An intranet is always confined to one physical building or location."** As Section 9 shows, modern VPN-based remote access and Zero Trust architectures mean an "intranet," in practice, is increasingly a logical boundary based on identity and authorization rather than a purely physical network location — an employee working from home, connected via VPN, is logically "inside" the intranet even though they're nowhere near the office building.

---

## 13. Connections Backward and Forward

This chapter used Chapter 3's precise definition of "network" and Chapter 4's LAN/WAN vocabulary to show that the Internet is infrastructure, the Web is one application running on that infrastructure, and an intranet is the same technology deployed with restricted access — and it previewed, without yet formalizing, the layering principle (fully developed starting Chapter 24) that is *why* these three things can be described, built, and diagnosed independently of each other. Chapter 6 now builds the full picture this chapter has been assuming piece by piece: exactly how millions of separately-owned LANs, connected through countless independently-operated WAN links, add up to the single (uncoordinated, unowned) structure everyone calls "the Internet" — the first complete mental model this course will hand you before Volume 2 tells the real history of how it came to exist.

---

## 14. Interview Questions & Model Answers

**Q1 (Beginner): What is the fundamental difference between "the Internet" and "the Web"?**

*Model answer:* The Internet is the underlying global network infrastructure — physical links, addressing, and routing — that lets any connected computer reach any other. The Web is one specific application that runs on top of that infrastructure, consisting of hypertext documents (web pages), a request/response protocol (HTTP), and a universal addressing scheme for documents (URLs). The Internet existed and functioned for over two decades before the Web was invented in 1989–1991, and many non-Web applications (email, messaging, video calls, gaming) use the Internet today without using the Web at all.

**Q2 (Beginner): Give two examples of applications that use the Internet but are not part of the Web.**

*Model answer:* Email (using protocols like SMTP and IMAP rather than HTTP) and real-time voice/video calling apps (which typically use specialized low-latency protocols, often built on UDP, rather than HTTP) are both common examples. Instant messaging apps and online multiplayer games are two more.

**Q3 (Intermediate): What is an intranet, and how does it differ technologically from the public Internet and Web?**

*Model answer:* An intranet is a private network, typically belonging to a single organization, that generally uses the exact same underlying protocols and technologies as the public Internet and Web — TCP/IP addressing, HTTP, browsers, web servers. It differs not in the technology used, but in access control: an intranet's resources are restricted to a specific set of users or devices (e.g., only computers inside a company's office network, or connected via VPN), whereas public Internet/Web resources are reachable by anyone.

**Q4 (Intermediate): A user reports "the internet is down," but a video call works fine while a specific website fails to load. What does this tell you, and what would you check next?**

*Model answer:* This strongly suggests the user's actual Internet connectivity is working — a functioning video call demonstrates a working connection capable of carrying real-time data across the network. The problem is more likely isolated to the specific website: possibly a DNS resolution issue, a problem with that site's server or hosting provider, or an application-level issue rather than a network connectivity issue. The next steps would be to test raw connectivity to a known-reliable address (e.g., pinging a public DNS resolver like 1.1.1.1), check whether the specific domain name resolves at all, and try loading a different, known-good website to further isolate whether the issue is with connectivity in general or with just the one site.

**Q5 (Advanced): Explain, using the concept of layering (previewed in this chapter, formalized in Chapter 24), why the Web could be invented and deployed globally without requiring any changes to the Internet's core addressing and routing infrastructure.**

*Model answer:* Networking is architected in layers, where each layer exposes a defined interface to the layers above and below it without requiring those other layers to understand its own internal workings. The Internet's addressing and routing layers were already solving "how does one computer's data reach another computer's address" by the time the Web was invented. HTTP and the Web were built as an application that simply uses the connectivity those lower layers already provide — a web browser doesn't need to know how routing works, and a router doesn't need to know or care that the data it's forwarding happens to be a web page rather than an email or a game's data. This separation of concerns is precisely what allowed Tim Berners-Lee to design and deploy the Web on top of an existing, unmodified Internet, and is the same architectural property that later allowed entirely new applications (streaming video, real-time gaming, cloud computing) to be layered on top of the same underlying Internet without redesigning its core.

**Q6 (Advanced): How does a VPN change which resources a laptop can reach, without changing anything about the resources themselves? Relate this to the difference between "technology" and "access control" this chapter emphasizes.**

*Model answer:* A VPN creates an encrypted tunnel that makes a remote device appear, from a networking perspective, to be located inside a private network (like a company's intranet), regardless of its actual physical location or the public network it's using to reach the internet. This changes the device's reachability — private, non-publicly-routable addresses (Chapter 40) inside the company's network become reachable through the tunnel — without changing anything about the internal resources themselves: the same HTTP servers, the same protocols, the same web pages. This is a direct, practical demonstration of this chapter's core claim that intranets differ from the public Internet/Web purely in access control, not underlying technology; a VPN is simply a mechanism for granting that access control conditionally, based on identity and authentication, rather than physical network location.

**Q7 (Intermediate): What does a URL scheme like `https://`, `mailto:`, or `ftp://` actually tell a computer, and why is this a real instance of Chapter 1's shared-code idea?**

*Model answer:* The URL scheme is an explicit, in-band declaration of which protocol the requester wants to use, allowing a browser or operating system to dispatch the rest of the URL to the correct handler before any network connection is made — `https://` triggers a TLS-wrapped HTTP request, `mailto:` hands off to an email client entirely outside the Web, and `ftp://` invokes a different file-transfer protocol altogether. This is a direct instance of Chapter 1's shared-code idea because the scheme only works because every browser and operating system has agreed in advance on what each prefix means; an unrecognized scheme cannot be acted on at all, exactly matching Chapter 1's point that a signal (here, a URL string) only carries meaning when sender and receiver share a prior agreement about its interpretation.

---

## 15. Exercises

### Easy

1. In your own words, explain why "the internet is down" and "the website is down" can be, and usually are, two different problems.
2. List three applications you use regularly that use the Internet but are not part of the Web (i.e., don't primarily work through a browser using HTTP).
3. Explain the difference between an intranet and a public website, focusing on what is different and what is identical between them.
4. Using Section 4's URL scheme table, explain what would happen if you clicked a `mailto:` link on a device with no email application configured at all, and relate your answer to Chapter 1's shared-code failure modes.

### Medium

1. Walk through Section 10's three-step diagnostic scenario using a real website and your own device (pick a working site, so nothing is actually broken, and simply confirm each step succeeds). Write down the actual output you observe at each step.
2. A friend says, "our company intranet uses completely different technology from the public internet, that's why outsiders can't access it." Correct this misconception in 3–4 sentences, using this chapter's explanation of access control versus technology.
3. Explain why the invention of the Web in 1989–1991, more than 20 years after ARPANET began operating, is strong historical evidence that the Internet and the Web are architecturally separate things, rather than the same system that simply gained a new name over time.
4. Using Section 8's two-state example, explain in your own words why `hr-portal.company.local` is unreachable in State 1 but reachable in State 2, without anything about the HR portal server itself changing.

### Hard

1. Research (briefly) what an "extranet" is, and explain, using this chapter's framework of "same technology, different access rules," how an extranet relates to both an intranet and the public Internet/Web.
2. Consider a company that runs its internal HR portal on the same physical servers and using the same HTTP protocol as its public marketing website, differing only in network-level access restrictions (the HR portal is only reachable from inside the company's network). Using vocabulary from this chapter and Chapter 4, describe what would need to be true about the company's network for this access restriction to actually work (you're not expected to know the specific mechanism yet — VPNs are covered in Chapter 85, and private addressing in Chapter 40 — reason about what property the restriction needs to have).
3. This chapter claims the Internet "doesn't particularly care what runs on top of it," then Section 9 introduces Zero Trust architecture as a shift away from purely network-location-based trust. Construct an argument for why Zero Trust represents a genuine evolution in how the intranet/Internet boundary is enforced, rather than simply a renaming of VPN-based access control — identify at least one concrete difference in approach.
4. Using Section 4's URL scheme table, propose what a hypothetical new scheme (e.g., `myapp://`) would need every relevant browser and operating system to agree on in advance before it could be used reliably, and relate your answer to Chapter 1's shared-code problem.

---

## 16. Summary

| Term | Meaning |
|---|---|
| The Internet | The global network infrastructure (addressing, routing, physical links) connecting computers worldwide; a WAN, per Chapter 4 |
| The Web (World Wide Web) | One application running on top of the Internet: hypertext documents (pages), HTTP (request/response protocol), and URLs (addressing for documents) |
| HTTP | The protocol browsers and servers use to request and deliver web pages (fully covered starting Chapter 70) |
| Intranet | A private network using the same underlying technology as the public Internet/Web, restricted to a specific organization's users |
| Extranet | A limited, controlled extension of an intranet to specific external parties |
| VPN (preview) | A tunnel that makes a remote device behave as if it's located inside a private network, regardless of physical location |
| Zero Trust (preview) | A security approach that verifies every request individually by identity, rather than trusting anything simply because it's "inside" the network perimeter |
| Layering (preview) | The architectural principle that lets an application (like the Web) be built on top of existing network infrastructure without modifying that infrastructure |
| "Internet down" vs. "website down" | Two genuinely different failure classes, diagnosable independently by testing raw connectivity, name resolution, and a known-good alternative site |

This chapter separated the Internet (infrastructure) from the Web (one application on top of it) from an intranet (the same technology, restricted access) — three ideas people conflate constantly, now precisely distinguished. Chapter 6 builds the complete mental model of how the Internet's infrastructure is actually structured: home networks, ISPs, and the fact that no single organization owns or controls the whole thing.
