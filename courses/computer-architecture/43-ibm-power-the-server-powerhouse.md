# Chapter 43: IBM POWER — The Server Powerhouse

When most people think of server CPUs, they think Intel Xeon or AMD EPYC. But for decades, IBM's POWER architecture has powered the world's most demanding workloads — global financial trading systems, airline reservation databases, mainframe-class reliability in a UNIX server. The POWER ISA is the only architecture besides x86-64 and ARM64 with serious commercial deployment in enterprise servers. Understanding POWER teaches us what "server-class" really means: extreme reliability, extreme I/O bandwidth, extreme thread counts, and the willingness to pay $50,000+ per socket for it.

## Table of Contents

1. [POWER's Origins: RISC before RISC](#1-powers-origins-risc-before-risc)
2. [POWER Architecture Fundamentals](#2-power-architecture-fundamentals)
3. [POWER10 — The Modern POWER Chip](#3-power10--the-modern-power-chip)
4. [SMT8 — Eight Threads Per Core](#4-smt8--eight-threads-per-core)
5. [OpenPOWER — Open Sourcing POWER](#5-openpower--open-sourcing-power)
6. [Who Uses POWER and Why](#6-who-uses-power-and-why)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. POWER's Origins: RISC before RISC

IBM's POWER (Performance Optimization With Enhanced RISC) architecture predates the RISC revolution as a mainstream concept. The **IBM 801** (1975) was one of the earliest RISC-like designs — developed by John Cocke at IBM Research, it used a load-store architecture with hardwired logic and no microcode.

The 801's ideas were productized as the ROMP (Research OPD MicroProcessor, 1981) and then as the **RIOS (RISC Instruction Set)** in the RT PC (1986). The first commercial POWER chip was in the **IBM RS/6000** workstation (1990) — which ran AIX (IBM's Unix).

**PowerPC (1991)**: IBM, Apple, and Motorola formed the AIM alliance to create PowerPC — a version of POWER for personal computers. Apple used PowerPC in Macs from 1994 until the Intel transition in 2006. Motorola's PowerPC descendants live on in automotive and aerospace applications today.

**POWER vs PowerPC**: IBM kept the high-end POWER line for workstations/servers; PowerPC was for PCs and embedded. They share ancestry but diverged in features.

### Quick Check
> 1. What year was the first commercial POWER chip, and what was it used for?
> 2. What was the AIM alliance and what did it produce?
> 3. Why does Apple no longer use PowerPC?

---

## 2. POWER Architecture Fundamentals

POWER64 (64-bit POWER ISA, the modern form) is a RISC, load-store architecture with 64-bit registers. It has features beyond x86 and ARM:

**Register file**: 32 × 64-bit general-purpose registers (GPRs), 32 × 64-bit FP registers (FPRs), 32 × 128-bit VMX/Altivec SIMD registers, 32 × 512-bit VSX (Vector Scalar eXtension) registers.

**Condition Register (CR)**: 8 × 4-bit condition field groups — multiple outstanding compare results simultaneously.

**Branch prediction**: POWER10 has one of the most sophisticated branch predictors in any production CPU — over 20 cycles of branch history, large BTB.

**Memory model**: POWER has a famously **weak memory consistency model** (POWER uses a relaxed consistency model where stores and loads can be reordered significantly). Programs must use explicit memory barriers (`sync`, `isync`, `lwsync`) for synchronization. This provides maximum hardware optimization freedom at the cost of more complex software.

```
POWER memory ordering:
  weak = hardware can reorder almost everything
  programmer (or compiler) must insert:
    sync  = full barrier (like x86 MFENCE)
    lwsync = lightweight sync (no store-load reordering but lighter than sync)
    isync = instruction sync
  
  This is unusual: even ARM is stricter than POWER in memory ordering
```

**Big-endian vs Little-endian**: POWER supports both, but historically has been big-endian. POWER8+ can operate in either mode. AIX is big-endian; Linux on POWER is typically little-endian (matching x86 software expectations).

### Quick Check
> 1. How many vector registers does POWER have, and how wide are they?
> 2. What is the POWER memory model and why is it "weak"?
> 3. Why does the weak memory model require software memory barriers?

---

## 3. POWER10 — The Modern POWER Chip

**POWER10 (2021)**: IBM's latest generation POWER chip.

```
POWER10 specifications:
  Cores: up to 15 per chip (from 12-core SC9 to 15-core DC chip variants)
  Threads per core: up to 8-way SMT
  Max threads per chip: 15 × 8 = 120 hardware threads
  
  Process: Samsung 7nm
  Die size: ~602 mm²
  Transistors: ~18 billion
  
  L1 cache: 48KB data / 48KB instruction per core
  L2 cache: 2MB per core
  L3 cache: 8MB per core → up to 120MB total L3
  
  Memory: up to 4TB DDR4/DDR5 per socket (via 8 memory channels)
  Memory bandwidth: 409 GB/s
  PCIe: PCIe 5.0
  
  Power: 250-300W typical (server-class TDP)
  
  Price: IBM Power S1014/S1022/S1024 servers: $10,000–$100,000+
```

**OpenCAPI / OMI**: IBM's memory interface. OMI (Open Memory Interface) allows DRAM to be attached to the POWER chip via a PCIe Gen 4-like serial link, supporting very large memory configurations.

**Memory Inception (POWER10 feature)**: Allows one POWER system to use the memory of another POWER system over the network (Ethernet/InfiniBand) as part of its address space — essentially distributed memory coherence across physical servers. Useful for in-memory databases and HPC.

### Quick Check
> 1. What is the maximum thread count per POWER10 chip?
> 2. What is "Memory Inception" in POWER10 and what workloads benefit from it?
> 3. Why are POWER servers so expensive compared to x86 servers?

---

## 4. SMT8 — Eight Threads Per Core

The most distinctive POWER feature: **8-way SMT** — 8 hardware threads per physical core.

In Chapter 27, we discussed Intel's 2-way SMT (Hyper-Threading). POWER goes much further:

```
POWER10 core with SMT8:
  Physical core: wide OOO machine, deep pipeline
  8 logical threads running simultaneously:
    Thread 0, 1, 2, 3, 4, 5, 6, 7 — each with own:
      - PC (program counter)
      - Architectural registers (32 GPRs × 8 copies)
      - Status/control registers
    
    Shared:
      - Execution units
      - ROB (partitioned between threads)
      - L1/L2 cache
      - Branch predictor
  
  OS sees: 15 physical cores × 8 threads = 120 logical CPUs per socket
```

**SMT4 mode**: When fewer threads are needed, POWER can run in SMT4 (4 threads/core) mode, giving each thread more hardware resources. Single-thread performance improves in SMT1 mode (all hardware resources to one thread).

**Why 8-way SMT?**: IBM's server workloads include many concurrent clients (hundreds of database sessions, thousands of web requests). SMT8 allows the hardware to context-switch far less — keeping 8 contexts alive simultaneously, each making progress even if some stall on memory. For throughput-oriented workloads (not single-thread latency), SMT8 provides near-linear scaling.

**The tradeoff**: Single-thread performance per logical CPU is lower than Intel or AMD (because resources are shared 8 ways). POWER10 is not competitive with Intel/AMD for single-threaded workloads at equivalent price. It is competitive for highly multithreaded database/transaction workloads.

### Quick Check
> 1. How many hardware threads does each POWER10 core support?
> 2. What workloads benefit most from 8-way SMT?
> 3. When would you prefer SMT1 mode on POWER vs SMT8?

---

## 5. OpenPOWER — Open Sourcing POWER

In 2013, IBM made a stunning move: open-sourced the POWER ISA through the **OpenPOWER Foundation** (now part of the Linux Foundation).

Companies can now:
- Implement POWER-compatible chips without paying IBM
- Modify the ISA and create extensions
- Access IBM's full architectural documentation

**Why IBM opened POWER**:
1. POWER was losing ground to x86 and ARM. Opening the ISA could attract more implementations.
2. The OpenPOWER ecosystem could create complementary chips (network controllers, accelerators) that enhance IBM Power systems.
3. Potential revenue from consulting, systems, and software rather than just silicon.

**OpenPOWER contributions:**
- **Raptorcs TALOS II**: Community-built, fully open-source POWER9 workstation. The only modern server-class computer with completely open hardware from firmware to chip
- **Google**: Designed OpenPOWER-compatible controllers for their data centers
- **NVIDIA**: NVLink was first demonstrated on OpenPOWER systems (before x86 NVLink)
- **IBM Open-POWER ISA 3.1 (2020)**: The ISA is now fully royalty-free and managed by the OpenPOWER Foundation

### Quick Check
> 1. What did IBM open-source in 2013 and why?
> 2. What is the Raptor Talos II and why is it significant?
> 3. How does the OpenPOWER Foundation relate to the Linux Foundation?

---

## 6. Who Uses POWER and Why

POWER systems are niche but deeply entrenched in specific industries:

**Financial services**: Major banks (Citigroup, Deutsche Bank, Goldman Sachs) run IBM Power for transaction processing, risk calculations, and core banking. Reliability and RAS (Reliability, Availability, Serviceability) features justify the premium.

**Airline reservations**: Amadeus, Sabre — global reservation systems handling millions of transactions per minute — run on IBM Power. Downtime is catastrophic; Power's RAS features (hot-swap memory, predictive failure analysis, redundant paths) make it worth the cost.

**Linux + POWER HPC**: Some national laboratories use POWER for HPC. Summit (ORNL, 2018) and Sierra (LLNL, 2018) used IBM POWER9 + NVIDIA V100 GPUs — the first time IBM and NVIDIA combined to build the #1 supercomputer.

**IBM AIX customers**: AIX is IBM's UNIX, only runs on POWER. Large enterprises that have run AIX for 30 years don't migrate easily — the software, procedures, and staff expertise are entrenched.

**IBM i (OS/400)**: An older but actively used IBM OS that runs on Power hardware. Used in mid-sized businesses for ERP and accounting. Remarkably backward-compatible — programs written in the 1980s still run.

```
POWER10 vs Intel Xeon Platinum comparison:
  
  POWER10 (15-core, SMT8, 120 threads): $50,000+ per socket
  Intel Xeon Platinum 8592+ (64-core, SMT2, 128 threads): ~$12,000 per socket
  
  At equivalent thread count:
    POWER10: $50,000 for 120 threads
    Intel: ~$25,000 for 2 sockets × 128 threads = 256 threads
  
  POWER wins on: RAS, memory capacity, AIX compatibility, IBM support
  Intel wins on: price/performance, x86 software ecosystem, core count per $
```

### Quick Check
> 1. What type of enterprise workloads specifically benefit from POWER's RAS features?
> 2. What supercomputer used IBM POWER9 + NVIDIA V100?
> 3. Why don't companies running IBM AIX just migrate to Linux on x86?

---

## Summary

- **IBM POWER** originated from the IBM 801 (1975), one of the earliest RISC designs.
- **POWER ISA**: 64-bit RISC, load-store, 32 GPRs + 32 VSX vector registers, weak memory model, big-endian historically.
- **POWER10**: up to 15 cores, 8-way SMT (120 threads/socket), Samsung 7nm, 250W+, $50K+ per socket.
- **SMT8** (8 simultaneous threads per core) enables massive throughput for database and transaction workloads at the cost of single-thread performance.
- **OpenPOWER** (2013): IBM open-sourced the POWER ISA. Full ISA documentation free, royalty-free implementations allowed.
- POWER is entrenched in financial services, airlines, and enterprises running IBM AIX/IBM i.
- Summit and Sierra supercomputers used POWER9 + NVIDIA V100 — world's fastest in 2018–2020.

---

## Exercises

### Easy
1. What is 8-way SMT? Compare to Intel's 2-way Hyper-Threading in terms of hardware threads per physical core.
2. What is AIX and why does its existence help lock customers into POWER hardware?
3. What is the OpenPOWER Foundation and what did IBM contribute to it?

### Medium
4. A database handles 10,000 concurrent connections, each with a dedicated thread. Each thread stalls 90% of the time waiting for I/O. Compare thread saturation: (a) Intel Xeon with 128 physical cores × SMT2 = 256 logical CPUs. How many threads can execute simultaneously? (b) POWER10 with 15 physical cores × SMT8 = 120 logical CPUs. How many execute simultaneously? (c) At 10% CPU utilization per thread when active: which system is underloaded/overloaded? (d) Which system is more cost-effective for this workload?
5. POWER's weak memory model allows stores to be reordered past loads. Write pseudocode demonstrating why this can cause a correctness bug in a producer-consumer program without memory barriers, then show the correct version with `lwsync`.
6. Summit supercomputer (2018) used 4,608 nodes, each with 2× IBM POWER9 CPUs + 6× NVIDIA V100 GPUs. (a) Total CPU cores? (b) Total GPU cores (V100: 5120 CUDA cores)? (c) Peak FLOPS (V100: 14 TFLOPS FP64, POWER9: ~1 TFLOPS FP64)? (d) Why use 6 GPUs per 2 CPUs rather than more CPUs?

### Hard
7. IBM POWER vs x86 total cost of ownership (TCO) analysis for a financial trading system: Input: 500 concurrent database connections, 24/7 availability, 5-year depreciation, 99.9999% uptime required (~30 seconds downtime per year). Compare: (a) IBM Power S1024 ($45K, AIX, 4×SMT8 cores, hot-swap memory, predictive failure, IBM support contract $20K/year), (b) Dell PowerEdge R760 (2× Intel Xeon Platinum, $15K, Linux, standard ECC memory, Dell ProSupport $3K/year). Calculate 5-year TCO including: hardware, support, 1 hour of downtime cost ($500K for financial trading), probability of unplanned outage.
8. The OpenPOWER ecosystem vs ARM: Both are "open" in some sense — RISC-V is the most open (ISA + implementations free), ARM has architecture licenses, POWER is now ISA-free. Rank these from most to least "open" and analyze what specifically "openness" means in each case: (a) ISA openness (can anyone implement?), (b) implementation openness (are reference designs open?), (c) ecosystem openness (tools, OS, software?), (d) governance openness (who controls the ISA evolution?). Which type of openness matters most for long-term adoption?
