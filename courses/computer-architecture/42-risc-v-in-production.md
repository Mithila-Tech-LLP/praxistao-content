# Chapter 42: RISC-V in Production

RISC-V (pronounced "risk five") is a free, open-source instruction set architecture. Unlike x86 (Intel's proprietary ISA) or ARM (licensed from ARM Holdings for fees), RISC-V belongs to nobody — and therefore everybody. Any company, university, or individual can implement RISC-V without paying royalties or signing license agreements. This radical openness is disrupting the chip industry. RISC-V chips are now shipping in billions of microcontrollers, running inside SSD controllers, graphics cards, and server management chips. China is betting on RISC-V to escape dependence on ARM and x86. And a growing ecosystem of open-source CPU designs means that for the first time, hardware can be as open as Linux.

## Table of Contents

1. [Why RISC-V Was Created](#1-why-risc-v-was-created)
2. [RISC-V Architecture Overview](#2-risc-v-architecture-overview)
3. [RISC-V in Commercial Chips](#3-risc-v-in-commercial-chips)
4. [SiFive — The ARM of RISC-V](#4-sifive--the-arm-of-risc-v)
5. [RISC-V in Western Digital HDDs/SSDs](#5-western-digital--all-in-on-risc-v)
6. [Alibaba T-Head — The Open-Source Titan](#6-alibaba-t-head--the-open-source-titan)
7. [RISC-V in the GPU: NVIDIA and AMD](#7-risc-v-in-gpus-nvidia-and-amd)
8. [RISC-V for India and the SHAKTI Connection](#8-risc-v-for-india-and-the-shakti-connection)
9. [Challenges for RISC-V](#9-challenges-for-risc-v)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why RISC-V Was Created

In 2010, Krste Asanović and David Patterson at UC Berkeley were designing a new chip for their ParLab research project. They evaluated existing ISAs:

- **x86**: Proprietary, legacy-laden, complex, Intel/AMD would never open-source it
- **MIPS**: Licensed, being sold to MIPS Technologies, uncertain future
- **ARM**: Licensed (royalties), ARM Holdings controls the ISA evolution
- **SPARC**: Sun's ISA, Oracle now controls it, not ideal
- **PowerPC**: IBM's ISA, complicated history

None worked for open academic research. They designed a new ISA from scratch: **RISC-V**. The Roman numeral V (5) indicates it was the fifth RISC design from Berkeley (RISC-I through RISC-IV came before).

**Key RISC-V design principles:**
1. **Free and open**: No royalties, no license agreements, published specification
2. **Clean, simple, extensible**: Base ISA is small (47 instructions for RV32I); optional extensions add functionality
3. **Modern**: Designed knowing lessons from decades of ISA mistakes
4. **Not backward-compatible with anything**: No legacy baggage

**The RISC-V Foundation (now RISC-V International)**: Founded 2015, incorporated in Switzerland (politically neutral), governs the ISA specification. Over 3,000 member organizations. Open governance — no single company controls RISC-V's future.

### Quick Check
> 1. Why did Berkeley researchers decide to create a new ISA instead of using an existing one?
> 2. What does the "V" in RISC-V stand for?
> 3. What is RISC-V International and why is it incorporated in Switzerland?

---

## 2. RISC-V Architecture Overview

RISC-V is a **RISC, load-store architecture** (like ARM) with clean 32-bit fixed-length instructions (for the base ISA). It is designed to be modular:

**Base ISAs:**
- **RV32I**: 32-bit base integer ISA (32 registers, 32-bit addresses) — the minimum
- **RV64I**: 64-bit base integer ISA — for 64-bit systems
- **RV128I**: 128-bit (experimental, not widely implemented)

**Standard Extensions (designated by letters):**

| Extension | Adds |
|-----------|------|
| M | Integer Multiply/Divide |
| A | Atomic instructions (for multiprocessor synchronization) |
| F | Single-precision Floating-Point |
| D | Double-precision Floating-Point |
| C | Compressed instructions (16-bit encoding for code density) |
| G | Shorthand for IMAFD (all of the above) |
| V | Vector operations |
| H | Hypervisor support |
| Zicsr | Control and Status Register instructions |
| Zifencei | Instruction-Fetch Fence |

A "general-purpose" embedded RISC-V chip typically implements **RV32IMAC** (base + multiply + atomic + compressed). A server-class chip implements **RV64GC**.

**RISC-V register file (RV64I):**
- 32 integer registers: x0–x31 (x0 is hardwired to zero)
- Calling convention assigns roles: a0–a7 (arguments), t0–t6 (temps), s0–s11 (saved), ra (return address), sp (stack pointer)

**RISC-V vs ARM comparison:**
- Both are RISC, load-store
- ARM has more complex addressing modes (stronger base)
- RISC-V is simpler and modular — easier to implement and verify
- RISC-V has no PC-register (unlike ARM), no condition codes (unlike ARM A32), no predicated execution
- RISC-V is free; ARM charges license fees

### Quick Check
> 1. What does "RV64GC" mean in RISC-V nomenclature?
> 2. How many registers does RISC-V have? What is special about x0?
> 3. What is the advantage of RISC-V's modular extension approach for embedded chip designers?

---

## 3. RISC-V in Commercial Chips

RISC-V has moved from academic curiosity to commercial reality faster than most expected. Key production deployments:

**Consumer devices:**
- Every **Android phone** with a Qualcomm Snapdragon 7/8 Gen 1+ or MediaTek Dimensity 9000+: RISC-V core inside for power management
- **Espressif ESP32-C3** (2020): RISC-V microcontroller for Wi-Fi/Bluetooth IoT. Billions shipped.
- **Espressif ESP32-P4** (2024): Dual-core RISC-V at 400 MHz, targeting edge AI
- **GigaDevice GD32VF103**: RISC-V microcontroller for industrial
- **StarFive VisionFive 2**: Single-board computer with RISC-V SoC (JH7110), runs Linux

**Semiconductor IP:**
- **Western Digital**: Replaced MIPS cores in HDD/SSD controllers with RISC-V (Chapter 5 below)
- **Google**: Uses RISC-V in Pixel phone security chip (Titan M2)
- **NVIDIA**: RISC-V cores inside NVIDIA GPUs for power management/firmware
- **AMD**: Uses RISC-V in some internal management functions
- **Intel**: Has RISC-V IP via SiFive acquisition attempt (blocked), owns RISC-V cores via other means

**China-specific:**
- U.S. export controls on ARM and x86 to certain Chinese entities push China toward RISC-V
- Alibaba's T-Head division (covered below)
- Multiple Chinese startups: Nuclei System Technology, Andes Technology

### Quick Check
> 1. Name three consumer products that contain a RISC-V processor.
> 2. Why is RISC-V strategically important for China?
> 3. How are NVIDIA and AMD using RISC-V, even though their main processors are not RISC-V?

---

## 4. SiFive — The ARM of RISC-V

**SiFive** (founded 2015) is the most prominent RISC-V commercial IP company — often called "the ARM of RISC-V."

SiFive's founders (Krste Asanović and Yunsup Lee) are the original RISC-V inventors. SiFive licenses RISC-V CPU IP — similar to ARM's business model but:
- Royalty-free (RISC-V base ISA), fees for SiFive's implementation IP
- More flexible licensing terms than ARM
- Open-source starter cores (Freedom E310, used in HiFive boards)

**SiFive's IP portfolio:**
- **U-series** (application processors): U54/U54-MC, U74, P670 — 64-bit Linux-capable, in-order, OOO
- **E-series** (embedded): E24, E31, E34 — 32-bit, tiny, for microcontrollers
- **S-series** (server): S76 — high-performance server core
- **X280**: High-performance with RISC-V V extension (vector)

**SiFive vs ARM for IP licensing:**
- SiFive cores are typically 20-30% cheaper to license than equivalent ARM Cortex designs
- RISC-V ISA is free — chips can evolve the ISA without paying ARM for new extensions
- But: ARM's ecosystem (software, tools, OS support) is massively more mature

**The Intel acquisition attempt**: In 2021, Intel offered $2 billion to acquire SiFive — a sign of how seriously Intel views RISC-V as a strategic threat (or opportunity). The deal fell through after SiFive raised VC funding instead. Intel has since invested in RISC-V through other means.

### Quick Check
> 1. What is SiFive and how does its business model compare to ARM Holdings?
> 2. What is the SiFive U-series used for vs the E-series?
> 3. Why did Intel try to acquire SiFive?

---

## 5. Western Digital — All-In on RISC-V

**Western Digital** made a dramatic RISC-V bet in 2017: replace all the MIPS and ARM processors inside their HDD and SSD controllers with RISC-V. They ship ~1 billion RISC-V cores per year.

**Why WD switched:**
1. WD needed custom ISA extensions for storage-specific operations (cryptography, XOR-intensive LDPC error correction, command processing)
2. MIPS required paying license fees
3. ARM licenses came with restrictions on custom extensions
4. RISC-V allowed WD to add custom instructions freely, change the design, and not pay royalties

**SweRV EH1/EH2 cores**: WD designed the SweRV — a high-performance, dual-issue in-order RISC-V core optimized for storage controller workloads. They open-sourced the SweRV design via the CHIPS Alliance.

```
SweRV EH1 (open-sourced by Western Digital):
  RV32IMC
  Dual-issue in-order pipeline
  9-stage pipeline
  Designed at 28nm for storage workloads
  ~4 DMIPS/MHz
```

This open-sourcing was significant: a production-grade, commercially deployed RISC-V core released for anyone to use. Universities and startups now use SweRV as a starting point.

### Quick Check
> 1. What business reason drove Western Digital to switch from MIPS/ARM to RISC-V?
> 2. What is the SweRV and why is it significant that WD open-sourced it?
> 3. How many RISC-V cores per year does WD ship?

---

## 6. Alibaba T-Head — The Open-Source Titan

Alibaba's chip design division, **T-Head Semiconductor** (founded 2018), has created some of the highest-performance RISC-V chips in the world.

**XuanTie (玄铁) C910 (2019)**:
- 64-bit RISC-V
- 3-issue OOO execution
- 12-stage pipeline
- Supports Linux, Zephyr, FreeRTOS
- Commercially used in smart cameras, industrial devices

**XuanTie C920 (2022)**:
- Improved OOO
- RISC-V Vector (RVV) extension support
- Used in Alibaba's Yitian 710 server chip (for AliCloud infrastructure)

**Yitian 710 (2021)**:
- 128-core RISC-V + T-Head custom ISA extensions
- TSMC 5nm
- Used in Alibaba Cloud's Yitian ECS instances
- Claimed: competitive with ARM Neoverse N1 cores

**XuanTie open-sourcing**: Alibaba open-sourced several T-Head RISC-V cores on GitHub:
- T-Head XuanTie E902, E906, E907, C906, C910
- Available to Chinese and international developers

**China's RISC-V strategy**: U.S. export controls restrict Chinese companies from using ARM (ARM Holdings has dual US-UK compliance requirements) and from accessing advanced manufacturing. RISC-V, owned by nobody, bypasses this — any Chinese company can implement RISC-V without needing permission from ARM or Intel.

```
Chinese RISC-V ecosystem (2024):
  Hardware: T-Head/Alibaba, SpacemiT, Andes Technology, Nuclei System
  Tools: Open-source GCC/LLVM toolchain
  OS: Linux (full support), RTOS (FreeRTOS, RT-Thread)
  Support: China's MIIT invested ¥40B ($5.5B) in RISC-V ecosystem development
```

### Quick Check
> 1. What is the Yitian 710 and where is it deployed?
> 2. Why is RISC-V strategically important for China's semiconductor independence?
> 3. What is T-Head and who owns it?

---

## 7. RISC-V in GPUs: NVIDIA and AMD

An interesting place to find RISC-V: inside NVIDIA and AMD GPUs — not as the main compute engine, but as internal management cores.

**NVIDIA Falcon → RISC-V**: NVIDIA's GPUs contain many embedded processor cores that handle firmware, power management, display control, and security. The original Falcon ISA (NVIDIA's proprietary) is being migrated to RISC-V. NVIDIA uses RISC-V in:
- **GSP (GPU System Processor)**: manages power states, PCIe communication
- **FECS/GPCCS**: graphical context switch controllers in recent GPU architectures
- Security microcontrollers

The RISC-V adoption simplifies NVIDIA's internal tool chain and allows using the standard RISC-V GCC/LLVM toolchain instead of proprietary tools.

**AMD**: Uses RISC-V in power management controllers in some recent GPU designs (similar reasoning — open toolchain, no licensing).

This is a pattern across the industry: even companies that use ARM or x86 for main application processors often use RISC-V for small embedded management tasks, where RISC-V's simplicity, toolchain maturity, and zero cost are compelling.

### Quick Check
> 1. What does NVIDIA use RISC-V for inside its GPUs?
> 2. Why would NVIDIA use RISC-V for internal microcontrollers instead of its own Falcon ISA?
> 3. What does this use pattern tell us about where RISC-V is most competitive today?

---

## 8. RISC-V for India and the SHAKTI Connection

India has strategic ambitions to build domestic semiconductor capability. The **SHAKTI processor** from IIT Madras is India's most prominent RISC-V chip — and it will be covered in detail in Chapters 62–68. But the RISC-V connection is worth noting here:

SHAKTI chose RISC-V specifically because:
1. **No royalties**: Critical for a publicly-funded academic project
2. **Freedom to customize**: India's security and defense applications need custom extensions without ARM's approval
3. **Open ecosystem**: Tools (RISC-V GCC, LLVM, spike simulator) are free
4. **Independence**: Not dependent on any foreign company's licensing decisions

India's semiconductor mission (ISM, 2021) specifically supports RISC-V development as a way to build indigenous IP. The SHAKTI program has taped out 7 generations of RISC-V chips, with the latest fabricated at Intel and GlobalFoundries. Full details in Chapter 62+.

### Quick Check
> 1. Why did IIT Madras choose RISC-V for the SHAKTI processor?
> 2. What is India's semiconductor strategy related to RISC-V?
> 3. How many SHAKTI chip generations have been fabricated?

---

## 9. Challenges for RISC-V

RISC-V is growing fast but faces real challenges:

**Software ecosystem maturity**:
- Linux: excellent support (RISC-V is a first-class Linux architecture since 5.10)
- Android: official support added but fewer devices → less testing
- Windows: Microsoft ported Windows 11 ARM version to RISC-V (experimental, 2023)
- Proprietary software: Adobe, Microsoft Office, many commercial apps not yet RISC-V native

**Fragmentation risk**:
- The extension system allows infinite customization — but each custom extension breaks binary compatibility
- A binary compiled for RV64GCV may not run on RV64GC (no vector)
- ARM solved this with strict compatibility requirements; RISC-V's freedom creates fragmentation

**ISA stability**:
- RISC-V ratification process can be slow
- Some extensions (Vector, Hypervisor) took years to ratify
- Software can't rely on an extension until it's ratified and deployed at scale

**No single "flagship" high-performance core**:
- ARM has Cortex-X4 (4 GHz, world-class IPC); RISC-V's best open-source cores (SiFive P670) are 2–3 generations behind
- Commercial RISC-V high-performance cores exist but lack Apple Firestorm-level investment

**Toolchain maturity**:
- RISC-V GCC/LLVM are production quality
- But: profiling tools, OS/hypervisor support, debug tooling (JTAG, OpenOCD) less polished than ARM equivalents

Despite these challenges, RISC-V's trajectory is strongly upward — particularly in embedded/IoT and China, where the cost and freedom advantages outweigh maturity gaps.

### Quick Check
> 1. What is the "fragmentation risk" with RISC-V's extension system?
> 2. What major operating systems officially support RISC-V?
> 3. What is the gap between the best RISC-V application processors and ARM Cortex-X/Apple Firestorm?

---

## Summary

- **RISC-V** is a free, open-source ISA born at UC Berkeley (2010), governed by RISC-V International. No royalties, no license restrictions, free to customize.
- Architecture: modular base ISA (RV32I/RV64I) + optional standard extensions (M, A, F, D, C, V, H, ...). A typical embedded config: RV32IMAC.
- **Production deployments**: WD storage controllers (1B+ cores/year), Espressif ESP32 (IoT), NVIDIA/AMD management cores, Google Titan M2, Chinese commercial chips.
- **SiFive**: the leading commercial RISC-V IP vendor, licensing cores like ARM does but cheaper and with RISC-V's freedom.
- **Western Digital** open-sourced SweRV — production RISC-V core for storage.
- **Alibaba T-Head** builds high-performance RISC-V for AliCloud (Yitian 710, 128-core server chip).
- **China** is investing heavily in RISC-V as a path to semiconductor independence from ARM and x86.
- **SHAKTI (India)**: open-source RISC-V processor family from IIT Madras — detailed in Chapters 62–68.
- Challenges: software ecosystem, fragmentation risk, no flagship high-performance core yet.

---

## Exercises

### Easy
1. What does RISC-V stand for, and why was it created?
2. What does "RV64GC" mean? Expand each component.
3. Why does Western Digital ship ~1 billion RISC-V cores per year?

### Medium
4. Compare RISC-V vs ARM for a startup designing an IoT chip: (a) licensing cost (assume ARM Cortex-M4 license = $2M, plus $0.10/chip royalty; RISC-V = $0 for ISA + $500K for SiFive E34 IP), (b) time to market (ARM: proven ecosystem; RISC-V: may need more software work), (c) flexibility (ARM: limited custom instruction support; RISC-V: add any extension), (d) at what production volume does RISC-V's zero royalty outweigh the higher NRE?
5. RISC-V's modular extension system: A chip implements RV64IMACFV (base + multiply + atomic + compressed + float + vector). (a) What workloads does the V extension enable? (b) A program compiled for RV64IMAC runs on this chip — does the vector unit ever get used? (c) Can the chip run RISC-V Linux without the F extension? (d) If a manufacturer adds a custom "crypto" extension (custom instruction format), can that chip run standard RISC-V binaries?
6. The "fragmentation problem": ARM's AArch64 has a stable ISA that guarantees any AArch64 binary runs on any AArch64 chip (within capability sets). RISC-V's custom extensions break this. Describe how the RISC-V hardware multi-extension detection mechanism (RISC-V ISA extension discovery via `misa` CSR and device trees) partially solves this and what its limitations are.

### Hard
7. Alibaba's Yitian 710 (128-core RISC-V, 5nm) vs ARM Neoverse N2 (128-core hypothetical): Analyze the competitive position: (a) What ISA advantage does Yitian 710 have over ARM for Alibaba's workloads? (b) What ecosystem disadvantage does Yitian 710 have? (c) Why might Alibaba build their own chip instead of buying ARM-based Ampere Altra servers? (d) Model the long-term strategic value: if Alibaba ships 1M servers/year and saves $100/server in license fees, what is the 5-year NPV of the RISC-V investment?
8. Design a minimal RISC-V SoC for a smart thermometer: reads temperature sensor (I2C), displays on small LCD, connects to Wi-Fi for cloud reporting, runs on battery for 2 years. Specify: (a) which RISC-V extensions are needed (justify each), (b) memory requirements (flash for program, SRAM for runtime), (c) power management strategy (what ISA features help with sleep modes), (d) compare to using an ARM Cortex-M33 — what would you gain/lose with RISC-V vs ARM for this specific application?
