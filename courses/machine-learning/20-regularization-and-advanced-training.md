# Chapter 20: Regularization, Batch Normalization, and Advanced Training Techniques

> **"A model that fits your training data perfectly is often a model that has memorized noise. The goal is not to minimize training loss — it is to minimize validation loss. Everything in this chapter is about that distinction."**

---

## Table of Contents
1. [Overfitting vs Underfitting — Review and Deepening](#1-overfitting-vs-underfitting)
2. [L1 and L2 Regularization in Neural Networks](#2-l1-and-l2-regularization-in-neural-networks)
3. [Dropout — Theory and Practice](#3-dropout--theory-and-practice)
4. [Batch Normalization — Deep Dive](#4-batch-normalization--deep-dive)
5. [Layer Normalization and Other Norm Types](#5-layer-normalization-and-other-norm-types)
6. [Early Stopping](#6-early-stopping)
7. [Data Augmentation](#7-data-augmentation)
8. [Learning Rate Schedules](#8-learning-rate-schedules)
9. [Mixed Precision Training](#9-mixed-precision-training)
10. [Gradient Clipping](#10-gradient-clipping)
11. [Model Checkpointing](#11-model-checkpointing)
12. [Debugging Training — Diagnosing What's Wrong](#12-debugging-training--diagnosing-whats-wrong)
13. [Summary and What's Next](#13-summary)

---

## 1. Overfitting vs Underfitting

### Review

```
BIAS-VARIANCE TRADEOFF for Deep Networks:

Underfitting (high bias):
  Training loss: HIGH
  Validation loss: HIGH
  Gap (val - train): SMALL
  
  Model is too simple for the problem.
  Root cause: not enough capacity (too few parameters, too shallow)
              OR  not trained long enough
              OR  learning rate too small

Overfitting (high variance):
  Training loss: LOW
  Validation loss: HIGH
  Gap (val - train): LARGE
  
  Model has memorized the training set.
  Root cause: too many parameters relative to data
              OR too many epochs without early stopping
              OR no regularization

Perfect generalization:
  Training loss: LOW
  Validation loss: LOW (close to training)
  Gap: SMALL (< 5% typically acceptable)
```

### Loss Curves Diagnostic

```
Epoch
  │
  │ Loss
  │
  │  train: ─────────────────────────────╮
  │                                       ╰────────
  │  val:   ─────────────────╮
  │                          ╰──────────────────────  ← overfitting starts here
  │
  └──────────────────────────────────────────────────► epoch
  │          Phase 1           Phase 2
  │       (both decreasing) (train↓, val↑ = overfitting)
  │
  Solution for Phase 2:
    - Add dropout
    - Add L2/weight decay
    - Reduce model size
    - More data / augmentation
    - Early stopping (save model at Phase 1/2 boundary)
```

### Why Deep Networks Overfit Easily

A ResNet-50 has 25.6 million parameters. If your dataset has 10,000 samples, that's 2,560 parameters per sample. The model can trivially memorize every training example. Regularization is not optional — it's the core of making deep learning work on real data.

---

## 2. L1 and L2 Regularization in Neural Networks

### L2 Regularization (Weight Decay)

Add a penalty to the loss for large weights:

```
L_total = L_task + λ · Σ_i (w_i)²

Gradient of penalty: 2λ · w_i
Update rule: w ← w - α · (∂L_task/∂w + 2λ · w)
           = w · (1 - 2αλ) - α · ∂L_task/∂w
```

The term `(1 - 2αλ)` **shrinks the weight toward zero** at every step — this is "weight decay." It's equivalent to L2 regularization when using SGD.

**Effect**: encourages smaller, more diffuse weights. Prevents any single weight from becoming very large and dominating predictions.

### L1 Regularization

```
L_total = L_task + λ · Σ_i |w_i|

Gradient (where w ≠ 0): λ · sign(w_i)
```

L1 pushes weights toward exactly zero, creating **sparse weight matrices**. Less common in neural networks (hard to compute with batches) but theoretically useful for feature selection.

### The Adam/AdamW Subtlety

There's an important bug in the naive combination of L2 regularization with Adam:

```
Naive Adam + L2:
  gradient includes 2λw term
  Adam's adaptive scaling (divide by √v + ε) also scales the weight decay!
  → different features get different effective weight decay
  → NOT the same as true L2 regularization

AdamW (correct weight decay):
  weight_decay applied DIRECTLY to weights, AFTER the adaptive gradient step
  θ ← θ - α·[m̂/(√v̂+ε)] - α·λ·θ
  → uniform weight decay regardless of gradient history
  → true L2 regularization in the limit
```

```python
import torch.optim as optim

# WRONG: L2 in Adam (weight decay is scaled by adaptive LR)
optimizer_wrong = optim.Adam(model.parameters(), lr=1e-3, weight_decay=1e-4)

# CORRECT: AdamW (decoupled weight decay)
optimizer_correct = optim.AdamW(model.parameters(), lr=1e-3, weight_decay=1e-4)

# For SGD: weight_decay IS equivalent to L2 (no adaptive scaling)
optimizer_sgd = optim.SGD(model.parameters(), lr=0.01,
                           momentum=0.9, weight_decay=1e-4)

# Manual L2 in loss (for full control):
def compute_l2_loss(model, lambda_l2=1e-4):
    l2 = sum(p.pow(2.0).sum() for p in model.parameters())
    return lambda_l2 * l2

# Usage in training loop:
loss = criterion(outputs, targets) + compute_l2_loss(model)
```

### Practical Guidance

| Scenario | Recommendation |
|----------|---------------|
| Transformers / modern networks | Use AdamW with weight_decay=0.01-0.1 |
| CNNs with SGD | Use weight_decay=1e-4 in SGD |
| Adam optimizer | Use AdamW instead — more correct |
| L1 regularization | Rarely needed in NNs; use for explicit sparsity only |

---

## 3. Dropout — Theory and Practice

### What Dropout Does

```
During training:
  For each training step, each neuron is independently zeroed out
  with probability p (typical: 0.5 for FC, 0.1-0.3 for conv layers).
  
  Remaining neurons are SCALED UP by 1/(1-p) to maintain expected value.
  This is called "inverted dropout" and is what PyTorch uses.

  Example with p=0.5:
    Inputs: [1.0, 2.0, 3.0, 4.0]
    Mask:   [1,   0,   1,   0]    (randomly sampled)
    Output: [2.0, 0.0, 6.0, 0.0] (multiply surviving by 1/(1-0.5)=2)
  
  The scaling ensures: E[output] = E[input]  regardless of p

During evaluation:
  ALL neurons active, NO scaling needed (scaling was done during training).
  Output is deterministic.
```

### Why Dropout Works — Two Interpretations

**Interpretation 1: Prevents co-adaptation**

Neurons in a layer cannot rely on specific other neurons being present. Each neuron must learn a feature that is individually useful. This forces **redundant representations** — the network learns multiple ways to represent the same information.

**Interpretation 2: Ensemble of subnetworks**

Each dropout mask creates a different subnetwork. With p=0.5 and n neurons, you have 2^n possible subnetworks. Dropout trains an exponential ensemble and averages at test time.

```
Subnetwork 1 (mask [1,0,1,0]):  ○─[neurons 1,3]─○ → prediction 1
Subnetwork 2 (mask [0,1,0,1]):  ○─[neurons 2,4]─○ → prediction 2
Subnetwork 3 (mask [1,1,0,0]):  ○─[neurons 1,2]─○ → prediction 3
...
2^n subnetworks sampled during training

At test time: use all neurons with scale (1-p)  ≈ average of all subnetwork predictions
```

### Dropout in Practice

```python
import torch
import torch.nn as nn

# ─── Standard Dropout ─────────────────────────────────────────────
class MLPWithDropout(nn.Module):
    def __init__(self, input_dim, hidden_dim, output_dim, dropout_p=0.5):
        super().__init__()
        
        self.fc1 = nn.Linear(input_dim, hidden_dim)
        self.fc2 = nn.Linear(hidden_dim, hidden_dim)
        self.fc3 = nn.Linear(hidden_dim, output_dim)
        
        # Dropout applied between layers (not before first or after last)
        self.drop1 = nn.Dropout(p=dropout_p)
        self.drop2 = nn.Dropout(p=dropout_p)
        
        self.relu = nn.ReLU()
    
    def forward(self, x):
        x = self.drop1(self.relu(self.fc1(x)))  # dropout after activation
        x = self.drop2(self.relu(self.fc2(x)))
        return self.fc3(x)    # NO dropout before final output layer

# ─── Dropout for Conv Layers ──────────────────────────────────────
class ConvWithDropout(nn.Module):
    def __init__(self):
        super().__init__()
        self.conv1 = nn.Conv2d(3, 64, 3, padding=1)
        self.conv2 = nn.Conv2d(64, 128, 3, padding=1)
        
        # Dropout2d zeroes entire feature maps (not individual pixels)
        self.drop_spatial = nn.Dropout2d(p=0.1)
        
        self.pool = nn.AdaptiveAvgPool2d(1)
        self.fc = nn.Linear(128, 10)
        self.drop_fc = nn.Dropout(p=0.5)  # higher dropout for FC layer
    
    def forward(self, x):
        x = self.drop_spatial(torch.relu(self.conv1(x)))
        x = self.drop_spatial(torch.relu(self.conv2(x)))
        x = self.pool(x).flatten(1)
        x = self.drop_fc(x)
        return self.fc(x)

# ─── Critical: train vs eval mode ────────────────────────────────
model = MLPWithDropout(784, 512, 10)

model.train()
out_train = model(x)    # dropout active — stochastic output

model.eval()
out_eval = model(x)     # dropout disabled — deterministic output
```

### Dropout Values by Layer Type

| Layer Type | Typical Dropout p |
|------------|-------------------|
| Fully connected (hidden) | 0.3 – 0.5 |
| Fully connected (classifier head) | 0.2 – 0.5 |
| Conv layers | 0.05 – 0.15 |
| Embedding layers | 0.1 – 0.3 |
| LSTM recurrent connections | 0.2 – 0.5 |
| Transformer attention | 0.1 (standard) |

---

## 4. Batch Normalization — Deep Dive

### The Problem Batch Norm Solves

During training, as lower layers update their weights, the distribution of their outputs changes. Upper layers then receive inputs from a shifting distribution — they're essentially chasing a moving target. This is called **Internal Covariate Shift**.

```
Epoch 1:  Layer 1 outputs → distribution centered at 2.3, std=1.1
Epoch 2:  Layer 1 outputs → distribution centered at 1.8, std=0.9  ← shifted!
Epoch 3:  Layer 1 outputs → distribution centered at 3.1, std=1.4  ← shifted again!

Layer 2 must constantly readjust to this changing distribution.
→ Slower learning, requires smaller learning rates, more sensitive to initialization.
```

### How Batch Normalization Works

For each mini-batch B of size m, for each feature dimension:

```
Step 1: Compute batch statistics
  μ_B = (1/m) Σ_i x_i             (batch mean)
  σ²_B = (1/m) Σ_i (x_i - μ_B)²  (batch variance)

Step 2: Normalize
  x̂_i = (x_i - μ_B) / √(σ²_B + ε)  (ε = 1e-5, numerical stability)
  x̂_i now has mean ≈ 0, variance ≈ 1

Step 3: Scale and shift (with LEARNABLE parameters γ, β)
  y_i = γ · x̂_i + β

γ (gamma): learned scale   — initialized to 1
β (beta):  learned shift   — initialized to 0
```

**Why the learnable γ and β?** After normalization, the network might need a non-zero mean or non-unit variance. γ and β allow the network to undo the normalization if that's optimal. The network can learn "normalize to N(0,1)" (γ=1, β=0) or "normalize to N(5, 2)" (γ=2, β=5) or anything in between.

### Batch Norm During Inference

At inference time, we may have batch size = 1 (or even batch size = 0 in some scenarios). We can't compute batch statistics from 1 sample!

Solution: **running statistics** accumulated during training:

```
During training (each batch):
  running_mean = (1 - momentum) * running_mean + momentum * μ_B
  running_var  = (1 - momentum) * running_var  + momentum * σ²_B
  (momentum typically 0.1 in PyTorch, so 0.9 weight on running mean)

During inference:
  x̂_i = (x_i - running_mean) / √(running_var + ε)
  y_i  = γ · x̂_i + β
```

This is why `model.eval()` must be called during inference — BatchNorm switches from batch statistics to running statistics.

```python
import torch
import torch.nn as nn

# ─── Batch Normalization in PyTorch ──────────────────────────────

# For fully connected layers (1D):
bn1d = nn.BatchNorm1d(
    num_features=256,   # number of features/neurons
    eps=1e-5,           # ε for numerical stability
    momentum=0.1,       # moving average factor
    affine=True,        # if True, learnable γ and β (almost always True)
    track_running_stats=True   # if True, track running mean/var for inference
)
# Input: (batch, features)  or  (batch, features, seq_len)

# For conv layers (2D):
bn2d = nn.BatchNorm2d(
    num_features=64    # number of channels
)
# Input: (batch, channels, H, W)
# Normalizes over batch and spatial (H, W) dimensions

# Typical placement: after Conv/Linear, BEFORE activation
conv_bn_relu = nn.Sequential(
    nn.Conv2d(3, 64, 3, padding=1, bias=False),  # bias=False — redundant with BN
    nn.BatchNorm2d(64),
    nn.ReLU(inplace=True)
)

# Check running statistics:
print(bn2d.running_mean.shape)   # (64,) — one per channel
print(bn2d.running_var.shape)    # (64,)
print(bn2d.weight.shape)         # (64,) — this is γ
print(bn2d.bias.shape)           # (64,) — this is β

# During training: running stats update automatically
# During eval: running stats are frozen (model.eval() does this)
```

### Benefits of Batch Normalization

```
1. Enables much higher learning rates:
   Without BN: LR > 0.01 often causes divergence
   With BN:    LR = 0.1 or higher is often stable

2. Reduces sensitivity to initialization:
   Activations are renormalized each layer → bad init is corrected after first BN

3. Acts as regularization:
   Batch statistics are noisy (they depend on which samples are in the batch)
   This noise is like adding random noise to activations → mild regularization
   In practice, you can reduce dropout when using BN

4. Speeds up training:
   ~5-10x fewer epochs needed to converge in some experiments
   Each epoch is slightly slower (extra computation) but far fewer epochs total
```

### Limitations of Batch Normalization

```
1. Requires batch size > 1 for training (problematic with batch size = 1)
2. Doesn't work well for:
   - RNNs (variable sequence lengths → statistics change at each step)
   - Very small batches (statistics are too noisy)
   - Reinforcement learning (single experience replay)
3. Running statistics must be matched between training and inference environments
4. Creates "train/test discrepancy" if not handled carefully
```

---

## 5. Layer Normalization and Other Norm Types

### Layer Normalization

Normalizes over the **feature** dimension instead of the **batch** dimension:

```
BatchNorm:
  For each mini-batch → normalize across the BATCH (per-feature statistics)
  Works for images, classification — where batch axis makes sense

LayerNorm:
  For each sample → normalize across the FEATURE dimension (per-sample statistics)
  No batch dependency → works with batch_size=1 → works for RNNs and Transformers
```

```
For an input x of shape (batch, features):
  LayerNorm computes:
    μ = mean(x, dim=-1)          ← mean of ALL features for this ONE sample
    σ² = var(x, dim=-1)
    x̂ = (x - μ) / √(σ² + ε)
    y = γ · x̂ + β

For Transformers, input x has shape (batch, seq_len, d_model):
  LayerNorm normalizes over d_model (last dim) for each (batch, seq) position
```

```python
import torch.nn as nn

# LayerNorm:
ln = nn.LayerNorm(normalized_shape=256)  # normalize over last 256 dimensions
# Input: (batch, 256) → normalizes each batch sample independently

# For Transformers:
ln_transformer = nn.LayerNorm(normalized_shape=512)  # d_model
# Input: (batch, seq_len, 512) → normalizes each token embedding independently

# Same behavior in train AND eval mode — no running stats needed
```

### Normalization Type Comparison

```
                         BatchNorm     LayerNorm    InstanceNorm   GroupNorm
─────────────────────────────────────────────────────────────────────────────
Normalize over:          Batch+Spatial  Features     Spatial        Feature groups
Works with batch=1?      No            Yes           Yes            Yes
Sequence models?         No            Yes           No             Limited
Conv layers?             Excellent     Sometimes     Style transfer OK
Transformer?             No            Yes (default) No             No
Running stats needed?    Yes           No            No             No
─────────────────────────────────────────────────────────────────────────────
```

**The rule**: Use BatchNorm for CNNs. Use LayerNorm for Transformers and RNNs.

---

## 6. Early Stopping

Early stopping is the simplest and most effective regularization technique for neural networks. The idea: stop training when validation loss stops improving.

```
          Train loss        Validation loss
                                   ↑
                              Best model
                             (stop here)
                                   │
 ────────────────────────────────────────────────────────►
                                   │
                             Overfitting
                             zone begins

Implementation:
  patience: number of epochs to wait for improvement
  min_delta: minimum improvement to count as improvement
  best_loss: track the best validation loss seen
  counter: how many epochs since last improvement
```

```python
class EarlyStopping:
    """
    Stop training if validation loss doesn't improve for 'patience' epochs.
    Optionally save the best model.
    """
    
    def __init__(self, patience=7, min_delta=0.0001, checkpoint_path='best_model.pth'):
        self.patience = patience
        self.min_delta = min_delta
        self.checkpoint_path = checkpoint_path
        
        self.best_loss = float('inf')
        self.counter = 0
        self.should_stop = False
        self.best_epoch = 0
    
    def __call__(self, val_loss, model, epoch):
        """
        Call at the end of each epoch.
        Returns True if training should stop.
        """
        if val_loss < self.best_loss - self.min_delta:
            # Improvement found
            self.best_loss = val_loss
            self.counter = 0
            self.best_epoch = epoch
            
            # Save best model
            torch.save(model.state_dict(), self.checkpoint_path)
            print(f"  Checkpoint saved (val_loss={val_loss:.6f})")
            
            return False  # don't stop
        
        else:
            # No improvement
            self.counter += 1
            print(f"  No improvement ({self.counter}/{self.patience})")
            
            if self.counter >= self.patience:
                print(f"\nEarly stopping at epoch {epoch}")
                print(f"Best epoch: {self.best_epoch}, best val_loss: {self.best_loss:.6f}")
                return True  # stop!
            
            return False


# Usage in training loop:
early_stopping = EarlyStopping(patience=7, checkpoint_path='best_model.pth')

for epoch in range(max_epochs):
    train_loss = train_epoch(model, train_loader, optimizer)
    val_loss, val_acc = evaluate(model, val_loader)
    
    if early_stopping(val_loss, model, epoch):
        break

# Load best model:
model.load_state_dict(torch.load('best_model.pth'))
```

---

## 7. Data Augmentation

### Why Augmentation Regularizes

The model sees augmented versions of training samples — effectively increasing the dataset size. Crucially, the labels stay the same despite the transformations. This teaches the model:

```
"A cat is a cat regardless of whether it's flipped horizontally"
→ Forces horizontal flip invariance

"A cat is a cat even with slightly different brightness"
→ Forces illumination invariance

"A cat is still a cat even if part of it is obscured"
→ Forces partial occlusion robustness
```

### Image Augmentation Pipeline

```python
import torchvision.transforms as transforms
import torchvision.transforms.v2 as T  # newer API (PyTorch 2.0+)

# Standard pipeline for ImageNet training:
train_transforms = transforms.Compose([
    transforms.RandomResizedCrop(224, scale=(0.08, 1.0)),
    transforms.RandomHorizontalFlip(),
    transforms.RandAugment(num_ops=2, magnitude=9),   # automatic augmentation
    transforms.ToTensor(),
    transforms.Normalize(mean=[0.485, 0.456, 0.406], std=[0.229, 0.224, 0.225]),
    transforms.RandomErasing(p=0.1)
])

# Minimal (for small datasets):
minimal_transforms = transforms.Compose([
    transforms.RandomHorizontalFlip(),
    transforms.ColorJitter(brightness=0.2, contrast=0.2),
    transforms.ToTensor(),
    transforms.Normalize(...)
])
```

### MixUp and CutMix

```python
def apply_mixup(inputs, targets, alpha=0.2, num_classes=10):
    """
    MixUp augmentation: blend two images and their one-hot labels.
    
    inputs:  (batch, C, H, W)
    targets: (batch,) — integer class indices
    Returns: mixed inputs, mixed one-hot targets
    """
    import torch
    import numpy as np
    
    # Sample mixing coefficient
    lam = np.random.beta(alpha, alpha) if alpha > 0 else 1.0
    
    batch_size = inputs.size(0)
    index = torch.randperm(batch_size)  # random permutation for pairing
    
    # Mix inputs
    mixed_inputs = lam * inputs + (1 - lam) * inputs[index]
    
    # Convert to one-hot and mix labels
    one_hot = torch.zeros(batch_size, num_classes).scatter_(
        1, targets.unsqueeze(1), 1.0
    )
    mixed_labels = lam * one_hot + (1 - lam) * one_hot[index]
    
    return mixed_inputs, mixed_labels

def mixup_loss(criterion, outputs, mixed_labels):
    """Compute loss for mixed labels (soft cross-entropy)."""
    # Standard CrossEntropyLoss doesn't handle soft labels
    # Use KLDivLoss or manual computation:
    log_probs = torch.log_softmax(outputs, dim=1)
    return -(mixed_labels * log_probs).sum(dim=1).mean()
```

### Text Augmentation

```python
import random

def random_deletion(tokens, p=0.1):
    """Randomly delete tokens with probability p."""
    if len(tokens) == 1:
        return tokens
    return [t for t in tokens if random.random() > p]

def random_swap(tokens, n=1):
    """Randomly swap n pairs of tokens."""
    tokens = tokens.copy()
    for _ in range(n):
        i, j = random.sample(range(len(tokens)), 2)
        tokens[i], tokens[j] = tokens[j], tokens[i]
    return tokens

# Back-translation: translate to German, then back to English
# Requires a translation model or API (e.g., Google Translate, Helsinki-NLP)
# "This movie was great" → "Dieser Film war toll" → "This film was great"
# Creates paraphrase with same meaning but different surface form
```

---

## 8. Learning Rate Schedules

### The Standard Scheduler Toolkit in PyTorch

```python
import torch.optim as optim

optimizer = optim.AdamW(model.parameters(), lr=1e-3)

# ─── StepLR ─────────────────────────────────────────────────────
# Reduce LR by gamma every step_size epochs
scheduler = optim.lr_scheduler.StepLR(optimizer, step_size=30, gamma=0.1)
# LR: [1e-3 for 0-29] → [1e-4 for 30-59] → [1e-5 for 60-89]

# ─── MultiStepLR ────────────────────────────────────────────────
# Reduce LR at specific epochs
scheduler = optim.lr_scheduler.MultiStepLR(optimizer, milestones=[30, 60, 80], gamma=0.1)

# ─── ExponentialLR ─────────────────────────────────────────────
# Exponential decay each epoch: lr = lr * gamma
scheduler = optim.lr_scheduler.ExponentialLR(optimizer, gamma=0.95)

# ─── CosineAnnealingLR ──────────────────────────────────────────
# Cosine decay from lr_max to eta_min over T_max epochs
scheduler = optim.lr_scheduler.CosineAnnealingLR(
    optimizer, T_max=100, eta_min=1e-6
)

# ─── CosineAnnealingWarmRestarts ────────────────────────────────
# Cosine with periodic restarts (SGDR)
scheduler = optim.lr_scheduler.CosineAnnealingWarmRestarts(
    optimizer, T_0=10, T_mult=2   # restart every 10, 20, 40, 80... epochs
)

# ─── ReduceLROnPlateau ──────────────────────────────────────────
# Reduce when a metric stops improving
scheduler = optim.lr_scheduler.ReduceLROnPlateau(
    optimizer,
    mode='min',       # monitoring a metric to minimize (e.g., val_loss)
    factor=0.5,       # new_lr = lr * factor
    patience=5,       # wait 5 epochs with no improvement
    min_lr=1e-7,      # don't go below this
    verbose=True
)
# Usage: scheduler.step(val_loss)  ← call AFTER validation, not after every batch

# ─── OneCycleLR (Leslie Smith's 1cycle policy) ──────────────────
# Triangle: LR increases from base_lr → max_lr → final_lr
scheduler = optim.lr_scheduler.OneCycleLR(
    optimizer,
    max_lr=0.01,                         # peak LR
    steps_per_epoch=len(train_loader),
    epochs=20,
    pct_start=0.3,                       # fraction of training for warmup
    div_factor=25,                       # base_lr = max_lr / div_factor
    final_div_factor=1e4                 # final_lr = base_lr / final_div_factor
)
# OneCycleLR is called PER BATCH, not per epoch:
# for batch in dataloader: ... optimizer.step(); scheduler.step()
```

### Warmup + Cosine Decay (Transformer Standard)

```python
import math
import torch.optim as optim

def get_cosine_schedule_with_warmup(optimizer, num_warmup_steps, num_training_steps,
                                     min_lr_ratio=0.0):
    """
    Cosine annealing with linear warmup.
    Used as the default for BERT, RoBERTa, and most transformer training.
    
    num_warmup_steps: number of steps for linear warmup
    num_training_steps: total number of training steps
    """
    def lr_lambda(current_step):
        if current_step < num_warmup_steps:
            # Linear warmup
            return float(current_step) / float(max(1, num_warmup_steps))
        
        # Cosine decay from 1.0 to min_lr_ratio
        progress = float(current_step - num_warmup_steps) / \
                   float(max(1, num_training_steps - num_warmup_steps))
        cosine_decay = max(min_lr_ratio, 0.5 * (1.0 + math.cos(math.pi * progress)))
        return cosine_decay
    
    return optim.lr_scheduler.LambdaLR(optimizer, lr_lambda)


# Usage:
steps_per_epoch = len(train_loader)
total_steps = steps_per_epoch * n_epochs
warmup_steps = steps_per_epoch * 2   # warm up for 2 epochs

scheduler = get_cosine_schedule_with_warmup(
    optimizer,
    num_warmup_steps=warmup_steps,
    num_training_steps=total_steps
)

# LR profile:
# Epoch 0-1: linearly increases from 0 to lr_max (warmup)
# Epoch 2+:  cosine decay from lr_max to lr_min
```

### LR Finder

```python
def find_best_lr(model, train_loader, criterion, device,
                  start_lr=1e-7, end_lr=1.0, num_steps=200):
    """
    Sweep LR exponentially. Plot loss vs LR.
    Best LR is typically 10x before the steepest descent.
    """
    model.train()
    optimizer = optim.SGD(model.parameters(), lr=start_lr)
    
    # Exponential schedule from start_lr to end_lr over num_steps
    lr_multiplier = (end_lr / start_lr) ** (1.0 / num_steps)
    scheduler = optim.lr_scheduler.ExponentialLR(optimizer, gamma=lr_multiplier)
    
    lrs = []
    losses = []
    best_loss = float('inf')
    smoothed_loss = 0.0
    beta = 0.98   # for smoothing
    
    for step, (inputs, targets) in enumerate(train_loader):
        if step >= num_steps:
            break
        
        inputs, targets = inputs.to(device), targets.to(device)
        optimizer.zero_grad()
        outputs = model(inputs)
        loss = criterion(outputs, targets)
        
        # Smooth loss for cleaner curve
        smoothed_loss = beta * smoothed_loss + (1 - beta) * loss.item()
        debiased = smoothed_loss / (1 - beta ** (step + 1))
        
        lrs.append(scheduler.get_last_lr()[0])
        losses.append(debiased)
        
        if debiased < best_loss:
            best_loss = debiased
        
        # Stop if loss is too high (diverged)
        if debiased > 4 * best_loss:
            break
        
        loss.backward()
        optimizer.step()
        scheduler.step()
    
    # Plot: import matplotlib; plt.semilogx(lrs, losses); plt.show()
    # Best LR: ~10x before the minimum loss
    
    return lrs, losses
```

---

## 9. Mixed Precision Training

### FP16 vs BF16 — The Complete Picture

```
FP32 (default):
  1 sign + 8 exponent + 23 mantissa = 32 bits
  Range: ±3.4×10^38
  Used for: optimizer states, gradient accumulation

FP16 (half precision):
  1 sign + 5 exponent + 10 mantissa = 16 bits
  Range: ±65504  ← SMALL! Overflows if values > 65504
  Precision: ~3 decimal digits
  Problem: gradient underflow (very small gradients → 0)
  Fix: gradient scaling

BF16 (brain float16):
  1 sign + 8 exponent + 7 mantissa = 16 bits
  Range: same as FP32 (same exponent bits!)
  Precision: ~2 decimal digits
  No overflow problem
  Requires: A100, H100, or Apple M1/M2
  Preferred over FP16 when available
```

### Complete AMP Implementation

```python
import torch
import torch.cuda.amp as amp

def train_with_amp(model, train_loader, optimizer, criterion, device):
    """Training with Automatic Mixed Precision."""
    
    model.train()
    
    # GradScaler: maintains a scale factor to prevent FP16 underflow
    scaler = amp.GradScaler(
        init_scale=2**16,         # initial scale (65536)
        growth_factor=2.0,        # double scale every growth_interval steps
        backoff_factor=0.5,       # halve if overflow detected
        growth_interval=2000,     # steps between scale increases
        enabled=True              # False to disable AMP entirely
    )
    
    for inputs, targets in train_loader:
        inputs = inputs.to(device, dtype=torch.float32)
        targets = targets.to(device)
        
        optimizer.zero_grad()
        
        # AMP context: automatically chooses FP16/BF16 for eligible ops
        with amp.autocast(device_type='cuda', dtype=torch.float16):
            outputs = model(inputs)       # forward in FP16
            loss = criterion(outputs, targets)  # loss computation
        
        # Scale the loss BEFORE backward (prevents FP16 underflow)
        scaler.scale(loss).backward()
        
        # Unscale gradients before clipping (MUST do this order)
        scaler.unscale_(optimizer)
        
        # Clip gradients (after unscaling)
        torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
        
        # Step (internally unscales if not already done, checks for infs/NaNs)
        scaler.step(optimizer)
        
        # Update scale factor for next iteration
        scaler.update()
    
    # For BF16 on supported hardware:
    with amp.autocast(device_type='cuda', dtype=torch.bfloat16):
        outputs = model(inputs)

# Memory savings with AMP:
# Model parameters: FP32 (kept as is for stability)
# Activations: FP16 during forward/backward (~2x memory reduction)
# Optimizer states: FP32 (Adam stores m and v in FP32)
# Net result: ~1.5-2x less GPU memory for activations
# Speed: ~1.5-3x faster on modern GPUs with Tensor Cores
```

---

## 10. Gradient Clipping

```python
import torch

# Method 1: Clip by global norm (RECOMMENDED)
# Computes L2 norm of all parameter gradients combined
# If total norm > max_norm: scale ALL gradients proportionally
max_norm = 1.0

torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=max_norm)

# Method 2: Clip by value (less common)
# Clips individual gradient values to [-clip_value, +clip_value]
torch.nn.utils.clip_grad_value_(model.parameters(), clip_value=0.5)

# Always do this BEFORE optimizer.step():
loss.backward()
torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
optimizer.step()

# Monitor gradient norm to decide if clipping is helping:
def get_grad_norm(model):
    """Compute global gradient norm for monitoring."""
    total_norm = 0.0
    for p in model.parameters():
        if p.grad is not None:
            total_norm += p.grad.data.norm(2).item() ** 2
    return total_norm ** 0.5

# In training loop:
loss.backward()
grad_norm = get_grad_norm(model)
torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)
optimizer.step()

# Log: grad_norm before and after clipping
# If grad_norm before >> max_norm frequently: you may need smaller LR
```

**When to use gradient clipping:**
- Always for RNNs/LSTMs (common to have exploding gradients in recurrence)
- For Transformers (standard to clip at 1.0)
- For any model where you observe NaN losses

---

## 11. Model Checkpointing

```python
import torch
import os
from pathlib import Path

class ModelCheckpoint:
    """
    Save model checkpoints during training.
    Saves the best model and optionally the most recent checkpoint.
    """
    
    def __init__(
        self,
        checkpoint_dir='./checkpoints',
        monitor='val_loss',
        mode='min',           # 'min' if lower is better (loss), 'max' if higher (accuracy)
        save_best_only=True,
        verbose=True
    ):
        self.checkpoint_dir = Path(checkpoint_dir)
        self.checkpoint_dir.mkdir(parents=True, exist_ok=True)
        
        self.monitor = monitor
        self.mode = mode
        self.save_best_only = save_best_only
        self.verbose = verbose
        
        self.best_value = float('inf') if mode == 'min' else float('-inf')
    
    def _is_improvement(self, current_value):
        if self.mode == 'min':
            return current_value < self.best_value
        else:
            return current_value > self.best_value
    
    def save(self, model, optimizer, scheduler, epoch, metrics: dict):
        """
        Save checkpoint if monitored metric improved.
        Returns True if saved.
        """
        current_value = metrics.get(self.monitor)
        if current_value is None:
            raise ValueError(f"Metric '{self.monitor}' not found in metrics dict")
        
        if self._is_improvement(current_value):
            self.best_value = current_value
            
            checkpoint = {
                'epoch': epoch,
                'model_state_dict': model.state_dict(),
                'optimizer_state_dict': optimizer.state_dict(),
                'scheduler_state_dict': scheduler.state_dict() if scheduler else None,
                'metrics': metrics,
            }
            
            path = self.checkpoint_dir / 'best_model.pth'
            torch.save(checkpoint, path)
            
            if self.verbose:
                print(f"  Checkpoint saved: {self.monitor}={current_value:.6f}")
            
            return True
        
        return False
    
    @staticmethod
    def load(checkpoint_path, model, optimizer=None, scheduler=None):
        """Load model (and optionally optimizer/scheduler) from checkpoint."""
        checkpoint = torch.load(checkpoint_path, map_location='cpu')
        
        model.load_state_dict(checkpoint['model_state_dict'])
        
        if optimizer is not None and 'optimizer_state_dict' in checkpoint:
            optimizer.load_state_dict(checkpoint['optimizer_state_dict'])
        
        if scheduler is not None and checkpoint.get('scheduler_state_dict'):
            scheduler.load_state_dict(checkpoint['scheduler_state_dict'])
        
        epoch = checkpoint.get('epoch', 0)
        metrics = checkpoint.get('metrics', {})
        
        print(f"Loaded checkpoint from epoch {epoch}: {metrics}")
        
        return epoch, metrics


# Usage:
checkpoint = ModelCheckpoint(
    checkpoint_dir='./checkpoints',
    monitor='val_loss',
    mode='min',
    save_best_only=True
)

for epoch in range(n_epochs):
    train_loss = train_epoch(...)
    val_loss, val_acc = evaluate(...)
    
    checkpoint.save(
        model, optimizer, scheduler, epoch,
        metrics={'val_loss': val_loss, 'val_acc': val_acc, 'train_loss': train_loss}
    )
```

---

## 12. Debugging Training — Diagnosing What's Wrong

### Diagnostic Table

| Symptom | Possible Cause | Fix |
|---------|---------------|-----|
| Loss is NaN from the start | LR too large, or log(0) in loss | Reduce LR; add eps to log |
| Loss is NaN after N steps | Exploding gradient | Gradient clipping; reduce LR |
| Loss not decreasing at all | LR too small | Increase LR (10x) |
| Loss decreasing slowly | LR too small, bad init | LR finder; He init |
| Training acc >> val acc | Overfitting | Dropout, L2, augmentation, more data |
| Val acc >> train acc | (Rare) Wrong split | Check train/val sets are correct |
| Loss oscillates wildly | LR too large | Reduce LR |
| Loss decreases but accuracy doesn't improve | Class imbalance | Weighted loss; check labels |
| GPU memory error | Batch too large, model too big | Reduce batch; gradient accum; AMP |

### Gradient Flow Visualization

```python
import matplotlib.pyplot as plt
import numpy as np

def plot_gradient_flow(model):
    """
    Plot the gradient flow in the model.
    Helps detect vanishing/exploding gradients.
    Call after loss.backward(), before optimizer.step().
    """
    mean_grads = []
    max_grads = []
    layers = []
    
    for name, param in model.named_parameters():
        if param.requires_grad and param.grad is not None:
            layers.append(name)
            mean_grads.append(param.grad.abs().mean().item())
            max_grads.append(param.grad.abs().max().item())
    
    plt.figure(figsize=(16, 6))
    plt.subplot(1, 2, 1)
    plt.bar(range(len(mean_grads)), mean_grads, alpha=0.7)
    plt.xticks(range(len(layers)), layers, rotation=90)
    plt.title('Mean Absolute Gradient')
    plt.ylabel('Gradient magnitude')
    plt.yscale('log')
    
    plt.subplot(1, 2, 2)
    plt.bar(range(len(max_grads)), max_grads, alpha=0.7, color='orange')
    plt.xticks(range(len(layers)), layers, rotation=90)
    plt.title('Max Absolute Gradient')
    plt.yscale('log')
    
    plt.tight_layout()
    plt.show()

# Use in training loop to check gradient health:
loss.backward()
plot_gradient_flow(model)    # call this occasionally (not every step — too slow)
optimizer.step()
```

### Activation Statistics Monitoring

```python
def monitor_activations(model, sample_input, device='cpu'):
    """
    Monitor activation statistics through all layers.
    Helps catch: dead ReLUs, vanishing/exploding activations.
    """
    stats = {}
    hooks = []
    
    def make_hook(name):
        def hook(module, input, output):
            if isinstance(output, torch.Tensor):
                stats[name] = {
                    'mean': output.mean().item(),
                    'std': output.std().item(),
                    'frac_zero': (output == 0).float().mean().item(),  # for ReLU
                    'has_nan': torch.isnan(output).any().item()
                }
        return hook
    
    for name, module in model.named_modules():
        if isinstance(module, (nn.ReLU, nn.Linear, nn.Conv2d, nn.BatchNorm2d)):
            hooks.append(module.register_forward_hook(make_hook(name)))
    
    model.eval()
    with torch.no_grad():
        model(sample_input.to(device))
    
    for h in hooks:
        h.remove()
    
    print("\nActivation Statistics:")
    print(f"{'Layer':<40} {'Mean':>8} {'Std':>8} {'%Zero':>8} {'NaN':>5}")
    print("─" * 70)
    for name, s in stats.items():
        print(f"{name:<40} {s['mean']:>8.3f} {s['std']:>8.3f} "
              f"{s['frac_zero']:>7.1%} {str(s['has_nan']):>5}")

# Red flags to look for:
# frac_zero > 0.5 for ReLU: possible dying ReLU problem
# std near 0: activations collapsed → add BatchNorm or fix initialization
# std >> 10: activations exploding → reduce LR or add gradient clipping
# has_nan = True: numerical instability → reduce LR, add gradient clipping
```

### Overfit a Single Batch — The Essential Sanity Check

Before any full training run, verify your code is correct:

```python
def sanity_check_single_batch(model, train_loader, optimizer, criterion, device,
                               n_steps=100, target_loss=0.01):
    """
    Overfit a single batch to verify:
    1. Forward pass is correct
    2. Backward pass computes gradients
    3. Model has enough capacity for this task
    4. Loss function is implemented correctly
    """
    model.train()
    
    # Get ONE batch and keep using it
    inputs, targets = next(iter(train_loader))
    inputs = inputs.to(device)
    targets = targets.to(device)
    
    print("Sanity check: overfitting single batch...")
    
    for step in range(n_steps):
        optimizer.zero_grad()
        outputs = model(inputs)
        loss = criterion(outputs, targets)
        loss.backward()
        optimizer.step()
        
        if step % 10 == 0:
            print(f"  Step {step:3d}: loss = {loss.item():.6f}")
        
        if loss.item() < target_loss:
            print(f"\nSanity check PASSED at step {step}")
            print(f"Model can overfit a single batch → forward/backward are correct")
            return True
    
    print(f"\nSanity check FAILED: loss={loss.item():.6f} still too high")
    print("Possible issues:")
    print("  - LR too small (loss barely moved)")
    print("  - LR too large (loss diverged)")
    print("  - Wrong loss function")
    print("  - Bug in forward pass")
    return False

# ALWAYS run this before long training runs
# If a model can't overfit 1 batch, there's a bug — not a hyperparameter issue
```

---

## 13. Summary

```
REGULARIZATION AND ADVANCED TRAINING
│
├── Diagnosing overfitting:
│     train_loss << val_loss → overfitting
│     both losses high → underfitting
│
├── L2 / Weight Decay:
│     Penalty λΣw² → shrinks weights toward zero
│     Use AdamW for correct decoupled weight decay
│     Typical λ: 1e-4 to 0.1
│
├── Dropout:
│     Training: randomly zero p fraction of neurons
│     Inference: all neurons active (model.eval()!)
│     Typical p: 0.1-0.3 (conv), 0.3-0.5 (FC)
│     model.train() vs model.eval() — critical!
│
├── Batch Normalization:
│     Normalize activations per mini-batch → μ=0, σ=1
│     Then scale+shift with learnable γ, β
│     Training: use batch stats; Inference: use running stats
│     Key benefit: higher LR, less sensitivity to init
│     Use for CNNs; use LayerNorm for Transformers/RNNs
│
├── Early Stopping:
│     Monitor val_loss; stop when no improvement for patience epochs
│     Save best model at improvement
│
├── Data Augmentation:
│     Images: crop, flip, color jitter, mixup, cutmix
│     Text: back-translation, swap, deletion
│
├── LR Schedules:
│     Cosine annealing: smooth, effective
│     Warmup + decay: essential for Transformers
│     ReduceLROnPlateau: adaptive, safe default
│
├── Mixed Precision (AMP):
│     FP16/BF16 forward, FP32 optimizer
│     GradScaler prevents underflow
│     ~2x memory, ~1.5-3x speed
│
└── Debugging:
      Overfit 1 batch first
      Check gradient flow (vanishing/exploding)
      Monitor activation statistics
      model.train() vs eval() is the most common bug
```

### Quick Decision Guide

```
Model overfitting? (val >> train)
  → More data first (always best)
  → Augmentation
  → Dropout (p=0.3-0.5 for FC, 0.1 for conv)
  → L2 regularization (AdamW with weight_decay)
  → Early stopping
  → Smaller model (fewer layers/neurons)

Model underfitting? (both losses high)
  → More epochs
  → Larger LR
  → More capacity (more layers/neurons)
  → Better initialization (He for ReLU)
  → BatchNorm (enables higher LR)

Training unstable? (loss oscillates or NaN)
  → Reduce LR
  → Gradient clipping (max_norm=1.0)
  → BatchNorm
  → Check for log(0) in loss
  → Smaller batch size (more gradient noise = regularization but also noise)
```

---

## Mini Projects

### Mini Project 1: Dropout Rate Sweeper

Train the same network with dropout rates from 0 to 0.8 and find the sweet spot.

**Objective:** See empirically how dropout rate affects bias-variance tradeoff.

```python
import torch
import torch.nn as nn
import numpy as np
import matplotlib.pyplot as plt
from sklearn.datasets import make_classification
from torch.utils.data import TensorDataset, DataLoader, random_split

np.random.seed(42); torch.manual_seed(42)
X_np, y_np = make_classification(n_samples=1000, n_features=20, n_informative=10,
                                   n_redundant=5, random_state=42)
X = torch.FloatTensor((X_np - X_np.mean(0)) / (X_np.std(0) + 1e-8))
y = torch.LongTensor(y_np)
ds = TensorDataset(X, y)
n = len(ds)
tr_ds, val_ds = random_split(ds, [int(0.8*n), n-int(0.8*n)],
                                generator=torch.Generator().manual_seed(42))
tr_loader  = DataLoader(tr_ds,  batch_size=64, shuffle=True)
val_loader = DataLoader(val_ds, batch_size=256)

class MLP(nn.Module):
    def __init__(self, dropout_p=0.0):
        super().__init__()
        self.net = nn.Sequential(
            nn.Linear(20, 256), nn.ReLU(), nn.Dropout(dropout_p),
            nn.Linear(256, 256), nn.ReLU(), nn.Dropout(dropout_p),
            nn.Linear(256, 128), nn.ReLU(), nn.Dropout(dropout_p),
            nn.Linear(128, 2)
        )
    def forward(self, x): return self.net(x)

def evaluate(model, loader):
    model.eval()
    correct = total = 0
    with torch.no_grad():
        for X_b, y_b in loader:
            correct += (model(X_b).argmax(1) == y_b).sum().item()
            total   += len(y_b)
    return correct / total

dropout_rates = [0.0, 0.1, 0.2, 0.3, 0.5, 0.7, 0.8]
results = {}

for p in dropout_rates:
    torch.manual_seed(42)
    model = MLP(dropout_p=p)
    opt = torch.optim.Adam(model.parameters(), lr=0.001)
    crit = nn.CrossEntropyLoss()
    tr_accs, val_accs = [], []
    for epoch in range(50):
        model.train()
        for X_b, y_b in tr_loader:
            opt.zero_grad(); crit(model(X_b), y_b).backward(); opt.step()
        tr_accs.append(evaluate(model, tr_loader))
        val_accs.append(evaluate(model, val_loader))
    results[p] = (tr_accs, val_accs)
    print(f"  p={p:.1f}: train={tr_accs[-1]:.3f}, val={val_accs[-1]:.3f}, gap={tr_accs[-1]-val_accs[-1]:.3f}")

fig, axes = plt.subplots(1, 3, figsize=(16, 5))
fig.suptitle("Dropout Rate Study", fontsize=13, fontweight='bold')
colors = plt.cm.viridis(np.linspace(0, 1, len(dropout_rates)))

for (p, (tr, val)), color in zip(results.items(), colors):
    axes[0].plot(range(50), tr,  '--', color=color, alpha=0.6, linewidth=1)
    axes[0].plot(range(50), val, '-', color=color, linewidth=1.5, label=f'p={p}')
axes[0].set_title("Val Accuracy Curves (dashed=train, solid=val)")
axes[0].set_xlabel("Epoch"); axes[0].set_ylabel("Accuracy"); axes[0].legend(fontsize=7); axes[0].grid(True, alpha=0.3)

final_tr  = [results[p][0][-1] for p in dropout_rates]
final_val = [results[p][1][-1] for p in dropout_rates]
gap = [t - v for t, v in zip(final_tr, final_val)]
x = np.arange(len(dropout_rates))
axes[1].bar(x - 0.15, final_tr,  0.3, label='Train', color='steelblue', alpha=0.8)
axes[1].bar(x + 0.15, final_val, 0.3, label='Val',   color='tomato',    alpha=0.8)
axes[1].set_xticks(x); axes[1].set_xticklabels([f'p={p}' for p in dropout_rates], rotation=30)
axes[1].set_ylabel('Final Accuracy'); axes[1].set_title('Final Train vs Val Accuracy')
axes[1].legend(); axes[1].grid(True, alpha=0.3, axis='y')

best_idx = np.argmax(final_val)
axes[2].bar(x, gap, color=['green' if g < 0.05 else 'orange' if g < 0.1 else 'red' for g in gap], alpha=0.8)
axes[2].set_xticks(x); axes[2].set_xticklabels([f'p={p}' for p in dropout_rates], rotation=30)
axes[2].axhline(0.05, color='orange', linestyle='--', label='5% gap (mild overfit)')
axes[2].set_ylabel('Train - Val Gap'); axes[2].set_title('Overfitting Gap by Dropout Rate')
axes[2].legend(); axes[2].grid(True, alpha=0.3, axis='y')
axes[2].annotate(f'Best val\np={dropout_rates[best_idx]}', xy=(best_idx, gap[best_idx]),
                  xytext=(best_idx+0.5, gap[best_idx]+0.02), fontsize=9,
                  arrowprops=dict(arrowstyle='->'))
plt.tight_layout()
plt.savefig("dropout_sweeper.png", dpi=150)
plt.show()
print(f"\nBest dropout rate: {dropout_rates[best_idx]} (val_acc={final_val[best_idx]:.3f})")
```

---

### Mini Project 2: Batch Normalization Effect Visualizer

Train networks with and without BatchNorm and visualize how it changes the gradient flow and convergence.

**Objective:** Understand why BatchNorm became standard — not just accuracy but training stability.

```python
import torch
import torch.nn as nn
import numpy as np
import matplotlib.pyplot as plt
from torch.utils.data import TensorDataset, DataLoader, random_split
from sklearn.datasets import make_classification

np.random.seed(42); torch.manual_seed(42)
X_np, y_np = make_classification(n_samples=1500, n_features=30, n_informative=15,
                                   n_redundant=5, random_state=42)
X = torch.FloatTensor(X_np)  # intentionally NOT pre-normalized
y = torch.LongTensor(y_np)
ds = TensorDataset(X, y)
n = len(ds)
tr_ds, val_ds = random_split(ds, [int(0.8*n), n-int(0.8*n)],
                                generator=torch.Generator().manual_seed(42))
tr_loader  = DataLoader(tr_ds,  batch_size=64, shuffle=True)
val_loader = DataLoader(val_ds, batch_size=256)

def make_network(use_batchnorm, use_dropout=True, hidden=256):
    layers = []
    prev = 30
    for next_dim in [hidden, hidden, hidden//2, hidden//4]:
        layers.append(nn.Linear(prev, next_dim))
        if use_batchnorm:
            layers.append(nn.BatchNorm1d(next_dim))
        layers.append(nn.ReLU())
        if use_dropout:
            layers.append(nn.Dropout(0.2))
        prev = next_dim
    layers.append(nn.Linear(prev, 2))
    return nn.Sequential(*layers)

def train_and_track(use_batchnorm, label, n_epochs=40, lr=0.01):
    torch.manual_seed(42)
    model = make_network(use_batchnorm)
    opt   = torch.optim.SGD(model.parameters(), lr=lr, momentum=0.9)
    crit  = nn.CrossEntropyLoss()
    tr_losses, val_losses, val_accs, grad_norms = [], [], [], []

    for epoch in range(n_epochs):
        model.train()
        epoch_grad_norms = []
        for X_b, y_b in tr_loader:
            opt.zero_grad()
            loss = crit(model(X_b), y_b)
            loss.backward()
            # Track gradient norm before clipping
            total_norm = sum(p.grad.norm().item()**2 for p in model.parameters()
                             if p.grad is not None) ** 0.5
            epoch_grad_norms.append(total_norm)
            opt.step()

        tr_loss = crit(model(X), y).item()
        model.eval()
        with torch.no_grad():
            val_loss = sum(crit(model(X_b), y_b).item() for X_b, y_b in val_loader)/len(val_loader)
            correct  = sum((model(X_b).argmax(1)==y_b).sum().item() for X_b, y_b in val_loader)
        tr_losses.append(tr_loss)
        val_losses.append(val_loss)
        val_accs.append(correct / len(val_ds))
        grad_norms.append(np.mean(epoch_grad_norms))

    print(f"  {label}: final val_acc={val_accs[-1]:.3f}")
    return tr_losses, val_losses, val_accs, grad_norms

print("Training with SGD (lr=0.01) — BatchNorm makes a big difference here!")
no_bn   = train_and_track(False, "No BatchNorm")
with_bn = train_and_track(True,  "With BatchNorm")

fig, axes = plt.subplots(2, 2, figsize=(14, 9))
fig.suptitle("BatchNorm Effect: Convergence, Stability, and Gradient Flow", fontsize=13, fontweight='bold')
epochs = range(1, 41)

axes[0,0].plot(epochs, no_bn[0],   'r--', label='Train (No BN)',   linewidth=1.5)
axes[0,0].plot(epochs, no_bn[1],   'r-',  label='Val (No BN)',     linewidth=2)
axes[0,0].plot(epochs, with_bn[0], 'b--', label='Train (With BN)', linewidth=1.5)
axes[0,0].plot(epochs, with_bn[1], 'b-',  label='Val (With BN)',   linewidth=2)
axes[0,0].set_title("Loss Curves"); axes[0,0].set_xlabel("Epoch"); axes[0,0].legend(fontsize=8); axes[0,0].grid(True, alpha=0.3)

axes[0,1].plot(epochs, no_bn[2],   'r-', label='No BN',   linewidth=2)
axes[0,1].plot(epochs, with_bn[2], 'b-', label='With BN', linewidth=2)
axes[0,1].set_title("Validation Accuracy"); axes[0,1].set_xlabel("Epoch"); axes[0,1].legend(); axes[0,1].grid(True, alpha=0.3)

axes[1,0].plot(epochs, no_bn[3],   'r-', label='No BN',   linewidth=2, alpha=0.7)
axes[1,0].plot(epochs, with_bn[3], 'b-', label='With BN', linewidth=2, alpha=0.7)
axes[1,0].set_title("Gradient Norm per Epoch\n(higher = better flow, spiky = instability)")
axes[1,0].set_xlabel("Epoch"); axes[1,0].legend(); axes[1,0].grid(True, alpha=0.3)
axes[1,0].set_yscale('log')

# Activation distribution comparison
torch.manual_seed(42)
net_no_bn   = make_network(False)
net_with_bn = make_network(True)
X_sample = X[:100]
activations_no_bn, activations_bn = [], []
def get_act_hook(lst):
    def hook(m, inp, out):
        if isinstance(m, nn.ReLU):
            lst.append(out.detach().flatten().numpy())
    return hook
for m in net_no_bn.modules():   m.register_forward_hook(get_act_hook(activations_no_bn)) if isinstance(m, nn.ReLU) else None
for m in net_with_bn.modules(): m.register_forward_hook(get_act_hook(activations_bn)) if isinstance(m, nn.ReLU) else None
with torch.no_grad():
    net_no_bn(X_sample); net_with_bn(X_sample)

axes[1,1].hist(np.concatenate(activations_no_bn), bins=50, alpha=0.6, color='red',  label='No BN',   density=True)
axes[1,1].hist(np.concatenate(activations_bn),    bins=50, alpha=0.6, color='blue', label='With BN', density=True)
axes[1,1].set_title("Activation Distribution (all ReLU layers combined)\n(BN keeps activations well-spread)")
axes[1,1].set_xlabel("Activation Value"); axes[1,1].legend(); axes[1,1].grid(True, alpha=0.3)
plt.tight_layout()
plt.savefig("batchnorm_effect.png", dpi=150)
plt.show()
print("Saved: dropout_sweeper.png, batchnorm_effect.png")
```

---

## Exercises

1. **Dropout experiment**: Train a 5-layer MLP on MNIST with p=0 (no dropout), p=0.3, p=0.5, p=0.7. Plot training and validation curves for all four. Find the optimal p.

2. **Batch Norm modes**: Train a ResNet on CIFAR-10. At epoch 5, call model.eval() during training (forgetting to call model.train() again). Observe how training loss behaves and explain why.

3. **Weight decay tuning**: Train the MNIST MLP with AdamW and weight_decay ∈ {0, 1e-5, 1e-4, 1e-3, 1e-2}. Plot val accuracy vs weight_decay. Is there a clear optimal?

4. **LR schedule comparison**: Train ResNet-18 on CIFAR-10 with: (a) constant LR=0.01; (b) StepLR (decay by 0.1 every 30 epochs); (c) CosineAnnealingLR; (d) OneCycleLR. Compare final accuracy and convergence speed.

5. **AMP benchmark**: Take a ResNet-50. Measure forward pass time and GPU memory with and without torch.cuda.amp.autocast. Report speedup and memory reduction.

6. **Early stopping**: Train a model for 200 epochs with patience=10. Verify that the saved checkpoint corresponds to the true best epoch (compare manually by inspecting val_loss at each epoch).

7. **Gradient flow**: After 1 backward pass through a 10-layer MLP with sigmoid activations, plot the gradient norm at each layer. Repeat with ReLU. Quantify the gradient decay ratio per layer.

---

**Chapter Summary:**

Regularization bridges the gap between training and validation performance. L2 regularization (weight decay) penalizes large weights; AdamW implements it correctly with decoupled decay. Dropout randomly zeros neurons during training, preventing co-adaptation and effectively training an ensemble of subnetworks — model.eval() must be called during inference to disable it. Batch normalization normalizes activations per mini-batch and exposes learnable scale/shift parameters; during inference it uses running statistics accumulated during training — another reason model.eval() is critical. Layer normalization is the equivalent for Transformers and RNNs where batch statistics are unreliable. Learning rate schedules (cosine annealing, warmup, ReduceLROnPlateau) are often as important as the choice of optimizer. Mixed precision training (FP16/BF16 via torch.cuda.amp) halves memory and significantly accelerates GPU training with minimal accuracy cost.

---

**What's Next →** [Chapter 21: Project — Image Classifier with FastAPI Deployment](./21-project-image-classifier.md)

*All the techniques in chapters 15–20 come together in a single, complete, production-grade project. We build an image classifier with ResNet50 transfer learning, train it end-to-end with all the regularization and training tricks, evaluate it with Grad-CAM, and deploy it as a REST API with FastAPI.*
