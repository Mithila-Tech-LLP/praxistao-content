# Chapter 65: SHAKTI's Security Extensions

SHAKTI's most distinctive contribution is its hardware security extensions. While ARM and Intel also implement hardware security features, SHAKTI's approach is notable for three reasons: the extensions are formally verified (mathematical proof of correctness), the code is openly auditable (no black-box trust required), and the focus on post-quantum cryptography places SHAKTI ahead of most commercial processors. This chapter dives deep into each security extension.

## Table of Contents

1. [Why Hardware Security?](#1-why-hardware-security)
2. [Physical Memory Protection (PMP)](#2-physical-memory-protection-pmp)
3. [Tagged Memory Architecture (SHAKTI-T)](#3-tagged-memory-architecture-shakti-t)
4. [Post-Quantum Cryptography Accelerators](#4-post-quantum-cryptography-accelerators)
5. [Secure Boot and Chain of Trust](#5-secure-boot-and-chain-of-trust)
6. [Information Flow Control](#6-information-flow-control)
7. [Side-Channel Attack Mitigations](#7-side-channel-attack-mitigations)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Hardware Security?

Software security can be bypassed if the hardware itself is compromised. Hardware security provides guarantees that software cannot override:

**Threat model for SHAKTI's target applications:**
- Government secure communications: adversary has physical access to device; must not be able to extract keys or private data
- Military systems: supply chain attack may introduce backdoors; need formal verification of no hidden functions
- Satellite/space: cosmic ray single-event upsets; need fault detection
- Financial/regulatory: need audit trail that cannot be modified by software

**The security argument for open hardware**: A closed-source chip (ARM, Intel) requires trusting the vendor. A government with adversarial relationships with Western chip companies cannot trust Western black-box chips for cryptographic operations. SHAKTI's open RTL means any qualified team can audit the design and verify it contains no backdoors.

```
Security layers and hardware guarantees:
  
  Application layer:  Software can lie, be exploited, contain bugs
  OS layer:           OS can be compromised (Spectre, Meltdown)
  Hypervisor layer:   Hypervisor can be compromised
  Hardware layer:     Cannot be bypassed by software if correct
                      (BUT hardware must be trusted to be correct)
  
  SHAKTI's answer: Make hardware auditable (open source) + formally verified
```

### Quick Check
> 1. Why can't software alone provide strong security guarantees?
> 2. What is the "trust" argument for open hardware in security applications?
> 3. What are SHAKTI's primary threat models?

---

## 2. Physical Memory Protection (PMP)

**PMP** is a standard RISC-V security feature that SHAKTI fully implements. It is the hardware enforcement of memory isolation.

**How PMP works:**
- 16 PMP entries, each defining: base address, size, permissions (R/W/X), and privilege level
- Before any memory access, hardware checks all 16 entries
- If no entry permits the access at the current privilege level → access fault exception

```
PMP configuration example:
  
  Entry 0: addr=0x80000000, size=1MB, R/W/X, M-mode only
           (secure monitor code and data, accessible only from M-mode)
  
  Entry 1: addr=0x80100000, size=64MB, R/W/X, all modes
           (shared Linux memory, accessible from S and U mode)
  
  Entry 2: addr=0xC0000000, size=4MB, R only, all modes
           (read-only firmware region)
  
  Entry 3: addr=0xFF000000, size=4KB, R/W, M-mode only
           (secure key storage, inaccessible to OS)
  
  If Linux kernel (S-mode) tries to access 0xFF000000 → fault!
  Even if kernel is compromised, it cannot read the key storage.
```

**PMP use in SHAKTI secure systems:**
- Secure monitor (M-mode) configures PMP at boot
- PMP protects the secure monitor itself from modification by OS
- PMP protects key storage regions
- PMP enforces user-space process isolation (reinforces MMU isolation)

**PMP vs MMU**: MMU (Memory Management Unit) provides virtual-to-physical address translation and isolation — but MMU is software-managed (OS controls page tables). PMP is hardware-managed from M-mode — even a compromised OS cannot change PMP settings.

### Quick Check
> 1. What is PMP and which privilege level configures it?
> 2. Why is PMP stronger than MMU-based isolation for security?
> 3. If a PMP entry protects a key storage region at 0xFF000000: what happens if the Linux kernel tries to read from it?

---

## 3. Tagged Memory Architecture (SHAKTI-T)

**SHAKTI-T** is the tagged memory variant of SHAKTI — the processor class most specifically designed for information security applications.

**Tags: what they are and how they work:**
Every 64-bit word in memory has a 4-bit metadata "tag" stored alongside it (conceptually; physically it's part of the cache/memory system). Tags flow through all memory operations:

```
Tag propagation rules:
  
  LOAD instruction: loads value AND its tag into register
    register.tag = memory[addr].tag
  
  STORE instruction: stores value WITH register's tag
    memory[addr].tag = register.tag
  
  ALU instruction: tag of result derived from tags of operands
    Policy examples:
      Conservative: result.tag = max(a.tag, b.tag)
      Taint:        result.tag = a.tag | b.tag (any tainted input → tainted output)
      Custom:       user-defined propagation policy
  
  Control flow: if an indirect jump uses a tainted address → trap
```

**Tag-based memory safety:**
```
Example: Buffer overflow prevention with tags:
  
  Stack frame:
    buf[0..63]    tag=HEAP_DATA (0x1)
    buf[64..127]  tag=HEAP_DATA (0x1)
    ...
    saved_RIP     tag=RETURN_PTR (0x5)  ← return address
  
  strcpy(buf, evil_input):  // evil_input is 200 bytes
    writes buf[0..199] → all tagged HEAP_DATA by the copy operation
    
    When writes reach the saved_RIP location:
      saved_RIP.tag becomes HEAP_DATA (0x1)
      
    On function return: ret loads address with tag=HEAP_DATA
      Hardware check: tag must be RETURN_PTR (0x5)
      Mismatch → TRAP (ILLEGAL_INSTRUCTION_ACCESS) before execution
      
    Exploit PREVENTED at hardware level.
```

**Tag policies**: Different tag values can enforce different policies:
- Tag 0: uninitialized/default
- Tag 1: untrusted (external network data, user input)
- Tag 2: trusted OS data
- Tag 3: kernel pointers
- Tag 4: return addresses (cannot be derived from non-tag-4 data)
- Tag 5: encryption keys (cannot be stored to untrusted memory)
- Tag 6–15: application-defined

**SHAKTI-T formal verification**: The IIT Madras team used the **Coq proof assistant** to formally verify:
1. Tag propagation correctness (tags are correctly propagated through all instructions)
2. Non-interference: high-tag data cannot flow to low-tag observable outputs

This formal proof is a strong guarantee — not just testing but mathematical certainty for the verified properties.

### Quick Check
> 1. What is a "tag" in SHAKTI's tagged memory and what operations affect it?
> 2. How does tag propagation prevent the JOP/ROP attack class?
> 3. What is "non-interference" and why is its formal proof important?

---

## 4. Post-Quantum Cryptography Accelerators

SHAKTI-T includes hardware accelerators for NIST PQC (Post-Quantum Cryptography) standard algorithms:

**Why hardware acceleration for PQC?**
- CRYSTALS-Kyber (key encapsulation) and CRYSTALS-Dilithium (digital signatures) are computationally heavier than RSA/ECC
- Software PQC on a 500 MHz embedded processor: Kyber keygen ~100ms (too slow for real-time)
- Hardware accelerator: Kyber keygen ~1ms (100× speedup)

**CRYSTALS-Kyber hardware module:**
Kyber's dominant operation is polynomial multiplication in the ring Zq[X]/(X^n+1) where n=256, q=3329:
- NTT (Number Theoretic Transform): like FFT but over integers mod q
- Each NTT: 256-point transform, 7 butterfly stages × 128 operations = 896 modular multiplications
- Hardware: 4 parallel multiplier-accumulate units, dedicated NTT butterfly datapath

```
SHAKTI-T Kyber accelerator:
  
  RISC-V instruction: KYBER.KEYGEN rd, rs1, rs2
  
  Hardware:
    1. Load random seed from rs1 address (32 bytes from TRNG)
    2. Hash seed (SHAKE128/SHA3)
    3. NTT module: generate public matrix A (256×256 polynomials)
    4. Sample secret s and error e
    5. NTT(s), NTT(e)
    6. Compute public key: t = NTT(A)·NTT(s) + e
    7. Store key pair at addresses in rd, rs2
    
  Execution: ~50,000 cycles at 500 MHz = 100 µs (vs 50M cycles in software)
```

**CRYSTALS-Dilithium for signing:**
Digital signatures for authentication. Dilithium is used for:
- Firmware signing (verify bootloader integrity)
- Command authentication (verify commands to ISRO satellite are from authorized source)
- Document signing in government applications

**TRNG (True Random Number Generator):**
PQC requires high-quality randomness. SHAKTI-T includes a hardware TRNG based on ring oscillator entropy sources and AES-CTR-DRBG (Deterministic Random Bit Generator). The TRNG is formally specified to provide 256 bits of entropy.

### Quick Check
> 1. What is NTT and why is hardware acceleration needed for it?
> 2. What is CRYSTALS-Dilithium used for and why is it important for satellite applications?
> 3. What is a TRNG and why is it necessary for cryptographic operations?

---

## 5. Secure Boot and Chain of Trust

**Secure boot** ensures that the chip only executes verified, authenticated software from power-on. This prevents an adversary with physical access from loading malicious firmware.

**SHAKTI secure boot flow:**
```
Power-on:
  1. Reset vector (M-mode) → ROM (read-only, immutable)
  
  ROM contains:
    - Public key hash (burned into OTP at manufacturing)
    - Minimal crypto code for signature verification
  
  Boot sequence:
    ROM loads first-stage bootloader from flash
    Verifies signature using Dilithium public key (from OTP)
    If valid: jump to bootloader
    If invalid: hang (do not boot untrusted code)
    
    First-stage bootloader:
    Verifies second-stage bootloader signature
    (chain of trust: each stage verifies next)
    
    → Linux kernel with verified boot
    → Application with attestation
```

**OTP (One-Time Programmable memory):** Fuse bits that can be written once and never changed. Used to permanently store:
- Root public key hash (for signature verification)
- Debug interface lock bit (disable JTAG after manufacturing test)
- Security policy bits (boot in secure mode only)

**Attestation**: SHAKTI can cryptographically prove to a remote server what software is running. The secure monitor holds a device-specific private key (in PMP-protected memory); it signs a measurement (hash) of the current software stack. The remote server verifies the signature and measurement → knows exactly what is running.

### Quick Check
> 1. What is "chain of trust" in secure boot?
> 2. What is OTP memory and what is it used for in secure boot?
> 3. What is attestation and why is it useful for cloud/IoT deployments?

---

## 6. Information Flow Control

**IFC (Information Flow Control)** is the formal framework for reasoning about how information flows through a system. SHAKTI-T implements hardware IFC using its tag system.

**Lattice-based IFC**: Tags form a security lattice:
```
High-security (tag=SECRET)
       ↑
Mid-security (tag=CONFIDENTIAL)
       ↑
Low-security (tag=PUBLIC)

Information can flow from lower to higher security:
  PUBLIC data can be used to compute on SECRET data
  
Information CANNOT flow from higher to lower security:
  SECRET data cannot be written to PUBLIC memory → confidentiality
  PUBLIC data cannot control SECRET execution flow → integrity
```

**Non-interference**: The formal property that high-security inputs do not affect low-security outputs. If an adversary can observe only "public" outputs, they cannot infer "secret" inputs. SHAKTI-T's formal verification proves this property holds for the ISA semantics.

**Implicit flows**: The dangerous corner case — even a 1-bit branch on secret data leaks information:
```c
if (secret_key & 1) {  // conditional on secret bit
    public_output = 1;  // public observation changed!
}
// This is a covert channel even though we never stored secret to public
```
Hardware IFC must track implicit flows (control dependencies) as well as explicit flows (data dependencies). SHAKTI-T addresses this by tagging the program counter's conditional expression.

### Quick Check
> 1. What is a security lattice and what does it enforce?
> 2. What is "non-interference" as a formal property?
> 3. What is an "implicit flow" and why is it dangerous?

---

## 7. Side-Channel Attack Mitigations

Hardware security must also address side-channel attacks — attacks that infer secret data from physical observables (timing, power, electromagnetic radiation).

**Spectre/Meltdown mitigations in SHAKTI (in-order):**
Because SHAKTI C-class and E-class are in-order processors, they are NOT vulnerable to Spectre v1 (speculative store bypass) — the vulnerability arises from out-of-order speculative execution. SHAKTI's in-order design is a security advantage for embedded/IoT uses.

**Cache timing attacks**: An attacker measures how long memory accesses take; if a secret determines which cache line is loaded, access time reveals the cache line index. SHAKTI-T mitigates this:
- Cache partitioning: each security domain (M-mode, S-mode, U-mode) gets separate L1 cache ways
- Cache flush on privilege level change
- Constant-time mode: all cache misses take the same time (sacrifices performance for security)

**Power analysis (DPA — Differential Power Analysis)**: Correlate power consumption with secret key bits during cryptographic operations. SHAKTI-T's PQC accelerator uses:
- Masked implementations: randomize intermediate values so power is uncorrelated with secret
- Dual-rail logic in critical circuits (constant power regardless of data)

**Electromagnetic side-channel**: Similar to power analysis but using EM emissions. Physical countermeasures (not digital) needed for highest assurance.

### Quick Check
> 1. Why is SHAKTI's in-order design an advantage against some Spectre variants?
> 2. What is DPA (Differential Power Analysis) and how does SHAKTI-T counter it?
> 3. What is "cache partitioning" and how does it prevent cache timing attacks?

---

## Summary

- **SHAKTI's security extensions** go beyond most commercial processors: formally verified, open-source, designed for high-assurance strategic applications.
- **PMP (Physical Memory Protection)**: RISC-V standard, hardware-enforced memory regions, configured by M-mode, immune to OS compromise.
- **Tagged memory (SHAKTI-T)**: 4-bit tags per word, propagated through computation, prevents buffer overflows/type confusion/code injection at hardware level. Formally verified in Coq.
- **PQC accelerators**: Hardware NTT for CRYSTALS-Kyber and Dilithium — 100× speedup, constant-time execution.
- **Secure boot**: Chain of trust from OTP-burned root key through signed bootloader stages to OS.
- **IFC**: Tag-based lattice, non-interference property formally proven.
- **Side-channel mitigations**: In-order design (Spectre resistant), cache partitioning, masked PQC.

---

## Exercises

### Easy
1. What is PMP and what problem does it solve that the MMU cannot?
2. What is a tag in SHAKTI-T and how does it prevent buffer overflow exploits?
3. What is secure boot and why is it important for satellite applications?

### Medium
4. Tag-based exploit prevention: A web server is running on SHAKTI-T. Network data (from user input) has tag=UNTRUSTED. A buffer overflow vulnerability exists in the HTTP parser. (a) The attacker sends a crafted request that overflows the parser's buffer and overwrites a function pointer on the stack. With tags: what tag does the overflowed data have? What happens when the function pointer is called? (b) ROP (Return-Oriented Programming) attacks reuse existing code gadgets by overwriting return addresses. Return addresses have tag=RETURN_PTR. If the overflow overwrites a return address with tag=UNTRUSTED data: what does SHAKTI-T do? (c) What is an attack that tag-based protection CANNOT prevent? (hint: attacker uses a logic bug, not a memory safety bug)
5. PQC performance: SHAKTI-T targets 500 MHz. Hardware Kyber keygen runs in 50,000 cycles. (a) Time for one Kyber keygen? (b) A TLS 1.3 handshake requires: 1 keygen + 2 encapsulations + 2 decapsulations. Assuming each operation takes similar cycles: total handshake time? (c) For an IoT device making 10 connections/minute: can SHAKTI-T handle this? (d) On an ARM Cortex-M4 at 168 MHz without hardware acceleration, software Kyber keygen takes 200 ms. Compare with SHAKTI-T. (e) Why is 100× faster PQC important for TLS connections in high-concurrency servers?
6. Formal verification scope: The SHAKTI team formally verified tag propagation for ALU operations. (a) What is the Coq proof assistant and what guarantees does a Coq-verified proof provide? (b) The proof covers: all arithmetic instructions, logical instructions, load/store. What is NOT covered by the current proof? (hint: think about multi-core and caches) (c) A critic argues: "The formal proof is for the ISA specification, not the RTL implementation. The RTL could have a bug not reflected in the spec." Is this a valid concern, and how would you address it?

### Hard
7. Information flow security analysis: SHAKTI-T is used to protect a cryptocurrency key (tag=SECRET). The wallet application performs: (a) ECDSA signing (using secret key k): result is a public signature. The signature reveals the public key, not k itself. How does IFC handle this? (b) An attacker can make the wallet sign arbitrary messages. After observing enough signatures, can they infer k? (this is a different vulnerability — not an IFC issue). (c) The wallet uses a random nonce during signing. If the TRNG is compromised (weak random), k can be recovered from 2 signatures. How should SHAKTI-T's TRNG interact with the IFC framework? (d) Design a threat model: what attacks does SHAKTI-T's hardware security prevent vs what requires protocol-level defense?
8. SHAKTI-T for satellite application: ISRO uses SHAKTI-T in a satellite commanding system. The system receives uplink commands from a ground station. Threat: adversary intercepts the uplink channel and injects forged commands. Security requirements: (a) Command authentication using Dilithium (verify each command is signed by authorized ground station). Sketch the verification flow using SHAKTI-T's Dilithium accelerator. (b) Replay prevention: how does the system detect a recorded authentic command being replayed hours later? (hint: sequence numbers, timestamps). (c) If the satellite is physically captured by an adversary: what does SHAKTI-T's secure boot guarantee about the software running on it? What does it NOT guarantee? (d) Post-mission analysis requires downloading telemetry. Telemetry contains classified sensor data (SECRET tag). Design a protocol using SHAKTI-T's PMP and PQC to securely transfer telemetry without leaking data to the downlink physical layer.
