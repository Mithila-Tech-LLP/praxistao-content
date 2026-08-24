# Chapter 23: Cast Irons — Iron's Other Face

> **"For every turbine blade made from a nickel superalloy with painstaking care, there are a million cast iron cylinder blocks made in automated foundries at $1/kg. Cast iron is not a lesser material — it is a different tool. Engine blocks, water pipes, cookware, brake rotors, machine bases: cast iron's combination of castability, vibration damping, and wear resistance makes it irreplaceable for the applications it serves."**

---

## Table of Contents

1. [Cast Iron vs. Steel — The Fundamental Difference](#1-cast-iron-vs-steel--the-fundamental-difference)
2. [The Carbon Equivalent and Castability](#2-the-carbon-equivalent-and-castability)
3. [Grey Cast Iron — The Workhorse](#3-grey-cast-iron--the-workhorse)
4. [White Cast Iron — Hard and Brittle](#4-white-cast-iron)
5. [Ductile (Nodular) Cast Iron — Grey's Tougher Cousin](#5-ductile-nodular-cast-iron)
6. [Malleable Cast Iron](#6-malleable-cast-iron)
7. [Compacted Graphite Iron (CGI)](#7-compacted-graphite-iron-cgi)
8. [Heat Treatment of Cast Irons](#8-heat-treatment-of-cast-irons)
9. [Casting and Manufacturing](#9-casting-and-manufacturing)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Cast Iron vs. Steel — The Fundamental Difference

Cast iron is an iron-carbon alloy with **2.14–6.7% carbon** — far above the 2.14% limit for steel. The excess carbon exists as:
- **Graphite** (soft, flaky or nodular) → grey/ductile/malleable cast iron
- **Iron carbide Fe₃C** (hard, brittle) → white cast iron

**Why not use steel for everything?**
Cast iron advantages:
- Excellent castability: low melting point (~1,200°C vs 1,500°C for steel), good fluidity → fills complex molds
- Lower cost: no alloying additions, scrap-based melting
- Excellent machinability (especially grey): graphite is self-lubricating
- Good vibration damping (grey): graphite flakes absorb vibration energy → ideal for machine tool bases
- Wear resistance: white cast iron, or graphite lubrication in grey

Steel advantages over cast iron:
- Higher ductility and toughness
- Weldable
- Higher tensile strength (at same hardness)
- Better fatigue resistance

---

## 2. The Carbon Equivalent and Castability

The **carbon equivalent (CE)** accounts for both C and other elements (Si, P promote graphite formation):
```
CE = %C + (%Si + %P) / 3

CE < 4.26: hypoeutectic cast iron (most common)
CE = 4.26: eutectic composition (best castability)
CE > 4.26: hypereutectic cast iron
```

**Silicon is critical:** Si promotes graphite (vs. cementite) formation. Without Si, most cast irons would be white (all cementite). Typical grey iron: 1–3% Si.

**Solidification products:**
- With Si: stable equilibrium → graphite (soft, lubricating)
- Without Si: metastable (fast cooling): → Fe₃C (cementite, hard)

```
Fe-C metastable (iron-cementite) vs. stable (iron-graphite):

Metastable: γ → α + Fe₃C      ← forms at fast cooling or low Si
Stable:     γ → α + C (graphite) ← forms at slow cooling or high Si

Si and slow cooling promote the STABLE (graphite) reaction
```

---

## 3. Grey Cast Iron — The Workhorse

**Composition:** 2.5–4.0%C, 1–3%Si, 0.5–1.0%Mn, < 0.3%S, < 0.1%P

**Microstructure:** Graphite flakes in a matrix of ferrite + pearlite (or all pearlite, or all ferrite depending on Si content and cooling rate):

```
Grey cast iron microstructure (schematic):

    ─────────────────────────────────
   │  /  \  /  /  \  /  \  /  \  / │
   │ /    \/  /    \/    \/    \/  │ ← Type A graphite flakes
   │  \   /  \   /  \   /  \  /   │   (random orientation)
   │   \ / Pearlite/Ferrite matrix │
    ─────────────────────────────────
```

**Graphite flake types (ASTM A247):**
- Type A: uniform distribution, random orientation → best mechanical + damping properties
- Type B: rosette clusters → weaker
- Type C: kish graphite (hypereutectic) → very coarse
- Type D: interdendrite undercooled → medium properties
- Type E: preferred orientation → direction-dependent properties

**Properties:**
- Compressive strength: 570–1,300 MPa (MUCH higher than tensile → strong in compression)
- Tensile strength: 100–400 MPa (weak in tension — flakes are stress concentrators)
- Elongation: near 0% (essentially brittle in tension)
- Hardness: 140–340 HB
- Excellent damping: 100× better than steel at vibration absorption

**Grades (ASTM A48):**
- Class 20: σ_UTS = 138 MPa (low grade, decorative, pipe)
- Class 40: σ_UTS = 276 MPa (standard engineering, cylinder blocks, cylinder heads)
- Class 60: σ_UTS = 414 MPa (pearlitic matrix, machine tools)

**Applications:** Engine blocks, cylinder heads, brake rotors, machine tool beds, cookware, manhole covers, counterweights.

**Chilling:** The surface of a grey iron casting that contacts the cold mold can solidify as white iron ("chill zone") because it cools too fast for graphite to form. This surface chill zone is very hard and wear-resistant → exploited for cylinder bore liners.

---

## 4. White Cast Iron

**What it is:** All carbon exists as Fe₃C (cementite) — NO graphite. Very fast solidification or low Si forces metastable Fe₃C formation.

**Properties:**
- Very hard (700–1,000 HV, 65+ HRC)
- Very brittle (no ductility)
- Wear-resistant

**Composition:** Low Si (<1%), or rapid cooling of any grey iron composition.

**Applications:**
- Ball mill liners, crusher jaws, dredge pump liners → extreme abrasion
- Intermediate form for conversion to malleable cast iron (§6)
- Chilled gray iron castings (surface white, core grey)

**High-chromium white iron:** 12–28%Cr + C → Cr₇C₃ carbides instead of Fe₃C → even harder, better corrosion resistance than regular white iron. Used for cement and mining equipment.

---

## 5. Ductile (Nodular) Cast Iron

**The revolution (1948):** Adding small amounts of Mg (0.03–0.06%) or Ce to molten grey iron changes graphite from flakes to spheres (nodules) → transforms a brittle material into a tough, ductile one.

**Mechanism:** Mg (or Ce) poisons certain crystal faces of graphite → can't grow along hexagonal basal planes → instead grows as spheres to minimize surface energy.

```
Comparison of graphite morphologies:

Grey iron (flakes):         Ductile iron (nodules):
   ─────────────────────       ─────────────────────
  │ /  \  /  /  \  /    │     │  ●    ●    ●    ●  │
  │/    \/  /    \/      │     │     ●    ●    ●    │
  │  ─────   ─────       │     │  ●    ●    ●    ●  │
   ─────────────────────       ─────────────────────
   
   Flakes = STRESS            Nodules = much less
   CONCENTRATORS               stress concentration
   → brittle                   → ductile!
```

**Properties of ductile cast iron (ASTM A536):**
| Grade | σ_y (MPa) | σ_UTS (MPa) | Elongation (%) | Microstructure |
|-------|-----------|-------------|----------------|----------------|
| 60-40-18 | 275 | 415 | 18 | Ferrite |
| 65-45-12 | 310 | 450 | 12 | F+P |
| 80-55-06 | 380 | 552 | 6 | Pearlite |
| 100-70-03 | 483 | 690 | 3 | Martensite (Q&T) |
| 120-90-02 | 621 | 827 | 2 | Bainite (austempered) |

**Austempered ductile iron (ADI):** Heat treatment → bainitic matrix → 
σ_UTS up to 1,600 MPa with 1–10% elongation — competes with alloy steel!

**Applications:** Automotive crankshafts (replacing forged steel in many cases), gears, hydraulic cylinders, pipe fittings, agricultural equipment.

---

## 6. Malleable Cast Iron

**Process:** Cast as white cast iron → anneal at 900°C for hours → convert Fe₃C to temper carbon (graphite nodules, but less round than ductile iron):

```
White cast iron → anneal → Malleable cast iron

Stage 1 (900-970°C): Fe₃C → graphite nuclei (slow)
Stage 2 (700-730°C): remaining Fe₃C → ferrite + graphite (faster)

Total anneal: 48–70 hours (very slow → expensive)
```

**Why make it when ductile iron exists?** Malleable iron is OLDER (1722 vs. 1948) and can be made in thinner sections (white iron casts well in thin sections; spheroidizing with Mg can fail in thin sections). Also, some properties differ slightly.

**Grades:** Similar to ductile iron in strength/ductility, but lower Si content.

**Applications:** Pipe fittings, electrical hardware, farm equipment (historical), now mostly replaced by ductile iron.

---

## 7. Compacted Graphite Iron (CGI)

**A hybrid:** Graphite morphology between flakes and nodules → "wormlike" or "vermicular" graphite. Achieved with precise Mg content control.

**Properties:** Between grey and ductile iron:
- 2× tensile strength of grey iron
- Better thermal conductivity than ductile iron (graphite network connected)
- Better fatigue resistance than grey iron
- Still excellent castability

**The killer application: diesel engine cylinder blocks.**
Modern diesel engines need high combustion pressure (200+ bar vs. 150 bar for gasoline) → need a stronger, thicker cylinder block → cast iron too heavy if grey, thermal conductivity too low if ductile. CGI is ideal.

**Examples:** Ford 6.7L Power Stroke diesel, many European diesel engines now use CGI blocks vs. traditional grey iron.

---

## 8. Heat Treatment of Cast Irons

**Annealing grey iron:** Reduce hardness for machining (700°C, furnace cool). Converts some pearlite → ferrite.

**Stress relief:** 500–600°C, slow cool — relieves casting stresses without changing structure.

**Hardening grey iron:** Austenitize (850°C) + quench → martensite in matrix but flakes unchanged → harder but brittle (rarely done).

**Austempering of ductile iron (ADI):**
```
Austenitize at 900°C → hold
    ↓
Quench into salt bath at 230–400°C
    ↓
Hold 1–4 hours → ausferrite (bainitic) matrix forms
    ↓
Air cool → no further transformation

Result: σ_UTS = 900–1,600 MPa, elongation = 1–10%
```

---

## 9. Casting and Manufacturing

**Foundry process:**
- Melt in cupola furnace (coke-fueled) or induction furnace
- Pour at 1,200–1,400°C (lower T than steel → easier equipment)
- Sand casting most common (sand mold around pattern → pour → break out → machine)
- Die casting: for small parts (white iron or Al, not generally grey iron)

**Solidification challenges:**
- Grey iron expands slightly on solidification (graphite lower density than liquid → expansion as graphite forms → fills mold well)
- This is why grey iron castings are so good — near-net-shape, minimal shrinkage
- Ductile iron: Mg reduces surface tension of graphite → must account for more shrinkage

---

## Summary

| Type | Carbon Form | Tensile Strength | Ductility | Key Property | Applications |
|------|------------|-----------------|-----------|-------------|--------------|
| Grey | Flakes | 100–400 MPa | ~0% | Damping, cheap, castable | Engine blocks, brakes |
| White | Fe₃C | Very high compressive | Brittle | Hardness, wear | Crusher liners |
| Ductile | Nodules | 415–1,600 MPa | 1–18% | Strength + toughness | Crankshafts, gears |
| Malleable | Temper C | 350–700 MPa | 5–12% | Thin sections | Fittings |
| CGI | Wormlike | 300–600 MPa | 2–5% | Strength + conductivity | Diesel engine blocks |

---

## Exercises

1. A grey cast iron has CE = 4.1 (3.4%C + 2.1%Si). Is it hypo-, hyper-, or eutectic? At slow cooling, will the solidified structure contain mostly graphite or cementite? If the cooling rate is doubled (thinner section), predict how the graphite morphology changes and how hardness changes.

2. An automotive crankshaft is currently grey cast iron (Grade 40, σ_UTS = 276 MPa). Due to higher engine power demands (50% more peak cylinder pressure), the engineer is considering upgrading to: (a) ductile iron Grade 80-55-06, (b) ADI Grade 900-650-09. For each: calculate the approximate weight reduction possible if the cross-section is scaled for the same bending strength (σ_y proportional), and identify any manufacturing changes needed.

3. Compare grey cast iron and ductile cast iron for engine crankshaft manufacture. Consider: (a) tensile strength and safety factor at peak combustion load, (b) fatigue resistance (grey iron flakes as crack initiators), (c) machining cost (grey iron's graphite lubrication advantage), (d) casting process modifications needed for ductile iron (Mg treatment). Conclude: which would you specify for a 500 hp performance engine?

4. White cast iron is used as the core of a bimetallic roll for a paper mill. The white iron core (Fe₃C matrix, 700 HV) provides wear resistance; the grey iron outer shell provides toughness. (a) How is this bimetallic structure made (hint: two-pour process or partial chilling)? (b) What happens at the interface between white and grey iron? (c) If the roll must be cut to final dimension, what machining challenges does the white iron core present?

5. CGI cylinder blocks for modern diesel engines require Mg content to be controlled within ±0.003% during production. (a) Why is this tolerance so tight (hint: low Mg → flake graphite; high Mg → nodular graphite; both bad for CGI)? (b) What manufacturing challenges arise from this tight control window (hint: Mg evaporates from melt → composition shifts over time)? (c) How can real-time graphite morphology be monitored in the foundry (thermal analysis method)?
