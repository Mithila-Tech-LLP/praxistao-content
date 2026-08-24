# Chapter 50: Quantum Computing — A Different Kind of Computer

Everything in this course so far — transistors, logic gates, CPUs, caches, memory — operates on classical bits: each bit is either 0 or 1. Quantum computers operate on qubits, where a bit can be 0, 1, or any quantum superposition of both simultaneously. This sounds like a gimmick, but for specific mathematical problems — factoring large numbers, simulating quantum systems, searching unsorted databases — quantum computers offer exponential speedup over any classical algorithm. This is not a replacement for classical computers. A quantum computer cannot browse the web or run Word. But for the specific problems it can solve, it can in hours what would take classical computers millions of years. This chapter explains the physics, the architecture, the limitations, and the state of the field in 2024.

## Table of Contents

1. [Classical vs Quantum Bits](#1-classical-vs-quantum-bits)
2. [Quantum Gates and Circuits](#2-quantum-gates-and-circuits)
3. [Quantum Algorithms That Matter](#3-quantum-algorithms-that-matter)
4. [Quantum Hardware](#4-quantum-hardware)
5. [Quantum Error Correction](#5-quantum-error-correction)
6. [Where Are We in 2024?](#6-where-are-we-in-2024)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Classical vs Quantum Bits

**Classical bit**: Binary. Exactly 0 or exactly 1 at all times. A register of 3 bits holds one value at a time — 000, 001, 010, ..., 111 (one of 8).

**Qubit (quantum bit)**: Can exist in a **superposition** of 0 and 1 simultaneously. The quantum state is: |ψ⟩ = α|0⟩ + β|1⟩, where α and β are complex numbers (amplitudes) satisfying |α|² + |β|² = 1.

When you **measure** the qubit, the superposition collapses to either 0 (with probability |α|²) or 1 (with probability |β|²). Before measurement, it holds both possibilities.

```
Classical register (3 bits):
  State: 011  ← exactly one of 8 possibilities
  
Quantum register (3 qubits):
  State: α₀|000⟩ + α₁|001⟩ + α₂|010⟩ + ... + α₇|111⟩
  
  All 8 values simultaneously, weighted by amplitudes.
  A 3-qubit register "holds" 8 values at once.
  A 300-qubit register holds 2^300 values simultaneously.
  (2^300 ≈ 10^90 — more than atoms in the observable universe)
```

**Entanglement**: Two qubits can be quantum-entangled: measuring one instantly determines the state of the other, no matter how far apart. Einstein called this "spooky action at a distance." Entanglement allows quantum algorithms to create correlations between qubits that classical computation cannot efficiently simulate.

**Interference**: Quantum states can constructively or destructively interfere, like waves. Quantum algorithms exploit this to amplify the probability of correct answers and cancel wrong answers. This is the core mechanism that makes quantum algorithms powerful.

**The catch**: You can only extract n classical bits from n qubits (one measurement per qubit). The quantum speedup must come from intelligent design of the quantum circuit — manipulating all 2ⁿ states simultaneously to increase the probability that measurement gives the right answer.

### Quick Check
> 1. What is superposition and how does a qubit differ from a classical bit?
> 2. A quantum computer has 50 qubits. In what sense does it "hold" 2⁵⁰ values simultaneously?
> 3. Why can't a quantum computer simply output all 2ⁿ values from a quantum register?

---

## 2. Quantum Gates and Circuits

Quantum computing uses **quantum gates** — unitary operations on qubits — analogous to classical logic gates.

**Key single-qubit gates:**

```
Hadamard (H) gate — creates superposition:
  H|0⟩ = (|0⟩ + |1⟩)/√2   ← equal superposition
  H|1⟩ = (|0⟩ - |1⟩)/√2
  
  Visually: rotates state by 45° on the Bloch sphere

Pauli-X gate (quantum NOT):
  X|0⟩ = |1⟩
  X|1⟩ = |0⟩

Pauli-Z gate (phase flip):
  Z|0⟩ = |0⟩
  Z|1⟩ = -|1⟩   ← flips the sign (phase) of |1⟩
  
Phase (S) and T gates: smaller rotations
```

**Key two-qubit gate:**
```
CNOT (Controlled-NOT):
  CNOT|00⟩ = |00⟩  (control=0, no flip)
  CNOT|01⟩ = |01⟩
  CNOT|10⟩ = |11⟩  (control=1, flip target)
  CNOT|11⟩ = |10⟩
  
  Analogous to XOR in classical logic, but quantum:
  If control qubit is in superposition, CNOT creates entanglement
```

**Creating a Bell state (maximally entangled pair):**
```
Start: |00⟩
Apply H to qubit 0: (|0⟩+|1⟩)/√2 ⊗ |0⟩ = (|00⟩+|10⟩)/√2
Apply CNOT: (|00⟩+|11⟩)/√2  ← Bell state!

Now if you measure qubit 0 = 0, qubit 1 is definitely 0.
If you measure qubit 0 = 1, qubit 1 is definitely 1.
This is entanglement.
```

**Quantum circuit model**: A quantum algorithm is expressed as a sequence of quantum gates on qubits, ending with measurements.

```
Example quantum circuit:
  q0: ─H──●──────M
          │
  q1: ────X──H──M
```

All quantum gates must be **reversible** (unitary): unlike classical NAND gates, quantum gates don't destroy information. The circuit must run forward AND be theoretically runnable backward. This is a fundamental physical constraint (unitarity of quantum mechanics).

### Quick Check
> 1. What does the Hadamard gate do?
> 2. What is a CNOT gate and what does it do to an entangled state?
> 3. Why must all quantum gates be reversible?

---

## 3. Quantum Algorithms That Matter

Not all problems are faster on quantum computers. The quantum speedup applies to a specific class of problems.

**Shor's Algorithm (1994) — integer factoring:**
- Classical best: sub-exponential (number field sieve) — still incredibly slow for large N
- Quantum: polynomial time O((log N)³) — exponentially faster
- Impact: RSA encryption is based on the hardness of factoring. A 2048-bit RSA key that would take classical computers millions of years to break could be broken by a quantum computer in hours.
- Status: Requires fault-tolerant quantum computer with ~4,000–20 million physical qubits to break 2048-bit RSA. Current machines have ~1,000 noisy qubits. Not a current threat but motivates post-quantum cryptography (NIST PQC standardization, 2024).

**Grover's Algorithm (1996) — unstructured search:**
- Classical: O(N) to search an unsorted list of N items
- Quantum: O(√N) — quadratic speedup (not exponential)
- Impact: Reduces 128-bit AES brute-force attack from 2¹²⁸ to 2⁶⁴ operations. This halves key security — meaning AES-256 is still secure against quantum, but AES-128 may not be.
- Application: Database search, NP-complete problem solving (modest speedup)

**Quantum simulation (Feynman's 1982 motivation):**
- Simulating quantum systems (molecules, materials) on classical computers is exponentially hard
- Quantum computers can simulate other quantum systems efficiently
- Applications: Drug discovery (simulate protein-drug interactions), materials science (high-temperature superconductors), battery chemistry
- This may be the first practical application of quantum computers

**Quantum ML (NISQ era):**
- Variational Quantum Eigensolvers (VQE), Quantum Approximate Optimization Algorithm (QAOA)
- Claims of quantum advantage are contested; limited proof of real-world advantage
- Active research area but no demonstrated quantum advantage over classical ML yet

### Quick Check
> 1. What is Shor's algorithm and why does it threaten RSA encryption?
> 2. What speedup does Grover's algorithm provide? Is it exponential?
> 3. What was Richard Feynman's original motivation for quantum computing?

---

## 4. Quantum Hardware

Building qubits is physically hard. A qubit must maintain quantum coherence (superposition) long enough to compute, while being isolated from environmental noise that would destroy it.

**Superconducting qubits** (IBM, Google, Rigetti):
- Josephson junction circuits cooled to ~15 millikelvin (colder than outer space, which is 2.7K)
- Qubit frequency: 4–8 GHz (microwave control)
- Coherence time: 50–500 µs
- Gate time: ~10–100 ns → can do ~1000–50,000 gates per coherence time
- Connectivity: nearest-neighbor on a 2D grid (limited connectivity requires SWAP gates)
- IBM Eagle (127 qubits, 2021), Osprey (433 qubits, 2022), Condor (1121 qubits, 2023)
- Google Sycamore (53 qubits, 2019 — claimed quantum supremacy)

```
Dilution refrigerator cooling system:
  Room temperature: 300K
  Outer stage: 4K (liquid helium)
  Inner stages: 800mK → 100mK → 20mK → 15mK (base)
  
  The quantum processor hangs at the bottom at 15mK
  Size: refrigerator-sized, weighs ~500kg
```

**Trapped ion qubits** (IonQ, Honeywell/Quantinuum):
- Individual atoms (typically ¹⁷¹Yb⁺ or ⁴⁰Ca⁺) levitated in electromagnetic traps
- Qubit encoded in hyperfine atomic states
- Coherence time: seconds to minutes (much longer than superconducting!)
- Gate time: 10–100 µs (10–100× slower than superconducting)
- Connectivity: All-to-all (any qubit can interact with any other via motional modes)
- Current: IonQ Forte (32 qubits), Quantinuum H2 (32 qubits, best error rates)

**Photonic qubits** (PsiQuantum, Xanadu):
- Qubits encoded in photon polarization or phase
- Room temperature (photons don't need cooling)
- High gate error rates with current photon sources/detectors
- PsiQuantum aims for fault-tolerant silicon photonics at wafer scale

**Neutral atom qubits** (QuEra, Pasqal, Atom Computing):
- Atoms trapped by optical tweezers (laser beams)
- High qubit count (256+ atoms), programmable connectivity
- QuEra Aquila: 256 qubits (2022)

```
Qubit technology comparison (2024):

                Superconducting  Trapped Ion   Neutral Atom
Coherence time  50–500 µs        1–100 seconds  1–100 seconds
Gate time       10–100 ns        10–100 µs      10–1000 µs
Gate error      0.1–1%           0.01–0.5%      1–5%
Qubit count     Up to 1,121      Up to 32       Up to 256
Connectivity    2D grid          All-to-all     Configurable
Temperature     15 mK            Room temp      Room temp
Company         IBM, Google      IonQ, Quantinuum  QuEra, Pasqal
```

### Quick Check
> 1. Why do superconducting quantum computers need to be cooled to 15 millikelvin?
> 2. What is "coherence time" and why does it limit quantum computation?
> 3. What advantage do trapped ion qubits have over superconducting qubits?

---

## 5. Quantum Error Correction

Current quantum computers are **NISQ (Noisy Intermediate-Scale Quantum)** devices: quantum noise causes qubit errors at 0.1–1% per gate. A 1000-gate algorithm with 1% error per gate will be wrong: (0.99)^1000 ≈ 4.3% success rate. Not useful.

**Fault-tolerant quantum computing** requires **quantum error correction (QEC)**: using many physical qubits to implement one logical qubit with a much lower effective error rate.

**The threshold theorem**: If physical qubit error rate is below a threshold (~1%), then using more physical qubits per logical qubit exponentially reduces the logical error rate.

**Surface code** (most promising QEC code):
- Arrange physical qubits in a 2D lattice
- Use surrounding qubits to detect errors without measuring the data qubits directly
- 1 logical qubit requires ~1000 physical qubits at current error rates
- To break RSA-2048: need ~4,000 logical qubits × 1,000 physical/logical = 4 million physical qubits
- Current best: IBM Condor at 1,121 physical qubits — 4000× short

```
Surface code logical qubit:
  
  d×d grid of physical qubits (d = code distance)
  d=7: 7×7 = 49 physical qubits → 1 logical qubit
  Logical error rate ≈ (p/p_th)^(d/2)
    p = physical error rate = 0.1%
    p_th = threshold ≈ 1%
    d = 7 → logical error ≈ 10^-7 per operation
  
  IBM's 2033 target: 100,000+ physical qubits → 100 logical qubits
```

### Quick Check
> 1. What is a NISQ device?
> 2. Why does the surface code require ~1000 physical qubits per logical qubit?
> 3. What is the "threshold theorem" in quantum error correction?

---

## 6. Where Are We in 2024?

**Quantum supremacy claims:**
- Google (2019): Sycamore 53-qubit chip performed a specific random circuit sampling task in 200 seconds that Google claimed would take Summit supercomputer 10,000 years. IBM disputed this, saying ~2.5 days with better classical algorithm.
- IBM (2023): Eagle 127-qubit computer performed "beyond-classical" computation on a specific physics simulation.
- These claims are specific to narrow artificial tasks — not general quantum advantage.

**Practical applications today (2024):**
- Quantum chemistry simulations on small molecules (H₂, LiH) — proofs of concept
- Quantum optimization (QAOA on small problems)
- Quantum key distribution (QKD) — quantum cryptography, not quantum computing
- Financial portfolio optimization (small scale)
- None of these beat classical computers in production use cases yet

**The timeline:**
- 2025–2030: NISQ with limited error correction, ~1,000–10,000 physical qubits
- 2030–2040: Early fault-tolerant systems, ~10,000–1,000,000 physical qubits, some quantum advantage for narrow tasks (chemistry, optimization)
- 2040+: Large-scale fault-tolerant quantum computers; cryptographically relevant (threatens RSA)

**Post-quantum cryptography (PQC)**: NIST finalized PQC standards in 2024 (ML-KEM/Kyber for key encapsulation, ML-DSA/Dilithium for signatures). These are classical algorithms believed to resist quantum attacks. Migration from RSA/ECC to PQC is underway.

### Quick Check
> 1. What was Google's "quantum supremacy" claim in 2019?
> 2. Are there any practical applications of quantum computers in production today (2024)?
> 3. What is post-quantum cryptography and why is it important now, before quantum computers can break RSA?

---

## Summary

- **Quantum computers** use qubits that can exist in superposition (both 0 and 1 simultaneously), exploit entanglement and interference to solve specific mathematical problems exponentially faster than classical computers.
- **Quantum gates** (Hadamard, CNOT, Pauli) manipulate qubit states. All gates must be reversible (unitary).
- **Key algorithms**: Shor's (threatens RSA, exponential speedup for factoring), Grover's (quadratic speedup for search), quantum simulation (chemistry/materials science).
- **Hardware types**: Superconducting (IBM, Google, 15mK cooling required), trapped ion (IonQ, room temp, slower gates but longer coherence), neutral atom (QuEra, programmable connectivity).
- **Error correction**: Current NISQ devices are too noisy for fault-tolerant computation. Surface code requires ~1000 physical qubits per logical qubit. Breaking RSA-2048 needs ~4 million physical qubits.
- **State of the art (2024)**: IBM has 1,121 physical qubits. No demonstrated practical quantum advantage over classical computers. Post-quantum cryptography migration underway as defense.

---

## Exercises

### Easy
1. What is superposition? Explain it to a non-physicist using an analogy.
2. What is quantum entanglement and what is "spooky action at a distance"?
3. Why is Shor's algorithm a threat to internet security (RSA encryption)?

### Medium
4. Qubit state probabilities: A qubit is in state |ψ⟩ = 0.6|0⟩ + 0.8|1⟩. (a) Verify this is a valid quantum state (|α|² + |β|² = 1). (b) What is the probability of measuring 0? (c) What is the probability of measuring 1? (d) After measuring 0, what is the new state of the qubit? (e) If you apply the Hadamard gate to |0⟩, what is the resulting state and what measurement probabilities does it have?
5. Error correction overhead: A fault-tolerant quantum computer needs 1,000 logical qubits to run Shor's algorithm on RSA-2048. Surface code requires 1,000 physical qubits per logical qubit at d=25 distance. (a) How many physical qubits are needed? (b) IBM adds 100 new physical qubits per year (rough estimate). How many years to reach this target? (c) If qubit error rates improve from 0.1% to 0.01%, recalculate how many physical qubits are needed per logical qubit. (d) How does this change the timeline?
6. Grover's algorithm and symmetric encryption: AES-128 key space is 2¹²⁸. Grover's algorithm reduces search to O(√N). (a) How many Grover iterations to search AES-128 key space? (b) What is this in terms of computational steps? (c) Compare to AES-256: is AES-256 "quantum safe" against Grover? (d) What is the practical time for a 2030-era quantum computer doing 1 billion Grover steps/second?

### Hard
7. Quantum simulation for drug discovery: A drug target protein has 100 electrons in its active site. Classical simulation of 100 electrons requires tracking 2¹⁰⁰ amplitude components (exponential in electrons). (a) Why is classical simulation exponential in the number of electrons? (hint: what must you track for a quantum system?) (b) A quantum computer can simulate this efficiently — what "data structure" does the quantum register naturally encode? (c) In 2024, the largest quantum chemistry simulation is on ~20 qubits (simulate ~20-electron systems). What qubit count is needed for 100-electron simulation? (d) What physical error rate and qubit count is required for this to be practically useful? Is this achievable by 2030?
8. Post-quantum cryptography transition: You are a security architect for a bank with 50 million customer records encrypted with RSA-2048. (a) Estimate when a cryptographically relevant quantum computer (CRQC) might exist (lower bound, upper bound). (b) What is the "harvest now, decrypt later" attack and why must you migrate before a CRQC exists? (c) NIST PQC finalist ML-KEM (Kyber-1024) uses ~2.4KB public keys vs RSA-2048 with 256-byte keys. What are the performance implications for TLS handshakes at 1 billion connections/day? (d) Describe a migration strategy: can you use hybrid RSA + PQC during transition without breaking security?
