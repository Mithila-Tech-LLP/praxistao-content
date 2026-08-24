# Chapter 19: RISC-V — The Open Source Revolution

Imagine a world where every programming language was owned by a single corporation, and you had to pay a royalty every time you wrote a line of code. That sounds absurd — but that was precisely the situation for processor architectures for decades. Every CPU design team, every startup, every nation building a chip had to pay tribute to Intel, AMD, or ARM. Then, in 2010, a group of researchers at UC Berkeley decided enough was enough. They would give the world a free ISA — not free as in cheap, but free as in freedom. That project became RISC-V, and it is quietly reshaping the entire computing landscape.

## Table of Contents

1. [The Problem With Proprietary ISAs](#1-the-problem-with-proprietary-isas)
2. [RISC-V Is Born: UC Berkeley 2010](#2-risc-v-is-born-uc-berkeley-2010)
3. [Design Principles: Clean Slate Thinking](#3-design-principles-clean-slate-thinking)
4. [Registers: The 32-Register File](#4-registers-the-32-register-file)
5. [The Base ISA: RV32I and RV64I](#5-the-base-isa-rv32i-and-rv64i)
6. [Standard Extensions: Building Blocks](#6-standard-extensions-building-blocks)
7. [Instruction Encoding: Elegance in 32 Bits](#7-instruction-encoding-elegance-in-32-bits)
8. [Privilege Levels: Security Layers](#8-privilege-levels-security-layers)
9. [Growing Adoption: Who Is Using RISC-V?](#9-growing-adoption-who-is-using-risc-v)
10. [Geopolitical Significance](#10-geopolitical-significance)
11. [Comparison: RISC-V vs ARM vs x86](#11-comparison-risc-v-vs-arm-vs-x86)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. The Problem With Proprietary ISAs

### What Is an ISA License?

Recall from Chapter 14 that an **Instruction Set Architecture (ISA)** is the contract between software and hardware — it defines the set of instructions a processor understands. Writing software for a processor requires knowing its ISA. Building a processor requires implementing its ISA.

Here is the catch: the two dominant ISAs — x86 and ARM — are **proprietary**. They are owned by corporations, and using them requires a license.

Think of it like this: imagine that the English alphabet was patented. Every time you wrote a letter, you owed a royalty to the patent holder. That is roughly the situation chip designers faced with processor ISAs.

```
x86 ISA:
  Owner:         Intel (AMD has a cross-license, no one else does)
  License cost:  Not available at any price to new entrants
  Restriction:   Intel can refuse licensing entirely
  Reality:       If you want to make an x86 chip, you cannot.
                 Intel has blocked every attempt.

ARM ISA:
  Owner:         ARM Holdings (now owned by SoftBank, then Nvidia tried to buy)
  License cost:  Architecture license: several million dollars upfront
                 Core license: varies by tier
  Royalty:       ~0.5-2% of chip selling price, per chip shipped
  Restriction:   ARM can revoke licenses (happened during Qualcomm dispute)
  Reality:       Every Apple, Qualcomm, Samsung chip pays ARM.
```

For a startup shipping 10 million chips at $2 each, a 1% ARM royalty is $200,000 per year — before making a single dollar of profit. For a country trying to build a domestic chip industry, paying a foreign company for every chip is a permanent dependency.

### The Research Problem

Beyond money, there was a deeper problem: **research**. Computer architecture professors at universities needed an ISA for their work. They needed to:

- Add new instructions to test new ideas
- Build simulated processors
- Publish papers with real implementations
- Teach students how ISAs work

Using x86 was legally impossible — Intel's license prohibited compatible designs. ARM required license negotiations that took months and cost money no research lab had. MIPS was declining. SPARC was controlled by Oracle. Every ISA had a gatekeeper.

Professor Krste Asanović at UC Berkeley described the situation in 2010: his research group spent more time fighting ISA licensing issues than doing actual architecture research.

### Quick Check 1.1

> 1. What is an ISA license and why do chip designers need one?
> 2. Why could university researchers not just use x86 for their architecture research?
> 3. If ARM charges 1% royalty on a $3 chip, how much does a company shipping 50 million chips per year pay?

---

## 2. RISC-V Is Born: UC Berkeley 2010

### The Team

In 2010, Professor Krste Asanović, graduate student Andrew Waterman, and Professor David Patterson began designing a new ISA from scratch. Patterson was already a legend in computer architecture — he was one of the original inventors of RISC in the 1980s and co-author of the most widely used computer architecture textbook in the world (the "Patterson and Hennessy" book).

The project was called **RISC-V** — pronounced "risk five." The name follows Berkeley's tradition:

```
RISC-I   (1982) — First Berkeley RISC processor
RISC-II  (1983) — Improved version
SOAR     (1984) — Smalltalk On A RISC
SPUR     (1988) — Symbolic Processing Using RISC
RISC-V   (2010) — The fifth generation, this time designed for real deployment
```

The Roman numeral V is intentional — it signals continuity with a 30-year research tradition, while starting completely fresh on the ISA design.

### The Goal: An ISA For Everyone

The Berkeley team had a striking insight: they were not just designing a research ISA. If they designed it carefully enough — cleanly enough, openly enough — it could be the ISA that the whole world used. Not because Berkeley forced anyone, but because it would simply be the best free choice available.

Their goals were ambitious:

```
  For embedded systems:    Tiny, simple, low-power (compete with ARM Cortex-M)
  For smartphones:         Efficient, 64-bit (compete with ARM Cortex-A)
  For servers:             High-performance, scalable (compete with x86)
  For research:            Extensible, formally specified
  For education:           Simple enough to teach in a semester
  License:                 Completely free, forever, for everyone
```

### Open Source Hardware

The key decision was licensing. The RISC-V ISA specification was released under a **Creative Commons license** — you can read it, implement it, build chips with it, sell those chips, and never pay a cent to anyone. The specification itself was eventually placed in the public domain.

This was unprecedented for a processor ISA. Think of it like the Linux kernel moment for hardware: suddenly, anyone anywhere in the world could build a complete processor without permission from any gatekeeper.

The analogy to open-source software is apt. Linux started as a hobby project by a Finnish student and eventually ran the majority of the world's servers. RISC-V started as a Berkeley research project and is following a similar trajectory — slowly at first, then suddenly everywhere.

### A Critical Timeline

```
2010 — RISC-V design begins at UC Berkeley
2011 — First RISC-V chip taped out (for class project: 5-month design)
2014 — RISC-V ISA specification v1.0 released publicly
2015 — RISC-V Foundation founded; SiFive founded (Berkeley spinout)
2016 — RISC-V Workshop attracts industry attention; NVIDIA adopts RISC-V
2018 — Linux kernel RISC-V support merged (kernel 4.15)
       GCC and LLVM support ratified
2019 — RISC-V Foundation becomes RISC-V International (moved to Switzerland)
       Western Digital announces 1 billion RISC-V cores deployed
2021 — RISC-V Vector (V) extension ratified
       ESP32-C3 (first mass-market RISC-V Wi-Fi chip) ships
2022 — Intel Foundry Services announces RISC-V services
       SiFive Performance P550 competes with ARM Cortex-A75
2023 — Google, Qualcomm, NVIDIA all shipping RISC-V in products
       India's SHAKTI program tapes out on commercial process
2024 — RISC-V chips appear in data centers; automotive adoption accelerates
```

In just 14 years, RISC-V went from a classroom project to running in data centers, smartphones, hard drives, and billions of IoT devices. This is the fastest adoption curve of any processor ISA in history.

### Quick Check 2.1

> 1. Who were the three main people behind RISC-V's creation?
> 2. What does "RISC-V" mean? Why is it called that?
> 3. What license allows anyone to implement RISC-V for free?

---

## 3. Design Principles: Clean Slate Thinking

When you design something from scratch — with no backwards compatibility requirements, no political constraints, no legacy customers — you have a rare opportunity. The RISC-V designers took full advantage of it. They studied every previous ISA, catalogued every mistake, and then made deliberate decisions about each one.

### Principle 1: Turing-Complete Base, Everything Else Optional

The base RISC-V ISA (called RV32I for 32-bit, RV64I for 64-bit) is the smallest set of instructions that makes a processor **Turing-complete** — capable of computing anything computable. Only 47 instructions.

Think of it like a kitchen with only the absolutely essential tools: one knife, one cutting board, one pan, one spatula. You can cook any recipe with just those. The fancy appliances (stand mixer, sous vide machine, immersion blender) are optional — you add them based on what you cook most.

The RISC-V extensions are those fancy appliances. A simple microcontroller might just need the base ISA. A machine learning accelerator needs the base plus vector extensions plus custom neural network instructions.

### Principle 2: Fixed-Length Instructions (With Optional Compression)

The base ISA uses fixed 32-bit instructions. Every instruction is exactly the same width. This makes decoding trivial:

```
Memory:  [  inst 0  ][  inst 1  ][  inst 2  ][  inst 3  ]
         32 bits      32 bits      32 bits      32 bits
         
No ambiguity. No variable-length decoding complexity.
```

Compare to x86, where an instruction can be 1 to 15 bytes long. Decoding x86 requires a complex state machine just to find where each instruction starts. That complexity costs transistors, power, and chip area.

If code density matters (embedded systems with limited flash memory), the optional **C extension** adds 16-bit compressed versions of the most common instructions. The hardware can mix 16-bit and 32-bit instructions freely — they're distinguishable by the two lowest bits.

### Principle 3: Load-Store Architecture

Like all RISC designs, RISC-V only operates on data in registers. To add two values in memory, you must:

1. Load value 1 from memory into a register
2. Load value 2 from memory into a register
3. Add the registers
4. Store the result back to memory

No instruction can directly add two memory locations. This simplifies the processor pipeline enormously — every execute stage works on registers, and memory access happens in dedicated load/store units.

```
RISC-V (Load-Store):                x86 (Register-Memory):
  lw   t0, 0(a0)    # load          addl (%rsi), %eax    # one instruction
  lw   t1, 4(a0)    # load          (but one very complex instruction)
  add  t2, t0, t1   # compute
  sw   t2, 8(a0)    # store
  
4 simple instructions                1 complex instruction
4 simple hardware operations         1 very complex hardware operation
```

The RISC approach generates more instructions, but each instruction executes in a single pipeline stage. The net performance is similar or better, with much simpler hardware.

### Principle 4: Orthogonal Design

RISC-V instructions are **orthogonal** — the fields of each instruction format mean the same thing across all instructions. The source register 1 is always bits [19:15]. The destination register is always bits [11:7]. The opcode is always bits [6:0].

This orthogonality means the decode logic is simple and regular. A hardware implementation can read all register addresses before fully decoding the instruction — important for performance optimization.

### Principle 5: Ratified Extensions Are Frozen

Once a RISC-V extension is ratified by RISC-V International, it is **frozen forever**. The instructions in the extension will never change meaning, never be deprecated, never be removed.

Compare to ARM, which deprecated the Thumb-1 ISA, changed behavior between ARMv7 and ARMv8, and has accumulated complex compatibility rules. Or x86, where instructions like `XCHG AX, AX` were originally a NOP but accumulated special behaviors over decades.

RISC-V's commitment to stability means code compiled today will run on RISC-V chips built in 20 years. For embedded systems that must operate for decades (industrial equipment, infrastructure), this is enormously valuable.

### Principle 6: Explicitly Reserved Custom Extension Space

The RISC-V ISA specification deliberately reserves opcode space for custom extensions. A company can add their own instructions — for cryptography, machine learning, signal processing — without conflicting with standard extensions:

```
Opcode space allocation:
  Standard:   Used by RISC-V International for standard extensions
  Reserved:   Will be used for future standard extensions
  Custom-0:   0x0B — Reserved for implementor-defined instructions
  Custom-1:   0x2B — Reserved for implementor-defined instructions
  Custom-2:   0x5B — Reserved for implementor-defined instructions
  Custom-3:   0x7B — Reserved for implementor-defined instructions
```

A chip maker can add up to hundreds of custom instructions in these spaces, build custom toolchain support, and ship specialized hardware — while still being fully compatible with standard RISC-V software.

### Quick Check 3.1

> 1. What does "Turing-complete" mean, and why does the RISC-V base ISA aim to be Turing-complete alone?
> 2. Why is fixed-length instruction encoding simpler for hardware than variable-length?
> 3. What does "ratified and frozen" mean for RISC-V extensions? Why is this good for embedded systems?

---

## 4. Registers: The 32-Register File

### The Hardware

RISC-V has 32 general-purpose integer registers, named **x0 through x31**. Each is 32 bits wide in RV32I, or 64 bits wide in RV64I. This is the same register count as MIPS and more than x86-64 (which has 16).

There is one magical register: **x0 is hardwired to the value zero**. You can read it any time and always get 0. If you write to x0, the write is silently discarded. This eliminates the need for a separate "load zero" instruction — just read x0.

```
Integer register file (RV32I):

  x0  │ 0x00000000 │  always zero (hardwired)
  x1  │            │
  x2  │            │  
  x3  │            │
  ...
  x31 │            │

  32 registers × 32 bits = 1,024 bits of fast storage
```

Why is hardwired zero useful? Many operations need a zero operand:

```riscv
# Move register (copy t1 into t0):
add  t0, t1, x0    # t0 = t1 + 0 = t1

# Negate register:
sub  t0, x0, t1    # t0 = 0 - t1 = -t1

# Compare to zero (branch if t0 == 0):
beq  t0, x0, label # branch if t0 == x0 == 0

# Clear register:
add  t0, x0, x0    # t0 = 0 + 0 = 0
```

All of these are "free" — they use the normal instructions with x0 as an operand, rather than requiring special-case instructions. This is the elegance of good ISA design: a single design choice (hardwired zero) eliminates many special-case instructions.

### The ABI Register Names

The hardware calls the registers x0-x31, but programmers use **ABI names** — human-readable names that describe each register's purpose. The ABI (Application Binary Interface) is the convention that all programs follow so they can call each other's functions correctly.

Think of the ABI as a formal agreement: "register x10 will be used for the first argument to a function, and x11 for the second." As long as everyone follows this agreement, functions written by different people (or different compilers) can call each other seamlessly.

```
Register  ABI Name  Role                          Saved by
────────  ────────  ────────────────────────────  ────────
x0        zero      Hardwired zero                —
x1        ra        Return address                Caller
x2        sp        Stack pointer                 Callee
x3        gp        Global pointer                —
x4        tp        Thread pointer                —
x5        t0        Temporary                     Caller
x6        t1        Temporary                     Caller
x7        t2        Temporary                     Caller
x8        s0/fp     Saved / Frame pointer         Callee
x9        s1        Saved                         Callee
x10       a0        Argument 0 / Return value     Caller
x11       a1        Argument 1 / Return value 2   Caller
x12       a2        Argument 2                    Caller
x13       a3        Argument 3                    Caller
x14       a4        Argument 4                    Caller
x15       a5        Argument 5                    Caller
x16       a6        Argument 6                    Caller
x17       a7        Argument 7 / Syscall number   Caller
x18       s2        Saved                         Callee
x19       s3        Saved                         Callee
x20       s4        Saved                         Callee
x21       s5        Saved                         Callee
x22       s6        Saved                         Callee
x23       s7        Saved                         Callee
x24       s8        Saved                         Callee
x25       s9        Saved                         Callee
x26       s10       Saved                         Callee
x27       s11       Saved                         Callee
x28       t3        Temporary                     Caller
x29       t4        Temporary                     Caller
x30       t5        Temporary                     Caller
x31       t6        Temporary                     Caller
```

### Understanding the Three Register Groups

**Argument registers (a0-a7):** These carry function arguments when calling a function, and return values when returning from one. If a function has 3 parameters, they go in a0, a1, and a2. The return value comes back in a0.

**Saved registers (s0-s11):** These are preserved across function calls. If you store a value in s0 before calling a function, it will still be there after the function returns — the called function (callee) is responsible for saving and restoring s registers if it uses them.

**Temporary registers (t0-t6):** These are scratch space. After calling a function, temporaries may have been overwritten — the caller cannot assume they are preserved. Good for intermediate computations that you only need within a single function.

The **ra (return address)** register holds the address to jump back to after a function completes. When you call a function with `jal ra, function_label`, it saves the next instruction's address into ra. The function ends with `ret` (which is actually `jalr x0, ra, 0` — jump to the address in ra, discarding the old ra into x0).

**The stack pointer (sp)** always points to the current top of the call stack. Every function that needs local variables decrements sp to allocate space on the stack, uses that space, and increments sp when done.

### Floating-Point Registers

The F and D extensions (single-precision and double-precision floating point) add a separate register file:

```
f0  – f31:    32 floating-point registers, each 64 bits wide (with D extension)
              (32 bits wide if only F extension is present)

ABI names:
  ft0-ft7   :  Floating-point temporaries (caller-saved)
  fs0-fs11  :  Floating-point saved registers (callee-saved)
  fa0-fa7   :  Floating-point arguments/return values
```

Floating-point registers are separate from integer registers because floating-point hardware operates differently — it has its own arithmetic units, its own rounding modes, and its own exception flags. Mixing them would complicate the pipeline.

### Quick Check 4.1

> 1. What is special about register x0? Give two examples of how this is useful.
> 2. What is the difference between argument registers (a0-a7) and saved registers (s0-s11)?
> 3. When a function is called with `jal ra, label`, what value is stored in ra?

---

## 5. The Base ISA: RV32I and RV64I

### RV32I: 32 Registers, 32-Bit Integers

**RV32I** is the smallest, most fundamental RISC-V specification: 32 integer registers, 32-bit wide, with 47 instructions. A complete Turing-complete processor.

The 47 instructions divide into clear categories:

```
Arithmetic and Logic:
  ADD   rd, rs1, rs2     rd = rs1 + rs2
  SUB   rd, rs1, rs2     rd = rs1 - rs2
  AND   rd, rs1, rs2     rd = rs1 & rs2   (bitwise AND)
  OR    rd, rs1, rs2     rd = rs1 | rs2   (bitwise OR)
  XOR   rd, rs1, rs2     rd = rs1 ^ rs2   (bitwise XOR)
  SLL   rd, rs1, rs2     rd = rs1 << rs2  (shift left logical)
  SRL   rd, rs1, rs2     rd = rs1 >> rs2  (shift right logical, fills 0)
  SRA   rd, rs1, rs2     rd = rs1 >> rs2  (shift right arithmetic, fills sign)
  SLT   rd, rs1, rs2     rd = (rs1 < rs2) ? 1 : 0  (signed compare)
  SLTU  rd, rs1, rs2     rd = (rs1 < rs2) ? 1 : 0  (unsigned compare)

Immediate Variants (I-suffix):
  ADDI  rd, rs1, imm     rd = rs1 + imm  (imm is 12-bit signed constant)
  ANDI  rd, rs1, imm     rd = rs1 & imm
  ORI   rd, rs1, imm     rd = rs1 | imm
  XORI  rd, rs1, imm     rd = rs1 ^ imm
  SLLI  rd, rs1, shamt   rd = rs1 << shamt
  SRLI  rd, rs1, shamt   rd = rs1 >> shamt (logical)
  SRAI  rd, rs1, shamt   rd = rs1 >> shamt (arithmetic)
  SLTI  rd, rs1, imm     rd = (rs1 < imm) ? 1 : 0
  SLTIU rd, rs1, imm     unsigned version

Load Instructions:
  LW    rd, offset(rs1)  rd = Memory[rs1 + offset]  (32-bit word)
  LH    rd, offset(rs1)  rd = sign-extend(Memory[rs1+offset], 16 bits)
  LHU   rd, offset(rs1)  rd = zero-extend(Memory[rs1+offset], 16 bits)
  LB    rd, offset(rs1)  rd = sign-extend(Memory[rs1+offset], 8 bits)
  LBU   rd, offset(rs1)  rd = zero-extend(Memory[rs1+offset], 8 bits)

Store Instructions:
  SW    rs2, offset(rs1) Memory[rs1+offset] = rs2  (32-bit word)
  SH    rs2, offset(rs1) Memory[rs1+offset] = rs2[15:0]  (16-bit)
  SB    rs2, offset(rs1) Memory[rs1+offset] = rs2[7:0]   (8-bit)

Branch Instructions (conditional):
  BEQ   rs1, rs2, label  if (rs1 == rs2) PC = label
  BNE   rs1, rs2, label  if (rs1 != rs2) PC = label
  BLT   rs1, rs2, label  if (rs1 < rs2, signed) PC = label
  BLTU  rs1, rs2, label  if (rs1 < rs2, unsigned) PC = label
  BGE   rs1, rs2, label  if (rs1 >= rs2, signed) PC = label
  BGEU  rs1, rs2, label  if (rs1 >= rs2, unsigned) PC = label

Jump Instructions (unconditional):
  JAL   rd, label        rd = PC+4; PC = label  (jump and link)
  JALR  rd, rs1, imm     rd = PC+4; PC = rs1+imm (indirect jump)

Upper Immediate:
  LUI   rd, imm          rd = imm << 12  (load 20-bit upper immediate)
  AUIPC rd, imm          rd = PC + (imm << 12)  (add upper imm to PC)

Memory Fence:
  FENCE  (ordering guarantee — flush memory operations before continuing)

System:
  ECALL   (transfer to OS/hypervisor — system call)
  EBREAK  (breakpoint — transfer to debugger)

CSR Access (Control and Status Registers):
  CSRRW, CSRRS, CSRRC, CSRRWI, CSRRSI, CSRRCI
  (read/write system registers for privilege and status)
```

### LUI and AUIPC: The Large Constant Problem

Notice that ADDI only has a 12-bit immediate field — it can only encode constants from -2048 to +2047. But programs need 32-bit constants (like memory addresses). How do you load a big number?

Answer: two instructions working together.

```riscv
# Load 32-bit constant 0xDEADBEEF into register t0:

lui  t0, 0xDEADB     # t0 = 0xDEADB000  (upper 20 bits)
addi t0, t0, 0xEEF   # t0 = 0xDEADB000 + 0xEEF = 0xDEADBEEF

# Wait — 0xEEF is 3823, which is > 2047. Sign-extension problem!
# Actually for 0xEEF (negative when sign-extended from 12 bits):
lui  t0, 0xDEADC     # t0 = 0xDEADC000  (adjusted upper)
addi t0, t0, 0xEEF - 0x1000  # negative offset corrects for sign extension
# Compilers handle this automatically.
```

**AUIPC** (Add Upper Immediate to PC) is used for **position-independent code** — code that works regardless of where it is loaded in memory. Instead of loading an absolute address, you load an address relative to the current instruction's location:

```riscv
# Get address of a data symbol (works regardless of where code is loaded):
auipc t0, %pcrel_hi(my_data)  # t0 = PC + upper offset to my_data
addi  t0, t0, %pcrel_lo(my_data)  # t0 = exact address of my_data
# (compiler fills in the %pcrel_hi and %pcrel_lo values)
```

This is how shared libraries work — they load at different addresses each run (for security via ASLR), and AUIPC + ADDI pairs give them self-relative addressing.

### RV64I: The 64-Bit Extension

RV64I extends RV32I to 64 bits:

- All 32 registers become 64 bits wide
- New load/store: **LD** (load doubleword, 64 bits), **SD** (store doubleword)
- New arithmetic: **ADDW, SUBW, SLLW, SRLW, SRAW** — these operate on the lower 32 bits and sign-extend the result to 64 bits (the "W" means "word" = 32 bits)

The W-suffix instructions are needed because C's `int` type is 32 bits even on 64-bit systems. When you add two `int` values, you want 32-bit wrap-around behavior, not 64-bit. The W instructions give you exactly that:

```riscv
# RV64I: adding two 32-bit integers (C's int type):
addw  t0, t1, t2    # 32-bit add, result sign-extended to 64 bits
                    # Overflow wraps at 2^32, not 2^64

# RV64I: adding two 64-bit integers (C's long or int64_t):
add   t0, t1, t2    # 64-bit add, overflow wraps at 2^64
```

### Quick Check 5.1

> 1. What is the purpose of the `FENCE` instruction in RV32I?
> 2. Why does RISC-V need two instructions (`LUI` + `ADDI`) to load a 32-bit constant?
> 3. In RV64I, what does the "W" suffix on `ADDW` mean, and why is it needed?

---

## 6. Standard Extensions: Building Blocks

One of RISC-V's most powerful ideas is **modularity through extensions**. Each extension adds a coherent set of instructions for a specific purpose. Chip designers choose the extensions their application needs and implement only those.

The standard extensions form an alphabet that describes a chip's capabilities:

### M Extension: Integer Multiply and Divide

**Why is this an extension and not in the base?** A minimal embedded microcontroller — like the kind running a thermostat or a light switch — may have only a few thousand gates. Hardware multiply requires around 1,000 extra gates. For a device that never needs fast multiply, that is pure waste. Software can emulate multiplication using shifts and adds, slowly, but it works.

For any serious computation, however, you need hardware multiply:

```riscv
# M extension instructions:
mul    t0, t1, t2     # t0 = lower 32 bits of (t1 × t2)
mulh   t0, t1, t2     # t0 = upper 32 bits of signed × signed product
mulhu  t0, t1, t2     # t0 = upper 32 bits of unsigned × unsigned product
mulhsu t0, t1, t2     # t0 = upper 32 bits of signed × unsigned product
div    t0, t1, t2     # t0 = t1 / t2  (signed integer division)
divu   t0, t1, t2     # t0 = t1 / t2  (unsigned division)
rem    t0, t1, t2     # t0 = t1 % t2  (signed remainder)
remu   t0, t1, t2     # t0 = t1 % t2  (unsigned remainder)
```

The four multiply variants (mul, mulh, mulhu, mulhsu) handle the full 64-bit product of two 32-bit numbers. Many algorithms — cryptography, big integer math — need the upper half of the product. These instructions provide it without needing a separate 64-bit register.

### A Extension: Atomic Operations

When multiple processor cores run simultaneously (multicore systems), they may try to modify the same memory location at the same time. Without special instructions, this causes **race conditions** — unpredictable results that depend on which core happens to execute first.

**Atomic operations** are instructions that perform a complete read-modify-write cycle without interruption. Think of an atomic operation like a locked box: you open it, take out the contents, put in new contents, and lock it again — all in one uninterruptible action.

```riscv
# Load-Reserved / Store-Conditional: the primitive building block
lr.w   t0, (a0)       # Load word from [a0], mark address as "reserved"
# ... do something ...
sc.w   t1, t2, (a0)   # Try to store t2 to [a0]
                       # Succeeds (t1=0) if reservation still valid
                       # Fails (t1=1) if another core touched [a0]

# Typical lock implementation:
lock_attempt:
  lr.w  t0, (a0)      # Load lock variable
  bnez  t0, lock_attempt  # If already locked (non-zero), retry
  li    t1, 1         # Value to store (locked = 1)
  sc.w  t2, t1, (a0)  # Attempt to set lock
  bnez  t2, lock_attempt  # If store failed, retry
  # Critical section — we have the lock
  ...
  sw    zero, 0(a0)   # Release lock (store 0)

# Atomic fetch-and-add (single instruction):
amoadd.w  t0, t1, (a0)  # t0 = old value of [a0]; [a0] += t1 (atomically)

# Other AMO instructions:
amoand.w  t0, t1, (a0)  # atomic AND
amoor.w   t0, t1, (a0)  # atomic OR
amoxor.w  t0, t1, (a0)  # atomic XOR
amoswap.w t0, t1, (a0)  # swap: t0 = old [a0]; [a0] = t1
amomax.w  t0, t1, (a0)  # atomic signed maximum
amominu.w t0, t1, (a0)  # atomic unsigned minimum
```

### F and D Extensions: Floating-Point

The F extension adds single-precision (32-bit) IEEE 754 floating-point. The D extension adds double-precision (64-bit). These use a separate 32-register floating-point file (f0-f31).

```riscv
# F extension (single-precision):
flw    ft0, 0(a0)     # load float from memory
fsw    ft0, 0(a1)     # store float to memory
fadd.s ft2, ft0, ft1  # ft2 = ft0 + ft1 (single-precision)
fmul.s ft2, ft0, ft1  # multiply
fdiv.s ft2, ft0, ft1  # divide
fsqrt.s ft0, ft1      # square root
fmadd.s ft0, ft1, ft2, ft3  # ft0 = ft1*ft2 + ft3 (fused multiply-add)

# D extension (double-precision):
fld    ft0, 0(a0)     # load double from memory
fadd.d ft2, ft0, ft1  # double-precision add
```

**Fused Multiply-Add (FMADD)** is the key to floating-point performance. It computes `a × b + c` in a single instruction with only one rounding step instead of two. This improves both speed and numerical accuracy. Modern CPUs can issue multiple FMADD instructions per cycle — they are the workhorse of floating-point computation.

### C Extension: Compressed Instructions

Most programs spend most of their time in a small set of common patterns: loading a small constant, adding a small offset, moving between nearby registers, branching a short distance. The C extension adds 16-bit versions of these common patterns:

```
Full 32-bit instructions     Compressed 16-bit equivalents
─────────────────────────    ─────────────────────────────
addi  t0, t0, 1      →      c.addi  t0, 1
lw    t0, 0(sp)      →      c.lwsp  t0, 0
sw    t0, 0(sp)      →      c.swsp  t0, 0
add   a0, a0, a1     →      c.add   a0, a1
beq   a0, zero, L    →      c.beqz  a0, L
jal   x0, L          →      c.j     L
```

The hardware decodes both 16-bit and 32-bit instructions in the same pipeline — the lowest two bits distinguish them (00, 01, 10 = 16-bit; 11 = 32-bit). Programs using the C extension are typically 25-30% smaller, which directly translates to fitting in less flash memory on embedded devices.

### V Extension: Vector Operations

The Vector extension — ratified in 2021 after years of debate — is RISC-V's most innovative extension. It takes a fundamentally different approach from ARM NEON or x86 AVX.

In those ISAs, the vector width is fixed: AVX-512 always operates on 512-bit vectors, NEON always on 128-bit. Code compiled for AVX-512 does not run on a CPU without AVX-512 support.

RISC-V vectors use **variable-length architecture (VLA)**. The hardware implements whatever vector width it wants — 128 bits, 256 bits, 512 bits, 1024 bits. The software asks "how many elements fit?" at runtime and loops accordingly:

```riscv
# Add two arrays of float32 values (a[] + b[] → c[]):
# a0 = pointer to a[], a1 = pointer to b[], a2 = pointer to c[]
# a3 = number of elements

loop:
  vsetvli   t0, a3, e32, m8, ta, ma  # Configure: 32-bit elements, 8 registers per group
                                       # t0 = actual elements processed this iteration
  vle32.v   v0,  (a0)                 # Load from a[] into vector register group v0
  vle32.v   v8,  (a1)                 # Load from b[] into v8
  vfadd.vv  v16, v0, v8              # v16 = v0 + v8 (element-wise float add)
  vse32.v   v16, (a2)                 # Store v16 to c[]
  sub       a3, a3, t0                # Remaining elements -= processed elements
  slli      t1, t0, 2                 # Byte offset = elements × 4 bytes
  add       a0, a0, t1                # Advance a[] pointer
  add       a1, a1, t1                # Advance b[] pointer
  add       a2, a2, t1                # Advance c[] pointer
  bnez      a3, loop                  # Continue if elements remain
```

This same code runs on any RISC-V V-extension chip, regardless of its vector width. A microcontroller with 128-bit vectors might process 4 floats per iteration; a server chip with 1024-bit vectors processes 32 per iteration. The code is identical.

### Extension Naming: Reading the Profile String

RISC-V chips are described by a **profile string** that lists their extensions:

```
RV32I     — 32-bit, integer base only
RV32IM    — 32-bit, integer + multiply
RV32IMAC  — 32-bit, integer + multiply + atomic + compressed
RV64GC    — 64-bit, G (= IMAFD) + compressed
            "G" is a shorthand for IMAFD: all the "general purpose" extensions
RV64GCV   — 64-bit, general + compressed + vector

Common embedded profile:   RV32IMAC  (microcontroller)
Common application profile: RV64GC   (Linux-capable)
```

When you see "this chip is RV64GC," you now know exactly what it can and cannot do.

### Quick Check 6.1

> 1. Why is hardware multiply in the M extension rather than the base ISA?
> 2. What is a race condition and how do atomic operations (A extension) prevent it?
> 3. What is the key advantage of RISC-V's variable-length vector (V) extension over ARM NEON's fixed 128-bit vectors?

---

## 7. Instruction Encoding: Elegance in 32 Bits

RISC-V's instruction encoding is a masterclass in regular design. Every 32-bit instruction falls into one of six formats:

```
R-type  (register-register operations):
 31      25 24  20 19  15 14  12 11   7 6      0
┌─────────┬──────┬──────┬──────┬──────┬────────┐
│ funct7  │  rs2 │  rs1 │funct3│  rd  │ opcode │
│  7 bits │5 bits│5 bits│3 bits│5 bits│ 7 bits │
└─────────┴──────┴──────┴──────┴──────┴────────┘
  Example: ADD rd, rs1, rs2

I-type  (immediate / loads / JALR):
 31                20 19  15 14  12 11   7 6      0
┌───────────────────┬──────┬──────┬──────┬────────┐
│    imm[11:0]      │  rs1 │funct3│  rd  │ opcode │
│    12 bits        │5 bits│3 bits│5 bits│ 7 bits │
└───────────────────┴──────┴──────┴──────┴────────┘
  Example: ADDI rd, rs1, immediate
           LW rd, offset(rs1)

S-type  (stores):
 31      25 24  20 19  15 14  12 11   7 6      0
┌─────────┬──────┬──────┬──────┬──────┬────────┐
│imm[11:5]│  rs2 │  rs1 │funct3│imm[4:0]│opcode│
│  7 bits │5 bits│5 bits│3 bits│5 bits│ 7 bits │
└─────────┴──────┴──────┴──────┴──────┴────────┘
  Example: SW rs2, offset(rs1)
  (offset is split across two fields — keeps rs1/rs2 in same position as R-type)

B-type  (conditional branches):
 31      25 24  20 19  15 14  12 11   7 6      0
┌─────────┬──────┬──────┬──────┬──────┬────────┐
│imm[12|10:5]│rs2│  rs1 │funct3│imm[4:1|11]│op│
└─────────┴──────┴──────┴──────┴──────┴────────┘
  Example: BEQ rs1, rs2, label
  (offset encodes jump distance in multiples of 2 bytes, always even)

U-type  (upper immediate):
 31                12 11   7 6      0
┌───────────────────┬──────┬────────┐
│    imm[31:12]     │  rd  │ opcode │
│    20 bits        │5 bits│ 7 bits │
└───────────────────┴──────┴────────┘
  Example: LUI rd, upper_immediate
           AUIPC rd, upper_immediate

J-type  (unconditional jump):
 31                12 11   7 6      0
┌───────────────────┬──────┬────────┐
│imm[20|10:1|11|19:12]│ rd │ opcode │
│    20 bits (shuffled)│5 bits│7 bits│
└───────────────────┴──────┴────────┘
  Example: JAL rd, label
```

### The Clever Bit Placement

Notice that in R-type, I-type, S-type, and B-type instructions:
- **rd** (destination) is always at bits [11:7]
- **rs1** (first source) is always at bits [19:15]
- **rs2** (second source) is always at bits [24:20]
- **funct3** is always at bits [14:12]
- **opcode** is always at bits [6:0]

This means a hardware decoder can read the register addresses from every instruction **before** fully decoding the opcode. This is critical for **out-of-order execution** pipelines, where the processor needs to know data dependencies between instructions as early as possible.

The B and J immediate bits look scrambled — why are they shuffled? Because the RISC-V designers wanted the sign bit of every immediate to always be at bit 31. Sign-extending an immediate (filling upper bits with the sign bit) is a very common operation, and placing the sign bit consistently at bit 31 means the sign extension hardware never needs instruction type information.

### Quick Check 7.1

> 1. In RISC-V's R-type, I-type, S-type, and B-type formats, where is rs1 always located?
> 2. Why is it useful that register addresses are at the same bit positions across instruction types?
> 3. Why is the sign bit of every RISC-V immediate placed at bit 31?

---

## 8. Privilege Levels: Security Layers

Modern processors must support multiple levels of trust. A user's application program should not be able to directly control hardware, read other programs' memory, or disable the CPU's safety features. The operating system kernel needs more power than applications. A hypervisor (software that runs multiple operating systems simultaneously) needs more power than the kernel.

RISC-V implements this hierarchy through **privilege levels**:

```
Higher privilege (more power, more trust):

  ┌─────────────────────────────────────────────────┐
  │  M-mode: Machine mode                           │
  │  Highest privilege. Direct hardware access.     │
  │  Used by: firmware, bootloader, OpenSBI         │
  │  Can access: all CSRs, all memory               │
  └─────────────────┬───────────────────────────────┘
                    │
  ┌─────────────────▼───────────────────────────────┐
  │  H-mode: Hypervisor mode (optional)             │
  │  Manages multiple virtual machines              │
  │  Used by: KVM, Xen, cloud hypervisors           │
  └─────────────────┬───────────────────────────────┘
                    │
  ┌─────────────────▼───────────────────────────────┐
  │  S-mode: Supervisor mode                        │
  │  Operating system kernel                        │
  │  Used by: Linux kernel, FreeBSD kernel          │
  │  Can access: supervisor CSRs, manage page tables│
  └─────────────────┬───────────────────────────────┘
                    │
  ┌─────────────────▼───────────────────────────────┐
  │  U-mode: User mode                              │
  │  Lowest privilege. Normal applications.         │
  │  Used by: your programs, web browsers, games    │
  │  Cannot: access hardware directly               │
  └─────────────────────────────────────────────────┘

Lower privilege (less power, less trust)
```

A simple embedded system implementing only M-mode is like a single-story building — everything runs in the same space with the same level of access. A server running Linux implements M+S+U: the firmware lives at M, the Linux kernel at S, and your applications at U.

### Control and Status Registers (CSRs)

System state — interrupt enable flags, exception handlers, page table pointers, performance counters — lives in **Control and Status Registers (CSRs)**. Each CSR has a 12-bit address. Access uses special instructions:

```riscv
csrrw  t0, csr, t1    # CSR Read/Write: t0 = old CSR value; CSR = t1
csrrs  t0, csr, t1    # CSR Read/Set: t0 = old CSR; CSR |= t1 (set bits)
csrrc  t0, csr, t1    # CSR Read/Clear: t0 = old CSR; CSR &= ~t1 (clear bits)
csrr   t0, csr        # Pseudoinstruction: just read CSR (csrrs t0, csr, x0)
csrw   csr, t0        # Pseudoinstruction: just write CSR (csrrw x0, csr, t0)
```

Key CSRs in M-mode:

```
mstatus   —  Machine status: global interrupt enable, privilege mode, etc.
mtvec     —  Machine trap-vector: address of exception/interrupt handler
mepc      —  Machine exception PC: the PC that caused the exception
mcause    —  Machine cause: code explaining what caused the trap
mscratch  —  Scratch space for M-mode handler
mhartid   —  Hardware thread ID (0 for first core, 1 for second, etc.)
misa      —  ISA features this hart implements (read to detect extensions)
```

Key CSRs in S-mode (Linux uses these):

```
sstatus   —  Supervisor status (subset of mstatus)
stvec     —  Supervisor trap-vector
sepc      —  Supervisor exception PC
scause    —  Supervisor cause
satp      —  Supervisor address translation and protection
             Contains page table root address and mode (Sv39, Sv48, Sv57)
```

### Traps: How Exceptions and System Calls Work

When something unexpected happens — a system call, a hardware error, an interrupt — the processor takes a **trap**:

```
Normal execution:     user program runs in U-mode

System call (ECALL):
  1. User program loads syscall number into a7
  2. User program loads arguments into a0-a5
  3. User program executes ECALL instruction
  4. Hardware saves current PC to sepc (or mepc)
  5. Hardware sets cause in scause (or mcause): "environment call from U-mode"
  6. Hardware switches to S-mode (or M-mode)
  7. Hardware jumps to stvec (the kernel's trap handler address)
  8. Kernel reads a7 (which syscall?), dispatches to correct handler
  9. Kernel executes syscall, puts return value in a0
  10. Kernel executes SRET instruction
  11. Hardware restores U-mode, jumps to saved PC + 4 (next instruction)
  12. User program continues, sees return value in a0

Page fault (memory access to unmapped address):
  Steps 4-7 same, but scause says "load page fault" or "store page fault"
  Kernel maps the page (or terminates the program if it's invalid)
  Kernel executes SRET to retry the instruction
```

### Virtual Memory: Sv39 and Sv48

S-mode enables virtual memory through the **satp** CSR. When satp is configured, every memory access by U-mode (and S-mode) goes through page table translation.

RISC-V supports multiple virtual address sizes:
- **Sv32**: 32-bit virtual addresses (for RV32), 4 GiB address space, 2-level page table
- **Sv39**: 39-bit virtual addresses (for RV64), 512 GiB address space, 3-level page table
- **Sv48**: 48-bit virtual addresses (for RV64), 256 TiB address space, 4-level page table
- **Sv57**: 57-bit virtual addresses (for RV64), 128 PiB address space, 5-level page table

Linux on RISC-V uses Sv39 by default (the same address space size as Linux on ARM64).

### Quick Check 8.1

> 1. What are the four RISC-V privilege levels, from highest to lowest privilege?
> 2. Which privilege level does the Linux kernel run in?
> 3. What happens to the program counter (PC) when a user program executes ECALL?

---

## 9. Growing Adoption: Who Is Using RISC-V?

RISC-V is no longer a research project or a niche experiment. It is running in products you interact with every day — in your Wi-Fi router, your hard drive, your cloud provider's data center, and possibly your next laptop.

### Western Digital: The Billion-Core Announcement

In 2019, Western Digital — the company that makes the hard drives and SSDs in hundreds of millions of devices — made a startling announcement: they had already deployed **one billion RISC-V cores** and planned to transition all future storage controllers to RISC-V.

Why? Western Digital ships roughly 350 million hard drives and SSDs per year. Each storage device contains one or more microcontroller cores for data management. At a 1% ARM royalty on a $10 controller chip, that is $35 million per year in royalties — just to ARM. Replacing ARM with internally-designed RISC-V cores eliminates that cost entirely.

Western Digital also open-sourced their SweRV core design, contributing it back to the RISC-V ecosystem.

### NVIDIA: RISC-V Inside the GPU

NVIDIA's graphics cards are controlled by dozens of small embedded processors. These handle power management, security, firmware updates, display management, and other control functions. NVIDIA called their original embedded processor the "Falcon" (with a custom ISA).

Starting around 2015, NVIDIA began replacing Falcon processors with RISC-V cores. Today, every modern NVIDIA GPU — the RTX 4090, the A100 data center GPU, the H100 — contains multiple RISC-V cores. When your GPU updates its firmware or manages thermal limits, a RISC-V processor is doing that work.

NVIDIA chose RISC-V because the open ecosystem means better tools: GCC, LLVM, debuggers, simulators — all work without proprietary toolchain licensing.

### SiFive: The ARM of the RISC-V World

**SiFive** was founded in 2015 by the original RISC-V creators — Krste Asanović, Andrew Waterman, and Yunsup Lee — as a commercial vehicle for RISC-V IP. SiFive's model mirrors ARM's: they design RISC-V processor cores and license them.

SiFive products:
- **U74**: Application processor competing with ARM Cortex-A55
- **P550**: High-performance core competing with ARM Cortex-A75, used in Intel Horse Creek development board
- **E31**: Embedded core for microcontrollers
- Development boards: HiFive Unmatched (desktop Linux board), HiFive1 (Arduino-form-factor RISC-V board)

The paradox: the RISC-V ecosystem has attracted commercial IP vendors, some of whom charge for their implementations — but you can always choose a free open-source core like the VexRiscv, CVA6, or Ibex instead.

### Alibaba T-Head: Open-Source RISC-V from China

Alibaba's chip division, **T-Head**, designed the XuanTie C906 — a 64-bit RISC-V processor targeting embedded and AIoT (AI + IoT) applications. In a landmark move, Alibaba open-sourced the complete chip design (RTL — Register-Transfer Level, the hardware description language source code) of the C906.

This means anyone can download, synthesize, and manufacture the C906 design. It has been incorporated into dozens of low-cost RISC-V SoCs sold by Chinese manufacturers, some retailing for under $2.

The **AllWinner D1** chip, containing the T-Head C906, powered a $15 Linux-capable single-board computer (the Nezha board) — demonstrating that a capable 64-bit Linux RISC-V system could reach consumer pricing.

### Google: OpenTitan Security Chip

Google's **OpenTitan** project is building an open-source silicon **Root of Trust** — a tamper-resistant security chip that verifies the integrity of a system's firmware before booting. It uses a RISC-V core (the Ibex core from ETH Zurich).

Root of trust chips sit at the foundation of secure computing: servers, laptops, and phones use them to verify that the software chain from firmware up to the operating system has not been tampered with. Google chose RISC-V because the open design allows external security auditors to verify the implementation — impossible with a proprietary ISA.

### Espressif ESP32-C3: RISC-V in Your Wi-Fi Devices

Espressif Systems — the company that makes the ESP8266 and ESP32 chips used in millions of IoT projects — switched their new chip family to RISC-V:

```
ESP32-C3 (2020):
  Core:       Single RV32IMC @ 160 MHz
  RAM:        400 KB SRAM
  Flash:      Up to 4 MB (external)
  Wireless:   Wi-Fi 802.11n + Bluetooth 5.0
  Price:      $0.50-1.00 in quantity
  
ESP32-C6 (2022):
  Core:       Single HP core RV32IMC @ 160 MHz
              + Low-power RV32IMC @ 20 MHz
  Wireless:   Wi-Fi 6, Bluetooth 5, Zigbee/Thread
  
ESP32-H2 (2023):
  Core:       RV32IMC @ 96 MHz
  Wireless:   Bluetooth 5.3, Zigbee, Thread
  Target:     Smart home mesh networking
```

Every smart light bulb, every home automation sensor, every cheap IoT gadget using these chips contains a RISC-V processor. Hundreds of millions of units have shipped.

### India: SHAKTI at IIT Madras

India's government has invested heavily in domestic chip design to reduce dependence on foreign semiconductor companies. The flagship project is **SHAKTI** at IIT Madras, which has produced multiple RISC-V processor classes:

- **Class-C (Chromite)**: Microcontroller class, RV32IMC, comparable to ARM Cortex-M3
- **Class-E (Elektra)**: Embedded, RV64IMACF
- **Class-I (Iravati)**: Application processor, RV64GC, Linux-capable
- **Class-M (Moushika)**: High-performance, out-of-order execution

SHAKTI successfully taped out on commercial CMOS processes, and a SHAKTI processor booted Linux in 2019. This made India one of the few countries in the world to have designed a domestic processor capable of running a modern operating system.

### Quick Check 9.1

> 1. Why did Western Digital switch from ARM to RISC-V in their storage products?
> 2. What is the ESP32-C3 and where is it commonly used?
> 3. What is the SHAKTI project and which institution leads it?

---

## 10. Geopolitical Significance

### The Chip as a Strategic Asset

Semiconductors are no longer just consumer products — they are strategic national assets. The race to control chip manufacturing and design is as important to 21st-century geopolitics as oil was to the 20th century.

This reality became undeniable in 2020, when the United States government placed sweeping export controls on semiconductor technology to Chinese companies:

- **Huawei** was cut off from Google's Android license (software) and from TSMC's advanced manufacturing (hardware)
- **SMIC** (China's largest chipmaker) was restricted from receiving advanced US semiconductor equipment
- **CXMT** and other Chinese memory makers faced restrictions

The ARM ISA is a British/US technology. ARM Holdings was owned by a Japanese company (SoftBank), which tried to sell it to NVIDIA (American). Any ARM license could potentially be subject to US technology export controls. For Chinese companies trying to build chips, ARM is a vulnerability.

### RISC-V as Geopolitical Insurance

RISC-V is governed by a Swiss non-profit. The ISA specification is in the public domain. There is no single country that can revoke RISC-V access.

This makes RISC-V enormously attractive for:

```
China:
  Goal:        Chip sovereignty, reduce US/UK dependency
  Strategy:    National RISC-V research programs, T-Head open-sourcing,
               dozens of domestic RISC-V startups funded
  Reality:     ARM is licensed from a company that could be cut off;
               RISC-V requires no license from anyone

European Union:
  Goal:        Reduce dependency on US and Asian chip ecosystems
  Strategy:    RISC-V as foundation for European processor programs
               (European Processor Initiative, EPAC project)
  Programs:    HiPEAC, EPI (European Processor Initiative)

India:
  Goal:        Domestic semiconductor capability, SHAKTI program
  Strategy:    RISC-V as foundation for all government chip programs

US:
  Goal:        Domestic chip production, CHIPS Act
  Strategy:    RISC-V supported as it keeps option value open
  Concern:     Chinese RISC-V advancement could be dual-use (defense)
```

### The Open Question: Is RISC-V Truly Free?

There is a subtle complication. The ISA specification is free. But building a working chip requires more than an ISA:

1. **EDA tools** (electronic design automation — software to design chips): Synopsys and Cadence, the dominant EDA vendors, are American companies. Their tools are subject to US export controls. Without EDA tools, you cannot design a chip.

2. **Process technology**: TSMC (Taiwan), Samsung (South Korea), and Intel (US) have the most advanced manufacturing processes. Access to leading-edge fabrication is controlled.

3. **IP blocks**: A chip needs more than a CPU core — it needs memory controllers, USB interfaces, DDR controllers, etc. Many of these IP blocks come from proprietary vendors.

RISC-V removes the ISA licensing barrier. It does not remove all barriers. A country aiming for complete chip sovereignty still needs its own EDA tools (China is developing Empyrean EDA), its own fabs (SMIC), and its own IP libraries.

Still, removing the ISA barrier is significant. RISC-V ensures that no country can be frozen out of processor design at the ISA level.

### Quick Check 10.1

> 1. Why is ARM considered a geopolitical vulnerability for Chinese chip companies?
> 2. Why is RISC-V International headquartered in Switzerland?
> 3. Besides the ISA, what other components of chip design are subject to export controls?

---

## 11. Comparison: RISC-V vs ARM vs x86

### The Three Architectures Side by Side

| Feature | x86-64 | AArch64 (ARMv8/9) | RISC-V (RV64GC) |
|---------|--------|-------------------|-----------------|
| Year introduced | 1978 (x86) / 2003 (x86-64) | 2011 | 2010 |
| Owner | Intel / AMD | ARM Holdings | Public domain |
| License cost | Not available | Millions + royalties | Free |
| Royalty per chip | N/A (only Intel/AMD make x86) | 0.5–2% of chip price | Zero |
| Instruction width | Variable: 1–15 bytes | Fixed: 32 bits | 32 bits (16 with C ext.) |
| General-purpose registers | 16 | 31 (x0 is zero in some contexts) | 32 (x0 hardwired zero) |
| Instruction count (base ISA) | ~1,500+ | ~800+ | 47 |
| SIMD/vector | SSE/AVX/AVX-512 | NEON, SVE | V extension |
| Floating-point | x87, SSE, AVX | VFP, NEON | F+D extensions |
| Atomics | LOCK prefix, CMPXCHG | LDXR/STXR, LSE | A extension (LR/SC, AMO) |
| Compressed instructions | No (but variable length gives density) | No base, Thumb2 in AArch32 | C extension |
| Virtual memory | x86-64 paging (up to 5-level) | ARMv8 paging (up to 4-level) | Sv39/Sv48/Sv57 |
| Privilege levels | Ring 0-3 (+ VMX) | EL0-EL3 | U/S/H/M |
| Formal specification | No | No | Yes (Sail language) |
| Open source cores | No | No (Cortex designs are proprietary) | Many: VexRiscv, CVA6, Ibex, Rocket |
| Custom extensions | No | Allowed by license | Built into spec |
| Legacy compatibility burden | Enormous (40+ years) | Moderate (AArch32 support) | Minimal (2010 clean slate) |
| Code density | Best (variable-width) | Good | Good with C ext. |
| Toolchain maturity | Excellent | Excellent | Good (GCC/LLVM full support) |
| Server market share | ~95% | ~5%, growing fast (AWS Graviton) | <1%, growing |
| Mobile market share | ~0% | ~100% | <1% |
| Embedded market share | ~5% | ~80% | ~15%, fastest growing |
| Notable users | Intel, AMD CPUs in PCs/servers | Apple, Qualcomm, Samsung | WD, NVIDIA, Espressif, SiFive |

### Where RISC-V Wins

```
Best use cases for RISC-V today:

1. EMBEDDED / IoT:
   Advantage: No royalties matter enormously at billions-of-chips scale
   Competitors: ARM Cortex-M (dominant), RISC-V growing fast
   Example: ESP32-C3 ($0.50) vs comparable ARM chip ($0.70+)
   
2. CUSTOM ACCELERATORS:
   Advantage: Custom extension space built into spec
   Examples: ML inference chips, crypto chips, DSPs
   Competitors: None — ARM and x86 don't allow arbitrary custom extensions
   
3. ACADEMIC / RESEARCH:
   Advantage: Free, simple, well-documented, easy to simulate
   Competitors: None — RISC-V has completely replaced MIPS/DLX in curricula
   
4. CHIP-SOVEREIGN NATIONS:
   Advantage: No foreign gatekeeper; full design freedom
   Competitors: None — only RISC-V offers this
   
5. SECURE / SAFETY-CRITICAL:
   Advantage: Formal Sail specification allows mathematical verification
   Competitors: None at this level of rigor
```

### Where RISC-V Is Still Catching Up

```
Areas where ARM and x86 still lead (as of 2024):

1. HIGH-PERFORMANCE COMPUTING:
   Gap: Apple M4, AMD Zen 5, Intel Meteor Lake have decades of microarchitecture
        optimization. Best RISC-V chips are roughly Cortex-A75 class.
   Closing: SiFive P670, Ventana Veyron V2, Sophgo SG2042

2. ECOSYSTEM MATURITY:
   Gap: ARM has 30 years of drivers, BSPs, libraries, certified software
   Closing: Linux/Android support complete; IoT stacks maturing

3. GPU COMPUTE:
   Gap: NVIDIA CUDA, AMD ROCm dominate; no RISC-V GPU compute
   Note: RISC-V is used inside GPUs as management cores, not for shaders

4. SERVER MARKET:
   Gap: x86-64 has >95% market share backed by decades of ecosystem
   Closing: Alibaba, Ventana targeting hyperscale workloads
```

### The Long Game

RISC-V does not need to win every market segment immediately. Its strategy is attrition: start in segments where no-royalty matters most (IoT, embedded, sovereign chips), build ecosystem maturity, and gradually expand upmarket.

The Linux analogy is instructive: in 1995, Linux ran on servers where free software mattered most. In 2010, it ran on phones. In 2025, it runs everywhere. RISC-V is following a similar trajectory, 15 years behind Linux but with the same structural advantages of openness.

### Quick Check 11.1

> 1. How many instructions are in the RISC-V RV32I base ISA compared to x86 and ARM?
> 2. In which market segment is RISC-V growing fastest, and why?
> 3. What advantage does RISC-V's formal Sail specification give it over ARM and x86 in safety-critical applications?

---

## Summary

- **The Problem**: x86 and ARM were proprietary ISAs — expensive to license, politically controlled, inaccessible to researchers and small companies.

- **RISC-V's Birth**: Created at UC Berkeley in 2010 by Krste Asanović, Andrew Waterman, and David Patterson as an open-source ISA — free for anyone to implement, extend, or manufacture chips with.

- **Design Principles**: Clean-slate RISC design; minimal base ISA (47 instructions); modular optional extensions; fixed 32-bit instruction width; ratified extensions frozen forever; reserved opcode space for custom extensions.

- **Registers**: 32 integer registers (x0-x31); x0 hardwired to zero; ABI names divide them into argument (a0-a7), saved (s0-s11), temporary (t0-t6), return address (ra), and stack pointer (sp).

- **Base ISAs**: RV32I for 32-bit systems; RV64I for 64-bit, adding 64-bit registers and W-suffix instructions for 32-bit word operations with sign extension.

- **Standard Extensions**: M (multiply/divide), A (atomics for multicore), F (single-precision float), D (double-precision float), C (16-bit compressed instructions, 25-30% code size savings), V (variable-length vector). Profile string like RV64GC describes a chip's capabilities.

- **Instruction Encoding**: Six formats (R/I/S/B/U/J); register addresses always at the same bit positions; sign bit always at bit 31.

- **Privilege Levels**: U (user), S (supervisor/OS), H (hypervisor), M (machine/firmware). System state accessed via CSRs. ECALL transfers from U-mode to S/M-mode for system calls.

- **Adoption**: Western Digital (1B+ cores in storage), NVIDIA (GPU management), Espressif (ESP32-C3 Wi-Fi chips), SiFive (commercial IP), Alibaba T-Head (XuanTie C906, open-sourced), Google (OpenTitan security), India SHAKTI (IIT Madras).

- **Geopolitics**: RISC-V International in Switzerland; ISA in public domain; no country can revoke access. Critical for Chinese, Indian, and European chip sovereignty efforts. EDA tools and fabs remain gatekeepers even if ISA is free.

- **Comparison**: x86 dominates servers, ARM dominates mobile, RISC-V growing fastest in embedded/IoT and custom accelerators. RISC-V's unique advantages: no royalties, custom extension space, formal Sail specification, and geopolitical freedom.

---

## Exercises

### Easy

1. **Extension identification**: A RISC-V chip is described as "RV32IMAC." What does each letter stand for? What capabilities does this chip have? What capabilities does it lack compared to RV64GC?

2. **Register roles**: List the ABI name, hardware name, and purpose for: the stack pointer, the return address register, the first function argument register, and the always-zero register.

3. **Privilege mapping**: Match each piece of software to its RISC-V privilege level: (a) a web browser, (b) the Linux kernel's scheduler, (c) the system firmware (OpenSBI), (d) a Python script, (e) a virtual machine manager (KVM).

### Medium

4. **Royalty economics**: A company ships 200 million Wi-Fi chips per year at $3.00 each. ARM charges a 1.2% royalty. A RISC-V core implementation requires $8M upfront design cost and $3M/year ongoing toolchain and support costs, versus $1M/year for an ARM ecosystem. Calculate: (a) annual ARM royalty, (b) annual RISC-V break-even cost (amortized over 5 years), (c) the annual net saving with RISC-V, and (d) after how many years does RISC-V become cheaper if volume grows 20% per year?

5. **ISA design tradeoff**: The RISC-V base ISA does not include hardware multiply (M extension) or floating-point (F/D extensions). Describe two concrete scenarios where including these in the base would be wasteful, and two scenarios where their absence is a significant inconvenience. What does this reveal about the philosophy of modular ISA design?

6. **Instruction encoding puzzle**: You are implementing a RISC-V hardware decoder. Given the 32-bit instruction `0x00A50533`, decode it step by step: (a) What are bits [6:0] (opcode)? (b) What are bits [11:7] (rd)? (c) What are bits [19:15] (rs1)? (d) What are bits [24:20] (rs2)? (e) What instruction format is this? (f) Look up opcode 0x33 and funct3/funct7 in the RISC-V spec to identify the exact instruction.

### Hard

7. **Geopolitical scenario analysis**: The year is 2026. The US government extends semiconductor export controls to prohibit licensing of ISAs to entities on a restricted list — including a major Asian technology company. (a) How does this affect that company's ability to use ARM? (b) How does this affect their ability to use RISC-V? (c) What parts of chip production would still be affected by export controls even if they use RISC-V? (d) Design a "chip sovereignty roadmap" for a hypothetical nation trying to achieve domestic semiconductor capability using RISC-V. What does it need to build or acquire, and in what order?

8. **Custom extension design**: You are designing a RISC-V chip for a cryptocurrency application (not mining — think hardware wallet key management and signing). The performance bottleneck is SHA-256 computation and elliptic curve point multiplication over secp256k1. (a) Design two custom RISC-V instructions — one for SHA-256 round function and one for 256-bit modular multiplication — specifying the instruction format, operands, and semantics. (b) What opcode space would you use? (c) Estimate the speedup versus software on RV32IM. (d) What are the security risks of adding hardware cryptographic instructions, and how might an attacker exploit timing side-channels in your implementation?

9. **Variable-length vector deep dive**: Write out the complete inner loop for matrix-vector multiplication `y = A × x` where A is an N×N matrix of float32 values and x, y are N-element float32 vectors, using RISC-V V-extension instructions. The loop should work correctly on any RISC-V V implementation regardless of hardware vector width. Annotate each instruction explaining what it does and why. Then analyze: what is the arithmetic intensity (floating-point operations per memory byte loaded) of this kernel, and at what hardware vector width does it become memory-bandwidth-bound rather than compute-bound?
