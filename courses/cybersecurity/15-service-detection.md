# Chapter 15: Service Detection and Banner Grabbing

*Knowing a port is open is just the start. Service detection tells you WHAT is running — the application, version, and configuration. Version information drives vulnerability identification.*

---

## Why Service Detection Matters

```
Port 22 open → Could be OpenSSH 9.0 (safe) or OpenSSH 7.2 (CVE-2016-10009)
Port 80 open → Apache 2.4.49 (CVE-2021-41773 — path traversal!) or nginx 1.22 (safe)
Port 21 open → vsftpd 2.3.4 (backdoor!) or ProFTPD 1.3.7 (safe)
```

Without version detection, you're blind to which vulnerabilities apply.

---

## Banner Grabbing

Many services immediately send a banner when you connect — announcing themselves.

```bash
# Manual banner grabbing
nc -v 192.168.1.1 22    # SSH banner
nc -v 192.168.1.1 21    # FTP banner
nc -v 192.168.1.1 25    # SMTP banner

# With timeout
echo "" | nc -v -w 3 192.168.1.1 80
```

### Common Banners

```
SSH:   SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.6
FTP:   220 vsftpd 3.0.5
SMTP:  220 mail.example.com ESMTP Postfix
HTTP:  HTTP/1.1 200 OK
       Server: Apache/2.4.54 (Debian)
MySQL: J\x00\x00\x00\n8.0.33
POP3:  +OK Dovecot ready.
IMAP:  * OK Dovecot ready.
```

---

## HTTP Service Detection

HTTP gives the most info — server header, powered-by header, error pages.

```go
package main

import (
    "crypto/tls"
    "fmt"
    "net/http"
    "strings"
    "time"
)

type WebInfo struct {
    URL         string
    Status      int
    Server      string
    PoweredBy   string
    Title       string
    Headers     map[string]string
    TLSVersion  string
    Certificate string
}

func detectWeb(targetURL string) (*WebInfo, error) {
    client := &http.Client{
        Timeout: 10 * time.Second,
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
        CheckRedirect: func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        },
    }
    
    resp, err := client.Get(targetURL)
    if err != nil {
        return nil, err
    }
    defer resp.Body.Close()
    
    info := &WebInfo{
        URL:       targetURL,
        Status:    resp.StatusCode,
        Server:    resp.Header.Get("Server"),
        PoweredBy: resp.Header.Get("X-Powered-By"),
        Headers:   make(map[string]string),
    }
    
    for k, v := range resp.Header {
        info.Headers[k] = strings.Join(v, ", ")
    }
    
    if resp.TLS != nil {
        versions := map[uint16]string{
            tls.VersionTLS10: "TLS 1.0 (deprecated!)",
            tls.VersionTLS11: "TLS 1.1 (deprecated!)",
            tls.VersionTLS12: "TLS 1.2",
            tls.VersionTLS13: "TLS 1.3",
        }
        info.TLSVersion = versions[resp.TLS.Version]
        if len(resp.TLS.PeerCertificates) > 0 {
            cert := resp.TLS.PeerCertificates[0]
            info.Certificate = fmt.Sprintf("%s (expires %s)",
                cert.Subject.CommonName, cert.NotAfter.Format("2006-01-02"))
        }
    }
    
    return info, nil
}

func printWebInfo(info *WebInfo) {
    fmt.Printf("URL: %s\n", info.URL)
    fmt.Printf("Status: %d\n", info.Status)
    if info.Server != "" {
        fmt.Printf("Server: %s\n", info.Server)
    }
    if info.PoweredBy != "" {
        fmt.Printf("Powered-By: %s\n", info.PoweredBy)
    }
    if info.TLSVersion != "" {
        fmt.Printf("TLS: %s\n", info.TLSVersion)
        fmt.Printf("Cert: %s\n", info.Certificate)
    }
    
    // Security headers check
    secHeaders := []string{
        "Strict-Transport-Security",
        "Content-Security-Policy",
        "X-Frame-Options",
        "X-Content-Type-Options",
    }
    for _, h := range secHeaders {
        if _, ok := info.Headers[h]; !ok {
            fmt.Printf("[MISSING] %s\n", h)
        }
    }
}

func main() {
    info, err := detectWeb("https://example.com")
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    printWebInfo(info)
}
```

---

## Nmap Service Detection

Nmap has 10,000+ service fingerprints. Use it for accurate version detection.

```bash
# Version detection (-sV)
nmap -sV 192.168.1.1

# Intensive version scan
nmap -sV --version-intensity 9 192.168.1.1

# With NSE scripts
nmap -sV -sC 192.168.1.1   # default scripts

# Specific service scripts
nmap --script ssh-hostkey 192.168.1.1
nmap --script http-headers 192.168.1.1
nmap --script ftp-anon 192.168.1.1         # anonymous FTP login?
nmap --script smb-vuln-ms17-010 192.168.1.1  # EternalBlue!
nmap --script mysql-empty-password 192.168.1.1
nmap --script ssl-cert,ssl-enum-ciphers 192.168.1.1

# Full recon scan
nmap -sV -sC -O -p- --min-rate 5000 192.168.1.1 -oA scan_results
```

### NSE Script Categories

```bash
# Categories: auth, broadcast, default, discovery, exploit, fuzzer, 
#             intrusive, malware, safe, version, vuln

nmap --script vuln 192.168.1.1    # run all vuln detection scripts
nmap --script malware 192.168.1.1  # check for malware indicators
nmap --script auth 192.168.1.1     # check for auth bypass

# List all scripts
ls /usr/share/nmap/scripts/
```

---

## Go: Full Service Fingerprinter

```go
package main

import (
    "bufio"
    "crypto/tls"
    "fmt"
    "net"
    "strings"
    "time"
)

type ServiceResult struct {
    Port    int
    Proto   string
    Service string
    Version string
    Banner  string
    Info    map[string]string
}

// Service detection functions per port
var detectors = map[int]func(string, int) ServiceResult{
    21:   detectFTP,
    22:   detectSSH,
    25:   detectSMTP,
    80:   detectHTTP,
    443:  detectHTTPS,
    3306: detectMySQL,
}

func grabBanner(host string, port int, send string, timeout time.Duration) string {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", addr, timeout)
    if err != nil {
        return ""
    }
    defer conn.Close()
    conn.SetDeadline(time.Now().Add(timeout))
    
    if send != "" {
        conn.Write([]byte(send))
    }
    
    scanner := bufio.NewScanner(conn)
    if scanner.Scan() {
        return scanner.Text()
    }
    return ""
}

func detectSSH(host string, port int) ServiceResult {
    banner := grabBanner(host, port, "", 3*time.Second)
    result := ServiceResult{Port: port, Proto: "tcp", Service: "SSH"}
    
    if strings.HasPrefix(banner, "SSH-") {
        result.Banner = banner
        // SSH-2.0-OpenSSH_8.9p1 Ubuntu-3ubuntu0.6
        parts := strings.SplitN(banner, "-", 3)
        if len(parts) >= 3 {
            result.Version = parts[2]
        }
    }
    return result
}

func detectFTP(host string, port int) ServiceResult {
    banner := grabBanner(host, port, "", 3*time.Second)
    result := ServiceResult{Port: port, Proto: "tcp", Service: "FTP"}
    result.Banner = banner
    
    // Detect vsftpd 2.3.4 backdoor
    if strings.Contains(banner, "vsftpd 2.3.4") {
        result.Info = map[string]string{
            "VULNERABLE": "vsftpd 2.3.4 contains a backdoor (CVE-2011-2523)",
        }
    }
    return result
}

func detectSMTP(host string, port int) ServiceResult {
    banner := grabBanner(host, port, "", 3*time.Second)
    return ServiceResult{
        Port:    port,
        Proto:   "tcp",
        Service: "SMTP",
        Banner:  banner,
    }
}

func detectHTTP(host string, port int) ServiceResult {
    req := "HEAD / HTTP/1.0\r\nHost: " + host + "\r\n\r\n"
    banner := grabBanner(host, port, req, 5*time.Second)
    result := ServiceResult{Port: port, Proto: "tcp", Service: "HTTP"}
    
    // Parse Server header
    if strings.Contains(banner, "Server:") {
        result.Version = strings.TrimPrefix(banner, "Server: ")
    }
    result.Banner = banner
    return result
}

func detectHTTPS(host string, port int) ServiceResult {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := tls.DialWithDialer(
        &net.Dialer{Timeout: 5 * time.Second},
        "tcp", addr,
        &tls.Config{InsecureSkipVerify: true},
    )
    result := ServiceResult{Port: port, Proto: "tcp", Service: "HTTPS"}
    
    if err != nil {
        return result
    }
    defer conn.Close()
    
    if certs := conn.ConnectionState().PeerCertificates; len(certs) > 0 {
        result.Version = fmt.Sprintf("cert: %s (expires %s)",
            certs[0].Subject.CommonName,
            certs[0].NotAfter.Format("2006-01-02"))
    }
    return result
}

func detectMySQL(host string, port int) ServiceResult {
    addr := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", addr, 3*time.Second)
    result := ServiceResult{Port: port, Proto: "tcp", Service: "MySQL"}
    
    if err != nil {
        return result
    }
    defer conn.Close()
    
    // MySQL sends initial handshake on connect
    buf := make([]byte, 256)
    n, err := conn.Read(buf)
    if err == nil && n > 5 {
        // Version string starts at offset 5
        end := n
        for i := 5; i < n; i++ {
            if buf[i] == 0 {
                end = i
                break
            }
        }
        if end > 5 {
            result.Version = string(buf[5:end])
        }
    }
    return result
}

func detectService(host string, port int) ServiceResult {
    if detector, ok := detectors[port]; ok {
        return detector(host, port)
    }
    // Generic banner grab for unknown ports
    banner := grabBanner(host, port, "\r\n", 2*time.Second)
    return ServiceResult{Port: port, Proto: "tcp", Banner: banner}
}

func main() {
    host := "192.168.1.1"
    openPorts := []int{21, 22, 25, 80, 443, 3306}
    
    fmt.Printf("Service detection on %s\n", host)
    fmt.Println(strings.Repeat("─", 60))
    
    for _, port := range openPorts {
        r := detectService(host, port)
        fmt.Printf("%5d/%-4s %-10s %s\n",
            r.Port, r.Proto, r.Service, r.Version)
        if r.Banner != "" && r.Version == "" {
            fmt.Printf("            Banner: %s\n", r.Banner[:min(len(r.Banner), 80)])
        }
        for k, v := range r.Info {
            fmt.Printf("            [%s] %s\n", k, v)
        }
    }
}

func min(a, b int) int { if a < b { return a }; return b }
```

---

## Summary

| Method | Use case | Accuracy |
|--------|---------|----------|
| Banner grabbing | Quick, passive | Good for verbose services |
| HTTP headers | Web servers | Excellent (Server header) |
| TLS cert | HTTPS services | Good (cert CN = domain) |
| MySQL handshake | Databases | Exact version |
| Nmap -sV | All ports | Best available |

---

## Exercises

1. Extend the service fingerprinter to detect PostgreSQL (port 5432) and Redis (port 6379)
2. Write a tool that scans a /24 network and identifies all web servers, outputting their versions
3. Add CVE lookup: for detected service versions, check a local JSON file of known CVEs
4. Build an SMTP enumerator that tries VRFY and EXPN commands to enumerate users
