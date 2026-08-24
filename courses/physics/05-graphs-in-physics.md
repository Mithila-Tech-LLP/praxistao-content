# Chapter 05: Graphs in Physics

> **"A graph is worth a thousand data points. It transforms numbers into understanding — showing trends, relationships, and stories hidden in raw data."**

---

## Table of Contents

- [5.1 Why Graphs Matter in Physics](#51-why-graphs-matter-in-physics)
- [5.2 How to Draw a Proper Graph](#52-how-to-draw-a-proper-graph)
- [5.3 Position-Time Graphs](#53-position-time-graphs)
- [5.4 Velocity-Time Graphs](#54-velocity-time-graphs)
- [5.5 Acceleration-Time Graphs](#55-acceleration-time-graphs)
- [5.6 Comparing Motion Types on Graphs](#56-comparing-motion-types-on-graphs)
- [5.7 Finding Gradient (Slope)](#57-finding-gradient-slope)
- [5.8 Finding Area Under a Graph](#58-finding-area-under-a-graph)
- [5.9 Non-Linear Graphs](#59-non-linear-graphs)
- [5.10 Linearizing Data](#510-linearizing-data)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 5.1 Why Graphs Matter in Physics

Imagine you are given this table of data from an experiment where a ball was dropped:

| Time (s) | Height (m) |
|----------|-----------|
| 0        | 45.0      |
| 1        | 40.1      |
| 2        | 25.4      |
| 3        | 0.9       |

Looking at the numbers, you can see the ball is falling. But you cannot immediately see *how fast* it is accelerating, or whether the relationship is linear or curved.

Now look at the same data on a graph:

```
HEIGHT (m)
50 |*
45 | *
40 |   *
35 |
30 |
25 |    *
20 |
15 |
10 |
 5 |
 0 |          *
   +-------------> TIME (s)
   0  1  2  3
```

Instantly, the curve tells you the ball is accelerating — the height drops faster and faster. The visual pattern reveals information the table hides.

This is why physicists use graphs constantly. A graph:
- Shows the **shape** of a relationship (linear, curved, exponential)
- Makes **trends** immediately visible
- Allows you to **extract information** (by reading off the slope or area)
- Reveals **anomalies** — points that don't fit the pattern
- Communicates results to others at a glance

---

## 5.2 How to Draw a Proper Graph

A graph in physics must meet specific standards to be meaningful and readable.

### The Essential Elements

```
ANATOMY OF A PHYSICS GRAPH

                    TITLE
     Position of Car vs. Time
          
y-axis   Position (m)    <- axis label + unit in brackets
label ->  |
         |         /
      30 |        /
         |       /
      20 |      /    <- plotted data points
         |     /        connected with a line
      10 |    /         of best fit
         |   /
       0 +-----------> Time (s)   <- x-axis label + unit
         0   1   2   3
         
         <----- origin should be clearly shown
```

### Rules for Drawing Graphs

**1. Label both axes** with the quantity AND its unit in brackets.
- Correct: "Velocity (m/s)"
- Wrong: "Velocity" or "m/s"

**2. Choose an appropriate scale**
- Use the data to fill most of the graph
- Use easy-to-read scales: 1, 2, 5, 10 per division — NOT 3, 7, or 9 per division
- The scale does not have to start at zero (unless that is important)

**3. Plot points clearly**
- Use small, neat crosses (x) or dots with a circle
- Do not connect points with a jagged zigzag line

**4. Draw a line or curve of best fit**
- For a linear relationship: draw a single straight line through the points, balanced so some points are above and some below
- For a curved relationship: draw a smooth curve through the points
- Do NOT force the line through the origin unless you have good reason

**5. Write a title** — what was measured against what

**6. Include a key** if there are multiple data sets on the same graph

### Common Mistakes to Avoid

```
WRONG                          RIGHT
---------                      ---------
   y                              y (m)
   |                              |
   |  *  *                        | * *
   | *  *                         |* *
   +--->                          +------> t (s)
      x                              
      
(no units, no title)           (labeled axes, units shown)
```

---

## 5.3 Position-Time Graphs

The **position-time graph** (also called an **x-t graph**) shows where an object is at each moment in time.

The x-axis is time, the y-axis is position.

### What Different Shapes Mean

**Shape 1: Flat horizontal line — Object is at rest**

```
Position (m)
20 |--------------------
   |
   |
 0 +--------------------> Time (s)
```
The position does not change. The object is stationary.

---

**Shape 2: Straight diagonal line (positive slope) — Constant velocity moving forward**

```
Position (m)
30 |              /
20 |            /
10 |          /
 0 +-----------> Time (s)
```
Equal increases in position for equal time intervals — constant velocity.

---

**Shape 3: Straight diagonal line (negative slope) — Constant velocity moving backward**

```
Position (m)
30 |  \
20 |    \
10 |      \
 0 +-----------> Time (s)
```
The object is moving in the negative direction at constant speed.

---

**Shape 4: Curved line (concave up) — Speeding up (positive acceleration)**

```
Position (m)
40 |              /
20 |          /
10 |       /
 5 |     /
 0 +-----------> Time (s)
```
The position increases by more and more each second — the object accelerates.

---

**Shape 5: Curved line (concave down) — Slowing down**

```
Position (m)
40 |        ----------
30 |      /
20 |    /
 0 +-----------> Time (s)
```
The object slows to a stop.

### The KEY RULE: Slope of x-t graph = velocity

```
         change in position    delta_x
slope = -------------------- = ------- = velocity
         change in time         delta_t
```

A steeper slope = greater velocity.
A flat slope = zero velocity.
A negative slope = negative velocity (moving backward).
A changing slope = changing velocity = acceleration.

### Worked Example 5.1

From the graph below, find the velocity between t = 1s and t = 4s.

```
Position (m)
20 |           *
15 |         *
10 |       *
 5 |     *
 0 +---*------ Time (s)
   0  1  2  3  4
```

**Solution:**
- At t = 1s: position = 5 m
- At t = 4s: position = 20 m
- Change in x = 20 - 5 = 15 m
- Change in t = 4 - 1 = 3 s
- Velocity = 15/3 = **5 m/s**

---

## 5.4 Velocity-Time Graphs

The **velocity-time graph** (v-t graph) shows how fast an object is moving at each moment.

The x-axis is time, the y-axis is velocity.

### What Different Shapes Mean

**Shape 1: Flat horizontal line — Constant velocity (zero acceleration)**

```
Velocity (m/s)
20 |--------------------
   |
   |
 0 +--------------------> Time (s)
```

---

**Shape 2: Straight line with positive slope — Constant acceleration (speeding up)**

```
Velocity (m/s)
30 |              /
20 |            /
10 |          /
 0 +-----------> Time (s)
```
Equal increases in velocity each second — constant acceleration.

---

**Shape 3: Straight line with negative slope — Constant deceleration (slowing down)**

```
Velocity (m/s)
30 |  \
20 |    \
10 |      \
 0 +-----------> Time (s)
```

---

**Shape 4: Starting below zero — Moving in negative direction**

```
Velocity (m/s)
 0 +-----------> Time (s)
-5 |--------
-10|
```

### KEY RULE 1: Slope of v-t graph = acceleration

```
         change in velocity    delta_v
slope = -------------------- = ------- = acceleration
         change in time         delta_t
```

Positive slope = positive acceleration.
Negative slope = deceleration.
Zero slope = constant velocity.

### KEY RULE 2: Area under v-t graph = displacement

This is extremely important.

```
Velocity (m/s)
20 |*--------*
   |         |
   |         |
 0 +-----------> Time (s)
   0    5    10

Area = base x height = 10 x 20 = 200 m (displacement)
```

For a triangle (accelerating from rest):

```
Velocity (m/s)
20 |          *
   |         /|
   |        / |
 0 +-------/---> Time (s)
   0       5

Area = 1/2 x base x height = 1/2 x 5 x 20 = 50 m
```

### Worked Example 5.2

A car accelerates from rest. Velocity at t = 0 is 0 m/s, at t = 6s is 24 m/s. v-t graph is a straight line.

(a) Find the acceleration.
(b) Find the displacement.

**Solution:**
(a) Acceleration = slope = change in v / change in t = 24/6 = **4 m/s²**

(b) Displacement = area = 1/2 x base x height = 1/2 x 6 x 24 = **72 m**

### Worked Example 5.3

```
Velocity (m/s)
20 |    *-------*
   |   /         \
   |  /           \
10 | /             
   |/               
 0 +-------------------> Time (s)
   0   4   7   9
```

(a) Acceleration in first 4s
(b) Acceleration from 4s to 7s
(c) Total displacement

**Solution:**
(a) slope from 0 to 4s = 20/4 = **5 m/s²**
(b) slope from 4s to 7s = 0 (flat) = **0 m/s²**
(c) Total area:
   - Triangle (0 to 4s): 1/2 x 4 x 20 = 40 m
   - Rectangle (4s to 7s): 3 x 20 = 60 m
   - Triangle (7s to 9s): 1/2 x 2 x 20 = 20 m
   - Total = **120 m**

---

## 5.5 Acceleration-Time Graphs

The **acceleration-time graph** (a-t graph) shows how acceleration changes over time.

### Flat horizontal line — Constant acceleration

```
Acceleration (m/s^2)
 4 |--------------------
   |
 0 +--------------------> Time (s)
```

This is the most common case (e.g., free fall, constant force problems).

### KEY RULE: Area under a-t graph = change in velocity

```
Area = a x delta_t = delta_v
```

### Connection Between All Three Graphs

```
x-t GRAPH              v-t GRAPH             a-t GRAPH

Slope gives v    -->   Slope gives a    -->   (nothing above)

(nothing below)  <--   Area gives dx   <--   Area gives dv
```

---

## 5.6 Comparing Motion Types on Graphs

| Motion Type      | x-t graph       | v-t graph      | a-t graph       |
|-----------------|-----------------|----------------|-----------------|
| At rest          | Flat line        | Flat at v=0    | Flat at a=0     |
| Constant +vel.   | Straight up      | Flat above 0   | Flat at a=0     |
| Constant +accel. | Curve up         | Straight up    | Flat above 0    |
| Decelerating     | Curve (flatter)  | Straight down  | Flat below 0    |
| Free fall (down) | Curve down       | Straight down  | Flat at -g      |

### Visual Comparison: Ball thrown upward

```
x-t GRAPH             v-t GRAPH            a-t GRAPH

x                     v                    a
 |    *               |  *                  |
 |   * *              | * \                 0|-------------------
 |  *   *             |/   \               -|
 | *     *            0     \            -10|-------------------
 |*       *           |      \               +---> t
 +----> t             +-------*-> t
                      
Parabola              Straight line         Flat line at -g
ball rises/falls      crossing zero         gravity constant
                      at top
```

---

## 5.7 Finding Gradient (Slope)

### Method: Rise Over Run

```
         rise    change_y    y2 - y1
slope = ------ = -------- = ---------
         run     change_x    x2 - x1
```

### Step-by-Step Method

1. Choose two points on the line that are **far apart** (for accuracy)
2. Draw a right-angle triangle between them
3. Read off the vertical change and horizontal change
4. Divide: slope = change_y / change_x
5. Include units: slope units = y-axis units / x-axis units

```
y (m/s)
30 |              B(6, 30)
   |           /
   |          /  change_y = 30-10 = 20
20 |        /
   |       /
10 | A(2, 10)
   |   change_x = 6-2 = 4
 0 +--------------------> x (s)
   
slope = 20/4 = 5 m/s^2
```

### For a Curved Line: Tangent Method

On a curved graph, the slope at any point equals the slope of the **tangent line** drawn at that point (a straight line that just touches the curve there).

```
x (m)
   |          /  <- tangent line at point P
   |        /
   |      /  <- slope of tangent = velocity at P
   |    P
   |   /
   |  /  <- the curve
   | /
   +--------> t (s)
```

Draw the tangent by eye, then calculate its gradient using the method above.

---

## 5.8 Finding Area Under a Graph

The area under a graph tells you the **accumulated quantity** — what builds up over time.

| Graph Type | Area Represents |
|-----------|-----------------|
| v-t graph | displacement (m) |
| a-t graph | change in velocity (m/s) |
| F-t graph | impulse (N·s) |
| P-t graph | energy (J) |

### Breaking Irregular Shapes into Simple Ones

```
v (m/s)
     *---------*
    /|          |
   / |          |\
  /  |          | \
 /   |          |  \
*----|----------|----|---> t(s)
0    2          7    9

Area = triangle + rectangle + triangle
     = 1/2(2)(20) + (5)(20) + 1/2(2)(20)
     = 20 + 100 + 20
     = 140 m
```

### Worked Example 5.4

A force F = 15 N acts for 4 seconds, then drops to 5 N for another 3 seconds.

```
F (N)
15 |*-----------*
   |            |
 5 |            *----------*
   |
 0 +--------------------> t(s)
   0    4        7
```

- Area from 0 to 4s: 15 x 4 = 60 N·s
- Area from 4s to 7s: 5 x 3 = 15 N·s
- Total impulse = **75 N·s**

---

## 5.9 Non-Linear Graphs

### Parabola: y = kx² (y proportional to x²)

Appears in: displacement vs. time during acceleration, kinetic energy vs. velocity.

```
y
   |           *
   |         *
   |        *
   |      *
   |    *
   |  *
   |*
   +---------------> x
```

### Hyperbola: y = k/x (y proportional to 1/x)

Appears in: Boyle's Law (pressure vs. volume), gravitational force vs. distance.

```
y
   |*
   | *
   |  *
   |    *
   |      *
   |           *
   |                  *
   +---------------> x
```

### Exponential decay: y = y0 x e^(-kt)

Appears in: radioactive decay, capacitor discharging, cooling.

```
y
   |*
   | *
   |  **
   |    **
   |      ***
   |         *****
   |               *********
   +------------------------------> t
```

### Sine wave: y = A sin(wt)

Appears in: oscillations, AC current, sound waves.

```
y
   |  *     *     *
  A|*   * *   * *   *
   |      *     *
   +----------------------> t
  -A
```

---

## 5.10 Linearizing Data

If you suspect a relationship like y = kx², you cannot directly find k from a y-vs-x graph (it is a curve). The solution is to **linearize** the data — transform it so you get a straight line.

### The Method

1. Identify the expected relationship (e.g., y = kx²)
2. Instead of plotting y vs x, plot y vs x²
3. The result should be a straight line with gradient k

### Example: Period of a Pendulum

The period T of a pendulum obeys: T = 2π√(L/g), which rearranges to:

T² = (4π²/g) × L

So if you plot **T² vs L**, you get a straight line with gradient = 4π²/g.

```
Before linearizing:         After linearizing:
T vs L                      T^2 vs L

T|                           T^2|           /
 |         curve               |          /
 |       /                     |        /  slope = 4pi^2/g
 |     /                       |      /
 +-------> L                   +-------> L
```

From the gradient, you can calculate g.

### When to Linearize

| Relationship | Plot | Expected Gradient |
|-------------|------|-------------------|
| y = kx²     | y vs x² | k |
| y = k/x     | y vs 1/x | k |
| y = k√x     | y vs √x | k |
| y = ke^(ax) | ln(y) vs x | a |

---

## Summary

- **Graphs** show the shape and nature of physical relationships visually
- A **position-time graph**: slope = velocity; flat = rest, straight = constant velocity, curved = acceleration
- A **velocity-time graph**: slope = acceleration; area under curve = displacement
- An **acceleration-time graph**: area = change in velocity; flat = constant acceleration
- The **gradient** (slope) of a straight line = rise/run = change_y / change_x; includes units
- For **curved lines**, the instantaneous slope is found by drawing a tangent
- The **area under a graph** = accumulated quantity (displacement under v-t, impulse under F-t)
- Common curve shapes: parabola (y ∝ x²), hyperbola (y ∝ 1/x), exponential decay, sine wave
- **Linearizing**: if y = kx², plot y vs x² to get a straight line with gradient k

---

## Key Equations

```
Gradient of straight line:
  slope = change_y / change_x = (y2-y1) / (x2-x1)

From x-t graph:
  velocity = slope = change_x / change_t

From v-t graph:
  acceleration = slope = change_v / change_t
  displacement = area under graph

From a-t graph:
  change in velocity = area under graph

Linearizing:
  If y = kx^2  ->  plot y vs x^2,  slope = k
  If y = k/x   ->  plot y vs 1/x,  slope = k
```
