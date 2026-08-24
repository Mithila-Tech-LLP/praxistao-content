# Chapter 51: Threat Intelligence — Knowing Your Enemy

*Threat intelligence transforms raw data (IPs, file hashes, domain names) into actionable knowledge: who is attacking you, how, and what to do about it.*

---

## Intelligence Levels

```
Strategic (executive level)
- "State actors from Country X target our industry"
- "Ransomware groups increasingly target healthcare"
- Informs: security investments, policies, risk decisions

Operational (security team level)
- "Lazarus Group is currently using spear-phishing with PDF attachments"
- "This APT pivots via SMB after initial access"
- Informs: threat hunting hypotheses, detection focus

Tactical (analyst level — most common)
- IP addresses of known C2 servers
- File hashes of known malware
- Domain names associated with phishing
- Informs: blocking rules, SIEM detections, alerts
```

---

## Indicators of Compromise (IoCs)

```
Types:
├── IP addresses      — C2 servers, attacker infrastructure
├── Domain names      — malicious domains, phishing sites
├── File hashes       — known malware (MD5, SHA-256)
├── URLs              — phishing URLs, malware download sites
├── Registry keys     — malware persistence locations
├── Mutex names       — process identifiers used by specific malware
├── Email addresses   — phishing sender addresses
└── SSL cert hashes   — malicious server certificates
```

---

## Threat Intelligence Sources

```bash
# Free sources:
# AlienVault OTX — community threat intel
# VirusTotal — file/URL/IP reputation
# AbuseIPDB — IP reputation crowdsourced
# URLhaus — malicious URL database
# ThreatFox — IoCs from Abuse.ch
# Feodo Tracker — botnet C2 IPs
# Shodan Monitor — monitor your IPs for changes

# Paid / Premium:
# Recorded Future
# Mandiant Advantage
# CrowdStrike Falcon X
# Threatconnect
# MISP (free, self-hosted sharing platform)

# Government:
# CISA AIS (Automated Indicator Sharing)
# NCSC UK feeds
# ANSSI (France)
```

---

## MITRE ATT&CK for Threat Intelligence

ATT&CK maps adversary behaviors beyond IoCs:

```python
# Threat group profile example:
# APT28 (Fancy Bear — Russian GRU)
# Uses: T1566.001 (Spearphishing), T1078 (Valid Accounts)
# Tools: X-Agent (malware), Mimikatz, Empire
# Targets: Government, Military, Aerospace

# Check if your environment is tested against APT28:
# ATT&CK Navigator: navigator.attack.mitre.org
# Color techniques red if detected, yellow if partial
```

---

## STIX/TAXII — Intelligence Standards

```
STIX (Structured Threat Information Expression)
├── Format for describing threats in JSON
├── Objects: Indicator, Malware, Threat Actor, Campaign, etc.

TAXII (Trusted Automated eXchange of Intelligence Information)
├── Protocol for sharing STIX data
├── Collections, channels, subscriptions

MISP (Malware Information Sharing Platform)
├── Open source threat intel sharing
├── Self-hosted
├── Syncs with multiple feeds
```

---

## Go: IoC Feed Checker

```go
package main

import (
    "bufio"
    "encoding/json"
    "fmt"
    "net/http"
    "os"
    "strings"
    "time"
)

type AbuseIPResponse struct {
    Data struct {
        AbuseConfidenceScore int    `json:"abuseConfidenceScore"`
        CountryCode          string `json:"countryCode"`
        UsageType            string `json:"usageType"`
        ISP                  string `json:"isp"`
        TotalReports         int    `json:"totalReports"`
    } `json:"data"`
}

type ThreatFeedEntry struct {
    Value    string
    Type     string
    Tags     []string
    Severity string
}

// Load IoCs from file (one per line)
func loadIoCFeed(path string) []ThreatFeedEntry {
    f, err := os.Open(path)
    if err != nil {
        return nil
    }
    defer f.Close()
    
    var entries []ThreatFeedEntry
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := strings.TrimSpace(scanner.Text())
        if line == "" || strings.HasPrefix(line, "#") {
            continue
        }
        entries = append(entries, ThreatFeedEntry{
            Value: line,
            Type:  detectIoCType(line),
        })
    }
    return entries
}

func detectIoCType(value string) string {
    if strings.Count(value, ".") == 3 {
        // Likely IP (simplified)
        return "ip"
    }
    if strings.HasPrefix(value, "http") {
        return "url"
    }
    if len(value) == 64 {
        return "sha256"
    }
    if len(value) == 32 {
        return "md5"
    }
    return "domain"
}

func checkIPReputation(ip, apiKey string) *AbuseIPResponse {
    client := &http.Client{Timeout: 5 * time.Second}
    
    req, _ := http.NewRequest("GET",
        "https://api.abuseipdb.com/api/v2/check?ipAddress="+ip, nil)
    req.Header.Set("Key", apiKey)
    req.Header.Set("Accept", "application/json")
    
    resp, err := client.Do(req)
    if err != nil {
        return nil
    }
    defer resp.Body.Close()
    
    var result AbuseIPResponse
    json.NewDecoder(resp.Body).Decode(&result)
    return &result
}

func checkNetworkTrafficAgainstFeed(logFile string, feed []ThreatFeedEntry) {
    f, _ := os.Open(logFile)
    defer f.Close()
    
    // Build lookup map for O(1) checking
    iocMap := make(map[string]ThreatFeedEntry)
    for _, entry := range feed {
        iocMap[entry.Value] = entry
    }
    
    scanner := bufio.NewScanner(f)
    lineNum := 0
    for scanner.Scan() {
        lineNum++
        line := scanner.Text()
        
        for ioc, entry := range iocMap {
            if strings.Contains(line, ioc) {
                fmt.Printf("[MATCH] Line %d: IoC type=%s value=%s found in: %s\n",
                    lineNum, entry.Type, ioc, line[:min(len(line), 100)])
            }
        }
    }
}

func min(a, b int) int {
    if a < b {
        return a
    }
    return b
}

func main() {
    // Example IoC feed (normally loaded from MISP, OTX, etc.)
    feed := []ThreatFeedEntry{
        {Value: "185.220.101.1", Type: "ip", Tags: []string{"tor-exit"}, Severity: "medium"},
        {Value: "malicious-domain.ru", Type: "domain", Tags: []string{"c2"}, Severity: "high"},
        {Value: "a665a45920422f9d417e4867efdc4fb8", Type: "md5", Tags: []string{"ransomware"}, Severity: "critical"},
    }
    
    fmt.Printf("Loaded %d IoCs\n", len(feed))
    
    // In production: 
    // feed = loadIoCFeed("/etc/goshield/ioc-feed.txt")
    // checkNetworkTrafficAgainstFeed("/var/log/firewall.log", feed)
    
    for _, entry := range feed {
        fmt.Printf("  [%s] %s: %s\n", entry.Severity, entry.Type, entry.Value)
    }
}
```

---

## Diamond Model of Intrusion Analysis

```
        Adversary
       /          \
      /            \
Capability ---- Victim
      \            /
       \          /
      Infrastructure
```

For every attack, analyze all four elements:
- **Adversary:** Who is attacking? (APT group, criminal, insider)
- **Capability:** What tools/techniques do they use?
- **Infrastructure:** What servers/domains do they use?
- **Victim:** Why this target? What do they want?

---

## Summary

| Intel Level | Audience | Example | Freshness |
|-------------|----------|---------|-----------|
| Strategic | Executives | "APT28 targets energy" | Months |
| Operational | Security teams | "Using spearphishing + PDF" | Weeks |
| Tactical (IoCs) | Analysts, tools | "C2 IP: 1.2.3.4" | Hours/days |

---

## Exercises

1. Sign up for a free MISP instance (misp-project.org) or use their demo — explore how threat sharing works
2. Fetch the Feodo Tracker botnet C2 IP list and write a Go script to check if any of your network connections match
3. Look up an APT group in ATT&CK (e.g., APT29/Cozy Bear) — map their techniques to your detection gaps
4. Integrate the Go IoC checker into GoShield as a real-time network connection checker
