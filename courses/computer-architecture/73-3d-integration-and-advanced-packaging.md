# Chapter 73: 3D Integration and Advanced Packaging

In the previous chapter we saw that 3D integration is one of the four key post-Moore strategies. It deserves its own chapter because it is the most commercially impactful right now — already present in every NVIDIA GPU, every AMD EPYC server processor, and every Apple M-series chip. Advanced packaging has become the key competitive differentiator in the semiconductor industry: TSMC's CoWoS and SoIC capabilities are as important to winning customers as their transistor density. This chapter explains the physics, the technologies, and the business strategy behind 3D integration.

## Table of Contents

1. [Why 3D Integration — The Limits of 2D](#1-why-3d-integration--the-limits-of-2d)
2. [Packaging Hierarchy: Wire Bond → Flip Chip → Advanced](#2-packaging-hierarchy-wire-bond--flip-chip--advanced)
3. [Intel's EMIB and Foveros](#3-intels-emib-and-foveros)
4. [TSMC's CoWoS and SoIC](#4-tsmcs-cowos-and-soic)
5. [AMD's 3D V-Cache and EPYC Architecture](#5-amds-3d-v-cache-and-epyc-architecture)
6. [UCIe: The Open Chiplet Standard](#6-ucie-the-open-chiplet-standard)
7. [Thermal Challenges in 3D](#7-thermal-challenges-in-3d)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why 3D Integration — The Limits of 2D

**Die yield economics (review from Chapter 60):**
- A single die area of 800mm² at TSMC N5 with defect density D₀=0.15/cm²:
  - Yield ≈ e^(-D₀ × A) = e^(-0.15 × 8) ≈ 30%
  - You discard 70% of your silicon
- Four chiplets of 200mm² each:
  - Yield per chiplet ≈ e^(-0.15 × 2) ≈ 74%
  - Combined yield ≈ 74%⁴ ≈ 30% — same? NO: you test before assembly
  - After binning out bad chiplets: ~74%⁴ / assembly_yield ≈ 55–65%
  - Yield advantage: ~2× fewer wasted chiplets

**Process optimization (chiplets enable heterogeneity):**
```
Monolithic SoC problem:
  CPU logic → needs N3 for performance
  SRAM → N3 has bad density vs N7; uses huge area
  Analog PHY → N3 analog is expensive; better at N28
  DRAM → DRAM has its own process; can't mix with logic
  
  Compromise: everything at N5 (mediocre for everything)
  
Chiplet solution:
  CPU logic chiplets → TSMC N3 (optimal for logic)
  Cache SRAM dies → TSMC N7 (better SRAM density)
  Analog I/O die → N28 (cheap, mature analog)
  HBM memory stacks → DRAM-specific process
  
  Each component at its optimal process → much better overall
```

**3D integration bandwidth**: interconnecting chiplets at sub-10µm bump pitch gives 10–100× more wires than PCIe, approaching on-chip wire density. Chapter 60 covered this; we focus on manufacturing details here.

### Quick Check
> 1. Why does chiplet architecture improve yield economics compared to monolithic SoC?
> 2. Why can't you mix logic, SRAM, analog, and DRAM on one process node optimally?
> 3. What is the key bandwidth advantage of chiplet-to-chiplet interconnects vs PCIe?

---

## 2. Packaging Hierarchy: Wire Bond → Flip Chip → Advanced

**Wire bond (1960s–present):**
```
  ┌──────────────────────────────────┐  ← package lid
  │   ┌────────────────┐             │
  │   │  silicon die   │             │
  │   └───────┬────────┘             │
  │            │ gold wires (bond)   │
  │  ┌─────────┴──────────────────┐  │
  │  │      substrate / leadframe │  │
  └──┴────────────────────────────┴──┘
                 ↕
             PCB traces

Bond wire length: 1–5mm → inductance, resistance, signal degradation
I/O density: ~100–500 bond pads per die
Bandwidth: limited (MHz range)
Cost: cheap, used for simple chips today
```

**Flip chip (1990s–present, mainstream in high-performance):**
```
  ┌────────────────────────────────┐  ← package
  │      ┌──────────────┐         │
  │      │  silicon die │ ← upside-down
  │      │ (flip chip)  │
  │      └──────┬───────┘
  │             │ controlled-collapse solder bumps (C4)
  │   ┌─────────┴────────────────┐
  │   │  organic substrate (BGA) │
  │   └─────────────┬────────────┘
  └─────────────────┼────────────┘
                    │ BGA solder balls
                   PCB
                   
Bump pitch: 100–150µm (C4 bumps)
I/O density: thousands of bumps
Bandwidth: 100s of GB/s
Cost: moderate; used in every modern CPU and GPU
```

**Advanced packaging (2012–present):**
Fine-pitch bumps, silicon interposers, TSVs, hybrid bonding — the focus of this chapter.

**Key metric progression:**
```
Technology          Bump pitch    I/O density    Bandwidth/mm
──────────────────────────────────────────────────────────
Wire bond           N/A           <1,000/die     ~50 GB/s
Flip chip C4        150µm         5,000–50,000   ~1 TB/s
µC4/micro-bump      40–55µm       100,000+       ~5 TB/s
EMIB silicon bridge 10–55µm       varies         ~10 TB/s/mm
CoWoS interposer    40µm          millions       ~10 TB/s/mm
SoIC hybrid bond    <10µm         billions       ~100 TB/s/mm²
```

### Quick Check
> 1. What is the difference between wire bond and flip chip packaging?
> 2. What is "bump pitch" and why does a smaller bump pitch enable more bandwidth?
> 3. What is the bandwidth advantage of hybrid bonding (SoIC) over C4 flip chip?

---

## 3. Intel's EMIB and Foveros

**EMIB (Embedded Multi-die Interconnect Bridge, 2017):**
Intel's approach to connecting chiplets at high bandwidth without an expensive full silicon interposer.

```
EMIB cross-section:
  
  ┌──────────────────┐     ┌──────────────────┐
  │   chiplet A      │     │   chiplet B       │
  │   (CPU die)      │     │   (HBM/IO die)   │
  └─────────┬────────┘     └────────┬──────────┘
            │ µ-bumps               │ µ-bumps
  ──────────┤              ────────┤
  organic   │    ┌──────────┐      │   organic
  substrate │    │ EMIB     │      │   substrate
  ──────────┴────┤ silicon  ├──────┴──────────────
                 │ bridge   │
                 │ (small   │
                 │ chip)    │
                 └──────────┘
                 
EMIB key facts:
- A small silicon chip embedded IN the organic substrate (not on top)
- Bridge area: ~15–25mm²
- Bump pitch on bridge: 55µm (vs 130µm on organic substrate)
- Bandwidth: ~2 TB/s across the bridge
- Cost: much cheaper than full silicon interposer
- Intel Stratix 10 FPGA (2016): first EMIB product (HBM + FPGA)
- Intel Ponte Vecchio (2022): 47 tiles connected with EMIB + Foveros
```

**Foveros (2019, first 3D logic-on-logic stacking):**
Unlike EMIB (which connects chiplets side-by-side in 2D), Foveros stacks dies vertically:

```
Foveros cross-section:
  
  ┌────────────────────────┐
  │  Top die (chiplet)     │  ← processor cores, GPU, etc.
  │  (active silicon)      │
  └──────────┬─────────────┘
             │ µ-bumps (36µm pitch)
  ┌──────────┴─────────────┐
  │  Base die (active)     │  ← I/O, power delivery, PHY
  │  (logic + drivers)     │
  └──────────┬─────────────┘
             │ C4 bumps (normal)
         package substrate
         
Intel Meteor Lake (Core Ultra, 2023): first consumer product with Foveros
- Top die: CPU compute tiles (Intel 4 process)
- Base die: I/O tile (Intel 22FFL process)
- Benefit: CPU tiles fabricated at cutting-edge process; I/O at mature process
```

**Foveros Direct (2022): hybrid bonding**
- Bump pitch: 10µm (vs 36µm for regular Foveros)
- Direct copper-to-copper bonding — no solder
- 36× more I/O density than regular Foveros
- Bandwidth: 200 TB/s/mm² (estimated)

### Quick Check
> 1. What is EMIB and how does it differ from a full silicon interposer?
> 2. What is Foveros and what consumer product first used it?
> 3. What is Foveros Direct and what improvement does it offer over regular Foveros?

---

## 4. TSMC's CoWoS and SoIC

TSMC has the world's most advanced packaging capabilities, which are a key reason why NVIDIA, AMD, and Apple use TSMC as their foundry.

**CoWoS (Chip on Wafer on Substrate, 2012):**
```
CoWoS-S (silicon interposer variant) cross-section:
  
  ┌──────────┐  ┌──────────┐  ┌───────────┐
  │  GPU die │  │  HBM #1  │  │  HBM #2   │  ← chiplets
  └────┬─────┘  └─────┬────┘  └─────┬─────┘
       │ µ-bumps      │             │
  ─────┴──────────────┴─────────────┴──── ← silicon interposer
  ┌────────────────────────────────────┐
  │     Silicon interposer             │
  │     (passive, just wires)          │
  │     2.5D integration               │
  └───────────────┬────────────────────┘
                  │ C4 bumps
              substrate / PCB
              
NVIDIA H100 uses CoWoS-S:
- Interposer size: ~850mm² (larger than the GPU die itself)
- HBM3 stacks: 5 × (819 GB/s) = 3.35 TB/s total bandwidth
- Interposer made at TSMC N40 (40nm, doesn't need advanced node)
- Bump pitch between GPU and HBM: 25µm (on silicon interposer)
- This is why H100 is so fast — 3.35 TB/s is only possible with CoWoS
```

**CoWoS variants:**
- **CoWoS-S**: Silicon interposer (passive redistribution layer) — highest bandwidth
- **CoWoS-R**: RDL (Redistribution Layer) interposer using organic material — cheaper, for less bandwidth-critical applications
- **CoWoS-L**: Local silicon interconnect (embedded silicon bridges, like EMIB) — intermediate

**SoIC (System on Integrated Chip, 2022):**
TSMC's 3D face-to-face stacking using hybrid bonding:

```
SoIC cross-section:
  
  ┌────────────────────────┐
  │  Top die               │  ← memory or logic
  │  (face down)           │
  └──────────┬─────────────┘
             │ copper-copper hybrid bonding
             │ pitch: <10µm → <3µm being developed
  ┌──────────┴─────────────┐
  │  Bottom die            │  ← logic die
  │  (face up)             │
  └────────────────────────┘
  
SoIC-X: 3D stacking (dies bonded through TSVs, vertical current flow)
SoIC-F: face-to-face (copper pads on top die bond to pads on bottom die)

Bandwidth: 100–1,000× more than C4 flip chip
AMD 3D V-Cache uses SoIC-F:
- 64MB L3 SRAM die stacked face-to-face on top of CPU die
- 9µm bump pitch (initial), going to <3µm
- 2 TB/s bandwidth between cache die and CPU die
```

**TSMC advanced packaging roadmap:**
- 2024: CoWoS capacity expansion (for NVIDIA Blackwell, H200 demand)
- 2025: CoWoS-L with embedded bridges
- 2026: SoIC-X with through-silicon vias for 3D logic stacking
- 2028: system-level 3D integration (multiple SoIC stacks on one CoWoS interposer)

### Quick Check
> 1. What is a silicon interposer in CoWoS and what is its role?
> 2. Why does NVIDIA H100 use CoWoS and what bandwidth does it achieve?
> 3. What is the difference between CoWoS (2.5D) and SoIC (3D)?

---

## 5. AMD's 3D V-Cache and EPYC Architecture

AMD's chiplet strategy (Chapter 60) uses advanced packaging for both EPYC (server) and Ryzen (desktop) processors. 3D V-Cache extends this to on-chip SRAM stacking.

**AMD 3D V-Cache (Ryzen 7 5800X3D, November 2022):**
- Problem: gaming performance often limited by L3 cache miss rate, not CPU compute
- Solution: stack 64MB extra SRAM directly on top of CPU die using TSMC SoIC-F
- Result: 96MB total L3 (32MB on CPU die + 64MB stacked) for one CPU chiplet
- Performance: 10–25% gaming improvement vs non-3D model

```
AMD Ryzen 9 7950X3D die structure:
  
  [CCD die (Zen 4 cores)]  [CCD die (regular, no V-Cache)]
  ┌─────────────────────┐  ┌──────────────────────────────┐
  │ 64MB SRAM cache die │  │  8 Zen 4 cores               │
  │ (stacked on top)    │  │  32MB L3 cache               │
  ├─────────────────────┤  └──────────────────────────────┘
  │ 8 Zen 4 cores       │
  │ 32MB L3 cache       │
  └─────────────────────┘  ← CPU chiplets at TSMC N5
                ↕ Infinity Fabric                
  ┌───────────────────────────────────────────────┐
  │  IOD: I/O die (TSMC N6)                       │
  │  PCIe 5.0, USB 4, DDR5 controller, etc.       │
  └───────────────────────────────────────────────┘
```

**EPYC Genoa (AMD, 2022):**
- Up to 12 CPU chiplets (CCD) per package
- Each CCD: 8 Zen 4 cores, TSMC N5
- 1 IOD: I/O functions, TSMC N6
- Total: 96 cores per socket (12 × 8)
- Compared to Intel Xeon 8592+: 60 cores on monolithic tile
- AMD's chiplet advantage: density, yield, power efficiency

**EPYC with 3D V-Cache (Genoa-X):**
- Select CCDs have 64MB extra SRAM stacked
- Total L3: up to 1.1GB (12 CCDs × 96MB)
- Targets: HPC simulations, EDA workloads, databases — all highly cache-sensitive
- 74% performance improvement in some FEM simulation workloads

### Quick Check
> 1. What is AMD 3D V-Cache and what problem does it solve?
> 2. How many chiplets does an AMD EPYC Genoa have and what are their roles?
> 3. Why does stacking SRAM on the CPU die improve gaming performance?

---

## 6. UCIe: The Open Chiplet Standard

**UCIe (Universal Chiplet Interconnect Express)** is an open standard for chiplet-to-chiplet interconnects — the "PCIe of chiplets."

**Why UCIe matters:**
- Currently, chiplet interconnects are proprietary: AMD Infinity Fabric, Intel EMIB+AIB, Apple proprietary
- A company designing a custom chiplet can't plug it into another company's SoC without proprietary agreement
- UCIe aims to create an industry standard so chiplets from different vendors can interconnect

**UCIe 1.0 (March 2022, founding members: Intel, AMD, ARM, Google, Meta, Microsoft, TSMC, Samsung, Qualcomm):**

```
UCIe protocol stack:
  
  ┌────────────────────────────────────┐
  │  Protocol layer: PCIe 6.0 or CXL  │  ← existing protocols
  ├────────────────────────────────────┤
  │  Die-to-die adapter layer          │  ← flow control, framing
  ├────────────────────────────────────┤
  │  Physical layer                    │  ← bump layout, timing
  └────────────────────────────────────┘
  
UCIe 1.0 specs:
  Standard package (55µm bump pitch): 16 GB/s/mm, 0.5 pJ/bit
  Advanced package (25µm bump pitch): 96 GB/s/mm, 0.25 pJ/bit
  Ultra-fine (10µm, UCIe 2.0 target): 1,300 GB/s/mm = 1.3 TB/s/mm
  
Compare to PCIe 5.0 (board-level): 16 GB/s per lane direction (much less dense)
```

**UCIe use cases:**
- CPU companies buy GPU chiplets from different vendors
- AI accelerator chiplets from startups plug into standard CPU platforms
- Memory chiplets (HBM-like) with standardized interface
- "Chiplet marketplace": buy-design-integrate independently

**Status (2024):**
- UCIe 2.0 specification in development (sub-10µm pitch, higher density)
- Few products shipping with UCIe 1.0 compliance yet (mostly proprietary variants compatible in spirit)
- TSMC CoWoS supports UCIe physical layer
- Intel's Open Platform approach includes UCIe support in future platforms

### Quick Check
> 1. What is UCIe and what problem does it solve for the chiplet ecosystem?
> 2. What is the bandwidth of UCIe 1.0 advanced package and how does it compare to PCIe?
> 3. Why do companies like Google and Meta care about UCIe?

---

## 7. Thermal Challenges in 3D

Stacking dies creates a fundamental problem: heat must escape through fewer surfaces.

```
Thermal path in 2D flip chip:
  
  Die → solder bumps → substrate → PCB → heatsink (from top)
  
  Thermal resistance: ~0.1°C/W at die surface
  CPU junction temperature target: <100°C
  
Thermal path in 3D stacked die:
  
  Top die → bonding layer → bottom die → substrate → heatsink
  
  Extra thermal resistance at each stacking interface
  Top die runs hotter: heat can't escape through the top
```

**AMD 3D V-Cache thermal solution:**
The SRAM cache die on top of the CPU die creates a thermal sandwich:
- SRAM generates less heat (lower power density than CPU cores)
- But it blocks heat from CPU cores below
- AMD's solution: lower the clock speed of cores under the cache die slightly; other cores run at full speed

**NVIDIA H100 SXM5 thermal:**
- 700W TDP on 814mm² die
- Power density: 860 mW/mm² — extremely high
- Liquid cooling mandatory for H100 SXM5
- CoWoS interposer: 800µm thick silicon + HBM stacks add thermal resistance to backside
- Specialized heat spreaders designed with vapor chambers

**3D stacking power delivery:**
Power must reach the top die through the bottom die (or through TSVs). This adds IR drop (resistance × current = voltage drop across the supply network).

**Backside Power Delivery Network (BSPDN):**
Intel's PowerVia (in Intel 18A process) and TSMC's proposed similar approach:
- Traditional: power rails on the front side compete with signal wires for routing resources
- BSPDN: power rails on the back side of the wafer → more space for signal routing on front, less IR drop

```
Backside power delivery (PowerVia):
  
  Front side (top):
    ┌──────────────────────────────┐
    │ Signal wires M1–M12          │  ← no power rails here
    │ Transistors (gates, source,  │
    │ drain at bottom of front)    │
    └──────────────┬───────────────┘
                   │ TSVs (power)
    ┌──────────────┴───────────────┐
    │ Back side (bottom):          │
    │ Power delivery network       │  ← VDD, VSS rails
    │ (thick copper, low R)        │
    └──────────────────────────────┘
    
Benefits: 5–6% performance improvement, lower power, better signal routing
```

### Quick Check
> 1. Why does 3D die stacking create thermal challenges?
> 2. How does AMD handle the thermal problem of 3D V-Cache stacking?
> 3. What is backside power delivery and what does it improve?

---

## Summary

- **Why 3D integration**: Die yield economics (smaller chiplets yield better), process optimization (each chiplet at optimal node), bandwidth density (fine-pitch bumps approach on-chip wire density).
- **Packaging hierarchy**: Wire bond → Flip chip (C4) → Advanced (EMIB, CoWoS, SoIC, hybrid bonding). Each step: higher I/O density, higher bandwidth, higher cost.
- **Intel EMIB**: Small silicon bridge embedded in organic substrate for high-bandwidth die-to-die connections. Cheaper than full interposer.
- **Intel Foveros**: 3D logic-on-logic stacking with µ-bumps. Foveros Direct uses copper-copper hybrid bonding at 10µm pitch.
- **TSMC CoWoS**: Silicon interposer in the package. CoWoS-S used in H100 (3.35 TB/s HBM3). CoWoS-L uses embedded bridges.
- **TSMC SoIC**: 3D face-to-face hybrid bonding. <10µm pitch. Used in AMD 3D V-Cache (96MB L3, 2 TB/s bandwidth).
- **AMD 3D V-Cache**: 64MB SRAM stacked on CPU die. 25% gaming improvement, 74% HPC improvement.
- **UCIe**: Open chiplet interconnect standard. Enables mix-and-match chiplets from different vendors.
- **Thermal**: 3D stacking adds thermal resistance. Solutions: lower power density on stacked dies, backside power delivery, liquid cooling.

---

## Exercises

### Easy
1. What is the yield advantage of chiplets over monolithic SoC?
2. What is the difference between CoWoS and SoIC?
3. What is UCIe and why do companies like Google want it to exist?

### Medium
4. Bandwidth calculation for H100 SXM5: TSMC CoWoS-S silicon interposer, 5 HBM3 stacks. (a) Each HBM3 stack: 1024-bit interface at 6.4 Gbps/pin → bandwidth per stack? (b) Total for 5 stacks? (c) If HBM3 pitch on interposer is 25µm, and each stack occupies a 10mm × 15mm area (150mm²), how many bump connections per stack? (d) A CoWoS interposer wire resistance: 0.05Ω/mm. For a signal traveling 20mm across interposer at 10 Gbps: is propagation delay or resistance the primary concern? (e) NVIDIA plans HBM4 at 1.6 TB/s per stack. With 6 HBM4 stacks: total bandwidth? What does this require from the interposer?
5. 3D V-Cache performance model: AMD Ryzen 7 7800X3D (8-core, 96MB L3 with V-Cache) vs Ryzen 7 7700X (8-core, 32MB L3 without). (a) In a game engine: 60% of L3 misses go to DRAM at 70ns; 40% are secondary hits at 5ns. Without V-Cache: assume 2% miss rate at L3. With V-Cache: assume 0.5% miss rate (more data fits). (b) Game loop: 1M L3 accesses per frame. Calculate average memory access cycles for each config (assume 3.8 GHz base clock). (c) If memory is the bottleneck for 30% of frames, and V-Cache reduces latency by 60% on those frames: expected FPS improvement? (d) For a database workload (32GB working set): V-Cache (96MB) still misses. Will V-Cache help this workload? Why?
6. UCIe economics: A chip startup wants to build an AI accelerator chiplet using UCIe 1.0 advanced package spec. (a) Required bandwidth: 4 TB/s between accelerator chiplet and memory chiplet (2× HBM-like). UCIe advanced: 96 GB/s/mm. How wide must the UCIe interface be (mm)? (b) The accelerator die is 200mm². If the UCIe interface is on one side of the die (side length = √200 ≈ 14mm): does the required width fit? (c) Power: UCIe advanced = 0.25 pJ/bit. At 4 TB/s: interface power consumption in watts? (d) If the accelerator die is fabricated at TSMC N5 and the memory die at TSMC N7: total tape-out cost (N5 NRE ~$30M, N7 NRE ~$15M)? What production volume justifies this?

### Hard
7. Hybrid bonding physics and limits: SoIC-F hybrid bonding (copper-copper, 9µm pitch). (a) Bonding area: copper pad of 4µm diameter, 9µm pitch → pad density = 1/(9µm)² = ? pads/mm². For 200mm² die: total bonds? (b) Bonding temperature: copper thermocompression bonding requires 250–300°C. The BEOL (back-end-of-line) metal stack can handle 400°C. DRAM cells are temperature-sensitive. Explain why this temperature constraint limits what you can stack on top of DRAM using SoIC. (c) Pad alignment: as pitch shrinks to 3µm (next gen), alignment accuracy requirement is <0.5µm (6-sigma). Wafer-to-wafer bonding alignment is ~0.5µm today. How does this affect yield as pitch decreases? (d) Electrical properties: at 9µm pitch, pad capacitance is ~2 fF, resistance ~0.1Ω. Data rate limit per pad from RC delay? If goal is 10 Gbps/pad: is this feasible? (e) Future: at 1µm pitch (possible with CFET-era processes), calculate pad density, total bonds for 200mm² die, and theoretical bandwidth at 5 Gbps/pad.
8. Full system analysis — NVIDIA Blackwell B100: Blackwell uses CoWoS-L with two GPU dies reticle-stitched into one logical GPU. (a) TSMC N4 process: 825mm² per die. Two dies stitched = ~1650mm²(+stitching overhead). At D₀=0.08/cm² (N4 yield): yield per die? Yield for stitched pair (assume 20% overhead from stitching defects)? (b) HBM3E specification: 1.2 TB/s per stack, 7 stacks = 8.4 TB/s. At 1024 bits wide, what data rate per pin to achieve 1.2 TB/s? (c) Power: B100 has ~1000W TDP on 2× 825mm² = 1650mm². Power density in W/mm²? Compare to H100 (700W/814mm²). What does this require from the cooling system? (d) NVIDIA uses a custom NVLink 5.0 die at 1.8 TB/s per GPU for multi-GPU connectivity. How many B100 GPUs can be in one NVLink domain (NVSwitch fabric capacity)? What bandwidth does each GPU see to the rest of the cluster? (e) Cost: a DGX B100 system (8× B100 + NVSwitch + InfiniBand) costs ~$150,000–$200,000. For a hyperscaler building 100,000 GPU clusters: infrastructure cost at each price point?
