# Chapter 64: SHAKTI Processor Families — From Embedded to Server

This chapter provides a detailed look at each SHAKTI processor class — its specifications, target applications, tape-out history, and the actual chips that have been produced. Understanding the SHAKTI family shows both the technical achievements and the honest gaps between academic research goals and commercial production requirements.

## Table of Contents

1. [E-Class — The Embedded Workhorse](#1-e-class--the-embedded-workhorse)
2. [C-Class — The Linux-Capable Processor](#2-c-class--the-linux-capable-processor)
3. [I-Class — Industrial Performance](#3-i-class--industrial-performance)
4. [M-Class — Mobile Ambitions](#4-m-class--mobile-ambitions)
5. [S-Class — The Server Target](#5-s-class--the-server-target)
6. [SHAKTI SoC — System Integration](#6-shakti-soc--system-integration)
7. [Performance Summary and Roadmap](#7-performance-summary-and-roadmap)
8. [Exercises](#exercises)

---

## 1. E-Class — The Embedded Workhorse

**E-class** (Embedded class) is the simplest SHAKTI processor, designed for microcontroller applications — the equivalent of an ARM Cortex-M0 or M3.

**Specifications:**
- ISA: RV32IMAC (32-bit, Integer + Multiply + Atomic + Compressed)
- Pipeline: 3-stage in-order (Fetch → Decode+Execute → Memory+Writeback)
- Frequency: 100–200 MHz at TSMC 22nm
- No data cache (uses tightly-coupled SRAM — deterministic timing)
- Instruction memory: accessed directly (optional small cache)
- Area: ~15,000–50,000 gate equivalents
- Power: 0.5–5 mW at typical application workloads

**Applications:**
- IoT sensor nodes (temperature, humidity, motion sensors)
- Microcontroller replacement for government/strategic systems where ARM license is undesirable
- Educational boards (SHAKTI E-class based FPGA development kits)
- ISRO small satellite systems (for attitude control, sensor interfaces)

**Tape-out**: E-class variants have been taped out at TSMC 180nm (older process, easy to work with), and later at TSMC 22nm via Intel partnership.

**Comparison to ARM Cortex-M3:**
- Cortex-M3 at 120 MHz: ~120 MIPS, Dhrystone 1.25 DMIPS/MHz
- SHAKTI E-class at 150 MHz: ~0.6 DMIPS/MHz → ~90 MIPS effective
- Performance comparable; difference is free license, open source, and security extensions

```
E-class pipeline:
  
  Stage 1: Fetch       Read instruction from ITCM (instruction TCM)
  Stage 2: Decode+Exec Decode instruction, execute ALU, compute address
  Stage 3: Mem+WB      Access data memory, write result to register file
  
  Simple design: no pipeline hazards requiring stalls
  (except memory accesses: 1 stall cycle for load-use hazard)
```

### Quick Check
> 1. What is TCM (Tightly Coupled Memory) and why does the E-class use it instead of a cache?
> 2. What RISC-V ISA extensions does the E-class implement?
> 3. Compare E-class performance to ARM Cortex-M3.

---

## 2. C-Class — The Linux-Capable Processor

**C-class** (Commercial class) is the flagship processor of the SHAKTI family — the first to run Linux and the most prominently deployed.

**Specifications:**
- ISA: RV64IMAC (64-bit) or RV64IMAFDC (with floating point)
- Pipeline: 5-stage in-order (Fetch/Decode/Execute/Memory/Writeback)
- Caches: 16KB–32KB I-cache, 8KB–16KB D-cache (4-way set-associative)
- Optional L2 cache: 256KB–1MB
- Branch predictor: Gshare predictor (2-bit saturating counters, global history)
- Hardware page table walker (for Linux virtual memory)
- Frequency: 400–600 MHz at TSMC 22nm
- Area: ~150,000–250,000 gate equivalents (without caches)
- Power: 50–150 mW active

**Linux boot**: C-class successfully boots upstream Linux kernel (kernel 5.10+). The Linux port uses:
- RISC-V SBI (Supervisor Binary Interface) for firmware services
- Device drivers for SHAKTI SoC peripherals
- U-Boot bootloader

**Deployed applications:**
- ISRO satellite ground station processors (replacing foreign processors in control systems)
- DRDO (Defence Research and Development Organisation) evaluation for defense electronics
- Indian government-mandated "secure computing" applications
- IIT Madras student research and education

**The Mindgrove SC23 SoC**: Mindgrove Technologies (IIT Madras spinoff) built a commercial SoC using a SHAKTI C-class core:
- 4× SHAKTI C-class cores
- Neural network accelerator
- 32KB L2 cache
- UART, SPI, I2C, GPIO peripherals
- TSMC 22nm process
- Target: industrial IoT, edge inference

```
C-class pipeline stages (detailed):
  
  Fetch:    PC → I-cache lookup → 64-bit fetch buffer
            Branch predictor updates PC speculatively
  
  Decode:   Decode RISC-V instruction
            Read two source registers from register file
            Compute branch target (PC-relative)
  
  Execute:  ALU computes result
            Branch condition evaluated → update/correct PC
            Load/store address computed
  
  Memory:   D-cache access for load/store
            TLB lookup for virtual memory (C-class has MMU)
  
  Writeback: Write result to register file
             Commit instruction state
```

### Quick Check
> 1. What is the Mindgrove SC23 SoC and what processor does it use?
> 2. What is the RISC-V SBI and why is it needed to boot Linux?
> 3. What is a hardware page table walker and why is it required for Linux?

---

## 3. I-Class — Industrial Performance

**I-class** (Industrial class) targets applications needing more performance than C-class — factory automation, networked industrial controllers, medical devices.

**Specifications:**
- ISA: RV64IMAFDC + V (vector extension, in development)
- Pipeline: 7-stage in-order pipeline
- Dual-issue capability (limited): can issue two arithmetic instructions per cycle
- Larger caches: 32KB I-cache, 32KB D-cache, 512KB L2
- FPU: IEEE 754 single + double precision, hardware divide
- Frequency: 600–800 MHz at TSMC 22nm
- Power: 100–300 mW

**RISC-V Vector Extension (V):**
The RISC-V V extension adds vector instructions (SIMD) with configurable vector lengths. Unlike ARM Neon (fixed 128-bit) or AVX-512 (fixed 512-bit), RISC-V V uses a "vector length agnostic" (VLA) model:
- `vl` register specifies how many vector elements are processed this instruction
- Hardware can choose to process 128, 256, 512, or more bits per instruction
- Software written for one vector length works on any RISC-V V hardware

**I-class status (2024)**: I-class is more advanced research than C-class. Partial implementations exist; full silicon tape-out is targeted but not yet produced as of 2023. The vector extension adds significant verification complexity.

### Quick Check
> 1. What is the RISC-V V (vector) extension and what is unique about its "vector length agnostic" model?
> 2. How is the I-class different from C-class in terms of pipeline complexity?
> 3. What is the target application domain for I-class?

---

## 4. M-Class — Mobile Ambitions

**M-class** (Mobile class) is SHAKTI's attempt at out-of-order execution — comparable to ARM Cortex-A55 class performance.

**Specifications (design target):**
- ISA: RV64IMAFDC + RISC-V P (Packed SIMD, similar to ARM Neon)
- Pipeline: Out-of-order, 2–4 issue width
- ROB (Reorder Buffer): 32–64 entries
- Caches: 32KB L1-I, 32KB L1-D, 512KB–1MB L2
- Frequency target: 1–1.5 GHz at TSMC 28nm
- Power: 500–1000 mW

**Out-of-order in SHAKTI:**
M-class implements the standard OOO pipeline: Fetch → Decode → Rename → Issue Queue → Execute → Commit (ROB). This is significantly more complex than the in-order C-class:
- Register renaming (physical register file, 64+ physical registers)
- Load-store queue for memory order disambiguation
- Multiple execution units in parallel (integer, FPU, load/store)

**Status (2024)**: M-class is in research/design phase. A partial RTL exists but silicon tape-out has not been announced. The team acknowledges this is the most complex processor in the SHAKTI family.

**SHAKTI's strategy**: Rather than rushing a mediocre M-class, the team is focusing on:
1. Perfecting C-class for production deployment
2. Using Mindgrove to commercialize C-class-based SoCs
3. Developing M-class as a research platform for India's CPU design talent pipeline

### Quick Check
> 1. What is the ROB (Reorder Buffer) in an out-of-order processor?
> 2. Why is M-class development more complex than C-class?
> 3. What is SHAKTI's deployment strategy while M-class is under development?

---

## 5. S-Class — The Server Target

**S-class** (Server class) is the long-term goal — a RISC-V server processor competitive with ARM Neoverse N1 or Intel Xeon at the low end.

**Target specifications:**
- ISA: RV64IMAFDC + V (vector) + custom security extensions
- Pipeline: Superscalar OOO, 4–6 issue width
- ROB: 192+ entries
- Caches: 64KB L1-I, 64KB L1-D, 2–4MB L2, 16–32MB L3
- Multi-core: 2–8 cores on one chip
- Frequency target: 2+ GHz
- Process: TSMC 22nm or newer

**Challenges ahead:**
1. **Multi-core coherence**: S-class needs cache coherence protocol (MESI or similar). IIT Madras researchers have studied this but implementation is complex.
2. **Memory consistency model**: RISC-V uses a weak memory model (RVWMO). Server software (databases, virtualization) must be carefully ported.
3. **I/O subsystem**: PCIe, SATA, Ethernet, DDR4 controllers — all major projects individually.
4. **Power management**: Complex P-state and C-state management for power efficiency.

**Timeline**: Honest assessment from the SHAKTI team (2023): S-class silicon is a 5–10 year project. The focus for the next 3 years is stabilizing C-class deployment and commercializing via Mindgrove.

### Quick Check
> 1. What is the target frequency and ISA width for SHAKTI S-class?
> 2. What is the RISC-V weak memory ordering model (RVWMO)?
> 3. What are the two biggest technical challenges for S-class development?

---

## 6. SHAKTI SoC — System Integration

The processor core is only part of the system. SHAKTI SoCs integrate peripherals for real-world deployment:

**SHAKTI SoC (reference design):**
```
SHAKTI SoC block diagram:
  
  ┌─────────────────────────────────────────────┐
  │                SHAKTI SoC                    │
  │                                              │
  │  ┌──────────┐  ┌──────────┐                 │
  │  │ SHAKTI   │  │ SHAKTI   │  (C-class cores) │
  │  │ Core 0   │  │ Core 1   │                 │
  │  └────┬─────┘  └────┬─────┘                 │
  │       │             │                        │
  │  ┌────┴─────────────┴──────┐                 │
  │  │      AXI4 Interconnect  │                 │
  │  └──┬─────┬──────┬──────┬──┘                 │
  │     │     │      │      │                    │
  │  ┌──┴──┐ ┌┴───┐ ┌┴───┐ ┌┴───┐               │
  │  │SRAM │ │UART│ │SPI │ │I2C │               │
  │  │512KB│ │    │ │    │ │    │               │
  │  └─────┘ └────┘ └────┘ └────┘               │
  │                                              │
  │  DDR controller, GPIO, Debug (JTAG), PLL     │
  └─────────────────────────────────────────────┘
```

**Software stack:**
- SHAKTI SDK: bare-metal C library, startup code, peripheral drivers
- U-Boot (bootloader): standard U-Boot with SHAKTI patches
- Linux kernel: upstream RISC-V kernel with SHAKTI-specific device drivers
- Yocto Linux distribution: full embedded Linux for SHAKTI

**FPGA prototyping**: All SHAKTI processors are available as FPGA bitstreams for Xilinx and Intel FPGA boards — students and researchers can run SHAKTI on affordable FPGA hardware before committing to ASIC.

### Quick Check
> 1. What peripherals does a typical SHAKTI SoC include?
> 2. What is U-Boot and why is it needed?
> 3. How can students experiment with SHAKTI without access to custom ASIC?

---

## 7. Performance Summary and Roadmap

```
SHAKTI Processor Family Summary (2024):

Class   ISA       Pipeline    Freq       Power   Status
E       RV32IMAC  3-stage     200MHz     2mW     Silicon (22nm, Intel tape-out)
C       RV64IMAC  5-stage     500MHz     100mW   Silicon + production (SC23 SoC)
I       RV64IMAFDC+V 7-stage  800MHz     300mW   RTL complete, tape-out planned
M       RV64IMAFDC  OOO 2-4W  1GHz       800mW   Research/partial RTL
S       RV64IMAFDC+V OOO 4-6W  2GHz      5W      Design phase

Tape-outs completed:
  2017: Riscy (E-class precursor) at Intel
  2018: SHAKTI C-class 64-bit, TSMC 22nm
  2021: SHAKTI C-class with PMP, security extensions
  2022: SHAKTI-T (tagged memory), IIT-Intel Research
  2023: Mindgrove SC23 (C-class based SoC), TSMC 22nm
```

**Roadmap:**
- 2024–2025: I-class tape-out, vector extension validation
- 2025–2027: M-class tape-out, SoC for ISRO advanced satellites
- 2027–2030: S-class design completion, multi-chip SMP system
- 2030+: Potential India-fab SHAKTI on domestically produced silicon (if India Semiconductor Mission succeeds)

---

## Exercises

### Easy
1. Which SHAKTI class can run Linux and what pipeline does it use?
2. What is the Mindgrove SC23 SoC?
3. How can students run SHAKTI without custom hardware?

### Medium
4. SHAKTI E-class for ISRO: An ISRO small satellite needs 3 microcontrollers: attitude control (real-time, 50 MIPS minimum, <5mW), telemetry (IoT, 10 MIPS, <2mW), command & control (secure, supports crypto, <10mW). (a) Which SHAKTI class fits each application? (b) What RISC-V security extension (PMP) is critical for command & control? (c) Total power for all 3 processors? (d) If ARM Cortex-M4 fits the same requirements: what is the licensing cost and strategic risk vs SHAKTI?
5. C-class vs Cortex-A53 application choice: A startup building a smart factory IoT gateway chip needs: 200 MIPS processing, run embedded Linux, connect to industrial Ethernet (EtherCAT), <200mW power. (a) Can SHAKTI C-class at 500 MHz (IPC ~0.7) meet the 200 MIPS requirement? (b) What is the estimated die area for C-class with 32KB caches at TSMC 22nm? (c) ARM Cortex-A55 licenses at ~$1M: at 100,000 units, per-chip license cost? (d) Which is cheaper for 100,000 units, and what quality/ecosystem trade-offs exist?
6. SHAKTI roadmap analysis: The SHAKTI team targets M-class at 1 GHz by 2025. (a) What is the performance gap between M-class (1 GHz, IPC ~1.0) and ARM Cortex-A55 (1.8 GHz, IPC ~1.3)? (b) If SHAKTI closes 50% of the gap every 3 years, when does it match Cortex-A55 performance? (c) But ARM is also improving: Cortex-A720 at 3.5+ GHz. Is "catching ARM" a realistic goal, or should SHAKTI define its own success metrics?

### Hard
7. Multi-core SHAKTI cache coherence: The SHAKTI S-class target includes 8 cores with shared L3 cache. (a) What cache coherence protocol (MESI or MOESI) would you implement and why? (hint: SHAKTI's tagged memory adds complexity) (b) When two cores share a cache line with different tag annotations, what coherence rule must be followed? (c) RISC-V uses a relaxed memory model (RVWMO). What fence instructions must the coherence protocol interact with? (d) Design a 2-core test scenario that would trigger a cache coherence bug if the protocol is implemented incorrectly.
8. SHAKTI competitiveness plan: As an advisor to the SHAKTI program, design a 5-year strategic plan: (a) Which market segments should SHAKTI target immediately (2024–2026) given current C-class performance? (b) What capability must be demonstrated to win ISRO's next generation satellite processor contract? (c) If Mindgrove needs to grow from 0 to 1M chips shipped by 2028: what SoC design, manufacturing partner, and distribution channel is needed? (d) What is the "minimum viable SHAKTI" that could be deployed in 1 million Indian government/defense applications, and what is the economic impact?
