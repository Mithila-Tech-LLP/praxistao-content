# Chapter 11: Stress, Strain, and the Tensile Test

> **"The tensile test is the single most important mechanical test in all of engineering. In 10 minutes and with one small specimen, it tells you the stiffness, strength, ductility, and energy absorption of a material. Every design begins with these numbers."**

---

## Table of Contents

1. [Force, Stress, and Strain — Precise Definitions](#1-force-stress-and-strain--precise-definitions)
2. [Engineering Stress and Strain](#2-engineering-stress-and-strain)
3. [The Stress-Strain Curve — Reading the Story](#3-the-stress-strain-curve--reading-the-story)
4. [Elastic Region — Hooke's Law and Young's Modulus](#4-elastic-region--hookes-law-and-youngs-modulus)
5. [The Yield Point — When Plastic Deformation Begins](#5-the-yield-point--when-plastic-deformation-begins)
6. [Plastic Region — Strain Hardening and UTS](#6-plastic-region--strain-hardening-and-uts)
7. [True Stress and True Strain](#7-true-stress-and-true-strain)
8. [Ductility Measures](#8-ductility-measures)
9. [Hardness — The Quick Alternative](#9-hardness--the-quick-alternative)
10. [Resilience and Toughness](#10-resilience-and-toughness)
11. [Anisotropy and Texture Effects](#11-anisotropy-and-texture-effects)
12. [How Temperature Affects the Curve](#12-how-temperature-affects-the-curve)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. Force, Stress, and Strain — Precise Definitions

When you pull or push on a material, you apply a **force**. But force alone doesn't tell you whether the material will yield — a 1-tonne force will easily deform a toothpick but barely affect a steel rod. What matters is force per unit area.

**Stress (σ):** Force per unit area.
```
σ = F / A₀

Units: N/m² = Pa (Pascal)
       In practice: MPa (megapascals) = N/mm²
       1 MPa = 1 N/mm² = 145 psi
```

**Strain (ε):** Fractional change in length (dimensionless).
```
ε = (L - L₀) / L₀ = ΔL / L₀

Dimensionless. Sometimes expressed as % (multiply by 100).
```

**Types of stress:**
- **Tensile stress (+)**: pulls material apart (stretching)
- **Compressive stress (−)**: pushes material together (squeezing)
- **Shear stress (τ)**: one surface slides parallel to another

The **shear modulus** G relates shear stress to shear strain: τ = Gγ.

For isotropic materials: G = E / [2(1+ν)] where ν is Poisson's ratio.

---

## 2. Engineering Stress and Strain

The **standard tensile test** uses a **dog-bone** specimen:

```
    ┌──────────────────────────────────────────┐
    |Grips|  Shoulder  |  Gauge length (L₀)  |  Shoulder  |Grips|
    └──────────────────────────────────────────┘
                       ↑
               This is where deformation
               and fracture occur
```

Standard dimensions (ASTM E8):
- Gauge length L₀: 50 mm or 25 mm (must report which)
- Cross-sectional area A₀: calculated from diameter measurement

**Engineering stress** uses the **original** cross-section area A₀:
```
σ_eng = F / A₀
```

**Engineering strain** uses the **original** gauge length L₀:
```
ε_eng = (L - L₀) / L₀
```

Why "engineering"? Because in practice, you measure the applied force and original dimensions only. The actual area changes during deformation, but tracking that continuously was traditionally impractical.

---

## 3. The Stress-Strain Curve — Reading the Story

The tensile test produces a **stress-strain curve** that tells the complete story of material response:

```
Stress (MPa)
  │                    UTS (Ultimate Tensile Strength)
  │              ╭───────────●──╮
  │          ╭──╯    Strain       ╲   Fracture point
  │      ────●   hardening   ╲  ●
  │ Upper yield ──╯               ╲  ←  Necking region
  │ Lower yield ──────────         ╲
  │              Lüders             ╲
  │              band               ╲
  │                                  ●
  │         Elastic     Plastic
  │         region      region
  │
  └────────────────────────────────────→ Strain
       ↑           ↑           ↑
   ε = 0      Yield strain   Fracture strain
               (~0.002)         (e.g., 0.20)
```

---

## 4. Elastic Region — Hooke's Law and Young's Modulus

In the **elastic region**, stress and strain are linearly proportional. This is **Hooke's Law:**
```
σ = E × ε

E = Young's Modulus (modulus of elasticity) — the slope of the elastic line
```

**Young's Modulus** measures stiffness — resistance to elastic deformation:

| Material | E (GPa) | Notes |
|----------|---------|-------|
| Diamond | 1000 | Stiffest material known |
| Tungsten (W) | 411 | Stiffest common metal |
| Steel (Fe) | 200–210 | Design reference |
| Nickel | 207 | Close to steel |
| Titanium | 116 | High strength/weight, lower E |
| Aluminum | 69 | Light, flexible |
| Magnesium | 45 | Very light, quite flexible |
| Rubber | 0.001–0.01 | Very elastic (low E) |

**Physical meaning of E:** It reflects bond stiffness — how steeply the energy rises when atoms are displaced from their equilibrium separation. Strong, short bonds (covalent, transition metals with d-d bonding) → high E.

E is a material **constant** at a given temperature — you cannot change it by heat treatment, cold working, or alloying (for the same base metal). It only changes with:
1. Temperature (E decreases ~20% from 25°C to 700°C for Ni)
2. Crystal orientation (elastic anisotropy — important for single crystals)

**Elastic deformation is reversible.** Remove the load → the material returns to its original dimensions. This is because you're just stretching atomic bonds, not permanently rearranging atoms.

**Poisson's Ratio (ν):**
When you stretch a material longitudinally, it contracts laterally:
```
ν = −(lateral strain) / (axial strain) = −ε_lateral / ε_axial

Typical metals: ν ≈ 0.25–0.35
Rubber: ν ≈ 0.50 (nearly incompressible)
```

---

## 5. The Yield Point — When Plastic Deformation Begins

At some critical stress, the behavior changes from elastic to **plastic** — atoms permanently rearrange (via dislocation motion), and the material does not spring back when load is removed.

### 5.1 0.2% Offset Yield Strength

For most metals, the transition from elastic to plastic is gradual (no sharp yield point). By convention, the **0.2% offset yield strength** (σ_y or σ_{0.2}) is defined:

Draw a line parallel to the elastic slope but starting at ε = 0.002 (0.2%). Where this line intersects the stress-strain curve is the yield strength.

```
Stress
  │              ╭────── curve
  │        ╭────╯
  │   ╭────╯
  │  ─●─── (offset line, same slope as elastic, starts at ε=0.002)
  │         ●← 0.2% offset yield point
  └──────────────────────→ Strain
     0   0.002
```

This gives a reproducible, standard measure of the stress at which "permanent deformation begins."

### 5.2 Upper and Lower Yield Point

Some materials (mild steel, certain Al alloys) show a **sharp yield point** with an upper and lower yield stress:

```
  Stress
  │   ●Upper yield point
  │╭──╯ Sudden drop!
  ││ Lower yield point ─────────────────── (flat plateau = Lüders bands propagating)
  │
  └──────────────────────────→ Strain
```

The upper yield point occurs because dislocations are initially **pinned** by Cottrell atmospheres of interstitial carbon. The stress must first unpin them → sudden drop → lower yield point at which free dislocations propagate.

This effect is important for mild steel forming — Lüders bands create surface marks on pressed steel panels ("stretcher strains"). Avoided by skin-passing (small pre-strain) or by alloying to reduce Cottrell atmospheres.

### Typical Yield Strengths:

| Material | σ_y (MPa) |
|----------|----------|
| Pure aluminum | 35 |
| 6061-T6 Al | 276 |
| Mild steel (Fe-0.2% C) | 220 |
| HSLA steel | 500–700 |
| Ti-6Al-4V | 880 |
| 4340 steel (quenched + tempered) | 1470 |
| CMSX-4 superalloy at 750°C | 900 |
| Maraging steel (300 grade) | 2000 |

---

## 6. Plastic Region — Strain Hardening and UTS

Once the yield point is passed, permanent (plastic) deformation begins. As deformation continues:
- More dislocations are created
- They tangle and impede each other (forest hardening)
- Stress **must increase** to continue deformation → **strain hardening** (work hardening)

The stress continues to rise until the **Ultimate Tensile Strength (UTS):**
```
UTS = Maximum engineering stress on the stress-strain curve = F_max / A₀
```

At UTS, a **neck** begins to form — the cross-section locally contracts. From this point, the **force** actually decreases (even though the true stress in the necked region continues to rise), because the area is decreasing faster than the material hardens.

**Strain hardening exponent (n):**
In the plastic region, the true stress-true strain relationship is often described by:
```
σ_true = K × ε_true^n

n = strain hardening exponent
K = strength coefficient
```

| Material | n |
|----------|---|
| Annealed copper | 0.54 |
| Annealed Al | 0.20 |
| Annealed low-C steel | 0.21 |
| Cold-worked copper | ~0.12 |
| Superalloy (aged) | ~0.05 |

High n → material hardens rapidly → better formability (harder to neck locally in sheet forming).

---

## 7. True Stress and True Strain

Engineering stress/strain uses original dimensions. **True stress and strain** use instantaneous dimensions:

```
σ_true = F / A_instantaneous = σ_eng × (1 + ε_eng)
ε_true = ln(L/L₀) = ln(1 + ε_eng)
```

These are more physically meaningful, especially at large strains.

**Volume conservation** in plastic deformation (metals don't change volume plastically):
```
A₀ × L₀ = A × L    →   A = A₀ × L₀ / L = A₀ / (1 + ε_eng)
```

The true stress-strain curve continues to rise all the way to fracture, while the engineering curve drops after UTS due to necking. In research and simulation, true stress-strain data is used. In engineering specifications, engineering values are reported.

---

## 8. Ductility Measures

Two standard measures:

**Percent Elongation (%EL):**
```
%EL = (L_f - L₀) / L₀ × 100%
```
L_f = gauge length after fracture (pieces pushed back together and measured)

**Percent Reduction in Area (%RA):**
```
%RA = (A₀ - A_f) / A₀ × 100%
```
A_f = cross-sectional area at the fracture surface

| Material | %EL | %RA |
|----------|-----|-----|
| Pure aluminum | 40 | 90 |
| 6061-T6 Al | 12 | 50 |
| Low-C steel (annealed) | 37 | 70 |
| 4340 steel Q+T | 12 | 50 |
| Ti-6Al-4V | 14 | 36 |
| CMSX-4 at RT | ~2 | — |
| Grey cast iron | ~0 | 0 |

Grey cast iron fractures without measurable plastic deformation — **brittle** material. Copper fractures after 40% elongation — **ductile**. Knowing the difference is critical for design (brittle materials must not be used in tension; ductile materials give warning before fracture).

---

## 9. Hardness — The Quick Alternative

The tensile test destroys the specimen. **Hardness testing** is non-destructive (or nearly so) and much faster.

**Rockwell Hardness (HRC, HRB):**
- Measures depth of indentation under a fixed load
- HRC: 150 kg, diamond cone → hard steels, superalloys
- HRB: 100 kg, steel ball → softer materials (Al, Cu, soft steels)
- Fast, simple, most common in production

**Vickers Hardness (HV or HVN):**
- Diamond pyramid indenter, small load (1–30 kg)
- Measure diagonal of square indent under microscope
- `HV = 1854 × F / d²` (F in grams, d in μm)
- More precise; works at microstructural scale (micro-Vickers)
- Used for superalloy characterization

**Brinell Hardness (HB):**
- Steel ball indenter, 3000 kg load
- Larger indent → averages more microstructure
- `HB = 2F / (πD(D - √(D²-d²)))` where D = ball diameter, d = indent diameter

**Hardness-Tensile Strength correlation (for steel):**
```
UTS (MPa) ≈ 3.45 × HB     (approximate; works well for carbon steels)
UTS (MPa) ≈ 3.3 × HV
```

This correlation allows rapid estimation of strength from hardness. In production control of turbine blades, Vickers hardness mapping across the blade cross-section checks for adequate heat treatment uniformity.

---

## 10. Resilience and Toughness

**Resilience** = energy stored elastically per unit volume (area under elastic portion of curve):
```
Modulus of Resilience = σ_y² / (2E)   (J/m³)
```
High resilience → good springs (springs must deflect elastically, not plastically).

**Toughness** = total energy absorbed per unit volume before fracture = area under the entire stress-strain curve.

```
Toughness ≈ σ_y × ε_f × correction factor ≈ (σ_y + UTS)/2 × ε_f
```

**High toughness requires BOTH high strength AND high ductility** — which is why the "strength-toughness trade-off" is the central design challenge in structural metallurgy.

| Material | UTS (MPa) | %EL | Toughness |
|----------|----------|-----|-----------|
| Mild steel | 400 | 37 | High |
| 4340 Q+T | 1600 | 12 | Moderate |
| Maraging 300 | 2000 | 10 | Moderate |
| Grey cast iron | 150 | 0 | Very low |
| 7075-T6 Al | 570 | 11 | Moderate |

Grey cast iron has high compression strength but near-zero toughness — catastrophic brittle fracture under tension or impact. That's why it's used in compressive applications (engine blocks) but not in tension-critical structures.

---

## 11. Anisotropy and Texture Effects

For **polycrystalline metals** with random grain orientations, mechanical properties are approximately isotropic (same in all directions). But:

**After rolling or forging**, grains align in preferred orientations (**texture**). Properties differ with direction:
- Rolling direction (RD): often highest strength, lowest ductility
- Transverse direction (TD): intermediate
- Through-thickness (S): lowest ductility (short transverse grain boundary faces)

This is why aircraft structural specifications quote separate properties for L, LT, and ST directions.

**For single-crystal superalloys**, anisotropy is profound:
- E[001] ≈ 125 GPa (most compliant — preferred for blade primary axis)
- E[111] ≈ 294 GPa (stiffest)
- This 2.4× variation in stiffness with direction is critical for blade design and life

A single-crystal blade oriented near [001] experiences lower thermal stresses for the same temperature gradient, extending creep life significantly.

---

## 12. How Temperature Affects the Curve

**Effects of increasing temperature:**
- **E decreases**: atomic bond softens at higher vibration amplitudes (E of Ni drops from 207 GPa at 20°C to ~140 GPa at 900°C)
- **σ_y decreases**: thermal energy activates dislocation motion at lower stress
- **%EL increases**: more ductile at high temperature (easier atom movement)
- **Strain rate sensitivity increases**: creep becomes significant (deformation depends on how fast you apply load)

```
Stress
  │  Room temperature:    ╭────● UTS
  │                  ╭────╯
  │──────────────────╯
  │                     
  │  High temperature:  ╭──●
  │                ╭───╯       (lower, softer, more ductile)
  │───────────────╯
  │
  └────────────────────────────→ Strain
```

For turbine blade materials, the key performance region is **750–1050°C** — where conventional yield and tensile properties are less important than **creep resistance**. The tensile test at high temperature gives useful data but the creep test (Chapter 16 and Chapter 41) is the primary design basis.

---

## Summary

- **Stress** = F/A₀ (MPa); **strain** = ΔL/L₀ (dimensionless). Elastic → proportional; plastic → permanent.
- **Hooke's Law**: σ = Eε in elastic region. **Young's modulus E** = material stiffness constant; cannot be changed by heat treatment.
- **Yield strength (σ_y)**: 0.2% offset convention; when permanent deformation begins; dislocations start to move.
- **UTS**: maximum engineering stress; onset of necking; materials fail here in tension.
- **%EL and %RA**: ductility measures. High ductility → safe (warning before fracture).
- **Strain hardening**: σ rises in plastic region as more dislocations form and tangle.
- **True stress**: uses instantaneous area; rises all the way to fracture.
- **Hardness** (Rockwell, Vickers, Brinell): rapid strength proxy; HV×3.3 ≈ UTS for steel.
- **Toughness** = area under curve = needs both strength AND ductility simultaneously.
- **Temperature**: E and σ_y decrease; ductility increases. At high T, creep governs.

**Next chapter:** The tensile test shows that metals deform plastically at stresses far below theoretical. The mechanism is dislocation motion — and to engineer strength, we must understand how dislocations move, and how to stop them.

---

## Exercises

1. A Ti-6Al-4V specimen (A₀ = 50 mm², L₀ = 50 mm) fractures at F = 44 kN after elongating to L_f = 57 mm. Calculate: (a) UTS, (b) %EL, (c) engineering strain at fracture.

2. A steel spring must elastically store 0.5 J/m³ of energy without yielding. If E = 200 GPa, what minimum yield strength is needed? (Use modulus of resilience = σ_y²/2E)

3. A 4340 steel has HB = 450. Estimate its UTS using the approximation UTS ≈ 3.45 × HB. How does this compare to the tabulated value of ~1600 MPa?

4. Explain why a single-crystal Ni superalloy turbine blade is oriented with [001] parallel to the centrifugal force direction, given that E[001] = 125 GPa and E[111] = 294 GPa. How does lower E reduce thermal fatigue damage?

5. Compare the toughness (qualitatively) of: (a) pure copper (UTS=210 MPa, %EL=40%), (b) 4340 Q+T steel (UTS=1600 MPa, %EL=12%), (c) grey cast iron (UTS=150 MPa, %EL=0%). Which has the highest toughness? The lowest? What does this mean for design?
