# Chapter 58: Turbine Blade Failure Modes — How Blades Die

> **"Every turbine blade failure is a detective story. The blade came out of service after 12,000 hours. It lost 2mm of wall thickness on the pressure side leading edge. The TBC is missing in a 10cm² patch on the suction side. There is a hairline crack at the cooling hole exit on row 3. The fracture face shows Stage II fatigue striations converging from two origins — one at the cooling hole, one at a TBC-spall pit. Reading this evidence tells you exactly what happened, in what order, and what changed in the engine operating profile that allowed it."**

---

## Table of Contents

1. [Overview of Blade Life and Failure Modes](#1-overview-of-blade-life-and-failure-modes)
2. [Creep — Dimensional Distortion and Rupture](#2-creep--dimensional-distortion-and-rupture)
3. [Thermal Fatigue (LCF) — Cycle-by-Cycle Damage](#3-thermal-fatigue-lcf--cycle-by-cycle-damage)
4. [High-Cycle Fatigue (HCF) — Vibration-Driven Cracking](#4-high-cycle-fatigue-hcf--vibration-driven-cracking)
5. [Oxidation — Wall Thinning](#5-oxidation--wall-thinning)
6. [Hot Corrosion — Accelerated Attack](#6-hot-corrosion--accelerated-attack)
7. [TBC Spallation — Loss of Thermal Protection](#7-tbc-spallation--loss-of-thermal-protection)
8. [Erosion — Particle Impact Wear](#8-erosion--particle-impact-wear)
9. [FOD and DOD — Foreign and Domestic Object Damage](#9-fod-and-dod--foreign-and-domestic-object-damage)
10. [Interaction Effects — How Failures Cascade](#10-interaction-effects--how-failures-cascade)
11. [Failure Investigation — Reading the Evidence](#11-failure-investigation--reading-the-evidence)
12. [Life Management and Blade Retirement](#12-life-management-and-blade-retirement)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. Overview of Blade Life and Failure Modes

A modern HPT blade is designed for a **life limit** — typically 30,000 engine flight hours (EFH) or 20,000–30,000 cycles (takeoff-landing cycles), whichever comes first.

In practice, most blades are removed from service for **inspection** at intermediate intervals (4,000–10,000 EFH), examined, and either returned to service, repaired, or scrapped. Very few blades actually reach their design life limit before removal for inspection and repair reasons.

**Failure modes in rough order of frequency:**

| Rank | Failure Mode | When It Typically Occurs |
|------|-------------|--------------------------|
| 1 | Creep elongation / tip clearance loss | Gradual; 10,000–30,000 EFH |
| 2 | Thermal fatigue (LCF) cracking | 10,000–20,000 cycles |
| 3 | TBC spallation | 8,000–15,000 EFH |
| 4 | Oxidation / wall thinning | Gradual; 5,000+ EFH |
| 5 | Hot corrosion | Depends on fuel/environment |
| 6 | High-cycle fatigue | Can be sudden; resonance-dependent |
| 7 | Erosion | Gradual; depends on dust ingestion |
| 8 | FOD | Sudden; event-driven |

**Key distinction: gradual vs. sudden:**
- Gradual failures accumulate slowly → detectable by inspection → managed by condition monitoring
- Sudden failures (FOD, HCF) can cause immediate engine damage if undetected → must be prevented by design

---

## 2. Creep — Dimensional Distortion and Rupture

### What Happens

Under sustained centrifugal stress at high temperature, the blade slowly elongates and distorts (Chapter 16). Three aspects matter in service:

**Tip elongation:**
The blade tip moves outward radially due to creep elongation. As the blade gets longer:
- Tip clearance (gap between blade tip and engine casing) decreases
- When tip clearance = 0: blade tip rubs against casing → catastrophic damage
- In some engines, small amounts of tip rub are acceptable; but sustained rubbing generates heat → local overtemperature → accelerated creep

**Airfoil untwist:**
Creep can alter the airfoil twist angle (the aerodynamic angle varies from root to tip). As the blade untwists:
- Aerodynamic efficiency decreases
- Mass flow through the turbine changes
- Eventually: engine performance out of spec

**Section thinning:**
Sustained stress + temperature → some creep occurs by grain boundary sliding in the γ channels between γ′ cubes → sections can thin.

### Creep Rupture

If the blade is severely overtemperature (cooling failure, TBC spallation, engine over-speed) and reaches the tertiary creep stage (Chapter 16), rupture occurs:
- Blade fractures at or near the root (highest stress section)
- Fragment is released at very high velocity into the engine
- Can cause immediate catastrophic engine failure (loss of engine)

This is the design-driver for "safe life" limits: blades are removed BEFORE their Larson-Miller parameter budget is consumed.

### Typical Numbers

At normal operation (CMSX-4, 980°C, 200 MPa):
- Expected secondary creep rate: ~10⁻¹⁰ per second
- Over 30,000 EFH (~10⁸ seconds): total creep strain ~1% → ~0.8mm tip elongation on 80mm blade
- This is acceptable; blade is designed with this elongation budget

---

## 3. Thermal Fatigue (LCF) — Cycle-by-Cycle Damage

### The Physical Process

Every engine startup-to-shutdown constitutes one low-cycle fatigue (LCF) cycle (Chapter 54). The temperature gradient through the blade wall creates thermal strain — and this strain reverses during heating vs. cooling:

```
LCF mechanism:

Heating:   alloy expands FASTER than TBC → TBC in COMPRESSION
           (T_alloy > T_TBC surface during transient)

Cooling:   alloy contracts FASTER than TBC → TBC in TENSION
           metal interior contracts → outer surface and TBC in tension

After many cycles: cracks nucleate at:
  - Film cooling hole exits (stress concentrator)
  - TBC-spalled regions (newly exposed metal in tension)
  - Blade platform edges
  - Tip cap cooling holes
```

### Film Cooling Hole Cracking

Film cooling holes are the most common LCF initiation sites on the blade airfoil:
```
                TBC
              ▓▓▓▓▓▓▓▓▓▓▓▓▓
    ──────────────────────────────────────── Outer blade surface
              ↑   ↑   ↑   ↑   ↑
           Film holes (0.3-0.5mm dia.)
    ──────────────────────────────────────── Inner passage wall

At each film hole:
  - Local stress concentration factor Kt ≈ 3–4 (for a hole in a plate)
  - Combined thermal + centrifugal stress
  - Crack nucleates at hole exit on suction side (highest tensile stress)
  - Crack propagates circumferentially around the hole → "horse-shoe" crack
  - If two adjacent holes link: wall section between them can fail
```

**Preventing film hole cracking:**
- Shaped holes (fan/laidback) have lower Kt than cylindrical holes
- Proper hole angle → film flow reduces surface temperature → reduces thermal stress
- Crystal orientation: secondary crystal orientation affects how much stress the holes experience

### Platform and Root Fatigue

The blade platform (horizontal section between airfoil and root) and fir-tree root are also LCF-critical:
- Platform edges: thin sections, geometric stress concentrations, platform-to-platform fretting
- Fir-tree root: contact stresses between blade and disk slot → combined fretting + LCF

---

## 4. High-Cycle Fatigue (HCF) — Vibration-Driven Cracking

### Source of Excitation

HPT blades rotate past stationary upstream vanes. Each time a blade passes the wake of a vane, it receives an aerodynamic excitation impulse. This occurs at blade-passing frequency:

```
Excitation frequency = N_vanes × Ω / 60

For N_vanes = 24, Ω = 10,000 RPM:
f = 24 × 10,000/60 = 4,000 Hz
```

If any natural vibration frequency of the blade coincides with this excitation → **resonance** → HCF loading.

### Campbell Diagram

The **Campbell diagram** shows natural frequencies vs. rotation speed and identifies resonance crossings:

```
Campbell Diagram (schematic):

Natural
frequency
(Hz)
  |          ×  ×  ×  ← crossing at operating speed = CRITICAL!
  |         ×  1st bending mode
3000|        ×
  |    ×──────────────←  "engine order" lines (Nω)
  |   ×  2nd engine order
2000|  ×     
  |  × 1st engine order (= 1×)
1000|×  
  |
  └────────────────────── Rotation speed (RPM)
     4000   8000   12000
               ↑
          Operating speed
```

Any intersection of a natural frequency line with an engine order line at the operating speed is a potential HCF concern.

**Design goal:** Ensure no natural frequencies coincide with engine orders at the operating speed. This is done through:
- Blade geometry design (change natural frequencies)
- Detuning: deliberately making adjacent blades slightly different → break up traveling wave resonance
- Dampers: friction dampers under blade platform (energy dissipation → reduce resonance amplitudes)

### HCF Failure Pattern

HCF cracks grow faster than LCF cracks per cycle (because there are 10⁹ HCF cycles vs. 30,000 LCF cycles). But HCF stress amplitudes are lower (50–150 MPa vs. 200–300 MPa for LCF).

HCF failure pattern:
- Crack initiates at stress concentration (cooling hole, surface defect, erosion pit)
- Propagates at near-constant rate per cycle → catastrophically fast in total time (millions of cycles per minute)
- Often appears as a flat fracture face perpendicular to the stress direction (transgranular)
- Fatigue striations under SEM: parallel lines spaced at ~10⁻⁶ to 10⁻⁴ mm per cycle

---

## 5. Oxidation — Wall Thinning

### The Process

Despite TBC protection, oxidation attack occurs wherever TBC is absent or compromised:
- TBC-spalled regions
- Blade tip (tip rubs remove TBC)
- Film cooling hole walls (cool air exiting → oxygen access at the hole wall)
- Bond coat/TGO system consuming Al from bond coat

### External Oxidation Rate

Without TBC, at 1,100°C:
```
Parabolic oxidation law: Δm² = k_p × t (where k_p is parabolic rate constant)
```

For Ni superalloy with MCrAlY bond coat at 1,100°C:
- k_p ≈ 0.05 (mg/cm²)²/h
- After 5,000 hours: Δm = √(0.05 × 5,000) ≈ 15.8 mg/cm²
- Converted to thickness: ~15 μm of metal consumed

This 15 μm is not catastrophic on a 2.5 mm wall. But:
- Thinning means higher stress in remaining wall → accelerated creep
- Oxidation combined with hot corrosion → rates 10–100× faster

### Internal Oxidation

In the cooling passages, oxygen from the cooling air can react with exposed alloy surfaces (especially at cooling holes). At 700°C (coolant temperature), oxidation is slow — but at 900°C exit temperatures, it becomes significant over 30,000 hours.

Al and Ti diffuse to the surface → preferential oxidation → Al depletion from γ′ → local microstructure degradation → reduced local strength.

---

## 6. Hot Corrosion — Accelerated Attack

Hot corrosion is the most aggressive form of blade surface attack. As described in Chapter 54, it occurs when Na₂SO₄ (from NaCl in combustion + SO₃ from sulfur in fuel) deposits on the blade surface.

### Type I Hot Corrosion (800–950°C)

Characteristic features:
- Gray/black corrosion pits penetrating into the alloy
- Sulfide inclusions deep in the corroded zone (below the scale) — the "sulfidation" signature
- Extensive Cr depletion in the alloy below the scale

Mechanism: Na₂SO₄ (molten at 884°C) dissolves Cr₂O₃ scale → fresh alloy exposed → unprotected alloy oxidizes and sulfidizes at catastrophic rate.

### Type II Hot Corrosion (650–800°C)

Pitting corrosion: discrete, sharp-edged pits rather than uniform attack. Characteristic:
- NiSO₄ forms as a liquid deposit
- Reaction front moves inward at each pit

### Hot Corrosion Prevention

1. **Low-sulfur fuel:** ASTM limits S to < 3,000 ppm; military specs sometimes tighter.
2. **MCrAlY bond coat:** Cr₂O₃ is the primary defense against Type I; Al₂O₃ for Type II.
3. **Chromia-former coatings:** Some blades use a CrAlY overlay (higher Cr than MCrAlY) for naval/marine environments with high NaCl ingestion.
4. **Engine air filtration:** Prevents NaCl ingestion at ground level.

---

## 7. TBC Spallation — Loss of Thermal Protection

As described in Chapter 57, TBC spallation causes an immediate blade metal temperature increase:

```
Consequences of local TBC loss (0.25 mm TBC spall):

Before spall: T_metal = 980°C
After spall:  T_metal jumps to ~1,100–1,120°C (TBC thermal resistance removed)

At T_metal = 1,100°C instead of 980°C:
  Creep rate: ~10× faster (exponential T dependence)
  Oxidation rate: ~5× faster
  Time to blade replacement: significantly shortened
```

TBC spallation is often a precursor to other failures rather than a primary failure itself. Once a patch of TBC is gone:
1. Metal temperature rises locally
2. Oxidation accelerates at the spall edge
3. Thermal gradient through remaining TBC increases → drives further delamination at the spall edge
4. Spall expands → more metal exposed → more overtemperature

**Detection:** Spallation is detectable by thermal imaging of the engine in service (hot spots visible on thermal cameras at the exhaust). Modern engines have real-time exhaust thermal monitoring.

---

## 8. Erosion — Particle Impact Wear

### Mechanism

Fine particles (sand, dust, ingested solids) impact the blade airfoil surface at high velocity. The impact:
- Removes TBC material (ceramic TBC is relatively hard but brittle → particle impacts cause chipping)
- Removes metal from blade tip and leading edge
- Roughens the blade surface → increases aerodynamic drag → reduces efficiency

### Blade Tip Erosion

The blade tip is particularly vulnerable:
- Tip speed: ~500 m/s
- Any particle trajectory that crosses the tip hits at ~500 m/s relative velocity
- Progressive tip rounding → tip clearance increases → hot gas leaks over tip → efficiency loss

Modern HPT tips are protected by:
- Abrasive tip (small ceramic particles brazed into tip cavity) → abrades the shroud instead of the blade tip
- Hard face coating on the tip (MCrAlY + CrC)

### Leading Edge Erosion

Larger particles can erode the leading edge stagnation region. Even 0.1 mm of LE erosion can:
- Change the blade's aerodynamic profile significantly
- Thin the leading edge wall → higher thermal stress
- Expose fresh alloy → accelerated oxidation

---

## 9. FOD and DOD — Foreign and Domestic Object Damage

**Foreign Object Damage (FOD):** Impact from items NOT part of the engine — birds, tools left in the intake, runway debris. FOD events create immediate, severe damage: blade nicks, cracks, or complete blade section loss.

**Domestic Object Damage (DOD):** Fragments originating FROM the engine — shed TBC flakes, broken compressor blade fragments, loose nuts/bolts in the air stream. DOD can propagate through multiple turbine stages.

Both can cause:
- Local cracks → HCF crack propagation → rapid failure
- Blade imbalance → vibration → HCF on other blades
- Complete blade section loss → immediate engine failure

FOD/DOD are rare but catastrophic — they are the leading cause of in-flight engine failures that require emergency landings.

**Mitigation:**
- Foreign Object Debris monitoring programs (FOD walks before operations)
- Self-check inspection on engine start
- Strict torque specifications to prevent nut/bolt loosening
- Progressive FOD arrest designs: blade containment rings prevent fragment exit from engine nacelle

---

## 10. Interaction Effects — How Failures Cascade

In practice, blade failures are rarely single-mechanism. Common cascades:

**TBC spallation → oxidation → fatigue:**
1. TBC spall exposes bare MCrAlY bond coat
2. MCrAlY oxidizes at accelerated rate (no TBC insulation → hot bond coat)
3. Al depletes from MCrAlY → can no longer form protective Al₂O₃ → Ni/Co oxides form → scale spalls
4. Fresh alloy exposed → hot corrosion if sulfur present
5. Corroded, oxidized surface is rough, pitted → HCF crack nucleation from corrosion pits

**Cooling hole HCF → LCF link-up:**
1. HCF crack initiates at cooling hole exit under vibratory loading
2. Crack grows slowly to ~1 mm depth
3. LCF cycling (each flight) drives crack deeper through the wall
4. Crack penetrates wall → cooling flow disrupted → local overtemperature
5. Local creep accelerated → wall section fails

**Creep elongation → tip rub → fatigue:**
1. Creep elongation over 15,000 hours reduces tip clearance by 0.5 mm
2. Blade tip now contacts the casing lightly
3. Tip contact creates dynamic load → HCF loading on blade
4. HCF crack initiates at tip → propagates downward along blade
5. Loss of blade section → imbalance → cascade failure

---

## 11. Failure Investigation — Reading the Evidence

When a blade is removed and examined, the failure analyst looks for:

### Macroscopic Examination

| Observation | What It Means |
|-------------|---------------|
| Blade tip rounded/blunted | Tip erosion or tip rub |
| TBC missing in defined patch | TBC spallation (delamination or impact-caused) |
| Dark gray/black pits on blade surface | Hot corrosion (sulfidation) |
| Blade elongated > 0.5% over design | Significant creep accumulation |
| Crack at film hole exit | LCF or HCF initiation at stress concentration |
| Blade section missing | FOD/DOD or HCF fracture at resonance |

### Microscopic Evidence

**Fatigue striations:** Parallel lines on fracture face, spacing = crack growth rate per cycle. HCF: 10⁻⁴ μm/cycle (very fine). LCF: 10⁻¹ μm/cycle (coarser).

**Creep voids:** In SX blade, creep voids appear as elongated cavities in the γ channels (from dislocation accumulation). Their density indicates creep damage fraction.

**TGO thickness:** Thick TGO (> 8 μm) indicates extensive service time at temperature. TGO morphology (smooth vs. rumpled) indicates whether temperature was excessive.

**γ′ coarsening:** As-new blade γ′ size: 400–500 nm. After heavy service: can grow to 1–5 μm (rafted microstructure). Coarsening ratio indicates service temperature history.

**EBSD/X-ray Laue:** Checks if recrystallization occurred (should be absent in a properly processed SX blade). Any recrystallization → new grain boundaries → stress concentration → failure initiation.

---

## 12. Life Management and Blade Retirement

### Condition-Based Maintenance

Modern engines use a combination of:
1. **Fixed life limits**: Manufacturer-defined maximum service life (hours or cycles) — the absolute hard limit
2. **Condition monitoring**: Inspection at each shop visit → remove blades that exceed damage criteria
3. **Predicted remaining life**: Computations using flight data (temperature exposure, cycle count) to predict remaining creep/LCF life

### Repair vs. Replace Decision

Blades removed at shop visits are typically examined and sorted:
- **Return to service as-is**: No damage beyond acceptable limits
- **Repair**: Weld repair of cracks < threshold size, TBC reapplication, HVOF re-coating, tip brazing
- **Scrap**: Beyond repairability (large cracks, excessive creep elongation, section loss)

**Economics:** A new HPT blade costs $10,000–50,000. A blade repair costs $2,000–5,000. If repair extends life by 8,000 hours, repair is strongly preferred.

### Retirement Criteria (Typical)

| Damage Type | Acceptance Limit | Retirement Criterion |
|-------------|-----------------|----------------------|
| Creep elongation | < 0.3% strain | > 0.5% strain |
| Tip wear | < 0.5 mm | > 1.0 mm |
| Crack length | < 1 mm (non-critical location) | > 2 mm or any crack in critical zone |
| Wall thinning | < 10% nominal | > 20% nominal |
| TBC missing area | < 5 cm² | > 10 cm² or any miss over LE/PS high heat zone |
| Corrosion pit depth | < 0.2 mm | > 0.5 mm |

---

## Summary

- **Creep**: gradual elongation and distortion under centrifugal load at temperature; rupture if temperature limit exceeded; life-limited by Larson-Miller parameter budget
- **LCF**: thermal cycling (one per flight) drives fatigue cracks from stress concentrations (film holes, platforms, root); ~30,000 cycle design life
- **HCF**: vibration at blade-passing frequency; sudden failure if resonance crossings occur at operating speed; prevented by Campbell diagram design and dampers
- **Oxidation**: wall thinning at ~15 μm/5,000h; accelerated without TBC; parabolic kinetics
- **Hot corrosion**: catastrophic sulfidation at 800–950°C (Type I) or pitting at 650–800°C (Type II); fuel sulfur + NaCl + high temperature
- **TBC spallation**: loss of 100–150°C of temperature protection; immediate creep/oxidation acceleration; detectable by thermal imaging
- **Erosion**: particle impact at 500 m/s tip speed; tip rounding, LE erosion
- **FOD/DOD**: sudden, catastrophic; rare but dominant in-flight emergency driver
- **Cascades**: most real failures involve multiple interacting mechanisms
- **Life management**: condition monitoring + fixed limits + repair vs. replace economics

**Next chapter:** How failed blades are repaired — the specific processes (welding, TBC replacement, tip brazing, dimensional restoration) and the metallurgical challenges of repairing a material with no grain boundaries.

---

## Exercises

1. An HPT blade is at 980°C with a Larson-Miller parameter (LMP = T(log t_r + 20), T in K) budget of 32,000. (a) How many hours of life does this represent? (b) If the engine is accidentally over-tempered to 1,020°C for 100 hours, what LMP is consumed in that event? (c) By how much does this reduce the remaining life at 980°C?

2. A blade is inspected after 12,000 hours. Fracture mechanics analysis of a detected 0.8 mm crack at a film hole shows it will grow to the critical size (2.0 mm) by: (a) HCF alone after 5×10⁷ more cycles at Δσ = 80 MPa (use Paris law da/dN = C×ΔK^m, C = 10⁻¹², m = 3, ΔK = Δσ√(πa)), and (b) LCF alone after N_LCF more flights at Δσ_LCF = 250 MPa (same Paris parameters). Which mechanism is life-limiting? Should the blade be repaired now or returned to service?

3. A Campbell diagram shows the 1st bending mode at 3,200 Hz at operating speed of 10,000 RPM. The engine has 24 HPT vanes upstream. (a) What is the 1st engine order excitation frequency at 10,000 RPM? (b) What engine order excites the 1st bending mode at operating speed? (c) Is this a problematic resonance crossing? What would you do to avoid it?

4. TBC spallation occurs on 5 cm² of the blade pressure side at 10,000 EFH. The local blade metal temperature rises from 980°C to 1,080°C. Using the Arrhenius relationship for creep rate (ε̇ ∝ exp(-Q/RT), Q = 300 kJ/mol), calculate: (a) the ratio of creep rate at 1,080°C to 1,080°C at 980°C, (b) If the remaining blade life at 980°C was 20,000 EFH, estimate the life reduction caused by the 100°C overtemperature in the spalled region.

5. A hot corrosion pit 0.4 mm deep is found on the blade pressure side at mid-airfoil. The blade wall is 2.0 mm nominal thickness. (a) What fraction of wall thickness has been consumed? (b) The original stress at this location was 200 MPa. What is the local stress after pitting (assuming stress concentrates proportionally to remaining area)? (c) This blade is in service at 950°C and 200 MPa nominal stress — after pitting, does the local LMP still fall within CMSX-4 capability at the new stress level? (Use: for CMSX-4, 1000h rupture at 950°C: 350 MPa; at 1000h, 950°C, 270 MPa: marginal.)
