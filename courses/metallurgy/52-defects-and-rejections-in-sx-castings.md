# Chapter 52: Defects and Rejections in Single Crystal Castings

> **"Making a single crystal turbine blade requires growing a crystal the size of your palm at millimeters-per-hour over 30 centimeters of directional solidification, without a single stray grain nucleating, without the crystal losing its preferred direction, without pores collapsing into voids, and without any of the dozens of other possible defects arising. A 15–20% rejection rate is considered excellent in this industry. Understanding defects is not just academic — it is the key to improving yield and reducing cost."**

---

## Table of Contents

1. [Overview — Defect Categories in SX Castings](#1-overview--defect-categories-in-sx-castings)
2. [Stray Grains (Spurious Grains)](#2-stray-grains-spurious-grains)
3. [Low-Angle Grain Boundaries (LAGBs)](#3-low-angle-grain-boundaries-lagbs)
4. [Freckles](#4-freckles)
5. [Solidification Porosity](#5-solidification-porosity)
6. [Hot Tearing and Cracking](#6-hot-tearing-and-cracking)
7. [Misorientation and Off-Axis Growth](#7-misorientation-and-off-axis-growth)
8. [Ceramic Inclusions and Core Breakout](#8-ceramic-inclusions-and-core-breakout)
9. [Surface Defects and Re-Cast Layer](#9-surface-defects-and-re-cast-layer)
10. [Inspection Methods — Detecting Defects in SX Blades](#10-inspection-methods--detecting-defects-in-sx-blades)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Overview — Defect Categories in SX Castings

SX casting defects fall into two categories:
- **Microstructural defects:** stray grains, LAGBs, freckles, porosity (affect mechanical properties directly)
- **Geometric defects:** misorientation, core shift, surface roughness, dimensional non-conformance

**Defect impact on blade life:**

| Defect | Effect on fatigue | Effect on creep | Detectability |
|--------|-----------------|----------------|---------------|
| Stray grain | Severe (grain boundary = crack initiation site) | Moderate | Laue X-ray (surface); UT (internal) |
| LAGB | Moderate (depends on angle) | Moderate | EBSD; synchrotron Laue |
| Freckle | Severe (compositional variation, low m.p. phases) | Moderate | X-ray fluoroscopy, sectioning |
| Shrinkage pore | Severe (stress concentrator) | Moderate | X-ray RT, FPI on surface |
| Gas pore | Moderate | Minor | X-ray RT |
| Hot tear | Severe (pre-existing crack) | — | FPI, X-ray |
| Misorientation > 10° | Moderate (increased thermal stress) | Moderate | Laue |
| Ceramic inclusion | Severe (hard particle, stress concentrator) | Minor | X-ray RT |

---

## 2. Stray Grains (Spurious Grains)

**Definition:** An equiaxed grain with random orientation nucleating ahead of the advancing SX front during solidification.

**Why stray grains form:**
Three mechanisms:
1. **Thermal perturbation:** Local undercooling ahead of the DS/SX front → nucleation of new grains (especially at geometric features: platform corners, firtree roots)
2. **Ceramic fragment:** A fragment of the ceramic shell mold or ceramic core breaks off → acts as nucleant for a new grain in the melt
3. **Thermosolutal convection:** At high thermal gradients, convective currents in the melt carry solute-depleted liquid ahead of the front → local constitutional supercooling → nucleation

**Where stray grains prefer to form:**
```
Turbine blade cross-section:

Leading edge: HIGH RISK
   → acute angle → large undercooled zone
   → ceramic core nearby

Platform area: MODERATE RISK
   → abrupt change in cross-section → local hot spot

Trailing edge: MODERATE RISK
   → thin, rapid solidification

Airfoil body: LOW RISK
   → controlled Bridgman growth
```

**Consequences:**
- A single stray grain creates a HIGH-ANGLE grain boundary (HAGB) across part of the blade → stress concentrator → fatigue crack initiation → 2–4× reduction in fatigue life

**Prevention:**
- Optimize G/R ratio (high G, low R → avoid undercooled zone ahead of front)
- Use grain selector (spiral or pigtail: only one grain survives the tortuous path)
- Minimize G-R variations at platform/leading edge by mold geometry design
- Reduce ceramic mold roughness (smoother surface → fewer nucleation sites)
- Use helium backfill gas during solidification → higher thermal conductivity → higher G in melt

**Detection:**
- Laue X-ray: double spot pattern (if stray grain is surface-accessible by X-rays)
- Fluorescent penetrant + gentle etch: grain boundaries visible on blade surface
- Sectioning (destructive): used for process qualification, not production

---

## 3. Low-Angle Grain Boundaries (LAGBs)

**Definition:** A boundary between two regions of the SAME crystal with slightly different orientations (< ~15°). Built from arrays of dislocations.

**How LAGBs form in SX:**
1. **Local cooling rate variation:** Thermal gradients across the blade cross-section are not perfectly uniform → different parts of the crystal solidify at slightly different orientations
2. **Post-solidification plastic deformation:** During cooling (while the blade contracts), stresses can cause local plastic deformation → dislocation substructure → sub-grain boundaries
3. **Thermomechanical interactions with mold:** The ceramic mold constrains the blade as it contracts → differential stress → LAGBs near the surface

**LAGB misorientation ranges:**
- < 1°: essentially invisible; no significant mechanical effect
- 1–5°: detectable by EBSD; slight fatigue life reduction
- 5–15°: clearly a LAGB; fatigue life moderately reduced; creep: crack follows LAGB
- > 15°: treated as a HAGB → automatic reject

**LAGB detection:**
- Laue diffraction: misorientation < 2–3° is too small to detect from spot splitting → LAGB can PASS Laue and be missed
- EBSD: detects misorientation ≥ 0.1° → gold standard for LAGB detection, but requires polished cross-section → semi-destructive (requires sectioning)
- Synchrotron white beam Laue (research): non-destructive 3D LAGB mapping through blade volume

**Specification:**
- Most OEM specs: LAGBs > 5° misorientation require rejection or disposition per engineering review
- Some specs: LAGBs > 2° if > 5 mm long in critical airfoil section → rejection

---

## 4. Freckles

**Freckles** are rows of fine equiaxed grains aligned along the casting direction, with slightly different composition from the surrounding crystal.

**Formation mechanism (thermosolutal convection):**
1. During solidification, interdendritic liquid is enriched in heavy, low-freezing-point elements (Re, W → k < 1)
2. This heavy, enriched liquid is denser than the surrounding bulk → tends to sink
3. BUT: this liquid is also hotter (lower solidification T) → lower viscosity → convective instability
4. If Rayleigh number Ra > critical value → convective fingers develop → "chimneys" of interdendritic liquid
5. Bulk liquid flows DOWN, fresh liquid flows UP through chimneys → channels in the solid → freckling

**Freckle composition:**
- Higher in: Al, Ti (elements with k > 1 in interdendritic; wait — actually freckles are enriched in low-k elements for DS Ni alloys)
- Lower in: Cr, Co (elements with k > 1 that deplete in interdendritic)
- Contains eutectic phases (γ′ + γ eutectic) and TCP phases if composition is extreme

**Freckle morphology:**
```
Blade cross-section:

   ┌──────────────────┐
   │         ●        │   ← Freckle streak (aligned along DS direction)
   │    ●              │   consists of equiaxed grains with off-composition
   │●                 │
   └──────────────────┘
   
Appear as dark tracks on etched cross-section;
visible on X-ray fluoroscopy as density variation
```

**Freckle susceptibility factor:**
Empirical prediction: Freckling index FI = G^(-1.5) × (ΔρΔT) where Δρ = density difference, ΔT = solidification range.

**Prevention:**
- Higher thermal gradient G → suppresses convection (higher G → shorter diffusion length → solidification more stable)
- Faster withdrawal rate R → less time for convection to develop (but R must balance with G for good SX)
- Alloy composition: reducing the density inversion (Re reduction, or higher Co which partitions more evenly)
- Taller melt column is worse (more buoyancy-driven flow)

**Consequences:**
- Freckles have equiaxed grain microstructure → HAGB interfaces → fatigue crack initiation
- Off-composition regions → different γ/γ′ phase balance → different creep properties → local stress concentration
- Automatic rejection criterion at most OEMs

---

## 5. Solidification Porosity

Two types (same as in conventional castings, Ch 33, but vacuum casting under controlled conditions):

**Shrinkage porosity:**
- Location: last-to-freeze regions (interdendritic spaces, areas remote from feeding)
- Morphology: irregular, interconnected channels
- Cause: liquid metal contracts 3–5% on solidification → insufficient feeding
- Prevention: proper riser and gate design; controlled withdrawal rate; vacuum prevents gas interference

**Gas porosity:**
- Spherical pores
- Cause: dissolved gases (O₂, N₂, H₂) come out of solution during solidification
- In VIM casting under vacuum: O₂ and N₂ are removed by the vacuum → gas porosity is much less common than in air-cast alloys
- H₂ from moisture in ceramic mold can still cause small spherical pores

**Pore size and location matter:**
- Surface pore: more dangerous (can be decorated with oxide → harder to close by HIP)
- Internal pore: can be closed by HIP (hot isostatic pressing) if surface is sealed
- HIP closes pores by creep at high pressure → reduces porosity to < 0.05 vol%
- Standard practice: HIP all cast SX blades (ASTM specification level 2 → < 0.5% porosity)

---

## 6. Hot Tearing and Cracking

**Hot tears** form during cooling when the semi-solid blade contracts against the ceramic mold:
- Same mechanism as in conventional castings (Ch 08): strain during mushy zone → crack along solidification front
- SX alloys are more susceptible than conventional alloys because: no grain boundaries → less ductility in the semi-solid state (grain boundaries can slide in conventional; SX cannot)

**Locations:**
- Platform edges (geometric stress concentration)
- Firtree root features (blade attachment to disk)
- Thin-to-thick transitions

**Detection:**
- FPI (fluorescent penetrant): surfaces only
- X-ray: sometimes visible if enough separation
- Sectioning: destructive cross-section reveals internal tears

**Prevention:**
- Mold compliance (use ceramic with lower stiffness near stress concentration zones)
- Controlled temperature profile around blade during cooling
- Modified firtree geometry (stress-optimized through finite element analysis)

---

## 7. Misorientation and Off-Axis Growth

**Primary misorientation > 10° from [001]:**
- Laue measurement: single spot pattern but shifted
- Rejection criterion: > 10° primary OR > 10° secondary deviation

**Causes:**
- Grain selector design: double helix, pigtail, or spiral selector may not consistently select [001]
- Casting tilt: if mold is slightly tilted relative to gravity during Bridgman withdrawal, growth direction tilts
- Hot zone asymmetry: if thermal field in the Bridgman furnace hot zone is asymmetric, growth direction deviates

**Off-axis growth (twist and tilt):**
Even within a nominally SX blade, the crystal axis may slightly twist (rotate about the growth direction) or tilt over the blade length. This is mapped by multi-point Laue measurements.

---

## 8. Ceramic Inclusions and Core Breakout

**Ceramic shell inclusions:**
During pouring, ceramic shell fragments can detach and be incorporated into the casting:
- Zirconia (ZrO₂) inclusions: visible by X-ray (high density → bright spots)
- Alumina (Al₂O₃): less visible by X-ray
- Effect: hard ceramic particle → stress concentrator → fatigue crack initiation site

**Core shift and core breakout:**
The ceramic core that defines internal cooling channels must be held precisely in position:
- Core shift: ceramic core moves during pouring → cooling channel wall thickness non-uniform → over-thin wall → overheating in service
- Core breakout: part of the ceramic core fractures and shifts → cooling channel geometry changes
- Detection: X-ray fluoroscopy of finished blade (shadow shows channel geometry)

**Core-to-wall tolerance:** 0.25 mm minimum wall thickness → if core shifts 0.3 mm → wall thickness 0 mm → rejection.

---

## 9. Surface Defects and Re-Cast Layer

**Oxidized surface:** If the blade is briefly exposed to air during handling → thin oxide scale → must be removed

**Re-cast layer from EDM drilling of cooling holes:**
(Covered in detail in Ch 37) — recast layer = automatic cleaning by HF acid or electrochemical treatment per specification.

**Ceramic wash (dressing) on airfoil:**
After investment casting, ceramic shell residue must be completely removed by shot blast + acid leach. Incomplete removal → ceramic bits on surface → FPI false indication or actual inclusion.

---

## 10. Inspection Methods — Detecting Defects in SX Blades

**Standard inspection sequence for SX HPT blades (aerospace):**

```
Step 1: Visual inspection (100%)
  → Surface finish, geometry, surface tears, ceramic residue

Step 2: X-ray radiography (100%)
  → Porosity, ceramic inclusions, core geometry (multiple orientations)

Step 3: Laue diffraction (100%)
  → Primary + secondary orientation, stray grains

Step 4: Fluorescent penetrant inspection (FPI) (100%)
  → Surface cracks, tears, open pores at surface

Step 5: Coordinate measuring machine (CMM) (100%)
  → Dimensional compliance (airfoil profile, wall thickness)

Step 6: Flow test (100%)
  → Cooling air flow through internal passages (verifies channel geometry)

Step 7: (Optional) HIP treatment → then repeat X-ray (for porosity closure verification)

Step 8: Chemical analysis coupon (periodic) → confirm alloy composition
Step 9: Mechanical property coupon (periodic) → confirm tensile/creep properties
```

**Statistical process control:**
OEMs track defect rates by: defect type, position in blade, furnace/mold lot, alloy heat → identify root cause trends.

---

## Summary

| Defect | Root Cause | Life Impact | Detection | Prevention |
|--------|-----------|------------|-----------|-----------|
| Stray grain | Nucleation ahead of SX front | Fatal (fatigue) | Laue, FPI | Optimize G/R, grain selector |
| LAGB > 5° | Thermal/stress during solidification | Moderate | EBSD | Uniform thermal field |
| Freckle | Thermosolutal convection | Severe | X-ray, sectioning | Higher G, alloy comp |
| Shrinkage pore | Insufficient feeding | Moderate-severe | X-ray RT | Riser design, HIP |
| Hot tear | Constraint during solidification | Fatal | FPI, X-ray | Mold compliance |
| Misorientation > 10° | Grain selector, furnace asymmetry | Moderate | Laue | Improved selector, furnace calibration |
| Ceramic inclusion | Shell fragment | Severe | X-ray RT | Shell handling, mold inspection |

---

## Exercises

1. During process qualification of a new SX casting furnace, 50 blades are cast and inspected. Results: 3 blades have stray grains, 2 have freckles, 4 have primary misorientation > 10°, 5 have porosity > level 2 (before HIP), 1 has a hot tear. (a) Calculate overall casting yield (fraction passing all criteria before HIP). (b) After HIP, the 5 porosity-failed blades are re-inspected: 4 pass (pores closed), 1 fails (surface-connected pore not closed). Calculate yield after HIP. (c) To improve yield, propose which defect to target first (highest impact). For stray grains: describe two specific casting parameter changes to reduce stray grain frequency. (d) The foundry estimates: each percentage point of yield improvement saves $500,000/year (on 10,000 blades/year at $800 each). Calculate the annual saving if stray grain frequency drops from 6% to 2%.

2. Freckle analysis in a DS Ni superalloy casting: (a) In unidirectional solidification, the critical Rayleigh number for freckling is Ra_c ≈ 0.1. Ra = (Δρ×g×K×h) / (μ×κ) where Δρ = 200 kg/m³ (density difference between interdendritic and bulk liquid), g = 9.8 m/s², K = 10⁻¹¹ m² (permeability of mushy zone), h = 0.05 m (mushy zone height), μ = 0.005 Pa·s (viscosity), κ = 5×10⁻⁶ m²/s (thermal diffusivity). Calculate Ra. Is freckling predicted? (b) To reduce Ra below Ra_c, what parameter would you change? By what factor must K change to bring Ra to Ra_c (keeping all other parameters constant)? (c) K ∝ λ₂² (secondary dendrite arm spacing squared). How much must λ₂ decrease to achieve the required K reduction? What casting parameter controls λ₂?

3. LAGB formation and impact on fatigue: (a) A LAGB consists of edge dislocations with Burgers vector b = 0.25 nm, spaced by D = b/θ (θ = misorientation angle in radians). Calculate dislocation spacing for θ = 2°, 5°, 10°. (b) In fatigue, the stress intensity at a grain boundary is approximately K_I = σ√(πd/2) where d is the grain (or sub-grain) size. For a LAGB at 2° misorientation, d = 100 μm; for 5°, d = 40 μm; for 10°, d = 20 μm. Calculate K_I for each at σ = 400 MPa. Compare to K_Ic = 45 MPa√m. (c) If K_I > K_Ic, the LAGB cracks immediately. If K_I > 0.6×K_Ic, fatigue crack initiation is likely early. Which misorientation angles cause concern? (d) Industry data shows: LAGB < 2°: fatigue life = 1.0×; 2–5°: 0.7×; 5–10°: 0.4×; > 10°: 0.1×. Calculate the fleet-wide impact if 5% of blades have LAGBs in the 5–10° range.

4. X-ray radiography for SX blade inspection: The specification allows maximum pore diameter 0.5 mm in the airfoil. (a) Standard X-ray RT has sensitivity of ~1% of material thickness to density changes. For a blade wall thickness of 4 mm, what is the minimum density change detectable? (b) A spherical pore of 0.3 mm diameter in 4 mm steel: what percentage of the wall thickness is the pore diameter? Is it above the 1% sensitivity limit? (c) To detect a 0.3 mm pore more reliably, what technique would you use? (d) X-ray CT can detect pores ≥ 50 μm. Calculate how many CT scan hours would be needed for a fleet of 2,000 blades × 80 HPT blades/engine if each CT scan takes 4 hours. Is 100% CT inspection economically feasible? At what critical component level would you use it?

5. Economic analysis of casting defects: OEM produces 15,000 SX HPT blades per year. Material cost per blade: $1,200 (CMSX-4 + investment). Processing cost per blade: $600 (casting, HIP, heat treatment). Inspection cost: $400 per blade (all NDT). Current defect rates and consequences: stray grain 4% → scrap ($2,200 lost). Misorientation 3% → rework ($150) or scrap (50% of these, $2,200). Porosity 5% → 80% closed by HIP (additional $200 HIP cost), 20% scrapped. Hot tear 1% → scrap. (a) Calculate total annual scrap/rework cost. (b) The foundry invests $2M in new helium-injection Bridgman furnaces that reduce stray grain to 1%. Calculate annual savings from reduced stray grain rejection. What is the payback period on the $2M investment? (c) Alternatively, they invest $500K in improved grain selector design that reduces misorientation > 10% to 0.5%. Calculate annual savings. Which investment has better ROI?
