# Chapter 51: DC Circuits

> **"A circuit is like a water system — current flows like water, voltage is the pressure, and resistance is the pipe's narrowness. Once you see the analogy, circuits click."**

---

## Table of Contents

1. [Building a Circuit: The Basics](#building-a-circuit-the-basics)
2. [Circuit Diagram Symbols](#circuit-diagram-symbols)
3. [Series Circuits](#series-circuits)
4. [Parallel Circuits](#parallel-circuits)
5. [Mixed Series-Parallel Circuits](#mixed-series-parallel-circuits)
6. [Internal Resistance of a Battery](#internal-resistance-of-a-battery)
7. [Short Circuits and Open Circuits](#short-circuits-and-open-circuits)
8. [Voltage Dividers](#voltage-dividers)
9. [Why House Wiring Uses Parallel](#why-house-wiring-uses-parallel)
10. [Worked Examples](#worked-examples)
11. [Summary](#summary)
12. [Key Equations](#key-equations)

---

## Building a Circuit: The Basics

Every DC circuit is built from the same fundamental elements:

1. A **source of EMF** (battery, power supply) — provides the energy to drive current
2. **Conductors** (wires) — carry the current with negligible resistance
3. **Load components** (resistors, bulbs, motors) — consume electrical energy
4. Optionally: **switches**, **meters**, and **protection devices** (fuses)

A circuit must form a complete, unbroken loop for current to flow. This is one of the most important concepts in all of electronics: **current needs a complete path**.

```
COMPLETE CIRCUIT (current flows):

         +----[ R ]----+
         |             |
        [+]           [-]
         |  BATTERY   |
         +------------+
         
         Electrons flow from - through R back to +
         Conventional current flows from + through R back to -


INCOMPLETE CIRCUIT (no current flows):

         +----[ R ]----
         |             (broken wire - no path!)
        [+]           
         |  BATTERY   
         +------------+
         
         No complete loop = no current
```

---

## Circuit Diagram Symbols

When drawing circuits, we use standardized symbols. Here is the complete set you need to know:

```
BATTERY (single cell):
        |  |
    ----+  +----     (long line = positive, short line = negative)
    
BATTERY (multiple cells):
        |  |  |  |
    ----+  +--+  +----

RESISTOR:
    ---[===]---
    or  ---/\/\/---

VARIABLE RESISTOR (rheostat):
    ---[===]-->---
    
LAMP / BULB:
    ---(X)---

SWITCH (open):
         /
    ---o  o---

SWITCH (closed):
    ---o--o---

AMMETER (measures current, in series):
    ---[A]---

VOLTMETER (measures voltage, in parallel):
    ---[V]---
    (with connecting wire going to both sides of component)

WIRE JUNCTION (wires connected):
        |
    ----●----
        |

WIRES CROSSING (NOT connected):
        |
    ----+----    or    ---+---
        |                  |
                      (no dot = no connection)

FUSE:
    ---[===]---   (with zig-zag or straight line inside box)

LED:
    ---[>|]---   (arrow shows conventional current direction)

CAPACITOR:
    ---| |---

EARTH / GROUND:
       ---
        -
         -
```

### How to Draw a Good Circuit Diagram

Follow these conventions when drawing circuits:

1. Draw wires as straight horizontal or vertical lines only (no diagonals)
2. Place components along straight sections
3. Make right-angle corners
4. Use dots to show junctions (connected wires)
5. Label all components with their symbol (R₁, R₂) and value
6. Indicate the positive terminal of the battery with a +

---

## Series Circuits

In a **series circuit**, components are connected one after another in a single path. There is only one route for the current to follow.

```
SERIES CIRCUIT:

    +----[R₁]----[R₂]----[R₃]----+
    |                             |
   [+]         BATTERY           [-]
    |                             |
    +-----------------------------+
    
    Only ONE path for current.
    Same current I passes through every component.
```

### Rules for Series Circuits

**Rule 1: Current is the same everywhere**

Because there is only one path, every charge that flows through R₁ must also flow through R₂ and R₃. No charge accumulates or gets lost.

```
I_total = I₁ = I₂ = I₃ = ... (all equal)
```

**Rule 2: Voltages add up**

Each resistor has a voltage drop across it (V = IR). These drops must add up to the total battery voltage, because all the energy given by the battery gets used up across the resistors.

```
V_total = V₁ + V₂ + V₃ + ...
```

**Rule 3: Resistances add up**

Since R_total = V_total / I, and V_total = V₁ + V₂ + V₃ = IR₁ + IR₂ + IR₃ = I(R₁ + R₂ + R₃):

```
R_total = R₁ + R₂ + R₃ + ...
```

Adding more resistors in series always increases the total resistance.

### Series Circuit: Voltage Distribution

The voltage across each resistor is proportional to its resistance. A larger resistor gets a larger share of the battery voltage.

```
VOLTAGE DIVIDER IN SERIES:
     
    EMF = 12V
    R₁ = 2Ω  -->  V₁ = ?
    R₂ = 4Ω  -->  V₂ = ?

    R_total = 2 + 4 = 6Ω
    I = 12/6 = 2A

    V₁ = I × R₁ = 2 × 2 = 4V
    V₂ = I × R₂ = 2 × 4 = 8V
    Check: V₁ + V₂ = 4 + 8 = 12V ✓
```

Notice that R₂ is twice R₁, so V₂ is twice V₁. This proportional sharing of voltage is what makes voltage dividers work.

### Practical Applications of Series Circuits

- **String of old Christmas lights** (if one bulb goes out, they all go out — because there's only one path)
- **Fuses and circuit breakers** in series with appliances
- Battery cells connected in series to increase voltage (two 1.5V cells in series = 3V)

---

## Parallel Circuits

In a **parallel circuit**, components are connected across the same two nodes, providing multiple paths for current to flow.

```
PARALLEL CIRCUIT:

    +--------+--------+--------+
    |        |        |        |
   [+]      [R₁]     [R₂]    [R₃]
    |        |        |        |
   [-]      [-]      [-]      [-]
    |        |        |        |
    +--------+--------+--------+
    
    MULTIPLE PATHS for current.
    Same voltage across every component.
```

### Rules for Parallel Circuits

**Rule 1: Voltage is the same across all branches**

All components connected in parallel share the same two nodes. By definition, the voltage between those two nodes is the same for every component.

```
V_total = V₁ = V₂ = V₃ = ... (all equal)
```

**Rule 2: Currents add up**

Current from the battery splits into branches. All branch currents add up to the total current drawn from the battery.

```
I_total = I₁ + I₂ + I₃ + ...
```

**Rule 3: Reciprocal resistance formula**

Since I_total = I₁ + I₂ + I₃ = V/R₁ + V/R₂ + V/R₃ = V(1/R₁ + 1/R₂ + 1/R₃):

```
    1          1      1      1
-------  =  ----- + ----- + ----- + ...
R_total      R₁     R₂     R₃
```

Adding more resistors in parallel always **decreases** the total resistance (more paths = less total opposition).

### Special Case: Two Resistors in Parallel

For exactly two resistors in parallel, the formula simplifies to:

```
         R₁ × R₂
R_total = --------
          R₁ + R₂
         
("Product over sum" — only valid for exactly TWO resistors)
```

### Special Case: Identical Resistors in Parallel

If n identical resistors each of value R are in parallel:

```
        R
R_total = ---
        n
```

For example: three 30Ω resistors in parallel = 30/3 = 10Ω

### Practical Applications of Parallel Circuits

- **Household wiring** (every socket at the same voltage, devices work independently)
- **Battery cells in parallel** (increases capacity/current, not voltage)
- **Computer circuit boards** (components at same supply voltage)

---

## Mixed Series-Parallel Circuits

Most real circuits combine series and parallel connections. The approach to solve these is: **simplify step by step**, replacing groups of components with their equivalent resistance.

```
GENERAL STRATEGY:
1. Identify groups of components that are purely series or purely parallel
2. Replace each group with its single equivalent resistance
3. Repeat until you have one equivalent resistance
4. Use Ohm's Law to find total current and voltage
5. Work backwards to find current/voltage in individual components
```

### Example Mixed Circuit

```
CIRCUIT:
              R₂ = 6Ω
    +---[R₁]---+---+---+
    |    4Ω   |   |    |
   [+]        |  [R₂]  [R₃]
    |         |   |   8Ω|
   [-]        +---+---+  
    |                  |
    +----- 12V --------+
    
    R₂ and R₃ are in parallel.
    Their combination is in series with R₁.
```

**Step 1: Simplify the parallel combination (R₂ parallel R₃)**

```
        R₂ × R₃     6 × 8     48
R_23 = --------- = ------- = ---- = 3.43Ω  ... wait let me redo this cleanly:
        R₂ + R₃     6 + 8     14
       
Actually: 48/14 = 24/7 ≈ 3.43 Ω

Let's use R₂ = 6Ω and R₃ = 12Ω for cleaner numbers:

        6 × 12      72
R_23 = -------- = ----- = 4 Ω
        6 + 12      18
```

The circuit is now just R₁ = 4Ω in series with R_23 = 4Ω, with 12V supply.

**Step 2: Find total resistance**

```
R_total = R₁ + R_23 = 4 + 4 = 8 Ω
```

**Step 3: Find total current**

```
I_total = V / R_total = 12 / 8 = 1.5 A
```

---

## Internal Resistance of a Battery

A real battery is not a perfect voltage source. It has its own internal resistance r due to the chemical materials inside and the connections.

```
MODEL OF A REAL BATTERY:

    |---[  r  ]---[+ EMF -]---|
    |                         |
    External circuit          
    
    r = internal resistance
    EMF = electromotive force (ideal voltage)
```

When current I flows from the battery:
- Voltage is lost across the internal resistance: V_lost = I × r
- The voltage actually available to the external circuit (**terminal voltage**) is less than the EMF

```
Terminal Voltage = EMF - (I × r)

V_terminal = ε - I × r

Where:
  ε   = EMF (Volts)
  I   = current (Amperes)
  r   = internal resistance (Ohms)
```

### What Happens as Current Increases

```
TERMINAL VOLTAGE vs CURRENT:

Voltage (V)
  ^
  |
ε |----___
  |       ---___
  |             ---___    (linear decrease)
  |                   ---___
  |                         ---
  +-----------------------------> Current (I)
  0
  
  When I = 0 (open circuit):  V_terminal = ε  (maximum voltage)
  When I is large:            V_terminal drops significantly
  When R_external = 0 (short): I = ε/r (very large current)
```

### Why Does a Battery "Die"?

As a battery discharges, the chemical reactants are consumed. This causes the internal resistance r to increase over time. Eventually, even though the battery may still have some EMF, the internal resistance drop becomes so large that the terminal voltage is too low to power devices.

A "dead" AA battery may still measure 1.4V with a voltmeter (no current, so no voltage drop), but when you put it in a device that draws significant current, the terminal voltage collapses.

### Measuring Internal Resistance

To measure a battery's internal resistance:
1. Measure the open-circuit EMF (V when no current flows)
2. Connect a known external resistance R and measure terminal voltage V
3. Calculate: r = (ε - V) / I = (ε - V) × R / V

---

## Short Circuits and Open Circuits

### Open Circuit

An **open circuit** is a break in the circuit — a gap where the conducting path is interrupted. No current can flow.

```
OPEN CIRCUIT:
    +----[R₁]----  X  ----+    (X = break in wire)
    |                     |
   BATTERY                
   
   No current flows: I = 0
   Voltage appears across the gap
```

If you measure voltage across an open circuit with a voltmeter, you will read the full supply voltage at the break point (all the EMF "piles up" at the gap).

### Short Circuit

A **short circuit** is an unintended low-resistance path that bypasses a component.

```
SHORT CIRCUIT:
    +----[R]----+
    |           |
   BATTERY  [R is bypassed by 
    |         direct wire connection]
    |           |
    +-----------+
    
    Nearly all current bypasses R.
    Very large current flows through short.
```

Short circuits are dangerous because:
1. Very large current flows (I = V/R with tiny R)
2. I²R heating is enormous in the short-circuited path
3. This can start fires or damage components

This is why **fuses** and **circuit breakers** exist: they break the circuit if current exceeds a safe level.

### Fuse Mechanism

```
FUSE:
    ----[thin wire]----
    
    Normal operation: thin wire carries current, stays intact
    Fault (too much current): thin wire heats up rapidly, MELTS, breaks circuit
    
    Once blown, must be replaced (single use).
    
CIRCUIT BREAKER:
    Same principle, but uses an electromagnetic or thermal mechanism.
    Can be reset (reusable).
```

---

## Voltage Dividers

A **voltage divider** is two resistors in series used to produce an output voltage that is a fraction of the input voltage.

```
VOLTAGE DIVIDER CIRCUIT:

        V_in
         |
        [R₁]
         |
         +----> V_out
        [R₂]
         |
        GND

    V_out = V_in × R₂ / (R₁ + R₂)
```

### Derivation

Since R₁ and R₂ are in series:
```
I = V_in / (R₁ + R₂)

V_out = I × R₂ = V_in × R₂ / (R₁ + R₂)
```

### Examples of Voltage Divider Use

```
Example 1: Half the voltage
    R₁ = R₂ = 10 kΩ, V_in = 12V
    V_out = 12 × 10/(10+10) = 12 × 0.5 = 6V

Example 2: One third the voltage  
    R₁ = 20 kΩ, R₂ = 10 kΩ, V_in = 9V
    V_out = 9 × 10/(20+10) = 9 × 0.333 = 3V
```

### Voltage Divider with a Sensor

A clever application: replace one resistor with a **thermistor** or **LDR** (light dependent resistor). As the sensor's resistance changes with temperature or light, V_out changes. A microcontroller can read V_out to measure temperature or light level.

```
TEMPERATURE SENSOR CIRCUIT:

    V_in = 5V
     |
    [R_fixed = 10kΩ]
     |
     +--------> to analog input (measures V_out)
    [R_thermistor]  (R changes with temperature)
     |
    GND

    Hot temperature:  R_thermistor decreases --> V_out decreases
    Cold temperature: R_thermistor increases --> V_out increases
```

---

## Why House Wiring Uses Parallel

Every socket and light fitting in your house is wired in **parallel** — connected directly to the same live and neutral wires. This design choice has several crucial advantages:

### Reason 1: All Devices Get Full Voltage

In a parallel circuit, every branch sees the full supply voltage. In the UK, every socket delivers 230V. If wiring were series, each device would only get a fraction of the total.

### Reason 2: Devices Work Independently

In a parallel circuit, each device has its own branch. Switching off one device (opening that branch) does not affect others.

```
PARALLEL HOUSE WIRING:
                                    
    Mains ====+=========+=========+===...
     230V     |         |         |
             [L₁]      [R₁]      [R₂]
             lamp     cooker    TV
              |         |         |
    Neutral ==+---------+---------+---...
    
    Turn off lamp: R₁ and R₂ keep working.
    Kettle blows fuse: only kettle branch breaks.
    All appliances get 230V regardless.
```

### Reason 3: Adding More Devices Doesn't Affect Others

Adding a new device in parallel adds a new branch without changing the voltage or current in existing branches. If you connect more appliances in series, each new device would reduce the current for all others.

### Reason 4: Fault Isolation

Each circuit has its own fuse or circuit breaker. A fault in one room only trips that circuit, not the whole house.

---

## Worked Examples

### Worked Example 1: Series Circuit

**Problem:** Three resistors R₁ = 2Ω, R₂ = 3Ω, R₃ = 5Ω are connected in series to a 20V battery. Find: (a) total resistance, (b) current, (c) voltage across each resistor.

```
CIRCUIT DIAGRAM:
    +--[2Ω]--[3Ω]--[5Ω]--+
    |                      |
   20V                    
    |                      |
    +----------------------+
```

**Solution:**

```
(a) Total resistance:
    R_total = R₁ + R₂ + R₃
    R_total = 2 + 3 + 5 = 10 Ω

(b) Current (same throughout series circuit):
    I = V / R_total = 20 / 10 = 2 A

(c) Voltage across each resistor:
    V₁ = I × R₁ = 2 × 2 = 4 V
    V₂ = I × R₂ = 2 × 3 = 6 V
    V₃ = I × R₃ = 2 × 5 = 10 V
    
    Check: V₁ + V₂ + V₃ = 4 + 6 + 10 = 20 V ✓
```

### Worked Example 2: Parallel Circuit

**Problem:** Three resistors R₁ = 6Ω, R₂ = 12Ω, R₃ = 4Ω are connected in parallel to a 12V battery. Find: (a) total resistance, (b) current through each resistor, (c) total current.

**Solution:**

```
(a) Total resistance:
    1/R_total = 1/R₁ + 1/R₂ + 1/R₃
    1/R_total = 1/6 + 1/12 + 1/4
    
    Convert to common denominator (12):
    1/R_total = 2/12 + 1/12 + 3/12 = 6/12 = 1/2
    
    R_total = 2 Ω

(b) Voltage is same for all branches (V = 12V)
    I₁ = V/R₁ = 12/6  = 2.0 A
    I₂ = V/R₂ = 12/12 = 1.0 A
    I₃ = V/R₃ = 12/4  = 3.0 A

(c) Total current:
    I_total = I₁ + I₂ + I₃ = 2 + 1 + 3 = 6 A
    
    Check: I_total = V/R_total = 12/2 = 6 A ✓
```

### Worked Example 3: Mixed Series-Parallel Circuit

**Problem:** R₁ = 4Ω is in series with a parallel combination of R₂ = 6Ω and R₃ = 12Ω. The battery is 24V. Find the current through each resistor.

```
CIRCUIT:
    +---[R₁=4Ω]---+--------+
    |              |        |
   24V            [R₂=6Ω] [R₃=12Ω]
    |              |        |
    +--------------+--------+
```

**Solution:**

```
Step 1: Simplify parallel combination R₂ || R₃
    1/R_23 = 1/6 + 1/12 = 2/12 + 1/12 = 3/12 = 1/4
    R_23 = 4 Ω

Step 2: Total resistance (R₁ in series with R_23)
    R_total = R₁ + R_23 = 4 + 4 = 8 Ω

Step 3: Total current (through R₁)
    I_total = V / R_total = 24 / 8 = 3 A
    So I₁ = 3 A (all current goes through R₁)

Step 4: Voltage across parallel combination
    V_23 = I_total × R_23 = 3 × 4 = 12 V
    
    (Alternatively: V_23 = 24 - V₁ = 24 - (3×4) = 24 - 12 = 12 V ✓)

Step 5: Current through R₂ and R₃ (both see V_23 = 12V)
    I₂ = V_23 / R₂ = 12 / 6  = 2 A
    I₃ = V_23 / R₃ = 12 / 12 = 1 A
    
    Check: I₂ + I₃ = 2 + 1 = 3 A = I₁ ✓
```

### Worked Example 4: Battery with Internal Resistance

**Problem:** A battery has EMF 9V and internal resistance 1Ω. It is connected to an external resistor of 8Ω. Find: (a) current, (b) terminal voltage, (c) voltage "lost" to internal resistance.

**Solution:**

```
(a) Current:
    Total resistance = R_external + r = 8 + 1 = 9 Ω
    I = EMF / R_total = 9 / 9 = 1 A

(b) Terminal voltage:
    V_terminal = EMF - I × r
    V_terminal = 9 - (1 × 1) = 9 - 1 = 8 V

(c) Voltage lost internally:
    V_lost = I × r = 1 × 1 = 1 V
    
    Check: V_terminal + V_lost = 8 + 1 = 9 V = EMF ✓
```

### Worked Example 5: Voltage Divider Design

**Problem:** Design a voltage divider to produce 3.3V from a 5V supply using standard resistors. Choose R₁ and R₂ such that R₁ + R₂ = 10 kΩ.

**Solution:**

```
We need:
    V_out = V_in × R₂ / (R₁ + R₂)
    3.3 = 5 × R₂ / 10000

Solve for R₂:
    R₂ / 10000 = 3.3 / 5 = 0.66
    R₂ = 6600 Ω = 6.6 kΩ

Therefore:
    R₁ = 10000 - 6600 = 3400 Ω = 3.4 kΩ

Standard resistors closest: R₁ = 3.3 kΩ, R₂ = 6.8 kΩ
    V_out = 5 × 6800/(3300+6800) = 5 × 6800/10100 = 5 × 0.673 = 3.37V ≈ 3.3V ✓
```

---

## Summary

- **Series circuit**: single current path; R_total = R₁+R₂+R₃; current same everywhere; voltages add up
- **Parallel circuit**: multiple current paths; 1/R_total = 1/R₁+1/R₂+1/R₃; voltage same everywhere; currents add up
- Adding resistors in **series** always **increases** total resistance; adding in **parallel** always **decreases** it
- **Mixed circuits**: simplify by identifying series and parallel groups, replace with equivalent resistance, repeat
- A real battery has **internal resistance** r; terminal voltage = EMF - Ir, which drops as current increases
- **Open circuit**: break in path, no current, full voltage appears across the gap
- **Short circuit**: unintended low-resistance path, very high current, dangerous heating — fuses protect against this
- **Voltage divider**: V_out = V_in × R₂/(R₁+R₂) — used to create lower voltage from higher supply, and with sensors
- **House wiring is parallel** so every device gets full voltage, devices work independently, and faults are isolated

---

## Key Equations

```
Series Circuit:
    R_total = R₁ + R₂ + R₃ + ...
    I_total = I₁ = I₂ = I₃   (same current everywhere)
    V_total = V₁ + V₂ + V₃   (voltages add)

Parallel Circuit:
    1/R_total = 1/R₁ + 1/R₂ + 1/R₃ + ...
    V_total = V₁ = V₂ = V₃   (same voltage everywhere)
    I_total = I₁ + I₂ + I₃   (currents add)
    
Two resistors in parallel (shortcut):
    R_total = (R₁ × R₂) / (R₁ + R₂)

Battery Terminal Voltage:
    V_terminal = EMF - I × r
    
    EMF = electromotive force (V)
    r   = internal resistance (Ω)
    I   = current (A)

Voltage Divider:
    V_out = V_in × R₂ / (R₁ + R₂)
```
