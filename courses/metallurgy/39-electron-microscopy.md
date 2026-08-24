# Chapter 39: Electron Microscopy — Seeing at the Nanoscale

> **"The optical microscope let us see grains. The scanning electron microscope lets us see fracture surfaces, grain boundary films, and individual precipitates at 100,000×. The transmission electron microscope lets us see individual dislocations, stacking faults, and atomic columns. These are not just better microscopes — they are different windows into different scales of structure, each essential for understanding why metals behave the way they do."**

---

## Table of Contents

1. [Why Electron Microscopy?](#1-why-electron-microscopy)
2. [Scanning Electron Microscopy (SEM) — Principles](#2-scanning-electron-microscopy-sem--principles)
3. [Signals in SEM: SE, BSE, EDS, EBSD](#3-signals-in-sem-se-bse-eds-ebsd)
4. [SEM Sample Preparation and Applications](#4-sem-sample-preparation-and-applications)
5. [Focused Ion Beam (FIB) — Cross-Sectioning and TEM Prep](#5-focused-ion-beam-fib--cross-sectioning-and-tem-prep)
6. [Transmission Electron Microscopy (TEM) — Principles](#6-transmission-electron-microscopy-tem--principles)
7. [TEM Imaging Modes: BF, DF, HRTEM, HAADF-STEM](#7-tem-imaging-modes-bf-df-hrtem-haadf-stem)
8. [TEM Analysis: Diffraction, Phase Identification, Dislocations](#8-tem-analysis-diffraction-phase-identification-dislocations)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Electron Microscopy?

**Resolution comparison:**
```
Tool            | Resolution | What you see
─────────────────────────────────────────────────
Naked eye       | ~0.1 mm    | Grains in coarse metals
Optical (LOM)   | ~0.3 μm    | Grains, phases, cracks, pearlite
SEM             | ~1–5 nm    | Fracture surfaces, particles, thin films, EBSD
TEM             | ~0.1 nm    | Dislocations, precipitates, interfaces, atoms
```

The electron wavelength (De Broglie): λ = h/mv
- At 10 keV: λ = 0.012 nm
- At 200 keV (TEM): λ = 0.0025 nm
- Much shorter than visible light (400–700 nm) → much higher resolution

**Trade-off:** Higher resolution requires:
- Higher vacuum (electrons scatter in air)
- Thin samples (TEM: < 100 nm for electron transmission)
- Careful preparation

---

## 2. Scanning Electron Microscopy (SEM) — Principles

**How SEM works:**
```
Electron gun (thermionic or field emission)
  ↓ high voltage (1–30 kV)
Electromagnetic lenses (focus beam)
  ↓ fine spot (1–5 nm for FE-SEM, 5–50 nm for thermionic)
Scanning coils (raster the beam across sample)
  ↓
Sample generates signals
  ↓
Detector converts signal to brightness → image on screen
```

**Key parameters:**
- Accelerating voltage (kV): higher → deeper interaction volume; lower → better surface sensitivity + less damage to soft materials
- Working distance (WD): distance from final lens to sample; short WD → high resolution; long WD → more depth of field (useful for rough surfaces)
- Beam current: higher → more signal but broader beam → compromise

**Interaction volume (pear-shaped):**
Electrons enter sample → scatter → generate signals from varying depths:
```
Surface
│←── BSE (~0.1–1 μm depth)
│
│←────────── X-rays (EDS) (~0.5–3 μm depth)
│
│←──────────────── Cathodoluminescence (deepest)
```

Secondary electrons (SE): from top ~5–10 nm → high-resolution surface topography.

---

## 3. Signals in SEM: SE, BSE, EDS, EBSD

### Secondary Electrons (SE)
- Low energy (< 50 eV) electrons ejected from sample atoms by inelastic scattering
- Come from very near surface (≤ 5 nm) → sharp topographic image
- **SE image:** bright edges + tops of ridges (more SE escape) → 3D appearance
- **Used for:** Fracture surface analysis, surface morphology, general imaging

### Backscattered Electrons (BSE)
- High energy electrons elastically reflected back from sample
- Come from larger depth (0.1–1 μm) → more blurry edges
- **Atomic number contrast:** heavier atoms → more BSE → brighter
- **BSE image:** compositional mapping (heavy phases = bright, light phases = dark)
- **Used for:** Phase identification, segregation, multi-phase microstructure overview

```
Example: In γ/γ′ nickel superalloy on BSE:
  γ matrix (Ni, Cr, Co): lower Z average → darker
  γ′ (Ni₃Al, with Cr, Co, Ta, W): varies — often brighter due to W, Ta
  Carbides (NbC, TaC): very bright (Nb Z=41, Ta Z=73)
```

### Energy-Dispersive X-ray Spectroscopy (EDS / EDX)
- Electron beam ionizes atoms → characteristic X-rays emitted
- Each element emits X-rays at characteristic energies (K-lines, L-lines, M-lines)
- EDS detector measures energy spectrum → identify + quantify elements

```
EDS process:
Electron ejects inner-shell electron
→ Outer electron drops to fill vacancy
→ Characteristic X-ray emitted (energy = ΔE between shells)
Fe Kα = 6.40 keV; Ni Kα = 7.47 keV; Cr Kα = 5.41 keV
```

**EDS capabilities:**
- Detect elements Z ≥ 4 (Be limit)
- Qualitative analysis: always possible
- Quantitative: < 1 wt% detection limit typically; < 0.1 wt% with WDS (wavelength-dispersive)
- Mapping: scan beam → build elemental maps pixel by pixel (15–60 min for large map)

**EDS limitations:**
- Resolution limited by interaction volume (same as BSE) → can't do < 0.5 μm features in bulk
- Light element accuracy poor (C, N, O: matrix correction issues)
- Overlapping peaks (e.g., Mn Kβ overlaps Fe Kα) need WDS to resolve

### Electron Backscatter Diffraction (EBSD)
**EBSD** determines crystallographic orientation of each grain in SEM:
- Beam hits tilted sample (70°) → diffracted electrons form Kikuchi bands (diffraction pattern)
- Computer matches pattern to known crystal structure → orientation of each pixel
- Scan across surface → **orientation map** with color coding

**What EBSD gives:**
- Grain orientation map (IPF map: inverse pole figure colored by orientation)
- Grain boundary character: LAGB (< 15°) vs HAGB (> 15°) — see Ch 09
- Texture (orientation distribution function ODF)
- Strain/KAM map: local misorientation → identifies deformed vs recrystallized regions

**Application:** Confirming SX (single crystal) orientation [001] ± 10° from the growth axis (Ch 49); detecting stray grains.

---

## 4. SEM Sample Preparation and Applications

**Sample preparation:**
- Bulk samples: just need to fit in stage + be electrically conductive
- Non-conductive samples: thin carbon or Au/Pd coat (sputter) to prevent charging
- Fracture surfaces: no preparation needed → examine as-fractured (preserve fracture surface!)
- Cross-sections: mount, polish (like LOM), light etch

**Key SEM applications in metallurgy:**

**Fracture surface analysis (fractography):**
- Identify fracture mode from morphology:
  ```
  Dimple fracture (MVC) → Ductile, overload
  Cleavage → Brittle, transgranular (river marks, facets)
  Fatigue → Beach marks (LOM), striations (SEM, spacing = da/cycle)
  Intergranular → Grain facets, smooth surfaces
  ```
- Measure fatigue striations → calculate crack growth rate (da/dN = striation spacing)
- Locate fracture origin (stress concentrator: inclusion, pore, defect)

**Phase identification (BSE + EDS):**
- Map phases in complex alloys
- Identify inclusions (alumina vs MnS vs TiN in steel)
- Measure coating layer thickness + composition

---

## 5. Focused Ion Beam (FIB) — Cross-Sectioning and TEM Prep

**FIB** uses a focused Ga⁺ ion beam to mill material (like precision microscopic CNC):

**Applications:**
1. **Cross-section any surface without mechanical preparation:**
   - Cut through coating, oxidized surface, exact crack tip → observe cross-section in-situ in SEM
   - Used for: TBC cross-section at exact crack tip; electrochemical cell at electrode surface
   
2. **TEM sample preparation (FIB lift-out):**
   ```
   1. Deposit protective Pt cap on area of interest (to protect surface during ion milling)
   2. Mill two trenches around the ROI → thin membrane (~10 μm × 10 μm × 5 μm)
   3. Attach micromanipulator needle → cut free membrane → lift out
   4. Mount on TEM grid
   5. Further thin to < 100 nm by low-energy ion beam (FEI)
   → Ready for TEM
   ```

**FIB damage:**
Ion beam implants Ga into sample, creates amorphous layer (~20 nm) → cleaned by low-energy Ar ion milling or UV cleaning before TEM.

---

## 6. Transmission Electron Microscopy (TEM) — Principles

**TEM requires thin (< 100 nm) samples** so electrons can transmit through:

**Electron gun:** field emission → coherent, monochromatic beam (200 keV typical)

**Optics:**
```
Condenser lens → focus on sample
Objective lens → forms image + diffraction pattern
Intermediate lens → switch between image and diffraction modes
Projector lens → magnify onto screen/camera
```

**Magnification:** 1,000× to 1,000,000×

**Sample preparation options (besides FIB):**
- Electropolishing: twin-jet polish thin foil; for metals only
- Ion milling: Ar ion beam thins from both sides; for all materials
- Ultramicrotomy: diamond knife sections; for soft materials (polymers, biological)

---

## 7. TEM Imaging Modes: BF, DF, HRTEM, HAADF-STEM

### Bright-Field (BF-TEM)
- Objective aperture selects only the transmitted beam (not diffracted beams)
- Regions diffracting strongly → appear dark (less transmitted intensity)
- **Shows:** Dislocations (as dark lines), stacking faults, precipitates, grain boundaries, strain contrast

### Dark-Field (DF-TEM)
- Objective aperture selects one diffracted beam
- Regions satisfying the diffraction condition for that beam → bright
- **Shows:** Specifically highlights phases diffracting at selected angle → excellent for: γ′ particles (select superlattice reflection), specific phase identification

```
Example: In Ni superalloy, select γ′ {001} superlattice reflection:
  → γ′ cubes appear bright; γ matrix = dark → direct map of γ′ distribution
```

### High-Resolution TEM (HRTEM)
- Interference of transmitted + multiple diffracted beams → atomic lattice fringes
- Shows atomic columns directly (at < 0.2 nm resolution)
- **Shows:** Interface atomic structure, dislocation core structure, stacking faults, oxide film structure
- Requires: very thin specimen (< 20 nm), high-stability microscope, advanced operator

### HAADF-STEM (High-Angle Annular Dark-Field STEM)
- Beam scanned (like SEM) but in transmission mode
- HAADF detector: only collects electrons scattered at very high angles (Rutherford scattering)
- Intensity ∝ Z² (atomic number squared) → very strong chemical contrast
- Called "Z-contrast imaging"
- **Shows:** Individual heavy atoms (even single Re, W atoms in Ni matrix), chemical segregation at atom level

---

## 8. TEM Analysis: Diffraction, Phase Identification, Dislocations

### Selected Area Electron Diffraction (SAED)
- Select small region with diffraction aperture (~100–500 nm diameter)
- Electrons diffract → diffraction pattern on screen
- Spot pattern → single crystal (each spot = a specific (hkl) reflection)
- Ring pattern → polycrystalline (many grains in illuminated area)

**Phase identification from SAED:**
- Measure d-spacings from spot distances
- Index all spots → identify phase by d-spacing + geometry

**Orientation relationship determination:**
- How do two phases orient relative to each other? Measure angle between diffraction spots in each phase → Kurdjumov-Sachs, Nishiyama-Wassermann, etc.

### Imaging Dislocations in TEM
**Two-beam condition:** Tilt to diffracting condition for specific g (diffraction vector):
```
Dislocation visible when: g·b ≠ 0 (g = diffraction vector, b = Burgers vector)
Dislocation invisible when: g·b = 0 → "extinction" condition
```

Using multiple two-beam conditions with different g → determine **b** (Burgers vector) of the dislocation.

**What you can see:**
- Individual dislocation lines, loops, tangles
- Partial dislocations (smaller b) → see separately from full dislocations
- Stacking faults: fringes between partial dislocations
- APB (antiphase boundaries) in ordered phases (γ′): dark/bright fringes in DF imaging

### EDS in TEM (STEM-EDS)
With STEM mode (focused probe, ~0.2 nm) + EDS:
- Can map elemental distributions at ATOMIC scale (each pixel = 0.2 nm)
- Reveals: Re/W core-to-surface segregation in CMSX-4 dendrites, Ti/Al partitioning at γ/γ′ interface

---

## Summary

| Technique | Resolution | Sample | Information |
|-----------|-----------|--------|------------|
| SEM (SE) | 1–5 nm | Bulk, conductive | Surface topography, fracture morphology |
| SEM (BSE) | 5–20 nm | Bulk | Atomic number contrast, phase map |
| EDS | 0.5–2 μm (bulk) | Bulk | Elemental composition, X-ray maps |
| EBSD | ~50 nm | Bulk, polished | Grain orientation, texture, LAGBs |
| FIB | ~10 nm | Any | Cross-section, TEM prep |
| TEM (BF) | 0.2 nm | < 100 nm foil | Dislocations, phases, interfaces |
| HRTEM | 0.1 nm | < 20 nm foil | Atomic structure, interfaces |
| HAADF-STEM | 0.08 nm | < 50 nm foil | Z-contrast, chemical mapping |
| SAED | — | TEM foil | Phase ID, orientation relationships |
| STEM-EDS | 0.5 nm | TEM foil | Atomic-scale chemistry |

---

## Exercises

1. A failed turbine blade (CMSX-4 SX) needs fractographic analysis. (a) The fracture surface is received directly. Should you mount it and etch it first? Why or why not? (b) SEM (SE mode) shows: a flat, faceted region (2 × 2 mm) with river marks; surrounding this: dimpled regions. Identify: fracture mode at the center (Mode I brittle crack? LCF? HCF?) and at the periphery. (c) At 10,000× in the faceted region, no striations are visible. In the dimpled region, striation spacing = 1.2 μm/cycle. Estimate the fatigue crack growth rate da/dN. Using Paris Law da/dN = C(ΔK)^m (C = 2×10⁻¹², m = 3, ΔK in MPa√m), calculate the ΔK at that crack growth rate. (d) EDS point analysis on a bright inclusion at the fracture origin shows: 50%Ti, 10%N, 40%C. What compound is this? What is its role as a fracture origin?

2. EBSD analysis of a partially recrystallized Ti-6Al-4V billet: (a) The EBSD map shows regions with high KAM (kernel average misorientation > 1°) and regions with low KAM (< 0.5°). Which regions are deformed (work hardened) and which are recrystallized? Why does KAM correlate with dislocation density? (b) The texture shows strong {0001}⟨11-20⟩ basal fiber texture in the deformed regions, and random texture in the recrystallized regions. What does this tell you about the forging direction and how recrystallization eliminates it? (c) The billet specification requires < 5% LAGB (grain boundary misorientation 2–5°). EBSD measurement shows 12% LAGB. What is the metallurgical implication for fatigue crack growth? (d) After full recrystallization anneal (980°C / 2 hours), repeat EBSD would show what changes?

3. TEM analysis of γ′ in CMSX-4 after service (1,000 hours at 950°C): (a) BF image shows cuboidal γ′ precipitates have become elongated into "rafts" (plates) perpendicular to the [001] growth direction. What driving force causes rafting? (tensile creep stress + γ/γ′ mismatch → directional coarsening). (b) Use DF imaging with the {001} superlattice reflection: describe what you would see (γ′ bright/dark, γ matrix brightness). (c) HRTEM at the γ/γ′ interface shows: no APB (antiphase boundary), coherent interface with matching lattice fringe spacing = 0.358 nm (γ) and 0.357 nm (γ′). Calculate the misfit δ = (a_γ - a_γ′)/a_γ. Is this consistent with a coherent interface? (d) After very long service (5,000 hours at 1,000°C), γ′ has dissolved. Using the Ostwald ripening equation r³(t) = r₀³ + Kt (Ch 21), if r₀ = 0.3 μm and K = 0.002 μm³/hour, at what time does r = 1 μm (the size where Orowan bypass dominates and cutting contribution disappears)?

4. FIB cross-section of a failed MCrAlY + YSZ TBC on an IN738 vane: (a) The TBC is spalled. FIB cross-section at the delamination crack shows: TBC still attached on one side, bare TGO on the other. Is the crack through the TBC, through the TGO, or at the TGO/bond coat interface? SEM (BSE) shows bright white layer at the crack surface. What is this bright layer? (b) EDS mapping of the TGO shows: mostly O + Al, but also Cr + Ni enriched at the inner TGO surface. What does Cr + Ni enrichment indicate? (refer to Ch 32: what happens when Al depletes in bond coat). (c) TGO thickness = 7 μm. Is this within the typical life-limit of 5–7 μm? (d) A secondary crack is visible parallel to the main spall, 10 μm below the TGO in the bond coat. What is this likely caused by? (rumpling? SRZ formation? Sigma-phase precipitation?)

5. TEM dislocation analysis in deformed IN718 (ε = 5% compression at RT): (a) BF image under g = [002] diffraction condition shows 45 dislocation lines visible per 1 μm². Under g = [1-10], 12 lines visible. Under g = [11-2], 0 lines visible. Using the invisibility criterion g·b = 0: determine the Burgers vector b of the dislocations (options: a/2[110], a/2[101], a/2[011] for FCC; which combination satisfies the visibility/invisibility observations?). (b) Some dislocations appear in pairs connected by bright/dark fringe contrast (APB). What does this indicate about which phase they are traveling through? (c) The dislocation density ρ = 45 × 10¹² /m² (estimated from BF). Using Taylor hardening τ = τ₀ + αGb√ρ (α = 0.5, G = 80 GPa, b = 0.25 nm), calculate the shear stress increase due to work hardening.
