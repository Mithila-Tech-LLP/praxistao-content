# Chapter 08: Equations of Motion (SUVAT)

> **"An object at rest stays at rest and an object in motion stays in motion — but if you want to know *where* it ends up, you need the equations of motion."**

---

## Table of Contents

1. [What Is This Chapter About?](#1-what-is-this-chapter-about)
2. [The Five SUVAT Variables](#2-the-five-suvat-variables)
3. [The Meaning of "Uniform Acceleration"](#3-the-meaning-of-uniform-acceleration)
4. [Deriving the Five Equations Step by Step](#4-deriving-the-five-equations-step-by-step)
5. [The Known/Unknown Method — Choosing the Right Equation](#5-the-knownunknown-method--choosing-the-right-equation)
6. [Worked Example 1: Car Accelerating from Rest](#6-worked-example-1-car-accelerating-from-rest)
7. [Worked Example 2: Ball Thrown Straight Up](#7-worked-example-2-ball-thrown-straight-up)
8. [Worked Example 3: Braking Distance of a Car](#8-worked-example-3-braking-distance-of-a-car)
9. [Worked Example 4: Stone Dropped from a Cliff](#9-worked-example-4-stone-dropped-from-a-cliff)
10. [Deriving SUVAT from Velocity-Time Graphs](#10-deriving-suvat-from-velocity-time-graphs)
11. [Common Mistakes to Avoid](#11-common-mistakes-to-avoid)
12. [Quick-Reference Cheat Sheet](#12-quick-reference-cheat-sheet)
13. [Practice Problems](#13-practice-problems)
14. [Summary](#14-summary)
15. [Key Equations](#15-key-equations)

---

## 1. What Is This Chapter About?

Imagine you are standing at the edge of a highway watching a car accelerate away from a red light. You notice:

- The car starts from a complete stop.
- It moves faster and faster.
- After a few seconds it is cruising at full speed.

Now ask yourself: *How far did the car travel in those first 5 seconds? What was its speed when it reached the 100-metre mark?*

These are not vague questions. They have exact, calculable answers — and the tool you use to calculate them is called the **SUVAT equations** (also called the **equations of uniform acceleration** or the **kinematic equations**).

SUVAT is simply a set of five formulas. Each formula connects five quantities that describe motion: displacement, initial velocity, final velocity, acceleration, and time. If you know three of those five quantities, you can always find the other two.

This chapter is one of the most practical in all of introductory physics. You will use these equations constantly — for cars on roads, balls in the air, rockets launching, cyclists braking. Master this chapter and a huge chunk of classical mechanics becomes straightforward.

Let us start from scratch.

---

## 2. The Five SUVAT Variables

The name **SUVAT** comes from the five letters that represent the five variables:

| Letter | Quantity | Meaning | Unit |
|--------|----------|---------|------|
| **s** | Displacement | How far the object moved, and in which direction | metres (m) |
| **u** | Initial velocity | How fast it was going at the START | metres per second (m/s) |
| **v** | Final velocity | How fast it is going at the END | metres per second (m/s) |
| **a** | Acceleration | How quickly the velocity is changing | metres per second squared (m/s²) |
| **t** | Time | How long the motion lasted | seconds (s) |

Let us understand each one carefully, because confusion about these variables is the number-one source of mistakes.

---

### 2.1 s — Displacement

You already know the word "distance." **Displacement** is similar but slightly different.

- **Distance** is how much ground you covered, full stop.
- **Displacement** is how far you are from your starting point, and in which direction.

```
Starting point                             Ending point
     |                                          |
     A ----------------------------------------> B
     
     <------------- displacement = 50 m -------->
     (you moved 50 m to the right)
```

If you walked 30 m to the right, then 10 m back to the left, your distance is 40 m but your displacement is only 20 m to the right.

In the SUVAT equations, **s** is always the displacement from the starting point to the ending point, measured in a straight line along the direction of motion.

**Sign convention:** We choose a positive direction at the start of any problem. Anything in that direction is positive (+). Anything opposite is negative (-).

```
Positive direction chosen: RIGHT
         
         +s (move right)       -s (move left)
              --->                   <---
```

---

### 2.2 u — Initial Velocity

**u** is the velocity of the object at the very beginning of the time period you are studying — the moment your "clock starts."

- If the car is stationary when the light turns green, u = 0 m/s.
- If the car is already doing 20 m/s when you start timing, u = 20 m/s.
- If you throw something downward, u is negative (if you chose upward as positive).

The subscript "initial" is what matters. It is not the velocity at some random middle point — it is the velocity at t = 0.

---

### 2.3 v — Final Velocity

**v** is the velocity at the end of the time period you are studying — when your "clock stops."

- It might be when the car hits the brakes.
- It might be when the ball reaches its highest point.
- It might be after exactly 5 seconds of travel.

Notice that "final" does not mean the object has stopped. It just means the end of the time window you have chosen.

---

### 2.4 a — Acceleration

**Acceleration** is the rate at which velocity changes. In plain English: how quickly is the object speeding up or slowing down?

```
Time:    t=0        t=1s       t=2s       t=3s
Velocity: 0 m/s --> 4 m/s --> 8 m/s --> 12 m/s

Velocity increases by 4 m/s every second.
Therefore: a = 4 m/s²  (pronounced "4 metres per second squared")
```

Acceleration can be:
- **Positive**: object is speeding up in the positive direction (or slowing down in the negative direction).
- **Negative**: object is slowing down in the positive direction (or speeding up in the negative direction).
- **Zero**: velocity is not changing at all.

When a car brakes, its acceleration is negative (opposite to the direction of motion). Physicists sometimes call this **deceleration**, but it is still just acceleration with a negative sign.

---

### 2.5 t — Time

**t** is simply the duration of the motion — how many seconds have passed from the start to the end of the time window.

Always use seconds in SUVAT. If a problem gives you minutes, convert first.

---

### Putting It Together: A Visual Summary

```
     MOMENT 1                                    MOMENT 2
     (clock starts)                              (clock stops)
         |                                           |
         |                                           |
         v                                           v
  =============== u = initial speed ====>====>====>===============
  |   Object   |  a = acceleration         |   Object   |
  |   starts   |  s = total displacement   |   ends up  |
  ===============                          ===============
  
  Time elapsed = t seconds
  Speed at end = v (final velocity)
```

---

## 3. The Meaning of "Uniform Acceleration"

Here is a critical detail: **the SUVAT equations only work when acceleration is constant (uniform)**.

**Uniform acceleration** means the acceleration does not change during the motion. Every second, the velocity increases (or decreases) by exactly the same amount.

```
UNIFORM ACCELERATION (SUVAT applies)
Velocity:  0   4   8   12  16  20  m/s
Time:      0   1   2   3   4   5   s
           +4  +4  +4  +4  +4  (same increase each second)


NON-UNIFORM ACCELERATION (SUVAT does NOT apply directly)
Velocity:  0   2   7   9   18  22  m/s
Time:      0   1   2   3   4   5   s
           +2  +5  +2  +9  +4  (different each second)
```

Real-world situations where acceleration is approximately uniform:
- A car accelerating on a straight, flat road (roughly)
- A ball falling under gravity (exactly uniform, a = 9.8 m/s² downward)
- A ball rolling down a smooth ramp
- A rocket during a short burn (approximately)

Real-world situations where acceleration is NOT uniform:
- A car in stop-and-go traffic
- A feather falling through air (air resistance changes with speed)
- A ball at the end of a rubber band

For this chapter, we assume uniform acceleration throughout.

---

## 4. Deriving the Five Equations Step by Step

This is the heart of the chapter. We are going to build all five equations from two simple definitions. You do not need to memorise the derivations, but understanding them will help you remember and apply the equations correctly.

Our two starting definitions:

**Definition 1: Acceleration = change in velocity divided by time**

```
a = (v - u) / t
```

**Definition 2: Displacement = average velocity × time**

```
s = ((u + v) / 2) × t
```

The second definition works because when acceleration is uniform, the velocity increases at a steady rate. This means the average velocity is exactly halfway between the start and end velocities.

```
velocity
  ^
  |        /
v |-------/
  |      /
  |     /
  |    /
u |---/
  |  /
  | /
  |/
  +-----------> time
  0            t

The average is exactly (u + v)/2 because the graph is a straight line.
The shaded triangle plus rectangle = total area = s
```

Now let us build the five equations.

---

### 4.1 Equation 1: v = u + at

This comes directly from rearranging the definition of acceleration.

We know:
```
a = (v - u) / t
```

Multiply both sides by t:
```
a × t = v - u
```

Add u to both sides:
```
u + at = v
```

Write it the standard way:
```
v = u + at
```

**What it says in plain English:** Your final speed equals your starting speed plus however much speed you gained (or lost) due to acceleration over time t.

Think of it like a bank account:
- You start with u pounds.
- You earn `a` pounds per second for `t` seconds.
- You end up with `v = u + at` pounds.

---

### 4.2 Equation 2: s = ut + ½at²

We take Equation 1 (v = u + at) and substitute it into Definition 2.

From Definition 2:
```
s = ((u + v) / 2) × t
```

Replace v with (u + at) from Equation 1:
```
s = ((u + (u + at)) / 2) × t
s = ((2u + at) / 2) × t
s = (u + ½at) × t
s = ut + ½at²
```

So:
```
s = ut + ½at²
```

**What it says in plain English:** Your displacement is made up of two parts:
1. `ut` — how far you would have gone at your original speed with no acceleration.
2. `½at²` — the extra distance gained (or lost) because of acceleration.

```
s = [  ut  ] + [ ½at² ]
     ^              ^
     distance at    bonus distance
     constant u     due to acceleration
```

If you start from rest (u = 0), the equation simplifies to:
```
s = ½at²
```

This is a beautiful result: the distance fallen by a dropped object is proportional to the square of the time. Double the time, quadruple the distance.

---

### 4.3 Equation 3: v² = u² + 2as

This equation is useful when you do not know the time. We derive it by eliminating t.

From Equation 1:
```
v = u + at
```

Rearrange for t:
```
t = (v - u) / a
```

Substitute into Definition 2:
```
s = ((u + v) / 2) × t
s = ((u + v) / 2) × ((v - u) / a)
s = (u + v)(v - u) / (2a)
```

Notice that (u + v)(v - u) = v² - u² (difference of squares). So:
```
s = (v² - u²) / (2a)
```

Multiply both sides by 2a:
```
2as = v² - u²
```

Rearrange:
```
v² = u² + 2as
```

**What it says in plain English:** Your final speed squared equals your initial speed squared plus twice the acceleration times the distance. This equation is incredibly useful for braking problems where time is unknown.

---

### 4.4 Equation 4: s = ½(u + v)t

This one is simply Definition 2 written explicitly.

We know average velocity = (u + v)/2. And displacement = average velocity × time.

```
s = ((u + v) / 2) × t
s = ½(u + v)t
```

**What it says in plain English:** The displacement equals the average of the start and end velocities, multiplied by time. This equation is useful when you know both velocities and the time, but not the acceleration.

---

### 4.5 Equation 5: s = vt - ½at²

This is like Equation 2 but written in terms of v (final velocity) instead of u (initial velocity).

From Equation 1:
```
v = u + at
```

Rearrange for u:
```
u = v - at
```

Substitute into Equation 2:
```
s = ut + ½at²
s = (v - at)t + ½at²
s = vt - at² + ½at²
s = vt - ½at²
```

So:
```
s = vt - ½at²
```

**What it says in plain English:** If you know the final velocity, this equation calculates displacement directly. It is less commonly used than Equation 2, but it is very handy when u is unknown and v is given.

---

### The Complete Set

Here they are, all five together:

```
+-------+---------------------+----------------------------+
|  No.  |      Equation       |  Missing variable          |
+-------+---------------------+----------------------------+
|   1   |  v = u + at         |  s (displacement)          |
|   2   |  s = ut + ½at²      |  v (final velocity)        |
|   3   |  v² = u² + 2as      |  t (time)                  |
|   4   |  s = ½(u + v)t      |  a (acceleration)          |
|   5   |  s = vt - ½at²      |  u (initial velocity)      |
+-------+---------------------+----------------------------+
```

The "missing variable" column tells you which variable does NOT appear in that equation. This is the key to choosing the right equation.

---

## 5. The Known/Unknown Method — Choosing the Right Equation

Here is a systematic method that never fails. Follow these steps every time:

**Step 1:** Write out all five SUVAT variables: s, u, v, a, t

**Step 2:** Fill in what you know. Write a question mark for what you want to find. Leave blank (or mark with a dash) what is neither given nor asked for.

**Step 3:** The equation you want is the one where all variables are either known or are the one unknown you want. Equivalently, use the "missing variable" column: find the variable that is both not given AND not asked for — that is the missing variable, and use the equation that omits it.

Let us practice this with a simple example before the full worked examples.

---

### Method Example: Finding the Right Equation

**Situation:** A cyclist starts at 5 m/s, accelerates at 2 m/s² for 4 seconds. How far does she travel?

**Step 1:** List variables:
```
s = ?     (what we want)
u = 5 m/s (given)
v = ---   (not given, not asked for)
a = 2 m/s²(given)
t = 4 s   (given)
```

**Step 2:** The variable that is neither given nor asked for is **v**.

**Step 3:** Use the equation that does not contain v. That is Equation 2: **s = ut + ½at²**

```
s = (5)(4) + ½(2)(4²)
s = 20 + ½(2)(16)
s = 20 + 16
s = 36 m
```

The cyclist travels **36 metres**.

---

### The Decision Table

| You know | You want | Use equation |
|----------|----------|--------------|
| u, a, t | v | v = u + at |
| u, a, t | s | s = ut + ½at² |
| u, v, t | s | s = ½(u+v)t |
| u, v, t | a | v = u + at → rearrange |
| u, a, s | v | v² = u² + 2as |
| u, v, s | a | v² = u² + 2as → rearrange |
| v, a, t | s | s = vt - ½at² |
| v, a, t | u | v = u + at → rearrange |

---

## 6. Worked Example 1: Car Accelerating from Rest

### Problem

A car starts from rest at a traffic light. It accelerates uniformly at 3 m/s² for 8 seconds. Calculate:

(a) The car's velocity after 8 seconds.  
(b) The distance the car travels in those 8 seconds.  
(c) The car's velocity when it has travelled exactly 48 metres.

---

### Setting Up

Always start by drawing a quick sketch and listing the SUVAT variables.

```
     REST                                         ?
      |                                           |
   u=0 m/s  ----a=3 m/s²---->  t=8 s  ---->  v=? m/s
      |                                           |
      <-------------- s = ? ---------------------->
```

**Given information:**
```
s = ?
u = 0 m/s    (starts from rest)
v = ?
a = 3 m/s²   (uniform acceleration)
t = 8 s
```

---

### Part (a): Velocity after 8 seconds

**What is missing/not asked for?** s is not needed here.  
**Equation to use:** v = u + at (missing s)

```
v = u + at
v = 0 + (3)(8)
v = 0 + 24
v = 24 m/s
```

**The car is travelling at 24 m/s after 8 seconds.** (That is about 86 km/h — a reasonable highway speed.)

---

### Part (b): Distance travelled in 8 seconds

**Now known:** u=0, a=3, t=8, v=24  
**What is missing/not asked for?** v is now known, but we need s.

We could use any equation that contains s. Let us use Equation 2 (simple to apply):
s = ut + ½at²

```
s = ut + ½at²
s = (0)(8) + ½(3)(8²)
s = 0 + ½(3)(64)
s = ½ × 192
s = 96 m
```

**The car travels 96 metres in 8 seconds.**

Let us verify using Equation 4 as a check:

```
s = ½(u + v)t
s = ½(0 + 24)(8)
s = ½(24)(8)
s = ½ × 192
s = 96 m  ✓
```

Great, same answer.

---

### Part (c): Velocity at the 48-metre mark

**Now the time window changes!** We want to know the speed at a specific distance.

```
s = 48 m   (now given)
u = 0 m/s
v = ?
a = 3 m/s²
t = ---     (not given, not asked)
```

**Missing variable is t.** Use Equation 3: v² = u² + 2as

```
v² = u² + 2as
v² = (0)² + 2(3)(48)
v² = 0 + 288
v² = 288
v  = √288
v  = 16.97 m/s   (approximately 17 m/s)
```

**The car is travelling at about 17 m/s when it reaches the 48-metre mark.**

---

### Sanity Check

Does this make sense? The car reaches 24 m/s at 96 m. At 48 m (exactly halfway), its speed is 17 m/s. Why is 17 m/s not exactly halfway between 0 and 24? Because the car speeds up more slowly at the beginning (when it is going slowly) and covers less distance per second. It is a square root relationship — the answer makes physical sense.

---

## 7. Worked Example 2: Ball Thrown Straight Up

This example introduces **negative acceleration** — one of the most important ideas in SUVAT.

### Problem

A person throws a ball straight upward with an initial velocity of 20 m/s. Assume g = 10 m/s² (gravitational acceleration, acting downward).

(a) How high does the ball go?  
(b) How long does it take to reach the highest point?  
(c) What is the ball's velocity after 3 seconds?  
(d) How long does it take to return to the thrower's hand?

---

### Setting Up — Choosing Positive Direction

This is critical. We must choose a positive direction before writing any equations.

**Choose: upward = positive (+)**

This means:
- Initial velocity is +20 m/s (upward)
- Gravity is -10 m/s² (downward, against the positive direction)
- When the ball is rising, velocity is positive
- When the ball is falling, velocity is negative

```
                        ★ Highest point (v = 0 here)
                       /\
                      /  \
                     /    \
                    /      \
                   / RISING \  FALLING
                  /          \
                 /            \
                /              \
               /                \
    ──────────────────────────────────────
    Thrower's hand (s=0, u=+20 m/s)
    
    ↑ positive direction
    ↓ negative direction (direction of gravity)
```

**Given information:**
```
s = ?
u = +20 m/s
v = ?
a = -10 m/s²  (gravity is DOWNWARD, so negative)
t = ?
```

---

### Part (a): Maximum height

At the highest point, the ball is momentarily stationary — it has stopped going up but not yet started coming down. Therefore:

**v = 0 at the highest point.**

```
v = 0 m/s    (at highest point)
u = 20 m/s
a = -10 m/s²
s = ?
t = ---       (not asked for here)
```

**Missing variable is t.** Use Equation 3: v² = u² + 2as

```
v² = u² + 2as
0² = 20² + 2(-10)(s)
0  = 400 - 20s
20s = 400
s = 20 m
```

**The ball reaches a maximum height of 20 metres.**

---

### Part (b): Time to reach highest point

```
v = 0 m/s
u = 20 m/s
a = -10 m/s²
t = ?
s = 20 m (now known, but we do not need it)
```

**Missing variable is s.** Use Equation 1: v = u + at

```
v = u + at
0 = 20 + (-10)(t)
0 = 20 - 10t
10t = 20
t = 2 s
```

**The ball takes 2 seconds to reach the top.**

---

### Part (c): Velocity after 3 seconds

Now we are asked about the ball's state at t = 3 s. Note that at t = 2 s the ball was at the top — so at t = 3 s, it has been falling for 1 second.

```
u = 20 m/s
a = -10 m/s²
t = 3 s
v = ?
s = ---
```

Use Equation 1: v = u + at

```
v = u + at
v = 20 + (-10)(3)
v = 20 - 30
v = -10 m/s
```

**The velocity after 3 seconds is -10 m/s.**

The negative sign means the ball is moving downward (in our convention, downward is negative). It is falling at 10 m/s. This makes perfect sense — the ball went up for 2 seconds, stopped, then fell. At t=3 s it has been falling for 1 second from rest, and under 10 m/s² acceleration it reaches 10 m/s downward. Perfect.

---

### Part (d): Time to return to the hand

When the ball returns to the thrower's hand, it is back at its starting height. Therefore s = 0.

```
s = 0
u = 20 m/s
a = -10 m/s²
t = ?
```

Use Equation 2: s = ut + ½at²

```
0 = 20t + ½(-10)t²
0 = 20t - 5t²
0 = t(20 - 5t)
```

This gives either t = 0 (the start — mathematically valid, physically the moment it was thrown) or:
```
20 - 5t = 0
5t = 20
t = 4 s
```

**The ball returns to the hand after 4 seconds.**

---

### Symmetry Check

The ball took 2 seconds to go up and 2 seconds to come back down (2 + 2 = 4 s total). This is a beautiful symmetry that occurs whenever air resistance is ignored. The flight is perfectly symmetric about the highest point.

```
Time:    0s        1s        2s        3s        4s
Height:  0 m  --> 15 m --> 20 m --> 15 m --> 0 m
Velocity: +20     +10       0        -10       -20

(velocities are mirror images about the midpoint at t=2s)
```

---

## 8. Worked Example 3: Braking Distance of a Car

### Problem

A car is travelling at 30 m/s (about 108 km/h) when the driver suddenly sees an obstacle and applies the brakes. The brakes produce a deceleration of 6 m/s². 

(a) How long does it take the car to stop?  
(b) What is the braking distance (distance to stop)?  
(c) If the car was only travelling at 15 m/s (half the speed), what would the braking distance be? What does this tell you?

---

### Setting Up

```
         u = 30 m/s                          v = 0 (stopped)
              |    BRAKING    a = -6 m/s²       |
    ══════════|═══════════════════════════════|══════
              <============ s = ? ============>
```

**Given information:**
```
s = ?
u = 30 m/s
v = 0 m/s    (car stops)
a = -6 m/s²  (deceleration means negative, opposing motion)
t = ?
```

---

### Part (a): Time to stop

**Missing variable:** s (not asked yet)  
**Use:** v = u + at

```
v = u + at
0 = 30 + (-6)(t)
0 = 30 - 6t
6t = 30
t = 5 s
```

**The car takes 5 seconds to stop.**

---

### Part (b): Braking distance at 30 m/s

**Now we know t = 5 s.** Use Equation 4 (simplest here):

```
s = ½(u + v)t
s = ½(30 + 0)(5)
s = ½ × 30 × 5
s = 75 m
```

Or alternatively, use Equation 3 (since t is not needed):

```
v² = u² + 2as
0² = 30² + 2(-6)(s)
0  = 900 - 12s
12s = 900
s  = 75 m  ✓
```

**The braking distance is 75 metres.** That is three-quarters the length of a football pitch. At 108 km/h, if you see an obstacle 70 m away, you cannot stop in time even with perfect brakes. This is why speed limits exist.

---

### Part (c): Braking distance at 15 m/s (half the speed)

```
v² = u² + 2as
0² = 15² + 2(-6)(s)
0  = 225 - 12s
12s = 225
s  = 18.75 m
```

**At half the speed, the braking distance is only 18.75 metres — that is one quarter of 75 m.**

| Speed | Braking Distance |
|-------|-----------------|
| 30 m/s | 75 m |
| 15 m/s | 18.75 m |

Halving the speed quartered the stopping distance. This is because braking distance is proportional to v² (from v² = u² + 2as → s = u² / (2×6)). Speed appears squared, so doubling speed quadruples braking distance.

This is one of the most important road safety facts there is. **Double your speed, quadruple your stopping distance.**

```
Speed × 2:   |---s---|  →  |---s---+---s---+---s---+---s---|
             18.75 m      75 m   (4 times longer!)
```

---

## 9. Worked Example 4: Stone Dropped from a Cliff

### Problem

A stone is dropped (not thrown) from the top of a cliff. It hits the water below after 4.5 seconds. Take g = 10 m/s² downward.

(a) What is the height of the cliff?  
(b) What is the stone's velocity just before it hits the water?  

A second stone is thrown downward from the same cliff with an initial speed of 8 m/s.

(c) How much faster does this second stone reach the water compared to the first?

---

### Setting Up

**Choose: downward = positive (+)** (convenient here since everything goes down)

```
     CLIFF TOP
        |
        * (stone released here, u=0)
        |
        |  a = +10 m/s² (gravity, downward = positive)
        |
        |  s = ? (height of cliff)
        |
        v
     ~~~~~  WATER  ~~~~~
```

**Stone 1:**
```
s = ?
u = 0 m/s   (dropped, not thrown)
v = ?
a = 10 m/s²
t = 4.5 s
```

---

### Part (a): Height of the cliff

**Missing variable:** v  
**Use:** s = ut + ½at²

```
s = ut + ½at²
s = (0)(4.5) + ½(10)(4.5²)
s = 0 + 5 × 20.25
s = 101.25 m
```

**The cliff is approximately 101 metres tall.** (About as tall as a 25-floor skyscraper.)

---

### Part (b): Velocity just before impact

```
v = u + at
v = 0 + (10)(4.5)
v = 45 m/s
```

**The stone hits the water at 45 m/s** (162 km/h). This illustrates why cliff diving is dangerous — even from a modest height, impact speed is enormous.

Let us verify with Equation 3:

```
v² = u² + 2as
v² = 0 + 2(10)(101.25)
v² = 2025
v  = 45 m/s  ✓
```

---

### Part (c): Time for the second stone

**Stone 2:**
```
s = 101.25 m   (same cliff)
u = 8 m/s      (thrown downward, same positive direction)
v = ?
a = 10 m/s²
t = ?
```

Use Equation 2: s = ut + ½at²

```
101.25 = 8t + ½(10)t²
101.25 = 8t + 5t²
```

Rearrange into standard quadratic form (at² + bt + c = 0):

```
5t² + 8t - 101.25 = 0
```

Multiply through by 4 to get nicer numbers:

```
20t² + 32t - 405 = 0
```

Using the quadratic formula: t = (-b ± √(b² - 4ac)) / (2a)

```
t = (-32 ± √(32² + 4 × 20 × 405)) / (2 × 20)
t = (-32 ± √(1024 + 32400)) / 40
t = (-32 ± √33424) / 40
t = (-32 ± 182.8) / 40
```

Taking the positive root (negative time is not physical):

```
t = (-32 + 182.8) / 40
t = 150.8 / 40
t ≈ 3.77 s
```

**Stone 2 reaches the water in about 3.77 seconds**, compared to Stone 1's 4.5 seconds. That is 0.73 seconds faster.

Even though the second stone was thrown downward with only 8 m/s extra speed, it arrived nearly 0.73 seconds sooner — because gravity amplifies that initial advantage during the fall.

---

## 10. Deriving SUVAT from Velocity-Time Graphs

There is a beautiful geometric way to understand the SUVAT equations using velocity-time (v-t) graphs. This deepens your understanding and gives you a powerful visual tool.

### Reading a v-t Graph

On a velocity-time graph:
- The **horizontal axis** is time (t)
- The **vertical axis** is velocity (v)
- The **gradient (slope)** of the line = acceleration
- The **area under the line** = displacement

```
velocity (m/s)
   ^
   |              /
 v |____________/
   |           /|
   |          / |
   |         /  |
   |        /   |
   |       /    |
 u |______/     |
   |     /      |
   |    /       |
   |   /        |
   +--|----------|--> time (s)
      0          t
```

This is the v-t graph for uniform acceleration from initial velocity u to final velocity v over time t.

---

### Gradient = Acceleration

The gradient (slope) of the line = rise / run = change in velocity / time taken.

```
gradient = (v - u) / t = a

Therefore: a = (v - u) / t
Rearranging: v = u + at     ← This is Equation 1!
```

---

### Area = Displacement

The area under the line from t=0 to t=T is the displacement.

The shape under the line is a **trapezium** (a shape with two parallel horizontal sides — one of length u at the bottom, one of length v at the top... actually let us look carefully):

```
velocity
   ^
 v |─────────────────/|
   |  (TRAPEZIUM)   / |
   |               /  |
   |              /   |
   |             /    |
 u |────────────/     |
   |           /      |
   |__________/       |
   +──────────────────+──> time
   0                  t
   
Area of trapezium = ½ × (sum of parallel sides) × height
                  = ½ × (u + v) × t
                  = displacement = s

So: s = ½(u + v)t    ← This is Equation 4!
```

We can also split the trapezium into a rectangle and a triangle:

```
velocity
   ^
 v |──────────────/|
   |           /▲▲|
   |         /  ▲▲|  ← triangle, height=(v-u), base=t
   |       /    ▲▲|     area = ½(v-u)t = ½(at)t = ½at²
   |     /      ▲▲|
 u |───/─────────▓|
   |  /   ▓▓▓▓▓▓▓|  ← rectangle, height=u, width=t
   |/    ▓▓▓▓▓▓▓▓|     area = u × t = ut
   +──────────────+──> time
   0              t
   
Total area = ut + ½at²
           = displacement = s

So: s = ut + ½at²    ← This is Equation 2!
```

This geometric approach confirms our algebraic derivations. The equations are not arbitrary formulas — they come directly from the geometry of how velocity changes over time.

---

### Finding Equation 3 Geometrically

We can also show v² = u² + 2as geometrically by noting:

Area of trapezium = ½(u+v)t

But from Equation 1: t = (v-u)/a

So: s = ½(u+v) × (v-u)/a = (v²-u²)/(2a)

Therefore: 2as = v² - u² → v² = u² + 2as ← Equation 3.

---

### Key Visual Insight: Negative Acceleration

When acceleration is negative (deceleration), the v-t graph slopes downward:

```
velocity
   ^
 u |──────\
   |       \
   |        \
   |         \
   |          \
 v |           \
   |            \
   +─────────────\──> time
   0              t
   
The area is still a trapezium, still = ½(u+v)t = s.
But the slope (gradient) is now negative: a = (v-u)/t < 0
```

And when the object comes to rest (v=0):

```
velocity
   ^
 u |──────\
   |       \
   |        \
   |         \
   |          \
 0 |───────────\────> time
   0            t
   
Area (triangle) = ½ × base × height = ½ × t × u = s
This matches s = ½(u+v)t with v=0: s = ½(u)t
```

---

## 11. Common Mistakes to Avoid

Over years of teaching physics, the same errors come up again and again. Here is a catalogue of the most common ones, with fixes.

---

### Mistake 1: Not Defining a Positive Direction

**Wrong approach:** Starting calculations without deciding which direction is positive.

**What happens:** You plug in 20 m/s for a velocity going upward and 10 m/s² for gravity also positive, without noticing they point in opposite directions. Your answer is completely wrong.

**Fix:** Always write at the top of your solution: "Let upward = positive" (or whatever you choose). Then check every quantity — does it point in the positive direction? If yes, write it as positive. If no, write it as negative.

---

### Mistake 2: Using a = 9.8 or 10 m/s² as Positive When It Should Be Negative

**Wrong approach:**
```
Ball thrown up, u = 20 m/s (upward = positive)
a = 9.8 m/s²  ← WRONG, gravity is downward!
```

**Fix:**
```
If upward = positive, then a = -9.8 m/s²
If downward = positive, then a = +9.8 m/s²
```

Gravity is 9.8 m/s² downward. Its sign depends entirely on which direction you chose as positive.

---

### Mistake 3: Confusing Displacement and Distance

**Wrong approach:** A ball travels 20 m upward and then 8 m back down. The question asks for displacement. Student writes s = 28 m.

**Fix:** Displacement = final position - starting position = 20 - 8 = 12 m upward. Distance = 20 + 8 = 28 m.

```
START ──────────────────────────────────────────► HIGHEST POINT
  0 m                                             20 m
  
  ◄───────────────────────────── falls 8 m ──────────────────
  
CURRENT POSITION = 20 - 8 = 12 m from start
Displacement = 12 m (upward)
Distance traveled = 20 + 8 = 28 m
```

---

### Mistake 4: Applying SUVAT to Non-Uniform Acceleration

**Wrong approach:** Using s = ut + ½at² for a car that accelerates at different rates at different moments.

**Fix:** Check whether the problem says "uniform acceleration" or "constant acceleration" before applying SUVAT. If acceleration changes, you cannot use these equations directly.

---

### Mistake 5: Forgetting That v = 0 at the Highest Point

**Wrong approach:** Student trying to find the time to reach the top of a throw, but they do not use v = 0.

**Fix:** At the highest point of any projectile thrown straight upward, the vertical velocity is always zero. The object momentarily stops before reversing direction. v = 0 at the top is a standard starting condition for many problems.

---

### Mistake 6: Solving the Quadratic and Taking the Wrong Root

**Wrong approach:** Getting two solutions from a quadratic (e.g., t = -2 and t = 5) and writing "t = -2 or t = 5."

**Fix:** Time cannot be negative in most physical situations. If one root is negative, discard it. The positive root is the physical answer.

Occasionally both positive roots are meaningful (e.g., a ball passing through a height of 10 m on the way up AND on the way down). In that case, examine the context to see which one the question refers to.

---

### Mistake 7: Wrong Units

**Wrong approach:** Mixing km/h and m/s in the same calculation.

**Fix:** Convert everything to base SI units first.
- 1 m/s = 3.6 km/h
- To convert km/h to m/s: divide by 3.6
- To convert m/s to km/h: multiply by 3.6

```
60 km/h = 60 / 3.6 = 16.67 m/s
30 m/s = 30 × 3.6 = 108 km/h
```

---

### Mistake 8: Not Checking the Answer for Reasonableness

**Fix:** After every calculation, ask yourself: "Does this make sense?"

- A car decelerating to rest in 0.001 seconds? Not reasonable.
- A ball reaching 300 m height when thrown at 20 m/s? Not reasonable (maximum is 20 m).
- A time of -3 seconds for when something happens? Not physically reasonable.

Develop the habit of estimating and checking.

---

### Mistake 9: Misidentifying u and v

**Wrong approach:** A ball is rolling at 5 m/s and slows to 2 m/s. Student writes u=2, v=5 without checking the direction of time.

**Fix:** u is ALWAYS the velocity at the START of your chosen time window. v is ALWAYS at the END. If the ball slows from 5 to 2, then u=5 and v=2.

---

### Mistake 10: Using the Wrong SUVAT Equation

**Wrong approach:** Just picking an equation randomly, or always using v = u + at.

**Fix:** Use the known/unknown method every time. List s, u, v, a, t. Mark what you know and what you want. Identify the missing variable. Use the equation that omits it.

```
ALWAYS DO THIS:
s = __     (known? needed?)
u = __     (known? needed?)
v = __     (known? needed?)
a = __     (known? needed?)
t = __     (known? needed?)

Missing variable (not given, not needed) = ___
Use equation without that variable.
```

---

## 12. Quick-Reference Cheat Sheet

```
╔══════════════════════════════════════════════════════════════════╗
║                    SUVAT CHEAT SHEET                             ║
╠══════════════════════════════════════════════════════════════════╣
║  VARIABLES                                                       ║
║  s = displacement (m)         u = initial velocity (m/s)         ║
║  v = final velocity (m/s)     a = acceleration (m/s²)            ║
║  t = time (s)                                                    ║
╠══════════════════════════════════════════════════════════════════╣
║  EQUATIONS          MISSING VARIABLE                             ║
║  v = u + at              s                                       ║
║  s = ut + ½at²           v                                       ║
║  v² = u² + 2as           t                                       ║
║  s = ½(u + v)t           a                                       ║
║  s = vt - ½at²           u                                       ║
╠══════════════════════════════════════════════════════════════════╣
║  METHOD                                                          ║
║  1. List all 5 variables.                                        ║
║  2. Mark known, unknown (want), and irrelevant.                  ║
║  3. Pick the equation that omits the irrelevant variable.        ║
║  4. Solve algebraically.                                         ║
╠══════════════════════════════════════════════════════════════════╣
║  SPECIAL CASES                                                   ║
║  Starting from rest: u = 0                                       ║
║  Coming to rest:     v = 0                                       ║
║  At highest point:   v = 0                                       ║
║  Gravity:            a = 9.8 m/s² (use 10 for estimates)         ║
╠══════════════════════════════════════════════════════════════════╣
║  SIGNS                                                           ║
║  Choose positive direction first. Everything opposite = negative.║
║  Deceleration = negative acceleration (opposing motion).         ║
╚══════════════════════════════════════════════════════════════════╝
```

---

## 13. Practice Problems

Try these yourself before checking the answers.

---

### Problem 1

A train accelerates from 5 m/s to 25 m/s at a constant rate over 40 seconds.
(a) What is the acceleration?  
(b) How far does the train travel during this time?

**Answers:** (a) a = 0.5 m/s²   (b) s = 600 m

---

### Problem 2

A motorcycle is travelling at 18 m/s when the rider brakes, producing a deceleration of 4.5 m/s².
(a) How far does the motorcycle travel before stopping?  
(b) How long does it take to stop?

**Answers:** (a) s = 36 m   (b) t = 4 s

---

### Problem 3

A ball is thrown upward with a velocity of 15 m/s. Take g = 10 m/s².
(a) What is the maximum height?  
(b) What is the velocity after 2.5 seconds?  
(c) When does the ball return to the starting height?

**Answers:** (a) s = 11.25 m   (b) v = -10 m/s (downward)   (c) t = 3 s

---

### Problem 4

A cheetah accelerates from rest to 30 m/s over a distance of 100 m.
(a) What is the cheetah's acceleration?  
(b) How long does this acceleration take?

**Answers:** (a) a = 4.5 m/s²   (b) t = 6.67 s

---

### Problem 5 — Challenge

A car is travelling at 72 km/h. The driver sees a traffic light turn red 80 m ahead and applies the brakes.

(a) Convert 72 km/h to m/s.  
(b) If the car stops exactly at the line, what deceleration was applied?  
(c) How long did the stop take?  

**Answers:** (a) 20 m/s   (b) a = 2.5 m/s² deceleration   (c) t = 8 s

---

## 14. Summary

This chapter introduced the SUVAT equations — the five fundamental equations of uniform acceleration. Here are the key takeaways:

- **SUVAT stands for:** Displacement (s), Initial velocity (u), Final velocity (v), Acceleration (a), Time (t).

- **SUVAT only applies when acceleration is constant (uniform).** If acceleration changes, these equations cannot be applied directly.

- **All five equations come from two simple definitions:** a = (v-u)/t and s = ((u+v)/2)×t. Everything else is algebraic rearrangement.

- **The five equations are:**
  - v = u + at
  - s = ut + ½at²
  - v² = u² + 2as
  - s = ½(u + v)t
  - s = vt - ½at²

- **To choose the right equation,** identify which variable is both not given and not needed — then use the equation that does not contain that variable.

- **Always define a positive direction** before writing any equations. Stick to it consistently throughout the problem.

- **Special conditions to recognise immediately:**
  - Object starts from rest: u = 0
  - Object comes to rest: v = 0
  - Object at the highest point of a throw: v = 0
  - Gravity: a = 9.8 m/s² directed downward

- **Braking distance is proportional to v²:** doubling speed quadruples stopping distance. This has critical road safety implications.

- **The v-t graph confirms the equations geometrically:** the gradient is acceleration and the area under the graph is displacement.

- **Common mistakes** include wrong signs for acceleration, confusing displacement and distance, using SUVAT for non-uniform acceleration, and taking the physically meaningless negative root of a quadratic.

---

## 15. Key Equations

```
SUVAT EQUATIONS (valid for constant acceleration only)

1.  v = u + at

2.  s = ut + ½at²

3.  v² = u² + 2as

4.  s = ½(u + v)t

5.  s = vt - ½at²

WHERE:
    s = displacement (metres, m)
    u = initial velocity (metres per second, m/s)
    v = final velocity (metres per second, m/s)
    a = acceleration (metres per second squared, m/s²)
    t = time (seconds, s)

DERIVED RELATIONS (useful to remember):

    a = (v - u) / t             (definition of acceleration)

    average velocity = (u + v) / 2   (for uniform acceleration only)

    Braking distance  ∝  v²     (proportional to square of initial speed)

STANDARD VALUES:
    g = 9.8 m/s²  (gravitational acceleration near Earth's surface)
    g ≈ 10 m/s²   (useful approximation for mental estimates)
```

---

*End of Chapter 08. In Chapter 09, we extend these ideas to motion in two dimensions — projectile motion — where we apply SUVAT separately to the horizontal and vertical components of motion.*
