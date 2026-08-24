# Chapter 70: Heterogeneous Computing

The most performant computing systems today are not homogeneous (all the same type of core) — they are heterogeneous: a CPU for general-purpose control flow, a GPU for parallel compute, an NPU for matrix operations, a DSP for signal processing, all on the same SoC or connected in the same system. Every modern smartphone SoC (Snapdragon, A18, Exynos), every data center AI accelerator system (DGX, AWS Inferentia), and the highest-performance HPC supercomputers use heterogeneous computing. This chapter explains the principles, challenges, and real-world examples.

## Table of Contents

1. [What Is Heterogeneous Computing?](#1-what-is-heterogeneous-computing)
2. [CPU + GPU: The Classic Heterogeneous Pair](#2-cpu--gpu-the-classic-heterogeneous-pair)
3. [Scheduling and Work Distribution](#3-scheduling-and-work-distribution)
4. [Memory in Heterogeneous Systems](#4-memory-in-heterogeneous-systems)
5. [SoC Heterogeneous Computing](#5-soc-heterogeneous-computing)
6. [AMD APU and Intel Iris Xe — Integrated Heterogeneous](#6-amd-apu-and-intel-iris-xe--integrated-heterogeneous)
7. [Programming Models: CUDA, OpenCL, oneAPI](#7-programming-models-cuda-opencl-oneapi)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is Heterogeneous Computing?

**Heterogeneous computing** uses multiple types of processors with different characteristics — optimized for different types of computation — within the same system. The key insight: no single processor design is optimal for all workloads.

```
The "right tool for the job" principle:

  Task                  Best processor type
  ─────────────────────────────────────────
  Sequential logic      CPU (fast single-thread, OOO)
  Data-parallel math    GPU (thousands of small cores)
  ML inference          NPU (systolic array, fixed precision)
  Signal processing     DSP (MACs, circular buffers)
  Encryption            Crypto accelerator (AES-NI)
  Video decode          VPU/Media engine
  Compression (zlib)    QAT (Intel QuickAssist)
  
  Using a CPU for video decode: 15% CPU utilization
  Using a VPU for video decode: 1% CPU utilization, hardware handles it
```

**Why not just use CPUs for everything?**
- A CPU core is over-engineered for simple parallel tasks: OOO execution, branch prediction, complex cache hierarchy — all wasted for SIMD-friendly workloads
- Energy efficiency: a GPU FLOP costs 10–50× less energy than a CPU FLOP
- Throughput: 1 CPU core vs 4,096 GPU cores for matrix multiply — GPU wins by orders of magnitude

**Why not just use GPUs for everything?**
- Poor at sequential code, control flow, irregular memory access (database queries, OS scheduling)
- High latency: GPU kernel launch overhead ~5µs; CPU instruction latency <1ns
- Power: GPU always-on is wasteful for infrequent computation

**The heterogeneous solution**: CPU for coordination and sequential code; GPU/NPU/DSP for the parallel/specialized heavy lifting.

### Quick Check
> 1. What is heterogeneous computing and why is it more efficient than using only CPUs?
> 2. What is GPU kernel launch overhead and why does it matter?
> 3. Name three processor types that might appear in a modern smartphone SoC.

---

## 2. CPU + GPU: The Classic Heterogeneous Pair

The CPU + discrete GPU combination is the dominant heterogeneous computing paradigm for desktop, server, and HPC.

**CPU role:**
- Run the OS and application main thread
- Handle I/O, system calls, networking
- Manage GPU memory allocations
- Dispatch GPU compute kernels
- Handle results and control flow decisions

**GPU role:**
- Execute massively parallel compute kernels
- Matrix multiply, convolutions, physics simulation, ray tracing
- Everything that maps to SIMT execution

**Typical deep learning training workflow:**
```
CPU:  Load minibatch from storage → preprocess (resize, normalize)
      ↓
GPU:  Forward pass (matrix multiply, activations) [100ms]
      Backward pass (gradient computation) [200ms]
      ↓
CPU:  Collect gradients → optimizer step → update parameters
      ↓
Repeat for 1 million iterations
```

CPU handles the data pipeline and orchestration. GPU handles 95%+ of actual compute.

**Data transfer bottleneck**: CPU and GPU have separate DRAM (in discrete GPU systems). Data must be transferred over PCIe:
- PCIe 5 × 16 lanes: 128 GB/s bidirectional
- GPU HBM bandwidth: 3.35 TB/s (H100)
- PCIe bandwidth is 26× lower than GPU memory bandwidth

This means: for small data or frequent transfers, PCIe is the bottleneck. Solutions:
- NVLink: 900 GB/s GPU-to-GPU (bypasses PCIe) for multi-GPU training
- Unified memory (CUDA managed memory): OS handles migration automatically
- Apple's Unified Memory Architecture: CPU and GPU share the same physical RAM (Chapter 36)

### Quick Check
> 1. What is the PCIe bandwidth bottleneck in CPU+GPU systems?
> 2. What is NVLink and what problem does it solve?
> 3. What is the CPU's role in a deep learning training loop?

---

## 3. Scheduling and Work Distribution

How does a heterogeneous system decide what runs where?

**Static partitioning**: The programmer explicitly decides. `cudaMalloc` allocates GPU memory; `kernel<<<grid, block>>>` launches GPU work. Simple, but requires programmer expertise. Used for production ML training code.

**Runtime scheduling**: The runtime system decides. OpenCL devices, Apple Metal's command queues, DirectX 12 command lists. More flexible but higher overhead.

**OS-level heterogeneous scheduling** (Apple Silicon, Qualcomm Snapdragon):
The OS scheduler is aware of P-cores (fast) and E-cores (efficient):
- Foreground interactive apps → P-cores (low latency)
- Background processes → E-cores (save power)
- ML inference → Neural Engine (automatically, via Core ML API)

**GPGPU dispatch model (CUDA):**
```
CPU code:                              GPU code (kernel):
  
// Launch 1M threads in 1K blocks         __global__ void add(float* a, b, c) {
add<<<1000, 1024>>>(a, b, c);              int i = blockIdx.x*blockDim.x + threadIdx.x;
                                            c[i] = a[i] + b[i];
// CPU continues while GPU runs           }
result = collect_from_gpu();
```

The GPU runs 1 million threads, each doing one add. CPU launched the work and continues. GPU notifies when done via synchronization (cudaDeviceSynchronize or streams/events).

**Task graphs**: Modern CUDA (via CUDA Graphs) and DirectML (via DirectML Meta-commands) allow pre-recording a graph of GPU operations, then replaying with minimal overhead. This amortizes kernel launch overhead for repetitive workloads (video decoding, inference serving).

### Quick Check
> 1. What is the difference between static and dynamic partitioning of work in heterogeneous systems?
> 2. How does Apple's OS scheduler use P-cores and E-cores differently?
> 3. What are CUDA task graphs and what problem do they solve?

---

## 4. Memory in Heterogeneous Systems

Memory management is the hardest part of heterogeneous programming:

**Discrete GPU (NVIDIA, AMD):**
- CPU DRAM and GPU DRAM are physically separate
- Data must be explicitly copied: `cudaMemcpy(gpu_buf, cpu_buf, size, H2D)`
- Unified Virtual Addressing (UVA): CPU and GPU share a virtual address space — but physical copies still happen
- CUDA Unified Memory: automatic migration on fault (OS migrates pages)

**Integrated GPU (Intel HD, AMD APU, Apple M1–M4):**
- CPU and GPU share the same physical DRAM
- No data copy needed — just pass a pointer
- Coherency: GPU and CPU see consistent data (if memory is coherent)
- Bandwidth sharing: CPU and GPU compete for the same memory bandwidth
- Apple's Unified Memory Architecture: special bandwidth + coherence management per access type

**Memory coherence in heterogeneous systems:**
```
Problem:
  CPU writes data to address A (to L1 cache, not yet to DRAM)
  GPU reads address A → reads old value from DRAM!
  
  Need: flush CPU caches before GPU access
  
Solution options:
  1. Explicit flush/invalidate (low-level, error-prone)
  2. Coherent interconnect: GPU snoops CPU cache (expensive)
  3. Non-coherent with barriers: programmer inserts flush calls
  4. Unified cache (Apple's approach on M-series): one cache hierarchy
```

**HSA (Heterogeneous System Architecture)**: An industry standard (AMD, ARM, QUALCOMM) for coherent CPU+GPU memory. "hUMA (heterogeneous Uniform Memory Access)" — GPU can access CPU memory and vice versa with the same pointer, with hardware-maintained coherence.

### Quick Check
> 1. Why do discrete GPU systems require explicit data copies between CPU and GPU?
> 2. What is CUDA Unified Memory and how does it simplify programming?
> 3. What does Apple's Unified Memory Architecture eliminate vs a discrete GPU system?

---

## 5. SoC Heterogeneous Computing

Mobile SoCs (Snapdragon, Apple A-series, Exynos, MediaTek Dimensity) are the most sophisticated heterogeneous computing platforms:

**Qualcomm Snapdragon 8 Gen 3 heterogeneous components:**
```
CPU cluster:
  1× Cortex-X4 (3.3 GHz, prime core) → single-threaded apps
  5× Cortex-A720 (3.15 GHz, perf cores) → multi-threaded apps
  2× Cortex-A520 (2.27 GHz, eff cores) → background tasks
  
GPU:
  Adreno 750 → gaming, UI rendering, parallel ML
  
NPU:
  Hexagon NPU (45 TOPS) → inference, photo AI, on-device LLM
  
DSP:
  Hexagon V75 (HVX) → audio, sensor fusion, always-on processing
  
ISP:
  Spectra ISP → camera pipeline, video encode/decode
  
Modem:
  Snapdragon X75 5G → cellular connectivity
  
GNSS:
  FastConnect GPS → location
```

All these compute units share system memory (LPDDR5X) via the AMBA NoC, with each unit having appropriate QoS (Quality of Service) settings for bandwidth and latency.

**Work allocation on a Snapdragon device:**
- Taking a photo: ISP, NPU (scene recognition, HDR), DSP (noise reduction), GPU (post-processing)
- Playing a game: GPU (rendering), CPU (game logic), NPU (AI NPCs)
- Navigation: GNSS, CPU (routing algorithm), DSP (audio turn-by-turn)
- Always listening for voice: Hexagon micro-DSP at <5mW (main CPU is asleep)

### Quick Check
> 1. How many different processor types are in a Qualcomm Snapdragon 8 Gen 3?
> 2. What is the role of the Hexagon micro-DSP in always-on voice recognition?
> 3. What is QoS (Quality of Service) in the context of an SoC bus fabric?

---

## 6. AMD APU and Intel Iris Xe — Integrated Heterogeneous

**AMD APU (Accelerated Processing Unit):**
AMD's brand for their integrated CPU+GPU chips. The AMD Ryzen AI 9 HX 370 (2024) includes:
- CPU: 12 cores Zen 5 (4 big + 8 efficient)
- GPU: Radeon 890M (16 RDNA 3.5 compute units, 768 shaders)
- NPU: XDNA 2 AI engine (50 TOPS INT4)
- Shared memory: LPDDR5X (shared pool, 51.2 GB/s)

Key: CPU and GPU share the same DRAM — no discrete GPU memory copies. The GPU performance is lower than a discrete GPU (bandwidth limited by shared LPDDR5X at 51 GB/s vs discrete GPU HBM at 3 TB/s).

**Intel Iris Xe / Arc:**
- Intel's integrated GPU on every Intel Core processor
- Xe-LP (low power): 80 EUs (Execution Units) in i7 mobile
- Xe-HPG (high performance graphics): in Arc A770 discrete GPU
- Intel maintains a single ISA across integrated and discrete Xe GPUs — same code runs on both

**Trade-offs: integrated vs discrete GPU:**
```
                    Integrated GPU      Discrete GPU
Memory bandwidth    20–100 GB/s         900–3,350 GB/s
Peak performance    0.5–10 TFLOPS       5–200 TFLOPS
Power               5–15W (GPU portion) 50–700W
Cost                Included in CPU     $150–$1500+
Latency (data)      Zero copy           PCIe copy overhead
Battery life        Excellent           Poor
Use case            Light gaming, CAD   ML training, heavy gaming, rendering
```

### Quick Check
> 1. What is an AMD APU and what does it integrate?
> 2. Why does an integrated GPU have lower performance than a discrete GPU despite sharing the processor?
> 3. What is Intel Xe and how does it span from integrated to discrete?

---

## 7. Programming Models: CUDA, OpenCL, oneAPI

**CUDA** (NVIDIA, 2007):
- Proprietary to NVIDIA GPUs
- C/C++ extensions: `__global__`, `__device__`, `__shared__`
- Most mature ecosystem: cuDNN, TensorRT, NCCL, cuBLAS
- De facto standard for ML training

**OpenCL** (Khronos Group, open standard):
- Works on CPUs, GPUs (NVIDIA, AMD, Intel), FPGAs
- More verbose than CUDA; slightly lower performance
- C-based kernel language
- Used where portability across hardware is required

**SYCL / Intel oneAPI:**
- Modern C++ (C++17) alternative to OpenCL
- Intel's primary programming model for CPU+GPU+FPGA
- DPC++ (Data Parallel C++) — based on SYCL
- Aims to replace CUDA for Intel hardware

**Metal (Apple):**
- Apple's GPU programming API (iOS, macOS)
- C++ Metal shading language
- Integrated with Core ML for neural network acceleration

**Abstraction layers:**
- PyTorch: uses CUDA/MPS/ROCm under the hood; python user doesn't see device code
- TensorFlow: similar abstraction
- OpenNN: portable neural network library using OpenCL

**The CUDA ecosystem lock-in**: Almost all ML research uses CUDA. Migrating to AMD ROCm or Intel oneAPI requires porting, and many CUDA-specific optimizations don't translate. NVIDIA's CUDA monopoly in ML training is partly a software ecosystem problem, not just hardware.

### Quick Check
> 1. What is CUDA and why is it the dominant GPU programming model?
> 2. What is OpenCL and when would you prefer it over CUDA?
> 3. What is Apple Metal and how does it relate to the iPhone's GPU?

---

## Summary

- **Heterogeneous computing**: use the right processor for each task — CPU for control, GPU for parallel math, NPU for ML, DSP for signals.
- **CPU+GPU**: classic pair. CPU dispatches kernels; GPU executes in SIMT mode. PCIe bandwidth is the bottleneck for discrete GPUs.
- **Memory**: discrete GPU requires explicit copies; integrated GPU shares DRAM; Apple UMA eliminates copies.
- **Scheduling**: static (programmer-explicit), runtime (OS-managed), or automatic (via high-level APIs like Core ML).
- **Mobile SoC**: most sophisticated example — 6+ different processor types all on one die, shared memory, OS-managed dispatch.
- **Integrated GPU**: AMD APU, Intel Iris Xe — eliminates copies but bandwidth-limited.
- **Programming**: CUDA (NVIDIA), OpenCL (portable), SYCL/oneAPI (Intel), Metal (Apple).

---

## Exercises

### Easy
1. Why is a GPU better than a CPU for matrix multiplication?
2. What is the PCIe bandwidth bottleneck in CPU+GPU systems?
3. What is CUDA and what does a `__global__` function mean?

### Medium
4. Roofline model: An H100 GPU has: memory bandwidth = 3.35 TB/s, peak FP16 performance = 989 TFLOPS. A matrix multiply of two 4096×4096 FP16 matrices: (a) FLOPs required: 2 × 4096³ operations (standard matrix multiply). (b) Memory bytes: input 2×(4096²×2 bytes), output 4096²×2 bytes. (c) Arithmetic intensity = FLOPs/bytes. (d) The roofline model says: performance = min(arithmetic intensity × bandwidth, peak FLOPs). Which bottleneck applies? (e) How does batch size change arithmetic intensity for matrix multiply?
5. Heterogeneous dispatch latency: A photo processing pipeline on a smartphone: (a) Face detection (NPU, 20ms), (b) HDR processing (ISP, 5ms), (c) Noise reduction (DSP, 10ms), (d) Color correction (GPU, 8ms), (e) UI rendering (GPU, 16ms). If all steps are sequential: total latency? If ISP and DSP run in parallel, then GPU: total latency? At what step does the pipeline need to synchronize the outputs?
6. Memory bandwidth allocation: A Snapdragon 8 Gen 3 has LPDDR5X at 51.2 GB/s. Simultaneous workload: (a) CPU: 2 active threads, each needing 4 GB/s → 8 GB/s total. (b) GPU: rendering at 12 GB/s. (c) NPU: inference at 20 GB/s. (d) DSP: audio at 2 GB/s. (e) Modem: 2 GB/s. Total: 44 GB/s. What headroom remains? If all units demand peak simultaneously (CPU 8, GPU 25, NPU 30, DSP 5, modem 3 = 71 GB/s): which gets throttled and by what mechanism?

### Hard
7. Unified memory performance: An AMD Ryzen AI APU with 51.2 GB/s shared LPDDR5X. A ResNet-50 inference task needs: (a) Load model weights: 100MB. (b) Inference: 4 GFLOPS per image, 100 images/second = 400 GFLOPS/s. (c) Memory access during inference: 200 MB/s bandwidth for activations. GPU ROPs: 8 TFLOPS. Does memory bandwidth or compute limit throughput? What is the bottleneck? (b) Compare to discrete RX 7600 (9 TFLOPS GPU, 288 GB/s GDDR6): how does the performance change? (c) For which model sizes does APU beat discrete GPU (hint: smaller models are memory-bandwidth limited)? (d) Power: APU runs at 25W total, discrete GPU at 120W. For 100 inferences/second: energy per inference comparison?
8. Cross-architecture ML optimization: A startup is building an edge AI chip for autonomous vehicles. They need: 50 TOPS INT8 AI inference, < 20W, real-time video (30fps 4K), sensor fusion (LiDAR, radar, cameras), Linux. (a) Design the heterogeneous compute architecture: what processor types do you need? (b) Select process node and area budget for each block. (c) The AI workload is 80% convolutions, 15% attention, 5% preprocessing. How do you allocate compute resources? (d) Memory: 16GB LPDDR5, or HBM? Trade-offs for automotive. (e) If PyTorch is the training framework: what compiler/runtime stack converts the trained model to run on your custom hardware? (hint: ONNX → custom backend; TVM; custom compiler)
