# Task 04 — Parser: Expressions

## Concept

Parsing `1 + 2 * 3` correctly as `1 + (2 * 3)` requires understanding **operator precedence**. The classic approach is a hand-written grammar with one function per precedence level. The **Pratt parser** (top-down operator precedence) does the same thing more elegantly: each token type carries a binding power, and a single loop drives the whole thing.

## Pratt Parser — How It Works

Define precedence levels as integers:

```go
const (
    PREC_LOWEST      = 1
    PREC_EQUALS      = 2  // ==  !=
    PREC_LESSGREATER = 3  // <  >  <=  >=
    PREC_SUM         = 4  // +  -
    PREC_PRODUCT     = 5  // *  /
    PREC_PREFIX      = 6  // -x  !x
    PREC_CALL        = 7  // f(...)
)
```

Map each infix token to its precedence:

```go
var precedences = map[TokenType]int{
    EQEQ:  PREC_EQUALS,
    BANGEQ: PREC_EQUALS,
    LT:    PREC_LESSGREATER,
    GT:    PREC_LESSGREATER,
    LTEQ:  PREC_LESSGREATER,
    GTEQ:  PREC_LESSGREATER,
    PLUS:  PREC_SUM,
    MINUS: PREC_SUM,
    STAR:  PREC_PRODUCT,
    SLASH: PREC_PRODUCT,
    LPAREN: PREC_CALL,
}
```

The core loop in `parseExpression`:

```go
func (p *Parser) parseExpression(prec int) Expression {
    left := p.parsePrefix()   // parse the left-hand side / prefix

    for prec < p.peekPrecedence() {
        // the next token is an infix operator with higher precedence
        // consume it and parse the right-hand side
        left = p.parseInfix(left)
    }

    return left
}
```

`peekPrecedence()` returns the precedence of the *current* token (not consumed yet), defaulting to `PREC_LOWEST`.

## Prefix Parse Functions

Register a function for each token that can start an expression:

| Token | Handler |
|-------|---------|
| `NUMBER`, `STRING`, `TRUE`, `FALSE`, `IDENT` | literal / identifier (from Task 03) |
| `MINUS`, `BANG` | `parsePrefixExpr` → `PrefixExpr{Op, Right}` |
| `LPAREN` | `parseGroupedExpr` → consume `(`, parse expr, consume `)` |
| `FN` | `parseFnLiteral` → `FnLiteral{Params, Body}` |
| `IF` | `parseIfExpr` → `IfExpr{Cond, Then, Else?}` |

## Infix Parse Functions

| Token | Handler |
|-------|---------|
| `PLUS MINUS STAR SLASH EQEQ BANGEQ LT GT LTEQ GTEQ` | `parseInfixExpr` → `InfixExpr{Left, Op, Right}` |
| `LPAREN` | `parseCallExpr` → `CallExpr{Fn, Args}` |

### parseInfixExpr

```go
func (p *Parser) parseInfixExpr(left Expression) Expression {
    op := p.advance().Literal          // consume the operator
    prec := precedences[TokenType(op)] // its precedence
    right := p.parseExpression(prec)   // parse right at same level (left-assoc)
    return &InfixExpr{Left: left, Op: op, Right: right}
}
```

### parseFnLiteral

Grammar: `fn ( params ) { body }`

1. Consume `FN`.
2. Consume `(`.
3. Parse comma-separated `IDENT` tokens until `)`.
4. Parse a `BlockStatement`.

### parseIfExpr

Grammar: `if ( cond ) { then } [ else { else } ]`

1. Consume `IF`.
2. Consume `(`.
3. Parse the condition expression.
4. Consume `)`.
5. Parse `then` `BlockStatement`.
6. If the next token is `ELSE`, consume it and parse the `else` `BlockStatement`.

### parseBlockStatement

1. Consume `{`.
2. Parse statements until `}` or `EOF`.
3. Consume `}`.
4. Return `BlockStatement`.

### parseCallExpr

Grammar: `fn ( args )`

1. `fn` is the already-parsed `left` expression.
2. Consume `(`.
3. Parse comma-separated expressions until `)`.
4. Consume `)`.
5. Return `CallExpr{Fn: left, Args: args}`.

## Examples

| Input | Expected AST String |
|-------|---------------------|
| `1 + 2 * 3` | `(1 + (2 * 3))` |
| `(1 + 2) * 3` | `((1 + 2) * 3)` |
| `-5` | `(-5)` |
| `!true` | `(!true)` |
| `a == b` | `(a == b)` |
| `fn(x) { x }` | `fn(x) { x }` |
| `add(1, 2)` | `add(1, 2)` |
| `if (x < 2) { x } else { 2 }` | `if (x < 2) { x } else { 2 }` |

## Tests to Pass

1. `"1 + 2 * 3"` → String `"(1 + (2 * 3))"`
2. `"(1 + 2) * 3"` → String `"((1 + 2) * 3)"`
3. `"-5 + 3"` → String `"((-5) + 3)"`
4. `"!true"` → `PrefixExpr{Op:"!", Right: BoolLiteral(true)}`
5. `"a == b"` → `InfixExpr{Left: Identifier(a), Op:"==", Right: Identifier(b)}`
6. `"fn(a, b) { a + b }"` → `FnLiteral{Params:["a","b"], ...}`
7. `"add(1, 2 + 3)"` → `CallExpr` with 2 args, second is `InfixExpr`
8. `"if (x < 2) { x } else { 2 }"` → `IfExpr` with both branches
9. `"if (true) { 1 }"` → `IfExpr` with nil Else
