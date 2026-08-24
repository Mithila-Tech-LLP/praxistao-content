# Chapter 55: Turbine Blade Loading — Forces, Stresses, and the Duty Cycle

> **"A turbine blade in a modern jet engine experiences simultaneously: centrifugal stress of 130 MPa from rotation at 10,000 rpm; thermal stress of 50–100 MPa from a 200°C temperature gradient through its wall; bending stress of 30–50 MPa from gas pressure; vibration stresses that must be kept below 20 MPa or the blade will fail by high-cycle fatigue in hours. Each of these stress sources is manageable individually. The challenge — and the art — is managing them simultaneously over 20,000 hours of engine operation."**

---

## Table of Contents

1. [Overview — Why Loading Analysis Matters](#1-overview--why-loading-analysis-matters)
2. [Centrifugal Stress](#2-centrifugal-stress)
3. [Gas Bending Load](#3-gas-bending-load)
4. [Thermal Gradient Stress](#4-thermal-gradient-stress)
5. [Low-Cycle Fatigue (LCF) — Engine Duty Cycle](#5-low-cycle-fatigue-lcf--engine-duty-cycle)
6. [High-Cycle Fatigue (HCF) — Vibration](#6-high-cycle-fatigue-hcf--vibration)
7. [Creep Under Combined Loading](#7-creep-under-combined-loading)
8. [Thermomechanical Fatigue (TMF)](#8-thermomechanical-fatigue-tmf)
9. [Foreign Object Damage (FOD)](#9-foreign-object-damage-fod)
10. [Combined Life Prediction — Interaction of Mechanisms](#10-combined-life-prediction--interaction-of-mechanisms)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Overview — Why Loading Analysis Matters

The turbine blade operates in one of the harshest mechanical environments in engineering. The blade designer must simultaneously satisfy:
- Structural integrity (no failure)
- Aerodynamic performance (correct airfoil geometry — even when hot)
- Thermal management (keep metal below maximum temperature)
- Weight minimization (reduces centrifugal load on disk)
- Manufacturing feasibility
- Service life (LCF, HCF, creep)

**The loading environment summary:**

| Load type | Magnitude | Frequency | Failure mode |
|-----------|---------|-----------|-------------|
| Centrifugal | 100–150 MPa (root) | Steady | Creep, LCF (net section fracture) |
| Gas bending | 30–60 MPa | Slowly varying | LCF, combined creep |
| Thermal | 50–150 MPa (gradient) | Per flight cycle | TMF, LCF (thermal fatigue) |
| Vibration | 10–50 MPa | HCF (10³–10⁵ Hz) | HCF (resonance) |
| Impact/FOD | Impulse | Single events | Dent, crack |

---

## 2. Centrifugal Stress

**Centrifugal stress is the primary mechanical load on a turbine blade.**

**For a simple uniform bar rotating about its root:**
```
σ_c(r) = ρ × ω² × ∫_r^R r dr = ρ × ω² × (R² - r²) / 2

where:
   ρ = blade density (kg/m³)
   ω = angular velocity (rad/s)
   r = radial position (m)
   R = blade tip radius (m)
```

**For a tapered blade (cross-section area A(r)):**
```
σ_c(r) = ρ × ω² × ∫_r^R A(r') r' dr' / A(r)
```

The centrifugal stress is MAXIMUM at the blade root (r = 0) and ZERO at the tip.

**Numerical example:**
- N = 10,000 rpm → ω = 2π × 10,000/60 = 1,047 rad/s
- CMSX-4 density ρ = 8,700 kg/m³
- Blade span = 0.12 m, mean radius R_mean = 0.35 m
- Blade mass = 250g = 0.250 kg, root area A_root = 400 mm² = 4×10⁻⁴ m²

```
Centrifugal force at root = m × ω² × R_mean = 0.250 × 1047² × 0.35 = 96,000 N

σ_centrifugal = F / A_root = 96,000 / 4×10⁻⁴ = 240 MPa
```

This is a high steady stress — the blade root must carry the centrifugal load of the entire blade over 20,000 hours at 950°C → creep is the primary concern.

**Effect of alloy density on centrifugal stress:**
Gen 5 alloys are ~7% denser than Gen 1 (9.2 vs 8.6 g/cm³). For the same geometry:
- Centrifugal stress increases proportionally
- This is why blade taper ratio (tip area / root area) is maximized → tip carries less weight

---

## 3. Gas Bending Load

**The gas flow exerts pressure on the airfoil surfaces** → net force on the blade → bending moment about the blade root.

**Aerodynamic force:**
- Lift force (perpendicular to airfoil) → bends blade in chord direction
- Drag force (parallel to flow direction) → small, usually negligible vs lift

**Bending stress at root:**
```
σ_bending = M × c / I

where:
   M = bending moment = F_aero × blade_span/2
   c = distance from neutral axis to airfoil surface
   I = second moment of area of root cross-section
```

**Typical values:**
- Gas bending force: 200–400 N per blade
- Bending stress at root: 30–60 MPa
- This is MUCH LESS than centrifugal stress but oscillates if there is unsteady aerodynamics

**Unsteady gas forces and resonance:**
The gas flow is not steady — the blade passes upstream stator vanes once per revolution → WAKE PASSING causes a periodic force on each blade:
```
Excitation frequency f_exc = N_stators × N_rotor (rpm) / 60

Example: 46 stator vanes, 10,000 rpm:
f_exc = 46 × 10,000 / 60 = 7,667 Hz
```

If this matches a blade natural frequency → resonance → HCF (see §6).

---

## 4. Thermal Gradient Stress

**Thermal gradients cause stress** because different parts of the blade expand by different amounts:

**Sources of thermal gradient:**
1. **Leading edge hotter than trailing edge:** Leading edge has first contact with hot gas → higher convective heat transfer coefficient
2. **Tip hotter than root:** Root is cooler (attached to cooler disk); tip is exposed to gas
3. **Outer wall hotter than inner wall:** Cooling air is inside; hot gas outside → temperature difference across 1–2 mm wall
4. **Cooling holes create local gradients:** High local cooling at hole exit → cold spots next to hot spots

**Thermal stress (simplified for constrained plate):**
```
σ_thermal = E × α × ΔT / (1 - ν)

For [001] CMSX-4: E = 130 GPa, α = 14×10⁻⁶/°C, ν ≈ 0.3
At ΔT = 100°C across wall:
σ_thermal = 130×10⁹ × 14×10⁻⁶ × 100 / 0.7 = 260 MPa
```

However, at high temperature, the material creeps under this stress → stress relaxes. The ACTUAL thermal stress is much less than this elastic estimate. The real thermal stress in a well-designed blade is 50–100 MPa.

**Key source of TMF damage:**
The thermal stress cycles with engine power changes:
- Take-off: maximum gas temperature → maximum thermal gradient → maximum thermal stress
- Cruise: steady state → thermal stress partially relaxed by creep
- Descent/landing: cooling → thermal contraction → reversed stress → net tensile
- Each flight cycle: one reversal → thermomechanical fatigue (TMF) crack initiation

---

## 5. Low-Cycle Fatigue (LCF) — Engine Duty Cycle

**LCF** occurs at the same frequency as engine power changes (< 10⁴ cycles total, high Δσ or Δε per cycle):

**Engine duty cycle:**
```
Max T₁ (take-off): σ_c + σ_bending + σ_thermal_max → total stress peak
         ↓
T₂ (cruise, steady): σ_c + σ_bending + σ_thermal_steady (lower)
         ↓
T₃ (descent): σ_c unchanged; σ_thermal reverses; cooling → net lower
         ↓
Ground idle / shutdown: σ_c → 0; cool → residual thermal
```

**LCF cycle counting:**
Each flight = one LCF cycle. Engine certified for N_LCF = 20,000–40,000 flights.

**Critical locations for LCF:**
1. Film cooling hole edges: stress concentration K_t ≈ 3–4; maximum stress ≈ K_t × σ_nominal
2. Platform fillets: geometric stress concentration
3. Root attachment firtree: highest stress amplitude (centrifugal load × stress concentration)

**Coffin-Manson for LCF (from Ch 14):**
```
Δε_p/2 = ε'_f × (2N)^c

For IN718 at 650°C: ε'_f = 0.46, c = -0.59
For CMSX-4 [001] at 850°C: different parameters but same form
```

**Life prediction:**
LCF life = f(Δσ, Δε, T, mean stress, environment) → complex models (e.g., Walker, SWT)

---

## 6. High-Cycle Fatigue (HCF) — Vibration

**HCF occurs at blade natural frequencies** — resonance driven by aerodynamic excitation:
- 10⁷–10⁹ cycles at relatively low Δσ (well below yield)
- Alternating stress of 20–50 MPa can cause failure if sustained

**Campbell diagram:**
The most important tool for HCF design:

```
Campbell Diagram:

Frequency
(Hz)
│         / 1st bending mode
│        /─────── RESONANCE CROSSINGS (circles!)
│       /
│      /   / 2nd bending
│     /   /
│    / X /  ← engine operating range (70–100% N₁)
│   /   /
│  /   /  
└──────────────────────────── Engine speed (rpm)

Horizontal lines: blade natural frequencies (constant vs. speed for uncoupled modes)
Slanted lines: excitation frequency = n × (rpm/60) × number of struts/vanes
Crossings within operating range: potential resonances → must clear or add damping
```

**Engine order (EO) excitation:**
- 1 EO = once per revolution (unbalance, one-off imperfection)
- N EO = N times per revolution (N = number of upstream vanes/stators, or N = number of downstream vanes)

**Natural frequencies of a turbine blade:**
1st bending: ~1,000–2,000 Hz (chordwise bending)
1st torsion: ~3,000–5,000 Hz
2nd bending: ~5,000–10,000 Hz
Axial stretch modes: > 20,000 Hz

**HCF safe stress:**
If a resonance CANNOT be avoided (speed range too wide), the stress AT resonance must be below the HCF limit:
```
σ_HCF_limit = σ_endurance / K_t + compensation for mean stress (Goodman)
```

**Damping mechanisms used in blades:**
1. **Friction damping (shrouds):** Blade tips contact adjacent blade's shroud → friction dissipates vibration energy
2. **Platform dampers:** Metal plate between blade roots → rubs → damping
3. **Internal damping:** Hysteresis in the alloy (very small for Ni superalloys at HCF frequencies)

---

## 7. Creep Under Combined Loading

**At service temperature, all stresses contribute to creep:**
- Primary: centrifugal stress (steady, high)
- Secondary: gas bending (steady, lower)
- Tertiary: thermal stress (partially relaxed)

**Centrifugal creep in [001] SX:**
Blade elongates under centrifugal stress → tip clearance (gap between blade tip and casing) closes → if blade elongates too much → tip rubs → catastrophic failure → must keep creep strain ≤ 0.5–1% over blade life

**Creep strain accumulation model:**
```
Δε_creep = ε̇_ss × Δt  (for steady-state portion)

Over 20,000 hours at 950°C / 130 MPa (CMSX-4):
ε̇_ss ≈ 1×10⁻⁵ /hour
Total ε_creep = 1×10⁻⁵ × 20,000 = 0.2% → within 0.5% limit
```

**Temperature sensitivity:**
A +25°C increase in TIT → k_p (TGO growth rate) increases → AND creep rate increases exponentially:
- If creep rate doubles → life halves
- This is why TIT increase of even 25°C requires new alloy generation or redesigned cooling

**Rafting effect on creep:**
After ~100 hours at 950°C: γ′ forms N-type rafts (plates perpendicular to [001] in tension):
- Stage I: creep rate DECREASES as rafts form (rafts are obstacles to dislocation passage)
- Stage II: steady state with raft structure
- Stage III: when rafts coarsen and lose coherency → accelerate → fracture

---

## 8. Thermomechanical Fatigue (TMF)

**TMF** is fatigue where temperature and mechanical strain cycle together:

**Two types:**
- **In-phase TMF (IP-TMF):** Maximum temperature coincides with maximum tensile stress → most damaging for oxidation-fatigue interaction
- **Out-of-phase TMF (OP-TMF):** Maximum temperature coincides with minimum (most compressive) stress → most damaging for mechanical cracking

**Turbine blade TMF cycle (typical):**
```
Ground idle → take-off (T increases, centrifugal stress increases):
  σ_mech increases (centrifugal)
  T increases (engine power up)
  This is approximately IN-PHASE → oxidation-fatigue

Take-off → cruise: stress stays high; T decreases slightly (cooling design)
Cruise → descent: T decreases rapidly (throttle back)
  σ_mech stays (centrifugal — engine still running)
  T drops → ε_thermal reversal
  This creates OUT-OF-PHASE condition during cooling
```

**TMF damage model (simplified):**
```
D_total = D_mech (fatigue) + D_ox (oxidation) + D_creep

1/N_f = [ε_mec / (2 × ε'_f)] ^(-1/c) + creep term + oxidation term
```

**Film cooling holes in TMF:**
Each film cooling hole in the airfoil is a stress concentration AND an oxidation initiation site:
- IP-TMF: tensile stress + high T → oxidation + fatigue interaction → crack initiation at hole
- K_t ≈ 3–4 at hole edge → local stress = 3–4 × nominal → dominant fatigue initiation site

---

## 9. Foreign Object Damage (FOD)

**FOD:** Hard debris (stones, nuts, bolts, bird fragments) ingested by the engine → hits rotating blades → causes dents, cracks, or leading edge damage.

**Types of FOD:**
- **Soft FOD (bird strike):** Large bird → high impact energy → deformation without sharp crack (bird is soft) → structural damage to fan but usually no immediate fracture of metal core
- **Hard FOD (stone, tool, bolt):** Creates sharp notch → stress concentrator → immediate HCF or TMF crack initiation

**FOD damage tolerance design:**
Blades are designed to be damage tolerant — they must survive a specific size FOD:
- Example specification: survive 1 mm × 5 mm notch at leading edge without fracture in 500 hours
- This requires: K_Ic high enough that a 1 mm defect doesn't immediately propagate

**FOD mitigation:**
- Inlet screen: mesh over engine inlet (helicopters, turboprops)
- Particle separators: swirl-type separator (debris centrifuged out before reaching core)
- FOD walk: crew inspects runway before every flight
- Blade design: leading edge may use softer, more ductile Ti-6Al-4V alloy (fan blades) vs brittle ceramics

**Shot peening after FOD damage:**
A blade that received FOD dent may be salvageable by:
- Blend/polish the dent to remove sharp edges (remove stress concentration)
- Shot peen → introduce compressive RS around the blend → delay crack initiation
- Re-measure geometry → if within tolerance after blend → accept for further service

---

## 10. Combined Life Prediction — Interaction of Mechanisms

**Real HPT blade life is controlled by SIMULTANEOUS action of:**

```
Life_total depends on:
├── LCF: from duty cycle (n_LCF cycles to failure N_f,LCF)
│         Miner's rule: D_LCF = n_LCF / N_f,LCF
├── HCF: from vibration (n_HCF cycles to failure N_f,HCF)
│         D_HCF = n_HCF / N_f,HCF
├── Creep: D_creep = ε̇_ss × t / ε_f (time fraction)
├── Oxidation: TGO growth → D_ox = (x_TGO/x_TGO,crit)^n
└── Combined TMF: interaction terms

Failure when: D_LCF + D_HCF + D_creep + D_ox ≥ 1.0
```

**Miner's Rule (linear damage summation):**
Simplest prediction model; often over-predicts life under IP-TMF conditions (interaction terms are non-linear).

**More advanced models:**
- Ostergren TMF model: accounts for stress-oxidation interaction
- Walker model (LCF): mean stress correction
- Robinson life fraction (creep): time fraction spent at each T-σ condition

**Design safety factor:**
All life predictions have large uncertainty → design includes safety factor of 2–4× on cycles:
- If analysis predicts N_f = 40,000 cycles → blade certified to N_LLP = 10,000–20,000 cycles
- After reaching LLP (life limit part): mandatory replacement regardless of apparent condition

---

## Summary

| Loading Type | Magnitude | Frequency | Primary Material Response | Design Approach |
|-------------|-----------|-----------|--------------------------|----------------|
| Centrifugal | 100–150 MPa | Steady | Creep (steady-state creep rate) | SX [001], high creep alloy |
| Gas bending | 30–60 MPa | Slowly varying | LCF when cycling | Low K_t at root fillets |
| Thermal gradient | 50–150 MPa | Per flight cycle | TMF, LCF at holes | Film cooling, [001] SX (low E) |
| Vibration | 10–50 MPa | HCF (kHz) | HCF (resonance fatigue) | Campbell diagram, dampers |
| FOD | Impulse | Rare events | Dent → HCF initiation | Damage tolerance, shot peen |

---

## Exercises

1. Centrifugal stress calculation for an HPT blade: N = 11,500 rpm, blade mass = 220g, root radius r_root = 0.32 m, tip radius r_tip = 0.44 m. (a) Calculate ω (rad/s). (b) The blade has constant cross-section area A = 350 mm². Using the simplified formula F_c = m × ω² × r_centroid (r_centroid = (r_root + r_tip)/2), calculate centrifugal force. (c) Calculate root centrifugal stress σ_c = F_c / A_root. (d) If the blade alloy is replaced with a Gen 5 alloy with density 9.3 g/cm³ vs Gen 2's 8.7 g/cm³, the blade mass increases proportionally. Recalculate F_c and σ_c. (e) The disk firtree attachment has design limit of 280 MPa bearing stress. Calculate safety factor for Gen 2 and Gen 5 blades.

2. Campbell diagram construction: A turbine blade has natural frequencies: 1st bending = 1,800 Hz, 1st torsion = 4,200 Hz, 2nd bending = 6,500 Hz. The engine operates from 8,000–11,000 rpm. Upstream stator count = 46. (a) Calculate excitation frequency for 1 EO, 2 EO, 3 EO, and 46 EO (stator wake) at both 8,000 and 11,000 rpm. (b) Which resonance crossings occur within the operating range? List the mode and EO for each crossing. (c) For the 1st bending/46 EO crossing: if the vibratory stress at resonance = 45 MPa and the HCF endurance limit (with mean stress correction) = 40 MPa, is this crossing safe? What options exist to make it safe? (d) At resonance, blade tip displacement = 1.5 mm. Calculate the maximum strain at the blade root (assume simple cantilever: ε_max = 6 × δ_tip × t / L² where t = 3 mm thickness, L = 120 mm blade length). Is this above the yield strain of CMSX-4 at that temperature?

3. Thermal fatigue life of film cooling holes: A trailing edge row of holes (0.5 mm diameter) experiences ΔT = 150°C per flight cycle. (a) Calculate thermal stress range Δσ_thermal = E × α × ΔT for [001] CMSX-4 (E = 130 GPa, α = 14×10⁻⁶/°C). (b) With K_t = 3.5 at hole edge, calculate local peak stress. (c) Using Basquin equation for HCF (σ_a = σ'_f × (2N)^b with σ'_f = 1,800 MPa, b = -0.1), and Coffin-Manson for LCF, and recognizing this is a LCF problem (< 10,000 flights), calculate the predicted cycles to crack initiation (ignore mean stress for simplicity). (d) After shot peening, the local RS = -400 MPa. Recalculate effective peak stress and revised fatigue life (use Goodman: σ_a/σ_end = 1 - σ_m/σ_UTS where σ_m = σ_mean + σ_RS).

4. Creep life management: An HPT blade (CMSX-4) operates at σ = 130 MPa, T = 950°C, ε̇_ss = 1×10⁻⁵ /hour. The blade is replaced when creep strain reaches 0.5%. (a) Calculate life based on steady-state creep only (ignore Stage I and III). (b) The engine experiences 30% of its time at higher power (T = 980°C). At 980°C, ε̇_ss = 3×10⁻⁵ /hour (Arrhenius: +30°C → creep rate ×3). The remaining 70% is at 950°C. Calculate the average creep rate and revised life. (c) The airline operator wants to extend blade life from 20,000 to 25,000 hours. What maximum steady-state TIT temperature could be tolerated if all else is equal? (use: life ∝ 1/ε̇_ss, and ε̇_ss ∝ exp(-Q/RT) to find ΔT for 20% reduction in creep rate) (d) Alternatively, if blade wall thickness is reduced 10% to allow more cooling air: the blade metal temperature drops 15°C. What is the new creep life?

5. FOD damage tolerance analysis: A Ti-6Al-4V fan blade (K_Ic = 44 MPa√m, σ_y = 900 MPa) suffers a 1.2 mm deep leading edge nick. (a) Using Griffith/LEFM: K_I = 1.1σ√(πa) (1.1 = geometry factor for edge crack), calculate K_I for σ = 300 MPa (centrifugal + bending). (b) Is K_I < K_Ic? Is immediate fracture predicted? (c) This is a stress-concentration question: the nick has root radius ρ = 0.1 mm. K_t ≈ 1 + 2√(a/ρ). Calculate K_t and local stress. Does local yielding occur? (d) After blend/polish: a = 0.3 mm (0.9 mm removed), ρ = 0.5 mm (smooth). Recalculate K_I and K_t. Is the blade now safe for return to service? (e) After re-blending, require shot peening (Almen intensity 8A). The compressive RS = -550 MPa to 200 μm depth. What is the minimum flaw size that could still propagate under fatigue loading after shot peening? (use K_threshold = 5 MPa√m for Ti-6Al-4V)
