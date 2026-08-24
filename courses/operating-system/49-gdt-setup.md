# Chapter 49: GDT Setup

> **"The Global Descriptor Table is the foundation of memory protection in x86 protected mode. Without it, your CPU can't tell code from data, kernel from user, valid memory from garbage. Setting up the GDT correctly is like pouring the concrete before raising the walls — everything else depends on it."**

---

## Table of Contents

1. [Why We Need to Set Up Our Own GDT](#1-why-we-need-to-set-up-our-own-gdt)
2. [Descriptor Format — The 8-Byte Entry](#2-descriptor-format--the-8-byte-entry)
3. [The Selectors We Need](#3-the-selectors-we-need)
4. [GDT in C — struct + lgdt](#4-gdt-in-c--struct--lgdt)
5. [Task State Segment (TSS) — Introduced Here](#5-task-state-segment-tss--introduced-here)
6. [Reloading Segment Registers](#6-reloading-segment-registers)
7. [Complete gdt.c / gdt.h](#7-complete-gdtc--gdth-summary)
8. [Testing the GDT](#8-testing-the-gdt)
9. [Summary](#summary)

---

## 1. Why We Need to Set Up Our Own GDT

GRUB sets up a minimal GDT to put the CPU in protected mode, but it's not guaranteed to be what we want. We need to set up our own GDT because:

```
Problems with GRUB's GDT:
  ✗ Location is unknown (GRUB may have it anywhere in memory)
  ✗ Does not include a TSS (Task State Segment) — needed for system calls and user mode
  ✗ No user-mode segments (DPL=3) — needed for Ring 3
  ✗ May conflict with memory we want to use for our kernel

Our GDT will have:
  ✓ Known, fixed location (in our .data section)
  ✓ 5 entries + 1 TSS entry
  ✓ Kernel code + data (Ring 0)
  ✓ User code + data (Ring 3) — needed for Ch 59
  ✓ TSS descriptor — needed for Ch 56 (context switching)
```

---

## 2. Descriptor Format — The 8-Byte Entry

Each GDT entry is exactly 8 bytes with a somewhat confusing layout (designed in the 80286 era and extended in the 80386):

```
63      56 55  52 51   48 47    40 39        16 15          0
┌─────────┬──────┬───────┬────────┬────────────┬────────────┐
│base[31:24]│flags│lim[19:16]│access │ base[23:0] │  limit[15:0]│
└─────────┴──────┴───────┴────────┴────────────┴────────────┘

Byte positions:
  Bytes 0-1: Limit[15:0]    — lower 16 bits of segment limit
  Bytes 2-4: Base[23:0]     — lower 24 bits of base address
  Byte 5:    Access byte    — type, privilege, present bits
  Byte 6:    Flags + Limit[19:16]
  Byte 7:    Base[31:24]    — upper 8 bits of base address

Access byte (byte 5):
  Bit 7:  P   = Present (must be 1 for valid descriptor)
  Bit 6-5: DPL = Descriptor Privilege Level (0=kernel, 3=user)
  Bit 4:  S   = Descriptor type (1=code/data, 0=system)
  Bit 3:  E   = Executable (1=code segment, 0=data segment)
  Bit 2:  DC  = Direction/Conforming
              Data: 0=grows up, 1=grows down
              Code: 0=only exec at DPL, 1=conforming (can call from lower privilege)
  Bit 1:  RW  = Readable/Writable (code: readable; data: writable)
  Bit 0:  A   = Accessed (CPU sets this when segment is accessed; we set to 0)

Flags byte (upper nibble of byte 6):
  Bit 7:  G  = Granularity (0=limit in bytes, 1=limit in 4KB pages)
  Bit 6: DB  = Default operation size (0=16-bit, 1=32-bit segment)
  Bit 5:  L  = Long mode (1=64-bit code segment; DB must be 0)
  Bit 4:  Reserved (set to 0)
```

**Kernel code segment (selector 0x08):**
```
P=1, DPL=0, S=1, E=1, DC=0, RW=1, A=0 → access = 1_00_1_1010 = 0x9A
G=1, DB=1, L=0 → flags = 1100 = 0xC
Base = 0x00000000, Limit = 0xFFFFF (with G=1: covers 4GB)
```

**Kernel data segment (selector 0x10):**
```
P=1, DPL=0, S=1, E=0, DC=0, RW=1, A=0 → access = 1_00_1_0010 = 0x92
G=1, DB=1 → flags = 0xC
Base = 0, Limit = 0xFFFFF
```

**User code segment (selector 0x18 | 3 = 0x1B):**
```
P=1, DPL=3, S=1, E=1, DC=0, RW=1, A=0 → access = 1_11_1_1010 = 0xFA
G=1, DB=1 → flags = 0xC
```

**User data segment (selector 0x20 | 3 = 0x23):**
```
P=1, DPL=3, S=1, E=0, DC=0, RW=1, A=0 → access = 1_11_1_0010 = 0xF2
```

---

## 3. The Selectors We Need

A **segment selector** is a 16-bit value loaded into CS/DS/SS/ES/FS/GS:

```
Selector format:
  Bits 15-3:  Index into GDT (or LDT) — 0-based
  Bit 2:      TI = Table Indicator (0=GDT, 1=LDT)
  Bits 1-0:   RPL = Requested Privilege Level (0-3)

Our GDT entries and their selectors:
  Index 0: Null descriptor   → selector 0x00 (never load into segment register)
  Index 1: Kernel code       → selector 0x08 (index=1, TI=0, RPL=0)
  Index 2: Kernel data       → selector 0x10 (index=2, TI=0, RPL=0)
  Index 3: User code         → selector 0x18 (index=3, TI=0, RPL=0)
                                              load with RPL=3: 0x18|3 = 0x1B
  Index 4: User data         → selector 0x20 (index=4, TI=0, RPL=0)
                                              load with RPL=3: 0x20|3 = 0x23
  Index 5: TSS descriptor    → selector 0x28 (index=5, TI=0, RPL=0)

Convention: kernel uses 0x08/0x10; user mode uses 0x1B/0x23
```

---

## 4. GDT in C — struct + lgdt

```c
/* include/gdt.h */
#pragma once
#include "stdint.h"

/* Segment selectors: */
#define GDT_KERNEL_CODE  0x08
#define GDT_KERNEL_DATA  0x10
#define GDT_USER_CODE    0x18   /* | 3 = 0x1B when loaded with RPL=3 */
#define GDT_USER_DATA    0x20   /* | 3 = 0x23 when loaded with RPL=3 */
#define GDT_TSS          0x28

/* Number of GDT entries (null + kernel code/data + user code/data + TSS): */
#define GDT_ENTRIES      6

/* A GDT entry (8 bytes): */
typedef struct {
    uint16_t limit_low;     /* Limit[15:0] */
    uint16_t base_low;      /* Base[15:0] */
    uint8_t  base_mid;      /* Base[23:16] */
    uint8_t  access;        /* Access byte */
    uint8_t  limit_high_flags; /* Limit[19:16] in low nibble, Flags in high nibble */
    uint8_t  base_high;     /* Base[31:24] */
} __attribute__((packed)) gdt_entry_t;

/* The GDT descriptor (48 bits) loaded by lgdt: */
typedef struct {
    uint16_t limit;  /* GDT size in bytes - 1 */
    uint32_t base;   /* Linear address of GDT */
} __attribute__((packed)) gdt_descriptor_t;

/* Task State Segment (basic layout, more in Ch 56): */
typedef struct {
    uint32_t prev_tss;
    uint32_t esp0;    /* Ring 0 stack pointer (for privilege level change) */
    uint32_t ss0;     /* Ring 0 stack segment */
    /* ... rest unused for now ... */
    uint8_t  unused[92];
    uint16_t io_map;  /* I/O permission bitmap offset */
} __attribute__((packed)) tss_t;

void gdt_init(void);
void gdt_set_kernel_stack(uint32_t esp0);
```

```c
/* kernel/gdt.c */
#include "gdt.h"
#include "string.h"  /* memset */

/* The actual GDT table: */
static gdt_entry_t gdt[GDT_ENTRIES];

/* GDT descriptor (loaded into GDTR register): */
static gdt_descriptor_t gdt_desc;

/* TSS: */
static tss_t tss;

/* Helper: encode one GDT entry: */
static void gdt_set_entry(int index, uint32_t base, uint32_t limit,
                           uint8_t access, uint8_t flags) {
    gdt[index].base_low    = (base  & 0x0000FFFF);
    gdt[index].base_mid    = (base  & 0x00FF0000) >> 16;
    gdt[index].base_high   = (base  & 0xFF000000) >> 24;
    gdt[index].limit_low   = (limit & 0x0000FFFF);
    /* Limit[19:16] in lower nibble + flags in upper nibble: */
    gdt[index].limit_high_flags = ((limit >> 16) & 0x0F) | (flags & 0xF0);
    gdt[index].access      = access;
}

/* Set up the TSS descriptor in the GDT: */
static void gdt_set_tss(int index, uint32_t base, uint32_t limit) {
    /* TSS descriptor has the same layout but different access byte: */
    /* Access = 0x89 = P=1, DPL=0, S=0 (system), type=0x9 (32-bit TSS available) */
    gdt_set_entry(index, base, limit, 0x89, 0x00);
}

/* Reload all segment registers (in assembly): */
extern void gdt_flush(uint32_t gdt_desc_addr);

/* Called from kernel_main to set up the GDT: */
void gdt_init(void) {
    /* Entry 0: Null descriptor (required — loading into segment reg causes GP fault): */
    gdt_set_entry(0, 0, 0, 0, 0);
    
    /* Entry 1: Kernel code segment (Ring 0, code, readable, 4GB flat): */
    /* access = 0x9A = 1001 1010:
       P=1, DPL=00, S=1, E=1, DC=0, RW=1, A=0 */
    /* flags  = 0xCF upper nibble = 0xC:
       G=1 (4KB granularity), DB=1 (32-bit), L=0, Reserved=0 */
    gdt_set_entry(1, 0, 0xFFFFF, 0x9A, 0xCF);
    
    /* Entry 2: Kernel data segment (Ring 0, data, writable, 4GB flat): */
    /* access = 0x92 = 1001 0010: P=1, DPL=00, S=1, E=0, DC=0, RW=1, A=0 */
    gdt_set_entry(2, 0, 0xFFFFF, 0x92, 0xCF);
    
    /* Entry 3: User code segment (Ring 3, code, readable, 4GB flat): */
    /* access = 0xFA = 1111 1010: P=1, DPL=11, S=1, E=1, DC=0, RW=1, A=0 */
    gdt_set_entry(3, 0, 0xFFFFF, 0xFA, 0xCF);
    
    /* Entry 4: User data segment (Ring 3, data, writable, 4GB flat): */
    /* access = 0xF2 = 1111 0010: P=1, DPL=11, S=1, E=0, DC=0, RW=1, A=0 */
    gdt_set_entry(4, 0, 0xFFFFF, 0xF2, 0xCF);
    
    /* Entry 5: TSS (Task State Segment): */
    memset(&tss, 0, sizeof(tss));
    tss.ss0    = GDT_KERNEL_DATA;  /* Ring 0 stack segment */
    tss.esp0   = 0;                /* Set properly in Chapter 56 */
    tss.io_map = sizeof(tss);      /* No I/O permission bitmap */
    gdt_set_tss(5, (uint32_t)&tss, sizeof(tss) - 1);
    
    /* Set up the GDT descriptor: */
    gdt_desc.limit = sizeof(gdt) - 1;
    gdt_desc.base  = (uint32_t)&gdt;
    
    /* Load the new GDT and reload all segment registers: */
    gdt_flush((uint32_t)&gdt_desc);
}

/* Update the Ring 0 stack pointer in the TSS (called before switching to user mode): */
void gdt_set_kernel_stack(uint32_t esp0) {
    tss.esp0 = esp0;
}
```

---

## 5. Task State Segment (TSS) — Introduced Here

The TSS is a special data structure the CPU uses when switching privilege levels:

```
When an interrupt fires while in Ring 3 (user mode):
  CPU needs to switch to Ring 0 (kernel mode)
  CPU needs a Ring 0 stack — where does it find it?
  Answer: the TSS!
  
TSS.esp0 = address of the top of the kernel stack
TSS.ss0  = segment selector for the kernel stack (our GDT_KERNEL_DATA)

When interrupt fires in Ring 3:
  1. CPU reads TSS from GDT (using TR register — loaded by ltr instruction)
  2. CPU loads ESP from TSS.esp0, SS from TSS.ss0
  3. CPU pushes Ring 3 SS:ESP + EFLAGS + CS:EIP onto the new Ring 0 stack
  4. CPU jumps to interrupt handler

We set TSS.esp0 properly in Chapter 56 (context switching).
Here we just create the TSS entry so it doesn't crash.
```

---

## 6. Reloading Segment Registers

After loading a new GDT, we MUST reload all segment registers — they cache the old descriptor values:

```nasm
; boot/gdt_flush.asm — Reload segment registers after GDT install

[GLOBAL gdt_flush]
[BITS 32]

gdt_flush:
    ; Get the address of the GDT descriptor from the stack:
    mov eax, [esp + 4]
    
    ; Load the new GDT:
    lgdt [eax]
    
    ; Reload CS via a far jump (only way to reload CS):
    jmp 0x08:.reload_cs   ; 0x08 = kernel code selector
    
.reload_cs:
    ; Reload all data segment registers:
    mov ax, 0x10          ; 0x10 = kernel data selector
    mov ds, ax
    mov es, ax
    mov fs, ax
    mov gs, ax
    mov ss, ax
    
    ; Load the TSS descriptor into TR (Task Register):
    mov ax, 0x28          ; 0x28 = TSS selector (index 5)
    ltr ax
    
    ret
```

---

## 7. Complete gdt.c / gdt.h Summary

After this chapter, our GDT provides:

```
Index  Selector  Name             Base   Limit    DPL
  0     0x00     Null             0      0        -
  1     0x08     Kernel Code      0      4GB      0    (Ring 0 code)
  2     0x10     Kernel Data      0      4GB      0    (Ring 0 data/stack)
  3     0x18     User Code        0      4GB      3    (Ring 3 code)
  4     0x20     User Data        0      4GB      3    (Ring 3 data/stack)
  5     0x28     TSS              &tss   ~108B    0    (Task State Segment)
```

The Makefile needs the new file:
```makefile
OBJS = boot/entry.o boot/gdt_flush.o kernel/gdt.o kernel/vga.o kernel/kernel.o
```

---

## 8. Testing the GDT

Add to `kernel_main`:
```c
#include "gdt.h"

void kernel_main(uint32_t magic, uint32_t mbi_ptr) {
    terminal_init();
    
    kprintf("Installing GDT... ");
    gdt_init();
    kprintf("OK\n");
    
    /* If we reach here without a General Protection Fault, the GDT is correct! */
    kprintf("GDT installed. Kernel code selector: 0x%x\n", GDT_KERNEL_CODE);
    kprintf("Kernel data selector: 0x%x\n", GDT_KERNEL_DATA);
    
    for (;;) {}
}
```

If you see "GDT installed." without a CPU exception (triple fault / QEMU restart), the GDT is set up correctly.

**Debug tip:** If QEMU restarts immediately, there's a fault during `gdt_flush`. Use GDB:
```bash
make debug
# (gdb) break gdt_flush
# (gdb) continue
# (gdb) stepi  ← step through assembly one instruction at a time
```

---

## Summary

| Concept | Description |
|---------|------------|
| GDT entry | 8-byte descriptor: base (32-bit), limit (20-bit), access byte, flags nibble |
| Null descriptor | GDT[0] must be zero; loading it into a segment register causes GP fault |
| Access byte | P=present, DPL=privilege (0-3), S=type, E=executable, RW=readable/writable |
| Granularity flag | G=1: limit in 4KB pages; G=0: limit in bytes. With G=1, limit=0xFFFFF = 4GB |
| DB flag | D/B=1: 32-bit segment (default operand/address size). Must be 1 for protected mode code/data |
| Kernel segments | DPL=0; cannot be accessed from Ring 3 (CPU raises GP fault) |
| User segments | DPL=3; can be accessed from Ring 3; kernel also uses them (DPL check allows downgrade) |
| Selector | 16-bit: index (13 bits) + TI (1 bit, 0=GDT) + RPL (2 bits) |
| lgdt | Load GDT register with a 48-bit descriptor (limit + base address) |
| Far jump | `jmp 0x08:label` — the only way to reload CS register |
| ltr | Load Task Register with TSS selector — required for privilege-level switching |
| TSS.esp0/ss0 | Ring 0 stack that CPU switches to when interrupt fires in Ring 3 |
