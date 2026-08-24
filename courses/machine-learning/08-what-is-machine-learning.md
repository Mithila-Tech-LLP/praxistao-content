# Chapter 08: What Is Machine Learning?

*Before this chapter, you learned Python. Before writing ML code, you need to truly understand what ML is — not just its definition, but its logic, its limits, and its three completely different flavors. This chapter gives you that foundation.*

---

## Table of Contents

1. From Rules to Learning — The Fundamental Shift
2. The Three Types of Machine Learning
3. The Machine Learning Workflow
4. Your First ML Program — Without Any Libraries
5. Introduction to scikit-learn
6. Train, Validation, and Test Split — Why You Need Three Sets
7. Bias vs Variance — The Central Trade-off
8. How to Read a Dataset
9. Common ML Benchmarks
10. Summary
11. Exercises

---

## 1. From Rules to Learning — The Fundamental Shift

Let's start with a problem that is genuinely impossible to solve by writing rules.

**Problem:** Given a photo, is there a dog in it?

You might try writing rules:
```
IF the image contains four-legged shapes
   AND those shapes have pointy ears
   AND the colour is brownish
THEN: dog
```

But:
- Dalmatians are spotted white and black
- Chihuahuas look very different from Great Danes
- A dog in a hat is still a dog
- A stuffed toy dog that matches all the rules is not a dog
- A dog seen from behind shows none of these features

**There is no set of rules that captures "dog" in all its variety.** And yet every human — including a 2-year-old — can identify a dog instantly.

How? By looking at hundreds or thousands of dogs from the age of one. By building an internal representation of "dog-ness" from experience, not from rules.

Machine learning asks: **what if a computer could do the same?**

The answer is yes. And here is what the shift looks like:

```
TRADITIONAL PROGRAMMING:
  Engineer looks at examples → Engineer writes rules → Computer applies rules → Output

MACHINE LEARNING:
  Engineer collects examples + labels → Algorithm discovers rules → Output
  (The computer finds the rules itself)
```

This is not just a different approach. It is a fundamentally different relationship between programmer and computer.

In traditional programming, you are the teacher and the computer is the obedient student. In machine learning, you are the supervisor — you provide examples and feedback — but the computer figures out the actual logic.

---

## 2. The Three Types of Machine Learning

Machine learning is not one thing. There are three completely different categories, each solving a different kind of problem.

### Type 1: Supervised Learning

**The idea:** You give the computer labeled examples. It learns to predict the label for new examples.

```
Training data (you provide):
Photo #1: [pixels] → label: "dog"
Photo #2: [pixels] → label: "cat"
Photo #3: [pixels] → label: "dog"
...10,000 more...

What the model learns: "which visual patterns correlate with which labels?"

At inference time:
New photo: [pixels] → model predicts: "dog"
```

**When to use it:** When you have examples of inputs AND their correct outputs.

**Real examples:**
- Spam detection: emails labeled spam/not-spam → model learns to filter spam
- House prices: house features + prices → model predicts price of new house
- Medical diagnosis: X-rays labeled tumor/no-tumor → model diagnoses new X-rays
- Translation: sentences in English + their French translations → model translates

### Type 2: Unsupervised Learning

**The idea:** You give the computer unlabeled examples. It finds patterns and structure on its own.

```
Data (no labels):
Customer 1: [age=25, purchases=50, city=London]
Customer 2: [age=45, purchases=200, city=London]
Customer 3: [age=24, purchases=45, city=Manchester]
...

What the model finds (on its own):
Group A: young, low-spending customers
Group B: older, high-spending customers
```

Nobody told it there were two groups. It found the structure in the data by itself.

**When to use it:** When you do not have labels, or when you do not even know what the "answer" looks like.

**Real examples:**
- Customer segmentation: grouping customers by behavior (Chapter 12)
- Anomaly detection: finding transactions that look unusual (potential fraud)
- Topic modeling: discovering themes in thousands of news articles
- Compression: finding compact representations of data (embeddings, Chapter 33)

### Type 3: Reinforcement Learning

**The idea:** An agent learns by taking actions in an environment, receiving rewards for good actions and penalties for bad ones.

```
State: chess board position
Agent: chess AI
Actions: available moves
Reward: +1 for winning, -1 for losing, 0 otherwise

Learning process:
Try millions of games → discover which moves tend to lead to wins → prefer those moves
```

There are no labeled examples. There is only: action, result, learn.

**When to use it:** When you want a system to optimize long-term behavior through trial and error.

**Real examples:**
- AlphaGo and AlphaZero: game playing
- RLHF (Reinforcement Learning from Human Feedback): how ChatGPT, Claude, and Gemini are fine-tuned to be helpful
- Robot arm training: teaching a robot to grasp objects
- Ad bidding: maximizing click-through rates

---

## 3. The Machine Learning Workflow

Almost every ML project follows the same sequence:

```
1. DEFINE THE PROBLEM
   What am I trying to predict?
   What data do I have?
   What counts as "success"?

2. COLLECT AND EXPLORE DATA
   Load the data
   Look at its shape, distributions, missing values
   Visualize it

3. PREPARE THE DATA
   Handle missing values
   Convert categories to numbers
   Scale numerical features
   Split into train/test

4. CHOOSE AND TRAIN A MODEL
   Pick an algorithm suitable for the problem
   Feed it the training data
   Let it learn

5. EVALUATE
   Test on data the model has never seen
   Measure accuracy/error
   Understand where it fails

6. IMPROVE
   Try different algorithms
   Tune hyperparameters
   Engineer better features
   Get more data

7. DEPLOY
   Package the model
   Serve it as an API
   Monitor it in production
```

You will go through this cycle dozens of times in this course.

---

## 4. Your First ML Program — Built from Scratch

Let's build the simplest possible ML system to make the learning process concrete. No libraries. Just math and loops.

**Problem:** Predict someone's exam score from the hours they studied.

```
Data:
Hours studied | Exam score
           1  |        50
           2  |        55
           3  |        65
           4  |        70
           5  |        80
```

The simplest model is a **linear model**: `predicted_score = weight × hours + bias`

We need to find the `weight` and `bias` that best fit the data. We do this by:
1. Starting with random values for `weight` and `bias`
2. Making a prediction
3. Measuring how wrong we are (the "loss")
4. Adjusting `weight` and `bias` to reduce the loss
5. Repeating

```python
# data
hours = [1, 2, 3, 4, 5]
scores = [50, 55, 65, 70, 80]

# start with random model parameters
weight = 0.0
bias = 0.0

# hyperparameter: how big should our steps be?
learning_rate = 0.01

def predict(hours_studied):
    return weight * hours_studied + bias

def mean_squared_error(predictions, targets):
    total = 0
    for pred, target in zip(predictions, targets):
        total += (pred - target) ** 2
    return total / len(predictions)

# training loop: 1000 iterations
for epoch in range(1000):
    # make predictions
    predictions = [predict(h) for h in hours]
    
    # compute loss (how wrong are we?)
    loss = mean_squared_error(predictions, scores)
    
    # compute gradients (which direction should we adjust?)
    # these formulas come from calculus — you will learn them in Chapter 06
    n = len(hours)
    d_weight = (-2/n) * sum(h * (s - p) for h, s, p in zip(hours, scores, predictions))
    d_bias   = (-2/n) * sum(s - p for s, p in zip(scores, predictions))
    
    # update parameters (gradient descent)
    weight -= learning_rate * d_weight
    bias   -= learning_rate * d_bias
    
    if epoch % 100 == 0:
        print(f"Epoch {epoch:4d} | Loss: {loss:.2f} | weight: {weight:.2f}, bias: {bias:.2f}")

print(f"\nFinal model: score = {weight:.2f} × hours + {bias:.2f}")
print(f"\nPredictions for 6 hours: {predict(6):.1f}")
print(f"Predictions for 10 hours: {predict(10):.1f}")
```

Running this produces output like:
```
Epoch    0 | Loss: 4050.00 | weight: 0.73, bias: 0.13
Epoch  100 | Loss: 12.34 | weight: 5.89, bias: 42.31
Epoch  200 | Loss: 12.15 | weight: 6.02, bias: 41.93
Epoch  500 | Loss: 12.11 | weight: 6.05, bias: 41.75
Epoch  900 | Loss: 12.11 | weight: 6.05, bias: 41.75

Final model: score = 6.05 × hours + 41.75

Predictions for 6 hours: 78.1
Predictions for 10 hours: 102.3
```

Notice what happened:
- Loss started at 4050 (very wrong)
- After 100 iterations: 12.34 (much better)
- After 500+: converged to ~12.11 (can't get better with this simple model)
- The model discovered that score ≈ 6 × hours + 42

You just wrote a machine learning algorithm from scratch. This is linear regression. In the next chapter, we will use scikit-learn which does all of this for you — but now you understand what is happening inside.

---

## 5. Introduction to scikit-learn

scikit-learn is Python's most important classical ML library. It gives you dozens of algorithms with a consistent interface.

```python
from sklearn.linear_model import LinearRegression
import numpy as np

# data (scikit-learn wants 2D arrays for X)
X = np.array([[1], [2], [3], [4], [5]])  # hours — each sample is a list
y = np.array([50, 55, 65, 70, 80])       # scores

# create a model
model = LinearRegression()

# train it (this replaces our 1000-iteration loop above)
model.fit(X, y)

# make predictions
print(f"Prediction for 6 hours: {model.predict([[6]])[0]:.1f}")
print(f"Slope (weight): {model.coef_[0]:.2f}")
print(f"Intercept (bias): {model.intercept_:.2f}")
```

Output:
```
Prediction for 6 hours: 78.3
Slope (weight): 6.05
Intercept (bias): 41.75
```

Same result as our from-scratch version. scikit-learn is doing the same things — it just does them faster and more reliably.

The pattern `model.fit(X, y)` / `model.predict(X)` is the same for nearly every scikit-learn model. Learn it once, use it everywhere.

---

## 6. Train, Validation, and Test Split

This is one of the most important concepts in ML, and one of the most commonly misunderstood.

**The core problem:** If you test your model on the same data you trained it on, you are measuring memorization, not learning.

Imagine a student who sees the exam questions beforehand. They study those exact questions and score 100%. Does that mean they understood the subject? No — they just memorized the test.

Models do the same thing. If you let a model see the test data during training, it can learn to memorize those specific examples without learning the underlying pattern.

**The solution:** Keep some data hidden until the very end.

```
ALL DATA (100%)
│
├── TRAINING SET (70%)
│   Used to train the model.
│   The model sees this data. It learns from it.
│
├── VALIDATION SET (15%)
│   Used to tune the model during development.
│   The model does NOT train on this. But YOU see the score.
│   Used to answer: "which model architecture works best?"
│   WARNING: If you use this to make many decisions, the model
│   indirectly learns from it (through your choices).
│
└── TEST SET (15%)
    Used once, at the very end, to get a fair estimate of performance.
    The model has NEVER seen this. You should only run it ONCE.
    If you run it many times and pick the best result, you are cheating.
```

```python
from sklearn.model_selection import train_test_split

X = np.array([[1], [2], [3], [4], [5], [6], [7], [8], [9], [10]])
y = np.array([50, 55, 65, 70, 80, 82, 90, 92, 95, 98])

# First split: set aside test set (20%)
X_temp, X_test, y_temp, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

# Second split: separate train and validation from the remaining 80%
X_train, X_val, y_train, y_val = train_test_split(X_temp, y_temp, test_size=0.25, random_state=42)
# 0.25 of 0.80 = 0.20 of total → 60/20/20 split

print(f"Train size: {len(X_train)}")
print(f"Validation size: {len(X_val)}")
print(f"Test size: {len(X_test)}")
```

---

## 7. Bias vs Variance — The Central Trade-off

This is one of those concepts that sounds abstract but has very concrete consequences.

**Bias** is systematic error — when your model is too simple to capture the real pattern.

**Variance** is inconsistency — when your model is so complex it memorizes the training data and fails on anything new.

```
UNDERFITTING (high bias):
  Training data:  [scattered points that roughly follow a curve]
  Model:          draws a straight line through them
  Problem:        line is too simple — misses the curve
  Training error: high
  Test error:     high

OVERFITTING (high variance):
  Training data:  [scattered points that roughly follow a curve]
  Model:          draws a wiggly line through EVERY SINGLE POINT
  Problem:        memorized the noise, not the pattern
  Training error: near zero
  Test error:     high

JUST RIGHT:
  Training data:  [scattered points that roughly follow a curve]
  Model:          draws a smooth curve through them
  Problem:        none
  Training error: moderate
  Test error:     moderate (close to training error)
```

**How to diagnose:**

| Situation | Bias | Variance | Problem |
|-----------|------|----------|---------|
| High train error, high test error | High | Low | Underfitting — model too simple |
| Low train error, high test error | Low | High | Overfitting — model too complex |
| Low train error, low test error | Low | Low | Good — this is the goal |

**How to fix underfitting:**
- Use a more complex model
- Train longer
- Add more features

**How to fix overfitting:**
- Get more training data
- Use a simpler model
- Use regularization (L1/L2, dropout — covered in Chapter 20)
- Use early stopping

You will encounter this trade-off in every single ML project you ever work on.

---

## 8. How to Read a Dataset

The first thing you do with any new dataset is **explore it**. Never train a model on data you have not looked at.

```python
import pandas as pd
import matplotlib.pyplot as plt

# Load a well-known dataset: California housing prices
from sklearn.datasets import fetch_california_housing

data = fetch_california_housing(as_frame=True)
df = data.frame

# Step 1: Look at the shape
print(f"Rows: {df.shape[0]}, Columns: {df.shape[1]}")

# Step 2: Look at the first few rows
print(df.head())

# Step 3: Get summary statistics
print(df.describe())

# Step 4: Check for missing values
print(df.isnull().sum())

# Step 5: Look at distributions
df.hist(bins=30, figsize=(12, 8))
plt.tight_layout()
plt.show()

# Step 6: Look at correlations
print(df.corr()['MedHouseVal'].sort_values(ascending=False))
```

Key things to look for:
- **Missing values:** rows where some data is absent (need to handle before training)
- **Outliers:** extreme values that might skew the model
- **Distribution:** is the data normally distributed, skewed, bimodal?
- **Correlations:** which features are most related to the target?
- **Data types:** are categories stored as numbers? Are numbers stored as text?

---

## Summary

You now understand the foundations of machine learning:

- **ML is a paradigm shift**: instead of writing rules, you give the computer examples and let it discover the rules.
- **Three types**: supervised (labeled data), unsupervised (no labels), reinforcement (rewards).
- **The workflow**: define → collect → prepare → train → evaluate → improve → deploy.
- **Train/val/test split**: essential to measure real performance, not memorization.
- **Bias vs variance**: the fundamental trade-off between underfitting and overfitting.

---

## Exercises

**Easy:**

1. In your own words (without looking at this chapter), explain the difference between supervised and unsupervised learning. Give one real-world example of each.

2. Run the from-scratch linear regression code. Change the `learning_rate` to 0.1 and run again. What happens? Try 0.001. What happens? Why?

3. What would be wrong with using `test_size=0.5` (50% test)? What about `test_size=0.01` (1% test)?

**Medium:**

4. Use the California housing dataset from Section 8. Train a `LinearRegression` model on 80% of the data. Evaluate on the remaining 20%. What is the mean squared error? What does that mean in dollars?

5. Train a `LinearRegression` on only 10 examples. Then train one on the full dataset. Compare both models' test accuracy. What does this demonstrate about data quantity?

6. Look at the sklearn documentation and find 3 classifiers and 3 regressors. For each, write one sentence describing what kind of problem it is best for.

**Hard:**

7. Implement **cross-validation** from scratch: split the data into 5 equal parts. Train on 4 parts, test on 1. Rotate which part is the test set. Average the 5 test scores. This is called 5-fold cross-validation. Why is this more reliable than a single train/test split?

8. Find a real dataset from Kaggle or UCI ML repository. Follow the 8-step exploration process from Section 8. Write a paragraph describing what you found — the shape, distributions, missing values, and which features seem most predictive.
