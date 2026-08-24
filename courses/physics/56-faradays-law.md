# Chapter 56: Faraday's Law of Induction

> **"Faraday's greatest discovery: you don't need a battery to make electricity. A changing magnetic field is enough. This single insight powers the entire modern world."**

---

## Table of Contents

- [Introduction: The Discovery That Powers the World](#introduction-the-discovery-that-powers-the-world)
- [Faraday's Original Experiments](#faradays-original-experiments)
- [Electromagnetic Induction — The Key Idea](#electromagnetic-induction--the-key-idea)
- [Magnetic Flux](#magnetic-flux)
- [Faraday's Law of Induction](#faradays-law-of-induction)
- [Factors Affecting Induced EMF](#factors-affecting-induced-emf)
- [The AC Generator](#the-ac-generator)
- [Sinusoidal Output of a Generator](#sinusoidal-output-of-a-generator)
- [Worked Example 1: EMF from Changing Flux](#worked-example-1-emf-from-changing-flux)
- [Worked Example 2: EMF of a Rotating Coil](#worked-example-2-emf-of-a-rotating-coil)
- [Worked Example 3: Number of Turns Required](#worked-example-3-number-of-turns-required)
- [Back-EMF in Electric Motors](#back-emf-in-electric-motors)
- [Eddy Currents](#eddy-currents)
- [Applications of Eddy Currents](#applications-of-eddy-currents)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: The Discovery That Powers the World

By 1820, scientists knew that electricity could create magnetism (Oersted's discovery). The obvious question was: can magnetism create electricity?

Michael Faraday, a bookbinder's apprentice who had taught himself science, spent over ten years investigating this question. He tried holding a magnet near a wire. Nothing happened. He tried connecting a battery to a coil next to another coil and measuring current in the second. Nothing. He was about to give up.

Then, in 1831, he noticed something. When he connected the battery — in that brief instant — the galvanometer in the second coil twitched. When he disconnected it, the needle twitched again. Not when the current was steady. Only when it **changed**.

Within days, he had done dozens of experiments. He found that moving a magnet in and out of a coil caused current. Not a stationary magnet. A **moving** one.

Faraday's insight: it is not the magnetic field itself, but the **changing** magnetic field, that produces electricity.

This discovery — **electromagnetic induction** — is the operating principle behind every electrical generator in the world. Every power station, every wind turbine, every hydroelectric dam uses Faraday's law to convert mechanical energy into electrical energy. Every transformer uses it. Every microphone uses it. The entire modern electrical grid exists because of what Faraday discovered on a rainy day in London in 1831.

---

## Faraday's Original Experiments

### Experiment 1: Moving Magnet

```
         FARADAY'S MAGNET EXPERIMENT
         
         Galvanometer
              G
             /|\
            / | \
           /  |  \
          ----+----
               |
              _|_
             /   \        Move magnet
            |     |     → into coil → deflection
            | N S |     → out of coil → opposite deflection
            |     |     → hold still → no deflection
             \___/
              ↕
             Coil
```

**Observation:** The galvanometer deflects only when the magnet is moving. The faster the motion, the larger the deflection.

### Experiment 2: Changing Current in Adjacent Coil

```
         TWO COILS — MUTUAL INDUCTION
         
         Primary coil              Secondary coil
              ↓                         ↓
    [Battery]──[Switch]──[Coil 1]    [Coil 2]──[Galvanometer G]
    
    (Coils placed side by side)
    
    Switch closes  → G deflects briefly → returns to zero
    Switch stays closed → G reads zero (steady current)
    Switch opens   → G deflects briefly in opposite direction
```

**Observation:** Changing current in Coil 1 changes its magnetic field. That changing field "reaches" Coil 2 and induces a brief current there.

### Experiment 3: Rotating Coil in Steady Field

Faraday found that rotating a coil in a fixed magnetic field produced a continuous current — the basic principle of the electric generator.

---

## Electromagnetic Induction — The Key Idea

The fundamental principle that Faraday discovered is:

> **A changing magnetic field induces an electromotive force (EMF) in a conductor.**

An **electromotive force (EMF)** is not truly a force — it is a voltage (measured in Volts) that drives current around a circuit, like a battery does. The key word is **changing**:

| Situation | Induced EMF? |
|-----------|-------------|
| Steady magnet near coil | No |
| Steady current in nearby coil | No |
| Moving magnet | Yes |
| Changing current in nearby coil | Yes |
| Coil rotating in steady field | Yes |
| Coil stationary in steady field | No |

The pattern is clear: whenever the **magnetic flux through a coil is changing**, there is an induced EMF.

---

## Magnetic Flux

Before we can state Faraday's Law precisely, we need to define **magnetic flux** (symbol Φ, pronounced "Phi").

Magnetic flux through a surface is a measure of "how much magnetic field passes through" that surface.

```
         MAGNETIC FLUX — VISUALIZING IT
         
         Case 1: B perpendicular to coil plane
         
              B →→→→→→→→→→→→
                ┌──────────┐
              B →→→→→→→→→→→→ →   (all field lines pass through)
                └──────────┘
              B →→→→→→→→→→→→
              
              Area = A, B uniform, angle = 0°
              Φ = B × A × cos(0°) = B × A    [maximum flux]
         
         Case 2: B parallel to coil plane
         
              B →→→→ ──────── →→→→
              B →→→→ |      | →→→→
              B →→→→ |      | →→→→
              B →→→→ ──────── →→→→
              
              (field lines skim across, none pass through)
              Φ = B × A × cos(90°) = 0       [zero flux]
         
         Case 3: B at angle θ to normal
         
              B →→→→\  ┌──────────┐
              B →→→→ \ │          │
              B →→→→  \│ θ        │
              B →→→→   └──────────┘
              
              Φ = B × A × cos(θ)             [general case]
```

**Magnetic flux formula:**
```
    Φ = B × A × cos(θ)
```

Where:
- **Φ** = magnetic flux (unit: Weber, Wb = T·m²)
- **B** = magnetic field strength (Tesla)
- **A** = area of the surface (m²)
- **θ** = angle between field direction and the **normal** (perpendicular) to the surface

When B is perpendicular to the surface (θ = 0°), flux is maximum.
When B is parallel to the surface (θ = 90°), flux is zero — field lines don't "cut through" the surface.

**Key insight:** Flux can change because B changes, because A changes, or because θ changes. All three methods can induce EMF.

---

## Faraday's Law of Induction

With flux defined, we can state **Faraday's Law** precisely:

```
    EMF = -N × ΔΦ / Δt
```

Where:
- **EMF** = induced electromotive force (Volts)
- **N** = number of turns in the coil
- **ΔΦ** = change in magnetic flux through one turn (Webers)
- **Δt** = time over which the change occurs (seconds)
- The **minus sign** comes from Lenz's Law (explained in Chapter 57)

In words: **The induced EMF equals the rate of change of total magnetic flux linkage (NΦ) through the coil.**

### What the Minus Sign Means

The minus sign is not just bookkeeping — it has physical meaning. It says the induced EMF acts in the direction that **opposes** the change in flux. This is Lenz's Law, covered in detail in Chapter 57.

### Why N Matters

A coil with N turns is like N single loops stacked together. Each turn has flux Φ passing through it. The total **flux linkage** is NΦ. A changing NΦ induces an EMF proportional to N. This is why generators and transformers use many-turn coils.

---

## Factors Affecting Induced EMF

From EMF = -N × ΔΦ/Δt, and Φ = B × A × cosθ, we can identify all the factors:

### 1. Speed of Change
```
    EMF = -N × ΔΦ / Δt
```
The faster the flux changes (smaller Δt), the larger the EMF. Move the magnet faster → bigger deflection. This is why fast-rotating generators produce more power.

### 2. Magnetic Field Strength (B)
Larger B → larger flux Φ → larger ΔΦ for same movement → larger EMF. Strong permanent magnets or powerful electromagnets increase generator output.

### 3. Number of Turns (N)
More turns → proportionally more EMF for the same rate of flux change. Doubling the turns doubles the EMF. This is the main design variable in generators and transformers.

### 4. Area of the Coil (A)
Larger coil area → larger flux → larger ΔΦ per rotation → larger EMF. But bigger coils are heavier and harder to spin — an engineering trade-off.

### 5. Angle of the Coil (θ)
Flux = B×A×cosθ. As a coil rotates, θ changes continuously. The rate of change of flux is largest when the coil is parallel to the field (θ = 90° — field passes through edge-on), and zero when the coil is perpendicular to the field. This produces sinusoidal AC output.

---

## The AC Generator

The **AC generator** (alternator) is the most important application of Faraday's Law. Every power station in the world uses this principle.

```
         AC GENERATOR — SCHEMATIC
         
                       N pole
                        |||
    External          __|_|__
    circuit          |       |
    (load)     ←     |  COIL |  ← rotates on axle
                    |_______|
    [~]              __|_|__
    ↑                   |||
    EMF output         S pole
    
    Slip rings and brushes
    carry current to/from rotating coil
    
         CUTAWAY VIEW — COIL POSITIONS
         
    Position 1:         Position 2:
    (coil ⊥ to B)       (coil ∥ to B)
    
      N         S         N    ↕     S
      |  ——|——  |         |   coil   |
      |  |   |  |         |  edge-on |
      |  ——|——  |         |          |
      
    Flux = max           Flux = 0
    Rate of change = 0   Rate of change = max
    EMF = 0              EMF = maximum
```

### The Coil Rotation

As the coil rotates at constant angular velocity ω (radians per second):
- The angle θ between the field and the coil normal changes: θ = ωt
- Flux: Φ = B × A × cos(ωt)
- EMF = -N × dΦ/dt = N × B × A × ω × sin(ωt)

This gives a sinusoidal (sine wave) output — the basis of **alternating current (AC)**.

---

## Sinusoidal Output of a Generator

```
         AC GENERATOR OUTPUT
         
    EMF
    (V)
     |
  +E₀|    *         *         *
     |  *   *     *   *     *   *
     | *     *   *     *   *     *
     |*       * *       * *       *
   --+---------------------------- time (s)
     |*       * *       * *       *→
     | *     *   *     *   *     *
     |  *   *     *   *     *   *
  -E₀|    *         *         *
     |
     
     T = period (time for one complete cycle)
     f = 1/T = frequency (cycles per second, Hz)
     ω = 2πf = angular frequency (radians/s)
     E₀ = peak EMF = N × B × A × ω
```

The output voltage varies as:
```
    EMF(t) = E₀ × sin(ωt)
    E₀ = N × B × A × ω   (peak EMF)
```

When does EMF = maximum? When the coil plane is **parallel** to the field (coil is edge-on to field lines — maximum rate of flux change).

When does EMF = zero? When the coil is **perpendicular** to the field (coil face-on — flux is maximum but rate of change is momentarily zero).

Note the subtle point: maximum flux → zero EMF. Maximum rate of change of flux → maximum EMF. It is the **rate of change** that matters, not the absolute value.

---

## Worked Example 1: EMF from Changing Flux

**Problem:** A circular coil has 200 turns and an area of 50 cm². It is placed perpendicular to a magnetic field. The field increases uniformly from 0.10 T to 0.50 T in 0.080 seconds. Calculate the induced EMF.

**Given:**
- N = 200 turns
- A = 50 cm² = 50 × 10⁻⁴ m² = 5.0 × 10⁻³ m²
- B₁ = 0.10 T (initial)
- B₂ = 0.50 T (final)
- Δt = 0.080 s
- θ = 0° (coil perpendicular to field → cos 0° = 1)

**Finding ΔΦ:**
```
    Φ₁ = B₁ × A = 0.10 × 5.0 × 10⁻³ = 5.0 × 10⁻⁴ Wb
    Φ₂ = B₂ × A = 0.50 × 5.0 × 10⁻³ = 25.0 × 10⁻⁴ Wb

    ΔΦ = Φ₂ - Φ₁ = (25.0 - 5.0) × 10⁻⁴ = 20.0 × 10⁻⁴ Wb = 2.0 × 10⁻³ Wb
```

**Applying Faraday's Law:**
```
    |EMF| = N × |ΔΦ / Δt|
    |EMF| = 200 × (2.0 × 10⁻³) / 0.080
    |EMF| = 200 × 0.025
    |EMF| = 5.0 V
```

**Answer:** The induced EMF is 5.0 V (magnitude; the sign/direction is determined by Lenz's Law).

**Check:** This is a realistic EMF. A hand-wound coil in a changing field can easily produce a few volts.

---

## Worked Example 2: EMF of a Rotating Coil

**Problem:** A rectangular coil has 150 turns, length 12 cm and width 8.0 cm. It rotates at 50 revolutions per second (50 Hz) in a uniform magnetic field of 0.20 T. Calculate:
(a) The peak EMF
(b) The EMF at the instant when the coil plane makes an angle of 30° with the field

**Given:**
- N = 150
- Length = 12 cm = 0.12 m, Width = 8.0 cm = 0.080 m
- f = 50 Hz → ω = 2πf = 2π × 50 = 100π ≈ 314 rad/s
- B = 0.20 T

**Part (a): Peak EMF**

First, find the coil area:
```
    A = length × width = 0.12 × 0.080 = 9.6 × 10⁻³ m²
```

Peak EMF:
```
    E₀ = N × B × A × ω
    E₀ = 150 × 0.20 × 9.6 × 10⁻³ × 100π
    E₀ = 150 × 0.20 × 9.6 × 10⁻³ × 314.16
    E₀ = 150 × 0.20 × 3.016
    E₀ = 150 × 0.6032
    E₀ = 90.5 V
```

**Part (b): EMF when coil plane makes 30° with field**

Careful with angles! If the coil **plane** makes 30° with the field direction, then the **normal** to the coil makes 90° - 30° = 60° with the field.

The EMF at angle θ (measured from the zero-EMF position, i.e., coil normal parallel to field):
```
    EMF = E₀ × sin(θ_from_normal)

    But we want the angle from the coil plane.
    If plane is at 30° to B, then normal is at 60° to B.

    EMF = E₀ × sin(60°)
    EMF = 90.5 × sin(60°)
    EMF = 90.5 × 0.866
    EMF = 78.4 V ≈ 78 V
```

**Discussion:** At this position the coil plane is nearly parallel to the field (30° away), so we are near maximum EMF. The maximum occurs when the coil plane is exactly parallel to B (the coil is edge-on to the field).

---

## Worked Example 3: Number of Turns Required

**Problem:** An engineer wants a generator that produces a peak EMF of 240 V. The coil rotates at 3000 rpm in a field of 0.15 T. The coil area is 200 cm². How many turns are needed?

**Given:**
- E₀ = 240 V
- Rotation speed = 3000 rpm = 3000/60 = 50 rev/s → f = 50 Hz
- ω = 2π × 50 = 100π rad/s ≈ 314 rad/s
- B = 0.15 T
- A = 200 cm² = 200 × 10⁻⁴ m² = 0.020 m²

**Rearranging peak EMF formula for N:**
```
    E₀ = N × B × A × ω
    
    N = E₀ / (B × A × ω)
    N = 240 / (0.15 × 0.020 × 314.16)
    N = 240 / (0.15 × 6.283)
    N = 240 / 0.9425
    N = 254.7
```

**Answer:** The engineer needs approximately **255 turns** (rounding up to the next whole number).

**Practical note:** Real generators use many more turns (thousands) but with smaller area coils, or use stronger magnetic fields from powerful electromagnets powered by part of the generator's output — these are called **exciter** systems.

---

## Back-EMF in Electric Motors

When a motor is running, the rotating coil inside is doing exactly what a generator coil does — moving in a magnetic field. This means it **generates its own EMF**. This self-generated voltage opposes the supply voltage (by Lenz's Law) and is called the **back-EMF** (or back-electromotive force).

```
         MOTOR CIRCUIT WITH BACK-EMF
         
    Supply voltage V_supply
    ↓
    ┌──────────────────────────────┐
    │                              │
    │   V_supply = V_R + V_back    │
    │                              │
    │   V_supply = I×R + E_back    │
    │                              │
    │   I = (V_supply - E_back)/R  │
    └──────────────────────────────┘
    
    R = resistance of coil windings
    E_back = back-EMF (proportional to motor speed)
```

### At Startup
When the motor first starts, it is not rotating, so back-EMF = 0.

```
    I_startup = V_supply / R   (no back-EMF to oppose current)
```

This is the **maximum current the motor ever draws**. For a large motor, this startup current can be 5–10 times the normal running current. This is why:
- Large motors need "soft starters" or star-delta starters to limit startup current
- Your house lights dim briefly when a large motor (refrigerator compressor) starts
- Fuses are chosen to survive brief startup surges

### At Full Speed
As the motor speeds up, back-EMF increases, reducing current:

```
    I_running = (V_supply - E_back) / R

    Since E_back is close to V_supply, I_running is small.
```

### When the Motor is Stalled
If the motor is mechanically stopped (e.g., jammed), back-EMF drops to zero. Current surges. Without protection, the motor burns out. This is why electric motors need fuses or thermal overload protection.

---

## Eddy Currents

When a **solid conductor** (not just a wire) is in a changing magnetic field, currents are induced throughout the bulk of the material. These are called **eddy currents** because they circulate in swirling patterns like eddies in water.

```
         EDDY CURRENTS IN A METAL PLATE
         
         Changing B field (decreasing)
                 ↓↓↓↓↓↓
         ┌──────────────────────┐
         │      ↗→→→↘          │
         │     ↑      ↓         │
         │      ↖←←←↙          │  ← eddy current loop
         │                      │
         │      ↗→→→↘          │
         │     ↑      ↓         │
         │      ↖←←←↙          │  ← another loop
         └──────────────────────┘
         
         These circulating currents dissipate energy as heat (I²R)
         They also create their own magnetic field opposing the change
```

**Eddy currents always oppose the change causing them** (Lenz's Law). This has both useful and unwanted consequences.

### Unwanted Effects
In **transformer cores**, eddy currents would waste enormous energy as heat. This is prevented by making the core from many thin **laminated sheets** separated by thin insulating layers. Eddy currents cannot flow between laminations, so they are restricted to tiny thin loops with high resistance — much less current, much less heating.

```
         LAMINATED CORE — CROSS SECTION
         
    Solid core:          Laminated core:
    
    ┌──────────┐         ┌─┬─┬─┬─┬─┬─┬─┬─┐
    │          │         │ │ │ │ │ │ │ │ │
    │  Large   │         │ │ │ │ │ │ │ │ │
    │  eddy    │         │ │ │ │ │ │ │ │ │
    │  current │         │ │ │ │ │ │ │ │ │
    │  loop    │         │ │ │ │ │ │ │ │ │
    └──────────┘         └─┴─┴─┴─┴─┴─┴─┴─┘
    (high current,       (tiny loops, high R
    large energy loss)   much less energy loss)
```

---

## Applications of Eddy Currents

### Induction Cooktops

```
         INDUCTION COOKTOP
         
    ┌─────────────────────────────────┐
    │    Glass ceramic surface         │
    ├─────────────────────────────────┤
    │    Copper coil below surface     │
    │    ~~~~~~~~~~~~~~~~~~~~~         │
    │    AC current at ~25 kHz         │
    └─────────────────────────────────┘
              ↕ changing magnetic field
    ┌─────────────────────────────────┐
    │    Iron/steel pot base           │
    │    (eddy currents induced here)  │
    │    → heats up from I²R           │
    └─────────────────────────────────┘
```

The rapidly alternating current (25,000 Hz) in the coil creates a rapidly changing magnetic field. This induces powerful eddy currents in the **metal pot**, heating it directly. The glass surface stays cool. Advantages: very efficient (no heat lost to surroundings), instant response, safe surface.

**Important:** Induction cooktops only work with magnetic materials (iron, steel). Aluminum, copper, and glass pots won't work on induction — they don't support strong eddy currents at useful frequencies.

### Metal Detectors

A metal detector sends out an alternating electromagnetic field. When a metal object enters this field, eddy currents are induced in the metal. These eddy currents create their own secondary magnetic field, which the detector picks up. The device signals that metal is present.

Different metals respond differently (they have different conductivities and magnetic properties), allowing sophisticated detectors to distinguish iron, gold, silver, and other metals.

### Magnetic Braking

Eddy currents create braking forces that oppose motion — useful where friction-based brakes would wear out.

```
         EDDY CURRENT BRAKE
         
         Rotating metal disk
         
                ___
               /   \
              | ↻   |  ← spinning
               \___/
         
         Electromagnet (not touching disk)
         
                N
               |||
           ───────────   ← pole pieces close to disk
               |||
                S
         
         Disk passes through field → eddy currents induced
         → eddy currents create force opposing motion → braking
         No physical contact! No wear!
```

Applications:
- Roller coaster braking systems
- Train braking (eddy current track brakes)
- Fairground ride controls
- Laboratory balances (damping oscillations)

### Induction Heating

Large eddy currents can heat metal objects rapidly and precisely. Used in:
- Melting metals in induction furnaces (no contact, no contamination)
- Hardening metal parts (heat just the surface, quench quickly)
- Welding (induction seam welding for metal pipes)
- Medical applications (heating cancerous tissue — hyperthermia treatment)

### Electromagnetic Damping in Meters

Old analog voltmeters and ammeters used a coil moving in a magnetic field. When the pointer swings toward the correct reading, eddy currents in a damping vane slow it down so it settles at the right value without oscillating. Without this damping, the pointer would swing back and forth for seconds before settling.

---

## Summary

- **Electromagnetic induction** is the production of an EMF by a changing magnetic field — discovered by Faraday in 1831.
- **Magnetic flux** Φ = B × A × cosθ measures how much magnetic field passes through a surface (unit: Weber, Wb).
- **Faraday's Law:** EMF = -N × ΔΦ/Δt — the induced EMF is proportional to the rate of change of magnetic flux and the number of turns.
- EMF is greater when: the flux changes faster, B is stronger, more turns are used, or the coil area is larger.
- An **AC generator** (alternator) rotates a coil in a magnetic field, changing flux continuously to produce a sinusoidal AC voltage.
- Peak EMF of a generator: **E₀ = N × B × A × ω**.
- **Back-EMF** is the EMF generated by a running motor opposing the supply voltage. It limits running current but disappears on startup, causing high startup currents.
- **Eddy currents** are circulating currents induced in bulk conductors by changing magnetic fields. They oppose the change causing them (Lenz's Law).
- Eddy currents in transformer cores are reduced by using **laminated cores**.
- Useful applications of eddy currents: induction cooktops, metal detectors, magnetic braking, induction heating.

---

## Key Equations

```
    Magnetic flux:
        Φ = B × A × cos(θ)
        (θ = angle between B and the normal to the area)

    Faraday's Law of Induction:
        EMF = -N × ΔΦ / Δt
        |EMF| = N × |ΔΦ / Δt|   (magnitude only)

    Peak EMF of an AC generator:
        E₀ = N × B × A × ω
        where ω = 2πf (angular frequency)

    Instantaneous EMF of generator:
        EMF(t) = E₀ × sin(ωt)

    Angular frequency:
        ω = 2 × π × f
        (f = frequency in Hz, ω in rad/s)

    Motor current with back-EMF:
        I = (V_supply - E_back) / R
```

**Units:**
- Magnetic flux Φ: Weber (Wb) = T·m²
- EMF: Volt (V)
- Angular frequency ω: rad/s
- Frequency f: Hertz (Hz) = cycles per second
