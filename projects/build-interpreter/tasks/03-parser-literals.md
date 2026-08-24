# Task 03 — Parser: Literals and Identifiers

## Concept

The **parser** consumes the flat list of tokens produced by the tokenizer and builds the AST. This task covers the simplest parsing: literals (numbers, strings, booleans), identifiers, `let` statements, and `return` statements. These have no operator precedence or recursion challenges — they are one token (or a short fixed sequence).

The technique used here and in Task 04 is **recursive descent parsing**, specifically the **Pratt parser** variant for expressions. Pratt parsing assigns a *precedence level* and a *parse function* to each token type, making operator precedence elegant without special-casing.

## Parser Structure

```go
type Parser struct {
    tokens  []Token
    pos     int
}

func NewParser(tokens []Token) *Parser

// Advance and return the current token, moving pos forward.
func (p *Parser) advance() Token

// Return the current token without consuming it.
func (p *Parser) current() Token

// Return the next token without consuming it.
func (p *Parser) peek() Token

// If the current token matches t, consume it; otherwise record an error.
func (p *Parser) expect(t TokenType) Token

func (p *Parser) ParseProgram() *Program
```

## Parsing Statements

`ParseProgram` repeatedly calls `parseStatement` until `EOF`.

`parseStatement` dispatches based on the current token:
- `LET` → `parseLetStatement`
- `RETURN` → `parseReturnStatement`
- anything else → `parseExpressionStatement`

### Let Statement

Grammar: `let <IDENT> = <expression> ;`

```go
func (p *Parser) parseLetStatement() *LetStatement
```

1. Consume `LET`.
2. Expect and consume an `IDENT` token — record its literal as `Name`.
3. Expect and consume `EQ`.
4. Parse the right-hand expression (`parseExpression`).
5. Optionally consume a trailing `SEMICOLON` if present.

### Return Statement

Grammar: `return <expression> ;`

1. Consume `RETURN`.
2. Parse the expression.
3. Optionally consume `SEMICOLON`.

### Expression Statement

Grammar: `<expression> ;`

1. Parse the expression.
2. Optionally consume `SEMICOLON`.
3. Wrap in `ExpressionStatement`.

## Parsing Expressions (this task: literals only)

For now implement only literal parsing in `parseExpression`. You will extend it in Task 04.

```go
func (p *Parser) parseExpression(precedence int) Expression
```

Handle:
- Current token is `NUMBER` → return `NumberLiteral`
- Current token is `STRING` → return `StringLiteral`
- Current token is `TRUE` or `FALSE` → return `BoolLiteral`
- Current token is `IDENT` → return `Identifier`

For any other token, return `nil` (or a parse error).

## Error Handling

Keep a slice of error messages: `errors []string`. When `expect` fails, append a message like `"expected = but got +"`. After parsing, callers can check `p.errors`. For now do not stop parsing on errors — keep going so tests can see partial output.

## Example

Input: `let x = 42`

Tokens: `LET IDENT("x") EQ NUMBER("42") EOF`

Expected AST:
```
Program{
  Statements: [
    LetStatement{Name: "x", Value: NumberLiteral{Value: 42}},
  ],
}
```

`prog.String()` → `"let x = 42;"`

## Tests to Pass

1. `"let x = 42"` → `LetStatement{Name:"x", Value: NumberLiteral(42)}`
2. `"let msg = \"hello\""` → `LetStatement{Name:"msg", Value: StringLiteral("hello")}`
3. `"let flag = true"` → `LetStatement{Name:"flag", Value: BoolLiteral(true)}`
4. `"return 99"` → `ReturnStatement{Value: NumberLiteral(99)}`
5. `"x"` (bare identifier) → `ExpressionStatement{Expr: Identifier{Name:"x"}}`
6. Multiple statements separated by semicolons all parsed into `Program.Statements`.
7. `p.Errors()` is empty for all valid inputs above.
