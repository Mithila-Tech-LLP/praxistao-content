# Chapter 35: ARM — The Architecture That Conquered Mobile

There are more ARM processors in the world than any other type of CPU. Every smartphone, every tablet, most IoT devices, most embedded systems, and now a significant slice of laptops, servers, and even desktops run on ARM. Apple's M-series MacBooks. Amazon's AWS Graviton. Apple's iPhone A18. Qualcomm's Snapdragon. All ARM. Yet ARM Holdings, the company that designs the architecture, has fewer than 7,000 employees and doesn't manufacture a single chip. ARM is the world's most successful architecture licensing business.

## Table of Contents

1. [ARM's Origins — Acorn and the BBC Micro](#1-arms-origins--acorn-and-the-bbc-micro)
2. [The ARM Business Model — Architecture Licensing](#2-the-arm-business-model--architecture-licensing)
3. [ARM Architecture Fundamentals](#3-arm-architecture-fundamentals)
4. [The Cortex-A Series — Application Processors](#4-the-cortex-a-series--application-processors)
5. [Cortex-M — Microcontrollers Everywhere](#5-cortex-m--microcontrollers-everywhere)
6. [ARMv8-A and AArch64 — The 64-bit Leap](#6-armv8-a-and-aarch64--the-64-bit-leap)
7. [ARM's Advantage: Power Efficiency](#7-arms-advantage-power-efficiency)
8. [ARM in Servers and Data Centers](#8-arm-in-servers-and-data-centers)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. ARM's Origins — Acorn and the BBC Micro

ARM stands for **Advanced RISC Machine** (originally Acorn RISC Machine). Its origin story begins with a tiny British computer company, not a semiconductor giant.

**Acorn Computers** made the BBC Micro (1981), a popular British home computer. In 1983, Acorn needed a faster processor for their next machine. Unable to find a suitable chip from existing vendors, they designed their own.

The first **ARM1** processor was designed by a 4-person team (Sophie Wilson and Steve Furber were the primary architects) in 1985. Key design decisions:
- Pure **RISC**: simple fixed-length instructions, load/store architecture
- Extremely simple: no microcode, minimal decode logic
- Small die area: ~25,000 transistors (the 8086 had 29,000 — but was much more complex)
- Power efficient by accident: the small chip consumed so little power it didn't even need a power supply on the prototype — it ran off leakage current from other chips

The simplicity was deliberate: the team was small, the tools were primitive, and small dies were cheap. The power efficiency was an emergent property of that simplicity.

**ARM Ltd. formed (1990)**: Acorn, Apple, and VLSI Technology jointly founded Advanced RISC Machines Ltd. Apple wanted the ARM for its Newton PDA. VLSI made chips. Acorn provided the IP. The new company would license the ARM architecture rather than make chips.

### Quick Check
> 1. Where did ARM get its name?
> 2. What was the key insight that made ARM power-efficient, even if unintentional?
> 3. Why did Acorn design their own processor instead of using an existing one?

---

## 2. The ARM Business Model — Architecture Licensing

ARM's business model is unique in the semiconductor industry: **ARM designs ISAs and CPU microarchitectures, then licenses them to companies who manufacture actual chips.**

**Two levels of licensing:**

**Architecture License**: The licensee can design their own CPU core that implements the ARM ISA but can have a completely custom microarchitecture. Examples:
- Apple (A-series, M-series): Uses AArch64 ISA but their own "Firestorm" / "Everest" microarchitecture
- Qualcomm (Oryon in Snapdragon X Elite): Custom microarchitecture
- Samsung (Exynos with custom Mongoose cores)
- Amazon (Graviton: uses ARM cores + custom SoC design)

**Processor License (IP License)**: The licensee gets ARM's RTL design — a ready-to-use CPU core. They can configure it (cache sizes, pipeline depth options) but can't redesign the core. Examples:
- Qualcomm Snapdragon (uses ARM Cortex-A720 for "efficiency" cores alongside custom cores)
- Most IoT/embedded chip vendors
- Raspberry Pi (Cortex-A72/A76)

```
ARM licensing hierarchy:
  ARM Holdings
       │
       ├── Architecture License → Apple, Qualcomm (custom cores using ARM ISA)
       │
       └── IP License → Qualcomm (Cortex efficiency cores), MediaTek, Samsung, Broadcom
                         (use ARM's ready-made Cortex-A/M/R designs directly)
```

**Revenue model**: ARM charges:
- Upfront license fee
- Per-chip royalty (a few cents per chip)

With 40+ billion ARM chips shipped per year, the royalties add up even at pennies per chip. In 2023, ARM had $2.68B revenue. TSMC manufactures chips for most ARM licensees; ARM gets a royalty from each chip regardless of who makes it.

### Quick Check
> 1. What is the difference between an "architecture license" and a "processor license" from ARM?
> 2. How does ARM make money? What are the two revenue streams?
> 3. Why is ARM's business model unusual in the semiconductor industry?

---

## 3. ARM Architecture Fundamentals

ARM is a **RISC (Reduced Instruction Set Computer)** architecture. Chapter 15 covered RISC vs CISC in detail; here are the ARM-specific characteristics:

**Load-Store Architecture**: Only LOAD and STORE instructions access memory. All arithmetic is register-to-register. Contrast with x86, which allows `ADD [memory], register` directly.

**Fixed-Length Instructions**: All ARM instructions are 32 bits wide (AArch32) or 64 bits wide (AArch64). This makes decode simple and predictable, unlike x86's variable 1–15 byte instructions.

**Thumb/Thumb-2**: A 16-bit compressed instruction encoding for code density (less memory = cheaper devices). Most ARMv7-A code is a mix of 16-bit Thumb and 32-bit ARM instructions (Thumb-2). AArch64 dropped this complication.

**Conditional execution (AArch32)**: In 32-bit ARM, most instructions can be conditionally executed based on flags, without a branch instruction. Reduces branch penalties. Example: `ADDNE R0, R1, R2` — add R1+R2 to R0, but only if the "not equal" flag is set.

**Register file**: 16 general-purpose 32-bit registers in AArch32 (R0–R15, where R15 = PC and R14 = LR); 31 general-purpose 64-bit registers (X0–X30) in AArch64.

**Calling convention (AArch64)**: The first 8 arguments in X0–X7, return value in X0, X8 as indirect return, X9–X15 caller-saved, X19–X28 callee-saved. Much cleaner than x86-64's mix of registers and stack.

**NEON and SVE**: ARM's SIMD extensions. NEON (128-bit, ARMv7-A/AArch64) is the baseline. SVE (Scalable Vector Extension, AArch64) allows variable-length vectors — the same code runs on SVE hardware with 128-bit, 512-bit, or 2048-bit registers.

### Quick Check
> 1. What does "load-store architecture" mean in ARM?
> 2. Why does fixed-length instructions make decode simpler for ARM than for x86?
> 3. What is SVE and why is "scalable" an advantage over fixed-width SIMD?

---

## 4. The Cortex-A Series — Application Processors

The **Cortex-A** series is ARM's family of high-performance application processors — designed for smartphones, tablets, laptops, and servers.

**Cortex-A naming**: Higher numbers are generally more recent (not necessarily better). A53 and A55 are "efficiency" cores (power-efficient, simpler OOO); A72/A76/A78/A710/A720 are "performance" cores (wider, deeper OOO).

Key milestones:

**Cortex-A9 (2008)**: First ARM OOO processor. 2-issue superscalar. Enabled ARM to compete with Intel Atom for lightweight laptops.

**Cortex-A15 (2011)**: 3-issue OOO, 60-entry ROB. Used in Samsung Galaxy S4 Exynos. First ARM that could genuinely run "real" desktop workloads.

**Cortex-A57 (2014)**: First AArch64 (64-bit) Cortex. 3-issue OOO, 128-entry ROB. Used in Nvidia Tegra X1, Samsung Exynos 7420, early Raspberry Pi 4 alternatives.

**Cortex-A72/A73 (2016/2017)**: Widely used, balanced performance and efficiency. Found in Raspberry Pi 4 (A72), many Android mid-range phones.

**Cortex-A76 (2018)**: Significant performance leap — designed to match Intel Core i5 in laptop workloads at a fraction of the power. This was ARM's "mobile can replace laptop" moment.

**Cortex-A77/A78 (2019/2020)**: Refinements of A76 with improved branch prediction, wider OOO.

**Cortex-X series (2020+)**: ARM's "premium" cores, wider than standard Cortex-A for maximum single-thread performance at higher power. Cortex-X1 in Samsung Galaxy S21, Cortex-X2/X3/X4 in subsequent flagship phones.

**big.LITTLE (2011+)**: ARM's heterogeneous computing solution for smartphones. Pair power-hungry "big" cores (Cortex-A57/A72/A76) with power-efficient "little" cores (Cortex-A53/A55). The OS scheduler routes tasks to appropriate cores. A text message notification wakes an A55; a game uses all A76 cores. The combination gives both performance and all-day battery life.

```
Typical ARM SoC topology (flagship 2023):
  Cortex-X3 (1 core):   3.36 GHz, 256KB L1, 8MB L2   — peak performance
  Cortex-A715 (4 cores): 2.96 GHz, 64KB L1, 2MB L2   — everyday tasks
  Cortex-A510 (4 cores): 2.27 GHz, 64KB L1, 512KB L2  — background tasks
  
  This is the DynamIQ arrangement in Snapdragon 8 Gen 2 configuration
```

### Quick Check
> 1. What is "big.LITTLE" and why does it matter for smartphone battery life?
> 2. What made the Cortex-A76 a milestone for ARM laptop competitiveness?
> 3. What is the Cortex-X series and how does it differ from standard Cortex-A?

---

## 5. Cortex-M — Microcontrollers Everywhere

If Cortex-A is the smartphone brain, **Cortex-M** is the nervous system of the embedded world. Cortex-M processors run in microcontrollers — the tiny computers inside light switches, medical devices, automotive systems, and industrial sensors.

**Cortex-M key characteristics:**
- Very small die area (Cortex-M0: ~12,000 gates — fits on tiny chips)
- Deterministic, low-latency interrupt handling (Nested Vectored Interrupt Controller, NVIC)
- No MMU (no virtual memory — runs directly on physical addresses, typical for embedded)
- No caches in simplest variants (predictable timing = critical for real-time systems)
- Ultra-low power: Cortex-M0+ consumes ~9 µA/MHz

**Cortex-M hierarchy:**
| Core | Width | FPU | DSP | Use Case |
|------|-------|-----|-----|----------|
| Cortex-M0/M0+ | 32-bit, 3-stage | No | No | Simplest IoT, button sensing |
| Cortex-M3 | 32-bit, 3-stage | No | Partial | General embedded |
| Cortex-M4 | 32-bit, 3-stage | Optional | Yes | Audio, motor control |
| Cortex-M7 | 32-bit, 6-stage, limited OOO | Yes | Yes | High-perf embedded |
| Cortex-M33 | 32-bit, 3-stage | Optional | Yes | IoT with security |
| Cortex-M55 | 32-bit | Yes | Yes | ML at the edge |

**Example chips using Cortex-M:**
- STM32 (ST Microelectronics): Cortex-M0/M3/M4/M7 in a huge family of microcontrollers for industrial and consumer electronics
- Arduino Due: Cortex-M3
- Nordic Semiconductor nRF52840: Cortex-M4 for Bluetooth LE wearables
- SAMD21 (in Arduino Zero, original Adafruit boards): Cortex-M0+
- Apple AirPods: Cortex-M (for always-on audio processing while the main SoC sleeps)

### Quick Check
> 1. Why don't Cortex-M processors have MMUs or caches in their simplest variants?
> 2. What is the NVIC and why is it important for real-time embedded systems?
> 3. Name a product you use daily that likely contains a Cortex-M processor.

---

## 6. ARMv8-A and AArch64 — The 64-bit Leap

In 2011, ARM released **ARMv8-A** — the most significant architectural revision since the original ARM. The key addition: **AArch64**, a completely new 64-bit execution state.

Unlike Intel's messy 16→32→64-bit evolution (which carried all the baggage forward), ARM designed AArch64 as a clean new ISA:

**AArch64 improvements over AArch32:**
- 31 general-purpose 64-bit registers (X0–X30) — vs 16 in AArch32
- Separate stack pointer (SP) register, not overloaded as R13
- No predicated execution (removed the conditional instruction encoding)
- No IT blocks (Thumb2's conditional blocks)
- Improved SIMD/FP (all ARMv8 FP/SIMD mandatory, not optional)
- PC is no longer accessible as a general-purpose register
- Simplified addressing modes

AArch64 was backward-compatible at the **binary level** via hardware: an ARMv8-A CPU can switch between AArch32 mode (running old 32-bit ARM code) and AArch64 mode (running 64-bit code) at runtime. This enabled a smooth transition.

**Apple was the first**: The iPhone 5S (2013) shipped with the **Apple A7** — the world's first consumer AArch64 processor. Apple was 18 months ahead of Android chips in 64-bit adoption. The A7 also signaled Apple's intent to build world-class custom ARM microarchitectures.

**ARMv9 (2021)**: Added SVE2 (Scalable Vector Extension 2), Realm Management Extension (hardware security domains), and improved ML capabilities. Cortex-X3 and Apple M2 are among the first ARMv9 implementations.

### Quick Check
> 1. How many general-purpose registers does AArch64 have? How many did AArch32 have?
> 2. Why was AArch64 designed as a clean new ISA rather than extending AArch32?
> 3. What was historic about the Apple A7 chip in the iPhone 5S?

---

## 7. ARM's Advantage: Power Efficiency

ARM's defining advantage is **power efficiency** — more computation per watt than x86.

**Why ARM is more power efficient:**
1. **RISC simplicity**: Simpler decode hardware, no CISC-to-µop translation overhead. Fewer transistors doing the same work.
2. **Load-store ISA**: Cleaner datapath with fewer memory-dependent instruction types.
3. **Mobile-first design**: ARM's reference designs have been optimized for power from the beginning, not retrofitted.
4. **Smaller pipeline stages**: Simpler decode → shorter critical paths → lower voltage for the same speed.

**The numbers (rough comparisons, 2023):**
```
  Apple M3 Pro (ARM): 11.5 TOPS/W AI performance, ~30W TDP
  Intel Core i9-13900H (x86): ~6 TOPS/W, 45W TDP
  
  ARM is approximately 2-3× more power efficient for similar workloads
  (though direct comparisons are complex due to different optimization targets)
```

**ARM's power efficiency is not absolute**: Apple's Firestorm (M1/M2 "big" core) consumes as much as 15-20W at peak — not fundamentally different from an Intel P-core. The advantage comes from the efficiency cores (Icestorm) consuming ~0.5–1W and the whole-system approach (tight integration, unified memory, efficient I/O).

The power advantage matters most in:
- **Smartphones**: Battery life is directly constrained by power
- **Laptops**: Fanless designs, all-day battery
- **Data centers**: PUE (Power Usage Effectiveness) — every watt of computing costs 1.3–1.8W total (compute + cooling). AWS Graviton saves meaningful money at scale.

### Quick Check
> 1. List three architectural reasons ARM tends to be more power-efficient than x86.
> 2. Why does power efficiency matter even more in a data center than in a laptop?
> 3. Is ARM always more power-efficient than x86? When might x86 be competitive?

---

## 8. ARM in Servers and Data Centers

The "ARM is only for phones" narrative ended around 2018–2020. ARM is now a serious data center architecture.

**Amazon Graviton (2018+)**: Amazon's internal ARM-based server CPU. Graviton3 (2021) uses Cortex-X2-like custom cores. AWS offers Graviton instances at ~20% lower price than equivalent x86 instances with comparable or better performance for many workloads (web servers, compiled code, databases).

**Ampere Altra (2020+)**: A fabless startup using ARM's Neoverse N-series architecture. Ampere Altra Max: 128 Cortex-A78 cores at 3.0 GHz. Oracle Cloud, Microsoft Azure use Ampere instances. The value proposition: many cores at good efficiency.

**ARM Neoverse**: ARM's server-focused CPU designs (Neoverse N-series for throughput, E-series for ultra-high core count, V-series for HPC). Used in AWS Graviton, Ampere Altra, and Chinese domestic server chips.

**Fujitsu A64FX**: The Fugaku supercomputer (Japan, #1 TOP500 in 2020) used custom ARM processors with SVE (512-bit vectors). 158,976 nodes × 48 cores = 7.6 million cores. Fugaku ran on ARM — a landmark moment for ARM HPC credibility.

**Microsoft Azure Cobalt (2024)**: Microsoft's first custom ARM server chip for Azure.

The server ARM market is growing but still dominated by x86 (Intel Xeon, AMD EPYC). Legacy enterprise software (Oracle databases, SAP HANA, Windows Server) runs on x86 by default. Open-source/cloud-native software (Linux, containers, Go, Python, Java) is easily ported to ARM. The transition will be gradual but is clearly underway.

### Quick Check
> 1. What is Amazon Graviton and why did Amazon design their own ARM server chip?
> 2. What was the significance of the Fugaku supercomputer using ARM?
> 3. What type of software transitions most easily to ARM servers, and what doesn't?

---

## Summary

- **ARM** began in 1985 at Acorn Computers — designed by a tiny team for simplicity, with power efficiency as an emergent benefit.
- ARM's **licensing model** means ARM Holdings designs architectures and cores, which chip companies (Apple, Qualcomm, Samsung) license to build products. Two levels: architecture license (custom microarchitecture) and IP license (ready-made Cortex cores).
- **Cortex-A** (smartphones/laptops), **Cortex-M** (microcontrollers), **Cortex-R** (real-time), **Cortex-X** (premium performance) serve different market segments.
- **AArch64** (ARMv8-A, 2011) was a clean 64-bit redesign — 31 GPRs, no legacy cruft, simpler addressing.
- ARM's core advantage: **power efficiency** — RISC simplicity, load-store ISA, mobile-first design philosophy.
- ARM is expanding beyond mobile: AWS Graviton, Ampere Altra, Fujitsu A64FX, and Microsoft Cobalt prove ARM's server viability.

---

## Exercises

### Easy
1. What is ARM's business model? How does ARM make money without making chips?
2. What is "big.LITTLE" in ARM terminology? Give a concrete example of how it helps smartphone battery life.
3. Why did AArch64 get 31 general-purpose registers instead of just extending AArch32's 16 registers?

### Medium
4. Analyze ARM's licensing model vs Intel's IDM (Integrated Device Manufacturer) model: (a) What are ARM's advantages (asset-light, scale to many customers)? (b) What are ARM's disadvantages (no control over manufacturing quality, no hardware revenue from chips)? (c) Under what market conditions could Intel's model outperform ARM's?
5. An IoT device has a 1000 mAh battery and uses a Cortex-M0+ running at 10MHz, consuming 100 µA. A developer wants to process sensor data every second. (a) How many cycles per second are available? (b) How many hours will the battery last at this duty cycle? (c) If the chip spends 99% of time in deep sleep (10 µA) and wakes up for 10ms every second, what is the effective current draw?
6. ARM SVE vectors are "scalable" — code compiled with SVE-128 runs correctly (but more slowly) on SVE-512 hardware, and the same compiled binary can leverage wider vectors on better hardware. Compare this to x86 AVX-512: code compiled targeting AVX-512 won't run at all on CPUs without AVX-512. What architectural principle makes SVE more portable, and what is the trade-off?

### Hard
7. Amazon Graviton3 vs Intel Xeon Platinum 8380 (Ice Lake) for a web serving workload: both at approximately the same price per instance. Research (or reason about): (a) How many physical cores does each have? (b) How does single-thread performance compare? (c) For a workload that scales with core count up to 64 cores, which wins and why? (d) What workloads favor the Intel Xeon over Graviton3?
8. Apple achieved battery life parity or superiority in M1 MacBooks vs Intel MacBooks, despite having comparable or higher single-thread performance. Analyze the system-level factors beyond just CPU microarchitecture efficiency: (a) unified memory (no CPU-GPU data copying), (b) Icestorm efficiency cores for background tasks, (c) on-chip media engine (avoiding CPU-based video decode), (d) tighter OS-hardware integration. Estimate the power savings from each factor for a typical "web browsing + video" usage pattern.
