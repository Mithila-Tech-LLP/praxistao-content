# Chapter 41: NPUs and AI Accelerators

The GPU was not designed for AI — it happened to be good at it because neural networks are fundamentally matrix multiplications, and GPUs are good at parallel arithmetic. But a chip optimized specifically for AI workloads — with custom data flows, memory hierarchies, and numerical formats tailored to neural networks — can outperform a GPU on AI tasks while consuming a fraction of the power. These chips are called **Neural Processing Units (NPUs)** or AI accelerators. Google's TPU, Cerebras, Graphcore, Groq, and the NPUs inside every smartphone are purpose-built for neural networks. This chapter explains how they work and why they matter.

## Table of Contents

1. [Why Dedicated AI Hardware?](#1-why-dedicated-ai-hardware)
2. [The Systolic Array — Heart of the TPU](#2-the-systolic-array--heart-of-the-tpu)
3. [Google TPU — The Original AI ASIC](#3-google-tpu--the-original-ai-asic)
4. [Smartphone NPUs](#4-smartphone-npus)
5. [Cerebras — The Wafer-Scale Engine](#5-cerebras--the-wafer-scale-engine)
6. [Groq — Deterministic LPU](#6-groq--deterministic-lpu)
7. [Graphcore — Intelligence Processing Unit](#7-graphcore--intelligence-processing-unit)
8. [AWS Trainium and Inferentia](#8-aws-trainium-and-inferentia)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Dedicated AI Hardware?

A GPU is a general parallel processor optimized for throughput. It must handle:
- Variable-length SIMT kernels
- Complex memory coherence
- Graphics-specific hardware (ray tracing, TMUs, ROPs)
- Flexible compute for diverse workloads

This generality costs power and silicon area. A chip designed only for neural network inference or training can be:
1. **More efficient**: Remove all non-AI hardware. No rasterizer, no ray tracing, no general cache coherence.
2. **More memory-bandwidth-efficient**: Design the data flow to maximize reuse of weights (reducing off-chip memory reads)
3. **More numerically optimized**: Implement INT8, INT4, FP8 natively with minimum overhead
4. **More power-efficient**: No power spent on unused features

The key operation in neural networks is **GEMM (General Matrix Multiply)**:
```
Y = W × X + b
where W = weight matrix, X = input (activation) matrix, b = bias

For a transformer with d_model=4096, the attention projection matrix is [4096 × 4096]
One forward pass through one transformer block: ~2 × 4096² × seq_len FLOPs
```

If you can optimize hardware specifically to do GEMM fast, you win.

### Quick Check
> 1. Name three things a GPU has that a dedicated AI accelerator doesn't need.
> 2. What is GEMM and why is it the dominant operation in neural networks?
> 3. What is the power efficiency advantage of specialized hardware vs general-purpose hardware?

---

## 2. The Systolic Array — Heart of the TPU

A **systolic array** is a hardware architecture for matrix multiplication that minimizes off-chip memory access by routing data through a grid of processing elements (PEs) in a wave-like pattern.

**Conceptual example**: Multiplying matrix A (M×K) by matrix B (K×N):

```
Systolic array (4×4 example):
     B[0][0]  B[0][1]  B[0][2]  B[0][3]   ← B values flow right
       ↓        ↓        ↓        ↓
A[0] →[PE] → [PE] → [PE] → [PE]          ← A values flow down
A[1] →[PE] → [PE] → [PE] → [PE]
A[2] →[PE] → [PE] → [PE] → [PE]
A[3] →[PE] → [PE] → [PE] → [PE]
       ↓        ↓        ↓        ↓
     C[0][0]  C[0][1]  C[0][2]  C[0][3]   ← output C accumulates

Each PE: multiply A[i] × B[j], add to accumulated result
```

**Why systolic arrays are efficient:**
- A values are loaded once and flow through all PEs in their row → each value does K multiply-adds before being discarded
- B values flow through all PEs in their column
- No PE needs to access memory during computation — all data flows from PE to PE
- The entire matrix multiply is done with 2 memory reads (A and B matrices) and 1 memory write (C matrix)

**Compare to CPU**: A CPU matrix multiply reads A and B values multiple times (from cache or memory) to compute each output element. A systolic array reads each input value once.

For a 256×256 systolic array: 65,536 PEs, each doing 1 MAC (multiply-accumulate) per cycle at 700 MHz = 46 TOPS — in a tiny die area, at low power.

### Quick Check
> 1. In a systolic array, how many times is each value of matrix A read from memory?
> 2. Why is the systolic array more memory-efficient than a CPU matrix multiply?
> 3. If a 128×128 systolic array runs at 1 GHz, what is its peak throughput in GOPS (giga-operations per second)?

---

## 3. Google TPU — The Original AI ASIC

**Google TPU v1 (2015, deployed 2016)**: Google's first custom AI accelerator, designed in just 15 months after Google realized their data centers would need 2× the number of servers to handle inference demand if using CPUs.

**TPU v1 architecture:**
- Designed specifically for **inference** (not training)
- 256×256 systolic array (65,536 multipliers)
- 8-bit integer arithmetic (INT8) — sufficient for inference, much cheaper than FP32
- 92 TOPS at 700 MHz
- 8GB on-chip SRAM (compared to off-chip DRAM on CPUs)
- 28W power consumption (vs 200W for the equivalent CPU setup)
- Connected to CPUs via PCIe — acts as a co-processor

**Key insight**: For inference with already-trained models, INT8 is almost as accurate as FP32 but 4× cheaper per operation and 4× cheaper per byte of memory. TPU v1 was designed around this.

**Results**: Google reported that TPU v1 achieved 15–30× better performance per watt than CPUs/GPUs for AlexNet and Inception inference.

**TPU v2 (2017)**: Added training support, floating-point cores, HBM memory, custom 128×128 systolic arrays.

**TPU v3 (2018)**: Liquid-cooled, 420 TFLOPS BF16 per chip.

**TPU v4 (2021)**: 275 TFLOPS BF16 (dense), OCS (optical circuit switches) for datacenter networking, up to 4096 chips in one "pod."

**TPU v5 (2024)**: Unknown full specs as of writing, but Google uses TPUs internally for training Gemini and serving all Google AI products.

```
Google's vertical integration advantage:
  Google designs: TPU hardware + TensorFlow/JAX frameworks + XLA compiler
  The compiler knows exactly what hardware capabilities exist
  The hardware is designed knowing what the compiler will generate
  Result: hardware and software are co-optimized from day 1
```

Google does not sell TPUs; they are only available via Google Cloud (TPU VM instances).

### Quick Check
> 1. What problem was TPU v1 designed to solve?
> 2. Why does TPU v1 use INT8 instead of FP32?
> 3. What is Google's vertical integration advantage in the TPU program?

---

## 4. Smartphone NPUs

Every flagship smartphone since 2018 has an NPU. The Apple Neural Engine (A12, 2018) started the trend; every other chip maker followed.

**Purpose**: Run neural networks locally (on-device) for:
- Face recognition (Face ID)
- Object detection in camera viewfinder
- Real-time portrait blur
- Noise cancellation in calls
- "Hey Siri" / "Hey Google" wake word detection (runs 24/7 at µW power)
- OCR, handwriting recognition
- On-device LLM inference (2024+)

**Architecture**: Most smartphone NPUs are custom systolic arrays or vector units:

| NPU | Device | Architecture | TOPS | Process |
|-----|--------|-------------|------|---------|
| Apple Neural Engine (A16) | iPhone 14 Pro | Custom | 17 TOPS | TSMC 4nm |
| Apple Neural Engine (A17 Pro) | iPhone 15 Pro | Custom | 35 TOPS | TSMC 3nm |
| Apple Neural Engine (M4) | MacBook | Custom | 38 TOPS | TSMC 3nm |
| Qualcomm Hexagon NPU (Snap 8 Gen 3) | Android flagship | Custom | 45 TOPS | TSMC 4nm |
| Samsung NPU (Exynos 2400) | Galaxy S24 (EU) | Custom | 14.7 TOPS | Samsung 4nm |
| Google Tensor G3 NPU | Pixel 8 | Custom | ~15 TOPS | Samsung 4nm |
| MediaTek APU 690 (Dimensity 9300) | Mid-high Android | Custom | 33 TOPS | TSMC 4nm |

**Always-on processing**: A separate, tiny micro-NPU often runs continuously at µW power levels (from a µA-scale power domain) to listen for wake words, track motion, or monitor sensor data without waking the main SoC. This is why "Hey Siri" works with the iPhone screen off and doesn't drain the battery noticeably.

**On-device LLM**: 2024 saw the first smartphones running small LLMs locally (2–4 billion parameters). Google Pixel 8 runs a pruned version of Gemini Nano. Apple Intelligence on A18 Pro runs ~3B parameter models. The NPU makes this power-feasible (inference at ~2–5W instead of 20–40W on a CPU).

### Quick Check
> 1. What does "always-on" NPU processing enable and how is it power-feasible?
> 2. Apple Neural Engine went from 5 TOPS (A12, 2018) to 38 TOPS (M4, 2024) — a 7.6× increase. What drives this improvement?
> 3. What is the smallest LLM that can run on a smartphone NPU and maintain usable quality?

---

## 5. Cerebras — The Wafer-Scale Engine

**Cerebras Systems** took a radical approach: instead of making a chip and connecting many chips together, make one chip the entire wafer.

**Cerebras Wafer-Scale Engine 2 (WSE-2, 2021):**
- Die size: 46,225 mm² (a standard wafer! vs NVIDIA A100's 826 mm² or H100's 814 mm²)
- 2.6 trillion transistors (vs 80B for H100)
- 850,000 AI-optimized cores
- 40GB on-die SRAM (vs H100's 80GB HBM — but on-die SRAM is 100× faster than HBM)
- 20 PB/s on-chip memory bandwidth (vs 3.35 TB/s for H100)
- TSMC 7nm
- Power: 23 kW per system (requires proprietary water cooling)

**The key insight**: For transformers, the bottleneck is often memory bandwidth (loading weight matrices). If the entire model fits in on-chip SRAM, the model can run without any off-chip memory accesses, eliminating the memory bottleneck entirely.

**Cerebras CS-3 (2024, WSE-3):**
- 4 trillion transistors
- TSMC 5nm
- 900,000 cores
- 44GB on-die SRAM
- 100 PB/s on-chip bandwidth

**Challenge**: Wafer-scale manufacturing has very low yield — any defect anywhere on the wafer kills a small-chip die, but a wafer-scale chip must use redundancy (spare cores to replace defective ones). Cerebras uses fault-tolerant routing and spare core allocation.

**Use case**: Cerebras systems are used by research labs (Argonne National Laboratory, National Energy Research Scientific Computing Center) for training medium-sized models (1-10B parameters) where the entire model fits in 44GB and training speed benefits from the massive on-chip bandwidth.

### Quick Check
> 1. What makes the Cerebras Wafer-Scale Engine unusual compared to all other chips?
> 2. What is the performance advantage of 44GB of on-chip SRAM over 80GB of HBM (off-chip)?
> 3. What manufacturing challenge does wafer-scale integration face?

---

## 6. Groq — Deterministic LPU

**Groq** (not related to xAI's Grok) designed a completely different approach: the **LPU (Language Processing Unit)**, focused on extremely low-latency, high-throughput LLM inference.

**Key Groq architectural insight**: A GPU's performance for LLM inference varies unpredictably — different batch sizes, sequence lengths, and memory access patterns cause variable performance. Groq's LPU is **deterministic**: the hardware executes operations in a fixed, predictable schedule. This enables:
- Predictable, consistent inference latency
- Extreme compiler optimization (the compiler knows exactly when each operation will execute)
- No hardware speculation, no out-of-order execution, no branch prediction

**GroqChip (TSP — Tensor Streaming Processor):**
- 220 MB SRAM per chip (enormous on-chip memory)
- 1 PB/s on-chip bandwidth per chip
- Functional units arranged in streams; each stream processes data in deterministic order
- No caches, no MMU, no complex hardware (simpler = faster)

**Groq LPU system**: Multiple TSP chips arranged in a rack, connected by high-bandwidth links. Inference is distributed across chips with model parallelism.

**Performance claim**: Groq demonstrates ~800 tokens/second throughput on Mixtral-8x7B — significantly faster than NVIDIA H100 for inference. This matters for real-time applications (chatbots, code completion) where latency is visible to users.

**Limitation**: Groq is inference-only; not suitable for training (no floating-point backward pass, no gradient accumulation). Also limited to models that fit in their memory configuration.

### Quick Check
> 1. What does "deterministic" mean in the context of Groq's LPU?
> 2. Why does Groq's approach achieve high throughput for LLM inference?
> 3. What workloads can Groq NOT handle that NVIDIA H100 can?

---

## 7. Graphcore — Intelligence Processing Unit

**Graphcore's IPU (Intelligence Processing Unit)** takes a graph-first approach to AI computation, reflecting that neural networks are computational graphs.

**IPU design philosophy**: Instead of organizing computation around vectors (SIMD) or matrix multiplications (systolic array), the IPU organizes computation around the **graph structure of neural networks**:
- Each processor in the IPU corresponds to a node in the computational graph
- Data flows along edges of the graph
- All local computation and data are in on-chip SRAM (no external DRAM)

**IPU Mk2 (2020):**
- 1472 IPU-processors, each with its own SRAM
- 900 MB total on-chip SRAM (distributed)
- ~250 TFLOPS FP16 (sparse)
- Bulk synchronous parallelism (BSP): alternating compute and communicate phases

**Applications**: Graph neural networks (GNNs), sparse computations, custom model architectures that don't map cleanly to matrix multiplies. IPUs struggle with dense transformer workloads that GPUs and TPUs dominate.

**Market reality**: Graphcore raised $700M+ and was highly hyped but failed to find widespread adoption in the transformer/LLM era. Acquired by SoftBank in 2024. A cautionary tale about betting on the "next big thing" in AI architecture that doesn't align with where the field actually goes.

### Quick Check
> 1. What is the key architectural metaphor behind the IPU design?
> 2. What type of neural network workload is the IPU designed to excel at?
> 3. Why might Graphcore have struggled as transformers became the dominant architecture?

---

## 8. AWS Trainium and Inferentia

**Amazon Web Services** built its own AI chips to reduce dependence on NVIDIA and offer lower-cost alternatives:

**AWS Inferentia (2019, 2nd gen Inf2 2023):**
- Purpose: LLM inference at scale for AWS customers
- Inferentia2 chip: 190 TOPS INT8, 380 TOPS INT4
- 32GB HBM2 per chip
- NeuronCore-v2 architecture (systolic arrays + vector engines)
- Connected via NeuronLink for multi-chip inference
- inf2.48xlarge instance: 12 Inferentia2 chips, 384GB HBM total

**AWS Trainium (2021, 2nd gen Trn2 2024):**
- Purpose: LLM training as an alternative to NVIDIA H100
- Trainium2 chip: 190 TFLOPS FP8, 380 TFLOPS FP4
- 96GB HBM3 per chip
- EFA (Elastic Fabric Adapter) for cluster-scale training across thousands of chips
- trn2.48xlarge instance: 16 Trainium2 chips, 1.5TB HBM total

**AWS Neuron SDK**: Programming framework for Inferentia and Trainium. Compiles PyTorch/JAX models to NeuronCore instructions. Less flexible than CUDA but sufficient for standard transformer models.

**The economics**: AWS Inferentia instances cost ~40% less than NVIDIA GPU instances for equivalent inference throughput on supported models. This makes LLM inference at scale significantly cheaper. Meta, Amazon's own services, and other large cloud consumers use AWS silicon.

**Other hyperscaler chips:**
- Google TPU (mentioned above)
- Microsoft Maia (2023): AI accelerator for Azure inference
- Meta MTIA (2023): Internal inference accelerator for Meta's recommendation systems

### Quick Check
> 1. What is the difference between AWS Inferentia and AWS Trainium?
> 2. Why does AWS build its own AI chips instead of just using NVIDIA?
> 3. What is the AWS Neuron SDK?

---

## Summary

- **NPUs and AI accelerators** are purpose-built for neural networks — removing all non-AI hardware, optimizing memory access, and using quantized arithmetic (INT8/FP8).
- **Systolic arrays** compute matrix multiplies by flowing data through grids of processing elements — each input value is read once and reused across many multiplications.
- **Google TPU** (2016+): first AI ASIC, Google-internal, powers all Google AI products. Vertically integrated with TensorFlow/JAX/XLA.
- **Smartphone NPUs** (Apple, Qualcomm, etc.): enable on-device AI inference with 1–45 TOPS at low power.
- **Cerebras WSE**: wafer-scale (entire silicon wafer = one chip), 44GB on-die SRAM, 100 PB/s bandwidth — eliminates memory bandwidth bottleneck.
- **Groq LPU**: deterministic execution, extreme inference throughput, specialized for LLM token generation.
- **Graphcore IPU**: graph-oriented AI processor. Good for sparse/graph workloads; struggled in the dense transformer era.
- **AWS Trainium/Inferentia**: hyperscaler-built AI chips for cost-efficient training and inference on cloud infrastructure.

---

## Exercises

### Easy
1. What is a systolic array and what operation is it optimized for?
2. Why do smartphone NPUs matter for battery life in AI applications?
3. How does Google TPU v1 use INT8 arithmetic for inference? What is the trade-off vs FP32?

### Medium
4. Cerebras WSE-2 has 40GB on-chip SRAM at 20 PB/s bandwidth. An H100 has 80GB HBM3 at 3.35 TB/s. For a 20GB model: (a) Can the model fit in WSE-2's on-chip SRAM? In H100's HBM? (b) For a workload that reads the entire model every 10ms (inference), what is the memory bandwidth required? Which chip supports this without bandwidth throttling? (c) What happens if the model is 60GB?
5. Systolic array efficiency: a 256×256 systolic array computing C = A × B where A is [1024×256] and B is [256×1024]: (a) How many multiply-accumulate operations total? (b) How many times is each element of A read from memory in a naive CPU implementation? In a systolic array? (c) What is the memory bandwidth reduction from using a systolic array?
6. Compare Groq's deterministic LPU to NVIDIA GPU for LLM inference: (a) Groq achieves 800 tokens/sec on Mixtral-8x7B; H100 achieves ~200 tokens/sec at batch=1. Why is Groq faster for low batch sizes? (b) At batch=32 (32 users simultaneously), H100 achieves ~2000 tokens/sec × 32 = effective 64K tokens/sec total; Groq remains at ~800 tokens/sec per user. Why doesn't Groq scale with batch size the same way? (c) For a chatbot serving 10,000 concurrent users, which hardware approach is more efficient?

### Hard
7. Design a custom AI accelerator for a specific workload: Vision Transformer (ViT) inference on 224×224 images at 1000 images/second with <10W power budget. The ViT has 12 transformer blocks, d_model=768, attention heads=12. (a) Estimate the FLOPs per image. (b) Design the systolic array size needed to achieve the throughput target. (c) How much on-chip SRAM is needed to hold weights without off-chip access? (d) What numerical format (FP32/FP16/INT8) would you use and why?
8. Analyze the AI accelerator startup landscape: Cerebras ($720M raised), Graphcore ($700M raised, then acquired), Groq ($300M+), SambaNova ($1.1B raised), Mythic AI, Hailo, etc. most have struggled to reach profitability despite strong hardware. (a) What makes it hard for startups to compete with NVIDIA despite building better hardware? (b) What makes it hard for hyperscaler chips (Google TPU, AWS Trainium) to displace NVIDIA for general AI workloads? (c) In what specific market niches do non-NVIDIA AI chips succeed?
