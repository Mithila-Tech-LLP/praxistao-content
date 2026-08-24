# Chapter 12: Concurrency Problems

> **"Concurrency bugs are the most treacherous in all of software. A race condition may occur once in a million runs — always when demonstrating to a customer or running in production. It doesn't show up in tests. It disappears when you add debugging statements. The only way to avoid them is to deeply understand WHY they happen."**

---

## Table of Contents

1. [The Fundamental Problem](#1-the-fundamental-problem)
2. [Race Conditions — When Timing Is Everything](#2-race-conditions--when-timing-is-everything)
3. [Atomicity Violations](#3-atomicity-violations)
4. [Order Violations](#4-order-violations)
5. [The Critical Section Problem](#5-the-critical-section-problem)
6. [The "Lost Update" Problem](#6-the-lost-update-problem)
7. [Memory Visibility and Cache Coherency](#7-memory-visibility-and-cache-coherency)
8. [The Producer-Consumer Problem](#8-the-producer-consumer-problem)
9. [The Readers-Writers Problem](#9-the-readers-writers-problem)
10. [The Dining Philosophers Problem](#10-the-dining-philosophers-problem)
11. [Detecting Race Conditions (Tools)](#11-detecting-race-conditions-tools)
12. [Summary](#summary)

---

## 1. The Fundamental Problem

Multiple threads share memory. They can read and write the same variables. When multiple threads access shared mutable state simultaneously, without coordination, the result depends on the exact scheduling order — which is non-deterministic.

**Why scheduling order is non-deterministic:**
- The OS timer fires at arbitrary moments
- Different CPUs run at slightly different speeds
- Cache behavior affects timing
- The exact sequence varies run-to-run

This non-determinism combined with shared state = **concurrency bugs**.

**The key principle:**
A sequence of operations that needs to happen together, without interruption, is called an **atomic operation**. The problem is that most high-level operations are NOT atomic — they're multiple machine instructions.

---

## 2. Race Conditions — When Timing Is Everything

A **race condition** occurs when the outcome of a program depends on the relative order or timing of two or more threads' operations.

**The classic counter example:**

```c
// Shared global variable
int counter = 0;

// Thread 1:                    // Thread 2:
void increment() {              void increment() {
    counter++;                      counter++;
}                               }
```

This looks harmless. But `counter++` is NOT atomic. It's three instructions:
```nasm
1. LOAD:  temp = counter       ; read counter into a register
2. ADD:   temp = temp + 1      ; add 1 to the register
3. STORE: counter = temp       ; write back to counter
```

With two threads, here's what can happen:

```
Time →
Thread 1: [LOAD: temp1=0][ADD: temp1=1]...........[STORE: counter=1]
Thread 2: .............[LOAD: temp2=0][ADD: temp2=1][STORE: counter=1]

Final result: counter = 1   ← WRONG! Should be 2!
```

Both threads read `0`, both add 1, both write `1`. One increment is **lost**.

**Another possible interleaving:**
```
Thread 1: [LOAD: temp1=0][ADD: temp1=1][STORE: counter=1]
Thread 2: ...............................[LOAD: temp2=1][ADD: temp2=2][STORE: counter=2]

Final result: counter = 2   ← CORRECT!
```

The outcome depends on timing. This is a race condition.

**A concrete demonstration in C:**
```c
#include <pthread.h>
#include <stdio.h>

int counter = 0;

void *increment(void *arg) {
    for (int i = 0; i < 1000000; i++) {
        counter++;    // RACE CONDITION!
    }
    return NULL;
}

int main() {
    pthread_t t1, t2;
    pthread_create(&t1, NULL, increment, NULL);
    pthread_create(&t2, NULL, increment, NULL);
    pthread_join(t1, NULL);
    pthread_join(t2, NULL);
    
    printf("Counter: %d\n", counter);  
    // Expected: 2,000,000
    // Actual: somewhere between 1,000,000 and 2,000,000 (varies every run!)
    return 0;
}
```

---

## 3. Atomicity Violations

An **atomicity violation** occurs when a compound operation (which should be atomic) is interrupted by another thread between its steps.

**Example: Check-then-act**
```c
// Thread 1 and Thread 2 both run this:
if (list_is_empty(list)) {
    // At this point, ANOTHER THREAD might add to the list!
    process_empty_list(list);
    // But the list is no longer empty!
}
```

This is a **TOCTOU (Time-of-Check, Time-of-Use)** race condition.

**Example: Singleton initialization**
```c
// Thread 1 and Thread 2 might both call get_singleton():
static Config *config = NULL;

Config *get_singleton() {
    if (config == NULL) {               // Step 1: check
        config = new_config();          // Step 2: allocate
        initialize_config(config);      // Step 3: initialize
    }
    return config;
}
```

Race:
```
Thread 1: config==NULL (true) → allocates → starts initializing (half done)
Thread 2: config!=NULL (partially initialized!) → returns it
Thread 2: uses half-initialized config → crash or wrong behavior!
```

---

## 4. Order Violations

An **order violation** occurs when the assumed execution order between two threads is not guaranteed.

```c
int x;
bool x_initialized = false;

// Thread 1:                    // Thread 2:
void init() {                   void use() {
    x = 42;                         while (!x_initialized) { }
    x_initialized = true;           printf("%d\n", x);
}                               }
```

Seems correct: Thread 2 spins until initialized, then reads `x`.

**But the compiler or CPU may reorder!**

The CPU (for performance) may execute `x_initialized = true` BEFORE `x = 42` — because from a single-thread perspective, the order doesn't matter (they don't depend on each other). Thread 2 then sees `x_initialized = true` but `x` is still garbage!

This requires **memory barriers/fences** to fix.

---

## 5. The Critical Section Problem

A **critical section** is a portion of code that accesses shared resources and must not be executed by more than one thread at a time.

**Formal requirements for a solution:**

1. **Mutual Exclusion:** Only one thread can be in the critical section at a time.
2. **Progress:** If no thread is in the critical section and some threads want to enter, then one of the waiting threads must eventually be allowed to enter. (Can't loop in "who goes first?" forever.)
3. **Bounded Waiting:** A thread waiting to enter the critical section must eventually get in. (No starvation — can't wait forever.)

**The structure of a critical section:**
```
Thread 1, 2, 3, ...:

ENTRY SECTION:
  // acquire the right to enter (might wait here)

CRITICAL SECTION:
  // access shared resource
  // MUST be executed by at most one thread at a time

EXIT SECTION:
  // release the right to enter (let others in)

REMAINDER SECTION:
  // rest of the code (can run concurrently)
```

---

## 6. The "Lost Update" Problem

The lost update is the most common race condition in real applications, especially databases:

**Example: Bank balance**
```
Initial balance: $100

Thread 1 (deposit $50):        Thread 2 (withdraw $30):
  balance = read_balance()        balance = read_balance()
  [both read 100]
  balance += 50                   balance -= 30
  write_balance(balance)          write_balance(balance)
  [Thread 1 writes 150]           [Thread 2 writes 70]

Final balance: $70   ← WRONG! Should be $120
Deposit was LOST.
```

**This is why databases use transactions with locks.** SQL's `BEGIN TRANSACTION / COMMIT` ensures that a read-modify-write sequence is atomic with respect to other transactions.

**Example in code:**
```c
// WRONG:
void deposit(int amount) {
    int balance = get_balance();   // read
    balance += amount;             // modify (not atomic with next step)
    set_balance(balance);          // write
}

// CORRECT (with atomic compare-and-swap or lock):
void deposit(int amount) {
    lock_acquire(&balance_lock);
    int balance = get_balance();
    balance += amount;
    set_balance(balance);
    lock_release(&balance_lock);
}
```

---

## 7. Memory Visibility and Cache Coherency

On a multi-core CPU, each core has its own L1/L2 cache. A write from Core 1 might not immediately be visible to Core 2.

**The cache coherency problem:**
```
Core 1 writes: x = 42;         (written to Core 1's L1 cache)
Core 2 reads:  printf("%d", x); (reads from Core 2's L1 cache → gets old value!)
```

Hardware cache coherency protocols (MESI) eventually propagate updates, but the exact timing is not instantaneous without explicit synchronization.

**Reordering by CPU:**
Modern CPUs execute instructions OUT OF ORDER for performance. Stores can be delayed, loads can be reordered ahead of prior stores. From a single thread's view, this is invisible (the CPU guarantees single-thread consistency). But between threads, another thread may see writes in a different order than written.

**Example:**
```c
// Shared variables
int flag = 0;
int data = 0;

// Thread 1:
data = 42;        // CPU might delay this store
flag = 1;         // flag becomes visible first

// Thread 2:
while (flag == 0) {}  // waits for flag
printf("%d\n", data); // might print 0, not 42!
```

**Solution: Memory barriers / fences**
A memory barrier instruction forces all prior stores to be visible before the barrier completes. Language-level: `std::atomic` in C++, `volatile` (only for ordering, not atomicity) in C.

```c
#include <stdatomic.h>

atomic_int flag = 0;
atomic_int data = 0;

// Thread 1:
atomic_store(&data, 42);    // with memory ordering
atomic_store(&flag, 1);     // visible AFTER data

// Thread 2:
while (atomic_load(&flag) == 0) {}
printf("%d\n", atomic_load(&data));  // guaranteed to see 42
```

---

## 8. The Producer-Consumer Problem

A classic concurrency problem (also called bounded-buffer problem):

**Setup:**
- A fixed-size buffer (like a circular queue)
- Producers add items to the buffer
- Consumers remove items from the buffer
- Producer must wait if buffer is full
- Consumer must wait if buffer is empty

```
Producer → [item] → [  BUFFER  ] → [item] → Consumer
              ↑full: producer waits    ↑empty: consumer waits
```

**The naive (buggy) solution:**
```c
int buffer[N];
int count = 0;  // number of items in buffer

// Producer:
while (true) {
    item = produce_item();
    while (count == N) {}  // busy wait if full ← wastes CPU!
    buffer[in] = item;
    in = (in + 1) % N;
    count++;               // RACE CONDITION!
}

// Consumer:
while (true) {
    while (count == 0) {}  // busy wait if empty ← wastes CPU!
    item = buffer[out];
    out = (out + 1) % N;
    count--;               // RACE CONDITION!
    consume_item(item);
}
```

Problems:
1. `count++` and `count--` are not atomic → race condition
2. Busy-waiting wastes CPU

**The correct solution uses semaphores (Chapter 13).**

---

## 9. The Readers-Writers Problem

**Setup:**
- Shared data (like a database)
- Multiple readers can read simultaneously (reading is safe in parallel)
- Only one writer at a time (writing must be exclusive)
- No reading while writing (writer must have exclusive access)

**Rules:**
1. Multiple readers can read simultaneously
2. Only one writer writes at a time
3. No readers while writing, no writers while reading

**Two classic variants:**

**Readers-first (first readers-writers problem):**
Readers take priority. If readers are reading, new readers can always join. Writers wait until ALL readers finish.
Problem: Writers can starve if there are always readers.

**Writers-first (second readers-writers problem):**
Writers take priority. If a writer is waiting, new readers cannot start.
Problem: Readers can starve.

```c
// Readers-first (conceptual):
int reader_count = 0;
lock_t reader_count_lock, write_lock;

void reader() {
    lock(reader_count_lock);
    reader_count++;
    if (reader_count == 1) lock(write_lock);  // first reader locks writers
    unlock(reader_count_lock);
    
    // READ DATA ...
    
    lock(reader_count_lock);
    reader_count--;
    if (reader_count == 0) unlock(write_lock);  // last reader unlocks writers
    unlock(reader_count_lock);
}

void writer() {
    lock(write_lock);       // exclusive access
    // WRITE DATA ...
    unlock(write_lock);
}
```

---

## 10. The Dining Philosophers Problem

A classic problem that elegantly illustrates deadlock.

**Setup:**
5 philosophers sit around a circular table. Each philosopher alternates between thinking and eating.
5 chopsticks lie between them (one between each pair of adjacent philosophers).
To eat, a philosopher needs BOTH the left AND right chopstick.

```
    Philosopher 0
   /              \
  C4              C0
 /                  \
P4                  P1
|  C3          C1   |
P3                  P2
  \                /
   C3 ------- C2
    Philosopher 3
```

**Naive algorithm (leads to deadlock):**
```
Each philosopher:
  1. Pick up left chopstick
  2. Pick up right chopstick
  3. Eat
  4. Put down right chopstick
  5. Put down left chopstick
  6. Think
  7. Go to 1
```

**Deadlock scenario:**
All 5 philosophers simultaneously pick up their LEFT chopstick.
Now everyone has their left chopstick and waits for the right one.
Nobody puts down their left chopstick. Deadlock! Everyone waits forever.

**Solutions:**
1. **Resource ordering:** Always pick up the lower-numbered chopstick first. Breaks the circular wait.
2. **Asymmetry:** One philosopher picks up right first, others pick up left first.
3. **Arbitration:** A waiter (semaphore/monitor) controls who can pick up chopsticks.
4. **Try-lock:** Pick up left, try to pick up right. If right unavailable, put down left and wait.

The dining philosophers problem models a general class of deadlock problems in real systems (where processes need multiple exclusive resources simultaneously).

---

## 11. Detecting Race Conditions (Tools)

Race conditions are hard to find by reading code. Use automated tools:

**ThreadSanitizer (TSan) — built into GCC/Clang:**
```bash
gcc -fsanitize=thread -g -o program program.c -pthread
./program

# Output if race detected:
# WARNING: ThreadSanitizer: data race (pid=1234)
#   Write of size 4 at 0x0000006020c0 by thread T1:
#     #0 increment /home/user/program.c:5
#   Previous write of size 4 at 0x0000006020c0 by thread T2:
#     #0 increment /home/user/program.c:5
```

**Helgrind — Valgrind tool:**
```bash
valgrind --tool=helgrind ./program
# Detects lock ordering violations, data races
```

**std::atomic in C++ and Java volatile:**
```cpp
// C++: use atomic types for shared variables
#include <atomic>
std::atomic<int> counter(0);

void increment() {
    counter++;  // atomic: LOAD+ADD+STORE as one uninterruptible operation
}
```

**Key takeaway:** Never manually inspect code for race conditions in complex code. Use thread sanitizers. They instrument the code to track every memory access and flag when two threads access the same memory without synchronization.

---

## Summary

| Problem | Cause | Symptom |
|---------|-------|---------|
| Race condition | Non-atomic compound operations on shared state | Wrong values, intermittent bugs |
| Lost update | Read-modify-write without synchronization | Increments lost, wrong totals |
| TOCTOU | Check and use not atomic | "Impossible" states, security holes |
| Order violation | Assumed execution order not guaranteed | Use before initialization, null deref |
| Memory visibility | CPU caches, reordering | Stale values seen across cores |
| Deadlock | Circular wait on resources | Process hangs forever |
| Starvation | One thread never gets access | Process makes no progress |

**The classification of concurrency bugs (Lu et al., ASPLOS 2008 study):**
- 97% of studied non-deadlock bugs: atomicity violations and order violations
- 3%: other patterns

**The root cause of ALL concurrency bugs:**
Shared mutable state accessed by multiple threads without proper synchronization.

**The solution (preview of Chapter 13):**
Synchronization primitives: mutexes, semaphores, condition variables, read-write locks, atomic operations.
