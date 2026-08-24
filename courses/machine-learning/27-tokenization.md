# Chapter 27: Tokenization — How Text Becomes Numbers

> **"Tokenization is the invisible foundation of every language model. Most LLM bugs and surprises trace back to tokenization. Every developer using AI APIs should understand it."**

---

## Table of Contents
1. [Why Tokenization Exists](#1-why-tokenization-exists)
2. [Character-Level Tokenization](#2-character-level-tokenization)
3. [Word-Level Tokenization](#3-word-level-tokenization)
4. [Subword Tokenization — The Sweet Spot](#4-subword-tokenization--the-sweet-spot)
5. [BPE (Byte Pair Encoding)](#5-bpe-byte-pair-encoding)
6. [WordPiece (BERT)](#6-wordpiece-bert)
7. [SentencePiece](#7-sentencepiece)
8. [tiktoken (OpenAI)](#8-tiktoken-openai)
9. [Special Tokens](#9-special-tokens)
10. [Training Your Own Tokenizer](#10-training-your-own-tokenizer)
11. [Tokenization Quirks and Gotchas](#11-tokenization-quirks-and-gotchas)
12. [Full Code: Build BPE Tokenizer from Scratch](#12-full-code-build-bpe-tokenizer-from-scratch)
13. [Summary](#13-summary)
14. [Exercises](#14-exercises)

---

## 1. Why Tokenization Exists

Neural networks operate on fixed-size numerical vectors. Text is a sequence of symbols with no natural numeric representation. Tokenization bridges this gap.

```
THE PROBLEM:

Raw text: "Hello, world!"
           ↓ ???
Neural network input: [vector_1, vector_2, ..., vector_n]
  where each vector is a sequence of floating-point numbers

TWO REQUIREMENTS:
  1. DISCRETE: text must become integers (which are looked up in embedding tables)
  2. FINITE VOCABULARY: can't have a different integer for every possible string
     → must limit to a fixed vocabulary of V tokens

THE CHALLENGE:
  Natural language has:
    English dictionary: ~170,000 words
    Proper nouns, technical terms, names: millions more
    Morphological variants: "play", "plays", "played", "playing", "playable", ...
    New words coined daily: "selfie", "COVID", "cryptocurrency"
    Multiple languages: 7,000+ languages
    Code, math, HTML, JSON: completely different symbol sets
    
  A vocabulary of 50,000 can't cover all of this!
  → Need a strategy that handles ANY text with a bounded vocabulary.
```

The challenge is to design a tokenization scheme that:
1. Has a bounded vocabulary (typically 30k-200k tokens)
2. Can represent ANY input text (no "unknown" tokens ideally)
3. Produces short sequences (shorter = less compute)
4. Groups meaningful linguistic units together
5. Is reversible: encode then decode returns the original text

---

## 2. Character-Level Tokenization

The simplest possible approach: each character is one token.

```
CHARACTER TOKENIZATION:

Vocabulary: all unique characters in training data
  For English: a-z, A-Z, 0-9, punctuation, space = ~100 characters
  With Unicode: thousands of possible characters
  
Example:
  Text: "hello"
  Tokens: ['h', 'e', 'l', 'l', 'o']
  IDs: [8, 5, 12, 12, 15]  (arbitrary IDs from vocab)
  
More complete example:
  Text: "The cat sat."
  Tokens: ['T','h','e',' ','c','a','t',' ','s','a','t','.']
  Length: 12 tokens
  
  (Compare: "The cat sat." → 4 tokens with word-level)

CODE:
  class CharTokenizer:
      def __init__(self, text):
          self.vocab = sorted(set(text))
          self.stoi = {c: i for i, c in enumerate(self.vocab)}
          self.itos = {i: c for i, c in enumerate(self.vocab)}
      
      def encode(self, text):
          return [self.stoi[c] for c in text]
      
      def decode(self, ids):
          return ''.join(self.itos[i] for i in ids)
  
  tokenizer = CharTokenizer("hello world")
  print(tokenizer.encode("hello"))  # [3, 1, 7, 7, 8]
  print(tokenizer.decode([3,1,7,7,8]))  # "hello"
```

### Pros and Cons

```
CHARACTER TOKENIZATION:

PROS:
  ✓ Tiny vocabulary (~100-1000 characters)
  ✓ Handles ANY text — unknown words, typos, emojis, any Unicode
  ✓ No OOV (out-of-vocabulary) problem ever
  ✓ Full reversibility
  ✓ Language-agnostic

CONS:
  ✗ Very long sequences: "transformer" = 11 tokens (vs 1-2 with subword)
    → Quadratic attention cost: 11² = 121 vs 2² = 4 (5.5× more compute!)
  ✗ Less semantic grouping: "un", "play", "able" have meaning; 'u','n' don't
  ✗ Model must learn to spell: learns "h"+"e"+"l"+"l"+"o" = word before
    learning word semantics (longer training, less efficient)
  ✗ Context window fills up quickly

WHEN USED:
  - Character-level language models (research)
  - Very small models where simplicity matters
  - Our TinyGPT example (Chapter 28) uses character-level for simplicity
  - DNA/protein sequence models (alphabet of 4 or 20 characters)
```

---

## 3. Word-Level Tokenization

Split text at word boundaries (spaces, punctuation):

```
WORD TOKENIZATION:

Vocabulary: all unique words in training corpus
  
  "The cat sat on the mat." 
  → ["The", "cat", "sat", "on", "the", "mat", "."]
  → IDs: [1432, 4218, 8901, 234, 1, 5643, 1]
  
  (Note: "The" and "the" are different tokens here — case-sensitive)
  
TOKENIZER RULES:
  Simple: split on whitespace
  Better: also split punctuation, handle contractions
    "don't" → ["don", "'", "t"]  or  ["don't"]  (depends on implementation)
    
TYPICAL VOCABULARY SIZES:
  English news corpus: 50,000-200,000 unique words
  Including names, technical terms: millions
  Practical limit: use top 50k by frequency, rest → [UNK]
```

### Pros and Cons

```
WORD TOKENIZATION:

PROS:
  ✓ Intuitive: one word = one token (usually)
  ✓ Short sequences: "transformer architecture" = 2 tokens
  ✓ Good semantic grouping: whole meaningful units

CONS:
  ✗ HUGE vocabulary: 50k+ words, larger embedding table
  ✗ OOV problem: "COVID-19", "BERT", "GPT", new words → [UNK]
  ✗ Morphology blind:
      "play", "plays", "playing", "played", "playable" → 5 separate tokens
      with no knowledge that they're related
  ✗ Language-specific: word boundaries differ by language
      Chinese, Japanese, Thai: no spaces between words!
  ✗ Rare words undertrained: "exsanguination" appears 3 times in corpus
      → word embedding barely trained, unreliable

OOV PROBLEM ILLUSTRATION:
  Training vocabulary: {the, cat, sat, on, mat, ...} (50k words)
  
  At inference: "The COVID-19 pandemic affected everyone."
  Tokenized: ["The", [UNK], "pandemic", "affected", "everyone", "."]
  
  "COVID-19" → [UNK]. The model has NO information about this.
  This was a real problem during 2020 with models trained before COVID.
```

---

## 4. Subword Tokenization — The Sweet Spot

The key insight that solved tokenization:

```
THE SUBWORD IDEA:
  Common words → single token: "the", "and", "is"
  Uncommon words → split into common subwords: 
    "unhappiness" → ["un", "##happiness"] or ["un", "happy", "##ness"]
  
  Vocabulary: ~32k-128k subword units
  
  PROPERTIES:
    - Common words: 1 token (efficient)
    - Rare words: 2-4 tokens (still manageable)
    - Truly unknown strings: split into characters (no UNK!)
    - Shared subwords: "play" appears in "player", "replay", "playful"
      → model learns shared morpheme representation
    
  EXAMPLE:
    "unhappiness"
    
    Word-level: [UNK] if not in vocabulary
    Character-level: ['u','n','h','a','p','p','i','n','e','s','s'] = 11 tokens
    Subword (BPE): ["un", "happiness"] = 2 tokens  ← SWEET SPOT!
    
    "electroencephalography" (technical word)
    Character-level: 22 tokens
    Subword (GPT-2): ["electro", "enc", "eph", "alo", "graph", "y"] = 6 tokens ← still manageable
```

---

## 5. BPE (Byte Pair Encoding)

BPE is the most widely used subword tokenization algorithm. It was originally a data compression algorithm, repurposed for NLP by Sennrich et al. (2016).

### Algorithm

```
BPE ALGORITHM:
Step 1: Initialize vocabulary with all individual characters.
Step 2: Count all consecutive character pair frequencies in corpus.
Step 3: Find the most frequent pair.
Step 4: Merge that pair into a new token, add to vocabulary.
Step 5: Update all occurrences in corpus.
Step 6: Repeat from Step 2 until vocabulary reaches target size.

EXAMPLE — Step by Step:

Initial corpus (4 words with frequencies):
  "low" ×5, "lower" ×2, "newest" ×6, "widest" ×3
  
  Word-ending marker: add </w> to end of each word
  (marks word boundaries, important for tokenization)
  
  Initial characters:
    "l o w </w>":5, "l o w e r </w>":2, "n e w e s t </w>":6, "w i d e s t </w>":3

ITERATION 1:
  Count all consecutive pairs:
    l-o: 7 (from "low"×5 + "lower"×2)
    o-w: 7
    e-s: 9 (from "newest"×6 + "widest"×3)
    s-t: 9
    w-</w>: 5
    ...
  
  Most frequent: "e-s" and "s-t" are tied at 9. Ties are broken by a fixed
  rule (e.g. first-encountered pair) — here we pick "e-s".
  
  MERGE "e" + "s" → "es"
  Updated corpus:
    "l o w </w>":5, "l o w e r </w>":2, 
    "n e w es t </w>":6, "w i d es t </w>":3
    
  Vocabulary now: l, o, w, </w>, e, r, n, s, t, i, d, es

ITERATION 2:
  Count pairs again:
    (es, t): 9 (from "newest" and "widest")
    (l, o): 7
    ...
  
  Most frequent: "es-t" = 9.
  MERGE "es" + "t" → "est"
  Updated corpus:
    "l o w </w>":5, "l o w e r </w>":2, 
    "n e w est </w>":6, "w i d est </w>":3

ITERATION 3:
  Pairs:
    (l, o): 7
    (o, w): 7
    (n, e): 6
    (e, w): 6
    (est, </w>): 9   ← highest
    
  MERGE "est" + "</w>" → "est</w>"
  Updated corpus:
    "l o w </w>":5, "l o w e r </w>":2, 
    "n e w est</w>":6, "w i d est</w>":3

ITERATION 4:
  Most frequent remaining pair: (l, o) = 7
  MERGE "l" + "o" → "lo"

ITERATION 5:
  Most frequent remaining pair: (lo, w) = 7
  MERGE "lo" + "w" → "low"

...continue until vocabulary size reached (e.g., 50,000)
```

### Final Tokenization

```
TO TOKENIZE NEW TEXT:

Given the learned merge rules:
  (e,s) → es
  (es,t) → est
  (lo) → lo
  (lo, w) → low
  (lo, wer) → lower
  ... etc

Algorithm:
  1. Split word into characters
  2. Apply merge rules in ORDER (the order they were learned!)
  3. Apply first rule that matches, keep applying until no more matches

"lowest":
  l o w e s t </w>
  Apply (e,s)→es:    l o w es t </w>
  Apply (es,t)→est:  l o w est </w>
  Apply (l,o)→lo:    lo w est </w>
  Apply (lo,w)→low:  low est </w>
  Result: ["low", "est</w>"]
  ← "lowest" = "low" + "est" (suffix) — morphologically meaningful!

"unhappiness":
  u n h a p p i n e s s </w>
  Apply relevant merges: ... → ["un", "happi", "ness</w>"]
```

### Byte-Level BPE

Modern models (GPT-2, GPT-4, LLaMA) use **byte-level BPE**:

```
BYTE-LEVEL BPE:
  Instead of starting with Unicode characters, start with bytes.
  
  Base vocabulary = 256 bytes (0x00 to 0xFF)
  BPE merges on top of bytes.
  
  ADVANTAGE:
    ANY Unicode character can be encoded as bytes (UTF-8).
    No [UNK] token possible — everything has a byte representation.
    Handles: emojis, Chinese characters, Arabic, code symbols, etc.
    
  ENCODING EXAMPLE:
    "hello" = [104, 101, 108, 108, 111] in bytes (ASCII values)
    After BPE merges: "hello" might be a single token if common enough.
    
  "😀" (emoji) = [0xF0, 0x9F, 0x98, 0x80] in UTF-8 bytes = [240, 159, 152, 128]
    → tokenized as up to 4 byte tokens: Ġ, ŀ (or whatever tokens these bytes merged into)
    
  "私は" (Japanese) = multiple UTF-8 bytes
    → tokenized via bytes, not assuming language-specific word boundaries
    
  WHY GPT-2 INTRODUCED THIS:
    With character-level BPE, rare Unicode characters are [UNK].
    With byte-level BPE, the absolute WORST case is 4 bytes per Unicode char.
    Every character is representable.
```

### BPE in Practice

```
MODELS USING BPE:
  GPT-2:        50,257 tokens (byte-level BPE)
  GPT-3:        same tokenizer as GPT-2
  GPT-4:        100,277 tokens (cl100k_base, tiktoken)
  LLaMA 1/2:    32,000 tokens (SentencePiece BPE)
  LLaMA 3:      128,256 tokens (tiktoken-style)
  Mistral:      32,000 tokens (SentencePiece BPE, same as LLaMA)

COMPRESSION RATIO (English text):
  Character-level: 1 char/token
  GPT-2 (50k): ~4 chars/token
  GPT-4 (100k): ~4-5 chars/token
  BERT (30k): ~4 chars/token

PRACTICAL IMPACT:
  "The quick brown fox jumps over the lazy dog"
  
  GPT-2 tokens:
  ['The', 'Ġquick', 'Ġbrown', 'Ġfox', 'Ġjumps', 'Ġover', 'Ġthe', 'Ġlazy', 'Ġdog']
  9 tokens (Ġ = space prefix in GPT-2)
  
  Character-level: 43 tokens
  Word-level: 9 tokens (same, coincidentally)
```

---

## 6. WordPiece (BERT)

WordPiece is similar to BPE but uses a different merge criterion:

```
WORDPIECE vs BPE:

BPE: merge the most FREQUENT pair
WordPiece: merge the pair with highest SCORE:
  
  score(a, b) = freq(a, b) / (freq(a) × freq(b))
  
  HIGH SCORE: a and b appear together MORE OFTEN than by chance
  LOW SCORE: a and b are common individually but not specifically together
  
  This slightly favors linguistically meaningful merges.
  "##ing" gets merged eagerly because it almost always follows verb stems.
  
BERT TOKENIZATION SPECIFICS:

  1. Lowercase input (BERT-uncased only)
  2. Normalize whitespace
  3. Split on whitespace to get "words"
  4. Apply WordPiece greedily to each word:
     Start: full word → try to find in vocabulary
     If not found: split off longest prefix that IS in vocabulary
     Mark continuation tokens with "##"
  
  GREEDY LONGEST MATCH:
    "playing" → try "playing" (not in vocab) → try "playin" → ... → "play" (found!)
    Remaining: "ing" → try "##ing" (found!)
    Result: ["play", "##ing"]
    
    "transformer":
    → try "transformer" (not found as single token)
    → "transform" (found!) + remaining "er"
    → "er" → try "##er" (found!)
    Result: ["transform", "##er"]
    
    "COVID":
    → try "COVID" (not found)
    → "CO" (not found)
    → "C" (found) + "OVID"
    → continue: ["C", "##O", "##V", "##I", "##D"]
    
  DIFFERENCE FROM BPE:
    WordPiece uses greedy forward matching (fast inference)
    BPE uses merge rules in training order (sequence of merges)
    Both produce similar segmentations in practice.
```

---

## 7. SentencePiece

Google's SentencePiece library provides language-agnostic tokenization:

```
SENTENCEPIECE:

KEY DIFFERENCE: Treats the input as a raw character stream.
  Does NOT assume spaces = word boundaries.
  Uses a special underscore ▁ (U+2581) to mark word beginnings.
  
  "Hello World" → "▁Hello▁World"
  
  This makes it work for ALL languages:
    Chinese: 我爱自然语言处理 (no spaces)
    Japanese: 私は機械学習が好きです (no spaces)
    Arabic: مرحبا (right-to-left)
    English: "Hello World" (spaces)
  
  All treated identically — just sequences of Unicode characters.

SENTENCEPIECE MODELS:
  1. BPE (default): same BPE algorithm, but on the raw char stream
  2. Unigram Language Model: alternative algorithm
     - Start with large vocabulary (all substrings up to length L)
     - Iteratively REMOVE tokens whose removal minimally hurts perplexity
     - More theoretically principled than BPE
     - T5, ALBERT, XLNet use Unigram

TOKENIZATION EXAMPLE:
  Input: "Hello World"
  (with space represented as ▁)
  
  After SentencePiece:
  ["▁Hello", "▁World"]
  
  More complex: "Unhappiness is real"
  ["▁Un", "happiness", "▁is", "▁real"]
  (▁ marks word start, continuation has no ▁)
  
MODELS USING SENTENCEPIECE:
  LLaMA 1/2:    SentencePiece BPE, 32k vocab
  T5:           SentencePiece Unigram, 32k vocab
  ALBERT:       SentencePiece Unigram, 30k vocab
  mT5:          SentencePiece Unigram, 250k vocab (multilingual)
  Gemma:        SentencePiece BPE, 256k vocab

USING SENTENCEPIECE:
  import sentencepiece as spm
  
  # Train
  spm.SentencePieceTrainer.train(
      input='corpus.txt',
      model_prefix='my_tokenizer',
      vocab_size=32000,
      model_type='bpe',  # or 'unigram'
      character_coverage=0.9995,  # coverage of training chars
      pad_id=0, unk_id=1, bos_id=2, eos_id=3,  # special tokens
  )
  
  # Load and use
  sp = spm.SentencePieceProcessor()
  sp.load('my_tokenizer.model')
  
  ids = sp.encode("Hello World")
  print(ids)  # [e.g., 4312, 1876]
  
  text = sp.decode(ids)
  print(text)  # "Hello World"
```

---

## 8. tiktoken (OpenAI)

OpenAI's fast tokenizer library, used for GPT-3.5 and GPT-4:

```
tiktoken:
  Language: Rust (with Python bindings) → very fast
  Algorithm: byte-level BPE
  Available tokenizers:
    "gpt2":           50,257 tokens (GPT-2, GPT-3)
    "p50k_base":      50,281 tokens (Codex)
    "cl100k_base":    100,277 tokens (GPT-3.5, GPT-4, text-embedding-ada-002)
    "o200k_base":     200,019 tokens (GPT-4o)

INSTALLATION AND BASIC USE:
  pip install tiktoken

EXAMPLES:
  import tiktoken
  
  # Load tokenizer
  enc = tiktoken.get_encoding("cl100k_base")   # GPT-4 tokenizer
  
  # Encode text to token IDs
  text = "Hello, world! This is GPT-4's tokenizer."
  tokens = enc.encode(text)
  print(f"Text: {text}")
  print(f"Token IDs: {tokens}")
  print(f"Number of tokens: {len(tokens)}")
  
  # Decode back to text
  decoded = enc.decode(tokens)
  print(f"Decoded: {decoded}")
  
  # Show individual tokens
  for token_id in tokens:
      token_bytes = enc.decode_single_token_bytes(token_id)
      token_text = token_bytes.decode('utf-8', errors='replace')
      print(f"  ID {token_id:6d} → {repr(token_text)}")

EXPECTED OUTPUT:
  Text: Hello, world! This is GPT-4's tokenizer.
  Token IDs: [9906, 11, 1917, 0, 1115, 374, 480, 2898, 12, 19, 596, 47058, 13]
  Number of tokens: 13
  
  ID   9906 → 'Hello'
  ID     11 → ','
  ID   1917 → ' world'
  ID      0 → '!'
  ID   1115 → ' This'
  ID    374 → ' is'
  ID    480 → ' G'
  ID   2898 → 'PT'
  ID     12 → '-'
  ID     19 → '4'
  ID    596 → "'s"
  ID  47058 → ' tokenizer'
  ID     13 → '.'

Note: 'G' and 'PT' are separate tokens! "GPT" is split.
Note: ' world' (with space) is ONE token (space is part of the token).

COUNTING TOKENS FOR API PRICING:
  OpenAI API: charged per token
  1000 tokens ≈ 750 words ≈ 1.5 pages of text
  
  def count_tokens(text: str, model: str = "gpt-4") -> int:
      """Count tokens before sending to API — avoid surprises!"""
      encoding_name = {
          "gpt-4": "cl100k_base",
          "gpt-3.5-turbo": "cl100k_base",
          "gpt-4o": "o200k_base",
          "gpt-2": "gpt2",
      }.get(model, "cl100k_base")
      enc = tiktoken.get_encoding(encoding_name)
      return len(enc.encode(text))
```

---

## 9. Special Tokens

Every tokenizer has special tokens with specific roles:

```
SPECIAL TOKENS REFERENCE TABLE:
═══════════════════════════════════════════════════════════════════════
Token      │ Used By   │ ID  │ Purpose
───────────────────────────────────────────────────────────────────────
[PAD]      │ BERT      │   0 │ Padding to equalize batch lengths
[UNK]      │ BERT      │ 100 │ Unknown token (rare with BPE)
[CLS]      │ BERT      │ 101 │ Classification token (whole-sentence rep)
[SEP]      │ BERT      │ 102 │ Separator between segments
[MASK]     │ BERT      │ 103 │ Masked token (MLM pre-training only)
<unk>      │ LLaMA 1   │   0 │ Unknown (rare)
<s>        │ LLaMA 1/2 │   1 │ Beginning of sequence
</s>       │ LLaMA 1/2 │   2 │ End of sequence
<pad>      │ LLaMA 2   │     │ Padding
<|endoftext|>│GPT-2   │50256│ End of document (sequence boundary)
<|fim_prefix|>│Code   │     │ Fill-In-Middle: prefix part
<|fim_middle|>│Code   │     │ Fill-In-Middle: middle (to fill)
<|fim_suffix|>│Code   │     │ Fill-In-Middle: suffix part
═══════════════════════════════════════════════════════════════════════

INSTRUCTION/CHAT TOKENS (varies by model):

LLaMA 2 Chat format:
  [INST] {user message} [/INST]  {assistant response} </s>
  <s>[INST] {user} [/INST]{assistant}</s><s>[INST]{user2}[/INST]{assistant2}</s>

Mistral/Mixtral format:
  <s>[INST] {system} [/INST]</s>
  [INST] {user} [/INST] {assistant} </s>

ChatML format (used by many open models, Llama-3):
  <|im_start|>system
  You are a helpful assistant.
  <|im_end|>
  <|im_start|>user
  Hello!
  <|im_end|>
  <|im_start|>assistant
  Hi there!
  <|im_end|>

LLaMA-3 format:
  <|begin_of_text|>
  <|start_header_id|>system<|end_header_id|>
  You are a helpful assistant.<|eot_id|>
  <|start_header_id|>user<|end_header_id|>
  Hello!<|eot_id|>
  <|start_header_id|>assistant<|end_header_id|>
  Hi there!<|eot_id|>

WHY CHAT TEMPLATES MATTER:
  If you use a Mistral model but format prompts as LLaMA:
  The model still "works" but quality drops significantly.
  The model was trained to expect SPECIFIC tokens in SPECIFIC positions.
  
  Using HuggingFace's tokenizer.apply_chat_template() handles this automatically.
```

---

## 10. Training Your Own Tokenizer

Why train a custom tokenizer?

```
WHEN TO TRAIN A CUSTOM TOKENIZER:
  1. Domain-specific vocabulary (medical, legal, code)
     "electroencephalography" → medical corpus: 1 token
                              → general GPT-2:   8 tokens
  2. New language not well covered (e.g., adding Swahili, Yoruba)
  3. Custom special tokens (proprietary formats)
  4. Efficiency: if your data is very different from training corpus
```

### Full Code: Train Custom Tokenizer

```python
"""
Train a BPE tokenizer from scratch using HuggingFace tokenizers library.
Then use it to encode/decode text.
"""

from tokenizers import Tokenizer, models, trainers, pre_tokenizers, decoders
from tokenizers.processors import TemplateProcessing
import json
import os
from typing import List, Iterator


# ── Step 1: Prepare training data ─────────────────────────────────────────────

def get_training_corpus(texts: List[str], batch_size: int = 1000) -> Iterator[List[str]]:
    """Yield batches of text for tokenizer training."""
    for i in range(0, len(texts), batch_size):
        yield texts[i:i + batch_size]


# Sample corpus for demonstration
SAMPLE_CORPUS = [
    "The quick brown fox jumps over the lazy dog.",
    "Machine learning is a subset of artificial intelligence.",
    "Transformers use attention mechanisms to process sequences.",
    "Natural language processing enables computers to understand text.",
    "Neural networks learn representations from data through training.",
    "The attention mechanism allows the model to focus on relevant parts.",
    "Deep learning models require large amounts of training data.",
    "Python is a popular programming language for machine learning.",
    "The encoder processes the input and the decoder generates output.",
    "Language models learn to predict the next word in a sequence.",
    "Tokenization converts text into numerical representations for models.",
    "The vocabulary size affects both model quality and computational cost.",
    "Pre-training on large corpora enables few-shot learning capabilities.",
    "Fine-tuning adapts a pre-trained model to specific downstream tasks.",
    "The loss function measures how well the model predicts target values.",
] * 100  # repeat to have more data for BPE learning


# ── Step 2: Train BPE Tokenizer ───────────────────────────────────────────────

def train_bpe_tokenizer(
    texts: List[str],
    vocab_size: int = 1000,
    min_frequency: int = 2,
    special_tokens: List[str] = None,
    save_path: str = None,
) -> Tokenizer:
    """
    Train a Byte-level BPE tokenizer on the provided texts.
    
    Args:
        texts:          List of training strings
        vocab_size:     Target vocabulary size (including special tokens)
        min_frequency:  Minimum pair frequency to merge
        special_tokens: Special tokens to add to vocabulary
        save_path:      Path to save tokenizer.json (optional)
    
    Returns:
        Trained HuggingFace Tokenizer object
    """
    if special_tokens is None:
        special_tokens = ["[UNK]", "[PAD]", "[BOS]", "[EOS]", "[MASK]"]
    
    # ── Initialize BPE model ──────────────────────────────────────────────────
    tokenizer = Tokenizer(models.BPE(
        unk_token="[UNK]",   # token to use for unknowns
    ))
    
    # ── Pre-tokenizer: how to split into initial tokens ───────────────────────
    # ByteLevel: split on bytes, handles all Unicode
    # GPT-2 compatible: space becomes Ġ prefix
    tokenizer.pre_tokenizer = pre_tokenizers.ByteLevel(add_prefix_space=False)
    
    # ── Trainer: BPE algorithm configuration ─────────────────────────────────
    trainer = trainers.BpeTrainer(
        vocab_size=vocab_size,
        min_frequency=min_frequency,
        special_tokens=special_tokens,
        show_progress=True,
    )
    
    # ── Train on corpus ───────────────────────────────────────────────────────
    print(f"Training BPE tokenizer on {len(texts)} texts...")
    tokenizer.train_from_iterator(
        get_training_corpus(texts),
        trainer=trainer,
    )
    
    # ── Post-processor and decoder ────────────────────────────────────────────
    tokenizer.post_processor = TemplateProcessing(
        single="[BOS] $A [EOS]",
        special_tokens=[
            ("[BOS]", tokenizer.token_to_id("[BOS]")),
            ("[EOS]", tokenizer.token_to_id("[EOS]")),
        ],
    )
    
    # Decoder: handles the byte-level encoding during decoding
    tokenizer.decoder = decoders.ByteLevel()
    
    # ── Save ──────────────────────────────────────────────────────────────────
    if save_path:
        tokenizer.save(save_path)
        print(f"Tokenizer saved to {save_path}")
    
    return tokenizer


# ── Step 3: Inspect and Use Tokenizer ─────────────────────────────────────────

def inspect_tokenizer(tokenizer: Tokenizer):
    """Show tokenizer properties and example encodings."""
    
    print("\n" + "=" * 60)
    print("TOKENIZER INSPECTION")
    print("=" * 60)
    
    vocab_size = tokenizer.get_vocab_size()
    print(f"\nVocabulary size: {vocab_size}")
    
    # Show some vocabulary entries
    vocab = tokenizer.get_vocab()
    print(f"\nSample vocabulary entries:")
    
    sorted_vocab = sorted(vocab.items(), key=lambda x: x[1])
    special_tokens = ["[UNK]", "[PAD]", "[BOS]", "[EOS]", "[MASK]"]
    
    # Print special tokens
    for token, idx in sorted_vocab[:len(special_tokens)]:
        print(f"  ID {idx:5d}: {repr(token)}")
    
    print("  ...")
    
    # Print some mid-vocab tokens
    mid_start = vocab_size // 4
    for token, idx in sorted_vocab[mid_start:mid_start+10]:
        print(f"  ID {idx:5d}: {repr(token)}")
    
    print("  ...")
    
    # Encoding examples
    print("\n" + "=" * 60)
    print("ENCODING EXAMPLES")
    print("=" * 60)
    
    test_texts = [
        "Hello, world!",
        "machine learning",
        "Tokenization is important",
        "transformers",
        "12345",
        "COVID-19",
        "😀",  # emoji
    ]
    
    for text in test_texts:
        encoding = tokenizer.encode(text)
        tokens = encoding.tokens
        ids = encoding.ids
        
        # Decode back
        decoded = tokenizer.decode(ids, skip_special_tokens=True)
        
        print(f"\n  Input:   {repr(text)}")
        print(f"  Tokens:  {tokens}")
        print(f"  IDs:     {ids}")
        print(f"  Decoded: {repr(decoded)}")
        print(f"  Length:  {len(ids)} tokens")
        
        # Check reversibility
        if decoded.strip() == text.strip():
            print(f"  Reversible: ✓")
        else:
            print(f"  Reversible: ✗ (got '{decoded}')")


# ── Step 4: Compare Tokenizers ────────────────────────────────────────────────

def compare_tokenizers_on_text(text: str):
    """
    Compare how different tokenizers handle the same text.
    Shows the difference in token count and vocabulary usage.
    """
    print("\n" + "=" * 60)
    print(f"TOKENIZER COMPARISON")
    print(f"Text: {repr(text)}")
    print("=" * 60)
    
    # 1. Character-level (manual)
    char_tokens = list(text)
    print(f"\nCharacter-level:")
    print(f"  Tokens: {char_tokens[:15]}{'...' if len(char_tokens) > 15 else ''}")
    print(f"  Count:  {len(char_tokens)}")
    
    # 2. Word-level (manual split)
    import re
    word_tokens = re.findall(r'\w+|[^\w\s]', text)
    print(f"\nWord-level:")
    print(f"  Tokens: {word_tokens}")
    print(f"  Count:  {len(word_tokens)}")
    
    # 3. BERT WordPiece
    try:
        from transformers import BertTokenizer
        bert_tok = BertTokenizer.from_pretrained("bert-base-uncased")
        bert_tokens = bert_tok.tokenize(text.lower())
        print(f"\nBERT WordPiece (uncased):")
        print(f"  Tokens: {bert_tokens}")
        print(f"  Count:  {len(bert_tokens)}")
    except:
        print("\nBERT tokenizer not available (pip install transformers)")
    
    # 4. GPT-2 / tiktoken
    try:
        import tiktoken
        enc = tiktoken.get_encoding("gpt2")
        gpt2_ids = enc.encode(text)
        gpt2_tokens = [enc.decode([t]) for t in gpt2_ids]
        print(f"\nGPT-2 BPE (tiktoken):")
        print(f"  Tokens: {[repr(t) for t in gpt2_tokens]}")
        print(f"  Count:  {len(gpt2_tokens)}")
    except:
        print("\ntiktoken not available (pip install tiktoken)")


# ── Step 5: Demonstrate Gotchas ───────────────────────────────────────────────

def demonstrate_tokenization_gotchas():
    """Show common tokenization surprises relevant to API users."""
    
    print("\n" + "=" * 60)
    print("TOKENIZATION GOTCHAS")
    print("=" * 60)
    
    try:
        import tiktoken
        enc = tiktoken.get_encoding("cl100k_base")  # GPT-4 tokenizer
        
        def show_tokens(text):
            ids = enc.encode(text)
            tokens = [enc.decode([t]) for t in ids]
            return tokens, len(ids)
        
        # Gotcha 1: Numbers
        print("\n1. NUMBERS — not always 1 token per digit:")
        for num in ["1", "100", "1000", "10000", "100000", "99", "2024"]:
            tokens, count = show_tokens(num)
            print(f"  '{num}' → {tokens} ({count} tokens)")
        
        # Gotcha 2: Case sensitivity
        print("\n2. CASE SENSITIVITY:")
        for word in ["hello", "Hello", "HELLO", "hElLo"]:
            tokens, count = show_tokens(word)
            print(f"  '{word}' → {tokens} ({count} tokens)")
        
        # Gotcha 3: Spaces matter
        print("\n3. LEADING SPACES (same word, different token):")
        for word in ["cat", " cat", "Cat", " Cat"]:
            tokens, count = show_tokens(word)
            ids = enc.encode(word)
            print(f"  {repr(word):10} → IDs: {ids} tokens: {tokens}")
        
        # Gotcha 4: Code and indentation
        print("\n4. CODE TOKENIZATION:")
        codes = [
            "def foo():",
            "    return 42",
            "    return 42    ",  # trailing spaces
            "\t\treturn 42",    # tabs
        ]
        for code in codes:
            tokens, count = show_tokens(code)
            print(f"  {repr(code):25} → {count} tokens: {tokens}")
        
        # Gotcha 5: Multilingual efficiency
        print("\n5. MULTILINGUAL EFFICIENCY (tokens per word concept):")
        texts = {
            "English": "The cat sat on the mat",
            "French":  "Le chat s'est assis sur le tapis",
            "Japanese": "猫がマットの上に座った",
            "Arabic":  "جلس القط على الحصيرة",
            "Code":    "def reverse(s): return s[::-1]",
        }
        for lang, text in texts.items():
            tokens, count = show_tokens(text)
            chars = len(text)
            print(f"  {lang:10}: {count:3d} tokens for '{text[:30]}...' "
                  f"({chars/count:.1f} chars/token)")
        
        # Gotcha 6: GPT vs BERT tokenizer differences
        print("\n6. TOKEN COUNT FOR API COST ESTIMATION:")
        long_text = """
        The transformer architecture, introduced in the landmark 2017 paper 
        "Attention Is All You Need" by Vaswani et al., revolutionized natural 
        language processing by replacing recurrent neural networks with 
        self-attention mechanisms. This enabled parallel training and better 
        modeling of long-range dependencies.
        """
        tokens_gpt4, count_gpt4 = show_tokens(long_text)
        words = len(long_text.split())
        
        print(f"  Text length: {len(long_text)} characters, ~{words} words")
        print(f"  GPT-4 tokens: {count_gpt4}")
        print(f"  Tokens/word: {count_gpt4/words:.2f}")
        print(f"  At $0.03/1k tokens (GPT-4): ${count_gpt4 * 0.03 / 1000:.4f}")
        
    except ImportError:
        print("Install tiktoken: pip install tiktoken")


# ── Main ─────────────────────────────────────────────────────────────────────

if __name__ == "__main__":
    # Train tokenizer
    tokenizer = train_bpe_tokenizer(
        texts=SAMPLE_CORPUS,
        vocab_size=500,   # small vocab for illustration
        min_frequency=2,
    )
    
    # Inspect it
    inspect_tokenizer(tokenizer)
    
    # Compare tokenizers
    compare_tokenizers_on_text("The transformer architecture changed everything")
    compare_tokenizers_on_text("unhappiness electroencephalography")
    
    # Show gotchas
    demonstrate_tokenization_gotchas()
```

---

## 11. Tokenization Quirks and Gotchas

Every developer using LLM APIs should know these:

```
GOTCHA 1: TOKEN BOUNDARIES ARE NOT WORD BOUNDARIES
  "GPT-4" in cl100k_base:
    ['G', 'PT', '-', '4']  — 4 tokens, not 1!
    ← 'GPT' happened to not be in the merge rules this way
    
  Impact: asking the model to "count letters in 'GPT'" is hard
  because it sees ['G','PT'] and 'PT' is one token — it can't easily
  see that 'P' and 'T' are separate letters.

GOTCHA 2: SPACES ARE PART OF TOKENS
  In GPT-2/GPT-4 (byte-level BPE):
    "cat" → ['cat']      (ID: 8415)
    " cat" → [' cat']    (ID: 3797) ← DIFFERENT TOKEN!
    
  This matters for text generation:
    After a comma, you expect " The" not "The" (different token!)
    Tokenizer handles this automatically, but it explains why
    "token probabilities" can look strange.

GOTCHA 3: NUMBERS ARE TOKENIZED INCONSISTENTLY
  In GPT-2 (cl100k_base):
    "1"     → ['1']        — 1 token
    "12"    → ['12']       — 1 token
    "123"   → ['123']      — 1 token
    "1234"  → ['1234']     — 1 token
    "12345" → ['123', '45'] — 2 tokens!
    
  This is why LLMs are bad at arithmetic:
    To add "12345 + 67890", the model must reason about TOKENS
    not digits. "12345" = ['123', '45'] — where does position 10000 go?
    
  GPT-4 improved significantly with TOOL USE (code interpreter).

GOTCHA 4: TOKEN COUNT ≠ WORD COUNT
  1 English word ≈ 1.3-1.5 tokens (average)
  Common words: usually 1 token
  Technical/rare words: 2-6 tokens
  
  Context window "128k tokens" ≠ 128k words
  At 1.3 tokens/word: 128k tokens ≈ 98k words ≈ 196 pages

GOTCHA 5: MULTILINGUAL INEFFICIENCY
  English: ~4-5 chars per token
  Japanese: ~1-2 chars per token  ← same context window = much less text!
  Arabic:   ~1-2 chars per token
  Code:     ~4-5 chars per token (similar to English)
  
  A GPT-4 call with the same text in Japanese costs 2-3× as much
  (in tokens) as the same content in English!
  
  This is an active area: multilingual models are trained with 
  larger vocabularies to improve non-English efficiency.

GOTCHA 6: CHAT TEMPLATES USE TOKENS TOO
  Every "<|system|>", "<|user|>", "<|assistant|>" = extra tokens
  For a system prompt of 500 tokens, these overhead tokens matter.
  
  Always count: system_tokens + history_tokens + user_tokens < max_context

GOTCHA 7: THE TOKENIZER MUST MATCH THE MODEL
  Using BERT tokenizer with GPT-2 model → nonsense
  Using GPT-2 tokenizer (50k vocab) with LLaMA (32k vocab) → index errors
  
  Always load the tokenizer that matches the model:
    tokenizer = AutoTokenizer.from_pretrained("mistralai/Mistral-7B-v0.1")
    # NOT: tokenizer = GPT2Tokenizer.from_pretrained("gpt2")

GOTCHA 8: PADDING AND MASKING
  When batching sequences of different lengths, shorter ones get padded.
  Padding tokens MUST be masked: attention_mask=0 for pad tokens.
  Forgetting the mask → model attends to [PAD] tokens → wrong results.
```

---

## 12. Full Code: Build BPE Tokenizer from Scratch

For deep understanding, here's BPE implemented from scratch in pure Python:

```python
"""
BPE (Byte Pair Encoding) Tokenizer from Scratch.
Pure Python implementation for educational understanding.
No external libraries required (except standard library).
"""

from collections import defaultdict
from typing import Dict, List, Tuple, Optional
import re


class BPETokenizer:
    """
    From-scratch BPE tokenizer implementation.
    
    Follows the algorithm from Sennrich et al. (2016):
    "Neural Machine Translation of Rare Words with Subword Units"
    """
    
    def __init__(self):
        self.vocab: Dict[str, int] = {}      # token → ID
        self.merges: List[Tuple[str,str]] = []  # ordered list of BPE merges
        self.special_tokens = {}
        
    def _get_word_counts(self, corpus: List[str]) -> Dict[str, int]:
        """
        Count word frequencies in corpus.
        Each word is represented as a space-separated character sequence.
        Add end-of-word marker </w>.
        """
        word_counts = defaultdict(int)
        for text in corpus:
            for word in text.lower().split():
                # Represent word as space-separated characters + </w>
                word_repr = ' '.join(list(word)) + ' </w>'
                word_counts[word_repr] += 1
        return dict(word_counts)
    
    def _get_pair_counts(
        self,
        word_counts: Dict[str, int],
    ) -> Dict[Tuple[str,str], int]:
        """Count all consecutive token pairs in current word representations."""
        pairs = defaultdict(int)
        for word, count in word_counts.items():
            symbols = word.split()
            for i in range(len(symbols) - 1):
                pairs[(symbols[i], symbols[i+1])] += count
        return dict(pairs)
    
    def _merge_pair(
        self,
        word_counts: Dict[str, int],
        pair: Tuple[str, str],
    ) -> Dict[str, int]:
        """
        Merge all occurrences of (a, b) into 'ab' in word representations.
        """
        new_word_counts = {}
        bigram = ' '.join(pair)
        replacement = ''.join(pair)
        
        for word, count in word_counts.items():
            # Replace "a b" with "ab" in the word representation
            new_word = word.replace(bigram, replacement)
            new_word_counts[new_word] = count
        
        return new_word_counts
    
    def train(
        self,
        corpus: List[str],
        vocab_size: int = 500,
        verbose: bool = True,
    ):
        """
        Train the BPE tokenizer on the given corpus.
        
        Args:
            corpus:     List of training texts
            vocab_size: Target vocabulary size
            verbose:    Print training progress
        """
        # Step 1: Build initial character vocabulary
        word_counts = self._get_word_counts(corpus)
        
        # Count initial characters
        char_counts = defaultdict(int)
        for word, count in word_counts.items():
            for char in word.split():
                char_counts[char] += count
        
        # Initialize vocabulary with all characters
        self.vocab = {char: idx for idx, char in enumerate(sorted(char_counts.keys()))}
        
        if verbose:
            print(f"Initial vocab size: {len(self.vocab)} characters")
        
        self.merges = []
        
        # Step 2: Merge pairs until vocab_size reached
        num_merges = vocab_size - len(self.vocab)
        
        for step in range(num_merges):
            # Count all pairs
            pair_counts = self._get_pair_counts(word_counts)
            
            if not pair_counts:
                break
            
            # Find most frequent pair
            best_pair = max(pair_counts, key=pair_counts.get)
            best_count = pair_counts[best_pair]
            
            if best_count < 2:  # stop if no pair appears more than once
                break
            
            # Merge the best pair
            word_counts = self._merge_pair(word_counts, best_pair)
            
            # Add new merged token to vocabulary
            new_token = ''.join(best_pair)
            self.vocab[new_token] = len(self.vocab)
            self.merges.append(best_pair)
            
            if verbose and step % 50 == 0:
                print(f"  Step {step:4d}: Merged {repr(best_pair)} → {repr(new_token)} "
                      f"(freq={best_count}), vocab_size={len(self.vocab)}")
        
        if verbose:
            print(f"\nTraining complete!")
            print(f"  Final vocab size: {len(self.vocab)}")
            print(f"  Total merges: {len(self.merges)}")
    
    def _tokenize_word(self, word: str) -> List[str]:
        """
        Tokenize a single word by applying BPE merges.
        Apply merges in the ORDER they were learned (crucial!).
        """
        # Start with character representation
        symbols = list(word) + ['</w>']
        
        # Apply merges in order
        for (first, second) in self.merges:
            i = 0
            new_symbols = []
            while i < len(symbols):
                if (i < len(symbols) - 1 and 
                    symbols[i] == first and 
                    symbols[i+1] == second):
                    new_symbols.append(first + second)
                    i += 2  # skip both symbols
                else:
                    new_symbols.append(symbols[i])
                    i += 1
            symbols = new_symbols
        
        return symbols
    
    def encode(self, text: str) -> List[int]:
        """
        Encode text to list of token IDs.
        
        Args:
            text: Input text
            
        Returns:
            List of integer token IDs
        """
        tokens = []
        for word in text.lower().split():
            word_tokens = self._tokenize_word(word)
            for token in word_tokens:
                if token in self.vocab:
                    tokens.append(self.vocab[token])
                else:
                    # Fall back to characters for unknown tokens
                    for char in token:
                        if char in self.vocab:
                            tokens.append(self.vocab[char])
                        else:
                            tokens.append(self.vocab.get('<unk>', 0))
        return tokens
    
    def decode(self, ids: List[int]) -> str:
        """
        Decode list of token IDs back to text.
        """
        id_to_token = {v: k for k, v in self.vocab.items()}
        tokens = [id_to_token.get(id_, '<unk>') for id_ in ids]
        
        # Join tokens and remove end-of-word markers
        text = ' '.join(tokens)
        text = text.replace('</w>', '')
        text = text.replace('  ', ' ').strip()
        return text
    
    def tokenize(self, text: str) -> List[str]:
        """Tokenize text, returning token strings (not IDs)."""
        tokens = []
        for word in text.lower().split():
            tokens.extend(self._tokenize_word(word))
        return tokens
    
    def get_vocab_size(self) -> int:
        return len(self.vocab)
    
    def print_vocab_sample(self, n: int = 20):
        """Print a sample of the vocabulary."""
        print(f"\nVocabulary sample (first {n} entries):")
        sorted_vocab = sorted(self.vocab.items(), key=lambda x: x[1])
        for token, idx in sorted_vocab[:n]:
            print(f"  ID {idx:4d}: {repr(token)}")


# ── Demo ─────────────────────────────────────────────────────────────────────

def run_demo():
    """Demonstrate the BPE tokenizer with a small corpus."""
    
    corpus = [
        "low lower lowest",
        "new newer newest",
        "wide wider widest",
        "fast faster fastest",
        "machine learning is great",
        "learning machine learning",
        "great machine learning models",
        "natural language processing",
        "language models process language",
        "attention mechanism in transformers",
        "transformers use self-attention layers",
        "the transformer model processes sequences",
    ] * 10  # repeat for more frequency data
    
    print("=" * 60)
    print("BPE TOKENIZER FROM SCRATCH")
    print("=" * 60)
    print(f"\nTraining on {len(corpus)} sentences...")
    
    tokenizer = BPETokenizer()
    tokenizer.train(corpus, vocab_size=100, verbose=True)
    
    tokenizer.print_vocab_sample(30)
    
    # Test encoding
    test_texts = [
        "low",
        "lowest",
        "newest",
        "machine learning",
        "natural language processing",
        "transformers",
        "unknown_word",  # novel word
    ]
    
    print("\n" + "=" * 60)
    print("TOKENIZATION EXAMPLES")
    print("=" * 60)
    
    for text in test_texts:
        tokens = tokenizer.tokenize(text)
        ids = tokenizer.encode(text)
        print(f"\n  Input:   {repr(text)}")
        print(f"  Tokens:  {tokens}")
        print(f"  IDs:     {ids}")
        print(f"  Length:  {len(ids)} tokens")
    
    # Show some BPE merges learned
    print("\n" + "=" * 60)
    print(f"FIRST 20 BPE MERGES LEARNED")
    print("=" * 60)
    for i, (a, b) in enumerate(tokenizer.merges[:20]):
        print(f"  Merge {i+1:3d}: '{a}' + '{b}' → '{a+b}'")


if __name__ == "__main__":
    run_demo()
```

---

## 13. Summary

```
TOKENIZATION METHODS COMPARISON:
═══════════════════════════════════════════════════════════════════
Method       │ Vocab   │ Seq Length │ OOV │ Languages │ Used By
─────────────────────────────────────────────────────────────────
Character    │ ~100    │ Very long  │ No  │ Any       │ Char-RNN
Word-level   │ 50k-1M  │ Short      │ Yes │ Latin     │ Old NLP
BPE          │ 32k-100k│ Medium     │ No  │ Any       │ GPT-*, LLaMA
WordPiece    │ 30k     │ Medium     │ No  │ Any (##)  │ BERT
SentencePiece│ 32k-250k│ Medium     │ No  │ Any (▁)   │ T5, LLaMA
═══════════════════════════════════════════════════════════════════

KEY FACTS FOR DEVELOPERS:
  - 1 English word ≈ 1.3-1.5 tokens
  - Context "128k tokens" ≈ 98k English words
  - Japanese/Chinese: ~2-4× more tokens per "word" concept
  - Numbers: 1-5 digits usually 1 token, longer numbers split
  - Always use model's own tokenizer
  - Check token count before API calls to estimate cost
  - Chat templates (system/user/assistant markers) cost tokens too
```

---

## Mini Projects

### Mini Project 1: BPE Tokenizer from Scratch

Implement Byte Pair Encoding — the algorithm behind GPT tokenization — and watch it merge character pairs.

**Objective:** Understand tokenization at the algorithm level, not just as a black box.

```python
import re
from collections import Counter, defaultdict
import matplotlib.pyplot as plt
import numpy as np

class BPETokenizer:
    def __init__(self):
        self.merges  = {}   # (a, b) → merged_token
        self.vocab   = {}   # token → id
        self.id2tok  = {}   # id → token

    def get_vocab(self, corpus):
        """Build initial word vocabulary with character-level splits."""
        vocab = Counter()
        for word in corpus.lower().split():
            word = ' '.join(list(word)) + ' </w>'  # mark end of word
            vocab[word] += 1
        return vocab

    def get_stats(self, vocab):
        """Count all adjacent pair frequencies."""
        pairs = Counter()
        for word, freq in vocab.items():
            tokens = word.split()
            for i in range(len(tokens)-1):
                pairs[(tokens[i], tokens[i+1])] += freq
        return pairs

    def merge_vocab(self, pair, vocab):
        """Merge the best pair in all words."""
        bigram  = re.escape(' '.join(pair))
        pattern = re.compile(r'(?<!\S)' + bigram + r'(?!\S)')
        return {pattern.sub(''.join(pair), word): freq for word, freq in vocab.items()}

    def train(self, corpus, n_merges=50, verbose=True):
        vocab = self.get_vocab(corpus)
        merge_history = []

        for step in range(n_merges):
            pairs = self.get_stats(vocab)
            if not pairs: break
            best_pair = max(pairs, key=pairs.get)
            best_freq = pairs[best_pair]
            vocab = self.merge_vocab(best_pair, vocab)
            merged = ''.join(best_pair)
            self.merges[best_pair] = merged
            merge_history.append((step+1, best_pair, best_freq, merged))
            if verbose and (step+1) % 10 == 0:
                print(f"  Step {step+1:3d}: merge {best_pair} → '{merged}' (freq={best_freq})")

        # Build vocabulary
        all_tokens = set()
        for word in vocab:
            all_tokens.update(word.split())
        self.vocab   = {tok: i for i, tok in enumerate(sorted(all_tokens))}
        self.id2tok  = {i: tok for tok, i in self.vocab.items()}
        return merge_history

    def tokenize(self, text):
        tokens = []
        for word in text.lower().split():
            word_tokens = list(word) + ['</w>']
            # Apply learned merges in order
            changed = True
            while changed:
                changed = False
                i = 0
                new_tokens = []
                while i < len(word_tokens):
                    if i < len(word_tokens)-1:
                        pair = (word_tokens[i], word_tokens[i+1])
                        if pair in self.merges:
                            new_tokens.append(self.merges[pair])
                            i += 2; changed = True; continue
                    new_tokens.append(word_tokens[i])
                    i += 1
                word_tokens = new_tokens
            tokens.extend(word_tokens)
        return tokens

    def encode(self, text):
        return [self.vocab.get(t, 0) for t in self.tokenize(text)]

    def decode(self, ids):
        return ''.join(self.id2tok.get(i, '?') for i in ids).replace('</w>', ' ').strip()

# Training corpus
corpus = """
the quick brown fox jumps over the lazy dog
the dog barked at the fox near the brown tree
a quick brown hare and a lazy brown bear
the fox was quick and the dog was slow
quick foxes and lazy dogs run in the morning
the cat sat on the mat and the fox sat on the log
the man ran past the fox and the dog chased the cat
""" * 5  # repeat to get better stats

print("=== Training BPE Tokenizer ===")
tokenizer = BPETokenizer()
merge_history = tokenizer.train(corpus, n_merges=60, verbose=True)

# Vocabulary analysis
print(f"\nVocabulary size: {len(tokenizer.vocab)}")
print(f"Number of merges: {len(tokenizer.merges)}")

# Test tokenization
test_sentences = [
    "the quick brown fox",
    "lazy dog jumps",
    "quickbrown",  # OOV-like test
]
print("\n=== Tokenization Examples ===")
for sent in test_sentences:
    toks = tokenizer.tokenize(sent)
    print(f"  '{sent}' → {toks} ({len(toks)} tokens)")

# Visualization
fig, axes = plt.subplots(2, 2, figsize=(15, 10))
fig.suptitle("BPE Tokenization Analysis", fontsize=14, fontweight='bold')

# Merge frequency over steps
freqs = [h[2] for h in merge_history]
axes[0, 0].plot(range(1, len(freqs)+1), freqs, 'b-o', markersize=3)
axes[0, 0].set_title("Merge Frequency per Step\n(decreasing = merging rarer pairs later)")
axes[0, 0].set_xlabel("Merge Step"); axes[0, 0].set_ylabel("Pair Frequency"); axes[0, 0].grid(True, alpha=0.3)
axes[0, 0].set_yscale('log')

# Token length distribution
token_lengths = [len(tok.replace('</w>', '')) for tok in tokenizer.vocab]
axes[0, 1].hist(token_lengths, bins=range(1, max(token_lengths)+2), color='steelblue', alpha=0.8, edgecolor='white')
axes[0, 1].set_title("Vocabulary Token Length Distribution")
axes[0, 1].set_xlabel("Token Length (chars)"); axes[0, 1].set_ylabel("Count"); axes[0, 1].grid(True, alpha=0.3)

# Tokens per word as vocabulary grows
n_merges_range = [0, 10, 20, 30, 40, 50, 60]
test_words = corpus.split()[:50]
tokens_per_word_history = []
for n in n_merges_range:
    if n <= len(merge_history):
        small_tok = BPETokenizer()
        small_tok.train(corpus, n_merges=n, verbose=False)
        avg_toks = np.mean([len(small_tok.tokenize(w)) for w in test_words[:20]])
        tokens_per_word_history.append(avg_toks)
axes[1, 0].plot(n_merges_range[:len(tokens_per_word_history)], tokens_per_word_history, 'g-o', linewidth=2)
axes[1, 0].axhline(1.0, color='red', linestyle='--', label='Word-level (ideal)')
axes[1, 0].set_title("Avg Tokens per Word vs # Merges\n(more merges → fewer, longer tokens)")
axes[1, 0].set_xlabel("# Merges"); axes[1, 0].set_ylabel("Avg Tokens per Word"); axes[1, 0].legend(); axes[1, 0].grid(True, alpha=0.3)

# Top 20 most common tokens in vocabulary
all_word_tokens = []
for word in corpus.split():
    all_word_tokens.extend(tokenizer.tokenize(word))
token_freqs = Counter(all_word_tokens)
top_tokens  = token_freqs.most_common(20)
tok_names = [t[0].replace('</w>', '↵') for t in top_tokens]
tok_counts = [t[1] for t in top_tokens]
axes[1, 1].barh(range(len(tok_names)), tok_counts[::-1], color='coral', alpha=0.8)
axes[1, 1].set_yticks(range(len(tok_names)))
axes[1, 1].set_yticklabels(tok_names[::-1], fontsize=9)
axes[1, 1].set_title("Top 20 Most Frequent Tokens")
axes[1, 1].set_xlabel("Frequency"); axes[1, 1].grid(True, alpha=0.3, axis='x')

plt.tight_layout()
plt.savefig("bpe_tokenizer.png", dpi=150)
plt.show()
print("Saved: bpe_tokenizer.png")
```

---

### Mini Project 2: Tokenizer Comparison Across Languages and Domains

Compare how different tokenizers (character, word, BPE) handle text from different domains.

**Objective:** Understand why tokenization affects model performance for code, math, and multilingual text.

```python
import re
from collections import Counter
import matplotlib.pyplot as plt
import numpy as np

# Try to use tiktoken (OpenAI) and transformers tokenizers
# Fall back to manual implementations if not installed

test_texts = {
    "English prose":   "The quick brown fox jumps over the lazy dog and ran away into the forest.",
    "Python code":     "def fibonacci(n):\n    if n <= 1:\n        return n\n    return fibonacci(n-1) + fibonacci(n-2)",
    "Math formula":    "The quadratic formula x = (-b ± √(b²-4ac)) / 2a solves ax² + bx + c = 0",
    "Email/URLs":      "Contact us at info@example.com or visit https://www.example.com/api/v2/users?id=123",
    "Repeated chars":  "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa",
    "Numbers":         "3.14159265 2.71828182 1.41421356 1.61803398 2.30258509 6.28318530",
}

# Tokenizer implementations
def char_tokenize(text):
    return list(text)

def word_tokenize(text):
    return re.findall(r'\w+|[^\w\s]', text)

def simple_bpe_tokenize(text, n_merges=30):
    """Quick BPE approximation for comparison."""
    tok = BPETokenizer()
    tok.train(text * 3, n_merges=n_merges, verbose=False)
    return tok.tokenize(text)

# Collect statistics
results = {}
for text_name, text in test_texts.items():
    char_toks = char_tokenize(text)
    word_toks = word_tokenize(text)
    bpe_toks  = simple_bpe_tokenize(text)

    results[text_name] = {
        'text_len':    len(text),
        'char_tokens': len(char_toks),
        'word_tokens': len(word_toks),
        'bpe_tokens':  len(bpe_toks),
        'char_compression': len(text) / len(char_toks),
        'word_compression': len(text) / len(word_toks),
        'bpe_compression':  len(text) / len(bpe_toks),
    }

# Try tiktoken if available
try:
    import tiktoken
    enc = tiktoken.get_encoding("cl100k_base")  # GPT-4's tokenizer
    for text_name, text in test_texts.items():
        tik_toks = enc.encode(text)
        results[text_name]['tiktoken_tokens'] = len(tik_toks)
        results[text_name]['tiktoken_compression'] = len(text) / len(tik_toks)
    HAS_TIKTOKEN = True
    print("Using tiktoken (GPT-4 tokenizer)")
except ImportError:
    HAS_TIKTOKEN = False
    print("tiktoken not available — install with: pip install tiktoken")

# Visualization
fig, axes = plt.subplots(2, 2, figsize=(16, 11))
fig.suptitle("Tokenizer Comparison: Efficiency Across Text Types", fontsize=13, fontweight='bold')

names = list(results.keys())
n = len(names)
x = np.arange(n)
width = 0.25

char_counts = [results[n]['char_tokens'] for n in names]
word_counts = [results[n]['word_tokens'] for n in names]
bpe_counts  = [results[n]['bpe_tokens']  for n in names]

axes[0, 0].bar(x - width, char_counts, width, label='Char', color='red',      alpha=0.8)
axes[0, 0].bar(x,          word_counts, width, label='Word', color='steelblue', alpha=0.8)
axes[0, 0].bar(x + width,  bpe_counts,  width, label='BPE',  color='green',    alpha=0.8)
if HAS_TIKTOKEN:
    tik_counts = [results[n]['tiktoken_tokens'] for n in names]
    axes[0, 0].bar(x + 2*width, tik_counts, width, label='GPT-4 (tiktoken)', color='purple', alpha=0.8)
axes[0, 0].set_xticks(x); axes[0, 0].set_xticklabels([n[:12] for n in names], rotation=30, ha='right', fontsize=7)
axes[0, 0].set_ylabel("Token Count"); axes[0, 0].set_title("Token Count by Tokenizer"); axes[0, 0].legend(fontsize=8); axes[0, 0].grid(True, alpha=0.3, axis='y')

# Chars per token (compression ratio)
char_cpr = [results[n]['char_compression'] for n in names]
word_cpr = [results[n]['word_compression'] for n in names]
bpe_cpr  = [results[n]['bpe_compression']  for n in names]
axes[0, 1].plot(names, char_cpr, 'r-o', label='Char',  markersize=6, linewidth=2)
axes[0, 1].plot(names, word_cpr, 'b-o', label='Word',  markersize=6, linewidth=2)
axes[0, 1].plot(names, bpe_cpr,  'g-o', label='BPE',   markersize=6, linewidth=2)
if HAS_TIKTOKEN:
    tik_cpr = [results[n]['tiktoken_compression'] for n in names]
    axes[0, 1].plot(names, tik_cpr, 'p-o', label='tiktoken', markersize=6, linewidth=2)
axes[0, 1].set_xticklabels([n[:12] for n in names], rotation=30, ha='right', fontsize=7)
axes[0, 1].set_ylabel("Chars per Token"); axes[0, 1].set_title("Compression Ratio\n(higher = more efficient)")
axes[0, 1].legend(fontsize=8); axes[0, 1].grid(True, alpha=0.3)

# Show actual tokenizations for one example
example_text = "def fibonacci(n):\n    return n if n <= 1 else fibonacci(n-1)"
char_toks  = char_tokenize(example_text)
word_toks  = word_tokenize(example_text)
bpe_toks   = simple_bpe_tokenize(example_text)

tokenizations = [('Character', char_toks[:30]), ('Word', word_toks[:20]), ('BPE', bpe_toks[:20])]
for row_idx, (name, toks) in enumerate(tokenizations):
    ax = axes[1, row_idx % 2] if row_idx < 2 else axes[1, 1]
    if row_idx == 2:
        ax = axes[1, 1]
    ax.axis('off')
    tok_str = f"Tokenizer: {name}\nInput: '{example_text[:40]}...'\n\n"
    tok_str += f"Tokens ({len(toks)} shown):\n"
    for i, t in enumerate(toks[:15]):
        tok_str += f"  [{i:2d}] '{t}'\n"
    if len(toks) > 15:
        tok_str += f"  ... ({len(toks) - 15} more)\n"
    tok_str += f"\nTotal: {len(char_tokenize(example_text))} char | "
    tok_str += f"{len(word_tokenize(example_text))} word | {len(bpe_toks)} BPE"
    ax.text(0.02, 0.95, tok_str, transform=ax.transAxes, fontsize=7,
            va='top', fontfamily='monospace',
            bbox=dict(boxstyle='round', facecolor='lightyellow', alpha=0.8))
    ax.set_title(f"{name} Tokenization of Python Code")

plt.tight_layout()
plt.savefig("tokenizer_comparison.png", dpi=150)
plt.show()

print("\nTokenization Summary for Python Code:")
print(f"  Character tokens: {len(char_tokenize(example_text))}")
print(f"  Word tokens:      {len(word_tokenize(example_text))}")
print(f"  BPE tokens:       {len(simple_bpe_tokenize(example_text))}")
if HAS_TIKTOKEN:
    print(f"  GPT-4 tokens:     {len(enc.encode(example_text))}")
```

---

## 14. Exercises

**Exercise 1**: Extend the `BPETokenizer` class to handle byte-level encoding (start from bytes 0-255 instead of characters). Test on text containing emojis and Japanese characters. Verify no [UNK] tokens appear.

**Exercise 2**: Use `tiktoken` to find the 10 English words that tokenize into the MOST tokens in the cl100k_base vocabulary. Find the 10 words that tokenize into the FEWEST tokens. What patterns do you notice?

**Exercise 3**: Measure the "tokenization efficiency" of 5 different languages (pick 5 from: English, French, Spanish, German, Chinese, Japanese, Arabic, Hindi). Use a paragraph of similar information content in each language. How does token count per character differ?

**Exercise 4**: Implement a function that checks if two strings are "tokenizer-equivalent" — that is, they tokenize to the same sequence. Use this to find 5 examples where small text changes (capitalization, spacing) change the tokenization.

**Exercise 5**: Train a SentencePiece tokenizer on a Python code corpus (use The Stack or any Python files). Compare the tokenization of common Python patterns ("def ", "return ", "import ") vs the GPT-2 tokenizer. Which handles Python code more efficiently?

**Exercise 6**: The "tokenization instability" problem: some models behave differently when the same content appears at different positions in the context (because the BPE merges depend on context boundaries). Write a test to demonstrate this: show a token sequence where inserting a space before the text changes its tokenization and therefore the model's behavior.

---

**Chapter Summary**: Tokenization converts raw text to integer sequences. Character-level is simple but creates very long sequences; word-level is intuitive but can't handle unseen words. Subword tokenization (BPE, WordPiece, SentencePiece) is the sweet spot: common words become single tokens, rare words split into common subwords, and any text can be encoded. BPE starts with characters, iteratively merges the most frequent consecutive pairs until the target vocabulary size. Byte-level BPE (GPT-2, LLaMA) starts from raw bytes (256 byte values), eliminating the unknown token entirely. Special tokens ([CLS], [SEP], [MASK], chat markers) add structure for specific tasks. Key developer gotchas: numbers tokenize inconsistently, spaces are part of tokens, multilingual text uses more tokens than English, and context window limits are in tokens not words.

**What's Next →** [Chapter 28: Building TinyGPT from Scratch](./28-building-tinygpt-from-scratch.md)

*We now have all the pieces. Let's build a complete GPT-style language model from scratch — tokenizer, model, training loop, and generation. You'll run it and watch it learn Shakespeare.*
