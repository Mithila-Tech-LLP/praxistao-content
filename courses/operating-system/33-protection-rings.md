# Chapter 33: Protection Rings and Privilege Levels

> **"The x86 ring system is a hardware-enforced wall between code that can do anything and code that is supposed to do almost nothing. When software respects this wall, systems are stable and secure. When exploits break through it — via kernel vulnerabilities or hypervisor escapes — the consequences are catastrophic."**

---

## Table of Contents

1. [Why Protection Rings?](#1-why-protection-rings)
2. [The x86 Ring Model](#2-the-x86-ring-model)
3. [Ring 0 — Kernel Mode](#3-ring-0--kernel-mode)
4. [Ring 3 — User Mode](#4-ring-3--user-mode)
5. [Privilege Level Enforcement](#5-privilege-level-enforcement)
6. [Transitioning Between Rings](#6-transitioning-between-rings)
7. [Hypervisor Rings — Ring -1](#7-hypervisor-rings--ring--1)
8. [ARM's Exception Levels](#8-arms-exception-levels)
9. [Summary](#summary)

---

## 1. Why Protection Rings?

Without hardware protection:
```
Any program can:
  - Write to any memory address (overwrite kernel)
  - Execute any CPU instruction (halt CPU, disable interrupts)
  - Access any I/O device directly (read your disk, send network packets)
  
One buggy or malicious program → entire system compromised
MS-DOS had no protection → one bad program crashed everything
```

With hardware protection rings:
```
User programs can ONLY:
  - Access their own memory
  - Execute unprivileged instructions
  - Request OS services via system calls (controlled entry points)
  
Kernel enforces ALL resource access
A buggy user program → only crashes itself (usually)
```

---

## 2. The x86 Ring Model

x86 defines 4 privilege levels (rings), implemented in hardware:

```
Ring 0 ─────────────────────────────────── Most Privileged
  │                                         (can do everything)
  │   Kernel code
  │   Device drivers (in monolithic kernels)
  │   
Ring 1 ─────────────────────────────────── (unused by modern OSes)
Ring 2 ─────────────────────────────────── (unused by modern OSes)
  │
Ring 3 ─────────────────────────────────── Least Privileged
        User applications                  (very limited)
```

**Why only Ring 0 and Ring 3 in practice?**
- Unix/Linux/Windows designed for 2 privilege levels (simpler)
- Rings 1 and 2 were meant for OS services less privileged than the kernel but more privileged than apps (like device drivers in OS/2)
- Modern OSes put all trusted code in Ring 0 and all user code in Ring 3

---

## 3. Ring 0 — Kernel Mode

**What Ring 0 code can do:**

**Privileged CPU instructions:**
```
LGDT    - Load Global Descriptor Table (set up segment table)
LLDT    - Load Local Descriptor Table
LIDT    - Load Interrupt Descriptor Table (set up interrupt handlers)
LTR     - Load Task Register (set current TSS)
MOV CR0 - Write to control registers (enable/disable paging, protected mode)
MOV CR3 - Switch page tables (change virtual address space)
MOV CR4 - Enable/disable CPU features (PAE, SMEP, SMAP)
RDMSR   - Read Model Specific Register (access CPU configuration)
WRMSR   - Write MSR (configure syscall entry point, CPU features)
IN/OUT  - Access I/O ports directly (talk to legacy hardware)
INVLPG  - Invalidate a TLB entry
CLTS    - Clear task-switched flag
HLT     - Halt the CPU (wait for interrupt)
STI     - Set interrupt flag (enable interrupts)
CLI     - Clear interrupt flag (disable interrupts)
```

**Attempting any of these from Ring 3:**
```
Attempts to execute privileged instruction from user mode
→ CPU detects: CPL (current privilege level) = 3
→ Instruction requires CPL = 0
→ General Protection Fault (#GP, interrupt 13)
→ OS delivers SIGSEGV to the process (or kills it)
```

**Memory access:**
Ring 0 can access all physical memory (including kernel-only pages with U/S=0 in PTE).

---

## 4. Ring 3 — User Mode

**What Ring 3 code can do:**

**Allowed instructions:**
- Arithmetic: ADD, SUB, MUL, DIV, etc.
- Data movement: MOV, PUSH, POP
- Control flow: JMP, CALL, RET, conditional branches
- String operations: MOVS, CMPS, SCAS
- FPU/SSE/AVX instructions
- RDTSC (read timestamp counter — usually allowed, sometimes restricted)

**Memory access:**
Ring 3 code can ONLY access pages where the U/S (User/Supervisor) bit = 1 in the page table entry. Kernel pages have U/S=0 → accessing them causes a Page Fault (#PF).

**Can NOT:**
- Access I/O ports directly (IOPL must be 3, usually not granted)
- Execute privileged instructions
- Modify page tables (would require writing to kernel memory)
- Change interrupt handlers (would require LIDT)

---

## 5. Privilege Level Enforcement

**CPL (Current Privilege Level):**
The CPU's current ring, stored in the lower 2 bits of the CS segment register:
- CS = 0x08 (kernel code segment): CPL = 0
- CS = 0x1B (user code segment): CPL = 3

**DPL (Descriptor Privilege Level):**
The required privilege level to access a segment or call gate, stored in the segment descriptor.

**RPL (Requested Privilege Level):**
In the selector loaded into a segment register, lower 2 bits.

**Access check:**
```
Can code with CPL access a resource with DPL?

For code (call/jump to segment):
  CPL ≤ DPL → allowed (numerically, Ring 0 = 0 is highest)

For data (access segment via DS/ES):
  max(CPL, RPL) ≤ DPL → allowed

Examples:
  Ring 3 code (CPL=3) trying to load kernel data segment (DPL=0):
    max(3, RPL) ≤ 0? → 3 > 0 → DENIED → #GP exception

  Ring 0 code (CPL=0) accessing user data segment (DPL=3):
    max(0, 3) ≤ 3? → 3 ≤ 3 → allowed (kernel can access user memory)
    BUT: SMAP (Supervisor Mode Access Prevention) blocks this by default!
    → Kernel must use special instructions (copy_to_user, copy_from_user)
```

**Page table check (different from segment check):**
Even if segment check passes, the page table independently enforces:
- U/S bit = 0: only Ring 0 can access
- R/W bit = 0: read-only (Ring 0 can still write, unless WP flag set in CR0)
- NX bit = 1: not executable (any ring)

---

## 6. Transitioning Between Rings

**Ring 3 → Ring 0 (user to kernel):**

Cannot jump directly (would be a security violation). Must use a controlled gate:

**1. SYSCALL/SYSRET (x86-64, fast path):**
```
syscall instruction:
  1. Save RIP (return address) in RCX
  2. Save RFLAGS in R11
  3. Load CS with kernel code selector (Ring 0) from IA32_STAR MSR
  4. Load RIP from IA32_LSTAR MSR (kernel entry point)
  5. Clear RFLAGS bits (IF, etc.)
  6. Stack not automatically switched — kernel must do it manually
  
sysret instruction (return path):
  1. Restore RIP from RCX
  2. Restore RFLAGS from R11
  3. Load CS with user code selector (Ring 3)
  4. Continue user code at saved RIP
```

**2. Interrupt gate (IRQ/exception):**
```
Interrupt fires:
  1. CPU reads IDT[vector] gate descriptor
  2. If DPL check passes (int 0x80 for syscall has DPL=3, hardware IRQs have DPL=0)
  3. If privilege change (Ring 3 → Ring 0): switch stack to kernel stack
     → read RSP0 from TSS (set up by kernel for this process)
  4. Push SS, RSP, RFLAGS, CS, RIP (and error code for some exceptions)
  5. Load CS = kernel code selector (Ring 0)
  6. Jump to handler address from IDT

iretq (return):
  1. Pop RIP, CS, RFLAGS, RSP, SS
  2. CS has Ring 3 selector → CPL becomes 3
  3. Continue user code
```

**Ring 0 → Ring 3 (kernel to user):**
The kernel uses SYSRET or IRETQ to return to user mode.

For the FIRST time launching a user process (before any syscall), the kernel manually constructs a fake interrupt frame and executes IRETQ to "return" to user mode for the first time.

---

## 7. Hypervisor Rings — Ring -1

**Virtualization problem:**
A virtual machine's guest OS kernel runs in Ring 0 and executes privileged instructions. But we want to run MULTIPLE guest OSes in one physical machine — only one kernel can really be in Ring 0.

**Hardware virtualization (VT-x / AMD-V):**
Intel and AMD added a new privilege level:
- **VMX root mode (Ring -1):** The hypervisor/host OS — can do anything
- **VMX non-root mode (Ring 0-3):** Guest OS runs here

```
Physical CPU:
  VMX Root Mode (Ring -1): Host OS / Hypervisor (KVM, VMware ESXi)
    VM 1 context:
      VMX Non-Root Ring 0: Guest OS kernel (Linux, Windows)
      VMX Non-Root Ring 3: Guest user applications
    VM 2 context:
      VMX Non-Root Ring 0: Another guest OS kernel
      VMX Non-Root Ring 3: Guest user apps
```

**How guest Ring 0 is controlled:**
When a guest kernel executes a privileged instruction (MOV CR3, VMWRITE, etc.):
1. CPU traps to the hypervisor (VM Exit)
2. Hypervisor handles: emulate the instruction safely
3. Hypervisor returns to guest (VM Entry)

Guest kernel thinks it's Ring 0. Hypervisor actually controls what it does.

**VMCS (VM Control Structure):**
A per-VM data structure controlling what triggers a VM exit:
```
VMCS fields:
  Pin-based controls: which interrupts cause VM exit
  Exec controls: which instructions cause VM exit (CR3 writes, port I/O, etc.)
  Exit controls: what state to save on exit
  Entry controls: what state to restore on entry
  Guest state: saved registers, segment state when in guest
  Host state: where to go on VM exit (hypervisor entry point)
```

---

## 8. ARM's Exception Levels

ARM64 uses a different model — **Exception Levels (EL)**:

```
EL3: Secure Monitor (TrustZone — runs secure OS for TEE)
EL2: Hypervisor (KVM on ARM, Xen)
EL1: OS Kernel (Linux kernel)
EL0: User applications

Normal world (non-secure):
  EL0 ↔ EL1: user ↔ kernel (via SVC instruction = ARM syscall)
  EL1 ↔ EL2: kernel ↔ hypervisor (via HVC instruction)
  
Secure world (TrustZone):
  EL0s: Trusted application (e.g., fingerprint processing)
  EL1s: Trusted OS (e.g., OP-TEE)
  EL3: Secure Monitor (switches between normal/secure worlds)
```

**TrustZone:**
ARM's security model splits the CPU into two "worlds":
- **Normal world:** Android OS + apps (untrusted)
- **Secure world:** Trusted Execution Environment (TEE) — handles keys, biometrics

When you verify a fingerprint on Android, your data travels to the secure world for processing. Even if Android is compromised, the biometric key material never leaves the secure world.

---

## Summary

| Concept | Description |
|---------|------------|
| Ring 0 | Kernel mode; can execute all instructions; unrestricted memory access |
| Ring 3 | User mode; limited instructions; can only access user-mapped pages |
| CPL | Current Privilege Level: bits 1:0 of CS register |
| DPL | Descriptor Privilege Level: required ring to access a segment |
| RPL | Requested Privilege Level: caller's privilege claim in a selector |
| Privileged instruction | Instruction that faults if CPL != 0 (LGDT, CR0 writes, CLI, HLT, etc.) |
| SYSCALL | x86-64 fast user→kernel transition; uses IA32_LSTAR MSR |
| SYSRET | x86-64 fast kernel→user transition |
| Interrupt gate | Controlled Ring 3→Ring 0 via IDT (saves full state) |
| IRETQ | Return from interrupt; restores CS (ring) + all saved state |
| VMX root | x86 hypervisor privilege level; controls guest VMs |
| VMX non-root | Guest VM context; hardware restricts privileged instructions |
| VM Exit | Trap from guest to hypervisor (on privileged instruction or event) |
| ARM EL0-EL3 | ARM exception levels: EL0=user, EL1=kernel, EL2=hypervisor, EL3=secure monitor |
| TrustZone | ARM security extension: normal world vs. secure world |
