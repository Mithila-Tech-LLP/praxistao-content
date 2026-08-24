# Chapter 14: Deadlocks — When Everyone Is Waiting

> **"A deadlock is a group of processes so thoroughly waiting for each other that none of them can ever continue. It's like two people trying to pass in a narrow hallway, each politely stepping aside — and both stepping to the same side, forever."**

---

## Table of Contents

1. [What Is a Deadlock?](#1-what-is-a-deadlock)
2. [The Four Necessary Conditions](#2-the-four-necessary-conditions)
3. [Resource Allocation Graphs](#3-resource-allocation-graphs)
4. [Types of Deadlock](#4-types-of-deadlock)
5. [Deadlock Prevention](#5-deadlock-prevention)
6. [Deadlock Avoidance — Banker's Algorithm](#6-deadlock-avoidance--bankers-algorithm)
7. [Deadlock Detection](#7-deadlock-detection)
8. [Deadlock Recovery](#8-deadlock-recovery)
9. [Deadlock in Practice — Real-World Examples](#9-deadlock-in-practice--real-world-examples)
10. [Livelock and Starvation](#10-livelock-and-starvation)
11. [Linux's Approach to Deadlock](#11-linuxs-approach-to-deadlock)
12. [Summary](#summary)

---

## 1. What Is a Deadlock?

A **deadlock** is a situation where a set of processes/threads are ALL blocked, EACH waiting for a resource held by ANOTHER in the same set. None can proceed. None can release what others need.

**Simple example:**
```
Thread 1 holds: Lock A
Thread 1 wants: Lock B (held by Thread 2)

Thread 2 holds: Lock B
Thread 2 wants: Lock A (held by Thread 1)

→ Deadlock: both wait forever
```

**Real-world analogy:**
Imagine a one-way road with cars stuck nose-to-tail in a circle. Every car is blocking the car behind it, and waiting for the car ahead to move. Nobody can move. Forever.

**Why deadlocks are serious:**
- The program appears frozen (processes running but nothing happens)
- No error message or crash — just silence
- System resources held forever (memory, files, locks)
- The only recovery is often killing processes or rebooting

---

## 2. The Four Necessary Conditions

Coffman et al. (1971) identified four conditions that must ALL be present simultaneously for a deadlock to occur:

**1. Mutual Exclusion:**
At least one resource must be held in a non-sharable mode — only one thread can use it at a time. Example: a mutex, a file open for writing, a printer.

**2. Hold and Wait:**
A thread is holding at least one resource AND waiting to acquire additional resources held by other threads.

**3. No Preemption:**
Resources cannot be forcibly taken from a thread — they can only be released voluntarily. A thread cannot have its mutex taken away.

**4. Circular Wait:**
A circular chain of two or more threads exists, where each thread is waiting for a resource held by the next thread in the chain:
```
T1 waits for T2's resource
T2 waits for T3's resource
T3 waits for T1's resource  ← forms a cycle
```

**All four conditions must hold simultaneously for deadlock.** Eliminating ANY ONE of them prevents deadlock.

---

## 3. Resource Allocation Graphs

A **Resource Allocation Graph (RAG)** visualizes resource-thread relationships:

**Notation:**
- Circle = thread (process)
- Square = resource
- Arrow from thread → resource: thread is REQUESTING this resource
- Arrow from resource → thread: resource is ASSIGNED to this thread

**No deadlock example:**
```
T1 ──request──► [R1] ──assigned──► T2
T2 ──request──► [R2] ──assigned──► T1

This is:
T1 requests R1 (held by T2)
T2 requests R2 (held by T1)
→ DEADLOCK (cycle: T1→R1→T2→R2→T1)
```

**Key insight:**
- If the resource allocation graph contains NO CYCLE → no deadlock
- If the graph contains a CYCLE:
  - With single-instance resources → deadlock
  - With multi-instance resources → deadlock POSSIBLE but not certain

---

## 4. Types of Deadlock

**Resource deadlock:**
Classic deadlock over physical or logical resources (locks, memory, devices).

**Communication deadlock:**
Processes waiting for messages from each other:
```
Process A: receive from B (blocking)
Process B: receive from A (blocking)
→ Neither will send because both are waiting to receive
```
Common in distributed systems with synchronous communication.

**Database deadlock:**
Transaction 1 holds a row lock for row X, wants row Y.
Transaction 2 holds row Y, wants row X.
Databases detect this automatically and abort one transaction (the "victim").

**Lock-ordering deadlock:**
Most common in multi-threaded programs:
```c
// Thread 1:
lock(mutex_a);
lock(mutex_b);  // → deadlock if Thread 2 locks in reverse order

// Thread 2:
lock(mutex_b);
lock(mutex_a);  // → T2 waits for A (held by T1), T1 waits for B (held by T2)
```

---

## 5. Deadlock Prevention

**Prevention** means designing the system so at least ONE of the four necessary conditions can NEVER hold.

**Attack Condition 1: Mutual Exclusion**
Make resources sharable. But: not everything can be shared (you can't share a mutex — that defeats the purpose). Only applicable for truly read-only resources.

**Attack Condition 2: Hold and Wait**
Ensure that a thread cannot hold any resource when it requests another.

Approach A: Request ALL resources before starting. If any are unavailable, request nothing and retry.
```c
// Thread must request ALL locks at once, or none:
if (trylock(A) && trylock(B)) {
    // have both
} else {
    // release any we got, retry
}
```
Problems: low resource utilization (hold A idle while waiting for B); starvation possible.

Approach B: Release all held resources before requesting a new one.
Often impractical (can't release a half-modified data structure).

**Attack Condition 3: No Preemption**
Allow resources to be forcibly taken.
- If a thread holding resources requests another that's unavailable: preempt it (save its state, release all its resources, restart later).
- Works well for CPU and memory (can be saved/restored).
- Works poorly for locks (partial updates → inconsistent state).

**Attack Condition 4: Circular Wait (Most Practical)**
Impose a TOTAL ORDERING on all resources. Require threads to request resources in ASCENDING order only.

```c
// Define a lock hierarchy:
// LOCK_ORDER: mutex_file (1) < mutex_network (2) < mutex_display (3)
// ALL code must ALWAYS acquire locks in this order

// Thread 1 and Thread 2:
lock(mutex_file);    // step 1: always first
lock(mutex_network); // step 2: always second
// Can never deadlock — both must follow the same order!
```

**Lock ordering is the standard, practical deadlock prevention technique in real systems.**

Example: Linux kernel defines lock ordering (documented in `Documentation/locking/lockdep-design.rst`). Violations are detected by lockdep (lock dependency checker).

---

## 6. Deadlock Avoidance — Banker's Algorithm

**Avoidance** doesn't prevent deadlock in advance — it examines each resource request and dynamically decides whether to grant it based on whether it might lead to deadlock.

**Safe state:** A state is "safe" if there exists a safe sequence — an ordering of all threads such that every thread can get the resources it needs and finish.

**Banker's Algorithm (Dijkstra, 1965):**
Before granting any resource request, the OS checks: if we grant this request, is the resulting state still safe?

**Setup:**
```
n threads, m resource types
Max[i][j]     = max demand of thread i for resource type j
Alloc[i][j]   = currently allocated
Need[i][j]    = Max[i][j] - Alloc[i][j]  (still needed)
Available[j]  = total available of resource type j
```

**Safety algorithm:**
```
Work = Available (copy)
Finish[i] = false for all i

Repeat:
  Find thread i such that:
    Finish[i] == false AND Need[i] <= Work
  If found:
    Work += Alloc[i]  (thread finishes, releases resources)
    Finish[i] = true
  Else:
    Break

If Finish[i] == true for ALL i: SAFE STATE
Else: UNSAFE STATE
```

**Resource request algorithm:**
```
Thread i requests Resources[i]
1. If Request[i] <= Need[i]: OK (valid request)
   Else: error (exceeds declared maximum)
2. If Request[i] <= Available: OK (resources exist)
   Else: wait
3. Pretend we grant: Available -= Request[i], Alloc[i] += Request[i]
4. Run safety algorithm:
   If safe: grant the request
   If unsafe: rollback, make thread wait
```

**Example:**
```
Resources: A(10), B(5), C(7)

Thread  Alloc(A,B,C)  Max(A,B,C)  Need(A,B,C)
T0      (0,1,0)       (7,5,3)     (7,4,3)
T1      (2,0,0)       (3,2,2)     (1,2,2)
T2      (3,0,2)       (9,0,2)     (6,0,0)
T3      (2,1,1)       (2,2,2)     (0,1,1)
T4      (0,0,2)       (4,3,3)     (4,3,1)

Available: (3,3,2)

Safe sequence: T1 → T3 → T4 → T0 → T2
T1 needs (1,2,2): Available (3,3,2) ≥ (1,2,2) → Run T1, release (3,2,2): Available=(6,3,2)
T3 needs (0,1,1): (6,3,2) ≥ (0,1,1) → Run T3, release (2,2,2): Available=(8,5,4)
... etc.
→ State is safe!
```

**Banker's Algorithm limitations:**
- Requires knowing maximum resource needs in advance (often impossible)
- Overhead for every request
- Only works for fixed number of threads and resources
- Not used in most real OSes (too impractical)
- **Used in:** academic study, some specialized systems

---

## 7. Deadlock Detection

Instead of prevention/avoidance, just let deadlocks happen — but detect them and recover.

**For single-instance resources (detect cycle in RAG):**
Run a cycle detection algorithm on the resource allocation graph.
Cycle → deadlock.
Time: O(n²) where n = number of threads.

**For multi-instance resources:**
Similar to Banker's Algorithm but detects existing deadlock, not potential:
```
Work = Available
Finish[i] = (Alloc[i] is all zero)  // threads that have no resources are done

Repeat:
  Find thread i: Finish[i] == false AND Request[i] <= Work
  If found: Work += Alloc[i], Finish[i] = true
  Else: break

If Finish[i] == false for any i: DEADLOCK
Those threads with Finish[i] == false are deadlocked.
```

**When to run detection:**
- After every resource request (expensive but immediate detection)
- Periodically (cheaper, but delayed detection)
- When CPU utilization falls below a threshold (heuristic: low CPU + no progress = suspect deadlock)

---

## 8. Deadlock Recovery

Once detected, recovery options:

**1. Process termination:**
- **Abort all deadlocked processes:** Drastic, but guaranteed to break deadlock. All their work is lost.
- **Abort one at a time, re-run detection:** Less drastic. Abort processes in order of priority, age, how much work done, resources held.

**Which process to abort?**
Choose the one with:
- Lowest priority
- Least time run (minimize wasted work)
- Fewest resources held (minimize impact)
- Not an interactive process (batch is safer to abort)

**2. Resource preemption:**
Select a victim, preempt some of its resources (take away from it), give to deadlocked processes.

**Challenges:**
- Rollback: after preemption, the victim must restart from a consistent state (checkpointing needed)
- Starvation: same process might always be preempted → must bound number of times a process can be the victim

**Database recovery:** Databases use this routinely. When a deadlock is detected, one transaction is rolled back (aborted), its locks released, and it's automatically retried.

---

## 9. Deadlock in Practice — Real-World Examples

**Example 1: Java HashMap (historic bug)**
In Java 1.8 and earlier, `HashMap.get()` could deadlock if multiple threads called `put()` simultaneously (causing infinite loops in the bucket chain). Fixed with `ConcurrentHashMap`.

**Example 2: Linux kernel lockdep findings**
The Linux kernel has a tool called `lockdep` that detects potential lock ordering violations at runtime. It has found hundreds of potential deadlocks in kernel code over the years, all fixed before becoming real bugs.

**Example 3: Database deadlock in production**
```sql
-- Transaction 1:
BEGIN;
UPDATE accounts SET balance = balance - 100 WHERE id = 1;  -- locks row 1
UPDATE accounts SET balance = balance + 100 WHERE id = 2;  -- waits for row 2

-- Transaction 2 (simultaneous):
BEGIN;
UPDATE accounts SET balance = balance - 50 WHERE id = 2;   -- locks row 2
UPDATE accounts SET balance = balance + 50 WHERE id = 1;   -- waits for row 1
-- → DEADLOCK: database detects and rolls back one transaction
```

PostgreSQL, MySQL, Oracle: all detect deadlocks automatically and return error 40P01 (PostgreSQL) or ER_LOCK_DEADLOCK (MySQL).

**Example 4: Network protocol deadlock**
```
Host A: send large data to B, then receive response
Host B: send large data to A, then receive

If both send buffers fill before reading begins → flow control blocks both senders
→ Communication deadlock
```
Fixed by using non-blocking I/O or separate threads for sending and receiving.

---

## 10. Livelock and Starvation

**Livelock:**
Processes are NOT blocked — they're actively running — but they make no forward progress because they keep reacting to each other.

```
Imagine two people in a narrow hallway:
  Person A steps right to let B pass
  Person B steps left to let A pass
  → Both moved to the same side. A steps left, B steps right.
  → Same side again. Both keep moving but never pass.
```

**Example in code:**
```c
// Thread 1 and Thread 2:
while (1) {
    if (try_lock(A)) {
        if (try_lock(B)) {
            break;  // success
        }
        unlock(A);
        sleep(1);   // politely wait
    }
}
```

If Thread 1 and Thread 2 both get A, both fail to get B, both release A, both sleep the same duration, both try again simultaneously... forever.

**Fix for livelock:** Add randomized backoff:
```c
sleep(rand() % MAX_WAIT_MS);  // random sleep breaks the symmetry
```

**Starvation:**
A process that COULD run never actually gets to run because other processes always take priority.

Example: In a priority scheduler, a low-priority process might never run if high-priority processes keep arriving.

**Fix:** Aging — gradually increase priority of waiting processes.

---

## 11. Linux's Approach to Deadlock

Linux uses several strategies:

**1. lockdep — Runtime deadlock detector:**
The kernel can be compiled with `CONFIG_LOCKDEP=y`. It instruments every lock acquisition/release and checks for potential lock ordering violations.

```
kernel: BUG: circular locking dependency detected!
task: systemd/1 is trying to acquire lock:
 (rcu_read_lock){....}, at: [<...>] ...
but task is already holding lock:
 (&mm->mmap_sem){++++}, at: [<...>] ...
which lock already depends on the new lock.
```

**2. Proof by lock ordering:**
Linux documents the strict ordering of its major locks. Any code that deviates is a bug. The lock ordering is maintained in `Documentation/locking/`.

**3. Lock validation in PREEMPT_RT:**
The real-time kernel patch makes all locks preemptible, which helps detect priority inversion issues.

**4. User-space:** Linux applications are responsible for their own deadlock prevention. No OS-level detection for user-space mutexes by default.

---

## Summary

| Strategy | Approach | Cost | Used When |
|----------|---------|------|-----------|
| Prevention | Eliminate one of 4 conditions | Design complexity | All new systems (esp. lock ordering) |
| Avoidance | Banker's algorithm | Runtime overhead per request | Specialized systems, batch |
| Detection | Detect cycles/starvation | Periodic overhead | Databases, some systems |
| Recovery | Abort/preempt | Wasted work, complexity | Databases (rollback), last resort |
| "Ostrich" | Ignore deadlocks | None | Most general-purpose OSes (rare occurrence) |

**The four necessary conditions for deadlock:**
1. Mutual exclusion (non-shareable resource)
2. Hold and wait (holding while requesting)
3. No preemption (can't take resource away)
4. Circular wait (cycle in wait-for graph)

**The practical rule:** In real systems, **lock ordering** (eliminating circular wait) is the most practical prevention strategy. Establish a global lock ordering and enforce it everywhere. This is what both the Linux kernel and most well-designed concurrent applications do.
