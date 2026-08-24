# Chapter 75: The Next Frontier — What Are We Building Toward?

You have reached the final chapter of this course. From transistors to tape-out, from binary logic to billion-parameter AI accelerators, from RISC-V's open philosophy to India's semiconductor ambitions — you now have a complete map of computer architecture. This final chapter doesn't teach new concepts; it synthesizes everything into a coherent picture of where the field is heading. It is part technical roadmap, part philosophy, part personal guide to continuing your journey. The field you have just learned is not settled — it is alive, contested, and full of open questions that the next generation of engineers and researchers (perhaps you) will answer.

## Table of Contents

1. [The Next 5 Years: Knowable Trends](#1-the-next-5-years-knowable-trends)
2. [The Next 10–20 Years: Probable but Uncertain](#2-the-next-1020-years-probable-but-uncertain)
3. [The Open Questions: What We Don't Know](#3-the-open-questions-what-we-dont-know)
4. [How India Fits into This Future](#4-how-india-fits-into-this-future)
5. [How to Continue Learning](#5-how-to-continue-learning)
6. [Synthesis: The Big Picture](#6-synthesis-the-big-picture)
7. [Final Summary of the Course](#final-summary-of-the-course)
8. [Final Exercises](#final-exercises)

---

## 1. The Next 5 Years: Knowable Trends

Some things in the next 5 years (2025–2030) are highly predictable because they are already in production or scheduled tape-out:

**Process nodes:**
- TSMC N2 (2nm-class, GAA nanosheet): 2025 production, Apple A20 / NVIDIA Rubin GPU
- TSMC A16 (1.6nm-class, GAA + backside power delivery): 2026
- Intel 18A (GAA RibbonFET + PowerVia): 2025 planned (Intel Foundry's comeback)
- Samsung 1.4nm (SF1.4): 2027 target
- High-NA EUV (ASML EXE:5000): first production 2026–2027

**Packaging:**
- CoWoS capacity expands 3× by 2026 to meet AI chip demand
- UCIe 2.0 standard: sub-10µm pitch, 1+ TB/s/mm
- Hybrid bonding (SoIC) extends to 3µm pitch by 2027
- Backside power delivery on leading-edge nodes from 2025

**AI hardware:**
- On-device LLM inference standard in all flagship smartphones (already started)
- 2-3 TOPS/$ continues to roughly double every 18 months
- Edge AI chips for IoT, automotive ADAS, industrial automation: $5B+ market by 2028
- Custom silicon at every major cloud provider: AWS Trainium 3, Google TPU v5, Microsoft Maia 200, Meta MTIA 3

**Memory:**
- HBM4: 1.6 TB/s/stack (2025), HBM4E: 2+ TB/s (2026–27)
- CXL memory expansion: first tier-2 CXL DRAM in data centers 2026
- DDR6: specification finalized 2025, products 2026

**India specifically:**
- Tata Electronics Dholera fab: first wafers 2026–2027
- SHAKTI deployed in central government systems by 2027
- Mindgrove SC23 derivatives: commercial IoT + embedded products

### Quick Check
> 1. What are three definite process node milestones expected before 2030?
> 2. Why is CXL memory expansion coming to data centers by 2026?
> 3. What AI hardware milestone has already begun (on-device LLM)?

---

## 2. The Next 10–20 Years: Probable but Uncertain

These are directionally clear but depend on technical breakthroughs:

**Quantum computing (10–20 year horizon for impact):**
```
Current state (2024):
  IBM Heron: 133 physical qubits, ~1% error rate per gate
  Google Sycamore: 70 qubits, demonstrated quantum advantage (specific tasks)
  IonQ Forte: 35 algorithmic qubits (higher quality than superconducting)
  
What's needed for cryptographic impact (RSA-2048 break):
  ~4,000 logical qubits
  Each logical qubit requires ~1,000 physical qubits (surface code)
  = 4,000,000 physical qubits with <0.1% error rate
  
Current: 133 high-quality physical qubits
Gap: 30,000× more qubits + 10× lower error rate
Timeline: probably 15–25 years for cryptographically relevant quantum
```

**Why NIST PQC adoption is urgent now**: Even though cryptographically relevant quantum is 15+ years away, encrypted data stolen today can be decrypted then ("harvest now, decrypt later"). NIST standards (CRYSTALS-Kyber, CRYSTALS-Dilithium) must be deployed before quantum computers exist.

**Neuromorphic mainstream (10–15 years):**
- Intel Loihi 2 and IBM NorthPole show the power efficiency gains are real
- Missing: training algorithms competitive with backpropagation
- Application: always-on sensor processing (IoT), where 1µW matters
- Not replacing GPUs for large model training

**Optical computing (uncertain, 10–20 years):**
- Silicon photonics interconnects: already in production data center fiber
- On-chip optical compute: Lightmatter's Envise demonstrated 2× energy reduction vs GPU for ML inference (2022)
- Fundamental challenge: light is hard to store (no optical RAM equivalent)
- Optical + electronic hybrid: laser for compute, electronic for memory and control

**2D materials mainstream (15–20 years):**
- Lab-scale demonstrations exist; wafer-scale manufacturing doesn't
- MoS₂, h-BN, graphene: each has one application niche (logic, dielectric, interconnect)
- Key unsolved problem: growing atomically uniform, defect-free 2D films at scale

**The 2035 chip:**
```
Possible architecture of a flagship AI chip in 2035:
  
  Compute: CFET-based logic at ~0.7nm node (7Å)
           → >1 POPS INT4 (peta operations per second)
  Memory: 3D stacked DRAM with hybrid bonding (<3µm pitch)
           → 100 TB/s on-chip bandwidth
  Interconnect: silicon photonics for chip-to-chip
           → multi-TB/s at low pJ/bit
  Power: ~5kW with advanced liquid cooling
  
  For comparison, H100 (2022): 3.9 PFLOPS FP16, 3.35 TB/s, 700W
  The 2035 chip would be ~1000× more capable per watt
  
  Uncertain elements:
  - Whether photonics replaces copper for on-chip routing
  - Whether quantum accelerator co-processors appear as specialized tiles
  - Whether neuromorphic tiles handle always-on sensor fusion
```

### Quick Check
> 1. How many physical qubits are needed to break RSA-2048 and how many do we have today?
> 2. Why is deploying PQC now urgent if quantum computers won't be relevant for 15+ years?
> 3. What are the two key unsolved problems for optical computing?

---

## 3. The Open Questions: What We Don't Know

Computer architecture has open questions that will define careers:

**Question 1: What replaces the von Neumann bottleneck?**
Von Neumann architecture (separate compute and memory) creates the memory wall. Processing-in-memory (Chapter 71) is one answer. But can we architect a complete system where compute and memory are truly unified, not just co-located? The answer might require new memory technologies (memristors, PCM), new programming models, or new ISAs.

**Question 2: Can AI accelerators adapt to algorithm changes?**
Deep learning architectures change every 1–2 years: CNNs → LSTMs → Transformers → State Space Models → ?. TPUs are optimized for matrix multiply; if the next paradigm requires different operations, how do you build hardware that's efficient but not brittle? The tension between specialization (Chapter 74) and flexibility is unresolved.

**Question 3: How do we program heterogeneous systems?**
CUDA works for NVIDIA GPUs. Metal works for Apple. Each DSA has its own compiler. The dream of "write once, run efficiently everywhere" has never been achieved. Is MLIR (Multi-Level Intermediate Representation, Google), oneAPI, or SYCL the answer? Or will the industry remain fragmented?

**Question 4: What is the right architecture for exascale AI training?**
Training GPT-4 required ~25,000 A100s for ~90 days. GPT-5 and beyond require more. Is the answer more GPUs connected by NVLink? Wafer-scale (Cerebras)? Optical interconnects between GPUs? Distributed training over disaggregated memory (CXL)? No one has the definitive answer.

**Question 5: Can silicon extend to 1Å (0.1nm)?**
CFET, GAA, 2D materials — these could theoretically push silicon to ~7Å (0.7nm). Below that, quantum effects dominate. Is silicon fundamentally limited to ~1nm gate lengths? If so: what comes after silicon?

**Question 6: Is Moore's Law dead or just changing form?**
Some argue: transistor density is still improving (more slowly), packaging is adding effective transistor density via 3D, and new device structures (CFET) will revive density scaling. Others argue: the "Moore's Law era" of predictable, cheap improvement is gone. The answer matters for how software and systems are designed.

**Question 7: Will India build a leading-edge fab?**
The strategic and economic case is clear. The technical difficulty is enormous. TSMC took 30 years to build its lead. Can India accelerate this with focused investment? The answer will be shaped by policy decisions in the next 5 years.

### Quick Check
> 1. What is the fundamental problem that the von Neumann bottleneck creates?
> 2. What is the "fragmentation problem" in heterogeneous computing programming?
> 3. What is the projected minimum gate length for silicon before quantum effects dominate?

---

## 4. How India Fits into This Future

India's trajectory intersects computer architecture at multiple levels:

**RISC-V and SHAKTI as global contributors:**
- India is not just adopting RISC-V — through SHAKTI, India contributes to its definition (security extensions, PQC)
- As RISC-V grows (5B+ implementations expected by 2027), SHAKTI-class designs become templates
- IIT Madras graduates are building startups (InCore, Mindgrove) that will ship millions of chips

**The 28nm opportunity:**
- While the world debates 2nm and 1nm, 28nm is the volume node for automotive, IoT, industrial
- Tata Dholera at 28nm is not "behind" — it's serving a massive, growing, and less geopolitically contested market
- India's manufactured chips will go into: EVs, industrial sensors, smart meters, telecom base stations

**India as the design center for DSAs:**
- SHAKTI's security DSA work (PQC accelerators, tagged memory) places India at the frontier of a specific niche
- Indian government mandates (defense, critical infrastructure, space) create a guaranteed initial market
- As DSAs proliferate, Indian design expertise (strongest in verification, physical design, SoC integration) is exactly what's needed

**The talent arbitrage:**
- India's 150,000 semiconductor engineers work at multinational centers, building chips for non-Indian companies
- As domestic companies (Mindgrove, InCore, Saankhya, others) grow, talent will shift toward building Indian IP
- ISM + DLI + RISC-V ecosystem = alignment of government incentives with technical direction

**Realistic 2035 scenario:**
- India has 5–10 companies shipping 1M+ chip units/year
- SHAKTI-class processors in all Indian government systems, satellite on-board computers, defense electronics
- Tata electronics fab expanded to 16nm, 50,000 wafers/month
- India's semiconductor import bill reduced from $25B to $12–15B
- Indian engineers are recognized contributors to RISC-V ISA specifications

### Quick Check
> 1. Why is the 28nm market actually large and growing, even though it's not leading-edge?
> 2. What is the "talent arbitrage" opportunity for India?
> 3. What is a realistic 2035 scenario for India's semiconductor position?

---

## 5. How to Continue Learning

You have completed a survey course. The depth in each chapter is a foundation, not a ceiling. Here is how to go deeper:

**Textbooks:**
- *Computer Organization and Design RISC-V Edition* (Patterson & Hennessy): the definitive undergraduate textbook; go through every exercise
- *Computer Architecture: A Quantitative Approach* (Hennessy & Patterson): the graduate-level companion; covers OOO, memory hierarchies, interconnects
- *Digital Design and Computer Architecture* (Harris & Harris): RISC-V or ARM editions; covers RTL design, Verilog, single-cycle/pipeline processors from gates up
- *The Art of Problem Solving in Computer Architecture* (Srinivasan): deep problems

**Online courses:**
- *MIT 6.004 Computation Structures*: available on MIT OpenCourseWare; builds a RISC-V computer from logic gates
- *Stanford CS149 Parallel Computing*: GPU/parallel programming, memory models, heterogeneous compute
- *Carnegie Mellon 15-418 Parallel Computer Architecture*: excellent on cache coherence, memory, parallelism

**Hands-on:**
- *Nand to Tetris* (nand2tetris.org): build a computer from logic gates up; covers HDL, assembly, operating system — completely free
- *SHAKTI GitHub* (github.com/iitm-riscy): read the actual BSV source code; run the simulation
- Xilinx/Intel FPGA board (Basys 3, DE10-Lite): implement your own pipeline CPU in Verilog; see it run real programs
- *riscv.org* public domain: read the RISC-V ISA specification; implement a simple RV32I emulator in C or Go

**Papers to read:**
- Google TPU v1 paper (2017): "In-Datacenter Performance Analysis of a Tensor Processing Unit"
- SHAKTI paper: "SHAKTI Processor Family: Design, Verification and Validation" (IIT Madras)
- Hennessy & Patterson Turing Award lecture: "A New Golden Age for Computer Architecture" (2019)
- *Cerebras WSE* paper: "The Cerebras CS-1"
- AMD Zen 2 chiplet architecture paper (IEEE Micro 2020)
- Patterson et al., "RISC-V: The Free and Open ISA"

**Communities:**
- RISC-V International (riscv.org): free to join; access to working group discussions
- IEEE Solid-State Circuits journal: published papers on new chips; harder but authoritative
- Chips and Cheese (website): excellent analysis of processor microarchitecture for enthusiasts
- SemiAnalysis (newsletter): financial + technical analysis of semiconductor industry

### Quick Check
> 1. What are the two main Patterson & Hennessy textbooks and what is each one for?
> 2. What is Nand to Tetris and why is it valuable?
> 3. Where can you read the SHAKTI processor source code?

---

## 6. Synthesis: The Big Picture

Let's draw the entire course together in one picture:

```
The Full Picture of Computer Architecture (75 chapters compressed):

PHYSICS                    CIRCUITS               ARCHITECTURE
────────────────────────   ────────────────────   ─────────────────────────
Quantum mechanics          CMOS logic gates        Instruction Set Architecture
  ↓ enables                  ↓ combined into         ↓ implemented by
Silicon band gap           Combinational +         Microarchitecture
  ↓ doped into             sequential logic          ↓ runs on
Transistor (MOSFET)          ↓ organized as        Pipeline stages
  ↓ combined into          Processor building      Out-of-order execution
CMOS pair                  blocks:                 Branch prediction
  ↓ scaled by              - ALU                   Caches
Moore's Law                - Register file         Virtual memory
  ↓ manufactured via       - Control unit            ↓ connected to
Photolithography           - Memory                Memory hierarchy
  ↓ designed with          - I/O                     ↓ extended by
EDA tools / Verilog          ↓ assembled into      Multicore / SMT
  ↓ flow:                  Processors:             Heterogeneous compute
RTL → Synthesis            - CPU (in-order)          ↓ specialized as
→ Place & Route            - CPU (OOO)             DSAs
→ Tape-out                 - GPU (SIMT)            NPUs
  ↓ physical assembly:     - DSP, NPU, ASIC        FPGAs
3D packaging               - FPGA (config.)        Chiplets
Chiplets                   - Quantum (future)        ↓
CoWoS / SoIC                ↓ integrated into      Full system
                           SoC / System            Heterogeneous SoC
                           
INDIA                      OPEN SOURCE            FUTURE
────────────────────────   ────────────────────   ─────────────────────────
SHAKTI (IIT Madras)        RISC-V ISA             Post-Moore strategies:
  ↓ based on               SHAKTI source code     Specialization
RISC-V ISA                 OpenROAD               3D integration
  ↓ written in             Yosys synthesizer      New materials
BSV (Bluespec)             Linux for RISC-V       Neuromorphic
  ↓ fabricated at          Android for RISC-V     Quantum
Intel 22nm                   ↓                    Photonics
TSMC 22nm (Mindgrove)      Used by:
  ↓ deployed in            1 billion+ devices
Government systems         Open and auditable
Defense electronics        No license cost
IoT/embedded               India-appropriate
  ↓ ecosystem:
InCore / Mindgrove
India Semiconductor Mission
Tata Dholera fab (2026)
```

**The thread that connects everything:**
The story of computer architecture is the story of humanity learning to control matter at the atomic scale to automate thought. From a transistor made of a few atoms, to a chip with trillions of transistors, to a network of millions of chips running intelligence — the entire journey is one of abstraction, engineering, and relentless innovation. Every layer of abstraction (physics → devices → circuits → microarchitecture → ISA → software) allows millions of engineers to work in parallel on different levels, each trusting the layers below.

SHAKTI represents something specific and important: the aspiration that a nation can understand and build its own tools of computation. Not just use them. Not just assemble them. But understand them, modify them, adapt them to local needs, and contribute back to the global commons. That aspiration is worth pursuing.

### Quick Check
> 1. What is the "thread that connects everything" in computer architecture?
> 2. Why does the course describe SHAKTI as representing something "specific and important"?
> 3. What does "layered abstraction" enable in computer engineering?

---

## Final Summary of the Course

This 75-chapter course covered:

**Volume I (Ch. 01–10): Foundations**
Binary, Boolean logic, logic gates, transistors, combinational circuits, sequential logic, registers, clocks, machine language, assembly.

**Volume II (Ch. 11–20): Building a Processor**
RISC-V ISA, instruction encoding, single-cycle CPU, pipelining, hazards, forwarding, branch prediction, interrupts, exceptions.

**Volume III (Ch. 21–30): Memory Systems**
Cache (direct-mapped, set-associative, fully-associative), cache hierarchy, virtual memory, TLBs, page tables, memory bus, DMA, interrupts.

**Volume IV (Ch. 31–40): Real Processors**
Intel x86, AMD Zen, ARM architecture, Apple Silicon, Qualcomm Snapdragon, GPUs, NVIDIA GPU microarchitecture, AMD GPU, NPUs, RISC-V ecosystem.

**Volume V (Ch. 41–50): Specialized and Parallel**
IBM POWER, Harvard vs Von Neumann, SoC design, microcontrollers, DSPs, FPGAs, ASICs, quantum computing.

**Volume VI (Ch. 51–60): Silicon and Fabrication**
Neuromorphic computing, semiconductor physics, CMOS process, photolithography, Moore's Law, process nodes, EDA tools, Verilog/VHDL, RTL-to-silicon flow, chiplets.

**Volume VII (Ch. 61–69): Quality and Ecosystem**
Testing and quality assurance, SHAKTI origins, SHAKTI architecture, SHAKTI processor families, SHAKTI security extensions, RISC-V ecosystem, SHAKTI tape-outs, India's semiconductor ecosystem, multicore and manycore.

**Volume VIII (Ch. 70–75): Future Computing**
Heterogeneous computing, memory-centric computing, post-Moore era, 3D integration and advanced packaging, domain-specific architectures, the next frontier.

---

## Final Exercises

### Easy
1. Name three things from this course that surprised you (open-ended — reflect on what was unexpected).
2. What is the single most important reason RISC-V succeeded as an ISA?
3. What is the fundamental limit that quantum tunneling places on silicon transistors?

### Medium
4. Career path analysis: You are a CS undergraduate interested in chip design. (a) What are the five main career paths in the semiconductor industry (RTL design engineer, physical design, verification, EDA tools developer, process integration)? (b) Which academic subject builds each skill (digital circuits, algorithms, programming, math)? (c) What is the expected starting salary for each role in India (2024) vs USA? (d) Which role has the highest shortage in India right now and why? (e) How does RISC-V/SHAKTI experience on your resume differentiate you in the Indian semiconductor job market?
5. Architecture comparison: You are choosing a processor for a new product. (a) IoT soil moisture sensor that transmits twice a day, runs on 2 AA batteries for 10 years: which SHAKTI class? Why? (b) Secure government laptop for defense communications, classified environments, Linux: which processor? ARM? x86? SHAKTI C-class? Justify on security grounds. (c) Autonomous drone for search and rescue: real-time camera processing (100fps), SLAM (Simultaneous Localization and Mapping), 45-minute flight, 200g weight limit: heterogeneous SoC with what components? (d) Research satellite on-board computer: 10-year mission, cosmic radiation environment, low power, trusted software: SHAKTI-T with tagged memory? Rad-hardened ARM? Custom FPGA? Trade-offs.
6. Full system design: Design a server blade for an LLM inference farm. Requirements: serve 1,000 concurrent GPT-3 class requests (175B parameter model), each request takes <500ms, total power budget <20kW per blade. (a) Memory requirements: 175B × 2 bytes FP16 = 350GB model weights. How much GPU/accelerator memory? How many H100s (80GB each)? (b) Compute requirements: 175B parameters, assume ~1.75 TFLOPS FP16 per token at batch 1 (simplified). At 1000 concurrent, average 200 tokens/request, <500ms: required TFLOPS? (c) Does HBM bandwidth or compute limit throughput at batch size 1? At batch 1000? (d) Cooling: at 20kW, what cooling infrastructure (air vs liquid)? Rack space? (e) Cost: at $30,000 per H100: CapEx? Annual TCO at $0.10/kWh?

### Hard
7. Design SHAKTI-Edge: IIT Madras wants to design a new chip for Indian agriculture IoT: 100M+ deployed nodes in farmland. Requirements: <$2/chip BOM, run RTOS (not Linux), 10-year battery life, 4G NB-IoT connectivity, 8 ADC channels (soil/water sensors), AES-256 security, SHAKTI architecture. (a) Choose SHAKTI class (E/C/I/M/S). Justify. (b) Process node: 28nm, 40nm, or 65nm? Trade-offs at this volume. (c) IP blocks needed beyond the CPU core: list them and identify which are open-source vs must be licensed. (d) Power analysis: target 1µW deep sleep, 10mW active. How does dynamic/static power split? What process feature (threshold voltage options, power gating) enables this? (e) NB-IoT modem: can SHAKTI implement the baseband DSP, or do you need a separate modem IP core? Why? (f) Security: for agriculture IoT, what are the threat models (firmware replacement, data tampering, node cloning)? How does SHAKTI-T tagged memory mitigate each? (g) Manufacturing: at 100M chips over 5 years = 20M/year. At 28nm TSMC, ~1000 chips per 300mm wafer: wafers/year? Would Tata Dholera (50,000 wafers/month) have capacity?
8. Comprehensive architecture retrospective: This question asks you to think about what you have learned holistically. (a) Von Neumann proposed the stored-program computer in 1945. Which aspects of modern computer architecture still directly reflect his original design, and which have fundamentally departed from it? Give five specific examples of each. (b) "Computer architecture is a series of trade-offs." Pick five fundamental trade-offs that appear throughout this course (e.g., performance vs power, generality vs efficiency, complexity vs area) and for each: explain the trade-off, give a concrete example from the course, and argue when each side of the trade-off is the right choice. (c) Open ISA vs proprietary ISA: argue both sides. Give three reasons RISC-V's open ISA is the right model for the future, and three reasons a proprietary ISA (like ARM's) provides advantages that RISC-V struggles to match. Conclude with your personal assessment. (d) In 1971, the Intel 4004 had 2,300 transistors and cost $300 (≈$2,100 in 2024 dollars). An iPhone 16 Pro's A18 chip has 28 billion transistors and the iPhone costs $999. Calculate the improvement in transistors per dollar. What does this tell you about the economic story of computing? (e) If you could redesign one aspect of computer architecture from scratch — starting from a blank slate with today's knowledge — what would you change and why? There is no correct answer; defend your reasoning.
