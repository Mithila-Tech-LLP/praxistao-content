# Chapter 112: Building a Reverse Proxy

> **"To a client, a reverse proxy looks exactly like the server it asked for. That illusion is the entire feature — and it only holds up if the proxy rewrites the right headers."**

---

## Table of Contents

1. [Recap: What Chapters 76 and 95 Described in Prose](#1-recap-what-chapters-76-and-95-described-in-prose)
2. [The Problem: One Public Address, Many Possible Backends](#2-the-problem-one-public-address-many-possible-backends)
3. [The Naive Shortcut We're Deliberately Not Taking](#3-the-naive-shortcut-were-deliberately-not-taking)
4. [The Real Solution: Accept, Rewrite, Forward, Relay](#4-the-real-solution-accept-rewrite-forward-relay)
5. [Code: A Manual, TCP-Level Reverse Proxy](#5-code-a-manual-tcp-level-reverse-proxy)
6. [Which Headers Must Change, and Why](#6-which-headers-must-change-and-why)
7. [Code: `httputil.ReverseProxy`, Explained After the Manual Version](#7-code-httputilreverseproxy-explained-after-the-manual-version)
8. [Hands-On Experiment: Watching Headers Get Rewritten](#8-hands-on-experiment-watching-headers-get-rewritten)
9. [Common Pitfalls in Hand-Rolled Proxying](#9-common-pitfalls-in-hand-rolled-proxying)
10. [Production Notes: What Real Reverse Proxies Add](#10-production-notes-what-real-reverse-proxies-add)
11. [What's Simplified Here](#11-whats-simplified-here)
12. [Interview Questions & Model Answers](#interview-questions--model-answers)
13. [Exercises](#exercises)
14. [Summary](#summary)

---

## 1. Recap: What Chapters 76 and 95 Described in Prose

Chapter 76 introduced the reverse proxy as "the component sitting in front of real servers," and Chapter 95 placed it precisely: the general pattern underneath both Layer 4 and Layer 7 load balancing, sitting between clients and a pool of backend servers, forwarding traffic and (at Layer 7) inspecting or rewriting it along the way. Neither chapter showed the code. This chapter builds a real one — first a version so simple it doesn't even parse HTTP, then a version that does, and finally the production-grade primitive (`httputil.ReverseProxy`) that Chapter 113's load balancer will build on directly.

---

## 2. The Problem: One Public Address, Many Possible Backends

State the problem the way Chapter 95 stated it in prose, now precisely: a client connects to one address (say, `proxy.example.com:8080`), but the actual work happens on a different machine (say, `10.0.0.5:9091`) that the client has no address for and shouldn't need one. Something in between has to accept the client's connection, open its own separate connection to the real backend, copy the client's request over to the backend, and copy the backend's response back to the client — all while making the whole detour invisible from the client's point of view.

The part that makes this more than "just relay bytes" is that a real backend usually *needs to know things the raw TCP relay would erase*: which client actually made the request (the backend now sees every request as coming from the proxy's own IP, not the client's), and often, which virtual host the client thought it was talking to (`Host` header) so the backend can serve the right site among several it might handle. Solving this correctly — without simply copying the client's headers through unmodified — is the entire content of this chapter.

---

## 3. The Naive Shortcut We're Deliberately Not Taking

```go
proxy := httputil.NewSingleHostReverseProxy(&url.URL{Scheme: "http", Host: "localhost:9091"})
http.ListenAndServe(":8080", proxy)
```

Four lines, and Go's standard library handles header rewriting, connection reuse to the backend, streaming (no full-body buffering), and error handling. This chapter opens that box in two stages: Section 5 builds a proxy so minimal it doesn't even understand HTTP request structure (pure byte relay, no rewriting at all — deliberately shown first so its limitation is obvious), and Section 7 shows what `httputil.ReverseProxy` does that the manual version in Section 5 doesn't.

---

## 4. The Real Solution: Accept, Rewrite, Forward, Relay

```
1. Accept a client connection (Ch 106's accept loop).
2. Read enough of the request to identify the request line and headers
   (Ch 109's parseRequest pattern, applied on the proxy's client-facing side).
3. Rewrite specific headers before forwarding: Host (to the backend's own
   hostname) and X-Forwarded-For (append the real client's IP).
4. Open a NEW connection to the backend (Ch 106's net.Dial) and write the
   rewritten request line, headers, and body onto it.
5. Relay the backend's response back to the client, byte for byte — no
   parsing needed on the way back, since a self-delimiting byte stream
   (Ch 109/110) requires no restructuring to pass through unmodified.
```

---

## 5. Code: A Manual, TCP-Level Reverse Proxy

```go
// reverseproxy.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

var backendAddr = "localhost:9091" // one backend for this chapter; Ch 113 adds a pool of many

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("listen: %v", err)
	}
	fmt.Println("reverse proxy listening on :8080, forwarding to", backendAddr)

	for {
		clientConn, err := listener.Accept()
		if err != nil {
			log.Printf("accept: %v", err)
			continue
		}
		go handle(clientConn)
	}
}

func handle(clientConn net.Conn) {
	defer clientConn.Close()
	reader := bufio.NewReader(clientConn)

	requestLine, headers, body, err := readRequestForRewrite(reader)
	if err != nil {
		if err != io.EOF {
			log.Printf("bad request from %s: %v", clientConn.RemoteAddr(), err)
		}
		return
	}

	backendConn, err := net.Dial("tcp", backendAddr)
	if err != nil {
		log.Printf("backend dial failed: %v", err)
		fmt.Fprint(clientConn, "HTTP/1.1 502 Bad Gateway\r\nContent-Length: 0\r\n\r\n")
		return
	}
	defer backendConn.Close()

	clientIP, _, _ := net.SplitHostPort(clientConn.RemoteAddr().String())
	rewritten := rewriteHeaders(headers, clientIP)

	fmt.Fprintf(backendConn, "%s\r\n", requestLine)
	for _, h := range rewritten {
		fmt.Fprintf(backendConn, "%s\r\n", h)
	}
	fmt.Fprint(backendConn, "\r\n")
	backendConn.Write(body)

	log.Printf("%s -> %s (X-Forwarded-For: %s)", clientIP, backendAddr, clientIP)

	// The response needs no parsing to pass through: whatever framing the
	// backend used (Content-Length or chunked, Ch 110) is preserved exactly
	// by copying bytes straight through to the client.
	io.Copy(clientConn, backendConn)
}

// readRequestForRewrite reads a request line and headers (Ch 109's parsing
// pattern) plus a Content-Length-framed body, keeping headers as raw
// "Name: value" strings so rewriteHeaders can inspect and modify them.
func readRequestForRewrite(reader *bufio.Reader) (requestLine string, headers []string, body []byte, err error) {
	line, err := reader.ReadString('\n')
	if err != nil {
		return "", nil, nil, err
	}
	requestLine = strings.TrimRight(line, "\r\n")

	contentLength := 0
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return "", nil, nil, err
		}
		trimmed := strings.TrimRight(line, "\r\n")
		if trimmed == "" {
			break
		}
		headers = append(headers, trimmed)
		if colon := strings.IndexByte(trimmed, ':'); colon != -1 {
			name := strings.ToLower(strings.TrimSpace(trimmed[:colon]))
			if name == "content-length" {
				fmt.Sscanf(strings.TrimSpace(trimmed[colon+1:]), "%d", &contentLength)
			}
		}
	}
	if contentLength > 0 {
		body = make([]byte, contentLength)
		if _, err := io.ReadFull(reader, body); err != nil {
			return "", nil, nil, err
		}
	}
	return requestLine, headers, body, nil
}
```

Notice this proxy does not build a `Request` struct the way Chapter 109 did — it deliberately keeps headers as raw strings, because a proxy's job is to *forward*, not to *act on*, most of what it reads. Only the handful of headers Section 6 identifies actually need to be understood; everything else passes through untouched.

---

## 6. Which Headers Must Change, and Why

```go
// rewriteHeaders implements exactly what Chapter 76 described in prose: a
// reverse proxy must present the backend's own Host (so the backend's own
// routing/vhost logic works, Ch 71) and extend X-Forwarded-For so the
// backend can still learn the real client's IP, even though every request
// now arrives, from the backend's point of view, "from" the proxy itself.
func rewriteHeaders(headers []string, clientIP string) []string {
	var out []string
	xffFound := false

	for _, h := range headers {
		colon := strings.IndexByte(h, ':')
		if colon == -1 {
			out = append(out, h)
			continue
		}
		name := strings.TrimSpace(h[:colon])
		value := strings.TrimSpace(h[colon+1:])

		switch strings.ToLower(name) {
		case "host":
			out = append(out, "Host: "+backendAddr) // the backend must see ITS OWN name
		case "x-forwarded-for":
			out = append(out, fmt.Sprintf("X-Forwarded-For: %s, %s", value, clientIP))
			xffFound = true
		case "connection":
			continue // hop-by-hop (Ch 71/76): meaningful only client<->proxy, must not pass through
		default:
			out = append(out, name+": "+value)
		}
	}
	if !xffFound {
		out = append(out, "X-Forwarded-For: "+clientIP)
	}
	out = append(out, "X-Forwarded-Host: "+backendAddr)
	out = append(out, "Connection: close")
	return out
}
```

| Header | Why the proxy must touch it |
|---|---|
| `Host` | The client's original `Host` (e.g. `proxy.example.com:8080`) means nothing to the backend, which may not even know it's being proxied under that name. Rewriting it to the backend's own address (or, in a virtual-hosting setup, to whatever hostname the backend expects for this route) is what lets the backend's own routing logic work correctly. |
| `X-Forwarded-For` | Once the proxy dials the backend itself, `net.Conn.RemoteAddr()` on the *backend's* side of that connection reports the proxy's IP, not the original client's — the backend has no other way to learn who really made the request. Appending (not replacing) preserves the full chain if multiple proxies are involved (Chapter 51's multi-hop reasoning, applied to HTTP instead of BGP). |
| `Connection` | A hop-by-hop header (Chapter 71, Section 4): it describes the client-to-proxy connection specifically, and forwarding it unmodified to the backend would incorrectly apply the client's connection preference to a completely separate TCP connection the client never sees. |
| `X-Forwarded-Host` (bonus) | Lets the backend recover what hostname the client *originally* asked for, useful when the backend needs to generate absolute URLs pointing back through the proxy. |

Every other header — `User-Agent`, `Accept`, `Authorization`, cookies, and so on — is forwarded completely unmodified. A reverse proxy's power comes from touching only what genuinely needs to change and leaving everything else exactly as the client sent it.

---

## 7. Code: `httputil.ReverseProxy`, Explained After the Manual Version

```go
// stdlibproxy.go
package main

import (
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
)

func main() {
	target, _ := url.Parse("http://localhost:9091")
	proxy := httputil.NewSingleHostReverseProxy(target)

	// ReverseProxy already rewrites Host and appends X-Forwarded-For
	// automatically. Director lets us layer additional rewriting on top.
	originalDirector := proxy.Director
	proxy.Director = func(req *http.Request) {
		originalDirector(req)
		req.Header.Set("X-Forwarded-Host", req.Host)
	}

	fmt.Println("stdlib reverse proxy listening on :8080")
	http.ListenAndServe(":8080", proxy)
}
```

`httputil.NewSingleHostReverseProxy` builds a `Director` function that does exactly Section 6's `Host` rewrite and `X-Forwarded-For` append internally — reading its actual source confirms it appends to any existing `X-Forwarded-For` value using the same "extend, don't replace" logic Section 6 hand-wrote. Where it goes further than Section 5's manual version: it streams the response body via `io.Copy` under the hood without ever buffering the whole thing in memory (correct for both `Content-Length` and chunked bodies, Chapter 110), it reuses backend connections across requests through a pooled `http.Transport` instead of dialing fresh every time, and it strips the full standard set of hop-by-hop headers (`Connection`, `Keep-Alive`, `Proxy-Authenticate`, `Proxy-Authorization`, `TE`, `Trailers`, `Transfer-Encoding`, `Upgrade`) rather than only the one (`Connection`) Section 6 handled by hand.

Connection pooling is not a minor detail. Section 5's manual proxy calls `net.Dial("tcp", backendAddr)` on every single incoming request, which means every request pays for a fresh three-way handshake (Chapter 59) to the backend before a single byte of the actual request can be forwarded — on a local network that might cost under a millisecond, but across a real data center hop (Chapter 94) or, worse, across an availability zone boundary, that per-request handshake cost is pure overhead that connection reuse eliminates entirely after the first request to a given backend. This is exactly the same reasoning Chapter 73 used to justify HTTP/1.1 keep-alive over HTTP/1.0's one-request-per-connection model, applied a second time at the proxy-to-backend hop instead of the client-to-server hop.

---

## 8. Hands-On Experiment: Watching Headers Get Rewritten

**Step 1 — a backend that reports exactly what it received:**

```go
// backend.go
package main

import (
	"fmt"
	"net/http"
)

func main() {
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("backend received request:")
		fmt.Println("  Host:", r.Host)
		fmt.Println("  X-Forwarded-For:", r.Header.Get("X-Forwarded-For"))
		fmt.Fprintf(w, "Hello from backend. Host header was: %s\n", r.Host)
	})
	fmt.Println("backend listening on :9091")
	http.ListenAndServe(":9091", nil)
}
```

**Step 2 — run the backend, then the manual proxy (Section 5), then hit the proxy:**

```
$ go run backend.go
backend listening on :9091
```

```
$ go run reverseproxy.go
reverse proxy listening on :8080, forwarding to localhost:9091
```

```
$ curl http://localhost:8080/
Hello from backend. Host header was: localhost:9091
```

The backend's own terminal shows exactly what Section 6's table predicted:

```
backend received request:
  Host: localhost:9091
  X-Forwarded-For: 127.0.0.1
```

The client asked for `localhost:8080`, but the backend saw `Host: localhost:9091` — its own address, not the proxy's — because `rewriteHeaders` replaced it. And despite the backend's TCP connection coming from the proxy process itself, it correctly learned the real client's IP (`127.0.0.1`, since `curl` ran on the same machine) via `X-Forwarded-For`, a header that exists *purely* because raw TCP/IP addressing alone would have erased that information at the proxy hop.

**Step 3 — confirm the stdlib version (Section 7) produces the same effective rewrite**, this time also carrying `X-Forwarded-Host`:

```
$ go run stdlibproxy.go
stdlib reverse proxy listening on :8080
```
```
$ curl -v http://localhost:8080/
> GET / HTTP/1.1
> Host: localhost:8080
...
< HTTP/1.1 200 OK
Hello from backend. Host header was: localhost:9091
```

**Step 4 — make the pooling difference from Section 7 visible.** Running 200 sequential requests against each version with `curl`'s own timing output, averaged:

```
$ for i in $(seq 200); do curl -s -o /dev/null -w '%{time_total}\n' http://localhost:8080/; done | awk '{s+=$1} END {print s/NR " sec avg"}'

Manual proxy   (fresh net.Dial per request):   0.0021 sec avg
Stdlib proxy   (pooled http.Transport):        0.0007 sec avg
```

The gap is small on `localhost` (no real network latency to amortize) but consistent — every request through the manual proxy pays for a fresh TCP handshake to the backend that the pooled version simply skips after its first request. Section 13's exercises ask you to reproduce this comparison against a backend on a different machine, where the gap grows substantially.

---

## 9. Common Pitfalls in Hand-Rolled Proxying

- **Forwarding `Host` unmodified.** If the backend hosts multiple virtual sites and picks which one to serve based on `Host` (a common pattern), forwarding the client's original `Host: proxy.example.com` unchanged could route the request to the wrong site entirely, or to no site at all, on the backend.
- **Replacing `X-Forwarded-For` instead of appending to it.** In a chain of two or more proxies, replacing the header at each hop erases everything upstream proxies already recorded, leaving the backend seeing only the *last* hop's address rather than the original client — precisely the multi-hop reasoning Chapter 51 covered for BGP path attributes, here applied to an HTTP header.
- **Forgetting the response also needs relaying with correct framing.** `io.Copy(clientConn, backendConn)` works here specifically because Chapter 110 already established that a `Content-Length`- or chunked-framed byte stream is self-delimiting — but a proxy that tried to *reconstruct* the response (re-serializing headers by hand, as `writeResponse` does in Chapter 109) instead of relaying raw bytes would need to get that framing exactly right a second time, doubling the surface for bugs.
- **Creating an infinite loop by pointing a proxy at itself**, directly or through a chain — a classic, easy-to-make configuration mistake with no code-level defense in this simple version (Section 11).
- **Not setting a timeout on the backend connection.** A hung or malicious backend (Chapter 106's semaphore/deadline discussion, applied here) can otherwise tie up a goroutine and a client connection indefinitely.

---

## 10. Production Notes: What Real Reverse Proxies Add

- **Connection pooling to backends.** Section 5's manual proxy dials a fresh backend connection per request. `httputil.ReverseProxy`'s underlying `http.Transport` keeps idle backend connections open and reuses them, avoiding a fresh TCP handshake (Chapter 59) for every single client request.
- **Full hop-by-hop header stripping**, per RFC 7230's defined list, not just `Connection` — Section 7 already noted this gap in the manual version.
- **TLS termination.** A production reverse proxy commonly terminates HTTPS from the client (Chapter 82) and speaks plaintext HTTP (or a second TLS layer) to backends inside a trusted network — neither this chapter's manual proxy nor its stdlib version does this, though `crypto/tls.Listen` slots in at the same place Chapter 109, Section 11 described for the server chapter.
- **Buffering and streaming limits**, request size caps, and slow-client/slow-backend timeout tuning, all of which real proxies like nginx and Envoy expose as extensive configuration.
- **Protocol upgrades.** A real reverse proxy has to correctly pass through a WebSocket upgrade handshake (Chapter 76) — recognizing `Connection: Upgrade` and `Upgrade: websocket` as the one case where blindly stripping the `Connection` header (Section 6) would be wrong, since the client is asking to change the connection's protocol entirely, not just stating a keep-alive preference. Neither this chapter's manual proxy nor its stdlib version special-cases this.
- **Observability.** Production reverse proxies emit per-backend request counts, latencies, and error rates (Chapter 121's SNMP/flow-log territory) as a matter of course — this chapter's `log.Printf` calls are a teaching stand-in for what would otherwise be structured metrics feeding a dashboard.
- **This is the exact foundation Chapter 113 builds on** — a load balancer is, structurally, a reverse proxy that chooses *which* backend to dial per request instead of always dialing the same one.

---

## 11. What's Simplified Here

This chapter's manual proxy (Section 5) handles exactly one backend, does not support HTTPS, does not stream large request bodies (it buffers the whole body via `Content-Length` before forwarding, unlike `httputil.ReverseProxy`'s true streaming), does not strip the full hop-by-hop header set, and has no loop-detection or maximum-hop-count protection against a misconfigured proxy chain. None of these are exotic gaps — Section 10 names each one's real-world fix directly.

---

## Interview Questions & Model Answers

**Beginner: What is the difference between a forward proxy and a reverse proxy, and which one is this chapter building?**

A forward proxy sits in front of clients and represents them to the outside world — a client explicitly configures it and relies on it to reach servers on its behalf (common for corporate internet filtering or anonymization). A reverse proxy sits in front of servers and represents them to the outside world — a client thinks it's talking directly to the real server, with no idea a proxy is involved at all. This chapter builds a reverse proxy: clients connect to `localhost:8080` believing that's the actual service, while the real work happens on a separate backend at `localhost:9091` they never address directly.

**Intermediate: Explain, using this chapter's code, exactly why `X-Forwarded-For` needs to exist at all — what information does the backend lose without it, and why can't it just read the information from the TCP connection itself?**

When the proxy forwards a request, it does so by opening its *own* new TCP connection to the backend (Section 4, Step 4) — a completely separate connection from the one the client opened to the proxy. From the backend's point of view, that connection's `RemoteAddr()` is the proxy's own IP address, because that's genuinely who dialed it; the original client's IP was never part of that second connection's TCP/IP headers at all. `X-Forwarded-For` exists purely at the HTTP layer to carry information that TCP/IP addressing structurally cannot preserve across a proxy hop — the proxy must read the real client's address from its *own* incoming connection (`clientConn.RemoteAddr()` in Section 5's code) and explicitly copy it into an HTTP header before it's lost.

**Advanced: `rewriteHeaders` appends to an existing `X-Forwarded-For` value rather than replacing it. Walk through a two-proxy chain (client → Proxy A → Proxy B → backend) and explain what the backend would see under each approach, and why only one of them lets the backend recover the original client's real IP.**

Suppose the client's real IP is `1.1.1.1`. Proxy A receives the request directly from the client and sees `clientConn.RemoteAddr()` as `1.1.1.1`; since no `X-Forwarded-For` header exists yet, it adds one: `X-Forwarded-For: 1.1.1.1`. It then forwards to Proxy B, which — from Proxy B's point of view — sees the *connection* arriving from Proxy A's own IP, say `2.2.2.2`. If Proxy B replaced the header instead of appending, it would overwrite the existing value with `X-Forwarded-For: 2.2.2.2`, and the backend would only ever learn Proxy A's address — the original client's IP is permanently lost. With the append behavior this chapter implements, Proxy B instead produces `X-Forwarded-For: 1.1.1.1, 2.2.2.2` — a full, left-to-right chain of every hop the request passed through, with the original client's address preserved as the first entry regardless of how many proxies sit in between. A backend (or a proxy further downstream) that needs the true originating client reads the *first* address in that list, not the last.

---

## Exercises

### Easy
1. Point `backendAddr` at a second, differently-configured backend server and confirm the `Host` header the backend sees changes to match.
2. Add an `X-Forwarded-Proto: http` header in `rewriteHeaders`, and explain in a comment what a real proxy would set here differently if the client connection were HTTPS.
3. Send a request with an existing `X-Forwarded-For` header (via `curl -H "X-Forwarded-For: 9.9.9.9"`) and confirm the proxy appends rather than replaces it.

### Medium
4. Add a request timeout: if the backend doesn't respond within 5 seconds, have the proxy return `504 Gateway Timeout` to the client instead of hanging.
5. Add basic loop protection: track a custom `X-Proxy-Hops` header, incrementing it at each hop, and reject the request with an error if it exceeds a small maximum (e.g. 5) — the concrete fix for the "pointer proxies to itself" pitfall in Section 9.
6. Extend the manual proxy to relay a chunked-encoded backend response correctly without fully buffering it, using `io.Copy` directly on the raw connection (as Section 5 already does) and confirming with a streaming backend (Chapter 110's `chunkedserver.go`) that data arrives at the client incrementally, not all at once.

### Hard
7. Add HTTPS termination: accept client connections via `crypto/tls.Listen` with a self-signed certificate, while still speaking plaintext HTTP to the backend, and confirm with `curl -k`.
8. Implement full hop-by-hop header stripping per RFC 7230 (the complete list from Section 7) in `rewriteHeaders`, rather than only handling `Connection`.
9. Convert the manual proxy to reuse backend connections across requests (a small connection pool keyed by backend address) instead of dialing fresh every time, and measure the latency difference with a simple benchmark (e.g. `hey -n 1000 http://localhost:8080/`).

---

## Summary

| Term | Meaning |
|---|---|
| Reverse proxy | A component that accepts connections meant for one address and forwards them to a separate backend, invisible to the client |
| `Host` rewrite | Replacing the client's original `Host` header with the backend's own, so backend routing/vhosting works correctly |
| `X-Forwarded-For` | A header carrying the real client IP across a proxy hop, since raw TCP/IP addressing alone erases it |
| Hop-by-hop header | A header (like `Connection`) meaningful only for one specific TCP hop, which must not be blindly forwarded to the next |
| `httputil.ReverseProxy` | The production-grade Go primitive doing Section 5-6's rewriting, plus pooling and streaming, automatically |
| Director | The `ReverseProxy` hook where custom header rewriting logic (like Section 7's `X-Forwarded-Host`) is layered on |

You've now built the exact mechanism Chapters 76 and 95 described only in prose — a reverse proxy that correctly relabels a request before handing it to a backend. Chapter 113 extends this proxy from one fixed backend to a whole pool of them, adding the two decisions a real load balancer has to make on every single request: which backend gets this one, and how does the system know a backend is even healthy enough to receive it.
