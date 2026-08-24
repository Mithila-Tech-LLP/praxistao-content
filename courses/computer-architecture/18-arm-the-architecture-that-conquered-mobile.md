# Chapter 18: ARM — The Architecture That Conquered Mobile

ARM is the most widely deployed processor architecture in history. If you count every chip ever manufactured, ARM wins by billions — every iPhone, every Android phone, most IoT devices, billions of embedded chips, and now a growing number of laptops (Apple M-series) and servers (AWS Graviton, Ampere). This chapter traces ARM from its origins in a tiny British team to its global dominance.

## Table of Contents

1. [ARM's Origins: A Humble Beginning](#1-arms-origins-a-humble-beginning)
2. [The ARM Business Model](#2-the-arm-business-model)
3. [ARM7, ARM9, ARM11: The Classic Era](#3-arm7-arm9-arm11-the-classic-era)
4. [ARM Cortex-A: High-Performance Application Processors](#4-arm-cortex-a-high-performance-application-processors)
5. [AArch64: The Clean 64-bit ISA](#5-aarch64-the-clean-64-bit-isa)
6. [ARM Cortex-M: Microcontrollers](#6-arm-cortex-m-microcontrollers)
7. [ARM Cortex-R: Real-Time Processors](#7-arm-cortex-r-real-time-processors)
8. [NEON and Advanced SIMD](#8-neon-and-advanced-simd)
9. [TrustZone: Hardware Security](#9-trustzone-hardware-security)
10. [Apple Silicon: ARM Done Extraordinarily Well](#10-apple-silicon-arm-done-extraordinarily-well)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. ARM's Origins: A Humble Beginning

### The BBC Micro Connection

In 1983, Acorn Computers (Cambridge, UK) won a BBC contract to supply a home computer for the UK national computer literacy program. Their BBC Micro was successful, but Acorn needed a 32-bit processor for a successor machine.

No existing processor was suitable. The team — Sophie Wilson (ISA design) and Steve Furber (hardware) — decided to design their own. With just 4 engineers, they created the ARM1 (Acorn RISC Machine) in 1985:

```
ARM1 Specifications:
  Transistors:    25,000 (compare: Intel 8086 had 29,000 in 1978)
  Clock speed:    8 MHz
  Performance:    4-8 MIPS (million instructions per second)
  Power:          ~100 mW
  Process:        VLSI Technology 3-micron CMOS
```

The tiny transistor count was intentional. Acorn's RISC team had read Patterson and Sequin's RISC papers and embraced the simple-is-fast philosophy. The low transistor count meant low power consumption — a property that would prove invaluable decades later for mobile devices.

### ARM Ltd. Founded (1990)

Acorn, Apple, and VLSI Technology jointly founded Advanced RISC Machines Ltd. in 1990. Apple was developing the Newton PDA and needed an efficient processor. The ARM710 appeared in the Apple Newton MessagePad (1993).

The founding as an IP licensing company — not a chip manufacturer — was the key business innovation (see next section).

### Quick Check

> 1. What was the ARM1's transistor count and how does it compare to the Intel 8086?
> 2. Why was ARM's low power consumption important for its later success?
> 3. Who were the two main designers of the ARM1 ISA and hardware?

---

## 2. The ARM Business Model

### IP Licensing vs Manufacturing

ARM Ltd. does not manufacture processors. It designs ISAs and processor cores, then licenses them to chip makers who embed ARM IP in their own designs:

```
ARM Ltd.          →  licenses ISA/core  →  Apple
                                          Qualcomm
                                          Samsung
                                          MediaTek
                                          NXP
                                          Texas Instruments
                                          ... 200+ licensees
                   
Each licensee builds their own chip using ARM technology, sells to OEMs.
```

### License Types

**Architecture License**: The licensee can implement the ARM ISA in any way they want, including custom microarchitectures. Apple uses this — they design their own cores (Swift, Cyclone, Typhoon, Hurricane, Monsoon, Vortex, Firestorm, Avalanche, Everest, Sawtooth) that implement ARM ISA.

**Core License**: The licensee gets an ARM-designed core (Cortex-A76, Cortex-M4, etc.) that they can integrate into their SoC but cannot modify. Most companies use this — designing a core takes 100+ engineers and years.

**Per-Chip Royalty**: ARM charges a royalty per chip sold, typically 0.5-2% of chip selling price. With billions of ARM chips sold annually, royalties are enormous.

### Why This Worked

- Zero manufacturing risk for ARM
- Network effect: more licensees → more software → more value → more licensees
- Every layer of the ecosystem (chip makers, OEMs, app developers) benefits from ARM compatibility
- ARM ISA standardization means a single compiler (GCC, LLVM/Clang) works for all ARM chips

### Quick Check

> 1. What does "IP licensing" mean in the context of ARM?
> 2. What is the difference between an architecture license and a core license?
> 3. How does ARM make money per chip sold?

---

## 3. ARM7, ARM9, ARM11: The Classic Era

### ARM7TDMI (1994)

Became one of the best-selling processor cores in history. Used in:
- Game Boy Advance (2001)
- Nokia 3310 and hundreds of feature phones
- Early iPod
- Countless embedded systems

The "TDMI" suffix indicates extensions:
- T: Thumb 16-bit instruction set (code density)
- D: JTAG Debug support
- M: 8×8→16 bit fast multiply
- I: Embedded ICE debug module

### ARM9 (1997)

5-stage pipeline (vs ARM7's 3-stage), Harvard architecture (separate instruction and data caches), performance ~200 MIPS at 200 MHz. Used in Nintendo DS, original iPhone's baseband processor.

### ARM11 (2002)

8-stage pipeline, ARMv6 ISA with SIMD extensions (media processing), hardware divider, 16KB L1 caches. Used in the original iPhone (2007) at 620 MHz — the ARM11 gave the iPhone its smoothness compared to feature phones.

### The ARM ISA Versions

Confusingly, ARM uses two separate numbering systems:
- **Processor series**: ARM7, ARM9, ARM11, Cortex
- **ISA version**: ARMv4, ARMv5, ARMv6, ARMv7, ARMv8, ARMv9

ARM11 implements ARMv6. Cortex-A8 implements ARMv7-A. Apple M1 implements ARMv8.5-A.

### Quick Check

> 1. Name three commercial products that used the ARM7TDMI.
> 2. What does the ARM ISA version (ARMv6, ARMv7) specify, and how is it different from the processor series number?
> 3. What key feature did the ARM9 add over ARM7?

---

## 4. ARM Cortex-A: High-Performance Application Processors

The Cortex-A series targets application processors — the main CPU in smartphones, tablets, and (now) laptops.

### Cortex-A8 (2005): iPhone 3GS

The first major Cortex application processor. 13-stage pipeline, in-order, NEON SIMD, 512KB L2. Used in Apple A4 (iPhone 4).

### Cortex-A9 (2007): First Dual-Core ARM

Added out-of-order execution (OOO) — a major architectural leap. Superscalar (2 instructions/cycle). Used in Apple A5 (iPhone 4S, iPad 2), Qualcomm Snapdragon S3, Samsung Exynos 4210.

### Cortex-A15 (2010): ARMv7-A High Performance

4-wide OOO, deep 15-stage pipeline, 40-bit physical address space. Big.LITTLE pairing with Cortex-A7 for power efficiency (run A7 for simple tasks, wake A15 for demanding ones). Used in Samsung Exynos 5.

### Cortex-A53/A57 (2012): ARMv8-A (64-bit)

First 64-bit Cortex cores. A53: efficient in-order; A57: high-performance OOO. Became the foundation of mobile 64-bit computing. Apple had already introduced 64-bit ARM with the A7 (2013, iPhone 5s) using their own core design (Cyclone).

### Cortex-A76/A78/X1 (2018-2020): PC-class Performance

The Cortex-A76 (2018) was ARM's first explicit claim to "laptop-class" performance. 4-wide superscalar, 13-stage OOO pipeline, AVX equivalent SIMD. Performance per watt began seriously challenging Intel laptop processors.

### Performance-per-Watt Comparison

```
Year  ARM Core    Device        Performance  Power   Perf/Watt
────  ──────────  ────────────  ───────────  ─────   ─────────
2007  Cortex-A8   iPhone 3GS     ~600 MIPS   ~250mW   2.4 MIPS/mW
2013  Apple Cyc   iPhone 5s     ~2500 MIPS   ~500mW   5.0 MIPS/mW
2020  Apple Fire  iPhone 12     ~9000 MIPS   ~1W      9.0 MIPS/mW
2020  Apple M1    MacBook Air   ~24000 MIPS  ~5W      4.8 MIPS/mW (laptop)
```

### Quick Check

> 1. What key microarchitecture feature did Cortex-A9 add over Cortex-A8?
> 2. What does "big.LITTLE" mean in the context of ARM processors?
> 3. When did ARM (via Apple) first offer 64-bit ARM processors in smartphones?

---

## 5. AArch64: The Clean 64-bit ISA

ARM64 (officially called AArch64) is the 64-bit ISA introduced with ARMv8-A (2011) and is a significant improvement over 32-bit ARMv7:

### Register File

```
AArch64 registers:
  x0-x30    31 × 64-bit general-purpose registers
             (w0-w30 are 32-bit views of the same registers)
  xzr/wzr   Zero register (hardwired 0, like RISC-V x0)
  sp        Stack pointer (x31 used implicitly)
  pc        Program counter (not directly readable in most instructions)
  
  v0-v31    32 × 128-bit SIMD/FP registers
             (b0 = 8-bit, h0 = 16-bit, s0 = 32-bit, d0 = 64-bit views)
  
  System:   NZCV (flags), FPCR, FPSR, ELR_EL1, SPSR_EL1, ...
```

### Clean Design Improvements vs ARMv7

| Feature | ARMv7 (32-bit) | AArch64 (64-bit) |
|---------|---------------|-----------------|
| GPRs | 16 (r0-r15) | 31 (x0-x30) |
| Conditional execution | Every instruction | Branch instructions only |
| Instruction width | Variable (16+32 Thumb) | Fixed 32-bit |
| PC access | GPR (r15=PC!) | Restricted |
| Inline barrel shifter | All instructions | Separate shift instruction |
| Load/store multiple | LDMIA/STMDB | LDP/STP (pair loads) |

Removing conditional execution from all instructions was controversial but rational — modern out-of-order processors handle branches efficiently with prediction, so predication provides little benefit and adds complexity to the decoder.

### AArch64 Calling Convention (AAPCS64)

```
Arguments:   x0-x7 (8 integer args), v0-v7 (8 float/SIMD args)
Return:      x0-x1 (up to 128-bit), v0
Caller-saved: x0-x18 (caller saves if needed across calls)
Callee-saved: x19-x28, x29 (frame pointer), x30 (link register = ra)
Stack:       16-byte aligned, grows downward
```

### Quick Check

> 1. How many general-purpose registers does AArch64 have vs ARMv7?
> 2. What is a "zero register" in AArch64 and what is it called?
> 3. Why did AArch64 remove conditional execution from most instructions?

---

## 6. ARM Cortex-M: Microcontrollers

The Cortex-M family targets embedded microcontrollers: devices where code size and power matter far more than peak performance.

### Cortex-M Profile Variants

```
Core       ISA         Pipeline  Frequency  Target
─────────  ──────────  ────────  ─────────  ───────────────────────────────
Cortex-M0  ARMv6-M     2-stage    48 MHz    Ultra-low power IoT
Cortex-M0+ ARMv6-M     2-stage    48 MHz    Lowest power, 1-cycle GPIO
Cortex-M3  ARMv7-M     3-stage   120 MHz    General embedded
Cortex-M4  ARMv7E-M    3-stage   180 MHz    DSP ops, optional FPU
Cortex-M7  ARMv7E-M    6-stage   600 MHz    High performance embedded
Cortex-M33 ARMv8-M     3-stage   200 MHz    Secure IoT (TrustZone)
Cortex-M55 ARMv8.1-M   3-stage   400 MHz    ML at edge (Helium SIMD)
```

### Thumb-2: The ISA That Makes Cortex-M Practical

Cortex-M uses Thumb-2 ISA — a mix of 16-bit (Thumb) and 32-bit (ARM) instructions. This gives near-full ARM performance (~98%) with near-Thumb code density (~26% smaller than 32-bit ARM). Thumb-2 is why embedded ARM is so efficient — you get the benefits of both worlds.

### Cortex-M in Real Products

- Cortex-M0+: Arduino Zero, many sensors
- Cortex-M4: STM32F4 (Arduino-compatible, used in countless projects), Nordic nRF52 (Bluetooth)
- Cortex-M33: Nordic nRF9160 (LTE-M cellular IoT), STM32L5 (TrustZone security)
- Cortex-M7: STM32H7 (microcontroller running at 480 MHz — laptop performance in an embedded chip!)

### Quick Check

> 1. What is the target use case for Cortex-M vs Cortex-A processors?
> 2. What is Thumb-2 and why is it important for embedded systems?
> 3. Name a common microcontroller board that uses the Cortex-M4.

---

## 7. ARM Cortex-R: Real-Time Processors

Cortex-R is designed for **hard real-time** applications: systems where late responses are as bad as wrong responses.

Examples: automotive braking control, hard disk drive read/write controllers, medical device pacemakers, industrial robots.

Key features for real-time:
- **Tightly coupled memory (TCM)**: RAM directly accessible in 0 cycles (no cache misses possible)
- **Lockable caches**: critical code sections can be pinned in cache
- **ECC memory**: error-correcting code detects and corrects single-bit memory errors
- **Deterministic interrupt latency**: worst-case interrupt response time is bounded and specified
- **Dual-core lockstep (Cortex-R52)**: two cores execute the same instructions simultaneously; if outputs diverge, a hardware fault is detected

Cortex-R processors don't run Linux (which has unbounded latencies). They run RTOS (Real-Time Operating Systems) like FreeRTOS, QNX, SafeRTOS.

---

## 8. NEON and Advanced SIMD

ARM NEON (introduced in ARMv7-A, Cortex-A8, 2005) is ARM's SIMD extension:

```
NEON register file: 32 × 128-bit D registers (or 16 × 128-bit Q registers)
Operations per NEON instruction (Q register):
  16 × 8-bit integer operations
   8 × 16-bit integer operations  
   4 × 32-bit integer operations
   4 × 32-bit float operations
   2 × 64-bit integer operations
```

### NEON for Image Processing

```
# Process 16 pixels of an 8-bit grayscale image simultaneously
vld1.8  {d0, d1}, [r0]!     # load 16 bytes (pixels) from memory
vmull.u8 q0, d0, d1          # multiply 8 pairs of 8-bit values
vst1.8  {d0, d1}, [r1]!     # store 16 results
# 3 instructions, 16 pixels processed → 5.3× more throughput than scalar
```

NEON is critical for video codecs (H.264, H.265), image filters, audio codecs, and machine learning inference on phones.

### SVE and SVE2 (Scalable Vector Extension)

ARMv8.2-A introduced SVE — a radical departure from fixed-width SIMD:

```
SVE vectors: variable length (128 to 2048 bits, in 128-bit increments)
             The ISA doesn't specify the width — the hardware implements it

SVE registers: z0-z31 (vector), p0-p15 (predicate masks)
```

Why variable length? Code compiled for SVE runs on hardware with any SVE vector length. A chip with 256-bit SVE and a chip with 2048-bit SVE run the same code — the wider chip just processes more elements per instruction. No recompilation needed when hardware improves.

Compare to AVX-512 (x86): if a program uses AVX-512, it requires hardware with exactly 512-bit vectors. SVE is more forward-compatible.

### Quick Check

> 1. How many 8-bit integer operations can a single NEON instruction perform on Q registers?
> 2. What does "scalable" mean in Scalable Vector Extension?
> 3. Why is SVE more forward-compatible than AVX-512?

---

## 9. TrustZone: Hardware Security

ARM TrustZone (introduced ARMv6-Z, fully in ARMv7-A) is a hardware security feature built into the processor:

### Two Worlds

TrustZone divides the system into two isolated worlds:
- **Normal World**: runs the rich OS (Linux/Android/iOS) and user applications
- **Secure World**: runs a Trusted Execution Environment (TEE) with high-security code

```
Normal World                    Secure World
────────────────────────────    ──────────────────────────────────
Android OS                      OP-TEE (Open Source TEE OS)
Apps (Gmail, Chrome, games)     Trusted Apps (TAs):
                                  - Fingerprint template matching
                                  - DRM key storage
                                  - Mobile payment (Apple Pay, GPay)
                                  - Device attestation
                                  - Secure boot verification
```

### How TrustZone Works

A single physical CPU core runs both worlds, switching between them using the **Secure Monitor Call** (SMC) instruction:

```
Normal World CPU state:
  - NS bit in SCR_EL3 = 1 (non-secure)
  - Access to: Normal World DRAM, peripherals marked NS=1
  - Cannot access: Secure World memory, secure peripherals

Secure World CPU state:
  - NS bit = 0
  - Access to: ALL memory (can read Normal World memory too)
  - Peripherals: TrustZone hardware enforces access control

Transition:
  Normal → Secure: smc instruction → Secure Monitor → Secure OS
  Secure → Normal: eret instruction from Secure OS
```

### Real-World Use

When you pay with Apple Pay or Google Pay, your card token is stored in the Secure World. The payment app (Normal World) requests a payment, but the actual card number never leaves the Secure World — it's validated there and only an approval/denial is returned.

When a phone reads your fingerprint, the biometric template stored in Secure World is never accessible to Normal World apps — even a compromised Android OS cannot steal fingerprint data.

### Quick Check

> 1. What are the two worlds in ARM TrustZone?
> 2. What instruction switches from Normal World to Secure World?
> 3. Name two real applications of ARM TrustZone in smartphones.

---

## 10. Apple Silicon: ARM Done Extraordinarily Well

### The A7 Surprise (2013)

When Apple introduced the iPhone 5s with the Apple A7 processor, they did something unprecedented: a 64-bit ARM chip in a consumer phone in 2013, when the rest of the industry (Qualcomm, Samsung) was still 32-bit. Apple had implemented their own ARM core (Cyclone) using an architecture license.

The industry was caught completely off guard. ARM64 (AArch64) support wasn't even complete in Android until 2014.

### The M1 Revolution (2020)

The Apple M1 (November 2020) ended the assumption that laptops needed x86:

```
Apple M1 Specifications:
  Process:        TSMC 5nm
  Transistors:    16 billion
  CPU:            4 × Firestorm (performance) + 4 × Icestorm (efficiency)
  GPU:            8-core (7 in Air)
  Neural Engine:  16-core, 11 TOPS
  Memory:         Unified memory (up to 16GB), 68.25 GB/s bandwidth
  
Performance vs Intel Core i7 (MacBook Pro 2020, 15W TDP):
  Single-thread:  +50% faster
  Multi-thread:   +70% faster
  Power:          M1 at 7W = Intel at 28W (4× more efficient)
```

### Why M1 is So Fast

**1. Wide superscalar**: M1 Firestorm can decode 8 instructions per cycle (Intel: 4, AMD Zen 3: 4). This feeds more work into the execution engine.

**2. Huge reorder buffer**: 630 entries (Intel Skylake: 224, AMD Zen 3: 256). More entries means more instructions can be in-flight simultaneously, hiding memory latency.

**3. Unified memory architecture**: CPU, GPU, and Neural Engine all share the same physical DRAM, connected via extremely high-bandwidth buses. No data copying between CPU memory and "GPU memory" — they use the same pool.

**4. High-bandwidth cache**: L1 cache per Firestorm core: 192KB instruction + 128KB data (2× Intel). L2: 12MB shared. L3 equivalent: 8-24MB system cache.

**5. Memory bandwidth**: 68.25 GB/s (M1) vs ~50 GB/s (Intel 11th gen). Critical for ML workloads.

**6. Custom silicon for everything**: The M1 has dedicated hardware for video encode/decode, signal processing, cryptography, neural network inference — offloading these from the CPU cores.

### M-Series Evolution

```
Chip   Year  CPU Cores   GPU Cores  Transistors  Memory BW
─────  ────  ─────────── ─────────  ───────────  ─────────
M1     2020  4P+4E        8          16B           68 GB/s
M1 Pro 2021  8P+2E       16          33B           200 GB/s
M1 Max 2021  8P+2E       32          57B           400 GB/s
M1 Ultra 2022 16P+4E     64         114B           800 GB/s
M2     2022  4P+4E        10          20B           100 GB/s
M3     2023  4P+4E/...    10          25B           100 GB/s
M4     2024  4P+6E/...    10          28B           120 GB/s
```

Apple's approach (custom ARM cores, unified memory, tight hardware-software integration) has set the standard for what ARM processors can achieve — forcing Intel and AMD to dramatically improve their power efficiency.

---

## Summary

- ARM started in 1985 with 25,000 transistors as an academic RISC experiment; it now ships billions of chips per year.
- ARM's business model (IP licensing) created a huge ecosystem: 200+ licensees, every major phone chip, IoT device, and now laptops.
- Cortex-A: high-performance application processors (phones, laptops, servers). Cortex-M: microcontrollers (IoT, embedded). Cortex-R: real-time systems (automotive, hard disks).
- AArch64 (ARM64) is a clean 64-bit ISA with 31 GPRs, fixed 32-bit encoding, no condition codes on most instructions.
- NEON provides 128-bit SIMD; SVE provides scalable-width SIMD (128-2048 bits, same code runs on all widths).
- TrustZone provides hardware isolation between Normal World (Android/iOS) and Secure World (TEE for payments, biometrics).
- Apple Silicon (M1-M4) demonstrates what ARM can achieve with full hardware/software integration: 2-4× better power efficiency than x86 at equivalent performance levels.

---

## Exercises

### Easy

1. List the three ARM Cortex profiles and their target applications.

2. What is ARM's business model? Who designs ARM chips and who manufactures them?

3. How many GPRs does AArch64 have vs ARMv7?

### Medium

4. **Big.LITTLE scheduling**: A Cortex-A75 (big) core consumes ~3W at peak; a Cortex-A55 (LITTLE) core consumes ~0.3W. If an app spends 80% of time doing simple background tasks (battery check, location update) and 20% doing intensive computation (video decode), calculate: (a) total power with always-on big core, (b) total power with always-on LITTLE core, (c) total power with perfect big.LITTLE scheduling. What is the battery life ratio between (a) and (c)?

5. ARM TrustZone is available on Cortex-A but NOT on Cortex-M0. The Cortex-M33 adds TrustZone support. Research: what is the difference between TrustZone on Cortex-A (TEE) vs TrustZone on Cortex-M33 (TrustZone-M)? Why are two different implementations needed?

6. **Unified memory vs discrete GPU memory**: The Apple M1 uses unified memory shared between CPU and GPU. A traditional desktop uses separate CPU DRAM and discrete GPU VRAM. For a deep learning inference task that requires transferring 4GB of model weights from storage to GPU: (a) in a traditional system with PCIe 4.0 (32 GB/s peak), how long does the transfer take? (b) In the M1's unified memory, is there any transfer at all? What is the implication for latency-sensitive ML applications?

### Hard

7. **Architecture license deep dive**: Apple uses an ARM architecture license (they implement the ARM ISA but design their own microarchitecture). Qualcomm also used to use an architecture license for their Kryo cores. In 2019, ARM terminated Qualcomm's architecture license (due to the Qualcomm-NXP acquisition dispute), forcing Qualcomm to revert to standard Cortex cores. Research: (a) What does an ARM architecture license grant vs a core license? (b) What happened to Qualcomm's performance position relative to Apple after switching from custom Kryo to Cortex? (c) Why would a company pay extra for an architecture license vs just using ARM's Cortex designs? (d) RISC-V was designed specifically to avoid this licensing risk. How does this affect the market?

8. **Rosetta 2 performance analysis**: When Apple transitioned from Intel x86 to Apple Silicon (ARM), they provided Rosetta 2 binary translation. Research and measure (or find benchmarks): (a) What technique does Rosetta 2 use (JIT compilation, AOT compilation, or both)? (b) For CPU-bound code, what is the typical performance penalty of Rosetta 2 vs native ARM code? (c) For I/O-bound code? (d) Why does some x86 code actually run faster under Rosetta 2 on M1 than natively on the original Intel Mac? (hint: M1 hardware advantage). (e) What happens when x86 code uses AVX-512 under Rosetta 2 (ARM has no equivalent)?
