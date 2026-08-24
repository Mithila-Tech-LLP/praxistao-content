# Chapter 56: TLS/SSL — Securing Data in Transit

*TLS (Transport Layer Security) is what makes HTTPS secure. It provides encryption, authentication, and integrity for network communication. Every web developer, sysadmin, and security professional needs to understand it deeply.*

---

## TLS vs SSL

```
SSL 2.0 (1995) — broken, never use
SSL 3.0 (1996) — broken (POODLE), never use
TLS 1.0 (1999) — deprecated 2020
TLS 1.1 (2006) — deprecated 2020
TLS 1.2 (2008) — widely used, still acceptable
TLS 1.3 (2018) — current, significantly improved
```

---

## TLS Handshake (TLS 1.2)

```
Client                          Server
  |                               |
  |  ClientHello                  |
  |  (TLS version, random,        |
  |   cipher suites supported)    |
  |──────────────────────────────►|
  |                               |
  |  ServerHello                  |
  |  (chosen cipher suite,        |
  |   server random)              |
  |◄──────────────────────────────|
  |                               |
  |  Certificate (server's cert)  |
  |◄──────────────────────────────|
  |                               |
  |  Client verifies cert         |
  |  (checks against CA, etc.)    |
  |                               |
  |  ClientKeyExchange            |
  |  (pre-master secret,          |
  |   encrypted with server's key)|
  |──────────────────────────────►|
  |                               |
  |  Both derive session keys     |
  |  from pre-master secret       |
  |  + client random + server rand|
  |                               |
  |  ChangeCipherSpec             |
  |──────────────────────────────►|
  |                               |
  |  Finished (encrypted)         |
  |──────────────────────────────►|
  |                               |
  |  ChangeCipherSpec + Finished  |
  |◄──────────────────────────────|
  |                               |
  |  Encrypted application data   |
  |◄─────────────────────────────►|
```

### TLS 1.3 Improvements

```
- Handshake: 1 round-trip (vs 2 in 1.2)
- Forward Secrecy: mandatory (DHE/ECDHE only)
- Removed: RC4, MD5, SHA-1, DH <2048, export ciphers
- Encrypted: more of the handshake is encrypted
- 0-RTT resumption: resume without round-trip (with replay attack risk)
```

---

## Cipher Suites

```
TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384

Breaking it down:
TLS           — protocol
ECDHE         — key exchange (Elliptic Curve Diffie-Hellman Ephemeral)
RSA           — authentication (verify server identity)
AES_256_GCM   — symmetric encryption (bulk data)
SHA384        — MAC / PRF (integrity)

Good cipher suites (strong):
TLS_AES_256_GCM_SHA384         (TLS 1.3)
TLS_CHACHA20_POLY1305_SHA256   (TLS 1.3)
TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384 (TLS 1.2)

Bad cipher suites (weak):
TLS_RSA_WITH_AES_128_CBC_SHA   — no forward secrecy
TLS_RSA_WITH_RC4_128_MD5       — RC4 is broken
TLS_RSA_WITH_3DES_EDE_CBC_SHA  — SWEET32 attack
```

---

## TLS Attacks

### BEAST (2011)

```
- Affects TLS 1.0, CBC mode
- Block boundary attack on predictable IV
- Fix: TLS 1.1+, RC4 (then RC4 was broken too), ECDHE
```

### POODLE (2014)

```
- Affects SSL 3.0, CBC mode
- Padding oracle attack
- Fix: Disable SSL 3.0 entirely
```

### Heartbleed (2014) — CVE-2014-0160

```
- OpenSSL bug, not TLS protocol flaw
- Heartbeat extension: "send me X bytes back"
- Bug: sent back X bytes from memory, including private keys
- Impact: leaked server private keys from memory

Testing for it:
sslscan old-server.example.com
# Shows: Heartbleed: vulnerable
```

### BEAST / CRIME / BREACH

```
CRIME (2012): Compression + TLS → side channel on compressed secrets
BREACH (2013): HTTP compression leaks secrets via repeated requests
Fix: Disable TLS/HTTP compression for sensitive data
```

### Downgrade Attacks

```
FREAK: force server to use export-grade (weak) cryptography
Logjam: downgrade DH key exchange to 512-bit
DROWN: use SSLv2 to attack TLS sessions

Fix: Disable all legacy protocols, use strong ciphers only
```

---

## Configuring TLS Securely

### nginx

```nginx
server {
    listen 443 ssl http2;
    
    ssl_certificate /etc/ssl/certs/server.crt;
    ssl_certificate_key /etc/ssl/private/server.key;
    
    # Only TLS 1.2 and 1.3
    ssl_protocols TLSv1.2 TLSv1.3;
    
    # Strong cipher suites only
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA384:ECDHE-RSA-CHACHA20-POLY1305:ECDHE-RSA-AES128-GCM-SHA256;
    ssl_prefer_server_ciphers off;  # TLS 1.3 clients choose better anyway
    
    # DH params for DHE
    ssl_dhparam /etc/ssl/dhparam.pem;
    ssl_ecdh_curve X25519:secp384r1;
    
    # Session settings
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 1d;
    ssl_session_tickets off;  # disable for forward secrecy
    
    # HSTS
    add_header Strict-Transport-Security "max-age=63072000; includeSubDomains; preload";
    add_header X-Frame-Options DENY;
    add_header X-Content-Type-Options nosniff;
}

# Generate DH params
openssl dhparam -out /etc/ssl/dhparam.pem 4096
```

---

## Go: HTTPS Server with Proper TLS

```go
package main

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "time"
)

func secureHTTPSServer() *http.Server {
    tlsConfig := &tls.Config{
        MinVersion: tls.VersionTLS12,
        MaxVersion: tls.VersionTLS13,
        
        CipherSuites: []uint16{
            // TLS 1.3 (auto, can't be listed)
            // TLS 1.2 only:
            tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
            tls.TLS_ECDHE_RSA_WITH_CHACHA20_POLY1305,
            tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
        },
        
        CurvePreferences: []tls.CurveID{
            tls.X25519,
            tls.CurveP256,
        },
        
        PreferServerCipherSuites: false,  // let client choose in TLS 1.3
    }
    
    mux := http.NewServeMux()
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        // Security headers
        w.Header().Set("Strict-Transport-Security", "max-age=63072000; includeSubDomains; preload")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("Content-Security-Policy", "default-src 'self'")
        
        fmt.Fprintf(w, "Secure HTTPS server — TLS version: %s",
            tls.VersionName(r.TLS.Version))
    })
    
    return &http.Server{
        Addr:         ":443",
        Handler:      mux,
        TLSConfig:    tlsConfig,
        ReadTimeout:  5 * time.Second,
        WriteTimeout: 10 * time.Second,
        IdleTimeout:  120 * time.Second,
    }
}

func main() {
    server := secureHTTPSServer()
    fmt.Println("Starting HTTPS server on :443")
    if err := server.ListenAndServeTLS("server.crt", "server.key"); err != nil {
        panic(err)
    }
}
```

---

## Testing TLS Configuration

```bash
# testssl.sh — comprehensive TLS testing
./testssl.sh https://example.com
# Checks: protocols, ciphers, vulnerabilities (BEAST, POODLE, Heartbleed)

# sslscan
sslscan example.com:443

# nmap ssl scripts
nmap --script ssl-enum-ciphers -p 443 example.com
nmap --script ssl-heartbleed -p 443 example.com

# Online testers
# ssllabs.com — SSL Labs Server Test (gets A+ rating target)
# cryptcheck.fr
```

---

## HTTPS Headers

```
Strict-Transport-Security (HSTS)
   max-age=63072000 — remember this for 2 years
   includeSubDomains — applies to all subdomains
   preload — included in browser preload list

Content-Security-Policy (CSP)
   Prevents XSS by specifying where content can load from

X-Frame-Options
   Prevents clickjacking (iframe embedding)

X-Content-Type-Options: nosniff
   Prevents MIME-sniffing attacks
```

---

## Summary

| TLS Feature | TLS 1.2 | TLS 1.3 |
|-------------|---------|---------|
| Handshake RTT | 2 | 1 |
| Forward secrecy | Optional | Mandatory |
| RC4, 3DES | Allowed | Removed |
| 0-RTT resumption | No | Yes (with caveats) |
| Encrypted handshake | Partial | More encrypted |

---

## Exercises

1. Run `testssl.sh` against a server you control — achieve an A+ rating on SSL Labs
2. Configure nginx with TLS 1.2/1.3 only and strong ciphers — verify with sslscan
3. Build the Go HTTPS server and verify the TLS version and cipher suite in the browser
4. Research HSTS preload — what does submitting your domain to Chrome's preload list mean?
