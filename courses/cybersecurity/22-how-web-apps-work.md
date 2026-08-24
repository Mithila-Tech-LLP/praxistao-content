# Chapter 22: How Web Applications Work — The Attacker's Mental Model

*You can't hack what you don't understand. Before learning web attacks, you need a precise model of how modern web applications are built — every layer is a potential attack surface.*

---

## The Client-Server Model

```
Browser (Client)                    Web Server
     |                                   |
     |------ HTTP GET /login ----------->|
     |                                   | → runs code (PHP/Node/Go)
     |                                   | → queries database
     |                                   |
     |<----- HTTP 200 OK (HTML) ---------|
     |                                   |
     | render HTML                       |
     | execute JavaScript                |
     | make more requests                |
```

The browser is not a passive viewer. It:
1. Renders HTML (structure)
2. Applies CSS (appearance)
3. Executes JavaScript (behavior)
4. Makes additional HTTP requests (APIs, images, scripts)

**Attack surfaces:** Both the server side AND client side.

---

## HTTP Deep Dive

### Request Structure

```
POST /api/login HTTP/1.1           ← Method, path, version
Host: app.example.com              ← Required header
Content-Type: application/json     ← Body format
Authorization: Bearer eyJhbGc...   ← Auth token
Cookie: session=abc123             ← Session identifier
Content-Length: 42                 ← Body size

{"username":"admin","password":"secret"}   ← Body
```

### Response Structure

```
HTTP/1.1 200 OK                    ← Status line
Content-Type: application/json    ← Response format
Set-Cookie: session=xyz789; HttpOnly; Secure; SameSite=Strict
X-Content-Type-Options: nosniff
Strict-Transport-Security: max-age=31536000
Content-Length: 85

{"token":"eyJhbGc...","user":{"id":1,"role":"admin"}}
```

### HTTP Methods

| Method | Purpose | Body | Safe? | Idempotent? |
|--------|---------|------|-------|-------------|
| **GET** | Read resource | No | Yes | Yes |
| **POST** | Create / action | Yes | No | No |
| **PUT** | Replace resource | Yes | No | Yes |
| **PATCH** | Partial update | Yes | No | No |
| **DELETE** | Delete | No | No | Yes |
| **HEAD** | GET without body | No | Yes | Yes |
| **OPTIONS** | CORS preflight | No | Yes | Yes |

**Attack relevance:** CSRF attacks work because browsers automatically send cookies on ANY request, including malicious ones from other sites.

### Status Codes (Security Perspective)

| Code | Meaning | Attack relevance |
|------|---------|-----------------|
| 200 | OK | |
| 301/302 | Redirect | Open redirect attacks |
| 400 | Bad request | May reveal server-side validation |
| 401 | Unauthorized | Need to authenticate |
| 403 | Forbidden | Authenticated but not authorized |
| 404 | Not found | |
| 405 | Method not allowed | Useful for method enumeration |
| 500 | Server error | May leak stack traces |
| 503 | Service unavailable | DoS indicator |

**Pro tip:** 401 vs 403 difference is important:
- `401` = "I don't know who you are" → try authentication bypass
- `403` = "I know who you are, you can't do this" → try privilege escalation

---

## Cookies and Sessions

HTTP is stateless. Every request is independent. **Sessions** solve this.

### Session Flow

```
1. User logs in: POST /login (username + password)
2. Server creates session: session_id = random(128bit)
3. Server stores: sessions["abc123"] = {user_id: 42, role: "admin"}
4. Server sends: Set-Cookie: session=abc123; HttpOnly; Secure
5. Browser stores cookie
6. All future requests: Cookie: session=abc123
7. Server looks up session, knows it's user 42
```

### Cookie Security Attributes

| Attribute | What it does | Attack prevented |
|-----------|-------------|-----------------|
| `HttpOnly` | JavaScript cannot read this cookie | XSS session theft |
| `Secure` | Cookie only sent over HTTPS | Network eavesdropping |
| `SameSite=Strict` | Cookie not sent cross-site | CSRF attacks |
| `SameSite=Lax` | Cookie sent on top-level navigation | Most CSRF |
| `Domain=example.com` | Cookie scope | |
| `Path=/api` | Cookie scope | |
| `Expires=...` | Cookie lifetime | |

**Missing `HttpOnly`:** JavaScript can steal the cookie via XSS → attacker takes over session.
**Missing `Secure`:** Cookie transmitted in plaintext over HTTP → attacker on network intercepts.
**Missing `SameSite`:** Browser sends cookie with cross-site requests → CSRF possible.

---

## How Authentication Works

### Token Types

**Session cookies** (traditional):
```
Client stores: session_id (opaque random string)
Server stores: session data in database or memory
```

**JWT (JSON Web Token)** (modern):
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwicm9sZSI6ImFkbWluIn0.SflKxwRJSMeKKF2QT4fwpMeJf36POk6yJV_adQssw5c

= base64(header) + "." + base64(payload) + "." + signature

header:  {"alg": "HS256", "typ": "JWT"}
payload: {"sub": "1234567890", "role": "admin", "iat": 1516239022}
sig:     HMAC-SHA256(header.payload, secret_key)
```

**JWT attack — algorithm confusion:**
If server accepts `alg: none`, attacker can:
1. Decode JWT
2. Change `role: "user"` to `role: "admin"`
3. Set `alg: none`, remove signature
4. Server accepts it!

**JWT attack — weak secret:**
If HMAC secret is weak (`password`, `secret`, `123456`), attacker can brute-force it offline and forge any token.

---

## Same-Origin Policy (SOP)

The most important browser security policy.

**Origin = protocol + domain + port**
- `https://example.com:443` ← origin
- `https://example.com:443/page` ← same origin
- `http://example.com:443` ← DIFFERENT (protocol differs)
- `https://api.example.com` ← DIFFERENT (subdomain differs)

**Rule:** JavaScript on `https://evil.com` cannot read the response from `https://bank.com`.

**Why it matters:** Without SOP, any website could make requests to your bank and read your balance.

**CORS (Cross-Origin Resource Sharing):** Controlled exceptions to SOP.

```
Browser requests: GET https://api.example.com/data
                  Origin: https://myapp.com

Server responds:  Access-Control-Allow-Origin: https://myapp.com
                  ← now browser allows JavaScript to read this

Or dangerously:   Access-Control-Allow-Origin: *
                  ← any origin can read this
                  
Or very dangerously:
                  Access-Control-Allow-Origin: https://evil.com
                  Access-Control-Allow-Credentials: true
                  ← evil.com can read your data INCLUDING cookies
```

---

## Web Application Architecture

### Traditional (Multi-Page Application)

```
Browser → HTTP GET /products → Server → Database → HTML → Browser renders
```

Every page load: full HTML from server.

### Modern (Single-Page Application)

```
Browser → HTTP GET / → Server → index.html + JavaScript bundle

Then JavaScript takes over:
Browser JavaScript → HTTP GET /api/products → JSON → Update DOM
```

**Security difference:**
- Traditional: Business logic on server, client just displays
- SPA: More client-side logic → more JavaScript code → more client-side vulnerabilities

### API Architecture

```
Frontend (React/Vue/Angular)
         ↕ REST API (JSON)
Backend (Go/Node/Python)
         ↕ SQL/ORM
Database (PostgreSQL/MySQL)
```

**What attackers target:**
- **API endpoints:** Are they all authenticated? Do they validate input?
- **Business logic:** Can I call `/api/admin/deleteUser` as a regular user?
- **Data access:** Can I access another user's data by changing an ID?

---

## HTTPS and TLS

HTTPS = HTTP over TLS (Transport Layer Security). TLS provides:
1. **Encryption:** Nobody can read the traffic in transit
2. **Authentication:** Server proves its identity via certificate
3. **Integrity:** Data cannot be tampered with

### TLS Handshake (Simplified)

```
1. Client Hello: "I support TLS 1.3, here are my ciphers"
2. Server Hello: "Let's use TLS_AES_128_GCM_SHA256"
3. Server Certificate: "Here's my certificate proving I'm example.com"
4. Client: Verify certificate against trusted CAs
5. Key exchange: establish shared secret
6. Both sides: derive session keys
7. Encrypted data flows
```

### Certificate Validation

```
Certificate chain:
Your browser trusts:  Root CA (built into OS/browser)
Root CA signed:       Intermediate CA
Intermediate CA signed: example.com certificate
```

**SSL stripping attack:** Attacker intercepts HTTP (before redirect to HTTPS), user never gets encrypted connection. Prevented by HSTS (`Strict-Transport-Security` header).

**Certificate transparency:** All certificates are logged publicly. Attackers look for newly issued certificates on `crt.sh` to discover subdomains.

---

## The Browser's Security Model

### Content Security Policy (CSP)

Tells browser what scripts are allowed to run:

```
Content-Security-Policy: 
    default-src 'self';              # only own domain
    script-src 'self' cdn.com;       # scripts from own domain + cdn
    style-src 'self' 'unsafe-inline'; # styles + inline
    img-src *;                        # images from anywhere
    frame-ancestors 'none';           # can't be iframed (clickjacking)
```

**CSP misconfigurations attackers look for:**
- `'unsafe-inline'` → inline scripts allowed → XSS still works
- `'unsafe-eval'` → eval() allowed → XSS via eval
- Wildcards: `script-src *` → any domain → load malicious scripts

### HTTP Security Headers

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Frame-Options: DENY
X-Content-Type-Options: nosniff
Referrer-Policy: no-referrer-when-downgrade
Permissions-Policy: geolocation=(), camera=()
```

Use `curl -I https://example.com` to check which headers are present. Missing headers = potential vulnerabilities.

---

## Inspecting Web Apps Like an Attacker

### Browser DevTools

- **Network tab:** See every HTTP request/response, headers, cookies, timing
- **Application tab:** Cookies, LocalStorage, SessionStorage, Service Workers
- **Console:** Execute JavaScript, see errors
- **Sources:** Read JavaScript source, set breakpoints

### Burp Suite — The Pen Tester's Proxy

Burp Suite sits between your browser and the server, intercepting every request.

```
Browser → Burp Proxy → Target Server
                ↑
        You see and modify
        every request here
```

Key features:
- **Intercept:** Pause requests, modify, forward
- **Repeater:** Resend modified requests
- **Intruder:** Fuzzing/brute force
- **Scanner:** Automated vulnerability detection (Pro only)
- **Decoder:** Base64/URL encode-decode

**Setup:** Set browser proxy to `127.0.0.1:8080`, install Burp's CA certificate.

### Manual Recon Steps

```bash
# 1. Get the IP and initial info
dig example.com
whois example.com
curl -I https://example.com   # headers

# 2. Find subdomains
# (use subfinder, amass, or brute force)

# 3. Spider the site
# (use Burp Spider or wget --mirror)

# 4. Check robots.txt and sitemap
curl https://example.com/robots.txt
curl https://example.com/sitemap.xml

# 5. Look for common admin paths
for path in admin administrator wp-admin login dashboard api docs; do
    curl -s -o /dev/null -w "%{http_code} $path\n" https://example.com/$path
done

# 6. Check JavaScript files for secrets
curl -s https://example.com/app.js | grep -E "api_key|secret|password|token"
```

---

## Summary

| Component | What it does | Attack target |
|-----------|-------------|--------------|
| HTTP | Communication protocol | Request manipulation, injection |
| Cookies | State management | Session theft, CSRF |
| JWT | Stateless auth tokens | Algorithm confusion, weak secrets |
| SOP | Browser isolation | Bypassed by CORS misconfig |
| TLS | Encryption + auth | Stripping, cert manipulation |
| CSP | Script whitelisting | Misconfiguration enables XSS |

---

## Exercises

1. Use browser DevTools Network tab to capture a login request. What headers does the server send back? Is there a `Set-Cookie` with `HttpOnly`?
2. Find a site without `X-Frame-Options` — could it be clickjacked?
3. Decode a JWT from any website. What claims does it contain? (Use jwt.io)
4. Use `curl -I` to check security headers on 5 websites. Which ones are missing?
5. Set up Burp Suite, configure your browser, and capture a request. Modify a parameter and see the response change.
