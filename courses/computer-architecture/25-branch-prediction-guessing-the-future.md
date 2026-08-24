# Chapter 25: Branch Prediction — Guessing the Future

Every `if` statement, every loop, every function call that might return early — all of these generate branch instructions. A modern program executes a branch roughly every 5–7 instructions. If the processor had to stop and wait every time it encountered a branch, it would spend 15–20% of its time doing nothing. Branch prediction is the processor's way of making educated guesses about where execution is headed — and modern predictors are right more than 99% of the time.

## Table of Contents

1. [Why Branches Are So Expensive](#1-why-branches-are-expensive)
2. [Static Branch Prediction](#2-static-branch-prediction)
3. [Dynamic Prediction: 1-Bit and 2-Bit Predictors](#3-dynamic-prediction-1-bit-and-2-bit-predictors)
4. [Global History Predictors](#4-global-history-predictors)
5. [The Branch Target Buffer](#5-the-branch-target-buffer)
6. [Modern Tournament and Neural Predictors](#6-modern-tournament-and-neural-predictors)
7. [Spectre: When Prediction Becomes a Vulnerability](#7-spectre-when-prediction-becomes-a-vulnerability)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. Why Branches Are Expensive

Recall from the previous chapter: when a branch instruction enters the pipeline, the processor doesn't know where to fetch the next instruction until the EX stage (2–3 cycles later). On a 5-stage pipeline that's a 2-cycle penalty. On a modern deep-pipeline CPU (Intel Core: ~14 stages, Cortex-A77: ~13 stages), a mispredicted branch flushes 12–19 instructions in-flight — a **15–20 cycle penalty**.

```
Branch frequency in typical code: ~15-20% of instructions
Taken frequency (loops, function calls): ~60% of branches
Misprediction penalty: 15-20 cycles

Without prediction (always stall):
CPI impact = 0.15 × 2 cycles = 0.30 added to CPI

With 70% accurate prediction (30% misprediction rate):
CPI impact = 0.15 × 0.30 × 15 cycles = 0.68 added to CPI  ← worse than stalling!

With 99% accurate prediction (1% misprediction rate):
CPI impact = 0.15 × 0.01 × 15 cycles = 0.02 added to CPI  ← nearly free
```

The lesson: branch prediction only helps if it's very accurate. A bad predictor can perform worse than simply always stalling.

### Quick Check
> 1. If branches appear every 5 instructions and each misprediction costs 15 cycles, what is the CPI penalty at 5% misprediction rate?
> 2. Why does the penalty grow with pipeline depth?
> 3. What does "flushing in-flight instructions" cost in terms of work done?

---

## 2. Static Branch Prediction

**Static** prediction uses no runtime information — the prediction is fixed at compile time or based on simple rules.

### Always Not-Taken
Predict every branch is not taken; continue fetching sequentially. If wrong (branch was taken), flush and fetch from target. Accuracy: ~40% for typical code (since ~60% of branches are taken). Poor but simple.

### Always Taken
Predict every branch is taken; immediately start fetching from the branch target. Requires knowing the target address immediately (it's encoded in the instruction). Accuracy: ~60%.

### Backward Taken / Forward Not-Taken (BTFNT)
A smarter static rule based on branch direction:
- **Backward branches** (target address < current PC) are almost always loops — predict **taken**
- **Forward branches** (target address > current PC) are usually `if` bodies that exit on a condition — predict **not-taken**

Accuracy: ~65–75%. This is the prediction used by early RISC processors and still used as a fallback when no dynamic history exists.

### Compiler Hints
Some ISAs (IA-64/Itanium, some RISC-V extensions) let the compiler embed prediction hints in the branch instruction itself. The hardware uses the hint if no better information is available.

### Quick Check
> 1. Why does "backward taken / forward not-taken" work well for loop code?
> 2. A compiler knows that `if (error_flag)` is almost never true. Which static prediction should it hint?
> 3. What is the best possible accuracy with static prediction? Why can't static prediction do better?

---

## 3. Dynamic Prediction: 1-Bit and 2-Bit Predictors

**Dynamic** prediction learns from the branch's own history at runtime.

### 1-Bit Predictor (Last-Outcome Predictor)

A small table indexed by (part of) the branch's PC address. Each entry stores the last outcome: T (taken) or NT (not taken). Predict the same as last time.

```
Branch at 0x1000:   T T T T T NT T T T T
Prediction:        [T] T T T T T  NT T T T
Correct?:               ✓ ✓ ✓ ✓ ✗  ✗ ✓ ✓ ✓

Accuracy: 7/9 = 78%  (only 2 errors: at the NT and the first T after NT)
```

Works well for loops. Problem: for a loop that runs N times, the predictor mispredicts twice — at the loop exit (predicts T, was NT) and at the first iteration of the next loop invocation (predicts NT, was T). For small loops (e.g., N=2), that's a 50% error rate.

### 2-Bit Saturating Counter (Bimodal Predictor)

Each entry is a 2-bit counter with 4 states:
```
00 = Strongly Not-Taken (SNT)
01 = Weakly Not-Taken (WNT)
10 = Weakly Taken (WT)
11 = Strongly Taken (ST)
```

Predict taken if counter ≥ 10 (i.e., either Weakly or Strongly Taken). The counter increments on taken, decrements on not-taken, clamped at 00 and 11.

```
Transition diagram:
SNT ──T──► WNT ──T──► WT ──T──► ST
SNT ◄──NT── WNT ◄──NT── WT ◄──NT── ST
```

**Key advantage**: A single not-taken outcome doesn't immediately flip the prediction. The predictor must see two consecutive not-taken outcomes to switch from "predict taken" to "predict not-taken". This makes it robust against isolated anomalies and handles loop exit correctly (one NT at loop end doesn't flip the prediction for the next invocation).

A practical bimodal predictor uses 2K–4K entries (2K–4K × 2 bits = 512 bytes–1KB of state). Each branch maps to an entry via its low PC bits.

**Aliasing problem**: Multiple branches map to the same entry and interfere. Solutions: larger tables, better indexing.

### Quick Check
> 1. For a loop that runs exactly 3 times, how many mispredictions does a 1-bit predictor make per loop invocation? A 2-bit predictor?
> 2. Draw the state transitions for a 2-bit counter seeing the pattern: T T T NT T T T NT T T
> 3. Why can't you simply use a huge predictor table with one entry per PC address?

---

## 4. Global History Predictors

The bimodal predictor only uses each branch's own history. But branches are correlated — the outcome of branch B often depends on what branch A did earlier in the same code path.

```
Example:
  if (x > 5) flag = 1;   // branch A
  if (flag == 1) ...;    // branch B
```

Branch B's outcome is perfectly correlated with branch A's. A predictor that tracks global history can exploit this.

### Two-Level Adaptive Predictor (Yeh & Patt, 1991)

**Global History Register (GHR)**: A shift register holding the outcomes of the last N branches (taken=1, not-taken=0). N=12–16 bits is common.

**Pattern History Table (PHT)**: A 2ᴺ-entry table of 2-bit counters. The GHR value indexes into the PHT to select the prediction.

```
GHR: [1 0 1 1 0 1 ...] ← last N branch outcomes
           │
           ▼
PHT: [entry 101101...] → 2-bit counter → prediction
```

After each branch, the GHR shifts left and the actual outcome is shifted in. The corresponding PHT counter is updated.

**Local History Predictor**: Instead of one global history register, keep a separate history register per branch PC. Better for loops with regular patterns.

**Tournament Predictor**: Use both a global predictor and a local predictor, plus a "chooser" table that tracks which predictor has been more accurate for each branch. Use the better one. This is what AMD K8/K10 and DEC Alpha 21264 used — typically 4–8KB of total state.

### Quick Check
> 1. How does the global history register capture correlations between different branches?
> 2. If the GHR is 16 bits, how large is the Pattern History Table?
> 3. What is the advantage of a tournament predictor over either global or local alone?

---

## 5. The Branch Target Buffer

Knowing that a branch is taken isn't enough — you also need to know where it goes. The **Branch Target Buffer (BTB)** is a cache that maps branch PC addresses to their targets.

```
BTB (typically 512–4096 entries):
  PC → Target Address + Branch Type + Prediction State
```

When the fetch unit reads an instruction and the instruction's PC hits in the BTB, the fetch unit knows immediately:
1. This is a branch
2. It is predicted taken (or not-taken)
3. If taken, the target is X

This allows the CPU to redirect the fetch stream before the branch even reaches the decode stage — a huge win for deep pipelines.

**Indirect branches** (where the target address is in a register, like function pointers, switch/jump tables, virtual function calls) are harder because the target changes each time. The **Indirect Branch Predictor** (a separate structure, sometimes called the IBTB) tracks the history of targets for each indirect branch PC.

**Return Address Stack (RAS)**: Function returns (`RET`) are a special case — the return address is pushed when `CALL` executes and popped when `RET` executes. The processor maintains a small hardware stack (8–32 entries) mirroring the call stack. When a CALL is seen, the predicted return address is pushed. When RET is seen, it's popped. Accuracy for returns: ~99%+ (the RAS is almost always correct unless the call stack is deeper than the RAS).

### Quick Check
> 1. What two pieces of information does the BTB provide that the branch predictor alone doesn't?
> 2. Why are indirect branches harder to predict than direct branches?
> 3. Why is the Return Address Stack almost always correct?

---

## 6. Modern Tournament and Neural Predictors

### TAGE Predictor (Tagged Geometric History Length)

TAGE (Seznec, 2006) is the basis for virtually all modern high-performance branch predictors. It uses multiple predictor tables with geometrically increasing history lengths (e.g., 5, 10, 20, 40, 80, 160 bits). Each table is indexed by XOR of PC bits with global history bits of that length. The prediction comes from the table with the longest matching history.

Intel, AMD, ARM, and Apple all use TAGE-based predictors internally. TAGE achieves misprediction rates below 2% on SPEC CPU benchmarks.

### Perceptron/Neural Predictors

A perceptron is a simple neural network — a weighted sum of inputs thresholded to give a binary output. Neural branch predictors use a perceptron to compute the branch prediction from a feature vector of recent branch outcomes.

Samsung Exynos and some AMD Zen 4 cores use perceptron-based components. They are particularly good at capturing long-history correlations that TAGE-style predictors miss.

### State-of-the-Art Numbers

| Predictor | Misprediction Rate (SPEC CPU2017) |
|-----------|----------------------------------|
| Bimodal | ~8-10% |
| Tournament (Alpha 21264) | ~3-5% |
| TAGE | ~1.5-2% |
| Modern production CPU | ~0.5-1% |

The difference between 10% and 1% misprediction rates translates directly to 10-15% overall performance.

### Quick Check
> 1. What makes TAGE better than simple tournament predictors?
> 2. A predictor uses history lengths of 4, 8, 16, 32, 64 bits. How does it decide which history length to use for a given branch?
> 3. Why is 0% misprediction rate impossible in general? (Think about what would be required.)

---

## 7. Spectre: When Prediction Becomes a Vulnerability

Branch prediction has a dark side. In January 2018, researchers disclosed **Spectre** — a vulnerability that weaponizes speculative execution.

Here's the core idea:

```
// Victim code
if (index < array_size) {          // Branch A
    secret = array[index];         // Only executes if branch A taken
    dummy  = another_array[secret * 64];  // Loads data based on secret
}
```

An attacker trains the branch predictor to predict Branch A as "taken" (by running it many times with valid indices). Then the attacker passes an out-of-bounds `index`. The CPU speculatively executes the body (before verifying the bounds check). It reads `array[invalid_index]` — which loads a secret value. It then loads `another_array[secret * 64]` — which brings a cache line into the cache.

The speculative execution is eventually squashed (the bounds check fails). But the cache state is not undone — the cache line loaded based on the secret remains. The attacker can then time accesses to `another_array` to determine which cache line was loaded, and thus recover the secret byte.

This attack worked against every modern processor using speculative execution (Intel, AMD, ARM). Mitigations (software barriers, hardware mitigations in later CPUs) exist but many have performance costs.

Spectre demonstrates that the microarchitecture's invisible optimizations can create security vulnerabilities that are invisible at the ISA level.

### Quick Check
> 1. What is the two-step secret of Spectre: (1) how is the secret read? (2) how is it leaked?
> 2. Why doesn't squashing the speculatively executed instructions fix the problem?
> 3. Spectre showed that the ISA-level security model is not the whole story. What additional layer of security must chip designers consider?

---

## Summary

- Branches cause control hazards because the next instruction isn't known until the branch outcome is determined.
- **Static prediction** (always taken, BTFNT) achieves ~65-75% accuracy with no hardware overhead.
- **1-bit predictor** remembers the last outcome. **2-bit saturating counter** is more stable, especially for loops.
- **Global history predictors** track the sequence of outcomes across multiple branches, capturing inter-branch correlations.
- The **BTB** stores branch targets, enabling the fetch unit to redirect before decode. The **RAS** predicts function return addresses with near-perfect accuracy.
- **TAGE** predictors with multiple geometry history lengths achieve <2% misprediction on typical workloads.
- **Spectre** revealed that speculative execution can leak information through cache side channels — a fundamental tension between performance and security.

---

## Exercises

### Easy
1. A loop runs 100 times. Using a 2-bit predictor starting in state "Weakly Taken," how many mispredictions occur for this loop's back-edge branch?
2. Why does the BTFNT heuristic work well for `for` loops?
3. Draw the 2-bit counter state machine. Starting at "Strongly Not-Taken," trace the states through the pattern: NT T T T NT T T T.

### Medium
4. A branch predictor table has 1024 entries indexed by PC[11:2]. Branch A is at address 0x1000 and branch B is at address 0x5000. Do they alias? What about branch A at 0x1000 and branch C at 0x2000?
5. The Return Address Stack has 8 entries. A program has a recursive function that calls itself 15 times before returning. What happens to the RAS? How many mispredictions on return will occur?
6. Explain in detail why disabling branch prediction would make a program run faster on some branches but slower overall. Under what conditions would disabling prediction speed things up?

### Hard
7. Design a 2-level global history predictor with GHR length = 4. The PHT has 16 entries of 2-bit counters. Trace the predictor through this branch sequence: T NT T T NT T T NT (all from the same branch PC). Show the GHR state and the PHT entry used at each step. What is the accuracy?
8. Spectre requires an attacker to: (1) mistrain the branch predictor, (2) trigger speculative execution past a bounds check, (3) observe a cache timing side channel. Propose a hardware mitigation for each of these three steps and analyze the performance cost of each mitigation.
