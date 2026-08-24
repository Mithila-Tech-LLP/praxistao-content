# Chapter 25: Orbits and Satellites

> "The universe is under no obligation to make sense to you."
> — Neil deGrasse Tyson

---

## Table of Contents

- [Introduction: Why Don't Satellites Fall?](#introduction-why-dont-satellites-fall)
- [How Orbits Work: Falling Sideways Fast Enough](#how-orbits-work-falling-sideways-fast-enough)
- [The Circular Orbit Condition](#the-circular-orbit-condition)
- [Orbital Speed Formula](#orbital-speed-formula)
- [Orbital Period T](#orbital-period-t)
- [Kepler's Three Laws of Planetary Motion](#keplers-three-laws-of-planetary-motion)
- [Kepler's First Law: Orbits Are Ellipses](#keplers-first-law-orbits-are-ellipses)
- [Kepler's Second Law: Equal Areas in Equal Times](#keplers-second-law-equal-areas-in-equal-times)
- [Kepler's Third Law: T Squared Proportional to r Cubed](#keplers-third-law-t-squared-proportional-to-r-cubed)
- [Geostationary Satellites](#geostationary-satellites)
- [The International Space Station](#the-international-space-station)
- [Astronaut Weightlessness Explained](#astronaut-weightlessness-explained)
- [Escape Velocity Revisited](#escape-velocity-revisited)
- [Worked Example: ISS Orbital Period](#worked-example-iss-orbital-period)
- [Worked Example: Geostationary Orbit Radius](#worked-example-geostationary-orbit-radius)
- [Worked Example: Kepler's Third Law Applied](#worked-example-keplers-third-law-applied)
- [Common Mistakes](#common-mistakes)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: Why Don't Satellites Fall?

This is one of the most common questions people ask about space. If gravity is pulling the satellite toward Earth, why does it not fall?

The answer: it is falling. It falls constantly. But it also moves sideways fast enough that the Earth curves away beneath it at the same rate it falls. The satellite keeps missing the Earth.

Let that sink in. A satellite is an object that is perpetually falling toward Earth and perpetually missing it because it is moving sideways too fast. This is the essence of orbital mechanics.

In this chapter we will:
- Understand orbits physically and mathematically
- Derive the orbital speed and period formulas
- Explore Kepler's three laws (the rules of all orbiting objects)
- Understand geostationary satellites (GPS, weather, TV)
- Explain why astronauts feel weightless

---

## How Orbits Work: Falling Sideways Fast Enough

### Newton's Cannon: A Thought Experiment

Imagine you stand on a very tall mountain (tall enough to be above the atmosphere, so no air resistance). You fire a cannonball horizontally.

```
ASCII: Newton's Cannon Thought Experiment

                     ●──►  (slow)  falls quickly, lands nearby
                     ●──────►  (faster)  lands farther away
                     ●──────────────►  (even faster)  very far away
                     ●──────────────────────────────────►
                                                          (orbital speed)
                                                          Earth curves
                                                          away at same
                                                          rate as fall!
                     ●──────────────────────────────────────────────────►
                                                          (faster than orbital)
                                                          escapes Earth

         (●) = Earth's surface (curved)
```

At low speeds, the cannonball curves down and hits the ground.
At higher speeds, the ball travels farther before hitting.
At exactly the right speed, the Earth's surface curves away at exactly the same rate the ball falls. The ball never hits — it orbits.
At even higher speeds (escape velocity), the ball escapes entirely.

### The Key Insight: Curvature Matches Fall

Earth's surface curves down by approximately 5 metres for every 8,000 metres you travel horizontally. Gravity accelerates a falling object downward by 5 metres in the first second (s = ½gt² = ½ × 9.8 × 1² = 4.9 m ≈ 5 m).

So if you can move 8,000 metres horizontally in one second, you fall 5 metres and the ground drops 5 metres — you never get closer to the surface. This is orbital speed: 8 km/s.

---

## The Circular Orbit Condition

### Setting Up the Equation

For a circular orbit, the satellite moves at constant speed v in a circle of radius r around Earth. To travel in a circle, the satellite needs a centripetal (centre-seeking) force directed inward.

The only force acting on the satellite is gravity. So:

```
Gravitational force = Centripetal force needed

G × M_Earth × m / r² = m × v² / r

Where:
  m = mass of satellite
  v = orbital speed
  r = radius of orbit (from Earth's CENTRE)
```

### Solving for Orbital Speed

```
G × M × m / r² = m × v² / r

Note: the satellite mass m cancels from both sides!

G × M / r² = v² / r

Multiply both sides by r:

G × M / r = v²

v = √(G × M / r)
```

This is the **orbital speed formula**:

```
    v_orbital = √(G × M / r)

Where:
  G = 6.674 × 10⁻¹¹ N·m²/kg²
  M = mass of the central body (Earth, Moon, Sun, etc.)
  r = orbital radius (distance from centre of central body)
```

### Important Conclusions from v = √(GM/r)

1. **Orbital speed does not depend on the satellite's mass.** A 1 kg CubeSat and a 400,000 kg ISS at the same altitude have the same orbital speed.

2. **Higher orbits have SLOWER orbital speeds.** As r increases, v decreases. Satellites farther out orbit more slowly. This seems counterintuitive but is correct.

3. **For a given orbit, there is exactly one orbital speed.** If the satellite moves faster, it spirals outward. If it moves slower, it spirals inward and eventually falls.

```
ORBITAL SPEED vs. ALTITUDE:

Higher orbit  ─── slower orbital speed
               ●──────────────────►
               │ orbit radius = r₂
               │
               ●──────────────────────────────►
               │ orbit radius = r₁ < r₂
               │
     [EARTH]   │
               │ r₁ < r₂, but v₁ > v₂

Lower orbit ─── faster orbital speed
```

---

## Orbital Speed Formula

### Orbital Speed at Various Altitudes

Using v = √(GM_Earth / r), with G × M_Earth = 3.986 × 10¹⁴ m³/s²:

```
┌──────────────────────────────────────────────────────────────────┐
│              ORBITAL SPEEDS AT VARIOUS ALTITUDES                 │
├─────────────────────────────┬──────────────┬─────────────────────┤
│  Orbit / Altitude           │  r (km)      │  v_orbital (km/s)   │
├─────────────────────────────┼──────────────┼─────────────────────┤
│  Surface (no atmosphere)    │  6,371       │  7.91               │
│  Low Earth Orbit (LEO)      │  6,771       │  7.67               │
│  ISS (408 km)               │  6,779       │  7.66               │
│  GPS satellites             │  26,560      │  3.87               │
│  Geostationary (35,786 km)  │  42,164      │  3.07               │
│  Moon (average)             │  384,400     │  1.02               │
└─────────────────────────────┴──────────────┴─────────────────────┘
```

---

## Orbital Period T

### Period from Speed

The orbital period T is the time to complete one full orbit. For a circular orbit of radius r:

```
Circumference = 2πr
Speed = v

Time = Distance / Speed:
    T = 2πr / v

Substituting v = √(GM/r):
    T = 2πr / √(GM/r)
    T = 2πr × √(r/GM)
    T = 2π × r × r^(1/2) / √(GM)
    T = 2π × r^(3/2) / √(GM)
```

The **orbital period formula**:

```
    T = 2π × √(r³ / (G × M))

Or equivalently:
    T = 2π × r^(3/2) / √(GM)
```

### What This Tells Us

The period depends on r^(3/2) — the orbital radius to the power of 3/2 (or equivalently, r cubed then square-rooted). This is the mathematical basis of Kepler's Third Law.

Higher orbits → larger r → longer period. Geostationary satellites orbit once per day (24 hours). The Moon takes 27.3 days.

---

## Kepler's Three Laws of Planetary Motion

### Historical Context

Johannes Kepler (1571–1630) was a German mathematician and astronomer who worked before Newton's time. He did not know about gravity. He analysed decades of careful astronomical observations by Tycho Brahe and discovered three empirical rules that all planets follow.

When Newton invented calculus and derived his law of gravitation, he could explain all three of Kepler's laws from first principles. This was one of the greatest triumphs of classical physics.

---

## Kepler's First Law: Orbits Are Ellipses

### The Law

**All planets (and all orbiting bodies) move in ellipses, with the Sun (or central body) at one focus of the ellipse.**

Circles are a special case of ellipses (where both foci coincide at the centre).

### What Is an Ellipse?

An ellipse is an oval shape defined by two special points called **foci** (singular: **focus**). For every point on the ellipse, the sum of the distances to the two foci is constant.

```
ASCII ELLIPSE:

         ......
      .         .
    .     F₁  F₂  .
   .       ●   ●   .
    .               .
      .           .
         ......

F₁ = one focus (Sun or Earth)
F₂ = other focus (empty in space)

For a planet: r₁ + r₂ = constant at every point on the orbit
```

### Perihelion and Aphelion

```
ASCII: Elliptical Orbit with Sun at Focus

                ● ← perihelion (closest approach)
               / \
              /   \
             /     \
        ──────● ─────────────────● ──────
         [SUN]                    ← aphelion (farthest point)
             \     /
              \   /
               \ /
                ●

Perihelion: closest point to Sun (planet moves fastest)
Aphelion:   farthest point from Sun (planet moves slowest)
```

### Earth's Elliptical Orbit

Earth's orbit is very nearly circular but slightly elliptical:
- **Perihelion** (closest): ~147 million km (early January)
- **Aphelion** (farthest): ~152 million km (early July)
- Difference: about 3.3%

Interestingly, Earth is closest to the Sun in January (winter in the northern hemisphere). The seasons are caused by the tilt of Earth's axis, NOT by distance from the Sun.

---

## Kepler's Second Law: Equal Areas in Equal Times

### The Law

**A line connecting a planet to the Sun sweeps out equal areas in equal time intervals.**

### What This Means

```
ASCII: Kepler's Second Law

      B ─── A
     /   area₁ \
    / (same     \
   /   area as   \
  /    area₂)     \
──────── [SUN]
  \               /
   \   area₂    /
    \           / C
     \         /
      D ───────

Area swept from A to B (in time Δt) = Area swept from C to D (same Δt)

Since A-B is close to the Sun and C-D is far from the Sun:
  • Near the Sun: smaller arc A-B → but broader wedge → same area
  • Far from the Sun: larger arc C-D → but narrower wedge → same area

Planet moves FASTER near the Sun, SLOWER when far away.
```

### Physical Explanation (from Newton)

Kepler's Second Law is equivalent to conservation of angular momentum. In the absence of external torques, angular momentum L = m × v × r is conserved. When r decreases (closer to Sun), v must increase to keep L = mvr constant.

### Example: Earth's Speed Changes with Seasons

```
Earth's orbital speed:
  January (perihelion, closest): ~30.3 km/s  ← fastest
  July (aphelion, farthest):     ~29.3 km/s  ← slowest
  Difference: about 3.4%
```

The northern hemisphere winter (Earth at perihelion) is slightly shorter than northern hemisphere summer — about 3 days shorter, due to Earth moving faster. This is a tiny but measurable consequence of Kepler's Second Law.

---

## Kepler's Third Law: T Squared Proportional to r Cubed

### The Law

**The square of the orbital period of a planet is proportional to the cube of its average orbital radius.**

```
    T² ∝ r³

Or:   T² / r³ = constant  (same for all objects orbiting the same body)

More precisely:
    T² = (4π² / GM) × r³
```

### The Constant

The constant 4π²/(GM) depends only on the central body. All planets orbiting the Sun have the same ratio T²/r³ (using the same units). All moons of Earth have the same ratio T²/r³.

### Kepler's Third Law Table: Solar System Planets

```
┌──────────────────────────────────────────────────────────────────────────┐
│                KEPLER'S THIRD LAW VERIFICATION                           │
├──────────────────┬─────────────────┬─────────────────┬───────────────────┤
│  Planet          │  r (AU)         │  T (years)      │  T²/r³           │
├──────────────────┼─────────────────┼─────────────────┼───────────────────┤
│  Mercury         │  0.387          │  0.241          │  1.002            │
│  Venus           │  0.723          │  0.615          │  1.001            │
│  Earth           │  1.000          │  1.000          │  1.000            │
│  Mars            │  1.524          │  1.881          │  1.000            │
│  Jupiter         │  5.203          │  11.86          │  0.999            │
│  Saturn          │  9.539          │  29.46          │  1.000            │
│  Uranus          │  19.18          │  84.01          │  0.998            │
│  Neptune         │  30.06          │  164.8          │  0.999            │
└──────────────────┴─────────────────┴─────────────────┴───────────────────┘

1 AU = 1 Astronomical Unit = Earth-Sun distance = 1.496 × 10¹¹ m

Note: T²/r³ is remarkably constant for all planets!
```

### Using Kepler's Third Law to Find Orbital Radii

If you measure a planet's period T, you can immediately find its orbital radius r:

```
r³ = GM × T² / (4π²)

r = (GM × T² / 4π²)^(1/3)
```

This is how early astronomers mapped the solar system — by timing planetary orbits.

---

## Geostationary Satellites

### What Is a Geostationary Orbit?

A **geostationary satellite** orbits Earth at exactly the right altitude so that its orbital period equals Earth's rotation period (24 hours). As the Earth rotates, the satellite keeps pace — it appears to hover motionless above a fixed point on the equator.

```
ASCII: Geostationary Orbit

                    [GPS satellite]
                         ↑ appears fixed
                         │  in the sky
                         │
              ···········●···········
           ·                         ·
         ·                             ·
        ·         [EARTH]              ·
        ·           ●                  ·
         ·         / \                ·
           ·      /   \             ·
              ·  / you  \        ·
               ·(looking up)   ·
                  ···········

From the ground, the satellite appears stationary.
```

### Where Must the Geostationary Orbit Be?

We need T = 24 hours and we need to find r.

```
T = 2π × √(r³ / GM)

T² = 4π² × r³ / GM

r³ = GM × T² / (4π²)

r = (GM × T² / 4π²)^(1/3)

With T = 24 hours = 86,400 s and GM = 3.986 × 10¹⁴ m³/s²:

r³ = 3.986×10¹⁴ × (86,400)² / (4π²)
   = 3.986×10¹⁴ × 7.465×10⁹ / 39.48
   = 2.974×10²⁴ / 39.48
   = 7.534 × 10²² m³

r = (7.534 × 10²²)^(1/3)
  = 4.216 × 10⁷ m
  = 42,164 km from Earth's centre

Altitude above surface = 42,164 - 6,371 = 35,793 km ≈ 35,800 km
```

### Why 35,800 km?

The geostationary orbit is at approximately **35,800 km altitude** (often quoted as 36,000 km). This is about 6.6 Earth radii above the surface — not a specific choice, but a physical consequence of Earth's mass, radius, and rotation rate.

If Earth rotated faster, the geostationary orbit would be lower. If Earth rotated slower, it would be higher.

### Geostationary vs. Geosynchronous

- **Geosynchronous**: period = 24 hours, but orbit may be inclined to the equator → satellite traces a figure-8 path over the ground
- **Geostationary**: geosynchronous AND equatorial AND circular → satellite truly hovers over one point

### Applications of Geostationary Satellites

```
┌─────────────────────────────────────────────────────────────────────┐
│                  GEOSTATIONARY SATELLITE USES                       │
├───────────────────────────────┬─────────────────────────────────────┤
│  Application                  │  Why geostationary is ideal         │
├───────────────────────────────┼─────────────────────────────────────┤
│  Television broadcasting      │  Fixed dish antenna works           │
│  Weather monitoring           │  Same area of Earth always watched  │
│  Communication relay          │  Fixed pointing → simple hardware   │
│  Internet (some services)     │  Ground antenna doesn't need to     │
│                               │  track the satellite                │
├───────────────────────────────┼─────────────────────────────────────┤
│  DISADVANTAGE: Latency        │  Signal travels 35,800 km to sat    │
│  (delay in communication)     │  + 35,800 km back = 71,600 km      │
│                               │  at speed of light: 0.24 seconds   │
│                               │  round-trip delay is noticeable     │
└───────────────────────────────┴─────────────────────────────────────┘
```

### GPS Satellites Are NOT Geostationary

GPS (Global Positioning System) satellites orbit at about 20,200 km altitude — about half the geostationary distance. There are 24+ GPS satellites in medium Earth orbit, arranged so that at least 4 are visible from any point on Earth at any time.

---

## The International Space Station

### ISS Orbital Parameters

The International Space Station orbits at approximately:

```
Altitude:          408 km above Earth's surface
Orbital radius:    r = 6,371 + 408 = 6,779 km = 6.779 × 10⁶ m
Orbital speed:     ~7.66 km/s
Orbital period:    ~92 minutes (15-16 orbits per day)
```

### ISS Speed Calculation

```
v_ISS = √(GM / r)
      = √(3.986×10¹⁴ / 6.779×10⁶)
      = √(5.879×10⁷)
      = 7,667 m/s
      ≈ 7.67 km/s
      ≈ 27,600 km/h
      ≈ 17,100 mph
```

At this speed, the ISS travels the circumference of the Earth in about 92 minutes.

### The ISS Passes Over You Quickly

```
At 7.67 km/s, the ISS crosses the width of the continental USA
(4,500 km) in:

    time = 4,500 / 7.67 = 587 seconds ≈ 10 minutes

A typical ISS flyover lasts about 3-6 minutes as it crosses your sky.
```

---

## Astronaut Weightlessness Explained

### The Common Misconception

Many people think astronauts on the ISS float because they are "in zero gravity" or "too far from Earth's gravity." This is wrong.

As we showed in Chapter 24, at the ISS altitude of 408 km, gravity is still 88% of its surface value (g = 8.69 m/s²). The astronauts are definitely in Earth's gravitational field.

### The Real Reason: Freefall

The ISS and everything inside it (including astronauts) are in **freefall** — they are all falling toward Earth at exactly the same rate. There is no normal force (no floor pushing up on the astronauts), so they feel weightless.

Imagine being in an elevator and the cable snaps. For a brief terrifying moment, you and the elevator fall together. You would feel weightless — you float off the floor, your coffee floats out of the cup. Not because gravity vanished, but because the floor is no longer pushing up on you.

The ISS is doing this continuously — it falls toward Earth, but its sideways speed keeps it in orbit.

```
ELEVATOR ANALOGY:

Normal:                    Freefall:
   ─────────────            ─────────────
   │   (you)   │  ↑N        │   (you)   │  → falling
   │     ●     │            │     ●     │     together
   │           │  ↓W        │           │
   ─────────────            ─────────────
       ↓mg                      ↓mg
   Cable holds elevator up   Cable cut! Both fall together.
   Floor pushes on you (N)   No normal force → weightless feeling
   You feel your weight       You feel weightless (but gravity unchanged)
```

### The Term "Weightlessness"

Technically, astronauts are in a state of **microgravity** (not zero gravity). There are tiny variations in the gravitational field across the ISS that create very small relative forces, but for practical purposes, the environment mimics zero gravity.

### Why This Matters for Space Travel

The prolonged weightless environment causes real physiological effects:
- Muscles atrophy (no gravity to fight)
- Bones lose density
- Fluids shift toward the head (puffy face, thin legs)
- Cardiovascular deconditioning
- Vision problems (increased pressure on eyeballs)

ISS astronauts exercise 2 hours daily to slow these effects.

---

## Escape Velocity Revisited

In Chapter 24 we derived escape velocity v_escape = √(2GM/r). Now we can compare it to orbital velocity at the same radius:

```
v_orbital = √(GM/r)
v_escape  = √(2GM/r) = √2 × √(GM/r) = √2 × v_orbital

v_escape = √2 × v_orbital ≈ 1.414 × v_orbital
```

To escape from any orbit, you need to speed up by a factor of √2 (about 41.4%). This is why rockets leaving Earth orbit need only a relatively modest additional burn — they are already doing 7.67 km/s (orbital speed) and need to reach 11.2 km/s (escape speed).

```
The Three Speed Regimes:

v < v_orbital:   Object spirals inward, eventually falls to Earth
v = v_orbital:   Circular orbit maintained
v_orbital < v < v_escape:  Elliptical orbit (elongated oval)
v = v_escape:    Object just barely escapes (parabolic trajectory)
v > v_escape:    Object escapes with excess speed (hyperbolic trajectory)
```

---

## Worked Example: ISS Orbital Period

### Problem

Calculate the orbital period of the International Space Station, which orbits at an altitude of 408 km above Earth's surface. Express your answer in minutes.

### Given

```
Altitude h = 408 km = 4.08 × 10⁵ m
R_Earth = 6.371 × 10⁶ m
G = 6.674 × 10⁻¹¹ N·m²/kg²
M_Earth = 5.972 × 10²⁴ kg
```

### Step 1: Find Orbital Radius

```
r = R_Earth + h
  = 6.371 × 10⁶ + 4.08 × 10⁵
  = 6.371 × 10⁶ + 0.408 × 10⁶
  = 6.779 × 10⁶ m
```

### Step 2: Calculate GM

```
GM = 6.674 × 10⁻¹¹ × 5.972 × 10²⁴
   = 3.986 × 10¹⁴ m³/s²
```

### Step 3: Apply Period Formula

```
T = 2π × √(r³ / GM)

First calculate r³:
    r³ = (6.779 × 10⁶)³
       = 6.779³ × 10¹⁸
       = 311.3 × 10¹⁸
       = 3.113 × 10²⁰ m³

Now r³ / GM:
    r³ / GM = 3.113 × 10²⁰ / 3.986 × 10¹⁴
            = 7.81 × 10⁵ s²

Take square root:
    √(7.81 × 10⁵) = 883.6 s

Multiply by 2π:
    T = 2π × 883.6 = 5551 s
```

### Step 4: Convert to Minutes

```
T = 5551 s ÷ 60 = 92.5 minutes
```

### Result

```
Orbital Period of ISS ≈ 92.5 minutes  (or about 1 hour 32 minutes)

This means the ISS completes approximately 15-16 full orbits per day.
Astronauts on the ISS see about 15-16 sunrises and sunsets every day!
```

### ASCII Orbit Diagram with Period

```
                     ↑ North
            ──────────────────────────
           /          Earth           \
          /         ●  ●  ●            \
         /      ●           ●           \
        /   ●                   ●        \
       │  ●                       ●       │
       │●                           ●     │
       │ ●    [EARTH]               ●    │
       │  ●      ●                 ●     │
        \    ●               ●          /
         \      ●         ●            /
          \        ● ● ●              /
           \                          /
            ──────────────────────────

One full orbit = one complete circle = 92.5 minutes for ISS

Circumference = 2πr = 2π × 6.779×10⁶ = 4.259 × 10⁷ m = 42,590 km
Speed = 42,590 km / 92.5 min = 460 km/min = 7.67 km/s  ✓
```

---

## Worked Example: Geostationary Orbit Radius

### Problem

Calculate the altitude of a geostationary satellite above Earth's equator. (Verify our earlier calculation.)

### Given

```
T = 24 hours = 24 × 3600 = 86,400 s
GM_Earth = 3.986 × 10¹⁴ m³/s²
R_Earth = 6.371 × 10⁶ m
```

### Solution

```
From Kepler's Third Law:
    T² = 4π² × r³ / GM

    r³ = GM × T² / (4π²)
       = 3.986×10¹⁴ × (86,400)² / (4π²)

Calculate (86,400)²:
    (86,400)² = 7.465 × 10⁹ s²

Calculate numerator:
    3.986×10¹⁴ × 7.465×10⁹ = 2.975 × 10²⁴

Calculate denominator:
    4π² = 4 × 9.870 = 39.48

    r³ = 2.975×10²⁴ / 39.48 = 7.537 × 10²² m³

    r = (7.537 × 10²²)^(1/3)

Calculate cube root:
    7.537 × 10²² = 75.37 × 10²¹
    (75.37)^(1/3) ≈ 4.222
    (10²¹)^(1/3) = 10⁷

    r ≈ 4.222 × 10⁷ m = 42,220 km

Altitude:
    h = r - R_Earth = 42,220 - 6,371 = 35,849 km ≈ 35,800 km
```

### Interpretation

```
Geostationary orbit altitude ≈ 35,800 km
(often rounded to 36,000 km in popular sources)

This is:
  • 5.62 times Earth's radius above the surface
  • Orbital speed: v = 2πr/T = 2π×42,220 km / 86,400 s = 3.07 km/s
  • Much slower than ISS (7.67 km/s) because it is much farther out
```

---

## Worked Example: Kepler's Third Law Applied

### Problem

Mars orbits the Sun with an orbital period of 687 days. Earth's orbital radius is 1.496 × 10¹¹ m. Use Kepler's Third Law to find Mars's average orbital radius.

### Given

```
T_Earth = 365.25 days
r_Earth = 1.496 × 10¹¹ m
T_Mars = 687 days
r_Mars = ?
```

### Solution Using the Ratio Form

For any two planets orbiting the same star:

```
T_A² / r_A³ = T_B² / r_B³

Rearranging:
    r_B³ / r_A³ = T_B² / T_A²

    (r_B / r_A)³ = (T_B / T_A)²

    r_B / r_A = (T_B / T_A)^(2/3)

    r_Mars / r_Earth = (T_Mars / T_Earth)^(2/3)

    r_Mars / r_Earth = (687 / 365.25)^(2/3)

    r_Mars / r_Earth = (1.881)^(2/3)

Calculate (1.881)^(2/3):
  First: (1.881)^(1/3) = cube root of 1.881 ≈ 1.234
  Then:  (1.234)² ≈ 1.524

    r_Mars = 1.524 × r_Earth
           = 1.524 × 1.496×10¹¹
           = 2.280 × 10¹¹ m
           = 228 million km
```

### Verification

The actual average orbital radius of Mars is 1.524 AU = 2.280 × 10¹¹ m. Our calculation matches precisely.

```
r_Mars = 1.524 AU  (AU = astronomical unit = Earth-Sun distance)
T_Mars = 1.881 years

(T_Mars / T_Earth)² = 1.881² = 3.538
(r_Mars / r_Earth)³ = 1.524³ = 3.541  ≈  3.538  ✓
```

---

## Common Mistakes

### Mistake 1: Confusing Orbital Radius with Altitude

In all orbit formulas, r is measured from the **centre of the Earth**, not from the surface.

```
WRONG: For ISS at 408 km altitude, use r = 408 km
RIGHT: r = R_Earth + altitude = 6,371 + 408 = 6,779 km
```

### Mistake 2: Thinking Higher Orbits Are Faster

Higher orbits have LOWER orbital speed. v = √(GM/r) decreases as r increases.

```
Low orbit (400 km):   v ≈ 7.67 km/s  (fast)
Geostationary:         v ≈ 3.07 km/s  (slow)
```

### Mistake 3: Weightlessness = No Gravity

At ISS altitude, g = 8.69 m/s², not zero. Weightlessness is due to freefall (no contact forces), not the absence of gravity.

### Mistake 4: T in Wrong Units

When using T = 2π√(r³/GM), T comes out in seconds (SI units). Convert carefully.

```
1 day = 86,400 s
1 year = 3.156 × 10⁷ s
```

### Mistake 5: Using Wrong Period for Geostationary Orbit

Use 24 hours (sidereal day is actually 23 hours 56 minutes, but 24 hours gives the correct answer for the practical geostationary orbit since we typically use the solar day).

---

## Summary

- **Orbits are perpetual freefall** — a satellite moves sideways fast enough that the Earth curves away at the same rate it falls.

- **Orbital speed** v = √(GM/r). Higher orbits are slower. The satellite's mass cancels out — all satellites at the same altitude have the same orbital speed.

- **Orbital period** T = 2π√(r³/GM). Higher orbits take longer to complete.

- **Kepler's First Law**: Orbits are ellipses with the central body at one focus.

- **Kepler's Second Law**: A line from the planet to the Sun sweeps equal areas in equal times. This is conservation of angular momentum — planets speed up near the Sun, slow down when far away.

- **Kepler's Third Law**: T² ∝ r³. The ratio T²/r³ is the same for all bodies orbiting the same central mass.

- **Geostationary satellites** orbit at 35,800 km altitude with a 24-hour period, appearing stationary above a fixed equatorial point. Used for TV, weather, and communications.

- **ISS** orbits at 408 km, takes 92.5 minutes per orbit, moves at 7.67 km/s, and experiences about 15-16 sunrises per day.

- **Astronaut weightlessness** is NOT caused by zero gravity (g is 88% of normal). It is caused by freefall — there is no surface pushing up, so no normal force, so the astronaut feels weightless.

- **Escape velocity** = √2 × orbital velocity at the same distance.

---

## Key Equations

```
ORBITAL SPEED (circular orbit):
    v = √(G × M / r)
    (r = distance from centre of planet)

ORBITAL PERIOD:
    T = 2π × √(r³ / GM)
    T = 2π × r / v

KEPLER'S THIRD LAW:
    T² = (4π² / GM) × r³
    T²/r³ = constant  (same for all objects orbiting same body)

RATIO FORM (comparing two orbits around same body):
    (T_A / T_B)² = (r_A / r_B)³

RELATIONSHIP: ESCAPE vs ORBITAL:
    v_escape = √2 × v_orbital  (at same radius)

GEOSTATIONARY ORBIT:
    Set T = 24 hours = 86,400 s
    r = (GM × T² / 4π²)^(1/3)
    r ≈ 42,164 km from Earth's centre  ≈ 35,800 km altitude

USEFUL CONSTANTS:
    G × M_Earth = 3.986 × 10¹⁴ m³/s²
    G × M_Sun   = 1.327 × 10²⁰ m³/s²
    R_Earth = 6.371 × 10⁶ m
```
