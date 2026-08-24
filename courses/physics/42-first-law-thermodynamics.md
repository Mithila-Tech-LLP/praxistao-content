# Chapter 42: The First Law of Thermodynamics

> **"Energy can never be created or destroyed — only converted. The first law of thermodynamics is just this principle applied to heat engines, refrigerators, and every process in nature."**

---

## Table of Contents

- [42.1 Thermodynamic Systems](#421-thermodynamic-systems)
- [42.2 Internal Energy](#422-internal-energy)
- [42.3 Heat and Work](#423-heat-and-work)
- [42.4 The First Law of Thermodynamics](#424-the-first-law-of-thermodynamics)
- [42.5 Thermodynamic Processes](#425-thermodynamic-processes)
- [42.6 Work Done by a Gas](#426-work-done-by-a-gas)
- [42.7 Heat Capacity of Gases](#427-heat-capacity-of-gases)
- [42.8 Adiabatic Processes](#428-adiabatic-processes)
- [42.9 Cyclic Processes and Heat Engines](#429-cyclic-processes-and-heat-engines)
- [42.10 Worked Examples](#4210-worked-examples)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 42.1 Thermodynamic Systems

A **thermodynamic system** is any collection of matter we choose to study.

```
SYSTEM TYPES:

OPEN SYSTEM:
  - Can exchange matter AND energy with surroundings
  - Example: boiling pot of water (steam escapes, heat exchanges)

CLOSED SYSTEM:
  - Can exchange energy, but NOT matter
  - Example: gas in a cylinder with a piston

ISOLATED SYSTEM:
  - No exchange of energy OR matter
  - Example: thermos flask (approximately isolated)
  
  SYSTEM         SURROUNDINGS
  +------+       (everything else)
  |      |
  |      |   <-- boundary
  |      |
  +------+
```

**State functions** describe the current state: P, V, T, internal energy U.
**Path functions** depend on how the process happened: heat Q, work W.

---

## 42.2 Internal Energy

**Internal energy (U)** is the total energy stored inside a system — kinetic energy of all molecules (translational, rotational, vibrational) plus potential energy of intermolecular forces.

For an ideal gas, there are no intermolecular forces, so:

```
U = (kinetic energy of all molecules)
  = N × (3/2)kT   (monatomic)
  = (3/2)nRT       (monatomic)
  = (5/2)nRT       (diatomic, at moderate temperature)
```

Key property: **U depends only on temperature for an ideal gas.**

If T increases → U increases.
If T stays constant → U stays constant, regardless of what happens to P or V.

---

## 42.3 Heat and Work

There are two ways to change a system's internal energy:

**Heat (Q)**: energy transferred due to temperature difference.
- Positive Q: heat flows INTO the system (system absorbs heat)
- Negative Q: heat flows OUT of system (system releases heat)

**Work (W)**: energy transferred by mechanical means (compression/expansion).
- Positive W: work done BY the gas (gas expands, pushes piston out)
- Negative W: work done ON the gas (gas compressed)

```
WARNING: Different textbooks define W differently!

Convention A (physics): W = work done BY the gas
  -> expansion: W > 0

Convention B (engineering): W = work done ON the gas
  -> compression: W > 0

This chapter uses PHYSICS convention: W = work done BY the gas.
```

---

## 42.4 The First Law of Thermodynamics

The first law is just conservation of energy for a thermodynamic system:

```
ΔU = Q - W

or equivalently: Q = ΔU + W

where:
  ΔU = change in internal energy
  Q   = heat added to the system
  W   = work done BY the system
```

In words: "The heat added to a system equals the increase in internal energy plus the work done by the system."

```
ENERGY BUDGET:

  Q (heat in)
      |
      v
  [SYSTEM]  --> W (work out)
      |
      v
   ΔU (energy stored)

Q = ΔU + W

If you add heat Q:
  - some goes to increasing internal energy (ΔU)
  - some goes to doing work (W, like pushing a piston)
```

---

## 42.5 Thermodynamic Processes

### Isothermal Process (constant temperature)

- T = constant → ΔU = 0 (for ideal gas)
- Therefore: Q = W (all heat absorbed → work done by gas)

```
PV DIAGRAM:
P
|*
| *
|  *    <- isothermal curve (P ∝ 1/V, a hyperbola)
|    *
|      *
+---------> V
```

### Isovolumetric (Isochoric) Process (constant volume)

- V = constant → W = 0 (no expansion, no work)
- Therefore: Q = ΔU (all heat goes into internal energy)

```
PV DIAGRAM:
P
|
|     *  <- higher P after heating
|     |
|     *  <- lower P before heating
|
+---------> V
     ^
     constant V (vertical line)
```

### Isobaric Process (constant pressure)

- P = constant → W = PΔV
- Q = ΔU + PΔV (heat goes partly into internal energy, partly into work)

```
PV DIAGRAM:
P
|          
|  *--------*   <- same P throughout
|           
+---------> V
   V₁       V₂

Area under line = P × ΔV = work done
```

### Adiabatic Process (no heat exchange)

- Q = 0 → ΔU = -W
- If gas expands (W > 0) → ΔU < 0 → temperature falls
- If gas compressed (W < 0) → ΔU > 0 → temperature rises

```
EXAMPLES:
  Diesel engine: air compressed adiabatically → heats up → ignites diesel fuel
  Bike pump: feels hot when you pump fast (adiabatic compression)
  Expanding gas from aerosol: feels cold (adiabatic expansion)
```

---

## 42.6 Work Done by a Gas

When a gas expands against a piston:

```
     P
     |
     v
  ------
  |    |  <- gas
  |    |
  ------
     |
   piston
     |
     F

Force = P × A  (pressure × area)
Work = F × dx = P × A × dx = P × dV
```

For a constant pressure process:

```
W = P × ΔV = P × (V₂ - V₁)
```

For a varying pressure process:

```
W = area under the P-V graph (integral of P dV)
```

### PV Diagrams

The **PV diagram** (pressure-volume graph) is the most powerful tool in thermodynamics.

```
KEY RULE: Work = area under the curve on a PV diagram

Expansion (V increases, left to right): W > 0 (positive work, gas does work)
Compression (V decreases, right to left): W < 0 (negative work, work done ON gas)

P
|    *------*          <- isobaric expansion (shaded area = W)
|    |      |
|    |      |
+----+------+----> V
     V₁     V₂

W = P × (V₂ - V₁) = area of rectangle
```

---

## 42.7 Heat Capacity of Gases

### Heat Capacity at Constant Volume (Cv)

When gas is heated at constant volume (no expansion, no work):

```
Q = ΔU = nCvΔT

For monatomic gas:
  Cv = (3/2)R = 12.5 J/(mol·K)

For diatomic gas:
  Cv = (5/2)R = 20.8 J/(mol·K)
```

### Heat Capacity at Constant Pressure (Cp)

When gas is heated at constant pressure (gas expands, does work):

```
Q = nCpΔT

Cp = Cv + R   (Mayer's relation)

For monatomic: Cp = (5/2)R = 20.8 J/(mol·K)
For diatomic: Cp = (7/2)R = 29.1 J/(mol·K)
```

Why Cp > Cv? Because at constant pressure, some heat goes into doing work (expanding). You need more heat to get the same temperature rise.

```
RATIO γ = Cp/Cv:

  Monatomic: γ = (5/2R)/(3/2R) = 5/3 ≈ 1.67  (He, Ne, Ar)
  Diatomic:  γ = (7/2R)/(5/2R) = 7/5 = 1.4   (N₂, O₂, air)
  
γ appears in many important equations (speed of sound, adiabatic processes).
```

---

## 42.8 Adiabatic Processes

For an adiabatic process (Q = 0), the PV relationship follows:

```
PV^γ = constant

(instead of PV = constant for isothermal)

Also:
  TV^(γ-1) = constant
  P^(1-γ) × T^γ = constant
```

The adiabatic curve is **steeper** than the isothermal curve on a PV diagram:

```
P
|  * ← adiabatic (steeper)
|   *
|    ** ← isothermal
|      ***
|         ****
+--------------> V
```

Because in adiabatic expansion, the temperature drops (no heat input to maintain T).

---

## 42.9 Cyclic Processes and Heat Engines

A **cyclic process** returns the system to its initial state:

- After a full cycle: ΔU = 0 (same state → same U)
- Therefore: Q_net = W_net (net heat absorbed = net work done)

```
PV DIAGRAM OF A CYCLE:

P
|    *--------*
|   /          \
|  *            \
|   \           *
|    *----------*
+---------------------> V

The AREA ENCLOSED by the cycle = net work done per cycle.

Clockwise loop: W_net > 0 (engine does work)
Counter-clockwise loop: W_net < 0 (refrigerator/pump needs work)
```

### Efficiency of a Heat Engine

A heat engine absorbs heat Q_h from a hot source, does work W, and dumps heat Q_c to a cold sink:

```
HOT RESERVOIR (T_h)
      |
      | Q_h (heat absorbed)
      v
   [ENGINE] -----> W (work output)
      |
      | Q_c (heat dumped)
      v
COLD RESERVOIR (T_c)

Efficiency = W/Q_h = (Q_h - Q_c)/Q_h = 1 - Q_c/Q_h

Maximum possible efficiency (Carnot limit):
  η_max = 1 - T_c/T_h  (temperatures in Kelvin!)
```

No real engine can exceed the Carnot efficiency.

---

## 42.10 Worked Examples

### Worked Example 42.1

5 J of heat is added to a gas, and the gas does 3 J of work. Find the change in internal energy.

**Solution:**

ΔU = Q - W = 5 - 3 = **+2 J** (internal energy increases by 2 J)

### Worked Example 42.2

A gas is compressed from 10 L to 4 L at constant pressure of 200 kPa.

(a) Find the work done by the gas.
(b) If the internal energy increases by 1000 J, how much heat was exchanged?

**Solution:**

(a) W = P × ΔV = 200,000 × (0.004 - 0.010) = 200,000 × (-0.006) = **-1200 J**

(Work done BY gas = -1200 J, meaning 1200 J of work done ON gas.)

(b) ΔU = Q - W
    1000 = Q - (-1200)
    Q = 1000 - 1200 = **-200 J** (200 J leaves the gas)

### Worked Example 42.3 — Heat Engine

A heat engine operates between 600 K and 300 K.

(a) Maximum efficiency (Carnot)?
(b) If Q_h = 5000 J per cycle, what is maximum work output?

**Solution:**

(a) η = 1 - T_c/T_h = 1 - 300/600 = **0.5 = 50%**

(b) W_max = η × Q_h = 0.5 × 5000 = **2500 J**

---

## Summary

- **System types**: open (exchanges matter + energy), closed (energy only), isolated (neither)
- **Internal energy U**: total microscopic energy; for ideal gas, U = (f/2)nRT, depends only on T
- **Heat Q**: energy transfer due to temperature difference; positive = into system
- **Work W**: mechanical energy transfer; W = PΔV for constant pressure
- **First Law**: ΔU = Q - W (conservation of energy for thermodynamic systems)
- **Isothermal**: ΔT = 0, ΔU = 0, Q = W; uses PV = const
- **Isochoric**: ΔV = 0, W = 0, Q = ΔU
- **Isobaric**: ΔP = 0, W = PΔV, Q = ΔU + PΔV
- **Adiabatic**: Q = 0, ΔU = -W; PV^γ = constant; temperature changes
- **PV diagram**: work = area under curve; cycle area = net work
- **Heat engine efficiency**: η = W/Q_h; maximum (Carnot): η = 1 - Tc/Th

---

## Key Equations

```
First Law:
  ΔU = Q - W

Internal energy (ideal gas):
  ΔU = nCvΔT
  Cv = (3/2)R  (monatomic)
  Cv = (5/2)R  (diatomic)

Work:
  W = P × ΔV  (constant pressure)
  W = area under P-V curve

Heat capacities:
  Cp = Cv + R  (Mayer's relation)
  γ = Cp/Cv = 5/3 (monatomic), 7/5 (diatomic)

Process relationships:
  Isothermal:  PV = const; ΔU = 0
  Adiabatic:   PV^γ = const; Q = 0
  Isochoric:   ΔV = 0; W = 0
  Isobaric:    ΔP = 0; W = PΔV

Carnot efficiency:
  η_max = 1 - Tc/Th  (temperatures in Kelvin)
```
