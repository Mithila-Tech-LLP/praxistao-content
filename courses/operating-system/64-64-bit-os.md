# Chapter 64: 64-Bit OS and Long Mode

> **"The jump to 64-bit is not just about bigger numbers — it is about escaping the 4GB address space ceiling, gaining 16 general-purpose registers, and entering a world where the kernel no longer fights against a patchwork of memory tricks. Long mode is the native home of modern operating systems: macOS, Linux, Windows all run here. Understanding it completes your picture of how x86 truly works."**

---

## Table of Contents

1. [Why 64-Bit?](#1-why-64-bit)
2. [x86-64 vs x86-32 — Key Differences](#2-x86-64-vs-x86-32--key-differences)
3. [4-Level Paging — PML4](#3-4-level-paging--pml4)
4. [Long Mode Registers](#4-long-mode-registers)
5. [Entering Long Mode](#5-entering-long-mode)
6. [64-Bit GDT and Segments](#6-64-bit-gdt-and-segments)
7. [64-Bit System Calls — SYSCALL/SYSRET](#7-64-bit-system-calls--syscallsysret)
8. [64-Bit Calling Convention](#8-64-bit-calling-convention)
9. [Adapting TinyOS to 64-Bit](#9-adapting-tinyos-to-64-bit)
10. [Summary](#summary)

---

## 1. Why 64-Bit?

```
32-bit limitations:
  Virtual address space: 4GB maximum (2^32 bytes)
  Physical address space: 4GB (with PAE extension: 64GB, but awkward)
  Registers: 8 × 32-bit general purpose registers
  
  4GB is no longer enough:
    - A browser with 40 tabs can use > 4GB RAM
    - Video editing software needs > 4GB address space
    - Databases need > 4GB for in-memory indexes
    - Modern kernel itself may use > 1GB

64-bit (x86-64 / AMD64) improvements:
  Virtual address space: 48-bit addressing = 256 TB (more than enough)
  Physical address space: 52-bit physical = 4 PB
  Registers: 16 × 64-bit GPRs (RAX-R15) — twice as many, twice as wide
  New calling convention: first 6 args in registers (fewer stack operations)
  Simplified segmentation: segments mostly flat (base=0, limit ignored)
  
All modern desktop/server CPUs are 64-bit (since ~2003).
All modern OSes are 64-bit.
```

---

## 2. x86-64 vs x86-32 — Key Differences

```
Feature          32-bit (IA-32)       64-bit (x86-64)
---------------------------------------------------------------------------
Register width   32-bit               64-bit
GPR count        8 (EAX-EDI)          16 (RAX-R15)
Stack pointer    ESP                  RSP
Instruction ptr  EIP                  RIP (not directly accessible)
Addresses        32-bit (4GB max)     48-bit canonical (256TB)
Segment regs     Active (ES/DS used)  Mostly ignored (FS/GS still matter)
Page table       2-level (PD/PT)      4-level (PML4/PDPT/PD/PT)
FP/SSE          8 x87 FP registers   XMM0-XMM15 (128-bit SSE)

System calls     int $0x80            syscall/sysret instruction
Calling conv     cdecl (args on stack) System V AMD64 (args in regs: rdi,rsi,rdx,rcx,r8,r9)
PUSH/POP         32-bit               64-bit
RIP-relative    Not available        Available (used for position-independent code)
```

---

## 3. 4-Level Paging — PML4

In 64-bit mode, virtual addresses use 4 levels of page tables:

```
48-bit virtual address:
  ┌─────────┬──────────┬──────────┬──────────┬─────────────┐
  │ PML4[8] │ PDPT[8]  │  PD[8]   │  PT[8]   │ OFFSET[12]  │
  └─────────┴──────────┴──────────┴──────────┴─────────────┘
   bits[47:39] bits[38:30] bits[29:21] bits[20:12] bits[11:0]
    9 bits       9 bits      9 bits      9 bits      12 bits
    512 entries  512 entries 512 entries 512 entries  4KB

Translation:
  CR3 → PML4 (512 × 8-byte entries)
  PML4[bits[47:39]] → PDPT (Page Directory Pointer Table)
  PDPT[bits[38:30]] → PD (Page Directory)
  PD[bits[29:21]]   → PT (Page Table)
  PT[bits[20:12]]   → Physical page (4KB)
  + offset[11:0]    → Final physical address

Total addressable: 512 × 512 × 512 × 512 × 4KB = 512 TB (theoretical)
Practical: 256 TB (48-bit addresses; top and bottom 128 TB only)
```

**Canonical addresses:**
```
x86-64 requires canonical addresses:
  Bits [63:48] must be either all 0 or all 1 (sign-extended from bit 47)
  
Valid user space:  0x0000_0000_0000_0000 – 0x0000_7FFF_FFFF_FFFF  (128 TB)
Kernel space:      0xFFFF_8000_0000_0000 – 0xFFFF_FFFF_FFFF_FFFF  (128 TB)
The "hole":        0x0000_8000_0000_0000 – 0xFFFF_7FFF_FFFF_FFFF  (non-canonical)

Accessing a non-canonical address → General Protection Fault
```

**PML4 entry format (same as PTE but used for 64-bit):**
```c
typedef uint64_t pml4e_t;
typedef uint64_t pdpte_t;
typedef uint64_t pde64_t;
typedef uint64_t pte64_t;

/* 64-bit page table entry flags (same bit positions as 32-bit): */
#define PTE64_PRESENT   (1ULL << 0)
#define PTE64_WRITABLE  (1ULL << 1)
#define PTE64_USER      (1ULL << 2)
#define PTE64_ACCESSED  (1ULL << 5)
#define PTE64_DIRTY     (1ULL << 6)
#define PTE64_HUGE      (1ULL << 7)   /* 1GB pages in PDPT, 2MB pages in PD */
#define PTE64_GLOBAL    (1ULL << 8)
#define PTE64_NX        (1ULL << 63)  /* No-Execute: prevent code execution from data pages */

#define PTE64_ADDR_MASK  0x000FFFFFFFFFF000ULL  /* Bits [51:12] = physical frame address */
```

---

## 4. Long Mode Registers

```
64-bit General Purpose Registers (all 64-bit):
  RAX  RCX  RDX  RBX  RSP  RBP  RSI  RDI  (same names + R prefix)
  R8   R9   R10  R11  R12  R13  R14  R15   (8 new registers)

Sub-register access:
  RAX (64-bit) → EAX (low 32) → AX (low 16) → AH (bits 15:8), AL (bits 7:0)
  R8  (64-bit) → R8D (low 32) → R8W (low 16) → R8B (low 8)

IMPORTANT: Writing to 32-bit sub-register (e.g., EAX) ZERO-EXTENDS to 64-bit!
  mov eax, 1     ; RAX = 0x0000000000000001  (not just the low 32 bits!)
  mov ax, 1      ; RAX = 0x????????00000001  (high 32 bits unchanged)
  
Special registers:
  RIP = instruction pointer (read via RIP-relative addressing or CALL+POP trick)
  RFLAGS = 64-bit flags register
  
Segment registers:
  CS, SS, DS, ES: base and limit IGNORED in 64-bit (effectively base=0)
  FS, GS: still have programmable base via MSRs (IA32_FS_BASE, IA32_GS_BASE)
           Used by: Linux (FS = thread-local storage), Windows (GS = TEB)
```

---

## 5. Entering Long Mode

The sequence to switch from 32-bit protected mode to 64-bit long mode:

```nasm
; enter_long_mode.asm — Switch from protected mode to long mode

[BITS 32]

enter_long_mode:
    ; 1. Disable paging (if it was enabled in 32-bit mode):
    mov eax, cr0
    and eax, 0x7FFFFFFF   ; clear PG bit
    mov cr0, eax
    
    ; 2. Enable PAE (Physical Address Extension) — required for long mode:
    mov eax, cr4
    or  eax, (1 << 5)     ; set PAE bit (bit 5)
    mov cr4, eax
    
    ; 3. Load the PML4 page directory into CR3:
    mov eax, [pml4_address]
    mov cr3, eax
    
    ; 4. Set IA32_EFER.LME (Long Mode Enable) in the MSR:
    mov ecx, 0xC0000080   ; IA32_EFER MSR
    rdmsr                 ; read current value into EDX:EAX
    or  eax, (1 << 8)     ; set LME bit (Long Mode Enable)
    wrmsr                 ; write back
    
    ; 5. Enable paging AND protected mode (CR0.PG + CR0.PE):
    ;    Setting PG while LME is set activates Long Mode Active (LMA)
    mov eax, cr0
    or  eax, 0x80000001   ; set PG (bit 31) and PE (bit 0)
    mov cr0, eax
    
    ; 6. Far jump to a 64-bit code segment:
    ;    This flushes the pipeline and loads a 64-bit CS
    jmp 0x08:long_mode_entry   ; 0x08 = 64-bit kernel code selector in new GDT
    
    ; We are now in 64-bit long mode!

[BITS 64]
long_mode_entry:
    ; Set data segments (ignored in 64-bit, but must be valid):
    mov ax, 0x10
    mov ds, ax
    mov es, ax
    mov ss, ax
    
    ; Call 64-bit C kernel:
    call kernel_main_64
    
    cli
    hlt
```

---

## 6. 64-Bit GDT and Segments

The GDT is simpler in 64-bit (segments mostly ignored):

```
64-bit GDT entry differences:
  - L bit (bit 53) = 1 → 64-bit code segment
  - D/B bit (bit 54) = 0 → required when L=1
  - Base and limit are IGNORED for code/data segments in 64-bit mode
    (Except FS and GS — their bases are loaded via IA32_FS_BASE / IA32_GS_BASE MSRs)

64-bit GDT entries:
  Index 0: Null (required)
  Index 1: 64-bit kernel code  (L=1, D/B=0, P=1, DPL=0, S=1, E=1)
  Index 2: 64-bit kernel data  (same as 32-bit data — base/limit ignored anyway)
  Index 3: 64-bit user code    (L=1, D/B=0, P=1, DPL=3)
  Index 4: 64-bit user data    (P=1, DPL=3)
  Index 5: TSS (16-byte entry in 64-bit!)

64-bit TSS is 16 bytes (two GDT slots):
  RSP0-RSP2: Ring 0/1/2 stack pointers (64-bit)
  IST1-IST7: Interrupt Stack Tables (for safe double-fault handling)
  IOPB: I/O permission bitmap offset
```

```c
/* 64-bit kernel code segment flags:
   P=1, DPL=0, S=1, E=1, D/B=0, L=1 (64-bit), G=1
   access byte: 0x9A (same as 32-bit code)
   flags nibble: 0b1010 = G=1, L=1, D/B=0, AVL=0 → high nibble = 0xA */
void gdt64_set_code(gdt_entry64_t *e) {
    e->limit_low  = 0xFFFF;
    e->base_low   = 0;
    e->base_mid   = 0;
    e->access     = 0x9A;  /* P=1, DPL=0, S=1, E=1, RW=1 */
    e->flags_lim  = 0xAF;  /* G=1, L=1, D=0, limit[19:16]=0xF */
    e->base_high  = 0;
}
```

---

## 7. 64-Bit System Calls — SYSCALL/SYSRET

In 64-bit mode, `int $0x80` still works but the preferred method is the `syscall` instruction — it's much faster (no interrupt descriptor table lookup, no stack switch overhead):

```nasm
; User-side: make a system call
; System V AMD64 syscall convention:
;   RAX = syscall number
;   RDI = arg 1, RSI = arg 2, RDX = arg 3
;   R10 = arg 4 (NOT RCX — syscall clobbers RCX)
;   R8  = arg 5, R9 = arg 6
;   Return value in RAX (negative = error)

mov rax, 1          ; sys_write
mov rdi, 1          ; fd = stdout
lea rsi, [rel msg]  ; buf = message (RIP-relative addressing)
mov rdx, 13         ; count = 13
syscall             ; SYSCALL instruction
```

**How SYSCALL works:**
```
SYSCALL instruction:
  1. Saves RIP into RCX (return address)
  2. Saves RFLAGS into R11
  3. Loads RIP from IA32_LSTAR MSR (kernel syscall entry point)
  4. Loads CS from IA32_STAR MSR [bits 47:32]
  5. DOES NOT switch stacks (kernel must do this itself using IA32_GS_BASE or SWAPGS)
  6. DOES NOT save RSP (kernel must save/restore user RSP)
  
SYSRET instruction (return to user):
  1. Loads RIP from RCX
  2. Loads RFLAGS from R11 (with some bits forced)
  3. Loads CS back to user selector from IA32_STAR [bits 63:48]
  
MSR setup:
  wrmsr IA32_STAR,  (user_cs << 48) | (kernel_cs << 32)
  wrmsr IA32_LSTAR, syscall_entry_address
  wrmsr IA32_SFMASK, 0x200   ; clear IF on syscall entry (disable interrupts)
```

```c
/* Set up SYSCALL/SYSRET: */
void syscall_init_64(void) {
    /* IA32_STAR: [63:48]=user CS, [47:32]=kernel CS */
    uint64_t star = ((uint64_t)GDT_USER_CODE << 48) | ((uint64_t)GDT_KERNEL_CODE << 32);
    wrmsr(0xC0000081, star);
    
    /* IA32_LSTAR: 64-bit syscall entry point */
    wrmsr(0xC0000082, (uint64_t)syscall_entry64);
    
    /* IA32_FMASK: flags to clear on syscall entry (IF=bit 9) */
    wrmsr(0xC0000084, 0x200);
    
    /* Enable SCE (System Call Extensions) in IA32_EFER: */
    uint64_t efer;
    rdmsr(0xC0000080, &efer);
    wrmsr(0xC0000080, efer | 1);   /* Set SCE bit */
}
```

---

## 8. 64-Bit Calling Convention

System V AMD64 ABI (used by Linux and our 64-bit TinyOS):

```
Integer arguments (first 6): RDI, RSI, RDX, RCX, R8, R9
Floating point (first 8):    XMM0-XMM7
Return value:                RAX (and RDX for 128-bit returns)
Additional args:             passed on stack (right to left)

Callee-saved (must preserve): RBX, RBP, R12, R13, R14, R15
Caller-saved (may clobber):   RAX, RCX, RDX, RSI, RDI, R8, R9, R10, R11

Example: void memcpy(void *dst, const void *src, size_t n)
  In 32-bit: push n; push src; push dst; call memcpy
  In 64-bit: mov rdi, dst; mov rsi, src; mov rdx, n; call memcpy
  
64-bit is more efficient: less memory traffic (register args vs stack args)

Stack alignment: RSP must be 16-byte aligned before CALL
  (CALL pushes an 8-byte return address → RSP drops 8 bytes → misaligned)
  Function prologues often: sub rsp, 8 or push rbx to re-align
  SSE instructions require 16-byte alignment for MOVAPS (crash if misaligned!)
```

---

## 9. Adapting TinyOS to 64-Bit

Key changes needed to port our 32-bit TinyOS to 64-bit:

```
1. Boot sequence:
   - GRUB can load a 64-bit ELF (Multiboot2 supports long mode)
   - Or: boot in 32-bit, enter long mode ourselves (as shown above)
   
2. All pointers: uint32_t → uint64_t (or uintptr_t)
   struct vfs_node *node  → must be 64-bit on a 64-bit build
   
3. Page tables: 2-level → 4-level
   VMM rewritten for PML4/PDPT/PD/PT hierarchy
   
4. System calls: int $0x80 → syscall/sysret
   New stub assembly for SYSCALL handler
   
5. GDT: add L bit to code segments; TSS doubles in size
   
6. Context switch: save/restore 16 registers instead of 8
   RBX, RBP, R12, R13, R14, R15 are callee-saved

7. Calling convention: update all assembly <-> C interfaces
   Function args now in RDI, RSI, RDX, RCX, R8, R9
   
8. Addresses: identity map first 2MB (or use 2MB huge pages)
   KASLR: randomize kernel base address for security

9. Higher half kernel:
   Kernel lives at 0xFFFFFFFF80000000+ (last 2GB)
   User space at 0x0000000000000000 - 0x00007FFFFFFFFFFF
```

**The payoff:**
```
With a 64-bit kernel:
  Process virtual address spaces: up to 128 TB
  Physical RAM: up to 4 PB
  No more segmentation tricks, no more PAE complexity
  All modern tooling (compilers, debuggers, profilers) works natively
  You can run actual user programs compiled with modern GCC/Clang
```

---

## Summary

| Concept | Description |
|---------|------------|
| Long mode | x86-64 operating mode: 64-bit registers, 48-bit addresses, 4-level paging |
| IA32_EFER | MSR 0xC0000080; LME bit enables long mode; LMA bit confirms activation |
| PAE | Physical Address Extension: must be enabled (CR4.PAE=1) before entering long mode |
| 4-level paging | PML4 → PDPT → PD → PT → page; each level has 512 entries (9 bits) |
| PML4 | Page Map Level 4: root of 64-bit page table; physical address in CR3 |
| Canonical address | Bits [63:48] must be same as bit 47; non-canonical → GP fault |
| RAX-R15 | 16 64-bit GPRs in long mode (vs 8 in 32-bit) |
| L bit | GDT flag bit 53; set to 1 for 64-bit code segments; D/B must be 0 |
| FS/GS | Only segment registers with programmable base in 64-bit (via MSR) |
| SYSCALL | Fast system call instruction: saves RIP→RCX, RFLAGS→R11, jumps to LSTAR |
| IA32_LSTAR | MSR 0xC0000082: 64-bit syscall entry point address |
| System V AMD64 | 64-bit calling convention: first 6 args in RDI/RSI/RDX/RCX/R8/R9 |
| NX bit | PTE bit 63: No-Execute; prevents shellcode execution on data pages |
| 2MB huge pages | PDE with PS bit set; covers 2MB instead of mapping 512 PTEs |
