# Chapter 63: SMP and Multi-Core OS

> **"A single-core OS is already impressive. A multi-core OS is where the real challenges begin: two CPUs running simultaneously means races, deadlocks, cache coherency issues, and synchronization problems that simply don't exist when only one CPU is active. Understanding SMP is understanding why concurrency is hard — and how to do it right."**

---

## Table of Contents

1. [SMP — Symmetric Multi-Processing](#1-smp--symmetric-multi-processing)
2. [APIC — The Modern Interrupt Controller](#2-apic--the-modern-interrupt-controller)
3. [The APIC ID and AP Discovery](#3-the-apic-id-and-ap-discovery)
4. [Starting Application Processors (APs)](#4-starting-application-processors-aps)
5. [Per-CPU Data Structures](#5-per-cpu-data-structures)
6. [SMP-Safe Locking — Spinlocks](#6-smp-safe-locking--spinlocks)
7. [Ticket Locks and Fairness](#7-ticket-locks-and-fairness)
8. [SMP Scheduler Considerations](#8-smp-scheduler-considerations)
9. [IPI — Inter-Processor Interrupts](#9-ipi--inter-processor-interrupts)
10. [Summary](#summary)

---

## 1. SMP — Symmetric Multi-Processing

```
Single-core (what we've built so far):
  One CPU executes one instruction at a time
  "Parallelism" is simulated by fast context switching
  No true concurrency — only one piece of code runs at once
  No need for locking (interrupts off = no concurrent access)

Multi-core (SMP):
  Multiple CPUs (cores) share the same physical RAM
  Two CPUs can truly execute code simultaneously
  One CPU can be in the scheduler while another runs user code
  A global variable can be modified by CPU 0 and CPU 1 at the same time
  → Race conditions
  → Need explicit synchronization (locks, atomics)
  
Terminology:
  BSP  = Bootstrap Processor (the CPU that boots the system)
  AP   = Application Processor (additional CPUs, started by BSP)
  LAPIC = Local APIC (each CPU has its own)
  IOAPIC = I/O APIC (handles device interrupts, replaces 8259 PIC in SMP systems)
```

---

## 2. APIC — The Modern Interrupt Controller

```
In SMP systems, the 8259 PIC is replaced by the APIC (Advanced PIC):

LAPIC (Local APIC) — one per CPU core:
  Built into the CPU die since Pentium
  Handles:
    - Receiving IPIs (Inter-Processor Interrupts) from other CPUs
    - Local timer (each CPU can have its own timer)
    - Thermal/performance monitoring interrupts
    - Sending EOI (End of Interrupt) for local interrupts
  
  Memory-mapped at physical 0xFEE00000 (default)
  Each CPU's LAPIC has a unique ID (0, 1, 2, ...)
  
IOAPIC — usually one per system (or per PCI domain):
  Handles device IRQs (keyboard, disk, network, etc.)
  Replaces 8259 PIC master/slave chain
  Routes each device IRQ to a specific CPU's LAPIC
  
  Memory-mapped at physical 0xFEC00000 (default)
  Up to 24 redirection entries (IRQ routing table)
```

**LAPIC registers (memory-mapped at 0xFEE00000):**
```c
#define LAPIC_BASE      0xFEE00000

#define LAPIC_ID        0x0020   /* LAPIC ID register */
#define LAPIC_VERSION   0x0030   /* LAPIC version */
#define LAPIC_TPR       0x0080   /* Task Priority Register */
#define LAPIC_EOI       0x00B0   /* End of Interrupt (write 0 to acknowledge) */
#define LAPIC_SIVR      0x00F0   /* Spurious Interrupt Vector Register */
#define LAPIC_ICR_LO    0x0300   /* Interrupt Command Register (low 32 bits) */
#define LAPIC_ICR_HI    0x0310   /* Interrupt Command Register (high 32 bits) */
#define LAPIC_TIMER     0x0320   /* LVT Timer */
#define LAPIC_TIMER_INIT 0x0380  /* Timer Initial Count */
#define LAPIC_TIMER_CUR  0x0390  /* Timer Current Count */
#define LAPIC_TIMER_DIV  0x03E0  /* Timer Divide Config */

/* Read/write LAPIC MMIO: */
static inline uint32_t lapic_read(uint32_t reg) {
    return *((volatile uint32_t *)(LAPIC_BASE + reg));
}

static inline void lapic_write(uint32_t reg, uint32_t val) {
    *((volatile uint32_t *)(LAPIC_BASE + reg)) = val;
}

/* Send EOI to local LAPIC: */
static inline void lapic_eoi(void) {
    lapic_write(LAPIC_EOI, 0);
}

/* Get this CPU's LAPIC ID: */
uint8_t lapic_get_id(void) {
    return (uint8_t)(lapic_read(LAPIC_ID) >> 24);
}

/* Initialize the local APIC: */
void lapic_init(void) {
    /* Enable LAPIC (set bit 8 of SIVR) and set spurious interrupt vector: */
    lapic_write(LAPIC_SIVR, 0x1FF);  /* Vector 0xFF = spurious, bit 8 = enable */
    
    /* Set task priority to 0 (accept all interrupts): */
    lapic_write(LAPIC_TPR, 0);
}
```

---

## 3. The APIC ID and AP Discovery

APs are discovered via the ACPI MADT (Multiple APIC Description Table) or the MP tables:

```c
/* Simplified: we use CPUID to get the number of logical processors: */
static uint32_t cpu_count = 1;  /* At least BSP */
static uint8_t  cpu_lapic_ids[16];  /* Up to 16 CPUs */

void smp_detect_cpus(void) {
    /* Method 1: CPUID leaf 0xB (x2APIC topology) */
    /* Method 2: Parse ACPI MADT */
    /* Method 3: Parse Intel MP tables */
    
    /* Simplified: use CPUID leaf 4 for core count hint: */
    uint32_t eax, ebx, ecx, edx;
    __asm__ volatile(
        "cpuid"
        : "=a"(eax), "=b"(ebx), "=c"(ecx), "=d"(edx)
        : "a"(1)
    );
    
    /* EBX[23:16] = maximum number of addressable logical processor IDs: */
    uint8_t max_logical = (ebx >> 16) & 0xFF;
    if (max_logical == 0) max_logical = 1;
    
    kprintf("SMP: %u logical CPU(s) detected\n", max_logical);
    cpu_count = max_logical;
    
    /* BSP is CPU 0: */
    cpu_lapic_ids[0] = lapic_get_id();
}
```

---

## 4. Starting Application Processors (APs)

APs start in real mode and must be guided through the same mode switches as the BSP:

```c
/* AP Startup procedure (INIT-SIPI-SIPI):
   1. BSP sends INIT IPI to target AP → AP resets
   2. BSP sends SIPI (Startup IPI) → AP starts executing at vector×0x1000
   3. BSP sends second SIPI if AP didn't start
   4. AP runs 16-bit startup code, enables protected mode, calls ap_main() */

/* Trampoline code: 16-bit real mode code for APs to boot from.
   Must be placed at a physical address below 1MB that is page-aligned. */
#define AP_TRAMPOLINE_ADDR  0x8000   /* Page-aligned address below 1MB */

void smp_start_ap(uint8_t lapic_id) {
    /* Send INIT IPI: */
    lapic_write(LAPIC_ICR_HI, (uint32_t)lapic_id << 24);
    lapic_write(LAPIC_ICR_LO, 0x00004500);  /* INIT, Level Assert */
    
    /* Wait 10ms: */
    /* (use PIT delay or timer_get_ticks comparison) */
    uint32_t start = timer_get_ticks();
    while (timer_get_ticks() - start < 1) yield();
    
    /* Send first SIPI: vector = AP_TRAMPOLINE_ADDR / 0x1000 = 8 */
    lapic_write(LAPIC_ICR_HI, (uint32_t)lapic_id << 24);
    lapic_write(LAPIC_ICR_LO, 0x00004600 | (AP_TRAMPOLINE_ADDR >> 12));
    
    /* Wait briefly, then send second SIPI: */
    start = timer_get_ticks();
    while (timer_get_ticks() - start < 1) yield();
    lapic_write(LAPIC_ICR_HI, (uint32_t)lapic_id << 24);
    lapic_write(LAPIC_ICR_LO, 0x00004600 | (AP_TRAMPOLINE_ADDR >> 12));
}

/* Called by each AP after starting up: */
void ap_main(void) {
    /* Initialize this AP's LAPIC: */
    lapic_init();
    
    /* Set up GDT for this CPU (same GDT as BSP): */
    /* lgdt [gdt_descriptor] — already done in trampoline */
    
    kprintf("CPU %u is online!\n", lapic_get_id());
    
    /* Enable interrupts: */
    __asm__ volatile("sti");
    
    /* Enter scheduler: */
    for (;;) {
        schedule();
        __asm__ volatile("hlt");
    }
}
```

---

## 5. Per-CPU Data Structures

Each CPU needs its own:
- `current_process` pointer (which process is running on THIS CPU)
- Kernel stack (for nested interrupts)
- TSS entry (for Ring 3 → Ring 0 transitions)
- LAPIC EOI handling

```c
/* Per-CPU data: */
typedef struct {
    uint8_t      cpu_id;            /* LAPIC ID */
    uint8_t      online;            /* 1 after CPU is initialized */
    process_t   *current_process;  /* Currently running process on this CPU */
    uint32_t     kernel_stack_top; /* This CPU's kernel stack */
    uint32_t     ticks;            /* Per-CPU tick counter */
    tss_t        tss;              /* Task State Segment for this CPU */
} cpu_t;

static cpu_t cpus[16];

/* Get current CPU's data using LAPIC ID: */
cpu_t *get_current_cpu(void) {
    uint8_t id = lapic_get_id();
    return &cpus[id];
}
```

The `current_process` pointer must now be per-CPU (not a single global):
```c
#define current_process (get_current_cpu()->current_process)
```

---

## 6. SMP-Safe Locking — Spinlocks

On a single-core OS, `cli` (disable interrupts) is sufficient to prevent races. On SMP, disabling interrupts on CPU 0 doesn't stop CPU 1 from running the same code simultaneously. We need real locks:

```c
/* A spinlock: busy-wait until the lock is available.
   "Spin" = keep trying in a tight loop. */

typedef struct {
    volatile uint32_t locked;  /* 0 = free, 1 = taken */
} spinlock_t;

#define SPINLOCK_INIT { .locked = 0 }

/* Acquire the lock (spin until free, then take it): */
static inline void spinlock_acquire(spinlock_t *lock) {
    /* Disable interrupts first (prevents deadlock if interrupt handler tries same lock): */
    __asm__ volatile("cli");
    
    /* Atomic compare-and-swap loop:
       xchg atomically exchanges memory with register.
       If old value was 1 (already locked): spin.
       If old value was 0 (was free): we now hold the lock (set it to 1). */
    while (1) {
        uint32_t old = 1;
        __asm__ volatile(
            "xchg %0, %1"
            : "+r"(old), "+m"(lock->locked)
        );
        if (old == 0) break;  /* We got it! (old value was 0 = it was free) */
        /* Otherwise: old was 1 (locked), spin */
        __asm__ volatile("pause");  /* Hint to CPU: this is a spin loop (saves power) */
    }
}

/* Release the lock: */
static inline void spinlock_release(spinlock_t *lock) {
    lock->locked = 0;
    __asm__ volatile("sti");  /* Re-enable interrupts */
    __asm__ volatile("" ::: "memory");  /* Compiler memory barrier */
}

/* Usage example: */
static spinlock_t run_queue_lock = SPINLOCK_INIT;

void add_to_run_queue(process_t *proc) {
    spinlock_acquire(&run_queue_lock);
    proc->next = run_queue;
    if (run_queue) run_queue->prev = proc;
    run_queue = proc;
    spinlock_release(&run_queue_lock);
}
```

**Why `xchg`?**
```
`xchg reg, [mem]` is an ATOMIC operation:
  - Read the memory value
  - Write the register value to memory
  - These two steps happen without any other CPU being able to interfere
  - On x86, xchg automatically implies a LOCK prefix (bus lock)
  
Without atomicity:
  CPU0: read locked=0  (intends to set to 1)
  CPU1: read locked=0  (also intends to set to 1)
  CPU0: write 1        (thinks it has the lock)
  CPU1: write 1        (also thinks it has the lock)
  → Both CPUs think they hold the lock — RACE CONDITION!
```

---

## 7. Ticket Locks and Fairness

Simple spinlocks suffer from "thundering herd" and unfairness under contention. Ticket locks provide FIFO ordering:

```c
typedef struct {
    volatile uint32_t next_ticket;  /* Next ticket to give out */
    volatile uint32_t now_serving;  /* Current ticket being served */
} ticketlock_t;

void ticketlock_acquire(ticketlock_t *lock) {
    uint32_t my_ticket;
    
    /* Atomically get a ticket number: */
    __asm__ volatile(
        "lock xaddl %0, %1"      /* my_ticket = lock->next_ticket; lock->next_ticket++ */
        : "=r"(my_ticket)
        : "m"(lock->next_ticket), "0"(1)
    );
    
    /* Wait until it's our turn: */
    while (lock->now_serving != my_ticket) {
        __asm__ volatile("pause");
    }
}

void ticketlock_release(ticketlock_t *lock) {
    lock->now_serving++;  /* Serve the next customer */
}
```

---

## 8. SMP Scheduler Considerations

```
Single-core scheduler problems were:
  1. Which process runs next?
  2. When to preempt?
  
SMP scheduler adds:
  3. Which CPU runs it on?
  4. How to avoid all CPUs picking the same process?
  5. How to load balance across CPUs?
  
Simple SMP approach: global run queue with a spinlock.
  - Any CPU that needs work: lock the queue, pick a process, unlock.
  - Works, but the queue lock becomes a contention bottleneck with many CPUs.

Production approach (Linux CFS): per-CPU run queues.
  - Each CPU has its own queue.
  - Load balancer periodically migrates tasks from busy to idle CPUs.
  - Much less contention, but more complex.
  
For TinyOS: global run queue + spinlock is sufficient.
  scheduler picks next process (round-robin) under run_queue_lock.
  gdt_set_kernel_stack() must update the CURRENT CPU's TSS, not a global one.
```

---

## 9. IPI — Inter-Processor Interrupts

IPIs are used for:
- **TLB shootdown**: When CPU 0 unmaps a page, it must tell CPU 1 to flush its TLB for that page
- **Process wakeup**: Notify a sleeping CPU that there's work to do
- **System shutdown**: Tell all CPUs to stop

```c
#define IPI_TLB_SHOOTDOWN  0xF0
#define IPI_SCHEDULE       0xF1
#define IPI_HALT           0xFF

/* Send IPI to a specific CPU: */
void smp_send_ipi(uint8_t dest_lapic_id, uint8_t vector) {
    lapic_write(LAPIC_ICR_HI, (uint32_t)dest_lapic_id << 24);
    lapic_write(LAPIC_ICR_LO, vector | 0x00004000);  /* Fixed, edge-triggered */
}

/* Send TLB shootdown IPI to all other CPUs: */
void smp_tlb_shootdown(uint32_t vaddr) {
    /* Store the address to flush: */
    volatile uint32_t flush_addr = vaddr;
    (void)flush_addr;
    
    /* Broadcast to all except self: */
    lapic_write(LAPIC_ICR_HI, 0);
    lapic_write(LAPIC_ICR_LO, IPI_TLB_SHOOTDOWN | 0x000C4000); /* All-ex-self, Fixed */
}

/* IPI handler — called when we receive a TLB shootdown IPI: */
void ipi_tlb_handler(registers_t *regs) {
    (void)regs;
    /* Flush our TLB: */
    flush_tlb_all();
    lapic_eoi();
}
```

---

## Summary

| Concept | Description |
|---------|------------|
| SMP | Symmetric Multi-Processing: multiple CPUs sharing one RAM; each equally capable |
| BSP | Bootstrap Processor: the one CPU that boots the system |
| AP | Application Processor: additional CPUs started by BSP via INIT-SIPI-SIPI |
| LAPIC | Local APIC: per-CPU interrupt controller at 0xFEE00000 |
| IOAPIC | Handles device IRQs in SMP; replaces 8259 PIC; routes to chosen LAPIC |
| LAPIC ID | Each CPU's unique identifier; read from LAPIC ID register |
| SIPI | Startup IPI: tells an AP which 4KB-aligned address to start executing at |
| Per-CPU data | Each CPU needs its own current_process, TSS, kernel stack, tick counter |
| Spinlock | Busy-wait lock using atomic xchg; disables interrupts during critical section |
| xchg | Atomic exchange: read+write in one indivisible operation (no CPU can interrupt between) |
| pause | x86 hint: "I'm spinning" — reduces power, improves hyper-threading efficiency |
| Ticket lock | FIFO spinlock: each CPU gets a numbered ticket, serves in order |
| TLB shootdown | When one CPU changes page mappings, all CPUs must flush their TLBs |
| IPI | Inter-Processor Interrupt: one CPU signals another via LAPIC ICR register |
| Memory barrier | Prevents compiler/CPU from reordering memory accesses around critical sections |
