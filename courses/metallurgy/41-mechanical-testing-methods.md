# Chapter 41: Mechanical Testing Methods — Measuring What Materials Can Do

> **"A tensile test tells you how a material behaves when you pull it slowly to death. A Charpy test tells you how it fails when hit suddenly at low temperature. A fatigue test tells you how it behaves over millions of repetitions. A creep test tells you what it does under constant stress at high temperature over years. Together these tests translate microstructure into engineering numbers — the bridge between what we make in the lab and what the designer puts in the stress analysis."**

---

## Table of Contents

1. [Tensile Testing — The Foundation](#1-tensile-testing--the-foundation)
2. [Hardness Testing](#2-hardness-testing)
3. [Impact Testing (Charpy, Izod)](#3-impact-testing-charpy-izod)
4. [Fatigue Testing — S-N and Strain-Life](#4-fatigue-testing--s-n-and-strain-life)
5. [Fracture Toughness Testing (K_Ic, J_Ic, CTOD)](#5-fracture-toughness-testing-k_ic-j_ic-ctod)
6. [Creep Testing](#6-creep-testing)
7. [High-Temperature Testing](#7-high-temperature-testing)
8. [Non-Destructive Testing (NDT) Overview](#8-non-destructive-testing-ndt-overview)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Tensile Testing — The Foundation

**The tensile test** is the most fundamental mechanical test: pull a standard specimen to fracture while measuring load and extension.

**Specimen geometry (ASTM E8):**
```
         ┌─────────────────────┐
Grip → ──┤  shoulder │  gauge  │ shoulder ├── Grip
         └─────────────────────┘
            gauge length L₀, diameter d₀

L₀ = 4d₀ (round) or 5.65√A₀ (flat)   (ASTM standard)
```

**Engineering stress-strain curve:**
```
σ = F / A₀   (engineering stress — uses ORIGINAL area)
ε = (L - L₀) / L₀   (engineering strain)
```

**True stress-strain:**
```
σ_true = F / A_instantaneous = σ_eng × (1 + ε_eng)
ε_true = ln(L/L₀) = ln(1 + ε_eng)
```

**Key measurements from tensile test:**

| Property | Symbol | Measurement method |
|---------|--------|--------------------|
| Young's modulus | E | Slope of linear elastic region (E = σ/ε) |
| 0.2% proof stress | σ_0.2 (or R_p0.2) | 0.2% offset from linear — used as yield strength |
| Ultimate tensile strength | UTS (σ_u) | Maximum stress on engineering curve |
| Fracture strain (ductility) | A% | (L_f - L₀)/L₀ × 100% |
| Reduction in area | RA% | (A₀ - A_f)/A₀ × 100% |
| Work hardening exponent | n | Slope of log-log σ vs ε in plastic range |
| Strength coefficient | K | σ = K × ε^n |

**Strain rate effect:**
Higher strain rate → higher apparent yield strength (rate-dependent plasticity). Standard tests done at ~10⁻⁴ s⁻¹ (quasi-static).

---

## 2. Hardness Testing

**Hardness** measures resistance to permanent indentation. Faster and cheaper than tensile testing; can be done on finished parts.

**Rockwell (HR, ASTM E18):**
- Applies minor load (98 N) + major load (varies by scale)
- Measures depth of indentation
- HRC (diamond, 150 kgf): for hard steels (25–70 HRC)
- HRB (1/16" ball, 100 kgf): for soft steels, Cu, Al (25–100 HRB)
- HRA: for hard coatings, cemented carbides

**Vickers (HV, ASTM E92):**
- Diamond pyramid indenter (136° included angle)
- HV = 0.1891 × F (kgf) / d² (mm²)
- Loads: 1–100 kgf (macrohardness); 0.001–1 kgf (microhardness)
- Consistent scale: soft to very hard
- Most versatile

**Brinell (HB, ASTM E10):**
- Steel ball (10 mm), 500–3,000 kgf load
- HB = 2F / (πD(D - √(D²-d²)))
- For soft to medium-hard materials; leaves large impression
- Used for: castings, forgings, raw stock

**Knoop (HK, ASTM E384):**
- Diamond elongated pyramid
- Very light loads (< 200g); elongated indentation
- Good for: very thin coatings, brittle materials, oriented crystals (elongated impression shows anisotropy)

**Conversion approximations (for steel):**
```
HV × 0.95 ≈ HB  (for HV < 300)
HRC ≈ (HV - 10) / 5.6  (approximate, valid 20–70 HRC)
UTS (MPa) ≈ 3.4 × HV  (very approximate; valid only for unnotched)
```

---

## 3. Impact Testing (Charpy, Izod)

**Impact tests** measure energy absorbed during rapid fracture — characterize toughness at high strain rate and in the presence of a notch.

**Charpy test (ASTM E23):**
```
         Pendulum
          \
           \  ← h₁ (initial height)
            \
             ┌──────┐
             │ notch│ ← specimen (10 × 10 × 55 mm with 2 mm V-notch)
             └──────┘
                  \
                   \  → h₂ (height after fracture)
Energy absorbed = mg(h₁ - h₂) (Joules)
```

**Charpy energy values:**
| Condition | Energy |
|-----------|--------|
| Brittle fracture (e.g., ferrite at -80°C) | 2–10 J |
| Ductile fracture (tough steel at 20°C) | 100–200 J |
| Transition region | 20–80 J |

**Ductile-to-Brittle Transition (DBTT):**
BCC metals show strong T-dependence:
- Above DBTT: high Charpy energy + fibrous (ductile) fracture surface
- Below DBTT: low Charpy energy + crystalline (cleavage) fracture surface
- DBTT criterion: 50% FATT (fracture appearance transition temperature) or specific J value (e.g., 27 J, 68 J)

**Effect of microstructure on DBTT:**
- Finer grain → lower DBTT (each grain boundary is an obstacle for cleavage crack)
- Higher C, P, N → higher DBTT
- Martensite + tempered: DBTT depends on tempering temperature
- Irradiation: raises DBTT (embrittlement in reactor pressure vessels → monitored by surveillance specimens)

**Izod test:** Similar principle but specimen is clamped vertically (less common; used in UK).

---

## 4. Fatigue Testing — S-N and Strain-Life

Covered extensively in Ch 14. Summary of testing methods:

**S-N (stress vs cycles):**
- Rotating beam (R.R. Moore): specimen rotates → fully reversed bending (R = -1)
- Servo-hydraulic frame: computer-controlled load; can apply any R-ratio
- Resonance testers: very high frequency (20 kHz) for > 10⁷ cycles (VHCF)

**Strain-life (ε-N):**
- Servo-hydraulic with extensometer: controls strain amplitude
- Get: Coffin-Manson coefficients (ε'_f, c) for plastic; Basquin (σ'_f, b) for elastic
- For LCF regime (turbine disks, < 10⁴ cycles)

**Fatigue crack growth (da/dN vs ΔK):**
- Compact tension (CT) specimen or center-crack panel (M(T))
- Pre-crack by fatigue → apply cycles → measure crack length (potential drop, crack gauge, compliance, visual)
- Plot da/dN vs ΔK → Paris Law constants C, m

**Key test standards:**
- ASTM E466: force-controlled axial fatigue
- ASTM E606: strain-controlled LCF
- ASTM E647: fatigue crack growth rate

---

## 5. Fracture Toughness Testing (K_Ic, J_Ic, CTOD)

**Linear Elastic Fracture Mechanics (K_Ic, ASTM E399):**

Specimen types: CT (compact tension), SEN (single edge notched), TPB (3-point bend)

```
Compact Tension (CT) specimen:
       ┌──────────────┐
   ○───┤  pre-crack   │
       └──────────────┘
       a (crack length)
       W (width)
       B (thickness)

K_I = F × f(a/W) / (B × √W)   where f(a/W) = geometry factor
```

**Validity requirements (plane strain):**
```
B, (W-a), a ≥ 2.5 × (K_Ic / σ_y)²

If not met → plane stress (apparent K > K_Ic) → test invalid, need thicker specimen
```

**K_Ic test procedure:**
1. Machine specimen, fatigue pre-crack to a/W ≈ 0.5
2. Load in tension → record load-displacement curve
3. Find P_Q from curve (5% secant method)
4. Calculate K_Q → verify validity criteria → K_Ic

**J-integral (J_Ic, ASTM E1820):**
For ductile materials that plastically deform before fracture (K_Ic test would need impossibly thick specimen):
- J measures energy per unit area of crack front
- J_Ic = (1-ν²)/E × K_Ic² (for small-scale yielding)
- Allows measurement on smaller specimens

**CTOD (Crack Tip Opening Displacement):**
- Used in pipeline, offshore structure, pressure vessel codes (BS7448, ISO 15653)
- δ = CTOD at onset of crack growth
- δ_c (critical) specifies material toughness for fitness-for-service assessment

---

## 6. Creep Testing

**Creep test:** Apply constant load/stress at constant temperature → measure strain as a function of time.

**Creep curve (strain vs time at constant σ and T):**
```
ε
│               Stage III: accelerating → fracture
│             /
│            /
│           /─────────────── Stage II: steady-state (linear)
│          /  
│        /
│────────  Stage I: primary (decreasing rate)
│ instantaneous elastic strain
└──────────────────────────────── t
```

**Three creep stages:**
- Stage I (primary): decreasing rate → work hardening dominates
- Stage II (secondary): constant minimum creep rate → steady state → work hardening = recovery
- Stage III (tertiary): accelerating → void coalescence, grain boundary sliding, necking

**Steady-state creep rate (Dorn equation):**
```
ε̇_ss = A × σⁿ × exp(-Q_c / RT)

where n = stress exponent (3–5 for dislocation creep; 1 for diffusion creep)
      Q_c = activation energy for dominant mechanism
      A = pre-exponential constant
```

**Creep mechanisms and deformation maps:**
At low σ/G, high T → Nabarro-Herring (bulk diffusion) or Coble (grain boundary diffusion) creep
At intermediate σ, high T → dislocation creep (n = 3–5)
At high σ, low T → dislocation glide (no T dependence)

**Stress rupture test:**
Similar to creep test but load applied until failure:
- Reports rupture time at given stress and temperature
- Larson-Miller parameter: LMP = T × (log t_r + C) where C ≈ 20 for Ni alloys
- Plot σ vs LMP → one master curve for all T/t combinations

---

## 7. High-Temperature Testing

**Challenges of testing at high temperature:**
- Oxidation of specimen surfaces → erroneous cross-section measurement
- Creep of grips, extensometer → need to account for
- Temperature uniformity across gauge length (< ±2°C specified in standards)
- Specimen-grips reaction at very high T

**High-temperature tensile test:**
- Induction or resistance furnace heating
- Alumina or high-purity Ni alloy grips
- Pyrometer or spot-welded thermocouple on specimen
- Temperature range: RT to 1,100°C for superalloy testing

**Creep testing at > 1,000°C:**
- Vacuum or controlled atmosphere (prevent oxidation)
- Dead-weight lever arm (constant load) or servo-controlled
- Laser extensometry (non-contact)

**Oxidation test:**
- Thermogravimetric analysis (TGA): sample on precision balance in furnace → record weight change vs T or time
- Cyclic oxidation: furnace + fan cooling; count cycles to specified weight loss
- Burner rig: combustion gas attacks sample at velocity → simulates engine environment

---

## 8. Non-Destructive Testing (NDT) Overview

NDT detects defects without damaging the part. Critical for aircraft engines (every blade inspected multiple times).

**Fluorescent Penetrant Inspection (FPI):**
- Low-viscosity fluorescent dye penetrates surface-breaking cracks (capillary action)
- Remove excess → developer draws dye out → UV light reveals cracks
- Sensitivity: surface cracks ≥ 25 μm wide
- Used for: all titanium and nickel alloy airfoils, 100% inspection

**Magnetic Particle Inspection (MPI):**
- Magnetize part → iron particles sprinkle on → field lines concentrate at cracks → particles pile up
- Only for ferromagnetic materials (steel, ferritic stainless)
- Detects: surface + subsurface cracks (≤ 3 mm deep)

**Ultrasonic Testing (UT):**
- Send ultrasonic pulse (0.5–25 MHz) → reflect from internal defects → time-of-flight = depth
- Sensitivity: ~ 0.5 mm defects (varies by frequency + geometry)
- Phased array UT (PAUT): computer-controlled beam steering → full 3D scan without moving probe
- Time-of-flight diffraction (TOFD): for accurate defect sizing (weld inspection)

**Radiographic Testing (RT):**
- X-ray or γ-ray through part → film or digital detector
- Dense material (inclusions, heavy segregation) → lighter on film; pores/cracks → darker
- Standard sensitivity: ~1% of material thickness
- Used for: casting inspection (turbine blade, structural castings), weld root pass

**Eddy Current Testing (ECT):**
- AC coil induces eddy currents in conductive material → cracks disrupt eddy current → impedance change in coil
- Surface + near-surface (< 2 mm)
- Fast, automated
- Used for: aircraft engine disk bores, fan blade surfaces

**Computed Tomography (CT):**
- X-ray CT: 2D X-ray images at multiple angles → reconstruct 3D volume
- Resolution: lab CT: 50–200 μm; synchrotron CT: < 1 μm
- Expensive, slow (hours per part) but gives complete 3D internal view
- Used for: turbine blade cooling channel inspection, PM part porosity mapping, ceramic TBC microstructure

---

## Summary

| Test | Standard | Measures | Key output |
|------|---------|---------|-----------|
| Tensile | ASTM E8 | σ-ε curve | E, σ_y, UTS, elongation, RA |
| Vickers hardness | ASTM E92 | Local hardness | HV; case depth |
| Charpy impact | ASTM E23 | Dynamic toughness | DBTT, absorbed energy |
| S-N fatigue | ASTM E466 | Stress vs cycles | Endurance limit, S-N curve |
| Strain-life fatigue | ASTM E606 | Strain vs cycles | ε'f, c, σ'f, b (LCF/HCF) |
| Crack growth | ASTM E647 | da/dN vs ΔK | Paris Law C, m |
| K_Ic | ASTM E399 | Fracture toughness | K_Ic (MPa√m) |
| Creep / rupture | ASTM E139 | Deformation at σ, T | ε̇_ss, t_rupture |
| FPI | ASTM E1417 | Surface cracks | Pass/fail, crack length |
| UT (PAUT) | ASTM E2700 | Internal defects | Defect size, depth |
| X-ray RT | ASTM E1742 | Porosity, inclusions | Accept/reject per code |

---

## Exercises

1. Tensile test of annealed 316L stainless (round bar, d₀ = 12.7 mm, L₀ = 50 mm): at fracture, F_max = 63.5 kN; F at 0.2% offset = 42.3 kN; L_f = 66.2 mm; d_f = 9.5 mm. (a) Calculate: σ_y (0.2% proof stress), UTS, elongation %, reduction in area %. (b) Convert to true stress and true strain at fracture (necking occurred — use d_f for true stress). Why is true stress at fracture higher than UTS? (c) Hollomon equation: fit σ_true = K × ε_true^n to data points at 5%, 10%, 20% uniform strain (read from curve approximation). Estimate n. (d) Calculate energy absorbed per unit volume up to fracture (area under engineering σ-ε curve = toughness). How does 316L compare to M2 tool steel (brittle, low area)?

2. Charpy impact testing of reactor pressure vessel steel (A508 Gr.2) over temperature range -100°C to +100°C gives data: -100°C: 8 J; -60°C: 22 J; -40°C: 48 J; -20°C: 105 J; 0°C: 168 J; +50°C: 185 J. (a) Plot the ductile-to-brittle transition curve. Determine the DBTT at: 27 J criterion and 50% FATT (fibrous fracture transition). (b) After 30 years of neutron irradiation (1 × 10¹⁹ n/cm²), the DBTT shifts +60°C (embrittlement). What new temperature corresponds to 27 J transition? (c) The design specification requires 27 J minimum at -20°C operating temperature. Does the irradiated vessel meet spec? What are the options? (d) Explain the mechanism of neutron irradiation embrittlement (hint: displacement damage, copper cluster precipitation).

3. Fracture toughness of 7075-T6 aluminum: CT specimen, W = 50 mm, B = 25 mm, a₀ (initial crack) = 25 mm. After fatigue pre-cracking, a = 26.5 mm. (a) Calculate a/W ratio and verify it meets the requirement 0.45 ≤ a/W ≤ 0.55. (b) The load-displacement curve gives P_Q = 18.5 kN (5% secant load). The geometry factor f(a/W = 0.53) = 9.59 from tables. Calculate K_Q in MPa√m. (c) Verify plane strain validity: B ≥ 2.5(K_Q/σ_y)² where σ_y = 503 MPa. Is B = 25 mm sufficient? (d) If not valid, what minimum thickness B_min is needed? (e) Compare to 7075-T73 (over-aged) with K_Ic = 32 MPa√m vs T6 (K_Ic = 24 MPa√m). What processing difference causes this improvement? (refer to Ch 24).

4. Creep test of CMSX-4 SX alloy at σ = 137 MPa, T = 1,050°C (a typical test condition for 2nd-gen SX alloys): (a) Results show Stage I (0–20 hours): rate decreasing from 5×10⁻⁵ /hr to 10⁻⁵ /hr; Stage II (20–500 hours): rate = 10⁻⁵ /hr; Stage III (>500 hours): rate accelerates → fracture at 800 hours. Calculate minimum creep rate in %/hour. (b) A Larson-Miller plot gives: at LMP = 24,000 for Ni SX alloys: σ_rupture = 137 MPa. LMP = T(K) × (log t_r + C) where C = 20. Verify: T = 1,050 + 273 = 1,323 K; t_r = 800 hours. Calculate LMP and compare to 24,000. (c) At T = 1,100°C, same σ = 137 MPa, using the same LMP = 24,000, calculate expected t_r. This shows the life penalty for +50°C. (d) Rafting microstructure forms after ~100 hours. Why does rafting (directionality of γ′ coarsening) initially reduce creep rate in Stage I, then eventually accelerate Stage III? (refer to Ch 48).

5. In-service turbine blade NDT program: after each engine overhaul (~25,000 cycles), all 80 HPT blades undergo: (a) FPI for surface cracks; (b) UT (phased array) for internal voids; (c) X-ray RT for cooling channel blockage; (d) Coordinate Measuring Machine (CMM) for airfoil geometry. (a) FPI sensitivity = 25 μm width. A crack of 50 μm × 2 mm is present in the blade. Will FPI detect it? What false-positive rate would you expect? (b) The UT system uses 15 MHz focused probe; sensitivity to voids ≥ 0.5 mm diameter. A pore of 0.3 mm is present. Will it be detected? What is the risk of a missed 0.3 mm pore in terms of fatigue life (use K_I = Fσ√πa with F=1, σ=300 MPa, K_Ic = 50 MPa√m; at what crack size does K_I reach K_Ic)? (c) If X-ray RT finds 15% blockage of a cooling channel: estimate the surface temperature increase for that section (use: Q = h×A×ΔT; if cooling flow drops 15%, ΔT increases 15%). Is this acceptable? (d) Propose a NDT detection probability (PD) curve (probability of detection vs defect size). What minimum detectable defect size would you specify for primary life-limiting defects in HPT blades?
