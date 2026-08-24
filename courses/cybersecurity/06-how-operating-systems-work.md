# Chapter 06: How Operating Systems Work — Processes, Syscalls, Memory

*Every hack ultimately involves tricking an operating system into doing something it shouldn't. Understanding how an OS actually works — what a process is, how memory is laid out, what syscalls do — is what separates security professionals from script kiddies.*

---

## What an Operating System Does

The OS is the master program that manages everything:

1. **Process management** — creates, schedules, kills processes
2. **Memory management** — gives each process its own virtual address space
3. **File system** — organizes storage into files and directories
4. **Device management** — talks to hardware through drivers
5. **Networking** — manages network interfaces and sockets
6. **Security** — enforces permissions, users, and access control

The OS is the only program that talks directly to hardware. Every other program goes through the OS via **system calls**.

---

## Processes

A process is a running program. It's not the program itself (that's a file on disk) — it's the program in execution.

```
Program (ELF/EXE file on disk)
    ↓  OS loads it
Process (in memory, running)
  ├── PID (Process ID) — unique number
  ├── PPID (Parent PID) — who spawned it
  ├── Memory space (virtual)
  ├── File descriptors (open files, sockets)
  ├── Credentials (UID, GID)
  └── Threads (one or more execution contexts)
```

### Process Tree

Every process (except init/PID 1) has a parent:

```
PID 1 (systemd/init)
├── PID 234 (sshd)
│   └── PID 891 (sshd: user session)
│       └── PID 892 (bash)
│           └── PID 1045 (vim)
├── PID 345 (nginx)
│   ├── PID 346 (nginx worker)
│   └── PID 347 (nginx worker)
└── PID 567 (cron)
    └── PID 902 (backup.sh)
```

**Why parent-child matters for security:**
- `outlook.exe` spawning `cmd.exe` = suspicious (spearphishing)
- `nginx` spawning `bash` = suspicious (webshell)
- Process injection makes a legit process spawn malicious code with no parent change

### Process States

```
Created → Ready → Running → Waiting → Terminated
              ↑_____________↓ (context switch)
```

- **Running:** Currently executing on CPU
- **Ready:** Waiting for CPU time
- **Waiting/Blocked:** Waiting for I/O, sleep, or event

---

## Virtual Memory

Every process gets its own virtual address space — as if it has the entire memory to itself. The OS maps virtual addresses to physical RAM.

```
Virtual Address Space of a Process (64-bit Linux):
┌─────────────────────────────┐ 0xFFFFFFFFFFFFFFFF
│   Kernel space              │ (kernel mapped here, inaccessible to user)
├─────────────────────────────┤ 0x7FFFFFFFFFFFFFFF
│   Stack (grows down ↓)      │ local variables, function frames
│   ↓                         │
│                             │
│   ↑                         │
│   Heap (grows up ↑)         │ malloc/new allocations
├─────────────────────────────┤
│   BSS segment               │ uninitialized global variables
├─────────────────────────────┤
│   Data segment              │ initialized global variables
├─────────────────────────────┤
│   Text segment (read-only)  │ program code (instructions)
└─────────────────────────────┘ 0x0000000000000000
```

**Why this matters for attacks:**
- **Stack overflow** — write beyond stack frame, overwrite return address, control execution flow
- **Heap overflow** — corrupt heap metadata to hijack allocations
- **Buffer overflow** — write past end of buffer into adjacent memory
- **Use-after-free** — use freed heap memory, now controlled by attacker

---

## System Calls

User programs can't touch hardware directly. They ask the kernel through **syscalls**.

```
User space program
    |
    | calls write(fd, buf, len)
    |
    ↓ int 0x80 / syscall instruction (trap to kernel)
Kernel space
    |
    | kernel validates, executes
    |
    ↓ returns result to userspace
```

### Key Syscalls for Security

```c
// Process
fork()          // create child process
execve()        // replace current process with new program
clone()         // create thread or process
ptrace()        // trace/debug another process (used for injection!)
kill(pid, sig)  // send signal to process

// Memory
mmap()          // map memory (also used to execute shellcode)
mprotect()      // change memory permissions (make memory executable)
brk()           // change heap size

// File
open()          // open file
read(), write() // I/O
stat()          // file metadata

// Network
socket()        // create socket
bind(), listen(), accept()  // server
connect()                   // client
send(), recv()              // data transfer

// Privilege
setuid()        // change user ID (used in privilege escalation)
getuid()        // get current UID
capget/capset() // capabilities
```

**Security relevance:**
- `execve` is how every shell command executes — monitoring it detects process creation
- `ptrace` is used for debuggers AND for process injection — suspicious in production
- `mmap` + `mprotect` with PROT_EXEC = creating executable memory (shellcode staging)

### Monitoring Syscalls

```bash
# strace — trace syscalls of a process
strace ls /tmp
strace -p 1234          # attach to running process
strace -e trace=network nc google.com 80  # only network syscalls

# auditd — kernel-level syscall auditing (permanent)
auditctl -a always,exit -F arch=b64 -S execve -k exec_events
ausearch -k exec_events | grep -A 5 "bash"
```

---

## User Space vs Kernel Space

```
User Space (Ring 3):          Kernel Space (Ring 0):
  Applications                  OS Kernel
  Libraries (libc)              Device Drivers
  Shell                         Hardware access
  Your tools                    
        |                              ↑
        +--- syscall (int 0x80) -------+
```

- **Ring 0:** Full hardware access, all instructions, no restrictions
- **Ring 3:** Restricted — cannot access hardware directly, cannot execute privileged instructions

**Privilege escalation** = getting from Ring 3 (user) to Ring 0 (kernel) level control. Either through:
- Kernel exploit (CVE in the kernel itself)
- SUID binary
- Sudo misconfiguration

---

## Signals

Signals are asynchronous notifications sent to processes.

```bash
kill -9 1234    # SIGKILL — force kill (unblockable)
kill -15 1234   # SIGTERM — graceful shutdown
kill -2 1234    # SIGINT — Ctrl+C
Ctrl+Z          # SIGTSTP — suspend
```

| Signal | Number | Default | Blockable? |
|--------|--------|---------|-----------|
| SIGHUP | 1 | Terminate | Yes |
| SIGINT | 2 | Terminate | Yes |
| SIGKILL | 9 | Terminate | **No** |
| SIGSEGV | 11 | Core dump | Yes |
| SIGTERM | 15 | Terminate | Yes |
| SIGSTOP | 19 | Stop | **No** |

**Security note:** SIGKILL (9) cannot be caught, blocked, or ignored. It's the only guaranteed way to stop a process.

---

## File Descriptors

Every open file, socket, or pipe is a file descriptor (integer).

```
0 = stdin  (standard input)
1 = stdout (standard output)
2 = stderr (standard error)
3+ = opened files, sockets, pipes
```

**Inheritance:** Child processes inherit parent's file descriptors. A fork of nginx inherits its listening socket.

**Limits:**
```bash
ulimit -n          # max open file descriptors (default: 1024)
cat /proc/1234/fd/ # all open FDs of a process
lsof -p 1234       # files opened by process 1234
```

**Security:** Malware often checks which FDs are open to detect sandboxes (sandboxes may have different FD patterns).

---

## Inter-Process Communication (IPC)

Processes need to talk to each other:

| Mechanism | Use case | Security concern |
|-----------|---------|-----------------|
| Pipes (`|`) | One-directional data flow | Restricted to related processes |
| Unix sockets | Fast local IPC | File permissions apply |
| Shared memory | Fastest IPC | Race conditions |
| Signals | Notifications | Limited data |
| D-Bus | Desktop IPC | Privilege escalation via D-Bus vulns |
| /proc | Kernel → userspace | Reading /proc can expose sensitive info |

---

## The /proc Filesystem

Linux exposes kernel internals as a virtual filesystem at `/proc`.

```bash
/proc/PID/cmdline    # full command line of process PID
/proc/PID/environ    # environment variables
/proc/PID/maps       # memory map (virtual address → physical)
/proc/PID/fd/        # open file descriptors
/proc/PID/status     # state, memory usage, uid/gid
/proc/PID/net/tcp    # network connections (for this process's namespace)
/proc/self/          # refers to current process
/proc/version        # kernel version
/proc/cpuinfo        # CPU info
/proc/meminfo        # memory stats
```

```bash
# Find injected memory regions (no file backing = suspicious)
cat /proc/1234/maps | grep -v "\.so\|\.py\|heap\|stack\|vvar\|vdso" | grep "rwx\|r-x"

# Read environment variables of running process
strings /proc/1234/environ | grep -i "password\|secret\|key\|token"
```

---

## Summary

| Concept | Security relevance |
|---------|------------------|
| Processes | Monitoring process creation detects execution of tools |
| Virtual memory | Stack/heap overflows exploit this layout |
| Syscalls | Monitoring execve/mmap/ptrace detects attacks |
| Ring 0 vs Ring 3 | Privilege escalation crosses this boundary |
| /proc | Goldmine for process info and attack investigation |
| Signals | SIGKILL for incident containment |

---

## Exercises

1. Run `strace /bin/ls` and count how many syscalls it makes. Which ones involve the filesystem?
2. Read `/proc/1/maps` (systemd). Identify the text, heap, and stack segments.
3. Write a Go program that reads its own `/proc/self/status` and prints its memory usage and PID.
4. Use `lsof -p $(pgrep sshd)` to list all files and sockets sshd has open. What do you see?
5. What happens when you `strace -e trace=execve bash -c "ls | grep tmp"`? How many execve calls?
