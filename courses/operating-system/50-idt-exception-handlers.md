# Chapter 50: IDT and Exception Handlers

> **"Without an IDT, any CPU exception — divide by zero, null pointer access, stack overflow — causes a triple fault and immediately reboots the machine. With an IDT, you intercept every exception, print a useful error message, and decide what to do. The IDT is how your kernel listens to the CPU."**

---

## Table of Contents

1. [What Exceptions Are](#1-what-exceptions-are)
2. [The x86 Exception Table](#2-the-x86-exception-table)
3. [The IDT Structure](#3-the-idt-structure)
4. [Gate Descriptors](#4-gate-descriptors)
5. [Exception Handler Stubs in Assembly](#5-exception-handler-stubs-in-assembly)
6. [The Common C Handler](#6-the-common-c-handler)
7. [Installing the IDT](#7-installing-the-idt)
8. [Complete idt.c / idt.h / isr.c](#8-complete-idth--isrh)
9. [Testing — Triggering Exceptions](#9-testing--triggering-exceptions)
10. [Summary](#summary)

---

## 1. What Exceptions Are

When the CPU encounters an error condition, it generates an **exception** — a special interrupt that the OS must handle:

```
Three types of CPU faults:
  Fault:   CPU saves state BEFORE the faulting instruction.
           Handler can fix the problem and restart the instruction.
           Example: page fault — handler maps the page, restarts the access.
           
  Trap:    CPU saves state AFTER the faulting instruction.
           Handler runs, then execution continues at next instruction.
           Example: INT3 breakpoint — debugger reads state, then continues.
           
  Abort:   Unrecoverable error. CPU state may be corrupted.
           Example: double fault — handler should crash/reboot the system.

If no handler is installed (no valid IDT entry):
  → CPU generates a "double fault" (exception 8)
  → If that's also unhandled: "triple fault"
  → Triple fault: CPU hard resets (QEMU restarts)
```

---

## 2. The x86 Exception Table

The CPU has 32 predefined exceptions (vectors 0-31):

| Vector | Mnemonic | Name | Error code? | Type |
|--------|----------|------|-------------|------|
| 0 | #DE | Divide Error | No | Fault |
| 1 | #DB | Debug | No | Fault/Trap |
| 2 | — | NMI | No | Interrupt |
| 3 | #BP | Breakpoint | No | Trap |
| 4 | #OF | Overflow | No | Trap |
| 5 | #BR | BOUND Range Exceeded | No | Fault |
| 6 | #UD | Invalid Opcode | No | Fault |
| 7 | #NM | Device Not Available | No | Fault |
| 8 | #DF | Double Fault | Yes (0) | Abort |
| 9 | — | Coprocessor Seg Overrun | No | Fault |
| 10 | #TS | Invalid TSS | Yes | Fault |
| 11 | #NP | Segment Not Present | Yes | Fault |
| 12 | #SS | Stack-Segment Fault | Yes | Fault |
| 13 | #GP | General Protection Fault | Yes | Fault |
| 14 | #PF | Page Fault | Yes | Fault |
| 15 | — | Reserved | — | — |
| 16 | #MF | x87 FPU Error | No | Fault |
| 17 | #AC | Alignment Check | Yes (0) | Fault |
| 18 | #MC | Machine Check | No | Abort |
| 19 | #XF | SIMD FP Exception | No | Fault |
| 20-31 | — | Reserved | — | — |

Vectors 32-255 are for hardware interrupts (IRQs) and software interrupts (int $N).

**Error code:** Some exceptions push a 32-bit error code onto the stack before calling the handler. Our stubs must account for this.

---

## 3. The IDT Structure

The **Interrupt Descriptor Table** (IDT) is an array of up to 256 gate descriptors:

```
IDT layout:
  Address:    IDTR (stored in IDTR register, loaded via lidt)
  Entries:    256 × 8 bytes = 2048 bytes total
  Entry N:    Describes how to handle interrupt/exception N
  
IDTR format (48 bits):
  [47:16] = Base address of IDT
  [15:0]  = Limit (size - 1 = 2047 for 256 entries)

How the CPU uses the IDT when exception N occurs:
  1. CPU looks up IDT[N]
  2. Reads gate descriptor → finds handler address + segment selector
  3. Pushes SS, ESP (if privilege change), EFLAGS, CS, EIP
  4. Some exceptions push error code
  5. Loads CS from gate descriptor (must be kernel code selector)
  6. Jumps to handler address
  7. Handler executes, eventually calls iret to restore saved state
```

---

## 4. Gate Descriptors

Each 8-byte gate descriptor:

```
63      48 47 46-45 44 43-40 39     32 31          16 15           0
┌──────────┬──┬─────┬──┬─────┬─────────┬─────────────┬─────────────┐
│handler[31:16]│P│ DPL │0│type │ reserved│   selector  │handler[15:0]│
└──────────┴──┴─────┴──┴─────┴─────────┴─────────────┴─────────────┘

Fields:
  handler[15:0]:   Lower 16 bits of handler function address
  selector:        Code segment selector (we use GDT_KERNEL_CODE = 0x08)
  type:            Gate type:
                     0xE = 32-bit Interrupt Gate (clears IF — disables interrupts)
                     0xF = 32-bit Trap Gate (does NOT clear IF)
  DPL:             Descriptor Privilege Level (0 for exceptions, 3 for int 0x80)
  P:               Present (must be 1)
  handler[31:16]:  Upper 16 bits of handler function address

For exceptions (vectors 0-31): use Interrupt Gate, DPL=0
For system call (vector 0x80):  use Interrupt Gate or Trap Gate, DPL=3
```

---

## 5. Exception Handler Stubs in Assembly

We need assembly stubs for each exception. The challenge: some push error codes, others don't. We normalize the stack layout:

```nasm
; boot/isr_stubs.asm — Exception handler stubs

[BITS 32]
[EXTERN isr_handler]   ; The C handler function
[GLOBAL isr_stubs]

; Macro for exceptions WITHOUT error code:
%macro ISR_NO_ERR 1
[GLOBAL isr%1]
isr%1:
    push byte 0         ; Push dummy error code (to keep stack layout uniform)
    push byte %1        ; Push interrupt number
    jmp isr_common_stub
%endmacro

; Macro for exceptions WITH error code (CPU already pushed it):
%macro ISR_ERR 1
[GLOBAL isr%1]
isr%1:
    push byte %1        ; Push interrupt number (error code already on stack)
    jmp isr_common_stub
%endmacro

; Define all 32 exception handlers:
ISR_NO_ERR  0   ; #DE Divide Error
ISR_NO_ERR  1   ; #DB Debug
ISR_NO_ERR  2   ; NMI
ISR_NO_ERR  3   ; #BP Breakpoint
ISR_NO_ERR  4   ; #OF Overflow
ISR_NO_ERR  5   ; #BR BOUND Range
ISR_NO_ERR  6   ; #UD Invalid Opcode
ISR_NO_ERR  7   ; #NM Device Not Available
ISR_ERR     8   ; #DF Double Fault (error code = 0)
ISR_NO_ERR  9   ; Coprocessor Overrun
ISR_ERR    10   ; #TS Invalid TSS
ISR_ERR    11   ; #NP Segment Not Present
ISR_ERR    12   ; #SS Stack Fault
ISR_ERR    13   ; #GP General Protection Fault
ISR_ERR    14   ; #PF Page Fault
ISR_NO_ERR 15   ; Reserved
ISR_NO_ERR 16   ; #MF x87 FPU Error
ISR_ERR    17   ; #AC Alignment Check
ISR_NO_ERR 18   ; #MC Machine Check
ISR_NO_ERR 19   ; #XF SIMD FP
ISR_NO_ERR 20
ISR_NO_ERR 21
ISR_NO_ERR 22
ISR_NO_ERR 23
ISR_NO_ERR 24
ISR_NO_ERR 25
ISR_NO_ERR 26
ISR_NO_ERR 27
ISR_NO_ERR 28
ISR_NO_ERR 29
ISR_NO_ERR 30
ISR_NO_ERR 31

; Common stub — called by all handlers above:
; Stack at this point (top → bottom):
;   [ESP+0]:  interrupt number
;   [ESP+4]:  error code (real or dummy 0)
;   [ESP+8]:  EIP (saved by CPU)
;   [ESP+12]: CS  (saved by CPU)
;   [ESP+16]: EFLAGS (saved by CPU)
;   [ESP+20]: ESP_old (only if privilege change Ring 3 → Ring 0)
;   [ESP+24]: SS_old  (only if privilege change)

isr_common_stub:
    ; Save all general-purpose registers:
    pusha               ; pushes EAX, ECX, EDX, EBX, ESP, EBP, ESI, EDI
    
    ; Save data segments, load kernel data segment:
    push ds
    push es
    push fs
    push gs
    mov ax, 0x10        ; kernel data segment selector
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax
    
    ; Call the C handler with a pointer to the register frame:
    push esp            ; argument: pointer to registers_t struct on stack
    call isr_handler
    add esp, 4          ; clean up argument
    
    ; Restore segments:
    pop gs
    pop fs
    pop es
    pop ds
    
    ; Restore general-purpose registers:
    popa
    
    ; Clean up interrupt number and error code:
    add esp, 8
    
    ; Return from interrupt (restores EIP, CS, EFLAGS):
    iret
```

---

## 6. The Common C Handler

```c
/* kernel/isr.c — Interrupt/exception C handler */

#include "isr.h"
#include "vga.h"

/* Exception names for printing: */
static const char *exception_names[] = {
    "Divide Error",         /* 0  */
    "Debug Exception",      /* 1  */
    "NMI",                  /* 2  */
    "Breakpoint",           /* 3  */
    "Overflow",             /* 4  */
    "BOUND Range Exceeded", /* 5  */
    "Invalid Opcode",       /* 6  */
    "Device Not Available", /* 7  */
    "Double Fault",         /* 8  */
    "Coprocessor Overrun",  /* 9  */
    "Invalid TSS",          /* 10 */
    "Segment Not Present",  /* 11 */
    "Stack-Segment Fault",  /* 12 */
    "General Protection",   /* 13 */
    "Page Fault",           /* 14 */
    "Reserved",             /* 15 */
    "x87 FPU Error",        /* 16 */
    "Alignment Check",      /* 17 */
    "Machine Check",        /* 18 */
    "SIMD FP Exception",    /* 19 */
};

/* The register state saved by our stub: */
typedef struct {
    uint32_t gs, fs, es, ds;             /* segment registers */
    uint32_t edi, esi, ebp, esp_dummy,   /* pusha (esp_dummy = stack at pusha time) */
             ebx, edx, ecx, eax;
    uint32_t int_no;                     /* interrupt number (pushed by stub) */
    uint32_t err_code;                   /* error code (pushed by CPU or dummy 0) */
    uint32_t eip, cs, eflags;           /* saved by CPU before calling handler */
    uint32_t user_esp, user_ss;         /* only present if Ring 3 → Ring 0 transition */
} __attribute__((packed)) registers_t;

/* The C exception handler called from assembly: */
void isr_handler(registers_t *regs) {
    uint32_t num = regs->int_no;
    
    /* Print a "Blue Screen of Death" style error: */
    terminal_set_color(VGA_COLOR_WHITE, VGA_COLOR_BLUE);
    
    kprintf("\n\n  *** KERNEL EXCEPTION ***\n\n");
    
    const char *name = (num < 20) ? exception_names[num] : "Unknown";
    kprintf("  Exception %u: %s\n", num, name);
    kprintf("  Error Code: 0x%x\n\n", regs->err_code);
    
    kprintf("  EIP: 0x%x    CS: 0x%x    EFLAGS: 0x%x\n",
            regs->eip, regs->cs, regs->eflags);
    kprintf("  EAX: 0x%x    EBX: 0x%x\n", regs->eax, regs->ebx);
    kprintf("  ECX: 0x%x    EDX: 0x%x\n", regs->ecx, regs->edx);
    kprintf("  ESP: 0x%x    EBP: 0x%x\n", regs->esp_dummy, regs->ebp);
    kprintf("  ESI: 0x%x    EDI: 0x%x\n", regs->esi, regs->edi);
    kprintf("  DS: 0x%x     SS: 0x%x\n",  regs->ds, regs->cs);
    
    /* Page fault extra info: */
    if (num == 14) {
        uint32_t cr2;
        __asm__ volatile("mov %%cr2, %0" : "=r"(cr2));
        kprintf("\n  Page fault at address: 0x%x\n", cr2);
        kprintf("  Cause: %s %s on a %s page\n",
                (regs->err_code & 4) ? "user" : "kernel",
                (regs->err_code & 2) ? "write" : "read",
                (regs->err_code & 1) ? "present" : "not-present");
    }
    
    kprintf("\n  System halted.\n");
    
    /* Halt — don't return from an unhandled exception: */
    for (;;) {
        __asm__ volatile("cli; hlt");
    }
}
```

---

## 7. Installing the IDT

```c
/* kernel/idt.c — Install the IDT */

#include "idt.h"
#include "gdt.h"

/* The IDT table: */
static idt_entry_t idt[256];

/* IDT descriptor: */
static idt_descriptor_t idt_desc;

/* External symbols for ISR stubs (from isr_stubs.asm): */
extern void isr0(void);  extern void isr1(void);  extern void isr2(void);
extern void isr3(void);  extern void isr4(void);  extern void isr5(void);
extern void isr6(void);  extern void isr7(void);  extern void isr8(void);
extern void isr9(void);  extern void isr10(void); extern void isr11(void);
extern void isr12(void); extern void isr13(void); extern void isr14(void);
extern void isr15(void); extern void isr16(void); extern void isr17(void);
extern void isr18(void); extern void isr19(void); extern void isr20(void);
extern void isr21(void); extern void isr22(void); extern void isr23(void);
extern void isr24(void); extern void isr25(void); extern void isr26(void);
extern void isr27(void); extern void isr28(void); extern void isr29(void);
extern void isr30(void); extern void isr31(void);

/* Helper: set one IDT gate: */
static void idt_set_gate(uint8_t n, void (*handler)(void),
                          uint16_t selector, uint8_t flags) {
    uint32_t base = (uint32_t)handler;
    idt[n].base_low  = base & 0xFFFF;
    idt[n].base_high = (base >> 16) & 0xFFFF;
    idt[n].selector  = selector;
    idt[n].zero      = 0;
    idt[n].flags     = flags;
}

void idt_init(void) {
    idt_desc.limit = sizeof(idt) - 1;
    idt_desc.base  = (uint32_t)&idt;
    
    /* Install all 32 exception handlers:
       flags = 0x8E = 1000 1110 = P=1, DPL=0, type=0xE (32-bit interrupt gate) */
    idt_set_gate(0,  isr0,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(1,  isr1,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(2,  isr2,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(3,  isr3,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(4,  isr4,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(5,  isr5,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(6,  isr6,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(7,  isr7,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(8,  isr8,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(9,  isr9,  GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(10, isr10, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(11, isr11, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(12, isr12, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(13, isr13, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(14, isr14, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(15, isr15, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(16, isr16, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(17, isr17, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(18, isr18, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(19, isr19, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(20, isr20, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(21, isr21, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(22, isr22, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(23, isr23, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(24, isr24, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(25, isr25, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(26, isr26, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(27, isr27, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(28, isr28, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(29, isr29, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(30, isr30, GDT_KERNEL_CODE, 0x8E);
    idt_set_gate(31, isr31, GDT_KERNEL_CODE, 0x8E);
    
    /* Load the IDT: */
    __asm__ volatile("lidt %0" : : "m"(idt_desc));
    
    /* Enable interrupts: */
    __asm__ volatile("sti");
}
```

---

## 8. Complete idt.h / isr.h

```c
/* include/idt.h */
#pragma once
#include "stdint.h"

typedef struct {
    uint16_t base_low;
    uint16_t selector;
    uint8_t  zero;
    uint8_t  flags;
    uint16_t base_high;
} __attribute__((packed)) idt_entry_t;

typedef struct {
    uint16_t limit;
    uint32_t base;
} __attribute__((packed)) idt_descriptor_t;

void idt_init(void);
void idt_set_gate(uint8_t n, void (*handler)(void), uint16_t sel, uint8_t flags);
```

```c
/* include/isr.h */
#pragma once
#include "stdint.h"

typedef struct {
    uint32_t gs, fs, es, ds;
    uint32_t edi, esi, ebp, esp_dummy, ebx, edx, ecx, eax;
    uint32_t int_no, err_code;
    uint32_t eip, cs, eflags;
    uint32_t user_esp, user_ss;
} __attribute__((packed)) registers_t;

void isr_handler(registers_t *regs);
```

---

## 9. Testing — Triggering Exceptions

Add tests in `kernel_main`:

```c
#include "gdt.h"
#include "idt.h"

void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    terminal_init();
    
    kprintf("Installing GDT... ");
    gdt_init();
    kprintf("OK\n");
    
    kprintf("Installing IDT... ");
    idt_init();
    kprintf("OK\n");
    
    kprintf("Testing breakpoint (INT3)... ");
    __asm__ volatile("int $3");  /* Software breakpoint — vector 3 */
    /* Handler runs, then execution continues here (it's a trap, not a fault): */
    kprintf("returned from breakpoint!\n");
    
    /* Test divide by zero (comment this in to see the blue screen): */
    /* int x = 1 / 0; */
    
    kprintf("\nAll exception tests passed!\n");
    
    for (;;) {}
}
```

Wait — the INT3 handler above will show the blue screen and halt because our handler always halts. For INT3 to return, we'd need a special handler that returns for vector 3. For now, just observe that the blue screen appears (proving the IDT works) instead of an instant QEMU reboot (triple fault).

To test return-from-interrupt, modify the handler to only halt for fatal exceptions:

```c
void isr_handler(registers_t *regs) {
    if (regs->int_no == 3) {
        /* Breakpoint: just print and return */
        kprintf("[Breakpoint at EIP=0x%x]\n", regs->eip);
        return;  /* iret in stub will continue execution */
    }
    /* ... blue screen for everything else ... */
}
```

---

## Summary

| Concept | Description |
|---------|------------|
| Exception | CPU-generated event on error: fault (before), trap (after), abort (fatal) |
| IDT | Interrupt Descriptor Table: 256 gate descriptors at a known address |
| Gate descriptor | 8-byte entry: handler address, code selector, type (interrupt/trap), DPL |
| Interrupt gate (0x8E) | Clears IF (disables interrupts during handler); use for exceptions |
| Trap gate (0x8F) | Does NOT clear IF; used for breakpoints (INT3) |
| Error code | Some exceptions push extra info (PF: protection violation bits; GP: segment selector) |
| ISR stub | Assembly wrapper: push error code + int number, save all regs, call C handler |
| pusha | Save all 8 GPRs at once (EAX, ECX, EDX, EBX, ESP, EBP, ESI, EDI) |
| iret | Restore EIP, CS, EFLAGS (+ ESP, SS if Ring 3→0 transition happened) |
| lidt | Load Interrupt Descriptor Table Register (takes 48-bit limit+base) |
| sti | Set Interrupt Flag (enable interrupts) — call AFTER IDT is loaded |
| Page fault CR2 | CR2 register holds the faulting linear address when #PF occurs |
| Double fault | Exception during exception handling — usually means bad IDT/GDT/stack |
