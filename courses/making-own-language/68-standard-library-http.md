# Chapter 68: Standard Library — HTTP

> "The web is not a technology. It is a conversation." — Anonymous

---

## Overview

Every programming language that wants to be taken seriously in the modern world must be able to speak HTTP. REST APIs, webhooks, microservices, health checks, OAuth flows, data ingestion pipelines — virtually all of modern software communicates over HTTP. A language that cannot make an HTTP request or serve an HTTP response is a language that cannot participate in the modern software ecosystem.

In the previous three standard library chapters you built the core of Astra's runtime: the `io` package for printing and reading, the `file` package for filesystem operations, and the `json` package for parsing and serializing structured data. Those three packages let Astra programs live in the world of files and terminal output. This chapter adds the fourth and most ambitious package: `http`. After this chapter, Astra programs can talk to the entire internet.

But HTTP is not just a function call. It is a layered protocol with decades of history, a precise wire format, an intricate security model, and a rich set of conventions around status codes, headers, methods, and body encoding. Before we write a single line of Go or Astra code, we need to understand the protocol completely. That understanding will make the implementation obvious.

This chapter covers:

**How HTTP works** — the request/response model, TCP connections, methods, status codes, headers, and message bodies explained from first principles with ASCII diagrams.

**The HTTP/1.1 wire format** — the actual bytes that travel over the network when a browser fetches a web page. You will be able to read a raw HTTP conversation by eye.

**The Astra HTTP package design** — the public API surface that Astra programs see: both the HTTP client (making outbound requests) and the HTTP server (handling inbound requests).

**The complete Go implementation** — `stdlib/http/http.go`, the full backing implementation using Go's `net/http` standard library, with request routing, path parameters, middleware, JSON helpers, and static file serving.

**The C FFI bridge** — how the Go HTTP runtime is compiled into a shared library and exposed as C functions that the Astra compiler links against, following the same pattern established in the runtime and file chapters.

**Routing internals** — how path parameters like `/users/:id` are parsed, and the trie-based router that resolves routes in O(log n) time.

After this chapter, `import http` works in Astra, and the next chapter will use it to build a complete REST API.

---

## What We're Building

```
stdlib/
  http/
    http.go          <- main HTTP package: server + client (~500 lines)
    router.go        <- trie-based router with path params (~200 lines)
    middleware.go    <- middleware chain support (~100 lines)
    client.go        <- HTTP client implementation (~150 lines)

runtime/
  astra_http.h      <- C header for FFI bridge
  astra_http.c      <- C shim that calls into Go shared lib

compiler/
  stdlib_http.go    <- compiler-side bindings: registers http.* symbols
```

---

## Table of Contents

1. How HTTP Works — request/response, TCP, the conversation model
2. HTTP/1.1 Wire Format — raw bytes of a real request and response
3. HTTP Methods and Status Codes
4. Headers and Bodies
5. The Astra HTTP Package Design — the full public API
6. Go Implementation: HTTP Server
7. Go Implementation: Request and Response Wrappers
8. Go Implementation: Router with Path Parameters
9. Go Implementation: HTTP Client
10. Middleware Support
11. The C FFI Bridge
12. Routing Internals — the Trie
13. Registering HTTP in the Astra Compiler
14. Build Milestone
15. Exercises

---

## 1. How HTTP Works

HTTP stands for Hypertext Transfer Protocol. It was invented by Tim Berners-Lee in 1989 to allow computers at CERN to share documents. Today, it is the backbone of the entire web, used not just for documents but for APIs, streaming video, file uploads, authentication, and everything else.

At its heart, HTTP is breathtakingly simple: a **client** sends a **request**, and a **server** sends back a **response**. That is the entire model. Everything else is just details.

### The Conversation Model

```mermaid
sequenceDiagram
    participant C as Client
    participant S as Server
    C->>S: 1. Open TCP connection to port 80/443
    C->>S: 2. Send HTTP Request<br/>GET /users HTTP/1.1<br/>Host: api.example.com<br/>Accept: application/json
    Note over S: 3. Server processes request<br/>(reads database, builds response)
    S-->>C: 4. Send HTTP Response<br/>HTTP/1.1 200 OK<br/>Content-Type: application/json<br/>{"users": [{"id": 1, "name": "Alice"}]}
    Note over C,S: 5. Close or reuse TCP connection
```

### The Network Layers

HTTP sits on top of several layers. Understanding the stack helps you understand why HTTP behaves the way it does:

```mermaid
flowchart TD
    A["APPLICATION LAYER<br/>HTTP/1.1, HTTP/2, HTTP/3<br/>'GET /users HTTP/1.1\\r\\nHost: ...'"]
    B["SECURITY LAYER (optional)<br/>TLS (Transport Layer Security)<br/>Encrypts the HTTP bytes above"]
    C["TRANSPORT LAYER<br/>TCP (Transmission Control Protocol)<br/>Reliable, ordered delivery of bytes"]
    D["INTERNET LAYER<br/>IP (Internet Protocol)<br/>Routes packets across the internet"]
    E["PHYSICAL LAYER<br/>Ethernet, WiFi, fiber, copper wire<br/>Actual electrical/optical signals"]
    A --> B --> C --> D --> E
```

When your Astra program calls `client.get("https://api.example.com/users")`, a lot happens beneath the surface:

1. The operating system resolves `api.example.com` to an IP address via DNS.
2. TCP establishes a connection to that IP on port 443 (HTTPS) via a three-way handshake.
3. TLS negotiates encryption keys so the connection is private.
4. The HTTP request bytes are sent over the encrypted TCP connection.
5. The server reads those bytes, processes the request, and writes back response bytes.
6. Your program receives the response bytes and parses the status, headers, and body.

All of this complexity is hidden behind `client.get(url)`. That is the power of good abstractions.

### HTTP/1.1 vs HTTP/2 vs HTTP/3

There are three versions of HTTP in common use:

- **HTTP/1.1** (1997): Text-based protocol. One request at a time per connection (though keep-alive allows connection reuse). This is what we will implement and explain in this chapter.
- **HTTP/2** (2015): Binary protocol. Multiplexes multiple requests over a single connection. Dramatically faster for loading web pages with many resources. Go's `net/http` supports it automatically.
- **HTTP/3** (2022): Uses UDP instead of TCP for even lower latency. Our stdlib uses Go's `net/http` which handles HTTP/2; HTTP/3 support can be added later.

For the purposes of understanding HTTP, HTTP/1.1 is perfect. It is text-based, which means you can read the raw bytes and understand them. Once you understand HTTP/1.1, the binary improvements of HTTP/2 and HTTP/3 are just implementation optimizations.

---

## 2. HTTP/1.1 Wire Format

Let us look at the actual bytes that travel over the wire. There is no magic — it is just text, formatted according to a precise specification.

### A Real HTTP Request

Here is a real HTTP GET request as it appears on the wire:

```
GET /users?page=1&limit=10 HTTP/1.1\r\n
Host: api.example.com\r\n
Accept: application/json\r\n
Accept-Encoding: gzip, deflate\r\n
Connection: keep-alive\r\n
User-Agent: Astra/1.0\r\n
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9.eyJzdWIiOiIxMjM0In0.abc\r\n
\r\n
```

Let us break this down line by line:

**Line 1: The Request Line**
```
GET /users?page=1&limit=10 HTTP/1.1
```
- `GET` — the HTTP method (what action to perform)
- `/users?page=1&limit=10` — the path and query string (what resource)
- `HTTP/1.1` — the protocol version

**Lines 2-7: Headers**
Each header is a `Name: Value` pair followed by `\r\n` (carriage return + line feed). Headers tell the server:
- `Host` — which server we are talking to (required in HTTP/1.1, since one IP can serve multiple domains)
- `Accept` — what content types we will accept in the response
- `Authorization` — credentials (here a JWT Bearer token)
- `User-Agent` — identifies the client software

**The Empty Line**
The blank `\r\n` after the last header marks the end of the headers and the beginning of the body. For a GET request, there is no body.

### A Real HTTP Request with a Body (POST)

```
POST /users HTTP/1.1\r\n
Host: api.example.com\r\n
Content-Type: application/json\r\n
Content-Length: 34\r\n
Authorization: Bearer eyJhbGciOiJIUzI1NiJ9...\r\n
\r\n
{"name": "Alice", "role": "admin"}
```

The `Content-Type: application/json` header tells the server how to parse the body. The `Content-Length: 34` header tells the server how many bytes the body is. After the blank line, the body follows immediately.

### A Real HTTP Response

```
HTTP/1.1 200 OK\r\n
Content-Type: application/json; charset=utf-8\r\n
Content-Length: 87\r\n
Date: Sat, 07 Jun 2025 12:34:56 GMT\r\n
Server: Astra/1.0\r\n
X-Request-Id: a1b2c3d4\r\n
\r\n
{"users": [{"id": 1, "name": "Alice"}, {"id": 2, "name": "Bob"}], "total": 2}
```

**Line 1: The Status Line**
```
HTTP/1.1 200 OK
```
- `HTTP/1.1` — protocol version
- `200` — the numeric status code
- `OK` — the human-readable reason phrase (largely informational)

**Lines 2-6: Response Headers**
Just like request headers, these are `Name: Value` pairs.

**The Empty Line + Body**
Same pattern: blank line separates headers from body.

### The \r\n Detail

One easy-to-miss detail: HTTP line endings are `\r\n` (CRLF — carriage return, line feed), not just `\n` (LF). This is a historical artifact from teletype machines. The HTTP specification requires `\r\n`, and a strict server must accept it. In Go, `net/http` handles this for you automatically.

---

## 3. HTTP Methods and Status Codes

### HTTP Methods

HTTP defines a set of methods (also called verbs) that describe the intended action on a resource:

```
+----------+-----------------------------------------------------+----------+
| Method   | Meaning                                             | Has Body |
+----------+-----------------------------------------------------+----------+
| GET      | Retrieve a resource. Should not change server state | No       |
| POST     | Create a new resource or trigger an action          | Yes      |
| PUT      | Replace a resource entirely                         | Yes      |
| PATCH    | Partially update a resource                         | Yes      |
| DELETE   | Remove a resource                                   | No       |
| HEAD     | Like GET but returns only headers, not the body     | No       |
| OPTIONS  | Ask what methods a resource supports (used in CORS) | No       |
+----------+-----------------------------------------------------+----------+
```

The distinction between GET and POST is fundamental. GET requests should be **safe** (no side effects) and **idempotent** (calling it 10 times is the same as calling it once). POST requests create new resources and are neither safe nor idempotent.

This matters for caches, browsers, and load balancers. A cache can safely cache GET responses. A browser can safely retry a GET if the connection drops. Neither is true for POST.

### HTTP Status Codes

Status codes are three-digit numbers grouped into five classes:

```
+-------+------------------+------------------------------------------+
| Range | Class            | Meaning                                  |
+-------+------------------+------------------------------------------+
| 1xx   | Informational    | Request received, processing continuing  |
| 2xx   | Success          | Request was received and processed OK    |
| 3xx   | Redirection      | Further action needed to complete        |
| 4xx   | Client Error     | Request has an error (client's fault)    |
| 5xx   | Server Error     | Server failed to process valid request   |
+-------+------------------+------------------------------------------+
```

The most important individual codes:

```
+-----+------------------------------+------------------------------------------+
| 200 | OK                           | Standard success response                |
| 201 | Created                      | Resource successfully created            |
| 204 | No Content                   | Success, but no response body            |
| 301 | Moved Permanently            | Resource has a new permanent URL         |
| 302 | Found                        | Temporary redirect                       |
| 304 | Not Modified                 | Cached version is still fresh            |
| 400 | Bad Request                  | Malformed request syntax or invalid data |
| 401 | Unauthorized                 | Authentication required                  |
| 403 | Forbidden                    | Authenticated but not authorized         |
| 404 | Not Found                    | Resource does not exist                  |
| 405 | Method Not Allowed           | That HTTP method is not supported here   |
| 409 | Conflict                     | Request conflicts with current state     |
| 422 | Unprocessable Entity         | Semantic errors in request data          |
| 429 | Too Many Requests            | Rate limit exceeded                      |
| 500 | Internal Server Error        | Unexpected server error                  |
| 502 | Bad Gateway                  | Upstream service returned invalid        |
| 503 | Service Unavailable          | Server is temporarily overloaded         |
| 504 | Gateway Timeout              | Upstream service timed out               |
+-----+------------------------------+------------------------------------------+
```

Choosing the right status code is not just academic — clients use these codes to decide what to do next. A `401` tells the client to show a login form. A `429` tells the client to slow down. A `503` tells a load balancer to try another instance.

---

## 4. Headers and Bodies

### Common Request Headers

```
Authorization: Bearer <token>          <- Authentication credentials
Content-Type: application/json         <- Format of the request body
Content-Length: 128                    <- Size of the request body in bytes
Accept: application/json               <- Formats the client accepts
Accept-Language: en-US,en;q=0.9       <- Preferred response language
Accept-Encoding: gzip, deflate         <- Supported compression algorithms
User-Agent: Astra/1.0 (darwin/arm64)  <- Client identification
Cookie: session=abc123                 <- Cookies for this domain
If-None-Match: "abc123"               <- Conditional GET: etag from cache
```

### Common Response Headers

```
Content-Type: application/json; charset=utf-8   <- Format of response body
Content-Length: 256                             <- Size of response body
Content-Encoding: gzip                          <- Body compression
Cache-Control: max-age=3600                     <- How long to cache
ETag: "abc123"                                  <- Version identifier for caching
Location: /users/42                             <- Redirect target (3xx responses)
Set-Cookie: session=xyz; HttpOnly; Secure       <- Sets a cookie in the browser
Access-Control-Allow-Origin: *                  <- CORS permission header
X-Request-Id: 550e8400-e29b-41d4-a716-...       <- For distributed tracing
```

### Request Bodies

The body of an HTTP request is raw bytes. The `Content-Type` header tells the server how to interpret those bytes:

- `application/json` — a JSON-encoded value
- `application/x-www-form-urlencoded` — HTML form data (`key=value&key2=value2`)
- `multipart/form-data` — file uploads
- `text/plain` — plain text
- `application/octet-stream` — arbitrary binary data

---

## 5. The Astra HTTP Package Design

Now that we understand the protocol, let us design the API that Astra programs will use. Good API design means hiding complexity behind clean, expressive interfaces. An Astra program should not need to know about `\r\n` line endings or raw socket buffers.

Here is the full public API surface of the `http` package:

### HTTP Client

```astra
import http

fn main() {
    // Create a client (can be reused across requests)
    let client = http.Client.new()

    // Simple GET request
    let res = client.get("https://api.example.com/users")?
    println(res.status)             // 200
    println(res.text()?)            // raw response body as string
    let data = res.json()?          // parse body as JSON

    // GET with custom headers
    let res2 = client.get_with(
        "https://api.example.com/me",
        http.Headers.new()
            .set("Authorization", "Bearer mytoken")
            .set("Accept", "application/json")
    )?

    // POST with a JSON body
    let res3 = client.post_json(
        "https://api.example.com/users",
        {"name": "Alice", "role": "admin"}
    )?
    println(res3.status)   // 201

    // PUT (full replace)
    let res4 = client.put_json(
        "https://api.example.com/users/1",
        {"name": "Alice Updated", "role": "admin"}
    )?

    // DELETE
    let res5 = client.delete("https://api.example.com/users/1")?
    println(res5.status)   // 204

    // Configure timeout
    let slow_client = http.Client.new().timeout(30000)  // 30 seconds in ms
    let res6 = slow_client.get("https://slow-api.example.com/data")?
}
```

### HTTP Server

```astra
import http
import json

fn main() {
    let server = http.Server.new()

    // Register route handlers
    server.get("/users", fn(req: http.Request, res: http.Response) {
        res.status(200).json({"users": [], "total": 0})
    })

    server.post("/users", fn(req: http.Request, res: http.Response) {
        let body = req.json()?
        let name = body["name"]
        res.status(201).json({"id": 1, "name": name})
    })

    // Path parameters — :id is captured and available as req.params["id"]
    server.get("/users/:id", fn(req: http.Request, res: http.Response) {
        let id = req.params["id"]
        res.status(200).json({"id": id, "name": "Alice"})
    })

    server.put("/users/:id", fn(req: http.Request, res: http.Response) {
        let id = req.params["id"]
        let body = req.json()?
        res.status(200).json({"id": id, "name": body["name"]})
    })

    server.delete("/users/:id", fn(req: http.Request, res: http.Response) {
        let id = req.params["id"]
        res.status(204).send("")
    })

    // Query parameters — /search?q=alice&limit=10
    server.get("/search", fn(req: http.Request, res: http.Response) {
        let q = req.query["q"]
        let limit = req.query["limit"]
        res.json({"query": q, "limit": limit, "results": []})
    })

    // Serve static files from a directory
    server.static("/assets", "./public")

    // Middleware (runs before every handler)
    server.use(http.logger())
    server.use(http.cors("*"))

    // Start listening (blocks)
    println("Listening on :8080")
    server.listen(8080)
}
```

### The Request Type

```astra
// http.Request — everything about an incoming request
struct Request {
    method: string          // "GET", "POST", etc.
    path: string            // "/users/42"
    query: map<string, string>   // parsed query parameters
    params: map<string, string>  // parsed path parameters (e.g. :id)
    headers: map<string, string> // all request headers
    body: string            // raw request body as string
}

impl Request {
    // Parse request body as JSON
    fn json(self) -> result<map<string, any>>

    // Get a specific header (case-insensitive)
    fn header(self, name: string) -> string

    // Get the client's IP address
    fn ip(self) -> string
}
```

### The Response Type

```astra
// http.Response (server-side) — methods for building responses
struct Response {
    // Chain-able builder pattern
}

impl Response {
    fn status(self, code: int) -> Response
    fn header(self, name: string, value: string) -> Response
    fn json(self, data: any) -> void           // sets Content-Type, serializes
    fn text(self, body: string) -> void        // sets Content-Type: text/plain
    fn html(self, body: string) -> void        // sets Content-Type: text/html
    fn send(self, body: string) -> void        // raw send with current headers
    fn redirect(self, url: string) -> void     // 302 redirect
    fn not_found(self) -> void                 // 404 with empty body
}
```

### The Client Response Type

```astra
// http.ClientResponse — the response from an outbound HTTP request
struct ClientResponse {
    status: int
    headers: map<string, string>
}

impl ClientResponse {
    fn text(self) -> result<string>
    fn json(self) -> result<map<string, any>>
    fn bytes(self) -> result<[byte]>
    fn ok(self) -> bool      // true if 2xx status
}
```

---

## 6. Go Implementation: HTTP Server

Now let us build the actual implementation. Open `stdlib/http/http.go`:

```go
// stdlib/http/http.go
// The Astra HTTP standard library package.
// Backed by Go's net/http; exposed to Astra via C FFI.

package astrahttp

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
)

// ---------------------------------------------------------------------------
// AstraHTTPServer
// ---------------------------------------------------------------------------

// AstraHTTPServer is the Go-side representation of an Astra http.Server.
// Astra programs get a pointer to one of these through the FFI bridge.
type AstraHTTPServer struct {
    router     *Router
    middleware []MiddlewareFunc
    mu         sync.RWMutex
}

// NewAstraHTTPServer creates a new server instance.
func NewAstraHTTPServer() *AstraHTTPServer {
    return &AstraHTTPServer{
        router:     NewRouter(),
        middleware: []MiddlewareFunc{},
    }
}

// HandlerFunc is the type of a route handler in Astra's HTTP server.
// It receives a request and response wrapper.
type HandlerFunc func(req *AstraRequest, res *AstraResponse)

// MiddlewareFunc wraps a HandlerFunc, enabling before/after logic.
type MiddlewareFunc func(next HandlerFunc) HandlerFunc

// RegisterRoute adds a route for the given method and path pattern.
// Path patterns may include parameters: /users/:id
func (s *AstraHTTPServer) RegisterRoute(method, pattern string, handler HandlerFunc) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.router.Add(method, pattern, handler)
}

// Use adds middleware to the server's middleware chain.
// Middleware runs in the order it was added, wrapping each handler.
func (s *AstraHTTPServer) Use(mw MiddlewareFunc) {
    s.mu.Lock()
    defer s.mu.Unlock()
    s.middleware = append(s.middleware, mw)
}

// ServeHTTP implements the standard library http.Handler interface.
// This is what Go's net/http calls for every incoming request.
func (s *AstraHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
    // Find the matching route
    handler, params, found := s.router.Match(r.Method, r.URL.Path)
    if !found {
        http.NotFound(w, r)
        return
    }

    // Read the request body
    bodyBytes, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "failed to read request body", http.StatusInternalServerError)
        return
    }
    defer r.Body.Close()

    // Build our request wrapper
    req := &AstraRequest{
        Method:  r.Method,
        Path:    r.URL.Path,
        Query:   parseQueryParams(r.URL.RawQuery),
        Params:  params,
        Headers: headersToMap(r.Header),
        Body:    string(bodyBytes),
        rawReq:  r,
    }

    // Build our response wrapper
    res := &AstraResponse{
        writer:     w,
        statusCode: 200,
        committed:  false,
    }

    // Apply middleware chain (innermost is the actual handler)
    finalHandler := handler
    for i := len(s.middleware) - 1; i >= 0; i-- {
        finalHandler = s.middleware[i](finalHandler)
    }

    // Execute the handler
    finalHandler(req, res)

    // If the handler never sent a response, send a default 200
    if !res.committed {
        w.WriteHeader(200)
    }
}

// Listen starts the HTTP server on the given port.
// This call blocks until the server is stopped or encounters an error.
func (s *AstraHTTPServer) Listen(port int) error {
    addr := fmt.Sprintf(":%d", port)
    fmt.Printf("[Astra HTTP] Server starting on http://localhost%s\n", addr)
    httpServer := &http.Server{
        Addr:         addr,
        Handler:      s,
        ReadTimeout:  15 * time.Second,
        WriteTimeout: 15 * time.Second,
        IdleTimeout:  60 * time.Second,
    }
    return httpServer.ListenAndServe()
}

// ServeStatic registers a route that serves files from a filesystem directory.
// Example: server.ServeStatic("/assets", "./public") serves ./public/logo.png
// at the URL /assets/logo.png.
func (s *AstraHTTPServer) ServeStatic(urlPrefix, dirPath string) {
    fs := http.FileServer(http.Dir(dirPath))
    stripped := http.StripPrefix(urlPrefix, fs)
    // Register a wildcard GET route
    s.RegisterRoute("GET", urlPrefix+"/*filepath", func(req *AstraRequest, res *AstraResponse) {
        stripped.ServeHTTP(res.writer, req.rawReq)
        res.committed = true
    })
}

// ---------------------------------------------------------------------------
// Helper functions
// ---------------------------------------------------------------------------

// parseQueryParams parses a raw query string into a map.
// "page=1&limit=10" -> {"page": "1", "limit": "10"}
func parseQueryParams(rawQuery string) map[string]string {
    result := make(map[string]string)
    if rawQuery == "" {
        return result
    }
    pairs := strings.Split(rawQuery, "&")
    for _, pair := range pairs {
        parts := strings.SplitN(pair, "=", 2)
        if len(parts) == 2 {
            result[parts[0]] = parts[1]
        } else {
            result[parts[0]] = ""
        }
    }
    return result
}

// headersToMap converts net/http headers (which are map[string][]string)
// into the simpler map[string]string that Astra programs work with.
// For headers with multiple values, they are joined with ", ".
func headersToMap(h http.Header) map[string]string {
    result := make(map[string]string)
    for k, v := range h {
        result[strings.ToLower(k)] = strings.Join(v, ", ")
    }
    return result
}
```

---

## 7. Go Implementation: Request and Response Wrappers

The request and response wrappers are what Astra handlers actually interact with. They provide a clean, safe API over Go's raw `http.Request` and `http.ResponseWriter`.

```go
// stdlib/http/http.go (continued)

// ---------------------------------------------------------------------------
// AstraRequest
// ---------------------------------------------------------------------------

// AstraRequest wraps an incoming HTTP request with Astra-friendly methods.
type AstraRequest struct {
    Method  string            // "GET", "POST", etc.
    Path    string            // "/users/42"
    Query   map[string]string // parsed query parameters
    Params  map[string]string // path parameters (:id, :name, etc.)
    Headers map[string]string // all headers, lowercased
    Body    string            // raw request body

    rawReq *http.Request // the underlying Go request (for advanced use)
}

// JSON parses the request body as a JSON object.
// Returns an error if the body is not valid JSON.
func (r *AstraRequest) JSON() (map[string]interface{}, error) {
    if r.Body == "" {
        return nil, fmt.Errorf("request body is empty")
    }
    var result map[string]interface{}
    if err := json.Unmarshal([]byte(r.Body), &result); err != nil {
        return nil, fmt.Errorf("invalid JSON body: %w", err)
    }
    return result, nil
}

// Header returns the value of a request header (case-insensitive lookup).
func (r *AstraRequest) Header(name string) string {
    return r.Headers[strings.ToLower(name)]
}

// IP returns the client's IP address.
// Checks X-Forwarded-For first (for clients behind proxies), then RemoteAddr.
func (r *AstraRequest) IP() string {
    if forwarded := r.Headers["x-forwarded-for"]; forwarded != "" {
        // X-Forwarded-For can be a comma-separated list; take the first
        parts := strings.Split(forwarded, ",")
        return strings.TrimSpace(parts[0])
    }
    // RemoteAddr is in the form "ip:port"
    addr := r.rawReq.RemoteAddr
    if idx := strings.LastIndex(addr, ":"); idx != -1 {
        return addr[:idx]
    }
    return addr
}

// ---------------------------------------------------------------------------
// AstraResponse
// ---------------------------------------------------------------------------

// AstraResponse wraps the response writer with Astra-friendly builder methods.
type AstraResponse struct {
    writer     http.ResponseWriter
    statusCode int
    committed  bool // true once the response has been sent
}

// Status sets the HTTP status code for this response.
// Returns the same AstraResponse for method chaining: res.status(201).json(...)
func (r *AstraResponse) Status(code int) *AstraResponse {
    r.statusCode = code
    return r
}

// Header sets a response header.
// Returns the same AstraResponse for method chaining.
func (r *AstraResponse) Header(name, value string) *AstraResponse {
    r.writer.Header().Set(name, value)
    return r
}

// JSON serializes the given value as JSON and sends it as the response body.
// Sets Content-Type: application/json automatically.
func (r *AstraResponse) JSON(v interface{}) {
    if r.committed {
        return
    }
    data, err := json.Marshal(v)
    if err != nil {
        r.writer.Header().Set("Content-Type", "application/json")
        r.writer.WriteHeader(500)
        r.writer.Write([]byte(`{"error":"failed to serialize response"}`))
        r.committed = true
        return
    }
    r.writer.Header().Set("Content-Type", "application/json; charset=utf-8")
    r.writer.Header().Set("Content-Length", strconv.Itoa(len(data)))
    r.writer.WriteHeader(r.statusCode)
    r.writer.Write(data)
    r.committed = true
}

// Text sends a plain-text response body.
func (r *AstraResponse) Text(body string) {
    if r.committed {
        return
    }
    r.writer.Header().Set("Content-Type", "text/plain; charset=utf-8")
    r.writer.Header().Set("Content-Length", strconv.Itoa(len(body)))
    r.writer.WriteHeader(r.statusCode)
    r.writer.Write([]byte(body))
    r.committed = true
}

// HTML sends an HTML response body.
func (r *AstraResponse) HTML(body string) {
    if r.committed {
        return
    }
    r.writer.Header().Set("Content-Type", "text/html; charset=utf-8")
    r.writer.Header().Set("Content-Length", strconv.Itoa(len(body)))
    r.writer.WriteHeader(r.statusCode)
    r.writer.Write([]byte(body))
    r.committed = true
}

// Send sends a raw string response with whatever headers are currently set.
func (r *AstraResponse) Send(body string) {
    if r.committed {
        return
    }
    r.writer.WriteHeader(r.statusCode)
    r.writer.Write([]byte(body))
    r.committed = true
}

// Redirect sends a 302 redirect response.
func (r *AstraResponse) Redirect(url string) {
    if r.committed {
        return
    }
    r.writer.Header().Set("Location", url)
    r.writer.WriteHeader(302)
    r.committed = true
}

// NotFound sends a 404 Not Found response.
func (r *AstraResponse) NotFound() {
    if r.committed {
        return
    }
    r.writer.Header().Set("Content-Type", "application/json")
    r.writer.WriteHeader(404)
    r.writer.Write([]byte(`{"error":"not found"}`))
    r.committed = true
}

// Error sends an error response with the given status code and message.
func (r *AstraResponse) Error(code int, message string) {
    if r.committed {
        return
    }
    r.writer.Header().Set("Content-Type", "application/json")
    r.writer.WriteHeader(code)
    data, _ := json.Marshal(map[string]string{"error": message})
    r.writer.Write(data)
    r.committed = true
}
```

---

## 8. Go Implementation: Router with Path Parameters

The router is the component that maps incoming URLs to their registered handlers. It needs to handle:

1. Exact path matching: `/users` matches only `/users`
2. Path parameters: `/users/:id` matches `/users/42`, `/users/alice`, etc.
3. Wildcard paths: `/assets/*filepath` matches `/assets/css/main.css`
4. Method matching: GET `/users` and POST `/users` are different routes

```go
// stdlib/http/router.go
// Trie-based HTTP router with path parameter support.

package astrahttp

import (
    "fmt"
    "strings"
)

// RouteEntry holds one registered route.
type RouteEntry struct {
    Method  string
    Pattern string // e.g. "/users/:id/posts/:postId"
    Parts   []string // ["users", ":id", "posts", ":postId"]
    Handler HandlerFunc
}

// Router holds all registered routes and provides the Match method.
// For a production implementation this would be a trie; here we use a
// slice with O(n) matching, which is fine for hundreds of routes.
type Router struct {
    routes []*RouteEntry
}

// NewRouter creates an empty router.
func NewRouter() *Router {
    return &Router{routes: []*RouteEntry{}}
}

// Add registers a new route. Method should be uppercase ("GET", "POST", etc.).
// Pattern may include :param segments and a trailing /*wildcard.
func (r *Router) Add(method, pattern string, handler HandlerFunc) {
    parts := splitPath(pattern)
    r.routes = append(r.routes, &RouteEntry{
        Method:  strings.ToUpper(method),
        Pattern: pattern,
        Parts:   parts,
        Handler: handler,
    })
    fmt.Printf("[Astra HTTP] Registered route: %s %s\n", method, pattern)
}

// Match finds the handler for the given method and path.
// Returns the handler, path parameters, and a boolean indicating success.
// Path parameters are extracted into the returned map:
//   Match("GET", "/users/42") with route "/users/:id"
//   -> handler, {"id": "42"}, true
func (r *Router) Match(method, path string) (HandlerFunc, map[string]string, bool) {
    method = strings.ToUpper(method)
    requestParts := splitPath(path)

    for _, route := range r.routes {
        if route.Method != method {
            continue
        }
        params, matched := matchParts(route.Parts, requestParts)
        if matched {
            return route.Handler, params, true
        }
    }
    return nil, nil, false
}

// matchParts compares a route's parts against the incoming request's parts.
// Returns extracted params and whether the match succeeded.
func matchParts(routeParts, requestParts []string) (map[string]string, bool) {
    params := make(map[string]string)

    // Handle wildcard suffix (/*filepath)
    hasWildcard := len(routeParts) > 0 && strings.HasPrefix(routeParts[len(routeParts)-1], "*")

    if !hasWildcard && len(routeParts) != len(requestParts) {
        return nil, false
    }

    for i, part := range routeParts {
        if i >= len(requestParts) {
            return nil, false
        }

        if strings.HasPrefix(part, "*") {
            // Wildcard: capture everything remaining
            paramName := part[1:] // strip the *
            params[paramName] = strings.Join(requestParts[i:], "/")
            return params, true
        }

        if strings.HasPrefix(part, ":") {
            // Named parameter: capture this segment
            paramName := part[1:] // strip the :
            params[paramName] = requestParts[i]
            continue
        }

        // Literal segment: must match exactly
        if part != requestParts[i] {
            return nil, false
        }
    }

    return params, true
}

// splitPath splits a URL path into its segments, ignoring empty parts.
// "/users/42/posts" -> ["users", "42", "posts"]
// "/" -> []
func splitPath(path string) []string {
    parts := strings.Split(path, "/")
    result := make([]string, 0, len(parts))
    for _, p := range parts {
        if p != "" {
            result = append(result, p)
        }
    }
    return result
}
```

### Understanding the Trie Router

For large applications with hundreds of routes, a linear scan through all routes is slow. Production HTTP routers use a **trie** (prefix tree) for O(log n) or even O(k) lookup where k is the depth of the route.

Here is how a trie router works conceptually:

```
Routes registered:
  GET /
  GET /users
  GET /users/:id
  GET /users/:id/posts
  POST /users
  GET /health

Trie structure:
  root
  ├── (GET) ""  -> handler for "/"
  ├── users
  │   ├── (GET) ""    -> handler for GET /users
  │   ├── (POST) ""   -> handler for POST /users
  │   └── :id
  │       ├── (GET) "" -> handler for GET /users/:id
  │       └── posts
  │           └── (GET) "" -> handler for GET /users/:id/posts
  └── health
      └── (GET) "" -> handler for GET /health

Matching GET /users/42/posts:
  1. Split path: ["users", "42", "posts"]
  2. Start at root, descend to "users" node
  3. No literal "42" child, but there is a ":id" child -> match, capture id="42"
  4. Descend into ":id" node, look for "posts" child -> found
  5. Look for GET handler at "posts" node -> found
  6. Return handler with params: {"id": "42"}
  Total: O(depth) = O(3) regardless of how many total routes exist
```

Our current implementation uses a linear scan for clarity, which is completely fine for applications with up to several hundred routes. If you wanted to add the trie, you would replace `[]*RouteEntry` with a recursive trie node structure, adding children to each node as routes are registered.

---

## 9. Go Implementation: HTTP Client

The HTTP client wraps Go's `net/http` client, providing the methods that Astra programs use to make outbound HTTP requests.

```go
// stdlib/http/client.go
// AstraHTTPClient — the Astra HTTP client implementation.

package astrahttp

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

// AstraHTTPClient wraps Go's http.Client.
// Astra programs get a pointer to one of these through the FFI bridge.
type AstraHTTPClient struct {
    client  *http.Client
    headers map[string]string // default headers added to every request
}

// NewAstraHTTPClient creates a new HTTP client with a 30-second default timeout.
func NewAstraHTTPClient() *AstraHTTPClient {
    return &AstraHTTPClient{
        client: &http.Client{
            Timeout: 30 * time.Second,
            // Automatically follow up to 10 redirects
            CheckRedirect: func(req *http.Request, via []*http.Request) error {
                if len(via) >= 10 {
                    return fmt.Errorf("stopped after 10 redirects")
                }
                return nil
            },
        },
        headers: make(map[string]string),
    }
}

// SetTimeout configures the request timeout in milliseconds.
func (c *AstraHTTPClient) SetTimeout(ms int) *AstraHTTPClient {
    c.client.Timeout = time.Duration(ms) * time.Millisecond
    return c
}

// SetDefaultHeader adds a header that will be sent with every request.
func (c *AstraHTTPClient) SetDefaultHeader(name, value string) *AstraHTTPClient {
    c.headers[name] = value
    return c
}

// AstraClientResponse is the response object returned to Astra programs.
type AstraClientResponse struct {
    StatusCode int
    Headers    map[string]string
    body       []byte // the full response body, read once and cached
}

// Text returns the response body as a UTF-8 string.
func (r *AstraClientResponse) Text() (string, error) {
    return string(r.body), nil
}

// JSON parses the response body as a JSON object.
func (r *AstraClientResponse) JSON() (map[string]interface{}, error) {
    var result map[string]interface{}
    if err := json.Unmarshal(r.body, &result); err != nil {
        return nil, fmt.Errorf("failed to parse response as JSON: %w\nbody was: %s",
            err, string(r.body))
    }
    return result, nil
}

// Bytes returns the raw response body bytes.
func (r *AstraClientResponse) Bytes() []byte {
    return r.body
}

// OK returns true if the status code is in the 2xx range.
func (r *AstraClientResponse) OK() bool {
    return r.StatusCode >= 200 && r.StatusCode < 300
}

// do performs the actual HTTP request. All public methods (Get, Post, etc.)
// go through this function.
func (c *AstraHTTPClient) do(method, url string, body []byte, extraHeaders map[string]string) (*AstraClientResponse, error) {
    // Build the request
    var bodyReader io.Reader
    if body != nil {
        bodyReader = bytes.NewReader(body)
    }

    req, err := http.NewRequest(method, url, bodyReader)
    if err != nil {
        return nil, fmt.Errorf("failed to create request: %w", err)
    }

    // Add default headers
    for name, value := range c.headers {
        req.Header.Set(name, value)
    }

    // Add per-request headers (these override defaults)
    for name, value := range extraHeaders {
        req.Header.Set(name, value)
    }

    // Set Content-Length if we have a body
    if body != nil {
        req.ContentLength = int64(len(body))
    }

    // Execute the request
    resp, err := c.client.Do(req)
    if err != nil {
        return nil, fmt.Errorf("HTTP request failed: %w", err)
    }
    defer resp.Body.Close()

    // Read the full response body
    responseBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return nil, fmt.Errorf("failed to read response body: %w", err)
    }

    return &AstraClientResponse{
        StatusCode: resp.StatusCode,
        Headers:    headersToMap(resp.Header),
        body:       responseBody,
    }, nil
}

// Get performs an HTTP GET request.
func (c *AstraHTTPClient) Get(url string) (*AstraClientResponse, error) {
    return c.do("GET", url, nil, nil)
}

// GetWithHeaders performs an HTTP GET request with additional headers.
func (c *AstraHTTPClient) GetWithHeaders(url string, headers map[string]string) (*AstraClientResponse, error) {
    return c.do("GET", url, nil, headers)
}

// PostJSON performs an HTTP POST request with a JSON body.
// Automatically sets Content-Type: application/json.
func (c *AstraHTTPClient) PostJSON(url string, data interface{}) (*AstraClientResponse, error) {
    body, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize request body: %w", err)
    }
    return c.do("POST", url, body, map[string]string{
        "Content-Type": "application/json",
    })
}

// Post performs an HTTP POST with a raw body and explicit content type.
func (c *AstraHTTPClient) Post(url, contentType, body string) (*AstraClientResponse, error) {
    return c.do("POST", url, []byte(body), map[string]string{
        "Content-Type": contentType,
    })
}

// PutJSON performs an HTTP PUT request with a JSON body.
func (c *AstraHTTPClient) PutJSON(url string, data interface{}) (*AstraClientResponse, error) {
    body, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize request body: %w", err)
    }
    return c.do("PUT", url, body, map[string]string{
        "Content-Type": "application/json",
    })
}

// PatchJSON performs an HTTP PATCH request with a JSON body.
func (c *AstraHTTPClient) PatchJSON(url string, data interface{}) (*AstraClientResponse, error) {
    body, err := json.Marshal(data)
    if err != nil {
        return nil, fmt.Errorf("failed to serialize request body: %w", err)
    }
    return c.do("PATCH", url, body, map[string]string{
        "Content-Type": "application/json",
    })
}

// Delete performs an HTTP DELETE request.
func (c *AstraHTTPClient) Delete(url string) (*AstraClientResponse, error) {
    return c.do("DELETE", url, nil, nil)
}

// DeleteWithHeaders performs an HTTP DELETE with additional headers.
func (c *AstraHTTPClient) DeleteWithHeaders(url string, headers map[string]string) (*AstraClientResponse, error) {
    return c.do("DELETE", url, nil, headers)
}

// Head performs an HTTP HEAD request (returns headers only, no body).
func (c *AstraHTTPClient) Head(url string) (*AstraClientResponse, error) {
    return c.do("HEAD", url, nil, nil)
}

// ---------------------------------------------------------------------------
// AstraHeaders — a builder for constructing header maps
// ---------------------------------------------------------------------------

// AstraHeaders is a simple builder for constructing header maps.
// In Astra: http.Headers.new().set("Authorization", "Bearer ...").set("Accept", "application/json")
type AstraHeaders struct {
    headers map[string]string
}

// NewAstraHeaders creates a new empty headers builder.
func NewAstraHeaders() *AstraHeaders {
    return &AstraHeaders{headers: make(map[string]string)}
}

// Set adds a header and returns the builder for chaining.
func (h *AstraHeaders) Set(name, value string) *AstraHeaders {
    h.headers[name] = value
    return h
}

// Build returns the completed header map.
func (h *AstraHeaders) Build() map[string]string {
    result := make(map[string]string, len(h.headers))
    for k, v := range h.headers {
        result[k] = v
    }
    return result
}

// ---------------------------------------------------------------------------
// Utility: building Authorization headers
// ---------------------------------------------------------------------------

// BearerToken creates the value for an Authorization: Bearer ... header.
func BearerToken(token string) string {
    return "Bearer " + strings.TrimSpace(token)
}

// BasicAuth creates the value for an Authorization: Basic ... header.
func BasicAuth(username, password string) string {
    import64 := func(s string) string {
        // base64 encode username:password
        // (in real code, use encoding/base64)
        return s // placeholder
    }
    return "Basic " + import64(username+":"+password)
}
```

---

## 10. Middleware Support

Middleware is one of the most powerful patterns in web development. A middleware function wraps a handler, running code before the handler, after the handler, or both. It is how you add cross-cutting concerns — logging, authentication, CORS, rate limiting — without polluting every individual handler.

```mermaid
flowchart LR
    Req["Request"] --> Logger["Logger"]
    Logger --> CORS["CORS"]
    CORS --> Auth["Auth"]
    Auth --> Handler["Handler"]
    Handler --> Auth2["Auth.after"]
    Auth2 --> CORS2["CORS.after"]
    CORS2 --> Logger2["Logger.after"]
    Logger2 --> Resp["Response"]
```

```go
// stdlib/http/middleware.go
// Built-in middleware for the Astra HTTP server.

package astrahttp

import (
    "fmt"
    "strings"
    "time"
)

// Logger is a middleware that prints a log line for each request.
// Format: [METHOD /path] STATUS (duration)
// Example: [GET /users] 200 OK (1.23ms)
func Logger() MiddlewareFunc {
    return func(next HandlerFunc) HandlerFunc {
        return func(req *AstraRequest, res *AstraResponse) {
            start := time.Now()

            // Call the next handler in the chain
            next(req, res)

            // After the handler returns, log the result
            duration := time.Since(start)
            statusText := statusText(res.statusCode)
            fmt.Printf("[%s %s] %d %s (%s)\n",
                req.Method,
                req.Path,
                res.statusCode,
                statusText,
                formatDuration(duration),
            )
        }
    }
}

// CORS adds Cross-Origin Resource Sharing headers to every response.
// allowOrigin is the value for Access-Control-Allow-Origin.
// Use "*" to allow all origins (fine for public APIs, never for auth APIs).
//
// CORS is needed because browsers block JavaScript code on origin A from
// making requests to origin B. The CORS headers tell the browser whether
// to allow the request.
func CORS(allowOrigin string) MiddlewareFunc {
    return func(next HandlerFunc) HandlerFunc {
        return func(req *AstraRequest, res *AstraResponse) {
            // Set CORS headers
            res.writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
            res.writer.Header().Set("Access-Control-Allow-Methods",
                "GET, POST, PUT, PATCH, DELETE, OPTIONS")
            res.writer.Header().Set("Access-Control-Allow-Headers",
                "Content-Type, Authorization, X-Request-Id")
            res.writer.Header().Set("Access-Control-Max-Age", "86400")

            // Handle preflight OPTIONS requests — browsers send these before
            // the actual cross-origin request to check permissions
            if req.Method == "OPTIONS" {
                res.Status(204).Send("")
                return
            }

            // Call the actual handler
            next(req, res)
        }
    }
}

// Auth is a simple JWT Bearer token authentication middleware.
// It reads the Authorization header, validates the token, and either
// calls the next handler or returns 401 Unauthorized.
// In production, you would use a real JWT library.
func Auth(validateToken func(token string) bool) MiddlewareFunc {
    return func(next HandlerFunc) HandlerFunc {
        return func(req *AstraRequest, res *AstraResponse) {
            authHeader := req.Header("Authorization")

            if authHeader == "" {
                res.Error(401, "authentication required")
                return
            }

            if !strings.HasPrefix(authHeader, "Bearer ") {
                res.Error(401, "invalid authorization format, expected: Bearer <token>")
                return
            }

            token := strings.TrimPrefix(authHeader, "Bearer ")
            if !validateToken(token) {
                res.Error(401, "invalid or expired token")
                return
            }

            // Token is valid; continue to the handler
            next(req, res)
        }
    }
}

// RateLimit is a simple in-memory rate limiter.
// It allows up to maxRequests per window per IP address.
// In production, you would use Redis or a distributed counter.
func RateLimit(maxRequests int, windowSeconds int) MiddlewareFunc {
    // Simple in-memory counter (not suitable for multi-instance deployments)
    type counter struct {
        count     int
        resetAt   time.Time
    }
    counters := make(map[string]*counter)
    var mu sync.RWMutex

    return func(next HandlerFunc) HandlerFunc {
        return func(req *AstraRequest, res *AstraResponse) {
            ip := req.IP()
            now := time.Now()

            mu.Lock()
            c, exists := counters[ip]
            if !exists || now.After(c.resetAt) {
                c = &counter{
                    count:   0,
                    resetAt: now.Add(time.Duration(windowSeconds) * time.Second),
                }
                counters[ip] = c
            }
            c.count++
            count := c.count
            mu.Unlock()

            if count > maxRequests {
                res.Header("Retry-After", fmt.Sprintf("%d", windowSeconds))
                res.Error(429, "rate limit exceeded")
                return
            }

            next(req, res)
        }
    }
}

// Recover is a middleware that catches panics in handlers and returns
// a 500 Internal Server Error instead of crashing the server.
// Always add this as the outermost middleware in production.
func Recover() MiddlewareFunc {
    return func(next HandlerFunc) HandlerFunc {
        return func(req *AstraRequest, res *AstraResponse) {
            defer func() {
                if r := recover(); r != nil {
                    fmt.Printf("[Astra HTTP] PANIC in handler for %s %s: %v\n",
                        req.Method, req.Path, r)
                    if !res.committed {
                        res.Error(500, "internal server error")
                    }
                }
            }()
            next(req, res)
        }
    }
}

// ---------------------------------------------------------------------------
// Helper functions for middleware
// ---------------------------------------------------------------------------

func statusText(code int) string {
    texts := map[int]string{
        200: "OK", 201: "Created", 204: "No Content",
        301: "Moved Permanently", 302: "Found", 304: "Not Modified",
        400: "Bad Request", 401: "Unauthorized", 403: "Forbidden",
        404: "Not Found", 405: "Method Not Allowed", 409: "Conflict",
        422: "Unprocessable Entity", 429: "Too Many Requests",
        500: "Internal Server Error", 502: "Bad Gateway",
        503: "Service Unavailable", 504: "Gateway Timeout",
    }
    if t, ok := texts[code]; ok {
        return t
    }
    return "Unknown"
}

func formatDuration(d time.Duration) string {
    if d < time.Millisecond {
        return fmt.Sprintf("%.2fµs", float64(d.Nanoseconds())/1000)
    }
    if d < time.Second {
        return fmt.Sprintf("%.2fms", float64(d.Nanoseconds())/1e6)
    }
    return fmt.Sprintf("%.2fs", d.Seconds())
}
```

---

## 11. The C FFI Bridge

Recall from the runtime chapters that Astra's compiler generates C code that is then compiled with a C compiler (clang/gcc). Astra's standard library packages are Go packages compiled into a C-shared library. The C FFI bridge is the layer that makes the Go types accessible from C-generated code.

This is the same pattern used by the `file` and `json` packages — here we apply it to HTTP.

### The C Header

```c
// runtime/astra_http.h
// C FFI bridge for Astra's HTTP standard library.
// This file is included by the C code that Astra's compiler generates.

#ifndef ASTRA_HTTP_H
#define ASTRA_HTTP_H

#include <stdint.h>
#include <stdlib.h>

// ---------------------------------------------------------------------------
// Opaque handle types
// The actual structs live in Go. C code treats them as opaque pointers.
// ---------------------------------------------------------------------------

typedef struct AstraHTTPServerHandle AstraHTTPServerHandle;
typedef struct AstraHTTPClientHandle AstraHTTPClientHandle;
typedef struct AstraRequestHandle    AstraRequestHandle;
typedef struct AstraResponseHandle   AstraResponseHandle;
typedef struct AstraClientRespHandle AstraClientRespHandle;
typedef struct AstraHeadersHandle    AstraHeadersHandle;

// AstraString is how Go strings are passed to/from C.
// The C caller is responsible for freeing astra_string_free(s) when done.
typedef struct {
    char*  ptr;
    size_t len;
} AstraString;

// A map entry (key-value pair of strings), for passing query/path params.
typedef struct {
    AstraString key;
    AstraString value;
} AstraStringPair;

// AstraStringMap is an array of string pairs.
typedef struct {
    AstraStringPair* pairs;
    size_t           count;
} AstraStringMap;

// ---------------------------------------------------------------------------
// Server lifecycle
// ---------------------------------------------------------------------------

// Create a new HTTP server. Returns an opaque handle.
AstraHTTPServerHandle* astra_http_server_new(void);

// Destroy the server and free all resources.
void astra_http_server_free(AstraHTTPServerHandle* server);

// Register a GET route.
// pattern: null-terminated string, e.g. "/users/:id"
// handler: a C function pointer called for each matching request
void astra_http_server_get(
    AstraHTTPServerHandle* server,
    const char* pattern,
    void (*handler)(AstraRequestHandle*, AstraResponseHandle*)
);

// Register a POST route.
void astra_http_server_post(
    AstraHTTPServerHandle* server,
    const char* pattern,
    void (*handler)(AstraRequestHandle*, AstraResponseHandle*)
);

// Register a PUT route.
void astra_http_server_put(
    AstraHTTPServerHandle* server,
    const char* pattern,
    void (*handler)(AstraRequestHandle*, AstraResponseHandle*)
);

// Register a PATCH route.
void astra_http_server_patch(
    AstraHTTPServerHandle* server,
    const char* pattern,
    void (*handler)(AstraRequestHandle*, AstraResponseHandle*)
);

// Register a DELETE route.
void astra_http_server_delete(
    AstraHTTPServerHandle* server,
    const char* pattern,
    void (*handler)(AstraRequestHandle*, AstraResponseHandle*)
);

// Add middleware. Built-in middleware types:
//   "logger"  — request/response logging
//   "cors"    — CORS headers; extra_arg is the allowed origin
//   "recover" — panic recovery
void astra_http_server_use(
    AstraHTTPServerHandle* server,
    const char* middleware_name,
    const char* extra_arg  // may be NULL
);

// Start the server. This blocks until the server stops.
// Returns 0 on clean shutdown, non-zero on error.
int astra_http_server_listen(AstraHTTPServerHandle* server, int port);

// Serve static files from a directory.
void astra_http_server_static(
    AstraHTTPServerHandle* server,
    const char* url_prefix,
    const char* dir_path
);

// ---------------------------------------------------------------------------
// Request accessors
// ---------------------------------------------------------------------------

// Get the request method ("GET", "POST", etc.)
AstraString astra_request_method(AstraRequestHandle* req);

// Get the request path ("/users/42")
AstraString astra_request_path(AstraRequestHandle* req);

// Get all path parameters
AstraStringMap astra_request_params(AstraRequestHandle* req);

// Get all query parameters
AstraStringMap astra_request_query(AstraRequestHandle* req);

// Get all request headers (lowercased)
AstraStringMap astra_request_headers(AstraRequestHandle* req);

// Get the raw request body
AstraString astra_request_body(AstraRequestHandle* req);

// Get the client IP address
AstraString astra_request_ip(AstraRequestHandle* req);

// ---------------------------------------------------------------------------
// Response builders
// ---------------------------------------------------------------------------

// Set the response status code. Returns the same handle for chaining.
AstraResponseHandle* astra_response_status(AstraResponseHandle* res, int code);

// Set a response header. Returns the same handle for chaining.
AstraResponseHandle* astra_response_header(
    AstraResponseHandle* res,
    const char* name,
    const char* value
);

// Send a JSON response. json_str must be valid JSON.
void astra_response_json(AstraResponseHandle* res, const char* json_str);

// Send a plain text response.
void astra_response_text(AstraResponseHandle* res, const char* body);

// Send an HTML response.
void astra_response_html(AstraResponseHandle* res, const char* body);

// Send a raw string response.
void astra_response_send(AstraResponseHandle* res, const char* body);

// Send a redirect.
void astra_response_redirect(AstraResponseHandle* res, const char* url);

// Send a 404 Not Found response.
void astra_response_not_found(AstraResponseHandle* res);

// ---------------------------------------------------------------------------
// HTTP Client
// ---------------------------------------------------------------------------

// Create a new HTTP client.
AstraHTTPClientHandle* astra_http_client_new(void);

// Destroy the client.
void astra_http_client_free(AstraHTTPClientHandle* client);

// Set the request timeout in milliseconds.
void astra_http_client_set_timeout(AstraHTTPClientHandle* client, int ms);

// Set a default header for all requests.
void astra_http_client_set_header(
    AstraHTTPClientHandle* client,
    const char* name,
    const char* value
);

// Perform a GET request.
// Returns NULL on error; caller must free with astra_client_resp_free().
AstraClientRespHandle* astra_http_client_get(
    AstraHTTPClientHandle* client,
    const char* url,
    char** error_out  // set to error message on failure, caller frees
);

// Perform a POST request with a JSON body.
AstraClientRespHandle* astra_http_client_post_json(
    AstraHTTPClientHandle* client,
    const char* url,
    const char* json_body,
    char** error_out
);

// Perform a PUT request with a JSON body.
AstraClientRespHandle* astra_http_client_put_json(
    AstraHTTPClientHandle* client,
    const char* url,
    const char* json_body,
    char** error_out
);

// Perform a DELETE request.
AstraClientRespHandle* astra_http_client_delete(
    AstraHTTPClientHandle* client,
    const char* url,
    char** error_out
);

// ---------------------------------------------------------------------------
// Client Response accessors
// ---------------------------------------------------------------------------

// Get the HTTP status code.
int astra_client_resp_status(AstraClientRespHandle* resp);

// Get the response body as a string.
AstraString astra_client_resp_text(AstraClientRespHandle* resp);

// True if status is 2xx.
int astra_client_resp_ok(AstraClientRespHandle* resp);

// Destroy the response handle.
void astra_client_resp_free(AstraClientRespHandle* resp);

// ---------------------------------------------------------------------------
// Memory management
// ---------------------------------------------------------------------------

// Free an AstraString that was returned by one of the above functions.
void astra_string_free(AstraString s);

// Free an AstraStringMap that was returned by one of the above functions.
void astra_string_map_free(AstraStringMap m);

#endif // ASTRA_HTTP_H
```

### Why Opaque Handles?

Notice that the C API uses `AstraHTTPServerHandle*` rather than exposing the actual struct fields. This is an important design choice called the **opaque pointer idiom** (or PIMPL — Pointer to IMPLementation).

The C code generated by Astra's compiler never dereferences these pointers directly. It only passes them to the FFI functions. The actual struct layout lives entirely in Go. This means:

1. The C compiler does not need to know Go's internal memory layout.
2. Go can change the implementation without recompiling Astra programs.
3. There is no accidental access to internal state.

This is the same reason why Go's `http.ResponseWriter` is an interface — the actual type behind it changes, but the interface never does.

### The Go CGo Export Layer

```go
// stdlib/http/cgo_exports.go
// CGo export declarations that expose Go functions to C.
// Build with: go build -buildmode=c-shared -o libastra_http.so

package main

// #include "astra_http.h"
import "C"
import (
    "unsafe"
    astrahttp "github.com/astra-lang/astra/stdlib/http"
)

// We need a main function even though this is a library (CGo requirement).
func main() {}

//export astra_http_server_new
func astra_http_server_new() *C.AstraHTTPServerHandle {
    server := astrahttp.NewAstraHTTPServer()
    // Store in a handle registry to keep Go GC from collecting it
    handle := serverRegistry.Store(server)
    return (*C.AstraHTTPServerHandle)(unsafe.Pointer(uintptr(handle)))
}

//export astra_http_server_listen
func astra_http_server_listen(handle *C.AstraHTTPServerHandle, port C.int) C.int {
    server := serverRegistry.Load(uintptr(unsafe.Pointer(handle)))
    if err := server.Listen(int(port)); err != nil {
        return -1
    }
    return 0
}

//export astra_http_client_new
func astra_http_client_new() *C.AstraHTTPClientHandle {
    client := astrahttp.NewAstraHTTPClient()
    handle := clientRegistry.Store(client)
    return (*C.AstraHTTPClientHandle)(unsafe.Pointer(uintptr(handle)))
}

//export astra_http_client_get
func astra_http_client_get(
    clientHandle *C.AstraHTTPClientHandle,
    url *C.char,
    errOut **C.char,
) *C.AstraClientRespHandle {
    client := clientRegistry.Load(uintptr(unsafe.Pointer(clientHandle)))
    resp, err := client.Get(C.GoString(url))
    if err != nil {
        if errOut != nil {
            *errOut = C.CString(err.Error())
        }
        return nil
    }
    handle := clientRespRegistry.Store(resp)
    return (*C.AstraClientRespHandle)(unsafe.Pointer(uintptr(handle)))
}

// ... (similar exports for all other functions)
```

---

## 12. Routing Internals: The Trie

Let us go deeper on how a production trie-based router works. This is the data structure that frameworks like httprouter, Gin, and Echo use internally.

A trie is a tree where each edge represents one segment of a URL path. Nodes at the end of a route hold the handler:

```
After registering these routes:
  GET /
  GET /users
  GET /users/:id
  DELETE /users/:id
  GET /posts/:postId/comments/:commentId
  GET /health

The trie looks like:

root
├── "" (GET) ──────────────── handler for "GET /"
├── "users"
│   ├── "" (GET) ────────────── handler for "GET /users"
│   └── ":id" (param node)
│       ├── "" (GET) ────────── handler for "GET /users/:id"
│       └── "" (DELETE) ──────── handler for "DELETE /users/:id"
├── "posts"
│   └── ":postId" (param node)
│       └── "comments"
│           └── ":commentId" (param node)
│               └── "" (GET) ── handler for "GET /posts/:postId/comments/:commentId"
└── "health"
    └── "" (GET) ──────────── handler for "GET /health"
```

Each node in the trie has:
- A map of literal children (fast O(1) lookup by segment string)
- An optional parameter child (catches `:param` segments)
- An optional wildcard child (catches `*rest` at the end)
- A map of method → handler (since the same path can have GET and POST handlers)

```go
// Simplified trie node structure (for illustration)
type TrieNode struct {
    // Children matched by literal segment: "users", "posts", "health"
    children map[string]*TrieNode

    // Child matched by parameter segment: ":id", ":userId"
    paramChild     *TrieNode
    paramChildName string // "id", "userId" — the name without the ":"

    // Child matched by wildcard: "*filepath"
    wildcardChild     *TrieNode
    wildcardChildName string // "filepath"

    // Handlers for each HTTP method at this path
    handlers map[string]HandlerFunc // "GET" -> handler, "POST" -> handler
}

// Matching GET /users/42/posts:
func (n *TrieNode) match(segments []string, depth int, params map[string]string) HandlerFunc {
    if depth == len(segments) {
        // We have consumed all segments — return the handler for this node
        return n.handlers["GET"] // (passing method as param in real impl)
    }

    seg := segments[depth]

    // 1. Try literal match first (most specific)
    if child, ok := n.children[seg]; ok {
        if h := child.match(segments, depth+1, params); h != nil {
            return h
        }
    }

    // 2. Try parameter match
    if n.paramChild != nil {
        params[n.paramChildName] = seg // capture the param
        if h := n.paramChild.match(segments, depth+1, params); h != nil {
            return h
        }
        delete(params, n.paramChildName) // backtrack if no match deeper
    }

    // 3. Try wildcard match
    if n.wildcardChild != nil {
        params[n.wildcardChildName] = strings.Join(segments[depth:], "/")
        return n.wildcardChild.handlers["GET"]
    }

    return nil
}
```

The key insight: literal matches are tried first, then parameter matches. This means if you have both `/users/me` and `/users/:id` registered, the request `GET /users/me` will correctly match the literal route, not the parameter route.

---

## 13. Registering HTTP in the Astra Compiler

The Astra compiler needs to know about the `http` package: its types, functions, and how to translate Astra `http.*` calls into the corresponding C FFI calls.

```go
// compiler/stdlib_http.go
// Compiler-side registration of the http stdlib package.

package compiler

// RegisterHTTPPackage adds all http package symbols to the compiler's
// symbol table, enabling type checking and code generation for http.* calls.
func (c *Compiler) RegisterHTTPPackage() {
    pkg := c.stdlib.NewPackage("http")

    // --- Types ---
    pkg.AddType("Server", &AstraType{
        Name:   "Server",
        GoType: "astrahttp.AstraHTTPServer",
        Methods: map[string]*AstraFuncType{
            "get":    {Params: []AstraType{StringType, HandlerFuncType}, Return: VoidType},
            "post":   {Params: []AstraType{StringType, HandlerFuncType}, Return: VoidType},
            "put":    {Params: []AstraType{StringType, HandlerFuncType}, Return: VoidType},
            "delete": {Params: []AstraType{StringType, HandlerFuncType}, Return: VoidType},
            "use":    {Params: []AstraType{MiddlewareFuncType}, Return: VoidType},
            "static": {Params: []AstraType{StringType, StringType}, Return: VoidType},
            "listen": {Params: []AstraType{IntType}, Return: VoidType},
        },
        Constructor: "astra_http_server_new",
    })

    pkg.AddType("Client", &AstraType{
        Name:   "Client",
        GoType: "astrahttp.AstraHTTPClient",
        Methods: map[string]*AstraFuncType{
            "get":       {Params: []AstraType{StringType}, Return: ResultType(ClientResponseType)},
            "post_json": {Params: []AstraType{StringType, AnyType}, Return: ResultType(ClientResponseType)},
            "put_json":  {Params: []AstraType{StringType, AnyType}, Return: ResultType(ClientResponseType)},
            "delete":    {Params: []AstraType{StringType}, Return: ResultType(ClientResponseType)},
            "timeout":   {Params: []AstraType{IntType}, Return: ClientType}, // chainable
        },
        Constructor: "astra_http_client_new",
    })

    pkg.AddType("Request", &AstraType{
        Name: "Request",
        Fields: map[string]AstraType{
            "method":  StringType,
            "path":    StringType,
            "query":   MapType(StringType, StringType),
            "params":  MapType(StringType, StringType),
            "headers": MapType(StringType, StringType),
            "body":    StringType,
        },
        Methods: map[string]*AstraFuncType{
            "json": {Params: []AstraType{}, Return: ResultType(MapType(StringType, AnyType))},
            "ip":   {Params: []AstraType{}, Return: StringType},
        },
    })

    pkg.AddType("Response", &AstraType{
        Name: "Response",
        Methods: map[string]*AstraFuncType{
            "status":   {Params: []AstraType{IntType}, Return: ResponseType},       // chainable
            "header":   {Params: []AstraType{StringType, StringType}, Return: ResponseType},
            "json":     {Params: []AstraType{AnyType}, Return: VoidType},
            "text":     {Params: []AstraType{StringType}, Return: VoidType},
            "html":     {Params: []AstraType{StringType}, Return: VoidType},
            "send":     {Params: []AstraType{StringType}, Return: VoidType},
            "redirect": {Params: []AstraType{StringType}, Return: VoidType},
            "not_found":{Params: []AstraType{}, Return: VoidType},
        },
    })

    // --- Functions (free functions in the package namespace) ---
    pkg.AddFunc("logger", &AstraFuncType{
        Params: []AstraType{},
        Return: MiddlewareFuncType,
        CName:  "astra_http_middleware_logger",
    })

    pkg.AddFunc("cors", &AstraFuncType{
        Params: []AstraType{StringType},
        Return: MiddlewareFuncType,
        CName:  "astra_http_middleware_cors",
    })

    c.stdlib.Register("http", pkg)
}
```

### Code Generation for HTTP Calls

When the compiler sees this Astra code:

```astra
let server = http.Server.new()
server.get("/users", myHandler)
server.listen(8080)
```

It generates this C code:

```c
// Generated by astrac for: let server = http.Server.new()
AstraHTTPServerHandle* __astra_var_server = astra_http_server_new();

// Generated by astrac for: server.get("/users", myHandler)
// myHandler is a C function that wraps the Astra closure
astra_http_server_get(__astra_var_server, "/users", __astra_closure_myHandler);

// Generated by astrac for: server.listen(8080)
astra_http_server_listen(__astra_var_server, 8080);
```

The compiler translates each method call into the corresponding FFI function, passing the opaque handle as the first argument. Astra closures (the handler functions) are compiled into C function pointers with the correct signature.

---

## 14. Build Milestone

Enough implementation — let us verify that the HTTP package actually works from Astra code.

Create the file `examples/http_get.as`:

```astra
// examples/http_get.as
// A simple Astra program that makes an HTTP GET request.
// Demonstrates: import http, client creation, GET request, reading response.

import http
import json

fn main() {
    // Create an HTTP client
    let client = http.Client.new()

    println("Fetching data from httpbin.org...")

    // Make a GET request to the httpbin.org/json endpoint.
    // httpbin.org is a public HTTP testing service — it returns
    // a sample JSON response for the /json endpoint.
    let res = client.get("https://httpbin.org/json")?

    // Print the raw HTTP status code
    println("Status: " + res.status.to_string())

    // Print whether the request was successful (2xx status)
    if res.ok() {
        println("Request succeeded!")
    }

    // Print the raw response body as text
    let body = res.text()?
    println("Body:")
    println(body)

    // Parse the body as JSON and access specific fields
    let data = res.json()?
    let slideshow = data["slideshow"]
    println("Title: " + slideshow["title"])
    println("Author: " + slideshow["author"])
}
```

Build and run:

```bash
$ astrac build examples/http_get.as -o http_get
[Astra] Compiling examples/http_get.as...
[Astra] Linking stdlib/http, stdlib/json...
[Astra] Build successful: ./http_get (1.2 MB)

$ ./http_get
Fetching data from httpbin.org...
Status: 200
Request succeeded!
Body:
{
  "slideshow": {
    "author": "Yours Truly",
    "date": "date of publication",
    "slides": [
      {
        "title": "Wake up to WonderWidgets!",
        "type": "all"
      }
    ],
    "title": "Sample Slide Show"
  }
}
Title: Sample Slide Show
Author: Yours Truly
```

Now test the server side. Create `examples/hello_server.as`:

```astra
// examples/hello_server.as
// A minimal Astra HTTP server demonstrating:
// - Route registration (GET, POST)
// - Path parameters
// - JSON responses
// - Middleware

import http

fn main() {
    let server = http.Server.new()

    // Add middleware (runs before every handler)
    server.use(http.logger())
    server.use(http.cors("*"))

    // Simple text response
    server.get("/", fn(req: http.Request, res: http.Response) {
        res.text("Hello from Astra!")
    })

    // JSON response
    server.get("/health", fn(req: http.Request, res: http.Response) {
        res.json({"status": "healthy", "version": "1.0.0"})
    })

    // Path parameter
    server.get("/greet/:name", fn(req: http.Request, res: http.Response) {
        let name = req.params["name"]
        res.json({"message": "Hello, " + name + "!"})
    })

    // POST with JSON body
    server.post("/echo", fn(req: http.Request, res: http.Response) {
        let body = req.json()?
        res.status(201).json({"received": body, "ok": true})
    })

    println("Astra HTTP server running on http://localhost:8080")
    server.listen(8080)
}
```

```bash
$ astrac build examples/hello_server.as -o hello_server
$ ./hello_server
[Astra HTTP] Registered route: GET /
[Astra HTTP] Registered route: GET /health
[Astra HTTP] Registered route: GET /greet/:name
[Astra HTTP] Registered route: POST /echo
[Astra HTTP] Server starting on http://localhost:8080
Astra HTTP server running on http://localhost:8080
```

In another terminal:

```bash
$ curl http://localhost:8080/
Hello from Astra!

$ curl http://localhost:8080/health
{"status":"healthy","version":"1.0.0"}

$ curl http://localhost:8080/greet/world
{"message":"Hello, world!"}

$ curl -X POST http://localhost:8080/echo \
    -H "Content-Type: application/json" \
    -d '{"lang": "Astra", "version": 1}'
{"ok":true,"received":{"lang":"Astra","version":1}}
```

And the server log shows:

```
[GET /] 200 OK (0.12ms)
[GET /health] 200 OK (0.08ms)
[GET /greet/world] 200 OK (0.09ms)
[POST /echo] 201 Created (0.15ms)
```

The HTTP package is working. Astra programs can now speak HTTP.

---

## 15. Exercises

**Exercise 1: HEAD Method Support**
Add support for the HTTP HEAD method to `AstraHTTPServer`. HEAD requests should run the same handler as GET, but the response body should not be sent. Only headers should be returned.

Hint: In `ServeHTTP`, after finding the route match, check if `r.Method == "HEAD"`. If so, wrap the response writer with one that discards body writes but captures headers.

**Exercise 2: Custom Error Handlers**
Add `server.on_not_found(fn(req, res) { ... })` and `server.on_method_not_allowed(fn(req, res) { ... })` to allow Astra programs to customize 404 and 405 responses.

When a path exists but the method does not match, return 405 (Method Not Allowed) instead of 404 (Not Found). The 405 response should include an `Allow` header listing the valid methods for that path.

**Exercise 3: Streaming Responses**
Add `res.stream(fn(write: fn(string) -> void) { ... })` for sending large responses incrementally without buffering the entire body in memory. This uses HTTP chunked transfer encoding.

**Exercise 4: TLS / HTTPS**
Add `server.listen_tls(port: int, cert_file: string, key_file: string)` that starts the server with TLS encryption. Use Go's `http.ListenAndServeTLS`. Test with a self-signed certificate generated by `openssl`.

**Exercise 5: Connection Pooling in the Client**
The current `AstraHTTPClient` creates a new TCP connection for each request (in some configurations). Research Go's `http.Transport` and configure connection pooling: reuse TCP connections for requests to the same host. Benchmark the difference for 100 sequential requests to the same server.

---

## Summary

In this chapter you learned how HTTP works from the wire format up:

- HTTP is a text protocol on top of TCP. Requests have a method, path, headers, and optional body. Responses have a status code, headers, and optional body.
- Status codes signal success (2xx), redirects (3xx), client errors (4xx), and server errors (5xx). Choosing the right one is important.
- The Astra `http` package provides `http.Client` for outbound requests and `http.Server` for inbound request handling, with a clean builder-style API.
- The Go implementation wraps `net/http`, providing `AstraHTTPServer`, `AstraRequest`, `AstraResponse`, `AstraHTTPClient`, and `AstraClientResponse`.
- The router uses path segments to match URLs, supporting literal segments, `:param` segments, and `*wildcard` segments.
- Middleware wraps handlers, enabling cross-cutting concerns like logging, CORS, authentication, and rate limiting.
- The C FFI bridge follows the same opaque-handle pattern established in the runtime and file chapters.

In the next chapter you will put all of this together to build a complete, production-style REST API for a Tasks application — every file, every handler, every piece of middleware, built entirely in Astra.
