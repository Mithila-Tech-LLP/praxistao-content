# Chapter 41: Gas Laws

> **"A few simple laws connect pressure, volume, temperature, and amount of gas — and they work for every gas in the universe. These are the tools chemists and engineers use every day."**

---

## Table of Contents

- [41.1 Introduction: State Variables of a Gas](#411-introduction-state-variables-of-a-gas)
- [41.2 Boyle's Law: Pressure and Volume](#412-boyles-law-pressure-and-volume)
- [41.3 Charles's Law: Volume and Temperature](#413-charless-law-volume-and-temperature)
- [41.4 Gay-Lussac's Law: Pressure and Temperature](#414-gay-lussacs-law-pressure-and-temperature)
- [41.5 The Combined Gas Law](#415-the-combined-gas-law)
- [41.6 Avogadro's Law](#416-avogadros-law)
- [41.7 The Ideal Gas Law](#417-the-ideal-gas-law)
- [41.8 Dalton's Law of Partial Pressures](#418-daltons-law-of-partial-pressures)
- [41.9 Standard Conditions](#419-standard-conditions)
- [41.10 Worked Problems](#4110-worked-problems)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 41.1 Introduction: State Variables of a Gas

The **state** of a gas is completely described by four variables:

```
FOUR STATE VARIABLES:

  P = Pressure  (Pa = N/m²)
  V = Volume    (m³)
  T = Temperature (K — must use KELVIN, not Celsius)
  n = Amount    (moles)

Change ANY one of these, and you change the others.
```

### The Kelvin Scale

The gas laws require **absolute temperature** in Kelvin:

```
T(K) = T(°C) + 273.15

Examples:
  0°C   = 273 K   (freezing water)
  100°C = 373 K   (boiling water)
  -273°C = 0 K    (absolute zero — coldest possible)
  25°C  = 298 K   (room temperature)
```

**Why Kelvin?** Because at 0 K, molecular motion stops. The gas laws describe averages of molecular behavior, and negative temperatures would give negative energies — physically meaningless.

---

## 41.2 Boyle's Law: Pressure and Volume

**Boyle's Law (1662):** At constant temperature and amount of gas, pressure and volume are inversely proportional.

```
P ∝ 1/V    (constant T and n)

Or:   P × V = constant

Or:   P₁V₁ = P₂V₂
```

### Physical Explanation

If you compress a gas (decrease V), molecules hit the walls more often per second → more pressure.

```
BIG VOLUME:                 SMALL VOLUME:
  +--------+                +---+
  | • •  • |                |•• |
  |  •   • |                |• •|
  | •  •   |                +---+
  +--------+
  
  Molecules hit walls      Molecules hit walls
  less often               much more often
  Lower pressure           Higher pressure
```

### Graph of Boyle's Law

```
P                           P
 |  *                        |          *
 |   *                       |       *
 |    **                     |    *
 |       **                  | *
 |          ****             +---------> 1/V
 +---------> V              
 
  (P vs V: hyperbola)       (P vs 1/V: straight line through origin)
```

The graph of P vs 1/V is a straight line through the origin — a good linearization to confirm Boyle's Law experimentally.

### Worked Example 41.1

A gas occupies 4 L at 2 atm pressure. What volume will it occupy at 5 atm? (temperature constant)

**Solution:**

P₁V₁ = P₂V₂
2 × 4 = 5 × V₂
V₂ = 8/5 = **1.6 L**

---

## 41.3 Charles's Law: Volume and Temperature

**Charles's Law (1787):** At constant pressure and amount, volume is directly proportional to absolute temperature.

```
V ∝ T    (constant P and n)

Or:   V/T = constant

Or:   V₁/T₁ = V₂/T₂
```

### Physical Explanation

Increasing temperature → molecules move faster → hit walls harder → if pressure stays constant, walls must move out → volume increases.

### Graph of Charles's Law

```
V
 |            /
 |          /
 |        /
 |      /
 |    /
 |  /
 | /
 +-----------> T (in Kelvin)
  0

Straight line through origin (when T is in Kelvin).
If plotted vs Celsius, line intersects x-axis at -273°C.
```

The x-intercept at -273°C was historically used to estimate absolute zero!

### Worked Example 41.2

A balloon contains 2 L of air at 20°C. What is its volume at 60°C? (pressure constant)

**Solution:**

T₁ = 20 + 273 = 293 K
T₂ = 60 + 273 = 333 K

V₁/T₁ = V₂/T₂
2/293 = V₂/333
V₂ = 2 × 333/293 = **2.27 L**

---

## 41.4 Gay-Lussac's Law: Pressure and Temperature

**Gay-Lussac's Law (1808):** At constant volume and amount, pressure is directly proportional to absolute temperature.

```
P ∝ T    (constant V and n)

Or:   P/T = constant

Or:   P₁/T₁ = P₂/T₂
```

### Physical Explanation

Increasing temperature → molecules move faster → hit walls harder → pressure increases.

### Practical Example

Car tires: pressure increases as you drive (friction heats the air inside). Always check tire pressure when cold!

Aerosol cans: have "do not incinerate" warning because heating increases pressure and can cause explosion.

### Worked Example 41.3

A gas in a sealed container is at 1.5 atm and 27°C.

What pressure does it exert at 127°C?

**Solution:**

T₁ = 27 + 273 = 300 K
T₂ = 127 + 273 = 400 K

P₁/T₁ = P₂/T₂
1.5/300 = P₂/400
P₂ = 1.5 × 400/300 = **2.0 atm**

---

## 41.5 The Combined Gas Law

Combining Boyle's, Charles's, and Gay-Lussac's laws:

```
(P × V) / T = constant

Or:   P₁V₁/T₁ = P₂V₂/T₂
```

Use when P, V, AND T all change, but n stays constant.

### Worked Example 41.4

A gas occupies 500 mL at 25°C and 100 kPa.

What is its volume at 150°C and 200 kPa?

**Solution:**

T₁ = 298 K, T₂ = 423 K
P₁ = 100 kPa, P₂ = 200 kPa
V₁ = 500 mL

P₁V₁/T₁ = P₂V₂/T₂

V₂ = V₁ × (P₁/P₂) × (T₂/T₁)
   = 500 × (100/200) × (423/298)
   = 500 × 0.5 × 1.42
   = **355 mL**

---

## 41.6 Avogadro's Law

**Avogadro's Law (1811):** At constant temperature and pressure, equal volumes of any gas contain equal numbers of molecules.

```
V ∝ n    (constant T and P)

Or:   V/n = constant

Or:   V₁/n₁ = V₂/n₂
```

This means 1 mole of ANY ideal gas at the same T and P occupies the same volume.

At STP (standard temperature and pressure: 0°C, 1 atm):
- 1 mole of any ideal gas = 22.4 L

At room temperature (25°C, 1 atm):
- 1 mole = 24.5 L

---

## 41.7 The Ideal Gas Law

Combining all four gas laws into one:

```
PV = nRT

where:
  P = pressure (Pa)
  V = volume (m³)
  n = amount of gas (moles)
  R = 8.314 J/(mol·K) = universal gas constant
  T = temperature (K)
```

This is the most important equation in the gas laws.

### Alternative Forms

```
Using number of molecules N instead of moles n:
  PV = NkT

where k = R/N_A = 1.38 × 10⁻²³ J/K (Boltzmann's constant)

Using density ρ:
  P = ρRT/M

where M = molar mass (kg/mol)
```

### Worked Example 41.5

Find the volume of 2 moles of ideal gas at 25°C and 1 atm (101,325 Pa).

**Solution:**

T = 298 K, n = 2 mol, P = 101,325 Pa

V = nRT/P = 2 × 8.314 × 298 / 101,325 = 4955/101,325 ≈ **0.0489 m³ = 48.9 L**

(2 × 24.5 L = 49 L at room conditions — checks out)

### Worked Example 41.6

A container of volume 10 L holds gas at 300 K and pressure 2 × 10⁵ Pa.

How many moles of gas are in the container?

**Solution:**

n = PV/(RT) = (2×10⁵ × 10×10⁻³) / (8.314 × 300) = 2000 / 2494 ≈ **0.80 mol**

---

## 41.8 Dalton's Law of Partial Pressures

In a mixture of gases that don't react with each other, each gas behaves independently:

```
P_total = P₁ + P₂ + P₃ + ...

where P₁, P₂, etc. are the partial pressures each gas would exert if it alone occupied the container.
```

The partial pressure of gas i is:

```
P_i = n_i × RT / V = x_i × P_total

where x_i = n_i / n_total is the mole fraction
```

### Example: Air

Air is roughly 78% N₂, 21% O₂, 1% Ar by mole fraction.

At atmospheric pressure (101.3 kPa):
- Partial pressure of N₂ ≈ 0.78 × 101.3 = **79.0 kPa**
- Partial pressure of O₂ ≈ 0.21 × 101.3 = **21.3 kPa**
- Partial pressure of Ar ≈ 0.01 × 101.3 = **1.0 kPa**

---

## 41.9 Standard Conditions

Two common sets of standard conditions:

```
STP (Standard Temperature and Pressure):
  T = 273.15 K (0°C)
  P = 101,325 Pa (1 atm)
  Molar volume = 22.4 L/mol

SATP (Standard Ambient Temperature and Pressure):
  T = 298.15 K (25°C)
  P = 100,000 Pa (100 kPa, approximately 1 bar)
  Molar volume = 24.8 L/mol
```

---

## 41.10 Worked Problems

### Worked Example 41.7 — Finding Molar Mass from Density

The density of an unknown gas at STP is 2.86 g/L. Find its molar mass.

**Solution:**

At STP, 1 mol occupies 22.4 L.

Molar mass = density × molar volume = 2.86 g/L × 22.4 L/mol = **64.1 g/mol**

(This is sulfur dioxide, SO₂: M = 32 + 2×16 = 64 g/mol ✓)

### Worked Example 41.8 — Combined Law Application

A diver's tank contains air at 200 atm and 15°C. The volume of the tank is 12 L.

How many moles of air does it contain?

**Solution:**

P = 200 × 101,325 = 20,265,000 Pa
V = 12 L = 0.012 m³
T = 15 + 273 = 288 K

n = PV/(RT) = (20,265,000 × 0.012) / (8.314 × 288)
  = 243,180 / 2394
  ≈ **101.6 mol**

(This is equivalent to 2285 L at STP — enough air to breathe for a significant time.)

---

## Summary

- All gas laws require **Kelvin** temperature: T(K) = T(°C) + 273
- **Boyle's Law**: PV = constant (constant T, n) → P₁V₁ = P₂V₂
- **Charles's Law**: V/T = constant (constant P, n) → V₁/T₁ = V₂/T₂
- **Gay-Lussac's Law**: P/T = constant (constant V, n) → P₁/T₁ = P₂/T₂
- **Combined Gas Law**: P₁V₁/T₁ = P₂V₂/T₂ (constant n)
- **Avogadro's Law**: V ∝ n (constant T, P); 1 mol at STP = 22.4 L
- **Ideal Gas Law**: PV = nRT; R = 8.314 J/(mol·K)
- **Dalton's Law**: P_total = sum of partial pressures
- Ideal gas law is an approximation; real gases deviate at high P or low T

---

## Key Equations

```
Boyle's Law (constant T, n):
  P₁V₁ = P₂V₂

Charles's Law (constant P, n):
  V₁/T₁ = V₂/T₂

Gay-Lussac's Law (constant V, n):
  P₁/T₁ = P₂/T₂

Combined Gas Law (constant n):
  P₁V₁/T₁ = P₂V₂/T₂

Ideal Gas Law:
  PV = nRT
  R = 8.314 J/(mol·K)
  (or PV = NkT, k = 1.38 × 10⁻²³ J/K)

Molar volume at STP (0°C, 1 atm):
  V_m = 22.4 L/mol

Dalton's Law:
  P_total = P₁ + P₂ + P₃ + ...

Temperature conversion:
  T(K) = T(°C) + 273.15
```
