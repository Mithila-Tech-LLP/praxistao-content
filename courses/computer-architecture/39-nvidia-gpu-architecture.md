# Chapter 39: NVIDIA GPU Architecture — From GeForce to Hopper

NVIDIA has dominated the GPU market for nearly three decades. What began as a graphics card company is now the most valuable chip company in the world — at its 2024 peak, NVIDIA's market capitalization exceeded $3 trillion, making it the third-largest company by market cap globally. The reason: NVIDIA bet on GPU computing for AI a decade before AI became the defining technology trend, building not just better hardware but a software ecosystem (CUDA) that made switching to competitors extraordinarily painful. This chapter traces NVIDIA's GPU architecture from the original GeForce through Volta's AI pivot to the Blackwell generation.

## Table of Contents

1. [NVIDIA's History and the GPU Era](#1-nvidias-history-and-the-gpu-era)
2. [The Tesla/Fermi Era — General-Purpose GPU Computing](#2-the-teslafermi-era--general-purpose-gpu-computing)
3. [Maxwell, Pascal, and Turing — Modern Foundations](#3-maxwell-pascal-and-turing--modern-foundations)
4. [Ampere — The AI Datacenter Chip](#4-ampere--the-ai-datacenter-chip)
5. [Hopper — The H100 Revolution](#5-hopper--the-h100-revolution)
6. [Blackwell and Beyond](#6-blackwell-and-beyond)
7. [NVLink and Multi-GPU Systems](#7-nvlink-and-multi-gpu-systems)
8. [CUDA Software Ecosystem](#8-cuda-software-ecosystem)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. NVIDIA's History and the GPU Era

NVIDIA was founded in 1993 by Jensen Huang, Chris Malachowsky, and Curtis Priem. Its first major product was the NV1 (1995) — a graphics/sound/game chip that failed commercially. The **RIVA 128 (1997)** was the first hit: the first single-chip 3D accelerator with hardware transform and lighting, at 3 million polygons per second.

**GeForce 256 (1999)**: NVIDIA coined the term "GPU" (Graphics Processing Unit) to differentiate its chip — which handled geometry transformation and lighting in hardware — from earlier "graphics accelerators" that only rasterized.

**The key decisions that defined NVIDIA:**

**Decision 1 (2006): CUDA**. Jensen Huang decided to invest in making NVIDIA's GPU programmable for general computing. CUDA launched with the GeForce 8800. Most analysts at the time considered this a risky distraction from graphics. In retrospect, it was the most consequential software investment in chip history.

**Decision 2 (2012): Bet on deep learning**. After seeing Alex Krizhevsky train AlexNet on two GTX 580 GPUs (winning ImageNet by a huge margin), NVIDIA pivoted resources toward AI/deep learning workloads. The Tesla K40 accelerator was the first product in this new direction.

**Decision 3 (2016): Tensor Cores**. NVIDIA invented specialized matrix-multiply hardware (Tensor Cores) for neural networks in the Volta architecture. This was hardware differentiation that competitors haven't fully replicated.

NVIDIA is fabless — manufacturing is entirely at TSMC and Samsung.

### Quick Check
> 1. What did NVIDIA mean when it coined the term "GPU" for the GeForce 256?
> 2. What was CUDA and why was it a risky decision in 2006?
> 3. When did NVIDIA first bet on deep learning as a major use case?

---

## 2. The Tesla/Fermi Era — General-Purpose GPU Computing

**Tesla (2007)**: Not the same as Tesla cars — NVIDIA's compute GPU brand. The Tesla architecture introduced:
- 128 scalar shader processors per SM (vs earlier vector units)
- Shared memory (programmer-controlled cache within an SM)
- CUDA support (first fully CUDA-capable architecture)

**Fermi (2010)**: Major architecture overhaul that made GPUs suitable for double-precision (FP64) scientific computing:
- 512 CUDA cores per die
- ECC memory (error-correcting code — essential for scientific reliability)
- C++ exceptions and function pointers in CUDA kernels
- L2 cache shared across all SMs

Fermi demonstrated that GPUs could replace traditional scientific workstations for simulation. National laboratories (ORNL, LLNL) began deploying GPU-accelerated supercomputers.

**Kepler (2012)**: Power efficiency focus. Dynamic Parallelism (GPU kernels can launch other GPU kernels). Used in the Titan supercomputer at ORNL — the fastest computer in the world in 2012.

### Quick Check
> 1. What is "shared memory" in NVIDIA GPUs and why is it programmer-controlled rather than hardware-managed?
> 2. What is ECC memory and why is it important for scientific computing?
> 3. What is "Dynamic Parallelism" in Kepler?

---

## 3. Maxwell, Pascal, and Turing — Modern Foundations

**Maxwell (2014)**: Dramatic power efficiency improvement — 2× performance per watt over Kepler. Key architectural change: rebalanced CUDA core to control logic ratio. GTX 980 was the most efficient GPU of its era.

**Pascal (2016)**: A generation that defined the deep learning training era before Tensor Cores:
- 16nm FinFET (TSMC) — major process improvement
- NVLink 1.0: first GPU-to-GPU high-bandwidth interconnect (80 GB/s bidirectional)
- HBM2 memory on Tesla P100 (the first production GPU with HBM2)
- Pascal P100: 10.6 TFLOPS FP16 — the workhorse of early deep learning training (AlphaGo was trained on P100s)
- GP102: consumer GTX 1080/1080 Ti — for years the benchmark for gaming performance

**Volta (2017)**: The AI pivot:
- **Tensor Cores**: 4×4 matrix multiply units. 120 TFLOPS FP16 (Tensor Core) vs 14 TFLOPS FP32 (CUDA cores)
- NVLink 2.0 (300 GB/s bidirectional)
- V100: 5120 CUDA cores, 640 Tensor Cores, 32GB HBM2
- No consumer GPU — Volta was exclusively for AI/HPC markets

**Turing (2018)**: Brought Tensor Cores and RT Cores to consumers:
- **RT Cores**: dedicated ray-triangle intersection hardware
- **DLSS (Deep Learning Super Sampling)**: use AI (running on Tensor Cores) to upscale lower-resolution frames to 4K — trading GPU compute for real-time quality
- GeForce RTX 2080 Ti: 4352 CUDA cores, 544 Tensor Cores, 68 RT Cores
- Turing was the first consumer GPU with hardware ray tracing

### Quick Check
> 1. What are Tensor Cores and in which NVIDIA architecture did they debut?
> 2. What is DLSS and how do Tensor Cores enable it?
> 3. Why was Pascal P100 significant for early deep learning training?

---

## 4. Ampere — The AI Datacenter Chip

**Ampere (2020)**: NVIDIA's biggest architecture jump in years, driven by the explosive demand for AI training:

**A100 (Data Center):**
- 6912 CUDA cores (FP32)
- 432 Tensor Cores (3rd gen)
  - 312 TFLOPS FP16 (Tensor Core)
  - 77.6 TFLOPS FP64 (Tensor Core, 2× vs V100)
  - **19.5 TFLOPS FP32** (CUDA cores, conventional)
- 80GB HBM2e, 2 TB/s bandwidth
- NVLink 3.0 (600 GB/s bidirectional)
- Multi-Instance GPU (MIG): partition one A100 into up to 7 independent GPU instances for smaller workloads
- 7nm Samsung process (then migrated to TSMC 7nm for A100 SXM4)

**GA102 (Consumer — RTX 3090/3090 Ti):**
- 10496 CUDA cores
- 328 Tensor Cores (3rd gen)
- 24GB GDDR6X, 936 GB/s
- 4th gen NVLink (limited use in consumer)

**The Tensor Core generational evolution:**
```
Volta Tensor Core (1st gen): FP16 × FP16 → FP32 accumulation, 4×4 matrix
Turing Tensor Core (2nd gen): Added INT8/INT4 support
Ampere Tensor Core (3rd gen): Added TF32 (19-bit), BF16, INT8, added FP64
Hopper Tensor Core (4th gen): Added FP8 (8-bit float), transformer engine
```

**TF32 (TensorFloat-32)**: NVIDIA's own numerical format — 19 bits total (1 sign, 8 exponent, 10 mantissa). Gives the dynamic range of FP32 with the storage of ~FP16, and uses Tensor Core acceleration. Drop-in replacement for FP32 in AI training with ~10× speedup and minimal accuracy loss. Used transparently by PyTorch/TensorFlow when targeting A100.

**Multi-Instance GPU (MIG)**: One A100 can be partitioned into 7 separate GPU instances, each with its own:
- CUDA cores (subset)
- Tensor cores (subset)
- HBM (1/7 of total)
- PCIe/NVLink bandwidth (partitioned)

Enabling cloud providers to rent A100s at 1/7, 1/4, 1/2, and full capacity. AWS, Azure, and GCP all offer MIG-based instances.

### Quick Check
> 1. What is MIG (Multi-Instance GPU) and why is it important for cloud providers?
> 2. What is TF32 and why did NVIDIA invent it?
> 3. A100 achieves 312 TFLOPS in FP16 Tensor Core mode but only 19.5 TFLOPS in FP32 CUDA core mode. What makes Tensor Cores so much faster for this operation?

---

## 5. Hopper — The H100 Revolution

**Hopper (2022)**: H100 — the most commercially important chip in history for AI, selling for $30,000–$40,000 each, with demand far exceeding supply in 2023–2024.

**H100 SXM Specifications:**
- **4th gen Tensor Cores**: 989 TFLOPS FP16 (sparse), 495 TFLOPS FP16 (dense)
- **FP8 support**: ~2000 TFLOPS effective for inference
- 80GB HBM3, **3.35 TB/s memory bandwidth**
- 18432 CUDA cores
- 132 SMs, each with 4 Tensor Cores
- NVLink 4.0: 900 GB/s bidirectional per GPU
- PCIe 5.0 host interface
- TSMC 4nm process

**The Transformer Engine**: a new hardware feature in Hopper. LLMs use transformer architecture. The Transformer Engine automatically selects FP8 or FP16 for each layer of a transformer based on numerical stability requirements, maximizing Tensor Core throughput while preserving accuracy.

```
GPT-4 (estimated 1.8T parameters) training:
  Compute: ~2 × 10^25 FLOPs (rough estimate)
  H100 cluster: ~25,000 H100s (OpenAI's reported training cluster)
  At 495 TFLOPS per H100: ~2.5 months of continuous training
  (actual training time varies with batch size, data passes, restarts)
```

**DGX H100 system:**
- 8 × H100 SXM GPUs
- 640GB total HBM3
- All-to-all NVLink (900 GB/s between any pair of GPUs)
- 2 TB/s NVSwitches (connects all 8 GPUs as one coherent GPU memory pool)

**The supply chain crisis**: In 2023, H100 supply was severely constrained. NVIDIA was TSMC's biggest customer for N4 (4nm) capacity. HBM3 supply was limited to Samsung and SK Hynix. The result: H100 secondary market prices reached $40,000–$70,000. Cloud providers pre-purchased entire allocations. NVIDIA's revenue grew 265% YoY in FY2024.

### Quick Check
> 1. What is the Transformer Engine in Hopper and what workload is it designed for?
> 2. H100 has 3.35 TB/s memory bandwidth. Why is this important for LLM inference?
> 3. What is an NVSwitch and what does it enable for multi-GPU systems?

---

## 6. Blackwell and Beyond

**Blackwell (2024)**: NVIDIA's next generation, announced March 2024.

**B100/B200/GB200:**
- 208B transistors (vs H100's 80B) — larger die on TSMC 4NP
- 4.5× faster FP4 inference vs H100
- FP8 training: 20 PFLOPS per GPU (20× over H100)
- **HBM3E**: B200 has 192GB HBM3e at 8 TB/s (vs 80GB/3.35 TB/s for H100)
- **NVLink-C2C**: connects two Blackwell dies in the GB200 (Grace Blackwell Superchip) — CPU + 2 GPUs in one package
- GB200 NVL72: rack-scale computing — 36 GB200 superchips = 72 B200 GPUs + 36 Grace CPUs

**Grace CPU**: NVIDIA's own ARM CPU for HPC/AI, based on Neoverse cores. Designed as the "host" CPU paired with Blackwell in GB200 systems.

**NVIDIA's long-term strategy**: Move beyond selling GPUs to selling complete AI systems — full racks, networking (Quantum InfiniBand), software (CUDA, cuDNN, NIM inference microservices, DGX Cloud). Transform from a chip company to an AI computing infrastructure company.

### Quick Check
> 1. What is the GB200 "Grace Blackwell Superchip"?
> 2. Blackwell B200 has 192GB HBM3e. If a model requires 200GB, can it fit? What option does NVIDIA offer?
> 3. Why is NVIDIA adding ARM CPUs (Grace) to its GPU systems?

---

## 7. NVLink and Multi-GPU Systems

When a single GPU isn't enough — for large AI models that don't fit in one GPU's memory, or for training that needs more throughput — multiple GPUs must cooperate. PCIe alone isn't fast enough:

```
PCIe 4.0 x16: 32 GB/s bidirectional (total)
NVLink 4.0:   900 GB/s bidirectional (per GPU-to-GPU link!)
              28× faster than PCIe
```

**NVLink** is NVIDIA's proprietary high-speed GPU-to-GPU interconnect. It uses dedicated metal-to-metal connections on the GPU package or through specially designed "NVLink bridges" on consumer GPUs.

**NVSwitch**: An NVIDIA chip that acts as a crossbar switch for NVLink — allows any GPU to communicate with any other GPU at full NVLink bandwidth simultaneously. The DGX H100 uses NVSwitches to create a fully non-blocking all-to-all interconnect between 8 GPUs.

**Tensor Parallelism**: A large transformer model is split across multiple GPUs — each GPU holds part of the weight matrices. During forward pass, intermediate results are all-reduced across GPUs via NVLink. This requires enormous inter-GPU bandwidth — only practical with NVLink, not PCIe.

**NVLink vs Infinity Fabric vs CXL**: Each interconnect serves a different purpose:
- NVLink: NVIDIA GPU-to-GPU, proprietary
- AMD Infinity Fabric: AMD GPU/CPU, proprietary
- CXL (Compute Express Link): Open standard, PCIe-based, lower bandwidth than NVLink but enables CPU-GPU memory coherence

### Quick Check
> 1. Why is NVLink so much faster than PCIe for GPU-to-GPU communication?
> 2. What is tensor parallelism and why does it require NVLink speeds?
> 3. What is an NVSwitch and why is it needed in an 8-GPU DGX system?

---

## 8. CUDA Software Ecosystem

NVIDIA's hardware dominance is reinforced by software moat — switching to AMD or Intel GPUs isn't just a hardware swap:

**CUDA Core:**
- nvcc: CUDA C++ compiler
- CUDA runtime + driver API
- Profiling: Nsight Systems (system-level), Nsight Compute (kernel-level)

**cuDNN (CUDA Deep Neural Network library)**: Hand-optimized primitives for conv, attention, normalization. Used by every major framework. ~10× faster than naive CUDA for convolutions.

**cuBLAS**: Optimized BLAS routines (matrix multiply, GEMM). The backbone of all neural network training.

**NCCL (NVIDIA Collective Communication Library)**: Optimized all-reduce, broadcast, gather operations over NVLink and InfiniBand for multi-GPU/multi-node training.

**TensorRT**: NVIDIA's inference optimization engine. Takes a trained model, optimizes operator fusion, precision reduction, and kernel selection for maximum inference throughput. 2–10× faster than raw PyTorch inference.

**AI Frameworks:**
```
PyTorch → torch.compile → inductor → triton → CUDA kernels → H100
TensorFlow → XLA → CUDA/cuDNN → H100
JAX → XLA → CUDA → H100
```

**The switching cost**: AMD's ROCm (their CUDA equivalent) supports most of the same APIs. But:
- cuDNN has years of hand-tuning for NVIDIA hardware; rocDNN is less optimized
- Some closed-source models and inference frameworks are CUDA-only
- Enterprise AI teams have CUDA expertise; retraining for ROCm is a cost
- NVIDIA has exclusive features (Transformer Engine, NVLink scale) not available anywhere else

### Quick Check
> 1. What is cuDNN and why does it matter more than raw FLOPS for neural network performance?
> 2. What is TensorRT and when would you use it instead of PyTorch?
> 3. What makes it difficult for AI researchers to switch from NVIDIA to AMD GPUs?

---

## Summary

- NVIDIA's GPU evolution: GeForce (1999, gaming) → CUDA (2006, general compute) → Tesla/Fermi/Kepler (scientific) → Volta (AI with Tensor Cores) → Ampere (A100, AI training workhorse) → Hopper (H100, LLM era) → Blackwell (B200, hyperscale AI).
- **Tensor Cores** (debuted Volta, 2017) provided 8–16× speedup for matrix multiply — the fundamental AI operation.
- **Hopper's H100** is the dominant LLM training GPU: 989 TFLOPS FP16, 3.35 TB/s HBM3, NVLink 4.0.
- **NVLink** (up to 900 GB/s, H100) enables multi-GPU systems with full bandwidth — essential for tensor-parallel LLM training.
- **CUDA** ecosystem (cuDNN, TensorRT, NCCL) creates a software moat: hardware competitors face both a hardware and a software switching cost.
- NVIDIA is transitioning from chip company to AI infrastructure platform.

---

## Exercises

### Easy
1. What are Tensor Cores and why are they more efficient than CUDA cores for matrix multiplication?
2. List three NVIDIA-specific libraries (not hardware) that contribute to NVIDIA's ecosystem advantage.
3. What is DLSS and how does it use GPU hardware for a graphics application?

### Medium
4. An A100 has 19.5 TFLOPS FP32 (CUDA cores) and 312 TFLOPS FP16 (Tensor Cores). A ResNet-50 training step requires 8.3 × 10^9 floating-point operations (FP32). (a) How long does one training step take on A100 FP32 CUDA cores? (b) If converted to FP16 with mixed precision, and assuming the workload is 80% matrix multiply (Tensor Core) and 20% element-wise (CUDA core), what is the effective throughput? (c) What is the speedup?
5. H100 memory bandwidth is 3.35 TB/s. A Llama-2 70B model has 70 billion FP16 parameters (140GB). For token generation (one forward pass per token), each token requires reading all weights once. (a) How long does one forward pass take if bandwidth-limited? (b) How many tokens per second can one H100 generate? (c) How does batching (generating for multiple users simultaneously) improve GPU utilization?
6. Design a hypothetical training cluster for a 100B parameter model (200GB at FP16). You need to train with 8-way tensor parallelism (each GPU holds 1/8 of the model). (a) What is the minimum VRAM per GPU needed? (b) What NVLink bandwidth is needed if each forward pass exchanges 200GB of activations between GPU pairs? (c) How many H100s would you need to achieve 10 PFLOPS training throughput?

### Hard
7. NVIDIA's CUDA moat vs hardware performance: AMD MI300X has 192GB HBM3 and 5.3 TB/s bandwidth vs H100's 80GB/3.35 TB/s — clear hardware advantage for large models. Yet NVIDIA holds ~90% market share in AI training. Model the total cost of switching for an AI startup: (a) engineer-months to port CUDA kernels to HIP/ROCm, (b) performance overhead from using non-optimized libraries, (c) risk of missing features (NVLink scale, Transformer Engine), (d) support and ecosystem network effects. When does the AMD hardware advantage outweigh the switching cost?
8. FP8 vs FP16 vs FP32 in neural network training: analyze the numerical trade-offs: (a) draw the floating point number lines for FP32, FP16, FP8 (E4M3), showing exponent range and mantissa precision, (b) for gradient descent, which quantities are most sensitive to numerical precision (gradients vs weights vs activations), (c) explain loss scaling for FP16 training — why is it needed and how does it work, (d) for FP8 training (H100 Transformer Engine), what additional techniques are needed beyond loss scaling?
