# Chapter 44: Gamma Prime (γ′) — The Strengthening Phase That Powers Jets

> **"Gamma prime is perhaps the most important microstructural feature in all of engineering. A tiny ordered cube, 200–500 nm on a side, that sits coherently in a nickel matrix, grows stronger as temperature rises, and resists the passage of dislocations in ways that would be impossible in any disordered phase. It is the reason jet engines work."**

---

## Table of Contents

1. [The Two-Phase γ/γ′ System — Overview](#1-the-two-phase-γγ-system--overview)
2. [Crystal Structure of γ′ — The L1₂ Structure](#2-crystal-structure-of-γ--the-l1-structure)
3. [Coherency — Why γ′ Sits So Perfectly in the Matrix](#3-coherency--why-γ-sits-so-perfectly-in-the-matrix)
4. [The Anomalous Yield Strength Effect — Stronger at Higher T](#4-the-anomalous-yield-strength-effect--stronger-at-higher-t)
5. [How γ′ Blocks Dislocations](#5-how-γ-blocks-dislocations)
6. [γ′ Volume Fraction — More Is (Usually) Better](#6-γ-volume-fraction--more-is-usually-better)
7. [γ′ Size and Morphology — The Shape Matters](#7-γ-size-and-morphology--the-shape-matters)
8. [γ′ Coarsening — The Enemy of Creep Resistance](#8-γ-coarsening--the-enemy-of-creep-resistance)
9. [γ″ (Gamma Double Prime) — IN718's Strengthener](#9-γ-gamma-double-prime--in718s-strengthener)
10. [Controlling γ′ Through Heat Treatment](#10-controlling-γ-through-heat-treatment)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. The Two-Phase γ/γ′ System — Overview

A nickel superalloy contains two major phases in its microstructure:

**γ (gamma) matrix:**
- Disordered FCC solid solution of Ni with Cr, Co, Mo, W, Re, Ru, etc. dissolved
- Continuous phase — provides the matrix
- Relatively soft at temperature; provides ductility and toughness

**γ′ (gamma prime) precipitates:**
- Ordered L1₂ compound: primarily **Ni₃Al** with Ti, Ta, Nb, Hf substituting for Al
- Small, cuboidal particles, 50–600 nm on a side depending on heat treatment
- Discrete particles dispersed in the γ matrix
- **Hard, strong at high temperature** — the primary strengthener

The γ/γ′ system is unique to Ni-base alloys. The near-perfect lattice match between γ and γ′ means they coexist with extraordinary microstructural stability. This is the materials science foundation of the jet engine.

**Typical composition split between phases:**
- γ phase is enriched in: Cr, Co, Mo, W, Re, Ru (refractories and large atoms)
- γ′ phase is enriched in: Al, Ti, Ta, Nb, Hf (smaller atoms that fit into Ni₃Al structure)

In CMSX-4 at 980°C:
- γ volume fraction: ~40%
- γ′ volume fraction: ~60%
- γ′ particle size: 0.4–0.5 μm (cuboidal, in channels ~100 nm wide)

---

## 2. Crystal Structure of γ′ — The L1₂ Structure

The **L1₂ crystal structure** (also written as Cu₃Au type) is an **ordered** version of FCC:

```
     L1₂ (ordered):           Disordered FCC:
     
        Al    Ni    Al           X     X     X
       / |   / |   / |          / |   / |   / |
      Ni─●──Ni─●──Ni─●         X─●──X─●──X─●
     /|  |  /  |  /  |        /  |  /  |  /  |
    Ni─●──Ni─●──Ni─●         X─●──X─●──X─●
       
     ● = Ni at face centers     ● = any atom randomly
     All corner atoms = Al      Face/corner = statistical average

   In Ni₃Al: Al at cube corners, Ni at face centers
   (or equivalently: Ni at face centers of every unit cell)
```

**Key features of L1₂:**
- **Ordered**: specific sites are occupied by specific atom types
- **Coherent with FCC γ**: same crystal structure, nearly same lattice parameter
- **Stoichiometry**: 3 Ni atoms per 1 Al atom (Ni₃Al)

**Modifications:** Real γ′ in superalloys is not pure Ni₃Al. Al sites are partially replaced by:
- Ti (forms Ni₃Ti, L1₂): common in older alloys (Waspaloy, Udimet 720)
- Ta (forms Ni₃Ta components): more temperature-stable than Al alone
- Nb (with Ti): both in some alloys
- Hf: small amount improves γ′ stability and grain boundary

The more Ta and less Ti in γ′, the higher the solvus temperature (temperature at which γ′ dissolves) → better high-temperature stability.

---

## 3. Coherency — Why γ′ Sits So Perfectly in the Matrix

**Coherency** means the crystal planes of γ′ are continuous with those of the γ matrix — no abrupt change in plane spacing at the interface:

```
γ matrix (FCC):   • • • • • • • • • •
                  • • • • • • • • • •
                  • •[γ′ particle]• •
                  • • • • • • • • • •

Coherent interface: crystal planes bend slightly but remain continuous
→ no broken bonds across interface
→ very low interfacial energy (~15 mJ/m²)
→ γ′ doesn't grow and coarsen as fast as incoherent precipitates
→ particle remains very stable
```

**Lattice parameter mismatch (misfit, δ):**
```
δ = (a_γ′ - a_γ) / a_γ

For most Ni superalloys: δ = −0.05% to −0.5% (γ′ is slightly smaller than γ)
Positive misfit: γ′ slightly larger than γ
```

The misfit creates **coherency strains** — elastic strain fields around each γ′ particle. These strain fields extend into the γ matrix and are crucial for:
1. Determining γ′ shape (cube or sphere — see §7)
2. Hardening (dislocations must cut through the strain field)
3. Raft formation at high temperature under stress

**Signing of misfit matters:**
- Negative misfit (γ′ smaller): particles elongate perpendicular to tensile stress at high T ("N-type rafts" form)
- Positive misfit: particles elongate parallel to tensile stress

Modern alloy design carefully controls misfit to optimize both room-temperature strength and high-temperature raft morphology. Small negative misfit (δ ≈ −0.05% to −0.10%) is typical for optimized 2nd/3rd generation SX alloys.

---

## 4. The Anomalous Yield Strength Effect — Stronger at Higher T

This is one of the most remarkable properties in all of materials science.

**For most metals:** yield strength decreases monotonically with temperature.  
**For Ni₃Al (γ′):** yield strength INCREASES from room temperature to about 800°C, then decreases.

```
Yield Strength
  │          
  │          ●← Peak strength ~760–800°C
  │      ●   ●
  │   ●         ●
  │ ●               ●
  │                     ●  ← strength decreasing above peak
  │
  └────────────────────────────→ Temperature
   RT    200   400   600   800   1000°C
```

This is called the **anomalous yield strength increase** or the **flow stress anomaly**.

**Mechanism — Kear-Wilsdorf lock:**

In γ′ (L1₂), primary slip occurs on {111}⟨110⟩ systems (like FCC). But a dislocation moving on {111} in γ′ creates an **antiphase boundary (APB)** — a local region where like-atoms are on wrong sites — because the L1₂ ordering means that moving exactly one Burgers vector breaks the ordered pattern.

To restore order, dislocations move in **pairs** (superpartials). The APB energy between the pair drives them together. But cross-slip of part of the pair onto {100} planes (which have lower APB energy in L1₂) creates a **Kear-Wilsdorf lock** — a segment of dislocation that cannot move on either plane.

As temperature rises:
- More dislocations cross-slip → more locks → harder γ′

Above the peak temperature:
- Thermal energy enables lock unlocking → softening dominates

**Practical consequence:**
γ′ doesn't just resist plastic deformation at room temperature — it becomes even more resistant as temperature rises to 800°C. This is exactly the regime where turbine blades operate. The alloy gets harder in the operating range. This is an extraordinarily valuable property found in essentially no other strengthening system.

---

## 5. How γ′ Blocks Dislocations

Two mechanisms operate depending on γ′ particle size and temperature:

### 5.1 Cutting (small γ′, low T)

When γ′ is small (< ~100 nm) or when stress is high, dislocations cut through the γ′ particles:

```
γ matrix: → → → → → [γ′] → → → → → →
Dislocation:         ↓ cuts through
             → → → ↗        ↘ → → →
```

Cutting requires:
- Overcoming the **antiphase boundary energy** (typically 200–250 mJ/m² for Ni₃Al)
- Creating a pair of superpartials (second dislocation follows closely behind)

The cutting stress increases with increasing γ′ volume fraction and size (up to a critical size). This is why **aging** grows γ′ initially improves strength.

### 5.2 Orowan Bypassing (large γ′, high T)

When γ′ is large (> ~200–300 nm) or widely spaced, it's harder to cut and dislocations bow around the particles:

```
γ channel:   →  →  →  →  →  →
                    ╭─────────╮ ← γ′ particle
Dislocation: ──────→ bows around particle ←──────
                    ╰─────────╯
                       → loop left behind
```

The Orowan stress: `σ_Orowan = M × 0.84 Gb/λ`

where λ is the inter-particle channel width (spacing between γ′ particles in the γ channel).

**Narrower channels = higher Orowan stress = stronger alloy.**

For typical CMSX-4 at 750°C:
- γ′ cube size: ~0.5 μm
- Channel width (γ channel): ~0.1 μm
- Orowan stress: ~500–700 MPa contribution to yield strength

As temperature rises above 850°C, dislocations prefer to **climb** over γ′ (using vacancies) rather than bypass by Orowan. At this point, minimizing climb rate becomes critical — this is why Re and other sluggish-diffusing elements (that slow climb) are so important.

---

## 6. γ′ Volume Fraction — More Is (Usually) Better

The **volume fraction of γ′ (f_γ′)** is one of the most important alloy design parameters:

```
Creep strength ↑ with f_γ′ (up to ~70%)
Oxidation resistance ↓ with high Al (needed to form more γ′)  [but TBC + aluminide coat compensate]
Density ↑ with more Ta/Re (heavy elements that stabilize γ′ and γ)
Ductility ↓ with very high f_γ′
```

| Alloy Generation | f_γ′ (approx.) | γ′ Chemistry |
|-----------------|----------------|-------------|
| Early polycrystal (IN100) | 55–60% | Ni₃(Al,Ti) |
| 1st gen SX (PWA 1480) | 55% | Ni₃(Al,Ti,Ta) |
| 2nd gen SX (CMSX-4) | 65% | Ni₃(Al,Ti,Ta) + Re in γ |
| 3rd gen SX (CMSX-10) | 65% | Higher Ta; more stable |
| 4th gen SX (TMS-138) | 65–70% | Ru addition; TCP suppression |

Increasing from 55% to 65% γ′ contributed significantly to creep improvement from 1st to 2nd generation SX alloys.

Above ~70% γ′, the "matrix" is insufficient to provide a continuous network → alloy can't accommodate thermal stresses → brittle. So there's a practical upper limit.

---

## 7. γ′ Size and Morphology — The Shape Matters

γ′ particles adopt different shapes depending on misfit and volume fraction:

```
Misfit/f_γ′:    Low misfit, low f_γ′      High misfit, high f_γ′
                     Spherical                  Cuboidal

                       ○                       □ □ □
                       ○ ○                     □ □ □
                       ○                       □ □ □

                Best for moderate T        Best for HPT blade applications
```

**At operating temperature (800–1000°C) under stress:**
Cuboidal γ′ particles **raft** — they merge into plate-like structures perpendicular to the applied tensile stress. This is called "N-type rafting" for negative misfit:

```
Before rafting:              After rafting (at T under stress):
□ □ □ □                     ══════════════
□ □ □ □           →          ══════════════   ← long γ′ rafts
□ □ □ □                     ══════════════
                                 ↑
                          γ channels are long and narrow
```

**Does rafting help or hurt creep?**
- Rafting creates a more tortuous path for dislocation climb → initially slows creep (good)
- Fully rafted structure can become a barrier to further deformation
- Eventually, raft disintegration and topological inversion → accelerated creep (bad)

Engineers design for a specific raft morphology and the blade is retired before raft breakdown occurs. Raft monitoring via TEM on service-exposed blades is a research tool for understanding component condition.

---

## 8. γ′ Coarsening — The Enemy of Creep Resistance

Even within a single-phase field, γ′ particles slowly **coarsen** over time — larger particles grow at the expense of smaller ones. This is **Ostwald ripening**:

Driving force: minimizing total interfacial energy. Large particles have less surface/volume ratio → lower total energy.

**Lifshitz-Slyozov-Wagner (LSW) theory:**
```
r³(t) - r³(0) = (8/9) × (D × γ_int × V_m) / (R T) × t

r = average particle radius at time t
D = diffusion coefficient (key!)
γ_int = γ/γ′ interfacial energy
V_m = molar volume
```

Key insight: **r³ ∝ t** (cube law) → the particle radius grows slowly but eventually becomes significant.

**Effect of coarsening on strength:**
As γ′ coarsens:
- γ channels become wider (lower Orowan stress)
- More particles can be cut (fewer Orowan bypasses needed)
- Creep rate increases
- **Blade life decreases**

**How coarsening is slowed:**
1. **Decrease interfacial energy** γ_int: control misfit (keep small negative misfit)
2. **Decrease diffusion coefficient** D: add elements that diffuse slowly in Ni (Re is the most effective — Re has extremely low D_Re in Ni, and Re clusters in γ preferentially at γ/γ′ interfaces slowing the diffusion that enables coarsening)
3. **Keep temperature as low as possible**: D ∝ exp(-Q/RT) — even 20°C lower dramatically slows coarsening

This is the **third reason why rhenium is so valuable**: beyond its direct solid-solution strengthening and dislocation climb retardation, Re dramatically slows γ′ coarsening, extending blade life at temperature.

---

## 9. γ″ (Gamma Double Prime) — IN718's Strengthener

IN718 (the most used superalloy by weight) uses a different precipitate: **γ″ (Ni₃Nb)**.

| Property | γ′ (Ni₃Al) | γ″ (Ni₃Nb) |
|----------|-----------|-----------|
| Structure | L1₂ (ordered FCC) | D0₂₂ (ordered bct) |
| Shape | Cuboidal | Disc-shaped (oblate spheroid) |
| Misfit with γ | Small (~−0.05%) | Large (~+2.9%) |
| Max service T | ~1000°C | ~650°C |
| Anomalous yield | Yes | Less pronounced |
| Stability | High | Metastable (converts to δ) |

**Why γ″ is limited to lower temperatures:**
γ″ is metastable — above about 650°C, it converts to the stable orthorhombic **δ phase (Ni₃Nb)** which provides little strengthening. This is why IN718 is limited to ~650°C service, making it unsuitable for HPT blades but excellent for cooler disks and compressor components.

The large misfit (+2.9%) of γ″ with the γ matrix creates very strong hardening (misfit × elastic modulus contributions) but also drives rapid transformation to δ above the stability limit.

---

## 10. Controlling γ′ Through Heat Treatment

The γ/γ′ microstructure is almost entirely controlled by heat treatment:

**Step 1 — Solution Treatment:**
Heat to above γ′ solvus (typically 1270–1340°C for most SX alloys) → all γ′ dissolves.
- Purpose: homogenize chemistry, dissolve any coarse as-cast γ′ (which is irregular)
- Single-crystal specific: solution treatment must NOT recrystallize the crystal

**Step 2 — Primary Aging (Precipitation):**
Cool rapidly (controlled cool or air cool) to precipitate γ′ at the optimal size.  
Typical: 1100–1130°C for 4–6 hours → coarse primary γ′ (~0.5 μm cubes)
- Controls the coarse γ′ population

**Step 3 — Secondary Aging (fine γ′):**
870°C for 20 hours (typical) → precipitates a fine secondary γ′ population within the γ channels.
- Fills the γ channels with fine γ′ → narrows effective channel width → raises Orowan stress

**Optimized bimodal γ′ distribution:**
The best creep resistance comes from a **bimodal** distribution:
- Large (primary) γ′: ~0.5 μm, force Orowan bypassing → long dislocation paths
- Small (secondary) γ′: ~50–100 nm, within channels, further block movement

This bimodal structure is the current state of the art in Ni SX superalloy heat treatment.

---

## Summary

| Topic | Key Point |
|-------|-----------|
| γ matrix | Disordered FCC; Cr, Co, W, Mo, Re, Ru dissolved |
| γ′ precipitate | Ordered L1₂ (Ni₃Al); cuboidal; 50–600 nm |
| Coherency | γ/γ′ planes continuous; low interfacial energy; very stable |
| Misfit (δ) | Small negative; controls shape, rafting, coarsening |
| Anomalous yield | γ′ strengthens up to ~800°C (Kear-Wilsdorf lock); unique in materials science |
| Dislocation blocking | Small γ′: cutting (APB energy); large γ′: Orowan bypassing (channel width) |
| Volume fraction | ~60–70% γ′ in best SX alloys; more = stronger up to limit |
| Coarsening | r³ ∝ t (LSW law); slowed by low interfacial energy + Re (slow D) |
| Rafting | γ′ merges into plates under stress at high T; initially helps, eventually harmful |
| γ″ | Ni₃Nb in IN718; disc-shaped, large misfit; metastable above 650°C; converts to δ |
| Heat treatment | Solution + primary age (coarse γ′) + secondary age (fine γ′) → bimodal optimum |

**Next chapter:** With the γ/γ′ microstructure established, we look at the surface of the superalloy — oxidation and hot corrosion, which attack without regard for microstructural excellence. The alloy's bulk properties are irrelevant if the surface is consumed.

---

## Exercises

1. In a Ni superalloy with 65% γ′, the channel width between γ′ particles is 0.1 μm. Calculate the Orowan stress using σ_Or = 0.84Gb/λ, where G = 80 GPa and b = 0.25 nm. Compare to the room-temperature yield strength of annealed copper (~70 MPa). What does this tell you about the strengthening?

2. γ′ coarsening follows r³(t) = r³(0) + K×t. If at 950°C, initial γ′ radius = 200 nm and after 1000h it reaches 240 nm, calculate K. Predict the γ′ size after 10,000h at 950°C. At 1000°C, K doubles (faster diffusion). Predict size after 1,000h at 1000°C.

3. Explain the Kear-Wilsdorf lock mechanism in your own words. Why does it cause yield strength to INCREASE with temperature up to ~800°C? What happens above that temperature?

4. Compare the heat treatment of γ′ in CMSX-4 (SX blade alloy) vs. γ″ in IN718 (disk alloy): (a) solution temperatures, (b) aging temperatures, (c) maximum service temperatures. Why is IN718 not used for rotating HPT blades?

5. In positive-misfit alloys, γ′ rafts form parallel to the stress direction. In negative-misfit alloys, γ′ rafts form perpendicular. Explain qualitatively why this is, using the elastic strain energy of the coherent interface.
