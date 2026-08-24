# Chapter 52: SOC Operations — The Security Command Center

*A Security Operations Center (SOC) is the team and technology responsible for monitoring, detecting, and responding to security incidents 24/7. Understanding how a SOC works is essential whether you're defending one or red-teaming against one.*

---

## SOC Structure

```
SOC Team Hierarchy:
├── Tier 1 — Alert Triage Analysts
│   ├── Monitor security tools 24/7
│   ├── First responders to alerts
│   ├── Investigate and close false positives
│   └── Escalate real incidents to Tier 2
│
├── Tier 2 — Incident Responders
│   ├── Deep dive on escalated incidents
│   ├── Threat hunting
│   ├── Malware analysis
│   └── Remediation guidance
│
├── Tier 3 — Threat Hunters / Senior Analysts
│   ├── Proactive hunting for hidden threats
│   ├── Detection engineering (write new rules)
│   ├── Red team collaboration
│   └── Tooling development
│
└── SOC Manager / CISO
    ├── Metrics and reporting
    ├── Budget and staffing
    └── Strategic direction
```

---

## SIEM (Security Information and Event Management)

The SIEM is the brain of the SOC:

```
Log sources              SIEM                  Analysts
──────────             ──────────           ────────────────
Firewalls           →  Collect & normalize  →  Alert queue
Web proxies         →  Correlate events     →  Investigation
Windows event logs  →  Apply rules          →  Incident creation
Linux syslog        →  Machine learning     →  Escalation
EDR alerts          →  User behavior        →
Cloud logs          →  Dashboards           →
```

### Popular SIEMs

```
Open Source:
- Elastic SIEM (ELK: Elasticsearch + Logstash + Kibana)
- Wazuh (fork of OSSEC, modern)
- Graylog

Commercial:
- Splunk (most popular, expensive)
- Microsoft Sentinel (cloud-native)
- IBM QRadar
- LogRhythm
- Exabeam
```

---

## Setting Up a Home SOC (ELK Stack)

```bash
# docker-compose.yml for ELK
version: '3'
services:
  elasticsearch:
    image: docker.elastic.co/elasticsearch/elasticsearch:8.10.0
    environment:
      - discovery.type=single-node
      - xpack.security.enabled=false
    ports: ["9200:9200"]
    
  logstash:
    image: docker.elastic.co/logstash/logstash:8.10.0
    volumes:
      - ./logstash.conf:/usr/share/logstash/pipeline/logstash.conf
    ports: ["5044:5044"]
    
  kibana:
    image: docker.elastic.co/kibana/kibana:8.10.0
    ports: ["5601:5601"]
    depends_on: [elasticsearch]

# logstash.conf
input {
  beats { port => 5044 }
  syslog { port => 5140 type => "syslog" }
}
filter {
  if [type] == "syslog" {
    grok {
      match => { "message" => "%{SYSLOGTIMESTAMP:timestamp} %{WORD:host} %{GREEDYDATA:msg}" }
    }
  }
}
output {
  elasticsearch {
    hosts => ["elasticsearch:9200"]
    index => "logs-%{+YYYY.MM.dd}"
  }
}

# Send logs from Linux hosts with Filebeat
filebeat.inputs:
- type: log
  paths:
    - /var/log/auth.log
    - /var/log/syslog
output.logstash:
  hosts: ["siem:5044"]
```

---

## Detection Engineering

Writing detection rules (SIEM correlation rules):

```
SIEM Rule: Brute Force SSH Login
─────────────────────────────────
Trigger: More than 5 failed SSH logins from same IP in 2 minutes
Logic:
  source.ip:* AND
  event.outcome: failure AND
  process.name: sshd
  | stats count() by source.ip
  | where count > 5

Alert: HIGH — Potential SSH Brute Force
Recommended Action: Block IP, investigate source

---

SIEM Rule: New Admin Account Created
─────────────────────────────────────
Trigger: New user added to Administrators group
Logic (Windows):
  event.code: 4728 OR 4732  (member added to group)
  winlog.event_data.GroupName: Administrators
  NOT winlog.event_data.SubjectUserName: scheduled_task_service

Alert: HIGH — Privileged Account Created
Recommended Action: Verify change was authorized

---

SIEM Rule: Log Cleared
──────────────────────
Trigger: Security event log cleared
Logic:
  event.code: 1102  (Windows Security log cleared)
  OR event.code: 104 (Windows System log cleared)

Alert: CRITICAL — Evidence Tampered With
Recommended Action: Immediately start IR process
```

---

## SOC Metrics

```
Key metrics tracked by SOC managers:
├── Mean Time to Detect (MTTD)       ← how long before we find an incident?
├── Mean Time to Respond (MTTR)      ← how long to contain after detection?
├── Mean Time to Contain (MTTC)
├── False Positive Rate              ← percentage of alerts that aren't real
├── Alert volume per analyst per day ← if >100/analyst, something's wrong
├── Detection coverage               ← % of ATT&CK techniques we can detect
└── MTTR by severity                 ← critical incidents handled faster
```

---

## Go: SIEM-like Log Correlator

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "strings"
    "time"
)

type Alert struct {
    Severity    string
    RuleName    string
    Description string
    Evidence    []string
    Time        time.Time
}

type CorrelationRule struct {
    Name      string
    Pattern   *regexp.Regexp
    Threshold int
    Window    time.Duration
    Severity  string
}

type Correlator struct {
    rules   []CorrelationRule
    events  map[string][]time.Time  // key → event timestamps
    alerts  []Alert
}

func NewCorrelator() *Correlator {
    c := &Correlator{
        events: make(map[string][]time.Time),
    }
    
    // Define correlation rules
    c.rules = []CorrelationRule{
        {
            Name:      "SSH Brute Force",
            Pattern:   regexp.MustCompile(`Failed password for .* from ([\d.]+)`),
            Threshold: 5,
            Window:    2 * time.Minute,
            Severity:  "HIGH",
        },
        {
            Name:      "Port Scan",
            Pattern:   regexp.MustCompile(`Connection from ([\d.]+) port \d+: Connection refused`),
            Threshold: 20,
            Window:    30 * time.Second,
            Severity:  "MEDIUM",
        },
    }
    
    return c
}

func (c *Correlator) ProcessLine(line string, t time.Time) {
    for _, rule := range c.rules {
        match := rule.Pattern.FindStringSubmatch(line)
        if len(match) < 2 {
            continue
        }
        
        key := rule.Name + ":" + match[1]  // rule + source IP
        
        // Add event
        c.events[key] = append(c.events[key], t)
        
        // Prune old events outside window
        cutoff := t.Add(-rule.Window)
        var recent []time.Time
        for _, et := range c.events[key] {
            if et.After(cutoff) {
                recent = append(recent, et)
            }
        }
        c.events[key] = recent
        
        // Check threshold
        if len(recent) >= rule.Threshold {
            c.alerts = append(c.alerts, Alert{
                Severity:    rule.Severity,
                RuleName:    rule.Name,
                Description: fmt.Sprintf("%d events in %s from %s", len(recent), rule.Window, match[1]),
                Time:        t,
            })
            // Reset to avoid duplicate alerts
            c.events[key] = nil
            fmt.Printf("[ALERT %s] %s: %s\n", rule.Severity, rule.Name,
                fmt.Sprintf("%d events in %s from %s", len(recent), rule.Window, match[1]))
        }
    }
}

func (c *Correlator) ProcessFile(path string) {
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("Error: %v\n", err)
        return
    }
    defer f.Close()
    
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        // Simplified timestamp parsing (real implementation would parse actual log timestamps)
        c.ProcessLine(line, time.Now())
        _ = strings.TrimSpace(line)
    }
}

func main() {
    c := NewCorrelator()
    c.ProcessFile("/var/log/auth.log")
    
    fmt.Printf("\nTotal alerts: %d\n", len(c.alerts))
}
```

---

## Analyst Day-to-Day

```
Morning Handoff (from night shift):
- Review incidents still open
- Check metric dashboards
- Note any ongoing attacks

Alert Triage (continuous):
1. Open alert
2. Look at the raw log — does this make sense?
3. Gather context: What is this host? Who uses it? 
4. Is this a known false positive?
5. Investigate: check related events, pivot on IPs/hashes/users
6. Decision: close (FP) or escalate (TP)

Typical ratio: 80-95% of alerts are false positives
Goal: find the 5-20% that are real as fast as possible
```

---

## Summary

| SOC Role | Main Focus | Key Tools |
|----------|-----------|-----------|
| Tier 1 | Alert triage, 24/7 coverage | SIEM, ticketing |
| Tier 2 | Incident investigation | EDR, threat intel, forensics |
| Tier 3 | Threat hunting, rule writing | SIEM queries, custom tooling |
| Detection Eng | Write better rules | SIEM, sigma rules, ATT&CK |

---

## Exercises

1. Deploy ELK Stack with Docker Compose and forward your Linux system's auth.log into it
2. Write 5 detection rules in Kibana for common attack patterns (brute force, large file transfer, new user created)
3. Practice with Blue Team Labs Online (blueteamlabs.online) — SOC analyst challenges
4. Calculate your home SIEM's false positive rate — how many real vs. false alerts per day?
