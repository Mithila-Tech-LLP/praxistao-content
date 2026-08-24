# Chapter 08: Why Go for Security — Speed, Concurrency, and Power

*Go is now the dominant language for security tooling. Nmap, Metasploit, and legacy tools are being replaced by Go-based alternatives. This chapter explains why, and gets you set up.*

---

## Why Go for Security Tools?

### 1. Single Binary Deployment

Go compiles to a single, statically linked binary with no dependencies.

```bash
GOOS=linux GOARCH=amd64 go build -o scanner-linux ./cmd/scanner
GOOS=windows GOARCH=amd64 go build -o scanner.exe ./cmd/scanner
GOOS=darwin GOARCH=arm64 go build -o scanner-mac ./cmd/scanner
```

You can build a tool on your laptop and drop it on any target — no need to install Go, no libraries needed, no runtime.

**Security implication:** Malware written in Go is increasingly common because it compiles to one file, cross-compiles easily, and has no dependencies to detect.

### 2. Built-in Concurrency — Goroutines

Port scanning, network scanning, password cracking — all involve doing many things simultaneously. Go makes this natural:

```go
// Scan 1,000 ports using 1,000 concurrent goroutines
for port := 1; port <= 1000; port++ {
    go scanPort(target, port)  // launches 1,000 goroutines
}
```

Go's goroutines are lightweight (a few KB of stack, vs MB for threads). You can run tens of thousands simultaneously. This makes Go-based scanners dramatically faster than Python equivalents.

### 3. Excellent Networking Libraries

Go's standard library includes:
- `net` — TCP/UDP/Unix socket connections
- `net/http` — Full HTTP client and server
- `crypto/tls` — TLS/SSL
- `encoding/binary` — Binary protocols
- `syscall` — Direct OS system calls

For raw packets, the `google/gopacket` library provides Wireshark-level packet manipulation.

### 4. Performance

Go is compiled (not interpreted like Python) and close to C performance. For:
- Password cracking: billions of hashes per second
- Port scanning: scan 65,535 ports in seconds
- Log analysis: process GB of logs in seconds

### 5. The Go Security Ecosystem

Major security tools written in Go:
- **Nuclei** (ProjectDiscovery) — vulnerability scanner
- **Amass** — network mapping / OSINT
- **gobuster** — directory/DNS brute-forcer
- **subfinder** — subdomain discovery
- **httpx** — fast HTTP probing
- **ffuf** — web fuzzer
- **Trivy** — container vulnerability scanner
- **GoSec** — Go source code security scanner
- **Falco** — runtime security (cloud-native)
- **CrowdStrike Falcon agent** (partially Go)
- **Kubernetes** — container orchestration (Go)

---

## Setting Up Go

### Installation

```bash
# Download from golang.org
wget https://go.dev/dl/go1.21.5.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.5.linux-amd64.tar.gz

# Add to PATH
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc

# Verify
go version  # go version go1.21.5 linux/amd64
```

### Project Structure

```bash
# Create a new project
mkdir security-tools
cd security-tools
go mod init github.com/yourusername/security-tools

# Project structure
security-tools/
├── go.mod              # Module definition and dependencies
├── go.sum              # Dependency checksums
├── cmd/
│   ├── scanner/
│   │   └── main.go     # Port scanner tool
│   └── sniffer/
│       └── main.go     # Network sniffer
├── pkg/
│   ├── network/        # Reusable networking code
│   └── crypto/         # Cryptography utilities
└── internal/           # Private packages (not importable externally)
```

---

## Go Fundamentals for Security

This is a quick but complete Go reference. If you're already a programmer, skim this. If not, study it carefully.

### Variables and Types

```go
package main

import "fmt"

func main() {
    // Explicit type declaration
    var port int = 80
    var hostname string = "192.168.1.1"
    var isOpen bool = true
    
    // Short declaration (type inferred)
    target := "192.168.1.1"
    maxThreads := 100
    
    // Multiple assignment
    ip, port := "10.0.0.1", 443
    
    // Constants
    const MaxPorts = 65535
    
    // Byte slice (raw binary data — essential for security tools)
    rawPacket := []byte{0x41, 0x42, 0x43}   // "ABC"
    
    fmt.Printf("Target: %s, Port: %d, Open: %v\n", hostname, port, isOpen)
    
    // Suppress unused variable warning
    _ = target
    _ = maxThreads
    _ = ip
    _ = rawPacket
}
```

### Security-Relevant Types

```go
// Net package types
import "net"

var ip net.IP = net.ParseIP("192.168.1.1")
var cidr *net.IPNet
ip, cidr, _ = net.ParseCIDR("192.168.1.0/24")

// Check if IP is in network
fmt.Println(cidr.Contains(net.ParseIP("192.168.1.100")))  // true

// Convert IP to 4-byte representation
ipv4 := ip.To4()  // [192, 168, 1, 1]
```

### Functions

```go
package main

import (
    "fmt"
    "net"
    "time"
)

// Returns (bool, error) — idiomatic Go error handling
func isPortOpen(host string, port int, timeout time.Duration) (bool, error) {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return false, nil  // Port is closed/filtered — not an error for us
    }
    conn.Close()
    return true, nil
}

func main() {
    open, err := isPortOpen("google.com", 443, 2*time.Second)
    if err != nil {
        fmt.Println("Error:", err)
        return
    }
    if open {
        fmt.Println("Port 443 is OPEN")
    }
}
```

### Slices and Maps

```go
// Slices — dynamic arrays
ports := []int{22, 80, 443, 3389, 8080}
ports = append(ports, 8443)  // Add element

// Iterate
for i, port := range ports {
    fmt.Printf("Port[%d] = %d\n", i, port)
}

// Maps — key-value store
// serviceNames maps port → service name
serviceNames := map[int]string{
    21:   "FTP",
    22:   "SSH",
    25:   "SMTP",
    53:   "DNS",
    80:   "HTTP",
    443:  "HTTPS",
    3306: "MySQL",
    3389: "RDP",
}

// Lookup with existence check
if name, exists := serviceNames[22]; exists {
    fmt.Printf("Port 22 = %s\n", name)
}

// Map of results
openPorts := make(map[int]bool)
openPorts[80] = true
openPorts[443] = true
```

### Goroutines and Channels

This is where Go shines for security tools:

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

func scanPort(host string, port int, results chan<- int, wg *sync.WaitGroup) {
    defer wg.Done()
    
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, 500*time.Millisecond)
    if err == nil {
        conn.Close()
        results <- port  // Send open port to channel
    }
}

func main() {
    target := "scanme.nmap.org"  // A legal test target
    results := make(chan int, 100)
    var wg sync.WaitGroup
    
    // Launch 1000 goroutines simultaneously
    for port := 1; port <= 1000; port++ {
        wg.Add(1)
        go scanPort(target, port, results, &wg)
    }
    
    // Close channel when all goroutines finish
    go func() {
        wg.Wait()
        close(results)
    }()
    
    // Collect results
    for port := range results {
        fmt.Printf("Port %d: OPEN\n", port)
    }
}
```

This scans 1,000 ports concurrently. Without goroutines (single-threaded), at 500ms timeout this would take 500 seconds. With goroutines, it takes ~0.5 seconds.

### Error Handling

Go doesn't have exceptions. Errors are values returned from functions:

```go
conn, err := net.Dial("tcp", "192.168.1.1:80")
if err != nil {
    // Handle error
    fmt.Println("Connection failed:", err)
    return
}
defer conn.Close()  // Ensure connection is closed when function returns

// Read data
buf := make([]byte, 1024)
n, err := conn.Read(buf)
if err != nil {
    fmt.Println("Read failed:", err)
    return
}
fmt.Printf("Received %d bytes: %s\n", n, buf[:n])
```

### Working with Raw Bytes

Security tools constantly work with raw binary data:

```go
package main

import (
    "encoding/binary"
    "encoding/hex"
    "fmt"
)

func main() {
    // Create a raw packet (simplified TCP header)
    packet := make([]byte, 20)
    
    // Write source port (2 bytes, big-endian) at offset 0
    binary.BigEndian.PutUint16(packet[0:2], 12345)  // source port
    
    // Write dest port (2 bytes, big-endian) at offset 2
    binary.BigEndian.PutUint16(packet[2:4], 80)     // HTTP
    
    // Write sequence number (4 bytes) at offset 4
    binary.BigEndian.PutUint32(packet[4:8], 0xDEADBEEF)
    
    fmt.Println("Raw packet (hex):", hex.EncodeToString(packet))
    
    // Parse a hex string back to bytes
    data, _ := hex.DecodeString("48656c6c6f")
    fmt.Println("Decoded:", string(data))  // "Hello"
    
    // Read from raw bytes
    srcPort := binary.BigEndian.Uint16(packet[0:2])
    dstPort := binary.BigEndian.Uint16(packet[2:4])
    fmt.Printf("Src: %d, Dst: %d\n", srcPort, dstPort)
}
```

---

## Installing Security Libraries

```bash
# Add gopacket for raw packet manipulation (like Wireshark)
go get github.com/google/gopacket

# Add cobra for CLI tools
go get github.com/spf13/cobra

# Add color output
go get github.com/fatih/color
```

---

## Your Go Security Toolkit — Directory Layout

Create this structure for the course:

```bash
mkdir -p ~/security-tools/{cmd/{scanner,sniffer,fuzzer,cracker},pkg/{network,crypto,dns},goshield/{agent,server,detector}}
cd ~/security-tools
go mod init github.com/yourusername/security-tools
```

All the tools we build will live in `cmd/`. The `goshield/` directory is for the enterprise EDR project we'll build in Part 9.

---

## Summary

| Go Feature | Security Tool Application |
|------------|--------------------------|
| Goroutines | Concurrent port/web scanning |
| Channels | Worker pool pattern for scanning |
| `net` package | TCP/UDP connections, DNS |
| `[]byte` | Raw packet manipulation |
| Static binary | Deploy tools without runtime dependencies |
| Cross-compilation | Build Linux tool from Windows/Mac |

---

## Exercises

1. Install Go on your system and run `go version`
2. Create the directory structure above and run `go mod init`
3. Write a Go program that takes an IP and port from command-line args and reports if the port is open
4. Modify it to scan ports 1-1024 concurrently and print all open ports
5. Time the difference: scan ports with 1 goroutine vs 1000 goroutines. What's the speedup?
