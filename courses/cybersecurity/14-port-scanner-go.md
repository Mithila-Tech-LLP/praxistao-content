# Chapter 14: Your First Security Tool — Port Scanner in Go

*A port scanner is one of the most fundamental security tools. Nmap is the industry standard. In this chapter, you'll build a complete, production-quality port scanner from scratch in Go. By the end, you'll understand exactly how Nmap works — because you'll have built it.*

---

## What Port Scanning Does

Port scanning reveals:
- **Which ports are open** (services are listening)
- **What service** might be running on each port
- **What version** of that service (banner grabbing)
- **What OS** is running (TTL analysis, TCP stack fingerprinting)

This information is essential for both attackers (finding attack surface) and defenders (verifying firewall rules, discovering unauthorized services).

**Legal note:** Only scan systems you own or have explicit written permission to scan. Unauthorized scanning is illegal in most jurisdictions.

`scanme.nmap.org` is maintained by the Nmap team specifically for testing — you can scan it legally.

---

## Building the Scanner — Step by Step

### Step 1: Basic Port Check

```go
// file: cmd/scanner/main.go
package main

import (
    "fmt"
    "net"
    "os"
    "strconv"
    "time"
)

func checkPort(host string, port int, timeout time.Duration) bool {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return false
    }
    conn.Close()
    return true
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: scanner <host> [start-port] [end-port]")
        os.Exit(1)
    }
    
    host := os.Args[1]
    startPort, endPort := 1, 1024
    
    if len(os.Args) >= 4 {
        startPort, _ = strconv.Atoi(os.Args[2])
        endPort, _ = strconv.Atoi(os.Args[3])
    }
    
    timeout := 1 * time.Second
    fmt.Printf("Scanning %s ports %d-%d...\n\n", host, startPort, endPort)
    
    for port := startPort; port <= endPort; port++ {
        if checkPort(host, port, timeout) {
            fmt.Printf("Port %d: OPEN\n", port)
        }
    }
}
```

Test it:
```bash
go run cmd/scanner/main.go scanme.nmap.org 1 100
```

This works — but it's sequential. Scanning 65,535 ports at 1 second timeout = 18 hours. We need concurrency.

---

### Step 2: Concurrent Scanning with Worker Pool

The naive approach (one goroutine per port) can overwhelm the target with 65,535 simultaneous connections. Instead, we use a **worker pool** pattern:

```go
// file: cmd/scanner/main.go
package main

import (
    "fmt"
    "net"
    "os"
    "sort"
    "strconv"
    "sync"
    "time"
)

// ScanResult holds the result for one port
type ScanResult struct {
    Port   int
    Open   bool
    Banner string
}

// worker processes ports from a job channel
func worker(host string, timeout time.Duration, jobs <-chan int, results chan<- ScanResult, wg *sync.WaitGroup) {
    defer wg.Done()
    
    for port := range jobs {
        address := fmt.Sprintf("%s:%d", host, port)
        conn, err := net.DialTimeout("tcp", address, timeout)
        
        if err != nil {
            results <- ScanResult{Port: port, Open: false}
            continue
        }
        
        // Try to grab a banner (service response)
        banner := ""
        conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
        buf := make([]byte, 256)
        n, err := conn.Read(buf)
        if err == nil && n > 0 {
            banner = string(buf[:n])
        }
        conn.Close()
        
        results <- ScanResult{Port: port, Open: true, Banner: banner}
    }
}

func scan(host string, startPort, endPort, workerCount int, timeout time.Duration) []ScanResult {
    jobs := make(chan int, workerCount)
    results := make(chan ScanResult, endPort-startPort+1)
    
    var wg sync.WaitGroup
    
    // Start workers
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go worker(host, timeout, jobs, results, &wg)
    }
    
    // Send jobs
    go func() {
        for port := startPort; port <= endPort; port++ {
            jobs <- port
        }
        close(jobs)
    }()
    
    // Wait for all workers to finish, then close results
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    var openPorts []ScanResult
    for result := range results {
        if result.Open {
            openPorts = append(openPorts, result)
        }
    }
    
    // Sort by port number
    sort.Slice(openPorts, func(i, j int) bool {
        return openPorts[i].Port < openPorts[j].Port
    })
    
    return openPorts
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: scanner <host> [start] [end] [workers]")
        os.Exit(1)
    }
    
    host := os.Args[1]
    startPort, endPort, workerCount := 1, 1024, 500
    timeout := 1 * time.Second
    
    if len(os.Args) >= 4 {
        startPort, _ = strconv.Atoi(os.Args[2])
        endPort, _ = strconv.Atoi(os.Args[3])
    }
    if len(os.Args) >= 5 {
        workerCount, _ = strconv.Atoi(os.Args[4])
    }
    
    start := time.Now()
    fmt.Printf("Scanning %s ports %d-%d with %d workers...\n\n", host, startPort, endPort, workerCount)
    
    results := scan(host, startPort, endPort, workerCount, timeout)
    
    elapsed := time.Since(start)
    
    if len(results) == 0 {
        fmt.Println("No open ports found.")
    } else {
        fmt.Printf("%-10s %-20s %s\n", "PORT", "SERVICE", "BANNER")
        fmt.Println("─────────────────────────────────────────────────")
        for _, r := range results {
            service := getServiceName(r.Port)
            banner := r.Banner
            if len(banner) > 50 {
                banner = banner[:50] + "..."
            }
            // Remove newlines from banner
            for i := 0; i < len(banner); i++ {
                if banner[i] == '\n' || banner[i] == '\r' {
                    banner = banner[:i]
                    break
                }
            }
            fmt.Printf("%-10d %-20s %s\n", r.Port, service, banner)
        }
    }
    
    fmt.Printf("\nScanned %d ports in %s\n", endPort-startPort+1, elapsed.Round(time.Millisecond))
}

// getServiceName returns the common service name for well-known ports
func getServiceName(port int) string {
    services := map[int]string{
        21:   "FTP",
        22:   "SSH",
        23:   "Telnet",
        25:   "SMTP",
        53:   "DNS",
        80:   "HTTP",
        110:  "POP3",
        143:  "IMAP",
        443:  "HTTPS",
        445:  "SMB",
        993:  "IMAPS",
        995:  "POP3S",
        1433: "MSSQL",
        3306: "MySQL",
        3389: "RDP",
        5432: "PostgreSQL",
        5900: "VNC",
        6379: "Redis",
        8080: "HTTP-Alt",
        8443: "HTTPS-Alt",
        9200: "Elasticsearch",
        27017: "MongoDB",
    }
    if name, ok := services[port]; ok {
        return name
    }
    return "unknown"
}
```

Run:
```bash
go run cmd/scanner/main.go scanme.nmap.org 1 10000 1000
```

Sample output:
```
Scanning scanme.nmap.org ports 1-10000 with 1000 workers...

PORT       SERVICE              BANNER
─────────────────────────────────────────────────────
22         SSH                  SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
80         HTTP                 HTTP/1.0 302 Found

Scanned 10000 ports in 4.231s
```

---

### Step 3: UDP Scanning

TCP scanning is straightforward — you get a response. UDP is harder because UDP doesn't have a connection. We send a packet; if we get an ICMP "port unreachable" back, the port is closed. No response might mean open or filtered.

```go
func scanUDP(host string, port int, timeout time.Duration) string {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("udp", address, timeout)
    if err != nil {
        return "closed"
    }
    defer conn.Close()
    
    // Send an empty probe
    conn.SetDeadline(time.Now().Add(timeout))
    _, err = conn.Write([]byte{})
    if err != nil {
        return "closed"
    }
    
    // Try to read (ICMP unreachable would cause an error here on closed UDP ports)
    buf := make([]byte, 1024)
    _, err = conn.Read(buf)
    if err != nil {
        // Error might mean closed (ICMP) or filtered
        return "open|filtered"
    }
    
    return "open"
}
```

Note: Full UDP scanning with ICMP parsing requires raw sockets, which need root privileges. We'll cover this in Chapter 19.

---

### Step 4: OS Detection via TTL

```go
import "golang.org/x/net/icmp"

// Rough OS detection based on TTL value
func guessOS(ttl int) string {
    switch {
    case ttl <= 64:
        return "Linux/Unix (TTL <= 64)"
    case ttl <= 128:
        return "Windows (TTL <= 128)"
    case ttl <= 255:
        return "Network device / old Unix (TTL <= 255)"
    default:
        return "Unknown"
    }
}
```

---

### Step 5: Full Feature Scanner

Let's add command-line flags for a production-quality experience:

```go
package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "sort"
    "strings"
    "sync"
    "time"
)

type Config struct {
    Host        string
    StartPort   int
    EndPort     int
    Workers     int
    Timeout     time.Duration
    GrabBanners bool
    TopPorts    bool
}

// Top 20 most common ports to scan quickly
var topPorts = []int{21, 22, 23, 25, 53, 80, 110, 135, 139, 143, 443, 445, 993, 995, 1723, 3306, 3389, 5900, 8080, 8443}

func main() {
    config := &Config{}
    
    flag.StringVar(&config.Host, "host", "", "Target host (required)")
    flag.IntVar(&config.StartPort, "start", 1, "Start port")
    flag.IntVar(&config.EndPort, "end", 1024, "End port")
    flag.IntVar(&config.Workers, "workers", 500, "Number of concurrent workers")
    flag.DurationVar(&config.Timeout, "timeout", 1*time.Second, "Connection timeout")
    flag.BoolVar(&config.GrabBanners, "banner", true, "Attempt banner grabbing")
    flag.BoolVar(&config.TopPorts, "top", false, "Scan top 20 common ports only")
    
    flag.Parse()
    
    if config.Host == "" {
        // Try positional arg
        if flag.NArg() > 0 {
            config.Host = flag.Arg(0)
        } else {
            fmt.Fprintln(os.Stderr, "Error: host is required")
            fmt.Println("Usage: scanner -host <target> [options]")
            flag.PrintDefaults()
            os.Exit(1)
        }
    }
    
    // Resolve hostname
    ips, err := net.LookupHost(config.Host)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error resolving %s: %v\n", config.Host, err)
        os.Exit(1)
    }
    
    fmt.Printf("GoScanner — Port Scanner\n")
    fmt.Printf("Target: %s (%s)\n", config.Host, strings.Join(ips, ", "))
    
    var results []ScanResult
    start := time.Now()
    
    if config.TopPorts {
        fmt.Printf("Scanning top %d ports...\n\n", len(topPorts))
        results = scanPorts(config.Host, topPorts, config.Workers, config.Timeout, config.GrabBanners)
    } else {
        portsToScan := endPort - startPort + 1
        fmt.Printf("Scanning %d ports (%d-%d) with %d workers...\n\n", 
                   portsToScan, config.StartPort, config.EndPort, config.Workers)
        
        portList := make([]int, 0, portsToScan)
        for p := config.StartPort; p <= config.EndPort; p++ {
            portList = append(portList, p)
        }
        results = scanPorts(config.Host, portList, config.Workers, config.Timeout, config.GrabBanners)
    }
    
    elapsed := time.Since(start)
    printResults(results, elapsed)
}

func scanPorts(host string, ports []int, workerCount int, timeout time.Duration, grabBanners bool) []ScanResult {
    jobs := make(chan int, workerCount)
    results := make(chan ScanResult, len(ports))
    var wg sync.WaitGroup
    
    for i := 0; i < workerCount; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for port := range jobs {
                result := checkPortWithBanner(host, port, timeout, grabBanners)
                results <- result
            }
        }()
    }
    
    go func() {
        for _, port := range ports {
            jobs <- port
        }
        close(jobs)
    }()
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var open []ScanResult
    for r := range results {
        if r.Open {
            open = append(open, r)
        }
    }
    
    sort.Slice(open, func(i, j int) bool {
        return open[i].Port < open[j].Port
    })
    
    return open
}

func checkPortWithBanner(host string, port int, timeout time.Duration, grabBanner bool) ScanResult {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return ScanResult{Port: port, Open: false}
    }
    
    result := ScanResult{Port: port, Open: true}
    
    if grabBanner {
        conn.SetReadDeadline(time.Now().Add(500 * time.Millisecond))
        buf := make([]byte, 256)
        n, err := conn.Read(buf)
        if err == nil && n > 0 {
            banner := strings.TrimSpace(string(buf[:n]))
            // Clean up banner — take only first line, remove special chars
            if idx := strings.IndexAny(banner, "\r\n"); idx > 0 {
                banner = banner[:idx]
            }
            result.Banner = banner
        }
    }
    
    conn.Close()
    return result
}

func printResults(results []ScanResult, elapsed time.Duration) {
    if len(results) == 0 {
        fmt.Println("No open ports found.")
    } else {
        fmt.Printf("Open ports: %d\n\n", len(results))
        fmt.Printf("%-8s %-20s %s\n", "PORT", "SERVICE", "BANNER")
        fmt.Println(strings.Repeat("─", 70))
        for _, r := range results {
            service := getServiceName(r.Port)
            fmt.Printf("%-8d %-20s %s\n", r.Port, service, r.Banner)
        }
    }
    fmt.Printf("\nCompleted in %s\n", elapsed.Round(time.Millisecond))
}
```

---

## Sample Output

```
GoScanner — Port Scanner
Target: scanme.nmap.org (45.33.32.156)
Scanning 1000 ports (1-1000) with 500 workers...

Open ports: 2

PORT     SERVICE              BANNER
──────────────────────────────────────────────────────────────────────
22       SSH                  SSH-2.0-OpenSSH_6.6.1p1 Ubuntu-2ubuntu2.13
80       HTTP                 

Completed in 3.456s
```

---

## What We've Built vs Nmap

| Feature | Our Scanner | Nmap |
|---------|-------------|------|
| TCP connect scan | ✅ | ✅ |
| Service name detection | ✅ | ✅ |
| Banner grabbing | ✅ | ✅ |
| Concurrent scanning | ✅ | ✅ |
| UDP scanning | Partial | ✅ |
| SYN scan (raw packets) | ❌ | ✅ (requires root) |
| OS fingerprinting | Partial (TTL) | ✅ |
| Script engine | ❌ | ✅ |

We'll add raw packet scanning (SYN scan) in Chapter 19. The core logic is identical to Nmap.

---

## Build and Test

```bash
# Build
go build -o goscanner ./cmd/scanner

# Scan top ports
./goscanner -host scanme.nmap.org -top

# Scan specific range
./goscanner -host scanme.nmap.org -start 1 -end 65535 -workers 1000

# Disable banner grabbing (faster)
./goscanner -host scanme.nmap.org -banner=false -end 1024

# Cross-compile for Windows
GOOS=windows go build -o goscanner.exe ./cmd/scanner
```

---

## Exercises

1. Add a `-timeout` flag and observe how different timeouts affect speed vs accuracy
2. Add IPv6 support (use `net.DialTimeout("tcp6", ...)`)
3. Add output to a file (JSON or CSV) with `-output results.json`
4. Add CIDR range input (`-host 192.168.1.0/24`) to scan an entire subnet
5. Compare your scanner's speed to Nmap on the same target. What's the difference?
