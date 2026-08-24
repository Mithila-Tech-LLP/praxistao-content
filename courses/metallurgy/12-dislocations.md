# Chapter 12: Dislocations — The Defects That Make Metals Malleable

> **"The theoretical shear strength of a perfect iron crystal is about 20,000 MPa. The actual yield strength of annealed iron is 130 MPa — 150 times lower. This astonishing difference is explained by one type of defect: the dislocation. A dislocation allows plastic deformation to occur row-by-row rather than all-at-once — the same reason a heavy carpet can be moved by pushing a wrinkle across it."**

---

## Table of Contents

1. [The Gap Between Theory and Reality](#1-the-gap-between-theory-and-reality)
2. [What Is a Dislocation?](#2-what-is-a-dislocation)
3. [Edge Dislocations](#3-edge-dislocations)
4. [Screw Dislocations](#4-screw-dislocations)
5. [Mixed Dislocations and the Burgers Vector](#5-mixed-dislocations-and-the-burgers-vector)
6. [Dislocation Motion — Glide and Climb](#6-dislocation-motion--glide-and-climb)
7. [Slip Systems and Critical Resolved Shear Stress](#7-slip-systems-and-critical-resolved-shear-stress)
8. [Dislocation Density](#8-dislocation-density)
9. [Dislocation Interactions and Work Hardening](#9-dislocation-interactions-and-work-hardening)
10. [Partial Dislocations and Stacking Faults](#10-partial-dislocations-and-stacking-faults)
11. [Dislocations and Creep](#11-dislocations-and-creep)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. The Gap Between Theory and Reality

**Theoretical shear strength (perfect crystal):**
A perfect crystal resists shear by requiring ALL atomic bonds in a plane to break simultaneously:
```
τ_theoretical = G / (2π)  ≈  G/6   (Frenkel 1926)

For iron: G = 80 GPa → τ_theory ≈ 13 GPa (about 0.1G)
```

**Actual measured yield strength:**
- Annealed iron whiskers (very pure, no dislocations): ~7 GPa — close to theory!
- Normal annealed iron: 130–200 MPa — 50× lower than theory
- Cold-worked iron: 500–700 MPa

**Explanation:** Real metals contain dislocations. Dislocations move under low applied stress by breaking only ONE ROW of bonds at a time — like moving a carpet wrinkle. The theoretical calculation assumes breaking an ENTIRE plane simultaneously.

---

## 2. What Is a Dislocation?

A **dislocation** is a one-dimensional crystal defect — a linear boundary between a region that has slipped and a region that has not yet slipped.

Two types:
- **Edge dislocation:** Extra half-plane of atoms inserted in the crystal
- **Screw dislocation:** Crystal sheared in a helical arrangement

Both types are defined by:
- **Line direction ξ:** The direction the dislocation line runs
- **Burgers vector b:** The magnitude and direction of the lattice distortion

---

## 3. Edge Dislocations

An edge dislocation can be visualized as an extra half-plane of atoms inserted into the crystal:

```
Edge dislocation (looking along dislocation line, into page):

Perfect crystal:          Crystal with edge dislocation:
                          
  ─○─○─○─○─○─            ─○─○─○─○─○─
  │ │ │ │ │ │            │ │ │ │ │ │
  ─○─○─○─○─○─            ─○─○─○─○─○─
  │ │ │ │ │ │            │ │ ┬ │ │ │
  ─○─○─○─○─○─            ─○─○─┤─○─○─   ← extra half-plane (T symbol = ⊥)
  │ │ │ │ │ │            │ │ │ │ │ │       dislocation line is ⊥ to page
  ─○─○─○─○─○─            ─○─○─○─○─○─
```

**Key feature of edge dislocation:**
- Burgers vector **b** is perpendicular to the dislocation line **ξ**
- Motion under shear stress: the extra half-plane moves in direction of **b**
- Can climb (move perpendicular to slip plane) by adding/removing vacancies — important for creep

**Stress field around edge dislocation:**
The lattice is compressed above the slip plane (where the extra half-plane is) and tensioned below. This stress field extends ~5–10 nm in all directions. Other dislocations and solute atoms interact with this stress field:
- Large solute atoms (size mismatch positive) → attracted to tensile region → drag on dislocation = solid solution strengthening
- Opposite-sign dislocations attract and annihilate (reduce dislocation density)
- Same-sign dislocations repel

---

## 4. Screw Dislocations

A screw dislocation can be visualized by cutting a crystal partway and shearing the top half by one lattice spacing:

```
Screw dislocation (axonometric view):

     ○───○───○
    /│  /│  /│
   / │ / │ / │
  ○───○───○  │
  │  ○───○───○
  │ /│  /│  /
  │/ │ / │ /
  ○───●───○    ← ● = core of screw dislocation (at center)
      │
   dislocation line runs VERTICALLY through the crystal
   
If you trace a circuit around the screw dislocation,
you advance ONE lattice spacing along the axis per circuit
→ helical path (like a screw thread, hence the name)
```

**Key feature of screw dislocation:**
- Burgers vector **b** is PARALLEL to the dislocation line **ξ**
- Can glide on any plane containing the dislocation line (not restricted to one slip plane)
- This makes cross-slip possible → screw dislocations can "switch" slip planes
- Cross-slip is important in work hardening and high-temperature deformation

---

## 5. Mixed Dislocations and the Burgers Vector

Real dislocations are usually **mixed** — partly edge, partly screw character. The **Burgers vector b** fully characterizes a dislocation:

### Burgers Circuit

To find **b**: draw a closed circuit in a PERFECT crystal around the same topology as the dislocation. The closure failure (gap vector) = **b**.

### Magnitude of b

For FCC metals: **b** = (a/2)⟨110⟩, |b| = a/√2
For BCC metals: **b** = (a/2)⟨111⟩, |b| = a√3/2
For HCP metals: **b** = (a/3)⟨11-20⟩

**Dislocation energy:** U = αGb²L (per unit length, L)
- Lower |b| → lower energy → preferred dislocation Burgers vector
- Frank's rule: a dislocation will only dissociate if the new dislocations have lower b² sum

### Frank's Rule
A reaction b₁ → b₂ + b₃ is favored if:
```
|b₁|² > |b₂|² + |b₃|²
```

---

## 6. Dislocation Motion — Glide and Climb

### Glide (Conservative Motion)

Dislocation moves on its slip plane — the plane containing both **b** and **ξ**. This conserves the number of atoms (no vacancies needed). Driven by applied shear stress. This is the mechanism for plastic deformation at low T.

**Peierls-Nabarro stress** — minimum stress to move a dislocation through a perfect lattice:
```
τ_PN = (2G)/(1-ν) × exp(-2πa/b(1-ν))

Very sensitive to b/a ratio:
- Large b → harder to move → stronger (one reason HCP metals like Be are harder to deform)
- Wide slip plane spacing (large a) → easier to move
```

### Climb (Non-Conservative Motion)

Edge dislocations can move perpendicular to their slip plane by ABSORBING or EMITTING vacancies. Climb is a diffusion-controlled process — requires vacancy migration → only significant at elevated temperatures (T > 0.3 T_m).

```
Dislocation climb (edge dislocation moves UP):

Before: ─○─○─┴─○─○─    ← extra half-plane ends here
        ─○─○─────○─○─   ← slip plane

Vacancy migrates to dislocation core:

After:  ─○─○─┬─○─○─    ← extra half-plane now one step shorter
        ─○─○─┴─○─○─    ← new slip plane (dislocation climbed UP one step)
```

**Dislocation climb is critical for creep** (Ch 16): at 900–1,100°C in turbine blades, dislocations climb over obstacles (γ′ particles, other dislocations) at a rate limited by vacancy diffusion → power-law creep. Re in CMSX-4 retards climb → better creep.

---

## 7. Slip Systems and Critical Resolved Shear Stress

Dislocations glide on specific crystallographic planes in specific directions:

**Slip system** = {slip plane} + ⟨slip direction⟩

```
FCC slip systems (most important for metals):
  Slip plane: {111} (4 planes)
  Slip direction: ⟨110⟩ (3 per plane)
  Total: 4 × 3 = 12 slip systems ← most slip systems → most ductile

BCC slip systems:
  Slip plane: {110}, {112}, {123} (pencil glide)
  Slip direction: ⟨111⟩
  Total: up to 48 ← many, but harder to activate

HCP slip systems:
  Slip plane: (0001) basal + {10-10} prismatic + {10-11} pyramidal
  Only 3 independent slip systems easily activated → more brittle
```

**Schmid's Law and Critical Resolved Shear Stress (CRSS):**

A slip system activates when the resolved shear stress on it reaches τ_CRSS:
```
τ_resolved = σ × cos(φ) × cos(λ) = σ × m

where φ = angle between stress axis and slip plane normal
      λ = angle between stress axis and slip direction
      m = Schmid factor = cos(φ) × cos(λ)

Maximum m = 0.5 (when φ = λ = 45°)
Minimum m = 0 (φ = 0° or λ = 0°)
```

**For a turbine blade under [001] tensile loading (Schmid factor calculation):**
- [001] loading in FCC: Schmid factor for (111)[10-1] = 0.41
- [111] loading in FCC: Schmid factor maximum = 0.27 (multiple systems activate simultaneously → hardest orientation)
- [001] is the softest orientation for FCC → SX blades oriented [001] → lower tensile yield, but better creep ductility (why [001] is chosen: see Ch 51)

---

## 8. Dislocation Density

**Dislocation density ρ:** Total length of dislocation lines per unit volume (m/m³ = m⁻²):

| Material state | ρ (m⁻²) |
|---------------|----------|
| Well-annealed metals | 10⁸ – 10¹⁰ |
| Lightly cold worked | 10¹⁰ – 10¹¹ |
| Heavily cold worked | 10¹² – 10¹⁴ |
| SX turbine blade (service) | 10¹² – 10¹³ |
| Metallic whiskers (near-perfect) | < 10⁵ |

**Dislocation density and strength:**
Paradoxically:
- Very low ρ (near-perfect crystal): very high strength (whisker strength ≈ theoretical)
- Low to moderate ρ (annealed): LOW strength (dislocations move easily)
- High ρ (cold worked): HIGH strength (dislocations interfere with each other)

The minimum in strength occurs at intermediate dislocation density — annealed metals are "soft" because dislocations can move freely on clean slip planes.

---

## 9. Dislocation Interactions and Work Hardening

As deformation proceeds, dislocation density increases (from Frank-Read sources multiplying dislocations). As density increases, dislocations interact and block each other:

**Frank-Read Source:** A dislocation segment pinned at both ends bows out under stress, eventually wrapping around and generating a new dislocation loop — multiplication mechanism.

```
Frank-Read source (successive stages):

    ──A───────────────B──    (pinned at A and B, length L)
        ↓ shear stress
    ──A──○─────────○──B──    (bows out)
        ↓
    ──A─○─────────────○─B──   (bows further)
        ↓
    A────────────────────B    (wraps around, generating new loop)
       ╰─────────────╯       (new dislocation loop emitted)
```

**Taylor hardening:** The flow stress increases with dislocation density:
```
τ = τ₀ + αGb√ρ

where α ≈ 0.2–0.5 (interaction constant)
      G = shear modulus, b = Burgers vector
      ρ = dislocation density
```

**Work hardening rate:** σ = σ_0 + K×ε^n (power-law hardening), where n = strain hardening exponent (0.1–0.5 typically). Higher n → more work hardening → better formability (delays necking in tensile test).

**Jogs and Kinks:** Steps on a dislocation line:
- Kink: step in dislocation line within slip plane → mobile, aids glide
- Jog: step in dislocation line out of slip plane → impedes glide of screw dislocations → requires climb to move → strengthening at low T

---

## 10. Partial Dislocations and Stacking Faults

In FCC metals, a perfect dislocation (a/2)⟨110⟩ can dissociate into two partial dislocations:
```
Perfect:   b = (a/2)[1-10]   (energy ∝ b² = a²/2)

Dissociates to two Shockley partials:
  b₁ = (a/6)[2-1-1]  
  b₂ = (a/6)[1-2 1]
  
Check: (a/6)[2-1-1] + (a/6)[1-2 1] = (a/6)[3-3 0] = (a/2)[1-10] ✓
Energy: b₁² + b₂² = a²/3 + a²/3 = 2a²/3 < a²/2 ← lower energy → favored
```

Between the two partial dislocations: a **stacking fault** — a region where the FCC stacking sequence ABCABC... is locally disrupted to ABCABABC... (like HCP stacking). 

**Stacking fault energy (SFE) γ_SF** determines the equilibrium separation of partial dislocations:
```
d = Gb₁b₂ / (2πγ_SF)
```

- Low SFE (Cu, austenitic SS, Co): wide stacking fault ribbon → dislocations separated → cross-slip difficult → more work hardening → important for creep resistance in superalloys
- High SFE (Al, Ni): narrow stacking fault → cross-slip easy → dislocations can rearrange → more efficient recovery → easier to hot work

---

## 11. Dislocations and Creep

At high temperature (Ch 16), dislocations overcome obstacles by thermal activation:

**Climb-controlled creep (power law creep):**
Dislocations pile up against γ′ obstacles in Ni superalloys. To continue deformation, they must climb over the γ′ particles. Climb requires vacancy diffusion. This controls the creep rate:
```
ε̇_creep = A × D_L × (σ/E)^n × exp(-Q/RT)

where D_L = lattice diffusion coefficient
      n = 3–8 (power law exponent)
      Q = activation energy for lattice diffusion
```

**How CMSX-4 resists creep:**
- Re partitions to γ matrix → reduces vacancy diffusion coefficient → slows climb
- γ′ particles force dislocations to climb long distances (γ′ rafts, Ch 44)
- High γ′ volume fraction (65%) → dislocations have almost nowhere to go except γ channels

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Why metals are weak | Dislocations allow row-by-row shear; σ_actual << σ_theoretical |
| Edge dislocation | Extra half-plane; b ⊥ ξ; moves parallel to b; climbs by vacancy exchange |
| Screw dislocation | b ∥ ξ; can cross-slip between planes; important for work hardening |
| Burgers vector | Closure failure of Burgers circuit; b = a/2⟨110⟩ for FCC |
| Slip systems | FCC: 12 (most ductile); HCP: 3 (most brittle) |
| CRSS / Schmid factor | Slip activates when τ_resolved = τ_CRSS; m = cos(φ)cos(λ) |
| Work hardening | τ = τ₀ + αGb√ρ; Frank-Read sources multiply dislocations |
| Partial dislocations | FCC dissociates into Shockley partials; stacking fault between them |
| Creep | Dislocation climb over γ′ obstacles; controlled by diffusion; Re retards this |

---

## Exercises

1. Iron is BCC with a = 2.87 Å. Calculate: (a) Burgers vector magnitude |b| for the (a/2)[111] dislocation, (b) the theoretical shear strength τ_th = G/(2π) using G = 80 GPa, (c) actual yield stress of annealed iron ≈ 130 MPa, (d) by what factor does the dislocation reduce the required stress? Express as a ratio.

2. A copper single crystal (FCC, a = 3.62 Å, G = 48 GPa) is loaded in tension along [123]. The primary slip system is (111)[10-1]. Calculate: (a) Schmid factor m = cos(φ)cos(λ) for this loading/slip system combination, (b) the tensile stress needed to initiate slip if τ_CRSS = 1.0 MPa for annealed Cu, (c) after deformation, the dislocation density is 10¹² m⁻². Using τ = τ₀ + αGb√ρ with α = 0.3, what shear stress is now required?

3. In FCC nickel (a = 3.52 Å), a perfect dislocation b = (a/2)[1-10] dissociates into Shockley partials. The stacking fault energy is γ_SF = 125 mJ/m² for Ni. Using d = Gb₁b₂/(2πγ_SF) (approximate) with G = 76 GPa: (a) calculate b₁ and b₂ magnitudes for (a/6)[2-1-1] and (a/6)[1-2 1], (b) estimate the equilibrium separation d between partial dislocations, (c) explain why Co additions to Ni reduce γ_SF and what effect this has on creep resistance.

4. A cold-worked steel has a dislocation density ρ = 5 × 10¹³ m⁻². Using the Taylor hardening law τ = αGb√ρ with α = 0.3, G = 80 GPa, b = 2.48 × 10⁻¹⁰ m: (a) calculate the shear flow stress contribution from dislocations, (b) estimate tensile yield strength using σ = 3τ (Taylor factor ≈ 3 for polycrystal), (c) this steel is then annealed at 700°C until ρ drops to 10¹⁰ m⁻². What is the new yield strength? (d) Calculate percentage strength loss. Is this consistent with the soft annealed state concept?

5. Explain in terms of dislocation physics why: (a) FCC metals are generally more ductile than HCP metals (consider slip system number and Schmid factor availability), (b) fine grain size increases strength but might decrease creep resistance at 900°C (hint: LAGB and HAGB effects on dislocation sources vs. grain boundary sliding), (c) adding Mo to nickel solid solution increases creep resistance more than adding Cr (consider atomic size mismatch and its effect on dislocation climb).
