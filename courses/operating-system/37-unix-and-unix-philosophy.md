# Chapter 37: Unix and the Unix Philosophy

> **"This is the Unix philosophy: Write programs that do one thing and do it well. Write programs that work together. Write programs that handle text streams, because that is a universal interface. — Doug McIlroy (1978)"**

---

## Table of Contents

1. [The Birth of Unix](#1-the-birth-of-unix)
2. [The Unix Philosophy](#2-the-unix-philosophy)
3. [Everything Is a File](#3-everything-is-a-file)
4. [The Process Model](#4-the-process-model)
5. [Pipes — The Unix Glue](#5-pipes--the-unix-glue)
6. [The Shell — A Programmable Interface](#6-the-shell--a-programmable-interface)
7. [The Unix Kernel Architecture](#7-the-unix-kernel-architecture)
8. [Unix Variants — BSD, System V, POSIX](#8-unix-variants--bsd-system-v-posix)
9. [Unix's Legacy](#9-unixs-legacy)
10. [Summary](#summary)

---

## 1. The Birth of Unix

**1969 — Bell Labs, Murray Hill, New Jersey:**
Ken Thompson, Dennis Ritchie, and others at AT&T Bell Labs had been working on the Multics operating system — an ambitious, complex system. When AT&T withdrew from Multics, Thompson and Ritchie rebuilt their ideas in a simpler, more elegant form.

Thompson wrote the first Unix in assembly language on a PDP-7 computer. The goals were deliberately modest:
- An environment to develop programs comfortably
- File system, process model, simple commands
- Small enough for one person to understand

**1972-1973:**
Dennis Ritchie invented the C programming language specifically to rewrite Unix. Unix became the first OS written mostly in a high-level language — enabling it to be **portable**. The kernel was ~10,000 lines of C.

**1976-1979: BSD (Berkeley Software Distribution):**
UC Berkeley received Unix source code and made significant improvements, eventually releasing BSD Unix with new features (virtual memory, TCP/IP networking, fast file system).

**1983:**
Richard Stallman starts the GNU Project to create a free Unix-compatible OS.

**1991:**
Linus Torvalds writes Linux — a free Unix-compatible kernel. GNU + Linux = a complete free Unix-like OS.

**The key insight:**
Unix survived decades not because of any single feature, but because of a **design philosophy** that led to small, composable, interoperable tools.

---

## 2. The Unix Philosophy

Doug McIlroy summarized the Unix philosophy in three rules:

**Rule 1: Do one thing and do it well.**
```
ls:   list files (just listing, nothing else)
wc:   count words/lines/bytes
sort: sort lines
grep: search for patterns
cut:  select columns
head: show first N lines
tail: show last N lines
```
Each tool is small, focused, and completely mastered.

**Rule 2: Write programs that work together.**
```bash
# Programs communicate via text streams (stdout/stdin)
ls /etc | grep conf | wc -l
# How many .conf files are in /etc?
# ls: lists files, grep: filters by "conf", wc: counts

ps aux | grep python | awk '{print $2}' | xargs kill
# Kill all python processes
```

**Rule 3: Write programs that handle text streams.**
Text is the universal interface:
- Human-readable
- No special format negotiation
- Every tool can read text, every tool can write text
- Pipes carry text between programs

**Anti-Unix (rule violations):**
```
Bad: One monolithic program that does file management, email, web browsing, music...
Good: Separate programs, each expert in their domain, composable

Bad: Binary-only output that requires special parsers
Good: Column-formatted text that awk/grep/cut can parse

Bad: Graphical configuration that can't be scripted
Good: Text config files that can be automated with sed/awk/python
```

---

## 3. Everything Is a File

**Unix's most powerful abstraction:** Everything is a file descriptor.

```
Regular files:           /home/user/document.txt
Directories:             /home/user/
Devices:                 /dev/sda, /dev/tty0
Sockets (Unix domain):   /tmp/myprogram.sock
Pipes:                   /tmp/mypipe (named pipe)
Kernel data:             /proc/cpuinfo, /proc/1234/maps
Hardware config:         /sys/block/sda/queue/scheduler
Network sockets:         fd = socket(AF_INET, SOCK_STREAM, 0)
```

**Practical consequence:**
A program written to process file input automatically works on:
- Regular files: `cat /etc/passwd`
- Pipe input: `ps aux | grep python`
- Network data: `curl http://example.com | grep title`
- Terminal input: `read input_line`

```bash
# The same program, different inputs, no code changes:
sort /var/log/access.log         # sort a file
cat bigfile | sort               # sort from stdin
curl -s api.example.com | sort   # sort network data
```

**File descriptors are the universal handle:**
```c
// These all use the same read/write API:
int file_fd    = open("/etc/passwd", O_RDONLY);
int socket_fd  = accept(server_fd, NULL, NULL);
int pipe_fd[2]; pipe(pipe_fd);
int dev_fd     = open("/dev/urandom", O_RDONLY);

// All use the same read():
read(file_fd,   buf, sizeof(buf));
read(socket_fd, buf, sizeof(buf));
read(pipe_fd[0], buf, sizeof(buf));
read(dev_fd,    buf, sizeof(buf));
```

---

## 4. The Process Model

**Unix's elegant process model:**

**fork() + exec():**
```c
pid_t pid = fork();
if (pid == 0) {
    // Child: replace with new program
    execve("/bin/ls", argv, envp);
} else {
    // Parent: wait for child
    waitpid(pid, &status, 0);
}
```

**Why fork + exec instead of spawn?**
- Simple: two independent operations, each simple
- Flexible: between fork and exec, you can: redirect files, close handles, change permissions
- Universal: the shell implements pipelines and redirections using fork+exec+dup2

**Everything starts from init:**
```
PID 1: init/systemd
├── PID 2: kthreadd (kernel threads)
├── PID 100: sshd
│   └── PID 150: sshd (connection for alice)
│       └── PID 151: bash (alice's shell)
│           └── PID 200: ls /home (alice ran ls)
├── PID 300: cron
│   └── PID 400: backup.sh (cron ran this at midnight)
└── PID 500: getty (login prompt)
    └── PID 501: login
        └── PID 502: bash (logged-in user's shell)
```

Every process is a descendant of PID 1. A process always has a parent. When a parent dies, orphan processes are adopted by PID 1 (systemd).

---

## 5. Pipes — The Unix Glue

**Pipes are the mechanism for combining programs:**

```bash
# Sort all logged-in users by username:
who | awk '{print $1}' | sort | uniq

# 1. who: lists logged-in users
# 2. awk: extracts the username field
# 3. sort: alphabetically sorts usernames
# 4. uniq: removes duplicates

# Each program is simple; the combination is powerful
```

**The pipe system call:**
```c
int pipefd[2];
pipe(pipefd);  // pipefd[0] = read end, pipefd[1] = write end

pid_t pid = fork();
if (pid == 0) {
    // Child: write to pipe
    close(pipefd[0]);                 // close read end
    dup2(pipefd[1], STDOUT_FILENO);   // redirect stdout to pipe
    close(pipefd[1]);
    execve("/usr/bin/ls", ...);       // ls output goes to pipe
} else {
    // Parent: read from pipe
    close(pipefd[1]);                 // close write end
    dup2(pipefd[0], STDIN_FILENO);    // redirect stdin from pipe
    close(pipefd[0]);
    execve("/usr/bin/wc", ...);       // wc reads from pipe
    waitpid(pid, ...);
}
```

This is exactly what the shell does when you type `ls | wc`.

**Named pipes (FIFOs):**
```bash
mkfifo /tmp/mypipe
ls > /tmp/mypipe &    # write to named pipe (backgrounds)
wc < /tmp/mypipe      # read from named pipe
```

---

## 6. The Shell — A Programmable Interface

**The Unix shell is both:**
1. A command interpreter (type commands, they run)
2. A scripting language (automate sequences of commands)

**Key shells:**
```
sh:   Bourne Shell (original Unix shell, 1979) — simple, portable
csh:  C Shell — C-like syntax, history, job control
ksh:  Korn Shell — extended sh with csh features
bash: Bourne-Again SHell — most common on Linux, extends sh
zsh:  Z Shell — bash-compatible with many extensions
fish: Friendly Interactive Shell — modern, user-friendly
```

**Shell as glue:**
```bash
#!/bin/bash
# Backup all changed files in the last 24 hours:

BACKUP_DIR="/backup/$(date +%Y%m%d)"
mkdir -p "$BACKUP_DIR"

find /home -mtime -1 -type f | while read file; do
    # Preserve directory structure:
    dest="$BACKUP_DIR/$file"
    mkdir -p "$(dirname "$dest")"
    cp "$file" "$dest"
done

echo "Backup complete: $BACKUP_DIR"
```

**Signal handling in the shell:**
```bash
# Ctrl+C sends SIGINT to the foreground process group
# Ctrl+Z suspends (SIGTSTP) the foreground process
# Ctrl+\ sends SIGQUIT (core dump)
# bg: send suspended process to background
# fg: bring background process to foreground
# kill %1: send SIGTERM to job 1
```

---

## 7. The Unix Kernel Architecture

**Classic Unix kernel structure (V7 Unix, 1979):**

```
User space:
  Shell, utilities (ls, cat, grep), user programs

System call interface:
  ~60 system calls: open, read, write, close, fork, exec, wait, pipe, ...

Kernel:
  File system:
    Directory operations, inode management, buffer cache
  Process management:
    Scheduling (simple priority), fork/exec/wait
  Memory management:
    Swapping (swap entire processes to disk!)
  Device management:
    Block devices, character devices, interrupt handlers
  
Hardware:
  CPU, disk, terminal, clock
```

**Key difference from modern Linux:**
- V7 Unix: no demand paging (swap whole processes, not pages)
- No virtual memory in early Unix
- Single CPU only (SMP came later with Mach/BSD/System V)
- Minimal security model (just setuid bits, no capabilities)

---

## 8. Unix Variants — BSD, System V, POSIX

**The Unix family tree:**

```
AT&T Unix V7 (1979)
├── BSD (Berkeley Software Distribution)
│   ├── 4.1BSD, 4.2BSD (added TCP/IP, 1983)
│   ├── 4.3BSD → 4.4BSD (freely redistributable, 1994)
│   ├── FreeBSD (modern, server/desktop)
│   ├── OpenBSD (security-focused)
│   ├── NetBSD (most portable)
│   └── macOS / iOS (Darwin kernel = XNU = Mach + FreeBSD userspace)
│
├── System V (AT&T commercial Unix)
│   ├── Solaris (Sun/Oracle)
│   ├── HP-UX (Hewlett-Packard)
│   ├── AIX (IBM)
│   └── IRIX (Silicon Graphics)
│
└── Linux (1991, inspired by Unix but written from scratch)
    ├── Debian/Ubuntu
    ├── Red Hat/Fedora/CentOS
    ├── Android (Linux kernel + Google's userspace)
    └── many others
```

**POSIX (Portable Operating System Interface):**
IEEE standard defining the Unix API so applications can compile on any conforming OS:
```
POSIX.1: Core API (file operations, processes, signals)
POSIX.2: Shell and utilities (sh, grep, sed, awk)
POSIX threads: pthreads API
```

A program written to POSIX can compile on Linux, macOS, FreeBSD, Solaris — the same source code.

---

## 9. Unix's Legacy

**What Unix contributed to computing:**

1. **Pipes and composable tools:** Programs as building blocks; each does one thing well
2. **Everything is a file:** Uniform interface for devices, sockets, kernel data
3. **C and portability:** OS written in a high-level language; could be ported to new hardware
4. **Hierarchical file system:** Directories containing files; mount points composing namespaces
5. **Fork/exec process model:** Simple, flexible, elegant
6. **Text as the universal interface:** Human-readable, automatable, composable
7. **The shell:** Programming and administration in one tool

**Operating systems influenced by Unix:**
- Linux: GPL Unix clone
- macOS/iOS: BSD/Mach hybrid
- Android: Linux kernel
- Plan 9: "what if we took Unix philosophy further?" (everything is truly a file, distributed)
- QNX: real-time microkernel, POSIX-compatible
- Solaris: AT&T System V successor, ZFS

---

## Summary

| Concept | Description |
|---------|------------|
| Unix philosophy | Do one thing well; work together via text streams |
| Everything is a file | Files, devices, sockets, pipes all share the fd interface |
| Fork + exec | Process creation: clone then replace with new program |
| Pipe | One-way channel connecting stdout of one process to stdin of another |
| Shell | Command interpreter + scripting language; glues programs together |
| Signal | Async notification to a process (SIGTERM, SIGINT, SIGKILL, SIGHUP) |
| POSIX | Standard API ensuring portability across Unix variants |
| BSD | Berkeley Unix branch; contributed TCP/IP, fast file system, vm |
| System V | AT&T commercial Unix; contributed System V IPC, SVR4 features |
| GNU project | Free implementations of Unix utilities (gcc, glibc, bash, coreutils) |
| Linux | Free Unix-like kernel; combined with GNU tools gives a complete OS |
