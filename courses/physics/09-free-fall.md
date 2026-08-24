# Chapter 09: Free Fall

> **"And yet it moves — all things fall at the same rate, whether a feather or a cannonball, when air steps aside."**
> — inspired by Galileo Galilei

---

## Table of Contents

1. [What Is Free Fall?](#1-what-is-free-fall)
2. [Galileo's Famous Experiment — Busting a 2000-Year-Old Myth](#2-galileos-famous-experiment--busting-a-2000-year-old-myth)
3. [The Acceleration Due to Gravity: g](#3-the-acceleration-due-to-gravity-g)
4. [The SUVAT Equations for Free Fall](#4-the-suvat-equations-for-free-fall)
5. [An Object Dropped From Rest](#5-an-object-dropped-from-rest)
6. [An Object Thrown Upward](#6-an-object-thrown-upward)
7. [The Symmetry of Free Fall](#7-the-symmetry-of-free-fall)
8. [Weightlessness — What It Actually Means](#8-weightlessness--what-it-actually-means)
9. [Air Resistance and Terminal Velocity](#9-air-resistance-and-terminal-velocity)
10. [Worked Examples](#10-worked-examples)
11. [Summary](#11-summary)
12. [Key Equations](#12-key-equations)

---

## 1. What Is Free Fall?

Imagine you are standing on the roof of a tall building. You hold a ball out in front of you — and let go.

The ball drops. Fast. Faster and faster. Until it hits the ground far below.

Now ask yourself: what is making the ball speed up as it falls? The answer is gravity — the invisible force that pulls every object with mass toward the centre of the Earth.

**Free fall** is the motion of an object that is falling under the influence of gravity alone, with no other forces acting on it. In free fall, only one thing is happening to the object: gravity is pulling it downward, giving it a constant acceleration.

That last word is the key. **Constant acceleration** — the object speeds up by the same amount every single second. It does not gradually slow down and stop. It does not sometimes speed up and sometimes slow down. It gets faster and faster in a perfectly predictable, mathematically regular way.

### The Ideal World vs. the Real World

Physics textbooks, including this one, usually start with a simplified model called **ideal free fall** (also called "free fall in a vacuum"). In this model:

- There is no air. No wind resistance. No drag.
- The only force is gravity.
- The object falls in a straight vertical line.

This is an idealization — real falling objects do have air resistance. But the ideal model is surprisingly accurate for dense, compact objects (like a metal ball or a stone) falling short distances at everyday speeds. We will deal with air resistance properly in Section 9.

For now, let us build the clean, simple picture first. Once you understand ideal free fall, the realistic version with air resistance is just a small modification on top.

### Why Does Free Fall Matter?

Free fall is one of the most important ideas in all of physics, for several reasons:

1. It was the first motion that scientists described mathematically with precision. Galileo's work on falling objects in the 1600s marked the birth of modern experimental science.

2. It connects directly to the universal law of gravitation — the same physics that keeps the Moon orbiting the Earth and the Earth orbiting the Sun.

3. It is the foundation of **projectile motion** — the topic that comes immediately after this chapter, which explains how a thrown ball, a fired arrow, or a kicked football moves through the air.

4. It reveals deep truths about the nature of mass and weight — truths so surprising that Einstein built his entire theory of general relativity on top of them.

So take this chapter seriously. Every concept here will pay dividends later.

---

## 2. Galileo's Famous Experiment — Busting a 2000-Year-Old Myth

### What Everyone "Knew" Before Galileo

For nearly 2,000 years, educated people believed what the ancient Greek philosopher **Aristotle** had written about falling objects: heavy things fall faster than light things.

It seemed obvious. Common sense. Drop a rock and a leaf, and the rock clearly hits the ground first. Drop a cannon ball and a feather — the cannon ball wins easily.

So nobody questioned it. Aristotle said it, generation after generation repeated it, and it was simply accepted as fact.

The problem? Nobody actually tested it carefully.

### Galileo Steps In (around 1590)

**Galileo Galilei** (1564–1642), an Italian scientist, was not satisfied with "common sense." He wanted to understand nature through careful observation and experiment — a revolutionary idea at the time.

He reasoned about it logically first. Aristotle's claim leads to a contradiction:

Suppose a heavy stone (mass 10 kg) falls faster than a light stone (mass 1 kg). Now tie them together with a string. What happens?

- Aristotle would say: the heavy stone wants to fall fast, the light stone wants to fall slow. Together, they should fall at an intermediate speed — slower than the heavy stone alone.
- But the tied-together system has a combined mass of 11 kg — heavier than either stone. So Aristotle would also say it should fall faster than the heavy stone.

Both conclusions follow from Aristotle's rule. They completely contradict each other. The rule cannot be right.

### The Leaning Tower of Pisa

According to legend, Galileo climbed the famous Leaning Tower of Pisa and dropped two cannonballs of different weights simultaneously from the top.

```
         LEANING TOWER OF PISA
              _____
             |     |
             |     |    <- Top (~56 m high)
           .-'     '-.
          |           |
          |           |
          |  GALILEO  |   <-- drops both balls
          |     *     |
          |    O O    |  <- heavy ball and light ball
          |           |
          |           |
          |           |
          |___________|
              _____
              |   |
          ____|___|____   <- Ground level
          
   Heavy ball: 10 kg       Light ball: 1 kg
   
   What happened?
   
   Both balls hit the ground AT THE SAME TIME!
```

The crowd watching below reportedly gasped. Two balls of very different weights — released at exactly the same moment — hit the ground simultaneously.

Heavy objects do NOT fall faster than light objects (ignoring air resistance).

### Why Does This Happen?

This is the deep beautiful result: gravity pulls harder on heavier objects (more mass = more gravitational force), BUT heavier objects also need more force to accelerate them (more mass = more inertia, more resistance to acceleration).

These two effects cancel perfectly. The result is that every object — regardless of mass — accelerates at the exact same rate under gravity.

This "coincidence" is actually one of the most profound facts in physics. Einstein spent years trying to understand it deeply, and eventually it became the cornerstone of his theory of general relativity.

### The Moon Drop Experiment

In 1971, astronaut David Scott stood on the surface of the Moon (which has no atmosphere) and dropped a hammer and a feather at the same time. With millions watching on television, both objects hit the ground at exactly the same time.

Galileo was right.

---

## 3. The Acceleration Due to Gravity: g

### The Value of g

Near the surface of the Earth, every freely falling object accelerates downward at:

```
g = 9.8 m/s²   (accurate value)
g ≈ 10 m/s²    (easy value for calculations)
```

This constant **g** is called the **acceleration due to gravity** or sometimes the **gravitational acceleration**.

What does g = 9.8 m/s² actually mean?

It means: every second an object is in free fall, its downward speed increases by 9.8 metres per second.

| Time of fall | Speed gained |
|---|---|
| After 0 seconds | 0 m/s (just released) |
| After 1 second | 9.8 m/s ≈ 10 m/s |
| After 2 seconds | 19.6 m/s ≈ 20 m/s |
| After 3 seconds | 29.4 m/s ≈ 30 m/s |
| After 4 seconds | 39.2 m/s ≈ 40 m/s |
| After 5 seconds | 49.0 m/s ≈ 50 m/s |

Notice the pattern: every second, the speed increases by approximately 10 m/s. After 5 seconds of free fall, the object is moving at 50 m/s — roughly 180 km/h!

### g Varies Slightly with Location

The value g = 9.8 m/s² is an average. In reality, g varies slightly:

| Location | g (m/s²) | Reason |
|---|---|---|
| At the equator | 9.78 | Farther from Earth's centre (Earth bulges here); also centrifugal effect |
| At the poles | 9.83 | Closer to Earth's centre |
| At sea level | 9.81 | Standard reference |
| At the top of Mt. Everest | 9.77 | High altitude = farther from centre |
| On the Moon | 1.62 | Moon is much less massive |
| On Mars | 3.72 | Mars is less massive |
| On Jupiter | 24.8 | Jupiter is much more massive |

For this course, we always use g = 9.8 m/s² (or 10 m/s² to keep arithmetic simple).

### Direction of g

This is important: **g always acts downward**, toward the centre of the Earth.

When we set up equations, we need to choose a positive direction. There are two common conventions:

**Convention A (downward positive):**
- Downward = positive
- g = +9.8 m/s²
- Objects dropped from rest: velocity becomes more positive over time

**Convention B (upward positive):**
- Upward = positive
- g = -9.8 m/s²
- Objects dropped from rest: velocity becomes more negative (more downward) over time

In this chapter, we will use **Convention B (upward positive)** because it matches our everyday intuition: throwing something up gives it positive velocity, and it comes back down. We will write g = -9.8 m/s² in our equations.

Always state your sign convention at the start of a problem. This avoids the most common mistakes in free fall calculations.

---

## 4. The SUVAT Equations for Free Fall

If you completed Chapter 07 (or Chapter 08) on kinematics, you already know the SUVAT equations — the five equations that describe motion with constant acceleration.

Free fall is simply constant acceleration with a = -g = -9.8 m/s² (upward positive convention).

### Refresher: the SUVAT Variables

| Symbol | Meaning | Unit |
|---|---|---|
| s | displacement (distance moved, with direction) | metres (m) |
| u | initial velocity (speed at the start) | m/s |
| v | final velocity (speed at the end) | m/s |
| a | acceleration | m/s² |
| t | time | seconds (s) |

For free fall problems, we simply substitute a = -9.8 m/s² (or a = -10 m/s² for easy numbers).

### The Five SUVAT Equations

```
(1)  v = u + at
(2)  s = ut + (1/2)at²
(3)  v² = u² + 2as
(4)  s = (u + v)/2 × t
(5)  s = vt - (1/2)at²
```

For free fall specifically, these become (with a = -g):

```
(1)  v = u - gt
(2)  s = ut - (1/2)gt²
(3)  v² = u² - 2gs
(4)  s = (u + v)/2 × t
(5)  s = vt + (1/2)gt²
```

### Which Equation to Use?

A simple strategy:

1. List the variables you know (s, u, v, a, t).
2. Identify the variable you want to find.
3. Pick the equation that contains those variables and doesn't involve any unknown you don't need.

Let us see this in action over the next two sections.

---

## 5. An Object Dropped From Rest

### Setting Up the Problem

"Dropped from rest" means the object starts with zero velocity: u = 0.

The only force acting on it is gravity, so a = -g = -9.8 m/s² (upward positive).

The object falls downward, so displacement will be negative (downward direction in our convention).

Let us define:
- Upward = positive
- The starting point = position zero (s = 0 at t = 0)
- The object is released with u = 0

### Speed After Each Second

Using equation (1): v = u - gt = 0 - gt = -gt

| Time (s) | Velocity (m/s) | Speed (magnitude) |
|---|---|---|
| 0 | 0 | 0 m/s |
| 1 | -9.8 | 9.8 m/s downward |
| 2 | -19.6 | 19.6 m/s downward |
| 3 | -29.4 | 29.4 m/s downward |
| 4 | -39.2 | 39.2 m/s downward |
| 5 | -49.0 | 49.0 m/s downward |

The negative sign just means "downward." The speed (how fast it is moving regardless of direction) increases by 9.8 m/s every second.

### Distance Fallen After Each Second

Using equation (2): s = ut - (1/2)gt² = 0 - (1/2)(9.8)t² = -4.9t²

The distance fallen (taking magnitude) = 4.9t²

| Time (s) | Distance fallen (m) | Approx. (using g=10) |
|---|---|---|
| 0 | 0 | 0 m |
| 1 | 4.9 | 5 m |
| 2 | 19.6 | 20 m |
| 3 | 44.1 | 45 m |
| 4 | 78.4 | 80 m |
| 5 | 122.5 | 125 m |

Notice something important: the distances are NOT equal each second. In the first second, the ball falls 5 m. In the second second (between t=1 and t=2), it falls 15 m. In the third second (between t=2 and t=3), it falls 25 m.

The distances in successive seconds go: 5, 15, 25, 35, 45, ... — always increasing by 10 m more each second. This is the signature of constant acceleration.

### ASCII Diagram: Height vs. Time for a Dropped Ball

```
Height above ground
(m)
 |
100 |  * (start here, t=0)
    |
 90 |
    |
 80 |
    |
 70 |
    |
 60 |
    |
 50 |
    |   * (t=1, at 95.1 m if started from 100m)
 40 |
    |
 30 |
    |
 20 |       * (t=2, at 80.4 m)
    |
 10 |
    |               * (t=3, at 55.9 m)
  0 |__________________________* (t=4.5s, hits ground at ~99m)
          1    2    3    4    5  --> Time (s)

Note: curve bends DOWN sharply — this is parabolic motion.
The ball falls further in each successive second.
```

A more accurate picture:

```
  t=0   t=1   t=2   t=3
   |     |     |     |
   *                        <- height 100 m (dropped)
   |
   |  *                     <- height 95.1 m (after 1s)
   |
   |        *               <- height 80.4 m (after 2s)
   |
   |              *         <- height 55.9 m (after 3s)
   |
   |                   *    <- height 21.6 m (after 4s)
   |
   ___________________________  <- ground (at about t=4.52s)
```

### Worked Example 1: Ball Dropped from a Bridge

**Problem:** A ball is dropped (released from rest) from a bridge that is 45 metres above the river. 
(a) How long does it take to hit the water?  
(b) What is its speed just before it hits?

**Given information:**
- u = 0 (dropped from rest)
- a = -9.8 m/s² (gravity, downward)
- s = -45 m (displacement is 45 m downward = negative in our convention)
- g ≈ 10 m/s² for easy calculation

**Part (a): Find t**

Known: s, u, a. Want: t. Use equation (2): s = ut + (1/2)at²

With upward positive:
```
s = ut - (1/2)gt²
-45 = (0)(t) - (1/2)(10)(t²)
-45 = -5t²
45 = 5t²
t² = 45/5 = 9
t = √9 = 3 seconds
```

The ball takes **3 seconds** to hit the water.

**Part (b): Find v**

Use equation (1): v = u - gt
```
v = 0 - (10)(3)
v = -30 m/s
```

The negative sign means downward. The speed is **30 m/s** (about 108 km/h — fast enough to be dangerous!).

**Check using equation (3):**
```
v² = u² - 2gs
v² = 0 - 2(-10)(-45)
v² = -900   ← Hmm, that should be positive...
```

Wait — let us be careful with signs. With upward positive, g = +9.8. So:
```
v² = u² - 2(9.8)(s)    where s = -45
v² = 0 - 2(9.8)(-45)
v² = 0 + 882
v² = 882
v = 29.7 m/s ≈ 30 m/s  ✓
```

The answer checks out.

---

## 6. An Object Thrown Upward

Now for a more interesting situation: what if someone throws a ball straight up into the air?

This is the same physics — gravity acting downward with a = -9.8 m/s². The only difference is the starting velocity: now u is positive (upward).

### What Happens Step by Step

When you throw a ball upward:

1. The ball leaves your hand with some upward velocity u (positive).
2. As it rises, gravity pulls it back — the velocity decreases every second.
3. At the very top of its path, the velocity is zero for a brief instant. The ball has stopped rising but not yet started falling.
4. Then it starts falling back down — velocity becomes negative (downward).
5. It returns to your hand (or the ground) with the same speed as it was thrown (if it returns to the same height).

```
Ball thrown upward — velocity at each moment:

          TOP (v = 0)
            *
           /|\
          / | \
         /  |  \
        /   |   \
       /    |    \
      *     |     *
     ↑      |      ↓
  v = +u    |    v = -u
  (going up)|  (coming down)
            |
   THROW    |   CATCH
   (here)   | (or lands)
            |
        t=0   t=T (total time)

The ball is in free fall the ENTIRE time —
both going up AND coming down.
Gravity acts downward the whole time.
```

### Finding the Maximum Height

At the highest point, the velocity v = 0.

Use equation (1): v = u - gt
```
0 = u - g × t_top
t_top = u / g
```

So the time to reach the top = initial speed divided by g.

Use equation (3) to find the maximum height (s_max):
```
v² = u² - 2g × s_max
0 = u² - 2g × s_max
s_max = u² / (2g)
```

### Worked Example 2: Ball Thrown Upward

**Problem:** A ball is thrown straight upward with an initial velocity of 20 m/s.  
(a) How long does it take to reach the highest point?  
(b) What is the maximum height reached?  
(c) What is the total time before the ball returns to the thrower's hand?

**Given:** u = +20 m/s, a = -g = -10 m/s², return to starting height so final s = 0

**Part (a): Time to reach top**

At the top, v = 0. Use v = u + at (with a = -g):
```
v = u - gt_top
0 = 20 - 10 × t_top
10 × t_top = 20
t_top = 2 seconds
```

**Part (b): Maximum height**

Use s_max = u² / (2g):
```
s_max = (20)² / (2 × 10)
s_max = 400 / 20
s_max = 20 metres
```

The ball rises 20 metres above the throwing point.

**Part (c): Total flight time**

We will cover this in detail in the next section, but: by symmetry, the ball takes the same time to come down as to go up.

Total time = 2 × t_top = 2 × 2 = **4 seconds**

We can verify this. When the ball returns to s = 0:
```
s = ut - (1/2)gt²
0 = 20t - (1/2)(10)t²
0 = 20t - 5t²
0 = t(20 - 5t)
t = 0  OR  t = 4 seconds
```

t = 0 is when the ball was thrown. t = 4 seconds is when it returns. Confirmed.

### Velocity at Every Moment

```
Velocity of ball thrown at u = +20 m/s, g = 10 m/s²:

 v (m/s)
 +20 |*
     | \
 +15 |  \
     |   \
 +10 |    \
     |     \
  +5 |      \
     |       \
   0 |--------*--------  <- top of flight (t = 2s)
     |         \
  -5 |          \
     |           \
 -10 |            \
     |             \
 -15 |              \
     |               \
 -20 |                *
     |________________________
        0   1   2   3   4    t (s)
        
Notice: The line is perfectly straight (constant slope = -g).
The velocity decreases at exactly 10 m/s every second.
```

---

## 7. The Symmetry of Free Fall

One of the most beautiful and useful properties of free fall is its **symmetry**.

### Going Up Takes the Same Time as Coming Down

If an object is thrown upward and returns to the same height it was thrown from, then:

- **Time going up = Time coming down**
- **Speed at any height on the way up = Speed at the same height on the way down**

This is because the motion is perfectly symmetric about the highest point.

```
HEIGHT vs TIME for ball thrown upward:

Height
(m)
  20 |          *         <- maximum height (t = 2s)
     |        /   \
  15 |       /     \
     |      /       \
  10 |     /         \
     |    /           \
   5 |   /             \
     |  /               \
   0 | /                 \
     |________________________
     0    1    2    3    4    time (s)
     
     <- Going up ->|<- Coming down ->
     
     Perfectly symmetric: the shape is a parabola.
     Left half = mirror image of right half.
```

### Speed at Every Height is the Same Up and Down

Suppose our ball is at a height of 15 m. Let us find its speed:

On the way up:
```
v² = u² - 2gs
v² = (20)² - 2(10)(15)
v² = 400 - 300 = 100
v = 10 m/s  (upward, positive)
```

On the way down (same height, same equation):
```
v² = u² - 2gs = 100
v = -10 m/s  (downward, negative)
```

Same speed (10 m/s), opposite direction. The magnitude is the same.

### Why Does This Symmetry Exist?

Because energy is conserved. When the ball is going up, it is trading kinetic energy (speed) for gravitational potential energy (height). When it comes back down, it trades that potential energy back into kinetic energy. At the same height, it has given up (or gained back) the same amount of energy — so the speeds are equal.

We will learn more about this in the Energy chapter. For now, just use the symmetry rule: **time up = time down, speed at same height is same**.

### Practical Use of the Symmetry Rule

If you know the time to reach the top (t_top), you immediately know:

- Total flight time = 2 × t_top
- The ball returns with speed equal to the launch speed (just downward instead of upward)

This cuts your calculation work in half for many problems.

---

## 8. Weightlessness — What It Actually Means

### You Have Felt Weightlessness

Have you ever been in a lift (elevator) that is going down? There is a brief moment, just after it starts moving, when you feel lighter. Your stomach seems to float. This is a taste of partial weightlessness.

Or think about a rollercoaster going over the top of a hill. At the very peak, if the coaster is going fast enough, you feel like you are floating out of your seat.

These are everyday examples of reduced apparent weight — and they point to something deep.

### Weight vs. Mass

First, a critical distinction:

- **Mass** is the amount of matter in an object. It never changes. It is measured in kilograms (kg).
- **Weight** is the gravitational force pulling an object downward. Weight = mass × g. It is measured in Newtons (N).

But here is the thing: when we say we "feel" our weight, we are not actually feeling gravity directly. We feel the **normal force** — the push from the floor or chair that holds us up.

Stand on a scale. The scale reading is not measuring how hard gravity pulls you down. It is measuring how hard the floor is pushing you up.

When those two are equal (standing still), the scale reads your weight. But when the floor stops pushing as hard, the scale reads less — and you feel lighter.

### Free Fall and the Apparent Weight

When you are in true free fall — nothing holding you up, only gravity acting on you — the floor (or chair, or anything beneath you) is also falling at the same rate you are.

Nothing is pushing up on you. You feel no weight at all.

This is genuine weightlessness.

```
ASTRONAUT IN ORBIT:

                 Earth
                  ___
                 /   \
                | Earth|
                 \___/
                
           Astronaut
             \O/
              |    <- always falling toward Earth
             /\    <- but also moving sideways so fast
                      that the ground curves away beneath
                      them at the same rate they fall.
                      
They are in continuous free fall — 
falling toward Earth but always missing it.

They feel ZERO weight because nothing is holding them up.
```

### Orbit Is Continuous Free Fall

This is one of the most mind-bending facts in physics: **astronauts in orbit are not in zero gravity**. Earth's gravity still pulls on them quite strongly (only slightly weaker than at the surface).

They feel weightless because they are in continuous free fall around the Earth. They are falling toward the Earth the entire time, but they are also moving sideways so fast (about 7,900 m/s!) that the Earth curves away beneath them at the same rate they fall.

They are perpetually falling and perpetually missing the ground. That is an orbit.

The International Space Station (ISS) is at about 400 km altitude. At that height, g is about 8.7 m/s² — only slightly less than at the surface. The astronauts are definitely in a strong gravitational field. They feel weightless only because they are in free fall.

### Microgravity in Practice

Aerospace engineers train astronauts using a specially modified aircraft nicknamed the "Vomit Comet." The plane flies in large parabolic arcs. At the top of each arc, the plane is in free fall for about 25 seconds, and the people inside experience weightlessness — floating around the cabin.

This is just free fall inside a falling container. Same physics, different scale.

---

## 9. Air Resistance and Terminal Velocity

So far, we have assumed no air resistance. In the real world, air resistance is significant, and it changes the story dramatically for lightweight or large objects.

### What Is Air Resistance?

**Air resistance** (also called **drag**) is the force that air exerts on a moving object in the opposite direction to its motion. When you fall downward, air pushes upward on you.

Air resistance depends on:
- The shape of the object (streamlined vs. flat)
- The size (surface area) of the object
- The speed of the object — faster = more drag
- The density of the air

The faster you go, the more drag you experience. This is the crucial point.

### The Race to Balance

When you first jump out of a plane (with a parachute, please!), you start with zero speed. Gravity is pulling you down at 9.8 m/s², and air resistance is very small (because you are slow).

As you speed up, the air resistance increases. The net downward force (gravity minus air resistance) gets smaller. Your acceleration gets smaller.

Eventually, you are going so fast that air resistance exactly equals the gravitational force. Net force = 0. Acceleration = 0. You stop speeding up.

You have reached **terminal velocity**.

```
FORCES ON A SKYDIVER

At the start (just jumped):        Later (speeding up):     At terminal velocity:

   [Air resistance = 0]               [Air resistance]          [Air resistance]
                                           ^                          ^
                                           |                          |
    Skydiver                           Skydiver                   Skydiver
   \  O  /                            \  O  /                    \  O  /
    \ | /                              \ | /                      \ | /
      |                                  |                          |
      v                                  v                          |  <-- no net force!
   [Gravity]                          [Gravity]                  [Gravity]
   
   Net force = Gravity                 Net force =                Net force = 0
   Acceleration = g                    Gravity - Drag             Acceleration = 0
                                       (less than g)              Constant speed!
```

### Terminal Velocity Values

| Object | Terminal Velocity |
|---|---|
| Skydiver (face down, spread out) | ~55 m/s (200 km/h) |
| Skydiver (head down, streamlined) | ~90 m/s (325 km/h) |
| Human with parachute deployed | ~5–6 m/s (safe landing speed) |
| Raindrop | ~9 m/s |
| Ping pong ball | ~9 m/s |
| Feather | ~0.5 m/s |
| Large hailstone | ~30 m/s |

The parachute works by massively increasing surface area, which greatly increases air resistance at any given speed. This lowers the terminal velocity from a fatal ~55 m/s to a survivable ~5 m/s.

### Why Galileo Was Right (With a Caveat)

Galileo's result — all objects fall at the same rate — is exactly correct in a vacuum (no air). In practice, the effect of air resistance depends on the ratio of drag force to weight.

A feather has very little weight but a large surface area relative to its weight. Air resistance dominates easily. It falls slowly.

A cannonball has enormous weight and relatively small surface area. Air resistance is negligible compared to gravity. It falls almost at the theoretical free fall rate.

This is why in everyday experience, it looks like heavy things fall faster — because the heavy things we compare (cannonballs, rocks, apples) also happen to be dense and compact, with small air resistance. If you compared two balls of the same size but different masses (one solid lead, one hollow plastic), in air the lead ball would win — not because of its weight directly, but because its weight far outpaces its air resistance.

### Velocity vs. Time with Air Resistance

```
Speed
(m/s)

  60 |              ..........------------  <- terminal velocity (55 m/s)
     |          ...
  50 |        ..
     |      ..
  40 |     .
     |    .
  30 |   .           Free fall WITHOUT air resistance:
     |  .
  20 | .            /  <- speed keeps increasing at constant rate
     |.            /     (straight line, slope = g)
  10 |            /
     |           /  <- free fall line
   0 |__________/________________________________
     0    2    4    6    8    10   12   14    time (s)

Solid curved line = real fall with air resistance (flattens out)
Dashed line = ideal free fall (straight, no limit)
```

With air resistance, the object approaches terminal velocity asymptotically — getting closer and closer but technically never quite reaching it.

---

## 10. Worked Examples

### Worked Example 3: Dropping a Ball from a Building

**Problem:** A stone is dropped from the top of a building. It hits the ground 4 seconds later. 
(a) How tall is the building?  
(b) What is the stone's velocity just before it hits?

**Given:** u = 0, t = 4 s, g = 9.8 m/s² (use 10 m/s²)

**Part (a):** Find the distance fallen (= height of building).

Use s = ut + (1/2)at² with downward positive (a = +g for this approach):
```
s = (0)(4) + (1/2)(10)(4²)
s = 0 + (1/2)(10)(16)
s = 0 + 80
s = 80 metres
```

The building is **80 metres** tall.

**Part (b):** Find velocity just before impact.

Use v = u + at (downward positive):
```
v = 0 + (10)(4)
v = 40 m/s  (downward)
```

Or about 144 km/h. This illustrates why falling from great heights is so dangerous.

**Sanity check with v² = u² + 2as:**
```
v² = 0 + 2(10)(80) = 1600
v = 40 m/s ✓
```

---

### Worked Example 4: Ball Thrown Straight Up from a Rooftop

**Problem:** From the top of a 20 m building, a ball is thrown straight up at 15 m/s. Find:  
(a) The maximum height above the ground.  
(b) The time for the ball to reach maximum height.  
(c) The total time until the ball hits the ground at the base.  
(d) The speed at which it hits the ground.

**Setup:** Let upward be positive. The ball starts at height 20 m above ground.

Taking the rooftop as the origin (s = 0 at the throw point):
- u = +15 m/s
- a = -10 m/s²
- The ground is at s = -20 m (20 m below the origin)

```
                                *  <- maximum height (s = s_max)
                              /   \
                            /       \
                          /           \
         * (s = 0)       /             \
    ROOFTOP  |----------/               \
    20 m     |                           \
    above    |                            \
    ground   |                             \
             |                              *
             |_______________________________  <- ground (s = -20)
```

**Part (a):** Maximum height above the throw point

At the top, v = 0:
```
v² = u² + 2as_max
0 = (15)² + 2(-10)(s_max)
0 = 225 - 20 × s_max
s_max = 225/20 = 11.25 m
```

Maximum height above the ground = 20 + 11.25 = **31.25 metres**

**Part (b):** Time to reach maximum height

```
v = u + at_top
0 = 15 + (-10)(t_top)
10 × t_top = 15
t_top = 1.5 seconds
```

**Part (c):** Total time until ball hits ground

The ball hits the ground when s = -20 (our origin is the rooftop):
```
s = ut + (1/2)at²
-20 = 15t + (1/2)(-10)t²
-20 = 15t - 5t²
5t² - 15t - 20 = 0
t² - 3t - 4 = 0
(t - 4)(t + 1) = 0
t = 4  or  t = -1
```

We discard t = -1 (negative time has no physical meaning). The ball hits the ground at **t = 4 seconds**.

**Part (d):** Speed at impact

```
v = u + at
v = 15 + (-10)(4)
v = 15 - 40
v = -25 m/s
```

Speed at impact = **25 m/s** (downward). That is about 90 km/h.

---

### Worked Example 5: Finding g from a Drop Experiment

**Problem:** A student drops a ball and times it with a stopwatch. The ball falls 19.6 m and takes 2.00 seconds. Calculate the value of g from this experiment.

**Given:** u = 0, s = 19.6 m (downward), t = 2.00 s

Use s = ut + (1/2)gt² (taking downward as positive):
```
19.6 = (0)(2) + (1/2)(g)(2²)
19.6 = 0 + (1/2)(g)(4)
19.6 = 2g
g = 19.6/2
g = 9.8 m/s²
```

The experiment gives g = **9.8 m/s²** — exactly the accepted value!

---

### Worked Example 6: The Coin and the Feather

**Conceptual question:** A coin and a feather are dropped inside a vacuum chamber (all air removed). Which hits the bottom first?

**Answer:** They hit at exactly the same time.

In a vacuum, there is no air resistance. Both objects experience only gravity. Since both have the same acceleration g downward regardless of their mass, they fall together.

This is Galileo's result, confirmed every time it is demonstrated. The mass of the object does not appear in any of the free fall equations (s = (1/2)gt², v = gt, etc.) — mass has been divided out because it appears in both the gravitational force and the inertia, and they cancel.

---

### Worked Example 7: Catching Your Own Throw

**Problem:** You throw a ball straight up at 30 m/s and catch it at the same height.  
(a) How high does it go?  
(b) How long is it in the air?  
(c) What is its speed when you catch it?

Using g = 10 m/s², upward positive:

**Part (a):** Maximum height
```
s_max = u²/(2g) = (30)²/(2 × 10) = 900/20 = 45 metres
```

**Part (b):** Total time in air
```
t_top = u/g = 30/10 = 3 seconds to reach top
Total time = 2 × 3 = 6 seconds
```

**Part (c):** Speed on return

By symmetry, the speed when caught = launch speed = **30 m/s** (but now downward).

Verify with v = u + at:
```
v = 30 + (-10)(6) = 30 - 60 = -30 m/s
Speed = 30 m/s ✓
```

---

### Worked Example 8: Two Balls, One Dropped, One Thrown Up

**Problem:** Ball A is dropped from rest from the top of a cliff. At the same moment, Ball B is thrown straight up from the bottom of the cliff at 20 m/s. The cliff is 30 m high. Do they meet? If so, where and when?

**Setup:**
- Ball A: starts at height 30 m, u = 0, falls under gravity.
- Ball B: starts at height 0 m, u = +20 m/s, rises under gravity.
- Taking ground as origin, upward positive, g = 10 m/s².

Height of Ball A at time t:
```
h_A = 30 - (1/2)(10)t² = 30 - 5t²
```

Height of Ball B at time t:
```
h_B = 20t - (1/2)(10)t² = 20t - 5t²
```

They meet when h_A = h_B:
```
30 - 5t² = 20t - 5t²
30 = 20t
t = 1.5 seconds
```

Height at meeting:
```
h_B = 20(1.5) - 5(1.5)²
h_B = 30 - 5(2.25)
h_B = 30 - 11.25
h_B = 18.75 metres above the ground
```

They meet at **18.75 m above the ground** at **t = 1.5 seconds**.

---

## Comprehensive ASCII Diagram: Height vs. Time for a Thrown Ball

This diagram shows the complete journey of a ball thrown upward at 20 m/s from the ground (using g = 10 m/s²).

```
Height
(metres)
  20 |                 *
     |              *     *
     |           *           *
  15 |         *               *
     |        *                 *
     |       *                   *
  10 |      *                     *
     |     *                       *
     |    *                         *
   5 |   *                           *
     |  *                             *
     | *                               *
   0 |*_________________________________*____
     0    0.5   1.0   1.5   2.0   2.5   3.0  4.0
                          time (seconds)
     
     Ball thrown up at 20 m/s, g = 10 m/s²
     t = 0: thrown from ground (height = 0)
     t = 2: reaches maximum height of 20 m
     t = 4: returns to ground
     
     Shape is a PARABOLA — 
     symmetric about the peak at t = 2s

Corresponding VELOCITY:

v (m/s)
  +20 |*
      | *
  +10 |  *
      |   *
    0 |----*---- <- top, v = 0 (t = 2s)
      |     *
  -10 |      *
      |       *
  -20 |        *
      |__________________________
      0    1    2    3    4   time (s)
      
      Velocity decreases linearly from +20 to -20.
      Slope of this line = -g = -10 m/s²
      
Corresponding ACCELERATION:

a (m/s²)
      |
  +10 |   (never here during free fall)
      |
    0 |_________________________________
      |
  -10 |*-*-*-*-*-*-*-*-*-*-*-*-*-*-*-*  <- constant -10 m/s²
      |   CONSTANT throughout the flight
      |
      0    1    2    3    4   time (s)
      
      Acceleration is the same at every point,
      even at the top where velocity = 0.
      At the top, v = 0 but a = -g ≠ 0.
      The ball is NOT "weightless" at the top
      in any real sense — gravity is still acting.
```

---

### Key Insight: Zero Velocity Does Not Mean Zero Acceleration

Many students make this mistake: they think that at the top of a thrown ball's path, since the velocity is zero, the acceleration must also be zero. The ball is momentarily still — so surely nothing is happening?

Wrong. Completely wrong.

Acceleration is the rate of change of velocity. Even when velocity is zero, it can still be changing. At the very top:

- The velocity is zero.
- But the velocity is changing at a rate of -9.8 m/s² — it is about to become negative (downward).
- Gravity is pulling on the ball just as hard as at any other point.

Think of it this way: throw a ball up and watch it at the exact moment it reaches the top. It is momentarily still in your hand — but if you let it go, it immediately starts falling. Gravity never switched off.

---

### Common Mistakes to Avoid

Here is a list of the most frequent errors students make in free fall problems:

| Mistake | The Correct Thinking |
|---|---|
| "At the top, acceleration = 0" | Acceleration = -g throughout. Only velocity = 0 at the top. |
| "Heavy objects fall faster" | In vacuum, all objects have the same acceleration g. |
| "Upward velocity means going up forever" | Gravity constantly reduces the upward velocity. |
| "Dropped object has zero total velocity" | It starts with zero velocity but gains speed immediately. |
| "Free fall only means falling downward" | An object thrown upward is also in free fall (only gravity acting). |
| Forgetting the sign convention | Always state: which direction is positive? Stick to it throughout. |
| Using g = 9.8 but forgetting direction | a = -9.8 m/s² (upward positive) or +9.8 m/s² (downward positive) |

---

### Comparing Free Fall on Different Planets

We close with a fun comparison. If you dropped a ball from 80 metres on different planets, how long would it take to hit the ground?

Using t = sqrt(2s/g):

| Planet/Moon | g (m/s²) | Time to fall 80 m |
|---|---|---|
| Moon | 1.62 | 9.9 seconds |
| Mars | 3.72 | 6.6 seconds |
| Earth | 9.81 | 4.0 seconds |
| Saturn | 10.44 | 3.9 seconds |
| Jupiter | 24.79 | 2.5 seconds |
| Sun | 274 | 0.76 seconds |

On the Moon, the ball would take nearly 10 seconds to fall 80 metres — drifting down in slow motion by our standards. On Jupiter, it would slam down in 2.5 seconds. On the surface of the Sun, it would cover 80 metres in less than a second.

---

## 11. Summary

- **Free fall** is motion under gravity alone, with no air resistance. The only force acting is gravity.

- **Galileo proved** that all objects fall at the same rate regardless of mass (in the absence of air resistance). This contradicted 2,000 years of Aristotelian dogma.

- The **acceleration due to gravity** near Earth's surface is g = 9.8 m/s² ≈ 10 m/s², directed downward. It is the same for all objects.

- Free fall is just constant acceleration with a = g = 9.8 m/s² downward. The SUVAT equations apply directly: just substitute a = -g (upward positive) or a = +g (downward positive).

- **Object dropped from rest** (u = 0): gains 9.8 m/s of downward speed every second; distance fallen = (1/2)gt².

- **Object thrown upward**: decelerates at g on the way up, stops momentarily at the peak (v = 0), then accelerates downward at g on the way down.

- **Symmetry of free fall**: time going up equals time coming down; speed at any given height is the same on the way up and the way down.

- At the **peak** of a thrown ball's path: velocity = 0, but acceleration = -g (gravity is still acting).

- **Weightlessness** is not the absence of gravity. It is being in free fall with nothing supporting you. Astronauts in orbit are in continuous free fall (and feel weightless), even though Earth's gravity is nearly as strong there as at the surface.

- **Air resistance** creates a drag force opposing motion. As a falling object speeds up, drag increases until it balances gravity — the object then falls at constant **terminal velocity**. Parachutes exploit this: large area = large drag = low terminal velocity.

- The mass of an object does NOT appear in the equations for free fall. Mass cancels out perfectly — which is why all objects fall at the same rate.

---

## 12. Key Equations

```
ACCELERATION DUE TO GRAVITY
g = 9.8 m/s²  (≈ 10 m/s² for easy calculations)
a = -g  (upward positive convention)
a = +g  (downward positive convention)

FREE FALL SUVAT EQUATIONS
(using upward positive, so a = -g)

v = u - gt
s = ut - (1/2)gt²
v² = u² - 2gs
s = (u + v)/2 × t

OBJECT DROPPED FROM REST (u = 0)
v = -gt             <- velocity after time t
s = -(1/2)gt²       <- displacement (downward, so negative)
distance fallen = (1/2)gt²

OBJECT THROWN UPWARD WITH SPEED u
Time to reach top:         t_top = u / g
Maximum height:            s_max = u² / (2g)
Total flight time:         T = 2u / g = 2 × t_top
Speed at same height:      same going up as coming down

WEIGHT (force of gravity on an object)
W = m × g

TERMINAL VELOCITY (conceptual)
Reached when: Drag force = Weight
At terminal velocity: net force = 0, acceleration = 0
```

---

*End of Chapter 09: Free Fall*

*Next chapter: Chapter 10 — Projectile Motion (combining horizontal and vertical free fall)*
