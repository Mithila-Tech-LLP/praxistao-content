# Chapter 33: Casting Processes — From Liquid Metal to Shape

> **"Every complex metal part you have ever seen — an engine block, a turbine casing, a hip implant, a jet blade — began as liquid metal poured into a void. The void defines the shape; the cooling defines the microstructure; the design of the void defines the quality. Casting is the most direct route from liquid metal to near-net shape, and it has been practiced for 8,000 years. Modern variants — vacuum investment casting for single crystals — would astonish ancient bronzesmiths, but the principle is identical."**

---

## Table of Contents

1. [Why Casting?](#1-why-casting)
2. [Sand Casting](#2-sand-casting)
3. [Die Casting](#3-die-casting)
4. [Investment Casting (Lost-Wax)](#4-investment-casting-lost-wax)
5. [Centrifugal Casting](#5-centrifugal-casting)
6. [Continuous Casting of Steel](#6-continuous-casting-of-steel)
7. [Vacuum Casting for Superalloys](#7-vacuum-casting-for-superalloys)
8. [Solidification Defects and Quality Control](#8-solidification-defects-and-quality-control)
9. [Summary](#summary)
10. [Exercises](#exercises)

---

## 1. Why Casting?

Casting allows:
- Complex internal geometry (cores create hollow sections, cooling passages)
- Near-net shape (minimal machining needed)
- Large single pieces (engine blocks, turbine casings, ship propellers)
- Low material waste (vs. machining a shape from a billet)
- Economical at all scales (from 1 prototype to millions of parts)

**Limitations:**
- Surface finish inferior to machining
- Solidification defects possible (porosity, shrinkage, segregation)
- Mechanical properties typically lower than wrought (large grain size, potential porosity)
- Thin sections challenging (solidification before fill)

---

## 2. Sand Casting

**The oldest and most versatile casting process.** Used for everything from 10-gram decorative pieces to 100-tonne industrial parts.

**Process:**
```
1. Pattern (wood, plastic, metal)
2. Cope (top) + drag (bottom) flask filled with sand + binder
3. Compact around pattern
4. Remove pattern → hollow mold cavity
5. Insert cores (for internal passages)
6. Close mold
7. Pour liquid metal
8. Allow solidification
9. Shake out (break mold)
10. Cut off gates/risers
11. Clean, machine
```

**Sand types:**
- Green sand (silica sand + bentonite + water): most common; can be reused; decent surface finish
- No-bake (chemically bonded): better dimensional accuracy; higher strength; not reclaimable as easily
- Shell mold (sand + thermosetting resin baked on heated pattern): good surface finish; dimensional accuracy; smaller parts

**Key features:**
- Riser design: liquid reservoir feeding solidification shrinkage
- Gating system: runners, sprues, gates to control filling velocity and avoid turbulence
- Core prints: sand cores held in mold by recesses in mold cope/drag

**Applications:** Engine blocks, gearboxes, valve bodies, agricultural machinery, pump housings.

**Limitations:** Poor surface finish (Ra ~6–25 μm); moderate dimensional tolerance (±1 mm); not suitable for thin walls < 3 mm.

---

## 3. Die Casting

**High-pressure, high-speed injection of liquid metal into a reusable steel die.** Best for high-volume, small-to-medium parts.

**Process:**
- Aluminum, zinc, or magnesium alloys (not steel — too hot for steel die life)
- Metal injected at 20–140 MPa and high velocity (30–120 m/s)
- Die: hardened H13 tool steel; water cooled; 50,000–1,000,000 shots life

**Types:**
- **Hot chamber:** Metal chamber is submerged in melt → fast cycle; for low-melting alloys (Zn, Mg)
- **Cold chamber:** Ladled into shot chamber → then injected; for higher-T alloys (Al)

**Advantages:**
- Excellent dimensional accuracy (±0.05 mm)
- Good surface finish (Ra 1–3 μm)
- Very fast cycle time (5–60 seconds)
- Thin walls possible (0.5–1 mm)
- High volume economical

**Limitations:**
- Porosity from trapped air (rapid fill traps gas) → reduces fatigue life; cannot weld
- Alloy selection limited (Al, Zn, Mg, Cu)
- High tooling cost ($50,000–$500,000 per die)

**Applications:** Automotive (engine covers, brackets, transmission housings), electronic housings, power tools.

---

## 4. Investment Casting (Lost-Wax)

**The premier near-net-shape casting process for complex shapes.** Produces excellent surface finish and dimensional accuracy with no parting line.

**Process:**
```
1. WAX PATTERN INJECTION
   Wax (or plastic) injected into metal die → wax pattern of part
   
2. ASSEMBLY
   Multiple patterns assembled on wax tree (sprue + gates)
   
3. CERAMIC SHELL BUILDING
   Tree dipped in ceramic slurry (colloidal silica + zircon flour)
   Stuccoed with coarser ceramic particles
   Dried (24 hours per layer)
   Repeat 7–15 times → shell 5–10 mm thick
   
4. DEWAX
   Flash-fire (~180°C) autoclave: wax melts and drains out
   Leaves ceramic mold with fine internal detail
   
5. BURNOUT
   Heat to 1,050°C → burn out wax residue, sinter ceramic, preheat mold
   
6. POUR
   Liquid metal poured into hot mold
   For Ni superalloys: VACUUM casting (prevent oxidation)
   
7. KNOCKOUT
   Cool → break ceramic shell
   Cut off sprue and gates
   
8. FINISHING
   Blasting, machining datum surfaces, NDT inspection
```

**Why investment casting excels:**
- No parting line → undercuts possible → complex airfoil shapes
- Excellent surface finish (Ra 1–3 μm directly as cast)
- Dimensional tolerance ±0.1 mm
- Very thin walls possible (0.5–1 mm for DS/SX blade cooling passages)
- Ceramic core for internal cooling passages (colloidal silica + fused silica core)

**Applications:** Turbine blades (ALL HPT blades worldwide are investment cast), golf club heads, dental crowns, orthopedic implants, firearms components.

**Ceramic core for turbine blades:**
The internal cooling passages in HPT blades are created by a pre-formed ceramic core:
- Core: fused silica + colloidal silica binder → fired
- Core positioned in wax die → wax injected around core
- After metal pouring: core removed by chemical leaching (NaOH or KOH solution)
- Core defines the cooling channel geometry to < 0.5 mm precision

---

## 5. Centrifugal Casting

**Metal is poured into a rotating mold.** Centrifugal force:
- Forces metal against outer mold wall
- Forces less-dense inclusions and porosity to the inner surface (which is machined away)
- Produces dense, clean outer region

**True centrifugal casting (hollow cylinders):**
- Pipe, tubes: metal fills rotating cylinder → hollow pipe with good outer surface
- Centrifugal force creates back-pressure → dense metal against mold
- Used for: cast iron pipes (water/sewer), bi-metallic cylinder liners, pressure vessels

**Semicentrifugal casting:** Mold is symmetric but gravity-fed from center; used for wheels.

**Centrifuging (gang casting):** Multiple molds arranged around a central sprue; centrifugal force feeds molds; used for small precision parts (dental castings, jewelry).

---

## 6. Continuous Casting of Steel

**The dominant method for producing steel slabs, blooms, and billets.** Over 96% of world steel production is continuously cast.

```
Continuous casting schematic:

Ladle (liquid steel)
    ↓
Tundish (buffer vessel, distributes to molds)
    ↓
Copper mold (water-cooled, oscillating) ← solidification begins here
    ↓
Strand (partially solidified, liquid core)
    ↓
Secondary cooling (water sprays on strand)
    ↓
Pinch rolls (withdraw strand)
    ↓
Cut to length → slabs, blooms, billets
```

**What happens:**
- Steel solidifies in mold (skin ~20 mm thick)
- Strand emerges with solid shell + liquid core
- Secondary cooling: spray water hardens the strand progressively
- Complete solidification: 3–10 meters below mold depending on section size
- Speed: slabs at 1–2 m/min; billets at 3–6 m/min

**Key metallurgical issue: segregation**
- Negative segregation at surface (first to solidify = depleted in low-k elements)
- Positive segregation at centerline (last to solidify = enriched)
- Centerline porosity: gas + shrinkage at last point of solidification
- Managed by: electromagnetic stirring (homogenize liquid core), soft reduction (squeeze strand as it solidifies)

---

## 7. Vacuum Casting for Superalloys

**Turbine blades and other critical Ni superalloy parts require vacuum melting and casting.**

**Why vacuum?**
- Ni alloys contain Ti, Al, Hf (very reactive with O₂, N₂)
- Even 10 ppm O₂ → TiO₂ or Al₂O₃ inclusions → stress concentrators → fatigue failure
- Atmosphere must be < 10⁻³ Pa (10⁻⁵ bar) during melting + pouring

**Vacuum Induction Melting (VIM):**
```
- Induction coil in vacuum vessel
- Charge: pre-alloyed master heats or component elements
- Melt under vacuum → pour into investment mold (also in vacuum)
- For DS/SX: mold is in Bridgman furnace within the vacuum system
```

**Vacuum Arc Remelting (VAR):**
- Secondary melting to improve homogeneity and reduce macro-segregation
- Used for: Ti alloys (required), high-grade Ni superalloys (optional for special applications)
- Electric arc between consumable electrode and water-cooled Cu mold → controlled solidification

**DS and SX casting:** Covered in detail in Ch 48–49. Key point: controlled withdrawal rate and thermal gradient in Bridgman furnace defines grain structure.

---

## 8. Solidification Defects and Quality Control

**Shrinkage porosity:**
- Liquid → solid: ~3–5% volume decrease
- Last-to-freeze regions don't get fed by liquid → voids
- Detected by: X-ray radiography (shows as dark spots)
- Fixed by: riser design, DS (freeze from bottom to top)

**Gas porosity:**
- Dissolved H₂ (Al), CO (steel), N₂ come out of solution
- Spherical pores (smooth walls = gas vs. rough walls = shrinkage)
- Prevented by: vacuum casting, degassing (rotary degassing for Al)

**Misruns:** Metal solidifies before filling the mold completely:
- Cause: too low temperature, slow pour, cold mold
- Fix: higher T, faster pour, preheated mold

**Cold shuts:** Two metal streams meet but don't fuse:
- Cause: metal has cooled before meeting
- Detected by: surface inspection, X-ray

**Inclusions:** Trapped oxide films, slag particles, core material:
- Most dangerous for fatigue life
- Detected by: X-ray, ultrasonic testing

**Hot tears:** Solidification cracking (see Ch 08) — rupture while semi-solid

**Stray grains in DS/SX:** Equiaxed grains that nucleate ahead of the DS front:
- Detected by: Laue X-ray diffraction
- Cause: ceramic fragment in mold, rough surface, thermosolutal convection

---

## Summary

| Process | Materials | Tolerance | Surface Ra | Volume | Applications |
|---------|-----------|-----------|-----------|-------|-------------|
| Sand casting | Any metal | ±1 mm | 6–25 μm | Low–high | Engine blocks, large parts |
| Die casting | Al, Zn, Mg | ±0.05 mm | 1–3 μm | Very high | Automotive housings |
| Investment (lost-wax) | Any, especially Ni | ±0.1 mm | 1–3 μm | Low–medium | Turbine blades, implants |
| Centrifugal | Fe, Ni, Cu | ±0.5 mm | 3–6 μm | Medium | Pipes, cylinder liners |
| Continuous casting | Steel | — | — | Very high | Steel slabs, billets |
| Vacuum DS/SX | Ni superalloys | ±0.2 mm | 1–3 μm | Low | Turbine HPT blades |

---

## Exercises

1. A sand casting produces an aluminum alloy (356-T6) pump housing with shrinkage porosity in the thick boss section. (a) Draw schematically where the porosity is expected to form (hint: last to solidify = furthest from chill, most isolated). (b) Design a riser: calculate minimum riser volume for a 200 cm³ hot spot (Al shrinks 3.5%); riser itself has 15% waste (must be bigger). (c) Alternatively, specify where to place a chill to change the solidification direction. (d) What post-casting process eliminates residual porosity and allows T6 heat treatment?

2. A die-cast automotive bracket (Al A380, 1.5 kg) has 0.8% porosity detected by Archimedes method. (a) Calculate the number and total volume of voids assuming 0.8% volume. (b) The fatigue limit of a sound casting = 90 MPa; with 0.8% porosity = 65 MPa. Why do pores reduce fatigue so significantly? (c) The die is being redesigned with vacuum-assisted die casting (VADC). How does this reduce porosity? (d) Can the existing porous parts be HIP'ed? What limitations apply (wall thickness, pore surface contamination)?

3. Compare investment casting vs. machining from billet for a titanium hip implant (Ti-6Al-4V, mass = 80g): (a) Investment cast: 10% scrap metal in gates/risers; material cost $35/kg. Machining from billet: 85% material removal; billet cost $35/kg. Calculate material cost per part for each. (b) Machining cycle time = 2 hours at $150/h. Casting cycle time = 4 hours total / 20 parts = 12 min/part at $200/h (foundry overhead). Which process has lower cost per part? (c) Why is investment casting preferred for implants despite the machined surface finish being better?

4. A steel continuous casting machine produces 300 mm × 400 mm blooms at 0.8 m/min. (a) Calculate the production rate in tonnes/hour (steel density = 7.87 g/cm³). (b) Centerline segregation index SI (ratio of centerline composition to nominal) for Mn = 1.15 in this casting. What microstructure and property non-uniformity results from this? (c) Soft reduction (squeezing the strand in the last 20% of solidification) can reduce SI to 1.05. Calculate the improvement in homogeneity. (d) What heat treatment would be used after bloom rolling to reduce residual segregation?

5. During DS turbine blade casting, a stray grain appears in the airfoil section. Describe: (a) the three most likely causes of stray grain formation, (b) how the grain selector (spiral) prevents stray grains in the main airfoil but cannot prevent them once past the selector, (c) the X-ray Laue diffraction method used to detect stray grains, (d) the economic impact: if 8% of blades have stray grains and must be scrapped, and each blade costs $3,000 in materials + $2,000 in processing, calculate the cost per accepted blade.
