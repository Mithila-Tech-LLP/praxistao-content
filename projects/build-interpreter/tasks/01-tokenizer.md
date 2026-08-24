# Task 01 — Tokenizer

## Concept

**Lexical analysis** is the first stage of any interpreter. Its job is simple: read a raw string of characters and group them into meaningful chunks called **tokens**. A token is the smallest unit of meaning in a language — a number, a keyword, an operator, an identifier.

Think of it like reading a sentence. "let x = 10" is not four characters separated by spaces — it is four distinct tokens: a keyword (`let`), a name (`x`), an operator (`=`), and a number (`10`). The tokenizer gives every subsequent stage clean, typed input to work with instead of raw text.

## Token Types

Your tokenizer must recognise the following token types:

| Group | Types |
|-------|-------|
| Literals | `NUMBER`, `STRING`, `TRUE`, `FALSE` |
| Identifier | `IDENT` |
| Arithmetic | `PLUS`, `MINUS`, `STAR`, `SLASH` |
| Comparison | `EQEQ`, `BANGEQ`, `LT`, `GT`, `LTEQ`, `GTEQ` |
| Assignment | `EQ` |
| Logical prefix | `BANG` |
| Delimiters | `LPAREN`, `RPAREN`, `LBRACE`, `RBRACE`, `COMMA`, `SEMICOLON` |
| Keywords | `LET`, `FN`, `IF`, `ELSE`, `RETURN`, `WHILE` |
| Sentinel | `EOF` |

## Data Structures

```go
type TokenType string

const (
    NUMBER    TokenType = "NUMBER"
    STRING    TokenType = "STRING"
    IDENT     TokenType = "IDENT"
    PLUS      TokenType = "+"
    MINUS     TokenType = "-"
    STAR      TokenType = "*"
    SLASH     TokenType = "/"
    EQ        TokenType = "="
    EQEQ      TokenType = "=="
    BANG      TokenType = "!"
    BANGEQ    TokenType = "!="
    LT        TokenType = "<"
    GT        TokenType = ">"
    LTEQ      TokenType = "<="
    GTEQ      TokenType = ">="
    LPAREN    TokenType = "("
    RPAREN    TokenType = ")"
    LBRACE    TokenType = "{"
    RBRACE    TokenType = "}"
    COMMA     TokenType = ","
    SEMICOLON TokenType = ";"
    LET       TokenType = "let"
    FN        TokenType = "fn"
    IF        TokenType = "if"
    ELSE      TokenType = "else"
    RETURN    TokenType = "return"
    WHILE     TokenType = "while"
    TRUE      TokenType = "true"
    FALSE     TokenType = "false"
    EOF       TokenType = "EOF"
)

type Token struct {
    Type    TokenType
    Literal string
    Line    int
}
```

## Your Task

Implement:

```go
func Tokenize(input string) []Token
```

Rules:
- Skip whitespace (spaces, tabs, `\r`, `\n`). Increment the line counter when you see `\n`.
- Skip comments: `//` causes the tokenizer to skip everything until the end of that line.
- Numbers can be integers or decimals: `42`, `3.14`. Scan digits, optionally a `.` followed by more digits.
- Strings are delimited by double quotes: `"hello world"`. The literal stored in the token should not include the quotes.
- Identifiers are `[a-zA-Z_][a-zA-Z0-9_]*`. After scanning, check if the identifier is a reserved keyword and use the keyword token type if so.
- Two-character operators (`==`, `!=`, `<=`, `>=`) must be matched before their single-character prefixes.
- Always append an `EOF` token at the end.

## Example

Input:
```
let x = 10 + 2.5
```

Expected output (Literal values shown):
```
LET("let") IDENT("x") EQ("=") NUMBER("10") PLUS("+") NUMBER("2.5") EOF("")
```

Input with a comment:
```
let y = 5 // this is ignored
let z = y
```

Expected: tokens for both `let` statements, nothing from the comment.

## Hints

- Use an index variable `pos int` and advance through `input` character by character.
- A helper `peek() byte` that looks one character ahead (without consuming) is very useful for two-character tokens.
- Map keywords at the end of identifier scanning: `keywords := map[string]TokenType{"let": LET, "fn": FN, ...}`.

## Tests to Pass

1. Single token types: each operator and delimiter individually.
2. Keywords: `let`, `fn`, `if`, `else`, `return`, `while`, `true`, `false` are not scanned as `IDENT`.
3. Numbers: integers and decimals both produce `NUMBER` tokens.
4. Strings: `"hello"` produces `STRING` with Literal `hello` (no quotes).
5. Comments: characters after `//` up to `\n` are skipped entirely.
6. Line tracking: the `Line` field increments for each `\n`.
7. A full expression: `let x = 10 + 2.5` produces the 7 tokens shown above.
