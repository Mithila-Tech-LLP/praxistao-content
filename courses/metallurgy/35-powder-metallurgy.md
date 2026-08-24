# Chapter 35: Powder Metallurgy — Building Parts Grain by Grain

> **"Powder metallurgy gives engineers a capability that casting and forging cannot: the ability to create alloys and composites that would segregate, react, or otherwise be impossible to make as a liquid. It is not a process for simple shapes — it is a process for impossible ones: tungsten carbide cutting tools, oxide-dispersion-strengthened superalloys, self-lubricating bearings, and the very finest-grained turbine disk alloys that hold jet engines together at speed."**

---

## Table of Contents

1. [Why Powder Metallurgy?](#1-why-powder-metallurgy)
2. [Powder Production Methods](#2-powder-production-methods)
3. [Powder Characterization](#3-powder-characterization)
4. [Compaction](#4-compaction)
5. [Sintering](#5-sintering)
6. [Hot Isostatic Pressing (HIP)](#6-hot-isostatic-pressing-hip)
7. [PM Superalloys for Turbine Disks](#7-pm-superalloys-for-turbine-disks)
8. [Cemented Carbides (WC-Co)](#8-cemented-carbides-wc-co)
9. [Oxide-Dispersion Strengthened (ODS) Alloys](#9-oxide-dispersion-strengthened-ods-alloys)
10. [Summary](#summary)
11. [Exercises](#exercises)

---

## 1. Why Powder Metallurgy?

Powder metallurgy (PM) routes produce metal parts by:
1. Making metal powder
2. Shaping the powder (compaction, injection, spray)
3. Densifying the powder (sintering, HIP, or combined)

**Why choose PM over casting/forging?**
- Alloy compositions impossible by melting (e.g., Cu-W is immiscible as a liquid; ODS materials with fine oxide dispersions at nm scale)
- Complex shapes with close tolerances (near-net-shape, minimal machining)
- Graded compositions (density gradient, functionally graded materials)
- Very fine grain size (rapid solidification of each droplet → fine structure)
- Controlled porosity (e.g., self-lubricating bearings at specific density)
- Better composition homogeneity (no dendritic segregation; each powder particle = a tiny ingot)

**PM limitations:**
- High cost of powder
- Part size typically limited (large HIP containers can be used but expensive)
- Residual porosity (if not fully densified by HIP)
- Complex shapes with thin sections can have ejection problems

---

## 2. Powder Production Methods

### Gas Atomization
**Most common for PM superalloys:**

```
Liquid metal stream → Inert gas jets (Ar, N₂) → Rapid solidification → Spherical powder particles

Cooling rate: 10³–10⁵ °C/s (each particle ~10–200 μm solidifies in milliseconds)
```

Results:
- Spherical particles (good flowability for die filling, HIP canning)
- Fine, uniform microstructure (no macro-segregation)
- Can be done under vacuum (VIGA — Vacuum Induction Gas Atomization) for reactive alloys
- Particle size: 10–300 μm, controlled by gas pressure + liquid stream diameter

### Centrifugal Atomization (Plasma Rotating Electrode Process — PREP)
- Bar electrode spun at high speed in inert atmosphere
- Plasma arc melts bar tip → centrifugal force atomizes → spherical, satellite-free powder
- Used for titanium (no contamination from nozzles)

### Water Atomization
- Cheaper than gas; irregular particle shape; good for iron, steel, Cu PM parts
- Particles often spongy → better green strength after compaction

### Chemical Reduction
- Oxides or salts reduced to metal powder (e.g., iron from iron oxide by H₂)
- Nickel powder from Ni carbonyl, Ni(CO)₄ → Ni + 4CO at 230°C (INCO process)
- Very fine (< 5 μm), irregular particles

### Mechanical Comminution
- Ball milling, rod milling: brittle materials (WC, oxides)
- High-energy ball milling: fracture + cold welding of particles → grain refinement → possible to make metastable phases

---

## 3. Powder Characterization

**Particle size distribution:** Measured by sieve analysis (> 45 μm), laser diffraction (0.1–2,000 μm):
- D10, D50, D90 values characterize the distribution
- Narrow distributions preferred for uniform sintering

**Particle shape:** Spherical (gas atomized) → good flowability; irregular (water atomized) → better interlocking in green compact

**Apparent density:** Density when powder freely poured (no packing) — spherical: ~60% of theoretical; irregular: 25–40%

**Flowability (Hall flowmeter):** Time for 50g to flow through a standard 2.5mm orifice; spherical > 30 g/s; irregular: sometimes won't flow at all

**Surface area (BET method):** Fine powders have large surface area → more surface oxide contamination → can affect sintering

**Purity / oxygen content:** Critical for PM superalloys — oxygen forms oxides at grain boundaries → embrittlement → must keep O₂ below 50 ppm in PM superalloy powder

---

## 4. Compaction

**Purpose:** Shape the powder and increase density before sintering.

### Die Compaction (Uniaxial)
Powder filled into rigid die, upper punch compresses → "green compact":
- Green density: 75–90% theoretical
- Friction from die walls causes density gradient (higher at punch faces, lower in middle)
- Lubricant (zinc stearate, wax) added to powder to reduce friction

**Isostatic pressing (CIP — Cold Isostatic Pressing):**
- Powder in rubber bag → pressure applied by water/oil from all sides
- No density gradient → more uniform green compact
- Can make complex shapes (long, L-shaped, tubes)
- Used for: WC-Co milling inserts, large superalloy preforms

### Metal Injection Molding (MIM)
- Fine powder (< 20 μm) + thermoplastic binder → feedstock
- Injected into complex die (like plastic injection molding)
- Debind (thermal or solvent) → sinter
- Complex 3D shapes; typical density: 96–99%
- Applications: watch parts, orthodontic brackets, small structural parts

### Tape Casting
- Powder suspended in slurry → cast into thin tape → stack → sinter
- Used for: ceramic substrates, solid oxide fuel cell components

---

## 5. Sintering

**Sintering:** Heating the green compact below the melting point → diffusion bonds particles → densification.

**Stages of sintering:**
```
Stage 1 (neck formation):
  Particles contact each other → surface diffusion forms "necks"
  Particle centers come closer → density increases slightly
  
Stage 2 (pore rounding):
  Pores become rounded (cylindrical → spherical)
  Grain boundaries form across necks
  Density: 70–92%
  
Stage 3 (pore elimination):
  Grain boundaries migrate → pores pinch off
  Isolated spherical pores → slow elimination
  Final density: 92–99% (difficult to reach 100% without HIP)
```

**Sintering parameters:**
- Temperature: typically 0.7–0.9 T_m(K)
- Time: 30 min to several hours
- Atmosphere: H₂, N₂, H₂-N₂ mix, vacuum — MUST prevent oxidation
- Heating rate: faster = less grain growth but may not de-bind uniformly

**Liquid phase sintering (LPS):**
- Small amount of liquid phase forms → capillary forces + enhanced diffusion → faster densification
- WC-Co: Co melts (T_sinter ~1,350°C, T_m Co = 1,495°C but liquidus with WC much lower)
- Cu-Ni-Fe automotive PM: small Cu-rich liquid aids densification

---

## 6. Hot Isostatic Pressing (HIP)

**HIP** applies simultaneous heat and isostatic pressure using inert gas (Ar) to fully densify PM parts:

```
HIP canister (steel or glass):
  ┌─────────────────────┐
  │   Powder or          │
  │   pre-sintered part  │
  └─────────────────────┘
           ↓
  High-pressure Ar (100–200 MPa)
  + High temperature (800–1,250°C)
  → Closes all pores → 100% density
```

**Why HIP achieves full density while sintering cannot:**
Sintering: gas in closed pores cannot escape → creates back-pressure resisting densification.
HIP: external pressure (100-200 MPa) >> internal pore pressure → pores collapse by creep + diffusion.

**HIP parameters for PM superalloys:**
- IN718 powder: 1,163°C / 103 MPa / 4 hours → 100% density
- René 88DT: 1,170°C / 103 MPa / 3 hours
- Result: ultrafine grain (ASTM 10–12, d ≈ 5–10 μm) vs conventional ingot (ASTM 6–8)

**HIP also used to:**
- Heal castings: close residual porosity in investment cast superalloy parts (without powder)
- Densify cermets, ceramics, composites

---

## 7. PM Superalloys for Turbine Disks

**Why turbine disks require PM route:**
- Disk alloys (René 95, René 88DT, Udimet 720Li, LSHR) contain >40% γ′:
  - High Al+Ti content → very high segregation coefficient → conventional large ingots crack on quench
  - Ingot route limited to ~20 cm diameter for these alloys
  - PM route: each gas-atomized powder particle ~50–100 μm → no macro-segregation regardless of final part size

**PM disk manufacturing route:**
```
VIM melt master alloy → Gas atomize (VIGA) → Screen to correct size fraction
→ Load into steel HIP can → Evacuate + weld → HIP (100% density)
→ Remove can → Forge/extrude → Isothermal forge disk → Solution treat + age
```

**PM superalloy grain size control:**
- Powder is refined (ASTM 12–14 grain, d = 3–6 μm) after HIP
- Supersolvus solution treatment (T > γ′ solvus) → abnormal grain growth → ASTM 6–8 → improves creep resistance
- Subsolvus solution treatment → retain fine grain → better fatigue life
- Dual microstructure disk: coarse grain rim (creep) + fine grain bore (fatigue) by controlled forging/heat treatment

---

## 8. Cemented Carbides (WC-Co)

**The most commercially important PM product:** ~20% of all PM production by value.

**Composition:**
- WC particles (hard phase, HV 2,400): 70–95%
- Co binder (tough phase): 5–30%

**Manufacturing:**
```
WC + Co powder → Ball mill (intimately mix) → Die press or CIP
→ Pre-sinter → Machine to near-net shape → LPS sinter at 1,350°C
→ HIP (optional, for better density) → Grind/EDM to final shape
```

**Microstructure after sintering:**
- WC: angular particles, 0.5–5 μm
- Co: continuous network between WC grains → provides toughness

**Properties:**
| Grade | WC% | Co% | HV | K_Ic (MPa√m) | Application |
|-------|-----|-----|----|--------------|-------------|
| Fine | 94 | 6 | 1,800 | 10 | Cutting inserts |
| Coarse | 85 | 15 | 1,200 | 18 | Mining drill bits |
| Tough | 80 | 20 | 1,000 | 22 | Cold heading tools |

**Why WC-Co works:** WC hardness gives wear resistance; Co binder absorbs crack energy → tough composite. Neither alone is useful (brittle WC; too-soft Co).

**Coating:** CVD TiC + Al₂O₃ + TiN on WC-Co inserts → improves oxidation resistance + hardness at cutting T.

---

## 9. Oxide-Dispersion Strengthened (ODS) Alloys

**Problem with conventional superalloys:** Above ~1,100°C, γ′ precipitates coarsen rapidly → lose strength. Maximum use T for CMSX-4 is ~1,060°C.

**ODS concept:** Disperse insoluble oxide particles (Y₂O₃, Al₂O₃, La₂O₃, ~5–50 nm) in metal matrix → obstacles to dislocation motion at all temperatures → strength maintained up to 1,300–1,400°C.

**Why oxides don't coarsen:**
- Oxides are thermodynamically stable at high T (Ellingham: very negative ΔG_f)
- Insoluble in metal matrix → no Ostwald ripening
- Pinning force: F = γ_gb × π × r (Zener pinning, Ch 09)

**ODS alloy production (ONLY possible by PM):**
```
1. Mechanical Alloying (MA):
   Metal powder + Y₂O₃ powder → Ball mill (high energy) → 20+ hours
   → MA powder: Y₂O₃ uniformly broken down to 2–5 nm + dispersed in matrix
   
2. Degassing: heat powder in vacuum (to remove H₂O, CO)
   
3. Can and extrude: HIP or hot extrude to consolidate
   
4. Anneal: elongated grains from extrusion + high-T anneal
   → COARSE elongated grains (grain boundaries parallel to stress)
   
5. Finish machine
```

**Commercial ODS alloys:**
| Alloy | Base | Y₂O₃% | Max T (°C) | Applications |
|-------|------|-------|-----------|-------------|
| MA754 | Ni-20Cr | 0.6 | 1,200 | Vane airfoils, combustors |
| MA758 | Ni-30Cr | 0.6 | 1,200 | Hot gas parts |
| MA956 | Fe-20Cr-4Al | 0.5 | 1,350 | Combustor liners |
| PM2000 | Fe-20Cr-5Al | 0.5 | 1,350 | Gas turbine combustors |
| ODS-EUROFER | Fe-9Cr | 0.3 | 700 | Fusion reactor blanket |

**Limitation:** Anisotropic properties; difficult to join (welding destroys the elongated grain structure near the weld); high cost of MA processing.

---

## Summary

| PM Process | Density | Application | Advantage |
|-----------|---------|-------------|-----------|
| Die compact + sinter | 85–95% | Iron/steel parts (gears, bearings) | High volume, low cost |
| CIP + sinter | 90–97% | WC-Co tools, Ti preforms | Uniform density |
| HIP | ~100% | PM superalloy disks, cemented carbides | Full density + fine grain |
| Mechanical alloying | ~100% after HIP | ODS alloys | Only route for nm oxide dispersion |
| Metal injection molding | 96–99% | Complex small parts | Injection molding geometry |

---

## Exercises

1. Nickel superalloy powder (René 88DT) is gas-atomized at a rate of 100 kg/hour. (a) The powder size distribution is log-normal: D10 = 50 μm, D50 = 120 μm, D90 = 250 μm. What fraction of the powder lies in the 50–250 μm range? (b) For PM disk manufacturing, only the 50–150 μm fraction is used. The rest is recycled or scrapped. If 60% of the total powder is in the 50–150 μm range, calculate usable powder yield per day (assuming 16 hours production). (c) Why is it critical to handle this powder under inert atmosphere? What oxide forms if oxygen is present? What property degradation occurs?

2. WC-Co cutting insert: 90 wt% WC + 10 wt% Co. After sintering at 1,350°C for 1 hour: (a) Show that sintering T (1,350°C) is below Co melting point (1,495°C) but above the WC-Co eutectic (~1,320°C). Why does liquid phase form below T_m(Co)? (b) The sintered density = 14.35 g/cm³; theoretical maximum = 14.8 g/cm³. Calculate % porosity. (c) HIPping at 1,250°C / 100 MPa reduces porosity to < 0.1%. Calculate the pressure required to close a pore of radius 10 μm using P = 2γ/r (capillary pressure). If γ_WC-Co interface ≈ 1 J/m², is the HIP pressure (100 MPa) >> capillary back-pressure?

3. ODS MA956 alloy (Fe-20Cr-4Al + 0.5% Y₂O₃) has oxide particles 10 nm in diameter at a spacing of 100 nm. Using Orowan bypass mechanism: τ_bypass = Gb/(2πλ) where G = 80 GPa, b = 0.25 nm, λ = particle spacing = 100 nm. (a) Calculate τ_bypass in MPa. (b) At 1,200°C, G reduces to 50 GPa. Recalculate τ. (c) Compare to CMSX-4 superalloy at 1,200°C where γ′ strengthening has essentially disappeared. Why does ODS maintain strength while γ′ does not? (d) What is the Zener pinning force on a grain boundary for a grain boundary energy γ_gb = 0.5 J/m², particle radius r = 5 nm, volume fraction f_v = 0.005?

4. A PM turbine disk (René 95) requires: σ_y = 1,000 MPa at 650°C, grain size ASTM 12 (d = 5.6 μm). Compare to a conventionally cast + forged disk of the same alloy (grain size ASTM 6, d = 90 μm): (a) Calculate the Hall-Petch strengthening difference (k_y = 500 MPa·μm^(1/2) for nickel superalloys). (b) Why can't René 95 be made by conventional ingot metallurgy in large diameters? (Hint: consider the segregation coefficient of Al, Ti, Nb in Ni alloys and what happens to a 200 mm diameter ingot on quenching from forging temperature.) (c) The PM disk costs 3× more per kg than the cast/forged equivalent. If the disk weighs 50 kg, PM price = $600/kg, cast/forged = $200/kg. Calculate cost difference. Is the cost premium justified by the mechanical property improvement?

5. Mechanical alloying for ODS: Ni + 1% Y₂O₃ starting powders are ball-milled for 0, 4, 8, 16, and 24 hours. TEM measurements after milling: Y₂O₃ particle size decreases from 500 nm → 200 → 80 → 15 → 8 nm. (a) After HIP at 1,150°C/100 MPa, the 24-hour milled material has grain size 500 nm, while the 0-hour has grain size 20 μm. Explain why the oxide dispersion controls grain size during HIP. Use Zener pinning: d_max = (4/3)r/f_v where f_v = 0.007. Calculate predicted maximum grain size for 8 nm particles. (b) Why does longer milling time produce smaller oxide particles? What is the mechanism? (c) The milling process can introduce contamination from mill balls (WC, steel). What effect would WC contamination (density 15.7 g/cm³) have in a Ni-based ODS alloy (density ~8.9 g/cm³)?
