# Chapter 29: Interrupts — How Hardware Talks to the CPU

> **"Without interrupts, the CPU would have to constantly ask every device 'do you have data for me?' — like a waiter running to every table every second asking if they want water. Interrupts flip this model: devices tap the CPU on the shoulder when they need attention. This simple mechanism is the foundation of all I/O in a modern computer."**

---

## Table of Contents

1. [The Problem Interrupts Solve](#1-the-problem-interrupts-solve)
2. [Types of Interrupts](#2-types-of-interrupts)
3. [The Interrupt Request (IRQ) System](#3-the-interrupt-request-irq-system)
4. [The Programmable Interrupt Controller (8259 PIC)](#4-the-programmable-interrupt-controller-8259-pic)
5. [The APIC — Modern Interrupt Controller](#5-the-apic--modern-interrupt-controller)
6. [The Interrupt Descriptor Table (IDT)](#6-the-interrupt-descriptor-table-idt)
7. [Interrupt Handling Flow](#7-interrupt-handling-flow)
8. [Interrupt Service Routines (ISRs)](#8-interrupt-service-routines-isrs)
9. [Top Halves and Bottom Halves](#9-top-halves-and-bottom-halves)
10. [Non-Maskable Interrupts and Exceptions](#10-non-maskable-interrupts-and-exceptions)
11. [Summary](#summary)

---

## 1. The Problem Interrupts Solve

**Polling (bad approach):**
```c
// CPU constantly asking devices:
while (1) {
    if (keyboard_has_key()) process_key();
    if (disk_read_complete()) process_data();
    if (network_packet_arrived()) process_packet();
    if (timer_fired()) update_scheduler();
    // ... repeat 1 million times per second
}
// CPU wastes ALL its time asking devices that are usually idle
```

**Interrupts (good approach):**
```
CPU runs normal code
...
Keyboard sends electrical signal on IRQ1 line
CPU: "Interrupt! Stop what I'm doing, handle keyboard"
Execute keyboard handler
Return to previous code exactly where I stopped
...
Continues normal execution
```

The CPU does useful work until a device needs attention. This is called **interrupt-driven I/O**.

---

## 2. Types of Interrupts

Linux (and most OSes) distinguish three types of "interrupt-like" events:

**Hardware interrupts (IRQs):**
Asynchronous signals from external devices:
- Timer chip (every ~1ms on x86) → drives the scheduler
- Keyboard → key press/release
- Mouse → movement/click
- Network card → packet arrived
- Disk → read/write complete
- USB → device event

**Software interrupts (syscalls):**
User programs invoke the kernel via:
- `int 0x80` (legacy 32-bit Linux syscall mechanism)
- `syscall` instruction (x86-64 fast path)
- Triggered intentionally by software, not hardware

**Exceptions (CPU exceptions/traps/faults):**
CPU detects an error condition:
- Division by zero (#DE, vector 0)
- Invalid opcode (#UD, vector 6)
- General Protection Fault (#GP, vector 13) — segmentation violation
- Page fault (#PF, vector 14) — memory not mapped
- Double fault (#DF, vector 8) — fault while handling another fault
- Machine check (#MC, vector 18) — hardware error (ECC, etc.)

---

## 3. The Interrupt Request (IRQ) System

On x86, hardware interrupts come through **IRQ lines** (Interrupt Request lines):

```
Classic IRQ assignments (legacy PC):
IRQ 0:  System timer (PIT — Programmable Interval Timer, ~1000 Hz)
IRQ 1:  PS/2 Keyboard
IRQ 2:  Cascade (used by 8259 to chain second PIC)
IRQ 3:  COM2 serial port
IRQ 4:  COM1 serial port
IRQ 5:  LPT2 (parallel port) or sound card
IRQ 6:  Floppy controller
IRQ 7:  LPT1 (parallel port / printer)
IRQ 8:  Real-time clock (RTC)
IRQ 9:  Available / ACPI
IRQ 10: Available / NIC
IRQ 11: Available / USB / NIC
IRQ 12: PS/2 Mouse
IRQ 13: Floating point unit
IRQ 14: Primary IDE controller
IRQ 15: Secondary IDE controller
```

Modern systems use **MSI (Message Signaled Interrupts)** — devices write to a special memory address to generate interrupts instead of using dedicated wires. Allows hundreds of IRQs for a single PCIe device.

---

## 4. The Programmable Interrupt Controller (8259 PIC)

The original IBM PC used the **Intel 8259 PIC** to manage IRQs. Two 8259 chips were chained together (cascade), giving 15 usable IRQ lines (16 minus cascade line):

```
        CPU
         │
     ┌───▼────┐
     │ Master  │  ← IRQ 0-7 (timer, keyboard, cascade, serial, floppy, LPT, RTC...)
     │  8259   │
     │ INT→CPU │
     └────┬────┘
          │ IRQ2 (cascade)
     ┌────▼────┐
     │  Slave  │  ← IRQ 8-15 (RTC, available, mouse, FPU, IDE...)
     │  8259   │
     └─────────┘
```

**Interrupt Mask Register (IMR):**
Each PIC has an 8-bit register where you can **mask** (disable) individual IRQ lines:
```c
// Mask IRQ1 (keyboard) in the master PIC:
outb(0x21, inb(0x21) | 0x02);  // set bit 1

// Unmask IRQ1:
outb(0x21, inb(0x21) & ~0x02); // clear bit 1
```

**End Of Interrupt (EOI):**
After handling an IRQ, the OS must tell the PIC "I'm done":
```c
// Signal EOI to master PIC:
outb(0x20, 0x20);

// For IRQ 8-15 (slave PIC), send EOI to BOTH:
outb(0xA0, 0x20);  // slave EOI
outb(0x20, 0x20);  // master EOI
```
Without EOI, the PIC won't deliver any more interrupts of that priority!

**PIC interrupt vectors:**
The 8259 is programmed to map IRQ0-7 to interrupt vectors 0x20-0x27, and IRQ8-15 to 0x28-0x2F (to avoid conflicts with CPU exception vectors 0-19).

---

## 5. The APIC — Modern Interrupt Controller

Modern x86 systems use the **APIC (Advanced Programmable Interrupt Controller)**:

**Local APIC (LAPIC):** One per CPU core. Receives and prioritizes interrupts for that core.

**I/O APIC:** One or more per system. Routes hardware IRQs to one or more LAPICs.

```
Hardware IRQ 0 (timer) → I/O APIC → routes to → CPU 0 LAPIC → CPU 0 handles it
Hardware IRQ 1 (keyboard) → I/O APIC → routes to → CPU 0 LAPIC → CPU 0 handles it
Network card MSI → I/O APIC → can be routed to any CPU for load balancing
```

**Advantages over 8259:**
- SMP support: route interrupts to specific CPUs
- Priority: higher-priority interrupts preempt lower ones
- MSI support: 200+ interrupt vectors per PCIe device
- Inter-Processor Interrupts (IPI): one CPU can interrupt another
  - Used for TLB shootdowns: "Hey CPU3, I changed a page table entry, invalidate your TLB"
  - Used for scheduler IPI: "Wake up and reschedule"

**LAPIC timer:**
The LAPIC has a built-in timer (used instead of the old PIT for per-CPU scheduling):
```c
// Set LAPIC timer to fire every 1ms (one-shot mode):
lapic_write(LAPIC_TIMER_DIVIDE, 0x3);   // divide by 16
lapic_write(LAPIC_TIMER_INITIAL, 1000); // initial count
lapic_write(LAPIC_LVT_TIMER, 0x20000 | TIMER_VECTOR); // periodic mode + vector
```

---

## 6. The Interrupt Descriptor Table (IDT)

The CPU needs to know: **when interrupt N fires, which code should I run?**

The answer is the **IDT (Interrupt Descriptor Table)** — an array of 256 gate descriptors:

```
IDT[0]  = descriptor for exception 0 (divide by zero) → #DE handler
IDT[1]  = descriptor for debug exception → #DB handler
IDT[2]  = descriptor for NMI → NMI handler
IDT[3]  = descriptor for breakpoint (int 3) → #BP handler
...
IDT[13] = descriptor for GPF → #GP handler
IDT[14] = descriptor for page fault → #PF handler
...
IDT[32] = descriptor for IRQ0 (timer) → timer_interrupt handler
IDT[33] = descriptor for IRQ1 (keyboard) → keyboard_interrupt handler
...
IDT[128] = descriptor for int 0x80 (syscall, 32-bit legacy) → syscall handler
```

**Gate descriptor format (64-bit, 16 bytes):**
```
Gate descriptor:
  bits 15:0    Handler offset [15:0]
  bits 31:16   Code segment selector (CS, usually kernel CS)
  bits 34:32   IST (Interrupt Stack Table — for separate stack on NMI/DF)
  bits 39:35   Reserved = 0
  bits 43:40   Type: 0xE = interrupt gate, 0xF = trap gate
  bit  44      S = 0 (system descriptor)
  bits 46:45   DPL (privilege level that can call this via int N)
  bit  47      P = 1 (present)
  bits 63:48   Handler offset [31:16]
  bits 95:64   Handler offset [63:32]
  bits 127:96  Reserved = 0
```

**Interrupt gate vs Trap gate:**
- **Interrupt gate**: CPU clears the IF flag (disables further interrupts) when handler runs
- **Trap gate**: CPU does NOT clear IF (interrupts remain enabled)
- Exception handlers use trap gates; hardware IRQ handlers use interrupt gates

**Loading the IDT:**
```c
struct idtr {
    uint16_t limit;  // IDT size - 1 = 256*16 - 1 = 4095
    uint64_t base;   // address of IDT array
} __attribute__((packed));

struct idtr idt_register = { sizeof(idt) - 1, (uint64_t)idt };
asm volatile("lidt (%0)" : : "r"(&idt_register));
```

---

## 7. Interrupt Handling Flow

**What the CPU does automatically on interrupt:**

```
1. CPU detects interrupt signal on its INTR pin (or internal exception)
2. CPU checks IF flag — if 0 (disabled), defer until reenabled (unless NMI)
3. CPU gets interrupt vector number N (from PIC/APIC/exception logic)
4. CPU reads IDT[N] → get handler address and segment selector
5. CPU saves current state to stack:
   a. If privilege change (Ring 3→0): switch to kernel stack (RSP from TSS.RSP0)
   b. Push SS, RSP (user stack)
   c. Push RFLAGS
   d. Push CS
   e. Push RIP (return address)
   f. For some exceptions: push error code
6. CPU loads new CS from IDT entry (sets CPL to Ring 0)
7. CPU jumps to handler address
8. Handler runs (the OS interrupt service routine)
9. Handler executes iret instruction:
   a. Pops RIP, CS, RFLAGS from stack
   b. If returning to Ring 3: pops RSP, SS (switches back to user stack)
   c. Continues execution at the interrupted instruction
```

**Stack frame during interrupt handler:**
```
(Stack grows down)
┌─────────────┐  ← SS:RSP before interrupt
│   SS        │  (saved if Ring 3 → Ring 0 change)
│   RSP       │
│   RFLAGS    │
│   CS        │
│   RIP       │  ← return address
│  [Error Code│  (only for some exceptions)
└─────────────┘  ← RSP during interrupt handler
```

---

## 8. Interrupt Service Routines (ISRs)

An **ISR (Interrupt Service Routine)** is the handler function that runs when an interrupt fires.

**Constraints on ISR code:**
- Must be FAST — while handling one interrupt, others may be deferred
- Cannot block (sleep) — must return quickly
- Cannot call most kernel functions that might sleep
- Must save and restore any registers it modifies (beyond what CPU auto-saves)

**Generic ISR stub (x86-64 NASM):**
```nasm
; All ISR stubs push the vector number and jump to common handler
isr32:              ; IRQ0 (timer), mapped to vector 32
    push 0          ; no error code for this interrupt → push dummy
    push 32         ; interrupt vector number
    jmp isr_common

isr33:              ; IRQ1 (keyboard), vector 33
    push 0
    push 33
    jmp isr_common

isr_common:
    ; Save all general-purpose registers
    push rax
    push rcx
    push rdx
    push rbx
    push rbp
    push rsi
    push rdi
    push r8
    push r9
    push r10
    push r11
    push r12
    push r13
    push r14
    push r15

    ; Call C interrupt handler with pointer to stack frame
    mov rdi, rsp    ; first argument = pointer to saved registers
    call interrupt_handler

    ; Restore registers
    pop r15
    pop r14
    ... (reverse order)
    pop rax

    add rsp, 16     ; discard error code + vector number
    iretq           ; return from interrupt (pops RIP, CS, RFLAGS, RSP, SS)
```

**C handler:**
```c
typedef struct {
    uint64_t r15, r14, r13, r12, r11, r10, r9, r8;
    uint64_t rdi, rsi, rbp, rbx, rdx, rcx, rax;
    uint64_t vector, error_code;
    uint64_t rip, cs, rflags, rsp, ss;  // pushed by CPU
} interrupt_frame_t;

void interrupt_handler(interrupt_frame_t *frame) {
    switch (frame->vector) {
        case 32: timer_interrupt(frame); break;
        case 33: keyboard_interrupt(frame); break;
        // ...
    }
    // Send EOI to PIC/APIC
    if (frame->vector >= 32 && frame->vector < 48) {
        if (frame->vector >= 40) outb(0xA0, 0x20);  // slave EOI
        outb(0x20, 0x20);                            // master EOI
    }
}
```

---

## 9. Top Halves and Bottom Halves

**The problem:**
Some interrupt handlers need to do a lot of work (e.g., network driver receives a packet and must parse TCP headers, do checksums, wake blocked processes). But ISRs must be fast!

**Solution: split interrupt handling into two parts:**

**Top half (hard IRQ handler):**
- Runs with interrupts disabled (or at high priority)
- Does the minimum: acknowledge the interrupt, save data from hardware
- Must complete in microseconds
- Schedules "bottom half" work

**Bottom half (deferred work):**
- Runs later, with interrupts enabled
- Does the heavy lifting: protocol processing, allocating memory, waking processes
- Lower priority than hardware interrupts

**Linux deferred work mechanisms:**

**Softirqs:**
```
Fixed set of high-priority deferred work types:
HI_SOFTIRQ:       High-priority tasklets
TIMER_SOFTIRQ:    Software timers
NET_TX_SOFTIRQ:   Network transmission
NET_RX_SOFTIRQ:   Network reception ← most common
BLOCK_SOFTIRQ:    Block device I/O completion
TASKLET_SOFTIRQ:  Tasklets
SCHED_SOFTIRQ:    Scheduler
```

**Tasklets (built on softirqs):**
Dynamically allocated, runs once, on the CPU that scheduled it, can be re-enabled:
```c
// Declare a tasklet:
void my_tasklet_handler(unsigned long data);
DECLARE_TASKLET(my_tasklet, my_tasklet_handler, 0);

// In top half (IRQ handler):
tasklet_schedule(&my_tasklet);  // queue it for later
// ... return quickly

// Later, kernel runs my_tasklet_handler() in softirq context
```

**Work queues:**
Runs in kernel thread context — can sleep, can do disk I/O, can call any kernel function:
```c
INIT_WORK(&my_work, my_work_handler);
schedule_work(&my_work);
```

**Example: Network receive path:**
```
1. Network card gets a packet → raises IRQ
2. Top half (hard IRQ):
   - Acknowledge interrupt to NIC
   - Copy packet from NIC DMA buffer to kernel skbuff
   - Schedule NET_RX_SOFTIRQ
   - Return (done in ~5 microseconds)
3. Bottom half (NET_RX_SOFTIRQ):
   - Parse Ethernet header
   - Parse IP header, verify checksum
   - Parse TCP header, find matching socket
   - Copy data to socket receive buffer
   - Wake up process waiting on this socket
   (Takes ~50-100 microseconds — OK to be slow)
```

---

## 10. Non-Maskable Interrupts and Exceptions

**NMI (Non-Maskable Interrupt):**
A special interrupt that CANNOT be disabled (even with CLI instruction):
- Vector 2
- Used for hardware failures: ECC memory error, watchdog timer, bus error
- Also used by debuggers (NMI watchdog detects hung CPUs)
- Handler must be extremely careful (can run at any time, even during another interrupt)

```bash
# Linux NMI watchdog — detects hung CPUs:
cat /proc/sys/kernel/nmi_watchdog   # 1 = enabled

# If a CPU is stuck (holding a spinlock for too long):
# NMI fires → NMI handler → detects the hang → prints stack trace → reboots or panics
```

**CPU Exceptions:**
```
#DE  (0):  Divide Error — division by zero, division overflow
#DB  (1):  Debug — step trap, breakpoint (also used by gdb)
#BP  (3):  Breakpoint — int 3 instruction (0xCC byte — gdb's breakpoint)
#OF  (4):  Overflow — INTO instruction with OF set
#BR  (5):  Bound Range Exceeded — BOUND instruction
#UD  (6):  Undefined Opcode — invalid instruction
#NM  (7):  Device Not Available — FPU not present or disabled
#DF  (8):  Double Fault — fault during fault handler; pushes error code 0
#TS  (10): Invalid TSS — bad TSS descriptor
#NP  (11): Segment Not Present — segment descriptor P=0
#SS  (12): Stack Fault — stack segment violation
#GP  (13): General Protection Fault — most common: null deref, ring violation
#PF  (14): Page Fault — access to unmapped or protected memory
#MF  (16): x87 FPU Error — floating point exception
#AC  (17): Alignment Check — misaligned memory access with AC flag
#MC  (18): Machine Check — hardware error (ECC, CPU internal error)
#XF  (19): SIMD FP Exception — SSE floating point exception
```

**#GP (General Protection Fault)** is the "catch-all" for privilege violations:
- Accessing a kernel address from user mode
- Executing a privileged instruction from Ring 3
- Loading an invalid segment register
- Using a misaligned address for aligned-only instruction

---

## Summary

| Concept | Description |
|---------|------------|
| Hardware interrupt | Asynchronous signal from external device requesting CPU attention |
| Software interrupt | int N instruction; user space → kernel transition (syscalls) |
| Exception | CPU-detected error: division by zero, page fault, GPF, etc. |
| IRQ | Interrupt Request line; hardware device signals via numbered line |
| PIC (8259) | Legacy interrupt controller; chains master+slave for 15 IRQ lines |
| APIC | Modern interrupt controller; supports SMP, MSI, 256 vectors per CPU |
| IDT | Interrupt Descriptor Table: maps vector number → handler address |
| Gate descriptor | 16-byte IDT entry: handler address, segment, DPL, type |
| ISR | Interrupt Service Routine: handler called when interrupt fires |
| EOI | End of Interrupt signal: tells PIC/APIC to allow more interrupts |
| Top half | Fast ISR: acks hardware, saves data, schedules bottom half |
| Bottom half | Deferred work: softirq, tasklet, work queue — can be slower |
| NMI | Non-Maskable Interrupt: cannot be disabled; used for hardware failures |
| #PF | Page Fault exception; vector 14; OS handles demand paging + protection |
| #GP | General Protection Fault; vector 13; privilege violation |
| iret/iretq | Return from interrupt; restores CPU state saved by interrupt entry |
