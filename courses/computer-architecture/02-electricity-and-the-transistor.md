# Chapter 02: Electricity and the Transistor — The Switch That Changed Everything

Every calculation your computer performs, every pixel on your screen, every network packet that reaches your phone — all of it ultimately comes down to tiny switches made from silicon turning on and off. Understanding those switches is the foundation of everything in this course.

## Table of Contents

1. [What Is Electricity?](#1-what-is-electricity)
2. [Voltage, Current, and Resistance](#2-voltage-current-and-resistance)
3. [Conductors, Insulators, and Semiconductors](#3-conductors-insulators-and-semiconductors)
4. [The P-N Junction — The Key to Everything](#4-the-p-n-junction)
5. [The Transistor — A Controlled Switch](#5-the-transistor--a-controlled-switch)
6. [MOSFET: The Transistor Inside Your CPU](#6-mosfet-the-transistor-inside-your-cpu)
7. [Why Transistors Are So Important](#7-why-transistors-are-so-important)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is Electricity?

Everything around you is made of atoms. An atom has a nucleus (protons and neutrons) surrounded by electrons. Protons carry a positive electric charge. Electrons carry a negative electric charge. Under normal conditions, an atom has equal numbers of protons and electrons, so it is electrically neutral.

In some materials — particularly metals like copper — the outermost electrons are loosely bound to their atoms. They can be pushed from atom to atom by an external force, like water molecules flowing through a pipe. This flow of electrons is what we call **electric current**.

```
Imagine a pipe full of marbles:

  ────────[●●●●●●●●●●]────────►
  Push here                   marbles exit here

If you push a marble in one end, a marble pops out the other almost instantly.
Electrons in a wire behave similarly.
```

The "push" that drives electrons through a wire comes from a difference in electric potential between two points — we call this difference **voltage**.

### The Water Analogy

The relationship between voltage, current, and resistance maps almost perfectly onto water flowing through a pipe:

| Electricity | Water |
|---|---|
| Voltage (V) | Water pressure (difference in height) |
| Current (A, amperes) | Flow rate (liters per second) |
| Resistance (Ω, ohms) | Pipe narrowness (friction) |

A tall tank of water (high pressure → high voltage) pushes water through a pipe faster (more current) unless the pipe is very narrow (high resistance → less current).

### Quick Check

> 1. What particles carry electric current through a wire?
> 2. Using the water analogy, what does voltage correspond to?
> 3. If you increase the resistance in a circuit without changing the voltage, what happens to the current?

---

## 2. Voltage, Current, and Resistance

### Voltage

Voltage is the "electrical pressure" — the potential energy difference between two points that drives electrons to flow. Measured in **volts (V)**. A AA battery is 1.5V. A laptop power supply might be 19V. Your home wall outlet is 230V (India) or 120V (USA).

In digital electronics, two voltage levels are used to represent the two binary values:
- **Logic HIGH (1)**: typically 0.8V to 3.3V depending on the technology
- **Logic LOW (0)**: typically 0V to 0.4V

Modern CPU cores operate at around **0.7–1.0V** (this low voltage is key to low power consumption).

### Current

Current is the flow rate of electrons. Measured in **amperes (A)**. A smartphone charges at about 3-20A. A CPU core might draw 30-50A at high load (at ~1V, that is 30-50 watts).

### Resistance

Resistance opposes current flow. Measured in **ohms (Ω)**. Ohm's Law: **V = I × R** (Voltage = Current × Resistance).

### Power

Power is the rate of energy consumption. **P = V × I** (Power = Voltage × Current), measured in **watts (W)**. This is why reducing voltage in CPUs reduces power dramatically: cutting voltage from 1.2V to 0.8V (a 33% reduction) while keeping the same current cuts power by 33%. In practice, reducing voltage also allows reducing current, so total power can drop by 50-60% — this is why modern phones get so much battery life despite fast processors.

### Quick Check

> 1. Using Ohm's Law (V = IR), if a circuit has 5V and 10Ω resistance, what is the current?
> 2. A CPU core runs at 1.0V and draws 40A. What is its power consumption in watts?
> 3. Why does reducing CPU voltage from 1.2V to 0.9V reduce power by more than just the 25% voltage drop?

---

## 3. Conductors, Insulators, and Semiconductors

Not all materials behave the same way when voltage is applied to them. There are three categories:

### Conductors

Conductors allow electrons to flow freely. Metals are conductors because their outermost electrons are loosely bound and free to move. Examples: copper (used in wires), gold (used for chip contacts because it does not corrode), aluminum (used in chip interconnects).

A copper wire has resistance of about 0.017 Ω per meter — very low, so very little voltage is needed to drive current through it.

### Insulators

Insulators resist electron flow. Their electrons are tightly bound and cannot move freely. Examples: glass, rubber, plastic, silicon dioxide (SiO₂). The plastic coating on a wire prevents electric shock and stops current from jumping between adjacent wires.

### Semiconductors

Semiconductors are the critical middle ground. At low temperatures, they behave like insulators. But you can make them conduct by:
1. **Heating them** (adding thermal energy frees electrons)
2. **Shining light on them** (photon energy frees electrons — this is how solar cells work)
3. **Adding impurities** (called doping — this is what makes transistors possible)

Silicon (Si) is the semiconductor the electronics industry is built on. Why silicon?
- Abundant (the second most common element in Earth's crust)
- Forms a strong, stable crystal structure
- Can be purified to extreme levels (99.9999999% pure)
- Forms an excellent insulating oxide (SiO₂) naturally
- Has an ideal bandgap for electronic switching at room temperature

```
     Conductors              Semiconductors           Insulators
 (copper, aluminum)      (silicon, germanium)      (glass, rubber)
 
 ───────────────┼──────────────────────────────────┼──────────────
 Very low          ← resistance →                    Very high
 (electrons flow freely)           (electrons cannot flow)
```

### Quick Check

> 1. Why is copper used for wires but silicon is used for transistors?
> 2. What is doping and why does it matter for making transistors?
> 3. Silicon dioxide (SiO₂, i.e., quartz/sand) is an insulator. Silicon (Si) is a semiconductor. Why does this matter for chip manufacturing?

---

## 4. The P-N Junction

The magic of semiconductors becomes practical through **doping**: intentionally adding impurities to change silicon's electrical properties.

### N-type Silicon

Silicon has 4 outer electrons. Phosphorus has 5 outer electrons. Add a tiny amount of phosphorus to silicon (one phosphorus atom per million silicon atoms) and the extra electron from each phosphorus atom is free to move. These free electrons are **negative charge carriers** — hence "N-type."

### P-type Silicon

Boron has 3 outer electrons. Add boron to silicon and each boron atom creates a "hole" — a missing electron that neighboring electrons can jump into, making the hole appear to move. These mobile holes are **positive charge carriers** — hence "P-type."

### The Junction

Place N-type and P-type silicon in contact. Electrons from the N-side diffuse toward the P-side; holes from the P-side diffuse toward the N-side. They combine and annihilate each other in a thin region called the **depletion region**. This creates an electric field that eventually stops further diffusion.

```
  N-type silicon  │  P-type silicon
  ────────────────┼────────────────
   free electrons │ holes (missing electrons)
          →  →  [depletion region]  ←  ←
   − − − − − │//////////////////│ + + + + +
             │   electric field  │
             │      points →     │
```

### The Diode Behavior

Apply voltage with + on the P-side (forward bias): the electric field is overcome, electrons flow, current passes. Apply voltage with + on the N-side (reverse bias): the electric field grows stronger, no electrons flow, no current.

**A P-N junction is a diode** — a one-way valve for current. This is not yet a switch (which needs three terminals — control, input, output). But the P-N junction is the foundation for the transistor.

### Quick Check

> 1. What makes N-type silicon "N-type"?
> 2. What is a "hole" in P-type silicon?
> 3. A diode allows current in one direction. How is this different from what we need for a transistor (a switch)?

---

## 5. The Transistor — A Controlled Switch

The transistor was invented in December 1947 at Bell Laboratories by William Shockley, John Bardeen, and Walter Brattain. The three received the Nobel Prize in Physics in 1956.

The key insight: if you sandwich a thin layer of P-type silicon between two N-type regions (or vice versa), you get a device with three terminals where a small current at one terminal controls a large current between the other two.

### The BJT (Bipolar Junction Transistor)

The first transistors were Bipolar Junction Transistors (BJTs). An NPN transistor has:
- **Base** (B): the control terminal — connected to the thin P-type layer in the middle
- **Collector** (C): one N-type region — where current enters
- **Emitter** (E): the other N-type region — where current exits

```
         Collector (C)
              │
              │
         ─────┼─────  ← thin P-type base region
              │
         Base (B) ───► small control current
              │
         Emitter (E)
              │
```

When a small current flows into the Base, it forward-biases the Base-Emitter junction. This allows a large current to flow from Collector to Emitter — **the small base current controls the large collector current**. This is amplification. When the base current is zero, no collector current flows — the transistor is "off." When sufficient base current flows, maximum collector current flows — the transistor is "on."

A transistor can thus function as:
- **A switch**: fully on or fully off (digital electronics)
- **An amplifier**: partially on, with output proportional to input (analog electronics)

In computers, transistors are used as switches.

### Quick Check

> 1. What are the three terminals of a BJT transistor?
> 2. How does a small base current control a large collector-emitter current?
> 3. In digital circuits, transistors are used as switches, not amplifiers. What does "fully on" and "fully off" mean for a transistor?

---

## 6. MOSFET: The Transistor Inside Your CPU

The transistors in modern computers are not BJTs — they are MOSFETs: **Metal-Oxide-Semiconductor Field-Effect Transistors**. Every modern processor, memory chip, and digital circuit is built from MOSFETs.

### Why MOSFET Instead of BJT?

BJTs require current at the base to control the switch. MOSFETs are controlled by **voltage** at the gate (not current). Since voltage control requires essentially zero power (static, not flowing), MOSFETs are far more power-efficient. This matters enormously when you have 20 billion of them on one chip.

### MOSFET Structure

```
                     Gate (G)
                       │
        ┌──────────────▼──────────────┐
        │         Gate Oxide (SiO₂)   │
        │  ┌───────────────────────┐  │
        │  │       P-substrate     │  │
        │  │  N+  ╔═════════╗  N+  │  │
        │  └──┬───╨─────────╨───┬──┘  │
        └─────┘                 └─────┘
         Source (S)           Drain (D)
```

The MOSFET has three terminals:
- **Gate (G)**: the control terminal — insulated from the silicon by a thin oxide layer
- **Source (S)**: where current enters
- **Drain (D)**: where current exits

Between Source and Drain is a P-type semiconductor. When a positive voltage is applied to the Gate, it attracts electrons from the P-type bulk, creating a thin N-type **channel** that connects Source to Drain. Current can now flow. When the gate voltage is zero, no channel forms, no current flows.

The "Metal-Oxide-Semiconductor" name refers to the structure: Metal gate on top, Oxide (SiO₂) insulator in the middle, Semiconductor (silicon) substrate at the bottom.

### NMOS and PMOS

There are two types of MOSFETs:
- **NMOS** (N-channel): Gate voltage HIGH turns it ON (positive gate voltage attracts electrons to form N-channel)
- **PMOS** (P-channel): Gate voltage LOW turns it ON (negative gate voltage attracts holes to form P-channel)

### CMOS: The Genius of Complementary Pairs

**CMOS** (Complementary MOS) pairs one NMOS and one PMOS transistor together. This is how every logic gate in every modern chip is built.

A CMOS inverter (NOT gate):

```
     VDD (power supply, e.g., 1V)
       │
    ┌──┴──┐
    │PMOS │  ← turns ON when input is LOW
    └──┬──┘
       │
       ├──── Output
       │
    ┌──┴──┐
    │NMOS │  ← turns ON when input is HIGH
    └──┬──┘
       │
     GND (0V)

When input is HIGH:  PMOS OFF, NMOS ON  → Output connected to GND → Output is LOW
When input is LOW:   PMOS ON, NMOS OFF  → Output connected to VDD → Output is HIGH
```

The brilliant property of CMOS: **in a stable state, one transistor is always off**. Current only flows during the brief moment of switching. This means CMOS draws nearly zero power when nothing is changing — making it dramatically more power-efficient than earlier technologies.

### The Modern FinFET

As transistors shrank below 22nm, the flat planar MOSFET stopped working well — the gate could not control the thin channel effectively. Intel introduced the **FinFET** (Fin Field-Effect Transistor) in 2011 at 22nm. The channel is a thin vertical fin of silicon, and the gate wraps around three sides of it — giving much better control.

```
     Traditional Planar MOSFET:     Modern FinFET:
     
        Gate                              Gate
     ───────────                          ┌─┐
     ══════════════  oxide          ─────│F│─────  gate wraps
     ──────────────  channel        ─────│I│─────  around fin
     S          D                  S ───│N│─── D
     ────────────────               ─────│S│─────
                                        └─┘
```

Today's most advanced transistors use **GAA (Gate-All-Around) nanosheets** — the gate wraps completely around all four sides of a thin horizontal nanosheet of silicon. TSMC's N2 node (2025) uses this technology.

### Quick Check

> 1. What is the key advantage of MOSFET over BJT for use in computers?
> 2. In CMOS, when input is HIGH, which transistor is ON: NMOS or PMOS?
> 3. Why was the FinFET invented, and what problem did it solve?

---

## 7. Why Transistors Are So Important

The transistor is arguably the most important invention in human history. Here is why:

### The Numbers Tell the Story

| Year | Transistors per chip | Process node |
|------|---------------------|--------------|
| 1971 | 2,300 (Intel 4004) | 10,000 nm |
| 1982 | 134,000 (Intel 286) | 1,500 nm |
| 1993 | 3.1 million (Pentium) | 800 nm |
| 2000 | 42 million (Pentium 4) | 180 nm |
| 2012 | 2.27 billion (Core i7) | 22 nm |
| 2020 | 15 billion (Apple M1) | 5 nm |
| 2024 | 28 billion (Apple M4) | 3 nm |

In 53 years, transistor count increased by a factor of ~12 million. Speed increased by a factor of ~10,000. Cost per transistor fell by a factor of a billion.

### The Economic Consequence

Today, a single transistor switching costs less than $0.0000000001 (a ten-billionth of a cent). Performing one billion arithmetic operations on a modern CPU costs about one-tenth of a cent in electricity. The entire Netflix library — hundreds of thousands of hours of video — is encoded in bits and stored on drives made of trillions of transistors costing a few million dollars total.

This cost collapse is unprecedented in human history. No other technology has ever improved by 12 million times in 50 years. That extraordinary improvement is why software is eating the world.

### Why the Transistor Enables the Computer

The computer does not need transistors specifically — it needs switches. Any reliable, fast, miniaturizable switch would do. Transistors happen to be the best switches available:
- **Fast**: modern transistors switch in picoseconds (10⁻¹² seconds)
- **Small**: a modern transistor's gate is ~5nm — 10,000x smaller than a human hair
- **Reliable**: transistors can switch trillions of times without wearing out
- **Cheap**: billions can be manufactured simultaneously on a single wafer
- **Low power**: CMOS draws power only when switching

No other switching technology comes close on all five dimensions simultaneously.

### Quick Check

> 1. By approximately what factor did transistor count per chip increase from 1971 to 2024?
> 2. A transistor gate today is about 5nm. How many times smaller is this than a human hair (approximately 50,000nm wide)?
> 3. What are the five properties that make transistors better than any alternative switching technology?

---

## Summary

- Electricity is the flow of electrons driven by a voltage difference. **Voltage** is pressure, **current** is flow rate, **resistance** opposes flow. **P = V × I** is power.
- Materials are conductors (electrons flow freely), insulators (electrons don't flow), or semiconductors (can be made to do either).
- Silicon is the semiconductor of choice because it is abundant, pure-able, and forms an excellent insulating oxide.
- **Doping** adds impurities to create N-type silicon (extra electrons) or P-type silicon (holes). A P-N junction is a diode — current flows one way only.
- A **transistor** adds a third terminal (gate or base) that controls current flow between the other two — it is a voltage-controlled switch.
- **MOSFET** transistors are controlled by gate voltage (not current), making them far more power-efficient than BJTs.
- **CMOS** pairs NMOS and PMOS transistors so exactly one is off at all times — drawing power only during switching.
- **FinFET** (and now GAA nanosheet) transistors wrap the gate around more sides of the channel for better control at small sizes.
- The transistor is the most important invention in history: transistor count has grown by 12 million times in 53 years, enabling the modern digital world.

---

## Exercises

### Easy

1. A circuit runs at 5V and 2A. How much power (in watts) does it consume? If you cut the voltage in half to 2.5V but keep the same current, what is the new power?

2. Draw a simple CMOS inverter (like the one in the chapter) with PMOS on top and NMOS on bottom. For each of the four possible states (Input=HIGH/LOW, which transistors are ON/OFF), fill in the output voltage.

3. What is the difference between NMOS and PMOS? Make a small table with the gate voltage condition that turns each ON.

### Medium

4. Modern CPUs operate at around 1V rather than the 5V of chips from the 1980s. Using P = V × I, if current draw is proportional to frequency, and a 1980s chip ran at 5V/100MHz while a modern chip runs at 1V/4GHz: (a) How many times higher is the frequency? (b) If current scales proportionally to frequency, how many times higher is the modern chip's current? (c) Compute the ratio of power consumption. Does the voltage reduction help or hurt?

5. The MOSFET's gate is insulated from the silicon by a thin oxide layer. This means essentially zero DC current flows into the gate. Why is this important for using MOSFETs as switches in a computer? Compare to a BJT where base current is required.

6. FinFET transistors have the gate wrapping around three sides of a silicon fin. The earlier planar MOSFET had the gate on only one side. Why does wrapping the gate around more sides of the channel improve transistor performance at small sizes? Think about what the gate is trying to do (control the channel) and what happens when the transistor is very thin.

### Hard

7. A modern TSMC 3nm chip has approximately 290 million transistors per mm². An Intel 4004 (1971, 10µm process) had 2,300 transistors on a 12mm² die. 
   (a) Calculate the transistor density of the Intel 4004 in transistors per mm².  
   (b) Calculate the ratio of modern to 1971 transistor density.  
   (c) If Moore's Law predicts doubling every 2 years, and 53 years have passed since 1971, how many doublings have occurred? What is the predicted density ratio?  
   (d) How does the actual ratio compare to the predicted ratio?

8. CMOS draws power only when transistors switch. But modern CPUs have transistors switching billions of times per second. Research the concept of "dynamic power" (P = α × C × V² × f, where α is activity factor, C is capacitance, V is voltage, f is frequency). Explain each term: what is activity factor? What is capacitance? Why is voltage squared? If you double the frequency (f) while keeping everything else constant, how does power change? If you also reduce voltage by 20%, what is the new power?
