# Chapter 53: Investment Casting of Turbine Blades — The Complete Process

> **"A turbine blade is born as a wax sculpture smaller than your palm. It gets dressed in fourteen layers of ceramic slurry — each one painstakingly applied, stuccoed, and dried for 12 hours. The wax is melted away. The ceramic is fired to 1,050°C. Then, under high vacuum, a few hundred grams of superalloy at 1,450°C is poured into the ceramic shell, the entire mold is set on a water-cooled chill plate, and an electric furnace withdraws upward at 5 mm per minute. Thirty centimeters later, 30 minutes later, a single crystal grows. The ceramic is broken away. What remains is, atom by atom, the closest thing to perfect metal engineering humanity has yet achieved."**

---

## Table of Contents

1. [Overview of the Complete Process](#1-overview-of-the-complete-process)
2. [Ceramic Core Manufacturing](#2-ceramic-core-manufacturing)
3. [Wax Pattern Injection and Assembly](#3-wax-pattern-injection-and-assembly)
4. [Ceramic Shell Building](#4-ceramic-shell-building)
5. [Dewax, Burnout, and Sintering](#5-dewax-burnout-and-sintering)
6. [Vacuum Melting and Pouring](#6-vacuum-melting-and-pouring)
7. [Directional and Single Crystal Solidification (Bridgman)](#7-directional-and-single-crystal-solidification-bridgman)
8. [Knockoff, Cleaning, and Core Removal](#8-knockoff-cleaning-and-core-removal)
9. [Heat Treatment of SX Blades](#9-heat-treatment-of-sx-blades)
10. [Dimensional and NDT Inspection](#10-dimensional-and-ndt-inspection)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Overview of the Complete Process

**Complete process flow:**

```
CERAMIC CORE MANUFACTURING (weeks before blade casting)
  ↓
WAX PATTERN INJECTION (days before)
  ↓
CERAMIC SHELL BUILDING (7–14 days: 14 layers)
  ↓
DEWAX → BURNOUT → PREHEAT (hours)
  ↓
VACUUM INDUCTION MELT + POUR + BRIDGMAN DS/SX (30–60 min)
  ↓
COOL DOWN + KNOCKOFF (hours)
  ↓
ABRASIVE BLAST + ACID LEACH (core removal) (days)
  ↓
HIP (hours)
  ↓
SOLUTION HEAT TREATMENT + AGING (hours)
  ↓
NDT INSPECTION: X-ray, Laue, FPI, CMM, flow test (days)
  ↓
COATING: aluminide/MCrAlY + YSZ TBC (days)
  ↓
FINAL INSPECTION → APPROVED BLADE READY FOR ENGINE
```

Total cycle time: 4–8 weeks for a production SX turbine blade.

---

## 2. Ceramic Core Manufacturing

The ceramic core defines the internal cooling passages. It is the most dimensionally critical component of the blade, because the cooling channel geometry determines the blade's thermal capability.

**Core composition:**
- Fused silica (SiO₂) — primary filler (low CTE, thermally stable)
- Colloidal silica binder (SiO₂ sol) — holds particles together before firing
- Plasticizer (wax or thermoplastic) — for injection molding

**Core manufacturing process:**
```
1. Mix core material (SiO₂ + binder + plasticizer)
2. Inject into precision metal die → core blank (1–5 mm thick, complex internal geometry)
3. Cure (binder crosslinks)
4. Fire at 1,150–1,200°C → sinter SiO₂ → remove organic binder
5. Dimensional inspection (CMM): ± 0.1 mm tolerance
6. Visual inspection for cracks, chips
7. Load into wax die
```

**Core geometry defines:**
- Radial cooling holes spacing
- Film cooling hole exit positions (must align with drilling pattern later)
- Leading edge impingement jet holes
- Trailing edge slot exits
- Internal cross-bridges (serpentine coolant path)

**Core breakage during casting:**
The ceramic core must survive:
- Thermal shock when hot superalloy contacts it (~1,350°C)
- Differential thermal contraction as blade cools (metal contracts ~1.5%, ceramic ~0.5%)
- Chemical reaction with superalloy at casting temperature

**Core leaching (removal after casting):**
SiO₂ core is dissolved by KOH (potassium hydroxide) solution at 150°C, 150 psi autoclave:
- KOH + SiO₂ → K₂SiO₃ + H₂O (silicate dissolved)
- Pure SiO₂ dissolves without attacking Ni alloy
- Multiple leach cycles (4–6 hours each × 2–3 cycles) for complete removal
- Final rinse to ensure no K residue (K causes hot corrosion if left)

---

## 3. Wax Pattern Injection and Assembly

**Wax injection:**
- Precision-heated wax (typically petroleum wax + polymer blend, T_inject = 60–70°C)
- Injected under pressure (0.1–0.5 MPa) into metal die with ceramic core pre-positioned
- Die includes: airfoil form, root firtree, tip details, integral platform

**Wax pattern quality:**
- Dimensional accuracy: ±0.1–0.2 mm (wax die + wax shrinkage correction)
- Surface finish: Ra ≤ 0.8 μm (very smooth die → smooth casting)
- Core position check: use X-ray on wax pattern to verify ceramic core position before continuing

**Tree assembly:**
Multiple wax patterns are assembled onto a central wax sprue:
```
Central sprue (wax cylinder, 50 mm dia)
  → Wax runners branching to each blade pattern
  → 4–12 blades per tree depending on size
  
Tree assembly considerations:
  - Gates at top (under-filling is better than over-filling with heavy alloy)
  - Runner angle: gravity feed from top
  - No sharp angles (turbulence → inclusions)
  - Riser design: not typical for SX castings (directional solidification acts as its own riser)
```

**Wax pattern inspection:**
- Visual: surface finish, completeness
- X-ray: ceramic core position (key inspection — if core is off-position, the whole blade will be rejected after casting, wasting all subsequent processing cost)

---

## 4. Ceramic Shell Building

**Goal:** Build a ceramic shell around the wax tree that is:
- Strong enough to hold molten metal (must not fail at 1,450°C under metallostatic pressure)
- Has the correct thermal properties for DS/SX solidification
- Has adequate permeability (no trapped gas)
- Replicates wax surface to < 0.5 μm Ra

**Shell building cycle (for each of 14 layers):**
```
1. DIP: Immerse wax tree in ceramic slurry (15–30 seconds)
2. DRAIN: Remove excess slurry (orientation optimized)
3. STUCCO: Sprinkle ceramic aggregate (different size for different layers)
4. DRY: 8–12 hours at controlled T/humidity (critical)
```

**Slurry composition:**
- Prime coats (layers 1–2): Fine zirconia (ZrO₂) + colloidal silica binder
  → Fine because they replicate the surface → Ra must be < 0.8 μm
  → Zirconia provides: chemical stability with Ni alloy, low reactivity
- Backup coats (layers 3–14): Coarser fused silica + colloidal silica
  → Build structural strength
  → Fused silica (thermal shock resistant, low CTE)

**Stucco materials:**
- Prime layers: fine zirconia or fused alumina (< 120 mesh = 0.125 mm)
- Backup layers: coarser fused silica (25–50 mesh = 0.3–0.7 mm)

**Drying conditions:**
- Temperature: 22–24°C (room T) or slight warming
- Humidity: 40–60% RH (critical — too dry = cracks; too wet = incomplete drying → failure during burnout)
- Each layer must fully gel before next layer → colloidal silica gelation
- Infrared (IR) drying can accelerate to 4–6 hours per layer

**Final shell structure:**
```
Wax (to be removed)
    ↓
Prime layer 1: 0.3 mm fine ZrO₂ (replicates surface)
Prime layer 2: 0.3 mm fine ZrO₂
Backup layers 3–14: each ~0.5 mm coarser SiO₂
    ↓
Final shell thickness: 6–10 mm total
```

---

## 5. Dewax, Burnout, and Sintering

**Dewax (wax removal):**
Flash autoclave (steam autoclave, or Speedwax):
- Steam atmosphere at 165°C + 5 atm
- Wax melts instantly (before it can expand and crack the shell)
- KEY INSIGHT: If wax is heated slowly, it expands (thermal expansion) before melting → cracks the shell
- Flash autoclave heats ALL surfaces simultaneously → wax melts from outside before it expands → drains out without cracking shell

**Burnout:**
After dewax, organic residue (wax remnants, mold release agents) remains in shell:
- Heat to 1,050°C in air furnace
- Burnout removes: residual wax (oxidizes to CO₂), polymer binder (burns out)
- Also: sinters the ceramic shell (partial sintering → stronger)
- Duration: 1–2 hours at 1,050°C

**Pre-heat before pouring:**
Mold is maintained at 1,050°C until poured:
- Pre-heated mold ensures metal fills completely (no premature solidification)
- For SX Bridgman: mold placed in furnace hot zone before pouring

---

## 6. Vacuum Melting and Pouring

**Why vacuum:**
CMSX-4 contains 5.6%Al, 1%Ti — both form oxides/nitrides instantly in air → inclusions → fatigue failure.

**Vacuum Induction Melting (VIM) for casting heats:**
```
Vacuum chamber (< 10⁻² Pa = 10⁻⁴ mbar)
  ↓
Induction coil melts superalloy charge in alumina crucible
  (T = 1,430–1,450°C — 30–50°C above liquidus)
  ↓
Pre-heat ceramic mold in Bridgman hot zone (within the same vacuum chamber)
  ↓
Tilt crucible → pour into ceramic mold (through funnel → sprue → gates → blades)
  ↓
Immediately begin Bridgman withdrawal
```

**Charge material:**
- Virgin alloy: primary VIM-melted master alloy
- Return: clean scrap from gates, risers, rejected blades (controlled ratio)
- Composition check: optical emission spectrometry of melt sample before casting

**Pouring:**
- Superheat: 30–50°C above liquidus → complete fill without premature solidification
- Too hot: excessive mold erosion (ceramic attack by hot metal)
- Too cold: misruns (metal solidifies before filling)
- Pour rate: controlled by nozzle stopper rod → turbulence-free, laminar flow

---

## 7. Directional and Single Crystal Solidification (Bridgman)

This is the core process that makes DS and SX blades possible. The complete Bridgman process was covered in Ch 48 and 49. Summary here for complete process context:

**Standard Bridgman process:**
```
Hot zone (heating elements, 1,550°C)
    ↓
Baffle (separates hot and cold zones)
    ↓
Ceramic mold with liquid superalloy
    ↓
Withdrawal at constant rate V_w (3–10 mm/min for SX)
    ↓
Cold zone / radiation baffle
    ↓
Water-cooled copper chill plate at bottom
```

**Key variables controlled in real-time:**
- Hot zone temperature → controls superheat → controls liquidus isotherm flatness
- Withdrawal rate → controls solidification velocity R (not directly G)
- Heater power profile (multi-zone furnaces have 5–10 independent zones)
- Vacuum level → must stay < 10⁻² Pa throughout

**Grain selector function:**
```
Wax tree has "pigtail" or "spiral" selector at the bottom of each blade:
  - A tortuous (winding) passage ← only ONE grain can grow through
  - At Bridgman start: multiple grains nucleate on the chill plate
  - As withdrawal proceeds: most grains compete and die
  - After passing through the grain selector: only ONE grain survives = the SX
  
The selector CANNOT control which crystallographic orientation survives —
it only ensures a single grain.
The orientation that survives depends on which grain happened to have its
fastest growth direction (near-[001] for FCC) aligned with the heat flow.
```

**Typical process parameters for CMSX-4 SX:**
- Withdrawal rate: 5 mm/min
- Hot zone temperature: 1,540°C
- Vacuum: < 5 × 10⁻³ Pa
- Total blade solidification time: ~25–40 min (blade span = 120–200 mm)
- Cooling from liquidus to room temperature: 2–8 hours

---

## 8. Knockoff, Cleaning, and Core Removal

**Shell removal (knockout):**
After solidification is complete and blade is cooled:
1. Water jet or mechanical vibration knocks off the ceramic shell (brittle after thermal cycling)
2. Shot blast with alumina shot to remove remaining shell residue
3. Visual inspection for ceramic fragments embedded in surface

**Gate removal:**
Cut off the wax tree runner connections with abrasive cut-off wheel or EDM:
- Cut close to blade → minimal excess material (machined off later)
- Leave root and tip datum surfaces for CMM fixture

**Core removal (acid leaching):**
KOH autoclave leaching: 150°C / 150 psi / multiple cycles → dissolves SiO₂ core completely.
After leaching: rinse in hot water → verify no K residue (K causes hot corrosion if retained).

**Confirmation of core removal:**
- Flow test: air forced through cooling passages → measure flow rate (verifies all channels open)
- X-ray: verify no ceramic core fragments remain inside

---

## 9. Heat Treatment of SX Blades

**Three-step heat treatment (for CMSX-4):**

**Step 1 — Solution treatment (homogenization):**
- Purpose: dissolve all γ′ + homogenize dendritic segregation of Re, W, Ta
- Temperature: 1,290–1,300°C (just below incipient melting temperature ~1,315°C)
- Time: 4–6 hours
- Atmosphere: Ar or vacuum (prevent oxidation at high T)
- Risk: incipient melting (if temperature too high → liquid forms at grain boundaries → blade scrapped)
- Result: all γ′ dissolved, dendritic banding eliminated, uniform composition

**Step 2 — First age (primary aging):**
- Temperature: 1,100–1,140°C
- Time: 4 hours
- Purpose: Nucleate coarse γ′ cubes (~450 nm) with correct γ/γ′ misfit
- Controls: γ′ size and volume fraction

**Step 3 — Second age:**
- Temperature: 870°C
- Time: 20 hours
- Purpose: Nucleate fine secondary γ′ (~50 nm) in the γ channels between primary γ′
- Result: bi-modal γ′ distribution → optimizes both creep (coarse) and fatigue (fine) resistance

**Heat treatment tolerance:**
Temperature control: ±3°C in certified heat treatment furnace (AMS 2750)
Atmosphere: must maintain vacuum or Ar ≤ 10⁻² Pa to prevent oxidation

---

## 10. Dimensional and NDT Inspection

Covered in detail in Ch 52. Complete sequence summary:

**100% inspection for EVERY blade:**
1. X-ray (2 orientations) → porosity, inclusions, core geometry
2. Laue diffraction → orientation, stray grains
3. FPI → surface cracks, tears
4. CMM → 200+ points on airfoil profile, root firtree dimensions
5. Flow test → cooling air flow rate and distribution
6. Visual → surface finish, coatings (if applicable)

**Acceptance criteria defined by:**
- Customer engineering drawing (geometric tolerances)
- AMS specifications (material, heat treatment, NDT acceptance)
- OEM-specific internal specifications (supplementary requirements)

---

## Summary

| Step | Critical Variable | Defect if Wrong | Detection |
|------|-----------------|----------------|-----------|
| Ceramic core | Dimensional accuracy | Core shift → thin wall → overheating | CMM, X-ray on wax |
| Wax injection | Core position | Hot tear, misrun | X-ray on wax |
| Shell building | Humidity, dry time | Shell crack → metal leak | Visual |
| Dewax | Flash-autoclave speed | Shell crack | Visual |
| VIM pouring | Superheat, vacuum level | Inclusions, misrun, porosity | X-ray RT |
| Bridgman | G/V ratio, withdrawal rate | Stray grain, freckle, LAGB | Laue, EBSD |
| Core removal | KOH completeness | Core remnant → blocked cooling | Flow test, X-ray |
| Solution HT | T control (± 3°C) | Incipient melt (scrap) or under-solution (segregation) | Metallographic, composition |
| Aging | T, time | Wrong γ′ size/fraction | Metallographic, SEM |

---

## Exercises

1. Ceramic shell design for a new blade geometry with an unusually thin trailing edge (0.5 mm): (a) The trailing edge solidifies very quickly in Bridgman. What defect does rapid solidification here risk? (b) Prime coat uses ZrO₂ slurry — what property of ZrO₂ makes it preferred over Al₂O₃ for the prime coat in contact with CMSX-4? (c) The backup coats use fused silica. If the shell develops a crack through all 14 layers during burnout, what would happen when liquid metal is poured? (d) Shell permeability is required to allow gases to escape as metal fills. If the shell is too dense (no permeability), what defect results in the casting?

2. VIM pouring analysis for a 12-blade tree: total metal required = 2.4 kg (12 × 200g blades + gates/sprues). VIM furnace holds 3 kg CMSX-4 charge. (a) Superheat temperature = 40°C above T_liquidus = 1,380°C → T_pour = 1,420°C. If the mold is at 1,050°C (preheat), calculate the temperature drop as metal hits the mold (assume equal heat capacities, rough estimate: ΔT = (T_metal - T_mold) × m_metal / (m_metal + m_mold) where m_mold ≈ 0.5 kg ceramic). (b) If superheat drops to 10°C (T_pour = 1,390°C): what is the risk? Why is it critical to maintain superheat especially at the trailing edge? (c) The entire pour happens in < 10 seconds. Calculate the average pour rate in kg/s. Why is slow pouring undesirable (despite reduced turbulence)?

3. Bridgman process parameter selection: Target solidification conditions for SX CMSX-4 blade: G/R = 20,000°C·s/mm² is needed for single-crystal solidification (no stray grains). (a) If the thermal gradient G = 40°C/mm (fixed by furnace design), what maximum withdrawal rate R (mm/min) can be used? (b) If R is increased to 8 mm/min while keeping G = 40°C/mm, does G/R meet the requirement? What defect would result? (c) A second furnace has G = 60°C/mm. What withdrawal rate can now be used? What is the advantage of higher G? (d) Increasing G requires higher power input or better insulation. A new hot zone design increases G from 40 to 55°C/mm. Calculate the allowed withdrawal rate increase and its effect on throughput (blades per hour).

4. Solution heat treatment criticality: CMSX-4 incipient melting temperature T_IMT = 1,315°C. Solution treatment target = 1,295°C ± 3°C. (a) If the furnace control sensor drifts +25°C → actual temperature = 1,320°C > T_IMT. Describe what happens to the blade microstructure. Which specific features of the blade (think: high-segregation interdendritic regions) melt first? (b) After incipient melting, the blade is cooled. What microstructure feature appears in the formerly-melted regions? Why is this fatal for the blade? (c) Incipient melting detection: the blade surface develops an irregular "orange peel" appearance. Can this be detected by visual inspection? What NDT method would reveal the internal melt? (d) To reduce IMT risk, the solution treatment is done in multiple steps: first at 1,265°C/2h → then 1,285°C/2h → then 1,295°C/6h. Explain why this step-up sequence reduces IMT risk.

5. Complete process yield calculation: For a batch of 120 investment casting shells (10 trees × 12 blades per tree): Stage-by-stage rejection rates: (a) Wax/shell stage: 2% trees rejected for shell cracks → how many blades lost? (b) Casting stage: 5% trees fail (misrun, metal leak) → how many blades? (c) Laue inspection: 6% of individual blades rejected (stray grain 4%, misorientation 2%) (d) X-ray: 5% of remaining blades rejected (porosity — but 80% of these can be recovered by HIP) (e) HIP: 80% of X-ray-failed blades pass after HIP. (f) FPI: 1% of remaining blades rejected. (g) CMM/flow: 2% of remaining blades rejected. Calculate the final blade yield (number passing all stages). Calculate cost of each blade that passes if total processing cost = $800/blade regardless of rejection point (including inspection and handling of rejected parts).
