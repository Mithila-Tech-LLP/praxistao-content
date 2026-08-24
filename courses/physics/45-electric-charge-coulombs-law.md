# Chapter 45: Electric Charge and Coulomb's Law

> **"The electric force is 10³⁹ times stronger than gravity. If you could extract the electrons from just two grams of copper and place them a meter apart, the repulsion would be strong enough to lift the entire Earth."**

---

## Table of Contents

- [45.1 What is Electric Charge?](#451-what-is-electric-charge)
- [45.2 Charge Quantization](#452-charge-quantization)
- [45.3 Conductors and Insulators](#453-conductors-and-insulators)
- [45.4 Charging Methods](#454-charging-methods)
- [45.5 Coulomb's Law](#455-coulombs-law)
- [45.6 Comparing Electric and Gravitational Forces](#456-comparing-electric-and-gravitational-forces)
- [45.7 Superposition Principle](#457-superposition-principle)
- [45.8 Electric Force Problems](#458-electric-force-problems)
- [Summary](#summary)
- [Key Equations](#key-equations)

---

## 45.1 What is Electric Charge?

Electric charge is a fundamental property of matter, just like mass.

```
THE TWO TYPES OF CHARGE:

POSITIVE (+)         NEGATIVE (-)
  proton               electron
  
Rule: LIKE charges REPEL, UNLIKE charges ATTRACT

  +  +                +  -
  |  |                |  |
  v  v                ^ v (attracted)
  (push apart)        
```

Every normal atom is electrically neutral: equal numbers of protons (+) and electrons (-).

When charges are separated, we say the object is **charged** or **electrified**.

### The Source of Charge in Matter

```
ATOMIC STRUCTURE:

  Nucleus: protons (+) and neutrons (0) — tightly bound, rarely moved
  
  Electron cloud: electrons (-) — loosely bound, can move
  
  When objects get "charged", it's always ELECTRONS being moved.
  (Protons stay fixed in their nuclei — too heavy, too tightly bound)
```

---

## 45.2 Charge Quantization

Electric charge comes in discrete chunks — it is **quantized**.

The smallest possible charge is the **elementary charge**:

```
e = 1.602 × 10⁻¹⁹ C  (C = coulomb, the SI unit of charge)

Electron: charge = -e = -1.602 × 10⁻¹⁹ C
Proton:   charge = +e = +1.602 × 10⁻¹⁹ C
Neutron:  charge = 0
```

Any measured charge is a whole-number multiple of e:

```
q = n × e    where n = 0, ±1, ±2, ±3, ...

You can never have a charge of 0.5e or 1.3e — only integer multiples.
```

### How Big is a Coulomb?

One coulomb is a HUGE amount of charge.

- 1 C = 6.24 × 10¹⁸ elementary charges
- The charge on a typical balloon rubbed on hair: ~1 μC = 10⁻⁶ C
- Current of 1 ampere: 1 C flowing past per second

---

## 45.3 Conductors and Insulators

**Conductors**: materials in which charges (electrons) can move freely.
- Examples: metals (copper, silver, aluminum), salt water, carbon (graphite)
- Metals have "free electrons" not bound to any particular atom

**Insulators**: materials in which charges cannot move.
- Examples: rubber, plastic, glass, wood, dry air, ceramics
- All electrons are tightly bound to their atoms

```
CONDUCTOR:        INSULATOR:
  [metal]           [rubber]
  
  + + + + +         + + + + +
  . . . . .         - - - - -  (stuck in place)
  - - - - -  (electrons can drift)
  
When charge is placed on a conductor, it spreads out.
When placed on an insulator, it stays put.
```

**Semiconductors** (silicon, germanium): in between — can be made to conduct or insulate by adding impurities or applying voltage. The basis of all electronics.

---

## 45.4 Charging Methods

### Charging by Friction (Triboelectric Effect)

When two materials are rubbed together, electrons transfer from one to the other.

```
BALLOON + HAIR:
  
  Before rubbing:  both neutral
  
  After rubbing:
    Hair: lost electrons → positively charged
    Balloon: gained electrons → negatively charged
    
    Hair and balloon attract each other!
```

Different materials have different "electron affinities." The triboelectric series lists them:

Glass → Wool → Fur → Silk → Rubber → Plastic

Glass rubbed with silk: glass becomes +, silk becomes -

### Charging by Conduction

Touch a charged object to an uncharged one — charge flows through the contact:

```
CHARGED ROD (+++) touches NEUTRAL BALL:
  
  Touch!     Charge flows:    Separate:
  
  +++  (o)    +  (+)           + (++)
  
  Ball becomes positively charged (same sign as rod).
```

### Charging by Induction

Charge a conductor WITHOUT touching it — by bringing a charged object nearby:

```
Step 1: Bring + charged rod near a neutral conductor.
        Electrons in conductor attracted toward rod.
  
  +++            [- - - + + +]
  rod            conductor
  
Step 2: Ground the far side (electrons flow into/out of ground).
  
  +++            [- - -      ]
                       |
                      ground
                      
Step 3: Remove ground. Conductor now has net negative charge.
Step 4: Remove rod. Charge redistributes. Conductor remains negative!
```

Key: the conductor ends up with charge OPPOSITE to the inducing charge.

---

## 45.5 Coulomb's Law

In 1785, Charles-Augustin de Coulomb measured the electric force between charged objects.

**Coulomb's Law:** The magnitude of the electric force between two point charges is proportional to the product of the charges and inversely proportional to the square of the distance between them.

```
         k × |q₁| × |q₂|
F = -----------------------
              r²

where:
  F  = magnitude of force (N)
  q₁, q₂ = charges (C)
  r  = distance between charges (m)
  k  = Coulomb's constant = 8.99 × 10⁹ N·m²/C²
```

Often written as k = 1/(4πε₀), where ε₀ = 8.85 × 10⁻¹² C²/(N·m²) is the permittivity of free space.

### Direction of the Force

```
SAME SIGN (repulsion):          OPPOSITE SIGN (attraction):

  +q₁ ← F₁₂    F₂₁ → +q₂        +q₁ → F₁₂    F₂₁ ← -q₂
  
  Forces point away from each other    Forces point toward each other
```

The forces are equal in magnitude and opposite in direction (Newton's Third Law).

### The Inverse-Square Law

Like gravity, Coulomb's force follows an inverse-square law:

```
Double the distance → force becomes 1/4
Triple the distance → force becomes 1/9
Half the distance → force becomes 4×

F ∝ 1/r²
```

### Worked Example 45.1

Two charges: q₁ = +3 μC, q₂ = -5 μC, separated by r = 0.2 m.

Find the magnitude and direction of the force.

**Solution:**

F = k|q₁||q₂|/r² = 8.99×10⁹ × 3×10⁻⁶ × 5×10⁻⁶ / (0.2)²

F = 8.99×10⁹ × 15×10⁻¹² / 0.04

F = 134.85×10⁻³ / 0.04 = **3.37 N**

Since q₁ is positive and q₂ is negative, the force is **attractive** — they pull toward each other.

---

## 45.6 Comparing Electric and Gravitational Forces

Both Coulomb's law and Newton's law of gravity have the same mathematical form (inverse-square laws). But they differ dramatically in strength.

```
COULOMB'S LAW:           NEWTON'S GRAVITY:

      k q₁q₂                    G m₁m₂
F = --------             F = ----------
       r²                          r²

k = 8.99 × 10⁹           G = 6.67 × 10⁻¹¹

RATIO of electric to gravitational force for two protons:

  k e²                  8.99×10⁹ × (1.6×10⁻¹⁹)²
 -------  =           ---------------------------------   ≈ 10³⁶
  G mp²               6.67×10⁻¹¹ × (1.67×10⁻²⁷)²

The electric force is about 10³⁶ times stronger than gravity!
```

Why doesn't electricity dominate our everyday experience? Because matter is almost always electrically neutral — positive and negative charges cancel nearly perfectly. Gravity dominates at large scales because mass only comes in one sign, so it always adds up.

---

## 45.7 Superposition Principle

When multiple charges are present, the total force on any one charge is the **vector sum** of all the individual Coulomb forces from each other charge.

```
SUPERPOSITION PRINCIPLE:
  F_total = F₁ + F₂ + F₃ + ...
  (vector addition)
```

### Worked Example 45.2

Three charges in a line:
- q₁ = +4 μC at x = 0
- q₂ = -2 μC at x = 0.3 m
- Find the force on q₂ due to q₁.

(Already done above — one charge, one force.)

Now add a third:
- q₃ = +1 μC at x = 0.6 m

Find the net force on q₂ at x = 0.3 m.

**Solution:**

Force on q₂ from q₁ (distance = 0.3 m):
F₁₂ = k|q₁||q₂|/r² = 8.99×10⁹ × 4×10⁻⁶ × 2×10⁻⁶ / 0.09 = **0.799 N**
Direction: q₁(+) attracts q₂(-) → force points left (toward q₁)

Force on q₂ from q₃ (distance = 0.3 m):
F₃₂ = k|q₃||q₂|/r² = 8.99×10⁹ × 1×10⁻⁶ × 2×10⁻⁶ / 0.09 = **0.200 N**
Direction: q₃(+) attracts q₂(-) → force points right (toward q₃)

Net force on q₂:
Taking right as positive:
F_net = +0.200 - 0.799 = **-0.599 N** (0.599 N pointing left)

---

## 45.8 Electric Force Problems

### Worked Example 45.3 — Charges in 2D

```
        q₂ = +3 μC
          |
          | 0.4 m
          |
q₁ = -2 μC ---- 0.3 m ---- q₃ = ?
```

Find the force on q₁ due to q₂ and q₃ if q₃ = +5 μC.

**Forces on q₁:**

From q₂ (0.4 m away, opposite signs → attractive):
F₁₂ = k × 2×10⁻⁶ × 3×10⁻⁶ / (0.4)² = 8.99×10⁹ × 6×10⁻¹² / 0.16 = **0.337 N upward**

From q₃ (0.3 m away, opposite signs → attractive):
F₁₃ = k × 2×10⁻⁶ × 5×10⁻⁶ / (0.3)² = 8.99×10⁹ × 10×10⁻¹² / 0.09 = **0.999 N rightward**

Net force magnitude:
F = sqrt(0.337² + 0.999²) = sqrt(0.114 + 0.998) = sqrt(1.112) ≈ **1.05 N**

Direction: angle θ = arctan(0.337/0.999) ≈ 18.6° above horizontal (toward right and up)

---

## Summary

- **Electric charge**: fundamental property; two types (positive, negative); like repels, unlike attracts
- **Quantization**: q = ne; e = 1.602 × 10⁻¹⁹ C; charge comes in integer multiples of e
- **Conductors**: electrons move freely (metals); **insulators**: electrons fixed (rubber, glass)
- **Charging by friction**: electron transfer between rubbed surfaces
- **Charging by conduction**: charge flows at direct contact (same sign)
- **Charging by induction**: charge separation without contact (opposite sign)
- **Coulomb's Law**: F = kq₁q₂/r²; k = 8.99 × 10⁹ N·m²/C²
- Inverse-square law: F ∝ 1/r²; doubling distance → force reduced by 4×
- Electric force is 10³⁶ times stronger than gravity for two protons
- **Superposition**: net force = vector sum of all individual Coulomb forces

---

## Key Equations

```
Coulomb's Law:
  F = k × |q₁| × |q₂| / r²
  k = 8.99 × 10⁹ N·m²/C²
  k = 1 / (4πε₀)
  ε₀ = 8.85 × 10⁻¹² C²/(N·m²)

Elementary charge:
  e = 1.602 × 10⁻¹⁹ C

Charge quantization:
  q = n × e   (n is integer)

Superposition:
  F_total = F₁ + F₂ + F₃ + ...  (vector sum)
```
