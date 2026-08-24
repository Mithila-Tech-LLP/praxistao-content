# Chapter 10: Go Data Structures — Slices, Maps, Structs, Interfaces

*Security tools deal with huge amounts of data: millions of log events, thousands of IP addresses, hundreds of open ports. The right data structure makes your tool fast and correct. The wrong one makes it slow and buggy.*

---

## Slices — The Core Collection

A slice is a dynamic array. It's the most-used data structure in Go security tools.

```go
package main

import (
    "fmt"
    "sort"
    "strings"
)

func main() {
    // Create
    var ports []int                      // nil slice
    ports = append(ports, 22, 80, 443)   // add elements
    
    ips := []string{"10.0.0.1", "10.0.0.2", "10.0.0.3"}
    
    // Access
    first := ips[0]           // "10.0.0.1"
    last  := ips[len(ips)-1]  // "10.0.0.3"
    sub   := ips[1:]          // ["10.0.0.2", "10.0.0.3"]
    
    // Length vs capacity
    results := make([]string, 0, 1000)  // len=0, cap=1000 (pre-allocate)
    fmt.Println(len(results), cap(results))  // 0 1000
    
    // Filter pattern (used everywhere in security tools)
    openPorts := []int{22, 80, 443, 8080, 3306}
    var webPorts []int
    for _, p := range openPorts {
        if p == 80 || p == 443 || p == 8080 {
            webPorts = append(webPorts, p)
        }
    }
    
    // Deduplication (common in IP/domain lists)
    seen := make(map[string]bool)
    var unique []string
    for _, ip := range []string{"1.1.1.1", "2.2.2.2", "1.1.1.1"} {
        if !seen[ip] {
            seen[ip] = true
            unique = append(unique, ip)
        }
    }
    fmt.Println(unique)  // [1.1.1.1 2.2.2.2]
    
    // Sort
    sort.Ints(openPorts)
    sort.Strings(ips)
    
    // Custom sort (by port number, then protocol)
    type Service struct { Port int; Proto string }
    services := []Service{{443, "HTTPS"}, {22, "SSH"}, {80, "HTTP"}}
    sort.Slice(services, func(i, j int) bool {
        return services[i].Port < services[j].Port
    })
    
    // Contains check (no built-in, use loop or map)
    target := "10.0.0.2"
    found := false
    for _, ip := range ips {
        if ip == target { found = true; break }
    }
    
    // Join for output
    fmt.Println(strings.Join(ips, ", "))
    
    _ = first; _ = last; _ = sub; _ = found
}
```

---

## Maps — Fast Lookups

Maps provide O(1) average lookup. Essential for IP counting, caching, deduplication.

```go
package main

import (
    "fmt"
    "sort"
)

// Count occurrences — core pattern for log analysis
func countIPs(logs []string) map[string]int {
    counts := make(map[string]int)
    for _, ip := range logs {
        counts[ip]++
    }
    return counts
}

// Top N by count
func topN(counts map[string]int, n int) []string {
    type kv struct { key string; val int }
    var pairs []kv
    for k, v := range counts {
        pairs = append(pairs, kv{k, v})
    }
    sort.Slice(pairs, func(i, j int) bool {
        return pairs[i].val > pairs[j].val
    })
    var result []string
    for i := 0; i < n && i < len(pairs); i++ {
        result = append(result, fmt.Sprintf("%s (%d)", pairs[i].key, pairs[i].val))
    }
    return result
}

// Nested maps — e.g., port → service → version
type ServiceInfo struct {
    Name    string
    Version string
    Banner  string
}

func main() {
    // Port → service info
    discovered := map[int]ServiceInfo{
        22:  {Name: "SSH",   Version: "OpenSSH 8.9"},
        80:  {Name: "HTTP",  Version: "nginx 1.18"},
        443: {Name: "HTTPS", Version: "nginx 1.18"},
    }
    
    // Safe lookup
    if svc, ok := discovered[22]; ok {
        fmt.Printf("SSH: %s\n", svc.Version)
    }
    
    // Delete
    delete(discovered, 443)
    
    // Iterate
    for port, svc := range discovered {
        fmt.Printf("%d: %s %s\n", port, svc.Name, svc.Version)
    }
    
    // Log analysis
    accessLog := []string{
        "1.2.3.4", "5.6.7.8", "1.2.3.4", "1.2.3.4", "9.9.9.9",
    }
    counts := countIPs(accessLog)
    fmt.Println("Top IPs:", topN(counts, 3))
    
    // Set using map[T]struct{} (zero memory overhead)
    blocklist := map[string]struct{}{
        "185.220.101.5": {},
        "45.33.32.156":  {},
    }
    ip := "185.220.101.5"
    if _, blocked := blocklist[ip]; blocked {
        fmt.Println(ip, "is blocked")
    }
}
```

---

## Structs — Modeling Security Data

```go
package main

import (
    "fmt"
    "time"
)

// Hierarchical struct — mirrors real security events
type Base struct {
    ID        string
    Timestamp time.Time
    AgentID   string
    Hostname  string
}

type Alert struct {
    Base
    RuleID      string
    RuleName    string
    Severity    string  // low, medium, high, critical
    Description string
    MITRE       string
    Resolved    bool
    ResolvedAt  *time.Time  // nil if not resolved
}

type ScanResult struct {
    Host     string
    Port     int
    Protocol string
    State    string
    Service  string
    Version  string
    Banner   string
    ScanTime time.Duration
}

func (r ScanResult) String() string {
    return fmt.Sprintf("%s:%d/%s %s (%s %s)",
        r.Host, r.Port, r.Protocol, r.State, r.Service, r.Version)
}

func (r ScanResult) IsVulnerable() bool {
    // Simple version-based check
    vulnerable := map[string]string{
        "vsftpd":   "2.3.4",
        "OpenSSH":  "7.6",
    }
    if maxVer, ok := vulnerable[r.Service]; ok {
        return r.Version <= maxVer
    }
    return false
}

// Embedding — composition over inheritance
type TLSInfo struct {
    Version     string
    CipherSuite string
    Certificate string
    ExpiresAt   time.Time
}

type HTTPSResult struct {
    ScanResult           // embed all ScanResult fields
    TLS        TLSInfo
    Headers    map[string]string
    StatusCode int
}

func main() {
    result := ScanResult{
        Host:     "192.168.1.1",
        Port:     22,
        Protocol: "tcp",
        State:    "open",
        Service:  "SSH",
        Version:  "OpenSSH 7.4",
    }
    
    fmt.Println(result)
    fmt.Println("Vulnerable:", result.IsVulnerable())
    
    // Slice of structs — report
    results := []ScanResult{
        {Host: "10.0.0.1", Port: 22, State: "open", Service: "SSH", Version: "OpenSSH 8.9"},
        {Host: "10.0.0.1", Port: 80, State: "open", Service: "HTTP"},
        {Host: "10.0.0.2", Port: 3306, State: "open", Service: "MySQL", Version: "5.7.39"},
    }
    
    fmt.Println("\nOpen services:")
    for _, r := range results {
        flag := ""
        if r.IsVulnerable() {
            flag = " [VULNERABLE]"
        }
        fmt.Printf("  %s%s\n", r, flag)
    }
}
```

---

## Interfaces — Polymorphism for Security Tools

Interfaces let you write code that works with multiple types. This powers flexible, extensible security tools.

```go
package main

import (
    "fmt"
    "io"
    "os"
    "strings"
)

// Scanner interface — any scan target implements this
type Scanner interface {
    Scan() ([]ScanResult, error)
    Name() string
}

type ScanResult struct {
    Finding  string
    Severity string
    Detail   string
}

// TCP port scanner
type PortScanner struct {
    Host  string
    Ports []int
}

func (p *PortScanner) Name() string { return "Port Scanner" }
func (p *PortScanner) Scan() ([]ScanResult, error) {
    // ... actual scanning logic
    return []ScanResult{{
        Finding:  fmt.Sprintf("Open port on %s:22", p.Host),
        Severity: "info",
    }}, nil
}

// Web vulnerability scanner
type WebScanner struct {
    URL string
}

func (w *WebScanner) Name() string { return "Web Scanner" }
func (w *WebScanner) Scan() ([]ScanResult, error) {
    return []ScanResult{{
        Finding:  fmt.Sprintf("Missing security headers on %s", w.URL),
        Severity: "medium",
    }}, nil
}

// Run any scanner — works for all types
func RunAll(scanners []Scanner) {
    for _, s := range scanners {
        fmt.Printf("[*] Running %s...\n", s.Name())
        results, err := s.Scan()
        if err != nil {
            fmt.Printf("  Error: %v\n", err)
            continue
        }
        for _, r := range results {
            fmt.Printf("  [%s] %s\n", strings.ToUpper(r.Severity), r.Finding)
        }
    }
}

// io.Writer interface — redirect output anywhere
type JSONReporter struct {
    w io.Writer
}

func NewJSONReporter(w io.Writer) *JSONReporter { return &JSONReporter{w} }
func (r *JSONReporter) Report(findings []ScanResult) {
    for _, f := range findings {
        fmt.Fprintf(r.w, `{"severity":%q,"finding":%q}`+"\n", f.Severity, f.Finding)
    }
}

func main() {
    scanners := []Scanner{
        &PortScanner{Host: "10.0.0.1", Ports: []int{22, 80, 443}},
        &WebScanner{URL: "http://10.0.0.1"},
    }
    RunAll(scanners)
    
    // Output to file or stdout
    reporter := NewJSONReporter(os.Stdout)
    reporter.Report([]ScanResult{
        {Finding: "SQL injection found", Severity: "critical"},
    })
    
    // Same reporter works with any writer
    var buf strings.Builder
    reporter2 := NewJSONReporter(&buf)
    reporter2.Report([]ScanResult{{Finding: "XSS", Severity: "high"}})
    fmt.Println(buf.String())
}
```

---

## JSON — Reading and Writing Security Tool Data

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type Vulnerability struct {
    CVE         string    `json:"cve"`
    Service     string    `json:"service"`
    Host        string    `json:"host"`
    Port        int       `json:"port"`
    Severity    string    `json:"severity"`
    Description string    `json:"description,omitempty"`
    Discovered  time.Time `json:"discovered"`
}

func main() {
    // Marshal (struct → JSON)
    vuln := Vulnerability{
        CVE:        "CVE-2021-44228",
        Service:    "Log4j",
        Host:       "192.168.1.100",
        Port:       8080,
        Severity:   "critical",
        Discovered: time.Now(),
    }
    
    data, err := json.MarshalIndent(vuln, "", "  ")
    if err != nil {
        panic(err)
    }
    fmt.Println(string(data))
    
    // Unmarshal (JSON → struct)
    jsonStr := `{"cve":"CVE-2021-44228","service":"Log4j","host":"10.0.0.5","port":8443,"severity":"critical"}`
    var v2 Vulnerability
    if err := json.Unmarshal([]byte(jsonStr), &v2); err != nil {
        panic(err)
    }
    fmt.Printf("Parsed: %s on %s:%d\n", v2.CVE, v2.Host, v2.Port)
    
    // Stream JSON (large reports — memory efficient)
    enc := json.NewEncoder(os.Stdout)
    for _, v := range []Vulnerability{vuln, v2} {
        enc.Encode(v)  // writes one JSON line per call
    }
    
    // Parse unknown JSON structure
    var raw map[string]interface{}
    json.Unmarshal([]byte(`{"severity":"high","data":{"port":443}}`), &raw)
    fmt.Println(raw["severity"])  // high
}
```

---

## Summary

| Structure | Best for | Avoid when |
|-----------|---------|-----------|
| `[]T` slice | Ordered lists, port lists, results | Lookup by key |
| `map[K]V` | Counts, caches, dedup | Ordered iteration |
| `struct` | Structured data, events | Simple key-value |
| `interface` | Pluggable scanners, reporters | Simple single-type code |

---

## Exercises

1. Write a function `ParseTargets(input string)` that accepts `"192.168.1.0/24"`, `"192.168.1.1-10"`, or `"192.168.1.1,1.1.1.1"` and returns a `[]string` of individual IPs.
2. Build a `VulnDB` type backed by a `map[string][]Vulnerability` (keyed by host). Add methods: `Add`, `GetByHost`, `GetBySeverity`, `Summary`.
3. Write a `Scanner` interface with two implementations: `MockScanner` (returns hardcoded results) and `FileScanner` (reads results from JSON file). Write a `RunAndReport` function that accepts any `Scanner`.
4. Build a log line struct for nginx access logs. Parse 1000 lines and find the top 5 IPs by request count.
