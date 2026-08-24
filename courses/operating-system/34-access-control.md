# Chapter 34: Access Control — Who Can Access What

> **"An operating system is ultimately a trust authority. It must decide: can this process access this file? Can this user run this program? Can this network packet reach this socket? The policies that govern these decisions — DAC, MAC, RBAC, ACLs — shape the security boundary between users, between processes, and between systems."**

---

## Table of Contents

1. [The Access Control Problem](#1-the-access-control-problem)
2. [DAC — Discretionary Access Control](#2-dac--discretionary-access-control)
3. [Unix Permission Model](#3-unix-permission-model)
4. [ACLs — Fine-Grained Permissions](#4-acls--fine-grained-permissions)
5. [MAC — Mandatory Access Control](#5-mac--mandatory-access-control)
6. [SELinux — Security-Enhanced Linux](#6-selinux--security-enhanced-linux)
7. [AppArmor — Profile-Based MAC](#7-apparmor--profile-based-mac)
8. [RBAC — Role-Based Access Control](#8-rbac--role-based-access-control)
9. [Capabilities — Fine-Grained Privileges](#9-capabilities--fine-grained-privileges)
10. [Seccomp — Syscall Filtering](#10-seccomp--syscall-filtering)
11. [Summary](#summary)

---

## 1. The Access Control Problem

Every OS must answer:

```
When process P tries to perform operation O on resource R:
  Should the OS allow or deny this?

Example decisions:
  Can alice's shell (PID 2345) read /etc/shadow? → NO (root-only)
  Can apache (UID 33) write to /var/www/html? → YES (apache owns it)
  Can a JavaScript engine read /etc/passwd? → NO (sandboxed)
  Can a database process bind to port 80? → NO (only root by default)
```

**Three questions every access control system must answer:**
1. **Authentication:** Who are you? (Verify identity)
2. **Authorization:** What are you allowed to do? (Check permissions)
3. **Audit:** What did you do? (Logging)

---

## 2. DAC — Discretionary Access Control

**Discretionary Access Control (DAC):** The owner of a resource decides who can access it. The OS enforces the owner's decision.

**"Discretionary":** The resource owner has discretion (freedom) to set permissions as they like. The OS doesn't impose any policy beyond "only the owner can change permissions."

**DAC is what Unix and Windows use by default.**

**Strengths:**
- Simple and flexible
- Users have full control over their own files

**Weaknesses:**
- Malware running as the user has ALL the user's permissions
- No protection against "confused deputy" attacks
- Can't prevent a process from doing things the user is allowed to do

```
Attack scenario (DAC weakness):
  alice's account is tricked into running malware
  Malware runs as alice (UID 1000)
  Malware has READ access to all alice's files (SSH keys, documents, passwords)
  Nothing stops it — alice's files allow alice to read them
  DAC: "alice authorized to read these files" — malware IS alice
```

---

## 3. Unix Permission Model

Unix uses a simple 9-bit permission model:

```
-rw-r--r-- 1 alice users 2048 Jul 1 file.txt

Type: - (regular file)
Owner permissions: rw- (read + write)
Group permissions: r-- (read only)
Others permissions: r-- (read only)

Owner: alice (UID 1000)
Group: users (GID 100)
```

**Permission bits:**
```
Permission check order:
  1. Is process UID == file owner UID? → apply owner permissions
  2. Is process GID (or supplementary GIDs) == file group GID? → apply group permissions
  3. Otherwise → apply others permissions
  
Note: these are OR gates — if you ARE the owner, only owner bits apply (not group or others)
If alice is the owner: alice's access is always determined by owner bits only
```

**Special bits:**

**SETUID bit (s on owner execute position):**
```
-rwsr-xr-x 1 root root /usr/bin/passwd

When any user executes passwd, the process runs with root's UID (EUID=0), not the caller's UID.
This allows passwd to write /etc/shadow (root-only) on the user's behalf.
The executable is trusted to do only what it should.

Dangerous: if a setuid program has a bug, attackers exploit it to get root.
```

**SETGID bit:**
Same but sets Effective GID instead of UID.
On directories: new files inherit the directory's group (useful for shared directories).

**Sticky bit (t on others execute):**
```
drwxrwxrwt 1 root root /tmp

In /tmp, anyone can create files. But without sticky bit, anyone can delete anyone else's files.
With sticky bit: you can only delete/rename your OWN files in this directory.
(Used on /tmp, /var/tmp to prevent users from deleting each other's temp files)
```

**umask:**
The umask masks out permissions when creating files:
```bash
umask 022
# New files: 0666 & ~022 = 0644 (rw-r--r--)
# New dirs:  0777 & ~022 = 0755 (rwxr-xr-x)

umask 027
# New files: 0666 & ~027 = 0640 (rw-r-----)
# New dirs:  0777 & ~027 = 0750 (rwxr-x---)
```

---

## 4. ACLs — Fine-Grained Permissions

Unix permissions only allow owner + group + others. What if you need:
- "Alice AND bob can read, but not other users"
- "The 'developers' group can read-write, 'testers' group can read-only"

**POSIX ACLs** extend Unix permissions with per-user and per-group entries:

```bash
# Set ACL on a file:
setfacl -m u:bob:rw- /home/alice/shared.txt    # bob can read+write
setfacl -m g:developers:r-- /home/alice/code/   # developers group: read only

# View ACL:
getfacl /home/alice/shared.txt
# file: /home/alice/shared.txt
# owner: alice
# group: alice
# user::rw-          ← alice (owner): rw
# user:bob:rw-       ← bob specifically: rw
# group::r--         ← alice's group: read only
# mask::rw-          ← maximum allowed for named entries
# other::r--         ← everyone else: read only

# The "mask" is the maximum permission for any named entry (ACL mask)
# Actual permission = ACL entry AND mask
```

**Windows NTFS ACL** (see Chapter 26) — even more powerful: per-user/group per-operation allow/deny entries.

---

## 5. MAC — Mandatory Access Control

**Mandatory Access Control (MAC):** The OS enforces a system-wide security policy that the resource owner CANNOT override.

**"Mandatory":** Permissions are IMPOSED by the system administrator (or security policy). Even the root user can't grant permissions that violate the MAC policy.

**MAC is used for:**
- Military systems (top secret / secret / unclassified labels)
- High-security servers (limit damage from compromised root)
- Mobile devices (Android uses SELinux; apps can't access each other's data even if running as same user)

**Classic MAC models:**

**Bell-LaPadula (confidentiality):**
```
Subjects (processes) and objects (files) have security levels:
  TOP_SECRET > SECRET > CONFIDENTIAL > UNCLASSIFIED

Rules:
  No Read Up:   Process can't read above its level (can't leak classified data to uncleared process)
  No Write Down: Process can't write below its level (can't leak classified data to lower-level)
```

**Biba (integrity):**
```
Rules:
  No Write Up:   Low-integrity process can't corrupt high-integrity data
  No Read Down:  High-integrity process can't be contaminated by low-integrity data
```

---

## 6. SELinux — Security-Enhanced Linux

**SELinux (Security-Enhanced Linux)** is the most widely deployed MAC system, originally developed by the NSA. Used by default in Fedora, Red Hat, CentOS, Android.

**Core concept: Every process and file has a label.**

```
File label:   system_u:object_r:httpd_content_t:s0
Process label: system_u:system_r:httpd_t:s0

Format: user:role:type:level
  user:  SELinux user (separate from Unix UID)
  role:  What roles this user can assume
  type:  The most important part — what TYPE this resource is
  level: MLS/MCS sensitivity level (optional)
```

**Type Enforcement (TE):**
SELinux policy says: "process of type X CAN/CANNOT perform operation O on resource of type Y"

```bash
# Example policy rule:
allow httpd_t httpd_content_t:file { read getattr };
#     ^          ^                   ^
# apache type  web content type    allowed operations

# This means:
# Apache (httpd_t) can READ files labeled httpd_content_t
# Apache CANNOT access files labeled user_home_t, shadow_t, etc.

# Even if a hacker compromises Apache, they can only access what httpd_t is allowed to
```

**Real-world example:**
```bash
# Apache is running. Hacker exploits a PHP RCE vulnerability.
# Without SELinux: hacker can do anything apache can — read /etc/passwd, etc.
# With SELinux: hacker is still constrained to httpd_t policy
#   Can read web files → yes (httpd_content_t allowed)
#   Can read /etc/passwd → NO (denied — shadow_t, passwd_file_t not in httpd_t policy)
#   Can write to /var/log/messages → NO
#   Can fork() and exec() a shell → NO (shell_exec_t not in httpd_t policy)
```

**Checking SELinux:**
```bash
# SELinux status:
getenforce
# Enforcing (active) / Permissive (log but don't block) / Disabled

# View labels:
ls -Z /var/www/html/
# -rw-r--r--. root root system_u:object_r:httpd_content_t:s0 index.html

# View process labels:
ps -eZ | grep httpd
# system_u:system_r:httpd_t:s0 ... httpd

# Check audit log for denials:
ausearch -m avc -ts recent
# type=AVC msg=audit(...): avc: denied { read } for pid=1234 comm="httpd"
#   path="/etc/shadow" scontext=httpd_t tcontext=shadow_t tclass=file permissive=0
```

---

## 7. AppArmor — Profile-Based MAC

**AppArmor** is an easier-to-use alternative to SELinux. Uses profiles defined by executable path.

**Example profile for nginx:**
```
#include <tunables/global>

/usr/sbin/nginx {
    #include <abstractions/base>
    
    capability net_bind_service,  # bind to port 80
    
    /var/www/html/ r,             # read web content directory
    /var/www/html/** r,           # read all files in it
    /var/log/nginx/ rw,           # read+write log directory
    /var/log/nginx/*.log rw,      # read+write log files
    /etc/nginx/ r,                # read config directory
    /etc/nginx/** r,              # read all config files
    
    /proc/*/fd/ r,                # process monitoring
    
    network tcp,                  # allow TCP networking
    deny /etc/shadow r,           # explicitly deny shadow file
}
```

**AppArmor vs SELinux:**
| Feature | SELinux | AppArmor |
|---------|---------|----------|
| Policy model | Type enforcement (all objects labeled) | Path-based profiles |
| Granularity | Very fine | File-path granularity |
| Complexity | Very complex | Easier to write/understand |
| Default distros | RHEL, Android, CentOS | Ubuntu, Debian, SUSE |
| Coverage | Everything labeled | Only profiled programs |

---

## 8. RBAC — Role-Based Access Control

**RBAC (Role-Based Access Control):** Permissions are assigned to roles, and users are assigned to roles. Users get permissions through their roles.

```
Traditional DAC:
  alice → [can read /etc/passwd, can write /home/alice/, can run apache...]
  bob   → [can read /etc/passwd, can write /home/bob/, ...]
  
RBAC:
  Role "webadmin" → [can restart nginx, can read/write /var/www/html/, can view logs]
  Role "dba" → [can read/write /var/lib/postgresql/, can run psql]
  
  alice ∈ {webadmin, dba}  → gets all permissions of both roles
  bob   ∈ {webadmin}        → gets webadmin permissions only
```

**Linux sudo (simplified RBAC):**
```bash
# /etc/sudoers — define what commands each user/group can run as root:
alice ALL=(ALL) /usr/bin/systemctl restart nginx
# alice can run "systemctl restart nginx" as root, nothing else

%webteam ALL=(ALL) NOPASSWD: /usr/bin/tail -f /var/log/nginx/*.log
# webteam group: run tail on nginx logs without password

# Check what you can sudo:
sudo -l
```

**Linux groups as lightweight RBAC:**
```bash
# Add user to docker group (allows running docker without sudo):
usermod -aG docker alice

# Add user to audio group (allows accessing /dev/dsp, /dev/snd):
usermod -aG audio alice

# Show group memberships:
id alice
# uid=1000(alice) gid=1000(alice) groups=1000(alice),100(users),999(docker),29(audio)
```

---

## 9. Capabilities — Fine-Grained Privileges

**The problem with root:**
Traditionally, processes are either unprivileged (uid=0's permissions are needed) or root (all privileges). No middle ground.

Running a web server as root just to bind to port 80 is risky — if it's compromised, the attacker has full root.

**Linux capabilities** divide root's power into ~41 independent capabilities:

```
CAP_BIND_SERVICE:  Bind to ports < 1024
CAP_KILL:          Send signals to any process
CAP_NET_ADMIN:     Network configuration (iptables, interface setup)
CAP_SYS_PTRACE:    ptrace any process (debugger access)
CAP_SYS_ADMIN:     Broad administrative powers (mount, sethostname, etc.)
CAP_SETUID:        setuid() to change UID
CAP_CHROOT:        chroot() to change root
CAP_SYS_RAWIO:     Raw I/O to devices (/dev/mem, iopl)
CAP_SYS_TIME:      Set system time
CAP_SYS_MODULE:    Load/unload kernel modules
CAP_DAC_OVERRIDE:  Bypass file permission checks
CAP_DAC_READ_SEARCH: Bypass file read permission and directory search
```

**Usage:**
```bash
# Give nginx just the capability to bind to port 80 (not full root):
setcap cap_net_bind_service=+ep /usr/sbin/nginx

# Now nginx can start and bind to port 80 without being root
# If compromised: attacker has no other root capabilities

# View process capabilities:
cat /proc/$(pgrep nginx)/status | grep Cap
# CapPrm: 0000000000000400  ← bit 10 = CAP_NET_BIND_SERVICE
# CapEff: 0000000000000400

# Docker uses capabilities:
docker run --cap-drop ALL --cap-add NET_BIND_SERVICE nginx
# Start with NO capabilities, add only what's needed
```

---

## 10. Seccomp — Syscall Filtering

**Seccomp (Secure Computing Mode)** filters which system calls a process can make.

**Why it matters:**
Even with capabilities and MAC, a compromised process can still call dangerous syscalls. Seccomp restricts what syscalls are even available.

**Strict mode (original seccomp):**
Only allows `read`, `write`, `_exit`, `sigreturn`. Anything else → SIGKILL. Used for sandboxed computation.

**seccomp-BPF (modern, flexible):**
A BPF (Berkeley Packet Filter) program filters each syscall:
```c
// Example: filter to allow read, write, exit, mmap, but deny execve:
struct sock_filter filter[] = {
    BPF_STMT(BPF_LD + BPF_W + BPF_ABS, offsetof(struct seccomp_data, nr)),
    BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_read,   5, 0),
    BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_write,  4, 0),
    BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_exit,   3, 0),
    BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_mmap,   2, 0),
    BPF_JUMP(BPF_JMP + BPF_JEQ + BPF_K, __NR_execve, 0, 1),
    BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),   // kill process
    BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_ALLOW),  // allow syscall
    BPF_STMT(BPF_RET + BPF_K, SECCOMP_RET_KILL),   // kill (execve)
};
```

**Practical usage:**
```bash
# Docker applies seccomp by default — restricts ~44 dangerous syscalls:
docker run --security-opt seccomp=default nginx

# View Docker's default seccomp policy:
docker info --format '{{.SecurityOptions}}'
# name=seccomp,profile=default

# Chrome, Firefox, VSCode all use seccomp to sandbox their renderer processes
# Android uses seccomp for every app process
```

---

## Summary

| Concept | Description |
|---------|------------|
| DAC | Owner controls access; flexible but attacker inherits user's permissions |
| MAC | System-enforced policy; can't be overridden even by owner |
| Unix permissions | 9 bits: owner rwx, group rwx, others rwx + SETUID/SETGID/sticky |
| SETUID | Executable runs with owner's UID; allows controlled privilege escalation |
| umask | Masks out default permissions when creating new files |
| POSIX ACL | Per-user and per-group permissions beyond basic Unix model |
| SELinux | Type enforcement MAC; every process/file has a label; policy controls transitions |
| AppArmor | Path-based MAC profiles; easier than SELinux; default on Ubuntu |
| RBAC | Permissions assigned to roles; users granted roles |
| Linux groups | Lightweight RBAC: add user to group to grant group permissions |
| Capabilities | Divide root power into ~41 independent capabilities |
| Seccomp | Filter which syscalls a process can make; defense-in-depth for sandboxes |
| sudo | Allow non-root users to run specific commands as root |
