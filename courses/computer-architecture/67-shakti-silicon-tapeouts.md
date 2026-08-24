# Chapter 67: SHAKTI Silicon Tapeouts

The true test of any chip design is silicon — physical chips that work in the real world. SHAKTI has progressed from paper designs to silicon tape-outs to deployed systems. This chapter documents the actual chips that have been produced, the tape-out milestones, the test results, and the lessons learned from going from BSV code to working silicon.

## Table of Contents

1. [The Tape-Out Journey](#1-the-tape-out-journey)
2. [Riscy — The First Chip (2017)](#2-riscy--the-first-chip-2017)
3. [SHAKTI C-Class Silicon (2018)](#3-shakti-c-class-silicon-2018)
4. [SHAKTI-T Security Processor (2021–2022)](#4-shakti-t-security-processor-20212022)
5. [Mindgrove SC23 SoC (2023)](#5-mindgrove-sc23-soc-2023)
6. [Test Results and Performance Data](#6-test-results-and-performance-data)
7. [Lessons Learned](#7-lessons-learned)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Tape-Out Journey

Moving from RTL (Verilog/BSV code) to working silicon involves:
1. Design completion (RTL freeze)
2. Synthesis and place-and-route
3. Sign-off (DRC, LVS, timing)
4. Tape-out (GDSII to foundry)
5. Fabrication (6–8 weeks)
6. Packaging
7. Test and bring-up

For SHAKTI, each tape-out is a milestone — proving that the design concepts (security extensions, RISC-V compliance) actually work in physical silicon. Academic chip design programs sometimes produce chips that only partially work; SHAKTI's documented chips show increasing maturity.

**The Intel-IIT Madras partnership**: Intel established a research center at IIT Madras (IITM-Intel Research Center) that provided access to Intel's 22nm FinFET process and EDA tools for academic tape-outs. This was crucial for SHAKTI's early silicon — without it, the cost of tape-out ($1–5M for a full run) would have been prohibitive for an academic program.

**Mindgrove as commercialization vehicle**: Tape-outs at TSMC (through Mindgrove) represent the transition from academic to commercial quality. Mindgrove's SC23 was taped out through commercial channels with production-quality design constraints.

### Quick Check
> 1. What was the role of the Intel-IIT Madras Research Center in SHAKTI's early tape-outs?
> 2. Why is physical silicon necessary to validate a chip design?
> 3. What is Mindgrove's role in SHAKTI's commercialization?

---

## 2. Riscy — The First Chip (2017)

**Riscy** was the first chip produced by the IIT Madras SHAKTI team — technically a precursor to SHAKTI E-class, designed primarily as a learning exercise and proof of concept.

**Specifications:**
- ISA: RV32IM (minimal RISC-V: integers + multiply, no compressed)
- Pipeline: 3-stage in-order
- No cache (SRAM interface directly)
- Process: Intel 22nm FinFET (Intel's IITM Research Center access)
- Area: ~2mm²
- Frequency: ~100 MHz

**What was demonstrated:**
- First time an IIT Madras team completed the full RTL → silicon → test cycle
- RISC-V compliance: passed basic RISC-V compliance tests
- Linux: NOT bootable (missing MMU, interrupt controller)
- Primary purpose: validate the design flow, build team experience

**Technical approach (BSV → Verilog → synthesis):**
The SHAKTI team developed their BSV-to-Verilog synthesis flow using the Bluespec compiler. This flow was new and required debugging tool integration issues. Riscy validated that BSV-generated Verilog was synthesizable and timing-closeable at 22nm.

**Outcome**: Riscy worked but with limited functionality. More importantly, the tape-out experience taught the team:
- BSV synthesis quirks at 22nm
- Intel 22nm library characterization
- Post-silicon debug methodology

```
Riscy comparison to commercial parts:
  
  ARM Cortex-M0 (competitor):
  - 0.4mm² at 40nm, 40+ DMIPS at 48 MHz, 98µW/MHz
  
  Riscy (SHAKTI precursor):
  - ~2mm² at 22nm, ~50 DMIPS at 100 MHz, higher power
  - 5× larger area, lower efficiency
  - Expected for a first academic design
  
  Gap explained by: no area optimization, no power analysis, 
  academic-standard DFT vs ARM's 10-year optimized flow
```

### Quick Check
> 1. What was the primary purpose of the Riscy chip?
> 2. Why couldn't Riscy boot Linux?
> 3. What did the Riscy tape-out teach the SHAKTI team?

---

## 3. SHAKTI C-Class Silicon (2018)

The SHAKTI C-class 64-bit processor — the first "real" SHAKTI chip capable of running Linux — was taped out in 2018.

**C-class silicon specifications:**
- ISA: RV64IMAC (64-bit, integers, multiply, atomic, compressed)
- Pipeline: 5-stage in-order
- I-cache: 16KB 4-way set-associative
- D-cache: 8KB 4-way set-associative
- Hardware TLB (Translation Lookaside Buffer): 32-entry for virtual memory
- Hardware PTW (Page Table Walker): for 3-level RISC-V Sv39 virtual memory
- Interrupt controller: PLIC (Platform-Level Interrupt Controller) compatible
- Process: Intel 22nm FinFET (via IITM-Intel Research Center)
- Die area: approximately 4mm²
- Frequency: ~400–500 MHz

**First Linux boot**: The team successfully booted a stripped-down Linux kernel on SHAKTI C-class in 2019 — a major milestone. The Linux port required:
- RISC-V Linux kernel configured for SHAKTI peripheral map
- Bootloader (BBL — Berkeley Boot Loader, later U-Boot)
- Device drivers for SHAKTI's UART and interrupt controller

**Compliance testing**: C-class passed the RISC-V Foundation's compliance test suite for RV64I, M, A, C extensions. This is the official certification that SHAKTI C-class is a correct RISC-V implementation.

**Application demonstrations:**
- Booted Linux 5.4 kernel
- Ran Dhrystone and CoreMark benchmarks
- Demonstrated UART, SPI, GPIO peripherals working
- Demonstrated PMP (Physical Memory Protection) enforcement

```
C-class performance (measured silicon results):
  
  Benchmark:    SHAKTI C-class    Cortex-A53 (reference)
  Dhrystone:    ~0.6 DMIPS/MHz    ~2.3 DMIPS/MHz
  CoreMark:     ~1.5/MHz          ~3.9/MHz
  Clock:        450 MHz           1500 MHz (28nm)
  
  SHAKTI lags by 3–4× in efficiency
  Primary causes: simple branch predictor, no forwarding optimization,
                  non-optimized standard cell library usage
```

### Quick Check
> 1. What process node was the C-class taped out on?
> 2. What does it mean that C-class passed RISC-V compliance tests?
> 3. What is the Dhrystone benchmark and how does C-class compare to Cortex-A53?

---

## 4. SHAKTI-T Security Processor (2021–2022)

**SHAKTI-T** is the most technically significant SHAKTI chip — adding tagged memory, formally verified security properties, and post-quantum cryptography to the C-class foundation.

**SHAKTI-T specifications:**
- ISA: RV64IMAC + custom security extensions (tagged memory, PQC)
- Pipeline: 5-stage (C-class base) + tag propagation logic
- Tagged memory: 4-bit tag per 64-bit word, propagated through all operations
- PQC accelerators: CRYSTALS-Kyber hardware (NTT engine), CRYSTALS-Dilithium
- TRNG: True Random Number Generator (ring oscillator based)
- Process: Intel 22nm (via IITM-Intel partnership) + TSMC 22nm test version
- Area: ~6mm² (larger than C-class due to tag hardware and PQC accelerators)

**Security verification:**
The team completed formal verification of the tagged memory ISA semantics using Coq:
- Proved tag propagation correctness for all arithmetic/logical instructions
- Proved non-interference: high-security tags do not influence low-security observables
- The proof covers the ISA specification (not the full RTL — a distinction the team acknowledges)

**Hardware tag infrastructure:**
```
SHAKTI-T cache with tags:
  
  Standard L1 D-cache (8KB):       SHAKTI-T L1 D-cache:
  Cache line: [512 bits data]       Cache line: [512 bits data | 32 bits tags]
                                    Each 64-bit word has 4-bit tag
                                    32-bit tag storage per 512-bit data line
                                    5.9% cache overhead for tags
  
  Register file (standard):         Register file (SHAKTI-T):
  32 × 64-bit registers             32 × 64-bit registers + 32 × 4-bit tag registers
  = 256 bytes                       = 260 bytes (1.6% overhead)
```

**PQC accelerator results (measured):**
- Kyber keygen: software = 200M cycles, hardware accelerator = 1.8M cycles (111× speedup)
- Dilithium sign: software = 180M cycles, hardware = 2.5M cycles (72× speedup)
- Power overhead: PQC accelerators add ~15mW at 500 MHz

**Deployment**: SHAKTI-T was evaluated by ISRO and DRDO for classified applications. The formal verification and open RTL audit allowed these agencies to perform independent security verification — something impossible with commercial black-box chips.

### Quick Check
> 1. What does SHAKTI-T add compared to the base C-class?
> 2. What is the speedup of the Kyber hardware accelerator over software?
> 3. Why can ISRO audit SHAKTI-T in a way they cannot audit ARM chips?

---

## 5. Mindgrove SC23 SoC (2023)

**Mindgrove Technologies**, the IIT Madras spinoff, produced the first commercially-oriented SHAKTI-based SoC.

**SC23 specifications:**
- Processor: 4× SHAKTI C-class cores (quad-core)
- L1 caches: 32KB I-cache, 32KB D-cache per core
- L2 cache: 512KB unified shared
- Neural network accelerator: INT8 matrix multiply unit
- Peripherals: UART (4×), SPI (2×), I2C (2×), GPIO (32-bit), PWM, ADC (12-bit), USB OTG
- Memory: LPDDR4 controller (dual-channel, 4GB max)
- Process: TSMC 22nm
- Package: FCBGA (Flip-Chip Ball Grid Array), 529 balls
- Die area: ~25mm²
- Power: ~1–2W typical (quad-core active)
- Operating voltage: 0.9V–1.1V

**Target applications:**
- Industrial IoT gateway (factory automation, smart building)
- Edge inference (on-device ML inference via NN accelerator)
- Medical device controllers
- Defense electronics (non-classified tier)

**Software support:**
- Linux 5.15 LTS (long-term support)
- Mindgrove BSP (Board Support Package)
- OpenCV + ONNX runtime ported for NN accelerator
- FreeRTOS for real-time applications

**Development board**: Mindgrove released the "Secure32" development board with SC23, USB-C power, JTAG debug, onboard flash, 2GB LPDDR4, and expansion headers (Arduino-compatible) — enabling engineers to develop on SHAKTI hardware.

```
SC23 vs commercial alternatives:
  
              SC23              NXP i.MX8M Mini   Raspberry Pi CM4
Processor     4× SHAKTI C      4× Cortex-A53     4× Cortex-A72
Frequency     500 MHz           1.6 GHz           1.5 GHz
Performance   ~1.2 GIPS total   ~8 GIPS           ~12 GIPS
Security      PMP, future PQC   TrustZone          TrustZone
Process       TSMC 22nm         GloFo 28nm         TSMC 28nm
Open source   Yes (SHAKTI IP)   No (ARM IP)        No (ARM IP)
License fee   None              ARM license        ARM license
```

### Quick Check
> 1. What is the Mindgrove SC23's processor configuration?
> 2. What is the neural network accelerator in SC23?
> 3. How does SC23 compare in performance to the Raspberry Pi Compute Module 4?

---

## 6. Test Results and Performance Data

Collected performance data from SHAKTI silicon across the tape-out history:

```
SHAKTI tape-out results summary:

Chip       Process   Freq    Power   Dhrystone    Status
Riscy      Intel 22  100MHz  5mW     ~60 DMIPS    Working (limited)
C-class    Intel 22  450MHz  120mW   ~270 DMIPS   Working, Linux-capable
SHAKTI-T   Intel 22  500MHz  140mW   ~300 DMIPS   Working, security ext.
SC23       TSMC 22   500MHz  1.8W    ~1200 DMIPS  Production (quad-core)
           (quad-core)
```

**Yield analysis (SHAKTI team reports):**
- C-class tape-out: ~30 chips received, ~22 fully functional (73% yield) — acceptable for academic run
- SHAKTI-T: ~20 chips received, ~16 functional (80% yield)
- SC23: Mindgrove commercial quality: target >90% yield (exact data not published)

**Frequency characterization:**
SHAKTI C-class was characterized at 400–500 MHz. Some chips achieved 550 MHz with slight voltage bump (+5% Vdd). No chips ran stably above 600 MHz — consistent with design target.

**Memory performance:**
- L1 hit rate: >95% for typical workloads (cache hierarchy design validated)
- TLB miss rate: ~1% for Linux workloads (normal)
- D-cache miss penalty: 12–15 cycles (to external SRAM) — dominant performance bottleneck

### Quick Check
> 1. What yield did SHAKTI C-class achieve on its first silicon run?
> 2. What is the dominant performance bottleneck identified in SHAKTI C-class?
> 3. At what frequency does SHAKTI C-class plateau?

---

## 7. Lessons Learned

The SHAKTI team has been unusually transparent about lessons from their tape-out journey:

**Lesson 1: BSV generates verbose Verilog**
BSV's compiler generates clean but not minimal Verilog. The synthesized netlist from BSV had 20–30% more area than hand-optimized Verilog would produce. The team is developing a BSV optimization pass.

**Lesson 2: Standard cell library optimization matters**
First tape-outs used Intel's generic library without optimization. Later runs used timing-optimized cell variants, improving critical path by 15–20%. Commercial chips invest heavily in cell library tuning.

**Lesson 3: Post-silicon bring-up is hard without lab infrastructure**
IIT Madras had limited bench equipment for analog/RF debugging. The team invested in oscilloscopes, logic analyzers, and in-circuit emulators after the first tape-out. Bring-up speed improved significantly.

**Lesson 4: PMP adds meaningful overhead**
PMP check logic adds ~2% area and ~0.5% timing overhead. This is acceptable for security applications but the team needed to pipeline the PMP check to avoid degrading frequency.

**Lesson 5: Tagged memory tag arrays need specific SRAM compilers**
Storing 4-bit tags per 64-bit word requires custom SRAM cells (not standard 6T SRAM). The team worked with TSMC to generate a custom tag-embedded memory compiler — this is now documented as a reusable IP block.

**Lesson 6: Commercial tape-out requires more than correct RTL**
Mindgrove's SC23 required: DFT (scan chains, MBIST), ATPG pattern generation, characterization across PVT corners, complete datasheet documentation, package design, board design, and customer-facing support documentation. These post-RTL activities consumed as much effort as the RTL itself.

### Quick Check
> 1. What is the main area overhead from using BSV vs hand-optimized Verilog?
> 2. What infrastructure gap hurt SHAKTI's early post-silicon bring-up?
> 3. What non-RTL activities did Mindgrove need for the SC23 commercial tape-out?

---

## Summary

- **Riscy (2017)**: First SHAKTI silicon, Intel 22nm, 100 MHz, proof of concept, limited functionality.
- **C-class (2018)**: 64-bit, Linux-capable, 450 MHz, Intel 22nm. Passed RISC-V compliance tests.
- **SHAKTI-T (2021–22)**: Tagged memory + PQC accelerators, formally verified security. 500 MHz, 140mW.
- **Mindgrove SC23 (2023)**: Quad-core commercial SoC, TSMC 22nm, 500 MHz, 1.8W. Production-quality design.
- **Performance gap**: C-class/SC23 runs at 3–4× lower efficiency than ARM Cortex-A53 at same process node.
- **Lessons**: BSV overhead, standard cell library optimization, post-silicon infrastructure, commercial DFT requirements, custom tag SRAM needed.

---

## Exercises

### Easy
1. What was the primary purpose of the Riscy chip?
2. What RISC-V extensions does the SHAKTI C-class implement?
3. What makes the SC23 a "production-quality" chip vs an academic research chip?

### Medium
4. Tape-out economics: The SHAKTI C-class tape-out at Intel 22nm used Intel's academic program (estimated at $500K for multi-project wafer slot). Commercial equivalent: TSMC 22nm: ~$15M for a custom tape-out. (a) If the team needed 100 working chips for testing: how many wafers at 73% yield and 100 chips/wafer? (b) Commercial equivalent cost? (c) Why does the academic partnership with Intel enable research that $50M of government funding alone might not? (d) What are the downsides of using Intel 22nm instead of TSMC 22nm for academic work?
5. SHAKTI-T performance analysis: 500 MHz, Kyber hardware: 1.8M cycles per keygen, Dilithium: 2.5M cycles per sign. (a) Time for keygen? (b) Time for sign? (c) A SHAKTI-T-based secure communication system handles 100 TLS connections/second. Each connection needs: 1 keygen + 1 encapsulation. Can SHAKTI-T handle this? (d) For a satellite receiving 10 authenticated commands per second: sign verification (Dilithium verify) takes 2.2M cycles. Can SHAKTI-T verify all commands in real-time?
6. SC23 thermal analysis: Quad-core SC23 at 1.8W typical, operating at 0.9V. Mounted in industrial enclosure (thermal resistance PCB-to-air: 30°C/W). Ambient temperature: 70°C. (a) Processor temperature at full load? (b) Maximum allowable power if T_junction < 125°C? (c) To stay at 105°C (with margin) at 70°C ambient: what is the power budget? (d) If two SC23 chips are used (8-core system): total power? Is passive cooling sufficient?

### Hard
7. Tagged memory SRAM design: Standard 6T SRAM cell stores 1 bit. SHAKTI-T stores 4-bit tag per 64-bit word = 6.25% overhead. (a) A standard 8KB D-cache has 64KB × 8 bits = 512Kbits SRAM. With tags: how much additional SRAM? (b) Tag access must be synchronous with data access (same cycle). Design a memory layout: can you store tags in the same SRAM array or a separate array? Trade-offs? (c) If tags are stored in a separate SRAM with its own port: what does this mean for cache controller complexity? (d) Intel SGX (Secure Guard Extensions) uses separate encrypted memory with integrity tags. Compare SHAKTI-T's approach to Intel SGX: scope of protection, performance overhead, threat model differences.
8. Post-silicon silicon failure analysis: A SHAKTI-T chip fails RISC-V compliance test `rv64ui-p-lb` (load byte test). JTAG shows the CPU boots, basic instructions work, but load-byte returns wrong value with `lb` instruction. (a) What is the `lb` instruction in RISC-V and how should it work? (b) What logic blocks are involved: decoder, ALU, D-cache, sign extension? (c) Design a minimal test sequence (using JTAG to set registers and PC, then single-step) to isolate whether the bug is in: (i) instruction decode, (ii) address calculation, (iii) byte select from cache, (iv) sign extension. (d) This bug passed RTL simulation. How is it possible that a bug exists in silicon but not simulation? (hint: corner cases, PVT variation, initialization state)
