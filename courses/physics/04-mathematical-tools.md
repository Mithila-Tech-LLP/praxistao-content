# Chapter 04: Mathematical Tools for Physics

> **"Mathematics is the language in which the universe is written. You don't need to be a mathematician to do physics — you just need the right tools."**

---

## Table of Contents

- [4.1 Why Math Is Physics's Language](#41-why-math-is-physicss-language)
- [4.2 Algebra — Rearranging Equations](#42-algebra--rearranging-equations)
- [4.3 The Substitution Method](#43-the-substitution-method)
- [4.4 Proportionality](#44-proportionality)
- [4.5 The Pythagorean Theorem](#45-the-pythagorean-theorem)
- [4.6 Basic Trigonometry](#46-basic-trigonometry)
- [4.7 Angles — Degrees and Radians](#47-angles--degrees-and-radians)
- [4.8 Reading and Interpreting Graphs](#48-reading-and-interpreting-graphs)
- [4.9 Estimation and Fermi Problems](#49-estimation-and-fermi-problems)
- [4.10 Order of Magnitude](#410-order-of-magnitude)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 4.1 Why Math Is Physics's Language

Physics is about understanding how the universe works. We observe, measure, and look for patterns. But patterns expressed only in words are vague. "The harder you throw a ball, the faster it goes" is true — but it doesn't tell you exactly how fast. "v = u + at" is also true, and it tells you exactly how fast, for any specific values.

Every physics equation is a **precise, compact statement about reality**. It is a relationship that holds true everywhere in the universe under specified conditions.

Consider Newton's second law: **F = m × a**

In words, this says: "The force on an object equals its mass multiplied by its acceleration."

But as an equation, it says much more:
- If you double the force (and keep mass constant), acceleration doubles
- If you double the mass (and keep force constant), acceleration halves
- Given any two of the three values, you can calculate the third
- It holds on Earth, the Moon, Mars, and in galaxies billions of light-years away

This is why physics is inseparable from mathematics. Math is not a barrier to physics — it is the tool that makes precise understanding possible.

### What You Actually Need

The reassuring truth: most physics at the introductory level uses only:
- Arithmetic (adding, subtracting, multiplying, dividing)
- Algebra (rearranging equations)
- A little geometry (areas, triangles)
- Basic trigonometry (sine, cosine, tangent)
- Graph reading

You do not need calculus to understand most of what follows. We will build each tool from scratch.

---

## 4.2 Algebra — Rearranging Equations

**Algebra** is the art of rearranging equations to isolate an unknown variable. This is arguably the single most important mathematical skill in all of physics. If you can confidently rearrange equations, you can solve almost any problem.

### The One Golden Rule

**Whatever you do to one side of an equation, you must do to the other side.**

An equation is like a perfect balance scale. Both sides are equal. If you add 5 to the left side, you must add 5 to the right side, or the balance tips and the equation becomes false.

```
         LEFT SIDE = RIGHT SIDE

              5 kg  =  5 kg
             ======    ======
              |||        |||
              |||        |||

Add 2 kg to both sides:
       5 + 2 = 5 + 2
           7 = 7     ✓  (still balanced)

Add 2 kg only to left:
       5 + 2 = 5
           7 ≠ 5     ✗  (broken!)
```

**Allowed operations** (apply to BOTH sides equally):
- Add any number
- Subtract any number
- Multiply by any number (except zero)
- Divide by any number (except zero)
- Square both sides
- Take the square root of both sides

### Example 1: Rearranging F = m × a

**Find: a (acceleration)**

We want a alone on one side. It is currently multiplied by m.

```
Start:     F = m × a

Goal:      isolate a

Step:      divide both sides by m
           F / m = (m × a) / m

           On the right, m/m = 1, so:
           F / m = a

           Flip it to write a on the left:
           a = F / m
```

**Find: m (mass)**

```
Start:     F = m × a

Goal:      isolate m

Step:      divide both sides by a
           F / a = (m × a) / a
           F / a = m

           Therefore:
           m = F / a
```

### Example 2: Rearranging v = u + at

This equation appears constantly in kinematics.

**Find: t (time)**

```
Start:     v = u + at

Goal:      isolate t

Step 1:    subtract u from both sides (to get 'at' alone)
           v − u = u + at − u
           v − u = at

Step 2:    divide both sides by a
           (v − u) / a = at / a
           (v − u) / a = t

           Therefore:
           t = (v − u) / a
```

**Find: u (initial velocity)**

```
Start:     v = u + at

Step 1:    subtract at from both sides
           v − at = u

           Therefore:
           u = v − at
```

### Example 3: Rearranging E = ½mv²

**Find: v (velocity) — this requires square root**

```
Start:     E = ½ × m × v²

Goal:      isolate v

Step 1:    multiply both sides by 2 (to clear the ½)
           2E = 2 × ½ × m × v²
           2E = m × v²

Step 2:    divide both sides by m
           2E / m = v²

Step 3:    take the square root of both sides
           √(2E / m) = √(v²)
           √(2E / m) = v

           Therefore:
           v = √(2E / m)
```

### A Rearrangement Strategy

When you see an equation and need to rearrange it, work backwards from what's blocking the variable:

```
WHAT'S BLOCKING IT?        WHAT DO YOU DO?
--------------------------+---------------------------
Added to the variable     → Subtract it from both sides
Subtracted from variable  → Add it to both sides
Multiplied with variable  → Divide both sides by it
Dividing the variable     → Multiply both sides by it
Variable is squared (x²)  → Take square root of both sides
Variable is in square     → Square both sides
root (√x)                 →
```

### Practice: Rearrange These

```
P = F / A          Find F, then find A.
v² = u² + 2as      Find a, then find s.
F = G × m₁ × m₂ / r²    Find r.
```

Answers:
```
F = P × A
A = F / P

a = (v² − u²) / (2s)
s = (v² − u²) / (2a)

r² = G × m₁ × m₂ / F
r = √(G × m₁ × m₂ / F)
```

---

## 4.3 The Substitution Method

Once you can rearrange an equation, the next step is **substitution** — plugging in known numbers to get a numerical answer.

### The Method

```
STEP 1: Write down what you know (list all given values with units)
STEP 2: Write down what you want to find
STEP 3: Identify the equation that connects them
STEP 4: Rearrange the equation to isolate the unknown
STEP 5: Substitute the numbers WITH units
STEP 6: Calculate
STEP 7: Check — does the unit of the answer make sense?
         Does the numerical value seem reasonable?
```

### Worked Example 1: Finding Force

A car of mass 1200 kg accelerates at 3 m/s². What force acts on it?

```
Step 1: Known: m = 1200 kg, a = 3 m/s²
Step 2: Find: F (force)
Step 3: Equation: F = m × a
Step 4: Already rearranged — F is alone
Step 5: F = 1200 kg × 3 m/s²
Step 6: F = 3600 kg·m/s² = 3600 N
Step 7: Units are N (Newton) ✓  Size is reasonable for a car ✓

ANSWER: F = 3600 N = 3.6 kN
```

### Worked Example 2: Finding Time

A car accelerates from 0 to 30 m/s with an acceleration of 5 m/s². How long does it take?

```
Step 1: Known: u = 0 m/s (starts from rest), v = 30 m/s, a = 5 m/s²
Step 2: Find: t (time)
Step 3: Equation: v = u + at
Step 4: Rearrange for t:
         v − u = at
         t = (v − u) / a

Step 5: t = (30 − 0) m/s / (5 m/s²)
Step 6: t = 30 / 5 s = 6 s
         (Note: (m/s) / (m/s²) = (m/s) × (s²/m) = s ✓)
Step 7: Units are seconds ✓  6 seconds to reach 108 km/h — fast car ✓

ANSWER: t = 6 s
```

### Worked Example 3: Finding Speed from Kinetic Energy

A 2 kg ball has kinetic energy 100 J. What is its speed?

```
Step 1: Known: m = 2 kg, E = 100 J
Step 2: Find: v (speed)
Step 3: Equation: E = ½mv²
Step 4: Rearrange for v:
         2E = mv²
         v² = 2E / m
         v = √(2E / m)

Step 5: v = √(2 × 100 / 2)
Step 6: v = √(200 / 2) = √100 = 10

         Check units: √(J / kg) = √(kg·m²/s² / kg) = √(m²/s²) = m/s ✓

ANSWER: v = 10 m/s
```

### Common Substitution Errors to Avoid

```
ERROR 1: Forgetting to convert units before substituting.
  Wrong: using speed in km/h in an equation that needs m/s.
  Fix: always convert to SI units first.

ERROR 2: Substituting before rearranging.
  It's much harder to rearrange after substituting numbers.
  Always rearrange algebraically first.

ERROR 3: Losing track of negative signs.
  If a car decelerates, a is negative. Keep the sign.

ERROR 4: Not checking units in the answer.
  If you get N when you expect m/s, something went wrong.
```

---

## 4.4 Proportionality

**Proportionality** describes how two quantities relate to each other when one changes. Understanding proportionality lets you make predictions without knowing the exact equation.

### Direct Proportion: y ∝ x

If y is **directly proportional** to x, doubling x doubles y. Tripling x triples y.

Written: y ∝ x   (the ∝ symbol means "is proportional to")

With a constant of proportionality k: y = k × x

```
GRAPH: Direct proportion — straight line through the origin

      y
      |              /
      |            /
      |          /
      |        /       slope = k
      |      /
      |    /
      |  /
      |/
      +----------------> x
      0

Key features:
  - Passes through origin (0, 0)
  - Straight line
  - Steeper line = larger k
```

Physics examples of direct proportion:
- Force and acceleration (F ∝ a when mass is constant): double the force, double the acceleration
- Hooke's Law — spring extension and force (F = kx): pull twice as hard, spring stretches twice as far
- Ohm's Law — voltage and current (V = IR): double the voltage, double the current (at constant R)

### Inverse Proportion: y ∝ 1/x

If y is **inversely proportional** to x, doubling x halves y. Tripling x reduces y to one-third.

Written: y ∝ 1/x    or equivalently    y × x = constant

With constant k: y = k / x

```
GRAPH: Inverse proportion — hyperbola (curve, never touches axes)

      y
      |
   10 |*
      |
    5 | *
      |  *
    3 |    *
    2 |       *
    1 |             *    *    *
      +--1---2---3---4---5----> x
      0

Key features:
  - Never touches either axis (as x → 0, y → infinity; as x → ∞, y → 0)
  - Not a straight line — it's a hyperbola
  - As x increases, y decreases
```

Physics examples of inverse proportion:
- Force and acceleration (a ∝ 1/m when force is constant): double the mass, halve the acceleration
- Gravitational force and distance squared: F ∝ 1/r² (inverse square law)
- Gas pressure and volume (Boyle's Law): P ∝ 1/V at constant temperature

### Square Proportion: y ∝ x²

If y is proportional to **x squared**, doubling x quadruples y. Tripling x multiplies y by nine.

Written: y ∝ x²

With constant k: y = k × x²

```
GRAPH: Square proportion — parabola

      y
      |                   *
   16 |
      |               *
    9 |
      |           *
    4 |       *
      |   *
    1 | *
      +--1---2---3---4----> x
      0

Key features:
  - Passes through origin (0, 0)
  - Curves upward (each step gets bigger)
  - Symmetrical around y-axis
```

Physics examples of square proportion:
- Kinetic energy and speed (E_k = ½mv²): double the speed, quadruple the kinetic energy
- Distance fallen and time (s = ½gt²): in twice the time, something falls four times as far
- Gravitational force and distance (F ∝ 1/r²): the inverse square law

### Comparing the Three Types

```
When x doubles (×2):
+----------------------+------------------+--------------+
|  Type                |  What happens    |  Equation    |
|                       |  to y           |  form        |
+----------------------+------------------+--------------+
| Direct (y ∝ x)       |  y doubles (×2) |  y = kx      |
| Square (y ∝ x²)      |  y × 4          |  y = kx²     |
| Inverse (y ∝ 1/x)    |  y halves (÷2)  |  y = k/x     |
| Inv. square (y∝1/x²) |  y × 1/4        |  y = k/x²    |
+----------------------+------------------+--------------+
```

### Worked Example: Using Proportionality

A car travelling at 20 m/s has kinetic energy 80,000 J. What is its kinetic energy at 40 m/s?

```
Kinetic energy: E = ½mv²   →   E ∝ v²

If v doubles: E multiplies by 2² = 4

New energy = 80,000 × 4 = 320,000 J

ANSWER: 320,000 J (= 320 kJ)

(No need to find mass separately — proportionality gives us the answer directly.)
```

---

## 4.5 The Pythagorean Theorem

The **Pythagorean theorem** is one of the most ancient and useful results in all of mathematics. In physics, it appears whenever we deal with vectors at right angles — which is extremely common.

### The Theorem

In any right-angled triangle:

```
         c (hypotenuse)
        /|
       / |
      /  |  b (opposite)
     /   |
    /    |
   /_____|
  a (adjacent)
  
  Right angle is at the bottom right corner.

  a² + b² = c²

  where c is the HYPOTENUSE (the side opposite the right angle — always the longest side)
        a and b are the other two sides
```

To find the hypotenuse: c = √(a² + b²)

To find a side: a = √(c² − b²)

### Application 1: Resultant Velocity of a Boat Crossing a River

A boat travels straight across a river at 4 m/s. The river current flows at 3 m/s downstream. What is the boat's actual speed (its resultant velocity)?

```
The velocities are perpendicular to each other:

      Current → 3 m/s
      ===========================
      |                         |
 4    |   ↑ boat                |
 m/s  |   |                     |
      |   |  but it also        |
      |   | drifts right        |
      ===========================

Draw the velocity triangle:

           R (resultant)
          /|
         / |
      R /  | 4 m/s (boat's own speed)
       /   |
      /____|
       3 m/s (river current)

The two velocities are the two legs of a right triangle.
The resultant is the hypotenuse.

Using Pythagoras:
R² = 3² + 4²
R² = 9 + 16
R² = 25
R = √25 = 5 m/s

ANSWER: The boat moves at 5 m/s relative to the riverbank.
(This is the famous 3-4-5 right triangle!)
```

### Application 2: Net Force From Two Perpendicular Forces

Two forces act on an object: 60 N to the right and 80 N upward. What is the resultant force?

```
         ↑ 80 N
         |
         |
         |
         +--------→ 60 N

The forces form two sides of a right triangle.
The resultant force F_net is the hypotenuse.

F_net² = 60² + 80²
F_net² = 3600 + 6400
F_net² = 10,000
F_net = √10,000 = 100 N

ANSWER: F_net = 100 N

(Another famous right triangle: 60-80-100, or simplified 3-4-5 scaled by 20)
```

### Application 3: Distance in Two Dimensions

A person walks 12 m east and then 5 m north. How far are they from the starting point?

```
       N ↑
         |
         |  5 m
         |
         +--------→ E
              12 m

Distance = √(12² + 5²)
         = √(144 + 25)
         = √169
         = 13 m

ANSWER: 13 m from the starting point.
```

---

## 4.6 Basic Trigonometry

**Trigonometry** is the mathematics of triangles, specifically the relationships between angles and the lengths of sides. In physics, it is essential for working with forces, velocities, and any quantity that has both a size and a direction.

### Right Triangle — The Key Diagram

Every trigonometric ratio is defined for a right-angled triangle. Consider this triangle with angle θ (theta) at the bottom left:

```
              C
             /|
            / |
           /  |
       H  /   |  O
   (hyp) /    |  (opposite
         / θ  |   side)
        /_____|
       A       B
              (right angle at B)

  H = Hypotenuse  (side opposite the right angle — always longest)
  O = Opposite    (side opposite the angle θ — the one you're focused on)
  A = Adjacent    (side next to the angle θ, not the hypotenuse)
```

### SOHCAHTOA

The three main trig ratios are defined as:

```
SOH:  sin(θ) = Opposite / Hypotenuse    =  O / H

CAH:  cos(θ) = Adjacent / Hypotenuse    =  A / H

TOA:  tan(θ) = Opposite / Adjacent      =  O / A

Memory trick:  S O H  C A H  T O A
               → SOH-CAH-TOA
               → "Some Old Hens Can Always Hop Through Any Obstacle"
```

### Standard Values to Memorise

```
+-------+----------+----------+----------+
| Angle |  sin(θ)  |  cos(θ)  |  tan(θ)  |
+-------+----------+----------+----------+
|   0°  |    0     |    1     |    0     |
|  30°  |   0.5    |  0.866   |  0.577   |
|  45°  |  0.707   |  0.707   |    1     |
|  60°  |  0.866   |   0.5    |  1.732   |
|  90°  |    1     |    0     |    ∞     |
+-------+----------+----------+----------+

Exact fractions:
sin(30°) = 1/2,   cos(30°) = √3/2,  tan(30°) = 1/√3
sin(45°) = √2/2,  cos(45°) = √2/2,  tan(45°) = 1
sin(60°) = √3/2,  cos(60°) = 1/2,   tan(60°) = √3
```

Memory trick for 30°, 45°, 60°:

```
sin:  1/2,  √2/2,  √3/2    ← numerators go √1, √2, √3
cos:  √3/2, √2/2,  1/2     ← exactly reversed
```

### Finding an Angle — Inverse Trig

If you know the ratio but want the angle, use the **inverse trig functions**:

```
If sin(θ) = 0.5,   then   θ = sin⁻¹(0.5) = 30°
If cos(θ) = 0.707, then   θ = cos⁻¹(0.707) = 45°
If tan(θ) = 1.732, then   θ = tan⁻¹(1.732) = 60°
```

On a calculator: look for buttons labeled sin⁻¹, cos⁻¹, tan⁻¹ (often accessed with SHIFT or 2nd).

### Worked Physics Example 1: Resolving a Force

A rope pulls a crate along the floor at an angle of 30° above the horizontal. The tension in the rope is 100 N. What are the horizontal and vertical components of this force?

```
                   T = 100 N
                   ↗ at 30° above horizontal
            30°  /
           ─────/──────────→  horizontal

The force has two components:
   Horizontal: F_x = T × cos(30°)
   Vertical:   F_y = T × sin(30°)

   F_x = 100 × cos(30°) = 100 × 0.866 = 86.6 N
   F_y = 100 × sin(30°) = 100 × 0.5   = 50 N

ANSWER: Horizontal component = 86.6 N
        Vertical component = 50 N
```

### Worked Physics Example 2: Inclined Plane

A block sits on a slope that makes a 40° angle with the horizontal. The block has mass 5 kg. What component of gravity acts along the slope (trying to pull it downslide)?

```
              /|
             / |
            /  |
           / 40|°
          /____|

Gravity (W = mg) acts straight DOWN.
The slope makes 40° with horizontal.

The angle between W and the "along-slope" direction is (90° − 40°) = 50°.
But more directly: the component along the slope = mg × sin(40°)

  W_along = m × g × sin(40°)
          = 5 × 9.8 × 0.643
          = 31.5 N

ANSWER: 31.5 N acts along the slope, tending to slide the block downward.
```

### Worked Physics Example 3: Finding an Angle

A force of 60 N to the right and 80 N upward act on an object. The resultant is 100 N (from Pythagoras). What angle does the resultant make with the horizontal?

```
         ↑ 80 N
         |
         |
  100 N ↗|
      θ  |
      ───+───────→ 60 N

tan(θ) = opposite / adjacent = 80 / 60 = 1.333

θ = tan⁻¹(1.333) = 53.1°

ANSWER: The resultant force acts at 53.1° above the horizontal.
```

---

## 4.7 Angles — Degrees and Radians

Angles can be measured in two common systems: **degrees** and **radians**. Degrees are familiar from everyday life. Radians are the natural mathematical unit and are used in advanced physics and calculus.

### Degrees

A full revolution is divided into **360 degrees (°)**. This number 360 is historical — it was convenient for the ancient Babylonians (360 is divisible by many numbers: 2, 3, 4, 5, 6, 8, 9, 10, 12, 15, 18, 20, 24...).

```
Full circle  = 360°
Half circle  = 180°  (straight angle)
Quarter      = 90°   (right angle)
```

### Radians

A **radian** is defined by the geometry of a circle itself. If you take the radius of a circle and wrap it along the circumference, it subtends (covers) exactly 1 radian at the centre.

```
           ____
         /      \
        /    r   \
       |          |
       |    r___  |     arc length = radius
       |   /    | |     → angle = 1 radian
        \ /     /
         \____/

1 radian ≈ 57.3°
```

Since the full circumference = 2πr, wrapping all the way around = 2π radians.

```
2π radians = 360°
π radians  = 180°
1 radian   = 180°/π ≈ 57.3°
```

### Converting Between Degrees and Radians

```
Degrees to radians:  radians = degrees × (π / 180)
Radians to degrees:  degrees = radians × (180 / π)
```

Common angles:

```
+----------+-----------+----------------------------+
|  Degrees |  Radians  |   Exact radian expression  |
+----------+-----------+----------------------------+
|    0°    |   0       |   0                        |
|   30°    |  0.524    |   π/6                      |
|   45°    |  0.785    |   π/4                      |
|   60°    |  1.047    |   π/3                      |
|   90°    |  1.571    |   π/2                      |
|  180°    |  3.142    |   π                        |
|  270°    |  4.712    |   3π/2                     |
|  360°    |  6.283    |   2π                       |
+----------+-----------+----------------------------+
```

### When to Use Radians

In introductory physics, angles are usually given in degrees and trig functions on your calculator can handle degrees. Radians become essential when:
- Calculating arc length: s = r × θ (θ in radians)
- Circular motion: angular velocity ω = angle / time
- Using calculus in physics

---

## 4.8 Reading and Interpreting Graphs

Graphs are the visual language of physics. Every graph tells a story. Learning to read that story quickly is a crucial skill.

### The Anatomy of a Good Graph

```
TITLE: Distance vs Time for a Moving Car

  d (m)
  |  ↑ y-axis
  |  labelled with
  |  quantity AND unit
80|              *
  |           *
60|        *
  |     *
40|  *
  | *
20|*
  +--+--+--+--+--+--→  t (s)
  0  2  4  6  8  10
     ↑
     x-axis labelled
     with quantity AND unit

Every axis must have:  quantity name + unit in brackets
Every graph must have: a title
Data points are shown as dots, not lines
Best-fit line is drawn smoothly through the data
```

### The Gradient (Slope) — The Most Important Feature

The **gradient** (or **slope**) of a graph tells you the **rate of change** of the y-variable with respect to the x-variable.

```
Gradient = rise / run = Δy / Δx

"Rise" = change in vertical (y) value
"Run" = change in horizontal (x) value
```

To find the gradient:
1. Pick TWO points on the line (NOT data points — use the line itself)
2. Read off their coordinates
3. Calculate Δy = y₂ − y₁
4. Calculate Δx = x₂ − x₁
5. Divide: gradient = Δy / Δx

```
Example: A line passes through (2, 20) and (8, 80)

Gradient = (80 − 20) / (8 − 2)
         = 60 / 6
         = 10

If this is a d-t graph, gradient = velocity = 10 m/s
```

### What Different Gradients Mean

```
POSITIVE gradient:  y increases as x increases
NEGATIVE gradient:  y decreases as x increases (e.g., deceleration)
ZERO gradient:      y is constant — horizontal line
LARGE gradient:     fast rate of change (steep line)
SMALL gradient:     slow rate of change (shallow line)
INFINITE gradient:  vertical line — instantaneous change
```

### The Area Under a Graph — The Second Key Feature

The **area under a curve** (between the curve and the x-axis) represents an **accumulated quantity** — what you get when you add up (integrate) y over the x range.

```
Example: On a velocity-time graph
  y-axis = velocity (m/s)
  x-axis = time (s)
  Area = velocity × time = distance

  A rectangle of width 5s and height 10 m/s:
  Area = 5 × 10 = 50 m    (50 metres of displacement)
```

This is one of the most powerful ideas in physics: the area under a v-t graph is the displacement. The area under an F-t graph is impulse (change in momentum). The area under a P-V graph is work done.

### Common Graph Shapes and Their Physics Meaning

```
SHAPE               WHAT IT MEANS IN PHYSICS
--------------------+-----------------------------------------------
Horizontal line     | Constant value (not changing)
Straight line       | Constant rate of change (linear relationship)
through origin      |
Straight line       | Linear relationship with non-zero intercept
not through origin  |
Curve upward        | Rate of change is increasing (e.g., accelerating)
(parabola)          |
Curve downward      | Rate of change is decreasing (e.g., decelerating)
Hyperbola           | Inverse proportion (y = k/x)
S-shaped curve      | Growth with a limit (common in biology/thermodynamics)
```

### Worked Example: Finding Velocity From a d-t Graph

```
d (m)
|
50|               *
  |           *
40|       *
  |   *
30| *
  |
  +--2--4--6--8--→ t (s)

Two points on the line: (2, 30) and (8, 50)

Gradient = (50 − 30) / (8 − 2)
         = 20 / 6
         = 3.33 m/s

ANSWER: The object moves at a constant velocity of 3.33 m/s.
```

### Worked Example: Finding Distance From a v-t Graph

```
v (m/s)
|
20|******* (constant 20 m/s for 5 seconds)
  |       |
10|       |
  |       |
  +--1--2--3--4--5--→ t (s)
  |_______|
   Area = 20 × 5 = 100 m

ANSWER: The object travelled 100 m in 5 seconds.
```

---

## 4.9 Estimation and Fermi Problems

One of the most impressive skills a physicist has is the ability to make a **reasonable estimate** of a quantity even when exact data is unavailable. These are called **Fermi problems** after the physicist Enrico Fermi, who was famous for making rapid, surprisingly accurate estimates.

The technique relies on:
1. Breaking a complex unknown into simpler knowns
2. Making reasonable assumptions
3. Multiplying everything together
4. Checking whether the answer makes sense

The goal is not to be exactly right — it is to be within a factor of 10 (one order of magnitude).

### The Fermi Method

```
STEP 1: Identify the unknown quantity
STEP 2: Break it down into factors you can estimate
STEP 3: Estimate each factor individually
STEP 4: Multiply all the factors together
STEP 5: Check the result against any known references
```

### Classic Example 1: How Many Piano Tuners Are in Chicago?

This is the most famous Fermi problem.

```
What we need to find: number of piano tuners in Chicago

Break it down:

  Factor 1: Population of Chicago
  Estimate: about 3 million people = 3 × 10⁶

  Factor 2: How many people per piano?
  Estimate: about 1 piano per 20 households
  Average household: 2–3 people, say 2.5
  So 1 piano per 50 people

  Number of pianos in Chicago = 3,000,000 / 50 = 60,000 pianos

  Factor 3: How often does a piano get tuned?
  Estimate: once a year, on average

  Factor 4: How many pianos can one tuner tune per day?
  A tuning takes about 2 hours. Working day = 8 hours.
  So 4 pianos per day.

  Factor 5: Working days per year = ~250

  One tuner can tune: 4 × 250 = 1,000 pianos per year

  Number of tuners needed = 60,000 pianos / 1,000 per tuner
                          = 60 piano tuners

Actual answer: roughly 100–150 piano tuners in Chicago.
We got 60. That's within a factor of 2. Excellent estimate!
```

### Classic Example 2: How Many Litres of Water Do You Drink in a Lifetime?

```
Break it down:

  Factor 1: Daily water intake
  Estimate: about 2 litres per day (drinking water, coffee, juice, etc.)

  Factor 2: Lifespan
  Estimate: about 75 years

  Factor 3: Days in a year = 365

  Total = 2 L/day × 365 days/year × 75 years
        = 2 × 365 × 75
        = 54,750 litres
        ≈ 55,000 litres = 5.5 × 10⁴ litres

ANSWER: Roughly 55,000 litres (about 55 cubic metres) in a lifetime.
That is enough to fill a small swimming pool.
```

### Classic Example 3: What Is the Mass of All Humans on Earth?

```
Break it down:

  Factor 1: World population = about 8 billion = 8 × 10⁹

  Factor 2: Average mass per person
  Estimate: about 65 kg (considering all ages and sizes globally)

  Total mass = 8 × 10⁹ × 65 kg
             = 5.2 × 10¹¹ kg
             ≈ 5 × 10¹¹ kg (500 billion kg = 500 million tonnes)

ANSWER: About 5 × 10¹¹ kg — roughly 500 million tonnes.

Cross-check: A large oil supertanker holds about 300,000 tonnes.
This is equivalent to about 1,600 supertankers full of people. Plausible.
```

### Why Physicists Estimate

Estimation is valuable because:
- It gives you a quick sanity check on calculated answers
- It helps you understand the scale of a problem before solving it
- It builds physical intuition — a feel for what "big" and "small" mean
- It is often all you need to make a decision

Famous example: On the first nuclear bomb test (Trinity, 1945), Fermi estimated the yield of the explosion by dropping small pieces of paper and watching how far the shock wave blew them. He estimated 10 kilotons. The actual yield was 18.6 kilotons. Not bad for torn-up paper!

---

## 4.10 Order of Magnitude

**Order of magnitude** is a way of comparing quantities by how many powers of ten they differ. Two quantities that differ by one order of magnitude differ by a factor of about 10. Two orders of magnitude is a factor of 100. And so on.

### What "Order of Magnitude" Means

```
Order of magnitude = the power of 10 in scientific notation

Examples:
  5 × 10²  → order of magnitude: 2   (roughly 100)
  3 × 10⁷  → order of magnitude: 7   (roughly 10,000,000)
  8 × 10⁻³ → order of magnitude: −3  (roughly 0.001)
```

When comparing two quantities, the number of orders of magnitude difference tells you how many times bigger one is than the other:

```
Difference of 1 order of magnitude = factor of ~10
Difference of 2 orders of magnitude = factor of ~100
Difference of 6 orders of magnitude = factor of ~1,000,000
```

### A Scale of Masses

```
Object                     Mass (kg)            Order of Magnitude
--------------------------+--------------------+--------------------
Electron                  | 9.1 × 10⁻³¹ kg    |  −31
Proton                    | 1.7 × 10⁻²⁷ kg    |  −27
Water molecule (H₂O)      | 3.0 × 10⁻²⁶ kg    |  −26
Bacterium                 | 1 × 10⁻¹⁵ kg      |  −15
Mosquito                  | 2.5 × 10⁻⁶ kg     |   −6
Penny                     | 2.5 × 10⁻³ kg     |   −3
Human adult               | 7 × 10¹ kg        |    1
Car                       | 1.5 × 10³ kg      |    3
Blue whale                | 1.5 × 10⁵ kg      |    5
Great Pyramid             | 6 × 10⁹ kg        |    9
Earth                     | 6 × 10²⁴ kg       |   24
Sun                       | 2 × 10³⁰ kg       |   30
```

From electron to Sun: 30 − (−31) = 61 orders of magnitude. The Sun is 10⁶¹ times more massive than an electron.

### A Scale of Lengths

```
Object                     Length (m)           Order of Magnitude
--------------------------+--------------------+--------------------
Quark                     | ~10⁻¹⁸ m           |  −18
Proton diameter           | ~10⁻¹⁵ m           |  −15
Atom diameter             | ~10⁻¹⁰ m           |  −10
DNA strand width          | 2 × 10⁻⁹ m         |   −9
Red blood cell            | 8 × 10⁻⁶ m         |   −6
Human hair width          | 7 × 10⁻⁵ m         |   −5
Human height              | ~1.7 × 10⁰ m       |    0
Eiffel Tower              | 3.3 × 10² m        |    2
Mount Everest             | 8.8 × 10³ m        |    3
Earth radius              | 6.4 × 10⁶ m        |    6
Earth-Moon distance       | 3.8 × 10⁸ m        |    8
Earth-Sun distance        | 1.5 × 10¹¹ m       |   11
Milky Way diameter        | 10²¹ m             |   21
Observable universe       | 8.8 × 10²⁶ m       |   26
```

### Worked Example: Comparing Sizes

How many orders of magnitude larger is the Sun than a proton?

```
Diameter of Sun: 1.4 × 10⁹ m   → order of magnitude: 9
Diameter of proton: 1.7 × 10⁻¹⁵ m → order of magnitude: −15

Difference: 9 − (−15) = 24 orders of magnitude

The Sun is roughly 10²⁴ times larger than a proton.
```

### Using Orders of Magnitude as a Sanity Check

After calculating an answer, ask: "Does this have roughly the right order of magnitude?"

```
Example: You calculate the speed of a car as 250 m/s.
   A reasonable car speed is about 30 m/s (about 100 km/h).
   250 m/s is nearly 10 times too large — probably an error.
   
   (250 m/s ≈ 900 km/h — faster than most aircraft)
   Go back and check your calculation.
```

This quick check catches most arithmetic errors in physics problems.

---

## Summary

- **Every physics equation is a precise statement about reality.** Mathematics is the tool that makes physics precise, predictive, and universal.

- **Algebra** is the most important skill: to find an unknown variable, do the same operation to both sides. Rearrange symbolically BEFORE substituting numbers.

- **The substitution method** has a clear procedure: list knowns, identify the equation, rearrange for the unknown, substitute with units, calculate, check.

- **Proportionality** describes how two quantities are related:
  - Direct (y ∝ x): doubling x doubles y — straight line through origin
  - Inverse (y ∝ 1/x): doubling x halves y — hyperbola
  - Square (y ∝ x²): doubling x quadruples y — parabola

- **The Pythagorean theorem** (a² + b² = c²) is used whenever two perpendicular quantities combine — resultant velocities, net perpendicular forces, distances in 2D.

- **Trigonometry** (SOHCAHTOA): sin = O/H, cos = A/H, tan = O/A. Used to find components of vectors and angles in right triangles. Know values for 30°, 45°, 60°. Use inverse trig to find angles.

- **Degrees and radians** are two ways to measure angles. 2π radians = 360°. Convert with: radians = degrees × (π/180).

- **Graph reading**: gradient = Δy/Δx = rate of change. Area under curve = accumulated quantity. Both are powerful tools for extracting information from graphs.

- **Fermi estimation** breaks a complex unknown into estimable factors. The goal is to be within one order of magnitude — it builds physical intuition.

- **Order of magnitude** = power of ten in scientific notation. Comparing orders of magnitude tells you how many times bigger one quantity is than another.

---

## Key Equations

```
ALGEBRA RULES:
    Whatever you do to one side, do to the other.
    To isolate x multiplied by k: divide both sides by k → x = .../k
    To isolate x + k: subtract k from both sides → x = ... - k
    To isolate x²: take square root of both sides → x = √(...)

COMMON REARRANGEMENTS:
    F = m × a    →   a = F/m,  m = F/a
    v = u + at   →   t = (v−u)/a,  u = v−at
    E = ½mv²     →   v = √(2E/m),  m = 2E/v²
    P = F/A      →   F = P×A,  A = F/P
    v² = u² + 2as → a = (v²−u²)/(2s), s = (v²−u²)/(2a)

PYTHAGOREAN THEOREM:
    c = √(a² + b²)   [hypotenuse from two legs]
    a = √(c² − b²)   [leg from hypotenuse and other leg]

TRIGONOMETRY (SOHCAHTOA):
    sin(θ) = Opposite / Hypotenuse  →  O = H × sin(θ)
    cos(θ) = Adjacent / Hypotenuse  →  A = H × cos(θ)
    tan(θ) = Opposite / Adjacent    →  θ = tan⁻¹(O/A)

ANGLE CONVERSION:
    radians = degrees × (π / 180)
    degrees = radians × (180 / π)
    Full circle: 360° = 2π radians

PROPORTIONALITY:
    y ∝ x     →   y = kx         (direct)
    y ∝ 1/x   →   y = k/x        (inverse)
    y ∝ x²    →   y = kx²        (square)

GRAPH GRADIENT:
    gradient = Δy / Δx = (y₂ − y₁) / (x₂ − x₁)
    = rate of change of y with respect to x
```
