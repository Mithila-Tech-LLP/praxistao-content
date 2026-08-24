# Chapter 09: Registers — The CPU's Working Memory

> *"A chef with a huge pantry but no cutting board is useless. The work happens on the board — the pantry just holds what you are not using right now."*

Registers are the fastest storage in a computer — tiny, permanent memory cells built directly into the processor chip, just a handful of wire lengths from the ALU. Every addition you ask the CPU to perform, every comparison that decides which branch your program takes, every byte that moves across the screen — all of it flows through registers first. Yet a modern CPU has fewer than 50 registers you can name. Understanding why that number is so small, and how engineers squeeze every last drop of performance from so few slots, is one of the most illuminating puzzles in computer architecture.

---

## Table of Contents

1. [Why Registers Exist — The Speed Argument](#1-why-registers-exist)
2. [The Chef's Cutting Board: An Analogy](#2-the-chefs-cutting-board-an-analogy)
3. [How Many Registers? Counting Across Architectures](#3-how-many-registers)
4. [Register Width: 8, 16, 32, and 64 Bits](#4-register-width)
5. [General-Purpose Registers](#5-general-purpose-registers)
6. [Special-Purpose Registers](#6-special-purpose-registers)
7. [The Zero Register: Hardwired to Nothing](#7-the-zero-register-hardwired-to-nothing)
8. [Register Numbers in Machine Code](#8-register-numbers-in-machine-code)
9. [The Register File: Multi-Ported SRAM](#9-the-register-file-multi-ported-sram)
10. [Register Spilling: When You Run Out of Space](#10-register-spilling-when-you-run-out-of-space)
11. [A Comparison Table: x86-64, ARM64, and RISC-V](#11-a-comparison-table-x86-64-arm64-and-risc-v)
12. [Register Renaming: A Preview](#12-register-renaming-a-preview)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. Why Registers Exist

Imagine you are a calculator. Someone hands you two numbers, you add them, and hand back the result. Simple. Now imagine that between every keystroke, you have to walk to a warehouse two blocks away, find the number on a shelf, carry it back, do one operation, then walk the number back to the warehouse. That is what computing without registers would look like.

The speed hierarchy of memory in a modern computer is brutal and wide:

```
Storage Type           Access Time    Capacity (typical)    Relative Speed
────────────────────────────────────────────────────────────────────────────
Registers              ~0.3 ns        32 × 64 bits = 256B   1× (baseline)
L1 Cache               1–3 ns         32–64 KB               5–10× slower
L2 Cache               5–10 ns        256 KB – 1 MB          17–33× slower
L3 Cache               20–40 ns       8–64 MB                67–133× slower
Main Memory (RAM)      50–100 ns      8–32 GB                167–333× slower
NVMe SSD               50–100 µs      512 GB – 2 TB          ~200,000× slower
Hard Disk Drive        5–10 ms        1–4 TB                 ~20,000,000× slower
```

Registers are **300 times faster than RAM**. They operate at the same speed as the CPU's clock — one read or write per clock cycle, no waiting. RAM, by contrast, takes 50 to 100 nanoseconds to respond: on a 3 GHz CPU that is 150 to 300 clock cycles of doing nothing, just waiting for data to arrive.

The solution is to keep the values you are currently working with in registers. Load a value from RAM once, put it in a register, use it ten times in a row — all ten uses happen at register speed. This single optimization is responsible for an enormous fraction of modern CPU performance.

### The Numbers in Instructions

Every instruction the CPU executes must tell the hardware where to find its operands. In a register-based CPU, operands are named by register numbers:

```
RISC-V instruction:   ADD  x5, x1, x2
                            ↑   ↑   ↑
                         dest src1 src2

Meaning: x5 = x1 + x2
```

The values in x1 and x2 are already inside the chip, a few transistors away from the adder. The result goes back into x5, also inside the chip. No memory bus is involved. No waiting.

### Quick Check 1.1

> 1. How many times faster are registers compared to RAM?
> 2. If a CPU runs at 3 GHz (one cycle = 0.33 ns), and RAM takes 100 ns to respond, how many CPU cycles are wasted waiting for RAM?
> 3. What does "keeping a value in a register" mean physically?

---

## 2. The Chef's Cutting Board: An Analogy

Picture a professional chef cooking a complex dish. The kitchen has:

- A **cutting board** right in front of the chef — tiny, immediately accessible. This is where the onion being chopped sits right now.
- A **countertop** within arm's reach, holding the mise en place — small bowls of prepped ingredients. This is the cache.
- A **pantry** across the kitchen, stocked with everything the recipe might need. This is RAM.
- A **warehouse** down the street that receives deliveries. This is the SSD.

The chef does not run to the pantry for each knife stroke. The current onion sits on the cutting board. The next three ingredients are staged on the countertop. The pantry is consulted only when the countertop needs restocking.

Registers are the cutting board. The chef (CPU) can only actively work on what is on the board. Everything else must be fetched — and fetching takes time. The art of CPU design is ensuring the cutting board always holds exactly the right ingredients, and the board is never so large that finding something on it takes time of its own.

---

## 3. How Many Registers?

Every architecture makes a different choice about how many registers to expose to programs. This is not arbitrary — it reflects deep tradeoffs.

### The Count Across Architectures

| Architecture       | Integer Registers | Width     | Notes                                        |
|--------------------|-------------------|-----------|----------------------------------------------|
| Intel 8086 (1978)  | 8                 | 16-bit    | AX, BX, CX, DX, SI, DI, SP, BP              |
| x86-32 (IA-32)     | 8                 | 32-bit    | EAX–EBP; same 8 registers, wider             |
| x86-64 (AMD64)     | 16                | 64-bit    | RAX–R15; added R8–R15 in 64-bit mode         |
| ARM32 (ARMv7)      | 16                | 32-bit    | r0–r15; r13=SP, r14=LR, r15=PC              |
| ARM64 (AArch64)    | 31 + XZR          | 64-bit    | x0–x30 + zero register; SP separate          |
| MIPS               | 32                | 32 or 64  | $0 always zero; $31 = return address         |
| RISC-V (RV32I/64I) | 32                | 32 or 64  | x0–x31; x0 hardwired zero                    |
| SPARC              | 32 (per window)   | 32 or 64  | Register windows give 128+ physical registers |
| PowerPC / POWER    | 32                | 32 or 64  | r0–r31                                       |

### Why Not Thousands of Registers?

If registers are so fast, why not have 1,024 of them? Three reasons:

**1. Instruction encoding space.** Every instruction must encode which registers it uses. With 32 registers, you need 5 bits per register field. With 1,024 registers, you need 10 bits per field — and a typical instruction uses three register fields. That is 30 bits just for register names, leaving little room for the opcode and immediate values in a 32-bit instruction.

**2. Silicon area and power.** The register file — the hardware block containing all registers — must support multiple simultaneous reads and writes. Every additional register and every additional port makes the file larger, slower, and hungrier for power. Register files scale poorly: doubling the number of ports can quadruple the area.

**3. Context switch cost.** When the operating system switches from one program to another, it must save all registers to memory and load the new program's registers. More registers means more saving and restoring on every context switch — a constant overhead.

The sweet spot, empirically found by decades of architecture research, is 16–32 integer registers. RISC architectures (RISC-V, ARM, MIPS, PowerPC) converge on 32. x86 settled at 16 for historical reasons (it grew from 8 in the 32-bit era).

### Quick Check 3.1

> 1. x86-64 has 16 integer registers. How many bits are needed to address one of 16 registers?
> 2. RISC-V has 32 registers. How many bits to address one of 32?
> 3. Name two reasons why having 1,024 registers would cause problems.

---

## 4. Register Width

A register's **width** is the number of bits it can hold. Width has evolved with CPU generations:

```
Era           Width    Max integer value held
──────────────────────────────────────────────────────────────
Early CPUs    8-bit    255 (0xFF)
16-bit era    16-bit   65,535 (0xFFFF)
32-bit era    32-bit   4,294,967,295 (0xFFFFFFFF) ≈ 4 billion
64-bit era    64-bit   18,446,744,073,709,551,615 ≈ 18 quintillion
```

Width matters for two things: the largest integer value you can compute without overflow, and the maximum memory address the CPU can form (since a memory address is just a very large number).

### The 4 GB Wall

A 32-bit pointer can address at most 2^32 = 4,294,967,296 bytes = 4 GB of memory. By the mid-2000s, programs were outgrowing 4 GB. The move to 64-bit CPUs (and 64-bit registers) lifted the ceiling to 2^64 bytes — a number so large it is effectively infinite for current purposes (about 16 exabytes).

### Sub-Register Accessors

x86-64 allows you to access portions of a register by different names:

```
RAX (64 bits): [63..................32|31.........16|15..8|7..0]
                                      EAX (32-bit)
                                                     AX (16-bit)
                                                     AH   AL (8-bit)
```

Writing to EAX **zeroes the upper 32 bits** of RAX (by design, to avoid partial-register stalls). Writing to AX does NOT change bits 32–63. This asymmetry is a famous source of subtle bugs in x86 assembly.

ARM64 does the same more cleanly: writing to w0 (the 32-bit view of x0) zeroes the upper 32 bits of x0, consistently.

### Quick Check 4.1

> 1. Why did the industry move from 32-bit to 64-bit registers?
> 2. In x86-64, if you write the value 5 to AL, does RAX change completely, or just its lowest 8 bits?
> 3. A 16-bit register holds a maximum unsigned value of 65,535. What is the maximum unsigned value in a 32-bit register?

---

## 5. General-Purpose Registers

**General-purpose registers (GPRs)** are the workhorses — you can put any integer value in them and use them for any computation. Most programs spend the vast majority of their time reading and writing GPRs.

Despite the "general-purpose" name, every architecture assigns **soft conventions** to each register: informal agreements between compilers, operating systems, and libraries about what each register is used for. These conventions are called the **ABI** (Application Binary Interface) or **calling convention**.

### Why Have Conventions?

When function A calls function B, registers are shared physical hardware. Both functions want to use RAX for their own computations. Without an agreement, A and B would silently overwrite each other's values.

Calling conventions solve this by splitting registers into two categories:

**Caller-saved (volatile):** The calling function cannot trust these registers to survive a function call. If function A uses RAX and then calls function B, A must save RAX to the stack before the call and restore it after. Function B is free to clobber RAX.

**Callee-saved (non-volatile):** The called function must preserve these registers. If function B wants to use RBX, it must push RBX to the stack at the start and pop it at the end. Function A can trust RBX to be unchanged after any function call.

```
RISC-V ABI Register Summary
─────────────────────────────────────────────────────────────────
Register   ABI Name   Role                          Saved By?
─────────────────────────────────────────────────────────────────
x0         zero       Hardwired 0 (reads 0, ignores writes)
x1         ra         Return address                Caller
x2         sp         Stack pointer                 Callee
x3         gp         Global pointer                —
x4         tp         Thread pointer                —
x5–x7      t0–t2      Temporaries                   Caller
x8         s0 / fp    Saved reg / Frame pointer     Callee
x9         s1         Saved register                Callee
x10–x11    a0–a1      Function arguments / return   Caller
x12–x17    a2–a7      Function arguments            Caller
x18–x27    s2–s11     Saved registers               Callee
x28–x31    t3–t6      Temporaries                   Caller
─────────────────────────────────────────────────────────────────
```

### Example: A Function Call in Action

```assembly
# RISC-V: function add_three(a, b, c) returns a + b + c
# Arguments arrive in a0, a1, a2 (x10, x11, x12)
# Return value goes in a0 (x10)

add_three:
    add  a0, a0, a1     # a0 = a + b
    add  a0, a0, a2     # a0 = (a + b) + c
    ret                 # return (return address is in ra = x1)
```

No saving needed — this function only touches caller-saved registers (a0, a1, a2) and returns in a0. Clean and fast.

### Quick Check 5.1

> 1. What is the difference between a "caller-saved" and a "callee-saved" register?
> 2. In RISC-V, the register a0 holds the return value of a function. Which registers hold the arguments?
> 3. If a function uses register s2 (x18, a callee-saved register), what must it do before returning?

---

## 6. Special-Purpose Registers

Beyond the general-purpose bank, every CPU has registers dedicated to specific roles. These are not just conventions — the hardware is wired to read or write these registers automatically during specific operations.

### Program Counter (PC)

The **Program Counter** (also called the Instruction Pointer, IP, in Intel's world) holds the address of the next instruction to fetch. It is updated automatically after every instruction:

```
Normal flow:         PC = PC + 4  (for 32-bit instructions like RISC-V)
Branch taken:        PC = branch_target_address
Function call:       PC = function_address
Exception/Interrupt: PC = exception_handler_address
```

On most architectures, you cannot name the PC in a general-purpose instruction like `ADD`. You change it only through control-flow instructions (branches, jumps, calls). This separation is intentional — it makes the CPU's fetch-decode-execute loop cleaner and safer.

```
                      ┌─────────────────────────────────┐
                      │           CPU Datapath           │
                      │                                  │
  ┌────────┐  fetch   │  ┌─────────────────────────────┐│
  │  RAM   │◄─────────┼──│  PC (Program Counter)        ││
  │        │  instr   │  │  next_address = 0x1000_04B8  ││
  └────────┘ ─────────►  └─────────────────────────────┘│
                      │              │ +4                 │
                      │              ▼                    │
                      │  ┌─────────────────────────────┐│
                      │  │  Instruction Register (IR)   ││
                      │  │  holds current instruction   ││
                      │  └─────────────────────────────┘│
                      └─────────────────────────────────┘
```

### Stack Pointer (SP)

The **Stack Pointer** tracks the top of the call stack — a region of memory that grows downward (toward lower addresses) as functions are called. Each function call pushes a **stack frame** containing:

- The return address
- Saved registers the function will clobber
- Local variables that could not fit in registers

```
High addresses
┌──────────────────┐
│   main()'s frame │  ← SP was here before main called foo()
├──────────────────┤
│   foo()'s frame  │  ← SP was here before foo called bar()
├──────────────────┤
│   bar()'s frame  │  ← SP points here now (top of stack)
└──────────────────┘
Low addresses
```

When bar() returns, SP moves back up, and its frame is effectively gone (the data is still there in memory, but SP now points above it, so the next function call will overwrite it).

### Frame Pointer / Base Pointer (FP / BP)

The **Frame Pointer** (called rbp in x86, x29/fp in ARM64, s0/fp in RISC-V) points to a fixed location within the current stack frame — usually just above the local variables. While SP fluctuates as things are pushed and popped, FP stays fixed for the duration of a function, making it easy to access local variables at predictable offsets:

```assembly
# x86-64: reading local variable at offset -8 from frame pointer
mov rax, [rbp - 8]     # first local variable
mov rbx, [rbp - 16]    # second local variable
mov rcx, [rbp + 8]     # first argument (above saved rbp)
```

Modern compilers often omit FP (a technique called *frame pointer omission*, enabled by -fomit-frame-pointer in GCC) to free up one extra register. Debuggers use FP to walk the call stack; without it, stack unwinding is harder but not impossible.

### Link Register / Return Address Register (LR / RA)

In ARM and RISC-V, when you call a function, the return address is saved in a dedicated register (x30/lr in ARM64, x1/ra in RISC-V) rather than automatically pushed to the stack. This makes **leaf functions** — functions that do not call other functions — extremely fast:

```assembly
# RISC-V leaf function (no function calls inside)
square:
    mul  a0, a0, a0     # a0 = a0 × a0
    ret                 # PC = ra (return to caller)
    # Never touched the stack!
```

If a function does call another function, it must first save ra to the stack (because the inner call will overwrite ra with a new return address).

### Instruction Register (IR)

After the CPU fetches an instruction from memory, it stores it in the **Instruction Register** while the decoder figures out what opcode it is and what registers it touches. This is an internal register — programs cannot read or write it directly.

### Status / Flags Register

The **flags register** (called RFLAGS in x86, CPSR/NZCV in ARM, implicitly set by compare instructions in RISC-V) stores single-bit outcomes from ALU operations. The major flags:

```
Flag   Meaning
────────────────────────────────────────────────────────
Z      Zero flag: result was zero
N      Negative flag: result was negative (sign bit = 1)
C      Carry flag: unsigned overflow (carry out of MSB)
V / O  Overflow flag: signed overflow
```

These flags drive conditional branches:

```assembly
# x86-64: "if (a == b) jump to label"
cmp  rax, rbx       # subtract rax - rbx, set flags, discard result
je   equal_label    # jump if Zero flag is set
```

RISC-V takes a different approach: no flags register. Instead, comparison and branch are combined in one instruction:

```assembly
# RISC-V: "if (x1 == x2) jump to label"
beq  x1, x2, equal_label    # branch if equal
```

This is cleaner — no hidden state flowing between instructions — but requires dedicated branch instructions for each comparison type (beq, bne, blt, bge, bltu, bgeu).

### Control and Status Registers (CSRs) in RISC-V

RISC-V defines a separate address space of **Control and Status Registers** for privileged operations:

```
CSR Name    Number   Purpose
────────────────────────────────────────────────────────────
mstatus     0x300    Machine status (global interrupt enable, privilege mode)
mie         0x304    Machine interrupt enable (per-interrupt enable bits)
mtvec       0x305    Machine trap vector (address of interrupt handler)
mepc        0x341    Machine exception program counter (where to return after trap)
mcause      0x342    Machine cause register (why did the trap happen?)
cycle       0xC00    Cycle counter (how many cycles have elapsed)
instret     0xC02    Instructions retired counter
```

These cannot be accessed with normal load/store instructions. You need special `csrr` (CSR read) and `csrw` (CSR write) instructions, and accessing most of them from user mode causes a privilege trap.

### Quick Check 6.1

> 1. What does the Program Counter hold, and when does it change?
> 2. Why does the stack pointer decrease (move to lower addresses) when a function is called?
> 3. What is the difference between caller-saved and callee-saved registers, and how does the Link Register (RA) relate to this?

---

## 7. The Zero Register: Hardwired to Nothing

RISC-V and MIPS both include a register hardwired to the value zero: **x0** in RISC-V, **$zero** in MIPS, **xzr/wzr** in ARM64. This sounds like a waste of a register slot — why permanently dedicate one of your precious 32 registers to a constant?

The answer: **zero appears everywhere in code**, and having it as a register simplifies the instruction set enormously.

### Uses of Zero

```assembly
# Initialize a register to zero
add  x5,  x0, x0       # x5 = 0 + 0 = 0

# Copy one register to another (no dedicated MOV instruction needed)
add  x5,  x6, x0       # x5 = x6 + 0 = x6

# Unconditional branch (always-true comparison)
beq  x0,  x0, label    # x0 == x0 always, so always branch

# Discard a result (write to x0 — hardware ignores it)
add  x0,  x1, x2       # compute x1+x2, throw away result (side effects only)

# Test if a value is zero
beq  x5,  x0, is_zero  # if x5 == 0, branch to is_zero

# Negate: compute 0 - x5 (two's complement negate)
sub  x5,  x0, x5       # x5 = 0 - x5 = -x5
```

Each of these would require a separate instruction format or a special opcode without the zero register. With it, a single ADD and BEQ instruction handle all these cases. The ISA becomes **orthogonal** — a small set of instructions that compose cleanly.

### Hardware Implementation

The zero register is trivially implemented: the read path for x0 is a wire tied permanently to zero. Any write to x0 is silently ignored — the decoder simply does not enable the write port for that register.

```
Register File Write Logic:

write_enable[i] = decoded_dest == i  AND  write_enable_overall
                                     AND  (i != 0)   ← always block writes to x0
```

One extra AND gate. Almost no cost. Enormous benefit to ISA cleanliness.

### Quick Check 7.1

> 1. How would you copy register x3 into x7 in RISC-V, using only ADD and the zero register?
> 2. What happens in hardware when a RISC-V instruction writes to x0?
> 3. ARM64 has xzr (zero register). What happens when you write to xzr?

---

## 8. Register Numbers in Machine Code

When the CPU executes `ADD x5, x1, x2`, how does the instruction encoding specify which registers to use? The answer is simple: register numbers are encoded as binary integers in the instruction word.

### The Bit Math

RISC-V has 32 registers (x0–x31). To identify one register, you need enough bits to count to 31:

```
2^5 = 32, so 5 bits can specify any register from 0 to 31.
```

A standard RISC-V R-type instruction encoding uses exactly this:

```
Bits:   31..25   24..20   19..15   14..12   11..7    6..0
Field:  funct7   rs2      rs1      funct3   rd       opcode
Width:  7 bits   5 bits   5 bits   3 bits   5 bits   7 bits
                 ↑        ↑                 ↑
               source2  source1          destination
               register register         register
```

For `ADD x5, x1, x2` (x5 = x1 + x2):
- rd = 5 = `00101` (destination: x5)
- rs1 = 1 = `00001` (source 1: x1)
- rs2 = 2 = `00010` (source 2: x2)
- funct7 = `0000000`, funct3 = `000`, opcode = `0110011` (integer register-register op)

```
31      25 24   20 19   15 14  12 11    7 6      0
0000000  | 00010 | 00001 | 000 | 00101 | 0110011
funct7     rs2     rs1   f3     rd      opcode
```

That 32-bit pattern uniquely encodes "add x5 = x1 + x2". Every field is at a fixed bit position, so the decoder can extract all three register numbers simultaneously with simple wiring.

### x86-64: The Cost of Complexity

x86-64 has 16 registers. That is only 4 bits per register. But x86 instructions are variable-length (1 to 15 bytes), and registers are encoded in a complex combination of ModRM bytes, REX prefixes, and SIB bytes:

```
REX prefix: 0100WRXB
                W=1: 64-bit operation
                R: extends ModRM reg field by 1 bit (for R8-R15)
                X: extends SIB index field
                B: extends ModRM r/m field (for R8-R15)

ModRM byte: [mod(2) | reg(3) | r/m(3)]
```

The REX prefix adds a 4th bit to the 3-bit register fields, expanding them from 0–7 (original 8 registers) to 0–15 (all 16 registers). This was the clever trick AMD used to double x86's register count in 64-bit mode without breaking backward compatibility.

### ARM64: Clean and Regular

ARM64 instructions are uniformly 32 bits. Each register field is exactly 5 bits (for x0–x30 plus xzr/sp):

```
ADD x5, x1, x2 in ARM64:
Bits:   31..29  28..24   23..22  20..16  15..10   9..5    4..0
Field:  sf+opc  10001    shift   Rm       imm6    Rn       Rd
        1 0001  0 0000    00    00010    000000   00001   00101
```

Five bits for Rd (destination), five for Rn (source 1), five for Rm (source 2). Uniform, predictable, easy to decode in hardware.

### Quick Check 8.1

> 1. How many bits are needed to address one of 32 RISC-V registers?
> 2. A RISC-V R-type instruction has three 5-bit register fields. How many bits total do they consume?
> 3. Why does x86-64 need a REX prefix byte to access registers R8 through R15?

---

## 9. The Register File: Multi-Ported SRAM

The register file is the hardware block that physically implements all the CPU's registers. It is not just a collection of flip-flops — it is a carefully engineered memory array designed for extremely fast, multi-simultaneous access.

### What "Multi-Ported" Means

Executing `ADD x5, x1, x2` in a single clock cycle requires:
1. Reading x1 (64 bits) simultaneously with...
2. Reading x2 (64 bits)
3. Writing the result back to x5 (64 bits)

These three operations must happen in one clock cycle. The register file therefore needs **two read ports** and **one write port** — a three-port structure.

```
                   ┌─────────────────────────────────────┐
Read Address 1 ───►│  5-to-32           ┌──────────────┐ │
  (5 bits)         │  Decoder 1    ────►│  x0 (64 bit) │ │
                   │               ────►│  x1 (64 bit) │ │──► Read Data 1
Read Address 2 ───►│  5-to-32      ────►│  x2 (64 bit) │ │    (64 bits)
  (5 bits)         │  Decoder 2    ────►│  x3 (64 bit) │ │
                   │               ────►│  ...         │ │──► Read Data 2
Write Address  ───►│  5-to-32      ────►│  x31(64 bit) │ │    (64 bits)
  (5 bits)         │  Decoder 3         └──────────────┘ │
Write Data     ───►│  (64 bits)                          │
Write Enable   ───►│                                     │
                   └─────────────────────────────────────┘
```

### Built from SRAM

Each register is a row of **SRAM cells** (Static Random Access Memory). Unlike the DRAM in your computer's RAM, SRAM does not need to be periodically refreshed — it holds its value as long as power is applied. SRAM is faster but uses more transistors per bit (6 transistors per bit vs. 1 for DRAM), which is why you cannot put 32 GB of SRAM in a computer affordably.

The register file is typically the fastest SRAM on the chip — designed for minimum latency above all else, accepting the area cost.

### The Area Cost of Additional Ports

Each additional read or write port approximately doubles the area of the register file. This is because each port requires its own set of sense amplifiers, address decoders, and bit lines running the full length of the array:

```
Ports   Approximate Area (relative to 1 port)
────────────────────────────────────────────────
1 read  1×
3 ports (2R + 1W)   ~4–6×
6 ports (4R + 2W)   ~12–16×
```

This is why superscalar CPUs — which can issue multiple instructions per cycle and therefore need more simultaneous register accesses — face a serious area problem as they try to go wider. A CPU that issues 8 instructions per cycle might need 16 read ports and 8 write ports. That register file would be enormous.

### The Solution: Banked Register Files

High-end CPUs often **bank** the register file: split it into multiple identical copies (banks), each serving different execution units. The same value may be replicated across banks. This adds write complexity (all banks must be updated on writes) but allows reads to proceed in parallel without port contention.

### Quick Check 9.1

> 1. Why does executing `ADD x5, x1, x2` require three simultaneous register file operations?
> 2. What type of memory (SRAM or DRAM) is used for register files, and why?
> 3. If a 3-port register file has area 4× a 1-port file, roughly what area would a 6-port file have?

---

## 10. Register Spilling: When You Run Out of Space

A function with 50 local variables cannot keep all of them in registers if the CPU only has 32. The compiler must decide which variables live in registers and which live in memory. Variables temporarily evicted to memory are said to be **spilled**.

### The Cost of Spilling

A register access takes one cycle. A memory access — even to L1 cache — takes 4–5 cycles. A cache miss takes 100+ cycles. Spilling a hot variable (one accessed frequently) can turn a loop that executes in 1 cycle per iteration into one that executes in 5 or more.

### What Gets Spilled?

Compilers use **liveness analysis**: a variable is *live* at a given point in the program if its current value might be read in the future. Variables that are live for a long time compete for registers throughout their lifetime. When two variables are simultaneously live and there are not enough registers, one must be spilled.

The compiler chooses which variable to spill based on heuristics:
- **Least Recently Used**: spill the variable that has not been accessed for the longest time
- **Lowest Use Frequency**: spill the variable used least often (especially useful to spill variables not inside loops)
- **Shortest Remaining Lifetime**: spill the variable that will become dead soonest

```
Live ranges (simplified example with 3 registers available):

Variable:    a  b  c  d  e
Line 1:      ████
Line 2:      ████ ████
Line 3:      ████ ████ ████
Line 4:           ████ ████ ████
Line 5:                ████ ████ ████ ← 3 vars live, need 3 registers — OK
Line 6:                     ████ ████ ████ ← a is already dead, d can use its slot
```

### Register Allocation as Graph Coloring

The classic formulation: create an **interference graph** where each variable is a node, and draw an edge between two variables if they are simultaneously live (cannot share a register). Then **color the graph** with k colors (k = number of registers) such that no two adjacent nodes share a color.

If the graph can be k-colored, the coloring is the register assignment. If not, some nodes must be spilled until the graph becomes k-colorable.

```
Variables a, b, c, d with live ranges:
a and b overlap → edge a–b
b and c overlap → edge b–c
c and d overlap → edge c–d
a and d do NOT overlap → no edge

Interference graph:    a — b — c — d

Colors (registers):    R1  R2  R1  R2   ← only 2 registers needed!
                                         (a and c don't interfere)
```

---

## 11. A Comparison Table: x86-64, ARM64, and RISC-V

Here is a comprehensive side-by-side comparison of the three most important architectures today:

### Architecture Overview

| Feature                     | x86-64 (Intel/AMD)           | ARM64 (AArch64)              | RISC-V (RV64I)              |
|-----------------------------|------------------------------|------------------------------|-----------------------------|
| Integer GPRs                | 16                           | 31 (+ XZR)                   | 32 (including x0=zero)      |
| Register width              | 64-bit (with 8/16/32 views)  | 64-bit (with 32-bit w0–w30)  | 64-bit (with 32-bit subset) |
| Zero register               | No (use XOR rax,rax trick)   | Yes: XZR / WZR               | Yes: x0                     |
| Dedicated PC                | Yes: RIP (not a GPR)         | Yes: PC (not a GPR)          | Yes: PC (not a GPR)         |
| Flags register              | RFLAGS (ZF, CF, OF, SF, PF)  | NZCV flags (separate)        | None (compare+branch fusion) |
| Link register               | No (return addr on stack)    | X30 (LR)                     | x1 (RA)                     |
| Stack pointer               | RSP (GPR slot 4)             | SP (separate from x31)       | x2 (SP, soft convention)    |
| Frame pointer               | RBP (soft convention)        | X29 (FP, soft convention)    | x8 / s0 (soft convention)   |
| Instruction size            | Variable (1–15 bytes)        | Fixed 32-bit                 | Fixed 32-bit (or 16-bit C)  |
| Bits to encode a register   | 4 bits (3+REX extension)     | 5 bits                       | 5 bits                      |
| Calling convention name     | System V AMD64 ABI (Linux)   | AAPCS64                      | RISC-V psABI                |
| Argument registers          | rdi, rsi, rdx, rcx, r8, r9   | x0–x7                        | a0–a7 (x10–x17)             |
| Return value registers      | rax (+ rdx for 128-bit)      | x0 (+ x1 for 128-bit)        | a0 (x10) + a1 (x11)         |
| Callee-saved registers      | rbx, rbp, r12–r15            | x19–x28, x29, x30            | s0–s11 (x8–x9, x18–x27)     |

### General-Purpose Register Names

| Number | x86-64           | ARM64      | RISC-V ABI | RISC-V Role              |
|--------|------------------|------------|------------|--------------------------|
| 0      | RAX              | x0         | zero       | Hardwired 0              |
| 1      | RCX              | x1         | ra         | Return address           |
| 2      | RDX              | x2         | sp         | Stack pointer            |
| 3      | RBX              | x3         | gp         | Global pointer           |
| 4      | RSP (stack ptr)  | x4         | tp         | Thread pointer           |
| 5      | RBP (frame ptr)  | x5         | t0         | Temporary                |
| 6      | RSI              | x6         | t1         | Temporary                |
| 7      | RDI              | x7         | t2         | Temporary                |
| 8      | R8               | x8         | s0 / fp    | Saved / Frame pointer    |
| 9      | R9               | x9         | s1         | Saved register           |
| 10     | R10              | x10        | a0         | Argument / return value  |
| 11     | R11              | x11        | a1         | Argument / return value  |
| 12     | R12              | x12        | a2         | Argument                 |
| 13     | R13              | x13        | a3         | Argument                 |
| 14     | R14              | x14        | a4         | Argument                 |
| 15     | R15              | x15        | a5         | Argument                 |
| 16–28  | —                | x16–x28    | a6–a7, s2–s9| Arguments / Saved       |
| 29     | —                | x29 (FP)   | s10        | Saved register           |
| 30     | —                | x30 (LR)   | s11        | Saved register           |
| 31     | —                | xzr / SP   | t6 / x31   | Temporary / zero         |

*(x86-64 only has registers 0–15; entries 16–31 are not applicable)*

### Quick Check 11.1

> 1. ARM64 and RISC-V both have 31–32 general-purpose registers. x86-64 has 16. What is one instruction-encoding reason x86 did not add more registers in the 32-bit era?
> 2. RISC-V has no flags register. How does it implement "branch if x5 > x6"?
> 3. In ARM64, x30 is the link register. What does this mean for leaf functions?

---

## 12. Register Renaming: A Preview

Modern CPUs execute instructions **out of order** — they look ahead in the instruction stream and execute later instructions early if their dependencies are ready. This requires more physical registers than the architecture exposes to programs.

### The False Dependency Problem

Consider this instruction sequence:

```
① ADD  x1, x2, x3    # x1 = x2 + x3
② SUB  x1, x4, x5    # x1 = x4 - x5  (overwrites x1!)
③ MUL  x6, x1, x7    # x6 = x1 × x7  (uses x1 from ②)
④ AND  x8, x1, x9    # x8 = x1 AND x9 (uses x1 from ②)
```

Instruction ③ must wait for ② to finish (true dependency — it needs the result). But what about ①? It also writes to x1. Instruction ③ does NOT need ①'s x1 — it needs ②'s x1. Yet naively, the CPU might think ③ and ④ need to wait until after both ① and ② complete.

This is a **Write-After-Write (WAW) hazard**: two instructions both write to the same architectural register. It creates a false dependency — ③ is forced to wait even though it does not actually need ①'s result.

### Register Renaming Solution

The CPU maintains a **Register Alias Table (RAT)** that maps each architectural register to a physical register:

```
Before executing ①:
  RAT: x1 → phys_r17,  x2 → phys_r22,  x3 → phys_r08, ...

Decode ①: ADD x1, x2, x3
  CPU allocates a new physical register for x1's new value: phys_r55
  Renamed: ADD phys_r55, phys_r22, phys_r08
  RAT updated: x1 → phys_r55

Decode ②: SUB x1, x4, x5
  CPU allocates another new physical register: phys_r61
  Renamed: SUB phys_r61, phys_r19, phys_r33
  RAT updated: x1 → phys_r61

Decode ③: MUL x6, x1, x7
  RAT lookup for x1: → phys_r61  (from ②, correct!)
  Renamed: MUL phys_r40, phys_r61, phys_r27
  ③ only depends on ② — it can execute as soon as ② finishes

Decode ④: AND x8, x1, x9
  RAT lookup for x1: → phys_r61  (same as ③, correct!)
  ③ and ④ can execute in parallel (different physical destinations)!
```

Now ① is completely independent from ②③④. All four instructions can be in flight simultaneously, limited only by true data dependencies, not register name reuse.

Modern CPUs (Intel Alder Lake, Apple M2, AMD Zen 4) have 280–360 physical integer registers mapped to 16 or 32 architectural registers. This pool of extra physical registers enables the instruction window to be large enough to find parallelism across dozens of instructions simultaneously.

We will explore register renaming in full detail in the chapter on out-of-order execution (Chapter 26). For now, the key insight is: **registers are not just hardware memory cells — they are a naming abstraction, and the CPU can maintain many simultaneous "versions" of the same architectural register name.**

### Quick Check 12.1

> 1. What is a WAW (Write-After-Write) hazard?
> 2. If x86-64 has 16 architectural registers but a modern Intel CPU has 280 physical registers, how many extra physical registers are available for renaming?
> 3. After register renaming, can instructions ① and ② in the example above execute in parallel?

---

## Summary

Registers are the most intimate memory in a computer — so close to the arithmetic units that they add no measurable latency to computation. Here is what this chapter covered:

**Why they exist:** Registers bridge the 300× speed gap between the CPU and RAM. A value loaded from RAM into a register can be used hundreds of times without further memory traffic.

**The analogy:** Registers are the chef's cutting board — tiny, immediately accessible, holding only what is being actively worked on. RAM is the pantry. The SSD is the warehouse.

**How many:** x86-64 has 16, ARM64 has 31 (+XZR), RISC-V has 32 (including hardwired zero). The count reflects tradeoffs between instruction encoding width, silicon area, context-switch cost, and compiler register allocation effectiveness.

**Register width:** Modern 64-bit CPUs use 64-bit registers. Sub-register access (32-bit, 16-bit, 8-bit views) exists for compatibility but can cause subtle hazards in x86.

**General-purpose registers** follow **calling conventions** (ABI) that specify which registers pass arguments, return values, and must be preserved across calls. This enables modular, interoperable code.

**Special-purpose registers:**
- **PC (Program Counter):** address of the next instruction
- **SP (Stack Pointer):** top of the call stack
- **FP (Frame Pointer):** fixed reference point within a stack frame
- **LR/RA (Link Register):** return address for fast leaf function calls
- **IR (Instruction Register):** holds the current instruction during decode
- **Flags/Status Register:** zero, negative, carry, overflow bits from ALU ops
- **CSRs:** privileged control registers in RISC-V

**Zero register:** RISC-V x0 and ARM64 xzr are hardwired to zero. Writes are silently discarded. This simplifies the ISA enormously, enabling MOV, NOP, and unconditional branches using ordinary ADD and BEQ instructions.

**Machine code encoding:** 5 bits address one of 32 registers. A standard RISC-V instruction uses three 5-bit register fields (source1, source2, destination) totaling 15 bits of a 32-bit instruction.

**The register file** is multi-ported SRAM. A 3-port file (2 reads + 1 write) executes one instruction per cycle. More ports cost exponentially more area — a key limit on superscalar width.

**Register spilling** occurs when the compiler runs out of registers. The compiler's register allocator models the problem as graph coloring, assigning colors (registers) to nodes (variables) and spilling when the graph cannot be colored.

**Register renaming** (preview): Modern CPUs rename architectural registers to a larger pool of physical registers, eliminating false dependencies and enabling out-of-order execution across a window of many instructions simultaneously.

---

## Exercises

### Easy

1. **Bit counting.** A RISC-V instruction has three register fields (rd, rs1, rs2), each 5 bits wide. How many bits do they consume in a 32-bit instruction? What percentage of the instruction word is register addressing?

2. **Zero register uses.** Write RISC-V assembly for each of the following using only ADD, SUB, and BEQ instructions (plus the zero register x0):
   - Copy x3 into x7
   - Set x5 to zero
   - Branch to label `done` if x10 equals x11
   - Negate x6 (compute -x6, store back in x6)

3. **Calling convention.** In the RISC-V ABI, function foo() calls bar(). foo() uses register t0 and s2. bar() uses register t1 and s3.
   - Which registers must foo() save before calling bar()?
   - Which registers must bar() save before using its registers?
   - Explain the difference in terms of caller-saved vs. callee-saved.

### Medium

4. **Instruction decoding.** The following RISC-V machine code word is given in binary:
   ```
   0000000  00011  00010  000  00001  0110011
   funct7   rs2    rs1    f3   rd     opcode
   ```
   - What are rs1, rs2, and rd as decimal register numbers?
   - The opcode `0110011` with funct3=`000` and funct7=`0000000` is ADD. Write the assembly instruction.
   - If funct7 were `0100000` instead (same opcode, same funct3), the instruction is SUB. Write that assembly instruction.

5. **Register file ports.** A superscalar CPU issues 4 instructions per cycle. Each instruction needs 2 source reads and 1 destination write.
   - How many read ports and write ports does the register file need in the worst case?
   - If a 3-port register file has area A, and each additional port doubles area, estimate the area of this register file relative to a 3-port design.
   - Suggest one alternative to using a single large register file with 12 ports.

6. **x86 partial register hazard.** Consider this x86-64 sequence:
   ```assembly
   mov  rax, 0xFFFFFFFFFFFFFFFF    ; rax = all ones (64 bits)
   mov  ax,  0x1234                ; write to lower 16 bits
   mov  rax into memory            ; what value is stored?
   ```
   - What is the value of RAX after the second instruction?
   - Now repeat with `mov eax, 0x00001234`. What is the value of RAX? (Hint: writing to EAX zeroes the upper 32 bits.)

### Hard

7. **Graph coloring register allocation.** Consider the following pseudocode:
   ```
   a = input1
   b = input2
   c = a + b
   d = a × 3
   e = c - d
   f = b + e
   output = f × d
   ```
   - Draw the live range table: for each variable, mark the lines where it is live (from its definition to its last use).
   - Build the interference graph: draw an edge between any two variables that are simultaneously live on at least one line.
   - What is the minimum number of registers (colors) needed?
   - If only 3 registers are available, which variable must be spilled? Justify your choice.

8. **Register renaming deep dive.** A CPU has 16 architectural registers (r0–r15) and 64 physical registers (p0–p63). The RAT maps each architectural register to a physical register. Initially, ri maps to pi for all i.

   The following instructions are decoded in order:
   ```
   ① ADD  r1, r2, r3     # r1 = r2 + r3
   ② MUL  r4, r1, r5     # r4 = r1 × r5  (true dependency on ①)
   ③ SUB  r1, r6, r7     # r1 = r6 - r7  (WAW with ①)
   ④ OR   r8, r1, r9     # r8 = r1 OR r9  (depends on ③'s r1)
   ⑤ AND  r10, r1, r11   # r10 = r1 AND r11  (also depends on ③'s r1)
   ```

   The CPU allocates physical registers p16, p17, p18, p19, p20 for new values.

   - Show the RAT state after each instruction is decoded (all 5 instructions).
   - Which instructions can execute in parallel?
   - When instruction ① retires (commits), what happens to the old p1? When can it be returned to the free list?
   - If the CPU mispredicts a branch after ②, how must the RAT be restored?

9. **Architectural design tradeoff essay.** RISC-V has 32 registers while x86-64 has 16.
   - Quantify the instruction encoding cost: how many bits per instruction does each ISA use for three-register instructions?
   - Explain how more registers reduce spilling but increase context-switch cost. Give a formula for context-switch overhead if saving each register takes T nanoseconds.
   - A server process switches context 1,000 times per second. Estimate the difference in context-switch CPU time between a 16-register and 32-register architecture, assuming T = 1 ns per register save/restore and a 3 GHz CPU.
   - Is 32 registers the right number, or should future ISAs have 64? Argue both sides.
