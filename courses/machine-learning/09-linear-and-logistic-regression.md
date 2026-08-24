# Chapter 09: Linear and Logistic Regression — Learning from Labeled Data

> **"All models are wrong, but some are useful."**
> — George Box

---

## Table of Contents
1. [Linear Regression](#1-linear-regression)
2. [Loss Functions for Regression](#2-loss-functions-for-regression)
3. [Ordinary Least Squares — The Closed-Form Solution](#3-ordinary-least-squares--the-closed-form-solution)
4. [Gradient Descent for Linear Regression](#4-gradient-descent-for-linear-regression)
5. [Linear Regression from Scratch](#5-linear-regression-from-scratch)
6. [Assumptions and Diagnostics](#6-assumptions-and-diagnostics)
7. [Polynomial Regression](#7-polynomial-regression)
8. [Regularization: Ridge, Lasso, Elastic Net](#8-regularization-ridge-lasso-elastic-net)
9. [Evaluating Regression Models](#9-evaluating-regression-models)
10. [Logistic Regression — Classification with Probabilities](#10-logistic-regression--classification-with-probabilities)
11. [Binary Cross-Entropy Loss](#11-binary-cross-entropy-loss)
12. [Multi-class Logistic Regression: Softmax](#12-multi-class-logistic-regression-softmax)
13. [Classification Evaluation Metrics](#13-classification-evaluation-metrics)
14. [Full Example: Credit Card Default Prediction](#14-full-example-credit-card-default-prediction)
15. [Summary](#15-summary)
16. [Exercises](#16-exercises)

---

## 1. Linear Regression

### The Problem

Linear regression solves the simplest possible supervised learning problem: given some input features, predict a **continuous number**.

Real examples:
- Given a house's square footage, number of bedrooms, and location → predict its sale price
- Given today's temperature, humidity, and season → predict tomorrow's temperature
- Given a company's marketing spend → predict quarterly sales

The "linear" in linear regression means the model assumes the output is a **linear combination** of the inputs.

### The Linear Model

For a single feature:

```
ŷ = w₁x₁ + w₀

Where:
  ŷ  = predicted value (y-hat)
  x₁ = input feature
  w₁ = weight (slope) — learned from data
  w₀ = bias/intercept — learned from data
```

For multiple features (the general case):

```
ŷ = w₀ + w₁x₁ + w₂x₂ + ... + wₚxₚ

In vector notation:
ŷ = Xw + b

Where:
  X ∈ ℝ^(n×p)  = feature matrix (n samples, p features)
  w ∈ ℝ^p      = weight vector
  b ∈ ℝ        = bias (scalar)
  ŷ ∈ ℝ^n      = predictions

Or more compactly, if we add a column of 1s to X:
ŷ = Xw           (bias absorbed into weight vector)
```

### What the Weights Mean

Weights are the model's learned understanding of the data. For the house price model:

```
price = w₀ + w₁ × (sq_ft) + w₂ × (bedrooms) + w₃ × (bathrooms)

If learned weights are:
  w₀ = 50,000   (base price for any house)
  w₁ = 120      (each sq ft adds $120)
  w₂ = 15,000   (each bedroom adds $15,000)
  w₃ = 25,000   (each bathroom adds $25,000)

Then: 1500 sq ft, 3 beds, 2 baths →
  price = 50,000 + 120×1500 + 15,000×3 + 25,000×2
        = 50,000 + 180,000 + 45,000 + 50,000
        = $325,000
```

The question is: **how do we learn these weights from data?** That requires defining what "good weights" means — which requires a loss function.

---

## 2. Loss Functions for Regression

A **loss function** measures how wrong the model's predictions are. Training is the process of finding weights that minimize the loss.

### Mean Squared Error (MSE) — The Standard

```
       1   n
MSE = ─── Σ (yᵢ - ŷᵢ)²
       n  i=1

Where:
  n   = number of training samples
  yᵢ  = true value for sample i
  ŷᵢ  = model's predicted value for sample i
  (yᵢ - ŷᵢ) = residual (the error)
```

**Why squared?**

1. **Sign cancellation:** Without squaring, positive errors and negative errors cancel out. A model that predicts +100 on half the samples and -100 on the other half has zero average error but is terrible. Squaring makes all errors positive.

2. **Differentiability:** MSE is smooth and differentiable everywhere — this is essential for gradient descent. The absolute value function has a kink at zero that makes optimization harder.

3. **Penalizes large errors more:** A prediction that's off by 10 contributes 100 to the loss. A prediction off by 1 contributes 1. This makes the model prioritize getting big outliers right. Whether this is desirable depends on your problem.

```
COMPARING RESIDUALS vs SQUARED RESIDUALS
──────────────────────────────────────────────────────────
Error = 1  →  Squared error =   1    (same magnitude)
Error = 2  →  Squared error =   4    (2× error, 4× loss)
Error = 5  →  Squared error =  25    (5× error, 25× loss)
Error = 10 →  Squared error = 100   (10× error, 100× loss)

Large errors are punished QUADRATICALLY.
This is great when outliers are actually wrong predictions to avoid.
This is bad when your data has legitimate outliers (use MAE instead).
```

### Mean Absolute Error (MAE)

```
       1   n
MAE = ─── Σ |yᵢ - ŷᵢ|
       n  i=1
```

- More robust to outliers than MSE (errors scale linearly, not quadratically)
- Not differentiable at zero (subgradient methods needed)
- Better when outliers in your data are legitimate (e.g., some houses legitimately sell for 10× the median)

### Huber Loss — Best of Both

```
         ⎧ ½(yᵢ - ŷᵢ)²           if |yᵢ - ŷᵢ| ≤ δ
L_δ(y, ŷ) = ⎨
         ⎩ δ|yᵢ - ŷᵢ| - ½δ²    if |yᵢ - ŷᵢ| > δ
```

Huber loss is quadratic for small errors (smooth, easy to optimize) and linear for large errors (robust to outliers). The `δ` parameter controls the threshold. This is the default loss in many gradient boosting implementations.

| Loss | Outlier sensitivity | Differentiable | Best for |
|------|---------------------|----------------|----------|
| MSE | High | Yes, everywhere | When outliers are errors, want smooth optimization |
| MAE | Low | No (at 0) | When outliers are genuine data points |
| Huber | Moderate (tunable) | Yes (everywhere) | General purpose, robust regression |

---

## 3. Ordinary Least Squares — The Closed-Form Solution

For linear regression with MSE loss, there is an **exact analytical solution** — no iterative optimization needed.

### Derivation

We want to minimize:

```
       1   n
MSE = ─── Σ (yᵢ - ŷᵢ)²
       n  i=1

In matrix form (absorbing bias into w):
MSE = (1/n) ||y - Xw||²

= (1/n)(y - Xw)ᵀ(y - Xw)

Expand:
= (1/n)(yᵀy - yᵀXw - wᵀXᵀy + wᵀXᵀXw)

= (1/n)(yᵀy - 2wᵀXᵀy + wᵀXᵀXw)
  (used: yᵀXw = wᵀXᵀy since both are scalars)
```

Take the gradient with respect to w and set to zero:

```
∂MSE/∂w = (1/n)(-2Xᵀy + 2XᵀXw) = 0

Xᵀy = XᵀXw

w = (XᵀX)⁻¹ Xᵀy   ← THE NORMAL EQUATION
```

This is the **Normal Equation** (also called the OLS solution). It gives the globally optimal weights in a single matrix operation.

### Implementation

```python
import numpy as np

def ols_solution(X, y):
    """
    Ordinary Least Squares (closed-form solution)
    w = (X^T X)^(-1) X^T y
    """
    # Add bias column (column of 1s) to X
    ones = np.ones((X.shape[0], 1))
    X_b = np.hstack([ones, X])   # shape: (n, p+1)

    # Normal equation
    w = np.linalg.inv(X_b.T @ X_b) @ X_b.T @ y
    return w

# Example
X = np.array([[1400], [1800], [2100], [2400], [950]])
y = np.array([245000, 312000, 380000, 435000, 178000])

w = ols_solution(X, y)
print(f"Bias (intercept): ${w[0]:,.2f}")
print(f"Weight (per sq ft): ${w[1]:,.2f}")

# Predict on new house
new_house = np.array([[1650]])
ones = np.ones((1, 1))
new_house_b = np.hstack([ones, new_house])
predicted_price = new_house_b @ w
print(f"\nPredicted price for 1650 sq ft: ${predicted_price[0]:,.2f}")
```

### When OLS Fails

The Normal Equation requires computing `(XᵀX)⁻¹`. This fails when:

1. **Singular matrix:** `XᵀX` is not invertible. This happens when:
   - More features than samples (p > n) — underdetermined system
   - Two features are perfectly correlated (multicollinearity)

2. **Computational cost:** For p features, the matrix inversion costs O(p³). For p = 100,000 features (text data), this is infeasible.

3. **Memory:** Storing an n×p matrix when n is millions of rows.

When OLS fails, we use **gradient descent** instead.

---

## 4. Gradient Descent for Linear Regression

Gradient descent is the workhorse optimization algorithm for nearly all of machine learning, from linear regression to deep neural networks.

### The Core Idea

```
GRADIENT DESCENT — THE HILL METAPHOR
──────────────────────────────────────────────────────────
Imagine you are on a hilly landscape (the loss surface).
Your goal: find the lowest valley (minimum loss).

Strategy: at each step, look around, find which direction
goes most steeply downhill, take a small step in that direction.
Repeat until you reach the bottom.

"Gradient" = direction of steepest INCREASE in loss
We move in the NEGATIVE gradient direction to go downhill.

w_new = w_old - α × ∂Loss/∂w

Where α (alpha) is the learning rate: how big a step to take.
```

### Deriving the Gradients for Linear Regression

Loss function (MSE for one sample):

```
L = (y - ŷ)²  where  ŷ = wx + b

∂L/∂w = 2(y - ŷ)(-x)   [chain rule: ∂L/∂ŷ × ∂ŷ/∂w]
       = -2(y - ŷ)x
       = -2(error)(x)

∂L/∂b = 2(y - ŷ)(-1)
       = -2(y - ŷ)
       = -2(error)
```

For the full dataset (averaging over all n samples):

```
∂MSE/∂w = -(2/n) Σᵢ (yᵢ - ŷᵢ)xᵢ

∂MSE/∂b = -(2/n) Σᵢ (yᵢ - ŷᵢ)

Update rules:
  w ← w - α × (∂MSE/∂w)
  b ← b - α × (∂MSE/∂b)
```

### The Three Flavors of Gradient Descent

```
BATCH vs STOCHASTIC vs MINI-BATCH
──────────────────────────────────────────────────────────
Batch Gradient Descent (BGD):
  - Compute gradient using ALL n training samples
  - Very accurate gradient estimate
  - Very slow for large datasets
  - Guaranteed convergence to minimum (for convex problems)

Stochastic Gradient Descent (SGD):
  - Compute gradient using 1 random sample per update
  - Very fast but very noisy
  - Can escape local minima (noise is sometimes helpful)
  - Convergence is erratic but eventually reaches near-minimum

Mini-Batch Gradient Descent:
  - Compute gradient using a random batch of B samples (B=32 to 256)
  - Best of both worlds: fast AND relatively stable
  - Standard in deep learning (this is what "batch_size" means)
  - Parallelizable on GPUs
```

---

## 5. Linear Regression from Scratch

Let's implement gradient descent linear regression fully in NumPy — no sklearn. Understanding this code deeply is more valuable than knowing a dozen sklearn APIs.

```python
import numpy as np
import matplotlib.pyplot as plt

class LinearRegressionGD:
    """
    Linear Regression using Gradient Descent.
    Implemented from scratch to understand the internals.
    """

    def __init__(self, learning_rate=0.01, n_iterations=1000):
        self.lr = learning_rate
        self.n_iterations = n_iterations
        self.weights = None    # w vector
        self.bias = None       # b scalar
        self.loss_history = [] # track training progress

    def fit(self, X, y):
        """
        Train the model using batch gradient descent.
        X: (n_samples, n_features)
        y: (n_samples,)
        """
        n_samples, n_features = X.shape

        # Initialize weights to zero (or small random values)
        self.weights = np.zeros(n_features)
        self.bias = 0.0

        for iteration in range(self.n_iterations):
            # Forward pass: compute predictions
            # ŷ = Xw + b
            y_pred = X @ self.weights + self.bias   # shape: (n_samples,)

            # Compute residuals (errors)
            residuals = y_pred - y   # shape: (n_samples,)

            # Compute MSE loss (for monitoring)
            mse = np.mean(residuals ** 2)
            self.loss_history.append(mse)

            # Compute gradients
            # ∂MSE/∂w = (2/n) X^T (ŷ - y)
            # ∂MSE/∂b = (2/n) Σ(ŷ - y)
            dw = (2 / n_samples) * (X.T @ residuals)
            db = (2 / n_samples) * np.sum(residuals)

            # Update parameters (gradient descent step)
            self.weights -= self.lr * dw
            self.bias    -= self.lr * db

            # Print progress every 100 iterations
            if (iteration + 1) % 100 == 0:
                print(f"  Iter {iteration+1:4d}/{self.n_iterations}, "
                      f"MSE: {mse:.4f}")

        return self

    def predict(self, X):
        """Make predictions: ŷ = Xw + b"""
        return X @ self.weights + self.bias

    def r2_score(self, X, y):
        """R² score (coefficient of determination)"""
        y_pred = self.predict(X)
        ss_res = np.sum((y - y_pred) ** 2)   # Residual sum of squares
        ss_tot = np.sum((y - y.mean()) ** 2) # Total sum of squares
        return 1 - (ss_res / ss_tot)


# ============================================================
# Test with synthetic data
# ============================================================
np.random.seed(42)

# Generate data: y = 3x₁ + 2x₂ - 5 + noise
n_samples = 200
X = np.random.randn(n_samples, 2)
true_weights = np.array([3.0, 2.0])
true_bias = -5.0
noise = np.random.randn(n_samples) * 0.5

y = X @ true_weights + true_bias + noise

# Split into train/test
split = int(0.8 * n_samples)
X_train, X_test = X[:split], X[split:]
y_train, y_test = y[:split], y[split:]

# Feature scaling (important for gradient descent!)
mean = X_train.mean(axis=0)
std = X_train.std(axis=0)
X_train_scaled = (X_train - mean) / std
X_test_scaled  = (X_test  - mean) / std

# Train
print("Training Linear Regression with Gradient Descent:")
model = LinearRegressionGD(learning_rate=0.1, n_iterations=500)
model.fit(X_train_scaled, y_train)

# Evaluate
train_r2 = model.r2_score(X_train_scaled, y_train)
test_r2  = model.r2_score(X_test_scaled, y_test)
print(f"\nTrain R²: {train_r2:.4f}")
print(f"Test  R²: {test_r2:.4f}")
print(f"\nLearned weights: {model.weights}")   # Should be near [3, 2] after scaling
print(f"Learned bias:    {model.bias:.4f}")    # Should be near -5

# ============================================================
# Now compare with sklearn
# ============================================================
from sklearn.linear_model import LinearRegression as SklearnLR

sk_model = SklearnLR()
sk_model.fit(X_train_scaled, y_train)
sk_r2 = sk_model.score(X_test_scaled, y_test)
print(f"\nSklearn R²: {sk_r2:.4f}  (our model: {test_r2:.4f})")
# Should be nearly identical
```

### Using Sklearn for Linear Regression

```python
from sklearn.linear_model import LinearRegression, SGDRegressor
from sklearn.datasets import load_boston  # deprecated → use California Housing
from sklearn.datasets import fetch_california_housing
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import mean_squared_error, r2_score
import numpy as np

# Load real dataset
data = fetch_california_housing()
X, y = data.data, data.target

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

scaler = StandardScaler()
X_train_s = scaler.fit_transform(X_train)
X_test_s  = scaler.transform(X_test)

# Method 1: OLS (closed-form, good for medium-size datasets)
ols = LinearRegression()
ols.fit(X_train_s, y_train)
y_pred_ols = ols.predict(X_test_s)

print("=== LinearRegression (OLS) ===")
print(f"Coefficients: {ols.coef_.round(3)}")
print(f"Intercept: {ols.intercept_:.3f}")
print(f"MSE:  {mean_squared_error(y_test, y_pred_ols):.4f}")
print(f"RMSE: {mean_squared_error(y_test, y_pred_ols, squared=False):.4f}")
print(f"R²:   {r2_score(y_test, y_pred_ols):.4f}")

# Method 2: SGDRegressor (stochastic gradient descent, scales to large data)
sgd = SGDRegressor(
    loss='squared_error',      # MSE loss
    learning_rate='invscaling', # Decreasing LR schedule
    eta0=0.01,                 # Initial learning rate
    max_iter=1000,
    random_state=42
)
sgd.fit(X_train_s, y_train)
y_pred_sgd = sgd.predict(X_test_s)
print("\n=== SGDRegressor ===")
print(f"R²: {r2_score(y_test, y_pred_sgd):.4f}")
```

---

## 6. Assumptions and Diagnostics

Linear regression is a statistical model with formal assumptions. Violating them doesn't make the model "invalid" for ML prediction, but understanding them helps diagnose problems.

### The Four Assumptions

```
ASSUMPTIONS OF LINEAR REGRESSION
──────────────────────────────────────────────────────────
1. LINEARITY
   The relationship between features and target is linear.
   Violation: S-curved or cyclic patterns
   Diagnose: plot predicted vs actual, residuals vs fitted
   Fix: feature engineering (log, polynomial), different model

2. INDEPENDENCE
   Each training sample is independent of others.
   Violation: time series data (sequential observations are correlated)
   Fix: time series models (ARIMA, LSTM)

3. HOMOSCEDASTICITY (constant variance)
   The variance of errors is the same for all predicted values.
   Violation: residuals fan out (errors increase with prediction size)
   Diagnose: residual plot (should show uniform scatter)
   Fix: log transform of y, WLS (Weighted Least Squares)

4. NORMALITY OF RESIDUALS
   Residuals are normally distributed.
   Violation: heavy-tailed or skewed residuals
   Diagnose: Q-Q plot, histogram of residuals
   Note: Less critical for prediction; more critical for confidence intervals
```

### Diagnostic Plots

```python
import numpy as np
import matplotlib.pyplot as plt
from sklearn.linear_model import LinearRegression
from scipy import stats

def plot_regression_diagnostics(model, X_test, y_test):
    """Create 4 standard regression diagnostic plots."""
    y_pred = model.predict(X_test)
    residuals = y_test - y_pred
    standardized_residuals = (residuals - residuals.mean()) / residuals.std()

    fig, axes = plt.subplots(2, 2, figsize=(12, 10))

    # Plot 1: Residuals vs Fitted
    axes[0, 0].scatter(y_pred, residuals, alpha=0.5, s=10)
    axes[0, 0].axhline(y=0, color='red', linestyle='--')
    axes[0, 0].set_xlabel('Fitted Values')
    axes[0, 0].set_ylabel('Residuals')
    axes[0, 0].set_title('Residuals vs Fitted')
    # Want: random scatter around 0 (no pattern)
    # Problem pattern: funnel shape = heteroscedasticity
    # Problem pattern: curve = non-linearity

    # Plot 2: Q-Q Plot (test normality of residuals)
    stats.probplot(residuals, dist="norm", plot=axes[0, 1])
    axes[0, 1].set_title('Q-Q Plot (Residual Normality)')
    # Want: points close to the diagonal line
    # Deviations at tails = heavy-tailed distribution

    # Plot 3: Scale-Location (test homoscedasticity)
    axes[1, 0].scatter(y_pred, np.sqrt(np.abs(standardized_residuals)), alpha=0.5, s=10)
    axes[1, 0].set_xlabel('Fitted Values')
    axes[1, 0].set_ylabel('√|Standardized Residuals|')
    axes[1, 0].set_title('Scale-Location')
    # Want: flat, horizontal spread
    # Problem: increasing spread = heteroscedasticity

    # Plot 4: Histogram of residuals
    axes[1, 1].hist(residuals, bins=30, edgecolor='black', color='steelblue')
    axes[1, 1].set_xlabel('Residual')
    axes[1, 1].set_ylabel('Count')
    axes[1, 1].set_title('Residual Distribution')
    # Want: bell-shaped, centered at 0

    plt.tight_layout()
    plt.savefig('regression_diagnostics.png', dpi=100)
    plt.show()
```

---

## 7. Polynomial Regression

The key insight: **polynomial regression is still linear regression** — we just engineer polynomial features first.

If the data follows `y ≈ 3x² + 2x + 1`, a degree-2 polynomial model is:

```
ŷ = w₀ + w₁x + w₂x²

Create new features: x₁ = x, x₂ = x²
Then fit: ŷ = w₀ + w₁x₁ + w₂x₂  ← still linear in the weights!
```

```python
import numpy as np
from sklearn.preprocessing import PolynomialFeatures
from sklearn.linear_model import LinearRegression
from sklearn.pipeline import Pipeline
from sklearn.model_selection import train_test_split
from sklearn.metrics import r2_score

# Generate nonlinear data
np.random.seed(42)
X = np.random.uniform(-3, 3, 100).reshape(-1, 1)
y = 0.5 * X.squeeze()**3 - 2 * X.squeeze() + np.random.randn(100) * 0.5

X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Compare linear vs polynomial models
results = []
for degree in [1, 2, 3, 5, 10]:
    pipe = Pipeline([
        ('poly', PolynomialFeatures(degree=degree, include_bias=False)),
        ('model', LinearRegression())
    ])
    pipe.fit(X_train, y_train)
    train_r2 = r2_score(y_train, pipe.predict(X_train))
    test_r2  = r2_score(y_test,  pipe.predict(X_test))
    results.append((degree, train_r2, test_r2))
    print(f"Degree {degree:2d}: Train R²={train_r2:.4f}, Test R²={test_r2:.4f}")

# Expected output (rough):
# Degree  1: Train R²=0.62xx, Test R²=0.57xx  ← underfitting
# Degree  2: Train R²=0.70xx, Test R²=0.67xx  ← still underfitting
# Degree  3: Train R²=0.96xx, Test R²=0.95xx  ← GOOD FIT (true degree is 3)
# Degree  5: Train R²=0.97xx, Test R²=0.94xx  ← slight overfitting
# Degree 10: Train R²=0.99xx, Test R²=0.70xx  ← clear overfitting
```

This demonstrates the **bias-variance tradeoff** in action:
- Degree 1: high bias, can't fit the curve
- Degree 3: just right (happens to match the true data-generating process)
- Degree 10: low training error, high test error — memorized noise

---

## 8. Regularization: Ridge, Lasso, Elastic Net

Regularization is the most important technique for preventing overfitting in linear models. It works by adding a penalty to the loss function that discourages large weights.

### Why Large Weights Cause Overfitting

When a model overfits, it typically has very large weights. Large weights mean the model is extremely sensitive to small changes in input — a hallmark of overfitting.

```
OVERFITTED POLYNOMIAL:
  ŷ = 5000x¹⁰ - 9999x⁹ + ... (huge weights)
  Tiny change in x → huge change in ŷ → brittle, not generalizable

REGULARIZED MODEL:
  Penalize large weights → prefer smaller, simpler models
  Smaller weights → smoother, more generalizable predictions
```

### L2 Regularization / Ridge Regression

```
Ridge Loss = MSE + λ × ||w||²
           = MSE + λ × Σᵢ wᵢ²

Where:
  λ (lambda) = regularization strength (hyperparameter)
  ||w||²     = L2 norm squared = sum of squared weights
```

Effect: all weights are **shrunk toward zero** but never exactly zero. The closed-form solution changes to:

```
w = (XᵀX + λI)⁻¹ Xᵀy

Adding λI to XᵀX ensures the matrix is always invertible!
→ Ridge also "fixes" the singular matrix problem
```

**Key properties:**
- Smoothly shrinks all weights
- Never produces exactly zero weights (all features kept)
- Great when you believe all features are somewhat relevant
- Works well for correlated features

### L1 Regularization / Lasso

```
Lasso Loss = MSE + λ × ||w||₁
           = MSE + λ × Σᵢ |wᵢ|

Where ||w||₁ = L1 norm = sum of absolute values of weights
```

**Key property: Lasso drives some weights to EXACTLY zero.** This performs automatic feature selection — irrelevant features get weight = 0 and are effectively removed from the model.

**Why does L1 produce sparsity?** The geometric intuition:

```
GEOMETRIC INTUITION: WHY L1 IS SPARSE
──────────────────────────────────────────────────────────
The MSE loss contours are ellipses (in weight space).
The constraint region is the shape of the regularization term.

L2 constraint: ||w||² ≤ t  → a CIRCLE (in 2D)
L1 constraint: ||w||₁ ≤ t  → a DIAMOND (in 2D)

The optimal solution is where the MSE ellipse first
touches the constraint region.

For L2 (circle): can touch anywhere on the smooth curve
For L1 (diamond): almost always touches at a CORNER
                  At corners, one weight = 0 → sparsity!

        w₂                     w₂
        │    ╭─────╮            │
        │   ╱       ╲           │     ╱╲
        │  │    ●    │          │    ╱  ╲
        │   ╲       ╱           │   ╱    ╲
        │    ╰─────╯            │  ╱  ●   ╲
        └──────────── w₁        │  ╲      ╱
        L2: smooth contact      │   ╲    ╱
        w₁ ≠ 0, w₂ ≠ 0          │    ╲  ╱
                                │     ╲╱←corners
                                └──────────── w₁
                                L1: contact at corner
                                w₁ = 0 or w₂ = 0
```

### Elastic Net — Best of Both

```
Elastic Net Loss = MSE + λ₁||w||₁ + λ₂||w||²
                 = MSE + α × [ρ||w||₁ + (1-ρ)||w||²/2]

Where:
  α = overall regularization strength
  ρ (rho) = l1_ratio in sklearn: 0 = pure Ridge, 1 = pure Lasso
```

Elastic Net handles the case where Lasso is too aggressive (drops useful correlated features).

### Implementation

```python
import numpy as np
from sklearn.linear_model import Ridge, Lasso, ElasticNet, RidgeCV, LassoCV
from sklearn.model_selection import train_test_split
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import r2_score
from sklearn.datasets import fetch_california_housing

X, y = fetch_california_housing(return_X_y=True)
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

scaler = StandardScaler()
X_train_s = scaler.fit_transform(X_train)
X_test_s  = scaler.transform(X_test)

# ---- Ridge ----
# alpha in sklearn = λ in math
# Higher alpha = stronger regularization = smaller weights
ridge = Ridge(alpha=1.0)
ridge.fit(X_train_s, y_train)
print(f"Ridge R²: {r2_score(y_test, ridge.predict(X_test_s)):.4f}")
print(f"Weights (non-zero count): {np.sum(ridge.coef_ != 0)}")  # All non-zero

# ---- Lasso ----
lasso = Lasso(alpha=0.01)  # smaller alpha = less regularization for Lasso
lasso.fit(X_train_s, y_train)
print(f"\nLasso R²: {r2_score(y_test, lasso.predict(X_test_s)):.4f}")
print(f"Weights: {lasso.coef_.round(3)}")
print(f"Non-zero weights: {np.sum(lasso.coef_ != 0)}")  # Some may be zero!

# ---- Choosing λ with Cross-Validation ----
# RidgeCV tries many alpha values automatically
alphas = np.logspace(-4, 4, 100)  # 100 values from 0.0001 to 10000
ridge_cv = RidgeCV(alphas=alphas, cv=5)  # 5-fold cross validation
ridge_cv.fit(X_train_s, y_train)
print(f"\nBest alpha (RidgeCV): {ridge_cv.alpha_:.4f}")
print(f"RidgeCV R²: {r2_score(y_test, ridge_cv.predict(X_test_s)):.4f}")

lasso_cv = LassoCV(alphas=alphas, cv=5, max_iter=5000)
lasso_cv.fit(X_train_s, y_train)
print(f"\nBest alpha (LassoCV): {lasso_cv.alpha_:.6f}")
print(f"LassoCV R²: {r2_score(y_test, lasso_cv.predict(X_test_s)):.4f}")

# ---- Elastic Net ----
from sklearn.linear_model import ElasticNetCV
enet_cv = ElasticNetCV(
    l1_ratio=[0.1, 0.5, 0.7, 0.9, 0.95, 0.99, 1.0],  # Try different L1/L2 mixes
    cv=5,
    max_iter=5000
)
enet_cv.fit(X_train_s, y_train)
print(f"\nElastic Net: alpha={enet_cv.alpha_:.4f}, l1_ratio={enet_cv.l1_ratio_:.2f}")
print(f"Elastic Net R²: {r2_score(y_test, enet_cv.predict(X_test_s)):.4f}")
```

---

## 9. Evaluating Regression Models

```
REGRESSION EVALUATION METRICS
──────────────────────────────────────────────────────────────
       1   n                           sensitive to outliers
MAE = ─── Σ|yᵢ-ŷᵢ|                    interpretable (same units as y)
       n

       1   n
MSE = ─── Σ(yᵢ-ŷᵢ)²                  penalizes outliers more
       n

RMSE = √MSE                           same units as y, penalizes outliers

MAPE = (100/n) Σ |(yᵢ-ŷᵢ)/yᵢ|        percentage error, not defined at y=0

       Σ(yᵢ-ŷᵢ)²             SS_res
R² = 1 - ─────────── = 1 - ──────
         Σ(yᵢ-ȳ)²              SS_tot

  R² = 1.0 → perfect fit
  R² = 0.0 → no better than predicting the mean
  R² < 0.0 → worse than predicting the mean (bad model)
```

**Adjusted R²** penalizes for adding useless features:

```
           (n-1)(1 - R²)
R²_adj = 1 - ────────────────
              n - p - 1

Where p = number of features

R²_adj increases only if new feature genuinely improves the fit.
Use this when comparing models with different numbers of features.
```

```python
from sklearn.metrics import mean_absolute_error, mean_squared_error, r2_score
import numpy as np

def regression_report(y_true, y_pred, n_features=None, label=""):
    """Print comprehensive regression metrics."""
    mae   = mean_absolute_error(y_true, y_pred)
    mse   = mean_squared_error(y_true, y_pred)
    rmse  = np.sqrt(mse)
    r2    = r2_score(y_true, y_pred)
    n     = len(y_true)

    # MAPE (exclude zero targets)
    nonzero = y_true != 0
    mape = np.mean(np.abs((y_true[nonzero] - y_pred[nonzero]) / y_true[nonzero])) * 100

    # Adjusted R²
    if n_features is not None:
        adj_r2 = 1 - (1 - r2) * (n - 1) / (n - n_features - 1)
    else:
        adj_r2 = None

    print(f"\n{'='*40}")
    print(f"Regression Metrics: {label}")
    print(f"{'='*40}")
    print(f"  MAE:   {mae:.4f}")
    print(f"  MSE:   {mse:.4f}")
    print(f"  RMSE:  {rmse:.4f}")
    print(f"  MAPE:  {mape:.2f}%")
    print(f"  R²:    {r2:.4f}")
    if adj_r2 is not None:
        print(f"  Adj R²:{adj_r2:.4f}")
```

---

## 10. Logistic Regression — Classification with Probabilities

Despite its name, logistic regression is a **classifier**, not a regressor. It predicts the **probability** that an input belongs to a class.

### Why Not Just Use Linear Regression for Classification?

Suppose we try linear regression on a binary classification problem (y = 0 or 1):

```
PROBLEM WITH LINEAR REGRESSION FOR CLASSIFICATION
──────────────────────────────────────────────────────────────
Linear regression output: ŷ = wx + b  can be ANY real number

But probabilities must be between 0 and 1!
  ŷ = 1.7 → "170% probability of spam?" — makes no sense
  ŷ = -0.3 → "negative probability?" — impossible

Also: decision boundary is problematic.
If we set threshold at 0.5:
  Adding a new data point far from the boundary shifts the line
  and can misclassify previously-correct points.

We need a function that:
  1. Takes any real number as input
  2. Outputs a value strictly between 0 and 1
  3. Is smooth and differentiable
  4. Has a nice derivative (for gradient descent)
```

### The Sigmoid Function

The sigmoid function is the answer:

```
        1
σ(z) = ────────    where z = w^T x + b
        1 + e^(-z)

Properties:
  σ(z) → 0 as z → -∞
  σ(z) → 1 as z → +∞
  σ(0) = 0.5 exactly
  σ'(z) = σ(z)(1 - σ(z))   ← nice derivative!

Plot:
  1.0 ─────────────────────────╮────
  0.9                      ╭───╯
  0.8                  ╭───╯
  0.7              ╭───╯
  0.6          ╭───╯
  0.5 ─────────┼──────────────────── (z=0, σ=0.5)
  0.4      ───╯
  0.3   ───╯
  0.2 ──╯
  0.1 ╯
  0.0 ──────╯
      ────────────────────────────
      -6  -4  -2  0  2  4  6   z

  Saturates near 0 and 1 ("squashes" any number to (0,1))
  Gradient ≈ 0 at extremes (vanishing gradient issue for deep nets)
```

### The Logistic Regression Model

```
Linear combination: z = w^T x + b

Apply sigmoid: P(y=1|x) = σ(z) = 1/(1 + e^(-z))

If P(y=1|x) ≥ 0.5 → predict class 1
If P(y=1|x) < 0.5 → predict class 0

The decision boundary is where P = 0.5, i.e., where z = 0:
  w^T x + b = 0
  This is a HYPERPLANE in the feature space.
  In 2D with 2 features: a line.
```

```
DECISION BOUNDARY IN 2D
──────────────────────────────────────────────────────────────
x₂
│
│   ●   ●                     ○   ○
│     ●  ●                  ○   ○
│       ●  ●               ○   ○
│         ● \             ○ ○
│            \ Decision  ○
│             \ boundary/
│              \/       /
│               ●      ○
│                     ○
└─────────────────────────────── x₁

● = class 0, ○ = class 1
The line w₁x₁ + w₂x₂ + b = 0 separates the classes.
```

---

## 11. Binary Cross-Entropy Loss

We cannot use MSE for logistic regression — the combination of sigmoid + MSE creates a non-convex loss surface with many local minima. We need **Binary Cross-Entropy** (also called Log Loss).

### Derivation from Maximum Likelihood Estimation

The model outputs a probability. We want to find weights that make the **observed data as probable as possible**. This is Maximum Likelihood Estimation (MLE).

For binary classification:
```
P(y=1|x) = σ(wx+b) = p̂

Likelihood for one sample:
  L(w) = p̂^y × (1-p̂)^(1-y)

  When y=1: L = p̂     (want p̂ close to 1)
  When y=0: L = (1-p̂) (want p̂ close to 0)

For all n samples (assuming independence):
  L(w) = Π p̂ᵢ^yᵢ × (1-p̂ᵢ)^(1-yᵢ)

Log-likelihood (easier to work with, equivalent):
  ℓ(w) = Σ [yᵢ log(p̂ᵢ) + (1-yᵢ) log(1-p̂ᵢ)]

MAXIMIZE log-likelihood = MINIMIZE negative log-likelihood:
  BCE = -(1/n) Σ [yᵢ log(p̂ᵢ) + (1-yᵢ) log(1-p̂ᵢ)]
```

This is **Binary Cross-Entropy**. The "cross-entropy" name comes from information theory — it measures the "surprise" of the model's predictions given the true labels.

### Gradient of BCE Loss

```
∂BCE/∂w = (1/n) Σ (p̂ᵢ - yᵢ) xᵢ
∂BCE/∂b = (1/n) Σ (p̂ᵢ - yᵢ)

This is the SAME form as linear regression's MSE gradient!
The elegance of cross-entropy with sigmoid:
  residual = (p̂ - y) which is intuitive
```

Update rules (identical in form to linear regression!):

```python
# Logistic Regression gradient descent - the core update
residuals = y_pred_proba - y_true     # ŷ - y
dw = (1/n) * X.T @ residuals
db = (1/n) * np.sum(residuals)
w -= lr * dw
b -= lr * db
```

---

## 12. Multi-class Logistic Regression: Softmax

For K > 2 classes, we use the **Softmax** function, which generalizes sigmoid:

```
                   e^(z_k)
Softmax(z)_k = ──────────────
               Σⱼ e^(z_j)

Where z_k = w_k^T x + b_k  (a separate linear model per class)

Properties:
  All outputs positive (e^x > 0 always)
  All outputs sum to 1 → valid probability distribution
  Largest z_k gets highest probability
```

The loss is **Categorical Cross-Entropy**:

```
CE = -(1/n) Σᵢ Σₖ yᵢₖ log(p̂ᵢₖ)

Where yᵢₖ = 1 if sample i belongs to class k, else 0 (one-hot encoding)
```

---

## 13. Classification Evaluation Metrics

This section covers the metrics you will use every day in real ML work.

### The Confusion Matrix

For binary classification, every prediction falls into one of four categories:

```
CONFUSION MATRIX
──────────────────────────────────────────────────────────
                  Predicted Negative    Predicted Positive
                  (model says No)       (model says Yes)
                ┌─────────────────────┬────────────────────┐
Actual Negative │  True Negative (TN) │ False Positive (FP)│
  (actually No) │  "correctly said No"│ "wrongly said Yes" │
                │                     │  → Type I Error    │
                ├─────────────────────┼────────────────────┤
Actual Positive │ False Negative (FN) │ True Positive (TP) │
 (actually Yes) │ "wrongly said No"   │ "correctly said Yes│
                │  → Type II Error    │                    │
                └─────────────────────┴────────────────────┘

Medical example (cancer detection):
  TP: Patient has cancer, model says cancer        ← best case
  TN: No cancer, model says no cancer              ← correct
  FP: No cancer, but model says cancer             ← unnecessary stress, more tests
  FN: Has cancer, but model says no cancer         ← DANGEROUS: missed diagnosis
```

### Core Metrics

```
ACCURACY = (TP + TN) / (TP + TN + FP + FN)
  "What fraction of all predictions were correct?"
  Problem: misleading when classes are imbalanced!
  Example: 99% healthy, 1% sick → model predicting "healthy" always = 99% accuracy

PRECISION = TP / (TP + FP)
  "Of everything I predicted as positive, how many were actually positive?"
  "When I say positive, how often am I right?"
  High precision → few false alarms
  Use when FP is costly: spam filter (legitimate emails flagged as spam)

RECALL = TP / (TP + FN)   (also called Sensitivity, True Positive Rate)
  "Of all the actual positives, how many did I catch?"
  High recall → few misses
  Use when FN is costly: disease detection (missing a sick patient)

F1 SCORE = 2 × (Precision × Recall) / (Precision + Recall)
  Harmonic mean of precision and recall.
  Punishes extreme imbalance between precision and recall.
  F1 = 1.0 only when both precision AND recall are 1.0

SPECIFICITY = TN / (TN + FP)
  "Of all actual negatives, how many did I correctly identify?"
  = 1 - False Positive Rate

TRUE POSITIVE RATE (TPR)  = Recall = TP / (TP + FN)
FALSE POSITIVE RATE (FPR) = FP / (FP + TN) = 1 - Specificity
```

### The Precision-Recall Tradeoff

```
THRESHOLD AND THE TRADEOFF
──────────────────────────────────────────────────────────────
Default: predict "positive" if P(y=1|x) ≥ 0.5

If you LOWER the threshold (e.g., to 0.3):
  → More samples predicted positive
  → Recall INCREASES (catch more true positives)
  → Precision DECREASES (more false positives)

If you RAISE the threshold (e.g., to 0.7):
  → Fewer samples predicted positive
  → Precision INCREASES (only confident positives predicted)
  → Recall DECREASES (miss more true positives)

You CANNOT simultaneously maximize both precision and recall.
The right balance depends on your problem.

Spam filter:  False positive (flagging good email) is bad → high precision
Cancer screen: False negative (missing cancer) is bad → high recall
```

### ROC Curve and AUC

The ROC (Receiver Operating Characteristic) curve plots TPR vs FPR at every possible threshold from 0 to 1.

```
ROC CURVE
──────────────────────────────────────────────────────────────
TPR  │                                       ╭────────────
(    │                               ╭───────╯
Recall)                       ╭──────╯
1.0  │                   ╭────╯
     │              ╭────╯
     │         ╭────╯
     │    ╭────╯      ← AUC = area under this curve
     │ ╭──╯            AUC = 1.0: perfect model
0.5  ├─╯               AUC = 0.5: random guessing (diagonal line)
     │╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱╱ diagonal = random classifier
0.0  └─────────────────────────────────────────────────────
     0.0        0.5         1.0     FPR

AUC (Area Under the ROC Curve):
  AUC = 1.0 → Perfect model
  AUC = 0.9 → Excellent
  AUC = 0.8 → Good
  AUC = 0.7 → Acceptable
  AUC = 0.5 → Random guessing
  AUC < 0.5 → Worse than random (flip predictions!)
```

**What AUC means:** The probability that the model ranks a random positive sample higher than a random negative sample. A model-level metric independent of threshold.

### Precision-Recall Curve

Better than ROC curve for **imbalanced datasets**:

```
PRECISION-RECALL CURVE
──────────────────────────────────────────────────────────────
Precision
1.0  │  ╮
     │   ╲
     │    ╲
0.5  │     ╲──────╮
     │            ╲
     │              ╲──────
0.0  └─────────────────────────── Recall
     0.0         0.5         1.0

High-right corner = ideal (high precision AND high recall)
AUC-PR = area under this curve → better metric for imbalanced data

ROC AUC can be misleadingly optimistic when negatives dominate.
Use Precision-Recall AUC when:
  - Class imbalance is severe (1% positive rate)
  - False positives and false negatives have different costs
```

### Complete Metrics Implementation

```python
import numpy as np
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, f1_score,
    roc_auc_score, average_precision_score, confusion_matrix,
    classification_report, roc_curve, precision_recall_curve
)

def classification_report_full(y_true, y_pred, y_proba=None, label=""):
    """
    Comprehensive classification evaluation.
    y_proba: probability of positive class, shape (n,)
    """
    print(f"\n{'='*50}")
    print(f"Classification Report: {label}")
    print(f"{'='*50}")

    # Confusion matrix
    cm = confusion_matrix(y_true, y_pred)
    tn, fp, fn, tp = cm.ravel()
    print(f"\nConfusion Matrix:")
    print(f"  True Neg (TN): {tn:4d}  |  False Pos (FP): {fp:4d}")
    print(f"  False Neg (FN):{fn:4d}  |  True Pos  (TP): {tp:4d}")

    # Core metrics
    acc  = accuracy_score(y_true, y_pred)
    prec = precision_score(y_true, y_pred)
    rec  = recall_score(y_true, y_pred)
    f1   = f1_score(y_true, y_pred)
    spec = tn / (tn + fp) if (tn + fp) > 0 else 0  # specificity

    print(f"\nCore Metrics:")
    print(f"  Accuracy:    {acc:.4f}")
    print(f"  Precision:   {prec:.4f}  (of predicted positives, how many correct?)")
    print(f"  Recall:      {rec:.4f}  (of actual positives, how many caught?)")
    print(f"  Specificity: {spec:.4f}  (of actual negatives, how many correctly identified?)")
    print(f"  F1 Score:    {f1:.4f}  (harmonic mean of precision & recall)")

    if y_proba is not None:
        auc_roc = roc_auc_score(y_true, y_proba)
        auc_pr  = average_precision_score(y_true, y_proba)
        print(f"\nProbability-based Metrics:")
        print(f"  AUC-ROC: {auc_roc:.4f}  (ranking quality, threshold-independent)")
        print(f"  AUC-PR:  {auc_pr:.4f}  (better for imbalanced data)")
```

---

## 14. Full Example: Credit Card Default Prediction

A complete end-to-end example tying together everything in this chapter.

```python
# =============================================================================
# FULL EXAMPLE: Credit Card Default Prediction
# =============================================================================
# Dataset: UCI Credit Card Default (or synthetic version below)
# Problem: Predict if a customer will default on credit card payment next month
# This is a binary classification problem with class imbalance

import numpy as np
import pandas as pd
from sklearn.datasets import make_classification
from sklearn.model_selection import train_test_split, cross_val_score, StratifiedKFold
from sklearn.preprocessing import StandardScaler
from sklearn.linear_model import LogisticRegression
from sklearn.metrics import (
    classification_report, confusion_matrix, roc_auc_score,
    roc_curve, precision_recall_curve, average_precision_score,
    accuracy_score, f1_score
)
import warnings
warnings.filterwarnings('ignore')

# =============================================================================
# STEP 1: Generate synthetic dataset (mirrors real credit data structure)
# =============================================================================
np.random.seed(42)

X, y = make_classification(
    n_samples=10000,
    n_features=20,
    n_informative=10,    # 10 features actually matter
    n_redundant=5,       # 5 are linear combinations of informative ones
    n_clusters_per_class=1,
    weights=[0.78, 0.22],  # 78% no default, 22% default (imbalanced)
    random_state=42
)

feature_names = [f'feature_{i}' for i in range(X.shape[1])]
df = pd.DataFrame(X, columns=feature_names)
df['default'] = y

print("Dataset shape:", df.shape)
print(f"\nDefault rate: {y.mean():.1%} ({y.sum()} out of {len(y)})")
print("Class distribution:", np.bincount(y))

# =============================================================================
# STEP 2: Stratified Train/Test Split
# =============================================================================
X_train, X_test, y_train, y_test = train_test_split(
    X, y,
    test_size=0.2,
    random_state=42,
    stratify=y  # CRITICAL: preserve class proportions in imbalanced data
)

print(f"\nTrain size: {len(X_train)}, default rate: {y_train.mean():.1%}")
print(f"Test size:  {len(X_test)},  default rate: {y_test.mean():.1%}")

# =============================================================================
# STEP 3: Feature Scaling
# =============================================================================
scaler = StandardScaler()
X_train_s = scaler.fit_transform(X_train)  # Fit on train only
X_test_s  = scaler.transform(X_test)       # Transform test with train stats

# =============================================================================
# STEP 4: Model 1 — Basic Logistic Regression
# =============================================================================
print("\n=== Model 1: Basic Logistic Regression ===")

lr_basic = LogisticRegression(max_iter=1000, random_state=42)
lr_basic.fit(X_train_s, y_train)

y_pred_basic = lr_basic.predict(X_test_s)
y_proba_basic = lr_basic.predict_proba(X_test_s)[:, 1]

print(f"Accuracy: {accuracy_score(y_test, y_pred_basic):.4f}")
print(f"AUC-ROC:  {roc_auc_score(y_test, y_proba_basic):.4f}")
print(f"F1 Score: {f1_score(y_test, y_pred_basic):.4f}")
print("\nClassification Report:")
print(classification_report(y_test, y_pred_basic, target_names=['no default', 'default']))

# =============================================================================
# STEP 5: Model 2 — Class-Weighted (handles imbalance)
# =============================================================================
print("\n=== Model 2: Logistic Regression with class_weight='balanced' ===")

lr_balanced = LogisticRegression(
    class_weight='balanced',  # Upweight minority class in loss
    max_iter=1000,
    random_state=42
)
lr_balanced.fit(X_train_s, y_train)
y_pred_balanced = lr_balanced.predict(X_test_s)
y_proba_balanced = lr_balanced.predict_proba(X_test_s)[:, 1]

print(f"Accuracy: {accuracy_score(y_test, y_pred_balanced):.4f}")
print(f"AUC-ROC:  {roc_auc_score(y_test, y_proba_balanced):.4f}")
print(f"F1 Score: {f1_score(y_test, y_pred_balanced):.4f}")
print("\nClassification Report:")
print(classification_report(y_test, y_pred_balanced, target_names=['no default', 'default']))

# =============================================================================
# STEP 6: Model 3 — L2 Regularized (tune C parameter)
# =============================================================================
print("\n=== Model 3: Regularized Logistic Regression (C tuning) ===")

# C is the inverse of regularization strength:
# Small C → strong regularization → simpler model
# Large C → weak regularization → complex model
# Default C=1.0

for C in [0.01, 0.1, 1.0, 10.0, 100.0]:
    lr = LogisticRegression(C=C, max_iter=1000, random_state=42, class_weight='balanced')
    lr.fit(X_train_s, y_train)
    proba = lr.predict_proba(X_test_s)[:, 1]
    auc = roc_auc_score(y_test, proba)
    print(f"  C={C:6.2f} → AUC-ROC: {auc:.4f}")

# =============================================================================
# STEP 7: Cross-Validation
# =============================================================================
print("\n=== Cross-Validation (5-Fold Stratified) ===")

cv = StratifiedKFold(n_splits=5, shuffle=True, random_state=42)
lr_final = LogisticRegression(C=1.0, class_weight='balanced', max_iter=1000, random_state=42)

# Score on multiple metrics
from sklearn.model_selection import cross_validate

cv_results = cross_validate(
    lr_final, X_train_s, y_train,
    cv=cv,
    scoring=['roc_auc', 'f1', 'precision', 'recall'],
    return_train_score=True
)

for metric in ['roc_auc', 'f1', 'precision', 'recall']:
    val_scores = cv_results[f'test_{metric}']
    print(f"  {metric:10s}: {val_scores.mean():.4f} ± {val_scores.std():.4f}")

# =============================================================================
# STEP 8: Threshold Optimization
# =============================================================================
print("\n=== Threshold Optimization ===")
# Default threshold is 0.5, but this may not be optimal

# Train final model
lr_final.fit(X_train_s, y_train)
y_proba_final = lr_final.predict_proba(X_test_s)[:, 1]

# Find threshold that maximizes F1 score
from sklearn.metrics import f1_score

thresholds = np.arange(0.1, 0.9, 0.05)
best_threshold = 0.5
best_f1 = 0

for thresh in thresholds:
    y_pred_thresh = (y_proba_final >= thresh).astype(int)
    f1 = f1_score(y_test, y_pred_thresh)
    if f1 > best_f1:
        best_f1 = f1
        best_threshold = thresh

print(f"Default threshold (0.5) F1: {f1_score(y_test, (y_proba_final >= 0.5).astype(int)):.4f}")
print(f"Optimal threshold ({best_threshold:.2f}) F1: {best_f1:.4f}")

y_pred_optimal = (y_proba_final >= best_threshold).astype(int)
print(f"\nFinal model with optimal threshold ({best_threshold:.2f}):")
print(classification_report(y_test, y_pred_optimal, target_names=['no default', 'default']))

# =============================================================================
# STEP 9: Model Coefficients (Interpret the Model)
# =============================================================================
print("\n=== Model Interpretation ===")
coef_df = pd.DataFrame({
    'feature': feature_names,
    'coefficient': lr_final.coef_[0]
})
coef_df['abs_coef'] = coef_df['coefficient'].abs()
coef_df = coef_df.sort_values('abs_coef', ascending=False)

print("Top 10 most important features:")
for _, row in coef_df.head(10).iterrows():
    direction = "↑ increases" if row['coefficient'] > 0 else "↓ decreases"
    print(f"  {row['feature']}: {row['coefficient']:+.4f}  "
          f"({direction} default probability)")

print("\nNote: Coefficients on scaled features — larger magnitude = more important")
```

---

## 15. Summary

```
CHAPTER 09 KEY CONCEPTS
─────────────────────────────────────────────────────────────

LINEAR REGRESSION:
  Model: ŷ = Xw + b
  Loss: MSE = (1/n)Σ(y - ŷ)²
  Solution: OLS w = (X^T X)^(-1) X^T y  (or gradient descent)
  Evaluation: R², RMSE, MAE
  Regularization: Ridge (L2), Lasso (L1), Elastic Net

GRADIENT DESCENT:
  w ← w - α × ∂Loss/∂w
  Learning rate α: too high → diverge, too low → slow
  Variants: Batch, SGD, Mini-batch

REGULARIZATION:
  Ridge: adds λ||w||²  → shrinks weights, no sparsity
  Lasso: adds λ||w||₁  → sparsity (feature selection)
  C parameter in sklearn = 1/λ (INVERSE of regularization)
  Choose with cross-validation

LOGISTIC REGRESSION:
  Model: P(y=1|x) = σ(w^T x + b) where σ is sigmoid
  Loss: Binary Cross-Entropy (derived from MLE)
  NOT a regression model — it's a classifier!

CLASSIFICATION METRICS:
  Accuracy:  (TP+TN)/N  — misleading for imbalanced data
  Precision: TP/(TP+FP) — how precise when you say positive
  Recall:    TP/(TP+FN) — how many positives did you catch
  F1:        harmonic mean of precision and recall
  AUC-ROC:   threshold-independent ranking quality
  Use AUC-PR for heavily imbalanced data
```

### Decision Guide: Which Loss / Model to Use

| Situation | Use |
|-----------|-----|
| Regression, outliers are errors | MSE + Linear Regression |
| Regression, robust to outliers needed | MAE or Huber loss |
| High-dimensional regression, multicollinearity | Ridge |
| Feature selection via regression | Lasso |
| Binary classification | Logistic Regression + Cross-Entropy |
| Imbalanced classes | class_weight='balanced' + Recall/F1/AUC-PR |
| Need probability calibration | Logistic Regression (well-calibrated) |

---

## Mini Projects

### Mini Project 1: Interactive House Price Predictor (1.5 hours)

**Goal:** Build a Ridge regression model on housing data, add interactive sliders to predict price on-the-fly.

```python
import pandas as pd
import numpy as np
from sklearn.datasets import fetch_california_housing
from sklearn.linear_model import Ridge
from sklearn.preprocessing import StandardScaler
from sklearn.pipeline import make_pipeline
from sklearn.model_selection import train_test_split
from sklearn.metrics import mean_absolute_error, r2_score
import matplotlib.pyplot as plt

# Load data
housing = fetch_california_housing(as_frame=True)
df = housing.frame

print(df.head())
print(f"\nTarget: Median house value (100k$)")
print(f"Features: {list(housing.feature_names)}")

# Train model
X, y = df[housing.feature_names], df["MedHouseVal"]
X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

model = make_pipeline(StandardScaler(), Ridge(alpha=1.0))
model.fit(X_train, y_train)

y_pred = model.predict(X_test)
print(f"\nTest MAE: ${mean_absolute_error(y_test, y_pred)*100_000:.0f}")
print(f"Test R²: {r2_score(y_test, y_pred):.3f}")

# Prediction vs Actual plot
fig, axes = plt.subplots(1, 2, figsize=(14, 5))

axes[0].scatter(y_test, y_pred, alpha=0.3, s=10)
lim = [0, y_test.max()]
axes[0].plot(lim, lim, 'r--', label='Perfect prediction')
axes[0].set_xlabel("Actual Price ($100k)"), axes[0].set_ylabel("Predicted")
axes[0].set_title(f"Predictions vs Actual (R²={r2_score(y_test, y_pred):.3f})")
axes[0].legend()

# Residuals
residuals = y_test - y_pred
axes[1].hist(residuals, bins=50, edgecolor='black')
axes[1].axvline(0, color='red', linestyle='--')
axes[1].set_xlabel("Residual (Actual - Predicted)")
axes[1].set_title("Residual Distribution")

plt.tight_layout()
plt.savefig("house_price_prediction.png")
plt.show()

# Interactive prediction (manual slider via print)
feature_medians = X_train.median()

def predict_price(**feature_overrides):
    features = feature_medians.copy()
    for k, v in feature_overrides.items():
        features[k] = v
    price = model.predict([features])[0]
    print(f"Predicted price: ${price * 100_000:,.0f}")
    return price

print("\n--- Interactive Predictor ---")
predict_price(MedInc=8.0, HouseAge=10, AveRooms=7)
predict_price(MedInc=2.0, HouseAge=40, AveRooms=4)
```

---

### Mini Project 2: Regularization Path Visualizer (1 hour)

**Goal:** Show how Ridge and Lasso shrink coefficients as regularization strength increases — see feature elimination happen live.

```python
import numpy as np
import matplotlib.pyplot as plt
from sklearn.datasets import load_diabetes
from sklearn.linear_model import Ridge, Lasso
from sklearn.preprocessing import StandardScaler

# Load dataset
X, y = load_diabetes(return_X_y=True)
feature_names = load_diabetes().feature_names
X_scaled = StandardScaler().fit_transform(X)

alphas = np.logspace(-3, 4, 100)

ridge_coefs = []
lasso_coefs = []

for alpha in alphas:
    ridge = Ridge(alpha=alpha).fit(X_scaled, y)
    lasso = Lasso(alpha=alpha, max_iter=10000).fit(X_scaled, y)
    ridge_coefs.append(ridge.coef_)
    lasso_coefs.append(lasso.coef_)

ridge_coefs = np.array(ridge_coefs)
lasso_coefs = np.array(lasso_coefs)

fig, axes = plt.subplots(1, 2, figsize=(14, 6))
cmap = plt.cm.Set1(np.linspace(0, 1, len(feature_names)))

for i, name in enumerate(feature_names):
    axes[0].plot(alphas, ridge_coefs[:, i], color=cmap[i], label=name)
    axes[1].plot(alphas, lasso_coefs[:, i], color=cmap[i], label=name)

for ax, title in zip(axes, ["Ridge (L2) Path", "Lasso (L1) Path"]):
    ax.set_xscale('log')
    ax.axhline(0, color='black', linewidth=0.5)
    ax.set_xlabel("Regularization Strength (α, log scale)")
    ax.set_ylabel("Coefficient value")
    ax.set_title(title)
    ax.legend(loc='upper right', fontsize=8, ncol=2)

plt.suptitle("Regularization Paths: How Coefficients Shrink with α")
plt.tight_layout()
plt.savefig("regularization_paths.png")
plt.show()

# Count active features (non-zero) in Lasso
print("\nLasso feature sparsity:")
for alpha_val, coef in zip([0.01, 0.1, 1.0, 10.0], 
                            [Lasso(alpha=a).fit(X_scaled, y).coef_ for a in [0.01, 0.1, 1.0, 10.0]]):
    active = (coef != 0).sum()
    print(f"  α={alpha_val:.2f}: {active}/{len(feature_names)} features active")
```

---

### Mini Project 3: Polynomial Degree Selector with Cross-Validation (45 min)

**Goal:** Automatically find the best polynomial degree for your regression problem using cross-validation.

```python
import numpy as np
import matplotlib.pyplot as plt
from sklearn.preprocessing import PolynomialFeatures
from sklearn.linear_model import LinearRegression, Ridge
from sklearn.pipeline import make_pipeline
from sklearn.model_selection import cross_val_score

np.random.seed(42)
n = 80
x = np.linspace(-3, 3, n)
y = 0.5*x**3 - 2*x + np.random.randn(n) * 1.5  # True: cubic

X = x.reshape(-1, 1)
max_degree = 12
degrees = range(1, max_degree + 1)

train_scores, cv_scores = [], []

for d in degrees:
    model = make_pipeline(PolynomialFeatures(d), Ridge(alpha=0.1))
    cv = cross_val_score(model, X, y, cv=5, scoring='neg_mean_squared_error')
    model.fit(X, y)
    train_mse = -cross_val_score(model, X, y, cv=5, scoring='neg_mean_squared_error').mean()
    cv_mse = -cv.mean()
    train_scores.append(train_mse)
    cv_scores.append(cv_mse)

best_degree = degrees[np.argmin(cv_scores)]

fig, axes = plt.subplots(1, 2, figsize=(14, 5))

axes[0].plot(degrees, train_scores, 'b-o', label="Training MSE")
axes[0].plot(degrees, cv_scores, 'r-o', label="CV MSE (5-fold)")
axes[0].axvline(best_degree, color='green', linestyle='--',
                label=f"Best degree={best_degree}")
axes[0].set_xlabel("Polynomial Degree"), axes[0].set_ylabel("MSE")
axes[0].set_title("Finding Optimal Polynomial Degree")
axes[0].legend()

x_plot = np.linspace(-3, 3, 200).reshape(-1, 1)
best_model = make_pipeline(PolynomialFeatures(best_degree), Ridge(alpha=0.1)).fit(X, y)
overfit_model = make_pipeline(PolynomialFeatures(12), LinearRegression()).fit(X, y)

axes[1].scatter(x, y, alpha=0.5, label="Data")
axes[1].plot(x_plot.flatten(), best_model.predict(x_plot), 'g-', lw=2.5,
             label=f"Degree {best_degree} (best)")
axes[1].plot(x_plot.flatten(), overfit_model.predict(x_plot), 'r--', lw=1.5,
             label="Degree 12 (overfit)")
axes[1].set_ylim(-20, 20)
axes[1].set_title("Best vs Overfit Model")
axes[1].legend()

plt.tight_layout()
plt.savefig("polynomial_degree_selection.png")
plt.show()
print(f"Best polynomial degree: {best_degree}")
```

---

## 16. Exercises

**Exercise 1:** Implement logistic regression from scratch using gradient descent. The skeleton:
- `sigmoid(z)` function
- `fit(X, y)` with binary cross-entropy loss
- `predict_proba(X)` returns probabilities
- `predict(X)` returns 0/1 predictions with threshold=0.5
- Verify your implementation gives the same AUC-ROC as sklearn's LogisticRegression on the Iris dataset (setosa vs non-setosa).

**Exercise 2:** The C parameter in sklearn's LogisticRegression is the **inverse** of regularization strength. On the breast cancer dataset:
- Train models with C = [0.001, 0.01, 0.1, 1, 10, 100]
- For each C, record: number of non-zero coefficients, training accuracy, test accuracy
- Plot all three as a function of log(C)
- Answer: what happens to model complexity as C increases?

**Exercise 3:** Class imbalance experiment. Use `make_classification(weights=[0.95, 0.05])` to create a 95/5 split. Compare:
- Model A: standard logistic regression, accuracy metric
- Model B: standard logistic regression, AUC-ROC metric
- Model C: class_weight='balanced', AUC-ROC metric
Which model would you deploy for a fraud detection system? Justify your answer.

**Exercise 4:** Polynomial regression and overfitting. Generate data from `y = sin(x) + noise`. Fit polynomial features of degree 1 through 20. For each:
- Compute training MSE and test MSE
- Plot both on the same graph
- Identify the optimal degree
- Use Ridge regression on degree-15 features — does regularization recover the good performance?

**Exercise 5:** Derive the gradient of Binary Cross-Entropy with respect to the weights `w`. Start from `BCE = -(1/n) Σ [y log σ(wx+b) + (1-y) log(1 - σ(wx+b))]`. Use the chain rule and the sigmoid derivative `σ'(z) = σ(z)(1-σ(z))`. Show that the result simplifies to `(1/n) X^T (σ(Xw+b) - y)`.

---

**Next Chapter →** [Chapter 10: Decision Trees, Random Forests, and Gradient Boosting](./10-decision-trees-random-forests-boosting.md)

*Linear models assume the world is linear — a big assumption. Tree-based methods make no such assumption and often win on tabular data. Let's understand why.*
