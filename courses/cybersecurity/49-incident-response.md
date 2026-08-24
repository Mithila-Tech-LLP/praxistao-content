# Chapter 49: Incident Response — Handling Breaches Like a Pro

*An incident isn't if, it's when. Incident response (IR) is the structured process for detecting, containing, eradicating, and recovering from security breaches. Done well, it limits damage and restores operations. Done poorly, it makes everything worse.*

---

## The IR Lifecycle (NIST SP 800-61)

```
        ┌─────────────────────────────────────────────┐
        │                                             │
        ▼                                             │
  1. Preparation     2. Detection &    3. Containment,
                        Analysis       Eradication &
                                       Recovery
                                              │
                                              ▼
                                    4. Post-Incident
                                       Activity
```

### Phase 1: Preparation

Before an incident happens:
- Create and test IR plan
- Build IR team (assign roles)
- Set up logging and monitoring
- Define "crown jewels" (critical assets)
- Establish communication channels
- Practice tabletop exercises

### Phase 2: Detection and Analysis

```bash
# Signs of compromise (IoCs):
# - Unusual outbound traffic
# - New user accounts created
# - Scheduled tasks / cron jobs added
# - New services installed
# - Antivirus disabled
# - Large data transfers
# - Off-hours logins
# - New SSH authorized_keys

# Initial triage
who               # who's logged in right now
last              # recent logins
netstat -tulpn    # what's listening?
netstat -an       # all connections
ps aux            # running processes
crontab -l        # current user cron
ls /etc/cron.d/   # system cron

# Check logs
tail -100 /var/log/auth.log
grep "Failed password" /var/log/auth.log | awk '{print $11}' | sort | uniq -c | sort -rn
grep "Accepted password" /var/log/auth.log
journalctl -xe --since "1 hour ago"
```

### Phase 3: Containment

**Short-term containment (stop bleeding):**
```bash
# Isolate the host (remove from network)
# If possible, DO NOT shutdown — memory evidence is lost
# Take a memory snapshot if possible

# Firewall rules to isolate
iptables -I INPUT -j DROP
iptables -I OUTPUT -j DROP
iptables -I INPUT -s MANAGEMENT_IP -j ACCEPT

# Block attacker IP (if known)
iptables -I INPUT -s ATTACKER_IP -j DROP
iptables -I OUTPUT -d ATTACKER_IP -j DROP

# Change compromised credentials
passwd compromised_user
# Rotate API keys, service account passwords
```

**Eradication:**
```bash
# Find and remove malware
find / -newer /tmp/reference_time -type f 2>/dev/null   # recently created files
find / -name ".*" -type f 2>/dev/null                    # hidden files
find /tmp /var/tmp -executable -type f                   # executables in temp
chkrootkit                                               # check for rootkits
rkhunter --check                                         # another rootkit check

# Remove persistence
crontab -r                          # remove cron
systemctl disable malicious.service
systemctl stop malicious.service
rm /etc/systemd/system/malicious.service
```

### Phase 4: Post-Incident

```
Lessons Learned Report:
1. What happened? (timeline)
2. What was the root cause?
3. What worked well in the response?
4. What didn't work?
5. What changes are needed?

Metrics to track:
- Mean Time to Detect (MTTD)
- Mean Time to Respond (MTTR)
- Mean Time to Contain (MTTC)
```

---

## IR Playbooks

A playbook is a pre-written procedure for a specific incident type:

### Ransomware Playbook

```
DETECT
- Alert: unusual file extensions (.encrypted, .locked) appearing en masse
- Alert: Volume Shadow Copy deletion (vssadmin delete shadows)
- Alert: processes encrypting large numbers of files rapidly

IMMEDIATE (first 15 minutes)
1. Do NOT pay the ransom immediately
2. Isolate affected systems from network
3. Alert IR team lead and management
4. Preserve logs and memory if possible
5. Check backup status (are they intact? offline?)

INVESTIGATE
1. Identify patient zero — which host was first affected?
2. Determine variant (ID Ransomware website)
3. Check for decryptor (nomoreransom.org)
4. Review network logs for C2 communication
5. Determine if data was exfiltrated before encryption

RECOVER
1. Restore from clean backups
2. Test restored systems before reconnecting to network
3. Force password resets for all affected accounts
4. Patch the initial access vulnerability
```

---

## Go: IR Triage Tool

```go
package main

import (
    "fmt"
    "net"
    "os"
    "os/exec"
    "path/filepath"
    "strings"
    "time"
)

type TriageReport struct {
    Hostname    string
    TimeStarted time.Time
    Connections []string
    SuspiciousProcs []string
    NewCronJobs []string
    SuspiciousFiles []string
}

func getConnections() []string {
    // Parse /proc/net/tcp for active connections
    data, err := os.ReadFile("/proc/net/tcp")
    if err != nil {
        return nil
    }
    
    var conns []string
    for _, line := range strings.Split(string(data), "\n")[1:] {
        fields := strings.Fields(line)
        if len(fields) < 4 {
            continue
        }
        if fields[3] == "01" { // ESTABLISHED state
            remoteHex := fields[2]
            // Parse hex IP:port
            if ip, port := parseHexAddr(remoteHex); port != 0 {
                // Filter out loopback and private
                if !isPrivate(ip) {
                    conns = append(conns, fmt.Sprintf("%s:%d", ip.String(), port))
                }
            }
        }
    }
    return conns
}

func parseHexAddr(hexAddr string) (net.IP, int) {
    parts := strings.Split(hexAddr, ":")
    if len(parts) != 2 {
        return nil, 0
    }
    
    var ipBytes [4]byte
    fmt.Sscanf(parts[0], "%08x", (*uint32)(nil))
    // Simplified — real implementation parses little-endian hex
    return net.IP{0, 0, 0, 0}, 0
}

func isPrivate(ip net.IP) bool {
    privateRanges := []string{"10.", "192.168.", "172.16.", "127."}
    s := ip.String()
    for _, r := range privateRanges {
        if strings.HasPrefix(s, r) {
            return true
        }
    }
    return false
}

func checkSuspiciousFiles() []string {
    var suspicious []string
    
    // Recently created executables in temp dirs
    dirs := []string{"/tmp", "/var/tmp", "/dev/shm"}
    for _, dir := range dirs {
        filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
            if err != nil || info.IsDir() {
                return nil
            }
            // Executable file in temp dir
            if info.Mode()&0111 != 0 {
                suspicious = append(suspicious,
                    fmt.Sprintf("Executable in temp: %s (size: %d, modified: %s)",
                        path, info.Size(), info.ModTime().Format("2006-01-02 15:04:05")))
            }
            return nil
        })
    }
    
    // Hidden files
    homeEntries, _ := os.ReadDir("/home")
    for _, u := range homeEntries {
        userHome := filepath.Join("/home", u.Name())
        entries, _ := os.ReadDir(userHome)
        for _, e := range entries {
            if strings.HasPrefix(e.Name(), ".") && e.Name() != ".bashrc" &&
                e.Name() != ".bash_history" && e.Name() != ".profile" &&
                e.Name() != ".ssh" && e.Name() != ".config" {
                suspicious = append(suspicious,
                    fmt.Sprintf("Unusual hidden file: %s", filepath.Join(userHome, e.Name())))
            }
        }
    }
    
    return suspicious
}

func runCommand(name string, args ...string) string {
    out, err := exec.Command(name, args...).Output()
    if err != nil {
        return ""
    }
    return string(out)
}

func main() {
    hostname, _ := os.Hostname()
    report := &TriageReport{
        Hostname:    hostname,
        TimeStarted: time.Now(),
    }
    
    report.SuspiciousFiles = checkSuspiciousFiles()
    
    fmt.Printf("=== IR TRIAGE REPORT ===\n")
    fmt.Printf("Host: %s\n", report.Hostname)
    fmt.Printf("Time: %s\n\n", report.TimeStarted.Format(time.RFC3339))
    
    // Current connections
    fmt.Println("RECENT LOGINS:")
    fmt.Println(runCommand("last", "-10"))
    
    fmt.Println("LISTENING SERVICES:")
    fmt.Println(runCommand("ss", "-tulpn"))
    
    if len(report.SuspiciousFiles) > 0 {
        fmt.Println("SUSPICIOUS FILES:")
        for _, f := range report.SuspiciousFiles {
            fmt.Printf("  [!] %s\n", f)
        }
    }
}
```

---

## Summary

| IR Phase | Key Activities | Tools |
|----------|---------------|-------|
| Preparation | IR plan, logging, training | Playbooks, SIEM |
| Detection | Alert triage, log analysis | SIEM, EDR alerts |
| Containment | Isolate host, block attacker | Firewall, EDR |
| Eradication | Remove malware, patch vuln | AV, manual analysis |
| Recovery | Restore backups, verify | Backups, monitoring |
| Post-incident | Lessons learned, metrics | Reports |

---

## Exercises

1. Build a personal IR playbook for a ransomware incident in a small company
2. Practice the IR triage tool on a Linux VM — what does it report on a clean system?
3. Download a CTF memory dump and use volatility3 to find malicious processes
4. Participate in a free tabletop exercise (CISA has free materials) — practice the decision-making process
