# Chapter 41: Context-Free Grammars — Describing Programming Language Syntax

> "A grammar is a theory of a language. Write the grammar, and you have written the specification of every possible program."
> — Every programming language designer

---

## Overview

We have established that regular languages — recognized by DFAs and described by regular expressions — are powerful enough to handle individual tokens but not enough to handle the structure of programs. The fundamental limitation is nesting: a regular language cannot track matching delimiters, recursive structure, or balanced pairs. For this, we need **context-free grammars** (CFGs).

This chapter is the longest and most mathematically rich in the volume. We define CFGs formally, study derivations and parse trees, tackle the thorny problem of ambiguity (and the famous "dangling else" problem), and build up a complete formal grammar for the entire Astra language. We also cover FIRST and FOLLOW sets — the mathematical tools that make LL parsing work — and connect everything back to the recursive descent parser you have already implemented.

---

## What We're Building

By the end of this chapter, you will have the **complete formal grammar of Astra** written in EBNF — a single document that precisely specifies every syntactically valid Astra program. This grammar is not just a theoretical artifact; it is the specification from which the recursive descent parser is derived.

---

## Table of Contents

1. Why Regular Languages Aren't Enough — The Formal Proof
2. CFG Definition: Variables, Terminals, Productions, Start Symbol
3. BNF and EBNF Notation
4. Derivations: Leftmost, Rightmost, and Parse Trees
5. Ambiguous Grammars — The Dangling Else Problem
6. Disambiguating Grammars: Precedence and Associativity
7. Grammar Transformations: Left Recursion Removal
8. Left Factoring
9. LL and LR Grammars
10. FIRST and FOLLOW Sets
11. Context-Free Language Properties
12. Astra Build Milestone: The Complete Astra Grammar
13. Exercises
14. Summary

---

## 1. Why Regular Languages Aren't Enough — The Formal Proof

We claimed in Chapter 38 that balanced parentheses cannot be recognized by a DFA. Now we prove it formally using the Pumping Lemma.

**Theorem**: The language L = {(ⁿ)ⁿ : n ≥ 0} (n open parens followed by n close parens) is NOT regular.

**Proof**:
Assume L is regular. By the Pumping Lemma, there exists a pumping length p ≥ 1 such that any string w ∈ L with |w| ≥ p can be written as w = xyz where:
- |xy| ≤ p
- |y| ≥ 1  
- ∀k ≥ 0: xyᵏz ∈ L

Consider the string w = (ᵖ)ᵖ ∈ L. Clearly |w| = 2p ≥ p.

Write w = xyz with |xy| ≤ p. Since the first p characters of w are all open parens, and |xy| ≤ p, the substring xy consists entirely of '(' characters. Therefore y = (ᵐ for some m ≥ 1.

Now pump k = 2: xy²z = xyyZ = (ᵖ⁺ᵐ)ᵖ.

But (ᵖ⁺ᵐ)ᵖ has p+m open parens and p close parens, and m ≥ 1, so p+m ≠ p. This string is NOT in L.

This contradicts the Pumping Lemma. Therefore L is NOT regular. QED.

**Corollary**: Since Astra programs can contain arbitrarily deeply nested blocks `{ { { ... } } }`, and matching these braces requires recognizing a language equivalent to {(ⁿ)ⁿ}, the Astra parser CANNOT be implemented as a DFA or NFA. It requires a pushdown automaton — which is what a context-free grammar parser is.

### What Goes Wrong with a DFA

Here is the intuition. To match balanced parens, you need to know the current **nesting depth** when you see ')'. With a DFA, you would need one state per depth level. But depth is unbounded (Astra allows arbitrarily nested functions, blocks, etc.). A DFA has a fixed, finite number of states. Contradiction — no DFA can handle unbounded depth.

A stack solves this: push on '(', pop on ')'. If the stack is empty at the end and you got a matching pop for every push, the string is balanced. A stack can grow arbitrarily — it is not limited to a fixed number of configurations.

---

## 2. CFG Definition: Variables, Terminals, Productions, Start Symbol

**Definition**: A **Context-Free Grammar** (CFG) G = (V, Σ, R, S) where:

- **V**: A finite set of **variables** (also called **non-terminals** or **grammar symbols**). These represent syntactic categories: `Expression`, `Statement`, `Block`, etc.
- **Σ**: A finite set of **terminals** (the actual tokens/symbols that appear in strings). This is disjoint from V.
- **R**: A finite set of **production rules** (also called **productions** or **rules**). Each rule has the form A → α where A ∈ V and α ∈ (V ∪ Σ)*.
- **S ∈ V**: The **start symbol**.

The word "context-free" in the name refers to the form of production rules: the left-hand side is always a single variable A, never a string of context around A. This contrasts with context-sensitive grammars where rules look like αAβ → αγβ — A can only be rewritten when surrounded by the context α and β.

### Context-Free vs. Context-Sensitive

```
Context-FREE rule:    A → γ
  (A can be rewritten anywhere, regardless of context)

Context-SENSITIVE rule:  αAβ → αγβ
  (A can only be rewritten when surrounded by α and β)
```

Context-free rules are much easier to parse. The context-free property is what allows top-down recursive descent parsing: when the parser sees a non-terminal, it can choose a rule to expand based only on the non-terminal itself (and a few tokens of lookahead), not on what came before it.

### Simple CFG Example: Balanced Parentheses

```
G = (V, Σ, R, S) where:
  V = {S}
  Σ = {'(', ')'}
  R = {
    S → '(' S ')'    -- rule 1: wrap another balanced string
    S → S S          -- rule 2: two balanced strings concatenated
    S → ε            -- rule 3: empty string
  }
  Start: S
```

```
L(G) = {ε, (), (()), ()(), (())(), ...}  -- all balanced paren strings
```

---

## 3. BNF and EBNF Notation

**BNF (Backus-Naur Form)**: The traditional notation for CFGs, invented by John Backus and Peter Naur for the ALGOL 60 language definition.

```
BNF notation:
<non-terminal> ::= alternative₁ | alternative₂ | ...
```

Terminals are written as-is (often quoted) or in UPPERCASE. Non-terminals are written in angle brackets or lowercase.

```
<expression> ::= <expression> '+' <term>
               | <expression> '-' <term>
               | <term>

<term>       ::= <term> '*' <factor>
               | <term> '/' <factor>
               | <factor>

<factor>     ::= '(' <expression> ')'
               | INTEGER
```

**EBNF (Extended BNF)**: Adds convenient shorthand for optional elements, repetition, and grouping. Defined by ISO/IEC 14977.

```
EBNF additions:
  [r]     -- optional: zero or one occurrence of r  (same as r?)
  {r}     -- repetition: zero or more occurrences of r  (same as r*)
  (r)     -- grouping
  r, s    -- concatenation (comma explicit)
  r | s   -- alternation
  r -s    -- r except s (subtraction)
  'x'     -- terminal symbol (quoted)
```

In modern notation (and throughout this book), we use a hybrid:
- Non-terminals in `lowercase_with_underscores`
- Terminals in `'quoted'` or `UPPERCASE` (for token types)
- `?` for optional, `*` for zero-or-more, `+` for one-or-more (regex style)
- `|` for alternation
- `()` for grouping

```
expression := expression '+' term
            | expression '-' term
            | term

-- Equivalent in EBNF using * :
expression := term (('+' | '-') term)*
```

The EBNF form with `*` is more concise and directly maps to `for` loops in a recursive descent parser.

---

## 4. Derivations: Leftmost, Rightmost, and Parse Trees

A **derivation** is a sequence of steps showing how a grammar generates a string. At each step, one non-terminal is replaced by the right-hand side of one of its production rules.

**Notation**: α ⇒ β (α directly derives β in one step)
**Notation**: α ⇒* β (α derives β in zero or more steps)

### Leftmost Derivation

Always expand the leftmost non-terminal. This corresponds to top-down, left-to-right parsing (LL parsing, recursive descent).

```
Grammar:
  E → E '+' T | T
  T → T '*' F | F
  F → '(' E ')' | INT

Leftmost derivation of "2 + 3 * 4":

E
⇒ E '+' T            [E → E + T, expand leftmost E]
⇒ T '+' T            [E → T]
⇒ F '+' T            [T → F]
⇒ INT '+' T          [F → INT]
⇒ INT '+' T '*' F    [T → T * F]
⇒ INT '+' F '*' F    [T → F]
⇒ INT '+' INT '*' F  [F → INT]
⇒ INT '+' INT '*' INT [F → INT]
```

### Rightmost Derivation

Always expand the rightmost non-terminal. Bottom-up parsers (LR parsers) produce the reverse of a rightmost derivation.

```
Rightmost derivation of "2 + 3 * 4":

E
⇒ E '+' T
⇒ E '+' T '*' F
⇒ E '+' T '*' INT
⇒ E '+' F '*' INT
⇒ E '+' INT '*' INT
⇒ T '+' INT '*' INT
⇒ F '+' INT '*' INT
⇒ INT '+' INT '*' INT
```

### Parse Trees

A **parse tree** (also called a **concrete syntax tree** or **derivation tree**) is a tree that shows the structure of a derivation:
- Internal nodes are non-terminals (V)
- Leaf nodes are terminals (Σ) or ε
- Each internal node and its children correspond to one production rule application
- Reading the leaves left-to-right gives the original string

```
Parse tree for "2 + 3 * 4" using the expression grammar:

              E
           /  |  \
          E  '+'   T
          |      / | \
          T     T '*'  F
          |     |      |
          F     F     INT(4)
          |     |
        INT(2) INT(3)
```

The structure reveals that multiplication binds tighter than addition: the `*` subtree is a child of a `T`, while `+` connects two `E` nodes. This is operator precedence encoded in the grammar structure.

---

## 5. Ambiguous Grammars — The Dangling Else Problem

A grammar G is **ambiguous** if there exists a string w ∈ L(G) with **two distinct parse trees** (equivalently, two distinct leftmost or rightmost derivations).

Ambiguity is a problem: if a string has two parse trees, it has two possible meanings. A compiler cannot be deterministic with an ambiguous grammar.

### Classic Arithmetic Ambiguity

Consider this overly simple expression grammar:

```
E → E '+' E | E '*' E | '(' E ')' | INT
```

The string "1 + 2 * 3" has TWO parse trees:

```
Parse tree 1: addition first, then multiply
       E
     / | \
    E  '*' E
   /|\      \
  E '+' E   INT(3)
  |     |
INT(1) INT(2)
Meaning: (1 + 2) * 3 = 9

Parse tree 2: multiplication first (mathematically correct)
       E
     / | \
    E  '+' E
    |    / | \
  INT(1) E '*' E
         |     |
        INT(2) INT(3)
Meaning: 1 + (2 * 3) = 7
```

This grammar is ambiguous. To fix it, we must restructure the grammar to encode precedence directly.

### The Dangling Else Problem

The most famous grammar ambiguity in programming languages:

```
statement := 'if' expr 'then' statement
           | 'if' expr 'then' statement 'else' statement
           | other
```

With this grammar, `if a then if b then s1 else s2` has two parses:

```
Parse 1: else belongs to inner if
  if a then (if b then s1 else s2)

Parse 2: else belongs to outer if
  if a then (if b then s1) else s2
```

Most languages (C, Java, Astra) use the "dangling else binds to nearest if" rule — Parse 1 is correct. But the grammar as written is ambiguous.

**Astra's solution**: Use mandatory braces for `if` bodies. This eliminates the ambiguity completely:

```ebnf
if_stmt := 'if' expression block ('else' (block | if_stmt))?
block   := '{' statement* '}'
```

Since `block` requires `{` and `}`, there is no ambiguity about where the body ends.

```astra
// Astra: no dangling else possible!
if condition {
    if inner_condition {
        // inner body
    }
} else {
    // outer else — unambiguous
}
```

---

## 6. Disambiguating Grammars: Precedence and Associativity

The standard technique for disambiguating arithmetic expressions is to create one non-terminal per **precedence level**, with higher-precedence operators at lower levels (closer to the leaves).

### Precedence Levels in Astra

```
Lowest precedence (evaluated last):
  Assignment:  =
  Logical OR:  ||
  Logical AND: &&
  Equality:    == !=
  Comparison:  < > <= >=
  Addition:    + -
  Multiplication: * / %
  Unary:       - !
  Call/member: f() x.y x[i]
Highest precedence (evaluated first):
  Primary:     literals, identifiers, grouped expressions
```

Each level becomes a non-terminal:

```ebnf
expression  := assignment
assignment  := IDENT '=' assignment | logical_or

logical_or  := logical_and ('||' logical_and)*
logical_and := equality   ('&&' equality)*
equality    := comparison (('==' | '!=') comparison)*
comparison  := addition   (('<' | '>' | '<=' | '>=') addition)*
addition    := multiply   (('+' | '-') multiply)*
multiply    := unary      (('*' | '/' | '%') unary)*
unary       := ('!' | '-') unary | call
call        := primary ('(' args? ')' | '.' IDENT | '[' expression ']')*
primary     := INT | FLOAT | STRING | 'true' | 'false' | IDENT
             | '(' expression ')'
```

This grammar is **unambiguous** because:
1. Each operator appears at exactly one precedence level
2. Lower-precedence operators bind the tree at higher levels
3. The `*` repetition in each level encodes left-associativity (e.g., `a+b+c` = `(a+b)+c`)

### Encoding Associativity

**Left-associative** operators (most arithmetic operators): use left recursion or `(op term)*` pattern.
```
addition := multiply (('+' | '-') multiply)*
-- "a + b + c" is parsed as:
--   multiply = a, then + multiply = b, then + multiply = c
--   result: left-to-right evaluation (((a)+b)+c)
```

**Right-associative** operators (assignment, exponentiation): use right recursion.
```
assignment := IDENT '=' assignment | logical_or
-- "a = b = c" is parsed as:
--   IDENT = assignment → IDENT = (IDENT = logical_or)
--   result: right-to-left evaluation (a = (b = c))
```

---

## 7. Grammar Transformations: Left Recursion Removal

**Left recursion** occurs when a non-terminal A can derive a string starting with A:
```
A → Aα | β
```

This is problematic for **top-down (LL) parsers** because parsing A requires parsing A first — infinite recursion.

The standard algorithm to remove left recursion:

```
Given: A → Aα | β

Replace with:
  A  → β A'
  A' → α A' | ε
```

Where A' (A-prime) is a new non-terminal for the "tail" part.

### Example: Expression Grammar

Original (left-recursive):
```
E → E '+' T | T
T → T '*' F | F
F → '(' E ')' | INT
```

After removing left recursion:
```
E  → T E'
E' → '+' T E' | ε

T  → F T'
T' → '*' F T' | ε

F  → '(' E ')' | INT
```

Or equivalently in EBNF (which makes the transformation invisible):
```
E → T ('+' T)*
T → F ('*' F)*
F → '(' E ')' | INT
```

The EBNF `*` operator naturally represents the right-recursive tail without actually using recursion — it corresponds to a `for` loop in the recursive descent parser.

### Why Astra Uses EBNF With `*`

The Astra grammar uses EBNF notation for binary operators:
```ebnf
logical_or := logical_and ('||' logical_and)*
```

This is equivalent to the left-recursive:
```
logical_or := logical_or '||' logical_and | logical_and
```

But the EBNF form is immediately parseable without left-recursion removal:

```go
func (p *Parser) parseLogicalOr() ast.Expression {
    left := p.parseLogicalAnd()
    for p.check(TOKEN_OR) {
        op := p.advance()
        right := p.parseLogicalAnd()
        left = &ast.BinaryExpr{Left: left, Op: op, Right: right}
    }
    return left
}
```

The `for` loop IS the EBNF `*`, which IS the left-recursion elimination, all in one.

---

## 8. Left Factoring

**Left factoring** is another grammar transformation, used when two rules for the same non-terminal share a common prefix — causing the parser to not know which rule to apply until it has read past the prefix.

**Problem**:
```
stmt → 'if' expr 'then' stmt 'else' stmt
      | 'if' expr 'then' stmt
```

Both alternatives start with `'if' expr 'then' stmt`. A top-down parser cannot decide which rule to use without reading past the common prefix.

**Left factoring solution**: Factor out the common prefix:
```
stmt    → 'if' expr 'then' stmt else_part
else_part → 'else' stmt | ε
```

Now the parser reads `if expr then stmt`, then looks at the next token: if it is `else`, apply `else_part → 'else' stmt`; otherwise apply `else_part → ε`.

**In Astra's grammar**: The mandatory braces solve left-factoring problems naturally. Since every block starts with `{`, the parser can always determine which alternative to use by looking at the next token.

---

## 9. LL and LR Grammars

### LL Grammars (Top-Down Parsing)

**LL(k)** grammars can be parsed top-down with k tokens of lookahead. The two L's stand for: Left-to-right scan, Leftmost derivation.

Properties:
- LL(1): Only one token of lookahead needed — very efficient, simple to implement
- Every LL(1) grammar is unambiguous
- LL(1) grammars cannot have left recursion (by definition — they must be left-recursion-free)
- Most programming language constructs can be written as LL(1) or LL(k)
- Recursive descent parsers implement LL(k) parsing

The Astra parser is an LL(1) recursive descent parser (with occasional LL(2) for things like `->` vs `-`).

### LR Grammars (Bottom-Up Parsing)

**LR(k)** grammars can be parsed bottom-up with k tokens of lookahead. The L stands for Left-to-right scan, R for Rightmost derivation (in reverse).

Properties:
- LR(0): No lookahead — very weak
- SLR(1): Simple LR with 1 token lookahead
- LALR(1): Look-Ahead LR — used by yacc/bison
- LR(1): Full LR with 1 token lookahead — most powerful
- Every LL(1) grammar is LR(1), but not vice versa
- LR parsers can handle left-recursive grammars directly

The tradeoff: LR parsers are more powerful but harder to write by hand (they require parser tables). Most hand-written parsers are LL/recursive descent. yacc, bison, and many parser generators produce LR parsers.

---

## 10. FIRST and FOLLOW Sets

FIRST and FOLLOW sets are the mathematical tools for building LL(1) parsing tables. They tell the parser: given a non-terminal and the current lookahead token, which production rule should we apply?

### FIRST Sets

**FIRST(α)** is the set of terminals that can begin a string derivable from α.

```
FIRST computation rules:
1. If α = ε,              then FIRST(α) = {ε}
2. If α = a (terminal),   then FIRST(α) = {a}
3. If α = A (non-terminal):
   For each rule A → β₁ | β₂ | ... | βₙ:
     Add FIRST(β₁) \ {ε} to FIRST(A)
     Add FIRST(β₂) \ {ε} to FIRST(A)
     ...
     If ε ∈ FIRST(βᵢ) for all i, add ε to FIRST(A)
4. If α = Xβ:
   Add FIRST(X) \ {ε} to FIRST(α)
   If ε ∈ FIRST(X), also add FIRST(β) to FIRST(α)
```

### Example: FIRST sets for Astra expression grammar

```
Grammar:
  expression := assignment
  assignment := IDENT '=' assignment | logical_or
  logical_or := logical_and ('||' logical_and)*
  ...
  primary    := INT | FLOAT | STRING | IDENT | '(' expression ')'
              | 'true' | 'false'

FIRST(primary) = {INT, FLOAT, STRING, IDENT, '(', TRUE, FALSE}
FIRST(unary)   = {'!', '-'} ∪ FIRST(primary)
               = {!, -, INT, FLOAT, STRING, IDENT, (, true, false}
FIRST(multiply) = FIRST(unary) (since multiply := unary ...)
...and so on up the chain...
FIRST(expression) = FIRST(assignment) = FIRST(logical_or) = FIRST(unary)
                  = {!, -, INT, FLOAT, STRING, IDENT, (, true, false}
```

### FOLLOW Sets

**FOLLOW(A)** is the set of terminals that can appear immediately after A in some sentential form. FOLLOW is used for ε-productions: when A → ε is a rule, the parser applies it if the lookahead is in FOLLOW(A).

```
FOLLOW computation rules:
1. $ (end-of-input) ∈ FOLLOW(S) (for start symbol S)
2. For each rule A → αBβ:
   Add FIRST(β) \ {ε} to FOLLOW(B)
   If ε ∈ FIRST(β), also add FOLLOW(A) to FOLLOW(B)
3. For each rule A → αB:
   Add FOLLOW(A) to FOLLOW(B)
```

### Example: FOLLOW sets for Astra

```
FOLLOW(expression) = {RPAREN, RBRACE, COMMA, RBRACKET, EOF, ')', '}'}
  (expression can be followed by: closing paren, closing brace, comma, etc.)

FOLLOW(statement)  = {'}'} ∪ FOLLOW(block)
  (statement is inside block, followed by more statements or '}')

FOLLOW(params)     = {')'}
  (parameter list is followed by closing paren)
```

### LL(1) Parsing Table

The parsing table T[A, a] tells the parser: "when you want to expand non-terminal A and the lookahead is terminal a, which rule should you apply?"

```
T[A, a] = { A → α } if:
  a ∈ FIRST(α)  (a can start what α generates)
  OR
  ε ∈ FIRST(α) AND a ∈ FOLLOW(A)  (α can be empty and a follows A)

A grammar is LL(1) iff no entry T[A, a] has more than one rule.
```

For Astra's expression grammar, each non-terminal has disjoint FIRST sets for its alternatives, making it LL(1).

---

## 11. Context-Free Language Properties

**Closure properties** of context-free languages:
- Union: If L₁ and L₂ are CFLs, so is L₁ ∪ L₂ (add rule S → S₁ | S₂)
- Concatenation: If L₁ and L₂ are CFLs, so is L₁·L₂ (add rule S → S₁ S₂)
- Kleene star: If L is a CFL, so is L* (add rule S → SS | ε)
- Homomorphism: If L is a CFL, so is h(L)
- Reverse: If L is a CFL, so is L^R

**NOT closed under**:
- Intersection: L₁ ∩ L₂ may not be CFL. Example: {aⁿbⁿcⁿ} = {aⁿbⁿ} ∩ {aⁿcⁿ} — each factor is CFL but intersection is not.
- Complement: L may not be CFL.

**The Pumping Lemma for CFLs** (for proving a language is not context-free):

For any CFL L, there exists p ≥ 1 such that any string w ∈ L with |w| ≥ p can be written as w = uvxyz where:
- |vy| ≥ 1 (v or y is non-empty)
- |vxy| ≤ p
- ∀k ≥ 0: uvᵏxyᵏz ∈ L

We can pump v and y simultaneously. Used to prove {aⁿbⁿcⁿ} is not context-free (need to pump both the a's-b's match AND the b's-c's match simultaneously — but they are in different places in the string).

---

## 12. Astra Build Milestone: The Complete Astra Grammar

Here is the complete, formal context-free grammar for the Astra programming language written in EBNF. This is the authoritative specification of every syntactically valid Astra program.

```ebnf
(* ============================================================ *)
(* ASTRA LANGUAGE FORMAL GRAMMAR                                *)
(* EBNF Notation:                                               *)
(*   A := B C       concatenation                               *)
(*   A := B | C     alternation                                 *)
(*   A := B?        optional (zero or one)                      *)
(*   A := B*        repetition (zero or more)                   *)
(*   A := B+        repetition (one or more)                    *)
(*   'x'            literal token                               *)
(*   UPPER           token type from lexer                      *)
(* ============================================================ *)

(* ---- TOP LEVEL ---- *)

program     := declaration* EOF

declaration := fn_decl
             | struct_decl
             | import_decl
             | var_decl

(* ---- IMPORT ---- *)

import_decl := 'import' STRING

(* ---- FUNCTION DECLARATION ---- *)

fn_decl     := 'fn' IDENT '(' params? ')' ('->' type)? block

params      := param (',' param)*
param       := IDENT ':' type

(* ---- STRUCT DECLARATION ---- *)

struct_decl := 'struct' IDENT '{' field* '}'
field       := IDENT ':' type ';'?

(* ---- BLOCK ---- *)

block       := '{' statement* '}'

(* ---- STATEMENTS ---- *)

statement   := var_decl
             | if_stmt
             | while_stmt
             | for_stmt
             | return_stmt
             | expr_stmt

var_decl    := 'let' IDENT (':' type)? '=' expression

if_stmt     := 'if' expression block ('else' (block | if_stmt))?

while_stmt  := 'while' expression block

for_stmt    := 'for' IDENT 'in' expression '..' expression block

return_stmt := 'return' expression?

expr_stmt   := expression

(* ---- EXPRESSIONS (precedence from lowest to highest) ---- *)

expression  := assignment

assignment  := IDENT '=' assignment
             | logical_or

logical_or  := logical_and ('||' logical_and)*

logical_and := equality ('&&' equality)*

equality    := comparison (('==' | '!=') comparison)*

comparison  := addition (('<' | '>' | '<=' | '>=') addition)*

addition    := multiply (('+' | '-') multiply)*

multiply    := unary (('*' | '/' | '%') unary)*

unary       := ('!' | '-') unary
             | call

call        := primary (call_suffix)*

call_suffix := '(' args? ')'
             | '.' IDENT
             | '[' expression ']'

args        := expression (',' expression)*

primary     := INTEGER
             | FLOAT
             | STRING
             | 'true'
             | 'false'
             | IDENT
             | '(' expression ')'
             | struct_literal
             | list_literal

struct_literal := IDENT '{' (IDENT ':' expression (',' IDENT ':' expression)*)? '}'

list_literal   := '[' (expression (',' expression)*)? ']'

(* ---- TYPES ---- *)

type        := 'int'
             | 'float'
             | 'string'
             | 'bool'
             | IDENT
             | list_type
             | fn_type

list_type   := 'List' '<' type '>'

fn_type     := 'fn' '(' (type (',' type)*)? ')' ('->' type)?
```

### Why This Grammar Is LL(1)

Let us verify the key properties:

**1. No left recursion**: Every non-terminal that could be left-recursive uses the `*` pattern instead:
```ebnf
logical_or := logical_and ('||' logical_and)*
-- Instead of: logical_or := logical_or '||' logical_and | logical_and
```

**2. No ambiguity in expressions**: Each operator appears at exactly one precedence level. The grammar is a cascade from `expression` down through decreasing precedence to `primary`.

**3. Disjoint FIRST sets**: For each non-terminal with alternatives:

```
FIRST(statement):
  var_decl:    {'let'}
  if_stmt:     {'if'}
  while_stmt:  {'while'}
  for_stmt:    {'for'}
  return_stmt: {'return'}
  expr_stmt:   FIRST(expression) = {INT, FLOAT, STRING, IDENT, '(', '-', '!', 'true', 'false', '['}

These sets are all disjoint → LL(1) for statements ✓

FIRST(primary):
  INTEGER: {INTEGER}
  FLOAT:   {FLOAT}
  STRING:  {STRING}
  'true':  {TRUE}
  'false': {FALSE}
  IDENT:   {IDENT}
  '(' expr ')': {'('}
  struct_lit: {IDENT} ← CONFLICT with plain IDENT!
  list_lit:   {'['}
```

Notice that `IDENT` and `struct_literal` both start with `IDENT`. This requires LL(2) lookahead: if the current token is `IDENT` and the NEXT token is `{`, parse a struct literal; otherwise parse a plain identifier. This is a common technique called **lookahead extension** — almost all "LL(1)" grammars require occasional 2-token lookahead in practice.

### The Grammar Implemented as Recursive Descent

Here is how the Astra grammar maps directly to recursive descent parser functions:

```go
// Each non-terminal becomes a function

// program := declaration* EOF
func (p *Parser) parseProgram() *ast.Program {
    program := &ast.Program{}
    for !p.isAtEnd() {
        decl := p.parseDeclaration()
        if decl != nil {
            program.Declarations = append(program.Declarations, decl)
        }
    }
    return program
}

// declaration := fn_decl | struct_decl | import_decl | var_decl
func (p *Parser) parseDeclaration() ast.Declaration {
    switch {
    case p.check(TOKEN_FN):
        return p.parseFnDecl()
    case p.check(TOKEN_STRUCT):
        return p.parseStructDecl()
    case p.check(TOKEN_IMPORT):
        return p.parseImportDecl()
    case p.check(TOKEN_LET):
        return p.parseVarDecl()
    default:
        p.error("expected declaration")
        p.synchronize()
        return nil
    }
}

// fn_decl := 'fn' IDENT '(' params? ')' ('->' type)? block
func (p *Parser) parseFnDecl() *ast.FnDecl {
    p.expect(TOKEN_FN)
    name := p.expect(TOKEN_IDENT)

    p.expect(TOKEN_LPAREN)
    var params []ast.Param
    if !p.check(TOKEN_RPAREN) {
        params = p.parseParams()
    }
    p.expect(TOKEN_RPAREN)

    var returnType ast.Type
    if p.match(TOKEN_ARROW) {
        returnType = p.parseType()
    }

    body := p.parseBlock()

    return &ast.FnDecl{
        Name:       name.Lexeme,
        Params:     params,
        ReturnType: returnType,
        Body:       body,
    }
}

// block := '{' statement* '}'
func (p *Parser) parseBlock() *ast.Block {
    p.expect(TOKEN_LBRACE)
    var stmts []ast.Statement
    for !p.check(TOKEN_RBRACE) && !p.isAtEnd() {
        stmt := p.parseStatement()
        if stmt != nil {
            stmts = append(stmts, stmt)
        }
    }
    p.expect(TOKEN_RBRACE)
    return &ast.Block{Statements: stmts}
}

// statement := var_decl | if_stmt | while_stmt | for_stmt | return_stmt | expr_stmt
func (p *Parser) parseStatement() ast.Statement {
    switch {
    case p.check(TOKEN_LET):
        return p.parseVarDecl()
    case p.check(TOKEN_IF):
        return p.parseIfStmt()
    case p.check(TOKEN_WHILE):
        return p.parseWhileStmt()
    case p.check(TOKEN_FOR):
        return p.parseForStmt()
    case p.check(TOKEN_RETURN):
        return p.parseReturnStmt()
    default:
        return p.parseExprStmt()
    }
}

// expression := assignment (entry point for expression parsing)
func (p *Parser) parseExpression() ast.Expression {
    return p.parseAssignment()
}

// assignment := IDENT '=' assignment | logical_or
func (p *Parser) parseAssignment() ast.Expression {
    // LL(2): check if this is assignment
    if p.check(TOKEN_IDENT) && p.checkNext(TOKEN_EQ) {
        name := p.advance() // consume IDENT
        p.advance()         // consume '='
        value := p.parseAssignment() // right-associative
        return &ast.Assignment{Name: name.Lexeme, Value: value}
    }
    return p.parseLogicalOr()
}

// logical_or := logical_and ('||' logical_and)*
// NOTE: The '*' in EBNF becomes a for loop in Go
func (p *Parser) parseLogicalOr() ast.Expression {
    left := p.parseLogicalAnd() // FIRST(logical_and) — enter sub-level
    for p.check(TOKEN_OR) {    // while '||' is next...
        op := p.advance()       // consume '||'
        right := p.parseLogicalAnd()
        left = &ast.BinaryExpr{Left: left, Op: op.Lexeme, Right: right}
    }
    return left
}

// logical_and := equality ('&&' equality)*
func (p *Parser) parseLogicalAnd() ast.Expression {
    left := p.parseEquality()
    for p.check(TOKEN_AND) {
        op := p.advance()
        right := p.parseEquality()
        left = &ast.BinaryExpr{Left: left, Op: op.Lexeme, Right: right}
    }
    return left
}

// equality := comparison (('==' | '!=') comparison)*
func (p *Parser) parseEquality() ast.Expression {
    left := p.parseComparison()
    for p.checkAny(TOKEN_EQEQ, TOKEN_NEQ) {
        op := p.advance()
        right := p.parseComparison()
        left = &ast.BinaryExpr{Left: left, Op: op.Lexeme, Right: right}
    }
    return left
}

// comparison := addition (('<' | '>' | '<=' | '>=') addition)*
func (p *Parser) parseComparison() ast.Expression {
    left := p.parseAddition()
    for p.checkAny(TOKEN_LT, TOKEN_GT, TOKEN_LTE, TOKEN_GTE) {
        op := p.advance()
        right := p.parseAddition()
        left = &ast.BinaryExpr{Left: left, Op: op.Lexeme, Right: right}
    }
    return left
}

// addition := multiply (('+' | '-') multiply)*
func (p *Parser) parseAddition() ast.Expression {
    left := p.parseMultiply()
    for p.checkAny(TOKEN_PLUS, TOKEN_MINUS) {
        op := p.advance()
        right := p.parseMultiply()
        left = &ast.BinaryExpr{Left: left, Op: op.Lexeme, Right: right}
    }
    return left
}

// multiply := unary (('*' | '/' | '%') unary)*
func (p *Parser) parseMultiply() ast.Expression {
    left := p.parseUnary()
    for p.checkAny(TOKEN_STAR, TOKEN_SLASH, TOKEN_PERCENT) {
        op := p.advance()
        right := p.parseUnary()
        left = &ast.BinaryExpr{Left: left, Op: op.Lexeme, Right: right}
    }
    return left
}

// unary := ('!' | '-') unary | call
func (p *Parser) parseUnary() ast.Expression {
    if p.checkAny(TOKEN_BANG, TOKEN_MINUS) {
        op := p.advance()
        operand := p.parseUnary() // right-recursive: !!x is !(!x)
        return &ast.UnaryExpr{Op: op.Lexeme, Operand: operand}
    }
    return p.parseCall()
}

// call := primary call_suffix*
// call_suffix := '(' args? ')' | '.' IDENT | '[' expression ']'
func (p *Parser) parseCall() ast.Expression {
    expr := p.parsePrimary()
    for {
        switch {
        case p.check(TOKEN_LPAREN):
            p.advance()
            var args []ast.Expression
            if !p.check(TOKEN_RPAREN) {
                args = p.parseArgs()
            }
            p.expect(TOKEN_RPAREN)
            expr = &ast.CallExpr{Callee: expr, Args: args}
        case p.check(TOKEN_DOT):
            p.advance()
            field := p.expect(TOKEN_IDENT)
            expr = &ast.MemberExpr{Object: expr, Field: field.Lexeme}
        case p.check(TOKEN_LBRACKET):
            p.advance()
            index := p.parseExpression()
            p.expect(TOKEN_RBRACKET)
            expr = &ast.IndexExpr{Object: expr, Index: index}
        default:
            return expr
        }
    }
}

// primary := INTEGER | FLOAT | STRING | 'true' | 'false' | IDENT
//          | '(' expression ')' | struct_literal | list_literal
func (p *Parser) parsePrimary() ast.Expression {
    switch {
    case p.check(TOKEN_INTEGER):
        tok := p.advance()
        return &ast.IntLiteral{Value: parseInt(tok.Lexeme)}
    case p.check(TOKEN_FLOAT):
        tok := p.advance()
        return &ast.FloatLiteral{Value: parseFloat(tok.Lexeme)}
    case p.check(TOKEN_STRING):
        tok := p.advance()
        return &ast.StringLiteral{Value: unquote(tok.Lexeme)}
    case p.check(TOKEN_TRUE):
        p.advance()
        return &ast.BoolLiteral{Value: true}
    case p.check(TOKEN_FALSE):
        p.advance()
        return &ast.BoolLiteral{Value: false}
    case p.check(TOKEN_IDENT):
        // LL(2) lookahead: struct literal or plain identifier?
        if p.checkNext(TOKEN_LBRACE) {
            return p.parseStructLiteral()
        }
        tok := p.advance()
        return &ast.IdentExpr{Name: tok.Lexeme}
    case p.check(TOKEN_LPAREN):
        p.advance()
        expr := p.parseExpression()
        p.expect(TOKEN_RPAREN)
        return &ast.GroupExpr{Expr: expr}
    case p.check(TOKEN_LBRACKET):
        return p.parseListLiteral()
    default:
        p.error("expected expression")
        return nil
    }
}
```

Notice the beautiful correspondence:
- Each non-terminal → one function
- `A := B C` → call `parseB()` then call `parseC()`
- `A := B | C` → switch on FIRST(B) vs FIRST(C)
- `A := B*` → `for p.check(...)` loop
- `A := B?` → `if p.check(...) { ... }`

The grammar IS the parser. The formal specification and the implementation are the same thing, written in two different languages.

---

## 13. Exercises

**Exercise 41.1** — Formal CFG
Write a formal CFG G = (V, Σ, R, S) for the language of all non-empty comma-separated lists of identifiers, surrounded by parentheses. Examples: "(x)", "(a, b)", "(x, y, z)".

**Exercise 41.2** — Derivations
Using the Astra grammar from the milestone:
a) Write a leftmost derivation of `let x = 3 + 4 * 2`
b) Draw the parse tree for `3 + 4 * 2`
c) Verify that the parse tree correctly encodes operator precedence

**Exercise 41.3** — Ambiguity Detection
Is the following grammar for Astra if-statements ambiguous?
```
stmt := 'if' expr block else_clause | other
else_clause := 'else' stmt | ε
```
If yes, give a string with two parse trees. If no, prove it is not ambiguous. How does Astra's actual grammar (with mandatory braces) avoid this?

**Exercise 41.4** — Left Recursion Removal
The grammar `E → E '+' E | E '*' E | INT` has two problems: left recursion AND ambiguity. Rewrite it to an unambiguous, left-recursion-free LL(1) grammar that correctly encodes multiplication having higher precedence than addition.

**Exercise 41.5** — FIRST and FOLLOW
Compute FIRST and FOLLOW for all non-terminals in the simplified Astra statement grammar:
```
statement := var_decl | if_stmt | return_stmt | expr_stmt
var_decl  := 'let' IDENT '=' expression
if_stmt   := 'if' expression block
return_stmt := 'return' expression?
expr_stmt := expression
```

**Exercise 41.6** — Grammar Extension
Add a `match` statement to the Astra grammar. The syntax should be:
```astra
match value {
    INT => expression
    STRING => expression
    _ => expression
}
```
Write the EBNF production rules for `match_stmt` and add it to the `statement` rule. Update the FIRST set for `statement`.

**Exercise 41.7** — LL(1) Table Construction
Build the LL(1) parsing table for this minimal Astra grammar fragment:
```
S → 'fn' IDENT '(' L ')' B
L → IDENT | IDENT ',' L | ε
B → '{' '}'
```
Draw the table with rows for each non-terminal and columns for each terminal.

**Exercise 41.8** — Grammar Ambiguity in Practice
Consider this attempt to add operator exponentiation to Astra:
```ebnf
multiply := unary (('*' | '/' | '%' | '**') unary)*
```
Is `**` right-associative or left-associative with this rule? (2**3**2 = ?) Show how to change the grammar to make `**` right-associative.

---

## 14. Summary

| Concept | Definition | Astra Relevance |
|---|---|---|
| CFG | (V, Σ, R, S) — generates context-free languages | Complete Astra grammar specification |
| Variable (non-terminal) | Syntactic category (Expression, Statement) | Each has a parser function |
| Terminal | Actual token from lexer | Appears in `p.expect()` calls |
| Production rule | A → α | Maps to parser function body |
| BNF | `<A> ::= B \| C` | Traditional grammar notation |
| EBNF | Adds `?, *, +` to BNF | Direct source of parser code |
| Derivation | Sequence of rule applications | The parser's execution trace |
| Leftmost derivation | Expand leftmost non-terminal first | What recursive descent does |
| Parse tree | Tree showing derivation structure | The AST (roughly) |
| Ambiguity | Two parse trees for one string | Bug in grammar — must eliminate |
| Dangling else | Classic ambiguity in if-else | Solved by mandatory braces in Astra |
| Operator precedence | Encoded by grammar hierarchy levels | Why expression grammar cascades |
| Left recursion | A → Aα — causes infinite recursion | Must eliminate for LL parsers |
| Left factoring | Factor common prefix of alternatives | Makes FIRST sets disjoint |
| LL(1) | Parseable top-down with 1 lookahead | Astra parser is LL(1) + rare LL(2) |
| FIRST(α) | Terminals that can start α | Used for rule selection |
| FOLLOW(A) | Terminals following A | Used for ε-production selection |

**The key insight of this chapter**: The Astra grammar is not just documentation — it is the program from which the parser is derived. Each production rule becomes a function, each alternative becomes a branch, each `*` becomes a loop. The grammar's formal properties (LL(1), unambiguous, left-recursion-free) directly determine that the parser is correct, deterministic, and runs in O(n) time. Formal grammars are the specification language for parsers.

In Chapter 42, we formalize the machine that the parser implements: the Pushdown Automaton.
