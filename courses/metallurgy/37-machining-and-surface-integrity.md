# Chapter 37: Machining and Surface Integrity

> **"A turbine blade may spend 400 hours being cast, heat treated, coated, and inspected — but the film cooling holes, the trailing edge slots, and the blade-to-disk attachment firtree are created in the last few hours of machining. The surface integrity of those machined features — the residual stress, the microstructure below the surface, the roughness — determines the fatigue life of the blade. A machined surface is not just a dimension; it is a metallurgical state."**

---

## Table of Contents

1. [What Is Machining?](#1-what-is-machining)
2. [Chip Formation and Cutting Mechanics](#2-chip-formation-and-cutting-mechanics)
3. [Tool Materials and Wear](#3-tool-materials-and-wear)
4. [Surface Integrity — The Critical Concept](#4-surface-integrity--the-critical-concept)
5. [Residual Stresses from Machining](#5-residual-stresses-from-machining)
6. [Microstructural Damage in the Machined Layer](#6-microstructural-damage-in-the-machined-layer)
7. [Non-Traditional Machining: EDM, ECM, Laser](#7-non-traditional-machining-edm-ecm-laser)
8. [Grinding](#8-grinding)
9. [Shot Peening — Beneficial Compressive Residual Stress](#9-shot-peening--beneficial-compressive-residual-stress)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. What Is Machining?

Machining removes material by cutting, to achieve dimensional accuracy and surface finish. For metals:
- Turning (rotating workpiece, fixed tool)
- Milling (rotating tool, fixed or traversing workpiece)
- Drilling (rotating drill, axial feed)
- Grinding (abrasive wheel: many tiny cutting edges)

**Why machining matters metallurgically:**
- The machined surface is NOT the same as the bulk material
- Cutting deforms, heats, and may transform the near-surface layer
- Residual stresses (compressive or tensile) left in the surface affect fatigue life by orders of magnitude
- Surface roughness acts as a stress concentrator (K_t ∝ √(a/ρ) for notch root radius ρ)

---

## 2. Chip Formation and Cutting Mechanics

**Orthogonal cutting model:**
```
Tool (rake angle α)
  \  ← feed
   \
    ──────────────────────────────── workpiece
    Shear zone (φ = shear angle)
```

**Three chip types:**
1. **Continuous chip:** Ductile material + sharp tool + good lubricant → smooth chip, good surface
2. **Segmented/serrated chip:** Hard materials (Ti-6Al-4V, Ni superalloys), adiabatic shear bands form → serrated chips, intermittent cutting forces, risk of built-up edge
3. **Discontinuous chip:** Brittle material (cast iron) or unfavorable conditions → broken chips, rougher surface

**Cutting forces:**
- F_c (cutting force, tangential): power consumption → F_c = K_s × t × f (K_s = specific cutting force in N/mm²; t = depth; f = feed)
- Specific cutting force K_s depends on material:
  | Material | K_s (N/mm²) |
  |----------|------------|
  | Aluminum alloys | 700–1,000 |
  | Steel (medium C) | 2,000–3,000 |
  | Stainless steel | 2,500–3,500 |
  | Nickel superalloys | 3,000–5,000 |
  | Titanium alloys | 2,000–3,000 |

**Heat generation in cutting:**
- Most energy → heat (~98%): generated at (a) shear zone, (b) tool-chip interface, (c) tool-workpiece interface
- Tool temperature can reach 800–1,200°C in the cutting zone
- High T → tool wear; high T in workpiece surface layer → phase transformations, tensile residual stress

---

## 3. Tool Materials and Wear

**Tool materials (in order of hardness / heat resistance):**
```
High Speed Steel (HSS, M2): σ_y = 4,000 MPa, red hardness to 600°C → Al, soft steel
     ↓
Cemented Carbide (WC-Co): HV 1,600–2,000, to 800°C → steel, cast iron (coated with TiN/Al₂O₃)
     ↓
Cermet (TiC-Ni): clean cuts, good surface finish, hard materials
     ↓
Ceramics (Al₂O₃, Si₃N₄): very hard, to 1,200°C → cast iron, hardened steel (no shock)
     ↓
CBN (Cubic Boron Nitride): second hardest known, to 1,000°C → hardened steels, chilled CI
     ↓
Diamond (PCD): hardest, to 800°C → Al, non-ferrous only (reacts with Fe!)
```

**Tool wear mechanisms:**
- **Flank wear (VB):** Abrasion of tool flank against machined surface → dimensional changes
- **Crater wear:** Chemical diffusion of tool into chip at high T (major failure mode for uncoated carbide)
- **Notch wear:** At depth-of-cut boundary (hard oxide scale on Ti or Ni alloys)
- **Built-up edge (BUE):** Work material welds to tool → periodic break-off → poor surface finish
- **Thermal cracking:** Cyclic thermal stresses in interrupted cuts → comb cracks on tool face

**Taylor's tool life equation:**
```
V × T^n = C

where V = cutting speed (m/min)
      T = tool life (min)
      n = material constant (0.1–0.4 for carbide)
      C = constant
      
Higher V → shorter tool life exponentially (n governs sensitivity)
```

**Coatings on cemented carbide:**
- CVD TiN + TiC + Al₂O₃ multilayer (3–15 μm): reduces crater wear, chemical diffusion
- PVD TiAlN: better toughness at interrupted cuts (bending fatigue of coating)
- Commercial insert life: 10–30 min per cutting edge at recommended speeds

---

## 4. Surface Integrity — The Critical Concept

**Surface integrity:** The combined physical, mechanical, metallurgical, and chemical condition of the machined surface.

**Components of surface integrity:**
```
Surface integrity
├── Surface roughness (Ra, Rz, Rmax)
│    → stress concentration for fatigue cracks
├── Surface residual stresses
│    → compressive: beneficial (fatigue); tensile: harmful
├── Work hardening layer
│    → increased hardness + reduced ductility
├── Microstructural alterations
│    → white layers, tempered zones, recrystallized zones
└── Surface contamination
     → absorbed elements, lubricant residue, oxide scale
```

**Why surface integrity matters for turbine blades:**
- Film cooling holes: K_t ≈ 3–4 at hole edge
- If machining introduces TENSILE residual stress: fatigue initiation at hole accelerates
- If shot peened or ECM machined: compressive stress → delays crack initiation
- Difference in fatigue life: up to 5× depending on machining process and condition

---

## 5. Residual Stresses from Machining

**Two competing mechanisms leave residual stress:**

1. **Mechanical effect** (dominant at low cutting T, sharp tools):
   - Plastic deformation stretches surface laterally → compressed by surrounding elastic material on springback → **compressive residual stress**

2. **Thermal effect** (dominant at high cutting T, dull tools, high speed):
   - Surface heats rapidly → expands → compressed plastically
   - On cooling → surface wants to contract but is constrained → **tensile residual stress**

**Resultant surface stress = thermal - mechanical effects:**
- Sharp tool + low speed + flood coolant → mechanical dominates → compressive RS → GOOD
- Dull tool + high speed + dry cut → thermal dominates → tensile RS → BAD for fatigue

**Residual stress measurement:**
- X-ray diffraction (XRD): non-destructive; measures strain in crystal lattice from peak shift; standard method for surface RS
- Hole drilling: drill small hole → strain gauge rosette measures strain relief → calculate RS
- Neutron diffraction: measures RS in depth (non-destructive but requires neutron source)

**Example:**
| Process condition | Surface RS (MPa) | Fatigue life ratio |
|-------------------|-----------------|-------------------|
| Conventional turning (sharp) | -200 (compressive) | 1.0 (baseline) |
| Turning (worn tool) | +400 (tensile) | 0.4 |
| Grinding (aggressive) | +600 (tensile) | 0.25 |
| Shot peening after turning | -600 to -800 | 2.5–5.0 |

---

## 6. Microstructural Damage in the Machined Layer

### White Layer (WL)
The most serious microstructural damage from machining:
- Appears **white** in optical micrograph (doesn't etch with nital) → hence "white layer" or "white etching layer"
- Very hard (1,000–1,200 HV), brittle
- **Formation mechanisms:**
  - Rapid heating → austenite → rapid cooling → martensite (untempered: very hard, brittle)
  - Severe plastic deformation → extremely fine-grained ferrite (nanocrystalline)
  - Phase transformation: depends on alloy

**White layers in:**
- Hard steel (bearing races, gear teeth after aggressive grinding): most concern; white layer = tensile RS + hard brittle → fatigue crack origin
- Nickel superalloys (IN718 after aggressive turning): amorphous or nanocrystalline Ni
- Ti-6Al-4V: difficult to observe but adiabatic shear bands related

**White layer depth:** 1–50 μm depending on aggressiveness. Even 1 μm of white layer can cut fatigue life by 30%.

### Over-tempered (Over-tempered zone)
Below the white layer in hardened steels: the next zone was tempered (softened) by the heat:
- Darker in SEM backscatter (lower hardness)
- Reduced hardness compared to bulk → soft layer

### Work-hardened layer
In austenitic stainless, Ni alloys, Al alloys:
- Machining deforms surface → increased dislocation density → hardened layer 10–100 μm deep
- Residual stress profile: usually compressive near surface → tensile subsurface

---

## 7. Non-Traditional Machining: EDM, ECM, Laser

Non-traditional machining processes are essential for hard-to-machine alloys (nickel superalloys, tungsten carbide) and complex geometries (turbine blade cooling holes).

### Electrical Discharge Machining (EDM)
**Principle:** Spark erosion — electric sparks between electrode and workpiece erode both:
```
Electrode (tool) ─────────────
                  ← spark gap (0.02–0.5 mm), filled with dielectric fluid
Workpiece ────────────────────
```
Each spark melts + vaporizes a tiny crater; thousands of sparks per second → material removal.

**Characteristics:**
- Any electrically conductive material can be machined regardless of hardness
- No mechanical force on workpiece (suitable for thin walls, delicate parts)
- Surface: recast layer (melted + resolidified) → white layer, tensile RS, microcracks
- Film cooling holes in turbine blades: EDM drill (rotating electrode tube)

**EDM surface integrity:**
- Recast layer: 1–30 μm of resolidified material → brittle, tensile residual stress
- Must be removed by acid etching (aerospace specifications: max 0.001" recast allowed)
- Subsequent electrochemical polishing can remove the recast layer

### Electrochemical Machining (ECM)
**Principle:** Reverse electroplating — workpiece is anode, metal dissolves away:
```
Tool (cathode, -) ─── electrolyte flow (NaCl, NaNO₃) ─── workpiece (anode, +)
                        Metal dissolves away at anode (oxidation)
```

**Characteristics:**
- NO thermal damage (cool electrochemical process)
- No tool wear
- Excellent surface integrity: compressive RS, no white layer
- Used for: turbine blade firtree (dovetail attachment), film cooling holes (STEM — shaped tube electrochemical machining)
- **STEM drilling:** For turbine blade cooling holes — Ni alloy blade → hundreds of holes, each 0.4–1 mm dia, 40–80 aspect ratio → only STEM can achieve this without thermal damage

### Laser Machining
- Pulsed laser: drill holes (percussion drilling) or cut
- Thermal process → heat-affected zone, recast layer
- Fast; suitable for non-conductive ceramics (EDM doesn't work)
- Hole quality: inferior to ECM but faster setup

---

## 8. Grinding

**Grinding** uses abrasive wheels to remove small amounts of material with high dimensional accuracy and surface finish. It is often the FINAL manufacturing operation.

**Thermal risk in grinding:**
- Grinding specific energy (J/mm³) is very high → much more heat per unit volume than turning
- Surface can reach 700–1,000°C even with coolant
- Result: tensile residual stress, white layer, if abusive: workpiece burns → blue oxidation → dimensional change

**Grinding burn indicators:**
- Surface blue/yellow/brown oxide color (oxidation → overheating)
- Magnetic testing (in steels): loss of magnetism where martensite reverts to austenite
- XRD: tensile RS vs expected compressive RS

**Grinding parameters for good surface integrity:**
- Lighter down-feed
- Higher wheel speed + lower table speed
- Sharp, properly dressed wheel
- Flood coolant (not mist)
- Avoid redress sparks without coolant

**CBN wheels for hardened steels:**
- CBN (cubic boron nitride): maintains sharpness >> aluminum oxide
- Lower grinding temperature → less thermal damage
- Required for aerospace landing gear, bearing races

---

## 9. Shot Peening — Beneficial Compressive Residual Stress

**Shot peening** deliberately introduces compressive residual stress to improve fatigue life:
```
High-velocity steel shot (0.2–2 mm diameter) → impact workpiece surface
  → plastic deformation of surface layer
  → compressed by surrounding elastic material
  → NET: compressive residual stress layer (~200–600 μm deep, -200 to -800 MPa)
```

**Why compressive RS is beneficial:**
Fatigue cracks initiate and propagate when local stress exceeds threshold. With compressive RS:
```
σ_local = σ_applied + σ_residual
If σ_residual is compressive (negative) → σ_local < σ_applied → harder to initiate/grow crack
```

**Almen strip test:** Measure arc height of thin steel strip after peening → quantify intensity (coverage and energy)

**Fatigue life improvement from shot peening:**
- Carburized gears: 2–3× improvement in bending fatigue life
- Aircraft landing gear (steel): 2× improvement
- IN718 turbine disk: compressive RS offsets 25–30 MPa of mean stress (Goodman diagram shift)

**Controlled (precision) shot peening (CASP):**
Computer-controlled position + particle velocity for uniform, reproducible coverage. Mandatory specification for:
- Turbine blade airfoil surfaces
- Compressor disk firtree features (Ti-6Al-4V)
- Aircraft engine connecting rods

**Laser peening:**
- Intense laser pulse → plasma → pressure wave → compressive RS
- Deeper compressive layer (2–5 mm) than conventional shot peening (~0.5 mm)
- No cold work of surface (pure pressure wave) → better surface finish maintained
- Cost: high (laser equipment), but justified for critical components

---

## Summary

| Process | Surface RS | Microstructural Damage | Use case |
|---------|-----------|------------------------|---------|
| Sharp turning | Compressive | Minor work hardening | General machining |
| Worn tool turning | Tensile | White layer possible | Avoid |
| Aggressive grinding | Tensile (burn) | White layer, tempered zone | Avoid; use CBN/careful params |
| EDM | Tensile | Recast layer | Hard alloys; must etch to remove recast |
| ECM/STEM | Compressive | None | Preferred for turbine holes |
| Shot peening | Compressive (deep) | Work hardening | Post-machining beneficial treatment |
| Laser peening | Compressive (very deep) | Minimal | Critical components |

---

## Exercises

1. Turning of IN718 turbine disk slot using carbide insert (WC-Co, TiAlN coated): (a) The cutting speed is 30 m/min (low, typical for superalloys). Why is superalloy machining limited to much lower cutting speeds than steel? Consider tool wear mechanisms and thermal conductivity (Ti: 6.7 W/m·K vs steel: 50 W/m·K). (b) After 12 cuts, the insert shows 0.3 mm flank wear (VB_max = 0.3 mm is the rejection criterion). Using Taylor's equation with n = 0.15, C = 100, calculate the tool life at V = 30 m/min. At V = 50 m/min. (c) Surface roughness target Ra ≤ 0.8 μm. At V = 30 m/min, Ra = 0.4 μm. At V = 50 m/min (worn tool), Ra = 2.0 μm. Explain why worn tools produce rougher surfaces. (d) XRD shows: V = 30 m/min, sharp tool → RS = -150 MPa. V = 50 m/min, worn tool → RS = +320 MPa. What is the fatigue life implication? (see Ch 14, Paris Law, for context).

2. EDM drilling of film cooling holes in CMSX-4 SX blade: 100 holes × 0.5 mm diameter × 20 mm deep. (a) Describe the material removal mechanism (spark erosion). Why is this process preferred over conventional drilling for these holes? (b) SEM analysis shows 15 μm recast layer and microcracks from the recast layer extending 20 μm into the base material. Why is this recast layer harmful for fatigue life? (c) The specification requires recast layer ≤ 2 μm. What post-process removes the recast layer? (d) Alternatively, STEM (electrochemical) drilling produces the same holes with zero recast. Why is ECM fundamentally free from thermal damage? (e) Cost comparison: EDM drilling + electrochemical polish = $85/blade; STEM = $130/blade. From a fatigue life standpoint (blade life 40,000 cycles EDM; 90,000 cycles STEM), what is the cost per fatigue cycle for each route? Which is more economical?

3. Shot peening of Ti-6Al-4V compressor blade: Almen intensity 6A, coverage 100%, steel shot S110. (a) Post-peening XRD profile (depth vs RS): surface RS = -380 MPa; maximum compressive depth = 150 μm; crossover to tensile RS at 180 μm depth. Draw this RS profile schematically. (b) In the Goodman diagram, the mean stress axis is offset by the RS value. If applied mean stress σ_m = 300 MPa and alternating stress σ_a = 200 MPa, with unpeened endurance limit = 450 MPa and UTS = 960 MPa: check if the condition is safe on the Goodman diagram for unpeened. (c) With shot peening, the effective mean stress = σ_m + σ_residual = 300 + (-380) = -80 MPa (compressive). Recheck on Goodman diagram. Is the peened component safe? (d) Why would deeper laser peening (compressive to 2 mm depth) be preferred over shot peening for this blade if it experiences bird-strike FOD (Foreign Object Damage) creating a surface dent of 0.5 mm depth?

4. White layer formation in gear grinding: bearing quality steel 52100 (1.0%C, 1.5%Cr, 0.25%Mn) hardened to 62 HRC. (a) During aggressive grinding, the surface reaches 900°C momentarily. Describe the phase transformation sequence: what happens to the martensite at 900°C? On rapid cooling (coolant), what new phase forms? (b) The new phase has hardness 900 HV. What is this phase? Why is it so hard? What is different about its residual stress state vs the surrounding tempered martensite? (c) White layer thickness = 20 μm. After 10⁶ bending fatigue cycles, bearing shows early spalling. Explain how the white layer initiates the spall. (d) How would switching from conventional Al₂O₃ grinding wheels to CBN wheels reduce white layer formation? What property of CBN makes the difference?

5. Surface integrity qualification for aerospace turbine disk (IN718). The specification requires: (a) Ra ≤ 0.8 μm, (b) no white layer (recast) ≥ 0.3 μm, (c) compressive residual stress in first 100 μm, (d) no re-hardening zones. Describe the complete inspection procedure: (i) which measurement technique for each specification item, (ii) sampling plan (every part, sampling), (iii) consequence of non-conformance (scrap or rework?), (iv) corrective process change if white layer is found on multiple consecutive parts.
