# Chapter 50: Ohm's Law and Resistance

> **"Ohm's Law is the most useful equation in all of electronics. Master it and you can analyze almost any circuit."**

---

## Table of Contents

1. [Introduction: What Is Resistance?](#introduction-what-is-resistance)
2. [Ohm's Law: V = IR](#ohms-law-v--ir)
3. [The Microscopic Picture of Resistance](#the-microscopic-picture-of-resistance)
4. [Ohmic vs Non-Ohmic Materials](#ohmic-vs-non-ohmic-materials)
5. [I-V Characteristic Graphs](#i-v-characteristic-graphs)
6. [Resistivity: R = ρL/A](#resistivity-r--ρla)
7. [Resistivity of Common Materials](#resistivity-of-common-materials)
8. [Temperature Effects on Resistance](#temperature-effects-on-resistance)
9. [Superconductivity](#superconductivity)
10. [Practical Resistance in Circuits](#practical-resistance-in-circuits)
11. [Worked Examples](#worked-examples)
12. [Summary](#summary)
13. [Key Equations](#key-equations)

---

## Introduction: What Is Resistance?

When you push a trolley along a smooth floor, it rolls easily. Push the same trolley through thick mud, and it slows down drastically. The mud resists the motion.

Electric current encounters the same kind of opposition as it flows through materials. **Resistance** is the property of a material that opposes the flow of electric current. Every material has some resistance — some materials have very low resistance (conductors like copper), others have extremely high resistance (insulators like rubber).

**Resistance** is given the symbol R and measured in **Ohms** (Ω, the Greek capital letter omega).

The concept of resistance was developed by Georg Simon Ohm, a German physicist who published his findings in 1827. His work was initially dismissed by the German scientific establishment — one critic called it "a tissue of naked fantasy" — but was later recognized as foundational to electrical science. The unit of resistance was named in his honor.

---

## Ohm's Law: V = IR

**Ohm's Law** states that for many conductors, the current through the conductor is proportional to the voltage across it, provided that physical conditions (especially temperature) remain constant.

Mathematically:

```
V = I × R

or equivalently:

        V                   V
I = -------       R = -------
        R                   I
```

Where:
- V = voltage (potential difference) across the component, in Volts (V)
- I = current through the component, in Amperes (A)
- R = resistance of the component, in Ohms (Ω)

### The Ohm's Law Triangle

Many students find this memory triangle helpful:

```
         +-------+
         |   V   |
         +---+---+
         |   |   |
         | I | R |
         +---+---+

Cover up what you want to find:
- Cover V: V = I × R
- Cover I: I = V / R
- Cover R: R = V / I
```

### Intuition for Ohm's Law

Think of it like a water analogy:
- **Voltage** (V) is like **water pressure** — it drives the flow
- **Current** (I) is like the **flow rate** — how much water passes per second
- **Resistance** (R) is like the **narrowness of the pipe** — it restricts flow

```
HIGH VOLTAGE, LOW RESISTANCE:     HIGH VOLTAGE, HIGH RESISTANCE:
      
  ++++++++++++                      ++++++++++++
  (High pressure)                   (High pressure)
  =====[wide pipe]=====             =====[narrow pipe]=====
  Large flow rate (high I)          Small flow rate (low I)
```

Double the voltage (pressure): double the current.
Double the resistance (narrower pipe): halve the current.

### What Counts as 1 Ohm?

One Ohm is defined as the resistance that allows exactly 1 Ampere to flow when 1 Volt is applied.

```
1 Ω = 1 V / 1 A
```

Practical resistances range widely:

```
Component                     Typical Resistance
--------------------------------------------------
Copper wire (1m, 1mm dia)     about 0.02 Ω   (very low)
LED                           about 10-100 Ω
Resistor (typical)            1 Ω to 10 MΩ
Human body (dry skin)         10,000 to 100,000 Ω
Human body (wet skin)         1,000 Ω (dangerous!)
Glass                         about 10^12 Ω  (very high)
Rubber                        about 10^13 Ω  (extremely high)
```

---

## The Microscopic Picture of Resistance

Why do materials resist current? The answer lies in what happens when electrons move through a crystal lattice.

### Electron Collisions with Ions

When a battery is connected to a metal wire, it creates an electric field throughout the wire. Free electrons start to accelerate in response to this field. However, they do not accelerate indefinitely. The crystal lattice of the metal is made of positive ions vibrating around fixed positions:

```
METAL LATTICE (Side View):

        Cu+   Cu+   Cu+   Cu+
         |     |     |     |
    ---> e⁻  collides!     e⁻ --->
         |     |     |     |
        Cu+   Cu+   Cu+   Cu+

    Electron accelerates, then collides with ion.
    After collision: electron moves slower again.
    Then accelerates again. Then collides. And so on.
```

This cycle of acceleration and collision is what creates **electrical resistance**. The energy transferred in each collision heats the metal — this is why wires and resistors warm up when current flows.

### What Affects the Collision Rate?

1. **Temperature**: Higher temperature means ions vibrate more violently and take up more space, increasing the collision rate and thus the resistance.

2. **Impurities and defects**: Impurity atoms in the lattice disrupt the regular crystal structure and increase collisions. Pure copper has lower resistance than copper with impurities.

3. **Length of wire**: More length means more collisions. Resistance scales with length.

4. **Cross-sectional area**: Wider wire has more paths for electrons to travel, reducing resistance.

---

## Ohmic vs Non-Ohmic Materials

Not all materials obey Ohm's Law. Materials are classified as **ohmic** or **non-ohmic** based on whether their resistance stays constant.

### Ohmic Materials (Resistors, Metals at Constant Temperature)

An **ohmic conductor** obeys Ohm's Law: its resistance R stays constant regardless of the voltage applied. The I-V graph is a straight line through the origin.

```
OHMIC MATERIAL - I-V GRAPH:

Current (I)
  ^
  |         /
  |        /
  |       /  Straight line
  |      /   through origin
  |     /
  |    /    slope = 1/R
  |   /
  |  /
  | /
  |/
  +-------------> Voltage (V)
  0

Steeper slope = lower resistance (more current for same voltage)
Shallower slope = higher resistance (less current for same voltage)
```

Examples: Metal resistors (at constant temperature), nichrome wire (heating element material)

### Non-Ohmic Materials

A **non-ohmic device** does not have constant resistance. The relationship between V and I is not linear.

#### Filament Bulb (Tungsten Wire)

At low temperatures, tungsten has low resistance. As more current flows, the filament heats up significantly (to over 2000°C when glowing!). At higher temperatures, the resistance increases dramatically.

```
FILAMENT BULB - I-V GRAPH:

Current (I)
  ^
  |       /----
  |      /    (curve flattens as
  |     /      R increases with temp)
  |    /
  |   /
  |  /
  | / (steep at start - low R when cold)
  |/
  +-------------> Voltage (V)
  0
```

The curve gets flatter at higher voltages because R increases as the filament gets hotter.

#### Diode (Semiconductor)

A **diode** is a semiconductor device that allows current to flow easily in one direction (forward bias) but barely at all in the other direction (reverse bias). Its I-V graph is dramatically non-linear.

```
DIODE - I-V GRAPH:

Current (I)
  ^
  |              /
  |             / (forward bias:
  |            /   large current once
  |           /    threshold voltage ~0.7V reached)
  |          /
  |         /
--|---------|--+---> Voltage (V)
  |     ~0.7V  |
  |            | (reverse bias: almost no current)
  |            |
              (tiny leakage current
               or breakdown)
```

Below the threshold voltage (~0.6-0.7V for silicon), almost no current flows. Above it, current increases rapidly. In reverse direction, current is essentially blocked — this is what makes diodes useful as rectifiers (converting AC to DC).

#### Thermistor

A **thermistor** is a resistor made from semiconductor material whose resistance changes significantly with temperature.

```
THERMISTOR - RESISTANCE vs TEMPERATURE:

Resistance (R)
  ^
  |\
  | \    NTC thermistor:
  |  \   Negative Temperature Coefficient
  |   \  (R decreases as temperature increases)
  |    \
  |     \---___
  |           ---___
  +-------------------------> Temperature (T)
  
PTC thermistor (less common):
Resistance increases with temperature.
```

**NTC (Negative Temperature Coefficient) thermistors** are used in:
- Thermometers (smartphone temperature sensors)
- Resettable fuses (trip when too hot)
- Circuit protection (limit current surge on startup)

---

## I-V Characteristic Graphs

I-V characteristic graphs (current vs voltage) are the standard way to describe how a component behaves.

### Key Features to Read from an I-V Graph

```
1. SLOPE = 1/R
   (Steeper slope means lower resistance)

2. LINEARITY:
   - Straight line through origin = ohmic (constant R)
   - Curved line = non-ohmic (variable R)

3. SYMMETRY:
   - Symmetric about origin = same R in both directions (resistor)
   - Asymmetric = different behavior in + and - directions (diode)

4. THRESHOLD:
   - Sudden onset of conduction = threshold voltage (diode, LED)
```

### Reading Resistance from an I-V Graph

To find resistance at any point on a non-linear I-V graph, take the V and I values at that point and apply R = V/I.

```
Example: At point (V = 4V, I = 2A) on a filament bulb graph:
    R = V/I = 4/2 = 2 Ω

At a different point (V = 8V, I = 3A):
    R = V/I = 8/3 = 2.67 Ω

The resistance increased as voltage (and thus temperature) increased.
```

---

## Resistivity: R = ρL/A

The resistance of a piece of wire depends on:
1. What it is made of (its **resistivity**)
2. Its length (longer = more resistance)
3. Its cross-sectional area (wider = less resistance)

These are combined in the resistivity formula:

```
        ρ × L
R  =  --------
          A
```

Where:
- R = resistance (Ω)
- ρ (rho) = resistivity of the material (Ω·m)
- L = length of the conductor (m)
- A = cross-sectional area (m²)

### Physical Reasoning

**Length**: Doubling the length is like connecting two identical resistors in series — resistance doubles.

**Area**: Doubling the area is like connecting two identical resistors in parallel — resistance halves.

```
EFFECT OF LENGTH:
   [===R===]  length L, resistance R
   
   [===R===][===R===]  length 2L, resistance 2R
   (twice as long = twice the resistance)

EFFECT OF AREA:
   [===R===]  area A, resistance R
   
   [===R===]  area 2A = two wires side by side
   [===R===]  resistance = R/2
   (twice the area = half the resistance)
```

---

## Resistivity of Common Materials

**Resistivity** ρ is a property of the material itself (not the wire dimensions). It is measured in Ohm-meters (Ω·m).

```
Material              Resistivity (Ω·m)         Category
-----------------------------------------------------------
Silver                1.59 × 10^-8              Best conductor
Copper                1.72 × 10^-8              Excellent conductor
Gold                  2.44 × 10^-8              Good conductor (corrosion-resistant)
Aluminium             2.82 × 10^-8              Good conductor (lightweight)
Tungsten              5.6 × 10^-8               High melting point
Nichrome              1.10 × 10^-6              Resistor/heating element
Carbon (graphite)     3-60 × 10^-5              Semi-conductor / electrode
Germanium             4.6 × 10^-1               Semiconductor
Silicon               6.4 × 10^2                Semiconductor
Glass                 10^10 to 10^14            Insulator
Rubber                10^13 to 10^15            Excellent insulator
```

Key observations:
- Silver is technically the best conductor, but copper is used in wires because it is cheaper
- Gold is used in connectors because it does not tarnish/corrode
- Nichrome (nickel-chromium alloy) is used in heating elements (toasters, electric hobs) because it has high resistance and does not melt
- Semiconductors span many orders of magnitude between conductors and insulators

---

## Temperature Effects on Resistance

Temperature profoundly affects resistance, but in different ways for different materials.

### Metals: Resistance Increases with Temperature

In metals, the dominant effect is thermal vibration of the crystal lattice. As temperature rises:
- Ions vibrate more vigorously
- Electrons collide with them more frequently
- Current flow is impeded more
- **Resistance increases**

```
METAL - R vs T GRAPH:

Resistance (R)
  ^
  |          /
  |         /
  |        /  Linear increase
  |       /   for metals
  |      /
  |     /
  |    /
  +-------------> Temperature (T)
  0
  
  R increases with temperature (positive TCR)
```

The relationship is approximately linear:

```
R_T = R_0 × (1 + α × ΔT)

R_T = resistance at temperature T
R_0 = resistance at reference temperature (usually 20°C)
α   = temperature coefficient of resistance (per °C)
ΔT  = change in temperature from reference
```

For copper: α ≈ 0.0039 per °C

### Semiconductors: Resistance Decreases with Temperature

In semiconductors like silicon, the situation is reversed. At room temperature, relatively few electrons have enough energy to break free and conduct. As temperature increases:
- More electrons gain enough energy to become conduction electrons
- More charge carriers become available
- Current flows more easily
- **Resistance decreases**

```
SEMICONDUCTOR - R vs T GRAPH:

Resistance (R)
  ^
  |\
  | \
  |  \  Exponential decrease
  |   \  for semiconductors
  |    \
  |     \__
  |        ---___
  +-------------> Temperature (T)
  0
```

This is why semiconductors heat up more when they carry current, and why they need thermal management (heatsinks, fans) in computers.

### The Temperature Coefficient Summary

```
Material Type          Effect of Heating    Temperature Coefficient
-------------------------------------------------------------------
Metals                 R increases          Positive (α > 0)
Semiconductors         R decreases          Negative (α < 0)
Special alloys         R barely changes     Near zero (α ≈ 0)
(Manganin, constantan)
```

Manganin and constantan are special alloys whose resistance is nearly independent of temperature. They are used to make precision resistors and measuring equipment.

---

## Superconductivity

In 1911, Heike Kamerlingh Onnes cooled mercury to about 4K (four degrees above absolute zero) and made an astonishing discovery: the resistance dropped to exactly zero. Not almost zero — zero.

This phenomenon is called **superconductivity**. Below a material-specific **critical temperature** T_c, the material becomes a perfect conductor.

### Why Superconductors Have Zero Resistance

In a superconductor, electrons pair up into **Cooper pairs**. These pairs behave quantum mechanically in a way that lets them flow through the crystal lattice without collisions — they pass right through without interacting with the ions. There is nothing to impede them, so there is no resistance.

```
RESISTANCE vs TEMPERATURE for a SUPERCONDUCTOR:

Resistance (R)
  ^
  |
  |     Metal behavior (linear increase)
  |    /
  |   /
  |  /
  | /
  |/
  +---------+--------> Temperature
            ^
        T_c (critical
         temperature)
         
At T_c, resistance drops suddenly to ZERO.
Below T_c, R = 0.0000... (not approximately zero — EXACTLY zero)
```

### Practical Applications of Superconductors

- **MRI machines**: Superconducting coils create the powerful magnetic field needed to image the body
- **Maglev trains**: Superconducting magnets allow trains to levitate with no friction
- **Particle accelerators**: CERN's Large Hadron Collider uses superconducting magnets to guide protons
- **Power transmission**: Superconducting cables lose no energy to heating

### The Challenge: They Need to Be Very Cold

Traditional superconductors require liquid helium cooling (4K = -269°C), which is expensive. Researchers are working on **high-temperature superconductors** — materials that superconduct at higher temperatures. The current record is around 250K (-23°C), still very cold, but a huge improvement. Room-temperature superconductivity remains one of the most sought-after goals in physics.

---

## Practical Resistance in Circuits

### Resistor Color Codes

Commercial resistors are labeled with colored bands that encode their resistance value:

```
COLOR CODE:
Black  = 0      Green  = 5
Brown  = 1      Blue   = 6
Red    = 2      Violet = 7
Orange = 3      Grey   = 8
Yellow = 4      White  = 9

READING A 4-BAND RESISTOR:
  Band 1: First digit
  Band 2: Second digit
  Band 3: Multiplier (10^n)
  Band 4: Tolerance (Gold = ±5%, Silver = ±10%)

Example: Red-Violet-Orange-Gold
  Red = 2, Violet = 7, Orange = 10^3 multiplier
  R = 27 × 10^3 = 27,000 Ω = 27 kΩ ±5%
```

### Fixed vs Variable Resistors

**Fixed resistors**: Constant resistance. Used to set current levels, voltage dividers.

**Variable resistors (potentiometers/rheostats)**: Resistance can be changed by turning a knob or sliding a contact. Used in volume controls, dimmer switches.

```
RHEOSTAT SYMBOL:
    ---/\/\/--->
           ^
           | (arrow = adjustable contact)
```

---

## Worked Examples

### Worked Example 1: Using Ohm's Law to Find Voltage

**Problem:** A resistor of 470 Ω has a current of 15 mA flowing through it. What is the voltage across it?

**Solution:**

```
Given:
  R = 470 Ω
  I = 15 mA = 15 × 10^-3 A = 0.015 A

Formula:
  V = I × R

Calculation:
  V = 0.015 × 470
  V = 7.05 V

Answer: The voltage across the resistor is 7.05 V
```

### Worked Example 2: Calculating Resistance from V and I

**Problem:** A lamp has a voltage of 12 V across it and a current of 0.5 A through it. What is its resistance?

**Solution:**

```
Given:
  V = 12 V
  I = 0.5 A

Formula:
  R = V / I

Calculation:
  R = 12 / 0.5
  R = 24 Ω

Answer: The lamp's resistance is 24 Ω
```

### Worked Example 3: Wire Resistance Using Resistivity

**Problem:** A copper wire is 5 m long and has a circular cross-section with diameter 1.5 mm. Calculate its resistance.
(Resistivity of copper: ρ = 1.72 × 10^-8 Ω·m)

**Solution:**

```
Given:
  ρ = 1.72 × 10^-8 Ω·m
  L = 5 m
  d = 1.5 mm = 1.5 × 10^-3 m
  r = d/2 = 0.75 × 10^-3 m

Step 1: Calculate cross-sectional area
  A = π × r²
  A = π × (0.75 × 10^-3)²
  A = π × 5.625 × 10^-7
  A = 1.767 × 10^-6 m²

Step 2: Apply resistivity formula
  R = ρL / A
  R = (1.72 × 10^-8 × 5) / (1.767 × 10^-6)
  R = (8.6 × 10^-8) / (1.767 × 10^-6)
  R = 0.0487 Ω

Answer: The wire has a resistance of approximately 0.049 Ω (very low, as expected for copper)
```

### Worked Example 4: Temperature Effect on Resistance

**Problem:** A copper wire has resistance 20 Ω at 20°C. What is its resistance at 80°C?
(Temperature coefficient of copper: α = 0.0039 /°C)

**Solution:**

```
Given:
  R_0 = 20 Ω (at T_0 = 20°C)
  T = 80°C
  α = 0.0039 /°C

Step 1: Find temperature change
  ΔT = T - T_0 = 80 - 20 = 60°C

Step 2: Apply temperature formula
  R_T = R_0 × (1 + α × ΔT)
  R_T = 20 × (1 + 0.0039 × 60)
  R_T = 20 × (1 + 0.234)
  R_T = 20 × 1.234
  R_T = 24.68 Ω

Answer: The resistance at 80°C is approximately 24.7 Ω
        (an increase of about 23% for a 60°C temperature rise)
```

This explains why resistance thermometers work — they measure temperature by measuring resistance.

---

## Summary

- **Resistance** R is the opposition to current flow, measured in Ohms (Ω)
- **Ohm's Law**: V = IR — voltage, current, and resistance are related; current is proportional to voltage for ohmic materials
- Resistance arises microscopically from **electrons colliding with lattice ions** as they drift through the metal
- **Ohmic materials** (metal resistors at constant temperature) have constant R; their I-V graph is a straight line through the origin
- **Non-ohmic materials** have R that varies with conditions: filament bulbs (R increases as it heats up), diodes (very asymmetric), thermistors (R changes with temperature)
- **Resistivity formula** R = ρL/A: resistance scales with length L, inversely with area A, and depends on the material's resistivity ρ (Ω·m)
- In **metals**, resistance increases with temperature (more lattice vibration); in **semiconductors**, resistance decreases with temperature (more charge carriers freed)
- **Superconductors** have exactly zero resistance below their critical temperature — electrons form Cooper pairs that travel with no collisions

---

## Key Equations

```
Ohm's Law:
    V = I × R
    I = V / R
    R = V / I

Units: 1 Ω = 1 V/A

Resistivity:
    R = ρ × L / A

    ρ = resistivity (Ω·m)
    L = length (m)
    A = cross-sectional area (m²)

Temperature Coefficient:
    R_T = R_0 × (1 + α × ΔT)

    R_0 = resistance at reference temperature
    α   = temperature coefficient (/°C)
    ΔT  = temperature change (°C)

Key Values:
    Copper resistivity:         ρ = 1.72 × 10^-8 Ω·m
    Copper temp coefficient:    α = 3.9 × 10^-3 /°C
    Silicon resistivity:        ρ ≈ 6.4 × 10^2 Ω·m
```
