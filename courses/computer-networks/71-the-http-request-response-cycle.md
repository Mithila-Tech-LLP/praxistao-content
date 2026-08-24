# Chapter 71: The HTTP Request/Response Cycle — Methods, Headers, Status Codes

> **"HTTP is deliberately, almost stubbornly simple: one side sends a plain-text message that says what it wants, the other side sends back a plain-text message that says what happened. Everything the Web does — search, shopping, banking, video — is built on top of that one exchange, repeated billions of times a second."**

---

## Table of Contents

1. [Where Chapter 70 Left Off](#1-where-chapter-70-left-off)
2. [The Problem: Agreeing on a Message Format](#2-the-problem-agreeing-on-a-message-format)
3. [A Naive Attempt, and Why Structure Is Needed](#3-a-naive-attempt-and-why-structure-is-needed)
4. [The Anatomy of an HTTP Request](#4-the-anatomy-of-an-http-request)
5. [The Anatomy of an HTTP Response](#5-the-anatomy-of-an-http-response)
6. [A Real, Fully Captured Request/Response Pair](#6-a-real-fully-captured-requestresponse-pair)
7. [HTTP Methods — What Each One Actually Means](#7-http-methods--what-each-one-actually-means)
8. [Safe, Idempotent, and Cacheable — The Properties That Matter](#8-safe-idempotent-and-cacheable--the-properties-that-matter)
9. [Headers as Metadata](#9-headers-as-metadata)
10. [Status Codes, Grouped by Class](#10-status-codes-grouped-by-class)
11. [A Sequence Diagram of One Full Cycle](#11-a-sequence-diagram-of-one-full-cycle)
12. [Deep Dive: How a Server Actually Parses a Request](#12-deep-dive-how-a-server-actually-parses-a-request)
13. [Hands-On: Talking Raw HTTP by Hand](#13-hands-on-talking-raw-http-by-hand)
14. [Code: A Minimal HTTP Server and Client in Go](#14-code-a-minimal-http-server-and-client-in-go)
15. [Common Misconceptions](#15-common-misconceptions)
16. [Production Notes](#16-production-notes)
17. [Interview Questions & Model Answers](#17-interview-questions--model-answers)
18. [Exercises](#18-exercises)
19. [Summary](#19-summary)

---

## 1. Where Chapter 70 Left Off

Chapter 70 ended with a TCP connection open (and, for HTTPS, a TLS session established on top of it) but not one byte of HTTP had been sent yet. All that work — DNS resolution, the three-way handshake, the TLS handshake — produced exactly one thing: a reliable, ordered, private pipe between the browser and the server. This chapter is about what actually gets written into that pipe.

## 2. The Problem: Agreeing on a Message Format

A TCP connection (Chapter 59) is just a byte stream — it has no concept of "messages" at all, only bytes flowing in each direction. If the client writes bytes into the connection and the server reads bytes out of it, both sides need to agree, in advance, on exactly how to interpret that stream as *requests* and *responses*: where one message ends and the next begins, what a "request" even means, and how to communicate outcomes like "the page doesn't exist" or "you're not authorized" without inventing a new format for every possible situation.

This is precisely the kind of problem Chapter 24 introduced in the abstract — a shared protocol above a shared channel — and HTTP is the concrete answer for the Web.

## 3. A Naive Attempt, and Why Structure Is Needed

Suppose you tried to invent this yourself. The absolute simplest idea: the client just writes the path it wants (`/index.html`) into the connection, and the server writes back the file's raw bytes. This is, in fact, almost exactly what the very first version of HTTP (retroactively called HTTP/0.9, 1991) did — no headers, no status line, no method other than an implicit "give me this":

```
Client sends:  GET /index.html

Server sends:  <html>...raw file contents...</html>
              (then closes the connection — that's how the client
               knows the response is complete)
```

This works for "serve me a static file" and nothing else. It cannot express:

- **What kind of content is this?** (HTML? An image? JSON? The client has to guess from the bytes or the file extension.)
- **Did it work?** (There's no way to say "not found" except sending back something that looks like an error page, which a program can't distinguish from a real page.)
- **Extra information about the request or response** (When was this last modified? Should the browser cache it? What language does the client prefer?)
- **Anything other than "fetch a file."** (Submit a form? Delete something? Update a record?)

Every one of these gaps is exactly what real HTTP (from 1.0 onward) fixes, using three additions: a **method** (what kind of action), **headers** (structured metadata, as key-value pairs), and a **status code** (a standardized, machine-readable outcome).

## 4. The Anatomy of an HTTP Request

An HTTP/1.1 request is plain, human-readable ASCII text (the body may be binary, but the request line and headers are always text), with this exact structure:

```
METHOD SP request-target SP HTTP-version CRLF        ← Request line
Header-Name: header-value CRLF                        ← Headers, one per line
Header-Name: header-value CRLF
...
CRLF                                                   ← Blank line: end of headers
[optional message body]                                ← Body (if present)
```

`SP` is a literal space, `CRLF` is a literal carriage-return + line-feed (`\r\n` — not just `\n`; this is specified precisely and real servers reject malformed line endings). Concretely:

```
GET /products/search?q=laptops&sort=price HTTP/1.1
Host: shop.example.com
User-Agent: Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7)
Accept: text/html,application/xhtml+xml
Accept-Language: en-US,en;q=0.9
Connection: keep-alive

```

Notice the request-target here is exactly the path and query string from the URL Chapter 70 decomposed (`/products/search?q=laptops&sort=price`) — the scheme, host, and port from the URL are *not* repeated in the request line; the host instead travels in the mandatory `Host` header (this single header is what makes virtual hosting, mentioned in Chapter 70 Section 2, possible: one IP address, one TCP connection, but the server still learns which of its many hosted domains the client meant). There is no body on this GET request, so the message ends at the blank line.

## 5. The Anatomy of an HTTP Response

The response mirrors the request's shape, replacing the request line with a **status line**:

```
HTTP-version SP status-code SP reason-phrase CRLF     ← Status line
Header-Name: header-value CRLF                         ← Headers
Header-Name: header-value CRLF
...
CRLF                                                    ← Blank line: end of headers
[optional message body]                                 ← Body
```

Concretely, answering the request above:

```
HTTP/1.1 200 OK
Content-Type: text/html; charset=utf-8
Content-Length: 1274
Date: Sun, 09 Aug 2026 09:12:03 GMT
Server: nginx/1.25.3
Cache-Control: no-store

<!DOCTYPE html><html><head><title>Search results</title></head>
<body><h1>12 laptops found</h1>...</body></html>
```

`200` is the status code, `OK` is the reason phrase (a human-readable label with no machine meaning — a client is required to key off the numeric code, not the text, which can even be omitted or customized by non-compliant servers). `Content-Length: 1274` tells the client exactly how many bytes of body to expect, which — on a persistent connection (Chapter 73) — is precisely how the client knows where this response ends and where a *second* response on the same connection would begin.

## 6. A Real, Fully Captured Request/Response Pair

Here is an unedited-in-substance capture of what a browser actually sends and receives for a simple page load, as you'd see it with `curl -v` or a raw `nc` (netcat) session against a real server (details like the exact `Date` and `ETag` values will differ on every real run, but the shape is exact):

```
> GET / HTTP/1.1
> Host: example.com
> User-Agent: curl/8.4.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Encoding: gzip
< Accept-Ranges: bytes
< Age: 201573
< Cache-Control: max-age=604800
< Content-Type: text/html; charset=UTF-8
< Date: Sun, 09 Aug 2026 09:14:22 GMT
< Etag: "3147526947+gzip"
< Expires: Sun, 16 Aug 2026 09:14:22 GMT
< Last-Modified: Thu, 17 Oct 2019 07:18:26 GMT
< Server: ECAcc (dcd/7D5A)
< Vary: Accept-Encoding
< Content-Length: 648
<
<!doctype html>
<html>
<head>
    <title>Example Domain</title>
    ...
</head>
<body>
<div>
    <h1>Example Domain</h1>
    <p>This domain is for use in illustrative examples in documents...</p>
</div>
</body>
</html>
```

(`>` lines are what `curl` sent; `<` lines are what the server sent back — this convention is used throughout this chapter and the rest of Volume 11.) Every header here will be explained in Section 9 and Chapter 72; for now, notice how much richer this is than the naive HTTP/0.9 attempt from Section 3 — content type, encoding, length, caching hints, a validator (`ETag`), and a real status code, all present in one small text block.

## 7. HTTP Methods — What Each One Actually Means

A method is the verb of the request — it tells the server *what kind of action* this is, independent of the path. HTTP defines a fixed vocabulary rather than letting every application invent its own verbs, precisely so that generic infrastructure (caches, proxies, browsers) can apply universal rules based on the method alone, without understanding the application.

| Method | Semantic meaning | Typical use |
|---|---|---|
| `GET` | Retrieve a representation of a resource. Must not change server state. | Loading a page, fetching an API resource |
| `HEAD` | Exactly like GET, but return only headers, no body. | Checking if a resource exists, its size, or its last-modified time, without downloading it |
| `POST` | Submit data to be processed by the resource, often creating something new whose identity the server decides. | Submitting a form, creating a new record via an API |
| `PUT` | Replace the resource at this exact URL entirely with the supplied representation. | Uploading a file to a known path; "set this resource to be exactly this" |
| `PATCH` | Apply a partial modification to a resource. | Updating one field of a record without resending the whole thing |
| `DELETE` | Remove the resource at this URL. | Deleting a record |
| `OPTIONS` | Ask what methods/capabilities are supported for a resource, without performing any action. | CORS preflight checks (a browser asking a cross-origin server "am I allowed to POST here?") |
| `CONNECT` | Establish a tunnel to the server, typically for HTTPS through a proxy. | An HTTP proxy tunneling a TLS connection through to the real destination |
| `TRACE` | Echo back the received request, for diagnostic loop-back testing. | Debugging what intermediaries changed about a request (rarely used; often disabled for security reasons) |

**Choosing between PUT and PATCH, concretely:** if a client sends the *entire* representation of a user record and means "this is the complete, final state of user 42," that's `PUT /users/42`. If a client sends `{"email": "new@example.com"}` and means "just change the email field, leave everything else alone," that's `PATCH /users/42`. Using `POST /users/42/update` for either is common in practice but semantically muddier — it hides the actual operation inside the path instead of expressing it through the method, which is exactly the kind of ad-hoc-verb problem the fixed method vocabulary exists to avoid.

**Choosing between POST and PUT for creation:** `POST /users` (no ID in the path) is correct when the *server* assigns the new resource's identity (e.g., an auto-incrementing user ID) — the client doesn't know the final URL until the server responds. `PUT /users/42` is correct when the *client* already knows and dictates the exact resource identity up front.

## 8. Safe, Idempotent, and Cacheable — The Properties That Matter

Three formal properties, defined by the HTTP specification per method, determine how infrastructure (browsers, proxies, retry logic) is allowed to treat a request automatically:

- **Safe** — the method is not expected to cause any observable side effect on the server; it's purely a read. A safe method can be prefetched, retried freely, or followed speculatively by a crawler without concern for consequences.
- **Idempotent** — making the same request multiple times has the same effect as making it once. This matters enormously for retries: if a network glitch makes a client unsure whether a request actually reached the server, it is *safe to blindly retry* an idempotent request, but not a non-idempotent one.
- **Cacheable** — a response to this method may, under the right headers (Chapter 72), be stored and reused for a later identical request instead of going to the server again.

| Method | Safe? | Idempotent? | Cacheable (by default)? |
|---|---|---|---|
| `GET` | Yes | Yes | Yes |
| `HEAD` | Yes | Yes | Yes |
| `OPTIONS` | Yes | Yes | No |
| `PUT` | No | Yes | No |
| `DELETE` | No | Yes | No |
| `POST` | No | No | No (unless explicit caching headers say otherwise) |
| `PATCH` | No | No | No |

Note the sharp, easy-to-miss distinction: `DELETE` is **not safe** (it clearly changes server state) but **is idempotent** — deleting an already-deleted resource still results in "the resource does not exist," the same end state, even if the second call returns `404` instead of `200`. `PUT` is idempotent for the identical reason: replacing a resource with the exact same representation twice in a row leaves the resource in the same final state either way. `POST`, by contrast, is neither — calling `POST /orders` twice with an identical body typically creates *two* separate orders, which is precisely why double-clicking a "Submit Order" button is a real, historically expensive bug class, and why idempotency keys (a client-generated unique ID attached to a POST, which the server deduplicates against) are a common production pattern specifically to make an inherently non-idempotent method behave safely under retries.

## 9. Headers as Metadata

Headers are the mechanism that carries everything the request/response line itself doesn't — think of the request line and status line as the envelope's destination and outcome stamp, and headers as everything written on the outside of the envelope in the margins. They are grouped, informally, into categories:

```
General headers   — apply to the message itself, either direction
                     Date, Connection, Cache-Control

Request headers    — describe the client or what it will accept
                     Host, User-Agent, Accept, Accept-Language,
                     Accept-Encoding, Authorization, Cookie (Ch. 72)

Response headers    — describe the server or the response's circumstances
                     Server, Set-Cookie (Ch. 72), WWW-Authenticate,
                     Location (used with redirects, Section 10)

Representation headers — describe the body/payload itself
                     Content-Type, Content-Length, Content-Encoding,
                     Content-Language, ETag, Last-Modified (Ch. 72)
```

A few headers deserve individual attention because of how much depends on them:

- **`Host`** — mandatory in HTTP/1.1 (Section 4). Without it, a server hosting multiple domains on one IP has no way to know which site the client means.
- **`Content-Type`** — tells the receiver how to interpret the body's bytes: `text/html`, `application/json`, `image/png`, `application/x-www-form-urlencoded`, `multipart/form-data` (used for file uploads). Getting this wrong is one of the most common real integration bugs — a server expecting JSON but receiving `Content-Type: text/plain` may refuse to parse a perfectly valid JSON body.
- **`Content-Length`** — the exact byte count of the body. Essential for framing multiple requests/responses on one persistent connection (Chapter 73) — without it (or its alternative, chunked transfer encoding), the receiver has no way to know where the message ends short of the connection closing.
- **`Accept` / `Accept-Language` / `Accept-Encoding`** — the client stating its preferences: which content types it can render, which human language it prefers, which compression algorithms (gzip, br) it can decompress. The server uses these to pick the best available representation — this negotiation is literally called "content negotiation."
- **`Authorization`** — carries credentials, most commonly `Authorization: Bearer <token>` for API tokens or `Authorization: Basic <base64(user:pass)>` for HTTP Basic auth.
- **`Referer`** (spelled with the historical typo baked permanently into the spec) — the URL of the page that linked to this request, used for analytics and — as Chapter 70 warned — a real, ongoing source of accidental credential/token leakage when sensitive data ends up in a URL.

## 10. Status Codes, Grouped by Class

Every response begins with a three-digit status code. The first digit defines the *class* of outcome — a client that has never seen a specific code before can still react sensibly based on its class alone, which is exactly the kind of forward-compatible design HTTP relies on throughout.

```
1xx  Informational  — request received, provisional response, keep going
2xx  Success         — the request was received, understood, and accepted
3xx  Redirection     — further action is needed to complete the request
4xx  Client Error    — the request has a problem the client caused
5xx  Server Error    — the server failed to fulfill a valid request
```

**1xx — Informational**

| Code | Meaning |
|---|---|
| `100 Continue` | Server has received request headers and the client should proceed to send the body (used with `Expect: 100-continue` for large uploads, to avoid sending a huge body the server will reject anyway) |
| `101 Switching Protocols` | Server agrees to switch protocols, e.g. upgrading a connection to a WebSocket (Chapter 76) |

**2xx — Success**

| Code | Meaning |
|---|---|
| `200 OK` | Standard success; body contains the requested representation |
| `201 Created` | A new resource was created (typical response to a successful POST), often with a `Location` header pointing to it |
| `202 Accepted` | Request accepted for processing, but not completed yet (async work) |
| `204 No Content` | Success, but there is deliberately no body (common for a successful DELETE) |

**3xx — Redirection**

| Code | Meaning |
|---|---|
| `301 Moved Permanently` | Resource has a new permanent URL; clients should update bookmarks/links. Cacheable. |
| `302 Found` | Resource is temporarily at a different URL; original URL should still be used in the future |
| `303 See Other` | Fetch the result via GET from a different URL (common after processing a POST, to avoid resubmission on refresh) |
| `304 Not Modified` | The cached copy the client already has is still valid — no body sent at all (Chapter 72 covers this in depth) |
| `307 Temporary Redirect` | Like 302, but explicitly guarantees the method and body are preserved on the redirected request |
| `308 Permanent Redirect` | Like 301, but guarantees method/body preservation |

**4xx — Client Error**

| Code | Meaning |
|---|---|
| `400 Bad Request` | Malformed request syntax the server cannot process |
| `401 Unauthorized` | Authentication is required and missing/invalid (despite the name, this is about *authentication*, not authorization) |
| `403 Forbidden` | Server understood the request but refuses to authorize it (the client's identity is known, but not permitted) |
| `404 Not Found` | No resource exists at this URL |
| `405 Method Not Allowed` | The resource exists, but doesn't support this method |
| `409 Conflict` | Request conflicts with the current state of the resource |
| `429 Too Many Requests` | Client has sent too many requests in a given time (rate limiting) |

**5xx — Server Error**

| Code | Meaning |
|---|---|
| `500 Internal Server Error` | Generic catch-all: something broke on the server, unrelated to what the client sent |
| `502 Bad Gateway` | A server acting as a proxy/gateway got an invalid response from an upstream server |
| `503 Service Unavailable` | Server is temporarily unable to handle the request (overloaded, in maintenance) |
| `504 Gateway Timeout` | A proxy/gateway didn't get a timely response from an upstream server |

**A precise, commonly confused distinction: 401 vs. 403.** `401 Unauthorized` really means "I don't know who you are, or your credentials are invalid — authenticate and try again." `403 Forbidden` means "I know exactly who you are, and you are not allowed to do this — authenticating again won't help." A logged-out user hitting an admin page gets `401`; a logged-in non-admin user hitting the same page gets `403`.

## 11. A Sequence Diagram of One Full Cycle

```mermaid
sequenceDiagram
    participant C as Client
    participant S as Server

    Note over C,S: TCP connection already open (Ch. 59), TLS if applicable (Ch. 82)
    C->>S: GET /products/search?q=laptops HTTP/1.1<br/>Host: shop.example.com<br/>Accept: text/html
    Note right of S: Server parses request line + headers,<br/>routes to application logic,<br/>builds a response
    S-->>C: HTTP/1.1 200 OK<br/>Content-Type: text/html<br/>Content-Length: 1274<br/><br/>&lt;html&gt;...&lt;/html&gt;
    Note over C,S: Connection stays open (Ch. 73) for the next request,<br/>or closes if Connection: close was sent
```

## 12. Deep Dive: How a Server Actually Parses a Request

At the byte level, a server reading from its TCP socket runs roughly this state machine on every request:

```
1. Read bytes until a CRLF is found → that's the request line.
   Split on spaces → method, request-target, HTTP-version.

2. Keep reading lines, each split on the first ':' → header name, header value.
   Stop when a line is empty (just CRLF) → headers are done.

3. Determine if a body is present and how long it is:
   - If Content-Length is present, read exactly that many bytes.
   - If Transfer-Encoding: chunked is present (HTTP/1.1 only, Ch. 73),
     read length-prefixed chunks until a zero-length chunk is seen.
   - If neither is present, there is no body (typical for GET/HEAD/DELETE).

4. Dispatch the (method, path, headers, body) to application logic.

5. Application logic produces a status code, response headers, and a body.

6. Serialize the status line, headers, blank line, and body back onto
   the same TCP connection, in the exact text format from Section 5.
```

This is why `Content-Length` (or chunked encoding) is not a nicety but a structural requirement for anything beyond a single request on a single connection: step 3 is the *only* thing standing between the server correctly framing one message and reading garbage from the start of the next one — this exact framing problem becomes central in Chapter 73's discussion of persistent connections.

## 13. Hands-On: Talking Raw HTTP by Hand

You can speak HTTP/1.1 yourself, by hand, with nothing but a raw TCP tool — this is the clearest possible proof that HTTP really is just structured text over a byte stream:

```
$ nc example.com 80
GET / HTTP/1.1
Host: example.com
Connection: close

HTTP/1.1 200 OK
Content-Type: text/html; charset=UTF-8
Content-Length: 648

<!doctype html>
<html>
...
```

(Type the `GET`, `Host`, and `Connection` lines, then press Enter twice — the blank line is what tells the server "headers are done, this is a complete request.") `openssl s_client` does the equivalent over TLS for an HTTPS site:

```
$ openssl s_client -connect example.com:443 -quiet
GET / HTTP/1.1
Host: example.com
Connection: close

HTTP/1.1 200 OK
...
```

**Inspecting only the response headers, without downloading the body**, using `HEAD` (Section 7):

```
$ curl -I https://example.com/
HTTP/2 200
content-type: text/html; charset=UTF-8
content-length: 648
last-modified: Thu, 17 Oct 2019 07:18:26 GMT
```

**Sending a POST with a JSON body and inspecting both directions:**

```
$ curl -v -X POST https://httpbin.org/post \
  -H "Content-Type: application/json" \
  -d '{"item":"laptop","qty":1}'

> POST /post HTTP/1.1
> Host: httpbin.org
> Content-Type: application/json
> Content-Length: 26
>
< HTTP/1.1 200 OK
< Content-Type: application/json
<
{
  "json": {"item": "laptop", "qty": 1},
  "headers": {"Content-Type": "application/json", ...}
}
```

## 14. Code: A Minimal HTTP Server and Client in Go

```go
// server.go — a server that speaks exactly the text this chapter described
package main

import (
	"fmt"
	"net/http"
)

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Method: %s  Path: %s  Query: %s\n", r.Method, r.URL.Path, r.URL.RawQuery)

	switch r.Method {
	case http.MethodGet:
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, "You asked for: %s\n", r.URL.Path)
	case http.MethodPost:
		w.WriteHeader(http.StatusCreated) // 201
		fmt.Fprintln(w, "Resource created")
	default:
		w.WriteHeader(http.StatusMethodNotAllowed) // 405
	}
}

func main() {
	http.HandleFunc("/", handler)
	http.ListenAndServe(":8080", nil)
}
```

```go
// client.go — a client making both a GET and a POST against it
package main

import (
	"fmt"
	"net/http"
	"strings"
)

func main() {
	resp, _ := http.Get("http://localhost:8080/products/search?q=laptops")
	fmt.Println("GET status:", resp.Status) // "200 OK"
	resp.Body.Close()

	resp2, _ := http.Post("http://localhost:8080/orders", "application/json",
		strings.NewReader(`{"item":"laptop"}`))
	fmt.Println("POST status:", resp2.Status) // "201 Created"
	resp2.Body.Close()
}
```

`http.Get` and `http.Post` in Go's standard library are thin conveniences over exactly the request line, headers, and body format Sections 4–5 described by hand — `net/http` builds the same raw text this chapter has been showing, then writes it to the TCP connection Chapter 70 established.

## 15. Common Misconceptions

- **"GET requests can't have a body."** The HTTP specification does not forbid it, but it explicitly states a GET body "has no defined semantics" — meaning servers, proxies, and caches are free to ignore it, and many implementations (including some HTTP libraries) strip it entirely. In practice: don't rely on a GET body; use query parameters or switch to POST.
- **"A 200 status code always means the operation the user wanted actually succeeded."** It means the *HTTP request itself* was handled successfully — some APIs (a design many consider an anti-pattern) return `200 OK` with a JSON body like `{"error": "insufficient funds"}`, conflating transport-level success with application-level outcome. A well-designed API uses the status code to reflect application outcome too (e.g. `402` or `409`), not just "the server didn't crash."
- **"PUT and POST are interchangeable ways to send data."** They express different semantics (Section 7): PUT means "this is the complete state of a specific, client-named resource," POST means "process this, and you (the server) decide what happens and what identity results." Using the wrong one doesn't break the wire protocol, but it breaks the guarantees (idempotency, Section 8) that infrastructure and other engineers rely on.
- **"Reason phrases (`OK`, `Not Found`) are meaningful to programs."** They exist for humans reading raw traffic. A conformant client must act on the numeric code (`200`, `404`) — the phrase can legally be blank, non-English, or nonstandard, and some servers customize it for style without it changing meaning at all.
- **"Headers are case-sensitive."** Header *names* are explicitly case-insensitive by spec (`Content-Type`, `content-type`, and `CONTENT-TYPE` are identical) — this is why HTTP/2 (Chapter 74) can safely lowercase all header names as part of its binary framing without changing meaning.

## 16. Production Notes

- **Idempotency keys are a real, widely used production pattern.** Payment APIs (Stripe being the most commonly cited example) require clients to attach a unique `Idempotency-Key` header to POST requests specifically so that network retries — which are unavoidable at scale — don't double-charge a customer; the server deduplicates by key and returns the original response for a repeat.
- **`Content-Length` mismatches are a real attack surface.** Request smuggling attacks exploit disagreements between a front-end proxy and a back-end server about where one request ends and the next begins — often by exploiting ambiguity between `Content-Length` and `Transfer-Encoding: chunked` when both are present. This is why modern servers are strict about rejecting requests with conflicting or duplicate framing headers rather than trying to guess intent.
- **Status code discipline matters for automated clients far more than for humans.** A monitoring system, a retry library, or a load balancer's health check reasons about a service almost entirely through status code classes (Section 10) — silently returning `200` for internal errors (to "keep the frontend happy") breaks alerting, breaks automatic retries, and breaks load balancer health checks that were designed around HTTP's actual semantics.
- **`HEAD` is genuinely useful in production, not just an academic curiosity.** CDNs and crawlers use `HEAD` to check whether a resource has changed (via `ETag`/`Last-Modified`, Chapter 72) or to get its size before deciding whether to fetch it — cheaper than a full GET when only the metadata is needed.

## 17. Interview Questions & Model Answers

**Beginner: "What's the difference between PUT and PATCH?"**

"PUT replaces the entire resource at a given URL with the representation you send — it's idempotent, because sending the same full replacement twice leaves the resource in the same state. PATCH applies a partial update, changing only the fields you include, and is not guaranteed to be idempotent depending on how the patch is expressed (e.g., 'increment this counter' is not idempotent even as a PATCH)."

**Intermediate: "Why is it safe to automatically retry a failed GET request but not a failed POST request?"**

"GET is defined as safe and idempotent — it's not supposed to change server state at all, so repeating it any number of times has no additional effect beyond the first. POST is neither safe nor idempotent by default — if a POST to create an order times out after the server actually processed it, blindly retrying could create a second, duplicate order. That's exactly why some APIs introduce idempotency keys: a client-supplied unique ID that lets the server recognize and deduplicate a retried POST safely, giving POST idempotency it doesn't have by default."

**Advanced: "A client gets a 502 Bad Gateway. Walk through what that tells you about where in the request path the failure occurred, and contrast it with a 504 and a plain connection-refused error."**

"A 502 means some intermediary — a reverse proxy, load balancer, or API gateway — successfully received a response from an upstream server, but that response was invalid or malformed, so the intermediary itself is generating the 502 rather than passing through what it got. That tells you the intermediary is up and the upstream connection was made, but the upstream server misbehaved at the protocol level. A 504 Gateway Timeout, by contrast, means the intermediary connected to the upstream but got no response at all within its timeout — the upstream may be alive but slow, hung, or unreachable past the intermediary. A raw connection-refused error, with no HTTP response at all, means the failure happened before any HTTP exchange could even begin — nothing was listening on that port, which points at the target process being down entirely rather than a request-handling problem. In practice, distinguishing these three during an incident tells you whether to look at the proxy layer, the upstream's health, or the upstream process's existence, respectively."

## 18. Exercises

### Easy

1. Write out, by hand, the exact raw HTTP/1.1 request text for a GET to `/api/users/42` on host `api.example.com`, including a `Host` header and an `Accept: application/json` header.
2. Classify each of these status codes by class (1xx-5xx) and state in one sentence what each means: `201`, `304`, `403`, `429`, `503`.
3. For each of GET, POST, PUT, DELETE: state whether it is safe and whether it is idempotent.

### Medium

4. An API endpoint `POST /login` returns `200 OK` with a body of `{"success": false, "error": "bad password"}` on a failed login attempt. Explain what's semantically wrong with this design and propose a better status code.
5. A client sends a request with both `Content-Length: 50` and a body that is actually 80 bytes long. Walk through, using Section 12's parsing state machine, what a strict server should do and why this scenario is relevant to real security vulnerabilities.
6. Explain, using the definitions from Section 8, why `DELETE /orders/42` is idempotent even though calling it a second time returns `404 Not Found` instead of the `204 No Content` the first call returned.

### Hard

7. Design (in words) an HTTP method-routing table for a REST-style API managing a `/articles` resource: which method + path combinations should exist for listing, creating, reading one, fully replacing one, partially updating one, and deleting one. Justify each method choice against Section 7's semantics.
8. A proxy in front of your API interprets `Transfer-Encoding: chunked` but your backend only understands `Content-Length`, and an attacker sends a request with both headers set to conflicting values. Explain, at a mechanical level, how this discrepancy could let an attacker "smuggle" a second hidden request past the proxy's inspection and directly to the backend.
9. You're designing a payment API's `POST /charges` endpoint. Explain exactly how an `Idempotency-Key` header should change the server's behavior on a retried request with the identical key, and what the server must store to make that guarantee — tying your answer to why POST is neither safe nor idempotent by default.

## 19. Summary

| Term | Meaning |
|---|---|
| Request line | `METHOD path HTTP-version` — the first line of an HTTP request |
| Status line | `HTTP-version status-code reason-phrase` — the first line of an HTTP response |
| Header | A `Name: value` metadata line describing the message, sender, or body |
| Method | The verb of a request (GET, POST, PUT, PATCH, DELETE, HEAD, OPTIONS) |
| Safe | Method causes no server-side side effects (GET, HEAD, OPTIONS) |
| Idempotent | Repeating the request has the same effect as doing it once (GET, PUT, DELETE — not POST/PATCH) |
| Content-Length | Header giving the exact byte length of the body, used to frame messages |
| 1xx/2xx/3xx/4xx/5xx | Informational / Success / Redirection / Client Error / Server Error status classes |
| 401 vs 403 | "I don't know who you are" vs. "I know who you are and you can't do this" |
| Idempotency key | Client-generated unique ID that makes a non-idempotent POST safe to retry |

HTTP's request/response cycle is stateless by design — every request in Section 6's capture stands entirely alone, with no memory of any request before it. Chapter 72 confronts exactly why that statelessness becomes a problem the moment you need a login or a shopping cart, and shows the mechanism (cookies) the Web built to add state back on top of a protocol that was never designed to have any.
