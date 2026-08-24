# Chapter 45: x86 Architecture for OS Developers

> **"You don't need to know every x86 instruction to write an OS — but you do need to understand the CPU's state machine: what mode it boots in, how it enters protected mode, which registers control memory and protection, and what happens on every interrupt. These 30 registers and 10 instructions are the vocabulary of OS development."**

---

## Table of Contents

1. [x86 CPU Modes](#1-x86-cpu-modes)
2. [General-Purpose Registers (32-bit)](#2-general-purpose-registers-32-bit)
3. [Special Registers — EFLAGS, EIP, ESP](#3-special-registers--eflags-eip-esp)
4. [Control Registers — CR0, CR3, CR4](#4-control-registers--cr0-cr3-cr4)
5. [Segment Registers in Protected Mode](#5-segment-registers-in-protected-mode)
6. [The Stack in x86](#6-the-stack-in-x86)
7. [Calling Convention (cdecl)](#7-calling-convention-cdecl)
8. [Essential Instructions for OS Dev](#8-essential-instructions-for-os-dev)
9. [Inline Assembly in C](#9-inline-assembly-in-c)
10. [Memory-Mapped I/O vs Port I/O](#10-memory-mapped-io-vs-port-io)
11. [Summary](#summary)

---

## 1. x86 CPU Modes

x86 has three main operating modes:

**Real Mode (16-bit):**
```
CPU starts here on every power-on
16-bit registers (AX, BX, CX, DX, SI, DI, SP, BP)
Physical addressing: segment:offset → segment × 16 + offset
Max addressable: 1MB (20-bit address bus)
No protection: any program can access any memory
No virtual memory
BIOS runs in real mode; we use BIOS interrupts for early I/O
```

**Protected Mode (32-bit):**
```
What modern 32-bit OSes use
32-bit registers (EAX, EBX, ECX, EDX, ESI, EDI, ESP, EBP)
Virtual memory via paging (4GB virtual address space)
Segment descriptors define code/data segments with DPL
Hardware-enforced Ring 0-3 protection
Must set up GDT before entering
```

**Long Mode (64-bit):**
```
64-bit extension by AMD (amd64), adopted by Intel
64-bit registers (RAX, RBX, ..., R8-R15)
64-bit virtual address space (48-bit usable = 256TB)
Segmentation mostly disabled (flat model)
4-level page tables
Requires: protected mode + PAE + enabling IA32_EFER.LME bit
```

**Switching modes:**
```
Real Mode → Protected Mode:
  1. Set up GDT (lgdt instruction)
  2. Set bit 0 (PE bit) in CR0
  3. Far jump to flush instruction pipeline and reload CS

Protected Mode → Long Mode:
  1. Set up 64-bit GDT
  2. Enable PAE (CR4.PAE = 1)
  3. Load PML4 address into CR3
  4. Set IA32_EFER.LME = 1 (MSR 0xC0000080)
  5. Enable paging (CR0.PG = 1)
  6. Far jump to 64-bit code segment
```

We will work in 32-bit protected mode for TinyOS (easier to start with).

---

## 2. General-Purpose Registers (32-bit)

```
32-bit:  EAX  EBX  ECX  EDX  ESI  EDI  ESP  EBP
16-bit:   AX   BX   CX   DX   SI   DI   SP   BP
8-bit hi: AH   BH   CH   DH
8-bit lo: AL   BL   CL   DL
```

**Traditional uses (also conventions):**
```
EAX: Accumulator — arithmetic, return value from functions
EBX: Base — often base pointer for data structures (callee-saved)
ECX: Counter — loop counters, repeat count for string ops
EDX: Data — multiplications/divisions (EDX:EAX for 64-bit results)
ESI: Source Index — source address for string/memory operations
EDI: Destination Index — destination address for string ops
ESP: Stack Pointer — points to top of stack (current frame)
EBP: Base Pointer — base of current stack frame (call stack frame)
```

**Extra registers in 64-bit mode:**
```
R8-R15: general-purpose 64-bit (R8D/R8W/R8B for 32/16/8-bit access)
RAX,RBX,RCX,RDX,RSI,RDI,RSP,RBP: 64-bit versions of classic registers
```

---

## 3. Special Registers — EFLAGS, EIP, ESP

**EFLAGS register:**
32 bits, each bit is a condition or control flag:
```
Bit  Name   Description
0    CF     Carry Flag — set when arithmetic produces carry/borrow
1    (res)  Always 1
2    PF     Parity Flag — set if least-significant byte has even parity
4    AF     Auxiliary Carry — used for BCD arithmetic
6    ZF     Zero Flag — set when result is zero
7    SF     Sign Flag — set when result is negative (MSB=1)
8    TF     Trap Flag — if set, CPU single-steps (raises #DB after each instruction)
9    IF     Interrupt Flag — if 0, maskable interrupts are ignored
10   DF     Direction Flag — controls string instruction direction (0=up, 1=down)
11   OF     Overflow Flag — set when signed overflow occurs
12-13 IOPL  I/O Privilege Level — minimum CPL to access I/O ports
14   NT     Nested Task — set if task is nested (old feature)
16   RF     Resume Flag — used with debug registers
17   VM     Virtual-8086 Mode — run 16-bit code in protected mode
18   AC     Alignment Check — enable alignment checks
21   ID     CPUID support flag
```

**Key flags for OS developers:**
- **IF (bit 9):** `cli` clears it (disable interrupts); `sti` sets it (enable interrupts)
- **TF (bit 8):** Used by debuggers to single-step
- **ZF (bit 6):** `jz` jumps if ZF=1; `test eax, eax; jz null_ptr` checks if null

**EIP (Instruction Pointer):**
Current instruction address. Modified by JMP, CALL, RET, and on interrupt (automatically saved/restored).

---

## 4. Control Registers — CR0, CR3, CR4

**CR0 (Control Register 0):**
Controls basic CPU features:
```
Bit 0  PE:  Protected Mode Enable (must set to enter 32-bit mode)
Bit 1  MP:  Monitor Coprocessor (FPU)
Bit 2  EM:  Emulation (no FPU, emulate)
Bit 3  TS:  Task Switched (set on task switch, cleared by CLTS)
Bit 4  ET:  Extension Type (FPU type)
Bit 5  NE:  Numeric Error (FPU exception handling)
Bit 16 WP:  Write Protect (even Ring 0 obeys R/W bit in page tables if set)
Bit 18 AM:  Alignment Mask
Bit 29 NW:  Not WriteThrough
Bit 30 CD:  Cache Disable
Bit 31 PG:  Paging Enable (must set to enable virtual memory)
```

```nasm
; Enable protected mode:
mov eax, cr0
or eax, 1       ; set PE bit
mov cr0, eax

; Enable protected mode AND paging:
mov eax, cr0
or eax, 0x80000001   ; set PE (bit 0) and PG (bit 31)
mov cr0, eax
```

**CR2 (Page Fault Linear Address):**
When a page fault occurs, CR2 contains the virtual address that caused the fault.
```c
// In page fault handler (ISR #14):
uint32_t fault_addr;
asm("mov %%cr2, %0" : "=r"(fault_addr));
printf("Page fault at address: 0x%08X\n", fault_addr);
```

**CR3 (Page Directory Base Register):**
Physical address of the page directory (or PML4 in 64-bit mode). Changing CR3 switches the entire virtual address space.
```nasm
mov eax, page_directory_physical_addr
mov cr3, eax   ; switch to new page tables (also flushes TLB)
```

**CR4:**
Extended CPU features:
```
Bit 0  VME:  Virtual 8086 mode extensions
Bit 1  PVI:  Protected-mode virtual interrupts
Bit 2  TSD:  Timestamp disable
Bit 3  DE:   Debug extensions
Bit 4  PSE:  Page Size Extension (4MB pages)
Bit 5  PAE:  Physical Address Extension (enables 64-bit paging in 32-bit mode)
Bit 6  MCE:  Machine Check Enable
Bit 7  PGE:  Page Global Enable (global TLB entries)
Bit 9  OSFXSR: SSE instructions support
Bit 20 SMEP: Supervisor Mode Execution Prevention
Bit 21 SMAP: Supervisor Mode Access Prevention
```

---

## 5. Segment Registers in Protected Mode

In 32-bit protected mode, segment registers hold **selectors** (not base addresses like in real mode):

```
CS — Code Segment: CPU fetches instructions using CS:EIP
DS — Data Segment: default for most memory access
SS — Stack Segment: stack operations (PUSH, POP, CALL, RET)
ES — Extra Segment: string operation destination (MOVS, STOS)
FS, GS — General purpose (often used for TLS)
```

**Selector format:**
```
15        3  2  1 0
[  Index    |TI|RPL]
  13 bits    1  2

Index: index into GDT (TI=0) or LDT (TI=1)
RPL: Requested Privilege Level (0-3)

Examples:
  0x08: index=1, TI=0, RPL=0 → GDT[1], ring 0 (kernel code)
  0x10: index=2, TI=0, RPL=0 → GDT[2], ring 0 (kernel data)
  0x18: index=3, TI=0, RPL=0 → GDT[3], ring 3 code (but selector says ring 0!)
  0x23: index=4, TI=0, RPL=3 → GDT[4], ring 3 (user data)
```

**Flat segmentation (what we'll use):**
All segments have base=0, limit=4GB — effectively disabled:
```
GDT[0]: null descriptor (required)
GDT[1]: kernel code: base=0, limit=4GB, type=code, DPL=0, CS=0x08
GDT[2]: kernel data: base=0, limit=4GB, type=data, DPL=0, DS=0x10
GDT[3]: user code:   base=0, limit=4GB, type=code, DPL=3, CS=0x18|3=0x1B
GDT[4]: user data:   base=0, limit=4GB, type=data, DPL=3, DS=0x20|3=0x23
GDT[5]: TSS descriptor: base=&tss, type=TSS, DPL=0
```

---

## 6. The Stack in x86

x86 stack grows DOWNWARD:
```
High address → low address as you push things

Initial: ESP = 0x7FFFF000

PUSH EAX:
  ESP -= 4
  [ESP] = EAX value
  
POP EBX:
  EBX = [ESP]
  ESP += 4

Call sequence:
  Before CALL:  ESP → [args if any] → caller's locals → saved frame...
  
  CALL target:
    Push EIP (return address) → ESP -= 4, [ESP] = EIP
    Jump to target
    
  Function prologue (generated by compiler):
    PUSH EBP        ; save old frame pointer
    MOV EBP, ESP    ; set new frame pointer
    SUB ESP, 16     ; allocate 16 bytes of locals
    
  Stack during function:
    [EBP + 8]  first argument
    [EBP + 4]  return address
    [EBP]      saved EBP
    [EBP - 4]  local variable 1
    [EBP - 8]  local variable 2
```

---

## 7. Calling Convention (cdecl)

**cdecl** is the standard C calling convention on x86-32:

```
Arguments: pushed onto stack RIGHT to LEFT before CALL
Return value: in EAX (or EDX:EAX for 64-bit values)
Caller cleans up stack (adds back the argument bytes after CALL)
Callee-saved: EBX, ESI, EDI, EBP (function must preserve these)
Caller-saved: EAX, ECX, EDX (function can modify freely)
```

```c
// C call: add(3, 4)
// Generated assembly:
push 4          ; second argument
push 3          ; first argument
call add        ; push return address, jump to add
add esp, 8      ; caller cleans up (2 args × 4 bytes = 8)

// Inside add():
push ebp        ; save caller's frame pointer
mov ebp, esp
mov eax, [ebp+8]  ; first argument (3)
add eax, [ebp+12] ; second argument (4)
pop ebp
ret             ; return, EAX = 7
```

---

## 8. Essential Instructions for OS Dev

```nasm
; Privilege/mode control:
cli              ; clear interrupt flag (disable maskable interrupts)
sti              ; set interrupt flag (enable maskable interrupts)
hlt              ; halt CPU until next interrupt

; GDT/IDT loading:
lgdt [gdtr]      ; load GDT register from memory descriptor
lidt [idtr]      ; load IDT register
ltr ax           ; load task register (TSS selector)

; Far jump (flushes pipeline, reloads CS):
jmp 0x08:protected_mode_start    ; jump to new CS:EIP

; Register moves:
mov eax, cr0     ; read CR0
mov cr0, eax     ; write CR0
mov eax, cr3     ; read page directory base
mov cr3, eax     ; set page directory (flushes TLB)

; I/O:
in al, 0x60      ; read 1 byte from port 0x60 (keyboard data)
out 0x20, al     ; write 1 byte to port 0x20 (PIC command)
in ax, 0x1F0     ; read 2 bytes from port 0x1F0 (IDE data)

; String operations (useful for memcpy/memset):
; ESI = source, EDI = destination, ECX = count
rep movsb        ; copy ECX bytes from [ESI] to [EDI], incrementing both
rep stosb        ; fill ECX bytes at [EDI] with AL
rep stosd        ; fill ECX dwords at [EDI] with EAX

; Interrupt return:
iret             ; return from interrupt (pops EIP, CS, EFLAGS [and ESP, SS if privilege change])
iretd            ; same as iret in 32-bit mode

; CPUID:
mov eax, 0       ; get vendor ID
cpuid
; EBX:EDX:ECX = vendor string ("GenuineIntel" or "AuthenticAMD")

; MSR:
mov ecx, 0xC0000080  ; IA32_EFER (Extended Feature Enable Register)
rdmsr               ; read MSR: result in EDX:EAX
wrmsr               ; write MSR: value from EDX:EAX
```

---

## 9. Inline Assembly in C

GCC allows embedding assembly inside C code:

```c
// Read from I/O port:
static inline uint8_t inb(uint16_t port) {
    uint8_t value;
    asm volatile ("inb %1, %0" : "=a"(value) : "Nd"(port));
    return value;
}

// Write to I/O port:
static inline void outb(uint16_t port, uint8_t value) {
    asm volatile ("outb %0, %1" : : "a"(value), "Nd"(port));
}

// Read CR2 (page fault address):
static inline uint32_t read_cr2(void) {
    uint32_t val;
    asm volatile ("mov %%cr2, %0" : "=r"(val));
    return val;
}

// Load GDT:
static inline void lgdt(void *gdtr) {
    asm volatile ("lgdt (%0)" : : "r"(gdtr));
}

// Invalidate a TLB entry for a virtual address:
static inline void invlpg(void *addr) {
    asm volatile ("invlpg (%0)" : : "r"(addr) : "memory");
}

// Disable/enable interrupts:
static inline void disable_interrupts(void) {
    asm volatile ("cli");
}

static inline void enable_interrupts(void) {
    asm volatile ("sti");
}
```

**Inline assembly syntax:**
```c
asm volatile ("instruction" : outputs : inputs : clobbers);

Constraints:
  "=a" → output to EAX
  "a"  → input from EAX
  "r"  → any general-purpose register
  "Nd" → immediate or DX register (for port numbers)
  "memory" → clobber: compiler must reload all memory after this
  
volatile: prevents compiler from optimizing away or moving the asm
```

---

## 10. Memory-Mapped I/O vs Port I/O

**Port I/O (PIO):**
x86 has a separate 64KB I/O address space. Use `in`/`out` instructions:
```c
// Legacy devices (PIC, PIT, PS/2 keyboard, CMOS RTC):
uint8_t data = inb(0x60);   // keyboard data port
outb(0x20, 0x20);           // send EOI to master PIC
```

**Memory-Mapped I/O (MMIO):**
Device registers appear as normal memory addresses. Use regular load/store:
```c
// APIC registers (at 0xFEE00000 by default):
volatile uint32_t *lapic = (volatile uint32_t *)0xFEE00000;
uint32_t lapic_id = lapic[0x20 / 4];  // LAPIC ID register at offset 0x20
lapic[0xB0 / 4] = 0;                  // send EOI to APIC (write 0 to offset 0xB0)

// VGA text mode (at 0xB8000):
volatile char *vga = (volatile char *)0xB8000;
vga[0] = 'H';  vga[1] = 0x0F;   // white 'H' on black background
vga[2] = 'i';  vga[3] = 0x0F;   // white 'i'
```

**`volatile` is critical:**
Without `volatile`, the compiler might optimize away your writes (assuming the memory location doesn't change between reads or writes to the same location). Device registers change independently of your code — always use `volatile` for MMIO.

---

## Summary

| Concept | Description |
|---------|------------|
| Real mode | 16-bit boot mode; no protection; BIOS lives here; 1MB max |
| Protected mode | 32-bit; segmentation + paging; Ring 0-3 protection |
| Long mode | 64-bit; flat segmentation; 4-level paging |
| CR0.PE | Protected Mode Enable bit; set to enter protected mode |
| CR0.PG | Paging Enable; set to activate virtual memory |
| CR3 | Page Directory Base Register; physical address of page directory |
| CR2 | Fault address; CPU stores faulting virtual address here on page fault |
| EFLAGS.IF | Interrupt enable flag; `cli` clears (disables), `sti` sets (enables) |
| Selector | 16-bit value in segment register: index into GDT + RPL |
| cdecl | C calling convention: args right-to-left, caller cleans up, EAX = return value |
| inb/outb | Read/write to x86 I/O port space (64KB separate address space) |
| volatile | Tell compiler: this memory changes without my explicit action; don't optimize |
| MMIO | Device registers mapped to physical memory; access with normal load/store |
| `lgdt` | Load GDT register; takes pointer to {limit, base} struct |
| `lidt` | Load IDT register |
| `iret` | Return from interrupt; restores EIP, CS, EFLAGS (and ESP, SS if privilege changed) |
