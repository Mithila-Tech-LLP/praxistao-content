# Chapter 08: What Is a Process?

> **"A program is a recipe. A process is the act of cooking it. A recipe sitting in a drawer does nothing. But once you start following it — with specific ingredients, a specific pot, at a specific stove — something is actually happening. That active, running instance is the process."**

---

## Table of Contents

1. [Program vs. Process — The Key Distinction](#1-program-vs-process--the-key-distinction)
2. [What a Process Contains](#2-what-a-process-contains)
3. [Process Address Space Layout](#3-process-address-space-layout)
4. [The Process Control Block (PCB)](#4-the-process-control-block-pcb)
5. [Process Creation — fork() and exec()](#5-process-creation--fork-and-exec)
6. [Process Identifiers (PID, PPID)](#6-process-identifiers-pid-ppid)
7. [The Process Tree](#7-the-process-tree)
8. [Process Termination](#8-process-termination)
9. [Orphan and Zombie Processes](#9-orphan-and-zombie-processes)
10. [Multiple Processes vs. Multiple Threads](#10-multiple-processes-vs-multiple-threads)
11. [How the OS Sees Processes (Linux Internals)](#11-how-the-os-sees-processes-linux-internals)
12. [Summary](#summary)

---

## 1. Program vs. Process — The Key Distinction

This distinction is fundamental:

**Program:**
- A file on disk (an executable binary)
- Static — just sitting there
- Contains: machine code instructions + data + metadata
- Can be run by multiple users simultaneously without issues
- Examples: `/bin/ls`, `/usr/bin/python3`, `MyApp.exe`

**Process:**
- A program in execution
- Dynamic — actively using CPU, memory, and other resources
- Has its own memory space, file descriptors, state
- Two processes running the same program are completely independent
- Examples: your browser tab (process), each terminal window (process)

**Analogy:** A musical score (sheet music) is a program. An orchestra playing that score right now is a process. 100 orchestras around the world can play the same score simultaneously — they're 100 separate processes.

**You can run the same program multiple times:**
```bash
# Three separate processes, all running the same program
python3 script1.py &   # PID 1234
python3 script2.py &   # PID 1235
python3 script3.py &   # PID 1236
```
Each gets its own memory, its own variables, its own state. Modifying a variable in one doesn't affect the others.

---

## 2. What a Process Contains

A process is more than just the executing code. It includes:

**1. Code (Text segment):**
The compiled machine instructions. Usually read-only (can be shared between multiple instances of the same program).

**2. Data (Data segment):**
Global and static variables initialized at compile time.
```c
int global_counter = 0;    // in data segment
static int max = 100;      // in data segment
```

**3. BSS:**
Uninitialized global/static variables (OS zeroes this at startup).
```c
int uninitialized_array[1000];  // in BSS — zeroed by OS
```

**4. Heap:**
Dynamically allocated memory (`malloc`, `new`, `mmap`).
Grows upward as more memory is allocated.

**5. Stack:**
Local variables, function arguments, return addresses.
One stack per thread. Grows downward.

**6. Kernel state (not in process's own memory, but belongs to it):**
- Open file descriptors (table of open files)
- Current working directory
- Signal handlers
- User/group IDs (who owns this process)
- CPU registers (saved when process isn't running)
- Process state (running, sleeping, stopped, zombie)
- Parent process ID

---

## 3. Process Address Space Layout

Each process has its own **virtual address space** — a private view of memory that looks like the whole machine belongs to it.

On a 64-bit Linux system, a typical process's address space looks like:

```
Virtual Address
0xFFFF_FFFF_FFFF_FFFF ┐
                       │  Kernel (mapped into every process,
0xFFFF_8000_0000_0000 ─┘  but inaccessible from user mode)

0x7FFF_FFFF_FFFF_FFFF ┐
                       │  Stack (grows downward ↓)
        ↓              │  (local variables, function call frames)
        ↓              │
     [free space]      │
        ↑              │
        ↑              │  Memory-mapped files, shared libraries
0x7F000_0000_0000 ─────┘

        ↑              
     [free space]      

        ↑
        │  Heap (grows upward ↑)
        │  (malloc/new allocations)
0x0000_0000_0060_0000 ┐
                       │  BSS (uninitialized globals, zeroed)
0x0000_0000_0040_1000 ─┤  Data (initialized globals)
0x0000_0000_0040_0000 ─┘  Text (code — read only)

0x0000_0000_0000_0000 ── NULL (unmapped — accessing crashes the program)
```

**This is virtual — not real!**
The OS (via the page table) maps these virtual addresses to actual physical RAM. Two processes can both have a variable at virtual address `0x601000`, but they point to different physical pages.

**Checking a process's memory map:**
```bash
cat /proc/1234/maps    # Linux: show memory map of process 1234

# Output looks like:
# 55a8a3c00000-55a8a3c01000 r-xp 00000000 08:01 789012  /bin/cat
# 55a8a3e01000-55a8a3e02000 rw-p 00001000 08:01 789012  /bin/cat
# 7f4b5c000000-7f4b5c200000 r--p 00000000 08:01 123456  /lib/x86_64-linux-gnu/libc.so.6
# ...
```

---

## 4. The Process Control Block (PCB)

The **PCB (Process Control Block)** is the kernel data structure that describes one process. Everything the OS knows about a process lives here.

In Linux, the PCB is `struct task_struct` (defined in `linux/sched.h`):

```c
struct task_struct {
    /* Process state */
    volatile long state;        // TASK_RUNNING, TASK_INTERRUPTIBLE, ...
    
    /* Scheduling */
    int prio;                   // priority
    unsigned int time_slice;    // remaining CPU time in current slice
    
    /* Identity */
    pid_t pid;                  // process ID
    pid_t tgid;                 // thread group ID (same as pid for main thread)
    
    /* Family */
    struct task_struct *parent;  // parent process
    struct list_head children;   // list of children
    struct list_head sibling;    // list of siblings
    
    /* Memory */
    struct mm_struct *mm;        // memory descriptor (page tables, etc.)
    
    /* File system */
    struct fs_struct *fs;        // root directory, working directory
    struct files_struct *files;  // open file descriptor table
    
    /* Signals */
    struct signal_struct *signal;  // signal handlers
    sigset_t blocked;              // blocked signals mask
    
    /* Credentials */
    kuid_t uid, euid;    // real + effective user ID
    kgid_t gid, egid;   // real + effective group ID
    
    /* CPU state (saved when not running) */
    struct thread_struct thread;  // saved CPU registers
    
    // ... hundreds more fields
};
```

**One `task_struct` per process (or thread in Linux).** The scheduler works by manipulating these structures — moving them between run queues, updating their state, saving/restoring their CPU registers.

---

## 5. Process Creation — fork() and exec()

Unix systems create new processes using `fork()` and `exec()`.

**`fork()` — Cloning a process:**

```c
#include <unistd.h>
#include <stdio.h>

int main() {
    printf("Before fork\n");
    
    pid_t pid = fork();
    
    if (pid == 0) {
        // CHILD process
        printf("Child: my PID = %d, parent PID = %d\n", 
               getpid(), getppid());
    } else if (pid > 0) {
        // PARENT process
        printf("Parent: my PID = %d, child PID = %d\n", 
               getpid(), pid);
    } else {
        perror("fork failed");
        return 1;
    }
    
    printf("This runs in BOTH parent and child\n");
    return 0;
}
```

Output:
```
Before fork
Parent: my PID = 1234, child PID = 1235
This runs in BOTH parent and child
Child: my PID = 1235, parent PID = 1234
This runs in BOTH parent and child
```

**What `fork()` does inside the kernel:**
1. Allocates a new `task_struct` for the child
2. Copies the parent's `task_struct` fields into the child's
3. Creates a new virtual address space for the child, initially **sharing** pages with parent using **copy-on-write (COW)**
4. Copies open file descriptors (child inherits parent's open files)
5. Returns twice: once in parent (returns child PID), once in child (returns 0)

**Copy-on-write (COW):** After `fork()`, parent and child SHARE the same physical pages. Only when one of them WRITES to a page does the OS actually make a copy (for the writer). This makes `fork()` very fast even for large processes.

**`exec()` — Loading a new program:**

```c
#include <unistd.h>

int main() {
    char *args[] = {"/bin/ls", "-la", "/tmp", NULL};
    char *env[]  = {NULL};
    
    execve("/bin/ls", args, env);
    
    // If execve() succeeds, this line is NEVER reached!
    // The current process is completely replaced by /bin/ls.
    perror("execve failed");
    return 1;
}
```

**What `exec()` does:**
1. Opens the new executable file
2. Validates it (check ELF header, permissions)
3. Allocates new memory regions (text, data, BSS, stack)
4. Copies code and initialized data from the file into memory
5. Sets up stack with argc, argv, environment variables
6. Resets signal handlers
7. Jumps to the entry point (`_start` in the executable)

The process ID stays the SAME. The same process now runs different code.

**The classic shell pattern:**
```c
// How a shell runs a command like "ls -la /tmp":
pid_t pid = fork();
if (pid == 0) {
    // Child: replace ourselves with the requested program
    execve("/bin/ls", (char*[]){"ls", "-la", "/tmp", NULL}, environ);
    exit(1);  // only reached if execve fails
} else {
    // Parent: wait for child to finish
    int status;
    waitpid(pid, &status, 0);
}
```

---

## 6. Process Identifiers (PID, PPID)

**PID (Process ID):**
Every process has a unique integer ID. On Linux:
- PID 1: init/systemd (the first user-space process)
- PIDs are assigned sequentially, then reused after they're freed
- Default maximum: 32768 (configurable in `/proc/sys/kernel/pid_max`)

**PPID (Parent Process ID):**
Every process (except PID 1) has a parent. The PPID is the parent's PID.

```bash
$ ps aux | head -10
USER       PID %CPU %MEM    VSZ   RSS TTY      STAT START   TIME COMMAND
root         1  0.0  0.1 169344 13092 ?        Ss   07:23   0:01 /sbin/init
root         2  0.0  0.0      0     0 ?        S    07:23   0:00 [kthreadd]
...
user      1234  0.5  2.1 456789 87654 pts/0    S    09:15   0:12 python3 myscript.py

$ echo $$   # current shell's PID
5678

$ ps -p 5678 -o pid,ppid,cmd
  PID  PPID CMD
 5678  5677 bash
```

---

## 7. The Process Tree

Every process was created by another process (except PID 1). This forms a tree:

```bash
$ pstree -p

systemd(1)─┬─cron(456)
            ├─NetworkManager(789)
            ├─sshd(1234)─┬─sshd(2345)───bash(2346)───vim(2347)
            │             └─sshd(2348)───bash(2349)───python(2350)
            ├─Xorg(3456)
            └─gnome-shell(4567)─┬─chrome(5678)─┬─chrome(5679)
                                 │              ├─chrome(5680)
                                 │              └─chrome(5681)
                                 └─terminal(6789)───bash(6790)
```

**This tree structure is important because:**
- When you kill a process, you might want to kill its entire subtree (all children)
- When a parent exits, its children either become orphans (adopted by PID 1) or are killed
- Signals can be sent to process groups (collections of processes in the tree)

---

## 8. Process Termination

A process can terminate in several ways:

**Normal exit:**
```c
// Option 1: return from main()
int main() {
    return 0;   // exit code 0 = success
}

// Option 2: call exit()
#include <stdlib.h>
exit(0);    // also runs atexit() handlers and flushes stdio buffers

// Option 3: call _exit() (raw syscall, no cleanup)
#include <unistd.h>
_exit(0);
```

**Killed by signal:**
```bash
kill -9 1234    # SIGKILL — immediate, cannot be ignored
kill -15 1234   # SIGTERM — polite request to terminate (can be handled)
kill -2 1234    # SIGINT — same as Ctrl+C
```

**Crash:**
- Segmentation fault (SIGSEGV) — accessed memory it shouldn't
- Illegal instruction (SIGILL) — executed bad instruction
- Floating point error (SIGFPE) — division by zero

**Exit status:**
Every process exits with a status code (0–255). 0 = success, non-zero = error.
The parent can retrieve this via `wait()`:
```c
int status;
waitpid(child_pid, &status, 0);
if (WIFEXITED(status)) {
    printf("Child exited with code %d\n", WEXITSTATUS(status));
} else if (WIFSIGNALED(status)) {
    printf("Child killed by signal %d\n", WTERMSIG(status));
}
```

---

## 9. Orphan and Zombie Processes

**Orphan Process:**
When a parent process exits before its child, the child becomes an orphan.

Linux automatically adopts orphans: PID 1 (systemd/init) becomes the new parent. This is called "reparenting."

```
Before: shell(1234) → python(5678)
Shell exits.
After: systemd(1) → python(5678)
```

Orphans are harmless — systemd will `wait()` for them when they exit.

**Zombie Process:**
When a child exits but its parent hasn't called `wait()` yet, the child becomes a "zombie."

```
Child calls exit(0).
Kernel can't fully clean up child's task_struct yet.
Why? The parent might want to read the exit code.
Kernel keeps a "zombie" entry (minimal, just the exit code).
Status shows: Z (zombie)
```

Zombies can't be killed (they're already dead!). They go away when the parent calls `wait()`.

**Resource leak if zombies accumulate:**
If a program creates children and never `wait()`s for them, they pile up as zombies, slowly consuming PID table slots. Eventually, no new PIDs can be allocated → system can't create new processes.

**Fix: always `wait()` for your children:**
```c
// Using SIGCHLD signal to clean up asynchronously:
signal(SIGCHLD, SIG_IGN);   // simple: tell OS to auto-reap zombies

// Or explicitly:
waitpid(-1, NULL, WNOHANG); // non-blocking wait for any child
```

---

## 10. Multiple Processes vs. Multiple Threads

When you want multiple tasks running simultaneously, you can use:

**Multiple Processes:**
- Each process has its own address space (isolated)
- Communication via IPC (pipes, sockets, shared memory)
- One process crashing doesn't affect others
- Higher overhead (separate memory space, more context switch work)
- Example: Chrome uses separate processes for each tab — one tab crashing doesn't crash others

**Multiple Threads:**
- Threads within a process share the same address space
- Communication by sharing memory variables (fast but requires synchronization)
- One thread crashing can crash the whole process
- Lower overhead (shared address space, lighter context switch)
- Example: A web server uses threads to handle concurrent connections

**When to use each:**
```
Use multiple processes when:
  - Isolation is critical (browser sandboxing, microservices)
  - Tasks are independent with little shared state
  - Security boundary needed between tasks

Use multiple threads when:
  - Tasks share lots of data
  - Low-latency communication needed
  - Memory is limited (threads are lighter)
  - Example: parallel computation on shared data
```

We'll cover threads in depth in Chapter 11.

---

## 11. How the OS Sees Processes (Linux Internals)

**Process list:**
All processes are in a circular doubly-linked list of `task_struct` pointers. The `init_task` (PID 0 — the idle process) is the list head.

```c
// Iterate over all processes in the kernel:
struct task_struct *task;
for_each_process(task) {
    printk("PID: %d, Name: %s\n", task->pid, task->comm);
}
```

**Run queues:**
Each CPU core has a run queue — a sorted data structure of runnable processes waiting for CPU time. The scheduler picks the highest-priority runnable task from the run queue.

**Procfs interface:**
Linux exposes process information through the `/proc` pseudo-filesystem:
```bash
ls /proc/1234/      # directory for process PID 1234
# cmdline          — command and arguments
# status           — process status (state, memory, UID, etc.)
# maps             — virtual memory map
# fd/              — symbolic links to open file descriptors
# exe              — symlink to the executable file
# cwd              — symlink to current working directory
# mem              — process memory (can read with ptrace permission)
# net/             — network statistics

cat /proc/1234/status
# Name: python3
# Pid: 1234
# PPid: 5678
# VmSize: 456789 kB   (virtual memory)
# VmRSS: 87654 kB     (physical RAM used)
# ...
```

---

## Summary

| Concept | Definition |
|---------|-----------|
| Program | Static executable file on disk |
| Process | Running instance of a program with its own memory, state |
| Address space | Each process's private virtual memory view |
| PCB / task_struct | Kernel data structure describing one process |
| fork() | Create a copy of the current process |
| exec() | Replace current process with a new program |
| PID | Unique process identifier |
| PPID | Parent process ID |
| Process tree | Hierarchy of processes; everything descends from PID 1 |
| Zombie | Process that exited but parent hasn't wait()ed |
| Orphan | Process whose parent exited; adopted by PID 1 |
| Copy-on-write | fork() shares pages; copies only on write (efficient) |

**The most important thing to understand:** A process is the OS's unit of resource ownership. Memory belongs to a process. File descriptors belong to a process. CPU time is given to a process. When a process dies, the OS reclaims all of its resources. This is what prevents resource leaks at the system level.
