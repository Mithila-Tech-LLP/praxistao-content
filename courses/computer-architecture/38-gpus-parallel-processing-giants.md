# Chapter 38: GPUs — Parallel Processing Giants

For decades, CPUs improved by getting faster at individual tasks — deeper pipelines, better branch prediction, more out-of-order execution. GPUs took the opposite approach: instead of one incredibly fast, clever processor, build thousands of simpler processors and run them all in parallel. This design made GPUs unbeatable for graphics — where you need to shade millions of pixels independently and simultaneously. Then the world discovered that the same massive parallelism is perfect for machine learning, scientific simulation, cryptocurrency mining, and video encoding. The GPU is no longer just a graphics card. It is the engine of the AI revolution.

## Table of Contents

1. [Why Graphics Needs Parallelism](#1-why-graphics-needs-parallelism)
2. [GPU Architecture vs CPU Architecture](#2-gpu-architecture-vs-cpu-architecture)
3. [SIMT — Single Instruction Multiple Threads](#3-simt--single-instruction-multiple-threads)
4. [The Graphics Pipeline](#4-the-graphics-pipeline)
5. [GPU Memory: GDDR and HBM](#5-gpu-memory-gddr-and-hbm)
6. [General-Purpose GPU Computing (GPGPU)](#6-general-purpose-gpu-computing-gpgpu)
7. [The GPU in the AI Era](#7-the-gpu-in-the-ai-era)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Graphics Needs Parallelism

Consider rendering a single frame in a video game at 1920×1080 resolution. That's 2,073,600 pixels. Each pixel must be:
1. Shaded (determine its color based on geometry, textures, lights, shadows)
2. Processed for post-effects (anti-aliasing, bloom, depth of field)
3. Written to the framebuffer

At 60 frames per second, you have 16.7 milliseconds per frame. With 2 million pixels per frame: you must process ~120 million pixels per second. Each pixel requires multiple texture lookups, interpolations, and arithmetic operations.

A CPU with 8 cores could process these pixels serially — but at maybe 500 million arithmetic operations per second per core, it would take far too long. The key insight: **each pixel's shading is independent of all other pixels**. This is **data parallelism** — the same computation (shade this pixel) applied to millions of independent data items simultaneously.

This is fundamentally different from the kind of parallelism CPUs exploit (instruction-level parallelism within a single sequential program). Data-parallel workloads don't need:
- Complex out-of-order execution
- Branch prediction (shaders rarely branch unpredictably)
- Large caches (each pixel accesses different data — poor temporal locality)

They do need:
- Thousands of simple arithmetic units
- High memory bandwidth (millions of texture reads per frame)
- Fixed-function hardware for common operations (texture sampling, rasterization)

### Quick Check
> 1. Why is pixel shading a perfect example of "data parallelism"?
> 2. At 1080p, 60 fps, with 100 shader operations per pixel: how many operations per second are needed?
> 3. What CPU features are unnecessary for GPU-style data-parallel work?

---

## 2. GPU Architecture vs CPU Architecture

The fundamental difference between a CPU and a GPU is how they allocate transistors:

```
CPU (e.g., Intel Core i9, 8 cores):
  ┌──────────────────────────────────────────────────────────┐
  │  ~60% die area: caches (L1/L2/L3)                        │
  │  ~20% die area: OOO execution logic (ROB, RS, rename)    │
  │  ~10% die area: branch predictor, fetch logic            │
  │  ~10% die area: actual ALUs (arithmetic units)           │
  └──────────────────────────────────────────────────────────┘
  
  8 cores × ~4 FP operations/cycle × 4 GHz = 128 GFLOPS peak

GPU (e.g., NVIDIA RTX 4090):
  ┌──────────────────────────────────────────────────────────┐
  │  ~70% die area: shader cores (16,384 CUDA cores)         │
  │  ~15% die area: memory interface (ROPs, TMUs)            │
  │  ~10% die area: fixed-function units (RT cores, Tensor)  │
  │   ~5% die area: control logic, caches                    │
  └──────────────────────────────────────────────────────────┘
  
  16,384 cores × 2 FP ops/cycle × 2.52 GHz = 82.6 TFLOPS peak
```

The GPU has ~640× more peak FP32 throughput than a high-end desktop CPU. But it achieves this by:
1. Being much larger (RTX 4090 die: 608 mm², vs Intel i9 die: ~200mm²)
2. Having tiny, simple shader cores without the complex OOO machinery
3. Running at lower clock speeds (2.52 GHz vs 5.5 GHz for CPU)
4. Using extreme levels of thread-level parallelism to hide memory latency

**CPU vs GPU: a metaphor**: The CPU is a Formula 1 racing car — one incredibly fast, sophisticated machine. The GPU is a bus convoy — thousands of simpler vehicles, each slower individually but collectively carrying far more passengers simultaneously.

### Quick Check
> 1. What fraction of a CPU's die area is dedicated to actual arithmetic units vs. control logic?
> 2. How does a GPU hide memory latency if it has smaller caches than a CPU?
> 3. Why does a GPU run at lower clock speeds than a CPU?

---

## 3. SIMT — Single Instruction Multiple Threads

GPUs use a execution model called **SIMT (Single Instruction, Multiple Threads)** — related to but distinct from CPU's SIMD (Single Instruction, Multiple Data).

**SIMD (CPU)**: One instruction, one thread, processes multiple data items in parallel using wide registers (e.g., AVX-512 adds 16 floats at once).

**SIMT (GPU)**: Groups of threads (typically 32 — called a **warp** in NVIDIA terminology, a **wavefront** in AMD terminology) execute the same instruction simultaneously, each with their own registers and memory addresses.

```
SIMT warp execution:
  Thread 0:  ADD R0, R1, R2    (R0=0x1000, R1=pixel[0])
  Thread 1:  ADD R0, R1, R2    (R0=0x1040, R1=pixel[1])
  ...
  Thread 31: ADD R0, R1, R2    (R0=0x17C0, R1=pixel[31])
  
  All 32 threads execute the same ADD instruction,
  but each thread has its own set of registers and memory addresses.
  
  Hardware: one instruction fetch + decode, 32 ALU executions in parallel
```

**Divergence problem**: When threads in a warp take different branches (`if thread.id > 16 { ... }`), the hardware must execute both paths sequentially (masking the threads that should not execute each path). This is **warp divergence** and is a major GPU performance problem.

```
Warp divergence:
  if (x > 0.5f) {           ← threads 0-20 take this path
      result = expensive();  ← cycles 1-100: only threads 0-20 active
  } else {                   ← threads 21-31 take this path
      result = cheap();      ← cycles 101-110: only threads 21-31 active
  }
  // Total: 110 cycles, but useful work only 50% of the time
  // Without divergence: 100 cycles (all 32 threads execute expensive())
```

**Hiding memory latency with thread parallelism**: When a warp is waiting for a memory load (which may take 500+ cycles on GDDR), the GPU switches instantly to another warp that is ready to execute. This **warp switching** is essentially free (all warps' register files are stored on-chip). A GPU might have 32–64 warps active simultaneously, so while most wait for memory, some are always executing.

This is the opposite of CPU philosophy: CPUs have big caches to avoid memory latency; GPUs have many threads to hide it.

### Quick Check
> 1. What is a "warp" in NVIDIA GPU terminology?
> 2. Explain warp divergence. When does it occur and what is its performance impact?
> 3. How does SIMT's thread-level parallelism hide memory latency differently from CPU caches?

---

## 4. The Graphics Pipeline

A GPU processes 3D graphics through a **fixed sequence of stages** — the **graphics pipeline**. Modern GPUs use a mix of fixed-function hardware (for efficiency) and programmable shader stages (for flexibility).

```
Vertex Shader → (Hull Shader → Tessellation → Domain Shader) → Geometry Shader
→ Rasterization (fixed) → Fragment/Pixel Shader → Render Output Units (ROPs)

Modern simplified pipeline:
  1. Vertex Shader (programmable): transform 3D vertices → 2D screen coordinates
  2. Primitive Assembly (fixed): assemble vertices into triangles
  3. Rasterization (fixed): determine which pixels each triangle covers
  4. Fragment Shader (programmable): calculate pixel color (texture sampling, lighting)
  5. Depth Test (fixed): discard pixels hidden behind other geometry
  6. ROP (Render Output Unit): blend, write final pixel to framebuffer
```

**Key hardware units:**

**SM (Streaming Multiprocessor)** [NVIDIA]: the basic GPU compute unit. Contains:
- 128 CUDA cores (FP32 ALUs)
- 4 Tensor cores (matrix multiply hardware)
- 1 RT core (ray tracing hardware)
- L1 cache / shared memory (64-96KB)
- Warp schedulers (4 per SM, each handling 16 warps)

An RTX 4090 has 128 SMs × 128 CUDA cores = 16,384 CUDA cores total.

**TMU (Texture Mapping Unit)**: Specialized hardware for texture sampling — takes a UV coordinate and returns the interpolated texel (texture pixel) value. Uses bilinear/trilinear filtering. Hundreds of TMUs in a high-end GPU; each handles one texture sample per cycle.

**ROP (Render Output Unit)**: Handles the final stage — blending transparency, writing pixels to the framebuffer, depth testing. Bandwidth to VRAM is the bottleneck here.

**Ray Tracing Cores** (NVIDIA Turing+, 2018): Dedicated hardware for ray-triangle intersection tests — the fundamental operation in ray tracing. Replaces a software operation that took dozens of shader instructions with single-cycle hardware.

### Quick Check
> 1. What does a vertex shader do? What does a fragment shader do?
> 2. What is a "Streaming Multiprocessor" (SM) in NVIDIA terminology?
> 3. What are dedicated ray tracing cores for, and why were they added as fixed-function hardware?

---

## 5. GPU Memory: GDDR and HBM

GPU performance is often memory-bandwidth limited. The shaders can produce results faster than they can read textures and write framebuffers. The GPU memory system is designed for maximum bandwidth:

**GDDR (Graphics Double Data Rate)**:

| Version | Bandwidth (per chip) | Total Bandwidth (typical high-end GPU) |
|---------|---------------------|----------------------------------------|
| GDDR5 (2008) | ~30 GB/s per chip | ~240 GB/s (8 chips) |
| GDDR6 (2018) | ~56 GB/s per chip | ~448 GB/s |
| GDDR6X (2020) | ~84 GB/s per chip | ~672 GB/s (RTX 3090 Ti) |
| GDDR7 (2024) | ~192 GB/s per chip | ~1.5 TB/s potential |

The RTX 4090 uses 24GB GDDR6X with a 384-bit bus at 21 Gbps → **1008 GB/s** total bandwidth.

**HBM (High Bandwidth Memory)**:
Used in professional GPUs (A100, H100) and some high-end consumer GPUs:

| HBM | Bandwidth (per stack) |
|-----|----------------------|
| HBM1 (2015) | 128 GB/s |
| HBM2 (2016) | 256 GB/s |
| HBM2E (2019) | 460 GB/s |
| HBM3 (2022) | 819 GB/s per stack |
| HBM3E (2024) | 1.2 TB/s per stack |

NVIDIA H100 SXM: 5 active HBM3 stacks totaling **3.35 TB/s** total bandwidth. More than 3× the RTX 4090. (The per-stack figures above are peak numbers for each HBM generation; a shipping product often clocks its stacks somewhat below the peak, which is why 5 stacks land at 3.35 TB/s rather than 5 × 819 GB/s.)

**VRAM capacity**:
- Consumer GPU (RTX 4090): 24GB GDDR6X — adequate for gaming and many AI tasks
- Professional GPU (NVIDIA H100 SXM): 80GB HBM3 — for large AI models (GPT-4 has ~1.8 trillion parameters)
- Multi-GPU (NVLink + NVSwitch): connect 8 H100s for 640GB total VRAM visible as one pool

### Quick Check
> 1. Why does GDDR use a wider bus (384-bit) rather than just faster clock speeds to achieve high bandwidth?
> 2. What is the main advantage of HBM over GDDR in terms of bandwidth?
> 3. Why does an AI researcher care about GPU VRAM capacity?

---

## 6. General-Purpose GPU Computing (GPGPU)

The insight that GPUs could do non-graphics computation was transformative. In 2006, NVIDIA released **CUDA (Compute Unified Device Architecture)** — a programming model and compiler that allowed programmers to write C-like code that runs on the GPU's shader cores.

**GPGPU programming model:**
- CPU is the "host" — orchestrates, manages data movement
- GPU is the "device" — executes data-parallel "kernels"
- Explicit memory management: copy data from CPU RAM to GPU VRAM, compute, copy back

```c
// Simple vector addition CUDA kernel
__global__ void vectorAdd(float* a, float* b, float* c, int n) {
    int idx = blockIdx.x * blockDim.x + threadIdx.x;  // compute thread's index
    if (idx < n) {
        c[idx] = a[idx] + b[idx];
    }
}

// Launch: 1M threads organized in blocks of 256
vectorAdd<<<(1000000/256), 256>>>(d_a, d_b, d_c, 1000000);
```

The programmer specifies a **kernel** (the per-thread function) and a **grid** (number of threads and their organization into blocks). The GPU scheduler maps blocks to SMs and threads to CUDA cores.

**GPGPU applications:**
- **Scientific computing**: molecular dynamics, climate modeling, fluid simulation
- **Machine learning**: training and inference for neural networks (Chapter 39)
- **Cryptocurrency**: SHA-256 hash computation (Bitcoin mining), Ethash (Ethereum)
- **Video transcoding**: GPU-accelerated H.264/H.265 encoding
- **Database acceleration**: PostgreSQL with pg_strom, RAPIDS cuDF

**OpenCL**: NVIDIA's competitor to CUDA, supported by AMD, Intel, and others. More portable but less optimized than CUDA for NVIDIA hardware.

**ROCm**: AMD's CUDA equivalent for Radeon GPUs. The ecosystem is growing but still less mature than CUDA.

### Quick Check
> 1. What does CUDA stand for and what did it enable?
> 2. In CUDA programming, what is the difference between the "host" and "device"?
> 3. Explain what a "kernel" and a "grid" are in CUDA terminology.

---

## 7. The GPU in the AI Era

The rise of deep learning transformed the GPU market. Training a neural network involves billions of matrix multiply operations — the same operations GPUs are designed to execute in parallel. An operation that takes hours on a CPU takes minutes on a GPU.

**Why GPUs excel at neural network training:**
1. Matrix multiply is data parallel: every element of the output can be computed independently
2. GPUs have enormous floating-point throughput (TFLOPS)
3. High memory bandwidth: model weights must be streamed through the GPU rapidly
4. Tensor Cores (NVIDIA Volta+, 2017): dedicated mixed-precision matrix multiply units (FP16 inputs, FP32 accumulation)

```
NVIDIA Tensor Core (A100):
  Performs a 4×4 matrix multiply per cycle: 4×4×4 = 64 multiply-accumulate ops
  512 Tensor Cores × 64 ops × 2 GHz = 312 TFLOPS (FP16 Tensor Core throughput)
  vs 19.5 TFLOPS for standard FP32 CUDA cores
  
  The Tensor Core is ~16× faster for matrix multiply than CUDA cores
```

**The NVIDIA monopoly**: NVIDIA's CUDA ecosystem — libraries (cuBLAS, cuDNN, TensorRT), tools (Nsight), frameworks (PyTorch, TensorFlow use CUDA) — has created an enormous software moat. Even when AMD's hardware (MI300X) matches or exceeds H100 in raw FLOPS or memory bandwidth, the software ecosystem pulls AI researchers back to NVIDIA.

**The AI training stack:**
- Framework: PyTorch / TensorFlow
- Compiler: torch.compile / XLA (generates optimized GPU code)
- Libraries: cuBLAS (matrix multiply), cuDNN (convolutions, attention), NCCL (multi-GPU communication)
- Hardware: A100/H100/H200 GPU with NVLink for multi-GPU, InfiniBand for multi-node

**NVIDIA GPU generations (AI focus):**
| GPU | Year | FP16 Tensor TFLOPS | HBM | Process |
|-----|------|---------------------|-----|---------|
| V100 | 2017 | 112 TFLOPS | 32GB HBM2 | TSMC 12nm |
| A100 | 2020 | 312 TFLOPS | 80GB HBM2e | Samsung 7nm |
| H100 SXM | 2022 | 989 TFLOPS (FP16) | 80GB HBM3 | TSMC 4nm |
| H200 | 2024 | 989 TFLOPS + 141GB HBM3e | 141GB HBM3e | TSMC 4nm |
| B100/B200 | 2024 | FP4 and FP8 dominant | 192GB HBM3e | TSMC 4nm |

The H100 ($30,000–$40,000 per GPU) is the foundation of ChatGPT, Claude, and every major large language model's training infrastructure.

### Quick Check
> 1. What are NVIDIA Tensor Cores and why are they 16× faster than CUDA cores for matrix multiplication?
> 2. Why does NVIDIA have such a dominant position in AI training hardware?
> 3. What is the memory capacity of an NVIDIA H100 SXM, and why does AI training require so much GPU memory?

---

## Summary

- Graphics rendering requires processing millions of independent pixels simultaneously — **data parallelism** that GPUs are designed for.
- GPUs allocate most transistors to arithmetic units (not caches or OOO logic), trading single-thread cleverness for massive parallel throughput.
- **SIMT** execution: 32 threads (a warp) execute the same instruction simultaneously; diverging branches cause performance penalties.
- GPU memory latency is hidden by **warp switching** — when one warp waits for memory, another ready warp executes.
- The **graphics pipeline** (vertex shader → rasterization → fragment shader → ROP) processes triangles to pixels. Ray tracing cores accelerate ray-triangle intersection.
- GPU memory: **GDDR** (high bandwidth, moderate capacity, consumer GPUs) vs **HBM** (very high bandwidth, larger capacity, professional GPUs).
- **GPGPU via CUDA**: GPUs became general computing engines for ML, scientific computing, and more.
- **Tensor Cores** accelerate matrix multiply for AI training — the H100 achieves ~1000 TFLOPS in FP16, making modern LLM training possible.

---

## Exercises

### Easy
1. Explain data parallelism in the context of pixel shading. Why are pixels independent?
2. A GPU has 10,000 shader cores running at 2 GHz, each doing 2 FP32 ops/cycle. What is the peak FP32 throughput?
3. What is warp divergence and when does it occur?

### Medium
4. A GPU with 10,000 CUDA cores at 1.5 GHz processes data from 16GB GDDR6X at 768 GB/s. A workload performs 2 FP32 operations per 4-byte value loaded. (a) Is this workload compute-bound or memory-bandwidth-bound? (b) What is the arithmetic intensity (FLOPs/byte)? (c) What intensity would make it compute-bound?
5. A GPU shader program has the following:
   ```
   if (thread_id % 4 == 0) { heavyCompute(); }  // ~100 cycles
   else { lightCompute(); }                       // ~10 cycles
   ```
   In a 32-thread warp: (a) how many threads take each branch? (b) how many cycles does the warp spend total with divergence? (c) what is the efficiency (useful work / total cycles)?
6. Memory bandwidth vs arithmetic speed trade-off in GPU design: if you double the number of shader cores but keep memory bandwidth the same, does performance double for (a) a compute-bound shader (doing 100 FP ops per 4B loaded), (b) a memory-bound shader (doing 2 FP ops per 4B loaded)?

### Hard
7. Large Language Model inference on GPU: GPT-3 has 175 billion FP16 parameters (~350GB). An NVIDIA A100 has 80GB HBM2e. (a) How many A100s are needed just to hold the weights? (b) At 312 TFLOPS FP16 per A100, and assuming a token generation requires 2 FLOPs per parameter, how many tokens/second can one A100 generate? (c) With tensor parallelism across 8 A100s (sharing the model), what is the interconnect bandwidth requirement if weights must be redistributed every forward pass?
8. The CUDA programming model vs CPU threading: write pseudocode (not real code) for parallel matrix multiplication in both CUDA (kernel + grid launch) and CPU pthreads/OpenMP. Describe: (a) how work is divided in each model, (b) how memory is accessed in each model, (c) what synchronization primitives each uses, (d) why the GPU approach scales to 10,000 threads but the CPU approach caps out at ~32 useful threads for a single NUMA node.
