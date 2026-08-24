# Chapter 34: Forming Processes — Shaping Metal by Deformation

> **"Rolling, forging, extrusion, drawing — these are the primary vocabulary of industrial metalworking. The same piece of aluminum, depending on whether you roll it, extrude it, or forge it, will have different grain structures, different mechanical properties, and different preferred orientations. Processing is microstructure; microstructure is properties. Understanding forming processes is understanding how we sculpt metal at the atomic scale."**

---

## Table of Contents

1. [Overview — Why Forming?](#1-overview--why-forming)
2. [Hot vs. Cold Working](#2-hot-vs-cold-working)
3. [Rolling](#3-rolling)
4. [Forging](#4-forging)
5. [Extrusion](#5-extrusion)
6. [Drawing (Wire, Rod, Tube)](#6-drawing-wire-rod-tube)
7. [Sheet Metal Forming (Stamping, Deep Drawing)](#7-sheet-metal-forming-stamping-deep-drawing)
8. [Formability and the Forming Limit Diagram](#8-formability-and-the-forming-limit-diagram)
9. [Texture and Anisotropy from Forming](#9-texture-and-anisotropy-from-forming)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Overview — Why Forming?

Forming processes shape metal by plastic deformation — WITHOUT removing material (unlike machining). Advantages:
- Low material waste (vs. machining)
- Microstructure improvement: wrought structure better than cast (eliminated casting porosity, refined grain, broken up dendritic segregation)
- Work hardening (cold work) can strengthen the material
- Grain flow can be aligned for optimal fatigue resistance (forged vs. machined thread root in a bolt)
- Very high production rates possible (rolling mills run at km/hour)

---

## 2. Hot vs. Cold Working

**Hot working:** T > 0.6 T_m — recrystallization occurs simultaneously with deformation:
- No work hardening (soft product)
- Large deformations possible
- Low flow stress (less force required)
- Oxidation (controlled atmosphere often needed)
- Decarburization risk for steels

**Cold working:** T < 0.3 T_m — recrystallization does NOT occur:
- Work hardens → increased yield strength
- Better surface finish
- Close dimensional tolerance
- Requires more force than hot working
- Ductility decreases (must anneal between passes for large total deformation)

**Warm working:** T = 0.3–0.6 T_m — intermediate properties:
- Some recovery (partial softening)
- Better surface than hot, less force than cold
- Used for: titanium (warm forging at 450°C), steel (warm rolling for HSLA)

---

## 3. Rolling

**Rolling** reduces thickness by passing metal between rollers:

```
Rolling geometry:

    ←───────── Δh = h₀ - h₁ ─────────→
    
    ○─────────────────────────────────○
    │  h₀  │→  ←  →  ←  →  ←  │ h₁ │
    ○─────────────────────────────────○
    
    Reduction = (h₀ - h₁) / h₀ × 100%
    Elongation = L₁ / L₀ = h₀ / h₁ (volume conservation)
```

**Rolling Mills:**
- Two-high: simplest; reversible or non-reversible
- Four-high: work rolls (small) + backup rolls (larger, prevent deflection); plate and sheet mills
- Cluster mill (Sendzimir): very small work rolls, multiple backups; for hard thin strip
- Tandem mill: multiple stands in sequence; each reduces thickness; used for continuous production

**Hot Rolling:**
- Slabs → plate → sheet at 1,100–1,300°C for steel
- Billets → bar, rod, shapes
- Austenite grain is flattened and elongated → recrystallizes between passes if long enough pause
- TMCP (Thermomechanical Controlled Processing): roll in austenite pancake + immediate quench → very fine transformed ferrite → HSLA steel

**Cold Rolling:**
- After hot rolling + pickling (acid descale)
- Produces: automotive body sheet (0.5–2 mm), foil (0.006 mm for packaging)
- Work hardens → anneal to restore ductility → temper roll to final temper/flatness

**Ring Rolling:**
Hollow ring blank rotated between rolls → ring grows in diameter while section reduces:
- Produces: bearing rings, gear rings, flanges, jet engine disk preforms
- Seamless rings (no welded joint) → better fatigue life

---

## 4. Forging

**Forging** compresses metal between dies to produce a specific shape. Produces the best mechanical properties of any forming process (refined grain, favorable grain flow, no porosity).

**Types:**

**Open-die forging (cogging, upsetting):**
- Simple flat or shaped dies
- Shape controlled by die motion and operator skill
- Large forgings: crankshafts, shafts, rings
- Produces billet/preform for subsequent closed-die forging

**Closed-die (impression die) forging:**
- Metal fills a cavity defined by upper and lower dies
- Flash (excess metal) squirts out around periphery and is trimmed
- Net shape (near-net-shape): forging matches final part closely
- Examples: connecting rods, gears, titanium airframe bulkheads, crankshafts

**Isothermal forging:**
- Die heated to same temperature as workpiece
- No heat loss → enables more uniform deformation, thinner sections
- Used for: Ti-6Al-4V bulkheads, superalloy disks (IN718 at 950°C with nickel dies)

**Grain flow (fiber flow) in forgings:**
```
Machined bolt (from bar):      Forged bolt:
  Grain flow:                    Grain flow:
    ─ ─ ─ ─ ─ ─ ─                ─────────────
    ─ ─ ─ ─ ─ ─ ─                 ────────────
    Thread root: grain            ╰ ─ thread root: grain
    lines cut by machining →       lines FOLLOW contour →
    stress concentrators           higher fatigue strength
```

This is why critical fasteners (engine bolts, landing gear fittings) must be forged, not machined from bar.

**Forging temperatures (representative):**
| Material | Temperature (°C) |
|----------|----------------|
| Low-C steel | 900–1,200 |
| Stainless steel | 1,050–1,200 |
| Aluminum alloys | 350–475 |
| Ti-6Al-4V | 870–980 (below β-transus) |
| IN718 superalloy | 950–1,050 |

---

## 5. Extrusion

**Metal is pushed through a shaped die** to produce a constant cross-section product:

```
Extrusion:

Billet → │ Ram │→ Die → Extruded profile
          ↓ pressure
```

**Direct extrusion:** Billet moves forward; friction between billet and container wall.
**Indirect extrusion:** Die moves backward through stationary billet; no billet-wall friction → lower force.
**Hydrostatic extrusion:** Fluid pressure around billet → no billet-container friction → enables very high reductions; for brittle materials (W, Mo, Be at RT).

**What can be extruded:**
- Aluminum: complex profiles (H, T, I, tubular shapes for windows, automotive)
- Copper: rods, tubes, cables
- Steel: hot-extruded bars, hollow sections (requires refractory glass lubricant at 1,200°C)
- Plastics: most polymers extruded (not metal process but analogous)

**Extrusion ratio (ER):** A_billet / A_product
- Aluminum: ER up to 100× (flows easily)
- Steel: ER max ~10 (much harder to extrude)
- Ti: ER up to 40× (hot extrusion)

**Why extrusion is valuable:**
- Complex cross-sections in single pass (what would take multiple rolling passes)
- No parting line → no flash → near-net-shape
- Compressive stress state during deformation → excellent grain refinement + high ductility product

---

## 6. Drawing (Wire, Rod, Tube)

**Pulling metal through a die** to reduce cross-section:

```
Wire drawing:

Feed wire → │Die (converging angle)│→ Reduced wire (pulled by draw block)
```

**Wire drawing:** Extensive cold work in multiple passes; hardwire to 4,000 MPa possible (piano wire, tire cord):
- Each pass: 10–25% area reduction
- Work hardens → anneal between passes if large total reduction needed
- Final properties: high strength (cold work) or soft (annealed) depending on temper

**Tube drawing:**
- Over a mandrel (reduces wall thickness, controls bore)
- On a plug (controls both OD and ID)
- Sinking (no mandrel — OD reduced but wall can thicken)

**Rod and bar drawing:**
- Improves surface, tightens tolerances from hot rolling
- Cold-drawn steel bar: σ_y 10–20% higher than hot-rolled equivalent

---

## 7. Sheet Metal Forming (Stamping, Deep Drawing)

**Sheet metal is formed into 3D shapes** using dies and punches:

**Stamping (pressing):** Shallow, complex shapes (auto body panels):
- Blank positioned over female die
- Punch descends → metal stretches and draws over punch
- Springback: metal returns partway to original shape on removal → must overbend

**Deep drawing:** Cups, shells, cans (beverage cans are drawn aluminum):
```
  Blank (flat)
  Drawn by punch through draw ring
  Blank-holding force prevents wrinkling
  
  Maximum drawing ratio = D_blank / D_punch ≤ 2.5 (for Al, single draw)
```

**Hemming, flanging, bending:** Straightforward bending operations for brackets, flanges.

**Stretch forming:** Sheet stretched over form block for large curved panels (aircraft skins):
- Pure tensile deformation (no compressive wrinkling)
- Used for large curvature skins (wing, fuselage)

---

## 8. Formability and the Forming Limit Diagram

**Forming Limit Diagram (FLD):** Maps combinations of major strain ε₁ and minor strain ε₂ at which necking/fracture occurs:

```
FLD:

ε₁ (major strain)
  │         SAFE ZONE
0.5│─────────────────────────────────
  │         / Forming Limit Curve
0.4│        /   (FLC)
  │       /
0.3│──────────────────────────────────
  │              FAILURE ZONE
  │
  └──────────────────── ε₂ (minor strain) →
     -0.4    0    +0.4
```

**Uses of FLD:**
- Design stamp tooling to stay below FLC
- Identify critical regions (necks form at peak of FLC)
- Compare materials (Al = lower FLC than steel at same strength)

**Work hardening exponent n (from tensile test):**
```
σ = K × ε^n   (Hollomon equation)
```
Higher n → better formability (more uniform strain distribution, delays necking).

---

## 9. Texture and Anisotropy from Forming

After heavy deformation, grains rotate → preferred crystallographic orientation (texture):

**Rolling texture in FCC metals:**
- {110}⟨112⟩ and {123}⟨634⟩ (Copper/S component)
- Results in anisotropy: different strength in rolling, transverse, and through-thickness directions

**Deep drawing texture:**
- {111} planes parallel to sheet surface (ideal drawing texture — resists thinning)
- "Earing" defect: cups develop scalloped top edge from anisotropic behavior
- Controlled by: annealing to develop {111} recrystallization texture in IF steel

**Wire drawing texture:**
- ⟨110⟩ fiber texture (FCC), ⟨110⟩ + ⟨100⟩ (BCC)
- Makes drawn wire anisotropic in strength vs. transverse

**Titanium alloys:** Very strong texture from hot working:
- α-phase textures affect: anisotropy in fatigue crack growth (cracks grow faster in basal plane)
- Must control forging direction for orthopedic implants to align texture with loading direction

---

## Summary

| Process | Deformation | Grain Effect | Main Products | Key Advantage |
|---------|-------------|-------------|---------------|---------------|
| Hot rolling | Compression | Recrystallized | Plate, sheet, bar | High deformation, no work hardening |
| Cold rolling | Compression | Elongated, work hardened | Sheet, foil, strip | Close tolerance, high strength |
| Forging | Compression | Refined, grain flow | Critical structural parts | Best mechanical properties |
| Extrusion | Compression | Fine grain | Complex profiles | Complex cross-section in one pass |
| Wire drawing | Tension | Elongated, work hardened | Wire, rod | Very high strength (cold work) |
| Deep drawing | Tension+compression | Mixed | Cups, cans, automotive | Complex shapes, thin walls |

---

## Exercises

1. A steel billet (200 × 200 mm, 2 m long) is hot rolled to a plate (200 × 10 mm). (a) Calculate the reduction ratio (width stays constant, thickness changes). (b) What is the final plate length? (c) What microstructure would you expect in the hot-rolled plate (before any heat treatment)? (d) If the rolling is done using TMCP (controlled rolling + rapid cooling), describe how the final microstructure and mechanical properties differ from conventional hot rolling.

2. Compare the fatigue life of: (a) a machined 4340 steel connecting rod (cut from bar stock, grain lines run across the critical fillet radius), and (b) a closed-die forged 4340 connecting rod (grain lines follow the contour). Both have identical composition and heat treatment: Q&T to σ_y = 900 MPa. Testing shows the forged part has 2.5× the fatigue life. Explain the metallurgical reason for this difference in terms of: (i) grain boundary orientation relative to stress, (ii) work hardening benefits, (iii) surface finish.

3. Aluminum 6061 extrusion through a die with a complex T-section profile: (a) The extrusion ratio is 35 (billet area / product area). Calculate the exit velocity if the billet enters at 0.02 m/s. (b) At the extrusion temperature (520°C), the flow stress of 6061 ≈ 30 MPa. Estimate the extrusion force for a 150mm diameter billet (pressure ≈ 10 × flow stress × reduction ratio, rough estimate). (c) After extrusion, the profile is naturally aged (T4). Why doesn't 6061 achieve T6 properties through natural aging alone? (d) What temper designation results from direct aging after extrusion?

4. Deep drawing of aluminum beverage cans (3104 Al alloy, 0.3 mm sheet): (a) A blank of 150mm diameter is drawn into a cup of 65mm diameter. Calculate the drawing ratio (DR = D_blank/D_punch). Is this within the limit of 2.5 for aluminum? (b) For a two-draw process: first draw DR₁ = 2.0, second draw reduces diameter from 100mm to 65mm. What is DR₂? (c) Describe the earing problem and how crystallographic texture in the sheet causes it. (d) How is the sheet designed (annealing conditions) to minimize earing?

5. Cold drawing of steel wire: starting rod = 8 mm diameter, final wire = 2 mm diameter. (a) Calculate total area reduction, and if done in 5 passes of equal reduction each, calculate reduction per pass. (b) The starting hardness = 180 HV; final hardness = 580 HV after full cold work. What metallurgical mechanism causes this increase? (c) The wire is used as spring wire requiring σ_y > 1,500 MPa. What annealing between passes is needed (full recrystallize, or just stress relief)? (d) At what point would you expect forming limit fracture during drawing? Use tensile test data: UTS = 600 MPa on initial rod.
