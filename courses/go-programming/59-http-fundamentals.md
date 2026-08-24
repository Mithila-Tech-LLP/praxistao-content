# Chapter 59: HTTP Fundamentals and net/http

Go's `net/http` package is one of the best standard library HTTP implementations in any language. It's production-grade out of the box — the same package powers some of the largest Go services in the world. Understanding it deeply means you can debug any HTTP framework because they're all just wrappers around `net/http`.

## Table of Contents

1. [HTTP Basics Recap](#1-http-basics-recap)
2. [Your First HTTP Server](#2-your-first-http-server)
3. [Request and Response](#3-request-and-response)
4. [ServeMux — Routing](#4-servemux--routing)
5. [HTTP Client](#5-http-client)
6. [Timeouts and Connection Management](#6-timeouts-and-connection-management)
7. [File Server and Static Assets](#7-file-server-and-static-assets)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. HTTP Basics Recap

**HTTP request anatomy:**
```
GET /users/42?include=profile HTTP/1.1
Host: api.example.com
Authorization: Bearer eyJ...
Accept: application/json
Content-Type: application/json

{"key": "value"}
```

**HTTP response anatomy:**
```
HTTP/1.1 200 OK
Content-Type: application/json
Content-Length: 42
X-Request-ID: abc123

{"id": 42, "name": "Alice"}
```

**Status code groups:**
- `2xx` Success: 200 OK, 201 Created, 204 No Content
- `3xx` Redirect: 301 Permanent, 302 Temporary, 304 Not Modified
- `4xx` Client Error: 400 Bad Request, 401 Unauthorized, 403 Forbidden, 404 Not Found, 409 Conflict, 422 Unprocessable Entity, 429 Too Many Requests
- `5xx` Server Error: 500 Internal Server Error, 502 Bad Gateway, 503 Service Unavailable

---

## 2. Your First HTTP Server

```go
package main

import (
    "fmt"
    "log"
    "net/http"
)

func main() {
    // A handler is anything with ServeHTTP(ResponseWriter, *Request):
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
    })

    http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        fmt.Fprint(w, `{"status":"ok"}`)
    })

    log.Println("listening on :8080")
    log.Fatal(http.ListenAndServe(":8080", nil))  // nil = use DefaultServeMux
}
```

**The Handler interface:**
```go
type Handler interface {
    ServeHTTP(ResponseWriter, *Request)
}

// HandlerFunc is a function that implements Handler:
type HandlerFunc func(ResponseWriter, *Request)
func (f HandlerFunc) ServeHTTP(w ResponseWriter, r *Request) { f(w, r) }
```

**Using a custom server (always in production):**
```go
srv := &http.Server{
    Addr:         ":8080",
    Handler:      mux,
    ReadTimeout:  5 * time.Second,
    WriteTimeout: 10 * time.Second,
    IdleTimeout:  120 * time.Second,
}
log.Fatal(srv.ListenAndServe())
```

---

## 3. Request and Response

### Reading the Request
```go
func handleRequest(w http.ResponseWriter, r *http.Request) {
    // Method and URL:
    fmt.Println(r.Method)        // "GET", "POST", etc.
    fmt.Println(r.URL.Path)      // "/users/42"
    fmt.Println(r.URL.RawQuery)  // "include=profile&sort=name"

    // Query parameters:
    q := r.URL.Query()
    include := q.Get("include")   // "profile"
    page := q.Get("page")         // "" if missing

    // Path values (Go 1.22+):
    userID := r.PathValue("id")   // from pattern "/users/{id}"

    // Headers:
    auth := r.Header.Get("Authorization")
    contentType := r.Header.Get("Content-Type")

    // Body (closing is handled by net/http for server handlers).
    // NOTE: the body is a stream — it can only be read ONCE. Pick either
    // io.ReadAll OR json.NewDecoder below, not both:
    body, err := io.ReadAll(r.Body)
    if err != nil {
        http.Error(w, "failed to read body", http.StatusBadRequest)
        return
    }

    // JSON decoding directly from body (alternative to io.ReadAll):
    var payload struct {
        Name  string `json:"name"`
        Email string `json:"email"`
    }
    dec := json.NewDecoder(r.Body)
    dec.DisallowUnknownFields()
    if err := dec.Decode(&payload); err != nil {
        http.Error(w, "invalid JSON: "+err.Error(), http.StatusBadRequest)
        return
    }

    // Form data (call ParseForm first):
    r.ParseForm()
    name := r.FormValue("name")

    // Multipart form data (files):
    r.ParseMultipartForm(10 << 20)  // 10 MB max in memory
    file, header, err := r.FormFile("upload")

    fmt.Println(include, page, userID, auth, contentType, name, file, header)
}
```

### Writing the Response
```go
func handleResponse(w http.ResponseWriter, r *http.Request) {
    // Set headers BEFORE WriteHeader or first Write:
    w.Header().Set("Content-Type", "application/json")
    w.Header().Set("X-Request-ID", "abc123")
    w.Header().Add("Vary", "Accept-Encoding")

    // Write status code (default is 200 if not called):
    w.WriteHeader(http.StatusCreated)  // 201

    // Write body:
    json.NewEncoder(w).Encode(map[string]any{
        "id":   42,
        "name": "Alice",
    })
}

// Common response helpers:
func writeJSON(w http.ResponseWriter, status int, v any) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
    writeJSON(w, status, map[string]string{"error": msg})
}

// Redirect:
http.Redirect(w, r, "/new-location", http.StatusMovedPermanently)

// Simple text error:
http.Error(w, "not found", http.StatusNotFound)
```

### Cookies
```go
// Set a cookie:
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    "abc123",
    Path:     "/",
    HttpOnly: true,
    Secure:   true,  // HTTPS only
    SameSite: http.SameSiteLaxMode,
    MaxAge:   3600,  // Seconds
})

// Read a cookie:
cookie, err := r.Cookie("session")
if err == http.ErrNoCookie {
    http.Error(w, "no session", http.StatusUnauthorized)
    return
}
fmt.Println(cookie.Value)
```

---

## 4. ServeMux — Routing

### Go 1.22 Enhanced ServeMux
```go
// Go 1.22 added method-based routing and path parameters:
mux := http.NewServeMux()

// Method + path pattern:
mux.HandleFunc("GET /users", listUsers)
mux.HandleFunc("POST /users", createUser)
mux.HandleFunc("GET /users/{id}", getUser)
mux.HandleFunc("PUT /users/{id}", updateUser)
mux.HandleFunc("DELETE /users/{id}", deleteUser)

// Path wildcard (matches /files/ and everything after):
mux.HandleFunc("GET /files/{path...}", serveFile)

func getUser(w http.ResponseWriter, r *http.Request) {
    id := r.PathValue("id")  // Extract {id}
    fmt.Fprintf(w, "user: %s", id)
}
```

### Handler vs HandlerFunc
```go
// HandlerFunc — for single functions:
mux.HandleFunc("/foo", myFunc)

// Handle — for types implementing http.Handler:
mux.Handle("/metrics", promhttp.Handler())
mux.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
```

### Sub-routing with http.StripPrefix
```go
api := http.NewServeMux()
api.HandleFunc("GET /users", listUsers)
api.HandleFunc("POST /users", createUser)

// Mount api under /api/v1/:
mux.Handle("/api/v1/", http.StripPrefix("/api/v1", api))
```

---

## 5. HTTP Client

```go
// Default client — don't use in production (no timeouts!):
resp, err := http.Get("https://api.example.com/users")

// Custom client with timeouts:
client := &http.Client{
    Timeout: 10 * time.Second,
    Transport: &http.Transport{
        DialContext: (&net.Dialer{
            Timeout:   5 * time.Second,
            KeepAlive: 30 * time.Second,
        }).DialContext,
        MaxIdleConns:        100,
        IdleConnTimeout:     90 * time.Second,
        TLSHandshakeTimeout: 10 * time.Second,
    },
}

// GET request:
resp, err := client.Get("https://api.example.com/data")
if err != nil {
    return fmt.Errorf("GET failed: %w", err)
}
defer resp.Body.Close()

if resp.StatusCode != http.StatusOK {
    body, _ := io.ReadAll(resp.Body)
    return fmt.Errorf("unexpected status %d: %s", resp.StatusCode, body)
}

var result map[string]any
json.NewDecoder(resp.Body).Decode(&result)

// POST JSON:
payload := map[string]string{"name": "Alice"}
body, _ := json.Marshal(payload)

req, err := http.NewRequestWithContext(ctx, "POST",
    "https://api.example.com/users",
    bytes.NewReader(body),
)
req.Header.Set("Content-Type", "application/json")
req.Header.Set("Authorization", "Bearer "+token)

resp, err = client.Do(req)
// ...

// Helper for JSON requests:
func postJSON(ctx context.Context, client *http.Client, url string, payload, result any) error {
    body, err := json.Marshal(payload)
    if err != nil { return err }

    req, err := http.NewRequestWithContext(ctx, "POST", url, bytes.NewReader(body))
    if err != nil { return err }
    req.Header.Set("Content-Type", "application/json")

    resp, err := client.Do(req)
    if err != nil { return err }
    defer resp.Body.Close()

    if resp.StatusCode >= 400 {
        bodyBytes, _ := io.ReadAll(resp.Body)
        return fmt.Errorf("status %d: %s", resp.StatusCode, bodyBytes)
    }
    return json.NewDecoder(resp.Body).Decode(result)
}
```

---

## 6. Timeouts and Connection Management

```go
// Four timeout layers — each protects against different failures:
srv := &http.Server{
    ReadTimeout:       5 * time.Second,   // Time to read entire request (headers + body)
    ReadHeaderTimeout: 2 * time.Second,   // Time to read headers only
    WriteTimeout:      10 * time.Second,  // Time to write entire response
    IdleTimeout:       120 * time.Second, // Keep-alive connection idle time
}

// Context timeout in handler — for downstream calls:
func slowHandler(w http.ResponseWriter, r *http.Request) {
    ctx, cancel := context.WithTimeout(r.Context(), 3*time.Second)
    defer cancel()

    result, err := fetchFromDB(ctx)
    if err != nil {
        if errors.Is(err, context.DeadlineExceeded) {
            http.Error(w, "upstream timeout", http.StatusGatewayTimeout)
            return
        }
        http.Error(w, "internal error", http.StatusInternalServerError)
        return
    }
    writeJSON(w, http.StatusOK, result)
}
```

### Graceful Shutdown
```go
func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("GET /health", healthHandler)

    srv := &http.Server{
        Addr:         ":8080",
        Handler:      mux,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
    }

    // Start server in background:
    go func() {
        log.Println("starting on :8080")
        if err := srv.ListenAndServe(); err != http.ErrServerClosed {
            log.Fatalf("ListenAndServe: %v", err)
        }
    }()

    // Wait for interrupt signal:
    quit := make(chan os.Signal, 1)
    signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
    <-quit

    // Graceful shutdown — let in-flight requests finish:
    ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer cancel()
    if err := srv.Shutdown(ctx); err != nil {
        log.Fatalf("shutdown: %v", err)
    }
    log.Println("server stopped cleanly")
}
```

---

## 7. File Server and Static Assets

```go
// Serve directory:
http.Handle("/static/", http.StripPrefix("/static/",
    http.FileServer(http.Dir("./web/static"))))

// Serve embedded files (Go 1.16+):
//go:embed web/static
var staticFiles embed.FS

// The embedded FS keeps the "web/static" prefix, so strip it with fs.Sub:
staticFS, _ := fs.Sub(staticFiles, "web/static")
http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(staticFS))))

// Single file download:
http.HandleFunc("/download", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Disposition", `attachment; filename="report.pdf"`)
    http.ServeFile(w, r, "./reports/report.pdf")
})

// In-memory file:
http.HandleFunc("/robots.txt", func(w http.ResponseWriter, r *http.Request) {
    w.Header().Set("Content-Type", "text/plain")
    fmt.Fprint(w, "User-agent: *\nDisallow: /admin/")
})
```

---

## Summary

- `net/http` is production-grade — use a custom `*http.Server` with timeouts set
- `Handler` interface: `ServeHTTP(ResponseWriter, *Request)` — everything is a handler
- Go 1.22 `ServeMux` supports `"METHOD /path/{param}"` patterns natively
- Always `defer resp.Body.Close()` for HTTP client responses
- Set `ReadTimeout`, `WriteTimeout`, and `IdleTimeout` on every server
- Use `r.Context()` to propagate cancellation from client disconnects to downstream calls
- Graceful shutdown: `signal.Notify` + `srv.Shutdown(ctx)` with a deadline

---

## Exercises

### Easy
1. Write an HTTP server with three endpoints: `GET /` returns HTML with a form, `POST /echo` returns the submitted form data as JSON, `GET /health` returns `{"status":"ok","uptime":"Xs"}`. Use `time.Since(startTime)` for uptime.
2. Write an HTTP client function `FetchJSON[T any](ctx context.Context, client *http.Client, url string) (T, error)` that fetches a URL and decodes the JSON body into type T. Handle non-2xx status codes as errors.
3. Demonstrate graceful shutdown: start the server, send a request that sleeps for 5 seconds, trigger SIGINT while it's running, verify the in-flight request completes before the server exits.

### Medium
4. Build a **simple reverse proxy**: any request to `/proxy/{path...}` forwards to `https://httpbin.org/{path}` with the same method, headers, and body. Copy the upstream response status, headers, and body back to the client. Use `httputil.ReverseProxy` or implement manually with the `http.Client`.
5. Implement **request rate limiting** using a token bucket without external libraries. Create a middleware `RateLimit(rps float64) func(http.Handler) http.Handler` that allows `rps` requests per second per IP. Store buckets in a `sync.Map`. Return `429 Too Many Requests` when the bucket is empty.
6. Build a **file upload server**: `POST /upload` accepts multipart/form-data with a file field, validates the file is an image (check magic bytes: JPEG = `FF D8 FF`, PNG = `89 50 4E 47`), saves it to a temp directory with a UUID filename, and returns the download URL. `GET /files/{name}` serves the file.

### Hard
7. **HTTP/2 server push**: Enable HTTP/2 by calling `srv.ListenAndServeTLS` with a self-signed cert. For `GET /`, push the CSS and JS files before sending the HTML. HTTP/2 push hint: check for `http.Pusher` interface: `if pusher, ok := w.(http.Pusher); ok { pusher.Push("/style.css", nil) }`. Note: major browsers (Chrome, Firefox) have removed HTTP/2 push support, so verify the pushed streams with `curl --http2 -v` or `nghttp -v` instead of a browser — and know that in modern production code, `Link: rel=preload` headers (103 Early Hints) have replaced push.
8. Build a **connection pool monitor**: wrap the `http.Transport` with a custom `RoundTripper` that tracks in-flight requests per host, total bytes sent/received, and connection establishment time. Expose these metrics at `GET /debug/connections`. Use `atomic.Int64` for thread-safe counters.
