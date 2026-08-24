# Chapter 12: Go Networking — TCP, UDP, HTTP, Raw Sockets

*Go's `net` package is one of the most powerful in the standard library. From raw TCP connections to full HTTP clients with custom headers, it gives you everything needed to build professional network security tools.*

---

## TCP Connections

```go
package main

import (
    "bufio"
    "fmt"
    "net"
    "time"
)

// Basic TCP client
func tcpConnect(host string, port int, timeout time.Duration) (net.Conn, error) {
    address := fmt.Sprintf("%s:%d", host, port)
    return net.DialTimeout("tcp", address, timeout)
}

// Banner grabbing — read initial server response
func grabBanner(host string, port int) (string, error) {
    conn, err := tcpConnect(host, port, 3*time.Second)
    if err != nil {
        return "", err
    }
    defer conn.Close()
    
    // Short read deadline
    conn.SetReadDeadline(time.Now().Add(2 * time.Second))
    
    scanner := bufio.NewScanner(conn)
    if scanner.Scan() {
        return scanner.Text(), nil
    }
    return "", scanner.Err()
}

// TCP server
func startEchoServer(port int) {
    listener, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
    if err != nil {
        panic(err)
    }
    defer listener.Close()
    fmt.Printf("Echo server on :%d\n", port)
    
    for {
        conn, err := listener.Accept()
        if err != nil {
            return
        }
        go handleConn(conn)
    }
}

func handleConn(conn net.Conn) {
    defer conn.Close()
    remote := conn.RemoteAddr().String()
    fmt.Printf("Connection from %s\n", remote)
    
    scanner := bufio.NewScanner(conn)
    for scanner.Scan() {
        line := scanner.Text()
        conn.Write([]byte("Echo: " + line + "\n"))
    }
}
```

---

## UDP — DNS Query from Scratch

```go
package main

import (
    "encoding/binary"
    "fmt"
    "net"
    "strings"
)

// Build a minimal DNS query packet
func buildDNSQuery(domain string, qtype uint16) []byte {
    var buf []byte
    
    // Header
    buf = append(buf,
        0x12, 0x34,  // transaction ID
        0x01, 0x00,  // flags: standard query, recursion desired
        0x00, 0x01,  // questions: 1
        0x00, 0x00,  // answers: 0
        0x00, 0x00,  // authority: 0
        0x00, 0x00,  // additional: 0
    )
    
    // Question: encode domain as labels
    for _, label := range strings.Split(domain, ".") {
        buf = append(buf, byte(len(label)))
        buf = append(buf, []byte(label)...)
    }
    buf = append(buf, 0x00)  // end of domain
    
    // Query type (A=1, MX=15, AAAA=28, TXT=16)
    qtypeBuf := make([]byte, 2)
    binary.BigEndian.PutUint16(qtypeBuf, qtype)
    buf = append(buf, qtypeBuf...)
    buf = append(buf, 0x00, 0x01)  // class IN
    
    return buf
}

func dnsLookup(domain, server string) {
    conn, err := net.Dial("udp", server+":53")
    if err != nil {
        panic(err)
    }
    defer conn.Close()
    
    query := buildDNSQuery(domain, 1)  // A record
    conn.Write(query)
    
    response := make([]byte, 512)
    n, _ := conn.Read(response)
    
    fmt.Printf("DNS response for %s: %d bytes\n", domain, n)
    fmt.Printf("Raw: %x\n", response[:n])
    // Parsing the response is an exercise — see net.Resolver for the easy way
}

// Easy way: Go's built-in resolver
func easyDNS(domain string) {
    ips, err := net.LookupHost(domain)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    fmt.Printf("%s → %v\n", domain, ips)
    
    mx, _ := net.LookupMX(domain)
    for _, m := range mx {
        fmt.Printf("MX: %s (pref %d)\n", m.Host, m.Pref)
    }
    
    txts, _ := net.LookupTXT(domain)
    for _, txt := range txts {
        fmt.Println("TXT:", txt)
    }
}
```

---

## HTTP Client — Security Tool Essentials

```go
package main

import (
    "crypto/tls"
    "fmt"
    "io"
    "net/http"
    "net/url"
    "strings"
    "time"
)

// Configured HTTP client for security tools
func newSecurityClient(followRedirects bool) *http.Client {
    transport := &http.Transport{
        TLSClientConfig: &tls.Config{
            InsecureSkipVerify: true,  // accept self-signed certs (for testing)
            MinVersion:         tls.VersionTLS10,  // test older TLS too
        },
        MaxIdleConns:    100,
        IdleConnTimeout: 30 * time.Second,
    }
    
    client := &http.Client{
        Transport: transport,
        Timeout:   10 * time.Second,
    }
    
    if !followRedirects {
        client.CheckRedirect = func(req *http.Request, via []*http.Request) error {
            return http.ErrUseLastResponse
        }
    }
    
    return client
}

// Check security headers
func checkHeaders(targetURL string) {
    client := newSecurityClient(false)
    
    resp, err := client.Get(targetURL)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    defer resp.Body.Close()
    
    fmt.Printf("Status: %d %s\n", resp.StatusCode, resp.Status)
    
    secHeaders := map[string]string{
        "Strict-Transport-Security": "HSTS",
        "X-Frame-Options":           "Clickjacking protection",
        "X-Content-Type-Options":    "MIME sniffing protection",
        "Content-Security-Policy":   "CSP",
        "X-XSS-Protection":          "XSS filter",
    }
    
    for header, desc := range secHeaders {
        if val := resp.Header.Get(header); val != "" {
            fmt.Printf("  [PRESENT] %s: %s\n", desc, val)
        } else {
            fmt.Printf("  [MISSING] %s\n", desc)
        }
    }
    
    // Check cookies
    for _, cookie := range resp.Cookies() {
        fmt.Printf("  Cookie: %s", cookie.Name)
        if cookie.HttpOnly { fmt.Print(" [HttpOnly]") }
        if cookie.Secure   { fmt.Print(" [Secure]") }
        if cookie.SameSite != 0 { fmt.Printf(" [SameSite=%v]", cookie.SameSite) }
        fmt.Println()
    }
}

// POST form data (login brute force)
func postForm(targetURL string, params map[string]string) (*http.Response, error) {
    client := newSecurityClient(false)
    
    form := url.Values{}
    for k, v := range params {
        form.Set(k, v)
    }
    
    req, err := http.NewRequest("POST", targetURL, strings.NewReader(form.Encode()))
    if err != nil {
        return nil, err
    }
    req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
    req.Header.Set("User-Agent", "Mozilla/5.0 (compatible; SecurityScanner/1.0)")
    
    return client.Do(req)
}

// POST JSON (API testing)
func postJSON(targetURL, body string, headers map[string]string) (int, string, error) {
    client := newSecurityClient(true)
    
    req, err := http.NewRequest("POST", targetURL, strings.NewReader(body))
    if err != nil {
        return 0, "", err
    }
    req.Header.Set("Content-Type", "application/json")
    for k, v := range headers {
        req.Header.Set(k, v)
    }
    
    resp, err := client.Do(req)
    if err != nil {
        return 0, "", err
    }
    defer resp.Body.Close()
    
    respBody, _ := io.ReadAll(resp.Body)
    return resp.StatusCode, string(respBody), nil
}

func main() {
    // Check a URL's security headers
    checkHeaders("https://example.com")
    
    // Test a login form
    resp, _ := postForm("http://dvwa.local/login.php", map[string]string{
        "username": "admin",
        "password": "password",
        "Login":    "Login",
    })
    if resp != nil {
        fmt.Printf("Login attempt: %d\n", resp.StatusCode)
        // 302 redirect usually = success; 200 = probably failed (stays on login)
    }
}
```

---

## HTTP Server — Building Security Tool APIs

```go
package main

import (
    "encoding/json"
    "fmt"
    "log"
    "net/http"
)

type ScanRequest struct {
    Target string `json:"target"`
    Ports  []int  `json:"ports"`
}

type ScanResponse struct {
    Target string   `json:"target"`
    Open   []int    `json:"open_ports"`
    Error  string   `json:"error,omitempty"`
}

func handleScan(w http.ResponseWriter, r *http.Request) {
    if r.Method != http.MethodPost {
        http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
        return
    }
    
    var req ScanRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        http.Error(w, "Bad request: "+err.Error(), http.StatusBadRequest)
        return
    }
    
    // ... actual scan logic here
    
    w.Header().Set("Content-Type", "application/json")
    json.NewEncoder(w).Encode(ScanResponse{
        Target: req.Target,
        Open:   []int{22, 80, 443},
    })
}

func main() {
    mux := http.NewServeMux()
    mux.HandleFunc("/scan", handleScan)
    mux.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
        w.Write([]byte(`{"status":"ok"}`))
    })
    
    fmt.Println("Security API server on :8080")
    log.Fatal(http.ListenAndServe(":8080", mux))
}
```

---

## Netcat-like Tool in Go

```go
package main

import (
    "flag"
    "fmt"
    "io"
    "net"
    "os"
)

func main() {
    listen := flag.Bool("l", false, "Listen mode")
    port   := flag.String("p", "4444", "Port")
    flag.Parse()
    
    if *listen {
        // Server mode
        ln, err := net.Listen("tcp", ":"+*port)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        fmt.Fprintf(os.Stderr, "Listening on :%s\n", *port)
        conn, err := ln.Accept()
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        defer conn.Close()
        fmt.Fprintf(os.Stderr, "Connection from %s\n", conn.RemoteAddr())
        relay(conn)
    } else {
        // Client mode
        args := flag.Args()
        if len(args) < 1 {
            fmt.Fprintln(os.Stderr, "Usage: gonc [-l -p port] [host]")
            os.Exit(1)
        }
        conn, err := net.Dial("tcp", args[0]+":"+*port)
        if err != nil {
            fmt.Fprintln(os.Stderr, err)
            os.Exit(1)
        }
        defer conn.Close()
        relay(conn)
    }
}

// Bidirectional relay: stdin/stdout ↔ connection
func relay(conn net.Conn) {
    done := make(chan struct{})
    go func() {
        io.Copy(conn, os.Stdin)   // stdin → connection
        conn.(*net.TCPConn).CloseWrite()
        close(done)
    }()
    io.Copy(os.Stdout, conn)    // connection → stdout
    <-done
}
```

---

## Summary

| Task | Go package/approach |
|------|-------------------|
| TCP connect | `net.DialTimeout` |
| TCP listen | `net.Listen` + `Accept` |
| UDP send/recv | `net.Dial("udp", ...)` |
| DNS query | `net.LookupHost`, `net.LookupMX` |
| HTTP client | `net/http.Client` with custom Transport |
| HTTP server | `net/http.ServeMux` |
| TLS inspection | `resp.TLS.PeerCertificates` |
| Raw packets | `golang.org/x/net/ipv4` + raw socket |

---

## Exercises

1. Write a Go function that resolves a hostname and returns whether each IP is in a private range
2. Build an HTTP header checker that takes a list of URLs from stdin and outputs a CSV: URL, has-HSTS, has-CSP, has-X-Frame-Options
3. Implement a TCP reverse shell (for authorized lab use only): client connects back to attacker, attacker types commands, output returned
4. Write a DNS brute-forcer that tries subdomains from a wordlist and prints which ones resolve
