# Chapter 22: Momentum and Impulse

> **"Momentum is the 'quantity of motion' in an object. It's why a slow-moving truck is harder to stop than a fast-moving bicycle — and why rockets work in space."**

---

## Table of Contents

- [22.1 What is Momentum?](#221-what-is-momentum)
- [22.2 Newton's Second Law in Terms of Momentum](#222-newtons-second-law-in-terms-of-momentum)
- [22.3 Impulse](#223-impulse)
- [22.4 The Impulse-Momentum Theorem](#224-the-impulse-momentum-theorem)
- [22.5 Conservation of Momentum](#225-conservation-of-momentum)
- [22.6 Elastic and Inelastic Collisions](#226-elastic-and-inelastic-collisions)
- [22.7 Explosions](#227-explosions)
- [22.8 2D Momentum Problems](#228-2d-momentum-problems)
- [22.9 Real-World Applications](#229-real-world-applications)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 22.1 What is Momentum?

**Momentum** is the product of an object's mass and velocity.

```
p = m × v

where:
  p = momentum (kg·m/s)
  m = mass (kg)
  v = velocity (m/s)
```

Momentum is a **vector** — it has both magnitude and direction.

### Why Momentum Matters

Compare these two objects:

```
TRUCK                          BICYCLE
Mass = 10,000 kg               Mass = 10 kg
Speed = 5 m/s                  Speed = 10 m/s

p = 10000 × 5                  p = 10 × 10
p = 50,000 kg·m/s              p = 100 kg·m/s

The truck has 500x more momentum!
```

The truck is much harder to stop — it has much more momentum.

### Visualizing Momentum

```
SMALL FAST OBJECT              LARGE SLOW OBJECT
   v = 30 m/s                     v = 3 m/s
   ---->                          ------>
  [  ] m = 1 kg                 [       ] m = 10 kg
  p = 30 kg·m/s                 p = 30 kg·m/s

Equal momenta — both equally difficult to stop!
```

---

## 22.2 Newton's Second Law in Terms of Momentum

The original form of Newton's Second Law (as Newton himself wrote it):

```
         dp       change in momentum
F_net = ---- = ---------------------------
         dt        change in time
```

This is more general than F = ma. When mass is constant:

```
       m × change_v
F = ------------------- = m × a
        change_t
```

So F = ma is a special case of the more general momentum form.

The momentum form applies even when mass changes (like a rocket burning fuel).

---

## 22.3 Impulse

**Impulse** is the product of a force and the time over which it acts.

```
J = F × t

where:
  J = impulse (N·s or kg·m/s)
  F = average force (N)
  t = time of contact (s)
```

Impulse is also a vector.

### Impulse on a Force-Time Graph

The impulse is the **area under a Force-Time graph**:

```
Force (N)
 |
F|  *---------*
 | /           \
 |/             \
 +-------------------> time (s)
    t1           t2

Impulse = area under the curve
        ≈ F_avg × (t2 - t1)
```

For a constant force: Impulse = F × Δt (a rectangle).
For a varying force: Impulse = area (which may be a triangle, trapezoid, or irregular shape).

### Why Impulse Matters: Airbags and Safety

```
CRASH WITHOUT AIRBAG:
  Force acts for very short time (0.01 s)
  Same momentum change, shorter time
  Force = delta_p / delta_t = HUGE force on body
  
CRASH WITH AIRBAG:
  Force acts over longer time (0.1 s)
  Same momentum change, longer time
  Force = delta_p / delta_t = 10x SMALLER force on body
  
Same impulse, longer time, smaller force = safer!
```

---

## 22.4 The Impulse-Momentum Theorem

The **Impulse-Momentum Theorem** states:

```
Impulse = Change in momentum

J = delta_p = m × v_final - m × v_initial

F × delta_t = m × delta_v
```

This connects force, time, mass, and velocity in one powerful equation.

### Worked Example 22.1

A cricket ball of mass 0.16 kg is bowled at 30 m/s. The batsman hits it back at 40 m/s.

The contact time is 0.002 s. Find the average force on the ball.

```
BEFORE                    AFTER
   --->  30 m/s              40 m/s  <---
  [ball]                    [ball]
  m = 0.16 kg               m = 0.16 kg
```

**Solution:**

Taking "towards batsman" as positive:
- v_initial = +30 m/s (ball coming in)
- v_final = -40 m/s (ball going back)

Change in momentum = m(v_f - v_i) = 0.16 × (-40 - 30) = 0.16 × (-70) = -11.2 kg·m/s

Impulse = change in momentum = -11.2 N·s

Average force = impulse / time = -11.2 / 0.002 = **-5600 N** (towards batsman)

Magnitude of force = **5600 N** (about 570 kg weight!)

---

## 22.5 Conservation of Momentum

**The Law of Conservation of Momentum:** In a closed system (no external forces), the total momentum is constant.

```
Before collision:    total p = p1 + p2

After collision:     total p = p1' + p2'

Therefore:           p1 + p2 = p1' + p2'

m1×v1 + m2×v2 = m1×v1' + m2×v2'
```

### Why is Momentum Conserved?

By Newton's Third Law, when objects A and B collide:
- Force on B from A = -Force on A from B
- Both forces act for the same time
- So Impulse on B = -Impulse on A
- Momentum gained by B = momentum lost by A
- Total momentum unchanged

```
COLLISION FORCES:
   [A] ---FAB---> [B]
   [A] <---FBA--- [B]

FAB = -FBA (Newton's 3rd Law)

Both act for same time dt:
Impulse on B = FAB × dt = delta_pB
Impulse on A = FBA × dt = -FAB × dt = -delta_pB

So: delta_pA + delta_pB = 0
Total momentum conserved!
```

---

## 22.6 Elastic and Inelastic Collisions

### Types of Collision

```
ELASTIC COLLISION:
  - Momentum conserved
  - Kinetic energy conserved
  - Objects bounce off each other perfectly
  - Rare in practice (billiard balls approximately elastic)

INELASTIC COLLISION:
  - Momentum conserved
  - Kinetic energy NOT conserved (some becomes heat, sound, deformation)
  - Objects may deform or stick together

PERFECTLY INELASTIC COLLISION:
  - Momentum conserved
  - Maximum kinetic energy lost
  - Objects stick together and move as one
```

### Elastic Collision Example

```
BEFORE:
  [A: 2kg, 6m/s] -->         <-- [B: 3kg, 4m/s]

AFTER:
  [A: 2kg, ?] <--         --> [B: 3kg, ?]
```

Setting up conservation equations:
- Conservation of momentum: m_A × v_A + m_B × v_B = m_A × v_A' + m_B × v_B'
- Conservation of KE: 1/2 × m_A × v_A² + 1/2 × m_B × v_B² = 1/2 × m_A × v_A'² + 1/2 × m_B × v_B'²

For elastic collisions between equal masses: objects exchange velocities!

### Worked Example 22.2 — Perfectly Inelastic Collision

A 1500 kg car moving at 20 m/s east hits a stationary 1000 kg truck. They stick together.

Find their combined velocity after the collision.

**Solution:**

```
BEFORE:
  [car: 1500kg, 20m/s] --> [truck: 1000kg, 0m/s]

AFTER:
  [combined: 2500kg, v?] -->
```

Conservation of momentum:
- Before: p = 1500 × 20 + 1000 × 0 = 30,000 kg·m/s
- After: p = 2500 × v
- 2500v = 30,000
- v = **12 m/s east**

### Checking Kinetic Energy

- KE before = 1/2 × 1500 × 20² = 300,000 J
- KE after = 1/2 × 2500 × 12² = 180,000 J
- KE lost = 120,000 J (became heat, sound, deformation)

---

## 22.7 Explosions

An explosion is the reverse of a perfectly inelastic collision. Objects start together and fly apart.

### Key Point: Total momentum before = total momentum after

If the system starts at rest: total momentum = 0, so the momenta of the pieces must sum to zero.

```
EXPLOSION FROM REST:
  [OBJECT at rest]
  
       |
       v  BANG!
       
  [piece1]         [piece2]
  <-- v1            v2 -->
  
  m1 × v1 = m2 × v2 (magnitudes)
  
  Momenta equal and opposite!
```

### Worked Example 22.3 — Rocket

A rocket in space (total mass 5000 kg) fires its engine. In 1 second, it ejects 50 kg of exhaust gas backward at 400 m/s.

Find the change in velocity of the rocket.

**Solution:**

Initial momentum = 0 (at rest)

Final momentum:
- Exhaust: 50 × (-400) = -20,000 kg·m/s
- Rocket: remaining mass = 4950 kg, velocity = v

Conservation: 0 = 4950v + (-20,000)
4950v = 20,000
v = **+4.04 m/s** (forward)

This is how rockets work in space — no need for air to push against!

---

## 22.8 2D Momentum Problems

Momentum is a vector. In 2D problems, you must apply conservation in both x and y directions separately.

```
BEFORE:
  Ball A: mass m1, velocity v1 at angle 0 (horizontal)
  Ball B: mass m2, at rest

  v1 -->
  [A]     [B] stationary

AFTER:
  Ball A moves at angle θ1 above horizontal
  Ball B moves at angle θ2 below horizontal

           /  [A] at angle θ1
  [A] -->
           \  [B] at angle θ2
```

Setting up equations:
- x-direction: m1×v1 = m1×v1'×cos(θ1) + m2×v2'×cos(θ2)
- y-direction: 0 = m1×v1'×sin(θ1) - m2×v2'×sin(θ2)

Solve the two simultaneous equations for the unknowns.

---

## 22.9 Real-World Applications

### Car Safety Systems

```
CRUMPLE ZONES:
  Purpose: increase collision time
  Effect: same impulse over longer time = smaller force
  
SEATBELTS:
  Purpose: apply force gradually over longer time
  Effect: reduce peak force on chest/neck
  
AIRBAGS:
  Purpose: increase time for head to stop
  Effect: dramatically reduce peak force on skull
```

### Ball Sports

In cricket, tennis, and golf, coaches teach players to "follow through". This extends the contact time and increases impulse (change in momentum), giving the ball more speed.

### Jet Propulsion

Jet engines and rockets work by the conservation of momentum:
- Exhaust gases expelled backward at high speed
- Equal and opposite momentum given to aircraft/rocket forward

### Ballistic Pendulum

A classic physics experiment to measure bullet speed:

```
BULLET FIRED
into hanging block:

   v (fast)
   .---> [block at rest]
         |
         | (string)
         |
         ^
         |
    [block swings up
     to height h]

Conservation of momentum:
  m_bullet × v = (m_bullet + m_block) × V

Conservation of energy (swing):
  1/2(M)V^2 = M×g×h

Combining: v = (M/m_bullet) × sqrt(2gh)
```

---

## Summary

- **Momentum**: p = mv (vector, kg·m/s)
- Newton's 2nd Law: F = dp/dt (rate of change of momentum)
- **Impulse**: J = F × Δt = area under F-t graph = change in momentum
- **Impulse-Momentum Theorem**: F × Δt = m × Δv
- **Conservation of Momentum**: in a closed system, total p before = total p after
- **Elastic collision**: both momentum AND kinetic energy conserved
- **Inelastic collision**: momentum conserved, kinetic energy NOT conserved
- **Perfectly inelastic**: objects stick together, maximum KE lost
- **Explosions**: start from rest (p=0), pieces fly apart with equal and opposite momenta
- Safety devices (airbags, crumple zones) increase collision time to reduce peak force

---

## Key Equations

```
Momentum:
  p = m × v

Newton's 2nd Law:
  F = dp/dt = m × a  (when mass is constant)

Impulse:
  J = F × delta_t  (for constant force)
  J = area under F-t graph  (for varying force)

Impulse-Momentum Theorem:
  J = delta_p = m × (v_f - v_i)
  F × delta_t = m × delta_v

Conservation of Momentum (1D):
  m1×v1 + m2×v2 = m1×v1' + m2×v2'

Elastic collision (equal masses):
  Objects exchange velocities

Kinetic Energy:
  KE = (1/2) × m × v^2
```
