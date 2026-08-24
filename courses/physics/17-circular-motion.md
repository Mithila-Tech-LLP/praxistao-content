# Chapter 17: Circular Motion

> "The universe is under no obligation to make sense to you." — Neil deGrasse Tyson

---

## Table of Contents

1. [Moving in a Circle — What Makes It Special?](#1-moving-in-a-circle--what-makes-it-special)
2. [Centripetal Acceleration — Always Pointing Inward](#2-centripetal-acceleration--always-pointing-inward)
3. [Centripetal Force — The Net Inward Force](#3-centripetal-force--the-net-inward-force)
4. [The Critical Point — Centripetal Force Is NOT a New Force](#4-the-critical-point--centripetal-force-is-not-a-new-force)
5. [Period, Frequency, and Angular Velocity](#5-period-frequency-and-angular-velocity)
6. [What Provides Centripetal Force in Real Situations?](#6-what-provides-centripetal-force-in-real-situations)
7. [The "Centrifugal Force" Myth — Why You Feel Pushed Outward](#7-the-centrifugal-force-myth--why-you-feel-pushed-outward)
8. [Worked Examples](#8-worked-examples)
9. [Banking of Roads — Curves Without Relying on Friction](#9-banking-of-roads--curves-without-relying-on-friction)
10. [Summary](#summary)
11. [Key Equations](#key-equations)

---

## 1. Moving in a Circle — What Makes It Special?

Imagine swinging a ball on a string above your head in a horizontal circle. The ball moves at a steady pace — maybe 3 meters per second, round and round. If you measured its speed with a speedometer at any instant, you would always get the same number. The speed is constant.

But here is the critical insight: **velocity is NOT the same as speed.**

**Speed** is how fast you are moving — a single number (3 m/s).
**Velocity** is speed PLUS direction — a vector (3 m/s northward, or 3 m/s eastward, etc.).

When an object moves in a circle, its speed stays constant but its **direction changes continuously**. Every fraction of a second, the ball is pointing in a slightly different direction. This means velocity is changing even though speed is not.

And if velocity is changing — that means there IS acceleration (remember: acceleration = change in velocity / time).

This is the surprising truth about circular motion: **an object moving in a circle at constant speed is ALWAYS accelerating.** Not because it is speeding up or slowing down, but because its direction is constantly changing.

```
                 TOP
                  ↓ velocity points LEFT here
                 
    LEFT ←  [ ball ]  → RIGHT
    velocity                velocity
    points                  points
    DOWN                    UP
    
                  ↑ velocity points RIGHT here
                 BOTTOM

    The velocity vector rotates as the ball moves around the circle.
    At every point, velocity is TANGENT to the circle
    (perpendicular to the radius at that point).
```

---

## 2. Centripetal Acceleration — Always Pointing Inward

We established that circular motion involves acceleration. But in which direction?

Consider two adjacent moments of time in the ball's journey around the circle:

```
    Position 1:    Ball at right side of circle
                   Velocity: pointing upward (tangent to circle)
    
    Position 2:    Ball slightly above right side
                   Velocity: pointing slightly left of up (still tangent)

    Change in velocity = velocity_2 - velocity_1
    
    This change in velocity points... TOWARD THE CENTER of the circle!
```

No matter where the ball is on the circle, the change in velocity — and therefore the acceleration — always points **toward the center**.

This inward-pointing acceleration is called **centripetal acceleration** (from Latin: "centrum" = center, "petere" = to seek). It literally means "center-seeking acceleration."

### The Formula

**a_c = v² / r**

Where:
- a_c = centripetal acceleration (m/s²)
- v = speed of the object (m/s)
- r = radius of the circular path (m)

```
              * ← ball here
             /|
            / |
           /  | r (radius)
          /   |
  center *----+
          
    Centripetal acceleration points from the ball TOWARD the center:
    
              * 
             /↙ a_c (points toward center)
            /
           /
  center *
```

### Why Does Faster Speed Mean More Acceleration?

If you swing the ball faster (larger v), it changes direction more rapidly. The acceleration must be larger to keep changing the direction so quickly. And v is squared in the formula, so doubling the speed quadruples the centripetal acceleration.

### Why Does Larger Radius Mean Less Acceleration?

On a very large circle, the path curves gently. The direction changes slowly. So the acceleration needed is smaller. Larger r means smaller a_c.

---

## 3. Centripetal Force — The Net Inward Force

Newton's Second Law says: if there is acceleration, there must be a net force in the same direction as the acceleration.

If centripetal acceleration points toward the center, then the **net force** on the object must also point toward the center.

This inward net force is called **centripetal force**:

**F_c = m × a_c = m × v² / r**

Where:
- F_c = centripetal force (Newtons, N)
- m = mass of the object (kg)
- v = speed (m/s)
- r = radius of circle (m)

### The Direction Always Points INWARD

```
        The centripetal force at various points on a circle:

               ↓ (force points down = toward center)
        
       →  [ball] at top
      [ball]                    The force always
        ←  [ball] at bottom     points inward,
               ↑                toward center
              
        Top of circle:    force points downward (toward center)
        Right of circle:  force points leftward (toward center)
        Bottom of circle: force points upward (toward center)
        Left of circle:   force points rightward (toward center)
```

---

## 4. The Critical Point — Centripetal Force Is NOT a New Force

This is possibly the most important thing to understand in this entire chapter. Read carefully.

**Centripetal force is NOT a new, separate force of nature.** It is not like gravity or friction or tension. It is a LABEL we give to whatever force (or combination of forces) happens to be pointing toward the center and providing the inward acceleration needed for circular motion.

Think of it this way. When you are on a merry-go-round and you label the force keeping you going in a circle as "centripetal force," what IS that force physically? It is the normal force from the seat, or the friction from the floor, or your grip on the handle rail. The label "centripetal force" describes the ROLE of the force (inward, center-seeking), not its NATURE (what type of force it is).

### An Analogy

Imagine you are at a party. Someone asks "Who is the host?" You point to a person: "That person." The word "host" describes their ROLE at the party, not their nature. They are still a human — they happen to be playing the role of "host."

Similarly, "centripetal force" describes the ROLE of some real force. That real force happens to be providing the inward acceleration required for circular motion.

```
  Situation                        What IS the centripetal force?
  ---------------------------------------------------------------
  Ball on a string (horizontal)    Tension in the string
  Planet orbiting the Sun          Gravity from the Sun
  Car turning on a flat road       Friction from the road
  Roller coaster going in a loop   Normal force from the track
  Electron orbiting nucleus        Electrostatic (Coulomb) force
  ---------------------------------------------------------------
  In ALL cases: F_c = mv²/r, but the physical source varies.
```

**Common mistake:** Students sometimes draw an FBD and add a separate "centripetal force" arrow in addition to all the real forces. This is wrong. The centripetal force IS the net real force pointing inward. Do not double-count it.

---

## 5. Period, Frequency, and Angular Velocity

To describe how fast something is going around a circle, we have several related quantities.

### Period (T)

The **period** T is the time it takes to complete one full revolution around the circle. Unit: seconds (s).

If a car goes around a circular track once every 20 seconds, T = 20 s.

### Frequency (f)

The **frequency** f is how many revolutions happen per second. Unit: Hertz (Hz) = revolutions per second.

  f = 1 / T

If T = 20 s, then f = 1/20 = 0.05 Hz (five hundredths of a revolution per second).

If T = 0.5 s (very fast), then f = 2 Hz (two full circles per second).

### Linear Speed (v)

In one full revolution, the object travels a distance equal to the circumference of the circle: C = 2πr.

It takes time T to do this, so:

  v = distance / time = 2πr / T = 2πr × f

### Angular Velocity (ω — pronounced "omega")

Angular velocity describes how fast the angle is changing. Instead of measuring speed in meters per second (how much distance per second), we measure in **radians per second** (how much angle per second).

One full revolution = 360° = 2π radians.

  ω = 2π / T = 2πf    (unit: radians per second, rad/s)

The relationship between linear speed and angular velocity:

  v = ω × r

### Summary Table

```
Quantity          Symbol  Unit          Formula
------------------------------------------------------
Period            T       seconds (s)   T = 1/f = 2π/ω
Frequency         f       Hz (1/s)      f = 1/T = ω/(2π)
Angular velocity  ω       rad/s         ω = 2πf = 2π/T
Linear speed      v       m/s           v = 2πr/T = ωr
Centripetal accel a_c     m/s²          a_c = v²/r = ω²r
Centripetal force F_c     N             F_c = mv²/r = mω²r
```

### Worked Example 5.1 — Spinning Object

**Problem:** A ball moves in a horizontal circle of radius 0.8 m and completes one revolution every 2 seconds. Find: period, frequency, angular velocity, linear speed, and centripetal acceleration.

  Period: T = 2 s

  Frequency: f = 1/T = 1/2 = 0.5 Hz

  Angular velocity: ω = 2π/T = 2π/2 = π ≈ 3.14 rad/s

  Linear speed: v = 2πr/T = 2π × 0.8 / 2 = 2.51 m/s

  Centripetal acceleration: a_c = v²/r = (2.51)² / 0.8 = 6.30 / 0.8 ≈ 7.88 m/s²

---

## 6. What Provides Centripetal Force in Real Situations?

Let us analyze several classic scenarios, drawing FBDs and writing the centripetal force equation for each.

---

### Scenario A — Ball on a String (Horizontal Circle)

A ball of mass m swings in a horizontal circle of radius r on a string at speed v.

```
          [Ball]
         /
        / string (length = r)
       /
   [Hand at center]
   
   View from above (top-down):
   
                T ← (string pulls ball toward center)
         [Ball]
         
   The tension T is the centripetal force.
```

**Setting up the equation:**

The tension T points toward the center. This provides the centripetal force:

  T = mv²/r

Note: In a truly horizontal circle, the string also has to support the ball's weight vertically. So in real life, the string is not perfectly horizontal — it hangs at a slight angle. For simplicity, we often treat these as horizontal first.

---

### Scenario B — Planet Orbiting the Sun

A planet of mass m orbits the Sun at radius r (distance from Sun to planet) and orbital speed v.

```
                   (planet)
                  *
                 /
                / r
               /
          (SUN) *
          
    Gravity from the Sun = centripetal force
    F_gravity = mv²/r
    
    Also: F_gravity = GMm/r² (Newton's Law of Gravitation)
    
    Setting equal: GMm/r² = mv²/r → v = √(GM/r)
    (planets closer to the Sun move faster — consistent with Kepler's Laws)
```

---

### Scenario C — Car Turning on a Flat Road

A car of mass m takes a circular turn of radius r at speed v.

```
   View from above (top-down):
   
       Path of car: curved left
       
       [CAR] →→ direction of motion
         ↑
         |
       Friction from road pointing TOWARD CENTER (left in this case)
       
   FBD of car from the front:
   
       ↑ N (normal force, upward)
       |
      [CAR]
       |
      mg (downward)
   
   Vertical: N = mg (no vertical acceleration)
   
   Horizontal (toward center): f = mv²/r
   
   Maximum friction: f_max = μs × N = μs × mg
   
   So maximum speed: μs × mg = mv²/r
                     v_max² = μs × g × r
                     v_max = √(μs × g × r)
```

If you drive too fast around a corner, friction cannot provide enough centripetal force and the car slides outward (the dreaded "skidding out").

---

### Scenario D — Roller Coaster at the Bottom of a Loop

At the bottom of a vertical loop, the car is moving fast and the track pushes UP on the car. The net upward force provides centripetal acceleration (which points UPWARD toward the center of the loop at this moment).

```
   CENTER OF LOOP
        *
        ↑  (center is above the car at the bottom)
        |
       [CAR] at bottom of loop
        |  ← centripetal direction is UPWARD
        
   FBD of car+person at the bottom:
   
        ↑ N (track pushes up — this is a big force!)
        |
      [CAR]
        |
        ↓ mg (gravity pulls down)
        
   Net upward force = centripetal force:
   N - mg = mv²/r
   N = mg + mv²/r = m(g + v²/r)
```

N > mg. The rider feels heavier than normal at the bottom of a loop. This is why you feel pressed into your seat at the bottom of a roller coaster dip.

---

### Scenario E — Roller Coaster at the TOP of a Loop

At the top of the loop, the center is BELOW the car. Centripetal direction is DOWNWARD.

```
      [CAR] at TOP of loop
        |  ← centripetal direction is DOWNWARD (toward center below)
        ↓
        *
   CENTER OF LOOP
        
   FBD at the top:
   
        ↑ N (track pushes down on car if car is on inside of loop)
        
   Wait — at the top of a loop, the track is ABOVE the car (the car is on the inside
   of the loop). The track can push the car DOWNWARD (normal force downward).
   Gravity also pulls down.
   
   Both N and mg point downward (toward center):
   
      [CAR]
       ↓ N + mg = mv²/r
       
   Therefore: N + mg = mv²/r
              N = mv²/r - mg = m(v²/r - g)
```

For the car to stay on the track at the top, N must be ≥ 0 (the track can push but not pull):

  m(v²/r - g) ≥ 0
  v²/r ≥ g
  v ≥ √(g × r)

So the **minimum speed at the top of a loop** is:

  v_min = √(g × r)

Below this speed, the car would leave the track and fall. Roller coaster designers always ensure the car is going faster than this at the top.

---

## 7. The "Centrifugal Force" Myth — Why You Feel Pushed Outward

You have probably felt it: when a car turns sharply left, you feel pushed to the right — toward the outer edge of the turn. This feels like a force pushing you outward. It is commonly called **centrifugal force** ("centrifugal" means "center-fleeing").

But here is the truth: **centrifugal force is NOT a real force.** It does not appear in a free body diagram drawn from an inertial (non-accelerating) reference frame. No object is actually exerting this outward force on you.

### What Is Actually Happening?

Newton's First Law says: an object in motion tends to stay moving in a straight line. When the car turns left, your body WANTS to keep going straight (forward). The car turns, but you resist turning. The car door — which is now curving with the car — comes around and pushes you INWARD (to the left).

```
   Car turns left:
   
   BEFORE:  Car going straight → [CAR] →
   
   DURING:  Car turns left.
            
            Your body wants to keep going straight.
            The RIGHT-SIDE CAR DOOR swings into you.
            The door pushes you LEFT (inward, toward the center of the turn).
            
            [You] ← force from door  (real, inward force)
```

You push back on the door — that is Newton's Third Law. You feel the door pressing on you. That is the real force.

What you INTERPRET as "being pushed outward" is actually your inertia resisting the inward push. You are not being pushed to the right; the car is pulling you to the left, and your body is lagging behind.

### Why Do Physicists Sometimes Use Centrifugal Force?

In a rotating reference frame (i.e., if you are analyzing things from inside the spinning car), it is mathematically convenient to introduce a "fictitious" centrifugal force to make Newton's Laws work in that frame. Engineers do this all the time. But it is not a real force — it is a mathematical correction for being in an accelerating frame.

### Everyday Examples of "Centrifugal" Sensations

- You lean outward when a bus turns. Your body resists the inward turn.
- Clothes spin to the outside of a washing machine drum. The drum pushes inward; the clothes resist turning and fly outward when the drum stops (or through holes if the drum is perforated).
- A car negotiates a curve — passengers feel pressed against the outer door.

In all cases, the physical explanation is the same: inertia resisting the change in direction.

---

## 8. Worked Examples

### Worked Example 8.1 — Ball on a String

**Problem:** A ball of mass 0.5 kg is attached to a string of length 1.2 m. The ball swings in a horizontal circle at a speed of 3 m/s. Find the tension in the string.

**Setup:**
The tension T provides the centripetal force. Treat the string as horizontal for this calculation.

```
   [Center (hand)]---string (r = 1.2 m)---[Ball, m = 0.5 kg]
   
   T ← provides centripetal force toward center
   
   F_c = mv²/r
   T = mv²/r
```

**Calculation:**
  T = m × v² / r
  T = 0.5 × (3)² / 1.2
  T = 0.5 × 9 / 1.2
  T = 4.5 / 1.2
  T = 3.75 N

**Answer:** The tension in the string is 3.75 N.

**Check:** What if the speed doubles to 6 m/s? T = 0.5 × 36 / 1.2 = 15 N. Tension quadruples when speed doubles (because v is squared). This is why it is hard to swing a ball very fast — the string tension grows rapidly.

---

### Worked Example 8.2 — Car on a Circular Bend

**Problem:** A car takes a circular bend of radius 50 m on a flat, dry road. The coefficient of static friction between the tires and road is μs = 0.4. Find the maximum safe speed.

**Setup:**
On a flat road, friction provides the centripetal force. Maximum friction sets the maximum speed.

```
   View from above:
   
   [CAR] moving in a curve of radius r = 50 m
   
   Friction f points toward center of curve.
   
   Maximum friction: f_max = μs × N = μs × mg (since N = mg on flat road)
   
   At maximum speed: f_max = mv²/r
   μs × mg = mv²/r
   μs × g = v²/r
   v² = μs × g × r
   v_max = √(μs × g × r)
```

**Calculation:**
  v_max = √(μs × g × r)
  v_max = √(0.4 × 9.8 × 50)
  v_max = √(196)
  v_max = 14 m/s

**Convert to km/h:** 14 × 3.6 = 50.4 km/h

**Answer:** The maximum safe speed is 14 m/s (about 50 km/h).

**What happens if you exceed this speed?** Friction cannot provide enough centripetal force. The car's inertia carries it outward — it slides wide of the turn. This is "skidding out."

**Check:** What if the road is icy (μs = 0.05)?
  v_max = √(0.05 × 9.8 × 50) = √24.5 ≈ 4.95 m/s ≈ 18 km/h.
  Drastically slower! This is why icy roads are so dangerous on curves.

---

### Worked Example 8.3 — Roller Coaster Loop

**Problem:** A roller coaster loop has a radius of 8 m. What is the minimum speed the car must have at the top of the loop to stay on the track?

**Setup:**
At the top of the loop, the car is on the inside of the track. Both gravity mg and the normal force N from the track point downward (toward the center of the loop). For minimum speed, the normal force = 0 (car barely maintains contact).

```
   TOP OF LOOP:
   
      [CAR]  ← on the inside of the loop, track is above car
       ↓ N (track pushes down)
       ↓ mg (gravity)
       
   Center of loop is BELOW:  both forces point toward center.
   
   N + mg = mv²/r
   
   At minimum speed: N = 0
   mg = mv²_min / r
   g = v²_min / r
   v_min = √(g × r)
```

**Calculation:**
  v_min = √(g × r)
  v_min = √(9.8 × 8)
  v_min = √78.4
  v_min ≈ 8.85 m/s

**Answer:** The minimum speed at the top of the loop is approximately 8.85 m/s (about 31.8 km/h).

If the car is going slower than 8.85 m/s at the top, gravity exceeds what is needed for circular motion at that radius, and the car falls away from the track — not a good day for the passengers.

---

## 9. Banking of Roads — Curves Without Relying on Friction

### The Problem With Flat Roads

On a flat road, the only force available to push a turning car toward the center of the curve is friction. In wet or icy conditions, friction is reduced, and the maximum safe speed drops dangerously.

Highway engineers solve this problem by **banking** the road — tilting it inward toward the center of the curve. A banked road allows cars to turn at higher speeds, partly or entirely without needing friction.

### How Banking Works

```
   Cross-section of a banked road (looking along the road):
   
   Outer edge                        Inner edge
   (higher)                          (lower)
        \                           /
         \                         /
          \                       /
           \       [CAR]         /
            \       ↑ N         /
             \      |          /
              \     |  angle θ/
               \    |        /
                \   |       /
                 \  |      /
                  \ |     /
                   \|θ___/
                    center of curve is to the right
                    
   N is perpendicular to the banked surface.
   N has two components:
      Vertical:   N cos θ  (this must balance gravity: N cos θ = mg)
      Horizontal: N sin θ  (this points toward center — it IS the centripetal force!)
```

**Vertical equilibrium (no vertical acceleration):**
  N cos θ = mg
  N = mg / cos θ

**Horizontal centripetal force:**
  N sin θ = mv²/r

**Divide horizontal by vertical:**
  (N sin θ) / (N cos θ) = (mv²/r) / mg
  tan θ = v² / (rg)

**Ideal banking angle:**
  tan θ = v² / (r × g)

At this angle, the car takes the curve with NO friction needed. The normal force alone provides all the centripetal force.

If the car goes faster than the ideal speed, it needs friction pointing inward (down the bank). If it goes slower, it needs friction pointing outward (up the bank).

---

### Worked Example 9.1 — Banking Angle

**Problem:** A highway curve has a radius of 200 m. At what angle should the road be banked for a design speed of 25 m/s?

  tan θ = v² / (r × g)
  tan θ = (25)² / (200 × 9.8)
  tan θ = 625 / 1960
  tan θ = 0.319
  θ = arctan(0.319) ≈ 17.7°

**Answer:** The road should be banked at approximately 17.7° for vehicles traveling at 25 m/s (90 km/h).

---

### Why Banked Roads Are Safer

On a flat road with friction: v_max = √(μsgr). This depends heavily on μs, which varies with weather.

On a properly banked road: v_ideal = √(rg tan θ). This is fixed by the road geometry — it does not depend on friction at all. For small deviations from v_ideal, friction provides a safety margin.

That is why banked tracks in Formula 1 racing and NASCAR allow much higher speeds through turns than flat roads would allow.

---

## Summary

- An object moving in a circle at constant speed has a **changing velocity** (direction changes), so there IS acceleration.

- **Centripetal acceleration** always points toward the center of the circle: a_c = v²/r.

- **Centripetal force** is the net inward force required to maintain circular motion: F_c = mv²/r.

- Centripetal force is NOT a new type of force. It is always provided by a real physical force: tension, gravity, friction, or normal force depending on the situation.

- **Period T** = time for one revolution. **Frequency f** = 1/T. **Angular velocity** ω = 2πf. **Linear speed** v = ωr = 2πr/T.

- **Ball on a string:** T = mv²/r (tension provides centripetal force).

- **Car on flat road:** Maximum safe speed v_max = √(μsgr). Friction provides centripetal force.

- **Roller coaster bottom:** N = m(g + v²/r) — you feel heavier.

- **Roller coaster top:** N = m(v²/r - g) — minimum speed v_min = √(gr) to stay on track.

- **"Centrifugal force"** is not real — it is your inertia resisting the change in direction. The real force is always inward (from the car door, track, rope, etc.).

- **Banking of roads:** The ideal bank angle satisfies tan θ = v²/(rg). At this angle, no friction is needed to navigate the curve.

---

## Key Equations

```
Centripetal acceleration:
   a_c = v² / r = ω² × r

Centripetal force:
   F_c = m × v² / r = m × ω² × r

Period and frequency:
   f = 1 / T
   ω = 2π × f = 2π / T
   v = 2πr / T = ω × r

Ball on a string (horizontal circle):
   T = m × v² / r

Car on flat road (maximum safe speed):
   v_max = √(μs × g × r)

Roller coaster at bottom of loop:
   N = m × (g + v²/r)

Roller coaster at top of loop:
   N = m × (v²/r - g)
   Minimum speed: v_min = √(g × r)

Ideal banking angle:
   tan θ = v² / (r × g)
```
