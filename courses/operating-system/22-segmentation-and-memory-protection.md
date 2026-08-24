# Chapter 22: Segmentation and Memory Protection

> **"Segmentation was the original solution to organizing a program's memory into logical sections — code, data, stack. Paging replaced it for flexibility, but segments didn't disappear. They live on in the x86 GDT, in the concept of code vs. data protections, and in the way the OS enforces boundaries between kernel and user, between processes, and between different types of data."**

---

## Table of Contents

1. [Segmentation — Logical Memory Regions](#1-segmentation--logical-memory-regions)
2. [Hardware Segmentation (x86)](#2-hardware-segmentation-x86)
3. [The GDT — Global Descriptor Table](#3-the-gdt--global-descriptor-table)
4. [Segment Descriptor Format](#4-segment-descriptor-format)
5. [Privilege Levels and Segment Protection](#5-privilege-levels-and-segment-protection)
6. [Segmentation vs. Paging](#6-segmentation-vs-paging)
7. [Modern OSes: Segmentation in Practice](#7-modern-oses-segmentation-in-practice)
8. [Memory Protection: The Combined Picture](#8-memory-protection-the-combined-picture)
9. [NX/XD Bit and W^X Policy](#9-nxxd-bit-and-wx-policy)
10. [ASLR — Address Space Layout Randomization](#10-aslr--address-space-layout-randomization)
11. [Summary](#summary)

---

## 1. Segmentation — Logical Memory Regions

**The idea:**
A program naturally divides into logical sections:
- Code (text): instructions — read-only, executable
- Data: global variables — read-write, not executable
- Stack: local variables — read-write, grows dynamically
- Heap: dynamic allocations — read-write, grows dynamically

**Segmentation** gives each of these sections its own **segment** — a named, typed memory region with its own base address, size limit, and access permissions.

```
Segment register → segment descriptor → {base, limit, type, DPL}

Physical address = base + offset
Access check: offset < limit? type matches access type? DPL allows?
```

**Why it was invented:**
In early 16-bit x86 (real mode), the address bus was 20 bits (1MB) but registers were 16 bits (64KB). Segments allowed addressing the full 1MB:
```
Physical address = segment_register × 16 + offset
```

This is real-mode segmentation — just arithmetic, no protection.

---

## 2. Hardware Segmentation (x86)

In 32-bit **protected mode**, x86 has a real hardware segmentation mechanism with protection:

**Segment registers:**
```
CS (Code Segment):    points to segment containing the current code
DS (Data Segment):    points to segment for general data access
SS (Stack Segment):   points to segment for stack operations (push/pop)
ES, FS, GS:           extra segments (general purpose)
```

Each segment register holds a **selector** — a 16-bit value:
```
Selector (16 bits):
┌────────────────────┬─┬──┐
│  Index (13 bits)   │TI│RPL│
└────────────────────┴─┴──┘
  Bits 15-3           Bit 2  Bits 1-0
  Index into GDT      Table  Requested Privilege Level
  or LDT              0=GDT  0=Ring0 to 3=Ring3
```

**Address translation with segmentation:**
```
1. CPU executes instruction at CS:EIP
   CS selector → look up segment descriptor in GDT
   Segment descriptor: {base=0x08000000, limit=0xFFFFF, type=code, DPL=3}
   
2. Check: RPL in CS ≥ DPL? (can we access this segment?)
   Check: EIP < limit? (within bounds?)
   Check: access type (execute) matches segment type (code)?
   
3. Physical (linear) address = base + EIP = 0x08000000 + EIP
   This linear address then goes through paging (if enabled)
```

**Two levels of address translation (x86 protected mode):**
```
Virtual address → [Segmentation] → Linear address → [Paging] → Physical address
(segment:offset)                  (base + offset)              (frame + offset)
```

Modern 64-bit x86 (long mode) largely bypasses segmentation — segment bases are forced to 0 (except FS and GS which are used for TLS). So in 64-bit mode: linear address = virtual address (segmentation is a no-op).

---

## 3. The GDT — Global Descriptor Table

The **GDT (Global Descriptor Table)** is an array of segment descriptors stored in memory. The CPU register **GDTR** points to it.

**Loading the GDT:**
```c
struct gdtr {
    uint16_t limit;   // GDT size - 1
    uint32_t base;    // linear address of the GDT
} __attribute__((packed));

struct gdtr gdt_register = { sizeof(gdt) - 1, (uint32_t)gdt };
asm volatile("lgdt (%0)" : : "r"(&gdt_register));
```

**GDT layout for a simple OS:**
```
GDT[0]: Null descriptor (required by spec — never used)
GDT[1]: Kernel code segment (Ring 0, executable, 0-4GB)
GDT[2]: Kernel data segment (Ring 0, read-write, 0-4GB)
GDT[3]: User code segment   (Ring 3, executable, 0-4GB)
GDT[4]: User data segment   (Ring 3, read-write, 0-4GB)
GDT[5]: TSS descriptor      (Task State Segment — for kernel stack on interrupt)
```

**Segment selectors for these:**
```c
#define SEG_KERNEL_CODE  0x08  // index 1, GDT, Ring 0
#define SEG_KERNEL_DATA  0x10  // index 2, GDT, Ring 0
#define SEG_USER_CODE    0x1B  // index 3, GDT, Ring 3
#define SEG_USER_DATA    0x23  // index 4, GDT, Ring 3
// Note: user selectors have RPL=3 (bits 1-0 = 11)
```

---

## 4. Segment Descriptor Format

Each GDT entry is 8 bytes with this format:

```
64-bit descriptor (8 bytes):
Byte 7: Base[31:24]  ┐
Byte 6: G, DB, L, AVL, SegLimit[19:16]  ├── Upper 4 bytes
Byte 5: P, DPL, S, Type  │
Byte 4: Base[23:16]  ┘

Byte 3: Base[15:8]   ┐
Byte 2: Base[7:0]    ├── Lower 4 bytes
Byte 1: SegLimit[15:8]│
Byte 0: SegLimit[7:0] ┘

Fields:
Base:  32-bit base address of the segment
Limit: 20-bit limit (segment size - 1; if G=1, unit is 4KB pages → 4GB max)
G:     Granularity (0 = byte granularity, 1 = 4KB page granularity)
DB:    0 = 16-bit, 1 = 32-bit operand/address size
L:     0 = compatibility mode, 1 = 64-bit code segment
P:     Present bit (1 = valid descriptor)
DPL:   Descriptor Privilege Level (0-3)
S:     0 = system descriptor (TSS, LDT, gate), 1 = code or data
Type:  code/data type and access bits
```

**Common segment types:**
```
Type (4 bits) for S=1 (code/data):
0b0010: Read-Write data segment (no execute)
0b1010: Execute-Read code segment (no write)
0b1000: Execute-Only code segment
```

**Writing a code segment descriptor in C:**
```c
// Create a descriptor for kernel code segment:
// Base = 0, Limit = 0xFFFFF (full 4GB with G=1), DPL=0
// S=1 (code/data), Type=0xA (execute/read), P=1, G=1, DB=1

uint64_t kernel_code_descriptor = 
    (0x0000000000000000ULL)  // base 0
    | (0x000FFFFFULL)         // limit 0xFFFFF
    | (1ULL << 43)            // type bit 3: code segment
    | (1ULL << 41)            // type bit 1: readable
    | (1ULL << 44)            // S=1: code/data descriptor
    | (0ULL << 45)            // DPL=0
    | (1ULL << 47)            // P=1: present
    | (0xFULL << 48)          // limit upper 4 bits
    | (1ULL << 54)            // DB=1: 32-bit
    | (1ULL << 55);           // G=1: 4KB granularity
```

---

## 5. Privilege Levels and Segment Protection

x86 has 4 privilege levels (rings):

```
Ring 0: Kernel — highest privilege (can execute all instructions)
Ring 1: Unused in modern OSes
Ring 2: Unused in modern OSes
Ring 3: User applications — lowest privilege
```

**Protection rules:**
- Code in Ring 3 CANNOT jump to a Ring 0 code segment (direct jump)
- Code in Ring 3 CANNOT load a Ring 0 data segment into DS/SS
- Exception: controlled entry via **call gate**, **interrupt gate**, or **syscall** instruction

**When a user process tries to access kernel memory:**
```
MOV EAX, [kernel_address]
  → CPU checks: DS descriptor DPL = 0 (kernel)
  → CS CPL = 3 (user is in Ring 3)
  → CPL > DPL → ACCESS DENIED
  → General Protection Fault (#GP) → OS kills the process
```

**The TSS (Task State Segment):**
When an interrupt occurs while in Ring 3 (user mode), the CPU needs to switch to Ring 0 (kernel stack). Where is the kernel stack?

The TSS descriptor in the GDT points to the TSS, which contains the kernel stack pointer (ESP0):
```c
struct tss {
    // ...
    uint32_t esp0;  // kernel stack pointer for Ring 0
    uint16_t ss0;   // kernel stack segment
    // ...
};
```

On every user → kernel transition (interrupt, syscall), the CPU reads ESP0 from the TSS and sets up the kernel stack.

---

## 6. Segmentation vs. Paging

| Feature | Segmentation | Paging |
|---------|-------------|--------|
| Unit size | Variable (can be any size) | Fixed (4KB pages) |
| External fragmentation | Yes (variable sizes) | No (fixed size frames) |
| Internal fragmentation | No (exact fit) | Yes (up to 4KB per allocation) |
| Sharing | Easy (share a segment) | Fine-grained (share pages) |
| Protection | Per-segment (code/data/stack) | Per-page (R/W/X/U/K) |
| Address space | Limited to segment model | Flexible virtual address space |
| Hardware support | Complex (segment registers, GDT) | Simpler (one register CR3) |

**Modern verdict:** Paging won. Modern OSes use paging for memory management.

Segmentation survives as:
- The Ring 0/3 protection mechanism (code segment DPL)
- The TSS for kernel stack on ring transitions
- FS/GS for Thread Local Storage
- The conceptual boundary between code, data, and stack (enforced by page permissions, not segment limits)

---

## 7. Modern OSes: Segmentation in Practice

**Linux 64-bit:**
```c
// All segment bases = 0 in Linux x86-64
// Segmentation is effectively disabled (flat address space)
// Protection is done entirely by paging (page table permission bits)
// Exceptions:
//   FS segment: base = Thread Local Storage area (per-thread)
//   GS segment: base = per-CPU data (kernel) or per-thread data (user)
```

**Windows 64-bit:**
Same approach: flat segmentation (all bases = 0), TLS via GS.

**Why keep segment registers at all?**
- FS/GS for TLS: `mov rax, fs:[0x28]` reads thread's stack canary
- CS DPL: still enforces Ring 0 vs Ring 3 transition (kernel vs user mode)
- IA32_LSTAR MSR: `syscall` instruction uses this to jump to kernel (replaces old segment-based call gates)

---

## 8. Memory Protection: The Combined Picture

Modern x86-64 uses multiple layers of protection together:

```
Layer 1: Segmentation
  - CS DPL = 0 (kernel code) vs. DPL = 3 (user code)
  - Prevents user code from executing as kernel code

Layer 2: Paging permission bits
  - U/S bit: kernel pages inaccessible from Ring 3
  - R/W bit: read-only pages (code, COW)
  - NX bit: non-executable pages (stack, heap)

Layer 3: SMEP (Supervisor Mode Execution Prevention)
  - CPU feature: Ring 0 code CANNOT execute pages with U/S=1
  - Prevents kernel from executing user-supplied code (exploit mitigation)

Layer 4: SMAP (Supervisor Mode Access Prevention)
  - CPU feature: Ring 0 code CANNOT read user pages (unless STAC instruction used)
  - Prevents kernel bugs from reading arbitrary user memory

Layer 5: KPTI (Kernel Page Table Isolation)
  - Separate page tables for user and kernel (Meltdown mitigation)
  - User page table has minimal kernel mappings
```

Together, these layers make it extremely difficult for exploits to gain kernel privileges.

---

## 9. NX/XD Bit and W^X Policy

**NX (No-Execute) bit** in page table entries prevents executing code from data pages.

**W^X (Write XOR Execute) policy:**
At any time, a page is either:
- Writable (can be modified) OR
- Executable (can contain code to run)
- NEVER BOTH simultaneously

```
Memory region     W   X   Reason
Kernel code       0   1   Can execute, can't modify
User code         0   1   Execute only (integrity)
Stack             1   0   Can write locals, can't execute code
Heap              1   0   Can allocate data, can't execute shellcode
Mapped libraries  0   1   Shared read-only executables
JIT buffer        0   1   First write code (W), then flip to X before execution
```

**JIT compilation and W^X:**
Just-in-Time compilers (JavaScript engines, Java JVM) need to generate and execute code:
```
1. Allocate anonymous page (PROT_WRITE | PROT_READ only, NOT EXEC)
2. Write machine code into the page
3. mprotect() the page to PROT_READ | PROT_EXEC (remove write)
4. Execute the code
```

Steps 2 and 4 are never simultaneous → W^X preserved.

Modern iOS, macOS, and hardened Linux configurations enforce strict W^X.

---

## 10. ASLR — Address Space Layout Randomization

**ASLR (Address Space Layout Randomization)** randomizes the base addresses of key regions:

Without ASLR:
```
Every process:
  Stack: 0x7fffffffffff000 (attacker knows!)
  Heap: 0x555555558000 (attacker knows!)
  libc: 0x7ffff7a00000 (attacker knows!)
```

With ASLR:
```
Process run 1:
  Stack: 0x7f3bcd124000
  Heap: 0x562a73cc0000
  libc: 0x7f12ab3d0000

Process run 2:
  Stack: 0x7e1fab239000  ← different each time
  Heap: 0x5a17bd210000
  libc: 0x7e89ac140000
```

Attackers often need to know the address of useful code (gadgets for ROP, shellcode) or data (function pointers). ASLR makes these addresses unpredictable.

**Entropy of ASLR:**
On 64-bit Linux: ~28 bits of stack randomization, ~17 bits of heap, ~28 bits of library placement. Brute-force is infeasible.

**Checking ASLR:**
```bash
cat /proc/sys/kernel/randomize_va_space
# 0 = disabled
# 1 = partial (stack, libraries)
# 2 = full (stack, heap, libraries, brk) ← default
```

**ASLR bypass techniques:**
- Information leaks: if attacker can read a pointer, they can deduce base address
- Partial overwrite: overwrite only lower bytes of address (offset within page)
- Heap spraying: fill heap with shellcode → statistics overcome entropy
→ These are why ASLR is necessary but not sufficient: ASLR + NX + stack canaries + RELRO together provide defense in depth.

---

## Summary

| Concept | Definition |
|---------|-----------|
| Segment | Named memory region with base, limit, type, DPL |
| GDT | Global Descriptor Table: array of segment descriptors |
| GDTR | CPU register pointing to GDT location |
| Selector | 16-bit value in segment register: index + table + RPL |
| DPL | Descriptor Privilege Level: 0 (kernel) to 3 (user) |
| CPL | Current Privilege Level: privilege of currently running code |
| TSS | Task State Segment: holds kernel stack pointer for Ring 3→0 transitions |
| NX bit | No-Execute: prevents executing code from data pages |
| W^X | Write XOR Execute: page is never both writable and executable |
| ASLR | Randomizes base addresses; attacker can't predict addresses |
| SMEP | Kernel can't execute user-mode pages |
| SMAP | Kernel can't access user-mode memory (without STAC) |
| KPTI | Separate page tables for user/kernel; mitigates Meltdown |

**The protection philosophy:**
Defense in depth. No single mechanism is sufficient. NX prevents code injection. ASLR makes addresses unpredictable. Stack canaries detect overflows. Segmentation enforces ring levels. Together, they make exploitation extremely difficult — though determined attackers continue to find creative bypasses.
