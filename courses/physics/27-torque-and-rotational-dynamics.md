# Chapter 27: Torque and Rotational Dynamics

> "Give me a lever long enough and a fulcrum on which to place it, and I shall move the world."
> — Archimedes

---

## Table of Contents

- [Introduction: Why Things Spin](#introduction-why-things-spin)
- [What Is Torque?](#what-is-torque)
- [The Lever Arm](#the-lever-arm)
- [The Torque Formula](#the-torque-formula)
- [Direction of Torque](#direction-of-torque)
- [Why Door Handles Are Far From Hinges](#why-door-handles-are-far-from-hinges)
- [Net Torque](#net-torque)
- [Newton's Second Law for Rotation](#newtons-second-law-for-rotation)
- [Moment of Inertia: Rotational Mass](#moment-of-inertia-rotational-mass)
- [Moment of Inertia for Common Shapes](#moment-of-inertia-for-common-shapes)
- [Rotational Kinetic Energy](#rotational-kinetic-energy)
- [Rolling Without Slipping](#rolling-without-slipping)
- [Worked Example 1: Torque Wrench](#worked-example-1-torque-wrench)
- [Worked Example 2: Spinning Disk Acceleration](#worked-example-2-spinning-disk-acceleration)
- [Worked Example 3: Rolling Cylinder Down a Ramp](#worked-example-3-rolling-cylinder-down-a-ramp)
- [Worked Example 4: Balancing a Seesaw](#worked-example-4-balancing-a-seesaw)
- [Real-World Applications](#real-world-applications)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: Why Things Spin

You have pushed a door open thousands of times in your life. But have you ever noticed something strange about where door handles are placed? They are almost always on the side opposite the hinges — as far from the pivot point as possible. Why?

Have you ever tried to loosen a very tight bolt? If you use a short wrench, you strain and struggle. Switch to a longer wrench and the bolt comes free with ease. Why does the length of the wrench matter?

Have you watched a figure skater pull their arms in and suddenly spin much faster? Or seen a tightrope walker holding a long pole to stay balanced? All of these phenomena involve the same concept: **torque**.

In the previous chapters, we studied how forces cause objects to accelerate in a straight line — **linear motion**. Now we turn to **rotational motion**: the physics of spinning, turning, and twisting. It turns out that rotation has its own complete set of laws that mirror Newton's laws for straight-line motion almost perfectly.

By the end of this chapter, you will understand:
- What torque is and how to calculate it
- Why leverage matters (and how Archimedes was right about moving the world)
- How Newton's second law applies to spinning objects
- What **moment of inertia** is and why a hollow sphere is harder to spin than a solid one
- How a rolling ball uses two kinds of kinetic energy at once

Let us begin.

---

## What Is Torque?

Imagine you are trying to open a very heavy jar lid. You grip the lid and twist. The force you apply causes the lid to **rotate** about its center. This tendency of a force to cause rotation is called **torque**.

**Torque** (symbol: τ, the Greek letter "tau") is the rotational equivalent of force. Just as a force causes a change in linear motion (acceleration in a straight line), torque causes a change in rotational motion (angular acceleration — spinning faster or slower).

But here is the key insight: **not all forces create torque equally**. A force applied far from a pivot point creates more torque than the same force applied close to the pivot point.

Think about it physically:

```
                    LEVER DIAGRAM
                    
    Force F applied here                Force F applied here
         |                                    |
         v                                    v
    ===========================O============
    |<---- long arm -------->| |<-short->|
                             ^
                         PIVOT (fulcrum)

    Same force F, but longer arm = MORE torque = easier to lift
```

The "arm" — the distance from the pivot to where the force is applied — is what makes the difference.

---

## The Lever Arm

The **lever arm** (also called the **moment arm**) is the perpendicular distance from the pivot point (also called the **axis of rotation** or **fulcrum**) to the line of action of the force.

This is a crucial definition. The lever arm is NOT simply the distance from pivot to where you push. It is the **perpendicular** distance from the pivot to the **line** along which the force acts.

Let us look at an example. Suppose you push on a wrench that is attached to a bolt:

```
    WRENCH DIAGRAM - Different force angles
    
    Case 1: Force perpendicular to wrench (maximum torque)
    
         F
         |  (force pointing upward)
         |
    -----+----------O  (bolt/pivot)
    |<--- r = 30 cm -->|
    
    Lever arm = r = 30 cm  (full length)
    
    -----------------------------------------------
    
    Case 2: Force at an angle θ
    
              F
             /
            / θ
    -------+----------O  (bolt/pivot)
    
    The perpendicular component is F×sin(θ)
    Lever arm = r×sin(θ)  (reduced!)
    
    -----------------------------------------------
    
    Case 3: Force along the wrench (zero torque!)
    
    F ------>----------O  (bolt/pivot)
    
    Force points directly at the pivot.
    Lever arm = 0
    No rotation happens at all!
```

This diagram reveals something important: **only the component of force perpendicular to the lever arm produces torque**. A force aimed directly at the pivot point produces zero torque no matter how hard you push.

---

## The Torque Formula

Now we can write the formula for torque:

```
    τ = F × r × sin(θ)
```

Where:
- τ = torque (measured in Newton-meters, N·m)
- F = magnitude of the force applied (in Newtons, N)
- r = distance from the pivot to the point where force is applied (in meters, m)
- θ = angle between the force vector and the lever (the radius vector pointing from pivot to force application point)

Let us check the extreme cases:
- When θ = 90°: sin(90°) = 1, so τ = F × r × 1 = F × r (maximum torque — force is perfectly perpendicular)
- When θ = 0°: sin(0°) = 0, so τ = 0 (no torque — force points directly at pivot)
- When θ = 180°: sin(180°) = 0, so τ = 0 (no torque — force points directly away from pivot)

This makes perfect intuitive sense.

An equivalent way to write the formula:

```
    τ = F_perp × r
    
    where F_perp = F × sin(θ) is the component of force perpendicular to r
    
    OR:
    
    τ = F × r_perp
    
    where r_perp = r × sin(θ) is the perpendicular distance from pivot to line of force
```

Both forms give the same answer. Use whichever is easier for a given problem.

**Units of torque**: Newton-meters (N·m). Note: this is the same unit as energy (Joules), but torque and energy are different physical quantities. Be careful not to confuse them.

---

## Direction of Torque

Torque is not just a number — it has a direction. In two-dimensional problems, we use a simple sign convention:

- **Counterclockwise (CCW) torque** is positive (+)
- **Clockwise (CW) torque** is negative (-)

```
    TORQUE DIRECTION CONVENTION
    
         CCW = Positive             CW = Negative
         
           F                              F
           |                              |
           | (+)                          | (-)
    -------O                    -------O
    
    (force on left side          (force on right side
     of pivot pushes             of pivot pulls down —
     lever to rotate CCW)        lever rotates CW)
```

In three dimensions, torque is actually a **vector** that points along the axis of rotation (using the right-hand rule), but for most problems in this chapter, the sign convention above is sufficient.

---

## Why Door Handles Are Far From Hinges

Now we can answer the question from the introduction. A door rotates about its hinges. The hinges act as the **pivot point**.

```
    DOOR DIAGRAM (top view)
    
    HINGE (pivot)
    |
    |========================| <- door
    |                        |
    |                    [HANDLE]
    |                        |
    |========================|
    
    |<------ r = 0.9 m ----->|
    
    The handle is placed at maximum distance from the hinge
    to maximize torque for a given force.
```

If you push with force F = 10 N perpendicular to the door:

- Handle at r = 0.9 m from hinge: τ = 10 × 0.9 × sin(90°) = **9 N·m**
- Same force at r = 0.1 m from hinge: τ = 10 × 0.1 × sin(90°) = **1 N·m**

The same 10 N force creates 9 times more torque when applied at the edge compared to near the hinge! This is why pushing a door near the hinge is so much harder. The door needs a minimum torque to overcome its friction with the latch; pushing near the hinge simply may not generate enough torque.

This is also why:
- Wrenches have long handles
- Steering wheels are large in diameter
- Bicycle pedals are attached to cranks (not directly to the axle)
- Bottle openers work using leverage

---

## Net Torque

Just as multiple forces can act on an object (we calculate net force), multiple torques can act on a rotating object. The **net torque** is the sum of all individual torques, with attention to sign (direction).

```
    τ_net = τ_1 + τ_2 + τ_3 + ...
```

When τ_net = 0, the object is in **rotational equilibrium** — it is either not rotating, or rotating at a constant rate.

**Example: Seesaw Balance**

```
    SEESAW DIAGRAM
    
         Child A (30 kg)              Child B (20 kg)
              |                              |
              v (weight = 300 N)             v (weight = 200 N)
    ==========+==============================+=========
    |<-- 2 m -->|          ^           |<------ 3 m ----->|
                        PIVOT
    
    Torque from Child A (clockwise, negative):
    τ_A = -(300 N)(2 m) = -600 N·m
    
    Torque from Child B (counterclockwise, positive):
    τ_B = +(200 N)(3 m) = +600 N·m
    
    τ_net = -600 + 600 = 0
    
    The seesaw is balanced!
```

---

## Newton's Second Law for Rotation

In Chapter 4, we learned Newton's second law for linear motion:

```
    F_net = m × a
    (net force = mass × linear acceleration)
```

There is a perfect rotational analogue:

```
    τ_net = I × α
    (net torque = moment of inertia × angular acceleration)
```

Where:
- τ_net = net torque in N·m
- I = **moment of inertia** in kg·m² (rotational equivalent of mass)
- α = **angular acceleration** in radians per second squared (rad/s²)

This equation is the cornerstone of rotational dynamics. It tells us:
- Greater net torque → greater angular acceleration (spins up faster)
- Greater moment of inertia → smaller angular acceleration (harder to spin up)

The analogy between linear and rotational motion is beautiful:

```
    LINEAR MOTION          ROTATIONAL MOTION
    ───────────────────────────────────────
    Force F         ↔      Torque τ
    Mass m          ↔      Moment of inertia I
    Acceleration a  ↔      Angular acceleration α
    Velocity v      ↔      Angular velocity ω
    Displacement x  ↔      Angle θ
    
    F = ma          ↔      τ = Iα
    KE = ½mv²       ↔      KE_rot = ½Iω²
    p = mv          ↔      L = Iω
```

---

## Moment of Inertia: Rotational Mass

**Moment of inertia** (symbol: I) is the rotational equivalent of mass. In linear motion, mass tells us how hard it is to change an object's linear velocity. Moment of inertia tells us how hard it is to change an object's **rotational velocity**.

But moment of inertia is more complicated than mass. It depends not just on how much mass an object has, but on **where that mass is located relative to the axis of rotation**.

Here is the key insight: **mass farther from the axis contributes more to the moment of inertia**.

For a single point mass m at distance r from the axis:

```
    I = m × r²
```

This r² dependence is critical. If you double the distance of a mass from the axis, the moment of inertia **quadruples**.

For a system of multiple point masses:

```
    I = m_1 × r_1² + m_2 × r_2² + m_3 × r_3² + ...
    
    I = Σ (m_i × r_i²)
```

**Physical intuition**: Think of two dumbbells.

```
    DUMBBELL COMPARISON
    
    Short dumbbell (masses close to center):
    
    O--[===]--O
    |<-0.2m->|
    
    Easy to spin — masses are close to the axis
    I is small
    
    Long dumbbell (masses far from center):
    
    O--------[===]--------O
    |<------1.0 m------->|
    
    Hard to spin — masses are far from the axis
    I is large
    
    Same total mass, but very different moments of inertia!
```

---

## Moment of Inertia for Common Shapes

For extended objects (not just point masses), we must add up the contributions from every tiny piece of mass across the whole object. This requires calculus, but the results are well-known formulas:

```
    MOMENT OF INERTIA FORMULAS
    (m = total mass, R = radius, L = length)
    
    ┌─────────────────────────────────────────────────────────┐
    │ Shape                      │ Axis              │  I     │
    ├─────────────────────────────────────────────────────────┤
    │ Point mass                 │ Distance r away   │ mr²    │
    │ Thin ring / hollow cylinder│ Central axis      │ mR²    │
    │ Solid disk / solid cylinder│ Central axis      │ ½mR²   │
    │ Solid sphere               │ Through center    │ 2/5 mR²│
    │ Hollow sphere (thin shell) │ Through center    │ 2/3 mR²│
    │ Thin rod                   │ Through center    │1/12 mL²│
    │ Thin rod                   │ Through one end   │ 1/3 mL²│
    └─────────────────────────────────────────────────────────┘
```

**Key observations from this table:**

1. A **hollow cylinder** (I = mR²) is harder to spin than a **solid cylinder** (I = ½mR²) of the same mass and radius. The hollow cylinder has all its mass at maximum distance from the axis.

2. A **hollow sphere** (I = 2/3 mR²) is harder to spin than a **solid sphere** (I = 2/5 mR²). Same reason.

3. A rod is much harder to spin about its end (1/3 mL²) than about its center (1/12 mL²) — the mass at the ends contributes much more when the pivot is at one end.

```
    VISUALIZING SOLID vs HOLLOW CYLINDER
    
    Solid cylinder:               Hollow cylinder (pipe):
    
       ___________                   _____________
      |   x x x   |                 |  |       |  |
      | x x x x x |                 |  |       |  |
      | x x x x x |                 |  |       |  |
      |   x x x   |                 |  |_______|  |
       ___________                   _____________
    
    Mass spread throughout          Mass concentrated at rim
    I = ½mR²                        I = mR²
    Easier to spin                  Harder to spin
```

---

## Rotational Kinetic Energy

A spinning object has kinetic energy due to its rotation. This is called **rotational kinetic energy**:

```
    KE_rot = ½ × I × ω²
```

Where:
- I = moment of inertia (kg·m²)
- ω = angular velocity (radians per second, rad/s)

Compare this to linear kinetic energy: KE_linear = ½mv². The structure is identical — just replace mass with moment of inertia and velocity with angular velocity.

**Example**: A flywheel (solid disk) of mass 20 kg and radius 0.5 m spins at 100 rad/s.

```
    I = ½mR² = ½ × 20 × (0.5)² = ½ × 20 × 0.25 = 2.5 kg·m²
    
    KE_rot = ½ × I × ω² = ½ × 2.5 × (100)² = ½ × 2.5 × 10000 = 12,500 J
```

That is 12.5 kilojoules of stored energy — enough to power a light bulb for over three minutes. This is why flywheels are used as energy storage devices in hybrid vehicles and power plants.

---

## Rolling Without Slipping

When a ball or cylinder rolls across a floor **without slipping**, it has two kinds of motion simultaneously:
1. **Translation** — the center of mass moves forward with velocity v
2. **Rotation** — the object spins about its own center with angular velocity ω

```
    ROLLING OBJECT DIAGRAM
    
               v (forward motion of center)
               -->
          ___________
         /    (+)    \      The center moves right at speed v
        |      O      |     The object spins clockwise
         \___________/      The contact point has zero velocity
         
         ^
         |
    (contact point — zero velocity for rolling without slipping)
    
    Key relationship: v = ω × R
    (linear speed = angular speed × radius)
```

The condition for **rolling without slipping** is:

```
    v = ω × R
```

This links the translational speed to the rotational speed. If a ball of radius R = 0.1 m rolls at v = 2 m/s, then ω = v/R = 2/0.1 = 20 rad/s.

The **total kinetic energy** of a rolling object is the sum of translational and rotational KE:

```
    KE_total = KE_translational + KE_rotational
    
    KE_total = ½mv² + ½Iω²
```

Using v = ωR (so ω = v/R):

```
    KE_total = ½mv² + ½I(v/R)²
    
    KE_total = ½mv² + ½(I/R²)v²
    
    KE_total = ½v²(m + I/R²)
```

For a solid cylinder (I = ½mR²):

```
    KE_total = ½v²(m + ½mR²/R²)
             = ½v²(m + ½m)
             = ½v²(3m/2)
             = ¾mv²
```

Compare to a sliding block (no rotation): KE = ½mv². The rolling cylinder needs more energy to reach the same speed — or equivalently, it moves slower than a sliding block after both descend the same ramp.

---

## Worked Example 1: Torque Wrench

**Problem**: A mechanic uses a torque wrench to tighten a bolt. The wrench handle is 0.35 m long. The mechanic applies a force of 80 N perpendicular to the wrench. What torque is applied to the bolt?

**Solution**:

Given:
- r = 0.35 m
- F = 80 N
- θ = 90° (force is perpendicular to wrench)

```
    DIAGRAM:
    
         F = 80 N (upward)
         |
         |
    =====+========O  (bolt)
    |<-- r = 0.35 m -->|
```

Using the torque formula:

```
    τ = F × r × sin(θ)
    τ = 80 N × 0.35 m × sin(90°)
    τ = 80 × 0.35 × 1
    τ = 28 N·m
```

**Answer**: The torque applied is 28 N·m.

**Follow-up**: The mechanic's service manual says the bolt requires 35 N·m. Can the mechanic achieve this with the same 80 N force?

```
    τ = F × r × sin(θ)
    35 = 80 × r × sin(90°)
    35 = 80 × r
    r = 35/80 = 0.4375 m
```

The mechanic needs a wrench at least 43.75 cm long. Time to get a longer wrench!

---

## Worked Example 2: Spinning Disk Acceleration

**Problem**: A solid disk (like a merry-go-round) has mass 50 kg and radius 1.5 m. A child pushes tangentially (perpendicular to the radius) with a force of 30 N. 

(a) What torque does the child apply?
(b) What angular acceleration does the disk experience?
(c) Starting from rest, what is the angular velocity after 4 seconds?

**Solution**:

Given:
- m = 50 kg
- R = 1.5 m
- F = 30 N (tangential, so θ = 90°)

**(a) Torque:**

```
    τ = F × R × sin(θ)
    τ = 30 × 1.5 × sin(90°)
    τ = 30 × 1.5 × 1
    τ = 45 N·m
```

**(b) Moment of inertia and angular acceleration:**

The disk is a solid disk, so:

```
    I = ½mR²
    I = ½ × 50 × (1.5)²
    I = ½ × 50 × 2.25
    I = 56.25 kg·m²
```

Using Newton's second law for rotation:

```
    τ_net = I × α
    45 = 56.25 × α
    α = 45 / 56.25
    α = 0.8 rad/s²
```

**(c) Angular velocity after 4 seconds:**

Using the rotational kinematic equation (analogous to v = v₀ + at):

```
    ω = ω₀ + α × t
    ω = 0 + 0.8 × 4
    ω = 3.2 rad/s
```

**Answer Summary**:
- Torque = 45 N·m
- Angular acceleration = 0.8 rad/s²
- Angular velocity after 4 s = 3.2 rad/s

To put this in perspective, 3.2 rad/s corresponds to about 30.6 revolutions per minute — a gentle but noticeable spin.

---

## Worked Example 3: Rolling Cylinder Down a Ramp

**Problem**: A solid cylinder of mass 2 kg and radius 0.1 m is released from rest at the top of a ramp that is 3 m long and makes an angle of 20° with the horizontal. The cylinder rolls without slipping. What is its speed at the bottom of the ramp?

```
    RAMP DIAGRAM:
    
    *
    * \
    *  \  3 m (ramp length)
    *   \
    *    \  20°
    *     \________
    
    Height h = 3 × sin(20°) = 3 × 0.342 = 1.026 m
```

**Solution**:

We use conservation of energy. At the top, all energy is gravitational potential energy. At the bottom, it is split between translational and rotational kinetic energy.

Step 1: Find the height

```
    h = L × sin(θ)
    h = 3 × sin(20°)
    h = 3 × 0.342
    h = 1.026 m
```

Step 2: Write the energy equation

```
    PE_initial = KE_trans + KE_rot
    mgh = ½mv² + ½Iω²
```

Step 3: Substitute I = ½mR² and ω = v/R

```
    mgh = ½mv² + ½(½mR²)(v/R)²
    mgh = ½mv² + ¼mv²
    mgh = (3/4)mv²
```

The mass m cancels:

```
    gh = (3/4)v²
    v² = (4/3)gh
    v² = (4/3) × 9.8 × 1.026
    v² = (4/3) × 10.055
    v² = 13.41
    v = √13.41
    v = 3.66 m/s
```

**Answer**: The cylinder reaches the bottom at 3.66 m/s.

**Comparison**: If the cylinder were instead a **hollow cylinder** (I = mR²), the calculation changes:

```
    mgh = ½mv² + ½(mR²)(v/R)²
    mgh = ½mv² + ½mv²
    mgh = mv²
    v² = gh = 9.8 × 1.026 = 10.055
    v = 3.17 m/s
```

The hollow cylinder reaches the bottom more slowly (3.17 m/s vs 3.66 m/s) because more of its mass is at the rim, giving it greater rotational inertia and more energy going into rotation rather than forward motion.

This is a real experiment you can try at home: race a solid ball against a hollow ball down a ramp. The solid ball always wins!

---

## Worked Example 4: Balancing a Seesaw

**Problem**: Three children sit on a seesaw pivoted at its center. Child A (weight 400 N) sits 2 m to the left. Child B (weight 250 N) sits 1.5 m to the right. Where must Child C (weight 300 N) sit for the seesaw to be balanced?

```
    SEESAW DIAGRAM:
    
         A (400N)    pivot    B (250N)    C (300N)
             |          |        |           ?
             v          |        v           v
    =========+============================+
    |<- 2m ->|          |    |<- 1.5m ->|<- x ->|
                        ^
                      PIVOT
```

**Solution**:

For rotational equilibrium, net torque = 0.

Taking counterclockwise as positive, and measuring from the pivot:

```
    Torque from A (CCW, left side, positive): τ_A = +400 × 2 = +800 N·m
    Torque from B (CW, right side, negative): τ_B = -250 × 1.5 = -375 N·m
    Torque from C (CW, right side, negative): τ_C = -300 × x
    
    For balance: τ_net = 0
    
    800 - 375 - 300x = 0
    425 = 300x
    x = 425/300
    x = 1.417 m
```

**Answer**: Child C must sit approximately 1.42 m to the right of the pivot.

---

## Real-World Applications

### The Human Body as a Lever System

Your body is full of lever systems. When you hold a weight in your hand with your forearm horizontal:

```
    FOREARM LEVER DIAGRAM
    
    Bicep muscle
    force (upward)
         |
         |  2 cm from elbow
    =====|=================================
    ELBOW                          HAND
    (pivot)                        holds 5 kg weight
    |<------ 35 cm (to hand) ----->|
    
    To hold 5 kg (weight ≈ 50 N):
    τ_weight = 50 × 0.35 = 17.5 N·m (clockwise)
    
    τ_muscle = F_bicep × 0.02 (counterclockwise)
    
    For balance: F_bicep × 0.02 = 17.5
    F_bicep = 17.5 / 0.02 = 875 N
    
    Your bicep exerts 875 N just to hold 50 N!
    That is why small lever arms in the body require huge muscle forces.
```

### Torque in Engines

Car and motorcycle engines produce **torque** (not just power). When you see "200 N·m of torque" in a car specification, that is the twisting force the engine applies to the drivetrain. Higher torque = better acceleration, especially from a standstill. This is why diesel trucks, which need to move heavy loads, are built for high torque at low engine speeds.

### Balance and Stability

A tightrope walker carries a long pole not for looks, but for physics. The long pole increases the walker's moment of inertia about the vertical axis. A large I means that any accidental rotation (starting to tip) produces only a small angular acceleration (τ = Iα, large I → small α for same τ). This gives the walker more time to correct their balance.

### Gyroscopes

Though we will discuss these more in the next chapter, note that gyroscopes resist changes to their axis of rotation because of their rotational inertia (moment of inertia). This makes them useful in aircraft navigation systems, phones, and satellite attitude control.

---

## Summary

- **Torque** (τ) is the rotational equivalent of force — it is the tendency of a force to cause rotation about a pivot point.

- The torque formula is: τ = F × r × sin(θ), where r is the distance from pivot to force application and θ is the angle between the force and the radius.

- Maximum torque occurs when the force is perpendicular to the lever (θ = 90°). Zero torque occurs when the force points directly at or away from the pivot (θ = 0° or 180°).

- Door handles are placed far from hinges to maximize torque for a given push force.

- **Newton's Second Law for Rotation**: τ_net = I × α. Net torque causes angular acceleration.

- **Moment of inertia** (I) is the rotational equivalent of mass. It depends on both the amount of mass and how far that mass is from the axis: I = Σ(mr²).

- Objects with mass concentrated far from the axis (hollow cylinder, hollow sphere) have larger moments of inertia than objects with mass spread throughout (solid cylinder, solid sphere) of the same total mass.

- **Rotational kinetic energy**: KE_rot = ½Iω²

- A rolling object has both translational KE (½mv²) and rotational KE (½Iω²). Total KE = ½mv² + ½Iω².

- For rolling without slipping: v = ω × R links the linear and angular velocities.

- A solid sphere or cylinder rolls down a ramp faster than a hollow one because its smaller I means less energy goes into rotation.

---

## Key Equations

```
    TORQUE:
    τ = F × r × sin(θ)
    
    Units: Newton-meters (N·m)
    
    ─────────────────────────────────────────
    
    NEWTON'S SECOND LAW FOR ROTATION:
    τ_net = I × α
    
    ─────────────────────────────────────────
    
    MOMENT OF INERTIA (POINT MASS):
    I = m × r²
    
    COMMON SHAPES:
    Solid disk/cylinder:  I = ½mR²
    Hollow ring/cylinder: I = mR²
    Solid sphere:         I = (2/5)mR²
    Hollow sphere:        I = (2/3)mR²
    Rod (center axis):    I = (1/12)mL²
    Rod (end axis):       I = (1/3)mL²
    
    ─────────────────────────────────────────
    
    ROTATIONAL KINETIC ENERGY:
    KE_rot = ½ × I × ω²
    
    ─────────────────────────────────────────
    
    TOTAL KE FOR ROLLING WITHOUT SLIPPING:
    KE_total = ½mv² + ½Iω²
    
    ROLLING CONDITION:
    v = ω × R
    
    ─────────────────────────────────────────
    
    ROTATIONAL KINEMATIC EQUATIONS:
    (analogous to linear equations)
    
    ω = ω₀ + α×t
    θ = ω₀×t + ½×α×t²
    ω² = ω₀² + 2×α×θ
    
    ─────────────────────────────────────────
    
    ANALOGY TABLE:
    Linear           Rotational
    F = ma    ↔      τ = Iα
    KE = ½mv² ↔      KE = ½Iω²
    p = mv    ↔      L = Iω
```
