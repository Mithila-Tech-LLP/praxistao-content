# Chapter 07: File Systems, Permissions, and Users

*File permissions are the primary access control mechanism on Linux. Misconfigurations — a world-writable cron script, an SUID binary, a readable private key — are among the most common privilege escalation vectors.*

---

## The Linux File System Hierarchy

```
/                   Root of everything
├── bin/            Essential binaries (ls, cp, mv)
├── sbin/           System binaries (root only: iptables, fdisk)
├── etc/            Configuration files
│   ├── passwd      User accounts (world-readable)
│   ├── shadow      Password hashes (root only)
│   ├── sudoers     Sudo configuration
│   └── crontab     System cron jobs
├── home/           User home directories
│   └── alice/      Alice's files
├── root/           Root user's home
├── var/            Variable data (logs, databases)
│   ├── log/        System logs
│   └── www/        Web server files
├── tmp/            Temporary files (world-writable, no persist)
├── usr/            User programs
│   ├── bin/        User utilities
│   └── share/      Shared data
├── proc/           Virtual: kernel/process info
├── sys/            Virtual: hardware/kernel params
├── dev/            Device files
└── lib/            Libraries
```

### Security-Critical Paths

| Path | Risk |
|------|------|
| `/etc/passwd` | User list — world-readable but no hashes |
| `/etc/shadow` | Password hashes — should be root:shadow only |
| `/etc/sudoers` | Sudo rules — misconfiguration = root |
| `/tmp`, `/var/tmp` | World-writable — malware staging area |
| `/home/*/.ssh/` | Private keys — should be 700 |
| `/var/www/` | Web content — writable = webshell |
| `/root/` | Root's home — should be 700 |

---

## File Permissions

Every file has three permission sets: **owner**, **group**, **others**.

```bash
ls -la /etc/passwd
-rw-r--r-- 1 root root 2847 Jan 1 12:00 /etc/passwd
↑↑↑↑↑↑↑↑↑↑
│└─┬─┘└┬┘└┬┘
│  │   │  └── others: r-- (read only)
│  │   └───── group:  r-- (read only)
│  └───────── owner:  rw- (read + write)
└──────────── file type: - (regular file)
```

**File types:**
- `-` regular file
- `d` directory
- `l` symbolic link
- `c` character device
- `b` block device
- `s` socket
- `p` named pipe

### Permission Bits

| Symbol | Octal | Meaning (file) | Meaning (directory) |
|--------|-------|---------------|---------------------|
| `r` | 4 | Read file | List directory |
| `w` | 2 | Write file | Create/delete files |
| `x` | 1 | Execute file | Enter directory |
| `-` | 0 | No permission | No permission |

```bash
chmod 644 file    # rw-r--r--  (owner rw, group r, others r)
chmod 755 dir/    # rwxr-xr-x  (owner rwx, group rx, others rx)
chmod 600 .ssh/id_rsa  # rw-------  (only owner can read)
chmod 700 .ssh/   # rwx------  (only owner can enter)
chmod 777 /tmp/   # rwxrwxrwx  (everyone can do anything — dangerous)
```

---

## Special Permission Bits

### SUID (Set User ID) — bit 4000

When a file with SUID is executed, it runs as the **file owner**, not the person running it.

```bash
ls -la /usr/bin/passwd
-rwsr-xr-x 1 root root ... /usr/bin/passwd
    ^ 's' = SUID set

# /usr/bin/passwd is owned by root with SUID
# When any user runs passwd, it runs as root
# (needs root to write /etc/shadow)
```

**Finding SUID binaries:**
```bash
find / -perm -4000 -type f 2>/dev/null
```

**Dangerous SUID binaries** (check GTFOBins):
```bash
# If /usr/bin/find has SUID (root)
find / -exec /bin/bash -p \; -quit   # root shell!

# If /usr/bin/vim has SUID
vim -c ':!/bin/bash'                 # root shell!
```

### SGID (Set Group ID) — bit 2000

Runs as file's group. Also: files created in an SGID directory inherit the directory's group.

```bash
ls -la /usr/bin/write
-rwxr-sr-x 1 root tty ... /usr/bin/write
       ^ 's' in group position = SGID
```

### Sticky Bit — bit 1000

On directories: only file owner can delete their own files (even if others have write).

```bash
ls -la /tmp
drwxrwxrwt 20 root root ...
         ^ 't' = sticky bit

# Even though /tmp is world-writable
# You can only delete YOUR OWN files
```

---

## Users and Groups

### /etc/passwd Format

```
username:password:UID:GID:comment:home:shell
root:x:0:0:root:/root:/bin/bash
alice:x:1000:1000:Alice:/home/alice:/bin/bash
www-data:x:33:33:www-data:/var/www:/usr/sbin/nologin
```

- **UID 0** = root (any user with UID 0 has root powers — not just "root")
- **x** in password field = hash is in /etc/shadow
- `/usr/sbin/nologin` shell = cannot login interactively (service account)

### /etc/shadow Format

```
username:$hash_type$salt$hash:last_changed:min:max:warn:inactive:expire
alice:$6$salt$longhash...:19000:0:99999:7:::
```

Hash prefixes:
- `$1$` = MD5 (weak, crackable fast)
- `$5$` = SHA-256
- `$6$` = SHA-512 (strong)
- `$y$` = yescrypt (strong, new)

**Cracking with hashcat:**
```bash
hashcat -m 1800 hashes.txt /usr/share/wordlists/rockyou.txt  # SHA-512
hashcat -m 500 hashes.txt /usr/share/wordlists/rockyou.txt   # MD5
```

### Groups

```bash
cat /etc/group
# group:password:GID:members
sudo:x:27:alice,bob     # members can use sudo
docker:x:999:alice      # members can run Docker (= root equivalent!)
www-data:x:33:          # web server group
```

**High-value groups:**
- `sudo` / `wheel` — can run sudo
- `docker` — can mount host filesystem = root
- `shadow` — can read /etc/shadow
- `disk` — raw disk access = read any file

---

## Access Control Lists (ACLs)

Standard permissions are coarse (owner/group/others). ACLs add fine-grained control.

```bash
# View ACLs
getfacl /var/www/html

# Set ACL: give alice read-only on a root-owned file
setfacl -m u:alice:r /root/secret.txt

# Remove ACL
setfacl -x u:alice /root/secret.txt
```

ACLs are often overlooked — check them when standard permissions seem correct but access is unexpected.

---

## Capabilities

Linux capabilities divide root's powers into smaller pieces. Programs can have capabilities without full root.

```bash
# View capabilities
getcap /usr/bin/ping
# /usr/bin/ping = cap_net_raw+ep

# Dangerous capabilities
cap_setuid      # change UID — can become root
cap_net_raw     # raw sockets — can sniff
cap_dac_override # bypass file permissions
cap_sys_ptrace  # ptrace any process — can inject

# Example: python3 with cap_setuid
python3 -c 'import os; os.setuid(0); os.system("/bin/bash")'
# If python3 has cap_setuid+ep = instant root
```

---

## Symlinks and Hard Links

```bash
# Hard link: two names, same inode
ln file1.txt hardlink.txt
# Both point to same data; deleting one doesn't delete data

# Soft/symbolic link
ln -s /etc/passwd /tmp/passwd_link
# /tmp/passwd_link → /etc/passwd
```

**Security risks:**
- **Symlink attacks:** Create symlink in world-writable directory pointing to sensitive file. If root program follows symlink and writes to it...
```bash
# Classic /tmp race condition attack
ln -sf /etc/shadow /tmp/logfile
# If root program writes logs to /tmp/logfile...
# ... it's writing to /etc/shadow
```

Modern kernels mitigate this with `fs.protected_symlinks=1`.

---

## File Attributes

Beyond permissions, files have attributes:

```bash
lsattr /etc/passwd          # list attributes
chattr +i /etc/passwd       # immutable — even root can't modify/delete!
chattr -i /etc/passwd       # remove immutable
chattr +a logfile.txt       # append-only — can add but not overwrite

# Attackers set +i on backdoors to prevent removal
# Defenders set +i on critical configs
```

---

## Umask — Default Permissions

`umask` sets which bits are removed from default permissions when creating files.

```bash
umask 022   # most common default
# New files: 666 - 022 = 644 (rw-r--r--)
# New dirs:  777 - 022 = 755 (rwxr-xr-x)

umask 077   # more secure (no group/other access)
# New files: 666 - 077 = 600 (rw-------)
# New dirs:  777 - 077 = 700 (rwx------)
```

---

## Auditing Permissions — Security Checklist

```bash
# Files writable by everyone (risky)
find / -perm -002 -type f 2>/dev/null | grep -v /proc | grep -v /sys

# Directories writable by everyone
find / -perm -002 -type d 2>/dev/null | grep -v /proc | grep -v /sys

# All SUID binaries
find / -perm -4000 -type f 2>/dev/null

# All SGID binaries
find / -perm -2000 -type f 2>/dev/null

# Files with no owner (leaked from deleted user)
find / -nouser -type f 2>/dev/null

# SSH key permissions (should be 600)
find /home -name "id_rsa" -not -perm 600 2>/dev/null

# World-readable sensitive files
ls -la /etc/shadow /etc/sudoers /etc/ssh/sshd_config
```

---

## Summary

| Concept | Attack relevance |
|---------|----------------|
| SUID binaries | Common privesc — run as owner (root) |
| World-writable dirs | Malware staging, symlink attacks |
| /etc/shadow | Offline password cracking |
| Docker group | Equivalent to root access |
| Capabilities | Selective root powers without full root |
| Symlinks | Race condition attacks in /tmp |

---

## Exercises

1. List all SUID binaries on your system. Research each on GTFOBins — which can be exploited?
2. Check if `/tmp` has the sticky bit set. Why is this important?
3. What UID does `www-data` have on your system? What groups?
4. Find all world-writable files in `/etc`. What could go wrong if any of these are modified?
5. Write a Go program that takes a file path and prints its permissions, owner, group, and special bits (SUID/SGID/sticky).
