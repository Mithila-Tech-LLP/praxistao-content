# Chapter 59: Turbine Blade Repair — Restoring Single Crystals

> **"Repairing a turbine blade is harder than making it. When you weld a single crystal, the heat-affected zone recrystallizes — you've introduced grain boundaries into a material whose entire purpose was to eliminate them. When you strip and reapply TBC, the blade has already experienced 15,000 hours of service stress, and the alloy chemistry near the surface is no longer nominal. Repair metallurgy is applied metallurgy at the frontier of what's possible."**

---

## Table of Contents

1. [Economics of Repair vs. Replace](#1-economics-of-repair-vs-replace)
2. [Inspection and Damage Assessment](#2-inspection-and-damage-assessment)
3. [Stripping the Old Coating](#3-stripping-the-old-coating)
4. [TBC Reapplication](#4-tbc-reapplication)
5. [Tip Repair — Brazing and Rebuilding](#5-tip-repair--brazing-and-rebuilding)
6. [Crack Repair — Welding Challenges in SX](#6-crack-repair--welding-challenges-in-sx)
7. [Recrystallization — The Enemy of Repair](#7-recrystallization--the-enemy-of-repair)
8. [Diffusion Bonding for Segment Replacement](#8-diffusion-bonding-for-segment-replacement)
9. [Re-Solution and Re-Aging After Repair](#9-re-solution-and-re-aging-after-repair)
10. [Repair Qualification and Return to Service](#10-repair-qualification-and-return-to-service)
11. [Life Extension Strategies](#11-life-extension-strategies)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. Economics of Repair vs. Replace

A new HPT blade costs $10,000–50,000. A blade repair costs $2,000–8,000 and can restore 60–80% of remaining design life.

**Repair economics:**
- An engine overhaul visits the shop every 4,000–8,000 flight hours
- Each HPT blade is inspected at each shop visit
- ~60–70% of blades can be repaired; ~30–40% must be scrapped

**Business case for repair:**
```
New blade:    $30,000
Repair cost:  $4,000 (extends life ~6,000 EFH)
Repair value per life-hour: $4,000 / 6,000 = $0.67/EFH

New blade cost per life-hour: $30,000 / 30,000 EFH = $1.00/EFH

Saving per repaired blade: ~$0.33/EFH × 6,000 EFH = ~$2,000 per repair
With 80 blades per engine: $160,000 savings per shop visit per engine
```

For an airline with 300 engines undergoing overhaul every 6,000 EFH, the fleet repair savings are ~$50M per year vs. all-new blades.

This is why a dedicated blade repair industry exists — MRO (Maintenance, Repair, Overhaul) shops specialized in turbine blade repair include Chromalloy, Sulzer, Praxair Surface Technologies, StandardAero, and engine OEM service centers.

---

## 2. Inspection and Damage Assessment

Before any repair, each blade undergoes comprehensive inspection to determine repairability:

### Visual and Dimensional

- Macroscopic examination under magnification (10–50×): cracks, TBC condition, erosion, hot corrosion
- Dimensional measurement: tip length (creep elongation), chord width, wall thickness via ultrasound
- Tip clearance measurement (blade elongation directly relates to disk tip clearance budget)

### Non-Destructive Evaluation (NDE)

**Fluorescent penetrant inspection (FPI):**
- Blade submerged in fluorescent dye penetrant → dye seeps into surface-breaking cracks
- Washed → developer applied → UV light reveals dye-filled cracks
- Sensitivity: detects cracks > 0.05 mm

**Thermal imaging:**
- Blade heated uniformly; thermal camera images surface
- Hot spots = regions without TBC (TBC has thermal resistance → TBC-free spot heats faster)
- Identifies TBC delamination/spallation areas before they're visible to the naked eye

**X-ray / CT scan:**
- Detects internal defects: porosity, cracks, ceramic core remnants
- Checks cooling hole integrity (blockage, connected cracks)
- Identifies casting defects that may have been masked by coating

**EBSD / Laue diffraction:**
- Checks crystallographic orientation of SX blade
- Critical if repair weld or thermal treatment might have caused recrystallization
- Any area with grain boundary detected → rejection or further evaluation

**Metallographic sectioning (destructive — used for life assessment:**
- Sample taken from a sentinel blade (same batch, same flight history)
- Measure: TGO thickness, γ′ size (coarsening = service temperature indicator), bond coat Al remaining, creep void density, SRZ depth
- Results used to calibrate remaining life models for the entire blade population

### Accept/Reject Decision

Blades are sorted into categories based on inspection:
1. **Serviceable as-is**: no structural damage, coating adequate → return to service
2. **Repairable**: damage within repair process capability → repair and return
3. **Scrap**: beyond repair capability (excessive creep, extensive cracks, SX recrystallization, etc.)

---

## 3. Stripping the Old Coating

Before recoating, old coatings must be removed completely — any residual contamination prevents proper adhesion of new coatings.

### TBC Removal

**Method 1 — Grit blasting:**
Aluminum oxide grit blasted at 4–5 bar removes the porous YSZ top coat. Controlled carefully to remove TBC without attacking bond coat.

**Method 2 — Chemical dissolution:**
Certain NaOH-based solutions selectively dissolve YSZ at elevated temperature. Used when grit blast would damage blade geometry.

### Bond Coat Removal

**Critical challenge:** Remove MCrAlY or aluminide completely without attacking CMSX-4 substrate.

**Aluminide removal (chemical):**
Proprietary acid solutions (phosphoric/chromic acid mixtures) dissolve the β-NiAl aluminide. The β phase is selectively attacked because it has different electrochemical potential from the γ/γ' alloy — selective electrochemical dissolution.

**MCrAlY removal:**
More aggressive chemical stripping or controlled grit blasting. The key: stop exactly at the bond coat/alloy interface — do not remove alloy material.

**Dimensional consequence:** Each strip/recoat cycle removes a small amount of substrate material (~10–20 μm). After 2–3 repair cycles, wall thickness has decreased by 30–60 μm → must verify wall thickness still meets minimum structural requirement.

**Useful life for re-coating:** Most blades can undergo 2–3 TBC strip/recoat cycles before wall thinning and alloy property degradation make further repair uneconomical.

---

## 4. TBC Reapplication

After stripping, the blade surface condition is different from new:
- Surface has been chemically stripped → residual contamination possible
- Near-surface alloy chemistry modified by service (Al depletion, interdiffusion)
- Surface roughness different from new blade

**Pre-coat steps:**
1. **Cleaning:** Ultrasonic cleaning in alkaline solution → solvent cleaning → dry
2. **Surface inspection:** Verify bond coat surface is free of contamination, old TGO fragments
3. **Possible diffusion anneal:** If alloy near-surface has been significantly modified, may need brief anneal to re-homogenize
4. **New bond coat application:** LPPS or HVOF MCrAlY, or new Pt-aluminide sequence
5. **Pre-oxidation:** Same as on new blade — establish α-Al₂O₃ TGO before TBC
6. **EB-PVD TBC:** Same process as new blade

**Life expectancy of repaired TBC:**
Typically 80–90% of new blade TBC life. The main limitation:
- Near-surface alloy Al depletion → bond coat Al reservoir supplemented by alloy is less
- Some TGO residuals from previous cycle → TGO/TBC interface less perfect

---

## 5. Tip Repair — Brazing and Rebuilding

Blade tips wear in service from:
- Tip rubs (when clearance closes → blade tip contacts casing)
- Oxidation/hot corrosion at the tip (high local temperature)
- Erosion

A worn tip increases tip clearance → more hot gas leaks over the tip → turbine efficiency loss → fuel burn increase.

### Tip Braze Repair

For tips worn 0.5–3 mm:
1. Machine the worn tip surface flat (remove damaged material)
2. Apply braze alloy (nickel-base braze, e.g., BNi-2: Ni-7Cr-4.5Si-3B-3Fe) as paste
3. Align a CMSX-4 or nickel superalloy tip cap on top
4. Vacuum braze at 1,000–1,100°C → braze alloy melts → flows into joint → solidifies → bonds tip cap to blade

**Brazing challenges:**
- Braze temperature must not exceed CMSX-4 γ′ solvus (>1,316°C → γ′ dissolves) but must fully melt braze alloy
- B and Si in braze alloy (used as melt-point depressants) — these can diffuse into the single crystal → local grain nucleation (Si promotes stacking faults)
- Joint must withstand centrifugal load equivalent to ~5,000 N at 10,500 RPM

**Modern development:** "Wide-gap brazes" using powder mixes with controlled melt point can fill larger gaps (up to 1.5 mm) with good mechanical properties.

### Tip Welding

For severely damaged tips or multiple tip cap loss: direct laser welding of CMSX-4 filler material.

Challenge: CMSX-4 is NOT weldable by conventional means — see §6.

Some advanced repair shops use laser powder deposition (LPD/LENS process) with special modified alloy powder (lower Re, more ductile, designed for weldability) for tip build-up. The deposited zone has a polycrystalline microstructure — this is acceptable at the tip (lower stress region, blade root carries most centrifugal load).

---

## 6. Crack Repair — Welding Challenges in SX

### Why SX Blades Cannot Be Conventionally Welded

Welding involves:
1. Melting the base metal
2. Re-solidification from the melt

In a polycrystalline metal, solidification from welds nucleates many grains → grain boundaries form → this is acceptable (weld microstructure is always polycrystalline).

In a SINGLE CRYSTAL blade, any re-solidification from a weld pool:
- Creates multiple new grains (random nucleation in molten pool)
- These grains form grain boundaries AT the repair location
- A grain boundary in a single-crystal blade is a DEFECT → reduces creep life → may require blade rejection

Furthermore, CMSX-4 is susceptible to **hot cracking** (solidification cracking): as the weld pool solidifies, the last-to-solidify interdendritic regions (rich in low-melting eutectics) are in tension → crack formation during solidification. This is characteristic of alloys with a wide freezing range (ΔT₀ large).

**Summary: CMSX-4 is not weldable in the classical sense for structural locations.**

### Brazing (for Cracks < 0.5 mm)

Very tight cracks (< 0.1–0.5 mm) can be filled by brazing. High-temperature vacuum braze flows into the crack by capillary action. Acceptable for:
- Non-structural crack locations (blade tip, platform edges)
- Low-stress areas
- Cracks not crossing the primary stress direction

### Activated Diffusion Healing (ADH) / Transient Liquid Phase (TLP) Bonding

A more advanced repair for structural cracks in SX blades:

1. Special repair paste applied into the crack: base-alloy powder + melt-point depressant (B, Si, or Mn to depress liquidus to 1,100–1,200°C)
2. Vacuum heat treatment at 1,100–1,250°C → paste melts → liquid fills crack by capillary action
3. The melt-point depressants diffuse away into the base alloy over time → the filler composition evolves toward the base alloy
4. After diffusion (several hours): filler is now same composition as base alloy → solidus temperature matches base alloy → filler is fully solid at service temperature
5. If conditions are controlled (low undercooling, very slow solidification) → epitaxial solidification on SX walls → no grain boundaries formed

This is the gold standard for SX crack repair. Success rate: ~70–80% for narrow cracks (< 0.3 mm) with careful process control.

---

## 7. Recrystallization — The Enemy of Repair

**Recrystallization** occurs when a heavily deformed SX blade region is heated above the recrystallization temperature (~950°C for CMSX-4). The deformed single crystal structure "resets" → new polycrystalline grains nucleate and grow → grain boundaries form in what was a perfect crystal.

**Triggers for recrystallization during repair:**
1. **Mechanical damage during cleaning:** Grit blasting too aggressively → work-hardening near surface → deformed surface layer → recrystallizes on subsequent heating
2. **Vibratory polishing:** Creates surface damage layer
3. **Straightening operations:** If blade is slightly bent from creep → straightening introduces local plastic strain → recrystallization during subsequent heat treatment
4. **Machining during tip repair:** Cutting/grinding the tip → subsurface plastic deformation layer

**Prevention:**
- Chemical rather than mechanical cleaning wherever possible
- Minimal abrasive action
- If mechanical treatment is necessary: use gentle methods (hand polishing with 600-grit SiC, not power-blasting)
- Verify SX integrity by Laue diffraction before re-coating

**Detection:**
- Metallographic examination of stripped blade surface
- EBSD mapping shows any recrystallized grains (identified by orientation change)
- Any recrystallized grain at a structural location → blade must be scrapped or repaired with TLP bonding

---

## 8. Diffusion Bonding for Segment Replacement

For heavily damaged sections (e.g., the entire leading edge eroded, a large section oxidized away), it may be possible to replace a segment:

1. Machine out the damaged section (e.g., remove 10 mm from the tip + 15 mm of leading edge)
2. Prepare a new replacement segment (same alloy, same crystallographic orientation — seed crystal grown to match)
3. Diffusion bond the replacement segment to the remaining blade:
   - Align crystal orientations within 2–3° using Laue diffraction
   - Apply bonding "interlayer" (very thin layer of lower-melting alloy)
   - Hold in vacuum press at 1,200°C, ~10 MPa for 4–6 hours → diffusion bond formed

**Challenge:** Matching the crystallographic orientation of the new segment to the existing blade within acceptable tolerance (< 5° misorientation at bond line). This requires specialized orientation measurement and alignment fixtures.

**Research status:** Successful in laboratory for model alloys. Limited production application. Most heavily damaged blades are scrapped rather than segment-repaired.

---

## 9. Re-Solution and Re-Aging After Repair

Any thermal repair process affects the γ/γ′ microstructure:
- Braze heating (1,100°C): partial γ′ dissolution → re-precipitation occurs but distribution will be non-optimal
- If braze temperature approaches γ′ solvus (~1,316°C for CMSX-4): significant γ′ dissolution → coarsened γ′ on re-cooling

For best mechanical properties after repair, re-solution and re-aging is desirable:

**Standard re-solution + re-age:**
1. **Solution treat**: 1,290–1,310°C, 4h → dissolves any coarsened γ′, equalizes chemistry
2. **Primary age**: 1,140°C, 6h → fine γ′ precipitates
3. **Secondary age**: 870°C, 20h → bimodal γ′ distribution

**Challenge:** The coating (if already applied) may not survive these temperatures. Options:
1. Re-solution BEFORE coating → then coat → but any repair welding done before coating may suffer from the thermal treatment
2. Re-solution at lower temperature (< γ′ solvus) → partial recovery only
3. Accept suboptimal γ′ microstructure if repair braze temperature was low

In practice, for blades with TBC already applied, the heat treatment is omitted or kept to the minimum needed for braze flow.

---

## 10. Repair Qualification and Return to Service

Aviation authority (FAA/EASA) regulations require:
- Each repair method must be qualified by the Original Equipment Manufacturer (OEM — GE, R-R, P&W) or an FAA-approved Designated Engineering Representative (DER)
- Qualification involves extensive testing: fatigue specimens from repaired blades, creep tests, engine test articles with repaired blades
- Repair must restore blade to "Serviceable Limits" — defined in the Engine Maintenance Manual

**Post-repair inspection:**
- Full visual + FPI after all repair processes
- Dimensional verification (wall thickness, tip length, chord)
- X-ray (braze voids, crack-fill quality)
- Laue diffraction if any thermal process could have caused recrystallization
- Thermal imaging to verify new TBC quality

**Documentation:**
Each blade repair generates a traveler document tracking every operation, every measurement, every material used. This traceability is required by aviation regulations and ties back to the blade's serial number and its entire history.

---

## 11. Life Extension Strategies

Beyond standard repair, some advanced strategies can extend blade life beyond nominal limits:

### Enhanced Inspection Intervals

If advanced inspection techniques can detect damage earlier and with better sensitivity (e.g., phased array ultrasound, X-ray tomography at overhaul), blades can be returned to service with greater confidence → fewer unnecessary scraps → more efficient fleet utilization.

### Re-Sensitization Heat Treatment

A controlled heat treatment after strip-coating (without new coating) can partially recover creep-depleted microstructure by re-precipitating fine γ′ and allowing some recovery of dislocation density. Not a full rejuvenation, but extends second-life useful period.

### Condition-Based Maintenance (CBM)

Rather than fixed overhaul intervals, modern engines carry sensors that record:
- Metal temperature (pyrometers in the engine)
- Cycle counts with severity weighting (longer, hotter cycles consume more life)
- Vibration levels (HCF risk indicator)
- Pressure ratios (indicates blade fouling)

This flight-by-flight data feeds life prediction models → each blade gets an individual, calculated remaining life → only blades actually approaching their limit are pulled for inspection.

Airlines using CBM can extend average time-on-wing by 10–20% vs. fixed-interval maintenance, while maintaining safety margins.

---

## Summary

- **Repair economics**: repair costs $2,000–8,000 vs. $30,000+ new; 60–70% of blades are repairable → $50M+/year savings for a large airline fleet
- **Inspection determines repairability**: FPI, thermal imaging, X-ray, Laue diffraction identify damage type and location
- **Stripping**: chemical removal of coatings without damaging SX substrate; 2–3 cycles possible before wall thinning is excessive
- **TBC reapplication**: same process as new blade, ~80–90% of new life restored
- **Tip braze repair**: vacuum braze with nickel-base braze alloy; CMSX-4 tip cap; handles 0.5–3mm wear
- **Crack repair**: conventional welding not possible (creates grain boundaries in SX). TLP bonding → epitaxial re-solidification → no grain boundary → the only structural SX crack repair
- **Recrystallization**: the key threat to SX integrity during any repair operation; prevented by avoiding mechanical damage; detected by Laue diffraction
- **Re-solution + re-aging**: restores γ′ microstructure after thermal repair processes; may conflict with coating sequence
- **CBM**: flight-data-driven individual blade life management → 10–20% time-on-wing extension

---

## Exercises

1. A blade has experienced 15,000 EFH with TGO thickness now at 6 μm. Strip/recoat repairs: (a) the new TGO starts at zero after recoating. If TGO grows at rate h² = 0.004t (μm², h), how long until TBC spallation (TGO = 8 μm) after the repair? (b) How does this compare to first-life TBC life? (c) After 2 strip/recoat cycles, the bond coat Al has dropped from 10% to 7%. If below 4% means insufficient protection, how many more cycles are possible?

2. TLP bonding: a crack in CMSX-4 is 0.15 mm wide. Repair paste has 80% CMSX-4 powder + 20% boron-containing depressant (2% B in Ni). During healing, B diffuses out of the filler at 1,180°C with D_B = 2×10⁻¹³ m²/s. (a) After 6 hours at 1,180°C, how far has B diffused (use x = 2√(Dt))? (b) If B concentration drops from 2% to < 0.1% over distance d, is this condition met after 6 hours? (c) Why must B be removed for good mechanical properties?

3. Recrystallization risk: a grit-blast operation creates a surface deformation layer to depth δ = 15 μm. The recrystallization temperature of deformed CMSX-4 is ~900°C. During pre-oxidation at 1,050°C: (a) will the surface layer recrystallize? (b) If recrystallization produces grains 20 μm in diameter, how many grain boundaries would span the 15 μm depth? (c) A subsequent solution treatment at 1,300°C — above recrystallization T but still below solidus — runs for 4 hours. Will the recrystallized zone grow? Use grain boundary migration data: v = M × P where M ≈ 10⁻¹² m/(Pa·s) at 1,300°C and P (driving pressure) ≈ 2γ/r ≈ 10⁵ Pa.

4. Economic analysis: an airline has 250 engines each with 80 HPT blades. Average blade cost is $25,000 (new). Shop visit interval: 5,000 EFH. At each visit: 65% of blades are repaired (cost $4,000 each, extends life 5,000 EFH), 25% are scrapped and replaced ($25,000 each), 10% are returned as-is. Calculate: (a) cost per shop visit per engine, (b) cost per EFH of operation for HPT blades, (c) if repair success rate improves to 75%, how much does this save per engine per year (assume 5,000 EFH/year)?

5. Tip braze fillet stress: A new CMSX-4 tip cap (5 mm tall, 20 mm × 5 mm cross-section, density 8.7 g/cm³) is brazed to a blade rotating at 10,200 RPM with the tip at r = 500 mm from the rotation axis. (a) Calculate the centrifugal force on the tip cap. (b) The braze joint area is 20 mm × 2 mm (the braze fillet perimeter cross-section). What shear stress acts on the braze? (c) Braze joint shear strength (at service temperature, 850°C) ≈ 150 MPa. Is the joint adequate? What safety factor does this give?
