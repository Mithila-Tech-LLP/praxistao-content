# Chapter 11: Threads — Lightweight Processes

> **"A process is a house. A thread is a person living in it. Multiple people (threads) can live in the same house (process), sharing rooms, furniture, and the kitchen. But each person has their own bedroom (stack) and can think independently (separate execution flow)."**

---

## Table of Contents

1. [Why Threads Exist](#1-why-threads-exist)
2. [What Is a Thread?](#2-what-is-a-thread)
3. [Process vs. Thread — Side by Side](#3-process-vs-thread--side-by-side)
4. [What Threads Share, What They Don't](#4-what-threads-share-what-they-dont)
5. [Thread States and Lifecycle](#5-thread-states-and-lifecycle)
6. [User Threads vs. Kernel Threads](#6-user-threads-vs-kernel-threads)
7. [Threading Models (1:1, M:1, M:N)](#7-threading-models)
8. [POSIX Threads (pthreads) — The Unix Threading API](#8-posix-threads-pthreads)
9. [Thread Local Storage (TLS)](#9-thread-local-storage-tls)
10. [Context Switching Between Threads](#10-context-switching-between-threads)
11. [Threads in the Linux Kernel](#11-threads-in-the-linux-kernel)
12. [Threads in Other Systems (Java, Go, Windows)](#12-threads-in-other-systems)
13. [Summary](#summary)

---

## 1. Why Threads Exist

**Problem 1: Blocking I/O wastes CPU**
A process reads from the network. While waiting, the CPU is idle. Could be doing other useful work.

**Problem 2: No parallelism within a process**
A CPU-intensive program can only use ONE core if it's single-threaded. On a 16-core machine, 15 cores sit idle.

**Problem 3: Multiple `fork()`s are expensive**
Processes have separate memory spaces. Communication between processes requires IPC (pipes, sockets) — expensive. Memory is not shared — duplicated.

**Threads solve all three:**
1. Thread 1 blocks on network → Thread 2 does computation → CPU stays busy
2. Threads share the process's memory → they can parallelize work on shared data
3. Thread creation is ~10× cheaper than process creation (no new address space)

---

## 2. What Is a Thread?

A **thread** is a unit of execution within a process.

Every process starts with one thread (the "main thread"). It can create additional threads that run concurrently.

Each thread has its own:
- **Program Counter (RIP):** Which instruction this thread is currently executing
- **Stack:** Local variables and function call frames for THIS thread
- **Register state:** All CPU registers (RAX, RBX, ..., RSP, RFLAGS)

All threads in the same process share:
- **Code (text segment):** All threads can call any function in the process
- **Heap:** malloc() from one thread can be free()d by another
- **Global/static variables:** Shared variables are visible to all threads
- **Open file descriptors:** If Thread 1 opens a file, Thread 2 can use the fd
- **Page table (virtual address space):** Same virtual addresses for all threads

```
PROCESS
├── Virtual Address Space (shared by all threads)
│   ├── Text (code)
│   ├── Data (globals)
│   ├── Heap (shared)
│   ├── Shared libraries
│   ├── Thread 1 Stack [private to Thread 1]
│   ├── Thread 2 Stack [private to Thread 2]
│   └── Thread 3 Stack [private to Thread 3]
│
├── Open file descriptors (shared)
│
├── Thread 1: { RIP=0x401100, RSP=stack1_top, regs... }
├── Thread 2: { RIP=0x402200, RSP=stack2_top, regs... }
└── Thread 3: { RIP=0x403300, RSP=stack3_top, regs... }
```

---

## 3. Process vs. Thread — Side by Side

| Feature | Process | Thread |
|---------|---------|--------|
| Memory | Separate address space | Shared address space |
| Communication | IPC (pipes, sockets) — expensive | Shared memory — cheap |
| Creation time | ~1–10 ms (fork + copy) | ~100 μs (just a new stack + TCB) |
| Context switch cost | Higher (flush TLB, new page table) | Lower (same page table, just registers) |
| Isolation | Crash in P1 doesn't affect P2 | Crash in one thread can crash entire process |
| Security | Process A can't read process B | All threads share everything — no isolation |
| Parallelism | Yes (separate processes) | Yes (multiple threads on multiple cores) |
| Synchronization | Not needed (separate memory) | Needed! (shared variables can race) |

---

## 4. What Threads Share, What They Don't

**Shared among all threads in a process:**

```c
// These are shared — ALL threads see the same values:
int global_count = 0;           // global variable
static char buffer[1024];       // static variable

// If Thread 1 does: global_count = 42;
// Thread 2 immediately sees: global_count == 42

// Heap allocations are shared:
char *ptr = malloc(100);  // Thread 1 allocates
// Thread 2 can use `ptr` too! Same address, same memory.
```

**Private per thread:**

```c
// Each thread has its own stack:
void my_function() {
    int local_var = 10;   // local variable — on THIS thread's stack
    // Another thread calling the same function has its OWN local_var!
    // Changes here don't affect other threads.
}

// Thread-local storage (TLS):
__thread int per_thread_counter = 0;  // Each thread has its own copy
```

**Why shared state is both useful and dangerous:**

Useful: Thread 1 adds items to a linked list, Thread 2 processes them. They share the list in the heap.

Dangerous: If both modify the list simultaneously, without coordination → data corruption. This is a **race condition** (covered in Chapter 12).

---

## 5. Thread States and Lifecycle

Each thread has its own state (independent of other threads in the same process):

```
            ┌──────────────┐
            │    READY     │
            │  (runnable)  │
            └──────┬───────┘
   scheduler picks │         timer/preempt
                   ▼         ─────────────►
            ┌──────────────┐
            │   RUNNING    │
            └──────┬───────┘
              │    │    │
              │    │    └── blocking call → BLOCKED
              │    │
              │    └────── pthread_exit() → TERMINATED
              │
              └─ creates new thread: new thread starts READY
```

**One thread blocking doesn't block the whole process:**
This is the key advantage of multi-threading.

```
Thread 1: reads from disk → BLOCKED
Thread 2: continues processing on CPU → RUNNING
Thread 3: waiting for lock → BLOCKED

Process is still "alive" — Thread 2 is running!
```

Compare with single-threaded: read from disk → entire process blocked.

---

## 6. User Threads vs. Kernel Threads

**Kernel thread:**
The OS kernel knows about this thread. It has a PCB/TCB entry in the kernel. The kernel scheduler manages it directly.

- Can run on multiple CPU cores simultaneously
- If one blocks (on I/O, lock, etc.), others can still run
- Context switch requires kernel involvement (~1 μs)
- More expensive to create (~100 μs)

**User thread (green thread):**
Managed entirely in user space by a user-space library. The kernel doesn't know about them — it just sees one or a few kernel threads.

- Faster to create (no kernel call, just allocate a stack)
- Fast context switch (no kernel involvement)
- All user threads on one kernel thread → can only use ONE core
- If any user thread blocks on I/O → ALL user threads blocked (kernel doesn't know there are others)

**Examples of user threads:**
- Go goroutines (multiplexed over kernel threads by the Go runtime)
- Erlang processes (thousands of lightweight processes)
- Python's asyncio tasks (coroutines — not even real threads, cooperative switching)
- Java virtual threads (JDK 21+)

---

## 7. Threading Models

Three ways to map user threads to kernel threads:

**1:1 (One-to-One):**
```
User Thread 1 ─── Kernel Thread 1 ─── CPU Core 1
User Thread 2 ─── Kernel Thread 2 ─── CPU Core 2
User Thread 3 ─── Kernel Thread 3 ─── CPU Core 3
```
- Each user thread is a kernel thread
- True parallelism on multiple cores
- Blocking in one thread doesn't affect others
- Heavier (can't create millions of threads)
- **Used by:** Linux pthreads, Windows threads, macOS

**M:1 (Many-to-One):**
```
User Thread 1 ─┐
User Thread 2 ─┤── Kernel Thread 1 ─── CPU Core 1
User Thread 3 ─┘
```
- Many user threads mapped to one kernel thread
- Fast user-space scheduling
- Only one core used
- One blocking call blocks ALL user threads
- **Used by:** Early Java green threads, very old systems

**M:N (Many-to-Many / Hybrid):**
```
User Thread 1 ─┐     Kernel Thread 1 ─── CPU Core 1
User Thread 2 ─┼──◄──Kernel Thread 2 ─── CPU Core 2
User Thread 3 ─┤     Kernel Thread 3 ─── CPU Core 3
User Thread 4 ─┘
(4 user threads on 3 kernel threads)
```
- Many user threads on many (but potentially fewer) kernel threads
- User-space scheduler decides which user threads run on which kernel threads
- Parallelism up to the number of kernel threads
- Complex to implement correctly
- **Used by:** Go runtime, Erlang, JVM with virtual threads

---

## 8. POSIX Threads (pthreads)

pthreads is the standard threading API on Unix/Linux. Let's learn it properly.

**Creating a thread:**
```c
#include <pthread.h>
#include <stdio.h>
#include <stdlib.h>

// Thread function: takes void*, returns void*
void *worker_function(void *arg) {
    int thread_num = *(int *)arg;
    printf("Thread %d: hello from thread!\n", thread_num);
    
    // Return value (can pass data back to main thread)
    int *result = malloc(sizeof(int));
    *result = thread_num * 100;
    return result;
}

int main() {
    pthread_t thread1, thread2;
    int arg1 = 1, arg2 = 2;
    
    // Create thread 1
    pthread_create(&thread1, NULL, worker_function, &arg1);
    
    // Create thread 2
    pthread_create(&thread2, NULL, worker_function, &arg2);
    
    printf("Main thread: both threads started\n");
    
    // Wait for threads to finish and get return values
    void *ret1, *ret2;
    pthread_join(thread1, &ret1);
    pthread_join(thread2, &ret2);
    
    printf("Thread 1 returned: %d\n", *(int*)ret1);
    printf("Thread 2 returned: %d\n", *(int*)ret2);
    
    free(ret1);
    free(ret2);
    return 0;
}
```

**Compile and run:**
```bash
gcc -o program program.c -pthread
./program
# Output (order varies):
# Main thread: both threads started
# Thread 1: hello from thread!
# Thread 2: hello from thread!
# Thread 1 returned: 100
# Thread 2 returned: 200
```

**Thread creation attributes:**
```c
pthread_attr_t attr;
pthread_attr_init(&attr);

// Set stack size (default is usually 8MB)
pthread_attr_setstacksize(&attr, 256 * 1024);  // 256KB

// Set detached (don't need to join; auto-cleanup on exit)
pthread_attr_setdetachstate(&attr, PTHREAD_CREATE_DETACHED);

pthread_create(&thread, &attr, function, arg);
pthread_attr_destroy(&attr);
```

**Detached vs. Joinable threads:**
- **Joinable** (default): main thread calls `pthread_join()` to wait for it and get return value
- **Detached**: runs independently; resources auto-cleaned when it exits; can't be joined

**Thread termination:**
```c
// Option 1: Return from thread function (normal)
return result;

// Option 2: pthread_exit() — can be called from anywhere, including called functions
pthread_exit(result);

// Option 3: pthread_cancel() — request another thread to terminate
// The canceled thread may not stop immediately (cancellation points)
pthread_cancel(thread_id);

// Main thread exits → ALL threads die with it
// exit() → kills entire process including all threads
```

---

## 9. Thread Local Storage (TLS)

Sometimes you want a variable that looks global but is actually separate for each thread.

**Example use case:** `errno` — the error code from system calls. If Thread 1 calls `read()` and gets ENOENT, and Thread 2 simultaneously calls `open()` and gets EACCES, they must not interfere. `errno` is thread-local.

**C compiler extension:**
```c
__thread int per_thread_counter = 0;   // GCC/Clang

// Each thread has its own counter:
// Thread 1: per_thread_counter = 5  (doesn't affect Thread 2)
// Thread 2: per_thread_counter = 9  (doesn't affect Thread 1)
```

**pthreads manual TLS:**
```c
pthread_key_t key;

// Create a key (once, at startup)
pthread_key_create(&key, free);  // free() will be called when thread exits

// Store a value for this thread
int *data = malloc(sizeof(int));
*data = 42;
pthread_setspecific(key, data);

// Get the value (from any thread — each gets its own)
int *my_data = pthread_getspecific(key);
printf("My data: %d\n", *my_data);
```

---

## 10. Context Switching Between Threads

Context switching between threads in the same process is cheaper than between processes.

**Process context switch:**
1. Save all CPU registers of current thread → memory
2. Switch to kernel stack
3. **Flush TLB (Translation Lookaside Buffer)** ← expensive!
4. Switch CR3 (page table base) to new process
5. Load saved registers of new process
6. Return to user mode

**Thread context switch (same process):**
1. Save all CPU registers of current thread → memory
2. Switch to kernel stack (if scheduler needs to run)
3. **No TLB flush** — same address space!
4. **No CR3 switch** — same page table!
5. Load saved registers of new thread
6. Return to user mode

The TLB flush is the expensive part. Without it, thread switches are 3–10× faster than process switches.

**Approximate costs:**
- Thread context switch: ~1–3 μs
- Process context switch: ~5–20 μs

---

## 11. Threads in the Linux Kernel

In Linux, threads and processes are treated nearly identically. The kernel has no separate "thread" concept — it just has **tasks** (`task_struct`).

`pthread_create()` internally calls `clone()` system call:
```c
// Process creation (fork):
clone(SIGCHLD, 0);  // new address space

// Thread creation (pthread_create):
clone(CLONE_VM | CLONE_FS | CLONE_FILES | CLONE_SIGHAND | ..., stack_addr);
// CLONE_VM: share virtual memory (address space)
// CLONE_FS: share filesystem (cwd, root)
// CLONE_FILES: share file descriptor table
```

**Process vs. thread in Linux from the kernel's view:**
- Both are `task_struct` entries
- Threads have CLONE_VM set → share `mm_struct` (memory descriptor)
- Processes have their own `mm_struct` (separate address space)

**Thread IDs in Linux:**
```c
getpid()   // Returns the THREAD GROUP ID (same for all threads in a process)
gettid()   // Returns the TASK ID (unique per thread)

// Example:
// Main thread: getpid()=1234, gettid()=1234
// Thread 2:    getpid()=1234, gettid()=1235
// Thread 3:    getpid()=1234, gettid()=1236
```

```bash
$ ps -eLf   # Show threads (each task)
UID        PID  PPID   LWP  C NLWP STIME TTY   TIME CMD
user      1234  5678  1234  5    4  ...   pts/1  ... python3 script.py
user      1234  5678  1235  5    4  ...   pts/1  ... python3 script.py
user      1234  5678  1236  5    4  ...   pts/1  ... python3 script.py
user      1234  5678  1237  5    4  ...   pts/1  ... python3 script.py
# LWP = Light Weight Process (Linux thread ID)
# NLWP = Number of threads
```

---

## 12. Threads in Other Systems

**Java Threads:**
```java
// Create a thread
Thread t = new Thread(() -> {
    System.out.println("Hello from thread!");
});
t.start();
t.join();  // Wait for it to finish
```
Modern Java (JDK 21+) supports "virtual threads" — lightweight threads implemented in the JVM (like Go goroutines) that don't require a 1:1 OS thread mapping.

**Go Goroutines:**
```go
// Creating a goroutine (ultra-cheap! ~2KB stack, unlike 8MB for pthreads)
go func() {
    fmt.Println("Hello from goroutine!")
}()

// Go runtime multiplexes goroutines onto OS threads (M:N model)
// Can have millions of goroutines concurrently
```

**Windows Threads (Win32 API):**
```c
HANDLE hThread = CreateThread(
    NULL,           // default security
    0,              // default stack size
    ThreadFunction, // function pointer
    lpParam,        // argument
    0,              // creation flags
    &dwThreadId     // thread ID output
);
WaitForSingleObject(hThread, INFINITE);
CloseHandle(hThread);
```

---

## Summary

| Concept | Definition |
|---------|-----------|
| Thread | Unit of execution within a process; shares memory with other threads |
| Stack | Each thread's private local variable storage |
| Shared state | Heap, globals, file descriptors — all threads in a process see the same |
| Kernel thread | Thread the OS scheduler knows about; can use multiple cores |
| User thread | Thread managed by user-space library; faster but limited |
| 1:1 model | One user thread = one kernel thread (Linux, Windows, macOS) |
| M:N model | M user threads on N kernel threads (Go, Erlang, Java virtual threads) |
| pthread_create | Standard Unix API to create a thread |
| pthread_join | Wait for a thread to finish; get its return value |
| TLS | Thread-local storage: each thread has its own copy of the variable |
| Clone syscall | Linux creates threads with `clone(CLONE_VM|...)` |

**The golden rule of threading:** Threads are efficient but require careful synchronization. Shared mutable state accessed by multiple threads without synchronization leads to bugs that are extremely hard to reproduce and debug. This is why the next two chapters (12 and 13) on concurrency problems and synchronization are critical.
