# Chapter 14: Fatigue — Failure Under Repeated Loading

> **"Over 80% of all mechanical failures in service are fatigue failures — and yet they often occur at stresses far below the static yield strength. A paperclip snapped 100 times or a crack growing invisibly through a turbine blade for 5,000 flights — both are fatigue. The tragedy is that fatigue cracks are invisible until they are nearly through the part. The art is detecting them before that point, and the science is preventing them from initiating at all."**

---

## Table of Contents

1. [What Is Fatigue?](#1-what-is-fatigue)
2. [The S-N Curve — Stress vs. Cycles to Failure](#2-the-s-n-curve--stress-vs-cycles-to-failure)
3. [Fatigue Limit and Endurance Limit](#3-fatigue-limit-and-endurance-limit)
4. [Fatigue Stages: Initiation, Propagation, Final Fracture](#4-fatigue-stages-initiation-propagation-final-fracture)
5. [Crack Initiation Mechanisms](#5-crack-initiation-mechanisms)
6. [Crack Propagation — Paris Law](#6-crack-propagation--paris-law)
7. [Fracture Mechanics Approach — ΔK and K_Ic](#7-fracture-mechanics-approach--δk-and-k_ic)
8. [Factors Affecting Fatigue Life](#8-factors-affecting-fatigue-life)
9. [Low-Cycle vs. High-Cycle Fatigue](#9-low-cycle-vs-high-cycle-fatigue)
10. [Fatigue in Turbine Blades](#10-fatigue-in-turbine-blades)
11. [Fatigue Testing and Detection](#11-fatigue-testing-and-detection)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. What Is Fatigue?

**Fatigue** is failure caused by repeated cyclic loading, at stresses typically below the static yield strength. Each cycle causes a tiny increment of damage — microplastic deformation or crack advance — until the cumulative damage causes fracture.

**Classic fatigue scenario:**
- A rotating shaft experiences tension on one side and compression on the other with every revolution
- After millions of revolutions, a small crack initiates at a stress concentration (notch, surface defect)
- The crack grows a tiny bit each cycle
- When the crack reaches a critical length, the remaining cross-section fails catastrophically

**Historical significance:** The Comet airliner disasters (1954) were caused by fatigue cracks growing from rectangular window corners. This tragedy established aviation's rigorous damage-tolerant design philosophy.

---

## 2. The S-N Curve — Stress vs. Cycles to Failure

The fundamental experimental characterization of fatigue is the **S-N curve** (Wöhler curve):

```
S-N Curve (semi-log plot):

Stress
amplitude
(MPa)
  600 │●
     │  ●
  400 │   ●●
     │      ●●
  200 │─────────●●●●●─────────── endurance limit (steel)
     │                          (no failure below this stress)
   0 │
    ─┴──────────────────────────────────────────── log(N)
     10³    10⁴    10⁵    10⁶    10⁷    10⁸
```

**Stress ratio R = σ_min / σ_max**
- R = -1: fully reversed (tension then equal compression)
- R = 0: zero-to-tension
- R > 0: tension-tension

**Mean stress effect (Goodman diagram):**
Tensile mean stress REDUCES fatigue life; compressive mean stress INCREASES fatigue life:
```
Goodman relation:
σ_a / σ_e + σ_m / σ_uts = 1  (for infinite life)

where σ_a = stress amplitude
      σ_e = endurance limit (R=-1)
      σ_m = mean stress
      σ_uts = ultimate tensile strength
```

---

## 3. Fatigue Limit and Endurance Limit

**Steel:** Has a true **fatigue limit** — below a certain stress amplitude, the S-N curve becomes horizontal and the steel will NEVER fail, no matter how many cycles:
```
Fatigue limit (steel) ≈ 0.45–0.5 × σ_UTS (for smooth specimens)
Typically: 200–500 MPa for structural steels
```

**Aluminum, titanium, nickel alloys:** Do NOT have a clear fatigue limit. The S-N curve continues to slope downward even at 10⁹ cycles. An "endurance strength" at a specific cycle count (10⁷ or 10⁸) is specified instead.

```
Material          Endurance at 10⁷ cycles / σ_UTS
────────────────────────────────────────────────
Low-carbon steel:  0.45–0.50
High-strength steel: 0.35–0.45
Aluminum alloys:   0.25–0.35
Titanium alloys:   0.45–0.55
Nickel superalloys: 0.40–0.50
```

---

## 4. Fatigue Stages: Initiation, Propagation, Final Fracture

Fatigue failure proceeds in three stages:

```
Fatigue life timeline:

[─────────── N_i ────────────][──── N_p ────][─]
  Stage I: Crack initiation    Stage II:      III: 
  (no visible crack)           Propagation   Final
                               (Paris Law)   fracture
  90–99% of fatigue life       ~1–9%         instant
  for smooth specimens
```

**Stage I (Initiation):** Microplastic deformation forms persistent slip bands (PSBs) at the surface. These create surface intrusions/extrusions (like accordion folds) that become crack nuclei. Very sensitive to surface finish and stress concentrations.

**Stage II (Propagation):** The crack propagates perpendicular to the maximum tensile stress. Characteristic **fatigue striations** form on the fracture surface — one striation per cycle. The crack advances a tiny distance (nanometers to micrometers) per cycle.

**Stage III (Final fracture):** When the crack reaches the critical length K_I = K_Ic, the remaining ligament fails instantly (overload fracture). Fracture surface shows fibrous zone (fatigue) and a coarser zone (final fracture).

---

## 5. Crack Initiation Mechanisms

### Persistent Slip Bands (PSBs)

In single-phase metals under cyclic loading:
1. Dislocations organize into wall-ladder structures (PSBs)
2. PSBs are softer than the surrounding matrix → concentrate plastic strain
3. Cyclic plastic strain pumps material OUT of the surface (extrusions) and IN (intrusions)
4. Intrusions are stress concentration points → crack nuclei after 10⁴–10⁵ cycles

### Notch Effect (Stress Concentration Factor K_t)

Holes, grooves, and notches concentrate stress:
```
K_t = σ_local / σ_nominal

For a circular hole: K_t = 3
For an elliptical hole, semi-axes a×b: K_t = 1 + 2a/b
For film cooling holes (turbine blades): K_t ≈ 3–4
```

Fatigue notch factor K_f (actual reduction in endurance strength) is slightly lower than K_t due to notch sensitivity:
```
K_f = 1 + q × (K_t - 1)    where q = notch sensitivity factor (0–1)
```

### Inclusions and Second-Phase Particles

Hard inclusions (MnS in steel, Al₂O₃ in superalloys) are incompatible in strain with the matrix → debond or crack on first few cycles → become crack nuclei.

---

## 6. Crack Propagation — Paris Law

Once a crack is initiated, it propagates according to the **Paris Law** (1963):

```
da/dN = C × (ΔK)^m

where: da/dN = crack growth rate (m/cycle)
       ΔK = stress intensity factor RANGE = K_max - K_min (MPa√m)
       C = material constant
       m = Paris exponent (typically 2–4 for metals)
```

```
Log(da/dN) vs Log(ΔK) diagram:

da/dN         │               /
(m/cycle)     │            / ← Region II (Paris Law)
              │         /    da/dN = C(ΔK)^m
10⁻⁶         │      /
              │   /  ← Region I (threshold, ΔK_th)
10⁻⁹         │ / (no crack growth below ΔK_th)
              │/
              └─────────────────── ΔK (MPa√m) →
                   ΔK_th         K_Ic/F
```

**Fatigue threshold ΔK_th:** Below this ΔK, cracks do not propagate. Typical values:
- High-strength steel: ΔK_th ≈ 4–8 MPa√m
- Titanium alloys: ΔK_th ≈ 3–5 MPa√m
- Aluminum alloys: ΔK_th ≈ 1–3 MPa√m

**Integrating Paris Law:** Total cycles from initial crack a_i to critical crack a_f:
```
N = ∫(from a_i to a_f) da / [C × (ΔK)^m]

For ΔK = ΔσF√(πa):

N = (a_f^(1-m/2) - a_i^(1-m/2)) / [(1-m/2) × C × (Δσ×F)^m × π^(m/2)]
```

---

## 7. Fracture Mechanics Approach — ΔK and K_Ic

(Full treatment of fracture mechanics is in Ch 15; key concepts applied to fatigue here)

**Stress intensity factor K_I:**
```
K_I = F × σ × √(πa)

where F = geometry factor (~1 for through crack in infinite plate)
      σ = applied stress
      a = crack half-length
```

**Damage-tolerant design:** Instead of trying to prevent all cracks (safe-life design), damage-tolerant design assumes cracks exist and determines the inspection interval:

```
Design logic:
1. Assume worst-case initial crack size a_i (NDE detection limit)
2. Calculate K as function of a
3. Integrate Paris Law to find cycles from a_i to a_f (where K → K_Ic)
4. Inspection interval = N_calculated / 2 (safety factor of 2)
```

**Used for:** Aircraft structures, pressure vessels, nuclear components, jet engine disks.

---

## 8. Factors Affecting Fatigue Life

### Surface Finish

Fatigue cracks almost always initiate at the surface. Surface condition is CRITICAL:

```
Surface roughness effect on fatigue limit (relative to polished surface):
  Mirror polished:        1.0  (reference)
  Fine machined:          0.9
  Turned/ground:          0.8
  Hot-rolled:             0.7
  Corroded (salt spray):  0.5–0.3
```

**Shot peening:** Blasting surface with steel shot → surface work hardening + compressive residual stress. The compressive stress reduces effective ΔK at surface cracks → 20–40% improvement in fatigue life. Used extensively on turbine compressor blades and connecting rods.

### Mean Stress

Tensile mean stress = harmful (adds to applied stress)
Compressive mean stress = beneficial (closes cracks, reduces ΔK)

### Environment

Corrosion fatigue: corrosive environment + cyclic stress → cracks initiate faster (pits as nuclei) and grow faster (active dissolution at crack tip). Stainless steel in seawater: 60–80% reduction in fatigue life.

### Frequency

At low frequency: more time for environment interaction → worse
At high frequency: adiabatic heating possible (polymers), but metals mostly unaffected

### Temperature

Elevated temperature: lower fatigue life (less matrix strength, oxidation, creep-fatigue interaction)
Low temperature: generally higher fatigue life (higher yield strength)

---

## 9. Low-Cycle vs. High-Cycle Fatigue

**Low-Cycle Fatigue (LCF):** N < 10⁴–10⁵ cycles; stresses exceed yield strength (local plasticity); characterized by strain amplitude ε_a:

**Coffin-Manson equation:**
```
ε_a^p = ε_f' × (2N)^c

where ε_a^p = plastic strain amplitude
      ε_f' = fatigue ductility coefficient (~true fracture strain)
      c = Coffin-Manson exponent (−0.5 to −0.7)
      2N = reversals to failure
```

**High-Cycle Fatigue (HCF):** N > 10⁵ cycles; stresses below yield strength; characterized by stress amplitude σ_a.

### Combined Basquin + Coffin-Manson (Strain-Life curve):
```
ε_a = σ_f'/E × (2N)^b + ε_f' × (2N)^c

where the first term = elastic strain (Basquin, HCF regime)
      second term = plastic strain (Coffin-Manson, LCF regime)
```

```
Strain-Life curve:

Total strain │   ε_a = elastic + plastic
amplitude    │  ╲  total = sum
             │   ╲ elastic (slope b)
             │    ╲
             │     ╲──────────────
             │         plastic (slope c, steeper in LCF)
             │
             └──────────────────────────── log(2N)
                    transition life N_t
```

---

## 10. Fatigue in Turbine Blades

Turbine blades experience BOTH LCF and HCF simultaneously:

**LCF (Low Cycle Fatigue):**
- 1 LCF cycle per flight (start → takeoff → cruise → landing)
- Thermal cycling: blade heats from RT to 950°C and back
- Each cycle generates ~ε = 0.1–0.5% thermal strain
- Design life: 30,000–50,000 flights
- Critical location: film cooling hole edges (K_t ≈ 3–4, as covered in Ch 58)

**HCF (High Cycle Fatigue):**
- Blade-passing frequency: 80 blades × 3,000 RPM = 4,000 Hz
- Each flight: 4,000 cycles/second × 3.5 hours = 50 million cycles per flight
- If a resonance is hit → 10⁹ cycles in minutes → sudden catastrophic failure
- Avoided by: Campbell diagram (resonance map), intentional damping, mistuning

```
Campbell diagram for HCF avoidance:

RPM
  │  
  │   / 6E (6 per revolution)
  │  / 
  │ / ── ── ──  ←  mode 1 natural frequency (horizontal line)
  │ / 
  └─────────────────── frequency (Hz)
        
Crossing = resonance. Design to avoid crossings in operating range.
```

**Material response:**
- Creep-fatigue interaction: at 950°C, each LCF cycle involves both cyclic strain AND creep relaxation
- Environmental fatigue: oxidation at crack tip → faster crack growth
- This is why turbine blade life is monitored to engine-cycle count, not just hours

---

## 11. Fatigue Testing and Detection

**Common fatigue tests:**
- **Rotating bending:** Classic R.R. Moore test; R = -1; good for S-N curves
- **Axial fatigue:** Servo-hydraulic test machine; any R ratio
- **Resonance fatigue:** Vibrate specimen at natural frequency (HCF, fast)
- **Thermal mechanical fatigue (TMF):** Simultaneous thermal cycling + mechanical load (turbine simulation)

**Non-destructive evaluation (NDE) for cracks:**
- Fluorescent penetrant inspection (FPI): cracks as small as 0.1 mm at surface
- Eddy current: surface and near-surface cracks; used for disk bores
- Ultrasonic testing: internal cracks; ≥0.5 mm typically
- X-ray radiography: internal porosity and large cracks
- **ACFM (Alternating Current Field Measurement):** Good for weld inspection, complex geometries

---

## Summary

| Concept | Key Point |
|---------|-----------|
| S-N curve | Stress amplitude vs. cycles to failure; tested at R = -1 typically |
| Fatigue limit | Steel has one (~0.5 σ_UTS); Al/Ni/Ti do NOT |
| Crack initiation | PSBs at surface; K_t at notches; inclusions |
| Paris Law | da/dN = C(ΔK)^m; m ≈ 2–4; integrates to give crack propagation life |
| LCF | N < 10⁵, plastic strain dominant, Coffin-Manson |
| HCF | N > 10⁶, elastic strain, endurance strength |
| Mean stress | Tensile = harmful; compressive = beneficial; shot peening exploits this |
| Turbine blades | 1 LCF cycle/flight (30,000 life); 10⁹ HCF cycles/flight (avoid resonance) |

---

## Exercises

1. A steel has σ_UTS = 900 MPa and fatigue limit ≈ 0.47 × σ_UTS. A component has a semicircular surface notch with K_t = 2.5 (notch sensitivity q = 0.85). Calculate: (a) K_f, (b) notch-corrected fatigue limit, (c) if mean stress σ_m = 150 MPa, use the Goodman relation to find the maximum allowable stress amplitude for infinite life.

2. A turbine disk material follows Paris Law: da/dN = 3 × 10⁻¹¹ × (ΔK)^3 (SI units, m and MPa√m). Initial crack a_i = 0.5 mm (NDE limit), critical crack a_f = 8 mm (where K_max → K_Ic = 60 MPa√m at σ_max = 600 MPa). Calculate: (a) K_max at a_i and a_f (K = σ√(πa)), (b) cycles to propagate from a_i to a_f using integrated Paris Law, (c) inspection interval (N/2 safety factor).

3. A rotating shaft (R = -1) has S-N curve following Basquin's equation: σ_a = σ_f'×(2N)^b with σ_f' = 1,200 MPa, b = -0.10. Calculate: (a) stress amplitude for N = 10⁵ cycles, (b) N for σ_a = 350 MPa, (c) stress amplitude for N = 10⁷ cycles. Is there an endurance limit in this equation? (d) If mean stress σ_m = 200 MPa and σ_UTS = 1,000 MPa, apply Goodman correction to the 10⁷ cycle endurance.

4. Shot peening introduces compressive residual stress of −400 MPa at the surface of a Ti-6Al-4V blade. The applied cyclic stress is σ_max = +600 MPa, σ_min = −100 MPa. Calculate: (a) original R ratio without shot peening, (b) effective R ratio after shot peening (add compressive residual stress to both max and min), (c) use Goodman diagram to explain qualitatively why this extends fatigue life.

5. A turbine blade film cooling hole has K_t = 3.5 and stress concentration at the hole edge. The LCF stress cycle is σ_max = 350 MPa, σ_min = 30 MPa. Material: σ_UTS = 1,150 MPa, σ_y = 1,000 MPa. (a) Calculate K_f assuming q = 0.9. (b) Check whether local stress at hole exceeds σ_y (local stress = K_t × σ_max). If so, the Goodman relation breaks down — why? (c) Suggest two design modifications to reduce the stress concentration at film cooling holes.
