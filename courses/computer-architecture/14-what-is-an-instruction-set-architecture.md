# Chapter 14: What Is an Instruction Set Architecture?

> "The ISA is the most important design decision in computer architecture — it is a choice you may have to live with for fifty years."
> — David Patterson, co-designer of RISC-V

Imagine two chefs — one in Tokyo, one in Paris — who have never met. They can still follow the same recipe because they share a common language for measurements, techniques, and ingredients. The recipe is the contract. Neither chef needs to know the other's kitchen layout, their brand of oven, or how many sous-chefs they employ. The recipe abstracts all of that away.

That is precisely what an Instruction Set Architecture (ISA) does for hardware and software. It is the recipe — the shared contract — that lets a program written by a software engineer run on a chip designed by a hardware engineer, even when neither has ever communicated directly.

This chapter unpacks that contract in full: what it covers, why it persists for decades, how it differs from the actual chip implementation, and why the choice of ISA is one of the most economically and technically consequential decisions in all of computing.

---

## Table of Contents

1. [The Contract Between Hardware and Software](#1-the-contract-between-hardware-and-software)
2. [What an ISA Specifies](#2-what-an-isa-specifies)
3. [ISA vs Microarchitecture: Same Language, Different Dialects](#3-isa-vs-microarchitecture)
4. [The Application Binary Interface (ABI): ISA's Sibling](#4-the-application-binary-interface)
5. [Instruction Formats](#5-instruction-formats)
6. [Addressing Modes](#6-addressing-modes)
7. [Privilege Levels and Exception Handling](#7-privilege-levels-and-exception-handling)
8. [Why ISAs Last for Decades: The Lock-In Effect](#8-why-isas-last-for-decades)
9. [x86: Forty-Six Years of Dominance](#9-x86-forty-six-years-of-dominance)
10. [RISC-V: A Case Study in Clean Design](#10-risc-v-a-case-study-in-clean-design)
11. [The Major ISA Families](#11-the-major-isa-families)
12. [Summary](#12-summary)
13. [Exercises](#13-exercises)

---

## 1. The Contract Between Hardware and Software

### What the Contract Actually Says

When you compile a C program on Linux, the compiler transforms your source code into a sequence of binary numbers. Each binary number is an instruction — a single command to the processor. The processor reads each instruction, decodes what it means, and executes it.

For this to work, the compiler and the processor must agree on a common language:

- What does the bit pattern `0x00500293` mean?
- How many registers are there, and how wide are they?
- Which register holds the stack pointer?
- How do you call a function?
- When something goes wrong, what happens?

The ISA answers all of these questions. It is a written specification — typically hundreds of pages long — that defines the processor's behavior from the perspective of a programmer. Think of it as the processor's public API.

```
      SOFTWARE WORLD                           HARDWARE WORLD
  ─────────────────────────                ──────────────────────────
  Compiler, OS, Applications               CPU chip design, wafer fab

        |                                          |
        |         THE ISA CONTRACT                 |
        |   ┌────────────────────────┐             |
        └──>│  Instruction encoding  │<────────────┘
            │  Register definitions  │
            │  Memory model          │
            │  Addressing modes      │
            │  Privilege levels      │
            │  Exception behavior    │
            └────────────────────────┘

  Both sides honor the contract.
  Neither side needs to know the other's internals.
```

This separation is enormously powerful. A compiler engineer can generate correct x86 machine code without knowing whether the target chip uses 4 stages or 20 stages in its pipeline. A chip designer can redesign the entire internal execution engine without breaking a single existing application.

### The IBM System/360: The First Formal ISA

Before 1964, IBM sold computers one model at a time. Buying a faster machine meant rewriting all your software — because each model had its own unique instruction set. A small IBM 1401 ran programs incompatible with a large IBM 7094.

IBM's engineers, led by Gene Amdahl (of Amdahl's Law fame) and Fred Brooks (of The Mythical Man-Month), proposed something radical for the System/360: a single, documented ISA shared across all machines in the product line, from the cheapest to the most powerful.

A slow, inexpensive Model 20 and a fast, expensive Model 91 ran the same programs. The same binary worked on every machine in the family. Customers could upgrade hardware without rewriting software.

This idea — separating the ISA from its implementation — became the defining principle of all modern processor design. Every ISA you will encounter today descends conceptually from this 1964 decision.

### Quick Check

> 1. What is an ISA, and why is it described as a "contract"?
> 2. If an ISA is the "what," what is the "how"?
> 3. What problem did the IBM System/360's shared ISA solve for IBM's customers?

---

## 2. What an ISA Specifies

A complete ISA specification is comprehensive. Here are the major components.

### 2.1 The Instruction Set

The core of the ISA: the operations the processor can perform. Each instruction has:

- **Opcode**: a binary code identifying the operation (ADD, SUB, LOAD, STORE, BRANCH...)
- **Operands**: which registers or memory locations to read and write
- **Encoding**: exactly how the bits of the instruction are laid out

Different ISAs have vastly different numbers of instructions. x86 has thousands when you count all variants. RISC-V's base integer ISA has 47 instructions. More is not always better — more instructions mean more complex hardware and larger specification documents.

### 2.2 Registers

The ISA specifies:

- How many general-purpose registers exist (x86: 16, ARM64: 31, RISC-V: 32)
- How wide each register is (32-bit or 64-bit)
- Any special-purpose registers visible to the programmer: the Program Counter (PC), Stack Pointer (SP), flags/status registers
- Floating-point registers (often a separate bank)
- Whether registers have special conventions (e.g., in RISC-V, register x0 is hardwired to always contain zero — reads return 0, writes are ignored)

### 2.3 The Memory Model

The memory model specifies how the processor views and accesses memory:

- **Address space size**: how much memory can be addressed? 32-bit ISA: 4 GB maximum. 64-bit ISA: theoretically 16 exabytes (in practice, 128 TB or 256 TB depending on the implementation)
- **Byte ordering (endianness)**: when a 32-bit integer is stored in memory, which byte comes first? Big-endian (most significant byte first, as in SPARC) or little-endian (least significant byte first, as in x86)? RISC-V supports both but defaults to little-endian
- **Alignment requirements**: must a 4-byte integer be stored at an address divisible by 4? Many RISC ISAs require alignment; x86 allows unaligned access (at a performance penalty)
- **Memory consistency model**: if processor A writes to address X and then address Y, is processor B guaranteed to see those writes in the same order? Weak ordering (ARM, RISC-V) vs strong ordering (x86) has major implications for concurrent software

### 2.4 Data Types

What kinds of data does the processor natively support?

| Data Type | Size | Example Use |
|-----------|------|-------------|
| Byte | 8 bits | Characters, small integers |
| Halfword | 16 bits | Unicode code points, small values |
| Word | 32 bits | Standard integer, single-precision float |
| Doubleword | 64 bits | Large integers, double-precision float |
| IEEE 754 half-precision | 16 bits | Neural network weights |
| IEEE 754 single-precision | 32 bits | General floating-point computation |
| IEEE 754 double-precision | 64 bits | Scientific computing, databases |
| Vector (SIMD) | 128–512 bits | Media, ML, scientific workloads |

### 2.5 Addressing Modes

How does the processor compute the memory address of a load or store? We cover this in depth in Section 6.

### 2.6 Exception and Interrupt Model

What happens when something goes wrong — or when hardware needs the CPU's attention?

- **Exceptions** (synchronous): caused by the current instruction — division by zero, illegal instruction, memory access violation, page fault
- **Interrupts** (asynchronous): caused by external hardware — timer, keyboard, network card, disk

The ISA specifies: what causes each type, where the processor looks up the handler address (the interrupt vector table), which registers are saved automatically, and how the processor returns from the handler.

### 2.7 Privilege Levels

Not all code should be able to do everything. The OS kernel must be able to control hardware; user applications should not be able to corrupt the kernel. The ISA defines privilege levels and which instructions are only executable at higher privilege. (Covered fully in Section 7.)

### Quick Check

> 1. Name five distinct things that an ISA specification must define.
> 2. What does "endianness" mean, and why does the ISA need to specify it?
> 3. What is the difference between an exception and an interrupt?

---

## 3. ISA vs Microarchitecture

This distinction is one of the most important in all of computer architecture — and one of the most commonly confused.

### The Analogy: Language vs Dialect

Consider the English language. The rules of English grammar — its vocabulary, syntax, and meaning — are the "ISA." A person speaking in a London accent and a person speaking in a Texas accent are both speaking English. Their internal "implementation" (mouth shape, vocal patterns, pace) differs dramatically, but they can communicate because they share the same language.

The microarchitecture is the accent — the implementation. The ISA is the language — the specification.

```
                        x86-64 ISA
              ┌─────────────────────────────┐
              │  "ADD EAX, EBX means:        │
              │   EAX ← EAX + EBX"          │
              │  (defined in the spec)       │
              └─────────────────────────────┘
                           |
         ┌─────────────────┼─────────────────┐
         |                 |                 |
   Intel Core i3      Intel Core i9     AMD Ryzen 9
   (Raptor Lake)      (Raptor Lake)     (Zen 4)
   2-wide pipeline   6-wide pipeline   8-wide pipeline
   ~15W TDP          ~125W TDP         ~105W TDP
   ~$150             ~$550             ~$420
   
   ALL THREE:
   - Run the same Windows/Linux/macOS software
   - Execute the same binary programs
   - Implement the same x86-64 ISA
   
   NONE of them share their internal design.
```

### What Microarchitecture Adds

The microarchitecture is where all the performance magic happens. Given the same ISA, different microarchitectures can vary in:

| Property | Simple Microarch | High-Performance Microarch |
|----------|-----------------|---------------------------|
| Pipeline stages | 5 | 20+ |
| Issue width | 1 instruction/cycle | 6+ instructions/cycle |
| Execution order | In-order | Out-of-order |
| Cache levels | L1 only | L1 + L2 + L3 (32 MB+) |
| Branch predictor | Simple, 1-bit | Neural network-like, 99%+ accurate |
| Power use | 0.1 W (embedded) | 300 W (server) |

The ISA says "ADD must produce the correct sum." The microarchitecture decides whether to compute that ADD in 1 cycle or pipelined over 4 cycles, whether to predict ahead, whether to reorder other instructions around it.

### Why This Matters in Practice

Intel releases a new microarchitecture roughly every two years (Skylake → Kaby Lake → Coffee Lake → Rocket Lake → Alder Lake → Raptor Lake). Each generation improves performance by 5-20% through microarchitectural innovations. But every generation still runs the same x86-64 ISA — so all existing software runs without recompilation.

This is the ISA's greatest gift: hardware and software can evolve independently.

### Quick Check

> 1. What is the analogy between a spoken language and an ISA?
> 2. Can two chips with the same ISA but different microarchitectures have different power consumption? Give an example.
> 3. What enables Intel to release a new microarchitecture every two years without breaking existing software?

---

## 4. The Application Binary Interface

### ISA Is Not Enough

The ISA tells you what instructions exist and how they are encoded. But it does not tell you how two compiled programs should talk to each other. Consider:

- When a program calls `printf()`, which register holds the first argument? The second?
- When a function returns an integer, which register holds the return value?
- Where does the stack grow? From high addresses to low, or low to high?
- What is the format of an executable file that the OS should load?
- How does a user program ask the kernel to perform a file read?

None of this is in the ISA. This is the domain of the **Application Binary Interface (ABI)** — a second layer of specification that builds on the ISA.

### The ABI as the "Social Contract"

If the ISA is the grammar of the language, the ABI is the social norms — the conventions everyone follows so things actually work together. An ABI specifies:

```
  ABI Layer (builds on ISA)
  ─────────────────────────────────────────────────────
  1. Calling conventions:
     - Which registers pass function arguments?
     - Which register holds the return value?
     - Which registers must the callee preserve (callee-saved)?
     - Which registers can the callee clobber (caller-saved)?
  
  2. Stack layout:
     - Stack growth direction (almost always downward)
     - Stack alignment requirements (x86-64 Linux: 16-byte aligned)
     - Stack frame structure (saved frame pointer, local variables)
  
  3. System call interface:
     - How is a system call invoked? (SYSCALL instruction on x86-64)
     - Which register holds the system call number?
     - Where are arguments passed for system calls?
  
  4. Object and executable file format:
     - ELF (Executable and Linkable Format) on Linux
     - PE/COFF on Windows
     - Mach-O on macOS
  
  5. Data layout:
     - Sizes and alignments of all data types
     - Struct padding and packing rules
  ─────────────────────────────────────────────────────
```

### A Concrete Example: x86-64 System V ABI (Linux)

When you call a function like `sum(a, b, c, d, e, f)` on Linux x86-64:

- First 6 integer arguments: passed in registers RDI, RSI, RDX, RCX, R8, R9
- Return value: in register RAX (and RDX for 128-bit returns)
- Additional arguments: pushed onto the stack

This convention is fixed for all code on Linux x86-64. If a library is compiled with one compiler and your application with another, they can still call each other — because both follow the same ABI.

### ABI vs ISA: Different Problems, Different Scope

| Aspect | ISA | ABI |
|--------|-----|-----|
| Defined by | Chip designer (Intel, ARM Ltd, RISC-V Foundation) | OS vendor, standards body (e.g., psABI) |
| Covers | Hardware-visible state and behavior | Software-to-software conventions |
| Changes how often | Rarely (decades) | Occasionally (OS version boundaries) |
| Breaking it affects | All software on that hardware | Software across compilation units |

The same ISA can have multiple ABIs. Linux on x86-64 uses the System V ABI. Windows on x86-64 uses the Microsoft x64 ABI. Both run on the same chip — but a function compiled for Linux calling conventions cannot be called directly by code compiled for Windows conventions.

### Quick Check

> 1. What does ABI stand for, and what layer does it sit at relative to the ISA?
> 2. Name three things an ABI specifies that the ISA does not.
> 3. Can two programs compiled on the same ISA but following different ABIs interoperate directly? Why or why not?

---

## 5. Instruction Formats

Instructions are binary numbers. The ISA specifies exactly how the bits of each instruction encode its meaning — the **instruction format**.

### Fixed-Length vs Variable-Length Instructions

```
FIXED-LENGTH (RISC — RISC-V, ARM64, MIPS):
Each instruction is exactly 32 bits.

  ┌────────────────────────────────┐
  │          32 bits               │
  └────────────────────────────────┘
  ┌────────────────────────────────┐
  │          32 bits               │
  └────────────────────────────────┘
  ┌────────────────────────────────┐
  │          32 bits               │
  └────────────────────────────────┘

Pros: simple decode, always aligned, pipeline-friendly
Cons: wastes space for simple operations


VARIABLE-LENGTH (CISC — x86):
Instructions range from 1 byte to 15 bytes.

  ┌──┐
  │1B│         (e.g., NOP, RET)
  └──┘
  ┌────────┐
  │  3B    │   (e.g., ADD EAX, EBX)
  └────────┘
  ┌──────────────────────────────┐
  │            15B               │  (rare, prefix-heavy instruction)
  └──────────────────────────────┘

Pros: compact code, common ops are short
Cons: complex decode, must find boundaries between instructions
```

### RISC-V Instruction Formats

RISC-V has six instruction formats, all exactly 32 bits. The formats are cleverly designed so that the register fields (rs1, rs2, rd) are always in the same bit positions regardless of format — allowing the register file to be read before the full instruction is decoded.

```
R-type (register-register operations: ADD, SUB, AND, OR...):
31      25 24  20 19  15 14  12 11   7 6      0
┌─────────┬──────┬──────┬──────┬──────┬────────┐
│ funct7  │ rs2  │ rs1  │funct3│  rd  │ opcode │
└─────────┴──────┴──────┴──────┴──────┴────────┘
   7 bits   5 bits  5 bits 3 bits 5 bits  7 bits

I-type (immediate operations, loads: ADDI, LW, JALR...):
31              20 19  15 14  12 11   7 6      0
┌─────────────────┬──────┬──────┬──────┬────────┐
│   imm[11:0]     │ rs1  │funct3│  rd  │ opcode │
└─────────────────┴──────┴──────┴──────┴────────┘
      12 bits        5bits  3bits  5bits   7bits

S-type (stores: SW, SH, SB...):
31    25 24  20 19  15 14  12 11   7 6      0
┌───────┬──────┬──────┬──────┬──────┬────────┐
│imm[11:5]│ rs2 │ rs1  │funct3│imm[4:0]│opcode│
└───────┴──────┴──────┴──────┴──────┴────────┘

Note: the immediate is split across two fields to keep rs1, rs2, rd
in the same bit positions as in R-type. Hardware designers love this.
```

### What an Assembly Instruction Looks Like

Here is a small RISC-V sequence to add two numbers and store the result:

```
# RISC-V assembly:  c = a + b  (a in x10, b in x11, result goes to x12)
add  x12, x10, x11     # x12 = x10 + x11

# Load a value from memory (address in x5, result to x6):
lw   x6, 0(x5)         # x6 = Memory[x5 + 0]

# Store a value to memory (x7 → address x5 + 4):
sw   x7, 4(x5)         # Memory[x5 + 4] = x7

# Branch if equal:
beq  x10, x11, label   # if x10 == x11, jump to label
```

These assembly mnemonics are a human-readable form of the binary encoding. The assembler converts each mnemonic into its 32-bit binary representation following the format diagrams above.

### Quick Check

> 1. What is the main advantage of fixed-length instructions for pipeline hardware?
> 2. In RISC-V, why are the rs1, rs2, and rd fields always at the same bit positions?
> 3. What does an assembler do to a mnemonic like `add x12, x10, x11`?

---

## 6. Addressing Modes

An **addressing mode** specifies how the processor computes the memory address of a load or store. Think of addressing modes as different ways to describe where something is stored in a vast warehouse.

### Register (Direct)

The operand is a register — no memory access at all. Most fast because registers are on-chip.

```
ADD x3, x1, x2       ; x3 = x1 + x2  (all registers, no memory)
```

### Immediate

The operand is a constant embedded directly in the instruction.

```
ADDI x1, x1, 5       ; x1 = x1 + 5  (5 is inside the instruction bits)
```

### Base + Offset (Displacement)

The most common memory addressing mode. The address is computed as a base register value plus a constant offset.

```
LW x5, 12(x6)        ; x5 = Memory[x6 + 12]
```

This mode is ideal for accessing struct fields. If x6 points to a struct, each field is at a fixed offset from the struct's start. It also handles local variables (base = stack pointer, offset = variable's position on the stack).

```
   x6 → ┌──────────────┐  offset 0:  name (first field)
         │  name        │
         ├──────────────┤  offset 8:  age (second field)
         │  age         │
         ├──────────────┤  offset 12: score (third field)
         │  score       │  ← LW x5, 12(x6) fetches this
         └──────────────┘
```

### Register Indirect

Address is simply the value in a register — effectively base+offset with offset=0.

```
LW x5, 0(x6)         ; x5 = Memory[x6]
```

### Indexed (Base + Register)

Address is a base register plus an index register (optionally scaled). x86 supports this natively. RISC-V must compute the address with an explicit ADD first.

```
# x86: load from array[i] where array base in RBX, index in RCX, 4 bytes each
MOV EAX, [RBX + RCX*4]

# RISC-V equivalent:
SLLI x7, x9, 2       ; x7 = x9 * 4  (shift left by 2 = multiply by 4)
ADD  x7, x6, x7      ; x7 = x6 + x7  (compute effective address)
LW   x5, 0(x7)       ; x5 = Memory[x7]
```

### PC-Relative

Address is the current Program Counter value plus a signed offset. Essential for position-independent code — code that can be loaded at any address in memory because all jumps and data references are relative, not absolute.

```
# RISC-V branch (always PC-relative):
BEQ x1, x2, +100     ; if x1 == x2, jump to (current_PC + 100)
```

Shared libraries (.so files on Linux, .dll on Windows) use PC-relative addressing so they can be loaded at any memory address by the OS.

### Quick Check

> 1. Which addressing mode is best for accessing a field of a struct? Explain why.
> 2. Why is PC-relative addressing important for shared libraries?
> 3. How does RISC-V handle indexed array access without a native indexed addressing mode?

---

## 7. Privilege Levels and Exception Handling

### Why Privilege Separation Exists

Imagine a bank where every customer also has the keys to the vault. The ISA's privilege system prevents this. Without privilege levels, any program could:

- Read another program's memory (steal passwords, private keys)
- Modify the OS kernel (install malware with full system access)
- Directly control hardware (send malicious commands to the network card)

Privilege levels create a protected boundary between user code and kernel code.

### x86-64 Rings

x86 defines four privilege rings (0 through 3). Modern OSes use only two:

```
       Ring 0 (Kernel Mode)
   ┌───────────────────────────────┐
   │  OS kernel, device drivers    │
   │  Can: access any hardware     │
   │  Can: modify page tables      │
   │  Can: execute ALL instructions│
   └───────────────────────────────┘
              |   SYSCALL / SYSRET
   ┌───────────────────────────────┐
   │  Ring 3 (User Mode)           │
   │  Applications: Chrome, Word   │
   │  Cannot: access hardware      │
   │  Cannot: modify page tables   │
   │  Cannot: execute privileged   │
   │          instructions         │
   └───────────────────────────────┘
```

### RISC-V Privilege Modes

RISC-V defines three privilege modes, each useful for different system software layers:

| Mode | Name | For |
|------|------|-----|
| M-mode | Machine | Firmware, bootloader (highest trust) |
| S-mode | Supervisor | Operating system kernel |
| U-mode | User | Applications (lowest trust) |

An embedded RISC-V chip might only implement M-mode (no OS needed). A desktop RISC-V system implements all three.

### Exception Handling

When a program causes an exception — dividing by zero, accessing an unmapped memory page, executing a privileged instruction from user mode — the processor:

1. Saves the current Program Counter (where were we?)
2. Saves key register state
3. Switches to a higher privilege level
4. Looks up the exception handler address (from the trap/interrupt vector table)
5. Jumps to the handler

The handler (in the OS kernel) diagnoses the problem and either fixes it (a page fault can be resolved by loading the page from disk) or terminates the offending process.

The ISA must specify exactly what state is saved, where the handler address comes from, and what the handler must do to return cleanly. In RISC-V, the `mtvec` register holds the machine-mode trap handler address. In x86, the IDT (Interrupt Descriptor Table) holds handler addresses for all 256 possible interrupt/exception vectors.

### Quick Check

> 1. Why does user-mode code not need to access hardware directly?
> 2. In RISC-V, what is M-mode used for that S-mode is not?
> 3. What are the five steps a processor takes when an exception occurs?

---

## 8. Why ISAs Last for Decades

Of all the engineering decisions in computing, ISA choices have the most extraordinary staying power. The x86 ISA was defined in 1978 — software compiled that year can still run on chips manufactured today, 46 years later.

### The Software Investment Flywheel

Every year that an ISA is in use, more software is compiled for it. More software means more value in that ISA. More value means more incentive to keep the ISA alive. This is a flywheel that becomes nearly impossible to stop.

```
         More software compiled
               ↗         ↖
More value in ISA         More hardware sold
               ↖         ↗
         More hardware supports ISA
```

Consider what x86 has accumulated over 46 years:

- Every Windows application ever written (billions of programs)
- Decades of OS kernel code, device drivers, libraries
- Compilers, interpreters, virtual machines all targeting x86
- Engineers trained specifically in x86 assembly and architecture
- Textbooks, courses, certification programs about x86
- An enormous ecosystem of debugging and profiling tools

### Binary Compatibility: The Invisible Promise

Users never see the ISA — they just expect their software to work. When Intel released the Pentium 4 in 2000, customers expected their Windows 98 software from 1998 to run. It did. When Intel released the Core 2 in 2006, the same promise held. When AMD introduced x86-64 in 2003, they deliberately maintained full 32-bit x86 compatibility — breaking it was unthinkable.

The cost of migration becomes clear when you look at what transitions require:

| What Must Change | Cost |
|-----------------|------|
| Compiler retargeting | Months to years of engineering |
| OS kernel rewrite | Years of engineering |
| All applications recompiled | Rebuild entire software ecosystem |
| User training for new behaviors | Enormous, largely invisible |
| Testing all software on new ISA | Enormous, largely invisible |

### The Itanium Catastrophe: A Warning

In 1999, Intel and HP introduced the Itanium processor with IA-64 — a brand-new ISA based on VLIW (Very Long Instruction Word) architecture. It was technically radical and offered massive performance potential for the right workloads.

It failed almost completely.

The problem: IA-64 was not backward-compatible with x86. Every existing program had to be recompiled. Compilers for IA-64 were harder to write (the ISA required compilers to do much of the scheduling work). Users refused to migrate.

Meanwhile, AMD extended the existing x86 ISA to 64 bits with AMD64 (now called x86-64). It preserved complete backward compatibility with all 32-bit x86 software. Intel was forced to license and adopt AMD's extension. By 2010, Itanium was all but dead. By 2021, Intel had discontinued it entirely.

**The lesson**: technical superiority does not overcome the gravitational pull of an established ISA.

### Quick Check

> 1. Describe the "software investment flywheel" in your own words.
> 2. Why did the Itanium processor fail despite being technically advanced?
> 3. What does "binary compatibility" mean, and why do users care about it?

---

## 9. x86: Forty-Six Years of Dominance

### From 8086 to Infinity

In 1978, Intel introduced the 8086 processor with a 16-bit ISA: 4 general-purpose 16-bit registers (AX, BX, CX, DX), a 20-bit segmented address space (1 MB), and an instruction set designed to be an improvement over the 8080 that had powered the Altair 8800.

Over the next four decades, Intel grew x86 while maintaining backward compatibility at every step:

```
1978: Intel 8086          — 16-bit, 1 MB address space
1982: Intel 80286         — protected mode, 16 MB
1985: Intel 80386         — 32-bit (IA-32), 4 GB, flat memory model
1993: Intel Pentium       — superscalar, FPU on chip
1999: Intel Pentium III   — SSE (streaming SIMD), 128-bit vectors
2000: Intel Pentium 4     — hyperthreading, SSE2
2003: AMD Opteron         — AMD64 (x86-64): 64-bit registers, 48-bit virtual address
2007: Intel Core 2        — Intel adopts AMD64 (renames it EM64T, then Intel 64)
2011: Intel Sandy Bridge  — AVX: 256-bit SIMD vectors
2013: Intel Haswell       — AVX2, FMA
2017: Intel Skylake-X     — AVX-512: 512-bit SIMD vectors
2022: Intel Raptor Lake   — Hybrid cores, still fully x86-64 compatible
```

Every step added new instructions (never removed old ones). Software from 1985 that ran on the 80386 still runs on a 2024 Intel processor. This is not an accident — it is the explicit design goal of x86 compatibility.

### The Economic Machine

x86 dominance is self-reinforcing through several economic mechanisms:

**Developer tooling**: Microsoft Visual Studio, GCC, Clang, LLVM — all have had decades of x86 optimization. The x86 code generation paths are mature, well-tuned, and reliable.

**Operating systems**: Windows's dominant market position means the vast majority of enterprise software is x86. Linux runs on everything, but most commercial applications are packaged for x86.

**Enterprise lock-in**: large organizations run x86 software written 10-20 years ago. Migrating to a new ISA would require testing and validating hundreds of applications — a project that costs millions and takes years.

**The virtualization ecosystem**: VMware, Hyper-V, KVM — the entire enterprise virtualization industry is built around x86 hardware. These products cannot be easily ported to another ISA without years of engineering.

### x86's Quirks: The Price of History

x86's backward compatibility comes at a cost. The ISA carries 46 years of historical decisions that made sense in their era:

- **Segmented memory model**: the 16-bit 8086 used a segment:offset addressing scheme to access more than 64 KB. Modern x86 uses flat addressing, but the segment registers (CS, DS, ES, FS, GS, SS) still exist and must be handled by hardware
- **Variable-length instructions**: from 1 to 15 bytes long, making the decode stage extremely complex. Modern x86 processors translate x86 instructions into an internal RISC-like micro-operation (uop) format before execution
- **Flag register**: a 64-bit register (RFLAGS) with individual bits encoding the state of the last arithmetic operation. Condition codes on flags is a complex dependency tracking problem for out-of-order execution
- **Limited registers**: x86-64 has 16 general-purpose registers; ARM64 has 31; RISC-V has 32. More registers reduce the need for spilling values to memory and improve performance

Despite all of this, x86 hardware is so sophisticated that it often outperforms theoretically cleaner ISAs in practice — because the hardware budget that would go into a simpler ISA gets spent on larger caches, better branch predictors, and more execution units.

### Quick Check

> 1. In which year was the x86 ISA defined, and how has it changed since then?
> 2. Why does x86 use variable-length instructions, and what problem does this cause?
> 3. What economic mechanisms keep x86 dominant beyond just technical merit?

---

## 10. RISC-V: A Case Study in Clean Design

### The Academic Origin

In 2010, Krste Asanovic and David Patterson at UC Berkeley needed an ISA for research. Every existing ISA had problems:

- x86: too complex, too many historical warts, patent and licensing issues
- ARM: proprietary, licensing fees, complex extension ecosystem
- MIPS: aging, proprietary, no longer actively developed
- SPARC: declining, proprietary

So they designed a new one from scratch: RISC-V (pronounced "RISC five" — the fifth RISC design from Berkeley). They made it open-source — free to implement without royalties, free to extend, free to use in any product.

### The Design Philosophy

RISC-V was designed knowing everything that had been learned about ISA design since 1964:

**Base + Extensions**: a small, clean base integer ISA (47 instructions) plus well-defined optional extensions:

```
  RISC-V Extension Naming:
  
  RV64GC  (a typical configuration)
  │  │  └── C: Compressed (16-bit versions of common instructions)
  │  └───── G: General (= IMAFD bundled)
  │
  └──────── 64-bit base integer ISA
  
  Individual extensions:
  I  — Base integer instructions (47 instructions)
  M  — Integer multiplication and division
  A  — Atomic memory operations
  F  — Single-precision floating-point
  D  — Double-precision floating-point
  C  — Compressed 16-bit instructions
  V  — Vector (SIMD) operations
  H  — Hypervisor support
```

This modularity means a tiny embedded RISC-V processor (RV32I, 32 registers, integer only) and a massive server RISC-V processor (RV64GCVH) share the same base ISA and can run the same base integer code.

**Clean register conventions**:

```
  RISC-V Register Conventions (ABI names):
  
  x0  (zero)  — always zero; writes ignored
  x1  (ra)    — return address
  x2  (sp)    — stack pointer
  x3  (gp)    — global pointer
  x4  (tp)    — thread pointer
  x5-x7   (t0-t2)   — temporaries (caller-saved)
  x8  (s0/fp) — saved register / frame pointer (callee-saved)
  x9  (s1)    — saved register (callee-saved)
  x10-x11 (a0-a1)   — function arguments / return values
  x12-x17 (a2-a7)   — function arguments
  x18-x27 (s2-s11)  — saved registers (callee-saved)
  x28-x31 (t3-t6)   — temporaries (caller-saved)
```

**No status flags**: Instead of a FLAGS register that every arithmetic instruction updates, RISC-V branches compare two registers directly (`BEQ x1, x2, label` — branch if equal). This eliminates a major source of data dependencies and simplifies out-of-order execution.

**Hardwired zero register**: x0 always reads as 0. This elegantly encodes several useful pseudo-instructions:

```
  MV x5, x6     →   ADD  x5, x0, x6    (copy: add zero)
  LI x5, 42     →   ADDI x5, x0, 42    (load immediate)
  NOP           →   ADDI x0, x0, 0     (write to x0, does nothing)
  J  label      →   JAL  x0, label     (jump, discard return address)
  RET           →   JALR x0, x1, 0     (jump to return address, discard)
```

These are not special instructions — they are just regular instructions that happen to use x0 cleverly.

### RISC-V's Growing Ecosystem

As of 2024, RISC-V processors appear in:

- **Embedded**: SiFive's microcontrollers, Western Digital's hard drive controllers, ESP32-C3 WiFi chips
- **High-performance**: SiFive P670, China's Xuantie C910, Alibaba's Yitian 710 (partially RISC-V inspired)
- **Research**: every computer architecture research paper that needs a real ISA uses RISC-V
- **Space**: NASA and ESA are exploring RISC-V for radiation-hardened processors
- **India's national processor**: the SHAKTI program is building indigenous RISC-V processors for national infrastructure independence

RISC-V's political advantage: it belongs to no single company. A country or company that wants to be independent of US export controls (which can restrict access to x86 and ARM licenses) can implement RISC-V without licensing from any American corporation.

### Quick Check

> 1. Why was RISC-V designed at UC Berkeley rather than being adopted from an existing ISA?
> 2. What does the "extension" system in RISC-V allow that monolithic ISAs like x86 cannot?
> 3. How does the hardwired x0 register enable several useful pseudo-instructions?

---

## 11. The Major ISA Families

Here is a comparison of the most important ISAs in computing history:

| ISA | Year | Bits | Registers | Instruction Width | License | Philosophy |
|-----|------|------|-----------|------------------|---------|-----------|
| x86-64 | 1978 / 2003 | 64 | 16 GPR | 1–15 bytes (variable) | Proprietary (Intel/AMD) | CISC — maximize compatibility |
| ARM64 (AArch64) | 1985 / 2011 | 64 | 31 GPR | Fixed 32-bit | Proprietary (Arm Ltd) | RISC — power efficiency |
| RISC-V (RV64G) | 2010 | 64 | 32 GPR | Fixed 32-bit (or 16-bit compressed) | Open source | RISC — clean, modular, open |
| MIPS64 | 1981 | 64 | 32 GPR | Fixed 32-bit | Proprietary (Wave Computing) | RISC — academic clean design |
| POWER10 | 1990 | 64 | 32 GPR | Fixed 32-bit | Proprietary (IBM) | RISC — high-performance servers |
| SPARC V9 | 1987 | 64 | 32+128 GPR | Fixed 32-bit | Open (Oracle) | RISC — register windows |
| IA-64 (Itanium) | 1999 | 64 | 128 GPR | 128-bit bundles | Discontinued | VLIW — compiler does scheduling |

### Where Each ISA Dominates Today

```
  ISA              Primary Market (2024)
  ──────────────   ──────────────────────────────────────────
  x86-64           PC desktops, laptops, servers, cloud
  ARM64            Mobile (iOS, Android), Apple M-series,
                   ARM-based cloud (AWS Graviton, Ampere Altra)
  RISC-V           Embedded, IoT, research, emerging markets
  POWER            IBM enterprise servers (financial, government)
  MIPS             Legacy routers, network equipment (declining)
  SPARC            Legacy Sun/Oracle servers (effectively retired)
  IA-64            Discontinued (last Itanium shipped 2021)
```

### The Philosophical Divide

The major ISA families represent two philosophical schools:

**CISC (Complex Instruction Set Computer)**: one instruction should do as much work as possible. Minimize the number of instructions a program needs. x86 is the archetypal CISC ISA — `REP MOVS` copies an entire block of memory in one instruction.

**RISC (Reduced Instruction Set Computer)**: keep instructions simple and uniform. Each instruction should complete in one clock cycle (in a simple implementation). The compiler, not the hardware, handles complexity. RISC-V, ARM, MIPS, SPARC, and POWER are RISC ISAs.

In practice, the distinction has blurred. Modern x86 processors internally translate CISC instructions into RISC-like micro-operations before execution. And modern RISC processors have gained many complex instructions (SIMD, cryptography, atomic operations) that are not simple at all.

### Quick Check

> 1. Which ISA dominates mobile computing? Which dominates desktop and server?
> 2. What is the philosophical difference between CISC and RISC?
> 3. Why is the CISC vs RISC distinction less meaningful in modern high-performance processors?

---

## 12. Summary

- The **Instruction Set Architecture (ISA)** is the contract between hardware and software — a complete specification of what the processor does, visible to the programmer.

- An ISA specifies: the instruction set (opcodes and encodings), registers (count, width, special roles), the memory model (address space, endianness, consistency), data types (integer widths, floating-point formats, SIMD), addressing modes, exception and interrupt handling, and privilege levels.

- **ISA vs Microarchitecture**: the ISA is the "what" — the abstract specification. The microarchitecture is the "how" — the concrete implementation. Multiple microarchitectures can implement the same ISA with vastly different performance and power characteristics. The Intel Core i3 and i9 both implement x86-64 but have very different performance.

- The **ABI (Application Binary Interface)** builds on the ISA to define calling conventions, stack layout, system call interfaces, and executable file formats. It is the contract that allows separately compiled programs to interoperate.

- **Fixed-length instructions** (RISC: 32-bit per instruction) simplify hardware decode and enable clean pipeline stages. **Variable-length instructions** (CISC: 1–15 bytes in x86) are more compact but require complex decode hardware.

- **Addressing modes** specify how to compute memory addresses: register, immediate, base+offset (dominant for struct access), indexed (for array access), and PC-relative (for position-independent code).

- **Privilege levels** separate kernel code (full hardware access) from user code (restricted). An illegal privilege escalation triggers an exception, which the OS handles.

- ISAs **last for decades** because of the software investment flywheel: more software increases the ISA's value, which increases hardware sales, which funds further development. Breaking compatibility means abandoning this ecosystem. The Itanium failure illustrates the cost of ignoring this inertia.

- **x86's 46-year dominance** is driven not just by technical merit but by backward compatibility, developer tooling, OS lock-in, and enterprise dependencies. x86 carries decades of complexity as the price of this compatibility.

- **RISC-V** is a clean, open-source, modular RISC ISA designed in 2010 with knowledge of every lesson from prior ISAs. Its base+extension model, open licensing, and growing ecosystem make it the ISA most likely to challenge ARM and x86 in new markets.

---

## 13. Exercises

### Easy

1. List five things that an ISA specification must cover. For each one, explain in one sentence why it must be in the ISA rather than left to the microarchitecture.

2. An Intel Core i5 and an Intel Core i9 both run Windows 11. What do they share that makes this possible? What do they not share?

3. Explain the base+offset addressing mode. Draw a diagram showing how it is used to access the second field of a struct where the first field is 8 bytes wide.

### Medium

4. x86 allows `ADD [0x1000], EAX` — add EAX to the 32-bit value at memory address 0x1000 and write the result back. RISC-V requires three separate instructions: `LW t0, 0(t1)`, then `ADD t0, t0, a0`, then `SW t0, 0(t1)`. Compare the two approaches: which has smaller code size? Which is easier to pipeline? Why might hardware designers prefer the RISC-V approach despite the larger program?

5. The RISC-V base integer ISA (RV32I) has exactly 47 instructions. x86-64 has several thousand. Why do more instructions not necessarily make an ISA better? What are the costs of a larger instruction set for compiler writers and hardware designers?

6. A program compiled for Linux x86-64 passes function arguments in registers RDI, RSI, RDX, RCX, R8, R9 (System V ABI). A program compiled for Windows x86-64 passes arguments in RCX, RDX, R8, R9 (Microsoft ABI). Both programs run on the same physical chip. What problem arises if a Linux-compiled library is called directly from a Windows program? How do real systems solve this?

### Hard

7. **ISA design exercise**: Suppose you are designing a new ISA for a deeply embedded processor with very limited area (small chip, low cost). You can only afford 4 registers. Design the instruction format: how many bits for the opcode? How many bits per register specifier? How do you handle constants (immediates)? What tradeoffs does a 4-register ISA force on compiler writers? Research the concept of a "stack machine" ISA as an alternative approach.

8. **The lock-in problem**: When Apple switched from Intel x86 to ARM-based Apple Silicon in 2020, they included Rosetta 2, a binary translation system that runs x86-64 code on ARM at roughly 70% of native ARM performance. (a) Explain at a high level how binary translation works: what must the translator discover about the x86 code to generate correct ARM code? (b) Why does binary translation incur a performance penalty? (c) What property of the x86 ISA makes binary translation harder than if the source ISA were RISC-V? (d) Why did Apple succeed where Intel/HP failed with Itanium? What was different about their situation?

9. **RISC-V extension design**: RISC-V allows custom extensions (custom opcode ranges reserved for non-standard instructions). A company building a processor for neural network inference wants to add a "matrix multiply accumulate" instruction: `MMA fd, fs1, fs2, fs3` — multiply the 4x4 matrix in floating-point registers fs1...fs4 by the matrix in fs5...fs8 and accumulate into fd...fd+15. What ISA-level questions must this extension answer? (Hint: consider encoding, register conventions, exception behavior on invalid inputs, interaction with context switches, and how the OS would save/restore the new architectural state.) Why do these questions make custom extensions harder to design correctly than they first appear?
