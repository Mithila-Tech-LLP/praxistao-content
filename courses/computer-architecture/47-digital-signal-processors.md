# Chapter 47: Digital Signal Processors

A DSP (Digital Signal Processor) is a specialized processor optimized for processing streams of sampled signals — audio, radio frequency, radar, sonar, image processing, telecommunications. While a general-purpose CPU can do signal processing, a DSP does it 10–100× more efficiently for the same power budget. DSPs are everywhere you don't see them: inside your phone's noise-cancelling headphones, in cellular base stations, in your car's radar, in hearing aids, in satellite receivers. Understanding DSPs bridges the gap between the arithmetic world of CPUs and the analog world of the physical signals they process.

## Table of Contents

1. [What Makes Signal Processing Special?](#1-what-makes-signal-processing-special)
2. [The Multiply-Accumulate (MAC) Operation](#2-the-multiply-accumulate-mac-operation)
3. [DSP Architecture](#3-dsp-architecture)
4. [Key DSP Algorithms](#4-key-dsp-algorithms)
5. [Texas Instruments TMS320 and C6000](#5-texas-instruments-tms320-and-c6000)
6. [Qualcomm Hexagon — The Mobile DSP](#6-qualcomm-hexagon--the-mobile-dsp)
7. [DSP vs CPU vs GPU for Signal Processing](#7-dsp-vs-cpu-vs-gpu-for-signal-processing)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Makes Signal Processing Special?

**Sampling**: The physical world produces continuous (analog) signals — sound pressure waves, electromagnetic fields, temperature. To process them digitally, we sample at regular intervals (e.g., audio at 44,100 Hz = 44,100 samples per second). The Nyquist theorem states we must sample at least twice the highest frequency component in the signal.

```
Audio signal processing pipeline:
  Microphone → ADC → [DSP processes digital samples] → DAC → Speaker
  
  At 48 kHz stereo: 96,000 samples/second
  With 16ms latency target: 1536 samples must be processed in 16ms
  Processing per sample: noise reduction, equalization, echo cancellation, compression
```

**Stream processing**: Unlike a database query (random access, variable length), signal processing is a continuous **stream** of fixed-size samples. Processing one sample requires a fixed number of operations. This regularity is what DSPs exploit.

**Key operations**:
- **Filtering**: Multiply each sample by coefficients, sum products → FIR/IIR filters
- **Correlation**: Compute similarity between signals → used in radar, CDMA
- **FFT (Fast Fourier Transform)**: Convert time-domain to frequency-domain → used everywhere
- **Convolution**: General operation underlying filtering, correlation, neural networks

### Quick Check
> 1. What is sampling and why is the Nyquist rate the minimum sampling frequency?
> 2. Why is signal processing a "stream" rather than random-access computation?
> 3. Name three applications where DSPs are used.

---

## 2. The Multiply-Accumulate (MAC) Operation

The core of almost every signal processing algorithm is the **MAC (Multiply-Accumulate)**:
```
accumulator += a[i] * b[i]
```

In FIR filter with N coefficients:
```
y[n] = h[0]×x[n] + h[1]×x[n-1] + h[2]×x[n-2] + ... + h[N-1]×x[n-N+1]
     = Σ h[k] × x[n-k]    for k = 0 to N-1
```

This is N multiply-accumulate operations to produce each output sample. At 48 kHz audio with a 256-tap filter: 256 × 48,000 = 12.3 million MACs per second.

**Why MACs are hard for general CPUs:**
1. Standard CPU: multiply and add are separate instructions, 2 ops + 2 cycles
2. DSP: MAC is a single instruction executing in 1 cycle
3. Standard CPU: accumulator overflow is a problem with fixed-point arithmetic
4. DSP: extended-precision accumulator prevents overflow during accumulation

**Fixed-point vs floating-point DSP arithmetic**:
- **Fixed-point**: integers with implicit decimal point (Q15: 16-bit signed with 15 fractional bits → range -1 to +1). Cheap hardware, needs careful programmer management of scale.
- **Floating-point**: IEEE 754, automatic range management. More expensive. Modern DSPs often support both.

### Quick Check
> 1. Write out the MAC operation mathematically and explain each term.
> 2. A 256-tap FIR filter at 48 kHz: how many MACs per second?
> 3. What is a Q15 fixed-point number? What value does 0x4000 represent in Q15?

---

## 3. DSP Architecture

DSPs include specialized hardware for the operations signal processing demands:

**Single-cycle MAC unit**: Dedicated multiplier + accumulator that completes a full MAC in 1 clock cycle. Not possible in standard ALUs without specific design.

**Dual Harvard memory**: Two separate data memory buses. A filter computation needs two operands per cycle — one coefficient h[k] and one data sample x[n-k]. With a single data bus, this requires 2 cycles. With dual data buses, both are fetched simultaneously.

```
Super-Harvard DSP memory architecture:
  Instruction memory bus → reads next instruction
  Data memory bus A      → reads coefficient h[k]
  Data memory bus B      → reads sample x[n-k]
  
  Three simultaneous memory accesses per cycle!
```

**Circular buffers in hardware**: Sample buffers are circular (new samples replace oldest). Hardware-managed circular addressing eliminates the software modulo operation that would be needed to wrap around a circular buffer.

**Zero-overhead loops**: A hardware loop counter decrements and branches automatically without consuming the CPU's branch instruction slots or pipeline bubbles.

**Bit-reversal addressing**: FFT algorithms access memory in bit-reversed order (input index 3 = 011 → bit-reverse → 110 = 6). Hardware bit-reverse addressing reduces FFT overhead.

**VLIW (Very Long Instruction Word)**: Many DSPs use VLIW — packing multiple operations into one wide instruction word. The compiler fills all slots for maximum parallelism. Simpler hardware than OOO but requires smart compiler.

```
TI C62x VLIW instruction word (256 bits!):
  [ALU op A][ALU op B][Load A][Load B][FP op][Mult][Branch][nop]
  All execute in parallel in one cycle
```

### Quick Check
> 1. Why does a DSP need two data memory buses instead of one?
> 2. What is a "circular buffer" and why does DSP hardware support it natively?
> 3. What is "bit-reversal addressing" and which algorithm requires it?

---

## 4. Key DSP Algorithms

**FIR (Finite Impulse Response) Filter**: Each output is a weighted sum of the last N input samples. Coefficients h[k] are fixed (or slowly updated). Always stable. Phase linear. Used for: equalizers, noise reduction, anti-aliasing.

**IIR (Infinite Impulse Response) Filter**: Output depends on both past inputs and past outputs — feedback. Fewer multiplies than FIR for same frequency selectivity. Can be unstable if not designed carefully. Used for: echo cancellation, audio EQ.

**FFT (Fast Fourier Transform)**: Computes the frequency content of a signal. An N-point FFT requires O(N log N) multiplies vs O(N²) for the naive DFT. Used for: spectrum analysis, OFDM (WiFi, LTE, 5G modulation), noise reduction, pitch detection.

**OFDM (Orthogonal Frequency Division Multiplexing)**: Divides a wide channel into many narrow subcarriers. Used in WiFi (802.11a/g/n/ac/ax) and LTE/5G. Requires FFT/IFFT at both transmitter and receiver — DSP's killer app in telecommunications.

**Convolution**: The fundamental operation. Neural network convolutions, digital filters, matched filtering in radar — all reduce to MAC operations. GPUs are good at convolutions too, which is why GPUs became ML training hardware.

**Correlation**: Measure how similar two signals are. In CDMA, each user has a unique spreading code; the receiver correlates the received signal with each code to decode each user's data.

### Quick Check
> 1. What is the difference between FIR and IIR filters? When would you use each?
> 2. Why is the FFT so much faster than the DFT for large N?
> 3. What is OFDM and why does it need FFT hardware?

---

## 5. Texas Instruments TMS320 and C6000

**Texas Instruments (TI)** is the dominant DSP manufacturer with decades of leadership.

**TMS320 family history:**
- TMS320C10 (1982): First TI DSP, 5 MIPS
- TMS320C50 (1993): Fixed-point, 100 MIPS, telecom applications
- TMS320C55x: Ultra-low power, used in MP3 players, cell phones (pre-smartphone era)
- TMS320C62x/C67x (1997): First VLIW DSP — C6000 architecture

**C6000 architecture** (C62x/C64x/C67x/C6748):
```
C6748 block diagram:
  2× VLIW data paths (A and B side)
  Each path: 4 functional units
    .L1/.L2: ALU (logic, arithmetic, long operations)
    .S1/.S2: Shifter/branch
    .M1/.M2: Multiplier (32×32-bit, 16×16-bit, SIMD)
    .D1/.D2: Data load/store
  
  8 functional units total per cycle
  VLIW: 256-bit instruction packet (8 × 32-bit operations)
  Peak: C6748 at 456 MHz = 3648 MIPS (8 operations × 456 MHz)
```

**C6000 toolchain**: TI's Code Composer Studio (CCX) compiler is specifically designed to maximize VLIW instruction packing. Good C code compiled with pragmas (`#pragma MUST_ITERATE`) can approach peak theoretical throughput.

**Current TI DSPs:**
- TMS320C66x: Floating-point + SIMD, used in radar, medical imaging
- C7x (KeyStone 3): 512-bit vector DSP, 40 GFLOPS, used in automotive ADAS
- AM243x, AM64x: DSP + ARM Cortex-A53 SoC for industrial/automotive

### Quick Check
> 1. What is VLIW and how does the C6000 architecture implement it?
> 2. A C6748 running at 456 MHz with 8 functional units: what is the peak theoretical MIPS?
> 3. What type of application uses TI C6000 DSPs today?

---

## 6. Qualcomm Hexagon — The Mobile DSP

Chapter 37 introduced the Hexagon briefly; here is more architectural depth.

**Hexagon architecture** (HVX DSP):
- VLIW processor, up to 4 operations per packet
- HVX (Hexagon Vector eXtensions): 128-byte (1024-bit) SIMD vectors — extremely wide
- Optimized for 8-bit operations (for CNN inference)
- Supports **threads**: 4 hardware threads for pipeline hiding of memory latency
- Separate L1 instruction and data caches; 32KB each

**HVX vector operations:**
```
Example: 8-bit convolutional filter on 1920×1080 image
  Each HVX operation processes 128 bytes = 128 pixels simultaneously
  One instruction computes 128 multiply-accumulate operations
  At 1 GHz: 128 billion 8-bit MACs/sec = 128 GOPS
  (Compare to Cortex-A76 with 128-bit NEON: 16 pixels/instruction)
```

**Always-on Hexagon**: A low-power Hexagon core (or micro-DSP) can run at <10 mW for always-on voice, motion detection, and sensor fusion. The main Hexagon core runs at 1GHz+ for camera and AI.

**Hexagon SDK**: Qualcomm provides an SDK for writing custom DSP code for Hexagon. Used by camera apps, audio vendors, and AI framework developers to write optimized kernels that run on Hexagon instead of the CPU.

### Quick Check
> 1. What is HVX and how many bytes does a single HVX vector operation process?
> 2. Why does the always-on Hexagon run at <10mW?
> 3. Why might a camera application developer write a custom Hexagon kernel?

---

## 7. DSP vs CPU vs GPU for Signal Processing

```
Comparison for audio/signal processing workloads:

                CPU             DSP             GPU
                (Cortex-A77)    (C6748)         (Adreno 750)
------------------------------------------------------------------------
MAC throughput  ~10 GOPS (NEON) ~3.6 GOPS       ~75 TOPS (INT8)
Power           4–8W             ~2W             15W
Latency         Low              Very low        Higher (setup overhead)
Programmability High             Medium          High
Cost            High             Medium          High
Real-time       Good             Excellent       Poor (irregular latency)
Code size       Small            Medium          Large (shader overhead)

Winner for:
  Audio filters: DSP (deterministic timing, low latency)
  OFDM modem:   DSP (bit-by-bit stream processing, precise timing)
  Image ML:     GPU/NPU (batch processing, not real-time stream)
  Radar FFT:    DSP (TI C66x specialized for radar)
  Video codec:  Dedicated hardware (VPU) or GPU
```

**Real-time requirements** favor DSPs:
- Audio processing must complete within a few milliseconds or there is audible glitching
- OFDM demodulation must process each symbol within the symbol period (~70µs for LTE)
- DSP architectures provide deterministic timing guarantees impossible in OOO CPUs

**Throughput requirements** favor GPUs/NPUs:
- A 100-layer CNN inference on a single image
- Training a neural network over millions of examples
- These can tolerate some latency variance (they're batch, not stream)

### Quick Check
> 1. For processing real-time audio (16ms latency requirement), why is a DSP better than a GPU?
> 2. What is the key performance metric for DSP workloads that differs from GPU workloads?
> 3. What is a VPU and when would you use it over a DSP?

---

## Summary

- **DSPs** are specialized processors for streaming signal processing — optimized for the MAC operation, fixed-function loops, and high memory bandwidth.
- Core features: single-cycle MAC, dual data buses (two operands simultaneously), hardware circular buffers, zero-overhead loops, VLIW instruction packing.
- Key algorithms: FIR/IIR filters, FFT, OFDM modulation/demodulation, convolution, correlation.
- **Texas Instruments C6000**: VLIW DSP architecture, 8 functional units per cycle, used in radar, medical imaging, base stations.
- **Qualcomm Hexagon**: Mobile DSP with 1024-bit HVX vector extension, used for camera, audio, and AI inference in Snapdragon.
- DSPs win over CPUs for real-time stream processing (audio, modem); GPUs win for batch ML inference; NPUs win for neural network inference specifically.

---

## Exercises

### Easy
1. What is a MAC operation? Write it as a mathematical formula.
2. Why does a DSP need two separate data memory buses when a CPU only has one?
3. Give three real-world applications where a DSP is used instead of a general-purpose CPU.

### Medium
4. FIR filter design: A 256-tap FIR filter runs at 48 kHz audio rate. (a) How many MACs per second? (b) A TI C6748 at 456 MHz with 2 MACs/cycle: how much CPU time does this filter use? (c) The same filter on an ARM Cortex-A53 at 1.5 GHz doing a MAC in 3 cycles: how much CPU time?
5. OFDM timing constraint: LTE uses 2048-point FFT, one OFDM symbol every 71.4µs. A DSP must complete the FFT computation within this window. A TI C66x at 1 GHz does 2048-point FFT in 10,000 cycles. (a) How long does the FFT take? (b) What fraction of the OFDM symbol period is used? (c) If you add 16QAM demodulation (additional 2048 operations), can you still meet the deadline?
6. Fixed-point vs floating-point audio processing: A reverb algorithm requires accumulating 512 samples with gains ranging from 0.001 to 0.999. In Q15 format (15 fractional bits, max value 1.0): (a) Can Q15 represent 0.001 accurately? (b) What is the quantization error? (c) During accumulation of 512 terms, by how many bits does the accumulator need to grow to avoid overflow? (d) Why might you prefer a 40-bit accumulator with 32-bit Q15 inputs?

### Hard
7. Cellular base station DSP: A 5G NR base station handles 100 simultaneous users, each sending 30.72 MHz bandwidth signals (numerology µ=3). Each sub-carrier requires an IFFT + OFDM modulation + beam-forming + channel coding. Estimate: (a) FFT size and rate for this bandwidth, (b) compute requirements per user per millisecond, (c) total compute for 100 users, (d) whether a TI C7x (200 GFLOPS FP32) or NVIDIA A30 GPU (10 TFLOPS FP32) is more appropriate, considering latency constraints.
8. Hardware vs software comparison: A software-defined radio (SDR) receiver written in Python on a laptop can demodulate FM radio in real-time, but cannot do LTE in real-time (too slow). An RTL-SDR dongle + GNSS-SDR library can do GPS in near-real-time with a fast CPU. Analyze: (a) what makes FM demodulatable in software but LTE not? (b) what FLOPS/GOPS are required for each? (c) how does a dedicated LTE modem chip (like Qualcomm SDR895) achieve real-time LTE processing that software cannot? (d) what will happen to DSPs as GPUs and NPUs become more powerful and power-efficient?
