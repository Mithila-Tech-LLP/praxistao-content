# Chapter 29: Fluid Mechanics

> "Water is the driving force of all nature."
> — Leonardo da Vinci

---

## Table of Contents

- [Introduction: The Physics of Flow](#introduction-the-physics-of-flow)
- [What Is a Fluid?](#what-is-a-fluid)
- [Density](#density)
- [Pressure](#pressure)
- [Pressure and Depth](#pressure-and-depth)
- [Why Your Ears Hurt When Diving](#why-your-ears-hurt-when-diving)
- [Pascal's Principle](#pascals-principle)
- [Hydraulic Systems: Force Multipliers](#hydraulic-systems-force-multipliers)
- [Archimedes' Principle and Buoyancy](#archimedes-principle-and-buoyancy)
- [Why Steel Ships Float](#why-steel-ships-float)
- [Fluid Flow: Streamlines and Continuity](#fluid-flow-streamlines-and-continuity)
- [Bernoulli's Principle](#bernoullis-principle)
- [Applications of Bernoulli's Principle](#applications-of-bernoullis-principle)
- [Viscosity and Real Fluids](#viscosity)
- [Turbulence](#turbulence)
- [Worked Example 1: Pressure at Depth](#worked-example-1-pressure-at-depth)
- [Worked Example 2: Hydraulic Lift](#worked-example-2-hydraulic-lift)
- [Worked Example 3: Floating Ship](#worked-example-3-floating-ship)
- [Worked Example 4: Bernoulli and Pipe Flow](#worked-example-4-bernoulli-and-pipe-flow)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: The Physics of Flow

Reach over and pour yourself a glass of water. Hold the glass. Think for a moment about what just happened.

Water flowed from the container through the air into your glass. The glass sank slightly in your hand when it filled — you felt the weight of the water. The water sits still now, but if you tilted the glass, it would flow immediately. And inside that glass, the pressure at the bottom is slightly higher than at the top.

Fluids are everywhere. Blood flows through your arteries. Air flows over airplane wings. Water flows through pipes in your walls. Oil flows through the Earth's crust. The atmosphere is a fluid. The Earth's outer core is a fluid (liquid iron). Even solid ice can flow like a fluid over geological timescales (that is what glaciers are).

**Fluid mechanics** is the branch of physics that studies how fluids (liquids and gases) behave when at rest (**fluid statics**) and when in motion (**fluid dynamics**).

In this chapter, you will learn:
- What pressure is and how it varies with depth (why ears hurt when diving)
- Pascal's Principle (how a small force becomes a large force in a hydraulic system)
- Archimedes' Principle (why a steel ship floats)
- Bernoulli's Principle (why airplane wings generate lift, and why a curveball curves)
- What viscosity is, and when flow becomes turbulent

By the end, you will never look at a running tap, an airplane, or a swimming pool the same way.

---

## What Is a Fluid?

A **fluid** is any substance that flows — that continuously deforms under an applied shear (sideways) force. Both liquids and gases are fluids.

The key difference between solids and fluids:
- A **solid** has a definite shape that resists deformation. You need a large force to change its shape, and it springs back.
- A **fluid** has no fixed shape. It conforms to the shape of its container. Any shear stress, no matter how small, will cause a fluid to flow (eventually).

The difference between liquids and gases:
- A **liquid** is nearly incompressible — it has a fixed volume (density barely changes with pressure). Water at 100 atm has almost the same density as at 1 atm.
- A **gas** is highly compressible — it expands to fill any container and its density changes significantly with pressure.

For most problems in this chapter, we treat fluids as **ideal fluids**: incompressible, with no viscosity (no internal friction), and with steady, smooth flow.

---

## Density

**Density** (symbol: ρ, the Greek letter "rho") measures how much mass is packed into a given volume:

```
    ρ = m / V
    
    density = mass / volume
```

Units: kg/m³ (kilograms per cubic meter) or g/cm³

**Common densities at room temperature:**

```
    ┌────────────────────────────────┬──────────────┐
    │ Substance                      │ Density       │
    ├────────────────────────────────┼──────────────┤
    │ Air (sea level)                │ 1.2 kg/m³    │
    │ Fresh water                    │ 1,000 kg/m³  │
    │ Sea water                      │ 1,025 kg/m³  │
    │ Ice                            │ 917 kg/m³    │
    │ Aluminum                       │ 2,700 kg/m³  │
    │ Steel / iron                   │ 7,800 kg/m³  │
    │ Copper                         │ 8,960 kg/m³  │
    │ Lead                           │ 11,340 kg/m³ │
    │ Gold                           │ 19,300 kg/m³ │
    │ Mercury (liquid metal)         │ 13,600 kg/m³ │
    └────────────────────────────────┴──────────────┘
```

Notice that ice (917 kg/m³) is less dense than water (1,000 kg/m³) — which is why ice floats. This unusual property of water has enormous consequences for life on Earth: lakes freeze from the top down, not from the bottom up, and fish can survive winter under the ice.

Steel is about 7.8 times denser than water. So how can a steel ship float? We will get to that — it is one of the most satisfying explanations in physics.

---

## Pressure

**Pressure** (symbol: P) measures how much force is applied per unit area:

```
    P = F / A
    
    pressure = force / area
```

Units: Pascals (Pa) = N/m², or atmospheres (atm), or pounds per square inch (psi).

```
    1 atmosphere (atm) = 101,325 Pa ≈ 10⁵ Pa
    1 atm ≈ 14.7 psi
    1 bar = 100,000 Pa (close to 1 atm)
```

Pressure is a **scalar** quantity — it has no direction. When a fluid exerts pressure on a surface, the force it exerts is always **perpendicular** (normal) to that surface.

**Why area matters**: The same force concentrated on a small area produces much higher pressure than the same force spread over a large area.

```
    PRESSURE AND AREA DIAGRAM
    
    Snowshoe:              High heel shoe:
    
    ┌─────────────────┐         │
    │                 │         │
    │   F (weight)    │         │  F (same weight)
    │                 │         │
    └─────────────────┘         └─
    
    Large area A               Small area A
    P = F/A is small            P = F/A is LARGE
    
    You walk on top              You sink into snow!
    of snow easily
    
    A 60 kg person on one heel (area ≈ 1 cm² = 10⁻⁴ m²):
    P = (60 × 9.8) / 10⁻⁴ = 5.9 × 10⁶ Pa ≈ 58 atm!
    
    A 60 kg person on one snowshoe (area ≈ 0.05 m²):
    P = (60 × 9.8) / 0.05 = 11,760 Pa ≈ 0.12 atm
```

This is why knives are sharpened to a thin edge (small area = high pressure to cut), elephants do not sink into soft ground (huge feet spread weight over large area), and snowshoes work.

---

## Pressure and Depth

In a static fluid (one that is not flowing), pressure increases with depth. The deeper you go, the more fluid is above you pushing down.

Consider a column of fluid of height h, cross-sectional area A, and density ρ:

```
    FLUID COLUMN DIAGRAM
    
    ─────────────────────── ← P₀ (pressure at top surface)
    |                     |
    |   fluid (density ρ) |  height h
    |                     |
    ─────────────────────── ← P (pressure at bottom)
    
    The weight of fluid in this column:
    W = m × g = (ρ × V) × g = (ρ × A × h) × g
    
    This weight is spread over area A:
    Extra pressure from the fluid = W / A = ρgh
    
    Total pressure at depth h:
    P = P₀ + ρgh
```

The **pressure-depth formula** is:

```
    P = P₀ + ρ × g × h
```

Where:
- P = pressure at depth h (Pa)
- P₀ = pressure at the surface (Pa) — usually atmospheric pressure = 101,325 Pa
- ρ = density of the fluid (kg/m³)
- g = 9.8 m/s² (acceleration due to gravity)
- h = depth below the surface (m)

**Key features of this formula:**

1. Pressure increases **linearly** with depth. Every extra meter of water adds ρgh = 1000 × 9.8 × 1 = 9,800 Pa ≈ 0.097 atm ≈ 1/10 atm.

2. Pressure depends only on depth, **not** on the horizontal position. At the same depth, the pressure is the same everywhere in a connected fluid at rest — no matter if the container is a narrow tube or a vast lake.

3. This is true for **horizontal surfaces** only. The pressure is the same at all points on any horizontal plane in a connected fluid.

```
    PRESSURE IS THE SAME AT ALL POINTS AT THE SAME DEPTH:
    
    ─────────────────────────────────────────────────
    |     same P here                               |
    |         ↓                                     |
    |    ─────────────────────────────────          |
    |    ← at the same depth everywhere →           |
    |                                               |
    ─────────────────────────────────────────────────
```

---

## Why Your Ears Hurt When Diving

Your eardrum separates your outer ear canal from your middle ear. The middle ear is normally at atmospheric pressure, connected to the outside via the Eustachian tube.

When you dive underwater, the pressure of the water on your eardrum increases with depth (P = P₀ + ρgh). But the middle ear is still at atmospheric pressure. The pressure difference across your eardrum:

```
    ΔP = ρ × g × h = 1000 × 9.8 × h
    
    At h = 1 m: ΔP = 9,800 Pa ≈ 0.097 atm (barely noticeable)
    At h = 3 m: ΔP = 29,400 Pa ≈ 0.29 atm (you notice it)
    At h = 10 m: ΔP = 98,000 Pa ≈ 0.97 atm (painful!)
```

At 10 m depth, the pressure on the outside of your eardrum is nearly double the pressure on the inside. This pressure difference pushes the eardrum inward, causing pain.

To equalize, you pinch your nose and gently blow (the Valsalva maneuver) — this forces air through the Eustachian tube into the middle ear, raising the inside pressure to match outside. You hear a pop as the pressures equalize. Experienced divers do this automatically as they descend.

---

## Pascal's Principle

In 1653, the French mathematician and physicist Blaise Pascal discovered something remarkable:

**Pascal's Principle**: A pressure change applied to an enclosed fluid is transmitted **undiminished** to every point in the fluid and to the walls of the container.

In other words, if you increase the pressure at one point in a sealed fluid, the pressure increases by the same amount everywhere.

```
    PASCAL'S PRINCIPLE DIAGRAM
    
    Before:                        After adding ΔP at left:
    
    ─────────────────────          ─────────────────────
    |    |         |    |          | ΔP |         | ΔP |
    | 5atm|        | 5atm|         | 7atm|        | 7atm|
    |    |         |    |          |    |         |    |
    ─────────────────────          ─────────────────────
    
    Adding 2 atm at the             The 2 atm increase appears
    input increases every           everywhere — the pressure
    point by 2 atm.                 change is "communicated"
                                    instantly throughout the fluid.
```

This might seem obvious, but its implications are enormous. It means a **small force over a small area can create a large force over a large area**, as long as the pressure change is transmitted through a fluid.

---

## Hydraulic Systems: Force Multipliers

A **hydraulic system** uses Pascal's Principle to multiply forces.

```
    HYDRAULIC LIFT DIAGRAM
    
         F₁ (small force)          F₂ (large force)
         |                              |
         ↓                              ↓
    ─────────────              ─────────────────────────
    |         |                |                       |
    |  SMALL  |                |       LARGE           |
    | PISTON  |                |       PISTON          |
    | A₁      |────────────────| A₂                    |
    |         |   FLUID        |                       |
    ─────────────              ─────────────────────────
    
    A₁ = small area (e.g., 0.01 m²)
    A₂ = large area (e.g., 0.5 m²)
    
    The pressure at both pistons must be equal (same depth,
    connected fluid, Pascal's Principle):
    
    P₁ = P₂
    F₁/A₁ = F₂/A₂
    F₂ = F₁ × (A₂/A₁)
```

This gives us the **hydraulic force equation**:

```
    F₁/A₁ = F₂/A₂
    
    OR equivalently:
    
    F₂ = F₁ × (A₂ / A₁)
```

The force is multiplied by the ratio of the areas. A small input force produces a proportionally large output force!

**But wait — is energy created from nothing?** No. Energy is conserved. The small piston must move a large distance, while the large piston moves a small distance.

```
    ENERGY CONSERVATION IN HYDRAULICS:
    
    Work in = Work out (assuming no friction)
    F₁ × d₁ = F₂ × d₂
    
    If F₂ = 50 × F₁, then d₂ = d₁/50
    
    You push a small force over a large distance
    to lift a large force over a small distance.
    
    No energy is created — it is just redistributed.
    This is a mechanical advantage, like a lever.
```

**Applications of hydraulic systems:**

1. **Car brakes**: You push the brake pedal with perhaps 150 N. The hydraulic system multiplies this to thousands of newtons at the brake calipers on all four wheels.

2. **Car jack**: A small pump handle allows one person to lift a 1,500 kg car.

3. **Backhoe / excavator**: The operator uses small joystick forces to control massive hydraulic cylinders that can lift tons of earth.

4. **Airplane landing gear**: Powerful hydraulic cylinders retract and extend the landing gear reliably.

5. **Dentist's chair**: The smooth up-and-down motion is often hydraulic.

---

## Archimedes' Principle and Buoyancy

Around 250 BCE, the Greek mathematician Archimedes supposedly had a revelation while stepping into an overfull bath: the water that spilled out was equal in volume to the part of his body submerged. He reportedly leapt from the bath and ran through the streets of Syracuse shouting "Eureka!" (I have found it!).

The principle he discovered:

**Archimedes' Principle**: Any object completely or partially submerged in a fluid experiences an upward **buoyant force** equal to the weight of the fluid displaced by the object.

```
    BUOYANCY DIAGRAM
    
    Object submerged in fluid:
    
    ─────────────────────────────── (surface)
    |                             |
    |      ┌───────────┐          |
    |      │           │          |
    |  F_b ↑  OBJECT  ↓ Weight   |
    |      │           │          |
    |      └───────────┘          |
    |                             |
    ─────────────────────────────── (bottom)
    
    F_buoyancy = weight of fluid displaced
    F_b = ρ_fluid × V_displaced × g
    
    Weight of object: W = ρ_object × V_object × g
    
    If F_b > W: object floats (accelerates upward)
    If F_b = W: object is in neutral buoyancy (hovers)
    If F_b < W: object sinks (accelerates downward)
```

The buoyant force formula:

```
    F_b = ρ_fluid × V_displaced × g
```

Where:
- F_b = buoyant force (N)
- ρ_fluid = density of the fluid (kg/m³)
- V_displaced = volume of fluid displaced = volume of submerged part of the object (m³)
- g = 9.8 m/s²

**Where does buoyancy come from?** It is actually a pressure difference. The water pressure on the bottom of a submerged object is greater than the water pressure on the top (since the bottom is deeper). This net upward pressure force is the buoyant force.

```
    PRESSURE ORIGIN OF BUOYANCY
    
    ─────────────────────── surface
    |                     |
    |  P_top (lower)  ↓  |  depth h₁
    |    ┌──────────┐    |
    |    │  CUBE    │    |
    |    └──────────┘    |
    |  P_bottom (higher) ↑ |  depth h₂
    |                     |
    
    P_top = P₀ + ρgh₁
    P_bottom = P₀ + ρgh₂
    
    Net upward force = (P_bottom - P_top) × A
    = ρg(h₂ - h₁) × A
    = ρg × (height of cube) × A
    = ρg × V_cube
    = weight of fluid displaced ✓
```

---

## Why Steel Ships Float

Steel has a density of about 7,800 kg/m³ — nearly 8 times denser than water. If you drop a solid steel ball into water, it sinks immediately. So how does a steel ship float?

The key is that a ship is **not solid steel**. A ship is a hollow steel shell enclosing a large volume of air. The ship's average density — total mass divided by total volume (including all the air inside) — is much less than water.

```
    SHIP CROSS-SECTION DIAGRAM
    
    ─────────────────────────── (waterline)
    |                         |
    | AIR                     |  above waterline (extra volume)
    |    ┌─────────────────┐  |
    |    │  CARGO, ENGINES │  |
    |    │                 │  │
    ─────┤  STEEL HULL     ├──┤── waterline
         │                 │
         └─────────────────┘
    
    Total mass = steel + cargo + engines + (tiny) air mass
    Total volume = volume enclosed by hull (mostly air!)
    
    Average density = total mass / total volume
    
    For the ship to float, we need:
    average density < density of water
    
    OR equivalently:
    weight of ship < weight of water the ship could displace
```

A ship floats at the waterline where the buoyant force (weight of displaced water) exactly equals the ship's total weight.

**Plimsoll Line**: Large cargo ships have a "Plimsoll line" (also called a load line) painted on the hull — a series of marks showing the maximum safe depth the ship can be loaded to in different water conditions (fresh water vs. salt water, cold vs. warm water — these have different densities). Loading below the Plimsoll line risks sinking.

**Submarines**: A submarine controls its buoyancy by filling or emptying ballast tanks with water. When ballast tanks are full of water, the submarine is denser than the surrounding water and sinks. When the tanks are blown out with compressed air, the sub becomes less dense than water and rises.

---

## Fluid Flow: Streamlines and Continuity

Before we get to Bernoulli, we need to understand how fluids flow.

**Streamlines** are imaginary lines that show the path a small fluid element would follow as it moves through the flow. In **steady flow** (flow that does not change with time), streamlines are fixed in space.

```
    STREAMLINES IN A PIPE NARROWING
    
    Wide section:              Narrow section:
    
    ──────────────────────────────────────────
    ────────────────────────────────────────
    ─────────────────────────────────────────
    ──────────────────────────────────────
                              ─────────────
                              ────────────
                              ─────────────
    ──────────────────────────────────────
    ─────────────────────────────────────────
    ────────────────────────────────────────
    ──────────────────────────────────────────
    
    Streamlines are spread out        Streamlines crowd together
    Flow is slow                      Flow is fast
```

**The Continuity Equation**: For an incompressible fluid (like water) flowing through a pipe, the same volume of fluid must pass through every cross-section per second (fluid does not pile up or disappear inside the pipe).

```
    CONTINUITY EQUATION:
    
    A₁ × v₁ = A₂ × v₂
    
    (cross-sectional area × flow speed = constant)
    
    This is conservation of mass for fluid flow.
```

This means: **faster in narrow sections, slower in wide sections**.

```
    PIPE NARROWING - CONTINUITY:
    
    A₁ = 0.1 m²              A₂ = 0.02 m²
    v₁ = 2 m/s               v₂ = ?
    
    A₁v₁ = A₂v₂
    0.1 × 2 = 0.02 × v₂
    0.2 = 0.02 × v₂
    v₂ = 10 m/s
    
    The flow speeds up by 5× when the pipe area shrinks by 5×.
```

---

## Bernoulli's Principle

**Bernoulli's Principle** describes the relationship between fluid speed and pressure:

> Where the speed of a fluid is high, the pressure is low.
> Where the speed of a fluid is low, the pressure is high.

The quantitative version is **Bernoulli's Equation**:

```
    P + ½ρv² + ρgh = constant
    
    OR between two points in a flow:
    
    P₁ + ½ρv₁² + ρgh₁ = P₂ + ½ρv₂² + ρgh₂
```

Where:
- P = fluid pressure (Pa)
- ρ = fluid density (kg/m³)
- v = fluid speed (m/s)
- g = 9.8 m/s²
- h = height above some reference level (m)

**Origin of Bernoulli's Equation**: It is actually a statement of energy conservation for fluid flow. Each term has units of energy per unit volume (J/m³ = Pa):
- P: pressure energy (energy stored in compressed fluid)
- ½ρv²: kinetic energy per unit volume
- ρgh: gravitational potential energy per unit volume

As fluid speeds up (½ρv² increases), pressure (P) must decrease to keep the total constant (at the same height).

**Intuitive explanation using the pipe:**

```
    BERNOULLI IN A PIPE (horizontal, so h is the same):
    
    Wide section:              Narrow section:
    v₁ = 2 m/s                v₂ = 10 m/s
    P₁ = HIGH                 P₂ = LOW
    
    ──────────────────────────────────────────
    │  slow flow              fast flow       │
    │  HIGH pressure          LOW pressure    │
    │  P₁ ────────────────────────────► P₂   │
    ──────────────────────────────────────────
    
    P₁ + ½ρv₁² = P₂ + ½ρv₂²
    
    Since v₂ > v₁:
    ½ρv₂² > ½ρv₁²
    So P₂ must be less than P₁. ✓
```

This seems counterintuitive at first. You might expect faster flow to mean higher pressure (like a firehose). But Bernoulli is about the pressure **within** the moving fluid itself — the internal pressure decreases as the fluid speeds up.

**Derivation hint**: Think of pushing a small slug of fluid from the wide part to the narrow part. To accelerate it, the pressure behind it must be higher than the pressure in front. So the slow region (behind) has higher pressure, and the fast region (in front) has lower pressure.

---

## Applications of Bernoulli's Principle

### 1. Airplane Wings (Lift)

An airplane wing (airfoil) is curved on top and flatter on the bottom:

```
    AIRPLANE WING CROSS-SECTION (AIRFOIL):
    
            ─────────────────────
         ─────────────────────────────────
        ───────────────────────────────────────
       ──────────────────────────────────────────
    ════════════════════════════════════════════════
    ──────────────────────────────────────────────
       ──────────────────────────────────────────
    
    Air going over the top must travel a longer path
    (the curved surface) in the same time
    → it moves FASTER over the top
    → LOWER pressure on top
    
    Air on the bottom moves slower
    → HIGHER pressure on bottom
    
    Net upward pressure difference = LIFT
```

More precisely: lift in airplane wings is actually more complex (involving angle of attack, circulation, and other factors), but the pressure difference due to different air speeds is the core mechanism.

### 2. The Curveball in Baseball

A spinning baseball creates the Magnus effect — a version of Bernoulli. The spinning ball drags air around with it:

```
    CURVEBALL SPIN EFFECT (top view):
    
          spin →
          ─────
         /     \
        │   ● → │  (ball moving right, spinning clockwise)
         \     /
          ─────
    
    Top of ball: ball surface moves with the airflow
                 → air speeds up on that side
                 → LOWER pressure on top
    
    Bottom of ball: ball surface moves against the airflow
                    → air slows down on that side
                    → HIGHER pressure on bottom
    
    Net force: downward (ball curves down!)
    
    A left-handed pitcher's curveball spins differently
    and curves the other way.
```

### 3. The Venturi Effect

A Venturi meter is a device used to measure fluid flow speed in a pipe. It has a narrowing (throat) where the fluid speeds up:

```
    VENTURI METER DIAGRAM
    
         h₁        h₂ (lower)
         |          |
         ↓ fluid    ↓
    ─────────────────────────────────────────
    
    Wide:  v₁ (slow), P₁ (high), fluid level h₁
                 ↘
    Narrow:     v₂ (fast), P₂ (low), fluid level h₂
                 ↗
    Wide:  v₁ (slow), P₁ (high)
    ─────────────────────────────────────────
    
    By measuring h₁ - h₂ (the height difference
    between the two manometer tubes), we can
    calculate the flow speed.
    
    This is used in carburetors, spray nozzles,
    and industrial flow meters.
```

### 4. Airplanes, Helicopters, and Racing Cars

The same principle that creates lift on airplane wings creates **downforce** on racing car spoilers (the spoiler is an inverted wing — faster air on the bottom, lower pressure on the bottom, net downward force). Downforce pushes the car into the road for better traction at high speeds.

Helicopter rotor blades are also airfoils. By varying the pitch (angle of attack) of the rotor blades, the pilot controls lift.

### 5. Atomizer / Spray Bottle

When you squeeze a spray bottle, air blows across the top of a vertical tube that dips into the liquid. The fast air over the tube mouth creates low pressure there. Atmospheric pressure (higher) on the liquid surface pushes liquid up the tube into the fast-moving air stream, where it is atomized into fine droplets.

---

## Viscosity

So far we have treated fluids as **ideal** — no internal friction. Real fluids have **viscosity**, which is resistance to flow (internal friction between layers of fluid).

**Viscosity** (symbol: η, the Greek letter "eta") measures how "thick" or resistant a fluid is:
- High viscosity: honey, motor oil, ketchup — flow reluctantly
- Low viscosity: water, acetone, air — flow easily

```
    VISCOSITY ILLUSTRATION
    
    Low viscosity (water):          High viscosity (honey):
    
    ─────────────────────          ─────────────────────
    ════════════════════           ══════════════════════
    ==================             ====================
    ──────────────────             ──────────────────────
    
    Speed profile (arrows):         Speed profile (arrows):
    ──────── long →                 ──── short →
    ──────── →                      ──── →
    ─── short →                     ── very short →
    (faster in center,              (same shape but much
    slower near walls)              less spread)
```

The layers of fluid near the walls move slowly (the wall drags them). Layers in the center move faster. Viscosity describes how strongly adjacent layers drag on each other.

Viscosity decreases with temperature for liquids (warm honey flows much more freely than cold honey) and increases with temperature for gases.

**Units of viscosity**: Pascal-seconds (Pa·s) or Poise (P). Water at 20°C has η ≈ 10⁻³ Pa·s = 1 millipascal-second. Honey has η ≈ 2-10 Pa·s (thousands of times more viscous than water).

---

## Turbulence

Fluid flow can be of two very different types:

**Laminar flow**: Smooth, orderly flow in parallel layers (streamlines). Like honey flowing slowly from a spoon.

**Turbulent flow**: Chaotic, irregular, swirling eddies. Like water going over a waterfall or flow past a blunt object at high speed.

```
    FLOW TYPES:
    
    Laminar (smooth):          Turbulent (chaotic):
    
    ──────────────────         ~~~~~~~~~~~~~~~~~
    ──────────────────         ~~~~ eddies ~~~~
    ──────────────────         ~~whorls~swirls~
    ──────────────────         ~~~~~~~~~~~~~~~~~
    ──────────────────         ~~~~~~~~~~~~~~~~~
    
    Orderly, predictable        Disordered, energy-dissipating
```

The transition from laminar to turbulent flow is described by the **Reynolds number** (Re):

```
    Re = (ρ × v × L) / η
    
    Where:
    ρ = fluid density
    v = flow speed
    L = characteristic length (e.g., pipe diameter)
    η = viscosity
    
    Re < ~2,300: laminar flow (in a pipe)
    Re > ~4,000: turbulent flow (in a pipe)
    2,300-4,000: transitional
```

Turbulence matters enormously in engineering. Turbulent flow creates much more drag on moving objects (cars, aircraft, ships). Aircraft designers shape wings and fuselages to maintain laminar flow as long as possible. Golf balls have dimples — which seem counterintuitive, but actually cause beneficial turbulence in the boundary layer that reduces the large-scale wake drag.

---

## Worked Example 1: Pressure at Depth

**Problem**: A scuba diver is swimming at a depth of 18 m in the ocean (density of sea water = 1,025 kg/m³). Atmospheric pressure at the surface is 1.013 × 10⁵ Pa.

(a) What is the absolute pressure at 18 m depth?
(b) What is the pressure in atmospheres?
(c) The diver's tank gauge reads the pressure above atmospheric (called gauge pressure). If the absolute tank pressure is 2.0 × 10⁷ Pa, what is the gauge pressure?

**Solution**:

**(a) Absolute pressure at 18 m:**

```
    P = P₀ + ρgh
    P = 1.013 × 10⁵ + 1025 × 9.8 × 18
    P = 1.013 × 10⁵ + 180,810
    P = 100,300 + 180,810
    P = 281,110 Pa
    P ≈ 2.81 × 10⁵ Pa
```

**(b) In atmospheres:**

```
    P = 2.81 × 10⁵ Pa ÷ 1.013 × 10⁵ Pa/atm
    P ≈ 2.77 atm
```

The diver experiences about 2.77 times atmospheric pressure — almost 3 atm.

**(c) Gauge pressure:**

```
    P_gauge = P_absolute - P_atmospheric
    P_gauge = 2.0 × 10⁷ - 1.013 × 10⁵
    P_gauge ≈ 2.0 × 10⁷ - 0.1013 × 10⁶
    P_gauge ≈ 19.9 × 10⁶ Pa
    P_gauge ≈ 1.99 × 10⁷ Pa
```

**Answers**:
- Absolute pressure: 2.81 × 10⁵ Pa
- 2.77 atm
- Gauge pressure: ≈ 2.0 × 10⁷ Pa (the atmospheric contribution is less than 1%)

---

## Worked Example 2: Hydraulic Lift

**Problem**: A hydraulic car lift at a garage has a small piston with area A₁ = 8.0 cm² = 8.0 × 10⁻⁴ m², and a large piston with area A₂ = 400 cm² = 0.04 m².

(a) If a mechanic pumps the small piston with a force F₁ = 120 N, what force is produced at the large piston?
(b) What is the pressure in the hydraulic fluid?
(c) A car weighing 14,000 N is on the lift. Is the force sufficient to lift it?
(d) If the small piston is pushed down 0.25 m, how far does the large piston rise?

**Solution**:

```
    HYDRAULIC LIFT DIAGRAM:
    
    F₁ = 120 N →   F₂ = ?
         ↓                ↑
    ─────────────────────────────────────────────
    | A₁=8cm²  |     fluid      | A₂=400cm²    |
    |   pump   |────────────────|   lift pad   |
    ─────────────────────────────────────────────
```

**(a) Force at large piston:**

```
    F₁/A₁ = F₂/A₂
    F₂ = F₁ × (A₂/A₁)
    F₂ = 120 × (0.04 / 8.0 × 10⁻⁴)
    F₂ = 120 × 50
    F₂ = 6,000 N
```

**(b) Pressure in the fluid:**

```
    P = F₁/A₁ = 120 / (8.0 × 10⁻⁴)
    P = 150,000 Pa = 1.5 × 10⁵ Pa
    
    Check: P = F₂/A₂ = 6000 / 0.04 = 150,000 Pa ✓
```

**(c) Is 6,000 N enough to lift a 14,000 N car?**

```
    No! 6,000 N < 14,000 N.
    
    To lift the car, the mechanic would need:
    F₁ = F_car × (A₁/A₂) = 14,000 × (8×10⁻⁴/0.04)
    F₁ = 14,000 × 0.02
    F₁ = 280 N
    
    The mechanic would need to push with 280 N.
    Most hydraulic lifts use a pumped hydraulic system
    that builds up pressure over multiple strokes.
```

**(d) Distance the large piston rises:**

Energy conservation (or equivalently, volume conservation — the same volume moves from one side to the other):

```
    V moved = A₁ × d₁ = A₂ × d₂
    d₂ = d₁ × (A₁/A₂)
    d₂ = 0.25 × (8×10⁻⁴ / 0.04)
    d₂ = 0.25 × 0.02
    d₂ = 0.005 m = 0.5 cm
```

Each pump stroke of 25 cm raises the car by only 0.5 cm. A mechanic must pump many times to raise a car. (This is why hydraulic car jacks have a pumping handle.)

---

## Worked Example 3: Floating Ship

**Problem**: A steel ship has a total mass of 50,000 tonnes (50,000,000 kg) including hull, engines, and cargo. 

(a) What volume of seawater must the ship displace to float?
(b) If the ship's hull is 200 m long and 30 m wide (rectangular cross-section for simplicity), how deep does the ship sit in the water (the draft)?
(c) The ship's total enclosed volume (below waterline) is 70,000 m³. What fraction of the hull is below water?

**Solution**:

**(a) Volume of water displaced:**

For the ship to float, the buoyant force must equal the ship's weight:

```
    F_b = W_ship
    ρ_seawater × V_displaced × g = m_ship × g
    
    V_displaced = m_ship / ρ_seawater
    V_displaced = 50,000,000 kg / 1025 kg/m³
    V_displaced = 48,780 m³
    V_displaced ≈ 48,800 m³
```

**(b) Draft of the ship:**

With a rectangular cross-section of 200 m × 30 m = 6,000 m²:

```
    V_displaced = length × width × draft
    48,800 = 200 × 30 × draft
    48,800 = 6,000 × draft
    draft = 48,800 / 6,000
    draft ≈ 8.1 m
```

The ship sits about 8.1 m below the waterline. (Real ships have more complex hull shapes, but this gives the right order of magnitude.)

**(c) Fraction of hull below water:**

```
    Fraction = V_displaced / V_total_hull
    Fraction = 48,800 / 70,000
    Fraction = 0.697 ≈ 70%
```

About 70% of the hull is below the waterline. The famous claim that only the "tip of the iceberg" is visible is true for icebergs too: ice (density 917 kg/m³) in seawater (1025 kg/m³) floats with (917/1025) ≈ 89.5% submerged. So about 10.5% sticks above water — you see only about 1/10 of an iceberg!

---

## Worked Example 4: Bernoulli and Pipe Flow

**Problem**: Water flows through a horizontal pipe that narrows from a diameter of d₁ = 10 cm to d₂ = 4 cm. The water pressure in the wide section is P₁ = 2.0 × 10⁵ Pa and the flow speed there is v₁ = 1.5 m/s. Density of water = 1000 kg/m³.

(a) Find the flow speed v₂ in the narrow section.
(b) Find the pressure P₂ in the narrow section.
(c) Is this the static pressure increase or decrease you would expect?

**Solution**:

**(a) Flow speed in narrow section (Continuity Equation):**

```
    Area = π × (d/2)²
    
    A₁ = π × (0.10/2)² = π × (0.05)² = π × 0.0025 = 7.854 × 10⁻³ m²
    A₂ = π × (0.04/2)² = π × (0.02)² = π × 0.0004 = 1.257 × 10⁻³ m²
    
    A₁ × v₁ = A₂ × v₂
    7.854 × 10⁻³ × 1.5 = 1.257 × 10⁻³ × v₂
    1.178 × 10⁻² = 1.257 × 10⁻³ × v₂
    v₂ = (1.178 × 10⁻²) / (1.257 × 10⁻³)
    v₂ = 9.37 m/s
    
    Quick check using area ratio: A₁/A₂ = (d₁/d₂)² = (10/4)² = 6.25
    v₂ = v₁ × 6.25 = 1.5 × 6.25 = 9.375 m/s ✓
```

**(b) Pressure in narrow section (Bernoulli's Equation):**

Since the pipe is horizontal, h₁ = h₂, so the ρgh terms cancel:

```
    P₁ + ½ρv₁² = P₂ + ½ρv₂²
    P₂ = P₁ + ½ρ(v₁² - v₂²)
    P₂ = 2.0 × 10⁵ + ½ × 1000 × ((1.5)² - (9.37)²)
    P₂ = 2.0 × 10⁵ + 500 × (2.25 - 87.8)
    P₂ = 2.0 × 10⁵ + 500 × (-85.55)
    P₂ = 2.0 × 10⁵ - 42,775
    P₂ = 200,000 - 42,775
    P₂ = 157,225 Pa
    P₂ ≈ 1.57 × 10⁵ Pa
```

**(c) Interpretation:**

The pressure dropped from 2.0 × 10⁵ Pa to 1.57 × 10⁵ Pa — a decrease of about 42,775 Pa (about 42% of an atmosphere). The flow speed increased from 1.5 m/s to 9.37 m/s. As predicted by Bernoulli: faster flow → lower pressure. This is precisely the principle behind Venturi flow meters.

---

## Summary

- A **fluid** (liquid or gas) takes the shape of its container and flows under any shear force.

- **Density**: ρ = m/V. Water has density 1,000 kg/m³; steel is about 7,800 kg/m³.

- **Pressure**: P = F/A. Units: Pascals (Pa) = N/m². Small area → large pressure for same force.

- **Pressure with depth**: P = P₀ + ρgh. Pressure increases linearly with depth. At the same depth, pressure is the same everywhere in a connected fluid.

- Your ears hurt when diving because the water pressure increases with depth, creating a pressure difference across your eardrum.

- **Pascal's Principle**: Pressure applied to an enclosed fluid is transmitted undiminished to all points. This enables hydraulic systems to multiply forces: F₁/A₁ = F₂/A₂.

- **Archimedes' Principle**: Buoyant force = weight of fluid displaced = ρ_fluid × V_displaced × g. Objects float if their average density is less than the fluid density.

- Steel ships float because they enclose large volumes of air — their average density (total mass / total volume) is less than water.

- **Continuity Equation**: A₁v₁ = A₂v₂. Fluid flows faster in narrow sections, slower in wide sections.

- **Bernoulli's Equation**: P + ½ρv² + ρgh = constant. Faster flow → lower pressure. Explains airplane lift, curveballs, Venturi meters, spray bottles.

- **Viscosity** is internal friction in a fluid. Liquids become less viscous at higher temperatures.

- **Turbulence** is chaotic flow that occurs at high Reynolds numbers (Re > ~4,000 for pipe flow). It dissipates more energy than laminar flow.

---

## Key Equations

```
    DENSITY:
    ρ = m / V
    Units: kg/m³
    
    ─────────────────────────────────────────
    
    PRESSURE:
    P = F / A
    Units: Pa = N/m²
    1 atm = 101,325 Pa
    
    ─────────────────────────────────────────
    
    PRESSURE WITH DEPTH:
    P = P₀ + ρgh
    
    ─────────────────────────────────────────
    
    PASCAL'S PRINCIPLE (HYDRAULIC SYSTEMS):
    P₁ = P₂
    F₁/A₁ = F₂/A₂
    F₂ = F₁ × (A₂/A₁)
    
    ─────────────────────────────────────────
    
    ARCHIMEDES' PRINCIPLE (BUOYANCY):
    F_b = ρ_fluid × V_displaced × g
    
    Float: ρ_object < ρ_fluid
    Sink:  ρ_object > ρ_fluid
    
    ─────────────────────────────────────────
    
    CONTINUITY EQUATION:
    A₁ × v₁ = A₂ × v₂
    (flow rate = constant for incompressible fluid)
    
    ─────────────────────────────────────────
    
    BERNOULLI'S EQUATION:
    P₁ + ½ρv₁² + ρgh₁ = P₂ + ½ρv₂² + ρgh₂
    
    Horizontal pipe (h₁ = h₂):
    P₁ + ½ρv₁² = P₂ + ½ρv₂²
    
    Faster flow → lower pressure
    
    ─────────────────────────────────────────
    
    REYNOLDS NUMBER (LAMINAR vs TURBULENT):
    Re = ρvL / η
    
    Re < 2300: laminar flow
    Re > 4000: turbulent flow
```
