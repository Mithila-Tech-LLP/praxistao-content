# Chapter 69: CORS and Security Headers in Go

Building a Go API that a browser can call safely requires understanding two things: **CORS**, which controls which origins can make cross-origin requests, and **security headers**, which tell browsers how to handle your content. Get these wrong and you either block legitimate frontend apps or leave your users exposed.

## Table of Contents

1. [The Same-Origin Policy](#1-the-same-origin-policy)
2. [How CORS Works](#2-how-cors-works)
3. [CORS Middleware in Go](#3-cors-middleware-in-go)
4. [Using the `rs/cors` Library](#4-using-the-rscors-library)
5. [Security Headers Middleware](#5-security-headers-middleware)
6. [Common Pitfalls](#6-common-pitfalls)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. The Same-Origin Policy

The **same-origin policy (SOP)** is a browser security mechanism: a page at `https://app.example.com` cannot make JavaScript `fetch()` calls to `https://api.other.com` by default.

Two URLs share an origin if and only if they have the same **scheme + host + port**:

```
https://app.example.com/page    and   https://app.example.com/api
  → same origin (same scheme, host, port 443)

https://app.example.com         and   http://app.example.com
  → different (scheme differs)

https://app.example.com         and   https://api.example.com
  → different (subdomain differs)

https://app.example.com         and   https://app.example.com:8080
  → different (port differs)
```

SOP protects users: if you are logged into `bank.com`, a malicious page at `evil.com` cannot read your bank data via JavaScript — even though your browser already holds the session cookie.

**CORS (Cross-Origin Resource Sharing)** is the opt-in mechanism that lets `api.example.com` say: "I allow requests from `app.example.com`." It is entirely browser-enforced — a `curl` request is never blocked by CORS.

---

## 2. How CORS Works

### Simple requests

A request is "simple" if it uses `GET`, `HEAD`, or `POST` with only safe headers and `Content-Type` of `text/plain`, `multipart/form-data`, or `application/x-www-form-urlencoded`. The browser sends it directly and checks the response headers.

```
Browser: GET https://api.example.com/products
         Origin: https://app.example.com

Server:  Access-Control-Allow-Origin: https://app.example.com
         (if this header is missing or doesn't match, browser blocks the response)
```

### Preflight requests

For any non-simple request (e.g. `POST` with `Content-Type: application/json`, `PUT`, `DELETE`, custom headers), the browser sends an `OPTIONS` preflight first:

```
Browser: OPTIONS https://api.example.com/orders
         Origin: https://app.example.com
         Access-Control-Request-Method: POST
         Access-Control-Request-Headers: Content-Type, Authorization

Server:  HTTP/1.1 204 No Content
         Access-Control-Allow-Origin: https://app.example.com
         Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
         Access-Control-Allow-Headers: Content-Type, Authorization
         Access-Control-Max-Age: 86400
         (browser caches this preflight for 86400 seconds = 24 hours)

Browser: (now sends the actual POST request)
```

### CORS headers explained

| Header | Direction | Meaning |
|---|---|---|
| `Access-Control-Allow-Origin` | Response | Which origin(s) may read the response |
| `Access-Control-Allow-Methods` | Response (preflight) | Allowed HTTP methods |
| `Access-Control-Allow-Headers` | Response (preflight) | Allowed request headers |
| `Access-Control-Allow-Credentials` | Response | Whether cookies/auth are included |
| `Access-Control-Max-Age` | Response (preflight) | How long to cache the preflight (seconds) |
| `Access-Control-Expose-Headers` | Response | Response headers the browser can read |
| `Origin` | Request | The requesting page's origin |
| `Access-Control-Request-Method` | Request (preflight) | The intended method |
| `Access-Control-Request-Headers` | Request (preflight) | The intended custom headers |

**`Allow-Credentials` and wildcard origins cannot be combined.** If you send `Access-Control-Allow-Credentials: true`, you must specify an exact origin in `Access-Control-Allow-Origin`, not `*`. The browser will block it otherwise.

---

## 3. CORS Middleware in Go

A hand-rolled CORS middleware for chi (or any `net/http` handler chain):

```go
package middleware

import (
    "net/http"
    "slices"
    "strings"
)

// CORSConfig configures which cross-origin requests are allowed.
type CORSConfig struct {
    // AllowedOrigins: list of origins, or ["*"] for any origin.
    // Use exact origins (no wildcards in hostnames) for credentialed requests.
    AllowedOrigins []string

    // AllowedMethods: HTTP methods to allow. Defaults to GET, POST, OPTIONS.
    AllowedMethods []string

    // AllowedHeaders: request headers to allow. Case-insensitive.
    AllowedHeaders []string

    // ExposedHeaders: response headers the browser JS can access.
    ExposedHeaders []string

    // AllowCredentials: send Access-Control-Allow-Credentials: true.
    // Cannot be combined with AllowedOrigins = ["*"].
    AllowCredentials bool

    // MaxAge: preflight cache duration in seconds. 0 means no header sent.
    MaxAge int
}

// DefaultCORSConfig returns a permissive config suitable for development.
func DefaultCORSConfig() CORSConfig {
    return CORSConfig{
        AllowedOrigins: []string{"*"},
        AllowedMethods: []string{"GET", "HEAD", "POST", "PUT", "PATCH", "DELETE", "OPTIONS"},
        AllowedHeaders: []string{"Accept", "Authorization", "Content-Type", "X-Request-ID"},
        MaxAge:         86400,
    }
}

// CORS returns a middleware that handles CORS headers and preflight requests.
func CORS(cfg CORSConfig) func(http.Handler) http.Handler {
    allowedMethodsStr := strings.Join(cfg.AllowedMethods, ", ")
    allowedHeadersStr := strings.Join(cfg.AllowedHeaders, ", ")
    exposedHeadersStr := strings.Join(cfg.ExposedHeaders, ", ")

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")

            // Not a cross-origin request — skip CORS handling
            if origin == "" {
                next.ServeHTTP(w, r)
                return
            }

            // Check if origin is allowed
            allowedOrigin := resolveOrigin(origin, cfg.AllowedOrigins, cfg.AllowCredentials)
            if allowedOrigin == "" {
                // Origin not allowed: don't set any CORS headers.
                // The browser will block the response.
                next.ServeHTTP(w, r)
                return
            }

            // Vary: Origin tells caches that responses differ by Origin
            w.Header().Add("Vary", "Origin")
            w.Header().Set("Access-Control-Allow-Origin", allowedOrigin)

            if cfg.AllowCredentials {
                w.Header().Set("Access-Control-Allow-Credentials", "true")
            }
            if exposedHeadersStr != "" {
                w.Header().Set("Access-Control-Expose-Headers", exposedHeadersStr)
            }

            // Preflight: handle OPTIONS before passing to the next handler
            if r.Method == http.MethodOptions {
                w.Header().Set("Access-Control-Allow-Methods", allowedMethodsStr)
                w.Header().Set("Access-Control-Allow-Headers", allowedHeadersStr)
                if cfg.MaxAge > 0 {
                    w.Header().Set("Access-Control-Max-Age",
                        fmt.Sprintf("%d", cfg.MaxAge))
                }
                w.WriteHeader(http.StatusNoContent) // 204: preflight succeeded
                return
            }

            next.ServeHTTP(w, r)
        })
    }
}

// resolveOrigin checks whether the request origin is allowed.
// Returns the value to use in Access-Control-Allow-Origin,
// or "" if the origin is not permitted.
func resolveOrigin(origin string, allowed []string, credentialed bool) string {
    for _, a := range allowed {
        if a == "*" {
            if credentialed {
                // Wildcard + credentials is invalid per spec: return exact origin instead
                return origin
            }
            return "*"
        }
        if strings.EqualFold(a, origin) {
            return origin // echo back the exact origin (case-insensitive match)
        }
    }
    return "" // not allowed
}
```

**Usage with chi:**

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "yourapp/middleware"
)

func main() {
    r := chi.NewRouter()

    // Development: allow all origins
    r.Use(middleware.CORS(middleware.DefaultCORSConfig()))

    // Production: specific origins only
    r.Use(middleware.CORS(middleware.CORSConfig{
        AllowedOrigins:   []string{"https://app.example.com", "https://admin.example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type"},
        ExposedHeaders:   []string{"X-Request-ID", "X-Total-Count"},
        AllowCredentials: true,
        MaxAge:           3600,
    }))

    r.Get("/api/products", listProducts)
    r.Post("/api/orders", createOrder)

    http.ListenAndServe(":8080", r)
}
```

**Preflight for `DELETE /api/orders/42`:**

```
Request:
  OPTIONS /api/orders/42
  Origin: https://app.example.com
  Access-Control-Request-Method: DELETE
  Access-Control-Request-Headers: Authorization

Response:
  HTTP/1.1 204 No Content
  Access-Control-Allow-Origin: https://app.example.com
  Access-Control-Allow-Methods: GET, POST, PUT, DELETE, OPTIONS
  Access-Control-Allow-Headers: Accept, Authorization, Content-Type
  Access-Control-Allow-Credentials: true
  Access-Control-Max-Age: 3600
  Vary: Origin
```

---

## 4. Using the `rs/cors` Library

For production use, [`rs/cors`](https://github.com/rs/cors) is the standard Go CORS library. It handles edge cases in the CORS spec that are easy to get wrong in a hand-rolled implementation (header canonicalization, `OPTIONS` passthrough, wildcard subdomain matching).

```bash
go get github.com/rs/cors
```

```go
package main

import (
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/rs/cors"
)

func main() {
    r := chi.NewRouter()

    // Attach your routes
    r.Get("/api/products", listProducts)
    r.Post("/api/orders", createOrder)

    // Wrap the router with the cors handler
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"https://app.example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type"},
        ExposedHeaders:   []string{"X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           86400,
        // Debug: true,  // log each CORS decision — useful during development
    })

    handler := c.Handler(r)
    http.ListenAndServe(":8080", handler)
}
```

**`rs/cors` extras over the hand-rolled version:**
- `AllowOriginFunc func(origin string) bool` for dynamic origin checks (e.g. look up in database)
- `AllowOriginRequestFunc` for per-request origin decisions
- `OptionsPassthrough bool` — pass `OPTIONS` to your handler instead of short-circuiting
- Correct handling of `Vary` headers for CDN caches

---

## 5. Security Headers Middleware

Beyond CORS, modern browsers support a suite of security headers that protect users from XSS, clickjacking, MIME sniffing, and data leakage.

```go
package middleware

import (
    "net/http"
)

// SecurityHeadersConfig controls which security headers are set.
type SecurityHeadersConfig struct {
    // HSTS: force HTTPS. maxAge in seconds. includeSubdomains and preload are optional.
    HSTSMaxAge            int
    HSTSIncludeSubdomains bool
    HSTSPreload           bool

    // ContentSecurityPolicy: restrict where resources can be loaded from.
    // An empty string skips the header.
    ContentSecurityPolicy string

    // ReferrerPolicy: controls how much of the URL is sent in Referer headers.
    ReferrerPolicy string

    // PermissionsPolicy: disable browser features you don't use.
    PermissionsPolicy string

    // FrameOptions: "DENY", "SAMEORIGIN", or "" to skip.
    FrameOptions string
}

// DefaultSecurityHeaders returns a safe default config for most REST APIs.
func DefaultSecurityHeaders() SecurityHeadersConfig {
    return SecurityHeadersConfig{
        HSTSMaxAge:            31536000, // 1 year
        HSTSIncludeSubdomains: true,
        HSTSPreload:           false,    // only set preload after testing thoroughly
        ContentSecurityPolicy: "default-src 'none'", // APIs don't serve HTML — block everything
        ReferrerPolicy:        "strict-origin-when-cross-origin",
        PermissionsPolicy:     "camera=(), microphone=(), geolocation=()",
        FrameOptions:          "DENY",
    }
}

// SecurityHeaders returns a middleware that sets security response headers.
func SecurityHeaders(cfg SecurityHeadersConfig) func(http.Handler) http.Handler {
    // Build the HSTS header value once at startup
    hsts := buildHSTS(cfg)

    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            h := w.Header()

            // HSTS: only sent over HTTPS (sending over HTTP is meaningless and can cause issues)
            if hsts != "" && r.TLS != nil {
                h.Set("Strict-Transport-Security", hsts)
            }

            // Prevent MIME type sniffing
            // Without this, a browser may execute a .txt file as JavaScript
            h.Set("X-Content-Type-Options", "nosniff")

            // Clickjacking protection
            if cfg.FrameOptions != "" {
                h.Set("X-Frame-Options", cfg.FrameOptions)
            }

            // Legacy XSS filter (mostly superseded by CSP, but still useful for IE)
            h.Set("X-XSS-Protection", "1; mode=block")

            // Referrer policy: control what URL is sent in Referer header
            if cfg.ReferrerPolicy != "" {
                h.Set("Referrer-Policy", cfg.ReferrerPolicy)
            }

            // CSP: define which sources are valid for scripts, styles, images, etc.
            if cfg.ContentSecurityPolicy != "" {
                h.Set("Content-Security-Policy", cfg.ContentSecurityPolicy)
            }

            // Permissions policy: disable browser features
            if cfg.PermissionsPolicy != "" {
                h.Set("Permissions-Policy", cfg.PermissionsPolicy)
            }

            next.ServeHTTP(w, r)
        })
    }
}

func buildHSTS(cfg SecurityHeadersConfig) string {
    if cfg.HSTSMaxAge == 0 {
        return ""
    }
    s := fmt.Sprintf("max-age=%d", cfg.HSTSMaxAge)
    if cfg.HSTSIncludeSubdomains {
        s += "; includeSubDomains"
    }
    if cfg.HSTSPreload {
        s += "; preload"
    }
    return s
}
```

**What each header does:**

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
  → Browser must use HTTPS for this domain for the next year, even if the user types http://
  → includeSubDomains: applies to *.example.com as well
  → preload: submit to browser preload lists — browsers enforce HTTPS before ever visiting

X-Content-Type-Options: nosniff
  → Don't guess content type — respect the Content-Type header
  → Prevents script injection via a .jpg that is actually JavaScript

X-Frame-Options: DENY
  → This page cannot be loaded in an <iframe> — prevents clickjacking
  → Superseded by CSP's frame-ancestors directive, but still needed for old browsers

X-XSS-Protection: 1; mode=block
  → Old IE/Chrome XSS filter — blocks the page if XSS is detected
  → Modern browsers use CSP instead; this header is mostly legacy

Referrer-Policy: strict-origin-when-cross-origin
  → Same-origin requests: full URL in Referer
  → Cross-origin HTTPS→HTTPS: only origin (no path)
  → Cross-origin HTTPS→HTTP: no Referer at all (don't leak URLs to downgrade)

Content-Security-Policy: default-src 'none'; script-src 'self'; style-src 'self'
  → The most powerful security header — see below

Permissions-Policy: camera=(), microphone=(), geolocation=()
  → Opt out of browser APIs — a compromised script can't silently access the camera
```

**Full router setup combining CORS + security headers:**

```go
package main

import (
    "fmt"
    "net/http"
    "github.com/go-chi/chi/v5"
    "github.com/rs/cors"
    "yourapp/middleware"
)

func main() {
    r := chi.NewRouter()

    // Security headers on every response
    r.Use(middleware.SecurityHeaders(middleware.DefaultSecurityHeaders()))

    // CORS: allow the frontend origin
    c := cors.New(cors.Options{
        AllowedOrigins:   []string{"https://app.example.com"},
        AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
        AllowedHeaders:   []string{"Authorization", "Content-Type", "X-Request-ID"},
        ExposedHeaders:   []string{"X-Request-ID"},
        AllowCredentials: true,
        MaxAge:           3600,
    })
    r.Use(c.Handler)

    r.Get("/api/health", func(w http.ResponseWriter, r *http.Request) {
        fmt.Fprintln(w, `{"status":"ok"}`)
    })
    r.Post("/api/orders", createOrder)

    http.ListenAndServe(":8080", r)
}
```

### CSP for a web application (not just an API)

For a Go service that also serves HTML (e.g. server-rendered templates), the CSP needs more nuance:

```go
// For a web app serving HTML with scripts and styles from the same origin:
csp := strings.Join([]string{
    "default-src 'self'",
    "script-src 'self'",           // only scripts from this origin
    "style-src 'self' 'unsafe-inline'", // allow inline styles (needed by many CSS-in-JS libs)
    "img-src 'self' data: https:",  // images: same origin, data: URIs, and any https
    "font-src 'self'",
    "connect-src 'self' https://api.example.com", // fetch/XHR to these origins
    "frame-ancestors 'none'",       // same as X-Frame-Options: DENY
    "form-action 'self'",           // form submissions only to this origin
    "base-uri 'self'",              // prevent base tag injection
    "upgrade-insecure-requests",    // rewrite http:// subresource URLs to https://
}, "; ")
```

---

## 6. Common Pitfalls

### HSTS locking users out

```
Problem:
  You set HSTS with includeSubDomains.
  You later try to run http://internal.example.com on port 80.
  Every browser that visited example.com will refuse to connect.

Fix:
  1. Test with a short max-age (e.g. 300s = 5 minutes) before committing
  2. Only add includeSubDomains when ALL subdomains serve HTTPS
  3. Never set preload unless you are absolutely certain — removal from preload
     lists takes months and cannot be done in an emergency
```

### CSP breaking things silently

```
Problem:
  You set a strict CSP. Users see a blank page. No error message.
  The browser blocked an inline script or a CDN resource.

Fix:
  1. Start with CSP in report-only mode:
     Content-Security-Policy-Report-Only: default-src 'self'; report-uri /csp-report

  2. Collect violation reports at /csp-report:

     type CSPReport struct {
         Report struct {
             DocumentURI        string `json:"document-uri"`
             BlockedURI         string `json:"blocked-uri"`
             ViolatedDirective  string `json:"violated-directive"`
             OriginalPolicy     string `json:"original-policy"`
         } `json:"csp-report"`
     }

     func cspReportHandler(w http.ResponseWriter, r *http.Request) {
         var report CSPReport
         json.NewDecoder(r.Body).Decode(&report)
         log.Printf("CSP violation: blocked %q on %q (directive: %s)",
             report.Report.BlockedURI,
             report.Report.DocumentURI,
             report.Report.ViolatedDirective)
         w.WriteHeader(http.StatusNoContent)
     }

  3. Fix violations, then switch from Report-Only to enforcing.
```

### Wildcard origin with credentials

```
Problem:
  Access-Control-Allow-Origin: *
  Access-Control-Allow-Credentials: true
  → Browser rejects with CORS error

Fix:
  Echo back the specific requesting origin:
  Access-Control-Allow-Origin: https://app.example.com
  Access-Control-Allow-Credentials: true

  In Go:
    origin := r.Header.Get("Origin")
    if isAllowed(origin) {
        w.Header().Set("Access-Control-Allow-Origin", origin) // exact echo
        w.Header().Set("Access-Control-Allow-Credentials", "true")
        w.Header().Add("Vary", "Origin") // critical: tell caches responses differ by origin
    }
```

### Missing `Vary: Origin`

```
Problem:
  A CDN (or browser) caches the first response for origin A.
  When origin B makes the same request, it gets origin A's CORS headers.
  The browser blocks it.

Fix:
  Always add Vary: Origin whenever you set Access-Control-Allow-Origin.
  This tells caches to store separate entries per Origin header value.
```

---

## Summary

- **Same-origin policy**: browsers block cross-origin JS requests by default. CORS is the opt-in override.
- **Preflight**: browsers send an `OPTIONS` request first for non-simple requests. Your server must respond with appropriate `Access-Control-Allow-*` headers.
- `Access-Control-Allow-Credentials: true` requires an exact origin (not `*`) and `Vary: Origin`.
- For production Go APIs: use `rs/cors` rather than hand-rolling; it handles spec edge cases.
- **Security headers** add defence in depth: `X-Content-Type-Options` stops MIME sniffing, `X-Frame-Options` prevents clickjacking, `HSTS` enforces HTTPS, `CSP` is the most powerful XSS defence.
- Test HSTS with a short `max-age` before committing. Use `Content-Security-Policy-Report-Only` to collect violations before enforcing CSP.
- Always set `Vary: Origin` when your CORS headers change based on the request origin.

---

## Exercises

### Easy

1. Using `curl`, manually test a preflight request against a local Go server:
   ```bash
   curl -X OPTIONS http://localhost:8080/api/orders \
     -H "Origin: https://app.example.com" \
     -H "Access-Control-Request-Method: POST" \
     -H "Access-Control-Request-Headers: Authorization, Content-Type" \
     -v
   ```
   Confirm the response contains all required `Access-Control-*` headers and a `204` status.

2. Implement a `CSPBuilder` struct with chainable methods:
   ```go
   csp := NewCSPBuilder().
       DefaultSrc("'self'").
       ScriptSrc("'self'").
       StyleSrc("'self'", "'unsafe-inline'").
       ImgSrc("'self'", "data:", "https:").
       Build()
   // → "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data: https:"
   ```

3. Write a middleware that sets `Cache-Control: no-store` on all `/api/*` routes (preventing browsers and proxies from caching API responses containing user data) but allows caching on `/static/*` routes.

### Medium

4. Extend the CORS middleware to support **wildcard subdomain matching**: if `AllowedOrigins` contains `"https://*.example.com"`, then `https://app.example.com` and `https://admin.example.com` should both be allowed. Be careful: a naive `strings.HasSuffix` check can be fooled by `https://evil-example.com`.

5. Implement the **CSP violation report endpoint** described in section 6. Parse the JSON report body, log structured fields (blocked URI, violated directive, document URI), and return a `204`. Wire it up at `POST /api/csp-report`. Write a test that posts a sample CSP report payload and verifies the handler logs the expected fields.

6. Create a **security header audit middleware**: in development mode, instead of setting headers, it scans the response headers after the handler runs and logs a warning for every missing security header. Use `httptest.ResponseRecorder` internally to capture the response, inspect headers, then replay to the real `ResponseWriter`.

### Hard

7. Build a **dynamic CORS allowlist backed by a database**: origins are stored in a `allowed_origins` table. The CORS `AllowOriginFunc` looks them up with a short TTL cache (use `sync.Map` with a per-key expiry). Write an admin API endpoint (`POST /admin/cors/origins`) to add new origins, and verify that a newly added origin is reflected in CORS responses within the cache TTL.

8. Implement a **nonce-based CSP** to allow inline scripts without `'unsafe-inline'`. For each request, generate a cryptographically random nonce (`crypto/rand`), inject it into the CSP header as `script-src 'nonce-<base64>'`, and store it in the request context. Write a template helper that reads the nonce from context and renders `<script nonce="...">` tags. Verify that a script tag without the nonce is blocked and one with the correct nonce runs successfully.
