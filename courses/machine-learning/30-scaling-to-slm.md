# Chapter 30: Scaling Up — From TinyGPT to SLM

> **"Scale is not just a technical problem — it's an engineering discipline. Understanding why bigger is harder helps you build smarter."**

---

## Table of Contents
1. [Why Scaling is Hard](#1-why-scaling-is-hard)
2. [The Llama Architecture Innovations](#2-the-llama-architecture-innovations)
3. [Flash Attention](#3-flash-attention)
4. [The KV Cache](#4-the-kv-cache)
5. [LoRA — Parameter-Efficient Fine-tuning](#5-lora--parameter-efficient-fine-tuning)
6. [QLoRA — 4-bit Training](#6-qlora--4-bit-training)
7. [Quantization Explained](#7-quantization-explained)
8. [Memory-Efficient Training Tricks](#8-memory-efficient-training-tricks)
9. [LoRA from Scratch](#9-lora-from-scratch)
10. [Fine-tuning with PEFT + TRL](#10-fine-tuning-with-peft--trl)
11. [Mini Projects](#11-mini-projects)
12. [Summary and Exercises](#12-summary-and-exercises)

---

## Before You Start

**Prerequisites:** Chapter 28 (TinyGPT), Chapter 17 (PyTorch)

```bash
pip install torch transformers peft trl bitsandbytes accelerate datasets
```

---

## 1. Why Scaling is Hard

```
THE THREE WALLS OF SCALING:

1. MEMORY WALL
   A 7B parameter model in FP32:
     7,000,000,000 × 4 bytes = 28 GB just for weights!
   
   But during training you also need:
   • Optimizer state (Adam): 2× the weights = 56 GB
   • Gradients: 1× the weights = 28 GB
   • Activations: depends on batch size and seq len
   
   Total: ~100-120 GB for a 7B model with Adam optimizer
   An A100 GPU has 80 GB. You need multiple GPUs!

2. COMPUTE WALL
   Attention is O(T²) in sequence length.
   Double your context window → 4× the attention compute.
   
   GPT-4 context: 128k tokens
   Naive attention at 128k: 128,000² = 16 billion operations
   PER LAYER, per forward pass.

3. COMMUNICATION WALL
   Multi-GPU training requires synchronizing gradients.
   With 8 GPUs on a DGX A100:
   • All-reduce bandwidth: ~600 GB/s (NVLink)
   • 28 GB of gradients per step: ~0.05 seconds/step just for sync
   • Becomes the bottleneck before compute!
```

**The solutions:** Better architectures (Llama), better attention (FlashAttention), better fine-tuning (LoRA), better quantization (INT4/INT8), better memory management (gradient checkpointing).

---

## 2. The Llama Architecture Innovations

Llama (Meta, 2023) introduced several improvements over vanilla GPT:

```
VANILLA GPT vs LLAMA ARCHITECTURE:

Feature              GPT-2/3              LLaMA/2/3
─────────────────────────────────────────────────────
Positional Encoding  Learned absolute     RoPE (Rotary)
Activation Function  GELU                 SwiGLU
Normalization        Post-LN              Pre-RMSNorm
Key-Value Heads      Same as Q heads      Grouped Query Attention (GQA)
Attention Bias       Yes                  No
─────────────────────────────────────────────────────
```

### RoPE — Rotary Position Encoding

Instead of adding a position vector to embeddings, RoPE rotates Q and K vectors based on position. The dot product between rotated Q and K naturally encodes relative position.

**Why it's better:** RoPE can generalize to longer sequences than seen during training (with some tricks), while learned positional embeddings are hard-capped at `max_seq_len`.

### SwiGLU Activation

Instead of: `FFN(x) = GELU(xW₁)W₂`

LLaMA uses: `FFN(x) = (swish(xW₁) ⊙ xW_gate)W₂`

The gating mechanism selectively passes information, giving the FFN more expressivity.

### RMSNorm

Instead of LayerNorm (compute mean and variance), RMSNorm only computes the root mean square:

```
LayerNorm: y = (x - mean) / sqrt(var + ε) * γ + β  (expensive)
RMSNorm:   y = x / sqrt(mean(x²) + ε) * γ          (cheaper, same quality)
```

### Grouped Query Attention (GQA)

```
MULTI-HEAD ATTENTION (standard):
Query heads:  8
Key heads:    8  (same as Q)
Value heads:  8  (same as Q)
KV cache size: proportional to 8 heads

GROUPED QUERY ATTENTION:
Query heads:  8
Key heads:    2  (4 Q heads share each K head)
Value heads:  2  (4 Q heads share each V head)
KV cache size: 4× smaller! Critical for long-context inference.
```

---

## 3. Flash Attention

The standard attention implementation is memory-bandwidth-bound:

```
STANDARD ATTENTION MEMORY ACCESS PATTERN:

GPU HBM (slow, large):  S, P, O matrices — T² memory!
GPU SRAM (fast, small): tiny scratchpad

Standard attention loads and stores T² values between HBM and SRAM.
For T=8192: 67M reads/writes — most GPU time is waiting for memory!

FLASHATTENTION solution:
• Process attention in TILES that fit in SRAM
• Never materialize the full T×T attention matrix
• Compute exact same result, just without writing T² to HBM

Result: 2-4× faster, 10-20× less memory for attention
Enables training with much longer sequences!
```

```python
# In PyTorch 2.0+, Flash Attention is built in:
import torch.nn.functional as F

# Instead of manual QKV computation, use:
output = F.scaled_dot_product_attention(
    query, key, value,
    attn_mask=None,
    dropout_p=0.0,
    is_causal=True   # applies causal mask automatically
)
# PyTorch automatically uses Flash Attention when available!
```

---

## 4. The KV Cache

Understanding the KV cache explains why chatbot responses start slow then speed up.

```
WITHOUT KV CACHE (naive inference):
Token 1: process "The" → compute K, V for position 0
Token 2: process "cat" → recompute K, V for positions 0, 1 (wasted!)
Token 3: process "sat" → recompute K, V for positions 0, 1, 2 (wasted!)
...
Token N: recompute everything from scratch. O(N²) total compute!

WITH KV CACHE:
Token 1: compute K₀, V₀ → STORE in cache
Token 2: compute K₁, V₁ → STORE in cache, REUSE K₀, V₀
Token 3: compute K₂, V₂ → STORE in cache, REUSE K₀, K₁, V₀, V₁
...
Token N: compute only Kₙ, Vₙ. REUSE all previous K, V. O(N) per step!

Cost: memory grows linearly with sequence length
      (this is why long contexts hit OOM errors)
```

```python
# KV cache in PyTorch (simplified concept):

class CachedAttentionHead(nn.Module):
    def __init__(self, d_model, head_size):
        super().__init__()
        self.query = nn.Linear(d_model, head_size, bias=False)
        self.key   = nn.Linear(d_model, head_size, bias=False)
        self.value = nn.Linear(d_model, head_size, bias=False)
        self.head_size = head_size
        
        # Cache storage
        self.k_cache = None
        self.v_cache = None
    
    def forward(self, x, use_cache=False):
        """x: (B, T, d_model) where T=1 during cached inference."""
        Q = self.query(x)
        K = self.key(x)
        V = self.value(x)
        
        if use_cache:
            if self.k_cache is None:
                self.k_cache = K
                self.v_cache = V
            else:
                # Append new K, V to cache
                self.k_cache = torch.cat([self.k_cache, K], dim=1)
                self.v_cache = torch.cat([self.v_cache, V], dim=1)
            
            K = self.k_cache
            V = self.v_cache
        
        scale = self.head_size ** -0.5
        scores = Q @ K.transpose(-2, -1) * scale
        weights = torch.softmax(scores, dim=-1)
        return weights @ V
    
    def reset_cache(self):
        self.k_cache = None
        self.v_cache = None
```

---

## 5. LoRA — Parameter-Efficient Fine-tuning

Full fine-tuning a 7B model requires 100+ GB of GPU memory. LoRA fine-tunes with as little as 8 GB.

```
THE LORA INSIGHT:

When you fine-tune a pre-trained model, the weight updates
tend to have LOW INTRINSIC RANK — they don't need to be
full-rank matrices.

Instead of updating W (d × d), approximate the update as:
  ΔW = B × A   where B is (d × r) and A is (r × d)
  
  r << d  (rank r is much smaller than d)
  
PARAMETER COUNT:
  Full fine-tune of W (4096 × 4096) = 16,777,216 parameters
  LoRA with r=16: B (4096×16) + A (16×4096) = 131,072 parameters
  
  Savings: 99.2% fewer trainable parameters!

During training: W₀ is frozen, only A and B are updated.
During inference: merge back: W = W₀ + BA

WHICH LAYERS TO APPLY LORA TO?
  Usually: query and value projections in attention
  Sometimes: all linear layers including FFN
  Rule of thumb: start with attention Q and V

HYPERPARAMETERS:
  r (rank): 4, 8, 16, 32, 64 — higher = more parameters = more capacity
  alpha: scaling factor, typically alpha = r or alpha = 2r
  dropout: 0.05-0.1 for regularization
```

---

## 6. QLoRA — 4-bit Training

QLoRA (2023) combined 4-bit quantization with LoRA, enabling fine-tuning of 65B models on a single 48GB GPU.

```
QLORA RECIPE:
1. Load model in 4-bit (NF4 quantization) → 4× memory savings
2. Add LoRA adapters (trainable, in FP16)
3. Train ONLY LoRA adapters
4. Quantized base model weights never change

MEMORY SAVINGS:
  Llama 7B in FP32:   28 GB
  Llama 7B in BF16:   14 GB
  Llama 7B in INT8:    7 GB
  Llama 7B in NF4:     3.5 GB  ← QLoRA base
  + LoRA adapters:    +0.5 GB
  + Optimizer state:  +1 GB (only for LoRA params!)
  Total: ~5 GB! Fits on consumer GPUs.
```

---

## 7. Quantization Explained

```
FLOATING POINT vs QUANTIZED:

FP32: 32 bits per number
  1 bit sign | 8 bits exponent | 23 bits mantissa
  Range: ±3.4 × 10^38
  Precision: ~7 decimal digits

INT8: 8 bits per number
  Range: -128 to 127
  
ABSMAX QUANTIZATION (simplest):
  Given weight tensor W with max absolute value M:
    W_quant = round(W × 127 / M)       (compress to INT8)
    W_dequant = W_quant × M / 127      (restore, with rounding error)
  
  Example:
    W = [0.5, -0.3, 1.2, -0.8]
    M = 1.2
    W_quant = round([53, -32, 127, -85])  ← INT8
    W_dequant ≈ [0.501, -0.302, 1.2, -0.803]  ← tiny error
  
NF4 (Normal Float 4): More sophisticated
  Uses a non-uniform quantization grid optimized for
  the normal distribution of neural network weights.
  4-bit with near FP16 quality!
```

---

## 8. Memory-Efficient Training Tricks

```python
# TRICK 1: Gradient Checkpointing
# Don't store ALL activations — recompute on backward pass
# Saves ~60% memory, costs ~30% more compute

model = GPT(config)
model.gradient_checkpointing_enable()  # HuggingFace API

# Manual PyTorch:
from torch.utils.checkpoint import checkpoint

class TransformerBlockWithCheckpoint(nn.Module):
    def forward(self, x, mask=None):
        return checkpoint(self._forward, x, mask)
    
    def _forward(self, x, mask):
        x = x + self.attn(self.ln1(x), mask)
        x = x + self.ffn(self.ln2(x))
        return x


# TRICK 2: Mixed Precision Training
from torch.cuda.amp import autocast, GradScaler

scaler = GradScaler()

for x, y in dataloader:
    with autocast():  # Forward pass in FP16
        logits, loss = model(x, y)
    
    # Backward pass with loss scaling (prevents underflow in FP16)
    scaler.scale(loss).backward()
    scaler.step(optimizer)
    scaler.update()


# TRICK 3: Gradient Accumulation
# Simulate larger batch size without more memory
accumulation_steps = 8  # Effective batch = batch_size × 8

for i, (x, y) in enumerate(dataloader):
    logits, loss = model(x, y)
    loss = loss / accumulation_steps
    loss.backward()
    
    if (i + 1) % accumulation_steps == 0:
        optimizer.step()
        optimizer.zero_grad()


# TRICK 4: 8-bit Adam (bitsandbytes)
import bitsandbytes as bnb

optimizer = bnb.optim.Adam8bit(
    model.parameters(),
    lr=2e-4
)
# 8-bit Adam: 4× less optimizer memory!
```

---

## 9. LoRA from Scratch

```python
import torch
import torch.nn as nn
import math

class LoRALayer(nn.Module):
    """
    Drop-in replacement for nn.Linear that adds LoRA.
    
    The original weight W₀ is frozen.
    Only A and B are trained.
    """
    
    def __init__(self, in_features: int, out_features: int, 
                 r: int = 16, alpha: float = 32):
        super().__init__()
        self.r = r
        self.alpha = alpha
        self.scaling = alpha / r
        
        # Frozen base weight (loaded from pretrained model)
        self.weight = nn.Parameter(
            torch.empty(out_features, in_features),
            requires_grad=False  # FROZEN
        )
        
        # LoRA matrices (trainable)
        self.lora_A = nn.Parameter(torch.empty(r, in_features))
        self.lora_B = nn.Parameter(torch.zeros(out_features, r))
        
        # Initialize: A with Kaiming, B with zeros
        # (B=0 ensures LoRA starts as identity — no disruption at start)
        nn.init.kaiming_uniform_(self.lora_A, a=math.sqrt(5))
        nn.init.zeros_(self.lora_B)
    
    def forward(self, x):
        # Base forward pass (frozen weights)
        base_out = x @ self.weight.T
        
        # LoRA forward pass (trained weights)
        lora_out = (x @ self.lora_A.T) @ self.lora_B.T
        lora_out = lora_out * self.scaling
        
        return base_out + lora_out
    
    def merge_weights(self):
        """Merge LoRA back into base weight for efficient inference."""
        self.weight.data += (self.lora_B @ self.lora_A) * self.scaling
        # After merging, LoRA is zero-valued
        self.lora_A.data.zero_()
        self.lora_B.data.zero_()


def add_lora_to_model(model, r=16, alpha=32, target_modules=('query', 'value')):
    """Replace target linear layers with LoRA versions."""
    for name, module in model.named_modules():
        if any(t in name for t in target_modules):
            if isinstance(module, nn.Linear):
                # Create LoRA replacement
                lora = LoRALayer(module.in_features, module.out_features, r, alpha)
                lora.weight.data = module.weight.data.clone()
                
                # Replace in parent module
                parent = model
                parts = name.split('.')
                for part in parts[:-1]:
                    parent = getattr(parent, part)
                setattr(parent, parts[-1], lora)
    
    # Count parameters
    trainable = sum(p.numel() for p in model.parameters() if p.requires_grad)
    total = sum(p.numel() for p in model.parameters())
    print(f"Trainable parameters: {trainable:,} / {total:,} ({100*trainable/total:.2f}%)")
    
    return model
```

---

## 10. Fine-tuning with PEFT + TRL

```python
# Complete QLoRA fine-tuning example using HuggingFace
# Fine-tune Llama-3.2-1B on a custom instruction dataset

from transformers import (
    AutoModelForCausalLM, 
    AutoTokenizer, 
    TrainingArguments,
    BitsAndBytesConfig
)
from peft import LoraConfig, get_peft_model, TaskType
from trl import SFTTrainer
from datasets import Dataset
import torch

# 1. Load model in 4-bit
bnb_config = BitsAndBytesConfig(
    load_in_4bit=True,
    bnb_4bit_quant_type="nf4",
    bnb_4bit_compute_dtype=torch.bfloat16,
    bnb_4bit_use_double_quant=True,  # QLoRA's double quantization
)

model_name = "meta-llama/Llama-3.2-1B"  # or "microsoft/phi-2"

model = AutoModelForCausalLM.from_pretrained(
    model_name,
    quantization_config=bnb_config,
    device_map="auto",
)
tokenizer = AutoTokenizer.from_pretrained(model_name)
tokenizer.pad_token = tokenizer.eos_token

# 2. Configure LoRA
lora_config = LoraConfig(
    r=16,
    lora_alpha=32,
    target_modules=["q_proj", "v_proj"],  # Apply to attention
    lora_dropout=0.05,
    bias="none",
    task_type=TaskType.CAUSAL_LM,
)

model = get_peft_model(model, lora_config)
model.print_trainable_parameters()
# Output: trainable params: 3,407,872 || all params: 1,239,842,816 (0.27%)

# 3. Prepare dataset
training_data = [
    {"text": "### Instruction: What is machine learning?\n### Response: Machine learning is..."},
    {"text": "### Instruction: Explain gradient descent.\n### Response: Gradient descent is..."},
    # Add more examples...
]
dataset = Dataset.from_list(training_data)

# 4. Training arguments
training_args = TrainingArguments(
    output_dir="./lora_output",
    num_train_epochs=3,
    per_device_train_batch_size=4,
    gradient_accumulation_steps=4,
    learning_rate=2e-4,
    fp16=True,
    logging_steps=50,
    save_steps=200,
    warmup_steps=100,
    lr_scheduler_type="cosine",
    report_to="none",
)

# 5. Train with SFTTrainer
trainer = SFTTrainer(
    model=model,
    train_dataset=dataset,
    args=training_args,
    dataset_text_field="text",
    max_seq_length=512,
)

trainer.train()

# 6. Save LoRA adapter
model.save_pretrained("./lora_adapter")

# 7. Load and merge for inference
from peft import PeftModel
base_model = AutoModelForCausalLM.from_pretrained(model_name)
merged_model = PeftModel.from_pretrained(base_model, "./lora_adapter")
merged_model = merged_model.merge_and_unload()  # Bake LoRA into weights
merged_model.save_pretrained("./merged_model")
```

---

## 11. Mini Projects

### Mini Project 1: LoRA Explorer

**What You'll Build:** Apply LoRA to TinyGPT, compare parameter counts and training dynamics.

```python
# Apply LoRA to your TinyGPT from chapter 28
model = TinyGPT(config)

# Before LoRA
total_before = sum(p.numel() for p in model.parameters())
print(f"Before LoRA: {total_before:,} total parameters")

# Add LoRA
model = add_lora_to_model(model, r=8, target_modules=('query', 'key', 'value'))

trainable = sum(p.numel() for p in model.parameters() if p.requires_grad)
print(f"After LoRA:  {trainable:,} trainable parameters ({100*trainable/total_before:.1f}%)")

# Train with LoRA and compare loss curve to full fine-tuning
```

**Bonus Challenge:** Experiment with different ranks (r=2, 4, 8, 16, 32). Find the minimum rank that achieves 95% of full fine-tuning performance.

---

### Mini Project 2: Memory Profiler

**What You'll Build:** A script that measures GPU memory at each stage of training.

```python
import torch

def profile_memory(model, batch, label=""):
    """Measure GPU memory before and after an operation."""
    if torch.cuda.is_available():
        torch.cuda.synchronize()
        mem_before = torch.cuda.memory_allocated() / 1024**3  # GB
        
        yield
        
        torch.cuda.synchronize()
        mem_after = torch.cuda.memory_allocated() / 1024**3
        print(f"  {label}: {mem_before:.2f} GB → {mem_after:.2f} GB (Δ{mem_after-mem_before:+.2f} GB)")
    else:
        print(f"  {label}: (CPU mode, no GPU memory tracking)")
        yield

# Usage:
with profile_memory(model, batch, "Forward pass"):
    logits, loss = model(batch_x, batch_y)

with profile_memory(model, batch, "Backward pass"):
    loss.backward()
```

---

### Mini Project 3: Quantization Comparator

**What You'll Build:** Compare model sizes and inference speed at different precisions.

```python
import torch
import time

def compare_precisions(model_path: str, test_input: torch.Tensor):
    """Compare FP32, FP16, and INT8 inference."""
    
    results = {}
    
    for precision in ['fp32', 'fp16']:
        model = TinyGPT(config)
        model.load_state_dict(torch.load(model_path))
        
        if precision == 'fp16':
            model = model.half()
        
        # Measure size
        param_size = sum(p.numel() * p.element_size() for p in model.parameters())
        results[precision] = {'size_mb': param_size / 1024**2}
        
        # Measure speed
        model.eval()
        times = []
        with torch.no_grad():
            for _ in range(100):
                start = time.perf_counter()
                _ = model(test_input)
                times.append(time.perf_counter() - start)
        
        results[precision]['ms_per_forward'] = sum(times) / len(times) * 1000
    
    print(f"\n{'Precision':<10} {'Size (MB)':>10} {'Speed (ms)':>12}")
    print("-"*35)
    for name, r in results.items():
        print(f"{name:<10} {r['size_mb']:>10.1f} {r['ms_per_forward']:>12.3f}")

# compare_precisions('tinygpt_best.pt', torch.randint(0, 65, (1, 128)))
```

---

## 12. Summary and Exercises

```
KEY TAKEAWAYS:
══════════════════════════════════════════════════════════
ARCHITECTURE (Llama innovations):
  • RoPE: relative position, better length generalization
  • SwiGLU: gated activation, more expressive FFN
  • RMSNorm: cheaper normalization, same quality
  • GQA: fewer KV heads, 4× KV cache reduction

EFFICIENCY:
  • Flash Attention: 2-4× faster, 10-20× less attn memory
  • KV Cache: O(N) vs O(N²) inference
  • Gradient checkpointing: 60% memory, 30% slower
  • Mixed precision (BF16): 2× memory, 2× faster

FINE-TUNING:
  • LoRA: W = W₀ + BA, 99% fewer trainable params
  • QLoRA: 4-bit base + FP16 LoRA, fits 7B on consumer GPU
  • Rule of thumb: start with LoRA on Q,V projections, r=16
══════════════════════════════════════════════════════════
```

**Exercise 1:** Implement RMSNorm from scratch and verify it produces the same output as PyTorch's `nn.LayerNorm` when the input has zero mean. Measure speed: is it actually faster?

**Exercise 2:** Visualize the KV cache growth. For a TinyGPT model generating 100 tokens, plot memory usage at each step. How many tokens can you generate before running out of memory on your hardware?

**Exercise 3:** Implement a "LoRA rank sweeper": train TinyGPT with LoRA at ranks r=1, 2, 4, 8, 16. For each rank, record: parameter count, validation loss after 1000 steps. Plot rank vs loss. Find the "knee" of the curve.

**Exercise 4:** The `alpha` parameter in LoRA scales the update: `ΔW = (alpha/r) × BA`. Experiment with alpha=r, alpha=2r, alpha=4r. How does the effective learning rate change? What's the best setting for your TinyGPT?

**Exercise 5:** Flash Attention in PyTorch: benchmark `F.scaled_dot_product_attention(is_causal=True)` vs your manual attention implementation from chapter 28. Measure time and memory for sequence lengths 128, 256, 512, 1024. What's the crossover point where Flash Attention becomes clearly faster?

---

← [Chapter 29: Training and Evaluating Your Language Model](./29-training-and-evaluating-lm.md) | [Chapter 31: Build a Story Generator](./31-project-build-story-generator.md) →
