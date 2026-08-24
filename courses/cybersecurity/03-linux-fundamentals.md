# Chapter 03: Linux — The Hacker's Operating System

*Linux runs 96% of the world's servers, all Android phones, most cloud infrastructure, and most IoT devices. It's the dominant target AND the dominant tool platform. You need to be fluent in Linux.*

---

## Why Linux for Security?

1. **It's everywhere:** Most targets run Linux (web servers, cloud instances, databases)
2. **All tools are on Linux:** Metasploit, Nmap, Wireshark, Burp Suite, Ghidra — the entire security tooling ecosystem
3. **You can see everything:** Open source, no hidden behavior
4. **File-based everything:** Devices, network sockets, processes — all accessible as files
5. **Kali Linux:** A Linux distribution with 600+ pre-installed security tools

---

## Linux Distributions

A "distribution" (distro) is a packaged version of Linux with specific tools and defaults.

| Distro | Use Case |
|--------|----------|
| **Kali Linux** | Penetration testing — has 600+ security tools |
| **ParrotOS** | Security, privacy, anonymity |
| **Ubuntu** | General use, servers |
| **Debian** | Stable, production servers |
| **CentOS/RHEL** | Enterprise servers |
| **Alpine** | Minimal, containers |
| **Arch Linux** | Advanced users who build their own |

**Recommendation:** Use **Kali Linux** in a VM (Virtual Machine) for this course. It comes with everything you need.

---

## The Linux File System Hierarchy

Linux organizes everything in a single tree starting at `/` (root):

```
/
├── bin/     → Essential user commands (ls, cat, cp)
├── boot/    → Boot files (kernel, grub)
├── dev/     → Device files (hard drives, terminals)
├── etc/     → Configuration files
│   ├── passwd    → User account info
│   ├── shadow    → Hashed passwords (root only)
│   ├── hosts     → Local DNS overrides
│   └── cron.d/   → Scheduled tasks
├── home/    → User home directories
│   └── aditya/   → Your files
├── lib/     → Shared libraries
├── mnt/     → Temporarily mounted filesystems
├── opt/     → Optional/third-party software
├── proc/    → Virtual filesystem — running processes
│   └── 1234/     → Process 1234's info
├── root/    → Root user's home directory
├── sbin/    → System administration commands
├── tmp/     → Temporary files (deleted on reboot)
├── usr/     → User programs
│   ├── bin/      → Most user commands
│   └── share/    → Documentation, shared data
└── var/     → Variable data (logs, databases)
    └── log/      → System logs
```

**Security-critical files:**
- `/etc/passwd` — User list (readable by all)
- `/etc/shadow` — Password hashes (root only — if you can read this, you can crack passwords)
- `/var/log/auth.log` — Authentication events (login attempts)
- `/var/log/syslog` — System events
- `/proc/[pid]/` — Running process information
- `/tmp/` — Writable by everyone — attacker favorite for dropping malware

---

## Linux Users and Permissions

### Users and Groups

Every file has an owner (user) and a group. Every process runs as a user.

- **root** (UID 0): The superuser. Has access to everything. This is the target of privilege escalation attacks.
- **Normal users** (UID 1000+): Limited access
- **System users** (UID 1-999): Services (www-data for Apache, mysql for MySQL)

Commands:
```bash
whoami              # Who am I?
id                  # My UID, GID, and groups
id username         # Another user's info
cat /etc/passwd     # List all users
```

### File Permissions

Every file has three permission sets: owner, group, others.

```
-rwxr-xr-- 1 root www-data 1234 Jan 1 12:00 script.sh
│││││││││
│││││││└└─ Others: r-x (read, no write, execute)  
│││││└└───  Group: r-x (read, no write, execute)
│││└└───── Owner: rwx (read, write, execute)
││└──────── File type: - (regular), d (directory), l (symlink)
```

Permission bits in octal:
- `r` = 4, `w` = 2, `x` = 1
- `rwx` = 7, `rw-` = 6, `r-x` = 5, `r--` = 4

So `rwxr-xr--` = 754 in octal.

```bash
chmod 755 script.sh    # rwxr-xr-x
chmod 644 file.txt     # rw-r--r--
chmod +x script.sh     # Add execute for all
chown root:www-data file.txt  # Change owner and group
```

### Special Permissions — SUID/SGID/Sticky Bit

**SUID (Set User ID):** When an executable with SUID is run, it runs as the file's owner, not the current user.

```bash
ls -la /usr/bin/sudo
-rwsr-xr-x 1 root root ... /usr/bin/sudo
#    ^ the 's' means SUID is set
```

`sudo` has SUID set and is owned by root — so when you run `sudo`, it briefly runs as root to perform privileged operations.

**This is critical for privilege escalation:** Finding SUID binaries owned by root that can be exploited gives you root access.

```bash
# Find all SUID binaries on the system
find / -perm -4000 -type f 2>/dev/null
```

---

## Essential Linux Commands

### Navigation
```bash
pwd              # Print Working Directory (where am I?)
ls               # List files
ls -la           # List with permissions, hidden files, sizes
cd /var/log      # Change directory
cd ~             # Go to home directory
cd ..            # Go up one level
```

### File Operations
```bash
cat /etc/passwd  # Print file contents
less /var/log/syslog  # Scroll through large files
head -20 file    # First 20 lines
tail -f /var/log/auth.log  # Follow file in real-time (great for watching logs)
cp source dest   # Copy
mv source dest   # Move/rename
rm file          # Delete
rm -rf dir/      # Delete directory recursively (dangerous!)
mkdir newdir     # Create directory
```

### Searching
```bash
grep "password" /etc/config  # Search for "password" in file
grep -r "password" /etc/     # Recursive search
grep -i "error" log.txt      # Case-insensitive
find / -name "*.conf"        # Find files by name
find / -user www-data        # Files owned by www-data
locate passwd                # Fast file search (uses index)
which nmap                   # Where is the nmap binary?
```

### Networking
```bash
ifconfig         # Network interfaces (older)
ip addr          # Network interfaces (modern)
ip route         # Routing table
netstat -tulpn   # Open ports and listening services
ss -tulpn        # Modern version of netstat
ping 8.8.8.8     # Test connectivity
traceroute google.com  # Route to destination
curl http://example.com    # HTTP request
wget http://example.com/file  # Download file
nmap -p 22,80 192.168.1.1   # Port scan
```

### Process Management
```bash
ps aux           # All running processes
ps aux | grep nginx  # Find nginx process
top              # Live process view
htop             # Better live process view
kill 1234        # Kill process by PID
kill -9 1234     # Force kill
pgrep nginx      # Find PID by name
lsof -p 1234     # Files opened by process 1234
lsof -i :80      # What process is using port 80?
```

### User Management
```bash
sudo command     # Run command as root
su - username    # Switch to user
useradd hacker   # Create user
passwd hacker    # Set password
usermod -aG sudo hacker  # Add to sudo group
```

---

## Bash Scripting Basics

Bash is the default shell in most Linux systems. Scripts automate tasks.

```bash
#!/bin/bash
# This is a comment

# Variables
TARGET="192.168.1.1"
PORT=80

# Conditional
if ping -c 1 $TARGET > /dev/null 2>&1; then
    echo "$TARGET is up"
else
    echo "$TARGET is down"
fi

# Loop over a list
for port in 22 80 443 8080; do
    if nc -z $TARGET $port 2>/dev/null; then
        echo "Port $port is OPEN"
    fi
done

# Loop with counter
for i in $(seq 1 10); do
    echo "Trying host 192.168.1.$i"
done

# Read from file
while IFS= read -r line; do
    echo "Checking: $line"
done < targets.txt

# Functions
scan_port() {
    local host=$1
    local port=$2
    nc -z -w1 $host $port 2>/dev/null
    return $?
}

scan_port 192.168.1.1 22 && echo "SSH open!"
```

---

## Linux for Security — Key Files to Know

```bash
# User accounts
cat /etc/passwd     # Format: user:x:uid:gid:info:home:shell
cat /etc/shadow     # Hashed passwords (root only)
cat /etc/group      # Group memberships

# Network configuration
cat /etc/hosts      # Local DNS overrides
cat /etc/resolv.conf  # DNS server configuration
cat /etc/iptables/rules.v4  # Firewall rules

# System logs
/var/log/auth.log   # Authentication (login attempts, sudo)
/var/log/syslog     # General system messages
/var/log/kern.log   # Kernel messages
/var/log/apache2/access.log  # Web server access log
/var/log/nginx/error.log     # Nginx errors

# Scheduled tasks (cron jobs)
crontab -l          # Current user's cron jobs
cat /etc/crontab    # System cron jobs
ls /etc/cron.d/     # Additional cron directories

# Running services
systemctl list-units --type=service --state=running
```

---

## Pipes and Redirection

Pipes connect commands: output of one becomes input of next.

```bash
# Find failed SSH logins
grep "Failed password" /var/log/auth.log | awk '{print $11}' | sort | uniq -c | sort -rn

# This does:
# grep "Failed password" → filter for failed logins
# awk '{print $11}'     → extract the IP address (11th field)
# sort                  → sort them
# uniq -c              → count unique IPs
# sort -rn             → sort by count, descending
```

Output:
```
    145 192.168.1.105
     67 10.0.0.45
     23 172.16.0.200
```
You've just found the IPs brute-forcing your SSH!

---

## Setting Up Kali Linux

1. Download Kali Linux ISO from kali.org (free)
2. Install VirtualBox or VMware (free virtualization)
3. Create a new VM, allocate 4GB RAM, 50GB disk
4. Install Kali from the ISO
5. Default credentials: `kali`/`kali`

**Or:** Use WSL2 on Windows (Windows Subsystem for Linux) for quick access.

---

## Summary

| Concept | Command | Security Use |
|---------|---------|-------------|
| File permissions | `ls -la`, `chmod` | Finding SUID bits, misconfigured files |
| User info | `id`, `whoami` | Understanding privilege context |
| Open ports | `ss -tulpn`, `netstat` | Service enumeration |
| Running processes | `ps aux`, `lsof` | Finding malicious processes |
| Log analysis | `tail -f`, `grep` | Detecting intrusions |
| Find files | `find /`, `locate` | Finding sensitive files, SUID binaries |

---

## Exercises

1. On your Linux system (or Kali VM), find all SUID binaries owned by root. How many are there?
2. Read `/etc/passwd`. What fields does each line have? What does the `x` in the password field mean?
3. Watch `/var/log/auth.log` in real time (`tail -f`). Try logging in with a wrong password from another terminal. What do you see?
4. Write a bash script that takes an IP range (like `192.168.1.1` to `192.168.1.254`) and pings each host, printing which ones are alive.
5. What processes are listening on ports on your machine? Use `ss -tulpn` and identify each service.
