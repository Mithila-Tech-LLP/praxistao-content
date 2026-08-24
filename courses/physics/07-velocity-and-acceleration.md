# Chapter 07: Velocity and Acceleration

> **"Nothing happens until something moves."**
> — Albert Einstein

---

## Table of Contents

1. [Why Velocity is More Than Just Speed](#1-why-velocity-is-more-than-just-speed)
2. [Velocity — The Full Picture](#2-velocity--the-full-picture)
3. [Average Velocity vs Instantaneous Velocity](#3-average-velocity-vs-instantaneous-velocity)
4. [What is Acceleration?](#4-what-is-acceleration)
5. [The Units of Acceleration — m/s²](#5-the-units-of-acceleration--ms)
6. [Positive and Negative Acceleration](#6-positive-and-negative-acceleration)
7. [The Sign Convention — Choosing a Direction](#7-the-sign-convention--choosing-a-direction)
8. [Uniform vs Non-Uniform Acceleration](#8-uniform-vs-non-uniform-acceleration)
9. [The Velocity-Time Graph](#9-the-velocity-time-graph)
10. [Acceleration Due to Gravity](#10-acceleration-due-to-gravity)
11. [The "Feel" of Acceleration](#11-the-feel-of-acceleration)
12. [Worked Examples](#12-worked-examples)
13. [Summary](#13-summary)
14. [Key Equations](#14-key-equations)

---

## 1. Why Velocity is More Than Just Speed

In the last chapter, we talked about **distance** and **displacement** — and discovered that they are different things. Distance counts every step you take, while displacement only cares about where you ended up relative to where you started.

Speed and velocity follow the same idea. They sound like the same thing in everyday English, but in physics, they mean something very different. This distinction is not just academic — it has real consequences. It is the difference between a GPS navigation system that actually gets you to your destination and one that just counts the total kilometres you drove.

Let's start with a simple story to make this clear.

---

### The Pizza Delivery Story

Imagine you work for a pizza delivery shop. One Friday evening, your manager gives you two assignments:

**Assignment A:** Deliver a pizza to a customer 3 km north of the shop.

**Assignment B:** Deliver a pizza to a customer who lives in a maze-like apartment complex. You have to ride 1 km north, then 1 km east, then 1 km south, then 1 km west — looping around — before finally arriving at the customer's door, which turns out to be directly next to the shop.

Both deliveries take 10 minutes.

Now ask yourself:

- What was your **speed** for each delivery?
- What was your **velocity** for each delivery?

For Assignment A:
- Distance travelled = 3 km
- Displacement = 3 km north
- Speed = 3 km / 10 min = 0.3 km/min
- Velocity = 0.3 km/min **northward**

For Assignment B:
- Distance travelled = 1 + 1 + 1 + 1 = 4 km
- Displacement = 0 km (you ended up where you started)
- Speed = 4 km / 10 min = 0.4 km/min
- Velocity = 0 km/min (zero! you went nowhere, net)

This is the core idea. **Velocity** tells you not just how fast you moved, but in which direction — and it is calculated using displacement, not distance.

---

## 2. Velocity — The Full Picture

### The Definition

**Velocity** is the rate of change of displacement. In simpler terms: it tells you how quickly your position is changing, and in which direction.

The formula is:

```
v = displacement / time

or

v = d / t
```

where:
- `v` = velocity (metres per second, m/s)
- `d` = displacement (metres, m) — with direction included
- `t` = time taken (seconds, s)

Because velocity includes a direction, it is what physicists call a **vector quantity**. Speed, on the other hand, has no direction — it is a **scalar quantity**.

---

### Velocity vs Speed — Side by Side

| Property | Speed | Velocity |
|---|---|---|
| What it measures | How fast | How fast + which direction |
| Based on | Distance | Displacement |
| Type | Scalar (just a number) | Vector (number + direction) |
| Example | 60 km/h | 60 km/h northward |
| Can it be negative? | No | Yes (if moving "backward") |
| Can it be zero while moving? | No | Yes (if net displacement = 0) |

---

### ASCII Diagram: Speed vs Velocity

```
SPEED (scalar)                VELOCITY (vector)

    Car                            Car
     |                              |
     |  Odometer reads:             |  GPS reads:
     |  60 km/h                     |  60 km/h ---> (eastward)
     |                              |
  (just a number)             (number + arrow)
```

The arrow is everything. Without knowing which way you are going, you cannot predict where you will be.

---

### Real-World Velocity Examples

Let's get a feel for common velocities:

| Situation | Approximate Velocity |
|---|---|
| Person walking | 1.4 m/s forward |
| Cyclist on flat road | 5–8 m/s forward |
| Car on highway | 25–33 m/s forward |
| Commercial aeroplane | 250 m/s eastward (approx) |
| Sound in air | 343 m/s in any direction |
| Earth orbiting the Sun | ~30,000 m/s (30 km/s) |

---

### Worked Example 2.1 — Calculating Velocity

**Problem:** A cyclist starts at a park bench and rides 800 metres due east in 100 seconds. What is the cyclist's velocity?

**Solution:**

Step 1: Identify what we know.
- Displacement = 800 m east
- Time = 100 s

Step 2: Write the formula.
```
v = displacement / time
```

Step 3: Substitute the values.
```
v = 800 m / 100 s
v = 8 m/s
```

Step 4: Include the direction.
```
v = 8 m/s east
```

**Answer:** The cyclist's velocity is **8 m/s eastward**.

---

## 3. Average Velocity vs Instantaneous Velocity

This section introduces one of the most important distinctions in physics. You have probably noticed that when you are in a car, the speedometer needle is always moving slightly — you are never travelling at exactly one speed the whole time.

This gives rise to two different types of velocity.

---

### Average Velocity

**Average velocity** is the total displacement divided by the total time taken. It gives you a single number that summarises the whole journey.

```
v_avg = total displacement / total time
```

Think of it like your report card grade for the whole semester — it does not tell you how you did on any single test, but it gives a general picture.

---

### Instantaneous Velocity

**Instantaneous velocity** is your velocity at one specific moment in time. It is what the speedometer (combined with your direction) is showing you right now.

Imagine taking a photograph of a moving car. The photograph freezes one instant. The velocity of the car at that frozen instant — that is instantaneous velocity.

In everyday driving, your instantaneous velocity is constantly changing as you speed up, slow down, and turn corners.

---

### ASCII Diagram: Average vs Instantaneous

```
Velocity
(m/s)
  |
30|          * (peak — instantaneous velocity here)
  |        /   \
20|      /       \
  |    /           \-------
10|  /
  |/
  +-------------------------> Time (s)
  0   2   4   6   8   10

Average velocity = area under curve / total time
                = not the peak, not the low point
                = somewhere in the middle
```

---

### Worked Example 3.1 — Average Velocity

**Problem:** Riya drives her car from home to the market, which is 6 km north. She then drives 2 km south to pick up her friend. The whole trip takes 30 minutes (0.5 hours). What is her average velocity?

**Solution:**

Step 1: Calculate the total displacement.
- She goes 6 km north, then 2 km south.
- Net displacement = 6 km north − 2 km south = 4 km north

Step 2: Note the total time.
- Total time = 30 minutes = 0.5 hours

Step 3: Apply the formula.
```
v_avg = total displacement / total time
v_avg = 4 km / 0.5 hours
v_avg = 8 km/h northward
```

**Answer:** Riya's average velocity is **8 km/h northward**, even though she drove a total distance of 8 km.

(Note: if we had used speed, it would be 8 km / 0.5 h = 16 km/h — a completely different answer!)

---

### When Average and Instantaneous are the Same

There is one special case where average velocity equals instantaneous velocity: when the object moves at a perfectly **constant** velocity the entire time. If a train moves north at exactly 20 m/s for 5 minutes without any change, every snapshot of that journey shows the same velocity — 20 m/s north.

This special case is called **uniform motion**, and it is the foundation for understanding the more interesting case: **acceleration**.

---

## 4. What is Acceleration?

Here is a question: can something be moving and still be accelerating?

Yes — absolutely. In fact, you experience acceleration every single day: every time a bus starts moving from a stop, every time a car brakes, every time a lift starts going up.

---

### The Core Idea

Think about this carefully. When a car is sitting at a red light, its velocity is 0 m/s. The light turns green, and 5 seconds later, the car is moving at 20 m/s forward.

The car's velocity **changed**. It went from 0 m/s to 20 m/s in 5 seconds.

**Acceleration** is the rate at which velocity changes. It measures how quickly velocity is changing — whether it is speeding up, slowing down, or even changing direction.

The formal definition:

> **Acceleration** is the change in velocity divided by the time taken for that change.

---

### The Formula

```
a = (v - u) / t
```

where:
- `a` = acceleration (metres per second squared, m/s²)
- `v` = final velocity (metres per second, m/s)
- `u` = initial velocity (metres per second, m/s) — the letter `u` is standard for "initial velocity" in physics
- `t` = time taken for the change (seconds, s)

The quantity `(v - u)` is the **change in velocity** — sometimes written as `Δv` (delta-v, where the Greek letter delta means "change in").

So you will also see:

```
a = Δv / t
```

Both mean exactly the same thing.

---

### Why `u` for Initial Velocity?

This is a fair question. The letter `u` for initial velocity comes from old British physics notation, and it has stuck around for over a century. You will see it in textbooks all over the world. Just remember:

- `u` = initial (starting) velocity — where you begin
- `v` = final velocity — where you end up

An easy memory trick: **U** comes before **V** in the alphabet, just like initial comes before final.

---

### ASCII Diagram: Acceleration Concept

```
TIME:     0s        1s        2s        3s        4s        5s
          |         |         |         |         |         |
          
VELOCITY: 0 m/s --> 4 m/s --> 8 m/s --> 12 m/s-> 16 m/s-> 20 m/s

          [+4 m/s]  [+4 m/s]  [+4 m/s]  [+4 m/s]  [+4 m/s]
          
          Each second, velocity increases by the SAME amount.
          This is UNIFORM ACCELERATION = 4 m/s²
```

---

## 5. The Units of Acceleration — m/s²

Let's take a moment to understand why the units of acceleration look so strange: **metres per second squared (m/s²)**.

It seems weird. What does "metres per second squared" even mean physically?

Let's build it from the formula:

```
a = change in velocity / time

  = (m/s) / s

  = m/s²
```

So m/s² literally means "metres per second, per second."

Here is the everyday way to read it: if an object has an acceleration of 3 m/s², it means its velocity **increases by 3 metres per second, every second**.

---

### Reading m/s² in Plain English

| Acceleration | What it means |
|---|---|
| 1 m/s² | Every second, velocity goes up by 1 m/s |
| 5 m/s² | Every second, velocity goes up by 5 m/s |
| 10 m/s² | Every second, velocity goes up by 10 m/s |
| −3 m/s² | Every second, velocity goes DOWN by 3 m/s (slowing) |

---

### Worked Example 5.1 — Calculating Acceleration

**Problem:** A car is stopped at a traffic light. When the light turns green, the car reaches a velocity of 15 m/s after 5 seconds. What is the car's acceleration?

**Solution:**

Step 1: Identify the known values.
- Initial velocity (u) = 0 m/s (car was stopped)
- Final velocity (v) = 15 m/s
- Time (t) = 5 s

Step 2: Write the formula.
```
a = (v - u) / t
```

Step 3: Substitute values.
```
a = (15 - 0) / 5
a = 15 / 5
a = 3 m/s²
```

**Answer:** The car accelerates at **3 m/s²**. This means every second, the car's velocity increases by 3 m/s.

Let us verify:
- After 1s: 0 + 3 = 3 m/s ✓
- After 2s: 3 + 3 = 6 m/s ✓
- After 3s: 6 + 3 = 9 m/s ✓
- After 4s: 9 + 3 = 12 m/s ✓
- After 5s: 12 + 3 = 15 m/s ✓

---

## 6. Positive and Negative Acceleration

So far we have only seen acceleration that makes things go faster. But acceleration can also make things go slower. And this is where the sign — positive (+) or negative (−) — becomes crucial.

---

### Positive Acceleration

When an object **speeds up** in its direction of travel, its acceleration is **positive**.

Example: A rocket launching upward. It starts at rest and gets faster and faster as it rises. Its acceleration is positive (upward, same direction as its velocity).

```
ROCKET LAUNCH

    ^  ^  ^         (velocity arrows getting longer = speeding up)
    |  |  |
  [ ROCKET ]
     Fire!
     
Velocity: 0 --> 100 --> 250 --> 500 m/s
Acceleration: POSITIVE (velocity increasing)
```

---

### Negative Acceleration (Deceleration)

When an object **slows down**, its acceleration is **negative**. This is also called **deceleration**.

**Deceleration** just means negative acceleration — the velocity is decreasing. Many people use "deceleration" in everyday speech, but physicists prefer to say "negative acceleration" because it is more precise.

Example: A car braking hard before a red light.

```
CAR BRAKING

--> --> --> -->    (car moving right, but arrows shrinking = slowing down)
  [ CAR ]

Velocity: 30 m/s --> 20 m/s --> 10 m/s --> 0 m/s
Acceleration: NEGATIVE (velocity decreasing)
```

---

### Worked Example 6.1 — Negative Acceleration (Braking)

**Problem:** A bus is travelling at 24 m/s on a highway. The driver sees a child on the road and brakes hard, coming to a complete stop in 6 seconds. What is the bus's acceleration?

**Solution:**

Step 1: Identify known values.
- Initial velocity (u) = 24 m/s (forward)
- Final velocity (v) = 0 m/s (stopped)
- Time (t) = 6 s

Step 2: Apply the formula.
```
a = (v - u) / t
a = (0 - 24) / 6
a = -24 / 6
a = -4 m/s²
```

**Answer:** The bus's acceleration is **−4 m/s²**. The negative sign tells us the bus is slowing down. Every second, its velocity decreases by 4 m/s.

Verification:
- After 1s: 24 − 4 = 20 m/s ✓
- After 2s: 20 − 4 = 16 m/s ✓
- After 3s: 16 − 4 = 12 m/s ✓
- After 4s: 12 − 4 = 8 m/s ✓
- After 5s: 8 − 4 = 4 m/s ✓
- After 6s: 4 − 4 = 0 m/s ✓ (stopped)

---

### An Important Nuance: Negative Acceleration Does Not Always Mean Slowing Down

This is a subtlety worth noting, even for beginners. Negative acceleration means the velocity is changing in the negative direction. Whether the object is speeding up or slowing down depends on the sign of the velocity too.

| Velocity | Acceleration | What Happens |
|---|---|---|
| Positive (+) | Positive (+) | Object speeds up in positive direction |
| Positive (+) | Negative (−) | Object slows down (decelerates) |
| Negative (−) | Negative (−) | Object speeds up in negative direction |
| Negative (−) | Positive (+) | Object slows down (decelerates) |

Do not worry too much about the bottom two rows for now — we will revisit them when we set up the sign convention properly.

---

## 7. The Sign Convention — Choosing a Direction

Physics involves objects moving in different directions. To handle this mathematically, we need to agree on what counts as "positive" and what counts as "negative."

**Sign convention** is the process of defining a positive direction for a problem.

---

### How to Set Up a Sign Convention

The good news is: **you get to choose**. There is no universal rule. You pick the positive direction at the start of the problem, and then you stick with it throughout.

Common conventions:
- For horizontal motion: rightward = positive, leftward = negative
- For vertical motion: upward = positive, downward = negative

```
HORIZONTAL SIGN CONVENTION

      Negative              Positive
<--- (−) ---[  OBJECT  ]--- (+) --->
      Left                  Right

VERTICAL SIGN CONVENTION

            (+)
             ^
             |  Upward = positive
             |
          [Object]
             |
             |  Downward = negative
             v
            (−)
```

---

### Why This Matters

Once you set your convention, you must be consistent. If you say rightward is positive, then:
- A car moving right has **positive velocity**
- A car moving left has **negative velocity**
- A car speeding up to the right has **positive acceleration**
- A car slowing down (while moving right) has **negative acceleration**

---

### Worked Example 7.1 — Using Sign Convention

**Problem:** A ball rolls to the right at 10 m/s. A gust of wind pushes back on it, causing it to decelerate at 2 m/s². How long does it take for the ball to stop?

**Set up:** Define rightward as positive.
- Initial velocity (u) = +10 m/s (moving right)
- Acceleration (a) = −2 m/s² (opposing rightward motion)
- Final velocity (v) = 0 m/s (stopped)

**Formula:**
```
a = (v - u) / t

Rearrange to find t:
t = (v - u) / a
t = (0 - 10) / (-2)
t = (-10) / (-2)
t = 5 s
```

**Answer:** The ball stops after **5 seconds**.

---

### Worked Example 7.2 — Elevator Going Up

**Problem:** An elevator starts from rest at the ground floor and accelerates upward at 1.5 m/s² for 4 seconds. What is the elevator's velocity after 4 seconds?

**Set up:** Define upward as positive.
- Initial velocity (u) = 0 m/s (starts at rest)
- Acceleration (a) = +1.5 m/s² (moving up, so positive)
- Time (t) = 4 s

**Formula:**
```
a = (v - u) / t

Rearrange to find v:
v = u + (a × t)
v = 0 + (1.5 × 4)
v = 6 m/s
```

**Answer:** After 4 seconds, the elevator is moving upward at **6 m/s**.

---

## 8. Uniform vs Non-Uniform Acceleration

Not all acceleration is the same throughout a journey. Physics distinguishes between two important types.

---

### Uniform Acceleration (Constant Acceleration)

**Uniform acceleration** means the acceleration stays the same throughout the motion. The velocity changes by the same amount every second.

This is the easiest type to work with mathematically, and it is what the standard kinematics formulas (which you will learn in the next few chapters) are designed for.

```
UNIFORM ACCELERATION EXAMPLE

Time (s):    0    1    2    3    4    5
Velocity:    0    3    6    9   12   15  (m/s)
Change:        +3   +3   +3   +3   +3

Constant change each second = UNIFORM acceleration = 3 m/s²
```

Real-world examples of approximately uniform acceleration:
- A ball rolling down a gentle, straight slope
- A skydiver in the first second of freefall
- A car with cruise-control acceleration engaged

---

### Non-Uniform Acceleration

**Non-uniform acceleration** means the acceleration is changing over time. The velocity does not increase (or decrease) by the same amount every second.

This is the more common real-world situation. Think of how you drive: sometimes you press the gas pedal gently, sometimes firmly. The acceleration is constantly varying.

```
NON-UNIFORM ACCELERATION EXAMPLE

Time (s):    0    1    2    3    4    5
Velocity:    0    2    5    9   12   15  (m/s)
Change:        +2   +3   +4   +3   +3

Changing amounts each second = NON-UNIFORM acceleration
```

---

### ASCII Comparison Diagram

```
UNIFORM ACCELERATION              NON-UNIFORM ACCELERATION
(constant, equal jumps)           (variable jumps)

Velocity                          Velocity
 |                                 |
 |      /                          |         ___
 |    /                            |       /
 |  /                              |     /
 | /                               |   _/
 |/                                |  /
 +-----------> Time                +-----------> Time
 
 Straight line                     Curved line
```

This is a very important visual concept. We will return to it when we discuss velocity-time graphs in Section 9.

---

### Comparison Table

| Feature | Uniform Acceleration | Non-Uniform Acceleration |
|---|---|---|
| Acceleration value | Constant (same every second) | Varies (different each second) |
| Velocity change | Equal intervals each second | Unequal intervals |
| v-t graph shape | Straight line | Curved line |
| Real-world example | Ball on smooth ramp | Car in city traffic |
| Maths difficulty | Simpler (use standard formulas) | More complex (needs calculus) |

---

## 9. The Velocity-Time Graph

One of the most powerful tools in kinematics is the **velocity-time graph** (often called a v-t graph). Once you understand it, you can "read" the motion of any object just by looking at its shape.

---

### Setting Up the Graph

- **Horizontal axis (x-axis):** Time in seconds
- **Vertical axis (y-axis):** Velocity in m/s

Every point on the graph represents the velocity of the object at that particular instant.

---

### What the Slope Tells You

Here is the key insight of the v-t graph:

> **The slope (steepness) of the line on a v-t graph equals the acceleration.**

This makes perfect sense from the formula:

```
a = (v - u) / t = change in velocity / change in time = slope of the v-t graph
```

Let's break down the different shapes a v-t graph can take.

---

### Case 1: Horizontal Line — Zero Acceleration

```
Velocity (m/s)
  |
20|-----------------------------  <- constant velocity
  |
  |
  |
  +------------------------------> Time (s)
  0   2   4   6   8   10

Slope = 0
Acceleration = 0 m/s²
Object is moving at constant velocity.
```

---

### Case 2: Straight Line Sloping Upward — Constant Positive Acceleration

```
Velocity (m/s)
  |                     *
20|                   *
  |                 *
15|               *
  |             *
10|           *
  |         *
 5|       *
  |     *
  |   *
  | *
  +------------------------------> Time (s)
  0   2   4   6   8   10

Slope is positive and constant.
Acceleration = positive constant value.
Object is uniformly accelerating.
```

---

### Case 3: Straight Line Sloping Downward — Constant Negative Acceleration

```
Velocity (m/s)
  |
20| *
  |   *
15|     *
  |       *
10|         *
  |           *
 5|             *
  |               *
  |                 *
  +------------------------------> Time (s)
  0   2   4   6   8   10

Slope is negative and constant.
Acceleration = negative constant value.
Object is uniformly decelerating.
```

---

### Case 4: Curved Line — Non-Uniform Acceleration

```
Velocity (m/s)
  |                     *
20|                  *
  |               *
15|           *
  |        *
10|     *
  |   *
 5|  *
  | *
  +------------------------------> Time (s)
  0   2   4   6   8   10

The curve is getting less steep over time.
Acceleration is decreasing (still positive, but reducing).
This is non-uniform acceleration.
```

---

### How to Calculate Acceleration From a v-t Graph

To find acceleration from a straight-line v-t graph, use the slope formula:

```
acceleration = rise / run
             = (change in velocity) / (change in time)
             = (v2 - v1) / (t2 - t1)
```

This is exactly the same as our acceleration formula, just written in graph terms.

---

### Worked Example 9.1 — Reading a v-t Graph

**Problem:** A car's velocity is measured every 2 seconds. The data is shown below. Find the acceleration.

```
Time (s):    0    2    4    6    8
Velocity:    5   11   17   23   29  (m/s)
```

**Solution:**

Step 1: Check if velocity changes by equal amounts. 
- 0→2s: change = 6 m/s
- 2→4s: change = 6 m/s
- 4→6s: change = 6 m/s
- 6→8s: change = 6 m/s
Equal changes → this is uniform acceleration.

Step 2: Calculate acceleration using any two points.
Let's use (t=0, v=5) and (t=8, v=29):

```
a = (v2 - v1) / (t2 - t1)
a = (29 - 5) / (8 - 0)
a = 24 / 8
a = 3 m/s²
```

**Answer:** The car has a constant acceleration of **3 m/s²**.

---

### Worked Example 9.2 — Finding Acceleration From Two Graph Points

**Problem:** A motorbike has a velocity of 20 m/s at t = 2s, and a velocity of 8 m/s at t = 7s. Calculate its acceleration.

**Solution:**

```
a = (v - u) / t
a = (8 - 20) / (7 - 2)
a = -12 / 5
a = -2.4 m/s²
```

**Answer:** The motorbike's acceleration is **−2.4 m/s²**. The negative sign tells us it is slowing down.

---

### A Complete v-t Graph Journey

Let's trace a car's entire trip and read every section:

```
Velocity (m/s)
  |
25|          ***********
  |        *             *
20|      *               *
  |    *                   *
15|  *                     *
  |*                         *
10|                           *-----------
  |                           
 5|
  |
  +--+--+--+--+--+--+--+--+--+---> Time (s)
  0  2  4  6  8  10 12 14 16 18 20

Section A (0→4s):   Straight line up      → Uniform positive acceleration
Section B (4→12s):  Flat at top           → Zero acceleration (constant velocity 25 m/s)

  Wait — actually let me draw this more simply:

SECTION    TIME     VELOCITY       SLOPE        MOTION
  A        0→5s     0 to 20 m/s   Positive     Speeding up
  B        5→10s    20 m/s (flat) Zero         Constant speed
  C        10→14s   20 to 0 m/s   Negative     Slowing down (braking)
  D        14→18s   0 m/s (flat)  Zero         Stopped
```

```
Velocity (m/s)
  |
20|     ___________
  |   /             \
15|  /               \
  | /                 \
10|/                   \
  |                     \
 5|                      \________
  |
  +--+--+--+--+--+--+--+--+--+--> Time (s)
  0  2  4  6  8  10 12 14 16 18

Section A (0-4s):   Accelerating (slope up)
Section B (4-10s):  Constant velocity (flat line)
Section C (10-14s): Decelerating (slope down)
Section D (14-18s): Stopped (flat at zero)
```

Reading this graph, a physicist can reconstruct the entire story of this car's journey without anyone telling them. That is the power of the v-t graph.

---

## 10. Acceleration Due to Gravity

Now we come to perhaps the most famous acceleration in all of physics — the acceleration caused by Earth's gravity.

---

### The Basic Idea

Drop anything — a pen, a book, a ball. It falls. And as it falls, it goes faster and faster. It does not fall at a constant speed. Gravity is accelerating it downward.

This is called the **acceleration due to gravity**, and it has a special symbol: **g**.

On Earth's surface, this value is approximately:

```
g ≈ 9.8 m/s²  (downward)
```

Sometimes, to make calculations easier, we round it to 10 m/s².

What this means practically: every second an object is in freefall, its downward velocity increases by 9.8 m/s.

---

### Freefall Velocity Table

If you drop an object from a great height (ignoring air resistance), here is how fast it would be going:

| Time After Release | Velocity (downward) |
|---|---|
| 0 seconds | 0 m/s |
| 1 second | 9.8 m/s |
| 2 seconds | 19.6 m/s |
| 3 seconds | 29.4 m/s |
| 4 seconds | 39.2 m/s |
| 5 seconds | 49.0 m/s |

---

### ASCII Diagram: Object in Freefall

```
t = 0s   [ Ball released from rest ]
           |
           v  (0 m/s)

t = 1s   
           |
           vvv  (9.8 m/s)

t = 2s   
           |
           vvvvvv  (19.6 m/s)

t = 3s   
           |
           vvvvvvvvv  (29.4 m/s)
           
Each second, the arrow gets longer — the ball is going faster.
The rate of increase is constant: 9.8 m/s per second.
```

---

### Sign Convention for Gravity

If you choose downward as positive:
- `g = +9.8 m/s²`

If you choose upward as positive (more common for problems involving throwing things up):
- `g = −9.8 m/s²` (gravity acts opposite to the positive direction)

You will always see both conventions used. The key is to be consistent within a single problem.

---

### A Quick Freefall Calculation

**Problem:** You drop a stone from a bridge. After 3 seconds, how fast is the stone falling? (Ignore air resistance. Use g = 9.8 m/s².)

**Solution:**
- u = 0 m/s (dropped from rest)
- a = 9.8 m/s² (downward)
- t = 3 s

```
v = u + (a × t)
v = 0 + (9.8 × 3)
v = 29.4 m/s
```

The stone is falling at **29.4 m/s** after 3 seconds.

(This is about 106 km/h — the speed of a car on a motorway. Gravity is powerful.)

---

### Why g is the Same for All Objects

Galileo proved this over 400 years ago by dropping objects of different masses from the Tower of Pisa (or so the legend goes). A heavy cannonball and a light wooden ball hit the ground at the same time. They had the same acceleration — g — regardless of their mass.

This seems counterintuitive, but it is one of the most thoroughly tested facts in all of physics. The reason involves a beautiful balance between gravitational force and inertia — which you will learn about in the chapters on Newton's Laws.

---

## 11. The "Feel" of Acceleration

Numbers on paper are one thing. Let us connect acceleration to things you feel in your body.

---

### 1 m/s² — Gentle

A typical family car on a smooth road, driven gently. You feel a soft push back into your seat. If you have a coffee in the cupholder, it might slosh gently but will not spill.

In 10 seconds: velocity goes from 0 to 10 m/s (about 36 km/h). Perfectly comfortable.

---

### 3–4 m/s² — Noticeable

A sportier car when the driver gives it moderate throttle. You feel clearly pushed back into your seat. The acceleration used in Worked Example 5.1 (a car going from 0 to 15 m/s in 5 seconds) was 3 m/s².

This is also roughly the deceleration of a normal car braking at a moderate pace. When you feel yourself pressed against your seatbelt as a car brakes — that is around 3–4 m/s².

---

### 8–10 m/s² — Intense

A high-performance sports car at maximum acceleration, or a subway/metro train during emergency braking. This is approaching the force of gravity itself (9.8 m/s²). You are pressed firmly back into your seat, or forward against your seatbelt.

---

### g = 9.8 m/s² — Gravity Itself

When you are in freefall — skydiving, for example — you experience exactly 1g of acceleration. Interestingly, you do not feel it because your whole body is accelerating at the same rate (this is why astronauts in orbit feel "weightless" even though gravity is still pulling them).

---

### Rollercoaster and Rocket Comparisons

| Vehicle / Situation | Approximate Acceleration |
|---|---|
| Walking pace increase | ~0.5 m/s² |
| Typical city car | 1–3 m/s² |
| Sports car 0-100 km/h | 5–8 m/s² |
| Gravity (freefall) | 9.8 m/s² |
| Fighter jet at takeoff | ~30 m/s² |
| Space Shuttle launch | ~30 m/s² (at peak) |
| Crash (no airbag) | 100–500 m/s² |

---

### What Acceleration Feels Like in Your Body

When you accelerate forward, you feel pushed back. This is because your seat is pushing you forward, but your body (and especially the fluid in your inner ear) tries to stay behind due to inertia.

When you decelerate (brake), you feel thrown forward. Your belt holds you back, but your body wants to keep going.

When you go around a corner at speed, you feel pushed to the outside of the turn. This is also a form of acceleration — changing direction means changing velocity, which is acceleration.

```
FEELING ACCELERATION

Forward acceleration:     Body feels pushed BACKWARD
                          [ seat ]---> pushes you
                          [ you  ] <--- feels this
                          
Braking (deceleration):   Body feels pushed FORWARD
                          [ belt ]<--- holds you back
                          [ you  ] ---> wants to keep going
                          
Turning left:             Body feels pushed to the RIGHT
                          (toward outside of the turn)
```

---

## 12. Worked Examples

Let us now bring everything together with a set of worked examples that cover the key concepts from this chapter.

---

### Worked Example 12.1 — Rocket Launch

**Problem:** A model rocket is launched straight up from rest. Its engine provides constant acceleration. After 8 seconds, the rocket is moving upward at 120 m/s. 

(a) What is the rocket's acceleration?  
(b) What is the rocket's velocity after 5 seconds?

**Solution:**

Set up: upward = positive.

Given:
- u = 0 m/s (starts at rest)
- v = 120 m/s (after 8 seconds)
- t = 8 s

**(a) Find acceleration:**

```
a = (v - u) / t
a = (120 - 0) / 8
a = 120 / 8
a = 15 m/s²
```

The rocket's acceleration is **15 m/s²** upward.

**(b) Find velocity after 5 seconds:**

Rearrange the formula to find v:
```
a = (v - u) / t
v = u + (a × t)
v = 0 + (15 × 5)
v = 75 m/s
```

After 5 seconds, the rocket is moving at **75 m/s upward**.

```
ROCKET VELOCITY OVER TIME

Velocity
(m/s)         *  (120 m/s at 8s)
120|         *
   |        *
   |       *
75 |.....*
   |    *
   |   *
   |  *
   | *
   |*
   +----+----+----+----> Time (s)
   0    2    4   5 6    8
```

---

### Worked Example 12.2 — The Braking Car

**Problem:** A car is travelling on a motorway at 30 m/s. The driver sees an accident ahead and brakes hard. The car comes to a complete stop in 10 seconds.

(a) What is the deceleration?
(b) What is the velocity after 6 seconds?
(c) Sketch the velocity-time graph.

**Solution:**

Set up: forward (initial direction of travel) = positive.

Given:
- u = 30 m/s
- v = 0 m/s
- t = 10 s

**(a) Find acceleration:**

```
a = (v - u) / t
a = (0 - 30) / 10
a = -30 / 10
a = -3 m/s²
```

The acceleration is **−3 m/s²** (a deceleration of 3 m/s²).

**(b) Find velocity after 6 seconds:**

```
v = u + (a × t)
v = 30 + (-3 × 6)
v = 30 - 18
v = 12 m/s
```

After 6 seconds, the car is still moving at **12 m/s** forward.

**(c) Sketch:**

```
Velocity (m/s)
  |
30| *
  |  *
25|   *
  |    *
20|     *
  |      *
15|       *
  |        *
10|         *
  |          *
 5|           *
  |            *
  |             *
  +--+--+--+--+--+--+---> Time (s)
  0  2  4  6  8  10

Straight line from (0, 30) to (10, 0).
Slope = -3 m/s² = constant deceleration.
```

---

### Worked Example 12.3 — The Elevator

**Problem:** An elevator starts from rest and accelerates upward at 2 m/s² for 4 seconds, then travels at constant velocity for 8 seconds, then decelerates at 2 m/s² until it stops.

(a) What is the maximum velocity?
(b) How long does the deceleration phase last?
(c) Sketch the full velocity-time graph.

**Solution:**

**Phase 1 — Acceleration:**
- u = 0 m/s, a = 2 m/s², t = 4 s
```
v = u + (a × t)
v = 0 + (2 × 4)
v = 8 m/s
```
Maximum velocity = **8 m/s**

**Phase 2 — Constant velocity:**
- v = 8 m/s for 8 seconds. Nothing to calculate.

**Phase 3 — Deceleration:**
- u = 8 m/s (starts phase at 8 m/s)
- v = 0 m/s (stops)
- a = −2 m/s²

```
t = (v - u) / a
t = (0 - 8) / (-2)
t = -8 / -2
t = 4 s
```

Deceleration phase lasts **4 seconds**.

**Full journey duration:** 4 + 8 + 4 = **16 seconds**

**The v-t Graph:**

```
Velocity (m/s)
  |
 8|     ***************
  |   *                 *
 6|  *                   *
  | *                     *
 4|*                       *
  |                         *
 2|                          *
  |                           *
  +--+--+--+--+--+--+--+--+--+-> Time (s)
  0  2  4  6  8  10 12 14 16

Phase 1 (0-4s):   Slope up (+2 m/s²)
Phase 2 (4-12s):  Flat (0 m/s²)
Phase 3 (12-16s): Slope down (-2 m/s²)
```

---

### Worked Example 12.4 — Finding Initial Velocity

**Problem:** A ball is rolling across a floor. Its final velocity is 3 m/s. It has been decelerating at 0.5 m/s² for 6 seconds due to friction. What was its initial velocity?

**Solution:**

Given:
- v = 3 m/s
- a = −0.5 m/s² (decelerating due to friction)
- t = 6 s

Rearrange to find u:
```
a = (v - u) / t

Multiply both sides by t:
a × t = v - u

Rearrange:
u = v - (a × t)
u = 3 - (-0.5 × 6)
u = 3 - (-3)
u = 3 + 3
u = 6 m/s
```

**Answer:** The ball's initial velocity was **6 m/s**.

---

### Worked Example 12.5 — A Falling Object

**Problem:** A stone is dropped from a cliff. Using g = 9.8 m/s² and ignoring air resistance:

(a) What is the stone's velocity after 4 seconds?
(b) If the stone hits the ground with a velocity of 49 m/s, how long did it fall?

**Solution:**

Set up: downward = positive (so g = +9.8 m/s²).

Initial velocity u = 0 (dropped from rest).

**(a) Velocity after 4 seconds:**

```
v = u + (a × t)
v = 0 + (9.8 × 4)
v = 39.2 m/s
```

The stone is falling at **39.2 m/s** after 4 seconds.

**(b) Time to reach 49 m/s:**

```
a = (v - u) / t
t = (v - u) / a
t = (49 - 0) / 9.8
t = 49 / 9.8
t = 5 s
```

The stone fell for **5 seconds** before hitting the ground.

---

### Worked Example 12.6 — Interpreting a Complicated v-t Graph

**Problem:** A cyclist's velocity is recorded in the table below. Describe the motion and calculate the acceleration in each phase.

| Time (s) | Velocity (m/s) |
|---|---|
| 0 | 0 |
| 5 | 10 |
| 10 | 10 |
| 15 | 10 |
| 20 | 4 |

**Solution:**

Phase 1 (t = 0 to 5s):
```
a = (v - u) / t = (10 - 0) / 5 = 2 m/s²
```
Cyclist is accelerating at 2 m/s² (speeding up).

Phase 2 (t = 5 to 15s):
```
a = (v - u) / t = (10 - 10) / 10 = 0 m/s²
```
Cyclist is moving at constant velocity (on a flat road at steady pace).

Phase 3 (t = 15 to 20s):
```
a = (v - u) / t = (4 - 10) / 5 = -6/5 = -1.2 m/s²
```
Cyclist is decelerating at 1.2 m/s² (perhaps going uphill).

```
Velocity (m/s)
  |
10|     *-----------*
  |   *               *
 8|  *                 *
  | *                   *
 6|*                     *
  |                       *
 4|                        *
  |
  +--+--+--+--+--+--+--+--+-> Time (s)
  0  2  4  6  8  10 12 14 16 18 20

Phase 1: Slope up (a = +2 m/s²)
Phase 2: Flat (a = 0 m/s²)
Phase 3: Slope down (a = −1.2 m/s²)
```

---

## 13. Summary

Let us collect everything we have learned in this chapter into a clear set of takeaways.

### Key Concepts

- **Velocity** is speed with direction. It is calculated using displacement, not distance: `v = displacement / time`. It is a **vector quantity**.

- **Speed** is a scalar (no direction). **Velocity** is a vector (has direction). They are different things, even though we often use the words interchangeably in everyday life.

- **Average velocity** is total displacement divided by total time. It summarises a whole journey with one number.

- **Instantaneous velocity** is the velocity at one specific moment in time — what the speedometer reads right now.

- **Acceleration** is the rate of change of velocity: `a = (v - u) / t`. It tells you how quickly velocity is changing.

- The standard letters used are: **u** for initial velocity, **v** for final velocity, **a** for acceleration, **t** for time.

- Acceleration has units of **m/s²** (metres per second squared), which means "metres per second, per second."

- **Positive acceleration** means the object is gaining velocity in the positive direction.

- **Negative acceleration** (also called **deceleration**) means the velocity is decreasing — the object is slowing down (when moving in the positive direction).

- You must set up a **sign convention** at the start of every problem — choose which direction is positive and stick to it.

- **Uniform acceleration** means the acceleration is constant — velocity changes by equal amounts every second. The v-t graph is a straight line.

- **Non-uniform acceleration** means the acceleration varies — velocity changes by different amounts at different times. The v-t graph is curved.

- On a **velocity-time graph**, the **slope** (steepness) of the line equals the acceleration. Upward slope = positive acceleration. Downward slope = negative acceleration. Flat line = zero acceleration.

- The **acceleration due to gravity** near Earth's surface is **g ≈ 9.8 m/s²** downward. All objects fall with this same acceleration (ignoring air resistance), regardless of mass.

- A typical car accelerating from a stop uses about 1–4 m/s². A falling object gains speed at 9.8 m/s². A rocket at launch can reach 30 m/s² or more.

---

## 14. Key Equations

```
VELOCITY
─────────────────────────────────────────────
v = d / t

where:
  v = velocity (m/s)
  d = displacement (m)
  t = time (s)


ACCELERATION
─────────────────────────────────────────────
a = (v - u) / t

where:
  a = acceleration (m/s²)
  v = final velocity (m/s)
  u = initial velocity (m/s)
  t = time (s)


REARRANGEMENTS OF THE ACCELERATION FORMULA
─────────────────────────────────────────────
To find final velocity:
  v = u + (a × t)

To find initial velocity:
  u = v - (a × t)

To find time:
  t = (v - u) / a


ACCELERATION DUE TO GRAVITY
─────────────────────────────────────────────
g ≈ 9.8 m/s²  (downward, near Earth's surface)
g ≈ 10 m/s²   (rounded value, for quick estimates)


SLOPE OF A VELOCITY-TIME GRAPH
─────────────────────────────────────────────
a = (v2 - v1) / (t2 - t1)

(This is just the standard rise/run slope formula applied to a v-t graph.)


AVERAGE VELOCITY
─────────────────────────────────────────────
v_avg = total displacement / total time
```

---

*End of Chapter 07. In the next chapter, we will use these ideas to derive the famous **equations of motion** — five powerful formulas that let you solve almost any kinematics problem involving uniform acceleration.*
