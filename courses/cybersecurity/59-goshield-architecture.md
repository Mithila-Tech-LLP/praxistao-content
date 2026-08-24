# Chapter 59: GoShield Architecture — Designing the Enterprise EDR

*Good architecture is the difference between a toy project and a production system. This chapter designs GoShield as if it were a real enterprise product — decisions that scale to hundreds of thousands of endpoints.*

---

## Requirements

Before designing, define what GoShield must do:

**Functional requirements:**
1. Collect file, process, and network events from Linux and Windows endpoints
2. Ship events to a central server with minimal performance impact on the endpoint
3. Apply detection rules and generate alerts in near-real-time
4. Allow security analysts to search and hunt through historical events
5. Automatically respond to threats (kill process, quarantine)
6. Expose an API for integration with SIEMs and ticketing systems

**Non-functional requirements:**
- **Throughput:** Handle 10,000 events/second per server instance
- **Latency:** Alert within 10 seconds of suspicious event
- **Agent overhead:** Less than 2% CPU, less than 50MB RAM on endpoint
- **Storage:** Retain 90 days of events per endpoint
- **Reliability:** Agent must survive server unavailability, buffer events

---

## System Design

```
┌─────────────────────────────────────────────────────────────────┐
│                     ENTERPRISE NETWORK                          │
│                                                                 │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌──────────┐       │
│  │ Laptop 1 │  │ Server A │  │ Server B │  │ Desktop  │       │
│  │ [Agent]  │  │ [Agent]  │  │ [Agent]  │  │ [Agent]  │       │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────┬─────┘       │
│       │              │              │              │             │
│       └──────────────┴──────┬───────┴──────────────┘           │
│                              │ TLS/gRPC                         │
│                    ┌─────────▼──────────┐                       │
│                    │  GoShield Server   │                       │
│                    │                   │                       │
│                    │  ┌─────────────┐  │                       │
│                    │  │ API Server  │  │◀─── Analyst Browser   │
│                    │  ├─────────────┤  │                       │
│                    │  │  Ingestion  │  │                       │
│                    │  ├─────────────┤  │                       │
│                    │  │ Detection   │  │                       │
│                    │  │  Engine     │  │                       │
│                    │  ├─────────────┤  │                       │
│                    │  │  Alerting   │─────▶ Slack/Email/SIEM  │
│                    │  ├─────────────┤  │                       │
│                    │  │  Storage    │  │                       │
│                    │  │(SQLite/PG)  │  │                       │
│                    │  └─────────────┘  │                       │
│                    └───────────────────┘                       │
└─────────────────────────────────────────────────────────────────┘
```

---

## Repository Structure

```
goshield/
├── go.mod
├── go.sum
├── README.md
│
├── cmd/
│   ├── agent/
│   │   └── main.go          # Agent entry point
│   └── server/
│       └── main.go          # Server entry point
│
├── pkg/
│   ├── events/
│   │   └── types.go         # Event type definitions
│   │
│   ├── agent/
│   │   ├── agent.go         # Agent core
│   │   ├── filewatcher.go   # File system monitoring
│   │   ├── procwatcher.go   # Process monitoring
│   │   ├── netwatcher.go    # Network monitoring
│   │   └── sender.go        # Event shipping to server
│   │
│   ├── server/
│   │   ├── server.go        # HTTP server setup
│   │   ├── api.go           # REST API handlers
│   │   ├── ingestion.go     # Event ingestion
│   │   └── storage.go       # Database layer
│   │
│   ├── detector/
│   │   ├── engine.go        # Detection engine
│   │   ├── rule.go          # Rule type and parser
│   │   └── evaluator.go     # Rule evaluation logic
│   │
│   └── alerting/
│       ├── alerter.go       # Alert dispatch
│       ├── webhook.go       # Webhook alerts
│       └── email.go         # Email alerts
│
├── rules/
│   ├── process/
│   │   ├── psh_encoded.yaml
│   │   └── reverse_shell.yaml
│   ├── file/
│   │   ├── ransomware.yaml
│   │   └── sensitive_files.yaml
│   └── network/
│       ├── beaconing.yaml
│       └── dns_tunneling.yaml
│
├── web/
│   ├── index.html           # Dashboard
│   ├── app.js
│   └── style.css
│
└── Dockerfile
```

---

## Core Data Structures

```go
// pkg/events/types.go
package events

import (
    "crypto/rand"
    "encoding/hex"
    "time"
)

// NewID generates a random event ID
func NewID() string {
    b := make([]byte, 8)
    rand.Read(b)
    return hex.EncodeToString(b)
}

// Severity levels
type Severity string

const (
    SeverityInfo     Severity = "info"
    SeverityLow      Severity = "low"
    SeverityMedium   Severity = "medium"
    SeverityHigh     Severity = "high"
    SeverityCritical Severity = "critical"
)

// EventType categorizes events
type EventType string

const (
    EventTypeProcess EventType = "process"
    EventTypeFile    EventType = "file"
    EventTypeNetwork EventType = "network"
    EventTypeAlert   EventType = "alert"
)

// Base is embedded in all events
type Base struct {
    ID        string    `json:"id" db:"id"`
    AgentID   string    `json:"agent_id" db:"agent_id"`
    Hostname  string    `json:"hostname" db:"hostname"`
    OS        string    `json:"os" db:"os"`
    Timestamp time.Time `json:"timestamp" db:"timestamp"`
    Type      EventType `json:"type" db:"type"`
}

// ProcessEvent captures process lifecycle events
type ProcessEvent struct {
    Base
    Action      string `json:"action"`       // "create" | "terminate"
    PID         int    `json:"pid"`
    PPID        int    `json:"ppid"`
    Name        string `json:"name"`
    CommandLine string `json:"command_line"`
    Username    string `json:"username"`
    ExePath     string `json:"exe_path"`
    SHA256      string `json:"sha256"`
    IsElevated  bool   `json:"is_elevated"`  // running as root/admin?
}

// FileEvent captures file system events
type FileEvent struct {
    Base
    Action    string `json:"action"`     // "create" | "write" | "delete" | "rename"
    Path      string `json:"path"`
    NewPath   string `json:"new_path"`   // for rename
    SHA256    string `json:"sha256"`
    Size      int64  `json:"size"`
    PID       int    `json:"pid"`
    Process   string `json:"process"`
    Extension string `json:"extension"`
}

// NetworkEvent captures connection events  
type NetworkEvent struct {
    Base
    Action     string `json:"action"`      // "connect" | "listen" | "close"
    Protocol   string `json:"protocol"`    // "tcp" | "udp"
    SrcIP      string `json:"src_ip"`
    SrcPort    int    `json:"src_port"`
    DstIP      string `json:"dst_ip"`
    DstPort    int    `json:"dst_port"`
    PID        int    `json:"pid"`
    Process    string `json:"process"`
    Domain     string `json:"domain"`      // resolved DNS name
    Country    string `json:"country"`     // GeoIP
    IsPrivate  bool   `json:"is_private"`  // private IP range?
}

// Alert is generated by the detection engine
type Alert struct {
    Base
    RuleID      string    `json:"rule_id"`
    RuleName    string    `json:"rule_name"`
    Severity    Severity  `json:"severity"`
    Description string    `json:"description"`
    EventID     string    `json:"event_id"`    // triggering event
    EventType   EventType `json:"event_type"`
    MITRE       string    `json:"mitre"`       // ATT&CK technique
    Resolved    bool      `json:"resolved"`
    ResolvedAt  time.Time `json:"resolved_at,omitempty"`
    Notes       string    `json:"notes"`
}

// AgentStatus tracks agent health
type AgentStatus struct {
    AgentID     string    `json:"agent_id"`
    Hostname    string    `json:"hostname"`
    OS          string    `json:"os"`
    Version     string    `json:"version"`
    LastSeen    time.Time `json:"last_seen"`
    Online      bool      `json:"online"`
    EventsTotal int64     `json:"events_total"`
}
```

---

## Configuration Files

```yaml
# agent-config.yaml
agent:
  id: ""                     # auto-generated if empty
  server: "https://goshield.company.com:8443"
  api_key: "YOUR_API_KEY_HERE"
  
monitoring:
  file_watch:
    enabled: true
    paths:
      - /etc
      - /bin
      - /usr/bin
      - /home
      - /tmp
      - /var/tmp
    exclude_patterns:
      - "*.log"
      - "*.pid"
    
  process_watch:
    enabled: true
    
  network_watch:
    enabled: true
    ignore_private_ips: false

transport:
  batch_size: 100            # Send events in batches of 100
  flush_interval: 5s         # Or every 5 seconds, whichever comes first
  retry_interval: 30s        # Retry if server unavailable
  max_buffer: 10000          # Buffer up to 10000 events offline
```

```yaml
# server-config.yaml
server:
  listen: ":8443"
  tls:
    cert: "/etc/goshield/server.crt"
    key: "/etc/goshield/server.key"
  api_key: "CHANGE_THIS_SECRET"

storage:
  driver: "sqlite"           # or "postgres"
  path: "/var/lib/goshield/events.db"
  retention_days: 90

detection:
  rules_dir: "/etc/goshield/rules"
  
alerting:
  webhooks:
    - name: "Slack"
      url: "https://hooks.slack.com/services/..."
      on_severity: ["high", "critical"]
  email:
    smtp_host: "smtp.gmail.com:587"
    from: "goshield@company.com"
    to: ["security@company.com"]
    on_severity: ["critical"]
```

---

## Module Setup

```bash
# Initialize the GoShield project
mkdir goshield && cd goshield
go mod init github.com/yourname/goshield

# Install dependencies
go get github.com/fsnotify/fsnotify    # File system watching
go get github.com/shirou/gopsutil/v3  # Process and network info
go get github.com/mattn/go-sqlite3    # SQLite storage
go get gopkg.in/yaml.v3               # YAML config/rules
go get github.com/gorilla/mux         # HTTP routing
go get github.com/rs/zerolog          # Fast structured logging
go get github.com/spf13/viper         # Configuration
```

---

## Agent Startup Sequence

```go
// cmd/agent/main.go
package main

import (
    "log"
    "os"
    "os/signal"
    "syscall"
    
    "github.com/yourname/goshield/pkg/agent"
)

func main() {
    configPath := "/etc/goshield/agent-config.yaml"
    if len(os.Args) > 1 {
        configPath = os.Args[1]
    }
    
    a, err := agent.New(configPath)
    if err != nil {
        log.Fatalf("Failed to create agent: %v", err)
    }
    
    // Graceful shutdown on SIGTERM/SIGINT
    sigCh := make(chan os.Signal, 1)
    signal.Notify(sigCh, syscall.SIGTERM, syscall.SIGINT)
    
    // Start the agent (non-blocking)
    if err := a.Start(); err != nil {
        log.Fatalf("Failed to start agent: %v", err)
    }
    
    log.Printf("GoShield Agent started (ID: %s)", a.ID())
    
    // Wait for shutdown signal
    sig := <-sigCh
    log.Printf("Received signal %v, shutting down...", sig)
    
    a.Stop()
    log.Println("Agent stopped cleanly")
}
```

---

## Summary: What We're Building Next

| Chapter | What gets built |
|---------|----------------|
| 60 | Agent: File Integrity Monitoring (fsnotify) |
| 61 | Agent: Process Monitoring (gopsutil) |
| 62 | Agent: Network Connection Monitoring |
| 63 | Agent: System Call Monitoring (Linux-specific) |
| 64 | Server: Event ingestion and SQLite storage |
| 65 | Detection Engine: YAML rules + evaluation |
| 66 | Alerting: Webhook and email |
| 67 | REST API + Web Dashboard |
| 68 | Packaging: Docker + deployment |
| 69 | Advanced: Behavioral analysis |

The full GoShield codebase will be a working enterprise EDR system you can deploy on real servers.

---

## Exercises

1. Create the directory structure and go.mod file for GoShield
2. Install all listed dependencies (`go get`)
3. Create the events/types.go file with the type definitions above
4. Write a simple test: create a `ProcessEvent`, marshal it to JSON, print it
5. Design the database schema: what tables do you need to store events, alerts, and agent status?
