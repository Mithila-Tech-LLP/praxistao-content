# Chapter 31: Springs and Pendulums

> "To every action there is an equal and opposite reaction — and to every displacement, a restoring force." — inspired by Newton

---

## Table of Contents

1. What Is Oscillation?
2. Hooke's Law and the Spring Force
3. The Spring-Mass System
4. Period of a Spring-Mass System
5. What the Period Depends On (and Does NOT Depend On)
6. Worked Example 1 — Finding the Period
7. Worked Example 2 — Finding the Spring Constant
8. The Simple Pendulum
9. Period of a Simple Pendulum
10. Why Mass Does Not Matter for a Pendulum
11. Grandfather Clocks and Precision Timekeeping
12. Pendulum on the Moon vs. Earth
13. Worked Example 3 — Pendulum Period
14. Worked Example 4 — Moon Pendulum
15. Energy in Oscillating Systems
16. Damped Oscillations
17. Driven Oscillations and Resonance Preview
18. Summary
19. Key Equations

---

## 1. What Is Oscillation?

An **oscillation** is a back-and-forth motion that repeats regularly over time. You see oscillations everywhere:

- A child on a swing
- A guitar string vibrating after being plucked
- The second hand on a pendulum clock
- Your car suspension bouncing after a speed bump
- The atoms in every solid object jiggling about fixed positions
- Your heart beating (a mechanical oscillation!)

What all oscillations have in common is a **restoring force**: a force that always pushes or pulls the object back toward a central **equilibrium position**. When the object overshoots (which it always does if there is no friction), the restoring force reverses direction to pull it back the other way. This is why oscillations keep going.

### Vocabulary

- **Equilibrium position**: The resting spot — where the net force is zero.
- **Displacement (x)**: How far the object is from equilibrium at any instant.
- **Amplitude (A)**: The maximum displacement from equilibrium. This is the "size" of the oscillation.
- **Period (T)**: The time to complete one full back-and-forth cycle. Measured in seconds.
- **Frequency (f)**: How many complete cycles happen per second. Measured in hertz (Hz). f = 1/T.
- **Restoring force**: The force directed back toward equilibrium.

```
        Displacement
             ^
             |    /\          /\
             |   /  \        /  \
    ---------+--/----\------/----\----> time
             | /      \    /      \
             |/        \  /        \
                        \/

    A = amplitude (height of peak above center line)
    T = time from one peak to the next peak
```

---

## 2. Hooke's Law and the Spring Force

Robert Hooke (1635–1703) discovered that a spring exerts a force proportional to how much it is stretched or compressed. This is **Hooke's Law**:

```
    F = -kx
```

Where:
- **F** is the restoring force (in Newtons, N)
- **k** is the **spring constant** (in N/m) — a measure of stiffness
- **x** is the displacement from equilibrium (in meters, m)
- The **negative sign** is critical: the force is always opposite to the displacement

The negative sign is everything. If you pull the spring to the right (+x), the spring pulls you back to the left (-F). If you push the spring to the left (-x), the spring pushes you back to the right (+F). The spring always fights your displacement.

### Spring Constant k

The **spring constant k** tells us how stiff the spring is.

- Large k → stiff spring → needs large force for small stretch (e.g., car suspension spring: ~25,000 N/m)
- Small k → soft spring → small force stretches it a lot (e.g., slinky: ~1 N/m)
- A spring with k = 100 N/m needs 100 N to stretch it 1 meter, 50 N to stretch it 0.5 m, etc.

```
    Soft spring (small k):            Stiff spring (large k):

    ~~~~~~~~~~~o                       |||||||o
    
    Easy to stretch,                  Hard to stretch,
    bounces slowly                    bounces fast
```

### Hooke's Law Graphically

```
    Force
    ^
    |        .
    |       .
    |      .  <-- slope = k
    |     .
    |    .
    +----+---------> displacement x
    |   .
    |  .
    | .
    |.
    
    Force is proportional to displacement.
    Steeper slope = stiffer spring = larger k.
```

Hooke's Law holds as long as you do not stretch the spring beyond its **elastic limit**. Past the elastic limit the spring deforms permanently and Hooke's Law no longer applies.

---

## 3. The Spring-Mass System

Attach a mass m to a horizontal spring (spring constant k) on a frictionless surface. Pull the mass to the side by amplitude A and release it.

```
    Equilibrium:
    
    |---spring---[mass]
    |            ^
    |            equilibrium position
    
    Pulled to right (x = +A):
    
    |---spring-------[mass]
                     ^
                     F points LEFT (restoring)
    
    At equilibrium moving right:
    
    |---spring---[mass] -->
                     no net spring force, maximum speed
    
    At left extreme (x = -A):
    
    |--spring--[mass]
               ^
               F points RIGHT (restoring)
```

The mass oscillates back and forth forever on a frictionless surface. In real life, friction gradually reduces the amplitude until the mass stops at equilibrium.

### Why Equilibrium Is Passed

When the mass reaches equilibrium, the spring force is zero — but the mass is moving! Newton's First Law: a moving object keeps moving unless a force stops it. So the mass shoots past equilibrium. Then the spring starts pulling it back, slowing it down. The mass stops, reverses, and the whole cycle repeats.

---

## 4. Period of a Spring-Mass System

Using calculus (specifically solving the differential equation ma = -kx), we find that the period is:

```
    T = 2π × sqrt(m / k)
```

Or equivalently, the frequency is:

```
    f = (1 / 2π) × sqrt(k / m)
```

Where:
- T = period in seconds
- m = mass in kilograms
- k = spring constant in N/m
- sqrt means "square root"
- 2π ≈ 6.283

### Intuitive Reasoning

- **More mass (larger m) → longer period**: A heavy mass accelerates slowly, takes more time to complete a cycle.
- **Stiffer spring (larger k) → shorter period**: A stiffer spring applies stronger restoring force, accelerates the mass faster, completes cycles more quickly.

Both of these make physical sense.

---

## 5. What the Period Depends On (and Does NOT Depend On)

This is one of the most important and counterintuitive facts in oscillation physics:

### The period of a spring-mass system depends ONLY on:
1. **The mass (m)** — more mass = longer period
2. **The spring constant (k)** — stiffer spring = shorter period

### The period does NOT depend on:
1. **Amplitude (A)** — whether you pull the mass 1 cm or 10 cm, the period is the same!

This property is called **isochronism** (from Greek: iso = equal, chronos = time).

### Why Amplitude Doesn't Matter

When you pull the mass farther (larger A), two things happen simultaneously:
- The mass has to travel a longer distance in each cycle.
- But the spring force is also larger (F = -kx), so the mass moves faster.

These two effects exactly cancel out, and the period stays the same.

```
    Small amplitude:          Large amplitude:
    
    |--~--[m]                 |-------[m]
    
    Short distance to cover,  Long distance to cover,
    small force,              large force,
    moves slowly              moves quickly
    
    SAME PERIOD!
```

This result was revolutionary when Galileo first noticed it in pendulums. Before accurate clocks existed, this was nature's gift to timekeepers.

---

## 6. Worked Example 1 — Finding the Period

**Problem**: A mass of 0.5 kg is attached to a spring with spring constant k = 200 N/m. The mass is pulled 8 cm from equilibrium and released. What is the period of oscillation?

**Given**:
- m = 0.5 kg
- k = 200 N/m
- Amplitude A = 8 cm (note: we don't even need this!)

**Formula**:
```
    T = 2π × sqrt(m / k)
```

**Solution**:
```
    T = 2π × sqrt(0.5 / 200)
    T = 2π × sqrt(0.0025)
    T = 2π × 0.05
    T = 6.283 × 0.05
    T = 0.314 seconds
```

**Answer**: T ≈ 0.31 seconds. The mass completes about 3.2 oscillations per second.

**Check**: f = 1/T = 1/0.314 ≈ 3.2 Hz. Does this make sense? A fairly stiff spring (200 N/m) with a small mass (0.5 kg) — yes, it should oscillate fast.

**What if we used A = 16 cm instead?** The period would still be 0.314 seconds! Amplitude does not affect period.

---

## 7. Worked Example 2 — Finding the Spring Constant

**Problem**: A 2 kg block attached to a spring oscillates with a period of 1.2 seconds. What is the spring constant?

**Given**:
- m = 2 kg
- T = 1.2 s

**Formula**:
```
    T = 2π × sqrt(m / k)
```

**Solve for k**:
```
    T = 2π × sqrt(m / k)
    
    T / (2π) = sqrt(m / k)
    
    [T / (2π)]² = m / k
    
    k = m × [2π / T]²
    
    k = 2 × [2π / 1.2]²
    k = 2 × [6.283 / 1.2]²
    k = 2 × [5.236]²
    k = 2 × 27.4
    k = 54.8 N/m
```

**Answer**: k ≈ 54.8 N/m.

**Sanity check**: A 2 kg block oscillating once every 1.2 seconds. This is a moderately soft spring — sensible for a 2 kg mass.

---

## 8. The Simple Pendulum

A **simple pendulum** consists of a mass (called a **bob**) hanging from a string of length L, free to swing back and forth under gravity.

```
                 O  <-- pivot point
                 |
                 |  L (string length)
                 |
                \|/
                [m]  <-- bob
                
    When displaced:
    
                 O
                /|
               / |
              /  |
            [m]  |
            ^    |
            bob  equilibrium
```

When you pull the bob to the side and release it, gravity provides a restoring force. The component of gravity along the arc of the pendulum's path always points back toward the lowest point (equilibrium).

### The Restoring Force in a Pendulum

When the string makes angle θ with the vertical:
```
    Restoring force = -mg × sin(θ)
```

For **small angles** (less than about 15°), sin(θ) ≈ θ (in radians). This **small angle approximation** is what makes the math clean and gives us the simple period formula.

---

## 9. Period of a Simple Pendulum

For small-angle oscillations, the period of a simple pendulum is:

```
    T = 2π × sqrt(L / g)
```

Where:
- T = period in seconds
- L = length of the string in meters
- g = gravitational acceleration = 9.8 m/s²

### What the Period Depends On

The pendulum period depends on:
1. **String length (L)** — longer pendulum = longer period (swings more slowly)
2. **Gravitational acceleration (g)** — stronger gravity = shorter period

### What the Period Does NOT Depend On

1. **Mass of the bob (m)** — a heavy bob and a light bob on the same string swing with the same period
2. **Amplitude (A)** — for small angles, bigger swings still take the same time (the isochronism principle again!)

### Why Mass Doesn't Matter

This connects to Galileo's famous experiment at the Tower of Pisa. When you increase the mass:
- Gravity pulls harder (larger force)
- But the inertia (resistance to acceleration) is also larger
- These cancel exactly, just as in free fall

All masses fall at the same rate, and all masses swing on the same string with the same period.

```
    Light bob (m):            Heavy bob (10m):
    
    O                         O
    |                         |
    |  L                      |  L
    |                         |
   [m]                      [10m]
    
    SAME PERIOD!
    (same as two objects dropped from the same height — same fall time)
```

---

## 10. Why Grandfather Clocks Work

A **grandfather clock** uses a pendulum to regulate its timekeeping. Here's how:

1. The pendulum swings back and forth with a precise, constant period.
2. An **escapement mechanism** (a gear with special teeth) advances by exactly one tick with each swing.
3. Gears translate tick counts into second, minute, and hour hand positions.

For a grandfather clock, the most common design uses a pendulum with period T = 2 seconds (1 second for each swing, left and right).

```
    Required: T = 2 seconds
    
    T = 2π × sqrt(L / g)
    2 = 2π × sqrt(L / 9.8)
    2 / (2π) = sqrt(L / 9.8)
    0.318 = sqrt(L / 9.8)
    0.318² = L / 9.8
    0.101 = L / 9.8
    L = 0.993 m ≈ 1 meter
```

This is why grandfather clocks are about 1 meter from pivot to bob — the physics demands it! The tall case exists specifically to house the ~1 meter pendulum.

### Temperature Effects

Metal expands when heated. In summer, the pendulum gets slightly longer, so the clock runs slightly slow. In winter, it contracts and runs fast. Precision pendulum clocks use **temperature-compensating pendulums** made of materials that expand and contract to keep the effective length constant.

---

## 11. Pendulum on the Moon vs. Earth

The Moon's gravitational acceleration is about g_moon = 1.62 m/s², compared to g_earth = 9.8 m/s².

A pendulum on the Moon swings much more slowly because gravity is weaker. The restoring force is smaller, so the bob accelerates more slowly through each swing.

```
    T_moon / T_earth = sqrt(g_earth / g_moon)
                     = sqrt(9.8 / 1.62)
                     = sqrt(6.05)
                     ≈ 2.46
```

A pendulum that takes 1 second per swing on Earth would take 2.46 seconds per swing on the Moon. A grandfather clock brought to the Moon would run about 2.46 times slower than on Earth — losing about 1 hour 26 minutes every hour!

Astronauts on the Moon would see pendulums swinging in a dreamy, slow-motion way.

---

## 12. Worked Example 3 — Pendulum Period

**Problem**: A simple pendulum has a string length of 0.25 m. What is its period on Earth (g = 9.8 m/s²)?

**Given**:
- L = 0.25 m
- g = 9.8 m/s²

**Formula**:
```
    T = 2π × sqrt(L / g)
```

**Solution**:
```
    T = 2π × sqrt(0.25 / 9.8)
    T = 2π × sqrt(0.02551)
    T = 2π × 0.1597
    T = 6.283 × 0.1597
    T = 1.003 seconds
```

**Answer**: T ≈ 1.00 second.

**Interesting**: A 25 cm pendulum has almost exactly a 1 second period! This is why many wall clocks (without the long pendulum case) use a 25 cm pendulum — one swing (half period = 0.5 s) per tick.

---

## 13. Worked Example 4 — Moon Pendulum

**Problem**: An astronaut wants to build a clock on the Moon using a pendulum. She needs a pendulum with a period of 2.0 seconds. How long should the pendulum string be? (g_moon = 1.62 m/s²)

**Given**:
- T = 2.0 s
- g = 1.62 m/s²

**Solve for L**:
```
    T = 2π × sqrt(L / g)
    T / (2π) = sqrt(L / g)
    [T / (2π)]² = L / g
    L = g × [T / (2π)]²
    
    L = 1.62 × [2.0 / (2π)]²
    L = 1.62 × [2.0 / 6.283]²
    L = 1.62 × [0.3183]²
    L = 1.62 × 0.1013
    L = 0.164 m ≈ 16.4 cm
```

**Answer**: The pendulum only needs to be about 16 cm long to give a 2-second period on the Moon.

**Compare with Earth**: On Earth, the same 2-second pendulum requires about 1 meter (as we calculated for the grandfather clock). The Moon needs only 16 cm because gravity is so much weaker — the pendulum barely needs any length to swing slowly!

---

## 14. Energy in Oscillating Systems

In an ideal (frictionless) spring-mass or pendulum system, energy is constantly swapped between two forms:

- **Kinetic energy (KE)**: energy of motion, KE = (1/2)mv²
- **Potential energy (PE)**: stored energy. In a spring: PE = (1/2)kx². In a pendulum: PE = mgh (height above lowest point)

```
    Spring-mass energy at different positions:

    Position 1: x = A (maximum displacement)
    KE = 0 (stopped)         PE = (1/2)kA² (maximum)

    Position 2: x = 0 (equilibrium)
    KE = max (fastest)       PE = 0 (minimum)

    Position 3: x = -A (maximum displacement other side)
    KE = 0 (stopped)         PE = (1/2)kA² (maximum)

    Total energy E = KE + PE = (1/2)kA² = constant (conserved!)
```

The total mechanical energy equals (1/2)kA², which is why amplitude determines the total energy (even though it doesn't affect the period).

---

## 15. Damped Oscillations

In real systems, friction and air resistance always steal energy. The amplitude gradually decreases over time. This is **damping**.

```
    Displacement
    ^
    |  A
    |  /\       /\       /\
    | /  \   /\/  \/\  /\/  \
    |/    \/            \/    \
    +-------------------------> time
    
    Amplitude shrinks with each cycle (damped oscillation)
```

There are three types of damping, depending on how strong the damping is:

### Underdamped

The system oscillates but with decreasing amplitude. Most real oscillators (car suspension, musical instruments) are slightly underdamped.

```
    |  /\    /\   /\  /\
    | /  \  /  \/  \/  \
    |/    \/             \___
    +-------------------------> time
    Still oscillates, amplitude decays exponentially
```

### Critically Damped

The system returns to equilibrium as fast as possible without oscillating. This is the ideal for many engineering applications: car shock absorbers are designed to be close to critically damped so the car does not bounce repeatedly after a bump.

```
    |  /
    | /
    |/
    |  \_______
    +----------> time
    No oscillation, fastest return to equilibrium
```

### Overdamped

The system returns to equilibrium slowly without oscillating, but more slowly than critically damped. A door with very heavy hydraulic damping might be overdamped.

```
    |  /
    | /
    |/
    |    \________
    +--------------> time
    Slow return, no oscillation
```

---

## 16. Driven Oscillations and Resonance Preview

What happens when you apply a periodic (repeating) force to an oscillating system? This is a **driven oscillation**.

Think about pushing a child on a swing:
- If you push at random times: the child doesn't go very high
- If you push too fast or too slow: the swing loses energy
- If you push exactly once per period, matching the swing's natural rhythm: the child goes higher and higher

Every oscillating system has a **natural frequency** (f₀) — the frequency at which it oscillates freely when disturbed.

When the **driving frequency** equals the **natural frequency**, you get **resonance**. At resonance:
- The driving force always adds energy at the perfect time
- Amplitude builds up dramatically
- In undamped systems, amplitude would grow without limit!

```
    Amplitude
    ^
         |
         |   RESONANCE PEAK
         |       |
         |      /|\
         |     / | \
         |    /  |  \
         |___/   |   \___
    -----+-------+--------> driving frequency
                f₀
```

Resonance will be explored in full in Chapter 36. But here is a preview of why it matters:

- **Musical instruments**: Strings and air columns resonate at specific frequencies to produce musical notes.
- **Radio receivers**: Tuning a radio means adjusting a circuit to resonate at the desired station's frequency.
- **Tacoma Narrows Bridge (1940)**: Wind drove the bridge at its natural frequency, causing catastrophic oscillations that collapsed it.
- **Microwave ovens**: Microwaves drive water molecules at their resonant frequency, heating food.
- **MRI machines**: Radio waves drive hydrogen atoms at resonance to produce medical images.

---

## 17. Comparing Springs and Pendulums

| Property | Spring-Mass | Simple Pendulum |
|----------|-------------|-----------------|
| Restoring force | F = -kx | F = -mg·sin(θ) ≈ -mg·θ |
| Period formula | T = 2π√(m/k) | T = 2π√(L/g) |
| Depends on mass? | YES (more mass = longer T) | NO |
| Depends on amplitude? | NO | NO (small angles) |
| Depends on g? | NO | YES |
| What controls period? | m and k | L and g |

Both systems show **isochronism** — amplitude independence — and both have period proportional to 2π times the square root of (inertia / restoring factor).

---

## 18. Summary

- **Oscillation** is periodic back-and-forth motion driven by a restoring force that always points toward equilibrium.

- **Hooke's Law**: F = -kx. The spring force is proportional to displacement and always directed back toward equilibrium. The spring constant k measures stiffness.

- **Spring-mass period**: T = 2π√(m/k). Depends on mass and spring constant. Does NOT depend on amplitude.

- **Simple pendulum period**: T = 2π√(L/g). Depends on string length and gravitational acceleration. Does NOT depend on mass or amplitude (for small angles).

- **Isochronism**: The remarkable property that period is independent of amplitude. This is why pendulums make good clocks.

- **Grandfather clocks** use ~1 meter pendulums to achieve a 2-second period (1 second per swing).

- **Moon pendulums** swing more slowly because g is smaller. A 1-meter Earth pendulum would take ~2.46 seconds per swing on the Moon.

- **Energy** in oscillating systems is constantly exchanged between kinetic (motion) and potential (stored) energy. Total energy = (1/2)kA².

- **Damping** reduces amplitude over time. Three types: underdamped (oscillates, dies away), critically damped (fastest return, no oscillation), overdamped (slow return, no oscillation).

- **Resonance** occurs when driving frequency = natural frequency, producing large amplitude oscillations. Preview of Chapter 36.

---

## 19. Key Equations

```
Hooke's Law:
    F = -kx
    (F in N, k in N/m, x in m)

Spring-Mass Period:
    T = 2π × sqrt(m / k)
    (T in s, m in kg, k in N/m)

Spring-Mass Frequency:
    f = (1 / 2π) × sqrt(k / m)
    (f in Hz)

Period-Frequency Relationship:
    f = 1 / T        T = 1 / f

Simple Pendulum Period:
    T = 2π × sqrt(L / g)
    (T in s, L in m, g in m/s²)

Potential Energy in Spring:
    PE = (1/2) × k × x²
    (PE in J, k in N/m, x in m)

Kinetic Energy:
    KE = (1/2) × m × v²
    (KE in J, m in kg, v in m/s)

Total Mechanical Energy (spring-mass):
    E = (1/2) × k × A²
    (E in J, A = amplitude in m)

Comparing Pendulum on Different Planets:
    T_1 / T_2 = sqrt(g_2 / g_1)
```
