# Chapter 11: Support Vector Machines — Maximum Margin Classifiers

> **"SVMs are one of the most theoretically beautiful algorithms in machine learning. They come with strong guarantees about generalization, they have a clear geometric interpretation, and the kernel trick is simply magical."**

---

## Table of Contents
1. [The Geometric Intuition](#1-the-geometric-intuition)
2. [Hyperplanes and Margins](#2-hyperplanes-and-margins)
3. [Support Vectors — The Critical Points](#3-support-vectors--the-critical-points)
4. [Hard Margin SVM](#4-hard-margin-svm)
5. [Soft Margin SVM — Handling Real Data](#5-soft-margin-svm--handling-real-data)
6. [Hinge Loss — What SVMs Actually Minimize](#6-hinge-loss--what-svms-actually-minimize)
7. [The Kernel Trick — The Big Idea](#7-the-kernel-trick--the-big-idea)
8. [Popular Kernel Functions](#8-popular-kernel-functions)
9. [SVM for Regression (SVR)](#9-svm-for-regression-svr)
10. [Multi-class SVM](#10-multi-class-svm)
11. [Sklearn: SVC and SVR](#11-sklearn-svc-and-svr)
12. [When SVMs Shine vs Struggle](#12-when-svms-shine-vs-struggle)
13. [Logistic Regression vs SVM](#13-logistic-regression-vs-svm)
14. [Full Example: Text Classification (Spam Detection)](#14-full-example-text-classification-spam-detection)
15. [Summary](#15-summary)
16. [Exercises](#16-exercises)

---

## 1. The Geometric Intuition

Imagine you have two classes of points on a 2D plane — red and blue dots. You want to draw a line that separates them. There are infinitely many valid lines. Which one should you choose?

```
THREE VALID DECISION BOUNDARIES
────────────────────────────────────────────────────────────────────────
  ●  ●                ●  ●                 ●  ●
●       ●           ●   / ●             ●  ○/    ●
  ●  ●               ●/    ●              ●/   ●
──────/──────       ─────────────      /──────────────
 ○  /  ○           ○  ○   ○           ○  ○   ○
  ○/     ○           ○   ○               ○   ○

Line A: too close  Line B: too close   Line C: equidistant
to red points      to blue points      from both classes
                                       ← THIS IS WHAT SVM FINDS

● = red class, ○ = blue class

The SVM principle: choose the decision boundary that maximizes the
margin (the "breathing room") between the two classes.

Intuition: if the boundary is close to points, a slight perturbation
in new test data (new point lands near the boundary) causes misclassification.
The wider the margin, the more robust the classifier.
```

This is the core SVM idea: **maximum margin classification**.

---

## 2. Hyperplanes and Margins

### What is a Hyperplane?

A **hyperplane** is a subspace of one dimension less than the ambient space:
- In 2D: a line (1D boundary in 2D space)
- In 3D: a plane (2D boundary in 3D space)
- In N-D: an (N-1)-dimensional subspace

A hyperplane in N-dimensional space is defined by:

```
w^T x + b = 0

Where:
  x ∈ ℝ^N  = a point in the feature space
  w ∈ ℝ^N  = the normal vector (perpendicular to the hyperplane)
  b ∈ ℝ    = the bias (shifts the hyperplane)

For 2D: w₁x₁ + w₂x₂ + b = 0  (a line)
For 3D: w₁x₁ + w₂x₂ + w₃x₃ + b = 0  (a plane)

Points ABOVE the hyperplane: w^T x + b > 0 → predict class +1
Points BELOW the hyperplane: w^T x + b < 0 → predict class -1
```

### The Margin

The margin is the distance between the decision boundary and the nearest points from each class.

```
MARGIN GEOMETRY
────────────────────────────────────────────────────────────────────────
  x₂
  │      ●           ●
  │  ●       ●           margin = 2/||w||
  │    ●           ●     ←──────────────→
  │         ─────────────────────────────  ← decision boundary (w^T x + b = 0)
  │    ○      ─────────────────────────    ← w^T x + b = -1 (negative margin plane)
  │      ○  ○              ────────────    ← w^T x + b = +1 (positive margin plane)
  │    ○         ○
  │
  └──────────────────────────────────────── x₁

  Support vectors: the circled/bold points ON the margin planes
  Margin width = 2 / ||w||
  Maximizing margin = minimizing ||w||
```

We define the margin planes as:
- Positive class: w^T x + b ≥ +1 for all positive samples
- Negative class: w^T x + b ≤ -1 for all negative samples

The margin width is:

```
Margin = 2 / ||w||   (distance between the two margin planes)

||w|| = √(w₁² + w₂² + ... + w_N²) = Euclidean norm of w

To MAXIMIZE margin → MINIMIZE ||w|| → MINIMIZE (1/2)||w||²
(the ½ and squared are for mathematical convenience)
```

---

## 3. Support Vectors — The Critical Points

The **support vectors** are the training examples that lie ON the margin boundaries (where w^T x + b = ±1). They are the points that "support" the margin — remove any other point and the margin wouldn't change. Remove a support vector and the entire hyperplane might shift.

```
SUPPORT VECTORS HIGHLIGHTED
────────────────────────────────────────────────────────────────────────
  ●         ●              ← regular points (inside or far from margin)
      ●◎          ●        ← ◎ = support vector (ON the margin line +1)
                  ●
   ─────────────────────── decision boundary (w^T x + b = 0)
                  ─────
         ◎─────────────    ← ◎ = support vector (ON the margin line -1)
     ○         ○
  ○     ○

Key property:
  Only support vectors matter for defining the hyperplane.
  If you remove any non-support-vector point, the SVM solution stays identical.
  The number of support vectors is much smaller than n (often 5-20%).
  This is what makes SVMs memory-efficient at test time.
```

**Why this matters for generalization theory:** SVMs have a theoretical bound on generalization error related to the number of support vectors, not the number of dimensions. This is the geometric basis for SVMs being effective in high-dimensional spaces (like text classification where p >> n).

---

## 4. Hard Margin SVM

The **hard margin SVM** assumes the data is linearly separable — you can draw a line (or hyperplane) that perfectly separates the two classes with no misclassifications.

### Formal Optimization Problem

```
PRIMAL PROBLEM (Hard Margin SVM)
────────────────────────────────────────────────────────────────────────
Minimize:    (1/2) ||w||²

Subject to:  yᵢ(w^T xᵢ + b) ≥ 1    for all i = 1, ..., n

Where:
  yᵢ ∈ {-1, +1}  = class label
  xᵢ ∈ ℝ^p       = feature vector

Constraint interpretation:
  If yᵢ = +1: w^T xᵢ + b ≥ +1  (positive samples above positive margin)
  If yᵢ = -1: w^T xᵢ + b ≤ -1  (negative samples below negative margin)
  Compactly: yᵢ(w^T xᵢ + b) ≥ 1  (both in one inequality)

This is a QUADRATIC PROGRAMMING (QP) problem:
  - Quadratic objective: (1/2)||w||² is a quadratic function of w
  - Linear constraints: yᵢ(w^T xᵢ + b) ≥ 1
  - Convex problem → guaranteed global optimum
  - No local minima to get stuck in
```

### The Dual Formulation

The actual SVM solvers work in the **dual** form (derived via Lagrange multipliers):

```
DUAL PROBLEM
────────────────────────────────────────────────────────────────────────
Maximize:   Σᵢ αᵢ - (1/2) Σᵢ Σⱼ αᵢαⱼyᵢyⱼ (xᵢ^T xⱼ)

Subject to: Σᵢ αᵢyᵢ = 0,   αᵢ ≥ 0  for all i

Where αᵢ = Lagrange multiplier for training sample i

Key insight: the objective involves only DOT PRODUCTS xᵢ^T xⱼ
           → this is where the kernel trick enters!

After solving:
  w = Σᵢ αᵢyᵢxᵢ   (w is a linear combination of training samples)
  b = yⱼ - w^T xⱼ  (for any support vector j where αⱼ > 0)

Most αᵢ = 0 (non-support vectors don't contribute to w!)
```

---

## 5. Soft Margin SVM — Handling Real Data

Real data is almost never perfectly linearly separable. The **soft margin SVM** allows some points to violate the margin or even the decision boundary, but penalizes these violations.

### Slack Variables

```
SOFT MARGIN: ALLOWING VIOLATIONS
────────────────────────────────────────────────────────────────────────
Introduce slack variable ξᵢ ≥ 0 for each training sample:

  ξᵢ = 0:          Point correctly classified, outside margin
  0 < ξᵢ ≤ 1:      Point inside the margin but correctly classified
  ξᵢ > 1:          Point on the WRONG side of the decision boundary

  x₂
  │
  │   ●    ●                ●
  │ ●   ●          ξ=0 ↙       ● ← ξ > 0 (inside margin)
  │       ─────────────────────────  decision boundary
  │     ─────────────────────────    margin planes
  │  ○  ─────────
  │   ○   ○        ← ○ with ξ > 1 = on wrong side = misclassified
  │         ○
  └────────────────────────── x₁

Modified constraint: yᵢ(w^T xᵢ + b) ≥ 1 - ξᵢ

Modified objective:
  Minimize: (1/2)||w||² + C × Σᵢ ξᵢ

  First term:  maximize margin (geometric goal)
  Second term: minimize total violation (penalty for violations)
  C: the tradeoff hyperparameter
```

### The C Parameter

The C parameter is the most important hyperparameter in SVMs:

```
C PARAMETER EFFECTS
────────────────────────────────────────────────────────────────────────
Large C (e.g., 100):               Small C (e.g., 0.01):
  Penalize violations heavily        Allow many violations
  Few margin violations allowed      Prioritize wide margin
  → Narrow margin                    → Wide margin
  → Model fits training harder       → More regularization
  → High variance (overfitting)      → High bias (underfitting)

  ─────────────────────────         ─────────────────────────
  Narrow margin, more SVs            Wide margin, fewer SVs
  Fits training data tightly         May misclassify train points

Rule of thumb: start with C=1, tune with cross-validation
C = 1/λ where λ is the regularization strength (same idea as lasso/ridge!)
```

---

## 6. Hinge Loss — What SVMs Actually Minimize

The SVM objective can be rewritten in a loss function framework:

```
HINGE LOSS
────────────────────────────────────────────────────────────────────────
ℓ(yᵢ, ŷᵢ) = max(0, 1 - yᵢ(w^T xᵢ + b))

Where yᵢ ∈ {-1, +1} and ŷᵢ = w^T xᵢ + b (the "raw score")

Plot (for yᵢ = +1):
  Loss
  │
2 │╲
  │ ╲
1 │  ╲──────────────────────────
  │   ╲                          0 loss for all correctly classified
  │    ╲                         points with margin ≥ 1
0 │─────────────────────────────
  │                              ←─ yᵢ × ŷᵢ = 0 (on decision boundary)
  ─1    0    1    2    3   yᵢ × ŷᵢ

  Correctly classified by >1 margin: loss = 0
  On the boundary: loss = 1
  Misclassified: loss = 1 + |margin violation|

SVM LOSS = (1/2)||w||² + C × Σ max(0, 1 - yᵢ(w^T xᵢ + b))
           ─────────────────   ──────────────────────────────────
           Regularization        Hinge Loss (data-fitting term)
```

### Hinge Loss vs Logistic Loss

```
HINGE LOSS vs LOGISTIC LOSS
────────────────────────────────────────────────────────────────────────
Loss
│
│ ╲         ╲
│  ╲         ╲──── Logistic (log loss) — always positive, smooth
│   ╲
│    ╲──────── Hinge — exactly zero for correct + large margin
│
└──────────────────────────── yᵢ × ŷᵢ
        0         1

Key difference:
  Logistic loss: always positive, even for well-classified points
  Hinge loss: ZERO for correctly classified points with margin ≥ 1
             → SVM ignores non-support vectors!

This is why:
  Logistic Regression → gives calibrated probabilities
  SVM → finds the maximum margin boundary, no probability output
```

---

## 7. The Kernel Trick — The Big Idea

This is the most powerful concept in SVM theory — and arguably one of the most elegant ideas in all of machine learning.

### The Problem: Non-linearly Separable Data

```
XOR PROBLEM — NOT LINEARLY SEPARABLE
────────────────────────────────────────────────────────────────────────
x₂
│
1 │  ●        ○
  │
0 │  ○        ●
  └──────────────── x₁
     0        1

● = class +1 (when x₁ ≠ x₂)
○ = class -1 (when x₁ = x₂)

No straight line can separate these! A linear classifier fails.
We need a curved decision boundary.

Option 1: Use a more complex model
Option 2: Map to a higher-dimensional space where linear separation is possible
          → This is what kernels do
```

### The Feature Map Solution

Let's add a new feature: x₃ = x₁ × x₂ (the product of the original features).

```
WITH ADDED FEATURE x₃ = x₁ × x₂
────────────────────────────────────────────────────────────────────────
Original 2D → New 3D space

(0,0) → (0, 0, 0)   ○   class -1
(1,1) → (1, 1, 1)   ○   class -1
(0,1) → (0, 1, 0)   ●   class +1
(1,0) → (1, 0, 0)   ●   class +1

In 3D, the plane x₃ = 0.5 separates ● and ○ perfectly!
A linear classifier (a plane) in 3D = a nonlinear boundary in 2D.
```

**The key insight:** By mapping data to a higher-dimensional space with a function φ(x), we can make non-linearly separable data linearly separable, then apply a linear SVM.

### The Problem with Explicit Mappings

If the original feature space has p features, the degree-d polynomial feature map has O(p^d) features. For:
- p = 1000 features, d = 3: 10⁹ features
- p = 10000 (text), d = 2: 10⁸ features

This is computationally infeasible. We can't actually compute φ(x).

### The Kernel Trick: Never Compute φ(x) Explicitly

Recall the SVM dual objective:

```
Maximize: Σᵢ αᵢ - (1/2) Σᵢ Σⱼ αᵢαⱼyᵢyⱼ (xᵢ^T xⱼ)
```

The data appears ONLY as dot products xᵢ^T xⱼ. If we use the feature map φ, we'd compute φ(xᵢ)^T φ(xⱼ).

**The kernel trick:** If we can compute φ(xᵢ)^T φ(xⱼ) = k(xᵢ, xⱼ) efficiently WITHOUT explicitly computing φ, we can use any feature space, no matter how high-dimensional!

```
THE KERNEL TRICK
────────────────────────────────────────────────────────────────────────
k(xᵢ, xⱼ) = φ(xᵢ)^T φ(xⱼ)

We define a kernel function k that computes the dot product
in the high-dimensional space DIRECTLY from the original inputs.

Example: Polynomial kernel k(x, z) = (x^T z + 1)^d

For d=2, x = (x₁, x₂), z = (z₁, z₂):
  k(x, z) = (x₁z₁ + x₂z₂ + 1)²
           = x₁²z₁² + x₂²z₂² + 2x₁x₂z₁z₂ + 2x₁z₁ + 2x₂z₂ + 1

This is the dot product of:
  φ(x) = (x₁², x₂², √2 x₁x₂, √2 x₁, √2 x₂, 1)
  φ(z) = (z₁², z₂², √2 z₁z₂, √2 z₁, √2 z₂, 1)

Computing k(x,z) directly: 5 multiplications and additions
Computing φ(x)^T φ(z) explicitly: 6-dimensional dot product

For degree 3 with 1000 features:
  Direct kernel: 1001 operations
  Explicit mapping: 10⁹ operations

SAME RESULT, VASTLY DIFFERENT COMPUTATIONAL COST.
```

**Mercer's theorem** tells us exactly which functions are valid kernels (their Gram matrix must be positive semi-definite). The popular kernels below are all valid.

---

## 8. Popular Kernel Functions

### Linear Kernel

```
k(x, z) = x^T z

This is just the standard dot product. No feature mapping.
Use when: data is already linearly separable, or high-dimensional
          (text data: p >> n makes linear kernel effective).
```

### Polynomial Kernel

```
k(x, z) = (γ x^T z + r)^d

Parameters:
  d = degree (typically 2 or 3)
  γ = scale factor
  r = offset

Feature space: all monomials of degree ≤ d
Use when: non-linear, feature interactions matter
         (e.g., pixel interactions in images)
```

### RBF (Radial Basis Function) / Gaussian Kernel

The most widely used kernel:

```
k(x, z) = exp(-γ ||x - z||²)   where γ > 0

||x - z||² = squared Euclidean distance between x and z

Properties:
  k(x, z) ≈ 1 when x and z are close (similar points)
  k(x, z) ≈ 0 when x and z are far apart (dissimilar)

  Corresponds to an INFINITE-DIMENSIONAL feature space!
  (Taylor expansion of exp reveals infinite polynomial terms)

  γ parameter controls the "width" of the Gaussian:
  Large γ → narrow Gaussian → only very close points influence each other
           → complex, wiggly decision boundary → high variance
  Small γ → wide Gaussian → distant points influence each other
           → smooth decision boundary → high bias

Decision function:
  f(x) = Σᵢ αᵢyᵢ k(xᵢ, x) + b
        = Σᵢ αᵢyᵢ exp(-γ ||xᵢ - x||²) + b

  Like a weighted sum of Gaussians centered at each support vector!
```

```
RBF KERNEL VISUALIZATION
────────────────────────────────────────────────────────────────────────
Small γ (smooth):         Large γ (complex):
   x₂                        x₂
   │  ╭───────────────╮       │  ╭──╮  ╭──╮  ╭──╮
   │  │       ●       │       │  │● │  │●  │  │  │
   │  │     ●   ●     │       │  ╰──╯  ╰──╯  │●●│
   │  ╰───────────────╯       │        ╭──╮  ╰──╯
   │     ○  ○  ○  ○           │   ○    │  │
   └──────────────────── x₁   └──────────────────── x₁
   Few support vectors        Many support vectors
   Smooth boundary            Wiggly boundary → overfit
```

### Sigmoid Kernel

```
k(x, z) = tanh(γ x^T z + r)

Similar to neural network activation function.
Not always positive semi-definite (care needed).
Rarely used in practice.
```

---

## 9. SVM for Regression (SVR)

SVMs can also handle regression with the **ε-insensitive loss**:

```
SVR: ε-INSENSITIVE TUBE
────────────────────────────────────────────────────────────────────────
y (true)
│          ╭────────────────────╮ ε-tube (no loss inside tube)
│          │     ○              │
│          │          ○         │
│     ──────────────────────────────  predicted function
│          │       ○            │
│          ╰────────────────────╯
│   ○  ← outside tube → contributes to loss
└─────────────────────────────────── x

ε-insensitive loss:
  L_ε(y, ŷ) = max(0, |y - ŷ| - ε)

  Points inside the ε-tube: zero loss
  Points outside: loss proportional to distance from tube

Objective:
  Minimize: (1/2)||w||² + C × Σᵢ (ξᵢ + ξᵢ*)

Where ξᵢ, ξᵢ* are slack variables for points above/below the tube.

Parameters:
  ε: tube width (how much error we tolerate)
  C: regularization (same interpretation as SVC)
```

---

## 10. Multi-class SVM

SVMs are inherently binary classifiers. For K > 2 classes, two strategies:

```
ONE-VS-REST (OvR)              ONE-VS-ONE (OvO)
──────────────────             ──────────────────────
K binary classifiers           K(K-1)/2 binary classifiers
Class k vs. all others         Class i vs. Class j for all pairs

For K=10 classes:              For K=10 classes:
  10 classifiers                 45 classifiers

Prediction: highest             Prediction: majority vote
confidence score                among all pairwise classifiers

sklearn default: OvO            Faster for large n
(better accuracy,               (each classifier trains on
but more classifiers)           smaller subset of data)
```

---

## 11. Sklearn: SVC and SVR

```python
from sklearn.svm import SVC, SVR, LinearSVC
from sklearn.preprocessing import StandardScaler
from sklearn.model_selection import train_test_split, GridSearchCV
from sklearn.metrics import accuracy_score, roc_auc_score, classification_report
from sklearn.datasets import load_breast_cancer, make_circles
import numpy as np

# ====================================================================
# CRITICAL: SVMs require feature scaling
# ====================================================================
# SVM uses distances — features on different scales will dominate unfairly
# Always apply StandardScaler (or MinMaxScaler) before SVM

X, y = load_breast_cancer(return_X_y=True)
X_train, X_test, y_train, y_test = train_test_split(
    X, y, test_size=0.2, random_state=42, stratify=y
)

scaler = StandardScaler()
X_train_s = scaler.fit_transform(X_train)
X_test_s  = scaler.transform(X_test)

# ====================================================================
# 1. Linear SVC (for linearly separable / high-dimensional data)
# ====================================================================
# LinearSVC is faster than SVC(kernel='linear') for large datasets
# Uses liblinear instead of libsvm
linear_svc = LinearSVC(
    C=1.0,
    max_iter=10000,
    random_state=42
)
linear_svc.fit(X_train_s, y_train)
print(f"LinearSVC accuracy: {accuracy_score(y_test, linear_svc.predict(X_test_s)):.4f}")

# ====================================================================
# 2. SVC with RBF Kernel (nonlinear, most commonly used)
# ====================================================================
svc_rbf = SVC(
    C=1.0,           # Regularization: high C = narrow margin = more fit
    kernel='rbf',    # Kernel: 'linear', 'poly', 'rbf', 'sigmoid'
    gamma='scale',   # γ: 'scale' = 1/(n_features × X.var()), 'auto' = 1/n_features
    probability=True,  # Enable predict_proba (slower training)
    random_state=42
)
svc_rbf.fit(X_train_s, y_train)
y_proba = svc_rbf.predict_proba(X_test_s)[:, 1]
print(f"SVC RBF accuracy: {accuracy_score(y_test, svc_rbf.predict(X_test_s)):.4f}")
print(f"SVC RBF AUC-ROC:  {roc_auc_score(y_test, y_proba):.4f}")
print(f"Number of support vectors: {svc_rbf.n_support_}")

# ====================================================================
# 3. Hyperparameter Tuning: Grid Search over C and gamma
# ====================================================================
param_grid = {
    'C':     [0.01, 0.1, 1, 10, 100],
    'gamma': [0.001, 0.01, 0.1, 1, 'scale'],
    'kernel': ['rbf']
}

grid_search = GridSearchCV(
    SVC(probability=True, random_state=42),
    param_grid,
    cv=5,
    scoring='roc_auc',
    n_jobs=-1,
    verbose=0
)
grid_search.fit(X_train_s, y_train)

print(f"\nBest parameters: {grid_search.best_params_}")
print(f"Best CV AUC-ROC: {grid_search.best_score_:.4f}")
best_svc = grid_search.best_estimator_
print(f"Test AUC-ROC: {roc_auc_score(y_test, best_svc.predict_proba(X_test_s)[:, 1]):.4f}")

# ====================================================================
# 4. Visualizing the kernel effect: XOR data
# ====================================================================
np.random.seed(42)
n = 200
X_xor = np.random.randn(n, 2)
y_xor = np.array([1 if x[0]*x[1] > 0 else -1 for x in X_xor])  # XOR pattern

# Add noise
y_xor[np.random.choice(n, 20)] *= -1

X_xor_train, X_xor_test, y_xor_train, y_xor_test = train_test_split(
    X_xor, y_xor, test_size=0.3, random_state=42
)

# Linear SVM on XOR (should fail)
svc_lin = SVC(kernel='linear', C=1.0)
svc_lin.fit(X_xor_train, y_xor_train)
print(f"\nLinear SVM on XOR: {svc_lin.score(X_xor_test, y_xor_test):.4f}  (expect ~0.5)")

# Polynomial SVM on XOR (should work with degree 2)
svc_poly = SVC(kernel='poly', degree=2, C=1.0, gamma='scale')
svc_poly.fit(X_xor_train, y_xor_train)
print(f"Poly-2 SVM on XOR: {svc_poly.score(X_xor_test, y_xor_test):.4f}  (expect ~0.9+)")

# RBF SVM on XOR
svc_rbf2 = SVC(kernel='rbf', C=5.0, gamma='scale')
svc_rbf2.fit(X_xor_train, y_xor_train)
print(f"RBF SVM on XOR:  {svc_rbf2.score(X_xor_test, y_xor_test):.4f}  (expect ~0.9+)")

# ====================================================================
# 5. Non-linearly separable circles
# ====================================================================
X_circles, y_circles = make_circles(n_samples=400, noise=0.1, random_state=42)
X_c_train, X_c_test, y_c_train, y_c_test = train_test_split(
    X_circles, y_circles, test_size=0.3, random_state=42
)

scaler_c = StandardScaler()
X_c_train_s = scaler_c.fit_transform(X_c_train)
X_c_test_s  = scaler_c.transform(X_c_test)

# Linear SVM (should fail on circles)
svc_c_lin = SVC(kernel='linear').fit(X_c_train_s, y_c_train)
print(f"\nCircles - Linear SVM:  {svc_c_lin.score(X_c_test_s, y_c_test):.4f}  (expect ~0.5)")

# RBF SVM (should work)
svc_c_rbf = SVC(kernel='rbf', C=1.0, gamma=1.0).fit(X_c_train_s, y_c_train)
print(f"Circles - RBF SVM:   {svc_c_rbf.score(X_c_test_s, y_c_test):.4f}  (expect ~0.95+)")

# ====================================================================
# 6. SVR Example
# ====================================================================
from sklearn.svm import SVR
from sklearn.datasets import fetch_california_housing
from sklearn.metrics import r2_score, mean_squared_error

Xr, yr = fetch_california_housing(return_X_y=True)
# Use a subset for speed (SVM is O(n²) in training)
idx = np.random.choice(len(Xr), 2000, replace=False)
Xr, yr = Xr[idx], yr[idx]

Xr_train, Xr_test, yr_train, yr_test = train_test_split(
    Xr, yr, test_size=0.2, random_state=42
)

scaler_r = StandardScaler()
Xr_train_s = scaler_r.fit_transform(Xr_train)
Xr_test_s  = scaler_r.transform(Xr_test)

svr = SVR(
    kernel='rbf',
    C=10.0,       # Regularization
    gamma='scale',
    epsilon=0.1   # ε-tube width
)
svr.fit(Xr_train_s, yr_train)
print(f"\nSVR R²: {r2_score(yr_test, svr.predict(Xr_test_s)):.4f}")
print(f"SVR RMSE: {mean_squared_error(yr_test, svr.predict(Xr_test_s), squared=False):.4f}")
```

---

## 12. When SVMs Shine vs Struggle

### When SVMs Shine

```
SVM STRENGTHS
────────────────────────────────────────────────────────────────────────
1. HIGH-DIMENSIONAL DATA where p >> n:
   Text classification: 10,000+ features, thousands of samples
   Genomics: 100,000+ gene expression features
   SVMs have theoretical guarantees related to support vectors
   not to dimensionality of feature space.

2. SMALL DATASETS:
   SVMs can perform well with hundreds of samples
   whereas deep learning needs thousands or millions.

3. WELL-SEPARATED CLASSES:
   When classes have a clear margin, SVM finds it perfectly.

4. KERNEL TRICK enables:
   Non-linear boundaries without explicit feature engineering
   Can use domain-specific kernels (string kernels for sequences,
   graph kernels for molecules)

5. ROBUST TO OUTLIERS:
   Only support vectors matter — outliers far from boundary don't affect it.
   (Contrast: OLS linear regression is very sensitive to outliers)
```

### When SVMs Struggle

```
SVM WEAKNESSES
────────────────────────────────────────────────────────────────────────
1. LARGE DATASETS:
   Training complexity: O(n² p) to O(n³) for kernel SVM
   10,000 samples: fine
   100,000 samples: slow
   1,000,000 samples: impractical
   Solution: LinearSVC, approximate methods (SGDClassifier with hinge loss)

2. FEATURE SCALING REQUIRED:
   Unlike tree-based models, SVMs are distance-based
   Must scale all features before training

3. NO PROBABILITY OUTPUT (by default):
   SVC with probability=True uses Platt scaling (additional fitting)
   The probabilities are not always well-calibrated

4. HARD TO INTERPRET:
   The model is defined by support vectors and their weights
   No simple "feature importance" like trees provide

5. SENSITIVE TO KERNEL CHOICE AND HYPERPARAMETERS:
   Wrong C or gamma → bad results
   Requires careful cross-validation

6. MEMORY FOR KERNEL MATRIX:
   For n samples with RBF kernel: stores n×n kernel matrix
   At n=50,000: 50000² × 8 bytes = 20 GB!
```

---

## 13. Logistic Regression vs SVM

```
LOGISTIC REGRESSION vs SVM
────────────────────────────────────────────────────────────────────────
Property          | Logistic Regression  | SVM
──────────────────┼──────────────────────┼──────────────────────────────
Loss function     | Log loss             | Hinge loss
Output            | Probabilities (0-1)  | Class labels (or ±scores)
Probabilities     | Well-calibrated      | Requires Platt scaling
Decision boundary | Minimizes log loss   | Maximizes margin
Outlier effect    | Moderate             | Low (only SVs matter)
Kernel            | No (linear only)     | Yes (can go non-linear)
Training speed    | Fast (O(np))         | Slow (O(n² p))
Scale required    | Yes (recommended)    | Yes (required!)
Interpretability  | Coefficients         | Support vectors, harder
Feature import.   | Weights              | Weights (linear only)
When to prefer    | Large data, need     | Small/medium data,
                  | probabilities        | high-dim, non-linear
                  |                      | with kernels

Conceptual difference:
  LR: "minimize wrong predictions (with probability)"
  SVM: "draw the widest possible highway between classes"

In practice:
  Both perform similarly on linearly-separable high-dim data
  SVM + RBF kernel can solve problems LR (linear) can't
  LR + polynomial features ≈ SVM + polynomial kernel
  For most tabular data: tree ensembles beat both
```

---

## 14. Full Example: Text Classification (Spam Detection)

Text classification is where SVMs truly excel. Text features (TF-IDF) are high-dimensional and sparse — exactly the scenario where SVMs are theoretically strong.

```python
# =============================================================================
# FULL EXAMPLE: Email Spam Detection with TF-IDF + SVM
# =============================================================================

import numpy as np
import pandas as pd
from sklearn.datasets import fetch_20newsgroups
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.svm import LinearSVC, SVC
from sklearn.linear_model import LogisticRegression
from sklearn.model_selection import train_test_split, cross_val_score
from sklearn.metrics import (
    classification_report, accuracy_score, roc_auc_score, confusion_matrix
)
from sklearn.pipeline import Pipeline
import warnings
warnings.filterwarnings('ignore')

# =============================================================================
# STEP 1: Load Text Data
# =============================================================================
# 20 Newsgroups: 20,000 newsgroup posts across 20 categories
# We'll do binary classification: technology vs sports

print("Loading 20 Newsgroups dataset...")
categories = ['comp.graphics', 'sci.space', 'rec.sport.baseball', 'talk.politics.guns']

newsgroups_train = fetch_20newsgroups(subset='train', categories=categories, remove=('headers', 'footers', 'quotes'))
newsgroups_test  = fetch_20newsgroups(subset='test',  categories=categories, remove=('headers', 'footers', 'quotes'))

X_train_text = newsgroups_train.data
X_test_text  = newsgroups_test.data
y_train      = newsgroups_train.target
y_test       = newsgroups_test.target
target_names = newsgroups_train.target_names

print(f"Train: {len(X_train_text)} documents")
print(f"Test:  {len(X_test_text)} documents")
print(f"Categories: {target_names}")
print(f"Class distribution (train): {np.bincount(y_train)}")

print(f"\nSample document:")
print(X_train_text[0][:300])
print(f"Label: {target_names[y_train[0]]}")

# =============================================================================
# STEP 2: TF-IDF Feature Extraction
# =============================================================================
# TF-IDF: Term Frequency-Inverse Document Frequency
# For each word in each document:
#   TF = how often this word appears in THIS document
#   IDF = log(N / number of docs containing this word)
#        rare words get higher weight (they're more discriminative)
#   TF-IDF = TF × IDF

print("\n=== TF-IDF Feature Extraction ===")
tfidf = TfidfVectorizer(
    max_features=50000,     # Use top 50,000 words by frequency
    min_df=2,               # Ignore words appearing in <2 documents
    max_df=0.95,            # Ignore words in >95% of documents (stop words)
    sublinear_tf=True,      # Apply log(1 + TF) instead of raw TF
    strip_accents='unicode',
    analyzer='word',
    token_pattern=r'\b[a-zA-Z]{3,}\b',  # Words of 3+ letters only
    ngram_range=(1, 2)      # Unigrams AND bigrams (pairs of consecutive words)
)

X_train_tfidf = tfidf.fit_transform(X_train_text)  # shape: (n_train, n_vocab)
X_test_tfidf  = tfidf.transform(X_test_text)

print(f"TF-IDF matrix shape: {X_train_tfidf.shape}")
print(f"Sparsity: {1 - X_train_tfidf.nnz / (X_train_tfidf.shape[0] * X_train_tfidf.shape[1]):.3%}")
# Expected: ~99%+ sparse (most words don't appear in most documents)
# This is why SVMs work so well on text: they're efficient on sparse features

# Top discriminative words (highest IDF — rare, specific)
feature_names_tfidf = tfidf.get_feature_names_out()
idf_scores = tfidf.idf_
top_idx = np.argsort(idf_scores)[-20:][::-1]  # highest IDF
print(f"\nTop 20 discriminative words (highest IDF):")
print([feature_names_tfidf[i] for i in top_idx])

# =============================================================================
# STEP 3: Train Multiple Models
# =============================================================================
models = {
    'LinearSVC (C=0.1)': LinearSVC(C=0.1, max_iter=5000, random_state=42),
    'LinearSVC (C=1.0)': LinearSVC(C=1.0, max_iter=5000, random_state=42),
    'LinearSVC (C=10.)': LinearSVC(C=10.,  max_iter=5000, random_state=42),
    'LogisticReg':       LogisticRegression(C=1.0, max_iter=1000, random_state=42),
}

print("\n" + "="*55)
print(f"{'Model':30s} {'Accuracy':>10} {'Macro F1':>10}")
print("="*55)

results = {}
for name, model in models.items():
    model.fit(X_train_tfidf, y_train)
    y_pred = model.predict(X_test_tfidf)
    acc = accuracy_score(y_test, y_pred)
    from sklearn.metrics import f1_score
    f1  = f1_score(y_test, y_pred, average='macro')
    results[name] = (acc, f1, model, y_pred)
    print(f"{name:30s} {acc:10.4f} {f1:10.4f}")
print("="*55)

# =============================================================================
# STEP 4: Detailed Report on Best Model
# =============================================================================
best_name = max(results.items(), key=lambda x: x[1][1])[0]
best_acc, best_f1, best_model, best_pred = results[best_name]

print(f"\n=== Detailed Report: {best_name} ===")
print(classification_report(y_test, best_pred, target_names=target_names))

print("Confusion Matrix:")
cm = confusion_matrix(y_test, best_pred)
print(pd.DataFrame(cm, index=target_names, columns=[f'pred_{c}' for c in target_names]))

# =============================================================================
# STEP 5: Most Discriminative Features (SVM Interpretation)
# =============================================================================
if isinstance(best_model, LinearSVC):
    print("\n=== Most Discriminative Words per Class ===")
    for class_idx, class_name in enumerate(target_names):
        # Get coefficient for this class (weights in the hyperplane)
        if best_model.coef_.shape[0] == len(target_names):
            # OvR multi-class: one weight vector per class
            coef = best_model.coef_[class_idx]
        else:
            coef = best_model.coef_[0]

        top_pos_idx = np.argsort(coef)[-10:][::-1]
        print(f"\n{class_name}:")
        print(f"  Top words: {[feature_names_tfidf[i] for i in top_pos_idx]}")

# =============================================================================
# STEP 6: Build a Spam-like Binary Classifier Using Pipeline
# =============================================================================
# Binary classification: tech (comp.graphics + sci.space) vs sports/politics
y_train_binary = (y_train < 2).astype(int)   # 0=tech, 1=not-tech
y_test_binary  = (y_test < 2).astype(int)

# Build a Pipeline (preprocessing + model as one object)
spam_pipeline = Pipeline([
    ('tfidf', TfidfVectorizer(
        max_features=30000,
        min_df=2,
        max_df=0.9,
        sublinear_tf=True,
        ngram_range=(1, 2)
    )),
    ('svm', LinearSVC(C=1.0, max_iter=5000, random_state=42))
])

spam_pipeline.fit(X_train_text, y_train_binary)
y_pred_binary = spam_pipeline.predict(X_test_text)

print("\n=== Binary Classification (Tech vs Non-Tech) ===")
print(classification_report(
    y_test_binary, y_pred_binary,
    target_names=['technology', 'other']
))

# Cross-validation on the pipeline
cv_scores = cross_val_score(
    spam_pipeline, X_train_text, y_train_binary,
    cv=5, scoring='f1', n_jobs=-1
)
print(f"5-Fold CV F1: {cv_scores.mean():.4f} ± {cv_scores.std():.4f}")

# =============================================================================
# STEP 7: Predict on Custom Text
# =============================================================================
custom_texts = [
    "The space shuttle launched successfully from Kennedy Space Center",
    "The Yankees beat the Red Sox in extra innings last night",
    "GPU memory bandwidth is critical for rendering performance",
    "Should citizens be allowed to carry handguns for self-defense?",
]

print("\n=== Custom Text Predictions ===")
preds = spam_pipeline.predict(custom_texts)
label_map = {0: 'TECHNOLOGY', 1: 'OTHER'}
for text, pred in zip(custom_texts, preds):
    print(f"  {label_map[pred]}: {text[:60]}...")

# =============================================================================
# STEP 8: SVM vs Logistic Regression Comparison Summary
# =============================================================================
print("\n=== SVM vs Logistic Regression on Text ===")
lr_pipeline = Pipeline([
    ('tfidf', TfidfVectorizer(max_features=30000, min_df=2, sublinear_tf=True)),
    ('lr',   LogisticRegression(C=1.0, max_iter=1000, n_jobs=-1))
])
lr_pipeline.fit(X_train_text, y_train_binary)
lr_preds = lr_pipeline.predict(X_test_text)
lr_acc = accuracy_score(y_test_binary, lr_preds)
lr_f1  = f1_score(y_test_binary, lr_preds)

svm_acc = accuracy_score(y_test_binary, y_pred_binary)
svm_f1  = f1_score(y_test_binary, y_pred_binary)

print(f"{'Model':20s} {'Accuracy':>10} {'F1':>10}")
print(f"{'LinearSVC':20s} {svm_acc:10.4f} {svm_f1:10.4f}")
print(f"{'LogisticReg':20s} {lr_acc:10.4f} {lr_f1:10.4f}")
print("\nConclusion: On text data, LinearSVC often matches or beats LogisticReg")
print("Both work well because TF-IDF creates a high-dimensional linear space")
```

---

## 15. Summary

```
CHAPTER 11 KEY CONCEPTS
─────────────────────────────────────────────────────────────

SVM CORE IDEA:
  Find the decision boundary with maximum margin.
  Margin = 2/||w|| → maximize margin = minimize ||w||²

SUPPORT VECTORS:
  Points lying ON the margin planes (w^T x + b = ±1)
  Only these points define the model
  Small fraction of training data

HARD MARGIN:
  Assumes linearly separable data
  Minimize (1/2)||w||² subject to yᵢ(w^T xᵢ + b) ≥ 1

SOFT MARGIN:
  Allows violations via slack variables ξᵢ ≥ 0
  Minimize (1/2)||w||² + C × Σξᵢ
  C = inverse of regularization strength
  Large C → narrow margin → low bias, high variance

HINGE LOSS:
  ℓ = max(0, 1 - y(w^T x + b))
  Zero for correct predictions with sufficient margin
  Only support vectors contribute to training

KERNEL TRICK:
  Replace dot product xᵢ^T xⱼ with k(xᵢ, xⱼ) = φ(xᵢ)^T φ(xⱼ)
  Never compute φ explicitly!
  RBF kernel: k(x,z) = exp(-γ||x-z||²) — infinite-dim feature space

KEY HYPERPARAMETERS:
  C: margin width control (like 1/λ in linear models)
  gamma: RBF kernel width (large γ → complex boundary)
  kernel: 'linear', 'poly', 'rbf' (default: 'rbf')

WHEN TO USE SVM:
  + High-dimensional data (text, genomics)
  + Small/medium datasets
  + Non-linear problems with kernel
  - Large datasets (too slow)
  - Need probabilities
  - Need feature importance
```

---

## Mini Projects

### Mini Project 1: Kernel Trick Visualizer

Compare how linear, RBF, and polynomial kernels draw decision boundaries on the same 2D dataset.

**Objective:** Build intuition for why kernel choice matters.

```python
import numpy as np
import matplotlib.pyplot as plt
from sklearn.svm import SVC
from sklearn.datasets import make_moons, make_circles, make_classification
from sklearn.preprocessing import StandardScaler

def plot_decision_boundary(ax, clf, X, y, title):
    h = 0.02
    x_min, x_max = X[:, 0].min() - 0.5, X[:, 0].max() + 0.5
    y_min, y_max = X[:, 1].min() - 0.5, X[:, 1].max() + 0.5
    xx, yy = np.meshgrid(np.arange(x_min, x_max, h),
                         np.arange(y_min, y_max, h))
    Z = clf.predict(np.c_[xx.ravel(), yy.ravel()])
    Z = Z.reshape(xx.shape)
    ax.contourf(xx, yy, Z, alpha=0.3, cmap='RdYlBu')
    ax.scatter(X[:, 0], X[:, 1], c=y, cmap='RdYlBu', edgecolors='k', s=40)
    # Highlight support vectors
    sv = clf.support_vectors_
    ax.scatter(sv[:, 0], sv[:, 1], s=150, facecolors='none',
               edgecolors='black', linewidths=2, zorder=5, label='Support Vectors')
    ax.set_title(f"{title}\n(SVs: {len(sv)}, Acc: {clf.score(X, y):.2f})")
    ax.legend(fontsize=7)

# Three datasets that stress-test different kernels
datasets = [
    (make_moons(n_samples=200, noise=0.15, random_state=42), "Moons"),
    (make_circles(n_samples=200, noise=0.1, factor=0.5, random_state=42), "Circles"),
    (make_classification(n_samples=200, n_features=2, n_redundant=0,
                         n_informative=2, random_state=42), "Linear"),
]

kernels = [
    ('linear', {}),
    ('rbf', {'gamma': 'scale'}),
    ('poly', {'degree': 3, 'gamma': 'scale'}),
]

fig, axes = plt.subplots(3, 3, figsize=(15, 12))
fig.suptitle("SVM Kernel Comparison", fontsize=16, fontweight='bold')

for row, ((X, y), ds_name) in enumerate(datasets):
    scaler = StandardScaler()
    X = scaler.fit_transform(X)
    for col, (kernel, params) in enumerate(kernels):
        clf = SVC(kernel=kernel, C=1.0, **params)
        clf.fit(X, y)
        plot_decision_boundary(axes[row, col], clf, X, y,
                               f"{ds_name} | kernel='{kernel}'")

plt.tight_layout()
plt.savefig("kernel_comparison.png", dpi=150)
plt.show()
print("Saved: kernel_comparison.png")
```

**Key insight:** RBF handles non-linear data; linear is fastest for high-dimensional text; poly is great for image features.

---

### Mini Project 2: SMS Spam Classifier

Build a text classification pipeline using TF-IDF + SVM and compare different kernels.

**Objective:** Apply SVMs to a real NLP task, measure precision/recall tradeoffs.

```python
import numpy as np
import pandas as pd
from sklearn.svm import SVC, LinearSVC
from sklearn.feature_extraction.text import TfidfVectorizer
from sklearn.model_selection import train_test_split, cross_val_score
from sklearn.pipeline import Pipeline
from sklearn.metrics import classification_report, confusion_matrix
import matplotlib.pyplot as plt
import seaborn as sns

# Synthetic SMS dataset (replace with real SMS Spam Collection if available)
spam_messages = [
    "WINNER!! You have won a free prize. Call now to claim!",
    "Congratulations! You've been selected for a $1000 gift card.",
    "FREE entry in 2 a wkly comp. Txt WIN to 87121 NOW!",
    "You have been awarded a free ringtone, reply to claim",
    "Urgent! Your mobile account has been compromised. Call immediately.",
    "Claim your prize now! Limited time offer expires tonight.",
    "You won the lottery! Send your bank details to collect.",
    "FREE msg: Txt STOP to opt out. Ringtone offer valid 7 days.",
]

ham_messages = [
    "Hey, are you coming to the party tonight?",
    "Can you pick up some groceries on your way home?",
    "Meeting rescheduled to 3pm tomorrow. Let me know if that works.",
    "Happy birthday! Hope you have a great day.",
    "Just finished the report, will send it over soon.",
    "Did you see the game last night? It was incredible!",
    "Lunch at noon? There's a new Thai place nearby.",
    "Running 10 minutes late, go ahead and order for me.",
    "Thanks for helping out yesterday, really appreciated it.",
    "The kids are at school, I'll be home by 6.",
]

# Expand dataset with variations for more data
def augment(messages, n=50):
    return [m + f" ref:{i}" for i in range(n) for m in messages][:200]

all_spam = augment(spam_messages)
all_ham = augment(ham_messages)

texts = all_spam + all_ham
labels = ['spam'] * len(all_spam) + ['ham'] * len(all_ham)

X_train, X_test, y_train, y_test = train_test_split(
    texts, labels, test_size=0.2, random_state=42, stratify=labels
)

# Compare different SVM configurations
configs = {
    'LinearSVC': Pipeline([
        ('tfidf', TfidfVectorizer(ngram_range=(1, 2), max_features=5000)),
        ('clf', LinearSVC(C=1.0, max_iter=1000))
    ]),
    'SVC-RBF': Pipeline([
        ('tfidf', TfidfVectorizer(ngram_range=(1, 2), max_features=5000)),
        ('clf', SVC(kernel='rbf', C=10, gamma='scale', probability=True))
    ]),
    'SVC-Linear': Pipeline([
        ('tfidf', TfidfVectorizer(ngram_range=(1, 2), max_features=5000)),
        ('clf', SVC(kernel='linear', C=1.0, probability=True))
    ]),
}

fig, axes = plt.subplots(1, 3, figsize=(15, 4))

for idx, (name, pipeline) in enumerate(configs.items()):
    pipeline.fit(X_train, y_train)
    y_pred = pipeline.predict(X_test)
    cv_scores = cross_val_score(pipeline, texts, labels, cv=5, scoring='f1_macro')

    print(f"\n{'='*40}")
    print(f"Model: {name}")
    print(f"CV F1 Score: {cv_scores.mean():.3f} ± {cv_scores.std():.3f}")
    print(classification_report(y_test, y_pred))

    cm = confusion_matrix(y_test, y_pred, labels=['spam', 'ham'])
    sns.heatmap(cm, annot=True, fmt='d', cmap='Blues', ax=axes[idx],
                xticklabels=['spam', 'ham'], yticklabels=['spam', 'ham'])
    axes[idx].set_title(f"{name}\nCV F1: {cv_scores.mean():.3f}")
    axes[idx].set_xlabel('Predicted')
    axes[idx].set_ylabel('True')

plt.tight_layout()
plt.savefig("svm_spam_classifier.png", dpi=150)
plt.show()

# Show top spam-indicating words
print("\n--- Top Spam-Indicating Features (LinearSVC) ---")
linear_svc_pipe = configs['LinearSVC']
linear_svc_pipe.fit(X_train, y_train)
vectorizer = linear_svc_pipe.named_steps['tfidf']
clf = linear_svc_pipe.named_steps['clf']
feature_names = vectorizer.get_feature_names_out()
coef = clf.coef_[0]
top_spam = sorted(zip(coef, feature_names), reverse=True)[:10]
top_ham = sorted(zip(coef, feature_names))[:10]
print("Spam indicators:", [f for _, f in top_spam])
print("Ham indicators: ", [f for _, f in top_ham])
```

---

### Mini Project 3: C Parameter and Margin Visualizer

See how the regularization parameter C controls the bias-variance tradeoff in SVMs.

**Objective:** Understand soft-margin SVM and when to use high vs. low C.

```python
import numpy as np
import matplotlib.pyplot as plt
from sklearn.svm import SVC
from sklearn.datasets import make_classification
from sklearn.model_selection import validation_curve
from sklearn.preprocessing import StandardScaler

np.random.seed(42)
X, y = make_classification(n_samples=150, n_features=2, n_redundant=0,
                            n_informative=2, n_clusters_per_class=1,
                            flip_y=0.1, random_state=42)
scaler = StandardScaler()
X = scaler.fit_transform(X)

# Add a noisy outlier to make C choice interesting
X = np.vstack([X, [[0.5, 0.3]]])
y = np.append(y, 1)

fig, axes = plt.subplots(2, 4, figsize=(20, 10))
C_values = [0.01, 0.1, 1.0, 100.0]

def plot_svm_margin(ax, clf, X, y, C_val):
    h = 0.02
    x_min, x_max = X[:, 0].min() - 0.5, X[:, 0].max() + 0.5
    y_min, y_max = X[:, 1].min() - 0.5, X[:, 1].max() + 0.5
    xx, yy = np.meshgrid(np.arange(x_min, x_max, h),
                         np.arange(y_min, y_max, h))
    Z = clf.decision_function(np.c_[xx.ravel(), yy.ravel()])
    Z = Z.reshape(xx.shape)

    ax.contourf(xx, yy, Z, levels=[-np.inf, -1, 0, 1, np.inf],
                colors=['#FFAAAA', '#FFDDDD', '#DDDDFF', '#AAAAFF'], alpha=0.6)
    ax.contour(xx, yy, Z, levels=[-1, 0, 1],
               linestyles=['--', '-', '--'], colors=['red', 'black', 'blue'], linewidths=1.5)
    ax.scatter(X[:, 0], X[:, 1], c=y, cmap='RdBu', edgecolors='k', s=40, zorder=3)
    sv = clf.support_vectors_
    ax.scatter(sv[:, 0], sv[:, 1], s=200, facecolors='none',
               edgecolors='green', linewidths=2.5, zorder=4)
    margin = 2 / np.linalg.norm(clf.coef_)
    n_sv = len(sv)
    ax.set_title(f"C={C_val}\nMargin={margin:.2f}, SVs={n_sv}", fontsize=10)

for idx, C in enumerate(C_values):
    clf = SVC(kernel='linear', C=C)
    clf.fit(X, y)
    plot_svm_margin(axes[0, idx], clf, X, y, C)
    axes[0, idx].set_xlabel("Feature 1")
    if idx == 0:
        axes[0, idx].set_ylabel("Feature 2\n(Decision boundary view)")

# Bottom row: validation curve (C vs train/val accuracy)
C_range = np.logspace(-3, 3, 20)
train_scores, val_scores = validation_curve(
    SVC(kernel='rbf', gamma='scale'), X, y,
    param_name='C', param_range=C_range,
    cv=5, scoring='accuracy', n_jobs=-1
)

ax_curve = axes[1, 0]
ax_curve.semilogx(C_range, train_scores.mean(axis=1), 'b-o', label='Train', markersize=4)
ax_curve.fill_between(C_range,
                      train_scores.mean(1) - train_scores.std(1),
                      train_scores.mean(1) + train_scores.std(1), alpha=0.15, color='blue')
ax_curve.semilogx(C_range, val_scores.mean(axis=1), 'r-o', label='Validation', markersize=4)
ax_curve.fill_between(C_range,
                      val_scores.mean(1) - val_scores.std(1),
                      val_scores.mean(1) + val_scores.std(1), alpha=0.15, color='red')
best_C = C_range[val_scores.mean(axis=1).argmax()]
ax_curve.axvline(best_C, color='green', linestyle=':', label=f'Best C={best_C:.2f}')
ax_curve.set_xlabel("C (regularization)")
ax_curve.set_ylabel("Accuracy")
ax_curve.set_title("Validation Curve: C vs Accuracy (RBF)")
ax_curve.legend()
ax_curve.grid(True, alpha=0.3)

# Gamma validation curve
gamma_range = np.logspace(-4, 2, 20)
train_g, val_g = validation_curve(
    SVC(kernel='rbf', C=1.0), X, y,
    param_name='gamma', param_range=gamma_range,
    cv=5, scoring='accuracy', n_jobs=-1
)
ax_gamma = axes[1, 1]
ax_gamma.semilogx(gamma_range, train_g.mean(1), 'b-o', label='Train', markersize=4)
ax_gamma.semilogx(gamma_range, val_g.mean(1), 'r-o', label='Validation', markersize=4)
best_gamma = gamma_range[val_g.mean(1).argmax()]
ax_gamma.axvline(best_gamma, color='green', linestyle=':', label=f'Best γ={best_gamma:.3f}')
ax_gamma.set_xlabel("Gamma")
ax_gamma.set_ylabel("Accuracy")
ax_gamma.set_title("Validation Curve: Gamma vs Accuracy (RBF)")
ax_gamma.legend()
ax_gamma.grid(True, alpha=0.3)

# Number of support vectors vs C
n_svs = []
for C in C_range:
    clf_tmp = SVC(kernel='rbf', C=C, gamma='scale')
    clf_tmp.fit(X, y)
    n_svs.append(len(clf_tmp.support_vectors_))
ax_sv = axes[1, 2]
ax_sv.semilogx(C_range, n_svs, 'g-o', markersize=4)
ax_sv.set_xlabel("C")
ax_sv.set_ylabel("# Support Vectors")
ax_sv.set_title("C vs Support Vector Count\n(high C = fewer SVs = less regularization)")
ax_sv.grid(True, alpha=0.3)

axes[1, 3].axis('off')
summary = ("SVM Key Takeaways:\n\n"
           "Low C (e.g. 0.01):\n"
           "  • Wide margin\n  • More support vectors\n"
           "  • Tolerates misclassifications\n  • Better generalization\n\n"
           "High C (e.g. 100):\n"
           "  • Narrow margin\n  • Fewer support vectors\n"
           "  • Tries to classify all points\n  • Risk of overfitting\n\n"
           "Low Gamma:\n  • Smoother boundary\n  • Far points influence model\n\n"
           "High Gamma:\n  • Complex boundary\n  • Only nearby points matter")
axes[1, 3].text(0.05, 0.95, summary, transform=axes[1, 3].transAxes,
                fontsize=9, verticalalignment='top', fontfamily='monospace',
                bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))

plt.suptitle("SVM: C Parameter, Margin, and Hyperparameter Tuning", fontsize=14, fontweight='bold')
plt.tight_layout()
plt.savefig("svm_margin_visualizer.png", dpi=150)
plt.show()
print("Saved: svm_margin_visualizer.png")
```

---

## 16. Exercises

**Exercise 1:** Implement a 2D SVM visualization. Generate a linearly-separable 2D dataset with 100 points. Train an SVM with `kernel='linear'`. Write code to:
- Draw the training points (colored by class)
- Draw the decision boundary
- Draw the two margin planes
- Highlight the support vectors with larger markers
Observe how the support vectors change as you modify C.

**Exercise 2:** Kernel comparison experiment. Generate the classic "two concentric circles" dataset with `make_circles(noise=0.1)`. Train SVMs with:
- Linear kernel
- Polynomial kernel (degree 2, 3, 5)
- RBF kernel (gamma = 0.1, 1, 10)
Plot the test accuracy and decision boundary for each. Which kernel works best? Why?

**Exercise 3:** The C parameter and the margin. On a simple 2D dataset:
- For C values [0.01, 0.1, 1, 10, 100]:
  - Count the number of support vectors
  - Record training and test accuracy
  - Plot the decision boundary
- Observe: as C increases, what happens to the number of support vectors?
- Explain the relationship between C, number of SVs, and overfitting.

**Exercise 4:** Scaling experiment. Take the breast cancer dataset. Train SVC with RBF kernel:
- Without any feature scaling
- With StandardScaler
- With MinMaxScaler
- With RobustScaler
Compare test accuracies. How critical is scaling for SVMs? Write a short analysis.

**Exercise 5:** Text classification from scratch. Download or use a spam email dataset (e.g., SMS Spam Collection from UCI ML Repository). Build a pipeline:
1. TfidfVectorizer with optimal parameters
2. LinearSVC with optimal C
Use GridSearchCV to tune both the TF-IDF parameters (`max_features`, `ngram_range`) and the SVM parameter (`C`) jointly. Report your final F1 score.

---

**Next Chapter →** [Chapter 12: Unsupervised Learning](./12-unsupervised-learning.md)

*Everything so far has required labeled data. Most data in the world is unlabeled. Unsupervised learning finds structure where you have no targets — clustering customers, reducing dimensions, detecting anomalies.*
