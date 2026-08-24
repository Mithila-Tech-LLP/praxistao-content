# Chapter 25: GPT and Autoregressive Language Models

> **"We started building language models to understand language. Somewhere along the way, they started understanding everything else too."**

---

## Table of Contents
1. [The Language Modeling Objective](#1-the-language-modeling-objective)
2. [Autoregressive Generation](#2-autoregressive-generation)
3. [GPT-1: The Foundation](#3-gpt-1-the-foundation)
4. [GPT-2: Scaling and Controversy](#4-gpt-2-scaling-and-controversy)
5. [GPT-3: The First "It Just Works" LLM](#5-gpt-3-the-first-it-just-works-llm)
6. [Scaling Laws](#6-scaling-laws)
7. [InstructGPT and ChatGPT: From Completion to Assistant](#7-instructgpt-and-chatgpt-from-completion-to-assistant)
8. [GPT-4 and Beyond](#8-gpt-4-and-beyond)
9. [Emergent Abilities](#9-emergent-abilities)
10. [Open-Source Alternatives](#10-open-source-alternatives)
11. [Context Window Evolution](#11-context-window-evolution)
12. [Sampling Strategies](#12-sampling-strategies)
13. [HuggingFace: GPT-2 Generation](#13-huggingface-gpt-2-generation)
14. [Summary](#14-summary)
15. [Exercises](#15-exercises)

---

## 1. The Language Modeling Objective

A **language model** assigns a probability to a sequence of tokens. Formally:

```
P(w_1, w_2, ..., w_n) — probability of the whole sequence

By the chain rule of probability:
P(w_1, ..., w_n) = P(w_1) · P(w_2|w_1) · P(w_3|w_1,w_2) · ... · P(w_n|w_1,...,w_{n-1})
               = Π_{t=1}^{n} P(w_t | w_1, ..., w_{t-1})

TRAINING OBJECTIVE (maximum likelihood):
  Maximize the probability of the training data under the model.
  Equivalently, minimize the cross-entropy loss:
  
  L = -1/n Σ_{t=1}^{n} log P(w_t | w_1, ..., w_{t-1})

In plain English:
  At each position t, given all previous words, predict the NEXT word.
  Maximize the probability assigned to the correct next word.
```

This is beautiful in its simplicity:
- **No labels required**: the data labels itself — the "label" at position t is the word at position t+1
- **Universal**: every piece of text in the world is a valid training example
- **Scalable**: the internet contains trillions of words

The transformer's causal mask (Chapter 23) makes this efficient: compute predictions for all positions simultaneously, with the mask ensuring position t cannot see position t+1 or later.

```
TRAINING EFFICIENCY:

Sequence: "The cat sat on the mat"
           w1   w2  w3  w4   w5  w6

One forward pass computes predictions:
  Position 0 ("The"):  predict w2 = "cat"
  Position 1 ("cat"):  predict w3 = "sat"
  Position 2 ("sat"):  predict w4 = "on"
  Position 3 ("on"):   predict w5 = "the"
  Position 4 ("the"):  predict w6 = "mat"
  Position 5 ("mat"):  predict w7 = [EOS]

All 6 predictions computed in ONE forward pass (thanks to causal mask).
Compare to RNN: 6 sequential forward passes.
```

---

## 2. Autoregressive Generation

At inference time (generating new text), a language model works by repeatedly predicting the next token and appending it to the context:

```
AUTOREGRESSIVE GENERATION:

Step 0: Prompt = "Once upon a time"
         Tokens = [Once, upon, a, time]

Step 1: Model predicts next token
         Input:  [Once, upon, a, time]
         Output: probability over vocabulary
         Sample: "there" (prob 0.15, sampled)
         New context: [Once, upon, a, time, there]

Step 2: Input:  [Once, upon, a, time, there]
         Output: probability over vocabulary
         Sample: "was" (prob 0.42)
         New context: [Once, upon, a, time, there, was]

Step 3: Sample "a" (prob 0.31)
         ...

Continue until [EOS] token or max_length reached.

RESULT: "Once upon a time there was a..."

COST: N forward passes to generate N tokens.
  For 100-token completion: 100 forward passes through the model.
  Each forward pass: O(n² · d) where n grows each step.
  This is why generation is slow for large models.
  (KV Cache in Chapter 26 mitigates this.)
```

---

## 3. GPT-1: The Foundation

**GPT-1** (Generative Pre-Training) was published by OpenAI in June 2018 — the same year as BERT, and in many ways the "other side of the coin."

```
GPT-1 SPECIFICATIONS (2018):
  Architecture:  Transformer decoder-only (12 layers)
  Parameters:    117 million
  Training data: BooksCorpus (~800M words, long-form fiction)
  Context:       512 tokens
  d_model:       768
  n_heads:       12
  Training:      100 epochs of BooksCorpus
  GPU time:      ~1 month on 8 GPUs
```

### Pre-training → Fine-tuning Paradigm

GPT-1 introduced the pre-train→fine-tune approach that BERT later made famous:

```
PHASE 1: UNSUPERVISED PRE-TRAINING
  Objective: predict next token on BooksCorpus
  Result: model learns grammar, style, factual knowledge, reasoning patterns

PHASE 2: SUPERVISED FINE-TUNING
  For each downstream task, add a linear output head.
  Training: minimize task loss (e.g., classification cross-entropy)
  
  Key insight: using LM loss as AUXILIARY objective during fine-tuning helps:
    Total loss = L_task + λ · L_lm
    The LM loss prevents "catastrophic forgetting" of pre-trained features.
```

GPT-1 showed that:
1. Pre-training on unlabeled text → useful features for labeled tasks
2. A SINGLE model can be fine-tuned for many tasks
3. Larger models + more data → better results

But GPT-1's impact was modest compared to BERT (which appeared 4 months later). The real leap came with GPT-2.

---

## 4. GPT-2: Scaling and Controversy

**GPT-2** was released by OpenAI in February 2019. It was the first model where the public (and OpenAI itself) was surprised by the output quality.

```
GPT-2 SPECIFICATIONS (2019):
  Architecture:  Transformer decoder-only
  Parameters:    1.5 billion (largest model)
    - GPT-2 Small: 117M (same as GPT-1)
    - GPT-2 Medium: 345M
    - GPT-2 Large: 762M
    - GPT-2 XL: 1.5B
  Training data: WebText (~40GB, from Reddit posts with 3+ karma)
  Context:       1,024 tokens
  d_model:       1600 (XL)
  n_heads:       25 (XL)
  n_layers:      48 (XL)
```

### WebText: The Data Breakthrough

```
WEBTEXT CONSTRUCTION:
  - Crawled URLs from Reddit posts with ≥3 upvotes ("karma")
  - Only outbound links (not reddit.com links)
  - Removed Wikipedia (common benchmark contamination)
  - 45M links → 8M documents → 40GB text
  
  The Reddit karma filter was a cheap "human quality filter":
  Links that humans upvote tend to be well-written, interesting content.
  Better than random web crawl.
```

### Zero-Shot and Few-Shot Generalization

The big surprise with GPT-2: it generalized to tasks **without any fine-tuning**.

```
ZERO-SHOT GENERALIZATION:

Without any task-specific training:

GPT-2 on WMT14 English-French translation:
  Just ask: "Translate to French: 'The cat sat on the mat' ="
  Output: "Le chat s'est assis sur le tapis"
  BLEU score: 11.5 (not great, but remarkable for zero-shot!)

GPT-2 on reading comprehension (CoQA):
  Provide context + questions in the prompt.
  Output: answers extracted from context.
  Score: 55% F1 (comparable to some supervised models from 2 years prior!)

GPT-2 on natural language inference (RTE):
  Frame as text completion: "Premise. Hypothesis. True or False?"
  Score: 56.4%

None of these tasks were in the training data (format-wise).
The model was just completing text naturally.
```

### The Release Controversy

```
GPT-2 AND AI SAFETY DISCOURSE:

OpenAI's decision (February 2019):
  "Due to our concerns about malicious applications of the technology,
  we are not releasing the trained model. We are instead releasing
  a much smaller model..."
  
  Released GPT-2-small (117M) → not GPT-2 XL (1.5B)
  
Arguments for staged release:
  - Potential for mass generation of fake news, spam
  - Coordinated influence campaigns
  - Easier to create deepfake text than deepfake video
  
Arguments against staged release restriction:
  - Many research groups could train similar models
  - Academic research needs access
  - Security through obscurity doesn't work
  - Smaller model barely worse than full model
  
Outcome:
  - Full 1.5B model released ~9 months later (November 2019)
  - World did not end
  - Established pattern: OpenAI → slow/closed release
  - Also established: AI safety as a serious consideration in releases
  - This controversy shaped the discourse around Claude, Gemini safety policies
```

---

## 5. GPT-3: The First "It Just Works" LLM

**GPT-3** was published by OpenAI in May 2020. This was the first model where people widely felt that AI had crossed some qualitative threshold.

```
GPT-3 SPECIFICATIONS (2020):
  Parameters:    175 billion
  Training data: ~570GB (filtered Common Crawl 60%, WebText2 22%,
                          Books1 8%, Books2 8%, Wikipedia 3%)
  Context:       2,048 tokens
  d_model:       12,288
  n_heads:       96
  n_layers:      96
  d_ff:          49,152 (4 × 12288)
  Training:      ~$4.6 million compute cost (300B token-steps on V100s)
  Batch size:    3.2 million tokens (massive!)
```

### In-Context Learning

The most important capability that emerged from scale in GPT-3:

```
IN-CONTEXT LEARNING (ICL):

Instead of gradient updates (fine-tuning), provide examples IN THE PROMPT.
The model "learns" from prompt examples without updating any weights.

ZERO-SHOT (no examples):
  Prompt: "Translate English to French: 'sea otter' =>"
  Output: "loutre de mer"

ONE-SHOT (one example):
  Prompt: "Translate English to French:
           sea otter => loutre de mer
           peppermint =>"
  Output: "menthe poivrée"

FEW-SHOT (several examples):
  Prompt: "Translate English to French:
           sea otter => loutre de mer
           plush giraffe => girafe en peluche
           cheese => fromage
           cat =>"
  Output: "chat"  ← with fewer errors than zero-shot

PERFORMANCE IMPROVEMENT with examples:
  Task             | 0-shot | 1-shot | Few-shot
  TriviaQA         |  64.3% |  68.0% |   71.2%
  WebQS            |  14.4% |  25.3% |   41.5%
  CoQA             |  81.5% |  84.0% |   85.0%
  
Examples in the prompt = dramatic improvements!
```

### Why Does In-Context Learning Work?

This is still an active research question, but the best current understanding:

```
HYPOTHESIS 1: "Pattern Matching at Scale"
  The model has seen so many text patterns during training that it can
  recognize "this prompt format means: do task X".
  Few-shot examples are just establishing the format.
  
HYPOTHESIS 2: "Bayesian Inference"
  The model maintains implicit probability distributions.
  Examples update the model's implicit "prior" about the task.
  
HYPOTHESIS 3: "In-Context Gradient Descent" (Akyürek et al. 2022)
  The attention mechanism can implement gradient descent implicitly.
  Each few-shot example changes the "effective model" via attention.
  
THE STRANGE PART:
  ICL works even when labels in examples are wrong!
  If you give examples with flipped labels, performance is still good.
  This suggests the format/structure matters more than the content.
  
  "The format tells the model what kind of task to do,
  the examples don't provide as much information as we thought."
```

### Chain-of-Thought at Scale

```
CHAIN-OF-THOUGHT PROMPTING (Wei et al. 2022):

Standard few-shot:
  Q: Roger has 5 tennis balls. He buys 2 more cans of 3 balls each.
     How many tennis balls does he have?
  A: 11

Chain-of-thought few-shot:
  Q: Roger has 5 tennis balls. He buys 2 more cans of 3 balls each.
     How many tennis balls does he have?
  A: Roger started with 5 balls. 2 cans × 3 balls = 6 balls.
     5 + 6 = 11. The answer is 11.

Performance on GSM8K math benchmark:
  Standard few-shot (GPT-3 175B):     17% accuracy
  Chain-of-thought few-shot:           48% accuracy
  
Chain-of-thought only works at scale (>100B params).
Smaller models trying CoT actually get worse.
This suggests CoT is an emergent capability of large scale.
```

---

## 6. Scaling Laws

### Kaplan et al. 2020: The Original Scaling Laws

OpenAI researchers found that language model performance follows remarkably regular **power laws** with respect to scale:

```
SCALING LAW FORM:

L(N) = (N_c / N)^α_N    — loss as a function of model size N
L(D) = (D_c / D)^α_D    — loss as a function of dataset size D
L(C) = (C_c / C)^α_C    — loss as a function of compute C

Where:
  L = cross-entropy loss (lower = better)
  N = model parameters
  D = dataset tokens
  C = compute (FLOPs)
  α values ≈ 0.076 for N, ≈ 0.095 for D

In log-log space, this is a LINE:
  log(L) = α · log(N_c/N) = const - α · log(N)
  
EMPIRICALLY:
  Double model size → loss drops by 2^(-0.076) ≈ 5%
  10× model size → loss drops by 10^(-0.076) ≈ 15%
  
This is PREDICTABLE. You can forecast a model's performance before training!
```

### The Compute Allocation Question

```
GIVEN A FIXED COMPUTE BUDGET C = N × D × 6 FLOPs/token:
(the "6 FLOPs per parameter per token" is approximate for transformers)

How should you split between N (model size) and D (data size)?

KAPLAN et al. (2020): Scale N much more than D
  Optimal: N ∝ C^0.73, D ∝ C^0.27
  Interpretation: Make the model BIGGER, don't worry much about data.
  
  GPT-3 (175B params, 300B tokens) followed this.
  
CHINCHILLA SCALING LAWS (Hoffmann et al. 2022):
  Kaplan et al. were WRONG! Models were undertrained.
  
  Chinchilla: optimal is N ∝ D (equal scaling)
  Rule of thumb: ~20 tokens per parameter
  
  GPT-3 (175B params) needed 175B × 20 = 3.5 TRILLION tokens!
  It was trained on only 300B tokens → severely undertrained!
  
  A 70B model trained on 1.4T tokens (Chinchilla) beats a
  175B model trained on 300B tokens (GPT-3 setting).
```

### Chinchilla Optimal Models

```
CHINCHILLA TRAINING RUNS (to find optimal):
  - Train 70 models of varying N and D, with fixed C budget
  - Measure loss
  - Fit scaling law
  
OPTIMAL COMPUTE BUDGET ALLOCATION (Chinchilla):
  
  Budget C (FLOPs) | Optimal N | Optimal D
  ─────────────────────────────────────────
  1e18 FLOPs      | 9.7B      | 200B tokens
  1e19 FLOPs      | 30.7B     | 615B tokens
  1e20 FLOPs      | 97.2B     | 1.9T tokens
  1e21 FLOPs      | 308B      | 6.1T tokens
  1e22 FLOPs      | 975B      | 19.5T tokens

REAL MODELS AFTER CHINCHILLA:
  LLaMA 7B:  trained on 1T tokens (142 tokens/param)
  LLaMA 13B: trained on 1T tokens (77 tokens/param)
  LLaMA 65B: trained on 1.4T tokens (21 tokens/param) ← Chinchilla-optimal
  
  Mistral 7B: 7B params × 20 = should use ~140B tokens
```

---

## 7. InstructGPT and ChatGPT: From Completion to Assistant

### The Problem with Raw Language Models

```
RAW GPT-3 AS A "HELPFUL ASSISTANT":

User: "What is the capital of France?"

Expected output: "Paris"

Actual GPT-3 output (as a language model):
  "What is the capital of France?
   What is the capital of Germany?
   What is the capital of Spain?
   What is the capital of Italy?..."
  
  ← It completes the text as if it's a list of quiz questions!
  
Or:
  "What is the capital of France?
   This is a commonly asked geography question that appears on..."
  
  ← Treats it as the beginning of a news article!

A raw language model predicts the most likely continuation of text.
It doesn't "know" the user wants a helpful answer.
```

### RLHF: Teaching Models to Follow Instructions

**InstructGPT** (Ouyang et al. 2022) introduced the RLHF training approach:

```
RLHF PIPELINE (3 steps):

STEP 1: Supervised Fine-Tuning (SFT)
  - Hire contractors to write (instruction, ideal_response) pairs
  - ~13,000 demonstrations for InstructGPT
  - Fine-tune GPT-3 on these: learn to follow instruction format
  
  Examples:
    Instruction: "List 5 ideas for a birthday party"
    Response: "1. Outdoor picnic... 2. Movie marathon..."
    
    Instruction: "Explain photosynthesis like I'm 10"
    Response: "Plants are like little food factories..."
  
  Result: GPT-3-SFT — better at following instructions, but limited by
          the number of demonstrations we can collect.

STEP 2: Train Reward Model (RM)
  - Show human annotators multiple completions for same prompt
  - Annotators RANK them (A > B > C)
  - Train a reward model R(prompt, response) → scalar "quality score"
  - Much cheaper than writing demonstrations: just ranking!
  - ~33,000 comparisons used for InstructGPT
  
  The reward model learns what humans prefer:
    - Helpful, clear, accurate responses → high reward
    - Harmful, irrelevant, verbose responses → low reward

STEP 3: RL Fine-tuning with PPO
  - Use the SFT model as starting point
  - Generate responses for prompts using current LLM
  - Score with reward model
  - Update LLM with PPO to maximize reward
  - CRUCIAL: add KL penalty to keep LLM close to SFT baseline
  
  Total objective:
    objective(φ) = E[R(x,y)] - β · KL(π_φ(y|x) || π_SFT(y|x))
  
  The KL term prevents "reward hacking":
  Without it, the model might find adversarial strings that get high
  reward from the RM without being actually useful.
```

### Why InstructGPT Beats Larger Raw LMs

```
INSTRUCTGPT PAPER RESULTS:

Human preference studies:
  InstructGPT-1.3B preferred over GPT-3-175B: 88% of the time
  
  A 1.3 BILLION parameter model preferred over a 175 BILLION parameter model!
  
  This shows: alignment matters more than raw scale for user preferences.
  
  GPT-3-175B: great at completing text
  InstructGPT-1.3B: actually tries to help

OTHER IMPROVEMENTS:
  Truthfulness: InstructGPT more honest (fewer hallucinations)
  Toxicity: less likely to produce harmful content
  Instruction following: follows multi-part instructions better
  
  Tradeoff: InstructGPT slightly worse on NLP benchmarks
  (RLHF optimizes for human preference, not accuracy on benchmarks)
```

### ChatGPT

```
CHATGPT (November 2022):
  Technical details: not fully disclosed
  Base model: GPT-3.5 (a fine-tuned version of GPT-3)
  Training: RLHF similar to InstructGPT
  Key addition: MULTI-TURN CONVERSATION formatting
  
  Chat format (special tokens):
    <|system|>: You are a helpful assistant.
    <|user|>: What is the capital of France?
    <|assistant|>: Paris is the capital of France.
    <|user|>: And what's the population?
    <|assistant|>: Paris has a population of approximately 2.1 million...
  
  The model sees the full conversation history as its context.
  This enables multi-turn coherent dialogue.
  
  Public reception: 1M users in 5 days, 100M users in 2 months.
  Fastest growing consumer app in history (at the time).
  
  What it showed the world:
    - LLMs can be genuinely useful, not just technically impressive
    - The UX (conversational interface) matters enormously
    - RLHF is the missing piece for practical assistants
```

---

## 8. GPT-4 and Beyond

```
GPT-4 (March 2023):
  OpenAI stopped publishing technical details — no parameter count released.
  
  Confirmed:
    - Multimodal: accepts text AND images as input
    - System prompts: configure model behavior per application
    - Context: 8K tokens (later 32K, eventually 128K)
    - Improved reasoning, coding, instruction following
    - Better calibration (knows what it doesn't know)
  
  Rumored architecture (leaked/speculated):
    - Sparse Mixture of Experts (MoE)
    - 8 expert networks, each ~220B params
    - Total: ~1.8T params, but only ~220B active per forward pass
    - This would explain high quality with manageable inference cost
  
  Notable achievements:
    - Passes bar exam in top 10th percentile (GPT-3.5: bottom 10th)
    - AMC10/12 competition math: 40th percentile
    - Introductory ML course (MIT 6.036): 86% (A)
    - Medical licensing exam: passes all three parts

GPT-4 VISION:
  Input: text + image
  "What is in this image?" + [image]
  Can describe images, read text in images, understand charts

GPT-4 TURBO (November 2023):
  128K context window
  Knowledge cutoff extended to April 2023
  Cheaper API pricing
```

---

## 9. Emergent Abilities

One of the most fascinating and debated phenomena in large language model research:

```
EMERGENT ABILITIES:
  Capabilities that appear suddenly at a certain scale threshold,
  and are not predictable by extrapolating from smaller models.

EXAMPLES OF EMERGENT ABILITIES:
  
  Arithmetic:
    Model size < 50B: random on 2-digit arithmetic
    Model size > 100B: ~80%+ accuracy on 2-digit arithmetic
    (Not gradual improvement — sudden jump)
  
  Multi-step reasoning:
    < 50B: fails completely at multi-hop questions requiring 3+ steps
    > 100B: suddenly solves many such questions
  
  Code generation:
    < 5B: gibberish
    > 50B: syntactically valid, sometimes runnable code
    > 100B: often functionally correct code
  
  Analogical reasoning:
    < 50B: random performance on A:B::C:D analogies
    > 100B: ~80% accuracy
  
  Chain-of-thought:
    < 100B: CoT prompting makes performance WORSE
    > 100B: CoT prompting dramatically improves performance

WHY DO EMERGENT ABILITIES APPEAR?

HYPOTHESIS 1: "True emergence"
  The model must learn certain fundamental capabilities first.
  These serve as prerequisites for harder tasks.
  Once the prerequisites are learned, harder tasks unlock.
  Analogous to phase transitions in physics.

HYPOTHESIS 2: "Measurement artifact" (Schaeffer et al. 2023)
  Metrics that are non-linear (like exact-match accuracy) appear to
  show sharp emergence even when the underlying capability improves smoothly.
  With smoother metrics (like partial credit), emergence disappears.
  
  Example: If probability of each step correct = p,
    n-step problem accuracy = p^n
    This creates a dramatic "cliff" even as p improves linearly.

THE DEBATE CONTINUES:
  Some researchers: true emergence is real — capabilities unlock
  Others: it's all about metric choice
  
  Practical consequence either way: don't extrapolate small model 
  performance to large models. Larger models may surprise you.
```

---

## 10. Open-Source Alternatives

The LLM landscape democratized significantly after 2023:

### LLaMA (Meta AI, 2023)

```
LLaMA (Large Language Model Meta AI):
  Released: February 2023 (weights released for research)
  
  Sizes: 7B, 13B, 33B, 65B parameters
  Training data: 1-1.4T tokens
  Data: CommonCrawl, C4, GitHub, Wikipedia, Books, ArXiv, StackExchange
  
  LLaMA-2 (July 2023):
    Sizes: 7B, 13B, 70B
    Commercial license (can use in products!)
    RLHF fine-tuned chat models: LLaMA-2-chat
    Context: 4096 tokens (doubled from LLaMA-1)
    GQA (Grouped Query Attention) in 34B+ models
  
  LLaMA-3 (April 2024):
    Sizes: 8B, 70B
    Training data: 15 TRILLION tokens (10× LLaMA-2)
    Context: 8,192 tokens (base), 128k with RoPE extension
    Better tokenizer (128k vocabulary vs 32k)
    
    LLaMA-3 8B quality ≈ GPT-3.5 Turbo
    
ARCHITECTURE INNOVATIONS (vs original Transformer):
  - RoPE positional embeddings (rotary, not absolute)
  - RMSNorm instead of LayerNorm
  - SwiGLU activation in FFN
  - No bias in linear layers
  - Grouped Query Attention (GQA) for large models
  (All covered in detail in Chapter 30)
```

### Mistral 7B (Mistral AI, 2023)

```
MISTRAL 7B:
  Released: September 2023, fully open weights
  
  Architecture innovations:
    - Sliding Window Attention (SWA):
        Each token attends to at most W previous tokens (W=4096)
        Instead of full context attention
        Reduces memory from O(n²) to O(n·W)
        Handles long sequences more efficiently
    
    - Grouped Query Attention (GQA):
        32 query heads, 8 key/value heads (4:1 ratio)
        K/V cache is 4× smaller
        Faster generation
    
  Performance:
    Mistral-7B outperforms LLaMA-2-13B on most benchmarks
    2× fewer parameters, same or better quality
    Shows: architecture + data quality matters as much as size
  
  Fine-tuned variants (community):
    Mistral-7B-Instruct: basic instruction following
    Zephyr-7B: DPO fine-tuned, very good instruction following
    OpenHermes-2.5: aggressive fine-tuning, widely used
```

### Other Notable Open Models

```
MODEL LANDSCAPE TABLE (2024):
═══════════════════════════════════════════════════════════════════════
Model          | Org         | Params | Key Feature
───────────────────────────────────────────────────────────────────────
LLaMA 3.1     | Meta        | 8B-405B| Massive training data, open
Mistral 7B    | Mistral AI  | 7B     | Efficient attention, open
Mixtral 8×7B  | Mistral AI  | 46.7B* | Sparse MoE, 8 experts
Gemma 2B/7B   | Google      | 2-7B   | Distillation from large models
Phi-3         | Microsoft   | 3.8B   | Textbook-quality data
Qwen2         | Alibaba     | 0.5-72B| Strong multilingual
Command-R+    | Cohere      | 104B   | RAG-optimized
Falcon        | TII UAE     | 7-180B | Fully open, Apache 2.0
Deepseek-V3   | DeepSeek    | 685B*  | MoE, strong math/code
───────────────────────────────────────────────────────────────────────
* = total params; much fewer active per token (MoE architecture)
```

---

## 11. Context Window Evolution

```
CONTEXT WINDOW GROWTH:
═══════════════════════════════════════════════════════════════════
GPT-2 (2019):         1,024 tokens
GPT-3 (2020):         2,048 tokens
GPT-3.5 Turbo (2022): 4,096 tokens
GPT-4 (2023):         8,192 tokens
GPT-4-32k (2023):     32,768 tokens
Claude 1 (2023):       8,000 tokens
Claude 2 (2023):      100,000 tokens
GPT-4 Turbo (2023):   128,000 tokens
Gemini 1.5 Pro (2024): 1,000,000 tokens (1M!)
═══════════════════════════════════════════════════════════════════

TOKEN TO WORD APPROXIMATE CONVERSION:
  ~1.3 tokens per English word (GPT tokenizer)
  2,048 tokens ≈ 1,500 words ≈ 3 pages
  8,192 tokens ≈ 6,000 words ≈ 12 pages
  128k tokens ≈ 96,000 words ≈ 192 pages (a novel)
  1M tokens ≈ 750,000 words (Harry Potter + Fellowship of the Ring + more)
```

### Why Extending Context Is Hard

```
PROBLEM 1: POSITIONAL EMBEDDING EXTRAPOLATION
  Absolute learned PE (GPT-2): position 0..1023 have learned embeddings.
  If you try to input position 2000: model has NEVER seen that embedding.
  Performance degrades catastrophically.

PROBLEM 2: ATTENTION COMPUTATIONAL COST
  Standard attention: O(n²) per layer
  8k tokens: manageable
  128k tokens: 128² = 16M attention weights per head, per layer → expensive

SOLUTION 1: ROPE SCALING
  RoPE encodes position as rotation angle.
  Scale the angle: θ_scaled = θ / scale_factor
  
  "Position Interpolation" (Chen et al. 2023):
    Compress positions to fit original range.
    Trained on 2048 → extend to 8192 by dividing position by 4.
    Works with minimal fine-tuning.

SOLUTION 2: YaRN (Yet Another RoPE extensioN)
  More sophisticated RoPE scaling:
    Different dimensions of RoPE have different frequencies.
    High-freq dimensions: scale aggressively (they change fast)
    Low-freq dimensions: scale less (they change slowly)
  Better quality for very long contexts.

SOLUTION 3: FLASH ATTENTION
  Tiled computation: O(n) memory instead of O(n²)
  Makes n=128k attention computationally feasible.

PRACTICAL LIMITS:
  Memory scales as O(n) even with Flash Attention (for KV cache)
  For n=128k, kv cache alone can be 10-50GB depending on model size
```

---

## 12. Sampling Strategies

How we choose the next token from the probability distribution matters enormously for output quality.

### Greedy Decoding

```
GREEDY DECODING:
  At each step: token = argmax(softmax(logits))
  Always pick the most probable next token.
  
  Pros: Fast, deterministic
  Cons: Repetitive, boring, gets stuck in loops
  
  Example (GPT-2, greedy):
    Prompt: "The universe is"
    Output: "The universe is a very, very, very, very, very, very..."
    ← Loops because "very" is always the most probable continuation
```

### Temperature Scaling

```
TEMPERATURE SAMPLING:
  Divide logits by temperature T before softmax:
  
  P(token_i) = softmax(logits / T)_i
  
  T = 1.0: original distribution (no change)
  T → 0:   approaches greedy (all mass on highest logit)
  T → ∞:   approaches uniform (all tokens equally likely)

INTUITION:
  High T: FLATTEN the distribution (more random, more diverse)
  Low T:  SHARPEN the distribution (more focused, less creative)

EXAMPLE — "The sky is":
  Logits: {"blue": 4.2, "red": 1.1, "green": 0.8, "banana": -2.0}
  
  T=1.0: softmax → {blue: 0.80, red: 0.15, green: 0.05, banana: 0.00}
         Sample: mostly "blue", occasionally "red"
  
  T=0.5: softmax(logits/0.5) → {blue: 0.97, red: 0.02, green: 0.01, banana: 0.00}
         Very likely to pick "blue" → safe, predictable
  
  T=2.0: softmax(logits/2.0) → {blue: 0.52, red: 0.20, green: 0.18, banana: 0.10}
         More uniform → more variety, more surprises, more errors

RECOMMENDED DEFAULTS:
  Factual Q&A:       T = 0.3 - 0.5  (confident, accurate)
  Creative writing:  T = 0.8 - 1.2  (diverse, interesting)
  Code generation:   T = 0.0 - 0.2  (deterministic, correct)
  Brainstorming:     T = 1.0 - 1.5  (very diverse)
```

### Top-k Sampling

```
TOP-K SAMPLING:
  Only sample from the TOP K most probable tokens.
  
  1. Compute softmax probabilities
  2. Keep only top-k tokens
  3. Re-normalize to sum to 1
  4. Sample from this restricted distribution
  
  Example (k=3):
    All probs: {blue: 0.45, red: 0.20, green: 0.15, yellow: 0.10, banana: 0.05, ...}
    Top-3:     {blue: 0.45, red: 0.20, green: 0.15} → renormalize → {blue: 0.56, red: 0.25, green: 0.19}
    Sample from top 3 only.
  
  Problem: k is FIXED, but the right "vocabulary size" varies!
  
  When model is confident:
    {blue: 0.98, red: 0.01, green: 0.005, ...}
    k=50 forces choosing from 50 tokens even though 1 is clearly right.
    Risk: sampling nonsense low-prob token.
  
  When model is uncertain:
    {blue: 0.03, red: 0.03, yellow: 0.03, ...}
    k=50 still restricts to 50 out of 50,000 tokens.
    May be too restrictive.
    
  k=50 is commonly used; larger k = more diverse but more risky.
```

### Top-p (Nucleus) Sampling

```
TOP-P (NUCLEUS) SAMPLING:
  Instead of fixed k tokens, use the SMALLEST SET of tokens
  whose cumulative probability ≥ p.
  
  Algorithm:
    1. Sort tokens by probability (descending)
    2. Compute cumulative sum
    3. Find smallest prefix with cumulative prob ≥ p
    4. Sample from that set (renormalized)
  
  Example (p=0.9):
  
  Token  | Prob  | Cumulative
  ─────────────────────────────
  blue   | 0.60  | 0.60
  red    | 0.20  | 0.80
  green  | 0.10  | 0.90 ← threshold reached!
  yellow | 0.07  | 0.97  (excluded)
  banana | 0.03  | 1.00  (excluded)
  
  Top-p nucleus: {blue, red, green} (3 tokens)
  Renormalize and sample.
  
  WHEN MODEL IS CONFIDENT:
    {blue: 0.95, ...}
    Cumulative 0.9 reached after 1 token.
    Nucleus = just "blue" → nearly deterministic! ✓
  
  WHEN MODEL IS UNCERTAIN:
    {blue: 0.02, red: 0.02, green: 0.02, ...}
    Cumulative 0.9 requires 45 tokens.
    Nucleus = 45 tokens → more variety! ✓
  
  p=0.9 or p=0.95: standard defaults
  Adapts to model confidence automatically.
  Generally preferred over top-k.
```

### Min-p Sampling

```
MIN-P SAMPLING (newer, often better):
  
  Set a minimum probability threshold = p_base × max_prob
  
  Example (p_base = 0.05):
    max_prob = 0.60 (for "blue")
    threshold = 0.05 × 0.60 = 0.03
    
    Keep all tokens with prob > threshold:
    blue (0.60 > 0.03) ✓
    red  (0.20 > 0.03) ✓
    green (0.10 > 0.03) ✓
    yellow (0.07 > 0.03) ✓
    banana (0.03 = threshold) marginal...
  
  ADVANTAGE: threshold scales with model confidence.
    Confident model (max=0.95): threshold = 0.0475 → few tokens
    Uncertain model (max=0.05): threshold = 0.0025 → many tokens
  
  Generally produces more coherent text than top-p for creative tasks.
  p_base = 0.05 to 0.10 are common values.
```

### Repetition Penalty

```
REPETITION PENALTY:
  Prevents the model from getting stuck repeating phrases.
  
  For each token in recent context, divide its logit by penalty:
  
  logit[token] /= repetition_penalty   if token appeared recently
  
  penalty = 1.0: no effect
  penalty = 1.2: moderate discouragement of recent tokens
  penalty = 1.5: strong discouragement
  
  Can also use "presence penalty" (additive) and "frequency penalty"
  (scales with count of how often token appeared).

TYPICAL SETTINGS:
  ┌─────────────────────────────────────────────────────┐
  │ Task              │ Temperature │ Top-p │ Rep. Pen. │
  ├─────────────────────────────────────────────────────┤
  │ Factual answers   │ 0.3         │ 0.9   │ 1.0       │
  │ Creative writing  │ 0.9         │ 0.95  │ 1.1       │
  │ Code completion   │ 0.1         │ 1.0   │ 1.0       │
  │ Chat/assistant    │ 0.7         │ 0.9   │ 1.1       │
  │ Brainstorming     │ 1.2         │ 0.95  │ 1.2       │
  └─────────────────────────────────────────────────────┘
```

---

## 13. HuggingFace: GPT-2 Generation

```python
"""
GPT-2 Text Generation with Multiple Sampling Strategies.
Complete, runnable code with detailed explanations.
"""

import torch
import torch.nn.functional as F
from transformers import GPT2LMHeadModel, GPT2Tokenizer
from typing import Optional
import time


def load_gpt2(model_name: str = "gpt2") -> tuple:
    """
    Load GPT-2 model and tokenizer.
    
    Available sizes:
      "gpt2"        — 117M params (fastest)
      "gpt2-medium" — 345M params
      "gpt2-large"  — 762M params
      "gpt2-xl"     — 1.5B params (slowest, best quality)
    """
    print(f"Loading {model_name}...")
    tokenizer = GPT2Tokenizer.from_pretrained(model_name)
    model = GPT2LMHeadModel.from_pretrained(model_name)
    
    device = torch.device("cuda" if torch.cuda.is_available() else "cpu")
    model = model.to(device)
    model.eval()
    
    # GPT-2 doesn't have a pad token by default
    tokenizer.pad_token = tokenizer.eos_token
    
    print(f"  Loaded on {device}")
    return model, tokenizer, device


def generate_with_sampling(
    model: GPT2LMHeadModel,
    tokenizer: GPT2Tokenizer,
    device: torch.device,
    prompt: str,
    max_new_tokens: int = 100,
    temperature: float = 1.0,
    top_k: int = 0,           # 0 = disabled
    top_p: float = 1.0,        # 1.0 = disabled
    repetition_penalty: float = 1.0,  # 1.0 = no penalty
    seed: int = 42,
) -> str:
    """
    Generate text with manual sampling control.
    Implements each sampling strategy from scratch for education.
    """
    torch.manual_seed(seed)
    
    # Encode prompt
    input_ids = tokenizer.encode(prompt, return_tensors="pt").to(device)
    generated = input_ids.clone()
    
    with torch.no_grad():
        for step in range(max_new_tokens):
            # Forward pass: get logits for next token
            outputs = model(generated)
            next_token_logits = outputs.logits[:, -1, :]  # (1, vocab_size)
            
            # ── Repetition Penalty ─────────────────────────────────────────
            if repetition_penalty != 1.0:
                for prev_token in generated[0].tolist():
                    if next_token_logits[0, prev_token] > 0:
                        next_token_logits[0, prev_token] /= repetition_penalty
                    else:
                        next_token_logits[0, prev_token] *= repetition_penalty
            
            # ── Temperature ────────────────────────────────────────────────
            if temperature != 1.0:
                next_token_logits = next_token_logits / temperature
            
            # ── Top-k Filtering ────────────────────────────────────────────
            if top_k > 0:
                # Get the top-k logits
                top_k_logits, _ = torch.topk(next_token_logits, top_k)
                # Minimum logit in top-k
                min_top_k = top_k_logits[:, -1].unsqueeze(-1)
                # Replace all logits below the threshold with -inf
                filter_mask = next_token_logits < min_top_k
                next_token_logits[filter_mask] = float('-inf')
            
            # ── Top-p Filtering ────────────────────────────────────────────
            if top_p < 1.0:
                # Sort logits descending
                sorted_logits, sorted_indices = torch.sort(
                    next_token_logits, dim=-1, descending=True
                )
                sorted_probs = F.softmax(sorted_logits, dim=-1)
                
                # Cumulative sum
                cumulative_probs = sorted_probs.cumsum(dim=-1)
                
                # Remove tokens where cumulative prob EXCEEDS top_p
                # (shift by 1 to include the token that crossed threshold)
                sorted_remove = cumulative_probs - sorted_probs > top_p
                sorted_logits[sorted_remove] = float('-inf')
                
                # Unsort
                next_token_logits.scatter_(1, sorted_indices, sorted_logits)
            
            # ── Sample or Greedy ──────────────────────────────────────────
            probs = F.softmax(next_token_logits, dim=-1)
            
            if temperature == 0.0:  # greedy
                next_token = torch.argmax(probs, dim=-1, keepdim=True)
            else:
                next_token = torch.multinomial(probs, num_samples=1)
            
            # Check for EOS token
            if next_token.item() == tokenizer.eos_token_id:
                break
            
            # Append new token
            generated = torch.cat([generated, next_token], dim=1)
    
    # Decode only the NEW tokens (not the prompt)
    new_tokens = generated[0, input_ids.shape[1]:]
    return tokenizer.decode(new_tokens, skip_special_tokens=True)


def compare_sampling_strategies(prompt: str):
    """Compare all sampling strategies on the same prompt."""
    model, tokenizer, device = load_gpt2("gpt2")
    
    strategies = [
        {
            "name": "Greedy (T=0)",
            "params": dict(temperature=0.01, top_k=0, top_p=1.0),
        },
        {
            "name": "Temperature T=0.5",
            "params": dict(temperature=0.5, top_k=0, top_p=1.0),
        },
        {
            "name": "Temperature T=1.0",
            "params": dict(temperature=1.0, top_k=0, top_p=1.0),
        },
        {
            "name": "Temperature T=1.5 (creative)",
            "params": dict(temperature=1.5, top_k=0, top_p=1.0),
        },
        {
            "name": "Top-k (k=50, T=0.8)",
            "params": dict(temperature=0.8, top_k=50, top_p=1.0),
        },
        {
            "name": "Top-p (p=0.9, T=0.8)",
            "params": dict(temperature=0.8, top_k=0, top_p=0.9),
        },
        {
            "name": "Top-p + Repetition Penalty",
            "params": dict(temperature=0.8, top_k=0, top_p=0.9, repetition_penalty=1.2),
        },
    ]
    
    print(f"\n{'=' * 70}")
    print(f"SAMPLING STRATEGIES COMPARISON")
    print(f"Prompt: '{prompt}'")
    print(f"{'=' * 70}")
    
    for strategy in strategies:
        name = strategy["name"]
        params = strategy["params"]
        
        start = time.time()
        text = generate_with_sampling(
            model, tokenizer, device, prompt,
            max_new_tokens=80,
            seed=42,
            **params
        )
        elapsed = time.time() - start
        
        print(f"\n[{name}] ({elapsed:.2f}s)")
        print(f"  {text[:200]}")  # truncate for display


def beam_search_generation(prompt: str, num_beams: int = 5):
    """
    Beam search generation using HuggingFace's built-in implementation.
    Shows the difference between beam search and sampling.
    """
    model, tokenizer, device = load_gpt2("gpt2")
    
    input_ids = tokenizer.encode(prompt, return_tensors="pt").to(device)
    
    with torch.no_grad():
        beam_outputs = model.generate(
            input_ids,
            max_new_tokens=80,
            num_beams=num_beams,
            early_stopping=True,
            no_repeat_ngram_size=2,  # prevent 2-gram repeats
            num_return_sequences=3,   # return top 3 beam sequences
        )
    
    print(f"\n{'=' * 70}")
    print(f"BEAM SEARCH (num_beams={num_beams})")
    print(f"Prompt: '{prompt}'")
    print(f"{'=' * 70}")
    
    for i, beam_output in enumerate(beam_outputs):
        new_tokens = beam_output[input_ids.shape[1]:]
        text = tokenizer.decode(new_tokens, skip_special_tokens=True)
        print(f"\nBeam {i+1}: {text}")


def compute_perplexity(model, tokenizer, device, text: str) -> float:
    """
    Compute perplexity of a text under GPT-2.
    Lower = model "expected" this text more (it fits the model's distribution).
    """
    input_ids = tokenizer.encode(text, return_tensors="pt").to(device)
    
    with torch.no_grad():
        outputs = model(input_ids, labels=input_ids)
        loss = outputs.loss  # cross-entropy loss per token
    
    return torch.exp(loss).item()  # perplexity = e^loss


def inspect_token_probabilities(
    model: GPT2LMHeadModel,
    tokenizer: GPT2Tokenizer,
    device: torch.device,
    context: str,
    top_n: int = 10,
):
    """
    Show the probability distribution over next tokens given context.
    Great for understanding what the model "thinks" comes next.
    """
    input_ids = tokenizer.encode(context, return_tensors="pt").to(device)
    
    with torch.no_grad():
        outputs = model(input_ids)
        logits = outputs.logits[0, -1, :]  # logits for next token
    
    probs = F.softmax(logits, dim=-1)
    top_probs, top_indices = torch.topk(probs, top_n)
    
    print(f"\nNext token probabilities after: '{context}'")
    print(f"{'Rank':>5} {'Token':>15} {'Prob':>8} {'Bar'}")
    
    for rank, (prob, idx) in enumerate(zip(top_probs, top_indices)):
        token = tokenizer.decode([idx.item()])
        bar = "█" * int(prob.item() * 40)
        print(f"  {rank+1:>3}. {repr(token):>15} {prob.item():>7.4f}  {bar}")


# ── Main demo ─────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    # Load model once
    model, tokenizer, device = load_gpt2("gpt2")  # 117M params, fast
    
    prompt = "The invention of the transformer architecture in 2017 was"
    
    # Show token probabilities
    print("=" * 60)
    print("NEXT TOKEN PROBABILITY INSPECTION")
    inspect_token_probabilities(model, tokenizer, device, prompt)
    
    # Compare sampling strategies
    compare_sampling_strategies(prompt)
    
    # Beam search
    beam_search_generation(prompt, num_beams=4)
    
    # Perplexity examples
    print("\n" + "=" * 60)
    print("PERPLEXITY COMPARISON")
    texts = [
        "The cat sat on the mat.",           # natural English
        "The cat quantum banana exploded.",   # unnatural
        "The transformer paper was published in 2017.",  # factual
        "kfjhskjfhskjfhskjfh",              # noise
    ]
    
    for text in texts:
        ppl = compute_perplexity(model, tokenizer, device, text)
        print(f"  PPL={ppl:8.1f}: {text}")
```

---

## 14. Summary

```
GPT PROGRESSION:
═══════════════════════════════════════════════════════════════════
GPT-1 (2018, 117M):
  Proof of concept — pre-train → fine-tune works
  
GPT-2 (2019, 1.5B):
  Zero-shot generalization emerges
  Release controversy shapes AI safety discourse
  
GPT-3 (2020, 175B):
  In-context learning: examples in prompt teach task without gradient update
  Chain-of-thought reasoning begins to emerge
  Scale alone produces powerful capabilities
  
InstructGPT (2022, 1.3B via RLHF):
  RLHF aligns model to human preferences
  Smaller aligned model beats larger raw model
  
ChatGPT (2022):
  Conversational interface + RLHF = mainstream breakthrough
  100M users in 2 months
  
GPT-4 (2023):
  Multimodal, 128k context, near-human expert performance
  Architecture not disclosed
  
Open-source alternatives catch up:
  LLaMA, Mistral, Phi, Gemma — competitive with early GPT versions
═══════════════════════════════════════════════════════════════════
```

### Sampling Strategy Quick Reference

| Strategy | When to Use | Key Parameter | Quality |
|---------|-------------|---------------|---------|
| Greedy | Factual, code | - | Deterministic, repetitive |
| Temperature | General | T=0.7-1.0 | Good baseline |
| Top-k | General | k=50 | Good, but fixed vocabulary |
| Top-p (nucleus) | Creative, chat | p=0.9 | Best overall |
| Min-p | Creative | p_base=0.05 | Better than top-p for creative |
| Beam search | Translation | n_beams=4-8 | Best for fixed-length outputs |

---

## Mini Projects

### Mini Project 1: Miniature GPT from Scratch

Build a decoder-only GPT-style transformer and train it to generate coherent text.

**Objective:** Understand causal language modeling by implementing the masked self-attention and training loop.

```python
import torch
import torch.nn as nn
import torch.nn.functional as F
import numpy as np
import matplotlib.pyplot as plt
import math

# ─── Training text ────────────────────────────────────────────────────────────
text = (
    "in the beginning god created the heaven and the earth and the earth was without form "
    "and void and darkness was upon the face of the deep and the spirit of god moved upon "
    "the face of the waters and god said let there be light and there was light and god saw "
    "the light that it was good and god divided the light from the darkness and god called "
    "the light day and the darkness he called night and the evening and the morning were the "
    "first day and god said let there be a firmament in the midst of the waters and let it "
    "divide the waters from the waters and god made the firmament and divided the waters "
    "which were under the firmament from the waters which were above the firmament and it "
    "was so and god called the firmament heaven and the evening and the morning were the second day"
)

# Character-level tokenization
chars  = sorted(set(text))
vocab  = len(chars)
c2i    = {c: i for i, c in enumerate(chars)}
i2c    = {i: c for c, i in c2i.items()}
data   = torch.LongTensor([c2i[c] for c in text])

BLOCK   = 32   # context length
BATCH   = 16
D_MODEL = 64
N_HEADS = 4
N_LAYER = 3
FF_DIM  = 128

def get_batch():
    idx = torch.randint(0, len(data) - BLOCK, (BATCH,))
    X   = torch.stack([data[i:i+BLOCK]   for i in idx])
    y   = torch.stack([data[i+1:i+BLOCK+1] for i in idx])
    return X, y

# ─── Causal Self-Attention ─────────────────────────────────────────────────
class CausalSelfAttention(nn.Module):
    def __init__(self, d_model, n_heads, block_size):
        super().__init__()
        self.n_heads = n_heads
        self.d_k     = d_model // n_heads
        self.qkv  = nn.Linear(d_model, 3 * d_model)
        self.proj = nn.Linear(d_model, d_model)
        self.drop = nn.Dropout(0.1)
        # Causal mask: lower triangular
        mask = torch.tril(torch.ones(block_size, block_size))
        self.register_buffer('mask', mask)

    def forward(self, x):
        B, T, D = x.shape
        qkv = self.qkv(x).split(D, dim=2)
        def reshape(t): return t.view(B, T, self.n_heads, self.d_k).transpose(1, 2)
        Q, K, V = map(reshape, qkv)
        scores = Q @ K.transpose(-2, -1) / math.sqrt(self.d_k)
        scores = scores.masked_fill(self.mask[:T, :T] == 0, -1e9)
        attn   = self.drop(F.softmax(scores, dim=-1))
        out    = (attn @ V).transpose(1, 2).contiguous().view(B, T, D)
        return self.drop(self.proj(out))

class GPTBlock(nn.Module):
    def __init__(self, d_model, n_heads, ff_dim, block_size):
        super().__init__()
        self.ln1  = nn.LayerNorm(d_model)
        self.attn = CausalSelfAttention(d_model, n_heads, block_size)
        self.ln2  = nn.LayerNorm(d_model)
        self.ff   = nn.Sequential(
            nn.Linear(d_model, ff_dim), nn.GELU(), nn.Linear(ff_dim, d_model), nn.Dropout(0.1)
        )

    def forward(self, x):
        x = x + self.attn(self.ln1(x))
        x = x + self.ff(self.ln2(x))
        return x

class MiniGPT(nn.Module):
    def __init__(self, vocab, d_model, n_heads, ff_dim, n_layer, block_size):
        super().__init__()
        self.tok_emb = nn.Embedding(vocab, d_model)
        self.pos_emb = nn.Embedding(block_size, d_model)
        self.blocks  = nn.Sequential(*[GPTBlock(d_model, n_heads, ff_dim, block_size)
                                        for _ in range(n_layer)])
        self.ln      = nn.LayerNorm(d_model)
        self.head    = nn.Linear(d_model, vocab, bias=False)
        # Weight tying (same matrix for embed and unembed — common in LLMs)
        self.head.weight = self.tok_emb.weight

    def forward(self, x):
        B, T = x.shape
        pos  = torch.arange(T, device=x.device)
        h    = self.tok_emb(x) + self.pos_emb(pos)
        h    = self.ln(self.blocks(h))
        return self.head(h)

    @torch.no_grad()
    def generate(self, prompt_ids, n_new=100, temperature=0.8):
        self.eval()
        ctx = prompt_ids.clone()
        for _ in range(n_new):
            ctx_crop = ctx[-BLOCK:]
            logits   = self(ctx_crop.unsqueeze(0))[0, -1]
            probs    = F.softmax(logits / temperature, dim=0)
            next_tok = torch.multinomial(probs, 1)
            ctx      = torch.cat([ctx, next_tok])
        return ''.join(i2c[i.item()] for i in ctx)

torch.manual_seed(42)
model = MiniGPT(vocab, D_MODEL, N_HEADS, FF_DIM, N_LAYER, BLOCK)
n_params = sum(p.numel() for p in model.parameters())
print(f"MiniGPT: {n_params:,} parameters ({n_params/1e6:.3f}M)")

opt  = torch.optim.AdamW(model.parameters(), lr=3e-3, weight_decay=0.01)
crit = nn.CrossEntropyLoss()

losses, perplexities = [], []
n_steps = 1000
for step in range(n_steps):
    X_b, y_b = get_batch()
    logits = model(X_b)
    loss   = crit(logits.view(-1, vocab), y_b.view(-1))
    opt.zero_grad(); loss.backward(); torch.nn.utils.clip_grad_norm_(model.parameters(), 1.0); opt.step()
    losses.append(loss.item())
    if (step+1) % 200 == 0:
        ppl = math.exp(loss.item())
        perplexities.append(ppl)
        prompt = torch.LongTensor([c2i[c] for c in "and god said "])
        gen = model.generate(prompt, n_new=60, temperature=0.7)
        print(f"Step {step+1:4d}: loss={loss.item():.4f}, PPL={ppl:.1f}")
        print(f"  Generated: '{gen}'")

fig, axes = plt.subplots(1, 2, figsize=(13, 4))
fig.suptitle("MiniGPT Training from Scratch", fontsize=12, fontweight='bold')

# Smooth loss
window = 30
smoothed = np.convolve(losses, np.ones(window)/window, mode='valid')
axes[0].plot(losses, alpha=0.3, color='purple', linewidth=0.5)
axes[0].plot(range(window-1, len(losses)), smoothed, 'purple', linewidth=2)
axes[0].set_title("Training Loss (charLM)"); axes[0].set_xlabel("Step"); axes[0].grid(True, alpha=0.3)

# Temperature comparison
temps = [0.3, 0.7, 1.0, 1.5]
axes[1].axis('off')
prompt = torch.LongTensor([c2i.get(c, 0) for c in "and the "])
text_block = "Temperature Sampling Comparison:\n\n"
for t in temps:
    gen = model.generate(prompt, n_new=80, temperature=t)
    text_block += f"T={t:.1f}: {gen[8:60]}...\n\n"
axes[1].text(0.02, 0.95, text_block, transform=axes[1].transAxes, fontsize=8,
             va='top', fontfamily='monospace',
             bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))
axes[1].set_title("Text Generation at Different Temperatures")

plt.tight_layout()
plt.savefig("minigpt.png", dpi=150)
plt.show()
```

---

### Mini Project 2: Decoding Strategy Comparator

Compare greedy, top-k, top-p, and beam search decoding on the same prompt.

**Objective:** Understand that *how* you decode is as important as the model itself.

```python
import torch
import torch.nn.functional as F
import numpy as np

# Re-use MiniGPT from above (run the previous cell first)
# Or load a pretrained model:
# from transformers import GPT2LMHeadModel, GPT2Tokenizer
# model = GPT2LMHeadModel.from_pretrained('gpt2')
# tokenizer = GPT2Tokenizer.from_pretrained('gpt2')

def get_logits_from_context(model, context_ids, vocab_size):
    """Get next-token logits from context."""
    ctx = context_ids[-BLOCK:].unsqueeze(0)
    with torch.no_grad():
        logits = model(ctx)[0, -1]
    return logits

def greedy_decode(model, prompt, n_new=80):
    ctx = prompt.clone()
    for _ in range(n_new):
        logits = get_logits_from_context(model, ctx, vocab)
        next_tok = logits.argmax()
        ctx = torch.cat([ctx, next_tok.unsqueeze(0)])
    return ''.join(i2c[i.item()] for i in ctx)

def topk_decode(model, prompt, n_new=80, k=10, temperature=0.8):
    ctx = prompt.clone()
    for _ in range(n_new):
        logits = get_logits_from_context(model, ctx, vocab) / temperature
        top_vals, top_ids = logits.topk(k)
        probs = F.softmax(top_vals, dim=0)
        next_tok = top_ids[torch.multinomial(probs, 1)]
        ctx = torch.cat([ctx, next_tok])
    return ''.join(i2c[i.item()] for i in ctx)

def topp_decode(model, prompt, n_new=80, p=0.9, temperature=0.8):
    ctx = prompt.clone()
    for _ in range(n_new):
        logits = get_logits_from_context(model, ctx, vocab) / temperature
        probs_sorted, sorted_ids = logits.softmax(0).sort(descending=True)
        cumsum = probs_sorted.cumsum(0)
        mask = (cumsum - probs_sorted) < p
        filtered_probs = probs_sorted[mask]
        filtered_ids   = sorted_ids[mask]
        filtered_probs /= filtered_probs.sum()
        next_tok = filtered_ids[torch.multinomial(filtered_probs, 1)]
        ctx = torch.cat([ctx, next_tok])
    return ''.join(i2c[i.item()] for i in ctx)

def beam_search(model, prompt, n_new=80, beam_width=4):
    beams = [(prompt.clone(), 0.0)]  # (tokens, log_prob)
    for _ in range(n_new):
        all_candidates = []
        for seq, score in beams:
            logits = get_logits_from_context(model, seq, vocab)
            log_probs = F.log_softmax(logits, dim=0)
            top_vals, top_ids = log_probs.topk(beam_width)
            for val, idx in zip(top_vals, top_ids):
                candidate = (torch.cat([seq, idx.unsqueeze(0)]), score + val.item())
                all_candidates.append(candidate)
        # Keep top beam_width candidates
        all_candidates.sort(key=lambda x: x[1], reverse=True)
        beams = all_candidates[:beam_width]
    best_seq = beams[0][0]
    return ''.join(i2c[i.item()] for i in best_seq)

# Compare all decoding strategies
prompt_str = "and god said "
prompt_ids = torch.LongTensor([c2i.get(c, 0) for c in prompt_str])

print("Decoding Strategy Comparison:")
print(f"Prompt: '{prompt_str}'\n")

strategies = [
    ("Greedy",          lambda: greedy_decode(model, prompt_ids, 80)),
    ("Top-k (k=10)",   lambda: topk_decode(model, prompt_ids, 80, k=10, temperature=0.8)),
    ("Top-p (p=0.9)",  lambda: topp_decode(model, prompt_ids, 80, p=0.9, temperature=0.8)),
    ("Beam Search (4)", lambda: beam_search(model, prompt_ids, 80, beam_width=4)),
    ("Top-k (k=3)",    lambda: topk_decode(model, prompt_ids, 80, k=3, temperature=0.5)),
    ("High Temp (1.5)", lambda: topk_decode(model, prompt_ids, 80, k=30, temperature=1.5)),
]

results = []
for name, fn in strategies:
    output = fn()
    results.append((name, output[len(prompt_str):]))
    print(f"[{name}]")
    print(f"  {output[len(prompt_str):80]}")
    print()

# Visualize token diversity: count unique n-grams
def unique_bigrams(text):
    return len(set(zip(text[:-1], text[1:])))

diversity = [unique_bigrams(r[1]) for r in results]
names_short = [r[0] for r in results]

fig, axes = plt.subplots(1, 2, figsize=(13, 5))
fig.suptitle("Decoding Strategy Comparison", fontsize=12, fontweight='bold')
colors = ['gray', 'steelblue', 'orange', 'green', 'purple', 'red']
axes[0].bar(names_short, diversity, color=colors, alpha=0.8)
axes[0].set_title("Output Diversity\n(unique character bigrams — higher = more varied)")
axes[0].set_ylabel("Unique Bigrams"); axes[0].set_xticklabels(names_short, rotation=20, ha='right')
axes[0].grid(True, alpha=0.3, axis='y')

# Text display
axes[1].axis('off')
text_display = "Generated Continuations:\n\n"
for name, output in results:
    text_display += f"[{name}]\n{output[:50]}...\n\n"
axes[1].text(0.02, 0.97, text_display, transform=axes[1].transAxes,
             fontsize=7, va='top', fontfamily='monospace',
             bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))

plt.tight_layout()
plt.savefig("decoding_strategies.png", dpi=150)
plt.show()
```

---

## 15. Exercises

**Exercise 1**: Load GPT-2 and compute perplexity on 20 sentences from a textbook, 20 random sentences, and 20 sentences of Lorem Ipsum text. Plot the distributions. Which corpus has the lowest perplexity?

**Exercise 2**: Implement min-p sampling from scratch (based on the description in Section 12). Compare it to top-p with p=0.9 on 10 creative writing prompts. Which produces more coherent text?

**Exercise 3**: Demonstrate "prompt sensitivity" in GPT-2. Generate the same content using 5 differently-phrased prompts ("Tell me about", "Write about", "Describe", "Explain", "What is"). How much does the phrasing affect quality? This motivates "prompt engineering".

**Exercise 4**: Implement and analyze in-context learning with GPT-2. Create a sentiment classification task. Give 0, 1, 3, and 5 examples in the prompt. At how many examples does GPT-2 start to follow the format? Compare with GPT-3 (if you have API access).

**Exercise 5**: Implement beam search from scratch (without HuggingFace's generate()). Verify your outputs match HuggingFace's for beam_size=2.

**Exercise 6**: Study the Chinchilla scaling law. Given a compute budget of 10²⁰ FLOPs, what is the optimal model size and training tokens according to Chinchilla? How does the prediction compare to what LLaMA-2 actually used?

---

**Chapter Summary**: GPT models are autoregressive language models — decoder-only transformers trained to predict the next token. GPT-1 proved pre-training + fine-tuning works. GPT-2 showed zero-shot generalization and sparked release controversy. GPT-3 at 175B demonstrated in-context learning and emergent abilities like chain-of-thought. Scaling laws reveal predictable power-law improvements with model size and data — Chinchilla (2022) showed previous models were undertrained, optimal is ~20 tokens per parameter. InstructGPT and ChatGPT added RLHF — training the model to follow instructions via human preference feedback — turning a text completer into a helpful assistant. Open-source alternatives (LLaMA, Mistral) provide competitive models. Sampling strategies (temperature, top-k, top-p) control the creativity/accuracy trade-off during generation.

**What's Next →** [Chapter 26: Training Large Language Models](./26-training-large-language-models.md)

*How do you actually train a 175B parameter model? We'll cover the engineering challenges: distributed training, data pipelines, RLHF in detail, Constitutional AI, and production serving.*
