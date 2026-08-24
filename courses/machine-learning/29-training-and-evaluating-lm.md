# Chapter 29: Training and Evaluating Your Language Model

> "A model that fits the training data perfectly but fails on new text hasn't learned language — it has memorized noise."

---

## Table of Contents

1. [Before You Start](#before-you-start)
2. [The Training Loop Deep Dive](#1-the-training-loop-deep-dive)
3. [Understanding Loss Curves](#2-understanding-loss-curves)
4. [Perplexity — The Language Model Metric](#3-perplexity--the-language-model-metric)
5. [Diagnosing Overfitting in LMs](#4-diagnosing-overfitting-in-lms)
6. [Sampling Strategies](#5-sampling-strategies)
7. [Checkpointing and Resuming Training](#6-checkpointing-and-resuming-training)
8. [Evaluating Text Quality](#7-evaluating-text-quality)
9. [Complete Training Script](#8-complete-training-script)
10. [Mini Projects](#mini-projects)
11. [Exercises](#exercises)

---

## Before You Start

### What You Need

- Python 3.9+ with PyTorch 2.0+
- A basic TinyGPT model (from Chapter 28, or use the scaffold below)
- ~2 GB RAM (CPU training is fine for tiny models)

### Quick Setup

```bash
pip install torch numpy matplotlib tqdm
```

### What We Will Build

By the end of this chapter you will have:

- A complete training loop that works for any character-level or token-level language model
- Tools to visualize and diagnose your training runs
- Five different text-generation strategies with working code
- A checkpoint manager that saves your best model automatically

### The Mental Model

Think of training a language model as teaching someone to predict the next word in a sentence. Every time they guess wrong you correct them (compute loss), and they adjust their thinking slightly (backprop + optimizer step). After millions of corrections they get quite good. This chapter is about making that process efficient, measurable, and debuggable.

---

## 1. The Training Loop Deep Dive

### 1.1 From Raw Text to Training Batches

Before any gradient math happens, we need to slice our text corpus into input/target pairs. The key insight: the target is just the input shifted by one position.

```
Raw text:  "to be or not to be"
Tokens:     [t][o][ ][b][e][ ][o][r][ ][n][o][t]

Sequence of length 8:
  Input:   [t][o][ ][b][e][ ][o][r]
  Target:  [o][ ][b][e][ ][o][r][ ]

At every position i:
  given tokens[0..i], predict tokens[i+1]
```

ASCII diagram of how we slice a long text into overlapping windows:

```
Full corpus (length T):
[tok0][tok1][tok2][tok3][tok4][tok5][tok6][tok7][tok8][tok9]...

block_size = 4, batch_size = 2

Random offsets chosen: 2 and 6

Batch item 0 (offset=2):
  Input:  [tok2][tok3][tok4][tok5]
  Target: [tok3][tok4][tok5][tok6]

Batch item 1 (offset=6):
  Input:  [tok6][tok7][tok8][tok9]
  Target: [tok7][tok8][tok9][tok10]

Shape of x: (batch_size=2, block_size=4)
Shape of y: (batch_size=2, block_size=4)
```

### 1.2 Batch Generator Code

```python
import torch
import torch.nn as nn
import numpy as np

class TextDataset:
    """
    Converts a raw text string into batches of (input, target) pairs
    for language model training.
    """
    def __init__(self, text: str, block_size: int, vocab=None):
        # Build character vocabulary if not provided
        if vocab is None:
            chars = sorted(set(text))
            self.stoi = {ch: i for i, ch in enumerate(chars)}
            self.itos = {i: ch for ch, i in self.stoi.items()}
        else:
            self.stoi, self.itos = vocab

        self.vocab_size = len(self.stoi)
        self.block_size = block_size

        # Encode full text as integer tensor
        data = [self.stoi[c] for c in text if c in self.stoi]
        self.data = torch.tensor(data, dtype=torch.long)

        print(f"Dataset: {len(self.data):,} tokens | vocab_size={self.vocab_size}")

    def get_batch(self, batch_size: int, device: str = 'cpu'):
        """
        Sample random (input, target) pairs from the dataset.
        Returns tensors of shape (batch_size, block_size).
        """
        n = len(self.data) - self.block_size
        offsets = torch.randint(0, n, (batch_size,))

        x = torch.stack([self.data[i : i + self.block_size]     for i in offsets])
        y = torch.stack([self.data[i+1 : i + self.block_size+1] for i in offsets])

        return x.to(device), y.to(device)

    def encode(self, text: str) -> torch.Tensor:
        return torch.tensor([self.stoi[c] for c in text if c in self.stoi],
                            dtype=torch.long)

    def decode(self, tokens) -> str:
        if isinstance(tokens, torch.Tensor):
            tokens = tokens.tolist()
        return ''.join(self.itos.get(t, '?') for t in tokens)


# Demo
sample_text = "hello world hello deep learning hello"
dataset = TextDataset(sample_text, block_size=8)

x, y = dataset.get_batch(batch_size=2)
print("Input  x:", x)
print("Target y:", y)
print("Input text  :", dataset.decode(x[0]))
print("Target text :", dataset.decode(y[0]))
```

### 1.3 Teacher Forcing

When training sequence models, we use **teacher forcing**: at every step, we feed the *ground-truth* previous token as input, not the model's own (potentially wrong) prediction.

```
Without teacher forcing (inference mode):
  Step 1: model sees "to"  → predicts "b" (wrong! should be " ")
  Step 2: model sees "b"   → predicts "e" (now everything is off...)
  Error compounds over the sequence.

With teacher forcing (training mode):
  Step 1: model sees "to"  → predicts " " (correct token provided)
  Step 2: model sees "to " → predicts "b" (correct token provided)
  Each step trains independently — no error accumulation.

Teacher forcing lets us compute loss at ALL positions in parallel:
  logits shape:  (batch, seq_len, vocab_size)
  targets shape: (batch, seq_len)
  loss = cross_entropy(logits.view(-1, V), targets.view(-1))
```

### 1.4 Cross-Entropy Loss Calculation

```python
def compute_loss(model, x, y):
    """
    Forward pass + loss computation.
    
    x, y: (batch_size, seq_len) integer tensors
    Returns: scalar loss tensor
    """
    # Forward pass: get logits for every position
    logits = model(x)           # shape: (B, T, vocab_size)

    B, T, V = logits.shape

    # Reshape for cross-entropy:
    #   logits: (B*T, V) — one distribution per token
    #   targets: (B*T,)  — one ground-truth per token
    logits_flat  = logits.view(B * T, V)
    targets_flat = y.view(B * T)

    loss = nn.functional.cross_entropy(logits_flat, targets_flat)
    return loss


# Cross-entropy intuition:
# loss = -log(probability assigned to the correct token)
# 
# If the model assigns prob=0.9 to the correct token: loss = -log(0.9) ≈ 0.105  (good)
# If the model assigns prob=0.1 to the correct token: loss = -log(0.1) ≈ 2.303  (bad)
# Random model over 65 chars:  loss = -log(1/65) ≈ 4.17  (baseline)
```

### 1.5 The Complete Training Step

```python
import torch.optim as optim

def training_step(model, optimizer, dataset, batch_size, device):
    """
    One full training step: sample batch → forward → loss → backward → update.
    Returns scalar loss value.
    """
    model.train()

    # 1. Sample a batch
    x, y = dataset.get_batch(batch_size, device)

    # 2. Zero gradients from the previous step
    optimizer.zero_grad()

    # 3. Forward pass
    loss = compute_loss(model, x, y)

    # 4. Backward pass — compute gradients
    loss.backward()

    # 5. Gradient clipping (prevents exploding gradients)
    torch.nn.utils.clip_grad_norm_(model.parameters(), max_norm=1.0)

    # 6. Optimizer step — update weights
    optimizer.step()

    return loss.item()


@torch.no_grad()
def estimate_loss(model, dataset, batch_size, eval_iters, device):
    """
    Estimate loss on train/val split by averaging over multiple batches.
    Using @no_grad() saves memory and speeds things up during evaluation.
    """
    model.eval()
    losses = torch.zeros(eval_iters)
    for k in range(eval_iters):
        x, y = dataset.get_batch(batch_size, device)
        loss = compute_loss(model, x, y)
        losses[k] = loss.item()
    return losses.mean().item()
```

---

## 2. Understanding Loss Curves

The loss curve is your window into what the model is learning. You should log loss every N steps and plot it. Here's what different shapes mean:

### 2.1 Healthy Training

```
Loss
 5 |*
   | *
 4 |  **
   |    **
 3 |      ***
   |         ***
 2 |            ****
   |                *****
 1 |                     ********
   +---------------------------------> Steps
     0    2k   4k   6k   8k   10k

  Both train loss (solid) and val loss (dashed) decrease together.
  Val loss slightly higher than train loss — normal and healthy.
  The gap between them stays roughly constant — no overfitting.
```

### 2.2 Overfitting

```
Loss
 5 |*
   | *
 4 |  **  <- val loss stops improving
   |    **...........*****  <- val loss starts rising
 3 |      ***              \
   |         ****           \  <- GAP GROWING (overfitting)
 2 |             *****
   |                  *******  <- train loss keeps falling
 1 |
   +---------------------------------> Steps
     0    2k   4k   6k   8k   10k

  Train loss (bottom line) keeps decreasing.
  Val loss (top line) reaches a minimum then increases.
  Solution: add dropout, reduce model size, or stop early.
```

### 2.3 Loss Divergence

```
Loss
 5 |                    /
   |                   /
 4 |        /\/\/\    /
   |       /       \ /
 3 |      /         V   <- temporary crash
   |  /\/
 2 | /
   |/  <- starts fine
 1 +---------------------------------> Steps
     0    2k   4k   6k   8k   10k

  Loss oscillates wildly or spikes upward.
  Causes: learning rate too high, exploding gradients,
          bad data (NaN tokens), corrupted checkpoint.
  Solution: lower LR, add gradient clipping, check data.
```

### 2.4 Loss Plateauing

```
Loss
 5 |*
   | *
 4 |  ***
   |     *****
 3 |          ************************************
   |
 2 |
   +---------------------------------> Steps
     0    2k   4k   6k   8k   10k

  Loss drops fast at first then completely flattens.
  Causes: learning rate too low, model capacity too small,
          saturated activations, bad initialization.
  Solution: increase LR (or use LR scheduler), larger model.
```

### 2.5 Loss Logging Code

```python
import matplotlib.pyplot as plt
from collections import defaultdict

class LossLogger:
    def __init__(self):
        self.history = defaultdict(list)

    def log(self, step, **kwargs):
        """Log any named metric at a given step."""
        for name, value in kwargs.items():
            self.history[name].append((step, value))

    def plot(self, title="Training Curves", save_path=None):
        fig, ax = plt.subplots(figsize=(10, 5))

        colors = {'train_loss': 'blue', 'val_loss': 'orange'}
        for name, values in self.history.items():
            steps, losses = zip(*values)
            color = colors.get(name, None)
            ax.plot(steps, losses, label=name, color=color, linewidth=2)

        ax.set_xlabel("Training Steps")
        ax.set_ylabel("Cross-Entropy Loss")
        ax.set_title(title)
        ax.legend()
        ax.grid(True, alpha=0.3)

        if save_path:
            plt.savefig(save_path, dpi=150, bbox_inches='tight')
        plt.show()

    def detect_issues(self):
        """Simple heuristics to flag training problems."""
        issues = []
        if 'train_loss' in self.history and 'val_loss' in self.history:
            train_steps, train_losses = zip(*self.history['train_loss'])
            val_steps, val_losses = zip(*self.history['val_loss'])

            # Check for divergence
            if len(train_losses) > 5:
                recent_train = train_losses[-5:]
                if recent_train[-1] > recent_train[0] * 1.5:
                    issues.append("WARNING: Train loss is increasing — possible divergence")

            # Check for overfitting
            if len(val_losses) > 3:
                mid = len(val_losses) // 2
                if val_losses[-1] > min(val_losses) * 1.05:
                    issues.append("WARNING: Val loss is rising — possible overfitting")

            # Check for plateau
            if len(train_losses) > 10:
                recent = train_losses[-10:]
                spread = max(recent) - min(recent)
                if spread < 0.01:
                    issues.append("INFO: Train loss has plateaued — consider adjusting LR")

        return issues
```

---

## 3. Perplexity — The Language Model Metric

### 3.1 What Is Perplexity?

Cross-entropy loss is great for training, but hard to interpret. What does loss=2.3 *mean*? Perplexity makes this intuitive.

```
Perplexity (PPL) = exp(loss)

If loss = 0.0  → PPL = 1    → perfect: always predicts the right token
If loss = 2.3  → PPL = 10   → as uncertain as choosing from 10 equally likely words
If loss = 4.17 → PPL = 65   → as uncertain as random guessing over 65 characters
If loss = 6.9  → PPL = 1000 → terrible: highly uncertain model

Intuition:
  "A model with perplexity 100 is, on average, as uncertain
   as if it were choosing uniformly among 100 equally likely options."

For reference:
  Good char-level LM on Shakespeare: PPL ≈ 5-15
  GPT-2 on WikiText-103 (word-level): PPL ≈ 17-29
  Human-level on Penn Treebank:       PPL ≈ 70-75 (we do better!)
```

### 3.2 Perplexity Formula

```
Given a sequence of N tokens: t₁, t₂, ..., tₙ

PPL = exp( -1/N * Σᵢ log P(tᵢ | t₁,...,tᵢ₋₁) )
         [------ this is cross-entropy loss -----]

So:  PPL = exp(cross_entropy_loss)
```

### 3.3 Computing Perplexity

```python
import math

def compute_perplexity(model, dataset, batch_size=32,
                       eval_iters=20, device='cpu'):
    """
    Estimate model perplexity on the dataset.
    
    Returns both loss and perplexity for easy interpretation.
    """
    avg_loss = estimate_loss(model, dataset, batch_size, eval_iters, device)
    perplexity = math.exp(avg_loss)
    return avg_loss, perplexity


def perplexity_on_text(model, dataset, text, device='cpu'):
    """
    Compute exact perplexity on a specific text string.
    Useful for evaluation on a fixed test set.
    """
    model.eval()
    tokens = dataset.encode(text).unsqueeze(0).to(device)  # (1, T)

    total_loss = 0.0
    count = 0

    # Different chapters' GPTConfig use different field names for this
    # (max_seq_len in Chapter 28/31, max_seq_len here in Chapter 29's own
    # scaffold) — check both rather than assuming one.
    cfg = getattr(model, 'config', None)
    block_size = getattr(cfg, 'max_seq_len', None) or getattr(cfg, 'block_size', None) or 128

    with torch.no_grad():
        for i in range(0, len(tokens[0]) - 1, block_size):
            x = tokens[:, i : i + block_size]
            y = tokens[:, i+1 : i + block_size + 1]

            if x.shape[1] == 0:
                break

            logits = model(x)
            B, T, V = logits.shape

            # Trim y to match (in case we hit end of text)
            y = y[:, :T]

            loss = nn.functional.cross_entropy(
                logits.view(B*T, V), y.view(B*T)
            )
            total_loss += loss.item() * T
            count += T

    avg_loss = total_loss / count
    return avg_loss, math.exp(avg_loss)


# Usage example
# loss, ppl = compute_perplexity(model, val_dataset)
# print(f"Validation | Loss: {loss:.4f} | Perplexity: {ppl:.2f}")
```

### 3.4 Perplexity as a Training Signal

```python
class PerplexityTracker:
    """Track perplexity over training and flag improvements."""

    def __init__(self):
        self.best_ppl = float('inf')
        self.history = []

    def update(self, step, loss):
        ppl = math.exp(loss)
        self.history.append((step, ppl))
        improved = ppl < self.best_ppl
        if improved:
            self.best_ppl = ppl
        return ppl, improved

    def report(self, step, train_loss, val_loss):
        train_ppl = math.exp(train_loss)
        val_ppl   = math.exp(val_loss)
        val_ppl, improved = self.update(step, val_loss)
        marker = " <-- BEST" if improved else ""
        print(f"Step {step:6d} | "
              f"Train Loss: {train_loss:.4f} (PPL: {train_ppl:7.2f}) | "
              f"Val Loss: {val_loss:.4f} (PPL: {val_ppl:7.2f}){marker}")
```

---

## 4. Diagnosing Overfitting in LMs

### 4.1 The Overfitting Signature

Overfitting in language models looks like this:

```
Step  | Train Loss | Val Loss | Train PPL | Val PPL
------|------------|----------|-----------|--------
 1000 |   3.12     |  3.18    |   22.6    |  24.0    <- healthy gap
 2000 |   2.54     |  2.63    |   12.7    |  13.9    <- still OK
 3000 |   2.11     |  2.27    |    8.2    |   9.7    <- gap growing
 4000 |   1.78     |  2.31    |    5.9    |  10.1    <- val PPL going UP
 5000 |   1.42     |  2.48    |    4.1    |  11.9    <- definite overfit
```

### 4.2 Train/Val Split

```python
def make_train_val_split(text: str, val_fraction: float = 0.1):
    """
    Split text into train and validation portions.
    We use a contiguous split (not random) to preserve context.
    """
    n = len(text)
    split_idx = int(n * (1 - val_fraction))
    train_text = text[:split_idx]
    val_text   = text[split_idx:]
    print(f"Train: {len(train_text):,} chars | Val: {len(val_text):,} chars")
    return train_text, val_text


# Usage:
# train_text, val_text = make_train_val_split(full_text)
# train_dataset = TextDataset(train_text, block_size=128)
# val_dataset   = TextDataset(val_text, block_size=128,
#                             vocab=(train_dataset.stoi, train_dataset.itos))
```

### 4.3 Early Stopping

```python
class EarlyStopping:
    """
    Stop training when validation loss stops improving.
    Prevents wasting compute and reduces overfitting.
    """
    def __init__(self, patience: int = 5, min_delta: float = 0.001):
        self.patience  = patience
        self.min_delta = min_delta
        self.best_loss = float('inf')
        self.counter   = 0
        self.should_stop = False

    def step(self, val_loss: float) -> bool:
        """Call after each validation. Returns True if training should stop."""
        if val_loss < self.best_loss - self.min_delta:
            self.best_loss = val_loss
            self.counter   = 0
        else:
            self.counter += 1
            print(f"EarlyStopping: no improvement for {self.counter}/{self.patience} checks")
            if self.counter >= self.patience:
                print("EarlyStopping: stopping training.")
                self.should_stop = True
        return self.should_stop
```

### 4.4 Common Overfitting Remedies for LMs

```python
# 1. Dropout — randomly zero activations during training
#    Already built into most transformer implementations.
#    Typical values: 0.1 (small model) to 0.3 (large model)

# 2. Weight decay — L2 regularization via optimizer
optimizer = optim.AdamW(
    model.parameters(),
    lr=3e-4,
    weight_decay=0.1    # key: AdamW applies decay correctly
)

# 3. Dataset augmentation — if data is tiny, repeat with noise
#    Not common for LMs, but worth knowing about.

# 4. Reduce model size — fewer parameters = less capacity to overfit
#    Try halving n_embd or n_layer first.

# 5. Learning rate scheduling — cosine decay helps generalization
def get_cosine_lr(step, warmup_steps, max_steps, max_lr, min_lr=1e-5):
    if step < warmup_steps:
        return max_lr * step / warmup_steps
    if step > max_steps:
        return min_lr
    progress = (step - warmup_steps) / (max_steps - warmup_steps)
    return min_lr + 0.5 * (max_lr - min_lr) * (1 + math.cos(math.pi * progress))
```

---

## 5. Sampling Strategies

Once your model is trained, how do you generate text? There are several strategies, each with different tradeoffs between coherence and creativity.

### 5.1 The Sampling Setup

```python
import torch.nn.functional as F

def generate(model, dataset, prompt: str, max_new_tokens: int,
             strategy: str = 'greedy', **kwargs):
    """
    Master generation function. Dispatches to specific strategies.
    
    strategy: 'greedy' | 'temperature' | 'top_k' | 'top_p'
    kwargs: temperature, k, p (depending on strategy)
    """
    model.eval()
    device = next(model.parameters()).device

    # Encode prompt
    tokens = dataset.encode(prompt).unsqueeze(0).to(device)  # (1, T)

    block_size = 128  # adjust to your model's context length

    with torch.no_grad():
        for _ in range(max_new_tokens):
            # Crop context to block_size
            context = tokens[:, -block_size:]

            # Get logits for the last token position
            logits = model(context)          # (1, T, vocab_size)
            next_logits = logits[:, -1, :]   # (1, vocab_size)

            # Sample next token using chosen strategy
            if strategy == 'greedy':
                next_token = greedy_decode(next_logits)
            elif strategy == 'temperature':
                next_token = temperature_sample(next_logits, **kwargs)
            elif strategy == 'top_k':
                next_token = top_k_sample(next_logits, **kwargs)
            elif strategy == 'top_p':
                next_token = top_p_sample(next_logits, **kwargs)
            else:
                raise ValueError(f"Unknown strategy: {strategy}")

            # Append to sequence
            tokens = torch.cat([tokens, next_token], dim=1)

    return dataset.decode(tokens[0])
```

### 5.2 Greedy Decoding

Always picks the single most likely next token. Fast, deterministic, but often repetitive.

```
Logits: [-1.2,  0.5,  2.3, -0.8,  1.1]
         tok0  tok1  tok2  tok3  tok4

Greedy always picks tok2 (highest logit).

Problem: "the the the the the..." — repetition loops
```

```python
def greedy_decode(logits):
    """
    Pick the token with the highest probability.
    Returns shape (1, 1) for easy concatenation.
    """
    return logits.argmax(dim=-1, keepdim=True)  # (1, 1)


# Example output from a character model:
# Prompt: "To be"
# Greedy: "To be or not to be or not to be or not to be..."
#          ^--- repetition loop after "to be"
```

### 5.3 Temperature Scaling

Temperature T divides the logits before softmax, sharpening or flattening the distribution.

```
Logits before temperature: [1.0,  3.0,  2.0,  0.5]

T = 0.5 (low temperature — more conservative):
  Scaled logits: [2.0, 6.0, 4.0, 1.0]  (divided by 0.5 = multiplied)
  Probabilities: [0.02, 0.87, 0.10, 0.01]  <- very peaked, high confidence

T = 1.0 (no change):
  Probabilities: [0.13, 0.52, 0.27, 0.08]  <- standard

T = 2.0 (high temperature — more creative/random):
  Scaled logits: [0.5, 1.5, 1.0, 0.25]
  Probabilities: [0.17, 0.38, 0.26, 0.14]  <- flatter, more random

Analogy:
  T=0.1 → the model always says the "safe" word (like an anxious speaker)
  T=1.0 → normal sampling
  T=1.5 → the model gets adventurous (like free-form poetry)
  T=5.0 → the model is almost incoherent (like random noise)
```

```python
def temperature_sample(logits, temperature: float = 1.0):
    """
    Apply temperature scaling and sample.
    temperature > 1: more random/creative
    temperature < 1: more focused/conservative
    temperature = 1: standard sampling
    """
    if temperature <= 0:
        return greedy_decode(logits)

    # Scale logits
    scaled = logits / temperature

    # Convert to probabilities
    probs = F.softmax(scaled, dim=-1)  # (1, vocab_size)

    # Sample one token
    return torch.multinomial(probs, num_samples=1)  # (1, 1)
```

### 5.4 Top-k Sampling

Keep only the k most likely tokens, zero out the rest, then sample.

```
Top-k with k=3:

All logits: [0.1, 2.3, 1.8, 0.5, 3.1, 0.2, 1.2, 0.8]
            tok0 tok1 tok2 tok3 tok4 tok5 tok6 tok7

Top-3 tokens (by value): tok4(3.1), tok1(2.3), tok2(1.8)

Set all others to -infinity:
           [-inf, 2.3, 1.8, -inf, 3.1, -inf, -inf, -inf]

Now softmax and sample only from {tok1, tok2, tok4}.

Effect: removes the long tail of low-probability tokens
        (which often produce incoherent/garbage tokens)
```

```python
def top_k_sample(logits, k: int = 50, temperature: float = 1.0):
    """
    Keep only the top-k tokens, zero out the rest, then sample.
    """
    # Apply temperature first
    if temperature != 1.0:
        logits = logits / temperature

    # Find the k-th largest value
    top_k_values, _ = torch.topk(logits, k=min(k, logits.shape[-1]))
    threshold = top_k_values[:, -1].unsqueeze(-1)  # k-th value

    # Zero out (set to -inf) any token below the threshold
    filtered = logits.masked_fill(logits < threshold, float('-inf'))

    # Sample
    probs = F.softmax(filtered, dim=-1)
    return torch.multinomial(probs, num_samples=1)
```

### 5.5 Top-p (Nucleus) Sampling

Instead of a fixed k, keep the smallest set of tokens whose cumulative probability exceeds p.

```
Top-p with p=0.9:

Sorted probabilities (descending):
  tok4: 0.45  | cumsum: 0.45  <- include
  tok1: 0.25  | cumsum: 0.70  <- include
  tok2: 0.15  | cumsum: 0.85  <- include
  tok6: 0.08  | cumsum: 0.93  <- include (pushes us over 0.9)
  tok7: 0.04  | cumsum: 0.97  <- EXCLUDE
  tok3: 0.02  | cumsum: 0.99  <- EXCLUDE
  tok0: 0.01  | cumsum: 1.00  <- EXCLUDE

We sample only from {tok4, tok1, tok2, tok6}.

Key advantage over top-k:
  - When the model is confident, nucleus is small (sharp peak → 1-2 tokens)
  - When the model is uncertain, nucleus is larger (flat → many tokens)
  Top-k always uses the same number; top-p adapts dynamically.
```

```python
def top_p_sample(logits, p: float = 0.9, temperature: float = 1.0):
    """
    Nucleus sampling: keep the smallest set of tokens with cumulative
    probability >= p, then sample from that set.
    """
    # Apply temperature
    if temperature != 1.0:
        logits = logits / temperature

    # Sort logits in descending order
    sorted_logits, sorted_indices = torch.sort(logits, descending=True, dim=-1)
    cumulative_probs = torch.cumsum(F.softmax(sorted_logits, dim=-1), dim=-1)

    # Remove tokens that push cumulative probability above p
    # Shift right by 1 so we always keep at least the top token
    remove_mask = cumulative_probs - F.softmax(sorted_logits, dim=-1) > p
    sorted_logits[remove_mask] = float('-inf')

    # Unsort back to original order
    unsorted = sorted_logits.scatter(1, sorted_indices, sorted_logits)

    probs = F.softmax(unsorted, dim=-1)
    return torch.multinomial(probs, num_samples=1)


# Practical recommendation:
#   Most systems combine top-p + temperature:
#   temperature=0.8, top_p=0.9 gives good results on most tasks.
```

### 5.6 Comparing Sampling Strategies

```python
def compare_sampling_strategies(model, dataset, prompt,
                                  max_new_tokens=100, device='cpu'):
    """
    Generate text with all four strategies and print side by side.
    """
    strategies = [
        ('Greedy',              dict(strategy='greedy')),
        ('Temperature T=0.7',   dict(strategy='temperature', temperature=0.7)),
        ('Temperature T=1.4',   dict(strategy='temperature', temperature=1.4)),
        ('Top-k k=40',          dict(strategy='top_k',       k=40)),
        ('Top-p p=0.9',         dict(strategy='top_p',       p=0.9)),
    ]

    print(f"Prompt: '{prompt}'\n" + "="*60)
    for name, kwargs in strategies:
        result = generate(model, dataset, prompt, max_new_tokens, **kwargs)
        generated = result[len(prompt):]  # show only the new tokens
        print(f"\n[{name}]")
        print(generated)
        print("-"*60)
```

---

## 6. Checkpointing and Resuming Training

### 6.1 Why Checkpointing Matters

Training can crash. GPU can overheat. You might accidentally kill the process. Checkpointing saves your progress so you can resume without starting over. It also lets you track the *best* model (by validation loss) separately from the *latest* model.

```
Checkpoint directory structure:

checkpoints/
├── latest.pt          <- most recent model (for resuming)
├── best.pt            <- best validation loss model (for deployment)
├── step_001000.pt     <- periodic snapshots
├── step_002000.pt
└── training_log.csv   <- loss history
```

### 6.2 Saving a Checkpoint

```python
import os

def save_checkpoint(model, optimizer, step, train_loss, val_loss,
                    config, save_dir='checkpoints', tag='latest'):
    """
    Save a complete training checkpoint.
    Includes model weights, optimizer state, and metadata.
    """
    os.makedirs(save_dir, exist_ok=True)
    path = os.path.join(save_dir, f'{tag}.pt')

    checkpoint = {
        # Model state
        'model_state_dict': model.state_dict(),

        # Optimizer state (crucial for resuming correctly)
        'optimizer_state_dict': optimizer.state_dict(),

        # Training progress
        'step': step,
        'train_loss': train_loss,
        'val_loss': val_loss,

        # Model configuration (so we can recreate the architecture)
        'config': config,
    }

    torch.save(checkpoint, path)
    print(f"Saved checkpoint to {path} | step={step} | val_loss={val_loss:.4f}")
    return path
```

### 6.3 Loading a Checkpoint

```python
def load_checkpoint(path, model, optimizer=None):
    """
    Load a checkpoint. If optimizer is provided, restore its state too.
    Returns the step to resume from.
    """
    print(f"Loading checkpoint from {path}")
    checkpoint = torch.load(path, map_location='cpu')

    model.load_state_dict(checkpoint['model_state_dict'])

    if optimizer is not None and 'optimizer_state_dict' in checkpoint:
        optimizer.load_state_dict(checkpoint['optimizer_state_dict'])

    step      = checkpoint.get('step', 0)
    val_loss  = checkpoint.get('val_loss', float('inf'))

    print(f"Resumed from step={step} | val_loss={val_loss:.4f}")
    return step, val_loss


# Usage:
# model = TinyGPT(config)
# optimizer = optim.AdamW(model.parameters(), lr=3e-4)
#
# # Resume if checkpoint exists
# if os.path.exists('checkpoints/latest.pt'):
#     start_step, _ = load_checkpoint('checkpoints/latest.pt', model, optimizer)
# else:
#     start_step = 0
```

### 6.4 Best Model Tracker

```python
class BestModelTracker:
    """
    Automatically saves a checkpoint whenever validation loss improves.
    """
    def __init__(self, save_dir='checkpoints', patience=None):
        self.save_dir  = save_dir
        self.best_loss = float('inf')
        self.patience  = patience
        self.no_improve_count = 0
        os.makedirs(save_dir, exist_ok=True)

    def update(self, model, optimizer, step, train_loss, val_loss, config):
        """
        Call after every validation. Returns True if this is a new best.
        """
        if val_loss < self.best_loss:
            self.best_loss = val_loss
            self.no_improve_count = 0
            save_checkpoint(model, optimizer, step, train_loss, val_loss,
                           config, self.save_dir, tag='best')
            print(f"  >> New best model! val_loss={val_loss:.4f}")
            return True
        else:
            self.no_improve_count += 1
            if self.patience and self.no_improve_count >= self.patience:
                print(f"  >> No improvement for {self.patience} checks. Stop.")
                return None  # signal to stop
            return False
```

---

## 7. Evaluating Text Quality

### 7.1 Perplexity as the Primary Metric

For language models, perplexity is the standard metric because it directly measures how well the model predicts held-out text.

```python
def full_evaluation(model, train_dataset, val_dataset,
                    batch_size=32, eval_iters=50, device='cpu'):
    """
    Comprehensive evaluation: perplexity on train and val sets.
    """
    train_loss, train_ppl = compute_perplexity(
        model, train_dataset, batch_size, eval_iters, device)
    val_loss, val_ppl = compute_perplexity(
        model, val_dataset, batch_size, eval_iters, device)

    gap = val_loss - train_loss
    overfit_signal = "OVERFITTING" if gap > 0.5 else "OK"

    print(f"\n{'='*50}")
    print(f"EVALUATION RESULTS")
    print(f"{'='*50}")
    print(f"Train: loss={train_loss:.4f}  ppl={train_ppl:.2f}")
    print(f"Val:   loss={val_loss:.4f}  ppl={val_ppl:.2f}")
    print(f"Gap:   {gap:.4f}  [{overfit_signal}]")
    print(f"{'='*50}\n")

    return {
        'train_loss': train_loss, 'train_ppl': train_ppl,
        'val_loss': val_loss,     'val_ppl': val_ppl,
    }
```

### 7.2 BLEU Score (Brief Introduction)

BLEU (Bilingual Evaluation Understudy) measures overlap between generated text and reference text. It's primarily used for translation and summarization, not pure LM tasks, but it's worth knowing.

```python
from collections import Counter

def simple_bleu_1gram(reference: str, hypothesis: str) -> float:
    """
    1-gram BLEU: fraction of generated words that appear in reference.
    This is a simplified version for educational purposes.
    For production use: pip install nltk; nltk.translate.bleu_score
    """
    ref_words  = reference.split()
    hyp_words  = hypothesis.split()

    if not hyp_words:
        return 0.0

    ref_counts = Counter(ref_words)
    matches = 0
    for word in hyp_words:
        if ref_counts.get(word, 0) > 0:
            matches += 1
            ref_counts[word] -= 1

    return matches / len(hyp_words)


# Note: BLEU is less useful for open-ended generation.
# For character-level models, perplexity is the better metric.
# For instruction-following models: human eval or LLM-as-judge.
```

### 7.3 Qualitative Evaluation

Numbers alone are not enough. Always do a qualitative check:

```python
def qualitative_eval(model, dataset, prompts, max_new_tokens=200, device='cpu'):
    """
    Generate text from multiple prompts and print results.
    Use this alongside quantitative metrics.
    """
    print("\nQUALITATIVE EVALUATION")
    print("="*60)

    test_configs = [
        ('Conservative (T=0.7)',  dict(strategy='temperature', temperature=0.7)),
        ('Standard    (Top-p)',   dict(strategy='top_p',       p=0.9, temperature=0.9)),
        ('Creative    (T=1.3)',   dict(strategy='temperature', temperature=1.3)),
    ]

    for prompt in prompts:
        print(f"\nPrompt: '{prompt}'")
        print("-"*40)
        for name, kwargs in test_configs:
            result = generate(model, dataset, prompt, max_new_tokens, **kwargs)
            generated = result[len(prompt):]
            print(f"[{name}]: {generated[:100]}...")
        print()
```

---

## 8. Complete Training Script

This section ties everything together into a single, runnable training script.

```python
"""
complete_training.py

Full training script for a character-level language model.
Run: python complete_training.py
"""

import os
import math
import time
import torch
import torch.nn as nn
import torch.optim as optim
from dataclasses import dataclass
from typing import Optional


# ─── Minimal TinyGPT (self-contained) ───────────────────────────────────────

@dataclass
class GPTConfig:
    vocab_size:   int = 65
    max_seq_len:   int = 128
    d_model:       int = 128
    n_heads:       int = 4
    n_layers:      int = 4
    dropout:      float = 0.1


class CausalSelfAttention(nn.Module):
    def __init__(self, config):
        super().__init__()
        assert config.d_model % config.n_heads == 0
        self.n_heads  = config.n_heads
        self.d_model  = config.d_model
        self.head_dim = config.d_model // config.n_heads
        self.c_attn  = nn.Linear(config.d_model, 3 * config.d_model)
        self.c_proj  = nn.Linear(config.d_model, config.d_model)
        self.drop    = nn.Dropout(config.dropout)
        self.register_buffer('mask', torch.tril(
            torch.ones(config.max_seq_len, config.max_seq_len)
        ).view(1, 1, config.max_seq_len, config.max_seq_len))

    def forward(self, x):
        B, T, C = x.shape
        q, k, v = self.c_attn(x).split(self.d_model, dim=2)
        q = q.view(B, T, self.n_heads, self.head_dim).transpose(1, 2)
        k = k.view(B, T, self.n_heads, self.head_dim).transpose(1, 2)
        v = v.view(B, T, self.n_heads, self.head_dim).transpose(1, 2)
        att = (q @ k.transpose(-2, -1)) * (self.head_dim ** -0.5)
        att = att.masked_fill(self.mask[:,:,:T,:T] == 0, float('-inf'))
        att = torch.softmax(att, dim=-1)
        att = self.drop(att)
        out = (att @ v).transpose(1, 2).contiguous().view(B, T, C)
        return self.c_proj(out)


class MLP(nn.Module):
    def __init__(self, config):
        super().__init__()
        self.net = nn.Sequential(
            nn.Linear(config.d_model, 4 * config.d_model),
            nn.GELU(),
            nn.Linear(4 * config.d_model, config.d_model),
            nn.Dropout(config.dropout),
        )
    def forward(self, x):
        return self.net(x)


class Block(nn.Module):
    def __init__(self, config):
        super().__init__()
        self.ln1  = nn.LayerNorm(config.d_model)
        self.attn = CausalSelfAttention(config)
        self.ln2  = nn.LayerNorm(config.d_model)
        self.mlp  = MLP(config)

    def forward(self, x):
        x = x + self.attn(self.ln1(x))
        x = x + self.mlp(self.ln2(x))
        return x


class TinyGPT(nn.Module):
    def __init__(self, config: GPTConfig):
        super().__init__()
        self.config = config
        self.tok_emb = nn.Embedding(config.vocab_size, config.d_model)
        self.pos_emb = nn.Embedding(config.max_seq_len, config.d_model)
        self.drop    = nn.Dropout(config.dropout)
        self.blocks  = nn.Sequential(*[Block(config) for _ in range(config.n_layers)])
        self.ln_f    = nn.LayerNorm(config.d_model)
        self.head    = nn.Linear(config.d_model, config.vocab_size, bias=False)

    def forward(self, idx):
        B, T = idx.shape
        pos = torch.arange(T, device=idx.device)
        x = self.drop(self.tok_emb(idx) + self.pos_emb(pos))
        x = self.blocks(x)
        x = self.ln_f(x)
        return self.head(x)

    def num_params(self):
        return sum(p.numel() for p in self.parameters())


# ─── Training Configuration ─────────────────────────────────────────────────

@dataclass
class TrainConfig:
    # Data
    data_path:       str   = 'data/input.txt'
    val_fraction:    float = 0.1
    max_seq_len:      int   = 128

    # Model
    d_model:          int   = 128
    n_heads:          int   = 4
    n_layers:         int   = 4
    dropout:         float = 0.1

    # Training
    batch_size:      int   = 32
    max_steps:       int   = 5000
    lr:              float = 3e-4
    weight_decay:    float = 0.1
    warmup_steps:    int   = 100
    grad_clip:       float = 1.0

    # Evaluation
    eval_interval:   int   = 200
    eval_iters:      int   = 50
    save_interval:   int   = 500

    # I/O
    save_dir:        str   = 'checkpoints'
    resume:          bool  = True


# ─── Main Training Loop ──────────────────────────────────────────────────────

def train(cfg: TrainConfig):
    device = 'cuda' if torch.cuda.is_available() else 'cpu'
    print(f"Training on: {device}")

    # ── Load data ──
    with open(cfg.data_path, 'r', encoding='utf-8') as f:
        text = f.read()

    n = len(text)
    split = int(n * (1 - cfg.val_fraction))
    train_text = text[:split]
    val_text   = text[split:]

    train_ds = TextDataset(train_text, cfg.max_seq_len)
    val_ds   = TextDataset(val_text,   cfg.max_seq_len,
                           vocab=(train_ds.stoi, train_ds.itos))

    # ── Build model ──
    model_cfg = GPTConfig(
        vocab_size = train_ds.vocab_size,
        max_seq_len = cfg.max_seq_len,
        d_model     = cfg.d_model,
        n_heads     = cfg.n_heads,
        n_layers    = cfg.n_layers,
        dropout    = cfg.dropout,
    )
    model = TinyGPT(model_cfg).to(device)
    print(f"Model: {model.num_params():,} parameters")

    # ── Optimizer ──
    optimizer = optim.AdamW(
        model.parameters(),
        lr=cfg.lr,
        weight_decay=cfg.weight_decay,
        betas=(0.9, 0.95),
    )

    # ── Optional: resume from checkpoint ──
    start_step = 0
    best_tracker = BestModelTracker(cfg.save_dir)
    logger = LossLogger()

    if cfg.resume and os.path.exists(f'{cfg.save_dir}/latest.pt'):
        start_step, _ = load_checkpoint(
            f'{cfg.save_dir}/latest.pt', model, optimizer)

    # ── Training loop ──
    model.train()
    t0 = time.time()

    for step in range(start_step, cfg.max_steps + 1):
        # Update learning rate (cosine schedule with warmup)
        lr = get_cosine_lr(step, cfg.warmup_steps, cfg.max_steps, cfg.lr)
        for param_group in optimizer.param_groups:
            param_group['lr'] = lr

        # Evaluation
        if step % cfg.eval_interval == 0:
            train_loss = estimate_loss(model, train_ds, cfg.batch_size,
                                       cfg.eval_iters, device)
            val_loss   = estimate_loss(model, val_ds, cfg.batch_size,
                                       cfg.eval_iters, device)

            logger.log(step, train_loss=train_loss, val_loss=val_loss)

            elapsed = time.time() - t0
            print(f"step {step:5d} | "
                  f"train_loss: {train_loss:.4f} | "
                  f"val_loss: {val_loss:.4f} | "
                  f"ppl: {math.exp(val_loss):.2f} | "
                  f"lr: {lr:.2e} | "
                  f"elapsed: {elapsed:.1f}s")

            # Save best model
            best_tracker.update(model, optimizer, step, train_loss,
                                val_loss, model_cfg)

        # Periodic checkpoint (for resuming)
        if step % cfg.save_interval == 0 and step > 0:
            train_loss_quick = estimate_loss(model, train_ds, 16, 5, device)
            val_loss_quick   = estimate_loss(model, val_ds,   16, 5, device)
            save_checkpoint(model, optimizer, step,
                           train_loss_quick, val_loss_quick,
                           model_cfg, cfg.save_dir, 'latest')

        if step == cfg.max_steps:
            break

        # ── Training step ──
        x, y = train_ds.get_batch(cfg.batch_size, device)
        optimizer.zero_grad()
        logits = model(x)
        B, T, V = logits.shape
        loss = nn.functional.cross_entropy(logits.view(B*T, V), y.view(B*T))
        loss.backward()
        torch.nn.utils.clip_grad_norm_(model.parameters(), cfg.grad_clip)
        optimizer.step()

    print("\nTraining complete!")

    # Final evaluation
    final_val_loss = estimate_loss(model, val_ds, cfg.batch_size,
                                   cfg.eval_iters, device)
    print(f"Final val loss: {final_val_loss:.4f} | "
          f"Final PPL: {math.exp(final_val_loss):.2f}")

    # Plot loss curves
    logger.plot(save_path=f'{cfg.save_dir}/loss_curves.png')

    # Sample some text
    print("\nSample generations:")
    compare_sampling_strategies(model, train_ds, "To be or", 100, device)

    return model, train_ds


if __name__ == '__main__':
    cfg = TrainConfig()

    # For quick testing, use a built-in text:
    os.makedirs('data', exist_ok=True)
    if not os.path.exists(cfg.data_path):
        sample = """To be or not to be, that is the question.
Whether tis nobler in the mind to suffer
The slings and arrows of outrageous fortune,
Or to take arms against a sea of troubles.
""" * 200
        with open(cfg.data_path, 'w') as f:
            f.write(sample)
        print(f"Created sample data at {cfg.data_path}")

    model, dataset = train(cfg)
```

---

## Mini Projects

### Mini Project 1: Loss Curve Analyzer

Reads a training log CSV, plots curves, and flags potential issues.

```python
"""
loss_curve_analyzer.py

Usage:
  python loss_curve_analyzer.py --log checkpoints/training_log.csv
  python loss_curve_analyzer.py --demo   # generates synthetic data for demo
"""

import argparse
import csv
import math
import matplotlib.pyplot as plt
import matplotlib.patches as mpatches
import numpy as np
import os


def load_log(path: str):
    """Load training log CSV. Expected columns: step, train_loss, val_loss"""
    steps, train_losses, val_losses = [], [], []
    with open(path, 'r') as f:
        reader = csv.DictReader(f)
        for row in reader:
            steps.append(int(row['step']))
            train_losses.append(float(row['train_loss']))
            val_losses.append(float(row.get('val_loss', row['train_loss'])))
    return steps, train_losses, val_losses


def generate_demo_log(scenario: str = 'healthy') -> str:
    """Generate a synthetic training log for demonstration."""
    path = f'/tmp/demo_{scenario}_log.csv'
    steps = list(range(0, 5001, 100))
    rows = [['step', 'train_loss', 'val_loss']]

    for s in steps:
        t = s / 5000
        base = 4.5 * math.exp(-3 * t) + 0.8 + np.random.normal(0, 0.03)

        if scenario == 'healthy':
            train_loss = base
            val_loss   = base + 0.15 + np.random.normal(0, 0.02)
        elif scenario == 'overfit':
            train_loss = base
            val_loss   = (base + 0.15) if t < 0.4 else (1.8 + (t - 0.4) * 3)
        elif scenario == 'diverge':
            train_loss = base if t < 0.3 else base + 5 * (t - 0.3) ** 2
            val_loss   = train_loss + 0.2
        elif scenario == 'plateau':
            train_loss = 4.5 * math.exp(-10 * min(t, 0.15)) + 1.8
            val_loss   = train_loss + 0.15
        else:
            train_loss = base
            val_loss   = base + 0.1

        rows.append([s, round(train_loss, 4), round(max(val_loss, 0.5), 4)])

    with open(path, 'w', newline='') as f:
        writer = csv.writer(f)
        writer.writerows(rows)

    return path


def detect_issues(steps, train_losses, val_losses):
    """Analyze the curves and return a list of issue descriptions."""
    issues = []
    n = len(steps)

    if n < 5:
        return ["Not enough data points to analyze"]

    # 1. Divergence: train loss going up significantly
    recent_train = train_losses[max(0, n-5):]
    if recent_train[-1] > recent_train[0] * 1.3:
        issues.append({
            'type': 'DIVERGENCE',
            'severity': 'HIGH',
            'message': f'Train loss increased {((recent_train[-1]/recent_train[0])-1)*100:.0f}% in last 5 checkpoints',
            'suggestion': 'Lower learning rate, add gradient clipping, or check data for NaN values'
        })

    # 2. Overfitting: val loss rising while train loss falls
    best_val_idx = np.argmin(val_losses)
    if best_val_idx < n - 3:  # best was not near the end
        val_increase = val_losses[-1] - val_losses[best_val_idx]
        if val_increase > 0.1:
            issues.append({
                'type': 'OVERFITTING',
                'severity': 'MEDIUM',
                'message': (f'Val loss at step {steps[best_val_idx]} '
                           f'(best={val_losses[best_val_idx]:.4f}), '
                           f'now {val_losses[-1]:.4f} (+{val_increase:.4f})'),
                'suggestion': 'Increase dropout, reduce model size, or use early stopping'
            })

    # 3. Plateau: train loss not improving
    recent_train_spread = max(recent_train) - min(recent_train)
    if recent_train_spread < 0.02 and train_losses[-1] > 1.0:
        issues.append({
            'type': 'PLATEAU',
            'severity': 'LOW',
            'message': f'Train loss spread in last 5 steps: {recent_train_spread:.4f} (very flat)',
            'suggestion': 'Increase learning rate, use LR scheduler, or increase model capacity'
        })

    # 4. High generalization gap
    final_gap = val_losses[-1] - train_losses[-1]
    if final_gap > 1.0:
        issues.append({
            'type': 'LARGE_GAP',
            'severity': 'MEDIUM',
            'message': f'Val-train loss gap: {final_gap:.4f} (>1.0)',
            'suggestion': 'Model may be overfitting; check train/val split size'
        })

    if not issues:
        issues.append({
            'type': 'HEALTHY',
            'severity': 'OK',
            'message': 'No issues detected. Training looks healthy.',
            'suggestion': 'Continue training or check if you have reached the target loss.'
        })

    return issues


def plot_analysis(steps, train_losses, val_losses, issues, save_path=None):
    """Create a comprehensive training analysis plot."""
    fig, axes = plt.subplots(2, 2, figsize=(14, 10))
    fig.suptitle('Training Analysis Report', fontsize=16, fontweight='bold')

    # Plot 1: Loss curves
    ax1 = axes[0, 0]
    ax1.plot(steps, train_losses, 'b-', label='Train Loss', linewidth=2)
    ax1.plot(steps, val_losses,   'r--', label='Val Loss',   linewidth=2)
    ax1.set_xlabel('Step')
    ax1.set_ylabel('Cross-Entropy Loss')
    ax1.set_title('Loss Curves')
    ax1.legend()
    ax1.grid(True, alpha=0.3)

    # Mark best val
    best_idx = np.argmin(val_losses)
    ax1.axvline(x=steps[best_idx], color='green', linestyle=':', alpha=0.7,
               label=f'Best val @ step {steps[best_idx]}')
    ax1.legend()

    # Plot 2: Perplexity curves
    ax2 = axes[0, 1]
    train_ppl = [math.exp(min(l, 10)) for l in train_losses]
    val_ppl   = [math.exp(min(l, 10)) for l in val_losses]
    ax2.plot(steps, train_ppl, 'b-', label='Train PPL', linewidth=2)
    ax2.plot(steps, val_ppl,   'r--', label='Val PPL',   linewidth=2)
    ax2.set_xlabel('Step')
    ax2.set_ylabel('Perplexity')
    ax2.set_title('Perplexity (lower=better)')
    ax2.legend()
    ax2.grid(True, alpha=0.3)

    # Plot 3: Generalization gap
    ax3 = axes[1, 0]
    gap = [v - t for v, t in zip(val_losses, train_losses)]
    ax3.fill_between(steps, gap, alpha=0.4, color='orange', label='Val-Train Gap')
    ax3.plot(steps, gap, 'darkorange', linewidth=1.5)
    ax3.axhline(y=0, color='black', linestyle='-', linewidth=0.5)
    ax3.set_xlabel('Step')
    ax3.set_ylabel('Loss Gap')
    ax3.set_title('Generalization Gap (val - train)')
    ax3.legend()
    ax3.grid(True, alpha=0.3)

    # Plot 4: Issues summary
    ax4 = axes[1, 1]
    ax4.axis('off')
    ax4.set_title('Detected Issues', fontweight='bold')

    severity_colors = {'HIGH': 'red', 'MEDIUM': 'orange', 'LOW': 'gold', 'OK': 'green'}
    y_pos = 0.95
    for issue in issues:
        color = severity_colors.get(issue['severity'], 'black')
        ax4.text(0.05, y_pos, f"[{issue['severity']}] {issue['type']}",
                transform=ax4.transAxes, fontsize=10,
                fontweight='bold', color=color, va='top')
        y_pos -= 0.08

        # Wrap message
        msg = issue['message']
        ax4.text(0.05, y_pos, f"  {msg[:70]}",
                transform=ax4.transAxes, fontsize=8, color='black', va='top')
        y_pos -= 0.07

        ax4.text(0.05, y_pos, f"  Tip: {issue['suggestion'][:60]}",
                transform=ax4.transAxes, fontsize=8, color='gray', va='top',
                style='italic')
        y_pos -= 0.10

        if y_pos < 0.1:
            break

    plt.tight_layout()
    if save_path:
        plt.savefig(save_path, dpi=150, bbox_inches='tight')
        print(f"Saved analysis to {save_path}")
    plt.show()


def analyze(log_path: str, output_path: str = None):
    steps, train_losses, val_losses = load_log(log_path)
    issues = detect_issues(steps, train_losses, val_losses)

    print("\n" + "="*60)
    print("LOSS CURVE ANALYSIS")
    print("="*60)
    print(f"Total steps logged: {len(steps)}")
    print(f"Final train loss:   {train_losses[-1]:.4f}  (PPL: {math.exp(min(train_losses[-1],10)):.2f})")
    print(f"Final val loss:     {val_losses[-1]:.4f}    (PPL: {math.exp(min(val_losses[-1],10)):.2f})")
    print(f"Best val loss:      {min(val_losses):.4f}    @ step {steps[int(np.argmin(val_losses))]}")

    print("\nISSUES:")
    for issue in issues:
        print(f"  [{issue['severity']:6s}] {issue['type']:15s} | {issue['message']}")
        print(f"           Suggestion: {issue['suggestion']}")
    print("="*60)

    out = output_path or log_path.replace('.csv', '_analysis.png')
    plot_analysis(steps, train_losses, val_losses, issues, out)


if __name__ == '__main__':
    parser = argparse.ArgumentParser(description='Analyze training loss curves')
    parser.add_argument('--log',  type=str, help='Path to training log CSV')
    parser.add_argument('--demo', type=str, default='healthy',
                       choices=['healthy', 'overfit', 'diverge', 'plateau'],
                       help='Run with synthetic demo data')
    args = parser.parse_args()

    if args.log:
        analyze(args.log)
    else:
        print(f"Running demo scenario: '{args.demo}'")
        demo_path = generate_demo_log(args.demo)
        analyze(demo_path, f'/tmp/analysis_{args.demo}.png')
```

---

### Mini Project 2: Sampling Playground

Compare all sampling strategies side by side with a simple interactive interface.

```python
"""
sampling_playground.py

Generates text from a trained model using all sampling strategies
and displays them side by side for easy comparison.

Usage: python sampling_playground.py --checkpoint checkpoints/best.pt
"""

import argparse
import torch

def load_model_for_inference(checkpoint_path, device='cpu'):
    """Load a saved TinyGPT model for inference."""
    ckpt = torch.load(checkpoint_path, map_location=device)
    config = ckpt['config']
    model = TinyGPT(config).to(device)
    model.load_state_dict(ckpt['model_state_dict'])
    model.eval()
    return model, config


def sampling_playground(model, dataset, device='cpu'):
    """Interactive sampling playground."""
    print("\n" + "="*70)
    print("SAMPLING PLAYGROUND")
    print("Compare text generation strategies interactively")
    print("="*70)

    configs = [
        ('1', 'Greedy',             {'strategy': 'greedy'}),
        ('2', 'Temp=0.5 (safe)',    {'strategy': 'temperature', 'temperature': 0.5}),
        ('3', 'Temp=0.8 (balanced)',{'strategy': 'temperature', 'temperature': 0.8}),
        ('4', 'Temp=1.2 (creative)',{'strategy': 'temperature', 'temperature': 1.2}),
        ('5', 'Top-k=10',           {'strategy': 'top_k',  'k': 10}),
        ('6', 'Top-k=50',           {'strategy': 'top_k',  'k': 50}),
        ('7', 'Top-p=0.9',          {'strategy': 'top_p',  'p': 0.9}),
        ('8', 'Top-p + Temp=0.8',   {'strategy': 'top_p',  'p': 0.9, 'temperature': 0.8}),
        ('A', 'ALL strategies',     None),
    ]

    while True:
        print("\nChoose a strategy:")
        for key, name, _ in configs:
            print(f"  [{key}] {name}")
        print("  [Q] Quit")

        choice = input("\nChoice: ").strip().upper()
        if choice == 'Q':
            break

        prompt = input("Enter prompt (or press Enter for default): ").strip()
        if not prompt:
            prompt = "The"

        length_str = input("Generation length (default 150): ").strip()
        max_tokens = int(length_str) if length_str.isdigit() else 150

        print()
        selected = [c for c in configs if c[0] == choice]

        if not selected:
            print("Invalid choice.")
            continue

        if selected[0][0] == 'A':
            to_run = [(name, kw) for _, name, kw in configs if kw is not None]
        else:
            _, name, kwargs = selected[0]
            to_run = [(name, kwargs)]

        for name, kwargs in to_run:
            print(f"\n{'─'*60}")
            print(f"Strategy: {name}")
            print(f"{'─'*60}")
            try:
                result = generate(model, dataset, prompt,
                                 max_tokens, device=device, **kwargs)
                print(result)
            except Exception as e:
                print(f"Error: {e}")

        print()


if __name__ == '__main__':
    parser = argparse.ArgumentParser()
    parser.add_argument('--checkpoint', default='checkpoints/best.pt')
    parser.add_argument('--device', default='cpu')
    args = parser.parse_args()

    print(f"Loading model from: {args.checkpoint}")
    model, config = load_model_for_inference(args.checkpoint, args.device)
    print(f"Model loaded: {model.num_params():,} parameters")

    # We need a dataset for the vocabulary; reload from data
    with open('data/input.txt', 'r') as f:
        text = f.read()
    dataset = TextDataset(text, config.block_size)

    sampling_playground(model, dataset, args.device)
```

---

### Mini Project 3: Checkpoint Manager

Saves and loads checkpoints with versioning and metadata tracking.

```python
"""
checkpoint_manager.py

A checkpoint management system with versioning, metadata, and
automatic cleanup of old checkpoints.

Usage: python checkpoint_manager.py --action list --dir checkpoints/
"""

import os
import json
import time
import math
import shutil
import argparse
from pathlib import Path
from datetime import datetime


class CheckpointManager:
    """
    Manages model checkpoints with versioning.
    
    Features:
    - Saves checkpoints with unique version IDs
    - Tracks metadata (step, loss, perplexity, timestamp)
    - Keeps only the N best checkpoints (by val loss)
    - Keeps only the M most recent checkpoints
    - Provides CLI for browsing and loading checkpoints
    """

    def __init__(self, base_dir: str, keep_best: int = 3, keep_recent: int = 5):
        self.base_dir    = Path(base_dir)
        self.keep_best   = keep_best
        self.keep_recent = keep_recent
        self.registry    = {}
        self.base_dir.mkdir(parents=True, exist_ok=True)
        self._load_registry()

    def _registry_path(self):
        return self.base_dir / 'registry.json'

    def _load_registry(self):
        if self._registry_path().exists():
            with open(self._registry_path()) as f:
                self.registry = json.load(f)
        else:
            self.registry = {}

    def _save_registry(self):
        with open(self._registry_path(), 'w') as f:
            json.dump(self.registry, f, indent=2)

    def save(self, model, optimizer, step, train_loss, val_loss,
             config, notes: str = ''):
        """Save a new checkpoint and return its version ID."""
        version_id = f"v_{step:07d}_{int(time.time())}"
        path = self.base_dir / f"{version_id}.pt"

        # Save checkpoint
        ckpt = {
            'model_state_dict':     model.state_dict(),
            'optimizer_state_dict': optimizer.state_dict(),
            'step':       step,
            'train_loss': train_loss,
            'val_loss':   val_loss,
            'config':     config,
        }
        torch.save(ckpt, path)

        # Update registry
        self.registry[version_id] = {
            'path':      str(path),
            'step':      step,
            'train_loss': round(train_loss, 6),
            'val_loss':   round(val_loss, 6),
            'train_ppl':  round(math.exp(min(train_loss, 10)), 4),
            'val_ppl':    round(math.exp(min(val_loss, 10)), 4),
            'timestamp':  datetime.now().isoformat(),
            'notes':      notes,
        }
        self._save_registry()

        # Update best symlink
        self._update_best_symlink(version_id, val_loss)

        # Prune old checkpoints
        self._prune()

        print(f"Saved checkpoint: {version_id} | val_loss={val_loss:.4f}")
        return version_id

    def _update_best_symlink(self, version_id, val_loss):
        """Track and update a 'best.pt' pointer."""
        best_meta_path = self.base_dir / 'best_meta.json'
        best_loss = float('inf')
        if best_meta_path.exists():
            with open(best_meta_path) as f:
                meta = json.load(f)
                best_loss = meta.get('val_loss', float('inf'))

        if val_loss < best_loss:
            src = self.base_dir / f"{version_id}.pt"
            dst = self.base_dir / 'best.pt'
            shutil.copy2(src, dst)
            with open(best_meta_path, 'w') as f:
                json.dump({'version_id': version_id, 'val_loss': val_loss}, f)
            print(f"  >> New best! Copied to best.pt")

    def _prune(self):
        """Remove old checkpoints, keeping keep_best + keep_recent."""
        if len(self.registry) <= max(self.keep_best, self.keep_recent):
            return

        # Sort by val_loss (keep best N)
        by_loss = sorted(self.registry.items(), key=lambda x: x[1]['val_loss'])
        best_ids = {v for v, _ in by_loss[:self.keep_best]}

        # Sort by step (keep most recent M)
        by_step = sorted(self.registry.items(), key=lambda x: x[1]['step'], reverse=True)
        recent_ids = {v for v, _ in by_step[:self.keep_recent]}

        # Keep union of best and recent
        keep_ids = best_ids | recent_ids
        to_delete = [v for v in list(self.registry.keys()) if v not in keep_ids]

        for version_id in to_delete:
            path = Path(self.registry[version_id]['path'])
            if path.exists():
                path.unlink()
            del self.registry[version_id]

        if to_delete:
            print(f"Pruned {len(to_delete)} old checkpoint(s)")
        self._save_registry()

    def load(self, version_id: str = 'best'):
        """Load a checkpoint by version ID, or 'best' / 'latest'."""
        if version_id == 'best':
            path = self.base_dir / 'best.pt'
        elif version_id == 'latest':
            by_step = sorted(self.registry.items(), key=lambda x: x[1]['step'])
            if not by_step:
                raise FileNotFoundError("No checkpoints found")
            version_id = by_step[-1][0]
            path = Path(self.registry[version_id]['path'])
        else:
            if version_id not in self.registry:
                raise KeyError(f"Unknown version: {version_id}")
            path = Path(self.registry[version_id]['path'])

        ckpt = torch.load(path, map_location='cpu')
        print(f"Loaded checkpoint from {path}")
        return ckpt

    def list_checkpoints(self):
        """Print a table of all saved checkpoints."""
        if not self.registry:
            print("No checkpoints found.")
            return

        best_loss = min(v['val_loss'] for v in self.registry.values())

        print(f"\n{'Version ID':<30} {'Step':>8} {'Train Loss':>12} "
              f"{'Val Loss':>10} {'Val PPL':>8} {'Timestamp':<22} {'Notes'}")
        print("-" * 105)

        by_step = sorted(self.registry.items(), key=lambda x: x[1]['step'])
        for vid, meta in by_step:
            marker = " *BEST*" if meta['val_loss'] == best_loss else ""
            print(f"{vid:<30} {meta['step']:>8} {meta['train_loss']:>12.4f} "
                  f"{meta['val_loss']:>10.4f} {meta['val_ppl']:>8.2f} "
                  f"{meta['timestamp'][:19]:<22} {meta.get('notes','')}{marker}")
        print()


def main():
    parser = argparse.ArgumentParser(description='Checkpoint Manager CLI')
    parser.add_argument('--action', choices=['list', 'load', 'delete', 'compare'],
                       default='list')
    parser.add_argument('--dir',     default='checkpoints')
    parser.add_argument('--version', default='best',
                       help='Version ID (or "best"/"latest")')
    args = parser.parse_args()

    mgr = CheckpointManager(args.dir)

    if args.action == 'list':
        mgr.list_checkpoints()

    elif args.action == 'load':
        ckpt = mgr.load(args.version)
        print(f"Checkpoint info:")
        print(f"  Step:       {ckpt['step']}")
        print(f"  Train loss: {ckpt['train_loss']:.4f}")
        print(f"  Val loss:   {ckpt['val_loss']:.4f}")
        print(f"  Config:     {ckpt['config']}")

    elif args.action == 'compare':
        print("Checkpoint Comparison:")
        mgr.list_checkpoints()


if __name__ == '__main__':
    main()
```

---

## Exercises

### Exercise 1: Perplexity Calculation

**Task:** Given a model that assigns probabilities `[0.4, 0.1, 0.3, 0.2]` to four consecutive tokens, compute the perplexity of this 4-token sequence by hand.

```
Hint:
  loss = -1/4 * (log(0.4) + log(0.1) + log(0.3) + log(0.2))
  PPL  = exp(loss)

Expected answer: ~8.9
```

**Extension:** Implement a function `exact_perplexity(probs_list)` that takes a list of per-token probabilities and returns the sequence perplexity.

---

### Exercise 2: Implement Beam Search

Greedy search keeps only 1 candidate at a time. Beam search keeps the top B candidates.

**Task:** Implement a `beam_search(model, dataset, prompt, beam_width=5, max_tokens=50)` function.

```
At each step, each of the B beams expands to V candidates.
Keep only the top B candidates (by cumulative log-probability).

Tip: track beam as list of (cumulative_log_prob, token_list) tuples.
```

---

### Exercise 3: Overfitting Detection

**Task:** Write a function `detect_overfitting(train_losses, val_losses, window=5)` that returns `True` if the model is overfitting.

Your function should:
1. Compute the moving average of val loss over the last `window` checkpoints
2. Compare against the best val loss seen so far
3. Return True if the moving average is more than 5% above the best

---

### Exercise 4: Custom Learning Rate Schedule

**Task:** Implement a learning rate scheduler that:
- Warms up linearly from 0 to `max_lr` over `warmup_steps` steps
- Holds at `max_lr` for a `hold_steps` period  
- Decays cosinely to `min_lr` over the remaining steps

Plot the learning rate schedule for `warmup=500, hold=1000, total=5000`.

---

### Exercise 5: Repetition Penalty

A common problem with greedy and low-temperature sampling is token repetition ("the the the the"). Implement a repetition penalty that reduces the probability of recently seen tokens.

```python
def apply_repetition_penalty(logits, generated_tokens, penalty=1.3):
    """
    For each token that already appears in generated_tokens,
    divide its logit by `penalty` (if logit > 0)
    or multiply by `penalty` (if logit < 0).
    
    This makes repeated tokens less likely.
    """
    # Your implementation here
    pass
```

---

### Exercise 6: Length-Normalized Scoring

When comparing model outputs of different lengths, raw log-probability scores are unfair (longer sequences accumulate more penalty). Implement a `score_sequence(model, dataset, text, length_penalty=0.8)` function that returns the length-normalized log-probability score, suitable for comparing sequences.

---

## Summary

In this chapter we covered:

| Topic | Key Takeaway |
|-------|-------------|
| Training Loop | Teacher forcing enables parallel loss computation across the sequence |
| Loss Curves | Plot both train and val loss; growing gap = overfitting |
| Perplexity | PPL = exp(loss); a PPL of N means the model is as uncertain as choosing among N equally likely words |
| Overfitting | Watch for val loss increasing while train loss decreases |
| Greedy | Fast, deterministic, but repetitive |
| Temperature | Low = conservative, High = creative |
| Top-k | Hard cutoff at k most likely tokens |
| Top-p | Adaptive cutoff based on cumulative probability |
| Checkpointing | Save `state_dict` + optimizer state; track best val loss separately |
| BLEU | n-gram overlap metric; perplexity is better for pure LM evaluation |

---

## Navigation

← [Chapter 28: Building TinyGPT](28-building-tinygpt-from-scratch.md) | [Chapter 30: Scaling Up — From TinyGPT to SLM](30-scaling-to-slm.md) →
