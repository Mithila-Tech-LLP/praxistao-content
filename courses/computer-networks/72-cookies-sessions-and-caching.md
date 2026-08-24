# Chapter 72: Cookies, Sessions, and Caching

> **"HTTP was designed to forget you the instant your response finishes arriving. Everything that makes the modern Web feel personal — being logged in, a shopping cart that remembers what you added, a page that loads instantly the second time — is a deliberate patch bolted on top of a protocol that was never supposed to remember anything at all."**

---

## Table of Contents

1. [The Problem: HTTP Has No Memory](#1-the-problem-http-has-no-memory)
2. [Why HTTP Was Designed Stateless in the First Place](#2-why-http-was-designed-stateless-in-the-first-place)
3. [Naive Fixes, and Why They Don't Scale](#3-naive-fixes-and-why-they-dont-scale)
4. [Cookies — The Mechanism That Adds State](#4-cookies--the-mechanism-that-adds-state)
5. [Cookie Attributes in Depth](#5-cookie-attributes-in-depth)
6. [Sessions — The Server-Side Pattern Built on Cookies](#6-sessions--the-server-side-pattern-built-on-cookies)
7. [A Sequence Diagram of Login, Then a Return Visit](#7-a-sequence-diagram-of-login-then-a-return-visit)
8. [The Other Half of the Story: Why Caching Exists](#8-the-other-half-of-the-story-why-caching-exists)
9. [Cache-Control Directives](#9-cache-control-directives)
10. [Validators: ETag and Last-Modified](#10-validators-etag-and-last-modified)
11. [Conditional Requests in Practice](#11-conditional-requests-in-practice)
12. [The Cache Hierarchy: Browser, Proxy, and CDN](#12-the-cache-hierarchy-browser-proxy-and-cdn)
13. [Deep Dive: How a Browser Decides "Fresh" vs. "Stale"](#13-deep-dive-how-a-browser-decides-fresh-vs-stale)
14. [Hands-On: Watching Cookies and Caching Live](#14-hands-on-watching-cookies-and-caching-live)
15. [Code: Setting Cookies and Handling Conditional Requests in Go](#15-code-setting-cookies-and-handling-conditional-requests-in-go)
16. [Common Misconceptions](#16-common-misconceptions)
17. [Production Notes](#17-production-notes)
18. [Interview Questions & Model Answers](#18-interview-questions--model-answers)
19. [Exercises](#19-exercises)
20. [Summary](#20-summary)

---

## 1. The Problem: HTTP Has No Memory

Chapter 71 ended on a pointed observation: the request/response cycle it described is **stateless**. Every request Section 6 of that chapter captured stood completely alone — the server had no idea, and no way to find out, whether this request came from someone who had visited five seconds ago or someone visiting for the very first time. Send the same GET request twice, from the same browser, one second apart, and the server produces two responses with zero awareness they're related.

This is fine for serving a static page. It is catastrophically insufficient for almost everything else the Web actually does:

- **Logging in.** You type a password once. Every subsequent page load needs the server to know "this is the person who logged in a moment ago" — but each of those page loads is, to raw HTTP, a brand-new, unrelated request.
- **A shopping cart.** You add three items across five different page views. Something has to remember that across requests that HTTP itself treats as unrelated strangers.
- **Personalization.** A site remembering your language preference, your dark-mode setting, your "remember me" checkbox — none of it survives from one stateless request to the next unless something outside of bare HTTP carries it forward.

## 2. Why HTTP Was Designed Stateless in the First Place

This wasn't an oversight — it was a deliberate, load-bearing design decision, and understanding *why* explains why the fix (cookies) had to be added on top rather than baked into HTTP's core:

- **Scalability.** A stateless server doesn't need to remember anything about who it talked to a moment ago. Any request can be handled by any server, in any order, with no coordination — this is precisely what lets a website behind a load balancer (Chapter 95) route each request to whichever of hundreds of backend servers happens to be free, rather than requiring every request from one user to hit the exact same machine that "remembers" them.
- **Simplicity and robustness.** A server that crashes and restarts loses nothing — there was never any per-client memory to lose. Compare this to a stateful protocol where a server crash mid-conversation leaves the client in an undefined, half-remembered state.
- **The original use case didn't need it.** HTTP/0.9 and HTTP/1.0 (Chapter 73) were built to serve documents — academic papers linking to other academic papers. A library handing you a book doesn't need to remember your previous checkouts to hand you the next one.

The tension this created is exactly the one Section 1 laid out: statelessness is a *feature* for scalability and robustness, but a *problem* the moment applications need continuity across requests. The Web's answer was not to make HTTP itself stateful, but to add a narrow, opt-in mechanism that lets state live at the edges (the client, and a server-side store keyed by something the client carries) while HTTP's core request/response mechanics (Chapter 71) stay exactly as stateless and simple as before.

## 3. Naive Fixes, and Why They Don't Scale

Before cookies existed (they were introduced by Netscape in 1994), early Web developers solved statelessness with two workarounds, both of which are useful to understand because their failure modes explain exactly what cookies needed to get right:

**URL rewriting** — embed a unique identifier directly in every link on the page:

```
https://shop.example.com/cart?sessionid=8f14e45fceea167a5a36
```

Every internal link on every page has to be rewritten to carry this parameter forward. It breaks the instant a user copies a URL and sends it to a friend (now sharing your identity), breaks if a user opens two tabs with different intents, pollutes every browser history entry and server log with a sensitive identifier (echoing Chapter 70's warning about credentials in URLs), and requires the *entire site* — every single link, every single form — to participate correctly, with a single missed link silently losing the state.

**Hidden form fields** — embed the identifier as an invisible `<input type="hidden">` in every HTML form, submitted back on every POST:

```html
<input type="hidden" name="sessionid" value="8f14e45fceea167a5a36">
```

This works only for form submissions, not for plain GET navigation between pages (clicking a link carries no form data at all) — so it can carry state *forward* through a checkout flow but can't maintain identity across ordinary browsing.

Both approaches share the same fundamental flaw: they force *every page and every link* on the entire site to actively participate in carrying the identifier forward. What was actually needed was a mechanism where the **browser itself** automatically attaches the identifier to every subsequent request to that server, without any cooperation from the page's HTML at all. That is exactly what a cookie is.

## 4. Cookies — The Mechanism That Adds State

A cookie is a small piece of data that a server asks the browser to store, and which the browser then automatically re-sends on every future request to that same server (subject to rules covered in Section 5) — without any page-level participation required.

**The server sets a cookie** with a `Set-Cookie` response header:

```
HTTP/1.1 200 OK
Set-Cookie: session_id=8f14e45fceea167a5a36; Path=/; HttpOnly; Secure; SameSite=Lax
Content-Type: text/html

<html>...</html>
```

**The browser stores it**, associated with the domain that sent it, and **automatically attaches it** to every subsequent request to that domain, via a `Cookie` request header — no HTML, no JavaScript, no form field required:

```
GET /cart HTTP/1.1
Host: shop.example.com
Cookie: session_id=8f14e45fceea167a5a36
```

This single mechanism — set once, sent automatically forever after (until expiry or deletion) — is the entire foundation. Everything from a login session to an A/B testing bucket assignment to an analytics visitor ID is built on top of exactly this `Set-Cookie` / `Cookie` header pair. Multiple cookies are separated by `;` in the `Cookie` header:

```
Cookie: session_id=8f14e45fceea167a5a36; theme=dark; cart_count=3
```

## 5. Cookie Attributes in Depth

The bare `name=value` pair is only the payload. Everything else in a `Set-Cookie` header is an **attribute** controlling exactly when and where the browser is allowed to send that cookie back:

```
Set-Cookie: session_id=8f14e45fceea167a5a36; Domain=example.com; Path=/account;
            Expires=Wed, 09 Sep 2026 10:00:00 GMT; Max-Age=2592000;
            Secure; HttpOnly; SameSite=Strict
```

| Attribute | Meaning |
|---|---|
| `Domain` | Which domain(s) the cookie is sent to. Omitted = only the exact host that set it. Set explicitly (`Domain=example.com`) = also sent to all subdomains (`shop.example.com`, `api.example.com`). |
| `Path` | Restricts the cookie to requests whose path starts with this prefix. `Path=/account` means the cookie is sent for `/account/settings` but not `/blog`. |
| `Expires` | An absolute date/time after which the browser deletes the cookie. |
| `Max-Age` | A relative lifetime in seconds from now; takes precedence over `Expires` if both are present. Omitting both makes it a **session cookie** — deleted when the browser closes, not based on a timer. |
| `Secure` | Only send this cookie over HTTPS, never plain HTTP. Prevents a passive network eavesdropper (Chapter 77's threat model) from ever seeing it in transit. |
| `HttpOnly` | Cookie is invisible to JavaScript (`document.cookie` cannot read it) — sent only by the browser's own network layer. This is a direct, deliberate defense against cross-site scripting (XSS) attacks stealing session cookies, previewed fully in Chapter 83. |
| `SameSite` | Controls whether the cookie is sent on cross-site requests (a request originating from a different site than the one the cookie belongs to). |

**`SameSite` deserves its own breakdown**, because it's the newest and most consequential of these attributes, added specifically to blunt a class of attack:

```
SameSite=Strict
  Cookie is NEVER sent on a cross-site request, even when a user clicks
  a link from another site to yours. Most secure, but can break a
  "click a link from an email, arrive already logged in" flow.

SameSite=Lax  (the modern browser DEFAULT if unspecified)
  Cookie IS sent on top-level navigation (clicking a link, typing a
  URL) even cross-site, but NOT sent on background cross-site
  requests (images, iframes, fetch/XHR calls embedded in another site).
  Balances usability and security — this is why "SameSite=Lax by
  default" became standard practice starting with Chrome 80 (2020).

SameSite=None
  Cookie is sent on all requests, cross-site or not. Requires the
  Secure attribute to be present as well — browsers reject
  `SameSite=None` without `Secure`. Needed for legitimate cross-site
  use cases, like a payment widget embedded via iframe on another
  company's checkout page.
```

`SameSite` is one of the concrete defenses against Cross-Site Request Forgery (CSRF), an attack where a malicious site tries to trigger state-changing requests (like `POST /transfer-money`) to your bank while your browser still automatically attaches your valid session cookie — previewed in full in Chapter 83.

## 6. Sessions — The Server-Side Pattern Built on Cookies

A cookie by itself is just a string the browser echoes back. It carries no meaning until a server chooses to *interpret* that string as a lookup key — and that interpretation is what a **session** is.

The pattern, concretely:

```
1. User logs in with username/password (POST /login).
2. Server verifies credentials, then generates a large, unguessable
   random session ID (e.g. 128+ bits of entropy — NOT sequential,
   NOT predictable, or an attacker could simply guess another
   user's session ID).
3. Server stores session data (user ID, permissions, cart contents,
   whatever the application needs) in a server-side store — an
   in-memory hash table for a single server, or more commonly in
   production, a shared store like Redis or a database table, so
   any backend server behind a load balancer can look it up.
4. Server responds with Set-Cookie: session_id=<the random ID>; HttpOnly; Secure
5. Every subsequent request from that browser carries
   Cookie: session_id=<the random ID>
6. Server looks up that ID in its session store, retrieves "this is
   user 42, logged in, cart has 3 items," and treats the request
   as if it remembers the user — even though HTTP itself (Ch. 71)
   never provided any such memory.
```

The cookie holds only an opaque, meaningless-to-outsiders identifier; **the actual state lives on the server**, keyed by that identifier. This split matters: the cookie can be small (session IDs are typically 16–32 bytes) regardless of how much session data exists server-side, and the server can invalidate a session instantly (delete it from the store) — something that is much harder to do cleanly with the alternative pattern of putting all session data *inside* a signed token the client holds (a JWT, common in modern API design), where the server has to actively track revocations rather than simply forgetting a row.

**Session fixation and hijacking**, in brief (Chapter 83 covers attacks in full): if an attacker can predict or steal a valid session ID, they can impersonate that user without ever knowing their password — which is exactly why session IDs need cryptographically strong randomness, `HttpOnly` (blocking JavaScript-based theft via XSS), `Secure` (blocking network-based theft via eavesdropping), and why well-designed systems regenerate the session ID at the moment of login (to invalidate any ID an attacker may have set before authentication occurred).

## 7. A Sequence Diagram of Login, Then a Return Visit

```mermaid
sequenceDiagram
    participant B as Browser
    participant S as Server
    participant DB as Session Store

    B->>S: POST /login  (username, password)
    S->>S: Verify credentials
    S->>DB: Create session: {id: "8f14...", user: 42}
    S-->>B: 200 OK<br/>Set-Cookie: session_id=8f14...; HttpOnly; Secure; SameSite=Lax
    Note over B: Browser stores the cookie silently

    B->>S: GET /account<br/>Cookie: session_id=8f14...
    S->>DB: Look up session "8f14..."
    DB-->>S: {user: 42, logged in}
    S-->>B: 200 OK  (personalized account page for user 42)

    Note over B,S: Hours later, new tab, same browser
    B->>S: GET /cart<br/>Cookie: session_id=8f14...
    Note right of S: Same session ID, same lookup —<br/>server "remembers" without HTTP itself remembering anything
    S-->>B: 200 OK  (cart still has the items from earlier)
```

## 8. The Other Half of the Story: Why Caching Exists

Cookies solve "how does the server remember me." Caching solves an entirely different, equally fundamental problem: **most of what gets requested over and over hasn't actually changed**, and re-fetching it every single time wastes time, bandwidth, and server capacity for no benefit.

Consider loading one modern web page: the HTML itself might be small and genuinely dynamic, but it typically references dozens of other resources — a CSS file, a JavaScript bundle, a logo image, web fonts — almost none of which have changed since the last time this browser (or a million other browsers) fetched them. Re-running the full pre-HTTP journey from Chapter 70 (DNS, TCP handshake, TLS handshake) and the full request/response cycle from Chapter 71 for every one of those unchanged resources, on every single page view, is enormous waste:

- **Latency the user actually feels** — every skipped round trip is milliseconds saved, and a page with 50 resources has 50 chances to save time.
- **Bandwidth**, for both the user (mobile data costs, slow connections) and the server (egress costs at scale).
- **Server load** — a server that can answer "you already have this, nothing changed" for 90% of its traffic needs far less capacity than one regenerating every response from scratch.

HTTP caching is the mechanism for skipping unnecessary requests entirely, or — when a full skip isn't safe — replacing a full response with a tiny "nothing changed" confirmation instead.

## 9. Cache-Control Directives

The `Cache-Control` header, sent by the server in a response, is the primary mechanism controlling caching behavior. It carries one or more comma-separated directives:

| Directive | Meaning |
|---|---|
| `max-age=<seconds>` | This response may be reused, without contacting the server at all, for up to this many seconds. |
| `no-cache` | Confusingly named: the response *can* be cached, but must be **revalidated** with the server (Section 10/11) before each reuse — not "don't cache," but "don't reuse blindly." |
| `no-store` | The real "don't cache" — this response must not be stored anywhere at all, not even temporarily. Used for sensitive data (banking, personal info). |
| `private` | May be cached only by the end user's own browser, not by a shared intermediary like a CDN or corporate proxy (Section 12) — appropriate for personalized responses. |
| `public` | May be cached by any cache, including shared ones — appropriate for content identical for every user. |
| `must-revalidate` | Once the response becomes stale (past `max-age`), it must not be used at all without successful revalidation — no serving a stale copy even if the server is momentarily unreachable. |
| `immutable` | The resource will never change for the duration it's considered fresh — tells the browser not to even bother revalidating on a user-triggered reload. |

A realistic example, for a hashed, versioned JavaScript bundle that truly never changes once published:

```
Cache-Control: public, max-age=31536000, immutable
```

(31,536,000 seconds = 365 days — the standard "cache this basically forever" value, used together with content-hashed filenames like `app.a3f9c2.js` so that when the content *does* change, it gets a new URL entirely rather than trying to invalidate the old one.)

By contrast, a personalized account page:

```
Cache-Control: private, no-store
```

**`Expires` is the older, HTTP/1.0-era header** that `Cache-Control: max-age` effectively superseded — it specifies an absolute date/time rather than a relative duration. Modern responses often send both, for backward compatibility with very old caches that don't understand `Cache-Control`, but `max-age` takes precedence whenever both are present and understood.

## 10. Validators: ETag and Last-Modified

`max-age` answers "how long can I trust this without asking again," but eventually every cached copy goes stale, and the client needs to ask, "is my stale copy actually still correct?" **Validators** are what makes that question cheap to answer.

**`ETag`** ("entity tag") is an opaque identifier — usually a hash of the resource's content — that changes if and only if the content changes:

```
HTTP/1.1 200 OK
ETag: "3147526947+gzip"
Cache-Control: max-age=3600
```

**`Last-Modified`** is the older, coarser alternative — a timestamp of when the resource last changed:

```
HTTP/1.1 200 OK
Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
Cache-Control: max-age=3600
```

`ETag` is generally preferred where available because it detects *any* content change precisely (even a change that happens within the same second, which a timestamp-based check can't distinguish), while `Last-Modified` has only one-second resolution and can be fooled by a file being rewritten with identical content but a new timestamp.

## 11. Conditional Requests in Practice

Once a client holds a validator from a previous response, it can ask the server "has this changed since I last saw it?" using a **conditional request** — and if the answer is no, the server can respond with almost no data at all.

**With `ETag`, using `If-None-Match`:**

```
GET /app.js HTTP/1.1
Host: cdn.example.com
If-None-Match: "3147526947+gzip"
```

If the resource hasn't changed, the server responds:

```
HTTP/1.1 304 Not Modified
Date: Sun, 09 Aug 2026 09:20:11 GMT
Cache-Control: max-age=3600

```

No body at all — `304 Not Modified` (introduced in Chapter 71's status code table) carries only headers, telling the browser "the copy you already have is still correct, and here's a refreshed expiry time; go ahead and keep using it." If the resource *has* changed, the server simply responds normally with `200 OK`, a new `ETag`, and the full new body.

**With `Last-Modified`, using `If-Modified-Since`:**

```
GET /app.js HTTP/1.1
Host: cdn.example.com
If-Modified-Since: Thu, 17 Oct 2019 07:18:26 GMT
```

Same idea, coarser precision. Both mechanisms turn a potentially large response into a tiny header-only exchange whenever nothing has actually changed — this is precisely "the mechanism that makes the Web fast by skipping unnecessary [full] requests" that Section 8 promised: the request still happens (a conditional GET is still a full round trip), but the *expensive part* — re-sending the entire body — is skipped whenever possible.

## 12. The Cache Hierarchy: Browser, Proxy, and CDN

Caching doesn't happen in just one place. A single response, as it travels from an origin server back to a user, may pass through several independent caches, each capable of satisfying future requests without going any further upstream:

```
User's browser cache  (private — Section 9's `private` directive matters here)
        ↑
Corporate/ISP forward proxy cache  (shared — rare today with widespread HTTPS,
                                     since a proxy can't read/cache encrypted bodies
                                     it can't decrypt)
        ↑
CDN edge cache  (shared — Chapter 96 covers this in full; a `public` response
                 can be cached here and served to thousands of different users
                 from a server physically near them, without ever reaching
                 the origin)
        ↑
Origin server
```

This is exactly why the `private` vs. `public` distinction in Section 9 matters operationally: a `private` response (like a personalized account page) must never be cached at the CDN layer, because the *next* request to that same CDN edge might come from a completely different user who should never see someone else's personalized content. A `public` response (a shared logo image, a CSS file) benefits enormously from being cached at every layer, since the same bytes are correct for every single visitor.

## 13. Deep Dive: How a Browser Decides "Fresh" vs. "Stale"

Putting Sections 9–11 together, here is the actual decision a browser (or any HTTP cache) makes every time it's asked to fetch a resource it may have already cached:

```
1. Do I have a cached response for this exact request (URL + relevant
   headers, per the Vary header — see below)?
     No  → make a normal, full request. Cache the response per its
           Cache-Control headers if cacheable.
     Yes → go to step 2.

2. Is the cached response still "fresh" (within Cache-Control: max-age,
   or before its Expires date)?
     Yes → use it directly. ZERO network requests. This is the fastest
           possible outcome — Section 8's whole point.
     No  → the response is "stale." Go to step 3.

3. Does the stale response have a validator (ETag or Last-Modified)?
     No  → treat it as gone; make a normal full request.
     Yes → make a CONDITIONAL request (Section 11) with
           If-None-Match / If-Modified-Since.

4. Server responds:
     304 Not Modified → reuse the stale copy, refresh its freshness
                         lifetime, no body was transferred.
     200 OK (new body) → replace the cached copy entirely with this
                          new response.
```

The `Vary` header (seen in Chapter 71's captured example, `Vary: Accept-Encoding`) refines step 1's cache lookup: it tells the cache "this response differs depending on the value of these request headers," so a cache can't reuse a gzip-compressed response for a client that sent a different `Accept-Encoding` — this is exactly why the earlier capture included both `Content-Encoding: gzip` and `Vary: Accept-Encoding` together.

## 14. Hands-On: Watching Cookies and Caching Live

**Watch a server set a cookie and a client send it back**, in two separate raw requests:

```
$ curl -v -c cookies.txt https://httpbin.org/cookies/set/session_id/8f14e45f
> GET /cookies/set/session_id/8f14e45f HTTP/1.1
< HTTP/1.1 302 Found
< Set-Cookie: session_id=8f14e45f; Path=/

$ curl -v -b cookies.txt https://httpbin.org/cookies
> GET /cookies HTTP/1.1
> Cookie: session_id=8f14e45f
< HTTP/1.1 200 OK
{
  "cookies": {"session_id": "8f14e45f"}
}
```

(`-c cookies.txt` tells `curl` to save cookies from `Set-Cookie` responses; `-b cookies.txt` tells it to send them back — this is `curl` acting as a tiny, manual browser cookie jar.)

**Watch a conditional request happen for real**, against any server that supports `ETag`:

```
$ curl -v -o /dev/null https://example.com/style.css
< HTTP/1.1 200 OK
< ETag: "60d5-5e6e5e6e5e6e5"
< Cache-Control: max-age=3600

$ curl -v -o /dev/null -H 'If-None-Match: "60d5-5e6e5e6e5e6e5"' https://example.com/style.css
< HTTP/1.1 304 Not Modified
< ETag: "60d5-5e6e5e6e5e6e5"
```

Notice the second response has no body and a tiny header block — exactly the savings Section 11 described.

**Watch it happen in a real browser**, using developer tools: open the Network tab, reload a page you've visited before, and look at the "Size" column — entries showing `(disk cache)` or `(memory cache)` never touched the network at all (Section 13, step 2); entries with a real byte size but a fast time and status `304` used a conditional request (step 3–4).

## 15. Code: Setting Cookies and Handling Conditional Requests in Go

```go
package main

import (
	"crypto/sha256"
	"fmt"
	"net/http"
)

var sessionStore = map[string]int{} // sessionID -> userID, an in-memory session store

func loginHandler(w http.ResponseWriter, r *http.Request) {
	// (Assume credentials were already verified above this line.)
	sessionID := "8f14e45fceea167a5a36" // in reality: a cryptographically random value
	sessionStore[sessionID] = 42        // user 42 is now logged in

	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    sessionID,
		Path:     "/",
		HttpOnly: true, // JavaScript cannot read this cookie (Section 5)
		Secure:   true, // only sent over HTTPS
		SameSite: http.SameSiteLaxMode,
		MaxAge:   30 * 24 * 60 * 60, // 30 days
	})
	fmt.Fprintln(w, "Logged in")
}

func accountHandler(w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		http.Error(w, "not logged in", http.StatusUnauthorized) // 401, Ch. 71
		return
	}
	userID, ok := sessionStore[cookie.Value]
	if !ok {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}
	fmt.Fprintf(w, "Welcome back, user %d\n", userID)
}

func cachedAssetHandler(w http.ResponseWriter, r *http.Request) {
	body := []byte("body { color: black; }") // pretend this is a CSS file's bytes
	etag := fmt.Sprintf(`"%x"`, sha256.Sum256(body))

	w.Header().Set("ETag", etag)
	w.Header().Set("Cache-Control", "public, max-age=3600")

	if r.Header.Get("If-None-Match") == etag {
		w.WriteHeader(http.StatusNotModified) // 304 — no body sent (Section 11)
		return
	}
	w.Write(body)
}

func main() {
	http.HandleFunc("/login", loginHandler)
	http.HandleFunc("/account", accountHandler)
	http.HandleFunc("/style.css", cachedAssetHandler)
	http.ListenAndServe(":8080", nil)
}
```

## 16. Common Misconceptions

- **"Cookies are the same thing as sessions."** A cookie is just a small piece of data the browser stores and re-sends. A session is a *pattern* — server-side storage keyed by a value the cookie happens to carry. You can use cookies without sessions (e.g., storing a theme preference directly in the cookie's value, with no server-side lookup at all) and, less commonly, maintain session-like state without cookies (e.g., a token in a mobile app's `Authorization` header instead).
- **"`Cache-Control: no-cache` means the browser won't cache the response."** It means the opposite of what it sounds like: the response *can* be stored, but must be revalidated (Section 11) before every reuse. `no-store` is the directive that actually forbids storage.
- **"`HttpOnly` encrypts the cookie."** It does not — it only prevents JavaScript from reading the cookie's value via `document.cookie`. The cookie is still plainly visible in the `Cookie` header on the wire (unless `Secure` + HTTPS is also used) and to anyone with access to the browser's cookie storage on disk.
- **"A 304 response means the request failed."** It's a success signal, just one that says "nothing to send, use what you already have" rather than delivering fresh content — it belongs squarely in Chapter 71's 3xx-adjacent-but-actually-successful category of "the client's cache is doing its job."
- **"Deleting a cookie logs a user out everywhere."** Deleting a cookie in one browser removes that browser's copy of the session identifier, but the session itself may still be perfectly valid and present in the server-side session store (Section 6) — a proper logout must also tell the server to invalidate the session server-side, or the same session ID (if somehow captured by an attacker beforehand) would remain usable.

## 17. Production Notes

- **Session stores don't stay in a single process's memory in production.** A real deployment behind a load balancer (Chapter 95) needs every backend instance to be able to look up any session, so the in-memory map in Section 15's example code becomes a shared store like Redis or a database table in practice — this is a direct consequence of statelessness at the HTTP layer forcing state to live somewhere that scales independently of any one server.
- **Cache invalidation is famously one of the hardest problems in computing, and this is exactly where it shows up.** Setting `max-age=31536000` on a resource that later needs to change urgently (a security fix in a JS bundle) means clients with a cached copy simply won't ask again for up to a year — the standard mitigation is content-hashed filenames (Section 9), so a changed file gets a brand-new URL and the old cached entry becomes irrelevant rather than needing active invalidation.
- **Third-party cookies are being phased out across major browsers**, specifically because a cookie set by an embedded ad or tracking script on one site, then sent back on every other site embedding the same script, enables cross-site tracking that `SameSite=Lax`-as-default (Section 5) and, more aggressively, full third-party cookie blocking are designed to prevent — an active, ongoing shift in browser privacy engineering as of the mid-2020s.
- **`Secure` and `HttpOnly` should be treated as close to mandatory for any session cookie in production**, not optional hardening — omitting `Secure` means a session ID could be sent in the clear over an accidental plain-HTTP request, and omitting `HttpOnly` means a single XSS vulnerability (Chapter 83) anywhere on the site can exfiltrate every visiting user's session identifier via `document.cookie`.

## 18. Interview Questions & Model Answers

**Beginner: "Why is HTTP called stateless, and what problem does that create?"**

"Each HTTP request/response cycle is handled independently — the server retains no memory of any previous request from the same client once a response is sent. This creates a problem the moment an application needs continuity across requests, like staying logged in or keeping items in a shopping cart, because without an added mechanism, every single page load looks to the server like a brand-new, unrelated visitor."

**Intermediate: "Explain how a login session actually works, end to end, using cookies."**

"After verifying credentials, the server generates a large random session ID and stores session data — like the user's identity — in a server-side store keyed by that ID. It sends the ID back to the browser in a Set-Cookie header, ideally with HttpOnly and Secure attributes. The browser automatically attaches that cookie to every subsequent request to the same site via the Cookie header. On each request, the server looks up the session ID in its store, retrieves the associated user data, and treats the request as authenticated — even though the underlying HTTP protocol itself never remembers anything between requests. The state lives in the cookie's opaque reference plus the server-side store, not in HTTP itself."

**Advanced: "Explain the full lifecycle of a conditional GET request, and why it's more useful than either 'always cache blindly' or 'never cache at all.'"**

"A conditional GET sits between those two extremes. On the first request, the server returns a full response along with a validator — an ETag or a Last-Modified timestamp — plus a Cache-Control max-age telling the client how long to trust the response without asking again. Once that freshness window expires, instead of either blindly reusing a possibly-stale copy or re-downloading the entire resource from scratch, the client sends the validator back in an If-None-Match or If-Modified-Since header. If the resource hasn't actually changed, the server responds 304 Not Modified with no body at all — just refreshed headers — letting the client keep using its existing cached copy. If it has changed, the server sends a full 200 response with new content and a new validator. This gets you correctness — you never serve genuinely stale data past what the validator confirms — at a fraction of the bandwidth cost of a full re-fetch, which is the entire reason caching headers exist beyond a simple time-based expiry."

## 19. Exercises

### Easy

1. Write the exact `Set-Cookie` header a server should send to store a 30-day session cookie named `sid` with value `abc123`, restricted to HTTPS, inaccessible to JavaScript, and sent on top-level cross-site navigation but not background cross-site requests.
2. Explain the difference between `Cache-Control: no-cache` and `Cache-Control: no-store` in one sentence each.
3. A client sends `If-None-Match: "v1"` and the server's current ETag for that resource is `"v2"`. What status code should the server return, and what should the response body contain?

### Medium

4. A shopping site sets a session cookie without the `Secure` attribute, and a user occasionally accesses the site over plain HTTP on public Wi-Fi. Explain the specific risk this creates, referencing Chapter 77's threat model of a passive network eavesdropper.
5. Explain why a CDN caching a `private` response for one logged-in user and then serving it to a different logged-in user would be a serious bug, and identify which Cache-Control directive exists specifically to prevent it.
6. A single-page app stores its JWT auth token in `localStorage` instead of an HttpOnly cookie, reasoning "so our JavaScript can read it and attach it to API calls." Identify the security trade-off this introduces compared to an HttpOnly session cookie.

### Hard

7. Design a cache-invalidation strategy for a company's main JavaScript bundle that is deployed multiple times per day, using content-hashed filenames and long `max-age`, and explain precisely why this avoids needing to actively invalidate anything in any downstream cache (browser, CDN).
8. A user reports "I logged out, but when I click back in my browser, the account page still shows my data." Using Sections 6, 12, and 13, explain the two independent, non-mutually-exclusive causes this could have (one about the session store, one about caching) and how you'd distinguish them.
9. An attacker embeds `<img src="https://bank.example.com/transfer?to=attacker&amount=1000">` on a malicious page. If bank.example.com's transfer endpoint incorrectly accepts GET requests and the user's browser has a valid, non-`SameSite=Strict` session cookie for the bank, trace exactly why this request succeeds, and identify two independent fixes — one from this chapter (a cookie attribute) and one from Chapter 71 (a method-semantics fix).

## 20. Summary

| Term | Meaning |
|---|---|
| Statelessness | HTTP's design property: no request remembers any prior request by default |
| Cookie | Small data set by a server (`Set-Cookie`) and auto-returned by the browser (`Cookie`) |
| Session | Server-side pattern: a cookie carries an opaque ID that keys server-stored state |
| HttpOnly | Cookie attribute blocking JavaScript access — mitigates XSS-based cookie theft |
| Secure | Cookie attribute restricting transmission to HTTPS only |
| SameSite | Cookie attribute controlling cross-site sending — mitigates CSRF |
| Cache-Control | Response header directing how/whether a response may be cached and reused |
| max-age | Cache-Control directive: seconds a response can be reused with zero network requests |
| ETag / Last-Modified | Validators identifying a specific version of a resource |
| Conditional request | A GET carrying If-None-Match/If-Modified-Since, answered with 304 if unchanged |
| 304 Not Modified | Success status meaning "use your cached copy, nothing has changed" |

Cookies and caching both quietly assume something Chapter 71 didn't dwell on: that a browser can keep a single connection to a server open and reuse it efficiently across many requests, rather than paying for a brand-new TCP (and possibly TLS) handshake every single time. Chapter 73 finally examines that assumption directly — starting from HTTP/1.0's costly one-request-per-connection design, through HTTP/1.1's fix, and into the surprisingly stubborn head-of-line-blocking problem that fix didn't actually solve.
