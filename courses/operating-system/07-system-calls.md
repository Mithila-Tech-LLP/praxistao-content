# Chapter 07: System Calls — Talking to the Kernel

> **"A system call is a knock on the kernel's door. The kernel opens the door, checks your ID, does what you asked, closes the door, and sends you the result. You never actually enter the kernel's house — you just pass notes through the mail slot."**

---

## Table of Contents

1. [The Problem System Calls Solve](#1-the-problem-system-calls-solve)
2. [What Is a System Call?](#2-what-is-a-system-call)
3. [How a System Call Works — Step by Step](#3-how-a-system-call-works--step-by-step)
4. [System Call Numbers — The Catalog](#4-system-call-numbers--the-catalog)
5. [System Calls in C (The High-Level View)](#5-system-calls-in-c-the-high-level-view)
6. [System Calls in Assembly (The Low-Level View)](#6-system-calls-in-assembly-the-low-level-view)
7. [The System Call Dispatch Table](#7-the-system-call-dispatch-table)
8. [Categories of System Calls](#8-categories-of-system-calls)
9. [Important Unix/Linux System Calls](#9-important-unixlinux-system-calls)
10. [System Call Performance — Why They're Not Free](#10-system-call-performance--why-theyre-not-free)
11. [strace — Spying on System Calls](#11-strace--spying-on-system-calls)
12. [Windows System Calls (NT Native API)](#12-windows-system-calls-nt-native-api)
13. [Summary](#summary)

---

## 1. The Problem System Calls Solve

A program runs in user space. It needs to:
- Read a file (the disk is hardware — only the kernel can touch it)
- Allocate memory (the physical pages are managed by the kernel)
- Send data over the network (the network card is hardware)
- Create a new process (process creation requires kernel data structures)

But user space programs CANNOT do any of these things directly. They don't have permission.

**The dilemma:**
- Programs need kernel services to do anything useful
- Programs must not have unrestricted access to the kernel
- Solution: a controlled, audited mechanism to ask the kernel for services

That mechanism is the **system call**.

---

## 2. What Is a System Call?

A **system call (syscall)** is a special function that causes the CPU to switch from user mode (Ring 3) to kernel mode (Ring 0) in a controlled, safe way.

The transition is controlled:
- The CPU only jumps to one specific kernel location (the syscall entry point)
- The kernel validates the request before doing anything
- The kernel returns to user mode when done
- User code can never stay in kernel mode — it's always a round trip

**Analogy:**
Think of a bank. You (user process) can't walk into the vault (kernel). Instead, you fill out a form (set up register arguments), slide it through the window (execute the syscall instruction), the teller (kernel handler) processes it, and slides the result back. You never entered the secure area.

---

## 3. How a System Call Works — Step by Step

Let's trace a `read()` system call in detail on x86-64 Linux:

**User code (in C):**
```c
ssize_t bytes = read(fd, buffer, count);
```

**What actually happens:**

```
Step 1: C library wraps the syscall
  glibc's read() function sets up registers:
    RAX = 0 (system call number for read on x86-64)
    RDI = fd (first argument)
    RSI = buffer (second argument)
    RDX = count (third argument)
  
  Then executes: syscall instruction

Step 2: CPU executes the syscall instruction
  CPU hardware does ALL of this automatically:
  a. Saves current RIP (return address) in RCX
  b. Saves current RFLAGS in R11
  c. Switches to kernel mode (Ring 0)
  d. Loads kernel stack pointer from MSR_LSTAR
  e. Jumps to the kernel syscall entry point (also from MSR_LSTAR)

Step 3: Kernel syscall entry (assembly, arch/x86/entry/entry_64.S in Linux)
  a. Save all user registers to the kernel stack (entire CPU state)
  b. Switch to kernel stack (if not already)
  c. Call do_syscall_64(regs, nr)

Step 4: Kernel syscall dispatcher (C code)
  a. nr = RAX = 0 (the syscall number)
  b. Look up sys_call_table[0] → sys_read
  c. Call sys_read(fd, buffer, count)

Step 5: sys_read() executes
  a. Validate fd: is it a valid open file descriptor for this process?
  b. Check permissions: can this process read this file?
  c. Read data from the file (maybe from page cache, maybe from disk)
  d. Copy data to user buffer (buffer is a user-space address)
  e. Return number of bytes read (or negative error code)

Step 6: Return to user space
  a. Put return value into RAX
  b. Restore user registers from kernel stack
  c. Execute sysret instruction (inverse of syscall)
  CPU hardware does:
    - Restores RIP from RCX (saved return address)
    - Restores RFLAGS from R11
    - Switches back to user mode (Ring 3)

Step 7: Back in C library
  glibc checks if RAX is negative (error)
  If yes: sets errno = -RAX, returns -1
  If no: returns RAX (bytes read)

Step 8: Back in user code
  bytes = return value from glibc read()
```

**The total cost:** This round trip takes about 100–500 nanoseconds. Cheap for reading a disk (milliseconds anyway), but significant if done millions of times per second.

---

## 4. System Call Numbers — The Catalog

Every system call has a unique integer number. The program puts this number in the appropriate register before triggering the syscall.

**x86-64 Linux syscall ABI:**
- Syscall number: `RAX`
- Arguments: `RDI, RSI, RDX, R10, R8, R9` (up to 6 args)
- Return value: `RAX`
- Trigger: `syscall` instruction

**x86-32 Linux syscall ABI (legacy):**
- Syscall number: `EAX`
- Arguments: `EBX, ECX, EDX, ESI, EDI, EBP`
- Return value: `EAX`
- Trigger: `int 0x80` instruction

**Selected x86-64 Linux syscall numbers:**

| Number | Name | What it does |
|--------|------|-------------|
| 0 | `read` | Read from file descriptor |
| 1 | `write` | Write to file descriptor |
| 2 | `open` | Open a file |
| 3 | `close` | Close a file descriptor |
| 4 | `stat` | Get file metadata |
| 9 | `mmap` | Map file/memory into address space |
| 11 | `munmap` | Unmap memory |
| 12 | `brk` | Change heap size |
| 20 | `writev` | Write multiple buffers (scatter-gather) |
| 22 | `pipe` | Create a pipe |
| 32 | `dup` | Duplicate file descriptor |
| 39 | `getpid` | Get current process ID |
| 41 | `socket` | Create a network socket |
| 49 | `bind` | Bind socket to address |
| 56 | `clone` | Create a new thread/process |
| 57 | `fork` | Fork (create child process) |
| 59 | `execve` | Execute a new program |
| 60 | `exit` | Terminate process |
| 61 | `wait4` | Wait for child process |
| 231 | `exit_group` | Terminate all threads in process |

The full Linux syscall table has ~300+ entries. The number must never change once defined — that would break all existing binaries.

---

## 5. System Calls in C (The High-Level View)

In practice, you never write syscall assembly yourself. The C standard library (glibc on Linux) provides wrapper functions:

```c
#include <unistd.h>
#include <sys/types.h>
#include <sys/stat.h>
#include <fcntl.h>

int main() {
    // open() is a glibc wrapper around the open (or openat) syscall
    int fd = open("hello.txt", O_RDONLY);
    if (fd < 0) {
        perror("open failed");
        return 1;
    }
    
    char buf[100];
    // read() is a glibc wrapper around the read syscall
    ssize_t n = read(fd, buf, sizeof(buf) - 1);
    
    if (n > 0) {
        buf[n] = '\0';
        // write() is a glibc wrapper around the write syscall
        write(1, buf, n);    // 1 = stdout
    }
    
    // close() → close syscall
    close(fd);
    return 0;
}
```

**What glibc's `read()` wrapper looks like internally:**
```c
// Simplified — real glibc is more complex
ssize_t read(int fd, void *buf, size_t count) {
    ssize_t result = syscall(SYS_read, fd, buf, count);
    if (result < 0) {
        errno = -result;  // kernel returns negative errno
        return -1;
    }
    return result;
}
```

**The `syscall()` function** in glibc is the generic wrapper:
```c
#include <unistd.h>
#include <sys/syscall.h>

long result = syscall(SYS_write, 1, "hello\n", 6);
```

---

## 6. System Calls in Assembly (The Low-Level View)

Let's write a complete "Hello, World" program that makes syscalls directly, without using any C library:

```nasm
; hello_syscall.asm
; nasm -f elf64 hello_syscall.asm -o hello.o
; ld hello.o -o hello
; ./hello

section .data
    message db "Hello, World!", 10  ; "Hello, World!\n"
    msg_len equ $ - message         ; length = 14

section .text
    global _start

_start:
    ; write(1, message, msg_len)
    ; syscall 1 = write
    ; arg1 (rdi) = 1 (stdout file descriptor)
    ; arg2 (rsi) = message address
    ; arg3 (rdx) = length
    mov rax, 1          ; syscall number: write
    mov rdi, 1          ; arg1: fd = 1 (stdout)
    mov rsi, message    ; arg2: pointer to message
    mov rdx, msg_len    ; arg3: length = 14
    syscall             ; execute syscall
    
    ; exit(0)
    ; syscall 60 = exit
    ; arg1 (rdi) = 0 (exit code)
    mov rax, 60         ; syscall number: exit
    mov rdi, 0          ; arg1: exit code = 0
    syscall             ; execute syscall
```

This program has NO C library, NO runtime. Just two system calls. It compiles to a ~100 byte binary.

**Understanding `_start`:**
The linker expects `_start` as the entry point when no C library is used. When using C (`int main()`), the C library provides `_start`, which sets up argc/argv, calls global constructors, then calls `main()`.

---

## 7. The System Call Dispatch Table

In the Linux kernel, the system call table is:
```c
// arch/x86/entry/syscalls/syscall_64.tbl (abbreviated)
// number  abi    name      entry point
0          common  read      sys_read
1          common  write     sys_write
2          common  open      sys_open
3          common  close     sys_close
...
```

This generates a table `sys_call_table[]`:
```c
// Conceptually:
typedef long (*syscall_fn_t)(long, long, long, long, long, long);

syscall_fn_t sys_call_table[NR_syscalls] = {
    [0] = sys_read,
    [1] = sys_write,
    [2] = sys_open,
    // ...
};
```

The kernel dispatcher:
```c
long do_syscall_64(unsigned long nr, struct pt_regs *regs) {
    if (nr >= NR_syscalls)
        return -ENOSYS;   // error: no such syscall
    
    return sys_call_table[nr](regs->di, regs->si, regs->dx,
                              regs->r10, regs->r8, regs->r9);
}
```

**When we build our OS, we'll implement this exact pattern** — a dispatch table indexed by syscall number.

---

## 8. Categories of System Calls

System calls fall into logical groups:

**1. Process Control:**
```
fork()     — create a copy of the current process
exec()     — replace current process with a new program
exit()     — terminate
wait()     — wait for a child to finish
getpid()   — get process ID
kill()     — send a signal to a process
```

**2. File Management:**
```
open()     — open or create a file
read()     — read bytes from fd
write()    — write bytes to fd
close()    — close file descriptor
lseek()    — move read/write position
stat()     — get file metadata
unlink()   — delete a file
rename()   — rename a file
mkdir()    — create directory
rmdir()    — remove directory
```

**3. Device Management:**
```
ioctl()    — device-specific control
mmap()     — map device into memory
read/write — also work on device fds
```

**4. Information Maintenance:**
```
getpid()   — process ID
getuid()   — user ID
time()     — current time
gettimeofday() — precise time
uname()    — OS/kernel info
sysinfo()  — RAM, load average
```

**5. Communication:**
```
pipe()     — create an anonymous pipe
socket()   — create a network socket
bind()     — bind socket to address
connect()  — connect to remote socket
send()/recv() — send/receive data
shmget()   — create shared memory segment
mq_open()  — create message queue
```

---

## 9. Important Unix/Linux System Calls

**`fork()` — Creating a new process:**
```c
pid_t pid = fork();
if (pid == 0) {
    // This code runs in the CHILD process
    printf("I am the child, my PID is %d\n", getpid());
} else if (pid > 0) {
    // This code runs in the PARENT process
    printf("I am the parent, my child's PID is %d\n", pid);
} else {
    // fork() returned -1: error
    perror("fork failed");
}
```
`fork()` creates an exact copy of the calling process. Parent gets the child's PID. Child gets 0. Both processes continue running from the next line after `fork()`.

**`exec()` — Replacing a process with a new program:**
```c
execve("/bin/ls", (char*[]){"/bin/ls", "-la", NULL}, environ);
// If execve() succeeds, this line is NEVER reached
// The process is completely replaced by /bin/ls
perror("execve failed");
```
`exec()` doesn't create a new process — it replaces the current process's code, data, and stack with a new program. Combined with `fork()`, this is how shells run programs: `fork()` creates a child, child calls `exec()` with the command.

**`mmap()` — Map memory:**
```c
// Map a file into memory
void *addr = mmap(NULL, file_size, PROT_READ, MAP_SHARED, fd, 0);
// Now you can access file contents as if they're in memory!
char first_byte = ((char*)addr)[0];

// Allocate anonymous memory (this is what malloc uses internally)
void *mem = mmap(NULL, 4096, PROT_READ|PROT_WRITE, 
                  MAP_PRIVATE|MAP_ANONYMOUS, -1, 0);
```

**`write()` and `read()` work on everything:**
```c
// Write to a regular file
write(file_fd, "hello\n", 6);

// Write to stdout (fd = 1)
write(1, "hello\n", 6);

// Write to a network socket
write(socket_fd, "HTTP/1.1 200 OK\r\n", 17);

// Write to a pipe
write(pipe_fd[1], "message", 7);
```
Same `write()` system call for all of them. The OS routes it to the right handler.

---

## 10. System Call Performance — Why They're Not Free

A system call is NOT free. Each one costs roughly:

| CPU | Typical syscall cost |
|-----|---------------------|
| Intel i9 (modern) | ~60–100 nanoseconds |
| With Spectre/Meltdown mitigations | ~200–1000 nanoseconds |
| Old hardware | ~500 ns+ |

**Sources of overhead:**
1. **Mode switch:** CPU must save user registers, switch stacks, load kernel stack
2. **Cache flushing:** Spectre mitigations (KPTI on vulnerable CPUs) flush TLB on every switch
3. **Cache misses:** Kernel code is not in the L1/L2 cache when programs are running
4. **Security checks:** Kernel must validate every argument before trusting it

**Strategies to reduce syscall overhead:**

**Batching:** Do work in big chunks, not tiny ones.
```c
// Bad: 1000 syscalls
for (int i = 0; i < 1000; i++) write(fd, &data[i], 1);

// Good: 1 syscall
write(fd, data, 1000);
```

**io_uring (Linux 5.1+):**
A new Linux kernel feature where user space can submit batches of I/O operations to a shared ring buffer, and the kernel processes them asynchronously. Dramatically reduces syscall overhead for I/O-intensive programs.

**vDSO (Virtual Dynamic Shared Object):**
Some system calls don't actually need kernel mode! `gettimeofday()` and `clock_gettime()` — instead of a full syscall, Linux maps a small piece of kernel memory into user space that contains the current time. User code reads it directly, no mode switch needed.

---

## 11. strace — Spying on System Calls

`strace` is an invaluable debugging tool that shows you every system call a program makes:

```bash
$ strace ls /tmp

execve("/bin/ls", ["ls", "/tmp"], 0x7fff...env[]) = 0
brk(NULL)                               = 0x55a8a3d2d000
mmap(NULL, 8192, PROT_READ|PROT_WRITE, MAP_PRIVATE|MAP_ANONYMOUS, -1, 0) = 0x7f4...
access("/etc/ld.so.preload", R_OK)      = -1 ENOENT (No such file or directory)
openat(AT_FDCWD, "/etc/ld.so.cache", O_RDONLY|O_CLOEXEC) = 3
fstat(3, {st_mode=S_IFREG|0644, st_size=155786, ...}) = 0
mmap(NULL, 155786, PROT_READ, MAP_PRIVATE, 3, 0) = 0x7f4...
close(3)                                = 0
openat(AT_FDCWD, "/lib/x86_64-linux-gnu/libselinux.so.1", O_RDONLY|O_CLOEXEC) = 3
...
getdents64(3, 0x55a8a3d2e4d0, 32768)    = 120
...
write(1, "file1.txt\nfile2.txt\n", 20)  = 20
close(3)                                = 0
exit_group(0)                           = ?
+++ exited with 0 +++
```

This is deeply educational — you can see EXACTLY what the OS does when `ls` runs.

**Practical uses:**
```bash
# Count syscalls for a command
strace -c ls /tmp

# Trace only file-related syscalls
strace -e trace=file ls /tmp

# Attach to a running process (PID 1234)
strace -p 1234

# Write trace to a file
strace -o trace.txt ls /tmp
```

---

## 12. Windows System Calls (NT Native API)

Windows has its own syscall mechanism, but it's intentionally hidden.

**Windows syscall mechanism:**
- The NT kernel has a "Native API" with about 400 system services
- Names start with `Nt` or `Zw` (e.g., `NtCreateFile`, `NtReadFile`)
- Numbers change between Windows versions (not stable like Linux!)
- The numbers are NOT public API — Microsoft provides Win32 API instead

**The Windows API stack:**
```
User Code:
  CreateFile("hello.txt", ...)    ← Win32 API (documented, stable)
       ↓
  kernel32.dll → CreateFileW()
       ↓
  ntdll.dll → NtCreateFile()      ← NT Native API (semi-documented)
       ↓
  syscall instruction             ← kernel mode transition
       ↓
  Windows NT Kernel: NtCreateFile() implementation
```

**Why Windows syscall numbers are unstable:**
Microsoft doesn't want developers using the native API directly (they'd bypass compatibility layers). By changing the numbers, they ensure only ntdll.dll (which they control) makes direct syscalls. This is different from Linux's philosophy of stable ABI.

---

## Summary

| Concept | Explanation |
|---------|------------|
| System call | Controlled way for user space to request kernel services |
| Syscall instruction | x86-64: `syscall`; x86-32: `int 0x80`; causes mode switch |
| Syscall number | Integer in RAX/EAX identifying which service to invoke |
| Syscall arguments | Up to 6, in registers (RDI, RSI, RDX, R10, R8, R9 on x86-64) |
| Return value | RAX: ≥0 = success, negative = -errno error code |
| glibc wrapper | C library function (read, write, open) that handles the syscall |
| sys_call_table | Kernel's array of function pointers indexed by syscall number |
| strace | Tool to trace all syscalls a process makes |
| vDSO | Some syscalls mapped to user-readable memory — no mode switch |
| io_uring | New Linux mechanism: batch I/O requests without many syscalls |

**The three most important syscalls to understand early:**
1. `write(fd, buf, count)` — output (all I/O goes through this)
2. `fork()` — create processes (how shells run programs)
3. `execve(path, argv, envp)` — load and run programs

Master these three and you understand how the shell, every program launcher, and most of user-space computing work.
