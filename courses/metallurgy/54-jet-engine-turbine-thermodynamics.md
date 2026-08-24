# Chapter 54: The Jet Engine — Thermodynamics and Why Blades Get So Hot

> **"The turbine blade is the most demanding component in engineering. It operates at the intersection of four extreme environments simultaneously: the highest stress, the highest temperature, the most chemically aggressive atmosphere, and the most demanding fatigue loading of any structural component in service. Understanding WHY this is so starts with the thermodynamics of the gas turbine cycle."**

---

## Table of Contents

1. [The Brayton Cycle — Engine Thermodynamics Fundamentals](#1-the-brayton-cycle--engine-thermodynamics-fundamentals)
2. [The Turbofan Engine — Architecture Overview](#2-the-turbofan-engine--architecture-overview)
3. [What Happens at the High-Pressure Turbine](#3-what-happens-at-the-high-pressure-turbine)
4. [How Hot Is the Gas? — Temperature Numbers](#4-how-hot-is-the-gas--temperature-numbers)
5. [Why Higher Temperature Means Better Efficiency](#5-why-higher-temperature-means-better-efficiency)
6. [The Stress Environment — Centrifugal and Bending Loads](#6-the-stress-environment--centrifugal-and-bending-loads)
7. [The Chemical Environment — Oxidation and Hot Corrosion](#7-the-chemical-environment--oxidation-and-hot-corrosion)
8. [The Fatigue Environment — Every Takeoff-Landing Cycle](#8-the-fatigue-environment--every-takeoff-landing-cycle)
9. [The Impossible Requirement — Metal Hotter Than Its Melting Point](#9-the-impossible-requirement--metal-hotter-than-its-melting-point)
10. [How Blade Metal Temperature Is Controlled](#10-how-blade-metal-temperature-is-controlled)
11. [The Complete Loading State — Summary](#11-the-complete-loading-state--summary)
12. [Summary](#summary)
13. [Exercises](#exercises)

---

## 1. The Brayton Cycle — Engine Thermodynamics Fundamentals

A jet engine operates on the **Brayton cycle** (also called the gas turbine cycle). It is a continuous-flow thermodynamic cycle:

### The Four Processes

```
Brayton Cycle (ideal):

Pressure
    |
    |    2────3   ← constant-pressure heat addition (combustion)
    |   /      \
    |  /        \
    | 1          4  ← constant-pressure heat rejection (exhaust)
    └────────────────── Entropy
    
1→2: Isentropic compression (compressor)
2→3: Constant pressure heat addition (combustion)
3→4: Isentropic expansion (turbine)
4→1: Constant pressure heat rejection (exhaust to atmosphere)
```

**Key parameters:**
- **T₁**: Inlet temperature (ambient, ~288 K / 15°C at sea level)
- **T₃**: Turbine inlet temperature (TIT) — the critical design variable
- **Pressure ratio** r_p = P₂/P₁ — typically 30–50:1 in modern engines

### Thermal Efficiency of the Ideal Brayton Cycle

```
η_Brayton = 1 - T₁/T₂ = 1 - 1/(r_p)^((γ-1)/γ)

where:
  r_p = pressure ratio (P₂/P₁)
  γ = ratio of specific heats (~1.4 for air)
```

For a modern engine with r_p = 40:
```
η = 1 - 1/(40)^(0.4/1.4) = 1 - 1/(40)^0.286 = 1 - 1/3.72 = 0.73 = 73% (ideal)
```

Real engines achieve ~55–60% thermal efficiency (largest modern gas turbines) to ~45–50% (aviation turbofans in cruise).

### Why Temperature Matters for Efficiency

The actual work output of the turbine also depends on T₃:
```
Specific work = c_p × (T₃ - T₄) - c_p × (T₂ - T₁)
              = c_p × T₁ × [(T₃/T₁) × (1 - 1/r_p^((γ-1)/γ)) - (r_p^((γ-1)/γ) - 1)]
```

Higher T₃ (turbine inlet temperature) → larger numerator → more work output per unit of compressed air.

**Every 25°C increase in T₃ improves engine efficiency by ~0.5–1%** (through both better efficiency and more work output). Over a long-haul flight with 100 tonnes of fuel, this means hundreds of kilograms of fuel saved.

This is the direct economic and performance incentive to push T₃ as high as possible — and that drives the entire field of superalloy development.

---

## 2. The Turbofan Engine — Architecture Overview

Modern commercial jets use **high-bypass-ratio turbofan** engines. Understanding the layout explains why the HPT blade is so stressed:

```
Large Turbofan Engine (e.g., GE90, CFM LEAP, Rolls-Royce Trent):

    ← Air flow direction →

  Fan    LPC   HPC  Combustor  HPT  LPT   Nozzle
  ╔═══╦═══╦════╦════╦══════╦═══╦══════╦═══╗
  ║ F ║ L ║    ║    ║  🔥  ║HPT║      ║   ║→ thrust
  ║ A ║ P ║ HPC║    ║COMB  ║   ║  LPT ║   ║
  ║ N ║ C ║    ║    ║      ║   ║      ║   ║
  ╚═══╩═══╩════╩════╩══════╩═══╩══════╩═══╝
  
  Bypass duct: ~85% of air bypasses the core → cool thrust
  Core:        ~15% of air goes through compressor→combustor→turbine
  
Key:
  Fan: large front fan, ~3m diameter, low pressure ratio (~1.6)
  LPC: Low-pressure compressor  
  HPC: High-pressure compressor (most compression, very high pressure)
  Combustor: Fuel + air → combustion → hot gas at T₃
  HPT: High-pressure turbine (SUBJECT OF THIS CHAPTER)
  LPT: Low-pressure turbine (drives fan via shaft)
```

### Shaft Arrangement

In a two-spool engine (most modern turbofans):
- **Inner shaft (high-pressure spool)**: HPC + HPT rotate together at 10,000–16,000 RPM
- **Outer shaft (low-pressure spool)**: Fan + LPC + LPT rotate together at 3,000–5,000 RPM

The HPT is directly coupled to the HPC — the HPT must extract exactly the work needed to drive the HPC (which takes 60–70% of the HPT work output). The remaining 30–40% of HPT work goes to driving the aircraft.

---

## 3. What Happens at the High-Pressure Turbine

The HPT is where the hottest, highest-pressure gas from the combustor first contacts metal. It is the most thermally loaded stage:

### Gas Path Through HPT

1. **HPT vanes (nozzle guide vanes, NGVs):** Stationary airfoils that accelerate and redirect combustion gas onto the HPT blade (rotor). NGVs operate at gas temperatures up to TIT (most critical stationary component).

2. **HPT blades (first-stage rotor):** Rotating airfoils that extract energy from the gas. The gas does work on the blade → gas slows and cools → blade spins → shaft turns.

3. **Second-stage HPT vanes and blades:** Continued expansion. Lower temperature than 1st stage.

### Energy Extraction Rate

A GE90-115B (Boeing 777 engine) HPT:
- Air mass flow: ~1,500 kg/s (core flow)
- TIT: ~1,650°C (gas temperature entering HPT, after combustion)
- HPT work output: ~50,000 kW (50 MW) from the 1st stage alone
- HPT rotation: ~10,000 RPM
- Number of 1st stage HPT blades: ~80

Each individual HPT blade extracts:
```
Power per blade = 50,000 kW / 80 blades = 625 kW
= 625,000 watts extracted from the gas per blade
```

This is more power per blade than a residential building consumes.

---

## 4. How Hot Is the Gas? — Temperature Numbers

### Turbine Inlet Temperature (TIT) vs. Metal Temperature

**Turbine Inlet Temperature (TIT)** — the gas temperature entering the HPT — is the engine's primary performance parameter. Modern large turbofans have TIT ≈ 1,500–1,700°C.

But the **blade metal temperature** (Tₘ) is much lower, thanks to cooling:
- Metal temperature at hottest point: 950–1100°C
- TIT: 1,500–1,700°C
- Difference: 400–600°C

This gap (TIT - Tₘ) is maintained by:
1. **Internal cooling**: cool compressor air (~600°C) flows through hollow blade passages → convects heat away from the inside
2. **Film cooling**: cool air exits through surface holes → forms a cool film protecting the blade surface from the hot gas
3. **Thermal barrier coating (TBC)**: ceramic insulating layer on blade surface → reduces heat flux into the metal

Without any cooling, a blade at TIT = 1,600°C would melt (Ni superalloy solidus: ~1340°C). Cooling is existential.

### Temperature Distribution in an HPT Blade

```
HPT Blade Temperature Distribution (schematic cross-section):

Outer surface (gas side):
  ┌──────────────────────────────────────────────────────┐
  │  TBC surface: 1,100–1,200°C (hot gas → TBC → metal)  │
  │  ░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░░  │ TBC layer
  │  Bond coat: 950–1,050°C                               │
  │  Metal (pressure side): 900–1,000°C                   │ Metal wall
  │  Metal (suction side): 850–950°C                      │
  │  Cooling hole / film: 750–850°C                       │
  │  Internal passage wall: 800–900°C                     │
  │  Cooling air: 600–700°C (from compressor)             │
  └──────────────────────────────────────────────────────┘
Inner surface (cooling air side)
```

The temperature gradient through the blade wall (typically 2–3 mm thick) can be 100–200°C → significant thermal stresses.

---

## 5. Why Higher Temperature Means Better Efficiency

Let's quantify this concretely:

**Carnot efficiency** (theoretical maximum for any heat engine):
```
η_Carnot = 1 - T_cold/T_hot
```

For a gas turbine: T_cold ≈ T_ambient (288 K), T_hot = TIT.

At TIT = 1,500°C = 1,773 K:
```
η_Carnot = 1 - 288/1773 = 0.84 (84%)
```

At TIT = 1,700°C = 1,973 K:
```
η_Carnot = 1 - 288/1973 = 0.85 (85%)
```

Real engine efficiency tracks this trend. The difference between 1,500°C and 1,700°C TIT:
- Carnot efficiency: +1% absolute (84% → 85%)
- Real engine efficiency improvement: ~2–4% absolute (because higher TIT also enables higher pressure ratios)
- Fuel burn on a 12-hour flight: -2,000 to -4,000 kg of fuel

For an airline with 300 aircraft flying 10 hours/day at $1/liter jet fuel, a 2% efficiency improvement saves ~$100 million per year.

**This is why airlines and engine manufacturers will pay $50,000 for a better HPT blade.**

---

## 6. The Stress Environment — Centrifugal and Bending Loads

### Centrifugal Stress

The rotating blade experiences centrifugal force pulling it radially outward:
```
Centrifugal force on blade = m × ω² × r_CG

where:
  m = blade mass (~0.5 kg for HPT blade)
  ω = angular velocity (10,000 RPM = 1,047 rad/s)
  r_CG = distance from rotation axis to blade center of mass (~450 mm)
```

```
F_centrifugal = 0.5 × (1047)² × 0.45 = 247,000 N ≈ 247 kN ≈ 25 tonnes force
```

This force is transmitted to the blade root and the disk attachment (fir-tree root). The blade root area is approximately 1,000 mm², so:
```
Blade root stress = F/A = 247,000 N / (1,000 × 10⁻⁶ m²) = 247 MPa
```

At 980°C, this 247 MPa centrifugal stress must be sustained for 30,000 hours without creep failure — that's the fundamental requirement for the blade material.

### Bending/Aerodynamic Loads

The airfoil shape experiences aerodynamic pressure differences between pressure side and suction side. This creates bending moment along the blade length → bending stresses of ~50–100 MPa superimposed on centrifugal load.

### Vibratory Loads (High-Cycle Fatigue)

Gas turbine blades experience resonance vibrations from:
- Wakes from upstream vanes (blade passing frequency)
- Combustor pressure fluctuations
- Engine orders (multiples of rotation speed)

Vibratory stress amplitudes: 50–150 MPa, at frequencies of 1,000–10,000 Hz.

Over 30,000 hours at 10,000 Hz: 10⁹–10¹² cycles → high-cycle fatigue (HCF) is critical.

### Combined Loading Summary

```
                    Primary creep stress:  ~200–250 MPa (centrifugal)
                    Bending stress:        ~50–100 MPa
                    Vibratory stress:      ~50–150 MPa
                    Thermal stress:        ~100–200 MPa (thermal gradient through wall)
                    ─────────────────────────────────────────
                    Peak stress:           ~400–600 MPa (at worst location, worst time)
```

The blade material must handle all of these simultaneously at ~980°C metal temperature.

---

## 7. The Chemical Environment — Oxidation and Hot Corrosion

The HPT blade operates in the combustion gas stream, which contains:

**Oxidizing species:**
- O₂: ~15% of combustion gas (from excess air)
- CO₂, H₂O: oxidation products; H₂O particularly accelerates oxidation

**Hot corrosion species:**
- NaCl: from ingested marine air → reacts with SO₃ → Na₂SO₄ deposits
- SO₂/SO₃: from sulfur in jet fuel (ASTM D1655 limits: < 3,000 ppm S)
- V₂O₅: from vanadium in some fuels (especially industrial gas turbines)

### Two Types of Hot Corrosion

**Type I (High-Temperature) Hot Corrosion:** 800–950°C
Na₂SO₄ deposits (liquid at this temperature) flux the protective Cr₂O₃/Al₂O₃ scales:
```
Na₂SO₄ + Cr₂O₃ → 2NaCrO₂ + SO₃
```
The protective scale dissolves → bare metal exposed → catastrophic oxidation rate.

**Type II (Low-Temperature) Hot Corrosion:** 650–800°C
Na₂SO₄ + NiO → liquid nickel sulfate compound; pitting corrosion.

Modern HPT blades are protected by:
1. **MCrAlY bond coat**: Alumina-forming metallic coating (Chapter 46)
2. **Thermal barrier coating (TBC)**: Ceramic insulator
These coatings extend blade life from ~5,000h to ~30,000h under equivalent conditions.

---

## 8. The Fatigue Environment — Every Takeoff-Landing Cycle

### Low-Cycle Fatigue (LCF)

Every flight is a thermal cycle:
1. **Start**: cold engine (~20°C) → hot operation (blade at 980°C) → ΔT ≈ 960°C
2. **Cruise**: steady state
3. **Throttle changes**: temperature fluctuations
4. **Shutdown**: hot blade cools back to ambient

Temperature differential through the blade wall: ~200°C
This creates cyclic thermal strain → **thermal fatigue**.

Cycles per flight: 1 (major cycle) + several throttle changes (minor cycles)

For a long-service commercial aircraft: 30,000 flights over a 25-year life → 30,000 major LCF cycles + millions of minor cycles.

**LCF life requirement for HPT blade: > 30,000 cycles.**

### The Elastic Anisotropy Advantage (Revisited)

As noted in Chapter 47, a [001]-oriented single crystal has E[001] = 125 GPa vs. polycrystal E = 207 GPa.

For the same temperature gradient ΔT and thermal expansion coefficient α:
```
Thermal stress σ = E × α × ΔT
SX blade: σ = 125 GPa × 13×10⁻⁶/K × 200K = 325 MPa
PC blade: σ = 207 GPa × 13×10⁻⁶/K × 200K = 538 MPa
```

The [001] SX blade has 40% lower thermal stress → far fewer LCF cycles to initiation → much longer life.

This alone, independent of creep improvement, justifies the use of single crystals in HPT blades that experience strong thermal cycling.

---

## 9. The Impossible Requirement — Metal Hotter Than Its Melting Point

Here is the most extraordinary fact about modern jet engine operation:

**The gas entering the HPT is at 1,500–1,700°C. The Ni superalloy solidus is ~1,340°C. The GAS IS 150–360°C HOTTER THAN THE BLADE'S MELTING POINT.**

Yet the blade survives for 30,000 hours.

This is only possible because:

1. **Cooling reduces metal temperature below the solidus:** Even though the GAS is above the blade's melting point, the METAL is kept at 950–1,100°C (below solidus) by internal and film cooling.

2. **The cooling flow is thermally and aerodynamically designed to prevent hot gas from directly contacting the metal surface:** Film cooling creates a thin cool air layer (~700°C) that insulates the blade surface from the 1,600°C gas stream.

3. **The TBC adds 100–200°C of thermal insulation:** Between the gas and the metal, a ceramic TBC layer (yttria-stabilized zirconia, YSZ, Chapter 57) reduces heat flux through the wall.

**Analogy:** The blade is essentially operating like a sealed refrigerator in a furnace. The furnace (gas stream) is much hotter than the refrigerator contents (blade metal), but as long as the cooling system works, the contents remain cold.

The critical risk: if cooling fails (blocked cooling holes, TBC spallation, cooling passage cracks), the blade temperature rises rapidly toward the gas temperature → blade melts → engine failure.

---

## 10. How Blade Metal Temperature Is Controlled

**The cooling air source:** Cooling air is bled from the HPC before combustion. It has not been heated — it exits the HPC at ~600°C and ~30–50 atm pressure. This high-pressure cool air is then routed through the turbine disk and into the blade via the blade root.

**Cooling effectiveness is measured by:**
```
η_cooling = (T_gas - T_metal) / (T_gas - T_coolant)

High η = blade much cooler than gas = good cooling
Low η = blade approaches gas temperature = poor cooling, risk of failure
```

Typical η for modern HPT blades: 0.6–0.7 (metal temperature is 60–70% of the way from gas temperature toward coolant temperature).

**Cost of cooling:** Each kg of cooling air bled from the HPC represents work done by the compressor that does NOT contribute to engine thrust → efficiency penalty. Modern cooling designs use ~15–20% of HPC mass flow for blade cooling. Every 1% reduction in cooling flow at equivalent blade temperature is a significant engine efficiency improvement.

This drives the competition between:
- **Better materials** (higher allowable metal T → less cooling needed)
- **Better cooling design** (more efficient cooling → can cool more at lower flow penalty)
- **Better TBC** (more insulation → less heat into metal → less cooling needed)

All three are active research areas.

---

## 11. The Complete Loading State — Summary

An HPT blade in service simultaneously experiences:

```
Environment           Magnitude              Failure Mode
──────────────────────────────────────────────────────────────
Creep stress          200–250 MPa, 980°C    Creep rupture
Thermal stress        100–200 MPa cycling   Thermal fatigue (LCF)
Vibratory stress      50–150 MPa, 10³-10⁴Hz High-cycle fatigue (HCF)
Aerodynamic bending   50–100 MPa            Fatigue + creep
Oxidizing gas         O₂/H₂O at 1600°C gas Wall thinning, scale spallation
Hot corrosion         Na₂SO₄ at 800-950°C   Pitting, rapid oxidation
Thermal gradient      ΔT=200°C through wall Thermal fatigue, TBC spallation
Thermal cycling       1 cycle/flight (LCF)  Fatigue crack initiation
Erosion               Ingested particles    Wall thinning
```

No other engineering component faces this combination of loads simultaneously. The turbine blade represents the absolute frontier of materials engineering capability.

---

## Summary

- **Brayton cycle efficiency** increases with turbine inlet temperature (TIT) — every 25°C improvement in TIT saves ~1% fuel burn
- **Modern TIT**: 1,500–1,700°C; HPT blade metal temperature: 950–1,100°C (kept below alloy solidus by cooling)
- **Centrifugal stress**: ~200–250 MPa at blade root, at 980°C, for 30,000 hours → defines creep requirement
- **Chemical attack**: oxidation (O₂, H₂O), hot corrosion (Na₂SO₄), both requiring protective coatings
- **Fatigue**: ~30,000 LCF cycles (one per flight) + 10⁹–10¹² HCF cycles from vibration
- **"Impossible" operating condition**: gas temperature above blade solidus is managed by cooling + TBC
- **[001] crystal orientation**: lowers thermal stress by 40% vs polycrystal → critical for thermal fatigue life
- **Everything follows from thermodynamics**: higher TIT → better efficiency → more airline profit → more R&D investment into better blade materials

**Next chapter:** How the blade is loaded mechanically — the detailed stress analysis of HPT blade airfoil design, and what it means for material property requirements.

---

## Exercises

1. A gas turbine has TIT = 1,650°C and ambient temperature = 15°C. Calculate: (a) the ideal Carnot efficiency, (b) the ideal Brayton efficiency with pressure ratio r_p = 45 (use η = 1 - T₁/T₂ where T₂ = T₁ × r_p^((γ-1)/γ), γ=1.4), (c) If raising TIT by 100°C increases actual efficiency by 0.8%, how much fuel (in kg) is saved on a 12-hour flight consuming 80,000 kg of fuel?

2. An HPT blade has: mass = 0.5 kg, rotation speed = 10,500 RPM, disk radius = 420 mm, blade length = 90 mm (so r_CG = 420 + 45 = 465 mm). Calculate: (a) ω in rad/s, (b) centrifugal force in kN and equivalent weight in tonnes, (c) if the blade root section has area 900 mm², what is the blade root stress in MPa?

3. The cooling effectiveness formula η = (T_gas - T_metal)/(T_gas - T_coolant) applies. For TIT = 1,650°C, cooling air at T_coolant = 620°C, and required blade metal temperature T_metal = 980°C: (a) what cooling effectiveness η is needed? (b) If η improves from 0.65 to 0.70, what is the new blade metal temperature? (c) How might a new coating allow η to decrease (less cooling flow) while maintaining the same metal temperature?

4. A blade experiences thermal stress of σ_thermal = E × α × ΔT where ΔT = 180°C through the blade wall, α = 13×10⁻⁶/K. Compare thermal stress for: (a) polycrystalline blade E = 207 GPa, (b) [001] SX blade E = 125 GPa, (c) [111] SX blade E = 294 GPa. Which orientation is worst for thermal fatigue and why? Why is [111] never used as the primary growth direction?

5. The blade root stress calculated above is sustained for 30,000 hours at 980°C. Using the Larson-Miller parameter (from Chapter 16) LMP = T(log t_r + C) with C = 20, T in Kelvin, t_r in hours: what is the LMP for these conditions? If CMSX-4 can sustain a stress of 200 MPa at LMP = 32,000, is this blade within safe design limits?
