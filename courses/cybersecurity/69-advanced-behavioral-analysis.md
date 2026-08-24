# Chapter 69: GoShield — Advanced Behavioral Analysis

*Signature-based detection catches known malware. Behavioral analysis catches unknown threats by recognizing attack patterns regardless of what tools are used.*

---

## Behavioral vs Signature Detection

```
Signature detection:
"This file has MD5 hash abc123 = WannaCry ransomware"
- Fast and certain when it matches
- Completely misses new/modified malware
- Attackers change one byte → new hash → bypass

Behavioral detection:
"A process is encrypting 1000+ files in 10 seconds"
- Catches new ransomware variants
- Catches fileless malware (no file to hash)
- Requires tuning to avoid false positives
```

---

## Behavioral Patterns to Detect

```
Process Behavior:
├── Web server spawning shell      → webshell/command injection
├── Office spawning PowerShell     → macro malware
├── Explorer spawning cmd.exe      → user double-clicked malware
├── Unusual parent-child chains    → process hollowing
└── Process reading LSASS memory   → credential theft

File System Behavior:
├── Mass file encryption           → ransomware
├── Shadow copy deletion           → ransomware (preparing)
├── SUID binary creation           → Linux privilege escalation
└── New cron jobs / services       → persistence

Network Behavior:
├── Regular beaconing              → C2 communication
├── Large outbound data            → exfiltration
├── Internal scanning              → lateral movement preparation
└── Connection to TOR exit nodes   → covert communication

Memory Behavior:
├── Anonymous RWX memory           → shellcode injection
├── Heap shellcode execution       → exploit
└── Unusual DLL loads              → DLL injection
```

---

## Go: Behavioral Analysis Engine

```go
package behavioral

import (
    "fmt"
    "sync"
    "time"
)

// Event types
type EventType string

const (
    EventProcessSpawn  EventType = "process_spawn"
    EventFileEncrypt   EventType = "file_encrypt"
    EventNetConnect    EventType = "net_connect"
    EventRegistryWrite EventType = "registry_write"
    EventShadowDelete  EventType = "shadow_delete"
)

type Event struct {
    Time       time.Time
    Type       EventType
    PID        int
    PPID       int
    ProcessName string
    ParentName  string
    Data       map[string]string
}

type BehaviorAlert struct {
    Severity    string
    Technique   string  // ATT&CK technique
    Description string
    Events      []Event
    Score       int
}

// Sliding window counter for frequency-based detection
type WindowCounter struct {
    mu      sync.Mutex
    events  map[string][]time.Time
    window  time.Duration
}

func NewWindowCounter(window time.Duration) *WindowCounter {
    return &WindowCounter{
        events: make(map[string][]time.Time),
        window: window,
    }
}

func (wc *WindowCounter) Add(key string) int {
    wc.mu.Lock()
    defer wc.mu.Unlock()
    
    now := time.Now()
    cutoff := now.Add(-wc.window)
    
    // Add new event
    wc.events[key] = append(wc.events[key], now)
    
    // Prune old events
    var recent []time.Time
    for _, t := range wc.events[key] {
        if t.After(cutoff) {
            recent = append(recent, t)
        }
    }
    wc.events[key] = recent
    
    return len(recent)
}

// Ransomware detector: file encryption rate
type RansomwareDetector struct {
    encryptCounter *WindowCounter
    threshold      int
}

func NewRansomwareDetector() *RansomwareDetector {
    return &RansomwareDetector{
        encryptCounter: NewWindowCounter(30 * time.Second),
        threshold:      50,  // 50 files encrypted in 30s = ransomware
    }
}

func (r *RansomwareDetector) Analyze(event Event) *BehaviorAlert {
    if event.Type != EventFileEncrypt {
        return nil
    }
    
    key := fmt.Sprintf("%d", event.PID)
    count := r.encryptCounter.Add(key)
    
    if count >= r.threshold {
        return &BehaviorAlert{
            Severity:    "CRITICAL",
            Technique:   "T1486 Data Encrypted for Impact",
            Description: fmt.Sprintf("Process %s (PID %d) encrypted %d files in 30 seconds — likely ransomware",
                event.ProcessName, event.PID, count),
            Events: []Event{event},
            Score:  100,
        }
    }
    return nil
}

// Webshell detector: web server spawning shell
type WebshellDetector struct{}

func (w *WebshellDetector) Analyze(event Event) *BehaviorAlert {
    if event.Type != EventProcessSpawn {
        return nil
    }
    
    webServers := map[string]bool{
        "apache2": true, "httpd": true, "nginx": true,
        "php-fpm": true, "php": true, "tomcat": true,
        "java": true, "uwsgi": true, "gunicorn": true,
    }
    
    shells := map[string]bool{
        "bash": true, "sh": true, "dash": true,
        "python": true, "python3": true, "perl": true,
        "ruby": true, "nc": true, "ncat": true,
    }
    
    if webServers[event.ParentName] && shells[event.ProcessName] {
        return &BehaviorAlert{
            Severity:    "HIGH",
            Technique:   "T1505.003 Web Shell",
            Description: fmt.Sprintf("Web server '%s' spawned shell '%s' (PID %d) — possible webshell execution",
                event.ParentName, event.ProcessName, event.PID),
            Events: []Event{event},
            Score:  80,
        }
    }
    return nil
}

// C2 beaconing detector: regular periodic connections
type BeaconingDetector struct {
    connectionTimes map[string][]time.Time  // dstIP → connection times
    mu              sync.Mutex
}

func NewBeaconingDetector() *BeaconingDetector {
    return &BeaconingDetector{
        connectionTimes: make(map[string][]time.Time),
    }
}

func (b *BeaconingDetector) Analyze(event Event) *BehaviorAlert {
    if event.Type != EventNetConnect {
        return nil
    }
    
    dstIP := event.Data["dst_ip"]
    if dstIP == "" {
        return nil
    }
    
    b.mu.Lock()
    defer b.mu.Unlock()
    
    now := time.Now()
    b.connectionTimes[dstIP] = append(b.connectionTimes[dstIP], now)
    
    times := b.connectionTimes[dstIP]
    
    // Need at least 5 connections to detect beaconing
    if len(times) < 5 {
        return nil
    }
    
    // Check if intervals are consistent (beacons are regular)
    recentTimes := times[len(times)-5:]
    intervals := make([]float64, 4)
    for i := 1; i < len(recentTimes); i++ {
        intervals[i-1] = recentTimes[i].Sub(recentTimes[i-1]).Seconds()
    }
    
    // Calculate standard deviation of intervals
    mean := 0.0
    for _, iv := range intervals {
        mean += iv
    }
    mean /= float64(len(intervals))
    
    variance := 0.0
    for _, iv := range intervals {
        diff := iv - mean
        variance += diff * diff
    }
    variance /= float64(len(intervals))
    
    stdDev := variance
    if stdDev < 0 {
        stdDev = -stdDev
    }
    // Simplified sqrt
    for i := 0; i < 10; i++ {
        if stdDev <= 0 {
            break
        }
        stdDev = (stdDev + variance/stdDev) / 2
    }
    
    // If standard deviation is <10% of mean → very regular → beaconing
    if mean > 30 && stdDev < mean*0.1 {
        return &BehaviorAlert{
            Severity:    "HIGH",
            Technique:   "T1071 Application Layer Protocol (C2)",
            Description: fmt.Sprintf("Regular beaconing to %s: interval=%.0fs (stddev=%.1fs) — possible C2",
                dstIP, mean, stdDev),
            Events: []Event{event},
            Score:  75,
        }
    }
    
    return nil
}

// Main behavior analysis engine
type BehaviorEngine struct {
    detectors []interface{ Analyze(Event) *BehaviorAlert }
    alerts    chan *BehaviorAlert
}

func NewBehaviorEngine() *BehaviorEngine {
    engine := &BehaviorEngine{
        alerts: make(chan *BehaviorAlert, 100),
    }
    engine.detectors = []interface{ Analyze(Event) *BehaviorAlert }{
        NewRansomwareDetector(),
        &WebshellDetector{},
        NewBeaconingDetector(),
    }
    return engine
}

func (e *BehaviorEngine) Process(event Event) {
    for _, detector := range e.detectors {
        if alert := detector.Analyze(event); alert != nil {
            select {
            case e.alerts <- alert:
            default:
                fmt.Println("[WARN] Alert channel full, dropping alert")
            }
        }
    }
}

func (e *BehaviorEngine) Alerts() <-chan *BehaviorAlert {
    return e.alerts
}
```

---

## MITRE ATT&CK Behavioral Coverage

```
T1055  Process Injection      → anonymous RWX memory, ptrace, CreateRemoteThread
T1486  Data Encrypted         → mass file encryption (ransomware)
T1505  Server Software Comp.  → web server spawning shell (webshell)
T1071  C2 Communication       → regular beaconing pattern
T1547  Boot Autostart         → new cron/service/registry run key
T1003  OS Credential Dumping  → LSASS memory access, /etc/shadow read
T1021  Remote Services        → unusual lateral movement via SMB/SSH
```

---

## Summary

| Behavioral Detector | Pattern | Technique |
|--------------------|---------|-----------|
| Ransomware | >50 file encryptions/30s | T1486 |
| Webshell | Web server → shell spawn | T1505 |
| C2 beaconing | Regular interval connections | T1071 |
| Credential theft | LSASS/shadow access | T1003 |
| Persistence | New cron/service | T1547 |

---

## Exercises

1. Integrate the behavior engine with GoShield's file monitor — trigger the ransomware detector by creating and "encrypting" 60 files quickly
2. Add a lateral movement detector: SSH connection from a host that doesn't normally SSH
3. Test the webshell detector by simulating `apache2 → bash` process spawn
4. Research UEBA (User and Entity Behavior Analytics) — how does it extend basic behavioral detection?
