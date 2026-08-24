package main

import "fmt"

// TokenType identifies the kind of a token.
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

// Token is a single lexical unit.
type Token struct {
	Type    TokenType
	Literal string
	Line    int
}

// keywords maps reserved words to their token types.
var keywords = map[string]TokenType{
	"let":    LET,
	"fn":     FN,
	"if":     IF,
	"else":   ELSE,
	"return": RETURN,
	"while":  WHILE,
	"true":   TRUE,
	"false":  FALSE,
}

// Tokenize turns the source string into a slice of tokens.
// The last token is always EOF.
func Tokenize(input string) []Token {
	// TODO: implement lexical analysis
	//
	// Suggested approach:
	//   pos  := 0        (current read position)
	//   line := 1        (current line number, starts at 1)
	//
	// Loop while pos < len(input):
	//   1. Skip whitespace. Increment line when you see '\n'.
	//   2. Skip comments: if current char is '/' and next is '/', advance
	//      past the end of the line.
	//   3. Try two-character tokens first (==, !=, <=, >=).
	//   4. Match single-character tokens.
	//   5. Scan numbers: digits (and optional '.' + more digits).
	//   6. Scan strings: open quote, content, close quote.
	//      Store the content WITHOUT the surrounding quotes.
	//   7. Scan identifiers / keywords: [a-zA-Z_][a-zA-Z0-9_]*
	//      Then check keywords map; use IDENT if not found.
	//
	// Append EOF at the end.

	var tokens []Token
	_ = tokens // remove when you start using it

	return append(tokens, Token{Type: EOF, Line: 0})
}

func main() {
	src := `let x = 10 + 2.5`
	tokens := Tokenize(src)
	for _, t := range tokens {
		fmt.Printf("%s(%q) line=%d\n", t.Type, t.Literal, t.Line)
	}
}
