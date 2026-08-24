# Chapter 10: The Control Unit — The Conductor

Imagine a great symphony orchestra. Dozens of musicians sit ready — violinists, cellists, horn players, percussionists. Each one is expert at their instrument. But without someone to coordinate them, they would play at random times, in random orders, creating noise instead of music. The **conductor** stands at the front, reads the score, and with precise gestures tells each section exactly when to play, how loud, and for how long.

The CPU has its own conductor: the **Control Unit**. The ALU can add and subtract. The registers can store values. The memory can hold programs. But without the Control Unit reading each instruction and issuing precise signals to every component, the CPU would be chaos. This chapter explores how the Control Unit works — how it decodes instructions, generates the signals that make everything happen in the right order, and why its design is one of the most fascinating challenges in computer architecture.

---

## Table of Contents

1. [The Control Unit's Role — What the Conductor Does](#1-the-control-units-role)
2. [Decoding Instructions — Reading the Score](#2-decoding-instructions--reading-the-score)
3. [Control Signals — The Conductor's Gestures](#3-control-signals--the-conductors-gestures)
4. [The Instruction Cycle — Fetch, Decode, Execute](#4-the-instruction-cycle--fetch-decode-execute)
5. [Timing Diagram — Watching the Signals Unfold](#5-timing-diagram--watching-the-signals-unfold)
6. [Different Instructions, Different Signals](#6-different-instructions-different-signals)
7. [Hardwired Control — The Reflex Conductor](#7-hardwired-control--the-reflex-conductor)
8. [Microprogrammed Control — The Sheet Music Conductor](#8-microprogrammed-control--the-sheet-music-conductor)
9. [Hardwired vs Microprogrammed — A Comparison](#9-hardwired-vs-microprogrammed--a-comparison)
10. [Modern CPUs — The Hybrid Approach](#10-modern-cpus--the-hybrid-approach)
11. [Summary](#summary)
12. [Quick Reference](#quick-reference)
13. [Exercises](#exercises)

---

## 1. The Control Unit's Role

The Control Unit (CU) is the brain of the CPU's brain. While the CPU as a whole processes data, the Control Unit is the part that tells every other part of the CPU what to do, moment by moment.

Here is the key insight: **the Control Unit does not compute anything**. It does not add numbers. It does not hold data. It does not make decisions about what program to run. Its sole job is to look at the current instruction and answer one question:

> "Given this instruction, what should every single component in the CPU do right now?"

Then it answers that question by generating **control signals** — electrical lines that are either high (1) or low (0) — that configure the CPU datapath for each instruction.

### The Orchestra Analogy in Detail

Let us extend the orchestra analogy to understand every part:

| Orchestra | CPU |
|-----------|-----|
| Conductor | Control Unit |
| Musical score | The program (instructions in memory) |
| Current measure being played | Current instruction (in the Instruction Register) |
| Individual musicians | ALU, registers, memory interface, MUXes |
| Conductor's gestures | Control signals |
| Beat of the baton | Clock signal |
| Rests (silence) | NOP instructions or pipeline stalls |

The conductor does not play any instrument. But the music cannot happen without the conductor. Similarly, the Control Unit does no computation, but every computation depends entirely on it.

### What the Control Unit Receives

```
Inputs to the Control Unit:
┌─────────────────────────────────────────────────────┐
│                                                     │
│  Instruction Register ──► opcode bits               │
│  (current instruction)     register field bits      │
│                            function code bits       │
│                                                     │
│  Processor Status ──────► privilege mode            │
│                            interrupt enable flag    │
│                            condition codes (N, Z, C, V) │
│                                                     │
│  Clock Signal ──────────► timing reference          │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### What the Control Unit Produces

```
Outputs from the Control Unit (Control Signals):
┌─────────────────────────────────────────────────────┐
│                                                     │
│  ──► ALU operation select (which operation to run)  │
│  ──► Register file read enables                     │
│  ──► Register file write enable                     │
│  ──► Memory read enable                             │
│  ──► Memory write enable                            │
│  ──► MUX select signals (which input to choose)     │
│  ──► PC update signal (next PC value choice)        │
│  ──► Sign extension control                         │
│                                                     │
└─────────────────────────────────────────────────────┘
```

### Quick Check 1

1. The Control Unit does not compute anything. What is its actual job?
2. What does the Control Unit receive as inputs? Name three types.
3. What is a control signal?

---

## 2. Decoding Instructions — Reading the Score

Before the Control Unit can generate signals, it must understand what instruction it is looking at. This process is called **instruction decoding**.

An instruction arrives as a binary number — just a sequence of 0s and 1s. The Control Unit must extract meaning from the bit pattern. This is like a conductor who receives sheet music as a sequence of symbols and must interpret them into gestures for the orchestra.

### The Instruction Format

Every instruction set architecture (ISA) defines an **instruction format** — a fixed layout specifying which bits carry which information. Let us use a simple 32-bit RISC-V instruction as an example.

For an R-type instruction (register-to-register operation):

```
 31      25  24    20  19    15  14  12  11     7  6      0
 ┌─────────┬────────┬────────┬───────┬─────────┬────────┐
 │ funct7  │  rs2   │  rs1   │funct3 │   rd    │ opcode │
 │ 7 bits  │ 5 bits │ 5 bits │3 bits │ 5 bits  │ 7 bits │
 └─────────┴────────┴────────┴───────┴─────────┴────────┘

 funct7 = function code 7 bits (sub-type of instruction)
 rs2    = source register 2 (second operand)
 rs1    = source register 1 (first operand)
 funct3 = function code 3 bits (sub-type of instruction)
 rd     = destination register (where result goes)
 opcode = operation code (general instruction category)
```

Example: The instruction `ADD x7, x5, x6` (add the values in registers x5 and x6, store result in x7):

```
funct7  = 0000000
rs2     = 00110  (x6 = register 6)
rs1     = 00101  (x5 = register 5)
funct3  = 000
rd      = 00111  (x7 = register 7)
opcode  = 0110011 (R-type arithmetic/logic)

Full instruction in binary:
0000000 00110 00101 000 00111 0110011
```

### Decoding Step by Step

The Control Unit decodes an instruction in a precise sequence:

```
Step 1: Extract opcode
        opcode = bits [6:0] = 0110011
        → "This is an R-type ALU instruction"

Step 2: Identify instruction format
        opcode 0110011 → R-type format
        → read rs2, rs1, rd fields; use funct3+funct7

Step 3: Identify specific instruction
        funct3 = 000, funct7 = 0000000
        → "Specifically, this is ADD (not SUB, AND, OR, XOR...)"

Step 4: Extract register addresses
        rs1 = 5 → will read register x5 for operand A
        rs2 = 6 → will read register x6 for operand B
        rd  = 7 → will write result to register x7

Step 5: Generate control signals
        RegWrite = 1    (we need to write x7)
        ALUSrc   = 0    (second operand from rs2, not immediate)
        ALUOp    = ADD  (tell ALU to add)
        MemRead  = 0    (no memory access)
        MemWrite = 0    (no memory access)
        MemToReg = 0    (write ALU result to register, not memory data)
        Branch   = 0    (not a branch)
```

Notice that by the end of decoding, the Control Unit has extracted everything it needs. It knows where to get the data (rs1, rs2), what to do with it (ADD), and where to put the result (rd = 7). The control signals it generates will configure every component of the CPU to carry this out.

### The Decoder Circuit

The first stage of decoding is a **decoder** — a circuit that takes a binary number as input and activates exactly one of many output lines. A 7-bit opcode decoder has 128 output lines (2^7 = 128), one for each possible opcode value. Only one is high at any moment.

```
Opcode Decoder:
                    ┌──────────────────┐
                    │                  │── line_0000000 (unused)
opcode[6:0] ───────►│   7-to-128       │── line_0000001 (unused)
(7 bits in)         │   Decoder        │── ...
                    │                  │── line_0110011 ← HIGH (R-type)
                    │                  │── ...
                    │                  │── line_1111111 (unused)
                    └──────────────────┘
Only line_0110011 is 1; all others are 0.
```

### Quick Check 2

1. What is the opcode field of an instruction used for?
2. In the ADD x7, x5, x6 example, what are rs1, rs2, and rd?
3. What does a decoder circuit do?

---

## 3. Control Signals — The Conductor's Gestures

The conductor of an orchestra has a repertoire of gestures: a sweeping motion to bring in the strings, a sharp cut to silence the brass, a gentle tap to call the woodwinds. Each gesture causes a specific response from specific musicians.

The Control Unit's "gestures" are control signals. Let us examine each one carefully and understand exactly what it controls.

### The CPU Datapath with Control Points

```
                                          Control Signals
                                          (each line is 0 or 1)
                           ┌─────────────────────────────────────┐
                           │           CONTROL UNIT              │
         Instruction ─────►│                                     │
                           │   Decode logic                      │
                           └──┬──┬──┬──┬──┬──┬──┬──┬────────────┘
                              │  │  │  │  │  │  │  │
                         RegWrite │ ALUSrc │MemRead│ Branch
                              │ ALUOp  │MemWrite│MemToReg
                              │  │  │  │  │  │  │  │
                              ▼  ▼  ▼  ▼  ▼  ▼  ▼  ▼
    ┌──────────┐             ┌──────────────────────────────────┐
    │          │◄───RegWrite─┤                                  │
    │ Register │             │         CPU Datapath             │
    │  File    │─────────────►  ALU ◄── ALUSrc                  │
    │          │             │  Memory interface ◄── MemRead    │
    └──────────┘             │                   ◄── MemWrite   │
                             │  MUX ◄── MemToReg                │
                             │  PC logic ◄── Branch             │
                             └──────────────────────────────────┘
```

### Control Signal Reference Table

| Signal | Bits | Values and Meaning |
|--------|------|--------------------|
| **RegWrite** | 1 | 1 = write result to destination register; 0 = do not write |
| **RegRead1** | 1 | 1 = read rs1 from register file |
| **RegRead2** | 1 | 1 = read rs2 from register file |
| **ALUSrc** | 1 | 0 = second ALU operand is rs2; 1 = second operand is sign-extended immediate |
| **ALUOp** | 3 | 000=ADD, 001=SUB, 010=AND, 011=OR, 100=XOR, 101=SLT, 110=SLL, 111=SRL |
| **MemRead** | 1 | 1 = read from data memory this cycle |
| **MemWrite** | 1 | 1 = write to data memory this cycle |
| **MemToReg** | 1 | 0 = write ALU result to register; 1 = write memory read data to register |
| **Branch** | 1 | 1 = this is a conditional branch instruction |
| **Jump** | 1 | 1 = this is an unconditional jump |
| **PCNext** | 2 | 00=PC+4, 01=branch target (if condition true), 10=jump target |
| **ImmSel** | 2 | selects how to extract and sign-extend the immediate value |

### Why MUXes Are Everywhere

A **multiplexer (MUX)** is a component that selects one of several inputs to pass through to its output, based on a select signal. MUXes are the "switches" that the Control Unit controls.

Example: The ALU needs a second operand. Where does it come from?
- For `ADD x7, x5, x6`: it comes from register x6 (ALUSrc = 0)
- For `ADDI x7, x5, 10`: it comes from the constant 10 embedded in the instruction (ALUSrc = 1)

```
        Register x6 value ──────────┐
                                    ├──► MUX ──► ALU second input
        Immediate value 10 ─────────┘    ▲
                                         │
                                    ALUSrc signal
                                    (0 or 1, from Control Unit)

When ALUSrc = 0: MUX passes the register value
When ALUSrc = 1: MUX passes the immediate value
```

The Control Unit flips this switch differently for every instruction type.

### Quick Check 3

1. What does the MemToReg signal control? Give an example of when it would be 0 vs 1.
2. What is a MUX and why does the Control Unit need to control them?
3. If Branch = 1, does the CPU always jump? What else needs to be true?

---

## 4. The Instruction Cycle — Fetch, Decode, Execute

The CPU operates in a repeating cycle called the **instruction cycle** (or fetch-decode-execute cycle). Every single instruction the CPU runs goes through this same sequence of steps. The Control Unit coordinates all three phases.

Think of it like a factory assembly line:
- The **fetch** station picks up the next item (instruction) from the warehouse (memory)
- The **decode** station reads the item's label and determines what to do with it
- The **execute** station actually does the work

### Phase 1: Fetch

During the fetch phase, the CPU retrieves the next instruction from memory.

```
FETCH PHASE:
┌─────────────────────────────────────────────────────────┐
│                                                         │
│  1. PC (Program Counter) holds address of next instr.   │
│     PC = 0x1000 (for example)                           │
│                                                         │
│  2. Memory Address Register ← PC                        │
│     MAR = 0x1000                                        │
│                                                         │
│  3. Memory Read: fetch instruction at address 0x1000    │
│                                                         │
│  4. Instruction Register ← memory[0x1000]               │
│     IR = 00000000011000101000001110110011               │
│          (binary encoding of ADD x7, x5, x6)           │
│                                                         │
│  5. PC ← PC + 4 (advance to next instruction)          │
│     PC = 0x1004                                         │
│                                                         │
│  Control signals active during FETCH:                   │
│    MemRead = 1 (read from instruction memory)           │
│    PCWrite = 1 (update PC to PC+4)                      │
│    IRWrite = 1 (load instruction into IR)               │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Phase 2: Decode

During the decode phase, the Control Unit examines the instruction in the IR and configures all components.

```
DECODE PHASE:
┌─────────────────────────────────────────────────────────┐
│                                                         │
│  1. Control Unit reads IR = ADD x7, x5, x6             │
│                                                         │
│  2. Extracts:                                           │
│     opcode = 0110011 → R-type ALU instruction           │
│     funct3 = 000, funct7 = 0000000 → ADD               │
│     rs1 = 5, rs2 = 6 (register addresses)              │
│     rd = 7 (destination register)                      │
│                                                         │
│  3. Register File: simultaneously reads x5 and x6      │
│     (register reads happen during decode in most CPUs)  │
│     A = Register[5] = value in x5                      │
│     B = Register[6] = value in x6                      │
│                                                         │
│  4. Control signals set for the execute phase:          │
│     ALUSrc = 0    (use register x6, not immediate)      │
│     ALUOp  = ADD  (add operation)                       │
│     RegWrite = 1  (will write result to x7)            │
│     MemRead = 0, MemWrite = 0 (no memory access)        │
│     MemToReg = 0  (write ALU result, not memory data)   │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### Phase 3: Execute

During the execute phase, the actual operation takes place — guided entirely by the control signals set during decode.

```
EXECUTE PHASE (for ADD x7, x5, x6):
┌─────────────────────────────────────────────────────────┐
│                                                         │
│  1. ALU receives:                                       │
│     Input A = value from x5 (via RegRead1)             │
│     Input B = value from x6 (ALUSrc=0 → use register)  │
│     Operation = ADD (from ALUOp control)                │
│                                                         │
│  2. ALU computes: result = A + B                       │
│                                                         │
│  3. MemToReg = 0: the result goes directly to register  │
│     (not memory read data)                              │
│                                                         │
│  4. RegWrite = 1: result is written to rd = x7         │
│     Register[7] ← result                               │
│                                                         │
│  5. Cycle complete. CU returns to FETCH for next instr. │
│                                                         │
└─────────────────────────────────────────────────────────┘
```

### The Cycle Visualized

```
     ┌────────┐     ┌────────┐     ┌─────────┐     ┌────────┐
     │        │     │        │     │         │     │        │
     │ FETCH  │────►│ DECODE │────►│ EXECUTE │────►│ FETCH  │────►...
     │        │     │        │     │         │     │        │
     └────────┘     └────────┘     └─────────┘     └────────┘
         ▲                                              │
         │                                              │
         └──────────────────────────────────────────────┘
                   (loop forever until halted)
```

Some architectures add a fourth phase — **Memory Access** — and a fifth — **Write Back** — giving the classic 5-stage pipeline: Fetch → Decode → Execute → Memory → Writeback. But conceptually, the Control Unit coordinates all of them.

### Quick Check 4

1. List the three phases of the instruction cycle.
2. During which phase does the Control Unit set the control signals?
3. What happens to the PC (Program Counter) during the fetch phase?

---

## 5. Timing Diagram — Watching the Signals Unfold

A **timing diagram** shows how signals change over time, in sync with the clock. The clock is a square wave — it alternates between 0 and 1 at a fixed frequency (e.g., 3 GHz = 3 billion cycles per second). Most CPU operations happen on the rising edge of the clock (the transition from 0 to 1).

Below is a timing diagram for a single complete fetch-execute cycle, executing the instruction `ADD x7, x5, x6`.

```
Clock Cycle:    |  T1  |  T2  |  T3  |  T4  |  T5  |
                |      |      |      |      |      |
Phase:          | Fetch| Dec  | Exec | WB   | Idle |
                |      |      |      |      |      |

CLOCK:     ____      ____      ____      ____      ____
              |____|      |____|      |____|      |____|
              ↑    ↑      ↑    ↑      ↑    ↑      ↑    ↑
             rise fall   rise fall   rise fall   rise fall

──────────────────────────────────────────────────────────────────────

IRWrite:   ____      ____
          |    |    |    |
     _____|    |____|    |________________________________________________
           (active: latch instr into IR at T1 and T2 rising edge)

MemRead    ____
(instr):  |    |
     _____|    |_________________________________________________________
           (active T1: fetch instruction from memory)

PC_Incr:   ____
          |    |
     _____|    |_________________________________________________________
           (active T1: PC ← PC+4)

RegRead:          ____
                 |    |
     ____________|    |___________________________________________________
                  (active T2: read x5 and x6 from register file)

ALUSrc:           0000     0000
(signal value)    (both cycles show 0 = register operand mode)

ALUOp:                    ADD
(signal value)            (active T3: ADD operation selected)

ALU_Active:               ____
                         |    |
     ____________________|    |____________________________________________
                          (T3: ALU computes x5 + x6)

MemRead    0         0         0         0
(data):    (never active — this is ADD, no memory read)

MemWrite:  0         0         0         0
           (never active — this is ADD, no memory write)

RegWrite:                           ____
                                   |    |
     ______________________________|    |____________________________________
                                    (T4: write result to x7)

MemToReg:  0         0         0         0
           (always 0: write ALU result, not memory data)

──────────────────────────────────────────────────────────────────────
Result:    Register x7 receives (x5 + x6) at the end of T4.
```

### Reading the Timing Diagram

A few key observations:

1. **IRWrite is active at T1** — the instruction is fetched from memory and stored in the IR.
2. **MemRead (instruction) is active at T1** — the memory interface is reading during fetch.
3. **RegRead is active at T2** — both source registers are read simultaneously during decode.
4. **ALUOp is set at T3** — the ALU performs its operation during execute.
5. **RegWrite is active at T4 (writeback)** — the result is stored back to the register file.
6. **MemRead/MemWrite (data) are never active** — ADD does not touch data memory.

This is the Control Unit conducting: each signal rises and falls at precisely the right time, ensuring every component acts on the right data at the right moment.

### For a Load Instruction

For comparison, here is the signal activity for `LW x5, 8(x2)` (load word: x5 = memory[x2 + 8]):

```
Clock Cycle:    |  T1  |  T2  |  T3  |  T4  |  T5  |
                | Fetch| Dec  | Exec | Mem  | WB   |

IRWrite:        ACTIVE|      |      |      |      |
MemRead (instr):ACTIVE|      |      |      |      |
PC_Incr:        ACTIVE|      |      |      |      |
RegRead:              |ACTIVE|      |      |      |
ALUSrc:               |      |  1   |      |      | ← immediate mode
ALUOp:                |      | ADD  |      |      | ← compute address
ALU_Active:           |      |ACTIVE|      |      |
MemRead (data):       |      |      |ACTIVE|      | ← memory access!
MemWrite:             |      |      |  0   |      | ← no write
MemToReg:             |      |      |      |  1   | ← use mem data
RegWrite:             |      |      |      |ACTIVE|
```

Compare: a Load instruction has ALUSrc=1 (to use the immediate offset), triggers MemRead in the memory phase, and sets MemToReg=1 to route the memory data into the register. The Control Unit makes all these different signals happen at the right times.

### Quick Check 5

1. What does the rising edge of the clock trigger in a synchronous CPU?
2. In the timing diagram for ADD, why is MemRead (data) never active?
3. For a Store instruction (SW x5, 8(x2)), which signal would be active in the memory phase instead of MemRead?

---

## 6. Different Instructions, Different Signals

One of the most elegant aspects of the Control Unit is how the same set of control signals can describe radically different behaviors simply by setting different values. Let us walk through several instruction types.

### Complete Control Signal Table

| Instruction Type | Example | RegWrite | ALUSrc | ALUOp | MemRead | MemWrite | MemToReg | Branch | Jump |
|-----------------|---------|----------|--------|-------|---------|----------|----------|--------|------|
| R-type (register) | ADD x7, x5, x6 | 1 | 0 | funct code | 0 | 0 | 0 | 0 | 0 |
| I-type (immediate) | ADDI x7, x5, 10 | 1 | 1 | funct code | 0 | 0 | 0 | 0 | 0 |
| Load | LW x5, 8(x2) | 1 | 1 | ADD | 1 | 0 | 1 | 0 | 0 |
| Store | SW x5, 8(x2) | 0 | 1 | ADD | 0 | 1 | X | 0 | 0 |
| Branch | BEQ x1, x2, label | 0 | 0 | SUB | 0 | 0 | X | 1 | 0 |
| Unconditional Jump | JAL x1, label | 1 | — | — | 0 | 0 | PC+4 | 0 | 1 |
| NOP | NOP | 0 | X | — | 0 | 0 | X | 0 | 0 |

*X = don't care (value does not matter because signal is not used)*

### Walking Through Each Type

**R-type (ADD x7, x5, x6):**
Both operands come from registers (ALUSrc=0). The ALU adds them. Result writes to the register file (RegWrite=1). No memory involved.

**I-type (ADDI x7, x5, 10):**
One operand is a register, the other is the constant 10 embedded in the instruction (ALUSrc=1, immediate). Result writes to register. This is how you add a constant to a register.

**Load (LW x5, 8(x2)):**
The ALU computes an address: x2 + 8. This is an ADD operation, so ALUOp=ADD and ALUSrc=1 (to use the immediate offset 8). The result is a memory address. MemRead=1 reads from that address. MemToReg=1 sends the read data (not the ALU result) to the register file. RegWrite=1 stores it in x5.

**Store (SW x5, 8(x2)):**
Like a load, the ALU computes x2 + 8. MemWrite=1 writes register x5 to that address. RegWrite=0 because we are storing, not loading — no register gets updated.

**Branch (BEQ x1, x2, label):**
The ALU subtracts x1 - x2. If the result is zero, the branch is taken (branch condition = "equal"). Branch=1 tells the PC logic to be ready to jump. The PC updates to the branch target only if the zero flag is set. RegWrite=0, MemRead=0, MemWrite=0 — no data changes, only the PC might change.

**Jump (JAL x1, label):**
Unconditionally jumps to the label address. Also saves PC+4 into x1 (the return address for function calls). Jump=1. RegWrite=1 to save PC+4 in x1.

### The Control Unit as a Function

The Control Unit can be thought of as a pure function:

```
control_signals = ControlUnit(instruction_bits)
```

Given the same instruction, it always produces the same control signals. There is no internal state, no memory of past instructions. It is a stateless combinational circuit mapping inputs to outputs — elegant in its simplicity.

### Quick Check 6

1. For a Store instruction, why is RegWrite = 0?
2. Why does a Load instruction use ALUOp = ADD?
3. What two conditions must both be true for a branch to actually change the PC?

---

## 7. Hardwired Control — The Reflex Conductor

Imagine a conductor who has rehearsed a single symphony so many times that every gesture is a reflex — no thought required, no score consulted. When the violins reach measure 47, the arm rises automatically. This is **hardwired control**.

In hardwired control, the mapping from instruction bits to control signals is built directly into combinational logic gates. There is no table lookup, no stored program, no ROM. The gates compute the control signals directly from the instruction bits in nanoseconds.

### How It Works

For each control signal, a boolean expression is derived from the instruction fields, and that expression is implemented in gates.

Example: When should RegWrite be 1?

```
RegWrite should be 1 for:
  R-type instructions:  opcode = 0110011
  I-type ALU:           opcode = 0010011
  Load instructions:    opcode = 0000011
  JAL:                  opcode = 1101111
  JALR:                 opcode = 1100111
  LUI:                  opcode = 0110111
  AUIPC:                opcode = 0010111

RegWrite = (opcode==0110011) OR (opcode==0010011) OR 
           (opcode==0000011) OR (opcode==1101111) OR
           (opcode==1100111) OR (opcode==0110111) OR
           (opcode==0010111)

In gates:
  Each (opcode==value) is implemented by a set of AND gates with 
  appropriate NOT gates on the inverted bits.
  The final RegWrite is the OR of all these terms.
```

### The Complete Hardwired Control Unit Structure

```
                ┌─────────────────────────────────────────────────┐
                │           HARDWIRED CONTROL UNIT                │
                │                                                 │
opcode[6:0] ───►│  ┌─────────┐                                   │──► RegWrite
funct3[2:0] ───►│  │ Decoder │──► 128 opcode lines               │──► ALUSrc
funct7[6:0] ───►│  └─────────┘          │                        │──► ALUOp[2:0]
                │                       ▼                         │──► MemRead
                │              AND/OR/NOT gate network             │──► MemWrite
                │              (one network per output signal)    │──► MemToReg
                │                                                 │──► Branch
                │                                                 │──► Jump
                └─────────────────────────────────────────────────┘

No clock needed. No memory. No ROM. Just gates.
Input arrives → gates compute → outputs appear (a few gate delays later).
```

### Speed of Hardwired Control

The propagation delay through this logic is typically 3-5 gate delays, meaning the control signals are ready within a few hundred picoseconds of the instruction arriving. On a 3 GHz CPU (clock period = 333 picoseconds), this must complete in a fraction of one clock cycle.

This extreme speed is why modern RISC processors (ARM, RISC-V) use hardwired control.

### Advantages of Hardwired Control

```
┌────────────────────────────────────────────────────────┐
│  HARDWIRED CONTROL ADVANTAGES                          │
│                                                        │
│  SPEED:        Fastest possible — just gate delays     │
│                No memory access, no table lookup       │
│                                                        │
│  AREA:         Compact for simple ISAs                 │
│                Only gates, no ROM needed               │
│                                                        │
│  POWER:        Lower power than ROM-based approach     │
│                Fewer transistors switching             │
│                                                        │
│  SIMPLICITY:   Transparent design — trace any signal   │
│                back to specific instruction bits       │
└────────────────────────────────────────────────────────┘
```

### Disadvantages of Hardwired Control

```
┌────────────────────────────────────────────────────────┐
│  HARDWIRED CONTROL DISADVANTAGES                       │
│                                                        │
│  INFLEXIBILITY: Once made, the gates cannot change.    │
│                 A bug in the control logic is         │
│                 permanent. The chip must be recalled.  │
│                                                        │
│  COMPLEXITY:    Adding new instructions requires       │
│                 new gate networks. For x86 with        │
│                 thousands of instructions, the gate    │
│                 network would be enormous.             │
│                                                        │
│  VERIFICATION:  Hard to verify correctness.            │
│                 Each gate connection must be checked.  │
└────────────────────────────────────────────────────────┘
```

### Quick Check 7

1. In hardwired control, how is the mapping from instruction to control signals implemented?
2. Why is hardwired control faster than microprogrammed control?
3. What happens if a bug is found in a hardwired control circuit after the chip is manufactured?

---

## 8. Microprogrammed Control — The Sheet Music Conductor

Now imagine a different conductor — one who, for each new piece of music, opens a binder and follows a precisely written sequence of instructions. This conductor can handle any new symphony simply by getting a new binder. This is **microprogrammed control**.

In microprogrammed control, each machine instruction is implemented as a short program — a **microprogram** — stored in a special ROM called the **control store**. Each entry in the control store (one **microinstruction**) directly specifies the control signals for one micro-cycle.

### The Concept of Micro-Operations

A micro-operation (μop) is the most primitive action the datapath can perform in one clock cycle:

```
Examples of micro-operations:
  MAR ← PC                (copy PC to Memory Address Register)
  MDR ← Memory[MAR]       (read memory into Memory Data Register)
  IR ← MDR                (copy memory data to Instruction Register)
  PC ← PC + 4             (advance program counter)
  A ← Reg[rs1]            (read source register 1 into temp A)
  B ← Reg[rs2]            (read source register 2 into temp B)
  ALUout ← A + B          (ALU computes the sum)
  Reg[rd] ← ALUout        (write result to destination register)
```

A complex instruction is just a sequence of these micro-operations. The microprogram for ADD x7, x5, x6 would be:

```
Microprogram for ADD:
  μop 1 (Fetch):    MAR ← PC
  μop 2 (Fetch):    MDR ← Memory[MAR]; PC ← PC + 4
  μop 3 (Fetch):    IR ← MDR
  μop 4 (Decode):   A ← Reg[rs1]; B ← Reg[rs2]
  μop 5 (Execute):  ALUout ← A + B
  μop 6 (WB):       Reg[rd] ← ALUout
  μop 7:            goto Fetch (start next instruction)
```

### The Control Store Architecture

```
Machine instruction arrives at Control Unit
           │
           ▼
┌──────────────────────┐
│  Starting Address    │
│  Lookup Table        │  Maps opcode → starting μop address
│  (like an index)     │
└──────────┬───────────┘
           │
           ▼ starting address (e.g., address 42 for ADD)
┌────────────────────────────────────────────────────────────┐
│                  CONTROL STORE (ROM)                       │
│                                                            │
│  Address │ RegWrite │ ALUSrc │ ALUOp │ MemR │ MemW │ Next  │
│  ──────────────────────────────────────────────────────── │
│    0     │    0     │   0    │  NOP  │  1   │  0   │  1    │ Fetch step 1
│    1     │    0     │   0    │  NOP  │  1   │  0   │  2    │ Fetch step 2
│    2     │    0     │   0    │  NOP  │  0   │  0   │  3    │ Fetch step 3
│   ...    │   ...    │  ...   │  ...  │ ...  │ ...  │ ...   │
│   42     │    0     │   0    │  NOP  │  0   │  0   │  43   │ ADD decode
│   43     │    0     │   0    │  ADD  │  0   │  0   │  44   │ ADD execute
│   44     │    1     │   0    │  NOP  │  0   │  0   │  0    │ ADD writeback
│   ...    │   ...    │  ...   │  ...  │ ...  │ ...  │ ...   │
└────────────────────────────────────────────────────────────┘
           │
           ▼ μPC (micro-program counter) steps through the rows

Each row is a **microinstruction** — all control signals for one micro-cycle.
```

### The μPC — The Micro-Program Counter

Just as the PC points to the next machine instruction, the **micro-program counter (μPC)** points to the next microinstruction in the control store. After each micro-cycle:
- μPC advances to the next row (for sequential μops)
- Or jumps to another address (for branches in the microprogram)
- Or returns to address 0 to begin the next instruction fetch

### Historical Significance

Microprogrammed control was invented by Maurice Wilkes at Cambridge in 1951 — one of the most important computer architecture ideas ever. It made complex instruction sets practical:

- **IBM System/360 (1964):** The first widely microprogrammed family of computers. Multiple models with different performance levels ran the same software because they shared the same microprogram.
- **Intel 8086 (1978):** x86's enormous instruction set was made manageable by microcode.
- **Bug fixing:** Intel has issued microcode updates to fix processor bugs after manufacture — something impossible with hardwired control. The Pentium FDIV bug (1994) famously required both a recall AND microcode fix.

### Advantages of Microprogrammed Control

```
┌────────────────────────────────────────────────────────┐
│  MICROPROGRAMMED CONTROL ADVANTAGES                    │
│                                                        │
│  FLEXIBILITY:  New instructions can be added by        │
│                updating the microcode ROM.             │
│                No new silicon required.                │
│                                                        │
│  REPAIRABILITY: Bugs can be fixed post-manufacture     │
│                 by updating the microcode ROM.         │
│                 (Intel ships microcode patches via OS) │
│                                                        │
│  COMPLEXITY:   Can implement arbitrarily complex       │
│                instructions. Good fit for CISC ISAs.   │
│                                                        │
│  SIMPLICITY:   The microprogram is easier to read      │
│                and verify than gate netlists.          │
└────────────────────────────────────────────────────────┘
```

### Disadvantages of Microprogrammed Control

```
┌────────────────────────────────────────────────────────┐
│  MICROPROGRAMMED CONTROL DISADVANTAGES                 │
│                                                        │
│  SPEED:        Every instruction requires multiple     │
│                ROM lookups. ROM access is slower than  │
│                combinational logic.                    │
│                                                        │
│  AREA:         ROM takes more silicon than gate        │
│                networks for simple ISAs.               │
│                                                        │
│  POWER:        ROM access consumes more power than     │
│                gate-based computation.                 │
└────────────────────────────────────────────────────────┘
```

### Quick Check 8

1. What is a microinstruction and what does each row in the control store contain?
2. What is the μPC (micro-program counter)?
3. Why can Intel release microcode updates to fix security vulnerabilities? What would happen if the control were hardwired instead?

---

## 9. Hardwired vs Microprogrammed — A Comparison

| Feature | Hardwired Control | Microprogrammed Control |
|---------|------------------|------------------------|
| **Implementation** | Combinational logic gates | ROM (control store) + μPC |
| **Speed** | Very fast (gate delays only) | Slower (ROM lookup per cycle) |
| **Flexibility** | Fixed at manufacture | ROM can be updated |
| **Bug fixing** | Impossible after manufacture | Yes — update the ROM |
| **ISA complexity** | Best for simple/regular ISAs | Can handle very complex ISAs |
| **Silicon area** | Smaller for RISC | Larger (needs ROM) |
| **Power** | Lower | Higher |
| **Design effort** | High (complex gate network) | Lower (write microprogram) |
| **Used by** | ARM, RISC-V, MIPS | Original x86, mainframes |
| **Analogy** | Trained reflex | Following a written script |

### The RISC Revolution

The RISC movement in the 1980s (Patterson at Berkeley with RISC, Hennessy at Stanford with MIPS) was largely motivated by the observation that:
1. Microcode is slow — every instruction takes many micro-cycles
2. Simple instructions with hardwired control execute faster
3. Compilers can decompose complex operations into fast simple instructions

This led to the design philosophy: **fewer, simpler instructions, executed faster** — rather than many complex instructions executed slowly via microcode.

---

## 10. Modern CPUs — The Hybrid Approach

The modern reality, particularly for x86 processors (Intel Core, AMD Ryzen), is that neither pure hardwired nor pure microprogrammed control is used. Instead, modern CPUs use a **hybrid approach** that gets the best of both worlds.

### The x86 Decode Stage

```
x86 Instruction (variable length: 1 to 15 bytes!)
        │
        ▼
┌───────────────────────────────────────────────────────┐
│           FRONTEND DECODER                            │
│                                                       │
│  ┌─────────────────┐    ┌────────────────────────┐   │
│  │ Simple Decoder  │    │ Microcode ROM          │   │
│  │ (hardwired)     │    │ (for complex instr.)   │   │
│  │                 │    │                        │   │
│  │ Handles ~95% of │    │ Handles complex/rare   │   │
│  │ instructions:   │    │ instructions:          │   │
│  │  MOV, ADD, SUB, │    │  CPUID, string ops,    │   │
│  │  CMP, JMP, etc. │    │  XSAVE, VM enters,     │   │
│  │                 │    │  FP transcendentals    │   │
│  └────────┬────────┘    └──────────┬─────────────┘   │
│           │                        │                  │
└───────────┼────────────────────────┼──────────────────┘
            │                        │
            ▼                        ▼
        1-4 μops                up to 100+ μops
            │                        │
            └───────────┬────────────┘
                        │
                        ▼
              ┌─────────────────┐
              │   μop Queue     │  (RISC-like micro-operations)
              └────────┬────────┘
                       │
                       ▼
              ┌─────────────────┐
              │  Out-of-Order   │  (hardwired execution engine)
              │  Execution Core │
              └─────────────────┘
```

The key insight: by the time instructions reach the execution core, they have been converted to simple RISC-like μops. The execution core is essentially a RISC processor running inside the x86 chip. The complex x86 instructions are "interpreted" at the frontend.

### Why This Works So Well

- The common case (simple instructions) is handled by fast hardwired decoders → 1-2 μops per instruction
- The rare case (complex instructions) is handled by microcode → more μops but these are rare
- The execution core sees only simple μops and can be optimized aggressively

This is why modern x86 processors can achieve throughputs close to ARM processors, despite the much more complex ISA. The hardware decodes away the complexity before the actual execution happens.

### The Microcode Update Mechanism

When Intel or AMD discovers a security vulnerability or CPU bug:
1. They write new microcode sequences that avoid the problem
2. The new microcode is distributed as part of OS updates (Windows Update, Linux kernel patches)
3. At boot time, the OS loads the new microcode into the processor's microcode RAM
4. The processor uses the patched microcode going forward

Famous examples:
- **Spectre and Meltdown (2018):** Microcode patches restricted speculative execution across privilege boundaries
- **Hyperthreading vulnerabilities (2019):** Microcode patches altered how hyperthreading shares the execution core
- **Intel TAA (TSX Asynchronous Abort, 2019):** Microcode updates disabled the affected feature

This is an incredible engineering achievement — the ability to partially reprogram processor behavior without replacing the chip.

### Quick Check 9

1. In a modern x86 processor, what does the hardwired decoder handle vs. the microcode ROM?
2. What are μops in a modern x86 processor?
3. How does a microcode update reach the processor at runtime?

---

## Summary

The Control Unit is the conductor of the CPU orchestra — it reads each instruction, decodes its meaning, and issues precise control signals that tell every other component exactly what to do and when.

**Key concepts covered:**

- The Control Unit does not compute. Its sole job is to generate control signals from instruction bits.

- **Instruction decoding** extracts the opcode, register addresses, and immediate values from the binary instruction word, identifying exactly what operation to perform.

- **Control signals** configure every part of the CPU datapath: ALU operation select, register read/write enables, memory read/write enables, MUX selects, and PC control.

- The **instruction cycle** (fetch → decode → execute) is coordinated by the Control Unit. During fetch, the instruction is retrieved from memory. During decode, the Control Unit sets all signals. During execute, the datapath carries out the operation guided by those signals.

- **Timing diagrams** show how control signals rise and fall in precise coordination with the clock, with different signals active at different phases of the instruction cycle.

- Different instruction types (R-type, I-type, Load, Store, Branch, Jump) require entirely different patterns of control signals — the Control Unit provides the right pattern for each.

- **Hardwired control** implements the instruction-to-signals mapping as combinational logic. It is extremely fast but inflexible — bugs cannot be fixed after manufacture. Preferred for RISC architectures.

- **Microprogrammed control** stores each instruction's implementation as a sequence of microinstructions in a ROM (the control store). Slower but flexible — the microcode can be updated. Historically used for CISC architectures.

- **Modern CPUs use a hybrid:** fast hardwired decoders for common instructions, microcode ROM for complex or rare instructions, producing RISC-like μops for an efficient out-of-order execution engine.

The elegance of the Control Unit lies in how such a conceptually simple idea — "given this instruction, set these signals" — enables such staggering complexity. Every program ever run on every x86 or ARM processor ultimately reduces to a Control Unit setting wires to 0 and 1 at exactly the right moments, billions of times per second.

---

## Quick Reference

### Control Signal Summary

| Signal | Width | What It Controls |
|--------|-------|------------------|
| RegWrite | 1 bit | Whether to write the result to a register |
| ALUSrc | 1 bit | ALU second operand: register or immediate |
| ALUOp | 3 bits | Which ALU operation |
| MemRead | 1 bit | Read from data memory |
| MemWrite | 1 bit | Write to data memory |
| MemToReg | 1 bit | Write-back source: ALU result or memory data |
| Branch | 1 bit | Conditional branch instruction |
| Jump | 1 bit | Unconditional jump instruction |
| PCNext | 2 bits | Next PC: PC+4, branch target, or jump target |

### Instruction Type Signal Patterns

| Type | RegWrite | ALUSrc | MemRead | MemWrite | MemToReg | Branch |
|------|----------|--------|---------|----------|----------|--------|
| R-type | 1 | 0 | 0 | 0 | 0 | 0 |
| I-type ALU | 1 | 1 | 0 | 0 | 0 | 0 |
| Load | 1 | 1 | 1 | 0 | 1 | 0 |
| Store | 0 | 1 | 0 | 1 | X | 0 |
| Branch | 0 | 0 | 0 | 0 | X | 1 |
| Jump | 1 | — | 0 | 0 | PC+4 | 0 |

### Hardwired vs Microprogrammed Quick Reference

| Aspect | Hardwired | Microprogrammed |
|--------|-----------|-----------------|
| Speed | Fastest | Slower |
| Flexibility | None | High |
| Best for | RISC | CISC |
| Bug fix after manufacture | No | Yes |

---

## Exercises

### Easy

1. Fill in the control signals for an `OR x3, x1, x2` instruction (R-type, bitwise OR of x1 and x2 stored in x3):
   - RegWrite = ?
   - ALUSrc = ?
   - ALUOp = ?
   - MemRead = ?
   - MemWrite = ?
   - MemToReg = ?
   - Branch = ?

2. For the instruction `SW x5, 12(x3)` (store word: write register x5 to memory address x3+12), explain in plain English what the ALU does and why it needs to run even though we are writing to memory, not computing anything.

3. Draw a simple timing diagram (just mark phases as ACTIVE or 0) for a Branch instruction `BEQ x1, x2, label` where x1 equals x2 (branch is taken). Which control signals are active in each phase? What happens to the PC?

### Medium

4. Consider a CPU with a hardwired control unit. A designer discovers that the RegWrite signal has a bug — it is incorrectly set to 1 for Store instructions (it should be 0). What are the consequences of this bug at runtime? What are the engineer's options for fixing it before shipping the chip vs. after the chip has been manufactured and shipped to customers?

5. The timing diagram section showed control signals for ADD and LW. Construct the timing diagram (list active signals per clock phase) for `SW x5, 8(x2)` (store word). A Store instruction needs four phases: Fetch, Decode, Execute (address calculation), Memory (write). Pay careful attention to which signals are active in the Memory phase vs. the Execute phase.

6. In a microprogrammed control unit, the starting address lookup table maps each opcode to a starting microinstruction address. If RISC-V has 47 unique instruction types and each microprogram averages 6 microinstructions:
   - How large does the starting address table need to be (how many entries)?
   - How many microinstructions are in the control store (minimum)?
   - If each microinstruction is 20 bits wide, how many bytes of ROM does the control store need?
   - Compare this to a modern Intel processor's microcode ROM, which is estimated to contain roughly 125,000 microinstructions.

### Hard

7. The hybrid decoder in modern x86 processors must handle variable-length instructions (x86 instructions range from 1 to 15 bytes). This creates a problem called the **decode bottleneck**:
   - Why is decoding variable-length instructions harder than fixed-length instructions?
   - Intel uses a technique called "instruction length decode" as a pre-decode step. What information must this step provide to the main decoder?
   - Modern Intel processors can decode up to 6 instructions per clock cycle. Given that x86 instructions are variable length (1-15 bytes), how does the processor determine where each instruction begins?
   - Why can RISC processors (with fixed 32-bit instructions) decode more easily and at higher throughput?

8. Design a microprogrammed control unit for a minimal 8-bit CPU with 4 instructions:
   - `ADD Rd, Rs` (Rd = Rd + Rs)
   - `LOAD Rd, [addr]` (Rd = memory[addr], where addr is in the instruction)
   - `STORE [addr], Rs` (memory[addr] = Rs)
   - `BEQ addr` (if Z-flag set, PC = addr)
   
   (a) Define control signals needed (name each, give bit width).
   (b) Write the microprogram for each instruction as a sequence of μops.
   (c) Draw the control store as a table with all microinstructions for all four instructions plus the shared fetch sequence.
   (d) How many microinstructions are needed in total?

9. The Pentium FDIV bug (1994): Intel's Pentium processor had a bug in its floating-point divide unit that produced slightly wrong answers for certain inputs. Intel initially claimed it only affected a tiny fraction of calculations.
   - The bug was in the lookup table used by the SRT division algorithm, not in the microcode. Given this, why could Intel NOT fix it with a microcode update?
   - Intel eventually recalled all Pentium chips at a cost of $475 million. How would the situation have been different if the divide unit had been implemented in microcode (as in older processors)?
   - This bug led Intel to improve its formal verification processes. What does "formal verification" mean in the context of processor design, and why is it harder for hardwired logic than for microcode?
