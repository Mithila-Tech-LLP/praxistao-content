# Chapter 63: The SHAKTI Architecture — Design Philosophy

SHAKTI is not a single processor but a family of processors designed with a common philosophy. Each member targets a different market segment — from low-power microcontrollers to high-performance server cores — all sharing the RISC-V ISA foundation and the SHAKTI design principles: open source, security-first, and India-relevant. This chapter explains the architectural decisions that define all SHAKTI processors, the microarchitectural choices, and how they compare to commercial alternatives.

## Table of Contents

1. [SHAKTI Design Philosophy](#1-shakti-design-philosophy)
2. [The RISC-V ISA in SHAKTI](#2-the-risc-v-isa-in-shakti)
3. [SHAKTI's Microarchitecture](#3-shaktis-microarchitecture)
4. [Security as a First-Class Design Goal](#4-security-as-a-first-class-design-goal)
5. [SHAKTI vs ARM Cortex Comparison](#5-shakti-vs-arm-cortex-comparison)
6. [The Development Methodology](#6-the-development-methodology)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. SHAKTI Design Philosophy

The SHAKTI design team articulated specific principles that guide every architectural decision:

**Principle 1: Open and auditable**
Every RTL file, every design document, every test bench is publicly available. The goal is not just to build a chip but to build a chip where anyone in the world can audit the design — crucial for security-sensitive applications that cannot trust black-box IP.

**Principle 2: Security by design, not afterthought**
Traditional chip design adds security features at the end ("tapeout minus 2 months, add the security core"). SHAKTI designed security extensions at the architectural level — tagged memory, physical memory protection, and post-quantum cryptographic accelerators are first-class citizens, not add-ons.

**Principle 3: India-relevant applications**
SHAKTI prioritizes applications important to India's strategic needs: satellite systems (ISRO), defense electronics (DRDO), government secure communications, and IoT for rural infrastructure. This means specific features: radiation tolerance (for space), tamper resistance (for defense), lightweight cryptography (for IoT).

**Principle 4: Academic rigor, production quality**
SHAKTI is an academic project that aims for industry quality — formal verification of security properties, complete verification testbenches, timing closure for real silicon. The project explicitly rejects the "good enough for research" standard.

**Principle 5: Ecosystem building**
SHAKTI publishes software stacks (Linux port, GCC SHAKTI patches, FreeRTOS port, SHAKTI SDK) alongside the hardware. A chip nobody can program is useless — the software stack is part of the product.

### Quick Check
> 1. What does "open and auditable" mean for a chip design?
> 2. Why is security as a "first-class design goal" different from adding security later?
> 3. Name two India-specific application domains that shaped SHAKTI's design.

---

## 2. The RISC-V ISA in SHAKTI

All SHAKTI processors implement RISC-V. The specific ISA extensions depend on the processor class:

**Base ISA**: SHAKTI processors implement the standard RISC-V base ISA:
- **RV32I**: 32-bit base integer ISA (E-class, C-class smaller)
- **RV64I**: 64-bit base integer ISA (C-class, S-class)

**Standard extensions used:**
```
Extension  Name            Description
I          Integer         Base integer (ALU, load/store, branches)
M          Multiply        Integer multiply and divide
A          Atomic          Atomic memory operations (for OS, SMP)
F          Single FP       32-bit floating point (IEEE 754)
D          Double FP       64-bit floating point
C          Compressed      16-bit compressed instructions (code size -30%)
Zicsr      CSR             Control and Status Registers (for OS/interrupt handling)
Zifencei   Instr fence     Instruction fetch barrier (for self-modifying code/JIT)
```

**SHAKTI-specific non-standard extensions:**
RISC-V's custom extension space (opcodes 0000 011 and 1111 011 are reserved for custom use). SHAKTI uses these for:
- **Tagged memory extension**: Tag bits on every 64-bit word for memory safety
- **Post-quantum cryptography accelerators**: Hardware acceleration for CRYSTALS-Kyber (key encapsulation), CRYSTALS-Dilithium (signatures)
- **Physical memory attribution**: Track ownership of memory regions for isolation

**RISC-V privilege levels in SHAKTI:**
```
Machine Mode (M-mode):    Highest privilege. Direct hardware access.
                          SHAKTI secure monitor runs here.
                          
Supervisor Mode (S-mode): OS kernel runs here (Linux, FreeRTOS supervisor)
                          
User Mode (U-mode):       Applications run here.
                          Cannot access privileged CSRs.
```

**Physical Memory Protection (PMP)**: SHAKTI implements RISC-V PMP — a hardware mechanism to restrict which physical addresses each privilege level can access. Even if the OS kernel is compromised, M-mode (secure monitor) retains control of PMP registers.

### Quick Check
> 1. What does the RISC-V 'C' extension provide and why is it useful for embedded systems?
> 2. What is Physical Memory Protection (PMP) in RISC-V?
> 3. What are the three privilege levels in RISC-V and what runs at each level?

---

## 3. SHAKTI's Microarchitecture

SHAKTI processors span a range from simple in-order to more complex pipelines:

**E-class (Embedded):**
- Simplest class: 3-stage in-order pipeline (Fetch/Decode+Execute/Memory+Writeback)
- No branch prediction (single-issue in-order)
- No caches (directly interfaces to SRAM or flash via tightly-coupled memory)
- RV32IMAC ISA
- Target: microcontrollers, IoT sensors, ultra-low-power applications
- Area: ~10K–50K gates
- Power: <1mW

**C-class (Commercial):**
- 5-stage in-order pipeline
- 16KB I-cache, 8KB D-cache
- RV64IMAC ISA (64-bit)
- Branch predictor (simple two-level)
- Hardware multiply/divide
- Area: ~100K–200K gates
- Target: embedded Linux applications, IoT gateways, ISRO satellite processors

**I-class (Industrial):**
- 7-stage in-order pipeline with some instruction-level parallelism
- Larger caches: 32KB I-cache, 32KB D-cache, optional L2
- FPU (F and D extensions)
- Target: industrial controllers, embedded networking

**M-class (Mobile):**
- Out-of-order execution (limited ROB)
- Hardware prefetcher
- NEON-like SIMD extensions
- Target: mobile/tablet SoCs

**S-class (Server):**
- Superscalar, out-of-order
- Large cache hierarchy (L1/L2/L3)
- Hardware page table walker
- Server-grade performance target
- Still in development as of 2023

```
SHAKTI processor family:

  Performance
       ↑
  S-class  ──────────────────────────── (server, OOO, >1GHz)
  M-class  ────────────────── (mobile, light OOO)
  I-class  ──────────── (industrial, deeper pipeline)
  C-class  ──────── (commercial, 64-bit, Linux capable)
  E-class  ─── (embedded, 3-stage, MCU)
       │
       └──────────────────────────────→ Power (mW)
         0.1mW         10mW        100mW
```

### Quick Check
> 1. What pipeline does the SHAKTI E-class use and what is it targeted for?
> 2. What makes the C-class the first "Linux-capable" SHAKTI processor?
> 3. What is the difference between in-order (C-class) and out-of-order (M-class/S-class) execution in SHAKTI?

---

## 4. Security as a First-Class Design Goal

SHAKTI's security architecture goes beyond what most commercial processors offer. This is driven by the program's focus on strategic applications where hardware-level security assurance is required.

**Tagged memory (Memory Tagging Architecture, MTA):**
Every 64-bit word in memory has an associated 4-bit tag. The processor propagates tags through computation:
- Tag 0: uninitialized/default
- Tag 1: trusted OS kernel data
- Tag 2: user application data
- Tag 3–15: application-defined

```
Tagged memory in SHAKTI:
  
  Memory word: [63:0] data | [3:0] tag
  
  Benefits:
  1. Buffer overflow detection: writing beyond allocated region
     changes tag → hardware trap
  2. Taint tracking: follow untrusted (tainted) data through computation
     → detect when tainted data is used as control flow target
  3. Type safety: different types tagged differently → type confusion attacks detected
  
  Example: C buffer overflow attempt:
    char buf[128];  // tag = USER_DATA
    int ret_addr;   // tag = RETURN_ADDR
    
    strcpy(buf, evil_input);  // if overflow overwrites ret_addr:
    return;                   // tag mismatch → hardware trap before exploit!
```

**SHAKTI-T (Tagged memory variant)**: A specific SHAKTI implementation with full tagged memory support, formally verified security properties, and designed for government/defense applications requiring information flow control.

**Secure enclave (SHAKTI-S)**: Isolation of a trusted execution environment within the processor, similar in concept to ARM TrustZone but implemented differently:
- M-mode (machine mode) = secure world
- S-mode (supervisor) = OS
- Physical memory ranges allocated to secure world are inaccessible from OS
- Secure monitor (SBI — Supervisor Binary Interface) manages the boundary

**Post-quantum cryptography acceleration:**
CRYSTALS-Kyber (KEM) and CRYSTALS-Dilithium (signatures) are NIST PQC standard algorithms. Hardware acceleration in SHAKTI-T provides:
- 10–100× speedup vs software implementation
- Constant-time execution (prevents timing side-channel attacks)
- These are integrated directly into the RISC-V instruction stream via custom instructions

**Timing side-channel mitigations:**
- Data cache partitioned for Spectre/Meltdown mitigation
- SHAKTI implements variants of RIDL/Fallout mitigations for in-order architectures
- FLUSH+RELOAD attack surface reduced by limiting shared cache state

### Quick Check
> 1. What is tagged memory and what security attacks does it prevent?
> 2. How does SHAKTI implement a secure enclave?
> 3. Why is post-quantum cryptography important for strategic applications?

---

## 5. SHAKTI vs ARM Cortex Comparison

How does SHAKTI compare to the ARM processors it hopes to supplement in India?

```
Comparison: SHAKTI C-class vs ARM Cortex-A53

                SHAKTI C-class    ARM Cortex-A53
ISA             RISC-V RV64IMAC   ARMv8-A
Pipeline        5-stage in-order  8-stage in-order
Frequency       ~500 MHz (22nm)   ~1.5 GHz (28nm)
IPC             ~0.5–0.8          ~1.2–1.5
Performance     ~0.5 GIPS         ~2 GIPS (28nm)
License         Open source       ARM license required
Security        PMP, tagged mem,  TrustZone, pointer auth
                PQC accel
Ecosystem       Growing           Mature (Android, Linux)
Maturity        Academic + prod   Production, mature
NRE for SoC     None (open IP)    $500K–$2M license

SHAKTI disadvantage: 3-4× lower performance per MHz vs Cortex-A53
SHAKTI advantage:    Free, auditable, security extensions, no license risk
```

**The performance gap**: SHAKTI C-class achieves ~500 MHz on TSMC 22nm. ARM Cortex-A53 achieves ~1.5 GHz on the same or worse process. The gap is real:
- SHAKTI team acknowledges this and views C-class as a learning step toward S-class
- For strategic applications (satellites, secure communications), security > raw speed
- For commercial markets (phones, laptops), the gap is too large to compete near-term

**Where SHAKTI competes today:**
- Government/defense applications where ARM license risk is unacceptable
- Space applications (ISRO) where radiation tolerance + known IP is required
- Education and research infrastructure
- IoT applications where cost (no license fee) > performance

### Quick Check
> 1. What is the performance gap between SHAKTI C-class and ARM Cortex-A53?
> 2. In what applications does SHAKTI's openness outweigh its performance disadvantage?
> 3. Why can't SHAKTI compete with ARM for consumer smartphones today?

---

## 6. The Development Methodology

SHAKTI is developed as a research project with production aspirations. The methodology:

**Language choice**: SHAKTI RTL is written in **BSV (Bluespec SystemVerilog)** — a high-level hardware description language from MIT/Bluespec Inc. BSV provides:
- Strong type system (reduces bugs vs Verilog)
- Guarded atomic actions (easier reasoning about concurrency)
- Modular, composable designs

BSV compiles to synthesizable Verilog, which then flows through standard EDA tools.

**Why BSV, not Verilog?**
Academic researchers prefer BSV because it is easier to reason about formally and has better abstraction. Industrial engineers prefer Verilog/SystemVerilog because the toolchain and expertise are mature. SHAKTI's BSV choice is a deliberate academic decision.

**Verification**: SHAKTI uses:
- Directed testbenches (unit tests for ALU, cache, branch predictor)
- Constrained random simulation with RISC-V compliance tests
- RISC-V formal verification tests (riscv-formal, SAIL model comparison)
- SAIL: formal mathematical model of RISC-V semantics; SHAKTI is verified against SAIL

**RISC-V compliance testing**: RISC-V International provides a compliance test suite. Any processor claiming RISC-V compatibility must pass this suite. SHAKTI processors have passed the compliance tests for their respective ISA profiles.

**Open source release**: All SHAKTI IP is released on GitHub (github.com/iitm-shakti):
- RTL source code (BSV)
- Synthesis scripts
- Testbenches
- Documentation
- Linux port and SDK

### Quick Check
> 1. What is BSV (Bluespec SystemVerilog) and why did SHAKTI choose it?
> 2. What is the RISC-V compliance test suite?
> 3. Where is SHAKTI's source code published?

---

## Summary

- **SHAKTI philosophy**: open/auditable, security-first, India-relevant, academic rigor, ecosystem building.
- **ISA**: RISC-V (RV32I or RV64I base + M, A, C, F, D extensions). Custom extensions for tagged memory and PQC.
- **Processor family**: E-class (3-stage, MCU, <1mW) → C-class (5-stage, 64-bit, Linux) → I-class → M-class → S-class (server, OOO, in development).
- **Security**: Tagged memory (detects buffer overflows, taint tracking), PMP (memory isolation), secure enclave (M-mode), PQC hardware accelerators.
- **vs ARM Cortex-A53**: SHAKTI is 3–4× slower per MHz but free, auditable, and has unique security extensions.
- **Development**: BSV language, RISC-V compliance tests, SAIL formal verification, GitHub open source.

---

## Exercises

### Easy
1. What is the RISC-V 'A' (atomic) extension and why is it needed for multi-core systems?
2. What is tagged memory and how does it prevent buffer overflow exploits?
3. Why does SHAKTI use BSV instead of Verilog?

### Medium
4. SHAKTI performance analysis: C-class at 500 MHz with IPC ≈ 0.6. (a) Calculate MIPS (million instructions per second). (b) An application requires 100 million instructions to complete. Execution time on SHAKTI C-class? (c) ARM Cortex-A53 at 1.5 GHz with IPC ≈ 1.3: execution time for same 100M instructions? (d) What is the performance ratio? (e) For an ISRO satellite application requiring <10W total power and 50 MIPS minimum: can SHAKTI C-class satisfy this? (estimate power from area and frequency)
5. Tagged memory protection: A C program has a stack buffer overflow vulnerability: `char buf[64]` followed by `int ret_addr`. (a) In a standard processor, show how writing 80 bytes to buf overwrites ret_addr. (b) In SHAKTI's tagged memory: tag for buf[0..63] = USER_STACK, tag for ret_addr = RETURN_TAG. When the write to buf+80 happens, what does the hardware do? (c) This prevents the exploit but what about legitimate programs that cast pointers across type boundaries (common in C systems code)? How does SHAKTI's tagged memory handle intentional type-unsafe operations? (d) Compare to ARM's Pointer Authentication Codes (PAC) — same goal, different mechanism. Which is more flexible?
6. RISC-V privilege levels: A Linux application, the Linux kernel, and the SHAKTI secure monitor all run simultaneously. Map each to a privilege level and explain: (a) What happens when the application calls `read(fd, buf, 1024)` (system call)? (b) What happens when the kernel needs to handle a page fault? (c) What happens when the secure monitor needs to control the PMP to protect secure memory from the kernel? (d) Why is the three-level hierarchy necessary rather than just two (trusted/untrusted)?

### Hard
7. SHAKTI security formal verification: The SHAKTI team claims tagged memory is formally verified. (a) What does "formally verified" mean for a hardware security property? (b) The property to verify: "A write to an address with tag X can never change the tag of an adjacent address." Express this as a temporal logic formula (hint: use ∀ time, ∀ address notation). (c) The proof must hold for all possible instruction sequences, including speculative execution. Why does speculation create a challenge for formal verification? (d) Compare to formal verification of software (Coq proofs, CompCert verified C compiler): what additional challenges does hardware verification have?
8. SHAKTI S-class design challenge: The SHAKTI team wants to design an S-class (server-grade) processor targeting 2 GHz at TSMC 28nm, with out-of-order execution, 64-entry ROB, 3-wide superscalar. (a) What pipeline stages would such a processor need? (b) RISC-V's tagged memory extension adds 4 tag bits per 64-bit word. For a 2MB L2 cache: how much extra SRAM is needed for tags? (c) Out-of-order execution with tagged memory: when an instruction reads a tagged value, executes, and produces a result — how should tags be propagated through the ROB? (d) If SHAKTI S-class is 50% of ARM Cortex-A72 performance (a reasonable target at same process): what applications in India's strategic domain would justify building 100,000 units annually?
