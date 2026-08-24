# Chapter 13: Model Evaluation, Feature Engineering, and Pipelines

> **"Feature engineering is the process of using domain knowledge to create features that make machine learning algorithms work. It is fundamental to the application of machine learning."**
> — Pedro Domingos

---

## Table of Contents
1. [The Evaluation Trap](#1-the-evaluation-trap)
2. [Cross-Validation — Deep Dive](#2-cross-validation--deep-dive)
3. [Regression Metrics — Full Reference](#3-regression-metrics--full-reference)
4. [Classification Metrics — Full Reference](#4-classification-metrics--full-reference)
5. [Handling Imbalanced Data](#5-handling-imbalanced-data)
6. [Probability Calibration](#6-probability-calibration)
7. [Hyperparameter Tuning](#7-hyperparameter-tuning)
8. [Feature Engineering — Numerical Features](#8-feature-engineering--numerical-features)
9. [Feature Engineering — Categorical Features](#9-feature-engineering--categorical-features)
10. [Handling Missing Values](#10-handling-missing-values)
11. [Feature Selection](#11-feature-selection)
12. [Dealing with Outliers](#12-dealing-with-outliers)
13. [Sklearn Pipelines — The Right Way](#13-sklearn-pipelines--the-right-way)
14. [ColumnTransformer — Mixed Data Types](#14-columntransformer--mixed-data-types)
15. [Saving and Loading Pipelines](#15-saving-and-loading-pipelines)
16. [Summary](#16-summary)
17. [Exercises](#17-exercises)

---

## 1. The Evaluation Trap

There is a consistent pattern in beginners' ML work: they report a high accuracy, deploy a model, and then it fails in production. The problem is almost always **evaluation methodology**, not the model itself.

### The Many Ways to Get Evaluation Wrong

```
COMMON EVALUATION MISTAKES
────────────────────────────────────────────────────────────────────────
1. WRONG METRIC
   Example: 99.5% accuracy on fraud detection sounds great.
   But if 99.5% of transactions are legitimate, a model that says
   "not fraud" for everything achieves 99.5% accuracy while catching
   zero fraud cases.
   Fix: use precision, recall, F1, AUC-PR for imbalanced data.

2. EVALUATING ON TRAINING DATA
   model.score(X_train, y_train) → inflated, meaningless
   Always evaluate on held-out test data.

3. TEST SET CONTAMINATION
   Fitting preprocessing (scaler, PCA, imputer) on all data including test.
   The model has "seen" test data statistics → inflated performance.
   Fix: fit preprocessors on training data only (use Pipelines).

4. LEAKY FEATURES
   Including features in X that are derived from y (the target).
   Example: predicting if a loan defaults, including "collection_calls"
   as a feature (you only call people who already defaulted!).
   Fix: careful temporal ordering, domain knowledge.

5. SINGLE TRAIN/TEST SPLIT
   One split has high variance → lucky or unlucky split.
   Fix: k-fold cross-validation.

6. IGNORING TEMPORAL STRUCTURE
   For time series: splitting randomly means using future data to
   predict the past! Fix: time-ordered splits.

7. SELECTION BIAS IN HYPERPARAMETER TUNING
   Tuning hyperparameters on the test set → test set contamination.
   Fix: tune on validation set, report test set only once.
```

---

## 2. Cross-Validation — Deep Dive

Cross-validation (CV) is the solution to the single-split variance problem. Instead of one train/test split, CV trains and evaluates the model multiple times on different subsets.

### K-Fold Cross-Validation

```
K-FOLD CROSS-VALIDATION (K=5)
────────────────────────────────────────────────────────────────────────
Full training data (never touch test set during CV)
[  Fold 1  |  Fold 2  |  Fold 3  |  Fold 4  |  Fold 5  ]

Round 1: [  VAL  |  train  |  train  |  train  |  train  ]  score₁
Round 2: [  train  |  VAL  |  train  |  train  |  train  ]  score₂
Round 3: [  train  |  train  |  VAL  |  train  |  train  ]  score₃
Round 4: [  train  |  train  |  train  |  VAL  |  train  ]  score₄
Round 5: [  train  |  train  |  train  |  train  |  VAL  ]  score₅

CV Score = mean(score₁, score₂, score₃, score₄, score₅)
CV Std   = std(score₁, ..., score₅)  ← uncertainty of the estimate

Benefits:
  Every sample is used for validation exactly once
  Every sample is used for training k-1 times
  More reliable estimate than single split
  Standard error ∝ 1/√k  (more folds → more precise)

Typical k: 5 or 10
  k=5: faster, slightly higher variance
  k=10: slower, slightly lower variance
  k=n (LOO): maximum use of data but very slow for large n
```

### Stratified K-Fold

For classification, standard k-fold may produce unbalanced folds. Stratified k-fold preserves the class proportions in each fold:

```
STRATIFIED K-FOLD (critical for imbalanced data)
────────────────────────────────────────────────────────────────────────
Dataset: 100 samples, 90 negative, 10 positive (10% positive rate)

Standard K-Fold (k=5):
  Fold 1 might have: 18 negative, 2 positive = 10% positive ✓
  Fold 2 might have: 20 negative, 0 positive = 0% positive  ✗ (no positives!)

Stratified K-Fold (k=5):
  Every fold guaranteed: ~18 negative, ~2 positive = 10% positive ✓

Result: more reliable CV estimates, especially when n_positive is small.
ALWAYS use stratified CV for classification.
```

### Time Series Cross-Validation

Regular K-fold MUST NOT be used for time series — it allows using future data to predict the past.

```
TIME SERIES CROSS-VALIDATION (Forward Chaining)
────────────────────────────────────────────────────────────────────────
Data ordered by time: [Jan  Feb  Mar  Apr  May  Jun  Jul  Aug]

Round 1: [TRAIN: Jan-Mar  | VAL: Apr     ]
Round 2: [TRAIN: Jan-Apr  | VAL: May     ]
Round 3: [TRAIN: Jan-May  | VAL: Jun     ]
Round 4: [TRAIN: Jan-Jun  | VAL: Jul     ]
Round 5: [TRAIN: Jan-Jul  | VAL: Aug     ]

The validation set is ALWAYS in the FUTURE relative to training.
Training set grows with each fold (expanding window).
Alternative: rolling window (fixed training size, slides forward).

Rule: in time series, NEVER let information from the future
      contaminate the past. Use sklearn.model_selection.TimeSeriesSplit.
```

### Implementation

```python
from sklearn.model_selection import (
    cross_val_score, cross_validate,
    KFold, StratifiedKFold, TimeSeriesSplit,
    LeaveOneOut, GroupKFold
)
from sklearn.ensemble import RandomForestClassifier
from sklearn.datasets import load_breast_cancer
import numpy as np

X, y = load_breast_cancer(return_X_y=True)
model = RandomForestClassifier(n_estimators=100, random_state=42)

# ── Basic K-Fold ──────────────────────────────────────────────
kf = KFold(n_splits=5, shuffle=True, random_state=42)
scores = cross_val_score(model, X, y, cv=kf, scoring='roc_auc', n_jobs=-1)
print(f"K-Fold CV AUC:   {scores.mean():.4f} ± {scores.std():.4f}")

# ── Stratified K-Fold ─────────────────────────────────────────
skf = StratifiedKFold(n_splits=5, shuffle=True, random_state=42)
scores_s = cross_val_score(model, X, y, cv=skf, scoring='roc_auc', n_jobs=-1)
print(f"Stratified CV:   {scores_s.mean():.4f} ± {scores_s.std():.4f}")

# ── Multiple Metrics at Once ──────────────────────────────────
results = cross_validate(
    model, X, y,
    cv=skf,
    scoring=['roc_auc', 'f1', 'precision', 'recall', 'accuracy'],
    return_train_score=True,  # Also compute training scores (detect overfit)
    n_jobs=-1
)

print("\nMultiple metrics:")
for metric in ['accuracy', 'f1', 'precision', 'recall', 'roc_auc']:
    train = results[f'train_{metric}']
    val   = results[f'test_{metric}']
    overfit_gap = train.mean() - val.mean()
    print(f"  {metric:10s}: train={train.mean():.4f}±{train.std():.4f}  "
          f"val={val.mean():.4f}±{val.std():.4f}  "
          f"gap={overfit_gap:.4f}")

# ── Leave-One-Out ─────────────────────────────────────────────
# LOO: k = n_samples (extreme case)
# Very slow for large n but maximum data efficiency for small n
if len(X) < 500:  # Only for small datasets
    loo = LeaveOneOut()
    loo_scores = cross_val_score(model, X, y, cv=loo, scoring='accuracy', n_jobs=-1)
    print(f"\nLOO CV Accuracy: {loo_scores.mean():.4f} ± {loo_scores.std():.4f}")
    # LOO gives essentially unbiased estimate but high variance

# ── Time Series Split ─────────────────────────────────────────
tss = TimeSeriesSplit(n_splits=5)
print("\nTime Series CV split sizes:")
for fold, (train_idx, val_idx) in enumerate(tss.split(X)):
    print(f"  Fold {fold+1}: train={len(train_idx)}, val={len(val_idx)}")
```

---

## 3. Regression Metrics — Full Reference

```python
import numpy as np
from sklearn.metrics import (
    mean_absolute_error, mean_squared_error, r2_score,
    mean_absolute_percentage_error
)

def comprehensive_regression_report(y_true, y_pred, model_name="", n_features=None):
    """
    Complete regression evaluation report.
    """
    n = len(y_true)
    residuals = y_true - y_pred

    mae   = mean_absolute_error(y_true, y_pred)
    mse   = mean_squared_error(y_true, y_pred)
    rmse  = np.sqrt(mse)
    r2    = r2_score(y_true, y_pred)

    # MAPE: undefined for y_true = 0
    nonzero = y_true != 0
    if nonzero.sum() > 0:
        mape = np.mean(np.abs(residuals[nonzero] / y_true[nonzero])) * 100
    else:
        mape = np.nan

    # Adjusted R²: penalizes adding useless features
    # Only useful when comparing models with different number of features
    if n_features is not None and n > n_features + 1:
        adj_r2 = 1 - (1 - r2) * (n - 1) / (n - n_features - 1)
    else:
        adj_r2 = None

    # Residual analysis
    res_mean  = residuals.mean()
    res_std   = residuals.std()
    res_skew  = ((residuals - res_mean)**3).mean() / res_std**3
    max_error = np.abs(residuals).max()

    print(f"\n{'='*50}")
    print(f"Regression Report: {model_name}")
    print(f"{'='*50}")
    print(f"\nPoint Metrics:")
    print(f"  MAE:   {mae:.4f}  (avg absolute error, same units as y)")
    print(f"  MSE:   {mse:.4f}  (penalizes large errors more)")
    print(f"  RMSE:  {rmse:.4f} (same units as y, most common)")
    print(f"  MAPE:  {mape:.2f}% (percentage error, undefined at y=0)")
    print(f"  R²:    {r2:.4f}  (1 = perfect, 0 = no better than mean)")
    if adj_r2 is not None:
        print(f"  Adj R²:{adj_r2:.4f} (penalizes extra features, for model comparison)")
    print(f"\nResidual Statistics:")
    print(f"  Mean:     {res_mean:.4f}  (should be ~0 for unbiased model)")
    print(f"  Std:      {res_std:.4f}")
    print(f"  Skewness: {res_skew:.4f}  (0 = symmetric, |>1| = very skewed)")
    print(f"  Max |err|:{max_error:.4f}")

    return {'mae': mae, 'mse': mse, 'rmse': rmse, 'mape': mape, 'r2': r2}
```

### When to Use Which Regression Metric

```
REGRESSION METRIC SELECTION GUIDE
────────────────────────────────────────────────────────────────────────
MAE (Mean Absolute Error):
  ✓ When outliers are natural/expected (house prices with mansions)
  ✓ When you want the metric in original units
  ✓ When all errors should be penalized equally
  ✗ Not differentiable at 0 (harder to optimize directly)

RMSE (Root Mean Squared Error):
  ✓ Most commonly reported metric
  ✓ Penalizes large errors (often what we want)
  ✓ Differentiable everywhere
  ✗ Sensitive to outliers (they contribute squared error)

MAPE (Mean Absolute Percentage Error):
  ✓ Scale-independent: compare across different targets
  ✓ Interpretable: "10% error on average"
  ✗ Undefined/explodes when true value is 0
  ✗ Asymmetric: 10% overestimate ≠ 10% underestimate
  ✗ Poor for targets near zero

R² (Coefficient of Determination):
  ✓ Interpretable: fraction of variance explained
  ✓ Scale-independent: compare models on different datasets
  ✓ Benchmark: compares to predicting the mean (R²=0)
  ✗ Can be negative! (model worse than constant mean)
  ✗ Always increases when you add features (use Adjusted R²)

Business Example:
  Predicting sales ($0 to $1,000,000):
    R² = 0.82  → 82% of variance explained
    RMSE = $45,000 → typical error of $45,000
    MAPE = 12%    → 12% relative error
    Which to report to management? MAPE or RMSE (intuitive)
    Which to optimize? MSE (differentiable) or MAE (robust)
```

---

## 4. Classification Metrics — Full Reference

```python
import numpy as np
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, f1_score,
    roc_auc_score, average_precision_score, log_loss,
    brier_score_loss, confusion_matrix, classification_report,
    roc_curve, precision_recall_curve
)

def comprehensive_classification_report(y_true, y_pred, y_proba=None, label=""):
    """
    Complete classification evaluation including all relevant metrics.
    y_proba: probability of positive class for binary, or (n, k) for multi-class
    """
    is_binary = len(np.unique(y_true)) == 2

    print(f"\n{'='*55}")
    print(f"Classification Report: {label}")
    print(f"{'='*55}")

    print(f"\nClass distribution:")
    classes, counts = np.unique(y_true, return_counts=True)
    for c, cnt in zip(classes, counts):
        print(f"  Class {c}: {cnt} samples ({cnt/len(y_true):.1%})")

    print(f"\nCore Metrics:")
    print(f"  Accuracy:    {accuracy_score(y_true, y_pred):.4f}")

    if is_binary:
        tn, fp, fn, tp = confusion_matrix(y_true, y_pred).ravel()
        print(f"\nConfusion Matrix:")
        print(f"  True Neg:   {tn:5d}  |  False Pos:  {fp:5d}")
        print(f"  False Neg:  {fn:5d}  |  True Pos:   {tp:5d}")

        prec = precision_score(y_true, y_pred)
        rec  = recall_score(y_true, y_pred)
        f1   = f1_score(y_true, y_pred)
        spec = tn / (tn + fp) if (tn + fp) > 0 else 0

        print(f"\n  Precision:   {prec:.4f}  TP/(TP+FP)")
        print(f"  Recall:      {rec:.4f}  TP/(TP+FN)")
        print(f"  Specificity: {spec:.4f}  TN/(TN+FP)")
        print(f"  F1 Score:    {f1:.4f}  harmonic mean(prec, recall)")

        if y_proba is not None:
            if y_proba.ndim == 2:
                y_proba_pos = y_proba[:, 1]
            else:
                y_proba_pos = y_proba

            auc_roc = roc_auc_score(y_true, y_proba_pos)
            auc_pr  = average_precision_score(y_true, y_proba_pos)
            ll      = log_loss(y_true, y_proba_pos)
            brier   = brier_score_loss(y_true, y_proba_pos)

            print(f"\nProbability-based Metrics:")
            print(f"  AUC-ROC:     {auc_roc:.4f}  (ranking quality, 0.5=random, 1=perfect)")
            print(f"  AUC-PR:      {auc_pr:.4f}  (better for imbalanced)")
            print(f"  Log Loss:    {ll:.4f}  (lower is better, 0=perfect)")
            print(f"  Brier Score: {brier:.4f}  (lower is better, MSE of probabilities)")
    else:
        # Multi-class
        print(f"\n{classification_report(y_true, y_pred)}")
        if y_proba is not None:
            auc_roc = roc_auc_score(y_true, y_proba, multi_class='ovr', average='macro')
            print(f"\nAUC-ROC (macro OvR): {auc_roc:.4f}")

# Multi-class averaging strategies
print("Multi-class F1 averaging:")
print("  'macro':    mean F1 over classes, equal weight per class")
print("              → sensitive to performance on rare classes")
print("  'micro':    global TP,FP,FN then compute F1")
print("              → dominated by frequent classes")
print("  'weighted': macro but weight by class support (sample count)")
print("              → appropriate when class imbalance is expected")
```

---

## 5. Handling Imbalanced Data

Class imbalance is extremely common in real problems: fraud detection (0.1% fraud), medical diagnosis (5% positive), spam detection (10% spam). Standard accuracy is useless; the model must be designed for the imbalance.

### Strategy 1: Class Weights

```python
from sklearn.linear_model import LogisticRegression
from sklearn.ensemble import RandomForestClassifier
from sklearn.utils.class_weight import compute_class_weight
import numpy as np

# Method 1: class_weight='balanced' (automatic)
# Weights: n_samples / (n_classes × np.bincount(y))
lr = LogisticRegression(class_weight='balanced', max_iter=1000)
rf = RandomForestClassifier(class_weight='balanced', n_estimators=100)

# Method 2: Manual class weights
y = np.array([0]*900 + [1]*100)  # 10% positive
weights = compute_class_weight('balanced', classes=np.unique(y), y=y)
class_weight_dict = dict(zip(np.unique(y), weights))
print(f"Auto class weights: {class_weight_dict}")
# Result: {0: 0.556, 1: 5.0}  (minority class gets 9x weight)

lr_manual = LogisticRegression(
    class_weight=class_weight_dict,
    max_iter=1000
)
```

### Strategy 2: SMOTE (Synthetic Minority Over-sampling)

```python
# pip install imbalanced-learn
try:
    from imblearn.over_sampling import SMOTE, ADASYN
    from imblearn.under_sampling import RandomUnderSampler
    from imblearn.pipeline import Pipeline as ImbPipeline
    from imblearn.combine import SMOTETomek

    # SMOTE: creates synthetic samples in minority class
    # For each minority sample, find k nearest minority neighbors,
    # create new samples along the line between them
    smote = SMOTE(
        sampling_strategy=0.5,  # After resampling, minority/majority = 0.5
        k_neighbors=5,
        random_state=42
    )

    # ADASYN: like SMOTE but generates more near the decision boundary
    adasyn = ADASYN(random_state=42)

    # Combine SMOTE + Tomek Links (removes borderline majority samples)
    smote_tomek = SMOTETomek(random_state=42)

    # IMPORTANT: Apply SMOTE only to training data, never test data
    # Use imblearn's Pipeline to avoid leakage
    from sklearn.preprocessing import StandardScaler
    from sklearn.datasets import make_classification

    X, y = make_classification(
        n_samples=1000, n_features=10,
        weights=[0.9, 0.1],  # 10% minority
        random_state=42
    )

    from sklearn.model_selection import train_test_split
    X_train, X_test, y_train, y_test = train_test_split(X, y, test_size=0.2, random_state=42)

    print(f"Before SMOTE: {np.bincount(y_train)}")
    X_resampled, y_resampled = smote.fit_resample(X_train, y_train)
    print(f"After SMOTE:  {np.bincount(y_resampled)}")

    # imblearn Pipeline (supports resampling steps)
    from sklearn.linear_model import LogisticRegression
    pipe = ImbPipeline([
        ('scaler', StandardScaler()),
        ('smote',  SMOTE(random_state=42)),   # Only applied during fit!
        ('model',  LogisticRegression(max_iter=1000))
    ])
    pipe.fit(X_train, y_train)
    print(f"\nSMOTE Pipeline predictions: {pipe.predict(X_test).sum()} positives")

except ImportError:
    print("imbalanced-learn not installed. pip install imbalanced-learn")
```

### Strategy 3: Threshold Optimization

```python
import numpy as np
from sklearn.metrics import f1_score, precision_score, recall_score, fbeta_score

def find_optimal_threshold(y_true, y_proba, metric='f1', beta=1.0):
    """
    Find the probability threshold that maximizes a given metric.
    
    For recall-heavy problems (medical), use metric='recall' or beta > 1.
    For precision-heavy problems (spam), use metric='precision' or beta < 1.
    F_beta: beta > 1 → recall more important; beta < 1 → precision more important.
    """
    thresholds = np.arange(0.01, 0.99, 0.01)
    best_score = -1
    best_threshold = 0.5

    for t in thresholds:
        y_pred = (y_proba >= t).astype(int)

        if metric == 'f1':
            score = f1_score(y_true, y_pred, zero_division=0)
        elif metric == 'fbeta':
            score = fbeta_score(y_true, y_pred, beta=beta, zero_division=0)
        elif metric == 'precision':
            score = precision_score(y_true, y_pred, zero_division=0)
        elif metric == 'recall':
            score = recall_score(y_true, y_pred, zero_division=0)

        if score > best_score:
            best_score = score
            best_threshold = t

    return best_threshold, best_score

# Example usage
# best_t, best_f1 = find_optimal_threshold(y_test, y_proba, metric='f1')
# y_pred_optimal = (y_proba >= best_t).astype(int)
```

---

## 6. Probability Calibration

A model's raw probability output should match empirical frequencies. If your model says "90% chance of fraud" for 100 transactions, about 90 of them should actually be fraud.

```python
from sklearn.calibration import CalibratedClassifierCV, CalibrationDisplay
from sklearn.ensemble import RandomForestClassifier
from sklearn.linear_model import LogisticRegression
from sklearn.svm import SVC

# RandomForest probabilities are often poorly calibrated
# (tend to be pushed toward 0 and 1 — overconfident)

rf = RandomForestClassifier(n_estimators=100, random_state=42)
rf.fit(X_train, y_train)
y_proba_raw = rf.predict_proba(X_test)[:, 1]

# Platt scaling (sigmoid calibration)
rf_platt = CalibratedClassifierCV(
    RandomForestClassifier(n_estimators=100, random_state=42),
    method='sigmoid',    # Platt scaling (logistic fit)
    cv=5
)
rf_platt.fit(X_train, y_train)
y_proba_platt = rf_platt.predict_proba(X_test)[:, 1]

# Isotonic regression calibration (more flexible but needs more data)
rf_isotonic = CalibratedClassifierCV(
    RandomForestClassifier(n_estimators=100, random_state=42),
    method='isotonic',
    cv=5
)
rf_isotonic.fit(X_train, y_train)
y_proba_iso = rf_isotonic.predict_proba(X_test)[:, 1]

# When calibration matters most:
# 1. Decision threshold is not 0.5 (you need actual probabilities)
# 2. Combining multiple model outputs
# 3. Expected Value calculations (cost-benefit analysis)
# Logistic Regression is usually well-calibrated without further calibration.
# SVMs and trees often need calibration.
```

---

## 7. Hyperparameter Tuning

### Grid Search

```python
from sklearn.model_selection import GridSearchCV, RandomizedSearchCV
from sklearn.ensemble import RandomForestClassifier
from sklearn.datasets import load_breast_cancer
from sklearn.model_selection import StratifiedKFold
import numpy as np
import time

X, y = load_breast_cancer(return_X_y=True)

# ── Grid Search: exhaustive ───────────────────────────────────
# Total combinations = product of all param list lengths
# 3 × 3 × 2 × 2 = 36 combinations × 5 folds = 180 fits
param_grid = {
    'n_estimators':   [100, 200, 300],
    'max_depth':      [None, 5, 10],
    'max_features':   ['sqrt', 'log2'],
    'min_samples_leaf': [1, 5],
}

grid_search = GridSearchCV(
    RandomForestClassifier(random_state=42),
    param_grid,
    cv=StratifiedKFold(n_splits=5, shuffle=True, random_state=42),
    scoring='roc_auc',
    n_jobs=-1,
    verbose=0,
    return_train_score=True
)

start = time.time()
grid_search.fit(X, y)
elapsed = time.time() - start

print(f"Grid Search: {len(grid_search.cv_results_['params'])} combinations "
      f"in {elapsed:.1f}s")
print(f"Best params: {grid_search.best_params_}")
print(f"Best CV AUC: {grid_search.best_score_:.4f}")
```

### Random Search (Often Better Than Grid)

```python
from scipy.stats import randint, uniform

# Random Search: sample from distributions instead of fixed grid
# With n_iter=50 and 5-fold CV: 250 fits
# Often finds better solutions than Grid Search with same budget
param_distributions = {
    'n_estimators':      randint(50, 500),    # Uniform integer
    'max_depth':         [None, 3, 5, 7, 10, 15],
    'max_features':      uniform(0.2, 0.8),   # Continuous [0.2, 1.0]
    'min_samples_split': randint(2, 20),
    'min_samples_leaf':  randint(1, 10),
    'bootstrap':         [True, False],
}

random_search = RandomizedSearchCV(
    RandomForestClassifier(random_state=42),
    param_distributions,
    n_iter=50,           # Number of random combinations to try
    cv=StratifiedKFold(n_splits=5, shuffle=True, random_state=42),
    scoring='roc_auc',
    n_jobs=-1,
    random_state=42,
    verbose=0
)

start = time.time()
random_search.fit(X, y)
elapsed = time.time() - start

print(f"Random Search: 50 combinations in {elapsed:.1f}s")
print(f"Best params: {random_search.best_params_}")
print(f"Best CV AUC: {random_search.best_score_:.4f}")

# Access full results
results_df = pd.DataFrame(random_search.cv_results_)
results_df = results_df.sort_values('rank_test_score')
print("\nTop 5 parameter combinations:")
top5_cols = ['rank_test_score', 'mean_test_score', 'std_test_score',
             'param_n_estimators', 'param_max_depth', 'param_max_features']
print(results_df[top5_cols].head(5).to_string(index=False))
```

### Bayesian Optimization with Optuna

```python
# pip install optuna
try:
    import optuna
    optuna.logging.set_verbosity(optuna.logging.WARNING)

    def objective(trial):
        """Objective function for Optuna to optimize."""
        # Suggest hyperparameters from distributions
        n_estimators = trial.suggest_int('n_estimators', 50, 500)
        max_depth    = trial.suggest_categorical('max_depth', [None, 3, 5, 7, 10])
        max_features = trial.suggest_float('max_features', 0.2, 1.0)
        min_samples_leaf = trial.suggest_int('min_samples_leaf', 1, 10)

        model = RandomForestClassifier(
            n_estimators=n_estimators,
            max_depth=max_depth,
            max_features=max_features,
            min_samples_leaf=min_samples_leaf,
            random_state=42
        )

        cv = StratifiedKFold(n_splits=5, shuffle=True, random_state=42)
        scores = cross_val_score(model, X, y, cv=cv, scoring='roc_auc', n_jobs=-1)
        return scores.mean()  # Optuna maximizes this by default

    study = optuna.create_study(direction='maximize')
    study.optimize(objective, n_trials=50, timeout=60)

    print(f"\nOptuna Best AUC: {study.best_value:.4f}")
    print(f"Best params: {study.best_params}")

except ImportError:
    print("Optuna not installed. pip install optuna")
    print("Bayesian optimization intelligently explores parameter space.")
```

---

## 8. Feature Engineering — Numerical Features

Feature engineering is often the highest-leverage activity in ML. Better features → better models, regardless of which model you use.

### Scaling

```python
from sklearn.preprocessing import (
    StandardScaler, MinMaxScaler, RobustScaler, MaxAbsScaler,
    QuantileTransformer, PowerTransformer
)
import numpy as np

# Generate data with different characteristics
np.random.seed(42)
X_demo = np.array([
    [100, 0.01, 50000],   # Feature 1: income, Feature 2: rate, Feature 3: salary
    [200, 0.05, 75000],
    [150, 0.02, 60000],
    [800, 0.10, 120000],  # Outlier
    [120, 0.03, 55000],
])

scalers = {
    'StandardScaler':    StandardScaler(),      # z = (x - mean) / std
    'MinMaxScaler':      MinMaxScaler(),         # z = (x - min) / (max - min) → [0, 1]
    'RobustScaler':      RobustScaler(),         # z = (x - median) / IQR → robust to outliers
    'MaxAbsScaler':      MaxAbsScaler(),         # z = x / max(|x|) → [-1, 1], good for sparse
}

for name, scaler in scalers.items():
    X_scaled = scaler.fit_transform(X_demo)
    print(f"\n{name}:")
    print(f"  Mean: {X_scaled.mean(axis=0).round(3)}")
    print(f"  Std:  {X_scaled.std(axis=0).round(3)}")
    print(f"  Min:  {X_scaled.min(axis=0).round(3)}")
    print(f"  Max:  {X_scaled.max(axis=0).round(3)}")

print("""
SCALER SELECTION GUIDE:
────────────────────────────────────────────────────────────────────────
StandardScaler:  Most common. Assumes roughly normal distribution.
                 Center to 0, scale to unit variance.

MinMaxScaler:    When you need bounded output [0,1] (image pixels, neural nets).
                 Very sensitive to outliers (outlier → everything compressed).

RobustScaler:    Use when you have outliers and can't remove them.
                 Uses median and IQR → outliers don't distort scaling.

MaxAbsScaler:    For sparse data (TF-IDF). Doesn't shift center (preserves zeros).

QuantileTransformer: Maps to uniform or normal distribution.
                     Useful when model assumes normality (e.g., linear models).
                     Nonlinear transformation → loses some interpretability.
""")
```

### Log Transform

```python
import numpy as np
import pandas as pd

# Right-skewed data is very common: income, prices, population
np.random.seed(42)
income = np.random.lognormal(mean=10.5, sigma=1.2, size=1000)
# This looks like: mostly $10k-100k with a long tail of millionaires

print(f"Income - mean: ${income.mean():,.0f}, std: ${income.std():,.0f}")
print(f"Income - skewness: {pd.Series(income).skew():.2f}")

# Log transform
log_income = np.log1p(income)  # log(1 + x) handles x=0 safely
print(f"\nLog(income) - mean: {log_income.mean():.3f}, std: {log_income.std():.3f}")
print(f"Log(income) - skewness: {pd.Series(log_income).skew():.2f}")  # Much closer to 0!

# sqrt transform: lighter than log, good for count data
sqrt_income = np.sqrt(income)
print(f"\nSqrt(income) skewness: {pd.Series(sqrt_income).skew():.2f}")

# Box-Cox transform: finds optimal power transformation
from sklearn.preprocessing import PowerTransformer
pt = PowerTransformer(method='box-cox')  # Requires positive values
income_boxcox = pt.fit_transform(income.reshape(-1, 1))
print(f"BoxCox lambda: {pt.lambdas_[0]:.3f}")  # Optimal power
```

### Polynomial and Interaction Features

```python
from sklearn.preprocessing import PolynomialFeatures
import numpy as np

X = np.array([[2, 3], [4, 5]])  # 2 features: x₁=2, x₂=3

# Degree 2 polynomial features
poly = PolynomialFeatures(degree=2, include_bias=False)
X_poly = poly.fit_transform(X)

print("Original features:    x₁, x₂")
print("Polynomial features:  x₁, x₂, x₁², x₁x₂, x₂²")
print(f"Input shape:  {X.shape}")
print(f"Output shape: {X_poly.shape}")
print(f"Feature names: {poly.get_feature_names_out(['x1', 'x2'])}")

# For p features, degree 2 → p + p(p+1)/2 new features
# For p=10: 10 → 65 features
# For p=100: 100 → 5150 features (can be expensive!)

# Manual interaction feature
df = pd.DataFrame({'age': [25, 35, 45], 'income': [40000, 70000, 90000]})
df['age_income']   = df['age'] * df['income']       # Interaction
df['income_per_age'] = df['income'] / df['age']     # Ratio
print(f"\nWith interaction features:\n{df}")
```

### Binning / Discretization

```python
from sklearn.preprocessing import KBinsDiscretizer
import numpy as np

# Convert continuous feature into categorical bins
age = np.array([15, 22, 35, 45, 55, 72, 29, 18, 63, 40]).reshape(-1, 1)

# Method 1: Equal-width bins (uniform)
kbd_uniform = KBinsDiscretizer(n_bins=4, encode='ordinal', strategy='uniform')
age_binned = kbd_uniform.fit_transform(age)
print("Equal-width bins:", age_binned.flatten().astype(int))
# [0, 1, 1, 2, 2, 3, 1, 0, 3, 2]

# Method 2: Equal-frequency bins (quantile)
kbd_quantile = KBinsDiscretizer(n_bins=4, encode='ordinal', strategy='quantile')
age_binned_q = kbd_quantile.fit_transform(age)

# Method 3: Manual binning with pandas
age_series = pd.Series(age.flatten())
age_labels = pd.cut(
    age_series,
    bins=[0, 18, 30, 50, 100],
    labels=['minor', 'young_adult', 'adult', 'senior'],
    right=True
)
print("Manual bins:", age_labels.values)
```

---

## 9. Feature Engineering — Categorical Features

### Encoding Strategies

```python
import pandas as pd
import numpy as np
from sklearn.preprocessing import LabelEncoder, OrdinalEncoder, OneHotEncoder

# Sample data
df = pd.DataFrame({
    'city':     ['London', 'Paris', 'London', 'Tokyo', 'Paris', 'Tokyo'],
    'size':     ['S', 'M', 'L', 'XL', 'S', 'M'],
    'premium':  ['Yes', 'No', 'Yes', 'Yes', 'No', 'No'],
    'price':    [100, 80, 150, 200, 70, 95]
})

# ── 1. Label Encoding: for ordinal categories ─────────────────
# ONLY use when the categories have a meaningful order!
# Wrong: encoding cities as 0,1,2 implies London < Paris < Tokyo!
# Right: encoding size S < M < L < XL

ordinal_enc = OrdinalEncoder(
    categories=[['S', 'M', 'L', 'XL']]  # Explicit order
)
df['size_encoded'] = ordinal_enc.fit_transform(df[['size']])
print("Size encoded:", df['size_encoded'].values)
# [0, 1, 2, 3, 0, 1] — preserves ordering!

# ── 2. One-Hot Encoding: for nominal categories ───────────────
# Creates k binary columns (or k-1 to avoid multicollinearity)
ohe = OneHotEncoder(
    drop='first',           # Drop first to avoid dummy variable trap
    sparse_output=False,    # Return dense array (not sparse matrix)
    handle_unknown='ignore' # Silently encode unknown categories as all-zeros
)
city_encoded = ohe.fit_transform(df[['city']])
city_columns = ohe.get_feature_names_out(['city'])
print(f"\nOHE columns: {city_columns}")
print(city_encoded)

# With pandas
city_dummies = pd.get_dummies(df['city'], prefix='city', drop_first=True)
print(f"\nPandas get_dummies:\n{city_dummies}")

# ── 3. Binary Encoding: for Yes/No ───────────────────────────
df['premium_binary'] = (df['premium'] == 'Yes').astype(int)

# ── 4. Target Encoding (Mean Encoding) ───────────────────────
# Replace category with mean target value per category
# DANGER: causes target leakage if not done carefully!
# Must compute target encoding on training data only
# Best practice: use k-fold target encoding to reduce leakage

def target_encode_safe(X_train, y_train, X_test, col, smoothing=10):
    """
    Target encoding with smoothing to prevent overfitting.
    Smoothing: pulls extreme estimates toward global mean.
    """
    global_mean = y_train.mean()

    # Compute per-category statistics from training data
    stats = pd.DataFrame({'y': y_train, 'cat': X_train[col]}).groupby('cat')['y'].agg(
        ['mean', 'count']
    )

    # Bayesian smoothing: blend category mean toward global mean
    # More samples → trust category mean more
    # Few samples → pull toward global mean (regularization)
    stats['encoded'] = (
        (stats['count'] * stats['mean'] + smoothing * global_mean)
        / (stats['count'] + smoothing)
    )

    # Apply to train and test
    X_train_enc = X_train[col].map(stats['encoded']).fillna(global_mean)
    X_test_enc  = X_test[col].map(stats['encoded']).fillna(global_mean)

    return X_train_enc, X_test_enc
```

### Frequency Encoding

```python
# Replace category with its frequency in training data
# Simple, handles high-cardinality well, no leakage

def frequency_encode(X_train, X_test, col):
    freq_map = X_train[col].value_counts(normalize=True)
    X_train_enc = X_train[col].map(freq_map).fillna(0)
    X_test_enc  = X_test[col].map(freq_map).fillna(0)
    return X_train_enc, X_test_enc

# High-cardinality categories (city names, product IDs):
# One-hot encoding → too many columns (thousands!)
# Better: target encoding or frequency encoding

print("""
ENCODING SELECTION GUIDE:
────────────────────────────────────────────────────────────────────────
Binary (Yes/No, Male/Female):  Map to 0/1 directly.

Ordinal (S<M<L<XL, Low<Med<High): OrdinalEncoder with explicit order.

Low-cardinality nominal (< 20 categories):
  OneHotEncoder → most common, no ordering assumed.

High-cardinality nominal (> 20 categories, e.g., zip codes, products):
  Target encoding with smoothing → encodes signal into single column.
  Frequency encoding → simple, no leakage, encodes popularity.
  Hashing: pd.util.hash_pandas_object → fixed-size representation.

Tree-based models:
  Often handle label encoding fine (can split on any value).
  OneHot still sometimes helps.

Linear models/SVMs/KNN:
  Always OneHot (or target encode) for nominal categories.
  OrdinalEncoder only for actually ordinal categories.
""")
```

---

## 10. Handling Missing Values

```python
import numpy as np
import pandas as pd
from sklearn.impute import SimpleImputer, KNNImputer, IterativeImputer
from sklearn.experimental import enable_iterative_imputer  # Must import this!

# Create dataset with missing values
np.random.seed(42)
df = pd.DataFrame({
    'age':    [25, np.nan, 35, 45, np.nan, 55],
    'income': [40000, 70000, np.nan, 90000, 50000, np.nan],
    'city':   ['London', np.nan, 'Paris', 'Tokyo', 'London', np.nan],
    'churn':  [0, 1, 0, 0, 1, 1]
})

print("Missing values per column:")
print(df.isnull().sum())
print(f"\nMissing percentage:")
print((df.isnull().sum() / len(df) * 100).round(1))

# ── Numerical Imputation ───────────────────────────────────────
X_num = df[['age', 'income']].values

# 1. Mean/Median/Mode (Simple)
imp_mean   = SimpleImputer(strategy='mean')
imp_median = SimpleImputer(strategy='median')
imp_const  = SimpleImputer(strategy='constant', fill_value=0)

print("\nMean imputation:")
print(imp_mean.fit_transform(X_num))

print("\nMedian imputation (more robust to outliers):")
print(imp_median.fit_transform(X_num))

# 2. KNN Imputation (uses k nearest neighbors)
# Imputed value = weighted mean of k nearest neighbors' values
imp_knn = KNNImputer(n_neighbors=2, weights='uniform')
print("\nKNN imputation:")
print(imp_knn.fit_transform(X_num))

# 3. Iterative Imputation (MICE: Multiple Imputation by Chained Equations)
# Fits a model to predict each missing feature from all others
# Multiple passes until convergence
imp_iter = IterativeImputer(max_iter=10, random_state=42)
print("\nIterative imputation (MICE):")
print(imp_iter.fit_transform(X_num).round(1))

# ── Categorical Imputation ─────────────────────────────────────
imp_cat = SimpleImputer(strategy='most_frequent')  # Mode
city_imputed = imp_cat.fit_transform(df[['city']])
print(f"\nCity after mode imputation: {city_imputed.flatten()}")

# ── Missing Indicator (add "is_missing" feature) ─────────────
from sklearn.impute import MissingIndicator

indicator = MissingIndicator(features='missing-only')
missing_flags = indicator.fit_transform(df[['age', 'income']])
print(f"\nMissing indicator flags shape: {missing_flags.shape}")
print(missing_flags)
# This adds information: "was this feature originally missing?"
# Sometimes missingness itself is informative (e.g., income not reported)

print("""
MISSING VALUE STRATEGY GUIDE:
────────────────────────────────────────────────────────────────────────
< 5% missing:   Mean/median imputation is usually fine.
5-20% missing:  KNN or IterativeImputer recommended.
> 20% missing:  Consider dropping the feature, or model-based imputation.

Always add a binary "is_missing" indicator feature alongside imputation!
The fact that a value is missing is itself informative.

NEVER impute test data using test set statistics!
Fit the imputer on training data, transform both train and test.
Use Pipelines to ensure this automatically.

For tree models (RF, XGBoost):
  Many tree implementations handle NaN natively.
  XGBoost learns the optimal direction for missing values.
  For sklearn trees: must impute or use a library that handles it.
""")
```

---

## 11. Feature Selection

Too many features can:
- Slow down training
- Introduce noise that hurts generalization
- Create the curse of dimensionality
- Make models harder to interpret

```python
import numpy as np
import pandas as pd
from sklearn.feature_selection import (
    VarianceThreshold, SelectKBest, chi2, mutual_info_classif, f_classif,
    RFE, SelectFromModel, RFECV
)
from sklearn.ensemble import RandomForestClassifier
from sklearn.linear_model import Lasso
from sklearn.preprocessing import MinMaxScaler

# ── 1. Filter Methods ─────────────────────────────────────────
# Do not use a model — compute a statistic per feature independently

from sklearn.datasets import load_breast_cancer
X, y = load_breast_cancer(return_X_y=True)
feature_names = load_breast_cancer().feature_names

# Variance threshold: remove near-constant features
vt = VarianceThreshold(threshold=0.01)
X_high_var = vt.fit_transform(X)
print(f"VarianceThreshold: {X.shape[1]} → {X_high_var.shape[1]} features")

# Chi-squared test (for non-negative integer/count features with classification)
# Tests if feature and target are statistically independent
X_pos = MinMaxScaler().fit_transform(X)  # chi2 requires non-negative
selector_chi2 = SelectKBest(chi2, k=10)
X_chi2 = selector_chi2.fit_transform(X_pos, y)
chi2_scores = selector_chi2.scores_
chi2_top = np.argsort(chi2_scores)[::-1][:10]
print(f"\nTop 10 features by chi-squared:")
for idx in chi2_top:
    print(f"  {feature_names[idx]:40s}: chi2={chi2_scores[idx]:.2f}")

# Mutual Information (model-free, captures non-linear relationships)
mi_scores = mutual_info_classif(X, y, random_state=42)
mi_top = np.argsort(mi_scores)[::-1][:10]
print(f"\nTop 10 features by Mutual Information:")
for idx in mi_top:
    print(f"  {feature_names[idx]:40s}: MI={mi_scores[idx]:.4f}")

# ── 2. Wrapper Methods ────────────────────────────────────────
# Use a model to evaluate feature subsets

# RFE: trains model on all features, removes least important, repeats
from sklearn.linear_model import LogisticRegression
rfe = RFE(
    estimator=LogisticRegression(max_iter=5000),
    n_features_to_select=10,   # Keep top 10 features
    step=1,                    # Remove 1 feature at a time
    verbose=0
)
rfe.fit(X, y)  # This is slow for large feature sets
selected_rfe = feature_names[rfe.support_]
print(f"\nRFE selected features: {selected_rfe}")

# RFECV: automatically finds optimal n_features via cross-validation
rfecv = RFECV(
    estimator=LogisticRegression(max_iter=5000),
    min_features_to_select=5,
    cv=5,
    scoring='roc_auc',
    n_jobs=-1
)
rfecv.fit(X, y)
print(f"\nRFECV optimal n_features: {rfecv.n_features_}")
print(f"Selected: {feature_names[rfecv.support_]}")

# ── 3. Embedded Methods ────────────────────────────────────────
# Feature selection happens inside the model training

# Lasso: L1 drives some weights to zero → automatic selection
from sklearn.linear_model import LassoCV
lasso = LassoCV(cv=5, random_state=42, max_iter=5000)
lasso.fit(X, y.astype(float))  # Lasso is for regression
selected_lasso = feature_names[np.abs(lasso.coef_) > 0]
print(f"\nLasso selected {len(selected_lasso)} features:")
print(selected_lasso)

# Tree Feature Importance
from sklearn.inspection import permutation_importance

rf = RandomForestClassifier(n_estimators=200, random_state=42, n_jobs=-1)
rf.fit(X, y)

# Impurity-based importance (fast but biased toward high-cardinality)
imp_importance = pd.DataFrame({
    'feature': feature_names,
    'importance': rf.feature_importances_
}).sort_values('importance', ascending=False)
print(f"\nTop 10 Random Forest feature importances:")
print(imp_importance.head(10).to_string(index=False))

# Permutation importance (slower but unbiased and model-agnostic)
# For each feature: shuffle its values, measure decrease in score
perm_result = permutation_importance(
    rf, X, y,
    n_repeats=10,
    random_state=42,
    n_jobs=-1
)
perm_importance = pd.DataFrame({
    'feature': feature_names,
    'importance_mean': perm_result.importances_mean,
    'importance_std':  perm_result.importances_std
}).sort_values('importance_mean', ascending=False)
print(f"\nTop 10 Permutation feature importances:")
print(perm_importance.head(10).to_string(index=False))

# SelectFromModel: keep features with importance > threshold
sfm = SelectFromModel(
    RandomForestClassifier(n_estimators=100, random_state=42),
    threshold='median'   # Keep features with importance > median importance
)
sfm.fit(X, y)
X_selected = sfm.transform(X)
selected_sfm = feature_names[sfm.get_support()]
print(f"\nSelectFromModel: {X.shape[1]} → {X_selected.shape[1]} features")
```

---

## 12. Dealing with Outliers

```python
import numpy as np
import pandas as pd

# ── Detecting Outliers ────────────────────────────────────────
def detect_outliers(X, feature_names=None):
    """
    Detect outliers using multiple methods.
    Returns a boolean mask: True = outlier.
    """
    n_samples, n_features = X.shape
    feature_names = feature_names or [f'f{i}' for i in range(n_features)]

    print("Outlier Detection Report:")
    for i, name in enumerate(feature_names):
        col = X[:, i]

        # Z-score method
        z_scores = np.abs((col - col.mean()) / col.std())
        z_outliers = z_scores > 3

        # IQR method
        Q1, Q3 = np.percentile(col, [25, 75])
        IQR = Q3 - Q1
        iqr_outliers = (col < Q1 - 1.5*IQR) | (col > Q3 + 1.5*IQR)

        print(f"  {name:20s}: Z-score outliers={z_outliers.sum():3d}, "
              f"IQR outliers={iqr_outliers.sum():3d}")

# ── Handling Outliers ─────────────────────────────────────────
def winsorize(X, lower_percentile=1, upper_percentile=99):
    """
    Capping/Winsorizing: clip outliers to percentile thresholds.
    Better than removing outliers (preserves all samples).
    """
    lower = np.percentile(X, lower_percentile, axis=0)
    upper = np.percentile(X, upper_percentile, axis=0)
    return np.clip(X, lower, upper)

# For most ML purposes:
# Tree-based models (RF, XGBoost): ROBUST to outliers (no scaling needed)
# Linear models, SVMs, KNN: sensitive to outliers → scale or winsorize

print("""
OUTLIER HANDLING STRATEGIES:
────────────────────────────────────────────────────────────────────────
1. Winsorize/Cap: clip to 1st-99th percentile.
   Use when: outliers are genuine but extreme values.

2. Log transform: natural compression of large values.
   Use when: data is right-skewed (prices, incomes).

3. Remove: if outliers are data errors.
   Risk: lose information, change distribution.

4. RobustScaler: scale using median and IQR instead of mean and std.
   Use when: can't remove outliers, want resistant scaling.

5. Use robust models: tree-based models are naturally robust.
   Random Forest, XGBoost ignore individual extreme values.

6. Add outlier indicator: binary flag for "was this sample an outlier?"
   Lets the model learn from the outlier nature itself.
""")
```

---

## 13. Sklearn Pipelines — The Right Way

The Pipeline is the most important sklearn tool for production ML. It solves the test set contamination problem by design.

### The Problem Without Pipelines

```
WITHOUT PIPELINE (WRONG WAY)
────────────────────────────────────────────────────────────────────────
Step 1: scaler.fit_transform(X_train)  → learn mean/std from train ✓
Step 2: model.fit(X_train_scaled)       → train model ✓
Step 3: scaler.transform(X_test)        → apply train's stats to test ✓
Step 4: model.predict(X_test_scaled)    → predict ✓

This SEEMS correct, but fails in cross-validation!
In cross-val, you call: cross_val_score(model, X, y)
This calls model.fit(X_train_fold) ONLY — doesn't scale!
You need to manually fit the scaler inside each fold.

COMMON MISTAKE:
  X_scaled = scaler.fit_transform(X)  ← fits on ALL data (including val/test!)
  cross_val_score(model, X_scaled, y)  ← test fold was already seen by scaler
  → inflated, misleading cross-validation scores

WITH PIPELINE (CORRECT WAY)
────────────────────────────────────────────────────────────────────────
pipe = Pipeline([('scaler', StandardScaler()), ('model', LogisticRegression())])

cross_val_score(pipe, X, y)  ← for each fold:
  1. Fit scaler on training folds only
  2. Transform training folds with learned scaler
  3. Fit model on transformed training folds
  4. Transform validation fold with SAME scaler (no fitting!)
  5. Evaluate on transformed validation fold
  → CORRECT: no leakage!
```

### Building Pipelines

```python
from sklearn.pipeline import Pipeline, make_pipeline
from sklearn.preprocessing import StandardScaler
from sklearn.impute import SimpleImputer
from sklearn.linear_model import LogisticRegression
from sklearn.ensemble import RandomForestClassifier
from sklearn.decomposition import PCA
from sklearn.model_selection import cross_val_score, GridSearchCV

# ── Simple Pipeline ────────────────────────────────────────────
pipe_simple = Pipeline([
    ('imputer', SimpleImputer(strategy='median')),
    ('scaler',  StandardScaler()),
    ('model',   LogisticRegression(max_iter=1000))
])

# make_pipeline: shorthand (auto-names steps by class name)
pipe_simple2 = make_pipeline(
    SimpleImputer(strategy='median'),
    StandardScaler(),
    LogisticRegression(max_iter=1000)
)

# ── Using the Pipeline ─────────────────────────────────────────
from sklearn.datasets import load_breast_cancer
from sklearn.model_selection import train_test_split

X, y = load_breast_cancer(return_X_y=True)
# Inject some missing values for demonstration
X_missing = X.copy().astype(float)
np.random.seed(42)
mask = np.random.random(X.shape) < 0.1  # 10% missing
X_missing[mask] = np.nan

X_train, X_test, y_train, y_test = train_test_split(
    X_missing, y, test_size=0.2, random_state=42
)

pipe_simple.fit(X_train, y_train)
print(f"Pipeline test accuracy: {pipe_simple.score(X_test, y_test):.4f}")

# Cross-validate the pipeline (no leakage by construction!)
cv_scores = cross_val_score(pipe_simple, X_missing, y, cv=5, scoring='roc_auc')
print(f"Pipeline CV AUC: {cv_scores.mean():.4f} ± {cv_scores.std():.4f}")

# ── Access Intermediate Steps ──────────────────────────────────
print(f"\nSteps: {pipe_simple.named_steps.keys()}")
print(f"Scaler mean: {pipe_simple.named_steps['scaler'].mean_[:5]}")
print(f"Model coef:  {pipe_simple.named_steps['model'].coef_[0][:5]}")

# ── Hyperparameter Tuning: Reference Pipeline Parameters ──────
# Use __ to reference nested parameters: stepname__param
param_grid = {
    'scaler': [StandardScaler(), MinMaxScaler()],  # Try different scalers
    'model__C': [0.01, 0.1, 1.0, 10.0],
    'model__penalty': ['l1', 'l2'],
    'model__solver': ['liblinear'],
}

grid = GridSearchCV(
    Pipeline([('scaler', StandardScaler()), ('model', LogisticRegression(max_iter=1000))]),
    param_grid,
    cv=5, scoring='roc_auc', n_jobs=-1
)
grid.fit(X_train, y_train)
print(f"\nGrid Search best: {grid.best_params_}")
print(f"Best CV AUC: {grid.best_score_:.4f}")

# ── PCA in Pipeline ───────────────────────────────────────────
pipe_pca = Pipeline([
    ('imputer', SimpleImputer()),
    ('scaler',  StandardScaler()),
    ('pca',     PCA(n_components=0.95)),  # Keep 95% variance
    ('model',   RandomForestClassifier(n_estimators=100, random_state=42))
])

pipe_pca.fit(X_train, y_train)
print(f"\nPCA pipeline n_components: {pipe_pca.named_steps['pca'].n_components_}")
print(f"PCA pipeline AUC: {pipe_pca.score(X_test, y_test):.4f}")
```

---

## 14. ColumnTransformer — Mixed Data Types

Real datasets have a mix of numerical and categorical features. ColumnTransformer applies different preprocessing to different columns.

```python
from sklearn.compose import ColumnTransformer, make_column_transformer, make_column_selector
from sklearn.pipeline import Pipeline
from sklearn.preprocessing import (
    StandardScaler, OneHotEncoder, OrdinalEncoder
)
from sklearn.impute import SimpleImputer
from sklearn.ensemble import RandomForestClassifier

# Example dataset with mixed types
import pandas as pd
import numpy as np

np.random.seed(42)
n = 500
df = pd.DataFrame({
    'age':       np.random.randint(18, 80, n).astype(float),
    'income':    np.random.lognormal(10.5, 1.2, n),
    'score':     np.random.normal(700, 100, n),
    'city':      np.random.choice(['London', 'Paris', 'Tokyo', 'NYC'], n),
    'education': np.random.choice(['HS', 'College', 'Graduate', 'PhD'], n),
    'employed':  np.random.choice(['Yes', 'No'], n),
})
# Add missing values
for col in ['age', 'income', 'city']:
    df.loc[np.random.choice(n, 30), col] = np.nan

y = (df['income'] > df['income'].median()).astype(int)

# Define feature types
numerical_features = ['age', 'income', 'score']
ordinal_features   = ['education']
nominal_features   = ['city']
binary_features    = ['employed']

# Build preprocessing for each type
numerical_transformer = Pipeline([
    ('imputer', SimpleImputer(strategy='median')),
    ('scaler',  StandardScaler())
])

ordinal_transformer = Pipeline([
    ('imputer', SimpleImputer(strategy='most_frequent')),
    ('encoder', OrdinalEncoder(categories=[['HS', 'College', 'Graduate', 'PhD']]))
])

nominal_transformer = Pipeline([
    ('imputer', SimpleImputer(strategy='most_frequent')),
    ('encoder', OneHotEncoder(drop='first', handle_unknown='ignore', sparse_output=False))
])

binary_transformer = Pipeline([
    ('encoder', OrdinalEncoder(categories=[['No', 'Yes']]))  # No=0, Yes=1
])

# ColumnTransformer: glues it all together
preprocessor = ColumnTransformer(transformers=[
    ('num',    numerical_transformer, numerical_features),
    ('ord',    ordinal_transformer,   ordinal_features),
    ('nom',    nominal_transformer,   nominal_features),
    ('bin',    binary_transformer,    binary_features),
], remainder='drop')   # 'drop' or 'passthrough' for unspecified columns

# Full pipeline: preprocessing + model
full_pipeline = Pipeline([
    ('preprocessor', preprocessor),
    ('model',        RandomForestClassifier(n_estimators=100, random_state=42))
])

# Train/test split on DataFrame
from sklearn.model_selection import train_test_split, cross_val_score

X = df.drop(columns=None)  # X is the full dataframe (pipeline handles columns)
X_train, X_test, y_train, y_test = train_test_split(df, y, test_size=0.2, random_state=42)

full_pipeline.fit(X_train, y_train)

print(f"Full pipeline test accuracy: {full_pipeline.score(X_test, y_test):.4f}")
cv_scores = cross_val_score(full_pipeline, df, y, cv=5, scoring='roc_auc', n_jobs=-1)
print(f"Full pipeline CV AUC: {cv_scores.mean():.4f} ± {cv_scores.std():.4f}")

# Inspect what came out of the preprocessor
X_train_transformed = full_pipeline.named_steps['preprocessor'].transform(X_train)
print(f"\nInput shape: {X_train.shape}")
print(f"After preprocessing shape: {X_train_transformed.shape}")
# Numerical (3) + ordinal (1) + OHE nominal (n-1) + binary (1)

# ── Auto Column Selector ──────────────────────────────────────
# make_column_selector auto-detects column types by dtype
auto_preprocessor = make_column_transformer(
    (StandardScaler(), make_column_selector(dtype_include=np.number)),
    (OneHotEncoder(handle_unknown='ignore'), make_column_selector(dtype_include=object))
)
print(f"\nAuto-selector: numerical → scale, categorical → OHE")
```

---

## 15. Saving and Loading Pipelines

```python
import joblib
import os

# Save the trained pipeline
model_path = '/tmp/churn_pipeline.joblib'
joblib.dump(full_pipeline, model_path)
print(f"Pipeline saved to: {model_path}")
print(f"File size: {os.path.getsize(model_path) / 1024:.1f} KB")

# Load and use
loaded_pipeline = joblib.load(model_path)
preds = loaded_pipeline.predict(X_test)
print(f"Loaded pipeline predictions: {preds[:5]}")

# The loaded pipeline has EVERYTHING:
# - Fitted imputer (knows training medians/modes)
# - Fitted scaler (knows training means and stds)
# - Fitted encoders (knows vocabulary from training)
# - Fitted model (knows weights/trees from training)
# No re-fitting needed!

# For production: also save feature column names
feature_info = {
    'numerical_features': numerical_features,
    'ordinal_features':   ordinal_features,
    'nominal_features':   nominal_features,
    'binary_features':    binary_features,
    'target_name':        'high_income',
    'training_date':      '2024-01-01',
    'model_version':      '1.0.0',
}
joblib.dump(feature_info, '/tmp/feature_info.joblib')

# Packaging for deployment
print("""
DEPLOYMENT CHECKLIST:
────────────────────────────────────────────────────────────────────────
1. Serialize model: joblib.dump(pipeline, 'model.joblib')
2. Document feature contract: what features the model expects
3. Version the model: include model_version in saved artifacts
4. Wrap in API: FastAPI/Flask endpoint that accepts JSON, returns prediction
5. Monitor: track input data distribution and prediction distribution
   (data drift: new customers behave differently → model degrades)

See Chapter 14 (Churn Project) for a complete FastAPI deployment example.
""")
```

---

## 16. Summary

```
CHAPTER 13 KEY CONCEPTS
─────────────────────────────────────────────────────────────

EVALUATION:
  Accuracy fails for imbalanced data → use F1, AUC-PR
  Test contamination → Pipelines
  Single split variance → Cross-validation
  Time series → TimeSeriesSplit (never random split)

CROSS-VALIDATION:
  k=5 or 10 (balance of speed and accuracy)
  Stratified for classification (preserves class ratios)
  cross_validate() for multiple metrics simultaneously
  Always CV the full Pipeline, not just the model

IMBALANCED DATA:
  class_weight='balanced' → easiest fix
  SMOTE → creates synthetic minority samples
  Threshold tuning → optimize for your metric, not 0.5

HYPERPARAMETER TUNING:
  GridSearch: exhaustive, O(total_combinations)
  RandomSearch: sample uniformly, often better with large grids
  Bayesian (Optuna): intelligent sampling, best for expensive models

FEATURE ENGINEERING:
  Scaling: StandardScaler general, RobustScaler for outliers
  Log transform: right-skewed data
  OneHotEncoding: nominal categories (never for ordered)
  Target encoding: high-cardinality (with smoothing!)
  Missing values: impute + add missing indicator
  Feature selection: permutation importance > impurity importance

SKLEARN PIPELINES:
  Pipeline: chained steps, fit/transform/predict as one object
  ColumnTransformer: different preprocessing per column type
  No leakage: preprocessor fit on train fold only in CV
  joblib: serialize and load trained pipelines

The three-line workflow:
  preprocessor = ColumnTransformer(...)
  pipeline = Pipeline([('prep', preprocessor), ('model', YourModel())])
  pipeline.fit(X_train, y_train)  # Everything handled
```

---

## Mini Projects

### Mini Project 1: Cross-Validation Strategy Comparator

Compare KFold, StratifiedKFold, TimeSeriesSplit, and Leave-One-Out on the same dataset.

**Objective:** See how CV strategy affects variance of the accuracy estimate.

```python
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.model_selection import (KFold, StratifiedKFold, TimeSeriesSplit,
                                     LeaveOneOut, cross_val_score)
from sklearn.ensemble import RandomForestClassifier
from sklearn.datasets import make_classification
from sklearn.preprocessing import StandardScaler
import warnings
warnings.filterwarnings('ignore')

np.random.seed(42)

# Dataset 1: Imbalanced (90/10 split)
X_imb, y_imb = make_classification(n_samples=300, n_features=10, weights=[0.9, 0.1],
                                    random_state=42)

# Dataset 2: Balanced
X_bal, y_bal = make_classification(n_samples=300, n_features=10, weights=[0.5, 0.5],
                                    random_state=42)

# Dataset 3: Time-series-like (features have temporal structure)
X_ts = np.column_stack([
    np.sin(np.linspace(0, 10, 300)) + np.random.randn(300)*0.2,
    np.cos(np.linspace(0, 10, 300)) + np.random.randn(300)*0.2,
    np.random.randn(300, 8)
])
y_ts = (X_ts[:, 0] > 0).astype(int)

cv_strategies = {
    'KFold(5)':          KFold(n_splits=5, shuffle=True, random_state=42),
    'StratifiedKFold(5)': StratifiedKFold(n_splits=5, shuffle=True, random_state=42),
    'TimeSeriesSplit(5)': TimeSeriesSplit(n_splits=5),
}

clf = RandomForestClassifier(n_estimators=50, random_state=42)

fig, axes = plt.subplots(1, 3, figsize=(16, 5))
datasets = [
    (X_imb, y_imb, "Imbalanced (90/10)"),
    (X_bal, y_bal, "Balanced (50/50)"),
    (X_ts, y_ts, "Temporal Data"),
]

for ax, (X, y, title) in zip(axes, datasets):
    results = {}
    for cv_name, cv in cv_strategies.items():
        try:
            scores = cross_val_score(clf, X, y, cv=cv, scoring='f1_macro', n_jobs=-1)
            results[cv_name] = scores
        except Exception as e:
            results[cv_name] = np.array([0.0])

    names  = list(results.keys())
    means  = [r.mean() for r in results.values()]
    stds   = [r.std() for r in results.values()]
    colors = ['steelblue', 'forestgreen', 'tomato']

    bars = ax.bar(names, means, color=colors, alpha=0.7, yerr=stds,
                  capsize=5, error_kw={'linewidth':2})
    ax.set_ylabel('F1 Score (macro)')
    ax.set_title(title)
    ax.set_ylim(0, 1.1)
    ax.set_xticklabels(names, rotation=15, ha='right', fontsize=8)
    for bar, mean, std in zip(bars, means, stds):
        ax.text(bar.get_x() + bar.get_width()/2, mean + std + 0.02,
                f'{mean:.3f}±{std:.3f}', ha='center', fontsize=7)
    ax.grid(True, alpha=0.3, axis='y')

plt.suptitle("Cross-Validation Strategy Comparison", fontsize=14, fontweight='bold')
plt.tight_layout()
plt.savefig("cv_strategy_comparison.png", dpi=150)
plt.show()

# Also show how fold size affects LOO instability
print("\n--- LOO vs KFold on small dataset ---")
X_small, y_small = make_classification(n_samples=50, n_features=5, random_state=42)
loo_scores = cross_val_score(clf, X_small, y_small, cv=LeaveOneOut(), scoring='accuracy')
kf_scores  = cross_val_score(clf, X_small, y_small, cv=KFold(5, shuffle=True, random_state=42), scoring='accuracy')
print(f"LOO:   mean={loo_scores.mean():.3f}, std={loo_scores.std():.3f} (high variance expected)")
print(f"KFold: mean={kf_scores.mean():.3f}, std={kf_scores.std():.3f}")
```

---

### Mini Project 2: Feature Selection Showdown

Compare Filter (SelectKBest), Wrapper (RFE), and Embedded (Lasso) methods side by side.

**Objective:** Understand that different selection methods pick different features and why.

```python
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
from sklearn.datasets import load_breast_cancer
from sklearn.feature_selection import (SelectKBest, f_classif, RFE,
                                       RFECV, mutual_info_classif)
from sklearn.linear_model import LogisticRegression, Lasso
from sklearn.ensemble import RandomForestClassifier
from sklearn.model_selection import cross_val_score, StratifiedKFold
from sklearn.preprocessing import StandardScaler
from sklearn.pipeline import Pipeline

data = load_breast_cancer()
X, y = data.data, data.target
feature_names = data.feature_names
scaler = StandardScaler()
X_scaled = scaler.fit_transform(X)

K = 10  # Number of features to select

# Method 1: Filter — ANOVA F-test
skb = SelectKBest(f_classif, k=K).fit(X_scaled, y)
filter_mask = skb.get_support()

# Method 2: Filter — Mutual Information
skb_mi = SelectKBest(mutual_info_classif, k=K).fit(X_scaled, y)
mi_mask = skb_mi.get_support()

# Method 3: Wrapper — RFE with Logistic Regression
rfe = RFE(LogisticRegression(max_iter=1000, C=0.1), n_features_to_select=K)
rfe.fit(X_scaled, y)
rfe_mask = rfe.support_

# Method 4: Embedded — Random Forest importance
rf = RandomForestClassifier(n_estimators=100, random_state=42)
rf.fit(X_scaled, y)
rf_importances = rf.feature_importances_
rf_mask = rf_importances >= np.sort(rf_importances)[-K]

methods = {
    'ANOVA F-test': filter_mask,
    'Mutual Info': mi_mask,
    'RFE (LogReg)': rfe_mask,
    'RF Importance': rf_mask,
}

# Agreement matrix
n_methods = len(methods)
method_names = list(methods.keys())
agreement = np.zeros((n_methods, n_methods))
for i, (n1, m1) in enumerate(methods.items()):
    for j, (n2, m2) in enumerate(methods.items()):
        agreement[i, j] = np.sum(m1 & m2)

fig, axes = plt.subplots(2, 2, figsize=(14, 10))

# Feature selection heatmap
selection_matrix = np.array([mask.astype(int) for mask in methods.values()])
im = axes[0, 0].imshow(selection_matrix, aspect='auto', cmap='Blues',
                        vmin=0, vmax=1.5)
axes[0, 0].set_xticks(range(len(feature_names)))
axes[0, 0].set_xticklabels([f[:10] for f in feature_names], rotation=90, fontsize=6)
axes[0, 0].set_yticks(range(n_methods))
axes[0, 0].set_yticklabels(method_names)
axes[0, 0].set_title(f"Feature Selection Matrix (top {K} features)")
plt.colorbar(im, ax=axes[0, 0], shrink=0.6)

# Add checkmarks
for i in range(n_methods):
    for j in range(len(feature_names)):
        if selection_matrix[i, j]:
            axes[0, 0].text(j, i, '✓', ha='center', va='center',
                            color='white', fontsize=8, fontweight='bold')

# Agreement between methods
im2 = axes[0, 1].imshow(agreement, cmap='YlOrRd', vmin=0, vmax=K)
axes[0, 1].set_xticks(range(n_methods))
axes[0, 1].set_xticklabels(method_names, rotation=20, ha='right', fontsize=8)
axes[0, 1].set_yticks(range(n_methods))
axes[0, 1].set_yticklabels(method_names, fontsize=8)
axes[0, 1].set_title("Method Agreement (shared features)")
for i in range(n_methods):
    for j in range(n_methods):
        axes[0, 1].text(j, i, f'{int(agreement[i, j])}', ha='center',
                        va='center', fontsize=10, color='black')
plt.colorbar(im2, ax=axes[0, 1], shrink=0.6)

# CV accuracy per method
cv = StratifiedKFold(n_splits=5, shuffle=True, random_state=42)
base_clf = LogisticRegression(max_iter=1000, random_state=42)

method_scores = {}
for method_name, mask in methods.items():
    X_sel = X_scaled[:, mask]
    scores = cross_val_score(base_clf, X_sel, y, cv=cv, scoring='accuracy')
    method_scores[method_name] = scores

# Baseline: all features
all_scores = cross_val_score(base_clf, X_scaled, y, cv=cv, scoring='accuracy')
method_scores['All Features'] = all_scores

names_ext = list(method_scores.keys())
means_ext = [s.mean() for s in method_scores.values()]
stds_ext  = [s.std() for s in method_scores.values()]
colors_ext = ['steelblue']*n_methods + ['gray']

bars = axes[1, 0].bar(names_ext, means_ext, color=colors_ext, alpha=0.8,
                       yerr=stds_ext, capsize=4)
axes[1, 0].set_ylabel('CV Accuracy')
axes[1, 0].set_title('Accuracy per Selection Method\n(LogReg on selected features)')
axes[1, 0].set_xticklabels(names_ext, rotation=15, ha='right', fontsize=8)
axes[1, 0].set_ylim(0.85, 1.02)
axes[1, 0].grid(True, alpha=0.3, axis='y')
for bar, mean, std in zip(bars, means_ext, stds_ext):
    axes[1, 0].text(bar.get_x() + bar.get_width()/2, mean + std + 0.003,
                    f'{mean:.3f}', ha='center', fontsize=8)

# Vote-based selection: features selected by 3+ methods
vote_counts = sum(mask.astype(int) for mask in methods.values())
y_pos = np.arange(len(feature_names))
axes[1, 1].barh(y_pos, vote_counts,
                color=[f'C{v}' if v > 0 else 'lightgray' for v in vote_counts], alpha=0.8)
axes[1, 1].axvline(3, color='red', linestyle='--', label='3-method threshold')
axes[1, 1].set_yticks(y_pos)
axes[1, 1].set_yticklabels([f[:18] for f in feature_names], fontsize=6)
axes[1, 1].set_xlabel('# Methods that selected this feature')
axes[1, 1].set_title(f'Feature Vote Count\n(consensus features in red zone)')
axes[1, 1].legend()
axes[1, 1].grid(True, alpha=0.3, axis='x')

plt.tight_layout()
plt.savefig("feature_selection_showdown.png", dpi=150)
plt.show()

consensus = [feature_names[i] for i, v in enumerate(vote_counts) if v >= 3]
print(f"\nConsensus features (selected by 3+ methods): {len(consensus)}")
for f in consensus:
    print(f"  • {f}")
```

---

### Mini Project 3: Learning Curve Diagnostic Dashboard

Diagnose whether a model suffers from high bias or high variance using learning curves.

**Objective:** Learn to prescribe the right fix (more data vs. more complexity vs. regularization).

```python
import numpy as np
import matplotlib.pyplot as plt
from sklearn.model_selection import learning_curve
from sklearn.linear_model import LogisticRegression, Ridge
from sklearn.svm import SVC
from sklearn.tree import DecisionTreeClassifier
from sklearn.datasets import load_breast_cancer, make_classification
from sklearn.preprocessing import StandardScaler
from sklearn.pipeline import Pipeline

data = load_breast_cancer()
X, y = data.data, data.target

def plot_learning_curve(ax, estimator, X, y, title, cv=5, scoring='accuracy'):
    train_sizes, train_scores, val_scores = learning_curve(
        estimator, X, y,
        train_sizes=np.linspace(0.1, 1.0, 10),
        cv=cv, scoring=scoring, n_jobs=-1, shuffle=True, random_state=42
    )
    train_mean = train_scores.mean(axis=1)
    train_std  = train_scores.std(axis=1)
    val_mean   = val_scores.mean(axis=1)
    val_std    = val_scores.std(axis=1)

    ax.plot(train_sizes, train_mean, 'b-o', label='Train', markersize=4)
    ax.fill_between(train_sizes, train_mean - train_std, train_mean + train_std,
                    alpha=0.1, color='blue')
    ax.plot(train_sizes, val_mean, 'r-o', label='Validation', markersize=4)
    ax.fill_between(train_sizes, val_mean - val_std, val_mean + val_std,
                    alpha=0.1, color='red')

    gap = train_mean[-1] - val_mean[-1]
    if val_mean[-1] < 0.75:
        diagnosis = "HIGH BIAS\n(underfitting)\nFix: more complexity"
        diag_color = 'orange'
    elif gap > 0.10:
        diagnosis = "HIGH VARIANCE\n(overfitting)\nFix: regularize / more data"
        diag_color = 'red'
    else:
        diagnosis = "GOOD FIT\nGap: {:.2f}".format(gap)
        diag_color = 'green'

    ax.text(0.98, 0.05, diagnosis, transform=ax.transAxes, fontsize=7,
            ha='right', va='bottom', color=diag_color,
            bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))

    ax.set_title(f"{title}\n(final val={val_mean[-1]:.3f}, gap={gap:.3f})", fontsize=9)
    ax.set_xlabel("Training Set Size")
    ax.set_ylabel(scoring.capitalize())
    ax.legend(fontsize=7)
    ax.grid(True, alpha=0.3)
    ax.set_ylim(0.5, 1.05)
    return val_mean[-1], gap

scaler = StandardScaler()
X_scaled = scaler.fit_transform(X)

models = [
    (Pipeline([('scaler', StandardScaler()),
               ('clf', LogisticRegression(C=0.001, max_iter=1000))]),
     "Logistic Regression\n(C=0.001 — high bias)"),
    (Pipeline([('scaler', StandardScaler()),
               ('clf', LogisticRegression(C=1.0, max_iter=1000))]),
     "Logistic Regression\n(C=1.0 — balanced)"),
    (Pipeline([('scaler', StandardScaler()),
               ('clf', DecisionTreeClassifier(max_depth=None))]),
     "Decision Tree\n(unbounded — high variance)"),
    (Pipeline([('scaler', StandardScaler()),
               ('clf', DecisionTreeClassifier(max_depth=3))]),
     "Decision Tree\n(depth=3 — pruned)"),
    (Pipeline([('scaler', StandardScaler()),
               ('clf', SVC(C=0.01, gamma='scale'))]),
     "SVM\n(C=0.01 — high bias)"),
    (Pipeline([('scaler', StandardScaler()),
               ('clf', SVC(C=100, gamma='scale'))]),
     "SVM\n(C=100 — high variance risk)"),
]

fig, axes = plt.subplots(2, 3, figsize=(16, 9))
fig.suptitle("Learning Curve Diagnostics: Bias vs. Variance", fontsize=14, fontweight='bold')

for ax, (model, title) in zip(axes.ravel(), models):
    plot_learning_curve(ax, model, X, y, title)

plt.tight_layout()
plt.savefig("learning_curve_diagnostics.png", dpi=150)
plt.show()
print("Saved: learning_curve_diagnostics.png")
print("\nDiagnostic Guide:")
print("  High Bias:     Train ≈ Val (both low)  → add features, increase model complexity")
print("  High Variance: Train >> Val             → add data, regularize, prune, dropout")
print("  Good Fit:      Train ≈ Val (both high)  → you're done, or try ensembles")
```

---

## 17. Exercises

**Exercise 1:** Pipeline cross-validation correctness. Generate a dataset with `make_classification`. Demonstrate the "leaky" approach (scale all data, then cross-validate model) vs the correct Pipeline approach. Compare the reported CV accuracies. How much does leakage inflate the score?

**Exercise 2:** Feature engineering for house prices. Download the Ames Housing dataset (available from Kaggle or sklearn's `fetch_openml('house_prices')`). Engineer at least 5 new features (age of house, has_garage, total_sf = basement_sf + first_floor_sf + second_floor_sf, etc.). Compare your engineered feature model's RMSE to the raw feature model.

**Exercise 3:** Imbalanced dataset experiment. Using `make_classification(weights=[0.97, 0.03])`:
- Train a logistic regression with default settings. Report accuracy and recall.
- Train with `class_weight='balanced'`. Report accuracy and recall.
- Train with SMOTE oversampling. Report accuracy and recall.
- Tune the decision threshold on model 1 to achieve recall > 0.8. What is the precision?
- Which approach would you choose for a medical diagnosis system?

**Exercise 4:** ColumnTransformer pipeline on Titanic data. Download the Titanic dataset. Build a ColumnTransformer that:
- Numerically imputes and scales: Age, Fare
- Mode-imputes and OHE: Embarked, Pclass
- Extracts title from Name (Mr, Mrs, Miss, Master, etc.) and OHE it
- Binarizes: Sex (male/female → 0/1)
Train a RandomForest through this pipeline. Cross-validate.

**Exercise 5:** Hyperparameter tuning comparison. On the breast cancer dataset:
- GridSearchCV over 3×3×2 = 18 combinations (note time taken)
- RandomizedSearchCV with n_iter=18 (same budget, wider distributions)
- Optuna with n_trials=18 (Bayesian)
Compare: best CV AUC from each, time taken, best parameters found.
Which approach finds the best result with the same computational budget?

---

**Next Chapter →** [Project: Customer Churn Prediction](./14-project-churn-prediction.md)

*Theory is complete. Now apply everything in a real project. We'll build a full ML system from raw data to deployed API.*
