# Chapter 53: The CMOS Process

Understanding how a chip is fabricated reveals why modern transistors are so precise, so small, and so expensive to develop. The CMOS manufacturing process is a sequence of over 1000 individual steps — depositing thin films, selectively removing material, implanting atoms, and growing oxides — repeated many times on a silicon wafer. Each step must be controlled to atomic precision. This chapter walks through the major process steps, explains the key equipment and materials, and connects the physical process to the transistor structures described in Chapter 52. By the end, you will understand how a logic gate goes from silicon sand to a working circuit.

## Table of Contents

1. [The Silicon Wafer](#1-the-silicon-wafer)
2. [Key Process Steps: Deposit, Pattern, Etch, Implant](#2-key-process-steps-deposit-pattern-etch-implant)
3. [Building a FinFET: The Modern Transistor](#3-building-a-finfet-the-modern-transistor)
4. [Interconnects: Metal Layers and Vias](#4-interconnects-metal-layers-and-vias)
5. [Planarization and CMP](#5-planarization-and-cmp)
6. [Front-End-of-Line vs Back-End-of-Line](#6-front-end-of-line-vs-back-end-of-line)
7. [Summary](#summary)
8. [Exercises](#exercises)

---

## 1. The Silicon Wafer

Everything starts with a **silicon wafer** — a thin, circular disc of extremely pure silicon crystal.

**From sand to wafer:**
1. Quartzite (SiO₂ sand) is reduced to metallurgical-grade silicon (98% pure) with carbon in an electric arc furnace
2. Metallurgical silicon is refined to polysilicon (>99.9999999% pure = 9N purity) via the Siemens process
3. The Czochralski process grows a single-crystal silicon ingot: a seed crystal is lowered into molten silicon (~1415°C) and slowly pulled out while rotating, forming a cylindrical ingot up to 300mm diameter × 1.5m long
4. The ingot is sliced into wafers (~0.775mm thick), lapped, polished to mirror finish

```
Wafer sizes (larger = more chips per wafer = lower cost per chip):
  1960s: 1 inch (25mm)
  1980s: 4 inch (100mm)
  1990s: 6 inch (150mm)
  2000s: 8 inch (200mm) — still used for analog, MEMS, power chips
  2010s: 12 inch (300mm) — current standard for leading edge
  Next:  18 inch (450mm) — proposed but not commercially adopted
  
  300mm wafer area = π×150² = 70,685 mm²
  5nm chip die ~80mm²: ~883 chips per wafer (ignoring edge loss)
```

**Wafer crystal orientation**: The crystal is cut along specific lattice planes (⟨100⟩ or ⟨110⟩ orientation) that affect transistor characteristics and cleaving direction.

**Wafer resistivity**: The base p-type doping level is precisely controlled (typically 1–10 Ω·cm) as it affects well implant and isolation characteristics.

### Quick Check
> 1. What is the Czochralski process and why does it produce a single-crystal ingot?
> 2. Why did the industry move from 6-inch to 8-inch to 12-inch wafers?
> 3. What purity level does semiconductor-grade silicon require?

---

## 2. Key Process Steps: Deposit, Pattern, Etch, Implant

The CMOS process repeats four fundamental operations many times:

### A. Thin Film Deposition

**Thermal oxidation**: Expose silicon wafer to O₂ or H₂O at 900–1200°C. Silicon reacts to form SiO₂ (silicon dioxide = glass) right on the surface. Used for gate oxide and field oxide.

**CVD (Chemical Vapor Deposition)**: Reactive gases flow over the heated wafer; they react and deposit a thin film on the surface. Used for polysilicon, silicon nitride (Si₃N₄), various dielectrics.

**PVD (Physical Vapor Deposition) / Sputtering**: A metal target is bombarded with argon ions; sputtered metal atoms deposit on the wafer. Used for metal layers (tungsten, titanium, tantalum).

**ALD (Atomic Layer Deposition)**: Deposit exactly one atomic layer per cycle (alternating precursor gases). Critical for gate dielectrics and barrier layers in advanced nodes. Achieves sub-1nm control.

### B. Photolithography (Patterning)

The key step that defines the circuit features. (Chapter 54 covers this in depth.) Short version:
1. Coat wafer with photoresist (light-sensitive polymer)
2. Expose with UV/EUV light through a patterned mask (reticle)
3. Develop: exposed resist dissolves → pattern on wafer

### C. Etching

**Wet etching**: Dip wafer in chemical bath. Isotropic (etches in all directions equally). Simple but undercutting causes feature rounding — limited to less precise steps.

**Dry etching / RIE (Reactive Ion Etching)**: Plasma of reactive gas ions bombard the wafer surface. Anisotropic (etches primarily vertically, leaving sidewalls intact). Required for sub-100nm features.

```
Wet etch (isotropic):           Dry etch (anisotropic):
  
  ┌──────┐                      ┌──────┐
  │resist│                      │resist│
  ──────────────────            ──────────────────
  │oxide to remove │            │oxide to remove │
  ──────────────────            ──────────────────
  │   substrate    │            │   substrate    │
  
  After wet etch:               After dry etch:
  ┌────┐ ┌────┐                 ┌──────┐
  └────┘ └────┘                 └──────┘
   undercutting                  vertical walls
```

### D. Ion Implantation

Dopant atoms (boron, phosphorus, arsenic) are ionized, accelerated to 10keV–1MeV, and fired into the wafer. Implanted depth is controlled by acceleration energy; dose by beam current and time. After implantation, a high-temperature anneal activates the dopants and repairs crystal damage.

```
Ion implanter:
  Dopant gas source → Ion beam extraction → Mass analyzer 
    → Acceleration → Deflection (to scan beam) → Wafer
  
  Typical: 1 second of implant → 10¹² to 10¹⁶ dopant atoms/cm²
  Energy range: 1 keV (shallow surface) to 1 MeV (deep well implants)
```

### Quick Check
> 1. What is the difference between wet etching and dry etching, and why is dry etching required for sub-100nm features?
> 2. What is ALD (Atomic Layer Deposition) and why is it needed for gate dielectrics?
> 3. What does ion implantation do and how is implant depth controlled?

---

## 3. Building a FinFET: The Modern Transistor

The planar MOSFET (Chapter 52) worked until about 22nm. Below that, the gate could no longer control the channel — short-channel effects caused excessive leakage. Intel introduced the **FinFET** (Fin Field-Effect Transistor) in 2011 (at 22nm) — a 3D transistor where the channel is a tall, thin "fin" of silicon that the gate wraps around from three sides.

```
Planar MOSFET vs FinFET:

Planar (top view + side):       FinFET (3D):
  
  Gate ─────────────               Gate
  ─────────────────              ┌──┐│┌──┐
  ─────────────────              │  ├┤┤  │
  S         D                   S  ├┴┤  D
  
  Gate controls channel          Gate wraps three sides of fin:
  from ONE side (top)            - Top
                                 - Left side
                                 - Right side
                                 
  Better electrostatic control → can turn OFF at smaller gate lengths
```

**FinFET fabrication (simplified):**

```
Step 1: Define fins
  - Grow hard mask (SiN) on p-Si substrate
  - Pattern fins using extreme UV lithography + spacer-based patterning
  - RIE etch to create silicon fins (height ~50nm, width ~6nm)
  
Step 2: Form well regions
  - Ion implant n-well and p-well dopants for NMOS/PMOS regions
  - Anneal to activate dopants

Step 3: Shallow Trench Isolation (STI)
  - Deposit SiO₂ to fill between fins
  - CMP to planarize
  - Recess oxide to expose fin sidewalls and top

Step 4: Dummy gate
  - Deposit polysilicon dummy gate
  - Pattern gate using lithography
  - Implant source/drain

Step 5: Replace dummy gate (replacement metal gate - RMG)
  - Remove polysilicon dummy gate
  - Deposit high-K dielectric (HfO₂, ~2nm equivalent oxide thickness)
  - Deposit metal gate (TiN work function metal + W fill)

Step 6: Source/drain
  - Epitaxial growth of SiGe (for pMOS) or Si:P (for nMOS) in source/drain
  - This strains the channel → higher carrier mobility
```

**Gate-All-Around (GAA)** (Samsung 3nm, Intel 20A/18A, TSMC N2): The next evolution after FinFET. The channel is a horizontal nanosheet (or nanowire) with the gate wrapping all four sides (top, bottom, left, right). Even better electrostatic control, enabling smaller gate lengths.

```
FinFET vs GAA nanosheet:

FinFET:              GAA:
  Gate                Gate (surrounds)
  ├──┤               ┌─────────────┐
  │fin│              │  nanosheet  │
  └──┘               └─────────────┘
  3 sides            4 sides
```

### Quick Check
> 1. Why did the planar MOSFET fail below 22nm and what does FinFET solve?
> 2. What is the "fin" in FinFET and how many sides does the gate wrap?
> 3. What is a Gate-All-Around (GAA) transistor and when was it introduced?

---

## 4. Interconnects: Metal Layers and Vias

Transistors are useless without connections. Modern chips have 10–15+ layers of metal interconnect, each separated by insulating dielectric, connected by vertical "vias."

```
Metal layer stack (cross-section):
  
  M10  ═══════════  ─═════════ (global power, clock distribution)
        │         │    │
  V9    ↕         ↕    ↕       (via layer)
  M9   ─══════════  ═══════    (semi-global routing)
        │              │
  V8    ↕              ↕
  M8   ═════  ─════════  ══    
  ...   (more layers)
  M3   ─  ─  ─  ─  ─  ─  ─    (local routing)
  V2   ↕  ↕  ↕  ↕  ↕  ↕
  M2   ─────────────────────
  V1   ↕  ↕  ↕
  M1   ────────────────────    (local connections between transistors)
  ─────────────────────────
  ◘ ◘ ◘ ◘ ◘ ◘ transistors ◘ ◘
```

**Copper interconnects**: Since 1997, copper (Cu) has replaced aluminum. Cu has ~40% lower resistivity, enabling faster signals and lower power. But copper diffuses into silicon (poisoning transistors) — so copper must be surrounded by a **barrier layer** (TaN/Ta) deposited by ALD.

**Low-K dielectrics**: The capacitance between adjacent metal wires causes RC delays and power consumption. Replacing SiO₂ (K=3.9) with porous low-K dielectrics (K=2.5–3.0) reduces capacitance by 30–40%.

**Interconnect resistance and RC delay**: As wires become thinner (sub-10nm), resistivity increases dramatically due to electron scattering at surfaces and grain boundaries. At 5nm pitch, copper wire resistance is 3–5× higher than bulk copper. This is now the dominant performance limiter, not transistor speed.

**Ruthenium (Ru) interconnects**: Recent advanced nodes use ruthenium for the tightest (local) metal layers — Ru maintains lower resistance at small dimensions than Cu.

### Quick Check
> 1. Why does a modern chip need 10–15 metal layers instead of just one?
> 2. Why was aluminum replaced by copper for interconnects in 1997?
> 3. What is a "via" in chip interconnects?

---

## 5. Planarization and CMP

After each deposition and etch cycle, the wafer surface becomes uneven — raised and lowered regions from deposited films and etch steps. Subsequent lithography requires a nearly flat surface for accurate focus. **CMP (Chemical Mechanical Planarization)** solves this.

**CMP process:**
- Wafer is pressed against a rotating polishing pad
- A slurry (fine abrasive particles + chemical etchant) is applied
- The mechanical abrasion + chemical dissolution planarizes the surface
- Hard areas polish slowly; soft areas polish faster → global flatness

```
Before CMP:                     After CMP:
  ╱╲╱╲╱╲ uneven surface           ─────────── flat surface
  ▓▓▓▓▓▓ deposited oxide          ▓▓▓▓▓▓▓▓▓▓ uniform oxide
```

CMP is used:
- After STI (Shallow Trench Isolation) oxide deposition
- After inter-layer dielectric (ILD) deposition before each new metal layer
- After tungsten plug (via) deposition
- After copper damascene fill (damascene = trench-fill-and-polish for copper)

CMP uniformity across the 300mm wafer is critical — non-uniform polishing causes some chips to be too thin (damaged circuits) or too thick (poor focus for lithography).

### Quick Check
> 1. What is CMP and why is it needed during chip fabrication?
> 2. What would happen to subsequent lithography steps if CMP were not used?
> 3. What is the "damascene" process for copper interconnects?

---

## 6. Front-End-of-Line vs Back-End-of-Line

The chip fabrication process is divided into two major phases:

**FEOL (Front-End-of-Line)**: All the steps that create the transistors — from the bare silicon wafer through source/drain formation and the gate stack. The result is a wafer covered in transistors, but with no connections between them.

**BEOL (Back-End-of-Line)**: All the steps that create the metal interconnect layers — depositing dielectrics, patterning trenches, filling with copper, polishing, and repeating for each metal layer. The result is a fully wired chip.

```
Process timeline:
  
  Bare wafer
  │
  ├── FEOL (Front-End-Of-Line) ~400–600 steps
  │     Well implants
  │     Shallow trench isolation (STI)
  │     Gate dielectric (high-K)
  │     Gate electrode (metal gate)
  │     Source/drain implant + anneal
  │     Silicide (contacts to transistors)
  │
  ├── BEOL (Back-End-Of-Line) ~400–700 steps
  │     Contact plug (tungsten)
  │     Metal 1 layer
  │     Via 1
  │     Metal 2 layer
  │     ... (repeat for each of 10–15 metal layers)
  │     Final passivation layer
  │
  └── Finished wafer → test → dicing → packaging
```

**FEOL temperature constraint**: Transistors are formed at high temperatures (up to 1050°C for anneals). Metal interconnects cannot withstand this — copper melts at 1085°C and aluminum at 660°C. BEOL is done at low temperatures (<400°C for copper-compatible processes) to prevent damage to existing layers.

**The challenge**: Every FEOL and BEOL process step accumulates. A 5nm chip has 3000+ process steps over 3–4 months of fabrication time. Each step adds defect risk. A single particle of dust or contamination can destroy a chip.

### Quick Check
> 1. What is FEOL and what structures does it create?
> 2. What is BEOL and what is its purpose?
> 3. Why must BEOL processing occur at lower temperatures than FEOL?

---

## Summary

- **Silicon wafer**: grows as a single crystal (Czochralski), 300mm diameter, >9N purity. Larger wafers = more chips = lower cost.
- **Core process steps**: thin film deposition (CVD/ALD), patterning (photolithography), etching (dry/anisotropic for small features), ion implantation (dopants).
- **FinFET**: 3D transistor with gate wrapping three sides of a silicon fin. Solves short-channel effects below 22nm. Replaced by Gate-All-Around (GAA) at 3nm/2nm.
- **Interconnects**: 10–15 copper metal layers separated by low-K dielectric. Vias connect layers. Copper barrier layers (TaN) required to prevent diffusion.
- **CMP**: polishes wafer flat after each deposition, enabling subsequent lithography steps.
- **FEOL vs BEOL**: FEOL creates transistors (high temperature), BEOL creates metal connections (low temperature, <400°C).

---

## Exercises

### Easy
1. What is the Czochralski process and what kind of silicon does it produce?
2. What is the difference between wet etching and dry etching?
3. What is FEOL and BEOL in chip fabrication?

### Medium
4. Defect density and yield: A process has a defect density of D = 0.05 defects/cm². Die area A = 50mm² = 0.5 cm². (a) Using the simple Poisson yield model Y = e^(-D×A), calculate yield. (b) If die area doubles to 100mm²: new yield? (c) If defect density improves from 0.05 to 0.02 defects/cm²: new yield at 100mm²? (d) A 300mm wafer has area 706.8 cm². At D=0.05, A=0.5cm²: how many wafers would you need to get 10,000 good dies?
5. Interconnect RC delay: Metal wire: width 10nm, height 20nm, length 1mm, resistivity 5×10⁻⁸ Ω·m (thin-film Cu, higher than bulk). Adjacent wire capacitance: 0.1 fF/µm. (a) Calculate wire resistance R = ρL/A. (b) Wire capacitance C = 0.1fF/µm × 1000µm. (c) RC time constant. (d) At what clock frequency does the RC delay equal one clock period? (e) Why do global wires (across the whole chip) use wider, taller metal (M10/M11) while local wires (M1/M2) are narrow?
6. FinFET geometry: An Intel 7nm FinFET has: fin height = 52nm, fin width = 6nm, gate length = 14nm, fin pitch = 30nm. (a) Calculate the fin aspect ratio (height/width). (b) The effective gate width = 2×height + width (two sides + top). Calculate W_eff. (c) Compare W_eff to a planar 14nm transistor's width = 14nm. How much more drive current does the FinFET provide? (d) A logic standard cell uses 4 fins. What is the total W_eff and effective width advantage over a single planar transistor?

### Hard
7. Replacement metal gate process: In the "gate-last" or "replacement metal gate" (RMG) process: a dummy polysilicon gate is used during FEOL, then removed and replaced with a high-K/metal gate in a late FEOL step. (a) Why can't you deposit the final HfO₂/TiN gate dielectric before source/drain implant and anneal? (hint: high-K materials are damaged by high temperatures) (b) The dummy gate is removed using an etch selective to silicon nitride spacers. What etch chemistry selectively removes polysilicon without etching SiN? (c) After removing the dummy gate, a 1nm EOT (Equivalent Oxide Thickness) HfO₂ is deposited by ALD. If HfO₂ has K=22 and SiO₂ has K=3.9, what is the physical thickness of HfO₂ that gives 1nm EOT? (hint: EOT = t_HK × K_SiO₂/K_HK) (d) Why does lower EOT improve transistor performance?
8. Process integration challenge: You are process integration engineer at a foundry. An engineer proposes integrating a copper BEOL interconnect layer between two FEOL steps that require annealing at 950°C. (a) Why is this catastrophically wrong? (b) What minimum temperature constraint does copper impose on subsequent steps? (c) What alternative metal could theoretically withstand 950°C? (d) Modern Intel "RibbonFET" (GAA) processes perform anneals after metal gate deposition. How do they manage the temperature constraint? (hint: millisecond laser annealing vs furnace annealing)
