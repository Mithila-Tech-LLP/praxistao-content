# Chapter 23: Collisions

> "In every interaction, something is always conserved. Nature keeps perfect books."
> — Anonymous

---

## Table of Contents

- [Introduction: The Universe Keeps Score](#introduction-the-universe-keeps-score)
- [What Is Momentum?](#what-is-momentum)
- [Conservation of Momentum: Proof from Newton's Third Law](#conservation-of-momentum-proof-from-newtons-third-law)
- [Types of Collisions](#types-of-collisions)
- [Elastic Collisions](#elastic-collisions)
- [Inelastic Collisions](#inelastic-collisions)
- [Perfectly Inelastic Collisions](#perfectly-inelastic-collisions)
- [Explosions: The Reverse Collision](#explosions-the-reverse-collision)
- [Collisions in 2D: A Brief Introduction](#collisions-in-2d-a-brief-introduction)
- [Newton's Cradle Explained](#newtons-cradle-explained)
- [Worked Example 1: Two Carts Colliding](#worked-example-1-two-carts-colliding)
- [Worked Example 2: Clay Ball Sticking](#worked-example-2-clay-ball-sticking)
- [Worked Example 3: Bullet Embedding in Block](#worked-example-3-bullet-embedding-in-block)
- [Worked Example 4: Explosion — Rocket Splits](#worked-example-4-explosion--rocket-splits)
- [Common Mistakes](#common-mistakes)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Introduction: The Universe Keeps Score

Imagine you are playing pool (billiards). The white cue ball rolls across the table and slams into a cluster of coloured balls. They scatter in every direction. Chaos, right?

Not quite. Underneath the apparent chaos, something remarkable happens: the total momentum of all the balls before the collision is exactly equal to the total momentum of all the balls after the collision. Every single time. Without exception.

This is one of the most powerful laws in all of physics: the **law of conservation of momentum**. It applies to billiard balls, car crashes, rocket launches, atomic collisions, and galaxies merging. It never breaks down.

In this chapter you will learn:

1. What momentum is and why it matters
2. Why momentum is always conserved (we will prove it from Newton's Third Law)
3. The three types of collisions and how to solve each one
4. How explosions are just collisions running in reverse
5. Why Newton's Cradle does what it does

No prior knowledge of collisions is required. If you know F = m × a and the idea of kinetic energy = ½mv², you have everything you need.

---

## What Is Momentum?

**Momentum** (symbol **p**) is the "quantity of motion" an object possesses. It depends on two things: how massive the object is, and how fast it is moving.

```
p = m × v
```

Where:
- p = momentum (kg·m/s)
- m = mass (kg)
- v = velocity (m/s)

Momentum is a **vector** — it has both size and direction. If an object moves to the right, its momentum points to the right. This matters enormously when objects collide head-on.

### Everyday Intuition for Momentum

Think about these two scenarios:

```
Scenario A: Tennis ball
  ●  ────────────────►  60 m/s
  mass = 0.06 kg
  momentum = 0.06 × 60 = 3.6 kg·m/s

Scenario B: Bowling ball
  ●  ──────────►  3 m/s
  mass = 7 kg
  momentum = 7 × 3 = 21 kg·m/s
```

The bowling ball has almost six times more momentum even though it moves twenty times slower. When either one hits your foot, you feel the difference.

### Why Momentum Is Different from Speed

Speed tells you how fast something moves. Momentum tells you how hard it is to stop. A massive truck moving at 10 m/s is much harder to bring to rest than a bicycle moving at 10 m/s. This "resistance to stopping" is exactly momentum.

### Units of Momentum

```
[p] = kg × (m/s) = kg·m/s

This can also be written as N·s (Newton-seconds)
because 1 N = 1 kg·m/s²
so 1 N·s = 1 kg·m/s² × s = 1 kg·m/s   ✓
```

---

## Conservation of Momentum: Proof from Newton's Third Law

Here is the beautiful part. We do not need to accept conservation of momentum as a mysterious rule handed down from physics gods. We can **prove** it directly from Newton's Third Law, which you already know.

### Newton's Third Law Reminder

Newton's Third Law states: when object A exerts a force on object B, object B simultaneously exerts an equal and opposite force on object A.

```
Object A  ←──── Force from B on A ────
Object B  ────── Force from A on B ───►

These forces are:
  • Equal in magnitude
  • Opposite in direction
  • Acting on DIFFERENT objects
```

### The Proof

Let two objects collide. Call them Object 1 and Object 2.

During the collision, Object 1 pushes on Object 2 with force **F₁₂** (force of 1 on 2).
By Newton's Third Law, Object 2 pushes back on Object 1 with force **F₂₁** = -**F₁₂**.

Both forces act for the same time interval Δt (they are simultaneous).

Now use Newton's Second Law in impulse form:

```
Impulse = Force × time = change in momentum

For Object 2:    F₁₂ × Δt = Δp₂  (momentum change of object 2)
For Object 1:    F₂₁ × Δt = Δp₁  (momentum change of object 1)

Since F₂₁ = -F₁₂:
    -F₁₂ × Δt = Δp₁
    F₁₂ × Δt = -Δp₁

Therefore:    Δp₂ = -Δp₁

Or rearranging:    Δp₁ + Δp₂ = 0

This means the total change in momentum is ZERO.
Total momentum does not change. It is conserved.  ✓
```

This is not a coincidence. This is a direct mathematical consequence of Newton's Third Law. As long as Newton's Third Law holds (and it holds for all contact forces), momentum is conserved.

### The Conservation Law Written Out

```
Total momentum before = Total momentum after

m₁v₁ + m₂v₂ = m₁v₁' + m₂v₂'

Where:
  v₁, v₂  = velocities BEFORE collision
  v₁', v₂' = velocities AFTER collision  (the prime symbol ' means "after")
```

### Important Condition: No External Forces

Momentum is conserved in a **closed system** — one with no net external force. During a collision, the collision forces (internal forces between the objects) are huge compared to external forces like friction or gravity, so for the brief instant of the collision, we can treat the system as closed.

```
VALID: Two carts on a frictionless track
VALID: Two cars in a crash (during the crash itself)
VALID: Atomic and nuclear collisions
VALID: Explosive separations

CAREFUL: If significant friction acts over a long time,
         momentum is NOT conserved.
```

---

## Types of Collisions

Not all collisions are the same. They differ in what happens to **kinetic energy** (KE = ½mv²).

Momentum is **always** conserved in any collision (given no external forces).
Kinetic energy is **sometimes** conserved and sometimes not.

```
┌─────────────────────────────────────────────────────────────┐
│                    TYPES OF COLLISIONS                      │
├──────────────────────┬──────────────┬───────────────────────┤
│      Type            │   Momentum   │   Kinetic Energy      │
├──────────────────────┼──────────────┼───────────────────────┤
│ Elastic              │  Conserved   │  Conserved            │
│ Inelastic            │  Conserved   │  NOT conserved (lost) │
│ Perfectly Inelastic  │  Conserved   │  Maximum loss         │
├──────────────────────┴──────────────┴───────────────────────┤
│ NOTE: "Explosion" is the reverse — KE is GAINED             │
└─────────────────────────────────────────────────────────────┘
```

The "missing" kinetic energy in inelastic collisions does not disappear — it converts to heat, sound, deformation of metal, etc. Energy is still conserved overall (First Law of Thermodynamics), but the *kinetic* energy decreases.

---

## Elastic Collisions

### What Makes a Collision Elastic?

An **elastic collision** is one where both momentum AND kinetic energy are perfectly conserved. The objects bounce off each other without any energy being lost to heat, sound, or deformation.

In practice, perfectly elastic collisions are rare. Billiard balls come very close. Gas molecules bouncing off each other are approximately elastic. At the atomic level, electron collisions can be perfectly elastic.

### The Two Equations for Elastic Collisions

Because both quantities are conserved, we have two equations:

```
Conservation of Momentum:
    m₁v₁ + m₂v₂ = m₁v₁' + m₂v₂'    ... (1)

Conservation of Kinetic Energy:
    ½m₁v₁² + ½m₂v₂² = ½m₁v₁'² + ½m₂v₂'²    ... (2)

Two equations. Two unknowns (v₁' and v₂'). Solvable!
```

### Derived Formulas for Elastic Collisions

After solving the two equations simultaneously (which involves some algebra), you get clean formulas:

```
        (m₁ - m₂)         2m₂
v₁' = ─────────── v₁  +  ──────── v₂
        (m₁ + m₂)         (m₁ + m₂)


         2m₁               (m₂ - m₁)
v₂' = ──────── v₁  +  ─────────── v₂
       (m₁ + m₂)         (m₁ + m₂)
```

These look scary. Let us check them in special cases.

### Special Case 1: Equal Masses (m₁ = m₂ = m)

```
v₁' = (m - m)/(m + m) × v₁ + 2m/(m + m) × v₂
    = 0 × v₁ + (2m/2m) × v₂
    = v₂

v₂' = 2m/(m + m) × v₁ + (m - m)/(m + m) × v₂
    = v₁ + 0
    = v₁
```

The objects swap velocities! This is exactly what billiard balls do — the rolling ball stops dead and the target ball rolls away at the original speed.

### Special Case 2: Target at Rest (v₂ = 0)

```
v₁' = (m₁ - m₂)/(m₁ + m₂) × v₁

v₂' = 2m₁/(m₁ + m₂) × v₁
```

If m₁ >> m₂ (heavy hits light), then v₁' ≈ v₁ (heavy barely slows) and v₂' ≈ 2v₁ (light shoots away at twice the speed).

If m₁ << m₂ (light hits heavy), then v₁' ≈ -v₁ (light bounces back!) and v₂' ≈ 0 (heavy barely moves).

```
ASCII: Light ball bouncing off heavy wall

    ●───►          ◄───●
    v               -v
    (before)        (after)

    Heavy wall barely moves.
```

---

## Inelastic Collisions

### The Real World

Most real-world collisions are **inelastic** — objects do not bounce off each other perfectly. Some kinetic energy is lost to:

- Heat (metal bending, friction)
- Sound (the "bang" you hear)
- Deformation (crumpled car doors)
- Internal vibrations (molecules jiggling)

In an inelastic collision:
- Momentum IS conserved ✓
- Kinetic energy is NOT conserved ✗ (some KE is lost)
- The objects do NOT stick together (that is perfectly inelastic)

```
BEFORE:                          AFTER:

   ●──────►     ●                  ●──►     ●──────►
   v₁           v₂=0               v₁'      v₂'
   (fast)       (still)            (slower) (moving)

Momentum: m₁v₁ = m₁v₁' + m₂v₂'   (conserved ✓)
KE:       ½m₁v₁² > ½m₁v₁'² + ½m₂v₂'²  (some lost ✓)
```

### How to Solve Inelastic Collision Problems

Unlike elastic collisions, you only have ONE conservation equation (momentum). So you need one more piece of information — usually the **coefficient of restitution** or the fact that it is perfectly inelastic (objects stick). We will focus on perfectly inelastic collisions.

---

## Perfectly Inelastic Collisions

### Objects Stick Together

A **perfectly inelastic collision** is the extreme case where the objects stick together and move as one combined mass after the collision. This is the maximum possible kinetic energy loss (while still conserving momentum).

Examples:
- A lump of clay hitting a stationary clay block
- A bullet embedding in a wooden block
- Two train cars coupling together
- A catcher catching a baseball (they move together briefly)

### The Equation

```
Before:
    Object 1: mass m₁, velocity v₁
    Object 2: mass m₂, velocity v₂

After:
    Combined mass (m₁ + m₂), moving together at velocity v_f

Conservation of Momentum:
    m₁v₁ + m₂v₂ = (m₁ + m₂) × v_f

Solve for v_f:
             m₁v₁ + m₂v₂
    v_f  =  ──────────────
               m₁ + m₂
```

This is just the **weighted average** of the two velocities. The more massive object pulls the final velocity toward its own initial velocity.

### How Much KE Is Lost?

```
KE lost = KE_before - KE_after

KE_before = ½m₁v₁² + ½m₂v₂²

KE_after = ½(m₁ + m₂)v_f²

KE lost = KE_before - KE_after
```

### ASCII Momentum Diagram for Perfectly Inelastic Collision

```
BEFORE:
─────────────────────────────────────────────►  (momentum scale)
│───────────────────────────│                   p₁ = m₁v₁
                             │──────────│        p₂ = m₂v₂ (if v₂=0, this is zero)

AFTER:
─────────────────────────────────────────────►
│──────────────────│                            p_f = (m₁+m₂)v_f

The total length (total momentum) is the SAME before and after.
KE depends on v², not v, so KE is NOT the same.
```

---

## Explosions: The Reverse Collision

### Running the Film Backward

An **explosion** is mathematically the reverse of a perfectly inelastic collision. Instead of two objects coming together, one object splits apart into two (or more) pieces.

The key insight: if the object starts at rest, the total initial momentum is ZERO. After the explosion, the total momentum must still be ZERO — the pieces fly off in opposite directions with equal and opposite momenta.

```
BEFORE:
    [====●====]
       at rest
       p_total = 0

AFTER:
◄─────────●         ●─────────►
   piece 1              piece 2
   momentum p₁          momentum p₂

For total momentum to stay zero:
    p₁ + p₂ = 0
    p₁ = -p₂
    m₁v₁ = -m₂v₂
```

### Examples of Explosions

1. **Gun firing a bullet**: the gun recoils backward as the bullet flies forward. The bullet has small mass but huge speed; the gun has huge mass but small recoil speed. Their momenta are equal and opposite.

2. **Rocket propulsion**: exhaust gases shoot backward, rocket shoots forward.

3. **Astronaut pushes off spacecraft**: astronaut goes one way, spacecraft goes the other.

4. **Grenade or bomb**: fragments scatter in all directions; if you add up all the momentum vectors, they cancel to zero.

### Where Does the Energy Come From?

In a collision, kinetic energy is LOST (to heat/deformation).
In an explosion, kinetic energy is GAINED (from chemical energy, compressed springs, nuclear energy, etc.).

The source of energy is whatever caused the explosion — gunpowder, compressed gas, chemical bonds.

### The Recoil Formula

```
If initial momentum = 0:

    m₁v₁ + m₂v₂ = 0

    v₁ = -m₂v₂ / m₁

The negative sign means they move in OPPOSITE directions.
```

---

## Collisions in 2D: A Brief Introduction

So far, we have only considered objects moving along a straight line (1D). In real life, objects often collide at angles.

The key principle: **momentum is conserved in every direction independently**.

This means we break momentum into x and y components and apply conservation separately to each.

```
            ↑ y

  Object 1 ──────────────►  Object 2 (at rest)
                  │
                  │  (collision at angle)
                  ↓

AFTER:
  Object 1: flies off at angle θ₁ above x-axis
  Object 2: flies off at angle θ₂ below x-axis
```

### The Two Component Equations

```
x-direction:  m₁v₁ + 0 = m₁v₁'cosθ₁ + m₂v₂'cosθ₂
y-direction:  0 + 0 = m₁v₁'sinθ₁ - m₂v₂'sinθ₂
```

The y-equation equals zero because the initial y-momentum was zero (both objects moving horizontally).

2D collisions appear in billiards, car accidents, particle physics, and sports. The same principles apply — just split into components.

---

## Newton's Cradle Explained

Newton's Cradle is that classic desktop toy with five metal balls hanging from strings. You pull one ball back and let it go. It swings down and hits the others. One ball flies off the other end. Pull two balls back, two fly off. Why?

```
        |  |  |  |  |
        |  |  |  |  |
        ●  ●  ●  ●  ●
   Ball 1  2  3  4  5

Step 1: Pull ball 1 back
     \  |  |  |  |
      \ |  |  |  |
       \●  ●  ●  ●  ●

Step 2: Release. Ball 1 swings down and hits ball 2.
        ●──►●  ●  ●  ●

Step 3: Ball 5 flies off the other end!
        ●  ●  ●  ●  ◄──●
                         ●
```

This seems magical. Why does ball 2 not fly off? Why ball 5?

### The Answer: Two Conservation Laws Working Together

We need BOTH momentum and kinetic energy to be conserved (the collisions between steel balls are nearly elastic).

For one ball hitting four at rest:
- Many combinations could conserve momentum alone
- Only ONE combination conserves both momentum AND kinetic energy simultaneously

That combination is: one ball flies off the other end.

### Why Two Balls In = Two Balls Out?

If you pull two balls (combined mass 2m, speed v) and they hit three at rest:
- Momentum: 2mv = ?
- KE: ½(2m)v² = ?

The only solution: two balls fly off at speed v. Any other combination fails to satisfy both equations.

The mathematics is elegant but the intuition is: nature must satisfy both bills simultaneously, and only one solution does.

---

## Worked Example 1: Two Carts Colliding

### Problem

Cart A has mass 2.0 kg and moves at 3.0 m/s to the right. Cart B has mass 3.0 kg and is stationary. They collide elastically. Find the velocities after the collision.

### Given Information

```
m₁ = 2.0 kg     v₁ = +3.0 m/s    (positive = rightward)
m₂ = 3.0 kg     v₂ = 0 m/s

Elastic collision: both momentum and KE are conserved.
```

### Before the Collision

```
        ──────────────────────────────────►
           Cart A               Cart B
          [  2 kg  ]──3 m/s──► [  3 kg  ] (at rest)

Total momentum before = 2.0 × 3.0 + 3.0 × 0 = 6.0 kg·m/s
Total KE before = ½ × 2.0 × 3.0² + 0 = 9.0 J
```

### Solution

Using the elastic collision formulas:

```
       (m₁ - m₂)            2m₂
v₁' = ─────────── × v₁  +  ─────── × v₂
       (m₁ + m₂)            (m₁+m₂)

     = (2.0 - 3.0)/(2.0 + 3.0) × 3.0  +  2×3.0/(2.0+3.0) × 0

     = (-1.0/5.0) × 3.0  +  0

     = -0.6 m/s


        2m₁               (m₂ - m₁)
v₂' = ──────── × v₁  +  ─────────── × v₂
      (m₁+m₂)             (m₁+m₂)

    = 2×2.0/(2.0+3.0) × 3.0  +  (3.0-2.0)/(5.0) × 0

    = (4.0/5.0) × 3.0  +  0

    = 2.4 m/s
```

### Results

```
Cart A: v₁' = -0.6 m/s  → moves LEFT (bounces back)
Cart B: v₂' = +2.4 m/s  → moves RIGHT
```

### After the Collision

```
        ◄──0.6 m/s──                       ──2.4 m/s──►
           Cart A                              Cart B
          [  2 kg  ]                          [  3 kg  ]
```

### Verification

```
Momentum check:
  Before: 2.0 × 3.0 + 3.0 × 0      = 6.0 kg·m/s
  After:  2.0 × (-0.6) + 3.0 × 2.4 = -1.2 + 7.2 = 6.0 kg·m/s  ✓

KE check:
  Before: ½ × 2.0 × 3.0²           = 9.0 J
  After:  ½ × 2.0 × 0.6² + ½ × 3.0 × 2.4²
        = ½ × 2.0 × 0.36 + ½ × 3.0 × 5.76
        = 0.36 + 8.64
        = 9.0 J  ✓  (elastic collision confirmed)
```

---

## Worked Example 2: Clay Ball Sticking

### Problem

A 0.5 kg ball of clay moves at 4.0 m/s and hits a stationary 1.5 kg clay block. They stick together. Find:
(a) The final velocity of the combined mass
(b) The kinetic energy lost in the collision

### Given Information

```
m₁ = 0.5 kg     v₁ = +4.0 m/s
m₂ = 1.5 kg     v₂ = 0 m/s

Perfectly inelastic collision: objects stick together.
```

### Before the Collision

```
        ●──4 m/s──►     [■■■■■■■■■]  (at rest)
       0.5 kg              1.5 kg

Momentum before = 0.5 × 4.0 + 1.5 × 0 = 2.0 kg·m/s
KE before = ½ × 0.5 × 4.0² = ½ × 0.5 × 16 = 4.0 J
```

### Solution: Part (a)

```
Conservation of momentum:
    m₁v₁ + m₂v₂ = (m₁ + m₂) × v_f

    0.5 × 4.0 + 1.5 × 0 = (0.5 + 1.5) × v_f

    2.0 = 2.0 × v_f

    v_f = 1.0 m/s
```

### After the Collision

```
               [●■■■■■■■■■]──1.0 m/s──►
                 2.0 kg combined mass
```

### Solution: Part (b)

```
KE after = ½ × (m₁ + m₂) × v_f²
         = ½ × 2.0 × 1.0²
         = 1.0 J

KE lost = KE before - KE after
        = 4.0 J - 1.0 J
        = 3.0 J

Percentage lost = 3.0/4.0 × 100% = 75%
```

### Interpretation

75% of the kinetic energy was converted to heat and sound when the clay deformed. This is typical for a perfectly inelastic collision — the maximum possible KE loss for a given initial condition.

```
MOMENTUM BAR CHART:
Before:  [═════════════════════════]  2.0 kg·m/s
After:   [═════════════════════════]  2.0 kg·m/s  (same ✓)

KE BAR CHART:
Before:  [████████████████]  4.0 J
After:   [████]  1.0 J
Lost:    [████████████]  3.0 J  → heat + sound
```

---

## Worked Example 3: Bullet Embedding in Block

### Problem

A 10-gram bullet travelling at 400 m/s embeds itself in a 2.0 kg wooden block sitting on a frictionless table. After the collision, how fast does the block (with bullet embedded) move?

Then, how much kinetic energy was lost?

### Given Information

```
m_bullet = 10 g = 0.010 kg     v_bullet = 400 m/s
m_block  = 2.0 kg              v_block  = 0 m/s

Perfectly inelastic (bullet embeds in block).
```

### Before the Collision

```
    ────────────────────────────────────►
    ·──400 m/s──►    [█████████████████]
    (bullet)             (block, still)
    0.010 kg               2.0 kg

Total momentum = 0.010 × 400 + 2.0 × 0 = 4.0 kg·m/s
Total KE = ½ × 0.010 × 400² = ½ × 0.010 × 160000 = 800 J
```

### Solution

```
Conservation of momentum (perfectly inelastic):
    m_bullet × v_bullet + m_block × 0 = (m_bullet + m_block) × v_f

    0.010 × 400 = (0.010 + 2.0) × v_f

    4.0 = 2.010 × v_f

    v_f = 4.0 / 2.010

    v_f ≈ 1.99 m/s  ≈ 2.0 m/s
```

### KE After

```
KE_after = ½ × 2.010 × 1.99²
         = ½ × 2.010 × 3.96
         = 3.98 J  ≈ 4.0 J

KE lost = 800 J - 4.0 J = 796 J

Percentage lost = 796/800 × 100% ≈ 99.5%
```

### After the Collision

```
    [·█████████████████]──2.0 m/s──►
    (bullet + block together)
    2.010 kg
```

### Why So Much KE Is Lost

Almost 99.5% of the kinetic energy disappeared! This energy went into:
- Heating the wood and bullet (significant)
- Deforming the wood fibres
- Sound (the "bang")

This is characteristic of a tiny, fast object hitting a large, stationary one. The bullet barely changes the block's velocity, but it loses almost all of its own KE.

```
KE BEFORE: ████████████████████████████████████  800 J
KE AFTER:  ■  4 J
LOST:      ████████████████████████████████████  796 J → HEAT + SOUND
```

---

## Worked Example 4: Explosion — Rocket Splits

### Problem

A rocket of mass 500 kg is floating in deep space, completely stationary. A small explosive charge fires, splitting it into two pieces:
- Piece A: 200 kg
- Piece B: 300 kg

After the explosion, Piece A moves at 15 m/s in the forward direction. Find the velocity of Piece B.

Then calculate the kinetic energy gained in the explosion and explain where it came from.

### Given Information

```
Initial state: total mass = 500 kg, at rest → v_initial = 0
Piece A: m_A = 200 kg, v_A = +15 m/s (forward)
Piece B: m_B = 300 kg, v_B = ?
```

### Before the Explosion

```
         [─────●─────]
          500 kg rocket
          at rest
          p_total = 0
```

### Solution

```
Conservation of momentum:
    p_initial = p_A + p_B

    0 = m_A × v_A + m_B × v_B

    0 = 200 × 15 + 300 × v_B

    0 = 3000 + 300 × v_B

    300 × v_B = -3000

    v_B = -10 m/s
```

The negative sign means Piece B moves in the **backward** direction.

### After the Explosion

```
◄────10 m/s────              ────15 m/s────►
    Piece B                      Piece A
    300 kg                       200 kg
```

### Verification of Momentum

```
p_A = 200 × 15  = +3000 kg·m/s
p_B = 300 × (-10) = -3000 kg·m/s
Total = 3000 + (-3000) = 0 kg·m/s  ✓
```

### KE Gained

```
KE before = 0 J (everything at rest)

KE after = ½ × 200 × 15² + ½ × 300 × 10²
         = ½ × 200 × 225  +  ½ × 300 × 100
         = 22500 J + 15000 J
         = 37500 J
         = 37.5 kJ

KE gained = 37.5 kJ - 0 = 37.5 kJ
```

### Where Did the Energy Come From?

The 37.5 kJ of kinetic energy came from the chemical potential energy stored in the explosive charge. Before the explosion, that energy was locked in chemical bonds. The explosion released it, converting it into the kinetic energy of the two pieces.

```
ENERGY TRANSFORMATION:

Chemical Energy    ──explosion──►    Kinetic Energy
    37.5 kJ                            37.5 kJ
    (in explosive)                  (in the two pieces)
```

---

## Common Mistakes

### Mistake 1: Forgetting Direction (Signs)

Momentum is a vector. If one object moves left and another moves right, they have opposite signs.

```
WRONG: p_total = 3 × 4 + 2 × 3 = 18 kg·m/s
       (both treated as positive)

RIGHT: p_total = 3 × (+4) + 2 × (-3) = 12 - 6 = 6 kg·m/s
       (leftward velocity is negative)
```

Always define a positive direction and stick to it.

### Mistake 2: Confusing Elastic and Inelastic

In an elastic collision, objects bounce off and KE is conserved.
In a perfectly inelastic collision, objects stick together. They do NOT necessarily stop.

### Mistake 3: Using KE Conservation for Inelastic Collisions

You CANNOT use KE conservation for inelastic collisions. Use ONLY momentum conservation.

### Mistake 4: Wrong Mass in Perfectly Inelastic

After objects stick, use the **combined** mass for the final KE calculation.

```
WRONG: KE_after = ½ × m₁ × v_f²
RIGHT: KE_after = ½ × (m₁ + m₂) × v_f²
```

### Mistake 5: Forgetting to Check Your Answer

Always verify by plugging back in. If momentum before ≠ momentum after, you made an error.

---

## Summary

- **Momentum** is defined as p = m × v. It is a vector quantity with units kg·m/s.

- **Conservation of momentum** follows directly from Newton's Third Law. When two objects collide, the equal and opposite internal forces produce equal and opposite changes in momentum that cancel out exactly.

- The conservation law states: m₁v₁ + m₂v₂ = m₁v₁' + m₂v₂'. Momentum before equals momentum after.

- **Elastic collisions** conserve both momentum and kinetic energy. Objects bounce off each other. Example: billiard balls, gas molecules.

- **Inelastic collisions** conserve momentum but lose kinetic energy (to heat, sound, deformation). Most real collisions are inelastic.

- **Perfectly inelastic collisions** are the extreme case where objects stick together: v_f = (m₁v₁ + m₂v₂) / (m₁ + m₂). Maximum KE is lost.

- **Explosions** are reverse collisions. If the initial object is at rest, the fragments fly off with equal and opposite momenta. The kinetic energy gained comes from a stored energy source (chemical, nuclear, spring).

- In **2D collisions**, momentum is conserved separately in each direction (x and y components).

- **Newton's Cradle** works because both momentum AND kinetic energy must be conserved simultaneously. Only the "one in, one out" solution satisfies both equations.

- KE lost in a collision does not disappear — it converts to heat, sound, and deformation energy (First Law of Thermodynamics is never violated).

---

## Key Equations

```
MOMENTUM:
    p = m × v                              (definition of momentum)
    [p] = kg·m/s = N·s

CONSERVATION OF MOMENTUM (any collision):
    m₁v₁ + m₂v₂ = m₁v₁' + m₂v₂'

ELASTIC COLLISION:
    Also: ½m₁v₁² + ½m₂v₂² = ½m₁v₁'² + ½m₂v₂'²  (KE also conserved)

    Velocity formulas:
        v₁' = [(m₁-m₂)/(m₁+m₂)] × v₁  +  [2m₂/(m₁+m₂)] × v₂
        v₂' = [2m₁/(m₁+m₂)] × v₁  +  [(m₂-m₁)/(m₁+m₂)] × v₂

PERFECTLY INELASTIC COLLISION:
    v_f = (m₁v₁ + m₂v₂) / (m₁ + m₂)

EXPLOSION (object starts at rest):
    0 = m₁v₁' + m₂v₂'
    v₁' = -m₂v₂' / m₁

KE LOST IN COLLISION:
    ΔKE = KE_after - KE_before   (negative means KE was lost)

IMPULSE-MOMENTUM THEOREM:
    F × Δt = Δp = m × Δv
```
