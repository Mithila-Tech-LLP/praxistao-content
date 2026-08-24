# Chapter 38: Thermal Expansion

> "Nature uses as little as possible of anything — but she also makes everything just a little bigger when it gets warm."
> — Johannes Kepler (paraphrased)

---

## Table of Contents

1. [Why Things Expand When Heated](#why-things-expand-when-heated)
2. [Linear Expansion](#linear-expansion)
3. [The Coefficient of Linear Expansion](#the-coefficient-of-linear-expansion)
4. [Area and Volume Expansion](#area-and-volume-expansion)
5. [Engineering Applications](#engineering-applications)
6. [Worked Example 1: Railway Track Gap](#worked-example-1)
7. [Bimetallic Strip and the Thermostat](#bimetallic-strip)
8. [Worked Example 2: Bimetallic Strip Deflection](#worked-example-2)
9. [Thermal Shock in Glass](#thermal-shock-in-glass)
10. [Anomalous Expansion of Water](#anomalous-expansion-of-water)
11. [Why Ice Floats](#why-ice-floats)
12. [Ecological Importance: Aquatic Life in Winter](#ecological-importance)
13. [Summary](#summary)
14. [Key Equations](#key-equations)

---

## Why Things Expand When Heated

In the previous chapter we learned that temperature is a measure of the average kinetic energy of particles. When you heat an object, you give its atoms and molecules more energy. They vibrate more vigorously. But here is the key: they do not vibrate symmetrically.

### The Asymmetric Potential Well

Atoms in a solid are held together by attractive forces (chemical bonds). The relationship between the force and the distance between atoms is not symmetric — it is easier to pull atoms apart slightly than to push them together. Physicists represent this with a diagram called a **potential energy well**.

```
Potential Energy of Bond vs Atom Separation

  Energy
    |
    |   \
    |    \
    |     \        (hard to push atoms
    |      \        together)
    |       \_ _ _
    |            --__
    |                ----_____     (easy to pull apart)
    |
    +────────────────────────────→ Atom Separation
              ^
              equilibrium position
```

When atoms vibrate with more energy (higher temperature), they spend more time on the "easy to pull apart" side. The **average** separation between atoms increases. This is thermal expansion.

This is why:
- Almost all materials expand when heated
- The expansion is approximately proportional to the temperature rise
- It is a bulk effect: every bond in the material expands slightly

The molecules in a solid are like balls connected by springs. At low temperature the balls barely move. At high temperature they vibrate wildly, and because the springs are asymmetric (easier to stretch than compress), the balls end up slightly farther apart on average.

```
Low Temperature:
    O─────O─────O─────O─────O
    (small vibration, close spacing)

High Temperature:
    O──────O──────O──────O──────O
    (large vibration, slightly larger average spacing)
```

---

## Linear Expansion

For objects where one dimension is much larger than the others (rods, beams, rails, wires), we focus on **linear expansion** — the change in length.

The change in length ΔL depends on three things:

1. The original length L₀ (a longer rod expands more in absolute terms)
2. The temperature change ΔT (more heating = more expansion)
3. The material (different substances expand at different rates)

The formula is:

```
ΔL = α × L₀ × ΔT

Where:
  ΔL = change in length (metres)
  α  = coefficient of linear expansion (per °C or per K)
  L₀ = original length (metres)
  ΔT = change in temperature (°C or K)

The new length after expansion:
  L = L₀ + ΔL = L₀ × (1 + α × ΔT)
```

Note: The degree-size of Celsius and Kelvin are identical, so ΔT has the same numerical value in either scale. You can use whichever is more convenient.

---

## The Coefficient of Linear Expansion

The **coefficient of linear expansion** α (Greek letter alpha) is a material property that tells you how much 1 metre of the material expands for each 1°C rise in temperature. It has units of 1/°C or °C⁻¹.

### Table of Coefficients of Linear Expansion

```
Material               | α (×10⁻⁶ per °C) | Notes
-----------------------|-------------------|------------------------
Invar (nickel-iron)    |       1.2         | Engineered for minimal expansion
Glass (borosilicate)   |       3.3         | Pyrex — used for lab glassware
Glass (ordinary)       |       8.5         | Window glass
Concrete               |      12           | Similar to steel — important!
Steel / Iron           |      12           | Used in reinforced concrete
Copper                 |      17           | Electrical wiring
Brass                  |      19           | Musical instruments, plumbing
Aluminium              |      23           | Aircraft, engine parts
Zinc                   |      26           |
Lead                   |      29           |
Rubber (approx.)       |      77           | Much larger than metals
```

Notice how **concrete and steel have very similar coefficients**. This is why steel-reinforced concrete works so well — both materials expand and contract together, preventing cracking.

Also notice how **Invar** has an exceptionally low coefficient. It was specifically engineered for precision instruments (clocks, measuring rules) that must not change size with temperature.

---

## Area and Volume Expansion

Linear expansion applies to one-dimensional changes. For surfaces (area) and 3D objects (volume), the coefficients are simply:

```
Area Expansion:
    ΔA = β × A₀ × ΔT
    Where β = 2α  (approximately)

Volume Expansion:
    ΔV = γ × V₀ × ΔT
    Where γ = 3α  (approximately)
```

The factor of 2 for area and 3 for volume comes from the fact that each dimension expands independently. If length, width, and height each increase by α × ΔT, then:

```
New length: L₀(1 + αΔT)
New width:  W₀(1 + αΔT)
New height: H₀(1 + αΔT)

New volume: L₀W₀H₀ × (1 + αΔT)³
           ≈ V₀ × (1 + 3αΔT)   for small αΔT

So γ = 3α
```

### Practical Implication: The Ball and Ring Experiment

A classic demonstration: take a metal ball and a metal ring. At room temperature, the ball barely fits through the ring. Heat the ball and it will no longer fit. Cool the ball (or heat the ring) and the fit loosens again. This demonstrates that both the ball and the ring expand in all dimensions.

```
Room temperature:            Ball just fits through ring.
    ╔═══╗
    ║ O ║   Ball diameter ≈ ring inner diameter
    ╚═══╝

After heating the ball:      Ball no longer fits.
    ╔═══╗
    ║(O)║   Ball diameter > ring inner diameter
    ╚═══╝
    
After heating the ring:      Ball fits through again.
    ╔═════╗
    ║  O  ║   Ring inner diameter > ball diameter
    ╚═════╝
```

---

## Engineering Applications

Thermal expansion is not just a curiosity — it has major real-world consequences in engineering. Ignoring it can lead to catastrophic failures.

### Bridge Expansion Joints

Steel bridges expand significantly in summer and contract in winter. A 1000-metre bridge made of steel (α = 12×10⁻⁶/°C) that experiences a temperature range of 50°C will change length by:

```
ΔL = 12×10⁻⁶ × 1000 × 50 = 0.6 metres
```

That is 60 centimetres of movement. If the bridge were built as a single rigid structure with no room to move, it would buckle in summer and crack in winter. 

**Expansion joints** are gaps (or sliding plates) built into bridges, roads, and buildings at regular intervals. They allow the structure to move without building up stress.

```
Bridge Expansion Joint:

     Road surface
────────────────────/\/\/\/\──────────────────────
                   ↑
              Expansion joint
              (gap or flexible connector)
              
In summer:    gap closes as metal expands →
In winter:    gap opens as metal contracts ←
```

### Railway Tracks

Old railway lines were laid in sections with small gaps between them. You could hear the clickety-clack as the wheels crossed each gap. Modern railways use **continuously welded rail** — sections welded together with no gaps — but special procedures are needed. The rails are installed at a controlled temperature and are pre-stressed so that at high temperatures they are under slight compression and at low temperatures under slight tension, rather than buckling or pulling apart.

```
Old railway (with gaps):
    ════════════════╡  ╞════════════════╡  ╞════════
                   gap                gap
                   
Hot day: gap closes (expansion)
Cold day: gap opens (contraction)
```

### Hot Water Pipes and Expansion Loops

Pipes carrying hot water or steam must accommodate expansion. A 10-metre copper pipe (α = 17×10⁻⁶/°C) carrying steam at 120°C in a room at 20°C expands by:

```
ΔL = 17×10⁻⁶ × 10 × 100 = 0.017 m = 1.7 cm
```

That might not sound like much, but if the pipe is rigidly fixed at both ends, enormous compressive stress builds up. Engineers install **expansion loops** or **flexible bellows** at intervals to absorb the movement.

### Tight Fits and Shrink Fitting

Engineers use thermal expansion deliberately to assemble parts that need very tight fits. A steel shaft can be cooled in liquid nitrogen until it shrinks enough to slide into a tight-fitting bearing housing. When it warms up, it expands to lock firmly in place. This is called **shrink fitting** and creates interference fits that are stronger than any mechanical fastener.

---

## Worked Example 1

**Problem:** A railway track section is 25.0 metres long at an installation temperature of 15°C. In summer, the track reaches 50°C. In winter, it drops to -10°C. Calculate the change in length for each season. (α for steel = 12×10⁻⁶/°C)

**Summer expansion (15°C to 50°C):**

```
ΔT_summer = 50 - 15 = +35°C

ΔL_summer = α × L₀ × ΔT
ΔL_summer = 12×10⁻⁶ × 25.0 × 35
ΔL_summer = 12×10⁻⁶ × 875
ΔL_summer = 0.0105 m = 10.5 mm
```

**Winter contraction (15°C to -10°C):**

```
ΔT_winter = -10 - 15 = -25°C

ΔL_winter = 12×10⁻⁶ × 25.0 × (-25)
ΔL_winter = -0.0075 m = -7.5 mm
```

**Answer:**
- In summer, the track expands by 10.5 mm
- In winter, the track contracts by 7.5 mm
- The total seasonal variation is 10.5 + 7.5 = **18 mm per 25-metre section**

If there are 40 sections per kilometre, the total movement per kilometre would be 40 × 18 mm = 720 mm = 72 cm — nearly three-quarters of a metre. This shows why expansion gaps or pre-stressing is essential.

---

## Bimetallic Strip

A **bimetallic strip** is one of the most elegant and practically useful devices based on thermal expansion. It consists of two strips of different metals bonded firmly together along their entire length.

When the temperature changes, both metals expand — but by different amounts because they have different coefficients of expansion. Since they are bonded together and cannot separate, the strip bends.

The metal with the **higher α** expands more and ends up on the **outside** of the curve (the longer side). The metal with the **lower α** ends up on the **inside** of the curve.

```
Bimetallic Strip — Straight at room temperature:

    ████████████████████████████████
    BRASS (α = 19×10⁻⁶/°C)
    ════════════════════════════════
    STEEL (α = 12×10⁻⁶/°C)
    ════════════════════════════════
    fixed end              free end

Bimetallic Strip — Heated above room temperature:

    ████████████████████████████████████╮
    BRASS (expanded more — outer curve)   ↓ BENDS DOWN
    ════════════════════════════════════╯
    ════════════════════════════════╯
    STEEL (expanded less — inner curve)

Free end has moved DOWNWARD because brass is on top.

Bimetallic Strip — Cooled below room temperature:

    ════════════════════════════════╮
    STEEL (contracted less — outer curve) ↑ BENDS UP
    ════════════════════════════════╯
    ████████████████████████████████╮
    BRASS (contracted more — inner curve) ↑
    Free end has moved UPWARD.
```

### Applications of the Bimetallic Strip

**Thermostats:** The bending strip makes or breaks an electrical contact. When the room gets too hot, the strip bends to open the contact and shut off the heater. When the room cools, it bends back to close the contact and turn the heater on. Simple, reliable, no electronics needed.

```
Thermostat Mechanism:
                    
    Cold:                           Hot:
    
    ────── electrical contact       ─────╮
    ══════ bimetallic strip ───◯         ╰─ bimetallic strip
    
    Circuit CLOSED → heater ON          ◯  Contact broken
                                        Circuit OPEN → heater OFF
```

**Oven thermometers:** A coiled bimetallic strip rotates a needle to indicate temperature.

**Fire alarms:** A bimetallic strip in the presence of heat bends to trigger an alarm circuit.

**Electrical circuit breakers:** Excess current heats a bimetallic strip, which bends to trip the circuit breaker and protect against overload.

**Mechanical clocks:** Early precision clocks used bimetallic pendulums to compensate for temperature-induced length changes.

---

## Worked Example 2

**Problem:** A bimetallic strip consists of a brass strip (α = 19×10⁻⁶/°C) and a steel strip (α = 12×10⁻⁶/°C) each 30 cm long, bonded together at 20°C. If the temperature rises to 80°C, calculate the difference in expansion between the two strips.

**Given:**
- L₀ = 0.30 m
- ΔT = 80 - 20 = 60°C
- α_brass = 19×10⁻⁶/°C
- α_steel = 12×10⁻⁶/°C

**Expansion of brass strip:**

```
ΔL_brass = α_brass × L₀ × ΔT
ΔL_brass = 19×10⁻⁶ × 0.30 × 60
ΔL_brass = 3.42×10⁻⁴ m = 0.342 mm
```

**Expansion of steel strip:**

```
ΔL_steel = α_steel × L₀ × ΔT
ΔL_steel = 12×10⁻⁶ × 0.30 × 60
ΔL_steel = 2.16×10⁻⁴ m = 0.216 mm
```

**Difference in expansion:**

```
ΔL_brass - ΔL_steel = 0.342 - 0.216 = 0.126 mm
```

This tiny difference of 0.126 mm causes the 30 cm strip to bend noticeably. Because the brass is longer, it curves toward the steel side — the free end moves toward the steel side of the strip.

**Answer:** The brass expands 0.126 mm more than the steel, causing the strip to bend with the brass on the outside of the curve.

---

## Thermal Shock in Glass

Ordinary glass is a poor thermal conductor and has a relatively high coefficient of expansion (α ≈ 8.5×10⁻⁶/°C). When you pour boiling water into a cold ordinary glass, the inner surface heats up and expands rapidly while the outer surface is still cold. This creates internal stress. If the stress is large enough, the glass cracks — this is called **thermal shock**.

**Borosilicate glass** (sold as Pyrex or Duran) has α ≈ 3.3×10⁻⁶/°C — only about 40% of ordinary glass. This means much less thermal stress for the same temperature change, so it can withstand boiling liquids, oven temperatures, and rapid changes without shattering. This is why laboratory glassware and baking dishes are made of borosilicate glass.

**Fused quartz** goes even further with α ≈ 0.59×10⁻⁶/°C. It can be heated to red heat and plunged into cold water without cracking.

```
Thermal Shock in Glass:

    Cold glass + hot liquid:
    
    Outer surface: still cold, original size
    ┌─────────────────────────┐
    │▓▓▓▓▓ GLASS WALL ▓▓▓▓▓▓│
    │                         │
    │         HOT             │  ← Inner surface: heated, expanding
    │        LIQUID           │
    │                         │
    └─────────────────────────┘
    
    Inner surface trying to expand → TENSION on outer surface
    
    Result: CRACK if stress exceeds glass strength
    
    Solution: Use borosilicate glass (low α) or pre-warm the glass
```

---

## Anomalous Expansion of Water

Water is one of the most remarkable substances in nature, and its thermal expansion behaviour is particularly unusual. Most substances contract continuously as they cool. Water mostly does too, but with a critical exception between 0°C and 4°C.

### Water's Density vs Temperature

```
Density of Water vs Temperature:

Density
(kg/m³)
  |
1000.00 |.......................●
        |                   .     .
 999.95 |                 .          .
        |               .               .
 999.85 |             .                    .
        |           .                         .
 999.50 |         .
        |
        +────────────────────────────────────────→
          0°C   4°C   10°C  20°C  30°C  40°C
          (ice)  ↑
            Maximum density at 4°C!
```

Water is **densest at 4°C**. Below 4°C, it actually becomes LESS dense as it cools. At 0°C it freezes into ice, which is significantly less dense than liquid water.

Density of ice ≈ 917 kg/m³
Density of water at 4°C ≈ 1000 kg/m³

Ice is about 8.3% less dense than liquid water.

### Why This Happens: The Hydrogen Bond Explanation

Water molecules (H₂O) form **hydrogen bonds** with each other. In liquid water, these bonds are constantly forming and breaking, so molecules can pack fairly closely together.

As water cools below 4°C, the hydrogen bonds start to become more permanent and arranged. They force the water molecules into the hexagonal ice crystal structure — a more open arrangement than liquid water. Each water molecule ends up with more space around it in ice than in liquid water. This open structure is less dense.

```
Liquid water:                Ice crystal structure:
(disordered, compact)        (hexagonal, open)

  H₂O H₂O H₂O              H₂O
H₂O H₂O H₂O H₂O         H₂O   H₂O
  H₂O H₂O H₂O          H₂O   O   H₂O
H₂O H₂O H₂O H₂O         H₂O   H₂O
                              H₂O

More molecules per volume     Fewer molecules per volume
→ Higher density              → Lower density
```

---

## Why Ice Floats

Because ice is less dense than liquid water, it floats on the surface. This is unusual — most solids sink in their own liquid.

You can see this in an ice cube floating in a glass of water. About 8% of the ice is above the water surface and 92% is submerged. (Icebergs follow the same rule — the famous "nine-tenths below the surface.")

```
Ice floating on water:

    ┌──────────────────────┐
    │  ████████████  ←─── above surface (≈8%)
    │──────────────────────│  ← water surface
    │  ████████████████    │
    │  ████████████████    │  ← submerged (≈92%)
    │  ████████████████    │
    │                      │
    │      liquid          │
    │       water          │
    └──────────────────────┘
```

---

## Ecological Importance

The fact that ice floats — and that water is densest at 4°C — has profound consequences for life on Earth.

### How Lakes Freeze (and Do Not Freeze Solid)

Consider a lake in winter as the air temperature drops below 0°C:

```
Stage 1: Air at 10°C, lake at 15°C — uniform mixing
Stage 2: Surface cools. 10°C water is denser than 15°C water → sinks.
         Warmer water rises. Continuous mixing until the entire lake 
         approaches 4°C.
         
Stage 3: Surface cools below 4°C. Water below 4°C is LESS dense 
         than 4°C water → it stays at the top! No more convection.
         
Stage 4: Surface reaches 0°C and freezes. Ice forms at top.

Stage 5: Ice layer insulates the water below. Water below stays 
         liquid at ≈4°C. Fish and aquatic life survive.
```

```
Winter Lake Cross-Section:

    ████████████ ice at 0°C ████████████  ← floats on top
    ─────────────────────────────────────
          water at 1°C
          water at 2°C
          water at 3°C              } liquid water — fish survive here
          water at 4°C  (densest — settles to bottom)
```

If water behaved normally — becoming denser as it cooled all the way to freezing — then the coldest water would always sink to the bottom. The lake would freeze from the bottom up, and all aquatic life would perish every winter. Life in cold climates would be extremely difficult or impossible.

The anomalous expansion of water is one of the reasons liquid water has been continuously available on Earth's surface, making complex aquatic life possible.

### The Pipe-Bursting Problem

The expansion of water as it freezes (density decreases → volume increases) creates tremendous pressure in enclosed spaces. A pipe full of water that freezes can crack even thick metal pipes, because water expanding from liquid to ice increases in volume by about 9%.

This is why:
- Plumbers insulate pipes that run through cold exterior walls
- You drain outdoor hoses before winter
- Antifreeze (ethylene glycol) is added to car radiators — it prevents water from forming the low-density ice crystal structure

---

## Summary

When objects are heated, their atoms vibrate more energetically. Because the interatomic potential is asymmetric (easier to stretch than compress), atoms spend more time at larger separations, and the average spacing increases. This is **thermal expansion**.

**Linear expansion** formula: ΔL = α × L₀ × ΔT

The **coefficient of linear expansion α** varies by material. Steel: 12×10⁻⁶/°C. Copper: 17×10⁻⁶/°C. Aluminium: 23×10⁻⁶/°C. Invar: 1.2×10⁻⁶/°C.

**Area expansion**: ΔA = 2α × A₀ × ΔT

**Volume expansion**: ΔV = 3α × V₀ × ΔT

**Engineering applications**: bridge expansion joints, railway track gaps, bimetallic strip thermostats, shrink fitting, expansion loops in pipes, borosilicate glass for thermal shock resistance.

The **bimetallic strip** bends when heated because two bonded metals with different α values expand by different amounts. The higher-α metal forms the outer curve. Used in thermostats, fire alarms, and circuit breakers.

**Water's anomalous expansion**: water is densest at 4°C and becomes less dense as it cools further, freezing into ice that is 8.3% less dense than liquid water. This causes ice to float, which insulates lakes in winter and allows aquatic life to survive. It is one of the key properties that makes liquid water suitable for life.

---

## Key Equations

```
Linear expansion:
    ΔL = α × L₀ × ΔT
    L  = L₀ × (1 + α × ΔT)

Area expansion:
    ΔA = 2α × A₀ × ΔT
    β  = 2α

Volume expansion:
    ΔV = 3α × V₀ × ΔT
    γ  = 3α

Common α values (×10⁻⁶ per °C):
    Steel / Iron    : 12
    Copper          : 17
    Aluminium       : 23
    Glass (ordinary): 8.5
    Glass (Pyrex)   : 3.3
    Invar           : 1.2

Density of ice ≈ 917 kg/m³
Density of water at 4°C ≈ 1000 kg/m³
Maximum density of water at 4°C
```

---

*Next Chapter: Chapter 39 — Heat Transfer: Conduction, Convection, and Radiation*
