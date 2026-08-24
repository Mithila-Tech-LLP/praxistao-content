# Chapter 65: The Future of Metallurgy — The 2000°C Challenge and Beyond

> **"We have been building on the same metallurgical foundations — melting, casting, forging, heat-treating — for 3,000 years. The jet engine turbine blade required us to push every one of these processes to the absolute limit. For the next generation of aircraft, hypersonic vehicles, fusion reactors, and space systems, we face temperatures and environments where all classical approaches fail. The future of metallurgy is the story of new paradigms, not incremental optimization."**

---

## Table of Contents

1. [The Temperature Frontier — Where We Are and Where We Need to Go](#1-the-temperature-frontier--where-we-are-and-where-we-need-to-go)
2. [Refractory High-Entropy Alloys — The Leading Metallic Candidate](#2-refractory-high-entropy-alloys--the-leading-metallic-candidate)
3. [Oxide-Dispersion-Strengthened Alloys — The Proven Approach](#3-oxide-dispersion-strengthened-alloys--the-proven-approach)
4. [MAX Phase Materials — A Bridge Between Metals and Ceramics](#4-max-phase-materials--a-bridge-between-metals-and-ceramics)
5. [Additive Manufacturing — Geometry Previously Impossible](#5-additive-manufacturing--geometry-previously-impossible)
6. [CALPHAD and AI-Driven Alloy Design](#6-calphad-and-ai-driven-alloy-design)
7. [Hypersonic Materials — Mach 5+ Challenges](#7-hypersonic-materials--mach-5-challenges)
8. [Space and Nuclear Applications](#8-space-and-nuclear-applications)
9. [Sustainable Metallurgy — The Green Transition](#9-sustainable-metallurgy--the-green-transition)
10. [The Research Landscape — Who Is Working on What](#10-the-research-landscape--who-is-working-on-what)
11. [Summary — The Core Challenges of 21st Century Metallurgy](#summary--the-core-challenges-of-21st-century-metallurgy)
12. [Exercises](#exercises)

---

## 1. The Temperature Frontier — Where We Are and Where We Need to Go

**Current state:** 
- Best Ni SX: ~1,100°C metal temperature (TMS-238)
- Best CMC (SiC/SiC): ~1,315°C with EBC

**What the future requires:**

| Application | Required T | Current Best | Gap |
|-------------|-----------|--------------|-----|
| Advanced military jet engines | 1,200°C+ | 1,100°C | 100°C |
| Hypersonic scramjet combustor | 1,600–2,000°C | CMC ~1,315°C | 285–685°C |
| Rocket nozzle (uncooled region) | 2,000–3,000°C | Rhenium (3,186°C MP) | Marginal |
| Fusion reactor first wall | 800–1,200°C + neutron irradiation | Ni alloys (with limits) | Neutron damage |
| Re-entry vehicle leading edge | 1,700–2,200°C | C/C composites | Oxidation |

The gap between "where we are" and "where we need to be" widens enormously for hypersonic and space applications.

**The fundamental physical limit:**

For any material operating at temperature T, the homologous temperature T/T_melt determines whether creep occurs. For T/T_melt < 0.4, creep is negligible for most materials. Therefore:

To operate at 2,000°C without creep → need T_melt > 5,000 K (impossible for metals; diamond melts at 3,550°C under pressure) OR accept creep and design against it (composites, coatings, cooling) OR use materials where T/T_melt is interpreted differently (oxides, carbides, nitrides).

---

## 2. Refractory High-Entropy Alloys — The Leading Metallic Candidate

As described in Chapter 61, RHEAs (based on W, Mo, Nb, Ta, Hf, Ti, Zr, V, Cr combinations) show yield strengths of 400–700 MPa at 1,200–1,600°C.

**What needs to be solved for RHEA to reach turbine blades:**

### 1. Oxidation — The Existential Challenge

Current RHEAs oxidize catastrophically above 800–900°C. The path forward involves RHEA + multilayer coating:

**Option A — Self-forming alumina scale:**
Add Al to RHEA at sufficient concentration → Al₂O₃ scale. But Al reduces melting point and strength. Currently: Al-containing RHEAs like (AlMo₀.₅NbTa₀.₅TiZr) form some Al₂O₃ but not consistently protective.

**Option B — External overlay EBC/TBC:**
Apply a Re/Ir metallic bond coat (stable to 2,200°C) + HfO₂-based or hafnon TBC. Research-stage concept; no production-ready system exists.

**Option C — Combined RHEA + CMC:** Use RHEA as a reinforcement in CMC matrix (RHEA fibers in SiC matrix) — the CMC provides oxidation protection via its EBC system, while RHEA reinforcement provides high-T strength. Most ambitious option; research-stage only.

### 2. Brittleness (DBTT)

The DBTT challenge for Mo/W-rich RHEAs at room temperature. Research directions:
- Al addition reduces DBTT (forms B2 ordered phase → ductile)
- TiZrHfNbTa alloys show DBTT near room temperature → promising candidates for ductile + high-T applications

### 3. Manufacturing at Scale

Current: small-scale arc melting, no castability. Needed: vacuum investment casting (requires inert crucible that survives 3,000°C melt). Mo/W melts attack all current ceramic crucibles. 

Active research: **cold-wall induction melting** (alloy is melted in its own skull — solid alloy around a liquid center — no crucible contact) → potential path to casting complex shapes.

**Timeline:** RHEAs in engine test articles: 2030–2035 (optimistic). Commercial engine use: 2040+.

---

## 3. Oxide-Dispersion-Strengthened Alloys — The Proven Approach

**ODS (Oxide Dispersion Strengthened) alloys** are an existing technology that extends metallic alloy capability:

**Concept:** Disperse a very fine, thermally stable oxide (Y₂O₃, Al₂O₃) as nano-particles (5–50 nm diameter) throughout a metallic matrix. These oxide particles:
- Do not dissolve at service temperatures (stable to >1,400°C for Y₂O₃)
- Pin dislocations against climb → dramatically reduce dislocation creep rate
- Survive thousands of hours at temperature (unlike γ′ which coarsens)

### MA956 and PM2000 — Commercial ODS Alloys

**MA956:** Fe-20Cr-4.5Al-0.5Ti-0.5Y₂O₃ (via mechanical alloying)
- Y₂O₃ dispersion: ~25 nm particles, 3×10²³ m⁻³ number density
- Excellent oxidation resistance (Al₂O₃ former)
- Creep resistance to 1,200°C
- Used in: combustor liners (rolls-Royce, Pratt & Whitney)

**Strength at 1,200°C (MA956 vs IN738 Ni superalloy):**
```
MA956 ODS:   ~100 MPa yield at 1,200°C (maintained for 1,000h)
IN738 Ni:    ~20 MPa (already near failure at 1,200°C)
```

5× stronger than conventional Ni at 1,200°C — but MA956 density = 7.25 g/cm³ (similar to Ni alloys), so no density benefit.

### ODS Ni Superalloys — The Premium Option

**MA6000:** Ni-15Cr-4.5Al-2Ti-2Mo-4W-2Ta-1.1Y₂O₃
- Combines γ′ precipitation strengthening (from Al, Ti, Ta) WITH oxide dispersion strengthening
- The Y₂O₃ particles pin dislocations even as γ′ coarsens at very high T
- Creep resistance to 1,150°C (vs. ~1,100°C for CMSX-4 with same composition base)

**Why ODS isn't universal:**

Mechanical alloying (the process needed to introduce Y₂O₃ uniformly):
- Cannot produce single crystals → grain boundaries remain → limits high-T performance
- Very expensive processing (ball milling + hot extrusion + rolling/HIP)
- Cannot form complex blade shapes easily

ODS is most appropriate for: combustor liners, vanes (stationary), industrial GT components — where shape complexity is lower and single crystal is not required.

---

## 4. MAX Phase Materials — A Bridge Between Metals and Ceramics

**MAX phases** are a class of ternary carbides and nitrides with the formula Mn+1AXn where:
- M = early transition metal (Ti, V, Cr, Nb, Ta, Mo...)
- A = group IIIA or IVA element (Al, Si, Ge, In, Sn...)
- X = C or N

**Examples:** Ti₃AlC₂, Ti₂AlC, Cr₂AlC, Ti₃SiC₂

**Why MAX phases are special — they combine metal and ceramic properties:**

| Property | Typical Metal | Typical Ceramic | MAX Phase |
|----------|---------------|-----------------|-----------|
| Electrical conductivity | High | Low | High (metallic) |
| Thermal conductivity | High | Low | High |
| Machinability | Good | Very poor | Good (like soft metal) |
| Oxidation resistance | Moderate | Excellent | Excellent |
| Fracture toughness | High | Very low | Moderate (5–15 MPa√m) |
| Density | High (Fe: 7.9, Ni: 8.9) | Low (SiC: 3.2) | Medium (Ti₃AlC₂: 4.2) |
| Ductility | High | None | Some (non-linear deformation) |

**Mechanism of unusual properties:**

MAX phases have a layered crystal structure alternating between MX layers (responsible for strength, conductivity) and A layers (responsible for ductility, kink resistance). Deformation occurs by kinking of the layers — a non-classical mechanism unique to MAX phases:

```
Undeformed:                  Kinked:
  ─────────────────────       ─────────────────────
  M₃Al layer (soft A)    →    ╲╲╲╲╲╲╲╲╲╲╲╲╲╲╲╲╲╲╲╲
  Ti₃C₂ layer (hard M)        ──────────────────────
  M₃Al layer               ╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱
  Ti₃C₂ layer                 ──────────────────────
```

Kink bands nucleate and propagate instead of cracks → energy absorption without complete fracture → "pseudoplastic" behavior.

### MAX Phase Applications

**Cr₂AlC** is particularly interesting for turbine applications:
- Forms protective α-Al₂O₃ scale above 900°C
- Good oxidation resistance to 1,300°C
- Density: 5.2 g/cm³ (40% lighter than Ni superalloys)
- Elastic modulus: 200–280 GPa (similar to metals)

**Research applications:**
- MAX phase coating (Cr₂AlC) on Ni superalloy blades — oxidation protection
- MAX phase insert in CMC leading edges — combined conductivity + oxidation resistance
- MAX phase self-healing: at high temperature, MAX phases can repair surface cracks by oxidation products filling the crack → "self-healing" oxides

Current TRL for turbine applications: 3–4 (laboratory demonstrations). Commercial use: 2030+ (coatings), 2035+ (structural applications).

---

## 5. Additive Manufacturing — Geometry Previously Impossible

**Additive manufacturing (AM)** builds parts layer-by-layer from powder or wire, enabling geometries impossible with conventional casting, forging, or machining.

### AM Processes for Metals

**Selective Laser Melting (SLM) / Laser Powder Bed Fusion (LPBF):**
- Metal powder spread in thin layers (~30–50 μm)
- Laser melts powder selectively → solidifies → next layer added
- Residual stresses are high (rapid cooling)
- Best surface finish; best feature resolution
- Limited build volume

**Electron Beam Melting (EBM):**
- Electron beam melts powder in vacuum
- Higher temperature environment → lower residual stresses vs. SLM
- Somewhat rougher surface
- Used for: titanium implants, complex turbine components (GE Additive)

**Directed Energy Deposition (DED) / LENS:**
- Powder or wire fed into a laser melt pool on an existing part
- Can add material to an existing component
- Used for: repair (blade tip, platform rebuild), hybrid manufacturing

### What AM Enables for Turbine Blades

**Conformal cooling holes:**
Traditional drill/EDM: only straight-line holes possible.
AM: hollow lattice structures, any geometry cooling channels, curved passages, **hollow struts** inside the blade interior → heat transfer surface area can increase 10× → dramatically better cooling with same cooling air mass flow.

**Lattice-core blades:**
Replace solid internal wall with optimized lattice structures:
- 30–50% weight reduction in internal structure
- Each lattice member sized for local stress → material where needed, gone where not needed
- Heat transfer optimized per location → less overtemperature everywhere

**Current AM blade status:**
- GE Aviation: fuel nozzles in LEAP engine (AM-made, 25% weight reduction, 5× durability)
- Siemens: AM turbine vanes in industrial gas turbines (SLM Inconel 718)
- HPT SX blades via AM: not yet (single crystal requires controlled solidification, not layer-by-layer)

**Can we AM a single crystal?**

Current research: directional solidification during LPBF by controlling scan pattern and heat extraction → columns with preferred orientation → "quasi-SX" microstructure.
- Achieved: columnar grains in [001] direction with ~10° scatter
- Challenge: secondary orientation uncontrolled; many low-angle boundaries

Prediction: 2030–2035 for first AM SX-like turbine blades (with appropriate post-process heat treatment).

---

## 6. CALPHAD and AI-Driven Alloy Design

The future of alloy development is computational-first:

### CALPHAD 2.0 — High-Throughput Screening

**CALPHAD** (Chapter 06) can calculate phase diagrams for any composition. Modern implementation:
- Automated calculations for 10⁶–10⁹ compositions per day
- Predict: phase stability, γ′ volume fraction, solidus, density, TCP risk
- Cross-reference with target property requirements → filter to candidate alloys

**Current databases:**
- Thermocalc TCNI9 (Ni alloys): covers ~30 elements
- Thermocalc TCREFRACTORY (refractory alloys): covers W, Mo, Nb, Ta, Re, V, Cr, Hf combinations

### Machine Learning for Property Prediction

**ML models trained on experimental databases** can predict:
- Yield strength, UTS, creep life from composition alone
- Oxidation rate from composition
- Castability index from composition

**Active databases:**
- AFLOW-LIB: DFT-calculated properties for ~3.5 million compounds
- NOMAD: experimental materials data repository
- ICSD: crystal structure database

**Example:** SIPFENN (DFT-informed neural network) can predict defect formation energies in multi-component alloys — relevant for predicting TCP nucleation tendency in RHEAs.

### The ICME Framework

**Integrated Computational Materials Engineering (ICME):** Design materials + processes + structures simultaneously using linked computational models:

```
Composition → CALPHAD → Phase stability
     ↓                        ↓
  DFT → Defect energetics → Dislocation mobility
     ↓                        ↓
  Phase field → Microstructure evolution
     ↓                        ↓
  Crystal plasticity → Local mechanical response
     ↓                        ↓
  FEM structural analysis → Component life prediction
```

Each model's output feeds the next → from atoms to engine service life in one computational chain.

**Current ICME capability:**
- Predict new alloy creep life within factor of 2 before synthesis (vs. trial-and-error)
- Reduce alloy development time from 10–20 years to 3–5 years
- Already applied: GE's LEAP engine alloy optimization, DARPA Accelerated Metallurgy program

---

## 7. Hypersonic Materials — Mach 5+ Challenges

**Hypersonic vehicles** (Mach 5+) face unique materials challenges:
- Aerodynamic heating: stagnation point temperature can reach 1,700–2,500°C (depending on speed and altitude)
- Duration: 5–30 minutes (HGV — Hypersonic Glide Vehicles) to hours (hypersonic cruise missiles)
- Oxidizing atmosphere (high-T air with dissociated oxygen and nitrogen)
- Structural loads + acoustic fatigue + vibration

### Ultra-High Temperature Ceramics (UHTC)

**Hafnium diboride (HfB₂)** and **zirconium diboride (ZrB₂)** composites:
- Melting point: HfB₂ 3,380°C, ZrB₂ 3,245°C
- Density: HfB₂ 11.2 g/cm³, ZrB₂ 6.1 g/cm³
- Excellent oxidation resistance (forms HfO₂/ZrO₂ + B₂O₃ glass scale at very high T)
- Good thermal conductivity (allows some heat flow away from leading edge)
- Very brittle → K_Ic = 2–4 MPa√m

**ZrB₂-SiC composite:** Adding 15–20% SiC to ZrB₂:
- Improves oxidation resistance (SiO₂ glass fills B₂O₃ evaporation pores)
- Increases fracture toughness slightly (SiC crack bridging)
- Used in: hypersonic leading edge test articles, re-entry vehicle nose caps (research level)

### Carbon-Carbon Composites (C/C)

**C/C = carbon fiber in carbon matrix:**
- Density: 1.7–2.0 g/cm³ (lightest high-T material)
- Maintains strength to 2,000°C+ in inert/reducing environments
- NO strength loss at high temperature (carbon strengthens on heating up to ~2,500°C)
- **Critical limitation:** Oxidizes above 400°C (carbon burns to CO₂)

**C/C-SiC:** Carbon fiber with part-SiC matrix → oxidation protection to ~1,600°C
Used in: brake discs (Formula 1, aircraft brakes — C/C only), thermal protection system of Space Shuttle (C/C-SiC nose and wing leading edges).

---

## 8. Space and Nuclear Applications

### Space Propulsion

**Nuclear thermal propulsion (NTP):**
Rocket that heats H₂ propellant by nuclear reactor → specific impulse 2× conventional chemical rocket → enables faster deep-space missions.

Material challenge: Fuel rods must survive 2,700°C in H₂ atmosphere for 1–10 hours.

**Leading material:** Tungsten-rhenium (W-Re) alloys:
- W-26Re: high T capability, not embrittled by neutron irradiation (unlike pure W)
- Carbide-coated W: UC (uranium carbide) fuel inside W-Re cladding
- T_max: ~2,700°C (tested in NERVA program, 1960s–1970s)

### Fusion Reactor Materials

**First wall (the plasma-facing component):**
- Plasma temperature: 150,000,000°C (but plasma touches nothing directly — magnetic confinement)
- First wall plasma-facing material operates at ~800–1,200°C
- Subject to 14 MeV neutron bombardment → atomic displacement, transmutation, helium bubble formation → embrittlement

**Current candidate: Tungsten (W)**
- Highest melting point (3,422°C)
- Low neutron activation (relatively)
- Low tritium retention
- **Challenge:** W becomes brittle after neutron irradiation (DBTT increases dramatically)

**Research direction:** W alloys with Re or ODS particles → stabilize against radiation embrittlement. ITER (the international fusion experiment in France) uses W plasma-facing components.

---

## 9. Sustainable Metallurgy — The Green Transition

The materials industry generates ~8% of global CO₂ emissions. The future of metallurgy must address sustainability:

### Low-Carbon Steel Production

**Electric arc furnace (EAF) + green hydrogen DRI:**
Traditional blast furnace: coke (coal) reduces iron ore → ~2 tonnes CO₂ per tonne steel.
New route: Green H₂ (electrolytic) reduces iron ore → direct reduced iron (DRI) → EAF melting.
CO₂ reduction: >90% vs. blast furnace.

**Current status:** H2GreenSteel (Sweden) has announced world's first commercial green steel production, targeting 2025. SSAB's HYBRIT process produced fossil-free steel in 2022.

### Recycling and Urban Mining

Rhenium, tantalum, cobalt (all critical for superalloys) are recovered from:
- Spent turbine blades (Re recycling: 80% of Re in retired blades can be recovered)
- Electronic waste (Co from batteries, Ta from capacitors)
- Mining waste (secondary recovery from Cu/Mo tailings)

**Circular economy for superalloys:** Engine OEMs (GE, RR, P&W) are developing "blade-to-blade" recycling where retired blades' alloy is remelted and recast into new blades → Re and Ru recovery → reduced primary mining requirement.

### Computational Design Reduces Trial-and-Error

Every "failed" alloy in traditional development:
- Consumes raw materials (Re at $3,500/kg, Ru at $15,000/kg)
- Uses energy for melting, heat treatment, testing

CALPHAD/ML-guided design: fewer experimental iterations → less material waste. A development program that used to test 200 alloys now tests 20 → 10× less material consumed in R&D.

---

## 10. The Research Landscape — Who Is Working on What

### Academic Centers

**USA:** MIT (DMSE), Caltech, Carnegie Mellon, Georgia Tech, UC Santa Barbara (the "Superalloy Capital" — originally Tresa Pollock's lab)

**UK:** University of Cambridge (David Dye, Mark Liddell groups), Imperial College London, Oxford

**Japan:** National Institute for Materials Science (NIMS) — world leader in SX superalloy development (TMS series), Kyoto University

**Germany:** DLR, Jülich Research Center

### Industrial R&D

**GE Aviation / GE Aerospace:** AM turbine components, CMC blades, next-gen alloys
**Rolls-Royce:** ODS combustors, advanced TBC, RHEA research (UK government-funded ADVENT program)
**Pratt & Whitney (RTX):** Advanced cooling architectures, Pt-aluminide coatings, CMC vane programs
**MTU Aero Engines:** EU TURBO-D program for CMC and RHEA turbine disk materials
**Safran:** LEAP engine CMC development, advanced forming processes

### National Programs

**USA:** DARPA (Materials Genome Initiative successor programs), DoD Next-generation turbine engine programs
**EU:** Clean Aviation program (successor to Clean Sky 2): targets 50% fuel reduction by 2050 → requires materials at ~1,300°C metal temperature
**China:** AECC (Aero Engine Corporation of China): major investment in domestic superalloy capability (STAL series), CMC development for COMAC aircraft
**India:** GTRE Kaveri engine: indigenous superalloy development

---

## Summary — The Core Challenges of 21st Century Metallurgy

### The Five Grand Challenges

**1. The 2,000°C Challenge**
Operating structural components at 2,000°C — required for next-gen propulsion. Solution path: RHEA + coating, or UHTC/C-C composites, or completely rethinking the cooling architecture. Timeline: 2040–2050.

**2. The Weight Challenge**
Every kg removed from a turbine component saves 5–10 kg in surrounding structure. CMC (density 2.7 g/cm³ vs. Ni at 8.7 g/cm³) is the solution in sight. Timeline: CMC dominance by 2040.

**3. The Design Speed Challenge**
Current alloy development: 10–20 years from concept to service. With ML + CALPHAD + ICME: 3–5 years. This is already happening — the methodology exists, validation time is the bottleneck.

**4. The Sustainability Challenge**
8% of global CO₂ from materials production. Green steel, green aluminum, Re/Co/Ru recycling, computational reduction in experimental waste → all active programs.

**5. The New Paradigm Challenge**
Classical alloying (primary metal + small additions) is mature. HEAs, MAX phases, CMCs, AM — all represent new paradigms. The next transformative material is probably not yet on any roadmap. It will be found in the unexplored 10²⁵ composition space using tools that don't yet exist.

---

### Where It All Started and Where It's Going

We began this course with the definition of metallurgy: the art and science of working with metals, from ore to engineered component.

We end having traced a journey from:
- The Bronze Age (copper + tin = bronze, 9000 BCE) 
- To the Iron Age (carbon in iron = steel, 1200 BCE)
- To the Alloy Age (Ni + 14 elements = CMSX-4, 1985 CE)
- To the Composite Age (SiC fibers + SiC matrix + EBC = GE9X blade, 2019 CE)
- And toward the Computational Materials Age (AI + CALPHAD + ICME = designed-on-screen, tested-in-simulation, 2030+)

The 2030°C challenge — building a material that can sustain structural loads at temperatures most metals cannot even approach — is the Mt. Everest of materials science. We've climbed ~2/3 of the mountain.

The tools to climb the rest — high-entropy alloys, ceramic composites, additive manufacturing, computational design — are in our hands.

---

## Exercises

1. Temperature scaling: Ni has melting point 1,455°C. CMSX-4 operates at T_H = T_metal/T_melt = 0.75 (75% of melting point). An RHEA with T_melt = 2,800°C operating at the same T_H = 0.75 would have what metal temperature capability? How does this compare to current Ni SX capability, and what would this mean for engine efficiency (1% efficiency per 25°C improvement)?

2. Additive manufacturing lattice blade: A conventional solid-wall HPT blade has a specific surface area (internal surface / volume) of 2,000 m⁻¹ for cooling. A lattice-core blade has specific surface area 20,000 m⁻¹. If heat transfer Q = h × A × ΔT and both have the same volume and same h, by what factor can the ΔT (temperature difference between blade wall and cooling air) be reduced with the lattice? What does this mean for blade metal temperature?

3. CALPHAD screening rate: A modern CALPHAD software can evaluate one alloy composition in 0.1 seconds. A 5-element RHEA system with elements W, Mo, Nb, Ta, Ti varied in 5 at% increments (0, 5, 10, 15, 20...45 at% for each, summing to 100%). (a) Estimate the number of possible compositions (hint: this is a stars-and-bars combinatorics problem: compositions of 5 parts summing to 100 in steps of 5). (b) How long would it take CALPHAD to screen all of them? (c) If only 0.1% pass the first filter (TCP stable, T_melt > 2,200°C), how many remain for experimental testing?

4. Sustainability calculation: An airline retires 1,000 CMSX-4 HPT blades per year (from 200 engines). Each blade contains 3 wt% Re and weighs 0.5 kg. (a) How many kg of Re are in the retired blades? (b) Global annual Re production is 60 tonnes. What fraction of world production does this represent? (c) If recycling captures 80% of Re from retired blades, and the same fraction of the global fleet does this, estimate the total Re recovery if there are 30,000 commercial aircraft with equivalent blade usage.

5. The 2,000°C challenge timeline: Research programs currently suggest: (1) RHEA oxidation-resistant alloys might be ready for ground engine tests by 2032, (2) CMC capability might reach 1,400°C by 2030, (3) UHTC-composite leading edges might be tested in scramjet engines by 2028. Write a 500-word scenario describing what a 2045 hypersonic commercial aircraft might look like — what materials would be used where, what temperature each section operates at, and what enabling technologies (not yet commercialized) would be required. Consider: powerplant (scramjet), airframe (high heat flux leading edges), thermal protection system, structural fuselage.
