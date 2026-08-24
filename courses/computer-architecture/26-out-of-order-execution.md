# Chapter 26: Out-of-Order Execution — Working Ahead

In-order execution is like a worker on an assembly line who stops the entire line when they can't complete their current task. If a part is missing, everyone waits. Out-of-order execution is like a smarter worker who, when blocked on one task, immediately picks up another ready task from a queue. The blocked task waits; work continues elsewhere. This insight — that instructions don't have to execute in the order they are written, as long as results are correct — is responsible for a large fraction of modern CPU performance.

## Table of Contents

1. [The Motivation: Wasted Execution Slots](#1-the-motivation-wasted-execution-slots)
2. [The Key Insight: Instructions Have Independent Work](#2-the-key-insight-instructions-have-independent-work)
3. [The Reorder Buffer — Maintaining Program Order](#3-the-reorder-buffer--maintaining-program-order)
4. [Reservation Stations — Waiting for Operands](#4-reservation-stations--waiting-for-operands)
5. [Register Renaming — Eliminating False Dependencies](#5-register-renaming--eliminating-false-dependencies)
6. [Tomasulo's Algorithm — Putting It Together](#6-tomasulos-algorithm--putting-it-together)
7. [The Out-of-Order Pipeline in Full](#7-the-out-of-order-pipeline-in-full)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Motivation: Wasted Execution Slots

Consider this instruction sequence:

```
I1: LOAD  R1, [R2]       ; Load from memory (takes 100+ cycles if cache miss)
I2: ADD   R3, R1, R4     ; Uses R1 — must wait for I1
I3: MUL   R6, R7, R8     ; Independent of I1 and I2!
I4: SUB   R9, R10, R11   ; Independent of everything above
I5: ADD   R12, R3, R6    ; Depends on I2 and I3
```

In an in-order processor:
- I1 starts a memory load
- I2 is stalled waiting for I1's result (100+ cycles if cache miss)
- I3 and I4 also stall, even though they are completely independent
- The ALU sits idle for 100+ cycles

This is catastrophically wasteful. I3 and I4 could execute immediately — they don't care about R1 at all.

Out-of-order (OOO) execution solves this: while I2 waits for I1's load, the processor looks ahead, finds I3 and I4 are ready, and executes them now. The ALU is kept busy. When I1's load finally completes, I2 executes with the result, and then I5 can follow.

### Quick Check
> 1. In the example above, which instructions are independent of the cache miss in I1?
> 2. In an in-order processor, how many cycles does the ALU sit idle waiting for the cache miss?
> 3. What fundamental property of I3 and I4 makes them safe to execute before I2?

---

## 2. The Key Insight: Instructions Have Independent Work

A program's instructions form a **dependency graph** — a directed graph where an edge from instruction A to instruction B means "B depends on A's output." Instructions with no incoming edges from pending instructions are **ready to execute**.

```
Dependency graph for the earlier example:

I1 (LOAD R1) ──→ I2 (ADD R3, R1)
                              │
I3 (MUL R6) ────────────────→│ I5 (ADD R12, R3, R6)
I4 (SUB R9)  [no dependencies — isolated]
```

I3 and I4 are ready immediately. I2 is ready only after I1. I5 is ready only after I2 and I3.

The critical path is I1 → I2 → I5. No matter how clever the hardware, these three instructions must execute in order. But I3 and I4 can overlap with I1's long latency.

**Instruction-Level Parallelism (ILP)** is the degree to which a program's instructions can be executed in parallel given their dependency constraints. Typical integer code has an ILP of 2-4; floating-point scientific code can have ILP > 10.

### Quick Check
> 1. Draw the dependency graph for: `A=B+C; D=A*E; F=G-H; I=D+F`
> 2. What is the critical path length in instruction count for the above?
> 3. If you had unlimited parallel execution units, what is the minimum number of cycles to execute the four instructions above?

---

## 3. The Reorder Buffer — Maintaining Program Order

Out-of-order execution creates a problem: instructions complete in the wrong order. If instruction 5 finishes before instructions 1, 2, 3, 4 — and instruction 1 causes a page fault — we need to discard 2, 3, 4, 5 and handle the fault. But we've already written 5's result to the register file. How do we undo that?

The answer is the **Reorder Buffer (ROB)**: a circular buffer that holds all in-flight instructions in program order, from oldest (head) to newest (tail).

```
ROB (circular buffer):
Head ──► [I1: LOAD  R1 | executing ] 
         [I2: ADD   R3 | waiting   ]
         [I3: MUL   R6 | done=42   ]
         [I4: SUB   R9 | done=17   ]
Tail ──► [I5: ADD  R12 | waiting   ]
```

Each ROB entry holds:
- The instruction
- The destination register
- The computed result (once available)
- A "done" flag
- Exception status

**Commit (Retire)**: Instructions commit only from the **head** of the ROB, in program order. When the head entry is done and has no exception, its result is written to the architectural register file and the entry is removed. Even if I3 and I4 finish first, they wait in the ROB until I1 and I2 have safely committed.

This gives us:
- **Out-of-order execution**: instructions execute whenever their operands are ready
- **In-order commit**: results become architectural state in program order

**Precise exceptions**: When an exception occurs, the ROB allows the processor to commit all instructions before the exception, flush all instructions after it, and invoke the exception handler — as if execution had been in-order all along. This is called a **precise exception model** and is required for correct OS operation.

### Quick Check
> 1. Why must the ROB commit from the head in program order?
> 2. What is a "precise exception" and why does it require in-order commit?
> 3. If the ROB has 256 entries, what is the maximum "instruction window" — the number of instructions in flight simultaneously?

---

## 4. Reservation Stations — Waiting for Operands

The ROB ensures correct ordering. But how does the hardware know when an instruction is ready to execute?

**Reservation Stations (RS)** are buffers attached to each execution unit. An instruction is issued into a reservation station with either its operand values (if available) or a tag identifying which ROB entry will produce the value.

```
Reservation Station for ALU:
Slot | Instruction | Src A value | Src A ready? | Src B value | Src B ready?
  1  | ADD R3,R1,R4 |     ?      |   No (tag=1) |    10       |    Yes
  2  | SUB R9,R10,R11|   55      |    Yes       |    22       |    Yes   ← ready!
  3  | ADD R12,R3,R6 |     ?      |   No (tag=2) |     ?      |   No (tag=3)
```

Slot 2 has both operands ready — it fires to the ALU immediately. When I1 eventually produces R1, it broadcasts the result on the **Common Data Bus (CDB)**. Any reservation station waiting for that tag captures the value and marks itself ready.

The key mechanism:
1. **Issue**: Instruction enters reservation station with operand values or forwarded tags
2. **Wake-up**: When an execution unit completes, it broadcasts (tag, value) on the CDB
3. **Select**: The reservation station wakeup logic fires the oldest ready instruction to the execution unit
4. **Execute**: The instruction executes on the available unit

### Quick Check
> 1. What information does a reservation station entry hold for an operand that isn't ready yet?
> 2. What is the Common Data Bus (CDB) and what does it broadcast?
> 3. If two reservation station entries both become ready in the same cycle, how does the hardware decide which one executes first?

---

## 5. Register Renaming — Eliminating False Dependencies

Remember WAR and WAW hazards from Chapter 24? These aren't true dependencies — they're caused by reuse of register names, not actual data dependencies. Out-of-order processors eliminate them through **register renaming**.

### The Problem

```
I1: ADD R1, R2, R3     ; produces R1 (v1)
I2: MUL R4, R1, R5     ; consumes R1 (v1) — true RAW dependency
I3: ADD R1, R6, R7     ; produces R1 (v2) — WAW with I1
I4: SUB R8, R1, R9     ; consumes R1 (v2) — depends on I3, not I1!
```

Without renaming: I3 can't execute until I1 commits (WAW). I4 depends on which version of R1 it gets. With only one R1 register, this is a mess.

### The Solution: Physical Registers

Modern CPUs have many more **physical registers** than the ISA specifies architectural registers. For example, x86-64 specifies 16 architectural registers; Intel Sunny Cove has 280 physical integer registers.

When an instruction writes a register, it is assigned a fresh physical register. The mapping from architectural register name to current physical register is maintained in the **Register Alias Table (RAT)** or **Register Renaming Table**.

```
Renaming the example:
I1: ADD P5,  P2, P3      ; R1 → P5 (update RAT: R1 = P5)
I2: MUL P6,  P5, P7      ; reads P5 (the R1 from I1) — true dependency preserved
I3: ADD P8,  P9, P10     ; R1 → P8 (update RAT: R1 = P8) — different physical reg!
I4: SUB P11, P8, P12     ; reads P8 (the R1 from I3) — correct dependency
```

Now I1 and I3 write to different physical registers (P5 and P8). They have no WAW hazard. I2 clearly depends on I1's output (P5), not I3's. All false dependencies are eliminated.

Register renaming also enables **speculation**: the processor can rename and execute instructions past a branch, using physical registers for speculative results. If the branch mispredicts, the physical registers allocated for speculative instructions are simply freed — no architectural state was corrupted.

### Quick Check
> 1. What is the difference between a true dependency and a false dependency?
> 2. How many physical registers does Intel Sunny Cove have for integer operations vs how many does x86-64 define?
> 3. If an instruction uses register R5 and you have 280 physical registers mapped to 16 architectural registers, what is the maximum number of "in-flight" versions of R5 possible simultaneously?

---

## 6. Tomasulo's Algorithm — Putting It Together

All these ideas — reservation stations, common data bus, register renaming — were invented together in 1967 by Robert Tomasulo at IBM for the IBM System/360 Model 91 floating-point unit. He received the Eckert–Mauchly Award in 1997 for this work. The algorithm bears his name.

Tomasulo's algorithm combines:
1. **Register renaming via reservation station tags** (each RS entry is a "virtual register")
2. **Dynamic scheduling** via reservation station ready logic
3. **CDB broadcast** to wake up waiting instructions

Modern OOO processors use a refined version with a physical register file (rather than distributing values in RS entries) and an explicit ROB for in-order commit — but the core ideas are unchanged from 1967.

```
Original Tomasulo Flow:
  Instruction decode → check RAT → issue to Reservation Station
       (with operand values or pending tags)
  ↓
  RS entry waits for tags to resolve via CDB
  ↓
  When ready: instruction fires to execution unit
  ↓
  Result broadcast on CDB → updates waiting RS entries + register file
```

### Quick Check
> 1. When was Tomasulo's algorithm invented, and for what machine?
> 2. What does "issue" mean in the context of Tomasulo's algorithm?
> 3. Modern CPUs use physical register files rather than distributing values in RS entries. What advantage does this provide?

---

## 7. The Out-of-Order Pipeline in Full

Putting it all together, a modern OOO pipeline has these stages:

```
Fetch → Decode → Rename/Dispatch → Issue → Execute → Writeback → Commit

In-order:         │                          │         In-order:
  Fetch           │   Out-of-order happens   │           Commit
  Decode          │   in the middle          │
  Rename          │                          │
```

**Fetch**: Fetch up to 4–8 instructions per cycle, using the BTB to predict branch targets.

**Decode**: Decode each instruction. For x86: split CISC instructions into RISC-like micro-ops (µops). A single x86 instruction like `ADD [mem], reg` might become 3 µops: LOAD, ADD, STORE.

**Rename/Dispatch**: Allocate ROB entries and physical registers; update the RAT; dispatch µops to reservation stations.

**Issue (Schedule)**: Examine all reservation stations; select ready µops and send to execution units (integer ALU, FPU, load/store unit, branch unit).

**Execute**: Compute result. Load/store accesses cache. Branch checks prediction.

**Writeback**: Broadcast result on CDB; write to physical register; mark ROB entry done.

**Commit (Retire)**: ROB head entries that are done commit in order — write to architectural state, free physical registers.

**Typical modern widths and sizes:**
| Resource | Typical Size |
|----------|-------------|
| Decode width | 4–6 instructions/cycle |
| Reorder Buffer | 256–600 entries |
| Reservation Stations | 60–120 entries per unit |
| Physical registers | 180–280 integer + similar FP |
| Issue ports | 10–16 execution ports |

Apple M1's Firestorm core has ~630 ROB entries — one of the largest in the industry, enabling very deep speculative execution.

### Quick Check
> 1. Why does an x86 decoder convert CISC instructions into RISC micro-ops?
> 2. A processor has a 256-entry ROB. If the average instruction takes 10 cycles from issue to commit, what is the maximum IPC it can sustain (assuming no other bottlenecks)?
> 3. What happens to all the in-flight instructions in the ROB when a branch misprediction is detected?

---

## Summary

- **In-order execution** stalls on any dependency, wasting execution unit cycles while operands are computed.
- **Out-of-order execution** uses a dependency graph to find and execute ready instructions, keeping execution units busy.
- The **Reorder Buffer (ROB)** holds all in-flight instructions in program order. Instructions execute out-of-order but commit in-order, ensuring precise exceptions and correct visible state.
- **Reservation Stations** hold instructions waiting for operands. The **Common Data Bus** broadcasts completed results, waking up waiting instructions.
- **Register renaming** maps architectural registers to physical registers, eliminating false WAR/WAW dependencies and enabling aggressive speculative execution.
- **Tomasulo's algorithm** (1967) first combined these ideas. Modern processors use the same principles, scaled up enormously.
- The Apple M1's Firestorm core has ~630 ROB entries — the deep instruction window enables hiding latency from even L2 cache misses.

---

## Exercises

### Easy
1. Given the instruction sequence: `A=B+C; D=E*F; G=A+D; H=I-J`, identify all true (RAW) dependencies and all false (WAR/WAW) dependencies.
2. The ROB has 200 entries. Instructions take an average of 8 cycles from issue to commit. What is the maximum sustainable IPC?
3. Explain in one paragraph what the Common Data Bus does and why it is necessary.

### Medium
4. Trace Tomasulo's algorithm through the first 4 instructions of the example from Section 1 (LOAD, ADD, MUL, SUB). Show the state of reservation stations and ROB after each instruction is issued.
5. An OOO processor has 16 architectural registers (x86-64) and 256 physical registers. What is the maximum number of in-flight instructions that write different registers simultaneously? Why can't you just keep allocating registers forever?
6. Register renaming eliminates WAR and WAW hazards. Show a concrete example of a WAR hazard in a register-file-only machine, then show how renaming eliminates it, using physical register names P1-P20.

### Hard
7. The ROB enables precise exceptions. Trace exactly what happens when instruction I3 causes a page fault while instructions I4, I5, and I6 have already completed out-of-order in the ROB. Which instructions commit? What happens to I4-I6's results? What state does the hardware present to the OS exception handler?
8. Modern CPUs decode x86 instructions into micro-ops. The instruction `LOCK XADD [mem], reg` (atomic add) requires: LOAD the memory value, ADD with register, STORE back — all atomically. How does the ROB handle the "atomic" requirement? Why can't you simply split this into 3 independent µops that might be interleaved with other threads' memory accesses?
