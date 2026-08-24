# Chapter 00: What Is AI? A Guide for Complete Beginners

*You do not need to know anything to read this chapter. No maths, no code, no prior experience. Just read.*

---

## Table of Contents

1. Before We Start — A Quick Promise
2. What Is a Program?
3. What Is AI?
4. The Story of AI — 70 Years in 5 Minutes
5. How Do Computers Learn?
6. What Can AI Do Today?
7. What Can AI NOT Do?
8. The Vocabulary You Need
9. How This Course Is Structured
10. Your First Experiment with AI
11. Summary
12. Exercises

---

## 1. Before We Start — A Quick Promise

You may have heard that AI is complicated. That it requires a maths degree, years of study, and expensive computers. That is not true.

Here is the truth: **the ideas behind AI are simple**. The maths that implements those ideas is learnable. And you can build real AI systems — chatbots, image recognizers, your own language model — without spending a single dollar.

This chapter will give you the mental model you need to understand everything that comes after it. If you read this chapter and understand it, the rest of the course will feel like filling in details, not encountering incomprehensible walls.

Let's begin.

---

## 2. What Is a Program?

Before we can understand AI, we need to understand what a normal computer program is.

A normal program is a set of instructions written by a human. The computer follows those instructions exactly, every single time.

Think of a recipe:

```
RECIPE FOR SCRAMBLED EGGS:
1. Crack 2 eggs into a bowl
2. Add a pinch of salt
3. Whisk until the yolk and white mix
4. Heat a pan on medium
5. Pour in the eggs
6. Stir gently until cooked
```

A recipe is a program. If you follow it, you get scrambled eggs. If a fire alarm goes off in step 4 and you skip step 5, you get no eggs. Recipes are rigid. They do not adapt.

Computer programs work the same way. A spell-checker program is a list of instructions: "look up each word in a dictionary, if it's not there, underline it in red." It does this exactly. Every time. With no creativity, no adaptation, no judgement.

This rigidity is useful for most things. But some problems are impossible to solve this way.

**Try to write a recipe for recognizing a cat:**

```
IF the thing has four legs
    AND the thing has pointy ears
    AND the thing has fur
    AND the thing meows
THEN it is a cat
```

Sounds good. But:
- A lion has four legs and pointy ears and fur
- A sphinx cat has no fur
- A cat in a coat is still a cat
- A stuffed animal cat might match all the rules

No matter how many rules you add, you will always find an edge case. There are an infinite number of ways a cat can look. No human can write a recipe for "is this a cat?" that works for all of them.

And yet — you can recognize a cat instantly, even if it is a breed you have never seen, in a photo you have never seen, at an angle you have never seen. You learned what a cat is by seeing thousands of cats.

**Machine learning asks: what if the computer could do the same?**

---

## 3. What Is AI?

AI stands for **Artificial Intelligence**. Intelligence that is artificial — made by us, rather than grown by biology.

Intelligence is a hard word to define. But for our purposes, let's say it means: *the ability to perceive a situation and respond appropriately to achieve a goal*.

Here are three types of things that exhibit intelligence:

| Who | How they got intelligent | Example |
|-----|--------------------------|---------|
| A human | Born with a brain, learned from experience for years | Recognizing your friend's face |
| An animal | Born with a brain, learned from experience | A dog recognizing its owner |
| An AI | Built by engineers, trained on data | A computer recognizing any face |

The key word in the AI row is **trained**. AI systems are not programmed with rules. They are trained with examples.

Here is the fundamental difference:

```
TRADITIONAL PROGRAM:
Engineers write rules.
Computer follows rules.
Computer produces output.

MACHINE LEARNING:
Engineers provide examples (input + correct output).
Computer figures out the rules itself.
Computer applies those rules to new input.
```

The computer does not know what a cat is before training. After seeing one million photos labeled "cat" or "not cat," it figures out the patterns — the shapes of ears, the textures of fur, the proportions of the face — entirely on its own. The engineers never told it any of this. It learned.

This is not magic. It is mathematics. And by the end of this course, you will understand that mathematics completely.

---

## 4. The Story of AI — 70 Years in 5 Minutes

### 1950s: The Dream

In 1950, a British mathematician named Alan Turing asked a question that became the founding question of AI: *Can machines think?*

He proposed the Turing Test: if a human, having a text conversation with something, cannot tell whether it is a human or a machine, the machine is intelligent.

At the time, the question was purely theoretical. Computers filled entire rooms and could barely do arithmetic.

### 1956: The Birth of AI

At a conference at Dartmouth College in 1956, researchers coined the term "Artificial Intelligence" and declared it would be solved within a generation. (It was not.)

The first AI programs were simple: chess-playing programs, theorem provers, programs that could play checkers. They impressed people but could not generalize — they could only do the exact task they were built for.

### 1960s–1970s: The First "AI Winter"

Researchers discovered that the problems they had called "almost solved" were actually very, very hard. Computers were not powerful enough, data was not available, and the algorithms were fundamentally limited.

Funding dried up. Progress stalled. This became known as the "AI winter" — a period when enthusiasm for AI collapsed.

### 1980s: Expert Systems

A new approach emerged: instead of teaching computers to learn, encode the knowledge of human experts as rules. "Expert systems" were built to diagnose diseases, configure computer orders, and analyze chemicals.

For a while, this worked. Companies poured billions into expert systems. But the same problem reappeared: the real world is too complex for any set of rules to capture.

### 1986: Backpropagation — The Algorithm That Changed Everything

In 1986, Geoff Hinton and colleagues published a paper describing **backpropagation** — an efficient algorithm for training neural networks.

Neural networks had been around since the 1950s, but nobody knew how to train deep ones efficiently. Backpropagation solved this. (You will learn exactly how it works in Chapter 13.)

This did not immediately cause a revolution — computers were still too slow, data was still too scarce. But the seed was planted.

### 1990s–2000s: The Machine Learning Era

As computers got faster and the internet created massive datasets, a practical approach emerged: **machine learning**. Instead of programming rules or building neural networks, use statistical algorithms that learn from data.

Support vector machines, decision trees, random forests — these algorithms could beat humans at specific tasks. In 1997, IBM's Deep Blue beat world chess champion Garry Kasparov. Not by thinking like a human, but by evaluating 200 million positions per second.

### 2012: The Deep Learning Revolution

In 2012, a neural network called AlexNet entered the ImageNet competition — a challenge to identify 1,000 types of objects in photos. It won by such a large margin that it shocked the entire field.

AlexNet used three things that finally came together:
1. **A deep neural network** — many layers of artificial neurons
2. **A GPU** — a graphics card that could run the computation fast enough
3. **A massive dataset** — millions of labeled images from the internet

This was the moment deep learning became dominant. Every image recognition system, speech recognition system, and recommendation algorithm you use today traces back to 2012.

### 2017: Transformers — The Architecture That Runs Everything

Google published a paper titled *"Attention is All You Need"* in 2017. It described a new architecture called the **transformer**.

Every modern AI language system — ChatGPT, Claude, Gemini, Llama — uses the transformer architecture. It is the foundation of everything. You will learn how it works in Chapter 19-20, and you will build one from scratch in Chapter 25.

### 2022–present: The AI Explosion

In November 2022, OpenAI released ChatGPT. Within two months, it had 100 million users — the fastest adoption of any technology in history.

What changed? Scale. Language models were trained on almost the entire internet, with hundreds of billions of parameters, for months on clusters of thousands of chips. The result was a system that could have a conversation, write code, explain concepts, and help with almost any task.

Since then, the field has moved at breathtaking speed. New models come out every few months, each one more capable than the last.

---

## 5. How Do Computers Learn?

Here is the simple explanation. You will learn the full details over the next 30 chapters — but you need a mental model now.

### Step 1: Data

Every learning process starts with data. Data means examples. If you want to teach a computer to recognize spam emails, you need thousands of emails labeled "spam" and "not spam." If you want to teach it to predict house prices, you need thousands of houses with their prices.

```
Example data for spam detection:

Email 1: "Congratulations! You have won £1,000,000!" → SPAM
Email 2: "Can we schedule a meeting Thursday?" → NOT SPAM
Email 3: "URGENT: Claim your prize NOW" → SPAM
Email 4: "Your invoice is attached" → NOT SPAM
...10,000 more examples
```

### Step 2: Model

A model is a mathematical function with many adjustable parameters. Think of a very complex equation with millions of knobs that you can turn. At the start, all the knobs are set randomly.

```
Model: a function with millions of knobs

Input: email text
Output: probability it is spam (0 to 1)

Before training:
  "Congratulations! You won!" → outputs 0.47 (wrong — should be near 1.0)
  "Can we meet Thursday?" → outputs 0.51 (wrong — should be near 0.0)
```

### Step 3: Learning

The model looks at each example, makes a guess, compares its guess to the correct answer, and adjusts its knobs to do better next time. This process repeats millions of times.

```
Example: "Congratulations! You won!" → model guesses 0.47
Correct answer: 1.0 (it IS spam)
Error: 0.53 (quite wrong)
Action: adjust knobs so words like "congratulations" and "won" push output higher

After seeing 10,000 examples:

"Congratulations! You won!" → outputs 0.97 ✓
"Can we meet Thursday?" → outputs 0.03 ✓
```

### Step 4: Inference

Once trained, the model can make predictions on new emails it has never seen:

```
New email (never seen during training):
"CLAIM YOUR MILLION DOLLAR PRIZE TODAY"
→ model outputs 0.99 → SPAM ✓
```

That is machine learning. The computer learned the pattern — it figured out that words like "prize", "million", "claim" tend to appear in spam — without anyone telling it this. It discovered the rule itself.

---

## 6. What Can AI Do Today?

Modern AI is astonishingly capable. Here is an honest summary of what it can do well today:

### Language
- Answer questions on almost any topic
- Write essays, code, emails, scripts, poetry
- Summarize long documents
- Translate between languages (often better than human translators for common languages)
- Have conversations indistinguishable from humans in many contexts

### Vision
- Identify objects, faces, animals in photos
- Read text in images (OCR)
- Describe what is happening in a photo
- Generate photorealistic images from text descriptions (DALL-E, Midjourney, Stable Diffusion)

### Code
- Write programs in any programming language
- Debug existing code
- Explain what code does
- Convert code between languages

### Reasoning
- Solve maths problems (though not perfectly)
- Plan multi-step tasks
- Draw logical conclusions from given information

### Science and Research
- Suggest drug candidates (AlphaFold predicted protein structures that took humans decades to find)
- Assist with literature review
- Generate and test hypotheses

---

## 7. What Can AI NOT Do?

It is equally important to know the limits. Modern AI systems:

**Cannot truly understand the world.** They are very sophisticated pattern-matchers. When an LLM appears to "understand" your question, it is doing something much closer to "this sequence of words is statistically likely to be followed by these other words" than "I genuinely comprehend what you are asking."

**Make things up (hallucinate).** Language models are not looking things up in a database. They are generating text that is statistically plausible. Sometimes that text is factually wrong. Always verify important factual claims.

**Cannot learn from your conversation.** By default, when you close a chat session with ChatGPT or Claude, the model forgets everything. It does not update its weights based on your conversation. Each session starts fresh.

**Are not conscious or sentient.** There is no "experience" of being ChatGPT. No pain, no joy, no goals of its own. These are tools.

**Cannot reliably count, do arithmetic, or reason about specific numbers.** Language models are probabilistic text predictors, not calculators. They can write code that does arithmetic, but they should not be trusted to calculate without a tool.

---

## 8. The Vocabulary You Need

These words will appear constantly in this course. Learn them now.

| Word | Plain English meaning |
|------|-----------------------|
| **Algorithm** | A step-by-step procedure for solving a problem |
| **Model** | A trained mathematical function that makes predictions |
| **Training** | The process of adjusting a model on data |
| **Parameter / Weight** | One of the millions of knobs inside a model |
| **Dataset** | A collection of examples used for training or testing |
| **Input / Features** | The information you give to a model |
| **Output / Prediction** | What the model produces |
| **Label / Target** | The correct answer, used during training |
| **Inference** | Using a trained model to make predictions on new data |
| **Loss** | A number measuring how wrong the model's predictions are |
| **Neural Network** | A type of model loosely inspired by the brain |
| **Deep Learning** | Training neural networks with many layers |
| **LLM** | Large Language Model — a model trained to understand and generate text |
| **Token** | A piece of text (roughly a word or part of a word) that LLMs process |
| **Prompt** | The input text you send to an LLM |
| **Context Window** | The maximum amount of text an LLM can "see" at once |
| **Fine-tuning** | Additional training on a specific dataset to specialize a model |
| **Embedding** | A list of numbers that represents the meaning of a word or sentence |
| **Vector** | A list of numbers |
| **Agent** | An AI system that can take actions in the world (search the web, run code, etc.) |

---

## 9. How This Course Is Structured

Here is the journey you are going on:

```
PART 0 — What Even Is AI?  ← YOU ARE HERE
  You understand what AI is. No maths, no code.

PART 1 — Python
  You learn to program. Python is the language AI is built in.

PART 2 — Classical ML
  You build simple models. You understand the learning process.

PART 3 — Deep Learning
  You learn about neural networks — the foundation of modern AI.

PART 4 — Transformers
  You understand the architecture behind GPT and Claude.

PART 5 — Build Your Own Language Model
  You write every line of a tiny GPT yourself.

PART 6 — Modern AI Toolkit
  You use real APIs, build RAG systems, fine-tune models.

PART 7 — Claude Code
  You use AI as a professional coding partner.

PART 8 — AI Agents
  You build agents that can act in the world.

PART 9–10 — Production + Projects
  You build complete systems you could deploy tomorrow.
```

Each part builds on the last. Do not skip ahead. A builder who skips foundations ends up with a building that falls down.

---

## 10. Your First Experiment with AI

Before you install anything, do this: have a real conversation with an AI. Go to claude.ai or chat.openai.com and try these prompts:

**Experiment 1 — Ask it to explain something:**
```
Explain how a rainbow forms, like I'm 10 years old.
```

Notice how it adjusts its language for your request.

**Experiment 2 — Ask it to write code:**
```
Write a Python program that prints "Hello, World!" 10 times.
```

Even if you do not know Python yet, you will be able to check this works by Chapter 01.

**Experiment 3 — Try to make it hallucinate:**
```
What is the capital of the fictional country Zagoria?
```

Notice what happens. Does it make up an answer? Does it admit it does not know?

**Experiment 4 — Give it a real task:**
```
I need to organize my bedroom. Make me a step-by-step plan.
```

Notice how it reasons through a practical problem.

**Experiment 5 — Ask it about itself:**
```
Do you have feelings? Are you conscious? What actually happens when you generate a response?
```

This is the most interesting experiment. The answer it gives — and the nuance in that answer — reveals a lot about what these systems are.

---

## Summary

Here is what you now know:

- A **normal program** follows rules written by humans. It cannot adapt.
- **Machine learning** is a different approach: give a computer examples, and let it figure out the rules.
- **Neural networks** are the most powerful type of ML model — especially for language, images, and sequences.
- **Transformers** are the architecture behind every modern AI language system.
- AI today can do remarkable things (language, vision, code) but also has clear limitations (hallucination, no real understanding, no persistent memory by default).
- The next 52 chapters will teach you everything — from Python basics to building your own language model to deploying production AI agents.

---

## Exercises

**Experiment (no coding required):**

1. Try all five AI experiments from Section 10. Write down what surprised you.

2. Ask an LLM to solve this riddle: "A farmer has 17 sheep. All but 9 die. How many sheep does he have?" Did it get it right? (The answer is 9.) Ask it to explain its reasoning.

3. Ask an LLM: "What is 127 × 83?" Then verify the answer with a calculator. Ask the LLM to solve it step by step and see if it gets each step right.

4. Think of one thing you do that you cannot write rules for — something you just "know" how to do. Could a machine learn to do it from examples? What data would it need?

5. Look up one of these: AlphaGo, AlphaFold, GitHub Copilot, DALL-E. Write two sentences explaining what it does and why it is impressive.

**Discussion questions (think about these as you go through the course):**

- If an AI can pass the Turing Test, does that mean it is intelligent? Or just good at faking intelligence?
- Should there be limits on what AI can do? Who should decide?
- If AI can write code as well as a programmer, what happens to programmers?
