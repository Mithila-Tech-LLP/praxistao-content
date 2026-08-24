package main

import (
	"fmt"
	"strconv"
	"strings"
)

// ============================================================
// Tokenizer
// ============================================================

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

var keywords = map[string]TokenType{
	"let": LET, "fn": FN, "if": IF, "else": ELSE,
	"return": RETURN, "while": WHILE, "true": TRUE, "false": FALSE,
}

func Tokenize(input string) []Token {
	var tokens []Token
	pos, line := 0, 1

	for pos < len(input) {
		ch := input[pos]
		if ch == ' ' || ch == '\t' || ch == '\r' {
			pos++
			continue
		}
		if ch == '\n' {
			line++
			pos++
			continue
		}
		if pos+1 < len(input) && ch == '/' && input[pos+1] == '/' {
			for pos < len(input) && input[pos] != '\n' {
				pos++
			}
			continue
		}
		if pos+1 < len(input) {
			two := string(input[pos : pos+2])
			switch two {
			case "==":
				tokens = append(tokens, Token{EQEQ, two, line})
				pos += 2
				continue
			case "!=":
				tokens = append(tokens, Token{BANGEQ, two, line})
				pos += 2
				continue
			case "<=":
				tokens = append(tokens, Token{LTEQ, two, line})
				pos += 2
				continue
			case ">=":
				tokens = append(tokens, Token{GTEQ, two, line})
				pos += 2
				continue
			}
		}
		switch ch {
		case '+':
			tokens = append(tokens, Token{PLUS, "+", line})
		case '-':
			tokens = append(tokens, Token{MINUS, "-", line})
		case '*':
			tokens = append(tokens, Token{STAR, "*", line})
		case '/':
			tokens = append(tokens, Token{SLASH, "/", line})
		case '=':
			tokens = append(tokens, Token{EQ, "=", line})
		case '!':
			tokens = append(tokens, Token{BANG, "!", line})
		case '<':
			tokens = append(tokens, Token{LT, "<", line})
		case '>':
			tokens = append(tokens, Token{GT, ">", line})
		case '(':
			tokens = append(tokens, Token{LPAREN, "(", line})
		case ')':
			tokens = append(tokens, Token{RPAREN, ")", line})
		case '{':
			tokens = append(tokens, Token{LBRACE, "{", line})
		case '}':
			tokens = append(tokens, Token{RBRACE, "}", line})
		case ',':
			tokens = append(tokens, Token{COMMA, ",", line})
		case ';':
			tokens = append(tokens, Token{SEMICOLON, ";", line})
		default:
			if ch >= '0' && ch <= '9' {
				start := pos
				for pos < len(input) && input[pos] >= '0' && input[pos] <= '9' {
					pos++
				}
				if pos < len(input) && input[pos] == '.' {
					pos++
					for pos < len(input) && input[pos] >= '0' && input[pos] <= '9' {
						pos++
					}
				}
				tokens = append(tokens, Token{NUMBER, input[start:pos], line})
				continue
			}
			if ch == '"' {
				pos++
				start := pos
				for pos < len(input) && input[pos] != '"' {
					pos++
				}
				tokens = append(tokens, Token{STRING, input[start:pos], line})
				pos++ // closing quote
				continue
			}
			if ch == '_' || (ch >= 'a' && ch <= 'z') || (ch >= 'A' && ch <= 'Z') {
				start := pos
				for pos < len(input) && (input[pos] == '_' || (input[pos] >= 'a' && input[pos] <= 'z') || (input[pos] >= 'A' && input[pos] <= 'Z') || (input[pos] >= '0' && input[pos] <= '9')) {
					pos++
				}
				word := input[start:pos]
				if tt, ok := keywords[word]; ok {
					tokens = append(tokens, Token{tt, word, line})
				} else {
					tokens = append(tokens, Token{IDENT, word, line})
				}
				continue
			}
		}
		pos++
	}
	return append(tokens, Token{Type: EOF, Line: line})
}

// ============================================================
// AST Nodes
// ============================================================

type Node interface {
	TokenLiteral() string
	String() string
}

type NumberLiteral struct{ Value float64 }

func (n *NumberLiteral) TokenLiteral() string { return fmt.Sprintf("%g", n.Value) }
func (n *NumberLiteral) String() string       { return fmt.Sprintf("%g", n.Value) }

type StringLiteral struct{ Value string }

func (s *StringLiteral) TokenLiteral() string { return s.Value }
func (s *StringLiteral) String() string       { return fmt.Sprintf("%q", s.Value) }

type BoolLiteral struct{ Value bool }

func (b *BoolLiteral) TokenLiteral() string { return fmt.Sprintf("%v", b.Value) }
func (b *BoolLiteral) String() string       { return fmt.Sprintf("%v", b.Value) }

type Identifier struct{ Name string }

func (i *Identifier) TokenLiteral() string { return i.Name }
func (i *Identifier) String() string       { return i.Name }

type PrefixExpr struct {
	Op    string
	Right Node
}

func (p *PrefixExpr) TokenLiteral() string { return p.Op }
func (p *PrefixExpr) String() string       { return fmt.Sprintf("(%s%s)", p.Op, p.Right.String()) }

type InfixExpr struct {
	Left  Node
	Op    string
	Right Node
}

func (i *InfixExpr) TokenLiteral() string { return i.Op }
func (i *InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Op, i.Right.String())
}

type LetStatement struct {
	Name  string
	Value Node
}

func (l *LetStatement) TokenLiteral() string { return "let" }
func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", l.Name, l.Value.String())
}

type ReturnStatement struct{ Value Node }

func (r *ReturnStatement) TokenLiteral() string { return "return" }
func (r *ReturnStatement) String() string       { return fmt.Sprintf("return %s", r.Value.String()) }

type BlockStatement struct{ Stmts []Node }

func (b *BlockStatement) TokenLiteral() string { return "{" }
func (b *BlockStatement) String() string {
	parts := make([]string, len(b.Stmts))
	for i, s := range b.Stmts {
		parts[i] = s.String()
	}
	return fmt.Sprintf("{ %s }", strings.Join(parts, "; "))
}

type Program struct{ Statements []Node }

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}
func (p *Program) String() string {
	parts := make([]string, len(p.Statements))
	for i, s := range p.Statements {
		parts[i] = s.String()
	}
	return strings.Join(parts, "\n")
}

// ============================================================
// Parser — literals, identifiers, let/return statements
// ============================================================

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens}
}

func (p *Parser) cur() Token {
	if p.pos < len(p.tokens) {
		return p.tokens[p.pos]
	}
	return Token{Type: EOF}
}

func (p *Parser) advance() Token {
	t := p.cur()
	p.pos++
	return t
}

func (p *Parser) expect(tt TokenType) bool {
	if p.cur().Type == tt {
		p.advance()
		return true
	}
	return false
}

func (p *Parser) ParseProgram() *Program {
	prog := &Program{}
	for p.cur().Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			prog.Statements = append(prog.Statements, stmt)
		}
		for p.cur().Type == SEMICOLON {
			p.advance()
		}
	}
	return prog
}

func (p *Parser) parseStatement() Node {
	switch p.cur().Type {
	case LET:
		return p.parseLetStatement()
	case RETURN:
		return p.parseReturnStatement()
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() Node {
	p.advance() // consume 'let'
	if p.cur().Type != IDENT {
		return nil
	}
	name := p.advance().Literal
	if !p.expect(EQ) {
		return nil
	}
	value := p.parseExpression()
	p.expect(SEMICOLON)
	return &LetStatement{Name: name, Value: value}
}

func (p *Parser) parseReturnStatement() Node {
	p.advance() // consume 'return'
	value := p.parseExpression()
	p.expect(SEMICOLON)
	return &ReturnStatement{Value: value}
}

func (p *Parser) parseExpressionStatement() Node {
	expr := p.parseExpression()
	p.expect(SEMICOLON)
	return expr
}

func (p *Parser) parseExpression() Node {
	return p.parsePrimary()
}

func (p *Parser) parsePrimary() Node {
	switch p.cur().Type {
	case NUMBER:
		t := p.advance()
		v, _ := strconv.ParseFloat(t.Literal, 64)
		return &NumberLiteral{Value: v}
	case STRING:
		t := p.advance()
		return &StringLiteral{Value: t.Literal}
	case TRUE:
		p.advance()
		return &BoolLiteral{Value: true}
	case FALSE:
		p.advance()
		return &BoolLiteral{Value: false}
	case IDENT:
		t := p.advance()
		return &Identifier{Name: t.Literal}
	}
	p.advance() // skip unknown token
	return nil
}

func main() {
	src := `let x = 42`
	tokens := Tokenize(src)
	parser := NewParser(tokens)
	prog := parser.ParseProgram()
	fmt.Println(prog.String())
}
