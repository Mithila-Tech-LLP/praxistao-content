# Chapter 24: Pipeline Hazards — When the Pipeline Stalls

The pipeline we built in Chapter 23 works beautifully when each stage takes exactly one clock cycle and each instruction is completely independent of the one before it. Reality is messier. Instructions depend on each other. Branches redirect the flow of execution. Two instructions might need the same hardware at the same time. These situations are called **hazards** — and every real processor must deal with them.

## Table of Contents

1. [What Is a Hazard?](#1-what-is-a-hazard)
2. [Data Hazards — When Instructions Depend on Each Other](#2-data-hazards--when-instructions-depend-on-each-other)
3. [Solutions: Stalling and Forwarding](#3-solutions-stalling-and-forwarding)
4. [Control Hazards — Branches and the Pipeline](#4-control-hazards--branches-and-the-pipeline)
5. [Structural Hazards — Resource Conflicts](#5-structural-hazards--resource-conflicts)
6. [The CPI Impact of Hazards](#6-the-cpi-impact-of-hazards)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. What Is a Hazard?

A **pipeline hazard** is any situation that prevents the next instruction from executing in its expected clock cycle. When a hazard occurs, the pipeline must either stall (wait) or take corrective action. Stalls reduce the ideal CPI of 1.0 — each stall cycle adds 1 to the CPI.

There are three classes of hazards:

| Type | Cause |
|------|-------|
| Data hazard | An instruction needs a result that hasn't been computed yet |
| Control hazard | A branch changes program flow; the pipeline doesn't know what to fetch next |
| Structural hazard | Two instructions need the same hardware resource simultaneously |

```
Ideal pipeline (no hazards):
Cycle:   1    2    3    4    5    6    7
I1:      IF   ID   EX   MEM  WB
I2:           IF   ID   EX   MEM  WB
I3:                IF   ID   EX   MEM  WB

Pipeline with a 2-cycle stall:
Cycle:   1    2    3    4    5    6    7    8    9
I1:      IF   ID   EX   MEM  WB
I2:           IF   ID  ---  ---  EX   MEM  WB
I3:                IF  ---  ---  ID   EX   MEM  WB
                        ↑↑ two bubble cycles inserted
```

### Quick Check
> 1. What is the ideal CPI of a pipelined processor?
> 2. If a processor has 15% of instructions causing 2-cycle stalls, what is the actual CPI?
> 3. Name the three classes of pipeline hazards.

---

## 2. Data Hazards — When Instructions Depend on Each Other

Data hazards occur when an instruction needs data that a previous instruction hasn't finished producing. There are three types, named by the order of Read (R) and Write (W) operations:

### RAW — Read After Write (True Dependency)

The most common and dangerous hazard. Instruction 2 reads a register that Instruction 1 is still writing.

```
I1: ADD R1, R2, R3    ; R1 = R2 + R3   (writes R1 in WB stage, cycle 5)
I2: SUB R4, R1, R5    ; R4 = R1 - R5   (reads R1 in ID stage, cycle 3)
```

In cycle 3, when I2 reaches the ID stage to read R1, I1 hasn't finished — it won't write R1 until cycle 5. I2 would read a stale value.

```
Cycle:   1    2    3    4    5
I1:      IF   ID   EX   MEM  WB   ← R1 written here
I2:           IF   ID   ←reads R1 here — TOO EARLY!
```

### WAR — Write After Read (Anti-Dependency)

Instruction 2 writes a register that Instruction 1 is still reading. Rare in simple 5-stage pipelines (ID comes before WB), but matters in out-of-order processors.

### WAW — Write After Write (Output Dependency)

Two instructions both write to the same register. Also rare in simple pipelines, important in out-of-order.

### Quick Check
> 1. In the RAW example above, which cycle does I1 write R1? Which cycle does I2 need to read R1?
> 2. Why is RAW the most dangerous hazard?
> 3. Can a WAR hazard occur in a simple 5-stage in-order pipeline? Why or why not?

---

## 3. Solutions: Stalling and Forwarding

### Solution 1: Pipeline Stalling (Bubbles)

The simplest solution: detect the hazard in the ID stage, freeze the instructions in IF and ID, and insert **bubbles** (NOP instructions) into the EX stage until the data is ready.

```
Cycle:   1    2    3    4    5    6    7
I1:      IF   ID   EX   MEM  WB
I2:           IF   ID  [NOP][NOP]  EX   MEM  WB
I3:                IF  [stall]     ID   EX   MEM
                   ↑ hazard detected: 2 bubbles inserted
```

Stalling works but costs 2 cycles per RAW hazard (for adjacent instructions). If 30% of instructions have a RAW dependency on the previous instruction, that's a 30% reduction in throughput.

### Solution 2: Data Forwarding (Bypassing)

A much smarter solution: instead of waiting for the result to be written to the register file, **route it directly** from the pipeline register where it sits to the EX stage input where it's needed.

```
Without forwarding:
  I1: ADD R1, R2, R3
  EX stage computes R1 = R2+R3, result sits in EX/MEM register...
  MEM stage passes result to MEM/WB register...
  WB stage writes to register file...  ← I2 can finally read R1

With forwarding:
  I1 EX computes result → forward directly to I2 EX input
  No stall needed!
```

The hardware adds multiplexers at the EX stage inputs. The hazard detection unit compares the destination register of instructions in later stages against the source registers of the instruction in ID. If a match is found, the forward path is selected instead of the register file output.

```
ALU input A normally comes from register file.
With forwarding MUX:
  ┌─ Register file output (normal path)
  ├─ EX/MEM pipeline register (forward from 1 instruction back)
  └─ MEM/WB pipeline register (forward from 2 instructions back)
```

Forwarding eliminates most RAW stalls. The remaining case: **load-use hazard**. A LOAD instruction reads from memory in the MEM stage. If the next instruction uses the loaded value, it needs it at the start of EX — but MEM hasn't produced it yet. One stall cycle is unavoidable.

```
I1: LOAD R1, [R2]     ; R1 loaded from memory (result available after MEM)
I2: ADD  R3, R1, R4   ; needs R1 at EX start — 1 cycle too early even with forwarding
                        → hardware inserts 1 bubble, compiler can reorder to avoid
```

### Quick Check
> 1. How many stall cycles does a RAW hazard between adjacent instructions require without forwarding? With forwarding?
> 2. What is a "load-use hazard" and why can't it be solved by forwarding alone?
> 3. What hardware does forwarding require that a basic pipeline doesn't have?

---

## 4. Control Hazards — Branches and the Pipeline

When the pipeline fetches a branch instruction, it doesn't yet know whether the branch is taken or not — that decision is made in the EX stage, two cycles after fetch. But the pipeline has already fetched the next two instructions. If the branch IS taken, those two instructions are wrong and must be discarded (flushed).

```
Cycle:   1    2    3    4    5    6
BEQ:     IF   ID   EX   MEM  WB
I+1:          IF   ID  [FLUSH if branch taken]
I+2:               IF  [FLUSH if branch taken]
                   branch target:
I_t:                   IF   ID   EX
```

If the branch is taken, 2 cycles are wasted (a 2-cycle branch penalty). With ~15-20% of instructions being branches and ~60% of those being taken, this penalty significantly reduces effective throughput.

### Solutions to Control Hazards

**1. Always stall**: Insert bubbles every time a branch is seen until the outcome is known. Simple but costly — wastes 2 cycles per branch regardless.

**2. Predict not-taken**: Continue fetching the instructions after the branch (as if it's not taken). If the branch IS taken, flush and start fetching from the target. Correct ~40% of the time for typical code. Penalty only on taken branches.

**3. Delayed branch** (used in MIPS): Define that the instruction immediately after a branch (the "delay slot") always executes, regardless of whether the branch is taken. The compiler fills the delay slot with useful work (an instruction from before the branch that doesn't affect the branch condition). Effectively hides the 1-cycle branch penalty. Ugly but practical for simple in-order pipelines.

**4. Dynamic branch prediction**: Hardware tracks past behavior of each branch and predicts whether it will be taken. Modern CPUs achieve >99% accuracy (Chapter 25 covers this in depth).

### Quick Check
> 1. Why does a branch cause a control hazard in a pipelined processor?
> 2. What does "flushing the pipeline" mean?
> 3. In delayed branching, what instruction goes in the delay slot, and when does it execute?

---

## 5. Structural Hazards — Resource Conflicts

A structural hazard occurs when two pipeline stages need the same hardware resource at the same time.

**Classic example**: A processor with a single unified memory (instructions and data in the same memory, accessed through the same port). In cycle 5, I1 is in the MEM stage (reading data) while I5 is in the IF stage (reading an instruction). They both need the memory at the same time — conflict!

```
Cycle:   5
I1:      MEM  ← needs memory to read data
I5:      IF   ← needs memory to fetch instruction
         ↑↑ both need the same memory port!
```

**Solutions**:
- **Separate instruction and data caches** (the standard solution in modern CPUs — this is the Modified Harvard architecture from Chapter 44). L1 is split: L1-I for instructions, L1-D for data. They can be accessed simultaneously.
- **Multi-ported memory**: Allow simultaneous accesses, but this is expensive in hardware.
- **Stalling**: The simpler but costlier solution — one access waits for the other.

Modern CPUs eliminate most structural hazards by duplicating hardware (multiple ALUs, separate caches, etc.), so structural hazards are rarely a concern in practice.

### Quick Check
> 1. Why does having separate instruction and data caches eliminate the most common structural hazard?
> 2. Name another potential structural hazard in a processor with only one multiply unit but instructions that both need to multiply.
> 3. How does the Modified Harvard architecture relate to structural hazards?

---

## 6. The CPI Impact of Hazards

The real cost of hazards is measured in CPI (Cycles Per Instruction). The ideal CPI for a 5-stage pipeline is 1.0. Hazards inflate it:

```
CPI_actual = CPI_ideal + stalls_per_instruction

Example workload analysis:
- 25% of instructions cause a RAW load-use hazard (1 stall each)
- 15% of instructions are branches with 50% taken (2 stalls for taken)
- Structural hazards: eliminated by split cache (0 stalls)

CPI = 1.0 + (0.25 × 1) + (0.15 × 0.5 × 2)
    = 1.0 + 0.25 + 0.15
    = 1.40

That's 40% worse than ideal — a significant impact!
```

With forwarding and branch prediction:
- Forwarding eliminates most RAW hazards (only load-use remains: 1 stall)
- Good branch prediction reduces taken-branch penalty from 2 to near 0

```
With forwarding + branch prediction (95% accuracy):
CPI = 1.0 + (0.05 × 1) + (0.15 × 0.05 × 15)  ← 15-cycle misprediction penalty
    = 1.0 + 0.05 + 0.11
    = 1.16
```

This is why branch prediction matters so much on deep-pipeline CPUs with large misprediction penalties.

### Quick Check
> 1. A program has 20% load instructions where 40% are followed immediately by a dependent instruction. What is the CPI contribution from load-use hazards?
> 2. Doubling the pipeline depth (from 5 to 10 stages) lets you double clock frequency but doubles the branch misprediction penalty. Is this a good trade?
> 3. Which hazard type is most important to solve in a deep pipeline CPU? Why?

---

## Summary

- **Data hazards** occur when an instruction needs a result still being computed. RAW (Read After Write) is the most common.
- **Stalling** solves data hazards by inserting bubble cycles — simple but costly.
- **Data forwarding** routes results directly from pipeline registers to where they're needed, eliminating most data hazard stalls. Load-use hazards still require 1 stall cycle.
- **Control hazards** occur at branches because the pipeline doesn't know which instructions to fetch. Solutions: stalling, predict-not-taken, delayed branches, dynamic prediction.
- **Structural hazards** occur when two instructions compete for the same hardware. Modern CPUs eliminate them by duplicating resources (split caches, multiple ALUs).
- Hazards inflate CPI above the ideal 1.0. Forwarding and branch prediction are the two most important mitigations.

---

## Exercises

### Easy
1. Draw a pipeline timing diagram for instructions I1–I4 where I2 has a RAW dependency on I1. Show the version with stalls and the version with forwarding.
2. A processor runs at 3 GHz with CPI = 1.3. What is its effective instruction throughput in MIPS (million instructions per second)?
3. What is a "bubble" in the pipeline and how does it differ from a real instruction?

### Medium
4. Given: 20% of instructions are loads, 30% of loads are immediately followed by a dependent instruction, 15% of instructions are branches, 60% of branches are taken. Calculate CPI assuming: (a) no forwarding and a 2-cycle data hazard stall; (b) forwarding (eliminating all but load-use stalls) and predict-not-taken with 2-cycle penalty on taken branches.
5. Explain why the compiler can eliminate some load-use hazards by reordering instructions. Give a concrete example with 4 instructions where reordering eliminates the stall without changing the result.
6. A MIPS processor uses delayed branches. The compiler is trying to fill the delay slot for: `BEQ R1, R2, target` followed by `ADD R3, R4, R5`. Can the ADD be moved to the delay slot? Under what conditions?

### Hard
7. Design the forwarding unit for a 5-stage pipeline. Specify the exact conditions (in terms of pipeline register contents) for: (a) forwarding from EX/MEM to EX input A, (b) forwarding from MEM/WB to EX input A. Handle the case where both forwarding conditions are true simultaneously (which takes priority?).
8. A speculative processor fetches instructions past a branch and begins executing them. If the branch is mispredicted, the processor must "undo" any side effects. List three types of side effects that must be carefully managed: memory writes, register writes, and I/O operations. For each, explain how the processor can undo them or prevent them until the branch is resolved.
