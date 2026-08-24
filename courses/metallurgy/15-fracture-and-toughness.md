# Chapter 15: Fracture and Toughness — When Metals Break

> **"A material can be incredibly strong — able to withstand enormous stress — yet still fail catastrophically in a split second. The Liberty Ship disasters of WWII, where T2 tankers split in half in icy harbors at stresses below their design limit, taught us this lesson in blood and steel. Strength is not the same as toughness. A glass rod is stronger than a copper wire of the same diameter — but the copper wire bends; the glass shatters."**

---

## Table of Contents

1. [Ductile vs. Brittle Fracture](#1-ductile-vs-brittle-fracture)
2. [Energy-Based View — Griffith Theory](#2-energy-based-view--griffith-theory)
3. [Linear Elastic Fracture Mechanics — K_I and K_Ic](#3-linear-elastic-fracture-mechanics--k_i-and-k_ic)
4. [Crack Opening Modes](#4-crack-opening-modes)
5. [Plastic Zone and Plane Strain vs. Plane Stress](#5-plastic-zone-and-plane-strain-vs-plane-stress)
6. [Toughness Testing — Charpy and K_Ic](#6-toughness-testing--charpy-and-k_ic)
7. [Ductile-to-Brittle Transition Temperature (DBTT)](#7-ductile-to-brittle-transition-temperature-dbtt)
8. [Fracture Mechanisms at the Micro Scale](#8-fracture-mechanisms-at-the-micro-scale)
9. [Fracture in Superalloys and Turbine Applications](#9-fracture-in-superalloys-and-turbine-applications)
10. [Design for Toughness](#10-design-for-toughness)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Ductile vs. Brittle Fracture

**Ductile fracture:** Large plastic deformation before fracture. Slow crack growth. Requires significant energy. Warning before failure.

**Brittle fracture:** Little or no plastic deformation. Rapid crack propagation. Low energy. No warning before failure.

```
Comparison of fracture surfaces:

Ductile fracture                 Brittle fracture
("cup-and-cone"):                (flat surface):
    
    ──────                           ────────
   ╱       ╲   ← cup               │        │ ← flat face
  │  ●●●●●  │    (internal void     │        │   with river marks
  │ ●●●●●●● │    coalescence)       │        │
   ╲       ╱   ← cone             ────────
    ──────
    
Fracture surface: rough,          Fracture surface: smooth,
fibrous, gray                     granular, shiny (intergranular)
                                  or faceted (transgranular cleavage)
```

**Macro indicators:**
| Feature | Ductile | Brittle |
|---------|---------|---------|
| Neck (reduction in area) | 30–80% RA | < 5% RA |
| Fracture surface | Gray, fibrous | Crystalline, flat, shiny |
| Energy absorbed | High (Charpy > 50 J) | Low (Charpy < 10 J) |
| Warning | Visible deformation | None |
| Temperature | Favored by high T | Favored by low T, high strain rate |

---

## 2. Energy-Based View — Griffith Theory

Alan Arnold Griffith (1921) showed why cracks propagate: a crack grows when the energy released (elastic strain energy) exceeds the energy consumed (new surface area creation):

```
Griffith criterion for crack extension:

Energy balance for a through crack of length 2a in infinite plate:

U_elastic = -σ² × π × a² × B / E  (energy RELEASED as crack opens)
U_surface = 4 × a × B × γ_s       (energy REQUIRED to create new surfaces)

(B = thickness, γ_s = surface energy)

For crack growth: dU_total/da < 0 (total energy must decrease)
  d/da [-σ²πa²B/E + 4aBγ_s] = 0

→ σ_f = √(2Eγ_s / πa)   ← Griffith equation for fracture stress
```

**Key insight:** Fracture stress σ_f ∝ 1/√a → larger cracks → lower fracture stress. A defect that is twice as long needs only 70% as much stress to propagate. This is why defect detection is so important.

**Griffith modification for metals (Orowan, Irwin):** Real metals do plastic work (G_p >> γ_s) before fracture:
```
σ_f = √(E × G_c / πa)

where G_c = critical strain energy release rate (J/m²)
For metals: G_c = 10–1,000 kJ/m² (vs γ_s = 1–10 J/m² for surface energy only)
```

---

## 3. Linear Elastic Fracture Mechanics — K_I and K_Ic

The stress field near a crack tip (Irwin 1957):
```
σ_ij = K_I / √(2πr) × f_ij(θ) + higher order terms

where r = distance from crack tip
      θ = angle from crack plane
      K_I = stress intensity factor
      f_ij(θ) = dimensionless angular functions
```

**K_I** completely describes the crack tip stress field (in LEFM):
```
K_I = F × σ × √(πa)

where F = geometry factor (dimensionless, order 1)
      σ = applied stress
      a = crack length
      
Units: MPa√m
```

**Critical stress intensity factor (fracture toughness) K_Ic:**
Material property — the value of K_I at which the crack propagates unstably:
```
K_Ic = F × σ_fracture × √(πa_critical)
```

**Fracture condition:** If K_I ≥ K_Ic → fracture occurs.

**Design equation:**
For a known K_Ic and a known (or inspected) crack size a:
```
σ_max = K_Ic / (F × √πa)   ← allowable stress

Or: a_max = (K_Ic / (F × σ))² / π   ← maximum tolerable crack size
```

---

## 4. Crack Opening Modes

Three fundamental crack opening modes:
```
Mode I (Opening):     Mode II (In-plane shear):   Mode III (Anti-plane):
     ↑ ↑ ↑                   → → →                    ↑ (front face)
  ───────────              ─────────────           ─────────────
  ───────────              ─────────────           ─────────────
     ↓ ↓ ↓                   ← ← ←                    ↓ (back face)
  
  Tensile opening          Sliding mode             Tearing mode
  Most critical for        (shear parallel          (shear out of
  fracture                 to crack)                plane)
```

**Mode I is most dangerous** — tensile stress opens crack → KI greatest concern. K_Ic is always Mode I.

Real cracks often have mixed-mode loading, but Mode I dominates for most engineering failures.

---

## 5. Plastic Zone and Plane Strain vs. Plane Stress

Near the crack tip, stresses exceed yield strength → a small plastic zone forms:

**Plastic zone radius:**
```
r_y = (1/2π) × (K_I / σ_y)²   (Irwin estimate)
```

**Plane stress (thin plate):** Free to deform in thickness direction → large plastic zone → high toughness
**Plane strain (thick plate):** Constrained in thickness direction → small plastic zone → LOW toughness

```
Toughness vs. thickness:

K_c
 │─────────── (plane stress, varies with B)
 │
 │           ─────────── K_Ic (plane strain, constant)
 │
 └──────────────────── thickness B →
               B_critical = 2.5 × (K_Ic/σ_y)²
```

**Plane strain condition (ASTM E399):** Thickness B and crack size a must satisfy:
```
B, a > 2.5 × (K_Ic / σ_y)²

If not: use plane stress K_c (higher, not conservative for thick parts)
```

This is why thick structures are MORE susceptible to brittle fracture than thin ones — a counterintuitive result!

---

## 6. Toughness Testing — Charpy and K_Ic

### Charpy Impact Test

A simple, widely used test (but not a fracture mechanics parameter):

```
Charpy test setup:

         ○ pivot
          \
           \  ← pendulum weight
            \
             \
         ──────────────
        |    Specimen   |  ← notched bar in V or U
         ──────────────
              ↓
        Energy absorbed = mgh_initial - mgh_final
```

**Charpy V-notch (CVN) energy** in Joules. Good for screening and transition temperature determination, but NOT directly applicable to design (not a geometry-independent property).

### K_Ic Testing (ASTM E399)

Pre-cracked compact tension (CT) or single-edge notched bend (SENB) specimens. Fatigue pre-cracked to sharp crack → loaded → K at fracture = K_Ic. Requires validity criteria (plane strain condition).

```
Typical K_Ic values:

Material                K_Ic (MPa√m)
───────────────────────────────────
Low-alloy steel         50–100
High-strength steel     50–120
Stainless steel (304)   200+
Titanium alloys         50–100
Aluminum alloys         25–45
Nickel superalloys      40–80
Ceramics                2–8
High-entropy alloy      100–200+ (cryogenic K_Ic)
```

---

## 7. Ductile-to-Brittle Transition Temperature (DBTT)

Many BCC metals (ferritic steel, Cr, Mo, W) undergo a dramatic change from ductile to brittle behavior at a transition temperature:

```
Charpy energy vs. temperature (steel):

Charpy
energy (J)
  200 │─────────────────────────────────  ← upper shelf (ductile)
     │                          /
  100│                        /
     │                      /  ← transition region
   50│                    /
     │                  /
   10│────────────────/              ← lower shelf (brittle)
     │
     └──────────────────────────── Temperature (°C)
          -100       DBTT       +100
```

**DBTT depends strongly on:**
- Carbon content: higher C → higher DBTT (brittler)
- Grain size: finer grain → lower DBTT (tougher) — Hall-Petch also helps toughness
- Impurities: P, S, Sb → higher DBTT
- Strain rate: faster loading → higher DBTT
- Irradiation: neutrons → raise DBTT (critical issue for nuclear pressure vessels)

**Engineering requirement (WWII Liberty ship lesson):**
Design DBTT to be well below the lowest expected service temperature. Steel plates used in T2 tankers had DBTT ≈ +5°C; ships operated in North Atlantic at −5°C → brittle fracture.

**Modern structural steels:** DBTT = −40 to −80°C through controlled grain refinement (TMCP), low sulfur, clean steel practice.

---

## 8. Fracture Mechanisms at the Micro Scale

### Ductile Fracture by Void Coalescence (Microvoid coalescence, MVC)

1. Void nucleation at inclusions or second-phase particles (debonding or particle cracking)
2. Void growth under triaxial stress
3. Void coalescence by necking of ligaments between voids
4. Dimple fracture surface (SEM shows spherical dimples)

**Sensitive to:** Inclusion density, cleanliness of steel → ultra-clean steel (ladle treatment) → better toughness.

### Brittle Transgranular Cleavage

Crack propagates along specific crystallographic planes (cleavage planes: {100} in BCC iron):

1. Dislocation pile-up creates local stress concentration
2. Stress reaches critical value for cleavage
3. Crack runs along cleavage plane at near-sound speed

**SEM appearance:** Flat facets with "river marks" pointing back to crack origin.

### Intergranular Fracture

Crack follows grain boundaries (weakened by:)
- Segregation of P, S, Sn (temper embrittlement)
- Grain boundary oxidation (high temperature)
- Hydrogen embrittlement (hydrogen concentrates at boundaries)
- Liquid metal embrittlement (e.g., liquid Bi on Cu, liquid Hg on Al)

**SEM appearance:** Faceted grain shapes visible; smooth surfaces.

---

## 9. Fracture in Superalloys and Turbine Applications

**Nickel superalloy toughness:**
- Good but not exceptional: K_Ic = 40–80 MPa√m (room temperature)
- Decreases slightly at high temperature (recovery of work hardening → larger plastic zone is actually beneficial, but phase instability → embrittlement at some temperatures)
- Single crystal blades: no grain boundary fracture paths → better high-temperature fracture behavior
- Sub-grain boundaries (LAGB < 5°) in SX: still acceptable toughness

**Critical failure scenarios:**
- LCF crack from film hole → propagates through blade thickness → blade tip loss
- FOD impact → immediate fracture if K_I at impact crack > K_Ic → immediate failure or fast-growing crack
- TBC spallation → increased metal temperature → accelerated oxidation → thin section → final fracture

**Fracture of TBC system:**
- YSZ toughness: K_Ic ≈ 1.5–2.0 MPa√m (brittle ceramic)
- TGO: K_Ic ≈ 3–5 MPa√m
- TBC failure by interface delamination: mode I + mode II mixed; G_c ≈ 30–80 J/m²
- Thermal cycling causes mode-mixity changes → delamination crack driven at peak temperature by compressive thermal stress → spallation

---

## 10. Design for Toughness

### Material Selection

- High-strength steel: lower K_Ic; careful design required
- Lower strength + better K_Ic (trade-off)
- Temper embrittlement: avoid 350–550°C range in service; control P, Sn, Sb content

### Geometry Design

- Avoid stress concentrations (smooth radii, no sharp corners)
- Use damage-tolerant design (accept that cracks will exist; inspect and manage)
- Design such that a_critical > NDE detection limit

### Residual Stress

- Compressive residual stress at surface → closes surface cracks → higher effective K_Ic
- Shot peening, autofrettage (overpressurize pressure vessel → compressive bore residual stress)

### Microstructure Control

- Cleaner steels (lower inclusions): better toughness
- Finer grain size: higher K_Ic + lower DBTT
- Avoid temper embrittlement: controlled heat treatment

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Griffith | σ_f = √(EG_c/πa); bigger crack → lower fracture stress |
| K_I | Crack tip stress intensity; K_I = Fσ√(πa); fracture when K_I ≥ K_Ic |
| K_Ic | Material property (plane strain fracture toughness); 40–100 MPa√m for metals |
| Plane strain vs stress | Thick parts → plane strain → lower K_c; counterintuitive |
| DBTT | BCC metals: brittle below DBTT; fine grain size lowers DBTT |
| MVC | Ductile fracture by void nucleation → growth → coalescence |
| Cleavage | Brittle fracture; {100} planes; river marks on fracture surface |
| Superalloys | K_Ic = 40–80 MPa√m; SX avoids GB fracture; TBC spallation drives failure |

---

## Exercises

1. An aircraft aluminum alloy (K_Ic = 33 MPa√m, σ_y = 460 MPa) has a detected edge crack of length a = 4 mm. (a) Calculate K_I for σ = 200 MPa (F = 1.12 for edge crack). (b) Is this safe? (c) What is the critical crack length a_c at σ = 200 MPa? (d) What is the minimum thickness B needed for valid plane strain condition?

2. A Charpy test on a structural steel shows CVN = 10 J at 0°C and 150 J at 100°C. DBTT ≈ 40°C (midpoint). The ship will operate in water at 2°C in winter. (a) Is this steel suitable? (b) Three microstructural modifications to lower DBTT: grain refinement, reduced C content, reduced P content. For each, explain the metallurgical mechanism.

3. Using the Griffith equation σ_f = √(2Eγ_s/πa), estimate the fracture stress for a glass rod (E = 70 GPa, γ_s = 0.5 J/m²) with a 0.1 μm crack. Compare to the theoretical cleavage strength ~E/10. What does this tell you about the effect of flaws on glass strength? Why is glass reinforcement with fibers (glass fiber composites) effective?

4. An SX Ni superalloy turbine blade has K_Ic = 60 MPa√m. The LCF stress cycle gives σ_max = 350 MPa. Using a_critical = (K_Ic/(Fσ))²/π with F = 1.0: (a) calculate a_critical. (b) NDE (fluorescent penetrant) can detect cracks ≥ 0.1 mm. Is a_critical larger than the detection limit? (c) If K_Ic drops to 40 MPa√m at 1,000°C, recalculate a_critical. What does this mean for hot section NDE requirements?

5. Explain why the DBTT of low-carbon steel can be moved from +10°C to −60°C by: (a) reducing grain size from 50 μm to 10 μm (discuss how grain boundaries arrest cleavage crack propagation), (b) reducing carbon from 0.3% to 0.1% (discuss how carbon affects dislocation mobility at low temperature and promotes pearlite which provides easy cleavage paths), (c) adding 0.1% Nb (controlled rolling → fine ferrite grain; NbC precipitates pin γ grain boundaries). Estimate the total shift in DBTT from all three changes.
