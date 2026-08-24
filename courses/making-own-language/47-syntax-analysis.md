# Chapter 47: Syntax Analysis — Parsing the Token Stream

> "Grammar gives shape to thought. The parser gives shape to programs."

---

## Overview

The lexer hands us a flat list of tokens. The parser's job is to discover the *structure* hidden in that list — to recognize that `let x = 2 + 3` is a variable declaration containing an addition expression, not just five tokens in a row. Structure is what allows the compiler to understand what a program *means*.

This chapter builds the Astra parser using two complementary techniques: **recursive descent** for statements and declarations, and **Pratt parsing** (top-down operator precedence) for expressions. Together these form the most powerful and readable approach used by real-world compilers including Go, Rust, and V8 JavaScript.

---

## What We're Building

The `parser/` package that takes a `[]lexer.Token` and produces an `*ast.Program`. The parser handles syntax errors gracefully, collecting them rather than crashing, and providing useful error messages.

---

## Table of Contents

1. What Is Parsing?
2. Parse Trees vs Abstract Syntax Trees
3. Top-Down Parsing: Recursive Descent
4. Bottom-Up Parsing: LR and Its Variants
5. Operator Precedence Parsing: Pratt Parsing
6. Error Handling and Recovery
7. Astra's Grammar (BNF)
8. Implementing Recursive Descent for Statements
9. Implementing Pratt Parsing for Expressions
10. Astra Build Milestone: Complete Parser
11. Exercises
12. Summary

---

## 1. What Is Parsing?

Parsing is the process of recognizing grammatical structure in a sequence of tokens. The parser answers: does this sequence of tokens form a valid Astra program? If so, what is its structure?

Think of diagramming a sentence in English class:

```
"The quick brown fox jumps over the lazy dog"

           Sentence
          /        \
    Noun Phrase   Verb Phrase
    /   |   \         |    \
  Det  Adj  Noun    Verb   Noun Phrase
  The quick fox    jumps  over the lazy dog
```

A parser does the same thing for programs. Instead of "noun phrase" and "verb phrase," it recognizes "variable declaration," "function call," "if statement," etc.

The grammar of a programming language is defined formally using **Backus-Naur Form (BNF)** or **Extended BNF (EBNF)**. Each rule describes how a construct can be composed of smaller constructs:

```
program     → declaration*
declaration → fnDecl | structDecl | importDecl
fnDecl      → "fn" IDENT "(" params? ")" ("->" type)? block
params      → param ("," param)*
param       → IDENT ":" type
block       → "{" statement* "}"
statement   → varDecl | ifStmt | forStmt | whileStmt | returnStmt | exprStmt
```

---

## 2. Parse Trees vs Abstract Syntax Trees

The **parse tree** (also called a *concrete syntax tree*) faithfully represents every symbol in the grammar — including all the punctuation, grouping constructs, and intermediate rule nodes that only exist to define structure.

The **abstract syntax tree** (AST) strips away syntactic sugar. It keeps only the semantically meaningful nodes.

Example: `2 + 3 * 4`

```
Parse Tree:                    AST:
expression                     BinaryExpr(+)
├─ expression                  ├─ IntLit(2)
│  └─ term                     └─ BinaryExpr(*)
│     └─ factor                   ├─ IntLit(3)
│        └─ NUMBER(2)             └─ IntLit(4)
├─ PLUS(+)
└─ term
   ├─ factor
   │  └─ NUMBER(3)
   ├─ STAR(*)
   └─ factor
      └─ NUMBER(4)
```

The parse tree is cluttered with intermediate nodes (`expression`, `term`, `factor`) that exist only to encode operator precedence in the grammar. The AST encodes precedence *structurally* — `*` is lower in the tree, meaning it binds tighter — without the intermediate nodes.

We build the AST directly during parsing, never materializing the parse tree. This is standard practice.

---

## 3. Top-Down Parsing: Recursive Descent

In **top-down parsing**, we start from the topmost grammar rule (`program`) and try to match it by expanding rules downward until we reach the tokens.

**Recursive descent** implements this with one Go function per grammar rule. Each function:
1. Looks at the current token (and sometimes the next one)
2. Decides which alternative of the rule applies
3. Consumes the expected tokens
4. Recursively calls other rule functions
5. Returns an AST node

```go
// Grammar rule: ifStmt → "if" expression block ("else" block)?
func (p *Parser) parseIfStatement() *ast.IfStatement {
    pos := p.expect(lexer.KW_IF).Pos   // consume "if"
    condition := p.parseExpression()    // recursive call for expression rule
    thenBlock := p.parseBlock()         // recursive call for block rule

    var elseBlock *ast.BlockStatement
    if p.check(lexer.KW_ELSE) {
        p.advance()                     // consume "else"
        elseBlock = p.parseBlock()
    }

    return &ast.IfStatement{
        Cond: condition,
        Then: thenBlock,
        Else: elseBlock,
        Pos:  pos,
    }
}
```

**Key operations in recursive descent:**
- `peek()` — look at the current token without consuming
- `check(type)` — true if current token has given type
- `advance()` — consume and return current token
- `expect(type)` — consume if type matches, else error
- `match(types...)` — if current matches any type, consume and return true

**The look-ahead issue:** Recursive descent works naturally for LL(1) grammars (where 1 token of lookahead is enough to decide which rule applies). Most programming constructs need only 1 token: `fn` means function declaration, `if` means if statement, `while` means while loop. Some need 2 (LL(2)): in Astra, `IDENT COLON` means a struct field, while `IDENT LPAREN` means a function call.

**Left recursion:** Recursive descent cannot directly handle *left-recursive* grammar rules:

```
expression → expression "+" term    ← left recursive! causes infinite loop
expression → term "+" expression    ← right recursive, but wrong associativity
```

Left recursion causes `parseExpression()` to call itself immediately without consuming any token, causing infinite recursion. This is why expressions need a different approach: Pratt parsing.

---

## 4. Bottom-Up Parsing: LR and Its Variants

While we won't implement LR parsing for Astra, understanding it helps you appreciate why Pratt parsing became popular.

**LR parsing** works bottom-up: it reads tokens from left to right (L), producing a rightmost derivation (R). It uses a stack and a *parsing table*.

Two operations:
- **Shift:** push the next token onto the stack
- **Reduce:** pop a sequence matching the right-hand side of a rule, push the rule's left-hand side

```
Input: 2 + 3
Stack: (empty)

Shift 2:    Stack: [2]
Reduce to factor: Stack: [factor(2)]
Reduce to term: Stack: [term(factor(2))]
Shift +:    Stack: [term, +]
Shift 3:    Stack: [term, +, 3]
Reduce to factor: Stack: [term, +, factor(3)]
Reduce to term: Stack: [term, +, term(factor(3))]
Reduce to expr: Stack: [expr(term + term)]
Accept!
```

**Variants:**
- **SLR (Simple LR):** smallest tables, weakest (can't handle some common grammars)
- **LALR(1):** more powerful, tables still compact, used by yacc/bison
- **Canonical LR(1):** most powerful, but table can be huge
- **Parser generators** (yacc, bison, ANTLR, tree-sitter) implement LALR or similar

LR parsers are powerful and handle a wider class of grammars. But they are hard to implement by hand, and error messages are harder to make readable. For Astra, recursive descent + Pratt is the better choice.

---

## 5. Operator Precedence Parsing: Pratt Parsing

**Pratt parsing** (invented by Vaughan Pratt in 1973, rediscovered by many) elegantly handles operator precedence and associativity within a top-down framework.

The core insight: each token type has a **binding power** (BP) — a number representing how tightly it binds to adjacent operands. Higher BP = tighter binding.

```
+   has BP = 6   (lower precedence)
*   has BP = 7   (higher precedence)
```

And binding power comes in two flavors:
- **Left binding power (lbp):** how strongly does this token bind to the expression on its LEFT?
- **Right binding power (rbp):** what minimum BP must the right operand have to be included?

The `parseExpression(minBP)` function:
1. Parse the leading operand (a "prefix" expression: literal, identifier, unary operator, `(expr)`)
2. Loop: peek at the next operator. If its lbp > minBP, consume it and parse the right side recursively with the operator's rbp.

```
Parsing: 2 + 3 * 4

parseExpression(0):
  left = 2
  peek: '+' (lbp=6 > 0), consume
  right = parseExpression(6):
    left = 3
    peek: '*' (lbp=7 > 6), consume
    right = parseExpression(7):
      left = 4
      peek: EOF (lbp=0 ≤ 7), stop
      return 4
    return BinaryExpr(3 * 4)
    peek: EOF (lbp=0 ≤ 6), stop
    return BinaryExpr(3 * 4)
  return BinaryExpr(2 + BinaryExpr(3 * 4))
                    ↑
               Correct! * binds tighter than +
```

**Right associativity:** For `=` (assignment), right-associative means `a = b = c` parses as `a = (b = c)`. This is achieved by calling `parseExpression(lbp - 1)` for the right side instead of `parseExpression(lbp)`.

**Pratt's power:** The same algorithm handles all of:
- Binary infix operators: `+`, `-`, `*`, `/`, `==`, `&&`
- Prefix operators: `-x`, `!x`
- Postfix operators: `x++` (though Astra doesn't have these)
- Function calls: `f(args)` — the `(` token triggers "call" parsing
- Array indexing: `a[i]` — the `[` token triggers "index" parsing
- Member access: `a.b` — the `.` token triggers "field" parsing

This is why Go, Rust, TypeScript, and many other modern compilers use Pratt parsing for expressions.

---

## 6. Error Handling and Recovery

When the parser encounters an unexpected token, it:

1. **Records the error** into the error list with position information
2. **Attempts to synchronize** — advance past tokens until a safe recovery point
3. **Continues parsing** from the recovery point

**Synchronization points** in Astra: the start of a new statement (`let`, `if`, `for`, `return`, `fn`) or a closing brace `}`. These tokens are unlikely to be the *cause* of an error and are a reliable restart point.

```go
func (p *Parser) synchronize() {
    p.advance() // skip the problematic token
    for !p.isAtEnd() {
        // A newline after a statement is a good recovery point
        if p.previous().Type == lexer.NEWLINE { return }
        // These tokens start new statements
        switch p.peek().Type {
        case lexer.KW_FN, lexer.KW_LET, lexer.KW_IF,
             lexer.KW_FOR, lexer.KW_WHILE, lexer.KW_RETURN,
             lexer.KW_STRUCT, lexer.RBRACE:
            return
        }
        p.advance()
    }
}
```

**Parse errors are not fatal:** We collect them and continue, then report all of them at the end before stopping compilation. This way users don't have to fix one error at a time.

---

## 7. Astra's Grammar (BNF)

Here is the complete formal grammar of Astra, which our parser implements:

```
program         → declaration* EOF

declaration     → fnDecl
                | structDecl
                | implDecl
                | importDecl
                | statement

fnDecl          → "fn" IDENT "(" paramList? ")" ("->" type)? block
paramList       → param ("," param)*
param           → IDENT ":" type

structDecl      → "struct" IDENT "{" structField* "}"
structField     → IDENT ":" type ";"

implDecl        → "impl" IDENT "{" fnDecl* "}"

importDecl      → "import" STRING_LIT

block           → "{" statement* "}"

statement       → varDecl
                | assignStmt
                | ifStmt
                | forStmt
                | whileStmt
                | returnStmt
                | breakStmt
                | continueStmt
                | exprStmt

varDecl         → "let" "mut"? IDENT (":" type)? "=" expression newline
assignStmt      → expression assignOp expression newline
assignOp        → "=" | "+=" | "-=" | "*=" | "/="

ifStmt          → "if" expression block ("else" (ifStmt | block))?
forStmt         → "for" IDENT "in" expression ".." expression block
whileStmt       → "while" expression block
returnStmt      → "return" expression? newline
breakStmt       → "break" newline
continueStmt    → "continue" newline
exprStmt        → expression newline

expression      → Pratt parsing (see precedence table below)

primary         → INT_LIT | FLOAT_LIT | STRING_LIT | BOOL_LIT
                | IDENT
                | "(" expression ")"
                | structLiteral
                | listLiteral
                | "-" expression       (unary)
                | "!" expression       (unary)

Infix:
IDENT "(" args ")"   → function call    (BP = 9)
expr "." IDENT       → field access     (BP = 9)
expr "[" expr "]"    → index            (BP = 9)
expr "*" expr        → multiply         (BP = 7)
expr "/" expr        → divide           (BP = 7)
expr "%" expr        → modulo           (BP = 7)
expr "+" expr        → add             (BP = 6)
expr "-" expr        → subtract         (BP = 6)
expr "<" expr        → less than        (BP = 5)
expr ">" expr        → greater than     (BP = 5)
expr "<=" expr       → less or equal    (BP = 5)
expr ">=" expr       → greater or equal (BP = 5)
expr "==" expr       → equal           (BP = 4)
expr "!=" expr       → not equal        (BP = 4)
expr "&&" expr       → logical and      (BP = 3)
expr "||" expr       → logical or       (BP = 2)

type            → "int" | "float" | "bool" | "string"
                | IDENT              (struct type)
                | "[" type "]"       (list type)
                | "fn" "(" types ")" ("->" type)?  (function type)
```

---

## 8. Implementing Recursive Descent for Statements

```go
// parser/parser.go

package parser

import (
    "fmt"
    "github.com/astra-lang/astrac/ast"
    "github.com/astra-lang/astrac/lexer"
)

// ParseError is a syntax error with position.
type ParseError struct {
    Message string
    Pos     lexer.Pos
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("%s: parse error: %s", e.Pos, e.Message)
}

// Parser transforms a token stream into an AST.
type Parser struct {
    tokens   []lexer.Token
    current  int
    errors   []*ParseError
    filename string
}

func New(tokens []lexer.Token, filename string) *Parser {
    return &Parser{tokens: tokens, filename: filename}
}

// Parse parses the entire program.
func (p *Parser) Parse() (*ast.Program, []*ParseError) {
    prog := &ast.Program{}
    p.skipNewlines()

    for !p.isAtEnd() {
        decl := p.parseDeclaration()
        if decl != nil {
            prog.Declarations = append(prog.Declarations, decl)
        }
        p.skipNewlines()
    }
    return prog, p.errors
}

func (p *Parser) parseDeclaration() ast.Declaration {
    defer func() {
        if r := recover(); r != nil {
            if pe, ok := r.(*ParseError); ok {
                p.errors = append(p.errors, pe)
                p.synchronize()
            } else {
                panic(r) // re-panic non-parse errors
            }
        }
    }()

    switch p.peek().Type {
    case lexer.KW_FN:
        return p.parseFnDeclaration()
    case lexer.KW_STRUCT:
        return p.parseStructDecl()
    case lexer.KW_IMPL:
        return p.parseImplDecl()
    case lexer.KW_IMPORT:
        return p.parseImportDecl()
    default:
        return p.parseStatement()
    }
}

// ── Function declarations ────────────────────────────────────────────────────

func (p *Parser) parseFnDeclaration() *ast.FnDeclaration {
    pos := p.expect(lexer.KW_FN).Pos
    name := p.expect(lexer.IDENT)
    p.expect(lexer.LPAREN)

    var params []ast.Param
    if !p.check(lexer.RPAREN) {
        params = p.parseParamList()
    }
    p.expect(lexer.RPAREN)

    var returnType ast.Type
    if p.check(lexer.ARROW) {
        p.advance() // consume "->"
        returnType = p.parseType()
    }

    body := p.parseBlock()
    return &ast.FnDeclaration{
        Name:   name.Lexeme,
        Params: params,
        Return: returnType,
        Body:   body,
        Pos:    pos,
    }
}

func (p *Parser) parseParamList() []ast.Param {
    var params []ast.Param
    params = append(params, p.parseParam())
    for p.check(lexer.COMMA) {
        p.advance() // consume ','
        if p.check(lexer.RPAREN) { break } // trailing comma
        params = append(params, p.parseParam())
    }
    return params
}

func (p *Parser) parseParam() ast.Param {
    name := p.expect(lexer.IDENT)
    p.expect(lexer.COLON)
    typ := p.parseType()
    return ast.Param{Name: name.Lexeme, Type: typ, Pos: name.Pos}
}

// ── Struct and impl declarations ─────────────────────────────────────────────

func (p *Parser) parseStructDecl() *ast.StructDecl {
    pos := p.expect(lexer.KW_STRUCT).Pos
    name := p.expect(lexer.IDENT)
    p.expect(lexer.LBRACE)
    p.skipNewlines()

    var fields []ast.Field
    for !p.check(lexer.RBRACE) && !p.isAtEnd() {
        fieldName := p.expect(lexer.IDENT)
        p.expect(lexer.COLON)
        fieldType := p.parseType()
        fields = append(fields, ast.Field{Name: fieldName.Lexeme, Type: fieldType, Pos: fieldName.Pos})
        p.skipSemicolonOrNewline()
        p.skipNewlines()
    }
    p.expect(lexer.RBRACE)
    return &ast.StructDecl{Name: name.Lexeme, Fields: fields, Pos: pos}
}

func (p *Parser) parseImplDecl() *ast.ImplDecl {
    pos := p.expect(lexer.KW_IMPL).Pos
    name := p.expect(lexer.IDENT)
    p.expect(lexer.LBRACE)
    p.skipNewlines()

    var methods []*ast.FnDeclaration
    for !p.check(lexer.RBRACE) && !p.isAtEnd() {
        methods = append(methods, p.parseFnDeclaration())
        p.skipNewlines()
    }
    p.expect(lexer.RBRACE)
    return &ast.ImplDecl{StructName: name.Lexeme, Methods: methods, Pos: pos}
}

func (p *Parser) parseImportDecl() *ast.ImportDecl {
    pos := p.expect(lexer.KW_IMPORT).Pos
    path := p.expect(lexer.STRING_LIT)
    p.expectNewline()
    return &ast.ImportDecl{Path: path.StringVal, Pos: pos}
}

// ── Statements ────────────────────────────────────────────────────────────────

func (p *Parser) parseStatement() ast.Statement {
    switch p.peek().Type {
    case lexer.KW_LET:
        return p.parseVarDecl()
    case lexer.KW_IF:
        return p.parseIfStatement()
    case lexer.KW_FOR:
        return p.parseForStatement()
    case lexer.KW_WHILE:
        return p.parseWhileStatement()
    case lexer.KW_RETURN:
        return p.parseReturnStatement()
    case lexer.KW_BREAK:
        pos := p.advance().Pos; p.expectNewline()
        return &ast.BreakStatement{Pos: pos}
    case lexer.KW_CONTINUE:
        pos := p.advance().Pos; p.expectNewline()
        return &ast.ContinueStatement{Pos: pos}
    default:
        return p.parseExprOrAssignStatement()
    }
}

func (p *Parser) parseVarDecl() *ast.VarDecl {
    pos := p.expect(lexer.KW_LET).Pos
    isMut := false
    if p.check(lexer.KW_MUT) {
        p.advance()
        isMut = true
    }
    name := p.expect(lexer.IDENT)

    var typ ast.Type
    if p.check(lexer.COLON) {
        p.advance()
        typ = p.parseType()
    }

    p.expect(lexer.ASSIGN)
    value := p.parseExpression()
    p.expectNewline()

    return &ast.VarDecl{
        Name:  name.Lexeme,
        Type:  typ,
        Value: value,
        IsMut: isMut,
        Pos:   pos,
    }
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
    pos := p.expect(lexer.KW_IF).Pos
    condition := p.parseExpression()
    thenBlock := p.parseBlock()

    var elseBlock *ast.BlockStatement
    p.skipNewlines()
    if p.check(lexer.KW_ELSE) {
        p.advance()
        if p.check(lexer.KW_IF) {
            // else if — wrap the nested if in a block
            nestedIf := p.parseIfStatement()
            elseBlock = &ast.BlockStatement{
                Statements: []ast.Statement{nestedIf},
                Pos:        nestedIf.Pos,
            }
        } else {
            elseBlock = p.parseBlock()
        }
    }

    return &ast.IfStatement{Cond: condition, Then: thenBlock, Else: elseBlock, Pos: pos}
}

func (p *Parser) parseForStatement() *ast.ForStatement {
    pos := p.expect(lexer.KW_FOR).Pos
    varName := p.expect(lexer.IDENT)
    p.expect(lexer.KW_IN)
    start := p.parseExpression()
    p.expect(lexer.RANGE)
    end := p.parseExpression()
    body := p.parseBlock()
    return &ast.ForStatement{Var: varName.Lexeme, Start: start, End: end, Body: body, Pos: pos}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
    pos := p.expect(lexer.KW_WHILE).Pos
    cond := p.parseExpression()
    body := p.parseBlock()
    return &ast.WhileStatement{Cond: cond, Body: body, Pos: pos}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
    pos := p.expect(lexer.KW_RETURN).Pos
    var value ast.Expression
    if !p.checkNewline() && !p.isAtEnd() {
        value = p.parseExpression()
    }
    p.expectNewline()
    return &ast.ReturnStatement{Value: value, Pos: pos}
}

func (p *Parser) parseExprOrAssignStatement() ast.Statement {
    expr := p.parseExpression()
    pos := expr.GetPos()

    // Check for assignment operators
    switch p.peek().Type {
    case lexer.ASSIGN, lexer.PLUS_EQ, lexer.MINUS_EQ, lexer.STAR_EQ, lexer.SLASH_EQ:
        op := p.advance().Type
        value := p.parseExpression()
        p.expectNewline()
        return &ast.AssignStmt{Target: expr, Op: string(op), Value: value, Pos: pos}
    }

    p.expectNewline()
    return &ast.ExprStatement{Expr: expr, Pos: pos}
}

func (p *Parser) parseBlock() *ast.BlockStatement {
    pos := p.expect(lexer.LBRACE).Pos
    p.skipNewlines()

    var stmts []ast.Statement
    for !p.check(lexer.RBRACE) && !p.isAtEnd() {
        stmt := p.parseStatement()
        if stmt != nil {
            stmts = append(stmts, stmt)
        }
        p.skipNewlines()
    }
    p.expect(lexer.RBRACE)
    return &ast.BlockStatement{Statements: stmts, Pos: pos}
}

// ── Type parsing ──────────────────────────────────────────────────────────────

func (p *Parser) parseType() ast.Type {
    switch p.peek().Type {
    case lexer.IDENT:
        name := p.advance()
        return &ast.NamedType{Name: name.Lexeme, Pos: name.Pos}
    case lexer.LBRACKET:
        p.advance() // consume '['
        elem := p.parseType()
        p.expect(lexer.RBRACKET)
        return &ast.ListType{Elem: elem}
    default:
        p.panicError(fmt.Sprintf("expected type, got %s", p.peek().Type))
        return nil
    }
}
```

---

## 9. Implementing Pratt Parsing for Expressions

```go
// parser/pratt.go

package parser

import (
    "github.com/astra-lang/astrac/ast"
    "github.com/astra-lang/astrac/lexer"
    "fmt"
)

// Binding power constants (precedence levels).
const (
    BP_NONE       = 0
    BP_ASSIGNMENT = 1  // =  (handled separately)
    BP_OR         = 2  // ||
    BP_AND        = 3  // &&
    BP_EQUALITY   = 4  // == !=
    BP_COMPARISON = 5  // < > <= >=
    BP_TERM       = 6  // + -
    BP_FACTOR     = 7  // * / %
    BP_UNARY      = 8  // ! - (prefix)
    BP_CALL       = 9  // () . []
    BP_PRIMARY    = 10
)

// leftBindingPower returns how tightly a token binds to its left operand.
// This determines operator precedence.
func leftBindingPower(typ lexer.TokenType) int {
    switch typ {
    case lexer.OR_OR:                         return BP_OR
    case lexer.AND_AND:                       return BP_AND
    case lexer.EQ_EQ, lexer.BANG_EQ:         return BP_EQUALITY
    case lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ: return BP_COMPARISON
    case lexer.PLUS, lexer.MINUS:             return BP_TERM
    case lexer.STAR, lexer.SLASH, lexer.PERCENT: return BP_FACTOR
    case lexer.LPAREN, lexer.DOT, lexer.LBRACKET: return BP_CALL
    default:                                  return BP_NONE
    }
}

// parseExpression is the main Pratt parsing loop.
// minBP is the minimum binding power the next operator must have to be consumed.
func (p *Parser) parseExpression() ast.Expression {
    return p.prattParse(BP_NONE)
}

func (p *Parser) prattParse(minBP int) ast.Expression {
    // Step 1: Parse the "null denotation" (prefix position)
    left := p.parsePrefix()

    // Step 2: Loop consuming infix operators with sufficient binding power
    for {
        op := p.peek()
        lbp := leftBindingPower(op.Type)
        if lbp <= minBP { break }

        p.advance() // consume the operator
        left = p.parseInfix(left, op)
    }
    return left
}

// parsePrefix handles expressions that can appear in prefix (left) position:
// literals, identifiers, grouping, unary operators.
func (p *Parser) parsePrefix() ast.Expression {
    tok := p.advance()

    switch tok.Type {
    case lexer.INT_LIT:
        return &ast.IntLiteral{Value: tok.IntVal, Pos: tok.Pos}

    case lexer.FLOAT_LIT:
        return &ast.FloatLiteral{Value: tok.FloatVal, Pos: tok.Pos}

    case lexer.STRING_LIT:
        return &ast.StringLiteral{Value: tok.StringVal, Pos: tok.Pos}

    case lexer.BOOL_LIT:
        return &ast.BoolLiteral{Value: tok.BoolVal, Pos: tok.Pos}

    case lexer.IDENT:
        ident := &ast.Identifier{Name: tok.Lexeme, Pos: tok.Pos}
        // Struct literal: MyStruct { field: value }
        if p.check(lexer.LBRACE) {
            return p.parseStructLiteral(tok)
        }
        return ident

    case lexer.LPAREN:
        // Grouped expression: (expr)
        inner := p.prattParse(BP_NONE)
        p.expect(lexer.RPAREN)
        return inner

    case lexer.LBRACKET:
        // List literal: [elem, elem, ...]
        return p.parseListLiteral(tok.Pos)

    case lexer.MINUS:
        // Unary negation: -expr
        operand := p.prattParse(BP_UNARY)
        return &ast.UnaryExpr{Op: "-", Operand: operand, Pos: tok.Pos}

    case lexer.BANG:
        // Logical not: !expr
        operand := p.prattParse(BP_UNARY)
        return &ast.UnaryExpr{Op: "!", Operand: operand, Pos: tok.Pos}

    default:
        p.panicError(fmt.Sprintf("expected expression, got %s %q", tok.Type, tok.Lexeme))
        return nil
    }
}

// parseInfix handles operators that appear between two expressions.
func (p *Parser) parseInfix(left ast.Expression, op lexer.Token) ast.Expression {
    switch op.Type {
    // Binary arithmetic and logical operators
    case lexer.PLUS, lexer.MINUS, lexer.STAR, lexer.SLASH, lexer.PERCENT,
         lexer.LT, lexer.GT, lexer.LT_EQ, lexer.GT_EQ,
         lexer.EQ_EQ, lexer.BANG_EQ,
         lexer.AND_AND, lexer.OR_OR:
        // Left-associative: use same binding power for right side
        right := p.prattParse(leftBindingPower(op.Type))
        return &ast.BinaryExpr{Left: left, Op: string(op.Type), Right: right, Pos: op.Pos}

    case lexer.LPAREN:
        // Function call: expr(args...)
        return p.parseCallExpr(left, op.Pos)

    case lexer.DOT:
        // Field access: expr.field
        field := p.expect(lexer.IDENT)
        return &ast.FieldAccess{Object: left, Field: field.Lexeme, Pos: op.Pos}

    case lexer.LBRACKET:
        // Index: expr[index]
        index := p.prattParse(BP_NONE)
        p.expect(lexer.RBRACKET)
        return &ast.IndexExpr{Object: left, Index: index, Pos: op.Pos}

    default:
        p.panicError(fmt.Sprintf("unexpected infix operator %s", op.Type))
        return nil
    }
}

// parseCallExpr parses the argument list of a function call.
// The '(' has already been consumed.
func (p *Parser) parseCallExpr(fn ast.Expression, pos lexer.Pos) *ast.CallExpr {
    var args []ast.Expression
    if !p.check(lexer.RPAREN) {
        args = append(args, p.prattParse(BP_NONE))
        for p.check(lexer.COMMA) {
            p.advance()
            if p.check(lexer.RPAREN) { break } // trailing comma
            args = append(args, p.prattParse(BP_NONE))
        }
    }
    p.expect(lexer.RPAREN)
    return &ast.CallExpr{Function: fn, Args: args, Pos: pos}
}

// parseStructLiteral parses: TypeName { field: value, field: value }
func (p *Parser) parseStructLiteral(nameTok lexer.Token) *ast.StructLiteral {
    p.expect(lexer.LBRACE)
    p.skipNewlines()

    var fields []ast.FieldInit
    for !p.check(lexer.RBRACE) && !p.isAtEnd() {
        name := p.expect(lexer.IDENT)
        p.expect(lexer.COLON)
        value := p.prattParse(BP_NONE)
        fields = append(fields, ast.FieldInit{Name: name.Lexeme, Value: value, Pos: name.Pos})
        if p.check(lexer.COMMA) { p.advance() }
        p.skipNewlines()
    }
    p.expect(lexer.RBRACE)
    return &ast.StructLiteral{TypeName: nameTok.Lexeme, Fields: fields, Pos: nameTok.Pos}
}

// parseListLiteral parses: [elem, elem, ...]
func (p *Parser) parseListLiteral(pos lexer.Pos) *ast.ListLiteral {
    var elems []ast.Expression
    if !p.check(lexer.RBRACKET) {
        elems = append(elems, p.prattParse(BP_NONE))
        for p.check(lexer.COMMA) {
            p.advance()
            if p.check(lexer.RBRACKET) { break }
            elems = append(elems, p.prattParse(BP_NONE))
        }
    }
    p.expect(lexer.RBRACKET)
    return &ast.ListLiteral{Elements: elems, Pos: pos}
}
```

### Parser Utility Methods

```go
// parser/utils.go
package parser

import (
    "fmt"
    "github.com/astra-lang/astrac/lexer"
)

func (p *Parser) peek() lexer.Token         { return p.tokens[p.current] }
func (p *Parser) previous() lexer.Token     { return p.tokens[p.current-1] }
func (p *Parser) isAtEnd() bool             { return p.peek().Type == lexer.EOF }
func (p *Parser) check(t lexer.TokenType) bool { return p.peek().Type == t }
func (p *Parser) checkNewline() bool        { return p.check(lexer.NEWLINE) }

func (p *Parser) advance() lexer.Token {
    if !p.isAtEnd() { p.current++ }
    return p.previous()
}

func (p *Parser) expect(t lexer.TokenType) lexer.Token {
    if p.check(t) { return p.advance() }
    p.panicError(fmt.Sprintf("expected %s but got %s %q", t, p.peek().Type, p.peek().Lexeme))
    return lexer.Token{}
}

func (p *Parser) expectNewline() {
    if p.checkNewline() || p.isAtEnd() || p.check(lexer.RBRACE) {
        if p.checkNewline() { p.advance() }
        return
    }
    p.panicError(fmt.Sprintf("expected newline or '}', got %s", p.peek().Type))
}

func (p *Parser) skipNewlines() {
    for p.check(lexer.NEWLINE) { p.advance() }
}

func (p *Parser) skipSemicolonOrNewline() {
    if p.check(lexer.SEMICOLON) || p.check(lexer.NEWLINE) { p.advance() }
}

func (p *Parser) panicError(msg string) {
    pe := &ParseError{Message: msg, Pos: p.peek().Pos}
    p.errors = append(p.errors, pe)
    panic(pe)
}

func (p *Parser) synchronize() {
    for !p.isAtEnd() {
        if p.previous().Type == lexer.NEWLINE { return }
        switch p.peek().Type {
        case lexer.KW_FN, lexer.KW_LET, lexer.KW_IF,
             lexer.KW_FOR, lexer.KW_WHILE, lexer.KW_RETURN,
             lexer.KW_STRUCT, lexer.RBRACE:
            return
        }
        p.advance()
    }
}
```

---

## 10. Exercises

1. **Trace parsing:** Manually trace the Pratt parser on the input `a + b * c - d`. Show the state of `left` and the current operator at each step of the loop. Verify that the resulting AST correctly represents `(a + (b * c)) - d`.

2. **Right associativity:** Astra's `=` is right-associative (`a = b = c` → `a = (b = c)`). How would you implement right associativity in the Pratt parser? What changes to `prattParse()` would be needed?

3. **Grammar ambiguity:** The grammar rule `statement → exprStmt | assignStmt` is ambiguous for the parser because both start with an expression. How does our `parseExprOrAssignStatement()` function resolve this ambiguity?

4. **Error recovery test:** Write an Astra program with a syntax error inside a for loop. Trace through what `synchronize()` would do. What would the parser report as errors?

5. **Ternary operator:** Design the Pratt parsing extension needed to add a ternary operator `cond ? then : else` to Astra. What is its binding power? Is it left- or right-associative?

6. **Parsing types:** The `parseType()` function handles `int`, `float`, `bool`, `string`, and `[T]`. Extend it to also handle function types like `fn(int, bool) -> string`.

7. **Grammar for match:** Design the BNF grammar for a `match` statement in Astra:
   ```astra
   match x {
       0 => print("zero")
       1 => print("one")
       _ => print("other")
   }
   ```
   Then sketch the `parseMatchStatement()` function.

8. **Parse a whole file:** Using the complete lexer from Chapter 46 and the parser skeleton from this chapter, parse the Astra sample program from the introduction. Verify each function declaration, struct, and impl block is parsed correctly.

---

## 11. Summary

| Concept | Key Idea |
|---|---|
| Parsing | Recognizing grammatical structure in a token stream |
| Parse tree | Full representation of every grammar symbol |
| AST | Simplified tree; only semantically meaningful nodes |
| Recursive descent | One function per grammar rule; naturally top-down |
| LL(1) | 1 token lookahead sufficient; handled by recursive descent |
| Pratt parsing | Binding power controls precedence; elegant expression parsing |
| Left binding power | How tight an operator binds to its left argument |
| Prefix | Literal, identifier, unary operator, grouped expression |
| Infix | Binary operator, call, field access, index |
| Error recovery | Panic mode: synchronize at statement boundaries; continue parsing |
| Grammar | Complete BNF rules define what is syntactically valid in Astra |

The parser is the bridge between text and meaning. It hands a structured AST to the semantic analyzer (Chapter 49), which checks that the meaning is coherent.
