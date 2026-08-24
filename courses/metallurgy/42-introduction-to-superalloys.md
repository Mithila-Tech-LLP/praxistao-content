# Chapter 42: Introduction to Superalloys — Where Normal Alloys Fail

> **"A superalloy is a material that exists at the edge of the possible. It must be strong where most metals are soft, stable where most phases transform, resistant to oxidation in conditions that destroy ordinary metals — and it must maintain these properties for 30,000 hours of operation. Nothing else in materials science comes close to this challenge."**

---

## Table of Contents

1. [What Is a Superalloy?](#1-what-is-a-superalloy)
2. [Why Normal Alloys Fail at High Temperature](#2-why-normal-alloys-fail-at-high-temperature)
3. [The Three Families of Superalloys](#3-the-three-families-of-superalloys)
4. [A Brief History — From Nimonic to TMS-238](#4-a-brief-history--from-nimonic-to-tms-238)
5. [Temperature Capability — The Evolution](#5-temperature-capability--the-evolution)
6. [Where Are Superalloys Used?](#6-where-are-superalloys-used)
7. [The Key Properties Required](#7-the-key-properties-required)
8. [Cost and Strategic Importance](#8-cost-and-strategic-importance)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What Is a Superalloy?

A **superalloy** is a high-performance alloy designed to operate at temperatures above about 540°C (1000°F), where it must maintain:
- High mechanical strength (especially creep and fatigue resistance)
- Good surface stability (resistance to oxidation and hot corrosion)
- Phase stability of the microstructure over long service lives (thousands of hours)

The name "superalloy" reflects that these materials transcend the performance limits of ordinary alloys. The three base elements are **nickel, cobalt, and iron** — but nickel-base superalloys dominate high-temperature aerospace applications.

**How "super" are they?**

A modern single-crystal Ni superalloy (e.g., CMSX-4):
- Maintains 700–800 MPa creep strength at 980°C
- Operates in gas turbine environments above its own melting point (with cooling)
- Resists oxidation at up to 1100°C (with coating)
- Provides ~30,000 hours of turbine blade service life

Compare to a good structural steel (4340):
- At 600°C, yield strength drops to ~400 MPa
- At 700°C, creep becomes significant
- At 800°C, oxidation becomes severe
- Would fail in minutes in a turbine environment

The gulf in performance is achieved through one of the most elegant and complex microstructures in all of engineering: the **gamma/gamma-prime (γ/γ′) two-phase system** in nickel.

---

## 2. Why Normal Alloys Fail at High Temperature

Three simultaneous failure modes destroy ordinary alloys in high-temperature service:

### 2.1 Thermal Softening (Loss of Yield Strength and Creep Resistance)

All metals soften as temperature rises. For most structural alloys:
- T_H = 0.5 → yield strength drops to ~25–50% of room-temperature value
- Creep becomes significant → sustained loads cause slow deformation
- Dislocation climb is activated → cannot sustain high stress

A standard 316 stainless steel has σ_y ≈ 450 MPa at room temperature and only ~70 MPa at 800°C. This 6× drop means parts designed for room temperature would yield immediately in a turbine.

### 2.2 Phase Instability

As temperature rises, phases transform:
- Carbon steels: pearlite → austenite at 727°C; then rapid oxidation
- Aluminum alloys: precipitates (GP zones, θ′) dissolve above ~180–200°C; T_H already 0.5 at 350°C
- Even stainless steels experience sigma phase embrittlement at 600–900°C

The microstructure designed at room temperature doesn't survive service conditions.

### 2.3 High-Temperature Oxidation and Hot Corrosion

At 800°C+, most metals oxidize aggressively:
- Carbon steels: form loose, non-protective FeO (scales spall off rapidly)
- Even Cr-bearing steels need > 13% Cr for minimal protection
- Nickel: forms NiO above 600°C; inadequate protection alone

Combustion gases also contain SO₂, NaCl, and vanadium compounds → "hot corrosion" — a catastrophic attack that destroys oxide scales and penetrates metal at rates 10–100× faster than pure oxidation.

### The Combined Attack

In a turbine engine, all three attack modes operate simultaneously under:
- 150–250 MPa centrifugal stress
- Temperature cycling (every takeoff-landing cycle is a thermal cycle)
- Oxidizing and hot corrosion atmosphere
- Vibration loads (high-cycle fatigue)

To survive this environment for 30,000 hours, materials must be specifically engineered for every one of these threats simultaneously. That is the challenge superalloys meet.

---

## 3. The Three Families of Superalloys

### 3.1 Nickel-Base Superalloys

The dominant family for the hottest applications.

**Why nickel?**
- FCC crystal structure (stable, no allotropic transformations up to melting)
- Nickel can dissolve a wide variety of alloying elements (Cr, Co, Al, Ti, Mo, W, Re, Ru, Ta, Nb, Hf)
- The unique γ/γ′ two-phase strengthening system is possible only with Ni (see Chapter 43-44)
- Good oxidation resistance (forms NiO) enhanced dramatically by Cr and Al additions
- Ferromagnetic (weakly) → no special thermal expansion issues like Fe alloys

**Temperature capability:** 950–1050°C (blade metal temperature, continuous); up to 1100°C with TBC + cooling.

**Examples:** CMSX-4, CMSX-10, René N5, René N6, PWA 1484, TMS-238, SRR99, RR3000.

### 3.2 Cobalt-Base Superalloys

Cobalt has a higher melting point than nickel (1495°C vs 1455°C). Why don't cobalt alloys dominate?

- Co-base alloys lack the γ′ strengthening mechanism (no ordered Co₃Al equivalent at useful scale)
- They are primarily solid-solution and carbide strengthened → lower specific strength
- Much better **hot corrosion resistance** (due to higher Cr content possible)
- Good weldability → used for stationary parts (vanes, combustors)
- High thermal conductivity → good for heat transfer applications

**Examples:** Haynes 188, Haynes 25 (L-605), Mar-M-509, FSX-414.
**Applications:** High-pressure turbine vanes (stationary), combustor liners, industrial gas turbine parts.

### 3.3 Iron-Base Superalloys

Iron-base superalloys are cheaper than Ni or Co alloys and used at lower temperatures:

- Based on austenitic stainless steels with Ni added to stabilize the FCC structure
- Strengthened by γ″ (Ni₃Nb in IN718), carbides, and solid solution
- Limited to ~ 700°C service
- Much cheaper and more weldable than Ni-base

**Examples:** Incoloy 901, A-286, Incoloy 800H, **Inconel 718** (IN718 — the most widely used superalloy by weight).

**Inconel 718 (IN718):**
- Fe-Ni-Cr base with ~53% Ni, 19% Cr, ~18% Fe, 5.1% Nb, 3% Mo, 0.9% Ti, 0.5% Al
- Strengthened by γ″ (Ni₃Nb, disc-shaped, forms at 620°C)
- Excellent weldability: unlike most Ni superalloys, IN718 can be TIG/electron-beam welded
- Used for turbine **disks**, compressor components, structural parts — NOT rotating turbine blades (too hot)
- The dominant superalloy by volume — used in almost every jet engine

---

## 4. A Brief History — From Nimonic to TMS-238

The history of superalloys is a story of continuously pushing temperature capability through alloy design, microstructure control, and processing innovation:

**1930s–1940s:** British work at Mond Nickel Company and later NGTE develops Nimonic 80 (80% Ni, 20% Cr + Ti, Al additions). First nickel superalloys for the Whittle jet engine. Temperature ~700°C.

**1940s–1950s:** Americans develop M-252, Waspaloy, René 41. Increased Al and Ti content → more γ′ → better creep. Casting techniques improve. 800°C capability reached.

**1950s–1960s:** Investment casting (lost-wax) enables complex blade geometry. Alloys like MAR-M200, IN100, B1900 push to 850°C. Directional solidification patented (B. Piearcey, Pratt & Whitney, 1966).

**1966:** First directionally solidified (DS) turbine blades in service (Pratt & Whitney JT9D). Columnar grains eliminate transverse grain boundaries. 30°C capability jump.

**1968:** First single-crystal (SX) blade grown in laboratory. PWA 1480 enters service in 1982 (F100 engine). No grain boundaries at all → major creep improvement.

**1983–1990s:** Second-generation SX alloys: 3% Re addition (PWA 1484, CMSX-4, René N5). Another ~30°C jump. Rhenium becomes a critical strategic material.

**1990s–2000s:** Third-generation SX alloys: 5–6% Re (CMSX-10, René N6, TMS-75). Better creep but density penalty and stability challenges with TCP phases.

**2000s–2010s:** Fourth generation: Ru added to stabilize against TCP phases while maintaining high Re content (TMS-138, EPM-102).

**2010s–present:** Fifth and sixth generation (TMS-238): further optimized Re/Ru/W/Mo/Ta balance; density-normalized creep strength continues to improve.

---

## 5. Temperature Capability — The Evolution

```
Turbine Inlet Temperature (°C) vs. Year of Introduction
1000 ─                          ●─── 6th gen SX (TMS-238)
 990 ─                       ●── 5th gen SX
 980 ─                    ●── 4th gen SX 
 970 ─                 ●── 3rd gen SX (CMSX-10)
 960 ─              ●── 2nd gen SX (CMSX-4, Re addition)
 950 ─           ●── 1st gen SX (PWA 1480)
 920 ─         ●── DS (directional solidification)
 890 ─      ●── Polycrystalline investment cast (IN100, MAR-M200)
 850 ─    ●── Wrought Waspaloy
 750 ─  ●── Nimonic 80 (1940s)
     1940  1950  1960  1970  1980  1990  2000  2010  2020
```

The x-axis represents actual metal temperature at the blade, accounting for cooling. Turbine inlet gas temperatures are ~300°C higher than the metal temperature.

Each major jump corresponds to:
- **1966**: Directional solidification (+30°C)
- **1982**: Single crystal (+30°C from DS)
- **1983**: Rhenium addition (+25°C from Gen 1 SX)
- **1990s**: 5-6% Re (+15°C from Gen 2)
- **2000s**: Ru addition (+15°C from Gen 3 stability issues)

Total improvement: ~250°C in alloy capability over 80 years — equivalent to half a Carnot cycle improvement in engine efficiency.

---

## 6. Where Are Superalloys Used?

### Jet Engines

The jet engine is the primary market for superalloys. In a typical large turbofan (GE90, Rolls-Royce Trent):

```
High-pressure turbine (HPT):
  ├── Blades (rotating): Gen 2–4 SX Ni superalloy, ~50 kg/engine
  ├── Vanes (stationary): DS or PC Ni superalloy, Co alloys
  └── Disk: PM IN718 or René 88DT, ~100 kg/disk

Combustor:
  └── Hastelloy X, Haynes 230, CM247LC sheet/castings

Low-pressure turbine (LPT):
  └── DS Mar-M247, IN100, IN792 blades; IN718 disk

Compressor:
  └── Ti alloys (low-T stages), IN718, A-286 (high-T stages)
```

A modern turbofan contains ~3,000–5,000 kg of superalloy components.

### Industrial Gas Turbines

Power generation turbines (GE Frame 7/9, Siemens SGT, Mitsubishi JAC) operate at similar or slightly lower temperatures than aircraft engines, but for much longer lives (>100,000 hours vs. ~30,000 for aircraft). They use large DS or coated PC Ni superalloy blades and IN718 disks.

### Rocket Engines

Combustion temperatures can exceed 3000°C — far above any metallic alloy. Turbopumps operate at cryogenic temperatures. Superalloys are used in:
- Turbopump rotor disks (IN718)
- Hot-gas manifolds (Inconel 625)
- Thrust chambers (nickel or Cu alloy with regenerative cooling)

### Other Applications

- Nuclear reactor components (IN600, IN625 for heat exchangers)
- Chemical processing (Hastelloy C-276 for aggressive acids)
- Oil and gas (IN718, IN625 for downhole tools)
- Medical implants (Co-Cr alloys for hip replacements)

---

## 7. The Key Properties Required

For a turbine blade material, the property requirements are extraordinarily demanding:

| Property | Requirement | Why |
|----------|------------|-----|
| Creep strength | >700 MPa at 980°C/1000h | Centrifugal loads at temperature |
| Fatigue strength | >300 MPa LCF at 950°C | Thermal cycling |
| Oxidation resistance | <10 mg/cm² after 500h at 1100°C | Combustion environment |
| Hot corrosion resistance | Survive Type I and II | Contaminated gas streams |
| Phase stability | No TCP phases after 1000h | Microstructure must not degrade |
| Density | <9 g/cm³ | Centrifugal stress = f(density × ω²r) |
| Castability | Shrink-free SX casting | Dimensional precision |
| Thermal conductivity | ~10–12 W/m·K | Heat removal from blade core |
| TBC compatibility | Good bond coat adhesion | Thermal protection system |

No single element addresses all these requirements — superalloy design is a multi-dimensional optimization across 10–15 alloying elements simultaneously. This is why CALPHAD and computational alloy design have become essential (Chapter 64).

---

## 8. Cost and Strategic Importance

Superalloys are expensive:

| Material | $/kg (approx.) |
|----------|---------------|
| Carbon steel | $0.50 |
| 316 Stainless steel | $5–8 |
| IN718 | $40–80 |
| CMSX-4 single crystal | $200–500/kg raw; $5,000–20,000 per finished blade |
| Re metal (key addition) | $2,500–4,000/kg |

A single high-pressure turbine blade can cost $10,000–50,000 fully finished. A commercial jet engine has 60–80 HPT blades per stage. The turbine section of one GE90 engine contains ~$8–15 million in superalloy components.

**Strategic importance of rhenium:**
Rhenium is produced primarily as a byproduct of molybdenum smelting from porphyry copper ores. World production is ~60 tonnes/year. 80% of rhenium goes into superalloys. The primary producers are Chile, Kazakhstan, USA, and Poland. Supply disruption (geopolitical risk) would halt jet engine manufacturing.

This is a genuine national security concern — the US maintains strategic reserves of rhenium, and significant research continues on reduced-Re alloys and Re-free compositions (using Ru for stability instead).

---

## Summary

- **Superalloys** = high-temperature alloys for T > 540°C service; must maintain strength, stability, and oxidation resistance simultaneously.
- **Three families**: Ni-base (hottest, γ/γ′ strengthening), Co-base (better hot corrosion, lower strength), Fe-base (cheaper, weldable, lower temperature, e.g., IN718).
- Normal alloys fail at high temperature by: thermal softening, phase instability, and oxidation/hot corrosion.
- **Temperature capability** has risen ~250°C since 1940 through alloy design (Ni-base), processing (DS, SX), and coatings.
- **IN718** is the most widely used superalloy by weight (disks, structure); CMSX-4 and its successors are the benchmark for blade (rotating HPT) applications.
- Key markets: jet engines (primary), industrial gas turbines, rockets, chemical processing.
- Cost: finished SX turbine blades cost $5,000–50,000 each. Rhenium at 3–6 wt% is a strategic critical material.

**Next chapter:** We go inside the Ni superalloy microstructure — the γ matrix, the γ′ precipitate, the carbides, and the harmful TCP phases — and understand how each element and each phase contributes to the extraordinary properties of these alloys.

---

## Exercises

1. At what homologous temperature (T/T_melt) does a Ni superalloy operate at 1000°C? At what T_H does steel at 600°C operate? Using the creep onset criterion (T_H > 0.4), which is more severely creep-limited?

2. Compare IN718 (Fe-Ni base) and CMSX-4 (Ni-base SX) for turbine blade use. List three reasons CMSX-4 is preferred for HPT blades and two reasons IN718 might be preferred for disk applications.

3. The temperature capability of turbine materials has improved ~250°C since 1940. Engine thermal efficiency (Brayton cycle) improves with turbine inlet temperature. If efficiency improves by ~1% per 25°C increase in TIT, roughly how much has materials improvement contributed to modern engine efficiency?

4. Rhenium world production is ~60 tonnes/year. If CMSX-4 is ~3 wt% Re and a large turbofan engine contains 50 kg of HPT blades, how many large turbofan engine sets (2 engines per aircraft) could be supplied per year from global Re production (assuming Re is entirely used for superalloys)? What does this tell you about supply chain risk?

5. Why is nickel (FCC, stable to melting point) a better superalloy base than iron (BCC→FCC transformation at 912°C)? What practical problem would an iron-base alloy face that a nickel-base avoids?
