# Chapter 43: Nickel Superalloy Composition — What Every Element Does

> **"A modern single-crystal nickel superalloy contains 10–12 elements. Every atom matters. Cr gives oxidation resistance but too much destabilizes γ′. Re dramatically improves creep but makes the alloy 20% denser and costs $3,000/kg. Ta strengthens γ′ but slows diffusion. Ru costs even more than Re but suppresses the TCP phases that Re causes. Each generation of alloy is a new negotiation between these competing demands."**

---

## Table of Contents

1. [The Composition Complexity — Why 10+ Elements?](#1-the-composition-complexity--why-10-elements)
2. [The γ Matrix — Primary Solid Solution Elements](#2-the-γ-matrix--primary-solid-solution-elements)
3. [The γ′ Phase — Precipitate-Forming Elements](#3-the-γ-phase--precipitate-forming-elements)
4. [Carbide-Forming Elements and Their Roles](#4-carbide-forming-elements-and-their-roles)
5. [Refractory Elements — The High-Temperature Strengtheners](#5-refractory-elements--the-high-temperature-strengtheners)
6. [Rhenium — The Transformative Addition](#6-rhenium--the-transformative-addition)
7. [Ruthenium — The Stabilizer](#7-ruthenium--the-stabilizer)
8. [Grain Boundary Elements (PC Alloys Only)](#8-grain-boundary-elements-pc-alloys-only)
9. [Harmful Phases — TCP and What Causes Them](#9-harmful-phases--tcp-and-what-causes-them)
10. [CMSX-4: A Complete Worked Example](#10-cmsx-4-a-complete-worked-example)
11. [Alloy Design Principles — The CALPHAD Approach](#11-alloy-design-principles--the-calphad-approach)
12. [Generational Comparison Table](#12-generational-comparison-table)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. The Composition Complexity — Why 10+ Elements?

A turbine blade must simultaneously:
- Resist creep for 30,000 hours at 980°C, 200 MPa
- Resist oxidation at up to 1100°C (with coating)
- Resist fatigue from 10,000+ thermal cycles
- Maintain a stable microstructure (no γ′ coarsening, no TCP phases)
- Be castable as a single crystal (no hot cracking, good fluidity)
- Have density low enough that centrifugal stress stays within limits

No single element can do all of this. A multi-element alloy is necessary, with each element making a specific contribution. The art of superalloy design is optimizing all elements simultaneously — adding any one too much harms something else.

**Structure of contributions:**

```
Function                →  Elements responsible
─────────────────────────────────────────────────────
γ matrix solid solution → Cr, Co, Mo, W, Re, Ru
γ′ stability & strength → Al, Ti, Ta, Nb
Oxidation resistance    → Cr, Al (Al₂O₃ scale formation)
Hot corrosion resistance→ Cr (primary)
Creep resistance        → Re, W, Mo, Ta (via multiple mechanisms)
Density reduction       → Al, Ti (low atomic weight)
Stability (anti-TCP)    → Ru, Co (reduce TCP tendency)
GB strengthening (PC)   → B, Zr, C, Hf
```

---

## 2. The γ Matrix — Primary Solid Solution Elements

The γ phase is a disordered FCC nickel-rich solid solution. Most alloying elements partition preferentially to γ rather than γ′.

### Chromium (Cr) — 5–15 wt%

**Roles:**
- **Oxidation resistance**: Cr₂O₃ scale forms at moderate temperatures; promotes Al₂O₃ formation by reducing O₂ partial pressure at the alloy surface (the "getter effect")
- **Hot corrosion resistance**: Cr is essential for Type I hot corrosion resistance (Na₂SO₄ deposit attack at 800–900°C). < 5% Cr → vulnerable to hot corrosion.
- **Solid solution strengthening** of γ matrix (moderate effect)

**Too much Cr:** Destabilizes γ/γ′ equilibrium at high temperatures → promotes TCP phases (σ phase, μ phase) → embrittlement. Modern SX alloys have reduced Cr (5–8%) vs. early PC alloys (15–20%) because: (1) coatings provide oxidation and corrosion protection now, and (2) high T operation makes TCP stability critical.

### Cobalt (Co) — 5–15 wt%

**Roles:**
- Partitions to γ matrix; reduces γ′ solvus temperature (allows solution treatment at higher temperature → more complete γ′ dissolution → better homogenization)
- Increases γ′ volume fraction by reducing Ni available to form matrix
- Improves hot corrosion resistance slightly
- Reduces stacking fault energy → affects deformation mechanism
- **Anti-TCP**: Co substitutes for Ni in γ and reduces the tendency to form TCP phases (Mo-rich σ phase, Cr-rich μ phase) relative to equivalent Re additions in Ni-only matrix

### Molybdenum (Mo) — 0–3 wt%

Strong solid solution strengthener in γ (large misfit with Ni lattice → strong elastic distortion → dislocation drag).

**Too much Mo:** Strong TCP phase former (σ phase). Mo-rich alloys need careful balance with Cr and W. In modern SX alloys, Mo is often kept < 2% with W doing more of the solid solution strengthening.

---

## 3. The γ′ Phase — Precipitate-Forming Elements

The γ′ phase (Ni₃Al-based, ordered L1₂ structure) is the primary strengthening phase. These elements all partition preferentially to γ′:

### Aluminum (Al) — 5–7 wt%

**Primary γ′ former**: Al is the most important γ′-forming element. Al substitutes for Ni in the Ni₃Al structure (Al sits on the corners of the L1₂ unit cell).

- Higher Al → more γ′ (higher volume fraction)
- Higher Al → better oxidation resistance (Al₂O₃ scale — the most protective oxide)
- γ′ volume fraction in modern SX alloys: 60–70% — remarkably high for a two-phase alloy

**Too much Al:** Alloy becomes difficult to cast (high viscosity), γ′ forms during solidification (not during aging), and the alloy can be brittle.

### Titanium (Ti) — 0–2 wt%

Substitutes for Al in γ′ → Ni₃(Al,Ti). Increases γ′ volume fraction and somewhat increases γ′ antiphase boundary energy (APB energy → APB energy determines how hard it is for dislocations to cut through γ′).

Higher Ti → higher APB energy → harder γ′ → better creep resistance at lower temperatures (< 800°C).

**Too much Ti:** Promotes η phase (Ni₃Ti, hexagonal) which is not as beneficial as γ′.
In modern SX alloys (designed for >950°C service), Ti is often < 1% or eliminated — Ta does a better job at high temperatures.

### Tantalum (Ta) — 5–9 wt%

**The key refractory γ′ former.** Ta substitutes strongly for Al in γ′ → Ni₃(Al,Ta).

- **Strengthens γ′** more than Ti (larger atom → larger misfit → higher APB energy → more resistance to dislocation cutting)
- **Slows diffusion**: Ta has very low diffusivity in Ni → significantly retards γ′ coarsening (Ostwald ripening rate ∝ diffusivity) → γ′ stays fine longer at temperature
- Large atomic radius → creates lattice strain in γ′ → solid solution strengthening OF γ′ itself
- Slightly negative lattice misfit contribution (helps achieve optimal negative δ in modern alloys)

Ta is the most important element in 2nd and later generation SX alloys for balancing strength and stability.

**Too much Ta:** Alloy density increases significantly (Ta = 16.6 g/cm³). Centrifugal stress in the blade is:
```
σ_centrifugal = ρ × ω² × r × L
```
Higher density → higher centrifugal stress → risk of blade separation. Modern alloys balance Ta content to avoid density > 8.7–8.9 g/cm³.

### Niobium (Nb) — Primarily in Fe-base alloys

In Ni-base SX alloys, Nb is largely replaced by Ta (stronger effect per atom). In IN718 (Fe-Ni base), Nb forms γ″ (Ni₃Nb, D0₂₂ structure) which is the primary strengthening phase.

---

## 4. Carbide-Forming Elements and Their Roles

### Carbon (C) — 0–0.15 wt%

In **polycrystalline** alloys, C forms carbides at grain boundaries:
- **MC carbides**: TiC, TaC, NbC, HfC — form at high temperature during solidification; block dislocation motion; good for creep
- **M₂₃C₆ carbides**: (Cr,Mo)₂₃C₆ — form at grain boundaries during aging; pin grain boundaries; good for creep and fatigue at grain boundaries
- **M₆C carbides**: W- and Mo-rich; secondary strengtheners

In **single crystal** alloys: no grain boundaries, so carbides provide no benefit. C content is reduced to < 0.005 wt% in modern SX alloys (compared to 0.05–0.15 wt% in PC alloys). Lower C → higher solidus → better solution heat treatment.

### Hafnium (Hf) — 0.1–1.5 wt% (in some alloys)

Hf partitions to grain boundaries in PC alloys → significant improvement in ductility and resistance to intergranular fracture. Also:
- Promotes fine carbide distribution
- Improves oxidation resistance (forms HfO₂ which pegs the α-Al₂O₃ scale → less scale spallation)
- Can cause casting difficulties (reactive, high melting point)

In SX alloys: Hf is sometimes added in small amounts for its oxidation resistance benefit, but grain boundary benefit is absent.

---

## 5. Refractory Elements — The High-Temperature Strengtheners

### Tungsten (W) — 4–12 wt%

One of the most effective solid solution strengtheners in γ. W sits in the Ni FCC matrix and creates large elastic distortion (W atom is 12% larger than Ni) → strong interaction with dislocation stress fields → dislocation drag → resistance to creep.

W also:
- Partitions to γ′ → strengthens γ′ (smaller effect than solid solution in γ)
- Increases liquidus temperature → allows higher use temperature
- Very low diffusivity → slows all diffusion-controlled processes (grain growth, phase coarsening, oxidation)

**Too much W:** Alloy density increases significantly (W = 19.3 g/cm³). σ/μ TCP phase formation risk.

---

## 6. Rhenium — The Transformative Addition

Rhenium (Re) is the most important single element addition in the history of superalloy development. The jump from 1st generation to 2nd generation SX alloys (PWA 1480 → CMSX-4) was driven entirely by adding 3 wt% Re.

**Re properties:** Melting point 3186°C (highest of all metals except W and Os). Atomic radius 137 pm (slightly larger than Ni's 124 pm). Very high bulk modulus. Extremely low diffusivity in Ni.

### Mechanism 1: Solid Solution Strengthening

Re is a very potent solid solution strengthener — among the most potent of all elements added to Ni. The combination of:
- Large size misfit with Ni
- Interaction between Re electrons and Ni d-band (electronic contribution)
- Very low diffusivity (Re atoms are "slow" → strong drag on dislocations trying to climb)

Results in dramatically increased matrix creep resistance.

### Mechanism 2: Retardation of Dislocation Climb

As discussed in Chapter 16, secondary creep in superalloys is limited by **dislocation climb** — dislocations in the γ channels climb over γ′ precipitates via diffusion.

Re has the lowest diffusivity of any common superalloy element in Ni:
```
Diffusivity ranking at 1050°C (low → high, high = fast diffuser):
Re < W < Mo < Ta < Cr < Co < Al < Ti < Ni (self-diffusion)
```

Re atoms segregate to dislocation cores (Cottrell-type atmosphere). When the dislocation tries to climb, it must drag its Re atmosphere with it. Re's low diffusivity means this drag is enormous → dislocation climb rate is dramatically reduced → creep rate is reduced.

### Mechanism 3: γ′ Coarsening Suppression

As shown in Chapter 44, γ′ coarsening (Ostwald ripening) follows:
```
r³(t) - r₀³ = K × t × exp(-Q_coarsen/RT)
```

K depends on the diffusivity of the rate-limiting element. Re (low D) is rejected from the γ′ and forms a diffusion barrier in the γ channel between precipitates. Atoms trying to diffuse from small γ′ particles to large ones must pass through Re-enriched γ channels → slowed coarsening → γ′ stays fine longer at temperature → maintained creep resistance.

### Re's Cost and Supply Problem

Rhenium:
- Global production: ~60 tonnes/year (2024)
- Cost: $2,500–4,000/USD per kg
- CMSX-4 contains 3 wt% Re; 4th+ generation alloys contain 5–6 wt% Re
- 80%+ of global Re production goes to superalloys

A commercial twin-engine wide-body aircraft (e.g., Boeing 777) has approximately 2 × 80 HPT blades × ~0.5 kg/blade = 80 kg of HPT blades. At 3% Re, that's 2.4 kg of Re per aircraft.

The global fleet of commercial aircraft (27,000+ aircraft) plus military plus industrial gas turbines represents many hundreds of tonnes of Re locked in service blades. This makes Re supply a national security concern.

---

## 7. Ruthenium — The Stabilizer

By the time 3rd generation alloys reached 6 wt% Re (CMSX-10, René N6), a new problem appeared: **TCP phase formation** (topologically close-packed phases: σ, μ, P phases). These brittle intermetallic phases precipitate within the blade after extended time at temperature, consuming γ′ formers and causing embrittlement.

The cause: high Re+W+Mo+Cr content makes the alloy composition inherently unstable toward TCP formation.

Ruthenium (Ru) was discovered in the 1990s to significantly suppress TCP formation while allowing Re to remain at high levels. Ru:
- Partitions to the γ matrix
- Has an electron configuration that increases the electron-per-atom ratio (e/a ratio) → shifts the alloy away from the TCP stability field
- Does NOT significantly contribute to TCP formation itself

**4th generation alloys** (TMS-138, EPM-102): 3–6% Re + 2% Ru → better stability than 3rd generation while maintaining high creep resistance.

Ru cost: $14,000–20,000/kg → even more expensive than Re. A blade with 2% Ru and 6% Re has extremely high alloying element cost.

---

## 8. Grain Boundary Elements (PC Alloys Only)

In **polycrystalline and DS alloys only** (not SX):

### Boron (B) — 0.005–0.03 wt%

Segregates to grain boundaries. Reduces grain boundary energy. Dramatically improves high-temperature ductility and fatigue crack growth resistance. Essential in PC alloys.

In SX alloys: harmful. Can form B-rich eutectic pockets (lower melting point → incipient melting during solution treatment → disrupts microstructure).

### Zirconium (Zr) — 0.01–0.1 wt%

Similar to B: segregates to grain boundaries, improves ductility. Removed in SX alloys.

### Summary: SX vs. PC alloy chemistry differences

| Element | PC / DS alloys | SX alloys | Reason for difference |
|---------|---------------|-----------|----------------------|
| B | 0.01–0.03% | < 0.001% | No grain boundaries in SX |
| C | 0.05–0.15% | < 0.005% | No grain boundaries in SX |
| Zr | 0.01–0.1% | < 0.01% | No grain boundaries in SX |
| Hf | 0.1–1.5% | 0–0.1% | Grain boundary benefit absent |
| Al | 4–6% | 5–7% | Higher SX window for HT |
| Re | 0% | 3–6% | Enables full creep benefit |
| Ta | 3–5% | 5–9% | More freedom without GB elements |

---

## 9. Harmful Phases — TCP and What Causes Them

### Topologically Close-Packed (TCP) Phases

TCP phases are intermetallic compounds with complex crystal structures that are brittle and consume beneficial solid solution elements from the γ matrix.

**Common TCP phases in Ni superalloys:**
- **σ phase**: (Cr,Mo,Re,W)x(Ni,Co)y — Tetragonal. Forms plates or needles. Extremely brittle. Cr/Mo/Re-rich.
- **μ phase**: Mo₆Co₇ type — Rhombohedral. Similar embrittlement.
- **P phase**: (Cr,Mo,Ni,Co) complex — Less common. High Re alloys.
- **Laves phase**: AB₂ type — Some alloy systems.

**How TCP forms:**
After extended time at temperature (>800°C), heavy refractory elements (Re, W, Mo, Cr) can precipitate from supersaturated γ matrix into TCP phases. The driving force is thermodynamic stability — these compositions prefer the TCP crystal structure to the disordered FCC matrix.

**Consequences:**
- TCP plates/needles are crack initiation sites (stress concentrators)
- Re/W/Mo removed from γ matrix → loss of solid solution strengthening → lower creep resistance
- Blade life reduced by 30–50% if TCP forms significantly

**Prevention:**
- CALPHAD alloy design to stay below the TCP stability boundary
- Ru addition (4th generation alloys)
- Not exceeding a critical Re+Cr+Mo+W equivalent composition
- Operating temperature control

---

## 10. CMSX-4: A Complete Worked Example

CMSX-4 (Cannon-Muskegon Single Crystal alloy #4) is the benchmark 2nd-generation SX alloy, developed by Cannon-Muskegon Corporation in the mid-1980s and used in virtually every major jet engine program from 1990 to the present.

**Nominal composition (wt%):**

| Element | wt% | Role |
|---------|-----|------|
| Ni | Balance (~66%) | FCC matrix base |
| Cr | 6.5% | Oxidation resistance, hot corrosion |
| Co | 9.6% | γ stability, solid solution, anti-TCP |
| Al | 5.6% | γ′ former (Ni₃Al), oxidation resistance |
| Ti | 1.0% | γ′ former, APB energy |
| Ta | 6.5% | γ′ strengthener, diffusion barrier |
| W | 6.4% | Solid solution strength, matrix stiffener |
| Mo | 0.6% | Solid solution (small amount) |
| Re | 3.0% | Creep resistance (dislocation climb, coarsening suppression) |
| Hf | 0.1% | Oxidation scale adhesion |
| **C** | **< 0.005%** | Intentionally minimized (no GBs) |
| **B** | **< 0.001%** | Intentionally minimized (no GBs) |

**Phase fractions at service temperature (~980°C):**
- γ (matrix): ~35 vol%
- γ′ (precipitate): ~65 vol%

**Crystal structure at service:**
```
CMSX-4 microstructure (service condition, cuboidal γ′):

  ┌────────────────────────────────────────────────┐
  │  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐      │
  │  │  γ′  │  │  γ′  │  │  γ′  │  │  γ′  │      │
  │  │cube  │  │cube  │  │cube  │  │cube  │      │
  │  └──────┘  └──────┘  └──────┘  └──────┘      │
  │  ←γ channel→                                  │
  │  ┌──────┐  ┌──────┐  ┌──────┐  ┌──────┐      │
  │  │  γ′  │  │  γ′  │  │  γ′  │  │  γ′  │      │
  │  └──────┘  └──────┘  └──────┘  └──────┘      │
  └────────────────────────────────────────────────┘
  
  γ′ cube size: ~400–500 nm (after aging)
  γ channel width: ~50–100 nm
  γ′/γ volume ratio: ~65/35
```

**CMSX-4 key properties:**
- Density: 8.7 g/cm³
- γ′ solvus: ~1316°C
- Solidus: ~1340°C
- Solution treatment: 1290–1320°C (below solidus, above γ′ solvus)
- 1000h rupture strength at 982°C: 200 MPa
- Oxidation rate at 1100°C: < 3 mg/cm² per 500h

---

## 11. Alloy Design Principles — The CALPHAD Approach

Modern superalloy design cannot rely on trial and error — the 10+ element alloy space is too large for experimental exploration. **CALPHAD** (CALculation of PHAse Diagrams) is the primary computational tool.

**CALPHAD approach:**
1. Build a thermodynamic database of Gibbs free energy functions for each phase in each binary and ternary subsystem, from experimental measurements
2. Sum these energies for the full multicomponent alloy (Ni-Cr-Co-Al-Ti-Ta-W-Mo-Re-Ru system)
3. Minimize the total Gibbs free energy → predict equilibrium phases and compositions at each temperature
4. Predict: γ′ volume fraction, γ′ composition, TCP stability, solidus temperature, density

**CALPHAD outputs for alloy design:**
- Can predict whether a new composition is TCP-stable
- Can predict γ′ volume fraction as a function of temperature → guides aging heat treatment
- Density calculation guides weight budget
- Phase stability map in (Re, W) or (Re, Cr) space shows "safe" alloy window

**Current state:** 4th-6th generation alloys (TMS-138, CMSX-10K, TMS-238) were designed using CALPHAD + limited experimental validation. What previously required 10 years of experimental alloy iterations now takes 1–2 years of computation + targeted experiments.

---

## 12. Generational Comparison Table

| Generation | Example Alloy | Re (wt%) | Ru (wt%) | T_capability | Key Advance |
|------------|--------------|----------|----------|--------------|-------------|
| PC | IN100 | 0 | 0 | ~930°C | Polycrystalline casting |
| PC → DS | MAR-M200DS | 0 | 0 | ~960°C | Directional solidification |
| 1st SX | PWA 1480 | 0 | 0 | ~990°C | Single crystal, no C/B/Zr |
| 2nd SX | CMSX-4, PWA1484, René N5 | 3 | 0 | ~1020°C | Re addition |
| 3rd SX | CMSX-10, René N6 | 5–6 | 0 | ~1045°C | More Re; TCP issues |
| 4th SX | TMS-138, EPM-102 | 3–6 | 2 | ~1060°C | Ru for TCP stability |
| 5th SX | CMSX-10K, TMS-196 | 5–6 | 3–6 | ~1080°C | Optimized Ru/Re balance |
| 6th SX | TMS-238 | ~6 | ~5 | ~1100°C | Highest known |

Notes:
- T_capability = estimated 1000h rupture temperature at 150 MPa stress (metal temperature)
- Turbine inlet GAS temperature is ~300°C higher than metal temperature (due to cooling)
- All 2nd+ generation alloys are single crystal; DS remains important for lower-temperature components

---

## Summary

- **Ni superalloys have 10–12 elements**, each with a specific role — no element is redundant.
- **γ matrix strengtheners**: Cr (oxidation, hot corrosion), Co (stability, anti-TCP), W/Mo/Re (solid solution creep resistance)
- **γ′ formers**: Al (primary, also oxidation), Ti (APB energy), Ta (best refractory γ′ former — strong, slow diffuser)
- **Carbide elements** (C, Hf): needed in PC alloys for grain boundary strengthening; minimized in SX alloys
- **Rhenium**: transformative addition — slows diffusion, retards dislocation climb, suppresses γ′ coarsening → 3 wt% Re gives +30°C capability (2nd generation)
- **Ruthenium**: stabilizes against TCP phases caused by high Re+W+Mo content → enables 4th+ generation alloys
- **SX vs PC chemistry**: SX alloys remove C/B/Zr (no grain boundaries needed), enabling higher solidus → better heat treatment window
- **TCP phases**: harmful intermetallics (σ, μ) forming from Re/W/Mo/Cr supersaturation → need CALPHAD to design around

**Next chapter:** How the γ′ phase works — its crystal structure, the anomalous yield effect, and the mechanisms that make it such an effective high-temperature strengthener.

---

## Exercises

1. CMSX-4 has density 8.7 g/cm³ vs. steel at 7.85 g/cm³. Calculate the centrifugal stress at the blade root for each material if a blade is 80 mm long, rotates at 10,000 RPM, and the disk radius is 400 mm. Use σ = ρ × ω² × L × (r + L/2) where ω is angular velocity in rad/s, L is blade length, r is disk radius. Which material sustains lower centrifugal stress?

2. CMSX-4 has 3% Re at $3,500/kg and density 8.7 g/cm³. A finished HPT blade has mass ~0.5 kg. Calculate: (a) the mass of Re in one blade, (b) the Re content value per blade at $3,500/kg, (c) if an engine has 80 HPT blades, what is the total Re value per engine?

3. An alloy designer wants to reduce TCP risk in a 6% Re alloy by substituting 2% Ru. The alloy originally has (wt%): Ni-9Co-6.5Cr-5.5Al-0.5Ti-9Ta-5W-6Re (balance Ni). After adding 2% Ru, what element(s) would you reduce to maintain total alloy weight fraction = 100%? What tradeoffs does each reduction create?

4. Explain why Ta is preferred over Ti as the γ′-strengthening element in modern high-temperature SX alloys (designed for 1000°C+ service), even though Ti also increases APB energy. Consider both the strengthening mechanism AND the coarsening resistance argument.

5. A new alloy has the following partition coefficients: Cr k=1.1, Re k=1.3, Ta k=0.8, Al k=0.9, Co k=1.0. In the as-cast condition with nominal composition 6.5Cr-3Re-6.5Ta-5.6Al-9.6Co, predict whether the dendrite cores or interdendritic regions will be enriched in: Re, Ta, Al. After standard solution heat treatment, these gradients are homogenized. Why is this homogenization more difficult in a high-Re alloy?
