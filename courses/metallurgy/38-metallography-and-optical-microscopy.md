# Chapter 38: Metallography and Optical Microscopy — Seeing the Microstructure

> **"All the theory about dislocations, grain boundaries, precipitates, and phases is just words until you look at them under a microscope. Metallography is the art of preparing a metal surface so perfectly that features as small as 1 micron — one thousandth of a millimeter — can be seen clearly. It has been practiced since 1863 when Henry Clifton Sorby first polished and etched steel, looked under a microscope, and saw pearlite."**

---

## Table of Contents

1. [Why Metallography?](#1-why-metallography)
2. [Sectioning and Mounting](#2-sectioning-and-mounting)
3. [Grinding and Polishing](#3-grinding-and-polishing)
4. [Etching](#4-etching)
5. [Optical Microscope — Principles and Operation](#5-optical-microscope--principles-and-operation)
6. [What You See: Interpreting Microstructures](#6-what-you-see-interpreting-microstructures)
7. [Quantitative Metallography](#7-quantitative-metallography)
8. [Hardness Testing in Context](#8-hardness-testing-in-context)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Metallography?

Metallography prepares metal samples so microstructural features can be observed and measured. It is the fundamental characterization tool connecting:
- **Processing → Microstructure:** Did the heat treatment produce the right phase?
- **Microstructure → Properties:** Is the grain size meeting specification?
- **Failure analysis:** What microstructure was present at the fracture origin?

Every alloy specification (AMS, ASTM, MIL-SPEC) for aerospace metals includes microstructural requirements: grain size, phase fractions, heat-affected zone, precipitate distribution. Metallography is the measurement tool that verifies these.

---

## 2. Sectioning and Mounting

**Step 1 — Select cutting plane:**
The plane of the section determines what you see. For a forging, cut both parallel and perpendicular to the working direction to reveal grain flow.

**Step 2 — Sectioning (abrasive cut-off):**
- Abrasive cut-off wheel (SiC or Al₂O₃) with flood coolant
- MUST use coolant: prevents heating → prevents microstructure change
- Sectioning force should be minimal to avoid mechanical deformation of surface

**Step 3 — Mounting:**
Small samples need to be mounted in a polymer for easier handling:
- **Hot mount (Bakelite, Duroplast):** Sample + powder in mount press → 150°C, 150 MPa → 10 min → solid mount
  - Requires heat + pressure → only for materials unaffected by 150°C
  - Not suitable for: Al alloys with T6/T7 temper (will over-age), PM samples (if pores matter)
- **Cold mount (epoxy, acrylic):** Mix two-component resin → pour over sample → cure at RT (hours)
  - No heat → safe for all materials
  - Epoxy infiltrates pores → essential for porous materials (PM parts, thermal spray coatings)

**Edge retention:** The sample edge must remain sharp (no rounding) for:
- Surface coatings (TBC, CVD, PVD)
- Carburized/nitrided cases
- Weld HAZ at fusion line

Electroplating a Ni layer on the surface BEFORE mounting → helps maintain edge sharpness.

---

## 3. Grinding and Polishing

**Goal:** Produce a scratch-free, flat, mirror-polished surface for etching.

**Grinding sequence:**
```
Coarse SiC paper (120–240 grit) → removes sectioning damage
  → rotate 90° each step
  → Medium SiC (320–400 grit) → finer scratches
  → Fine SiC (600–800–1200 grit) → near-scratch free
```
After each step: rotate sample 90° to ensure previous scratches are removed.

**Polishing sequence:**
```
3 μm diamond paste on cloth → removes SiC scratches
  → 1 μm diamond paste → near-mirror
  → 0.3 μm Al₂O₃ or SiO₂ suspension → mirror polish
  → Final: 0.05 μm colloidal silica (OP-S) → true mirror, also light chemical attack
```

**Common polishing artifacts to avoid:**
- **Relief:** Hard phases stand above soft matrix → height difference → scratches in softer phase
- **Scratches:** Usually from contaminated cloth or previous step not completed
- **Pull-out:** Inclusions or second phases fall out → appear as holes (can be mistaken for pores)
- **Smearing:** Soft phases (Pb in brass, MnS in steel) smear during polishing → don't show true morphology
- **Rounding of edges:** Excessive polishing time rounds the mount edges → coating/layer appears thicker than it is

**Automated polishing:** Computer-controlled force, speed, time → reproducible → required for quantitative work.

---

## 4. Etching

Etching reveals microstructure by selectively attacking the polished surface. The polished surface shows nothing (only reflections).

**How etching works:**
```
Un-etched → mirror → no contrast → no microstructure visible
Etched:
  - Grain boundaries: attacked more rapidly (high energy) → appear as dark lines (grooves scatter light)
  - Different phases: different dissolution rates → height differences → contrast
  - Crystallographic orientation: different rates on different planes → same phase, different grain = different gray level
```

**Common etchants:**

| Etchant | Composition | Used for |
|---------|------------|---------|
| Nital (2%) | 2 mL HNO₃ + 98 mL ethanol | Carbon steels, alloy steels |
| Nital (5%) | 5 mL HNO₃ + 95 mL ethanol | Stainless, tool steels |
| Aqua regia | 3 HCl : 1 HNO₃ | Nickel alloys, Cu alloys |
| Keller's | HF + HNO₃ + HCl in water | Aluminum alloys |
| Marble's | CuSO₄ + HCl + H₂O | Austenitic stainless, Cu alloys |
| Murakami's | KOH + K₃Fe(CN)₆ | Carbides, WC-Co |
| Weck's | 4g NH₄HF₂ + 0.5g K₂S₂O₅ in 100mL H₂O | Aluminum alloys (color etching) |

**Color etching (tint etching):**
Forms thin oxide films of varying thickness on different grains → interference colors → different orientations appear different colors (Beraha's reagent). Useful for: stainless steel grain visualization, Ti alloys, Al alloys.

**Electrolytic etching:**
Pass current through sample in etchant → controlled, reproducible → essential for:
- Austenitic stainless steel (oxalic acid electrolytic → reveals sensitization)
- Ni superalloys (phosphoric acid → reveals γ/γ′ contrast)

---

## 5. Optical Microscope — Principles and Operation

**Magnification range:** 50× to 1,000× (maximum useful = 1,000× due to visible light wavelength limit)

**Resolution limit:**
```
d_min = 0.61 × λ / NA

where λ = wavelength of light (0.4–0.7 μm visible)
      NA = numerical aperture of objective (max ~1.4 with oil immersion)

d_min = 0.61 × 0.55 μm / 1.4 ≈ 0.24 μm (240 nm)
```

At 1,000× magnification, this resolution corresponds to features ≥ 0.24 μm being distinguishable — sufficient for grains, phases, larger precipitates.

**Illumination modes:**

**Bright-field (BF):** Normal mode — light reflects off flat surface → bright; scratches, grain boundaries (depressed) → appear dark.

**Dark-field (DF):** Only scattered/diffracted light collected → flat surface = dark; scratches/edges = bright. Used for: detecting fine cracks, second-phase particles.

**Polarized light:** Differentiates anisotropic phases (different crystal structures):
- α-Ti (HCP): anisotropic → different orientations show different colors/brightness
- β-Ti (BCC): isotropic → doesn't respond to polarized light
- Extremely useful for: Ti alloys, Al alloys, graphite in CI, uranium oxides

**Differential Interference Contrast (DIC/Nomarski):** Creates 3D surface topography effect → reveals height differences after etching → excellent for: carburized layers, martensitic structures, coating cross-sections.

---

## 6. What You See: Interpreting Microstructures

**Steel microstructures (after nital etch, optical microscope):**

| Microstructure | Appearance | Notes |
|---------------|-----------|-------|
| Ferrite | Light (white) grains | BCC, low C |
| Cementite (Fe₃C) | White boundaries or plates | Very hard |
| Pearlite | Dark lamellar structure | Alternating ferrite + Fe₃C |
| Bainite | Dark acicular (needle-like) | Upper/lower bainite differ |
| Martensite | Dark lath or plate | Untempered = very hard |
| Tempered martensite | Dark, fine carbide + matrix | After tempering |
| Austenite (retained) | White islands | LOM can't easily distinguish from ferrite |
| Carbide | White (primary or alloy carbide) | Different etchants to distinguish |

**Nickel superalloy (γ/γ′, etched with aqua regia or Kalling's):**
- Matrix (γ): appears etched, darker
- Cuboidal γ′: light-colored cubes arranged in matrix
- GB carbides: discrete particles at grain boundaries
- TCP phases (if present): plate-like acicular features

**Aluminum alloys (Keller's etch):**
- Grain boundaries: revealed as dark lines
- MgSi₂ (β phase in 6xxx): dark striped/fibers precipitated within grain
- CuAl₂ (θ in 2xxx): white islands
- PFZ (precipitate-free zones) at grain boundaries: visible as brighter band

---

## 7. Quantitative Metallography

**Grain size measurement (ASTM E112):**

**Intercept method:**
```
Draw a test line of known length L across the micrograph
Count intercepts N (each grain boundary crossing = 1 intercept)
Mean intercept length = L / N
```

**Planimetric method (Jefferies):**
Count grains (N_inside + 0.5×N_on line) in known area → number of grains per mm²

**ASTM grain size number (G):**
```
N_A = 2^(G-1) grains per mm² at 100×

Alternatively: d(μm) = 15 × 2^(-G/2) (approximate)
```
| G | d (μm) | Notes |
|---|--------|-------|
| 1 | 240 | Very coarse |
| 4 | 84 | Coarse |
| 6 | 42 | Medium |
| 8 | 21 | Fine |
| 10 | 10 | Very fine |
| 12 | 5 | Ultrafine |

**Phase fraction measurement:**
- Point count method: superimpose grid; count hits on each phase → fraction = hits/total × 100%
- Lineal analysis: total length of test line in each phase / total length
- Image analysis software: threshold by gray level → automatic fraction

---

## 8. Hardness Testing in Context

Hardness testing is often performed alongside metallographic examination:

**Vickers hardness (HV):**
```
F (kgf applied to diamond pyramid)
HV = 0.1891 × F / d²   (d = mean diagonal, mm)
```
| Scale | Load range | Used for |
|-------|-----------|---------|
| HV0.001–HV0.05 | < 50g | Microhardness (individual phase, HAZ mapping) |
| HV1–HV10 | 1–10 kg | Layer hardness, small samples |
| HV30–HV100 | 30–100 kg | Bulk hardness |

**Microhardness traverses:**
- Measure HV every 0.1–0.5 mm across a section → hardness profile
- Case depth determination: carburized case depth = distance from surface to 550 HV (or 50 HRC)
- Weld HAZ hardness mapping: identify hard martensite zones (>380 HV → cold cracking risk)

**Nano-indentation (instrumented indentation):**
- Load-displacement curve → hardness + elastic modulus of individual phases
- Can measure: γ′ hardness, TBC hardness, individual grain orientation effects

---

## Summary

| Step | Purpose | Key Variables |
|------|---------|--------------|
| Sectioning | Select representative plane | Section orientation vs. microstructure |
| Mounting | Handle small samples; edge retention | Cold vs hot mount; epoxy for pores |
| Grinding | Remove damage; flatten | SiC grit sequence; rotation |
| Polishing | Mirror finish | Diamond size; OP-S final |
| Etching | Reveal microstructure | Etchant selection; time; temperature |
| Microscopy | Observe + measure | Mode (BF/DF/Pol/DIC), magnification |
| Quantitative | Grain size, phase fraction | ASTM E112; point count |

---

## Exercises

1. A carburized and hardened 8620 steel (0.20%C) gear tooth sample needs to be examined for: (a) case depth, (b) grain size in the case, (c) martensite morphology, (d) retained austenite estimate. Specify: (i) sectioning plane, (ii) mounting material (why hot or cold mount?), (iii) polishing sequence from rough grinding to final polish, (iv) etchant and concentration, (v) optical mode (BF or polarized?) for each measurement. For case depth: if the hardness profile shows 720 HV at the surface → 550 HV at 1.2 mm → 380 HV at 2.5 mm → bulk 280 HV; which depth is specified as the case depth?

2. An IN718 turbine disk (nickel superalloy) requires grain size verification: specified ASTM 8–10. After solution + age heat treatment, the sample is mounted, polished, and etched with electrolytic 10% H₂CrO₄ (chromic acid). At 200× (scale bar 100 μm), a test line of 50 mm crosses 28 grain boundaries. (a) Calculate the mean intercept length in actual mm. (b) Calculate the grain size in μm. (c) Convert to ASTM G number. (d) Does it meet the ASTM 8–10 specification? (e) The specification also requires no individual grain larger than ASTM 5. How would you measure the maximum grain size? Is the ASTM planimetric method appropriate here or should you use an image analysis approach?

3. A welded 304 stainless steel pipe (welded with 308L filler) is tested for sensitization using the ASTM A262 Practice A (oxalic acid electrolytic etch). The microstructure shows: (a) "step" structure in base metal = acceptable, (b) "ditch" structure 3 mm from fusion line = possible sensitization, (c) "dual" structure at fusion line = acceptable. What does each structure indicate? For the "ditch" structure zone: which temperature range caused sensitization during welding? (use Cr depletion model: Ch₂₃C₆ forms at 450–850°C). (c) If the pipe was welded with 316 (not 316L), calculate the maximum time it can be held at 700°C before sensitization occurs (use: C = 0.05%, carbide formula C + Cr in 7:1 ratio by weight, sensitization when Cr drops below 12% locally, initial Cr = 18%).

4. Quantitative metallography of a ductile cast iron (Grade 65-45-12, ASTM A536): Perform the following measurements at 100×: (a) Graphite nodule count (N per mm²): 22 nodules observed in a 500×500 μm² field. Calculate N_A. (b) Graphite area fraction: point count with 100 points → 14 land on graphite. Calculate f_graphite (%). (c) Nodularity: 19 of 22 nodules appear spherical; 3 are elongated. Calculate nodularity (%). (d) Pearlite fraction in matrix: 30 points on pearlite, 56 on ferrite out of 100 total. What are the specifications per ASTM A536 for Grade 65-45-12? Is this casting acceptable? (e) If nodularity drops below 85%, what property change occurs? (Hint: compare graphite flakes in gray iron vs nodules — revisit Ch 23).

5. Color metallography of duplex stainless 2205 using Beraha's tint etch: the microstructure shows dark (austenite, γ) and light (ferrite, α) phases. Image analysis of a 500×500 μm² field gives: ferrite area = 48,000 μm², austenite = 52,000 μm². (a) Calculate ferrite fraction. Specified range: 40–60% ferrite. Is this acceptable? (b) The ferrite grains have mean intercept length 12 μm; austenite islands 8 μm. Calculate ASTM G for each. (c) After improper heat treatment at 900°C for 2 hours, the microstructure shows a third white phase at ferrite/austenite grain boundaries. What phase is this? What is the consequence for: (i) corrosion resistance, (ii) toughness? (d) What etchant and microscope mode would you use to confirm the identity of this third phase?
