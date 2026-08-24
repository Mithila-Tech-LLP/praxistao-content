# Chapter 70: URLs and the Structure of a Web Request

> **"A URL looks like a string you type into a box. It is actually a complete set of instructions — which protocol to speak, who to ask for an address, which door to knock on, and which room to walk into once you're inside."**

---

## Table of Contents

1. [The Problem: Addressing One Resource Among Billions](#1-the-problem-addressing-one-resource-among-billions)
2. [A Naive Address Scheme, and Why It Fails](#2-a-naive-address-scheme-and-why-it-fails)
3. [URI, URL, and URN — Getting the Terms Straight](#3-uri-url-and-urn--getting-the-terms-straight)
4. [Anatomy of a URL, Piece by Piece](#4-anatomy-of-a-url-piece-by-piece)
5. [A Fully Annotated Real Example](#5-a-fully-annotated-real-example)
6. [Percent-Encoding — What Happens to "Illegal" Characters](#6-percent-encoding--what-happens-to-illegal-characters)
7. [Default Ports and the Scheme-Port Relationship](#7-default-ports-and-the-scheme-port-relationship)
8. [What Happens Before the First HTTP Byte Is Sent](#8-what-happens-before-the-first-http-byte-is-sent)
9. [A Sequence Diagram of the Full Pre-HTTP Journey](#9-a-sequence-diagram-of-the-full-pre-http-journey)
10. [Deep Dive: How a Browser Actually Parses What You Typed](#10-deep-dive-how-a-browser-actually-parses-what-you-typed)
11. [Relative URLs and Base Resolution](#11-relative-urls-and-base-resolution)
12. [Hands-On: Watching the Pre-Request Work Happen](#12-hands-on-watching-the-pre-request-work-happen)
13. [Code: Parsing a URL in Go](#13-code-parsing-a-url-in-go)
14. [Common Misconceptions](#14-common-misconceptions)
15. [Production Notes](#15-production-notes)
16. [Interview Questions & Model Answers](#16-interview-questions--model-answers)
17. [Exercises](#17-exercises)
18. [Summary](#18-summary)

---

## 1. The Problem: Addressing One Resource Among Billions

By Chapter 69 you have everything needed to get one machine talking to another: DNS (Chapters 66–69) turns a name into an IP address, and TCP (Chapters 57–65) turns that IP address plus a port into a reliable byte stream between two programs. But neither of those chapters answers a much more specific question:

*Which* byte stream, to *which* program, asking for *which* thing?

`google.com` resolves to an IP address. But Google's web server behind that IP address hosts search results, image search, maps tiles, account settings pages, and a thousand other things. A single IP address and a single TCP port (443, almost always) is the front door to an entire building with millions of rooms. Something has to specify:

- Which protocol's rules to speak once the door opens (HTTP? something else?)
- Which specific resource inside the building you want (a search results page? a JPEG? a JSON API endpoint?)
- Extra parameters that modify the request (search for "laptops", sorted by price)
- A location *within* that resource you want to jump straight to (a specific section of a long page)

The URL (Uniform Resource Locator) is the answer: one compact string format that identifies a protocol, a destination, a specific resource at that destination, and optional parameters and sub-location — all in one line you can type, click, bookmark, or paste into a chat message.

## 2. A Naive Address Scheme, and Why It Fails

Imagine designing this from scratch. Your first instinct might be: "just use the IP address and a file path."

```
142.250.80.46/search.html
```

This fails almost immediately, for reasons that should feel familiar from earlier chapters:

- **IP addresses change.** A server might move to new hardware with a new address (this is exactly the problem Chapter 66 introduced DNS to solve). Baking a raw IP into every link on the Web would break every link every time infrastructure changed.
- **One IP address can host many unrelated websites.** Modern web servers use *virtual hosting* — one machine, one IP, hundreds of independent domains, distinguished only by a name the client sends along with the request (the `Host` header, which Chapter 71 covers). An IP address alone doesn't tell the server which site you mean.
- **There's no room for a protocol choice.** Should the client speak HTTP? HTTPS? FTP? Something else entirely? A bare address plus path says nothing about that.
- **There's no room for parameters.** "Search for laptops, sorted by price" is not a static file — it needs a way to pass structured input alongside the target resource.

So the real design needs, at minimum: a way to say *which protocol*, a way to say *which name* (not a raw IP), a way to say *which resource at that name*, a way to pass *extra parameters*, and — it turns out — a way to say *which part of the resource*, for jumping straight to content without a new server round trip. That is exactly what a URL provides, and RFC 3986 (2005) is the standard that formalizes its grammar.

## 3. URI, URL, and URN — Getting the Terms Straight

Before going further, three terms that get used loosely in casual conversation but have a precise relationship in the spec:

```
URI (Uniform Resource Identifier)
  │  The umbrella term: any string that identifies a resource.
  │
  ├── URL (Uniform Resource Locator)
  │     A URI that also tells you HOW to get the resource — a location
  │     plus an access method. Almost everything you call "a URL" in
  │     daily life is technically a URL, which is a kind of URI.
  │     Example: https://example.com/page.html
  │
  └── URN (Uniform Resource Name)
        A URI that names a resource without saying how to fetch it —
        pure identity, no location. Rare in everyday browsing, common
        in specs and standards.
        Example: urn:isbn:9780132350884
```

Every URL is a URI. Not every URI is a URL. In this course, and in almost all networking and web-engineering conversation, "URL" is the term that matters — it's what a browser's address bar holds and what an HTTP request is built from.

## 4. Anatomy of a URL, Piece by Piece

Here is the full grammar, with every optional piece shown in brackets:

```
scheme://[userinfo@]host[:port][/path][?query][#fragment]
```

Let's take each component in isolation, in engineering terms, before reassembling a real example.

**Scheme** — the protocol the client will speak once connected. `http`, `https`, `ftp`, `mailto`, `ws`, `wss`, `file`. This is the single biggest decision in the URL: it determines default port, whether the connection will be encrypted (Chapter 82 previews TLS), and which "verb vocabulary" applies once connected.

**Authority** — everything between `//` and the next `/`, `?`, or `#`. It has three sub-parts:
- **Userinfo** (`user:password@`) — credentials embedded directly in the URL. Rare and actively discouraged today (Section 13 explains why), but part of the grammar since the beginning.
- **Host** — a domain name (resolved via DNS, Chapters 66–68) or a raw IP address (IPv4 or bracketed IPv6, e.g. `[2001:db8::1]`).
- **Port** — the TCP port to connect to. Optional; if omitted, the scheme's default port applies (Section 7).

**Path** — identifies a specific resource *within* the host, hierarchically, using `/` as a separator — historically mirroring a filesystem path, though on a modern web server it is usually a purely logical route handled by application code, not a real file.

**Query** — a `?`-prefixed string of key=value pairs joined by `&`, used to pass parameters that don't belong in the hierarchical path itself (search terms, pagination, filters, sort order).

**Fragment** — a `#`-prefixed identifier for a location *within* the resource, resolved entirely on the client side. This is the one piece of the URL, by design, that **never leaves the browser** — it is not sent to the server as part of the HTTP request at all (Section 8 explains why that matters).

## 5. A Fully Annotated Real Example

```
https://alice:s3cr3t@shop.example.com:8443/products/search?q=laptops&sort=price#reviews
└─┬──┘   └─┬─┘ └─┬──┘ └──────┬────────┘ └┬─┘ └────────┬─────────┘ └───────┬────────┘ └──┬───┘
  │        │     │           │           │            │                   │              │
scheme  userinfo userinfo   host        port         path                query        fragment
```

Field by field:

| Component | Value | Meaning |
|---|---|---|
| Scheme | `https` | Speak HTTP wrapped in TLS. Default port 443. Connection will be encrypted (Ch. 82). |
| Userinfo | `alice:s3cr3t` | Username `alice`, password `s3cr3t`, sent as HTTP Basic-style credentials embedded in the URL. |
| Host | `shop.example.com` | The domain name. Chapters 66–68's DNS machinery resolves this to an IP address before anything else can happen. |
| Port | `8443` | Connect to TCP port 8443, *not* the HTTPS default of 443 — the site is running on a non-standard port. |
| Path | `/products/search` | The specific resource being requested — logically, "the product search endpoint." |
| Query | `q=laptops&sort=price` | Two parameters: `q` (search term) = `laptops`, `sort` = `price`. |
| Fragment | `reviews` | After the page loads, the browser should scroll to the element with `id="reviews"`. Never sent over the network. |

Reading this URL end to end in plain English: *"Using HTTPS, authenticate as alice, connect to shop.example.com on port 8443, ask for the product-search resource with a search term of 'laptops' sorted by price, and once the page has loaded, jump to the reviews section."*

## 6. Percent-Encoding — What Happens to "Illegal" Characters

URLs are restricted to a subset of ASCII. Anything outside that subset — a literal space, a `&` inside a query value that isn't meant as a separator, a non-Latin character, a `#` that's meant to be data and not a fragment marker — must be **percent-encoded**: replaced with `%` followed by the two-digit hexadecimal value of the byte.

```
Character     Percent-encoded
' ' (space)   %20   (or '+' inside a query string, by web convention, not the URL spec itself)
'&'           %26
'#'           %23
'/'           %2F
'?'           %3F
'@'           %40
'é'           %C3%A9   (UTF-8 encodes 'é' as two bytes, each percent-encoded)
```

RFC 3986 defines two categories of characters:

```
Unreserved (always safe, never encoded):
  A-Z a-z 0-9 - . _ ~

Reserved (have special meaning as URL delimiters; must be encoded
  if you want them to appear as literal DATA rather than as a
  delimiter):
  : / ? # [ ] @ ! $ & ' ( ) * + , ; =
```

**Why this matters concretely:** if a search box lets a user type `Q&A` and that string is put directly into a query parameter, the unescaped `&` would be misread as separating two different parameters. The value must become `q=Q%26A` for the server to correctly reconstruct the literal string `Q&A`.

**Internationalized domain names (IDN)** face a related but distinct problem — domain names historically allowed only ASCII, but people register domains in Arabic, Chinese, Cyrillic, and other scripts. The fix is **Punycode**: `münchen.de` is transmitted over the wire as `xn--mnchen-3ya.de`. Browsers decode this back to the readable form for display (and flag suspicious lookalike domains — this is one of the defenses against homograph phishing attacks, previewed in Chapter 83).

## 7. Default Ports and the Scheme-Port Relationship

Chapter 57 established that ports identify *which program* on a machine a connection is for. Schemes have well-known default ports so that the vast majority of URLs never need to state one explicitly:

| Scheme | Default port | Notes |
|---|---|---|
| `http` | 80 | Plaintext |
| `https` | 443 | TLS-wrapped (Ch. 82) |
| `ftp` | 21 | File Transfer Protocol |
| `ws` | 80 | WebSocket (Ch. 76) |
| `wss` | 443 | WebSocket over TLS |
| `ssh` | 22 | Not typically a browsable URL scheme, but same idea |

If a URL omits the port, the client fills in the scheme's default. `https://example.com/` and `https://example.com:443/` are exactly the same request. Writing `:8443` (as in Section 5's example) overrides the default — common for development servers, internal tools, or services deliberately avoiding the standard port.

## 8. What Happens Before the First HTTP Byte Is Sent

This is the payoff of reading Chapters 57–69 first: by the time an HTTP request is actually transmitted, an enormous amount of work from *earlier* volumes has already happened, entirely invisibly. Walking through it in order:

**Step 1 — Parse the URL.** The browser breaks the typed (or clicked) string into scheme, host, port, path, query, and fragment, exactly as Section 4 described. This is pure string processing; nothing has touched the network yet.

**Step 2 — Check local caches, first.** Before doing anything else, the browser checks whether it already has a fresh, cached copy of this exact resource (Chapter 72 covers HTTP caching in depth). If there's a valid cache hit, the entire rest of this list can be skipped — that's the whole point of caching.

**Step 3 — Resolve the host to an IP address (DNS, Chapters 66–68).** The browser checks its own DNS cache, then the OS resolver's cache, then — on a miss — sends a query to a configured recursive resolver, which may itself walk the full root → TLD → authoritative hierarchy (Chapter 67) before returning an answer. This step ends with an IP address (or several, for load-balanced or geo-distributed services) and nothing else — no HTTP has happened.

**Step 4 — Open a TCP connection (three-way handshake, Chapter 59).** The client sends a SYN to the resolved IP address on the URL's port (Section 7). The server responds SYN-ACK, the client responds ACK. This costs one full round-trip time (RTT) before a single application byte can move in either direction, and it establishes the reliable, ordered byte stream that Chapters 60–64 describe TCP maintaining.

**Step 5 — (If `https`) Perform the TLS handshake.** Chapter 82 covers this fully; the short version previewed here: the client and server negotiate a cipher suite, the server proves its identity with a certificate (Chapter 81's PKI machinery), and both sides derive a shared symmetric key (Chapter 78) to encrypt everything that follows. In TLS 1.3 this costs one additional RTT on a fresh connection (zero additional RTTs on a resumed one) — meaning a first-ever HTTPS connection typically needs **2 full round trips** before any HTTP request leaves the client.

**Step 6 — Only now does HTTP begin.** The client writes an HTTP request (Chapter 71) into the now-open, now-encrypted TCP stream. Everything above this line was infrastructure; HTTP is the first layer in this entire sequence that has ever heard of the word "URL" as a whole, structured concept — TCP and DNS only ever saw the pieces (an address, a port) that were extracted from it back in Step 1.

**The fragment never appears anywhere above.** Because the fragment (Section 4) is resolved entirely client-side after the page has already loaded, it plays no role in DNS resolution, TCP connection, or the HTTP request line itself — a server serving `/products/search?q=laptops#reviews` receives a request for `/products/search?q=laptops` and has no way of knowing a fragment was even present, unless the client-side JavaScript reads it from `window.location.hash` and acts on it separately.

**Step 2, revisited — the cache hierarchy is deeper than "one cache."** By the time Chapter 72 covers HTTP caching properly, it's worth knowing here that "check local caches" in Step 2 is really several independent checks, roughly in this order, each one a chance to avoid every later step entirely:

```
1. In-memory browser cache (fastest — cleared on browser restart, in some browsers)
2. On-disk browser HTTP cache (survives restarts, keyed by full URL + Vary headers)
3. OS-level DNS resolver cache (separate from the browser's own DNS cache — Step 3)
4. A forward proxy or corporate NAT/cache, if the network is configured to use one
5. A CDN edge cache sitting in front of the origin server (Chapter 96) — invisible
   to the client, but from the origin server's point of view, this is also
   "a cache that might make Step 4 onward unnecessary"
```

A cache hit at layer 1 or 2 means the browser never even reaches Step 3 (DNS) — the whole pre-HTTP journey in this section is skipped, which is precisely why a page you've visited seconds ago reloads instantly while a page you're visiting for the first time visibly takes a moment.

## 9. A Sequence Diagram of the Full Pre-HTTP Journey

```mermaid
sequenceDiagram
    participant B as Browser
    participant C as OS/Browser Cache
    participant R as Recursive Resolver
    participant S as Server (shop.example.com)

    B->>B: Parse URL into scheme/host/port/path/query/fragment
    B->>C: Have I cached this exact resource, and is it still fresh?
    C-->>B: Cache miss (or stale) — must fetch
    B->>C: Do I already know the IP for shop.example.com?
    C-->>B: Not cached — need DNS
    B->>R: DNS query: A record for shop.example.com?
    R-->>B: 203.0.113.42 (Ch. 66-68)
    B->>S: TCP SYN to 203.0.113.42:443
    S-->>B: SYN-ACK
    B->>S: ACK  (connection established, Ch. 59)
    B->>S: TLS ClientHello
    S-->>B: ServerHello, Certificate, Finished (Ch. 82 preview)
    B->>S: TLS Finished (encrypted channel ready)
    B->>S: GET /products/search?q=laptops&sort=price HTTP/1.1  (Ch. 71 begins here)
    S-->>B: HTTP/1.1 200 OK + HTML body
    B->>B: Render page, then jump to #reviews locally
```

Notice that of the nine network-crossing steps in this diagram, only the last two are HTTP. Everything before them is Chapters 59–69's material, doing its job silently every single time you load a page.

## 10. Deep Dive: How a Browser Actually Parses What You Typed

Browsers apply a set of heuristics before even reaching Section 4's clean grammar, because what you type into an address bar is rarely a well-formed URL:

- **Missing scheme.** Typing `example.com` gets a scheme (`https`, almost universally today) silently prepended. Browsers historically tried `http` and upgraded on redirect; most now default straight to `https` and fall back to `http` only on failure, or rely on HSTS preload lists.
- **Search vs. navigation.** Typing `best laptops 2026` (with spaces, no dot, no scheme) is heuristically detected as a search query, not a URL, and is routed to the browser's configured search engine instead — `best laptops 2026` becomes something like `https://www.google.com/search?q=best+laptops+2026`.
- **Trailing dot and case normalization.** `Example.COM` is lowercased to `example.com` before DNS resolution — hostnames are case-insensitive by spec (RFC 4343), though paths are *not* automatically lowercased (`/Products` and `/products` can be genuinely different resources on a case-sensitive server).
- **Default path.** A bare `https://example.com` (no path at all) is normalized to `https://example.com/` — an empty path and a path of `/` are treated identically; the request line that goes out over the wire always has *some* path, and it defaults to `/`.

Once normalized, the string is run through the actual URL parser (implementing RFC 3986's formal grammar, or in browsers specifically, the WHATWG URL Living Standard, which is stricter and more precisely specified than the older RFC for web-compatibility reasons) to produce the structured scheme/host/port/path/query/fragment record Section 4 described.

## 11. Relative URLs and Base Resolution

Not every URL you encounter is a full, "absolute" one like the examples above. Look at the HTML for almost any real web page and you'll find links like:

```html
<a href="/about">About</a>
<a href="pricing.html">Pricing</a>
<a href="../images/logo.png">Logo</a>
<a href="?page=2">Next page</a>
<a href="#top">Back to top</a>
<a href="//cdn.example.com/lib.js">External script</a>
```

None of these are complete URLs on their own — none has a scheme, and most don't have a host. They are **relative references**, and resolving them requires a **base URL**: the URL of the document the link appears in. If this HTML was served from `https://shop.example.com/products/search`, the browser resolves each relative reference against that base according to a precise algorithm (RFC 3986 §5, "Reference Resolution"):

| Relative reference | Resolves against the base above to |
|---|---|
| `/about` | `https://shop.example.com/about` (leading `/` replaces the whole path) |
| `pricing.html` | `https://shop.example.com/products/pricing.html` (replaces only the last path segment) |
| `../images/logo.png` | `https://shop.example.com/images/logo.png` (`..` walks up one directory level) |
| `?page=2` | `https://shop.example.com/products/search?page=2` (same path, new query) |
| `#top` | `https://shop.example.com/products/search?q=laptops#top` (same everything, new fragment only) |
| `//cdn.example.com/lib.js` | `https://cdn.example.com/lib.js` (a **protocol-relative** URL — keeps the base's scheme, replaces everything from the host onward) |

This is not a cosmetic convenience — it is the mechanism that lets an entire website be moved from a staging domain to a production domain, or served over `http` in development and `https` in production, without rewriting a single internal link. Every internal `<a>`, `<img>`, `<script src>`, and `fetch()` call written as a relative reference automatically re-resolves against wherever the page itself happens to be served from.

**HTTP redirects (Chapter 71's 3xx status codes) can also be relative**, and the same base-resolution rule applies: a server responding with `Location: /login` is telling the client "resolve this against the URL you just requested," not handing over a complete URL.

## 12. Hands-On: Watching the Pre-Request Work Happen

You can observe every layer from Section 8 directly from a terminal.

**See the DNS step in isolation:**

```
$ dig +short shop.example.com
203.0.113.42
```

**See the TCP handshake and TLS handshake happen, then the HTTP request/response, with `curl`'s verbose timing:**

```
$ curl -v -o /dev/null https://example.com/
*   Trying 93.184.216.34:443...
* Connected to example.com (93.184.216.34) port 443 (#0)
* ALPN: offering h2, http/1.1
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (IN), TLS handshake, Finished (20):
* TLSv1.3 (OUT), TLS handshake, Finished (20):
* SSL connection using TLSv1.3 / TLS_AES_128_GCM_SHA256
* ALPN: server accepted h2
> GET / HTTP/2
> Host: example.com
> user-agent: curl/8.4.0
> accept: */*
>
< HTTP/2 200
< content-type: text/html
```

**Break down exactly where the time goes**, using `curl`'s timing variables:

```
$ curl -o /dev/null -s -w \
  "dns: %{time_namelookup}s  connect: %{time_connect}s  tls: %{time_appconnect}s  ttfb: %{time_starttransfer}s  total: %{time_total}s\n" \
  https://example.com/

dns: 0.012s  connect: 0.041s  tls: 0.089s  ttfb: 0.132s  total: 0.134s
```

Reading this: 12ms was spent on DNS (Section 8, Step 3), the TCP handshake finished by 41ms (Step 4), TLS finished by 89ms (Step 5), and the first byte of the HTTP response arrived at 132ms — meaning roughly 43ms was spent purely inside the HTTP request/response cycle (Chapter 71) itself, sitting on top of nearly 90ms of connection setup that had *nothing to do with HTTP*.

## 13. Code: Parsing a URL in Go

```go
package main

import (
	"fmt"
	"net/url"
)

func main() {
	raw := "https://alice:s3cr3t@shop.example.com:8443/products/search?q=laptops&sort=price#reviews"

	u, err := url.Parse(raw)
	if err != nil {
		panic(err)
	}

	fmt.Println("Scheme:  ", u.Scheme)
	fmt.Println("User:    ", u.User.Username())
	pw, _ := u.User.Password()
	fmt.Println("Password:", pw)
	fmt.Println("Host:    ", u.Hostname())
	fmt.Println("Port:    ", u.Port())
	fmt.Println("Path:    ", u.Path)
	fmt.Println("RawQuery:", u.RawQuery)
	fmt.Println("Fragment:", u.Fragment)

	q := u.Query()
	fmt.Println("q param:   ", q.Get("q"))
	fmt.Println("sort param:", q.Get("sort"))
}

// Output:
// Scheme:   https
// User:     alice
// Password: s3cr3t
// Host:     shop.example.com
// Port:     8443
// Path:     /products/search
// RawQuery: q=laptops&sort=price
// Fragment: reviews
// q param:    laptops
// sort param: price
```

Go's `net/url` package implements the RFC 3986 grammar directly, giving you the exact same field breakdown as Section 4 — programmatically, this is precisely the parsing step described as "Step 1" in Section 8.

**Building a URL safely (why you should never string-concatenate one):**

```go
package main

import (
	"fmt"
	"net/url"
)

func main() {
	base, _ := url.Parse("https://shop.example.com/products/search")

	q := base.Query()
	q.Set("q", "laptops & tablets") // contains a raw '&' — dangerous to concatenate by hand
	q.Set("sort", "price")
	base.RawQuery = q.Encode()

	fmt.Println(base.String())
}

// Output:
// https://shop.example.com/products/search?q=laptops+%26+tablets&sort=price
```

Notice `url.Values.Encode()` automatically percent-encodes the literal `&` inside the search term (Section 6) so it cannot be misread as a parameter separator, and encodes the space as `+` per the `application/x-www-form-urlencoded` convention. This is exactly why hand-building query strings with `fmt.Sprintf` or plain string concatenation is a recurring source of bugs (and, in server-side contexts, of injection vulnerabilities): a value containing `&`, `#`, or `%` silently corrupts the URL's structure unless it goes through an encoder.

**Resolving a relative reference against a base (Section 11), programmatically:**

```go
base, _ := url.Parse("https://shop.example.com/products/search")
rel, _ := url.Parse("../images/logo.png")
resolved := base.ResolveReference(rel)
fmt.Println(resolved.String())
// Output: https://shop.example.com/images/logo.png
```

This `ResolveReference` call is the exact algorithm every browser runs internally whenever it renders an `<img src="../images/logo.png">` tag found in the page at `/products/search`.

## 14. Common Misconceptions

- **"The fragment gets sent to the server."** It does not, ever, in a standard HTTP request. It is resolved entirely client-side after the response is received. Servers cannot see or log fragment values from an HTTP request line (though JavaScript on the page can read `location.hash` and send it separately, e.g. in an XHR call, if the application chooses to).
- **"Putting a password in the URL (`user:pass@host`) is a normal way to authenticate."** It is part of the grammar and technically works with some tools (like `curl` or `ftp`), but browsers actively discourage and sometimes strip it from the display, and it is a serious security anti-pattern: URLs are logged in browser history, server access logs, proxy logs, and `Referer` headers (Chapter 71), all of which would then contain a live password in plaintext.
- **"A URL and a file path are basically the same idea."** They *look* similar (`/products/search` resembles a directory structure) for historical and readability reasons, but on virtually all modern web servers the path is routed by application code to arbitrary logic — there is usually no file on disk at that path at all.
- **"Query parameter order or `+` vs `%20` for spaces doesn't matter."** Order is not guaranteed to be preserved by every server or framework (though most do preserve it) and should not be relied upon for meaning. `+` as a space is a convention specific to the `application/x-www-form-urlencoded` content type used in query strings and form bodies — it is *not* part of the general percent-encoding rules in Section 6, and using literal `+` outside that context means a literal plus sign, not a space.
- **"HTTPS URLs are always on port 443."** Only by *default* (Section 7) — port 8443, 8080-with-TLS, or any other port is entirely legal if explicitly stated in the URL.
- **"A relative URL like `/about` and `about` mean the same thing."** They don't. A leading `/` is an *absolute path reference* — it replaces the entire path of the base URL. Without the leading slash, the reference is resolved relative to the base's *last path segment* (Section 11) — `about` from a base of `/products/search` resolves to `/products/about`, not `/about`. This single character is one of the most common sources of broken links when a site's directory structure changes.
- **"The whole URL is case-insensitive, like domain names."** Only the scheme and host are case-insensitive by spec. The path, query, and fragment are case-sensitive on any server that chooses to treat them that way — and most do, because they usually map to case-sensitive route definitions or filenames.

## 15. Production Notes

- **URL length limits are real but not standardized.** RFC 3986 sets no hard maximum, but practical limits exist everywhere in the stack: most browsers cap total URL length around 32,000–64,000 characters (varies by browser and by whether it's in the address bar vs. a link), and — critically for API design — many HTTP servers, proxies, and CDNs reject request lines beyond roughly 8,000 bytes by default (`LimitRequestLine` in Apache defaults to 8190 bytes; many CDNs cap similarly). This is a real, practical reason large parameter sets belong in a POST body, not a giant query string.
- **URL canonicalization matters for caching and SEO.** `https://Example.com:443/a/../b?` and `https://example.com/b` may be treated as the same resource by a well-behaved cache or search engine crawler (via canonicalization rules: lowercase host, remove default port, resolve `.`/`..` path segments) — but are treated as *different* cache keys by a naive one, causing duplicate cache entries or duplicate-content SEO penalties.
- **IDN homograph attacks are an active, ongoing concern.** Because Punycode (Section 6) allows visually near-identical domains using different Unicode code points (a Cyrillic "а" versus a Latin "a"), browsers apply heuristics to decide when to *display* the Punycode form (`xn--...`) instead of the decoded Unicode, specifically to make lookalike phishing domains visible rather than deceptively rendered.
- **Trailing slashes are not cosmetic.** `/products` and `/products/` are, by the raw HTTP/URL spec, different paths, and many frameworks (and REST API designs) treat them as genuinely different routes — or redirect one to the other with a 301 (Chapter 71), which costs an extra round trip if a client habitually gets it wrong.
- **Query parameters are the most common injection surface in real applications.** Every value that arrives after `?` is attacker-controllable, untyped, and unauthenticated by the URL layer itself — a query parameter used to build a database query without parameterization (SQL injection) or echoed back into a page without escaping (reflected cross-site scripting, previewed in Chapter 83) is one of the oldest and still most common web vulnerabilities in production systems. The URL's job stops at *delivering* the string; validating and safely using it is entirely the application's responsibility.
- **Logging full request URLs is a common accidental data leak.** Because query strings are just as visible in access logs, browser history, and proxy logs as anything else in the URL, session tokens or API keys accidentally passed as query parameters (instead of headers or a request body) end up persisted in plaintext across every log line the request passes through — a mistake real incident postmortems cite regularly.

## 16. Interview Questions & Model Answers

**Beginner: "What are the components of a URL?"**

"A URL has a scheme (like `https`), an authority section made up of an optional userinfo, a host, and an optional port, then a path identifying the resource, an optional query string of key-value parameters, and an optional fragment for jumping to a location within the loaded resource. For example, in `https://shop.example.com:8443/search?q=phones#reviews`, the scheme is `https`, the host is `shop.example.com`, the port is `8443`, the path is `/search`, the query is `q=phones`, and the fragment is `reviews`."

**Intermediate: "Does the server ever see the URL fragment?"**

"No. The fragment is defined by the URL spec to be resolved entirely on the client. When the browser builds the actual HTTP request line, it includes only the path and query string — the fragment is stripped off before the request is ever sent. If an application needs the server to know about the fragment's value, it has to explicitly read it in JavaScript and transmit it some other way, such as an additional query parameter or a follow-up API call."

**Advanced: "Walk through everything that happens, in order, from typing a URL to the first HTTP request byte leaving the machine, and identify which chapters/layers each step belongs to."**

"First the browser parses the typed string into its URL components — scheme, host, port, path, query, fragment — this is pure string handling, no network yet. It then checks its HTTP cache for a valid cached response; on a miss, it resolves the host to an IP address via DNS, which itself may involve checking local caches before walking root, TLD, and authoritative servers. With an IP address in hand, it opens a TCP connection to that IP on the URL's port — or the scheme's default port if none was given — via the three-way handshake, which costs one round trip. If the scheme is `https`, a TLS handshake follows on top of that TCP connection, negotiating a cipher and verifying the server's certificate, costing roughly one more round trip in TLS 1.3. Only after all of that completes does the browser write the actual HTTP request — method, path, query string, and headers — into the now-encrypted stream. So a URL that looks like one atomic action is really: URL parsing, cache check, DNS resolution, TCP handshake, TLS handshake, and only then HTTP."

## 17. Exercises

### Easy

1. Decompose `http://blog.example.org:8080/2026/articles/networking?utm_source=newsletter#comments` into its six components (scheme, host, port, path, query, fragment).
2. Explain, in one sentence, why `https://example.com` and `https://example.com:443` are identical requests.
3. Percent-encode the string `hello world & goodbye` so it could be safely used as a single query parameter value.

### Medium

4. A user types `192.168.1.5:5000/api/users` into a browser with no scheme. Explain what the browser is likely to do with this string, and why the outcome differs from typing `my search terms` with no dots or scheme at all.
5. You run `curl -w "%{time_namelookup} %{time_connect} %{time_appconnect} %{time_starttransfer}\n" -o /dev/null -s https://yoursite.example` twice in a row. Explain why the second run's DNS and TLS-related timings are typically much smaller than the first, tying your answer back to specific chapters.
6. Explain why embedding a password in a URL's userinfo section (`https://user:pass@host/`) is considered a security anti-pattern, referencing at least two different places that URL would end up logged in plaintext.

### Hard

7. A CDN in front of your site canonicalizes cache keys by lowercasing the host and stripping default ports, but does *not* sort query parameters. Explain a scenario where this causes a cache miss for what is semantically the same request, and propose a fix.
8. Design (in words, no code needed) a URL parser's handling of `https://a.com/x/../y//z?`. What should the normalized path be, and what ambiguity exists around the trailing empty path segment and the empty query string? Compare to how `net/url` in Go or the WHATWG URL Standard handles each case.
9. A phishing domain uses Cyrillic characters that render identically to `paypal.com` in most fonts. Explain, at the level of Punycode encoding and browser heuristics, how a modern browser is meant to detect and surface this to a user — and why this defense is not perfect.

## 18. Summary

| Term | Meaning |
|---|---|
| URI | Umbrella identifier for any resource; URL and URN are both kinds of URI |
| URL | A URI that also specifies how to locate/access the resource |
| Scheme | The protocol to use (`http`, `https`, `ftp`, `ws`, `wss`) |
| Authority | Userinfo + host + port section of a URL |
| Userinfo | Deprecated-in-practice credentials embedded in a URL |
| Host | Domain name or IP address, resolved via DNS (Ch. 66–68) |
| Port | TCP port; defaults per scheme if omitted (Ch. 57) |
| Path | Hierarchical identifier for a specific resource on the host |
| Query | `?key=value&...` parameters passed to the server |
| Fragment | `#id` client-side-only location marker; never sent to the server |
| Percent-encoding | Escaping reserved/non-ASCII bytes as `%XX` hex |
| Punycode | ASCII-safe encoding of internationalized domain names |
| Pre-HTTP work | DNS resolution (Ch. 66–68) → TCP handshake (Ch. 59) → optional TLS handshake (Ch. 82) — all before one HTTP byte is sent |

A URL is the address on the envelope. Chapter 71 opens that envelope and shows exactly what's written on the letter inside — the plain-text structure of the HTTP request and response that finally travel across the connection this chapter spent so much effort setting up.
