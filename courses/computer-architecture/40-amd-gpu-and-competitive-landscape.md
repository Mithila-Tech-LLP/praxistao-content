# Chapter 40: AMD GPU and the Competitive Landscape

NVIDIA dominates the GPU market, but AMD has been the only serious challenger for decades. AMD's GPU division (originally ATI Technologies, acquired by AMD in 2006) has competed with NVIDIA in gaming graphics and increasingly in AI computing. In 2024, AMD's Radeon GPU line holds ~18% of discrete GPU market share in gaming; its Instinct MI300X has become the first credible alternative to NVIDIA H100 for AI inference. This chapter covers AMD GPU architecture, the RDNA gaming line, the CDNA compute line, and Intel's late entry into the discrete GPU market.

## Table of Contents

1. [ATI and the AMD Acquisition](#1-ati-and-the-amd-acquisition)
2. [RDNA — The Radeon Gaming Architecture](#2-rdna--the-radeon-gaming-architecture)
3. [CDNA — The Compute/AI Architecture](#3-cdna--the-computeai-architecture)
4. [MI300X — AMD's AI Challenger](#4-mi300x--amds-ai-challenger)
5. [ROCm — AMD's Software Answer to CUDA](#5-rocm--amds-software-answer-to-cuda)
6. [Intel Arc — The Late Entrant](#6-intel-arc--the-late-entrant)
7. [GPU Market Summary](#7-gpu-market-summary)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. ATI and the AMD Acquisition

ATI Technologies (1985–2006) was NVIDIA's primary rival through the 2000s. The **Radeon** brand launched in 2000, and the ATI-NVIDIA competition drove rapid advancement in graphics hardware. The Radeon 9700 Pro (2002) introduced DirectX 9 hardware and was widely considered the best GPU of its era.

AMD acquired ATI for **$5.4 billion** in 2006, one of the largest semiconductor acquisitions at the time. AMD's motivation:
1. Gain a GPU for AMD's "Fusion" (APU) strategy — CPU + GPU on one die
2. Compete with Intel's CPU+integrated graphics
3. Create a full graphics pipeline: AMD CPU + AMD GPU

The acquisition was financially painful during AMD's lean years (2011–2016) when AMD nearly went bankrupt. The GPU division (renamed AMD Radeon) kept AMD viable during those years with gaming revenue.

Post-acquisition, the Radeon brand continued with:
- HD 4000/5000/6000 series (2008–2010): Competitive with NVIDIA
- Radeon HD 7000 (2012): First GCN (Graphics Core Next) architecture
- RX series (2016+): Modern naming scheme

### Quick Check
> 1. Why did AMD acquire ATI in 2006?
> 2. What is the "Fusion" strategy that motivated the ATI acquisition?
> 3. What role did AMD's GPU business play during AMD's near-bankruptcy period (2011–2016)?

---

## 2. RDNA — The Radeon Gaming Architecture

**GCN (Graphics Core Next, 2012–2019)**: AMD's architecture that ran for 7 years. GCN was designed for compute (OpenCL) as well as graphics. The problem: GCN was not efficient at traditional graphics rasterization workloads compared to NVIDIA's Maxwell/Pascal.

**RDNA 1 (2019)**: Clean break from GCN for gaming-focused performance:
- New "Work Group Processor" instead of GCN's "Compute Unit"
- Dual-compute unit pairs sharing an instruction cache — better shader efficiency
- Increased clock frequencies (1800 MHz vs GCN's 1200 MHz)
- RX 5700 XT: competitive with NVIDIA RTX 2070 at lower price

**RDNA 2 (2020)**: Major leap:
- Infinity Cache: 128MB on-die SRAM cache for the GPU
- Hardware ray tracing (RT) acceleration — AMD's first
- Used in PlayStation 5 and Xbox Series X — giving AMD significant console market share
- RX 6900 XT: competitive with RTX 3090 in many games

```
Infinity Cache:
  Problem: GDDR6 bandwidth requires wide bus (256-bit) which is expensive
  Solution: Put a large on-die SRAM cache; most data hits in cache
            128MB cache + 128-bit bus achieves similar effective bandwidth
            as NVIDIA's 384-bit bus setup, at lower area/power cost
  
  Infinity Cache hit rate: ~60-80% for 1080p/1440p gaming workloads
  Effective bandwidth with cache: ~256 GB/s effective vs 128 GB/s GDDR6 raw
```

**RDNA 3 (2022)**: AMD's first chiplet GPU — the compute dies are separate from the MCD (Memory Cache Die) that holds the Infinity Cache and memory interfaces:
- 5nm compute chiplets + 6nm MCDs
- RX 7900 XTX: competitive with RTX 4080 but behind RTX 4090
- Doubled ray tracing throughput vs RDNA 2

**RDNA 4 (2024)**: RX 9000 series. Improved ray tracing, AI acceleration (new AI compute units), improved power efficiency.

### Quick Check
> 1. What is "Infinity Cache" and how does it allow AMD to use a narrower memory bus?
> 2. Why did AMD's RDNA 2 matter for console gaming?
> 3. What made RDNA 3 a "chiplet GPU"?

---

## 3. CDNA — The Compute/AI Architecture

While RDNA is AMD's gaming line, **CDNA (Compute DNA)** is AMD's separate HPC/AI compute architecture. CDNA is the successor to GCN's compute strengths, but focused exclusively on data center workloads (no display outputs, no ray tracing, no rasterization hardware).

**CDNA 1 (Instinct MI100, 2020):**
- 120 TFLOPS FP16
- 32GB HBM2
- First AMD chip with Matrix Core units (equivalent to NVIDIA Tensor Cores)
- Used in Frontier supercomputer (ORNL) — world's first exascale computer

**CDNA 2 (Instinct MI200 series, 2021):**
- Multi-chip module: two CDNA2 dies connected with Infinity Fabric
- MI250X: 383 TFLOPS FP16, 128GB HBM2e
- Frontier supercomputer (2022): 9,408 nodes × AMD EPYC CPU + 4 MI250X = 1.1 ExaFLOPS peak

**CDNA 3 (Instinct MI300 series, 2023):** AMD's most innovative chip in years:
```
MI300X architecture:
  CPU: None (GPU-only version of MI300)
       (MI300A has CPU+GPU in one package)
  
  GPU compute dies: 8 × CDNA3 XCDs (Compute Die), 5nm
  I/O die: 1 × IOD, 6nm
  Memory: 8 stacks HBM3, 192GB total, 5.3 TB/s bandwidth
  
  The 8 GPU dies are connected via Infinity Fabric and appear as one unified GPU
  192GB unified HBM3 — the most GPU memory of any single card (vs 80GB H100)
```

MI300X advantages over H100:
- **2.4× more memory** (192GB vs 80GB) — critical for large LLM inference
- **1.58× more memory bandwidth** (5.3 vs 3.35 TB/s)

MI300X disadvantages vs H100:
- Less mature ROCm software ecosystem
- NVIDIA NVLink superiority for multi-GPU training
- H100 leads in FLOPS for training workloads

MI300X found early adoption at **Microsoft Azure** (Azure ND MI300X instances) and some AI inference providers. For LLM inference of large models (70B+), MI300X's memory advantage is compelling.

### Quick Check
> 1. What is the difference between RDNA and CDNA architectures?
> 2. What is AMD's MI300X and what is its key advantage over NVIDIA H100?
> 3. Why was the Frontier supercomputer historically significant?

---

## 4. MI300X — AMD's AI Challenger

The MI300X (late 2023) represents AMD's most serious challenge to NVIDIA dominance in AI:

**MI300X key specs:**
- 192GB HBM3 (vs 80GB for H100 SXM5)
- 5.3 TB/s memory bandwidth (vs 3.35 TB/s for H100)
- 1307.4 TFLOPS FP16 (peak, sparse) — higher raw than H100
- 304W TDP in standard mode; 750W when overclocked (OAM)
- 5nm + 6nm multi-die

**Where MI300X wins:**
- LLM inference with large models (Llama-3 70B, Mixtral 8×22B): model fits entirely in one MI300X vs needing 2 H100s
- Memory-bound workloads: more bandwidth = faster weight streaming during inference
- Price/performance for inference: at equivalent memory capacity configurations, MI300X can be cost-competitive

**Where H100 still wins:**
- LLM training: NVLink scale (900 GB/s) vs MI300X's Infinity Fabric between nodes
- Ecosystem: CUDA libraries (cuDNN, NCCL) are more optimized
- Tensor Core precision flexibility (FP8, INT4 support maturity)
- Enterprise support and reliability data from years of datacenter deployment

**Real-world adoption (2024):**
- Microsoft: Azure ND MI300X v5 VMs — 8 MI300X per node, 1.5TB total HBM3
- Meta: Evaluating MI300X for inference
- AMD's AI revenue: grew from ~$400M in 2023 to ~$5B target for 2024

### Quick Check
> 1. For a 70B parameter LLM at FP16 (140GB memory), can it fit on one H100 (80GB)? Can it fit on one MI300X (192GB)?
> 2. What is MI300A and how does it differ from MI300X?
> 3. What is the main barrier to broader MI300X adoption despite its hardware advantages?

---

## 5. ROCm — AMD's Software Answer to CUDA

**ROCm (Radeon Open Compute platform)** is AMD's open-source alternative to CUDA — a complete software stack for GPU computing.

**ROCm components:**
- **HIP (Heterogeneous-computing Interface for Portability)**: CUDA-like programming API. HIP code looks almost identical to CUDA code. Most CUDA code can be converted to HIP with `hipify` tool.
- **rocBLAS**: OpenBLAS-compatible GPU BLAS library (matrix operations)
- **rocDNN**: Convolution and attention primitives (ROCm equivalent of cuDNN)
- **RCCL**: Collective communication library (ROCm equivalent of NCCL)
- **rocRAND, rocSPARSE**: Random number generation, sparse linear algebra

**PyTorch + ROCm**: PyTorch has official ROCm support. `pip install torch --index-url https://download.pytorch.org/whl/rocm6.1` installs a ROCm-enabled PyTorch.

**HIP portability:**
```
CUDA code:          HIP equivalent:
cudaMalloc          hipMalloc
cudaMemcpy          hipMemcpy
__global__          __global__ (same)
threadIdx.x         hipThreadIdx_x (or same in newer HIP)
cuBLAS              rocBLAS (mostly compatible API)
```

**The ROCm gap**: Despite strong effort, ROCm has gaps:
- Some cuDNN ops are not yet in rocDNN (or less optimized)
- Flash Attention (critical for transformer training efficiency) was CUDA-only for years; ROCm support arrived later
- Enterprise support: AMD has fewer GPU datacenter deployments and less experience with cluster-scale debugging
- Triton (the open-source GPU kernel compiler): initially CUDA-focused; ROCm backend added but less tuned

AMD's ROCm investment is growing. In 2024, AMD hired key engineers from NVIDIA and invested heavily in ROCm compatibility, catching up faster than in previous years.

### Quick Check
> 1. What is HIP and what CUDA feature does it correspond to?
> 2. How does the `hipify` tool help developers migrate from CUDA to ROCm?
> 3. What are the two largest practical gaps in ROCm vs CUDA?

---

## 6. Intel Arc — The Late Entrant

Intel made discrete GPUs in the 1990s (Intel740, 1998) and was out of the discrete market for 25 years. In 2022, **Intel Arc** returned:

**Arc Alchemist (2022–2023):**
- Xe-HPG architecture (derived from Intel's Xe-LP iGPU used in Tiger Lake laptops)
- Arc A770: 32 Xe-cores, 16GB GDDR6, 560 GB/s — competitive with RTX 3060 Ti in many workloads
- XeSS (Xe Super Sampling): Intel's version of DLSS, using AI upscaling with Xe Matrix Extensions hardware
- Troubled launch: driver maturity issues, poor DX9/DX11 game compatibility, confusing naming

**Arc Battlemage (2024):**
- Arc B580/B770: second generation, improved driver maturity
- B580 (12GB GDDR6, ~20 TFLOPS) at ~$250: strong value proposition vs comparable NVIDIA/AMD
- Better OpenGL/DX9 compatibility addressed major complaints

**Intel GPU Strategy:**
- Arc (gaming/consumer): compete in the mid-range price tier
- Iris Xe Max (integrated laptop GPU): shipping in many laptops
- Ponte Vecchio / Gaudi (data center): for HPC/AI workloads

**Gaudi** (HPC): Intel's AI accelerator, based on a completely different architecture from Arc. Gaudi 2 and Gaudi 3 target data center inference at lower cost than NVIDIA H100.

**Market reality**: Intel Arc holds a small percentage of discrete GPU market share. NVIDIA holds ~90%, AMD ~10%, Intel ~1-2%. Arc's challenge is overcoming NVIDIA's brand dominance and software ecosystem, plus AMD's gaming reputation.

### Quick Check
> 1. What is XeSS and which hardware unit enables it?
> 2. What was the primary complaint about Arc Alchemist at launch?
> 3. How does Intel Gaudi differ from Intel Arc?

---

## 7. GPU Market Summary

```
GPU Market (Discrete, Gaming) — approximate 2024 market share:
  NVIDIA:  ~90% (RTX 40 series dominates every tier above $300)
  AMD:     ~10% (Radeon RX 7000 series, strong value play)
  Intel:   ~1%  (Arc, growing slowly)
  
GPU Market (Data Center AI) — approximate 2024 market share:
  NVIDIA:  ~80% (H100, A100 installed base)
  AMD:     ~10% (MI300X growing, especially for inference)
  Google:  ~5%  (internal TPU usage, not sold externally in meaningful volume)
  Intel:   ~3%  (Gaudi 3, Ponte Vecchio)
  Other:   ~2%  (AWS Trainium, Cerebras, Groq, etc.)

GPU Market (Mobile, SoC-integrated) — approximate:
  Apple:   ~30% (all Apple devices)
  Qualcomm Adreno: ~30% (Android premium)
  ARM Mali: ~25% (Android mid-range, Samsung Exynos)
  MediaTek: ~10%
  Other:   ~5%
```

**Why NVIDIA's lead is so durable:**
1. CUDA ecosystem (10+ years of software investment)
2. NVLink/NVSwitch (no open standard equivalent)
3. Enterprise relationships and support
4. Process technology access (TSMC's best nodes)
5. Continuous hardware innovation (Tensor Cores, Transformer Engine)
6. Data flywheel: more deployments → more feedback → better drivers/software

**Why AMD can still compete:**
1. Open-source ROCm vs proprietary CUDA (matters for some customers)
2. Price competition on equivalent hardware tiers
3. Console gaming monopoly (PS5, Xbox Series X) keeps Radeon relevant
4. MI300X memory advantage for specific inference workloads
5. AMD's overall CPU+GPU integrated strategy (EPYC + Instinct in HPC)

### Quick Check
> 1. What percentage of the discrete GPU gaming market does NVIDIA hold?
> 2. Name three reasons NVIDIA's lead in AI GPUs is durable.
> 3. Why does AMD's console gaming position (PS5, Xbox) matter for its GPU business?

---

## Summary

- **AMD GPU** history: acquired from ATI (2006), evolved through GCN → RDNA (gaming) → CDNA (compute/AI).
- **RDNA architecture**: gaming-focused, featuring Infinity Cache (large on-die SRAM to reduce memory bandwidth needs), hardware ray tracing (RDNA 2+), and chiplet design (RDNA 3).
- **CDNA architecture**: compute/AI focused. MI100 (exascale pioneer), MI200 (Frontier supercomputer), MI300X (192GB HBM3, flagship AI challenger).
- **MI300X**: AMD's most credible H100 alternative — 2.4× more memory, 1.58× more bandwidth. Winning inference workloads. Still behind in multi-GPU training ecosystem.
- **ROCm**: AMD's open CUDA alternative. HIP for CUDA portability. Growing but still behind in library optimization and ecosystem maturity.
- **Intel Arc**: late entrant to discrete GPU market. Small share but growing with driver improvements. Gaudi for data center AI.
- NVIDIA holds ~80–90% of AI compute GPU market due to hardware + software moat.

---

## Exercises

### Easy
1. What is Infinity Cache and what problem does it solve for AMD's RDNA GPUs?
2. What is the difference between AMD RDNA (Radeon) and AMD CDNA (Instinct) architectures?
3. What is HIP and how does it relate to CUDA?

### Medium
4. AMD MI300X (192GB, 5.3 TB/s) vs NVIDIA H100 (80GB, 3.35 TB/s) for LLM inference: (a) Llama-3 70B at FP16 requires 140GB of memory for weights. Can it fit on each card? (b) If memory bandwidth determines tokens/second for a batch-size=1 inference, calculate tokens/second for each GPU assuming 2 FLOPs per weight access per token. (c) For a batch-size=32 inference (compute-bound), H100 has 989 TFLOPS FP16 vs MI300X's ~1307 TFLOPS. Which wins?
5. Intel Arc B580 launched at $249 with 12GB GDDR6 and ~20 TFLOPS. NVIDIA RTX 4060 is $299 with 8GB GDDR6 and ~15 TFLOPS. Pure hardware numbers favor Arc — yet NVIDIA outsells Arc significantly. List and analyze 5 non-hardware factors that explain this market outcome.
6. RDNA 2's Infinity Cache design: the cache has 128MB of SRAM at 1920 GB/s internal bandwidth (4× the GDDR6 bus). If gaming has a 70% cache hit rate with 128-bit GDDR6: (a) what is the effective bandwidth for cache hits (reads from Infinity Cache)? (b) for cache misses (from GDDR6)? (c) what is the weighted average effective bandwidth? (d) compare to NVIDIA RTX 3080's 320-bit GDDR6X at 760 GB/s (no Infinity Cache equivalent).

### Hard
7. ROCm adoption decision for a small AI startup: You're training a custom LLM. You have $500K budget for GPUs. H100 SXM is $30K, MI300X is $20K. Training requires 100 PFLOPS for 2 weeks. Analyze: (a) How many of each GPU do you need to meet the FLOPS target? (b) What is the hardware cost for each option? (c) Estimate engineering cost to port your PyTorch codebase from CUDA to ROCm (assume 2 engineers × 6 weeks × $10K/week loaded cost). (d) Factor in that ROCm Flash Attention is 80% efficiency vs CUDA (10% slower training). What is your total-cost-of-ownership decision?
8. The GPU memory hierarchy: explain why the combination of Infinity Cache (GPU L3 cache equivalent) + GDDR (main VRAM) + shared memory (per-SM) creates a multi-level hierarchy similar to CPU caches. For a shader computing 1000 iterations of a loop over a 4KB buffer: (a) where would this buffer ideally reside? (b) what happens if 32 threads in a warp each need a different 4KB buffer? (c) how does texture cache differ from shared memory for read-only data?
