# Chapter 28: Angular Momentum

> "The universe is under no obligation to make sense to you."
> — Neil deGrasse Tyson
>
> (But angular momentum conservation actually does make sense — and it is one of the most beautiful laws in all of physics.)

---

## Table of Contents

- [Introduction: The Mystery of the Spinning Skater](#introduction-the-mystery-of-the-spinning-skater)
- [What Is Angular Momentum?](#what-is-angular-momentum)
- [Angular Momentum of a Point Mass](#angular-momentum-of-a-point-mass)
- [Angular Momentum of a Rotating Object](#angular-momentum-of-a-rotating-object)
- [The Law of Conservation of Angular Momentum](#the-law-of-conservation-of-angular-momentum)
- [The Figure Skater: A Detailed Calculation](#the-figure-skater-a-detailed-calculation)
- [The Diver Tucking](#the-diver-tucking)
- [The Earth, the Moon, and the Tides](#the-earth-the-moon-and-the-tides)
- [Gyroscopes: Angular Momentum Resists Change](#gyroscopes-angular-momentum-resists-change)
- [Gyroscope Applications](#gyroscope-applications)
- [Precession: The Wobble of a Gyroscope](#precession-the-wobble-of-a-gyroscope)
- [Neutron Stars and Pulsars](#neutron-stars-and-pulsars)
- [Worked Example 1: Spinning Stool](#worked-example-1-spinning-stool)
- [Worked Example 2: Diver Tuck](#worked-example-2-diver-tuck)
- [Worked Example 3: Collapsing Star](#worked-example-3-collapsing-star)
- [Worked Example 4: Planet in Orbit](#worked-example-4-planet-in-orbit)
- [Angular Momentum and Quantum Mechanics](#angular-momentum-and-quantum-mechanics)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: The Mystery of the Spinning Skater

You have almost certainly seen this at the Olympics or in a skating show. A figure skater glides onto the ice, pushes off, and begins to spin slowly with arms spread wide. Then, in an instant, they pull their arms tightly into their body — and the spinning becomes incredibly fast, a blur of motion. When they extend their arms again, the spin slows to its original rate.

```
    FIGURE SKATER: ARMS OUT vs ARMS IN
    
    Arms Out (slow spin):          Arms In (fast spin):
    
         o                               o
        /|\                             \|/
       / | \                             |
      /  |  \                            |
     /   |   \                           |
    
    ω is small                     ω is large
    
    Mass is spread out             Mass is concentrated
    near the axis
    I is large                     I is small
    
    But I × ω stays the same!
```

No one pushed the skater. No external force acted. No engine provided energy. The skater's speed simply changed because they moved their own arms. How is this possible?

The answer is **conservation of angular momentum** — one of the deepest and most universal laws in physics. Angular momentum, like linear momentum and energy, is a conserved quantity. It cannot be created or destroyed — only transferred.

This same law explains why:
- A diver tucks to spin faster
- The Moon is slowly moving away from Earth
- A gyroscope in your phone knows which direction is up
- Neutron stars (tiny collapsed stars) spin hundreds of times per second

Let us understand it fully.

---

## What Is Angular Momentum?

In Chapter 7, we studied **linear momentum**: p = m × v. Linear momentum is the quantity of motion in a straight line. It is conserved when no net external force acts.

**Angular momentum** is the rotational equivalent of linear momentum. It is the quantity of rotational motion. It is conserved when no net external **torque** acts.

The symbol for angular momentum is **L** (capital L).

Angular momentum measures how much an object is "committed" to its current rotation. A fast-spinning, massive, spread-out object has a lot of angular momentum — it will strongly resist any attempt to stop its spin or change its axis.

---

## Angular Momentum of a Point Mass

Consider a single small mass m moving in a circle of radius r with speed v (and therefore angular velocity ω = v/r):

```
    POINT MASS IN CIRCULAR MOTION
    
              v (tangential speed)
              ↑
              |
         m ●--+
              |
              |<--- r --->|
                         *  (axis of rotation)
```

The angular momentum of this point mass is:

```
    L = m × v × r
```

Or, since v = ω × r:

```
    L = m × (ω × r) × r = m × r² × ω = I × ω
```

(where I = mr² is the moment of inertia of a point mass, from Chapter 27)

Units of angular momentum: kg·m²/s (kilogram-meters-squared per second)

**Direction**: Angular momentum is a vector. For circular motion in a plane, we use the sign convention from Chapter 27:
- Counterclockwise rotation → L is positive
- Clockwise rotation → L is negative

(In 3D, the angular momentum vector points along the axis of rotation, determined by the right-hand rule.)

---

## Angular Momentum of a Rotating Object

For an extended rotating object (not just a point mass), the angular momentum generalizes to:

```
    L = I × ω
```

Where:
- L = angular momentum (kg·m²/s)
- I = moment of inertia (kg·m²)
- ω = angular velocity (rad/s)

This is the master equation. Everything we need flows from it.

**Compare the parallel structure:**

```
    LINEAR MOMENTUM          ANGULAR MOMENTUM
    ─────────────────────────────────────────
    p = m × v        ↔       L = I × ω
    
    Newton's 2nd law:         Newton's 2nd for rotation:
    F = Δp/Δt        ↔       τ = ΔL/Δt
    
    Conservation:             Conservation:
    If F_net = 0,    ↔       If τ_net = 0,
    then p = constant         then L = constant
```

The analogy is perfect and beautiful.

---

## The Law of Conservation of Angular Momentum

**The Law of Conservation of Angular Momentum**:

> If no net external torque acts on a system, the total angular momentum of the system remains constant.

In equation form:

```
    If τ_net = 0, then L_initial = L_final
    
    I_initial × ω_initial = I_final × ω_final
```

This is one of the most powerful conservation laws in physics. Notice what it means:

If you somehow change the moment of inertia I (by moving mass closer to or farther from the axis of rotation), then the angular velocity ω **must** change in the opposite way to keep L constant.

- I increases → ω decreases (spins slower)
- I decreases → ω increases (spins faster)

This is exactly what the figure skater does. When they pull their arms in, they are reducing I. So ω increases automatically to conserve L.

**Key point**: The skater's muscles do work (they use energy to pull the arms in against centrifugal tendency). So kinetic energy is NOT conserved here — it actually increases! But angular momentum IS conserved. These are separate conservation laws.

---

## The Figure Skater: A Detailed Calculation

Let us put numbers to the figure skater problem.

**Setup**: Model the skater as a cylindrical body plus two point-mass arms.

```
    SKATER MODEL (simplified):
    
    Arms Out:                      Arms In:
    
    Body: cylinder                 Body: same cylinder
    m_body = 50 kg                 m_body = 50 kg
    R_body = 0.2 m                 R_body = 0.2 m
    
    Each arm: point mass           Each arm: same point mass
    m_arm = 4 kg                   m_arm = 4 kg
    r_arm_out = 0.8 m              r_arm_in = 0.2 m
    (arms extended)                (arms tucked in)
    
         o                              o
        /|\                            \|/
       / | \                            |
      4kg 4kg                          (arms at body radius)
      0.8m 0.8m                        0.2m 0.2m
```

**Step 1: Calculate I_initial (arms out)**

Moment of inertia of the body (solid cylinder):

```
    I_body = ½ × m_body × R_body²
    I_body = ½ × 50 × (0.2)²
    I_body = ½ × 50 × 0.04
    I_body = 1.0 kg·m²
```

Moment of inertia of both arms (point masses):

```
    I_arms_out = 2 × m_arm × r_arm_out²
    I_arms_out = 2 × 4 × (0.8)²
    I_arms_out = 2 × 4 × 0.64
    I_arms_out = 5.12 kg·m²
```

Total initial moment of inertia:

```
    I_initial = I_body + I_arms_out
    I_initial = 1.0 + 5.12
    I_initial = 6.12 kg·m²
```

**Step 2: Calculate I_final (arms in)**

```
    I_arms_in = 2 × m_arm × r_arm_in²
    I_arms_in = 2 × 4 × (0.2)²
    I_arms_in = 2 × 4 × 0.04
    I_arms_in = 0.32 kg·m²
    
    I_final = I_body + I_arms_in
    I_final = 1.0 + 0.32
    I_final = 1.32 kg·m²
```

**Step 3: Apply conservation of angular momentum**

Suppose the initial spin rate is ω_initial = 2 rad/s (slow spin with arms out).

```
    L_initial = L_final
    I_initial × ω_initial = I_final × ω_final
    
    6.12 × 2 = 1.32 × ω_final
    12.24 = 1.32 × ω_final
    ω_final = 12.24 / 1.32
    ω_final = 9.27 rad/s
```

**The skater speeds up from 2 rad/s to 9.27 rad/s — nearly a 5× increase in spin rate!**

This corresponds to going from about 19 RPM to about 89 RPM. That is the blurring speed you see in competitions.

**Step 4: Check the kinetic energy (should increase)**

```
    KE_initial = ½ × I_initial × ω_initial²
    KE_initial = ½ × 6.12 × (2)²
    KE_initial = ½ × 6.12 × 4
    KE_initial = 12.24 J
    
    KE_final = ½ × I_final × ω_final²
    KE_final = ½ × 1.32 × (9.27)²
    KE_final = ½ × 1.32 × 85.93
    KE_final = 56.71 J
```

The kinetic energy increased from 12.24 J to 56.71 J — an increase of 44.5 J. Where did this energy come from? The skater's muscles did work pulling the arms inward against the tendency of mass to fly outward. The skater gets tired doing this!

---

## The Diver Tucking

A springboard diver performs the same physics in the air.

```
    DIVER PHASES:
    
    Phase 1: Straight position (layout)
    
    ┌─────────────────────────────────────┐
    │                                     │
    │   ======o======                    │
    │   (arms and legs extended)          │
    │   I is large, ω is small           │
    │   Somersaults slowly               │
    └─────────────────────────────────────┘
    
    Phase 2: Tuck position
    
    ┌─────────────────────────────────────┐
    │                                     │
    │          (O)                        │
    │         /(|)\                       │
    │   (arms and legs pulled in)         │
    │   I is small, ω is large           │
    │   Somersaults rapidly              │
    └─────────────────────────────────────┘
    
    Phase 3: Straight again (entry)
    
    ┌─────────────────────────────────────┐
    │   ======o======                    │
    │   (extends before hitting water)   │
    │   I large again, ω small again     │
    │   Enters water with minimal splash │
    └─────────────────────────────────────┘
```

After the diver leaves the board, gravity acts through their center of mass and creates no torque about their center (the lever arm is zero). So angular momentum is conserved throughout the dive.

By tucking, the diver rapidly decreases I and therefore increases ω — completing multiple somersaults in the short time before hitting the water. By extending again before entry, they slow the rotation and enter cleanly, nearly vertical.

Experienced divers can precisely control how many somersaults and twists they complete by adjusting their body position in flight.

---

## The Earth, the Moon, and the Tides

One of the most surprising applications of angular momentum conservation is happening right now, gradually, over millions of years.

**The situation**: The Moon's gravity creates tides in Earth's oceans. As Earth rotates, these tidal bulges — the high tides — are dragged slightly ahead of the Moon (because Earth rotates faster than the Moon orbits).

```
    EARTH-MOON TIDAL DIAGRAM (exaggerated)
    
                  MOON
                   O
                  /
                 /  (Moon pulls the bulge backward)
    ====== >>>  /
   /  tidal   \/
   |  bulge    Earth (rotating CCW)
   \  ahead  /
    ======
    
    The tidal bulge is slightly ahead of the Moon
    The Moon's gravity pulls the bulge backward
    This acts as a braking torque on Earth's rotation
```

The Moon's gravity exerts a small **torque** on Earth's tidal bulge, slowing Earth's rotation slightly. This means:

- **Earth's rotation is slowing down**: About 1.4 milliseconds per century. A day is getting longer, very slowly.
- **Earth's angular momentum is decreasing**

But angular momentum of the **Earth-Moon system** must be conserved (the Sun's tidal effects are much smaller and we can initially ignore them). So where does the angular momentum go?

The Moon must gain angular momentum. And the only way the Moon can gain angular momentum in its orbit is to **move to a larger orbit** — farther from Earth.

- **The Moon is moving away from Earth**: About 3.8 centimeters per year.

This has been confirmed precisely by laser ranging experiments. The Apollo astronauts left retroreflectors on the Moon's surface, and scientists bounce lasers off them to measure the Moon-Earth distance to millimeter accuracy.

The same process happened on the Moon long ago: Earth's tidal forces slowed the Moon's rotation until it became **tidally locked** — one face always pointing toward Earth. This is why we always see the same side of the Moon.

---

## Gyroscopes: Angular Momentum Resists Change

A **gyroscope** is a rapidly spinning wheel or disk. When spinning, it has a large angular momentum vector pointing along its spin axis.

The key property of a gyroscope: **a spinning object resists changes to its rotation axis**.

Why? Because changing the direction of the angular momentum vector requires applying a torque. Newton's second law for rotation says τ = ΔL/Δt — to change L (either its magnitude or direction), you need a net torque. If no torque acts, L stays constant — both in magnitude AND direction.

```
    GYROSCOPE SPIN AXIS STABILITY
    
    Non-spinning top:              Spinning gyroscope:
    
         |                              |
         |  (tilts and falls)           | spin
         |                             / axis (stable!)
        / \                            /
       / o \                          / o
      /_____|                        /_____
                                     
    Gravity tips it over         Angular momentum resists tilting
    immediately                  — gyroscope maintains its axis
```

This resistance to orientation change is called **gyroscopic stability**. It is what makes:
- A spinning top stay upright (momentarily)
- A bicycle wheel stable while rolling (try spinning a bicycle wheel and holding it — you can feel it resist your attempts to tilt it)
- A thrown frisbee maintain a flat orientation
- A bullet stable in flight (rifled barrels spin the bullet)

---

## Gyroscope Applications

### Aircraft Navigation (Gyrocompass)

Before GPS, aircraft relied on gyroscopes to determine their heading. A gyroscope whose spin axis is aligned with True North will maintain that alignment even as the aircraft turns, because its angular momentum resists reorientation. The pilot can read the heading relative to the stable gyroscope axis.

```
    AIRCRAFT GYROCOMPASS (simplified):
    
    Before turn:            After aircraft turns right:
    
    Spin axis → North       Spin axis still → North
         |                       |
    [Aircraft]→             [Aircraft]↗
    
    The gyroscope "remembers" North
    even as the aircraft changes direction.
```

### Smartphones and Gaming Controllers

Your smartphone contains tiny **MEMS gyroscopes** (Micro-Electro-Mechanical Systems). These microscopic spinning structures (smaller than a human hair) detect rotation by measuring how the Coriolis effect deflects their vibrating elements.

When you rotate your phone, the gyroscope detects the change in orientation and updates the screen. This is why your phone knows when you are holding it horizontally vs. vertically, and why game controllers can detect tilting.

### Bicycle Stability

Ever wondered why a moving bicycle stays upright but falls over when stationary? Part of the answer (though not all of it) involves gyroscopic effects. The spinning wheels have angular momentum vectors pointing horizontally along the wheel axle. To tip the bicycle over, gravity would have to change the direction of that angular momentum vector — which requires a torque. This provides some stabilizing effect at speed.

```
    BICYCLE WHEEL ANGULAR MOMENTUM
    
    Side view:                Top view:
    
        ○                     L points this way
       / \                    ←─────────────────
      /   \                   (along wheel axle)
       ○○○
    
    The L vector points along the axle.
    Tipping the bike changes the direction of L.
    This requires torque — provided by the steering
    mechanism and rider's adjustments.
```

### Ships and Camera Stabilization

Large ships carry massive gyroscopes to reduce rolling in rough seas. The gyroscope's angular momentum resists the ship's tendency to roll from side to side.

Camera stabilizers (including those in action cameras) use fast-spinning gyroscopes or electronic gyroscopic sensors to detect and compensate for camera shake, producing smooth footage even during motion.

---

## Precession: The Wobble of a Gyroscope

If angular momentum resists change, why does a spinning top eventually wobble and fall? The answer is **precession** — one of the most counterintuitive phenomena in physics.

When a spinning gyroscope has its spin axis tilted, gravity pulls the top downward. But instead of simply tipping over (as a non-spinning object would), the gyroscope's axis sweeps out a **cone** shape. The top "wobbles" in a circle. This is precession.

```
    PRECESSION DIAGRAM
    
    The spin axis traces a cone:
    
         * * * * *
       *           *    <- precession circle
      *             *       (top of spin axis
      *      |      *        traces this path)
       *     |     *
         * * * * *
              |
             [TOP]
              |
           (floor)
    
    Instead of falling straight down (which would change
    L dramatically), the axis slowly circles around,
    which requires only a small, continuous torque.
```

The mathematics of precession:

```
    Precession angular velocity:
    
    Ω_precession = (m × g × r) / (I × ω_spin)
    
    Where:
    Ω = rate of precession (rad/s)
    m = mass of gyroscope
    g = gravitational acceleration
    r = distance from pivot to center of mass
    I = moment of inertia
    ω_spin = spin angular velocity
```

Key insight: **Faster spin = slower precession**. The faster the gyroscope spins (larger L), the more firmly it resists gravity's torque, and the more slowly its axis precesses. A very fast-spinning gyroscope barely wobbles.

**Earth's Precession**: Earth itself is a giant gyroscope. Its spin axis (pointing toward Polaris, the North Star) precesses slowly due to gravitational torques from the Sun and Moon acting on Earth's equatorial bulge. Earth's axis completes one full precession cycle in about **26,000 years**. This means that in 13,000 years, the North Star will be Vega (in the constellation Lyra), not Polaris. Ancient Egyptians and future humans see/will see different pole stars.

---

## Neutron Stars and Pulsars

The most dramatic demonstration of angular momentum conservation in the universe is the birth of a **neutron star**.

When a massive star (much bigger than our Sun) exhausts its nuclear fuel, its core collapses catastrophically in a supernova explosion. In just seconds, the core (which was about the size of Earth) collapses to a ball roughly 10-20 kilometers in diameter — about the size of a city.

```
    STELLAR COLLAPSE DIAGRAM
    
    Before collapse:              After collapse:
    
         * * * *                       *
       *         *                   * * *
      *           *    ──►           * * *   (neutron star)
      *    Sun    *                   * * *
      *   sized   *                     *
       *         *           (10-20 km diameter!)
         * * * *
    
    R ≈ 700,000 km              R ≈ 10 km
    ω ≈ slow                    ω ≈ enormous
    
    I = (2/5)mR²
    
    I changes by factor of (R_final/R_initial)²
    = (10 km / 700,000 km)²
    = (1.4 × 10⁻⁵)²
    ≈ 2 × 10⁻¹⁰
    
    I decreases by factor of ~5 billion!
    So ω must INCREASE by factor of ~5 billion!
```

Let us calculate:

A typical star rotates once per month ≈ 30 days ≈ 2.6 × 10⁶ seconds.

```
    ω_initial = 2π / T_initial = 2π / (2.6 × 10⁶) ≈ 2.4 × 10⁻⁶ rad/s
    
    Using L = Iω = (2/5)mR²ω = constant:
    
    R_initial² × ω_initial = R_final² × ω_final
    
    ω_final = ω_initial × (R_initial / R_final)²
    ω_final = 2.4 × 10⁻⁶ × (7 × 10⁵ km / 10 km)²
    ω_final = 2.4 × 10⁻⁶ × (7 × 10⁴)²
    ω_final = 2.4 × 10⁻⁶ × 4.9 × 10⁹
    ω_final ≈ 11,760 rad/s
    
    T_final = 2π / ω_final ≈ 0.00053 seconds
```

The neutron star would spin at about 1,900 revolutions per second!

Real observations of **pulsars** (neutron stars that emit beams of radio waves) show spin rates ranging from once per second down to hundreds of times per second, in perfect agreement with this prediction. The fastest known pulsar spins at 716 revolutions per second.

This is angular momentum conservation at its most spectacular — a force of nature operating on cosmic scales over the lifetime of a star.

---

## Worked Example 1: Spinning Stool

**Problem**: A physics student sits on a frictionless rotating stool, initially at rest, holding two 3 kg dumbbells. Each dumbbell is held at arm's length, 0.9 m from the rotation axis. The student's moment of inertia (body only) is I_body = 3.5 kg·m². 

A friend gives the student a gentle push, causing them to spin at ω_i = 1.5 rad/s with arms extended. The student then pulls the dumbbells in to 0.1 m from the axis. What is the new spin rate?

**Solution**:

Step 1: Calculate I_initial (dumbbells extended):

```
    I_dumbbells_out = 2 × m × r²
    I_dumbbells_out = 2 × 3 × (0.9)²
    I_dumbbells_out = 2 × 3 × 0.81
    I_dumbbells_out = 4.86 kg·m²
    
    I_initial = I_body + I_dumbbells_out
    I_initial = 3.5 + 4.86
    I_initial = 8.36 kg·m²
```

Step 2: Calculate I_final (dumbbells pulled in):

```
    I_dumbbells_in = 2 × 3 × (0.1)²
    I_dumbbells_in = 2 × 3 × 0.01
    I_dumbbells_in = 0.06 kg·m²
    
    I_final = 3.5 + 0.06
    I_final = 3.56 kg·m²
```

Step 3: Apply conservation of angular momentum:

```
    I_initial × ω_initial = I_final × ω_final
    8.36 × 1.5 = 3.56 × ω_final
    12.54 = 3.56 × ω_final
    ω_final = 12.54 / 3.56
    ω_final = 3.52 rad/s
```

**Answer**: The student spins up to 3.52 rad/s — more than double the original rate.

---

## Worked Example 2: Diver Tuck

**Problem**: A diver leaves the springboard in a "layout" (straight) position with angular velocity ω_i = 1.2 rad/s. In the layout position, I_layout = 15 kg·m². The diver tucks to I_tuck = 4 kg·m². 

(a) What is ω when tucked?
(b) How long does each position last if the total dive takes 1.5 seconds, with 0.5 s in layout, 0.8 s tucked, and 0.2 s in layout again for entry?
(c) How many somersaults are completed in each phase?

**Solution**:

**(a) Angular velocity in tuck:**

```
    I_layout × ω_layout = I_tuck × ω_tuck
    15 × 1.2 = 4 × ω_tuck
    18 = 4 × ω_tuck
    ω_tuck = 4.5 rad/s
```

**(b) Already given in the problem: 0.5 s, 0.8 s, 0.2 s**

**(c) Somersaults (angle rotated divided by 2π):**

Angle = ω × t

```
    Phase 1 (layout, 0.5 s):
    θ_1 = 1.2 rad/s × 0.5 s = 0.6 rad
    Somersaults = 0.6 / (2π) = 0.096 (about 1/10 of a flip)
    
    Phase 2 (tuck, 0.8 s):
    θ_2 = 4.5 rad/s × 0.8 s = 3.6 rad
    Somersaults = 3.6 / (2π) = 0.573 (about half a flip)
    
    Phase 3 (layout, 0.2 s):
    θ_3 = 1.2 rad/s × 0.2 s = 0.24 rad
    Somersaults = 0.24 / (2π) = 0.038 (very small)
    
    Total: about 0.707 somersaults — less than one full flip.
    (Increase tuck time or initial spin rate for a full flip.)
```

---

## Worked Example 3: Collapsing Star

**Problem**: A star with radius 6.0 × 10⁵ km rotates once every 25 days (like our Sun). It collapses to a neutron star of radius 12 km. Assume the mass is conserved and the star can be approximated as a uniform sphere.

(a) Find the initial angular velocity.
(b) Find the final angular velocity using conservation of angular momentum.
(c) Find the rotation period of the neutron star.

**Solution**:

**(a) Initial angular velocity:**

```
    T_initial = 25 days × 24 hrs/day × 3600 s/hr
    T_initial = 25 × 86,400 s
    T_initial = 2,160,000 s = 2.16 × 10⁶ s
    
    ω_initial = 2π / T_initial
    ω_initial = 2π / (2.16 × 10⁶)
    ω_initial = 2.91 × 10⁻⁶ rad/s
```

**(b) Conservation of angular momentum:**

For a solid sphere: I = (2/5)mR²

```
    I_initial × ω_initial = I_final × ω_final
    (2/5)mR_i² × ω_initial = (2/5)mR_f² × ω_final
    
    The (2/5)m cancels:
    R_i² × ω_initial = R_f² × ω_final
    
    ω_final = ω_initial × (R_i / R_f)²
    
    R_i = 6.0 × 10⁵ km = 6.0 × 10⁸ m
    R_f = 12 km = 1.2 × 10⁴ m
    
    (R_i / R_f) = (6.0 × 10⁸) / (1.2 × 10⁴) = 5.0 × 10⁴
    
    (R_i / R_f)² = (5.0 × 10⁴)² = 2.5 × 10⁹
    
    ω_final = 2.91 × 10⁻⁶ × 2.5 × 10⁹
    ω_final = 7,275 rad/s
```

**(c) Rotation period:**

```
    T_final = 2π / ω_final
    T_final = 2π / 7275
    T_final = 8.63 × 10⁻⁴ s
    T_final ≈ 0.86 milliseconds
```

**Answer**: The neutron star rotates about 1,160 times per second — consistent with the fastest known pulsars.

---

## Worked Example 4: Planet in Orbit

**Problem**: A planet moves in an elliptical orbit. At its closest approach to the Sun (**perihelion**), it is at distance r_p = 1.0 × 10¹¹ m and has speed v_p = 3.0 × 10⁴ m/s. At its farthest point (**aphelion**), it is at r_a = 2.0 × 10¹¹ m. What is the planet's speed at aphelion?

(This is Kepler's Second Law — the "equal areas in equal times" law — which is a direct consequence of angular momentum conservation!)

```
    ELLIPTICAL ORBIT DIAGRAM
    
           Aphelion (farthest)
                v_a = ?
              ────────►
    ·        /          \        · 
    ·       /            \       ·
    ·      /      ☀       \      ·
    ·       \  (Sun)      /      ·
    ·        \          /        ·
              ────────►
           Perihelion (closest)
                v_p = 3×10⁴ m/s
```

**Solution**:

For a planet orbiting the Sun, the only force is gravity — directed toward the Sun. This force passes through the Sun (the pivot point), so it creates **zero torque** about the Sun. Therefore, angular momentum is conserved.

At any point in the orbit, the planet's motion is (instantaneously) perpendicular to the radius (at perihelion and aphelion exactly — these are the turning points where v is perpendicular to r).

Using L = m × v × r (for perpendicular velocity):

```
    L_perihelion = L_aphelion
    m × v_p × r_p = m × v_a × r_a
    
    The mass m cancels:
    v_p × r_p = v_a × r_a
    
    v_a = (v_p × r_p) / r_a
    v_a = (3.0 × 10⁴ × 1.0 × 10¹¹) / (2.0 × 10¹¹)
    v_a = (3.0 × 10¹⁵) / (2.0 × 10¹¹)
    v_a = 1.5 × 10⁴ m/s
```

**Answer**: The planet moves at 1.5 × 10⁴ m/s at aphelion — exactly half its perihelion speed, because the distance doubled. This is Kepler's Second Law: a planet sweeps out equal areas in equal times. When far from the Sun, it moves slowly; when close, it moves fast. All from angular momentum conservation.

---

## Angular Momentum and Quantum Mechanics

Angular momentum is not just a classical concept. In quantum mechanics (the physics of atoms and subatomic particles), angular momentum is **quantized** — it can only take certain discrete values, not any arbitrary value.

Electrons in atoms have angular momentum that comes in whole-number multiples of a fundamental unit called **ħ** (h-bar, pronounced "h-bar"):

```
    L_quantum = n × ħ
    
    where ħ = h / (2π) = 1.055 × 10⁻³⁴ kg·m²/s
    (h is Planck's constant)
    
    n = 0, 1, 2, 3, ... (only these values are allowed)
```

This quantization of angular momentum is what creates the discrete energy levels in atoms, which in turn explains why atoms emit and absorb light at specific colors (wavelengths). The entire field of atomic physics — and therefore chemistry — rests on this quantum mechanical version of angular momentum.

Even more surprisingly, elementary particles like electrons have an intrinsic angular momentum called **spin** that has no classical analogue. An electron's spin is always ħ/2 — it can never be zero or changed. This permanent, irreducible angular momentum of fundamental particles is one of the most counterintuitive facts in all of physics.

---

## Summary

- **Angular momentum** (L) is the rotational equivalent of linear momentum: L = I × ω (for rotating objects) or L = m × v × r (for a point mass in circular motion).

- Units of angular momentum: kg·m²/s.

- **Newton's second law for angular momentum**: τ = ΔL/Δt. Torque causes a change in angular momentum.

- **Conservation of Angular Momentum**: If no net external torque acts on a system, the total angular momentum remains constant: L_initial = L_final, so I_i × ω_i = I_f × ω_f.

- A figure skater speeds up by pulling arms in (reducing I → increasing ω to keep L constant). Energy is NOT conserved (the skater does work); angular momentum IS.

- A diver tucks to spin faster and extends to slow before entry — same principle.

- The Moon is slowly moving away from Earth because angular momentum is transferred from Earth's rotation (which is slowing) to the Moon's orbital angular momentum.

- **Gyroscopes** resist changes to their spin axis because changing the direction of L requires a torque. This makes gyroscopes useful for navigation and stabilization.

- **Precession** is the slow sweeping of a gyroscope's axis in a cone, caused by gravity's torque acting on the spinning object. Faster spin = slower precession.

- **Neutron stars** (pulsars) spin hundreds of times per second because a collapsing star with I = (2/5)mR² conserves L: as R decreases by a factor of ~50,000, ω increases by ~2.5 × 10⁹.

- In quantum mechanics, angular momentum is quantized (comes in discrete multiples of ħ), and this quantization underlies atomic structure and chemistry.

---

## Key Equations

```
    ANGULAR MOMENTUM - POINT MASS:
    L = m × v × r
    (v perpendicular to r)
    
    ANGULAR MOMENTUM - ROTATING OBJECT:
    L = I × ω
    
    Units: kg·m²/s
    
    ─────────────────────────────────────────
    
    NEWTON'S SECOND LAW (ANGULAR FORM):
    τ = ΔL / Δt
    
    ─────────────────────────────────────────
    
    CONSERVATION OF ANGULAR MOMENTUM:
    (when τ_net = 0)
    
    L_initial = L_final
    I_i × ω_i = I_f × ω_f
    
    ─────────────────────────────────────────
    
    PRECESSION RATE:
    Ω = (m × g × r) / (I × ω_spin)
    
    ─────────────────────────────────────────
    
    ORBITAL ANGULAR MOMENTUM CONSERVATION:
    v_1 × r_1 = v_2 × r_2
    (perpendicular velocity × radius = constant)
    
    ─────────────────────────────────────────
    
    QUANTUM ANGULAR MOMENTUM:
    L = n × ħ
    ħ = 1.055 × 10⁻³⁴ kg·m²/s
    n = 0, 1, 2, 3, ...
    
    ─────────────────────────────────────────
    
    KEY ANALOGY:
    Linear momentum      Angular momentum
    p = mv        ↔      L = Iω
    F = Δp/Δt     ↔      τ = ΔL/Δt
    p conserved   ↔      L conserved
    when F=0             when τ=0
```
