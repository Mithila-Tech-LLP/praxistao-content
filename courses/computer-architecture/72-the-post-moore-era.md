# Chapter 72: The Post-Moore Era

For fifty years, Moore's Law and Dennard scaling together gave the industry a reliable "free lunch": wait two years and your software runs twice as fast on the same hardware — for the same price. That era is over. Transistors still get smaller (slowly), but they no longer get faster or cheaper at the same rate, and the power wall makes clock-speed scaling impossible. The post-Moore era demands entirely new thinking: instead of waiting for the next process node, we must innovate in architecture, packaging, specialization, and even fundamentally new computing substrates. This chapter maps the landscape of post-Moore strategies.

## Table of Contents

1. [What Moore's Law Actually Gave Us](#1-what-moores-law-actually-gave-us)
2. [Where It Broke Down](#2-where-it-broke-down)
3. [Post-Moore Strategy 1 — Specialization](#3-post-moore-strategy-1--specialization)
4. [Post-Moore Strategy 2 — 3D Integration](#4-post-moore-strategy-2--3d-integration)
5. [Post-Moore Strategy 3 — New Materials and Devices](#5-post-moore-strategy-3--new-materials-and-devices)
6. [Post-Moore Strategy 4 — Rethinking the Architecture](#6-post-moore-strategy-4--rethinking-the-architecture)
7. [The Timeline: When Does "Post-Moore" Fully Arrive?](#7-the-timeline-when-does-post-moore-fully-arrive)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Moore's Law Actually Gave Us

Moore's Law is often described as "transistors double every 2 years." But the true benefit came from a combination:

```
The "Moore's Law + Dennard Scaling" combination (1970–2004):

  Transistors double every 2 years (Moore's Law)
  +
  Voltage scales down with transistor size (Dennard Scaling)
  =
  Clock speed doubles every 2 years (1970→2004: MHz → GHz)
  Power stays constant (same die, faster clock)
  Cost per transistor halves every 2 years (more function per chip)

  Net effect for software: FREE performance improvements every 2 years
  
  Intel 4004 (1971): 2,300 transistors, 740 KHz, 1µm
  Intel Pentium 4 (2000): 42M transistors, 1.5 GHz, 180nm
  30 years: 18,000× more transistors, 2,000× faster clock
```

**Three distinct gifts:**
1. **Speed scaling**: Each generation faster clock (ended ~2004, Dennard breakdown)
2. **Cost scaling**: Cheaper per transistor (slowing since ~2018 for leading-edge)
3. **Density scaling**: More transistors per mm² (still happening, but slower)

This gave the software industry an implicit promise: write code; it will run faster every few years for free. This enabled the entire software industry — operating systems that grew more complex, databases with more features, games with better graphics — all without performance-aware development.

### Quick Check
> 1. What is the difference between Moore's Law and Dennard Scaling?
> 2. What three benefits did Moore's Law + Dennard Scaling give the industry?
> 3. Approximately when did clock-speed scaling end and why?

---

## 2. Where It Broke Down

**Dennard Scaling failure (~2004):**
- Threshold voltage (V_th) cannot scale below ~0.3V (leakage current becomes unmanageable)
- Leakage current ∝ e^(-V_th/kT) — exponentially worse as V_th decreases
- Intel hit the "power wall" with Pentium 4 Prescott (2004): 130W at 3.8 GHz on 90nm
- The CPU community pivoted to multicore (more cores at lower clock, same power)

**Cost scaling reversal (~2018):**
- TSMC N5 (5nm) costs ~$170M per mask set for tape-out; N7 was ~$100M
- Extreme process engineering (EUV, GAA, SADP) adds cost per transistor
- Industry inflection: N5 → N3 transistor density improved ~70%, but cost per transistor increased
- Smaller nodes now cost MORE per transistor for the first time in history

**Density scaling slowdown:**
- Historical: 2× density every ~2 years
- Current (2024): moving from N3 to N2 gives ~1.3× density improvement
- The cadence has slowed from 2-year cycles to 3-4 year cycles for meaningful nodes

```
The post-Moore gap:
  
  Expected (if trends continued):
    2024: 1nm node, $0.001/transistor, 10 GHz clocks
  
  Actual 2024:
    Best: TSMC N2 (first 2nm-class chips 2025)
    Cost: ~$0.005/transistor (5× more than projected)
    Clock: 4-5 GHz max (clock speed flat since 2004)
  
  The gap between expected and actual performance = 
  the opportunity that post-Moore strategies must fill
```

### Quick Check
> 1. What is the "power wall" and when did it emerge?
> 2. Why has cost per transistor stopped falling (and even increased) at leading-edge nodes?
> 3. What happened to the 2-year node cadence?

---

## 3. Post-Moore Strategy 1 — Specialization

If you can't make general-purpose faster, make special-purpose hardware for the most important workloads.

**The specialization insight**: A general-purpose CPU executes any instruction. For a specific workload (matrix multiply, video decode, encryption), a custom hardware block can be 10–1000× more efficient by removing all the general-purpose machinery.

```
General-purpose CPU executing matrix multiply:
  - Load instruction → decode → issue → execute unit
  - General ALU: any operation
  - Out-of-order engine: predicts which instructions to run next
  - Branch predictor, register file, TLB, cache hierarchy
  - All this for: one multiply-accumulate per cycle

Custom matrix multiply accelerator (systolic array):
  - Fixed grid of multiply-accumulate units
  - No instruction decode, no OOO, no branch prediction
  - Wires connect cells in a fixed pattern
  - 1024 MACs per cycle per row × 1024 rows = 1M MACs/cycle
  - 10,000× more MACs per unit of chip area
```

**Specialization examples:**
- Google TPU v4: 275 TOPS → 90% lower cost-per-inference vs GPU for TensorFlow workloads
- Apple Neural Engine: 38 TOPS at <5W → enables on-device ML inference on iPhone
- Intel QAT (QuickAssist): cryptography accelerator → 10× faster AES-256 than CPU
- Video codecs (HEVC/VP9/AV1 hardware decode): 50-200× lower power than CPU software decode
- Bitcoin ASIC (BM1397): SHA-256 hashing 100,000× more efficient than CPU

**The specialization trade-off:**
```
                CPU         GPU         DSA
Flexibility     Highest     High        Low/None
Efficiency      Lowest      Medium      Highest
Programmability Easiest     Medium      Hardest
Cost            Moderate    High        Variable (NRE + volume)
Failure risk    None        None        High (wrong workload = waste)
```

**Domain-Specific Architectures (DSA)**: coined by John Hennessy and David Patterson in their Turing Award lecture (2017). The era of specialization. Every major workload category (AI, networking, storage, crypto, graphics, codec) will have dedicated hardware. Chapter 74 covers DSAs in depth.

### Quick Check
> 1. Why is a custom matrix multiply accelerator so much more efficient than a CPU executing matrix multiply?
> 2. What is a Domain-Specific Architecture (DSA)?
> 3. What is the main risk when building a specialized accelerator?

---

## 4. Post-Moore Strategy 2 — 3D Integration

If we can't shrink transistors fast enough in 2D, go 3D: stack dies, integrate chiplets, and wire them together with extremely high-bandwidth interconnects.

**The 3D opportunity:**
- A 300mm wafer contains ~700 reticle-sized (858mm²) chips or ~100 large SoC die
- Instead of one monolithic SoC, build specialized chiplets and assemble them
- Each chiplet can use the optimal process node (logic at N3, DRAM at N10, analog at N28)
- Connect chiplets via fine-pitch bumps at near-on-chip bandwidth

**Bandwidth math for chiplet vs monolithic:**
```
On-chip global wire (same die):
  Bandwidth: ~100 TB/s for a cross-chip bus
  
Chiplet-to-chiplet (EMIB, 25µm pitch):
  UCIe 1.0: 1.3 TB/s/mm at 25µm pitch
  
Package-level (organic substrate):
  ~1 TB/s/mm at 55µm pitch
  
Board-level (PCIe 5.0 × 16):
  128 GB/s
  
3D stacking (HBM, hybrid bonding):
  Internal HBM bandwidth: 1-4 TB/s
  SoIC hybrid bonding: ~10 TB/s/mm²
```

The bandwidth cliff from on-chip to off-chip to board-level spans 3 orders of magnitude. 3D integration closes this gap by keeping more of the die interaction at near-on-chip bandwidth.

**3D integration approaches** (covered in depth in Chapter 73):
- **EMIB** (Intel Embedded Multi-die Interconnect Bridge): silicon bridge embedded in organic substrate
- **CoWoS** (TSMC Chip on Wafer on Substrate): silicon interposer
- **SoIC** (TSMC 3D stacking with hybrid bonding): face-to-face metal bonding at 1µm pitch
- **Foveros** (Intel 3D stacking): logic-on-logic with active base tile

**Post-Moore 3D economics:**
- Yield: a 400mm² die at 90% yield has ~70% individual die yield (Poisson with D₀=0.3/cm²). Four 100mm² chiplets: each at 97% yield → assembled yield 97%⁴ ≈ 88%. Yield economics favor chiplets.
- Cost: design once, assemble differently → AMD uses same CPU chiplets across EPYC 7003 (8 chiplets) and Ryzen (1 chiplet)

### Quick Check
> 1. What are the three reasons chiplets are economically attractive vs monolithic SoCs?
> 2. What is EMIB and how does it differ from SoIC?
> 3. What is the bandwidth difference between on-chip wires and off-chip PCIe links?

---

## 5. Post-Moore Strategy 3 — New Materials and Devices

Silicon MOSFET physics will eventually reach fundamental quantum limits. Research explores new materials and device structures.

**Gate-All-Around (GAA) and CFET:**
- FinFET (2011): three-sided gate wraps around fin
- GAA nanosheet (Samsung 3nm, TSMC 2nm): gate completely wraps nanosheet → better electrostatic control → less leakage
- CFET (Complementary FET): n-type and p-type transistors stacked vertically on the same footprint → 2× density vs lateral CMOS
- CFET is a key enabler for sub-1nm nodes (Intel's "angstrom era")

**2D materials:**
- Silicon MOSFET channel: 3D material → quantum tunneling limits at <2nm channel length
- **Graphene**: single atom thick, extremely high carrier mobility (10× silicon)
- **MoS₂ (Molybdenum Disulfide)**: 2D semiconductor, 0.67nm thick → thinnest possible channel
- Challenge: growing uniform, defect-free 2D material sheets at wafer scale
- IBM demonstrated a 40nm MoS₂ transistor (2022), proving sub-2nm gate length is feasible in 2D

**III-V semiconductor devices (GaAs, InP, GaN):**
- Much higher electron mobility than silicon
- Used in RF (phone modems), power electronics (GaN chargers), but not yet in logic
- Process integration with silicon CMOS is the challenge

**Carbon Nanotube FETs (CNFETs):**
- MIT demonstrated CNFET microprocessor in 2019 (16-bit, RISC-V instruction set, 14,000 transistors)
- Transistor current per width: 3× higher than silicon
- Challenge: CNT alignment, purity (metallic vs semiconducting tubes), wafer-scale manufacturing

**Spintronics:**
- Use electron spin (up/down) rather than charge (0/1) for information storage
- MRAM (Chapter 71) is the successful commercial product
- Spin-orbit torque (SOT) devices: faster writes than conventional MRAM
- Not yet competitive with CMOS logic

**The 1nm+ roadmap (ITRS / IRDS vision):**
```
Year  Node      Key innovation
────────────────────────────────────────────────────────
2022  3nm       GAA nanosheet (Samsung), EUV
2025  2nm       GAA nanosheet (TSMC N2), High-NA EUV ready
2027  1.4nm     Intel 14A (planned), CFET exploration
2030  1nm       CFET, High-NA EUV, possible 2D channel materials
2035  7Å (0.7nm) 2D channel (MoS₂), new memory structures?
```

### Quick Check
> 1. What is a CFET and how does it improve density over FinFET?
> 2. Why are 2D materials like MoS₂ attractive for future transistors?
> 3. What did MIT demonstrate with carbon nanotube FETs?

---

## 6. Post-Moore Strategy 4 — Rethinking the Architecture

The most radical post-Moore approaches don't just make existing architectures better — they propose entirely different computational substrates. (Quantum and neuromorphic computing were covered in Chapters 50–51; this section covers additional architectural innovations.)

**Approximate computing:**
- Many AI and signal processing applications don't need exact answers — approximate answers are acceptable
- **Stochastic computing**: encode values as bit-stream probability rather than binary number. Multiply becomes AND gate (probability × probability). Adder → multiplexer. Massive simplification.
- **Analog neural networks**: matrix-vector multiply in analog domain (power × bandwidth, no quantization overhead)
- Trade-off: lower precision → less energy per operation, but error management adds complexity

**Dataflow architectures:**
- Traditional von Neumann: program counter steps through instructions sequentially
- **Dataflow**: instructions fire when their inputs are ready (data-driven). Massive natural parallelism.
- **Spatial computing** (Intel/Altera Agilex, Wave Computing, Cerebras): entire network of processors connected by a fixed dataflow graph. No instruction fetch, no register file overhead.
- Cerebras Wafer-Scale Engine (WSE-3, 2023): entire 300mm wafer is one chip. 900,000 cores, 44GB on-chip SRAM. No off-chip memory needed for GPT-4 class models.

**Cerebras WSE-3 numbers:**
```
NVIDIA H100 SXM:              Cerebras WSE-3:
- Die: 814mm²                 - Die: 46,225mm² (entire wafer!)
- SRAM: 50MB                  - SRAM: 44GB (880× more)
- Cores: 16,896 CUDA cores    - Cores: 900,000 compute cores
- Power: 700W                 - Power: ~23,000W (liquid cooled)
- Peak AI: 3.9 PFLOPS         - Peak AI: ~4 PFLOPS at FP16
  
Key difference: WSE-3 has no HBM. The 44GB on-chip SRAM eliminates 
the memory bandwidth wall entirely for models that fit. For training 
LLMs with very long contexts: massive SRAM reduces stalls.
```

**Near-threshold computing:**
- Reduce supply voltage to just above threshold (~0.3V vs standard 0.8V)
- Dynamic power ∝ V² → 7× power reduction at the expense of ~10× slower operation
- Suitable for always-on sensor processing, IoT, wearables

**Photonic computing:**
- Use light instead of electrons to carry signals and perform computation
- Optical interconnects: already used inside data centers (fiber between servers)
- **Silicon photonics**: CMOS-compatible waveguides, modulators, photodetectors on silicon chips
- Optical matrix multiply: coherent light through a Mach-Zehnder interferometer mesh = analog matrix computation at the speed of light
- Companies: Lightmatter, Luminous Computing, Intel's silicon photonics division

### Quick Check
> 1. What is the Cerebras Wafer-Scale Engine and what problem does it solve?
> 2. What is near-threshold computing and what workloads is it suited for?
> 3. How can photonic computing perform matrix multiplication?

---

## 7. The Timeline: When Does "Post-Moore" Fully Arrive?

```
Post-Moore transition timeline:

  2004: End of clock scaling (Dennard breakdown, power wall)
  2012: Multicore mainstream (4-8 cores standard in laptops)
  2015: GPU compute mainstream (deep learning revolution begins)
  2017: Hennessy/Patterson "Golden Age of Computer Architecture" Turing lecture
  2018: Domain-specific accelerators mainstream (Google TPU, Apple Neural Engine)
  2020: Chiplets mainstream (AMD EPYC with 7nm chiplets on 12nm I/O die)
  2022: 3D integration commercial (TSMC CoWoS for H100, AMD 3D V-Cache)
  2024: GAA nanosheets (Samsung 3nm, TSMC 2nm announcement)
  2025: High-NA EUV first production (ASML EXE:5000)
  2027: CFET research → production readiness assessment
  2030: Possible transition to 2D channel materials or radical alternatives
  2035+: Quantum computing (error-corrected, fault-tolerant) starts impacting specific problems?
```

**What "post-Moore" means in practice:**
- Performance per dollar for general-purpose computing: flat or slowly improving
- Performance per dollar for specific workloads (AI, graphics, ML): still rapidly improving due to specialization
- Software optimization, algorithm improvement, and hardware specialization carry more of the performance improvement burden than they did in the "free lunch" era
- The best programmers and architects matter more than they did before — performance is no longer free

**The positive framing (Hennessy and Patterson):**
- We are entering a "Golden Age of Computer Architecture"
- Architecture innovation was suppressed when Dennard Scaling gave free speedups
- Now architects MUST innovate → the field is more intellectually active than ever
- RISC-V, SHAKTI, chiplets, neuromorphic, quantum, DSAs — all evidence of this renaissance

### Quick Check
> 1. When did each phase of Moore's Law breakdown occur (clock scaling, cost scaling, density scaling)?
> 2. Why do Hennessy and Patterson call this a "Golden Age" of computer architecture?
> 3. What will post-Moore mean for software development?

---

## Summary

- **What Moore's Law gave**: clock scaling (ended 2004), cost scaling (slowing since 2018), density scaling (slowing).
- **Dennard breakdown**: threshold voltage floor → leakage → power wall → end of clock scaling → multicore pivot.
- **Specialization (Strategy 1)**: 10–10,000× efficiency gains from domain-specific hardware (TPUs, Neural Engines, video codecs, crypto accelerators).
- **3D integration (Strategy 2)**: Chiplets, HBM, EMIB, CoWoS, SoIC — more transistors via packaging innovation, not necessarily smaller gates.
- **New materials (Strategy 3)**: GAA nanosheets (2nm), CFET (~1nm), 2D materials (MoS₂), carbon nanotubes — extending silicon CMOS.
- **Architecture rethinking (Strategy 4)**: Dataflow architectures, wafer-scale computing (Cerebras), near-threshold computing, photonics.
- **Timeline**: Fully "post-Moore" has been arriving gradually since 2004. The transition is a slowdown, not a cliff — but the implications for hardware and software are profound.
- **Golden Age**: Innovation is more important than ever because performance is no longer free.

---

## Exercises

### Easy
1. What did Dennard Scaling provide that Moore's Law alone did not?
2. Why did Intel pivot to multicore CPUs in 2004-2005?
3. What is a Domain-Specific Architecture (DSA) and give two examples?

### Medium
4. Cost-performance analysis of post-Moore strategies: A company needs 10 PFLOPS of AI inference at a cost under $10M/year (including power). (a) Option 1: General-purpose CPUs at $0.03/GFLOP/hour. Cost for 10 PFLOPS for 1 year? (b) Option 2: NVIDIA H100 at $30,000 each, 3.9 PFLOPS each. Number of GPUs? CapEx + at $4/hour cloud cost? (c) Option 3: Custom ASIC (NRE $50M, unit cost $500, 50 TOPS INT8 each). Number of chips? CapEx? Break-even vs GPU at what production volume? (d) Which option makes sense for: 1-year experiment, 3-year deployed service, 10-year deployed service?
5. Roadmap projection: Intel's 18A node (2025 planned): GAA RibbonFET, PowerVia backside power. (a) Historical density improvement per node: N5→N3: 1.7×, N3→N2: 1.3× projected. If Intel 18A is comparable to TSMC N2, what MTr/mm² should it achieve? (b) If H100 GPU transistors (80B on N4) were built on N2: how many transistors would fit in the same 814mm² die? (c) What performance improvement would those additional transistors deliver, if the architecture stays the same? If the architecture is redesigned to use them? (d) If CFET doubles density vs GAA at same node: what would a 2030 AI accelerator look like?
6. Specialization vs generality trade-off: A company is building a recommendation system. The bottleneck is embedding table lookup (10TB table, 1M queries/second). (a) General-purpose solution: 100 CPU servers × 512GB RAM each = 51TB RAM capacity. At $50,000/server: cost? (b) AxDIMM solution: 40 servers × 8 AxDIMM modules × 16GB each = 5.12TB RAM, but each AxDIMM does in-memory lookup. Hardware cost? (c) Custom ASIC (NRE $5M, $200/unit, 1TB on-chip DRAM, handles 1M lookups/second per chip): how many chips needed? Total cost? (d) The embedding table changes monthly (new products, updated embeddings). How does this affect the ASIC option?

### Hard
7. Cerebras WSE economics: Cerebras WSE-3 costs ~$4M per system. NVIDIA H100 costs $30,000 × 8 (in a DGX H100) = $240,000. (a) Compute parity: WSE-3 at 4 PFLOPS FP16, DGX H100 at 8×3.9 = 31.2 PFLOPS FP16. How many DGX systems equal one WSE-3 in compute? What is their total cost? (b) Memory bandwidth: WSE-3 has 44GB SRAM at ~20 PB/s internal bandwidth; DGX H100 has 640GB HBM3 at 8×3.35 = 26.8 TB/s + 600GB NVLink bandwidth between GPUs. For models that fit in WSE-3 SRAM: which is faster and by how much? (c) For models that don't fit in 44GB (e.g., GPT-4 class ~350GB): WSE-3 must partition across cores with limited bandwidth crossing. How does this change the comparison? (d) Power: WSE-3 uses ~23kW; DGX H100 uses ~10.2kW for similar compute performance (when DGX wins on compute). TCO over 3 years ($0.10/kWh): include CapEx, power, cooling for each option.
8. Design a 2030 AI training chip: You are designing a chip to train LLMs in 2030. Assume: TSMC N1.4 process (CFET, 1Å = 0.1nm features), High-NA EUV lithography, 3D stacking with SoIC. Target: 100 PFLOPS BF16, fit GPT-5 class model (500B parameters) on-chip, <10kW. (a) Transistor budget: at CFET N1.4, assume 800 MTr/mm² (projected from current trends). What die area for 200B transistors (H100 has 80B at 814mm²)? (b) On-chip SRAM: 500B parameters × 2 bytes/param (BF16) = 1TB needed. SRAM density at N1.4: ~100MB/mm². Area? Is this feasible? What technology supplement is needed? (hint: embedded DRAM, 3D stacking) (c) Power: P = αCV²f. From N3 to N1.4: V scales 0.8→0.55V, C scales with area, f scales 5→8 GHz. Estimate power reduction ratio vs a hypothetical N3 version of this chip. Can you hit 10kW? (d) The software stack: how does the compiler translate PyTorch code to run on your custom chip? What hardware features must be exposed (memory tiling, dataflow scheduling, quantization)?
