# Chapter 52: Claude Code — AI as Your Coding Partner

*AI is not just something you build. It is also a tool that makes you a better builder. This chapter teaches you to use Claude Code — Anthropic's AI coding tool — the way professional engineers actually use it.*

---

## Table of Contents

1. What Is Claude Code?
2. Installing and Getting Started
3. Your First Claude Code Session
4. How Claude Code Works — Under the Hood
5. The Most Useful Commands
6. Permission Modes — Controlling What Claude Can Do
7. CLAUDE.md — Teaching Claude About Your Project
8. Memory — Helping Claude Remember
9. MCP Servers — Connecting Claude to Your Tools
10. Hooks — Automating Your Workflow
11. Real ML Workflows with Claude Code
12. Claude Code vs Other AI Coding Tools
13. Summary
14. Exercises

---

## 1. What Is Claude Code?

Claude Code is an AI coding assistant that runs directly in your terminal. Unlike browser-based chatbots, Claude Code:

- **Lives in your project** — it can see your files, run commands, and edit code
- **Understands context** — it reads your entire codebase, not just what you paste
- **Takes actions** — it can write files, run tests, install packages, and debug errors
- **Persists memory** — it remembers things you tell it between sessions

Think of it as the difference between asking a consultant questions over email versus having them sit next to you at your desk with their hands on your keyboard.

A typical interaction might look like:

```
You:     "My training loop is running but the loss isn't going down. Can you 
          look at the code and figure out why?"

Claude:  [reads train.py, model.py, and your last training log]
         "I see the problem — you are applying your learning rate scheduler 
          before the optimizer step, not after. On line 47 of train.py, 
          move scheduler.step() to after optimizer.step(). Also, your 
          learning rate of 0.1 is very high for Adam — try 1e-4."

         Shall I make those changes?

You:     "Yes"

Claude:  [edits train.py, shows you the diff, runs the training for one 
          epoch to verify the loss now decreases]
```

---

## 2. Installing and Getting Started

### Prerequisites
- Node.js 18 or higher (download at nodejs.org)
- An Anthropic API key (get one at console.anthropic.com)

### Installation

```bash
# Install Claude Code
npm install -g @anthropic-ai/claude-code

# Verify installation
claude --version
```

### First time setup

```bash
# Navigate to your project folder
cd my-ml-project

# Start Claude Code
claude
```

The first time you run it, Claude Code will ask for your Anthropic API key. After that, it opens an interactive session in your terminal.

---

## 3. Your First Claude Code Session

```bash
# Start in a new ML project folder
mkdir sentiment-classifier
cd sentiment-classifier
claude
```

Once inside, you can just type naturally. Here is an example first session for an ML project:

```
> Create a new sentiment classification project. I want to classify movie 
  reviews as positive or negative using a fine-tuned BERT model. Set up 
  the project structure with a requirements.txt, a data loader, a model 
  file, and a training script.

Claude: I'll create a complete sentiment classification project structure...
[creates files, writes code, explains each choice]

> Run the training script on a small sample to make sure it works without errors.

Claude: [runs the script, sees an import error, fixes it automatically, runs again]
        The training loop ran successfully. Loss after 1 step: 0.693

> The dataset I have is in data/reviews.csv with columns "text" and "label". 
  Update the data loader to read that format.

Claude: [reads the current data loader, updates it to match your format]
```

---

## 4. How Claude Code Works — Under the Hood

When you send a message, Claude Code:

1. **Reads relevant files** — it decides which files are relevant to your question and reads them
2. **Plans the work** — for complex tasks, it thinks through the steps before acting
3. **Takes actions** — it edits files, runs commands, reads outputs
4. **Shows you what it did** — all changes are shown as diffs before they are applied (in careful mode) or applied immediately (in auto mode)
5. **Verifies** — it can run your tests or re-run your code to confirm the change worked

The model behind Claude Code is Claude — the same model you might use at claude.ai. The difference is that Claude Code has **tools**: the ability to actually read and write files, run shell commands, and search the internet.

```
Your message
     ↓
Claude (the model) reads your message + relevant files
     ↓
Claude decides what actions to take (tool calls)
     ↓
Tools execute (bash, file read, file write, web search)
     ↓
Results flow back to Claude
     ↓
Claude continues planning → more tools → response to you
```

---

## 5. The Most Useful Commands

### Starting a session

```bash
claude                    # interactive mode in current directory
claude "fix all the type errors in my code"   # one-shot mode
claude --print            # print response to stdout (good for scripts)
```

### Inside a session

```
/help           — show all available commands
/clear          — clear conversation history (start fresh)
/compact        — summarize conversation to save context space
/cost           — show how many tokens this session used
/model          — switch to a different model
Ctrl+C          — exit
```

### Referencing files

```
> Read src/train.py and explain what each section does

> Compare the performance in results/run1.json and results/run2.json 
  and tell me which hyperparameters made the difference

> Here is the error I am seeing: [paste error]
  Fix it
```

### Running code

```
> Run python train.py and tell me what the final validation accuracy is

> Run the full test suite and show me any failures

> Install all the dependencies from requirements.txt
```

---

## 6. Permission Modes — Controlling What Claude Can Do

Claude Code has different permission levels that control how much it can do without asking you first.

### Default mode
Claude asks before making any changes to files or running commands. Best when you are learning or working on important code.

```
> Update the model architecture to use 4 attention heads instead of 2.

Claude: I would like to modify src/model.py. Here are the changes:
[shows diff]
Approve? (y/n)
```

### Auto-approve mode (--dangerously-skip-permissions)

Claude takes all actions without asking. Best for trusted workflows where you know what you want.

```bash
claude --dangerously-skip-permissions "refactor all the training code to use PyTorch Lightning"
```

**Use this carefully.** It can delete files, install packages, and run any command.

### Restricted mode

You can configure which operations are allowed and which require approval in your project's settings.

```json
// .claude/settings.json
{
  "permissions": {
    "allow": [
      "Bash(python *)",
      "Bash(pip install *)",
      "Read(*)",
      "Write(src/*)"
    ],
    "deny": [
      "Bash(rm -rf *)",
      "Bash(git push *)"
    ]
  }
}
```

---

## 7. CLAUDE.md — Teaching Claude About Your Project

`CLAUDE.md` is a file you create in your project root. Claude Code reads it at the start of every session. It is your way of giving Claude a permanent briefing about your project.

A good `CLAUDE.md` answers:
- What is this project?
- What coding style do we use?
- What commands run tests, lint, build?
- What should Claude always/never do?
- What context does Claude need to be useful here?

### Example CLAUDE.md for an ML project

```markdown
# SentimentAI — Movie Review Classifier

## What This Is
A BERT-based sentiment classifier for movie reviews. Trained on IMDb dataset.
Goal: achieve >93% accuracy on the test set.

## Project Structure
- src/data.py       — dataset loading and preprocessing
- src/model.py      — BERT classifier model
- src/train.py      — training loop with early stopping
- src/evaluate.py   — evaluation metrics
- tests/            — unit tests
- notebooks/        — EDA and experiments

## Key Commands
- Train: python src/train.py --epochs 3 --batch-size 32
- Evaluate: python src/evaluate.py --checkpoint checkpoints/best.pt
- Test: pytest tests/ -v
- Lint: ruff check src/

## Style Guide
- Python 3.11, type hints everywhere
- Docstrings in Google style
- Line length: 100 characters
- Use PyTorch, NOT TensorFlow

## Current State
We are optimizing inference speed. The model accuracy is good (93.1%).
The bottleneck is tokenization — it is slow for batch inference.

## Do NOT
- Modify checkpoints/ (these are trained models)
- Change the vocabulary or tokenizer config
- Use any library that requires GPU (this deploys on CPU)
```

The better your `CLAUDE.md`, the less you need to re-explain context every session.

---

## 8. Memory — Helping Claude Remember

Claude Code has two memory mechanisms, and they work differently:

- **`CLAUDE.md` files** (project-level and user-level) are the authoritative way to persist context — instructions, conventions, project facts that should apply to every session. Edit them directly, or ask Claude to add something and it will write it into the relevant `CLAUDE.md` via the `/memory` command.
- **Auto-memory** is passive, not an explicit command: as Claude works, it may decide something is worth remembering for next time (stored under `~/.claude/projects/<project>/memory/`) — there's no reliable "say the magic word and it's guaranteed to persist" trigger. Don't rely on saying "Remember: X" as if it were a stored command; if something must persist, put it in `CLAUDE.md` or ask Claude to update its memory file explicitly.

```
> This project uses Python 3.11 and PyTorch 2.1 — please add that
  to CLAUDE.md so it's there next session.

Claude: Added to CLAUDE.md.

[Later, in a new session]

> Write a training script for the new image classifier.

Claude: [already knows the Python/PyTorch versions from CLAUDE.md]
```

### Types of things worth persisting to CLAUDE.md

```
- This project uses Python 3.11 and PyTorch 2.1.

- Training data lives at /data/processed/, not /data/raw/ —
  raw data has preprocessing errors.

- I'm new to transformers but comfortable with classical ML —
  explain transformer concepts from first principles.

- Always add type hints to functions. Never use the Any type.

- GPU is unavailable on this machine — use CPU-only training.
```

---

## 9. MCP Servers — Connecting Claude to Your Tools

MCP (Model Context Protocol) is how Claude Code connects to external tools and data sources. You write a small server that exposes tools, and Claude can then call those tools.

For ML, useful MCP servers include:
- **File system access** (built in): reading/writing any file
- **Weights & Biases**: query your experiment logs
- **Jupyter**: run notebook cells
- **Database**: query your feature store

### Installing an MCP server

```bash
# Add an MCP server to Claude Code — the -- separates Claude Code's own
# flags from the command that launches the server itself
claude mcp add my-wandb-server -- npx @wandb/mcp-server

# List installed MCP servers
claude mcp list
```

### Using MCP tools in conversation

```
> Show me the training curves for the last 5 experiments on wandb.

Claude: [calls the wandb MCP tool, fetches experiment data, 
        displays loss and accuracy curves, identifies the best run]
```

(Chapter 42 covers building your own MCP server from scratch.)

---

## 10. Hooks — Automating Your Workflow

Hooks are shell commands that run automatically when certain events happen in your Claude Code session. They let you integrate Claude Code into your development workflow.

```json
// .claude/settings.json
{
  "hooks": {
    "PostToolUse": [
      {
        "matcher": "Write",
        "hooks": [
          {
            "type": "command",
            "command": "jq -r '.tool_input.file_path' | xargs -I{} ruff check {} --fix"
          }
        ]
      }
    ]
  }
}
```

`matcher` is a regex against the **tool name** only (`"Write"`, `"Edit|Write"`, etc.) — it doesn't support path patterns. The hook command receives the tool call as JSON on stdin (fields like `tool_input.file_path`), not a shell variable — the snippet above pipes that JSON through `jq` to pull out the file path before running the linter on it. This hook runs your linter (`ruff`) every time Claude writes a file. Other useful hooks:

```json
{
  "hooks": {
    "PostToolUse": [
      {
        // Run tests after Claude writes or edits any file — filter to
        // src/ paths inside the hook command itself, since matcher can't
        "matcher": "Write|Edit",
        "hooks": [{"type": "command", "command": "pytest tests/ -x -q"}]
      },
      {
        // Format the written file with black
        "matcher": "Write",
        "hooks": [{"type": "command", "command": "jq -r '.tool_input.file_path' | xargs -I{} black {}"}]
      }
    ]
  }
}
```

---

## 11. Real ML Workflows with Claude Code

Here are complete workflows for common ML tasks.

### Workflow 1: Exploratory Data Analysis

```
> I just downloaded a new dataset. It's at data/raw/housing.csv.
  Do a complete EDA: load it, describe its shape, check for missing values,
  plot distributions for all numeric columns, compute correlations, 
  and write a summary of what you found.

Claude: [reads the file, writes an EDA script, runs it, saves charts to 
        eda_plots/, writes a summary in eda_notes.md]
```

### Workflow 2: Debugging a Training Run

```
> My training loss was going down for 10 epochs then suddenly jumped up.
  The logs are in logs/training_20240615.txt. The model is in src/model.py
  and the training script is src/train.py. What went wrong?

Claude: [reads all three files]
        "The spike at epoch 10 coincides with when your learning rate 
         scheduler multiplied the LR by 10. This is likely causing the 
         gradient to explode. I see you are using ReduceLROnPlateau with 
         mode='max' but your loss metric is min — this means the scheduler 
         is increasing LR when loss goes up instead of decreasing it."
        
        Shall I fix the scheduler configuration?
```

### Workflow 3: Adding a New Feature

```
> I want to add attention visualization to the model — 
  show which parts of the input the model attends to.
  Add this to the existing inference pipeline.

Claude: [reads model.py and inference.py]
        "I'll add attention map extraction to the model's forward method
         and create a visualize_attention() function. I'll also add a 
         new CLI flag --visualize to the inference script."
        
        [makes changes, runs a test inference, shows the output]
```

### Workflow 4: Code Review Before Committing

```
> Review the changes I am about to commit. Check for: security issues,
  performance problems, missing tests for new functions, and style violations.

Claude: [runs git diff, reads changed files]
        "Found 3 issues:
         1. src/data.py line 45: the file path is concatenated with + 
            instead of using os.path.join — this will break on Windows.
         2. src/model.py: the new encode() method has no unit test.
         3. src/train.py line 112: you are loading the dataset inside 
            the training loop — this reloads from disk every epoch. 
            Move it before the loop."
```

---

## 12. Claude Code vs Other AI Coding Tools

| Tool | Where it lives | Can act on files | Understands whole codebase | Best for |
|------|----------------|-----------------|---------------------------|----------|
| **Claude Code** | Terminal | Yes | Yes | Complex multi-file tasks, debugging |
| GitHub Copilot | IDE (VS Code, JetBrains) | Autocomplete only | Partially | Line-by-line suggestions while typing |
| Cursor | IDE | Yes | Yes | IDE-integrated coding with AI |
| ChatGPT | Browser | No (copy/paste) | No | Explaining concepts, one-off questions |
| Claude.ai | Browser | No (copy/paste) | No | Conversation, explaining, drafting |

**Use Claude Code when:** you have a complex task spanning multiple files, you need to debug by running code, or you want AI-assisted refactoring of an entire codebase.

**Use GitHub Copilot/Cursor when:** you want AI suggestions while typing, inline completions, and IDE integration.

**Use browser chatbots when:** you have a quick question, want to explain something, or are learning a concept.

---

## Summary

- Claude Code runs in your terminal and can read files, write code, and run commands.
- `CLAUDE.md` teaches Claude about your project — write a good one and every session is faster.
- The memory system lets Claude remember preferences across sessions.
- MCP servers extend Claude's capabilities to external tools (databases, APIs, services).
- Permission modes let you control how much Claude can do without asking.
- Hooks automate quality checks (linting, testing) every time Claude writes code.

---

## Exercises

**Easy:**

1. Install Claude Code. Navigate to any project you have (or create an empty folder) and start a session. Ask Claude to explain what Claude Code is. Observe how it responds.

2. Create a simple Python file with a bug (e.g., a division by zero error). Open Claude Code and ask it to find and fix the bug. Observe how it reads the file and makes the change.

3. Create a `CLAUDE.md` for a real or imaginary ML project. Include: project description, coding style preferences, key commands, and one thing Claude should never do.

**Medium:**

4. Start a Claude Code session in a new folder. Ask it to create a complete sentiment analysis project from scratch: data loading, model, training script, and a README. Evaluate the quality of the code it writes.

5. Take any of the ML programs you have written in previous chapters. Ask Claude Code to: add type hints, add docstrings, add unit tests, and reduce any repeated code. How long does this take vs doing it yourself?

6. Configure hooks in `.claude/settings.json` to automatically run `black` (formatter) and `pytest` (tests) every time Claude writes a Python file. Test that the hooks work.

**Hard:**

7. Build a complete ML workflow with Claude Code: start with a raw CSV dataset, ask Claude to do EDA, build a model, evaluate it, and document everything in a README. Your only inputs should be describing what you want — let Claude write all the code.

8. Look up the Claude Code documentation and find a feature not covered in this chapter. Implement it in your workflow and document what you learned.
