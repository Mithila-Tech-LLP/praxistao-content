# Chapter 07: Probability and Statistics for Machine Learning

> **"Uncertainty is not ignorance — it is the native language of intelligence. Every model output is a claim about probability, and every training objective is a statement about probability theory."**

---

## Table of Contents
1. [Why Probability in ML?](#1-why-probability-in-ml)
2. [Probability Basics](#2-probability-basics)
3. [Conditional Probability](#3-conditional-probability)
4. [Bayes' Theorem](#4-bayes-theorem)
5. [Independence](#5-independence)
6. [Random Variables](#6-random-variables)
7. [Expectation and Variance](#7-expectation-and-variance)
8. [Covariance and Correlation](#8-covariance-and-correlation)
9. [Key Probability Distributions](#9-key-probability-distributions)
10. [The Normal Distribution](#10-the-normal-distribution)
11. [The Multivariate Normal](#11-the-multivariate-normal)
12. [Central Limit Theorem](#12-central-limit-theorem)
13. [Maximum Likelihood Estimation (MLE)](#13-maximum-likelihood-estimation-mle)
14. [MAP Estimation and Regularization](#14-map-estimation-and-regularization)
15. [Information Theory for ML](#15-information-theory-for-ml)
16. [Hypothesis Testing in ML](#16-hypothesis-testing-in-ml)
17. [Confidence Intervals](#17-confidence-intervals)
18. [Correlation vs Causation](#18-correlation-vs-causation)
19. [Summary](#19-summary)
20. [Exercises](#20-exercises)

---

## 1. Why Probability in ML?

Every single ML model makes probabilistic claims, whether it admits it or not.

```
What models actually do:
─────────────────────────────────────────────────────────────────────
Classification output:
  "This email has probability 0.97 of being spam."
  Not: "This email IS spam."
  
Regression output:
  "Predicted price: $350,000 ± $15,000 (uncertainty!)"
  Not: "Price is exactly $350,000."

Training:
  "Find parameters that make the observed data MOST LIKELY."
  This IS probability theory — maximum likelihood estimation.
  
Evaluation:
  "Test accuracy 91.2% ± 0.3% (95% confidence interval)"
  Statistical testing is needed to compare models.

Language models (GPT, Claude):
  "Next token probabilities: ['cat' → 0.35, 'dog' → 0.28, ...]"
  Literally outputs a probability distribution over vocabulary.
─────────────────────────────────────────────────────────────────────
```

Understanding probability lets you answer:
- Why is cross-entropy the right loss function for classification?
- Why does L2 regularization correspond to a Gaussian prior?
- Why is the area under the ROC curve (AUC) a probability?
- When is your model's confidence actually calibrated?

---

## 2. Probability Basics

### Sample Space and Events

```
Definitions:
────────────────────────────────────────────────────────────────────
Sample Space Ω: all possible outcomes
  Coin flip: Ω = {H, T}
  Die roll:  Ω = {1, 2, 3, 4, 5, 6}
  A pixel value: Ω = {0, 1, ..., 255}
  A real-valued measurement: Ω = ℝ

Event A: a subset of Ω (a set of outcomes we care about)
  "rolling an even number": A = {2, 4, 6}

Probability P(A): a number in [0, 1]
  P(∅) = 0             (impossible event)
  P(Ω) = 1             (certain event)
  P(Aᶜ) = 1 - P(A)    (complement: probability of NOT A)
```

### Basic Probability Rules

```
Union: P(A ∪ B) = P(A) + P(B) - P(A ∩ B)
  "A or B" — subtract intersection to avoid double-counting

Intersection: P(A ∩ B)
  "A and B" — both occur simultaneously

De Morgan's Laws:
  P(Aᶜ ∩ Bᶜ) = P((A ∪ B)ᶜ) = 1 - P(A ∪ B)
  P(Aᶜ ∪ Bᶜ) = P((A ∩ B)ᶜ) = 1 - P(A ∩ B)
```

```python
import numpy as np
from scipy import stats

# Simulate probability with many coin flips
np.random.seed(42)
n_flips = 100000
flips = np.random.choice(['H', 'T'], size=n_flips)

p_heads = np.mean(flips == 'H')
print(f"P(Heads) ≈ {p_heads:.4f}  (true: 0.5000)")

# Union: P(A or B) for a die
n_rolls = 100000
rolls = np.random.randint(1, 7, n_rolls)

A = rolls % 2 == 0     # even
B = rolls > 4          # greater than 4

p_A = np.mean(A)
p_B = np.mean(B)
p_AB = np.mean(A & B)
p_A_or_B = np.mean(A | B)

print(f"\nDie roll: P(even) = {p_A:.4f}, P(>4) = {p_B:.4f}")
print(f"P(even ∩ >4) = {p_AB:.4f}  [should be P({6}) = 1/6 ≈ 0.1667]")
print(f"P(even ∪ >4) = {p_A_or_B:.4f}")
print(f"From formula: {p_A + p_B - p_AB:.4f}  (should match!)")
```

---

## 3. Conditional Probability

Conditional probability answers: "Given that I know B happened, what's the probability of A?"

```
Conditional probability:
  P(A | B) = P(A ∩ B) / P(B)   [assuming P(B) > 0]

Read as: "probability of A GIVEN B"

Intuition:
  Restricting our universe to only cases where B occurred,
  then asking: in that restricted universe, how often does A occur?
```

```python
import numpy as np

# Medical test example — foundation of Bayesian reasoning
# Suppose:
# - 1% of population has disease D
# - Test has 95% sensitivity: P(Test+ | D) = 0.95
# - Test has 90% specificity: P(Test- | No D) = 0.90
# Question: given a positive test, what's P(D | Test+)?

# Simulate a population
np.random.seed(42)
n = 1_000_000
has_disease = np.random.random(n) < 0.01      # 1% prevalence

# Conditional probabilities of test result
true_positive_rate  = 0.95   # P(Test+ | D)
false_positive_rate = 0.10   # P(Test+ | No D) = 1 - specificity

test_positive = np.where(
    has_disease,
    np.random.random(n) < true_positive_rate,   # sick and tested positive
    np.random.random(n) < false_positive_rate    # healthy and tested positive
)

# Among those who tested positive, what fraction actually have the disease?
test_pos_mask = test_positive
print(f"Total tested positive: {test_pos_mask.sum()}")
print(f"True positives (D ∩ Test+): {(has_disease & test_pos_mask).sum()}")

p_disease_given_pos = has_disease[test_pos_mask].mean()
print(f"\nP(D | Test+) = {p_disease_given_pos:.4f}")
# Surprisingly low! ~48% even with a good test!
# This is why mass screening programs need very careful analysis.
```

### The Law of Total Probability

```
If B₁, B₂, ..., Bₙ partition Ω (exhaustive, mutually exclusive):

P(A) = Σᵢ P(A | Bᵢ) P(Bᵢ)

Example: P(Test+) = P(Test+ | D)P(D) + P(Test+ | No D)P(No D)
                  = 0.95 × 0.01 + 0.10 × 0.99
                  = 0.0095 + 0.099 = 0.1085
```

---

## 4. Bayes' Theorem

Bayes' theorem is how you update your beliefs when new evidence arrives. It is the backbone of Bayesian machine learning.

```
Bayes' Theorem:
  P(A | B) = P(B | A) P(A) / P(B)

Named parts (in the ML context):
  P(A)       = Prior         — what we believed BEFORE seeing data
  P(B | A)   = Likelihood    — probability of data given hypothesis
  P(B)       = Evidence      — normalizing constant (marginal probability)
  P(A | B)   = Posterior     — updated belief AFTER seeing data
```

### Deriving Bayes' Theorem

```
From conditional probability:
  P(A | B) = P(A ∩ B) / P(B)    ... (1)
  P(B | A) = P(A ∩ B) / P(A)    ... (2)

From (2): P(A ∩ B) = P(B | A) P(A)

Substitute into (1):
  P(A | B) = P(B | A) P(A) / P(B)     ← Bayes' Theorem ✓
```

### Spam Filter Example

```python
import numpy as np

# Naive Bayes spam classifier
# Training data statistics
# P(spam) = 0.3 (30% of emails are spam)
# P("free" | spam) = 0.8  (80% of spam contains "free")
# P("free" | not spam) = 0.1  (10% of legitimate emails contain "free")

p_spam    = 0.30
p_not_spam = 0.70
p_free_given_spam     = 0.80
p_free_given_not_spam = 0.10

# Observed: email contains "free"
# What's P(spam | free)?

# P(free) = P(free|spam)P(spam) + P(free|not spam)P(not spam)
p_free = p_free_given_spam * p_spam + p_free_given_not_spam * p_not_spam
print(f"P('free' in email) = {p_free:.4f}")

# Bayes:
p_spam_given_free = (p_free_given_spam * p_spam) / p_free
print(f"P(spam | 'free') = {p_spam_given_free:.4f}")
# ≈ 0.774 — seeing 'free' raises spam probability from 30% to 77%

# Another word: "meeting" - P("meeting" | spam) = 0.05, P("meeting" | not spam) = 0.4
p_meeting_given_spam     = 0.05
p_meeting_given_not_spam = 0.40

# Naive Bayes: multiply likelihoods (assumes independence — the "naive" part)
# P(spam | free, meeting) ∝ P(spam) × P(free|spam) × P(meeting|spam)
p_spam_both = p_spam * p_free_given_spam * p_meeting_given_spam
p_notspam_both = p_not_spam * p_free_given_not_spam * p_meeting_given_not_spam

# Normalize
total = p_spam_both + p_notspam_both
p_spam_given_both = p_spam_both / total
print(f"\nP(spam | 'free' AND 'meeting') = {p_spam_given_both:.4f}")
# Meeting is very un-spam-like → posterior drops


# ── Full Naive Bayes classifier from scratch ──────────────────────────────
class NaiveBayesClassifier:
    """Simple Gaussian Naive Bayes for continuous features."""
    
    def fit(self, X, y):
        self.classes_ = np.unique(y)
        self.priors_ = {}
        self.means_  = {}
        self.stds_   = {}
        
        for cls in self.classes_:
            X_cls = X[y == cls]
            self.priors_[cls] = len(X_cls) / len(X)
            self.means_[cls]  = X_cls.mean(axis=0)
            self.stds_[cls]   = X_cls.std(axis=0) + 1e-9  # avoid division by zero
        
        return self
    
    def _gaussian_log_likelihood(self, X, mean, std):
        """log P(X | class) assuming Gaussian distribution per feature."""
        return -0.5 * np.sum(((X - mean) / std)**2 + np.log(2 * np.pi * std**2), axis=1)
    
    def predict_proba(self, X):
        log_probs = np.zeros((len(X), len(self.classes_)))
        
        for i, cls in enumerate(self.classes_):
            log_prior = np.log(self.priors_[cls])
            log_likelihood = self._gaussian_log_likelihood(
                X, self.means_[cls], self.stds_[cls]
            )
            log_probs[:, i] = log_prior + log_likelihood
        
        # Normalize (log-sum-exp for numerical stability)
        log_probs -= log_probs.max(axis=1, keepdims=True)
        probs = np.exp(log_probs)
        return probs / probs.sum(axis=1, keepdims=True)
    
    def predict(self, X):
        return self.classes_[np.argmax(self.predict_proba(X), axis=1)]

# Test on iris-like data
from sklearn.datasets import load_iris
from sklearn.model_selection import train_test_split

iris = load_iris()
X_train, X_test, y_train, y_test = train_test_split(
    iris.data, iris.target, test_size=0.2, random_state=42)

nb = NaiveBayesClassifier()
nb.fit(X_train, y_train)
preds = nb.predict(X_test)
accuracy = np.mean(preds == y_test)
print(f"\nNaive Bayes accuracy on Iris: {accuracy:.2%}")
```

---

## 5. Independence

Two events A and B are independent if knowing B tells you nothing about A.

```
Independence condition:
  P(A | B) = P(A)    ← knowing B doesn't change P(A)
  P(B | A) = P(B)    ← equivalently
  P(A ∩ B) = P(A) P(B)    ← multiplication rule for independent events

Examples:
  Independent:      Two coin flips
  NOT independent:  Weather today and weather tomorrow
                    Whether an email has "free" and whether it's spam
                    
Conditional independence:
  A ⊥ B | C  means  P(A | B, C) = P(A | C)
  "A and B are independent GIVEN C"
  This is the key assumption in Naive Bayes — word occurrences are
  independent GIVEN the class (spam/not spam).
```

```python
import numpy as np

np.random.seed(42)
n = 100000

# Independent events
coin1 = np.random.choice(['H', 'T'], n)
coin2 = np.random.choice(['H', 'T'], n)   # independent flip

p_c1_heads = np.mean(coin1 == 'H')
p_c2_heads_given_c1_heads = np.mean(coin2[coin1 == 'H'] == 'H')

print(f"P(coin2=H) = {np.mean(coin2 == 'H'):.4f}")
print(f"P(coin2=H | coin1=H) = {p_c2_heads_given_c1_heads:.4f}")
print(f"Independent: {np.isclose(np.mean(coin2=='H'), p_c2_heads_given_c1_heads, atol=0.01)}")

# Dependent events
temp_high = np.random.random(n) > 0.4   # 60% chance of hot day
# Ice cream sales depend on temperature
icecream = temp_high & (np.random.random(n) > 0.2)  # more likely when hot

p_icecream = np.mean(icecream)
p_icecream_given_hot = np.mean(icecream[temp_high])
print(f"\nP(ice cream sold) = {p_icecream:.4f}")
print(f"P(ice cream sold | hot day) = {p_icecream_given_hot:.4f}")
print(f"NOT independent: these differ!")
```

---

## 6. Random Variables

A **random variable** (RV) is a variable whose value is determined by a random process.

```
Discrete RV: takes countable values
  Examples: number of heads in 10 flips (0-10)
            which class a sample belongs to (0, 1, 2)
            word frequency in a document (0, 1, 2, 3, ...)

Continuous RV: takes any value in an interval
  Examples: temperature, income, pixel intensity
            weight of a person, time until next event

PMF (Probability Mass Function) — for discrete X:
  P(X = k) = probability that X takes value k
  Σₖ P(X = k) = 1

PDF (Probability Density Function) — for continuous X:
  f(x) dx = P(x ≤ X ≤ x + dx)  (probability of X in tiny interval)
  ∫ f(x) dx = 1  (total probability = 1)
  P(a ≤ X ≤ b) = ∫ₐᵇ f(x) dx

CDF (Cumulative Distribution Function):
  F(x) = P(X ≤ x) = ∫₋∞ˣ f(t) dt
```

---

## 7. Expectation and Variance

### Expectation (Mean)

```
Expectation E[X]: the average value of X, weighted by probability

Discrete:    E[X] = Σₖ k · P(X = k)
Continuous:  E[X] = ∫ x · f(x) dx

Linearity of expectation (always true, even for dependent variables!):
  E[aX + bY] = a·E[X] + b·E[Y]

Examples:
  Fair die: E[X] = 1·(1/6) + 2·(1/6) + ... + 6·(1/6) = 21/6 = 3.5
  Biased coin (P(H)=0.7): E[flips_until_heads] = 1/0.7 ≈ 1.43
```

### Variance

```
Variance Var(X): measure of spread around the mean

  Var(X) = E[(X - μ)²] = E[X²] - (E[X])²

Standard deviation σ(X) = √Var(X)  (in same units as X)

Key properties:
  Var(aX) = a² Var(X)    (scaling multiplies variance by a²)
  Var(X + c) = Var(X)    (shifting doesn't change variance)
  Var(X + Y) = Var(X) + Var(Y)  (only if X and Y are independent!)
```

```python
import numpy as np
from scipy import stats

# ── Expectation ────────────────────────────────────────────────────────────
die_outcomes = np.array([1, 2, 3, 4, 5, 6])
die_probs    = np.array([1/6] * 6)

E_die = np.sum(die_outcomes * die_probs)
print(f"E[fair die] = {E_die:.4f}")   # 3.5

# Verify by simulation
simulated_rolls = np.random.randint(1, 7, 100000)
print(f"E[die] by simulation = {simulated_rolls.mean():.4f}")

# ── Variance ────────────────────────────────────────────────────────────────
var_die = np.sum((die_outcomes - E_die)**2 * die_probs)
std_die = np.sqrt(var_die)
print(f"Var[die] = {var_die:.4f}")    # 35/12 ≈ 2.917
print(f"Std[die] = {std_die:.4f}")   # ≈ 1.708

# Sample variance (divides by n-1, not n, for unbiased estimate)
print(f"Simulated std = {simulated_rolls.std(ddof=1):.4f}")  # ddof=1 for unbiased

# ── Biased vs Unbiased estimator ──────────────────────────────────────────
# True population variance: Var(X) = E[(X-μ)²]
# Sample variance:          s² = (1/(n-1)) Σ(xᵢ - x̄)²
# The n-1 (Bessel's correction) makes the estimator unbiased

np.random.seed(42)
true_var = 4.0
samples = np.random.normal(0, 2, 1000)   # N(0, σ²=4)
biased_var   = np.var(samples, ddof=0)   # divides by n
unbiased_var = np.var(samples, ddof=1)   # divides by n-1 (better!)
print(f"\nTrue variance: {true_var}")
print(f"Biased estimate (÷n):   {biased_var:.4f}")
print(f"Unbiased estimate (÷n-1): {unbiased_var:.4f}")
```

---

## 8. Covariance and Correlation

```
Covariance: measures HOW MUCH two variables change together

Cov(X, Y) = E[(X - μX)(Y - μY)]
           = E[XY] - E[X]E[Y]

Positive: X and Y tend to increase together (income and education)
Negative: X goes up as Y goes down (temperature and heating bills)
Zero:     X and Y are uncorrelated (but not necessarily independent!)

Problem with covariance: units! Cov(height_m, weight_kg) vs Cov(height_cm, weight_g)
give wildly different numbers even for the same data.

Correlation (Pearson): standardized covariance, unit-free
  ρ(X,Y) = Cov(X,Y) / (σX · σY)

Range: [-1, +1]
  ρ = +1: perfect positive linear relationship
  ρ = -1: perfect negative linear relationship
  ρ = 0:  no linear relationship (might still have nonlinear relationship!)
```

```python
import numpy as np
import matplotlib.pyplot as plt

np.random.seed(42)
n = 500

# Generate correlated data
rho_target = 0.8
mean = [0, 0]
cov_matrix = [[1, rho_target],
              [rho_target, 1]]
X, Y = np.random.multivariate_normal(mean, cov_matrix, n).T

# Compute correlation
print(f"Target correlation: {rho_target}")
print(f"Computed Pearson r: {np.corrcoef(X, Y)[0,1]:.4f}")

# Full correlation matrix (used in EDA)
features = np.column_stack([
    X,                          # feature 1
    Y,                          # feature 2 (correlated with 1)
    np.random.randn(n),         # feature 3 (uncorrelated)
    2*X + np.random.randn(n),   # feature 4 (highly correlated with 1)
])

corr_matrix = np.corrcoef(features.T)
print("\nCorrelation matrix:")
print(corr_matrix.round(3))

# ── Visualization ─────────────────────────────────────────────────────────
fig, axes = plt.subplots(1, 4, figsize=(20, 4))
rhos = [0.9, 0.5, 0.0, -0.8]

for ax, r in zip(axes, rhos):
    cov = [[1, r], [r, 1]]
    x, y = np.random.multivariate_normal([0,0], cov, 300).T
    computed_r = np.corrcoef(x, y)[0,1]
    ax.scatter(x, y, alpha=0.4, s=20, color='steelblue')
    ax.set_title(f"ρ = {r} (computed: {computed_r:.2f})")
    ax.set_xlabel("X")
    ax.set_ylabel("Y")
    ax.grid(True, alpha=0.3)

plt.tight_layout()
plt.show()

# ── Warning: Correlation ≠ Linear Relationship Only ────────────────────────
# Anscombe's quartet style: nonlinear data with zero correlation
x_nl = np.linspace(-3, 3, 100)
y_nl = x_nl**2 + np.random.randn(100) * 0.5   # U-shaped, corr ≈ 0!
print(f"\nNonlinear (y=x²): correlation = {np.corrcoef(x_nl, y_nl)[0,1]:.4f}")
# Nearly 0, yet there's a perfect nonlinear relationship!
```

---

## 9. Key Probability Distributions

```
Distribution Guide for ML:
──────────────────────────────────────────────────────────────────────────
Distribution    Use case in ML
──────────────────────────────────────────────────────────────────────────
Bernoulli       Single binary outcome (coin flip, email spam/not)
Binomial        k successes in n independent Bernoulli trials
Uniform         Prior when all values are equally likely; initialization
Normal          Weights, activations, residuals; Central Limit Theorem
Multivariate    Feature space, Gaussian mixtures, Gaussian processes
  Normal
Categorical     Output of softmax (k mutually exclusive classes)
Beta            Probabilities as outputs (0 to 1), Bayesian conjugate
Dirichlet       Distribution over probability simplex, topic modeling
Exponential     Time until first event, inter-arrival times
Poisson         Count of events in fixed time/space
Log-normal      Income, price, multiplicative processes
──────────────────────────────────────────────────────────────────────────
```

```python
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

fig, axes = plt.subplots(2, 4, figsize=(20, 9))
axes = axes.flatten()

# ── Bernoulli ──────────────────────────────────────────────────────────────
ax = axes[0]
p = 0.7
k = np.array([0, 1])
pmf = np.array([1-p, p])
ax.bar(k, pmf, color=['#F44336', '#4CAF50'], width=0.5)
ax.set_title(f"Bernoulli(p={p})\n(email spam: 0=ham, 1=spam)")
ax.set_xticks([0, 1])
ax.set_xticklabels(['0 (no)', '1 (yes)'])
ax.set_ylabel("P(X=k)")

# ── Binomial ───────────────────────────────────────────────────────────────
ax = axes[1]
n_b, p_b = 20, 0.3
k = np.arange(0, n_b + 1)
pmf = stats.binom.pmf(k, n_b, p_b)
ax.bar(k, pmf, color='steelblue', alpha=0.8)
ax.set_title(f"Binomial(n={n_b}, p={p_b})\n(correct predictions in 20 tries)")
ax.set_xlabel("k successes")
ax.set_ylabel("P(X=k)")

# ── Uniform ────────────────────────────────────────────────────────────────
ax = axes[2]
a, b = 0, 10
x = np.linspace(-1, 11, 500)
ax.fill_between(x, stats.uniform.pdf(x, a, b-a), alpha=0.6, color='purple')
ax.set_title(f"Uniform({a}, {b})\n(uninformative prior)")
ax.set_xlabel("x")
ax.set_ylabel("f(x)")

# ── Exponential ────────────────────────────────────────────────────────────
ax = axes[3]
x = np.linspace(0, 10, 500)
for lam, color in [(0.5, 'blue'), (1.0, 'green'), (2.0, 'red')]:
    ax.plot(x, stats.expon.pdf(x, scale=1/lam), lw=2, label=f'λ={lam}', color=color)
ax.set_title("Exponential(λ)\n(time between events)")
ax.legend()

# ── Poisson ────────────────────────────────────────────────────────────────
ax = axes[4]
k = np.arange(0, 20)
for lam, color in [(1, 'blue'), (4, 'green'), (8, 'red')]:
    ax.plot(k, stats.poisson.pmf(k, lam), 'o-', lw=2, label=f'λ={lam}', color=color, ms=5)
ax.set_title("Poisson(λ)\n(count of rare events)")
ax.legend()

# ── Normal ─────────────────────────────────────────────────────────────────
ax = axes[5]
x = np.linspace(-4, 4, 500)
for mu, sigma, color in [(0, 1, 'blue'), (0, 2, 'red'), (1, 0.5, 'green')]:
    label = f'μ={mu}, σ={sigma}'
    ax.plot(x, stats.norm.pdf(x, mu, sigma), lw=2, label=label, color=color)
ax.set_title("Normal(μ, σ²)\n(weights, residuals, CLT limit)")
ax.legend()

# ── Beta ────────────────────────────────────────────────────────────────────
ax = axes[6]
x = np.linspace(0.001, 0.999, 500)
for a_b, b_b, color in [(2, 5, 'blue'), (0.5, 0.5, 'red'), (5, 2, 'green')]:
    label = f'α={a_b}, β={b_b}'
    ax.plot(x, stats.beta.pdf(x, a_b, b_b), lw=2, label=label, color=color)
ax.set_title("Beta(α, β)\n(probability values, Bayesian prior for p)")
ax.legend()

# ── Log-normal ─────────────────────────────────────────────────────────────
ax = axes[7]
x = np.linspace(0.01, 10, 500)
for s, color in [(0.5, 'blue'), (1.0, 'red'), (1.5, 'green')]:
    ax.plot(x, stats.lognorm.pdf(x, s), lw=2, label=f'σ={s}', color=color)
ax.set_title("Log-Normal(μ, σ)\n(income, price distributions)")
ax.legend()

fig.suptitle("Common Probability Distributions in ML", fontsize=15, fontweight='bold')
plt.tight_layout()
plt.show()
```

---

## 10. The Normal Distribution

The Normal (Gaussian) distribution is the most important distribution in ML. It appears everywhere due to the Central Limit Theorem.

```
Normal distribution: X ~ N(μ, σ²)

PDF: f(x) = (1/√(2πσ²)) exp(-(x-μ)²/(2σ²))

Parameters:
  μ = mean (center of the bell curve)
  σ = standard deviation (width of the bell curve)
  σ² = variance

The 68-95-99.7 Rule:
  P(μ - σ  ≤ X ≤ μ + σ)  ≈ 68.27%  (1σ interval)
  P(μ - 2σ ≤ X ≤ μ + 2σ) ≈ 95.45%  (2σ interval)
  P(μ - 3σ ≤ X ≤ μ + 3σ) ≈ 99.73%  (3σ interval)
  
  "5-sigma event": P outside 5σ ≈ 0.00006% (essentially impossible)

Standard Normal: Z ~ N(0, 1)
  Standardization: Z = (X - μ) / σ
  P(Z > 1.96) ≈ 0.025  (one-sided)
  P(-1.96 ≤ Z ≤ 1.96) ≈ 0.95  (two-sided 95% interval)
```

```python
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

# ── Standard Normal and Z-scores ────────────────────────────────────────────
mu, sigma = 170, 10     # height: mean 170cm, std 10cm
x = np.linspace(130, 210, 500)

fig, axes = plt.subplots(1, 2, figsize=(14, 5))

# Normal distribution
ax = axes[0]
pdf = stats.norm.pdf(x, mu, sigma)
ax.plot(x, pdf, 'b-', lw=2)

# Fill 1σ, 2σ, 3σ regions
for n_sig, alpha, label in [(3, 0.1, '99.7%'), (2, 0.2, '95.4%'), (1, 0.4, '68.3%')]:
    lower, upper = mu - n_sig*sigma, mu + n_sig*sigma
    x_fill = np.linspace(lower, upper, 300)
    ax.fill_between(x_fill, stats.norm.pdf(x_fill, mu, sigma), alpha=alpha,
                    color='steelblue', label=f'±{n_sig}σ = {label}')

ax.set_title(f"Normal Distribution N(μ={mu}, σ={sigma})\n"
             f"68-95-99.7 Rule", fontsize=12)
ax.set_xlabel("Height (cm)")
ax.set_ylabel("Density")
ax.legend()

# Z-score visualization
ax = axes[1]
z_values = np.linspace(-4, 4, 500)
ax.plot(z_values, stats.norm.pdf(z_values, 0, 1), 'b-', lw=2)
ax.fill_between(z_values, stats.norm.pdf(z_values, 0, 1),
                where=(z_values >= -1.96) & (z_values <= 1.96),
                alpha=0.4, color='green', label='95% interval: [-1.96, 1.96]')
ax.axvline(-1.96, color='red', ls='--', alpha=0.7)
ax.axvline( 1.96, color='red', ls='--', alpha=0.7)
ax.set_title("Standard Normal N(0,1)\n95% of data lies within ±1.96σ")
ax.set_xlabel("Z-score = (x - μ)/σ")
ax.legend()

plt.tight_layout()
plt.show()

# Computing probabilities
print("Normal distribution computations:")
height_dist = stats.norm(mu, sigma)
print(f"P(height > 180cm) = {height_dist.sf(180):.4f}")     # sf = 1 - cdf
print(f"P(160 < height < 185) = {height_dist.cdf(185) - height_dist.cdf(160):.4f}")
print(f"90th percentile height = {height_dist.ppf(0.90):.1f} cm")

# Z-scores
heights = np.array([155, 170, 185, 200])
z_scores = (heights - mu) / sigma
print(f"\nZ-scores for heights {heights}: {z_scores}")
```

---

## 11. The Multivariate Normal

```python
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

# ── Multivariate Normal: X ~ N(μ, Σ) ─────────────────────────────────────
# μ = mean vector (center)
# Σ = covariance matrix (shape, orientation, spread)

mu = np.array([0.0, 0.0])

# Three covariance matrices with different shapes
covs = {
    "Isotropic (I)": np.eye(2),                      # circle
    "Anisotropic diag": np.diag([3.0, 0.5]),         # elongated ellipse, aligned
    "Correlated": np.array([[2.0, 1.5], [1.5, 2.0]]) # rotated ellipse
}

fig, axes = plt.subplots(1, 3, figsize=(18, 5))

for ax, (label, Sigma) in zip(axes, covs.items()):
    # Sample from multivariate normal
    samples = np.random.multivariate_normal(mu, Sigma, 1000)
    
    # Compute PDF on a grid
    x_g = np.linspace(-5, 5, 100)
    y_g = np.linspace(-5, 5, 100)
    Xg, Yg = np.meshgrid(x_g, y_g)
    pos = np.dstack((Xg, Yg))
    rv = stats.multivariate_normal(mu, Sigma)
    Z = rv.pdf(pos)
    
    ax.contour(Xg, Yg, Z, levels=10, cmap='Blues')
    ax.scatter(samples[:, 0], samples[:, 1], alpha=0.2, s=10, color='steelblue')
    ax.set_title(f"{label}\nΣ = {Sigma.tolist()}", fontsize=11)
    ax.set_aspect('equal')
    ax.grid(True, alpha=0.3)

plt.tight_layout()
plt.show()

# ── Sampling from multivariate normal using Cholesky ──────────────────────
mu = np.array([2.0, 5.0])
Sigma = np.array([[4.0, 2.0],
                  [2.0, 3.0]])

L = np.linalg.cholesky(Sigma)
z = np.random.randn(1000, 2)          # standard normal samples
x = z @ L.T + mu                      # transform to N(μ, Σ)

print(f"Sample mean: {x.mean(axis=0)}")     # ≈ [2, 5]
print(f"Sample cov: \n{np.cov(x.T)}")       # ≈ [[4,2],[2,3]]
```

---

## 12. Central Limit Theorem

The Central Limit Theorem (CLT) is why the normal distribution appears everywhere.

```
Central Limit Theorem:
─────────────────────────────────────────────────────────────────────
Let X₁, X₂, ..., Xₙ be independent, identically distributed (iid)
random variables with mean μ and variance σ².

Then the SAMPLE MEAN X̄ₙ = (1/n)Σ Xᵢ, as n → ∞, has distribution:

  X̄ₙ → N(μ, σ²/n)

In other words:
  (X̄ₙ - μ) / (σ/√n) → N(0, 1)

The DISTRIBUTION OF THE ORIGINAL DATA doesn't matter!
Whether it's uniform, exponential, Poisson, binomial...
The SAMPLE MEAN will be approximately normal for large n.
─────────────────────────────────────────────────────────────────────
```

```python
import numpy as np
import matplotlib.pyplot as plt

np.random.seed(42)

# Demonstrate CLT for different source distributions
distributions = {
    "Uniform [0,1]":   lambda n: np.random.uniform(0, 1, n),
    "Exponential(1)":  lambda n: np.random.exponential(1, n),
    "Bernoulli(0.3)":  lambda n: (np.random.random(n) < 0.3).astype(float),
}

sample_sizes = [1, 5, 30, 100]
n_experiments = 10000

fig, axes = plt.subplots(len(distributions), len(sample_sizes),
                          figsize=(20, 12))

for row, (dist_name, sampler) in enumerate(distributions.items()):
    for col, n in enumerate(sample_sizes):
        ax = axes[row, col]
        
        # Generate n_experiments sample means, each from n samples
        sample_means = np.array([sampler(n).mean() for _ in range(n_experiments)])
        
        ax.hist(sample_means, bins=50, density=True, color='steelblue',
                alpha=0.7, edgecolor='white')
        
        # Overlay normal approximation
        mu_approx = sample_means.mean()
        std_approx = sample_means.std()
        x_fit = np.linspace(sample_means.min(), sample_means.max(), 200)
        from scipy import stats as scipy_stats
        ax.plot(x_fit, scipy_stats.norm.pdf(x_fit, mu_approx, std_approx),
                'r-', lw=2, label='Normal fit')
        
        if row == 0:
            ax.set_title(f"n = {n}", fontsize=12, fontweight='bold')
        if col == 0:
            ax.set_ylabel(dist_name.split('[')[0].strip(), fontsize=10)
        
        ax.tick_params(labelsize=8)

fig.suptitle("Central Limit Theorem: Sample Mean Distribution\n"
             "(Left: n=1 looks like source. Right: n=100 always looks Normal)",
             fontsize=14, fontweight='bold')
plt.tight_layout()
plt.show()

# MLimplication: even though individual data points might not be normal,
# the AVERAGE of many samples is, which justifies many statistical tests.
```

---

## 13. Maximum Likelihood Estimation (MLE)

MLE is the probabilistic foundation for why we minimize loss functions. It answers: "Given data, what parameter values are most likely?"

```
MLE Principle:
  Choose parameters θ* that maximize the probability of observing the data.

  θ* = argmax P(data | θ)
             θ

  = argmax ∏ᵢ P(xᵢ | θ)   [assuming iid samples]
  
  = argmax Σᵢ log P(xᵢ | θ)   [log-likelihood, easier to optimize]

  = argmin -Σᵢ log P(xᵢ | θ)   [minimizing negative log-likelihood]
```

### Linear Regression IS MLE

```
Assume: y = w^T x + ε,  where ε ~ N(0, σ²)

Then: P(y | x, w) = N(y; w^T x, σ²)

Log-likelihood:
  log P(y | x, w) = -1/(2σ²) (y - w^T x)² - constant

Maximizing log P = Minimizing (y - w^T x)²  = Minimizing MSE!

Conclusion: fitting linear regression with MSE IS maximizing the 
likelihood of data under a Gaussian noise assumption.
```

### Logistic Regression IS MLE with Bernoulli

```
Assume: y | x ~ Bernoulli(sigmoid(w^T x))

Log-likelihood:
  log P(y | x, w) = y log(σ(w^Tx)) + (1-y) log(1 - σ(w^Tx))

Maximizing log P = Minimizing Binary Cross-Entropy!

Conclusion: logistic regression with BCE loss IS MLE under Bernoulli model.
```

```python
import numpy as np
from scipy.optimize import minimize

# ── MLE for a Gaussian — estimate μ and σ ─────────────────────────────────
true_mu, true_sigma = 5.0, 2.0
data = np.random.normal(true_mu, true_sigma, 1000)

def negative_log_likelihood_gaussian(params, data):
    mu, log_sigma = params
    sigma = np.exp(log_sigma)   # ensure sigma > 0
    n = len(data)
    # NLL = n*log(sigma) + (1/(2σ²)) Σ(xᵢ - μ)²  + constant
    return n * log_sigma + 0.5 * np.sum((data - mu)**2) / sigma**2

# Optimize
result = minimize(
    negative_log_likelihood_gaussian,
    x0=[0.0, 0.0],   # initial guess: mu=0, log_sigma=0
    args=(data,),
    method='BFGS'
)
mu_mle, sigma_mle = result.x[0], np.exp(result.x[1])
print(f"MLE estimates: μ = {mu_mle:.4f} (true: {true_mu}), "
      f"σ = {sigma_mle:.4f} (true: {true_sigma})")

# Analytical MLE for Gaussian:
print(f"Analytical MLE: μ̂ = {data.mean():.4f}, σ̂ = {data.std():.4f}")
# The sample mean and std ARE the MLE estimates!

# ── Visualize likelihood surface ──────────────────────────────────────────
small_data = np.random.normal(3.0, 1.5, 30)   # small dataset to see clearly

mu_vals    = np.linspace(0, 6, 100)
sigma_vals = np.linspace(0.5, 4.0, 100)
Mu, Sigma  = np.meshgrid(mu_vals, sigma_vals)

# Log-likelihood for each (mu, sigma) pair
def log_likelihood_grid(mu, sigma, data):
    n = len(data)
    return (-n * np.log(sigma)
            - 0.5 * np.sum((data[:, None, None] - mu)**2, axis=0) / sigma**2)

LL = log_likelihood_grid(Mu, Sigma, small_data)

import matplotlib.pyplot as plt
fig, ax = plt.subplots(figsize=(8, 6))
contour = ax.contour(Mu, Sigma, LL, levels=20, cmap='viridis')
plt.colorbar(contour, ax=ax, label='Log-likelihood')

# Mark MLE (should be at data.mean(), data.std())
ax.scatter([small_data.mean()], [small_data.std()],
           color='red', s=100, zorder=5, label='MLE estimate')
ax.set_xlabel("μ (mean)")
ax.set_ylabel("σ (std)")
ax.set_title("Log-Likelihood Surface for Gaussian Model\n"
             "(Maximum at red point = MLE estimate)", fontsize=12)
ax.legend()
plt.tight_layout()
plt.show()
```

---

## 14. MAP Estimation and Regularization

MAP (Maximum A Posteriori) estimation adds a prior to MLE. It answers: "Given a prior belief AND data, what's the most probable parameter value?"

```
MAP:
  θ_MAP = argmax P(θ | data) = argmax P(data | θ) P(θ) / P(data)
                              = argmax [P(data | θ) P(θ)]
                              = argmax [log P(data | θ) + log P(θ)]
                                                ↑               ↑
                                           log-likelihood    log-prior
                                           = neg cross-entropy  = regularization!

Prior choices and their regularization equivalents:
─────────────────────────────────────────────────────────────────────────────
Gaussian prior: P(w) ∝ exp(-λ||w||²) → log P(w) = -λ||w||²
  MAP → minimize MSE + λ||w||² = Ridge Regression (L2 regularization)!

Laplace prior:  P(w) ∝ exp(-λ||w||₁) → log P(w) = -λ||w||₁
  MAP → minimize MSE + λ||w||₁ = Lasso Regression (L1 regularization)!

No prior (uniform): P(w) = constant → log P(w) = 0
  MAP → minimize MSE only = ordinary least squares = MLE

Regularization IS Bayesian — it encodes your prior belief that weights
should be small (Ridge) or sparse (Lasso).
─────────────────────────────────────────────────────────────────────────────
```

```python
import numpy as np

# Ridge regression via MAP (Gaussian prior)
def ridge_regression_map(X, y, lambda_reg=1.0):
    """
    MAP estimate with Gaussian prior: minimize ||Xw - y||² + λ||w||²
    Analytical solution: w = (X^T X + λI)^(-1) X^T y
    """
    n, p = X.shape
    A = X.T @ X + lambda_reg * np.eye(p)
    b = X.T @ y
    return np.linalg.solve(A, b)

np.random.seed(42)
n, p = 100, 50   # more features (p=50) than effective samples → need regularization!
X = np.random.randn(n, p)
true_w = np.zeros(p); true_w[:5] = [2.0, -1.0, 1.5, -0.5, 0.8]
y = X @ true_w + np.random.randn(n) * 0.5

# Split
X_train, X_test = X[:80], X[80:]
y_train, y_test = y[:80], y[80:]

results = {}
for lam in [0.0, 0.1, 1.0, 10.0, 100.0]:
    w = ridge_regression_map(X_train, y_train, lam)
    train_mse = np.mean((X_train @ w - y_train)**2)
    test_mse  = np.mean((X_test  @ w - y_test)**2)
    results[lam] = (train_mse, test_mse, np.linalg.norm(w))
    print(f"λ={lam:5.1f}: train_MSE={train_mse:.4f}, "
          f"test_MSE={test_mse:.4f}, ||w||={results[lam][2]:.4f}")

# λ=0 (OLS) overfits massively — test_MSE >> train_MSE
# Increasing λ reduces test_MSE (better generalization) up to a point
```

---

## 15. Information Theory for ML

Information theory connects probability to machine learning loss functions.

### Entropy

```
Entropy H(X): measures the average uncertainty/information content of X.

For discrete distribution P over k outcomes:
  H(X) = -Σₖ P(xₖ) log P(xₖ)   (in bits if log₂, nats if log)

Properties:
  H(X) ≥ 0  (always non-negative)
  H(X) = 0 iff P is deterministic (one outcome has probability 1)
  H(X) is MAXIMUM for uniform distribution: H = log k

Intuition:
  Low entropy = more predictable = less information
  High entropy = less predictable = more information

  "cat"    always follows "the"     → low entropy (certain)
  next word after "I"               → high entropy (many possibilities)
```

```python
import numpy as np

def entropy(probs, base=2):
    """Compute entropy H(X) from probability distribution."""
    probs = np.asarray(probs, dtype=float)
    probs = probs[probs > 0]   # ignore zero-probability events
    return -np.sum(probs * np.log(probs + 1e-10) / np.log(base))

# Entropy examples
deterministic = [1.0, 0.0, 0.0]          # one certain outcome
balanced      = [1/3, 1/3, 1/3]          # maximum entropy for 3 outcomes
skewed        = [0.9, 0.09, 0.01]        # mostly certain

print(f"Deterministic: H = {entropy(deterministic):.4f} bits")   # 0.0
print(f"Balanced:      H = {entropy(balanced):.4f} bits")         # log2(3) = 1.585
print(f"Skewed:        H = {entropy(skewed):.4f} bits")           # low

# Binary entropy (special case: k=2)
import matplotlib.pyplot as plt
p_vals = np.linspace(0.001, 0.999, 500)
h_binary = -p_vals * np.log2(p_vals) - (1-p_vals) * np.log2(1-p_vals)
plt.figure(figsize=(8, 4))
plt.plot(p_vals, h_binary, 'b-', lw=2)
plt.xlabel("P(X=1)")
plt.ylabel("Entropy (bits)")
plt.title("Binary Entropy H(p)\nMaximum at p=0.5 (most uncertain)")
plt.grid(True, alpha=0.3)
plt.axvline(0.5, color='red', ls='--', alpha=0.7)
plt.show()
```

### Cross-Entropy — The Classification Loss

```
Cross-entropy H(P, Q): measures how different distribution Q is from P.

H(P, Q) = -Σₖ P(xₖ) log Q(xₖ)

When P = true labels (one-hot) and Q = model predictions:
  H(P, Q) = -Σₖ yₖ log(ŷₖ)  =  -log(ŷ_correct_class)

So cross-entropy loss = -log(probability assigned to the correct class).

Intuition:
  If model assigns P=0.99 to correct class: loss = -log(0.99) = 0.01  (low)
  If model assigns P=0.01 to correct class: loss = -log(0.01) = 4.6   (high!)

Relation to entropy:
  H(P, Q) = H(P) + D_KL(P||Q)
  If P is deterministic (true labels), H(P) = 0
  → Cross-entropy = KL divergence from true to predicted distribution
```

```python
import numpy as np

def cross_entropy(y_true_probs, y_pred_probs, eps=1e-7):
    """H(P, Q) = -Σ P log Q"""
    return -np.sum(y_true_probs * np.log(y_pred_probs + eps))

def kl_divergence(P, Q, eps=1e-7):
    """D_KL(P||Q) = Σ P log(P/Q)"""
    return np.sum(P * np.log((P + eps) / (Q + eps)))

# Example: 3-class classification
y_true = np.array([1, 0, 0])   # true label: class 0 (one-hot)

predictions = {
    "Confident correct": np.array([0.98, 0.01, 0.01]),
    "Uncertain":         np.array([0.40, 0.35, 0.25]),
    "Confident wrong":   np.array([0.01, 0.98, 0.01]),
}

print("Cross-entropy loss vs. prediction quality:")
for name, y_pred in predictions.items():
    ce = cross_entropy(y_true, y_pred)
    kl = kl_divergence(y_true, y_pred)
    print(f"  {name:25s}: H(P,Q) = {ce:.4f},  KL = {kl:.4f}")
```

### KL Divergence

```
KL Divergence D_KL(P||Q): "extra bits needed to encode P using Q's code"

D_KL(P||Q) = Σ P(x) log(P(x)/Q(x))  ≥ 0  (always!)
             = 0 iff P = Q

Properties:
  - NOT symmetric: D_KL(P||Q) ≠ D_KL(Q||P)
  - "Information gain" from updating Q to P
  
ML uses:
  - VAEs (Variational Autoencoders): KL between learned and prior distribution
  - Knowledge distillation: student model matches teacher's distribution
  - Reinforcement learning: policy optimization with KL constraint (PPO)
```

```python
from scipy.special import kl_div

# KL divergence is NOT symmetric
P = np.array([0.9, 0.1])
Q = np.array([0.5, 0.5])

kl_PQ = kl_divergence(P, Q)
kl_QP = kl_divergence(Q, P)
print(f"D_KL(P||Q) = {kl_PQ:.4f}")
print(f"D_KL(Q||P) = {kl_QP:.4f}")
print(f"Not equal: {not np.isclose(kl_PQ, kl_QP)}")
```

---

## 16. Hypothesis Testing in ML

Statistical testing lets you decide if observed differences are real or just noise.

```
The framework:
  H₀: Null hypothesis — "no difference" (default assumption)
  H₁: Alternative hypothesis — "there IS a difference"
  
  p-value: probability of observing data AT LEAST AS EXTREME as yours,
           IF H₀ were true.
  
  Decision rule:
    If p-value < α (significance level, usually 0.05): reject H₀
    If p-value ≥ α: fail to reject H₀ (don't "accept" H₀!)

Common misinterpretation:
  WRONG: "p-value = probability that H₀ is true"
  RIGHT: "p-value = probability of seeing this data if H₀ were true"
```

```python
import numpy as np
from scipy import stats

np.random.seed(42)

# ── t-test: Is model A significantly better than model B? ─────────────────
# 30 runs of cross-validation for each model
model_A_scores = np.random.normal(0.87, 0.03, 30)   # mean=0.87, std=0.03
model_B_scores = np.random.normal(0.85, 0.03, 30)   # mean=0.85

print("Model comparison:")
print(f"Model A: mean={model_A_scores.mean():.4f}, std={model_A_scores.std():.4f}")
print(f"Model B: mean={model_B_scores.mean():.4f}, std={model_B_scores.std():.4f}")
print(f"Difference: {(model_A_scores - model_B_scores).mean():.4f}")

# Paired t-test (better when same data splits used)
t_stat, p_value = stats.ttest_rel(model_A_scores, model_B_scores)
print(f"\nPaired t-test: t={t_stat:.4f}, p={p_value:.6f}")
if p_value < 0.05:
    print("Model A is significantly better (p < 0.05)")
else:
    print("No significant difference detected")

# ── chi-squared test: Is class distribution as expected? ─────────────────
# Testing if model predictions are balanced across classes
observed = np.array([320, 180, 250, 250])   # predicted class counts
expected = np.array([250, 250, 250, 250])   # expected if uniform

chi2_stat, p_value_chi2 = stats.chisquare(observed, expected)
print(f"\nChi-squared test: χ²={chi2_stat:.4f}, p={p_value_chi2:.4f}")
if p_value_chi2 < 0.05:
    print("Predictions are NOT uniformly distributed across classes")

# ── Bootstrap confidence interval — model-free approach ───────────────────
def bootstrap_ci(scores, n_bootstrap=10000, ci=0.95, seed=42):
    """Compute bootstrap confidence interval for the mean."""
    rng = np.random.default_rng(seed)
    bootstrap_means = np.array([
        rng.choice(scores, len(scores), replace=True).mean()
        for _ in range(n_bootstrap)
    ])
    alpha = 1 - ci
    lower = np.percentile(bootstrap_means, 100*alpha/2)
    upper = np.percentile(bootstrap_means, 100*(1-alpha/2))
    return lower, upper

ci_lower, ci_upper = bootstrap_ci(model_A_scores)
print(f"\nModel A accuracy: {model_A_scores.mean():.4f} "
      f"(95% CI: [{ci_lower:.4f}, {ci_upper:.4f}])")
```

---

## 17. Confidence Intervals

```
Confidence interval (CI): a range of values likely to contain the true parameter.

"95% CI [0.87, 0.91]" means:
  If we repeated this experiment many times, 95% of the intervals computed
  this way would contain the true population parameter.

Common WRONG interpretation: "There's a 95% probability the true value is in [0.87, 0.91]"
(The interval is fixed; the true value is a constant, not a random variable.)

For a normal distribution, 95% CI for mean:
  x̄ ± 1.96 × (σ / √n)
  where σ/√n is the "standard error of the mean"
```

```python
import numpy as np
from scipy import stats

np.random.seed(42)

# ── t-confidence interval for model accuracy ──────────────────────────────
scores = np.random.normal(0.88, 0.04, 20)   # 20 cross-validation scores

n = len(scores)
mean = scores.mean()
se = scores.std(ddof=1) / np.sqrt(n)   # standard error

# 95% CI using t-distribution (appropriate for small samples)
t_critical = stats.t.ppf(0.975, df=n-1)   # 97.5th percentile of t(n-1)
ci_lower = mean - t_critical * se
ci_upper = mean + t_critical * se

print(f"Sample mean: {mean:.4f}")
print(f"Standard error: {se:.4f}")
print(f"95% CI: [{ci_lower:.4f}, {ci_upper:.4f}]")
print(f"Width: {ci_upper - ci_lower:.4f}")

# As n increases, CI shrinks
print("\nEffect of sample size on CI width:")
for n_samples in [5, 10, 30, 100, 1000]:
    scores_n = np.random.normal(0.88, 0.04, n_samples)
    se_n = scores_n.std(ddof=1) / np.sqrt(n_samples)
    t_crit = stats.t.ppf(0.975, df=n_samples-1)
    width = 2 * t_crit * se_n
    print(f"  n={n_samples:5d}: CI width ≈ {width:.4f}")
# CI shrinks proportionally to 1/√n
```

---

## 18. Correlation vs Causation

This is a critical mindset for any ML practitioner.

```
Correlation: two variables tend to move together.
Causation:   changes in one variable CAUSE changes in another.

Examples of spurious correlations:
  - Nicolas Cage movies released per year ↔ drowning deaths  (r=0.67!)
  - Per capita cheese consumption ↔ deaths by bed sheet tangling (r=0.95!)
  - Ice cream sales ↔ shark attacks  (both caused by hot weather)

Why ML models are correlation machines:
  ML models find correlations in training data.
  They do NOT understand causation.
  
  Example: A model learns that "hospital visits correlate with death."
  → Model predicts high mortality for hospital patients.
  → Solution is NOT "don't go to hospital"!
  
  Example: Model trained in winter might learn:
           "wearing heavy coats correlates with being sick."
           In summer, coats are absent → model confused.

  This is distribution shift — training distribution ≠ deployment distribution.

Testing for causation requires:
  - Randomized controlled experiments (A/B tests)
  - Natural experiments (exogenous changes)
  - Causal inference methods (Pearl's do-calculus, instrumental variables)
```

```python
import numpy as np
import matplotlib.pyplot as plt

np.random.seed(42)
n = 1000

# Confounded example: ice cream vs shark attacks (temperature is confounder)
temperature = np.random.normal(25, 5, n)   # temperature in Celsius
ice_cream_sales = temperature * 50 + np.random.randn(n) * 100 + 500
shark_attacks   = temperature * 2  + np.random.randn(n) * 5

print("Correlations:")
print(f"  temp vs ice_cream: {np.corrcoef(temperature, ice_cream_sales)[0,1]:.3f}")
print(f"  temp vs sharks:    {np.corrcoef(temperature, shark_attacks)[0,1]:.3f}")
print(f"  ice_cream vs sharks: {np.corrcoef(ice_cream_sales, shark_attacks)[0,1]:.3f}")

# Controlling for temperature removes the spurious correlation
# Residual analysis: regress out temperature
def residualize(X, confounder):
    """Remove linear effect of confounder from X."""
    slope = np.cov(X, confounder)[0,1] / np.var(confounder)
    return X - slope * confounder

ice_cream_residual = residualize(ice_cream_sales, temperature)
sharks_residual    = residualize(shark_attacks, temperature)
residual_corr = np.corrcoef(ice_cream_residual, sharks_residual)[0,1]
print(f"\nAfter controlling for temperature:")
print(f"  residual corr(ice_cream, sharks) = {residual_corr:.4f}  (near 0!)")
print(f"  The correlation was spurious — both caused by temperature (the confounder).")
```

---

## 19. Summary

```
Probability & Statistics for ML — Complete Reference
────────────────────────────────────────────────────────────────────────
FOUNDATIONS
  P(A|B) = P(A∩B)/P(B)          conditional probability
  Bayes: P(A|B) = P(B|A)P(A)/P(B)    posterior = likelihood × prior / evidence

EXPECTATION & VARIANCE
  E[X] = Σ x P(x)               expected value (weighted average)
  Var(X) = E[(X-μ)²]            spread around mean
  Cov(X,Y) = E[(X-μX)(Y-μY)]    how two RVs move together
  ρ(X,Y) = Cov / (σX σY)        correlation in [-1, +1]

KEY DISTRIBUTIONS
  Gaussian N(μ,σ²) → everywhere due to CLT
  Bernoulli(p) → binary outcomes
  Categorical → softmax outputs (multi-class)
  Uniform → uninformative prior, random initialization

INFORMATION THEORY
  Entropy H(P) = -Σ P log P             uncertainty of a distribution
  Cross-entropy H(P,Q) = -Σ P log Q     divergence between distributions
                                          = classification loss function!
  KL D_KL(P||Q) = Σ P log(P/Q)         VAE regularizer, policy optimization

MLE AND MAP
  MLE: maximize P(data|θ) → minimize -log P(data|θ)
       Linear regression MSE = Gaussian MLE
       Logistic regression BCE = Bernoulli MLE
  MAP: MLE + prior → MLE + regularization
       Gaussian prior → Ridge; Laplace prior → Lasso

TESTING
  p-value: P(seeing this extreme data | H₀ is true)
  CI: range containing true parameter with 95% confidence
  Bootstrap: model-free CIs by resampling

MINDSET
  Correlation ≠ Causation
  ML models find correlations; causation requires experiments
────────────────────────────────────────────────────────────────────────
```

---

## Mini Projects

### Mini Project 1: Naive Bayes Spam Classifier (1.5 hours)

**Goal:** Build a spam filter from scratch using Bayes' theorem applied to word frequencies. No sklearn — pure probability math.

```python
import re
import math
from collections import defaultdict
from typing import List, Tuple

class NaiveBayesSpamFilter:
    """Multinomial Naive Bayes for spam detection."""
    
    def __init__(self, alpha: float = 1.0):
        self.alpha = alpha  # Laplace smoothing
        self.log_prior = {}
        self.log_likelihoods = {}
        self.vocab = set()
    
    def tokenize(self, text: str) -> List[str]:
        return re.findall(r'\b\w+\b', text.lower())
    
    def fit(self, texts: List[str], labels: List[str]):
        """Train on labeled examples."""
        classes = set(labels)
        class_docs = defaultdict(list)
        
        for text, label in zip(texts, labels):
            class_docs[label].append(self.tokenize(text))
        
        n_total = len(texts)
        
        for cls in classes:
            # Prior: P(class)
            self.log_prior[cls] = math.log(len(class_docs[cls]) / n_total)
            
            # Word counts for this class
            word_counts = defaultdict(int)
            total_words = 0
            for doc in class_docs[cls]:
                for word in doc:
                    word_counts[word] += 1
                    total_words += 1
                    self.vocab.add(word)
            
            # Log-likelihood with Laplace smoothing
            vocab_size = len(self.vocab)
            self.log_likelihoods[cls] = {
                word: math.log((count + self.alpha) / (total_words + self.alpha * vocab_size))
                for word, count in word_counts.items()
            }
            # OOV token
            self.log_likelihoods[cls]['<OOV>'] = math.log(
                self.alpha / (total_words + self.alpha * vocab_size)
            )
    
    def predict(self, text: str) -> Tuple[str, dict]:
        """Predict class and return log-posteriors."""
        tokens = self.tokenize(text)
        scores = {}
        
        for cls in self.log_prior:
            score = self.log_prior[cls]
            ll = self.log_likelihoods[cls]
            for word in tokens:
                score += ll.get(word, ll.get('<OOV>', -10))
            scores[cls] = score
        
        predicted = max(scores, key=scores.get)
        # Convert log-scores to probabilities
        max_score = max(scores.values())
        exp_scores = {cls: math.exp(s - max_score) for cls, s in scores.items()}
        total = sum(exp_scores.values())
        probs = {cls: v/total for cls, v in exp_scores.items()}
        
        return predicted, probs


# Sample training data
train_data = [
    ("Get free money now! Click here for amazing offer!", "spam"),
    ("Win a million dollars today! Limited time offer!", "spam"),
    ("Congratulations! You've been selected for a prize!", "spam"),
    ("Buy cheap pills online. Lowest prices guaranteed!", "spam"),
    ("Hey, are we still on for lunch tomorrow?", "ham"),
    ("The meeting has been rescheduled to 3pm Friday", "ham"),
    ("Can you review my pull request when you get a chance?", "ham"),
    ("Here are the project files you requested", "ham"),
    ("Your account statement is ready to view", "ham"),
]

texts, labels = zip(*train_data)
clf = NaiveBayesSpamFilter(alpha=1.0)
clf.fit(list(texts), list(labels))

# Test
test_messages = [
    "FREE! Click to claim your prize money now!",
    "Let me know if Tuesday works for the meeting",
]

for msg in test_messages:
    pred, probs = clf.predict(msg)
    print(f"[{pred.upper()}] ({probs.get('spam', 0):.2%} spam) → {msg[:60]}")
```

---

### Mini Project 2: A/B Test Statistical Significance Calculator (1 hour)

**Goal:** Build an interactive A/B test calculator that uses hypothesis testing to tell you if a change is statistically significant.

```python
import numpy as np
from scipy import stats
import matplotlib.pyplot as plt

def ab_test_proportions(
    control_visitors: int, control_conversions: int,
    treatment_visitors: int, treatment_conversions: int,
    alpha: float = 0.05
) -> dict:
    """Test if two conversion rates are significantly different."""
    
    p_control = control_conversions / control_visitors
    p_treatment = treatment_conversions / treatment_visitors
    
    # Pooled proportion (under null hypothesis)
    p_pool = (control_conversions + treatment_conversions) / (control_visitors + treatment_visitors)
    
    # Standard error
    se = np.sqrt(p_pool * (1 - p_pool) * (1/control_visitors + 1/treatment_visitors))
    
    # Z-statistic
    z = (p_treatment - p_control) / se
    
    # Two-tailed p-value
    p_value = 2 * (1 - stats.norm.cdf(abs(z)))
    
    # Confidence interval for difference
    se_diff = np.sqrt(p_control*(1-p_control)/control_visitors + p_treatment*(1-p_treatment)/treatment_visitors)
    ci_lower = (p_treatment - p_control) - 1.96 * se_diff
    ci_upper = (p_treatment - p_control) + 1.96 * se_diff
    
    return {
        "control_rate": p_control,
        "treatment_rate": p_treatment,
        "relative_lift": (p_treatment - p_control) / p_control,
        "z_stat": z,
        "p_value": p_value,
        "ci_95": (ci_lower, ci_upper),
        "is_significant": p_value < alpha,
        "power": None  # Calculated below
    }

def calculate_sample_size(baseline_rate: float, min_detectable_effect: float, alpha: float = 0.05, power: float = 0.80) -> int:
    """How many visitors do we need per group?"""
    z_alpha = stats.norm.ppf(1 - alpha/2)   # 1.96
    z_beta = stats.norm.ppf(power)           # 0.84
    p2 = baseline_rate + min_detectable_effect
    p_avg = (baseline_rate + p2) / 2
    n = ((z_alpha * np.sqrt(2 * p_avg * (1-p_avg)) +
          z_beta * np.sqrt(baseline_rate*(1-baseline_rate) + p2*(1-p2)))**2
         / min_detectable_effect**2)
    return int(np.ceil(n))

# Example: e-commerce checkout button A/B test
result = ab_test_proportions(
    control_visitors=5000,   control_conversions=250,   # 5.0% rate
    treatment_visitors=5000, treatment_conversions=290,  # 5.8% rate
)

print("=== A/B Test Results ===")
print(f"Control:    {result['control_rate']:.2%}")
print(f"Treatment:  {result['treatment_rate']:.2%}")
print(f"Lift:       {result['relative_lift']:+.1%}")
print(f"Z-stat:     {result['z_stat']:.3f}")
print(f"P-value:    {result['p_value']:.4f}")
print(f"95% CI:     ({result['ci_95'][0]:+.3%}, {result['ci_95'][1]:+.3%})")
print(f"Significant: {'YES ✓' if result['is_significant'] else 'NO ✗'}")

# Sample size calculator
n = calculate_sample_size(baseline_rate=0.05, min_detectable_effect=0.01)
print(f"\nNeed {n:,} visitors per group to detect 1% lift with 80% power")
```

---

### Mini Project 3: Distribution Fitter (1 hour)

**Goal:** Given real data, automatically fit multiple distributions and pick the best fit using the Kolmogorov-Smirnov test.

```python
import numpy as np
import matplotlib.pyplot as plt
from scipy import stats

def fit_distributions(data: np.ndarray) -> list:
    """Fit multiple distributions to data, rank by K-S test."""
    candidates = [
        stats.norm, stats.lognorm, stats.expon, 
        stats.gamma, stats.beta, stats.weibull_min,
    ]
    
    results = []
    for dist in candidates:
        try:
            params = dist.fit(data)
            ks_stat, p_value = stats.kstest(data, dist.cdf, args=params)
            aic = 2 * len(params) - 2 * np.sum(dist.logpdf(data, *params))
            results.append({
                "name": dist.name,
                "params": params,
                "ks_stat": ks_stat,
                "p_value": p_value,
                "aic": aic,
            })
        except Exception:
            pass
    
    return sorted(results, key=lambda x: x["ks_stat"])

# Test with log-normal data (typical for income, stock prices)
data = np.random.lognormal(mean=2.0, sigma=0.5, size=1000)

fits = fit_distributions(data)
print("Distribution Fitting Results (ranked by K-S statistic):")
print(f"{'Distribution':<15} {'K-S stat':>10} {'p-value':>10} {'AIC':>12}")
print("-" * 50)
for r in fits[:5]:
    print(f"{r['name']:<15} {r['ks_stat']:>10.4f} {r['p_value']:>10.4f} {r['aic']:>12.1f}")

# Plot best fit vs data
best = fits[0]
dist = getattr(stats, best['name'])
params = best['params']

fig, axes = plt.subplots(1, 2, figsize=(12, 5))

# Histogram + PDF
x_range = np.linspace(data.min(), data.max(), 200)
axes[0].hist(data, bins=50, density=True, alpha=0.5, label="Data")
axes[0].plot(x_range, dist.pdf(x_range, *params), 'r-', lw=2,
             label=f"Best fit: {best['name']}")
axes[0].set_title(f"Best Fit: {best['name']} (K-S={best['ks_stat']:.4f})")
axes[0].legend()

# QQ-plot
theoretical_quantiles = dist.ppf(np.linspace(0.01, 0.99, 100), *params)
empirical_quantiles = np.percentile(data, np.linspace(1, 99, 100))
axes[1].scatter(theoretical_quantiles, empirical_quantiles, s=10)
lims = [min(theoretical_quantiles.min(), empirical_quantiles.min()),
        max(theoretical_quantiles.max(), empirical_quantiles.max())]
axes[1].plot(lims, lims, 'r--', label="Perfect fit line")
axes[1].set_xlabel("Theoretical quantiles"), axes[1].set_ylabel("Empirical quantiles")
axes[1].set_title("Q-Q Plot"), axes[1].legend()

plt.tight_layout()
plt.savefig("distribution_fit.png")
plt.show()
```

---

## 20. Exercises

**Exercise 1: Bayesian Update Chain**
A factory produces items where 2% are defective. Three quality tests are applied sequentially:
- Test 1: 90% sensitive, 95% specific
- Test 2: 85% sensitive, 92% specific
- Test 3: 95% sensitive, 99% specific

Using Bayes' theorem, compute the posterior probability of defect after each test comes back positive. Plot the probability as a function of the number of positive tests.

*Hint: After each test, the posterior becomes the prior for the next test.*

**Exercise 2: Distribution Fitting**
Generate a mixture of two Gaussians: `X = 0.4 * N(2, 1) + 0.6 * N(7, 1.5)`. Fit this using the EM algorithm (a classic ML algorithm — look it up!): iteratively compute "responsibilities" (which Gaussian each point belongs to) and update the mixture parameters. Plot the true distribution vs your fitted distribution after each EM iteration.

*Hint: Use `scipy.stats.norm.pdf()` for the E-step. Update means, variances, weights in M-step.*

**Exercise 3: MLE for Custom Distribution**
Derive analytically (on paper) the MLE estimator for the Poisson distribution parameter λ. Then verify numerically: generate 1000 Poisson(λ=3.5) samples and show that `λ̂_MLE = mean(data)`. Use both `scipy.optimize.minimize` and your analytical formula.

*Hint: Poisson PMF: P(X=k) = e^(-λ)λ^k/k!. Log-likelihood: l(λ) = -nλ + log(λ)Σxᵢ - Σlog(xᵢ!). Set dl/dλ = 0.*

**Exercise 4: Bootstrap vs t-test**
Run a simulation study comparing bootstrap CI to t-interval CI. For each of 1000 experiments: draw n=20 samples from a skewed distribution (Exponential(1)), compute the 95% CI using both methods, and check if the true mean (=1.0) is inside. Report the "coverage probability" for each method. Which is more accurate?

*Hint: Coverage probability should be close to 95% for a correct CI. Bootstrap is more accurate for non-normal data.*

**Exercise 5: Calibration Curve**
A classifier outputs probability scores for the positive class. The model is "calibrated" if when it predicts P(positive) = 0.7, about 70% of those examples are truly positive. Generate 1000 predictions with known calibration errors (e.g., overconfident: multiply logits by 2 before sigmoid). Plot the calibration curve (reliability diagram): x = predicted probability, y = actual fraction positive (binned). Then apply Platt scaling (fit a logistic regression on the scores) to correct calibration.

*Hint: Use `sklearn.calibration.calibration_curve()` or implement it: bin predictions into 10 bins, compute actual positive rate per bin.*

---

**What's Next →** Chapter 08: Your First ML Model — Linear Regression

*You've completed the seven foundational chapters. You now have the tools (Python, NumPy, Pandas, visualization), the mathematics (linear algebra, calculus, probability), and the conceptual framework to understand machine learning deeply. The next section begins with your first real model: linear regression — where everything you've learned in these seven chapters comes together in one beautiful algorithm.*
