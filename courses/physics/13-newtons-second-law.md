# Chapter 13: Newton's Second Law — F = ma

> **"F = ma. Three letters. The most powerful equation in classical physics. It tells you everything about how things move when forces act on them."**

---

## Table of Contents

1. [What Is This Chapter About?](#1-what-is-this-chapter-about)
2. [Recap: Newton's First Law and Net Force](#2-recap-newtons-first-law-and-net-force)
3. [Newton's Second Law — The Official Statement](#3-newtons-second-law--the-official-statement)
4. [Breaking Down the Variables](#4-breaking-down-the-variables)
5. [Defining the Newton (Unit of Force)](#5-defining-the-newton-unit-of-force)
6. [Three Ways to Use F = ma](#6-three-ways-to-use-f--ma)
7. [The Relationships: Doubling and Halving](#7-the-relationships-doubling-and-halving)
8. [Free Body Diagrams — The Essential Tool](#8-free-body-diagrams--the-essential-tool)
9. [The Five Common Forces You Will Draw](#9-the-five-common-forces-you-will-draw)
10. [Worked Example 1: Box on a Frictionless Surface](#10-worked-example-1-box-on-a-frictionless-surface)
11. [Worked Example 2: Box With Friction](#11-worked-example-2-box-with-friction)
12. [Worked Example 3: Accelerating a Car](#12-worked-example-3-accelerating-a-car)
13. [Worked Example 4: The Elevator Problem — Apparent Weight](#13-worked-example-4-the-elevator-problem--apparent-weight)
14. [Weight vs Normal Force — What Is the Difference?](#14-weight-vs-normal-force--what-is-the-difference)
15. [Forces in Two Directions — Horizontal and Vertical](#15-forces-in-two-directions--horizontal-and-vertical)
16. [Connecting F = ma to the Kinematic Equations (SUVAT)](#16-connecting-f--ma-to-the-kinematic-equations-suvat)
17. [Common Mistakes to Avoid](#17-common-mistakes-to-avoid)
18. [Practice Problems](#18-practice-problems)
19. [Summary](#19-summary)
20. [Key Equations](#20-key-equations)

---

## 1. What Is This Chapter About?

In Chapter 12 we looked at Newton's First Law — an object at rest stays at rest and an object in motion stays in motion, unless a **net force** acts on it. That chapter was about what happens when forces are balanced.

This chapter is about what happens when forces are **unbalanced**.

Specifically: if you push something, how fast does it speed up? If something heavy and something light experience the same push, which one accelerates more? How do engineers calculate the engine force needed to get a car to highway speed?

All of those questions are answered by a single equation:

```
F_net = m × a
```

This is **Newton's Second Law of Motion**, and it is arguably the most useful equation in all of classical physics. Engineers use it to design rockets, bridges, and car brakes. Doctors use it to understand forces on bones and joints. Animators use it to make movies look realistic.

By the end of this chapter you will be able to:

- State Newton's Second Law in words and as an equation
- Identify what F, m, and a each mean and what units they use
- Draw a Free Body Diagram for any object
- Solve problems involving forces, mass, and acceleration
- Calculate apparent weight in an accelerating elevator
- Connect the result of F = ma directly to SUVAT kinematics

---

## 2. Recap: Newton's First Law and Net Force

Before diving into the Second Law, let's cement two ideas from earlier.

**Idea 1: Inertia.**
Every object resists changes to its motion. A heavy boulder is harder to start moving and harder to stop than a tennis ball. This resistance is called **inertia**, and it is measured by the object's **mass**.

**Idea 2: Net Force.**
When multiple forces act on an object at once, what matters is their **vector sum** — the single combined force. We call this the **net force** (also written F_net or ΣF, where the Greek letter Sigma means "sum of").

```
If forces point the same direction, add them.
If they point opposite directions, subtract them.

Example:
   [30 N →]  [10 N ←]

   Net force = 30 - 10 = 20 N to the right
```

When F_net = 0, the object either stays still or moves at constant velocity (Newton's First Law).

When F_net ≠ 0 (unbalanced), the object **accelerates**. That is what Newton's Second Law describes.

---

## 3. Newton's Second Law — The Official Statement

**Newton's Second Law**: The acceleration of an object is directly proportional to the net force acting on it, and inversely proportional to its mass. The acceleration is in the same direction as the net force.

In equation form:

```
F_net = m × a
```

Or written with the summation symbol:

```
ΣF = m × a
```

Where:
- F_net = the total (net) force on the object, in Newtons (N)
- m = the mass of the object, in kilograms (kg)
- a = the acceleration of the object, in metres per second squared (m/s²)

This one equation tells you:

- A bigger force produces a bigger acceleration (they are **directly proportional**)
- A bigger mass produces a smaller acceleration for the same force (they are **inversely proportional**)

Think about it physically. Push a shopping cart (light). It accelerates easily. Push a loaded truck (heavy). It barely moves. Same force, very different accelerations. That is Newton's Second Law in action.

---

## 4. Breaking Down the Variables

Let's look carefully at each variable.

### F_net — Net Force

Force is a **vector** — it has both a magnitude (size) and a direction. F_net is the total force after you add all individual forces acting on the object, accounting for direction.

Units: **Newtons (N)**

You can also have forces in multiple directions. In that case you will work with separate equations for horizontal and vertical:

```
Horizontal:  ΣF_x = m × a_x
Vertical:    ΣF_y = m × a_y
```

### m — Mass

Mass is the amount of **matter** in an object. It is also the measure of inertia — how much the object resists being accelerated.

Units: **kilograms (kg)**

Important: mass is NOT the same as weight. Mass is a fixed property of the object. Weight depends on gravity and changes if you travel to the Moon. (We will cover this in Section 14.)

### a — Acceleration

Acceleration is the rate of change of velocity.

Units: **metres per second squared (m/s²)**

Acceleration can be:
- Positive (speeding up in the chosen positive direction)
- Negative (slowing down, or speeding up in the negative direction)
- Zero (constant velocity — no net force)

The direction of acceleration is always the same as the direction of F_net. If the net force points right, the object accelerates to the right.

---

## 5. Defining the Newton (Unit of Force)

The **Newton** (symbol: N) is the SI unit of force. It is defined directly from F = ma:

```
1 Newton = the force required to accelerate a mass of 1 kg at 1 m/s²

1 N = 1 kg × 1 m/s²
1 N = 1 kg·m/s²
```

To get a feel for it:
- 1 N ≈ the weight of a small apple (about 100 grams on Earth)
- 10 N ≈ the weight of a 1 kg bag of sugar (because g ≈ 10 m/s²)
- 600 N ≈ the weight of an average adult (about 60 kg)
- 1,000,000 N (1 MN) ≈ the thrust of a space rocket engine

Always check your units when using F = ma:

```
Force in Newtons = Mass in kg × Acceleration in m/s²

     [N]          =     [kg]    ×      [m/s²]
  [kg·m/s²]       =     [kg]    ×      [m/s²]   ✓
```

If your mass is in grams or your acceleration is in cm/s², convert first.

---

## 6. Three Ways to Use F = ma

F = ma is one equation with three variables. If you know any two, you can find the third.

### Form 1 — Find Force (given mass and acceleration)

```
F = m × a
```

Use when: you know how fast an object accelerates and you want to find the force that causes it.

Example: A 3 kg ball accelerates at 5 m/s². What force acted on it?
```
F = m × a = 3 × 5 = 15 N
```

### Form 2 — Find Acceleration (given force and mass)

```
a = F / m
```

Use when: you know the force applied and the mass, and you want to know how fast it will speed up.

Example: A 10 N force acts on a 2 kg object. What is its acceleration?
```
a = F / m = 10 / 2 = 5 m/s²
```

### Form 3 — Find Mass (given force and acceleration)

```
m = F / a
```

Use when: you observe the acceleration produced by a known force, and you want to find the object's mass.

Example: A 24 N force produces an acceleration of 6 m/s². What is the mass?
```
m = F / a = 24 / 6 = 4 kg
```

---

## 7. The Relationships: Doubling and Halving

Newton's Second Law contains two clean proportional relationships that are worth memorising.

### Relationship 1: Force and Acceleration are Directly Proportional (mass fixed)

```
a ∝ F    (when m is constant)
```

This means:
- Double the force → double the acceleration
- Triple the force → triple the acceleration
- Halve the force → halve the acceleration

```
EXPERIMENT: Same 2 kg box, different forces

  Force (N)  |  Mass (kg)  |  Acceleration (m/s²)
  -----------|-------------|---------------------
       5     |      2      |       2.5
      10     |      2      |       5.0      ← doubled force, doubled acceleration
      20     |      2      |      10.0      ← doubled again, doubled again
```

### Relationship 2: Mass and Acceleration are Inversely Proportional (force fixed)

```
a ∝ 1/m    (when F is constant)
```

This means:
- Double the mass → half the acceleration
- Triple the mass → one-third the acceleration
- Halve the mass → double the acceleration

```
EXPERIMENT: Same 12 N force, different masses

  Force (N)  |  Mass (kg)  |  Acceleration (m/s²)
  -----------|-------------|---------------------
      12     |      1      |      12.0
      12     |      2      |       6.0      ← doubled mass, halved acceleration
      12     |      4      |       3.0      ← doubled mass again, halved again
      12     |      6      |       2.0
```

This inverse relationship is why a heavy truck needs a much more powerful engine than a small car just to achieve the same acceleration.

---

## 8. Free Body Diagrams — The Essential Tool

A **Free Body Diagram** (FBD) is the most important problem-solving tool in dynamics. It is a drawing that shows every force acting on a single object as an arrow pointing in the direction of that force.

**Rules for a good Free Body Diagram:**

1. Draw the object as a simple box or dot (it does not need to look realistic)
2. Draw EVERY force as an arrow pointing away from the object
3. Label each arrow with the force's name and symbol
4. The length of an arrow represents the approximate magnitude of the force
5. Include only forces that act ON the chosen object — not forces the object exerts on other things
6. Choose a positive direction and label it

**The procedure:**

```
Step 1: Identify the object you are analyzing
Step 2: List every force acting ON it
Step 3: Draw the object as a box or dot
Step 4: Draw each force as a labelled arrow
Step 5: Add up forces in each direction to find F_net
Step 6: Apply F_net = m × a
```

Here is a basic FBD for a book sitting on a table:

```
             N (Normal force, upward)
             ↑
             |
    ---------+--------
    |                |   ← The book
    ---------+--------
             |
             ↓
             W = mg (Weight, downward)

Since the book is not accelerating: F_net = 0
So N = W = mg
```

Here is an FBD for a box being pushed to the right on a frictionless surface:

```
             N
             ↑
             |
    ---------+--------
    |                |----→  F_A (Applied force)
    ---------+--------
             |
             ↓
             W = mg

Horizontal:  F_net = F_A  →  a = F_A / m
Vertical:    N - W = 0    →  N = mg (no vertical acceleration)
```

---

## 9. The Five Common Forces You Will Draw

You will encounter the same five forces in almost every FBD at this level.

### Force 1: Weight (W or F_g)

Weight is the gravitational force pulling the object downward toward the centre of the Earth.

```
W = m × g

where:
  m = mass in kg
  g = acceleration due to gravity = 9.8 m/s² (we often round to 10 m/s²)

Direction: always straight DOWN (toward Earth's centre)
```

A 5 kg object has weight:  W = 5 × 9.8 = 49 N ≈ 50 N

### Force 2: Normal Force (N)

The **normal force** is the contact force a surface exerts on an object perpendicular (at 90°) to the surface. It prevents the object from sinking into the surface.

```
Direction: perpendicular to the surface, pointing AWAY from the surface

On a flat horizontal surface: N points straight UP
On a ramp: N points at 90° to the ramp surface
```

Normal force is NOT always equal to weight — it equals weight only when the object is on a flat surface and not accelerating vertically.

### Force 3: Friction Force (f)

Friction opposes motion (or attempted motion) and acts parallel to the surface.

```
Direction: opposite to the direction of motion (or attempted motion)
```

We will study friction in depth in Chapter 15.

### Force 4: Applied Force (F_A or F_push)

Any force you push or pull an object with.

```
Direction: in whatever direction the push/pull acts
```

### Force 5: Tension (T)

The pulling force transmitted through a rope, string, or cable.

```
Direction: along the rope, AWAY from the object (ropes can only pull, not push)
```

```
SUMMARY OF FORCES IN AN FBD:

            N (↑ perpendicular to surface)
            |
            |
 f (←) ----[BOX]---- F_A (→)
            |
            |
            W (↓ = mg)

Note: friction f points left because the applied force F_A pushes right.
```

---

## 10. Worked Example 1: Box on a Frictionless Surface

**Problem:** A 5 kg box sits on a perfectly frictionless horizontal surface. A horizontal force of 20 N is applied to it. Find the acceleration of the box.

**Step 1: Draw the Free Body Diagram**

```
            N
            ↑
            |
   ---------+--------
   |   5 kg          |----→  F_A = 20 N
   ---------+--------
            |
            ↓
         W = mg

(No friction force because the surface is frictionless)
(Positive direction: →)
```

**Step 2: Write the equations for each direction**

Vertical direction (y-axis):
```
ΣF_y = 0    (no vertical acceleration)
N - W = 0
N = W = mg = 5 × 9.8 = 49 N
```

Horizontal direction (x-axis):
```
ΣF_x = m × a_x
F_A = m × a
20 = 5 × a
```

**Step 3: Solve for acceleration**

```
a = F_A / m
a = 20 / 5
a = 4 m/s²
```

**Answer:** The box accelerates at **4 m/s²** to the right.

**Sense check:** 20 N pushing a 5 kg object. Does 4 m/s² feel right? If we applied 50 N to a 5 kg object we'd get 10 m/s², and if we applied 5 N we'd get 1 m/s². So 20 N giving 4 m/s² fits the pattern. ✓

---

## 11. Worked Example 2: Box With Friction

**Problem:** A 10 kg box is pushed along a horizontal floor with an applied force of 30 N to the right. The friction force acting on the box is 10 N. Find the net force and the acceleration of the box.

**Step 1: Draw the Free Body Diagram**

```
              N
              ↑
              |
   f = 10 N ←---------+----------→  F_A = 30 N
              | 10 kg  |
              ---------
              |
              ↓
           W = mg = 10 × 9.8 = 98 N

(Positive direction: → to the right)
```

**Step 2: Find the net force in the horizontal direction**

```
ΣF_x = F_A - f
ΣF_x = 30 - 10
ΣF_x = 20 N   (to the right)
```

The friction force opposes the motion, so it acts to the LEFT (negative direction), which is why we subtract it.

**Step 3: Apply Newton's Second Law**

```
F_net = m × a
20 = 10 × a
a = 20 / 10
a = 2 m/s²
```

**Answer:** The net force is **20 N to the right** and the acceleration is **2 m/s² to the right**.

**Compare with Example 1:** Without friction, the 30 N force on a 10 kg box would give a = 30/10 = 3 m/s². With friction eating up 10 N, we only get 2 m/s². Friction always reduces acceleration.

---

## 12. Worked Example 3: Accelerating a Car

**Problem:** What net force is needed to accelerate a 1200 kg car from rest to a velocity of 18 m/s in 6 seconds?

**Step 1: Find the required acceleration**

We are not given acceleration directly, but we know:
- Initial velocity: u = 0 m/s (starts from rest)
- Final velocity: v = 18 m/s
- Time: t = 6 s

Using the kinematic equation:
```
v = u + a × t
18 = 0 + a × 6
a = 18 / 6
a = 3 m/s²
```

**Step 2: Draw the Free Body Diagram**

```
              N
              ↑
              |
    ---------[CAR]----------→  F_net = ?
              |
              ↓
           W = mg
           = 1200 × 9.8
           = 11,760 N

(Simplified — not showing individual engine force and friction,
 just the net forward force)
```

**Step 3: Apply Newton's Second Law**

```
F_net = m × a
F_net = 1200 × 3
F_net = 3600 N
```

**Answer:** The car requires a **net force of 3600 N** to achieve that acceleration.

**Context:** 3600 N = 3.6 kN. A typical small car engine produces around 80–150 kN of maximum force — but much of that is used to overcome air resistance and friction. The NET force available is what actually causes acceleration.

---

## 13. Worked Example 4: The Elevator Problem — Apparent Weight

This is one of the most famous problems in introductory physics. It reveals the difference between **true weight** and **apparent weight**.

**Setup:** A 70 kg person stands on a bathroom scale inside an elevator. The elevator accelerates upward at 2 m/s².

*What does the scale read?*

First, a crucial clarification: a bathroom scale does NOT directly measure your weight (gravity). It measures the **normal force** that the scale pushes up on your feet. When you are stationary, N = W = mg, so the reading equals your weight. But when you accelerate, N ≠ W.

**Step 1: Draw the Free Body Diagram for the person**

```
              N (Normal force from scale — this is what scale READS)
              ↑
              |
         -----+-----
         | Person  |      a = 2 m/s² UPWARD (acceleration of elevator)
         -----+-----
              |
              ↓
           W = mg = 70 × 9.8 = 686 N

Positive direction: ↑ (upward)
```

**Step 2: Apply Newton's Second Law in the vertical direction**

The person accelerates upward at 2 m/s², so the net upward force must equal m × a:

```
ΣF_y = m × a
N - W = m × a
N - mg = m × a
N = mg + m × a
N = m(g + a)
N = 70 × (9.8 + 2)
N = 70 × 11.8
N = 826 N
```

**Answer:** The scale reads **826 N**, which is equivalent to 84.3 kg on Earth.

The person feels **heavier** than normal. This is called increased **apparent weight**.

**What if the elevator decelerates (slows down while going up)?**

Decelerating upward means acceleration is DOWNWARD (negative in our convention):

```
N = m(g + a)   where a = -2 m/s² (downward)
N = 70 × (9.8 - 2)
N = 70 × 7.8
N = 546 N
```

The scale reads 546 N. The person feels **lighter** than normal.

**What if the cable snaps? (Free fall: a = -9.8 m/s²)**

```
N = m(g + a)
N = 70 × (9.8 - 9.8)
N = 70 × 0
N = 0 N
```

The scale reads zero. The person is **weightless**. This is exactly what astronauts in orbit experience — they are in constant free fall around the Earth!

**Summary of Elevator Cases:**

```
Elevator State                 | Acceleration   | Scale Reading
-------------------------------|----------------|------------------
Stationary                     | a = 0          | N = mg   (normal)
Accelerating upward            | a = +2 m/s²    | N > mg   (heavier)
Decelerating upward            | a = -2 m/s²    | N < mg   (lighter)
Accelerating downward          | a = -2 m/s²    | N < mg   (lighter)
Decelerating downward          | a = +2 m/s²    | N > mg   (heavier)
Free fall (cable snapped)      | a = -g         | N = 0    (weightless)
```

The pattern: whenever the elevator accelerates in the same direction as gravity (downward), you feel lighter. Whenever it accelerates against gravity (upward), you feel heavier.

---

## 14. Weight vs Normal Force — What Is the Difference?

Students often confuse weight and normal force. Let's clarify once and for all.

**Weight (W = mg):**
- The gravitational force Earth exerts on the object
- Acts downward, toward Earth's centre
- Is a fixed property for a given mass and location
- Does NOT depend on what surface the object is on
- W = 5 kg × 9.8 = 49 N for a 5 kg object, regardless of where it is resting

**Normal Force (N):**
- The contact force a surface exerts on the object
- Acts perpendicular (normal) to the surface
- Adjusts to prevent the object from passing through the surface
- Depends on the situation (acceleration, other forces, surface angle)
- Can be greater than, equal to, or less than W

```
FLAT SURFACE, STATIONARY:

         N ↑
         |
    [----+----]
         |
         ↓ W = mg

ΣF_y = 0 → N = mg
N equals W. But this is a special case!


FLAT SURFACE, EXTRA DOWNWARD FORCE:

    Person pushing DOWN on box
         ↓ F_extra
         N ↑
         |
    [----+----]
         |
         ↓ W = mg

ΣF_y = 0 → N = mg + F_extra
N is GREATER than mg.


FLAT SURFACE, ACCELERATING UPWARD:

         N ↑
         |
    [----+----]  ↑ a
         |
         ↓ W = mg

ΣF_y = ma → N - mg = ma → N = m(g + a)
N is GREATER than mg.
```

**Rule of thumb:** N = mg only for an object on a flat, horizontal surface with no vertical acceleration and no other vertical forces. In every other situation, calculate N separately using Newton's Second Law.

---

## 15. Forces in Two Directions — Horizontal and Vertical

When forces act in both horizontal and vertical directions at once, you handle them separately.

**The Key Principle:** Newton's Second Law applies independently in each direction.

```
Horizontal:   ΣF_x = m × a_x
Vertical:     ΣF_y = m × a_y
```

If an object moves only horizontally (like a box sliding on a floor), then a_y = 0, so ΣF_y = 0. You use the vertical equation just to find N, then use N to find friction, then use friction in the horizontal equation.

**Example: Box being pulled at an upward angle**

Imagine a box pulled by a rope at 30° above the horizontal with force T = 40 N. The box is on a flat floor with mass 8 kg.

```
              T (at 30° above horizontal)
             /
            /  30°
   --------[BOX]---------
```

Step 1: Resolve T into components:

```
T_x = T × cos(30°) = 40 × 0.866 = 34.6 N  (horizontal, →)
T_y = T × sin(30°) = 40 × 0.5   = 20 N    (vertical, ↑)
```

Step 2: FBD with components:

```
              N + T_y = N + 20 N
              ↑
              |
T_x = 34.6 N →[BOX]
              |
              ↓
           W = mg = 8 × 9.8 = 78.4 N
```

Step 3: Vertical equation (no vertical acceleration):

```
ΣF_y = 0
N + T_y - W = 0
N + 20 - 78.4 = 0
N = 58.4 N
```

Note: the upward component of T reduces the normal force. This means friction would also be reduced (because friction depends on N).

Step 4: Horizontal equation (assuming no friction for simplicity):

```
ΣF_x = m × a_x
T_x = m × a
34.6 = 8 × a
a = 34.6 / 8
a = 4.3 m/s²
```

---

## 16. Connecting F = ma to the Kinematic Equations (SUVAT)

Newton's Second Law and the kinematic equations (from Chapter 08) are two halves of a complete system. They work together like this:

```
STEP 1: Use Newton's Second Law (Forces → Acceleration)

   F_net = m × a    →    a = F_net / m

STEP 2: Use SUVAT (Acceleration → Motion details)

   v = u + a × t         (find final velocity)
   s = u × t + ½ × a × t²  (find displacement)
   v² = u² + 2 × a × s   (find velocity without time)
```

**Combined Example:**

A 2 kg toy car is pushed from rest with a net force of 6 N. How far does it travel in 4 seconds?

Step 1 — Find acceleration:
```
a = F_net / m = 6 / 2 = 3 m/s²
```

Step 2 — Use SUVAT to find displacement:
```
Known: u = 0, a = 3 m/s², t = 4 s
s = u × t + ½ × a × t²
s = 0 × 4 + ½ × 3 × 4²
s = 0 + ½ × 3 × 16
s = 24 m
```

The car travels **24 metres** in 4 seconds.

This two-step approach — forces first, then kinematics — solves a huge range of real-world problems.

---

## 17. Common Mistakes to Avoid

### Mistake 1: Forgetting that F means NET force

The equation is F_**net** = ma, not just any old F. If a 30 N force and a 10 N friction force both act, F_net = 20 N — not 30 N.

```
WRONG:  a = 30 / 10 = 3 m/s²    (ignoring friction)
RIGHT:  a = (30 - 10) / 10 = 2 m/s²
```

### Mistake 2: Confusing mass and weight

Mass is in **kg**. Weight is in **N**. Never plug 70 kg in for a force — convert first: W = mg = 70 × 9.8 = 686 N.

### Mistake 3: Assuming N = mg always

Normal force only equals mg on a flat surface with no vertical acceleration and no other vertical forces. In elevators, on ramps, or with angled pushes, N will be different.

### Mistake 4: Wrong direction for friction

Friction always opposes the **direction of motion** (or attempted motion). If the object moves right, friction points left.

### Mistake 5: Not drawing an FBD

Skipping the Free Body Diagram almost always leads to errors. Always draw it, even for simple problems. It takes 30 seconds and prevents mistakes.

### Mistake 6: Adding forces as numbers without considering direction

Forces are vectors. A 10 N force to the right and a 10 N force to the left give F_net = 0, not 20 N.

---

## 18. Practice Problems

Try these problems on your own before checking answers.

**Problem 1 (Easy)**
A net force of 15 N acts on a 3 kg object. What is its acceleration?

**Problem 2 (Easy)**
A 1500 kg car accelerates at 2.5 m/s². What is the net force on it?

**Problem 3 (Medium)**
A 12 kg box is pushed with 60 N on a surface with 12 N of friction. Find the net force and acceleration.

**Problem 4 (Medium)**
A 500 kg elevator accelerates downward at 1.5 m/s². What does a 60 kg person's scale read inside it?
(Hint: acceleration is downward, so it opposes the normal force direction)

**Problem 5 (Medium)**
A force of 40 N produces an acceleration of 8 m/s² in an object. What is the object's mass?

**Problem 6 (Challenging)**
A 4 kg block is pushed from rest by a 20 N net force. Using both F = ma and SUVAT, find its velocity and displacement after 3 seconds.

---

**Answers:**

1. a = F/m = 15/3 = **5 m/s²**
2. F = ma = 1500 × 2.5 = **3750 N**
3. F_net = 60 - 12 = 48 N; a = 48/12 = **4 m/s²**
4. N = m(g - a) = 60 × (9.8 - 1.5) = 60 × 8.3 = **498 N** (lighter than normal)
5. m = F/a = 40/8 = **5 kg**
6. a = 20/4 = 5 m/s²; v = 0 + 5×3 = **15 m/s**; s = 0 + ½×5×9 = **22.5 m**

---

## 19. Summary

- **Newton's Second Law** states that the net force on an object equals its mass times its acceleration: F_net = m × a

- The **net force** is the vector sum of all forces acting on an object. You must account for direction when adding forces.

- The **Newton (N)** is the unit of force. 1 N = 1 kg·m/s². It is defined as the force that accelerates 1 kg at 1 m/s².

- F = ma has three forms:
  - F = m × a (find force)
  - a = F / m (find acceleration)
  - m = F / a (find mass)

- **Directly proportional**: doubling the net force (at constant mass) doubles the acceleration.

- **Inversely proportional**: doubling the mass (at constant force) halves the acceleration.

- A **Free Body Diagram** (FBD) is a drawing showing every force on an object as a labelled arrow. Always draw one before solving a dynamics problem.

- The five common forces: weight (W = mg, down), normal force (N, perpendicular to surface), friction (f, opposing motion), applied force (F_A), and tension (T).

- **Weight** (W = mg) is the gravitational force — always downward. **Normal force** (N) is the contact force from the surface — always perpendicular to it. They are equal only in the special case of a stationary object on a flat surface.

- In an elevator, **apparent weight** = N = m(g + a), where a is positive upward. Accelerating upward → heavier. Accelerating downward → lighter. Free fall → weightless.

- Newton's Second Law and the SUVAT kinematic equations work together: use F = ma to find acceleration, then use SUVAT to find velocity, displacement, or time.

- When forces act in two directions, apply ΣF_x = m × a_x and ΣF_y = m × a_y separately.

---

## 20. Key Equations

```
Newton's Second Law (vector form):
   F_net = m × a

In component form:
   ΣF_x = m × a_x
   ΣF_y = m × a_y

Three rearrangements:
   F = m × a          (find force)
   a = F / m          (find acceleration)
   m = F / a          (find mass)

Weight:
   W = m × g
   g = 9.8 m/s²  (use 10 m/s² for quick estimates)

Unit of force:
   1 N = 1 kg × 1 m/s²  =  1 kg·m/s²

Apparent weight in elevator:
   N = m × (g + a)    where a is positive upward, negative downward

Combined with SUVAT:
   Step 1:  a = F_net / m
   Step 2:  v = u + a × t
            s = u × t + ½ × a × t²
            v² = u² + 2 × a × s
```

---

*Next chapter: Chapter 14 — Newton's Third Law: Action and Reaction*
