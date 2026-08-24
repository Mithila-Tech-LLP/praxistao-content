# Chapter 53: Getting the Most from Claude Code

*Knowing how to install Claude Code is the easy part. This chapter teaches you the mental models and habits that separate engineers who get 10× productivity gains from those who get 10% gains.*

---

## Table of Contents

1. The Right Mental Model — What Claude Is and Is Not
2. Writing Prompts That Actually Work
3. The CLAUDE.md Deep Dive — Real Examples
4. Multi-Session Projects — Staying in Sync
5. IDE Integration — VS Code and JetBrains
6. Using Claude Code for ML-Specific Tasks
7. When NOT to Use Claude Code
8. Building Your AI-Assisted Development Workflow
9. Exercises

---

## 1. The Right Mental Model

The biggest mistake beginners make is treating Claude Code like a search engine — asking short questions and expecting answers.

The right mental model is: **Claude Code is a very capable but new team member who just joined your project.**

A new team member:
- Needs context about the codebase and goals
- Does their best work when you explain the "why," not just the "what"
- Should ask clarifying questions when the task is ambiguous
- Gets dramatically better results when you show examples of what you want
- Should not be trusted to deploy to production without review

### The context hierarchy

Claude Code has three levels of context it draws from:

```
1. AUTO-MEMORY (persists across all sessions)
   ~/.claude/projects/[project]/memory/
   Passive — Claude decides on its own what's worth saving as it works.
   Not a reliable place to put something you need remembered; for that,
   write it into CLAUDE.md (below) instead.

2. CLAUDE.md (project-level, persists across sessions)
   ./CLAUDE.md
   Your project briefing. Architecture, commands, style guide.

3. CONVERSATION (current session only)
   What you have said this session.
   Cleared when you /clear or start a new session.
```

Build from level 1 outward. Everything you tell Claude in conversation that you want to remember later, save to memory or CLAUDE.md.

---

## 2. Writing Prompts That Actually Work

### The bad prompt pattern

```
> Fix my model
```

This gives Claude almost no information. Fix what about the model? What is broken? What is the goal?

### The good prompt pattern

```
> My LSTM text classifier (src/model.py) is overfitting — training 
  accuracy is 97% but validation accuracy is only 71% after 20 epochs.
  The dataset has 5,000 examples per class across 3 classes.
  I have already tried reducing the learning rate. 
  What should I try next, and can you implement the most promising fix?
```

This gives:
- **What** is broken (overfitting)
- **Specific evidence** (train 97%, val 71%)
- **Context** (dataset size, architecture)
- **What you have already tried** (so Claude does not suggest it again)
- **What you want** (recommendation + implementation)

### The PREP framework

| Component | What it means | Example |
|-----------|---------------|---------|
| **Problem** | What is wrong or what do you want? | "Training loss is not decreasing" |
| **Relevant context** | Files, data, configuration | "Train.py uses AdamW, lr=0.001" |
| **Evidence** | Logs, error messages, numbers | "Loss stuck at 0.693 after 100 steps" |
| **Prior attempts** | What you have already tried | "I tried lowering the LR" |

### Showing examples

For style preferences, examples are far more effective than descriptions:

```
> Refactor data.py to use a cleaner coding style. 
  Here is an example of the style I want:
  
  [paste an example function written the way you like]
  
  Apply this style to all functions in data.py.
```

---

## 3. The CLAUDE.md Deep Dive — Real World Examples

A good CLAUDE.md is one of the highest-leverage things you can write. Here are complete examples.

### Example: ML Research Project

```markdown
# ImageNet Transfer Learning Research

## Goal
We are researching how few-shot fine-tuning performance varies with 
the depth of the layer we freeze in ResNet-50.

## Experimental Setup
- Base model: ResNet-50 pretrained on ImageNet
- Target dataset: CUB-200 (bird species classification, 200 classes)
- Training budget: max 30 epochs, early stop at patience=5
- Hardware: 1× A100 (available via RunPod, not local)

## Code Structure
- experiments/    — one folder per experiment, named YYYYMMDD_description
- src/model.py    — ResNet wrapper with configurable freeze depth
- src/train.py    — training script, reads config from experiments/*/config.yaml
- src/eval.py     — evaluation, saves results to experiments/*/results.json
- analysis/       — Jupyter notebooks for plotting and analysis

## Key Commands
- New experiment: python src/train.py --config experiments/my_config.yaml
- Full eval: python src/eval.py --experiment experiments/YYYYMMDD_name/
- Plot results: jupyter notebook analysis/compare_experiments.ipynb

## Important Invariants
- NEVER change a config.yaml after training starts (breaks reproducibility)
- All experiments MUST log to wandb project "imagenet-finetune-depth"
- results.json format must match the schema in analysis/schema.json

## Current Status (as of 2025-06-14)
- Freeze layers 0-2: done, accuracy 78.3%
- Freeze layers 0-4: done, accuracy 81.1%
- Freeze layers 0-6: IN PROGRESS
- Hypothesis: diminishing returns above layer 6

## Do NOT
- Modify trained checkpoints
- Run training locally (it will take 40 hours on CPU)
- Change the evaluation script without versioning the change
```

### Example: Production ML System

```markdown
# RecommendationEngine v2

## What This Is
Collaborative filtering recommendation system for an e-commerce platform.
Serves 50M users. Latency budget: < 20ms p99.

## Architecture
[diagram or description]
User events → Kafka → Feature pipeline → Redis feature store → 
Serving API → Response

## Key Technologies
- Python 3.11, FastAPI, PyTorch 2.1
- Redis for feature store (local: localhost:6379)
- PostgreSQL for user data (local: postgresql://localhost/recengine_dev)
- Docker for containerization

## Development Commands
- Start local services: docker compose up -d
- Run server: uvicorn app.main:app --reload
- Run tests: pytest tests/ -v
- Build: docker build -t recengine .

## Code Standards
- All endpoints must have integration tests
- Type hints required (enforced by mypy)
- No N+1 queries (use eager loading)
- All database queries logged with EXPLAIN ANALYZE in dev

## CRITICAL
- Never hardcode API keys — use environment variables from .env
- Latency regression: if any endpoint exceeds 15ms p50 in tests, 
  DO NOT merge
- The feature schema in schemas/features.py is shared with the 
  feature pipeline — changes require coordination with the data team
```

---

## 4. Multi-Session Projects — Staying in Sync

For a project you work on across many sessions, use this habit:

**Start of session:**
```
> Review the CLAUDE.md and the last 5 git commits. 
  Summarize where the project is and what we should focus on today.
```

**End of session:**
```
> Update the CLAUDE.md "Current Status" section to reflect 
  what we accomplished today. Also, is there anything from 
  this session I should remember permanently?
```

**When something important changes:**
```
> I switched the optimizer from SGD to AdamW today because AdamW
  converged 3x faster. Please add that to CLAUDE.md so any future
  training scripts default to AdamW.
```

---

## 5. IDE Integration

### VS Code

Install the "Claude Code" extension from the VS Code marketplace. This adds:

- **Command palette**: `Ctrl+Shift+P` → "Claude: Ask about selected code"
- **Inline suggestions**: select code, right-click → "Explain with Claude"
- **Terminal integration**: Claude Code runs in VS Code's integrated terminal

### JetBrains (PyCharm, IntelliJ)

Install the Claude Code plugin from the JetBrains marketplace.

### The best of both worlds

Run Claude Code in the terminal for complex tasks (refactoring, debugging, architecture). Use IDE integration for quick questions and inline explanations.

```
IDE strengths:          Claude Code terminal strengths:
- Autocomplete          - Multi-file refactoring
- Quick explanations    - Running and debugging
- Inline errors         - Complex workflows
- Git integration       - Autonomous tasks
```

---

## 6. Claude Code for ML-Specific Tasks

### Task: Hyperparameter search assistance

```
> My model is converging too slowly. Current config:
  - learning_rate: 1e-3
  - batch_size: 32
  - optimizer: SGD
  - epochs: 100
  - val accuracy after 20 epochs: 62%
  
  Based on the architecture (read src/model.py), suggest 3 configurations
  to try next. For each, explain why you expect it to help.
  Then run each one for 5 epochs and compare the validation loss.
```

### Task: Debugging NaN loss

```
> My training is producing NaN loss after about 200 steps. 
  Read train.py, model.py, and the last training log (logs/last_run.txt).
  Find out where the NaN originates and fix it.
```

### Task: Profiling and optimization

```
> My training loop is slow — 2 hours per epoch on a dataset that 
  should take 20 minutes. Profile the training loop and identify 
  the top 3 bottlenecks. Implement the biggest win.
```

### Task: Writing experiment reports

```
> I ran 5 experiments this week. The results are in experiments/.
  Write a summary report (markdown) that:
  - Lists each experiment and its key results
  - Identifies which change had the most impact
  - Recommends the 3 most promising directions for next week
```

---

## 7. When NOT to Use Claude Code

Claude Code is a powerful tool, but it is not the right tool for everything.

**Do not use Claude Code for:**

| Situation | Why | Alternative |
|-----------|-----|-------------|
| Quick syntax questions | Slower than a search | Google / Stack Overflow |
| Learning a new concept | Reading AI output ≠ understanding | Read a book or this course |
| Generating creative text | It can, but better tools exist | Claude.ai with more context |
| Safety-critical code without review | Claude makes mistakes | Always review and test |
| Replacing your own understanding | If you can't read the code Claude writes, you can't maintain it | Understand before using |

**The dependency trap:** It is easy to start relying on Claude for everything and stop building your own understanding. This is dangerous in ML — you will not understand why your model is failing if you never developed the intuition.

**Use Claude Code to accelerate your work, not to replace your thinking.**

---

## 8. Building Your AI-Assisted Development Workflow

Here is a workflow that professional ML engineers actually use:

```
MORNING STARTUP
1. cd to project
2. claude
3. "Review CLAUDE.md and our last commit. What should we work on today?"
4. Set the day's goal

FEATURE DEVELOPMENT
1. Explain the feature to Claude
2. Claude writes initial implementation
3. YOU read every line — understand it before moving on
4. Run tests
5. Review and refine together

DEBUGGING
1. Describe the symptom precisely
2. Paste the relevant error and context
3. Claude diagnoses and proposes fix
4. YOU understand the fix before applying it
5. Verify the fix works

END OF DAY
1. "Summarize what we built today and update CLAUDE.md"
2. git commit the CLAUDE.md update
3. Note any open questions for tomorrow
```

---

## Summary

- **Context is everything**: CLAUDE.md and memory are how you make Claude persistently useful.
- **Prompt quality scales with output quality**: be specific, give evidence, show examples.
- **IDE integration** handles quick tasks; terminal handles complex workflows.
- **Stay in the loop**: always read and understand code Claude writes — never ship what you cannot explain.
- **Update your CLAUDE.md** at the end of every session to keep context fresh.

---

## Exercises

**Easy:**

1. Take a project you are working on (or create a new ML project). Write a `CLAUDE.md` file. Include: what the project is, the key commands, the coding style, and the current status.

2. Try the PREP framework on a real problem you have. Write the prompt using all four components. Notice how much more useful the response is compared to a short question.

**Medium:**

3. At the start of a session, ask Claude to review your CLAUDE.md and the last 5 git commits, then tell you what the project status is. How accurate is Claude's summary?

4. Find a piece of code you wrote that could be cleaner. Show Claude an example of your preferred style and ask it to refactor the code to match. Review every change.

5. Ask Claude to write a failing test for a feature you have not built yet. Then build the feature. This is called test-driven development.

**Hard:**

6. Build a complete workflow where Claude Code helps you: do EDA on a new dataset, build a baseline model, try 3 improvements, write a report comparing all 4 versions, and commit everything with good commit messages. Document each step of your workflow.
