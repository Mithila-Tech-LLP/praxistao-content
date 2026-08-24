# Chapter 23: Pipelining — Doing Many Things At Once

Pipelining is one of the most fundamental techniques in computer architecture. It allows a processor to work on multiple instructions simultaneously by breaking instruction execution into stages and overlapping them. Without pipelining, a modern CPU running at 3 GHz would be perhaps 10-20× slower. This chapter covers the pipeline from first principles to real implementations.

## Table of Contents

1. [The Laundry Analogy](#1-the-laundry-analogy)
2. [The 5-Stage RISC Pipeline](#2-the-5-stage-risc-pipeline)
3. [Pipeline Timing and Speedup](#3-pipeline-timing-and-speedup)
4. [Data Paths and Control Signals](#4-data-paths-and-control-signals)
5. [Pipeline Registers](#5-pipeline-registers)
6. [Performance: CPI and Throughput](#6-performance-cpi-and-throughput)
7. [Real Processor Pipeline Depths](#7-real-processor-pipeline-depths)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Laundry Analogy

Imagine you have 4 loads of laundry to do. Each load requires:
1. Wash (45 min)
2. Dry (45 min)
3. Fold (45 min)

### Sequential (No Pipelining)

Do one complete load at a time:
```
Load 1: [Wash][Dry][Fold]
Load 2:              [Wash][Dry][Fold]
Load 3:                           [Wash][Dry][Fold]
Load 4:                                        [Wash][Dry][Fold]

Total time: 4 × 3 × 45 = 540 minutes
```

### Pipelined (Overlap Operations)

Start load 2 washing while load 1 is drying:
```
       45   90  135  180  225  270  315
Load 1 [W ] [D ] [F ]
Load 2      [W ] [D ] [F ]
Load 3           [W ] [D ] [F ]
Load 4                [W ] [D ] [F ]

Total time: 3 × 45 + 3 × 45 = 270 minutes (2× speedup!)
```

The speedup comes from **overlapping** independent stages. The washer, dryer, and folder are all busy simultaneously.

**Ideal speedup** = number of pipeline stages. With 3 stages, we get ~2× speedup (not 3×, because we need to "fill" the pipeline at the start).

**The key requirement**: each stage must be independent — stage 1 (wash) doesn't need stage 2 (dry) to be complete before it can start on the next item.

### Quick Check

> 1. What is the key insight behind pipelining?
> 2. With N pipeline stages, what is the ideal speedup over a sequential processor?
> 3. Why can't you get exactly N× speedup with N stages?

---

## 2. The 5-Stage RISC Pipeline

The classic RISC pipeline has 5 stages:

```
IF → ID → EX → MEM → WB

IF:  Instruction Fetch  — Read the instruction from memory
ID:  Instruction Decode — Decode opcode, read registers
EX:  Execute           — ALU operation or address calculation
MEM: Memory Access      — Load from or store to data memory
WB:  Write Back         — Write result to register file
```

### Timing Diagram

Each instruction moves through all 5 stages, taking 5 cycles from start to finish. But with pipelining, a new instruction starts every cycle:

```
Cycle:    1    2    3    4    5    6    7    8    9
──────────────────────────────────────────────────────────
Instr 1:  IF   ID   EX   MEM  WB
Instr 2:       IF   ID   EX   MEM  WB
Instr 3:            IF   ID   EX   MEM  WB
Instr 4:                 IF   ID   EX   MEM  WB
Instr 5:                      IF   ID   EX   MEM  WB
```

All 5 stages are active simultaneously from cycle 5 onward. The pipeline is **full**.

### What Each Stage Does

#### IF: Instruction Fetch

```
PC → Instruction Memory → Instruction Register

1. Read instruction at address PC from instruction cache/memory
2. Store instruction in IF/ID pipeline register
3. Update PC: PC = PC + 4  (for next sequential instruction)
   or PC = branch_target  (if branch resolved — complication for later)
```

#### ID: Instruction Decode

```
Instruction → Control Unit + Register File

1. Decode instruction fields: opcode, rs1, rs2, rd, immediate
2. Read rs1 and rs2 from the register file
3. Sign-extend the immediate
4. Generate control signals for EX, MEM, WB stages
5. Store decoded values in ID/EX pipeline register
```

#### EX: Execute

```
rs1, rs2 (or immediate), control signals → ALU → result

For R-type: result = rs1 OP rs2  (ADD, SUB, AND, OR, ...)
For I-type: result = rs1 OP imm  (ADDI, ANDI, ...)
For Load:   result = rs1 + imm   (effective address for memory)
For Store:  result = rs1 + imm   (effective address for memory)
For Branch: result = (rs1 == rs2)?  Compare result; also compute branch target = PC+imm
```

#### MEM: Memory Access

```
For Load:  data = Memory[EX_result]  → read data from cache at computed address
For Store: Memory[EX_result] = rs2   → write data to cache
For other instructions: pass EX_result through (no memory operation)
```

#### WB: Write Back

```
For Load:   Register[rd] = MEM_data    (write loaded value)
For ALU:    Register[rd] = EX_result   (write computed value)
For Store:  (no register write)
For Branch: (no register write)
```

### Quick Check

> 1. In the 5-stage pipeline, which stage reads from the register file?
> 2. Which instruction types use the MEM stage to actually access data memory?
> 3. At steady state (pipeline full), how many instructions are "in flight" simultaneously?

---

## 3. Pipeline Timing and Speedup

### Clock Period with Pipelining

Without pipelining, the clock period must accommodate the LONGEST single instruction path:

```
Single-cycle time:
  IF time:   200ps (instruction cache access)
  ID time:   100ps (register file read)
  EX time:   200ps (ALU: adder, shifter, logic)
  MEM time:  200ps (data cache access)
  WB time:   100ps (register file write)

Total (combinational path): 800ps
Clock period (single-cycle): 800ps
Frequency: 1.25 GHz
```

With pipelining, the clock period only needs to accommodate the SLOWEST STAGE:

```
Pipeline clock period = max(IF, ID, EX, MEM, WB) = max(200, 100, 200, 200, 100) = 200ps
Frequency: 5 GHz  (4× faster than single-cycle!)
```

### Pipeline Overhead: Register Setup Time

In practice, pipeline registers (the flip-flops between stages) add latency — each pipeline register adds ~20-30ps overhead. With 4 pipeline register boundaries:

```
Effective clock period = max_stage + register_overhead = 200 + 25 = 225ps
Frequency: ~4.4 GHz
```

Still much faster than the 1.25 GHz single-cycle processor.

### Speedup Formula

For a pipeline with N stages and k instructions:

```
Time (single-cycle) = k × (sum of all stage times)
Time (pipelined)    = (N + k - 1) × clock_period

Speedup = Time_single / Time_pipelined
        = k × T_total / ((N + k - 1) × T_stage_max)

For large k (many instructions): Speedup ≈ T_total / T_stage_max = N stages
```

The speedup approaches N for large numbers of instructions (when pipeline fill/drain overhead is amortized). This is **throughput** speedup — for a long stream of instructions.

**Latency** (single instruction): pipelining DOES NOT help. A single instruction still takes N cycles to complete (5 cycles for 5-stage). Pipelining helps throughput, not latency.

### Quick Check

> 1. If stage times are IF=200ps, ID=150ps, EX=300ps, MEM=200ps, WB=100ps, what is the pipeline clock period (ignoring register overhead)?
> 2. How does pipelining affect the latency of a single instruction vs the throughput of many instructions?
> 3. For k=1000 instructions with a 5-stage pipeline, what is the theoretical speedup vs a single-cycle design (assuming equal stage times)?

---

## 4. Data Paths and Control Signals

### Datapath Overview

The complete pipelined datapath passes values through pipeline registers from stage to stage:

```
        IF Stage           ID Stage          EX Stage         MEM Stage        WB Stage
    ─────────────    ──────────────────  ─────────────────  ──────────────   ──────────────
                    IF/ID Register       ID/EX Register    EX/MEM Register  MEM/WB Register
    ┌──────────┐    ┌───────────────┐   ┌──────────────┐  ┌─────────────┐  ┌─────────────┐
    │ PC → I$ │ →  │ Instruction   │ → │ RS1_data     │ → │ ALU_result  │ → │ Read_data   │
    └──────────┘    │ PC_plus_4     │   │ RS2_data     │  │ RS2_data    │  │ ALU_result  │
                    └───────────────┘   │ Immediate    │  │ Rd (dest)   │  │ Rd (dest)   │
                                        │ Rd (dest)    │  │ MemWrite    │  │ MemToReg    │
                                        │ Control sigs │  │ MemRead     │  │ RegWrite    │
                                        └──────────────┘  └─────────────┘  └─────────────┘
```

### Control Signal Propagation

Control signals are generated in ID but needed in later stages. They must travel through the pipeline registers:

```
Stage ID generates:
  → EX controls: ALUSrc (use immediate or rs2?), ALUOp (which operation?)
  → MEM controls: MemRead, MemWrite
  → WB controls: RegWrite, MemToReg (write ALU result or memory data?)

These control signals travel in the pipeline registers:
  ID/EX register: carries all control signals + data values
  EX/MEM register: carries MEM+WB controls + EX results
  MEM/WB register: carries WB controls + MEM results
```

The control signals "age" through the pipeline alongside the data they control — always arriving at the right stage at the right time.

### Quick Check

> 1. Why do control signals need to travel through pipeline registers?
> 2. What does the `MemToReg` control signal determine?
> 3. What value goes in the WB stage for a STORE instruction (is there anything to write back)?

---

## 5. Pipeline Registers

Pipeline registers are the flip-flops that hold values between stages. They are clocked on every cycle, capturing the output of one stage and feeding it to the next.

### What Each Pipeline Register Contains

**IF/ID register**:
```
Instruction[31:0]   — the fetched instruction
PC_plus_4           — PC+4 for branch/return calculations
```

**ID/EX register**:
```
RS1_data[63:0]      — value read from register rs1
RS2_data[63:0]      — value read from register rs2
Immediate[63:0]     — sign-extended immediate
rd[4:0]             — destination register number
rs1[4:0]            — source register 1 number (for forwarding)
rs2[4:0]            — source register 2 number (for forwarding)
PC_plus_4           — for PC-relative addressing
Controls:
  EX: ALUSrc, ALUOp[2:0]
  MEM: MemRead, MemWrite
  WB: RegWrite, MemToReg
```

**EX/MEM register**:
```
ALU_result[63:0]    — computed result (address for load/store, value for ALU)
RS2_data[63:0]      — data to write for STORE instructions
rd[4:0]             — destination register number
Controls:
  MEM: MemRead, MemWrite
  WB: RegWrite, MemToReg
```

**MEM/WB register**:
```
Read_data[63:0]     — data read from memory (for LOAD)
ALU_result[63:0]    — ALU result (for non-LOAD instructions)
rd[4:0]             — destination register
Controls:
  WB: RegWrite, MemToReg
```

### Why Track rd Through the Pipeline?

The destination register number (rd) travels all the way to the WB stage to tell the register file which register to update. It also enables **hazard detection** — identifying when a later instruction reads a register that an earlier instruction (still in the pipeline) is going to write.

### Quick Check

> 1. Why does the ID/EX register carry rs1 and rs2 register numbers (not just the data values)?
> 2. For a STORE instruction, what travels in EX/MEM.RS2_data and why?
> 3. What does MemToReg select at the WB stage?

---

## 6. Performance: CPI and Throughput

### CPI: Cycles Per Instruction

In an ideal pipelined processor with no hazards or stalls:

**CPI = 1.0**

Every cycle, one instruction completes (one instruction enters the WB stage). This is the ideal.

In practice, stalls (due to hazards — covered in Chapter 24) increase CPI:

```
CPI = 1 + stall_cycles_per_instruction

If 30% of instructions cause a 1-cycle stall:
CPI = 1 + 0.30 × 1 = 1.30
```

### IPC: Instructions Per Cycle

IPC = 1/CPI. For a simple in-order pipeline:
- Ideal: IPC = 1.0
- With 30% one-cycle stalls: IPC = 1/1.30 = 0.77

For superscalar processors (Chapter 26), IPC > 1 is possible.

### Performance Equation

```
Performance = Instructions × CPI × Clock_period
Time        = Instructions / (IPC × Frequency)

Speedup = (old_instructions × old_CPI × old_period) / (new_instructions × new_CPI × new_period)
```

Often, Instructions is roughly constant between architectures, so:
```
Speedup ≈ (old_CPI / new_CPI) × (old_period / new_period)
         = IPC_ratio × Frequency_ratio
```

### Example

Processor A: 1 GHz, CPI = 1.2
Processor B: 3 GHz, CPI = 1.8

```
Time_A = N × 1.2 / 1GHz = 1.2N ns
Time_B = N × 1.8 / 3GHz = 0.6N ns

Speedup(B/A) = 1.2 / 0.6 = 2.0×
```

Processor B is 2× faster despite higher CPI, because frequency increase (3×) dominates CPI increase (1.5×).

### Quick Check

> 1. What is CPI and what value does an ideal pipeline achieve?
> 2. If 20% of instructions cause 2-cycle stalls, what is the CPI?
> 3. Processor A: 2 GHz, CPI=1.5. Processor B: 4 GHz, CPI=2.0. Which is faster and by how much?

---

## 7. Real Processor Pipeline Depths

Different processors make different tradeoffs on pipeline depth:

```
Processor              Year  Stages  Frequency   Strategy
─────────────────────  ────  ──────  ─────────   ──────────────────────────────────
Berkeley RISC-I        1982    3      2 MHz       Proof of concept
MIPS R2000             1988    5     16 MHz       Classic 5-stage
Intel 486              1989    5     25 MHz       x86 with 5-stage pipeline
Intel Pentium          1993    5      66 MHz      Dual-issue 5-stage
Intel Pentium Pro (P6) 1995   14    200 MHz       First OOO x86; µop pipeline
Intel Pentium 4 (NB)   2000   20    1.5 GHz       NetBurst: very deep for frequency
Intel Pentium 4 (Prescott) 2004 31  3.8 GHz      Maximum depth; thermal disaster
Intel Core (Yonah)     2006   14    2.0 GHz       Back to moderate depth
Intel Sandy Bridge     2011   14    3.5 GHz       Stable architecture
Intel Skylake          2015   14    4.0 GHz       Incremental
Apple M1 (Firestorm)  2020   ~15    3.2 GHz       High IPC, moderate depth
AMD Zen 3              2020   19    5.0 GHz       Deeper for higher frequency
```

### The Pipeline Depth Tradeoff

**Deeper pipeline** (more stages, each shorter):
- Higher maximum clock frequency (each stage does less work)
- Higher penalty for branch mispredictions (more stages to flush)
- Higher latency per instruction
- Intel Prescott (31 stages) hit a thermal wall: too much power at 4 GHz+

**Shallower pipeline** (fewer stages, each longer):
- Lower maximum frequency
- Lower branch misprediction penalty
- Lower latency per instruction
- Better for irregular code, harder to exploit frequency

The "sweet spot" for modern processors is roughly 12-20 stages.

### Quick Check

> 1. What was the Intel Prescott pipeline depth and what problem did it cause?
> 2. Why does deeper pipelining allow higher clock frequency?
> 3. What is the main cost of a deeper pipeline in terms of control flow?

---

## Summary

- **Pipelining** overlaps multiple instruction executions by dividing instruction execution into stages. Like laundry with washer, dryer, and folder running simultaneously.
- The **classic 5-stage RISC pipeline**: IF (fetch), ID (decode), EX (execute), MEM (memory), WB (write back). Each stage takes one clock cycle.
- **Pipeline speedup**: theoretically N× for N stages. In practice, limited by the slowest stage and overhead from stalls.
- **Clock period** = slowest stage time (+ register overhead). Pipelining enables higher clock frequency.
- **Pipeline registers** between stages hold the values (instruction, data, register numbers, control signals) needed by later stages.
- **Ideal CPI = 1.0** for a simple in-order single-issue pipeline. Stalls increase CPI.
- **Performance = Instructions × CPI × Period**. Pipelining reduces Period (higher frequency) at the cost of some CPI increase (from hazards).
- **Pipeline depth tradeoff**: deeper = higher frequency + higher stall penalties. Shallow = lower frequency + lower stall costs. Modern sweet spot: 12-20 stages.

---

## Exercises

### Easy

1. With a 5-stage pipeline and stage times of IF=200ps, ID=150ps, EX=300ps, MEM=200ps, WB=100ps:
   (a) What is the clock period for a pipelined implementation?
   (b) What is the clock period for a single-cycle implementation?
   (c) What is the speedup for 1000 instructions?

2. Draw the pipeline timing diagram (like the one in Section 2) for 6 instructions in a 5-stage pipeline. How many cycles does it take for all 6 instructions to complete?

3. If a processor has CPI=2.0 at 4 GHz, and another has CPI=1.2 at 3 GHz, which executes 1 billion instructions faster? By how much?

### Medium

4. **Pipeline register design**: For the ID/EX pipeline register in a 5-stage RISC-V pipeline, list all fields that must be stored, including their sizes in bits. Calculate the total number of bits in the ID/EX register. (Assume 64-bit registers, 5-bit register addresses, 3-bit ALUOp, and 5 single-bit control signals.)

5. **Depth vs frequency tradeoff**: The Intel Pentium 4 "Prescott" (2004) had 31 pipeline stages and ran at 3.8 GHz. The Intel Core 2 (2006) had 14 pipeline stages and ran at 2.93 GHz (lower frequency!). Yet Core 2 outperformed Prescott on nearly every benchmark. Explain how this is possible given that Prescott had higher frequency.

6. **Pipeline stages for modern instructions**: A simple pipeline executes one instruction type per stage. But modern processors have:
   - Integer ALU (1 cycle)
   - Multiply (3 cycles latency, 1 cycle throughput)
   - Divide (20-90 cycles, not pipelined)
   - FP add (4 cycles)
   - Load (4 cycles total: address + cache)
   
   How does a simple 5-stage pipeline handle these variable-latency instructions? Does it stall? What is the impact on CPI? Research: how do modern out-of-order processors handle this without stalling (hint: instruction scheduling)?

### Hard

7. **Pipeline design exercise**: You are designing a 4-stage RISC-V pipeline for a low-power embedded processor:
   - Stage 1 (FD): Fetch and Decode
   - Stage 2 (EX): Execute (ALU + branch resolution)
   - Stage 3 (MEM): Memory access
   - Stage 4 (WB): Write back
   
   Merging IF and ID into one stage has implications:
   (a) When is the branch target computed? When does the next PC get updated?
   (b) What is the branch misprediction penalty (how many cycles flushed)?
   (c) For a 5-stage pipeline vs this 4-stage pipeline, which has lower CPI for branch-heavy code? (Assume 15% branches, 80% prediction accuracy)
   (d) Design the forwarding paths: which stages forward to which?

8. **Formal performance analysis**: Given a workload with instruction mix:
   - 40% ALU (1 cycle EX)
   - 25% Load (4 cycle total latency including cache)
   - 15% Store (1 cycle EX, then write)
   - 15% Branch (1 cycle EX, 3-cycle penalty if mispredicted, 90% prediction accuracy)
   - 5% Multiply (3 cycle EX latency)
   
   For a 5-stage in-order pipeline with full forwarding, calculate:
   (a) CPI contribution from load-use hazards (assume 50% of loads are followed immediately by a dependent instruction)
   (b) CPI contribution from branch mispredictions
   (c) CPI contribution from multiply latency (assume 40% of multiply results are used immediately by the next instruction)
   (d) Total CPI and IPC
   (e) If clock frequency is 3 GHz, what is the throughput in MIPS?
