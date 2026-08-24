# Chapter 50: Single Crystal Alloy Compositions — Four Generations of Progress

> **"The composition of a single crystal turbine blade alloy is the result of decades of computational and experimental alloy design. Every element has a purpose: Re slows creep, Ta and W strengthen the γ matrix, Al and Ti stabilize γ′, Cr and Co provide hot corrosion resistance, Hf improves oxidation scale adhesion, and a careful balance of all of them together achieves something impossible: a blade that operates for 20,000 hours at temperatures above its own melting point, alive only because of its coatings."**

---

## Table of Contents

1. [Why Alloy Generations?](#1-why-alloy-generations)
2. [Role of Each Alloying Element](#2-role-of-each-alloying-element)
3. [First-Generation SX Alloys (1980s)](#3-first-generation-sx-alloys-1980s)
4. [Second-Generation SX Alloys (1990s) — The Re Revolution](#4-second-generation-sx-alloys-1990s--the-re-revolution)
5. [Third-Generation SX Alloys (2000s) — More Re, More Complex](#5-third-generation-sx-alloys-2000s--more-re-more-complex)
6. [Fourth and Fifth Generation Alloys — Co-Based and Beyond](#6-fourth-and-fifth-generation-alloys--co-based-and-beyond)
7. [Alloy Design Trade-Offs: Strength vs. Oxidation vs. Stability](#7-alloy-design-trade-offs-strength-vs-oxidation-vs-stability)
8. [TMS-238 and Ru-Bearing Alloys](#8-tms-238-and-ru-bearing-alloys)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Alloy Generations?

SX alloy development is driven by turbine inlet temperature (TIT) increases:
- Every 25°C increase in TIT → ~1% improvement in thermal efficiency → significant fuel savings for airline fleets
- Each alloy generation buys ~25–50°C more temperature capability

**Performance metric:**
Stress to cause 1% creep in 100 hours (σ_1%/100h) at 850°C or 982°C:
- Gen 1 (MAR-M200, B1900): ~170 MPa at 982°C
- Gen 2 (CMSX-4, PWA 1484): ~220 MPa at 982°C
- Gen 3 (CMSX-10, René N6): ~250 MPa at 982°C
- Gen 4/5 (TMS-238, MC-NG): ~280 MPa + better oxidation

**How alloy generations are classified:**
By Re (rhenium) content and additional elements:
| Generation | Re% | Ru% | Key advance |
|-----------|-----|-----|-------------|
| 0 (equiaxed/DS) | 0 | 0 | Columnar DS grain structure |
| 1 (SX) | 0 | 0 | Single crystal (no grain boundaries) |
| 2 (SX) | ~3 | 0 | Re addition → 5°C capability increase |
| 3 (SX) | ~6 | 0 | More Re, more W, Ta |
| 4 (SX) | ~3 | ~3 | Ru counteracts Re TCP precipitation |
| 5 (SX) | ~6 | ~5 | Highest capability; supply chain risk |

---

## 2. Role of Each Alloying Element

**Ni (base):** γ matrix (FCC), γ′ host (Ni₃Al/Ti), forms L1₂ structure

**Al (8–12 at%):** Primary γ′ former (Ni₃Al); increases γ′ volume fraction; forms Al₂O₃ scale for oxidation protection. Too much → lower melting point + casting difficulty.

**Ti (1–4 at%):** Secondary γ′ former; increases γ′ solvus → higher operating T; solid solution strengthener. Too much → low ductility.

**Cr (3–8 wt%):** Hot corrosion resistance (Cr₂O₃ above; reduces fluxing by Na₂SO₄); too much → lowers γ′ solvus. Compromised in Gen 3+ for creep (Cr sacrificed for Re/W).

**Co (5–12 wt%):** Raises γ′ solvus → higher γ′ fraction; solid solution strengthener; hot corrosion resistance. Balances γ/γ′ misfit.

**W (3–10 wt%):** Strong solid solution strengthener in γ; partitions to γ; slows lattice diffusion (heavy atom, reduces D); raises density. Too much → TCP precipitation (σ, μ phases).

**Ta (4–12 wt%):** Partitions strongly to γ′; slows γ′ coarsening by reducing interdiffusion; improves oxidation (via TBC bond coat adherence, reactive element effect). High density (16.7 g/cm³) increases overall blade density.

**Re (0–6 wt%):** The key Gen 2/3 addition:
- Strongly partitions to γ matrix (k ≈ 0.12 → >99% in γ)
- Reduces D_eff in γ by clustering (Re clusters slow interdiffusion)
- Reduces γ/γ′ misfit → better coherency → slower dislocation climb
- Reduces diffusion in γ→ slows Orowan bypass rate and dislocation climb rate
- MAIN MECHANISM: slows all thermally-activated deformation → creep life ×2–3
- Penalty: TCP precipitation (σ, μ, P phases at high Re); density increase; cost

**Hf (0.05–0.2 wt%):** Reactive element; improves Al₂O₃ scale adhesion (anchor pegs); improves DS/SX castability.

**Ru (0–6 wt%):** Gen 4+ addition:
- Suppresses TCP precipitation (even at high Re) → allows more Re + W
- Strongly partitions to γ (like Re)
- Adds cost: Ru ~ $15,000/kg

**Mo (0–2 wt%):** Solid solution strengthener; cheaper than Re/W; but lowers oxidation resistance (MoO₃ volatile at >700°C).

---

## 3. First-Generation SX Alloys (1980s)

No Re. Derived from DS alloys (MAR-M200, IN792) by eliminating grain boundary strengtheners (C, B, Zr — removed to avoid grain boundary cracking in SX casting).

**PWA 1480 (Pratt & Whitney, ~1980):**
Ni-10Cr-5Co-4W-1.5Mo-5Al-1.5Ti-12Ta

**René N4 (GE Aviation, ~1984):**
Ni-9Cr-8Co-6W-1.5Mo-3.7Al-4.2Ti-4.8Ta-0.5Nb

**CMSX-2 (Cannon Muskegon, ~1981):**
Ni-8Cr-5Co-8W-0.6Mo-5.6Al-1Ti-6Ta

**Properties of Gen 1:**
- σ_1%/100h at 982°C: 160–175 MPa
- Density: ~8.6 g/cm³ (high Cr, lower Re)
- γ′ volume fraction: ~60–65%
- γ/γ′ misfit: -0.2 to -0.3% (coherent)

**Why Gen 1 succeeded:**
- SX eliminates ALL grain boundaries → no grain boundary creep, no HAZ cracking
- Higher Al/Ti (no B, Zr grain boundary requirements) → higher γ′ fraction than DS alloys
- Better oxidation (higher Al activity → Al₂O₃)

---

## 4. Second-Generation SX Alloys (1990s) — The Re Revolution

Addition of ~3 wt% Re → step change in creep performance.

**CMSX-4 (Cannon Muskegon, introduced ~1989, currently most widely used HPT alloy worldwide):**
Ni-6.5Cr-9Co-6W-0.6Mo-5.6Al-1Ti-6.5Ta-3Re-0.1Hf

Composition rationale:
- 3% Re → main creep improvement (reduces D_eff in γ)
- 6.5% Cr: hot corrosion resistance maintained (vs Gen 3 compromise)
- 6.5% Ta: oxidation + creep benefit
- 0.1% Hf: Al₂O₃ scale adherence
- γ′ fraction: ~70%

**PWA 1484 (Pratt & Whitney, ~1992):**
Ni-5Cr-10Co-6W-2Mo-5.6Al-8.7Ta-3Re-0.1Hf

**René N5 (GE Aviation, ~1994):**
Ni-7Cr-7.5Co-5W-1.5Mo-6.2Al-6.5Ta-3Re-0.15Hf

**Gen 2 properties:**
- σ_1%/100h at 982°C: 210–220 MPa (+25% vs Gen 1)
- Density: ~8.7 g/cm³
- γ′ fraction: ~68–72%
- TCP tendency: starts to appear with 3% Re at high T (long service) — mitigated by Cr, Co

**Why CMSX-4 dominates production:**
- First to market with balanced composition
- Thousands of hours validation across Rolls-Royce, Pratt & Whitney, CFM engines
- Proven casting yield, heat treatment protocol, coating compatibility

---

## 5. Third-Generation SX Alloys (2000s) — More Re, More Complex

~6 wt% Re (doubled from Gen 2). Maximum Re without Ru stabilization.

**CMSX-10 (Cannon Muskegon, ~1994):**
Ni-2Cr-3Co-5W-0.4Mo-5.7Al-0.2Ti-8Ta-6Re-0.03Hf

- Note: Cr reduced to 2% (TCP suppression, density reduction) → oxidation requires better coatings
- Co reduced to 3% (unusual, reduces density)

**René N6 (GE, ~1997):**
Ni-4.2Cr-12.5Co-6W-1.4Mo-5.75Al-7.2Ta-5.4Re-0.15Hf-0.05C

- Higher Co → better TCP resistance vs CMSX-10
- 0.05%C: minor carbide strengthening retained

**Gen 3 properties:**
- σ_1%/100h at 982°C: ~245–260 MPa
- Density: ~9.0 g/cm³ (Re + Ta increase density)
- Casting difficulty: increased microsegregation of Re → harder to homogenize
- TCP risk: σ, μ phases form after long service if T profile wrong

**Key development challenge for Gen 3:**
High Re content creates:
1. Increased microsegregation on solidification (Re partitions strongly to dendrite cores)
2. Larger required heat treatment window for homogenization (requires longer HT at higher T)
3. TCP precipitation risk (Co, Cr reduced to compensate → narrows hot corrosion resistance window)

---

## 6. Fourth and Fifth Generation Alloys — Co-Based and Beyond

**Problem with Gen 3:** At 6% Re, TCP phases (σ, μ, P) form during service → embrittlement.

**Solution (Gen 4): Add Ru to suppress TCP:**

**MC-NG (SNECMA/Onera, France, ~2005):**
Ni-4Cr-4Co-5W-1Mo-6Al-0.5Ti-5Ta-4Re-4Ru-0.1Hf

- 4% Ru + 4% Re → better stability than 6% Re alone
- Used in Safran (CFM LEAP engine) blades

**TMS-138 (Japan, NIMS, 2002):**
Ni-3.2Cr-5.8Co-5.9W-0.5Mo-5.8Al-5.6Ta-4.9Re-2Ru-0.1Hf

**Gen 5 — highest performance research alloys:**

**TMS-238 (NIMS, ~2010):**
Ni-4.6Cr-5.9Co-3.6W-0.5Mo-5.9Al-7.6Ta-6.4Re-5Ru-0.08Hf

**Gen 5 properties:**
- σ_1%/100h at 1,100°C: ~140 MPa (vs Gen 2: ~80 MPa) — dramatic improvement at very high T
- Density: ~9.1–9.3 g/cm³ (heavy elements penalize density)
- Cost: Re + Ru at these levels → very expensive
- Not yet in full production; primarily for advanced military/demonstrator engines

---

## 7. Alloy Design Trade-Offs: Strength vs. Oxidation vs. Stability

**The fundamental trade-offs:**

```
                   High creep strength
                         ↑
               Re, W, Ta, Mo
                         |
High oxidation ←─ High Cr, Al ─→ Too much reduces γ′ solvus → less strengthening
  resistance
                         |
                     High Re
                         |
            TCP phases form ←─── Need Ru or Cr to suppress
                                 (but Ru = expensive; Cr reduces creep)
```

**Key alloy design tool:** PHACOMP (Phase Computation):
- Calculates "electron vacancy number" N_v of the alloy
- If N_v > threshold → TCP phases form
- Used to screen compositions before casting

**Density-corrected specific strength:**
Since heavier elements (Re, W, Ta, Ru) increase density (and centrifugal stress = ρ × ω² × r), performance is sometimes reported as σ_creep / ρ (specific creep strength).

**Newer approach: ICME + CALPHAD:**
Use computational thermodynamics (Ch 64) to predict:
- Phase stability at temperature
- γ′ volume fraction
- Partitioning of each element
- TCP tendency (N_v equivalent from thermodynamics)

---

## 8. TMS-238 and Ru-Bearing Alloys

The science behind TMS-238 (one of the best characterized Gen 5 alloys):

**γ/γ′ misfit:** δ = -0.5% (negative; γ matrix has larger lattice than γ′)
- Under tensile creep stress → directional rafting perpendicular to [001]
- Negative misfit + tensile stress → N-type rafting
- Controlled rafting → improved creep resistance in primary stage

**Re clustering in γ:**
EXAFS and atom probe tomography show Re forms short-range ordered Re-rich clusters (diameter ~0.5 nm) in γ matrix:
- Clusters resist dislocation passage → additional strengthening
- Slow dislocation climb rate (diffusivity reduced by Re clusters)
- Slower dissolution of γ′ → better coarsening resistance

**Ru effect on TCP suppression:**
Ru strongly partitions to γ, like Re. But Ru actually REDUCES the tendency of Re to form TCP phases:
- Mechanism: Re + Ru → maintain γ phase stability (Ru provides electrons to balance Re electron structure)
- Result: can have 6% Re + 5% Ru → higher combined solid solution strengthening than 6% Re alone, without TCP

---

## Summary

| Generation | Representative | Re% | Ru% | σ_creep at 982°C | Key Innovation |
|-----------|---------------|-----|-----|-----------------|----------------|
| 1st | CMSX-2, PWA1480 | 0 | 0 | 175 MPa | Single crystal (no grain boundaries) |
| 2nd | CMSX-4, PWA1484, N5 | 3 | 0 | 220 MPa | Re addition: ×2.5 creep life |
| 3rd | CMSX-10, René N6 | 6 | 0 | 250 MPa | More Re: +25°C capability |
| 4th | MC-NG, TMS-138 | 3–5 | 2–4 | 265 MPa | Ru suppresses TCP at high Re |
| 5th | TMS-238 | 6+ | 5+ | ~280 MPa | Maximum creep; highest density |

---

## Exercises

1. Compare CMSX-4 (Gen 2) and CMSX-10 (Gen 3): (a) Both have nominally similar γ′-forming elements (Al, Ti, Ta). CMSX-10 has 6% Re vs CMSX-4's 3% Re. Using the Dorn equation ε̇_ss = A×σⁿ×exp(-Q_c/RT): if doubling Re reduces D_eff by 40%, estimate the reduction in ε̇_ss. (Assume n = 4 for dislocation creep; σ, T constant). (b) CMSX-10 has only 2% Cr vs CMSX-4's 6.5% Cr. At 750°C (Type II hot corrosion range), what risk does lower Cr pose? (c) Calculate the density of CMSX-4 (given: Ni=8.9, Cr=7.2, Co=8.9, W=19.3, Mo=10.2, Al=2.7, Ti=4.5, Ta=16.7, Re=21.0 g/cm³ densities, composition in wt%). Why does density increase from Gen 1 to Gen 3?

2. Element partitioning in CMSX-4: (a) After directional solidification, element partitioning between dendrite core and interdendritic region is characterized by k (partition coefficient). For CMSX-4: k_Re ≈ 0.72, k_W ≈ 0.79, k_Ta ≈ 1.35, k_Cr ≈ 1.18. Which elements segregate to dendrite cores? Which to interdendritic? (b) The as-cast homogenization treatment at 1,290°C/6h is required to homogenize Re and W. How does solutioning diffuse Re from the core? What homogenization time would be needed if you dropped the temperature to 1,250°C? (Use D_Re ∝ exp(-Q_c/RT) with Q = 280 kJ/mol.) (c) What happens if homogenization is incomplete (Re-rich dendrite cores remain)? What service failure mode results?

3. TCP phase stability analysis: Using PHACOMP-like analysis, an alloy with N_v > 4.55 is predicted to form TCP phases. Calculate N_v for: (a) CMSX-4: N_v = ΣX_i × V_i where V values: Ni=0, Cr=4.7, Co=3.7, W=6.7, Mo=4.7, Al=0, Ti=0, Ta=4, Re=5.5 (use wt fractions and V values approximately). (b) CMSX-10: recalculate with its composition. (c) Based on N_v criterion, which alloy is more susceptible to TCP formation? (d) How does adding Ru (V_Ru ≈ 5.0) to CMSX-10 reduce TCP tendency? (note: mechanism is more subtle than PHACOMP — briefly explain).

4. Density penalty from alloying: A turbine blade runs at 10,000 rpm, mid-span radius r = 0.3 m. (a) Calculate the centrifugal acceleration (a_c = ω²r) in units of g (gravitational acceleration). (b) The blade weighs 250g with Gen 2 alloy (density 8.7 g/cm³). Gen 5 alloy has density 9.3 g/cm³. Calculate new blade mass. (c) Centrifugal stress at blade root = m × a_c / A_root where A_root = 400 mm². Calculate the centrifugal stress for Gen 2 and Gen 5 blades. (d) If the disk material's bearing stress limit is 650 MPa, can both alloys be used? What blade design change could compensate for the Gen 5 density increase?

5. Economic analysis of alloy generations: (a) Re current market price ~$3,000/kg; Ru ~$15,000/kg. CMSX-4 blade weighs 250g with ~3 wt% Re. Calculate Re content and cost per blade. (b) Gen 5 alloy has 6% Re + 5% Ru. Calculate additional precious metal cost per blade vs CMSX-4. (c) A fleet of 2,000 aircraft × 2 engines × 60 HPT blades per engine = total blade count. Calculate total precious metal cost difference (Gen 5 vs CMSX-4) for the entire fleet (assume entire fleet replacement). (d) Gen 5 blade lasts 40,000 cycles vs Gen 2's 25,000 cycles. At $8,000/blade replacement cost (labor + blade), calculate total fleet maintenance cost over 100,000 flight cycles for each. Does the extra alloy cost justify the longer life?
