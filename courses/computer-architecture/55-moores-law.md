# Chapter 55: Moore's Law — The Driving Force

In 1965, Gordon Moore (co-founder of Intel) published a paper observing that the number of transistors on an integrated circuit was doubling roughly every year. He later revised this to every two years. This observation, retroactively called **Moore's Law**, became a self-fulfilling prophecy that drove the semiconductor industry for 50+ years — every chip maker, equipment company, software developer, and system designer planned their roadmaps around it. Moore's Law is not a law of physics; it is an empirical observation that sustained investment made true. This chapter covers what Moore's Law actually says, what it enabled, why it is slowing, and what comes next.

## Table of Contents

1. [Moore's Original Observation](#1-moores-original-observation)
2. [What Moore's Law Actually Enabled](#2-what-moores-law-actually-enabled)
3. [Dennard Scaling — The Real Engine](#3-dennard-scaling--the-real-engine)
4. [Why Moore's Law Is Slowing](#4-why-moores-law-is-slowing)
5. [What Continues vs What Stopped](#5-what-continues-vs-what-stopped)
6. [Beyond Moore: Alternative Paths](#6-beyond-moore-alternative-paths)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. Moore's Original Observation

In April 1965, Gordon Moore wrote "Cramming More Components onto Integrated Circuits" in Electronics Magazine. At the time, the most complex chip had 64 transistors. Moore noticed:

- 1959 (Fairchild first planar IC): 1 transistor
- 1960: 4 transistors
- 1961: 8 transistors
- 1962: 16 transistors
- 1963: 32 transistors
- 1964: 64 transistors

The count was roughly doubling every year. Moore projected this trend would continue, reaching 65,000 components on a chip by 1975.

```
Transistor count history (selected):
  
  Year    Chip                      Transistors
  1971    Intel 4004                4,300
  1978    Intel 8086                29,000
  1982    Intel 286                 134,000
  1989    Intel 486                 1,200,000
  1993    Intel Pentium             3,100,000
  2000    Intel Pentium 4           42,000,000
  2006    Intel Core 2 Duo          291,000,000
  2012    Intel Ivy Bridge (quad)   1,400,000,000
  2017    AMD EPYC Rome (7nm)       32,000,000,000
  2021    Apple M1 Max              57,000,000,000
  2024    Apple M4 Pro              28nm process... wait
          Nvidia Blackwell B200     208,000,000,000
  
  Trend: transistor count doubled every ~2 years for 50+ years
```

Moore revised his estimate to doubling every 2 years in 1975 — this is the version commonly cited. He never said speed doubles, or that price halves; the original observation was specifically about transistor count per unit area.

**What Moore's Law became**: Industry bodies (ITRS, later IRDS) turned Moore's Law into coordinated roadmaps. Every 2 years, a new process "node" would deliver ~2× transistor density at the same (or lower) cost per transistor. Chip designers, EDA tool companies, equipment makers, and materials suppliers all planned around this cadence.

### Quick Check
> 1. What exactly did Moore observe in 1965?
> 2. How many transistors did Moore predict chips would have by 1975?
> 3. Why is Moore's Law called a "self-fulfilling prophecy"?

---

## 2. What Moore's Law Actually Enabled

The transistor count doubling drove several downstream benefits:

**Cost per transistor**: More transistors per wafer = lower cost per transistor. Manufacturing costs per wafer grew slowly while die capacity grew exponentially → transistor cost fell ~30% per year for decades.

```
Cost per transistor trend:
  1968: $1,000 per transistor
  1978: $0.10 per transistor
  1988: $0.001 per transistor (1 mil)
  2000: $0.00001 per transistor
  2020: ~$0.0000001 ($10⁻⁷) per transistor
  
  50-year improvement: ~10 billion times cheaper per transistor
```

**CPU clock speed** (until ~2004): Smaller transistors switch faster. Combined with Dennard scaling (next section), each generation delivered faster clocks. The Pentium 4 at 3.8 GHz (2004) seemed like 10 GHz was imminent.

**Memory capacity**: DRAM capacity grew 4× every 3 years (faster than Moore's Law) for decades. Hard disk density followed similar trends.

**Cost of computing**: A computation that cost $1 in 1970 costs less than $10⁻¹² today — a 10¹² (one trillion) improvement. This is why software as a business became possible, smartphones became affordable, and AI training at scale is conceivable.

**Software complexity**: Software complexity grew to match available hardware. Operating systems, databases, browsers, compilers — all expanded to fill the available transistors and clock cycles. Software engineers could code inefficiently and hardware progress would bail them out next generation ("the free lunch").

### Quick Check
> 1. How many orders of magnitude cheaper has a transistor become from 1968 to 2020?
> 2. What is "the free lunch" that software engineers historically relied on from Moore's Law?
> 3. Besides transistor count, what other metrics benefited from Moore's Law scaling?

---

## 3. Dennard Scaling — The Real Engine

Moore's Law (transistor doubling) is only half the story. The real magic came from **Dennard Scaling** (1974, Robert Dennard, IBM). Dennard showed that if you scale transistor dimensions by factor k (<1):

- Transistor area: scales by k²  (2× denser per unit area)
- Gate delay: scales by k  (transistors switch k× faster)
- Supply voltage V: scales by k  (voltage reduced proportionally)
- Dynamic power per gate: P = CV²f, C scales k, V scales k, f scales 1/k → P scales k² 

This meant: **2× transistors per chip, each running k× faster, consuming the same total power as the previous generation**. Every process shrink was a gift: more transistors, faster, free power.

```
Dennard scaling example (k=0.7, ~30% shrink):
  
  Current gen:    Next gen (Dennard):   Change:
  L = 100nm       L = 70nm              0.7×
  V = 1.2V        V = 0.85V             0.7×
  f = 1 GHz       f = 1.4 GHz           1/k = 1.43×
  P_gate = C V²f  P_gate = 0.7C×0.5V²×1.4f = 0.5P   half power!
  
  Result: 2× transistors, 43% faster, same total chip power!
```

**Dennard scaling breakdown (~2004)**: Supply voltage could not continue scaling below ~1V because:
1. The transistor **threshold voltage** (V_th ~0.2–0.3V) cannot scale proportionally (thermal noise kT/q = 26mV is a hard floor)
2. **Subthreshold leakage**: reducing V_th causes exponentially more leakage current — chips leak power even when doing nothing
3. **Gate oxide leakage**: oxide cannot thin below ~1nm (quantum tunneling)

With voltage no longer scaling, power per unit area grew with each node: **the power wall**.

After 2004: clock speeds stopped scaling (10GHz Pentium 4 was throttled to 3.8GHz due to heat), and the industry pivoted to multi-core processors — more cores at the same clock = more throughput without more heat.

### Quick Check
> 1. What does Dennard scaling predict happens to total chip power across generations?
> 2. Why did Dennard scaling break down around 2004?
> 3. What did the chip industry do after Dennard scaling failed?

---

## 4. Why Moore's Law Is Slowing

Multiple physical and economic factors are slowing or stopping Moore's Law:

**Physical limits:**
- Transistor gate length is approaching ~2nm — less than 10 silicon atoms
- Quantum tunneling through gate oxide causes leakage
- Dopant atom variability causes transistor mismatch
- Lithography approaching fundamental diffraction limits even with EUV

**Economic limits:**
- Cost per transistor is no longer falling as steeply
- New fabs cost $20–30 billion to build
- EUV machine: $170M each; High-NA EUV: $350M each
- A 300mm wafer at 3nm (TSMC): ~$20,000 (vs ~$5,000 at 28nm)
- Mask set: $5–10 million per design

```
Cost trend reversal:
  
  Process node:  28nm    16nm    7nm     3nm
  Wafer cost:    $3K     $5K     $12K    $20K
  Die cost at    $0.20   $0.35   $0.80   $1.50
  80mm² die:
  
  For the first time: smaller node costs MORE per transistor at some die sizes.
  Moore's Law = more transistors for less money/power is breaking.
```

**Cadence slowdown:**
- Intel's 10nm: announced 2016, volume production 2019 (3 years late)
- TSMC: 10nm (2016) → 7nm (2018) → 5nm (2020) → 3nm (2022) → N2 (2025) — cadence slower
- DRAM scaling has essentially stopped at 10nm-class (manufacturers use "1α", "1β" naming to obscure)

**Node naming inflation**: "5nm" doesn't mean 5nm gate length (which would be ~20 silicon atoms). It is a marketing name. TSMC N5 has gate length ~12nm. Intel 5 has gate pitch ~27nm. Node names are competitive branding, not physical measurements.

### Quick Check
> 1. What physical limits are causing Moore's Law to slow?
> 2. Why has the cost per transistor stopped falling as reliably as it once did?
> 3. What is "node naming inflation" and why does it make comparisons between vendors difficult?

---

## 5. What Continues vs What Stopped

Moore's Law is not binary — some aspects continue while others have stopped:

**Stopped:**
- Dennard scaling (voltage scaling → power efficiency)
- Clock frequency scaling (plateaued at 3–5 GHz since 2004)
- Cost per transistor falling steeply
- Reliable 2-year cadence

**Continues (slower):**
- Transistor density still increasing (~1.5–1.7× per generation, not 2×)
- Wafer-level chip integration (chiplets, advanced packaging) adds effective transistor density
- Performance per watt still improving (through architecture + manufacturing)

**New "More Moore" innovations:**
- **CFET (Complementary FET)**: Stack nMOS on top of pMOS — 2× density with no layout change
- **GAA nanosheets**: Better electrostatics than FinFET — enables scaling to 1nm-class
- **Backside power delivery**: Move power rails to bottom of wafer → more routing on top
- **2D materials (MoS₂, WSe₂)**: Atomically thin semiconductors — transistors at ultimate physical scale

**"More than Moore":**
- Adding non-digital functionality (RF, MEMS, sensors, photonics) to silicon
- 3D stacking (SRAM on logic, DRAM on logic — Chapter 60/73)
- Heterogeneous integration (different dies for different functions)

### Quick Check
> 1. List two aspects of Moore's Law that have stopped and two that continue (slower).
> 2. What is CFET and how does it potentially double density without shrinking transistors?
> 3. What is "More than Moore" and how is it different from classical scaling?

---

## 6. Beyond Moore: Alternative Paths

With classical Moore's Law slowing, performance improvements must come from elsewhere:

**Architecture specialization (biggest current lever)**:
- CPUs with domain-specific accelerators (Neural Engine, Media Engine, Signal Processor on every SoC)
- GPUs → Tensor Cores (16× FLOPS for ML vs general FP)
- Dedicated silicon always wins vs general-purpose for a fixed function (Chapter 49 — ASICs)

**Advanced packaging / Chiplets**:
- Combine multiple dies in one package (Chapter 60/73)
- Connect SRAM, CPU, HBM in the same package with short interconnects
- Effective transistor density grows even without smaller process

**Algorithm and software improvement:**
- AlexNet (2012, CPU) → ResNet → EfficientNet: 10,000× fewer FLOPs for same accuracy at ImageNet
- Compilers, quantization, sparsity → 10–100× efficiency gains
- Software optimization often outpaces hardware scaling

**3D integration:**
- Stack SRAM on top of CPU logic (TSMC SoIC, Intel Foveros, Samsung X-Cube)
- AMD 3D V-Cache: 3× L3 cache increase via wafer-on-wafer stacking
- Eliminates cache bandwidth bottleneck that limits many workloads

**New computing paradigms:**
- In-memory computing (process data where it is stored — eliminates memory bus)
- Photonic interconnects and photonic computing (light instead of electrons)
- Neuromorphic and quantum computing (for specific workloads)

```
Performance improvement sources (approximate, 2020–2030 era):
  
  Chip scaling (Moore):      ~1.3× per generation
  Architecture:              ~1.5–2× per generation (AI accelerators)
  Advanced packaging:        ~1.2–1.5×
  Algorithm/SW:              ~1.5–5× for ML workloads specifically
  
  Combined: >3–5× system performance per generation still achievable
  Just not from one source anymore
```

### Quick Check
> 1. Name three alternative paths to performance improvement beyond classical Moore's Law.
> 2. Why is algorithm improvement sometimes more impactful than hardware improvement?
> 3. What is "3D integration" and how does it improve effective chip density?

---

## Summary

- **Moore's Law** (1965): transistor count on a chip doubles roughly every 2 years. An empirical observation that the industry maintained as a coordinated roadmap for 50+ years.
- **Cost implications**: transistor cost fell 10 billion times over 50 years. Computing became universally affordable.
- **Dennard scaling** (1974): as transistors shrink, each is faster and more power-efficient → clock speeds and power efficiency improved proportionally. Broke down ~2004 due to threshold voltage and leakage limits.
- **Slowdown**: Gate lengths approaching 2nm; voltage can't scale; wafer costs rising; 2-year cadence broken.
- **What continues**: Density still improving (1.5–1.7× per node), architecture specialization, advanced packaging, algorithms.
- **Beyond Moore**: ASICs, chiplets, 3D integration, algorithm optimization, new devices (GAA, CFET, 2D materials).

---

## Exercises

### Easy
1. What did Gordon Moore observe in 1965? Be specific about what was doubling.
2. What is Dennard scaling and how is it different from Moore's Law?
3. Give two reasons why Moore's Law is slowing down.

### Medium
4. Transistor cost calculation: Intel 4004 (1971): 4,300 transistors, die cost ~$50. Apple M4 (2024): ~20 billion transistors, die cost ~$500. (a) Calculate cost per transistor for 4004. (b) Calculate cost per transistor for M4. (c) What is the ratio? (d) The ratio of transistor counts is how many doublings? (e) If doubling period was 2 years: how many years should separate them? Does this match reality?
5. Dennard scaling math: Current process: V=1V, f=3GHz, C_gate = 1fF, transistor count = 1B. Next process: 0.7× linear scaling (k=0.7). (a) New voltage? (b) New frequency? (c) New transistor count (density scales as k²)? (d) New dynamic power per gate P = CV²f? (e) New total chip power? (f) If Dennard breaks (voltage stays at 1V): new total chip power?
6. Moore's Law inflection: The smartphone era began ~2007 (iPhone). (a) Using Moore's Law (2× transistors every 2 years), starting from iPhone's SoC (ARM11, ~100M transistors), predict transistor count in 2025 (18 years). (b) The actual 2024 Apple A18 Pro has ~16 billion transistors. Does this match Moore's Law prediction? (c) If performance/watt had scaled with Dennard scaling (unchanged), but transistor count grew as per Moore's Law: what would the CPU performance of a 2024 smartphone be vs a 2007 phone in theoretical GHz? Does this match reality?

### Hard
7. The end of Dennard scaling consequences: A 3GHz 130nm Pentium 4 core (2001) consumed ~80W. At 130nm, V = 1.5V. The 2020 Apple M1 at 5nm has ~100x more transistors per core area at 1.0V, 3.2GHz per Firestorm core, ~5W per big core. (a) Compare power density (W/mm²) between Pentium 4 (approximate die area: 145mm²) and M1 Firestorm core (estimate area ~3mm²). (b) If Dennard had continued from 130nm to 5nm (ratio ~26×): what voltage and frequency would we theoretically have? What would power density be? (c) What architectural and manufacturing innovations allowed Apple to achieve 5W/core vs the theoretical Dennard prediction? (d) Why is the Pentium 4 → M1 comparison the "end of the free lunch" story?
8. Economic analysis of Moore's Law: A startup designs a chip in 2015 (28nm, $10K/wafer, wafer holds 500 dies at 70% yield) and in 2024 (3nm, $20K/wafer, 200 dies at 60% yield). Both designs sell 1M chips. (a) 2015 chip cost per die? (b) 2024 chip cost per die (before NRE)? (c) If NRE is $5M (2015, 28nm) and $50M (2024, 3nm): amortize over 1M units for each. Total cost/chip each year? (d) Assuming the 3nm chip performs 10× better: is the cost premium justified, and for what product categories? (e) What does this analysis suggest about which applications should target cutting-edge nodes vs older, cheaper nodes?
