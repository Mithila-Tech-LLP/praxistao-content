# Chapter 09: Go Foundations — Variables, Functions, and Control Flow for Security Tools

*You don't need to be a Go expert to write security tools. But you do need to understand the language well enough to read, modify, and build on existing code. This chapter gives you that foundation.*

---

## Why Learn Go for Security?

- **Single binary**: Compile once, run everywhere — no dependencies on the target
- **Goroutines**: Concurrent scanning, parallel attacks without thread management hell
- **Standard library**: `net`, `crypto`, `os`, `syscall` — everything you need
- **Cross-compilation**: Build a Linux agent from your Mac
- **Speed**: Close to C, without C's memory dangers

---

## Installation and Setup

```bash
# Install Go
# Download from https://go.dev/dl/
# Or on Linux:
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Verify
go version    # go version go1.22.0 linux/amd64

# Create a project
mkdir my-security-tool
cd my-security-tool
go mod init github.com/yourname/my-security-tool
```

---

## Variables and Types

```go
package main

import "fmt"

func main() {
    // var keyword (explicit type)
    var port int = 443
    var hostname string = "example.com"
    var isOpen bool = false

    // Short declaration (inferred type)
    timeout := 5         // int
    protocol := "TCP"   // string
    alive := true       // bool

    // Constants
    const MaxPorts = 65535
    const Banner = "GoScanner v1.0"

    fmt.Println(port, hostname, isOpen)
    fmt.Println(timeout, protocol, alive)
    fmt.Println(MaxPorts, Banner)

    // Zero values (default when not initialized)
    var count int      // 0
    var name string    // ""
    var flag bool      // false
    fmt.Println(count, name, flag)
}
```

### Types Important for Security Tools

```go
// Integer types
var pid int = 1234           // platform size (32 or 64 bit)
var port uint16 = 443        // 0-65535, perfect for ports
var uid int32 = 1000         // user IDs

// Byte manipulation (critical for protocol parsing)
var b byte = 0xFF            // byte = uint8
var data []byte = []byte{0x47, 0x45, 0x54} // "GET"

// Strings are immutable byte sequences
s := "Hello"
bs := []byte(s)   // convert to bytes for manipulation
bs[0] = 'h'       // modify
s2 := string(bs)  // convert back

// Multiple assignment
srcIP, dstIP := "192.168.1.1", "10.0.0.1"
port, err := parsePort("443")
```

---

## Functions

```go
// Basic function
func greet(name string) string {
    return "Hello, " + name
}

// Multiple return values (Go's killer feature)
func connect(host string, port int) (bool, error) {
    if host == "" {
        return false, fmt.Errorf("empty host")
    }
    // ... connection logic
    return true, nil
}

// Using multiple returns
ok, err := connect("example.com", 443)
if err != nil {
    fmt.Println("Error:", err)
    return
}
if ok {
    fmt.Println("Connected!")
}

// Variadic functions (variable number of args)
func scanPorts(host string, ports ...int) {
    for _, port := range ports {
        fmt.Printf("Scanning %s:%d\n", host, port)
    }
}

scanPorts("192.168.1.1", 22, 80, 443, 8080)

// Named return values
func parseBanner(data []byte) (service string, version string) {
    // can use bare return
    service = "HTTP"
    version = "1.1"
    return
}
```

---

## Control Flow

```go
// If/else — always with braces
if port < 1024 {
    fmt.Println("Privileged port")
} else if port > 49151 {
    fmt.Println("Ephemeral port")
} else {
    fmt.Println("Registered port")
}

// If with init statement (very common in Go)
if conn, err := net.Dial("tcp", "example.com:80"); err != nil {
    fmt.Println("Failed:", err)
} else {
    defer conn.Close()
    fmt.Println("Connected:", conn.RemoteAddr())
}

// Switch
switch protocol {
case "TCP":
    fmt.Println("Reliable")
case "UDP":
    fmt.Println("Fast but unreliable")
case "ICMP":
    fmt.Println("Diagnostic")
default:
    fmt.Println("Unknown")
}

// Switch with no condition (acts like if/else chain)
switch {
case port == 22:
    fmt.Println("SSH")
case port == 80 || port == 8080:
    fmt.Println("HTTP")
case port == 443:
    fmt.Println("HTTPS")
}

// For loops (Go only has for, no while)
for i := 0; i < 10; i++ {
    fmt.Println(i)
}

// While equivalent
count := 0
for count < 100 {
    count++
}

// Infinite loop
for {
    // break to exit
    if shouldStop {
        break
    }
}

// Range — iterate over slices, maps, channels
ports := []int{22, 80, 443, 8080}
for i, port := range ports {
    fmt.Printf("ports[%d] = %d\n", i, port)
}

// Ignore index with _
for _, port := range ports {
    fmt.Println("Scanning:", port)
}
```

---

## Data Structures

### Slices (Dynamic Arrays)

```go
// Slice — the most important collection in Go
var ips []string                      // nil slice
ips = append(ips, "192.168.1.1")     // add element
ips = append(ips, "192.168.1.2")

// Literal
openPorts := []int{22, 80, 443}

// Make with capacity
results := make([]string, 0, 1000)  // len=0, cap=1000

// Slice operations
all := []int{1, 2, 3, 4, 5}
first := all[0]       // 1
last := all[len(all)-1] // 5
middle := all[1:4]    // [2, 3, 4]

// Filter pattern (common in security tools)
var suspicious []string
for _, ip := range ips {
    if isPrivate(ip) {
        suspicious = append(suspicious, ip)
    }
}
```

### Maps (Hash Tables)

```go
// Map — key/value store
portServices := map[int]string{
    22:   "SSH",
    80:   "HTTP",
    443:  "HTTPS",
    3306: "MySQL",
}

// Add/update
portServices[8080] = "HTTP-alt"

// Lookup with existence check
if service, ok := portServices[port]; ok {
    fmt.Println("Service:", service)
} else {
    fmt.Println("Unknown service")
}

// Iterate
for port, service := range portServices {
    fmt.Printf("%d → %s\n", port, service)
}

// Delete
delete(portServices, 8080)

// Count occurrences (common for log analysis)
ipCount := make(map[string]int)
for _, entry := range logEntries {
    ipCount[entry.IP]++
}
```

### Structs

```go
// Struct — group related data
type ScanResult struct {
    Host    string
    Port    int
    State   string  // "open", "closed", "filtered"
    Service string
    Banner  string
}

// Create instance
result := ScanResult{
    Host:    "192.168.1.1",
    Port:    22,
    State:   "open",
    Service: "SSH",
    Banner:  "OpenSSH_8.9",
}

// Access fields
fmt.Printf("%s:%d is %s (%s)\n", result.Host, result.Port, result.State, result.Service)

// Methods on structs
func (r ScanResult) IsVulnerable() bool {
    return r.State == "open" && r.Service == "Telnet"
}

// Pointer receiver (can modify the struct)
func (r *ScanResult) SetBanner(b string) {
    r.Banner = b
}
```

---

## Error Handling

Go doesn't have exceptions. Errors are values returned as the last return value.

```go
import (
    "errors"
    "fmt"
)

// Return error
func parsePort(s string) (int, error) {
    port, err := strconv.Atoi(s)
    if err != nil {
        return 0, fmt.Errorf("invalid port %q: %w", s, err)
    }
    if port < 1 || port > 65535 {
        return 0, fmt.Errorf("port %d out of range", port)
    }
    return port, nil
}

// Custom error type
type ScanError struct {
    Host    string
    Port    int
    Message string
}

func (e *ScanError) Error() string {
    return fmt.Sprintf("scan error on %s:%d: %s", e.Host, e.Port, e.Message)
}

// Sentinel errors
var ErrTimeout = errors.New("connection timeout")
var ErrRefused = errors.New("connection refused")

// Use errors.Is for checking
if errors.Is(err, ErrTimeout) {
    // handle timeout specifically
}

// The pattern: always check errors
conn, err := net.Dial("tcp", "192.168.1.1:80")
if err != nil {
    // DON'T ignore errors in security tools — they tell you things
    log.Printf("Connection failed: %v", err)
    return
}
defer conn.Close()
```

---

## Defer, Panic, and Recover

```go
// defer — runs when function exits (even on error)
func readFile(path string) {
    f, err := os.Open(path)
    if err != nil {
        return
    }
    defer f.Close()  // always closes, even if we return early
    
    // read file...
}

// Multiple defers execute LIFO (last in, first out)
func connect() {
    fmt.Println("connect")
    defer fmt.Println("cleanup 1")
    defer fmt.Println("cleanup 2")
    fmt.Println("working")
}
// Output: connect → working → cleanup 2 → cleanup 1

// panic — for truly unrecoverable errors
func mustHavePort(port int) {
    if port < 1 {
        panic("invalid port: must be > 0")
    }
}

// recover — catch a panic (use sparingly)
func safeOperation() {
    defer func() {
        if r := recover(); r != nil {
            fmt.Println("Recovered from panic:", r)
        }
    }()
    panic("something went wrong")
}
```

---

## Packages and Imports

```go
// All Go files start with package declaration
package scanner

import (
    "fmt"           // standard library
    "net"           // networking
    "os"            // operating system
    "time"          // time functions
    
    "github.com/some/external-package"  // third-party
)

// Exported names start with uppercase
func ScanHost(host string) []ScanResult {
    // exported — accessible from other packages
}

func checkPort(host string, port int) bool {
    // unexported — internal only
}
```

---

## A Complete Mini Security Tool

```go
package main

import (
    "flag"
    "fmt"
    "net"
    "os"
    "sync"
    "time"
)

type Result struct {
    Port  int
    Open  bool
    Error string
}

func scanPort(host string, port int, timeout time.Duration) Result {
    address := fmt.Sprintf("%s:%d", host, port)
    conn, err := net.DialTimeout("tcp", address, timeout)
    if err != nil {
        return Result{Port: port, Open: false, Error: err.Error()}
    }
    conn.Close()
    return Result{Port: port, Open: true}
}

func main() {
    host := flag.String("host", "127.0.0.1", "Target host")
    startPort := flag.Int("start", 1, "Start port")
    endPort := flag.Int("end", 1024, "End port")
    workers := flag.Int("workers", 100, "Number of concurrent workers")
    timeout := flag.Duration("timeout", 2*time.Second, "Connection timeout")
    flag.Parse()

    if *endPort < *startPort {
        fmt.Fprintln(os.Stderr, "end must be >= start")
        os.Exit(1)
    }

    fmt.Printf("Scanning %s ports %d-%d...\n", *host, *startPort, *endPort)

    jobs := make(chan int, *workers)
    results := make(chan Result, *workers)
    var wg sync.WaitGroup

    // Start workers
    for i := 0; i < *workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for port := range jobs {
                results <- scanPort(*host, port, *timeout)
            }
        }()
    }

    // Feed jobs
    go func() {
        for port := *startPort; port <= *endPort; port++ {
            jobs <- port
        }
        close(jobs)
    }()

    // Close results when done
    go func() {
        wg.Wait()
        close(results)
    }()

    // Collect and print
    openCount := 0
    for r := range results {
        if r.Open {
            fmt.Printf("[OPEN] %s:%d\n", *host, r.Port)
            openCount++
        }
    }

    fmt.Printf("\nFound %d open ports\n", openCount)
}
```

```bash
go run main.go -host 192.168.1.1 -start 1 -end 10000 -workers 500
```

---

## Summary

| Concept | Quick reference |
|---------|----------------|
| Variables | `x := 5` or `var x int = 5` |
| Multiple returns | `val, err := func()` |
| Error handling | `if err != nil { return err }` |
| Loops | `for`, `for i < n`, `for _, v := range` |
| Slices | `append(s, v)`, `s[1:3]` |
| Maps | `m["key"]`, `m["key"] = val`, `v, ok := m["key"]` |
| Structs | `type Name struct {}`, `Name{Field: value}` |
| Defer | Cleanup code that always runs |
| Goroutines | Chapter 11 |

---

## Exercises

1. Write a function `parseIPRange("192.168.1.0/24")` that returns all host IPs in the range as a `[]string`.
2. Write a function that reads a file of IP addresses (one per line) and returns them as a slice.
3. Extend the mini scanner above to also print the service name for known ports (use a map).
4. Add a `--json` flag to the scanner that outputs results as JSON instead of text.
5. Make the scanner save results to a file using `os.WriteFile`.
