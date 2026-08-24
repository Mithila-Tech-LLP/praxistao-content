# Chapter 18: Work and Energy

> "Energy cannot be created or destroyed; it can only be changed from one form to another." — Albert Einstein

---

## Table of Contents

1. [Two Ways to Use the Word "Work"](#1-two-ways-to-use-the-word-work)
2. [The Physics Definition of Work](#2-the-physics-definition-of-work)
3. [What Does the Angle θ Mean?](#3-what-does-the-angle-θ-mean)
4. [When Is Work Zero?](#4-when-is-work-zero)
5. [When Is Work Negative?](#5-when-is-work-negative)
6. [Worked Examples — Calculating Work](#6-worked-examples--calculating-work)
7. [Energy — The Capacity to Do Work](#7-energy--the-capacity-to-do-work)
8. [Kinetic Energy — The Energy of Motion](#8-kinetic-energy--the-energy-of-motion)
9. [The Work-Energy Theorem](#9-the-work-energy-theorem)
10. [Worked Examples — The Work-Energy Theorem](#10-worked-examples--the-work-energy-theorem)
11. [Why Work and Energy Are Two Faces of the Same Thing](#11-why-work-and-energy-are-two-faces-of-the-same-thing)
12. [Preview — Potential Energy and Mechanical Energy](#12-preview--potential-energy-and-mechanical-energy)
13. [Summary](#summary)
14. [Key Equations](#key-equations)

---

## 1. Two Ways to Use the Word "Work"

In everyday speech, "work" can mean almost anything effortful:
- "I worked hard today."
- "That plan won't work."
- "She works as a doctor."

In physics, the word "work" has a very specific, precise, and somewhat unexpected meaning. A doctor who sits still and thinks all day is not doing any "work" in the physics sense — even though it feels exhausting. A robot that pushes a crate across a warehouse floor IS doing work, even though the robot feels nothing.

This is one of those moments in physics where everyday language and physics language CONFLICT. Be aware of the conflict and use the physics definition consistently.

Here is the key difference:

```
   Everyday "work":    Any effort, mental or physical. Vague.
   
   Physics "work":     A specific quantity. Requires:
                       (1) A force applied to an object
                       (2) The object actually moves
                       (3) The force has a component in the direction of movement
```

If ANY of those three ingredients is missing, physics says no work is done.

---

## 2. The Physics Definition of Work

**Work** is defined as:

  W = F × d × cos θ

Where:
- W = work done (unit: Joule, J)
- F = magnitude of the applied force (unit: Newton, N)
- d = displacement of the object (unit: meters, m)
- θ = the angle between the force vector and the displacement vector

The unit of work is the **Joule** (symbol J):
  1 Joule = 1 Newton × 1 meter = 1 N·m

One Joule is roughly the work done lifting a small apple (about 100 grams) by 1 meter. It is a small amount of energy. Human activities typically involve thousands or millions of Joules.

### Visualizing the Formula

```
   Object moves this way: ----------------------->   displacement d
   
   Force is applied at angle θ to the displacement:
   
              Force F
             /
            /  θ (angle between force and displacement)
           /
          /___________________________________>
                        displacement d

   Only the COMPONENT of force in the direction of motion does work.
   That component = F × cos θ

   So:  W = (F cos θ) × d
```

This is equivalent to W = F × d × cos θ (same thing, different grouping).

---

## 3. What Does the Angle θ Mean?

The angle θ is the angle between two vectors:
- The direction of the applied force
- The direction of the object's displacement (the direction it moves)

Let us examine each case carefully.

### Case 1 — Force is in the SAME direction as motion (θ = 0°)

```
   F ---------> (force pointing right)
   Object -----> (moving right)
   
   θ = 0°
   cos 0° = 1
   W = F × d × 1 = F × d
```

This is the maximum work for a given force and displacement. Every bit of the force contributes to moving the object.

**Example:** Pushing a shopping cart horizontally with a horizontal push.

---

### Case 2 — Force is PERPENDICULAR to motion (θ = 90°)

```
   F ↑ (force pointing upward)
   Object -----> (moving horizontally)
   
   θ = 90°
   cos 90° = 0
   W = F × d × 0 = 0
```

**Zero work is done.** The force is entirely perpendicular to the motion — it does not help or hinder the movement at all.

**Example:** Gravity acts downward on a car driving horizontally. Gravity does zero work on the car while it drives on a flat road.

---

### Case 3 — Force is at an ACUTE ANGLE to motion (0° < θ < 90°)

```
         F
        /
       / θ
      /___>   displacement
   
   cos θ is between 0 and 1.
   W = F × d × cos θ   (positive work, less than maximum)
```

Only the horizontal component of the force (F cos θ) contributes to moving the object forward.

**Example:** Pulling a suitcase with a handle at an angle. The upward component of pull lifts slightly; only the horizontal component moves it forward.

---

### Case 4 — Force is in the OPPOSITE direction to motion (θ = 180°)

```
   <--------- F (force pointing left)
   Object -----> (moving right)
   
   θ = 180°
   cos 180° = -1
   W = F × d × (-1) = -Fd
```

The work is **negative**. The force opposes the motion. This is what friction does — it acts opposite to the direction of motion and does negative work, removing energy from the moving object and converting it to heat.

---

### The cos θ Rule Summarized

```
   θ = 0°    (same direction):      cos θ = +1.0   Maximum positive work
   θ = 30°:                          cos θ = +0.87  Positive work
   θ = 60°:                          cos θ = +0.50  Positive work
   θ = 90°   (perpendicular):        cos θ =  0.0   Zero work
   θ = 120°:                         cos θ = -0.50  Negative work
   θ = 150°:                         cos θ = -0.87  Negative work
   θ = 180°  (opposite direction):   cos θ = -1.0   Maximum negative work
```

---

## 4. When Is Work Zero?

Physics has a very unforgiving answer: work is zero in more situations than you might think.

### Situation 1 — The object does not move

If d = 0, then W = F × 0 × cos θ = 0.

A weightlifter holding a barbell completely still above their head is doing ZERO work on the barbell in the physics sense. Their muscles are working hard (biologically), but the barbell does not move, so no physics work is done on it.

```
   Weightlifter:
   
      [barbell]   ← held at rest
         |
      [arms]      ← exerting large upward force
      
   d = 0 (no displacement)
   W = F × 0 × cos θ = 0
```

Your body is doing metabolic work (burning calories, firing muscles), but in the physics sense, no work is being done on the barbell. The physics concept of work is about energy transfer to the object, not about how tired you get.

---

### Situation 2 — Force is perpendicular to motion (θ = 90°)

**Example 1 — Carrying a suitcase horizontally:**

```
   You carry a suitcase horizontally down a corridor.
   
   Force (you exert on suitcase):   F ↑ (upward, to support the weight)
   Displacement:                    -----> (horizontal)
   
   θ = 90° between upward force and horizontal motion
   W = F × d × cos 90° = 0
```

Gravity also does zero work here (gravity is downward, displacement is horizontal: also θ = 90°).

You are not doing any work on the suitcase in the horizontal direction! (You ARE doing work when you pick it up or put it down — that is vertical motion with a vertical force.)

**Example 2 — The Moon orbiting Earth:**

```
   Gravity pulls the Moon TOWARD Earth (inward, centripetal direction).
   The Moon moves ALONG its circular path (tangential direction).
   
   These two directions are always perpendicular: θ = 90° always.
   
   W_gravity = 0 for circular orbit.
```

Gravity does zero net work on the Moon in a circular orbit. The Moon's speed never changes because no net work is done. This is consistent with the fact that the Moon has been orbiting for billions of years with no energy input.

(Note: Real orbits are elliptical. In elliptical orbits, gravity does positive work when the planet moves inward and negative work when it moves outward. The net work over one full orbit is still zero.)

---

### Situation 3 — Normal force on a horizontal surface

When a box slides horizontally on a flat floor, the normal force from the floor acts upward. The box moves horizontally. Again θ = 90°, so the normal force does zero work. The normal force changes the direction the floor contacts the box — it does not change the box's energy.

---

## 5. When Is Work Negative?

Negative work means the force is removing energy from the object — slowing it down, or working against the direction of motion.

### Friction

When a box slides to the right across a rough floor:
- Kinetic friction points to the LEFT (opposing motion).
- Displacement is to the RIGHT.
- θ = 180° between friction force and displacement.
- W_friction = -f × d (negative work).

Friction steals kinetic energy from the box and converts it to heat. The box slows down. Energy is NOT destroyed — it just changes form from kinetic energy to thermal energy.

### Gravity on a Rising Object

When you throw a ball upward:
- Gravity points DOWN.
- Displacement points UP.
- θ = 180°.
- W_gravity = -mg × h (negative work).

Gravity does negative work on the rising ball, slowing it down. The ball loses kinetic energy and gains gravitational potential energy.

### Braking in a Car

When you apply the brakes, the braking force opposes the car's motion. The brakes do negative work on the car. The car's kinetic energy decreases (it slows down) and converts to heat (you can smell hot brakes after heavy braking).

---

## 6. Worked Examples — Calculating Work

### Worked Example 6.1 — Horizontal Push

**Problem:** You push a 10 kg box across a flat floor with a horizontal force of 40 N. You push it 5 m. How much work do you do on the box?

```
   F = 40 N -------->[BOX]-----------> displacement = 5 m
   
   θ = 0° (force and displacement in same direction)
```

  W = F × d × cos θ
  W = 40 × 5 × cos 0°
  W = 40 × 5 × 1
  W = 200 J

**Answer:** You do 200 J of work on the box.

---

### Worked Example 6.2 — Push at an Angle

**Problem:** Same 10 kg box. You push with 40 N at 30° below horizontal. Box moves 5 m horizontally. How much work do you do?

```
   Your push:
         \
          \ 40 N
           \
            \ 30°
             \_____________> displacement = 5 m
   
   θ = 30° (angle between force direction and displacement)
```

  W = F × d × cos θ
  W = 40 × 5 × cos 30°
  W = 40 × 5 × 0.866
  W = 200 × 0.866
  W = 173 J

**Answer:** You do 173 J of work — less than in Example 6.1, because part of your 40 N push is directed downward (into the floor) and does not contribute to horizontal movement.

**What force actually moves the box horizontally?**
  Horizontal component = F cos 30° = 40 × 0.866 = 34.6 N.
  Vertical component = F sin 30° = 40 × 0.5 = 20 N (pushes box into floor, just increases normal force).

Only the 34.6 N horizontal component does work on the box.

---

### Worked Example 6.3 — Lifting a Box (Work Against Gravity)

**Problem:** You carry a 20 kg box up a staircase. The vertical height gained is 3 m. How much work do you do against gravity?

```
   You carry the box:
   
        [BOX] at top
          |
          | h = 3 m (vertical height)
          |
        [BOX] at bottom
   
   Force needed to support box: F = mg = 20 × 9.8 = 196 N (upward)
   Displacement: 3 m upward
   θ = 0° between lifting force and upward displacement
```

  W = F × d × cos θ
  W = mg × h × cos 0°
  W = 196 × 3 × 1
  W = 588 J

**Answer:** You do 588 J of work against gravity.

**Important note:** This 588 J is the MINIMUM work you need to do — it is just enough to lift the box. In reality, you also do extra work overcoming friction in your body, maintaining balance, etc. But 588 J accounts purely for the work against gravity.

Also note: it does not matter whether you take the stairs in a straight line, zigzag, or take three separate flights. The work done against gravity depends only on the vertical height gained, not the path taken. This is a deep and important fact that we will explore further when we study potential energy.

---

### Worked Example 6.4 — Work Done by Friction

**Problem:** A box of mass 5 kg slides 3 m across a floor. The coefficient of kinetic friction is μk = 0.3. How much work does friction do on the box?

**Find the friction force:**
  N = mg = 5 × 9.8 = 49 N
  f = μk × N = 0.3 × 49 = 14.7 N

**Friction opposes motion:** θ = 180° between friction force and displacement.

  W_friction = F × d × cos θ
  W_friction = 14.7 × 3 × cos 180°
  W_friction = 14.7 × 3 × (-1)
  W_friction = -44.1 J

**Answer:** Friction does -44.1 J of work on the box. The negative sign means friction removes 44.1 J of kinetic energy from the box, slowing it down.

---

### Worked Example 6.5 — Multiple Forces: Net Work

**Problem:** A box is pulled 4 m horizontally by a force F = 30 N (horizontal). Friction opposes motion with f = 10 N. Find the net work done on the box.

**Work by F (θ = 0°):**
  W_F = 30 × 4 × cos 0° = 120 J

**Work by friction (θ = 180°):**
  W_f = 10 × 4 × cos 180° = -40 J

**Work by gravity (θ = 90°):**
  W_gravity = 0 J

**Work by normal force (θ = 90°):**
  W_N = 0 J

**Net work:**
  W_net = W_F + W_f + W_gravity + W_N
  W_net = 120 + (-40) + 0 + 0
  W_net = 80 J

**Answer:** The net work done on the box is 80 J. The box gains 80 J of kinetic energy.

---

## 7. Energy — The Capacity to Do Work

**Energy** is defined as the capacity (ability) to do work. An object that has energy can do work on other objects.

Energy is measured in the same units as work: Joules (J).

Energy is NOT a thing you can touch. It is a property of a system — a number that describes how much work the system is capable of doing.

### Forms of Energy

Energy comes in many forms, and it can convert from one form to another:

```
   Form               Example
   -----------------------------------------------------------------------
   Kinetic            A moving car, a flying ball, flowing water
   Gravitational      A raised weight, water behind a dam
   potential
   Elastic            A compressed spring, a stretched rubber band
   potential
   Chemical           Food, fuel, batteries, gunpowder
   Thermal            Hot steam, warm water, friction heat
   Electromagnetic    Light, radio waves, X-rays
   Nuclear            Energy in atomic nuclei (nuclear power plants)
   Sound              Vibrating air from a speaker
   -----------------------------------------------------------------------
```

All these forms are related. You can convert one to another. A car engine converts chemical energy (fuel) → thermal energy (combustion) → kinetic energy (car moving). A solar panel converts electromagnetic energy (light) → electrical energy → (eventually) kinetic or thermal energy.

### The Law of Conservation of Energy

**Energy cannot be created or destroyed. It can only be converted from one form to another.**

This is one of the most fundamental laws in all of physics. No exception has ever been found. The total energy of an isolated system is always constant.

When friction slows a box and the box "loses" kinetic energy, the energy does not disappear — it becomes thermal energy (the floor and box get slightly warmer). The total is conserved.

---

## 8. Kinetic Energy — The Energy of Motion

A moving object has energy by virtue of its motion. This is called **kinetic energy** (KE), from the Greek word "kinesis" meaning motion.

  KE = (1/2) × m × v²

Where:
- KE = kinetic energy (Joules)
- m = mass (kg)
- v = speed (m/s)

### Important Properties of Kinetic Energy

**1. It depends on v squared.**

If you double the speed, the kinetic energy quadruples (2² = 4). This is why car accidents at high speed are so much more dangerous than low-speed collisions. A car at 100 km/h has four times the kinetic energy of a car at 50 km/h.

```
   Speed (km/h)   Relative KE
   ------------------------------
   10             1× (reference)
   20             4×
   30             9×
   40             16×
   100            100×
```

**2. It is always positive (or zero).**

Mass is positive. v² is positive. So KE ≥ 0 always. An object at rest has KE = 0.

**3. It is a scalar.**

Kinetic energy has no direction — just a number (in Joules). A car moving east at 20 m/s has the same kinetic energy as a car moving west at 20 m/s.

---

### Quick Kinetic Energy Examples

```
   Object            Mass       Speed          KE
   -----------------------------------------------------------
   Slow walking person  70 kg    1.4 m/s (5 km/h)   69 J
   Running person       70 kg    6 m/s (21 km/h)    1260 J
   Family car           1200 kg  14 m/s (50 km/h)   117,600 J
   Family car           1200 kg  28 m/s (100 km/h)  470,400 J
   Rifle bullet         0.01 kg  900 m/s            4,050 J
   -----------------------------------------------------------
```

Interesting: a tiny bullet moving extremely fast can have more kinetic energy than a walking person even though the person is thousands of times heavier — because KE depends on v².

---

## 9. The Work-Energy Theorem

Here is one of the most powerful and elegant results in mechanics:

**The net work done on an object equals the change in its kinetic energy.**

  W_net = ΔKE = KE_final - KE_initial = (1/2)mv² - (1/2)mu²

Where:
- W_net = total (net) work done by ALL forces combined
- u = initial speed
- v = final speed

This is the **Work-Energy Theorem**.

### Deriving It (Optional — Follow the Algebra)

Starting from Newton's Second Law: F_net = ma

Using kinematics: v² = u² + 2as → as = (v² - u²)/2

Net work: W_net = F_net × d = ma × d = m × (as) = m × (v² - u²)/2 = (1/2)mv² - (1/2)mu²

So W_net = (1/2)mv² - (1/2)mu² = ΔKE. Proven.

### What This Theorem Tells Us

1. If net positive work is done on an object, its kinetic energy INCREASES (it speeds up).
2. If net negative work is done on an object, its kinetic energy DECREASES (it slows down).
3. If net zero work is done, kinetic energy is UNCHANGED (constant speed).

This gives us a powerful alternative method for solving dynamics problems. Instead of using F = ma (which requires finding acceleration and then using kinematics), we can use energy methods directly.

---

### The Work-Energy Theorem: A Conceptual Picture

```
   Think of kinetic energy as water in a bucket.
   
   Positive work (a push in the direction of motion) = adding water.
   Negative work (friction, braking) = removing water.
   
   The NET change in the water level = net work done.
   
   Initial    Work        Final
   bucket     done        bucket
   ------    ------       ------
   [  KE ]  + W_net  =  [ KE'  ]
   [ = 50J]  [+30J ]    [ = 80J]
   
   The bucket (kinetic energy) went from 50 J to 80 J.
   Exactly 30 J of work was done on the object.
```

---

## 10. Worked Examples — The Work-Energy Theorem

### Worked Example 10.1 — Finding Final Speed

**Problem:** A 2 kg ball starts from rest. A net work of 36 J is done on it. Find the final speed.

**Using the Work-Energy Theorem:**
  W_net = (1/2)mv² - (1/2)mu²
  36 = (1/2)(2)(v²) - (1/2)(2)(0²)
  36 = v² - 0
  v² = 36
  v = 6 m/s

**Answer:** The ball reaches a final speed of 6 m/s.

---

### Worked Example 10.2 — Finding How Far an Object Slides Before Stopping

**Problem:** A 3 kg box is moving at 8 m/s across a flat floor. The coefficient of kinetic friction is μk = 0.4. How far does the box slide before stopping?

**Step 1 — Find the friction force:**
  N = mg = 3 × 9.8 = 29.4 N
  f = μk × N = 0.4 × 29.4 = 11.76 N

**Step 2 — Work done by friction over distance d:**
  W_friction = -f × d = -11.76d     (negative because it opposes motion)

No other horizontal forces. Net work = -11.76d.

**Step 3 — Apply Work-Energy Theorem:**
  W_net = (1/2)mv² - (1/2)mu²
  -11.76d = (1/2)(3)(0)² - (1/2)(3)(8)²    [final speed v = 0]
  -11.76d = 0 - 96
  -11.76d = -96
  d = 96 / 11.76
  d ≈ 8.16 m

**Answer:** The box slides approximately 8.16 m before stopping.

**Check using kinematics (alternative approach):**
  a = -f/m = -11.76/3 = -3.92 m/s²  (deceleration)
  v² = u² + 2ad:  0 = 64 + 2(-3.92)d → d = 64/7.84 ≈ 8.16 m. Matches!

---

### Worked Example 10.3 — Box Being Pushed With Friction

**Problem:** A 5 kg box is pushed 6 m across a floor by a horizontal force F = 25 N. Friction force is 10 N. The box starts from rest. Find its final speed.

**Step 1 — Net work:**
  W_F = 25 × 6 = 150 J          (positive: force in direction of motion)
  W_friction = -10 × 6 = -60 J  (negative: friction opposes motion)
  W_net = 150 - 60 = 90 J

**Step 2 — Apply Work-Energy Theorem:**
  W_net = (1/2)mv² - (1/2)mu²
  90 = (1/2)(5)(v²) - 0
  90 = 2.5v²
  v² = 36
  v = 6 m/s

**Answer:** The box reaches a final speed of 6 m/s.

---

### Worked Example 10.4 — Ball Thrown Upward

**Problem:** A 0.2 kg ball is thrown straight up with an initial speed of 10 m/s. Using the Work-Energy Theorem, find how high it rises before momentarily stopping.

**Analysis:**
- Only gravity does work (no friction in the air for this problem).
- W_gravity = -mg × h = -0.2 × 9.8 × h = -1.96h  (negative: gravity opposes upward motion)

**Apply Work-Energy Theorem:**
  W_net = (1/2)mv² - (1/2)mu²
  -1.96h = (1/2)(0.2)(0)² - (1/2)(0.2)(10)²    [final v = 0 at top]
  -1.96h = 0 - 10
  -1.96h = -10
  h = 10 / 1.96
  h ≈ 5.1 m

**Answer:** The ball rises approximately 5.1 m.

**Check:** Using kinematics: v² = u² - 2gh → 0 = 100 - 2(9.8)h → h = 100/19.6 ≈ 5.1 m. Matches!

---

### Worked Example 10.5 — Car Braking

**Problem:** A 1000 kg car is traveling at 20 m/s. The driver applies the brakes, which exert a total braking force of 4000 N. How far does the car travel before stopping?

**Work by braking force over distance d:**
  W_brake = -4000 × d    (negative: opposes motion)

**Work-Energy Theorem (final speed = 0):**
  -4000d = 0 - (1/2)(1000)(20)²
  -4000d = -200,000
  d = 200,000 / 4000
  d = 50 m

**Answer:** The car takes 50 m to stop.

**What if the car was going 40 m/s (twice as fast)?**
  -4000d = -(1/2)(1000)(40)²
  -4000d = -800,000
  d = 200 m

Doubling the speed quadruples the stopping distance (from 50 m to 200 m). This is why speed limits are so important for road safety — at high speed, you need much more room to stop.

---

## 11. Why Work and Energy Are Two Faces of the Same Thing

In a deep sense, work and energy are not two separate concepts — they are the same concept viewed from two different angles.

**Energy** is how much work a system CAN do. It is a property stored in the system.
**Work** is how much energy IS transferred from one system to another.

When you push a box and do 200 J of work on it:
- You LOSE 200 J of chemical energy (from your muscles, which came from food).
- The box GAINS 200 J of kinetic energy.

The 200 J of work is the bridge: it is the energy flowing from you to the box.

```
   [YOU]  --200 J of work-->  [BOX]
   
   Your chemical energy -200 J        Box kinetic energy +200 J
   
   Total energy in the system is unchanged. It just moved.
```

This is why work and energy have the same unit (Joule) — because they are the same thing in different states. Work is energy in transit; energy is work in storage.

---

## 12. Preview — Potential Energy and Mechanical Energy

We have studied kinetic energy — the energy of motion. There is another equally important form: **potential energy** — energy stored in a system due to position or configuration.

**Gravitational potential energy** (GPE) is the energy stored when you lift an object against gravity:

  GPE = m × g × h

Where h is the height above some reference level (usually the ground).

When you lift a box to a height h, you do work W = mgh on it. That work is stored as gravitational potential energy. When the box falls, the potential energy converts back to kinetic energy.

**Mechanical energy** is the sum of kinetic energy and potential energy:

  Mechanical Energy = KE + PE = (1/2)mv² + mgh

In the absence of friction and air resistance, **mechanical energy is conserved**: the total stays constant as an object moves. Energy transforms back and forth between KE and PE, but their sum never changes.

We will explore this deeply in the next chapter. For now, appreciate the theme: energy is never destroyed. It only changes form.

```
   Ball thrown upward:
   
   At the bottom (just after throw):
      KE = (1/2)mv²    (maximum)
      PE = 0            (ground level)
      Total = KE + 0
   
   At the top (momentarily at rest):
      KE = 0            (v = 0 at top)
      PE = mgh          (maximum)
      Total = 0 + mgh
   
   Total energy is the same at both points!
   The energy just converted from kinetic to potential.
```

---

## Summary

- **Everyday "work"** means effort. **Physics "work"** means energy transferred by a force through displacement.

- **Work formula:** W = F × d × cos θ. Units: Joules (J). 1 J = 1 N × 1 m.

- **θ is the angle** between the force direction and the displacement direction.
  - θ = 0° (same direction): maximum positive work, W = Fd.
  - θ = 90° (perpendicular): zero work, W = 0.
  - θ = 180° (opposite direction): maximum negative work, W = -Fd.

- **Zero work** is done when: the object does not move (d = 0), or the force is perpendicular to the motion. Examples: carrying a suitcase horizontally, the Moon in a circular orbit.

- **Negative work**: force opposes motion. Examples: friction, air resistance, gravity on a rising object.

- **Energy** is the capacity to do work. Measured in Joules. Forms include kinetic, gravitational potential, chemical, thermal, electromagnetic.

- **Kinetic energy:** KE = (1/2)mv². Depends on v² — doubling speed quadruples KE.

- **Work-Energy Theorem:** W_net = ΔKE = (1/2)mv² - (1/2)mu². The net work done on an object equals the change in its kinetic energy.

- Work and energy are the same thing: work is energy in transit; energy is work in storage. They share the same unit (Joule) for this reason.

- **Conservation of energy:** energy cannot be created or destroyed. When an object "loses" energy, that energy goes somewhere else — often to thermal energy through friction.

---

## Key Equations

```
Work:
   W = F × d × cos θ
   Unit: Joule (J) = Newton × meter = N·m

Special cases:
   W = F × d           (when force is in direction of motion, θ = 0°)
   W = 0               (when force is perpendicular to motion, θ = 90°)
   W = -F × d          (when force opposes motion, θ = 180°)

Work against gravity (lifting):
   W = m × g × h

Work by friction:
   W_friction = -f × d = -μk × m × g × d

Kinetic energy:
   KE = (1/2) × m × v²

Work-Energy Theorem:
   W_net = KE_final - KE_initial
   W_net = (1/2) × m × v² - (1/2) × m × u²

Gravitational potential energy (preview):
   PE = m × g × h

Mechanical energy (preview):
   E = KE + PE = (1/2)mv² + mgh
```
