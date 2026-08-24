# Chapter 16: Creep — Slow Deformation at High Temperature

> **"Creep is the reason jet engine turbine blades have a finite life. At operating temperature, the blade slowly elongates, grain by grain, bond by bond, diffusion step by diffusion step — invisibly, inexorably, until it reaches its dimensional limit. Understanding creep is understanding the fundamental life-limiting mechanism of every hot-section component in an engine."**

---

## Table of Contents

1. [What Is Creep?](#1-what-is-creep)
2. [When Does Creep Matter?](#2-when-does-creep-matter)
3. [The Creep Curve — Three Stages](#3-the-creep-curve--three-stages)
4. [Mechanisms of Creep](#4-mechanisms-of-creep)
5. [Diffusion Creep — Nabarro-Herring and Coble](#5-diffusion-creep--nabarro-herring-and-coble)
6. [Dislocation Creep — Climb and Glide](#6-dislocation-creep--climb-and-glide)
7. [Grain Boundary Sliding](#7-grain-boundary-sliding)
8. [The Larson-Miller Parameter — Practical Life Prediction](#8-the-larson-miller-parameter--practical-life-prediction)
9. [Creep Rupture — When Creep Ends in Failure](#9-creep-rupture--when-creep-ends-in-failure)
10. [Designing Against Creep — The Metallurgist's Strategies](#10-designing-against-creep--the-metallurgists-strategies)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What Is Creep?

**Creep** is the time-dependent, permanent deformation of a material under constant stress, at temperatures above roughly **0.3–0.4 × T_melt** (the homologous temperature).

At room temperature, steel is stable under constant load — it doesn't slowly stretch. But put it at 600°C and apply a steady stress, and it will slowly elongate over days, weeks, and years. This is creep.

The key conditions:
1. **Elevated temperature** (T > 0.3 T_melt): atoms have enough thermal energy to diffuse and climb
2. **Sustained stress**: even below the yield strength, given enough time
3. **Time**: creep is a rate process, not instantaneous

### Homologous Temperature

The **homologous temperature** T_H = T/T_melt (where T is in Kelvin) normalizes temperature relative to melting point. Creep becomes significant at T_H > 0.3–0.4 for most metals:

| Metal | T_melt (K) | 0.4 T_melt (°C) | Creep significant above |
|-------|-----------|----------------|------------------------|
| Al | 933 | 100°C | 100°C |
| Pb | 600 | 27°C | Room temperature! |
| Fe | 1811 | 452°C | ~450°C |
| Ni | 1728 | 419°C | ~400°C |
| W | 3695 | 1208°C | >1200°C |

Lead pipes in old plumbing slowly creep under their own weight — that's why they develop bulges. This happens at room temperature because room temperature is 0.5 T_melt for lead.

Nickel superalloys in a jet engine operate at 850–1050°C, which is T_H = 0.61–0.75 — deep into the creep regime. Creep is not just a concern; it is the primary design-limiting mechanism.

---

## 2. When Does Creep Matter?

Creep is critical in:

| Application | Temperature | T_H |
|-------------|-------------|-----|
| Jet engine turbine blades | 850–1050°C | 0.65–0.75 |
| Gas turbine disks | 650–750°C | 0.55–0.60 |
| Nuclear fuel cladding | 300–400°C | 0.45–0.55 |
| Steam turbine blades | 550–620°C | 0.50–0.55 |
| Refinery pressure vessels | 450–560°C | 0.40–0.48 |
| Tungsten filaments | 2200°C | 0.60 |
| Lead roofing | Room temp | 0.50 |

Any structural component operating above T_H ≈ 0.4 must be creep-designed, not just statically stress-designed.

---

## 3. The Creep Curve — Three Stages

A **creep test** applies a constant tensile stress at constant temperature and measures strain vs. time. The result is the characteristic three-stage creep curve:

```
Strain (ε)
  │                                         ●── Fracture
  │                                    ─────
  │                                ─────    ← Stage III (tertiary/accelerating creep)
  │                          ──────          High ε_dot, damage accumulating
  │              ─────────────               Stage II (secondary/steady-state creep)
  │           ───                            Constant ε_dot (minimum creep rate)
  │       ──                                 Stage I (primary/transient creep)
  │  ──                                      Decreasing ε_dot
  │  ↑
  │ Instantaneous
  │ elastic strain
  │ at load application
  └─────────────────────────────────────────→ Time
```

### Stage I — Primary Creep

**Decelerating strain rate.** The material initially creeps fast but slows down. Work hardening and microstructural rearrangement occur (dislocations rearrange into lower-energy configurations). The rate continuously decreases.

### Stage II — Secondary (Steady-State) Creep

**Constant strain rate (ε̇_min).** Work hardening is exactly balanced by thermally activated recovery (dislocation climb, annihilation). This is the longest and most important stage.

The **minimum creep rate** ε̇_min is the key engineering parameter:
- All design life calculations use ε̇_min
- Blade life = (allowable creep strain) / ε̇_min

### Stage III — Tertiary Creep

**Accelerating strain rate** leading to fracture. Causes:
- Microstructural degradation (precipitate coarsening → loss of strength)
- Grain boundary cavitation and cracking
- Necking (reducing cross-section → increasing true stress)
- Oxidation penetration

When Stage III begins, the component is essentially failed from a life prediction standpoint, even if fracture hasn't yet occurred.

---

## 4. Mechanisms of Creep

Creep occurs by different atomic-scale mechanisms depending on temperature and stress:

```
Stress
  │ High        Power-law breakdown
  │
  │             Dislocation creep (power-law creep)
  │ Moderate    (glide + climb)
  │
  │             Grain boundary sliding
  │ Low
  │             Diffusion creep (Nabarro-Herring, Coble)
  │
  └──────────────────────────────────────────→ Temperature
          Low T            High T
```

Each mechanism has a characteristic rate equation (constitutive law). Knowing which mechanism dominates tells you how to design against it.

---

## 5. Diffusion Creep — Nabarro-Herring and Coble

At **low stress and high temperature**, creep occurs by diffusion of atoms themselves, driven by stress gradients.

**Applied stress creates a vacancy gradient:**
- Tensile stress regions: high vacancy concentration
- Compressive stress regions: low vacancy concentration
- Result: net atomic flux from compressive to tensile → material elongates

**Nabarro-Herring Creep (bulk/lattice diffusion):**
```
ε̇ = A_NH × (D_L × σ × Ω) / (k × T × d²)

where:
  D_L = lattice diffusion coefficient
  σ = applied stress
  Ω = atomic volume
  d = grain diameter
  k = Boltzmann constant
```

**Key insights:**
- Rate ∝ **D_L** (lattice diffusion coefficient) → increases exponentially with T
- Rate ∝ **σ** (linear stress dependence — not power law!)
- Rate ∝ **1/d²** → fine grains creep faster! This is the opposite of room-temperature where fine grains are stronger.

**Coble Creep (grain boundary diffusion):**
```
ε̇ = A_C × (D_GB × σ × Ω × δ) / (k × T × d³)

where:
  D_GB = grain boundary diffusion coefficient
  δ = grain boundary width (~0.5 nm)
```

Coble creep dominates at lower temperatures (relative to N-H) because grain boundary diffusion activates at lower T.

- Rate ∝ **1/d³** → even stronger grain size dependence
- Eliminating grain boundaries (single crystal) eliminates both N-H and Coble creep

**Why single crystals beat polycrystals in creep:**
In a single crystal, there are no grain boundaries — neither N-H nor Coble creep can operate. The only remaining mechanisms are dislocation-based, which requires higher stress to activate. This is the fundamental physical reason why single-crystal turbine blades have superior creep resistance.

---

## 6. Dislocation Creep — Climb and Glide

At **moderate to high stress and high temperature**, creep occurs by dislocation motion. The key process is **dislocation climb** — movement of dislocations perpendicular to their slip plane, enabled by vacancy diffusion.

**Why climb is special:**
At room temperature, dislocations can only glide along their slip plane. When they encounter an obstacle (precipitate, forest dislocation), they're stuck — they can't move sideways off the plane.

At high temperature, vacancies diffuse rapidly. By absorbing or emitting vacancies, a dislocation moves off its original slip plane (**climbs**). Once it's climbed past the obstacle, it can glide again until it hits the next obstacle. This cycle of glide → obstacle → climb → glide again is **power-law creep**.

```
Power-law creep:
ε̇ = A × D_L × (σ/E)^n × exp(-Q_c/RT)

where:
  n = stress exponent (typically 3–8 for metals)
  Q_c = creep activation energy (close to bulk diffusion activation energy)
  E = Young's modulus
```

The **stress exponent n** is crucial:
- n ≈ 1: diffusion creep (N-H or Coble)
- n ≈ 3: solid solution alloys (viscous glide)
- n ≈ 4–6: pure metals and simple alloys (dislocation climb)
- n ≈ 8–20: precipitation-hardened alloys and superalloys (Orowan bypass + climb)

Higher n means creep rate is more stress-sensitive — a 10% reduction in stress causes a much larger reduction in creep rate (proportional to σ^n). This is why operating turbine blades at slightly lower stress dramatically extends life.

**Effect of γ′ precipitates on dislocation creep:**
In Ni superalloys, dislocations must either cut through γ′ particles (which requires overcoming the antiphase boundary energy) or climb around them. At low temperature, they cut. At high temperature, they climb. Optimizing the γ′ size (coarse enough to force climb rather than cutting, fine enough to maximize Orowan stress) is a key alloy design parameter.

---

## 7. Grain Boundary Sliding

At high temperature, grain boundaries can **slide** — one grain shifts relative to its neighbor. This mechanism is accommodated by:
- Diffusion creep at the triple junctions where three grains meet
- Local deformation within grains to maintain continuity

Grain boundary sliding (GBS) leads to:
- Cavities forming at triple junctions and boundary ledges
- Eventual grain boundary cracking
- Reduced ductility and toughness

**Important:** GBS occurs at all grain boundaries — but it's most damaging at boundaries **perpendicular to the stress direction** (transverse grain boundaries). These boundaries get pulled apart by the tensile stress.

In a polycrystalline turbine blade, transverse grain boundaries are the primary sites of creep damage. **Directional solidification** (Chapter 48) eliminates transverse boundaries — all boundaries are parallel to the centrifugal force. **Single crystal** casting eliminates all boundaries.

Quantitatively, eliminating transverse boundaries by going from equiaxed → DS improves creep life by 50–100°C equivalent. Going from DS → single crystal improves by another 30–50°C.

---

## 8. The Larson-Miller Parameter — Practical Life Prediction

Creep testing at service temperature (900–1000°C) takes too long for every design iteration. The **Larson-Miller parameter** (LMP) allows accelerated testing at higher temperature to predict long-term behavior:

```
LMP = T × (log t_r + C)

where:
  T   = absolute temperature (°C + 273, in Rankine or Kelvin)
  t_r = time to rupture
  C   = material constant (typically ~20 for Ni superalloys)
```

The LMP is approximately constant for a given stress level:

```
At high T, short time:    T_high × (log t_short + C) ≈ LMP
At low T, long time:      T_low × (log t_long + C) ≈ LMP
```

**How it's used:**
1. Test alloy at 1100°C for 100 hours under stress σ → get LMP
2. Use the same LMP to predict life at 950°C (service temperature)
3. t_rupture at 950°C can be read from the LMP curve

```
LMP Curve:
  σ (MPa)
  │
  │ ────────────────────────────── (each curve = one alloy)
  │                               ╲
  │                                ╲
  │                                 ╲
  │                                  ╲
  │
  └─────────────────────────────────→ LMP = T(log t_r + C)
           Low T/long time        High T/short time
```

The LMP is used in every turbine blade material specification. Manufacturers plot LMP curves for their alloys and guarantee a minimum LMP value that ensures adequate creep life in service.

---

## 9. Creep Rupture — When Creep Ends in Failure

**Creep rupture** (stress rupture) is fracture that occurs at the end of Stage III. The time to rupture (t_r) at given stress and temperature is the critical design parameter.

**Rupture mechanisms:**
1. **Grain boundary cavitation**: vacancies coalesce at boundaries → voids → cracks → fracture along boundaries (intergranular fracture)
2. **Wedge cracking at triple junctions**: grain boundary sliding concentrates stress
3. **Denuded zones**: precipitates dissolve near grain boundaries → local softening
4. **Oxidation-creep synergy**: cracks at surface allow oxygen in → local embrittlement → deeper cracks → failure accelerated

**Creep rupture ductility** — the %EL at rupture under creep conditions — is generally much lower than tensile ductility. High-temperature failures of turbine blades are often intergranular with < 5% creep ductility.

---

## 10. Designing Against Creep — The Metallurgist's Strategies

Every strategy maps onto reducing one or more creep mechanisms:

| Strategy | Effect | Mechanism Targeted |
|----------|--------|-------------------|
| **Increase grain size** | Reduce N-H and Coble creep (∝ 1/d², 1/d³) | Diffusion creep |
| **Eliminate grain boundaries (DS/SX)** | Eliminate GBS and intergranular damage | All GB mechanisms |
| **Add precipitates (γ′, γ″)** | Block dislocation climb and glide | Dislocation creep |
| **Large, stable precipitates** | Force Orowan climbing (needs full climb) | Dislocation creep |
| **Add solid solution strengtheners** (Mo, W, Re) | Slow diffusion, increase lattice friction | Dislocation creep |
| **Reduce operating temperature** | Exponential reduction in D and ε̇ | All T-dependent |
| **Reduce operating stress** | Reduce ε̇ ∝ σ^n | All stress-dependent |
| **Add Re (rhenium)** | Reduce diffusivity by 2–5×; slow climb | Dislocation climb |
| **Add Ru (ruthenium)** | Stabilize γ′ and prevent TCP phases | Microstructure stability |

**The rhenium effect in superalloys:**
Rhenium at 3–6 wt% reduces the creep rate of Ni superalloys by 5–10× at 1000°C. The mechanism involves:
- Re clusters in the γ matrix → very slow diffusing species → slow climb
- Re increases the lattice friction to dislocation glide
- Re stabilizes γ′ precipitate morphology (resists coarsening)

This single element addition enabled the jump from 1st to 2nd generation single-crystal superalloys and roughly 30°C additional temperature capability.

---

## Summary

| Topic | Key Point |
|-------|-----------|
| What is creep | Time-dependent plastic deformation at T > 0.3–0.4 T_melt |
| Three stages | Primary (decelerating) → Secondary (steady-state, ε̇_min) → Tertiary (accelerating → rupture) |
| Diffusion creep | N-H (lattice diffusion, ∝1/d²) and Coble (GB diffusion, ∝1/d³); linear in stress |
| Dislocation creep | Climb + glide; power-law (ε̇ ∝ σ^n, n=4-8); dominant at moderate-high stress |
| GB sliding | Grain boundaries slide; transverse boundaries most damaging; eliminated by DS/SX |
| Larson-Miller | LMP = T(log t_r + C); extrapolates short tests to long service lives |
| Design strategies | Large grains; precipitates; DS/SX; solid solution; lower T and σ |
| Rhenium | 3–6% Re reduces creep rate 5–10×; enables Gen 1→2 SX jump |
| Turbine blades | T_H = 0.65–0.75; creep is the primary life-limiting mechanism |

**This chapter is the key to understanding why single-crystal superalloys exist.** Every strategy against creep ultimately points to the same conclusion: eliminate grain boundaries (remove N-H, Coble, GBS mechanisms) and engineer a perfect precipitate structure (maximize dislocation creep resistance). That is exactly what a Gen 4 single-crystal superalloy does.

**Next chapter:** We have the material. Now we learn to control it — heat treatment. TTT and CCT diagrams are the process engineer's version of the phase diagram: not what's stable at equilibrium, but what actually forms when we heat and cool at real rates.

---

## Exercises

1. A Ni superalloy turbine blade operates at 950°C under 150 MPa for 20,000 hours. The minimum creep rate ε̇_min = 2×10⁻¹⁰ s⁻¹. Calculate the total creep strain accumulated. If the allowable blade elongation is 2%, is the blade still within spec?

2. Diffusion creep rate varies as 1/d² (N-H) or 1/d³ (Coble). If you double the grain size from 0.5 mm to 1.0 mm, by what factor does N-H creep rate decrease? By what factor does Coble creep rate decrease?

3. A polycrystalline alloy has a stress exponent n=5. If you reduce the applied stress from 200 MPa to 180 MPa (a 10% reduction), by what factor does the creep rate change?

4. Use the Larson-Miller parameter: alloy A has C=20. A 1000-hour test at 1100°C under 200 MPa gives rupture. Calculate LMP. What life does this predict at 950°C under the same stress?

5. Explain, using mechanism arguments, why going from equiaxed polycrystal → directionally solidified (columnar) → single crystal each improves creep life, and which mechanisms are eliminated at each step.
