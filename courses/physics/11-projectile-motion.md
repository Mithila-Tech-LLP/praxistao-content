# Chapter 11: Projectile Motion

> **"A cannonball, a basketball, and a raindrop — they all obey the same equations once they're in the air. Gravity is the great equalizer."**

---

## Table of Contents

- [What Is a Projectile?](#what-is-a-projectile)
- [The Big Insight — Two Independent Motions](#the-big-insight--two-independent-motions)
- [Setting Up the Coordinate System](#setting-up-the-coordinate-system)
- [Projectile Launched Horizontally](#projectile-launched-horizontally)
- [Worked Example 1 — Ball Rolled Off a Table](#worked-example-1--ball-rolled-off-a-table)
- [Projectile Launched at an Angle](#projectile-launched-at-an-angle)
- [Breaking the Launch into Components](#breaking-the-launch-into-components)
- [Key Equations for Angled Projectiles](#key-equations-for-angled-projectiles)
- [Why 45 Degrees Gives Maximum Range](#why-45-degrees-gives-maximum-range)
- [Worked Example 2 — Ball Kicked at an Angle](#worked-example-2--ball-kicked-at-an-angle)
- [Worked Example 3 — Clearing a Crossbar](#worked-example-3--clearing-a-crossbar)
- [Symmetry of Projectile Motion](#symmetry-of-projectile-motion)
- [Real-World Projectiles](#real-world-projectiles)
- [Effect of Air Resistance](#effect-of-air-resistance)
- [The Full Parabolic Path](#the-full-parabolic-path)
- [Common Mistakes to Avoid](#common-mistakes-to-avoid)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## What Is a Projectile?

A **projectile** is any object that has been launched into the air and then moves only under the influence of gravity, with no engine or thrust pushing it after the initial launch.

Examples:
- A football kicked through the air
- A stone thrown from a cliff
- A cannonball fired from a cannon
- A basketball in flight between your hands and the hoop
- Water shooting from a garden hose
- A long jumper flying through the air
- A bullet after it leaves the barrel

What is NOT a projectile?
- A rocket with its engine burning (it has thrust)
- An airplane (it has lift and engine force)
- A ball sitting still on a shelf (it is not in flight)

The key defining feature: once it is launched, only gravity acts on it (assuming we ignore air resistance, which is a common and useful simplification for introductory physics).

---

## The Big Insight — Two Independent Motions

The secret to understanding projectile motion is the same independence principle from the last chapter:

> **Horizontal and vertical motions are completely independent of each other.**

Let us spell out what this means for a projectile:

### Horizontal Motion: Constant Velocity

Once the projectile is in the air, there is no horizontal force acting on it (gravity only pulls downward). With no horizontal force, there is no horizontal acceleration. So the horizontal velocity stays constant throughout the flight.

    Horizontal: constant velocity, no acceleration
    x = vₓ × t

### Vertical Motion: Constant Downward Acceleration

Gravity pulls the projectile downward with an acceleration of g = 9.8 m/s² (we often round to 10 m/s² for easier arithmetic). This is exactly the same as free fall, which you studied in Chapter 9.

    Vertical: accelerating downward at g = 9.8 m/s²
    y = vᵧ₀ × t - ½ × g × t²

(The minus sign is because gravity pulls downward, which is the negative y direction.)

### The Parabola

Because x grows linearly with time (x = vₓ × t) and y grows with t² (from the vertical acceleration), the path of a projectile is a **parabola** — that classic curved shape you see every time you throw something.

```
Height (y)
    |
    |          *
    |       *     *
    |     *         *
    |    *             *
    |   *                 *
    |  *                     *
    | *                         *
    |*                               *
    +-------------------------------------------  Horizontal (x)
    
    The parabolic path of a projectile.
    
    - Left half: projectile goes up AND forward
    - Right half: projectile goes down AND forward
    - Horizontal spacing: EQUAL intervals (constant horizontal speed)
    - Vertical spacing: UNEQUAL (changes due to gravity)
```

---

## Setting Up the Coordinate System

We use the standard convention:
- **Origin**: the launch point of the projectile
- **x-axis**: horizontal, pointing in the direction of launch
- **y-axis**: vertical, pointing upward
- Positive x: forward (in the launch direction)
- Positive y: upward
- Negative y: downward
- **g = 9.8 m/s²** is the magnitude of gravitational acceleration (always positive)

When we write the vertical acceleration in our equations, we write **-g** (negative because gravity acts downward, in the -y direction).

This setup means:
- At t = 0: position is (0, 0), at the launch point
- Horizontal velocity: vₓ = constant throughout flight
- Vertical velocity: changes every instant due to gravity

---

## Projectile Launched Horizontally

The simplest case: you throw an object horizontally from a height h. It has an initial horizontal velocity u and zero initial vertical velocity.

This is like rolling a ball off the edge of a table.

### The Equations

**Horizontal** (no acceleration, constant velocity u):

    x = u × t

**Vertical** (starts from rest vertically, accelerates downward at g):

    y = h - ½ × g × t²

(y is the height above the ground; it starts at h and decreases)

### Time to Hit the Ground

The projectile hits the ground when y = 0:

    0 = h - ½ × g × t²
    ½ × g × t² = h
    t² = 2h / g
    t = sqrt(2h / g)

### Horizontal Range

The horizontal distance travelled before hitting the ground:

    R = u × t = u × sqrt(2h / g)

### Observations

- A larger h (higher starting point) means more time in the air, so more range.
- A larger u (faster horizontal launch) means more range for the same h.
- g is fixed at 9.8 m/s² — you cannot change it on Earth.
- The time of flight depends only on h, not on u. Two balls rolled off the same height hit the ground at the same time, no matter their horizontal speed.

```
Table (height = h)
+----------------------*  ← Ball leaves edge here with speed u

                         *  ← After a moment
                         
                              *  ← Halfway down
                              
                                    *  ← Just before landing
                                    
Ground:--------------------+--------+--------+-------
                         (drops)    Range R = u × t
```

---

## Worked Example 1 — Ball Rolled Off a Table

**Problem:** A ball rolls off the edge of a table that is 1.2 m high. It leaves the edge horizontally at 3 m/s. Find:

(a) How long it takes to hit the floor
(b) How far from the table base it lands
(c) Its speed just before it hits the floor

### Given

- Height: h = 1.2 m
- Initial horizontal velocity: u = 3 m/s
- Initial vertical velocity: vᵧ₀ = 0 (launched horizontally)
- g = 9.8 m/s²

### Part (a): Time of Flight

    t = sqrt(2h / g)
      = sqrt(2 × 1.2 / 9.8)
      = sqrt(2.4 / 9.8)
      = sqrt(0.2449)
      ≈ 0.495 seconds

So the ball takes about 0.5 seconds to fall 1.2 m.

### Part (b): Horizontal Range

    R = u × t
      = 3 × 0.495
      ≈ 1.48 m

The ball lands about 1.48 m from the base of the table.

### Part (c): Speed Just Before Landing

When the ball lands, it has two components of velocity:

Horizontal component (unchanged throughout flight):
    vₓ = 3 m/s

Vertical component (gained by falling under gravity):
    vᵧ = g × t = 9.8 × 0.495 ≈ 4.85 m/s (downward)

Total speed on impact:
    v = sqrt(vₓ² + vᵧ²)
      = sqrt(3² + 4.85²)
      = sqrt(9 + 23.52)
      = sqrt(32.52)
      ≈ 5.70 m/s

Direction (angle below horizontal):
    tan(θ) = vᵧ / vₓ = 4.85 / 3 = 1.617
    θ = arctan(1.617) ≈ 58.3° below horizontal

```
    Table edge (height = 1.2 m)
    * ← launches at 3 m/s horizontal
     \
      \
       \   (parabolic path)
        \
         \
          * ← lands 1.48 m away
          
    Impact velocity: 5.70 m/s at 58.3° below horizontal
    (mostly downward by the time it lands)
```

---

## Projectile Launched at an Angle

Most real projectiles are launched at an angle — not straight up and not horizontally. A kick, a throw, a cannonball — they go up and forward at the same time.

This case is more interesting and requires a bit more setup, but the core method is the same: separate into x and y components.

---

## Breaking the Launch into Components

If a projectile is launched with speed u at angle θ above the horizontal:

**Initial horizontal velocity:**

    vₓ₀ = u × cos(θ)

**Initial vertical velocity:**

    vᵧ₀ = u × sin(θ)

These two components then behave independently throughout the flight.

```
         ^
         | vᵧ₀ = u sin θ
         |
         |   / ← launch velocity u
         |  /
         | /
         |/ θ
         +------------>
              vₓ₀ = u cos θ
              
    The launch velocity vector u has two components.
    vₓ₀ is the horizontal component (constant throughout flight).
    vᵧ₀ is the vertical component (decreases due to gravity).
```

### How the Velocity Changes During Flight

**Horizontal velocity:** stays constant at vₓ₀ = u cos(θ) throughout.

**Vertical velocity:** starts at vᵧ₀ = u sin(θ) and decreases at the rate of g m/s per second.

    At time t: vᵧ = u sin(θ) - g × t

At the top of the arc: vᵧ = 0 (the projectile is momentarily moving only horizontally).

After the top: vᵧ becomes negative (moving downward).

```
         At launch:        At peak:         On landing:
         
         ^                 →                \
         |\                                  \
         | \                                  v
         |  vₓ →           vₓ →               vₓ →
         |
         vᵧ upward         vᵧ = 0             vᵧ downward
```

---

## Key Equations for Angled Projectiles

### Position at Time t

    x(t) = vₓ₀ × t = u cos(θ) × t
    y(t) = vᵧ₀ × t - ½ × g × t² = u sin(θ) × t - ½ × g × t²

### Velocity at Time t

    vₓ(t) = u cos(θ)           ← constant
    vᵧ(t) = u sin(θ) - g × t  ← changes with time

### Time to Reach Maximum Height

At maximum height, vertical velocity = 0:

    0 = u sin(θ) - g × t_top
    t_top = u sin(θ) / g

### Maximum Height

Substitute t_top into the y equation:

    H = u sin(θ) × t_top - ½ × g × t_top²
    H = u sin(θ) × (u sin θ / g) - ½ × g × (u sin θ / g)²
    H = (u sin θ)² / g - ½ × (u sin θ)² / g
    H = ½ × (u sin θ)² / g
    H = (u sin θ)² / (2g)

### Total Time of Flight

Since the motion is symmetric (going up takes the same time as coming down), the total time is twice the time to the top:

    T = 2 × t_top = 2u sin(θ) / g

### Horizontal Range

The total horizontal distance travelled:

    R = vₓ₀ × T = u cos(θ) × 2u sin(θ) / g
    R = 2u² sin(θ) cos(θ) / g

Using the trigonometric identity: 2 sin(θ) cos(θ) = sin(2θ):

    R = u² sin(2θ) / g

This is the famous **range formula** for projectile motion.

---

## Why 45 Degrees Gives Maximum Range

The range formula is:

    R = u² sin(2θ) / g

For fixed launch speed u and fixed g, the range is largest when sin(2θ) is largest.

The sine function reaches its maximum value of 1 when its argument equals 90°.

    sin(2θ) = 1
    2θ = 90°
    θ = 45°

So the range is maximum when the launch angle is **45 degrees**.

```
Range R
  |                    * max range at 45°
  |                  *   *
  |                *       *
  |              *           *
  |            *               *
  |          *                   *
  |        *                       *
  |      *                           *
  +----+---+---+---+---+---+---+---+---+---> θ (degrees)
       10  20  30  40  50  60  70  80  90
       
    Range increases up to 45°, then decreases symmetrically.
    Notice: 30° and 60° give the same range (sin 60° = sin 120°).
    So do 20° and 70°, 15° and 75°, etc.
```

### Complementary Angles Give Equal Range

Another beautiful result: launch angles that add up to 90° give the same range.

- 30° and 60° give the same range
- 20° and 70° give the same range
- 10° and 80° give the same range

Why? Because sin(2θ) = sin(180° - 2θ). If θ₁ + θ₂ = 90°, then 2θ₁ + 2θ₂ = 180°, and sin(2θ₁) = sin(2θ₂).

The 30° and 60° trajectories have the same range but different shapes:
- 30° is lower and faster (less height, more horizontal)
- 60° is higher and slower (more height, less horizontal)

```
       Height
         |
         |             60° path (tall arc)
         |          ___---___
         |       __/         \__
         |      /               \
         |  ___/  30° path        \___
         | /  (lower, wider arc)      \
         |/                            \
         +---------------------------------  horizontal
              <-------- same range ------->
```

---

## Worked Example 2 — Ball Kicked at an Angle

**Problem:** A ball is kicked with an initial speed of 20 m/s at an angle of 30° above the horizontal. Find:

(a) The maximum height reached
(b) The total time of flight
(c) The horizontal range
(d) The speed at the highest point

### Given

- Launch speed: u = 20 m/s
- Launch angle: θ = 30°
- g = 9.8 m/s²

### Initial Velocity Components

    vₓ₀ = u cos(30°) = 20 × 0.866 = 17.32 m/s
    vᵧ₀ = u sin(30°) = 20 × 0.5 = 10 m/s

### Part (a): Maximum Height

    H = (u sin θ)² / (2g)
      = (20 × 0.5)² / (2 × 9.8)
      = (10)² / 19.6
      = 100 / 19.6
      ≈ 5.10 m

The ball reaches a maximum height of about 5.10 m.

### Part (b): Total Time of Flight

    T = 2u sin(θ) / g
      = 2 × 20 × 0.5 / 9.8
      = 20 / 9.8
      ≈ 2.04 seconds

The ball is in the air for about 2.04 seconds.

### Part (c): Horizontal Range

    R = u² sin(2θ) / g
      = 20² × sin(60°) / 9.8
      = 400 × 0.866 / 9.8
      = 346.4 / 9.8
      ≈ 35.35 m

The ball lands about 35.35 m away.

### Part (d): Speed at the Highest Point

At the highest point, the vertical velocity is zero. Only horizontal velocity remains:

    v = vₓ₀ = 17.32 m/s

The speed at the highest point equals the initial horizontal velocity component.

```
Height (m)
    |
  5 |          *  ← maximum height: 5.10 m
    |        *   *
    |       *     *
  3 |      *       *
    |     *         *
    |    *           *
  1 |   *             *
    |  *               *
    | *                 *
    |*                   *
    +--+---+---+---+-----+----> horizontal (m)
       5   10  15  20    35
    
    Launch: 30°, 20 m/s         Lands: 35.35 m away
    T = 2.04 s
```

---

## Worked Example 3 — Clearing a Crossbar

**Problem:** A footballer kicks a ball from the ground at a speed of 18 m/s at an angle of 40°. A goalpost crossbar is 3.0 m high and 28 m from the kick. Does the ball clear the crossbar?

### Given

- Launch speed: u = 18 m/s
- Launch angle: θ = 40°
- Crossbar distance: x = 28 m
- Crossbar height: 3.0 m
- g = 9.8 m/s²

### Step 1: Find Initial Components

    vₓ₀ = 18 × cos(40°) = 18 × 0.766 = 13.79 m/s
    vᵧ₀ = 18 × sin(40°) = 18 × 0.643 = 11.57 m/s

### Step 2: Find the Time When Ball Reaches x = 28 m

    x = vₓ₀ × t
    28 = 13.79 × t
    t = 28 / 13.79 ≈ 2.03 seconds

### Step 3: Find the Ball's Height at That Time

    y = vᵧ₀ × t - ½ × g × t²
      = 11.57 × 2.03 - ½ × 9.8 × (2.03)²
      = 23.49 - ½ × 9.8 × 4.12
      = 23.49 - 20.19
      = 3.30 m

### Step 4: Does It Clear?

The ball is at 3.30 m when it reaches the crossbar location. The crossbar is at 3.0 m.

Since 3.30 m > 3.0 m, **yes, the ball clears the crossbar** — but only just, by 0.30 m.

### A Useful Check: Is the Ball Still Going Up or Coming Down?

Total time of flight:
    T = 2 × 11.57 / 9.8 = 2.36 seconds

The ball arrives at the crossbar at t = 2.03 s, which is before the midpoint (T/2 = 1.18 s). Wait — 2.03 > 1.18, so the ball has already passed its peak and is coming DOWN when it reaches the crossbar.

Time to peak: t_top = 11.57 / 9.8 = 1.18 s. Since the ball arrives at x = 28 m at t = 2.03 s > 1.18 s, the ball is on the descending part of its arc. It just barely clears.

```
Height (m)
    |
  6 |          *
    |        *   *
    |      *       *
  3 |    *     ----+----  ← crossbar at 3.0 m, x = 28 m
    |   *       ball is at 3.30 m here (just clears)
    |  *               *
    | *                 *
    |*                   *
    +--+---+---+---+---+-+----> horizontal (m)
       5   10  15  20  28
    
    Ball clears crossbar by 0.30 m
```

---

## Symmetry of Projectile Motion

One of the most elegant properties of projectile motion (ignoring air resistance) is its perfect symmetry.

### The Arc Is a Perfect Parabola

The path is symmetric around the highest point. The left half is a mirror image of the right half.

### Speed on Landing Equals Speed on Launch

When the projectile returns to the same height it was launched from, it has exactly the same **speed** as when it launched. The direction is different (it is going downward instead of upward) but the magnitude is identical.

Why? Because the kinetic energy at launch and at landing must be equal (if the height is the same, the gravitational potential energy is the same, so kinetic energy is the same, so speed is the same). We will revisit this properly in the Energy chapters.

### Velocity Reversal

At launch: vᵧ₀ = u sin(θ) (upward)
At landing: vᵧ = -u sin(θ) (downward, same magnitude)
The horizontal component never changes: vₓ = u cos(θ) throughout.

### Going Up Takes as Long as Coming Down

The time to reach the top equals the time to fall back down:

    t_up = t_down = u sin(θ) / g

This is not obvious, but it follows directly from the symmetric shape of the parabola.

```
Going up and coming down are mirror images:

Launch                              Landing
   * →→→→→→→→→→→→→→→→→→→→→→→→→→→ *
    \                             /
     \                           /
      \       At peak:          /
       \      ← only vₓ →      /
        \                     /
         \                   /
          \                 /
           \               /
            *             *
            
  t_up = t_down
  Speed at landing = Speed at launch
```

---

## Real-World Projectiles

Let us look at some everyday projectile situations to build intuition.

### Basketball Free Throw

A free throw line is 4.6 m from the basket. The basket is 3.05 m high. A player releases the ball from about 2.0 m height.

The ball must travel 4.6 m horizontally and rise about 1.05 m (3.05 - 2.0 = 1.05 m) to reach the basket. Players typically release at around 50-60° to give the ball enough height to have a good angle going into the hoop.

### Javelin Throw

In theory, a javelin would travel farthest at 45°. But in practice, aerodynamics and the mechanics of throwing mean athletes release at around 30-35° for maximum distance. The javelin is also designed to be aerodynamic, so air resistance plays a bigger role than for a dense ball.

### Water from a Hose

When you hold a garden hose at 45°, the water stream travels the farthest. At a steeper angle, the water goes higher but not as far. At a shallower angle, it goes lower and not as far either. This is projectile motion — each water droplet is a tiny projectile once it leaves the nozzle.

### Long Jump

A long jumper takes off at a relatively shallow angle (about 20-25° in practice) because they cannot generate enough vertical takeoff speed to use a higher angle effectively. The horizontal component of their speed (from the run-up) is the dominant factor.

---

## Effect of Air Resistance

Everything we have done in this chapter assumes no air resistance. This is a good approximation for:
- Dense, heavy objects (a shot put, a baseball, a bowling ball)
- Slow-moving objects
- Short distances

But real projectiles do experience air resistance, which changes the motion significantly:

### How Air Resistance Affects the Path

1. **Reduced range**: Air drag opposes the motion, slowing the projectile horizontally.

2. **Asymmetric path**: The trajectory is no longer a perfect symmetric parabola. The descending portion is steeper and shorter than the ascending portion.

3. **Optimal angle less than 45°**: With air resistance, the ideal launch angle for maximum range is less than 45° (typically around 35-42° depending on the object).

4. **Reduced maximum height**: Air drag also opposes vertical motion, reducing how high the projectile goes.

```
Without air resistance (ideal):
    
         *
       *   *
      *     *    (symmetric parabola)
     *       *
    *         *
   *           *
  *             *
 +---------------*

With air resistance (real):

        *
       * *
      *   *
     *     *  (steeper descent)
    *       **
   *          **
  *              ***
 +--------------------*
 
 Landing point is closer (less range).
 Descent is steeper than ascent.
```

For this course, we always work with the ideal no-air-resistance case unless told otherwise.

---

## The Full Parabolic Path

Let us appreciate the complete picture of what happens to a projectile over its full flight.

```
HEIGHT
(metres)
                                                         
  H    |               *                               
       |             *   *                            
       |           *       *                         
       |         *           *                      
  H/2  |       *               *                  
       |      *                 *                
       |    *                     *             
       |   *                       *           
       |  *                         *         
       | *                            *      
       |*                               *   
  0    +---+---+---+---+---+---+---+---+----> DISTANCE
       0                                   R
       
  At start (0,0):
    - vₓ = u cos θ (forward, constant)
    - vᵧ = u sin θ (upward, decreasing)
    - Speed = u
    - Direction: θ above horizontal
    
  At peak (R/2, H):
    - vₓ = u cos θ (unchanged)
    - vᵧ = 0 (momentarily not moving vertically)
    - Speed = u cos θ (minimum speed)
    - Direction: horizontal
    
  At landing (R, 0):
    - vₓ = u cos θ (unchanged)
    - vᵧ = -u sin θ (downward, same magnitude as initial)
    - Speed = u (same as launch)
    - Direction: θ below horizontal
```

---

## Common Mistakes to Avoid

### Mistake 1: Using Total Speed Instead of Component

When finding time to reach a certain height, use the VERTICAL component of velocity, not the total speed. When finding horizontal distance, use the HORIZONTAL component.

### Mistake 2: Forgetting That Horizontal Velocity Is Constant

Students sometimes apply g to the horizontal direction. Gravity only accelerates the projectile vertically. Horizontal velocity never changes (in the no-air-resistance case).

### Mistake 3: Plugging in the Wrong Sign for g

If you define upward as positive y (which is standard), then gravity gives a negative vertical acceleration: -g = -9.8 m/s². If you write y = ½ × g × t² (without the minus sign), you will get the wrong answer for anything thrown upward.

### Mistake 4: Assuming Maximum Height Occurs at Halfway Through Time

The maximum height occurs at T/2 only if the projectile is launched AND lands at the same height. If the projectile is launched from a cliff and lands below the launch height, the peak occurs before T/2.

### Mistake 5: Forgetting the Initial Vertical Velocity for Angled Launches

If a ball is launched at 20 m/s at 30°, the initial vertical velocity is NOT 20 m/s — it is 20 × sin(30°) = 10 m/s. Only use the full speed u as the initial vertical velocity if the launch is straight up (90°).

### Mistake 6: Using the Range Formula When Launch and Landing Heights Differ

The formula R = u² sin(2θ) / g only works when the projectile lands at the same height it was launched from. If it lands higher or lower, you must use x = vₓ₀ × t and solve for t from the vertical equation.

---

## Summary

- A **projectile** is an object moving freely under gravity after being launched.

- The **independence principle** means horizontal and vertical motions are completely separate — you analyze each independently.

- **Horizontal motion**: constant velocity vₓ = u cos(θ). No acceleration. x = vₓ × t.

- **Vertical motion**: constant downward acceleration g = 9.8 m/s². Starts at vᵧ₀ = u sin(θ) for angled launches or vᵧ₀ = 0 for horizontal launches.

- For a **horizontally launched projectile** from height h:
  - Time to land: t = sqrt(2h / g)
  - Range: R = u × sqrt(2h / g)

- For an **angled projectile** at angle θ with speed u:
  - Max height: H = (u sin θ)² / (2g)
  - Time of flight: T = 2u sin(θ) / g
  - Range: R = u² sin(2θ) / g

- **Maximum range** occurs at θ = 45°.

- **Complementary angles** (adding to 90°) give equal range.

- The path is a **parabola** due to constant horizontal speed and accelerating vertical motion.

- **Symmetry**: the speed at landing equals the speed at launch (if landing at the same height). Going up takes as long as coming down.

- **Air resistance** breaks the symmetry, reduces range, and makes the descending path steeper than the ascending path.

---

## Key Equations

| Quantity | Formula |
|---|---|
| Initial horizontal velocity | vₓ₀ = u cos(θ) |
| Initial vertical velocity | vᵧ₀ = u sin(θ) |
| Horizontal position | x = vₓ₀ × t |
| Vertical position | y = vᵧ₀ × t - ½ × g × t² |
| Vertical velocity at time t | vᵧ = vᵧ₀ - g × t |
| Time to reach max height | t_top = u sin(θ) / g |
| Maximum height | H = (u sin θ)² / (2g) |
| Total time of flight | T = 2u sin(θ) / g |
| Range (same launch and landing height) | R = u² sin(2θ) / g |
| Time to fall from height h (horizontal launch) | t = sqrt(2h / g) |
| Range for horizontal launch from height h | R = u × sqrt(2h / g) |

---

*End of Chapter 11*
