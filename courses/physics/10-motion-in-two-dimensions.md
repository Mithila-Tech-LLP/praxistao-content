# Chapter 10: Motion in Two Dimensions

> **"The universe doesn't restrict movement to a single line. Real motion happens in two and three dimensions — but the secret is that you can always break it into simpler pieces."**

---

## Table of Contents

- [Why One Dimension Is Not Enough](#why-one-dimension-is-not-enough)
- [Setting Up a Coordinate System](#setting-up-a-coordinate-system)
- [The Independence Principle](#the-independence-principle)
- [Displacement in Two Dimensions](#displacement-in-two-dimensions)
- [Breaking Velocity into Components](#breaking-velocity-into-components)
- [Adding Vectors Using Components](#adding-vectors-using-components)
- [Relative Velocity in One Dimension First](#relative-velocity-in-one-dimension-first)
- [Relative Velocity in Two Dimensions](#relative-velocity-in-two-dimensions)
- [River Crossing Problems](#river-crossing-problems)
- [Wind and Aircraft Problems](#wind-and-aircraft-problems)
- [Worked Example 1 — Resultant Displacement from Two Journeys](#worked-example-1--resultant-displacement-from-two-journeys)
- [Worked Example 2 — River Crossing](#worked-example-2--river-crossing)
- [Worked Example 3 — Aircraft Heading Correction](#worked-example-3--aircraft-heading-correction)
- [Choosing a Convenient Coordinate System](#choosing-a-convenient-coordinate-system)
- [Common Mistakes to Avoid](#common-mistakes-to-avoid)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## Why One Dimension Is Not Enough

In the last few chapters, we studied motion in one dimension — things moving in a straight line, either forward or backward. That was a great starting point. You learned about displacement, velocity, and acceleration, all measured along a single number line.

But look around you. How does the real world actually move?

- A bird flies through the sky, moving forward AND upward at the same time
- A car drives around a bend, constantly changing direction
- You throw a ball and it curves through the air in a beautiful arc
- A river flows east while a boat tries to cross it heading north
- An airplane flies in wind that pushes it sideways

None of these can be described by a single number. They all happen in **two dimensions** — they involve motion in two perpendicular directions at once.

The good news is this: two-dimensional motion is not twice as complicated. It is just one-dimensional motion applied *twice* — once in each direction. The key insight, which we will come back to again and again, is that the two directions are completely independent of each other.

This chapter builds the tools you need to handle any 2D problem confidently.

---

## Setting Up a Coordinate System

Before you can describe 2D motion mathematically, you need to set up a **coordinate system** — a grid that gives every point in space a unique address.

The standard setup uses two perpendicular lines called **axes**:

- The **x-axis** runs horizontally (left-right)
- The **y-axis** runs vertically (up-down)

Where the two axes cross is called the **origin**, labelled O, and it has coordinates (0, 0).

```
         y
         |
         |
    +y   |   
         |
         |
---------O-----------  x
         |
    -y   |
         |
         
 (left = -x, right = +x, up = +y, down = -y)
```

Any point in the plane can be described as **(x, y)** — how far right and how far up from the origin.

### Choosing Positive Directions

By convention:
- Right is positive x, left is negative x
- Up is positive y, down is negative y

But you can choose any positive direction you like — as long as you stay consistent throughout the problem. Some problems are easier if you put the y-axis pointing horizontally and the x-axis pointing at an angle. We will see this later.

### An Example Grid

Imagine a city laid out on a grid. You start at the town center (the origin). You walk:
- 3 blocks east → your x position is +3
- 2 blocks north → your y position is +2
- Your position is the point **(3, 2)**

```
         N (y)
         |
    3    |    * (3,2)
    2    |   /
    1    |  /
         | /
------O--+--------  E (x)
         |  1  2  3
         |
```

The point (3, 2) is 3 blocks east and 2 blocks north of where you started.

---

## The Independence Principle

This is the most important idea in this chapter. Write it down, underline it, tattoo it on your brain:

> **The horizontal (x) motion and the vertical (y) motion are completely independent of each other.**

What does "independent" mean here? It means:
- What happens in the x-direction has NO effect on what happens in the y-direction
- You can analyze x and y completely separately, then combine your results at the end

This seems almost too simple to be important. But it is revolutionary. It means that any 2D problem can be split into two 1D problems, each of which you already know how to solve.

### A Demonstration

Imagine you are standing at the edge of a table. You have two identical balls. At exactly the same moment:
- Ball A: you drop straight down (no horizontal motion)
- Ball B: you roll off the edge horizontally at 2 m/s

Which ball hits the floor first?

Most people guess Ball A, because it goes straight down. But the answer is: **they hit the floor at exactly the same time**.

Why? Because both balls fall vertically with exactly the same acceleration due to gravity. The fact that Ball B is also moving horizontally makes NO difference to how fast it falls. The horizontal motion is independent of the vertical motion.

```
Table edge
----------*---------->
          | \
          |   \   Ball B path (parabola)
          |     \
          |       \
          |         *
          * Ball A lands    Ball B lands
          (same time)       (farther away)
```

This experiment was actually performed by Galileo. The result surprised people then, and it still surprises people today. But it works — every time — because of the independence principle.

---

## Displacement in Two Dimensions

**Displacement** in 2D is not just a number — it is a **vector**, meaning it has both a size and a direction.

If you start at point A and end at point B, your displacement is the straight-line arrow from A to B.

### Components of Displacement

Suppose you walk 5 m east and then 3 m north. Your total displacement has two **components**:
- x-component: Δx = +5 m (east)
- y-component: Δy = +3 m (north)

```
         N
         |
    3m   +-----* End point
    north|     
         |     
         O-----------  E
         Start  5m east
```

The arrow from the start to the end point is the displacement vector. It goes diagonally. But we can always describe it using its x and y components.

### The Magnitude of Displacement (Distance)

The magnitude (length) of the displacement vector is found using the **Pythagorean theorem**:

    magnitude = sqrt(Δx² + Δy²)

For our example:
    magnitude = sqrt(5² + 3²)
              = sqrt(25 + 9)
              = sqrt(34)
              ≈ 5.83 m

So you ended up 5.83 m from where you started, in a direction northeast.

### The Direction of Displacement

The angle θ that the displacement makes with the x-axis (east direction) is:

    tan(θ) = Δy / Δx = 3 / 5 = 0.6

    θ = arctan(0.6) ≈ 31°

So the displacement is 5.83 m at 31° north of east.

```
         |
    3    +......* 
         |     /
    2    |    /  5.83 m
         |   /  (magnitude)
    1    |  /
         | / θ = 31°
         O-----------
           1  2  3  4  5
```

---

## Breaking Velocity into Components

**Velocity** in 2D is also a vector. An object moving at speed v in a direction θ from the x-axis has:

- x-component of velocity: vₓ = v × cos(θ)
- y-component of velocity: vᵧ = v × sin(θ)

This is just the same as breaking any 2D vector into its components using trigonometry.

### Why This Works

Think of the velocity vector as the hypotenuse of a right triangle. The horizontal leg is vₓ and the vertical leg is vᵧ.

```
         |
         |  vᵧ = v sin θ
         | /
         |/  ← angle θ
         +--------
              vₓ = v cos θ
              
    The velocity vector v is the hypotenuse.
    Its components are vₓ (horizontal) and vᵧ (vertical).
```

### Example

A bird flies at 10 m/s at an angle of 40° above the horizontal.

- vₓ = 10 × cos(40°) = 10 × 0.766 = 7.66 m/s (horizontal)
- vᵧ = 10 × sin(40°) = 10 × 0.643 = 6.43 m/s (vertical, upward)

The bird moves 7.66 m/s forward and 6.43 m/s upward, simultaneously and independently.

### Rebuilding the Vector from Components

If you know vₓ and vᵧ, you can find the full velocity:

- Magnitude: v = sqrt(vₓ² + vᵧ²)
- Direction: θ = arctan(vᵧ / vₓ)

---

## Adding Vectors Using Components

Often you need to add two or more vectors together. The cleanest way is the **component method**:

1. Break each vector into x and y components
2. Add all x-components together: Rₓ = Aₓ + Bₓ + ...
3. Add all y-components together: Rᵧ = Aᵧ + Bᵧ + ...
4. Find the magnitude: R = sqrt(Rₓ² + Rᵧ²)
5. Find the direction: θ = arctan(Rᵧ / Rₓ)

### Example: Two Forces

A box is pulled by two ropes:
- Rope 1: 20 N at 0° (pure east)
- Rope 2: 15 N at 90° (pure north)

Step 1: Components of Rope 1:
    Aₓ = 20 × cos(0°) = 20 × 1 = 20 N
    Aᵧ = 20 × sin(0°) = 20 × 0 = 0 N

Step 2: Components of Rope 2:
    Bₓ = 15 × cos(90°) = 15 × 0 = 0 N
    Bᵧ = 15 × sin(90°) = 15 × 1 = 15 N

Step 3: Sum the components:
    Rₓ = 20 + 0 = 20 N
    Rᵧ = 0 + 15 = 15 N

Step 4: Magnitude of resultant:
    R = sqrt(20² + 15²) = sqrt(400 + 225) = sqrt(625) = 25 N

Step 5: Direction:
    θ = arctan(15/20) = arctan(0.75) ≈ 36.9° above east

```
         |
    15N  +.....* Resultant: 25 N at 36.9°
    Rope2|     /
         |    /
         |   /  25 N
         |  /
         | /θ
         O-----------
              20 N (Rope 1)
```

---

## Relative Velocity in One Dimension First

Before tackling 2D relative velocity, let us make sure the concept is clear in 1D.

**Relative velocity** is the velocity of one object as measured by an observer on another moving object.

### The Train and the Walker

Imagine a train moving east at 30 m/s. You are walking inside the train heading east at 2 m/s relative to the train.

What is your velocity relative to someone standing on the ground (a "stationary" observer)?

    v(you, ground) = v(you, train) + v(train, ground)
    v(you, ground) = 2 + 30 = 32 m/s east

Now suppose you walk WEST (backward) at 2 m/s relative to the train:

    v(you, ground) = -2 + 30 = 28 m/s east

The key formula for 1D relative velocity:

    v_AB = v_A - v_B

where v_AB means "velocity of A as measured by observer B."

```
Ground observer
|
|   Train moving at 30 m/s →
|   +---------------------------+
|   |  You walking at 2 m/s →  |
|   +---------------------------+
|
|   Your speed relative to ground = 30 + 2 = 32 m/s →
```

### Another Example

Two cars on a highway, both heading north:
- Car A: 25 m/s north
- Car B: 20 m/s north

Velocity of A relative to B:
    v_AB = v_A - v_B = 25 - 20 = 5 m/s north

From Car B's driver's perspective, Car A appears to slowly pull away at 5 m/s. This makes intuitive sense.

---

## Relative Velocity in Two Dimensions

Now we extend the idea to 2D. The formula is the same, but now the velocities are vectors with x and y components:

    v⃗_AB = v⃗_A - v⃗_B

Or equivalently:

    v⃗_A(ground) = v⃗_A(relative to medium) + v⃗_medium(ground)

The key step is to handle each component (x and y) separately.

---

## River Crossing Problems

**River crossing** is the classic 2D relative velocity problem. Let us build it carefully.

### The Setup

A river flows east (let us say this is the +x direction). A boat wants to cross the river. The river is 100 m wide (north to south, which is the +y direction).

The boat points directly north (straight across) and moves at 3 m/s relative to the water.
The river flows east at 4 m/s relative to the ground.

What actually happens?

```
North
  |
  |   Boat heads this way →
  |   but river pushes east ↓
  |
  |   Actual path of boat (diagonal)
  |  /
  | /
  |/ ← angle downstream
  *-----------------------
Start          River flows this way →
```

### The Independence Principle in Action

The boat's velocity relative to the ground has two independent components:

- North (y) component: 3 m/s (from the boat's engine, pointing straight north)
- East (x) component: 4 m/s (from the river current)

These act independently. The boat moves north AND gets carried east at the same time.

### Finding the Actual Velocity

The actual velocity vector of the boat relative to the ground:

    vₓ = 4 m/s (east, from current)
    vᵧ = 3 m/s (north, from boat's engine)

Magnitude:
    v = sqrt(4² + 3²) = sqrt(16 + 9) = sqrt(25) = 5 m/s

Direction (angle east of north):
    tan(θ) = vₓ / vᵧ = 4/3 = 1.333
    θ = arctan(1.333) ≈ 53° east of north

So the boat travels at 5 m/s, but its actual path is 53° east of north — it ends up downstream.

```
       N
       |
       |         * End point (downstream)
       |        /
       |       /  5 m/s (actual path)
       |      /
 100m  |     /
river  |    /  θ = 53°
width  |   /
       |  /
       | /
       |/
       *-----------  E
       Start       River flows east at 4 m/s
```

### How Long to Cross?

The time to cross the river depends ONLY on the north (y) component of velocity:

    time = width / vᵧ = 100 / 3 ≈ 33.3 seconds

### How Far Downstream?

The east (x) drift depends on the time and the east velocity:

    drift = vₓ × time = 4 × 33.3 ≈ 133 m downstream

Even though the boat aimed straight across, it ends up 133 m downstream from where it intended to land.

---

## Wind and Aircraft Problems

The same logic applies to an airplane flying in wind. The airplane has a velocity relative to the air, and the air (wind) has a velocity relative to the ground. The airplane's actual velocity over the ground is the vector sum.

### The Scenario

An airplane needs to fly due north at an airspeed of 200 m/s. But there is a wind blowing due east at 50 m/s.

If the pilot simply points the plane north, the plane will be carried east by the wind — it will NOT travel due north over the ground.

To travel due north, the pilot must angle the nose slightly west, to compensate for the eastward wind.

### Finding the Required Heading

Let θ be the angle west of north the pilot must point the aircraft.

We need the ground velocity to point due north (no east component).

    vₓ (ground) = vₓ (plane in air) + vₓ (wind) = 0

The wind pushes east at +50 m/s, so the plane must have a westward component of -50 m/s through the air:

    vₓ (plane in air) = -v_air × sin(θ) = -50
    v_air × sin(θ) = 50
    200 × sin(θ) = 50
    sin(θ) = 50/200 = 0.25
    θ = arcsin(0.25) ≈ 14.5° west of north

### Finding the Ground Speed

The northward component of the plane's actual velocity:

    vᵧ (ground) = v_air × cos(θ) = 200 × cos(14.5°) = 200 × 0.968 = 193.6 m/s

Even though the plane's airspeed is 200 m/s, it only makes 193.6 m/s of progress northward because some of its effort is spent fighting the crosswind.

```
         N
         |
         |  Plane's actual path
         |  (due north, 193.6 m/s)
         |
         |
Plane    | /  ← Plane points slightly west of north
nose     |/      (14.5° west)
heads    *
west     |\
of       |  \ Wind carries east
north    |    → 50 m/s
         |
```

---

## Worked Example 1 — Resultant Displacement from Two Journeys

**Problem:** A hiker starts at the origin. They walk 8 km at 60° north of east (60° above the +x axis), then walk 6 km due south (−y direction). Find the total displacement: magnitude and direction.

### Step 1: Break each journey into components

**Journey 1:** 8 km at 60° north of east

    x₁ = 8 × cos(60°) = 8 × 0.5 = 4 km (east)
    y₁ = 8 × sin(60°) = 8 × 0.866 = 6.93 km (north)

**Journey 2:** 6 km due south

    x₂ = 0 km (no east-west movement)
    y₂ = -6 km (south means negative y)

### Step 2: Sum the components

    Total x: Rₓ = x₁ + x₂ = 4 + 0 = 4 km (east)
    Total y: Rᵧ = y₁ + y₂ = 6.93 + (-6) = 0.93 km (north)

### Step 3: Find the magnitude

    R = sqrt(Rₓ² + Rᵧ²)
      = sqrt(4² + 0.93²)
      = sqrt(16 + 0.865)
      = sqrt(16.865)
      ≈ 4.11 km

### Step 4: Find the direction

    tan(θ) = Rᵧ / Rₓ = 0.93 / 4 = 0.2325
    θ = arctan(0.2325) ≈ 13.1° north of east

### Answer

The total displacement is approximately **4.11 km at 13.1° north of east**.

```
         N
         |
     7   |  * (end of journey 1)
         | /
     6   |/   Journey 1: 8 km at 60°
         |
     1   +--* Final position (4, 0.93)
         |
         O---------  E
              4
              
    Journey 2 goes from (4, 6.93) straight south 6 km to (4, 0.93).
```

---

## Worked Example 2 — River Crossing

**Problem:** A river is 80 m wide. It flows east at 2 m/s. A swimmer can swim at 1.5 m/s relative to the water. The swimmer heads straight north (directly across). Find:

(a) The time to cross
(b) How far downstream the swimmer ends up
(c) The actual speed of the swimmer relative to the ground

### Given Information

- River width: d = 80 m (north direction)
- River current: vᵣ = 2 m/s (east)
- Swimmer speed in water: v_s = 1.5 m/s (north, straight across)

### Part (a): Time to Cross

The swimmer's northward velocity is 1.5 m/s. The river current is east and does NOT affect north-south travel.

    time = distance / northward velocity
    time = 80 / 1.5 ≈ 53.3 seconds

### Part (b): Downstream Drift

During those 53.3 seconds, the current carries the swimmer east:

    drift = vᵣ × time = 2 × 53.3 ≈ 106.7 m downstream

So the swimmer ends up 106.7 m east of where they intended to land.

### Part (c): Actual Speed

    v_actual = sqrt(v_s² + vᵣ²)
             = sqrt(1.5² + 2²)
             = sqrt(2.25 + 4)
             = sqrt(6.25)
             = 2.5 m/s

The swimmer actually travels at 2.5 m/s (even though they can only swim at 1.5 m/s) because the current contributes to their speed.

```
       N
       |
       |                        * Swimmer ends up HERE
       |                       /
 80m   |                      /  2.5 m/s (actual path)
       |                     /
       |                    /
       |                   /  arctan(2/1.5) ≈ 53.1°
       |                  /    east of north
       |_________________/____________________________  E
   Start              106.7 m east (downstream drift)
```

---

## Worked Example 3 — Aircraft Heading Correction

**Problem:** A pilot wants to fly from City A to City B, which is located exactly 500 km due east (in the +x direction). The airplane's airspeed is 250 m/s. A wind is blowing from south to north at 60 m/s (in the +y direction).

If the pilot points directly east, the plane will drift north of City B. What heading should the pilot use to fly directly east and reach City B? What is the actual ground speed?

### Understanding the Problem

We need the plane's velocity over the ground to point purely east (no north-south component).

The wind pushes north at +60 m/s. So the plane must have a southward component of -60 m/s through the air to cancel the wind.

### Finding the Heading

Let θ be the angle south of east the pilot must point the nose.

The plane must have vᵧ = -60 m/s (southward) through the air to cancel the wind:

    v_air × sin(θ) = 60
    250 × sin(θ) = 60
    sin(θ) = 60/250 = 0.24
    θ = arcsin(0.24) ≈ 13.9° south of east

### Finding Ground Speed

The eastward component of the plane's airspeed:

    vₓ (air) = v_air × cos(θ) = 250 × cos(13.9°) = 250 × 0.971 = 242.7 m/s

The wind has no eastward component, so the ground speed is:

    v_ground = 242.7 m/s (due east)

Even though the airspeed is 250 m/s, the plane only makes 242.7 m/s of eastward progress because it angled slightly south.

### Time to Reach City B

    time = 500,000 m / 242.7 m/s ≈ 2060 seconds ≈ 34.3 minutes

```
         N (y)
         |
         |   Wind blows north
         |   ↑ 60 m/s
         |
         |
    A    |  Pilot aims slightly
    +----|--south of east-----> B
         |  to stay on path
         |
         | ← Nose points 13.9° below east
         |
```

---

## Choosing a Convenient Coordinate System

One of the skills you develop in 2D physics is choosing the best coordinate system for a given problem. The math is easier when you align the axes with the main directions of motion.

### Guidelines

1. **Align one axis with the main motion**: If a ball rolls along a flat surface, put x along the surface.

2. **Align one axis with a known force**: For inclined planes (covered later), it helps to tilt the axes with the slope.

3. **Put the origin at the starting point**: Then the initial position is (0, 0), which simplifies all displacement equations.

4. **Keep the positive y direction pointing up in most problems**: Gravity is then always -g (negative y), which is the standard convention.

5. **Think about symmetry**: If the problem has left-right symmetry, put the y-axis in the center.

### Example: Choosing Origin and Axes for River Problem

For a river crossing problem:
- Put the origin at the starting bank where the boat enters the water
- Put the y-axis pointing across the river (north)
- Put the x-axis pointing downstream (east)

Now:
- The boat's velocity is purely in the +y direction
- The river's velocity is purely in the +x direction
- The two components are naturally separated

This makes the equations much cleaner to write and solve.

---

## Common Mistakes to Avoid

### Mistake 1: Forgetting Independence

Students sometimes think that a larger horizontal velocity will make a falling object take longer to hit the ground. It does not. The vertical fall time is set by gravity and height alone. Horizontal speed has no effect on vertical fall time.

### Mistake 2: Adding Magnitudes Instead of Components

You cannot add vector magnitudes directly unless the vectors point in the same direction.

WRONG: A 3 m/s boat in a 4 m/s current does NOT give 3 + 4 = 7 m/s.
RIGHT: Find components first, then use the Pythagorean theorem: sqrt(3² + 4²) = 5 m/s.

### Mistake 3: Mixing Up Angles

When finding components:
- The component ALONG the angle uses cosine: v cos(θ)
- The component PERPENDICULAR to the reference uses sine: v sin(θ)

Always draw a diagram to verify which component is which.

### Mistake 4: Forgetting Signs

East is +x, west is -x, north is +y, south is -y. If a current pushes west, its x-component is negative. Always include the sign when adding components.

### Mistake 5: Using the Wrong Velocity for Time Calculations

In river crossing problems, the time to cross depends on the velocity component perpendicular to the river banks — NOT the total actual speed. And the downstream drift depends on the current speed times the crossing time.

---

## Summary

- **Two-dimensional motion** requires two numbers (x and y coordinates) to describe position, velocity, and displacement.

- The **coordinate system** uses x (horizontal) and y (vertical) axes with an origin at (0, 0). Right and up are positive; left and down are negative.

- The **independence principle** states that x and y motions are completely independent of each other. You can analyze each separately and combine the results at the end.

- Any 2D **displacement vector** can be broken into x and y components. The magnitude is sqrt(Δx² + Δy²) and the direction is arctan(Δy / Δx).

- Any 2D **velocity vector** at angle θ from the x-axis has components: vₓ = v cos(θ) and vᵧ = v sin(θ).

- To **add vectors**, break each into components, add the x-components, add the y-components, then find the magnitude and direction of the result.

- **Relative velocity in 1D**: v_AB = v_A − v_B. Walking on a moving train: your ground speed is the sum of your walking speed and the train's speed.

- **Relative velocity in 2D**: Apply the same formula as a vector equation, handling x and y components separately.

- In **river crossing problems**: the crossing time depends only on the velocity component pointing across the river; the downstream drift depends on the current speed times the crossing time.

- In **aircraft problems**: to fly a desired course, the pilot aims the nose to cancel any crosswind; the ground speed is found from the vector sum of airspeed and wind velocity.

- **Choosing a good coordinate system** (origin at the start, axes aligned with motion directions) simplifies the math greatly.

---

## Key Equations

| Quantity | Formula |
|---|---|
| Magnitude of 2D displacement | R = sqrt(Δx² + Δy²) |
| Direction of 2D displacement | θ = arctan(Δy / Δx) |
| x-component of velocity | vₓ = v × cos(θ) |
| y-component of velocity | vᵧ = v × sin(θ) |
| Rebuild magnitude from components | v = sqrt(vₓ² + vᵧ²) |
| Relative velocity (1D) | v_AB = v_A − v_B |
| Vector addition (x) | Rₓ = Aₓ + Bₓ |
| Vector addition (y) | Rᵧ = Aᵧ + Bᵧ |
| River crossing time | t = width / v_perpendicular |
| River downstream drift | drift = v_current × t |

---

*End of Chapter 10*
