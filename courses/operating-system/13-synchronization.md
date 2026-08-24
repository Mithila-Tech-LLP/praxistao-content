# Chapter 13: Synchronization — Locks and Semaphores

> **"A lock is a social contract between threads. It works only because every thread agrees to obey it. One thread that ignores the lock agreement destroys the entire safety guarantee for all threads. Correctness in concurrent programming requires collective discipline, not just individual cleverness."**

---

## Table of Contents

1. [The Goal: Mutual Exclusion](#1-the-goal-mutual-exclusion)
2. [Spinlocks — Busy-Wait Synchronization](#2-spinlocks--busy-wait-synchronization)
3. [Hardware Atomic Instructions](#3-hardware-atomic-instructions)
4. [Mutexes — Sleep-Wait Locking](#4-mutexes--sleep-wait-locking)
5. [Semaphores — Counting Locks](#5-semaphores--counting-locks)
6. [Condition Variables — Wait for a Condition](#6-condition-variables--wait-for-a-condition)
7. [Reader-Writer Locks](#7-reader-writer-locks)
8. [The Happens-Before Relationship](#8-the-happens-before-relationship)
9. [Lock-Free Data Structures](#9-lock-free-data-structures)
10. [Monitor Pattern](#10-monitor-pattern)
11. [Solving the Producer-Consumer Problem](#11-solving-the-producer-consumer-problem)
12. [Common Locking Mistakes](#12-common-locking-mistakes)
13. [Summary](#summary)

---

## 1. The Goal: Mutual Exclusion

The goal of synchronization is **mutual exclusion** — ensuring that at most one thread executes a critical section at a time.

```
Thread 1: ──────[critical section]──────────────
Thread 2: ─────────────────────[waits]──[critical section]──
Thread 3: ──────────────────────────────[waits]──────────[critical section]─
```

Additionally:
- **Progress:** If no thread is in the critical section, a waiting thread must eventually enter (no infinite postponement between threads).
- **Bounded waiting:** A thread waiting to enter will eventually get in (no starvation).

---

## 2. Spinlocks — Busy-Wait Synchronization

A **spinlock** is the simplest lock. A waiting thread spins (loops) checking the lock until it's free.

**Naive (broken) spinlock:**
```c
// This is WRONG:
int lock = 0;

void acquire(int *lock) {
    while (*lock == 1) {}  // spin until lock is free
    *lock = 1;             // take the lock
}

void release(int *lock) {
    *lock = 0;
}
```

**Why it's broken:** The check and set are two separate steps:
```
Thread 1: while(*lock==1) → exits loop (*lock is 0)
Thread 2: while(*lock==1) → exits loop (*lock is 0)  [both exit at same time!]
Thread 1: *lock = 1  ← both think they own the lock!
Thread 2: *lock = 1
```

**Correct spinlock requires an atomic read-modify-write instruction:**

```c
// Using C11 atomics:
#include <stdatomic.h>

atomic_flag lock = ATOMIC_FLAG_INIT;

void spin_lock(atomic_flag *lock) {
    while (atomic_flag_test_and_set(lock)) {
        // spin: test_and_set returns old value
        // if it returns 1: lock was taken, keep spinning
        // if it returns 0: lock was free, we just took it (returns 0 → exits loop)
    }
}

void spin_unlock(atomic_flag *lock) {
    atomic_flag_clear(lock);
}
```

`atomic_flag_test_and_set()` atomically: reads the flag and sets it to 1. Returns the OLD value. This is a single uninterruptible operation — no other thread can sneak in between.

**Spinlock usage:**
```c
atomic_flag counter_lock = ATOMIC_FLAG_INIT;
int counter = 0;

void increment() {
    spin_lock(&counter_lock);
    counter++;           // now safe — only one thread here at a time
    spin_unlock(&counter_lock);
}
```

**When to use spinlocks:**
- **Yes:** Critical sections that are very short (a few nanoseconds). Spinning is cheaper than sleeping if the wait is shorter than the cost of sleeping and waking.
- **Yes:** Kernel code that can't sleep (interrupt handlers, spinlock regions in the kernel).
- **No:** Long critical sections. Spinning wastes CPU time.
- **No:** User-space applications where the critical section might take milliseconds (use mutex instead).

---

## 3. Hardware Atomic Instructions

Hardware provides atomic instructions that are the building blocks of all synchronization:

**Test-and-Set (TAS):**
```c
// Atomically: reads the value, sets it to 1, returns old value
int test_and_set(int *addr) {
    // THIS IS ONE ATOMIC CPU INSTRUCTION — cannot be interrupted
    int old = *addr;
    *addr = 1;
    return old;
}
```

**Compare-and-Swap (CAS):**
```c
// Atomically: if *addr == expected, set *addr = desired, return true
// Otherwise, return false (someone else changed it first)
bool compare_and_swap(int *addr, int expected, int desired) {
    // ATOMIC:
    if (*addr == expected) {
        *addr = desired;
        return true;
    }
    return false;
}
```

CAS is the foundation of most lock-free data structures. The x86 instruction is `CMPXCHG`.

**Fetch-and-Add:**
```c
// Atomically: adds delta to *addr, returns OLD value
int fetch_and_add(int *addr, int delta) {
    // ATOMIC:
    int old = *addr;
    *addr += delta;
    return old;
}
```

x86 instruction: `LOCK XADD`

In C:
```c
#include <stdatomic.h>
atomic_int counter = 0;
atomic_fetch_add(&counter, 1);  // atomic increment
```

---

## 4. Mutexes — Sleep-Wait Locking

A **mutex** (mutual exclusion lock) is a sleeping lock: if the lock is taken, the waiting thread sleeps (is blocked), freeing the CPU for other work.

**Behavior:**
- `lock()`: If unlocked → take the lock. If locked → sleep until unlocked.
- `unlock()`: Release the lock; wake up one waiting thread.

**pthreads mutex:**
```c
#include <pthread.h>

pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;
int counter = 0;

void *increment(void *arg) {
    for (int i = 0; i < 1000000; i++) {
        pthread_mutex_lock(&lock);      // acquire lock (sleeps if taken)
        counter++;                       // critical section — only one thread here
        pthread_mutex_unlock(&lock);    // release lock (wakes waiting threads)
    }
    return NULL;
}
```

**Dynamic initialization:**
```c
pthread_mutex_t lock;
pthread_mutex_init(&lock, NULL);    // initialize
// ... use the lock ...
pthread_mutex_destroy(&lock);       // cleanup
```

**Mutex with attributes:**
```c
pthread_mutexattr_t attr;
pthread_mutexattr_init(&attr);

// Recursive mutex: same thread can lock it multiple times
pthread_mutexattr_settype(&attr, PTHREAD_MUTEX_RECURSIVE);

pthread_mutex_t lock;
pthread_mutex_init(&lock, &attr);
pthread_mutexattr_destroy(&attr);
```

**Mutex vs. Spinlock:**

| Aspect | Spinlock | Mutex |
|--------|---------|-------|
| Wait method | Busy-waits (burns CPU) | Sleeps (frees CPU) |
| Overhead | Very low if no contention | Kernel call required |
| Best for | Very short critical sections | Any length |
| Kernel interrupts | Safe (can't sleep in IRQ) | Can't use in interrupt context |

**RAII pattern in C++ (always unlock):**
```cpp
#include <mutex>

std::mutex lock;
int counter = 0;

void increment() {
    std::lock_guard<std::mutex> guard(lock);  // locks here
    counter++;
    // lock automatically released when guard goes out of scope
}  // ← unlocked here, even if exception thrown
```

---

## 5. Semaphores — Counting Locks

A **semaphore** generalizes a mutex. Instead of a binary (locked/unlocked), it has a non-negative integer counter.

**Operations:**
- `wait()` (also called `P()` or `down()`): Decrement counter. If counter < 0, block.
- `signal()` (also called `V()` or `up()`): Increment counter. If any threads waiting, wake one.

```c
struct semaphore {
    int count;
    queue_t waiting_threads;
};

void wait(semaphore *s) {
    s->count--;
    if (s->count < 0) {
        // block this thread
        enqueue(s->waiting_threads, current_thread);
        sleep_current_thread();
    }
}

void signal(semaphore *s) {
    s->count++;
    if (s->count <= 0) {
        // there are waiting threads (count was negative)
        thread_t *t = dequeue(s->waiting_threads);
        wake_thread(t);
    }
}
```

**Binary semaphore (mutex equivalent):**
Initial value = 1. Acts exactly like a mutex.

**Counting semaphore:**
Initial value = N. Allows N threads into the critical section simultaneously.

**pthreads semaphore:**
```c
#include <semaphore.h>

sem_t semaphore;
sem_init(&semaphore, 0, 1);    // initialize with value 1 (binary semaphore)

sem_wait(&semaphore);           // P operation (decrement, block if 0)
// critical section
sem_post(&semaphore);           // V operation (increment, wake waiter)

sem_destroy(&semaphore);
```

**Semaphore for signaling (not just mutual exclusion):**
```c
// Initial value = 0: Thread B waits until Thread A signals
sem_t done;
sem_init(&done, 0, 0);  // start at 0

// Thread A:
do_work();
sem_post(&done);        // signal Thread B that work is done

// Thread B:
sem_wait(&done);        // block until A posts
use_a_result();
```

This is a powerful use: semaphores as event signaling, not just mutual exclusion.

---

## 6. Condition Variables — Wait for a Condition

A **condition variable** allows a thread to sleep while WAITING FOR A SPECIFIC CONDITION to become true.

Condition variables are ALWAYS used WITH a mutex.

**Pattern:**
```c
// Thread waiting for a condition:
pthread_mutex_lock(&lock);
while (!condition_is_true) {
    pthread_cond_wait(&cond, &lock);
    // This atomically:
    // 1. Releases the lock
    // 2. Sleeps this thread
    // 3. When woken: re-acquires the lock
    // 4. Returns
}
// condition is now true, and we hold the lock
pthread_mutex_unlock(&lock);

// Thread that changes the condition:
pthread_mutex_lock(&lock);
make_condition_true();
pthread_cond_signal(&cond);   // wake ONE waiting thread
// or: pthread_cond_broadcast(&cond);  // wake ALL waiting threads
pthread_mutex_unlock(&lock);
```

**Why `while` not `if` around `cond_wait`?**
Because of **spurious wakeups**: `cond_wait` can return even when no signal was sent (hardware/OS anomaly). The condition MUST be re-checked. Always use `while (!condition)`.

**Example: Thread waiting for work:**
```c
pthread_mutex_t lock = PTHREAD_MUTEX_INITIALIZER;
pthread_cond_t work_cond = PTHREAD_COND_INITIALIZER;
int work_available = 0;

// Worker thread:
void *worker(void *arg) {
    while (1) {
        pthread_mutex_lock(&lock);
        while (!work_available) {
            pthread_cond_wait(&work_cond, &lock);
        }
        work_available = 0;
        pthread_mutex_unlock(&lock);
        do_work();
    }
}

// Producer thread:
void produce_work() {
    pthread_mutex_lock(&lock);
    work_available = 1;
    pthread_cond_signal(&work_cond);
    pthread_mutex_unlock(&lock);
}
```

---

## 7. Reader-Writer Locks

For data that is read frequently but written rarely, we can allow multiple readers simultaneously while ensuring exclusive access for writers.

**pthreads rwlock:**
```c
pthread_rwlock_t rwlock = PTHREAD_RWLOCK_INITIALIZER;

// Multiple readers can run simultaneously:
void read_data() {
    pthread_rwlock_rdlock(&rwlock);   // read lock
    // READ the shared data...
    pthread_rwlock_unlock(&rwlock);
}

// Writers get exclusive access:
void write_data() {
    pthread_rwlock_wrlock(&rwlock);   // write lock (exclusive)
    // WRITE the shared data...
    pthread_rwlock_unlock(&rwlock);
}
```

**When to use:**
- Data that is read 100× more than written (e.g., routing table, configuration)
- Reads are relatively long (reading many fields)

**Trade-off:** More complex, slightly more overhead than a simple mutex. Not always faster — depends on read-write ratio.

---

## 8. The Happens-Before Relationship

**Formal correctness** for concurrent code requires understanding the "happens-before" (HB) relationship.

Operation A "happens-before" operation B means: A's effects are visible to B.

**Rules:**
1. Within a thread: if A comes before B in code order → A HB B
2. `unlock()` HB `lock()` (acquiring a released lock)
3. `signal()` HB `wait()` returns
4. `thread_create()` HB the thread's first instruction
5. A thread's last instruction HB `join()` returns

**Memory ordering in C11:**
```c
#include <stdatomic.h>

atomic_int flag = 0;
int data = 0;

// Thread 1:
data = 42;                              // (A)
atomic_store_explicit(&flag, 1,         // (B) release: A happens-before B
                      memory_order_release);

// Thread 2:
while (atomic_load_explicit(&flag,      // (C) acquire: B happens-before C
                             memory_order_acquire) == 0) {}
printf("%d\n", data);                  // (D) C happens-before D → sees 42
```

The `release` store and `acquire` load create a happens-before: everything before the release is visible after the acquire.

---

## 9. Lock-Free Data Structures

**Lock-free** means progress is guaranteed even if individual threads are preempted — at least one thread makes progress.

The key building block is Compare-and-Swap (CAS):

**Lock-free stack (push):**
```c
typedef struct Node {
    int value;
    struct Node *next;
} Node;

atomic_uintptr_t top = 0;  // atomic pointer to top of stack

void push(int value) {
    Node *new_node = malloc(sizeof(Node));
    new_node->value = value;
    
    Node *old_top;
    do {
        old_top = (Node *)atomic_load(&top);
        new_node->next = old_top;
    } while (!atomic_compare_exchange_weak(&top, 
              (uintptr_t *)&old_top, 
              (uintptr_t)new_node));
    // CAS: if top is still old_top, set it to new_node
    // If another thread changed top, retry (the loop)
}
```

**ABA problem:** Lock-free code using CAS suffers from the ABA problem: A reads value A, B changes A→B→A, A's CAS succeeds even though the structure has changed. Solution: add a version/stamp counter to the value.

**Lock-free vs. mutex tradeoffs:**
- Lock-free can be faster under high contention (no thread blocking)
- Lock-free is MUCH harder to write correctly
- Lock-free doesn't prevent starvation (wait-free algorithms do)
- For most cases, a well-tuned mutex is simpler and fast enough

---

## 10. Monitor Pattern

A **monitor** combines a mutex with condition variables into a unified object:

```c
// Java-style monitor (language built-in):
class Buffer {
    int[] buf = new int[N];
    int count = 0;
    
    synchronized void put(int item) {    // implicitly acquires 'this' lock
        while (count == N) wait();       // wait if full
        buf[count++] = item;
        notifyAll();                     // wake any waiting readers
    }
    
    synchronized int get() {
        while (count == 0) wait();       // wait if empty
        int item = buf[--count];
        notifyAll();                     // wake any waiting writers
        return item;
    }
}
```

`synchronized` is syntactic sugar for mutex lock/unlock. `wait()` and `notifyAll()` are condition variable operations.

C++ equivalents:
```cpp
#include <mutex>
#include <condition_variable>

std::mutex m;
std::condition_variable not_full, not_empty;
```

---

## 11. Solving the Producer-Consumer Problem

Using semaphores (the textbook solution):

```c
#include <semaphore.h>

#define N 10
int buffer[N];
int in = 0, out = 0;

sem_t empty;  // counts empty slots (init = N)
sem_t full;   // counts full slots (init = 0)
sem_t mutex;  // mutual exclusion (init = 1)

// Initialize:
sem_init(&empty, 0, N);  // N empty slots
sem_init(&full, 0, 0);   // 0 full slots
sem_init(&mutex, 0, 1);  // binary mutex

// Producer:
void producer() {
    while (1) {
        int item = produce_item();
        
        sem_wait(&empty);    // decrement empty count; block if 0 (buffer full)
        sem_wait(&mutex);    // enter critical section
        
        buffer[in] = item;
        in = (in + 1) % N;
        
        sem_post(&mutex);    // leave critical section
        sem_post(&full);     // increment full count; wake waiting consumers
    }
}

// Consumer:
void consumer() {
    while (1) {
        sem_wait(&full);     // decrement full count; block if 0 (buffer empty)
        sem_wait(&mutex);    // enter critical section
        
        int item = buffer[out];
        out = (out + 1) % N;
        
        sem_post(&mutex);    // leave critical section
        sem_post(&empty);    // increment empty count; wake waiting producers
        
        consume_item(item);
    }
}
```

**Why the order of sem_wait(&empty) and sem_wait(&mutex) matters:**
If you swap them (mutex first, then empty), you can deadlock:
- Producer holds mutex, waits on empty
- Consumer waits on mutex (can't enter critical section to post empty)
- → Deadlock

Always: acquire counting semaphore BEFORE the mutex.

---

## 12. Common Locking Mistakes

**1. Forgetting to unlock (lock leaking):**
```c
void process() {
    pthread_mutex_lock(&lock);
    
    if (error_condition) {
        return;          // BUG: lock not released!
    }
    
    // process...
    pthread_mutex_unlock(&lock);
}
// Fix: unlock before every return, or use cleanup handlers
```

**2. Lock ordering inconsistency (deadlock):**
```c
// Thread 1:
lock(A); lock(B); // ...

// Thread 2:
lock(B); lock(A); // ...

// Deadlock: T1 holds A waiting for B, T2 holds B waiting for A
// Fix: always lock in the SAME ORDER: A then B, in EVERY thread
```

**3. Holding a lock while sleeping:**
```c
pthread_mutex_lock(&lock);
sleep(10);       // BUG: holding the lock for 10 seconds!
                 // No other thread can progress
pthread_mutex_unlock(&lock);
// Fix: use condition variables to sleep WITHOUT holding the lock
```

**4. Using destroyed/uninitialized locks:**
```c
pthread_mutex_t *lock = malloc(sizeof(pthread_mutex_t));
pthread_mutex_lock(lock);   // BUG: not initialized!
// Fix: always initialize before use
pthread_mutex_init(lock, NULL);
```

**5. Double-locking a non-recursive mutex (self-deadlock):**
```c
void foo() {
    pthread_mutex_lock(&lock);
    bar();                    // bar() also locks the same lock!
    pthread_mutex_unlock(&lock);
}

void bar() {
    pthread_mutex_lock(&lock);  // BUG: self-deadlock if called from foo()
    // ...
    pthread_mutex_unlock(&lock);
}
```

---

## Summary

| Mechanism | Type | Best For |
|-----------|------|---------|
| Spinlock | Busy-wait | Very short critical sections, kernel code |
| Mutex | Sleep-wait | General-purpose mutual exclusion |
| Semaphore | Count + signal | Resource counting, signaling between threads |
| Condition variable | Event wait | Waiting for a condition (used with mutex) |
| RW Lock | Read/write split | Read-mostly shared data |
| Atomic operations | Hardware | Simple counter updates, lock-free algorithms |

**The golden rules of locking:**
1. Always initialize before use
2. Always unlock for every lock, even on error paths
3. Minimize time spent holding a lock
4. Never hold a lock while sleeping
5. Acquire locks in a consistent order (prevents deadlock)
6. Use condition variables to wait for conditions, not spin loops
7. Prefer existing library patterns (producer-consumer, monitor) over DIY
