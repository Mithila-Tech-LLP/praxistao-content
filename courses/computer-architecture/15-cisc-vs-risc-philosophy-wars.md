# Chapter 15: CISC vs RISC — The Philosophy Wars

For decades, computer architects have debated two fundamentally different approaches to instruction set design: CISC (Complex Instruction Set Computing) and RISC (Reduced Instruction Set Computing). Understanding this debate is central to understanding why modern processors are designed the way they are, and why the winner of the debate is... complicated.

## Table of Contents

1. [The CISC Philosophy](#1-the-cisc-philosophy)
2. [The RISC Philosophy](#2-the-risc-philosophy)
3. [The Semantic Gap and Why It Matters](#3-the-semantic-gap)
4. [RISC Principles in Detail](#4-risc-principles-in-detail)
5. [The Academic Debate: Evidence and Counterevidence](#5-the-academic-debate-evidence-and-counterevidence)
6. [x86: CISC on the Outside, RISC on the Inside](#6-x86-cisc-on-the-outside-risc-on-the-inside)
7. [ARM: RISC Growing Up](#7-arm-risc-growing-up)
8. [RISC-V: RISC Pure from the Start](#8-risc-v-risc-pure-from-the-start)
9. [Who Won?](#9-who-won)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The CISC Philosophy

### The World Before RISC (1960s-1970s)

In the early era of computing, memory was expensive, slow, and scarce. A memory access could take 100-200 nanoseconds — enormously slow compared to the CPU. Programmers wrote in assembly language by hand. Every byte of code mattered. Programs that were shorter ran faster simply because fewer memory accesses were needed to fetch instructions.

In this environment, the dominant philosophy was: **pack as much work as possible into each instruction.** This minimizes:
1. Program size (fewer instructions needed to express the same computation)
2. Memory accesses for instruction fetching
3. The work programmers had to do writing assembly

The result was CISC: processors with large, complex instruction sets capable of sophisticated operations in a single instruction.

### Classic CISC Instructions

Consider VAX (Digital Equipment Corporation, 1977) — the canonical CISC architecture:

```
POLYD src, degree, tbladdr
```
This single instruction evaluates a polynomial: it computes `a0 + a1*x + a2*x² + ... + an*xⁿ` using Horner's method by reading n+1 coefficients from a table in memory. What a compiler needs 20+ simpler instructions to do, VAX does in one.

Intel x86 examples of complex instructions:
- `ENTER n, 0` — sets up a stack frame (adjusts RSP, saves RBP, links frame)
- `LOOP target` — decrements RCX and branches if non-zero
- `REPNE SCASB` — scan a string for a byte value (implements `strlen`-like functionality)
- `MUL [0x1234]` — multiply AX by the 16-bit value at memory address 0x1234

These instructions implement in hardware what a compiler would generate as loops or multi-instruction sequences.

### The CISC Logic

1. **Semantic density**: Each instruction does more work → programs are shorter
2. **Assembly programmer productivity**: Easier to write by hand
3. **Microcode**: Complex instructions are implemented in microcode (a second layer of simple internal instructions) — the hardware doesn't have to be truly complex at the gate level
4. **Memory bandwidth**: Shorter programs put less pressure on the instruction memory bus

### Quick Check

> 1. What two scarce resources drove the CISC philosophy?
> 2. What does the VAX POLYD instruction do? Why is this remarkable?
> 3. What is "microcode" and why did it enable complex instructions?

---

## 2. The RISC Philosophy

By the late 1970s, the hardware landscape had changed:

- **Compilers** were replacing assembly language programmers. Software was increasingly compiled, not hand-written in assembly.
- **Memory was getting cheaper and faster**, reducing the penalty of larger programs.
- **Research showed** that programs mostly used only a tiny fraction of the CISC instruction set.

In this environment, researchers at Berkeley and Stanford developed a radical alternative.

### The IBM 801 and the 1980 Studies

John Cocke at IBM Research (around 1974-1975) first demonstrated that a simple, fast processor could outperform complex processors. His IBM 801 project used a tiny instruction set — all instructions executed in one cycle — and achieved remarkable performance.

At Berkeley in 1980, David Patterson and Carlo Sequin coined the term RISC for their Berkeley RISC project. At Stanford, John Hennessy built MIPS (Microprocessor without Interlocked Pipeline Stages).

### The RISC Principle

RISC starts from a completely different premise: **make the common case fast.**

Research showed that in real programs:
- Simple instructions (ADD, LOAD, STORE, branch) account for 80%+ of executed instructions
- Complex instructions (string operations, polynomial evaluation) are almost never used
- Compilers don't generate complex instructions reliably — they prefer simpler ones they can schedule freely

If complex instructions slow down simple ones (because the processor pipeline must accommodate them), you're paying a cost on 80% of instructions to help 1%.

### The RISC Manifesto (Patterson and Sequin, 1980)

1. **Small, fixed-length instructions**: every instruction is 32 bits
2. **Load-store architecture**: only LOAD and STORE access memory; all ALU operations use registers
3. **Large, uniform register file**: many registers (32) of equal size and capability
4. **Single-cycle execution**: all instructions complete in one clock cycle (or close to it)
5. **Hardwired control**: no microcode — decode and control logic is purely combinational
6. **Compiler-friendly**: the ISA is designed for compilers, not assembly programmers

### Quick Check

> 1. What changed in the hardware landscape between the CISC era (1960s) and RISC era (1980)?
> 2. What fraction of instructions in real programs are simple operations (ADD, LOAD, STORE)?
> 3. List the six principles of RISC design.

---

## 3. The Semantic Gap

The central argument for CISC was the **semantic gap** — the distance between what high-level languages express and what simple hardware instructions can do.

A single line of C:
```c
for (int i = 0; i < n; i++) sum += array[i];
```

A CISC architecture can express this as just a few instructions (perhaps a MOVEM + loop instruction). A RISC architecture needs maybe 10 instructions.

But here's the RISC counterargument:

### The "RISC Advantage" Study (Patterson, 1985)

Patterson examined what code compilers actually generated for CISC vs RISC. He found:

| Factor | CISC Expectation | Reality |
|--------|-----------------|---------|
| Instructions per program | RISC has more | True: 2-4× more instructions |
| Cycles per instruction | Same | False: CISC CPI is 3-10×; RISC CPI ≈ 1 |
| Time per program = IPC × CPI | CISC wins | RISC wins overall by 2-3× |

The key insight: **CPI matters more than instruction count.**

CISC complex instructions have high CPI (cycles per instruction) because:
- They take many clock cycles to execute
- They make pipelining nearly impossible
- They require complex decode logic that slows all instructions

RISC sacrifices instruction count but achieves CPI ≈ 1 (and later < 1 with superscalar execution), winning on overall execution time.

### Quick Check

> 1. What is the "semantic gap" in the CISC vs RISC debate?
> 2. Which metric matters more for performance: instructions per program, or cycles per instruction?
> 3. Why do complex CISC instructions have high CPI?

---

## 4. RISC Principles in Detail

### Load-Store Architecture

RISC processors can ONLY access memory through dedicated LOAD and STORE instructions. ALL arithmetic/logic operations work exclusively on registers.

```
CISC (x86) style:
ADD [0x1234], EAX        ; one instruction: reads memory, adds, writes memory

RISC (RISC-V) style:
LW   t0, 0(a0)           ; load from address in a0 into t0
ADD  t0, t0, a1          ; add a1 to t0
SW   t0, 0(a0)           ; store t0 back to address in a0
```

Why this matters for hardware: the pipeline has dedicated stages (IF, ID, EX, MEM, WB). If any instruction can be a memory access, the MEM stage must be active for ALL instructions — wasting hardware. With load-store, only LW/SW use the MEM stage; ALU instructions skip it (or pass through with no memory operation).

### Large Register File

Early RISC machines (Berkeley RISC-I, 1982) had 138 registers using a clever **register window** scheme (see Chapter 9). Modern RISC-V has 32 × 64-bit integer registers and 32 × 64-bit floating-point registers. ARM64 has 31 × 64-bit general-purpose registers + 32 × 128-bit vector registers.

More registers → less register spilling to memory → fewer LOAD/STORE instructions → faster execution.

### Single-Cycle Execution and Pipelining

RISC instructions are designed so that every instruction completes in exactly one clock cycle. This makes them trivially pipelineable: you can overlap 5 instructions in the 5-stage pipeline (IF, ID, EX, MEM, WB) with no hardware complexity to handle variable-latency instructions.

With CISC variable-latency instructions, the pipeline must stall or use complex out-of-order logic to handle instructions that take 1 cycle vs. 30 cycles.

### Hardwired Control Unit

CISC processors used microprogrammed control (a ROM of microcode sequences). RISC processors use hardwired control: combinational logic that directly maps opcode bits to control signals. This is:
- Faster (no ROM lookup latency)
- Simpler (less hardware)
- More predictable (same latency every time)

### Compiler-Friendly Design

RISC ISAs are explicitly designed for compilers, not human assembly programmers:
- Orthogonal instruction set (every instruction works with every register)
- No special-purpose "hidden" registers that instructions modify implicitly
- Explicit branch delay slots (early RISC) or compiler-friendly no-delay branches (modern RISC-V)
- Calling convention specified in the ISA ABI, not enforced by hardware

### Quick Check

> 1. Why is load-store architecture important for pipelining?
> 2. How does having more registers reduce memory traffic?
> 3. What is the difference between microprogrammed and hardwired control? Which is faster?

---

## 5. The Academic Debate: Evidence and Counterevidence

### RISC Won the Academic Argument

By the late 1980s, academic consensus was clear: RISC was superior. The MIPS R2000 (1988) and SPARC chips demonstrated that RISC processors were 2-4× faster than equivalent CISC machines at the same clock speed.

### But x86 CISC Dominated the Market

Despite losing the technical debate, x86 dominated personal computing through the 1990s and 2000s. Why?

1. **IBM PC momentum**: the IBM PC (1981) standardized on x86. By the time RISC clearly won technically (late 1980s), there were tens of millions of x86 programs and billions in software investment.
2. **Microsoft DOS/Windows**: tied to x86 binary compatibility
3. **Intel's manufacturing**: Intel had the best chip fabs in the world and poured billions into making x86 faster
4. **Intel's secret weapon**: translate CISC to RISC internally (see next section)

### The Convergence

By the mid-1990s, both camps had adopted the best ideas from each other:
- x86 CPUs translated CISC instructions to internal RISC-like micro-operations at runtime
- RISC CPUs added more complex instructions (SIMD, cryptography) to close the semantic gap for specialized workloads

The boundary blurred. Today's debate is less "CISC vs RISC" and more "which ISA legacy carries less baggage?"

---

## 6. x86: CISC on the Outside, RISC on the Inside

The biggest secret in computer architecture: since the Intel Pentium Pro (1995), x86 processors have NOT been executing CISC instructions internally.

### Micro-Operation (µop) Translation

Modern Intel and AMD processors have a **front-end** that fetches and decodes x86 instructions, translating each CISC instruction into one or more RISC-like **micro-operations** (µops):

```
CISC instruction (from x86 program):
ADD [0x1234], EAX           ; add EAX to memory at 0x1234

Decoded into 3 µops (what the execution engine sees):
µop1: LOAD  t0, [0x1234]    ; load memory into temp register
µop2: ADD   t0, t0, EAX     ; add
µop3: STORE t0, [0x1234]    ; store result back
```

The µops are RISC-like: fixed-size, single operation, register-to-register. After translation, the execution engine is essentially a RISC out-of-order processor that happens to have a CISC decoder in front.

### Consequences

1. x86 programs run on what is effectively RISC hardware
2. The CISC complexity is absorbed in the front-end decoder
3. The execution core can be highly optimized for RISC-style µops
4. Software binary compatibility with 45 years of x86 programs is preserved

This is why Intel and AMD can squeeze so much performance from x86: the backend is a state-of-the-art RISC engine, and only the frontend maintains CISC compatibility.

```
What the programmer sees:         What the hardware actually does:
─────────────────────────────     ────────────────────────────────────────
x86-64 CISC instructions          µop-based superscalar RISC engine
• Variable length (1-15 bytes)    • Fixed-size µops
• Complex addressing modes        • Register-register operations
• Implicit register modifications • Explicit operands
• Memory operands in ALU instrs  • Load-store only
```

### Quick Check

> 1. What are micro-operations (µops) and when were they introduced to x86?
> 2. Why does the µop translation allow x86 processors to have RISC-like performance?
> 3. What does the decoder stage do in a modern x86 processor?

---

## 7. ARM: RISC Growing Up

ARM (Advanced RISC Machines, 1985) started as a pure RISC design for the BBC Micro computer — tiny team, tiny budget, elegant 26-bit instruction set.

### ARM's Journey

| Year | Development |
|------|-------------|
| 1985 | ARM1: 25,000 transistors, 8 MHz, 26-bit ISA |
| 1990 | ARM6: first widely-used embedded ARM |
| 1994 | ARM7TDMI: Game Boy Advance processor (2001) |
| 2001 | ARMv6: added SIMD (Jazelle/NEON precursor) |
| 2011 | ARMv8/AArch64: 64-bit extension (modern ARM64) |
| 2016 | ARMv8.2+: SVE (scalable vector extension) |
| 2020 | Apple M1: ARM64 dominant in laptops |

### ARMv7 vs ARMv8 (AArch64)

ARMv7 (32-bit ARM) accumulated complexity over 25 years: 16 registers, conditional execution on every instruction, Thumb/Thumb-2 variable-length encoding, inline barrel shifter, strange calling conventions. It became almost as complex as x86 in practice.

ARMv8 (AArch64, 2011) was a clean-break redesign:
- 31 × 64-bit general-purpose registers (x0-x30)
- 32 × 128-bit vector registers (v0-v31)
- Fixed 32-bit instruction encoding
- No conditional execution (except branch instructions)
- Clean calling convention

AArch64 is a much cleaner RISC ISA than AArch32 (32-bit ARM) — closer in spirit to RISC-V.

### Thumb Instruction Set

For embedded systems where code size matters, ARM added **Thumb**: a 16-bit encoding for the most common ARM instructions, achieving ~70% of 32-bit code density (i.e., 30% smaller programs). The processor switches between 16-bit Thumb mode and 32-bit ARM mode using a bit in the processor status register.

This is the RISC concession to CISC: adding variable-width encoding to reduce code size for constrained environments.

### Quick Check

> 1. How many transistors did the original ARM1 have compared to modern processors?
> 2. What is ARM Thumb and why was it added?
> 3. What is the main difference between ARMv7 and ARMv8 (AArch64) in terms of register count?

---

## 8. RISC-V: RISC Pure from the Start

RISC-V (2010, UC Berkeley, Krste Asanović's group) is the newest major RISC architecture and the most intellectually pure RISC design ever deployed at scale.

### Why Another RISC ISA?

By 2010, the RISC landscape had problems:
- **MIPS**: owned by MIPS Technologies, licensed expensively, declining
- **ARM**: licensed per chip by ARM Ltd., expensive royalties
- **SPARC**: locked to Sun/Oracle ecosystem, declining
- **PowerPC**: expensive, IBM-controlled

Academia and industry needed an open, free, modern RISC ISA for research and commercial use without legal or licensing encumbrance.

### RISC-V Design Principles

1. **Open and free**: no royalties, no licensing fees, full specification publicly available
2. **Modular**: base ISA + optional extensions
3. **Simple and clean**: no legacy baggage
4. **Scalable**: RV32 (32-bit), RV64 (64-bit), RV128 (128-bit, future)

### The Base ISA + Extensions Model

```
RV32I (base): 47 instructions — integer arithmetic, loads, stores, branches, JAL, JALR
              This alone is Turing-complete.

Standard extensions:
  M: integer multiplication and division (MUL, DIV, REM)
  A: atomic instructions (for multi-core synchronization)
  F: single-precision floating-point (IEEE 754 32-bit)
  D: double-precision floating-point (IEEE 754 64-bit)
  C: compressed 16-bit instructions (code size reduction, like ARM Thumb)
  V: vector operations (SIMD)
  H: hypervisor support
  B: bit manipulation
  K: cryptography

A common embedded profile: RV32IMAC
A common 64-bit application profile: RV64GC (where G = IMAFD)
```

### Why RISC-V Matters

- **India's SHAKTI processor** is RISC-V (see Volume 8)
- **Western Digital** builds its own RISC-V processor for storage controllers (saving on ARM royalties)
- **NVIDIA** uses RISC-V cores for microcontrollers in its GPUs
- **ESP32-C3** (popular Wi-Fi microcontroller) uses RISC-V
- **SiFive** (startup), Andes, StarFive produce commercial RISC-V chips
- **China's RISC-V push**: many Chinese chip companies use RISC-V to avoid geopolitical risks of ARM licensing

RISC-V is the most important ISA development since ARM64.

### Quick Check

> 1. Why was RISC-V created rather than using an existing RISC ISA?
> 2. What does the "G" in RV64GC stand for?
> 3. Name three companies using RISC-V in commercial products.

---

## 9. Who Won?

The CISC vs RISC debate has no clean winner. Instead:

### In Servers and Desktops: x86 and ARM
- x86 dominates servers (AMD EPYC, Intel Xeon) and desktops
- ARM is growing rapidly in servers (AWS Graviton, Ampere Altra) and dominates laptops (Apple M-series)

### In Mobile: ARM Won Decisively
Every major smartphone — iPhone, Android flagships — runs ARM. The power efficiency advantage of RISC's clean ISA (plus decades of ARM optimization for low power) made x86 uncompetitive in mobile. Intel tried repeatedly (Atom for phones, Medfield) and failed every time.

### In Embedded: ARM and RISC-V Dominate
The vast majority of embedded microcontrollers (IoT, appliances, cars, medical devices) run ARM Cortex-M or increasingly RISC-V. x86 is almost absent.

### In AI/ML: Specialized ISAs
Neither CISC nor RISC in the traditional sense — AI accelerators use custom ISAs optimized for matrix operations (systolic arrays, dot-product instructions, etc.).

### The Deeper Answer

The CISC vs RISC battle is largely over. The winning philosophy is:
- **Simple, regular, load-store ISA** (RISC principles)
- **Binary compatibility with existing software** (CISC's market moat for x86)
- **Complex optional extensions** for specific workloads (SIMD, crypto, AI) added to otherwise RISC-clean ISAs

Pure CISC died in the 1990s. Pure RISC evolved to accommodate the real world. Today's ISA design is pragmatic RISC with targeted complexity.

---

## Summary

- **CISC** was born from expensive memory: pack maximum work per instruction. Classic examples: VAX, x86, Motorola 68k.
- **RISC** was born from compiler research: simple, regular instructions with CPI≈1 beat complex instructions with high CPI. Classic examples: MIPS, SPARC, early ARM, RISC-V.
- **The key insight**: total execution time = instruction_count × cycles_per_instruction. RISC sacrifices instruction count to achieve CPI≈1 (vs CISC CPI=5-10). RISC wins.
- **x86 resolved the paradox** by translating CISC instructions to RISC-like µops internally (since Pentium Pro, 1995). CISC outside, RISC inside.
- **ARM** started pure RISC but added Thumb (code size), NEON (SIMD), SVE (AI), growing in complexity.
- **RISC-V** is the newest ISA: open-source, modular (base + extensions), royalty-free, rapidly growing.
- Modern ISA design is pragmatic: RISC principles for the core, targeted complex extensions for specific workloads.

---

## Exercises

### Easy

1. What are the three main principles that CISC architects prioritized in the 1960s-1970s?

2. List the six RISC principles (from the 1980 RISC manifesto). For each, explain WHY it contributes to performance.

3. What is "µop translation" in modern x86 processors? Why does it matter?

### Medium

4. **Counting cycles**: A program has these instruction mix statistics (from profiling):
   - 40% ADD/SUB (CPI = 1 for RISC, 1 for CISC)
   - 30% LOAD/STORE (CPI = 1 for RISC, 2 for CISC)
   - 20% multiply (CPI = 3 for RISC, 8 for CISC)
   - 10% complex string/memory ops (CPI = 5 for RISC, 10 for CISC)
   
   Also, the RISC processor needs 1.5× more instructions than CISC to express the same program. Compare total cycles for RISC vs CISC. Which wins?

5. ARM Thumb uses 16-bit encoding for a subset of ARM instructions and achieves ~70% code density. What is the tradeoff? If a cache line holds 64 bytes, how many instructions can fit in a cache line with 32-bit ARM vs 16-bit Thumb? Why does this matter?

6. RISC-V uses a modular ISA. The base RV32I has 47 instructions. The "C" (compressed) extension adds 16-bit instruction variants. The "V" vector extension adds hundreds of vector instructions. Does this modular growth make RISC-V a de facto CISC architecture over time? Argue both sides.

### Hard

7. **Semantic gap and compiler quality**: The CISC argument rests on the "semantic gap" — CISC instructions are closer to what high-level languages express. But in 1980, Patterson observed that compilers rarely used complex CISC instructions. Investigate: why would a compiler NOT use a hardware POLYD (polynomial evaluation) instruction? What would prevent a C compiler from recognizing a polynomial loop and emitting POLYD? What does this tell us about the assumption that hardware complexity reduces the semantic gap?

8. **The µop translation overhead**: Intel's Pentium Pro (1995) introduced the decode front-end that translates x86 to µops. This decoder adds: (a) latency (1-5 cycles to decode complex instructions), (b) die area (the decoder is ~10-15% of total chip area), and (c) power consumption. Yet the P6 microarchitecture was faster than all RISC competitors at the time. How? What design choices allowed Intel to win with this approach? (Consider: clock frequency, branch prediction, out-of-order execution, L2 cache, manufacturing process.) This research exercise shows that architectural decisions don't happen in isolation — system-level factors often dominate.
