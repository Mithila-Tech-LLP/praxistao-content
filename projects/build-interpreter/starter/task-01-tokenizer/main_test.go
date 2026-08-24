package main

import "testing"

func tokenTypes(tokens []Token) []TokenType {
	tt := make([]TokenType, len(tokens))
	for i, t := range tokens {
		tt[i] = t.Type
	}
	return tt
}

func assertTypes(t *testing.T, tokens []Token, expected []TokenType) {
	t.Helper()
	types := tokenTypes(tokens)
	if len(types) != len(expected) {
		t.Fatalf("token count: expected %d (%v), got %d (%v)", len(expected), expected, len(types), types)
	}
	for i, et := range expected {
		if types[i] != et {
			t.Errorf("token[%d]: expected %s, got %s", i, et, types[i])
		}
	}
}

func TestTokenize_LetAssign(t *testing.T) {
	tokens := Tokenize("let x = 42")
	assertTypes(t, tokens, []TokenType{LET, IDENT, EQ, NUMBER, EOF})

	if tokens[1].Literal != "x" {
		t.Errorf("IDENT literal: expected x, got %q", tokens[1].Literal)
	}
	if tokens[3].Literal != "42" {
		t.Errorf("NUMBER literal: expected 42, got %q", tokens[3].Literal)
	}
}

func TestTokenize_FnExpression(t *testing.T) {
	tokens := Tokenize("fn(a, b) { a + b }")
	assertTypes(t, tokens, []TokenType{
		FN, LPAREN, IDENT, COMMA, IDENT, RPAREN, LBRACE, IDENT, PLUS, IDENT, RBRACE, EOF,
	})
}

func TestTokenize_EqualityOp(t *testing.T) {
	tokens := Tokenize("1 == 2")
	assertTypes(t, tokens, []TokenType{NUMBER, EQEQ, NUMBER, EOF})
}

func TestTokenize_MultiCharOps(t *testing.T) {
	tokens := Tokenize("!= <= >=")
	assertTypes(t, tokens, []TokenType{BANGEQ, LTEQ, GTEQ, EOF})
}

func TestTokenize_Comments(t *testing.T) {
	tokens := Tokenize("// hello\n1")
	assertTypes(t, tokens, []TokenType{NUMBER, EOF})
	if tokens[0].Literal != "1" {
		t.Errorf("expected literal 1 after comment, got %q", tokens[0].Literal)
	}
}

func TestTokenize_Keywords(t *testing.T) {
	tokens := Tokenize("if true else false return while")
	assertTypes(t, tokens, []TokenType{IF, TRUE, ELSE, FALSE, RETURN, WHILE, EOF})
}

func TestTokenize_StringLiteral(t *testing.T) {
	tokens := Tokenize(`"hello world"`)
	assertTypes(t, tokens, []TokenType{STRING, EOF})
	if tokens[0].Literal != "hello world" {
		t.Errorf("string literal: expected 'hello world', got %q", tokens[0].Literal)
	}
}

func TestTokenize_FloatNumber(t *testing.T) {
	tokens := Tokenize("3.14")
	assertTypes(t, tokens, []TokenType{NUMBER, EOF})
	if tokens[0].Literal != "3.14" {
		t.Errorf("float literal: expected 3.14, got %q", tokens[0].Literal)
	}
}
