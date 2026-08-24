# Chapter 27: Superscalar Processors — Multiple Pipelines

A scalar processor issues one instruction per clock cycle. A superscalar processor issues multiple instructions per clock cycle simultaneously. Where a scalar processor has one lane on the highway, a superscalar has four, six, or eight lanes — all flowing at the same speed but carrying more traffic. Combined with out-of-order execution, superscalar processing is why modern CPUs can execute four or more instructions per cycle instead of the one that a naive pipeline would suggest.

## Table of Contents

1. [Scalar vs Superscalar](#1-scalar-vs-superscalar)
2. [Execution Units — The Specialized Workers](#2-execution-units--the-specialized-workers)
3. [IPC: Instructions Per Cycle](#3-ipc-instructions-per-cycle)
4. [The ILP Wall — The Limits of Superscalar](#4-the-ilp-wall--the-limits-of-superscalar)
5. [VLIW — The Compiler-Driven Alternative](#5-vliw--the-compiler-driven-alternative)
6. [SMT — Simultaneous Multithreading](#6-smt--simultaneous-multithreading)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Scalar vs Superscalar

A **scalar** processor fetches, decodes, and issues exactly one instruction per clock cycle. A **superscalar** processor fetches, decodes, and issues multiple instructions per clock cycle — typically 2, 4, 6, or 8.

```
Scalar pipeline (1 instruction/cycle):
Cycle:  1    2    3    4    5    6    7    8
I1:     IF   ID   EX   MEM  WB
I2:          IF   ID   EX   MEM  WB
I3:               IF   ID   EX   MEM  WB

4-wide Superscalar (ideal — no dependencies):
Cycle:  1    2    3    4    5
I1:     IF   ID   EX   MEM  WB
I2:     IF   ID   EX   MEM  WB
I3:     IF   ID   EX   MEM  WB
I4:     IF   ID   EX   MEM  WB
         └── all 4 in parallel ──┘
```

To issue 4 instructions per cycle, the processor needs 4× the hardware at each stage:
- **4 instruction fetch units** (or a wider fetch that reads 4 instructions at once)
- **4 decoders** working in parallel
- **Dispatch logic** that assigns 4 decoded instructions to appropriate reservation stations
- **Multiple execution units** (ALUs, FPUs, load/store units) to actually run 4 instructions per cycle
- **4 result write-back paths**

The dispatch/schedule logic is the most complex part: it must, every cycle, examine all ready instructions in all reservation stations and select the best combination to send to available execution units.

### Quick Check
> 1. What hardware must be duplicated to go from a 1-wide to a 4-wide superscalar processor?
> 2. Why can't you simply issue 4 copies of the same instruction every cycle to achieve 4× speedup?
> 3. A 4-wide superscalar issues on average 2.8 instructions per cycle on real code. Why can't it reach 4.0?

---

## 2. Execution Units — The Specialized Workers

Instructions are not equal — different operations require different hardware. A modern superscalar CPU has multiple specialized execution units, each connected to one or more "ports" through which instructions are dispatched.

**Intel Core (Sunny Cove / Willow Cove) execution ports:**

| Port | Units Available |
|------|----------------|
| Port 0 | ALU, multiply, AES, branch |
| Port 1 | ALU, shift, divide, FPU |
| Port 2 | Load address |
| Port 3 | Load address |
| Port 4 | Store data |
| Port 5 | ALU, shuffle, vector |
| Port 6 | ALU, branch |
| Port 7 | Store address |

With 8 ports, Intel can execute up to 8 micro-operations per cycle. The scheduler must match each ready µop to a port that supports that operation type.

**Apple M1 Firestorm execution ports (estimated):** 4 integer ALU, 2 load, 1 store, 2 FP/SIMD — approximately 9 total ports, 8-wide decode.

**Load/Store Units**: Usually 2–3 load units and 1–2 store units. Loads and stores access the L1 data cache. The Load-Store Unit (LSU) also checks for memory ordering violations — if a load speculatively reads a value that a later store should have written, the load must be re-executed.

**Floating-Point Units (FPUs)**: Handle IEEE 754 arithmetic on 32-bit and 64-bit values. FP operations have longer latencies than integer (typically 3–5 cycles for add/multiply, 10–20 cycles for divide/sqrt). Modern FPUs also handle SIMD operations.

**Branch Units**: Evaluate branch conditions and compare the prediction against the actual outcome. If mispredicted, flush the pipeline.

### Quick Check
> 1. Intel's Sunny Cove has 8 execution ports. Does this mean it can always execute 8 instructions per cycle?
> 2. Why are floating-point divide and sqrt operations so much slower than multiply?
> 3. What is a "memory ordering violation" and why does the LSU need to check for it?

---

## 3. IPC: Instructions Per Cycle

**IPC (Instructions Per Cycle)** is the fundamental metric of microarchitecture efficiency. Higher IPC = more work per clock cycle = faster programs at the same frequency.

```
Performance = Clock Frequency × IPC × (1 / Instructions per task)

Two ways to double performance:
1. Double frequency (hard: requires smaller transistors, more power)
2. Double IPC (hard: requires more complex hardware, wider superscalar)
```

### Theoretical vs Practical IPC

A 6-wide superscalar theoretically achieves 6 IPC. In practice on real workloads:

| Processor | Theoretical | SPEC INT | SPEC FP |
|-----------|-------------|----------|---------|
| Intel Skylake (4-wide) | 4 | ~2.0 | ~2.5 |
| AMD Zen 4 (6-wide) | 6 | ~2.5 | ~3.0 |
| Apple Firestorm (8-wide) | 8 | ~3.5 | ~4.5 |

(These sustained-IPC figures are rough averages across benchmark suites — individual programs range from below 1 to near the theoretical peak.)

The gap between theoretical and practical IPC is caused by:
1. **Dependencies** limit ILP in the instruction stream
2. **Cache misses** stall the pipeline
3. **Branch mispredictions** flush in-flight work
4. **Execution unit imbalance** — some ports are bottlenecks
5. **Decode/dispatch bandwidth** — can't always decode complex instructions fast enough

Apple Firestorm achieves the highest IPC of any processor in history for typical workloads — this is why M1 MacBooks feel so responsive despite "lower" clock speeds than competing Intel/AMD chips.

### Quick Check
> 1. A processor runs at 4 GHz with IPC = 3.5. What is its performance in billion instructions per second (BIPS)?
> 2. AMD Zen 4 achieves ~4.1 IPC on integer code with a 6-wide superscalar. What fraction of its theoretical maximum is this?
> 3. Apple Firestorm achieves higher IPC than Intel/AMD despite a similar process node. Name two architectural reasons this is possible.

---

## 4. The ILP Wall — The Limits of Superscalar

Around 2002–2005, processor designers discovered that making superscalar CPUs wider didn't help much anymore — the IPC barely increased. This became known as the **ILP wall** (Instruction-Level Parallelism wall).

### Why ILP Is Limited

Real programs have data dependencies. Even with perfect hardware and no structural hazards, the instruction stream itself limits how many instructions can execute simultaneously. Consider:

```
int sum = 0;
for (int i = 0; i < N; i++) {
    sum += array[i];      // sum depends on previous sum — long chain of dependencies!
}
```

Every iteration depends on the previous one through `sum`. No amount of superscalar width helps here — it's a serial chain. This is **loop-carried dependency**, common in reduction operations.

### Amdahl's Law Applied to ILP

Even if 90% of a program is parallelizable, the serial 10% limits maximum speedup to 10×. With 8-wide superscalar, the realistic IPC ceiling is much lower than 8 for most workloads.

### What Replaced Wide Superscalar

After the ILP wall became clear, the industry pivoted to:
1. **Multiple cores** (thread-level parallelism instead of instruction-level)
2. **SIMD/vector execution** (exploit data-level parallelism with explicit instructions)
3. **Specialized accelerators** (GPUs, NPUs — exploit application-specific parallelism)

Modern CPUs still get wider (6–10 wide), but the gains are incremental. The big wins now come from better branch prediction, larger instruction windows, smarter prefetching, and specialized hardware.

### Quick Check
> 1. What is the ILP wall?
> 2. Why does a loop like `sum += array[i]` have poor ILP?
> 3. If data-level parallelism (SIMD) can achieve better throughput than instruction-level parallelism for number crunching, why don't CPUs just focus on SIMD and ignore ILP?

---

## 5. VLIW — The Compiler-Driven Alternative

**VLIW (Very Long Instruction Word)** takes a completely different approach to superscalar execution: instead of the hardware dynamically finding parallel instructions, the **compiler** does it statically and explicitly encodes multiple operations in one wide instruction.

```
VLIW instruction (example with 4 slots):
[  ALU op  ][  FPU op  ][ Load/Store ][  Branch  ]

The compiler fills each slot with independent operations.
If a slot can't be filled usefully, it gets a NOP.
```

**Advantages**:
- Simpler hardware (no reservation stations, no dynamic scheduling, no OOO)
- Lower power consumption
- Potentially higher clock frequency

**Disadvantages**:
- The compiler must find parallelism at compile time — much harder than hardware doing it at runtime
- Code doesn't run well on different VLIW designs (binary compatibility breaks)
- Cache misses are hard to handle statically (the compiler can't predict runtime latencies)
- Code size bloat from NOPs

**VLIW in practice**: TI C6000 DSPs (very successful for signal processing with regular, predictable code). Intel Itanium (VLIW for server CPUs — a famous commercial failure: very fast for benchmarks, terrible on real code). Most Itanium customers eventually migrated back to x86.

**EPIC (Explicitly Parallel Instruction Computing)**: Intel/HP's refinement of VLIW for Itanium — allows the compiler to indicate which instruction groups are independent. Failed commercially by 2017.

VLIW works well for DSPs and GPU shaders (regular, compiler-friendly code) but has never succeeded for general-purpose CPUs.

### Quick Check
> 1. What is the key difference between a superscalar CPU and a VLIW CPU in terms of who finds parallelism?
> 2. Why did the Itanium VLIW processor fail commercially despite being technically impressive?
> 3. Why does VLIW work well for DSPs but not for general-purpose processors?

---

## 6. SMT — Simultaneous Multithreading

Even an 8-wide superscalar CPU often has execution units sitting idle — one thread simply doesn't have enough ILP to fill all ports every cycle. **SMT (Simultaneous Multithreading)** solves this by running two (or more) independent threads simultaneously on the same physical core.

```
Without SMT: Thread A uses ports 0,2,5 — ports 1,3,4,6,7 idle
With SMT:    Thread A uses ports 0,2,5 — Thread B uses ports 1,3,4,6,7
```

Intel calls their SMT implementation **Hyper-Threading (HT)**. Most Intel consumer CPUs since Pentium 4 (2002) support 2-way SMT (2 logical cores per physical core). AMD calls theirs **SMT** as well, with Zen CPUs supporting 2-way SMT.

### What's Shared vs Separate in SMT

| Component | With SMT |
|-----------|----------|
| Register file | **Duplicated** (each thread has its own architectural registers, mapped to shared physical registers) |
| ROB | **Shared** (partitioned between threads) |
| Reservation Stations | **Shared** (both threads' instructions compete for slots) |
| Execution units | **Shared** (threads time-share execution ports) |
| L1 I-Cache & D-Cache | **Shared** |
| Branch predictor | **Shared** (some context per thread) |

The two logical threads appear to the OS as two separate CPUs. Thread scheduling, context switches, and process isolation all work normally.

**SMT performance gain**: Typically 15–30% throughput improvement on throughput workloads (web servers, compilers, databases) for essentially free hardware cost. However:
- Single-threaded performance does not improve (and may slightly decrease due to resource sharing)
- SMT introduces security concerns (cache timing side channels between sibling threads — the Spectre/Meltdown vulnerabilities are amplified by SMT)

IBM POWER CPUs support up to **8-way SMT** (8 simultaneous threads per core). This is useful for server workloads where many threads share a single fat core, but single-thread performance per thread is lower.

### Quick Check
> 1. What is the performance benefit of SMT, and why doesn't it help single-threaded performance?
> 2. Which hardware components are shared between two SMT threads on the same core?
> 3. Why does SMT introduce security concerns?

---

## Summary

- A **superscalar** processor issues multiple instructions per cycle to multiple execution units. Modern CPUs are 4–8-wide.
- **Execution units** are specialized: integer ALUs, FPUs, load/store units, branch units. Multiple units of each type feed through dispatch ports.
- **IPC** is the key metric; typical values are 3–5 for modern high-performance CPUs, well below the theoretical maximum.
- The **ILP wall** limits superscalar scaling — most real programs don't have enough independent instructions to fill 6-8 execution slots every cycle.
- **VLIW** uses compiler-scheduled static parallelism — works for DSPs but failed for general-purpose CPUs (Itanium).
- **SMT (Hyper-Threading)** fills idle execution slots by running a second thread on the same core — typically 15–30% throughput gain at minimal hardware cost.

---

## Exercises

### Easy
1. A 4-wide superscalar runs at 3 GHz. In the best case, what is its peak throughput in GIPS (billion instructions per second)? What is a realistic throughput for integer code?
2. List three reasons why a 4-wide superscalar achieves less than 4.0 IPC on real code.
3. What does "Hyper-Threading" actually do? How does the CPU appear to the operating system?

### Medium
4. A program has the following mix: 40% integer instructions (can issue 3/cycle max), 30% floating-point (can issue 2/cycle max), 20% loads/stores (can issue 2/cycle max), 10% branches (1/cycle max). What is the maximum sustainable IPC assuming perfect scheduling?
5. Consider two processors: Processor A runs at 4 GHz with IPC = 2.0; Processor B runs at 3 GHz with IPC = 4.0. For a workload that runs 10 billion instructions, calculate the execution time for each. Which is faster?
6. Explain why VLIW code compiled for a 4-slot VLIW machine will not run correctly on a 6-slot VLIW machine (even if the ISA is "compatible").

### Hard
7. SMT performance depends heavily on workload type. Analyze how SMT would behave for: (a) a memory-bound workload where one thread spends 80% of cycles waiting for cache misses, (b) a compute-bound workload where one thread uses 90% of all execution ports, (c) a workload of two threads that both frequently access the same data (high L1 cache sharing). For each case, estimate whether SMT helps, hurts, or has negligible effect.
8. The ILP wall showed that wider superscalar yields diminishing returns. Research "memory-level parallelism" (MLP) — the ability to have multiple cache misses outstanding simultaneously. Explain how MLP is different from ILP, why out-of-order processors can exploit MLP even when ILP is low, and why MLP is just as important as ILP for modern workloads.
