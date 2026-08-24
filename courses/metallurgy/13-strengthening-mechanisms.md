# Chapter 13: Strengthening Mechanisms — How Metals Are Made Strong

> **"The yield strength of pure aluminum is 7 MPa — barely stronger than candle wax. The yield strength of 7075-T6 aluminum (Al-Zn-Mg-Cu alloy, peak aged) is 570 MPa — 80 times stronger. The same base metal; the difference is how we obstruct dislocation motion. Every strengthening mechanism in metallurgy is, at its core, a technique to make dislocations work harder."**

---

## Table of Contents

1. [The Unifying Principle](#1-the-unifying-principle)
2. [Grain Boundary Strengthening (Hall-Petch)](#2-grain-boundary-strengthening-hall-petch)
3. [Solid Solution Strengthening](#3-solid-solution-strengthening)
4. [Precipitation Hardening](#4-precipitation-hardening)
5. [Dispersion Strengthening](#5-dispersion-strengthening)
6. [Work Hardening (Strain Hardening)](#6-work-hardening-strain-hardening)
7. [Composite Strengthening — Fiber and Whisker Reinforcement](#7-composite-strengthening--fiber-and-whisker-reinforcement)
8. [Combining Mechanisms — Real Alloys Use All of Them](#8-combining-mechanisms--real-alloys-use-all-of-them)
9. [Strengthening vs. Temperature](#9-strengthening-vs-temperature)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. The Unifying Principle

Every strengthening mechanism in metallurgy works by creating **obstacles to dislocation motion**. The yield strength σ_y is the stress at which dislocations break through or bypass these obstacles:

```
Obstacle types and their mechanisms:

OBSTACLE                   MECHANISM TO RESIST
─────────────────────────────────────────────────────────────
Grain boundary             Slip must reinitiate in new grain
                           (Hall-Petch)
                           
Solute atoms               Stress field interaction with
(in solid solution)        dislocation strain field
                           (size + modulus mismatch)
                           
Second-phase precipitates  Coherency strain + Orowan
(coherent/semi-coherent)   bypassing (cutting vs. looping)
                           
Incoherent particles       Orowan bypassing
(dispersoids)              (dislocation loops around particle)
                           
Other dislocations         Intersection + jog formation +
(work hardening)           long-range stress fields
```

**Superposition of mechanisms:**
When multiple mechanisms are active, the contributions add (approximately):
```
σ_y = σ_0 + Δσ_HP + Δσ_SS + Δσ_P + Δσ_WH + ...

(sometimes added in quadrature: σ = √(σ₁² + σ₂² + ...))
```

---

## 2. Grain Boundary Strengthening (Hall-Petch)

Covered in Ch 09, but revisited here in the context of combined strengthening:

```
Δσ_HP = k_y / √d

Examples:
  Iron: k_y = 0.74 MPa√m
  Steel (0.3%C): k_y = 0.54 MPa√m
  Aluminum: k_y = 0.07 MPa√m (low — boundaries less effective in Al)
  Titanium: k_y = 0.40 MPa√m
```

**Practical limit:** Can't reduce grain size below ~1 μm by conventional processing. Below ~10–20 nm, the Hall-Petch relationship INVERTS (inverse Hall-Petch) because grain boundaries themselves become the weak zone, not barriers.

**Not effective at high temperature:** Grain boundaries become sources of weakness (grain boundary sliding, diffusional creep) above ~0.4 T_m → this is why single crystals are used for hot sections of turbines (Ch 47).

---

## 3. Solid Solution Strengthening

Solute atoms dissolved in the crystal lattice create local stress fields that impede dislocation motion.

### Size Mismatch (Paraffit Effect)

A solute atom larger or smaller than the solvent distorts the lattice:
- Large atom: compressive stress field around it
- Small atom: tensile stress field around it
- Edge dislocation has a tensile side (below extra half-plane) and a compressive side (above)
- Solute atoms migrate to minimize lattice strain energy → "atmosphere" around dislocation (Cottrell atmosphere)
- Dislocation must drag this atmosphere to move → higher stress needed

```
Hardening increment (per unit concentration):
Δτ ∝ G × |ε_s|^(3/2) × c^(1/2)

where ε_s = misfit strain = (r_solute - r_solvent)/r_solvent
      c = solute concentration
```

### Modulus Mismatch

If solute has different elastic modulus from solvent → local stiffness change → dislocation energy changes as it passes through → additional resistance.

### Combined Effect

```
Δσ_SS (approximate, per at% solute, for Ni base):

Mo: ~40 MPa/at%  (size mismatch ~12%)
W:  ~35 MPa/at%  (size mismatch ~10%)
Re: ~30 MPa/at%  (size mismatch ~9%)
Cr: ~10 MPa/at%  (size mismatch ~4%)
Co: ~5 MPa/at%   (size mismatch ~3%)
```

**Temperature independence:** Solid solution hardening is relatively temperature-insensitive (compared to precipitation hardening) because solute atoms are thermally stable at any temperature below the solidus. This makes it useful for high-temperature applications.

---

## 4. Precipitation Hardening

Fine precipitate particles embedded in the matrix provide very strong obstacles. This is the most powerful room-temperature strengthening mechanism and is the basis for all high-strength Al alloys, many steel grades, and nickel superalloys.

### Two Mechanisms: Cutting vs. Bypassing

The dominant mechanism depends on precipitate size and coherency:

**Particle Cutting (Shearing):** For small, coherent particles, dislocations CUT THROUGH the particle:
```
Mechanisms of particle cutting:
1. Coherency strain (lattice mismatch between particle and matrix)
2. Order hardening (dislocation must create an antiphase boundary in ordered γ′)
3. Stacking fault energy difference
4. Modulus mismatch

Δτ_cutting ∝ r × f^(1/2) × (particle size and volume fraction)
As particle size r INCREASES → Δτ_cutting INCREASES
```

**Orowan Bypassing:** For large, incoherent particles, dislocations LOOP AROUND the particle and leave a dislocation loop:
```
Orowan mechanism:
   
   ─────────●─────────   ← particle (●)
            ↓ stress
   ──────○     ○──────   ← dislocation bows around particle
         ↓
   ──────○──●──○──────   ← loops meet behind particle
         ↓
   ─────────●─────────   ← dislocation passes, loop left behind
         ↓ (loop tightens)
   ──────────────────    ← dislocation continues, loop remains around particle
   
Orowan stress: τ_Orowan = Gb / (2π × λ) where λ = interparticle spacing
```

**Critical Transition Radius r*:**

```
Cutting (r < r*):    Δτ increases with r   (dislocations prefer cutting smaller particles less)
Orowan (r > r*):     Δτ decreases with 1/λ (larger particles → wider spacing → easier bypass)

Peak strength at r = r* → this is the "peak-aged" condition!
```

```
Strength vs. aging time (precipitation hardening):

Strength
  │         peak aged
  │        ╱       ╲
  │       ╱         ╲  overaged (particles coarsened,
  │      ╱            ╲   larger spacing, Orowan bypassing)
  │     ╱               ─────
  │────╱  under-aged (small GP zones, cutting)
  │
  └──────────────── Aging time →
```

### Precipitation Hardening in Nickel Superalloys

The γ/γ′ system (Ch 44) exploits precipitation hardening with additional features:
- Ordered L1₂ structure → antiphase boundary (APB) energy → cutting requires superdislocation
- APB energy (γ_APB ≈ 180 mJ/m² for Ni₃Al) → large cutting resistance
- Orowan stress: τ = 0.84Gb/λ where λ is the γ channel width (~50 nm in CMSX-4)
- BOTH mechanisms active simultaneously

---

## 5. Dispersion Strengthening

Similar to precipitation hardening but with incoherent, thermally stable particles that do NOT dissolve on heating:

**Key difference from precipitation hardening:**
- Precipitates: thermodynamically stable only below solvus; coarsen quickly at high T
- Dispersoids: mechanically stable at much higher T; deliberately added as insoluble oxides/carbides

**Examples:**
- **TD-Nickel:** Thoria (ThO₂) dispersed in Ni → ThO₂ does not dissolve in liquid Ni → survives solidification → pins dislocations and grain boundaries at 1,200°C
- **ODS alloys:** MA956 (Fe-Cr-Al-Y₂O₃), MA6000 (Ni-Cr-Mo-Y₂O₃) → Y₂O₃ nano-dispersoids from mechanical alloying → creep resistance to 1,200°C (Ch 65)
- **Sintered aluminum powder (SAP):** Al₂O₃ particles in Al matrix

**Mechanism:** Pure Orowan bypassing (particles too large and incoherent to cut). The strength retention at high temperature is the key advantage:
```
Δσ_ODS ~ Gb / λ (same Orowan formula)

But λ stays small even at 1,100°C because dispersoids DON'T coarsen
(insoluble in matrix at any temperature)
```

---

## 6. Work Hardening (Strain Hardening)

As plastic deformation proceeds, dislocation density increases and dislocations interfere with each other:

```
σ = σ_0 + K × ε^n   (power-law work hardening)

where ε = plastic strain
      n = work hardening exponent (0 = no hardening; 1 = linear)
      K = strength coefficient

Typical values:
  Annealed copper: n ≈ 0.5, K ≈ 500 MPa
  Low-carbon steel: n ≈ 0.25, K ≈ 550 MPa
  Cold-worked steel: n ≈ 0.10, K ≈ 700 MPa
```

**Sources of work hardening:**
1. **Forest hardening:** Moving dislocations intersect "forest" dislocations → jog creation → energy cost
2. **Long-range back stress:** Piled-up dislocations create back stress opposing further dislocation motion
3. **Dislocation tangles and cells:** High-density regions form cells → dislocations must break through cell walls

**Practical use:**
- Cold-drawn wire (piano wire): extreme cold working → σ_y > 2,000 MPa
- Cold-rolled sheet: controlled work hardening for springback in forming
- Shot peening turbine blades: surface work hardening → compressive residual stress → improved fatigue life

**Removed by annealing:** High temperature restores soft annealed state through recovery (dislocation rearrangement) and recrystallization (new grain nucleation).

---

## 7. Composite Strengthening — Fiber and Whisker Reinforcement

**Not a dislocation mechanism** — but relevant for advanced systems:

In metal matrix composites (MMCs) and CMCs, a hard, strong reinforcement carries load alongside the matrix. Load transfer occurs through the interface:

```
Rule of mixtures (longitudinal loading):
σ_composite = V_f × σ_fiber + (1 - V_f) × σ_matrix

where V_f = volume fraction of reinforcement
```

For SiC/SiC composites (Ch 62):
- SiC fiber: σ_uts ≈ 3,000 MPa
- SiC matrix: σ ≈ 400 MPa
- Fiber bridging prevents crack propagation → K_Ic ≈ 15–25 MPa√m (much higher than monolithic SiC ≈ 3–5 MPa√m)

---

## 8. Combining Mechanisms — Real Alloys Use All of Them

Modern structural alloys typically use multiple mechanisms simultaneously:

### CMSX-4 Single Crystal Superalloy at 800°C:
| Mechanism | Contribution | Notes |
|-----------|-------------|-------|
| Solid solution (γ matrix) | ~300 MPa | Mo, W, Re, Co, Cr in FCC |
| γ′ precipitation | ~700 MPa | 65 vol%, coherent L1₂, APB cutting + Orowan |
| Dislocation density (sub-structure) | ~50 MPa | from prior solution treatment |
| **Total estimated** | **~1,050 MPa** | σ_y(800°C) ≈ 950 MPa measured |

### 7075-T6 Aluminum at 25°C:
| Mechanism | Contribution |
|-----------|-------------|
| Solid solution (Zn, Mg, Cu) | ~100 MPa |
| Grain boundary (d ≈ 10 μm) | ~50 MPa |
| MgZn₂ precipitation (peak aged) | ~420 MPa |
| **Total** | **~570 MPa** |

---

## 9. Strengthening vs. Temperature

All strengthening mechanisms weaken at high temperature, but at different rates:

```
Normalized strength vs. T/T_m:

Strength
  │
1.0│─── work hardening (→ anneals out above 0.3 T_m)
   │────── grain boundary (Hall-Petch) (→ GBS above 0.5 T_m)
   │─────────── solid solution (relatively stable to 0.8 T_m)
   │────────────────── γ′ precipitation (peaks at 0.55 T_m due to anomalous yield)
   │──────────────────────── ODS dispersoids (most stable, to 0.9 T_m)
   │
0  └────────────────── T/T_m ──────────────────────────────→
   0                   0.5                                  1.0
```

**For turbine blades at 0.70–0.78 T_m:**
- Work hardening: largely irrelevant (recovery fast)
- Hall-Petch: eliminated by using SX
- Solid solution: still active (~40%)
- γ′ precipitation: still active (~60%) — anomalous yield peak is beneficial
- ODS: would help but manufacturing is complex

---

## Summary

| Mechanism | Physical Basis | Best For |
|-----------|---------------|---------|
| Hall-Petch | Grain boundary blocks slip transmission | Room-T; HSLA steel, fine-grained alloys |
| Solid solution | Solute stress fields obstruct dislocations | All T; Mo/W/Re in Ni superalloys |
| Precipitation | APB energy + Orowan bypassing | RT to 0.6T_m; Al alloys, superalloys |
| Dispersion | Orowan bypass of insoluble particles | High T; ODS alloys (Y₂O₃) |
| Work hardening | Dislocation density → mutual interference | RT fabrication; springs, cables |
| σ_y peak | At r = r* (peak aged) | Optimize aging time/temperature |

---

## Exercises

1. An Al-4wt%Cu alloy (2024 aluminum, UNS A92024) starts with σ_y = 70 MPa (annealed). After peak aging: σ_y = 480 MPa. Estimate the contribution of precipitation hardening (assuming solid solution + grain boundary = 70 MPa). The θ′ (CuAl₂) precipitates have size 30 nm at peak age. If overaged to 200 nm diameter, Orowan stress ∝ 1/r: estimate new σ_y if all other contributions stay constant.

2. The Orowan stress for γ′ in a nickel superalloy with γ-channel width λ = 50 nm is τ_Orowan = Gb/(2πλ) × (1-ν) (simplified). Use G = 76 GPa, b = 2.5 × 10⁻¹⁰ m, ν = 0.3. Calculate: (a) τ_Orowan in MPa, (b) tensile yield stress contribution (σ = M × τ, M = 3 for polycrystal; for [001] SX, M ≈ 2.0), (c) after coarsening to λ = 200 nm (at 1,100°C after 1,000 hours), what is the new contribution? This illustrates why hot sections have a finite service life.

3. A steel alloy achieves yield strength through three mechanisms: (a) grain boundary: σ_HP = 0.74/√(15×10⁻⁶) MPa, (b) solid solution: 25 wt% Cr adding 5 MPa/wt%, and (c) work hardening: σ_WH = 300×ε^0.25 at ε = 0.15. Calculate σ_y assuming linear superposition. Compare to the rule that strength contributions add in quadrature and discuss which assumption is more conservative.

4. Explain the "inverse Hall-Petch" effect observed in nanocrystalline materials (grain size < 10 nm). Sketch strength vs. grain size from 1 mm to 1 nm on a log scale showing: (a) normal Hall-Petch regime, (b) transition around 10–20 nm, (c) inverse Hall-Petch regime. What physical mechanism causes the inversion, and why is this important for nano-crystalline coatings?

5. You are designing a nickel-base alloy for use at 1,050°C. Rank the following alloying additions by expected strengthening effectiveness at this temperature and explain your reasoning: (a) 0.5% C (interstitial), (b) 6% W (substitutional, large atom), (c) 0.5% B (grain boundary strengthening), (d) 5% Al (forms γ′), (e) 0.5% Y₂O₃ added via mechanical alloying. Consider both the mechanism and temperature stability.
