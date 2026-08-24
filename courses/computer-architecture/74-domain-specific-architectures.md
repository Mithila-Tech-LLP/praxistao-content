# Chapter 74: Domain-Specific Architectures

Chapter 72 introduced Domain-Specific Architectures (DSAs) as a key post-Moore strategy. This chapter goes deep: what they are, why they outperform general-purpose processors, the key examples across domains, and how they are designed. Hennessy and Patterson — the co-authors of the canonical computer architecture textbook and winners of the 2017 Turing Award — argued that the "Golden Age of Computer Architecture" is driven precisely by DSAs. This chapter shows you what that means in practice.

## Table of Contents

1. [What Is a DSA and Why Now?](#1-what-is-a-dsa-and-why-now)
2. [The Systolic Array — Foundation of AI Accelerators](#2-the-systolic-array--foundation-of-ai-accelerators)
3. [Google TPU — The Canonical DSA](#3-google-tpu--the-canonical-dsa)
4. [Neural Processing Units (NPUs) in Mobile SoCs](#4-neural-processing-units-npus-in-mobile-socs)
5. [Video Codec Accelerators](#5-video-codec-accelerators)
6. [Network and Storage Accelerators](#6-network-and-storage-accelerators)
7. [Crypto and Security Accelerators](#7-crypto-and-security-accelerators)
8. [How to Design a DSA](#8-how-to-design-a-dsa)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What Is a DSA and Why Now?

**Definition**: A Domain-Specific Architecture (DSA) is a processor or hardware accelerator designed specifically for one class of application, trading generality for efficiency.

**The efficiency spectrum:**
```
General-purpose CPU:
  - Can run any program
  - Complex microarchitecture: OOO, speculation, deep cache hierarchy
  - Efficiency: ~1 TOPS/W for INT8 AI inference

GPU:
  - Optimized for SIMT (same instruction, many threads) workloads
  - CUDA cores for general parallel + Tensor Cores for matmul
  - Efficiency: ~5–15 TOPS/W for INT8 AI inference

DSA (Google TPU v4):
  - Optimized specifically for matrix multiply + activation functions
  - Systolic array, no cache hierarchy, high internal bandwidth
  - Efficiency: ~100 TOPS/W for INT8 AI inference
  
Dedicated fixed-function block (AV1 video decode):
  - Executes exactly one algorithm (AV1 decode)
  - Hard-wired logic, not programmable
  - Efficiency: 1000–10000× more efficient than CPU software
  
    CPU    GPU    DSA   Fixed-function
     │      │      │         │
  ←──┴──────┴──────┴─────────┴────────→
  Most general                Most efficient
```

**Why DSAs are dominant now:**
1. **End of Dennard Scaling**: No free performance from process nodes → must architect for efficiency
2. **Workload consolidation**: Cloud providers (Google, Microsoft, Meta) run billions of identical inference requests → amortize DSA NRE cost
3. **Algorithm stability**: Matrix multiply, convolution, attention — these core ML operations have been stable enough to build hardware for
4. **Energy dominates cost**: Data center power is $0.05–0.12/kWh; efficient chips save millions of dollars/year

**Hennessy & Patterson's prediction (2017 Turing lecture):**
> "The next decade will see a Cambrian explosion of new computer architectures. It is probably the most exciting time in computer architecture since the 1980s."

Seven years later: Google TPU, Apple Neural Engine, AWS Inferentia/Trainium, Microsoft Maia, Meta MTIA, NVIDIA Tensor Cores, SHAKTI PQC accelerators, Qualcomm Hexagon NPU — the prediction has come true.

### Quick Check
> 1. What is the efficiency difference between a general-purpose CPU and a dedicated DSA for INT8 inference?
> 2. Why does workload consolidation (e.g., Google search inference) justify DSA NRE costs?
> 3. What year did Hennessy and Patterson give their Turing Award lecture predicting the DSA era?

---

## 2. The Systolic Array — Foundation of AI Accelerators

The **systolic array** is the core building block of most AI accelerators. The name comes from the heart: like blood pumping rhythmically through vessels, data flows rhythmically through an array of processing elements.

**Matrix multiply on a systolic array (simplified 3×3 example):**
```
Problem: multiply matrix A (3×3) by matrix B (3×3) to get C (3×3)

C[i][j] = Σ A[i][k] × B[k][j]    for k = 0, 1, 2

Systolic array structure:
  
     b00  b01  b02    ← B values flow rightward →
      ↓    ↓    ↓
a00→ PE  → PE  → PE   Row 0
      ↓    ↓    ↓
a10→ PE  → PE  → PE   Row 1
      ↓    ↓    ↓
a20→ PE  → PE  → PE   Row 2

Each PE (Processing Element):
  - Receives a_value from left
  - Receives b_value from above
  - Multiplies and adds to its accumulator: acc += a × b
  - Passes a_value to the right
  - Passes b_value downward

Cycles:
  Cycle 1: a00 enters PE[0][0], b00 enters PE[0][0]: acc[0][0] += a00×b00
  Cycle 2: a00→PE[0][1], a10→PE[1][0], b00→PE[1][0], b01→PE[0][1]: more multiplies
  ...after 5 cycles all 9 output values computed simultaneously
  
All 9 PEs working in parallel. No data loaded to/from memory during computation.
Data flows through, like blood through capillaries.
```

**Why systolic arrays are efficient:**
1. **Data reuse**: Each weight loaded once, used by all rows; each activation reused across columns
2. **No memory access during computation**: data flows through PEs; only loaded once
3. **Massive parallelism**: all PEs compute simultaneously
4. **Simple PEs**: just multiply-accumulate; no decode, no branch prediction, no caching

**Google TPU v1 systolic array**: 256×256 = 65,536 MAC units operating in parallel. One cycle = 65,536 multiply-accumulates.

### Quick Check
> 1. What does each PE (Processing Element) in a systolic array do?
> 2. Why is the systolic array so efficient compared to CPU matrix multiply?
> 3. How many MAC units does the Google TPU v1 systolic array have?

---

## 3. Google TPU — The Canonical DSA

Google's Tensor Processing Unit (TPU) is the most studied and documented commercial DSA. Google open-sourced significant architectural details in their 2017 paper.

**Why Google built the TPU:**
- 2013: Google Brain team projected that if all Android users used voice recognition for 3 minutes/day, Google would need to double its data center compute capacity
- CPU inference cost was too high — DSA would cost 10–30× less per inference

**TPU v1 (2015, in production; paper 2017):**
```
TPU v1 architecture:
  
  ┌────────────────────────────────────────────────────────┐
  │                   TPU v1                               │
  │                                                        │
  │  Host interface (PCIe)                                 │
  │         ↓                                              │
  │  ┌──────────────┐   ┌────────────────────────┐        │
  │  │  Weight FIFO │   │  Unified buffer (28MB) │        │
  │  │  (DRAM load) │   │  (activation storage)  │        │
  │  └──────┬───────┘   └────────────┬───────────┘        │
  │         │ weights                 │ activations        │
  │  ┌──────▼─────────────────────────▼───────────┐        │
  │  │     Matrix Multiply Unit (MXU)             │        │
  │  │     256×256 systolic array                 │        │
  │  │     65,536 INT8 MACs per cycle             │        │
  │  └───────────────────────┬────────────────────┘        │
  │                          │ partial sums                │
  │  ┌───────────────────────▼────────────────────┐        │
  │  │  Accumulators (4MB)                        │        │
  │  └───────────────────────┬────────────────────┘        │
  │                          │                             │
  │  ┌───────────────────────▼────────────────────┐        │
  │  │  Activation function unit                  │        │
  │  │  (ReLU, sigmoid, tanh, softmax)            │        │
  │  └────────────────────────────────────────────┘        │
  └────────────────────────────────────────────────────────┘

Specs:
  Process: 28nm (older — this is 2015!)
  Chip area: 331mm²
  Power: 40W
  INT8 performance: 92 TOPS
  Memory: 8GB LPDDR3 at 34 GB/s
```

**TPU vs GPU vs CPU (from the 2017 Google paper):**
- TPU ran production inference workloads (search ranking, image recognition, translation)
- Geometric mean performance per chip: TPU 80× faster than CPU, 29× faster than GPU
- Performance per watt: TPU 83× better than CPU, 30× better than GPU

**TPU v4 (2021):**
- 3D torus mesh network connecting 4096 TPU chips per "pod"
- 275 TOPS peak INT8, 137.5 TFLOPS BF16
- Custom interconnect (ICI, Inter-Chip Interconnect) at 1.2 TB/s between chips
- Google uses TPU pods for training LaMDA, PaLM, Gemini

### Quick Check
> 1. What was Google's original motivation for building the TPU?
> 2. What are the main components of the TPU v1 architecture?
> 3. How much faster was TPU v1 than a contemporary CPU per watt?

---

## 4. Neural Processing Units (NPUs) in Mobile SoCs

Every modern flagship smartphone SoC has a Neural Processing Unit — a DSA optimized for neural network inference at low power.

**Apple Neural Engine:**
- A11 Bionic (2017): 2-core, 600 GOPS, first consumer NPU
- A17 Pro (2023): 16-core, 35 TOPS
- M4 (2024): 38 TOPS

```
Apple Neural Engine characteristics:
  
  Architecture: grid of "compute units" (similar to systolic blocks)
  Data type: INT8, INT16, FP16, BF16 mixed precision
  On-chip SRAM: large (~40MB for M4) to avoid DRAM bottleneck
  Integration: unified memory (no copy between CPU and NPU)
  
  Core ML framework: converts PyTorch/TensorFlow models → 
                     Apple's internal format → ANE-optimized binary
  
  Energy: ~1W for ANE at 38 TOPS = 38 TOPS/W
  CPU equivalent: 35W × 38 TOPS/35W = comparable, 35× more power
  
  Use cases:
  - Face ID (depth map, face mesh recognition)
  - Photo computational HDR
  - "On-device LLM" (iPhone 15 Pro: Gemini Nano, 2B params, on ANE)
  - Siri speech processing
  - Portrait mode background segmentation
```

**Qualcomm Hexagon NPU (Snapdragon 8 Gen 3):**
- 45 TOPS
- CDSP (Compute DSP) + Scalar + HVX (Hexagon Vector eXtensions) + HMX (Hexagon Matrix eXtensions)
- HMX: dedicated 8K × 8K INT8/INT16 systolic array
- ONNX runtime + Qualcomm AI Engine Direct SDK

**Samsung/Exynos NPU:**
- Exynos 2400: 34.97 TOPS
- Used in Galaxy S24 (some markets)

**MediaTek APU (AI Processing Unit):**
- Dimensity 9300: 35 TOPS (on-device LLM capable)
- NeuroPilot SDK

**Why mobile NPUs matter now (2024):**
- On-device LLM inference: Gemini Nano (2B parameters) on Pixel 8, iPhone 15 Pro
- Privacy: data never leaves device
- Latency: no network round-trip
- Cost: no server inference cost per query

### Quick Check
> 1. What was the first consumer device with a dedicated NPU and what year?
> 2. Why does on-device LLM inference require a dedicated NPU rather than using the CPU?
> 3. What is the energy efficiency advantage of the Apple Neural Engine over CPU for inference?

---

## 5. Video Codec Accelerators

Video processing is among the highest-volume computational workloads globally — YouTube serves 500 hours of new video per minute, Netflix serves 300 million subscribers. Without hardware acceleration, this would be impossible.

**Why video is a perfect DSA domain:**
- Highly standardized algorithms (H.264, H.265/HEVC, AV1, VP9) with stable specifications
- Repetitive patterns: motion estimation, DCT/IDCT, entropy coding
- Hard real-time requirements (30/60/120fps playback)
- Power sensitivity: decode on mobile without draining battery

**Video compression key operations:**
```
Inter-frame prediction (motion compensation):
  Frame N-1: ┌───────────────────────────────┐
             │    ← moving car →             │
             └───────────────────────────────┘
  Frame N:   ┌───────────────────────────────┐
             │        ← moving car →         │
             └───────────────────────────────┘
  
  Instead of storing Frame N entirely:
  Store "car moved 12 pixels right" + residual difference
  
  Hardware motion estimation: test 64×64 grid of possible offsets
                              for each 16×16 block = 64² × block_count comparisons
  Full HD 1080p: 8100 blocks × 4096 offsets = 33M comparisons per frame
  At 60fps: 2B comparisons/second → need dedicated hardware
  
Transform coding (DCT/IDCT):
  Converts 8×8 pixel blocks to frequency domain
  High-frequency = fine detail (compress aggressively)
  Low-frequency = coarse shape (keep more precisely)
  Hardware DCT: butterfly network, fixed topology, very efficient
```

**Apple Media Engine (M2/M3/M4):**
- Dedicated hardware for: H.264/H.265/ProRes encode + decode, HEVC, ProRes RAW
- Power for 4K H.265 decode: ~0.5W (vs ~15W CPU soft decode)
- M4 can decode 8× 4K ProRes streams simultaneously via Media Engine

**AV1**: Open-source, royalty-free codec (Google/Alliance for Open Media). 50% better compression than H.264. Computationally complex (5–10× more than H.264 to decode in software). AV1 hardware decode is now in every new SoC: Apple M3+, Snapdragon 8 Gen 2+, Intel Arc, AMD RDNA 3+.

### Quick Check
> 1. Why does video decode require a dedicated hardware accelerator rather than software on a CPU?
> 2. What is motion estimation in video compression and why is it compute-intensive?
> 3. What is AV1 and why does it need hardware acceleration?

---

## 6. Network and Storage Accelerators

Data centers process billions of network packets and exabytes of storage I/O. Each operation is simple, but the volume demands dedicated silicon.

**SmartNIC / DPU (Data Processing Unit):**
```
Traditional NIC (Network Interface Card):
  Network packets → NIC chip → PCIe → CPU → process → PCIe → NIC → network
  
  CPU must handle: encryption/decryption, packet routing, firewall rules,
                   RDMA (Remote Direct Memory Access), storage virtualization
  
  For a 100 Gbps NIC: 100M packets/second × 1µs CPU time = 100% of one CPU core
  
DPU (Data Processing Unit): NVIDIA BlueField-3, AMD Pensando, Intel IPU
  Network packets → DPU (has its own ARM cores + ASIC accelerators)
  DPU handles: TLS encryption, vSwitch, storage stack
  CPU gets: clean payload, high-level events only
  
  Result: 80% reduction in host CPU cycles for network operations
          frees CPUs for application work
          
BlueField-3 specs:
  16× ARM Cortex-A78 cores (for programmable data plane)
  400 Gbps networking
  Hardware crypto (AES-256-GCM at 400 Gbps)
  RDMA, RoCE v2 offload
  ~65W total
```

**Storage accelerators:**
- NVMe SSDs: controller chip (DRAM-less: uses PCM or SLC cache) — ARM Cortex M-class DSP handling flash management, garbage collection, wear leveling
- Computational storage (CSDs): Seagate Lyve, Samsung SmartSSD — FPGA or ARM inside the SSD. Filter/aggregate data before sending to host.
  - Example: SELECT WHERE country='India' → scan 1TB SSD → normally sends 1TB over PCIe → computational storage executes predicate in SSD, sends 50MB back
- Compression: Intel QAT, NVIDIA DEFLATE — hardware gzip/zstd at 100+ GB/s

### Quick Check
> 1. What is a DPU and what CPU work does it offload?
> 2. What is computational storage and how does it reduce data movement?
> 3. What tasks does an NVMe SSD controller chip perform?

---

## 7. Crypto and Security Accelerators

Cryptographic operations are mathematically intensive and bottleneck many systems — TLS connections, VPNs, blockchain, secure boot, password hashing.

**AES-NI (Intel/AMD/ARM, 2010):**
- Hardware AES encryption/decryption in the CPU itself
- AES-256-GCM: 1 instruction per 16 bytes on AES-NI vs 100+ clock cycles on software
- Performance: 50 GB/s single-core AES vs 0.5 GB/s without AES-NI
- Used for: HTTPS, disk encryption, VPN

**SHA extensions (SHA-NI, SHA-2/SHA-3):**
- Similar idea for hash functions (used in TLS, certificates, signatures)

**PQC accelerators (SHAKTI and others):**
As covered in Chapter 65, SHAKTI includes:
- NTT (Number Theoretic Transform) accelerator for CRYSTALS-Kyber: 111× speedup
- Similar for Dilithium signature verification

**Bitcoin ASIC (Bitmain Antminer):**
- SHA-256 hashing in double-loop: hash(hash(block header + nonce)) < target
- CPU: ~10 MH/s (megahashes/second)
- GPU: ~500 MH/s
- Custom ASIC (TSMC 5nm): ~100 TH/s = 100,000,000 MH/s
- Efficiency: 10 billion× improvement from CPU to dedicated ASIC over 12 years
- The extreme example of what specialization can achieve

```
Bitcoin mining ASIC vs CPU (SHA-256):

CPU Intel i9:       ~10 MH/s      at 125W = 0.08 MH/J
GPU NVIDIA 4090:   ~500 MH/s      at 450W = 1.1 MH/J
ASIC Antminer S21: 200,000,000 MH/s at 3500W = 57,143 MH/J

ASIC vs CPU: ~700,000× more energy-efficient per hash
```

### Quick Check
> 1. What is AES-NI and how does it compare to software AES?
> 2. Why are Bitcoin mining ASICs so many orders of magnitude more efficient than CPUs?
> 3. What is the NTT accelerator in SHAKTI optimized for?

---

## 8. How to Design a DSA

Designing a DSA is a structured process. Here is the engineering approach:

**Step 1: Profile the workload**
- Identify the 20% of operations consuming 80% of compute and energy
- Use roofline analysis: is the bottleneck compute or memory bandwidth?
- Example finding: transformer inference → 70% of time in matrix multiply, 15% in softmax

**Step 2: Identify algorithmic structure**
- Matrix multiply: regular, data-independent access pattern → systolic array
- Graph traversal: irregular → hard to accelerate, maybe near-memory compute
- AES: fixed 10-round key schedule → unroll all rounds into pipeline

**Step 3: Design the data path**
- Precision: does INT8 work? (saves 2× area vs INT16, 4× vs FP32)
- Data flow: weight-stationary vs output-stationary vs row-stationary systolic array
- Memory hierarchy: how much SRAM needed to eliminate DRAM stalls?

**Step 4: Choose the implementation form**
- Pure ASIC: highest efficiency, lowest flexibility, high NRE (millions of $)
- FPGA: flexible, moderate efficiency, can update, high power
- DSP array with custom ISA (Hexagon): flexible within domain, middle ground

**Step 5: Write the software stack**
- Compiler: converts high-level graph (ONNX, TFLite) to hardware instructions
- Runtime: manages memory, schedules kernels, handles batching
- Hardware abstraction: programmers shouldn't write assembly for the DSA

**Step 6: Validate and iterate**
- The first version will have bottlenecks you didn't anticipate
- Post-silicon profiling: add performance counters (like a CPU PMU but for DSA)

**Common DSA pitfalls:**
- Optimizing for yesterday's algorithm (the neural network architecture you design for today may not be the one that matters in 3 years)
- Underestimating memory bandwidth requirements
- Forgetting about data movement cost (99% of energy can be in moving data, not computing)
- Inflexible enough to handle model updates (quantization changes, new activation functions)

### Quick Check
> 1. What is the first step in designing a DSA?
> 2. Why is precision (INT8 vs FP32) so important in DSA design?
> 3. Why must a DSA have a software stack (compiler + runtime)?

---

## Summary

- **DSA definition**: specialized processor designed for one workload domain; trades generality for 10–10,000× efficiency.
- **Why now**: end of Dennard Scaling, workload consolidation in cloud, stable ML algorithms, energy dominates cost.
- **Systolic array**: grid of simple MACs; data flows rhythmically. Foundation of all AI accelerators. 65,536 MACs/cycle in TPU v1.
- **Google TPU**: 80× faster than CPU, 30× faster than GPU per watt. Systolic array + unified buffer + host via PCIe. TPU v4 pods: 4096 chips, 3D torus.
- **Mobile NPUs**: Apple ANE (38 TOPS), Qualcomm Hexagon (45 TOPS), MediaTek APU (35 TOPS). Enable on-device LLM. Privacy + latency + cost.
- **Video codecs**: H.264/HEVC/AV1 hardware decode. 30× less power than software. 500 hours/min YouTube throughput impossible without it.
- **DPU**: Offloads network, crypto, storage from CPU. 80% CPU cycle savings for networking in data centers.
- **Crypto**: AES-NI (50 GB/s), Bitcoin ASIC (700,000× CPU efficiency), SHAKTI PQC (111× speedup for Kyber).
- **DSA design process**: profile → algorithm structure → datapath → implementation → software stack → validate.

---

## Exercises

### Easy
1. What is the energy efficiency advantage of a DSA over a CPU for the same workload?
2. What is a systolic array and what operation is it optimized for?
3. What is Apple's Neural Engine and what workloads does it run?

### Medium
4. Systolic array performance calculation: A 256×256 systolic array (like Google TPU v1) runs at 700 MHz. (a) MACs per second = 256² × 700M = ? TOPS. (b) For a convolution layer: 3×3 kernel, 256 input channels, 256 output channels, 224×224 input image (after im2col reshaping, becomes matrix multiply of shape [50,176, 2304] × [2304, 256]). Total MACs for this layer? (c) At 65,536 MACs/cycle and 700 MHz: cycles for this layer? Wall-clock time? (d) Memory: weight matrix [2304, 256] × 1 byte (INT8) = 590KB. Fits in 28MB unified buffer? (e) What if INT4 is used instead of INT8: same storage, double compute per MAC unit → effective throughput?
5. DPU vs CPU for TLS: A web server handles 50,000 HTTPS requests/second. Each TLS handshake: (a) RSA-4096 signature verification: 10ms on CPU × 50,000/second = ? CPU cores? (b) AES-256-GCM encryption: 1 GB/s per CPU core, average 64KB per request, 50,000 requests/s = 3.2 GB/s → ? CPU cores? (c) Total CPU cores for crypto alone at full load? (d) An Intel QAT accelerator card: 100 Gbps AES, 1M RSA ops/second, $800. What CPU core count does it replace? At $3,000/CPU: what is the CapEx savings? (e) Power: CPU core at 15W; QAT card at 50W. How many cores is QAT replacing in terms of power?
6. TPU cost model: Google deploys TPU v4 for Gemini inference. (a) At 275 TOPS/chip, 10M inference requests/second, average 100M operations each request: TOPS needed? How many TPU v4 chips? (b) TPU v4 cost estimated at $15,000 per chip. CapEx for this fleet? (c) Power: 175W per TPU v4 chip. Annual energy cost at $0.04/kWh for data center power? (d) NVIDIA H100: 3.9 PFLOPS FP16 (note: INT8 = 2× = 7.8 TOPS). Same TOPS requirement: how many H100s? CapEx at $30,000/H100? Annual energy at 700W each? (e) Conclusion: for which metric does TPU v4 win? GPU? When does the higher H100 compute matter?

### Hard
7. Design a Transformers inference DSA: You are designing a chip specifically for LLM inference (decoder-only transformer). Target: GPT-3 175B at 100 tokens/second, single chip. (a) Workload analysis: attention (Q×K^T and score×V) vs linear layers (FFN) breakdown — attention is 30% of FLOPs, FFN is 70%. Arithmetic intensity: attention O(n²)/O(n²d) FLOPs/bytes ≈ low (1–5 FLOPs/byte); FFN O(nd²)/O(nd) FLOPs/bytes ≈ medium-high. What does this imply for the design? (b) Memory: 175B parameters × 2 bytes (FP16) = 350GB. Your chip can have: Option A (512GB HBM3 at 3 TB/s) or Option B (1TB LPDDR5x at 400 GB/s). Which is better and why? (c) Compute: at batch size 1, attention is memory-bandwidth-bound. At 3 TB/s bandwidth and 350GB model: minimum time to do one forward pass? What batch size makes it compute-bound? (d) Core design: you have 200mm² at TSMC N3. Use: 40% for systolic array (1000 MACs/mm² at N3), 40% for SRAM (10 MB/mm² at N3), 20% for misc. How many TOPS? How much SRAM? Is it enough to avoid going to DRAM for KV cache (for 4K context length: 2×2×d×L×batch×2 bytes)? (e) Compiler: describe the 5 passes a compiler must do to convert a PyTorch model to run on this chip.
8. Bitcoin ASIC business case (2024): Bitmain Antminer S21 Pro: 234 TH/s SHA-256, at 3531W, TSMC 4nm, unit cost ~$3,000. (a) Bitcoin network hashrate: ~700 EH/s (exahashes/second). How many S21 Pro units make up the entire network? Total power? (b) Mining revenue: 3.125 BTC/block × 144 blocks/day = 450 BTC/day split among all miners. At $70,000/BTC: daily revenue for entire network? For one S21 Pro with 234 TH/s / 700 EH/s share? (c) Operating cost: one S21 Pro at 3531W × 24h = 84.7 kWh/day. At $0.05/kWh industrial rate: electricity cost? Daily profit? Payback period? (d) Why is this a perfect example of Dennard Scaling reversal: each generation of ASICs improves efficiency but operating cost (electricity) becomes the dominant expense rather than hardware cost? (e) In 2020, SHA-256 ASIC performance was ~110 TH/s at 3250W (Antminer S19). S21 Pro: 234 TH/s at 3531W. Compute the improvement in TH/J. Is this consistent with process scaling from 7nm to 4nm? What would you expect at TSMC 2nm?
