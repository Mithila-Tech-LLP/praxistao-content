# Chapter 29: API Security Testing — The Modern Attack Surface

*Modern applications are APIs first. REST APIs, GraphQL, and gRPC endpoints are now the primary attack surface — and they're often less well-tested than traditional web interfaces.*

---

## API vs Traditional Web Security

| Traditional web | API |
|----------------|-----|
| HTML forms | JSON/XML endpoints |
| Server-side rendering | Client-side SPA |
| Cookie session | JWT / API key |
| Redirect flows | Status codes + JSON |
| Less exposed endpoints | Dozens of documented endpoints |

APIs often have the same vulnerabilities as traditional web apps — but are easier to enumerate because they come with documentation.

---

## API Recon

### Finding API Documentation

```bash
# Common documentation paths
/api/docs
/swagger.json
/openapi.json
/api-docs
/swagger-ui.html
/graphql          # GraphQL endpoint
/graphiql         # GraphQL IDE

# Robots.txt and JS source files
cat app.js | grep -E "api|endpoint|/v[0-9]"

# Wayback machine for old endpoints
curl "http://web.archive.org/cdx/search/cdx?url=target.com/api/*&output=json"
```

### Exploring GraphQL

```bash
# GraphQL introspection — dump full schema
curl -X POST https://target.com/graphql \
    -H "Content-Type: application/json" \
    -d '{"query":"{ __schema { types { name fields { name } } } }"}'

# GraphQL Voyager — visualize schema
# InQL Burp extension — automatic introspection and query generation

# Test for introspection being enabled (often should be disabled in production)
curl -X POST https://target.com/graphql \
    -d '{"query":"{ __typename }"}'
```

---

## REST API Testing Methodology

```
1. Enumerate endpoints (docs, fuzzing, JS analysis)
2. Understand authentication (JWT, API key, OAuth)
3. Test for IDOR (change IDs, try other users' objects)
4. Test authorization (access admin endpoints as user)
5. Test input validation (injection, parameter tampering)
6. Test rate limiting (can you spam requests?)
7. Test for mass assignment (send extra fields)
8. Test verbose error messages (do errors leak info?)
```

### Practical API Testing with curl

```bash
# List all users (should be admin only?)
curl -H "Authorization: Bearer USER_TOKEN" https://api.target.com/v1/admin/users

# Access another user's data
curl -H "Authorization: Bearer MY_TOKEN" https://api.target.com/v1/users/OTHER_ID/profile

# Change another user's email (IDOR + privilege escalation)
curl -X PUT -H "Authorization: Bearer MY_TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"email": "attacker@evil.com"}' \
    https://api.target.com/v1/users/ADMIN_ID

# Test HTTP method differences
curl -X GET https://api.target.com/v1/users    # → 403
curl -X POST https://api.target.com/v1/users   # → 200? (method bypass)
curl -X OPTIONS https://api.target.com/v1/     # CORS configuration

# Test content type switching
curl -X POST -H "Content-Type: text/xml" \
    -d '<user><admin>true</admin></user>' \
    https://api.target.com/v1/users  # XML instead of JSON — different parser!
```

---

## API Key Security

```bash
# Common API key locations to check
# Headers
X-API-Key: xxxxx
Authorization: APIKey xxxxx
Authorization: Bearer xxxxx

# Query parameters (bad practice — shows in logs!)
?api_key=xxxxx
?token=xxxxx
?key=xxxxx

# Test API key exposure in JS source
grep -r "api_key\|apiKey\|API_KEY\|token" *.js

# Test scope: does key have more access than needed?
# Admin key shouldn't be used in frontend
```

---

## OAuth 2.0 Vulnerabilities

```
Authorization Code Flow:
1. User clicks "Login with Google"
2. Server redirects: GET /oauth/auth?client_id=X&redirect_uri=https://app.com/callback&state=ABC
3. Google asks user to authorize
4. Google redirects: https://app.com/callback?code=XYZ&state=ABC
5. Server exchanges code for token: POST /oauth/token {code: XYZ, client_id, client_secret}
6. Token returned, user logged in

Attack vectors:
- Open redirect in redirect_uri: &redirect_uri=https://evil.com (steal the code)
- CSRF: missing state parameter validation
- Token leakage via Referer header
- Scope escalation: requesting more permissions than needed
```

---

## Go: API Security Tester

```go
package main

import (
    "encoding/json"
    "fmt"
    "io"
    "net/http"
    "strings"
    "time"
)

type APITester struct {
    BaseURL string
    Token   string
    Client  *http.Client
}

func NewAPITester(baseURL, token string) *APITester {
    return &APITester{
        BaseURL: baseURL,
        Token:   token,
        Client:  &http.Client{Timeout: 10 * time.Second},
    }
}

func (t *APITester) request(method, path string, body string) (int, string, error) {
    var bodyReader io.Reader
    if body != "" {
        bodyReader = strings.NewReader(body)
    }
    
    req, err := http.NewRequest(method, t.BaseURL+path, bodyReader)
    if err != nil {
        return 0, "", err
    }
    
    req.Header.Set("Authorization", "Bearer "+t.Token)
    req.Header.Set("Content-Type", "application/json")
    
    resp, err := t.Client.Do(req)
    if err != nil {
        return 0, "", err
    }
    defer resp.Body.Close()
    
    respBody, _ := io.ReadAll(resp.Body)
    return resp.StatusCode, string(respBody), nil
}

// Test IDOR: try accessing IDs around a known ID
func (t *APITester) testIDOR(endpoint string, knownID int) {
    fmt.Printf("\n[*] Testing IDOR on %s (known ID: %d)\n", endpoint, knownID)
    
    for _, id := range []int{knownID - 1, knownID + 1, 1, 2, 9999} {
        path := fmt.Sprintf("%s/%d", endpoint, id)
        status, body, err := t.request("GET", path, "")
        if err != nil {
            continue
        }
        
        if status == 200 {
            preview := body
            if len(preview) > 100 {
                preview = preview[:100] + "..."
            }
            fmt.Printf("  [FOUND] ID %d: %s (status %d)\n", id, preview, status)
        }
    }
}

// Test for verbose error messages
func (t *APITester) testErrorMessages(endpoint string) {
    fmt.Printf("\n[*] Testing error verbosity on %s\n", endpoint)
    
    testInputs := []string{"'", "\"", "<script>", "../../../", "null", "{}"}
    
    for _, input := range testInputs {
        path := fmt.Sprintf("%s?id=%s", endpoint, input)
        status, body, err := t.request("GET", path, "")
        if err != nil {
            continue
        }
        
        // Check for common stack trace indicators
        leaks := []string{"stack trace", "at line", "SQL syntax", "ORA-", "SQLSTATE", "Exception"}
        for _, leak := range leaks {
            if strings.Contains(strings.ToLower(body), strings.ToLower(leak)) {
                fmt.Printf("  [LEAK] Input %q → Status %d → Contains %q\n", input, status, leak)
                break
            }
        }
    }
}

// Test HTTP method enumeration
func (t *APITester) testMethods(endpoint string) {
    fmt.Printf("\n[*] Testing HTTP methods on %s\n", endpoint)
    
    methods := []string{"GET", "POST", "PUT", "PATCH", "DELETE", "OPTIONS", "HEAD", "TRACE"}
    for _, method := range methods {
        status, _, err := t.request(method, endpoint, `{}`)
        if err != nil {
            continue
        }
        if status != 405 && status != 404 {
            fmt.Printf("  [%s] → %d\n", method, status)
        }
    }
}

// Check for JWT in response and analyze it
func (t *APITester) analyzeJWT(endpoint string) {
    fmt.Printf("\n[*] Checking for JWT exposure on %s\n", endpoint)
    
    _, body, _ := t.request("GET", endpoint, "")
    
    // Look for JWT pattern: xxx.yyy.zzz
    if idx := strings.Index(body, "eyJ"); idx != -1 {
        end := idx
        for end < len(body) && body[end] != '"' && body[end] != ' ' {
            end++
        }
        jwt := body[idx:end]
        parts := strings.Split(jwt, ".")
        if len(parts) == 3 {
            fmt.Printf("  [JWT FOUND] %s...\n", jwt[:min(20, len(jwt))])
            // Decode header
            if decoded := decodeBase64URL(parts[0]); decoded != "" {
                fmt.Printf("  [JWT Header] %s\n", decoded)
            }
            // Decode payload
            if decoded := decodeBase64URL(parts[1]); decoded != "" {
                fmt.Printf("  [JWT Payload] %s\n", decoded[:min(200, len(decoded))])
            }
        }
    }
}

func decodeBase64URL(s string) string {
    // Add padding
    switch len(s) % 4 {
    case 2: s += "=="
    case 3: s += "="
    }
    import_b64 := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
    _ = import_b64
    // Simplified — use encoding/base64 in real code
    return s
}

func main() {
    tester := NewAPITester("https://vulnerable-api.example.com/api/v1", "YOUR_JWT_TOKEN")
    
    tester.testIDOR("/users", 1001)
    tester.testErrorMessages("/users")
    tester.testMethods("/users")
    tester.analyzeJWT("/users/me")
}

func min(a, b int) int { if a < b { return a }; return b }
```

---

## API Rate Limiting Testing

```bash
# Send 100 requests quickly, see if rate limiting kicks in
for i in $(seq 1 100); do
    curl -s -o /dev/null -w "%{http_code}\n" \
        -H "Authorization: Bearer TOKEN" \
        https://api.target.com/v1/users
done

# With hey (Go HTTP benchmark tool)
hey -n 1000 -c 50 -H "Authorization: Bearer TOKEN" https://api.target.com/v1/users

# If no rate limiting → brute force, data scraping, IDOR mass exploitation
```

---

## Summary

| Test | What to look for |
|------|-----------------|
| Endpoint enumeration | Hidden admin, undocumented API versions |
| IDOR | Access other users' objects by changing IDs |
| Method bypass | GET blocked but POST/PUT works |
| JWT issues | Weak secret, alg:none, claims manipulation |
| Rate limiting | Can you spam requests without consequences? |
| Error verbosity | Stack traces, SQL errors, internal IPs in errors |
| Mass assignment | Extra fields accepted (role, is_admin, balance) |

---

## Exercises

1. Set up a vulnerable REST API (e.g., crAPI, pixi) and run through the full methodology
2. Use `mitmproxy` to intercept a mobile app's API traffic. What endpoints do you find?
3. Write a Go API client that tests all IDOR permutations for a given endpoint and known ID
4. Explore a public GraphQL API (GitHub uses GraphQL). What data is accessible?
