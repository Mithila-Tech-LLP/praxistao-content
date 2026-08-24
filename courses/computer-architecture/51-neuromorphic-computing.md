# Chapter 51: Neuromorphic Computing — Computing Like a Brain

The human brain performs remarkable feats — recognizing faces, understanding language, navigating unfamiliar environments — using about 20 watts of power. A GPU doing equivalent pattern recognition tasks uses 300–700 watts — orders of magnitude less efficient per task. The brain does this with biological neurons and synapses that operate fundamentally differently from transistors and logic gates. **Neuromorphic computing** attempts to build hardware that mimics the brain's architecture: asynchronous spiking neurons, event-driven processing, in-memory computation, and massive parallelism. This chapter explains the neuroscience inspiration, the hardware architectures, the software challenges, and why neuromorphic computing may become important for ultra-low-power AI at the edge.

## Table of Contents

1. [How the Brain Computes](#1-how-the-brain-computes)
2. [Spiking Neural Networks — The Brain's Language](#2-spiking-neural-networks--the-brains-language)
3. [Neuromorphic Hardware](#3-neuromorphic-hardware)
4. [Intel Loihi and IBM TrueNorth](#4-intel-loihi-and-ibm-truenorth)
5. [Programming Neuromorphic Systems](#5-programming-neuromorphic-systems)
6. [Neuromorphic vs Conventional AI](#6-neuromorphic-vs-conventional-ai)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. How the Brain Computes

The brain is made of approximately 86 billion neurons. Each neuron is a simple processing element that receives electrical signals from other neurons and, if the combined signal exceeds a threshold, "fires" — generating an electrical spike that travels along its axon to other neurons.

```
Single neuron (simplified):
  
  Input dendrites         Cell body (soma)       Output axon
  ─────────────────────► ┌──────────────────┐ ──────────────────►
  signals from 1000s     │  Integrate inputs │    spike if
  of upstream neurons    │  If sum > threshold│ threshold reached
                         │  → fire!           │
                         └──────────────────┘
                              │
                              └── Synaptic connections to
                                  1000s of downstream neurons
```

**Key differences from digital computers:**

1. **Asynchronous**: Neurons don't tick on a global clock. They fire when they have reason to fire. No synchronization overhead.

2. **Event-driven (sparse)**: Neurons are silent most of the time. A typical neuron fires at ~1–10 Hz on average out of a possible 1000 Hz maximum. Most computation happens in spikes — the silence is computationally free.

3. **Analog computation**: Membrane potential (the electrical charge building up) is a continuous analog value. The spike is digital (all-or-none), but the "when" encodes information.

4. **In-memory computing**: Computation and memory are not separated. Synapses (connections between neurons) store learned weights AND perform computation simultaneously. There is no Von Neumann bottleneck.

5. **Massively parallel**: 86 billion neurons × 7000 synapses each = ~6×10¹⁴ synaptic operations per second — at 20 watts.

```
Brain power efficiency:
  Brain: 20W for ~10^15 synaptic ops/sec = 5×10^13 SOPS/W
  GPU (H100): 700W for ~4×10^15 FP8 ops/sec = 5.7×10^12 OPS/W
  
  So the brain is ~10× more efficient per op — but even this
  comparison is misleading:
  GPU ops are full multiplications.
        Brain synaptic ops are simpler (threshold + increment).
        AND: Brain handles recognition tasks no GPU can do at equivalent power.
  
  The right comparison is task performance per watt, not raw ops.
```

### Quick Check
> 1. What is a neuron and how does it "compute"?
> 2. How is brain computation different from a CPU's clock-driven pipeline?
> 3. What is "in-memory computing" and why does it matter?

---

## 2. Spiking Neural Networks — The Brain's Language

Standard artificial neural networks (ANNs — used in GPT, image recognition) use floating-point numbers to represent activations: "this feature detector is 73.4% activated." A spiking neural network (SNN) uses **spikes**: discrete all-or-nothing events in time.

```
Standard ANN (rate-coded):
  neuron output = floating-point number (e.g., 0.734)
  
  Compute: output = activation_fn(Σ weight[i] × input[i])
  Hardware: matrix multiplication (what GPUs are built for)

Spiking Neural Network (spike-coded):
  neuron output = sequence of spike events: [t=1.2ms, t=5.7ms, t=12.1ms]
  
  Information encoding options:
    Rate coding: spike frequency = activation intensity
      (10 Hz = 10% activated, 100 Hz = 100% activated)
    Temporal coding: exact spike timing encodes information
    Population coding: which neurons fire together encodes a "pattern"
  
  Compute: accumulate spikes (integer adds), fire if threshold reached
  Hardware: needs efficient sparse event processing
```

**Leaky Integrate-and-Fire (LIF) neuron** — the most common SNN model:
```
Membrane potential update:
  V(t+1) = V(t) × (1 - 1/τ) + Σ w[i] × spike[i](t)
  
  If V > threshold: fire a spike, reset V to V_reset
  
  τ = time constant (how fast membrane potential "leaks" back to rest)
  w[i] = synaptic weight
  spike[i] = 1 if presynaptic neuron fired, 0 otherwise
  
This is: integer accumulation (sparse) + threshold compare
Not: floating-point multiply-accumulate (dense)
```

**STDP (Spike-Timing Dependent Plasticity)**: The brain's learning rule. If neuron A fires just before neuron B, strengthen the A→B synapse (Hebbian learning: "neurons that fire together, wire together"). If B fires before A, weaken the synapse. This is an unsupervised, local learning rule — no backpropagation required.

```
STDP rule:
  Δw = A_+ × exp(-Δt / τ_+)  if Δt > 0  (pre fires before post → strengthen)
  Δw = -A_- × exp(Δt / τ_-)  if Δt < 0  (post fires before pre → weaken)
  
  where Δt = t_post - t_pre
```

### Quick Check
> 1. What is the difference between how an ANN encodes information vs an SNN?
> 2. What is STDP and how does it relate to "neurons that fire together, wire together"?
> 3. Why are SNNs potentially more power-efficient than ANNs on dedicated hardware?

---

## 3. Neuromorphic Hardware

The challenge: silicon transistors are not naturally suited to implement asynchronous spiking neurons efficiently. Building neuromorphic hardware requires new circuit architectures.

**Key design principles for neuromorphic chips:**

1. **On-chip memory co-located with computation**: Synaptic weights stored in SRAM or analog memristors right next to the neuron circuits. No Von Neumann bottleneck.

2. **Event-driven circuits**: Circuits consume power only when a spike arrives. Idle neurons consume near-zero power.

3. **Massive parallelism**: Many small neuron/synapse circuits in parallel, not a few powerful cores.

4. **Analog or mixed-signal options**: Some neuromorphic chips use analog circuits for the membrane potential (mimicking analog membrane voltage), reducing power further.

**Memristors** (memory-resistors): Two-terminal devices whose resistance changes based on past current — a physical analog of synaptic weight. Research device; not in commercial production yet. Promising for analog in-memory computing.

```
Crossbar array with memristors (research concept):
  
       Word lines (neuron outputs) →
  ─────────●─────────●─────────●───── bit line 0 (sums)
            │         │         │
  ─────────●─────────●─────────●───── bit line 1
            │         │         │
           R₀₀       R₁₀       R₂₀
  
  Each node is a memristor with resistance R_ij = 1/w_ij
  Current through each column = Σ V_i / R_ij = Σ V_i × w_ij (Ohm's law!)
  
  → Analog matrix multiply in O(1) time, using no digital arithmetic
```

### Quick Check
> 1. What is the Von Neumann bottleneck and how does neuromorphic co-located memory solve it?
> 2. What is a memristor and what computational function does it perform?
> 3. What is "event-driven computing" and why does it save power?

---

## 4. Intel Loihi and IBM TrueNorth

**IBM TrueNorth (2014)**:
- 4096 neurosynaptic cores per chip
- 1 million neurons, 256 million synapses
- No fast global clock — cores are event-driven and asynchronous, synchronized only by a slow 1 kHz global time step
- ~70 mW typical — 46 billion synaptic operations per second per watt
- Each core: 256 neurons × 256 synaptic inputs per neuron
- Fabricated at Samsung 28nm, 5.4 billion transistors in 430 mm²
- Application: pattern recognition in 1 mW power envelope (IoT sensor classification)

```
TrueNorth chip:
  70mm² die, 4096 neurosynaptic cores arranged in 64×64 grid
  Each core: 256 neurons, 256×256 = 65,536 synaptic bits (SRAM)
  Cores communicate via spike events on mesh network (asynchronous)
  
  A spike event: (destination_neuron_id, time_step)
  No floating point — all operations are binary (+1/0 per spike) and threshold compare
```

**Intel Loihi (2018) and Loihi 2 (2021)**:
- Loihi: 128 neuromorphic cores per chip, 130,000 neurons, 130 million synapses
- Loihi 2: 128 cores, 1 million neurons, Intel 4 (7nm-class) process
- Supports on-chip learning (STDP implemented in hardware)
- Programmable neuron models (can implement LIF, adaptive threshold, etc.)
- Intel Hala Point (2024): 1152 Loihi 2 chips, 1.15 billion neurons — largest neuromorphic system
- Power: ~0.1–1 mW per core for event-driven workloads

**Comparison:**
```
                  TrueNorth      Loihi 2        GPU (RTX 4090)
Neurons           1M             1M             N/A (not spiking)
Power @ inference 70mW           <1W            450W
Programmability   Fixed LIF      Flexible        General
Learning on-chip  No             Yes (STDP)      Via backprop
Precision         Binary         Integer         FP32/FP8
ANN accuracy      Lower          Lower           High
```

### Quick Check
> 1. How many neurons does TrueNorth implement and at what power?
> 2. What advantage does Loihi 2 have over TrueNorth for learning?
> 3. Why is the power comparison between neuromorphic chips and GPUs unfair if you compare raw operations?

---

## 5. Programming Neuromorphic Systems

Programming neuromorphic chips is very different from programming CPUs or GPUs. You define networks of spiking neurons and synapses, not sequential algorithms.

**Intel's Lava framework** (open-source):
- Python-based programming model for Loihi
- Define Process objects (neurons) and Channel objects (synaptic connections)
- Compile to neuromorphic hardware or simulate on CPU

```python
# Lava: simple spiking neuron network
from lava.proc.lif.process import LIF
from lava.proc.dense.process import Dense
from lava.proc.io.source import RingBuffer

# Input: spike generator
spike_gen = RingBuffer(data=spike_data)  # preloaded spike patterns

# Dense synaptic layer: 100 input neurons → 10 output neurons
synapses = Dense(weights=np.random.rand(10, 100))

# 10 output LIF neurons
neurons = LIF(shape=(10,), vth=10, dv=0.1, du=0.1)

# Connect
spike_gen.s_out.connect(synapses.s_in)
synapses.a_out.connect(neurons.a_in)
```

**NxSDK** (Intel's lower-level Loihi SDK):
- Direct neuron/synapse configuration
- Hardware-specific optimization

**Converting ANNs to SNNs**: Researchers convert trained PyTorch/TensorFlow models to SNNs by replacing ReLU activations with spiking neurons. Rate coding maps activation value to spike frequency. Accuracy typically drops 1–5% vs the original ANN, but the SNN runs efficiently on neuromorphic hardware.

**Challenges:**
- Training SNNs with backpropagation is hard (spikes are non-differentiable — you can't compute gradients through a step function)
- Surrogate gradient methods approximate the gradient of the spike threshold
- STDP learning is unsupervised and hard to direct toward specific tasks

### Quick Check
> 1. What is Lava and what does it allow programmers to do?
> 2. How can a trained ANN model be deployed on neuromorphic hardware?
> 3. Why is training SNNs with backpropagation difficult?

---

## 6. Neuromorphic vs Conventional AI

```
Task comparison:

                        Neuromorphic SNN    GPU (conventional DNN)
──────────────────────────────────────────────────────────────────
ImageNet accuracy         ~70% (ANN→SNN)    95%+ (ResNet, ViT)
Power (inference)         1–100 mW          50–700W
Latency (first token)     Low (event-driven) Medium (batch GPU)
Continuous streaming      Excellent          Requires batching
On-chip learning          Limited (STDP)     No (need host GPU)
Temporal sequences        Natural (spikes)   Needs RNN/attention
Analog sensor interface   Natural           Requires ADC + preprocessing
Maturity                  Research/demo      Production ready
Software ecosystem        Minimal            Rich (PyTorch, CUDA)
```

**Where neuromorphic wins:**
- Always-on, ultra-low-power sensing (keyword detection, motion, gesture)
- Edge IoT with µW/mW budgets (impossible for GPUs)
- Event-camera data (cameras that output spikes, not frames — natural SNN input)
- Continuous learning without periodic training runs

**Where conventional AI wins:**
- Accuracy (GPT-4, AlphaFold — state of the art is all conventional DNN)
- Software tooling (PyTorch ecosystem is mature; Lava is not)
- Training (SNNs still lag for complex supervised learning tasks)

**Hybrid approach**: Low-power neuromorphic MCU (TrueNorth/Loihi) handles always-on sensing and wakes up a conventional CPU/GPU only when needed. This "neuromorphic front-end" reduces the system power dramatically.

### Quick Check
> 1. For which type of task does neuromorphic computing have a clear advantage over GPUs?
> 2. Why is the software ecosystem gap important for adoption of neuromorphic computing?
> 3. What is an "event camera" and why is it a natural fit for neuromorphic hardware?

---

## Summary

- **Neuromorphic computing** mimics the brain's architecture: spiking neurons, event-driven processing, in-memory computation, and massive parallelism.
- **Brain advantages**: 20W for complex intelligence, event-driven (sparse computation), in-memory (no Von Neumann bottleneck), asynchronous.
- **Spiking Neural Networks (SNNs)**: Use discrete spike events instead of floating-point activations. More power-efficient on specialized hardware, harder to train.
- **TrueNorth** (IBM): 1M neurons, 70mW, fixed LIF model, binary synapses.
- **Loihi 2** (Intel): 1M neurons, <1W, programmable neuron models, on-chip STDP learning.
- **Programming**: Lava framework (Intel), ANN-to-SNN conversion for deployment.
- **Current state**: Neuromorphic chips excel at ultra-low-power sensing and temporal tasks; conventional DNNs lead in accuracy, training, and software ecosystem.

---

## Exercises

### Easy
1. What is a neuron and how does it decide to "fire"?
2. What is the key difference between an ANN activation (floating-point) and an SNN spike?
3. Why does event-driven computing save power compared to clock-driven computing?

### Medium
4. Power comparison: A smart doorbell needs to run face detection 24/7. Compare: (a) ARM Cortex-A53 at 1GHz running MobileNet-V2: 500mW continuous. (b) TrueNorth SNN at 70mW: 70% accuracy. (c) "Neuromorphic wakeup" approach: TrueNorth always-on at 70mW, wakes Cortex-A53 only for confirmed faces (10 times/day, 100ms each): what is the average power? (d) Over one year (8760 hours), what is the energy consumption in Wh for each approach?
5. Spiking vs rate coding: A neuron needs to encode the value 0.75 (out of max 1.0). (a) In rate coding over 100ms window at max 100 Hz: how many spikes would represent 0.75? (b) In temporal coding (earlier spike = higher activation): if max activation fires at 1ms and min at 100ms, at what time does 0.75 fire? (c) Compare the information capacity: rate coding with 100Hz max over 100ms vs binary spike timing to 1ms precision. Which carries more information per neuron?
6. Memristor crossbar: A 4×4 memristor crossbar implements weight matrix W (4 output neurons × 4 inputs). Memristor resistances: R_ij = 1/w_ij kΩ, with w_ij between 0.1 and 1. Input voltages: [1V, 0V, 1V, 0V]. (a) Using Ohm's law (I = V/R = V × w), calculate the current into each output neuron. (b) What classical operation does this implement? (c) What is the power consumed? (d) Compare to a digital 4×4 matrix multiply: operations count and power if each multiplier uses 1pJ.

### Hard
7. SNN training with surrogate gradients: Standard backpropagation fails for SNNs because the spike function is a step function (derivative = 0 everywhere, except ∞ at the threshold). Surrogate gradient methods replace the true derivative of the spike function with a smooth approximation during the backward pass. (a) Write the forward pass for a LIF neuron including the spike decision. (b) What is the derivative of the Heaviside step function, and why can't you backpropagate through it? (c) Design a surrogate function (e.g., a sigmoid with a steep slope, or a piecewise linear function) and justify your choice. (d) How does this compare to training RNNs with truncated BPTT (backpropagation through time)?
8. Neuromorphic for autonomous edge AI: You are designing an insect-sized drone (weight: 5 grams, battery: 1Wh) that must navigate a room autonomously for 30 minutes using computer vision and path planning. (a) Power budget: 1Wh / 0.5h = 2W total system. Motors use 1.5W. What is the compute budget? (b) Compare platforms: ARM Cortex-M4 (30 DMIPS, 10mW), Intel Loihi 2 (1M neurons, ~100mW), MobileNet-V2 on Cortex-A53 (500mW). Which fits? (c) Design a hybrid compute architecture using neuromorphic for obstacle detection + minimal classical processor for navigation. (d) What are the software challenges: how do you train the SNN for obstacle detection when labeled spike-format training data doesn't exist?
