# Chapter 06: Memory from Logic — Flip-Flops and Registers

All the circuits we have built so far are **combinational**: the output depends only on the current inputs, with no memory of what came before. But a computer needs to remember things — the value of a variable between clock cycles, the current instruction, the state of the program. This chapter introduces **sequential logic**: circuits whose output depends on both current inputs AND past history.

## Table of Contents

1. [Combinational vs Sequential Logic](#1-combinational-vs-sequential-logic)
2. [The SR Latch — The Simplest Memory](#2-the-sr-latch--the-simplest-memory)
3. [The D Latch — Cleaner Memory](#3-the-d-latch--cleaner-memory)
4. [The Clock Signal](#4-the-clock-signal)
5. [The D Flip-Flop — Edge-Triggered Memory](#5-the-d-flip-flop--edge-triggered-memory)
6. [Registers — Multi-Bit Storage](#6-registers--multi-bit-storage)
7. [The Register File](#7-the-register-file)
8. [Counters — Registers That Count](#8-counters--registers-that-count)
9. [Shift Registers](#9-shift-registers)
10. [From Flip-Flops to SRAM](#10-from-flip-flops-to-sram)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Combinational vs Sequential Logic

**Combinational logic**: output is a pure function of current inputs. Change the input, the output immediately reflects it. No memory, no history.

```
Input ──► [Combinational Circuit] ──► Output
```

A full adder is combinational. Apply A=1, B=1 → Sum=0, Cout=1. Remove the inputs → output goes to 0 immediately. It "forgets" instantly.

**Sequential logic**: output depends on both current inputs AND past state. The circuit has memory — it "remembers" what happened before.

```
Input ──► [Sequential Circuit] ──► Output
                ▲         │
                └─────────┘
                  (feedback)
```

The key to sequential logic is **feedback**: the output is fed back as an input. This creates stable states that can hold information even after the input changes.

Think of a light switch with a "memory." You press it once → light turns on and STAYS on even after you release the button. Press again → light turns off and STAYS off. The circuit remembers its state. This is the essence of sequential logic.

### Quick Check

> 1. What is the fundamental difference between combinational and sequential logic?
> 2. How is feedback used in sequential logic circuits?
> 3. Give an example of a real-world device that is combinational and one that is sequential.

---

## 2. The SR Latch — The Simplest Memory

The **SR Latch** (Set-Reset Latch) is the simplest memory circuit. It can hold one bit of information indefinitely.

Built from two cross-coupled NAND gates:

```
S ──┐
    ├──[NAND]──── Q
    │        │
    │        └──────────────┐
    │                       │
    └──────────────┐        │
                   │        │
    ┌──────────────┘        │
    │                       │
R ──┤                       │
    ├──[NAND]──── Q̄ (NOT Q)
    │        │
    │        └──────────────┘
```

### How the SR Latch Works

**Set (S=0, R=1)**: Force Q to 1.
- The NAND gate with S input: 0 NAND X = 1 → Q = 1
- Q (now 1) feeds back to the other NAND gate
- The other NAND gate: R(=1) NAND Q(=1) = 0 → Q̄ = 0
- The latch is "set": Q=1, Q̄=0

**Reset (S=1, R=0)**: Force Q to 0.
- The NAND gate with R input: 0 NAND X = 1 → Q̄ = 1
- Q̄ (now 1) feeds back
- Other NAND gate: S(=1) NAND Q̄(=1) = 0 → Q = 0
- The latch is "reset": Q=0, Q̄=1

**Hold (S=1, R=1)**: Remember current state.
- Both NAND gates have one input = 1, so each output depends on the other's output
- The circuit is in a stable state — whatever Q was, it stays that way
- This is the "hold" or "store" state

**Forbidden (S=0, R=0)**: Invalid state.
- Both gates output 1 → Q=1 and Q̄=1 simultaneously (contradictory)
- When you release this state, which value does it settle to? Unpredictable — this is a hazard

```
Truth Table (NAND SR Latch):
┌─────┬─────┬─────────────────────────┐
│  S  │  R  │  Action                  │
├─────┼─────┼─────────────────────────┤
│  0  │  1  │  Set: Q=1               │
│  1  │  0  │  Reset: Q=0             │
│  1  │  1  │  Hold: Q unchanged      │
│  0  │  0  │  FORBIDDEN              │
└─────┴─────┴─────────────────────────┘
(Note: active-low NAND latch, so 0 means "active")
```

### Where SR Latches Appear

SR latches are used in debouncing mechanical switches. When a switch clicks, its contacts bounce 5-20 times before settling — a naive circuit would register many keypresses. An SR latch's stable states mean it captures the first clean signal and ignores the subsequent bounces.

### Quick Check

> 1. What does "Set" do to the SR latch's output?
> 2. What is the "forbidden" state and why is it forbidden?
> 3. In the "Hold" state, does the SR latch need any input to maintain its value?

---

## 3. The D Latch — Cleaner Memory

The SR latch has a problem: the forbidden state. The **D Latch** (Data Latch) solves this by having a single data input D and a control input Enable (EN):

```
D ──────┬──[AND]──── S ──┐
        │                 ├── SR Latch ── Q
NOT(D) ─┴──[AND]──── R ──┘
                    ↑
                   EN
```

When **EN=1 (transparent mode)**: Q follows D. If D=1, Q=1. If D=0, Q=0. The latch is "transparent" — data flows through.

When **EN=0 (latched mode)**: Q holds its last value. Changes to D have no effect.

The D latch eliminates the forbidden state: S and R can never both be 0 (because S=D AND EN, and R=NOT(D) AND EN — D and NOT(D) cannot both be 1 simultaneously).

### The Level-Triggered Problem

The D latch captures data whenever EN is high (transparent). This means data can change any time during the high phase — the "data window" is wide. In a complex circuit, this creates timing problems: if something changes D during the middle of the EN=1 phase, the latch catches an intermediate value.

The solution: edge-triggered flip-flops.

### Quick Check

> 1. How does the D latch prevent the forbidden state?
> 2. When EN=0, what happens to Q if D changes?
> 3. What is the "transparent mode" of a D latch?

---

## 4. The Clock Signal

Before we introduce flip-flops, we need to understand **clocks** — the heartbeat of a digital system.

A clock is a signal that alternates between 0 and 1 at a fixed frequency:

```
     ┌──┐  ┌──┐  ┌──┐  ┌──┐
  1  │  │  │  │  │  │  │  │
     │  │  │  │  │  │  │  │
  0 ─┘  └──┘  └──┘  └──┘  └──
  
     ↑  ↑  ↑  ↑  ↑  ↑  ↑  ↑
     Rising edges (0→1 transitions)
     
     Period = time for one complete cycle
     Frequency = 1 / Period (in Hertz)
```

A clock running at 3 GHz means 3,000,000,000 cycles per second — the signal rises and falls 3 billion times per second.

### Why Clocks?

Clocks are the great synchronizer of digital logic. Every flip-flop in the entire CPU updates at the same moment — when the clock edge occurs. This ensures that:

1. All circuits are stable before data is sampled (no partial values)
2. Data flows through one stage of logic per clock cycle in a predictable way
3. The CPU's state is well-defined at every clock edge

Without a clock, different parts of the circuit would change at different times due to different propagation delays, leading to chaos.

### Frequency and Period

```
Clock frequency:  3 GHz = 3 × 10⁹ Hz
Clock period:     1 / (3 × 10⁹) = 0.33 nanoseconds

In each 0.33ns clock period, the CPU must complete one or more pipeline stages.
Logic gates are measured in picoseconds (1ps = 10⁻¹²s).
A typical gate delay: 20-100ps.
0.33ns = 330ps → about 3-16 gate delays per clock cycle.
```

### Quick Check

> 1. A CPU runs at 4 GHz. What is the clock period in nanoseconds?
> 2. Why is having all flip-flops update at the same time important?
> 3. What is a "rising edge" of a clock signal?

---

## 5. The D Flip-Flop — Edge-Triggered Memory

The **D Flip-Flop** (DFF) is the most important memory element in digital design. It captures the value of its D input at the instant of a clock edge, and holds that value until the next clock edge.

```
D ──────┐
        │  D Flip-Flop
        ├────────────── Q  (captured on clock edge)
Clock ──┘
```

**On the rising edge of the clock (0→1 transition)**: Q takes the current value of D.
**At all other times**: Q holds its value, regardless of changes to D.

This is called **edge-triggered** because capture happens only at the transition (edge) of the clock, not during the high or low phase.

```
Timing Diagram:

Clock: ___|‾|_|‾|_|‾|_|‾|_
D:     _____|‾‾‾‾‾|___|‾‾‾
Q:     _________|‾‾‾|___|‾  ← Q captures D value at each rising clock edge
```

### Why Edge-Triggered Is Better

The latch is transparent for half the clock cycle — a large window where data can "leak through" at the wrong time. The flip-flop captures at one instant (the edge) — a tiny window that ensures data is stable before capture.

This is why real CPUs use D flip-flops (or JK flip-flops, T flip-flops), never latches for state storage.

### Building a D Flip-Flop

A D flip-flop is built from two D latches in a master-slave configuration:

```
D ──► [Master D Latch (EN=NOT CLK)] ──► [Slave D Latch (EN=CLK)] ──► Q
```

When CLK=0: Master latch is transparent (captures D), Slave is latched (holds output)
When CLK=1: Master latch is latched (holds), Slave becomes transparent (passes master's output to Q)

The output Q only changes at the rising edge (0→1), when Slave latch opens. Perfect timing control.

### Quick Check

> 1. When does a rising-edge D flip-flop update its output Q?
> 2. What is the difference between a D latch and a D flip-flop?
> 3. Why is edge-triggered behavior important in a CPU?

---

## 6. Registers — Multi-Bit Storage

A **register** holds multiple bits simultaneously. An 8-bit register is simply 8 D flip-flops connected in parallel — they all share the same clock signal, so all 8 bits update simultaneously.

```
8-bit Register:

D[7] ──► DFF ──► Q[7]
D[6] ──► DFF ──► Q[6]
D[5] ──► DFF ──► Q[5]   ← all 8 flip-flops share
D[4] ──► DFF ──► Q[4]      the same clock signal
D[3] ──► DFF ──► Q[3]
D[2] ──► DFF ──► Q[2]
D[1] ──► DFF ──► Q[1]
D[0] ──► DFF ──► Q[0]

Clock ──────────────── (connected to all 8 DFFs)
```

On each rising clock edge, all 8 output bits Q[7:0] simultaneously capture their corresponding D[7:0] input values.

Registers are the fastest storage in a computer — they sit directly inside the processor chip, just a few gate delays from the ALU. But they are few in number (a typical CPU has 16-32 general-purpose registers) because each flip-flop uses 6-10 transistors and fast flip-flop circuits are area-expensive.

### Key CPU Registers

Every CPU has special-purpose registers:

- **PC (Program Counter)**: holds the address of the next instruction to execute
- **IR (Instruction Register)**: holds the currently executing instruction
- **SP (Stack Pointer)**: holds the address of the top of the call stack
- **FLAGS / PSW (Program Status Word)**: holds the condition codes (Zero, Carry, Negative, Overflow flags)
- **General-Purpose Registers (R0-R31 in RISC-V, RAX-R15 in x86-64)**: hold temporary values during computation

### Quick Check

> 1. How many D flip-flops are in a 32-bit register?
> 2. When does a 64-bit register update all its bits?
> 3. What is the Program Counter (PC) and what does it hold?

---

## 7. The Register File

A CPU needs multiple registers and must be able to read two and write one in every clock cycle (to feed the two operands of an instruction and store the result). A **register file** is the organized collection of registers with the multiplexing circuitry to access any of them.

A 32-register, 64-bit register file (like in RISC-V):

```
Read Port 1 address (5 bits) ──► [5-to-32 Decoder] ──► Selects one of 32 registers
Read Port 2 address (5 bits) ──► [5-to-32 Decoder] ──► Selects one of 32 registers
Write address (5 bits)       ──► [5-to-32 Decoder] ──► Selects destination register
Write Enable                 ──► (enables write to selected register)
Write Data (64 bits)         ──► (data to store)
                                 
Read Data 1 (64 bits)        ◄── (value from selected register 1)
Read Data 2 (64 bits)        ◄── (value from selected register 2)
```

The register file has three "ports": two read ports and one write port. In every clock cycle, it can simultaneously:
- Read the value of any two registers (for the instruction's source operands)
- Write the result of the previous instruction to any register

Multi-ported register files are expensive in silicon area — each additional port roughly doubles the area. This is why CPU register counts are limited.

### Quick Check

> 1. In RISC-V, there are 32 general-purpose registers. How many bits are needed to specify which register to read or write?
> 2. What does the write enable signal do in a register file?
> 3. Why might a CPU designer want MORE register ports, and what limits how many you can have?

---

## 8. Counters — Registers That Count

A **counter** is a register with feedback logic that makes it increment (or decrement) on each clock edge. The simplest: a binary counter built from T flip-flops.

A **T flip-flop** (Toggle flip-flop) flips its output every time a clock edge arrives (if T=1). Chain several together:

```
4-bit binary counter:

Q[0] flips every clock edge        (toggles every 1 cycle)
Q[1] flips when Q[0] goes 0→1     (toggles every 2 cycles)
Q[2] flips when Q[1] goes 0→1     (toggles every 4 cycles)
Q[3] flips when Q[2] goes 0→1     (toggles every 8 cycles)

Counting sequence:
0000, 0001, 0010, 0011, 0100, 0101, 0110, 0111,
1000, 1001, 1010, 1011, 1100, 1101, 1110, 1111, 0000, 0001...
```

**Counter uses in CPUs**:
- The **Program Counter (PC)** is a counter that increments by 4 (the instruction size in bytes) on each fetch
- Hardware timers/watchdogs are counters
- Performance counters (counting cache misses, branch mispredictions)

### Quick Check

> 1. A 4-bit counter has been running and currently shows 1111. What will it show after the next clock edge?
> 2. Why does the PC in a 32-bit RISC CPU typically increment by 4, not by 1?

---

## 9. Shift Registers

A **shift register** stores n bits but moves (shifts) them one position left or right on each clock edge. It is a chain of D flip-flops where each flip-flop's output feeds the next one's input:

```
Serial input ──► [DFF₀] ──► [DFF₁] ──► [DFF₂] ──► [DFF₃] ──► Serial output
                   ↓           ↓           ↓           ↓
                  Q[0]        Q[1]        Q[2]        Q[3]
                 (parallel output)
```

Each clock edge shifts all bits one position to the right, and a new bit enters from the left.

**Shift register uses**:
- Serial-to-parallel conversion: receive data one bit at a time (e.g., from a serial port) and assemble it into a parallel word
- Parallel-to-serial conversion: send data bit by bit
- Arithmetic shifts: shifting a binary number left by n = multiplying by 2ⁿ; shifting right by n = dividing by 2ⁿ
- CRC (Cyclic Redundancy Check) computation for error detection

**CPU shift operations (SHL, SHR, SAR)**:
- SHL (shift left): multiply by 2 per shift, fast alternative to multiplication
- SHR (shift right logical): divide by 2, fills with zeros from the left
- SAR (shift right arithmetic): divide by 2 for signed numbers, fills with the sign bit

### Quick Check

> 1. If a shift register contains 1010 and you shift left by 1, what does it contain?
> 2. Multiplying by 8 using shift operations: shift left by how many positions?
> 3. What is the difference between SHR and SAR for the value 11110000?

---

## 10. From Flip-Flops to SRAM

The flip-flops we have discussed are how CPU registers work. But a CPU also needs more storage than a few dozen registers — it needs cache memory (kilobytes to megabytes) that must also be fast.

**SRAM** (Static Random Access Memory) is built from groups of 6 transistors per cell — effectively a cross-coupled inverter pair (like a latch) plus 2 access transistors. No refresh is needed (unlike DRAM) because the cross-coupled inverters actively hold their state as long as power is applied.

```
SRAM Cell (6-transistor, 6T):

        VDD
         │
    ┌───[P1]───┬───[P2]───┐
    │          │          │
    │   [N3]──►│          │◄──[N4]
    │          │          │
    ├───[N1]───┴───[N2]───┤
         │                │
        GND               GND
         
    BL (Bit Line)  ──── access via N3 ──── BL̄ (NOT Bit Line)
    WL (Word Line): activates N3 and N4 to connect cell to bit lines
```

L1 cache in a modern CPU is SRAM. It is fast (4-cycle access), but expensive in area (6 transistors per bit vs 1 transistor per bit for DRAM).

**DRAM** (Dynamic RAM) uses 1 transistor + 1 capacitor per bit. Much denser, but the capacitor leaks charge and must be "refreshed" every few milliseconds. This is main memory (your laptop's 16 GB RAM is DRAM).

| Type | Transistors/bit | Speed | Power | Density | Usage |
|------|----------------|-------|-------|---------|-------|
| Flip-flop register | ~6-10 | fastest | high | lowest | CPU registers |
| SRAM | 6 | fast | medium | low | CPU cache |
| DRAM | 1+1 cap | slower | low | high | Main memory |
| Flash | 1 | slowest | very low | highest | SSD storage |

### Quick Check

> 1. How many transistors does an SRAM cell use?
> 2. What is the main difference between SRAM and DRAM?
> 3. Why is DRAM used for main memory instead of SRAM, even though SRAM is faster?

---

## Summary

- **Combinational logic** produces outputs based only on current inputs. **Sequential logic** remembers past history through feedback.
- The **SR Latch** is the simplest memory circuit — two cross-coupled NAND gates that hold one bit. The forbidden state (both inputs active) is a design hazard.
- The **D Latch** improves the SR latch: one data input D and an Enable signal, no forbidden state. Transparent when Enable=1.
- The **clock signal** synchronizes all state changes in a CPU, ensuring stable timing.
- The **D Flip-Flop** captures D only at the clock edge — more reliable than a latch. It is the fundamental memory element in CPUs.
- An **n-bit register** is n D flip-flops sharing a clock, holding an n-bit word.
- The **register file** provides multiple registers with read/write ports, accessible by address.
- **Counters** increment on each clock edge — the PC is a counter.
- **Shift registers** chain flip-flops to shift bits, enabling fast multiplication/division and serial communication.
- **SRAM** (6T cell) is used for cache; **DRAM** (1T+1C) is used for main memory. Speed-density-power tradeoffs determine where each type is used.

---

## Exercises

### Easy

1. Trace through an SR latch starting in the Reset state (Q=0). Apply S=0, R=1 (Set). Then apply S=1, R=1 (Hold). What is Q at each step?

2. A D flip-flop has D=0 when the rising clock edge occurs. What happens to Q?

3. An 8-bit binary counter currently holds 11111110. After one clock edge, what does it hold?

### Medium

4. Draw a 4-bit shift register using D flip-flops. Starting state: 1010. Apply 3 clock edges. What is the state after each edge if the serial input is always 0?

5. A 32-register file needs to read two registers and write one per cycle. How many decoder circuits does it need? How many flip-flops (for 64-bit registers)?

6. Explain why the master-slave configuration of D flip-flops (master latch + slave latch) ensures that Q only changes at the rising clock edge and not during the high phase of the clock.

### Hard

7. Design a modulo-6 counter (a counter that counts 0, 1, 2, 3, 4, 5, 0, 1, ...). Start with a 3-bit binary counter (counts 0-7) and add combinational logic that detects the count of 6 (110 in binary) and resets the counter to 0 when it reaches 6. Draw the circuit. This is how ring counters in CPUs and clock dividers work.

8. **SRAM vs DRAM**: Research the refresh mechanism of DRAM. Why does DRAM need refreshing while SRAM does not? If DRAM must be refreshed every 64ms and a memory module has 8 billion rows, how many rows must be refreshed per second? At 2 rows refreshed per nanosecond, what fraction of memory bandwidth is consumed by refresh? How does this affect system performance?
