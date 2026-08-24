# Chapter 56: Context Switching

> **"Context switching is the illusion of parallelism on a single CPU. You save the current process's registers onto its kernel stack, load the next process's registers from its kernel stack, and suddenly the CPU is 'inside' a completely different program. Understand this mechanism and multitasking stops being magic."**

---

## Table of Contents

1. [What Is Context Switching?](#1-what-is-context-switching)
2. [What State Must Be Saved?](#2-what-state-must-be-saved)
3. [The Context Switch in Assembly](#3-the-context-switch-in-assembly)
4. [Switching Address Spaces](#4-switching-address-spaces)
5. [The TSS and Kernel Stack](#5-the-tss-and-kernel-stack)
6. [Putting It Together — switch_to()](#6-putting-it-together--switch_to)
7. [The Idle Process](#7-the-idle-process)
8. [Complete context.asm / switch.c](#8-complete-switchc)
9. [Testing Context Switching](#9-testing-context-switching)
10. [Summary](#summary)

---

## 1. What Is Context Switching?

```
Without context switching:
  One process runs until it voluntarily gives up the CPU (cooperative)
  or until it finishes
  Other processes must wait

With context switching (preemptive):
  Timer fires (IRQ0) → kernel saves current process state
  Kernel picks next process from run queue
  Kernel loads that process's state
  CPU resumes as if it was always in the new process
  
From the process's perspective: nothing happened.
Its registers, stack pointer, instruction pointer — all exactly where it left off.

From the kernel's perspective: just a stack swap.
```

---

## 2. What State Must Be Saved?

```
In x86 protected mode, during a kernel-kernel context switch (kernel threads):

MUST save:
  ESP — stack pointer (implicitly saved: we save the new ESP in pcb.kernel_esp)
  EIP — instruction pointer (saved as return address on the stack)
  
  Callee-saved registers (C calling convention):
    EBX, ESI, EDI, EBP
  
  These MUST be preserved by any C function — if we context switch
  in the middle of kernel code, these registers must be unchanged on return.

DO NOT need to explicitly save:
  EAX, ECX, EDX — caller-saved: C compiler already saves/restores these
  CS, DS, ES, FS, GS, SS — constant for kernel threads (all use kernel selectors)
  EFLAGS — we'll restore IF (interrupt flag) via sti/cli, not by saving

For user process → kernel (interrupt-driven):
  Everything is already saved on the kernel stack by the CPU + ISR stub:
  SS, ESP (old user), EFLAGS, CS, EIP — pushed by CPU
  EAX-EDI, DS-GS — pushed by our ISR stub
  This is the full saved state from which we can resume the user process.
```

---

## 3. The Context Switch in Assembly

```nasm
; boot/switch.asm — Context switch routine

[BITS 32]
[GLOBAL switch_context]

; void switch_context(uint32_t *old_esp, uint32_t new_esp, uint32_t new_cr3)
;   old_esp: address to store the current stack pointer
;   new_esp: the next process's kernel stack pointer
;   new_cr3: physical address of the next process's page directory (or 0 to keep current)

switch_context:
    ; Get arguments (cdecl: on stack after return address):
    ; [esp+4]  = old_esp (pointer to store current ESP)
    ; [esp+8]  = new_esp (load this as our new ESP)
    ; [esp+12] = new_cr3 (load this into CR3, or 0 to skip)
    
    ; --- SAVE current process state ---
    
    ; Push callee-saved registers onto the CURRENT stack:
    push ebp
    push ebx
    push esi
    push edi
    ; (EIP is already on the stack as the return address we'll use later)
    
    ; Save current ESP into old_esp:
    mov eax, [esp + 20]     ; old_esp argument (now at +20 because we pushed 4 regs)
    mov [eax], esp           ; *old_esp = current ESP
    
    ; --- SWITCH ADDRESS SPACE (if new_cr3 != 0) ---
    
    mov ecx, [esp + 28]     ; new_cr3 argument
    test ecx, ecx
    jz .skip_cr3
    
    mov cr3, ecx            ; Switch page directory → flushes non-global TLB
    
.skip_cr3:
    ; --- LOAD next process state ---
    
    ; Load new ESP:
    mov eax, [esp + 24]     ; new_esp argument
    mov esp, eax            ; NOW we're on the new process's stack
    
    ; Pop callee-saved registers from the NEW stack:
    pop edi
    pop esi
    pop ebx
    pop ebp
    
    ; Return — pops EIP from the new stack (= entry point or where new process was interrupted)
    ret
    ; From here, we're executing in the context of the new process!
```

---

## 4. Switching Address Spaces

```c
/* When switching between user processes, we must also change CR3: */

void switch_to(process_t *next) {
    if (!next) return;
    
    process_t *prev = current_process;
    current_process = next;
    next->state = PROC_RUNNING;
    
    /* Update TSS with the new process's kernel stack (for interrupt handling in Ring 3): */
    gdt_set_kernel_stack(next->kernel_stack + KERNEL_STACK_SIZE);
    
    if (prev) {
        prev->state = (prev->state == PROC_RUNNING) ? PROC_READY : prev->state;
        
        /* Context switch:
           - Save prev's ESP into prev->kernel_esp
           - Load next's ESP from next->kernel_esp
           - Switch page directory if different */
        uint32_t new_cr3 = (next->page_dir != prev->page_dir)
                           ? (uint32_t)next->page_dir : 0;
        
        switch_context(&prev->kernel_esp, next->kernel_esp, new_cr3);
        /* After switch_context returns, we're back in the PREV process
           (when it gets scheduled again in the future). */
    } else {
        /* First ever context switch — no "previous" process: */
        /* Just load the new state (can't save nothing): */
        uint32_t dummy;
        switch_context(&dummy, next->kernel_esp, (uint32_t)next->page_dir);
    }
}
```

---

## 5. The TSS and Kernel Stack

When an interrupt fires while executing user mode code (Ring 3), the CPU needs to find the kernel stack. It does this via the TSS:

```
Interrupt fires while in Ring 3:
  1. CPU reads current TSS (found via TR register)
  2. CPU loads SS0 and ESP0 from TSS → switches to kernel stack
  3. CPU pushes: old SS, old ESP, EFLAGS, old CS, old EIP onto kernel stack
  4. CPU jumps to IDT handler
  
This is why we call gdt_set_kernel_stack() before switching to a process:
  We update TSS.esp0 to point to the TOP of the new process's kernel stack.
  If an interrupt fires, the CPU will use THAT stack.
  
Without this: interrupt in Ring 3 would use the WRONG kernel stack
  → Corrupt the previous process's stack
  → Chaos

From gdt.c (already written in Ch 49):
  void gdt_set_kernel_stack(uint32_t esp0) {
      tss.esp0 = esp0;
  }
```

---

## 6. Putting It Together — switch_to()

The complete flow when the timer fires:

```
1. Timer IRQ fires (IRQ0)
2. CPU: push SS, ESP (if Ring 3), EFLAGS, CS, EIP
3. PIC: interrupt signal sent
4. Our IRQ stub: push dummy error, push vector, pusha, push segs
5. Call irq_handler() → calls timer_callback() → calls schedule()
6. schedule() calls switch_to(next_process)
7. switch_to() calls switch_context():
   - push EBP, EBX, ESI, EDI (save current process)
   - *old_esp = current ESP
   - mov esp, new_esp (switch stacks)
   - pop EDI, ESI, EBX, EBP (restore next process)
   - ret → jumps to next process's EIP
8. Now we're executing on the next process's stack
9. If this is first time running: EIP = entry_point
10. If resuming: unwind back through irq_handler → irq_common_stub → iret
11. iret restores: EIP, CS, EFLAGS, (old ESP, SS if Ring 3)
12. Process resumes exactly where it was interrupted
```

---

## 7. The Idle Process

When no process is runnable, the CPU must do something. The idle process:

```c
/* The idle process — runs when nothing else is ready: */
static void idle_process(void) {
    kprintf("[idle] Idle process started\n");
    while (1) {
        __asm__ volatile("sti; hlt");  /* enable interrupts, sleep until next IRQ */
    }
}

void process_init(void) {
    memset(process_table, 0, sizeof(process_table));
    
    /* Create the idle process (PID 0, lowest priority): */
    process_t *idle = pcb_alloc();
    idle->pid      = 0;
    idle->ppid     = 0;
    idle->state    = PROC_READY;
    idle->priority = 0;
    
    const char *idle_name = "idle";
    for (int i = 0; idle_name[i]; i++) idle->name[i] = idle_name[i];
    
    idle->kernel_stack = (uint32_t)kmalloc(KERNEL_STACK_SIZE);
    uint32_t *stack_top = (uint32_t *)(idle->kernel_stack + KERNEL_STACK_SIZE);
    *(--stack_top) = (uint32_t)idle_process;
    *(--stack_top) = 0; /* EBP */
    *(--stack_top) = 0; /* EBX */
    *(--stack_top) = 0; /* ESI */
    *(--stack_top) = 0; /* EDI */
    idle->kernel_esp = (uint32_t)stack_top;
    idle->page_dir   = kernel_page_dir;
    
    idle->next = run_queue;
    if (run_queue) run_queue->prev = idle;
    run_queue = idle;
    
    current_process = NULL;
}
```

---

## 8. Complete switch.c

```c
/* kernel/switch.c */

#include "process.h"
#include "gdt.h"
#include "vga.h"

/* Defined in boot/switch.asm: */
extern void switch_context(uint32_t *old_esp, uint32_t new_esp, uint32_t new_cr3);

void switch_to(process_t *next) {
    if (!next || next == current_process) return;
    
    process_t *prev = current_process;
    
    /* Update accounting: */
    next->state = PROC_RUNNING;
    
    /* Update TSS so Ring 3 → Ring 0 transition uses the correct kernel stack: */
    gdt_set_kernel_stack(next->kernel_stack + KERNEL_STACK_SIZE);
    
    /* Update current_process BEFORE the switch:
       Once we switch stacks, current_process must already be correct
       because any code running on the new stack will use it. */
    current_process = next;
    
    if (prev) {
        if (prev->state == PROC_RUNNING) prev->state = PROC_READY;
        
        uint32_t new_cr3 = (next->page_dir != prev->page_dir)
                           ? (uint32_t)next->page_dir : 0;
        
        switch_context(&prev->kernel_esp, next->kernel_esp, new_cr3);
        /* ← When this returns, we're back in PREV again (in the future) */
    } else {
        /* Bootstrap: no previous process */
        uint32_t dummy;
        uint32_t cr3 = (next->page_dir) ? (uint32_t)next->page_dir : 0;
        switch_context(&dummy, next->kernel_esp, cr3);
    }
}
```

---

## 9. Testing Context Switching

Context switching is tested together with the scheduler (Chapter 57). Here's what to expect:

```c
static int tick_a = 0, tick_b = 0;

static void proc_a(void) {
    while (1) {
        tick_a++;
        if (tick_a % 100 == 0) {
            kprintf("[A] tick=%d\n", tick_a);
        }
        /* Yield — voluntarily give up CPU: */
        schedule();
    }
}

static void proc_b(void) {
    while (1) {
        tick_b++;
        if (tick_b % 100 == 0) {
            kprintf("[B] tick=%d\n", tick_b);
        }
        schedule();
    }
}

/* If context switching works, output will interleave:
[A] tick=100
[B] tick=100
[A] tick=200
[B] tick=200
...
*/
```

---

## Summary

| Concept | Description |
|---------|------------|
| Context switch | Save current CPU state, load new process's CPU state |
| Callee-saved regs | EBX, ESI, EDI, EBP must be preserved across C function calls; we save/restore these |
| kernel_esp | Each PCB stores the process's kernel stack pointer (saved when not running) |
| switch_context | Assembly function: push regs onto old stack, save old ESP, load new ESP, pop regs, ret |
| CR3 switch | Load page directory address into CR3; must happen when switching user processes |
| TSS.esp0 | Kernel stack pointer in TSS; CPU reads this when interrupt fires in Ring 3 |
| gdt_set_kernel_stack | Update TSS.esp0 before switching to a process |
| Idle process | PID 0; runs when no other process is READY; hlt loop saves power |
| "Inside" the switch | After switch_context() returns, current CPU state matches the new process |
| First-time run | Entry point pre-loaded on stack as the "return address" of switch_context |
| Resume | On re-entry: switch_context() returns exactly where switch_to() was called before |
