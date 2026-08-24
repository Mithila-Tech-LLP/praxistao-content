# Chapter 10: CPU Scheduling — Who Runs Next?

> **"The scheduler is the OS's most critical algorithm. A bad scheduler wastes CPU time on the wrong processes. A great scheduler makes the machine feel instantly responsive even under load. It's the difference between a computer that feels snappy and one that feels sluggish."**

---

## Table of Contents

1. [The Scheduling Problem](#1-the-scheduling-problem)
2. [Scheduling Goals and Tradeoffs](#2-scheduling-goals-and-tradeoffs)
3. [When Does Scheduling Happen?](#3-when-does-scheduling-happen)
4. [Non-Preemptive vs. Preemptive Scheduling](#4-non-preemptive-vs-preemptive-scheduling)
5. [Algorithm 1: First-Come, First-Served (FCFS)](#5-algorithm-1-first-come-first-served-fcfs)
6. [Algorithm 2: Shortest Job First (SJF)](#6-algorithm-2-shortest-job-first-sjf)
7. [Algorithm 3: Round Robin (RR)](#7-algorithm-3-round-robin-rr)
8. [Algorithm 4: Priority Scheduling](#8-algorithm-4-priority-scheduling)
9. [Algorithm 5: Multilevel Queue Scheduling](#9-algorithm-5-multilevel-queue-scheduling)
10. [Algorithm 6: Multilevel Feedback Queue (MLFQ)](#10-algorithm-6-multilevel-feedback-queue-mlfq)
11. [Linux's CFS — Completely Fair Scheduler](#11-linuxs-cfs--completely-fair-scheduler)
12. [Windows Scheduling](#12-windows-scheduling)
13. [Real-Time Scheduling](#13-real-time-scheduling)
14. [Multicore Scheduling](#14-multicore-scheduling)
15. [Summary](#summary)

---

## 1. The Scheduling Problem

You have 4 CPU cores. You have 200 runnable processes. Which ones run?

This is the scheduling problem. The OS must decide:
- Which process runs on which CPU core, right now?
- For how long?
- What happens when a high-priority process becomes runnable?
- How do you prevent any process from starving (never getting CPU)?

**Why this is hard:**
Different processes have different needs:
- **Interactive processes** (text editor, browser): need to respond instantly to key presses (<10ms)
- **Batch processes** (video encoding, backups): don't need instant response; want maximum throughput
- **Real-time processes** (audio, video, control systems): need guaranteed response within hard deadlines
- **Background daemons** (system tasks): should not interfere with foreground work

One scheduler must balance all of these simultaneously.

---

## 2. Scheduling Goals and Tradeoffs

A perfect scheduler would optimize all of these:

| Goal | Definition | Who Cares |
|------|-----------|-----------|
| **CPU utilization** | Keep CPU busy as much as possible | Server/batch |
| **Throughput** | Jobs completed per unit time | Batch |
| **Turnaround time** | Total time from submit to finish | Batch |
| **Waiting time** | Time spent in ready queue | All |
| **Response time** | Time from request to first response | Interactive |
| **Fairness** | Each process gets reasonable CPU share | All |
| **Predictability** | Consistent response time | Interactive/RT |

**The fundamental tradeoff:**
- Maximizing throughput → run long jobs without interruption → poor responsiveness
- Maximizing responsiveness → switch frequently → overhead from context switches

Example:
- Context switch overhead: ~1–10 μs (saving/restoring registers, flushing TLB, cache effects)
- If you switch every 1ms, you waste 0.1–1% on overhead — acceptable
- If you switch every 10μs, you waste 10–100% on overhead — terrible

**No single algorithm is best for all workloads.** Real OSes use different strategies for different process classes.

---

## 3. When Does Scheduling Happen?

The scheduler runs in these situations:

1. **Process terminates:** A slot is freed; pick next process
2. **Process blocks:** Waiting for I/O; give CPU to someone else
3. **Timer interrupt fires:** Periodic opportunity to switch
4. **Process yields:** Voluntarily gives up CPU (`sched_yield()`)
5. **New process created:** Should new process run immediately? (maybe if higher priority)
6. **I/O completes:** Blocked process becomes READY; should it preempt current?
7. **Interrupt handled:** After any interrupt, scheduler can choose to switch

The most common trigger is the **timer interrupt** (also called the "scheduler tick"). This happens at a fixed frequency — typically 250 Hz (every 4ms) in modern Linux.

---

## 4. Non-Preemptive vs. Preemptive Scheduling

**Non-preemptive (cooperative):**
A running process keeps the CPU until it voluntarily gives it up (by blocking, yielding, or exiting). The OS cannot forcibly take the CPU.
- Simpler to implement
- Problem: one process can hog the CPU forever (infinite loop → system freeze)
- Used in: early Windows (1.x–3.x), old Mac OS, simple embedded systems

**Preemptive:**
The OS can forcibly take the CPU from a running process via a timer interrupt.
- More complex (need locks everywhere in the kernel)
- Essential for responsiveness and fairness
- Used in: ALL modern general-purpose OSes (Linux, Windows NT, macOS)

---

## 5. Algorithm 1: First-Come, First-Served (FCFS)

**Idea:** Processes run in the order they arrive. Simple queue.

**Example:**
```
Process  Arrival  Burst Time (how long it needs CPU)
P1       0        24ms
P2       1ms      3ms
P3       2ms      3ms

Gantt chart:
[────────── P1 (24ms) ──────────][──P2──][──P3──]
0                              24      27      30

Waiting times:
P1: 0ms   (ran immediately)
P2: 24 - 1 = 23ms   (waited for P1)
P3: 27 - 2 = 25ms   (waited for P1 + P2)
Average waiting time: (0 + 23 + 25) / 3 = 16ms
```

**Problem: Convoy Effect**
Short processes pile up behind long ones. Like a slow truck on a single-lane road.

```
If P2 and P3 had run first:
[P2][P3][─────── P1 ───────]
0   3   6                  30

Average waiting time: (4 + 0 + 1) / 3 = 1.67ms  ← MUCH better!
```

**FCFS is rarely used alone** because of the convoy effect.

---

## 6. Algorithm 2: Shortest Job First (SJF)

**Idea:** Always run the process with the shortest remaining CPU burst.

**Two variants:**
- **Non-preemptive SJF:** Once a process starts, it runs until it blocks or finishes (even if a shorter job arrives)
- **Preemptive SJF (SRTF — Shortest Remaining Time First):** If a new process arrives with a shorter burst than remaining time, preempt current process

**Example (non-preemptive SJF):**
```
Process  Arrival  Burst
P1       0        6ms
P2       2ms      8ms
P3       4ms      7ms
P4       5ms      3ms

At time 0: only P1 available → run P1
At time 6: P2(8), P3(7), P4(3) available → pick P4 (shortest)
At time 9: P2(8), P3(7) available → pick P3
At time 16: only P2 → run P2

Gantt: [──P1──][─P4─][───P3───][────P2────]
       0      6    9         16          24

Waiting times: P1=0, P2=16, P3=9, P4=6 (but P4 arrived at 5, so 6-5=1ms wait)
```

**Why SJF is optimal:** It minimizes average waiting time. Mathematical proof: running shorter jobs first reduces cumulative wait for all others.

**The fundamental problem:** How do you know how long a job will take?
- You can't know for certain
- Solution: predict based on past behavior (exponential averaging)
- `τ_{n+1} = α × t_n + (1-α) × τ_n`
  - `τ_{n+1}` = predicted next burst
  - `t_n` = actual last burst
  - `τ_n` = predicted last burst
  - `α` = weighting factor (0 < α < 1, typically 0.5)

**Not practical for interactive systems** (can't accurately predict burst lengths), but the concept influences real schedulers.

---

## 7. Algorithm 3: Round Robin (RR)

**Idea:** Give each process a fixed time slice (quantum). When the quantum expires, preempt it and put it at the back of the queue.

```
Quantum = 4ms

Processes: P1(24ms), P2(3ms), P3(3ms) all arrive at time 0

Gantt chart:
[P1:4ms][P2:3ms][P3:3ms][P1:4ms][P1:4ms][P1:4ms][P1:4ms][P1:4ms]
0      4      7     10     14     18     22     26     30

P1 runs 4ms, preempted → P2 runs 3ms (done) → P3 runs 3ms (done) → P1 runs 4ms → ... 
P1 total: 6 runs × 4ms = 24ms
```

**Waiting times:**
- P2: 4ms (waited for P1's first quantum)
- P3: 7ms (waited for P1 + P2)
- P1: 10ms + 10ms + 10ms... (gets interleaved)

**Effect of quantum size:**
- **Too small** (e.g., 1ms): Context switch overhead dominates. If context switch takes 1ms and quantum is 1ms, 50% time is wasted on switching!
- **Too large** (e.g., 1000ms): Degenerates into FCFS. Response time becomes terrible.
- **Typical**: 10–100ms quantum. Linux default: ~4ms (with CFS, it's not a fixed quantum but a proportional share).

**Rule of thumb:** Quantum should be larger than 80% of CPU bursts.

**Round Robin is good for:**
- Interactive systems (no process waits too long)
- Fair CPU sharing
- Default scheduling in many systems

---

## 8. Algorithm 4: Priority Scheduling

**Idea:** Each process has a priority. Always run the highest-priority runnable process.

```
Process  Priority  Burst
P1       3         10ms
P2       1 (high)  1ms
P3       4 (low)   2ms
P4       5 (low)   1ms
P5       2         5ms

Order: P2 (pri 1), P5 (pri 2), P1 (pri 3), P3 (pri 4), P4 (pri 5)

Gantt: [P2][──P5──][───P1───][P3][P4]
        0   1      6        16  18  19
```

**The starvation problem:**
Low-priority processes may never run if high-priority processes always arrive.

**Solution: Aging**
Gradually increase the priority of waiting processes. After a process waits for 1 second, increase its priority by 1. After 2 seconds, another +1. Eventually, even the lowest-priority process will reach the highest priority and run.

---

## 9. Algorithm 5: Multilevel Queue Scheduling

**Idea:** Instead of one run queue, have multiple queues for different process types. Each queue has its own scheduling algorithm and inter-queue scheduling.

```
Queue 0: Real-time processes    [FIFO or RR, highest priority]
          ──────────────────────
Queue 1: System processes       [RR with small quantum]
          ──────────────────────
Queue 2: Interactive processes  [RR with medium quantum]
          ──────────────────────
Queue 3: Batch processes        [FCFS, lowest priority]
```

Scheduler always checks Queue 0 first. Only if empty, check Queue 1. And so on.

**Problem:** Processes are classified once at creation and can't move between queues. A batch job that becomes interactive is stuck in the batch queue.

---

## 10. Algorithm 6: Multilevel Feedback Queue (MLFQ)

**The most practical classical scheduling algorithm.** Processes can move between queues based on their behavior.

**Rules:**
1. Multiple queues with different priorities (Q0 highest, Qn lowest)
2. Process starts in the highest-priority queue
3. If a process uses its full quantum → move down to lower queue (it's CPU-bound, give it lower priority)
4. If a process blocks before quantum expires → stay in current queue (it's I/O-bound, interactive — keep it high priority)
5. Periodically (every N ms): boost ALL processes back to highest queue (prevents starvation)

```
Q0: quantum=2ms  [new/recently-blocked processes]
Q1: quantum=4ms  [medium-term CPU users]
Q2: quantum=∞    [CPU-bound batch jobs]

New process P arrives:
  → Q0, gets 2ms
  If uses all 2ms (CPU-bound): → Q1, gets 4ms
  If uses all 4ms: → Q2, runs until blocking
  
P is interactive (keyboard input):
  → Q0, gets 2ms
  Blocks after 0.5ms (waiting for input): stays in Q0
  → Gets response quickly next time
  
Every 1 second: ALL processes → Q0  [aging/boost]
```

**Why MLFQ works well:**
- Interactive processes naturally float to the top (they block quickly)
- CPU-bound processes sink to the bottom (they don't need fast response)
- Adapts automatically without needing prior knowledge

---

## 11. Linux's CFS — Completely Fair Scheduler

Linux uses the **Completely Fair Scheduler (CFS)** since kernel 2.6.23 (2007). It's one of the best production schedulers.

**Core concept: virtual runtime (vruntime)**

Every process tracks how much virtual CPU time it has consumed:
```
vruntime += actual_runtime_ns × (1024 / nice_weight)
```

Where `nice_weight` depends on the process's nice value:
- nice = 0: weight = 1024 (default)
- nice = -5: weight = 2560 (runs ~2.5× more than nice=0)
- nice = +5: weight = 335 (runs ~3× less than nice=0)

**The key invariant:** CFS always runs the process with the LOWEST vruntime.

This ensures that over time, all processes get a proportional share of the CPU.

```
vruntime example with 3 processes (equal nice):
Time 0: P1=0, P2=0, P3=0

Run P1 for 10ms: P1.vruntime = 10, P2=0, P3=0
→ P2 has lowest vruntime, run P2 for 10ms: P2.vruntime = 10, P3=0
→ P3 has lowest, run P3 for 10ms: P3.vruntime = 10
→ All equal, run P1 again...

Result: Each process gets exactly 1/3 of CPU time → "completely fair"
```

**CFS internals:**
- Uses a **red-black tree** (sorted by vruntime) as the run queue
- Insertion/deletion: O(log n) time
- Finding next process: always the leftmost node: O(1)
- On a 4-core machine: one run queue per core

**What CFS doesn't do:**
CFS is designed for throughput and fairness. For hard real-time requirements, Linux has `SCHED_FIFO` and `SCHED_RR` policies that bypass CFS entirely.

**Seeing CFS in action:**
```bash
# Check scheduler stats per process
cat /proc/1234/sched
# nr_voluntary_switches:   1234  (blocked willingly - I/O bound)
# nr_involuntary_switches: 56    (preempted - CPU bound)
# se.vruntime:             123456789  (current vruntime)

# Check overall scheduler stats
cat /proc/schedstat
```

---

## 12. Windows Scheduling

Windows uses a **priority-based preemptive scheduler**.

**Priority levels:**
- 32 priority levels (0–31)
- 16–31: Real-time priorities (REALTIME_PRIORITY_CLASS)
- 1–15: Normal priorities (user processes)
- 0: System idle process

**Thread priority = base priority + dynamic boost**

Windows dynamically boosts priority:
- When a thread wakes from I/O: boost by 1-8 levels
- When foreground window's thread: boost by 2 levels
- When waiting on keyboard/mouse: boost by 6 levels

This makes Windows feel responsive to the user even under load.

**Quantum:**
- Server: 120ms per quantum (throughput-oriented)
- Workstation: 15ms per quantum (interactivity-oriented)

---

## 13. Real-Time Scheduling

For applications with timing guarantees:

**Rate Monotonic (RM):**
Assign priority based on period — shorter period = higher priority. Static assignment.
- Optimal for independent periodic tasks
- Utilization bound: n(2^{1/n} - 1) → approaches ln(2) ≈ 69% as n → ∞

**Earliest Deadline First (EDF):**
Always run the task whose deadline is soonest. Dynamic priority.
- Optimal: can achieve up to 100% utilization
- More complex to implement

**Linux PREEMPT_RT:**
A set of patches making Linux a hard RTOS:
- Converts spinlocks to sleeping locks
- Makes almost all of the kernel preemptible
- Reduces worst-case latency from milliseconds to microseconds

Used in: professional audio, industrial control, robotics.

---

## 14. Multicore Scheduling

Modern machines have multiple CPU cores. This introduces new challenges:

**Per-core run queues (Linux approach):**
Each core has its own run queue. This avoids locking (no shared run queue to contend over).

**Load balancing:**
If one core is idle and another has 10 runnable processes, the OS migrates processes to balance load. But migration has cost: the process's data is no longer in the migrated core's cache (cache cold).

```
Core 0: [P1][P2][P3]    Core 1: []
→ Load balancer migrates P3 from Core 0 to Core 1
→ P3 now on Core 1: cache miss for P3's data
→ But Core 1 was idle: worth it!
```

**NUMA (Non-Uniform Memory Access):**
On servers with multiple CPU sockets, RAM physically attached to Socket 0 is faster for Core 0 than for Core 1. The OS scheduler must be NUMA-aware:
- Prefer to run a process on the same socket as its memory
- Allocate memory on the same socket as the CPU running the process

Linux's scheduler has sophisticated NUMA topology awareness.

**CPU affinity:**
Pin a process to specific CPU cores:
```bash
# Run a process on core 2 and 3 only:
taskset -c 2,3 ./my_program

# Set affinity for running process:
taskset -p -c 0,1 1234
```

Used for real-time applications (avoid migration costs) and to isolate critical processes from interference.

---

## Summary

| Algorithm | When to Use | Advantage | Disadvantage |
|-----------|-------------|-----------|-------------|
| FCFS | Simple batch | Easy to implement | Convoy effect, bad response |
| SJF | Batch, minimal wait | Optimal average wait | Needs burst prediction, starvation |
| Round Robin | General-purpose interactive | Fair, good response | Context switch overhead |
| Priority | Mixed workloads | Urgent tasks run first | Starvation without aging |
| MLFQ | General-purpose | Self-tuning to workload | Complex implementation |
| CFS | Linux desktop/server | Proportional fairness | Not hard real-time |
| EDF | RTOS | Optimal utilization | Complex, needs deadline info |

**Key metrics:**

| Metric | Formula | Who cares |
|--------|---------|-----------|
| Turnaround time | Finish - Submit | Batch jobs |
| Waiting time | Time in ready queue | All |
| Response time | First response - Submit | Interactive |
| Throughput | Jobs / second | Servers |
| CPU utilization | % time CPU is busy | Servers |

**The most important rule:** The scheduler that works best is the one that best matches its workload. There is no universally optimal scheduler.
