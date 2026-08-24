# Chapter 71: Memory-Centric Computing

Traditional computer architecture keeps computation and memory separate: the CPU computes, memory stores data, and a bus carries data back and forth. This simple model worked for decades, but modern workloads — especially AI, database queries, genomics, and big data analytics — are fundamentally limited not by compute speed but by memory bandwidth. Moving data from memory to the processor and back uses more time and energy than the computation itself. Memory-centric computing (also called "near-memory computing" or "processing-in-memory") proposes a radical shift: bring the computation to the data.

## Table of Contents

1. [The Memory Wall](#1-the-memory-wall)
2. [Near-Memory Computing (NMC)](#2-near-memory-computing-nmc)
3. [Processing-in-Memory (PIM)](#3-processing-in-memory-pim)
4. [High Bandwidth Memory (HBM) and 3D Stacking](#4-high-bandwidth-memory-hbm-and-3d-stacking)
5. [Storage-Class Memory — Bridging DRAM and SSD](#5-storage-class-memory--bridging-dram-and-ssd)
6. [Compute Express Link (CXL)](#6-compute-express-link-cxl)
7. [Real Systems: Samsung HBM-PIM, UPMEM, Micron Automata](#7-real-systems-samsung-hbm-pim-upmem-micron-automata)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. The Memory Wall

**The memory wall** is the growing gap between processor speed and memory speed:

```
Year 1980: CPU speed ≈ memory speed (both grew together)

Year 2024:
  Processor (FLOPS/s) doubles every 2 years
  DRAM bandwidth doubles every 4-5 years
  
  Modern CPU peak FP64:  1,000+ GFLOPS
  Modern DRAM bandwidth: 50–200 GB/s
  
  Arithmetic intensity needed for full utilization:
    1,000 GFLOPS / 50 GB/s = 20 FLOPs/byte
  
  Many workloads do < 2 FLOPs/byte — memory-bound by 10×
```

**Energy breakdown for an ML operation:**
```
Operation              Energy (pJ)
─────────────────────────────────
32-bit ADD in CPU         1
SRAM read (near core)     5
L2 cache read            50
L3 cache read           100
DRAM read             2,000
Off-chip interconnect 5,000

DRAM access costs 2000× more energy than the operation itself!
```

**The roofline model** (introduced informally in Ch. 70) precisely captures this:
- X-axis: **arithmetic intensity** (FLOPs / bytes read from memory)
- Y-axis: **achieved performance** (GFLOPS)
- Roofline ceiling: min(arithmetic intensity × bandwidth, peak FLOPS)
- Most AI inference operators: low arithmetic intensity → memory-bandwidth-bound

**Memory-bound workloads** (arithmetic intensity < ~10 FLOPs/byte):
- LLM inference (attention is O(n²) memory, O(n²) compute → low intensity)
- Sparse neural networks
- Database scans and aggregations
- Genome sequence alignment
- Graph analytics (PageRank, BFS — irregular memory access)

**The fundamental question**: If memory access is 99% of the time and energy — why not do computation *inside* or *next to* the memory?

### Quick Check
> 1. What is the memory wall and when did it become a significant problem?
> 2. What is "arithmetic intensity" and how does it relate to memory-bound workloads?
> 3. How much more energy does a DRAM access consume vs a 32-bit ADD?

---

## 2. Near-Memory Computing (NMC)

**Near-memory computing** places compute logic close to DRAM — on the same package, or in the logic die beneath stacked DRAM — to reduce data movement distance.

```
Traditional:                       Near-Memory:
  
  ┌──────────┐                       ┌─────────────────────┐
  │   CPU    │                       │  DRAM array (top)   │
  └────┬─────┘                       ├─────────────────────┤
       │ memory bus (50 GB/s)        │  Logic layer (near) │
       │                             │  ┌────┐  ┌────┐     │
  ┌────┴─────┐                       │  │ALU │  │ALU │     │
  │   DRAM   │                       │  └────┘  └────┘     │
  └──────────┘                       └─────────────────────┘
       
  Data travels cm at 50 GB/s         Data travels µm at TB/s
```

**Types of near-memory placement:**

1. **In-package**: Logic chip and DRAM stack in the same package, connected via TSVs (Through-Silicon Vias) or HBM stacking. Bandwidth: 1–4 TB/s. Example: NVIDIA H100 has HBM3 at 3.35 TB/s.

2. **In-module**: Processing logic on the DRAM module PCB (e.g., near each DIMM). Less bandwidth gain but easier manufacturing.

3. **In-stacked memory**: Logic layer embedded within HBM stack, between DRAM dies. Most aggressive form — compute is literally inside the memory package.

**NMC advantages:**
- Internal bandwidth of HBM: 4 TB/s+ (vs 50 GB/s CPU→DRAM)
- Reduced data movement → lower energy
- Higher parallelism: many small processors per DRAM vault

**NMC challenges:**
- Logic in or near DRAM uses a different process node than CPU logic (DRAM is specialized — dense but slow transistors). This limits the complexity of near-memory logic.
- Heat: DRAM is sensitive to temperature. Adding computation generates heat, potentially causing DRAM errors.
- Programming: existing code doesn't run transparently on near-memory logic; programmer or compiler must explicitly schedule.

### Quick Check
> 1. What is near-memory computing and how does it address the memory wall?
> 2. Why is the logic layer in near-memory computing limited in complexity?
> 3. What is the bandwidth difference between the CPU-DRAM bus and internal HBM bandwidth?

---

## 3. Processing-in-Memory (PIM)

**Processing-in-memory (PIM)** goes further: computation happens inside the DRAM arrays themselves, using the analog properties of memory cells (for analog PIM) or small digital logic placed between array rows (for digital PIM).

**Digital PIM:**
- Small ALUs, adders, comparators placed in the "peripheral circuitry" of DRAM arrays
- Can scan rows, apply predicates, aggregate — without sending data off-chip
- Samsung HBM-PIM (2021): 16 floating-point SIMD units per HBM layer, each operating on their local DRAM vaults

**Analog PIM (Compute-in-Memory, CIM):**
- Use memory cell's analog physics to perform multiply-accumulate (MAC) operations
- Crossbar array: voltage on row wire, resistance of cell (programmed as matrix weight), current on column wire
- Column current = sum of (voltage × conductance) = dot product
- One clock cycle performs the entire matrix-vector multiply for that row/column combination

```
Crossbar array (analog PIM):
  
       x₁   x₂   x₃       ← input voltages (vector)
        │    │    │
  ──W₁₁─┼──W₁₂──W₁₃──→ I₁ (= W₁₁x₁ + W₁₂x₂ + W₁₃x₃)
  ──W₂₁─┼──W₂₂──W₂₃──→ I₂
  ──W₃₁─┼──W₃₂──W₃₃──→ I₃
        ↓    ↓    ↓
     column currents = matrix-vector product
  
  Done in ONE cycle, by physics (Ohm's Law + Kirchhoff's Current Law)
  
  Uses: resistive RAM (ReRAM), phase-change memory (PCM), 
        or flash cells as the programmable resistors
```

**Analog PIM challenges:**
- Precision: analog circuits are noisy; 4-bit or 8-bit precision is practical; 16-bit is hard
- Variability: each cell has different resistance; calibration required
- Programming memory weights takes time (erase/write cycles)
- Non-idealities: sneak path currents, IR drop, temperature sensitivity
- Competing with digital TPUs at 16-bit precision is difficult; works for inference where 8-bit is acceptable

**Applications that benefit most from PIM:**
- Transformer attention layers (O(n²) operations, low arithmetic intensity)
- Database: filter/aggregate on in-memory columns
- Genomics: pattern matching (Smith-Waterman alignment)
- Sparse matrix operations

### Quick Check
> 1. What is the difference between digital PIM and analog PIM?
> 2. How does a crossbar array perform matrix-vector multiplication?
> 3. Why is precision a challenge for analog PIM systems?

---

## 4. High Bandwidth Memory (HBM) and 3D Stacking

**HBM** is the current state-of-the-art for high-bandwidth DRAM, widely used in GPUs, FPGAs, and AI accelerators.

```
HBM3 stack structure:
  
  ┌─────────────────────────────────────────┐
  │ DRAM die 8 (top)                        │
  ├─────────────────────────────────────────┤
  │ DRAM die 7                              │
  ├─────────────────────────────────────────┤
  │ ... (8 dies total)                      │
  ├─────────────────────────────────────────┤
  │ DRAM die 1 (bottom)                     │
  ├─────────────────────────────────────────┤
  │ Base die (I/O + PHY + ECC + vaults)    │
  └───────────────────────────────────────┘
         ↕ TSVs (thousands of vertical wires)
  ┌─────────────────────────────────────────┐
  │ Silicon interposer                      │
  └─────────────────────────────────────────┘
         ↕ micro-bumps
  ┌─────────────────────────────────────────┐
  │ GPU / accelerator die                   │
  └─────────────────────────────────────────┘
```

**HBM specifications:**
```
Version  Year  Bandwidth/stack  Width  Capacity  Process
───────────────────────────────────────────────────────
HBM1     2015  128 GB/s         1024b  4GB        -
HBM2     2016  256 GB/s         1024b  8GB       20nm
HBM2E    2019  460 GB/s         1024b  16GB       -
HBM3     2022  819 GB/s         1024b  24GB      10nm
HBM3E    2024  1.2 TB/s         1024b  36GB       8nm

H100 uses 5× HBM3 stacks: 5 × 819 GB/s = ~3.35 TB/s total
A100 uses 5× HBM2E stacks: 5 × 460 GB/s = 2.0 TB/s total
```

**Why HBM is so fast**: 1024-bit wide interface per stack (vs 64-bit for regular DDR DRAM). The width × frequency gives the bandwidth. TSVs allow thousands of wires between DRAM and logic die, which simply is not possible with standard DRAM packages.

**Cost**: HBM is ~10× more expensive per GB than GDDR or DDR. Feasible only for AI accelerators and high-end GPUs, not general-purpose RAM.

**HBM-PIM (Samsung, 2021)**: Samsung integrated 2 GFLOPS of processing in each HBM-PIM stack. For a GPU with 5 stacks: 10 GFLOPS of in-memory compute + 1.2 TB/s bandwidth. First commercially available PIM product.

### Quick Check
> 1. What is HBM and why is it much faster than DDR DRAM?
> 2. How many HBM3 stacks does an NVIDIA H100 use and what total bandwidth does this provide?
> 3. What did Samsung add to HBM-PIM?

---

## 5. Storage-Class Memory — Bridging DRAM and SSD

The traditional memory hierarchy has a large gap between DRAM (fast, expensive, volatile) and SSD (slow, cheap, non-volatile):

```
DRAM:  ~100ns latency, $5–10/GB, volatile (data lost on power off)
SSD:   ~100µs latency (1000× slower), $0.05–0.20/GB, non-volatile
```

**Storage-class memory (SCM)** tries to fill this gap: cheaper than DRAM, faster than SSD, ideally non-volatile.

**Intel Optane / 3D XPoint (2017–2022):**
- Technology: PCM (Phase Change Memory) — uses heat to switch material between amorphous (high resistance = 0) and crystalline (low resistance = 1)
- Latency: ~300ns (3× slower than DRAM, but 300× faster than SSD)
- Persistence: data survives power-off
- Byte-addressable (unlike SSD/NVMe which are block-addressed)
- Intel Optane DIMM: used in the same memory slot as DRAM (DDR4 interface), supported on Xeon servers
- Intel discontinued Optane (2022–2023): market didn't adopt; NAND SSD got faster; economics didn't work

**MRAM (Magnetoresistive RAM):**
- Uses magnetic tunnel junctions (spin polarization) to store bits
- Fast (~10ns), truly non-volatile, unlimited endurance
- Challenge: low density vs DRAM (cells are larger)
- Current use: embedded non-volatile for microcontrollers (STMicroelectronics, Everspin)
- Not yet viable as main memory replacement

**FeRAM (Ferroelectric RAM):**
- Uses polarization of ferroelectric material
- Very fast writes (<100ns), non-volatile
- Low density, limited to small embedded uses

**CXL memory (Chapter 6)**: Not a new memory technology but a new protocol for connecting memory pools. Works with standard DRAM.

### Quick Check
> 1. What is the latency gap that storage-class memory tries to fill?
> 2. What technology did Intel Optane use and why was it different from DRAM?
> 3. Why did Intel discontinue Optane?

---

## 6. Compute Express Link (CXL)

**CXL (Compute Express Link)** is an interconnect standard that allows CPUs, GPUs, FPGAs, and memory devices to share a coherent memory space over PCIe-physical links. It's not a new memory technology — it's a protocol that enables new ways to compose systems.

**CXL layers:**
- **CXL.io**: PCIe-compatible I/O (device enumeration, register access, similar to PCIe)
- **CXL.cache**: Accelerator-to-CPU cache coherency (accelerator can cache CPU memory and vice versa)
- **CXL.mem**: CPU-to-device memory access (CPU accesses device's memory as if it were regular memory)

```
CXL memory expansion example:
  
  ┌─────────────────────────────────────────────────────┐
  │ CPU + 64GB DDR5 (local DRAM)                        │
  └──────────────────────────────────┬─────────────────┘
                                     │ CXL 2.0 link (PCIe 5.0 ×16)
  ┌──────────────────────────────────┴─────────────────┐
  │ CXL Memory Expansion Module: 256GB DRAM            │
  │ (Samsung CXL DIMM, Micron CZ120, etc.)             │
  └────────────────────────────────────────────────────┘
  
  CPU sees: 320GB total memory, coherent
  CXL DRAM latency: ~200ns (vs 70ns local) — 3× slower but cheaper
  Use case: memory-hungry apps (databases, LLM serving) that don't
            need the lowest latency on all their data
```

**CXL fabric (CXL 3.0, 2022):**
- Switch: multiple CPUs and devices share a CXL switch, creating a memory mesh
- Shared memory: CPU A and GPU B can access the same physical memory pool
- Memory pooling: unused memory from one server available to another (disaggregation)
- This enables **memory disaggregation**: memory as a separate infrastructure resource in data centers

**Why CXL matters for AI:**
- LLM serving requires huge KV cache (key-value attention cache): hundreds of GB per deployed model
- CXL memory expansion allows servers to have 1TB+ of accessible memory at low additional cost
- The ~200ns latency vs 70ns for local DRAM is acceptable for KV cache (since GPU compute time >> memory latency at these scales)

**CXL vendors**: Marvell, Astera Labs, Montage Technology (CXL controllers), Samsung/Micron/SK Hynix (CXL DRAM modules).

### Quick Check
> 1. What does CXL stand for and what problem does it solve?
> 2. What are the three CXL protocol layers (CXL.io, CXL.cache, CXL.mem) and what does each do?
> 3. Why is CXL memory useful for LLM inference serving?

---

## 7. Real Systems: Samsung HBM-PIM, UPMEM, Micron Automata

**Samsung HBM-PIM (Aquabolt-XL, 2021):**
- First commercially available PIM product
- 16 GFLOPS of floating-point compute inside HBM2 stack
- SIMD units in each DRAM vault's logic layer
- Optimization: softmax and GELU (neural network activations) — these are memory-bandwidth bound
- Result: 2.4× performance improvement for transformer inference vs standard HBM2, 70% energy reduction
- Deployed by Samsung in AI servers

**UPMEM PIM DRAM:**
- Standard DDR4 DIMM form factor with DPUs (DRAM Processing Units) inside
- Up to 2,500 DPUs per server (each DPU has 24MB SRAM + 64MB DRAM)
- DPU is a simple 32-bit processor (no out-of-order, no caches)
- Software: C code compiles to DPU ISA; host CPU dispatches tasks
- Use cases: genomics (BWA-MEM DNA alignment), database (sort, selection), fraud detection
- Benchmark: DNA alignment 19× faster than CPU-only with same total power

**Micron Automata Processor:**
- Implements non-deterministic finite automata (NFA) entirely in DRAM
- Uses DRAM cells to represent NFA states and transitions
- One DRAM access = one NFA state transition for millions of parallel patterns
- Use case: deep packet inspection (network security), pattern matching in genomics
- Not a general-purpose processor — highly specialized for regular expression / NFA evaluation

**AxDIMM (Samsung / SK Hynix):**
- Processing logic on the DIMM PCB (not inside the DRAM die itself)
- Simpler to manufacture than in-die PIM
- Accelerates DRAM-intensive operations: embedding table lookups (key bottleneck in recommendation systems like Meta's DLRM model)
- Meta uses billions of embedding lookups per day — each lookup is a random memory access with near-zero arithmetic intensity

**The broader PIM landscape (2024):**
- IBM: pioneered analog CIM for DNN inference
- Mythic: analog-digital hybrid AI inference chip (CIM with flash memory)
- Syntiant: ultra-low-power CIM for edge audio ML (on-device wake word)
- Various research projects from MIT, Stanford, ETH Zürich

### Quick Check
> 1. What did Samsung HBM-PIM achieve compared to standard HBM2?
> 2. What is UPMEM and what is its programming model?
> 3. What workload is AxDIMM optimized for and why?

---

## Summary

- **Memory wall**: Processor speed grows 2× faster than memory bandwidth. Memory access costs 2000× more energy than arithmetic. Most AI and big-data workloads are memory-bandwidth-bound.
- **Near-memory computing (NMC)**: Place compute logic in the logic layer of HBM or on-package near DRAM. Access internal TB/s bandwidth instead of CPU-DRAM bus GB/s bandwidth.
- **Processing-in-memory (PIM)**: Compute inside or directly adjacent to DRAM arrays. Digital (Samsung HBM-PIM) or analog crossbar arrays for neural network inference.
- **HBM**: Stacked DRAM with TSVs. 1024-bit interface, 3.35 TB/s for H100. The foundation for most NMC work.
- **Storage-class memory**: Fills the latency gap between DRAM and SSD. Intel Optane used PCM (~300ns). Optane discontinued; MRAM promising for embedded use.
- **CXL**: Coherent interconnect standard. Enables memory expansion, pooling, and disaggregation over PCIe links. Critical for memory-hungry AI serving workloads.
- **Real systems**: Samsung HBM-PIM (2.4× speedup, 70% energy savings), UPMEM (19× DNA alignment speedup), Micron Automata (pattern matching), AxDIMM (embedding lookups).

---

## Exercises

### Easy
1. What is the memory wall and why does it matter for AI workloads?
2. What is arithmetic intensity and how does it determine if a workload is memory-bound?
3. What is CXL and what new memory configurations does it enable?

### Medium
4. Roofline analysis for LLM inference: GPT-3 175B parameter model. One transformer layer has: (a) Weight matrices (175B FP16 params / 96 layers = 1.82B FP16 per layer = 3.64GB data). (b) For one token generation: each weight is read once, each multiply-add uses it once. Arithmetic intensity ≈ 1 FLOP/byte. (c) Roofline on H100: bandwidth 3.35 TB/s, peak FP16 989 TFLOPS. Which bound is tighter? (d) At batch size 1 (single user): performance in tokens/second? (e) At batch size 64: each weight is reused 64 times for 64 tokens simultaneously → arithmetic intensity 64× higher. Now which bound? This is why LLM serving providers maximize batch size.
5. PIM vs standard GPU comparison: A recommendation model does 1 trillion embedding lookups per day. Each lookup: 16 random DRAM reads of 256-byte embeddings. (a) Total bytes read: 1T × 16 × 256 = ? GB. (b) Standard GPU (A100, 2 TB/s): time in seconds? (c) AxDIMM (64× 16GB DIMMs, each at 50 GB/s = 3.2 TB/s total) can access embeddings locally: time? (d) Energy: standard DRAM access ≈ 2 pJ/bit; AxDIMM near-memory access ≈ 0.5 pJ/bit (4× less per bit). Total energy difference for the workload? (e) What limits AxDIMM from being used everywhere?
6. CXL memory tiering: A database server has: 128GB local DDR5 at 70ns latency ($1200), CXL-attached 512GB DRAM at 200ns latency ($1600). Workload: 60% hot data (accessed frequently, latency-sensitive), 40% cold data (accessed rarely). (a) Assign data to memory tier (hot → DDR5, cold → CXL). Does hot data fit in DDR5? (b) Performance impact: if 40% of accesses hit CXL (200ns), overall average latency vs all-DDR5 (70ns). (c) Cost comparison: 640GB all-DDR5 vs 128+512 GB tiered. (d) What hardware/OS mechanism manages tiering automatically? (hint: Intel Memory Mode, Linux demotion, NUMA policies)

### Hard
7. Design a PIM system for transformer attention: Transformer self-attention: for sequence length L=1024, embedding dim d=4096. (a) Memory footprint: Q, K, V matrices each L×d FP16 = 1024×4096×2 = 8MB each. Attention scores: L×L FP16 = 2MB. Total: 26MB per layer. (b) Operations: QKᵀ matmul (L×d × d×L = L² multiply-adds = 1M ops), softmax (L² ops), score×V (L²×d = 4B ops). Total: ~5B FP32 ops. (c) Arithmetic intensity for each step: compute/memory for QKᵀ, softmax, score×V. (d) Which step benefits most from PIM and why? (e) Design the PIM hardware: if HBM bandwidth is 3.35 TB/s and you place 8 SIMD units per vault (16 vaults per stack, 5 stacks = 640 SIMDs), each doing 100 GOPS: total PIM compute = ? TOPS. (f) For batch size 1, compare: standard GPU vs HBM-PIM for attention. For batch size 64?
8. Analog CIM accelerator design: Design an analog crossbar PIM chip for ResNet-50 inference. (a) ResNet-50: 25M parameters, mostly 3×3 convolutions mapped to matrix form (im2col). Matrix sizes range from 9×64 to 2304×512. Total MACs: 4.1 GFLOP. (b) Crossbar limitations: 128×128 tiles (precision: 4-bit weights via PCM cells, 8-bit inputs). How many tiles needed to store all weights? (c) Energy model: digital MAC = 1 pJ; analog crossbar MAC = 0.05 pJ. Total energy digital GPU vs analog CIM for one ResNet-50 inference? (d) Precision impact: ResNet-50 accuracy at 8-bit (INT8) ≈ 76.0% (vs FP32 baseline 76.1%). At 4-bit ≈ 75.2%. At 2-bit ≈ 68.3%. Is 4-bit (crossbar limitation) acceptable? (e) ADC cost: each column needs an ADC (analog-to-digital converter) to convert current to digital number. ADC is 90% of analog CIM chip area. For 128 columns × (number of tiles) ADCs: what limits scalability? (f) Variability: PCM cells drift by ±10% over time. How does this affect inference accuracy and how do you mitigate it?
