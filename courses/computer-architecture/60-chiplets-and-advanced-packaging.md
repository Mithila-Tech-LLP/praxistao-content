# Chapter 60: Chiplets and Advanced Packaging

For 50 years, the trend was clear: put more on one chip, make the chip bigger, and make every transistor smaller. This is changing. The new trend is **chiplets**: design multiple smaller dies (chiplets) and package them close together using advanced packaging techniques that approximate the performance of a monolithic die — at lower cost, higher yield, and with the ability to mix different manufacturing processes on different chiplets. AMD's EPYC processors, Apple's M1 Ultra, Intel's Ponte Vecchio GPU, and NVIDIA's HBM-equipped GPUs all use chiplet or advanced packaging techniques. This chapter explains what chiplets are, why they work, how they are packaged, and the key interconnect standards.

## Table of Contents

1. [Why Chiplets? The Economics of Die Size](#1-why-chiplets-the-economics-of-die-size)
2. [Die-to-Die Interconnects](#2-die-to-die-interconnects)
3. [AMD Chiplet Architecture — Zen + Infinity Fabric](#3-amd-chiplet-architecture--zen--infinity-fabric)
4. [Intel Advanced Packaging — EMIB, Foveros](#4-intel-advanced-packaging--emib-foveros)
5. [TSMC SoIC and CoWoS](#5-tsmc-soic-and-cowos)
6. [UCIe — Universal Chiplet Interconnect Express](#6-ucie--universal-chiplet-interconnect-express)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Why Chiplets? The Economics of Die Size

**Die size and yield**: As Chapter 55 explained, larger dies have exponentially worse yield. A 600mm² die at TSMC N5 (yield ~50–60%) costs dramatically more per working die than two 300mm² dies (yield ~80% each).

```
Monolithic vs chiplet yield math:
  
  Scenario: GPU with 80 billion transistors, N5 process (D=0.1 defects/cm²)
  
  Option A: Single 600mm² die:
    Y = e^(-D×A) = e^(-0.1 × 6) = e^(-0.6) = 0.549 = 54.9% yield
    Wafer: $17,000 / (117 dies/wafer × 0.549) = $265 per good die
    
  Option B: Four 150mm² chiplets:
    Y_each = e^(-0.1 × 1.5) = e^(-0.15) = 0.861 = 86.1% yield
    Each chiplet cost: $17,000 / (470 dies/wafer × 0.861) = $42 per chiplet
    Four chiplets: 4 × $42 = $168 for equivalent compute
    But: must assemble 4 chiplets (packaging cost ~$20) = $188 total
    
  Chiplets save: $265 - $188 = $77 per chip (29% savings)
```

**Process specialization**: Different chiplets can use different process nodes optimized for their function:
- Logic (CPU/GPU): leading edge node (3nm/5nm) for maximum performance
- SRAM (cache): a node optimized for SRAM density (often a generation behind)
- Analog (SerDes, PHY): mature node (28nm–65nm) — analog circuits don't scale well
- Memory (HBM): its own DRAM process

A monolithic chip must use ONE process for everything — compromising on all. Chiplets let you pick the best process per function.

### Quick Check
> 1. Why is chiplet yield better than monolithic die yield for the same total transistor count?
> 2. What is "process specialization" and why is it an advantage of chiplets?
> 3. At what die size does chiplet decomposition typically become economically worthwhile?

---

## 2. Die-to-Die Interconnects

The challenge with chiplets is connecting them. The connection between chiplets must approach the bandwidth and latency of on-chip connections — otherwise you lose the benefit of having fast inter-chiplet communication.

**Metrics for die-to-die interconnects:**
- **Bandwidth density**: GB/s per mm of connection width
- **Latency**: cycles/ns for a signal to cross the interface
- **Energy efficiency**: pJ per bit transferred
- **Pitch**: minimum distance between bumps/pads (smaller = more bumps = more bandwidth)

```
Comparison of interconnect options:

                Bandwidth density    Latency    Energy    Typical distance
PCIe 5 (off-pkg)  16 GB/s per lane   ~100ns     ~5 pJ/b   meters
NVLink 4 (on-PCB)  25 GB/s per lane   ~50ns      ~4 pJ/b   cm
Chip-to-chip (pkg) 1-2 TB/s total     ~1-5ns     ~1 pJ/b   mm
On-chip (SoC bus)  5-10 TB/s total    ~0.1ns     ~0.1 pJ/b µm

Goal: get die-to-die as close to "on-chip bus" as possible
```

**Packaging technologies determine die-to-die performance:**

**Organic substrate (flip-chip BGA)**: Traditional packaging. Dies placed on an organic substrate (PCB-like), connected via solder bumps. Pitch: 100–200µm. Low bandwidth, suitable for different chips that communicate via PCIe/NVLink.

**Si bridge (EMIB — Intel)**: A small silicon "bridge" chip embedded in the organic substrate, connecting two chiplets with fine-pitch connections (55µm bump pitch). Higher bandwidth than organic substrate for local chip-to-chip connections.

**Silicon interposer (CoWoS — TSMC)**: All chiplets are placed on a large passive silicon interposer wafer. The interposer has fine-pitch redistribution layers. Pitch: 10–40µm. High bandwidth (multiple TB/s) and low power. Used for HBM + GPU connections.

**3D stacking (Foveros — Intel, SoIC — TSMC)**: Stack one die on top of another, connected through fine-pitch Cu-Cu bonding. Pitch: 1–10µm. Extremely high bandwidth, ultra-short connections. Most expensive and complex.

### Quick Check
> 1. What is the key trade-off between chiplet packaging cost and die-to-die bandwidth?
> 2. What is a silicon interposer and how does it improve die-to-die bandwidth?
> 3. What is the difference between a silicon bridge (EMIB) and a full silicon interposer?

---

## 3. AMD Chiplet Architecture — Zen + Infinity Fabric

AMD pioneered the commercial adoption of chiplets with EPYC and Ryzen processors. The strategy:

**AMD EPYC (Rome, 3rd gen, 2019):**
- 8 × CPU chiplets (8 Zen 2 cores each = 64 cores total)
- 1 × I/O die (memory controllers, PCIe, Infinity Fabric hub)
- CPU chiplets: TSMC 7nm (compute-optimized)
- I/O die: GloFo/TSMC 12nm (I/O-optimized, cheaper for non-compute)

```
EPYC Rome package:
  
  ┌──────────────────────────────────────────────────┐
  │  CPU  │  CPU  │  CPU  │  CPU  │  CPU  │  CPU  │  │
  │  Die  │  Die  │  Die  │  Die  │  Die  │  Die  │  │
  │  8C   │  8C   │  8C   │  8C   │  8C   │  8C   │  │
  ├──────────────────────────────────────────────────┤
  │  CPU  │  CPU  │          I/O Die                 │  │
  │  Die  │  Die  │  (mem ctrl, PCIe, UPI, fabric)   │  │
  └──────────────────────────────────────────────────┘
  
  All chiplets connected via Infinity Fabric (Chapter 34)
  Each CPU die: one 8-core chiplet = one "CCD" (Core Complex Die)
```

**Why the I/O die is a separate chiplet:**
- Memory controllers (DDR4 PHY) use large, expensive circuitry that doesn't benefit from 7nm
- Running them at 7nm would waste expensive process node on non-scaling logic
- I/O die at 12nm: 1/3 the cost per mm² of 7nm

**Infinity Fabric bandwidth**: 16–36 GB/s per direction between CPU die and I/O die. Each CPU core has 8 channels × 32 GB/s = ~256 GB/s per die. With 8 CPU dies: 2 TB/s total internal bandwidth.

**AMD 3D V-Cache** (Zen 3, Zen 4): Stack an additional 64MB of SRAM directly on top of the CPU compute die using wafer-on-wafer bonding. The SRAM connects with 200 TB/s bandwidth (vs ~1 TB/s for HBM). Used for gaming CPUs where large cache dramatically improves performance.

### Quick Check
> 1. Why does AMD put the CPU cores on a different chiplet from the I/O die?
> 2. What manufacturing process does AMD use for CPU compute dies vs I/O die?
> 3. What is AMD 3D V-Cache and what problem does it solve?

---

## 4. Intel Advanced Packaging — EMIB, Foveros

Intel's strategy uses two complementary technologies:

**EMIB (Embedded Multi-die Interconnect Bridge):**
- A small silicon chip (~4mm²) embedded inside the organic substrate
- Chiplets on the package bridge their local connections through this embedded silicon
- EMIB has fine-pitch wiring (55µm bump pitch vs 100–200µm for organic substrate)
- Bandwidth: ~1 TB/s between adjacent chiplets
- Used in: Intel Stratix 10 FPGA, Ponte Vecchio GPU, Lunar Lake CPU

**Foveros (3D stacking):**
- Stack one die ("active interposer" or base tile) under another die
- Through-Silicon Vias (TSVs) or face-to-face bonding connects them
- Very high bandwidth (face-to-face: 10–20µm pitch, 100+ TB/s potential)
- Used in: Intel Lakefield (2020, first commercial Foveros product), Meteor Lake, Lunar Lake

```
Intel Lunar Lake (2024): Foveros example
  
  Top die (on-chip tile): CPU cores (Intel 3), GPU, NPU, media engine
  Bottom die (base tile): Intel 22nm — power delivery, I/O, SOC fabric
  
  Face-to-face Cu bonding: millions of connections, ~10µm pitch
  Benefit: move power delivery to a more cost-efficient node;
           use best node only for performance-critical logic
```

**Intel Ponte Vecchio (HPC GPU, 2022):**
- 47 chiplets in a single package!
- Uses both EMIB (horizontal) and Foveros (vertical stacking)
- 100 billion transistors total
- Chiplets from: Intel 7 (compute tiles), Intel 4 (Rambo cache), TSMC N7 (base tiles), Micron HBM2e
- Total package: 2,328 mm² — larger than any single silicon die

### Quick Check
> 1. What is EMIB and how does it differ from a full silicon interposer?
> 2. What is Foveros and what type of connection does it create?
> 3. How many chiplets does the Intel Ponte Vecchio GPU contain?

---

## 5. TSMC SoIC and CoWoS

TSMC provides two key advanced packaging platforms for its customers:

**CoWoS (Chip on Wafer on Substrate):**
- A large passive silicon interposer (as large as a full reticle field, ~800mm²)
- Multiple chiplets (HBM memory stacks + logic die) are placed on the interposer
- Interconnect pitch: 40µm (CoWoS-S) or 9µm (CoWoS-R with RDL interposer)
- Bandwidth: HBM to GPU via interposer: 3–5 TB/s (H100) or higher

```
NVIDIA H100 SXM package using CoWoS-S:
  
  Top view:
  ┌────────────────────────────────────────────────────────┐
  │                                                        │
  │  [HBM]  [HBM]     GH100 GPU Die      [HBM]  [HBM]   │
  │  stack  stack    (814mm², 80B trans)  stack  stack   │
  │                                                        │
  │   All connected through silicon interposer below       │
  └────────────────────────────────────────────────────────┘
  
  Silicon interposer: 1000+ mm², fine-pitch Cu redistribution layers
  GPU-to-HBM bandwidth: 3.35 TB/s (H100 80GB)
```

**SoIC (System on Integrated Chips):**
- TSMC's 3D stacking technology (face-to-face or face-to-back wafer bonding)
- Hybrid bonding: Cu-Cu direct bonding at 9µm pitch or sub-µm pitch
- Very high bandwidth: >1000 TB/s potential at tightest pitch
- Apple N3B uses SoIC-like stacking for M3 family's 3D SRAM layers

**The CoWoS supply constraint (2023):**
NVIDIA's explosive demand for H100/A100 GPUs stressed TSMC's CoWoS capacity. CoWoS requires large silicon interposers (difficult to manufacture at high yield) and specialized assembly. This limited H100 supply and drove up prices — a reminder that packaging can be as much a bottleneck as chip fabrication.

### Quick Check
> 1. What is CoWoS and what problem does it solve for GPU chips?
> 2. What is SoIC and how is it different from CoWoS?
> 3. Why was CoWoS capacity a bottleneck for H100 supply in 2023?

---

## 6. UCIe — Universal Chiplet Interconnect Express

The chiplet ecosystem has a problem: each company has its own proprietary interconnect (AMD Infinity Fabric, Intel EMIB/FLI, TSMC SoIC). A chiplet from Company A cannot plug into a package from Company B because the physical and electrical interfaces are different.

**UCIe (Universal Chiplet Interconnect Express)** aims to standardize chiplet interconnects — enabling a chiplet ecosystem where:
- Chip designers can mix chiplets from different vendors
- IP blocks can be reused across packages
- Smaller companies can design specialized chiplets sold to system integrators

**UCIe 1.0 (2022, consortium: Intel, AMD, ARM, Qualcomm, Samsung, TSMC, TSMC, Google, Meta, Microsoft):**
- Defines physical layer: bump pitch, pad design, connector pinout
- Defines link layer: packetization, flow control, error correction
- Standard packages: "Die-to-Die Adapter" at 25µm pitch and "Modular" at 100µm pitch
- Bandwidth: 1.3 TB/s per mm of interface width at 25µm pitch

```
UCIe chiplet ecosystem concept:
  
  CPU chiplet        Memory controller chiplet     I/O chiplet
  (TSMC 3nm)         (GlobalFoundries 12nm)         (TSMC 28nm)
      │                        │                        │
      └────────────────────────┴────────────────────────┘
                        UCIe interconnect
                    (standard physical interface)
  
  Any chiplet with UCIe interface can connect to any other
```

**Current state (2024)**: UCIe is still gaining adoption. Intel, AMD, and TSMC are all supporting it. Production chips using UCIe are beginning to appear but it is not yet as ubiquitous as PCIe.

**SiP (System-in-Package)** vs SoC: A SiP integrates multiple dies (possibly from different companies) in one package. The endpoint: a complete electronic system (CPU + memory + radio + sensors) in a package smaller than a SIM card.

### Quick Check
> 1. What problem does UCIe solve in the chiplet ecosystem?
> 2. What are the two physical configurations defined by UCIe 1.0?
> 3. What is a SiP and how is it different from an SoC?

---

## Summary

- **Chiplets** are smaller dies combined in one package. Better yield, process specialization, and modularity vs monolithic dies.
- **Die-to-die interconnects**: organic substrate (coarse), silicon bridge (EMIB, medium), silicon interposer (CoWoS, high bandwidth), 3D stacking (SoIC/Foveros, very high bandwidth).
- **AMD**: pioneered chiplets for x86. CPU compute dies on leading edge node (TSMC 7nm); I/O die on mature node. Infinity Fabric for die-to-die communication. 3D V-Cache adds SRAM on top.
- **Intel**: EMIB (in-package bridge for horizontal connection), Foveros (3D vertical stacking). 47-chiplet Ponte Vecchio GPU.
- **TSMC**: CoWoS (silicon interposer for HBM + GPU), SoIC (3D hybrid bonding for SRAM stacking).
- **UCIe**: industry standard for chiplet-to-chiplet interfaces; enables multi-vendor chiplet ecosystems.

---

## Exercises

### Easy
1. What is a chiplet and why does it improve manufacturing yield?
2. What is a silicon interposer and how does it improve die-to-die bandwidth?
3. What is AMD Infinity Fabric and what role does it play in chiplet-based EPYC processors?

### Medium
4. Yield economics: A GPU design requires 80B transistors at TSMC N5 (D=0.1 defects/cm²). Option A: one 500mm² monolithic die. Option B: four 125mm² chiplets assembled on an interposer. (a) Calculate yield for Option A using Y = e^(-D×A). (b) Yield for each chiplet in Option B. (c) Yield for getting all 4 chiplets working (yield^4). (d) Wafer cost: $17,000, 300mm wafer. Dies per wafer ≈ wafer area / die area. Cost per good die for each option (ignore packaging). (e) Assembly cost: $50 for chiplet package. Total cost per unit for B? Which is cheaper?
5. Die-to-die bandwidth requirements: An AMD EPYC CPU with 8 CPU chiplets (each 8-core Zen 4) and 1 I/O die. Each core needs peak memory bandwidth of 50 GB/s. (a) Total memory bandwidth needed for 64 cores? (b) If Infinity Fabric can provide 36 GB/s per direction per die, what is the total Infinity Fabric bandwidth from 8 CPU dies to I/O die? (c) Can the Infinity Fabric bottleneck memory bandwidth? (d) The I/O die supports DDR5 8-channel at 51.2 GB/s each = 409.6 GB/s total. Compare to the Infinity Fabric bandwidth.
6. HBM vs GDDR6 trade-offs: H100 GPU uses 6× HBM3 stacks, total 3.35 TB/s bandwidth. RTX 4090 uses GDDR6X, 21 Gbps × 384-bit bus = 1.008 TB/s. (a) Compare bandwidth and cost (HBM3 stack ~$500, GDDR6X board ~$100). (b) HBM requires CoWoS packaging at ~$300 extra; GDDR6X uses PCB traces at ~$10. Total memory subsystem cost ratio? (c) For AI training (needs maximum bandwidth): which is better? (d) For gaming (bandwidth less critical, cost matters): which is better? (e) At what price-per-bandwidth crossover does HBM make sense for non-AI?

### Hard
7. UCIe interface design: Design a UCIe-compliant interface between a CPU chiplet and an I/O chiplet. (a) UCIe 25µm standard bump pitch: if the interface width is 4mm, how many bumps are available? (b) UCIe uses differential signaling at 32 Gbps per lane. At 80% utilization (overhead): effective data bandwidth per lane? (c) Total bandwidth across the 4mm interface? (d) A 64-core CPU needs 512 GB/s to the I/O chiplet (for memory + PCIe + I/O). Can the 4mm UCIe interface satisfy this? (e) Power: each UCIe bump at 32 Gbps consumes ~1 pJ/bit. What is the power for the 4mm interface at peak bandwidth?
8. Chiplet security: Chiplets from multiple vendors in one package create a new security threat surface. (a) If the CPU chiplet and the security enclave chiplet are different dies from different companies: what new attack surfaces exist that don't exist in a monolithic SoC? (b) How can UCIe die-to-die traffic be encrypted to prevent snooping at the interposer? (c) If an adversary can modify the interposer (supply chain attack): what are the implications for secure enclaves that store private keys? (d) How do current secure chiplet designs (like ARM TrustZone over UCIe) handle this threat model?
