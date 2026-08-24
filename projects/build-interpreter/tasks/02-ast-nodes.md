# Task 02 — AST Nodes

## Concept

After tokenizing, the next step is to decide **what the tokens mean together**. A stream of tokens like `NUMBER(1) PLUS NUMBER(2) STAR NUMBER(3)` could mean either `(1+2)*3` or `1+(2*3)`. The **Abstract Syntax Tree (AST)** encodes the intended structure as a tree, where each node represents a syntactic construct and its children represent sub-expressions.

This task defines all the node types your language needs. You will not parse anything yet — you will only define the data structures and their `String()` methods. The parser (Tasks 03–04) will fill them in.

## The Node Interface

All AST nodes implement:

```go
type Node interface {
    TokenLiteral() string  // the literal of the first token of this node
    String() string        // a human-readable representation (for debugging)
}

type Expression interface {
    Node
    expressionNode()  // marker method to distinguish from Statement
}

type Statement interface {
    Node
    statementNode()   // marker method
}
```

## Expression Nodes

| Node | Fields | String() example |
|------|--------|-----------------|
| `NumberLiteral` | `Value float64` | `"42"` |
| `StringLiteral` | `Value string` | `"hello"` |
| `BoolLiteral` | `Value bool` | `"true"` |
| `Identifier` | `Name string` | `"x"` |
| `PrefixExpr` | `Op string`, `Right Expression` | `"(-5)"` |
| `InfixExpr` | `Left Expression`, `Op string`, `Right Expression` | `"(1 + 2)"` |
| `IfExpr` | `Cond Expression`, `Then *BlockStatement`, `Else *BlockStatement` | `"if (x) { ... } else { ... }"` |
| `FnLiteral` | `Params []string`, `Body *BlockStatement` | `"fn(a, b) { ... }"` |
| `CallExpr` | `Fn Expression`, `Args []Expression` | `"add(1, 2)"` |

## Statement Nodes

| Node | Fields | String() example |
|------|--------|-----------------|
| `LetStatement` | `Name string`, `Value Expression` | `"let x = 42;"` |
| `ReturnStatement` | `Value Expression` | `"return 7;"` |
| `ExpressionStatement` | `Expr Expression` | same as inner expr |
| `BlockStatement` | `Stmts []Statement` | `"{ let x = 1; x }"` |

## Program

```go
type Program struct {
    Statements []Statement
}

func (p *Program) TokenLiteral() string { ... }
func (p *Program) String() string       { ... }
```

`Program.String()` should concatenate the `String()` of all its statements.

## Your Task

Define each node struct and implement both `TokenLiteral()` and `String()` for each.

Guidelines for `String()`:
- `PrefixExpr`: wrap in parens — `"(-x)"`
- `InfixExpr`: wrap in parens — `"(1 + 2)"`
- `IfExpr`: `"if (<cond>) <then>"` and optionally `" else <else>"`
- `FnLiteral`: `"fn(a, b) <body>"`
- `CallExpr`: `"<fn>(arg1, arg2)"`
- `BlockStatement`: `"{ stmt1; stmt2 }"`
- `LetStatement`: `"let <name> = <value>;"`
- `ReturnStatement`: `"return <value>;"`

You do not need to store a `Token` field in each node for this project — `TokenLiteral()` can return a hardcoded string like `"let"` or the operator string.

## Example — Manual Construction

```go
prog := &Program{
    Statements: []Statement{
        &LetStatement{
            Name: "x",
            Value: &InfixExpr{
                Left:  &NumberLiteral{Value: 1},
                Op:    "+",
                Right: &NumberLiteral{Value: 2},
            },
        },
    },
}
fmt.Println(prog.String())
// Output: let x = (1 + 2);
```

## Tests to Pass

1. `NumberLiteral{Value: 3.14}.String()` → `"3.14"`
2. `BoolLiteral{Value: true}.String()` → `"true"`
3. `Identifier{Name: "foo"}.String()` → `"foo"`
4. `PrefixExpr{Op: "-", Right: &NumberLiteral{Value: 5}}.String()` → `"(-5)"`
5. `InfixExpr{Left: &Identifier{Name:"x"}, Op:"*", Right: &NumberLiteral{Value:2}}.String()` → `"(x * 2)"`
6. `LetStatement` with an `InfixExpr` value matches the example above.
7. `CallExpr{Fn: &Identifier{Name:"add"}, Args: [...]}.String()` → `"add(1, 2)"`
