# Chapter 54: Building the Astra Lexer — Turning Text Into Tokens

> "Lexing is easy to get mostly right and hard to get completely right." — Every compiler writer, eventually

---

## Overview

The lexer is the front door of your compiler. It reads raw text — a sequence of characters — and produces a stream of **tokens**. A token is a meaningful unit: a keyword, a number, a string, an operator. By the time the parser runs, it never sees individual characters. It sees tokens.

Think of the lexer as a translator between two worlds:

```
Raw Text World                    Token World
─────────────────────────────────────────────────────
"fn add(a: int) -> int {"   →    [FN] [IDENT "add"] [LPAREN]
                                 [IDENT "a"] [COLON] [IDENT "int"]
                                 [RPAREN] [ARROW] [IDENT "int"]
                                 [LBRACE]
```

The lexer does not understand *meaning*. It does not know that `fn` starts a function. It only knows that `fn` is a keyword token called `FN`. Understanding meaning is the parser's job.

In this chapter we build the complete Astra lexer in Go. By the end, you will have two files — `lexer/token.go` and `lexer/lexer.go` — that correctly tokenize every valid (and invalid) Astra program.

---

## Table of Contents

1. What Is a Token?
2. Designing the Token Type System
3. The Token Struct
4. Building the Lexer Struct
5. The Main Scan Loop
6. Scanning Identifiers and Keywords
7. Scanning Numbers
8. Scanning Strings
9. Scanning Operators
10. Handling Whitespace and Comments
11. Error Handling
12. Complete Implementation
13. Testing the Lexer

---

## 1. What Is a Token?

A token has three key properties:

```
Token
├── Type    — what kind of thing is this? (NUMBER, STRING, PLUS, FN, ...)
├── Lexeme  — the exact text from the source ("42", "hello", "+", "fn")
└── Position — where in the source (line 5, column 12)
```

Some tokens also carry a **literal value** — the semantic value extracted from the lexeme. For example, the token for `42` has lexeme `"42"` (a string) and literal `42` (an integer).

```
Source:    let x = 42 + 0xFF
           ─── ─ ─ ── ─ ────
Tokens:    LET IDENT ASSIGN INT PLUS INT
Lexemes:   let  x    =      42  +    0xFF
Literals:  -    -    -      42  -    255
```

---

## 2. Designing the Token Type System

Every distinct kind of token needs a name. Here are all the token types for Astra:

```
ASCII diagram: Token Type Taxonomy

┌─────────────────────────────────────────────────────────────┐
│                       TOKEN TYPES                           │
│                                                             │
│  LITERALS           IDENTIFIERS     KEYWORDS                │
│  ─────────          ───────────     ────────                 │
│  INT_LIT            IDENT           FN, LET, CONST          │
│  FLOAT_LIT                          IF, ELSE                 │
│  STRING_LIT         OPERATORS       FOR, WHILE, IN           │
│  TRUE, FALSE        ─────────       RETURN                   │
│                     PLUS            STRUCT, IMPL             │
│  DELIMITERS         MINUS           IMPORT                   │
│  ──────────         STAR            PUB, SELF                │
│  LPAREN             SLASH           AS, MATCH                │
│  RPAREN             PERCENT         ENUM, TRAIT              │
│  LBRACE             BANG            BREAK, CONTINUE          │
│  RBRACE             ASSIGN                                   │
│  LBRACKET           EQ              SPECIAL                  │
│  RBRACKET           NEQ             ──────                   │
│  COMMA              LT, GT          EOF                      │
│  COLON              LTE, GTE        ILLEGAL                  │
│  SEMICOLON          AND, OR                                  │
│  DOT                ARROW                                    │
│  DOTDOT             COLON_COLON     │
└─────────────────────────────────────────────────────────────┘
```

---

## 3. Complete Implementation

Let's build the complete files. We start with `lexer/token.go`.

### `lexer/token.go`

```go
// lexer/token.go
// Defines all token types and the Token struct for the Astra language.

package lexer

import "fmt"

// TokenType is the type of a single lexical token.
type TokenType int

const (
	// ── Literals ──────────────────────────────────────────────────────────────
	INT_LIT    TokenType = iota // 42, 0xFF, 0b1010
	FLOAT_LIT                   // 3.14
	STRING_LIT                  // "hello"
	TRUE                        // true
	FALSE                       // false

	// ── Identifier ────────────────────────────────────────────────────────────
	IDENT // foo, bar, my_var

	// ── Keywords ──────────────────────────────────────────────────────────────
	FN
	LET
	CONST
	IF
	ELSE
	FOR
	WHILE
	IN
	RETURN
	STRUCT
	IMPL
	IMPORT
	PUB
	SELF
	AS
	MATCH
	ENUM
	TRAIT
	BREAK
	CONTINUE

	// ── Operators ─────────────────────────────────────────────────────────────
	PLUS    // +
	MINUS   // -
	STAR    // *
	SLASH   // /
	PERCENT // %
	BANG    // !

	ASSIGN // =
	EQ     // ==
	NEQ    // !=
	LT     // <
	GT     // >
	LTE    // <=
	GTE    // >=
	AND    // &&
	OR     // ||

	// ── Multi-char Operators ──────────────────────────────────────────────────
	ARROW       // ->
	DOTDOT      // ..
	COLON_COLON // ::

	// ── Delimiters ────────────────────────────────────────────────────────────
	LPAREN    // (
	RPAREN    // )
	LBRACE    // {
	RBRACE    // }
	LBRACKET  // [
	RBRACKET  // ]
	COMMA     // ,
	COLON     // :
	SEMICOLON // ;
	DOT       // .

	// ── Special ───────────────────────────────────────────────────────────────
	EOF     // end of file
	ILLEGAL // unexpected character
)

// tokenTypeNames maps TokenType to a human-readable string for debugging.
var tokenTypeNames = map[TokenType]string{
	INT_LIT: "INT_LIT", FLOAT_LIT: "FLOAT_LIT", STRING_LIT: "STRING_LIT",
	TRUE: "TRUE", FALSE: "FALSE", IDENT: "IDENT",
	FN: "FN", LET: "LET", CONST: "CONST",
	IF: "IF", ELSE: "ELSE",
	FOR: "FOR", WHILE: "WHILE", IN: "IN",
	RETURN: "RETURN", STRUCT: "STRUCT", IMPL: "IMPL",
	IMPORT: "IMPORT", PUB: "PUB", SELF: "SELF",
	AS: "AS", MATCH: "MATCH", ENUM: "ENUM", TRAIT: "TRAIT",
	BREAK: "BREAK", CONTINUE: "CONTINUE",
	PLUS: "PLUS", MINUS: "MINUS", STAR: "STAR",
	SLASH: "SLASH", PERCENT: "PERCENT", BANG: "BANG",
	ASSIGN: "ASSIGN", EQ: "EQ", NEQ: "NEQ",
	LT: "LT", GT: "GT", LTE: "LTE", GTE: "GTE",
	AND: "AND", OR: "OR",
	ARROW: "ARROW", DOTDOT: "DOTDOT", COLON_COLON: "COLON_COLON",
	LPAREN: "LPAREN", RPAREN: "RPAREN",
	LBRACE: "LBRACE", RBRACE: "RBRACE",
	LBRACKET: "LBRACKET", RBRACKET: "RBRACKET",
	COMMA: "COMMA", COLON: "COLON", SEMICOLON: "SEMICOLON", DOT: "DOT",
	EOF: "EOF", ILLEGAL: "ILLEGAL",
}

func (t TokenType) String() string {
	if s, ok := tokenTypeNames[t]; ok {
		return s
	}
	return fmt.Sprintf("TokenType(%d)", int(t))
}

// keywords maps keyword strings to their token types.
// Using a map gives O(1) lookup — faster than a long if-else chain.
var keywords = map[string]TokenType{
	"fn":       FN,
	"let":      LET,
	"const":    CONST,
	"if":       IF,
	"else":     ELSE,
	"for":      FOR,
	"while":    WHILE,
	"in":       IN,
	"return":   RETURN,
	"struct":   STRUCT,
	"impl":     IMPL,
	"import":   IMPORT,
	"pub":      PUB,
	"self":     SELF,
	"as":       AS,
	"match":    MATCH,
	"enum":     ENUM,
	"trait":    TRAIT,
	"break":    BREAK,
	"continue": CONTINUE,
	"true":     TRUE,
	"false":    FALSE,
}

// Position records where in the source a token appears.
type Position struct {
	Line   int // 1-based line number
	Column int // 1-based column number
	Offset int // byte offset from start of source
}

func (p Position) String() string {
	return fmt.Sprintf("%d:%d", p.Line, p.Column)
}

// Token is a single lexical unit produced by the lexer.
type Token struct {
	Type    TokenType
	Lexeme  string      // the exact source text
	Literal interface{} // the semantic value (int64 for INT_LIT, float64 for FLOAT_LIT, string for STRING_LIT)
	Pos     Position
}

func (t Token) String() string {
	if t.Literal != nil {
		return fmt.Sprintf("Token(%s, %q, %v, %s)", t.Type, t.Lexeme, t.Literal, t.Pos)
	}
	return fmt.Sprintf("Token(%s, %q, %s)", t.Type, t.Lexeme, t.Pos)
}
```

### `lexer/lexer.go`

Now the main lexer implementation:

```go
// lexer/lexer.go
// The Astra lexer: converts source text into a stream of tokens.

package lexer

import (
	"fmt"
	"strconv"
	"strings"
)

// LexError represents a lexical error with position information.
type LexError struct {
	Pos     Position
	Message string
}

func (e *LexError) Error() string {
	return fmt.Sprintf("lex error at %s: %s", e.Pos, e.Message)
}

// Lexer holds the state of the scanning process.
type Lexer struct {
	source  string    // the complete source text
	tokens  []Token   // tokens accumulated so far
	errors  []LexError // errors encountered during scanning

	start   int // byte index of start of current token
	current int // byte index of current character being examined
	line    int // current line number (1-based)
	lineStart int // byte index of start of current line (for column calc)
}

// New creates a new Lexer for the given source string.
func New(source string) *Lexer {
	return &Lexer{
		source:  source,
		tokens:  []Token{},
		errors:  []LexError{},
		start:   0,
		current: 0,
		line:    1,
		lineStart: 0,
	}
}

// ScanAll scans the entire source and returns all tokens.
// Errors are accumulated; call Errors() to retrieve them.
func (l *Lexer) ScanAll() []Token {
	for !l.isAtEnd() {
		l.start = l.current
		l.scanToken()
	}
	// Always append an EOF token
	l.tokens = append(l.tokens, Token{
		Type:   EOF,
		Lexeme: "",
		Pos:    l.currentPos(),
	})
	return l.tokens
}

// Errors returns all lexical errors encountered during scanning.
func (l *Lexer) Errors() []LexError {
	return l.errors
}

// ─── Core helpers ─────────────────────────────────────────────────────────────

// isAtEnd returns true if we've consumed all source characters.
func (l *Lexer) isAtEnd() bool {
	return l.current >= len(l.source)
}

// advance consumes and returns the current character.
func (l *Lexer) advance() byte {
	ch := l.source[l.current]
	l.current++
	return ch
}

// peek returns the current character without consuming it.
// Returns 0 if at end of source.
func (l *Lexer) peek() byte {
	if l.isAtEnd() {
		return 0
	}
	return l.source[l.current]
}

// peekNext returns the character after the current one without consuming.
// Returns 0 if there is no next character.
func (l *Lexer) peekNext() byte {
	if l.current+1 >= len(l.source) {
		return 0
	}
	return l.source[l.current+1]
}

// match consumes the current character only if it equals expected.
// Returns true if it matched (and was consumed), false otherwise.
func (l *Lexer) match(expected byte) bool {
	if l.isAtEnd() || l.source[l.current] != expected {
		return false
	}
	l.current++
	return true
}

// currentPos returns the Position of the start of the current token.
func (l *Lexer) currentPos() Position {
	return Position{
		Line:   l.line,
		Column: l.start - l.lineStart + 1,
		Offset: l.start,
	}
}

// currentLexeme returns the source text for the current token.
func (l *Lexer) currentLexeme() string {
	return l.source[l.start:l.current]
}

// addToken adds a token of the given type (with no literal).
func (l *Lexer) addToken(typ TokenType) {
	l.tokens = append(l.tokens, Token{
		Type:   typ,
		Lexeme: l.currentLexeme(),
		Pos:    l.currentPos(),
	})
}

// addTokenLit adds a token with a literal value.
func (l *Lexer) addTokenLit(typ TokenType, literal interface{}) {
	l.tokens = append(l.tokens, Token{
		Type:    typ,
		Lexeme:  l.currentLexeme(),
		Literal: literal,
		Pos:     l.currentPos(),
	})
}

// addError records a lexical error.
func (l *Lexer) addError(msg string) {
	l.errors = append(l.errors, LexError{
		Pos:     l.currentPos(),
		Message: msg,
	})
	l.addToken(ILLEGAL)
}

// ─── Main dispatch ─────────────────────────────────────────────────────────────

// scanToken reads the next token from the source.
// This is the main dispatch function — it looks at the current character
// and decides what kind of token to scan.
func (l *Lexer) scanToken() {
	ch := l.advance()

	switch ch {
	// ── Single-character tokens ──────────────────────────────────────────────
	case '(':
		l.addToken(LPAREN)
	case ')':
		l.addToken(RPAREN)
	case '{':
		l.addToken(LBRACE)
	case '}':
		l.addToken(RBRACE)
	case '[':
		l.addToken(LBRACKET)
	case ']':
		l.addToken(RBRACKET)
	case ',':
		l.addToken(COMMA)
	case ';':
		l.addToken(SEMICOLON)
	case '%':
		l.addToken(PERCENT)
	case '*':
		l.addToken(STAR)

	// ── Potentially two-character tokens ────────────────────────────────────
	case '+':
		l.addToken(PLUS)
	case '!':
		if l.match('=') {
			l.addToken(NEQ)
		} else {
			l.addToken(BANG)
		}
	case '=':
		if l.match('=') {
			l.addToken(EQ)
		} else {
			l.addToken(ASSIGN)
		}
	case '<':
		if l.match('=') {
			l.addToken(LTE)
		} else {
			l.addToken(LT)
		}
	case '>':
		if l.match('=') {
			l.addToken(GTE)
		} else {
			l.addToken(GT)
		}
	case '-':
		if l.match('>') {
			l.addToken(ARROW)
		} else {
			l.addToken(MINUS)
		}
	case '.':
		if l.match('.') {
			l.addToken(DOTDOT)
		} else {
			l.addToken(DOT)
		}
	case ':':
		if l.match(':') {
			l.addToken(COLON_COLON)
		} else {
			l.addToken(COLON)
		}
	case '&':
		if l.match('&') {
			l.addToken(AND)
		} else {
			l.addError(fmt.Sprintf("unexpected character '&'; did you mean '&&'?"))
		}
	case '|':
		if l.match('|') {
			l.addToken(OR)
		} else {
			l.addError(fmt.Sprintf("unexpected character '|'; did you mean '||'?"))
		}

	// ── Slash: division or comment ───────────────────────────────────────────
	case '/':
		if l.match('/') {
			l.scanLineComment()
		} else if l.match('*') {
			l.scanBlockComment()
		} else {
			l.addToken(SLASH)
		}

	// ── String literals ──────────────────────────────────────────────────────
	case '"':
		l.scanString()

	// ── Whitespace ───────────────────────────────────────────────────────────
	case ' ', '\r', '\t':
		// ignore whitespace — do nothing

	case '\n':
		l.line++
		l.lineStart = l.current

	// ── Numbers ──────────────────────────────────────────────────────────────
	default:
		if isDigit(ch) {
			l.scanNumber(ch)
		} else if isLetter(ch) {
			l.scanIdentifierOrKeyword()
		} else {
			l.addError(fmt.Sprintf("unexpected character: %q", ch))
		}
	}
}

// ─── Scanning helpers ──────────────────────────────────────────────────────────

// scanLineComment skips everything from // to end of line.
func (l *Lexer) scanLineComment() {
	for !l.isAtEnd() && l.peek() != '\n' {
		l.advance()
	}
	// Don't emit a token — comments are ignored
}

// scanBlockComment skips everything between /* and */.
func (l *Lexer) scanBlockComment() {
	startLine := l.line
	startCol := l.start - l.lineStart + 1
	for !l.isAtEnd() {
		if l.peek() == '*' && l.peekNext() == '/' {
			l.advance() // consume *
			l.advance() // consume /
			return
		}
		if l.peek() == '\n' {
			l.line++
			l.lineStart = l.current + 1
		}
		l.advance()
	}
	// If we get here, the block comment was never closed
	l.errors = append(l.errors, LexError{
		Pos:     Position{Line: startLine, Column: startCol},
		Message: "unterminated block comment",
	})
}

// scanString scans a string literal (after the opening " has been consumed).
func (l *Lexer) scanString() {
	var sb strings.Builder
	for !l.isAtEnd() && l.peek() != '"' {
		ch := l.advance()
		if ch == '\n' {
			// Strings cannot span newlines in Astra
			l.addError("unterminated string literal: newline in string")
			return
		}
		if ch == '\\' {
			// Process escape sequence
			if l.isAtEnd() {
				l.addError("unterminated escape sequence at end of file")
				return
			}
			esc := l.advance()
			switch esc {
			case 'n':
				sb.WriteByte('\n')
			case 't':
				sb.WriteByte('\t')
			case 'r':
				sb.WriteByte('\r')
			case '\\':
				sb.WriteByte('\\')
			case '"':
				sb.WriteByte('"')
			case '0':
				sb.WriteByte(0)
			default:
				l.addError(fmt.Sprintf("unknown escape sequence: \\%c", esc))
				return
			}
		} else {
			sb.WriteByte(ch)
		}
	}

	if l.isAtEnd() {
		l.addError("unterminated string literal: reached end of file")
		return
	}

	l.advance() // consume the closing "
	l.addTokenLit(STRING_LIT, sb.String())
}

// scanNumber scans an integer or float literal.
// ch is the first digit already consumed.
func (l *Lexer) scanNumber(ch byte) {
	// Check for hex (0x...) or binary (0b...) prefix
	if ch == '0' {
		if l.peek() == 'x' || l.peek() == 'X' {
			l.advance() // consume 'x'
			l.scanHex()
			return
		}
		if l.peek() == 'b' || l.peek() == 'B' {
			l.advance() // consume 'b'
			l.scanBinary()
			return
		}
	}

	// Decimal integer or float
	for !l.isAtEnd() && (isDigit(l.peek()) || l.peek() == '_') {
		l.advance()
	}

	// Check for float: digits followed by '.' followed by more digits
	if l.peek() == '.' && isDigit(l.peekNext()) {
		l.advance() // consume '.'
		for !l.isAtEnd() && (isDigit(l.peek()) || l.peek() == '_') {
			l.advance()
		}
		// Parse the float
		lexeme := strings.ReplaceAll(l.currentLexeme(), "_", "")
		val, err := strconv.ParseFloat(lexeme, 64)
		if err != nil {
			l.addError(fmt.Sprintf("invalid float literal: %s", l.currentLexeme()))
			return
		}
		l.addTokenLit(FLOAT_LIT, val)
		return
	}

	// Parse the integer (strip underscores first)
	lexeme := strings.ReplaceAll(l.currentLexeme(), "_", "")
	val, err := strconv.ParseInt(lexeme, 10, 64)
	if err != nil {
		l.addError(fmt.Sprintf("invalid integer literal: %s", l.currentLexeme()))
		return
	}
	l.addTokenLit(INT_LIT, val)
}

// scanHex scans a hexadecimal literal after "0x" has been consumed.
func (l *Lexer) scanHex() {
	if !isHexDigit(l.peek()) {
		l.addError("expected hex digit after '0x'")
		return
	}
	for !l.isAtEnd() && (isHexDigit(l.peek()) || l.peek() == '_') {
		l.advance()
	}
	// Strip "0x" prefix and underscores, then parse
	raw := l.currentLexeme()[2:] // remove "0x"
	raw = strings.ReplaceAll(raw, "_", "")
	val, err := strconv.ParseInt(raw, 16, 64)
	if err != nil {
		l.addError(fmt.Sprintf("invalid hex literal: %s", l.currentLexeme()))
		return
	}
	l.addTokenLit(INT_LIT, val)
}

// scanBinary scans a binary literal after "0b" has been consumed.
func (l *Lexer) scanBinary() {
	if l.peek() != '0' && l.peek() != '1' {
		l.addError("expected binary digit (0 or 1) after '0b'")
		return
	}
	for !l.isAtEnd() && (l.peek() == '0' || l.peek() == '1' || l.peek() == '_') {
		l.advance()
	}
	raw := l.currentLexeme()[2:] // remove "0b"
	raw = strings.ReplaceAll(raw, "_", "")
	val, err := strconv.ParseInt(raw, 2, 64)
	if err != nil {
		l.addError(fmt.Sprintf("invalid binary literal: %s", l.currentLexeme()))
		return
	}
	l.addTokenLit(INT_LIT, val)
}

// scanIdentifierOrKeyword scans an identifier or keyword.
// The first character has already been consumed.
func (l *Lexer) scanIdentifierOrKeyword() {
	for !l.isAtEnd() && (isLetter(l.peek()) || isDigit(l.peek())) {
		l.advance()
	}

	word := l.currentLexeme()

	// Check if it's a keyword
	if tokType, isKeyword := keywords[word]; isKeyword {
		// true and false carry literal values
		if tokType == TRUE {
			l.addTokenLit(TRUE, true)
		} else if tokType == FALSE {
			l.addTokenLit(FALSE, false)
		} else {
			l.addToken(tokType)
		}
		return
	}

	// Otherwise it's a plain identifier
	l.addToken(IDENT)
}

// ─── Character classification helpers ─────────────────────────────────────────

func isDigit(ch byte) bool {
	return ch >= '0' && ch <= '9'
}

func isHexDigit(ch byte) bool {
	return (ch >= '0' && ch <= '9') ||
		(ch >= 'a' && ch <= 'f') ||
		(ch >= 'A' && ch <= 'F')
}

func isLetter(ch byte) bool {
	return (ch >= 'a' && ch <= 'z') ||
		(ch >= 'A' && ch <= 'Z') ||
		ch == '_'
}
```

---

## 4. Testing the Lexer

Go has excellent built-in testing. Here is the complete test file with ten table-driven test cases:

```go
// lexer/lexer_test.go
package lexer

import (
	"testing"
)

// tokenResult is a simplified token representation for test assertions.
type tokenResult struct {
	Type    TokenType
	Lexeme  string
	Literal interface{}
}

func tok(typ TokenType, lexeme string) tokenResult {
	return tokenResult{Type: typ, Lexeme: lexeme}
}

func tokLit(typ TokenType, lexeme string, lit interface{}) tokenResult {
	return tokenResult{Type: typ, Lexeme: lexeme, Literal: lit}
}

func TestLexer(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []tokenResult
		errCount int
	}{
		{
			// ── Test 1: Simple keywords ──────────────────────────────────────
			name:  "keywords",
			input: "fn let const if else for while in return struct impl",
			expected: []tokenResult{
				tok(FN, "fn"), tok(LET, "let"), tok(CONST, "const"),
				tok(IF, "if"), tok(ELSE, "else"),
				tok(FOR, "for"), tok(WHILE, "while"), tok(IN, "in"),
				tok(RETURN, "return"), tok(STRUCT, "struct"), tok(IMPL, "impl"),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 2: Integer literals (decimal, hex, binary) ───────────
			name:  "integer_literals",
			input: "42 0xFF 0b1010 1_000_000",
			expected: []tokenResult{
				tokLit(INT_LIT, "42", int64(42)),
				tokLit(INT_LIT, "0xFF", int64(255)),
				tokLit(INT_LIT, "0b1010", int64(10)),
				tokLit(INT_LIT, "1_000_000", int64(1000000)),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 3: Float literals ────────────────────────────────────
			name:  "float_literals",
			input: "3.14 2.718 0.001",
			expected: []tokenResult{
				tokLit(FLOAT_LIT, "3.14", float64(3.14)),
				tokLit(FLOAT_LIT, "2.718", float64(2.718)),
				tokLit(FLOAT_LIT, "0.001", float64(0.001)),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 4: String literals with escapes ──────────────────────
			name:  "string_literals",
			input: `"hello" "world\n" "tab\there" "quote\"here"`,
			expected: []tokenResult{
				tokLit(STRING_LIT, `"hello"`, "hello"),
				tokLit(STRING_LIT, `"world\n"`, "world\n"),
				tokLit(STRING_LIT, `"tab\there"`, "tab\there"),
				tokLit(STRING_LIT, `"quote\"here"`, `quote"here`),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 5: Operators (single and multi-char) ──────────────────
			name:  "operators",
			input: "+ - * / % == != <= >= && || -> .. :: !",
			expected: []tokenResult{
				tok(PLUS, "+"), tok(MINUS, "-"), tok(STAR, "*"),
				tok(SLASH, "/"), tok(PERCENT, "%"),
				tok(EQ, "=="), tok(NEQ, "!="),
				tok(LTE, "<="), tok(GTE, ">="),
				tok(AND, "&&"), tok(OR, "||"),
				tok(ARROW, "->"), tok(DOTDOT, ".."), tok(COLON_COLON, "::"),
				tok(BANG, "!"),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 6: A complete function declaration ───────────────────
			name:  "function_declaration",
			input: "fn add(a: int, b: int) -> int { return a + b }",
			expected: []tokenResult{
				tok(FN, "fn"), tok(IDENT, "add"),
				tok(LPAREN, "("),
				tok(IDENT, "a"), tok(COLON, ":"), tok(IDENT, "int"),
				tok(COMMA, ","),
				tok(IDENT, "b"), tok(COLON, ":"), tok(IDENT, "int"),
				tok(RPAREN, ")"),
				tok(ARROW, "->"), tok(IDENT, "int"),
				tok(LBRACE, "{"),
				tok(RETURN, "return"), tok(IDENT, "a"), tok(PLUS, "+"), tok(IDENT, "b"),
				tok(RBRACE, "}"),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 7: Comments are skipped ──────────────────────────────
			name: "comments",
			input: `// this is a comment
42 /* block
comment */ 100`,
			expected: []tokenResult{
				tokLit(INT_LIT, "42", int64(42)),
				tokLit(INT_LIT, "100", int64(100)),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 8: Boolean literals ──────────────────────────────────
			name:  "booleans",
			input: "true false",
			expected: []tokenResult{
				tokLit(TRUE, "true", true),
				tokLit(FALSE, "false", false),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 9: Struct definition ─────────────────────────────────
			name: "struct_definition",
			input: `struct Person {
    name: string
    age:  int
}`,
			expected: []tokenResult{
				tok(STRUCT, "struct"), tok(IDENT, "Person"),
				tok(LBRACE, "{"),
				tok(IDENT, "name"), tok(COLON, ":"), tok(IDENT, "string"),
				tok(IDENT, "age"), tok(COLON, ":"), tok(IDENT, "int"),
				tok(RBRACE, "}"),
				tok(EOF, ""),
			},
		},
		{
			// ── Test 10: Error recovery on illegal characters ─────────────
			name:     "error_recovery",
			input:    "let x = @",
			errCount: 1,
			expected: []tokenResult{
				tok(LET, "let"),
				tok(IDENT, "x"),
				tok(ASSIGN, "="),
				tok(ILLEGAL, "@"),
				tok(EOF, ""),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			l := New(tt.input)
			tokens := l.ScanAll()
			errors := l.Errors()

			// Check error count
			if len(errors) != tt.errCount {
				t.Errorf("expected %d errors, got %d: %v", tt.errCount, len(errors), errors)
			}

			// Check token count
			if len(tokens) != len(tt.expected) {
				t.Errorf("expected %d tokens, got %d", len(tt.expected), len(tokens))
				for i, tok := range tokens {
					t.Logf("  [%d] %s", i, tok)
				}
				return
			}

			// Check each token
			for i, want := range tt.expected {
				got := tokens[i]
				if got.Type != want.Type {
					t.Errorf("token[%d]: expected type %s, got %s", i, want.Type, got.Type)
				}
				if got.Lexeme != want.Lexeme {
					t.Errorf("token[%d]: expected lexeme %q, got %q", i, want.Lexeme, got.Lexeme)
				}
				if want.Literal != nil && got.Literal != want.Literal {
					t.Errorf("token[%d]: expected literal %v (%T), got %v (%T)",
						i, want.Literal, want.Literal, got.Literal, got.Literal)
				}
			}
		})
	}
}

// TestLineNumbers verifies that the lexer correctly tracks line and column.
func TestLineNumbers(t *testing.T) {
	input := "fn main() {\n    let x = 42\n}"
	l := New(input)
	tokens := l.ScanAll()

	// "let" should be on line 2
	for _, tok := range tokens {
		if tok.Type == LET {
			if tok.Pos.Line != 2 {
				t.Errorf("expected LET on line 2, got line %d", tok.Pos.Line)
			}
			if tok.Pos.Column != 5 {
				t.Errorf("expected LET at column 5, got column %d", tok.Pos.Column)
			}
			return
		}
	}
	t.Error("LET token not found")
}

// TestUnterminatedString verifies error handling for unclosed strings.
func TestUnterminatedString(t *testing.T) {
	l := New(`"hello world`)
	l.ScanAll()
	if len(l.Errors()) == 0 {
		t.Error("expected error for unterminated string, got none")
	}
}
```

---

## 5. How the Lexer Processes Input

Here is a step-by-step trace of scanning `fn add(a: int) -> int`:

```
ASCII Diagram: Lexer State Machine Trace

Source:  f n   a d d ( a :   i n t )   - >   i n t
         ↑
         current=0, start=0

Step 1: advance() → 'f'
        isLetter('f') → true
        scanIdentifierOrKeyword()
        consume n, space not letter → stop
        word = "fn" → keyword FN
        emit Token{FN, "fn", pos{1,1}}
         
         start=0, current=2
        
Step 2: advance() → ' ' (space)
        whitespace → skip
        
Step 3: advance() → 'a'
        isLetter → scanIdentifierOrKeyword()
        consume d, d → "add" not in keywords → IDENT
        emit Token{IDENT, "add", pos{1,4}}

Step 4: advance() → '('
        → emit Token{LPAREN, "(", pos{1,7}}

... and so on

Final token stream:
  Token{FN,     "fn",  pos{1,1}}
  Token{IDENT,  "add", pos{1,4}}
  Token{LPAREN, "(",   pos{1,7}}
  Token{IDENT,  "a",   pos{1,8}}
  Token{COLON,  ":",   pos{1,9}}
  Token{IDENT,  "int", pos{1,11}}
  Token{RPAREN, ")",   pos{1,14}}
  Token{ARROW,  "->",  pos{1,16}}
  Token{IDENT,  "int", pos{1,19}}
  Token{EOF,    "",    pos{1,22}}
```

---

## 6. The Keyword Lookup Strategy

When the lexer encounters a sequence of letters, it scans as many as possible (the "maximal munch" rule) and then checks whether the resulting word is a keyword.

```
ASCII Diagram: Identifier vs Keyword Resolution

"foobar" → scan all letters → "foobar" → NOT in keywords → IDENT

"fn"     → scan all letters → "fn" → IS in keywords → FN

"format" → scan all letters → "format" → NOT in keywords → IDENT
                                         (even though it starts with "for")

"for"    → scan all letters → "for" → IS in keywords → FOR
```

The maximal munch rule means we always try to consume as much as possible before deciding what a token is. This is why `format` does not produce `FOR` followed by `mat` — we consume all of `format` before checking.

---

## Summary Table

| Component | File | Lines | Purpose |
|---|---|---|---|
| TokenType enum | `lexer/token.go` | 1-80 | All token type names |
| Keywords map | `lexer/token.go` | 81-110 | Keyword lookup table |
| Position struct | `lexer/token.go` | 111-125 | Source location tracking |
| Token struct | `lexer/token.go` | 126-145 | A single token |
| Lexer struct | `lexer/lexer.go` | 1-50 | Scanning state |
| scanToken() | `lexer/lexer.go` | 51-130 | Main dispatch |
| scanString() | `lexer/lexer.go` | 131-175 | String literals |
| scanNumber() | `lexer/lexer.go` | 176-240 | Int/float literals |
| scanIdentifier() | `lexer/lexer.go` | 241-265 | Identifiers/keywords |
| Comment scanners | `lexer/lexer.go` | 266-295 | // and /* */ |
| Tests | `lexer/lexer_test.go` | 300 | 10 test cases |

---

## Exercises

1. **Multiline Strings**: Add support for multiline string literals using backtick delimiters (like Go's raw strings). When the lexer sees `` ` ``, it should consume everything until the next `` ` ``, including newlines. No escape sequences inside raw strings. Write the `scanRawString()` function.

2. **Better Error Recovery**: Currently, when the lexer encounters an illegal character, it emits one ILLEGAL token and moves on. Improve error recovery so that if there are multiple consecutive illegal characters, they are grouped into one ILLEGAL token rather than one per character.

3. **Column Tracking**: The current column tracking resets at each newline. Verify this works correctly by writing a test where a token appears at column > 40 on the third line of input.

4. **Number Validation**: The current `scanNumber` does not check for numbers that are too large for int64 (overflow). What happens when you lex `99999999999999999999`? Fix the error message to be clear about integer overflow.

5. **Unicode Identifiers**: Modify `isLetter()` to accept Unicode letters (rune values > 127). This requires changing the lexer to work with `rune` (int32) instead of `byte`. What complications arise with `peekNext()` when characters are multi-byte?

6. **Benchmark**: Write a Go benchmark (`func BenchmarkLexer(b *testing.B)`) that measures how many tokens per second the lexer can process. Generate a 10,000-line Astra source file (you can repeat the `fn add` example) and benchmark it.

7. **Trie-based Keyword Lookup**: The current implementation uses a hash map for keyword lookup. Implement a trie-based keyword lookup and benchmark it against the hash map. Which is faster for Astra's keyword set? Why?

8. **Lexer as Iterator**: The current API scans everything at once (`ScanAll()`). Implement a `Next() Token` method that returns one token at a time (lazy evaluation). The parser chapter uses this API. What state does the iterator need to maintain?
