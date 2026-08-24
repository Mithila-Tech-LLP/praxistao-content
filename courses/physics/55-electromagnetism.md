# Chapter 55: Electromagnetism

> **"In 1820, Oersted accidentally discovered that electricity and magnetism are connected. This accident changed the world."**

---

## Table of Contents

- [Introduction: The Accidental Discovery](#introduction-the-accidental-discovery)
- [Magnetic Fields Around Current-Carrying Wires](#magnetic-fields-around-current-carrying-wires)
- [The Right-Hand Rule](#the-right-hand-rule)
- [Magnetic Field of a Straight Wire](#magnetic-field-of-a-straight-wire)
- [Magnetic Field Inside a Solenoid](#magnetic-field-inside-a-solenoid)
- [Electromagnets](#electromagnets)
- [Magnetic Domains](#magnetic-domains)
- [Permanent Magnets](#permanent-magnets)
- [Demagnetization](#demagnetization)
- [Force Between Parallel Current-Carrying Wires](#force-between-parallel-current-carrying-wires)
- [The Ampere Definition](#the-ampere-definition)
- [Magnetic Permeability](#magnetic-permeability)
- [Worked Example 1: Field Around a Wire](#worked-example-1-field-around-a-wire)
- [Worked Example 2: Field Inside a Solenoid](#worked-example-2-field-inside-a-solenoid)
- [Worked Example 3: Force Between Parallel Wires](#worked-example-3-force-between-parallel-wires)
- [Applications of Electromagnetism](#applications-of-electromagnetism)
- [MRI Machines](#mri-machines)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: The Accidental Discovery

On a spring evening in 1820, Hans Christian Oersted was preparing a lecture demonstration for his students at the University of Copenhagen. He planned to show that electricity and magnetism were completely separate phenomena — the scientific consensus at the time.

During the setup, he noticed something strange. When he switched on the electric current through a wire, the compass needle nearby **twitched**. When he switched it off, it swung back. He tried it again. Same result. The current in the wire was affecting the compass.

Oersted had stumbled onto one of the most important discoveries in the history of physics: **electricity and magnetism are not separate — they are two faces of the same force**.

This discovery launched a revolution. Within weeks, André-Marie Ampere had worked out the mathematical laws. Within decades, Michael Faraday had discovered electromagnetic induction. By 1865, James Clerk Maxwell had unified the whole picture into four elegant equations. And by the twentieth century, electromagnetism had given us electrical generators, motors, radio, television, computers, and the entire infrastructure of modern civilization.

All from a twitching compass needle.

---

## Magnetic Fields Around Current-Carrying Wires

Before Oersted, scientists knew that bar magnets created magnetic fields — regions of space where magnetic forces act. They mapped these fields using iron filings or compass needles, and found that field lines go from north pole to south pole outside the magnet.

What Oersted discovered was that a **current-carrying wire also creates a magnetic field** around it. But this field has a fundamentally different shape from a bar magnet's field.

### The Shape of the Field

When current flows through a straight wire, the magnetic field lines form **concentric circles** around the wire. Imagine looking down the length of the wire from above — the field lines are rings, like tree rings, centered on the wire.

```
         Current flowing upward (out of page)
                      ⊙
        
              ←←←←←←←←←←←←←
           ↙                   ↖
          ↓       ⊙ wire        ↑
           ↘                   ↗
              →→→→→→→→→→→→→
        
         Field lines are circles around the wire
         (shown as arrows: counterclockwise when
          current comes toward you)
```

The key features:
- The field is **circular** — it wraps around the wire
- The field is **strongest close to the wire** and gets weaker farther away
- The field direction depends on the **direction of current flow**
- Unlike bar magnet fields, these field lines have **no beginning and no end** — they are complete closed loops

---

## The Right-Hand Rule

To figure out which way the circular field points, use the **right-hand rule**:

> Point your right thumb in the direction of conventional current flow (positive charge direction). Your fingers naturally curl in the direction of the magnetic field.

```
         RIGHT HAND RULE FOR A WIRE
         
              Right hand
              
                ___
               /   \
              |  ⊙  |  ← thumb points UP (current direction)
               \___/
               
              Fingers curl
              COUNTERCLOCKWISE
              (when viewed from above)
              
              This is the field direction.
```

**Important:** Conventional current flows from positive to negative (opposite to electron flow). Always use conventional current direction when applying the right-hand rule.

If you flip the current direction:
- Current flowing downward → fingers curl clockwise when viewed from above
- The field direction reverses completely

This is why reversing current in an electromagnet reverses its poles — a fact with enormous practical importance.

---

## Magnetic Field of a Straight Wire

The strength of the magnetic field around a straight wire was worked out by Biot and Savart (pronounced "Bee-oh Savar"). Their result is:

```
    B = μ₀ × I / (2 × π × r)
```

Where:
- **B** = magnetic field strength (in Tesla, T)
- **μ₀** (mu-naught) = permeability of free space = 4π × 10⁻⁷ T·m/A
- **I** = current in the wire (Amperes)
- **r** = perpendicular distance from the wire (metres)
- **π** ≈ 3.14159

### What This Equation Tells Us

**1. Field decreases with distance:** B is inversely proportional to r. Double the distance → half the field. Ten times the distance → one-tenth the field.

**2. Field increases with current:** Double the current → double the field. This is why electromagnets work — more current, stronger field.

**3. The constant μ₀:** This is the **permeability of free space**, a fundamental constant that describes how "easy" it is for magnetic fields to form in vacuum. Its value is exactly 4π × 10⁻⁷ T·m/A (approximately 1.257 × 10⁻⁶ T·m/A).

```
         FIELD STRENGTH VS DISTANCE FROM WIRE
         
    B(T)
    |
    |*
    | *
    |  *
    |   **
    |     ***
    |        ****
    |            ******
    |                  **********
    +-------------------------------- r (m)
    0   r₁  r₂   r₃   r₄   r₅
    
    Field drops as 1/r — rapidly at first, then slower
```

---

## Magnetic Field Inside a Solenoid

A **solenoid** is a coil of wire wound into a helix (spring-like shape). When current flows through the solenoid, the circular fields from each individual loop add together to create a nearly uniform field inside the coil.

```
         SOLENOID — CROSS SECTION VIEW
         
    ___  ___  ___  ___  ___  ___  ___
   |   ||   ||   ||   ||   ||   ||   |
   | ⊗ || ⊗ || ⊗ || ⊗ || ⊗ || ⊗ || ⊗ |   <- wire cross sections
   |___||___||___||___||___||___||___|      (current going INTO page)
   
   →→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→
   →→→→  UNIFORM FIELD INSIDE  →→→→→
   →→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→→
   
    ___  ___  ___  ___  ___  ___  ___
   |   ||   ||   ||   ||   ||   ||   |
   | ⊙ || ⊙ || ⊙ || ⊙ || ⊙ || ⊙ || ⊙ |   <- wire cross sections
   |___||___||___||___||___||___||___|      (current coming OUT of page)
   
   North pole on right (field exits right end)
   South pole on left  (field enters left end)
```

The field inside a solenoid is given by:

```
    B = μ₀ × n × I
```

Where:
- **B** = magnetic field strength inside the solenoid (Tesla)
- **μ₀** = 4π × 10⁻⁷ T·m/A
- **n** = number of turns per metre (turns/m)
- **I** = current (Amperes)

### Why Is the Field Uniform Inside?

Each turn of wire contributes a small circular field. Inside the solenoid, all these circular fields point in the same direction (along the axis) and add together. Outside the solenoid, the fields from neighboring turns point in opposite directions and largely cancel out.

The result: **strong, uniform field inside** — weak, messy field outside.

This uniformity is incredibly useful. Scientists use large solenoids to create precise, stable magnetic fields for experiments. MRI machines are essentially giant solenoids.

### Increasing the Number of Turns

If you double n (turns per metre), you double B. This is why:
- A solenoid wound with 1000 turns/m produces twice the field of one wound with 500 turns/m (same current, same dimensions)
- Packing more turns into the same length increases the field

---

## Electromagnets

An **electromagnet** is a solenoid with a **ferromagnetic core** (usually iron) inside it. The iron core dramatically increases the magnetic field — sometimes by a factor of thousands.

```
         ELECTROMAGNET
         
    Battery  ────────────────────────────
    (+)  |                               |
         |    ========================   |
         |   |  [Fe core inside coil] |  |
         +→→→|========================|→→+
         |    ========================   |
    (-)  |                               |
         ────────────────────────────────
         
         South pole        North pole
              ←                 →
         
         Without iron core: field = B₀
         With iron core: field = μᵣ × B₀
         (μᵣ for soft iron ≈ 1000-5000 !!)
```

### Why Does the Iron Core Help So Much?

This is where **magnetic domains** come in.

---

## Magnetic Domains

Iron (and other ferromagnetic materials like nickel and cobalt) has a special internal structure. The atoms in iron act like tiny bar magnets — each has a magnetic moment due to the spin of its electrons.

In an unmagnetized piece of iron, these atomic magnets are organized into **domains** — regions of thousands to millions of atoms where all the atomic magnets point in the same direction. However, different domains point in different directions, so they cancel out.

```
         MAGNETIC DOMAINS IN IRON
         
         UNMAGNETIZED:
         ┌──────────────────────────────┐
         │  →→→→  │  ←←←←  │  ↑↑↑↑  │
         │  →→→→  │  ←←←←  │  ↑↑↑↑  │
         │  →→→→  │  ←←←←  │  ↑↑↑↑  │
         ├─────────┼─────────┼─────────┤
         │  ↓↓↓↓  │  →→→→  │  ←←←←  │
         │  ↓↓↓↓  │  →→→→  │  ←←←←  │
         └──────────────────────────────┘
         Random domain orientations → net field = 0
         
         MAGNETIZED (in external field →→→):
         ┌──────────────────────────────┐
         │  →→→→  │  →→→→  │  →→→→  │
         │  →→→→  │  →→→→  │  →→→→  │
         │  →→→→  │  →→→→  │  →→→→  │
         ├─────────┼─────────┼─────────┤
         │  →→→→  │  →→→→  │  →→→→  │
         │  →→→→  │  →→→→  │  →→→→  │
         └──────────────────────────────┘
         All domains aligned → strong net field →→→
```

When you place iron in an external magnetic field (like inside a solenoid), the domains **align with the external field**. This alignment massively amplifies the total magnetic field. Once the external field is removed, some iron types (**soft iron**) lose their alignment quickly. Others (**hard iron** or steel) keep it — becoming permanent magnets.

---

## Permanent Magnets

A **permanent magnet** is a material where virtually all the magnetic domains are aligned and stay that way without any external field. The material has been "hardened" so domains cannot easily rotate.

Materials used for permanent magnets:
- **Steel** (iron + carbon) — traditional horseshoe magnets
- **Alnico** (aluminum, nickel, cobalt) — strong, temperature resistant
- **Ferrite** (iron oxide compounds) — cheap, used in fridge magnets
- **Neodymium** (Nd₂Fe₁₄B) — the strongest permanent magnets available today, used in hard drives, headphones, electric motors

```
         BAR MAGNET — DOMAIN VIEW
         
         ┌─────────────────────────────────────┐
     N   │ →→→ │ →→→ │ →→→ │ →→→ │ →→→ │ →→→ │   S
         └─────────────────────────────────────┘
         
         All domains point the same way.
         North pole is where field exits (→),
         South pole is where field enters (→ points toward S from outside).
```

---

## Demagnetization

A permanent magnet can lose its magnetism if the domain alignment is disrupted:

### 1. Heating (Above the Curie Temperature)
Each ferromagnetic material has a **Curie temperature** — the temperature above which thermal energy is so great that domains spontaneously randomize.

- Iron: Curie temperature ≈ 770°C
- Nickel: ≈ 358°C
- Neodymium magnets: ≈ 310°C (surprisingly low — keep them away from heat!)

When you heat a magnet above its Curie temperature and let it cool without an external field, the domains settle randomly. The magnet is demagnetized.

### 2. Hammering or Vibration
Strong mechanical shocks can shake the domains out of alignment. Repeatedly hitting a magnet with a hammer scrambles the domains and demagnetizes it.

### 3. Alternating Magnetic Fields
Applying a strong alternating magnetic field (one that reverses direction rapidly) and slowly decreasing its strength forces the domains to oscillate and eventually settle in random orientations. This is used industrially to degauss (demagnetize) materials.

---

## Force Between Parallel Current-Carrying Wires

Here is a remarkable fact: two parallel wires carrying current exert forces on each other — even though they are not touching.

This happens because each wire creates a magnetic field, and a current-carrying conductor in a magnetic field experiences a force (the motor effect, which we will study in the next chapter).

The result:
- **Currents flowing in the same direction → wires attract each other**
- **Currents flowing in opposite directions → wires repel each other**

```
         PARALLEL WIRES — SAME DIRECTION CURRENTS
         
         Wire 1          Wire 2
           ⊙               ⊙          (both out of page)
           
           |←──── r ────→|
           
    Wire 1 creates circular field.
    At Wire 2's location, Wire 1's field points upward (↑).
    Current in Wire 2 is out of page (⊙).
    Force = I × L × B, direction given by right hand:
    fingers point up (field), thumb out of page (current)
    → palm pushes LEFT → Wire 2 attracted to Wire 1. ✓
    
         PARALLEL WIRES — OPPOSITE CURRENTS
         
         Wire 1          Wire 2
           ⊙               ⊗          (opposite directions)
           
    Same logic gives: Wire 2 pushed RIGHT → repulsion. ✓
```

The force per unit length between two parallel wires separated by distance r, each carrying currents I₁ and I₂ is:

```
    F/L = μ₀ × I₁ × I₂ / (2 × π × r)
```

---

## The Ampere Definition

The attraction/repulsion between parallel wires is so precise and reproducible that it was used to **define the Ampere** (the SI unit of current):

> **One Ampere** is the constant current that, when flowing in two infinitely long, straight parallel wires 1 metre apart in vacuum, produces a force of exactly 2 × 10⁻⁷ Newtons per metre of wire.

(Note: In 2019, the SI was revised. The Ampere is now defined via the elementary charge e = 1.602176634 × 10⁻¹⁹ C. But the old definition explains why μ₀ = 4π × 10⁻⁷ exactly.)

---

## Magnetic Permeability

**Permeability** (symbol μ) measures how easily a magnetic field can be established in a material.

**μ₀** = permeability of free space (vacuum) = 4π × 10⁻⁷ T·m/A

**Relative permeability** μᵣ = μ/μ₀ (dimensionless ratio):
- Vacuum: μᵣ = 1
- Air: μᵣ ≈ 1.0000004 (essentially 1)
- Aluminum: μᵣ ≈ 1.00002 (slightly more than 1 — paramagnetic)
- Soft iron: μᵣ ≈ 1000–5000 (ferromagnetic!)
- Mumetal (special alloy): μᵣ ≈ 100,000

When you insert an iron core (μᵣ = 1000) into a solenoid, the field inside becomes:

```
    B = μ₀ × μᵣ × n × I = μ × n × I
```

The field is multiplied by μᵣ. This is why a small current through an iron-core electromagnet can produce a field thousands of times stronger than an air-core solenoid with the same current.

---

## Worked Example 1: Field Around a Wire

**Problem:** A long straight wire carries a current of 5.0 A. Calculate the magnetic field strength at a perpendicular distance of 4.0 cm from the wire.

**Given:**
- I = 5.0 A
- r = 4.0 cm = 0.040 m
- μ₀ = 4π × 10⁻⁷ T·m/A

**Formula:**
```
    B = μ₀ × I / (2 × π × r)
```

**Substituting:**
```
    B = (4π × 10⁻⁷ × 5.0) / (2 × π × 0.040)

    B = (4π × 5.0 × 10⁻⁷) / (2π × 0.040)

    The π cancels:
    B = (4 × 5.0 × 10⁻⁷) / (2 × 0.040)

    B = (20 × 10⁻⁷) / (0.080)

    B = 250 × 10⁻⁷ T

    B = 2.5 × 10⁻⁵ T = 25 μT (microTesla)
```

**Direction:** Use the right-hand rule with your thumb pointing in the current direction to find the field direction at each point.

**Context:** Earth's magnetic field is about 50 μT. So this wire's field at 4 cm is about half of Earth's field — significant enough to deflect a nearby compass needle, just as Oersted observed!

---

## Worked Example 2: Field Inside a Solenoid

**Problem:** A solenoid is 20 cm long and has 400 turns wound on it. A current of 2.5 A passes through it. The core is air.
(a) Calculate the number of turns per metre (n).
(b) Calculate the magnetic field inside.
(c) If the air core is replaced with soft iron (μᵣ = 2000), what is the new field?

**Given:**
- Length L = 20 cm = 0.20 m
- Total turns N = 400
- I = 2.5 A
- μ₀ = 4π × 10⁻⁷ T·m/A

**Part (a): Finding n**
```
    n = N / L = 400 / 0.20 = 2000 turns/m
```

**Part (b): Air-core field**
```
    B = μ₀ × n × I
    B = (4π × 10⁻⁷) × 2000 × 2.5
    B = 4π × 10⁻⁷ × 5000
    B = 4π × 5000 × 10⁻⁷
    B = 4 × 3.1416 × 5000 × 10⁻⁷
    B = 62,832 × 10⁻⁷
    B = 6.28 × 10⁻³ T ≈ 6.3 mT (millitesla)
```

**Part (c): Iron-core field**
```
    B_iron = μᵣ × B_air = 2000 × 6.28 × 10⁻³
    B_iron = 12.56 T ≈ 12.6 T
```

**Discussion:** 12.6 Tesla is an extremely powerful magnetic field — stronger than most research electromagnets. This illustrates why iron cores are so valuable. In practice, iron saturates (domains fully align) at around 1–2 T, limiting the achievable field, but the amplification is still enormous.

---

## Worked Example 3: Force Between Parallel Wires

**Problem:** Two long parallel wires are 8.0 cm apart. Wire A carries 3.0 A and Wire B carries 5.0 A, both in the same direction. Calculate the force per unit length on Wire B.

**Given:**
- r = 8.0 cm = 0.080 m
- I₁ = 3.0 A (Wire A)
- I₂ = 5.0 A (Wire B)
- μ₀ = 4π × 10⁻⁷ T·m/A

**Formula:**
```
    F/L = μ₀ × I₁ × I₂ / (2 × π × r)
```

**Substituting:**
```
    F/L = (4π × 10⁻⁷ × 3.0 × 5.0) / (2 × π × 0.080)

    F/L = (4π × 15 × 10⁻⁷) / (2π × 0.080)

    Cancelling π:
    F/L = (4 × 15 × 10⁻⁷) / (2 × 0.080)

    F/L = (60 × 10⁻⁷) / (0.16)

    F/L = 375 × 10⁻⁷ N/m

    F/L = 3.75 × 10⁻⁵ N/m = 37.5 μN/m
```

**Direction:** Since currents are in the same direction, the force is **attractive** — Wire B is pulled toward Wire A.

**Note:** This is a tiny force. That is why ordinary parallel wires in household wiring do not visibly attract each other. However, in power equipment carrying thousands of amperes, these forces become enormous and must be engineered around.

---

## Applications of Electromagnetism

### Electric Bell

```
         ELECTRIC BELL CIRCUIT
         
    Battery → Spring contact → Electromagnet → Bell
         ↑                           |
         |              Iron striker ←
         └──── (contact breaks when striker hits bell,
                field collapses, spring pulls back,
                contact remakes → oscillates!)
```

When the circuit is closed, the electromagnet attracts the iron striker, which hits the bell. But as the striker moves, it breaks the circuit. The electromagnet loses its field. A spring pulls the striker back. The circuit closes again. The cycle repeats dozens of times per second — producing the ringing sound.

### Relay

A **relay** uses a small current through an electromagnet to switch a large current in a separate circuit. Applications:
- Automobile starter motors (small signal from key switch controls huge current to motor)
- Safety systems (low-power sensor controls high-power equipment)
- Old telephone exchanges (electromechanical switching)

### Solenoid Valve

A **solenoid valve** uses an electromagnetic solenoid to move a metal plunger that opens or closes a valve. Found in:
- Washing machines (controlling water inlet)
- Car engines (fuel injectors)
- Industrial pipelines (remote-controlled valves)

### Loudspeaker

```
         LOUDSPEAKER — CROSS SECTION
         
         Permanent magnet
               ↓
          N ──────── S
         (gap with strong radial field)
              ↑
         Voice coil (connected to amplifier output)
              ↓
         Cone (paper or plastic)
         
         Signal current → force on coil → cone vibrates → sound
```

A permanent magnet creates a radial field in a gap. A coil of wire sits in this gap. When audio signal current flows through the coil, the varying current creates a varying force, pushing the cone in and out. The cone vibrates and pushes air — creating sound waves that match the electrical signal.

---

## MRI Machines

**Magnetic Resonance Imaging** (MRI) uses superconducting electromagnets — solenoids cooled to near absolute zero with liquid helium — to create extremely strong, precise magnetic fields (typically 1.5 T to 7 T).

```
         MRI MACHINE — SIMPLIFIED CROSS SECTION
         
         ┌─────────────────────────────────────────┐
         │     Superconducting solenoid coils       │
         │   ┌─────────────────────────────────┐   │
         │   │  Liquid helium cooling jacket    │   │
         │   │ ┌─────────────────────────────┐ │   │
         │   │ │       Patient bore          │ │   │
         │   │ │    ←←← B field →→→         │ │   │
         │   │ └─────────────────────────────┘ │   │
         │   └─────────────────────────────────┘   │
         └─────────────────────────────────────────┘
         
         Field uniformity: < 1 part per million over imaging volume!
```

**How it works:**
1. The strong magnetic field aligns hydrogen nuclei (protons) in the patient's body
2. A radio-frequency pulse knocks these protons out of alignment
3. As protons relax back to alignment, they emit radio signals
4. The signals are detected and mathematically reconstructed into an image
5. Different tissues (fat, muscle, water) have different relaxation times → different image contrast

MRI produces detailed images of soft tissue without using ionizing radiation (unlike X-rays). It is one of the most important medical diagnostic tools ever developed — made possible by the solenoid and the discovery that current creates magnetic fields.

---

## Summary

- **Oersted discovered** in 1820 that a current-carrying wire deflects a compass needle, proving electricity and magnetism are linked.
- A **straight current-carrying wire** creates circular magnetic field lines centered on the wire.
- The **right-hand rule** determines the field direction: thumb = current direction, fingers = field direction.
- The field around a straight wire: **B = μ₀I / (2πr)** — field decreases with distance from wire.
- A **solenoid** (helical coil) produces a uniform field inside: **B = μ₀nI** where n = turns per metre.
- An **electromagnet** adds a ferromagnetic core to a solenoid, amplifying the field by the relative permeability μᵣ (up to thousands for iron).
- **Magnetic domains** are regions in ferromagnetic materials where atomic magnetic moments are aligned. In unmagnetized material, domains point randomly.
- **Permanent magnets** have all domains aligned and locked in place.
- **Demagnetization** occurs by heating above the Curie temperature, strong mechanical vibration, or alternating fields.
- **Parallel wires with same-direction currents attract; opposite-direction currents repel.**
- The **Ampere** was historically defined using the force between parallel wires.
- **Permeability** μ measures how easily magnetic field forms in a material; μᵣ is the ratio relative to vacuum.
- **Applications** include electric bells, relays, solenoid valves, loudspeakers, and MRI machines.

---

## Key Equations

```
    Magnetic field around straight wire:
        B = μ₀ × I / (2 × π × r)

    Magnetic field inside a solenoid (air core):
        B = μ₀ × n × I
        where n = N / L (turns per metre)

    Magnetic field inside a solenoid (with core):
        B = μ₀ × μᵣ × n × I

    Force per unit length between parallel wires:
        F / L = μ₀ × I₁ × I₂ / (2 × π × r)
        (attractive if same direction, repulsive if opposite)

    Permeability of free space:
        μ₀ = 4π × 10⁻⁷ T·m/A ≈ 1.257 × 10⁻⁶ T·m/A

    Number of turns per metre:
        n = N / L
```

**Units reminder:**
- Magnetic field B: Tesla (T)
- Current I: Ampere (A)
- Distance r, L: metre (m)
- Permeability μ₀: T·m/A (or equivalently H/m, henry per metre)
- Force F: Newton (N)
- Force per unit length F/L: N/m
