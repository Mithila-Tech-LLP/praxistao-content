# Chapter 07: The Iron–Carbon Phase Diagram — Foundation of All Steel

> **"The iron-carbon phase diagram is the most studied, most useful, and most economically important phase diagram in human history. Every piece of steel ever made — every bridge, every ship hull, every surgical tool — was designed using this diagram, whether the engineer knew it or not."**

---

## Table of Contents

1. [Why Iron-Carbon?](#1-why-iron-carbon)
2. [The Key Phases — Meet the Players](#2-the-key-phases--meet-the-players)
3. [The Fe-C Phase Diagram — Reading It](#3-the-fe-c-phase-diagram--reading-it)
4. [Regions and Reactions at a Glance](#4-regions-and-reactions-at-a-glance)
5. [The Eutectoid Reaction — The Heart of Steel](#5-the-eutectoid-reaction--the-heart-of-steel)
6. [The Eutectic Reaction — Cast Iron Territory](#6-the-eutectic-reaction--cast-iron-territory)
7. [Cooling of Hypoeutectoid Steel (< 0.77% C)](#7-cooling-of-hypoeutectoid-steel--077-c)
8. [Cooling of Hypereutectoid Steel (> 0.77% C)](#8-cooling-of-hypereutectoid-steel--077-c)
9. [The Metastable Iron-Carbide Diagram vs. Stable Iron-Graphite](#9-the-metastable-fe-fec-diagram-vs-stable-fe-c-diagram)
10. [Steel vs. Cast Iron — The Carbon Line](#10-steel-vs-cast-iron--the-carbon-line)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Why Iron-Carbon?

Iron is the fourth most abundant element in the Earth's crust. Carbon is everywhere — in charcoal, coke, wood, limestone. The combination produces **steel**, which is strong, tough, hard (when treated), ductile (when annealed), and cheaply produced in enormous quantities.

Adding just **0.1 to 2.0% carbon** by weight transforms soft, weak pure iron into engineering steel. The reason this tiny amount of carbon matters so dramatically lies entirely in the iron-carbon phase diagram.

The diagram encodes:
- Where and how carbon dissolves in iron
- Where and how it precipitates out
- What microstructures form on heating and cooling
- Why quenching produces hard steel
- Why annealing produces soft steel
- Why cast iron is different from steel

---

## 2. The Key Phases — Meet the Players

Before reading the diagram, you must know the phases:

### α-Ferrite (α)
- Crystal structure: **BCC**
- Maximum carbon solubility: **0.022 wt% C** at 727°C (practically zero)
- Properties: soft (HB ~90), ductile, **magnetic** below 770°C
- Exists at: room temperature, up to 912°C
- In microstructure: light-colored equiaxed grains, often the majority phase in low-carbon steel

### δ-Ferrite (δ)
- Crystal structure: **BCC** (same as α, different temperature range)
- Stable from 1394°C to 1538°C (just below melting)
- Maximum carbon solubility: 0.09 wt% C
- Mostly important for solidification; rarely significant in final parts

### Austenite (γ)
- Crystal structure: **FCC**
- Maximum carbon solubility: **2.14 wt% C** at 1147°C — nearly 100× more than ferrite!
- Stable between 727°C (for eutectoid steel) and 1394°C
- Properties: soft, non-magnetic, highly formable
- Critical: **all steel heat treatment passes through austenite**

The dramatic difference in carbon solubility between BCC ferrite (0.022%) and FCC austenite (2.14%) is why heat treatment works. Heat the steel → transform to austenite → carbon dissolves → control what happens on cooling.

### Cementite (Fe₃C)
- Also called **iron carbide**
- Crystal structure: orthorhombic
- Composition: exactly 6.70 wt% C (25 at% C)
- Properties: extremely hard (HV ~1000), very brittle
- Appears in almost all steels and cast irons as the carbon-bearing phase

### Pearlite
- **Not a single phase** — it's a two-phase microstructure
- Alternating lamellae of ferrite and cementite
- Forms from austenite of eutectoid composition on slow cooling
- Name: looks like mother-of-pearl under optical microscope
- Properties: intermediate strength (~700 MPa UTS), good toughness

### Martensite
- **Not an equilibrium phase** — forms by diffusionless transformation on rapid quenching
- Crystal structure: **BCT** (body-centered tetragonal — BCC distorted by trapped carbon)
- Properties: extremely hard (HV can reach 800+), very brittle (in high carbon)
- The product of quenching; the hardest form of steel

### Bainite
- Also not an equilibrium phase; forms between pearlite temperatures and martensite start temperature (M_s)
- Upper bainite: feathery ferrite laths with cementite between
- Lower bainite: similar to tempered martensite; fine carbides within plates
- Properties: good combination of strength and toughness

---

## 3. The Fe-C Phase Diagram — Reading It

```
Temperature (°C)
1538 ┤ ─── δ ───● ────────── Liquid ───────────
     |          | 0.53%   ↑ Peritectic (1493°C)
1494 ┤─── δ ───●───── Liquid ──────────────────
     |      δ+γ     ↑ 0.17% C
     |              γ (austenite, FCC)
1147 ┤──────────────────────●───────────────────  ← Eutectic (1147°C)
     |     γ                  γ+Fe₃C              4.3% C
     |                        ← Fe₃C is cementite
     |
 912 ┤──●───────────────────────────────────────  ← A₃ line
     | α/γ boundary 
     |   γ (austenite)
 727 ┤───────────────●─────────────────────────  ← A₁ line (eutectoid, 727°C)
     | α+γ       eutectoid    γ+Fe₃C            0.77% C
     |
     |  α (ferrite, BCC)   α + Fe₃C (cementite)
     |
   0 ┤───────────────────────────────────────────
     0    0.022  0.77        2.14       4.3    6.70
                  ↑                            ↑
            eutectoid                    Fe₃C composition
            composition
     ←── steels ──→←────────── cast irons ──────────→
           (< 2.14% C)            (2.14–6.70% C)
```

### Critical Temperatures in Plain Carbon Steels:

| Temperature | Name | What Happens |
|-------------|------|--------------|
| 727°C | **A₁ (lower critical)** | Eutectoid reaction; pearlite ↔ austenite |
| 912°C | **A₃ (upper critical)** | α+γ → γ for hypoeutectoid; γ → α+γ on cooling |
| 1147°C | **A**cm or eutectic | Fe₃C dissolves into austenite / eutectic reaction |

---

## 4. Regions and Reactions at a Glance

| Region | Phases Present | Carbon Content |
|--------|---------------|----------------|
| Ferrite only (α) | BCC iron, ~0% C | 0–0.022% C, below 727°C |
| Austenite (γ) | FCC iron, dissolved C | 0–2.14% C, 727°C–1394°C |
| α + γ | Ferrite + Austenite | 0–0.77% C, between A₁ and A₃ |
| γ + Fe₃C | Austenite + Cementite | 0.77–2.14% C, between A₁ and Acm |
| α + Fe₃C | Ferrite + Cementite | 0–6.70% C, below A₁ |
| Liquid | Molten Fe-C | above liquidus |
| γ + Liquid | Mushy zone | around casting temperatures |

---

## 5. The Eutectoid Reaction — The Heart of Steel

The **eutectoid reaction** at 727°C, 0.77 wt% C is the most important reaction in the diagram:

```
γ (0.77% C)  →  α (0.022% C)  +  Fe₃C (6.70% C)
(austenite)      (ferrite)         (cementite)
```

On slow cooling, austenite of eutectoid composition simultaneously ejects ferrite (which can hold almost no carbon) and cementite (the carbon dump). They form alternating lamellae — **pearlite**:

```
Pearlite microstructure (cross-section):
 ▓ = ferrite (α)  □ = cementite (Fe₃C)
 
 ▓▓▓▓▓▓□□▓▓▓▓▓▓□□▓▓▓▓▓▓□□▓▓▓▓▓▓□□
 ▓▓▓▓▓▓□□▓▓▓▓▓▓□□▓▓▓▓▓▓□□▓▓▓▓▓▓□□
 ▓▓▓▓▓▓□□▓▓▓▓▓▓□□▓▓▓▓▓▓□□▓▓▓▓▓▓□□
 (lamellar alternation; spacing depends on cooling rate)
```

**Volume fractions in pearlite** (from lever rule at 727°C):
- Fraction cementite = (0.77 - 0.022) / (6.70 - 0.022) = 0.748/6.678 ≈ **11.2%**
- Fraction ferrite = **88.8%**

Pearlite is thus mostly ferrite with thin cementite lamellae — the cementite provides hardness while the ferrite provides toughness.

**Pearlite interlamellar spacing:**
- Coarse (slow cooling, high T): ~500 nm → soft
- Fine (faster cooling, lower T): ~100 nm → harder
- Extremely fine: **bainite** (different transformation mechanism, T < ~550°C)

---

## 6. The Eutectic Reaction — Cast Iron Territory

At 1147°C, 4.3 wt% C:
```
L (4.3% C)  →  γ (2.14% C)  +  Fe₃C (6.70% C)
(liquid)       (austenite)       (cementite)
```

This is the **cast iron eutectic**. The product is called **ledeburite** (austenite + cementite mixture). On further cooling, the austenite itself undergoes the eutectoid reaction → pearlite + more cementite.

Cast irons (2.14–6.70% C) have much more carbon than steels. The excess carbon can't dissolve — it forms large amounts of Fe₃C (or precipitates as graphite in grey/ductile irons). This makes cast irons:
- Easy to cast (low melting point, good fluidity)
- Hard and brittle (white cast iron with lots of cementite)
- Or — if carbon precipitates as graphite — reasonable toughness (grey, ductile)

---

## 7. Cooling of Hypoeutectoid Steel (< 0.77% C)

"Hypoeutectoid" = carbon content **below** 0.77% (on the left side of the eutectoid point).

**Example: 0.40 wt% C steel, slow cooling from 900°C**

```
Step 1: 900°C → fully austenitic (γ, 0.40% C). 
        All carbon in solution. Single phase.

Step 2: Cool to ~790°C → cross the A₃ line.
        α-ferrite begins forming at austenite grain boundaries.
        α has very little C (near 0%); C enriches the remaining γ.
        
Step 3: Cool from 790°C to 727°C.
        More α forms. Remaining γ becomes richer in C (moving along A₃ curve 
        toward 0.77%).
        At 727°C, remaining γ is exactly 0.77% C.
        
Step 4: At 727°C (A₁ line):
        Remaining γ → pearlite (via eutectoid reaction).
        
Final microstructure:
        Proeutectoid ferrite (white grains) + Pearlite colonies (dark patches)
```

The amount of each is given by the lever rule at 727°C:
```
Fraction pearlite = (0.40 - 0.022) / (0.77 - 0.022) = 0.378/0.748 ≈ 50.5%
Fraction proeutectoid ferrite ≈ 49.5%
```

A 0.40% C steel is roughly 50% ferrite, 50% pearlite → medium strength, good ductility.

---

## 8. Cooling of Hypereutectoid Steel (> 0.77% C)

"Hypereutectoid" = carbon content **above** 0.77% (right side of eutectoid).

**Example: 1.0 wt% C steel, slow cooling from 900°C**

```
Step 1: 900°C → fully austenitic.

Step 2: Cool to ~850°C → cross the Acm line.
        Proeutectoid cementite (Fe₃C) begins precipitating,
        preferentially at austenite grain boundaries.
        The network of grain-boundary cementite is a classic
        "necklace" or "cementite network" microstructure.
        
Step 3: Cool to 727°C.
        Remaining austenite has become 0.77% C.
        
Step 4: At 727°C:
        Remaining γ (0.77% C) → pearlite.
        
Final microstructure:
        Proeutectoid cementite films at grain boundaries
        + Pearlite colonies inside grains
```

The grain-boundary cementite network makes hypereutectoid steels brittle and difficult to machine. **Spheroidize annealing** (Chapter 18) dissolves this network and forms spherical cementite particles → much better machinability.

Tool steels (1.0–2.0% C) are hypereutectoid. They rely on fine carbide particles for wear resistance after heat treatment.

---

## 9. The Metastable Fe-Fe₃C Diagram vs. Stable Fe-C Diagram

There are actually **two** versions of the iron-carbon phase diagram:

**Metastable system: Fe-Fe₃C (cementite)**
- Cementite (Fe₃C) is actually metastable — given enough time at high temperature, it decomposes to iron + graphite
- But cementite persists indefinitely in most engineering steels → it is the practically relevant diagram
- This is the standard Fe-C phase diagram for steels

**Stable system: Fe-C (graphite)**
- The truly stable carbon phase in iron is **graphite** (pure carbon)
- At high silicon content (> 2%), graphite precipitates instead of cementite
- This is the basis of **grey cast iron** — silicon promotes graphite rather than cementite
- The solidification is "grey" (graphite is dark) rather than "white" (cementite is silver-white)

The difference:
```
Fe₃C (cementite) → 3 Fe + C (graphite) at high T or with Si additions
This reaction is the basis of all cast iron types:
- White cast iron: fast cooling → cementite (metastable)
- Grey cast iron: slow cooling + Si → graphite (stable)
- Ductile cast iron: Mg addition → spheroidal graphite instead of flakes
```

---

## 10. Steel vs. Cast Iron — The Carbon Line

The practical boundary is **2.14 wt% C** — the maximum carbon solubility in austenite:

| Carbon Range | Material | Character |
|-------------|----------|-----------|
| 0.022–0.30% C | Low carbon steel (mild steel) | Soft, ductile, weldable |
| 0.30–0.60% C | Medium carbon steel | Balance of strength and ductility |
| 0.60–1.40% C | High carbon steel | Hard, wear resistant |
| 1.40–2.14% C | Very high carbon steel / tool steel | Very hard, brittle |
| 2.14–4.30% C | Hypoeutectic cast iron | Hard, brittle (useful with graphite) |
| 4.30% C | Eutectic cast iron | Easy casting |
| 4.30–6.70% C | Hypereutectic cast iron | Brittle |

**Why 2.14% C is the boundary:**
- Below 2.14%: can fully austenitize → heat treat → wide range of microstructures → "steel"
- Above 2.14%: even at maximum T, excess carbon is present as cementite → cannot fully homogenize → "cast iron"

The ability to fully austenitize is what makes steels heat-treatable to martensite. Cast irons cannot be hardened by quenching in the same way.

---

## Summary

| Phase | Structure | Max C | Key Properties |
|-------|-----------|-------|----------------|
| α-Ferrite | BCC | 0.022% | Soft, ductile, magnetic |
| δ-Ferrite | BCC | 0.09% | High T, near solidification |
| Austenite (γ) | FCC | 2.14% | Heat treatment gateway; dissolves C |
| Cementite (Fe₃C) | Orthorhombic | 6.70% | Very hard, brittle |
| Pearlite | α + Fe₃C lamellae | (mixture) | Good strength and toughness |
| Martensite | BCT | (trapped) | Hardest; diffusionless transformation |

**Key reactions:**
- **Eutectoid (727°C, 0.77% C):** γ → α + Fe₃C → forms pearlite (the most important reaction)
- **Eutectic (1147°C, 4.3% C):** L → γ + Fe₃C → ledeburite (cast iron territory)
- **Peritectic (1493°C, 0.17% C):** L + δ → γ (solidification)

**Steels** = 0.022–2.14% C; can be fully austenitized; heat-treatable.  
**Cast irons** = 2.14–6.70% C; cannot be fully austenitized.

**Next chapter:** Phase diagrams tell you what's *possible* at equilibrium. But real solidification never reaches equilibrium — it's a kinetic process. Chapter 8 covers what actually happens when metal freezes: nucleation, dendrite growth, segregation, and porosity.

---

## Exercises

1. A steel with 0.20 wt% C is slowly cooled from 950°C to room temperature. (a) At what temperature does the first solid phase change occur? (b) What phase forms first? (c) What is the final microstructure? (d) Use the lever rule to find the fraction of pearlite.

2. A 1.2 wt% C hypereutectoid steel is slowly cooled. At 727°C, what fraction of the total microstructure is proeutectoid cementite vs. pearlite? (Lever rule, endpoints at 0.022% and 6.70%)

3. Explain why austenite (FCC) can dissolve 2.14% C while ferrite (BCC) dissolves only 0.022% C. Use crystal structure arguments (site sizes from Chapter 5).

4. A white cast iron contains 3.0 wt% C. It is cooled slowly from 1200°C. Describe every phase transformation it undergoes and the final room-temperature microstructure.

5. Tool steels are often given a "spheroidize anneal" before machining. What does this do to the cementite morphology? Why does this make the steel easier to machine? (Think: sharp plates vs. rounded particles and how cutting tools interact with each.)
