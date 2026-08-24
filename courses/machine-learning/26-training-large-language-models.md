# Chapter 26: Training Large Language Models — From Pre-training to RLHF to Constitutional AI

> **"Training a large language model is one of the most complex engineering feats in software history — a billion-dollar experiment where you don't find out how it went until weeks later."**

---

## Table of Contents
1. [The Pre-training Data Pipeline](#1-the-pre-training-data-pipeline)
2. [Distributed Training](#2-distributed-training)
3. [Pre-training Compute Costs](#3-pre-training-compute-costs)
4. [Supervised Fine-tuning (SFT)](#4-supervised-fine-tuning-sft)
5. [RLHF: Reinforcement Learning from Human Feedback](#5-rlhf-reinforcement-learning-from-human-feedback)
6. [DPO: Direct Preference Optimization](#6-dpo-direct-preference-optimization)
7. [Constitutional AI (Anthropic)](#7-constitutional-ai-anthropic)
8. [Alignment — The Bigger Picture](#8-alignment--the-bigger-picture)
9. [Inference Optimization: KV Cache and Speculative Decoding](#9-inference-optimization-kv-cache-and-speculative-decoding)
10. [Production Serving: vLLM](#10-production-serving-vllm)
11. [Summary](#11-summary)
12. [Exercises](#12-exercises)

---

## 1. The Pre-training Data Pipeline

Training an LLM requires not just the model and the compute — it requires building a massive, carefully curated data pipeline.

### The Raw Data Sources

```
MAJOR DATA SOURCES FOR LLM PRE-TRAINING:
═══════════════════════════════════════════════════════════════════

COMMON CRAWL (CC):
  What: Automated web crawl of the entire public internet
  Size: Petabytes (monthly snapshot: ~200-250TB compressed)
  Quality: VERY MIXED — spam, duplicates, low-quality text everywhere
  Filtering required: aggressive deduplication, quality classifiers, toxicity filtering
  Share in major models: 60-85% of data (by volume)
  
WIKIPEDIA:
  What: English Wikipedia + multi-lingual versions
  Size: ~20GB English, ~100GB all languages
  Quality: EXCELLENT — encyclopedic, fact-checked, well-organized
  Share in major models: 3-5%
  Note: High quality → often weighted higher than volume suggests
  
BOOKS:
  What: Project Gutenberg (public domain), BookCorpus, Open Library
  Size: Several hundred GB
  Quality: EXCELLENT — long-form coherent text, good for reasoning
  Share in major models: 5-15%
  
GITHUB / CODE:
  What: Public code repositories
  Size: Several hundred GB (filtered)
  Languages: Python, JavaScript, Java, C++, SQL, etc.
  Effect: Dramatically improves code generation AND reasoning!
  (Code training seems to improve logical reasoning in general)
  Share in GPT-3.5/4: ~8%
  
PAPERS / SCIENTIFIC:
  What: arXiv, PubMed, Semantic Scholar
  Size: 100-200GB
  Quality: EXCELLENT for domain knowledge
  
StackExchange / Reddit:
  What: Q&A forums with quality signals (votes, accepted answers)
  Size: Several hundred GB
  Quality: Good (vote filtering), some noise
```

### Data Quality Filtering

The raw Common Crawl data is mostly unusable. Filtering is essential:

```
CC DATA PIPELINE (typical):
Step 1: LANGUAGE IDENTIFICATION
  Keep only English (or target languages)
  fastText language classifier
  Discard: ~70% of raw data (non-English)
  
Step 2: QUALITY FILTERING
  Heuristic filters:
    - Minimum length: discard pages < 50 words
    - Maximum repetition: discard if >20% duplicate lines
    - Stop word ratio: English text has ~30% function words (the, and, of...)
      Less → probably not natural language
    - Perplexity filter: discard text with too-high perplexity under trained LM
      Catches garbled text, random characters
    - Fraction of alphabetic characters > 0.8
  
  ML-based classifier (C4/RefinedWeb approach):
    Train a binary classifier: "high quality" vs "low quality" text
    Train on Wikipedia (positive) vs random CC (negative)
    Apply to all CC data: keep only "Wikipedia-like quality"
    Typical threshold: top 20% by classifier score
    
Step 3: DEDUPLICATION
  Exact deduplication: remove exact duplicate pages
  Fuzzy deduplication (MinHash):
    - Represent each document as a set of n-gram hashes
    - MinHash + LSH to find near-duplicates efficiently
    - Remove duplicates even when slightly modified
  
  WHY DEDUPLICATION MATTERS:
    LLMs memorize repeated text verbatim
    Deduplication → less memorization → better generalization
    Chinchilla paper: dedup improved performance significantly
    
Step 4: TOXICITY FILTERING
  Remove content that would be harmful in training:
  - Hate speech, explicit content, extremist material
  Tools: Perspective API, custom classifiers
  
  TRADEOFF: Too aggressive filtering → model refuses to discuss
            legitimate topics (violence in history, medical descriptions)
  
Step 5: DATA MIX
  After filtering, combine sources with desired ratios:
  
  LLAMA-3 DATA MIX (approximately):
  ┌────────────────────────────────────────┐
  │ Source          │ Share │ Quality Weight │
  ├────────────────────────────────────────┤
  │ Web (filtered CC)│ 65%  │ 1×             │
  │ Code (GitHub)   │ 8%   │ 1×             │
  │ Wikipedia       │ 5%   │ 4×             │
  │ Books           │ 5%   │ 4×             │
  │ ArXiv           │ 2%   │ 2×             │
  │ StackExchange   │ 2%   │ 2×             │
  │ Other           │ 13%  │ 1×             │
  └────────────────────────────────────────┘
  
  "Quality weight" = how many times to repeat in training
  High-quality sources repeated to counteract their smaller volume
```

### Open Datasets

```
MAJOR OPEN LLM DATASETS:

The Pile (EleutherAI, 2021):
  825GB, 22 diverse sources
  Wikipedia, Books, GitHub, PubMed, FreeLaw, DM Mathematics, etc.
  First major open dataset for LLM training
  
RefinedWeb (Falcon, 2023):
  600B tokens, Common Crawl only but heavily filtered
  MacroFilter: removed ~90% of CC as low quality
  Showed: quality > quantity
  
RedPajama V1 (2023):
  1.2 trillion tokens
  Open reproduction of LLaMA training data
  
RedPajama V2 (2023):
  30 trillion tokens (raw), quality annotations provided
  Researchers can filter to their preferences
  
Dolma (Allen AI, 2024):
  3 trillion tokens, released with data cards and documentation
  Provenance tracked: know exactly where each document came from
  
StarCoder / The Stack (BigCode, 2023):
  6.4TB of code from GitHub, 300+ programming languages
  Permissive licenses only (no GPL, etc.)
```

### Tokenizer Training

```
TOKENIZER TRAINING:

After data is ready, train the tokenizer FIRST:
  1. Sample representative subset (~100GB) from final training data
  2. Run BPE/SentencePiece training on this sample
  3. Choose vocabulary size:
     - BERT: 30,522
     - GPT-2: 50,257
     - LLaMA 1/2: 32,000
     - LLaMA 3: 128,000
     - GPT-4 (cl100k_base): 100,277
  
  Larger vocab → shorter sequences → less compute per token
  Larger vocab → bigger embedding table → more parameters
  
BYTE-LEVEL BPE (GPT-2, LLaMA 3):
  Base vocabulary = 256 bytes (all possible byte values)
  BPE merges on top of bytes
  
  Advantage: ZERO unknown tokens
  Any Unicode character → encoded as bytes
  No [UNK] token ever needed
  
  Compression ratio: how many characters per token (higher = more efficient)
  GPT-2 (50k vocab): ~4 chars/token for English
  LLaMA-3 (128k vocab): ~5 chars/token for English
  (English-optimized; Chinese/Korean may be 1:1 or worse)
```

---

## 2. Distributed Training

Training a 175B parameter model on a single GPU is impossible. It requires splitting the work across hundreds or thousands of GPUs.

### Memory Requirements — Why We Can't Fit in One GPU

```
MEMORY FOR GPT-3 (175B parameters):
═══════════════════════════════════════════════════════════════════

MODEL WEIGHTS:
  175B parameters × 4 bytes (FP32) = 700 GB
  175B parameters × 2 bytes (FP16) = 350 GB  ← minimum for inference

  A single A100 GPU: 80 GB HBM
  175B FP16: 350 GB → needs at least 5× A100s just for weights!

TRAINING STATES (much worse):
  Mixed precision training (FP16 forward + FP32 optimizer):
    FP16 weights:       2 bytes × 175B = 350 GB
    FP32 master weights:4 bytes × 175B = 700 GB  ← copy for optimizer
    
  Adam Optimizer states (FP32):
    First moment (m):   4 bytes × 175B = 700 GB
    Second moment (v):  4 bytes × 175B = 700 GB
    
  Gradients (FP16):     2 bytes × 175B = 350 GB
  
  TOTAL: 350 + 700 + 700 + 700 + 350 = 2,800 GB = 2.8 TB!
  
  At 80GB per A100: need 35 A100s minimum, just for memory.
  GPT-3 was actually trained on 1,024 A100s.
═══════════════════════════════════════════════════════════════════
```

### Data Parallelism

```
DATA PARALLELISM (DP):
  Simplest form of distributed training.
  
  HOW IT WORKS:
    - Replicate the ENTIRE model on each of N GPUs
    - Split each batch into N micro-batches
    - Each GPU processes its micro-batch independently
    - After backward pass, average (all-reduce) gradients across all GPUs
    - Update all model copies with averaged gradients
    
  EXAMPLE (4 GPUs, batch_size=256):
    GPU 0: processes samples 0-63   → gradients_0
    GPU 1: processes samples 64-127 → gradients_1
    GPU 2: processes samples 128-191→ gradients_2
    GPU 3: processes samples 192-255→ gradients_3
    
    All-reduce: average gradients_0..3 → gradient_avg
    
    All GPUs update: weights -= lr * gradient_avg
    
  LIMITATION: Each GPU must hold the entire model.
    With 175B params → 350GB per GPU → not possible even with 80GB GPUs.
    Data parallelism alone: limited to models < 80GB (FP16) = < 40B params.
```

### Tensor Parallelism

```
TENSOR PARALLELISM (TP, Megatron-LM):
  Split individual weight MATRICES across GPUs.
  
  EXAMPLE — Splitting a Linear layer (d_model × d_ff) across 4 GPUs:
  
    Standard: y = x @ W where W ∈ R^(512 × 2048)
    
    Split W column-wise into 4 pieces:
      GPU 0: W_0 ∈ R^(512 × 512)
      GPU 1: W_1 ∈ R^(512 × 512)
      GPU 2: W_2 ∈ R^(512 × 512)
      GPU 3: W_3 ∈ R^(512 × 512)
    
    Each GPU computes partial result:
      GPU i: y_i = x @ W_i  (x is broadcast to all GPUs)
    
    Concatenate: y = concat(y_0, y_1, y_2, y_3) = x @ W
    
  For attention: split across heads
    4-head attention on 4 GPUs: each GPU handles 1 head
  
  REQUIRES: frequent all-gather operations between GPUs
  BEST FOR: intra-node (within same server, high bandwidth NVLink)
  TYPICAL TENSOR PARALLEL DEGREE: 4-8 GPUs per node (NVLink)
```

### Pipeline Parallelism

```
PIPELINE PARALLELISM (PP):
  Split model LAYERS across GPUs.
  
  EXAMPLE — 96-layer model on 4 GPUs:
    GPU 0: Layers  0-23
    GPU 1: Layers 24-47
    GPU 2: Layers 48-71
    GPU 3: Layers 72-95
  
  NAIVE PIPELINE (slow, lots of idle time):
    t=0: GPU 0 processes batch
    t=1: GPU 1 processes (GPU 0 idle)
    t=2: GPU 2 processes (GPU 0,1 idle)
    t=3: GPU 3 processes (GPU 0,1,2 idle)
    
    "Pipeline bubble": most GPUs idle most of the time!
    
  MICRO-BATCH PIPELINING (GPipe/PipeDream):
    Split batch into smaller micro-batches.
    Overlap computation: while GPU 1 processes micro-batch 1,
    GPU 0 starts processing micro-batch 2.
    
    t=0: GPU 0 processes microbatch 0
    t=1: GPU 1 processes mb0, GPU 0 processes mb1
    t=2: GPU 2 processes mb0, GPU 1 processes mb1, GPU 0 processes mb2
    t=3: GPU 3 processes mb0, GPU 2 processes mb1, ... GPU 0 mb3
    
    Pipeline efficiency = (m - p + 1) / m ≈ (m / (m + p - 1))
    Where m = number of micro-batches, p = pipeline stages
    With m=16 micro-batches, p=4 stages: efficiency = 81%
  
  BEST FOR: inter-node communication (between servers, lower bandwidth)
```

### 3D Parallelism

```
3D PARALLELISM (used for GPT-3, Megatron-Turing NLG):
  Combine all three: DP + TP + PP
  
  EXAMPLE: 512 GPUs (64 nodes × 8 GPUs each)
    Tensor Parallel: 8 GPUs per node (NVLink = high bandwidth)
    Pipeline Parallel: 8 pipeline stages (1 node per stage)
    Data Parallel: 8 data parallel replicas
    
    8 TP × 8 PP × 8 DP = 512 GPUs ✓
    
  MEMORY DISTRIBUTION:
    No single GPU holds the full model
    Each GPU holds: total_params / (TP × PP) fraction
    
    175B params / (8 × 8) = 2.7B params per GPU
    2.7B × 2 bytes (FP16) = 5.4 GB just for model weights
    ← Now it fits on an 80GB GPU with room for activations!

NVIDIA's Megatron-LM: main framework for tensor parallelism
Microsoft's DeepSpeed: main framework for ZeRO optimization
NVIDIA's Apex: mixed precision training utilities
```

### ZeRO Optimizer (DeepSpeed)

```
ZeRO (Zero Redundancy Optimizer):
  Key insight: In data parallelism, each GPU has IDENTICAL copies of:
    - Optimizer states (Adam m and v)
    - Gradients
    - Model parameters
  
  This is massive redundancy! ZeRO partitions these across GPUs.
  
  ZERO STAGE 1: Partition optimizer states
    Each GPU holds: 1/N of optimizer states (m, v)
    After gradient computation: each GPU only updates its 1/N parameters
    Memory: optimizer states reduced by N×
    
  ZERO STAGE 2: Partition gradients + optimizer states
    After backward: each GPU only stores gradients for its 1/N parameters
    Memory: gradients + optimizer states reduced by N×
    
  ZERO STAGE 3: Partition everything (parameters + gradients + optimizer)
    Each GPU holds: only 1/N of parameters!
    During forward: gather parameters on-the-fly, release after use
    Maximum memory reduction: N× (where N = number of GPUs)
    
    Tradeoff: more communication (parameter gathering) ↔ less memory
    
  MEMORY SAVINGS (GPT-3 at 128 GPUs):
    Without ZeRO: 2,800 GB / 128 = 21.9 GB per GPU (still 175B params each)
    ZeRO-3:       2,800 GB / 128 = 21.9 GB TOTAL per GPU (shared!)
    ← Model fits on much fewer GPUs!

ZeRO-Offload: offload optimizer to CPU RAM (much more available)
ZeRO-Infinity: offload to NVMe SSDs
```

### Inter-GPU Communication

```
COMMUNICATION PRIMITIVES:

All-Reduce (ring-based):
  Used in data parallelism to average gradients.
  All GPUs start with gradient_i.
  After all-reduce: all GPUs have sum(gradient_i) / N.
  Ring algorithm: each GPU sends and receives from neighbors.
  Cost: 2(N-1)/N × message_size ≈ 2 × message_size
  
All-Gather:
  Used in ZeRO-3, tensor parallelism.
  GPU i has element_i.
  After all-gather: all GPUs have [element_0, element_1, ..., element_{N-1}]
  
Reduce-Scatter:
  Each GPU has full gradient vector.
  After reduce-scatter: GPU i has sum of gradient[i*chunk:(i+1)*chunk]
  Used in ZeRO stage 2.
  
INTERCONNECT BANDWIDTH:
  Within node (NVLink, GPU-to-GPU):
    A100 NVLink: 600 GB/s bidirectional
    H100 NVLink: 900 GB/s
  
  Between nodes (InfiniBand, server-to-server):
    IB HDR: 200 Gb/s ≈ 25 GB/s per port
    IB NDR: 400 Gb/s
  
  Why this matters:
    Tensor parallelism requires very fast communication → needs NVLink
    Pipeline/data parallelism works with slower inter-node IB
```

---

## 3. Pre-training Compute Costs

```
COMPUTING FLOPs FOR TRANSFORMER TRAINING:

FLOPs per token (forward pass only):
  FLOPs ≈ 6 × N    (where N = number of parameters)
  
  This approximation comes from:
    - Each linear layer: 2 × input_dim × output_dim FLOPs per token
    - Attention: 4 × n_heads × d_k × seq_len per token (per layer)
    - 2 passes (forward + backward) doubles: multiply by 2
    - Total ≈ 6N per token

GPT-3 training compute:
  N = 175B parameters
  D = 300B tokens
  
  FLOPs = 6 × N × D = 6 × 175B × 300B = 3.15 × 10^23 FLOPs

Converting to GPU-hours:
  A100 80GB: ~312 TFLOPS (FP16) = 3.12 × 10^14 FLOPS
  
  GPU-hours = 3.15 × 10^23 / (3.12 × 10^14) / 3600 ≈ 280,000 GPU-hours
  
  At 1,024 A100s: 280,000 / 1,024 ≈ 273 hours ≈ 11.4 days
  
  Cost at $3/GPU-hour (2020 cloud pricing):
  280,000 × $3 = $840,000 ≈ $1M for just the final run
  
  With experiments, ablations, failed runs: ~$4.6M total (estimated)

LARGER MODELS:
  GPT-4 (estimated): $50-100M training run
  GPT-5 (estimated): $500M+
  
  Frontier model training is now a capital-intensive business.
```

---

## 4. Supervised Fine-tuning (SFT)

```
SFT OVERVIEW:
  After pre-training, the model knows language but doesn't "want" to help.
  SFT teaches it the instruction-following FORMAT.
  
INSTRUCTION DATASETS:
  
  FLAN (Google, 2021):
    1.8K tasks, 62 datasets formatted as instructions
    Templates: "Classify the sentiment of this review: {review}"
    Key finding: "instruction tuning" (training on task descriptions)
    dramatically improves zero-shot generalization
    
  Alpaca (Stanford, 2023):
    52K instruction pairs generated using GPT-3 ("self-instruct")
    Cost: ~$500 to generate using OpenAI API
    Quality: okay; some hallucinations from GPT-3 generator
    
  OpenHermes (2023):
    900K+ instruction pairs from various sources
    Higher quality than Alpaca
    
  UltraChat (2023):
    1.5M multi-turn conversations
    Covers Q&A, creative writing, task solving
    
  ORCA (Microsoft, 2023):
    GPT-4 generated "reasoning traces" — teaches step-by-step thinking
    Models trained on ORCA show much better reasoning

SFT TRAINING PROCEDURE:
  
  DATA FORMAT:
    System: "You are a helpful AI assistant."
    User: "Explain how photosynthesis works."
    Assistant: "Photosynthesis is the process by which plants..."
    
  LOSS:
    Only compute loss on ASSISTANT turns (not system/user tokens)
    System + User tokens → teacher-forced in, loss masked to 0
    
    This is crucial: we want the model to learn to generate good responses,
    not to "predict" the user's questions.
  
  HYPERPARAMETERS:
    Learning rate: 2e-5 (much lower than pre-training, e.g., 1e-4)
    Epochs: 1-3 (more → overfit to training format, lose generality)
    Batch size: 128-512 tokens per example
    Warmup: 3-5% of steps
    
  TYPICAL RESULT:
    SFT alone produces decent instruction-following behavior.
    But raw SFT models can still be verbose, sycophantic, or harmful.
    RLHF refines the behavior further.
```

---

## 5. RLHF: Reinforcement Learning from Human Feedback

### Step 1: Collect Human Preference Data

```
PREFERENCE DATA COLLECTION:
  
  Setup:
    Human annotators (contractors) see:
    - A PROMPT
    - TWO (or more) model COMPLETIONS for that prompt
    
    Annotators rank: Which completion is better?
    
  Example annotation:
    Prompt: "Write a Python function to reverse a string."
    
    Completion A:
      def reverse_string(s):
          return s[::-1]
    
    Completion B:
      Here is a Python function to reverse a string.
      The function takes a string as input and returns
      it reversed using Python's slice notation.
      def reverse(text):
          result = ""
          for char in text:
              result = char + result
          return result
      This function... [continues for 200 more words]
    
    Human judgment: A > B (concise, correct, Pythonic)
    
  WHAT ANNOTATORS RATE:
    - Helpfulness: Does it answer the question?
    - Harmlessness: Is it safe/appropriate?
    - Honesty: Is it accurate?
    - Format: Right length, appropriate style?
    
  ANNOTATION CHALLENGES:
    - Annotator subjectivity (humans disagree ~25% of the time)
    - Sycophancy: longer = better bias, more confident = better bias
    - Domain knowledge: annotators may not know if code is correct
    - Cost: good annotators are expensive ($15-30/hour)
```

### Step 2: Train Reward Model

```
REWARD MODEL TRAINING:
  
  Architecture:
    Same as the SFT model, but replace final language modeling head
    with a scalar head: Linear(d_model, 1)
    
    Input:  (prompt, completion) concatenated
    Output: scalar reward r ∈ R (higher = better completion)
  
  TRAINING OBJECTIVE (Bradley-Terry preference model):
    Given a pair (y_w, y_l) for prompt x (y_w = winner, y_l = loser):
    
    Loss = -log σ(r(x, y_w) - r(x, y_l))
    
    Where σ = sigmoid function.
    
    This maximizes the probability that the preferred completion
    gets higher reward.
    
  WHY THIS WORKS:
    σ(r_w - r_l) = probability that y_w is preferred
    Maximizing log of this probability = learning to rank correctly
    
    At convergence: r(x, y_w) >> r(x, y_l) for all training pairs
    The model generalizes to unseen prompts!
    
  REWARD MODEL QUALITY:
    Check: does it correctly rank held-out pairs?
    Typical accuracy: 60-70% (humans only agree ~75%, so near ceiling)
```

### Step 3: PPO Fine-tuning

```
PPO (PROXIMAL POLICY OPTIMIZATION):
  Most complex step: uses RL to update the language model.
  
  TERMINOLOGY MAPPING (RL → LLM):
    Policy π:           the language model we're training
    State s:            the current context (prompt + generated so far)
    Action a:           generating the next token
    Reward R:           score from reward model (at end of response)
    
  THE TRAINING LOOP:
    
  ┌─────────────────────────────────────────────────────────┐
  │ REPEAT for each batch:                                  │
  │                                                         │
  │ 1. SAMPLE prompts from prompt dataset                   │
  │                                                         │
  │ 2. GENERATE responses with CURRENT LLM (policy π_θ)    │
  │    (Using sampling — temperature=0.9 typically)         │
  │                                                         │
  │ 3. SCORE with reward model: r = R(prompt, response)     │
  │                                                         │
  │ 4. COMPUTE KL PENALTY:                                  │
  │    KL(π_θ(y|x) || π_ref(y|x))                          │
  │    (measure how far current model is from SFT model)    │
  │                                                         │
  │ 5. COMPUTE TOTAL REWARD:                                │
  │    r_total = r - β × KL_penalty                        │
  │    β typically = 0.01 - 0.1                             │
  │                                                         │
  │ 6. UPDATE LLM with PPO to maximize r_total              │
  │    PPO clips the policy update to prevent too-large steps│
  │                                                         │
  └─────────────────────────────────────────────────────────┘
  
  WHY THE KL PENALTY IS CRITICAL:
    Without it, the model would "reward hack":
    
    Reward hacking examples:
      - Summarization RM trained on human prefs:
        Model learns to paste the original text verbatim → perfect recall
        RM: "high similarity" → high reward
        Quality: terrible (not a summary!)
        
      - Helpful assistant RM:
        Model learns to add excessive qualifications and caveats
        RM: "seems careful and thoughtful" → high reward
        Quality: unhelpful hedge-filled responses
        
      - Safety RM:
        Model refuses everything: "I cannot help with that"
        RM: "never says anything harmful" → high reward
        Quality: completely useless
        
    KL penalty prevents too-large deviation from SFT model.
    Keeps the model "honest" — stays close to what it learned from humans.
    
  PPO CLIP OBJECTIVE:
    Standard policy gradient:
      maximize E[A_t × ∇log π_θ(a_t|s_t)]
      (increase probability of actions with positive advantage)
    
    PPO:
      maximize E[min(
        r_t(θ) × A_t,
        clip(r_t(θ), 1-ε, 1+ε) × A_t
      )]
      Where r_t(θ) = π_θ(a|s) / π_old(a|s)  (importance ratio)
      ε = 0.2 (clip range)
    
    The clip prevents any single update from changing the policy too much.
    Stabilizes training considerably vs vanilla policy gradient.
```

### RLHF Challenges

```
RLHF CHALLENGES IN PRACTICE:
  
  1. REWARD MODEL OVERFITTING:
     The policy (LLM) gets very good at maximizing the RM's score.
     RM may not perfectly represent human preferences → "Goodhart's Law"
     After too much RLHF: model score high on RM but quality drops.
     Fix: stop training before RM "saturates", monitor human evals
  
  2. KL TERM TUNING:
     β too small → reward hacking
     β too large → model barely changes from SFT baseline (RLHF does nothing)
     Finding the sweet spot requires careful experimentation
  
  3. SYCOPHANCY:
     Human annotators rate confident, agreeable answers higher.
     Model learns to agree with the user even when wrong.
     Model changes its answer if user pushes back, even if originally correct.
     Mitigation: diverse annotators, specific annotation guidelines
  
  4. PPO INSTABILITY:
     PPO is notoriously finicky to tune.
     Hyperparameters (clip ε, value function coefficient, entropy bonus)
     all interact in complex ways.
     Each new model/task requires re-tuning.
  
  5. COMPUTE COST:
     RLHF is ~3× more expensive than SFT (need to run LLM to generate,
     run RM to score, then update — vs just SFT forward+backward)
```

---

## 6. DPO: Direct Preference Optimization

DPO (Rafailov et al., 2023) eliminates the need for a separate reward model and PPO:

```
DPO MOTIVATION:
  RLHF is complex: train SFT → train RM → PPO RL loop.
  Three stages, each with their own hyperparameters and instabilities.
  
  Can we directly optimize on preference pairs without RL?
  
  Answer: Yes, because RLHF has a closed-form optimal solution
  that can be directly optimized!
  
THE DPO OBJECTIVE:
  Given preference dataset: {(x, y_w, y_l)} 
    x = prompt
    y_w = preferred (winning) completion
    y_l = dispreferred (losing) completion
  
  DPO Loss:
  
  L_DPO = -E[log σ(
    β log(π_θ(y_w|x) / π_ref(y_w|x)) - 
    β log(π_θ(y_l|x) / π_ref(y_l|x))
  )]
  
  WHERE:
    π_θ = model we're training
    π_ref = reference model (SFT model, frozen)
    β = temperature parameter (how much to deviate from reference)
    σ = sigmoid
  
  INTUITION:
    Maximize the probability that π_θ prefers y_w over y_l.
    Relative to the reference model's preferences.
    
    The β log(π_θ(y)/π_ref(y)) term = implicit reward
    This is exactly the reward the RLHF optimal policy corresponds to!
    
    DPO directly optimizes this without explicitly computing the reward.
    
  WHY IT WORKS:
    RLHF has a mathematically optimal policy:
      π*(y|x) ∝ π_ref(y|x) exp(r*(x,y) / β)
    
    Substituting back into the Bradley-Terry preference model:
      The reward cancels out and you get exactly the DPO objective.
    
    DPO ≡ RLHF under the Bradley-Terry model, without RL!
    
  IMPLEMENTATION:
    def dpo_loss(model, ref_model, prompt, chosen, rejected, beta=0.1):
        # Log probabilities under current model
        logp_chosen  = get_sequence_logprob(model,     prompt, chosen)
        logp_rejected= get_sequence_logprob(model,     prompt, rejected)
        
        # Log probabilities under reference model (frozen)
        with torch.no_grad():
            ref_logp_chosen  = get_sequence_logprob(ref_model, prompt, chosen)
            ref_logp_rejected= get_sequence_logprob(ref_model, prompt, rejected)
        
        # Implicit reward differences
        chosen_reward   = beta * (logp_chosen   - ref_logp_chosen)
        rejected_reward = beta * (logp_rejected - ref_logp_rejected)
        
        # DPO loss
        loss = -F.logsigmoid(chosen_reward - rejected_reward).mean()
        return loss
  
  ADVANTAGES OVER RLHF:
    - No reward model to train (saves 2 training stages)
    - No PPO instability
    - Simpler hyperparameter tuning (just β)
    - Often matches or beats RLHF quality
    
  DISADVANTAGES:
    - Requires offline preference pairs (no online generation)
    - Can't incorporate real-time human feedback
    - Less flexibility than full RL
    
  MODELS USING DPO:
    Zephyr-7B (HuggingFace, 2023): DPO on top of Mistral-7B
    Llama-3-Instruct (Meta, 2024): DPO + additional alignment
    Tulu-2 (Allen AI): DPO + RLHF comparison study
```

---

## 7. Constitutional AI (Anthropic)

```
CONSTITUTIONAL AI (CAI):
  Anthropic's alternative to RLHF that scales without human labeling.
  Introduced in "Constitutional AI: Harmlessness from AI Feedback" (2022).
  
THE PROBLEM WITH RLHF:
  Labeling is expensive and doesn't scale.
  Human annotators are biased, inconsistent.
  Hard to get coverage of rare harmful situations.
  
THE CONSTITUTION IDEA:
  Instead of human labels: define a set of PRINCIPLES ("the constitution")
  that the model should follow.
  
  Anthropic's constitution (simplified):
    - Be helpful: provide useful information and assistance
    - Be harmless: avoid causing harm to users, third parties, society
    - Be honest: don't deceive, don't claim to be human
    - Don't assist with weapons of mass destruction
    - Don't generate CSAM
    - Respect user autonomy
    - Be forthright: share information users would want
    - Don't take actions that concentrate power inappropriately
    ... (several dozen principles total)
  
CAI TRAINING PIPELINE:

PHASE 1: SUPERVISED LEARNING FROM AI FEEDBACK (SL-CAI)
  
  Step 1: Collect harmful completions
    Prompt the SFT model with "red-team" harmful prompts
    Get responses (which may be harmful)
    
  Step 2: AI Critique and Revision
    For each harmful response:
      a. Ask Claude (the critic) to critique based on the constitution:
         "Identify specific ways in which the response is harmful,
          according to principle: 'be harmless'"
      b. Ask Claude to revise the response:
         "Please rewrite the response to remove harmful content,
          keeping it as helpful as possible"
      c. May repeat critique+revision several times
      
  Step 3: Fine-tune on revised completions
    Train on the cleaned (prompt, revised_response) pairs
    
PHASE 2: RLHF FROM AI FEEDBACK (RLAIF)
  
  Step 4: Preference labels from AI
    Show Claude pairs of completions (original vs revised)
    Ask: "Which is better according to [principle]?"
    Get Claude's preference label
    
  Step 5: Train reward model on AI preferences
    Same as RLHF reward model, but AI-labeled data
    
  Step 6: RL training
    Same PPO loop as RLHF, using AI-labeled reward model
    
  THE RESULT:
    CAI-trained model (Claude-1) showed:
      - More harmless than HH-RLHF models
      - Less "over-refusal" (doesn't refuse benign requests)
      - More consistent in how it handles edge cases
      - More transparent about its reasoning (trained to explain)
    
  WHY CAI WORKS:
    1. Principles are explicit → easier to audit and update
    2. AI labels are infinite (generate as many as needed)
    3. AI labels can cover rare/extreme cases cheaply
    4. Constitution can be iterated without re-labeling everything
    5. Model behavior is more predictable (follows explicit rules)
    
  LIMITATIONS:
    - AI feedback quality bounded by quality of the "critic" model
    - Bootstrapping problem: who teaches the first Claude?
    - Constitution itself must be carefully designed
    - Less responsive to user preference nuances than human labels
```

---

## 8. Alignment — The Bigger Picture

```
WHAT IS ALIGNMENT?

"AI Alignment" = ensuring AI systems pursue goals that are consistent with
human values and intentions.

THE ALIGNMENT PROBLEM:
  As AI systems become more capable, small misalignments in their goals
  could have large and possibly irreversible consequences.
  
  Like hiring an employee:
    - Competent employee who disagrees with your values → dangerous
    - Incompetent employee who agrees → ineffective but manageable
    - Competent + well-aligned → ideal
  
  Scale makes this harder: a very powerful AI with even small misalignment
  can find ways to achieve its misaligned goals that humans can't prevent.

GOODHART'S LAW IN AI:
  "When a measure becomes a target, it ceases to be a good measure."
  
  Applied to RLHF:
    Target metric: reward model score
    Problem: model optimizes reward score, not underlying human values
    
    Famous example:
      Boat Racing game (Specification Gaming):
        Objective: win the race
        Agent discovered: spinning in circles collecting powerups
        is more rewarding than actually racing!
        Never learned to actually race.
    
    LLM examples:
      Reward model rewarded: "seems helpful"
      Model learned: verbose, confident-sounding, agreeable responses
      Not the same as: actually correct, genuinely helpful
      
    Reward model rewarded: "doesn't say harmful things"
    Model learned: refuse more requests ("if in doubt, refuse")
    Not the same as: being helpful to legitimate users

CURRENT ALIGNMENT APPROACHES:
  
  1. RLHF / DPO: align to human preferences (works but limited)
  2. Constitutional AI: align to explicit principles (scales better)
  3. Scalable Oversight: use AI to assist humans in evaluating AI
     (needed as AI becomes more capable than human evaluators)
  4. Interpretability: understand what's inside the model to verify alignment
     (Anthropic's "mechanistic interpretability" research)
  5. Debate: two AI systems argue, humans judge the winner
     (aligned AI wins debates with misaligned AI)
  6. Amplification: iteratively improve oversight quality

WHY THIS MATTERS FOR DEVELOPERS:
  Understanding alignment helps you:
  - Know why models refuse certain requests (alignment, not technical limitation)
  - Write better prompts (work with alignment, not against it)
  - Anticipate model behavior in edge cases
  - Understand AI safety news and debates
```

---

## 9. Inference Optimization: KV Cache and Speculative Decoding

### KV Cache

```
THE PROBLEM:
  Autoregressive generation:
  
  Step 1: input = [t1, t2, t3]
    Compute K1,V1 for t1; K2,V2 for t2; K3,V3 for t3
    Attention: t3 attends to K1,K2,K3
    
  Step 2: input = [t1, t2, t3, t4]
    RECOMPUTE K1,V1 for t1; K2,V2 for t2; K3,V3 for t3 ← WASTEFUL!
    Then compute K4,V4 for t4
    
  Without KV cache: O(n²) attention computations for generating n tokens.
  
THE KV CACHE SOLUTION:
  Cache the computed K and V matrices for all previous tokens.
  At each new step, only compute K, V for the NEW token.
  
  Step 1: Compute K1,V1, K2,V2, K3,V3. Cache them.
          Compute output for all 3 positions.
          
  Step 2: Load K1,V1, K2,V2, K3,V3 from cache.
          Compute K4,V4 for new token t4.
          Attend: t4 attends to K1,K2,K3,K4 from cache.
          
  Step 3: Load K1..K4. Compute K5,V5. Attend.
  
  Computation per step: O(n) instead of O(n²)
  Memory: must store K,V for all previous tokens
  
KV CACHE MEMORY:
  Per token, per layer, K and V:
    2 × n_layers × seq_len × n_heads × head_dim × dtype_bytes
    
  LLaMA-7B (FP16, seq_len=4096):
    2 × 32 × 4096 × 32 × 128 × 2 bytes = 2.1 GB
    
  LLaMA-70B (FP16, seq_len=4096):
    2 × 80 × 4096 × 64 × 128 × 2 bytes = 10.7 GB!
    
  For long contexts (128k tokens):
    LLaMA-70B KV cache at 128k: ~334 GB ← needs multiple GPUs just for KV!
    
  This is why long-context models need special memory management.
  
  SOLUTIONS:
    - Grouped Query Attention (GQA): share K,V across query heads → fewer heads to cache
    - Paged Attention (vLLM): manage KV cache like virtual memory
    - Quantized KV cache: INT8 for K,V → 2× memory reduction
    - Eviction policies: remove old/low-importance K,V entries
```

### Speculative Decoding

```
SPECULATIVE DECODING (Chen et al. 2023):
  
  PROBLEM: Large models are slow to generate (75B tokens/GPU/sec is fast).
  
  KEY INSIGHT:
    Most tokens are "easy" — the model is very confident.
    Only a few tokens are "hard" — the model needs to think.
    
    Can we use a SMALL FAST model (draft model) to predict multiple
    tokens, then verify them with the LARGE model?
    
  ALGORITHM:
    
    Setup:
      Large model M (175B params): accurate but slow
      Small model M' (7B params): fast but less accurate
      γ = number of tokens to speculate (e.g., 4)
    
    Each iteration:
    
      Step 1: DRAFT
        Use small model M' to autoregressively generate γ tokens:
        draft_tokens = M'(context)[:γ]
        
        This is fast: small model runs γ forward passes quickly.
      
      Step 2: VERIFY
        Run large model M on (context + draft_tokens) in ONE pass:
        
        Large model outputs logits for positions 1..γ+1 simultaneously.
        
        For each draft token t_i:
          acceptance_prob = min(1, P_M(t_i) / P_M'(t_i))
          Accept t_i with probability acceptance_prob.
          If rejected: sample a corrected token from M and STOP.
        
      Step 3: RESULT
        If all γ draft tokens accepted: γ tokens generated in one M pass!
        If k tokens accepted: k tokens generated + 1 correction.
        
      Expected number of tokens per large-model forward pass:
        ≈ γ × (acceptance_rate) + 1
        
      SPEEDUP:
        Easy text (common phrases): acceptance rate ~90% → ~4.5 tokens/pass
        Hard text (technical/creative): acceptance rate ~50% → ~2.5 tokens/pass
        Random: acceptance rate ~10% → ~1.1 tokens/pass (no speedup)
        
        Typical real-world speedup: 2-3× for code/common text
        
  DRAFT MODELS:
    Same family as main model, small version: LLaMA-7B for LLaMA-70B
    Universal draft model: can use any small model
    Medusa heads: multiple prediction heads on the large model itself
```

---

## 10. Production Serving: vLLM

```
vLLM (2023):
  Open-source LLM serving engine from UC Berkeley.
  
  PROBLEMS IN PRODUCTION:
    - Multiple users sending requests simultaneously
    - Each request needs its own KV cache
    - Requests have different lengths → memory fragmentation
    - Sequential processing: user 2 waits while user 1's request completes
    
  PAGED ATTENTION:
    Inspired by OS virtual memory management.
    
    Traditional KV cache: allocate max_seq_len memory upfront
      Request 1 max_seq_len = 2048 → allocate 2048 slots even if only uses 100
      Fragmented, wasteful
    
    Paged attention:
      Divide KV cache into fixed-size "pages" (e.g., 16 tokens per page)
      Each request's KV cache = list of pages (not contiguous memory!)
      
      Request 1: allocates page 0, page 5, page 12 (as needed, non-contiguous)
      Request 2: allocates page 1, page 3, page 7
      
      No internal fragmentation!
      Maximum memory utilization.
    
  CONTINUOUS BATCHING:
    Traditional: process one request completely, then next.
    Continuous batching: mix tokens from multiple requests in each batch!
    
    TIME STEP 1: [User1_token17, User2_token3, User3_token52]
    TIME STEP 2: [User1_token18, User2_token4, User4_token1]
    
    User 2 finishes at step 5 → User 4 can start at step 6 immediately
    
    RESULT: 10-100× better GPU utilization for serving
    Latency for user 1 unaffected by user 2's long request.
    
  vLLM PERFORMANCE:
    vs naive sequential serving:
      Throughput: 10-24× higher tokens/second
      Latency: similar or lower (due to better batching)
    
    A single A100 GPU can serve ~2,000 users/minute with Mistral-7B
    using vLLM (vs ~200 users/minute with naive serving)
    
  GETTING STARTED:
    pip install vllm
    
    from vllm import LLM, SamplingParams
    
    llm = LLM(model="mistralai/Mistral-7B-Instruct-v0.2")
    
    params = SamplingParams(temperature=0.8, top_p=0.95, max_tokens=200)
    
    outputs = llm.generate([
        "Write a haiku about transformers:",
        "Explain attention in one sentence:",
    ], params)
    
    for output in outputs:
        print(output.outputs[0].text)
```

---

## 11. Summary

```
LLM TRAINING PIPELINE:
═══════════════════════════════════════════════════════════════════════

DATA PIPELINE:
  Raw: Common Crawl, Wikipedia, Books, GitHub, Papers
  Filter: language ID, quality classifiers, deduplication, toxicity
  Mix: 65% web, 8% code, 5% Wikipedia, etc.
  Tokenize: BPE or SentencePiece, 32k-128k vocabulary

DISTRIBUTED TRAINING:
  Data Parallel: replicate model on N GPUs, average gradients
  Tensor Parallel: split layers across GPUs (intra-node, NVLink)
  Pipeline Parallel: split layer stacks across nodes
  ZeRO: partition optimizer states to save memory
  
  Cost: $1M - $100M for frontier models

ALIGNMENT:
  SFT: ~50k instruction examples, fine-tune with language modeling loss
  RLHF: human preferences → reward model → PPO RL loop + KL penalty
  DPO: skip reward model, directly optimize preference pairs
  Constitutional AI: use AI feedback based on written principles

INFERENCE:
  KV Cache: cache K,V for previous tokens → O(n) per step instead of O(n²)
  Speculative Decoding: draft with small model, verify with large → 2-3× speedup
  vLLM: paged attention + continuous batching → 10-24× serving throughput
═══════════════════════════════════════════════════════════════════════
```

---

## Mini Projects

### Mini Project 1: Memory Calculator and Training Cost Estimator

Build an interactive LLM training cost calculator to understand the resource requirements.

**Objective:** Make GPU memory math concrete — stop guessing, start calculating.

```python
import numpy as np
import matplotlib.pyplot as plt

def bytes_to_human(n):
    for unit in ['B', 'KB', 'MB', 'GB', 'TB']:
        if n < 1024: return f"{n:.1f} {unit}"
        n /= 1024
    return f"{n:.1f} PB"

def training_memory_estimate(
    n_params_B,          # model size in billions of parameters
    precision='bf16',    # 'fp32', 'fp16', 'bf16', 'int8', 'int4'
    batch_size=1,
    seq_len=2048,
    hidden_dim=4096,
    n_layers=32,
    include_optimizer=True,  # Adam states
    gradient_checkpointing=False
):
    n_params = n_params_B * 1e9
    bytes_per_param = {'fp32': 4, 'fp16': 2, 'bf16': 2, 'int8': 1, 'int4': 0.5}[precision]

    # Model weights
    model_mem = n_params * bytes_per_param

    # Optimizer states (Adam: 2x FP32 copies of params)
    optimizer_mem = n_params * 4 * 2 if include_optimizer else 0  # FP32 momentum + variance

    # Gradients (same dtype as weights usually, but at least FP16)
    grad_mem = n_params * bytes_per_param

    # Activations (rough estimate: seq_len * hidden * n_layers * 4 bytes per float)
    if gradient_checkpointing:
        # Only store activations at layer boundaries, recompute rest
        act_factor = np.sqrt(n_layers)
    else:
        act_factor = n_layers
    act_mem = batch_size * seq_len * hidden_dim * act_factor * 4  # fp32 activations

    total = model_mem + optimizer_mem + grad_mem + act_mem

    return {
        'model': model_mem,
        'optimizer': optimizer_mem,
        'gradients': grad_mem,
        'activations': act_mem,
        'total': total,
        'total_GB': total / 1e9,
    }

# Common LLM configurations
models = [
    ("GPT-2 (117M)",    0.117, 768,  12, 'fp32'),
    ("GPT-2-XL (1.5B)", 1.5,   1600, 48, 'bf16'),
    ("LLaMA-7B",        7.0,   4096, 32, 'bf16'),
    ("LLaMA-13B",       13.0,  5120, 40, 'bf16'),
    ("LLaMA-70B",       70.0,  8192, 80, 'bf16'),
    ("GPT-4 (~1.8T)",   1800,  12288,96, 'bf16'),  # estimated
]

# GPU specs
gpus = {
    'RTX 4090 (24GB)': 24,
    'A100 (40GB)': 40,
    'A100 (80GB)': 80,
    'H100 (80GB)': 80,
    'H200 (141GB)': 141,
}

fig, axes = plt.subplots(2, 2, figsize=(16, 11))
fig.suptitle("LLM Training: Memory and Compute Analysis", fontsize=14, fontweight='bold')

# Memory breakdown for each model
categories = ['Model Weights', 'Optimizer States', 'Gradients', 'Activations']
bar_colors  = ['steelblue', 'orange', 'green', 'purple']

model_names_short = [m[0].split('(')[0].strip() for m in models[:5]]
memory_breakdowns = [training_memory_estimate(m[1], m[4], hidden_dim=m[2], n_layers=m[3])
                      for m in models[:5]]

x = np.arange(len(models[:5]))
width = 0.6
bottoms = np.zeros(len(models[:5]))
for i, (cat, color) in enumerate(zip(categories, bar_colors)):
    keys = ['model', 'optimizer', 'gradients', 'activations']
    vals = [m[keys[i]] / 1e9 for m in memory_breakdowns]
    axes[0, 0].bar(x, vals, width, bottom=bottoms, label=cat, color=color, alpha=0.8)
    bottoms += np.array(vals)

axes[0, 0].set_xticks(x)
axes[0, 0].set_xticklabels(model_names_short, rotation=15, ha='right', fontsize=8)
axes[0, 0].set_ylabel("Memory (GB)")
axes[0, 0].set_title("Training Memory by Component (batch=1, seq=2048)")
axes[0, 0].legend(fontsize=8); axes[0, 0].grid(True, alpha=0.3, axis='y')
axes[0, 0].set_yscale('log')

# GPU count required
gpu_mem = {'A100-40GB': 40, 'A100-80GB': 80, 'H100-80GB': 80, 'RTX-4090': 24}
for gpu_name, gpu_gb in gpu_mem.items():
    n_gpus_needed = [max(1, int(np.ceil(m['total_GB'] / gpu_gb))) for m in memory_breakdowns]
    axes[0, 1].plot(model_names_short, n_gpus_needed, marker='o', linewidth=2,
                    label=gpu_name, markersize=6)

axes[0, 1].set_yscale('log')
axes[0, 1].set_ylabel("GPUs Required (log scale)"); axes[0, 1].set_xlabel("")
axes[0, 1].set_title("Minimum GPU Count for Training\n(no model parallelism)")
axes[0, 1].set_xticklabels(model_names_short, rotation=15, ha='right', fontsize=8)
axes[0, 1].legend(fontsize=7); axes[0, 1].grid(True, alpha=0.3)

# Gradient checkpointing savings
n_params_range = np.logspace(8, 11, 50)  # 100M to 100B
for hidden in [768, 4096, 8192]:
    n_layers = max(12, int(hidden / 100))
    regular_act = [1 * 2048 * hidden * n_layers * 4 / 1e9 for _ in n_params_range]
    ckpt_act    = [1 * 2048 * hidden * np.sqrt(n_layers) * 4 / 1e9 for _ in n_params_range]
    axes[1, 0].plot(n_params_range/1e9, regular_act, '--', linewidth=1.5,
                    label=f'd={hidden} no-ckpt', alpha=0.7)
    axes[1, 0].plot(n_params_range/1e9, ckpt_act, '-', linewidth=1.5,
                    label=f'd={hidden} ckpt', alpha=0.7)

axes[1, 0].set_xscale('log'); axes[1, 0].set_yscale('log')
axes[1, 0].set_xlabel("Model Size (B params)"); axes[1, 0].set_ylabel("Activation Memory (GB)")
axes[1, 0].set_title("Gradient Checkpointing:\nActivation Memory Savings")
axes[1, 0].legend(fontsize=6, ncol=2); axes[1, 0].grid(True, alpha=0.3)

# Cost estimation (approximate cloud pricing)
gpu_cost_per_hr = {'A100-40GB': 3.00, 'A100-80GB': 4.50, 'H100-80GB': 8.00}
tokens_per_gpu_per_second = {'A100-40GB': 800, 'A100-80GB': 1200, 'H100-80GB': 2500}
training_tokens = 1e12  # 1T tokens (LLaMA-2 style)

axes[1, 1].axis('off')
cost_text = "Cost Estimate: 1T Token Training\n"
cost_text += "=" * 45 + "\n\n"
for (gname, cost_hr), (_, tok_s) in zip(gpu_cost_per_hr.items(), tokens_per_gpu_per_second.items()):
    for model_name, n_params, *_ in models[:4]:
        n_gpus = max(1, int(np.ceil(
            training_memory_estimate(n_params, 'bf16')['total_GB'] / 80
        )))
        time_hrs = training_tokens / (tok_s * n_gpus * 3600)
        cost_usd  = time_hrs * n_gpus * cost_hr
        cost_text += f"{model_name[:12]:12s} | {gname[:10]:10s} | "
        cost_text += f"{n_gpus:3d} GPUs | ${cost_usd:,.0f}\n"
    cost_text += "\n"

axes[1, 1].text(0.02, 0.98, cost_text, transform=axes[1, 1].transAxes, fontsize=7,
                va='top', fontfamily='monospace',
                bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))
axes[1, 1].set_title("Training Cost Estimates (Approximate)")

plt.tight_layout()
plt.savefig("llm_memory_costs.png", dpi=150)
plt.show()

# Print summary table
print("\nMemory Summary (batch=1, seq=2048, bf16, with Adam):")
print(f"{'Model':<20} {'Total GB':>10} {'# A100-80GB':>12}")
print("-" * 45)
for (name, n_params, hidden, n_layers, prec) in models[:5]:
    mem = training_memory_estimate(n_params, prec, hidden_dim=hidden, n_layers=n_layers)
    n_gpus = max(1, int(np.ceil(mem['total_GB'] / 80)))
    print(f"{name:<20} {mem['total_GB']:>10.1f} {n_gpus:>12}")
```

---

## 12. Exercises

**Exercise 1**: Calculate the memory required to train a Mistral-7B model (7B params, 32 layers, 32 heads, head_dim=128) using mixed precision (FP16 weights + FP32 optimizer states). How many A100 80GB GPUs are needed just to fit the model and optimizer? With ZeRO-3 on 8 GPUs, how much memory per GPU?

**Exercise 2**: Implement a simple reward model using DistilBERT. Create a synthetic dataset of prompt-response pairs where "good" responses are concise and factual, "bad" responses are verbose and irrelevant. Train the reward model and evaluate its preference accuracy.

**Exercise 3**: Implement the DPO loss function from scratch. Use a small GPT-2 model as both the reference model and the model to be trained. Create 20 preference pairs (for a simple task like sentiment response) and run 10 DPO training steps. Observe how the model's log probabilities change.

**Exercise 4**: Implement a basic KV cache for a small transformer. Measure generation speed with and without KV cache for different sequence lengths. Plot the speedup ratio.

**Exercise 5**: Research Constitutional AI further. Write your own "constitution" (list of 10 principles) for a specific use case (e.g., a coding assistant, a medical information bot). Which principles would be most important? Which would be hardest to teach via RLHF vs CAI?

**Exercise 6**: Using vLLM (if GPU available) or the HuggingFace Inference API, compare the throughput (tokens/second) of serving Mistral-7B with batch_size=1 vs batching 16 requests. What is the speedup from batching?

---

**Chapter Summary**: Training an LLM is a pipeline spanning data curation (petabytes of filtered, deduplicated web text), distributed training (3D parallelism — data + tensor + pipeline — with ZeRO memory optimization), and alignment (SFT on instruction data → RLHF with reward model + PPO, or DPO directly on preference pairs). Constitutional AI (Anthropic) replaces human labels with AI critique based on explicit principles. The alignment problem — ensuring models pursue human-compatible goals — becomes more critical as capability increases, with Goodhart's Law ensuring reward metrics will be gamed if not carefully managed. At inference, KV caching eliminates redundant computation, speculative decoding uses a fast draft model for 2-3× speedup, and production serving with vLLM uses paged attention and continuous batching for 10-24× throughput improvement.

**What's Next →** [Chapter 27: Tokenization — How Text Becomes Numbers](./27-tokenization.md)

*Before we can build our own language model, we need to understand the foundational step: how raw text gets converted into the integer sequences that transformers actually process.*
