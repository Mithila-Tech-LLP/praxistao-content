# Chapter 07: What Is a CPU? Anatomy of a Processor

You have built the pieces: transistors become gates, gates become adders, flip-flops become registers. Now it is time to assemble them into the most important circuit ever built — the **Central Processing Unit**. This chapter gives you the complete picture of what a CPU is, what every component does, and how they connect into a working machine.

---

## Table of Contents

1. [The Kitchen Analogy — A Mental Model for the Whole CPU](#1-the-kitchen-analogy)
2. [What a CPU Actually Does](#2-what-a-cpu-actually-does)
3. [The ALU — The Chef Doing the Work](#3-the-alu--the-chef-doing-the-work)
4. [Registers — The Counter Space](#4-registers--the-counter-space)
5. [The Control Unit — The Recipe](#5-the-control-unit--the-recipe)
6. [The Program Counter and Instruction Register](#6-the-program-counter-and-instruction-register)
7. [The Clock — The Kitchen Timer](#7-the-clock--the-kitchen-timer)
8. [Machine Code — What Instructions Actually Look Like](#8-machine-code--what-instructions-actually-look-like)
9. [The Datapath — From Memory to ALU to Registers](#9-the-datapath--from-memory-to-alu-to-registers)
10. [IPC — Instructions Per Cycle](#10-ipc--instructions-per-cycle)
11. [The Big Picture — How All Components Interact](#11-the-big-picture--how-all-components-interact)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. The Kitchen Analogy

Before diving into transistors and binary, let us build a mental model that will stay with you through the entire book.

Imagine a **professional kitchen** during a busy dinner service. This kitchen is your CPU.

```
┌─────────────────────────────────────────────────────────────┐
│                    THE CPU KITCHEN                           │
│                                                              │
│  ┌──────────────┐    ┌─────────────────┐                    │
│  │   PANTRY     │    │  COUNTER SPACE  │    ┌────────────┐  │
│  │   (Memory /  │    │  (Registers)    │    │            │  │
│  │    RAM)      │    │                 │    │   CHEF     │  │
│  │              │    │  bowl1  bowl2   │◄──►│   (ALU)    │  │
│  │  Flour       │◄──►│  bowl3  bowl4   │    │            │  │
│  │  Eggs        │    │                 │    └────────────┘  │
│  │  Sugar       │    └─────────────────┘                    │
│  │  Butter      │                                           │
│  └──────────────┘    ┌─────────────────┐    ┌────────────┐  │
│                      │    RECIPE       │    │   TIMER    │  │
│                      │  (Control Unit) │    │  (Clock)   │  │
│                      │                 │    │            │  │
│                      │  Step 1: ...    │    │ tick tick  │  │
│                      │  Step 2: ...    │    │            │  │
│                      │  Step 3: ...    │    └────────────┘  │
│                      └─────────────────┘                    │
└─────────────────────────────────────────────────────────────┘
```

Here is the mapping:

| Kitchen Component | CPU Component | Role |
|-------------------|---------------|------|
| Chef | ALU (Arithmetic Logic Unit) | Does the actual work — chopping, mixing, cooking |
| Counter space | Registers | The small work surface right in front of the chef |
| Pantry | Memory (RAM) | Where all the ingredients (data) are stored |
| Recipe | Control Unit | The step-by-step instructions telling the chef what to do |
| Kitchen timer | Clock | The rhythmic signal that paces every action |
| Recipe book index | Program Counter | Tracks which step of the recipe you're on |
| Current recipe card | Instruction Register | The specific step currently being followed |

This analogy will snap everything into place as we explore each component. Keep it in mind.

### Why the Analogy Works

A good chef does not invent new techniques while cooking. They follow the recipe step by step: "add 200g flour, mix for 30 seconds, preheat oven to 180°C." The CPU is the same. It does not invent or improvise — it follows instructions mechanically, one at a time. The magic comes from *which* instructions the programmer wrote, not from any cleverness in the CPU itself.

The **counter space** (registers) is tiny — maybe four or five bowls can fit. If the chef needs more ingredients, they must fetch them from the pantry (memory) and put one of the current bowls away first. This back-and-forth between counter and pantry is one of the central performance challenges in CPU design.

The **timer** (clock) rings at regular intervals. The chef does not begin the next step until the timer rings. This keeps everything synchronized — nobody grabs an ingredient before the oven has preheated.

### Quick Check

> 1. In the kitchen analogy, what does the pantry represent in a real CPU?
> 2. Why is the counter space (registers) kept small rather than made as large as the pantry?
> 3. What would happen if the kitchen had no timer — why does the CPU need a clock?

---

## 2. What a CPU Actually Does

A CPU **executes programs**. A program is a sequence of instructions — "add these two numbers," "compare this value to zero," "store this result in memory," "jump to a different part of the program if the condition is true."

The CPU executes these instructions one by one (or, in modern CPUs, several at once), billions of times per second.

### A Profound Observation

The CPU does not understand what a program *means* at a high level. It does not know whether it is sorting a list, playing music, or drawing a 3D scene. It only executes the next instruction.

This is one of the most beautiful facts in all of computing: **intelligence emerges from the combination of dumb operations, not from any single smart one.** A single instruction is like a single brushstroke — meaningless in isolation, a masterpiece in combination.

### The CPU's Three Core Jobs

Every CPU, from the simplest microcontroller to the most powerful server chip, performs the same three steps for every instruction:

```
┌─────────────────────────────────────────────────────────────┐
│                                                              │
│    1. FETCH      2. DECODE      3. EXECUTE                  │
│                                                              │
│   Read the       Figure out      Carry out                  │
│   next           what the        the                        │
│   instruction    instruction     operation                  │
│   from memory    means                                      │
│                                                              │
│   (Go to the     (Read the       (Chef does                 │
│   right page     recipe step)    the step)                  │
│   in the book)                                              │
│                                                              │
└─────────────────────────────────────────────────────────────┘
```

These three steps — repeated billions of times per second — are the heartbeat of every program ever run. We will return to them in much greater detail in Chapter 11.

### Quick Check

> 1. What does it mean for a CPU to "execute a program"?
> 2. Name the three steps a CPU performs for every instruction.
> 3. If a CPU does not understand the meaning of a program, why does the program still do something useful?

---

## 3. The ALU — The Chef Doing the Work

The **Arithmetic Logic Unit** is the component that actually performs computations. Everything else in the CPU exists to feed data to the ALU and take results away from it.

Think of the ALU as the chef: the most skilled person in the kitchen, but entirely dependent on the support staff. The control unit tells the chef what to cook. The registers (counter space) hold the ingredients right at hand. Memory (pantry) holds the rest.

### What the ALU Receives and Produces

```
                Operand A                Operand B
                (ingredient 1)          (ingredient 2)
                    │                        │
                    ▼                        ▼
              ┌─────────────────────────────────┐
              │                                 │◄── Operation Code
              │            A L U                │    (which recipe step?)
              │                                 │
              └──────────────────┬──────────────┘
                                 │
                    ┌────────────┼────────────┐
                    ▼            ▼            ▼
                  Result      Zero         Other
                (n bits)      Flag         Flags
                                        (Carry, Neg, Overflow)
```

The ALU takes two input values (operands) and an operation code, and produces:
- A **result** — the computed value
- **Status flags** — condition bits that record properties of the result

### ALU Operations

A typical simple ALU handles these categories:

```
┌─────────────────────────────────────────────────────────────┐
│                    ALU OPERATIONS                            │
├──────────────────┬──────────────────┬───────────────────────┤
│   ARITHMETIC     │   LOGIC          │   SHIFTS              │
│                  │                  │                       │
│ ADD   A + B      │ AND   A & B      │ SHL   A << n          │
│ SUB   A - B      │ OR    A | B      │ SHR   A >> n          │
│ MUL   A × B      │ XOR   A ^ B      │ SAR   A >>> n         │
│ DIV   A ÷ B      │ NOT   ~A         │   (arithmetic right)  │
│ NEG   -A         │                  │                       │
└──────────────────┴──────────────────┴───────────────────────┘
```

### Status Flags — Why They Matter

After every ALU operation, a special register called the **status register** (or **flags register**) is updated with four key bits:

| Flag | Name | Meaning | Use |
|------|------|---------|-----|
| Z | Zero | Result is exactly zero | Testing if two values are equal |
| N | Negative | Result is negative (MSB = 1) | Testing if a value is less than zero |
| C | Carry | Unsigned result overflowed the word size | Multi-word arithmetic |
| V | Overflow | Signed result overflowed (result too large to fit) | Detecting signed arithmetic errors |

These flags are how the CPU implements `if` statements. When you write `if (a == b)` in C, the compiler generates a **compare** instruction that subtracts b from a and sets the flags, then a **branch** instruction that checks the Zero flag. If Z = 1, the branch is taken.

```
Without flags, branching would be impossible.
Without branching, loops and conditionals would be impossible.
Without loops and conditionals, general-purpose computing would be impossible.

Flags are tiny — just 4 bits — but they are the foundation of all conditional logic.
```

### Quick Check

> 1. What are the two inputs to an ALU, and what does it produce as output?
> 2. What does the Zero flag indicate, and how is it used in practice?
> 3. Why would an ALU that only did arithmetic (no logic operations) be less useful?

---

## 4. Registers — The Counter Space

**Registers** are the fastest storage in the entire computer. They sit inside the CPU chip itself, right next to the ALU. Reading from or writing to a register takes one clock cycle — or even less in some designs. There is zero travel time, because there is zero distance.

But they are precious and few. A typical CPU has between 8 and 32 general-purpose registers. Each holds one "word" — the CPU's native data width (32 bits on a 32-bit CPU, 64 bits on a 64-bit CPU).

### Why So Few?

Remember the kitchen analogy: the counter space in front of the chef is small. You can only have a few bowls out at once. If you could have an infinite counter, you would never need to go to the pantry. But counter space is expensive real estate in a commercial kitchen.

The same is true for registers:

- A single 64-bit register, built from D flip-flops, requires about 500–600 transistors (including the multi-port access circuitry)
- A 32-register file with 2 read ports and 1 write port requires tens of thousands of transistors
- Each transistor must be placed on the chip right next to the ALU — premium silicon space
- Doubling the number of registers roughly doubles this cost *and* can slow down the register file (larger = farther away = slower)

Compare that to RAM: a single DRAM cell needs just 1 transistor + 1 capacitor. 8 GB of RAM is 64 billion cells.

```
┌──────────────────────────────────────────────────────────────┐
│              MEMORY HIERARCHY — SPEED vs. SIZE               │
│                                                              │
│  REGISTERS      ◄── fastest ──►   ~32 locations (256 bytes) │
│     ↕            (sub-nanosecond)                           │
│  L1 CACHE                          ~64 KB                   │
│     ↕                                                       │
│  L2 CACHE                          ~512 KB                  │
│     ↕                                                       │
│  L3 CACHE                          ~16 MB                   │
│     ↕                                                       │
│  RAM (DRAM)     ◄── slowest ──►   ~8–32 GB                  │
│                   (nanoseconds)                              │
└──────────────────────────────────────────────────────────────┘
```

### General-Purpose Registers

Different ISAs name their registers differently, but the concept is universal:

```
RISC-V:    x0,  x1,  x2  ... x31   (x0 is hardwired to zero — always reads as 0)
ARM64:     x0,  x1,  x2  ... x30,  SP, PC
x86-64:    rax, rbx, rcx, rdx, rsi, rdi, rsp, rbp, r8 ... r15
MIPS:      $0 ($zero), $1 ($at), $2-$3 ($v0-$v1) ... $31 ($ra)
```

### Special-Purpose Registers

Beyond general-purpose registers, every CPU has special registers with fixed roles:

| Register | Full Name | Purpose |
|----------|-----------|---------|
| PC | Program Counter | Address of the next instruction to fetch |
| IR | Instruction Register | Holds the instruction currently being decoded/executed |
| SP | Stack Pointer | Tracks the top of the call stack in memory |
| FLAGS | Status / Flags Register | ALU condition codes (Z, N, C, V) |
| MAR | Memory Address Register | Address sent to memory during load/store |
| MDR | Memory Data Register | Data read from or about to be written to memory |

### Quick Check

> 1. Why do registers have lower latency (faster access time) than RAM?
> 2. In RISC-V, register x0 always contains the value 0. Why is a hardwired-zero register useful to have?
> 3. If a CPU had 512 registers instead of 32, what advantages and disadvantages might that create?

---

## 5. The Control Unit — The Recipe

The **Control Unit** is the CPU's recipe. It reads each instruction, figures out what it means, and generates the precise set of control signals that make the ALU, registers, and memory do the right thing at the right time.

Critically: **the control unit does no computation itself.** It is pure orchestration. It is the chef's instructions, not the chef's hands.

### Control Signals

The control unit generates binary signals that flow to every other component. Think of them as switches:

```
┌─────────────────────────────────────────────────────────────┐
│                  CONTROL SIGNALS                             │
│                                                              │
│  RegWrite  ─── 1 = write result to register                 │
│                0 = don't write                              │
│                                                              │
│  MemRead   ─── 1 = read from data memory                    │
│  MemWrite  ─── 1 = write to data memory                     │
│                                                              │
│  ALUOp     ─── 2-bit code: which operation?                 │
│            ─── 00 = ADD,  01 = SUB,  10 = AND,  11 = OR    │
│                                                              │
│  ALUSrc    ─── 0 = second ALU operand comes from register   │
│            ─── 1 = second operand is an immediate value     │
│                 (a constant embedded in the instruction)     │
│                                                              │
│  MemToReg  ─── 0 = write ALU result back to register        │
│            ─── 1 = write data loaded from memory to register│
│                                                              │
│  Branch    ─── 1 = this instruction might change the PC     │
│  Jump      ─── 1 = unconditionally change the PC            │
└─────────────────────────────────────────────────────────────┘
```

### Two Types of Control Units

**1. Hardwired Control Unit**
The opcode bits feed directly into combinational logic — a decoder circuit that maps each opcode to a fixed set of control signal values. This is fast (pure logic, no memory lookup) but inflexible. Changing an instruction's behavior requires redesigning the circuit.

**2. Microprogrammed Control Unit**
Each machine instruction maps to a sequence of simpler *micro-operations* stored in a tiny internal ROM called the **control store**. The opcode acts as an index into this ROM. This is slower (a ROM lookup takes time) but extremely flexible — you can add or modify instructions by changing the ROM contents. Many x86 processors use this approach internally.

```
HARDWIRED:
  Opcode[5:0] ──► Decoder Circuit ──► Control Signals
  (pure logic, one gate delay)

MICROPROGRAMMED:
  Opcode[5:0] ──► Micro-ROM Lookup ──► Microcode Sequence ──► Control Signals
  (ROM latency, but very flexible)
```

### Quick Check

> 1. Does the control unit perform arithmetic? What does it actually do?
> 2. What is the difference between a hardwired and a microprogrammed control unit?
> 3. If the ALUSrc signal is set to 1, what does that mean about where the ALU gets its second operand?

---

## 6. The Program Counter and Instruction Register

These two special registers manage the flow of instruction execution. They are so fundamental that it is worth studying them individually.

### The Program Counter (PC) — Your Page Number

The **Program Counter** holds the memory address of the next instruction to fetch. It is the CPU's bookmark in the program.

After fetching each instruction, the PC is automatically incremented to point to the next one. On a RISC architecture where every instruction is exactly 4 bytes (32 bits), the PC increments by 4:

```
Initial state:  PC = 0x1000

Step 1:  Fetch instruction at address 0x1000
         PC ← 0x1000 + 4 = 0x1004

Step 2:  Fetch instruction at address 0x1004
         PC ← 0x1004 + 4 = 0x1008

Step 3:  Fetch instruction at address 0x1008
         PC ← 0x1008 + 4 = 0x100C

... and so on, walking through memory one instruction at a time.
```

For **branch instructions** — the hardware implementation of `if` and `while` — the PC is set to a completely different address instead of incrementing:

```
Normal increment:
  PC = 0x1050 → fetch → PC = 0x1054

Branch taken (jump to address 0x2000):
  PC = 0x1050 → execute BEQ instruction → PC = 0x2000
  PC = 0x2000 → fetch → PC = 0x2004
```

This ability to set the PC to any address is what makes loops possible. A `for` loop is just a branch instruction that sets the PC back to the top of the loop body.

### The Instruction Register (IR) — The Current Recipe Card

After the CPU fetches an instruction from memory, the raw bytes are loaded into the **Instruction Register**. The control unit reads the IR to decode the instruction and generate the right control signals.

The IR holds the instruction for the duration of its execution — like keeping the recipe card in front of you while you follow the step.

```
Memory at 0x1000: 0000 0000 0101 0000 0000 0010 1001 0011
                                    │
                           (CPU fetches it)
                                    │
                                    ▼
Instruction Register (IR): 0000 0000 0101 0000 0000 0010 1001 0011
                                    │
                           (Control Unit reads it)
                                    │
                                    ▼
                           Opcode = 0010011 (ADDI)
                           rs1    = 00000   (register x0)
                           imm    = 000000000101 (immediate = 5)
                           rd     = 00101   (destination x5)
```

### Quick Check

> 1. In a 32-bit RISC CPU where all instructions are 4 bytes, by how much does the PC increment after fetching each instruction?
> 2. How does the CPU implement a `while` loop using the Program Counter?
> 3. What happens to the Instruction Register after the instruction finishes executing?

---

## 7. The Clock — The Kitchen Timer

The **clock** is a crystal oscillator — a piece of quartz that vibrates at a precise, stable frequency when voltage is applied. Each vibration produces a rising and falling edge on a wire that connects to almost every component in the CPU.

```
Clock signal:
    ┌──┐  ┌──┐  ┌──┐  ┌──┐  ┌──┐  ┌──┐
    │  │  │  │  │  │  │  │  │  │  │  │
────┘  └──┘  └──┘  └──┘  └──┘  └──┘  └──

    ↑     ↑     ↑     ↑     ↑     ↑
  rising edges — the CPU "ticks" forward
```

Every flip-flop in the CPU (in the registers, the pipeline stages, the memory buffers) updates on the rising edge. This ensures that all state changes happen simultaneously and in lock-step — nobody reads a value before it has been computed.

### Clock Frequency and What It Means

| Frequency | Period | Meaning |
|-----------|--------|---------|
| 1 MHz | 1,000 ns | 1 million clock edges per second |
| 1 GHz | 1 ns | 1 billion clock edges per second |
| 4 GHz | 0.25 ns | 4 billion clock edges per second |
| 5.8 GHz (modern top end) | ~0.17 ns | 5.8 billion edges per second |

Light travels about 5 cm in 0.17 ns. Modern CPUs are so fast that the speed of light itself becomes an engineering constraint — signals cannot travel far in the time available.

### Clock ≠ Performance (an important nuance)

A common misconception: "a 4 GHz CPU is faster than a 3 GHz CPU." Not necessarily! What matters is how much *work* is done per clock cycle, not just how many cycles happen per second. We will explore this more in Section 10.

### Quick Check

> 1. What physical component generates the clock signal in a CPU?
> 2. If a CPU runs at 3 GHz, how long is each clock period in nanoseconds?
> 3. Why does every flip-flop in the CPU need to be connected to the same clock signal?

---

## 8. Machine Code — What Instructions Actually Look Like

This is where abstract concepts become concrete reality. Every instruction the CPU executes is encoded as a binary number — a specific bit pattern stored in memory.

### Binary Instruction Format

A 32-bit RISC-V instruction is divided into fields. Each field has a specific location and width in the 32-bit word:

```
RISC-V R-Type Instruction (Register-Register operations like ADD):

 31      25  24    20  19   15  14  12  11    7  6       0
┌──────────┬─────────┬────────┬───────┬────────┬─────────┐
│  funct7  │   rs2   │  rs1   │funct3 │   rd   │ opcode  │
│  7 bits  │  5 bits │ 5 bits │3 bits │ 5 bits │  7 bits │
└──────────┴─────────┴────────┴───────┴────────┴─────────┘

funct7 : modifier to the operation (e.g., ADD vs SUB share an opcode)
rs2    : source register 2 (second operand)
rs1    : source register 1 (first operand)
funct3 : further specifies the operation
rd     : destination register (where result goes)
opcode : broad category of instruction (R-type, I-type, Load, Store, Branch...)
```

Let us decode a real instruction. The instruction `ADD x7, x5, x6` (add register x5 and x6, put result in x7):

```
Field    Value    Meaning
─────────────────────────────────────────────────────────────
funct7   0000000  ADD (not SUB)
rs2      00110    Register x6 (binary 6 = 00110)
rs1      00101    Register x5 (binary 5 = 00101)
funct3   000      ADD/SUB (not AND, OR, XOR...)
rd       00111    Register x7 (binary 7 = 00111)
opcode   0110011  R-type ALU operation

Full 32-bit binary:
0000000  00110  00101  000  00111  0110011
│─────│  │───│  │───│  │─│  │───│  │─────│
funct7   rs2    rs1   fn3   rd    opcode

Hexadecimal: 0x006283B3
```

### I-Type: Instructions with Immediate Values

Some instructions use a constant value embedded directly in the instruction, rather than reading from a register. These are called **I-type** (immediate-type):

```
RISC-V I-Type (e.g., ADDI — add immediate):

 31              20  19   15  14  12  11    7  6       0
┌──────────────────┬────────┬───────┬────────┬─────────┐
│   immediate      │  rs1   │funct3 │   rd   │ opcode  │
│     12 bits      │ 5 bits │3 bits │ 5 bits │  7 bits │
└──────────────────┴────────┴───────┴────────┴─────────┘

immediate : a constant value from -2048 to +2047
rs1       : source register
rd        : destination register

Example: ADDI x5, x0, 5  (load constant 5 into register x5)
immediate = 000000000101 (= 5 in binary)
rs1       = 00000 (register x0, which is always 0)
funct3    = 000 (ADDI)
rd        = 00101 (register x5)
opcode    = 0010011 (I-type ALU)
```

### Why Is This Encoding Important?

The instruction encoding is the contract between software and hardware. The compiler that translates your C code into machine code must produce exactly this binary format. The CPU's control unit is designed to decode exactly this format. If either side gets it wrong, the CPU executes the wrong operation.

This encoding is called the **Instruction Set Architecture (ISA)** — one of the most important concepts in all of computing. We explore it fully in Chapter 14.

### A Complete Example: 5 + 3

Here is the three-instruction sequence to compute 5 + 3 and hold the result in a register:

```
Assembly:          ADDI x5, x0, 5     ADDI x6, x0, 3     ADD x7, x5, x6

Hex:               0x00500293         0x00300313         0x006283B3

Binary:
Instruction 1:  0000 0000 0101  00000  000  00101  0010011
Instruction 2:  0000 0000 0011  00000  000  00110  0010011
Instruction 3:  0000000  00110  00101  000  00111  0110011

In memory (bytes, address increasing left to right):
0x1000: 93 02 50 00  (instruction 1, stored little-endian)
0x1004: 13 03 30 00  (instruction 2)
0x1008: B3 83 62 00  (instruction 3)
```

The CPU sees nothing but these bytes. It reads them, decodes the field patterns, and executes the operations. No concept of "addition" or "the number 5" — just bit patterns decoded by logic circuits.

### Quick Check

> 1. In a RISC-V R-type instruction, what does the `opcode` field identify?
> 2. What is the difference between an R-type and an I-type instruction?
> 3. Why does a 5-bit register field limit RISC-V to exactly 32 registers?

---

## 9. The Datapath — From Memory to ALU to Registers

The **datapath** is the collection of wires, multiplexers, registers, and the ALU that data flows through when an instruction executes. It is the physical path that data travels from memory to the CPU and back.

### A Single Instruction's Journey

Let us trace the path of `ADD x7, x5, x6` through a simple CPU:

```
STEP 1 — FETCH
══════════════
PC (contains 0x1008) ─────────────────────────► Memory
                                                  │
                                                  │ returns 32-bit instruction word
                                                  ▼
                                        Instruction Register (IR)
                                        IR ← 0x006283B3


STEP 2 — DECODE
═══════════════
IR (0x006283B3) ──────────────────────────────► Control Unit
                                                  │
  Control Unit decodes:                           │
  - opcode = R-type ALU                           │
  - funct3 = ADD                                  │
  - rs1 = x5, rs2 = x6, rd = x7                  │
                                                  │ generates control signals
                                                  ▼
                                         RegWrite  = 1
                                         ALUOp     = ADD
                                         ALUSrc    = 0 (register, not immediate)
                                         MemRead   = 0
                                         MemWrite  = 0
                                         MemToReg  = 0 (ALU result, not memory)


STEP 3 — EXECUTE
════════════════
Register File:
  Read x5 (contains value 5) ───────────────────► ALU Input A
  Read x6 (contains value 3) ───────────────────► ALU Input B
                                                    │
                                           ALU performs ADD
                                                    │
                                                    ▼
                                             Result = 8
                                             Flags: Z=0, N=0, C=0, V=0
                                                    │
                                                    │ (RegWrite = 1, rd = x7)
                                                    ▼
                                           Write 8 to Register x7


STEP 4 — INCREMENT PC
══════════════════════
PC ← PC + 4 = 0x100C   (ready for next instruction)
```

### The Datapath Diagram

```
           ┌────────────────────────────────────────────────────────┐
           │                     DATAPATH                           │
           │                                                        │
 ┌───┐     │  ┌─────────────┐   ┌──────────────────────────┐      │
 │ PC│────►│  │ Instruction │   │      Register File        │      │
 └───┘     │  │   Memory    │   │                          │      │
   ↑        │  │             │   │  32 registers × 64 bits  │      │
   │        │  └──────┬──────┘   │                          │      │
   │        │         │          │  Read Port 1  Read Port 2│      │
   │        │         ▼          │      ↓             ↓     │      │
   │        │  ┌─────────────┐   └──────┼─────────────┼─────┘      │
   │        │  │ Instruction │          │             │            │
  PC+4      │  │  Register   │   ┌──────┼─────────────┼──────┐     │
   │        │  │    (IR)     │   │      ▼     MUX     ▼      │     │
   │        │  └──────┬──────┘   │  ┌───────────────────┐   │     │
   │        │         │          │  │        ALU         │   │     │
   │        │         ▼          │  │                    │   │     │
   │        │  ┌─────────────┐   │  └────────┬──────────┘   │     │
   │        │  │   Control   │   │           │               │     │
   │        │  │    Unit     │   │       Result              │     │
   │        │  └──────┬──────┘   │           │               │     │
   │        │         │          │           ▼  MUX          │     │
   │        │   Control│Signals  │    ┌──────────────┐        │     │
   │        │         ├──────────┼───►│ Write Port   │        │     │
   │        │         │          │    └──────────────┘        │     │
   │        │         │          └────────────────────────────┘     │
   │        │         │                                             │
   │        │         ▼                                             │
   │        │  ┌─────────────────────────────────────────┐         │
   └────────┤  │              Data Memory                 │         │
(branch)    │  │         (for load/store ops)             │         │
            │  └─────────────────────────────────────────┘         │
            └────────────────────────────────────────────────────────┘
```

The **MUX** (multiplexer) boxes are critical — they are the switches controlled by the control unit. The ALUSrc MUX selects whether the ALU's second operand comes from a register or from an immediate value in the instruction. The MemToReg MUX selects whether the register file is written from the ALU output or from data read from memory.

### Quick Check

> 1. What is the "datapath"?
> 2. In a Load instruction (read from memory into a register), which signal determines that the register is written with the memory data rather than the ALU result?
> 3. After an ADD instruction completes, what two things happen to the program counter?

---

## 10. IPC — Instructions Per Cycle

**IPC** stands for **Instructions Per Cycle** — the number of instructions a CPU completes on average per clock cycle.

This is one of the most important metrics in CPU performance, and understanding it separates superficial understanding from deep understanding.

### The Performance Equation

A CPU's raw performance on any task depends on three factors:

```
                Total Instructions
Time = ─────────────────────────────────
         Clock Frequency × IPC

Or equivalently:
Time = (Total Instructions) × (Cycles Per Instruction) × (Clock Period)
       └──────────────────────┘   └──────────────────────┘   └──────────────┘
              from the program            IPC inverted         1 / frequency
```

This means:
- A higher clock frequency → fewer seconds per cycle → faster execution
- A higher IPC → fewer cycles per instruction → faster execution
- Better compiled code (fewer total instructions) → faster execution

All three factors matter. A CPU at 5 GHz with IPC = 1 does the same work per second as a CPU at 2.5 GHz with IPC = 2.

### Why IPC Matters More Than Clock Speed

In the early days of computing (1970s-1990s), CPUs completed one instruction per cycle at best — IPC = 1. Engineers competed by increasing clock frequencies.

Modern CPUs achieve IPC values of 4–6 for typical workloads through techniques like:

| Technique | What it does | IPC benefit |
|-----------|-------------|-------------|
| Pipelining | Overlaps stages of multiple instructions | ~1× (baseline) |
| Superscalar | Executes multiple instructions per cycle | 2–6× |
| Out-of-order execution | Reorders instructions to avoid waiting | 1.5–3× |
| Branch prediction | Guesses the next path to avoid stalls | 1.1–2× |
| Instruction-level parallelism | Finds independent instructions to run in parallel | varies |

We will cover all of these techniques in later chapters. For now, remember this:

**Clock frequency alone tells you very little. IPC tells you how efficiently the CPU uses each cycle. The product of the two tells you actual throughput.**

### A Simple Example

| CPU | Clock | IPC | Work/second |
|-----|-------|-----|-------------|
| A | 4 GHz | 1.0 | 4 billion instr/sec |
| B | 3 GHz | 2.0 | 6 billion instr/sec |
| C | 2 GHz | 3.5 | 7 billion instr/sec |

CPU C, with the lowest clock speed, does the most work per second — purely because of its higher IPC.

### Quick Check

> 1. Write out the performance equation relating time, total instructions, IPC, and clock frequency.
> 2. A CPU runs at 2 GHz with an IPC of 3. How many instructions does it complete per second?
> 3. Why did early CPU designers focus on increasing clock frequency rather than IPC?

---

## 11. The Big Picture — How All Components Interact

Now let us pull everything together into a single unified view. This is the moment when the kitchen analogy, the binary encodings, the datapath, and the IPC all connect.

### The Complete Interaction Diagram

```
┌─────────────────────────────────────────────────────────────────────────┐
│                           THE COMPLETE CPU                               │
│                                                                          │
│   ┌──────────────────────────────────────────────────────────────────┐  │
│   │                        CLOCK DOMAIN                              │  │
│   │        (Every component below pulses in sync with the clock)     │  │
│   └──────────────────────────────────────────────────────────────────┘  │
│                               ↓ clock                                    │
│                                                                          │
│   ┌──────┐  address   ┌────────────────┐   32-bit instruction            │
│   │  PC  │───────────►│  Instruction   │──────────────────────────┐     │
│   │      │            │    Memory      │                          │     │
│   └──┬───┘            └────────────────┘                          ▼     │
│      │ +4 (or                                               ┌──────────┐ │
│      │  branch                                              │    IR    │ │
│      │  target)                                             │ (Instr.  │ │
│      │                                                      │ Register)│ │
│      │                                                      └─────┬────┘ │
│      │                                                            │      │
│      │                                                            ▼      │
│      │                                                  ┌──────────────┐ │
│      │                                                  │   Control    │ │
│      │                                                  │    Unit      │ │
│      │                                                  │              │ │
│      │                                                  │ Reads opcode │ │
│      │                                                  │ Generates    │ │
│      │                                                  │ control wires│ │
│      │                                                  └──────┬───────┘ │
│      │                                                         │         │
│      │              Control Signals ──────────────────────────┤         │
│      │                ↙          ↘                            │         │
│      │    ┌──────────────────┐  ┌──────────────────────┐     │         │
│      │    │   Register File  │  │    Data Memory       │     │         │
│      │    │                  │  │    (for loads/stores) │     │         │
│      │    │ x0  x1  x2  ...  │  │                      │     │         │
│      │    │ x30 x31          │  └──────────┬───────────┘     │         │
│      │    └───────┬──────────┘             │                  │         │
│      │            │ Operand A              │                  │         │
│      │            │ Operand B              │                  │         │
│      │            ▼                        │                  │         │
│      │    ┌──────────────────┐             │                  │         │
│      │    │       ALU        │             │                  │         │
│      │    │                  │             │                  │         │
│      │    │  + - & | ^ << >> │             │                  │         │
│      │    │                  │             │                  │         │
│      │    └───────┬──────────┘             │                  │         │
│      │            │ Result                 │                  │         │
│      │            │ Flags (Z,N,C,V)        │                  │         │
│      │            ▼                        │                  │         │
│      │       ┌────────┐                    │                  │         │
│      │       │  MUX   │◄───────────────────┘                  │         │
│      │       └────┬───┘  (ALU result OR memory data?)         │         │
│      │            │                                            │         │
│      │            ▼                                            │         │
│      │    ┌──────────────────┐                                 │         │
│      │    │  Write Back to   │                                 │         │
│      │    │  Register File   │                                 │         │
│      │    └──────────────────┘                                 │         │
│      │                                                          │         │
│      └──────────────────────────────────────────────────────────┘         │
│                     (PC updated; cycle repeats)                            │
└────────────────────────────────────────────────────────────────────────────┘
```

### Walking Through a Full Execution Cycle

Let us trace the execution of `ADD x7, x5, x6` from start to finish:

```
TICK 1 — Fetch
─────────────────────────────────────────────────────────
  • PC = 0x1008
  • CPU drives address 0x1008 onto the address bus
  • Instruction memory returns 0x006283B3
  • Instruction Register ← 0x006283B3
  • PC ← 0x1008 + 4 = 0x100C (prepared for next cycle)

TICK 1 — Decode (happens simultaneously with fetch in simple designs)
─────────────────────────────────────────────────────────
  • Control Unit reads IR: 0x006283B3
  • Extracts opcode[6:0]  = 0110011 → R-type
  • Extracts funct3[14:12] = 000
  • Extracts funct7[31:25] = 0000000 → ADD (not SUB)
  • Extracts rs1[19:15]    = 00101 → x5
  • Extracts rs2[24:20]    = 00110 → x6
  • Extracts rd[11:7]      = 00111 → x7
  • Sets: RegWrite=1, ALUOp=ADD, ALUSrc=0, MemRead=0,
           MemWrite=0, MemToReg=0

TICK 1 — Execute
─────────────────────────────────────────────────────────
  • Register File reads x5 → value 5 → ALU Input A
  • Register File reads x6 → value 3 → ALU Input B
  • ALU computes 5 + 3 = 8
  • Flags: Z=0, N=0, C=0, V=0

TICK 1 — Write Back
─────────────────────────────────────────────────────────
  • MemToReg=0 → MUX selects ALU result (not memory data)
  • RegWrite=1 → Register File writes 8 to x7
  • Cycle complete!
```

### The Beauty of It All

Step back and appreciate what just happened. In less than one nanosecond:

1. A number was read from memory (the instruction encoding)
2. That number was decoded by combinational logic into a set of control signals
3. Two register values were read simultaneously
4. They were combined by the ALU — a circuit of a few thousand gates
5. The result was written back to a register
6. The program counter advanced to the next instruction

Then it happened again. And again. 4 billion times per second, reliably, for years without error.

No single component is intelligent. The transistor does not know it is part of an adder. The adder does not know it is part of an ALU. The ALU does not know it is part of a CPU. The CPU does not know it is computing a game physics simulation. But the combination — the hierarchy of emergent behavior — produces something that feels like intelligence.

This emergent complexity from simple, mechanical rules is one of the most profound ideas in all of science.

### Quick Check

> 1. In the execution cycle described above, what happens during "write back"?
> 2. For a Load instruction (not ADD), the MemToReg signal would be 1. What would be different in the execution?
> 3. What would a Branch instruction need to do differently to the PC compared to an ADD instruction?

---

## Summary

| Component | Kitchen Analogy | Role |
|-----------|-----------------|------|
| ALU | Chef | Performs arithmetic, logic, and comparison operations |
| Registers | Counter space | Fast, on-chip storage for current working values |
| Control Unit | Recipe | Decodes instructions and generates control signals |
| Program Counter | Page number | Tracks the address of the next instruction |
| Instruction Register | Current recipe card | Holds the instruction currently being executed |
| Clock | Kitchen timer | Synchronizes all components; drives the pace of execution |
| Memory | Pantry | Stores all program data and instructions |

Key takeaways from this chapter:

- A CPU executes programs by repeatedly **fetching, decoding, and executing** instructions — the three fundamental steps.
- The **ALU** does the computation; the **control unit** tells it what computation to do.
- **Registers** are fast but few; **memory** is large but slow. Managing the gap between them is a central challenge of CPU design.
- Machine code is a **binary encoding** where each instruction is divided into fields: opcode, register identifiers, and sometimes an immediate constant.
- The **datapath** is the physical path data travels: from memory to the instruction register, through the control unit for decoding, from register file to ALU, and back to the register file.
- **IPC (Instructions Per Cycle)** — not just clock frequency — determines how much real work a CPU does per second. Modern CPUs achieve IPC > 4 through techniques like pipelining and out-of-order execution.
- The emergent behavior of billions of transistors executing simple binary instructions produces everything from video games to artificial intelligence.

---

## Exercises

### Easy

1. Name all five core components of a CPU and provide a one-sentence description of each. Use the kitchen analogy in your answer.

2. A RISC-V CPU runs at 3.2 GHz with an average IPC of 2.5. How many instructions does it complete per second?

3. Decode this RISC-V R-type instruction into its fields. The 32-bit value is:
   ```
   0000000  00011  00010  000  00001  0110011
   ```
   What are the values of funct7, rs2, rs1, funct3, rd, and opcode? What operation does this perform?

### Medium

4. Trace through a complete execution cycle for the instruction `ADDI x3, x0, 42` (load the constant 42 into register x3). List what happens in each stage: Fetch, Decode, Execute, Write Back. Include what values are read from and written to each component.

5. **The IPC trade-off**: A CPU designer can either:
   - Option A: Increase clock frequency from 3 GHz to 4 GHz, but this reduces IPC from 3.0 to 2.0 (higher frequency means less time per cycle for complex operations)
   - Option B: Keep 3 GHz, add out-of-order execution to raise IPC from 3.0 to 4.5, but this increases chip area by 40%
   
   Calculate the instructions per second for both options. Which is faster? What non-performance factors might influence the choice?

6. The RISC-V immediate field in an I-type instruction is 12 bits. The value is sign-extended to 32 or 64 bits before use.
   - What is the range of values representable in a 12-bit signed field?
   - If a program needs to load the constant 5000 into a register (which exceeds this range), what instruction sequence would be needed?

### Hard

7. **Design your own instruction encoding**: Design a 16-bit instruction format (not 32-bit) for a minimal CPU with 8 registers. Your ISA needs to support at minimum: register-register ADD/SUB, load from memory, store to memory, and branch if equal. Show the bit field layout for each instruction type. What compromises did you have to make compared to 32-bit RISC-V?

8. **The critical path and clock speed**: The maximum clock frequency a CPU can run at is determined by its **critical path** — the longest chain of logic gates that data must pass through in one cycle.
   
   In a simple single-cycle CPU, the critical path passes through: instruction memory → register file read → ALU → data memory → register file write.
   
   - Explain why the critical path determines the maximum clock frequency.
   - The ALU's critical path is an 64-bit ripple carry adder. If each full adder stage has a propagation delay of 0.1 ns, what is the ALU's contribution to the critical path for a 64-bit addition?
   - Pipelining (Chapter 23) breaks the single-cycle into multiple shorter stages. If the critical path above takes 5 ns total and you split it into 5 equal pipeline stages of 1 ns each, what is the maximum clock frequency of the pipelined CPU? What is the trade-off?
