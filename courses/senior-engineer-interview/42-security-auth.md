# Chapter 42: Security — OWASP Top 10, JWT, OAuth2 & Authorization

Security is tested in two ways at senior interviews: knowledge of common vulnerabilities (what they are, how to prevent them) and practical authentication/authorization design. Both are required for senior roles at companies that handle user data.

## Table of Contents

1. [OWASP Top 10 — Critical Vulnerabilities](#1-owasp-top-10--critical-vulnerabilities)
2. [Sessions vs JWT — The Right Choice](#2-sessions-vs-jwt--the-right-choice)
3. [OAuth2 and OpenID Connect](#3-oauth2-and-openid-connect)
4. [Authorization — RBAC and ABAC](#4-authorization--rbac-and-abac)
5. [Secrets Management](#5-secrets-management)
6. [Security in Go](#6-security-in-go)
7. [Interview Questions & Model Answers](#7-interview-questions--model-answers)
8. [Summary](#summary)

---

## 1. OWASP Top 10 — Critical Vulnerabilities

### A01: Broken Access Control

```
Vulnerability: Users access resources they shouldn't.

Examples:
  GET /api/invoices/123     (user sees another user's invoice)
  PUT /api/admin/settings   (regular user can change admin settings)
  Path traversal: GET /files?name=../../etc/passwd

Fix:
  - Check authorization on EVERY request, not just login
  - Use parameterized authorization checks (not URL-based)
  - Principle of least privilege: users only see their own data

// In Go:
func getInvoice(w http.ResponseWriter, r *http.Request) {
    invoiceID := chi.URLParam(r, "id")
    userID := getUserIDFromToken(r)
    
    invoice, err := db.GetInvoice(r.Context(), invoiceID)
    if err != nil { ... }
    
    // ALWAYS verify ownership:
    if invoice.OwnerID != userID {
        http.Error(w, "forbidden", http.StatusForbidden)
        return
    }
}
```

### A02: Cryptographic Failures (formerly "Sensitive Data Exposure")

```
Vulnerability: Sensitive data stored or transmitted insecurely.

Examples:
  - Passwords stored as MD5 or SHA1 (easily cracked)
  - Credit card numbers stored in plaintext
  - HTTP (not HTTPS) for login pages
  - Weak encryption keys

Fix:
  - Passwords: bcrypt, scrypt, or Argon2 (NEVER MD5/SHA1/SHA256 directly)
  - Sensitive data: AES-256-GCM encryption at rest
  - Always HTTPS (HSTS header)
  - Store only what you need (minimize sensitive data)
  
// Go: password hashing with bcrypt
import "golang.org/x/crypto/bcrypt"

func hashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    return string(hash), err
}

func verifyPassword(hash, password string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}
```

### A03: Injection (SQL, NoSQL, Command)

```
SQL Injection:
  Attacker input: ' OR '1'='1
  Query becomes: SELECT * FROM users WHERE username = '' OR '1'='1' -- returns all rows!
  
Fix: ALWAYS use parameterized queries
  // BAD (vulnerable):
  db.Query("SELECT * FROM users WHERE email = '" + email + "'")
  
  // GOOD (safe):
  db.QueryContext(ctx, "SELECT * FROM users WHERE email = $1", email)

Command Injection:
  If user input is passed to exec.Command, attacker can inject shell commands
  
  // BAD:
  exec.Command("sh", "-c", "echo " + userInput)
  
  // GOOD: never pass user input to shell; use argument lists
  exec.Command("echo", userInput) // userInput is not interpreted as shell
```

### A07: Identification and Authentication Failures

```
Examples:
  - No brute-force protection on login (attackers can try millions of passwords)
  - Session tokens that don't expire
  - Passwords stored without salt (rainbow table attacks)
  - Missing 2FA for privileged accounts
  
Fix:
  - Rate limiting on login (max 5 attempts per minute per IP)
  - bcrypt with cost factor ≥ 10 (makes brute-force expensive)
  - Short session TTLs + refresh tokens
  - Require 2FA for admin/privileged actions
  - Never log passwords, tokens, or sensitive input
```

### A09: Security Logging and Monitoring Failures

```
Fix: log all authentication events with enough context to investigate
  - Login successes and failures (with IP, user agent, timestamp)
  - Access to sensitive resources
  - Privilege escalation attempts
  
// What to log:
{
  "event": "login_failed",
  "email": "alice@example.com",  // NOT the password!
  "ip": "192.168.1.1",
  "user_agent": "Mozilla/5.0...",
  "timestamp": "2024-01-15T10:30:00Z"
}

// What NOT to log:
{
  "password": "hunter2",          // NEVER
  "credit_card": "4111...",       // NEVER
  "session_token": "eyJhbGci..."  // NEVER
}
```

### CSRF (Cross-Site Request Forgery)

```
Attack: attacker tricks a browser into making a request to your site using the user's session cookie

Example: user logged into bank.com
  Attacker sends email with link to evil.com
  evil.com has: <form action="https://bank.com/transfer" method="POST">
                <input name="to" value="attacker_account">
                Browser sends the request WITH bank.com cookie

Fix:
  1. CSRF tokens: include a secret token in forms that is checked server-side
     Only your forms have the valid token, not attacker's forms
  
  2. SameSite cookies: Set-Cookie: session=abc; SameSite=Strict
     Browser won't send cookies on cross-site requests

  3. Verify Origin/Referer header for state-changing operations
```

### XSS (Cross-Site Scripting)

```
Attack: attacker injects JavaScript that runs in other users' browsers

Example: attacker submits a comment: <script>steal(document.cookie)</script>
  If stored without escaping, every user who sees the comment runs the script

Types:
  Stored XSS: malicious script persisted in database
  Reflected XSS: script in URL parameter reflected back in response
  DOM XSS: script injected into DOM by JavaScript

Fix:
  - Escape user content when rendering HTML: < → &lt; > → &gt;
  - Content Security Policy (CSP) header: limits which scripts can execute
  - Never use innerHTML with user content; use textContent
  - React/Vue: automatic HTML escaping by default (use dangerouslySetInnerHTML cautiously)
```

---

## 2. Sessions vs JWT — The Right Choice

### Sessions (Server-Side)

```
Flow:
  1. Login: server validates credentials, creates session in database/Redis
     session_id = random 32-byte token
     Redis: SET session:{session_id} {user_id, role, exp} EX 86400
  2. Server sends: Set-Cookie: session_id=abc; HttpOnly; Secure; SameSite=Strict
  3. Browser sends cookie with every request
  4. Server looks up session_id in Redis → gets user context

Pros:
  - Instant logout: delete session from Redis
  - Revoke specific sessions (logout from all devices)
  - Server controls session validity

Cons:
  - Requires shared storage (Redis) accessible to all servers
  - Additional Redis lookup per request (~0.5ms)
```

### JWT (JSON Web Token)

```
Structure: header.payload.signature
  header: {"alg": "HS256", "typ": "JWT"}
  payload: {"sub": "user123", "role": "admin", "exp": 1700000000}
  signature: HMAC-SHA256(base64url(header) + "." + base64url(payload), secret)

Flow:
  1. Login: server creates JWT, signs it with secret key
  2. Client stores JWT (localStorage or cookie)
  3. Client sends JWT in Authorization: Bearer <token> header
  4. Server verifies signature (no database lookup!) → extracts user context

Pros:
  - Stateless: no database lookup needed for auth
  - Works across microservices (all share the secret/public key)
  - Self-contained: payload carries user info

Cons:
  - Cannot revoke before expiry (once issued, valid until exp)
  - JWT payload is BASE64 encoded, not encrypted — don't put secrets in it
  - If secret is compromised, all tokens are compromised

The revocation problem:
  - Short expiry (15 min) + refresh tokens: logout invalidates refresh token
  - Blocklist: store revoked JWT IDs in Redis. Check on each request.
    (Reintroduces a database lookup — reduces the stateless benefit)
```

### Which to Choose?

```
Sessions: for web apps where server control over auth state is needed
  - Medical, financial apps where instant revocation is required
  - When "logout from all devices" must work immediately

JWT: for microservices (service-to-service auth), mobile apps, APIs
  - When services don't share a database
  - When you can tolerate up to N-minute stale auth state
  
In practice: many apps combine both:
  Short-lived JWT (15 min) for stateless requests
  Long-lived refresh token (30 days) stored in HttpOnly cookie
  Logout invalidates refresh token only; JWT expires naturally
```

---

## 3. OAuth2 and OpenID Connect

### OAuth2 — Authorization (not Authentication)

```
OAuth2 lets a user grant a third-party app limited access to their resources
WITHOUT sharing their password.

Example: Trello wants to access your Google Calendar
  1. User clicks "Connect Google Calendar" on Trello
  2. Trello redirects to Google: 
     https://accounts.google.com/oauth2/authorize?
       client_id=trello_id&
       redirect_uri=https://trello.com/callback&
       scope=calendar.readonly&
       response_type=code&
       state=random_csrf_token
  3. User sees Google's consent screen: "Trello wants to read your calendar"
  4. User approves → Google redirects to trello.com/callback?code=authcode&state=...
  5. Trello exchanges code for access_token + refresh_token (server-to-server)
  6. Trello uses access_token to call Google Calendar API on user's behalf

Key parties:
  Resource Owner: the user
  Client: Trello (the app requesting access)
  Authorization Server: Google accounts
  Resource Server: Google Calendar API
```

### OpenID Connect (OIDC) — Authentication on Top of OAuth2

```
OAuth2 is for authorization (access to APIs).
OIDC adds authentication (who is the user?).

OIDC adds:
  - ID Token (JWT): proves who the user is
  - /userinfo endpoint: get user profile (email, name, picture)
  - Standard scopes: "openid", "email", "profile"
  
When you see "Login with Google":
  That's OpenID Connect — Google proves identity to your app
  You get an ID token with the user's Google profile
  Your app creates a local account or finds existing one by email
```

---

## 4. Authorization — RBAC and ABAC

### RBAC (Role-Based Access Control)

```
Users are assigned roles. Roles have permissions.
  User → Role(s) → Permissions

Example:
  Role: "admin"      → permissions: ["read:all", "write:all", "delete:all"]
  Role: "editor"     → permissions: ["read:own", "write:own"]
  Role: "viewer"     → permissions: ["read:own"]
  
  User Alice → ["admin", "editor"]
  User Bob   → ["viewer"]

Implementation:
  Check at request time: user.HasPermission("write:invoices")
```

### ABAC (Attribute-Based Access Control)

```
Access based on attributes of user, resource, and environment.
More fine-grained than RBAC.

Policy example:
  "Users can read invoices IF invoice.customer_id == user.customer_id"
  "Managers can approve expenses IF expense.amount <= manager.approval_limit"
  "Admins can delete users IF request is from corporate IP AND is business hours"

Implementation: policy engine (OPA — Open Policy Agent, or custom)
```

### Open Policy Agent (OPA)

```rego
// OPA policy file (Rego language):
package invoices

default allow = false

allow {
    input.method == "GET"
    input.path == ["api", "invoices", invoice_id]
    input.user.customer_id == data.invoices[invoice_id].customer_id
}

allow {
    input.user.role == "admin"
}
```

---

## 5. Secrets Management

```
NEVER hardcode secrets in code or config files:
  db_password = "hunter2"  // will appear in git history FOREVER
  
Where to store secrets:
  Development: .env file (gitignored), or 1Password/Bitwarden
  Production: AWS Secrets Manager, HashiCorp Vault, GCP Secret Manager
  
HashiCorp Vault:
  Secrets stored encrypted, access controlled by policies
  Dynamic secrets: Vault generates a temporary DB password on demand
    (expires when lease expires — no long-lived credentials)
  Audit log of all secret accesses
  
In Go:
  // At startup, fetch secrets from Vault or AWS Secrets Manager
  secret, err := secretsManager.GetSecretValue(ctx, "/prod/db/password")
  // Use in-memory, never write to disk or logs
```

---

## 6. Security in Go

```go
// Cryptographically secure random token generation:
import "crypto/rand"

func generateSecureToken(n int) (string, error) {
    bytes := make([]byte, n)
    if _, err := rand.Read(bytes); err != nil {
        return "", err
    }
    return base64.URLEncoding.EncodeToString(bytes), nil
}

// Constant-time comparison to prevent timing attacks:
import "crypto/subtle"

// BAD: timing attack possible (comparison stops at first mismatch)
if token == expectedToken { /* ... */ }

// GOOD: always compares all bytes regardless of mismatch
if subtle.ConstantTimeCompare([]byte(token), []byte(expectedToken)) == 1 { /* ... */ }

// JWT with RS256 (asymmetric) — better than HS256 for microservices:
// Private key: only auth service signs tokens
// Public key: all services verify tokens (no need to share the secret)

// HTTPS-only cookie:
http.SetCookie(w, &http.Cookie{
    Name:     "session",
    Value:    sessionID,
    HttpOnly: true,   // not accessible via JavaScript
    Secure:   true,   // only sent over HTTPS
    SameSite: http.SameSiteStrictMode, // CSRF protection
    MaxAge:   86400,
})
```

---

## 7. Interview Questions & Model Answers

**Q: What is the difference between authentication and authorization?**

"Authentication is proving who you are — 'I am Alice, and here's my password/token to prove it.' Authorization is deciding what you're allowed to do — 'Alice is authenticated, but is she allowed to delete this invoice?' Both are separate concerns. JWT handles authentication (the token proves identity). RBAC or ABAC handles authorization (the system decides what Alice can do with that identity). A common mistake is combining them: just because someone is authenticated doesn't mean they should be authorized for everything."

**Q: How would you prevent SQL injection in Go?**

"Always use parameterized queries — never string concatenation. The database/sql package in Go supports parameterized queries via $1, $2 placeholders (PostgreSQL) or ? (MySQL). The SQL driver handles escaping — the user input is never interpreted as SQL syntax. I'd also use an ORM like GORM or sqlx which enforces parameterized queries by default. Code review checks and static analysis tools (go vet, sqlvet) can catch unsafe query construction. Additionally, the database user should have minimum required permissions — an app user should have INSERT/SELECT/UPDATE but not DROP TABLE."

---

## Summary

- **OWASP Top 10:** broken access control is #1 — always verify authorization per request.
- SQL injection: parameterized queries, always.
- XSS: escape user content. CSRF: SameSite cookies + CSRF tokens.
- **Passwords:** bcrypt/Argon2 with salt. NEVER MD5/SHA1.
- **Sessions:** server-side state, instant revocability. **JWT:** stateless, no revocation without blocklist.
- **OAuth2:** authorization (third-party API access). **OIDC:** authentication (who is the user).
- **RBAC:** roles → permissions. **ABAC:** attribute-based policies, more fine-grained.
- **Secrets:** AWS Secrets Manager / HashiCorp Vault in production. Never in code.
- Constant-time comparison for secrets. HttpOnly + Secure + SameSite for cookies.
