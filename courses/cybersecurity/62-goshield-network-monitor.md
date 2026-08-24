# Chapter 62: GoShield — Network Connection Monitor

*The network collector is the most critical GoShield component for detecting C2 (command-and-control) beaconing, lateral movement, and data exfiltration. Every suspicious connection tells a story.*

---

## What We're Building

A network connection monitor that:
1. Tracks all TCP/UDP connections (outbound and listening)
2. Detects new connections and closed connections
3. Resolves domain names for IP addresses (reverse DNS)
4. Identifies the process behind each connection
5. Fires alerts for suspicious patterns:
   - New outbound connections to rare or new destinations
   - High-frequency beaconing (C2 pattern)
   - Connections to known-bad ports (IRC, unusual high ports)
   - Large data transfers (potential exfiltration)

---

## Directory Structure

```
goshield/
├── agent/
│   ├── collectors/
│   │   ├── file.go      (Chapter 60)
│   │   ├── process.go   (Chapter 61)
│   │   └── network.go   ← this chapter
│   └── main.go
├── ...
```

---

## The Network Watcher

```go
// agent/collectors/network.go
package collectors

import (
    "context"
    "fmt"
    "net"
    "strings"
    "sync"
    "time"

    "github.com/goshield/agent/events"
    "github.com/shirou/gopsutil/v3/net"
    "github.com/shirou/gopsutil/v3/process"
)

type NetworkWatchConfig struct {
    PollInterval    time.Duration
    AlertCallback   func(events.NetworkEvent)
    EventCallback   func(events.NetworkEvent)
    
    // Alerting thresholds
    BeaconingThresholdConnections int           // Alert if >N connections to same IP in window
    BeaconingWindow               time.Duration  // Time window for beaconing detection
    SuspiciousPorts               []int          // Ports that are always suspicious
    
    // Resolved DNS cache TTL
    DNSCacheTTL time.Duration
}

func DefaultNetworkWatchConfig() NetworkWatchConfig {
    return NetworkWatchConfig{
        PollInterval:                  5 * time.Second,
        BeaconingThresholdConnections: 20,
        BeaconingWindow:               60 * time.Second,
        SuspiciousPorts:               []int{4444, 1337, 31337, 6667, 6668, 6669, 4321},
        DNSCacheTTL:                   5 * time.Minute,
    }
}

// connectionKey uniquely identifies a connection
type connectionKey struct {
    LocalAddr  string
    RemoteAddr string
    PID        int32
    Status     string
}

// dnsEntry is a cached DNS lookup
type dnsEntry struct {
    Domain  string
    Cached  time.Time
}

// beaconTracker tracks connection frequency to detect C2 beaconing
type beaconTracker struct {
    mu          sync.Mutex
    connections map[string][]time.Time  // remoteIP → connection timestamps
}

func (bt *beaconTracker) record(ip string, window time.Duration) int {
    bt.mu.Lock()
    defer bt.mu.Unlock()
    
    now := time.Now()
    times := bt.connections[ip]
    
    // Remove old entries outside window
    cutoff := now.Add(-window)
    var recent []time.Time
    for _, t := range times {
        if t.After(cutoff) {
            recent = append(recent, t)
        }
    }
    
    recent = append(recent, now)
    bt.connections[ip] = recent
    return len(recent)
}

type NetworkWatcher struct {
    config     NetworkWatchConfig
    ctx        context.Context
    cancel     context.CancelFunc
    
    // State: map of currently known connections
    knownConns map[connectionKey]bool
    mu         sync.RWMutex
    
    // DNS cache
    dnsCache   map[string]dnsEntry
    dnsMu      sync.RWMutex
    
    // Beaconing tracker
    beacon     *beaconTracker
}

func NewNetworkWatcher(config NetworkWatchConfig) *NetworkWatcher {
    ctx, cancel := context.WithCancel(context.Background())
    return &NetworkWatcher{
        config:     config,
        ctx:        ctx,
        cancel:     cancel,
        knownConns: make(map[connectionKey]bool),
        dnsCache:   make(map[string]dnsEntry),
        beacon: &beaconTracker{
            connections: make(map[string][]time.Time),
        },
    }
}

func (nw *NetworkWatcher) Start() {
    go nw.poll()
}

func (nw *NetworkWatcher) Stop() {
    nw.cancel()
}

func (nw *NetworkWatcher) poll() {
    ticker := time.NewTicker(nw.config.PollInterval)
    defer ticker.Stop()
    
    // First pass: establish baseline (don't alert for existing connections)
    nw.snapshot()
    
    for {
        select {
        case <-nw.ctx.Done():
            return
        case <-ticker.C:
            nw.detectChanges()
        }
    }
}

// snapshot captures current connections without alerting
func (nw *NetworkWatcher) snapshot() {
    conns, err := gopsutil_net.Connections("all")
    if err != nil {
        return
    }
    
    nw.mu.Lock()
    defer nw.mu.Unlock()
    
    for _, conn := range conns {
        key := nw.makeKey(conn)
        nw.knownConns[key] = true
    }
}

// detectChanges finds new and closed connections
func (nw *NetworkWatcher) detectChanges() {
    conns, err := gopsutil_net.Connections("all")
    if err != nil {
        return
    }
    
    nw.mu.Lock()
    defer nw.mu.Unlock()
    
    currentConns := make(map[connectionKey]bool)
    
    for _, conn := range conns {
        key := nw.makeKey(conn)
        currentConns[key] = true
        
        if !nw.knownConns[key] {
            // New connection
            event := nw.buildEvent(conn, "connect")
            nw.config.EventCallback(event)
            
            // Check for suspicious patterns
            nw.analyzeConnection(conn, event)
        }
    }
    
    // Find closed connections
    for key := range nw.knownConns {
        if !currentConns[key] {
            event := events.NetworkEvent{
                Base:      events.NewBase("network"),
                Action:    "disconnect",
                Protocol:  key.Status,
                SrcIP:     parseHost(key.LocalAddr),
                SrcPort:   parsePort(key.LocalAddr),
                DstIP:     parseHost(key.RemoteAddr),
                DstPort:   parsePort(key.RemoteAddr),
                PID:       int(key.PID),
            }
            nw.config.EventCallback(event)
        }
    }
    
    nw.knownConns = currentConns
}

func (nw *NetworkWatcher) makeKey(conn gopsutil_net.ConnectionStat) connectionKey {
    return connectionKey{
        LocalAddr:  fmt.Sprintf("%s:%d", conn.Laddr.IP, conn.Laddr.Port),
        RemoteAddr: fmt.Sprintf("%s:%d", conn.Raddr.IP, conn.Raddr.Port),
        PID:        conn.Pid,
        Status:     conn.Status,
    }
}

func (nw *NetworkWatcher) buildEvent(conn gopsutil_net.ConnectionStat, action string) events.NetworkEvent {
    remoteIP := conn.Raddr.IP
    domain := nw.reverseLookup(remoteIP)
    processName := nw.getProcessName(conn.Pid)
    
    proto := "TCP"
    if conn.Type == 2 {
        proto = "UDP"
    }
    
    return events.NetworkEvent{
        Base:        events.NewBase("network"),
        Action:      action,
        Protocol:    proto,
        SrcIP:       conn.Laddr.IP,
        SrcPort:     int(conn.Laddr.Port),
        DstIP:       remoteIP,
        DstPort:     int(conn.Raddr.Port),
        Domain:      domain,
        PID:         int(conn.Pid),
        Process:     processName,
    }
}

// reverseLookup resolves IP to hostname with caching
func (nw *NetworkWatcher) reverseLookup(ip string) string {
    if ip == "" || ip == "0.0.0.0" || ip == "::" {
        return ""
    }
    
    nw.dnsMu.RLock()
    if entry, ok := nw.dnsCache[ip]; ok {
        if time.Since(entry.Cached) < nw.config.DNSCacheTTL {
            nw.dnsMu.RUnlock()
            return entry.Domain
        }
    }
    nw.dnsMu.RUnlock()
    
    // Do the lookup
    names, err := net.LookupAddr(ip)
    domain := ""
    if err == nil && len(names) > 0 {
        domain = strings.TrimSuffix(names[0], ".")
    }
    
    nw.dnsMu.Lock()
    nw.dnsCache[ip] = dnsEntry{Domain: domain, Cached: time.Now()}
    nw.dnsMu.Unlock()
    
    return domain
}

func (nw *NetworkWatcher) getProcessName(pid int32) string {
    if pid == 0 {
        return "unknown"
    }
    p, err := process.NewProcess(pid)
    if err != nil {
        return "unknown"
    }
    name, err := p.Name()
    if err != nil {
        return "unknown"
    }
    return name
}

// analyzeConnection checks for suspicious patterns
func (nw *NetworkWatcher) analyzeConnection(conn gopsutil_net.ConnectionStat, event events.NetworkEvent) {
    // Skip listening and established local connections
    if conn.Raddr.IP == "" || conn.Raddr.IP == "0.0.0.0" || conn.Status == "LISTEN" {
        return
    }
    
    // Skip loopback
    if strings.HasPrefix(conn.Raddr.IP, "127.") || conn.Raddr.IP == "::1" {
        return
    }
    
    // Check 1: Suspicious destination ports
    for _, port := range nw.config.SuspiciousPorts {
        if int(conn.Raddr.Port) == port {
            nw.fireAlert(event, "NET-001", "Suspicious port connection",
                fmt.Sprintf("Process %s connected to suspicious port %d", event.Process, port),
                "high", "T1043")
            return
        }
    }
    
    // Check 2: Beaconing detection
    count := nw.beacon.record(conn.Raddr.IP, nw.config.BeaconingWindow)
    if count >= nw.config.BeaconingThresholdConnections {
        nw.fireAlert(event, "NET-002", "Potential C2 Beaconing",
            fmt.Sprintf("Process %s made %d connections to %s in %s (possible C2 beaconing)",
                event.Process, count, conn.Raddr.IP, nw.config.BeaconingWindow),
            "critical", "T1071")
        return
    }
    
    // Check 3: High-numbered destination ports (common for C2)
    if conn.Raddr.Port > 49151 && !isKnownEphemeral(conn.Laddr.Port) {
        nw.fireAlert(event, "NET-003", "Connection to ephemeral port",
            fmt.Sprintf("Outbound connection to high port %d (possible C2 or reverse shell)",
                conn.Raddr.Port),
            "medium", "T1095")
    }
    
    // Check 4: Unexpected outbound from server processes
    serverProcesses := map[string]bool{
        "nginx": true, "apache2": true, "httpd": true,
        "mysqld": true, "postgres": true, "redis-server": true,
    }
    if serverProcesses[event.Process] {
        nw.fireAlert(event, "NET-004", "Server process making outbound connection",
            fmt.Sprintf("Server process %s made unexpected outbound connection to %s:%d — possible webshell/compromise",
                event.Process, conn.Raddr.IP, conn.Raddr.Port),
            "high", "T1059")
    }
}

func (nw *NetworkWatcher) fireAlert(event events.NetworkEvent, ruleID, ruleName, description, severity, mitre string) {
    alertEvent := events.NetworkEvent{
        Base:     event.Base,
        Action:   event.Action,
        Protocol: event.Protocol,
        SrcIP:    event.SrcIP,
        SrcPort:  event.SrcPort,
        DstIP:    event.DstIP,
        DstPort:  event.DstPort,
        Domain:   event.Domain,
        PID:      event.PID,
        Process:  event.Process,
    }
    _ = alertEvent // callback would include rule info
    
    if nw.config.AlertCallback != nil {
        nw.config.AlertCallback(event)
    }
}

// isKnownEphemeral checks if a port is a client ephemeral port
func isKnownEphemeral(port uint32) bool {
    return port >= 32768 && port <= 60999
}

func parseHost(addr string) string {
    host, _, err := net.SplitHostPort(addr)
    if err != nil {
        return addr
    }
    return host
}

func parsePort(addr string) int {
    _, portStr, err := net.SplitHostPort(addr)
    if err != nil {
        return 0
    }
    var port int
    fmt.Sscanf(portStr, "%d", &port)
    return port
}
```

---

## Detection Rules for Network Events

```yaml
# network-rules.yaml

rules:
  - id: NET-001
    name: "Reverse Shell — Common Ports"
    severity: critical
    description: "Outbound connection to ports commonly used for reverse shells"
    event_type: network
    conditions:
      - field: action
        op: equals
        value: connect
      - field: dst_port
        op: equalsAny
        values: ["4444", "4445", "9001", "9999", "1337", "31337"]
    mitre: "T1059.004"
    
  - id: NET-002
    name: "DNS over Non-Standard Port"
    severity: high
    description: "DNS query to non-standard port — possible DNS tunneling (C2 over DNS)"
    event_type: network
    conditions:
      - field: dst_port
        op: equals
        value: "53"
        negate: true
      - field: protocol
        op: equals
        value: UDP
      - field: dst_port
        op: not
        value: "53"
    mitre: "T1071.004"
    
  - id: NET-003
    name: "IRC Traffic"
    severity: high
    description: "IRC traffic detected — old C2 technique, unusual in enterprise"
    event_type: network
    conditions:
      - field: dst_port
        op: equalsAny
        values: ["6667", "6668", "6669", "7000"]
    mitre: "T1095"
    
  - id: NET-004
    name: "TOR Exit Node Connection"
    severity: high
    description: "Connection to known TOR exit node or proxy"
    event_type: network
    conditions:
      - field: dst_port
        op: equalsAny
        values: ["9001", "9030", "9050", "9051"]
    mitre: "T1090.003"
    
  - id: NET-005
    name: "Large Outbound Transfer"
    severity: medium
    description: "Unusually large outbound network transfer — possible exfiltration"
    event_type: network
    threshold:
      field: process
      count: 1000    # 1000 connections from same process
      window: 60s    # in 60 seconds
    mitre: "T1048"
```

---

## Integration with Agent Main Loop

```go
// agent/main.go — adding network watcher
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/goshield/agent/collectors"
    "github.com/goshield/agent/events"
    "github.com/goshield/agent/transport"
)

func main() {
    serverURL := os.Getenv("GOSHIELD_SERVER")
    if serverURL == "" {
        serverURL = "http://localhost:8080"
    }
    
    client := transport.NewClient(serverURL)
    
    // File watcher
    fileConfig := collectors.DefaultFileWatchConfig()
    fileWatcher := collectors.NewFileWatcher(fileConfig)
    
    // Process watcher
    procConfig := collectors.DefaultProcessWatchConfig()
    procWatcher := collectors.NewProcessWatcher(procConfig)
    
    // Network watcher
    netConfig := collectors.DefaultNetworkWatchConfig()
    netConfig.EventCallback = func(e events.NetworkEvent) {
        if err := client.SendEvent(e); err != nil {
            log.Printf("Failed to send network event: %v", err)
        }
    }
    netConfig.AlertCallback = func(e events.NetworkEvent) {
        log.Printf("[ALERT] Network: %s → %s:%d (process: %s)",
            e.SrcIP, e.DstIP, e.DstPort, e.Process)
    }
    netWatcher := collectors.NewNetworkWatcher(netConfig)
    
    // Start all collectors
    fileWatcher.Start()
    procWatcher.Start()
    netWatcher.Start()
    
    log.Println("GoShield agent running. Ctrl+C to stop.")
    
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
    <-sigCh
    
    log.Println("Shutting down...")
    fileWatcher.Stop()
    procWatcher.Stop()
    netWatcher.Stop()
}
```

---

## Testing the Network Monitor

```bash
# Build the agent
cd goshield/agent
go build -o goshield-agent .

# Run as root (needed for network connection info)
sudo ./goshield-agent

# In another terminal, generate test connections
nc -zv 192.168.1.1 4444    # Should trigger NET-001 (suspicious port)
for i in $(seq 1 25); do   # Should trigger beaconing detection
    nc -zv 8.8.8.8 443 &
done

# Check what the agent reports
```

---

## Threat Scenarios Detected

| Threat | Indicator | Rule |
|--------|-----------|------|
| Reverse shell | Outbound to port 4444/1337/etc | NET-001 |
| C2 beaconing | High-frequency periodic connections | NET-002 |
| DNS tunneling | Large DNS queries, high frequency | NET-002 |
| IRC botnet | Port 6667/6668 connections | NET-003 |
| Webshell → C2 | nginx/apache making outbound connections | NET-004 |
| Data exfiltration | Massive outbound transfer volume | NET-005 |
| Lateral movement | Internal connections to unusual ports | Custom |

---

## Summary

The GoShield network collector:
1. Polls connection state every 5 seconds via `gopsutil`
2. Diffs against previous state to find new/closed connections
3. Resolves IPs to domain names with caching
4. Identifies the owning process for each connection
5. Runs detection rules for suspicious patterns
6. Reports events to the central server

Combined with file and process monitoring (Chapters 60-61), this gives complete visibility into the endpoint: **what files changed, what processes ran, and what network connections were made** — the three pillars of EDR visibility.

---

## Exercises

1. Add a `geo` field to network events using a local GeoIP database (MaxMind's free `GeoLite2-Country` database works offline)
2. Implement connection duration tracking — alert if a connection stays open for > 1 hour to an external IP
3. Add detection for connections to newly registered domains (domains registered < 30 days ago — common for C2 infrastructure)
4. Build a "connection baseline" that learns normal behavior for 24 hours, then alerts on deviations
