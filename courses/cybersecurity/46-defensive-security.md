# Chapter 46: Defensive Security — Defense in Depth

*Offense is about finding one way in. Defense is about closing every way in, detecting every attempt, and responding before damage is done. Defense in depth means no single control protects everything.*

---

## The Attacker's Kill Chain (What Defenders Stop)

The **MITRE ATT&CK Framework** maps adversary behavior into tactics and techniques:

```
Reconnaissance      → Gather target information
Resource Development → Set up infrastructure (C2 servers, domains)
Initial Access      → Get a foothold (phishing, exploit, stolen creds)
Execution           → Run malicious code
Persistence         → Survive reboots, re-infection
Privilege Escalation→ Get more access
Defense Evasion     → Avoid detection
Credential Access   → Steal passwords, hashes
Discovery           → Map the environment
Lateral Movement    → Move to other systems
Collection          → Gather data
Command and Control → Maintain contact with attacker infrastructure
Exfiltration        → Steal the data
Impact              → Destroy, encrypt, deface
```

**Key insight:** Stop the attacker at ANY stage and you win. Defenders don't need to stop everything — they need to stop or detect anything before impact.

---

## Defense in Depth

The principle: multiple overlapping layers, so that failing one layer doesn't mean total compromise.

```
┌─────────────────────────────────────────────────────┐
│  Physical Security (locked rooms, badge access)      │
│ ┌─────────────────────────────────────────────────┐ │
│ │  Perimeter (Firewall, IDS/IPS, DDoS protection) │ │
│ │ ┌───────────────────────────────────────────┐  │ │
│ │ │  Network (VLANs, NAC, encryption)          │  │ │
│ │ │ ┌───────────────────────────────────────┐ │  │ │
│ │ │ │  Endpoint (EDR, antivirus, patching)   │ │  │ │
│ │ │ │ ┌─────────────────────────────────┐   │ │  │ │
│ │ │ │ │  Application (WAF, SAST, DAST)  │   │ │  │ │
│ │ │ │ │ ┌───────────────────────────┐   │   │ │  │ │
│ │ │ │ │ │  Data (encryption at rest,│   │   │ │  │ │
│ │ │ │ │ │   DLP, backups)           │   │   │ │  │ │
│ │ │ │ │ └───────────────────────────┘   │   │ │  │ │
│ │ │ │ └─────────────────────────────────┘   │ │  │ │
│ │ │ └───────────────────────────────────────┘ │  │ │
│ │ └───────────────────────────────────────────┘  │ │
│ └─────────────────────────────────────────────────┘ │
└─────────────────────────────────────────────────────┘
```

---

## Perimeter Security

### Firewall Rules

Think of firewalls as allow-lists, not block-lists.

```
Default-deny: block everything, explicitly allow what's needed

GOOD rules:
ALLOW  any → 10.0.1.0/24:80,443        # Web servers receive HTTP/HTTPS
ALLOW  any → 10.0.1.0/24:25            # Mail servers receive email  
ALLOW  10.0.2.0/24 → 10.0.1.0/24:3306 # App servers reach database
DENY   any → any                        # Everything else blocked

BAD:
ALLOW  any → any:3306                   # Database exposed to internet!
ALLOW  any → 10.0.0.0/8:22             # SSH from anywhere!
```

### IDS vs IPS

| | IDS (Intrusion Detection) | IPS (Intrusion Prevention) |
|--|--------------------------|---------------------------|
| Action | Alerts only | Blocks traffic |
| Risk | No false-positive disruption | Can block legitimate traffic |
| Position | Out-of-band (tap/SPAN) | Inline (traffic passes through) |
| Tools | Zeek, Snort (IDS mode) | Snort/Suricata (IPS mode), Palo Alto |

**Snort/Suricata rule example:**
```
# Detect port scan
alert tcp any any -> $HOME_NET any (
    msg:"Port Scan Detected";
    flags:S;
    threshold:type both, track by_src, count 20, seconds 2;
    sid:9001;
    rev:1;
)

# Detect SQL injection attempt
alert http any any -> $HTTP_SERVERS $HTTP_PORTS (
    msg:"SQL Injection Attempt";
    content:"'OR 1=1";
    http_uri;
    nocase;
    sid:9002;
)
```

---

## Endpoint Security

### Hardening a Linux Server

```bash
# 1. Keep everything patched
apt update && apt upgrade -y
# Configure unattended-upgrades for security patches

# 2. Disable unnecessary services
systemctl list-units --type=service --state=running
systemctl disable bluetooth
systemctl disable cups
# If you don't need it, disable it — every service is attack surface

# 3. Configure SSH properly
cat /etc/ssh/sshd_config
# Recommended settings:
PermitRootLogin no           # Don't allow root SSH
PasswordAuthentication no    # Key-based only
MaxAuthTries 3               # Limit brute force
AllowUsers deploy admin      # Whitelist users
Port 2222                    # Non-standard port (minimal security, just reduces noise)

# 4. Firewall with iptables/ufw
ufw default deny incoming
ufw default allow outgoing
ufw allow 22/tcp
ufw allow 443/tcp
ufw enable

# 5. Limit SUID binaries
find / -perm -4000 -type f > /root/suid_baseline.txt
# Review the list, remove SUID from anything that doesn't need it

# 6. Configure audit logging
auditctl -w /etc/passwd -p wa -k passwd_changes
auditctl -w /etc/sudoers -p wa -k sudoers_changes
auditctl -a always,exit -F arch=b64 -S execve -k command_execution
```

### Hardening Windows

```powershell
# Enable Windows Defender
Set-MpPreference -DisableRealtimeMonitoring $false

# Enable audit logging
auditpol /set /category:"Logon/Logoff" /success:enable /failure:enable
auditpol /set /category:"Process Creation" /success:enable

# Disable SMBv1 (EternalBlue was SMBv1)
Set-SmbServerConfiguration -EnableSMB1Protocol $false

# Restrict PowerShell execution
Set-ExecutionPolicy AllSigned -Force

# Enable PowerShell logging (catches most living-off-the-land attacks)
# Group Policy: Administrative Templates → Windows Components → PowerShell
# Enable Script Block Logging

# AppLocker / WDAC — whitelist applications
# Only allowed executables run; prevents most malware
```

---

## Log Management

Logs are your digital eyewitness testimony. Without logs, you're blind.

### What to Log

**Authentication events:**
- Successful logins (when, from where, which user)
- Failed logins (brute force detection)
- Privilege escalation (sudo usage)

**Process execution:**
- What commands ran, with what arguments
- Parent-child relationships (PowerShell spawning cmd is suspicious)

**Network connections:**
- Outbound connections from servers (beaconing)
- DNS queries (domain generation algorithms)
- Data volumes (exfiltration)

**File system changes:**
- Changes to critical files (/etc/passwd, binaries)
- New executables created

### Linux Log Locations

```
/var/log/auth.log      → authentication events (SSH, sudo)
/var/log/syslog        → general system events
/var/log/kern.log      → kernel messages
/var/log/nginx/        → web server access/error
/var/log/mysql/        → database logs
/var/log/audit/audit.log → audit framework
```

### Practical Log Analysis

```bash
# Find all failed SSH logins
grep "Failed password" /var/log/auth.log | \
    awk '{print $11}' | sort | uniq -c | sort -rn

# Find successful logins after failures (successful brute force)
grep "Accepted" /var/log/auth.log | awk '{print $11}' | \
while read ip; do
    count=$(grep "Failed password.*$ip" /var/log/auth.log | wc -l)
    echo "$ip: $count failed before success"
done

# Detect privilege escalation
grep "sudo" /var/log/auth.log | grep -v "session opened"

# Find unusual parent-child in audit log
ausearch -sc execve | grep -A 5 "bash.*python"

# Timeline of events around an incident
grep "2024-01-15 14:" /var/log/auth.log /var/log/syslog | sort
```

---

## SIEM — Security Information and Event Management

A SIEM aggregates logs from all sources, normalizes them, and allows detection rules.

```
Servers → Log → Syslog/Beats Agent → SIEM → Detection Engine → Alerts → SOC
Windows →     →                   →      →                  →        →
Network →     →                   →      →                  →        →
Apps    →     →                   →      →                  →        →
```

**Open source options:**
- **Elastic Stack (ELK):** Elasticsearch + Logstash + Kibana
- **Wazuh:** Security-focused SIEM built on Elastic
- **Graylog:** Good for smaller teams

**Commercial:**
- Splunk, IBM QRadar, Microsoft Sentinel, Exabeam

### Detection Rules Example (Elastic KQL)

```
# SSH brute force detection
event.category: "authentication" AND event.outcome: "failure"
  | stats count by source.ip
  | where count > 20

# Suspicious process creation
process.name: "cmd.exe" AND process.parent.name: "outlook.exe"
# Outlook spawning CMD = spearphishing

# PowerShell with encoded commands (common malware technique)
process.name: "powershell.exe" AND
process.command_line: *-EncodedCommand*

# Lateral movement via PsExec
process.name: "psexesvc.exe" OR
event.action: "network_connection" AND 
process.name: "psexec.exe"
```

---

## Incident Response

When something bad happens, you need a plan.

### IR Phases

```
1. Preparation   → IR playbooks, tools ready, team trained
2. Identification → Determine if incident is real; scope it
3. Containment   → Stop the spread (isolate affected systems)
4. Eradication   → Remove the threat (malware, backdoors)
5. Recovery      → Restore systems to production
6. Lessons Learned → Post-mortem, improve defenses
```

### Triage: Is This Real?

```bash
# Check currently logged-in users
who
w
last | head -20

# Check recent file modifications
find / -newer /tmp/baseline -type f 2>/dev/null | head -50

# Check network connections
ss -tp | grep ESTABLISHED

# Check suspicious processes
ps aux | sort -rk3 | head -20    # CPU usage
ps aux | sort -rk4 | head -20    # memory usage

# Check cron for new jobs
diff /tmp/cron_baseline /etc/crontab

# Check for new SUID binaries
diff /tmp/suid_baseline <(find / -perm -4000 2>/dev/null)

# Capture memory (for forensics)
# Use LiME (Linux Memory Extractor) kernel module
```

### Containment Options

```bash
# Option 1: Isolate from network
# Firewall off or physically disconnect

# Option 2: Block specific IP
iptables -I INPUT -s 185.220.101.5 -j DROP
iptables -I OUTPUT -d 185.220.101.5 -j DROP

# Option 3: Kill suspicious process
kill -9 <PID>

# Option 4: Snapshot VM before cleanup (for forensics)
# Keep evidence!
```

---

## Zero Trust Architecture

Traditional security: "everything inside the firewall is trusted."

Zero Trust: **"never trust, always verify"**

Key principles:
1. **Verify explicitly** — authenticate and authorize every request, every time
2. **Least privilege** — minimum access needed to do the job
3. **Assume breach** — design as if attackers are already inside

```
Traditional model:
  Internet → Firewall → [TRUSTED ZONE — everything is trusted]

Zero Trust:
  Internet → Identity verification → MFA → Device health check
           → Just-in-time access → Encrypted connection → Specific resource
           → Audit every access
```

**Implementation building blocks:**
- **MFA everywhere** — even for VPN, admin consoles
- **PAM (Privileged Access Management)** — vault for admin passwords
- **Microsegmentation** — no lateral movement between segments
- **EDR on every endpoint**
- **All traffic encrypted**

---

## Security Checklist (Defense Quick Reference)

### For Any Server

```
☐ OS patched and auto-updating
☐ Only necessary ports open (firewall default-deny)
☐ SSH: key-based only, no root login
☐ All services running as least-privilege users
☐ No SUID binaries beyond system defaults
☐ Audit logging enabled and shipped off-host
☐ File integrity monitoring (AIDE or Wazuh)
☐ Outbound connections monitored
```

### For Web Applications

```
☐ WAF in front (rate limiting, SQL injection blocking)
☐ Security headers (CSP, HSTS, X-Frame-Options)
☐ HTTPS everywhere, HTTP redirects to HTTPS
☐ Parameterized queries (no SQL injection)
☐ Input validation on all fields
☐ Authentication: strong passwords + MFA
☐ Session cookies: HttpOnly + Secure + SameSite
☐ Dependency scanning (no known-vulnerable libraries)
☐ Regular DAST scanning
```

---

## Building a Simple Log Analyzer in Go

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"
)

var failedSSH = regexp.MustCompile(`Failed password.*from (\d+\.\d+\.\d+\.\d+)`)
var acceptedSSH = regexp.MustCompile(`Accepted.*from (\d+\.\d+\.\d+\.\d+)`)

type SSHStats struct {
    Failures map[string]int
    Successes map[string]int
}

func analyzeAuthLog(path string) (*SSHStats, error) {
    f, err := os.Open(path)
    if err != nil {
        return nil, err
    }
    defer f.Close()

    stats := &SSHStats{
        Failures:  make(map[string]int),
        Successes: make(map[string]int),
    }

    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        if m := failedSSH.FindStringSubmatch(line); m != nil {
            stats.Failures[m[1]]++
        }
        if m := acceptedSSH.FindStringSubmatch(line); m != nil {
            stats.Successes[m[1]]++
        }
    }
    return stats, scanner.Err()
}

type IPCount struct {
    IP    string
    Count int
}

func topN(counts map[string]int, n int) []IPCount {
    var list []IPCount
    for ip, count := range counts {
        list = append(list, IPCount{ip, count})
    }
    sort.Slice(list, func(i, j int) bool {
        return list[i].Count > list[j].Count
    })
    if len(list) > n {
        list = list[:n]
    }
    return list
}

func main() {
    logPath := "/var/log/auth.log"
    if len(os.Args) > 1 {
        logPath = os.Args[1]
    }

    stats, err := analyzeAuthLog(logPath)
    if err != nil {
        fmt.Fprintf(os.Stderr, "Error: %v\n", err)
        os.Exit(1)
    }

    fmt.Println("=== SSH Authentication Analysis ===\n")
    
    fmt.Println("Top 10 IPs with Failed Logins:")
    for _, item := range topN(stats.Failures, 10) {
        bar := strings.Repeat("#", min(item.Count/10, 50))
        fmt.Printf("  %-20s %5d %s\n", item.IP, item.Count, bar)
    }

    fmt.Println("\nSuccessful Logins from IPs with Prior Failures:")
    for _, item := range topN(stats.Successes, 20) {
        if failCount := stats.Failures[item.IP]; failCount > 0 {
            fmt.Printf("  [!] %-20s succeeded after %d failures\n", item.IP, failCount)
        }
    }
}

func min(a, b int) int {
    if a < b { return a }
    return b
}
```

---

## Summary

| Defense layer | Key tools/techniques |
|--------------|---------------------|
| Perimeter | Firewall, IDS/IPS, DDoS protection |
| Network | VLANs, encryption, NAC |
| Endpoint | EDR (GoShield!), patch management, hardening |
| Application | WAF, SAST, DAST, secure coding |
| Data | Encryption at rest/transit, DLP, backups |
| Identity | MFA, PAM, least privilege |
| Detection | SIEM, log analysis, detection rules |
| Response | IR playbook, forensics, containment |

---

## Exercises

1. Set up Wazuh on a VM. Connect a Linux agent. What alerts does it generate just from normal system activity?
2. Analyze an auth.log file. Find the top attacking IPs. Were any successful after failures?
3. Harden an Ubuntu server using the checklist above. Run Lynis and score your hardening.
4. Write a Snort/Suricata rule that detects HTTP requests with common SQL injection patterns.
5. Build an IR runbook for a "server compromised via web shell" scenario.
