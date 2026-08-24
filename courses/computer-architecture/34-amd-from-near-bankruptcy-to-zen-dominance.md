# Chapter 34: AMD — From Near-Bankruptcy to Zen Dominance

In 2015, AMD's stock traded below $2. The company had just reported its fourth consecutive year of losses. Its latest CPU (Bulldozer architecture) was widely mocked as a performance disaster. Its GPU division (Radeon) was losing ground to NVIDIA. Many analysts predicted AMD would be acquired or go bankrupt. Four years later, AMD's stock was above $50. By 2022, it hit $160. The turnaround — powered by a completely new architecture called Zen, a manufacturing partnership with TSMC, and a chiplet-based design philosophy — is one of the most dramatic comebacks in semiconductor history.

## Table of Contents

1. [AMD's Origins — Born as Intel's Clone](#1-amds-origins--born-as-intels-clone)
2. [K7/K8 — The Athlon Era of Dominance](#2-k7k8--the-athlon-era-of-dominance)
3. [Bulldozer — The Near-Death Architecture (2011–2017)](#3-bulldozer--the-near-death-architecture-20112017)
4. [Zen — The Comeback (2017)](#4-zen--the-comeback-2017)
5. [Chiplets and Infinity Fabric — The Design Revolution](#5-chiplets-and-infinity-fabric--the-design-revolution)
6. [Zen 2, 3, 4 — Relentless Improvement](#6-zen-2-3-4--relentless-improvement)
7. [AMD EPYC — Server Dominance](#7-amd-epyc--server-dominance)
8. [3D V-Cache — Stacking Memory on CPU](#8-3d-v-cache--stacking-memory-on-cpu)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. AMD's Origins — Born as Intel's Clone

AMD was founded in 1969, one year after Intel. Its early business was literally **second-sourcing Intel chips** — IBM required a second source supplier for the 8086 in its PCs, and AMD was licensed to manufacture Intel's x86 chips. AMD built its business on manufacturing Intel designs reliably.

In 1991, AMD and Intel's cross-licensing agreement expired, and AMD began designing its own x86-compatible chips. The **Am386** (1991) was AMD's first independently designed x86 processor — and it was significantly cheaper and often faster than Intel's 80386. AMD had found its strategy: compete on price and occasionally on performance.

**The x86 legal framework**: Intel and AMD signed cross-licensing agreements in 1976, 1982, and 1995. These agreements allowed both companies to make x86-compatible processors. Intel periodically tried to revoke AMD's licenses (over the 287 math coprocessor, over the K5 design) but courts generally upheld AMD's right to manufacture compatible chips.

### Quick Check
> 1. What is "second sourcing" and why did IBM require it for PC components?
> 2. How did AMD differentiate from Intel in the early 1990s?
> 3. What legal framework allows AMD to make x86-compatible chips?

---

## 2. K7/K8 — The Athlon Era of Dominance

AMD's golden era began with the **K7 (Athlon, 1999)** — the first time AMD unambiguously beat Intel's best chip at every clock speed.

The Athlon K7 was designed under Dirk Meyer (recruited from DEC's Alpha team) and featured:
- 9-issue superscalar out-of-order execution
- Full-speed on-die L2 cache (competitors had slower off-chip L2)
- 200 MHz system bus (double the Pentium III's FSB)
- 3DNow! SIMD instructions for 3D games

The Athlon beat the Pentium III and early Pentium 4 in gaming and floating-point workloads. This was the first time AMD had definitively better hardware than Intel — and the market responded. AMD's market share surged.

**K8 (Athlon 64, 2003)**: AMD's second great success. Key innovations:
- **First to extend x86 to 64 bits** (AMD64 / x86-64) — Intel was forced to adopt AMD's design
- **On-die memory controller**: AMD integrated the memory controller directly into the CPU before Intel did. This eliminated the separate Northbridge chip and reduced memory latency by ~50%.
- **HyperTransport**: AMD's high-speed CPU-to-chipset and CPU-to-CPU interconnect (before UPI/QPI from Intel)
- **64-bit computing at consumer prices** before Intel offered it

Athlon 64 FX was the fastest consumer processor available from 2003 to 2006 — AMD's high point.

### Quick Check
> 1. What made the Athlon K7 better than the Pentium III?
> 2. What did the K8 Athlon 64 introduce that Intel was later forced to adopt?
> 3. Why did integrating the memory controller into the CPU reduce latency?

---

## 3. Bulldozer — The Near-Death Architecture (2011–2017)

After the K8's success, AMD developed the **K10** (Phenom, Barcelona) for multi-core and servers. K10 was competitive but not dominant. Then came **Bulldozer** — a disaster.

Bulldozer (2011) introduced a radical new design concept: the **Module**. Instead of two complete independent cores, a module contained:
- 2 integer execution pipelines (shared front-end, shared floating-point unit)
- The idea: floating-point operations were rare enough that sharing one FPU between two integer cores was efficient

```
Bulldozer Module (2 "cores"):
  ┌──────────────────────────────────────┐
  │  Shared: Fetch, Decode, FPU, L2 cache│
  │  Core 0: Integer ALU cluster 0        │
  │  Core 1: Integer ALU cluster 1        │
  └──────────────────────────────────────┘
  
  FX-8150: 4 modules × 2 cores = "8 cores"
  (But only 4 FPUs — OS sees 8 cores, gets 4 FP units)
```

The Bulldozer's shared FPU was a fatal flaw — many modern workloads mix integer and floating-point. When both "cores" needed the FPU simultaneously, they contended. Actual performance per "core" was far less than Intel's.

Results: The FX-8150 at 3.6 GHz lost to the Intel Core i5-2500K (Sandy Bridge) at 3.3 GHz despite having "8 cores." Power consumption was 125W vs 95W. A humiliating defeat.

**Piledriver (2012), Steamroller (2014), Excavator (2015)**: Incremental refinements of Bulldozer that improved power efficiency but never fixed the fundamental architectural flaw. AMD's CPU division lost hundreds of millions of dollars.

**The AMD sell-off**: AMD sold its fabs to create GlobalFoundries (2009), becoming a fabless company. This was financially necessary but meant AMD was dependent on external foundries — initially GlobalFoundries, then TSMC.

### Quick Check
> 1. What was the "Module" concept in Bulldozer, and why did it backfire?
> 2. The FX-8150 was marketed as "8-core." Was this accurate? Why or why not?
> 3. Why did AMD sell its fabs, and what was the risk of doing so?

---

## 4. Zen — The Comeback (2017)

In 2012, AMD hired **Lisa Su** as CTO (CEO from 2014) and started a complete architectural overhaul. Dr. Jim Keller (who had designed the K7 and K8, then worked on Apple's A4/A5) was recruited to lead the new CPU architecture.

**Zen (2017)** — codename Summit Ridge, product name Ryzen 7 — was a clean-sheet design. No Bulldozer modules. No shared FPUs. Full SMT (Simultaneous Multithreading) — 2 hardware threads per core. A modern OOO engine:

```
Zen core features:
  - 6-wide front-end decode
  - 168-entry ROB
  - 4 integer, 2 FP, 2 load/store execution units
  - 18% better IPC than Excavator (Bulldozer generation)
  - Full SMT2 (2 logical threads per core)
  - 8MB L3 cache per CCX (4-core cluster)
```

**First-generation Ryzen results (2017):**
- Ryzen 7 1800X (8 cores, 16 threads): matched Intel i7-7700K in multithreaded work at $499 vs $339
- But: single-thread performance was ~10-15% behind Intel's Kaby Lake
- Memory latency issues due to Infinity Fabric (first generation)

The reception was remarkable: AMD had returned to competitiveness. Ryzen offered more cores at lower prices. For multi-threaded workloads (compilation, content creation, servers), Ryzen was a genuine Intel alternative.

### Quick Check
> 1. What was architecturally different about Zen compared to Bulldozer?
> 2. What was the IPC improvement of Zen over the previous Excavator architecture?
> 3. What weakness did first-generation Ryzen have compared to Intel?

---

## 5. Chiplets and Infinity Fabric — The Design Revolution

Zen introduced AMD's most enduring innovation: the **chiplet architecture**.

**Traditional monolithic design**: One large die contains all CPUs, I/O, memory controller. The entire die is manufactured on one process node. Problem: large dies have lower yields (more defects per die), higher costs, and tie all functions to one process node.

**AMD's chiplet approach (Zen 2+):**
```
Ryzen 9 3900X (Zen 2, 2019):
  ┌───────────────────────────────────────────────────────┐
  │              I/O Die (GlobalFoundries 12nm)            │
  │  Memory controllers, PCIe, USB, SATA, Infinity Fabric  │
  └──────┬─────────────────────────────┬──────────────────┘
         │                             │
  ┌──────┴──────┐               ┌──────┴──────┐
  │  CCD #1     │               │  CCD #2     │
  │  Zen 2, 7nm │               │  Zen 2, 7nm │
  │  8 cores    │               │  8 cores    │
  └─────────────┘               └─────────────┘
  
  Total: 16 cores, 32 threads
  
  CCD = Core Complex Die (compute)
  I/O Die = I/O Controller Die
```

**Infinity Fabric (IF)**: AMD's high-speed die-to-die and intra-die interconnect. Connects CCDs to the central I/O die. Scales from intra-chip (connecting CPU clusters to L3 cache) to multi-socket (EPYC server CPU-to-CPU links, called xGMI).

**Advantages of chiplets:**
1. **Higher yield**: Small 7nm compute dies are cheaper to manufacture than one large 7nm die
2. **Mixed process nodes**: Compute (TSMC 7nm) + I/O (GF 12nm) — use the best process for each function
3. **Scalability**: Add more CCDs for higher core counts without redesigning the I/O die
4. **Flexibility**: Same CCD works in Ryzen (consumer) and EPYC (server) — economies of scale

**Disadvantage**: CCD-to-CCD communication crosses the I/O die, adding ~50–80ns latency vs intra-CCD communication. Cache coherence across chiplets is more complex.

### Quick Check
> 1. What is a "chiplet" and what is the difference between a CCD and the I/O die in AMD Zen 2?
> 2. Why does using small chiplets improve manufacturing yield?
> 3. What is the latency cost of cross-CCD communication, and why does it matter for software?

---

## 6. Zen 2, 3, 4 — Relentless Improvement

Each Zen generation brought significant IPC improvements and manufacturing process advances:

**Zen 2 (TSMC 7nm, 2019):**
- Doubled L3 cache per CCD (from 8MB to 16MB)
- Improved branch predictor
- FP execution improvements
- First AMD CPU at TSMC 7nm
- Ryzen 9 3950X: 16 cores, ~30% IPC improvement vs Zen 1

**Zen 3 (TSMC 7nm enhanced, 2020):**
- Reorganized CCD: instead of two 4-core CCX groups each with their own L3, one unified 8-core cluster shares 32MB L3
- All 8 cores in a CCD can directly access all 32MB L3 — eliminated inter-CCX latency
- 19% IPC improvement vs Zen 2
- Ryzen 5000 series: Ryzen 9 5950X (16 cores) beat Intel at single-thread performance for the first time in the Ryzen era

**Zen 4 (TSMC 5nm, 2022):**
- DDR5 and PCIe 5.0 support
- Added AVX-512 (AMD's consumer chips finally got 512-bit SIMD)
- Integrated RDNA 2 iGPU in every desktop chip
- ~13% IPC improvement vs Zen 3
- Ryzen 9 7950X: 16 cores, 32 threads at 5.7 GHz boost

```
Zen IPC progression (relative to Zen 1 = 100%):
  Zen 1 (2017):   100%
  Zen 2 (2019):   ~130%
  Zen 3 (2020):   ~155%
  Zen 4 (2022):   ~175%
  
  From Bulldozer (2011) to Zen 4 (2022):
  Estimated 3-4× IPC improvement over 11 years
```

**Zen 5 (TSMC 4nm/3nm, 2024):**
- Major front-end redesign (wider decode: 4→8 instructions/cycle)
- Doubled AVX-512 throughput
- Significant L1 data cache capacity increase (48KB→48KB but much higher bandwidth)
- ~16% IPC improvement vs Zen 4

### Quick Check
> 1. What was the key architectural change in Zen 3 regarding L3 cache topology?
> 2. At what point did AMD's Ryzen beat Intel at single-thread performance?
> 3. Zen 4 added AVX-512. Why did AMD add a feature that was previously an Intel exclusive?

---

## 7. AMD EPYC — Server Dominance

AMD's server business was nearly dead by 2017. Intel Xeon dominated data centers with ~99% market share. The EPYC server CPU line — based on the same Zen chiplets but scaled further — changed that.

**EPYC Rome (Zen 2, 2019):**
```
EPYC 7742: 64 cores, 128 threads
  8 CCDs × 8 cores = 64 cores
  8 × 16MB L3 = 128MB L3 cache
  8 memory channels (DDR4-3200)
  384GB/s memory bandwidth
  128 PCIe 4.0 lanes
  
  vs Intel Xeon Platinum 8280 (Cascade Lake):
  28 cores, 56 threads — same generation, less than half the cores
```

EPYC Rome offered more cores, more memory bandwidth, and more PCIe lanes at lower prices. Major cloud providers (Amazon AWS, Google Cloud, Microsoft Azure, Oracle Cloud) adopted EPYC for their server fleets. By 2023, AMD's EPYC held ~20-25% server CPU market share — from near zero in 2017.

**EPYC Genoa (Zen 4, 2022):** Up to 96 cores. 12 DDR5 memory channels. PCIe 5.0.

**EPYC Turin (Zen 5, 2024):** Up to 192 cores (128 Zen 5 + 64 Zen 5c dense cores).

### Quick Check
> 1. What was AMD's server CPU market share before EPYC Rome, and approximately what did it reach by 2023?
> 2. EPYC Rome had 64 cores vs Intel's 28-core Xeon — yet both were made on "7nm" (TSMC) vs "14nm" (Intel). How did chiplets make this core count advantage possible?
> 3. Why do cloud providers care so much about cores per server (as opposed to single-thread performance)?

---

## 8. 3D V-Cache — Stacking Memory on CPU

In 2021, AMD announced **3D V-Cache** — a semiconductor packaging breakthrough that stacks additional SRAM directly on top of the CPU die using TSMC's SoIC (System on Integrated Chips) process.

```
Standard Ryzen 9 5950X (Zen 3):   32MB L3 per CCD × 2 CCDs = 64MB L3
Ryzen 9 5900X3D with V-Cache:     32MB native + 64MB stacked = 96MB L3 per CCD

Stacking uses TSVs (Through-Silicon Vias) — tiny metal pillars that pass
vertically through the silicon die, connecting layers.
```

**Why it matters**: The L3 cache is the last resort before DRAM. For workloads that don't fit in standard L3 but do fit in the 3D V-Cache-enhanced L3, cache hit rate jumps dramatically.

**Gaming performance**: Games have complex data access patterns — texture data, AI state, physics data. Ryzen 7 5800X3D (with V-Cache) beat Intel Core i9-12900K in gaming — a remarkable result for a CPU based on older microarchitecture.

**Limitations**: The stacked SRAM insulates the CCD — it can't cool as effectively through the top. Maximum clock speed is slightly lower on V-Cache SKUs vs standard.

**Zen 4 V-Cache (Ryzen 7 7800X3D, 2023)**: 96MB stacked V-Cache. 96 + native = 128MB per CCD. Dominant in gaming benchmarks, often beating the most expensive Intel CPUs.

### Quick Check
> 1. What is 3D V-Cache and how is it physically connected to the CPU die?
> 2. Why does the Ryzen 7800X3D have a lower maximum clock speed than the standard 7700X?
> 3. For which workloads does 3D V-Cache help the most, and for which does it provide little benefit?

---

## Summary

- AMD was founded as an Intel second-source manufacturer, became an independent x86 designer, and had a first dominance era with the **Athlon K7/K8** (1999–2006).
- **Bulldozer** (2011–2017) was a disastrous architecture with shared FPUs ("modules") that failed to compete with Intel Core.
- **Zen** (2017) was a complete clean-sheet redesign — competitive IPC, full SMT, proper OOO execution. A direct comeback.
- **Chiplets + Infinity Fabric**: AMD's packaging innovation — small compute dies (CCD) bonded to a central I/O die. Higher yield, mixed process nodes, scalable core counts.
- **Zen 2/3/4/5** brought 30%+ IPC improvements each generation, with Zen 3 finally beating Intel at single-thread performance.
- **EPYC** server chips used the chiplet advantage to offer 64–192 cores vs Intel's 28–60, taking market share in cloud and HPC.
- **3D V-Cache** stacks SRAM on the CPU die, tripling L3 cache capacity for gaming/HPC at the cost of clock speed.

---

## Exercises

### Easy
1. What was wrong with the Bulldozer "Module" architecture? How did Zen fix it?
2. Explain chiplets in one paragraph. Why are small dies cheaper to manufacture than large ones?
3. What is Infinity Fabric and what roles does it play within an AMD CPU?

### Medium
4. AMD EPYC Rome (2019) had 64 cores using 8 chiplets × 8 cores, manufactured on TSMC 7nm. An Intel Xeon at the same time had 28 cores on a monolithic Intel 14nm die. Estimate the die areas: if a single 8-core Zen 2 CCD is ~74mm² on 7nm, what is the total CCD area for EPYC? Compare this to a theoretical 64-core monolithic 7nm die (hint: scale linearly from the CCD area). What yield advantage do chiplets provide if defect density is 0.1 defects/cm²?
5. Zen 3 unified the 8-core cluster to share one 32MB L3. Zen 2 had two 4-core clusters each with 16MB L3, requiring cache-to-cache transfers between clusters. For a workload where 8 threads share 28MB of data: (a) Does the data fit in Zen 2's L3 per half-cluster? (b) Does it fit in Zen 3's unified 32MB L3? (c) What is the performance implication?
6. Compare Ryzen 7 7800X3D vs 7700X (same silicon, V-Cache vs no V-Cache). The 7700X boosts to 5.4 GHz; 7800X3D boosts to 5.0 GHz. For a game that fits its working set in the larger L3 and has 50% of its instructions memory-bound: estimate whether the cache advantage outweighs the clock speed disadvantage.

### Hard
7. The Infinity Fabric frequency scales with DDR memory speed in first-gen Zen (1:1 ratio). This caused slower DDR3200 (IF at 1600 MHz) to perform better than DDR3600 (IF would need to drop to 1:2 mode at 1600 MHz if the ratio wasn't met). Why does IF frequency affect performance beyond just memory bandwidth? (Hint: cross-CCD coherence traffic.) Design an experiment to quantify the IF frequency effect independently from memory bandwidth.
8. AMD's chiplet advantage depends on small die sizes. As transistors shrink, what happens to chiplet economics? Specifically: (a) if defect density improves 10× for each process node, how does the yield advantage of chiplets change? (b) At what point does die-to-die interconnect overhead (Infinity Fabric latency, power) outweigh the yield advantage? (c) How does this analysis change if chiplets enable mixed-node designs (compute on 3nm + I/O on 6nm)?
