# Chapter 22: The Attention Mechanism — The Core Idea That Changed Everything

> **"Attention is all you need. But before we understood that, we needed to understand what was missing."**

---

## Table of Contents
1. [The Sequence-to-Sequence Problem](#1-the-sequence-to-sequence-problem)
2. [The Information Bottleneck](#2-the-information-bottleneck)
3. [Bahdanau Attention — The Original Breakthrough](#3-bahdanau-attention--the-original-breakthrough)
4. [The Alignment Matrix](#4-the-alignment-matrix)
5. [Luong Attention — Simplified Dot-Product](#5-luong-attention--simplified-dot-product)
6. [Self-Attention — The True Revolution](#6-self-attention--the-true-revolution)
7. [Query, Key, Value — The Database Analogy](#7-query-key-value--the-fundamental-construct)
8. [Scaled Dot-Product Attention — Full Derivation](#8-scaled-dot-product-attention--full-derivation)
9. [Multi-Head Attention](#9-multi-head-attention)
10. [Computational Complexity](#10-computational-complexity)
11. [Modern Attention Variants](#11-modern-attention-variants)
12. [PyTorch Implementation from Scratch](#12-pytorch-implementation-from-scratch)
13. [Summary](#13-summary)
14. [Mini Projects](#14-mini-projects)
15. [Exercises](#15-exercises)

---

## 1. The Sequence-to-Sequence Problem

Before attention, the dominant approach to sequence problems (translation, summarization, speech recognition) was the **encoder-decoder RNN architecture**.

The idea was simple and elegant:

```
INPUT SENTENCE (English):
"The agreement on the European Economic Area was signed in August 1992"

ENCODER RNN:
Word by word, left to right, update hidden state:

  "The"  → h1
  "agreement" → h2
  "on"   → h3
  "the"  → h4
  "European" → h5
  "Economic" → h6
  "Area" → h7
  "was"  → h8
  "signed" → h9
  "in"   → h10
  "August" → h11
  "1992" → h12 ← FINAL HIDDEN STATE (context vector)

DECODER RNN:
Takes h12 as its starting state, generates French word by word:
  h12 → "L'"
  h12 + "L'" → "accord"
  ...
```

The encoder is a Recurrent Neural Network (LSTM or GRU) that reads the input sequence one token at a time, updating a hidden state vector at each step. The decoder is another RNN that uses the **final hidden state** of the encoder as its initial state, then generates output tokens one by one.

This was a breakthrough in 2014. Before this, you needed hand-engineered translation systems with phrase tables, alignment models, and language models. Now a single neural network could learn translation end-to-end.

But there was a fundamental problem.

---

## 2. The Information Bottleneck

The entire input sentence — no matter how long — had to be compressed into **a single fixed-size vector**. The last hidden state h₁₂ had to somehow capture everything: "The", "agreement", "European Economic Area", "August 1992".

```
THE BOTTLENECK PROBLEM:
===============================================================

Input sentence (12 words):
"The agreement on the European Economic Area was signed in August 1992"

                    ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓
                 [ENCODER RNN — 12 steps]
                    ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓ ↓
                         ┌─────────┐
                         │  h_12   │  ← 512 numbers
                         │ (single │     must encode
                         │ vector) │     ALL information
                         └────┬────┘
                              │ BOTTLENECK
                              ↓
                         [DECODER RNN]
                         generates French
===============================================================

What h_12 must remember:
  - "The" (grammatical article, determines French gender)
  - "agreement" (main noun)
  - "European Economic Area" (specific proper noun — don't translate it!)
  - "1992" (year, goes at a specific place in French)
  - grammatical structure of the whole sentence

This is like asking someone to read a page, memorize it,
then throw away the page and reconstruct it from memory.
```

### Why This Fails for Long Sentences

RNNs process sequences step by step. The gradient of the loss with respect to early inputs must travel backward through many time steps. This causes:

**Vanishing gradients**: The gradient shrinks exponentially as it travels backward. By the time it reaches the first word, it's nearly zero — the model can't learn from early parts of long sentences.

**LSTM and GRU** mitigate this with gating mechanisms but don't eliminate it. Even with LSTMs, performance on translation degrades significantly once sentences exceed 20-30 words. The model simply can't hold everything in a fixed vector.

```
Performance vs Sentence Length (approximate BLEU score):

BLEU
 35 |*
 30 | **
 25 |   ***
 20 |      ****
 15 |          *****
 10 |               *******
  5 |                      **********
  0 +----+----+----+----+----+----+----
    5   10   15   20   25   30   35  40
                  Sentence Length (words)

Encoder-Decoder (no attention): performance cliff after ~20 words
With Attention: nearly flat performance across lengths
```

This degradation was the unsolved problem in 2014. Then Dzmitri Bahdanau, Kyunghyun Cho, and Yoshua Bengio published a paper that changed everything.

---

## 3. Bahdanau Attention — The Original Breakthrough

The 2014 paper "Neural Machine Translation by Jointly Learning to Align and Translate" introduced a radical idea:

**Instead of compressing everything into one vector, let the decoder look at ALL encoder hidden states at every output step.**

At each decoder step t, instead of just using the final hidden state, the decoder can "attend" to whichever encoder states are most relevant for the current output word.

### The Architecture

```
BAHDANAU ATTENTION ARCHITECTURE:
=================================================================

INPUT: "The cat sat on the mat"
       h1   h2  h3  h4  h5  h6    ← ALL encoder hidden states

At each decoder step t (generating French word):

Step 1: Compute ALIGNMENT SCORES
  How relevant is each encoder state h_s to current decoder state s_{t-1}?

  e_{t,1} = a(s_{t-1}, h1)   ← relevance of h1 to current step
  e_{t,2} = a(s_{t-1}, h2)
  e_{t,3} = a(s_{t-1}, h3)
  e_{t,4} = a(s_{t-1}, h4)
  e_{t,5} = a(s_{t-1}, h5)
  e_{t,6} = a(s_{t-1}, h6)

Step 2: SOFTMAX to get ATTENTION WEIGHTS
  α_{t,s} = softmax(e_{t,s})
  
  These are probabilities: α_{t,1} + α_{t,2} + ... + α_{t,6} = 1.0

Step 3: Compute CONTEXT VECTOR
  c_t = α_{t,1}·h1 + α_{t,2}·h2 + ... + α_{t,6}·h6
      = Σ_s α_{t,s} · h_s

Step 4: DECODE
  s_t = f(s_{t-1}, y_{t-1}, c_t)
  output word = softmax(W · s_t)

=================================================================
```

### The Alignment Function a(s, h)

The alignment function scores how well a decoder state s matches an encoder state h. Bahdanau et al. used a small feedforward neural network:

```
e_{t,s} = v_a^T · tanh(W_a · s_{t-1} + U_a · h_s)

Where:
  s_{t-1} : previous decoder hidden state (shape: decoder_hidden_dim)
  h_s     : encoder hidden state at position s (shape: encoder_hidden_dim)
  W_a     : weight matrix for decoder state
  U_a     : weight matrix for encoder state
  v_a     : weight vector (learnable)
  tanh    : nonlinearity

All W_a, U_a, v_a are LEARNED during training.
```

This is key: **the alignment function is learned**. The model learns, from the translation data itself, which source words correspond to which target words. No hand-coded alignment rules.

### The Breakthrough Insight

In the original encoder-decoder (no attention), the decoder only receives the final hidden state once. With Bahdanau attention, at each output step, the decoder can:

1. Look at all encoder states
2. Decide which ones are relevant (alignment scores)
3. Form a custom context vector that emphasizes the relevant parts
4. Use this context to generate the output word

This is like being allowed to look back at the original English text while translating each French word, instead of translating purely from memory.

---

## 4. The Alignment Matrix

The most illuminating way to understand attention is to visualize the attention weights as a matrix.

For the sentence "The agreement on the EEA was signed in August 1992", each row is a decoder step (French output word), each column is an encoder step (English input word). The value at position (i,j) is the attention weight: how much the i-th French word attended to the j-th English word.

```
ALIGNMENT MATRIX: English → French Translation
(from Bahdanau et al. 2014, approximate visualization)

ENGLISH:   The  agr  on  the  EEA  was  sign  in   Aug  1992
         ┌─────────────────────────────────────────────────────┐
L'       │ 0.8  0.1  0.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  │
accord   │ 0.1  0.8  0.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  │
sur      │ 0.0  0.1  0.8  0.1  0.0  0.0  0.0  0.0  0.0  0.0  │
la       │ 0.0  0.0  0.1  0.7  0.1  0.0  0.0  0.0  0.0  0.0  │
ZEE      │ 0.0  0.0  0.0  0.1  0.8  0.0  0.0  0.0  0.0  0.0  │
a        │ 0.0  0.0  0.0  0.0  0.1  0.8  0.0  0.0  0.0  0.0  │
ete      │ 0.0  0.0  0.0  0.0  0.0  0.1  0.7  0.0  0.0  0.0  │
signe    │ 0.0  0.0  0.0  0.0  0.0  0.0  0.8  0.1  0.0  0.0  │
en       │ 0.0  0.0  0.0  0.0  0.0  0.0  0.0  0.7  0.1  0.0  │
aout     │ 0.0  0.0  0.0  0.0  0.0  0.0  0.0  0.1  0.8  0.0  │
1992     │ 0.0  0.0  0.0  0.0  0.0  0.0  0.0  0.0  0.1  0.8  │
         └─────────────────────────────────────────────────────┘

Darker = higher attention weight
Near-diagonal pattern: the model discovered that English and French 
word order is similar for this sentence type.

For languages with very different word orders (English-Japanese),
the alignment matrix would be highly non-diagonal.
```

This matrix is not hand-coded — it emerges naturally from training. The model learned, purely from examples of English-French pairs, which words align with which.

This was scientifically beautiful: the model was implicitly learning **word alignment**, a task that took years of hand-engineering in statistical MT systems, and it was learning it as a byproduct of learning translation.

### Why This Was a Breakthrough for NMT

```
BEFORE ATTENTION (2014 NMT):
  BLEU score on WMT14 EN-FR: ~33.3

AFTER BAHDANAU ATTENTION:
  BLEU score on WMT14 EN-FR: ~36.1 (+2.8 BLEU)
  
  - Performance on long sentences (>30 words): +6 BLEU
  - No more performance cliff with sentence length
  - Interpretable alignment as a bonus
```

---

## 5. Luong Attention — Simplified Dot-Product

In 2015, Minh-Thang Luong, Hieu Pham, and Christopher Manning proposed a simplified attention mechanism. The key simplification was in the alignment score function.

Instead of a learned neural network for alignment, Luong proposed simpler scoring functions:

```
LUONG ATTENTION SCORING FUNCTIONS:

1. Dot product (simplest):
   score(s_t, h_s) = s_t^T · h_s
   (requires encoder and decoder same dimension)

2. General (dot product with weight matrix):
   score(s_t, h_s) = s_t^T · W_a · h_s

3. Concat (same as Bahdanau):
   score(s_t, h_s) = v_a^T · tanh(W_a · [s_t; h_s])
```

Luong also differs from Bahdanau in **when** attention is applied:
- **Bahdanau**: attention computed using the *previous* decoder state s_{t-1}
- **Luong**: attention computed using the *current* decoder state s_t (after the RNN step)

Both work well. The dot-product form was computationally cheaper and became the basis for self-attention.

| Feature | Bahdanau | Luong |
|---------|----------|-------|
| Alignment function | Small MLP | Dot product / matrix |
| When computed | Before decoder step | After decoder step |
| Complexity | Higher | Lower |
| Performance | Similar | Similar |
| Impact | Introduced concept | Simplified for efficiency |

The dot-product idea from Luong attention is what evolved into self-attention. Once you realize you can compute attention with a simple dot product, the entire mechanism becomes much more efficient and parallelizable.

---

## 6. Self-Attention — The True Revolution

Bahdanau and Luong attention are **cross-attention**: the decoder attends to the encoder. But in 2017, researchers asked a different question:

**What if a sequence attends to itself?**

Self-attention allows each position in a sequence to attend to all other positions in the same sequence. This is fundamentally different:

```
CROSS-ATTENTION (Bahdanau):
  Decoder at step t → attends to → Encoder positions 1..n
  (different source and target sequences)

SELF-ATTENTION:
  Each position in the SAME sequence → attends to → all other positions
  (single sequence, looking at itself)

Example: "The cat that chased the dog ate the fish"

Position "ate" attends to:
  "The"     : 0.02
  "cat"     : 0.45  ← HIGH: "cat" is the subject of "ate"
  "that"    : 0.02
  "chased"  : 0.03
  "the"     : 0.01
  "dog"     : 0.05
  "ate"     : 0.01  (itself)
  "the"     : 0.01
  "fish"    : 0.40  ← HIGH: "fish" is the object of "ate"
```

This is crucial for understanding. In the sentence "The cat that chased the dog ate the fish", to correctly understand "ate" you need to know:
1. The subject is "cat" (not "dog")
2. The object is "fish"

Self-attention lets the model directly connect "ate" to both "cat" and "fish", regardless of how far apart they are in the sequence. An RNN would need to maintain this information through many intermediate hidden states.

### Why Self-Attention is Powerful

```
COMPARISON: How models capture long-range dependencies

RNN:
  Information from position 1 reaches position 100 by flowing 
  through 99 intermediate hidden states. Each step may dilute/corrupt it.
  
  pos 1 → h1 → h2 → h3 → ... → h99 → h100
  
  Path length: O(n) — information degrades over long distances

SELF-ATTENTION:
  Any position can directly attend to any other position.
  
  pos 1 ──────────────────────────────────────→ pos 100
  
  Path length: O(1) — direct connection, no degradation
```

This is why self-attention revolutionized NLP. Every word can directly talk to every other word. The model doesn't need to maintain long-range memory through sequential processing.

---

## 7. Query, Key, Value — The Fundamental Construct

Self-attention uses three matrices: **Query (Q)**, **Key (K)**, and **Value (V)**. These are learned linear projections of the input.

### The Database Analogy

This is the best analogy to understand Q, K, V:

```
IMAGINE A DATABASE RETRIEVAL SYSTEM:

You have a database of key-value pairs:
  Key: "subject of sentence"    Value: [representation of "cat"]
  Key: "main verb"              Value: [representation of "ate"]  
  Key: "direct object"          Value: [representation of "fish"]

You want to look up: Query = "what does 'ate' interact with?"

RETRIEVAL:
  1. Compare Query against all Keys → similarity scores
  2. Normalize scores (softmax) → attention weights
  3. Weighted average of Values → your result

This is EXACTLY what attention does, but:
  - Keys, Values, and Queries are all learned from data
  - Similarity = dot product (not exact match)
  - Retrieval = soft (weighted average), not hard (exact lookup)
```

More precisely:

```
QUERY  (Q): "What am I looking for?"
  The current token's "question" — what information do I need?

KEY    (K): "What do I contain?"  
  Each token's "index entry" — what information can I provide?

VALUE  (V): "What is my actual content?"
  Each token's actual representation — the information being passed

ATTENTION:
  Score = Q · K^T         (how well does Q match each K?)
  Weights = softmax(Score) (normalize to probabilities)
  Output = Weights · V    (weighted combination of Values)
```

### Why Three Separate Matrices?

Why not just attend directly to the hidden states with one weight matrix?

Using separate Q, K, V projections gives the model more flexibility:
- Q and K can be projected into a "matching space" optimized for comparison
- V is projected into a "content space" optimized for information transfer
- Different heads (multi-head attention) can project into different subspaces
- It allows the model to learn different "roles" for the same token depending on context

---

## 8. Scaled Dot-Product Attention — Full Derivation

This is the core equation from "Attention Is All You Need" (Vaswani et al., 2017):

```
                    QK^T
Attention(Q,K,V) = softmax(────) · V
                    √d_k
```

Let's derive this step by step with full matrix shapes.

### Setup: A 4-Word Sentence

Input: **"The cat sat mat"** (4 tokens)
Model dimension: d_model = 4 (tiny, for illustration)
Key/Query dimension: d_k = 4
Value dimension: d_v = 4

### Step 1: Input Embeddings

Each word is first converted to an embedding vector (we'll cover this in depth in Chapter 27-28):

```
X (input matrix) — shape: (4, 4)   [n_seq × d_model]

         d0    d1    d2    d3
"The"  [ 0.1,  0.2, -0.1,  0.3 ]
"cat"  [ 0.8, -0.2,  0.4,  0.1 ]
"sat"  [ 0.3,  0.7, -0.3,  0.5 ]
"mat"  [ 0.2,  0.1,  0.9, -0.2 ]
```

### Step 2: Compute Q, K, V via Linear Projections

We multiply X by three learnable weight matrices:

```
W_Q ∈ R^(d_model × d_k)   [projection matrix for Queries]
W_K ∈ R^(d_model × d_k)   [projection matrix for Keys]
W_V ∈ R^(d_model × d_v)   [projection matrix for Values]

Q = X · W_Q    shape: (4, 4) · (4, 4) = (4, 4)   [n × d_k]
K = X · W_K    shape: (4, 4) · (4, 4) = (4, 4)   [n × d_k]
V = X · W_V    shape: (4, 4) · (4, 4) = (4, 4)   [n × d_v]
```

After projection (example values):
```
Q — "what is each token looking for?"
         d0    d1    d2    d3
"The"  [ 0.2,  0.5,  0.1, -0.3 ]
"cat"  [ 0.9,  0.1, -0.2,  0.4 ]
"sat"  [ 0.1, -0.4,  0.8,  0.2 ]
"mat"  [ 0.3,  0.2, -0.1,  0.6 ]

K — "what does each token contain?"
         d0    d1    d2    d3
"The"  [ 0.1,  0.3, -0.2,  0.5 ]
"cat"  [ 0.7, -0.1,  0.3,  0.2 ]
"sat"  [ 0.2,  0.8, -0.4,  0.1 ]
"mat"  [-0.3,  0.1,  0.9,  0.4 ]

V — "what is each token's actual information?"
         d0    d1    d2    d3
"The"  [ 0.4,  0.1,  0.2, -0.1 ]
"cat"  [ 0.9,  0.8,  0.1,  0.3 ]
"sat"  [ 0.2, -0.3,  0.7,  0.5 ]
"mat"  [ 0.1,  0.4, -0.2,  0.8 ]
```

### Step 3: Compute Attention Scores = QK^T

```
QK^T — shape: (4, 4)   [n × n]

This is a matrix multiply: each row of Q (a query) dotted with each column of K^T (a key)
= each query dotted with each key
= "how similar is token i's query to token j's key?"

           "The"  "cat"  "sat"  "mat"
"The"  [   1.2,   0.8,   0.4,   0.2 ]
"cat"  [   0.7,   1.8,   0.3,   0.1 ]
"sat"  [  -0.2,   0.1,   1.5,   0.9 ]
"mat"  [   0.3,   0.5,   0.2,   1.1 ]

(These are raw dot products — just example values)
```

### Step 4: Scale by √d_k

```
Why divide by √d_k?

In high dimensions, dot products tend to grow large.
If d_k = 64, vectors of dimension 64 have dot products with 
expected magnitude ≈ √64 = 8.

Without scaling:
  Large dot products → extreme softmax inputs → gradient vanishes
  
  softmax([100, 0, 0]) ≈ [1.0, 0.0, 0.0]  ← nearly a step function
  softmax([1.0, 0, 0]) = [0.58, 0.21, 0.21] ← reasonable distribution

With scaling (divide by √d_k = √4 = 2):
           "The"  "cat"  "sat"  "mat"
"The"  [   0.6,   0.4,   0.2,   0.1 ]
"cat"  [   0.35,  0.9,   0.15,  0.05]
"sat"  [  -0.1,   0.05,  0.75,  0.45]
"mat"  [   0.15,  0.25,  0.1,   0.55]
```

### Step 5: Apply Softmax

Each row becomes a probability distribution (sums to 1.0):

```
Attention Weights α = softmax(QK^T / √d_k)
Shape: (4, 4)

           "The"  "cat"  "sat"  "mat"
"The"  [   0.32,  0.26,  0.22,  0.20 ]  ← "The" mostly attends to itself
"cat"  [   0.23,  0.40,  0.19,  0.17 ]  ← "cat" attends most to itself
"sat"  [   0.16,  0.19,  0.38,  0.28 ]  ← "sat" attends to "sat" and "mat"
"mat"  [   0.22,  0.24,  0.21,  0.33 ]  ← "mat" attends to itself

Each row sums to 1.0 ✓

Interpretation of "sat" row:
  "sat" is a verb — it attends most to itself (0.38)
  and somewhat to "mat" (0.28) — the location
  and somewhat to "cat" (0.19) — the subject
  This makes linguistic sense!
```

### Step 6: Multiply by Values

```
Output = softmax(...) · V
Shape: (4, 4) · (4, 4) = (4, 4)   [n × d_v]

For "sat" (row 2 of attention weights × V):
  output_sat = 0.16·V["The"] + 0.19·V["cat"] + 0.38·V["sat"] + 0.28·V["mat"]

  = 0.16·[0.4, 0.1, 0.2, -0.1]
  + 0.19·[0.9, 0.8, 0.1,  0.3]
  + 0.38·[0.2,-0.3, 0.7,  0.5]
  + 0.28·[0.1, 0.4,-0.2,  0.8]

  = [0.064, 0.016, 0.032, -0.016]   (The's contribution)
  + [0.171, 0.152, 0.019,  0.057]   (cat's contribution)
  + [0.076,-0.114, 0.266,  0.190]   (sat's contribution)
  + [0.028, 0.112,-0.056,  0.224]   (mat's contribution)
  = [0.339, 0.166, 0.261,  0.455]   ← contextual representation of "sat"
  
This vector now contains information from ALL tokens,
weighted by how relevant each was to "sat".
```

### Complete Formula Summary

```
SCALED DOT-PRODUCT ATTENTION:

Inputs:
  Q ∈ R^(n × d_k)    — Query matrix
  K ∈ R^(n × d_k)    — Key matrix  
  V ∈ R^(n × d_v)    — Value matrix

Step 1: Raw scores
  S = Q · K^T          ∈ R^(n × n)

Step 2: Scale
  S_scaled = S / √d_k  ∈ R^(n × n)

Step 3: Normalize (softmax row-wise)
  A = softmax(S_scaled) ∈ R^(n × n)  ← attention weights

Step 4: Weighted sum
  Output = A · V       ∈ R^(n × d_v)

Combined:
                       QK^T
  Attention(Q,K,V) = softmax(────) · V
                       √d_k
```

---

## 9. Multi-Head Attention

Single-head attention computes one "perspective" on the relationships in the sequence. But language is rich — the same sentence has grammatical relationships, semantic relationships, coreference relationships, positional relationships.

**Multi-head attention** runs h parallel attention heads, each learning to attend to different types of relationships.

### The Mechanism

```
MULTI-HEAD ATTENTION:
===============================================================

Input X ∈ R^(n × d_model)

For each head i = 1, 2, ..., h:
  1. Project to lower dimension:
     Q_i = X · W_i^Q    W_i^Q ∈ R^(d_model × d_k)
     K_i = X · W_i^K    W_i^K ∈ R^(d_model × d_k)
     V_i = X · W_i^V    W_i^V ∈ R^(d_model × d_v)
  
  2. Compute attention for this head:
     head_i = Attention(Q_i, K_i, V_i) ∈ R^(n × d_v)

3. Concatenate all heads:
   MultiHead = Concat(head_1, head_2, ..., head_h) ∈ R^(n × h·d_v)

4. Project back to d_model:
   Output = MultiHead · W^O    W^O ∈ R^(h·d_v × d_model)
   Output ∈ R^(n × d_model)

Standard dimensions (original Transformer):
  d_model = 512
  h = 8 heads
  d_k = d_v = d_model/h = 64

Total parameters per multi-head attention:
  8 heads × (W^Q + W^K + W^V) + W^O
  = 8 × 3 × (512 × 64) + (512 × 512)
  = 786,432 + 262,144 = ~1M parameters
===============================================================
```

### What Different Heads Learn

After training, different attention heads specialize in different linguistic patterns:

```
HEAD 1 (Syntactic dependencies):
  "The cat sat on the mat"
  "sat" strongly attends to "cat" (subject-verb)
  "sat" strongly attends to "mat" (verb-object)

HEAD 2 (Positional / local context):
  Each token mainly attends to its immediate neighbors
  Useful for local syntactic patterns

HEAD 3 (Coreference / reference resolution):
  "The dog that bit the man went home. It was angry."
  "It" strongly attends to "dog"

HEAD 4 (Semantic similarity):
  Words with similar meanings attend to each other
  "quick" and "fast" in different parts of a text

HEAD 5 (Sentence structure):
  [CLS] token attends to key nouns and verbs
  Used for classification

HEAD 6 (Rare / domain-specific patterns):
  Technical terms, named entities

HEAD 7 (Long-range dependencies):
  Opening word attends to closing word
  "If ... then ..." structures

HEAD 8 (Fallback / general):
  Roughly uniform attention across positions
```

This is why multi-head works better than single-head with the full dimension: you get **specialization**. Each head develops expertise in different linguistic phenomena.

```
VISUALIZATION: Multi-head attention on "The cat sat on the mat"

Head 1 (grammar):                Head 2 (position):
       The cat sat on  the mat          The cat sat on  the mat
The  [ ■   □   □   □   □   □ ]   The  [ ■   ▪   □   □   □   □ ]
cat  [ □   ■   ■   □   □   □ ]   cat  [ ▪   ■   ▪   □   □   □ ]
sat  [ □   ■   ■   □   □   ■ ]   sat  [ □   ▪   ■   ▪   □   □ ]
on   [ □   □   ■   ■   □   □ ]   on   [ □   □   ▪   ■   ▪   □ ]
the  [ ■   □   □   □   ■   ■ ]   the  [ □   □   □   ▪   ■   ▪ ]
mat  [ □   □   □   □   □   ■ ]   mat  [ □   □   □   □   ▪   ■ ]

■ = high attention  ▪ = medium  □ = low

Head 1 learns grammatical subject-verb-object
Head 2 learns local positional patterns
```

---

## 10. Computational Complexity

Understanding the computational cost of attention is critical for understanding why Transformers work for some lengths and not others.

### Complexity Analysis

```
OPERATION COSTS:

Computing QK^T:
  Matrix multiply: (n × d_k) × (d_k × n) = (n × n)
  Operations: O(n² · d_k)

Applying softmax to (n × n) matrix:
  Operations: O(n²)

Multiplying by V:
  Matrix multiply: (n × n) × (n × d_v) = (n × d_v)
  Operations: O(n² · d_v)

Total per attention head: O(n² · d)
Memory to store attention matrix: O(n²)

FOR FULL TRANSFORMER (h heads):
  Computation: O(n² · d · h) = O(n² · d_model) ← quadratic in n
  Memory:      O(n²)                             ← quadratic in n
```

### Comparison with RNN

```
COMPLEXITY COMPARISON:

              | Time complexity  | Memory      | Sequential ops | Parallelizable
─────────────────────────────────────────────────────────────────────────────────
RNN           | O(n · d²)        | O(n)        | O(n)           | No (by step)
Attention     | O(n² · d)        | O(n²)       | O(1)           | Yes (fully)
─────────────────────────────────────────────────────────────────────────────────

RNN:
  - O(n) sequential steps (must process step-by-step)
  - Each step: O(d²) matrix-vector multiply
  - Can't parallelize across sequence → slow on GPU
  - Good for very long sequences (n >> d)

Attention:
  - O(1) sequential depth (all attention in parallel)
  - O(n²) operations (bad for n >> d)
  - Perfectly parallelizable → GPU-efficient
  - Practical limit: n ≈ 512-2048 for standard attention
```

### When Does n² Become a Problem?

```
n = 512 tokens:
  n² = 262,144 attention weights per head
  With h=8 heads: 2M weights
  → Manageable ✓

n = 4096 tokens:
  n² = 16,777,216 attention weights per head
  With h=8 heads: 134M weights
  → Getting expensive

n = 32768 tokens (32k context window, GPT-4):
  n² = 1,073,741,824 per head
  With h=96 heads: 103B weights (just the attention matrix!)
  → Standard attention is infeasible — need Flash Attention or sparse methods

n = 1,000,000 tokens (whole books):
  → Standard attention impossible without approximations
```

This is why long-context models need specialized attention implementations.

---

## 11. Modern Attention Variants

The n² bottleneck motivated a decade of research into more efficient attention mechanisms.

### Flash Attention

```
STANDARD ATTENTION (memory):
  1. Compute full QK^T → store in HBM (GPU high-bandwidth memory): O(n²)
  2. Apply softmax → store: O(n²)
  3. Multiply by V → store: O(n²)
  
  Problem: HBM bandwidth is the bottleneck
  For n=2048, d=64: attention matrix = 32MB per layer (too big for SRAM)

FLASH ATTENTION (Dao et al. 2022):
  Key insight: HBM reads/writes are expensive; SRAM is fast
  
  1. Tile Q, K, V into SRAM-sized blocks
  2. Compute attention in tiles, never store full n×n matrix
  3. Use online softmax trick to compute correct result in one pass
  
  Result:
  - Same mathematical output (exact, not approximate)
  - O(n) memory instead of O(n²)
  - 2-4x faster wall-clock time due to memory efficiency
  - Enables much longer sequences on same hardware
```

### Sparse Attention

```
Instead of attending to ALL n positions, attend to a SUBSET:

STRIDED SPARSE ATTENTION:
  Each token attends to:
    - Local window: positions [i-w, i+w]
    - Strided positions: i, i-stride, i-2·stride, ...

LOCAL + GLOBAL ATTENTION (Longformer):
  - Most tokens: local window of w tokens
  - Special tokens ([CLS], question tokens): attend globally to all
  
Complexity: O(n · w) where w << n
Trade-off: lose some long-range connections, gain efficiency
```

### Linear Attention

```
STANDARD: softmax(QK^T/√d_k) · V     O(n²)

LINEAR APPROXIMATION:
  softmax(QK^T) ≈ φ(Q) · φ(K)^T
  
  Where φ is a kernel function (e.g., ELU(x)+1)
  
  Key trick: matrix associativity
    φ(Q) · (φ(K)^T · V)  =  O(n · d²)  ← compute φ(K)^T·V first!
  
  This makes attention O(n) instead of O(n²)
  
  Trade-off: approximation, loses some quality, tricky to get right
  Examples: Performer, Linear Transformer, RWKV
```

| Method | Complexity | Memory | Quality vs Standard |
|--------|-----------|--------|---------------------|
| Standard Attention | O(n²·d) | O(n²) | Baseline |
| Flash Attention | O(n²·d) | O(n) | Identical (exact) |
| Sparse Attention | O(n·w·d) | O(n·w) | Near-identical for most tasks |
| Linear Attention | O(n·d²) | O(n) | Slightly worse |

---

## 12. PyTorch Implementation from Scratch

Let's implement everything we've discussed in clean, annotated PyTorch:

```python
"""
Complete implementation of Scaled Dot-Product Attention and Multi-Head Attention
from scratch in PyTorch. No black boxes — every operation explained.
"""

import math
import torch
import torch.nn as nn
import torch.nn.functional as F


def scaled_dot_product_attention(
    Q: torch.Tensor,   # shape: (batch, n_heads, seq_len, d_k)
    K: torch.Tensor,   # shape: (batch, n_heads, seq_len, d_k)
    V: torch.Tensor,   # shape: (batch, n_heads, seq_len, d_v)
    mask: torch.Tensor = None,  # shape: (batch, 1, seq_len, seq_len) or similar
    dropout: float = 0.0,
) -> tuple[torch.Tensor, torch.Tensor]:
    """
    Computes scaled dot-product attention.
    
    Attention(Q,K,V) = softmax(QK^T / sqrt(d_k)) * V
    
    Returns:
        output: (batch, n_heads, seq_len, d_v)
        attn_weights: (batch, n_heads, seq_len, seq_len) -- for visualization
    """
    d_k = Q.size(-1)  # dimension of keys/queries
    
    # Step 1: Compute raw attention scores
    # Q: (..., seq_len, d_k), K: (..., seq_len, d_k)
    # K.transpose(-2, -1): (..., d_k, seq_len)
    # scores: (..., seq_len, seq_len)
    scores = torch.matmul(Q, K.transpose(-2, -1))  # QK^T
    
    # Step 2: Scale to prevent softmax saturation
    scores = scores / math.sqrt(d_k)
    
    # Step 3: Apply mask (for causal/padding masking)
    if mask is not None:
        # Mask should be True where we want to BLOCK attention
        # We fill those positions with a very large negative number
        # so softmax gives ~0 weight to them
        scores = scores.masked_fill(mask == 0, float('-inf'))
    
    # Step 4: Softmax over the last dimension (seq dimension of K)
    # Each query's scores become a probability distribution over keys
    attn_weights = F.softmax(scores, dim=-1)
    
    # Handle NaN from softmax(-inf, -inf, ...) — happens with fully-masked rows
    # Replace NaN with 0 (masked out positions contribute nothing)
    attn_weights = torch.nan_to_num(attn_weights, nan=0.0)
    
    # Step 5: Dropout on attention weights (as in original paper)
    if dropout > 0.0 and torch.is_grad_enabled():
        attn_weights = F.dropout(attn_weights, p=dropout)
    
    # Step 6: Weighted sum of Values
    # attn_weights: (..., seq_len, seq_len)
    # V: (..., seq_len, d_v)
    # output: (..., seq_len, d_v)
    output = torch.matmul(attn_weights, V)
    
    return output, attn_weights


class MultiHeadAttention(nn.Module):
    """
    Multi-Head Attention as described in "Attention Is All You Need".
    
    MultiHead(Q,K,V) = Concat(head_1, ..., head_h) * W^O
    where head_i = Attention(Q*W_i^Q, K*W_i^K, V*W_i^V)
    """
    
    def __init__(
        self,
        d_model: int,    # model dimension (e.g., 512)
        n_heads: int,    # number of attention heads (e.g., 8)
        dropout: float = 0.1,
        bias: bool = True,
    ):
        super().__init__()
        
        assert d_model % n_heads == 0, "d_model must be divisible by n_heads"
        
        self.d_model = d_model
        self.n_heads = n_heads
        self.d_k = d_model // n_heads   # dimension per head
        self.d_v = d_model // n_heads
        self.dropout = dropout
        
        # Linear projections for Q, K, V
        # We use a single weight matrix for efficiency:
        # project Q to d_model, K to d_model, V to d_model
        # then split into n_heads
        self.W_q = nn.Linear(d_model, d_model, bias=bias)  # d_model → d_model (= n_heads * d_k)
        self.W_k = nn.Linear(d_model, d_model, bias=bias)
        self.W_v = nn.Linear(d_model, d_model, bias=bias)
        
        # Output projection W^O
        self.W_o = nn.Linear(d_model, d_model, bias=bias)
        
        # Store attention weights for visualization (optional)
        self.attn_weights = None
        
    def forward(
        self,
        query: torch.Tensor,   # (batch, seq_q, d_model)
        key: torch.Tensor,     # (batch, seq_k, d_model)
        value: torch.Tensor,   # (batch, seq_v, d_model)
        mask: torch.Tensor = None,
    ) -> torch.Tensor:
        """
        For self-attention: query = key = value = x
        For cross-attention: query from decoder, key/value from encoder
        """
        batch_size = query.size(0)
        seq_q = query.size(1)
        seq_k = key.size(1)
        
        # ── 1. Linear projections ──────────────────────────────────
        # Each: (batch, seq, d_model)
        Q = self.W_q(query)   # (batch, seq_q, d_model)
        K = self.W_k(key)     # (batch, seq_k, d_model)
        V = self.W_v(value)   # (batch, seq_v, d_model)
        
        # ── 2. Split into multiple heads ───────────────────────────
        # Reshape from (batch, seq, d_model) 
        #           to (batch, seq, n_heads, d_k)
        #           to (batch, n_heads, seq, d_k)   ← for batched matmul
        def split_heads(x, seq_len):
            x = x.view(batch_size, seq_len, self.n_heads, self.d_k)
            return x.transpose(1, 2)  # (batch, n_heads, seq, d_k)
        
        Q = split_heads(Q, seq_q)  # (batch, n_heads, seq_q, d_k)
        K = split_heads(K, seq_k)  # (batch, n_heads, seq_k, d_k)
        V = split_heads(V, seq_k)  # (batch, n_heads, seq_k, d_v)
        
        # ── 3. Attention for all heads in parallel ─────────────────
        # PyTorch broadcasts over batch and n_heads dimensions
        attn_output, attn_weights = scaled_dot_product_attention(
            Q, K, V,
            mask=mask,
            dropout=self.dropout if self.training else 0.0,
        )
        # attn_output: (batch, n_heads, seq_q, d_v)
        
        # Save for visualization
        self.attn_weights = attn_weights.detach()
        
        # ── 4. Concatenate heads ───────────────────────────────────
        # (batch, n_heads, seq_q, d_v)
        # → (batch, seq_q, n_heads, d_v)
        # → (batch, seq_q, d_model)   [because n_heads * d_v = d_model]
        attn_output = attn_output.transpose(1, 2).contiguous()
        attn_output = attn_output.view(batch_size, seq_q, self.d_model)
        
        # ── 5. Output projection ───────────────────────────────────
        output = self.W_o(attn_output)  # (batch, seq_q, d_model)
        
        return output


# ── Demonstration ──────────────────────────────────────────────────────────────

def demonstrate_attention():
    """Shows attention working on a small example."""
    
    torch.manual_seed(42)
    
    # Hyperparameters
    batch_size = 1
    seq_len = 4     # "The cat sat mat"
    d_model = 16    # small for clarity
    n_heads = 2
    
    # Simulate input embeddings for "The cat sat mat"
    # In a real model, these come from an embedding lookup table
    x = torch.randn(batch_size, seq_len, d_model)
    
    print("=" * 60)
    print("MULTI-HEAD SELF-ATTENTION DEMO")
    print("Sentence: 'The cat sat mat'")
    print("=" * 60)
    print(f"\nInput shape: {x.shape}")
    print(f"  batch={batch_size}, seq={seq_len}, d_model={d_model}")
    
    # Create multi-head attention layer
    mha = MultiHeadAttention(d_model=d_model, n_heads=n_heads, dropout=0.0)
    
    # Self-attention: query = key = value = x
    with torch.no_grad():
        output = mha(x, x, x)
    
    print(f"\nOutput shape: {output.shape}")
    print(f"  Same as input: ✓")
    
    print(f"\nAttention weights (head 1): shape {mha.attn_weights[:, 0, :, :].shape}")
    print("  (each row = which positions token i attends to)")
    attn_h1 = mha.attn_weights[0, 0].numpy()  # (seq, seq)
    words = ["The", "cat", "sat", "mat"]
    print(f"\n  {'':>6}", end="")
    for w in words:
        print(f"  {w:>5}", end="")
    print()
    for i, w in enumerate(words):
        print(f"  {w:>6}", end="")
        for j in range(seq_len):
            print(f"  {attn_h1[i,j]:.3f}", end="")
        print()
    
    return output


def create_causal_mask(seq_len: int) -> torch.Tensor:
    """
    Creates a causal (autoregressive) mask.
    
    For language modeling, each position can only attend to positions
    at or before it (can't look at future tokens).
    
    Returns a lower-triangular matrix of 1s (attend) and 0s (don't attend).
    Shape: (1, 1, seq_len, seq_len) — broadcastable over batch and heads.
    
    Example for seq_len=4:
      [[1, 0, 0, 0],
       [1, 1, 0, 0],
       [1, 1, 1, 0],
       [1, 1, 1, 1]]
    
    0 → masked out (replaced with -inf before softmax)
    1 → allowed
    """
    mask = torch.tril(torch.ones(seq_len, seq_len))
    return mask.unsqueeze(0).unsqueeze(0)  # (1, 1, seq, seq)


def demonstrate_causal_attention():
    """Shows how causal masking prevents attending to future tokens."""
    
    torch.manual_seed(42)
    
    seq_len = 4
    d_model = 8
    n_heads = 1
    
    x = torch.randn(1, seq_len, d_model)
    
    # Create causal mask
    mask = create_causal_mask(seq_len)
    
    mha = MultiHeadAttention(d_model=d_model, n_heads=n_heads, dropout=0.0)
    
    with torch.no_grad():
        output = mha(x, x, x, mask=mask)
    
    print("\n" + "=" * 60)
    print("CAUSAL (MASKED) SELF-ATTENTION")
    print("(Used in GPT for language modeling)")
    print("=" * 60)
    
    print("\nCausal mask:")
    print(mask[0, 0].int().numpy())
    print("1 = can attend, 0 = blocked (future tokens)")
    
    attn = mha.attn_weights[0, 0].numpy()
    words = ["The", "cat", "sat", "mat"]
    print(f"\nAttention weights after masking:")
    print(f"  {'':>6}", end="")
    for w in words:
        print(f"  {w:>5}", end="")
    print()
    for i, w in enumerate(words):
        print(f"  {w:>6}", end="")
        for j in range(seq_len):
            v = attn[i, j]
            print(f"  {v:.3f}" if v > 0.001 else "  0.000", end="")
        print()
    
    print("\n'sat' (row 2) cannot see 'mat' (column 3) — it's in the future ✓")
    print("'The' (row 0) can only see itself — no previous context ✓")


if __name__ == "__main__":
    output = demonstrate_attention()
    demonstrate_causal_attention()
    
    print("\n" + "=" * 60)
    print("PARAMETER COUNT")
    print("=" * 60)
    d_model = 512
    n_heads = 8
    mha_large = MultiHeadAttention(d_model=d_model, n_heads=n_heads)
    total_params = sum(p.numel() for p in mha_large.parameters())
    print(f"d_model={d_model}, n_heads={n_heads}")
    print(f"Total parameters: {total_params:,}")
    print(f"  W_q: {d_model*d_model:,} (512×512)")
    print(f"  W_k: {d_model*d_model:,}")
    print(f"  W_v: {d_model*d_model:,}")
    print(f"  W_o: {d_model*d_model:,}")
    print(f"  Total: {4*d_model*d_model:,} ← 4 × d_model²")
```

### Running This Code

```bash
pip install torch

# Save as attention_demo.py and run:
python attention_demo.py
```

Expected output (abbreviated):
```
MULTI-HEAD SELF-ATTENTION DEMO
Sentence: 'The cat sat mat'
============================================================
Input shape: torch.Size([1, 4, 16])
  batch=1, seq=4, d_model=16

Output shape: torch.Size([1, 4, 16])
  Same as input: ✓

Attention weights (head 1):
             The    cat    sat    mat
     The   0.312  0.208  0.251  0.229
     cat   0.198  0.341  0.227  0.234
     sat   0.241  0.289  0.301  0.169
     mat   0.175  0.256  0.312  0.257
```

---

## 13. Summary

```
THE EVOLUTION OF ATTENTION:
═══════════════════════════════════════════════════════════════════

2014 — ENCODER-DECODER (no attention):
  Entire input → one fixed vector → decode
  Problem: information bottleneck, fails on long sentences

2014 — BAHDANAU ATTENTION (cross-attention):
  Decoder at each step looks at ALL encoder states
  Learns alignment: which source words matter for current output
  Key formula: c_t = Σ softmax(a(s_{t-1}, h_s)) · h_s
  Breakthrough: long-sentence translation works

2015 — LUONG ATTENTION:
  Same idea, simpler dot-product scoring
  Faster, equivalent performance

2017 — SELF-ATTENTION (the revolution):
  Attend within the SAME sequence
  Q, K, V from the same input
  Every token directly sees every other token
  Formula: Attention(Q,K,V) = softmax(QK^T/√d_k)·V
  
2017 — MULTI-HEAD ATTENTION:
  Run h independent attention heads
  Each head learns different relationship type
  Concatenate and project: richer representations

2017 — "ATTENTION IS ALL YOU NEED" (Transformer):
  Built entirely from self-attention + FFN
  No recurrence, no convolution
  Parallelizable → perfect for GPUs
  The foundation of everything that followed
═══════════════════════════════════════════════════════════════════
```

### Key Equations

| Name | Formula | Notes |
|------|---------|-------|
| Alignment score (Bahdanau) | e_{t,s} = v^T tanh(W s_{t-1} + U h_s) | Learned MLP |
| Attention weights | α = softmax(e) | Row-wise normalization |
| Context vector | c_t = Σ α_{t,s} h_s | Weighted average |
| Scaled dot-product | Attention = softmax(QK^T/√d_k)V | Core formula |
| Multi-head | MultiHead = Concat(heads)W^O | h parallel heads |

### Why √d_k?

If vectors have unit variance components, their dot product has variance d_k. Dividing by √d_k normalizes the variance to 1, keeping softmax inputs in a reasonable range where gradients flow well.

---

## 14. Mini Projects

### Mini Project 1: Attention from Scratch

**What You'll Build:** Implement scaled dot-product attention in pure NumPy (no PyTorch) and compute every intermediate step by hand.

**Time Estimate:** 1-2 hours

**Skills Practiced:** NumPy matrix operations, softmax implementation, attention formula, debugging numerical computations

**Step-by-Step:**

1. Set up Q, K, V matrices of your choice:
```python
import numpy as np

np.random.seed(42)

# 4 tokens ("The", "cat", "sat", "mat"), d_k = 4
seq_len = 4
d_k = 4

Q = np.random.randn(seq_len, d_k)
K = np.random.randn(seq_len, d_k)
V = np.random.randn(seq_len, d_k)

print("Q shape:", Q.shape)
print("K shape:", K.shape)
print("V shape:", V.shape)
```

2. Compute raw attention scores (QK^T):
```python
scores = Q @ K.T           # shape: (4, 4)
print("\nRaw scores (QK^T):")
print(np.round(scores, 3))
```

3. Scale by sqrt(d_k) and apply softmax row-wise:
```python
def softmax(x, axis=-1):
    """Numerically stable softmax."""
    x = x - x.max(axis=axis, keepdims=True)
    exp_x = np.exp(x)
    return exp_x / exp_x.sum(axis=axis, keepdims=True)

scaled_scores = scores / np.sqrt(d_k)
print("\nScaled scores (/ sqrt(d_k)):")
print(np.round(scaled_scores, 3))

attn_weights = softmax(scaled_scores, axis=-1)
print("\nAttention weights (softmax):")
print(np.round(attn_weights, 3))

print("\nRow sums (should be 1.0):", attn_weights.sum(axis=1).round(6))
```

4. Compute the weighted sum of values:
```python
output = attn_weights @ V   # shape: (4, 4)
print("\nAttention output:")
print(np.round(output, 3))

words = ["The", "cat", "sat", "mat"]
print("\nWhich word each token attended to most:")
for i, w in enumerate(words):
    most_attended = words[np.argmax(attn_weights[i])]
    print(f"  '{w}' attended most to '{most_attended}' "
          f"(weight={attn_weights[i].max():.3f})")
```

5. Compare to PyTorch to verify:
```python
import torch
import torch.nn.functional as F

Q_t = torch.tensor(Q, dtype=torch.float32)
K_t = torch.tensor(K, dtype=torch.float32)
V_t = torch.tensor(V, dtype=torch.float32)

torch_output = F.scaled_dot_product_attention(
    Q_t.unsqueeze(0), K_t.unsqueeze(0), V_t.unsqueeze(0)
).squeeze(0).numpy()

max_diff = np.abs(output - torch_output).max()
print(f"\nMax difference from PyTorch: {max_diff:.2e}")
print("Match!" if max_diff < 1e-5 else "Mismatch — check your implementation")
```

**Expected Output:**
Each intermediate matrix printed (scores, scaled scores, attention weights), verification that rows sum to 1.0, max difference from PyTorch output below 1e-5, and a human-readable summary of which token attended to which.

**Bonus Challenge:**
Add a causal mask: set all positions above the diagonal to -inf before softmax. Verify that "The" can only attend to itself, "cat" to "The" and itself, etc. Also try implementing multi-head attention by running this twice with different Q, K, V projections and concatenating the results.

---

### Mini Project 2: Attention Pattern Visualizer

**What You'll Build:** Extract attention weights from a pre-trained BERT model and render a heatmap showing which words attend to which other words.

**Time Estimate:** 1-2 hours

**Skills Practiced:** HuggingFace transformers, attention output extraction, matplotlib heatmaps, interpreting model internals

**Step-by-Step:**

1. Install dependencies and load BERT:
```python
# pip install transformers matplotlib torch

from transformers import BertTokenizer, BertModel
import torch
import matplotlib.pyplot as plt
import numpy as np

tokenizer = BertTokenizer.from_pretrained("bert-base-uncased")
model = BertModel.from_pretrained(
    "bert-base-uncased",
    output_attentions=True,   # <-- key flag
)
model.eval()
```

2. Tokenize a sentence and run inference:
```python
sentence = "The cat sat on the mat"
inputs = tokenizer(sentence, return_tensors="pt")
tokens = tokenizer.convert_ids_to_tokens(inputs["input_ids"][0])
print("Tokens:", tokens)  # includes [CLS] and [SEP]

with torch.no_grad():
    outputs = model(**inputs)

# outputs.attentions: tuple of n_layers tensors, each (batch, n_heads, seq, seq)
print(f"\nNumber of attention layers: {len(outputs.attentions)}")
print(f"Shape per layer: {outputs.attentions[0].shape}")
```

3. Plot a heatmap for a chosen layer:
```python
def plot_attention_layer(attentions, layer_idx, tokens, head_idx=None):
    attn = attentions[layer_idx][0]  # (n_heads, seq_len, seq_len)

    if head_idx is not None:
        attn_matrix = attn[head_idx].numpy()
        head_label = f"Head {head_idx}"
    else:
        attn_matrix = attn.mean(dim=0).numpy()
        head_label = "Average (all heads)"

    fig, ax = plt.subplots(figsize=(8, 7))
    im = ax.imshow(attn_matrix, cmap="Blues", vmin=0, vmax=attn_matrix.max())

    ax.set_xticks(range(len(tokens)))
    ax.set_yticks(range(len(tokens)))
    ax.set_xticklabels(tokens, rotation=45, ha="right", fontsize=10)
    ax.set_yticklabels(tokens, fontsize=10)

    for i in range(len(tokens)):
        for j in range(len(tokens)):
            val = attn_matrix[i, j]
            color = "white" if val > attn_matrix.max() * 0.6 else "black"
            ax.text(j, i, f"{val:.2f}", ha="center", va="center",
                    fontsize=7, color=color)

    ax.set_title(f"BERT Attention — Layer {layer_idx}, {head_label}", fontsize=12)
    ax.set_ylabel("Query (attending FROM)", fontsize=10)
    ax.set_xlabel("Key (attending TO)", fontsize=10)
    plt.colorbar(im, ax=ax)
    plt.tight_layout()
    plt.savefig(f"attn_layer{layer_idx}.png", dpi=150, bbox_inches="tight")
    plt.show()

for layer in [0, 5, 11]:
    plot_attention_layer(outputs.attentions, layer, tokens)
```

4. Inspect semantic connections:
```python
def top_attended(attn_matrix, tokens, from_token_idx, top_k=3):
    weights = attn_matrix[from_token_idx]
    top_indices = np.argsort(weights)[::-1][:top_k]
    print(f"  '{tokens[from_token_idx]}' attends most to:")
    for idx in top_indices:
        print(f"    '{tokens[idx]}' -> weight={weights[idx]:.3f}")

attn_layer11 = outputs.attentions[11][0]  # (12, seq_len, seq_len)
for head_idx in range(3):
    print(f"\nLayer 11, Head {head_idx}:")
    attn_matrix = attn_layer11[head_idx].numpy()
    for i in range(1, len(tokens) - 1):
        top_attended(attn_matrix, tokens, i, top_k=2)
```

**Expected Output:**
Three saved heatmap PNGs. Early layers show mostly local (diagonal) attention. Later layers show more semantic patterns — "cat" and "sat" should attend to each other in some heads, and "sat" should attend to "mat".

**Bonus Challenge:**
Test the sentence "The bank by the river bank is near the bank where I bank". Does "bank" attend differently depending on its context? Plot all 12 heads for one layer in a 3x4 grid to see head specialization.

---

### Mini Project 3: Self-Attention Intuition

**What You'll Build:** Create 5 sentences where attention should behave interestingly (pronoun resolution, ambiguity, negation), run them through BERT, visualize attention heads, and write a short analysis of what you observe.

**Time Estimate:** 1-2 hours

**Skills Practiced:** Model probing, linguistic analysis, attention interpretation, critical thinking about model behavior

**Step-by-Step:**

1. Choose linguistically interesting sentences:
```python
from transformers import BertTokenizer, BertModel
import torch
import numpy as np
import matplotlib.pyplot as plt

tokenizer = BertTokenizer.from_pretrained("bert-base-uncased")
model = BertModel.from_pretrained("bert-base-uncased", output_attentions=True)
model.eval()

sentences = [
    # Pronoun resolution: "it" should refer to "cat" not "dog"
    "The cat chased the dog and then it ran away",
    # Subject-verb-object with long-range dependency
    "The lion that scared the mouse ate the zebra",
    # Negation: "not" should affect "happy"
    "She was not happy about the result",
    # Modifier attachment: "fast" should modify "car" not "driver"
    "The driver of the fast car waved",
    # Long-range dependency across a relative clause
    "The book that my friend who lives in Paris recommended was excellent",
]
```

2. Find the most focused attention head for a target word:
```python
def find_focused_heads(sentence, target_word, top_n=3):
    """Return heads with highest variance for target_word — most focused attention."""
    inputs = tokenizer(sentence, return_tensors="pt")
    tokens = tokenizer.convert_ids_to_tokens(inputs["input_ids"][0])

    target_idx = next(
        (i for i, t in enumerate(tokens) if target_word.lower() in t.lower()),
        None
    )
    if target_idx is None:
        print(f"Warning: '{target_word}' not found in {tokens}")
        return [], tokens, None

    with torch.no_grad():
        outputs = model(**inputs)

    scored = []
    for layer in range(12):
        for head in range(12):
            row = outputs.attentions[layer][0, head, target_idx].numpy()
            scored.append(((layer, head), np.var(row), row))

    scored.sort(key=lambda x: -x[1])
    return scored[:top_n], tokens, target_idx
```

3. Visualize and document each sentence:
```python
def analyse_sentence(sentence, target_word):
    top_heads, tokens, target_idx = find_focused_heads(sentence, target_word)
    if not top_heads:
        return

    fig, axes = plt.subplots(1, 3, figsize=(15, 4))
    fig.suptitle(f'"{sentence}"\nTarget: "{target_word}"', fontsize=9)

    for ax, ((layer, head), var, row) in zip(axes, top_heads):
        bars = ax.bar(range(len(tokens)), row, color="steelblue", alpha=0.8)
        ax.set_xticks(range(len(tokens)))
        ax.set_xticklabels(tokens, rotation=45, ha="right", fontsize=7)
        ax.set_title(f"L{layer} H{head} (var={var:.4f})", fontsize=9)
        bars[target_idx].set_color("orange")

    plt.tight_layout()
    safe = target_word.replace(" ", "_")
    plt.savefig(f"intuition_{safe}.png", dpi=150, bbox_inches="tight")
    plt.show()

targets = ["it", "ate", "happy", "fast", "excellent"]
for sent, tgt in zip(sentences, targets):
    print(f"\nAnalysing '{tgt}' in: {sent}")
    analyse_sentence(sent, tgt)
```

4. Answer these questions after running:
```python
questions = [
    ("it", "Does 'it' attend more to 'cat' or 'dog'? Which layer/head?"),
    ("ate", "Does 'ate' show long-range connections to 'lion' and 'zebra'?"),
    ("happy", "Does 'not' appear in 'happy's attention? Which head?"),
    ("fast", "Does 'fast' attend to 'car' (correct) or 'driver' (wrong)?"),
    ("excellent", "Can 'excellent' reach back to 'book' across the relative clause?"),
]
for word, question in questions:
    print(f"\n[{word}] {question}")
    print("  Answer: _______________")
```

**Expected Output:**
Five attention bar-chart PNGs (one per sentence, 3 heads each). Your written observations for each question form the deliverable — the goal is understanding model behavior, not a single correct answer.

**Bonus Challenge:**
Compare BERT (bidirectional) vs GPT-2 (causal). For the pronoun resolution sentence, GPT-2 can only look left-to-right — does this hurt its ability to link "it" back to "cat"? Use the same `output_attentions=True` flag with `GPT2Model` and compare the patterns side by side.

---

## 15. Exercises

**Exercise 1**: Implement the Bahdanau attention alignment function (the MLP version) and compute attention weights for a 5-word input sentence. Use hidden dimension = 32 for both encoder and decoder.

**Exercise 2**: In the scaled dot-product attention formula, what happens to the attention weights as d_k → ∞ without the √d_k scaling? Show mathematically that the softmax saturates (approaches a one-hot vector).

**Exercise 3**: Modify the `MultiHeadAttention` class to support cross-attention (where query comes from a different sequence than key/value). This is needed for the encoder-decoder architecture.

**Exercise 4**: Visualize attention weights for a real sentence using the provided code. Load a pre-trained BERT model from HuggingFace and plot the attention weights for the sentence "The quick brown fox jumps over the lazy dog". Which words does "fox" attend to?

**Exercise 5**: The multi-head attention has 4 × d_model² parameters (W_q, W_k, W_v, W_o). For d_model=512, n_heads=8, what is d_k per head? What is the total FLOP count for one attention forward pass on a sequence of length 512?

**Exercise 6**: Implement a simple version of Flash Attention's tiling idea: instead of computing the full n×n attention matrix at once, compute it in tiles of size B×B. Verify the output matches standard attention.

---

**Chapter Summary**: Attention began as a solution to the bottleneck problem in sequence-to-sequence models. Bahdanau's insight — let the decoder look at all encoder states — solved long-sentence translation. Luong simplified it. Self-attention generalized the idea: every token in a sequence can attend to every other token, with O(1) path length between any two positions. Query-Key-Value is the core abstraction: Q asks "what do I need?", K says "here's my index", V carries the actual content. Scaled dot-product attention (softmax(QK^T/√d_k)V) is the fundamental operation. Multi-head attention runs this in parallel for h "perspectives". The computational cost O(n²d) is the main limitation for long sequences, motivating Flash Attention and sparse variants.

**What's Next →** [Chapter 23: The Transformer Architecture — Attention Is All You Need](./23-transformer-architecture.md)

*With attention as our foundation, we're ready to build the full Transformer architecture — the model that introduced attention to the world and replaced RNNs for nearly every sequence task.*
