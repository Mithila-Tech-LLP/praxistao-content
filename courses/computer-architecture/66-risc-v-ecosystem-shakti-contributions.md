# Chapter 66: RISC-V Ecosystem and SHAKTI's Contributions

SHAKTI exists within a larger and rapidly growing RISC-V ecosystem. Understanding how SHAKTI contributes to — and benefits from — this ecosystem is important for understanding its long-term viability. The RISC-V ecosystem in 2024 spans from the smallest IoT MCU to the largest HPC supercomputer, from open-source academic designs to billion-dollar commercial products. This chapter maps the ecosystem and SHAKTI's place within it.

## Table of Contents

1. [The RISC-V Ecosystem Overview](#1-the-risc-v-ecosystem-overview)
2. [RISC-V International — The Standard Body](#2-risc-v-international--the-standard-body)
3. [Commercial RISC-V: SiFive, Andes, Codasip](#3-commercial-risc-v-sifive-andes-codasip)
4. [Industrial Deployment: WD, Alibaba, NVIDIA](#4-industrial-deployment-wd-alibaba-nvidia)
5. [SHAKTI's Contributions to RISC-V](#5-shaktis-contributions-to-risc-v)
6. [Software Ecosystem: GCC, LLVM, Linux, Android](#6-software-ecosystem-gcc-llvm-linux-android)
7. [India's RISC-V Strategy](#7-indias-risc-v-strategy)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The RISC-V Ecosystem Overview

The RISC-V ecosystem is unusual: unlike ARM (one company's proprietary design) or x86 (two companies duopoly), RISC-V has thousands of implementors producing radically different products all using the same ISA.

```
RISC-V deployment spectrum (2024):

  Smallest                                                    Largest
  ─────────────────────────────────────────────────────────────────
  IoT MCU     Embedded    Mobile   Desktop/Server    HPC
  
  Examples:
  ESP32-C3    SiFive E21  Android  Alibaba Yitian     Frontier
  (1core,     (CPU IP     (2024    710               (RISC-V
   160MHz)     cores)     preview)  (128cores,3.4GHz) management)
               
  SHAKTI E    SHAKTI C    SHAKTI M  SHAKTI S          Future
```

**RISC-V members (2024)**: 4,000+ organizational members across 70+ countries. Includes:
- Big Tech: Google, Meta, Microsoft, Intel, Qualcomm, NVIDIA, Samsung, IBM
- Academia: MIT, Stanford, Berkeley, IIT Madras
- National labs: Lawrence Berkeley, Fermilab
- Countries with sovereign chip programs: India, China, EU, Russia, Brazil

**Why RISC-V succeeded**:
1. **Free ISA**: Zero cost to implement, zero royalties
2. **Modular**: Implement only the extensions you need
3. **Timing**: Released just as ARM's licensing model became controversial (Arm v8 restriction on old chips, Softbank acquisition uncertainty)
4. **Backed by Berkeley**: Academic credibility + industrial funding ($50M+ DARPA investment)
5. **China's forced hand**: US export restrictions made ARM risky for Chinese companies; RISC-V was the safe alternative

### Quick Check
> 1. How many organizational members does RISC-V International have?
> 2. Why did China's tech industry embrace RISC-V so strongly?
> 3. What makes RISC-V different from ARM in terms of business model?

---

## 2. RISC-V International — The Standard Body

**RISC-V International** is a non-profit membership organization based in Switzerland that owns and develops the RISC-V specification.

**Governance structure:**
- Board of Directors: elected by premium members (Google, Qualcomm, etc.)
- Technical Working Groups (WG): develop specifications for each extension
- IIT Madras/SHAKTI sits on several WGs

**Key specifications (ratified as of 2024):**
```
Specification          Version    Status      Description
RV32I/RV64I            2.1        Ratified    Base integer ISA
M (multiply)           2.0        Ratified    Integer multiply/divide
A (atomic)             2.1        Ratified    Atomic memory operations
F/D/Q (FP)             2.2        Ratified    Single/double/quad FP
C (compressed)         2.0        Ratified    16-bit compressed instructions
V (vector)             1.0        Ratified    Variable-length vectors
Zicsr (CSR)            2.0        Ratified    Control/Status Register ops
Sv39/Sv48/Sv57         —          Ratified    Virtual memory systems
Sm/Ss (supervisor)     1.12       Ratified    Privileged architecture
Zfinx                  1.0        Ratified    FP in integer registers
Hypervisor extension   1.0        Ratified    Virtualization (H extension)
Cryptography           1.0        Ratified    Bit manipulation + crypto
Zicntr (counters)      2.0        Ratified    Performance counters
```

**WIP extensions** (under development):
- P (SIMD/DSP packed): in development, SHAKTI team participates
- J (dynamically translated languages): JVM/WASM optimization
- T (tagged memory): SHAKTI's tagged memory work directly informs this WG

**The frozen ISA guarantee**: Once an extension is ratified, it cannot be changed in a backward-incompatible way. This is crucial for long-lived deployments (automotive, aerospace, medical).

### Quick Check
> 1. Where is RISC-V International based and why does this matter?
> 2. What does "ratified" mean for a RISC-V extension?
> 3. What is the RISC-V T (tagged memory) extension WG and what is SHAKTI's role?

---

## 3. Commercial RISC-V: SiFive, Andes, Codasip

While SHAKTI is open source, commercial RISC-V IP companies provide licensed cores with full commercial support:

**SiFive** (founded 2015 by RISC-V creators, San Jose):
- Business: License RISC-V processor cores (design IP)
- Products: E-series (embedded, MCU class), U-series (Linux-capable, application processor), P-series (DSP-enhanced, automotive)
- Notable: U54, U74 cores used in HiFive development boards; SiFive P650 for premium mobile
- TSMC HPC fabricated; available as hardened IP at multiple nodes
- Revenue: ~$100M (estimate), private company

**Andes Technology** (Taiwan):
- Products: N-class (embedded), A-class (application processor), D-class (DSP)
- Target: automotive, industrial, IoT
- Notable: A45 (Linux-capable, competitive with Cortex-A55)
- Many Chinese SoC companies license Andes cores

**Codasip** (Czech Republic):
- Unique: configurable RISC-V cores + EDA tools to customize ISA extensions
- Used by automotive (NXP, STMicro) for safety-critical applications

**Comparison SHAKTI vs SiFive:**
```
                SHAKTI (IIT Madras)    SiFive
Business model  Open source            Commercial license
ISA             RISC-V + security ext  RISC-V + SiFive extensions
Security        Tagged memory, PQC     Limited (standard PMP)
Performance     C-class ~500MHz        U74 ~1.5GHz
Support         Academic               Commercial SLA
Price           Free (IP)              $500K–$2M license
Target          Strategic/government   Commercial SoCs
```

### Quick Check
> 1. What is SiFive's business model and how does it differ from SHAKTI?
> 2. What makes Codasip unique among RISC-V IP vendors?
> 3. At what performance level does SiFive's U74 compete?

---

## 4. Industrial Deployment: WD, Alibaba, NVIDIA

The most significant validation of RISC-V's commercial viability comes from massive industrial deployments:

**Western Digital (WD)**: Announced in 2018 that they would transition all their storage processors (flash controllers, HDD firmware processors) to RISC-V. Target: 1 billion RISC-V cores per year. WD open-sourced their SweRV EH1 core design.
- Why WD? Storage controllers need billions of ultra-efficient compute cores; RISC-V eliminates ARM license fees across their entire product line

**Alibaba T-Head (DAMO Academy)**:
- XuanTie C910 (2019): 16-core RISC-V cluster for edge computing
- XuanTie C906: RISC-V core used in Allwinner D1 SoC (Linux-capable, <$5 SoC)
- **Yitian 710 (2021)**: 128-core server chip at TSMC N5, 3.4 GHz, for Alibaba Cloud
  - This is the most powerful commercial RISC-V chip ever built
  - 60 billion transistors, 512-bit vector extensions
  - Powers Alibaba Cloud compute instances

**NVIDIA**:
- Internal GPU management processors (PMC cores) are RISC-V
- Falcon security processor (in GeForce, Tesla GPUs) replaced with RISC-V
- NVIDIA's Hopper H100 GPU contains multiple RISC-V cores for management tasks

**Google**:
- Titan M2 security chip in Pixel phones contains RISC-V core
- OpenTitan project: open-source RISC-V root of trust chip design (Google + partners)

**Qualcomm**:
- RISC-V for sensor hub/microcontroller tier
- Hexagon DSP (future versions) may use RISC-V

```
RISC-V deployment scale (2024):
  Total RISC-V chips shipped: >10 billion (mostly ESP32, WD controllers)
  RISC-V cores: >100 billion total across all applications
  Market: $3B RISC-V IP market, growing 40% annually
```

### Quick Check
> 1. What is Western Digital's RISC-V commitment and why did they make it?
> 2. What is the Alibaba Yitian 710 and why is it historically significant?
> 3. What RISC-V role does NVIDIA use in its GPUs?

---

## 5. SHAKTI's Contributions to RISC-V

SHAKTI is not only a consumer of RISC-V — it has contributed back to the ecosystem:

**Tagged memory (T extension)**: SHAKTI's tagged memory work is the most detailed published implementation of fine-grained hardware tagging for RISC-V. The RISC-V T extension WG references SHAKTI-T as a leading implementation. SHAKTI's formal verification work (Coq proofs) has influenced the extension's security specification.

**Security Working Group participation**: IIT Madras faculty and students participate in the RISC-V security architecture committee, contributing to:
- Crypto extensions specification
- Privilege architecture security considerations
- Physical Memory Protection extension improvements

**RISC-V India chapter**: IIT Madras organized the RISC-V India Summit, bringing together Indian industry and academia around RISC-V. This built community and influenced government policy (India Semiconductor Mission now explicitly includes RISC-V).

**Open-source RISC-V test infrastructure**: SHAKTI released compliance test extensions and simulation infrastructure that the broader RISC-V community uses:
- BSV models of RISC-V behavior for formal verification
- FPGA emulation environments for SHAKTI/RISC-V development

**Academic publications**: 50+ papers published by IIT Madras on SHAKTI and RISC-V topics — microarchitecture, security, verification, VLSI design — making SHAKTI an internationally recognized research program.

**Bluespec BSV advocacy**: SHAKTI demonstrated that BSV (Bluespec SystemVerilog) is viable for production-quality RISC-V implementation, influencing several other academic RISC-V projects to adopt BSV.

### Quick Check
> 1. What RISC-V extension did SHAKTI's work directly influence?
> 2. What is the RISC-V India Summit and what does it accomplish?
> 3. How does SHAKTI's use of BSV benefit the broader RISC-V community?

---

## 6. Software Ecosystem: GCC, LLVM, Linux, Android

A processor is useless without software. RISC-V has mature software ecosystem support:

**Compilers:**
- GCC: full RISC-V backend, all standard extensions, `riscv64-linux-gnu-gcc`
- LLVM/Clang: full RISC-V support in upstream LLVM
- SHAKTI SDK: customized GCC/Clang for SHAKTI-specific extensions

**Operating systems:**
- Linux: upstream RISC-V support since kernel 4.15 (2018). Full SMP, NUMA, all major subsystems work
- FreeRTOS: supported for RISC-V M-class (embedded)
- Zephyr RTOS: full RISC-V support
- RT-Thread: Chinese RTOS, full RISC-V support

**Android on RISC-V:**
- Google added experimental RISC-V support to AOSP (Android Open Source Project) in 2022
- Alibaba T-Head demonstrated Android running on RISC-V at 2022 RISC-V Summit
- Full Android on RISC-V (production quality) expected ~2025–2026

**SHAKTI software:**
- Linux 5.10+ boots on SHAKTI C-class
- SHAKTI SDK: peripheral drivers, BSP, startup code
- GNU toolchain with SHAKTI extensions
- RISC-V QEMU: simulate SHAKTI processor in software

**Missing ecosystem (honest gaps):**
- No major commercial OS vendor (Windows on RISC-V: not announced)
- No RISC-V iPhone/Samsung equivalent (yet)
- Sparse commercial software (games, productivity apps) for RISC-V
- JVM/V8/LLM inference engines: RISC-V support present but less optimized than ARM/x86

### Quick Check
> 1. Which Linux kernel version first added upstream RISC-V support?
> 2. What is AOSP and what did Google add to it for RISC-V?
> 3. What are the two biggest software ecosystem gaps for RISC-V vs ARM?

---

## 7. India's RISC-V Strategy

India's government and industry have aligned around RISC-V as the foundation for domestic semiconductor capability:

**India Semiconductor Mission (ISM):**
- $10B PLI (Production Linked Incentive) scheme for fabs and design
- Explicitly includes RISC-V-based processor design as a priority
- SHAKTI listed as technology to be supported

**CDAC (Centre for Development of Advanced Computing):**
- Government agency developing RISC-V for strategic applications
- Vega processor: CDAC's RISC-V microprocessor (C906-based, 500 MHz)
- Working with ISRO for satellite processor qualification

**IIT ecosystem:**
- IIT Bombay: IITB-RISC (students design and tape out RISC-V chips as part of coursework)
- IIT Delhi: RISC-V research, memory systems
- IIT Kharagpur: RISC-V compiler and architecture research

**Industry:**
- Tata Electronics: building semiconductor fab in Dholera (Gujarat), includes design
- HCL Semiconductors: chip design services, exploring RISC-V
- InCore Semiconductors (spinoff from IIT Madras): commercial RISC-V cores, particularly for security applications

**The vision**: India in 2030 with a complete domestic RISC-V ecosystem:
- IIT Madras SHAKTI cores → licensed to Indian SoC companies
- Indian SoC companies → tapeout at Indian fab (Tata/ISMC)
- Indian-designed, Indian-fabricated chips → deployed in government and consumer electronics
- Reduction of $10–15B in annual semiconductor imports

### Quick Check
> 1. What is the India Semiconductor Mission and what is its total investment?
> 2. What is CDAC's role in the Indian RISC-V ecosystem?
> 3. What is the 2030 vision for India's RISC-V semiconductor ecosystem?

---

## Summary

- **RISC-V ecosystem**: 4,000+ members, 10B+ chips shipped, spans IoT to HPC. Free, open, modular ISA.
- **Commercial RISC-V**: SiFive (silicon IP), Andes (Taiwan/automotive), Codasip (configurable) — all licensing RISC-V cores.
- **Industrial deployment**: Western Digital (1B cores/year), Alibaba Yitian 710 (128-core server), NVIDIA (GPU management), Google (Titan M2 security).
- **SHAKTI contributions**: Tagged memory (T extension WG), security WG participation, India RISC-V community, BSV advocacy, open test infrastructure.
- **Software ecosystem**: GCC/LLVM complete, Linux since 4.15, Android RISC-V coming, FreeRTOS/Zephyr. Gaps: Windows, commercial apps.
- **India strategy**: ISM ($10B), CDAC (Vega processor), IIT ecosystem, Tata fab, InCore Semiconductors.

---

## Exercises

### Easy
1. What is RISC-V International and where is it based?
2. Why did Western Digital commit to 1 billion RISC-V cores per year?
3. What is the Alibaba Yitian 710 and why is it significant?

### Medium
4. Ecosystem comparison: Compare RISC-V, ARM, and x86 on the following dimensions for a startup wanting to build a smartwatch chip: (a) ISA license cost, (b) compiler/OS support quality, (c) design tool availability, (d) available off-the-shelf processor IP, (e) long-term supply chain risk. For each dimension: rank the three options and justify. (f) Overall recommendation for the startup.
5. RISC-V extension selection: You are designing an IoT sensor SoC for a smart agriculture system (soil moisture, temperature, GPS, cellular). Select the RISC-V extensions needed: (a) RV32I vs RV64I? Why? (b) M extension: needed? (c) A extension: needed? Why? (d) F/D extensions: needed? (e) C extension: needed? Why? (f) V extension: needed? (g) What custom extensions might you add? (h) Final ISA string (e.g., RV32IMC)?
6. India semiconductor supply chain: India's chip import bill: $25B/year. SHAKTI helps reduce this by substituting domestic designs. (a) If India deploys SHAKTI in 50% of government/defense applications (market ~$1B/year): what import substitution is achieved? (b) If Mindgrove ships 10M IoT chips/year at $3 average: revenue and import substitution? (c) What is India's most realistic near-term path to $5B/year chip export: design+export vs fab+export? (d) The China example: China went from $10B domestic production (2010) to ~$70B (2023) — but still imports $400B/year. What does this suggest about India's timelines?

### Hard
7. RISC-V vs ARM for India's national strategy: India's strategic applications (military, space, government) require processors that are: secure (no backdoors), reliable supply, auditable design. ARM offers: mature ecosystem, high performance, but requires license from UK company. RISC-V offers: open ISA, but implementations vary in quality and ecosystem maturity. (a) What are the specific strategic risks of ARM for each: a missile guidance system, a satellite attitude controller, a government secure communications device? (b) What would India need to build to have "RISC-V sovereignty" (not just ISA freedom but complete stack)? List all components. (c) Timeline: if India started this effort in 2024 with full government commitment ($5B/year), when could it achieve complete strategic processor independence? (d) Is complete independence the right goal, or is "managed interdependence" (keep ARM for commercial, use RISC-V for strategic) more practical?
8. SHAKTI's academic vs production tension: SHAKTI is simultaneously: (a) an academic research project (designed to teach chip design, publish papers), and (b) a production chip aimed at government deployment. These goals conflict. Analyze: (a) What design choices are made for academic reasons that harm production quality? (hint: BSV language, verification coverage gaps) (b) What production requirements (documentation, support, long-term availability) are hard for an academic project to meet? (c) How should IIT Madras and Mindgrove Technologies divide responsibilities to balance both goals? (d) Compare to how Berkeley (originally developed RISC-V) + SiFive (commercialized it) divided academic/commercial roles — is this the right model for India?
