# Chapter 40: Covering Tracks — Removing Evidence of Compromise

*Attackers remove evidence to delay detection and investigation. Defenders preserve evidence to understand what happened. Understanding both sides makes you a better incident responder.*

---

## What Evidence Exists?

```
Log Files           — auth.log, syslog, apache access.log, Windows Event Logs
Bash History        — ~/.bash_history
Command History     — .zsh_history, .fish_history
Web Logs            — /var/log/apache2/, /var/log/nginx/
Auth Logs           — /var/log/auth.log, /var/log/secure
Last Login          — /var/log/wtmp (who/last), /var/log/btmp (failed)
Process History     — /proc/, audit logs
Network Logs        — firewall logs, netflow, IDS alerts
File Access Times   — atime, mtime, ctime on files
Memory Artifacts    — crash dumps, swap files
```

---

## Linux Track Covering

```bash
# Clear bash history
history -c
history -w
cat /dev/null > ~/.bash_history
unset HISTFILE
# Or disable completely
export HISTFILE=/dev/null
export HISTSIZE=0

# Do NOT leave commands in history while working
HISTFILE=/dev/null bash   # spawn a no-history subshell

# Clear logs
> /var/log/auth.log
> /var/log/syslog
> /var/log/apache2/access.log
> /var/log/apache2/error.log

# More subtle: remove specific lines
sed -i '/192.168.1.200/d' /var/log/auth.log   # remove attacker IP
sed -i '/malicious/d' /var/log/apache2/access.log

# Last/who/lastb use binary log files
# Clear lastlog (shows last login for each user)
> /var/log/lastlog

# Clear wtmp (last/who history)
> /var/log/wtmp

# Clear btmp (failed logins)
> /var/log/btmp

# Restore file timestamps (after modifying logs)
touch -r /var/log/auth.log.1 /var/log/auth.log   # same timestamp as old backup
```

---

## Timestomping

```bash
# Change file timestamps to blend in
# -t format: CCYYMMDDhhmm.SS
touch -t 202201151200.00 /tmp/malware  # set to Jan 15, 2022

# Match another file's timestamps
touch -r /bin/ls /tmp/malware

# Python version
python3 -c "import os; os.utime('/tmp/malware', (1642252800, 1642252800))"
```

**Note:** ctime (change time) cannot be manipulated with standard tools — requires raw disk access. Forensic tools compare atime/mtime/ctime to detect timestomping.

---

## Windows Track Covering

```powershell
# Clear Windows Event Logs
wevtutil cl System
wevtutil cl Security
wevtutil cl Application
wevtutil el | foreach { wevtutil cl $_ }  # clear ALL logs

# PowerShell history
Remove-Item (Get-PSReadlineOption).HistorySavePath

# Prefetch (evidence of program execution)
Remove-Item C:\Windows\Prefetch\MALWARE.EXE-*.pf

# Windows logs clearing itself generates Event ID 1102 (Security) / 104 (System)
# Very suspicious → usually means active attack in progress
```

---

## Why Full Erasure Fails

Attackers rarely succeed in completely covering tracks:

```
Memory forensics      → RAM image captured before reboot
Network logs          → External firewall/router logs they don't control
SIEM                  → Logs already forwarded to remote server before deletion
Snapshot/backups      → Hypervisor/backup snapshots of disk state
Disk forensics        → Deleted files recoverable, journal shows changes
Correlation           → Small inconsistencies between log files reveal manipulation
```

---

## Detecting Log Tampering (GoShield)

```go
package main

import (
    "crypto/sha256"
    "fmt"
    "io"
    "os"
    "sync"
    "time"
)

type LogMonitor struct {
    files    map[string]string // path → last known hash
    mu       sync.Mutex
}

func NewLogMonitor() *LogMonitor {
    return &LogMonitor{
        files: make(map[string]string),
    }
}

func fileHash(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "", err
    }
    defer f.Close()
    
    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "", err
    }
    return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func (lm *LogMonitor) baseline(paths []string) {
    for _, p := range paths {
        if h, err := fileHash(p); err == nil {
            lm.mu.Lock()
            lm.files[p] = h
            lm.mu.Unlock()
        }
    }
}

func (lm *LogMonitor) check() []string {
    var alerts []string
    
    lm.mu.Lock()
    defer lm.mu.Unlock()
    
    for path, known := range lm.files {
        info, err := os.Stat(path)
        if err != nil {
            alerts = append(alerts, fmt.Sprintf("LOG DELETED: %s", path))
            continue
        }
        
        // Empty file = likely cleared
        if info.Size() == 0 {
            alerts = append(alerts, fmt.Sprintf("LOG CLEARED (empty): %s", path))
            continue
        }
        
        current, err := fileHash(path)
        if err != nil {
            continue
        }
        
        if current != known {
            // Log changed — is this expected rotation or tampering?
            // Real system: compare file size (cleared log would shrink unexpectedly)
            _ = current
            lm.files[path] = current
        }
    }
    return alerts
}

func (lm *LogMonitor) Monitor(paths []string) {
    lm.baseline(paths)
    ticker := time.NewTicker(30 * time.Second)
    
    for range ticker.C {
        alerts := lm.check()
        for _, a := range alerts {
            fmt.Printf("[ALERT] %s\n", a)
            // In GoShield: send to alerting engine
        }
    }
}

func main() {
    mon := NewLogMonitor()
    mon.Monitor([]string{
        "/var/log/auth.log",
        "/var/log/syslog",
        "/var/log/apache2/access.log",
    })
}
```

---

## Defensive Measures

```bash
# Forward logs to remote SIEM immediately (rsyslog)
# /etc/rsyslog.conf:
*.* @@siem.company.com:514     # TCP forwarding to remote syslog

# auditd — immutable audit log (cannot be cleared by non-root)
auditctl -e 2      # enable immutable mode (requires reboot to disable)

# Integrity monitoring
aide --init        # baseline
aide --check       # compare vs baseline

# Centralized log management:
# ELK Stack (Elasticsearch, Logstash, Kibana)
# Splunk
# Graylog
```

---

## Summary

| Evidence | How attackers remove it | How defenders preserve it |
|----------|------------------------|--------------------------|
| Auth logs | `> /var/log/auth.log` | rsyslog remote forwarding |
| Event logs | `wevtutil cl Security` | SIEM forwarding |
| Bash history | `HISTFILE=/dev/null` | auditd syscall logging |
| File timestamps | `touch -r` | Forensic image before investigation |
| Network activity | Can't — external logs | Firewall/router logs out of attacker reach |

---

## Exercises

1. Enable rsyslog forwarding on a Linux VM to a second VM — verify logs arrive even if the source is cleared
2. Research auditd — configure it to log all command execution and verify it can't be easily disabled
3. Build a GoShield check that detects when a log file shrinks in size unexpectedly
4. Research Event ID 1102 in Windows — what does it mean and how does it aid forensics?
