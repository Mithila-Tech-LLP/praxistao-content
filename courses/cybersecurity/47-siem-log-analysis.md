# Chapter 47: SIEM and Log Analysis — Finding Attacks in the Noise

*A SIEM receives millions of log events per day. 99.99% are benign. Your job is to find the 0.01% that matters — before the attacker achieves their objective. This is detection engineering.*

---

## The Scale Problem

A mid-sized company (1000 employees) generates:
- 500,000 Windows Security events/day
- 100,000 authentication events/day
- 50,000 network flow records/day
- 10,000 web proxy logs/day
- 5,000 firewall events/day

**= ~665,000 events per day**

No human can read these. You need:
1. **Collection:** Get all logs to one place
2. **Normalization:** Common format across different sources
3. **Detection rules:** Alert on patterns of interest
4. **Prioritization:** Not all alerts are equal

---

## SIEM Architecture

```
Log Sources                  Collection          SIEM Platform
─────────────               ────────────        ─────────────
Windows Event Logs ──────→  Winlogbeat ──────→
Linux Syslog ────────────→  Filebeat ────────→  Elasticsearch
Network Flows ───────────→  Logstash ────────→  + Detection Rules
Web Proxy ───────────────→  Fluentd ─────────→  + Kibana Dashboard
Firewall ────────────────→  Syslog ──────────→
EDR (GoShield) ──────────→  REST API ────────→
                                                      ↓
                                             Alerts → SOC Analyst
```

---

## Log Sources and What They Tell You

### Windows Event Logs

Windows logs everything to the Event Viewer. Critical event IDs:

```
Security Log:
  4624 — Successful logon (what account, from where, logon type)
  4625 — Failed logon (brute force detection)
  4648 — Logon with explicit credentials (pass-the-hash, runas)
  4672 — Special privileges assigned (admin logon)
  4688 — Process creation (enabled in policy)
  4698 — Scheduled task created (persistence)
  4720 — User account created
  4728 — Member added to security-enabled global group
  4776 — NTLM authentication attempt
  4768/4769 — Kerberos ticket requests (golden ticket detection)

System Log:
  7045 — Service installed (persistence mechanism)
  
PowerShell:
  4103 — Module logging
  4104 — Script block logging ← most important for detecting malicious PS
```

**Logon types:**
- Type 2: Interactive (keyboard)
- Type 3: Network (SMB, RDP headless)
- Type 10: Remote Interactive (RDP)
- Type 5: Service

Type 3 logons at 3AM from a new IP = investigate.

### Linux Logs

```
/var/log/auth.log    — SSH, sudo, PAM authentication
/var/log/syslog      — General system events
/var/log/kern.log    — Kernel messages (module loading, etc.)
/var/log/cron.log    — Cron job execution
/var/log/audit/audit.log — Linux Audit daemon (most granular)

# Audit daemon captures:
# - File opens/reads/writes
# - System calls
# - Network connections
# - Process execution
# - User/group changes
```

### Web Server Logs

```apache
# Nginx/Apache access log format:
192.168.1.100 - - [25/Jun/2024:12:34:56 +0000] "GET /admin HTTP/1.1" 403 512 "-" "Mozilla/5.0"

IP  -  user  [timestamp]  "method path protocol"  status  bytes  referrer  user-agent
```

**Patterns to detect:**
- HTTP 403/401 storms → brute force or directory scanning
- POST to `/xmlrpc.php` → WordPress attack
- Long query strings with `'`, `--`, `union`, `select` → SQL injection
- `.php` file requests on non-PHP sites → webshell hunting
- Unusual user agents → automated tools

### Network Flow Logs (NetFlow/sFlow)

```
src_ip:src_port → dst_ip:dst_port | protocol | bytes | packets | duration
192.168.1.100:54321 → 10.0.0.50:445 | TCP | 1234 | 10 | 0.5s
```

**Key insights:**
- Volume anomalies (100GB outbound at 2AM = exfiltration)
- Port scan patterns (one source, many destinations, many ports)
- Lateral movement (internal source → internal destinations, SMB/RDP)
- Beaconing (regular connections to external IP, fixed interval)

---

## Detection Rules — Writing Quality Alerts

### Anatomy of a Good Detection Rule

A good rule has:
1. **High precision:** Few false positives (analysts ignore noisy alerts)
2. **High recall:** Catches real attacks (not easily evaded)
3. **Context:** Enough info to act on immediately
4. **Documented rationale:** Why does this indicate an attack?

### Sigma Rules

Sigma is the de facto standard for writing portable detection rules.

```yaml
# Sigma rule for detecting PsExec lateral movement
title: PsExec Lateral Movement
id: 42e69bfd-9a73-4b29-9e5e-f7412c1f6c07
status: stable
description: Detects PsExec usage for lateral movement via Windows service creation
author: Detection Engineering Team
date: 2024/01/15
tags:
  - attack.lateral_movement
  - attack.t1021.002
logsource:
  product: windows
  service: system
detection:
  selection:
    EventID: 7045    # Service Install
    ServiceName|contains:
      - 'PSEXESVC'
      - 'psexec'
    ServiceFileName|contains: '\ADMIN$\'
  condition: selection
falsepositives:
  - Legitimate admin use of PsExec
level: high
```

### Splunk Detection Query

```spl
# Detect brute force: 10+ failed logins in 5 minutes from same IP
index=security EventCode=4625
| stats count by src_ip, _time span=5m
| where count > 10
| table _time, src_ip, count

# Detect successful login after failures (possible successful brute force)
index=security (EventCode=4624 OR EventCode=4625)
| stats sum(eval(EventCode=4625)) as failures, 
        sum(eval(EventCode=4624)) as successes 
        by src_ip user span=1h
| where failures > 5 AND successes > 0
| sort -failures

# Detect PowerShell with encoded commands (common malware technique)
index=sysmon EventCode=1
| search ParentCommandLine="*powershell*" AND CommandLine="*-enc*"
| table _time, Computer, ParentImage, CommandLine

# Detect process injection (suspicious process creating remote thread)
index=sysmon EventCode=8
| where SourceImage != TargetImage
| where NOT (SourceImage IN (
    "C:\\Windows\\System32\\svchost.exe",
    "C:\\Windows\\System32\\lsass.exe"
))
| table _time, Computer, SourceImage, TargetImage
```

### Elastic KQL / Kibana Rules

```json
{
  "name": "Windows - Lateral Movement via SMB",
  "query": "event.action: 'logged-in' AND winlog.event_id: 4624 AND winlog.logon.type: 3 AND NOT source.ip: '127.0.0.1'",
  "threshold": {
    "field": "source.ip",
    "value": 1
  },
  "filters": [
    {"range": {"@timestamp": {"gte": "now-5m"}}}
  ]
}
```

---

## Threat Hunting

Detection rules catch known-bad patterns. Threat hunting proactively looks for unknown-bad.

### Hunting Process

```
1. Hypothesis: "What if an attacker is using living-off-the-land techniques?"
2. Hunt: Look for PowerShell or wscript spawning from Office processes
3. Pivot: Find anomalies, investigate
4. Improve: Turn findings into detection rules
```

### Common Hunt Hypotheses

```
"Web server is compromised"
→ Hunt: Look for outbound connections from web server processes
→ Query: process:nginx AND network.direction:outbound AND NOT dst_port:80,443

"Credential theft happened"
→ Hunt: Look for LSASS memory access (process dump)
→ Sigma: EventCode=10 AND TargetImage contains lsass.exe AND GrantedAccess:0x1fffff

"Lateral movement is occurring"
→ Hunt: Internal source → internal destination, SMB, new connections
→ Query: src.subnet:10.0.0.0/8 AND dst.subnet:10.0.0.0/8 AND dst.port:445

"Beaconing C2 traffic"
→ Hunt: Regular intervals, same destination, small packets
→ Query: Connections to same external IP, standard deviation of interval < 30s
```

---

## Go: Log Parser and Alert Generator

```go
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"
    "time"
)

// NginxLog parses a standard nginx access log entry
type NginxLog struct {
    IP        string
    Timestamp time.Time
    Method    string
    Path      string
    Status    int
    Bytes     int
    UserAgent string
}

var nginxPattern = regexp.MustCompile(
    `^(\S+) \S+ \S+ \[([^\]]+)\] "(\S+) (\S+) [^"]+" (\d+) (\d+) "[^"]*" "([^"]*)"`,
)

// Detects SQL injection in paths/params
var sqliPatterns = []*regexp.Regexp{
    regexp.MustCompile(`(?i)(\bSELECT\b|\bUNION\b|\bINSERT\b|\bDROP\b|\bDELETE\b)`),
    regexp.MustCompile(`('|--|;|/\*|\*/|@@|char\(|exec\(|xp_)`),
    regexp.MustCompile(`(?i)(OR\s+1=1|AND\s+1=1|OR\s+'1'='1')`),
}

// Suspicious paths
var suspiciousPaths = []string{
    "/wp-admin", "/wp-login", "/phpmyadmin", "/admin",
    "/.env", "/.git", "/config", "/backup",
    "/xmlrpc.php", "/shell.php", "/webshell",
}

type Alert struct {
    Time     time.Time `json:"time"`
    Type     string    `json:"type"`
    Severity string    `json:"severity"`
    IP       string    `json:"ip"`
    Detail   string    `json:"detail"`
    Path     string    `json:"path"`
}

func analyzeLog(logFile string) {
    f, err := os.Open(logFile)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Cannot open log: %v\n", err)
        return
    }
    defer f.Close()

    // Track per-IP stats
    ipErrors := make(map[string]int)
    ipRequests := make(map[string]int)
    
    var alerts []Alert
    
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        m := nginxPattern.FindStringSubmatch(line)
        if m == nil {
            continue
        }
        
        ip := m[1]
        path := m[4]
        var status, bytes int
        fmt.Sscanf(m[5], "%d", &status)
        fmt.Sscanf(m[6], "%d", &bytes)
        
        ipRequests[ip]++
        if status >= 400 {
            ipErrors[ip]++
        }
        
        // Check 1: SQL injection in path
        for _, pattern := range sqliPatterns {
            if pattern.MatchString(path) {
                alerts = append(alerts, Alert{
                    Time:     time.Now(),
                    Type:     "SQL_INJECTION",
                    Severity: "high",
                    IP:       ip,
                    Detail:   "SQLi pattern in request path",
                    Path:     path,
                })
                break
            }
        }
        
        // Check 2: Suspicious paths
        for _, suspPath := range suspiciousPaths {
            if strings.Contains(strings.ToLower(path), suspPath) {
                alerts = append(alerts, Alert{
                    Time:     time.Now(),
                    Type:     "SUSPICIOUS_PATH",
                    Severity: "medium",
                    IP:       ip,
                    Detail:   "Access to sensitive path",
                    Path:     path,
                })
                break
            }
        }
    }
    
    // Check 3: High error rate per IP (scanning)
    for ip, errors := range ipErrors {
        total := ipRequests[ip]
        if total > 50 && errors > 30 {
            rate := float64(errors) / float64(total) * 100
            alerts = append(alerts, Alert{
                Time:     time.Now(),
                Type:     "DIRECTORY_SCAN",
                Severity: "medium",
                IP:       ip,
                Detail:   fmt.Sprintf("%.0f%% error rate (%d errors / %d requests)", rate, errors, total),
            })
        }
    }
    
    // Sort alerts by severity
    sort.Slice(alerts, func(i, j int) bool {
        sev := map[string]int{"critical": 4, "high": 3, "medium": 2, "low": 1}
        return sev[alerts[i].Severity] > sev[alerts[j].Severity]
    })
    
    // Output
    fmt.Printf("Found %d alerts\n\n", len(alerts))
    for _, a := range alerts {
        enc, _ := json.Marshal(a)
        fmt.Println(string(enc))
    }
}

func main() {
    if len(os.Args) < 2 {
        fmt.Println("Usage: loganalyzer <logfile>")
        return
    }
    analyzeLog(os.Args[1])
}
```

---

## Dashboard Metrics for a SOC

**What to show on a SOC dashboard:**

```
┌─────────────────────────────────────────────────────┐
│ SECURITY OPERATIONS CENTER                          │
│                                                     │
│ Active Alerts: 12 CRITICAL  34 HIGH  156 MEDIUM    │
│                                                     │
│ Top Attack Sources (24h):           Alert Timeline: │
│  185.220.101.5    342 events         [sparkline]   │
│  192.242.116.160  289 events                       │
│  45.33.32.156     201 events        MTTR: 4.2 hrs  │
│                                                     │
│ Attack Types:                       Systems Status: │
│  SQLi: 45%                           EDR: ✓ 847    │
│  Brute Force: 30%                    Firewall: ✓   │
│  Dir Scan: 15%                       SIEM: ✓       │
│  Other: 10%                          DNS: ✓        │
└─────────────────────────────────────────────────────┘
```

---

## Summary

| Component | Purpose | Key tools |
|-----------|---------|----------|
| Log collection | Get logs to SIEM | Beats, Fluentd, Syslog |
| SIEM platform | Store, search, alert | Splunk, Elastic, Wazuh |
| Detection rules | Alert on patterns | Sigma, SPL, KQL |
| Threat hunting | Find unknown threats | Analyst-driven queries |
| Dashboards | Operational visibility | Kibana, Grafana |

---

## Exercises

1. Set up Elastic Stack (Elasticsearch + Kibana). Ship logs from your Linux system using Filebeat. Write a KQL query to find failed SSH logins.
2. Write 3 Sigma rules: brute force, web directory scan, and command injection in web logs.
3. Run the Go log analyzer against a Nginx access log. How many alerts does it find?
4. Hunt for beaconing in network logs: write a script that finds IPs connected to at regular intervals
5. Build a Kibana dashboard showing top attacking IPs and alert categories
