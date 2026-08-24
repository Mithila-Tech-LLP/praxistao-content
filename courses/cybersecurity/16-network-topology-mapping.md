# Chapter 16: Network Topology Mapping — Discovering the Full Attack Surface

*A single IP tells you one target. A network map tells you the entire battlefield. Topology mapping uncovers every host, subnet, router, and trust relationship — the intelligence an attacker needs before striking.*

---

## What is Network Topology Mapping?

Network topology mapping answers:
- **What hosts exist?** (host discovery)
- **How are they connected?** (routing, VLANs)
- **What services are they running?** (port scanning)
- **What OS are they running?** (OS fingerprinting)
- **What trusts exist?** (which machines can reach which)

In a real engagement, you build this map incrementally:
1. Ping sweep → live hosts
2. Port scan → services on live hosts
3. OS fingerprint → attack surface per host
4. Service version → CVE matching
5. Route tracing → network structure

---

## Host Discovery

### ICMP Ping Sweep

```bash
# Classic ping sweep
for i in $(seq 1 254); do
    ping -c 1 -W 1 192.168.1.$i | grep "bytes from" &
done
wait

# Nmap ping sweep (faster)
nmap -sn 192.168.1.0/24

# No ICMP? Try TCP SYN to common ports
nmap -sn -PS22,80,443 192.168.1.0/24

# ARP scan (LAN only, most reliable on local network)
nmap -PR -sn 192.168.1.0/24
arp-scan -l         # local network
arp-scan 192.168.1.0/24
```

### Go: ICMP Ping Sweep

```go
package main

import (
    "fmt"
    "net"
    "sync"
    "time"
)

// TCP-based host discovery (works without root, avoids ICMP filtering)
func isHostUp(ip string, ports []int, timeout time.Duration) bool {
    for _, port := range ports {
        addr := fmt.Sprintf("%s:%d", ip, port)
        conn, err := net.DialTimeout("tcp", addr, timeout)
        if err == nil {
            conn.Close()
            return true
        }
    }
    return false
}

func discoverHosts(subnet string, workers int) []string {
    // Parse CIDR
    ip, ipNet, err := net.ParseCIDR(subnet)
    if err != nil {
        return nil
    }
    
    var targets []string
    for ip := ip.Mask(ipNet.Mask); ipNet.Contains(ip); incrementIP(ip) {
        targets = append(targets, ip.String())
    }
    
    // Common ports to probe for host discovery
    probePorts := []int{22, 80, 443, 445, 3389}
    
    jobs := make(chan string, workers)
    results := make(chan string, workers)
    var wg sync.WaitGroup
    
    for i := 0; i < workers; i++ {
        wg.Add(1)
        go func() {
            defer wg.Done()
            for target := range jobs {
                if isHostUp(target, probePorts, 500*time.Millisecond) {
                    results <- target
                }
            }
        }()
    }
    
    go func() {
        for _, t := range targets {
            jobs <- t
        }
        close(jobs)
    }()
    
    go func() {
        wg.Wait()
        close(results)
    }()
    
    var live []string
    for host := range results {
        live = append(live, host)
        fmt.Printf("[UP] %s\n", host)
    }
    return live
}

func incrementIP(ip net.IP) {
    for j := len(ip) - 1; j >= 0; j-- {
        ip[j]++
        if ip[j] != 0 {
            break
        }
    }
}

func main() {
    fmt.Println("Host discovery on 192.168.1.0/24")
    live := discoverHosts("192.168.1.0/24", 100)
    fmt.Printf("\nFound %d live hosts\n", len(live))
}
```

---

## OS Fingerprinting

Different operating systems implement TCP differently. Nmap uses these quirks to identify them.

```bash
# OS detection (requires root for raw sockets)
sudo nmap -O 192.168.1.1

# TTL-based quick guess
ping -c 1 192.168.1.1 | grep ttl
# TTL 64  → Linux/macOS
# TTL 128 → Windows
# TTL 255 → Cisco/network device

# Nmap OS detection with service scan
sudo nmap -O -sV 192.168.1.1
```

### OS Fingerprinting from TTL

```go
package main

import (
    "fmt"
    "os/exec"
    "regexp"
    "strconv"
    "strings"
)

var ttlPattern = regexp.MustCompile(`ttl=(\d+)`)

func guessOS(ip string) string {
    out, err := exec.Command("ping", "-c", "1", "-W", "1", ip).Output()
    if err != nil {
        return "unreachable"
    }
    
    m := ttlPattern.FindStringSubmatch(strings.ToLower(string(out)))
    if m == nil {
        return "unknown"
    }
    
    ttl, _ := strconv.Atoi(m[1])
    switch {
    case ttl <= 64:
        return "Linux/macOS/FreeBSD"
    case ttl <= 128:
        return "Windows"
    case ttl <= 255:
        return "Network Device (Cisco/Juniper)"
    default:
        return "unknown"
    }
}

func main() {
    hosts := []string{"192.168.1.1", "8.8.8.8", "1.1.1.1"}
    for _, h := range hosts {
        fmt.Printf("%-20s → %s\n", h, guessOS(h))
    }
}
```

---

## Route Tracing

Traceroute reveals network topology — routers between you and target.

```bash
traceroute 8.8.8.8            # UDP-based (Linux default)
traceroute -T -p 443 8.8.8.8  # TCP SYN (bypasses some firewalls)
mtr 8.8.8.8                   # live continuous traceroute

# For internal networks
traceroute 10.0.50.1           # discovers internal subnets
# Each hop = a router or gateway
# Multiple hops within 10.x.x.x = multiple internal subnets
```

Reading a traceroute for network mapping:
```
Hop 1:  192.168.1.1     → your default gateway (home router)
Hop 2:  100.64.0.1      → ISP gateway
Hop 3:  10.0.50.1       → internal corporate router!
Hop 4:  172.16.20.1     → DMZ gateway!
Hop 5:  target.corp.com → destination
```

Detecting `10.x`, `172.16.x`, `192.168.x` hops reveals internal network structure.

---

## SNMP Reconnaissance

SNMP (Simple Network Management Protocol) leaks enormous amounts of network info.

```bash
# Check if SNMP is open
nmap -sU -p 161 192.168.1.0/24

# Enumerate with community string "public" (common default)
snmpwalk -v 2c -c public 192.168.1.1

# Specific OIDs
snmpget -v 2c -c public 192.168.1.1 1.3.6.1.2.1.1.1.0  # sysDescr (OS info)
snmpget -v 2c -c public 192.168.1.1 1.3.6.1.2.1.1.5.0  # sysName (hostname)
snmpwalk -v 2c -c public 192.168.1.1 1.3.6.1.2.1.4.20  # IP addresses
snmpwalk -v 2c -c public 192.168.1.1 1.3.6.1.2.1.2.2   # interfaces
```

---

## Building a Network Map Report

```go
package main

import (
    "encoding/json"
    "fmt"
    "os"
    "time"
)

type Host struct {
    IP          string              `json:"ip"`
    Hostname    string              `json:"hostname,omitempty"`
    OS          string              `json:"os,omitempty"`
    OpenPorts   []PortInfo          `json:"open_ports"`
    DiscoveredAt time.Time          `json:"discovered_at"`
}

type PortInfo struct {
    Port    int    `json:"port"`
    Proto   string `json:"proto"`
    Service string `json:"service"`
    Version string `json:"version,omitempty"`
}

type NetworkMap struct {
    Subnet     string    `json:"subnet"`
    ScannedAt  time.Time `json:"scanned_at"`
    LiveHosts  int       `json:"live_hosts"`
    Hosts      []Host    `json:"hosts"`
}

func generateReport(subnet string, hosts []Host) {
    m := NetworkMap{
        Subnet:    subnet,
        ScannedAt: time.Now(),
        LiveHosts: len(hosts),
        Hosts:     hosts,
    }
    
    // JSON report
    data, _ := json.MarshalIndent(m, "", "  ")
    os.WriteFile("network-map.json", data, 0644)
    
    // Text summary
    fmt.Printf("\n=== Network Map: %s ===\n", subnet)
    fmt.Printf("Scanned: %s\n", m.ScannedAt.Format("2006-01-02 15:04:05"))
    fmt.Printf("Live hosts: %d\n\n", m.LiveHosts)
    
    for _, h := range hosts {
        fmt.Printf("Host: %s", h.IP)
        if h.Hostname != "" {
            fmt.Printf(" (%s)", h.Hostname)
        }
        if h.OS != "" {
            fmt.Printf(" [%s]", h.OS)
        }
        fmt.Println()
        for _, p := range h.OpenPorts {
            fmt.Printf("  %d/%s  %-12s %s\n", p.Port, p.Proto, p.Service, p.Version)
        }
    }
}

func main() {
    // Demo data
    hosts := []Host{
        {
            IP: "192.168.1.1", Hostname: "router.local", OS: "Linux",
            OpenPorts: []PortInfo{
                {22, "tcp", "SSH", "OpenSSH 8.9"},
                {80, "tcp", "HTTP", "nginx 1.22"},
            },
        },
        {
            IP: "192.168.1.100", OS: "Windows",
            OpenPorts: []PortInfo{
                {445, "tcp", "SMB", "Windows 10"},
                {3389, "tcp", "RDP", ""},
            },
        },
    }
    generateReport("192.168.1.0/24", hosts)
}
```

---

## Summary

| Task | Tool | Go equivalent |
|------|------|--------------|
| Host discovery | `nmap -sn` | TCP probe workers |
| OS fingerprinting | `nmap -O` | TTL analysis |
| Service detection | `nmap -sV` | Banner grabbing |
| Route mapping | `traceroute` | `exec.Command("traceroute")` |
| SNMP enum | `snmpwalk` | Go SNMP library |

---

## Exercises

1. Write a complete Go network scanner: takes CIDR, finds live hosts, detects OS via TTL, scans top 20 ports, outputs JSON report
2. Parse a traceroute output and identify private IP ranges — these reveal internal network structure
3. Build an ARP scanner using Go's `syscall` package to send raw ARP requests on the local network
4. Write a script that runs nmap, parses the XML output, and generates a vulnerability prioritization list
