# Chapter 11: Goroutines & The GMP Scheduler

The Go scheduler is one of the most commonly probed topics at senior Go interviews. Understanding the GMP model — how goroutines are scheduled across OS threads — explains why Go programs perform the way they do and why goroutines are so cheap.

## Table of Contents

1. [Why Goroutines Instead of Threads](#1-why-goroutines-instead-of-threads)
2. [The GMP Model](#2-the-gmp-model)
3. [Work Stealing](#3-work-stealing)
4. [What Happens When a Goroutine Blocks](#4-what-happens-when-a-goroutine-blocks)
5. [GOMAXPROCS](#5-gomaxprocs)
6. [Goroutine Lifecycle](#6-goroutine-lifecycle)
7. [Preemption](#7-preemption)
8. [Interview Questions & Model Answers](#8-interview-questions--model-answers)
9. [Summary](#summary)

---

## 1. Why Goroutines Instead of Threads

OS threads are expensive. Creating a thread typically costs ~1MB of stack space plus kernel overhead. Goroutines start with only ~2KB of stack and can grow dynamically. This enables Go programs to run hundreds of thousands of goroutines concurrently.

```
OS Thread:
- Stack size:     1-8 MB (fixed)
- Creation time:  ~1 microsecond (syscall required)
- Context switch: ~1-10 microseconds (kernel involved)
- Managed by:     OS scheduler

Goroutine:
- Stack size:     2KB initial (grows and shrinks dynamically up to ~1GB)
- Creation time:  ~100 nanoseconds (no syscall)
- Context switch: ~100 nanoseconds (userspace only)
- Managed by:     Go runtime scheduler
```

The Go runtime multiplexes many goroutines onto a small number of OS threads — typically equal to the number of CPU cores. This is M:N threading.

---

## 2. The GMP Model

The Go scheduler has three entities:

```
G — Goroutine
    A unit of concurrent computation. Has its own stack, program counter,
    and goroutine-local data. Created with the `go` keyword.

M — Machine (OS Thread)
    An actual OS thread. Goroutines run ON an M.
    There are typically GOMAXPROCS active Ms running Go code,
    plus additional Ms that may be blocked in syscalls.

P — Processor (Logical CPU)
    A resource required to run Go code. Holds a run queue of goroutines
    waiting to run. A G needs to be scheduled onto an M with a P to execute.

Relationship:
  P (run queue) → M (OS thread) → G (currently running goroutine)
  
  Each P has a local run queue (LRQ) of up to 256 goroutines.
  There is also a global run queue (GRQ) for overflow.
```

```
  +-------+    +-------+    +-------+
  |  P 0  |    |  P 1  |    |  P 2  |    (number of Ps = GOMAXPROCS)
  |  LRQ  |    |  LRQ  |    |  LRQ  |
  |[G][G] |    |[G][G] |    |[G]    |
  +---↑---+    +---↑---+    +---↑---+
      |             |             |
  +---M---+    +---M---+    +---M---+    (OS threads)
  |  M 0  |    |  M 1  |    |  M 2  |
  | (G 5) |    | (G 3) |    | (G 1) |    (currently running goroutine)
  +-------+    +-------+    +-------+
  
  Global Run Queue (GRQ): [G][G][G]...
```

### How Goroutines Are Scheduled

1. When you call `go f()`, a new goroutine G is created
2. G is placed in the current P's local run queue (LRQ)
3. When the running goroutine yields, blocks, or is preempted, the scheduler picks the next G from the LRQ
4. If LRQ is empty, steal goroutines from other Ps' LRQs or the GRQ

---

## 3. Work Stealing

If one P runs out of goroutines in its LRQ, it steals half the goroutines from another P's LRQ. This keeps all processors busy and avoids idle CPUs.

```
Before steal:
  P0 LRQ: [G1, G2, G3, G4]     P1 LRQ: []

After P1 steals from P0:
  P0 LRQ: [G1, G2]             P1 LRQ: [G3, G4]
```

Work stealing is why Go programs with many small goroutines can efficiently use all CPU cores without the programmer having to think about load distribution.

---

## 4. What Happens When a Goroutine Blocks

This is the crucial insight that makes Go's concurrency model work.

### Case 1: Blocking on a Channel or sync.Mutex

The goroutine is **parked** — removed from its P and placed in a wait queue. The P continues running other goroutines. When the channel becomes ready or mutex is released, the goroutine is moved back to a run queue.

```
Before:                    After G1 blocks on channel receive:
  P0: running G1           P0: running G2 (from LRQ)
                           Wait queue: G1 (waiting for channel)
```

The M is NOT blocked. The P continues running. This is why a thousand goroutines blocking on channels don't consume a thousand OS threads.

### Case 2: Blocking Syscall (e.g., reading a file, network I/O)

When Go knows a syscall will block (via `netpoller` for network I/O, or explicitly for blocking syscalls):

```
Before:                    G1 starts blocking syscall:
  P0: M0 running G1        P0 detaches from M0
                           M0 is now blocked in syscall (running G1)
                           P0 acquires M1 (or creates new M) and continues
                           After syscall completes: G1 put on P's run queue
```

For network I/O specifically, Go uses a non-blocking I/O + polling mechanism (`epoll` on Linux, `kqueue` on macOS). Goroutines waiting for network I/O are parked by the `netpoller` and woken up when data is available — so they never actually block an OS thread.

This is why Node.js-style async/await is unnecessary in Go — the scheduler handles it transparently.

### Case 3: CGo Calls

CGo calls always block an M. This is why heavy CGo usage can cause M proliferation.

---

## 5. GOMAXPROCS

`GOMAXPROCS` sets the number of Ps (logical processors). By default, it equals the number of available CPU cores.

```go
import "runtime"

// Get current GOMAXPROCS
procs := runtime.GOMAXPROCS(0) // 0 means "just query, don't change"

// Set GOMAXPROCS (rarely needed — let the runtime decide)
runtime.GOMAXPROCS(4) // use 4 Ps

// Check CPU count
cpus := runtime.NumCPU()

// Check number of goroutines currently running
goroutines := runtime.NumGoroutine()
```

**When should you change GOMAXPROCS?**

Almost never. The default (= CPU cores) is right for CPU-bound work. For I/O-bound work, Go handles it via the goroutine scheduler — you don't need more Ps. Reducing GOMAXPROCS to 1 can help reproduce race conditions during debugging.

**The GOMAXPROCS=1 container trap:** In Docker/Kubernetes, if the container is limited to 0.5 CPUs, `runtime.NumCPU()` may report the host machine's CPU count (e.g., 96) instead of the container limit. This causes Go to create 96 Ps but only get scheduled on 0.5 CPUs, causing context switching overhead. Fix: use the `automaxprocs` library.

```go
import _ "go.uber.org/automaxprocs" // automatically sets GOMAXPROCS to container quota
```

---

## 6. Goroutine Lifecycle

```go
// CREATED: go keyword creates a goroutine
go func() {
    // RUNNING: goroutine is executing on a P
    
    ch := make(chan int)
    // WAITING: goroutine parked, waiting for channel
    val := <-ch
    
    // RUNNABLE: goroutine is ready to run but waiting for a P
    _ = val
    
    // DEAD: goroutine function returned
}()
```

### Goroutine States

| State | Description |
|---|---|
| Gidle | Just created, not initialized |
| Grunnable | In a run queue, waiting for a P |
| Grunning | Currently executing on a P+M |
| Gsyscall | In a system call |
| Gwaiting | Blocked (channel, mutex, timer, etc.) |
| Gdead | Finished, resources may be reused |

---

## 7. Preemption

Early versions of Go (before 1.14) used **cooperative preemption**: goroutines yielded only at function call sites. A goroutine in a tight loop without function calls could block its M indefinitely.

```go
// This would starve other goroutines in pre-1.14 Go!
go func() {
    for { // tight loop, no function calls, no yield
        i++
    }
}()
```

Go 1.14 introduced **asynchronous preemption**: the runtime can preempt any goroutine at any safe point using OS signals (SIGURG). This ensures all goroutines get CPU time, even in tight loops.

**In interviews:** Mention this evolution. It shows depth in Go version history.

---

## 8. Interview Questions & Model Answers

**Q: Explain the GMP model.**

"G is a goroutine — a lightweight coroutine with ~2KB initial stack. M is an OS thread. P is a logical processor that holds a run queue of goroutines. Go programs have GOMAXPROCS Ps (default = CPU cores). Each P is attached to an M, which runs the G at the front of P's run queue. When a goroutine blocks on a channel or mutex, the scheduler parks it and the P continues with the next goroutine. When a goroutine makes a blocking syscall, the P detaches and gets a new M so other goroutines aren't blocked. This is how one program can have 100,000 goroutines on 8 CPU cores."

**Q: Why are goroutines cheaper than OS threads?**

"Three reasons: smaller initial stack (2KB vs ~1MB), no kernel involvement in creation or context switching (goroutine context switch is ~100ns vs ~10µs for OS threads), and cooperative/signal-based preemption rather than preemptive OS scheduling. The Go runtime manages scheduling in userspace using the GMP model, which avoids expensive system calls."

**Q: What happens when a goroutine makes a system call?**

"For network I/O, Go uses the netpoller — epoll/kqueue under the hood — which parks the goroutine without blocking an M. For blocking syscalls like file reads, the P detaches from the M before the syscall. The M blocks in the kernel, but the P immediately picks up another M (from a pool or newly created) to continue running other goroutines. When the syscall completes, the goroutine becomes runnable and eventually gets scheduled on a P again."

**Q: What is the difference between GOMAXPROCS and the number of goroutines?**

"GOMAXPROCS controls the number of Ps — the maximum number of goroutines that can execute CPU instructions simultaneously. The number of goroutines is unbounded — you can have millions of goroutines while GOMAXPROCS=1 allows only one to run at a time. The other goroutines wait in run queues or are parked. GOMAXPROCS is about CPU parallelism; goroutine count is about concurrent code organization."

---

## Summary

- Goroutines are ~2KB initial stack, created in userspace, ~100ns context switch. Much cheaper than OS threads.
- **GMP:** G=goroutine, M=OS thread, P=logical processor with a run queue.
- Work stealing: idle Ps steal goroutines from busy Ps' run queues.
- Channel/mutex block: goroutine is parked, P continues. No M is blocked.
- Blocking syscall: P detaches from blocked M, acquires a new M to continue.
- Network I/O: handled by netpoller (epoll/kqueue) — goroutines are never OS-blocked.
- GOMAXPROCS = number of Ps = max parallel goroutines. Default = CPU cores.
- Go 1.14+ uses asynchronous preemption via SIGURG — tight loops no longer starve other goroutines.
