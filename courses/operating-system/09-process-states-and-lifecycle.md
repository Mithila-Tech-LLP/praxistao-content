# Chapter 09: Process States and Lifecycle

> **"A process is never just 'running.' At any given moment, it is in one specific state — and the OS keeps precise track of which state every process is in. Understanding state transitions is understanding how an OS manages hundreds of concurrent tasks on one or a few CPUs."**

---

## Table of Contents

1. [Why Process States Exist](#1-why-process-states-exist)
2. [The Five Core Process States](#2-the-five-core-process-states)
3. [State Transition Diagram](#3-state-transition-diagram)
4. [Linux Process States in Detail](#4-linux-process-states-in-detail)
5. [Running vs. Runnable — A Common Confusion](#5-running-vs-runnable--a-common-confusion)
6. [What Causes State Transitions](#6-what-causes-state-transitions)
7. [The Blocked State — Waiting for Events](#7-the-blocked-state--waiting-for-events)
8. [The Zombie State Revisited](#8-the-zombie-state-revisited)
9. [Seeing Process States](#9-seeing-process-states)
10. [Process Priorities](#10-process-priorities)
11. [Process Lifecycle Example — Running a Shell Command](#11-process-lifecycle-example--running-a-shell-command)
12. [Summary](#summary)

---

## 1. Why Process States Exist

Imagine your computer running 200 processes on 4 CPU cores. Only 4 can actually run at once. What are the other 196 doing?

They're waiting — but waiting for different things:
- Some are waiting for their turn on the CPU
- Some are waiting for disk I/O to complete
- Some are waiting for network data to arrive
- Some are waiting for user input
- Some have finished but their parent hasn't collected their exit code yet

The OS must track each of these different waiting situations, because the correct response to each is different:

- A process waiting for **disk I/O** should not get CPU time until the disk signals it's done
- A process waiting for **CPU** should get on the run queue and be scheduled
- A process that has **exited** should not be scheduled at all, but must preserve its exit code

**Process states** are how the OS classifies exactly what each process is doing and why, so it can manage everything correctly.

---

## 2. The Five Core Process States

Any OS textbook describes these 5 fundamental states:

**1. NEW:**
The process is being created. The kernel is allocating memory, setting up the address space, copying data. The process doesn't yet have all resources ready to run.
- Duration: milliseconds
- Transition: → READY when initialization is complete

**2. READY (Runnable):**
The process has everything it needs. It's waiting for its turn on the CPU.
- Has memory, has its code loaded, has file descriptors open
- Just needs a CPU core
- Could be waiting a few microseconds (lightly loaded system) or many milliseconds (heavily loaded)
- Transition: → RUNNING when the scheduler picks it

**3. RUNNING:**
The process is actively executing on a CPU core. Instructions are being fetched and executed.
- Only ONE process per CPU core can be in this state
- On a 4-core machine, max 4 processes can be RUNNING simultaneously
- Transitions: → READY (preempted by timer), → BLOCKED (needs I/O), → TERMINATED (exits)

**4. BLOCKED (Waiting/Sleeping):**
The process cannot run right now because it's waiting for an event:
- Waiting for disk I/O to complete
- Waiting for network data
- Waiting for a lock to be released
- Waiting for a timer to expire
- Waiting for user input
- Blocked processes do NOT sit in the run queue — they can't use the CPU anyway
- Transition: → READY when the event they're waiting for occurs

**5. TERMINATED (Zombie):**
The process has finished executing (called `exit()` or was killed). But the kernel keeps a "zombie" entry until the parent collects the exit status via `wait()`.
- The process is dead — no more execution
- Kernel still holds: exit code, CPU time used, PID (so parent can call waitpid)
- Memory and resources have been freed
- Transition: → gone (entry removed) when parent calls `wait()`

---

## 3. State Transition Diagram

```
                ┌─────────────────────────────────────────────┐
                │                                             │
                ▼                                             │ fork() creates child
             ┌─────┐                                         │
             │ NEW │──── initialization ────────────────────►│
             └─────┘     complete                            │
                ▓                                             │
                ▓                                             │
                ▼                                             │
           ┌──────────┐    scheduler picks this process       │
           │  READY   │◄────────────────────────────────────  │
           │(Runnable)│                                       │
           └──────────┘                                       │
           ▲        │                                         │
           │        │ scheduler dispatches                    │
           │        ▼                                         │
           │   ┌─────────┐                                    │
           │   │ RUNNING │────────── exit() or killed ───────►│
           │   └─────────┘                                    │
           │      │    │                                   ┌──┴──────────┐
           │      │    │                                   │  TERMINATED  │
   I/O     │      │    │ I/O request,                     │   (Zombie)   │
 completes │      │    │ sleep(), wait(), lock             └─────────────┘
  (IRQ)    │      │    ▼                                        │
           │   ┌──────────┐                                     │ parent calls wait()
           │   │ BLOCKED  │                                     ▼
           └───│ (Waiting)│                                  [gone]
               └──────────┘
                    │
                    │ timer expires (sleep),
                    │ data arrives (network/disk),
                    │ lock released
                    └──────────────────────────────────────────►
                              (transitions to READY)
```

**Key rules:**
- BLOCKED → RUNNING is NOT allowed (must go through READY first)
- RUNNING → BLOCKED happens when the process needs to wait for something
- RUNNING → READY happens when the OS takes the CPU away (preemption)
- NEW → READY happens when `fork()` or process creation is complete

---

## 4. Linux Process States in Detail

Linux has more states than the basic 5. Here's what you actually see in `ps`:

| Letter | State | Description |
|--------|-------|-------------|
| `R` | Running | Currently running or on the run queue (ready) |
| `S` | Sleeping (interruptible) | Waiting for event; can be woken by signal |
| `D` | Sleeping (uninterruptible) | Waiting for I/O; cannot be woken by signal |
| `T` | Stopped | Process is stopped (Ctrl+Z, SIGSTOP) |
| `Z` | Zombie | Process terminated, waiting for parent's wait() |
| `I` | Idle | Kernel thread that is idle |

**The critical distinction: S vs. D**

**S (interruptible sleep):** Process is sleeping but can be woken by a signal.
- Example: `sleep 60` — waiting for 60 seconds timer
- Example: waiting for network data (which might never come)
- If you send SIGTERM, the process wakes up and handles the signal (can terminate gracefully)

**D (uninterruptible sleep):** Process is sleeping and CANNOT be woken by any signal.
- Example: waiting for physical disk I/O to complete
- This is necessary: you can't safely interrupt a disk write in progress
- A process stuck in D state (disk I/O hung) cannot be killed with `kill -9`
- "unkillable" processes in `ps` are always in D state with hung I/O

```bash
$ ps aux | awk '$8 ~ /D/'    # Find processes in uninterruptible sleep
# These are waiting for I/O (disk hung, NFS stale, etc.)
```

---

## 5. Running vs. Runnable — A Common Confusion

Linux shows both RUNNING and RUNNABLE processes with `R` in `ps`. This causes confusion.

**RUNNING:** Process has a CPU core and is actively executing instructions. On a 4-core machine, exactly 4 processes can be in this state.

**RUNNABLE:** Process has everything it needs and is waiting in the run queue for a CPU. Also shows as `R` in Linux.

**Example:**
```bash
$ stress --cpu 16 &    # Start 16 CPU-intensive processes
$ ps aux | grep stress | awk '{print $1, $8}'
# You'll see 16 processes with state R
# But only (number of cores) are actually running
# The rest are RUNNABLE — waiting in the run queue
```

The distinction matters for:
- Performance analysis: high run queue length = CPU bottleneck
- Understanding load average: Linux load average counts both RUNNING and RUNNABLE processes
- `uptime` shows load average: 1.0 on a single-core means the CPU is fully busy

---

## 6. What Causes State Transitions

**NEW → READY:**
- `fork()` system call: kernel creates new process, finishes initialization
- `exec()` finishes loading the new program
- `vfork()`, `clone()` for threads

**READY → RUNNING:**
- The OS scheduler selects this process from the run queue
- Happens on: timer interrupt (preemption of previous process), previous process blocks or exits

**RUNNING → READY (preemption):**
- Timer interrupt fires and scheduler decides to switch (time slice exhausted)
- Higher-priority process becomes runnable (preemptive scheduling)
- Process voluntarily yields with `sched_yield()`

**RUNNING → BLOCKED:**
- Process calls `read()` and data isn't in cache (needs disk I/O)
- Process calls `sleep()` or `nanosleep()`
- Process tries to lock a mutex that's held by another thread
- Process calls `wait()` for a child that hasn't exited yet
- Process calls `select()/poll()/epoll_wait()` waiting for I/O readiness
- Process accesses a page that needs to be loaded from swap (page fault)

**BLOCKED → READY:**
- Disk I/O completes (interrupt from disk controller)
- Network data arrives (interrupt from NIC)
- Timer expires (timer interrupt)
- Another thread releases the mutex
- Child process exits (parent's `wait()` can now return)
- Signal arrives (for interruptible sleep)

**RUNNING → TERMINATED:**
- Process calls `exit()`
- Process receives a fatal signal (SIGSEGV, SIGKILL, SIGTERM with default handler)
- Process is killed with `kill -9`

**TERMINATED → gone:**
- Parent calls `waitpid()` for this process
- If parent already exited: PID 1 calls `wait()` (after reparenting)

---

## 7. The Blocked State — Waiting for Events

The BLOCKED state is the most important state to understand for performance.

**When a process blocks on I/O:**
```
Process calls read(fd, buf, count):
  1. Kernel checks if data is available (in page cache)
  2. If NOT available:
     a. Kernel initiates I/O request to disk
     b. Moves process to blocked state: TASK_INTERRUPTIBLE or TASK_UNINTERRUPTIBLE
     c. Removes process from run queue
     d. Calls scheduler → another process gets the CPU
  3. Disk completes I/O (interrupt fires):
     a. Kernel interrupt handler runs
     b. Copies data to page cache
     c. Moves process back to READY state (adds to run queue)
  4. Scheduler eventually runs the process again
  5. read() returns with the data
```

While the process is blocked on I/O, the CPU is free to run other processes. This is why blocking I/O is acceptable in most programs — you're not wasting CPU, just waiting.

**Types of blocking:**

| What process waits for | How long typically | How kernel tracks it |
|------------------------|-------------------|---------------------|
| Disk I/O (SSD) | 100 μs | Wait queue in block I/O layer |
| Disk I/O (HDD) | 10 ms | Wait queue in block I/O layer |
| Network data | ms to seconds | Socket's wait queue |
| sleep(n) | n seconds | Timer wait queue |
| mutex | μs to ms | Lock's wait queue |
| child exit | seconds to forever | Wait queue in process |
| keyboard input | seconds to forever | tty wait queue |

Each device/subsystem has a **wait queue** — a list of processes sleeping on that event. When the event fires, the kernel walks the wait queue and moves processes to READY.

---

## 8. The Zombie State Revisited

A zombie process:
- Has called `exit()` (code no longer runs)
- Still occupies an entry in the process table (task_struct)
- Cannot be killed by `SIGKILL` (already dead)
- Its resources (memory, file descriptors) have been freed
- Only thing kept: exit code, PID, CPU usage stats

**Why zombies exist:**
The parent may call `wait()` AFTER the child exits. If the kernel immediately destroyed the process entry on `exit()`, the parent couldn't get the exit code.

```c
pid_t child = fork();
if (child == 0) {
    exit(42);       // Child exits immediately
}

// Parent does other work...
sleep(10);          // For 10 seconds, child is a zombie!

int status;
waitpid(child, &status, 0);  // Now child is reaped — zombie gone
printf("Exit code: %d\n", WEXITSTATUS(status));  // prints 42
```

**Zombie accumulation is a bug:**
```bash
$ ps aux | awk '$8 == "Z"'   # Find zombies
```

If you see many zombies, the parent process is not calling `wait()`. It's a programming error in the parent. The fix: use `SIGCHLD` signal handler or `waitpid(-1, NULL, WNOHANG)` in a loop.

---

## 9. Seeing Process States

**Using `ps`:**
```bash
ps aux           # All processes, extended info (BSR format)
ps -ef           # All processes, full-format (standard format)
ps -o pid,stat,cmd   # Custom format: PID, state, command

# Output column meanings:
# S  = state (S=sleeping, R=running, D=disk wait, Z=zombie, T=stopped)
# %CPU = CPU usage percentage
# %MEM = memory usage percentage
```

**Using `top` or `htop`:**
```bash
top         # Dynamic process list, sorts by CPU usage
htop        # Colored, interactive version of top (install separately)
```

**In top's output:**
```
Tasks: 312 total,   2 running, 310 sleeping,   0 stopped,   0 zombie
```
This shows: 2 processes currently running (or runnable), 310 sleeping.

**Examining a specific process:**
```bash
cat /proc/1234/status | grep State
# State: S (sleeping)
# State: R (running)
# State: D (disk sleep)
# State: Z (zombie)
```

**Watching state transitions (advanced):**
```bash
# Watch a process's state in real time:
watch -n 0.5 'cat /proc/1234/status | grep -E "State|VmRSS"'
```

---

## 10. Process Priorities

Not all processes are equal. Some are more important and should get more CPU time.

**Nice values (Unix priority):**
```
-20 = highest priority (runs most often)
  0 = default priority
+19 = lowest priority (runs least often)
```

Confusing name: "nicer" processes are MORE polite = LOWER priority.

```bash
# Run a process with low priority (nice = 10)
nice -n 10 make -j8      # CPU-intensive build, don't disturb others

# Change priority of running process (requires root for negative nice)
renice -n -5 -p 1234     # Increase priority of PID 1234
renice -n 15 -p 5678     # Decrease priority of PID 5678
```

**Real-time priorities (Linux):**
Beyond nice values, Linux supports real-time scheduling policies:
- `SCHED_FIFO`: runs until blocked or pre-empted by HIGHER priority RT process
- `SCHED_RR`: round-robin among same-priority RT processes

Real-time processes have higher priority than ALL normal processes.

```bash
# Check scheduling policy:
chrt -p 1234
# scheduling policy: SCHED_OTHER (normal)
# scheduling priority: 0

# Set real-time (careful — can lock up system if done wrong):
sudo chrt -f -p 90 1234  # set SCHED_FIFO priority 90
```

---

## 11. Process Lifecycle Example — Running a Shell Command

Let's trace the complete lifecycle when you type `ls -la` in a bash shell:

```
You press Enter in terminal
       │
       ▼
bash (PID 100, state: S — sleeping on input)
  wakes up: input event → state: R (running)
  parses command: "ls -la"
  
  calls fork()
       │
       ├──► bash (PID 100) continues
       │    state: S (sleeping, waiting for child)
       │    calls waitpid(child_pid, ...)
       │
       └──► new process (PID 101): state NEW
               kernel sets up PCB, address space
               state → READY
               scheduler runs it: state → RUNNING
               process calls execve("/bin/ls", ["-la"], env)
                 kernel replaces address space with ls binary
                 ls code starts: state still RUNNING
                 ls calls getdents() to read directory
                   kernel checks cache → need disk I/O
                   ls state → BLOCKED (uninterruptible D)
                   scheduler runs other processes
                   disk completes (interrupt)
                   ls state → READY
                   scheduler runs ls: state → RUNNING
                 ls calls write() to stdout
                   kernel writes to terminal buffer
                   (fast, no blocking needed)
                 ls returns from main()
                 exit(0) called
                   ls state → TERMINATED (zombie)
                   kernel frees memory, file descriptors
                   sends SIGCHLD to parent (bash)

bash (PID 100): wakes from waitpid()
  state: R (running)
  collects ls's exit code (0)
  ls zombie → gone
  bash prints new prompt
  bash calls read() waiting for next command
  state: S (sleeping)
  
You see: "$ " prompt again
```

---

## Summary

| State | Meaning | OS action |
|-------|---------|-----------|
| NEW | Being created | Allocate PCB, memory, set up address space |
| READY | Waiting for CPU | Sit in run queue; scheduler will pick it |
| RUNNING | On a CPU, executing | Execute instructions; can be preempted by timer |
| BLOCKED | Waiting for event | Removed from run queue; put on event wait queue |
| TERMINATED | Exited | Keep zombie entry until parent calls wait() |

| Transition | Triggered by |
|------------|-------------|
| NEW → READY | fork() completes |
| READY → RUNNING | Scheduler selects this process |
| RUNNING → READY | Timer interrupt (preemption) or voluntary yield |
| RUNNING → BLOCKED | I/O request, sleep(), mutex, wait() |
| BLOCKED → READY | I/O complete, timer fires, lock released |
| RUNNING → TERMINATED | exit(), SIGKILL, crash |
| TERMINATED → gone | Parent calls wait() |
