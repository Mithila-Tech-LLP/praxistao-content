# Chapter 57: Process Scheduler

> **"The scheduler decides which process runs next. Every operating system decision about fairness, responsiveness, throughput, and latency flows from the scheduler. A simple round-robin that gives each process equal time is enough to make multitasking work — and it's the perfect starting point."**

---

## Table of Contents

1. [What a Scheduler Does](#1-what-a-scheduler-does)
2. [Round-Robin Scheduling](#2-round-robin-scheduling)
3. [Priority and Time Slices](#3-priority-and-time-slices)
4. [When the Scheduler Runs](#4-when-the-scheduler-runs)
5. [The schedule() Function](#5-the-schedule-function)
6. [Voluntary Yield and Sleep](#6-voluntary-yield-and-sleep)
7. [Waking Up Processes](#7-waking-up-processes)
8. [Hooking into the Timer IRQ](#8-hooking-into-the-timer-irq)
9. [Complete scheduler.c / scheduler.h](#9-complete-schedulerh)
10. [Testing the Scheduler](#10-testing-the-scheduler)
11. [Summary](#summary)

---

## 1. What a Scheduler Does

```
The scheduler answers one question: "which process runs next?"

Inputs:
  run_queue:       list of PROC_READY processes
  current_process: currently running process
  timer tick:      periodic interrupt telling us time has passed
  
Outputs:
  A decision: run process X next
  An action:  call switch_to(X) to perform the context switch

Scheduler goals (often in tension):
  Fairness:       every process gets some CPU time
  Throughput:     maximize completed work per second
  Responsiveness: interactive programs feel snappy
  Low overhead:   scheduling decision should be fast (O(1) or O(log n))
  Real-time:      if we had hard RT tasks, meet their deadlines
```

---

## 2. Round-Robin Scheduling

Round-robin is the simplest fair scheduler:

```
Ready queue: [P1] → [P2] → [P3] → [P4]

Each process gets a time slice (e.g., 10ms).
Run P1 for 10ms → move P1 to end → run P2 for 10ms → move to end → ...

Timeline:
|── P1 ──|── P2 ──|── P3 ──|── P4 ──|── P1 ──|── P2 ──| ...
   10ms     10ms     10ms     10ms     10ms     10ms

Properties:
  ✓ No starvation (every process eventually runs)
  ✓ Fair (equal time for equal priority)
  ✓ Simple: O(1) next process selection
  ✗ Not optimal for interactive programs (latency = n × time_slice)
  ✗ Not optimal for batch jobs (doesn't reward short jobs)
```

**Our enhancement: priority-based time slices**

```
Higher priority → larger time slice:
  Priority 1 (lowest):  2 timer ticks
  Priority 5 (medium):  10 timer ticks
  Priority 10 (highest): 20 timer ticks
  
Still round-robin within priority levels, but high-priority processes
get more CPU time before being preempted.
```

---

## 3. Priority and Time Slices

```c
/* How many timer ticks each priority level gets: */
static uint32_t priority_to_ticks(uint32_t priority) {
    return priority * 2;  /* priority 1 → 2 ticks, priority 10 → 20 ticks */
}
```

---

## 4. When the Scheduler Runs

```
The scheduler is invoked in these situations:

1. Timer interrupt (preemptive):
   - IRQ0 fires every timer tick
   - Decrement current process's time_slice
   - If time_slice reaches 0: call schedule()
   
2. Voluntary yield:
   - Process calls yield() or schedule() directly
   - Used when waiting for I/O or wanting to be cooperative
   
3. Process blocks:
   - Process waits for I/O, mutex, semaphore
   - State changes to PROC_BLOCKED
   - Must call schedule() to give CPU to someone else
   
4. Process sleeps:
   - Process calls sleep(ms)
   - State changes to PROC_SLEEPING, sleep_until = now + ms
   - Call schedule()
   
5. Process exits:
   - State changes to PROC_ZOMBIE
   - Call schedule()
```

---

## 5. The schedule() Function

```c
/* kernel/scheduler.c */

#include "scheduler.h"
#include "process.h"
#include "switch.h"
#include "timer.h"
#include "vga.h"

/* Pick the next runnable process (round-robin): */
static process_t *pick_next(void) {
    if (!run_queue) return NULL;
    
    process_t *start = current_process ? current_process->next : run_queue;
    if (!start) start = run_queue;   /* wrap around */
    
    /* Scan from 'start' wrapping around: */
    process_t *candidate = start;
    do {
        if (candidate->state == PROC_READY) {
            return candidate;
        }
        candidate = candidate->next;
        if (!candidate) candidate = run_queue;
    } while (candidate != start);
    
    /* No ready process found — check if current is still running: */
    if (current_process && current_process->state == PROC_RUNNING) {
        return current_process;  /* Keep running */
    }
    
    /* Find idle (PID 0): */
    for (int i = 0; i < MAX_PROCESSES; i++) {
        if (process_table[i].pid == 0) return &process_table[i];
    }
    
    return NULL;  /* Should never happen if idle exists */
}

/* Main scheduling function: */
void schedule(void) {
    /* Disable interrupts during scheduler decision: */
    __asm__ volatile("cli");
    
    /* Check sleeping processes — wake them up if their time has come: */
    uint32_t now = timer_get_ticks();
    for (int i = 0; i < MAX_PROCESSES; i++) {
        process_t *p = &process_table[i];
        if (p->state == PROC_SLEEPING && now >= p->sleep_until) {
            p->state = PROC_READY;
        }
    }
    
    /* Pick next process: */
    process_t *next = pick_next();
    
    if (!next) {
        /* Nothing to run — shouldn't happen if idle process exists */
        __asm__ volatile("sti; hlt");
        return;
    }
    
    if (next == current_process) {
        /* Same process — reset its time slice and continue: */
        if (current_process) {
            current_process->time_slice = priority_to_ticks(current_process->priority);
        }
        __asm__ volatile("sti");
        return;
    }
    
    /* Reset the new process's time slice: */
    next->time_slice = priority_to_ticks(next->priority);
    
    /* Re-enable interrupts AFTER switch (switch_to will enable them): */
    /* Actually: switch_to performs the switch; interrupts are re-enabled
       when the new process resumes (via iret which restores EFLAGS.IF) */
    
    switch_to(next);
    
    /* When we get here, we've been context-switched back in.
       Re-enable interrupts: */
    __asm__ volatile("sti");
}
```

---

## 6. Voluntary Yield and Sleep

```c
/* Voluntarily give up the CPU: */
void yield(void) {
    if (current_process) {
        current_process->time_slice = 0;  /* Force reschedule */
    }
    schedule();
}

/* Sleep for 'ms' milliseconds: */
void sleep_ms(uint32_t ms) {
    if (!current_process) return;
    
    /* Assuming 100 Hz timer: 1 tick = 10ms */
    uint32_t ticks = (ms + 9) / 10;  /* round up */
    
    current_process->state      = PROC_SLEEPING;
    current_process->sleep_until = timer_get_ticks() + ticks;
    
    schedule();  /* Give up CPU until we're woken */
}

/* Block the current process (waiting for an event): */
void block_current(void) {
    if (current_process) {
        current_process->state = PROC_BLOCKED;
    }
    schedule();
}

/* Unblock a process (event has occurred): */
void unblock_process(process_t *proc) {
    if (proc && proc->state == PROC_BLOCKED) {
        proc->state = PROC_READY;
    }
}
```

---

## 7. Waking Up Processes

```c
/* Check and wake sleeping processes — called from timer handler: */
void scheduler_tick(void) {
    if (!current_process) return;
    
    /* Decrement the current process's time slice: */
    current_process->total_ticks++;
    
    if (current_process->time_slice > 0) {
        current_process->time_slice--;
    }
    
    /* Check sleeping processes: */
    uint32_t now = timer_get_ticks();
    for (int i = 0; i < MAX_PROCESSES; i++) {
        if (process_table[i].state == PROC_SLEEPING &&
            now >= process_table[i].sleep_until) {
            process_table[i].state = PROC_READY;
        }
    }
    
    /* If time slice expired, request a reschedule: */
    if (current_process->time_slice == 0) {
        schedule();
    }
}
```

---

## 8. Hooking into the Timer IRQ

Update the timer driver to call the scheduler on each tick:

```c
/* kernel/timer.c — update timer_callback: */

static void timer_callback(registers_t *regs) {
    (void)regs;
    ticks++;
    
    /* Call the scheduler tick handler: */
    scheduler_tick();  /* This may call schedule() and context switch! */
}
```

Now the scheduler is fully preemptive — every timer tick automatically switches to the next process if the time slice expired.

---

## 9. Complete scheduler.h

```c
/* include/scheduler.h */
#pragma once
#include "process.h"

void schedule(void);
void yield(void);
void sleep_ms(uint32_t ms);
void block_current(void);
void unblock_process(process_t *proc);
void scheduler_tick(void);

/* Access process_table from process.c: */
extern process_t process_table[];
```

**Full kernel_main initialization order:**
```c
void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    terminal_init();
    
    struct multiboot_info *mbi = (struct multiboot_info *)mbi_ptr;
    
    kprintf("GDT... ");   gdt_init();                           kprintf("OK\n");
    kprintf("IDT... ");   idt_init();                           kprintf("OK\n");
    kprintf("PIC... ");   pic_init(); pic_disable();            kprintf("OK\n");
    kprintf("PMM... ");   pmm_init(mbi->mmap_addr, mbi->mmap_length, mbi->mem_upper);
    kprintf("VMM... ");   vmm_init();                           kprintf("OK\n");
    kprintf("Heap... ");  heap_init();                          kprintf("OK\n");
    kprintf("Procs... "); process_init();                       kprintf("OK\n");
    kprintf("Timer... "); timer_init(100);                      kprintf("OK\n");
    kprintf("Kbd... ");   keyboard_init();                      kprintf("OK\n");
    
    /* Create test processes: */
    process_create("proc_a", proc_a, 5);
    process_create("proc_b", proc_b, 3);
    
    /* Start scheduling (context switch to first ready process): */
    kprintf("\nStarting scheduler...\n");
    __asm__ volatile("sti");
    schedule();   /* This will switch to the first process */
    
    /* Should never reach here: */
    for (;;) __asm__ volatile("hlt");
}
```

---

## 10. Testing the Scheduler

```c
static volatile int counter_a = 0;
static volatile int counter_b = 0;

static void proc_a(void) {
    while (1) {
        counter_a++;
        if (counter_a % 500 == 0) {
            terminal_set_color(VGA_COLOR_LIGHT_GREEN, VGA_COLOR_BLACK);
            kprintf("[A:%d] ", counter_a);
        }
        /* Yield every iteration to share CPU fairly: */
        yield();
    }
}

static void proc_b(void) {
    while (1) {
        counter_b++;
        if (counter_b % 500 == 0) {
            terminal_set_color(VGA_COLOR_LIGHT_BLUE, VGA_COLOR_BLACK);
            kprintf("[B:%d] ", counter_b);
        }
        yield();
    }
}
```

Expected output (interleaved, showing both processes advancing):
```
Starting scheduler...
[A:500] [B:500] [A:1000] [B:1000] [A:1500] [B:1500] ...
```

If you only see one process's output, context switching is not working.
If both appear but in long runs (1000 before switching), time slices are too long.

**Debugging tip:** Add a tick counter print in `scheduler_tick` to verify IRQ0 is firing:
```c
if (ticks % 100 == 0) kprintf(".");  /* Print a dot every second */
```

---

## Summary

| Concept | Description |
|---------|------------|
| Scheduler | Decides which process runs next |
| Round-robin | Each process gets equal time in rotation; no starvation |
| Time slice | Number of timer ticks a process runs before being preempted |
| Priority | Higher priority → larger time slice per round |
| schedule() | Find next PROC_READY process and call switch_to() |
| yield() | Process voluntarily gives up remaining time slice |
| sleep_ms(ms) | Block for N milliseconds; state = PROC_SLEEPING; sleep_until = now + ticks |
| block_current() | Block process waiting for an event; state = PROC_BLOCKED |
| unblock_process() | Mark a blocked process as PROC_READY |
| scheduler_tick() | Called every timer IRQ; decrements time_slice; triggers reschedule if expired |
| Preemptive | OS can forcibly switch processes without their cooperation (via timer) |
| Cooperative | Process must call yield() or block itself for switching to occur |
| pick_next() | Scan run queue from current process forward; wrap around; skip non-READY |
| Idle process | PID 0; always READY; runs when no other process can; uses hlt |
