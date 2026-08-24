# Chapter 60: The Evolution of Superalloys — Six Generations of Progress

> **"The story of superalloy evolution is one of the great engineering sagas of the twentieth century. In 1940, a turbine blade made of Nimonic 80 could survive 700°C. Today, a CMSX-4 single crystal survives 1,050°C metal temperature — while the gas around it is at 1,650°C. That 350°C improvement in metal temperature capability represents 80 years of work by thousands of metallurgists, crystallographers, thermodynamicists, and engineers. Each step was a new negotiation between competing demands: stronger vs. stable, creep-resistant vs. castable, oxidation-resistant vs. dense."**

---

## Table of Contents

1. [The Framework: What We're Measuring](#1-the-framework-what-were-measuring)
2. [Pre-Superalloy Era — Stainless Steels and Early Ni Alloys](#2-pre-superalloy-era--stainless-steels-and-early-ni-alloys)
3. [Polycrystalline Era: 1940s–1960s](#3-polycrystalline-era-1940s1960s)
4. [DS Revolution: 1966–1982](#4-ds-revolution-19661982)
5. [First-Generation SX: 1982–1985](#5-first-generation-sx-19821985)
6. [Second-Generation SX: Re Addition 1985–1995](#6-second-generation-sx-re-addition-19851995)
7. [Third-Generation SX: High-Re 1990s–2000s](#7-third-generation-sx-high-re-1990s2000s)
8. [Fourth-Generation SX: Ru Stabilization 2000s](#8-fourth-generation-sx-ru-stabilization-2000s)
9. [Fifth and Sixth Generation: Frontier Materials 2010–Present](#9-fifth-and-sixth-generation-frontier-materials-2010present)
10. [The Density Challenge](#10-the-density-challenge)
11. [Alloy Design Lessons — Patterns Across Generations](#11-alloy-design-lessons--patterns-across-generations)
12. [Summary Table — All Generations](#12-summary-table--all-generations)
13. [Summary](#summary)
14. [Exercises](#exercises)

---

## 1. The Framework: What We're Measuring

The single most important metric in turbine blade alloy development is:

**Specific creep life** = temperature at which the alloy sustains a given stress (typically 137–200 MPa) for a given time (typically 1,000 or 100 hours) without rupture.

This is measured as **stress rupture temperature at 1,000 hours, 137 MPa** — the higher, the better.

A secondary metric: **specific strength** = strength/density. Denser alloys create higher centrifugal stresses, so raw strength at density is what matters.

Also tracked: oxidation life, fatigue life, phase stability, castability, cost.

---

## 2. Pre-Superalloy Era — Stainless Steels and Early Ni Alloys

**1930s:** Turbine blades were forged from austenitic stainless steels (18-8, 310SS).
- T_max ≈ 650°C (blade metal temperature)
- Limited by rapid creep above 600°C, low oxidation resistance
- Frank Whittle's first jet engine (1937) used forged steel blades — they lasted only hours before creep failure

**Nimonic 75 (Mond Nickel/Rolls-Royce, 1940):**
- Ni-20Cr, small Ti
- First specifically designed jet engine alloy
- T_max ≈ 700°C
- Used in early Whittle W.2/700 and de Havilland Goblin

The key breakthrough: adding Ti and Al to Ni-Cr base → discovered the precipitation hardening by γ′ phase → this was the beginning of the superalloy concept.

---

## 3. Polycrystalline Era: 1940s–1960s

Each increment of this period achieved 20–40°C improvement in temperature capability:

### Nimonic 80A (1943)
**Composition:** Ni-19.5Cr-2Ti-1.5Al
- First deliberate γ′ strengthening alloy
- T_max ≈ 730°C metal temperature
- Used in de Havilland Ghost, Rolls-Royce Avon
- Al + Ti → γ′ (Ni₃(Al,Ti)) forms during aging → dramatic creep improvement vs. solid solution alloys

### Nimonic 90 (1945)
**Composition:** Ni-20Cr-16Co-2.5Ti-1.5Al
- Co addition: first use of Co in superalloy (Co reduces γ′ solvus slightly but improves hot corrosion resistance and allows more γ′ to form → strength)
- T_max ≈ 760°C

### Waspaloy (Pratt & Whitney, 1952)
**Composition:** Ni-19.5Cr-13.5Co-4.3Mo-3Ti-1.4Al + B, Zr
- Introduction of Mo as solid solution strengthener (Mo in γ matrix → strong dislocation drag)
- B and Zr at grain boundaries → grain boundary strengthening → ductility improvement
- T_max ≈ 820°C
- Still used today for disk applications (not blades)

### Mar-M200 (Martin Metals, 1960s)
**Composition:** Ni-9Cr-10Co-12.5W-5Al-2Ti-1Nb + Hf, B, Zr
- Major increase in Al (5%) → much more γ′ → volume fraction approaches 40–50%
- W replacing Mo: similar solid solution strengthening but higher melting point
- Hf: grain boundary ductilizer (Hf-rich carbides at GB)
- T_max ≈ 870°C

### IN100 (International Nickel / Pratt & Whitney, 1963)
**Composition:** Ni-10Cr-15Co-3Mo-5.5Al-4.7Ti-1V + B, Zr
- γ′ volume fraction ~60% (very high for PC alloy)
- V addition: reduces stacking fault energy → extends range of glide → useful for some deformation-based fatigue resistance
- T_max ≈ 900°C
- The high-water mark of polycrystalline investment-cast superalloys

### What Limited PC Alloys

By ~1960, polycrystalline alloy design was approaching fundamental limits:
- More γ′ → alloy becomes uncastable (freezing range too wide, mushy zone forms)
- More W, Mo → density too high, TCP phases form
- More Cr → alloy unstable toward σ phase
- All routes were blocked by the grain boundary failure mechanism — no matter how strong the grain interior, transverse GBs failed first at >900°C

The only path forward: change the microstructure, not just the composition. → Directional solidification.

---

## 4. DS Revolution: 1966–1982

**First DS blade in service: Pratt & Whitney JT9D, 1966**

Key alloy: **Mar-M200 DS** (same composition as PC version, but directionally solidified):
- No NEW alloy — same composition as Mar-M200 equiaxed
- +30°C temperature capability PURELY from microstructure change (elimination of transverse grain boundaries)
- T_max ≈ 960°C

**Lesson:** The ~30°C improvement from DS to SX (comparable to adding an entirely new refractory element) came from PROCESS innovation, not alloy innovation. This demonstrated that processing and microstructure could be as powerful as composition.

**Early DS-specific alloys (1970s):**

**CM247LC (Cannon-Muskegon, 1971):**
- Designed SPECIFICALLY for DS (PC alloys weren't optimized for DS casting)
- Reduced grain boundary elements (B, Zr, C lowered to enable DS casting without hot cracking)
- Optimized for columnar grain formation
- T_max ≈ 980°C DS

**PWA 1422:**
- Pratt & Whitney DS alloy
- Reduced Hf (Hf caused hot cracking in DS if too high) → adjusted composition for DS castability
- T_max ≈ 970°C

### What DS Changed in Alloy Design

With DS, grain boundary strengtheners (B, C, Zr) could be reduced because transverse boundaries were eliminated:
- Freed up composition "space" for other elements
- Slightly raised solidus → better heat treatment window
- Set the stage for the next step: eliminate boundaries entirely

---

## 5. First-Generation SX: 1982–1985

**First SX blade in commercial service: PWA 1480 in Pratt & Whitney F100 (military), 1982**

**PWA 1480:**
- **Ni-10Cr-5Co-4W-5Al-1.5Ti-12Ta** — remarkable for what's ABSENT: no B, C, Zr (grain boundary strengtheners removed entirely)
- No Re yet (1st generation predates Re)
- T_max ≈ 1,010°C (vs. ~960°C for best DS alloys)

**Why 1st gen SX gave +30–50°C over DS (same alloy):**
- No transverse OR longitudinal grain boundaries
- No grain boundary elements → solidus increased by ~40°C → solution treatment 40°C hotter → more complete γ′ dissolution → better homogenization
- The higher solidus was UNEXPECTED bonus of removing C/B/Zr from the alloy

**SRR99 (Rolls-Royce, 1984):**
- **Ni-8Cr-5Co-10W-5.5Al-2.2Ti-3Ta** — high W for solid solution strengthening
- First major UK SX alloy, designed for Rolls-Royce Spey engine upgrade
- T_max ≈ 1,005°C

**René N4 (GE, 1983):**
- **Ni-9.75Cr-7.5Co-1.5Mo-6W-3.7Al-3.7Ti-4Ta-0.5Nb** 
- GE's first SX blade alloy for CF6/F110
- Note: contains Nb (later replaced by Ta in 2nd gen)

**What 1st generation revealed:**

1. **SX advantage confirmed**: The 30–50°C improvement over DS was more than expected
2. **W limitation**: High W is heavy (W = 19.3 g/cm³) → density limit of ~8.5 g/cm³
3. **No Re needed... yet**: 1st gen achieved remarkable performance without Re. But the limit of W and Ta was being approached.
4. **Secondary orientation matters**: Researchers discovered that secondary crystal orientation (rotation around [001]) affects thermal fatigue life — this was not anticipated from the equiaxed alloy experience

**Research breakthrough during 1st gen era:**

Discovery that Re could substitute for some W with MUCH better creep benefit-to-density ratio → led directly to 2nd generation.

---

## 6. Second-Generation SX: Re Addition 1985–1995

The 2nd generation represents the single most important compositional innovation in superalloy history: the deliberate addition of rhenium (Re).

### CMSX-4 (Cannon-Muskegon, 1985)
**Composition (wt%):** Ni-6.5Cr-9.6Co-0.6Mo-6.4W-5.6Al-1Ti-6.5Ta-3Re-0.1Hf

Key changes from 1st gen (vs. CMSX-2 as reference):
- **Re: 3 wt%** (new element, was absent in 1st gen)
- W reduced from 8% to 6.4% (Re partially replaces W, less dense but more effective)
- Mo kept low (0.6%) — TCP phase concerns with high Re+Mo
- Hf added (0.1%) — oxidation scale adhesion improvement

**T_max:** ~1,020°C (+30°C vs. 1st gen at equivalent alloy/equivalent stress)

**Mechanisms of Re improvement (as discussed in Chapter 43):**
1. Potent solid solution strengthener (dislocation drag)
2. Retards dislocation climb (Re clouds on dislocations)
3. Suppresses γ′ coarsening (Re diffusion barrier in γ channels)

### PWA 1484 (Pratt & Whitney, 1988)
**Composition:** Ni-5Cr-10Co-2Mo-6W-5.6Al-9Ta-3Re
- Less Cr than CMSX-4 (5% vs 6.5%) → accepted reduced hot corrosion resistance in exchange for stability
- Higher Ta (9%) → very strong γ′ formers
- T_max ≈ 1,020°C

### René N5 (GE, 1990)
**Composition:** Ni-7Cr-7.5Co-2Mo-5W-6.2Al-7Ta-3Re-0.15Hf
- Higher Cr (7%) than PWA1484 → better hot corrosion balance
- Standard GE HPT blade alloy through GE90, CF6-80 programs
- T_max ≈ 1,020°C

### What 2nd Generation Revealed

**Re poisoning effect (Density):** Re is very dense (21.0 g/cm³). Adding 3% Re increases alloy density from ~8.5 to ~8.7 g/cm³. This 2% density increase means 2% more centrifugal stress → not ideal, but acceptable.

**Re TCP precipitation (future concern identified):** At 3% Re, TCP phases were manageable with existing Cr and Co levels. But preliminary experiments with 6% Re (what would become 3rd gen) showed unacceptable TCP formation — this was recognized as the challenge for future generations.

**Re supply vulnerability:** As Re usage increased, supply chain concerns emerged. Re is a byproduct of Cu/Mo mining — supply is inherently volatile.

---

## 7. Third-Generation SX: High-Re 1990s–2000s

The 3rd generation pursued the logical next step: if 3% Re gives +30°C, what does 6% Re give?

### CMSX-10 (Cannon-Muskegon, 1994)
**Composition:** Ni-2Cr-3Co-0.4Mo-5W-5.7Al-0.2Ti-8Ta-6Re-0.03Hf

Key changes:
- **Re: 6 wt%** (doubled from 2nd gen)
- **Cr: 2 wt%** (dramatically reduced — Cr promotes TCP formation when combined with high Re)
- **Co: 3 wt%** (also reduced — Co stabilizes some harmful phases with high Re)

**T_max:** ~1,045°C (+25°C vs. 2nd gen CMSX-4)

**Problems discovered:**
CMSX-10's very low Cr (2%) made it vulnerable to hot corrosion. More critically, the combination of 6% Re + low Cr + low Co created severe **TCP phase precipitation** during service at 900–1,050°C. σ and P phases nucleated in the γ matrix → embrittlement → premature blade retirement.

### René N6 (GE, 1994)
**Composition:** Ni-4.2Cr-12.5Co-1.4Mo-6W-5.75Al-7.2Ta-5.4Re-0.15Hf

N6 attempted to manage TCP with higher Co (12.5%) as a stabilizer — higher Co reduces the tendency for TCP. This worked somewhat better than CMSX-10 at maintaining hot corrosion resistance.

### The 3rd Generation Lesson

**3rd generation achieved the performance target (+25°C over 2nd gen) but at a cost:**
- TCP phase stability degraded → blades needed more frequent inspection
- Hot corrosion resistance compromised (low Cr)
- Density increased further (Re = 21 g/cm³ → 6% Re makes alloy significantly denser)
- Manufacturing yield decreased (higher Re → wider freezing range → more stray grain risk)

**The fundamental problem:** You can't keep adding Re indefinitely. Beyond ~6% Re, TCP formation becomes severe with ALL conventional Cr/Co/W/Mo combinations. A new approach was needed.

This drove the discovery of Ru (ruthenium) as the TCP stabilizer.

---

## 8. Fourth-Generation SX: Ru Stabilization 2000s

**Key discovery (2001–2003, multiple research groups simultaneously):**

Ruthenium (Ru) additions to high-Re alloys dramatically reduce TCP precipitation. Ru:
- Has a specific electron configuration (4d⁷5s¹ → filled d-band effects) that shifts alloy electron-per-atom ratio away from TCP stability field
- Does NOT form TCP phases itself (in contrast to Re, Mo, W)
- Allows retention of high Re content without TCP problems

### TMS-138 (Japan NIMS, 2003)
**Composition:** Ni-3.2Cr-5.8Co-2.8Mo-5.6W-5.8Al-0.2Ti-5.6Ta-5.6Re-2Ru

The first 4th-generation alloy with significant Ru:
- 2% Ru: stabilizes against TCP at 5.6% Re
- Higher Mo than typical (2.8%) — allowed because Ru suppresses Mo's TCP tendency
- T_max ≈ 1,060°C

### EPM-102 (US NIMS/NASA, 2003)
**Composition:** Ni-2Cr-16.5Co-2.9Mo-6W-5.9Al-8.3Ta-5.9Re-3Ru

Very high Co (16.5%) combined with Ru → particularly effective TCP suppression.

### TMS-162 (Japan NIMS, 2005)
4th gen optimized: higher Ru (6%) enables 6% Re without TCP concerns:
- T_max ≈ 1,070°C
- Density: ~8.95 g/cm³ (heavier due to Re + Ru)

### 4th Generation Challenges

**Ru is more expensive than Re:** Ru ≈ $14,000–20,000/kg vs. Re ≈ $3,500/kg. Adding 2–6% Ru to a blade that already has 5–6% Re → very high alloying element cost.

**Density continues to increase:** Re (21.0) + Ru (12.4) → both dense elements → alloy density → 8.9–9.1 g/cm³ for some 4th gen alloys. This is approaching practical limits for centrifugal stress.

---

## 9. Fifth and Sixth Generation: Frontier Materials 2010–Present

### TMS-238 (Japan NIMS, 2010s)
**Composition:** Ni-4.6Cr-6.5Co-2.3Mo-3.5W-5.9Al-0.1Ti-5.8Ta-6.4Re-5Ru

Key optimizations vs. TMS-162:
- More Cr (4.6 vs 3.2%) → better oxidation/hot corrosion balance
- Lower W (3.5 vs 5.6%) → lower density, offset by more Re
- More Ru (5%) → better TCP stability
- T_max ≈ 1,080–1,100°C (based on NIMS publications)
- Density: ~8.92 g/cm³ (well-controlled despite high Re+Ru)

TMS-238 is currently the highest-reported-performance SX alloy in peer-reviewed literature.

### CMSX-10K / CMSX-12 Family (Cannon-Muskegon, 2010s)

CMSX-10K: modified CMSX-10 with improved TCP stability (undisclosed exact composition — proprietary). Used in some military programs.

### RR3000 (Rolls-Royce, in development)

Details proprietary; known to be 4th+ generation SX with Ru.

### STAL-15 (China AECC, 2010s)

China's indigenous 4th-generation SX development. Significant investment in domestic alloy capability to reduce dependence on foreign alloys for COMAC C919 and military aircraft.

---

## 10. The Density Challenge

As alloys have added more refractory elements (Re, W, Ru), density has increased:

| Generation | Example | Density (g/cm³) |
|------------|---------|-----------------|
| 1st gen | PWA 1480 | 8.45 |
| 2nd gen | CMSX-4 | 8.70 |
| 3rd gen | CMSX-10 | 9.05 |
| 4th gen | TMS-138 | 8.95 |
| 5th/6th gen | TMS-238 | 8.92 |

**Why density matters:**
```
Centrifugal blade root stress = ρ × ω² × r × L

For 5% density increase (8.45 → 8.87 g/cm³):
Stress increase = 5% → blade root stress goes from 200 MPa to 210 MPa
```

This 10 MPa increase must be borne by the alloy. At 980°C and 210 MPa, the required creep life is slightly more than at 200 MPa — meaning the alloy's improved creep resistance is partially offset by its own increased density.

**Density-normalized creep performance** (strength/density) is the true metric of progress. 4th-generation alloys are better than 3rd-generation on this metric because Ru's density (12.4 g/cm³) is much lower than Re (21.0 g/cm³), and Ru allows Re content optimization.

**Future direction:** Interest in lower-density additions that can replace some Re/W:
- Molybdenum (Mo, 10.2 g/cm³) — lower than Re, but TCP concerns
- Co (8.9 g/cm³) — comparable to Ni, beneficial for other reasons
- Reduced W — already being pursued in latest generation
- Ru continues to be the key stabilizer enabling this density reduction while maintaining Re

---

## 11. Alloy Design Lessons — Patterns Across Generations

Looking across 80 years of alloy development, clear patterns emerge:

### Lesson 1: Process Innovation ≥ Composition Innovation

DS (+30°C) and SX (+30°C) together gave +60°C — comparable to ALL alloy composition improvements from 1st to 4th generation combined. Process breakthroughs are high-leverage.

### Lesson 2: Each New Element Brings New Problems

- Add more W → denser, TCP prone
- Add Re → denser, TCP prone, more expensive
- Add Ru → less TCP, but $14,000/kg
- Every "fix" creates a new constraint → alloy design is an evolving negotiation

### Lesson 3: The Rate of Improvement Is Slowing

| Period | Improvement |
|--------|-------------|
| PC → DS (1966) | +30°C |
| DS → 1st SX (1982) | +30°C |
| 1st → 2nd gen (1985) | +30°C |
| 2nd → 3rd gen (1994) | +25°C |
| 3rd → 4th gen (2003) | +15°C |
| 4th → 5th/6th gen (2010s) | +15–20°C |

Each increment takes more investment, more time, and more cost for less temperature improvement. This is the classic S-curve of technology maturation.

### Lesson 4: Systems Thinking Beats Isolated Property Optimization

The coating system, cooling architecture, and alloy composition must be co-optimized. A blade alloy with perfect creep resistance but poor coating compatibility is useless. The system performance (TIT capability) depends on all three subsystems.

### Lesson 5: Computational Design Has Become Essential

1st and 2nd generation alloys were designed by incremental empiricism (make alloy, test, iterate). 3rd generation onward: CALPHAD thermodynamic modeling was essential for predicting TCP stability in 10+ element systems. 4th generation: CALPHAD + DFT (density functional theory) predictions of Ru's TCP-suppression mechanism preceded experimental confirmation.

---

## 12. Summary Table — All Generations

| Gen | Year | Example | Key Feature | T_max (°C) | Density | Key Elements |
|-----|------|---------|-------------|------------|---------|-------------|
| PC (early) | 1940 | Nimonic 80 | First γ′ | ~730 | 8.2 | Ni-Cr-Ti-Al |
| PC (advanced) | 1963 | IN100 | 60% γ′ | ~900 | 8.5 | +Co, Mo, V |
| DS | 1966 | Mar-M200DS | Columnar grains | ~960 | 8.5 | Same PC + process |
| 1st SX | 1982 | PWA 1480 | No grain boundaries | ~1,010 | 8.45 | No B,C,Zr; high Ta |
| 2nd SX | 1985 | CMSX-4 | 3% Re | ~1,020 | 8.70 | 3% Re |
| 3rd SX | 1994 | CMSX-10 | 6% Re | ~1,045 | 9.05 | 6% Re, low Cr |
| 4th SX | 2003 | TMS-138 | Re + Ru | ~1,060 | 8.95 | Re + 2–3% Ru |
| 5th/6th SX | 2010s | TMS-238 | Optimized Re/Ru | ~1,080+ | 8.92 | Re + 5% Ru |

---

## Summary

- **80 years of development**: from 730°C (Nimonic 80) to 1,080°C+ (TMS-238) metal temperature capability — +350°C total improvement
- **Three eras of breakthrough**: (1) polycrystalline alloys maximizing γ′, (2) DS/SX removing grain boundaries, (3) Rhenium addition (2nd gen) + Ruthenium stabilization (4th gen)
- **Re is the single most important compositional advance**: 3 wt% Re gave +30°C in 1985 (2nd generation) — still unmatched by any other single element
- **Ru solved the 3rd-generation TCP problem**: enabled retention of high Re levels → paved the way for 4th–6th generation
- **Rate of improvement slowing**: each new generation costs more for fewer degrees of improvement → driving research into radically different materials (CMC, refractory HEA, etc.) covered in Chapters 61–65
- **Process remains as important as composition**: coating systems and cooling architecture co-develop with alloy improvements; blade temperature capability = alloy + coating + cooling combined

---

## Exercises

1. Temperature capability improvement rate: From the generation data in §12, calculate the average °C improvement per generation for: (a) PC to 1st SX, (b) 1st SX to 4th SX, (c) 4th SX to 6th SX. Is the trend accelerating or decelerating? At the current rate, what temperature capability might a hypothetical 8th generation alloy achieve in 2035?

2. The Larson-Miller parameter (LMP = T(log t_r + 20)) determines creep life. CMSX-4 (2nd gen) has LMP capability of 33,000 at 200 MPa. TMS-238 (6th gen) has LMP capability of 35,000 at 200 MPa. (a) For a blade operating at 980°C (1253 K) and 200 MPa, calculate life expectancy for each alloy in hours. (b) By how many hours has 6th gen improved blade life over 2nd gen at these conditions?

3. Re economics: CMSX-4 has 3 wt% Re, TMS-238 has 6.4 wt% Re + 5 wt% Ru. For a blade of mass 0.5 kg: (a) calculate Re mass in each blade type, (b) calculate Ru mass in TMS-238 blade, (c) at Re = $3,500/kg and Ru = $15,000/kg, calculate the alloying element cost for Re+Ru in each alloy, (d) what fraction of a $30,000 finished blade cost is attributable to Re+Ru alone?

4. Density constraint: An engine rotor is designed for blade root stress of 220 MPa maximum. Currently using CMSX-4 blades (density 8.70 g/cm³). A designer wants to switch to TMS-238 (density 8.92 g/cm³). (a) By what percentage does root stress increase? (b) Is the new stress still within the 220 MPa design limit? (c) If the engine disk must be redesigned to accommodate the new stress, which direction does the change go (more or less disk mass)?

5. CALPHAD prediction: A new alloy is proposed: Ni-5Cr-8Co-2Mo-5W-5.5Al-7Ta-5Re-3Ru. Using the alloy generation table as a guide and what you know about each element's role: (a) identify which generation this most closely resembles, (b) predict whether TCP phases are likely (compare Re+Mo+W+Cr content to known stable vs. unstable compositions), (c) estimate the density using weighted average of element densities, (d) what would you change to improve TCP stability further without sacrificing temperature capability?
