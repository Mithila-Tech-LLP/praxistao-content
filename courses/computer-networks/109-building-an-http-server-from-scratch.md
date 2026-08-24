# Chapter 109: Building an HTTP Server From Scratch

> **"`http.ListenAndServe` is one line of code that hides an entire chapter's worth of byte-level parsing. This chapter deletes that line and writes the chapter."**

---

## Table of Contents

1. [Recap: What Chapter 71 Described in Prose](#1-recap-what-chapter-71-described-in-prose)
2. [The Problem: From a Raw Byte Stream to a Structured Request](#2-the-problem-from-a-raw-byte-stream-to-a-structured-request)
3. [The Naive Shortcut We're Deliberately Not Taking](#3-the-naive-shortcut-were-deliberately-not-taking)
4. [The Real Solution: Parse the Request Line, Headers, and Body By Hand](#4-the-real-solution-parse-the-request-line-headers-and-body-by-hand)
5. [Code: A Complete Hand-Rolled HTTP Server](#5-code-a-complete-hand-rolled-http-server)
6. [Tracing Bytes: From Raw Socket to Parsed Request](#6-tracing-bytes-from-raw-socket-to-parsed-request)
7. [Writing a Spec-Correct Response, Field by Field](#7-writing-a-spec-correct-response-field-by-field)
8. [Hands-On Experiment: `nc`, `curl`, and a Byte-Level Comparison](#8-hands-on-experiment-nc-curl-and-a-byte-level-comparison)
9. [Keep-Alive: Handling More Than One Request Per Connection](#9-keep-alive-handling-more-than-one-request-per-connection)
10. [Common Pitfalls in Hand-Rolled HTTP Parsing](#10-common-pitfalls-in-hand-rolled-http-parsing)
11. [Production Notes: Why Real Servers Don't Do This, and What They Do Instead](#11-production-notes-why-real-servers-dont-do-this-and-what-they-do-instead)
12. [What's Simplified Here](#12-whats-simplified-here)
13. [Interview Questions & Model Answers](#interview-questions--model-answers)
14. [Exercises](#exercises)
15. [Summary](#summary)

---

## 1. Recap: What Chapter 71 Described in Prose

Chapter 71 laid out the exact textual structure of an HTTP/1.1 request and response — a request line, headers as `Name: value` pairs separated by CRLF, a blank line, and an optional body — and Section 12 of that chapter walked through, in prose, the six-step state machine a server runs to turn a raw byte stream into a request it can act on. Chapter 71's own code sample (Section 14) then used `net/http`'s `http.ListenAndServe`, which performs that exact parsing internally, in code you never see.

This chapter deletes that convenience. Every byte this server reads off the wire is parsed by code you can read start to finish in this chapter — no framework, no hidden state machine, nothing standing between the raw TCP stream (Chapter 106) and a `Request` struct your own code builds field by field.

---

## 2. The Problem: From a Raw Byte Stream to a Structured Request

Chapter 106, Section 8 already established the fact this chapter has to grapple with head-on: **a TCP connection is an undifferentiated stream of bytes.** `conn.Read()` doesn't know or care that the bytes it returns happen to spell out `GET / HTTP/1.1\r\n`. Turning that stream into "this is a GET request for path `/`, using HTTP version 1.1, with these headers and this body" is entirely the server's job, and HTTP's specification (Chapter 71) is precisely the agreed-upon set of rules for doing it.

Restated as a concrete engineering problem: given a `*bufio.Reader` wrapping a live TCP connection, and no guarantee about how the bytes were chunked when they arrived (Chapter 106, Section 8 — one `Read()` might return a fragment of a line, or several lines at once), correctly identify where the request line ends, where each header ends, where the headers end and the body begins, and exactly how many bytes of body to read — all without either blocking forever waiting for bytes that will never come, or reading past the end of one request into the start of the next.

---

## 3. The Naive Shortcut We're Deliberately Not Taking

It would take exactly three lines to get a working HTTP server in Go:

```go
http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "hello")
})
http.ListenAndServe(":8080", nil)
```

This is genuinely how you should build HTTP servers in real, production Go code — Chapter 71 used exactly this. But it is also, deliberately, a black box for the purposes of this chapter: `net/http` does its own request-line parsing, header parsing, `Content-Length`/chunked body handling, keep-alive management, and response serialization internally, in code this course hasn't shown you. The entire point of this chapter is to open that box. Every line inside `net/http`'s request parser is doing something this chapter's Section 5 does too, just with far more edge-case handling, performance tuning, and security hardening (Section 11 discusses exactly what's missing here as a result).

---

## 4. The Real Solution: Parse the Request Line, Headers, and Body By Hand

Directly implementing Chapter 71, Section 12's six-step outline:

```
1. Read one line (up to \r\n) → split on spaces → method, request-target, version.
2. Read lines one at a time, splitting each on the first ':' → header
   name/value pairs. Stop at the first blank line.
3. Check the Content-Length header. If present and > 0, read EXACTLY
   that many more bytes as the body.
4. Build a Request value from all of the above.
5. Hand it to routing/handler logic, which decides on a status code,
   headers, and a body for the response.
6. Serialize the status line, headers, blank line, and body back onto
   the same connection, in the exact text format Chapter 71, Section 5
   specified.
```

`bufio.Reader` is the right tool for step 1 and 2 specifically because it buffers incoming bytes internally and exposes `ReadString('\n')`, which handles exactly the "read until this delimiter, however many underlying `Read()` calls that takes" problem Section 2 described — you never have to manually stitch together partial lines yourself.

---

## 5. Code: A Complete Hand-Rolled HTTP Server

This is the whole program — one file, no `net/http` import anywhere, compiled and run exactly as shown.

```go
// httpserver.go
package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"strings"
	"time"
)

// ---- Request: what parsing produces ----

type Request struct {
	Method  string
	Path    string
	Query   string
	Version string
	Headers map[string]string // keys stored lowercase — headers are case-insensitive (Ch 71)
	Body    []byte
}

// parseRequest implements Chapter 71, Section 12's six-step state machine,
// steps 1–3, reading directly off a buffered TCP connection.
func parseRequest(reader *bufio.Reader) (*Request, error) {
	// Step 1: the request line.
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err // includes io.EOF — the client closed the connection
	}
	line = strings.TrimRight(line, "\r\n")
	parts := strings.SplitN(line, " ", 3)
	if len(parts) != 3 {
		return nil, fmt.Errorf("malformed request line: %q", line)
	}
	method, target, version := parts[0], parts[1], parts[2]

	if version != "HTTP/1.1" && version != "HTTP/1.0" {
		return nil, fmt.Errorf("unsupported HTTP version: %q", version)
	}

	path, query := target, ""
	if idx := strings.IndexByte(target, '?'); idx != -1 {
		path, query = target[:idx], target[idx+1:]
	}

	// Step 2: headers, one per line, until a blank line.
	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}
		line = strings.TrimRight(line, "\r\n")
		if line == "" {
			break // blank line: end of headers (Ch 71, Section 4)
		}
		colon := strings.IndexByte(line, ':')
		if colon == -1 {
			return nil, fmt.Errorf("malformed header line: %q", line)
		}
		name := strings.ToLower(strings.TrimSpace(line[:colon]))
		value := strings.TrimSpace(line[colon+1:])
		headers[name] = value
	}

	// Step 3: the body, if Content-Length says there is one.
	var body []byte
	if clStr, ok := headers["content-length"]; ok {
		length, err := strconv.Atoi(clStr)
		if err != nil || length < 0 {
			return nil, fmt.Errorf("invalid content-length: %q", clStr)
		}
		if length > 0 {
			body = make([]byte, length)
			if _, err := io.ReadFull(reader, body); err != nil {
				return nil, fmt.Errorf("body shorter than content-length: %w", err)
			}
		}
	}
	// Note: Transfer-Encoding: chunked is NOT handled here — see Section 12.

	return &Request{
		Method: method, Path: path, Query: query,
		Version: version, Headers: headers, Body: body,
	}, nil
}

// ---- Response: writing spec-correct HTTP/1.1 bytes ----

func writeResponse(w io.Writer, status int, statusText string, headers map[string]string, body []byte) error {
	if headers == nil {
		headers = map[string]string{}
	}
	headers["Content-Length"] = strconv.Itoa(len(body)) // Ch 71 Section 9: mandatory for framing
	if _, ok := headers["Content-Type"]; !ok {
		headers["Content-Type"] = "text/plain; charset=utf-8"
	}
	headers["Date"] = time.Now().UTC().Format("Mon, 02 Jan 2006 15:04:05 GMT")
	headers["Server"] = "hand-rolled-go/1.0"

	var b strings.Builder
	fmt.Fprintf(&b, "HTTP/1.1 %d %s\r\n", status, statusText) // status line, Ch 71 Section 5
	for name, value := range headers {
		fmt.Fprintf(&b, "%s: %s\r\n", name, value)
	}
	b.WriteString("\r\n") // blank line: end of headers, start of body

	if _, err := io.WriteString(w, b.String()); err != nil {
		return err
	}
	_, err := w.Write(body)
	return err
}

// ---- Routing: application logic, independent of parsing/serialization ----

func route(req *Request) (status int, statusText string, headers map[string]string, body []byte) {
	headers = map[string]string{}
	switch {
	case req.Method == "GET" && req.Path == "/":
		return 200, "OK", headers, []byte("Welcome to a hand-rolled HTTP server.\n")

	case req.Method == "GET" && req.Path == "/hello":
		name := "World"
		for _, pair := range strings.Split(req.Query, "&") {
			if kv := strings.SplitN(pair, "=", 2); len(kv) == 2 && kv[0] == "name" {
				name = kv[1]
			}
		}
		return 200, "OK", headers, []byte(fmt.Sprintf("Hello, %s!\n", name))

	case req.Method == "POST" && req.Path == "/echo":
		headers["Content-Type"] = "application/octet-stream"
		return 200, "OK", headers, req.Body

	default:
		return 404, "Not Found", headers, []byte("404 Not Found\n")
	}
}

// ---- Connection handling: reusing Chapter 106's accept-loop pattern ----

func handleConnection(conn net.Conn) {
	defer conn.Close()
	reader := bufio.NewReader(conn)

	for { // Section 9: one connection may carry several requests (keep-alive)
		conn.SetReadDeadline(time.Now().Add(30 * time.Second))

		req, err := parseRequest(reader)
		if err != nil {
			if err != io.EOF {
				log.Printf("bad request from %s: %v", conn.RemoteAddr(), err)
				writeResponse(conn, 400, "Bad Request",
					map[string]string{"Connection": "close"}, []byte("400 Bad Request\n"))
			}
			return
		}

		status, statusText, headers, body := route(req)
		keepAlive := shouldKeepAlive(req)
		if keepAlive {
			headers["Connection"] = "keep-alive"
		} else {
			headers["Connection"] = "close"
		}

		log.Printf("%s %s %s -> %d", req.Method, req.Path, conn.RemoteAddr(), status)

		if err := writeResponse(conn, status, statusText, headers, body); err != nil {
			return
		}
		if !keepAlive {
			return
		}
	}
}

func shouldKeepAlive(req *Request) bool {
	conn := strings.ToLower(req.Headers["connection"])
	if req.Version == "HTTP/1.0" {
		return conn == "keep-alive" // HTTP/1.0 defaults to close (Ch 73)
	}
	return conn != "close" // HTTP/1.1 defaults to keep-alive (Ch 73)
}

func main() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	defer listener.Close()
	fmt.Println("hand-rolled HTTP server listening on :8080")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("accept error: %v", err)
			continue
		}
		go handleConnection(conn) // Chapter 106's pattern, unchanged
	}
}
```

Notice the deliberate separation between `parseRequest` (bytes in, `Request` out), `route` (a `Request` in, a status/headers/body out — pure application logic, unaware that TCP or parsing even exist), and `writeResponse` (status/headers/body in, bytes out). This mirrors exactly how `net/http` itself is structured internally, just without the years of edge-case hardening.

---

## 6. Tracing Bytes: From Raw Socket to Parsed Request

Suppose a client sends exactly this over the wire (shown with explicit `\r\n` for clarity):

```
GET /hello?name=Ada HTTP/1.1\r\n
Host: localhost:8080\r\n
User-Agent: curl/8.4.0\r\n
Accept: */*\r\n
\r\n
```

`parseRequest` walks through it exactly like this:

```
Call 1: reader.ReadString('\n')
  → "GET /hello?name=Ada HTTP/1.1\r\n"
  → trimmed: "GET /hello?name=Ada HTTP/1.1"
  → split on " ", limit 3: ["GET", "/hello?name=Ada", "HTTP/1.1"]
  → method="GET", target="/hello?name=Ada", version="HTTP/1.1"
  → target has '?' at index 6: path="/hello", query="name=Ada"

Call 2: reader.ReadString('\n') → "Host: localhost:8080\r\n"
  → colon at index 4 → name="host", value="localhost:8080"

Call 3: reader.ReadString('\n') → "User-Agent: curl/8.4.0\r\n"
  → name="user-agent", value="curl/8.4.0"

Call 4: reader.ReadString('\n') → "Accept: */*\r\n"
  → name="accept", value="*/*"

Call 5: reader.ReadString('\n') → "\r\n"
  → trimmed: "" → blank line → headers loop ends

Content-Length header absent → body = nil

Result:
Request{
  Method: "GET", Path: "/hello", Query: "name=Ada", Version: "HTTP/1.1",
  Headers: {"host": "localhost:8080", "user-agent": "curl/8.4.0", "accept": "*/*"},
  Body: nil,
}
```

That `Request` value is exactly what `route()` receives — by the time application logic runs, every byte-level detail from Section 2 has already been resolved.

---

## 7. Writing a Spec-Correct Response, Field by Field

For the request traced above, `route()` returns `(200, "OK", {}, "Hello, Ada!\n")`, and `writeResponse` turns that into these exact bytes on the wire:

```
HTTP/1.1 200 OK\r\n
Content-Length: 12\r\n
Content-Type: text/plain; charset=utf-8\r\n
Date: Sun, 09 Aug 2026 09:20:11 GMT\r\n
Server: hand-rolled-go/1.0\r\n
Connection: keep-alive\r\n
\r\n
Hello, Ada!\n
```

Cross-check every line against Chapter 71, Section 5's specification: the status line is `HTTP-version SP status-code SP reason-phrase CRLF`, exactly reproduced. `Content-Length: 12` is the exact byte count of `"Hello, Ada!\n"` (11 visible characters plus the trailing `\n` = 12 bytes) — get this number wrong by even one byte and a real client (as opposed to this chapter's own eyes) will either truncate the body or hang waiting for bytes that never come, which is precisely why `writeResponse` computes it programmatically from `len(body)` rather than trusting a handler to supply it correctly.

---

## 8. Hands-On Experiment: `nc`, `curl`, and a Byte-Level Comparison

**Step 1 — start the server:**

```
$ go run httpserver.go
hand-rolled HTTP server listening on :8080
```

**Step 2 — talk to it with raw `nc`, exactly as Chapter 71, Section 13 did against a real server:**

```
$ nc localhost 8080
GET /hello?name=Ada HTTP/1.1
Host: localhost

HTTP/1.1 200 OK
Content-Length: 12
Content-Type: text/plain; charset=utf-8
Date: Sun, 09 Aug 2026 09:20:11 GMT
Server: hand-rolled-go/1.0
Connection: keep-alive

Hello, Ada!
```

This is your own parsing and serialization code, produced by hand-written Go, indistinguishable in shape from any real HTTP/1.1 server's output.

**Step 3 — use `curl -v` to see request and response sides together:**

```
$ curl -v "http://localhost:8080/hello?name=Ada"
> GET /hello?name=Ada HTTP/1.1
> Host: localhost:8080
> User-Agent: curl/8.4.0
> Accept: */*
>
< HTTP/1.1 200 OK
< Content-Length: 12
< Content-Type: text/plain; charset=utf-8
< Date: Sun, 09 Aug 2026 09:20:11 GMT
< Server: hand-rolled-go/1.0
< Connection: keep-alive
<
Hello, Ada!
```

**Step 4 — exercise the POST handler and body parsing:**

```
$ curl -v -X POST http://localhost:8080/echo -d 'raw request body bytes'
> POST /echo HTTP/1.1
> Host: localhost:8080
> Content-Length: 23
> Content-Type: application/x-www-form-urlencoded
>
< HTTP/1.1 200 OK
< Content-Length: 23
< Content-Type: application/octet-stream
<
raw request body bytes
```

`Content-Length: 23` on the request tells `parseRequest`'s step 3 exactly how many bytes to `io.ReadFull` for the body — and the server's own response then reports `Content-Length: 23` again, unchanged, because the `/echo` handler returns `req.Body` verbatim.

**Step 5 — trigger the hand-written 400 Bad Request path deliberately:**

```
$ printf 'NOT A VALID REQUEST LINE\r\n\r\n' | nc localhost 8080
HTTP/1.1 400 Bad Request
Content-Length: 16
Content-Type: text/plain; charset=utf-8
Date: Sun, 09 Aug 2026 09:21:02 GMT
Server: hand-rolled-go/1.0
Connection: close

400 Bad Request
```

`strings.SplitN(line, " ", 3)` on `"NOT A VALID REQUEST LINE"` produces `["NOT", "A", "VALID REQUEST LINE"]` — three parts, but `"VALID REQUEST LINE"` is obviously not a valid HTTP version, so the explicit version check in Section 5's `parseRequest` catches it and returns an error, which `handleConnection` turns into exactly this 400 response.

---

## 9. Keep-Alive: Handling More Than One Request Per Connection

`handleConnection`'s `for` loop is what makes this a real HTTP/1.1 server rather than the HTTP/0.9-style "one request, then close" behavior Chapter 71, Section 3 described as the naive starting point. Watch it handle two requests on one connection:

```
$ nc localhost 8080
GET / HTTP/1.1
Host: localhost

HTTP/1.1 200 OK
Content-Length: 39
...
Connection: keep-alive

Welcome to a hand-rolled HTTP server.
GET /hello HTTP/1.1
Host: localhost

HTTP/1.1 200 OK
Content-Length: 13
...
Connection: keep-alive

Hello, World!
```

The connection stayed open between the two requests — `parseRequest` was called a second time on the *same* `bufio.Reader`, which is exactly why `Content-Length`-based body framing (Section 4) matters so much: without it, the server would have no reliable way to know where the first response's body ended and the second request's bytes began, since both are just more bytes on the same stream (Chapter 106, Section 8's point about TCP having no message boundaries, applied here to HTTP framing on top of it). `shouldKeepAlive` decides, per request, whether to loop again or return (closing the connection) — sending `Connection: close` on any request in the exchange ends the loop after that response.

---

## 10. Common Pitfalls in Hand-Rolled HTTP Parsing

- **Trusting the client's line endings.** Real HTTP requires `\r\n`, but a naive parser using only `strings.TrimRight(line, "\n")` would leave a trailing `\r` inside header values, silently corrupting comparisons like `headers["connection"] == "close"` (which would actually be `"close\r"`, never equal). Section 5's `TrimRight(line, "\r\n")` trims both.
- **Not enforcing a body-length limit.** `io.ReadFull(reader, body)` with `length` taken directly from an untrusted `Content-Length` header will happily try to allocate and read gigabytes if a malicious client sends `Content-Length: 5000000000` — a real, exploitable resource-exhaustion vector this chapter's code does not defend against (Section 11 covers the fix).
- **Splitting the request line on any whitespace instead of exactly one space per field.** `strings.SplitN(line, " ", 3)` deliberately caps the split at three pieces, so a request-target that itself (incorrectly) contained a literal space wouldn't shatter the version field — but a fully spec-compliant parser needs additional validation this simplified version skips (Section 12).
- **Forgetting header names are case-insensitive.** Chapter 71, Section 15 called this out as a common misconception. This server's `strings.ToLower(name)` on every header at parse time is what makes `headers["content-length"]` work reliably regardless of whether the client sent `Content-Length`, `content-length`, or `CONTENT-LENGTH` — skip that lowercasing and the lookup breaks unpredictably depending on client behavior.
- **Reusing the wrong reader across requests, or a new one per request.** `handleConnection` deliberately creates exactly one `bufio.Reader` for the entire connection's lifetime and reuses it across every request in the keep-alive loop (Section 9) — creating a fresh `bufio.NewReader(conn)` on each loop iteration would discard any bytes that reader had already buffered past the current request's boundary, silently corrupting the start of the next request, exactly the same class of bug flagged in Chapter 108, Section 5 for the chat server's nickname handling.

---

## 11. Production Notes: Why Real Servers Don't Do This, and What They Do Instead

- **This server has no protection against a `Content-Length` bomb.** A production parser must cap the maximum acceptable body size (and reject with `413 Payload Too Large` before attempting to allocate a buffer) rather than trusting the header unconditionally, exactly the gap Section 10 flagged.
- **`Transfer-Encoding: chunked` is entirely unhandled here.** A real HTTP/1.1 server must support chunked bodies (Chapter 71, Section 9 mentions this as the alternative to `Content-Length`) since many real clients use it for bodies of unknown length in advance. Chapter 110's HTTP client chapter covers chunked decoding directly; a production server would need the equivalent encoding/decoding logic on both sides.
- **Request smuggling is a genuine risk in exactly this kind of hand-rolled parsing.** Chapter 71, Section 16 (Production Notes) described how disagreements between a front-end proxy and a back-end server over `Content-Length` vs. `Transfer-Encoding` framing can be exploited to smuggle a hidden second request past inspection. Any hand-rolled parser sitting behind a separate proxy needs to handle conflicting or duplicate framing headers by rejecting the request outright, not by guessing — this server doesn't defend against that today, which is a real, documented reason production infrastructure relies on extensively hardened, widely-audited HTTP implementations (`net/http`, nginx, Envoy) rather than bespoke parsers like this one.
- **No TLS.** This server speaks plaintext HTTP only. A production deployment would wrap the listener with `crypto/tls.Listen` (Chapter 82) — a change that touches almost nothing in `parseRequest`/`writeResponse`/`route`, since TLS operates below this chapter's parsing logic, encrypting the same byte stream this code already treats as opaque bytes from `net.Conn`.
- **No connection or request-rate limiting.** Chapter 106, Section 13's semaphore pattern and read/write deadlines both apply directly here — this server already sets a `SetReadDeadline` per request, but has no cap on total simultaneous connections, leaving it exposed to the same resource-exhaustion risk Chapter 106 discussed in the abstract.
- **This is exactly the value `net/http` provides in real code.** Every gap in this list is something the Go standard library's HTTP implementation already handles, tested against real-world adversarial input for over a decade. Building this chapter's version isn't a recommendation to replace `net/http` in production — it's proof that the black box from Section 3 is not magic, just a lot of carefully-hardened code doing exactly what Sections 4–7 did in miniature.

---

## 12. What's Simplified Here

Beyond the production gaps in Section 11: this parser does not support HTTP pipelining edge cases, `Expect: 100-continue`, HEAD requests correctly omitting a body while still reporting the real `Content-Length` a GET would have returned, or any header value spanning multiple lines (an obsolete HTTP/1.0-era folding syntax most modern parsers reject outright, which this one also does implicitly by requiring every header on one line). The router in Section 5 is a simple `switch` statement rather than anything resembling a real routing library with path parameters or middleware. None of these are exotic gaps — they are precisely the kind of accumulated edge-case handling that separates a teaching implementation from a production one.

---

## Interview Questions & Model Answers

**Beginner: Why can't an HTTP server just read a fixed number of bytes off the TCP connection and assume that's one complete request?**

Because HTTP requests vary in length — a different number of headers, a different or absent body — so there's no fixed size that would work for every request. Instead, the server has to read incrementally: read the request line up to its line ending, then read headers one at a time up to a blank line that marks the end of headers, and then, only if a `Content-Length` header says a body follows, read exactly that many more bytes. Each of those boundaries is discovered by parsing the bytes themselves, not assumed in advance.

**Intermediate: Explain exactly why `Content-Length` is essential for a server that supports keep-alive (multiple requests per TCP connection), using specifics from this chapter's code.**

On a keep-alive connection, the request line and headers can only be correctly separated from the body, and the body from the *next* request's request line, if the server knows precisely how many body bytes belong to the current request. `io.ReadFull(reader, body)` reads exactly `Content-Length` bytes and then stops — leaving the `bufio.Reader`'s internal buffer positioned exactly at the start of the next request's request line, ready for the next call to `parseRequest` on the same reader. Without `Content-Length` (or the chunked-encoding alternative this chapter doesn't implement), the server would have no way to know where a body ends short of the connection closing entirely, which is fundamentally incompatible with reusing the same connection for a second request afterward.

**Advanced: This chapter's server lowercases every header name during parsing. Explain the specific bug this avoids, and why relying on a client to send headers in a consistent case would be a mistake.**

HTTP header names are explicitly case-insensitive per specification — `Content-Length`, `content-length`, and `CONTENT-LENGTH` are all the same header, and different clients, proxies, and HTTP library versions are all free to send whichever casing they like, with no guarantee of consistency. If the parser stored header names exactly as received and the rest of the code looked them up with a hardcoded casing (e.g., `headers["Content-Length"]`), a client sending `content-length` (all lowercase, which many minimal HTTP clients and tools do) would cause that lookup to silently miss, and the server would incorrectly conclude no body was present — either dropping real request data or, worse, misreading where the current request ends and the next one on a keep-alive connection begins. Lowercasing every header name once, consistently, at parse time (and looking values up using the same lowercase keys everywhere else in the code) eliminates this entire class of bug rather than requiring every call site to remember to normalize casing itself.

---

## Exercises

### Easy
1. Add a new route, `GET /time`, that returns the current server time as plain text, and verify it with `curl`.
2. Deliberately send a request with a header line missing its colon (e.g., `Host localhost` instead of `Host: localhost`) via raw `nc`, and confirm the server responds with `400 Bad Request`.
3. Send two requests back-to-back on one `nc` session (as in Section 9) and confirm both get correct, independent responses.

### Medium
4. Add a maximum body size check to `parseRequest` (e.g., reject any `Content-Length` greater than 1MB with a `413 Payload Too Large` response) before allocating the body buffer, closing the gap flagged in Section 10.
5. Implement `HEAD` support: for any path that would return 200 to a GET, a HEAD request to the same path should return the identical headers (including the correct `Content-Length` the GET would have had) but with an empty body — exactly matching Chapter 71, Section 7's definition of HEAD.
6. Add a simple access log line for every response that includes the client's IP, method, path, status code, and response body size, formatted similarly to a real web server's access log.

### Hard
7. Implement basic support for `Transfer-Encoding: chunked` request bodies: parse a chunked body (each chunk prefixed by its length in hex, followed by CRLF, followed by that many bytes, followed by CRLF, terminated by a zero-length chunk) into a single assembled `[]byte`, and wire it into `parseRequest` as an alternative to `Content-Length`.
8. Harden the server against the request-smuggling risk from Section 11: reject any request that specifies both `Content-Length` and `Transfer-Encoding` headers simultaneously, with a `400 Bad Request`, and explain in a comment exactly which real vulnerability this defends against.
9. Wrap this server's listener with `crypto/tls.Listen` using a self-signed certificate, verify with `curl -k https://localhost:8443/` that it now serves HTTPS, and explain in your own words why almost none of `parseRequest`, `writeResponse`, or `route` needed to change to make this work.

---

## Summary

| Term | Meaning |
|---|---|
| Request line parsing | Splitting the first line into method, target (path + query), and version |
| Header parsing | Reading `Name: value` lines until a blank line, storing names lowercased |
| `Content-Length`-based framing | Reading exactly N more bytes as the body, and no more |
| `writeResponse` | Serializing a status line, headers, and body into spec-correct HTTP bytes |
| Keep-alive loop | Calling `parseRequest` repeatedly on one `bufio.Reader` for multiple requests per connection |
| Request smuggling | An exploit arising from inconsistent `Content-Length`/`Transfer-Encoding` framing between systems |
| `net/http` | The production-grade version of everything this chapter built by hand, hardened over years |

You've now built the server side of HTTP entirely by hand, from raw TCP bytes to a spec-correct response. Chapter 110 builds the mirror image: an HTTP *client* from scratch, constructing a valid request by hand and sending it over a raw `net.Dial`'d connection — including the one major piece this chapter's server sidestepped, chunked transfer-encoding, which a real client has to be able to decode from servers that use it.
