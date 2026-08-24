# Chapter 24: Universal Gravitation

> "Nature and Nature's laws lay hid in night. God said, 'Let Newton be!' and all was light."
> — Alexander Pope

---

## Table of Contents

- [Introduction: The Apple and the Moon](#introduction-the-apple-and-the-moon)
- [The Apple-Moon Connection](#the-apple-moon-connection)
- [Newton's Law of Universal Gravitation](#newtons-law-of-universal-gravitation)
- [The Gravitational Constant G](#the-gravitational-constant-g)
- [The Inverse-Square Law](#the-inverse-square-law)
- [Comparing Gravitational and Electric Forces](#comparing-gravitational-and-electric-forces)
- [Surface Gravity: Where Does g = 9.8 m/s² Come From?](#surface-gravity-where-does-g--98-ms-come-from)
- [Gravity Varies with Altitude](#gravity-varies-with-altitude)
- [Weight on Different Planets](#weight-on-different-planets)
- [Escape Velocity](#escape-velocity)
- [Worked Example 1: Force Between Two People](#worked-example-1-force-between-two-people)
- [Worked Example 2: Gravity at Altitude](#worked-example-2-gravity-at-altitude)
- [Worked Example 3: Weight on the Moon](#worked-example-3-weight-on-the-moon)
- [Worked Example 4: Escape Velocity from Earth](#worked-example-4-escape-velocity-from-earth)
- [Common Mistakes](#common-mistakes)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: The Apple and the Moon

The story goes like this. It is 1666. A young Isaac Newton, 23 years old, has fled Cambridge because of the plague and is sitting under an apple tree at his mother's farm in Woolsthorpe, Lincolnshire. An apple falls. Newton looks up at the Moon.

A thought strikes him: what if the same force pulling the apple to the ground is also pulling the Moon toward the Earth? What if gravity is not a local phenomenon — something that only happens near the Earth's surface — but a universal force that reaches out to the Moon, to the planets, to the stars?

This single idea, if correct, would unify the heavens and the Earth under one law. And Newton proved that it was correct.

The result — **Newton's Law of Universal Gravitation** — is one of the greatest achievements in the history of science. It explained:
- Why apples fall
- Why the Moon orbits the Earth
- Why the planets orbit the Sun
- Why tides rise and fall
- The precise orbits of every planet in the solar system

All from a single equation. Let us understand it completely.

---

## The Apple-Moon Connection

### The Crucial Question

When the apple falls, it accelerates toward Earth at 9.8 m/s². The Moon also "falls" toward Earth — it is constantly being pulled by Earth's gravity and curves toward the Earth (that is what an orbit is, as we will see in Chapter 25).

But the Moon is much farther away. Newton asked: how does gravity weaken with distance?

### The Numbers

```
Distance from Earth's centre to the surface:  r_Earth = 6,370 km
Distance from Earth's centre to the Moon:     r_Moon  = 384,000 km

Ratio: r_Moon / r_Earth = 384,000 / 6,370 ≈ 60
```

So the Moon is about 60 times farther from Earth's centre than the apple.

If gravity falls off as 1/r² (the **inverse-square law**), then the Moon's gravitational acceleration should be:

```
a_Moon = a_apple / 60²
       = 9.8 / 3600
       = 0.00272 m/s²
```

### Does the Moon Actually Accelerate at This Rate?

We can check independently. The Moon moves in a roughly circular orbit with:
- Radius: r = 3.84 × 10⁸ m
- Period: T = 27.3 days = 2.36 × 10⁶ s

For circular motion, centripetal acceleration = v²/r = (2πr/T)²/r = 4π²r/T²

```
a_Moon = 4π² × 3.84×10⁸ / (2.36×10⁶)²
       = 4 × 9.87 × 3.84×10⁸ / (5.57×10¹²)
       = 1.517×10¹⁰ / 5.57×10¹²
       = 0.00272 m/s²   ✓
```

The agreement is exact. Gravity does fall off as 1/r². Newton's insight was correct. The same force governs both the falling apple and the orbiting Moon.

### The Conceptual Leap

```
Earth                        Moon
  ●                            ○
  │                            │
  │◄──── 384,000 km ──────────►│
  │
  │ ← 6,370 km
  │
  apple (on surface)
  ◉

Same force. Same law. Different distances.
```

---

## Newton's Law of Universal Gravitation

### The Law

**Newton's Law of Universal Gravitation** states: every mass in the universe attracts every other mass with a force that is:
1. Proportional to the product of their masses
2. Inversely proportional to the square of the distance between them

Written as an equation:

```
         G × m₁ × m₂
F  =  ──────────────────
              r²

Where:
  F  = gravitational force (N)
  G  = gravitational constant = 6.674 × 10⁻¹¹ N·m²/kg²
  m₁ = mass of first object (kg)
  m₂ = mass of second object (kg)
  r  = distance between their centres (m)
```

This is the equation Newton published in his masterwork *Principia Mathematica* in 1687. It was the most powerful scientific equation of its time.

### Key Features

**"Universal"**: Every mass attracts every other mass. You are gravitationally attracted to the person sitting next to you, to the buildings around you, to the stars in other galaxies. Gravity never turns off.

**Product of masses**: Double either mass → double the force. Quadruple both masses → force increases 16×.

**Inverse-square of distance**: Double the distance → force becomes 1/4 as strong. Triple the distance → force becomes 1/9. This is the **inverse-square law**.

**Mutual attraction**: Object 1 pulls Object 2 toward it with force F. Simultaneously, Object 2 pulls Object 1 toward it with the same force F (Newton's Third Law). The forces are equal in magnitude, opposite in direction.

```
        ←────── F ──────          ────── F ──────►
    m₁ ●                                          ● m₂
        ├────────────── r ──────────────────────┤
```

---

## The Gravitational Constant G

### What Is G?

The **gravitational constant** G is a fundamental constant of nature, like the speed of light c. It tells us the "strength" of gravity in the universe.

```
G = 6.674 × 10⁻¹¹ N·m²/kg²
```

The tiny value of G (0.0000000000667...) tells us that gravity is an extremely weak force. You need astronomically large masses (like planets and stars) to produce noticeable gravitational effects.

### How Was G Measured?

Newton could not measure G directly — he did not know the mass of the Earth! He could only determine the ratio GM_Earth.

G was first measured by Henry Cavendish in 1798, over 100 years after Newton's Principia, using a delicate torsion balance experiment.

```
THE CAVENDISH EXPERIMENT:
                  ─────[─●─────────●─]─────
                          wire                 (torsion balance)
                  large balls placed near
                  small balls → tiny force
                  causes wire to twist
                  measure twist → calculate G
```

The experiment required enormous care to shield against air currents and vibrations. Even today, G is one of the least precisely known fundamental constants.

### Units of G

```
[G] = N·m²/kg²

Check: F = G × m₁ × m₂ / r²
       [F] = [G] × [m]² / [r]²
       N = (N·m²/kg²) × kg² / m²
       N = N  ✓
```

---

## The Inverse-Square Law

### What the Inverse-Square Law Means

The force decreases with the **square** of the distance. This has profound consequences.

```
ASCII INVERSE-SQUARE LAW DIAGRAM:

At distance r:         Force = F
──────────────────────────────────────────────────────────

At distance 2r:        Force = F/4
  ████████████████████████████████████████████████████
  Area = 4 times the area at r (force spread over 4× area)

At distance 3r:        Force = F/9
  ████████████████████████████████████████████████████████████████████████████
  Area = 9 times the area at r

Force lines "spread out" over a larger area as distance increases.
The same total "amount" of gravity is spread over a larger sphere.
```

### Visualisation

Imagine gravity radiating outward from a mass like light from a lamp. At distance r, the light hits a sphere of area 4πr². At distance 2r, the same light hits a sphere of area 4π(2r)² = 16πr². The same total light is spread over 4× the area, so the intensity (light per unit area) is 1/4. This is exactly the inverse-square law.

```
         Mass at centre
              ●
             /│\
            / │ \
           /  │  \
  r ──────●   │   ●
           \  │  /
            \ │ /
             \│/
              ●
         Sphere of radius r
         Surface area = 4πr²

  2r ─────────────────────●
                   Sphere of radius 2r
                   Surface area = 4π(2r)² = 16πr²

Same gravity, spread over 4× the area → ¼ the force.
```

### The Inverse-Square Table

```
┌────────────────┬─────────────────────────────────────┐
│   Distance     │   Gravitational Force               │
├────────────────┼─────────────────────────────────────┤
│       r        │   F                                 │
│      2r        │   F/4      (0.25 F)                 │
│      3r        │   F/9      (0.11 F)                 │
│      4r        │   F/16     (0.0625 F)               │
│      5r        │   F/25     (0.04 F)                 │
│     10r        │   F/100    (0.01 F)                 │
│    100r        │   F/10000  (0.0001 F)               │
└────────────────┴─────────────────────────────────────┘

Notice: gravity never reaches zero. It just gets smaller and smaller.
Even at 100× the distance, some gravitational force remains.
```

### Gravity Never Truly Disappears

Unlike a wall blocking sound, gravity extends to infinite distance. It gets weaker but never vanishes. In principle, every atom in your body exerts a gravitational pull on every atom in the Andromeda Galaxy (2.5 million light-years away). The force is absurdly small, but it is not zero.

---

## Comparing Gravitational and Electric Forces

Gravity is not the only force that obeys an inverse-square law. Electric forces (between charged particles) also obey an inverse-square law (**Coulomb's Law**):

```
           k × q₁ × q₂
F_electric = ──────────────────
                  r²

Where k = 8.99 × 10⁹ N·m²/C²  (Coulomb's constant)
```

The two laws have the same mathematical form! But they differ dramatically in strength.

### The Proton-Electron Comparison

Consider a proton and electron in a hydrogen atom, separated by r = 5.3 × 10⁻¹¹ m.

```
Gravitational force:
  m_proton = 1.67 × 10⁻²⁷ kg
  m_electron = 9.11 × 10⁻³¹ kg

  F_grav = G × m_p × m_e / r²
         = 6.674×10⁻¹¹ × 1.67×10⁻²⁷ × 9.11×10⁻³¹ / (5.3×10⁻¹¹)²
         = 6.674×10⁻¹¹ × 1.52×10⁻⁵⁷ / 2.81×10⁻²¹
         = 3.6 × 10⁻⁴⁷ N

Electric force:
  q_proton = +1.6 × 10⁻¹⁹ C
  q_electron = -1.6 × 10⁻¹⁹ C  (opposite signs → attractive)

  F_elec = k × |q₁| × |q₂| / r²
         = 8.99×10⁹ × (1.6×10⁻¹⁹)² / (5.3×10⁻¹¹)²
         = 8.99×10⁹ × 2.56×10⁻³⁸ / 2.81×10⁻²¹
         = 8.2 × 10⁻⁸ N

Ratio:
  F_elec / F_grav = 8.2×10⁻⁸ / 3.6×10⁻⁴⁷ ≈ 2.3 × 10³⁹
```

### What This Means

The electric force between a proton and an electron is **2.3 × 10³⁹ (a trillion trillion trillion)** times stronger than gravity between them.

```
┌──────────────────────────────────────────────────────────────────┐
│  FORCE COMPARISON (between proton and electron in hydrogen)      │
├──────────────────────────────────────────────────────────────────┤
│  Electric Force:     ████████████████████████████████  strong   │
│  Gravitational:      ·  (essentially invisible)         weak    │
│                                                                  │
│  Ratio: Electric is 2,300,000,000,000,000,000,000,000,          │
│         000,000,000,000,000 times stronger.                      │
└──────────────────────────────────────────────────────────────────┘
```

### Then Why Does Gravity Dominate at Large Scales?

Electric charges come in positive and negative. Large objects (like planets) are electrically neutral — equal positive and negative charges nearly cancel. So the electric force between planets is almost zero.

But mass only comes in one "sign" (there is no negative mass). Gravity always adds up. Every kilogram of a planet adds to its gravitational pull. So at planetary and cosmic scales, gravity wins by default.

```
ELECTRIC FORCE: ++ and -- attract, +- cancel out at large scale → ~0
GRAVITY:        All masses attract, nothing cancels, adds up at large scale → dominates
```

---

## Surface Gravity: Where Does g = 9.8 m/s² Come From?

We know that near Earth's surface, objects fall with acceleration g = 9.8 m/s². Now we can derive this from Newton's Law of Gravitation.

### Deriving g

Consider an object of mass m sitting on Earth's surface. Earth has mass M_Earth and radius R_Earth.

```
Gravitational force on the object:
    F = G × M_Earth × m / R_Earth²

This force causes acceleration g (Newton's Second Law):
    F = m × g

Setting them equal:
    m × g = G × M_Earth × m / R_Earth²

The mass m cancels:
    g = G × M_Earth / R_Earth²
```

This is the formula for surface gravity:

```
    g = G × M / R²
```

Where M is the mass of the planet and R is its radius.

### Calculating Earth's g

```
G = 6.674 × 10⁻¹¹ N·m²/kg²
M_Earth = 5.972 × 10²⁴ kg
R_Earth = 6.371 × 10⁶ m

g = (6.674 × 10⁻¹¹ × 5.972 × 10²⁴) / (6.371 × 10⁶)²
  = (3.985 × 10¹⁴) / (4.059 × 10¹³)
  = 9.81 m/s²   ✓
```

Exactly what we measure experimentally. The formula works.

### Why Does a Bowling Ball Fall at the Same Rate as a Feather?

Notice that in g = G × M_Earth / R_Earth², the mass m of the falling object cancelled out. The acceleration due to gravity does not depend on the mass of the falling object. A bowling ball and a feather fall with the same acceleration in a vacuum. (Air resistance slows the feather in air, but that is not gravity's fault.)

---

## Gravity Varies with Altitude

### Above the Surface

As you go higher above Earth's surface, the distance r increases, so gravity weakens.

At height h above the surface:

```
r = R_Earth + h

g_h = G × M_Earth / (R_Earth + h)²

Or equivalently:

g_h = g_surface × R_Earth² / (R_Earth + h)²

g_h = g_surface / (1 + h/R_Earth)²
```

### Gravity at Various Altitudes

```
┌─────────────────────────────────────────────────────────────────┐
│           GRAVITY vs. ALTITUDE ABOVE EARTH'S SURFACE           │
├───────────────────────────────┬─────────────────────────────────┤
│  Altitude (km)                │  g (m/s²)     % of surface g   │
├───────────────────────────────┼─────────────────────────────────┤
│  0 (surface)                  │  9.81          100%             │
│  10 (airplane altitude)       │  9.78           99.7%           │
│  100 (edge of space)          │  9.50           96.8%           │
│  400 (ISS orbit)              │  8.69           88.6%           │
│  1,000                        │  7.33           74.7%           │
│  5,000                        │  3.99           40.7%           │
│  10,000                       │  2.60           26.5%           │
│  36,000 (geostationary orbit) │  0.224           2.3%           │
│  384,000 (Moon distance)      │  0.0027          0.027%         │
└───────────────────────────────┴─────────────────────────────────┘
```

Note: even on the ISS at 400 km altitude, gravity is 88.6% of its surface value. Astronauts are NOT in zero gravity — they are almost fully in Earth's gravitational field. They feel weightless because they are in freefall (we will explain this in Chapter 25).

### Gravity Below the Surface

What happens if you go below Earth's surface (imagine a tunnel)? The gravity weakens because part of the Earth is now above you (pulling you upward). At the very centre of the Earth, gravity is zero — mass pulls equally from all directions.

```
GRAVITY vs DEPTH/HEIGHT RELATIVE TO EARTH'S CENTRE:

                        g
    ↑                   ↑
    │                   │  /\
    │               ───/  \───
    │ surface      /         \
    │ ────────────●           \─────────────►
    │             │           │               r
    │             │ (earth    │ (above surface)
    │             │  inside)  │
    │    linear   │           │  inverse-square
    │    increase │           │  decrease
    ↓             ↓           ↓

  Inside Earth: g increases linearly with r from centre.
  Outside Earth: g decreases as 1/r² with distance from centre.
```

---

## Weight on Different Planets

### What Is Weight?

**Weight** is the gravitational force a planet exerts on you. It depends on the planet's mass and radius, not just its mass alone.

```
W = m × g_planet = m × G × M_planet / R_planet²
```

Your mass never changes (it is the amount of matter in you). Your weight changes depending on where you are.

### The Planet Comparison Table

```
┌──────────────────────────────────────────────────────────────────────┐
│              SURFACE GRAVITY ON SOLAR SYSTEM BODIES                  │
├──────────────────┬───────────┬──────────────┬────────────────────────┤
│  Body            │  g (m/s²) │  Relative g  │  Your weight if 70 kg  │
├──────────────────┼───────────┼──────────────┼────────────────────────┤
│  Sun             │  274      │  27.94×      │  19,180 N (1956 kg)    │
│  Mercury         │  3.70     │  0.38×        │   259 N  (26 kg)       │
│  Venus           │  8.87     │  0.90×        │   621 N  (63 kg)       │
│  Earth           │  9.81     │  1.00×        │   687 N  (70 kg)       │
│  Moon            │  1.62     │  0.17×        │   113 N  (12 kg)       │
│  Mars            │  3.72     │  0.38×        │   260 N  (27 kg)       │
│  Jupiter         │  24.8     │  2.53×        │  1,736 N (177 kg)      │
│  Saturn          │  10.4     │  1.06×        │   729 N  (74 kg)       │
│  Uranus          │  8.87     │  0.90×        │   621 N  (63 kg)       │
│  Neptune         │  11.1     │  1.13×        │   779 N  (79 kg)       │
│  Pluto           │  0.62     │  0.063×       │    43 N  (4.4 kg)      │
└──────────────────┴───────────┴──────────────┴────────────────────────┘
```

### Interesting Observations from the Table

**Moon (0.17×)**: You would weigh only 17% as much on the Moon. This is why Apollo astronauts could leap so high — they still had Earth-trained muscles but only Moon-strength gravity.

**Mars (0.38×)**: Mars has slightly less than half Earth's gravity, despite being about half Earth's size. This is because Mars is much less dense than Earth.

**Jupiter (2.53×)**: Jupiter has 2.53 times Earth's gravity at its cloud tops. You would weigh 253% of your Earth weight. Standing on Jupiter would be physically crushing (and impossible for other reasons, since it has no solid surface).

**Saturn**: Despite being a gas giant 95 times Earth's mass, Saturn's gravity is only slightly stronger than Earth's. This is because Saturn is huge in volume and very low in density (it would float in water, famously).

**Pluto**: At only 6.3% Earth gravity, a person weighing 70 kg on Earth weighs about 4.4 kg on Pluto. You could jump roughly 16 times higher than on Earth.

---

## Escape Velocity

### The Question

If you throw a ball upward, it slows down, stops, and falls back. If you throw it harder, it goes higher before coming back. Is there a speed at which the ball would escape Earth's gravity entirely and never return?

Yes. This is called the **escape velocity**.

### Derivation

The escape velocity is the minimum speed needed to escape a planet's gravity well, given no further propulsion. We find it using energy conservation.

At the surface, the object has:
- Kinetic energy: KE = ½mv²
- Gravitational potential energy: PE = -G × M × m / r  (negative because it is a bound system)

```
To just barely escape, we need v_escape such that the object reaches
infinity (r → ∞) with zero velocity remaining.

At infinity, both KE and PE are zero.

By conservation of energy:
  Total energy at surface = Total energy at infinity

  ½mv² + (-GMm/r) = 0 + 0

  ½mv² = GMm/r

  v² = 2GM/r

  v_escape = √(2GM/r)
```

This is the **escape velocity formula**:

```
    v_escape = √(2GM/r)

Where:
  G = 6.674 × 10⁻¹¹ N·m²/kg²
  M = mass of the planet (kg)
  r = distance from the planet's centre (m)
    (usually taken as the planet's radius for surface escape velocity)
```

### Note: No Mass of Rocket

Just like g, the escape velocity does not depend on the mass of the escaping object. A tennis ball and a rocket need exactly the same escape velocity.

### Earth's Escape Velocity

```
v_escape = √(2 × 6.674×10⁻¹¹ × 5.972×10²⁴ / 6.371×10⁶)
         = √(2 × 3.985×10¹⁴ / 6.371×10⁶)
         = √(7.97×10¹⁴ / 6.371×10⁶)
         = √(1.251×10⁸)
         = 11,186 m/s
         ≈ 11.2 km/s
```

Earth's escape velocity is approximately **11.2 km/s** (about 40,000 km/h, or 25,000 mph).

### Escape Velocities of Solar System Bodies

```
┌──────────────────────────────────────────────────────┐
│        ESCAPE VELOCITIES                             │
├──────────────────┬───────────────────────────────────┤
│  Body            │  Escape Velocity (km/s)           │
├──────────────────┼───────────────────────────────────┤
│  Sun             │  617.5                            │
│  Earth           │  11.2                             │
│  Moon            │  2.38                             │
│  Mars            │  5.03                             │
│  Jupiter         │  59.5                             │
│  Pluto           │  1.23                             │
│  Black hole      │  > 300,000 (faster than light)    │
└──────────────────┴───────────────────────────────────┘
```

The Moon's low escape velocity (2.38 km/s) explains why the Moon has essentially no atmosphere — gas molecules at room temperature move fast enough to escape permanently.

---

## Worked Example 1: Force Between Two People

### Problem

Two people stand 2.0 metres apart. Person A has mass 70 kg, Person B has mass 80 kg. What is the gravitational force between them?

### Given

```
m₁ = 70 kg
m₂ = 80 kg
r = 2.0 m
G = 6.674 × 10⁻¹¹ N·m²/kg²
```

### Solution

```
F = G × m₁ × m₂ / r²

F = 6.674 × 10⁻¹¹ × 70 × 80 / (2.0)²

F = 6.674 × 10⁻¹¹ × 5600 / 4.0

F = 6.674 × 10⁻¹¹ × 1400

F = 9.34 × 10⁻⁸ N
```

### Interpretation

```
F = 0.0000000934 N

This is about 93 nanонewtons (93 nN).

For comparison:
  Weight of a grain of sand ≈ 0.000025 N  (250 times larger!)
  Force to lift a pencil ≈ 0.1 N          (1 million times larger!)
```

The force is utterly negligible in everyday life. You cannot feel it. This is why we never notice gravitational attraction between ordinary objects — G is simply too small. Only planet-sized masses produce noticeable gravity.

---

## Worked Example 2: Gravity at Altitude

### Problem

The International Space Station orbits at an altitude of approximately 408 km above Earth's surface. What is the acceleration due to gravity at this altitude?

### Given

```
h = 408 km = 408,000 m = 4.08 × 10⁵ m
R_Earth = 6.371 × 10⁶ m
g_surface = 9.81 m/s²
```

### Solution — Method 1 (ratio formula)

```
r = R_Earth + h = 6.371 × 10⁶ + 4.08 × 10⁵
  = 6.371 × 10⁶ + 0.408 × 10⁶
  = 6.779 × 10⁶ m

g_h = G × M_Earth / r²

But we know G × M_Earth = g_surface × R_Earth²

So:   g_h = g_surface × R_Earth² / r²

      g_h = 9.81 × (6.371 × 10⁶)² / (6.779 × 10⁶)²

      g_h = 9.81 × (6.371/6.779)²

      g_h = 9.81 × (0.9398)²

      g_h = 9.81 × 0.8832

      g_h = 8.66 m/s²
```

### Result

```
At ISS altitude (408 km):   g = 8.66 m/s²

This is 8.66 / 9.81 = 88.3% of surface gravity.

Gravity on the ISS is NOT zero — it is nearly 90% of what you feel on Earth.
Astronauts feel weightless because they are in freefall, not because gravity is absent.
```

---

## Worked Example 3: Weight on the Moon

### Problem

An astronaut has a mass of 80 kg. Calculate their weight on:
(a) Earth
(b) The Moon

Then calculate the gravitational force the Moon exerts on the astronaut.

### Given

```
m_astronaut = 80 kg
g_Earth = 9.81 m/s²
g_Moon = 1.62 m/s²
G = 6.674 × 10⁻¹¹ N·m²/kg²
M_Moon = 7.342 × 10²² kg
R_Moon = 1.737 × 10⁶ m
```

### Solution: Part (a)

```
W_Earth = m × g_Earth
        = 80 × 9.81
        = 784.8 N
        ≈ 785 N
```

### Solution: Part (b)

```
W_Moon = m × g_Moon
       = 80 × 1.62
       = 129.6 N
       ≈ 130 N
```

### Ratio

```
W_Moon / W_Earth = 130 / 785 = 0.166 ≈ 1/6

The astronaut weighs about one-sixth as much on the Moon.
```

### Verification Using Universal Gravitation

```
F = G × M_Moon × m / R_Moon²

F = 6.674×10⁻¹¹ × 7.342×10²² × 80 / (1.737×10⁶)²

F = 6.674×10⁻¹¹ × 5.874×10²⁴ / 3.017×10¹²

F = 3.919×10¹⁴ / 3.017×10¹²

F = 129.9 N ≈ 130 N  ✓
```

Both methods give the same answer.

---

## Worked Example 4: Escape Velocity from Earth

### Problem

(a) Calculate the escape velocity from Earth's surface.
(b) How does this compare to the speed needed to orbit at Earth's surface (ignore atmosphere)?
(c) If you fire a projectile upward at 8.0 km/s from Earth's surface, does it escape?

### Given

```
G = 6.674 × 10⁻¹¹ N·m²/kg²
M_Earth = 5.972 × 10²⁴ kg
R_Earth = 6.371 × 10⁶ m
```

### Part (a): Escape Velocity

```
v_escape = √(2GM/r)
         = √(2 × 6.674×10⁻¹¹ × 5.972×10²⁴ / 6.371×10⁶)
         = √(2 × 3.985×10¹⁴ / 6.371×10⁶)
         = √(1.251×10⁸)
         = 11,185 m/s
         ≈ 11.2 km/s
```

### Part (b): Orbital Velocity at Surface

For a circular orbit at Earth's surface (we derive this fully in Chapter 25):

```
v_orbital = √(GM/r)
           = √(3.985×10¹⁴ / 6.371×10⁶)
           = √(6.255×10⁷)
           = 7,909 m/s
           ≈ 7.9 km/s

Ratio: v_escape / v_orbital = 11.2 / 7.9 = √2 ≈ 1.414

In general: v_escape = √2 × v_orbital  (for circular orbit at same radius)
```

### Part (c): Projectile at 8.0 km/s

```
v_projectile = 8.0 km/s = 8000 m/s
v_escape = 11,185 m/s

8000 m/s < 11,185 m/s

The projectile does NOT escape Earth's gravity.
It will slow down, stop at some maximum altitude, and fall back.

Maximum altitude calculation (optional):
  Energy conservation:
    ½mv² - GMm/R = 0 - GMm/r_max

    ½v² - GM/R = -GM/r_max

    GM/r_max = GM/R - ½v²

    r_max = GM / (GM/R - ½v²)

  Let's calculate:
    GM = 3.985 × 10¹⁴
    GM/R = 3.985×10¹⁴ / 6.371×10⁶ = 6.255 × 10⁷
    ½v² = ½ × 8000² = 32,000,000 = 3.2 × 10⁷

    GM/r_max = 6.255×10⁷ - 3.2×10⁷ = 3.055 × 10⁷

    r_max = 3.985×10¹⁴ / 3.055×10⁷ = 1.304 × 10⁷ m

  Altitude = r_max - R_Earth
           = 1.304×10⁷ - 6.371×10⁶
           = 6.67 × 10⁶ m
           = 6,670 km above the surface
```

The projectile would reach 6,670 km before falling back — about the altitude of GPS satellites.

---

## Common Mistakes

### Mistake 1: Using Diameter Instead of Radius

The formula uses **r** = distance between the **centres** of the two objects. For objects on Earth's surface, r = R_Earth (Earth's radius), NOT Earth's diameter.

```
WRONG: r = 12,742 km  (diameter)
RIGHT: r = 6,371 km   (radius)
```

### Mistake 2: Forgetting to Square the Distance

F = GMm/r² — the distance is squared. This is the most common arithmetic error.

```
WRONG: F = 6.67×10⁻¹¹ × 1000 × 2000 / 5  =  2.67×10⁻⁵ N
RIGHT: F = 6.67×10⁻¹¹ × 1000 × 2000 / 5² =  5.34×10⁻⁶ N
```

### Mistake 3: Confusing g with G

Little g = 9.8 m/s² is the acceleration due to gravity at Earth's surface. It varies by location and altitude.
Big G = 6.674 × 10⁻¹¹ is the universal gravitational constant. It never changes anywhere in the universe.

### Mistake 4: Thinking Astronauts Are in Zero Gravity

At ISS altitude, gravity is 88% of surface value. Astronauts feel weightless because they are in freefall (free orbit), not because gravity is absent.

---

## Summary

- **Newton's Law of Universal Gravitation**: F = G × m₁ × m₂ / r². Every mass attracts every other mass.

- The **gravitational constant** G = 6.674 × 10⁻¹¹ N·m²/kg². Its tiny value explains why we only notice gravity from planet-sized objects.

- The **inverse-square law** means: double the distance → force drops to one-quarter. The force decreases rapidly with distance but never reaches zero.

- **Surface gravity** is derived as g = G × M / R² (mass of planet M, radius R). For Earth this gives 9.81 m/s², which matches observation.

- Gravity decreases above the surface as 1/(R + h)². At ISS altitude (408 km), g = 8.66 m/s², still 88% of surface value.

- **Weight** = m × g. Your weight changes on different planets; your mass does not.

- Weight on the Moon is 0.17× Earth weight; on Mars 0.38×; on Jupiter 2.53×.

- **Escape velocity** v = √(2GM/r) is the minimum speed to permanently escape a planet's gravity. For Earth, this is 11.2 km/s.

- Gravity is by far the weakest fundamental force — the electric force between a proton and electron is 10³⁹ times stronger. But gravity dominates at cosmic scales because there is no "negative mass" to cancel it.

---

## Key Equations

```
NEWTON'S LAW OF UNIVERSAL GRAVITATION:
    F = G × m₁ × m₂ / r²

GRAVITATIONAL CONSTANT:
    G = 6.674 × 10⁻¹¹ N·m²/kg²

INVERSE-SQUARE SCALING:
    F ∝ 1/r²
    If r doubles: F → F/4
    If r triples: F → F/9

SURFACE GRAVITY (acceleration):
    g = G × M / R²
    (M = planet mass, R = planet radius)

GRAVITY AT ALTITUDE h:
    g_h = G × M / (R + h)²
    g_h = g_surface × R² / (R + h)²

WEIGHT:
    W = m × g

ESCAPE VELOCITY:
    v_escape = √(2GM/r)
    For Earth's surface: v_escape ≈ 11.2 km/s

RELATIONSHIP (escape vs orbital velocity):
    v_escape = √2 × v_orbital   (same radius)
```
