# Chapter 03: Vectors and Scalars

> **"Direction is more important than speed. Many people are going nowhere fast."**
> — Unknown

---

## Table of Contents

1. [What Is a Physical Quantity?](#1-what-is-a-physical-quantity)
2. [Scalars: Quantities with Size Only](#2-scalars-quantities-with-size-only)
3. [Vectors: Quantities with Size AND Direction](#3-vectors-quantities-with-size-and-direction)
4. [The Core Difference: Why Direction Matters](#4-the-core-difference-why-direction-matters)
5. [Representing Vectors](#5-representing-vectors)
6. [Adding Vectors: The Graphical Method](#6-adding-vectors-the-graphical-method)
7. [Adding Vectors: The Algebraic Method](#7-adding-vectors-the-algebraic-method)
8. [Subtracting Vectors](#8-subtracting-vectors)
9. [Multiplying a Vector by a Scalar](#9-multiplying-a-vector-by-a-scalar)
10. [Resolving Vectors into Components](#10-resolving-vectors-into-components)
11. [Finding the Resultant: Magnitude and Direction](#11-finding-the-resultant-magnitude-and-direction)
12. [Unit Vectors i, j, and k](#12-unit-vectors-i-j-and-k)
13. [Practical Applications](#13-practical-applications)
14. [Summary](#14-summary)
15. [Key Equations](#15-key-equations)

---

## 1. What Is a Physical Quantity?

Every time you do something physical in the world — walk to school, pick up a bag, feel the warmth of sunlight — you are interacting with **physical quantities**. A physical quantity is anything in nature that can be **measured**.

Think about your morning routine:

- You wake up at **7:00 AM** — that is time, a physical quantity.
- The temperature outside is **22 degrees Celsius** — temperature is a physical quantity.
- You walk **500 metres** to the bus stop — that distance is a physical quantity.
- You push your door open — that push involves **force**, another physical quantity.

Physics is essentially the science of measuring these quantities and figuring out how they relate to each other. But not all physical quantities are the same kind of thing. Some of them are simple numbers. Others have an extra piece of information baked in — a direction.

This distinction — between quantities that have direction and those that do not — is one of the most important ideas in all of physics. Once you understand it, a huge amount of the physical world snaps into focus.

The two categories are called **scalars** and **vectors**.

---

## 2. Scalars: Quantities with Size Only

A **scalar** is a physical quantity that is fully described by just a single number (with a unit). That number tells you the **magnitude** — basically, how much or how big.

### 2.1 Understanding Magnitude

**Magnitude** just means the size or amount of something. When someone says "I have ₹500", the number 500 (with the unit rupees) completely describes how much money they have. You do not need to add "in which direction?" — the question makes no sense for money.

The same logic applies to many physical quantities.

### 2.2 Common Scalar Quantities

| Scalar Quantity | Symbol | Common Unit | Example |
|-----------------|--------|-------------|---------|
| Temperature | T | °C or K | 37°C (body temperature) |
| Mass | m | kg | 70 kg |
| Speed | v | m/s or km/h | 60 km/h |
| Time | t | seconds (s) | 45 minutes |
| Distance | d | metres (m) | 100 m |
| Energy | E | Joules (J) | 500 J |
| Power | P | Watts (W) | 100 W |
| Volume | V | m³ or litres | 2 litres |
| Density | ρ | kg/m³ | 1000 kg/m³ |
| Electric charge | q | Coulombs (C) | 1.6 × 10⁻¹⁹ C |

### 2.3 Everyday Examples of Scalars

**Temperature:** If someone says "it is 35°C today," you have everything you need to know about the temperature. You do not need to ask "in which direction?" It is just a number.

**Mass:** Your mass is, say, 60 kilograms. That single number says everything. (We will be very precise later — mass and weight are different, but both are scalars.)

**Speed:** Your car's speedometer reads 80 km/h. That single number tells you how fast you are moving. Note the word *speed* carefully — we will soon see that **velocity** is different from speed because velocity also includes direction.

**Time:** A race takes 9.58 seconds. Just a number.

**Energy:** A candy bar contains 200 kilocalories of energy. Just a number.

### 2.4 Worked Example: Adding Scalars

Scalars follow ordinary arithmetic. If you pour 500 ml of water from one bottle and 300 ml from another into a jug, you get:

```
Total volume = 500 ml + 300 ml = 800 ml
```

Simple addition. No direction needed. This is the behaviour of scalars — they add like regular numbers.

---

## 3. Vectors: Quantities with Size AND Direction

A **vector** is a physical quantity that requires both a magnitude AND a direction to be fully described.

This distinction sounds simple, but it changes everything about how we do calculations.

### 3.1 The Taxi Analogy

Imagine you are in an unfamiliar city and you call a taxi. You tell the driver: "Drive at 60 km/h." The driver nods and floors the accelerator.

But... in which direction? North? South? Into the river? Into the building in front of them?

Speed alone — 60 km/h — is not enough information to tell the taxi where to go. You need to say: "Drive at 60 km/h towards the airport, which is to the north-east." Now the driver has everything they need.

That complete instruction — both the size (60 km/h) and the direction (north-east) — is a **vector quantity**. In physics we call it **velocity**.

### 3.2 Common Vector Quantities

| Vector Quantity | Symbol | Common Unit | Needs Direction? |
|-----------------|--------|-------------|------------------|
| Velocity | v | m/s | Yes — e.g., "30 m/s north" |
| Force | F | Newtons (N) | Yes — e.g., "10 N downward" |
| Displacement | s or d | metres (m) | Yes — e.g., "5 m east" |
| Acceleration | a | m/s² | Yes — e.g., "9.8 m/s² downward" |
| Momentum | p | kg·m/s | Yes — same direction as velocity |
| Weight | W | Newtons (N) | Yes — always downward |
| Electric field | E | N/C | Yes — direction matters |
| Magnetic field | B | Tesla (T) | Yes — direction matters |

### 3.3 Key Vector Pairs to Remember

Physics has several pairs of quantities — one is a scalar, one is a vector — that students often confuse.

| Scalar Version | Vector Version | What Is Different |
|----------------|----------------|-------------------|
| Speed | Velocity | Velocity adds direction to speed |
| Distance | Displacement | Displacement adds direction to distance |
| Mass | Weight/Force | Weight is a force with direction (down) |
| Energy (scalar) | Momentum (vector) | Momentum has direction |

These pairs are critically important. Let us understand two of them in depth.

### 3.4 Distance vs. Displacement

**Distance** is how much total ground you cover. It is a scalar.
**Displacement** is how far you are from your starting point, in a straight line, and in which direction. It is a vector.

Here is a vivid example:

```
You start at your house (A).
You walk 3 km east to a shop (B).
Then you walk 3 km west back home (A).

        3 km east            3 km west
    A -----------> B -----------> A

Total distance walked = 3 km + 3 km = 6 km  (scalar)
Total displacement = 0 km  (you ended where you started!)
```

The distance is 6 km. The displacement is zero. Same trip — two completely different answers depending on whether you care about direction.

This is why athletes who run a full lap of a 400 m track travel a distance of 400 m but have a displacement of zero (they end up back where they started).

### 3.5 Speed vs. Velocity

**Speed** is how fast you are moving — just the magnitude. It is a scalar.
**Velocity** is how fast you are moving AND in what direction. It is a vector.

A car driving in a circle at a constant 60 km/h has:
- **Constant speed** (always 60 km/h)
- **Changing velocity** (the direction keeps changing as it goes around the circle)

This seems strange at first! How can something be changing if the number on the speedometer never changes? The answer: because velocity includes direction, and the direction is constantly changing as the car turns. This has real consequences — it means there is an acceleration even though the speed is constant (we will explore this in a later chapter).

---

## 4. The Core Difference: Why Direction Matters

The reason physicists insist on tracking direction is simple: **direction changes the result of calculations**.

### 4.1 A Tug-of-War Example

Imagine two people pulling a box.

**Scenario 1: Both pulling the same direction**

```
Person A          Box           Person B
(pulls right) --> [BOX] --> (pulls right)
    50 N                        30 N

Result: Total force = 50 + 30 = 80 N to the right
```

**Scenario 2: Pulling opposite directions**

```
Person A          Box           Person B
(pulls left) <-- [BOX] --> (pulls right)
    50 N                        30 N

Result: Total force = 50 - 30 = 20 N to the LEFT
```

Same magnitudes (50 N and 30 N), completely different results — just because the directions changed.

If you ignore direction and just add 50 + 30 = 80, you will get the wrong answer in Scenario 2. Physics would not work. Bridges would fall. Aeroplanes would crash. **Direction is not optional.**

### 4.2 The Navigation Example

A ship travels 4 km east, then 3 km north. Where does it end up?

If you just add 4 + 3 = 7, you might think it is 7 km from where it started. But that is wrong! Because the two legs of the journey were in different directions, we need to use the Pythagorean theorem:

```
                North
                  ^
                  |
                  |  3 km
                  |
                  +---------> East
             Start   4 km

Final distance from start = sqrt(4² + 3²)
                          = sqrt(16 + 9)
                          = sqrt(25)
                          = 5 km
```

The ship is only 5 km from its starting point, not 7 km. The direction of each journey segment matters enormously.

This is exactly what **vector addition** is — and we will work through it in full detail in Section 6.

---

## 5. Representing Vectors

Since vectors have both magnitude and direction, we need special tools to write them down and draw them.

### 5.1 Drawing Vectors as Arrows

The most natural way to represent a vector is with an **arrow**.

- The **length** of the arrow represents the magnitude (bigger magnitude = longer arrow)
- The **direction the arrow points** represents the direction of the vector

```
Small force (5 N):   -->

Larger force (15 N): ------------>

Force pointing up:   ^
                     |
                     |

Force pointing down: |
                     |
                     v

Force at an angle:   /
                    /
                   /
```

When drawing vector diagrams, we choose a scale. For example:
- 1 cm of arrow length = 5 N of force
- So a 10 N force would be drawn as a 2 cm arrow

### 5.2 Writing Vectors with Notation

In written physics, we need to distinguish vectors from scalars. Several notations are used:

**Bold type:** **F** for force (as a vector), F for the magnitude of force

**Arrow above:** F⃗ (an arrow drawn above the letter — common in handwriting)

**Component notation:** F = (Fx, Fy) — listing the x and y components

In this course, we will use bold (**v**, **F**, **a**) for vectors and plain letters (v, F, a) for magnitudes. When writing by hand, put an arrow above the letter.

### 5.3 Giving a Vector Its Full Description

A vector is not fully described unless you give both:
1. Its magnitude (with units)
2. Its direction (usually an angle or compass bearing)

**Examples:**
- "The wind blows at 20 m/s towards the north-east" ✓
- "The wind blows at 20 m/s" ✗ (incomplete — where is it blowing?)
- "The force is 50 N at 30° above the horizontal" ✓
- "The displacement is 8 km at a bearing of 045°" ✓

### 5.4 The Standard Reference System

In most physics problems, we use a coordinate system with:
- The **positive x-axis** pointing right (east)
- The **positive y-axis** pointing up (north)
- The **negative x-axis** pointing left (west)
- The **negative y-axis** pointing down (south)

```
              +y (North, Up)
               ^
               |
               |
-x  <----------+----------> +x
(West,Left)    |          (East, Right)
               |
               v
              -y (South, Down)
```

Angles are measured **counter-clockwise from the positive x-axis** as the standard convention, though navigation uses compass bearings (measured clockwise from north).

---

## 6. Adding Vectors: The Graphical Method

Adding two or more vectors graphically is called the **head-to-tail method** (or the triangle method). It is visual and intuitive — perfect for building understanding.

### 6.1 The Head-to-Tail Rule

To add two vectors **A** and **B**:

1. Draw vector **A** as an arrow.
2. Place the **tail** (start) of vector **B** at the **head** (tip) of vector **A**.
3. Draw a new arrow from the tail of **A** to the head of **B**. This new arrow is the **resultant** vector, **R** = **A** + **B**.

The **resultant** is the single vector that has the same effect as the original two vectors combined.

### 6.2 Example: Two Vectors in the Same Direction

A boy walks 3 m east, then 4 m east.

```
Vector A (3 m east):     --->      (3 m)
Vector B (4 m east):         ----> (4 m)

Head-to-tail:            --->---->
                         A    B

Resultant R:             -------->  (7 m east)
```

Result: 7 m east. Simple because directions are the same.

### 6.3 Example: Two Vectors at Right Angles

A person walks 4 m east, then 3 m north.

```
Draw A first (4 m east):

Start -----> (4 m east) End-of-A

Then place tail of B at head of A (3 m north):

         End-of-B
              ^
              |  3 m (north)
              |
Start ------> A-end

Now draw the resultant from Start to End-of-B:

         End
          ^   \
          |    \  <-- This diagonal arrow is R
        3 |     \
          |      \
Start ----+--------> 
              4

```

More clearly:

```
         +
         |\
    3 m  | \  R = ?
  (north)|  \
         |   \
         +----+
           4 m (east)

R = sqrt(4² + 3²) = sqrt(16 + 9) = sqrt(25) = 5 m
Direction = arctan(3/4) = 36.87° north of east (approximately 37°)
```

The resultant is 5 m at 37° north of east. This is the classic 3-4-5 right triangle.

### 6.4 Adding Three or More Vectors

The head-to-tail method extends to any number of vectors. You just keep chaining them:

```
Adding A, B, and C:

     --> (A)  ---> (B)  --> (C)
     
Chain them:
     
     --> ---> -->
     A    B    C

Resultant: from tail of A to head of C
```

**Example:** A bee flies 5 m east, then 2 m south, then 3 m west.

```
Step 1: Draw 5 m east
    --------> 

Step 2: At the head, draw 2 m south
    -------->
             |
             | 2 m

Step 3: At the new head, draw 3 m west
    -------->
             |
             | 2 m
             |
             <---  3 m

Step 4: Draw resultant from start to final position

Start        End
  +-------->  
  |        |
  |  Path  | 2 m
  |        |
  R        <--- 3 m
  |        
  Final

Net east: 5 - 3 = 2 m east
Net south: 2 m south

Resultant = sqrt(2² + 2²) = sqrt(8) = 2.83 m
Direction = 45° south of east (south-east)
```

### 6.5 The Parallelogram Method (Alternative)

Instead of head-to-tail, you can also use the **parallelogram method** when adding exactly two vectors from the same starting point:

1. Draw both vectors starting from the **same point**.
2. Complete the parallelogram (draw lines parallel to each vector).
3. The diagonal of the parallelogram is the resultant.

```
        B /|
         / |
        /  |
       /   |
      /    |
  A /______| 
    \       \
     \       \
      \       \  <- Resultant (diagonal)
       \       \
        \_______\

The diagonal from the common start point to the far corner
is the resultant R = A + B.
```

Both methods give the same answer — use whichever you find more intuitive.

---

## 7. Adding Vectors: The Algebraic Method

Graphical methods are useful for visualization, but they are slow and imprecise for calculations. The **algebraic method** breaks vectors into components (x and y parts) and adds those parts separately. This is the method used in real calculations.

### 7.1 What Are Components?

Any vector can be broken into a **horizontal part** (the x-component) and a **vertical part** (the y-component). These components are scalars — ordinary numbers that can be positive or negative.

```
Vector V at angle θ:

         V
        /|
       / |
      /  |  Vy (vertical component)
     /   |
    /θ   |
   +-----+
     Vx  
(horizontal component)

Vx = V × cos(θ)
Vy = V × sin(θ)
```

We will cover this in full detail in Section 10. For now, let us look at the algebraic addition method.

### 7.2 Steps for Algebraic Vector Addition

1. Break each vector into its x and y components.
2. Add all the x-components together to get the total Rx.
3. Add all the y-components together to get the total Ry.
4. Use the Pythagorean theorem to find the magnitude: R = sqrt(Rx² + Ry²)
5. Use trigonometry to find the direction: θ = arctan(Ry / Rx)

### 7.3 Worked Example: Adding Two Vectors Algebraically

**Problem:** Vector **A** has magnitude 6 m at 0° (due east). Vector **B** has magnitude 8 m at 90° (due north). Find the resultant.

**Step 1: Find components of A**
- Ax = 6 × cos(0°) = 6 × 1 = 6 m
- Ay = 6 × sin(0°) = 6 × 0 = 0 m

**Step 2: Find components of B**
- Bx = 8 × cos(90°) = 8 × 0 = 0 m
- By = 8 × sin(90°) = 8 × 1 = 8 m

**Step 3: Add components**
- Rx = Ax + Bx = 6 + 0 = 6 m
- Ry = Ay + By = 0 + 8 = 8 m

**Step 4: Find magnitude of resultant**
- R = sqrt(Rx² + Ry²) = sqrt(6² + 8²) = sqrt(36 + 64) = sqrt(100) = 10 m

**Step 5: Find direction**
- θ = arctan(Ry / Rx) = arctan(8/6) = arctan(1.333) ≈ 53.1°

**Answer: R = 10 m at 53.1° north of east**

```
         ^
       8 |   \
         |    \  R = 10 m
         |     \
         |  53.1°\
         +---------
              6
```

### 7.4 Worked Example: Two Vectors at Arbitrary Angles

**Problem:** Vector **P** = 10 N at 30° above horizontal. Vector **Q** = 8 N at 120° from positive x-axis. Find **P** + **Q**.

**Step 1: Components of P**
- Px = 10 × cos(30°) = 10 × 0.866 = 8.66 N
- Py = 10 × sin(30°) = 10 × 0.5 = 5.00 N

**Step 2: Components of Q** (120° means it points up and to the left)
- Qx = 8 × cos(120°) = 8 × (-0.5) = -4.00 N
- Qy = 8 × sin(120°) = 8 × 0.866 = 6.93 N

**Step 3: Add components**
- Rx = 8.66 + (-4.00) = 4.66 N
- Ry = 5.00 + 6.93 = 11.93 N

**Step 4: Magnitude**
- R = sqrt(4.66² + 11.93²) = sqrt(21.72 + 142.32) = sqrt(164.04) ≈ 12.81 N

**Step 5: Direction**
- θ = arctan(11.93 / 4.66) = arctan(2.56) ≈ 68.7° above horizontal

**Answer: R ≈ 12.81 N at 68.7° above the horizontal**

---

## 8. Subtracting Vectors

Vector subtraction follows naturally from addition. To subtract vector **B** from vector **A**:

**A** - **B** = **A** + (-**B**)

where **-B** is the vector **B** with its direction reversed (same length, opposite direction).

### 8.1 The Negative of a Vector

The **negative of a vector** has the same magnitude but points in exactly the opposite direction.

```
Vector B:        ------>   (pointing right, magnitude 5 m)

Negative of B:   <------   (pointing left, magnitude 5 m, i.e., -B)
```

If **B** = 5 m east, then **-B** = 5 m west.

### 8.2 Graphical Subtraction

To find **A** - **B**:
1. Reverse vector **B** to get **-B**
2. Add **A** and **-B** using the head-to-tail method

```
Example: A - B where A = 6 m east, B = 4 m east

-B = 4 m WEST

Head-to-tail of A and -B:

A (6 m east): ------>
-B (4 m west): <----

Combined (tail of A to head of -B):
      ------><----
         A    -B

Net result: -->    (2 m east)
```

A - B = 2 m east, which makes sense: 6 - 4 = 2.

### 8.3 Worked Example: Subtracting Vectors at a Right Angle

**Problem:** **A** = 5 m north, **B** = 3 m east. Find **A** - **B**.

This means find **A** + (**-B**).

**-B** = 3 m west

```
Vector A (5 m north):       ^
                             |
                             | 5
                             |
                             |
                             
Vector -B (3 m west):   <---
                          3

Head-to-tail (A first, then -B):

Start at bottom. Draw A upward:
           ^ A-end
           |
           | 5 m
           |
        A-start

Then from A-end, draw -B (3 m west):
    -B-end <--- A-end
    
Now resultant is from A-start to -B-end:

           ^ A-end ----> ...
           |
           | (this is not right, let me redo)

Let me draw more carefully:

        -B-end
           +<------+ A-end
                   |
                   | 5 m (A)
                   |
               A-start (= our origin)

Resultant goes from A-start to -B-end (diagonal):

        +<--------+
        |          |
        |    ?     | 5 m
        |          |
        +-----------
             3 m

Rx = 0 + (-3) = -3 m (3 m west)
Ry = 5 + 0 = 5 m (north)

|R| = sqrt((-3)² + 5²) = sqrt(9 + 25) = sqrt(34) ≈ 5.83 m
Direction = arctan(5/3) ≈ 59° west of north (or 121° from positive x-axis)
```

### 8.4 Algebraic Subtraction

Algebraically, subtract the components:

If **A** = (Ax, Ay) and **B** = (Bx, By), then:

**A** - **B** = (Ax - Bx, Ay - By)

**Example:**
- **A** = (8, 6)
- **B** = (3, 2)
- **A** - **B** = (8-3, 6-2) = (5, 4)
- Magnitude = sqrt(5² + 4²) = sqrt(25+16) = sqrt(41) ≈ 6.40

---

## 9. Multiplying a Vector by a Scalar

Sometimes we need to scale a vector up or down. This is done by **multiplying a vector by a scalar**.

### 9.1 What Happens When You Multiply

When you multiply a vector **v** by a scalar number k:

- If **k > 0**: The new vector points in the **same direction**, but its magnitude is k times larger.
- If **k < 0**: The new vector points in the **opposite direction**, and its magnitude is |k| times larger.
- If **k = 0**: The result is the zero vector (no magnitude, no direction).

```
Vector v (5 m east):        ----->

2 × v (10 m east):          --------->

0.5 × v (2.5 m east):       -->

-1 × v (5 m west):          <-----

-2 × v (10 m west):         <---------
```

### 9.2 Formal Definition

If **v** has components (vx, vy), then:

k × **v** = (k × vx, k × vy)

The magnitude of k × **v** = |k| × |**v**|

The direction is the same as **v** if k is positive, opposite if k is negative.

### 9.3 Real-World Example: Newton's Second Law

One of the most important equations in physics — Newton's second law — is a scalar multiplication:

**F** = m × **a**

Here, mass m is a scalar (it has no direction). Acceleration **a** is a vector. So force **F** is a vector pointing in the **same direction** as the acceleration, but scaled by the mass.

If a 3 kg ball accelerates at 4 m/s² northward:
- Force = 3 kg × 4 m/s² northward = **12 N northward**

The direction of the force is the same as the direction of the acceleration. The mass just scales the magnitude.

### 9.4 Worked Example: Scaling a Vector

**Problem:** A force vector **F** = (4, 3) N. Find 3**F** and -0.5**F**.

**3F:**
- Components: (3×4, 3×3) = (12, 9)
- Magnitude: sqrt(12² + 9²) = sqrt(144 + 81) = sqrt(225) = 15 N
- Direction: arctan(9/12) = arctan(0.75) = 36.87° above horizontal

**-0.5F:**
- Components: (-0.5×4, -0.5×3) = (-2, -1.5)
- Magnitude: sqrt((-2)² + (-1.5)²) = sqrt(4 + 2.25) = sqrt(6.25) = 2.5 N
- Direction: 180° + 36.87° = 216.87° (pointing backwards compared to original)

---

## 10. Resolving Vectors into Components

This section is perhaps the most practically important in the entire chapter. **Resolving** a vector means breaking it into its x (horizontal) and y (vertical) parts. This is the key skill that makes every vector calculation possible.

### 10.1 Why We Resolve Vectors

When a vector points at an angle, it is doing two things at once — it has both a horizontal effect and a vertical effect.

Think of a person pulling a suitcase with a rope:

```
                   /  <- rope at angle θ
                  /
                 /
                / θ
---------------+---- floor

The rope pull has:
- A horizontal part: pulling the suitcase forward
- A vertical part: trying to lift the suitcase up
```

To understand the motion of the suitcase, we need to separate these two effects. That is resolving the vector.

### 10.2 The Trigonometry Behind Components

Consider a vector **V** with magnitude V pointing at angle θ above the horizontal (positive x-axis).

```
         |  V
         | /
         |/  θ (angle from horizontal)
    Vy   +----------
      ^     Vx
      |
      |  (The vertical component Vy)
      |
   We can think of V as the hypotenuse of a right triangle:
   
         V (hypotenuse)
        /|
       / |
      /  | Vy (opposite side)
     /   |
    / θ  |
   +-----+
     Vx  (adjacent side)
```

From basic trigonometry (SOH-CAH-TOA):

- **Vx = V × cos(θ)** — the horizontal component
- **Vy = V × sin(θ)** — the vertical component

Where:
- V is the magnitude of the vector
- θ is the angle measured from the positive x-axis (horizontal)

### 10.3 SOH-CAH-TOA Reminder

If you have forgotten your trigonometry:

```
     Hypotenuse (H)
        /|
       / |
      /  | Opposite (O)
     /   |
    / θ  |
   +-----+
   Adjacent (A)

sin(θ) = Opposite / Hypotenuse = O/H  ← "SOH"
cos(θ) = Adjacent / Hypotenuse = A/H  ← "CAH"
tan(θ) = Opposite / Adjacent   = O/A  ← "TOA"
```

The horizontal side is always "Adjacent" to the angle θ, so we use **cosine for x**.
The vertical side is always "Opposite" to the angle θ, so we use **sine for y**.

Memory trick: **"x is for cosine"** (both start with... okay they do not, but cos gives the x-axis component — just memorise this one fact and the rest follows.)

### 10.4 Signs of Components

The angle θ can place our vector in any of the four quadrants. The signs of the components tell us direction:

```
        +y
        |
  II    |    I
  (x-,y+) | (x+,y+)
        |
--------+--------  +x
        |
  III   |    IV
  (x-,y-) | (x+,y-)
        |
        -y
```

| Quadrant | x-component | y-component | Example direction |
|----------|-------------|-------------|-------------------|
| I (0° to 90°) | Positive | Positive | North-East |
| II (90° to 180°) | Negative | Positive | North-West |
| III (180° to 270°) | Negative | Negative | South-West |
| IV (270° to 360°) | Positive | Negative | South-East |

### 10.5 Worked Example: Resolving a Force

**Problem:** A force of 20 N acts at 60° above the horizontal (to the right). Find its x and y components.

```
              /
             / 20 N
            /
           / 60°
          +------->
         
```

**Solution:**
- Fx = F × cos(60°) = 20 × 0.5 = **10 N** (horizontal, to the right)
- Fy = F × sin(60°) = 20 × 0.866 = **17.32 N** (vertical, upward)

**Check:** sqrt(Fx² + Fy²) = sqrt(10² + 17.32²) = sqrt(100 + 300) = sqrt(400) = 20 N ✓

### 10.6 Worked Example: Inclined Plane

**Problem:** A 50 N weight sits on a ramp inclined at 30° to the horizontal. Resolve the weight (acting straight down) into components parallel and perpendicular to the ramp.

```
         /|
        / |
       /  |
      /   |
     / W  |
    / |   |
   /  v   |
  /30°    |
 /________|

W = 50 N downward

Resolved into:
- W_parallel (down the ramp) = W × sin(30°) = 50 × 0.5 = 25 N
- W_perpendicular (into the ramp) = W × cos(30°) = 50 × 0.866 = 43.3 N
```

Note: For inclined plane problems, we sometimes resolve along the incline rather than along horizontal/vertical axes. The principle is the same — pick the axes that make your problem easiest.

### 10.7 Worked Example: Projectile at an Angle

**Problem:** A ball is kicked at 25 m/s at an angle of 40° to the ground. Find the horizontal and vertical components of its initial velocity.

```
        *
       / \
      /   \
     /     \
    / 40°   \
   +----------
   
v = 25 m/s at 40°
```

**Solution:**
- vx = 25 × cos(40°) = 25 × 0.766 = **19.15 m/s** (horizontal)
- vy = 25 × sin(40°) = 25 × 0.643 = **16.07 m/s** (vertical, upward)

The horizontal component tells us how fast the ball moves across the ground. The vertical component tells us how fast it initially moves upward. These two components behave independently (we will explore this fully in projectile motion chapters).

---

## 11. Finding the Resultant: Magnitude and Direction

After adding or resolving vectors, we often need to turn the components back into a single vector with a magnitude and direction. This is called finding the **resultant**.

### 11.1 Magnitude from Components

Given the x and y components of a vector (Rx and Ry), the magnitude is:

**R = sqrt(Rx² + Ry²)**

This is just the Pythagorean theorem — the components form the two legs of a right triangle, and the vector is the hypotenuse.

### 11.2 Direction from Components

The direction angle θ (measured from the positive x-axis) is:

**θ = arctan(Ry / Rx)**

You also need to check which quadrant you are in (see Section 10.4) because arctan alone gives an angle in the range -90° to +90°, and you need to adjust for angles in the second and third quadrants.

### 11.3 Quadrant Adjustment

| Rx | Ry | Quadrant | Adjustment |
|----|----|----------|------------|
| + | + | I | θ = arctan(Ry/Rx) — no change |
| - | + | II | θ = arctan(Ry/Rx) + 180° |
| - | - | III | θ = arctan(Ry/Rx) + 180° |
| + | - | IV | θ = arctan(Ry/Rx) + 360° (or just write it as negative) |

### 11.4 Comprehensive Worked Example

**Problem:** Three forces act on an object. Find the resultant force.
- **F1** = 10 N at 0° (due east)
- **F2** = 15 N at 120° (up and to the left)
- **F3** = 8 N at 270° (due south)

**Step 1: Resolve each vector into components**

F1 (10 N at 0°):
- F1x = 10 × cos(0°) = 10 × 1 = **10 N**
- F1y = 10 × sin(0°) = 10 × 0 = **0 N**

F2 (15 N at 120°):
- F2x = 15 × cos(120°) = 15 × (-0.5) = **-7.5 N**
- F2y = 15 × sin(120°) = 15 × 0.866 = **12.99 N**

F3 (8 N at 270°):
- F3x = 8 × cos(270°) = 8 × 0 = **0 N**
- F3y = 8 × sin(270°) = 8 × (-1) = **-8 N**

**Step 2: Sum the components**

Rx = F1x + F2x + F3x = 10 + (-7.5) + 0 = **2.5 N**
Ry = F1y + F2y + F3y = 0 + 12.99 + (-8) = **4.99 N**

**Step 3: Find magnitude**

R = sqrt(Rx² + Ry²) = sqrt(2.5² + 4.99²) = sqrt(6.25 + 24.9) = sqrt(31.15) ≈ **5.58 N**

**Step 4: Find direction**

θ = arctan(Ry / Rx) = arctan(4.99 / 2.5) = arctan(1.996) ≈ **63.4°**

Both Rx and Ry are positive, so we are in Quadrant I — no adjustment needed.

**Answer: The resultant force is approximately 5.58 N at 63.4° above the horizontal (north of east).**

```
       +y
        ^
        |   /
 4.99 N |  /  R ≈ 5.58 N
        | / 63.4°
        |/________
             2.5 N      +x
```

---

## 12. Unit Vectors i, j, and k

As you go further in physics, you will encounter a very convenient notation for vectors using **unit vectors**.

### 12.1 What Is a Unit Vector?

A **unit vector** is a vector with a magnitude of exactly **1** (no units — just direction). Unit vectors are used as "direction pointers" — they tell us which way something points without saying anything about the size.

The three standard unit vectors in three-dimensional space are:

| Symbol | Direction | What it points along |
|--------|-----------|----------------------|
| **i** (or î) | Positive x-axis | East / Right / Horizontal |
| **j** (or ĵ) | Positive y-axis | North / Up / Vertical |
| **k** (or k̂) | Positive z-axis | Out of the page / Depth |

```
             +y
              ^
              |  ĵ (unit vector in y direction)
              |
-x  <---------+---------> +x
              |         î (unit vector in x direction)
              |
              v
              -y
              
              (k̂ would be pointing out of the page toward you)
```

### 12.2 Writing Vectors Using Unit Vector Notation

Any vector in 2D can be written as:

**V** = Vx **i** + Vy **j**

This means: "Take Vx steps in the x-direction and Vy steps in the y-direction."

**Examples:**
- A displacement of 5 m east: **d** = 5 **i** + 0 **j** = 5**i**
- A velocity of 3 m/s north: **v** = 0 **i** + 3 **j** = 3**j**
- A force of 4 N east and 3 N north: **F** = 4**i** + 3**j**

In 3D, you add the z-component:
**V** = Vx **i** + Vy **j** + Vz **k**

### 12.3 Adding Vectors in Unit Vector Notation

Adding is very clean in this notation:

If **A** = 3**i** + 4**j** and **B** = 2**i** - 1**j**

Then **A** + **B** = (3+2)**i** + (4-1)**j** = 5**i** + 3**j**

You just add the coefficients of the same unit vectors separately. The i terms together, the j terms together, the k terms together.

### 12.4 Worked Example

**Problem:** Force **F1** = 6**i** + 8**j** N and force **F2** = -3**i** + 2**j** N. Find the resultant and its magnitude.

**Solution:**
- **R** = **F1** + **F2** = (6 + (-3))**i** + (8 + 2)**j** = **3i + 10j** N
- Rx = 3 N, Ry = 10 N
- |**R**| = sqrt(3² + 10²) = sqrt(9 + 100) = sqrt(109) ≈ **10.44 N**
- Direction = arctan(10/3) = arctan(3.33) ≈ **73.3°** above the positive x-axis

### 12.5 Unit Vectors in 3D (Brief Overview)

In three dimensions, the **k** component handles depth. For example:

A force pushing diagonally in 3D: **F** = 5**i** + 3**j** + 2**k** N

Magnitude: |**F**| = sqrt(5² + 3² + 2²) = sqrt(25 + 9 + 4) = sqrt(38) ≈ 6.16 N

The Pythagorean theorem simply extends to three dimensions — we will explore this more in later chapters.

---

## 13. Practical Applications

Vectors are not abstract mathematics — they are the language of the physical world. Here are some of the most important places where vectors appear in everyday life.

### 13.1 Navigation and GPS

Every time your phone's GPS navigates you somewhere, it is doing vector calculations constantly.

**Dead Reckoning Example:**

A ship starts at position A. It travels:
- 30 km at 040° (slightly east of north)
- Then 20 km at 130° (slightly south of east)

Where does it end up?

```
Step 1: Resolve first leg (30 km at 40° from north)
    Note: In navigation, angles are measured from NORTH (not x-axis)
    So the angle from east (our x-axis) = 90° - 40° = 50°
    
    x1 = 30 × cos(50°) = 30 × 0.643 = 19.28 km east
    y1 = 30 × sin(50°) = 30 × 0.766 = 22.98 km north

Step 2: Resolve second leg (20 km at 130° from north = 40° below east)
    Angle from x-axis = 90° - 130° = -40° (below horizontal)
    
    x2 = 20 × cos(-40°) = 20 × 0.766 = 15.32 km east
    y2 = 20 × sin(-40°) = 20 × (-0.643) = -12.86 km (southward)

Step 3: Total displacement
    X = 19.28 + 15.32 = 34.60 km east
    Y = 22.98 + (-12.86) = 10.12 km north

Step 4: Magnitude
    R = sqrt(34.60² + 10.12²) = sqrt(1197 + 102) = sqrt(1299) ≈ 36.04 km

Step 5: Direction (bearing from north)
    θ from east = arctan(10.12/34.60) = arctan(0.292) = 16.3°
    Bearing from north = 90° - 16.3° = 73.7°
    
    The ship ends up 36.04 km from start, at a bearing of about 074°.
```

### 13.2 Aeroplane and Wind Velocity

When an aeroplane flies, it must account for wind. The aeroplane's velocity relative to the ground is the **vector sum** of its airspeed and the wind velocity.

**Example:**

An aeroplane flies at 300 km/h due east (its speed through the air). There is a wind blowing at 50 km/h from south to north. What is the plane's actual velocity over the ground?

```
Plane's velocity (through air): 300 km/h east
Wind velocity: 50 km/h north

Vector addition:
    +y (north)
    ^
    | 50 km/h (wind)
    |
    +-----------> +x (east)
         300 km/h (plane)

Resultant velocity over ground:
Rx = 300 km/h
Ry = 50 km/h

Speed = sqrt(300² + 50²) = sqrt(90000 + 2500) = sqrt(92500) ≈ 304.1 km/h
Direction = arctan(50/300) = arctan(0.167) ≈ 9.5° north of east
```

The wind pushes the plane slightly north, so it arrives at a point slightly north of where it was aiming.

**The pilot must correct for this!** To fly straight east, the pilot must angle the nose slightly south to counteract the northward wind. Calculating this correction requires vector subtraction.

### 13.3 Force Diagrams (Free Body Diagrams)

In engineering and physics, we draw **free body diagrams** to show all the forces acting on an object. Each force is a vector — magnitude and direction.

**Example:** A traffic light hanging from two wires.

```
       Wall A                          Wall B
         \      wire 1    wire 2      /
          \___/        \___/
          T1             T2
          ^               ^
           \             /
            \  θ1     θ2 /
             \         /
              \       /
               [Light]
                  |
                  | W (weight downward)
                  v
```

For the light to be in equilibrium (not moving), all forces must balance:
- Sum of horizontal components = 0
- Sum of vertical components = 0

This gives us a system of equations to solve for the tensions T1 and T2.

**Simplified version:** A lamp of weight 100 N hangs from two equal wires, each at 45° to the vertical.

```
        T1        T2
         ^         ^
          \       /
           \ 45° / 45°
            \   /
             \ /
            [LAMP]
              |
           W = 100 N
              |
              v

Vertical equilibrium:
T1 × cos(45°) + T2 × cos(45°) = W
2T × 0.707 = 100
T = 100 / (2 × 0.707) = 70.7 N

Each wire carries 70.7 N.
```

### 13.4 Sport: Kicking a Ball

When a footballer kicks a ball at an angle, the ball's initial velocity is a vector. Breaking it into components allows us to predict:
- How far the ball travels horizontally (the horizontal component of velocity)
- How high it goes (the vertical component of velocity)

**Example:** Ball kicked at 20 m/s at 35° to the ground.

```
           *
         *   *
        *     *
       *       *
      *         *
     *           *
----*-----ground-----*----

Horizontal velocity = 20 × cos(35°) = 20 × 0.819 = 16.38 m/s
Vertical velocity = 20 × sin(35°) = 20 × 0.574 = 11.47 m/s

The ball keeps its horizontal velocity (16.38 m/s) throughout.
The vertical velocity decreases due to gravity (9.8 m/s²).
```

### 13.5 Ramps and Inclined Planes

When an object sits on a ramp, gravity pulls it straight down. But we resolve this into:
- A component parallel to the ramp (trying to slide the object down)
- A component perpendicular to the ramp (pressing the object into the surface)

This resolution is essential for calculating whether an object will slide, and how much friction is needed to stop it.

```
        /|
       / |
      /  |
     /   |
    /θ   |
   ------+

Weight W acts straight down.

W_parallel (down the slope) = W × sin(θ)
W_perpendicular (into slope) = W × cos(θ)

For θ = 30° and W = 200 N:
W_parallel = 200 × sin(30°) = 200 × 0.5 = 100 N
W_perpendicular = 200 × cos(30°) = 200 × 0.866 = 173.2 N
```

### 13.6 Forces in Equilibrium

An object is in **equilibrium** (not accelerating, not moving or moving at constant velocity) when all the vectors acting on it add up to zero. This principle is the foundation of structural engineering.

```
Example: A book resting on a table.

Forces:
- Weight W = 5 N downward (gravity pulling it down)
- Normal force N = 5 N upward (table pushing it up)

      N = 5 N upward
         ^
         |
       [book]
         |
         v
      W = 5 N downward
      
Net force = N + W = 5 N up + 5 N down = 0 N
(The vectors cancel. The book does not accelerate.)
```

For anything more complex — a bridge, a skyscraper, an aircraft wing — engineers use exactly this principle, just with many more vectors.

---

## 14. Summary

This chapter introduced one of the most fundamental distinctions in all of physics: the difference between scalar and vector quantities. Here are the key takeaways:

**What We Learned:**

- A **scalar** is a physical quantity described by magnitude (size) only. Examples: temperature, mass, speed, time, distance, energy.

- A **vector** is a physical quantity described by both magnitude AND direction. Examples: velocity, force, displacement, acceleration, momentum.

- The key pairs to remember: speed (scalar) vs. velocity (vector); distance (scalar) vs. displacement (vector).

- Vectors are represented as **arrows**, where length indicates magnitude and the arrow direction indicates the vector's direction.

- **Vector addition** using the head-to-tail method: draw vectors tip-to-tail; the resultant goes from the start of the first to the end of the last.

- **Algebraic addition**: break each vector into components (x and y), add the components separately, then reconstruct the resultant using the Pythagorean theorem and trigonometry.

- **Vector subtraction**: A - B = A + (-B), where -B has the same magnitude as B but the opposite direction.

- **Scalar multiplication**: multiplying a vector by a positive scalar stretches/shrinks it; multiplying by a negative scalar reverses its direction.

- **Resolving a vector** into components: Vx = V cos(θ) and Vy = V sin(θ), where θ is the angle from the positive x-axis.

- The **magnitude of a resultant**: R = sqrt(Rx² + Ry²)

- The **direction of a resultant**: θ = arctan(Ry / Rx), adjusted for the correct quadrant.

- **Unit vectors** i, j, k point along the positive x, y, and z axes respectively, each with magnitude 1. Any vector can be written as V = Vx i + Vy j + Vz k.

- Vector mathematics is the foundation for navigation, structural engineering, aeronautics, sports science, and virtually every other branch of applied physics.

**Key Insights:**

- You cannot add or subtract quantities that have different units or types — you can only add forces to forces, velocities to velocities, and so on.
- Direction is not optional in physics. The same magnitude in a different direction gives a completely different physical result.
- The component method (resolving into x and y) is the most powerful technique in 2D vector problems. Master it and most problems become straightforward.
- When all forces on an object sum to zero (as vectors), the object is in equilibrium.

---

## 15. Key Equations

Below are the important equations from this chapter. Memorise these — they will appear repeatedly throughout your physics journey.

**Resolving a vector into components (angle θ from positive x-axis):**

```
Vx = V × cos(θ)
Vy = V × sin(θ)
```

**Recovering magnitude and direction from components:**

```
V = sqrt(Vx² + Vy²)      [Pythagorean theorem]
θ = arctan(Vy / Vx)       [angle from positive x-axis, adjust for quadrant]
```

**Vector addition (component method):**

```
Rx = Ax + Bx + Cx + ...   [add all x-components]
Ry = Ay + By + Cy + ...   [add all y-components]
R  = sqrt(Rx² + Ry²)      [magnitude of resultant]
θ  = arctan(Ry / Rx)       [direction of resultant]
```

**Vector subtraction:**

```
A - B = A + (-B)
(Ax - Bx, Ay - By) = components of A - B
```

**Scalar multiplication:**

```
k × V = (k × Vx, k × Vy)
|k × V| = |k| × |V|
```

**Unit vector notation (2D):**

```
V = Vx i + Vy j
```

**Unit vector notation (3D):**

```
V = Vx i + Vy j + Vz k
|V| = sqrt(Vx² + Vy² + Vz²)
```

**Equilibrium condition (object not accelerating):**

```
Sum of all force vectors = 0
ΣFx = 0   AND   ΣFy = 0
```

**Trigonometry reference (SOH-CAH-TOA):**

```
sin(θ) = Opposite / Hypotenuse
cos(θ) = Adjacent / Hypotenuse
tan(θ) = Opposite / Adjacent
```

---

*In Chapter 04, we will use these vector skills to study motion in one and two dimensions — building up to understanding how objects move through the air, how cars accelerate on roads, and how physics describes the motion of every object in the universe.*
