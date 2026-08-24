# Chapter 56: Process Nodes — 5nm, 3nm, 2nm

When you hear "TSMC 3nm" or "Intel 20A," you are hearing a **process node name** — a shorthand for a specific manufacturing generation with specific transistor density, performance, and power characteristics. These names no longer mean the physical gate length is 5nm or 3nm. They are marketing labels that loosely indicate a generation of technology. This chapter explains what each node actually delivers, compares the leading foundries (TSMC, Samsung, Intel), traces the path from 28nm to 2nm, and looks at what comes after.

## Table of Contents

1. [What a Process Node Actually Means](#1-what-a-process-node-actually-means)
2. [The Leading Nodes: 7nm to 2nm](#2-the-leading-nodes-7nm-to-2nm)
3. [TSMC, Samsung, Intel — Who Is Leading?](#3-tsmc-samsung-intel--who-is-leading)
4. [28nm and Mature Nodes — Still Important](#4-28nm-and-mature-nodes--still-important)
5. [Gate-All-Around at 3nm/2nm](#5-gate-all-around-at-3nm2nm)
6. [What Comes After 2nm?](#6-what-comes-after-2nm)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. What a Process Node Actually Means

Before ~2012, process nodes were named after the physical gate length of the transistor. A "130nm process" really did have ~130nm gate length. After 2012, gate lengths stopped decreasing at the same pace as the node name, and the names became marketing shorthand.

**What the node name correlates with today:**
1. **Transistor density**: measured in Million Transistors per mm² (MTr/mm²) — the most physically meaningful metric
2. **Performance**: estimated clock frequency at a given power
3. **Power**: leakage current and dynamic power at that frequency

```
Node name vs actual metrics (TSMC):
  
  Node    MTr/mm²    Metal Pitch    Gate Pitch    Year
  28nm    40          90nm           108nm         2012
  16nm    65          80nm           90nm          2014
  10nm    100         64nm           66nm          2016
  7nm     91*         40nm           57nm          2018 (*complex cells)
  5nm     173         30nm           48nm          2020
  4nm     193         28nm           45nm          2022
  3nm     292         21nm           48nm          2022
  N2      ~400        est.           est.          2025

  * "7nm" metal pitch of 40nm is much larger than 7nm
```

**The gap between name and reality**: TSMC N5 has a contacted poly pitch (CPP, the transistor gate-to-gate distance) of 48nm — about 10× larger than the "5nm" name. The name was chosen to be competitive-sounding relative to Intel's simultaneous "10nm" node (which had similar density to TSMC 7nm). Node names became competitive branding.

**Why density metrics matter more than names**: Apple's M3 chip uses TSMC 3nm. NVIDIA's H100 uses TSMC 4nm. The H100 has ~80 billion transistors on ~814 mm² = ~98 MTr/mm². Apple M3 Ultra has ~192 MTr/mm². Comparing "3nm vs 4nm" without density numbers is meaningless.

### Quick Check
> 1. Why don't node names (like "5nm") correspond to the physical gate length anymore?
> 2. What metric is more meaningful than the node name for comparing process generations?
> 3. What does CPP (Contacted Poly Pitch) measure?

---

## 2. The Leading Nodes: 7nm to 2nm

**7nm (TSMC N7, 2018; Intel 10nm, 2019; Samsung 8nm)**:
- First node to require multi-patterning extensively
- TSMC N7 used ArF immersion + multi-patterning (no EUV)
- TSMC N7+ (2019): first EUV usage in production (limited layers)
- Density: ~91–100 MTr/mm²
- Products: iPhone XS (A12), AMD Ryzen 3000 (Zen 2), NVIDIA RTX 2080

**5nm (TSMC N5, 2020; Samsung 5nm, 2020)**:
- EUV adopted for multiple layers (TSMC N5: ~14 EUV layers)
- ~1.8× density vs N7
- FinFET transistors
- Products: Apple A14, Qualcomm Snapdragon 888, TSMC's own reference cells

**4nm (TSMC N4, 2022; Samsung 4nm, 2022)**:
- Incremental refinement of 5nm process
- ~1.1× density vs 5nm, better power
- Products: Apple A15 (N5), NVIDIA H100 (N4), Qualcomm Snapdragon 8 Gen 2 (N4P)

**3nm (TSMC N3, 2022; Samsung SF3, 2022)**:
- TSMC N3: FinFET at 3nm (not GAA like Samsung)
- Samsung SF3: first mass production with GAA (Multibridge Channel FET)
- TSMC N3E/N3P/N3X: enhanced versions for different power/performance trade-offs
- Density: ~292 MTr/mm² (N3E)
- Products: Apple A17 Pro, Apple M3 family (N3B — custom TSMC variant for Apple)

**2nm (TSMC N2, 2025; Intel 20A/18A, 2024-25)**:
- TSMC N2: First TSMC node with GAA (nanosheet transistors)
- Intel 18A: GAA ("RibbonFET") + backside power delivery (PowerVia), scheduled 2025
- Expected density: ~350–500 MTr/mm²
- Products: Apple A18, future Intel Core 2 generations

### Quick Check
> 1. At which node did TSMC first use EUV in production?
> 2. Which foundry first introduced GAA transistors in mass production and at what node?
> 3. What is TSMC N3B and why did Apple get a special variant?

---

## 3. TSMC, Samsung, Intel — Who Is Leading?

The three leading-edge foundries each have different strategies:

**TSMC** (Taiwan Semiconductor Manufacturing Company):
- **Business model**: Pure-play foundry. TSMC makes chips for Apple, NVIDIA, AMD, Qualcomm, Intel. Does not design chips.
- **Revenue**: ~$70B annually (2024). Makes ~90% of the world's most advanced chips.
- **Technology**: Ahead on density and yield. N3E and N2 are industry-leading.
- **Customers**: Apple (~25% of revenue), NVIDIA, AMD, Qualcomm, Intel Foundry
- **Risk**: Geographic concentration in Taiwan. TSMC Arizona (Fab 21) breaking ground; ~$65B investment, producing N4P (2024) and N2 (2026).

**Samsung Foundry**:
- **Business model**: IDM (Integrated Device Manufacturer) + foundry. Makes Exynos chips for Samsung phones AND foundry services for others.
- **Technology**: First to ship GAA at 3nm, but yield was initially poor. 4nm yield issues caused Qualcomm to move Snapdragon 8 Gen 2 from Samsung to TSMC.
- **Customers**: Google Tensor, Qualcomm (some products), NVIDIA (older products)
- **DRAM**: Samsung leads in DRAM manufacturing (not covered by node naming — different tech)

**Intel Foundry** (formerly Intel Foundry Services):
- **History**: Intel was the pioneer — built the first commercial transistors, led at every node from 1970s–2015. Then 10nm delays (announced 2016, mass production 2019) allowed TSMC/Samsung to leap ahead.
- **Recovery**: Intel Foundry is Intel's attempt to regain manufacturing leadership with 18A (GAA + backside power). Intel is also an ASML customer for High-NA EUV machines.
- **Customers**: Intel's own products (Core, Xeon), plus external foundry customers (Amazon AWS Graviton processors being evaluated on Intel 18A)

```
Leading-edge market share (~2024):
  TSMC:     ~60% of foundry revenue, >90% of sub-7nm wafers
  Samsung:  ~15% of foundry revenue
  Intel:    ~5% foundry revenue (mostly Intel's own products)
  GlobalFoundries, SMIC, others: remaining (at older nodes)
```

### Quick Check
> 1. What is the difference between a pure-play foundry (TSMC) and an IDM (Intel)?
> 2. Why did Qualcomm move some production from Samsung to TSMC?
> 3. What two new technologies does Intel 18A introduce?

---

## 4. 28nm and Mature Nodes — Still Important

Not every chip needs the most advanced node. The economics of mature nodes are entirely different, and they supply the vast majority of chips by volume.

**28nm** (TSMC introduced 2012):
- Still the most popular "value" node for cost-sensitive applications
- Simple process: planar MOSFET, single patterning, no EUV
- Wafer cost: ~$3,000–4,000 (vs $20,000 for 3nm)
- Used for: automotive ICs, microcontrollers, IoT chips, analog mixed-signal, display drivers, power management ICs, Wi-Fi chips

**Why 28nm is so popular:**
- Fully depreciated fab equipment (fabs built in 2012 are paid off)
- Mature process with high yield and known behavior
- Massive installed capacity
- Multi-year supply contracts → stable pricing
- Many applications don't need smaller transistors (power ICs, analog, RF)

**Automotive**: Car chips must pass AEC-Q100 qualification (automotive reliability), survive -40°C to +150°C, have 15-year supply guarantees. Advanced nodes have shorter production lifespans. Most automotive chips use 28nm–65nm.

**COVID chip shortage (2020–2022)**: The chip shortage was NOT about advanced 5nm chips — it was about 28nm, 40nm, 65nm mature node chips. Auto makers canceled orders in early COVID, foundries allocated capacity to electronics, then auto demand surged. It took 18–24 months to rebalance because adding mature node capacity requires building new fabs that cost $3–5B and take 18–24 months.

**GlobalFoundries, SMIC, UMC**: These foundries focus on 12nm–28nm and older, serving the mature market. Not competing at 3nm/5nm. SMIC is China's largest foundry, currently at 14nm, working toward 7nm (using multi-patterning without EUV).

### Quick Check
> 1. Why is 28nm still widely used despite 3nm being available?
> 2. Why was the 2020–2022 chip shortage primarily about mature nodes, not advanced nodes?
> 3. Which foundries focus on mature nodes?

---

## 5. Gate-All-Around at 3nm/2nm

The **Gate-All-Around (GAA)** transistor was introduced (in production) at Samsung's 3nm (SF3) and TSMC's N2. This is the successor to FinFET.

**Why GAA is needed at 2nm and below:**
FinFET wraps the gate around three sides of the fin. At very small gate lengths (<12nm), the depletion from source and drain reaches through the thin fin even with a 3-sided gate — the transistor cannot fully turn off.

GAA wraps the gate on ALL FOUR sides of a horizontal silicon nanosheet (or multiple stacked nanosheets), providing complete electrostatic control.

```
FinFET → GAA transition:

FinFET (cross-section through gate):     GAA nanosheet (cross-section):
  
  ┌─────────────────────────────┐       ──────── Gate ────────
  │         Gate                │       │      nanosheet 3    │
  │  ┌─────────────────────┐   │       │      nanosheet 2    │
  │  │       Fin            │   │       │      nanosheet 1    │
  │  └─────────────────────┘   │       ──────── Gate ────────
  └─────────────────────────────┘
  
  Gate wraps 3 sides             Gate wraps 4 sides of EACH nanosheet
  (top + 2 sidewalls)            Multiple nanosheets stacked
```

**TSMC N2 GAA (nanosheet):**
- Stacked nanosheet channels: typically 2–4 nanosheets, each ~5nm thick, ~20–25nm wide
- Gate length: ~12nm
- Benefits vs FinFET: better electrostatics → lower V_th → lower Vdd → lower power
- Performance: ~10–15% faster or ~25–30% lower power vs N3 FinFET

**Intel RibbonFET (GAA):**
- Intel's name for their GAA implementation (Intel 20A/18A)
- "Ribbon" nanosheets, stacked 2–4
- Combined with PowerVia (backside power delivery) for additional performance

**CFET (Complementary FET) — future:**
Stack an nMOS nanosheet on top of a pMOS nanosheet. This effectively combines the two transistors of a CMOS pair into a single vertical structure, doubling cell density. Expected at ~1nm-class nodes (2030+).

```
CFET concept:
  ┌─────────────────────────────┐
  │  pMOS nanosheets (top)      │  VDD
  │────────────────────────────│
  │  nMOS nanosheets (bottom)   │  GND
  └─────────────────────────────┘
  
  Both transistors share the same footprint → 2× density
```

### Quick Check
> 1. Why does FinFET fail at very small gate lengths, requiring GAA?
> 2. What is the key advantage of GAA over FinFET?
> 3. What is CFET and when is it expected to appear?

---

## 6. What Comes After 2nm?

The roadmap beyond 2nm requires a combination of new materials, new transistor designs, and 3D integration:

**1nm-class process nodes (~2027–2030):**
- TSMC A16 (Angstrom 16, announced 2024): 1.6nm-class, GAA nanosheets + backside power delivery
- Intel 14A (2027): Next after 18A
- 2D materials: MoS₂ transistors with 3-atom-thick channels; demonstrated in research, not production

**Backside power delivery:**
- In current chips, power rails (VDD/GND) are routed on the BEOL metal layers
- This competes for routing space with signal wires
- Backside power: move power distribution to a new metal layer on the WAFER BACKSIDE
- Frees up >20% of front-side routing → better performance and density
- Intel PowerVia (with 18A), TSMC A16 include backside power

**High-NA EUV impact:**
- ASML EXE:5000 (NA=0.55) machines arriving 2024–2026
- Single exposure down to ~8nm — eliminates multi-patterning for these layers
- Reduces process steps, improves yield

**2D materials (research):**
- MoS₂, WSe₂: single-layer transistors (3 atoms thick)
- MIT demonstrated 1nm MoS₂ transistors in 2021
- Challenge: cannot yet grow uniformly at wafer scale
- IBM: integrating 2D materials into CMOS process research

**Heterogeneous integration replaces pure scaling:**
As classical transistor scaling slows, the industry will increasingly rely on combining different chips optimized for different functions in advanced packages (chiplets, 3D stacking, wafer-on-wafer). This is covered in detail in Chapter 60 and 73.

### Quick Check
> 1. What is backside power delivery and what problem does it solve?
> 2. What are 2D materials in semiconductor context and what are the challenges?
> 3. What is TSMC A16 and what technologies does it introduce?

---

## Summary

- **Process node names** (5nm, 3nm) are marketing labels, not physical measurements. Transistor density (MTr/mm²) and contacted poly pitch (CPP) are the real metrics.
- **7nm to 3nm**: major leap, EUV adopted, density roughly doubled each generation.
- **3nm/2nm**: transition from FinFET to GAA (nanosheet/nanoribbon) transistors.
- **Leading foundries**: TSMC (~90% of sub-7nm volume, N2 in 2025), Samsung (GAA first but yield issues), Intel (recovering with 18A).
- **Mature nodes** (28nm): still critical for automotive, IoT, analog — cost matters more than performance.
- **Beyond 2nm**: backside power delivery, High-NA EUV, CFET, 2D materials, heterogeneous integration.

---

## Exercises

### Easy
1. What does "5nm" actually mean in a modern chip process node?
2. Why are mature nodes (28nm) still important even though 3nm is available?
3. What is GAA (Gate-All-Around) and how is it different from FinFET?

### Medium
4. Density comparison: TSMC N5 has 173 MTr/mm². TSMC N3E has 292 MTr/mm². (a) What is the density ratio N3/N5? (b) Apple M3 uses N3B. If the M2 die was 234mm² (N5), and M3 achieves same transistor count with 1.69× density: what is the M3 die area? (c) Smaller die means more dies per wafer: at 173 MTr/mm², die area 234mm². At 292 MTr/mm², die area 138mm². 300mm wafer: approximately how many more dies does M3 get per wafer? (Approximate: wafer area 70,685 mm², ignore edge loss.)
5. Process node cost modeling: A startup designs a chip requiring 20 billion transistors. (a) At TSMC N7 (91 MTr/mm²): what minimum die area? (b) At TSMC N5 (173 MTr/mm²): minimum die area? (c) Wafer costs: N7=$12K, N5=$17K. At 70% yield: cost per die for each? (d) NRE: N7=$30M, N5=$40M. For 1M units: total cost/chip for each? (e) Which node is cheaper for 1M units? What about 100K units?
6. FinFET vs GAA performance: A FinFET at TSMC N3 runs at 1V, achieves 3.5 GHz, leakage current 5nA/µm width. Switching to GAA N2: Vdd can be reduced to 0.75V (better electrostatics). (a) Using P = CV²f: if frequency stays at 3.5 GHz, what is the dynamic power ratio N2/N3? (b) If frequency increases to 4.5 GHz at 0.75V: dynamic power ratio vs N3 at 3.5 GHz/1V? (c) Subthreshold leakage scales as: I_leak ∝ exp(-Vth/VT) where VT = 26mV. If GAA allows Vth reduction from 0.35V to 0.25V: leakage ratio? (d) Why does lower Vth improve performance but worsen leakage, and how does GAA help manage this trade-off?

### Hard
7. TSMC vs Samsung vs Intel comparative analysis: In 2024, Intel is using Intel 4 process (7nm-class). TSMC is at N3E. Samsung at SF4. Analyze: (a) Obtain or estimate MTr/mm² for each (research or estimate from CPP/FP data). (b) Intel 4 delivers Intel Core Ultra with 125W TDP at N4-class. TSMC N3 delivers Apple M3 Max at 35W TDP. Why is there such a large difference despite similar transistor density? (hint: consider architecture, use case, thermal design, not just process) (c) Intel 18A promises 1.07× density improvement over Intel 3, plus RibbonFET + PowerVia. What technical claim is Intel making about relative performance vs TSMC N2? Is this plausible? (d) From a customer perspective (like AWS): why might you choose Intel Foundry for custom chips despite TSMC's track record?
8. Chinese semiconductor catch-up analysis: SMIC achieved 7nm-class chip production (Huawei Kirin 9000S, 2023) using ArF immersion + multi-patterning without EUV. (a) What multi-patterning technique (SADP, SAQP) would be needed to achieve 7nm-class pitch with 193nm immersion (NA=1.35, k₁=0.28)? (b) What is the additional number of process steps, cycle time penalty, and yield impact of this approach vs TSMC's EUV-based 7nm? (c) Can SMIC reach 5nm-class without EUV? What pitch would SASQP (self-aligned sextuple patterning) achieve? What are the practical limits? (d) The US government restricted ASML from selling even ArF immersion scanners (DUV) to China in late 2023. What impact does this have on SMIC's ability to expand 7nm capacity?
