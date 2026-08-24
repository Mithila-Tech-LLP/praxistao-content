# Chapter 17: Heat Treatment Fundamentals — Using Temperature to Control Microstructure

> **"Iron is the most versatile structural metal on earth, largely because of the iron-carbon system's remarkable phase transformations. The same lump of iron-carbon alloy, depending only on how you heat and cool it, can be soft enough to machine or hard enough to scratch glass, tough enough to absorb impacts or brittle enough to shatter like ceramic. Heat treatment is the metallurgist's art of converting potential into performance."**

---

## Table of Contents

1. [What Is Heat Treatment?](#1-what-is-heat-treatment)
2. [Iron's Allotropic Transformations](#2-irons-allotropic-transformations)
3. [Austenite — The Starting Phase for Most Heat Treatments](#3-austenite--the-starting-phase-for-most-heat-treatments)
4. [TTT Diagrams — Time-Temperature-Transformation](#4-ttt-diagrams--time-temperature-transformation)
5. [CCT Diagrams — Continuous Cooling Transformation](#5-cct-diagrams--continuous-cooling-transformation)
6. [Critical Cooling Rates and Hardenability](#6-critical-cooling-rates-and-hardenability)
7. [Furnaces, Atmospheres, and Temperature Uniformity](#7-furnaces-atmospheres-and-temperature-uniformity)
8. [Heat Treatment Process Overview](#8-heat-treatment-process-overview)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. What Is Heat Treatment?

Heat treatment is a controlled sequence of **heating, holding, and cooling** operations designed to achieve a specific microstructure — and therefore specific mechanical properties — in a metal.

**Why it works:** Metals can exist in different crystal structures (allotropes) and can dissolve different amounts of alloying elements at different temperatures. By controlling temperature and time, we control:
- Which crystal structure exists (austenite, ferrite, martensite)
- Whether alloy elements are dissolved or precipitated
- Grain size
- Internal (residual) stresses

**Most important system:** Iron-carbon (steels). Most steel properties are achieved through heat treatment, not just composition.

---

## 2. Iron's Allotropic Transformations

Pure iron exists in three crystal structures as temperature changes:

```
Temperature (°C)
     │
1538 │──────────── Liquid iron
     │
1394 │──────────── δ-iron (BCC) ← high-temperature BCC phase
     │
912  │──────────── γ-iron (FCC, "austenite") ← paramagnetic
     │
770  │──────────── Curie temperature (magnetic transition, no crystal change)
     │
 RT  │──────────── α-iron (BCC, "ferrite") ← ferromagnetic
```

**The critical transformation:**
- γ-iron (FCC) → α-iron (BCC) on cooling through 912°C (A₃ temperature in alloys)
- FCC has much larger octahedral sites → dissolves up to 2.14% C at 1,147°C
- BCC has very small octahedral sites → dissolves only 0.022% C at 727°C (A₁ temperature)
- This huge difference in C solubility is the engine of steel heat treatment

---

## 3. Austenite — The Starting Phase for Most Heat Treatments

**Austenitizing:** Heating steel above the A₁ (727°C for eutectoid) or A₃ (composition-dependent for hypoeutectoid) temperature → all carbon dissolves in FCC austenite:

```
Carbon solubility:
  α-ferrite: max 0.022% C at 727°C
  γ-austenite: max 2.14% C at 1,147°C
  
→ On austenitizing, ALL cementite (Fe₃C) dissolves into austenite
→ Uniform FCC solid solution
→ Subsequent cooling determines what phases form
```

**Typical austenitizing temperatures:**
- Plain carbon steels: 800–850°C (hypo), 820–840°C (hyper)
- Alloy steels: 830–950°C (higher to dissolve alloy carbides)
- Stainless steels (304): 1,050°C (dissolve all Cr₂₃C₆)
- Tool steels (H13): 1,000–1,020°C
- Nickel superalloys: 1,100–1,350°C (solution treat γ′)

**Grain growth during austenitizing:**
Higher temperature → faster grain growth (Ch 09). ASTM grain size G = 5–10 is target for most structural steels. Alloyed with Nb, V → microalloyed grain boundary pinning → fine austenite grain → fine final product.

---

## 4. TTT Diagrams — Time-Temperature-Transformation

When austenite is **rapidly quenched to a fixed temperature** and held isothermally, it transforms to equilibrium phases. The TTT diagram maps: how long it takes to start and finish the transformation at each temperature.

```
TTT diagram for eutectoid steel (0.8% C):

Temperature (°C)

  727 ─────────────────────────────────────────── A₁ (eutectoid)
      │  Austenite (stable)
      │
  600 │         ╮
      │        /│╰─ Pearlite start (nose = ~550°C, 1 sec)
  550 │    ╭──╯  ╰── Pearlite finish
      │   /         
  400 │  /  ← Bainite start (lower nose)
      │ │
  250 │ │   Bainite finish
      │ │
  200 │─┴──────────────────── M_s (martensite start)
      │
  100 │────────────────────── M_50 (50% martensite)
      │
  RT  │────────────────────── M_90
      
      ─┴────────────────────────────────── log(time) →
        0.1s  1s  10s  100s  1000s  10000s
```

**Reading the TTT diagram:**
- **Pearlite field (upper nose):** Austenite → ferrite + cementite lamellae. Coarser at higher T (close to 727°C), finer at lower T.
- **Bainite field (lower nose):** Austenite → bainite (complex ferrite + fine carbide; either upper bainite or lower bainite depending on T).
- **Martensite region:** Below M_s (~220–230°C for 0.8% C steel), austenite transforms diffusionlessly → martensite. NOT shown as a curve — martensite forms immediately when austenite is quenched below M_s. Temperature, not time, controls fraction transformed.

**Martensite start temperature M_s:**
```
M_s (°C) = 550 - 350×[C] - 40×[Mn] - 20×[Cr] - 17×[Ni] - 10×[Mo]
            (Andrews equation, wt%)
```

---

## 5. CCT Diagrams — Continuous Cooling Transformation

In practice, steel is continuously cooled (not isothermally held). CCT diagrams show transformation during continuous cooling:

```
CCT diagram (overlay multiple cooling curves on TTT-like diagram):

Temperature
  │ 727°C
  │─────────────────────────────────── A₁
  │     A    |Pearlite start
  │    u     |          |Pearlite finish
  │   s      |          
  │  t ──────────────── Bainite start
  │ e                   ──── Bainite finish
  │n
  │ ─────────────────────────────── M_s
  │
  └────────────────────────── log(time) →
      ↑ slow     ↑ fast       ↑ very fast
      (pearlite) (bainite)    (martensite)
      
     1°C/min    1°C/s        1,000°C/s
     FURNACE    AIR          WATER
     COOL       COOL         QUENCH
```

**Cooling curves determine microstructure:**
- Slow cool (furnace): fully pearlitic (equilibrium)
- Moderate cool (air): mixed pearlite + bainite
- Fast cool (oil quench): mostly martensite
- Very fast cool (water quench): fully martensite
- Cooling fast enough to miss the "nose" of both C curves → martensite

---

## 6. Critical Cooling Rates and Hardenability

**Critical cooling rate (CCR):** The minimum cooling rate to form 100% martensite in a steel. Depends strongly on alloy composition.

**Hardenability:** The ability of a steel to form martensite throughout its cross-section (not just at the surface). High hardenability = martensite forms even in thick sections at low cooling rates.

**Why hardenability matters:**
- A thick shaft quenched in water: surface cools fast → martensite; center cools slow → pearlite/bainite
- Only the martensite achieves the desired high hardness after tempering
- Alloying elements (Mn, Cr, Ni, Mo) SHIFT TTT curves to longer times → easier to form martensite

```
Jominy End Quench Test (ASTM A255):

           Water spray
              ↓
    ─────────────────────────────
    │ Standard bar, 100×25mm dia │
    ─────────────────────────────
    │← measure hardness every 2mm→│
    
Hardness vs. distance from quench end = hardenability curve
```

**Ideal critical diameter (D_I):** The diameter of a round bar that would form 50% martensite at the center when quenched in an "ideal" quench medium. D_I = 25–250 mm for various steels (pure iron ~0 mm, highly alloyed tool steel > 500 mm).

---

## 7. Furnaces, Atmospheres, and Temperature Uniformity

### Furnace Types

| Furnace | Temperature Range | Atmosphere | Use |
|---------|-----------------|------------|-----|
| Air furnace (box) | RT–1,200°C | Oxidizing (air) | Annealing, normalizing |
| Vacuum furnace | RT–1,400°C | Vacuum/inert gas | SX solution treat, Ti, precision alloys |
| Salt bath | 150–1,200°C | Molten salt | Marquenching, austempering, tool steel |
| Fluidized bed | 300–1,000°C | N₂, air | Fast, uniform heating |
| Pusher/roller | RT–1,200°C | Controlled | Continuous production (strip, wire) |
| Induction | Localized | None | Case hardening, specific zones |

### Atmosphere Control

Oxidizing atmosphere: surface scale forms → acceptable for some operations, not for precision parts.

Controlled/protective atmosphere (reducing or neutral):
- **Endothermic gas:** CO + H₂ + N₂ (cracked propane) → "endogas" — neutral to slightly carburizing
- **Nitrogen-hydrogen:** Bright annealing of stainless steel
- **Vacuum:** Best surface quality; essential for superalloys (prevents oxidation)

**Carbon potential control:** For carburizing, the atmosphere carbon potential must be set to the target surface carbon content. Measured by oxygen probe (dew point or O₂ sensor):
```
log(p_O₂) = A + B/T  (equilibrium between atmosphere and steel surface)
```

### Temperature Uniformity

±3–5°C typical for precision heat treatment (aerospace). ±10°C acceptable for structural steel.
Thermocouples at multiple locations + controller.
AMS 2750 pyrometry specification (aerospace): calibrated sensors, type B or R, quarterly verification.

---

## 8. Heat Treatment Process Overview

All the specific heat treatments in Chapters 18–21 use the same fundamental steps:

```
General heat treatment sequence:

1. HEAT: bring material from RT to target temperature
         Rate: usually slow enough to avoid thermal shock
         
2. SOAK: hold at temperature for sufficient time
         To: dissolve all phases, homogenize composition, 
             allow diffusion, achieve target grain size
         Time: minutes to hours depending on section size
         
3. COOL: controlled cooling to room temperature
         Rate: determines which transformation products form
         Media: furnace, air, oil, water, salt, polymer quench
         
4. OPTIONAL ADDITIONAL STEPS:
   Tempering, aging, stress relief, etc.
```

**The four fundamental steel heat treatments:**
| Treatment | Heat to | Soak | Cool | Result |
|-----------|---------|------|------|--------|
| Annealing | A₃+30°C | Long | Furnace cool | Soft, maximum ductility |
| Normalizing | A₃+50°C | 1 h/25mm | Air cool | Fine pearlite, improved toughness vs. annealing |
| Hardening | A₃+30°C | Short | Oil/water quench | Martensite (hard, brittle) |
| Tempering | 150–700°C | 1–2 h | Air cool | Martensite → tempered martensite (hard+tough) |

These are covered in detail in Ch 18, 19, and 20.

---

## Summary

| Concept | Key Point |
|---------|-----------|
| Iron allotropes | α-BCC (RT), γ-FCC (912–1394°C), δ-BCC (1394–1538°C) |
| Austenite | FCC; dissolves up to 2.14% C; starting point for most heat treatments |
| TTT diagram | Maps transformation time/temperature for isothermal hold; pearlite, bainite, martensite noses |
| CCT diagram | Maps transformation for continuous cooling; faster cool → more martensite |
| Hardenability | Ability to form martensite in thick sections; Mn/Cr/Mo/Ni increase it |
| Atmosphere | Vacuum for precision alloys; controlled gas for carburizing; inert for titanium |

---

## Exercises

1. Using the Andrews equation M_s = 550 - 350[C] - 40[Mn] - 20[Cr] - 17[Ni] - 10[Mo], calculate M_s for: (a) 0.4% C plain carbon steel, (b) 4340 steel (0.4%C-0.8%Mn-1.8%Ni-0.8%Cr-0.25%Mo), (c) 1% C high-carbon steel for bearing rings. How does the alloy steel (4340) compare to the plain carbon steel in terms of M_f (assume M_f ≈ M_s - 150°C)? What fraction of austenite might be retained at RT in the high-carbon steel if M_f < RT?

2. A steel TTT diagram shows the pearlite nose at 550°C, 1.5 seconds. If you must form 100% martensite, you need to cool from 750°C to 250°C faster than 1.5 seconds (to miss the pearlite nose). Calculate the required average cooling rate in °C/s. For a 25mm diameter bar with quench coefficient h = 1,000 W/m²K (water quench), the cooling rate at the center is roughly 30°C/s. Is full martensite achievable at the center?

3. Compare the CCT diagrams for: (a) 1020 plain carbon steel (0.2%C), (b) 4340 alloy steel (0.4%C, high alloy). Which has higher hardenability (Jominy D_I)? For a 50mm diameter shaft, what microstructure would you expect at the center after water quenching each steel? What implications does this have for choosing steel for a large shaft requiring high core hardness?

4. A vacuum furnace is used to solution treat single-crystal CMSX-4 blades at 1,310°C for 4 hours. (a) Why is vacuum atmosphere necessary (hint: Al and Ti oxidize rapidly at this temperature). (b) Temperature uniformity requirement is ±5°C for this process. If one thermocouple reads 1,315°C, what happens to the blade? (Hint: the γ+γ′ → γ solvus temperature is ~1,325°C; exceeding it causes incipient melting). (c) What NDE method would you use to detect incipient melting damage?

5. Sketch a TTT diagram for a hypothetical high-alloy tool steel and show how alloying elements have shifted the pearlite nose compared to plain carbon steel. Label: (a) A₁ temperature, (b) pearlite nose, (c) bainite nose, (d) M_s line. Mark the path on the diagram for: (a) furnace cooling (annealing), (b) air cooling (normalizing), (c) oil quenching (hardening). For each path, state the expected microstructure and typical hardness (HRC) range.
