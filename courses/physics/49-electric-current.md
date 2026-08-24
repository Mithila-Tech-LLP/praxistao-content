# Chapter 49: Electric Current

> **"Electricity is not magic — it is simply charges in motion. Once you understand that, you hold the key to the modern world."**

---

## Table of Contents

1. [What Is Electric Current?](#what-is-electric-current)
2. [Conventional Current vs Electron Flow](#conventional-current-vs-electron-flow)
3. [Drift Velocity vs Signal Speed](#drift-velocity-vs-signal-speed)
4. [Measuring Current: The Ampere](#measuring-current-the-ampere)
5. [Direct Current vs Alternating Current](#direct-current-vs-alternating-current)
6. [Current Density](#current-density)
7. [Charge Carriers in Different Materials](#charge-carriers-in-different-materials)
8. [How Batteries Produce Current](#how-batteries-produce-current)
9. [Circuit Symbols](#circuit-symbols)
10. [Worked Examples](#worked-examples)
11. [Summary](#summary)
12. [Key Equations](#key-equations)

---

## What Is Electric Current?

Every material around you — the phone in your hand, the wire in the wall, the water in a glass — is made of atoms. At the center of each atom is a nucleus surrounded by electrons. In most materials those electrons are locked tightly in place. But in some materials, especially **metals**, the outermost electrons are only loosely bound to their atoms. They can break free and wander randomly through the material like gas molecules bouncing around.

These free electrons are what make electricity possible.

**Electric current** is the net flow of electric charge through a cross-section of a conductor per unit time. When we say current is flowing through a wire, we mean that charge is moving in a preferred direction rather than just wandering randomly.

Think of it this way: even in a copper wire with no battery attached, electrons are constantly moving — but randomly, in all directions. The net charge crossing any cross-section of the wire per second is zero. The moment you connect a battery, it creates an electric field throughout the wire. Now the electrons still bounce around randomly, but they also drift very slowly in one direction. That tiny drift is what we call **electric current**.

### The Mathematical Definition

If a charge Q passes through a cross-section of a conductor in time t, the average current I is:

```
        Q
I  =  -----
        t
```

Where:
- I = current (measured in Amperes, A)
- Q = charge (measured in Coulombs, C)
- t = time (measured in seconds, s)

One Ampere means one Coulomb of charge passes through a cross-section every second.

```
1 A = 1 C/s
```

This is actually a very large amount of charge. One Coulomb contains about 6.24 × 10^18 electrons.

---

## Conventional Current vs Electron Flow

Here is one of the most confusing things in all of introductory physics: conventional current flows in the **opposite direction** to the actual movement of electrons.

How did this happen? In the 1700s, Benjamin Franklin decided to describe electricity as a fluid that flowed from positive to negative. He had no way to know about electrons — those were not discovered until 1897 by J.J. Thomson. By then, the convention of current flowing from positive to negative was so deeply embedded in science and engineering that it was kept.

When it was later discovered that in metals it is actually negatively charged electrons that carry the current, and that they flow from negative to positive, physics had two choices: change all existing textbooks and engineering diagrams, or keep the convention. The convention was kept.

### The Direction Diagram

```
CONVENTIONAL CURRENT DIRECTION (what we draw in circuits):
                         I  --->
    +  ==========================================  -
  (positive terminal)                         (negative terminal)
  BATTERY                                     BATTERY


ACTUAL ELECTRON MOVEMENT (what really happens):
                         e- <---
    +  ==========================================  -
  (positive terminal)                         (negative terminal)
  BATTERY                                     BATTERY
```

Electrons are negatively charged. They are attracted toward the positive terminal of the battery. So they flow from the negative terminal, through the external circuit, toward the positive terminal.

Conventional current is defined as flowing from the positive terminal, through the external circuit, toward the negative terminal — the opposite direction.

**Does it matter?** For most circuit analysis, it does not matter at all. You can analyze circuits using conventional current and get the right answers for current magnitude, voltage, and power. Where it matters is in understanding semiconductors and specific electronic devices like diodes where the actual charge carriers make a difference.

### Why Keep a Wrong Convention?

It is worth making peace with this. Consider: we say the Sun rises in the East, even though we know the Earth is rotating. We use relative language because it is useful. Similarly, conventional current is useful because it describes how circuits behave in a way that is consistent and self-contained.

Every circuit diagram in engineering uses conventional current direction. Every ammeter reads conventional current. Once you know the convention, it works perfectly.

---

## Drift Velocity vs Signal Speed

One of the most surprising facts about electric current: the electrons in a wire move incredibly slowly.

### Drift Velocity

**Drift velocity** (v_d) is the average velocity of charge carriers (electrons in a metal) in the direction of current flow.

In a typical household wire carrying 1 Ampere:
- Drift velocity of electrons ≈ 0.1 mm/s (one tenth of a millimeter per second)

At that speed, an electron would take about 3 hours to travel 1 meter along the wire!

Yet when you flip a light switch, the bulb turns on almost instantly. How?

### The Pipe Analogy

Imagine a long pipe completely filled with water. If you push more water in one end, water immediately comes out the other end. The water already in the pipe transmits the pressure almost instantly, even though no individual water molecule travels the whole length of the pipe quickly.

```
PIPE ANALOGY:
                   Push here
                      |
    ====|=============v============|====
    ====|=O=O=O=O=O=O=O=O=O=O=O=O=|====  --> Water comes out here
    ====|============================|====
    ^
    Piston
    
    Each O = water molecule. You push one in, one comes out the other end.
    The individual molecules barely move, but the effect travels instantly.
```

A wire filled with electrons works the same way. When the battery creates an electric field, it propagates through the wire at nearly the speed of light. This field pushes all the electrons simultaneously. Each electron only moves a tiny bit, but the net effect — current flowing in the circuit — happens essentially instantaneously.

**Signal speed**: The speed at which the electrical effect (the electromagnetic field) propagates through the wire: approximately 2 × 10^8 m/s (about two-thirds the speed of light in vacuum).

**Drift velocity**: The speed at which individual electrons actually move: approximately 0.1 mm/s.

```
SUMMARY OF SPEEDS:
Speed of light in vacuum:          3.0 × 10^8 m/s
Signal speed in copper wire:       ~2.0 × 10^8 m/s  (fast — why lights turn on instantly)
Speed of sound in air:             343 m/s           (much slower)
Walking speed:                     ~1.4 m/s          (very slow)
Electron drift velocity in wire:   ~0.0001 m/s       (incredibly slow)
```

---

## Measuring Current: The Ampere

The **Ampere** (symbol: A) is the SI unit of electric current. It is named after André-Marie Ampère, a French physicist and mathematician who did foundational work on electromagnetism in the early 19th century.

The formal definition of the Ampere was revised in 2019 as part of the SI redefinition. The modern definition fixes the value of the elementary charge e:

```
e = 1.602176634 × 10^-19 C (exactly, by definition)
```

One Ampere is therefore the current corresponding to the flow of approximately 6.24 × 10^18 elementary charges per second.

### Common Current Values

```
Device                          Typical Current
-------------------------------------------------
Small LED                       20 mA = 0.020 A
Phone charging                  0.5 - 2 A
Laptop charging                 3 - 5 A
Household light bulb            0.5 A (60W at 120V)
Electric kettle                 8 - 10 A
Car starter motor               100 - 200 A
Lightning bolt                  ~30,000 A (briefly!)
Nerve signal in your brain      ~1 nanoampere = 10^-9 A
```

### Measuring with an Ammeter

An **ammeter** measures current. It must be placed **in series** with the component — the current must actually flow through it. An ideal ammeter has zero resistance so it doesn't affect the circuit.

```
CORRECT: Ammeter in series
         +----[A]----[Bulb]----+
         |                    |
        [+]     Battery      [-]
         |                    |
         +--------------------+
         
WRONG: Ammeter in parallel (would short-circuit the bulb)
```

---

## Direct Current vs Alternating Current

### Direct Current (DC)

**Direct current** (DC) is current that flows in one direction only. The magnitude may vary, but the direction does not change. Batteries produce DC. Most electronic devices (phones, computers, LEDs) run on DC internally.

```
DC CURRENT vs TIME:
Current (A)
|
2 |------------------------------------
  |
1 |
  |
0 +-------------------------------------> Time
```

### Alternating Current (AC)

**Alternating current** (AC) is current that periodically reverses direction. It follows a sinusoidal pattern. The power grid in most countries delivers AC.

```
AC CURRENT vs TIME:
Current (A)
|      /\          /\
2 |    /  \        /  \
  |   /    \      /    \
0 |--/------\----/------\-->  Time
  |          \  /        \  /
-2|           \/          \/
  |<-- Period T -->|
```

**Frequency** is the number of complete cycles per second, measured in Hertz (Hz).

```
Europe, Asia, Africa, Australia:  50 Hz (50 cycles per second)
North America, some of Asia:      60 Hz (60 cycles per second)
```

At 50 Hz, the current direction reverses 100 times per second (it completes 50 full back-and-forth cycles).

### Why AC for Power Grids?

AC has a crucial practical advantage: it is easy to change its voltage using a transformer. High-voltage AC transmission (400,000 V in the UK) loses much less energy to heat in power lines than low-voltage transmission. Near your home, transformers step it down to safe levels (240V in Europe, 120V in the USA).

Nikola Tesla championed AC, while Thomas Edison promoted DC. The "War of Currents" in the 1880s was eventually won by AC for long-distance power distribution. However, modern high-voltage DC (HVDC) transmission is making a comeback for very long distances and undersea cables.

### Converting Between AC and DC

Most modern electronics use AC from the wall socket but require DC internally. The conversion is done by the power supply unit:

```
AC from wall --> Rectifier --> Filter --> Regulator --> DC output
(230V AC)    (flips -ve half)  (smooths)  (stabilizes)  (5V DC)
```

The charger or adapter you use for your phone is performing this conversion.

---

## Current Density

**Current density** J is the current per unit cross-sectional area of a conductor. It tells you how much current is packed into each square meter of the wire.

```
        I
J  =  -----
        A
```

Where:
- J = current density (A/m²)
- I = current (A)
- A = cross-sectional area (m²)

Current density matters because it determines how much heating occurs in the wire. Pack too much current into a thin wire and it overheats — this is the principle behind fuses and circuit breakers, which protect your home from electrical fires.

### Worked Example: Current Density in a Wire

A copper wire with a circular cross-section has a diameter of 2 mm and carries a current of 4 A. What is the current density?

```
Step 1: Find the cross-sectional area
   Radius r = diameter/2 = 2mm/2 = 1mm = 1 × 10^-3 m
   A = π × r² = π × (1 × 10^-3)² = π × 10^-6 ≈ 3.14 × 10^-6 m²

Step 2: Calculate current density
   J = I/A = 4 / (3.14 × 10^-6) ≈ 1.27 × 10^6 A/m²
```

This is a very high current density. In practice, copper wires are rated for maximum currents based on their diameter to prevent overheating.

---

## Charge Carriers in Different Materials

The electrons in metals are not the only way to carry electric current. Different materials have different charge carriers.

### Metals: Free Electrons

In metals like copper, silver, and aluminum, the outermost electrons (one or two per atom) are loosely bound. They form what is called a "sea" of free electrons that permeates the metal lattice. These are the **conduction electrons**.

```
METAL STRUCTURE:
    Cu²⁺  Cu²⁺  Cu²⁺  Cu²⁺
       e⁻     e⁻     e⁻
    Cu²⁺  Cu²⁺  Cu²⁺  Cu²⁺
       e⁻     e⁻     e⁻     e⁻
    Cu²⁺  Cu²⁺  Cu²⁺  Cu²⁺

    Cu²⁺ = copper ion (fixed in lattice)
    e⁻   = free electron (wanders through lattice)
```

When a battery is connected, these electrons drift slowly in one direction while still bouncing randomly off the copper ions.

### Electrolytes: Positive and Negative Ions

In solutions of salts, acids, or bases (**electrolytes**), the charge carriers are ions — both positive and negative. When you dissolve table salt (NaCl) in water, it splits into Na+ and Cl- ions.

```
IN ELECTROLYTE (e.g., saltwater):

Positive ions (Na+):  +  +  +     move toward negative electrode  <---
Negative ions (Cl-):  -  -  -     move toward positive electrode  --->

CONVENTIONAL CURRENT direction is NET result of both movements.
```

This is how car batteries work internally: sulfuric acid solution with lead sulfate ions carries charge between the lead plates.

### Semiconductors: Electrons AND Holes

In **semiconductors** like silicon, something remarkable happens. When an electron gains enough energy to break free and become a conduction electron, it leaves behind an empty space in the crystal lattice. This empty space behaves like a positive charge moving through the material — it is called a **hole**.

```
SEMICONDUCTOR:
Before:  Si-Si-Si-Si-Si  (all bonds occupied)
After:   Si-Si-[  ]-Si-Si  (electron escaped, hole remains)
                  ^
                  hole = missing electron = acts like +ve charge
```

In a semiconductor, current is carried by both electrons (moving opposite to conventional current) and holes (moving in the same direction as conventional current). This is fundamental to how transistors, diodes, and all modern electronics work.

### Plasmas: Ionized Gas

In a **plasma** — the fourth state of matter found in lightning, neon signs, and the Sun — atoms are stripped of electrons. The result is a soup of free electrons and positive ions, both acting as charge carriers.

### Summary Table of Charge Carriers

```
Material          Charge Carrier(s)           Example
-----------------------------------------------------------
Metals            Free electrons (negative)   Copper wire
Electrolytes      Positive and negative ions  Battery acid
Semiconductors    Electrons AND holes         Silicon chip
Plasma            Electrons + positive ions   Lightning
```

---

## How Batteries Produce Current

A **battery** converts **chemical energy** into **electrical energy**. Understanding this at a basic level demystifies where "electricity comes from."

### The Electrochemical Cell

The simplest battery is a voltaic cell (named after Alessandro Volta who invented it in 1800). It consists of:
1. Two different metals (**electrodes**) — typically zinc and copper
2. An **electrolyte** — a solution that conducts ions (typically sulfuric acid or saltwater)

```
SIMPLE VOLTAIC CELL:

      Zinc (Zn)          Copper (Cu)
      electrode          electrode
          |                  |
          |    Electrolyte   |
          |    (acid/salt    |
          |     solution)    |
    Zn²⁺ goes           Cu²⁺ deposits
    into solution         onto copper
    (oxidation)          (reduction)
          |                  |
          |   e⁻ flow        |
          +<-----------------+
              External
              circuit
```

### The Process Step by Step

1. **At the zinc electrode (negative terminal / anode):**
   Zinc atoms lose electrons and become Zn²⁺ ions that dissolve into the electrolyte.
   ```
   Zn --> Zn²⁺ + 2e⁻
   ```
   Electrons are left behind on the zinc electrode, making it negatively charged.

2. **At the copper electrode (positive terminal / cathode):**
   Copper ions from the solution gain electrons and deposit as solid copper.
   ```
   Cu²⁺ + 2e⁻ --> Cu
   ```
   Electrons are consumed here, so this electrode becomes positively charged relative to the zinc.

3. **The external circuit:**
   The difference in charge creates a **potential difference (voltage)** between the terminals. When a wire connects the terminals, electrons flow through it from zinc (–) to copper (+). This is the electric current.

4. **The electrolyte:**
   Ions move through the solution to balance the charge build-up internally, completing the circuit.

### EMF (Electromotive Force)

The **EMF** (Electromotive Force) of a battery, given the symbol ε (epsilon), is the energy given to each coulomb of charge as it passes through the battery. It is measured in Volts.

```
EMF is NOT a force! It is energy per charge.
ε = Energy supplied / Charge = E/Q

Units: Joules per Coulomb = Volts (V)
```

A fresh AA battery has an EMF of 1.5V. A car battery has 12V. A phone charger adapter outputs around 5V.

---

## Circuit Symbols

Before drawing or reading any circuit, you need to know the standard symbols. These are international conventions used in all physics and engineering.

```
STANDARD CIRCUIT SYMBOLS:

Battery (single cell):
    ---|+  -|---
    long line = positive terminal
    short line = negative terminal

Battery (multiple cells):
    ---|+  -|---|+  -|---

Resistor:
    ---[  R  ]---
    or
    ---/\/\/\/---

Wire (conductor):
    ------------

Junction (wires connected):
         |
    -----+-----
         |

Wires crossing (NOT connected):
         |
    -----+-----   (sometimes shown with a bridge: ---)
         |

Switch (open):
    ---  /  ---

Switch (closed):
    ---/----

Lamp/Bulb:
    ---(@)---

Voltmeter (in parallel):
    ---|V|---
    (high resistance device)

Ammeter (in series):
    ---|A|---
    (low resistance device)

Capacitor:
    ---| |---

Diode:
    --->|---
    (arrow shows direction of conventional current)

LED (Light Emitting Diode):
    --->|---> (with light arrows)

Earth/Ground:
    ---
      =
       =
        =
```

### How to Draw Circuit Diagrams

1. Draw wires as straight lines with right-angle corners
2. Use standard symbols at every component
3. Label each component (R₁, R₂, etc. and their values)
4. Indicate the current direction with arrows
5. Label voltages if known

---

## Worked Examples

### Worked Example 1: Calculating Current from Charge and Time

**Problem:** A charge of 180 C flows through a light bulb in 2 minutes. What is the current?

**Solution:**

```
Given:
  Q = 180 C
  t = 2 min = 2 × 60 = 120 s

Formula:
  I = Q/t

Calculation:
  I = 180 / 120
  I = 1.5 A

Answer: The current is 1.5 A
```

### Worked Example 2: How Many Electrons?

**Problem:** A current of 2 A flows through a wire for 30 seconds. How many electrons pass through a cross-section of the wire?

**Solution:**

```
Given:
  I = 2 A
  t = 30 s
  Charge of one electron = 1.6 × 10^-19 C

Step 1: Find total charge
  Q = I × t
  Q = 2 × 30
  Q = 60 C

Step 2: Find number of electrons
  Number of electrons = Q / (charge per electron)
  n = 60 / (1.6 × 10^-19)
  n = 3.75 × 10^20 electrons

Answer: 3.75 × 10^20 electrons (375 quintillion electrons!)
```

This shows just how tiny the electron's charge is. Even "small" currents involve astronomical numbers of electrons.

### Worked Example 3: Current Density Comparison

**Problem:** Two wires carry the same current of 5 A. Wire A has diameter 1 mm, Wire B has diameter 2 mm. Compare their current densities.

**Solution:**

```
Wire A:
  r_A = 0.5 mm = 5 × 10^-4 m
  A_A = π × (5 × 10^-4)² = π × 2.5 × 10^-7 = 7.85 × 10^-7 m²
  J_A = I / A_A = 5 / (7.85 × 10^-7) = 6.37 × 10^6 A/m²

Wire B:
  r_B = 1 mm = 1 × 10^-3 m
  A_B = π × (1 × 10^-3)² = π × 10^-6 = 3.14 × 10^-6 m²
  J_B = I / A_B = 5 / (3.14 × 10^-6) = 1.59 × 10^6 A/m²

Ratio: J_A / J_B = 6.37 / 1.59 ≈ 4

Answer: Wire A has 4 times the current density of Wire B.
        (Doubling the diameter quadruples the area, quartering the current density)
```

This is why thicker wires are used for high-current applications — they spread the current over more area, reducing heating.

---

## Summary

- **Electric current** I is the rate of flow of electric charge: I = Q/t, measured in Amperes (A)
- **Conventional current** flows from positive to negative terminal (outside the battery), but **actual electrons** flow from negative to positive — the directions are opposite due to historical convention
- **Drift velocity** of electrons in a wire is extremely slow (~0.1 mm/s), but the **signal speed** (electromagnetic field propagation) is close to the speed of light — which is why lights turn on instantly
- In **DC** (direct current), charge flows in one direction only; in **AC** (alternating current), the direction reverses periodically (50 Hz in Europe, 60 Hz in North America)
- **Current density** J = I/A is the current per unit cross-sectional area; higher current density means more heating
- Charge carriers vary by material: **free electrons** in metals, **ions** in electrolytes, **electrons and holes** in semiconductors, **both** in plasmas
- A **battery** converts chemical energy to electrical energy via electrochemical reactions; its **EMF** is the energy given per coulomb of charge
- Standard circuit symbols must be memorized to read and draw circuit diagrams

---

## Key Equations

```
Current:
    I = Q / t
    
    I = current (A)
    Q = charge (C)
    t = time (s)

Charge:
    Q = I × t

Number of electrons:
    n = Q / e
    
    e = 1.6 × 10^-19 C (charge of one electron)

Current Density:
    J = I / A
    
    J = current density (A/m²)
    A = cross-sectional area (m²)

Units Summary:
    1 Ampere = 1 Coulomb per second
    1 A = 1 C/s
```
