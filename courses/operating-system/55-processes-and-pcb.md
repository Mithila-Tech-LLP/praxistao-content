# Chapter 55: Processes and the PCB

> **"A process is the operating system's unit of work. It is not just code — it is code plus state: registers, memory, open files, identity, priority. The Process Control Block is the kernel's data structure for capturing all of that state. When you understand the PCB, you understand what 'a running program' truly means."**

---

## Table of Contents

1. [What Is a Process?](#1-what-is-a-process)
2. [Process States](#2-process-states)
3. [The Process Control Block — PCB](#3-the-process-control-block--pcb)
4. [Process IDs and the Process Table](#4-process-ids-and-the-process-table)
5. [Creating a Process](#5-creating-a-process)
6. [Process Stack Layout](#6-process-stack-layout)
7. [The Process List](#7-the-process-list)
8. [Complete process.c / process.h](#8-complete-processh-summary)
9. [Testing Process Creation](#9-testing-process-creation)
10. [Summary](#summary)

---

## 1. What Is a Process?

```
A process is a running instance of a program. It consists of:

  Virtual address space:
    Code segment (.text)      — the program's instructions
    Data segment (.data/.bss) — global/static variables
    Heap                      — dynamically allocated memory
    Stack                     — function calls, local variables
    Kernel-mapped region      — kernel code/data (Ring 0 only)
    
  CPU state (saved when not running):
    General-purpose registers (EAX, EBX, ... EDI)
    Stack pointer (ESP)
    Instruction pointer (EIP) — where to resume
    EFLAGS
    
  Kernel metadata:
    Process ID (PID)
    State: running, ready, blocked, sleeping, zombie
    Priority
    Parent PID
    Open file descriptors
    Working directory
    Accumulated CPU time
    
Two processes can run the same program (e.g., two bash instances)
but each has its own address space, stack, registers, and state.
They don't interfere with each other.
```

---

## 2. Process States

```
State machine:

         fork/exec
            │
            ▼
         [NEW] ──────────── kernel creates PCB
            │
            ▼
         [READY] ◄─────────── runnable; waiting for CPU
            │ ▲
  scheduled │ │ preempted / yield
            ▼ │
         [RUNNING] ──────── currently executing on CPU
            │
            ├──── I/O or sleep ──► [BLOCKED/SLEEPING]
            │                           │
            │                    event/wake ──► [READY]
            │
            └──── exit() ──────► [ZOMBIE]
                                      │
                                 parent wait() ──► freed

Our initial states:
  PROC_READY   = 0    Runnable, waiting for CPU
  PROC_RUNNING = 1    Currently executing
  PROC_BLOCKED = 2    Waiting for event (not yet scheduled)
  PROC_SLEEPING = 3   Waiting for timer
  PROC_ZOMBIE  = 4    Exited but not yet reaped by parent
  PROC_DEAD    = 5    Resources freed, slot available for reuse
```

---

## 3. The Process Control Block — PCB

The PCB holds everything the kernel needs to know about a process:

```c
/* include/process.h */
#pragma once
#include "stdint.h"

#define MAX_PROCESSES    64
#define KERNEL_STACK_SIZE 8192   /* 8KB kernel stack per process */
#define USER_STACK_SIZE  16384   /* 16KB user stack */
#define PROC_NAME_MAX    32

/* Process states: */
typedef enum {
    PROC_DEAD    = 0,
    PROC_READY   = 1,
    PROC_RUNNING = 2,
    PROC_BLOCKED = 3,
    PROC_SLEEPING = 4,
    PROC_ZOMBIE  = 5,
} proc_state_t;

/* Saved CPU registers (what we push/pop during context switch): */
typedef struct {
    /* Saved by our context switch code: */
    uint32_t edi, esi, ebx, ebp;   /* callee-saved registers */
    uint32_t eip;                   /* where to resume */
    /* Note: ESP is saved implicitly as pointer to this struct */
} cpu_state_t;

/* The PCB (Process Control Block): */
typedef struct process {
    /* Identity: */
    uint32_t    pid;
    uint32_t    ppid;           /* Parent PID */
    char        name[PROC_NAME_MAX];
    
    /* State: */
    proc_state_t state;
    int         exit_code;      /* Set when process exits */
    
    /* Scheduling: */
    uint32_t    priority;       /* Higher = more important (we use 1-10) */
    uint32_t    time_slice;     /* Remaining ticks this round */
    uint32_t    total_ticks;    /* Total CPU ticks used */
    uint32_t    sleep_until;    /* Timer tick to wake up (if PROC_SLEEPING) */
    
    /* Memory: */
    uint32_t   *page_dir;       /* Physical address of page directory */
    uint32_t    kernel_stack;   /* Bottom of kernel stack (virtual) */
    uint32_t    kernel_esp;     /* Current kernel stack pointer (when not running) */
    uint32_t    user_stack;     /* Bottom of user stack (virtual) */
    uint32_t    brk;            /* Current program break (heap top) */
    
    /* CPU state (saved when not running): */
    cpu_state_t cpu;
    
    /* Linked list for scheduler: */
    struct process *next;
    struct process *prev;
} process_t;
```

---

## 4. Process IDs and the Process Table

```c
/* kernel/process.c */

#include "process.h"
#include "pmm.h"
#include "vmm.h"
#include "heap.h"
#include "string.h"
#include "vga.h"
#include "timer.h"

/* Process table: fixed-size array */
static process_t process_table[MAX_PROCESSES];
static uint32_t  next_pid = 1;         /* PID 0 = idle, PIDs start at 1 */

/* Currently running process: */
process_t *current_process = NULL;

/* Scheduler run queue head: */
process_t *run_queue = NULL;

/* Get a free PCB slot: */
static process_t *pcb_alloc(void) {
    for (int i = 0; i < MAX_PROCESSES; i++) {
        if (process_table[i].state == PROC_DEAD) {
            memset(&process_table[i], 0, sizeof(process_t));
            return &process_table[i];
        }
    }
    return NULL; /* No free slots */
}

/* Initialize the process subsystem: */
void process_init(void) {
    memset(process_table, 0, sizeof(process_table));
    /* All slots start as PROC_DEAD (0 = dead): */
    current_process = NULL;
    run_queue       = NULL;
    kprintf("Process subsystem initialized.\n");
}
```

---

## 5. Creating a Process

```c
/* Create a new kernel process: */
process_t *process_create(const char *name, void (*entry_point)(void),
                           uint32_t priority) {
    /* 1. Allocate a PCB: */
    process_t *proc = pcb_alloc();
    if (!proc) {
        kprintf("process_create: no free PCB slots!\n");
        return NULL;
    }
    
    /* 2. Set identity: */
    proc->pid   = next_pid++;
    proc->ppid  = current_process ? current_process->pid : 0;
    proc->state = PROC_READY;
    proc->priority   = priority;
    proc->time_slice = priority * 2;  /* Higher priority → more CPU time */
    
    /* Copy name: */
    int i;
    for (i = 0; name[i] && i < PROC_NAME_MAX - 1; i++)
        proc->name[i] = name[i];
    proc->name[i] = '\0';
    
    /* 3. Create address space:
       For kernel processes: use kernel page directory (shared)
       For user processes:   create new page directory (Chapter 59) */
    proc->page_dir = kernel_page_dir;  /* kernel thread shares kernel mapping */
    
    /* 4. Allocate kernel stack: */
    proc->kernel_stack = (uint32_t)kmalloc(KERNEL_STACK_SIZE);
    if (!proc->kernel_stack) {
        proc->state = PROC_DEAD;
        return NULL;
    }
    
    /* Stack grows downward; kernel_esp points to TOP of stack: */
    uint32_t *stack_top = (uint32_t *)(proc->kernel_stack + KERNEL_STACK_SIZE);
    
    /* 5. Set up initial stack frame so context switch can resume at entry_point:
    
       When the scheduler first switches to this process, it pops:
         EDI, ESI, EBX, EBP (callee-saved registers)
         EIP  (return address = entry_point)
       
       We pre-populate the stack with these values: */
    
    *(--stack_top) = (uint32_t)entry_point; /* EIP: where to start */
    *(--stack_top) = 0;                     /* EBP: initial frame pointer */
    *(--stack_top) = 0;                     /* EBX */
    *(--stack_top) = 0;                     /* ESI */
    *(--stack_top) = 0;                     /* EDI */
    
    /* Save the stack pointer: */
    proc->kernel_esp = (uint32_t)stack_top;
    
    /* 6. Add to run queue: */
    proc->next = run_queue;
    proc->prev = NULL;
    if (run_queue) run_queue->prev = proc;
    run_queue  = proc;
    
    kprintf("Created process '%s' (PID=%u, priority=%u)\n",
            proc->name, proc->pid, proc->priority);
    
    return proc;
}

/* Terminate a process: */
void process_exit(int exit_code) {
    if (!current_process) return;
    
    current_process->state     = PROC_ZOMBIE;
    current_process->exit_code = exit_code;
    
    kprintf("Process '%s' (PID=%u) exited with code %d\n",
            current_process->name, current_process->pid, exit_code);
    
    /* Force a reschedule (scheduler will pick someone else): */
    schedule();
}
```

---

## 6. Process Stack Layout

Understanding the initial stack setup is key to context switching:

```
When process_create() returns, the kernel stack looks like:
(high address = bottom of stack allocation)

kernel_stack + KERNEL_STACK_SIZE  ←── bottom of allocated block (high address)
                                        (never used — just the allocation end)
kernel_esp ─────────────────────►  [EDI = 0        ]
                                   [ESI = 0        ]
                                   [EBX = 0        ]
                                   [EBP = 0        ]
                                   [EIP = entry_point]
kernel_stack ────────────────────  (top of allocation = low address)
                                   (unused space above ESP)

When scheduler does context switch (pop EDI, ESI, EBX, EBP, then ret):
  pop edi ← 0
  pop esi ← 0
  pop ebx ← 0
  pop ebp ← 0
  ret     ← pops EIP = entry_point, jumps there
  
The process starts running at entry_point() with a clean register state!
```

---

## 7. The Process List

```c
/* Find a process by PID: */
process_t *process_get(uint32_t pid) {
    for (int i = 0; i < MAX_PROCESSES; i++) {
        if (process_table[i].state != PROC_DEAD &&
            process_table[i].pid == pid) {
            return &process_table[i];
        }
    }
    return NULL;
}

/* Remove a process from the run queue: */
void run_queue_remove(process_t *proc) {
    if (proc->prev) proc->prev->next = proc->next;
    else            run_queue        = proc->next;
    
    if (proc->next) proc->next->prev = proc->prev;
    
    proc->next = proc->prev = NULL;
}

/* Print all processes (like 'ps' command): */
void process_print_all(void) {
    kprintf("PID  PPID  STATE    PRI  TICKS  NAME\n");
    kprintf("---- ----  -------  ---  -----  ----\n");
    
    static const char *state_names[] = {
        "DEAD   ", "READY  ", "RUNNING", "BLOCKED",
        "SLEEP  ", "ZOMBIE "
    };
    
    for (int i = 0; i < MAX_PROCESSES; i++) {
        process_t *p = &process_table[i];
        if (p->state == PROC_DEAD) continue;
        kprintf("%4u %4u  %s  %3u  %5u  %s\n",
                p->pid, p->ppid,
                state_names[p->state],
                p->priority,
                p->total_ticks,
                p->name);
    }
}
```

---

## 8. Complete process.h Summary

```c
/* include/process.h (additions) */

extern process_t *current_process;
extern process_t *run_queue;
extern uint32_t   kernel_page_dir[];

void       process_init(void);
process_t *process_create(const char *name, void (*entry)(void), uint32_t priority);
void       process_exit(int exit_code);
process_t *process_get(uint32_t pid);
void       run_queue_remove(process_t *proc);
void       process_print_all(void);
void       schedule(void);  /* Defined in Ch 57 */
```

---

## 9. Testing Process Creation

```c
/* Test processes: */
static void proc_a(void) {
    kprintf("[Process A] Hello from PID %u!\n", current_process->pid);
    kprintf("[Process A] Exiting.\n");
    process_exit(0);
}

static void proc_b(void) {
    kprintf("[Process B] Hello from PID %u!\n", current_process->pid);
    kprintf("[Process B] Exiting.\n");
    process_exit(0);
}

void kernel_main(...) {
    /* ... init ... */
    
    process_init();
    
    process_t *a = process_create("proc_a", proc_a, 5);
    process_t *b = process_create("proc_b", proc_b, 3);
    
    kprintf("\nProcess table:\n");
    process_print_all();
    
    /* The scheduler (Ch 57) will actually run these processes.
       For now, we just verify they're created: */
    kprintf("\nPID %u created: '%s' at entry 0x%x\n",
            a->pid, a->name, a->cpu.eip);
    
    for (;;) {}
}
```

Output (before scheduler is implemented):
```
Process subsystem initialized.
Created process 'proc_a' (PID=1, priority=5)
Created process 'proc_b' (PID=2, priority=3)

Process table:
PID  PPID  STATE    PRI  TICKS  NAME
---- ----  -------  ---  -----  ----
   1    0  READY      5      0  proc_a
   2    0  READY      3      0  proc_b

PID 1 created: 'proc_a' at entry 0x[...]
```

---

## Summary

| Concept | Description |
|---------|------------|
| Process | Running program instance: address space + CPU state + kernel metadata |
| PCB | Process Control Block: kernel struct holding all process info |
| PID | Process ID: unique integer identifier |
| Process states | Dead → Ready → Running ↔ Blocked/Sleeping → Zombie |
| cpu_state_t | Saved registers (EDI/ESI/EBX/EBP/EIP) restored on context switch |
| Kernel stack | Per-process private stack used when running kernel code on behalf of the process |
| Initial stack frame | Pre-populated on creation so context switch can "resume" at entry_point |
| Run queue | Linked list of PROC_READY processes waiting for CPU time |
| Priority | 1-10; higher → more CPU; affects time_slice |
| time_slice | Ticks remaining in current CPU burst before forced preemption |
| PROC_ZOMBIE | Process exited but parent hasn't called wait() yet; PCB slot still held |
| process_table | Fixed-size global array of PCBs (MAX_PROCESSES = 64) |
