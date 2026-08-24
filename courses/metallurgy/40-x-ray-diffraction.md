# Chapter 40: X-Ray Diffraction — Crystal Structure by Diffraction

> **"Every crystalline metal carries within its lattice a unique diffraction fingerprint. X-rays bounce off the atomic planes and interfere — constructively if the geometry is right, destructively if not. The result is a set of peaks at specific angles, as distinctive as a bar code. From those peaks we can read: which phases are present, how much of each, the stress in the lattice, and how the grains are oriented. X-ray diffraction is the most information-dense non-destructive measurement in materials science."**

---

## Table of Contents

1. [Bragg's Law — The Foundation](#1-braggs-law--the-foundation)
2. [The X-Ray Diffractometer](#2-the-x-ray-diffractometer)
3. [The Powder Diffraction Pattern](#3-the-powder-diffraction-pattern)
4. [Phase Identification (Qualitative)](#4-phase-identification-qualitative)
5. [Quantitative Phase Analysis (Rietveld)](#5-quantitative-phase-analysis-rietveld)
6. [Residual Stress by XRD (sin²ψ method)](#6-residual-stress-by-xrd-sinψ-method)
7. [Texture and Pole Figures](#7-texture-and-pole-figures)
8. [Single Crystal Orientation (Laue Diffraction)](#8-single-crystal-orientation-laue-diffraction)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Bragg's Law — The Foundation

X-rays (wavelength λ ≈ 0.05–0.25 nm) have wavelengths comparable to interatomic spacings in crystals → diffraction occurs.

**Bragg's Law:**
```
Constructive interference when:
   nλ = 2 × d_hkl × sin(θ)

where:
   n    = integer (diffraction order, usually 1)
   λ    = X-ray wavelength (fixed for given source, e.g., Cu Kα = 0.15406 nm)
   d_hkl = interplanar spacing of (hkl) planes
   θ    = Bragg angle (half the angle between incident and diffracted beam)
```

**Physical meaning:**
```
X-ray beam
     ↘ θ
────────────  plane 1 ←── d_hkl
────────────  plane 2

Path difference = 2 × d × sin(θ)
= nλ → constructive interference → peak in diffraction pattern
≠ nλ → destructive interference → no peak
```

**d-spacing for cubic crystal:**
```
d_hkl = a / √(h² + k² + l²)

where a = lattice parameter (e.g., γ-Fe FCC: a = 0.3591 nm; α-Fe BCC: a = 0.2866 nm)
```

This means each (hkl) reflection gives a different peak at a different 2θ angle.

---

## 2. The X-Ray Diffractometer

**Bragg-Brentano geometry (most common for powder XRD):**

```
X-ray source (Cu Kα anode, λ = 0.15406 nm)
        ↘ θ (source angle)
[Divergence slit] → sample surface (rotates at θ)
        ↗ θ (detector angle)
[Receiving slit] → detector (rotates at 2θ)
```

Source and detector BOTH move symmetrically: source at -θ, detector at +θ relative to sample = "θ/2θ scan."

**X-ray sources:**
- Laboratory: Cu Kα (most common), Mo Kα (higher energy, less fluorescence for Fe-containing alloys), Cr Kα (for residual stress)
- Synchrotron: tunable wavelength, extremely intense, very high resolution — for research

**Detector types:**
- Point detector (proportional counter): traditional, slow
- Linear detector (position-sensitive): fast (collects full angular range simultaneously)
- 2D detector: collects full Debye-Scherrer ring → texture analysis

---

## 3. The Powder Diffraction Pattern

A powder (or fine-grained) sample has grains in all orientations → every (hkl) family that satisfies Bragg's Law will diffract → multiple peaks.

**What each peak tells you:**
- **Peak position (2θ):** d-spacing → lattice parameter → identify phase
- **Peak intensity:** structure factor (which atoms, where in unit cell) + multiplicity → identify phase + composition
- **Peak breadth:** crystallite size (Scherrer: B = Kλ/(L·cos θ)) + strain broadening
- **Peak area:** amount of that phase

**Example: Iron powder pattern (Cu Kα):**
```
α-Fe (BCC, a = 0.2866 nm):
  2θ = 44.67° → (110)  [strongest peak]
  2θ = 65.0°  → (200)
  2θ = 82.3°  → (211)
  
γ-Fe would give different peaks (FCC, a = 0.3591 nm):
  2θ = 43.8°  → (111)
  2θ = 50.9°  → (200)
```

**Systematic absences (structure factor rules):**
Not all (hkl) reflections appear — some are extinct due to lattice type:
- FCC: only (all even) or (all odd) h,k,l reflections allowed
- BCC: only h+k+l = even reflections allowed

This is how you confirm FCC vs BCC: by the set of peaks that appear.

---

## 4. Phase Identification (Qualitative)

**Match measured d-spacings to the PDF (Powder Diffraction File, ICDD database):**
- Measure 2θ of all peaks → calculate d_hkl from Bragg's Law
- Match d-spacing list against database of 900,000+ known compounds
- Software (JADE, HighScore, Match!) does this automatically

**Detection limit:** ~1–2 wt% for phases with distinct peaks; lower sensitivity for phases that overlap with major phase.

**Common metallurgical applications:**
- Steel: confirm α-ferrite vs γ-austenite (after heat treatment or at elevated T)
- Stainless steel: detect σ-phase (embrittlement) or martensite
- TBC: confirm stable tetragonal YSZ (t') vs monoclinic ZrO₂ (dangerous: volume change on phase transformation → TBC cracking)
- Carbides in tool steels: distinguish M₇C₃ vs M₂₃C₆ vs MC carbides
- Corrosion products: identify oxide/hydroxide type on steel surface

---

## 5. Quantitative Phase Analysis (Rietveld)

**Rietveld method:** Fit a calculated diffraction pattern (from crystal structure models) to the measured pattern by least-squares refinement.

**What Rietveld refinement extracts:**
- Phase fractions (weight% of each phase)
- Lattice parameters (accurately, < 0.0001 nm)
- Crystallite size and strain (from peak broadening)
- Preferred orientation (texture)
- Atomic positions

**Process:**
```
Input: measured pattern + crystal structure models for each phase
→ Refine: scale factors, background, peak shape, lattice parameters
→ Output: Rwp (weighted residual) → quality of fit; phase fractions
```

**Application: Retained austenite in steel**
Retained austenite (γ, FCC) in hardened steel reduces fatigue life. XRD quantifies it:
- Standard method: ASTM E975
- Measure peak areas for α and γ reflections → calculate weight fraction
- Detection limit: ~0.5 vol%

**YSZ phase transformation in TBC:**
- Tetragonal YSZ (t'): stable in service → no volume change
- Monoclinic ZrO₂: transforms from tetragonal at high T (low Y₂O₃ content) → volume change → cracking
- XRD monitors m-ZrO₂ fraction in service-exposed TBCs

---

## 6. Residual Stress by XRD (sin²ψ method)

XRD measures elastic lattice strain → converts to stress using Hooke's Law.

**Principle:**
Stress in a surface deforms the crystal lattice (changes d-spacings). XRD measures this change.

**d vs ψ relationship:**
```
d_ψ = d₀ × [1 + (1+ν)/E × σ × sin²ψ - 2ν/E × σ]

where:
   d_ψ  = measured d-spacing at tilt angle ψ
   d₀   = stress-free d-spacing
   ψ    = tilt angle of sample (0° = surface normal to beam; varies to ~45°)
   σ    = surface stress (unknown)
   E, ν = elastic constants of the material
```

**Procedure:**
1. Measure peak position at multiple ψ angles (0°, 10°, 20°, 30°, 40°, 45°)
2. Plot d vs sin²ψ
3. Slope = (1+ν)/E × σ → solve for σ

```
d_ψ vs sin²ψ:
│                          /
│                        /
│                      /   slope = (1+ν)σ/E
│                    /
│──────────────────/
└──────────────────────── sin²ψ
0                    1
```

**Note:** If sin²ψ plot is NOT linear (curvature or oscillation), the sample has a stress gradient or texture — more complex analysis needed.

**Measurement depth:**
X-rays penetrate ~5–20 μm in steel (depends on λ and material). So XRD residual stress = average over this thin surface layer.

**Practical specifications:**
- Shot-peened compressor disk: require σ ≤ -200 MPa at surface, measured by XRD on every lot
- Critical aerospace forgings: XRD RS measured to verify compliance with AMS 2430 (shot peening spec)

---

## 7. Texture and Pole Figures

**Crystallographic texture** = preferred orientation of grains. XRD measures texture by measuring how intensities vary as the sample is tilted and rotated.

**Pole figure measurement:**
- Choose one (hkl) reflection
- Measure peak intensity as sample is tilted (χ) and rotated (φ)
- Plot intensity on a hemisphere → pole figure

```
Pole figure example: {001} for a rolled FCC metal

     N (χ=0°)
     ┃
──── ┿ ──── E (φ=0°)
     ┃
     S

High intensity at: center (ND) and edges (RD, TD) → shows cube texture
```

**Orientation Distribution Function (ODF):**
Full 3D description of texture using Euler angles (φ₁, Φ, φ₂). Generated from multiple pole figures.

**Key textures in metals:**
| Metal | Process | Texture | Effect |
|-------|---------|---------|--------|
| FCC (Al, Cu) | Rolling | Copper {112}⟨111⟩ + S component | Anisotropy; earing in deep drawing |
| BCC (steel) | Rolling | {001}⟨110⟩ + {112}⟨110⟩ | Mild steel: {111} for good drawability |
| HCP (Ti) | Forging | Basal texture ({0001} in compression plane) | Anisotropy in fatigue crack growth |
| SX Ni (casting) | DS/SX | {001} fiber texture | Lower elastic modulus in [001] |

---

## 8. Single Crystal Orientation (Laue Diffraction)

**Laue diffraction** uses white X-rays (broad wavelength range) on a stationary SINGLE crystal → each wavelength diffracts at different angles → produces spots on flat film/detector.

**Back-reflection Laue:**
```
White X-rays → single crystal → back-scattered diffraction → spots on film behind X-ray source
```

Each set of planes diffracts at exactly the right λ for the given geometry → spot positions reveal crystal orientation.

**Application: Turbine blade SX orientation verification**

Every single crystal turbine blade must have [001] within 10° of the blade growth axis. Laue diffraction is the standard method:

```
X-ray film placed between source and blade (or behind)
→ Expose for 10 minutes
→ Develop film → spot pattern
→ Measure spot positions → software calculates crystal orientation
→ Pass/fail vs. specification: [001] ± 10° primary, [010] ± 10° secondary
```

**What Laue cannot do:** Detect LAGBs within the crystal (sub-grain misorientation < 5°) — needs EBSD or higher-resolution synchrotron Laue.

**Stray grain detection:**
If a stray grain (equiaxed, random orientation) is present in the blade → TWO spot patterns on same film → automatic reject. This is why Laue inspection is 100% on SX blades.

---

## Summary

| XRD Application | Method | Information | Detection Limit |
|-----------------|--------|-------------|----------------|
| Phase identification | θ/2θ scan + PDF matching | Which phases are present | ~1–2 wt% |
| Phase quantification | Rietveld refinement | Weight fraction of each phase | ~0.5 vol% |
| Retained austenite | ASTM E975 | % γ in hardened steel | ~0.5 vol% |
| Residual stress | sin²ψ, multiple tilts | Surface stress (MPa) | ± 20 MPa |
| Lattice parameter | Rietveld / Bragg | a, b, c, α, β, γ | ± 0.0001 nm |
| Texture | Pole figures + ODF | Preferred orientation | Qualitative + quantitative |
| SX orientation | Laue diffraction | [hkl] axis, misorientation | Stray grains detected |
| TBC phase | θ/2θ | t' vs m-ZrO₂ fraction | ~1 vol% |

---

## Exercises

1. A sample contains α-Fe (BCC, a = 0.2866 nm) and γ-Fe (FCC, a = 0.3591 nm). Using Cu Kα (λ = 0.15406 nm): (a) Calculate 2θ for the (110) reflection of α-Fe and the (111) reflection of γ-Fe. Use Bragg's Law. (b) After carburizing, the FCC austenite lattice parameter increases to 0.3605 nm (C dissolved in interstitial sites). Calculate the new 2θ for the (111) γ reflection. What is the peak shift? (c) If peak shift of 0.1° corresponds to 50 MPa stress, estimate the stress in a sample showing a 0.05° shift. (d) After quenching, the γ → martensite (BCT). The c-axis of BCT martensite with 0.5%C = 0.296 nm, a-axis = 0.284 nm. Calculate the c/a ratio. What 2θ does the (002) tetragonal peak appear at?

2. Rietveld refinement of a duplex stainless steel (2205) sample after PWHT at 475°C for 2 hours gives: α-ferrite = 52 wt%, γ-austenite = 42 wt%, σ-phase (FeCr) = 6 wt%. (a) The σ-phase d-spacing shows it has tetragonal structure (a = 0.879 nm, c = 0.455 nm). What effect does σ-phase embrittlement have on impact toughness? (b) In the untreated sample, σ-phase = 0 wt%. The detection limit of the diffractometer is 0.5 wt%. What maximum σ content could have been present without detection? Is this a problem? (c) To dissolve σ-phase: heat to 1,050°C for 30 min. After this re-anneal, repeat XRD. Predict the new phase fractions. (d) Why is σ-phase formation temperature (475–900°C) called the "475°C embrittlement" range for ferritic/duplex stainless?

3. Residual stress measurement on a shot-peened Ti-6Al-4V compressor disk: using Cr Kα radiation (λ = 0.22897 nm), measuring the (213) reflection at multiple ψ angles: ψ = 0°: 2θ = 139.0°; ψ = 15°: 139.3°; ψ = 30°: 139.8°; ψ = 45°: 140.4°. (a) Calculate d-spacing at each ψ using Bragg's Law. (b) Plot d vs sin²ψ and calculate the slope. (c) Given E = 110 GPa and ν = 0.31 for Ti-6Al-4V, calculate the residual stress σ using: slope = (1+ν)/E × σ. (d) Is this compressive or tensile? Is it above the -200 MPa requirement for shot-peened disks? (e) At 50 μm depth (accessible by layer-removal + XRD), the stress is -350 MPa. At what depth does the stress cross from compressive to tensile? What determines this depth in shot peening?

4. Laue diffraction quality control of a CMSX-4 SX turbine blade: (a) The Laue film shows a clear single spot pattern. Software measures the primary [001] axis is 7° from the blade Z-axis (growth direction). The secondary [100] is 12° from the circumferential direction. The specification is: primary ≤ 10°, secondary ≤ 10°. Pass or fail? Why does secondary orientation matter? (b) A second blade shows TWO overlapping spot patterns — one with the correct orientation, one rotated 25°. What microstructural feature does this indicate? At what stage of the investment casting process did this defect form? (c) Laue inspection is performed at 100% (every blade). Why is sampling inspection (e.g., 1 in 20) not acceptable for this test? (d) A research program uses synchrotron Laue to detect sub-grain boundaries (LAGBs) within a "passing" SX blade. What misorientation angle would concern you for creep performance? (refer to Ch 47)

5. XRD texture analysis of cold-rolled 3104 aluminum alloy (for beverage cans): {111} pole figure shows high intensity at χ = 40–80° (a ring of high intensity), and {001} pole figure shows four lobes at ϕ = 0°, 90°, 180°, 270° at χ ≈ 35°. (a) This texture is called "copper rolling texture" ({112}⟨111⟩). Why does cold rolling produce preferred orientation? (b) The earing test (deep drawing a cup) shows 4 ears at 0°, 90°, 180°, 270° from the rolling direction. How does crystallographic texture cause earing? (hint: {111} grains deform more easily in deep drawing; uneven distribution → uneven cup height). (c) A target texture of {111} fiber (all grains with {111} parallel to sheet) would give 0 earing. What annealing treatment promotes this texture? (d) The beverage can wall has 15% earing. After design change: use lower cold-work ratio + specific annealing sequence → earing reduced to 2%. What is the economic benefit? (Assume 1 billion cans per year, each can 0.6g aluminum lighter with lower earing → calculate total Al saving in tonnes).
