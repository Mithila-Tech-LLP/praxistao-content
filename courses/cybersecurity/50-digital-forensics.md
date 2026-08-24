# Chapter 50: Digital Forensics — Following the Evidence

*Digital forensics is the science of finding, preserving, and analyzing digital evidence. It's used in incident response to understand what happened, in legal proceedings to prove it, and in threat hunting to find hidden attacks.*

---

## Forensics Principles

```
Core Principles:
1. Preserve Evidence First
   - Never work on original — work on forensic copies
   - Hash everything to prove integrity (MD5, SHA-256)
   - Chain of custody documentation

2. Order of Volatility (collect most volatile first)
   CPU registers/cache      ← most volatile (lost on reboot)
   RAM / Running processes
   Network connections
   Disk contents
   Logs
   Backup media            ← least volatile

3. Don't Alter Evidence
   - Mount disk images read-only
   - Document every action taken
   - Use write blockers for physical drives
```

---

## Disk Forensics

```bash
# Create a forensic image (bit-for-bit copy)
# Physical device → image file
dd if=/dev/sda of=/external/disk.img bs=4M status=progress

# Better: ewfacquire (EnCase format with hash verification)
ewfacquire /dev/sda -c best -S 4G -u

# Verify hash
md5sum /dev/sda
md5sum disk.img
# They must match!

# Mount image read-only
losetup -r -f disk.img    # attach to loop device
mount -r /dev/loop0 /mnt/forensics -o ro,noexec

# Autopsy — GUI forensic analysis tool
autopsy   # web interface at http://localhost:9999
```

---

## File System Artifacts

```bash
# Timeline analysis — when were files modified/accessed?
fls -r -m "/" disk.img | sort -k 3 > timeline.txt
mactime -b timeline.txt -d > readable_timeline.txt

# Recover deleted files
foremost -i disk.img -o recovered/   # carves by file header
photorec disk.img                    # GUI file carver

# Find files modified in last 24h
find /mnt/forensics -mtime -1 -type f 2>/dev/null

# The Sleuth Kit — file system analysis
fsstat disk.img     # file system info
istat disk.img 12   # inode info
icat disk.img 12    # read file by inode

# Look at unallocated space (deleted file content)
blkcat disk.img BLOCK_NUMBER
```

---

## Memory Forensics

```bash
# Acquire memory image
# Linux
sudo avml /mnt/usb/memory.lime

# Windows (WinPmem or DumpIt)
winpmem.exe memory.raw

# Analyze with Volatility 3
vol -f memory.raw imageinfo    # identify OS version

# Process analysis
vol -f memory.raw windows.pslist    # list all processes
vol -f memory.raw windows.pstree   # show parent-child relationships
vol -f memory.raw windows.cmdline  # command line for each process

# Network artifacts
vol -f memory.raw windows.netstat  # network connections at time of capture
vol -f memory.raw windows.netscan

# Find malware indicators
vol -f memory.raw windows.malfind  # suspicious memory regions
vol -f memory.raw windows.dlllist --pid 1234  # DLLs loaded by process

# Dump suspicious process
vol -f memory.raw windows.dumpfiles --pid 1234
```

---

## Log Analysis

```bash
# Windows Event Log analysis
# Key Event IDs:
# 4624 - Successful logon
# 4625 - Failed logon
# 4648 - Explicit credential logon (runas)
# 4688 - Process creation (with command line if auditing enabled)
# 4698 - Scheduled task created
# 4720 - User account created
# 7045 - New service installed
# 1102 - Audit log cleared (!)

# Parse with chainsaw (fast Windows event log analysis)
chainsaw hunt /path/to/evtx/ --rules rules/ --sigma sigma/ --mapping mapping.yml

# PowerShell
Get-WinEvent -LogName Security | Where-Object {$_.Id -eq 4625} | 
    Select-Object TimeCreated, Message | Export-Csv failed_logins.csv

# Linux auth log analysis
grep "Accepted " /var/log/auth.log | awk '{print $1,$2,$3,$9,$11}' | sort | uniq -c | sort -rn
grep "Failed " /var/log/auth.log | awk '{print $11}' | sort | uniq -c | sort -rn | head -20
```

---

## Network Forensics

```bash
# Analyze PCAP captures
tcpdump -r capture.pcap -nn
wireshark capture.pcap

# Extract files from PCAP
tcpflow -r capture.pcap -o extracted/
networkMiner capture.pcap  # GUI, extracts files, creds

# DNS analysis
tshark -r capture.pcap -Y "dns" -T fields -e dns.qry.name | sort | uniq -c | sort -rn

# Find exfiltration
tshark -r capture.pcap -Y "tcp.len > 1000 and ip.dst != 192.168.1.0/24"

# HTTP POST data
tshark -r capture.pcap -Y "http.request.method == POST" \
    -T fields -e ip.src -e http.request.uri -e urlencoded-form
```

---

## Go: Log Analyzer

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "regexp"
    "sort"
    "strings"
    "time"
)

type AuthEvent struct {
    Time     time.Time
    Type     string
    User     string
    SourceIP string
}

func analyzeAuthLog(path string) {
    f, err := os.Open(path)
    if err != nil {
        fmt.Printf("Cannot open %s: %v\n", path, err)
        return
    }
    defer f.Close()
    
    // Patterns
    successRe := regexp.MustCompile(`(\w+\s+\d+\s+\d+:\d+:\d+).*Accepted.*for (\w+) from ([\d.]+)`)
    failRe := regexp.MustCompile(`(\w+\s+\d+\s+\d+:\d+:\d+).*Failed.*for.*?(\w+) from ([\d.]+)`)
    
    failsByIP := make(map[string]int)
    failsByUser := make(map[string]int)
    successLogins := []string{}
    
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text()
        
        if m := failRe.FindStringSubmatch(line); len(m) > 0 {
            failsByIP[m[3]]++
            failsByUser[m[2]]++
        }
        
        if m := successRe.FindStringSubmatch(line); len(m) > 0 {
            successLogins = append(successLogins,
                fmt.Sprintf("%s | user: %s | from: %s", m[1], m[2], m[3]))
        }
    }
    
    fmt.Println("=== AUTHENTICATION LOG ANALYSIS ===\n")
    
    // Top failing IPs (potential brute force)
    type kv struct{ K string; V int }
    var pairs []kv
    for k, v := range failsByIP {
        pairs = append(pairs, kv{k, v})
    }
    sort.Slice(pairs, func(i, j int) bool { return pairs[i].V > pairs[j].V })
    
    fmt.Println("Top 10 IPs with failed logins (potential brute force):")
    for i, p := range pairs {
        if i >= 10 {
            break
        }
        alert := ""
        if p.V > 20 {
            alert = " [!!! BRUTE FORCE]"
        }
        fmt.Printf("  %4d attempts: %s%s\n", p.V, p.K, alert)
    }
    
    fmt.Println("\nTop targeted usernames:")
    pairs = nil
    for k, v := range failsByUser {
        pairs = append(pairs, kv{k, v})
    }
    sort.Slice(pairs, func(i, j int) bool { return pairs[i].V > pairs[j].V })
    for i, p := range pairs {
        if i >= 5 {
            break
        }
        fmt.Printf("  %4d attempts: %s\n", p.V, p.K)
    }
    
    fmt.Printf("\nSuccessful logins (%d total):\n", len(successLogins))
    for _, l := range successLogins {
        // Flag logins from high-fail IPs
        for _, ip := range pairs {
            if strings.Contains(l, ip.K) && ip.V > 10 {
                fmt.Printf("  [!] AFTER BRUTE FORCE: %s\n", l)
            }
        }
        fmt.Printf("  %s\n", l)
    }
}

func main() {
    analyzeAuthLog("/var/log/auth.log")
}
```

---

## Chain of Custody

```
DIGITAL EVIDENCE CHAIN OF CUSTODY FORM

Case Number: IR-2025-042
Evidence Item: Laptop SSD — Dell XPS 15
Collected by: J. Smith
Collection time: 2025-03-15 14:32 UTC
Collection method: Forensic imaging with write blocker
SHA-256 of original: a4b9c2d1e8f3...
SHA-256 of copy: a4b9c2d1e8f3... (MATCHES)

Transfer log:
2025-03-15 14:32 - Collected by J. Smith at Scene
2025-03-15 16:00 - Transferred to Forensics Lab, received by M. Jones
2025-03-16 09:00 - Analysis started by M. Jones
2025-03-20 17:00 - Analysis complete, stored in Evidence Room #4
```

---

## Summary

| Forensic Area | What you find | Key Tools |
|---------------|---------------|-----------|
| Disk forensics | Files, deleted data, timeline | Autopsy, The Sleuth Kit, dd |
| Memory forensics | Running processes, network state, keys | Volatility3, AVML |
| Log forensics | Login events, execution, tampering | Chainsaw, grep, SIEM |
| Network forensics | Traffic, credentials, exfil | Wireshark, tcpflow |

---

## Exercises

1. Download a forensic challenge from CyberDefenders — analyze a memory dump with Volatility3
2. Create a forensic image of a USB drive, hash it, and recover deleted files with foremost
3. Analyze auth.log from a honeypot — build a complete picture of who attacked it and when
4. Write the Go log analyzer above and test it on your own system's auth log
