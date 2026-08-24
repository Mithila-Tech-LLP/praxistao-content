---
title: Build an Interpreter
subtitle: Write a complete interpreter from scratch — lexer, parser, and evaluator
category: Systems Programming
difficulty: advanced
duration: 8-12 hours
accent: "#a78bfa"
technologies: [Go]
skills: [Lexing, Parsing, AST, Evaluation, Recursion, Closures]
prerequisites: [basic-programming, data-structures]
repo: build-interpreter
outcomes:
  - Tokenize source code into a stream of typed tokens
  - Build an Abstract Syntax Tree (AST) by recursive descent parsing
  - Evaluate arithmetic, comparison, and boolean expressions
  - Implement variables with lexical scoping
  - Support if/else conditionals and while loops
  - Define and call first-class functions with closures
---

## Overview

Every programming language you use — Go, Python, JavaScript — started as someone writing a lexer, a parser, and an evaluator. In this project you will do the same.

You will build an interpreter for a small but complete language that supports numbers, strings, booleans, variables, functions, closures, and control flow. Each task builds on the previous one. By the end you will have a working language you can extend however you like.

## The Language

```
// arithmetic
1 + 2 * 3       // 7
(1 + 2) * 3     // 9

// variables
let x = 10
let y = x + 5   // 15

// functions
let add = fn(a, b) { a + b }
add(3, 4)       // 7

// closures
let makeAdder = fn(n) { fn(x) { x + n } }
let add5 = makeAdder(5)
add5(3)         // 8
```
