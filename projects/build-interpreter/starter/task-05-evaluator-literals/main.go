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
				pos++
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

type FnLiteral struct {
	Params []string
	Body   *BlockStatement
}

func (f *FnLiteral) TokenLiteral() string { return "fn" }
func (f *FnLiteral) String() string {
	return fmt.Sprintf("fn(%s) %s", strings.Join(f.Params, ", "), f.Body.String())
}

type CallExpr struct {
	Fn   Node
	Args []Node
}

func (c *CallExpr) TokenLiteral() string { return "(" }
func (c *CallExpr) String() string {
	args := make([]string, len(c.Args))
	for i, a := range c.Args {
		args[i] = a.String()
	}
	return fmt.Sprintf("%s(%s)", c.Fn.String(), strings.Join(args, ", "))
}

type IfExpr struct {
	Condition Node
	Then      *BlockStatement
	Else      *BlockStatement
}

func (ie *IfExpr) TokenLiteral() string { return "if" }
func (ie *IfExpr) String() string {
	s := fmt.Sprintf("if (%s) %s", ie.Condition.String(), ie.Then.String())
	if ie.Else != nil {
		s += " else " + ie.Else.String()
	}
	return s
}

// ============================================================
// Parser
// ============================================================

const (
	PREC_LOWEST      = 1
	PREC_EQUALS      = 2
	PREC_LESSGREATER = 3
	PREC_SUM         = 4
	PREC_PRODUCT     = 5
	PREC_PREFIX      = 6
	PREC_CALL        = 7
)

var precedences = map[TokenType]int{
	EQEQ:   PREC_EQUALS,
	BANGEQ: PREC_EQUALS,
	LT:     PREC_LESSGREATER,
	GT:     PREC_LESSGREATER,
	LTEQ:   PREC_LESSGREATER,
	GTEQ:   PREC_LESSGREATER,
	PLUS:   PREC_SUM,
	MINUS:  PREC_SUM,
	STAR:   PREC_PRODUCT,
	SLASH:  PREC_PRODUCT,
	LPAREN: PREC_CALL,
}

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser { return &Parser{tokens: tokens} }

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

func (p *Parser) curPrec() int {
	if prec, ok := precedences[p.cur().Type]; ok {
		return prec
	}
	return PREC_LOWEST
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
	p.advance()
	if p.cur().Type != IDENT {
		return nil
	}
	name := p.advance().Literal
	if !p.expect(EQ) {
		return nil
	}
	value := p.parseExpression(PREC_LOWEST)
	p.expect(SEMICOLON)
	return &LetStatement{Name: name, Value: value}
}

func (p *Parser) parseReturnStatement() Node {
	p.advance()
	value := p.parseExpression(PREC_LOWEST)
	p.expect(SEMICOLON)
	return &ReturnStatement{Value: value}
}

func (p *Parser) parseExpressionStatement() Node {
	expr := p.parseExpression(PREC_LOWEST)
	p.expect(SEMICOLON)
	return expr
}

func (p *Parser) parseExpression(prec int) Node {
	left := p.parsePrefix()
	if left == nil {
		return nil
	}
	for p.curPrec() > prec {
		if p.cur().Type == LPAREN {
			left = p.parseCallExpr(left)
		} else {
			left = p.parseInfixExpr(left)
		}
	}
	return left
}

func (p *Parser) parsePrefix() Node {
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
	case MINUS:
		p.advance()
		return &PrefixExpr{Op: "-", Right: p.parseExpression(PREC_PREFIX)}
	case BANG:
		p.advance()
		return &PrefixExpr{Op: "!", Right: p.parseExpression(PREC_PREFIX)}
	case LPAREN:
		p.advance()
		expr := p.parseExpression(PREC_LOWEST)
		p.expect(RPAREN)
		return expr
	case FN:
		return p.parseFnLiteral()
	case IF:
		return p.parseIfExpr()
	}
	p.advance()
	return nil
}

func (p *Parser) parseInfixExpr(left Node) Node {
	op := p.advance().Type
	prec := precedences[op]
	right := p.parseExpression(prec)
	return &InfixExpr{Left: left, Op: string(op), Right: right}
}

func (p *Parser) parseCallExpr(fn Node) Node {
	p.advance() // consume '('
	var args []Node
	if p.cur().Type != RPAREN {
		args = append(args, p.parseExpression(PREC_LOWEST))
		for p.cur().Type == COMMA {
			p.advance()
			args = append(args, p.parseExpression(PREC_LOWEST))
		}
	}
	p.expect(RPAREN)
	return &CallExpr{Fn: fn, Args: args}
}

func (p *Parser) parseFnLiteral() Node {
	p.advance()
	if !p.expect(LPAREN) {
		return nil
	}
	var params []string
	if p.cur().Type != RPAREN {
		if p.cur().Type == IDENT {
			params = append(params, p.advance().Literal)
		}
		for p.cur().Type == COMMA {
			p.advance()
			if p.cur().Type == IDENT {
				params = append(params, p.advance().Literal)
			}
		}
	}
	p.expect(RPAREN)
	body := p.parseBlockStatement()
	return &FnLiteral{Params: params, Body: body}
}

func (p *Parser) parseIfExpr() Node {
	p.advance()
	p.expect(LPAREN)
	cond := p.parseExpression(PREC_LOWEST)
	p.expect(RPAREN)
	then := p.parseBlockStatement()
	var els *BlockStatement
	if p.cur().Type == ELSE {
		p.advance()
		els = p.parseBlockStatement()
	}
	return &IfExpr{Condition: cond, Then: then, Else: els}
}

func (p *Parser) parseBlockStatement() *BlockStatement {
	p.expect(LBRACE)
	block := &BlockStatement{}
	for p.cur().Type != RBRACE && p.cur().Type != EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Stmts = append(block.Stmts, stmt)
		}
		for p.cur().Type == SEMICOLON {
			p.advance()
		}
	}
	p.expect(RBRACE)
	return block
}

// ============================================================
// Object system
// ============================================================

type Object interface {
	Type() string
	Inspect() string
}

type NumberObj struct{ Value float64 }

func (n *NumberObj) Type() string    { return "NUMBER" }
func (n *NumberObj) Inspect() string { return fmt.Sprintf("%g", n.Value) }

type StringObj struct{ Value string }

func (s *StringObj) Type() string    { return "STRING" }
func (s *StringObj) Inspect() string { return s.Value }

type BoolObj struct{ Value bool }

func (b *BoolObj) Type() string    { return "BOOL" }
func (b *BoolObj) Inspect() string { return fmt.Sprintf("%v", b.Value) }

type NullObj struct{}

func (n *NullObj) Type() string    { return "NULL" }
func (n *NullObj) Inspect() string { return "null" }

// Singletons for common values.
var (
	TRUE_OBJ  = &BoolObj{Value: true}
	FALSE_OBJ = &BoolObj{Value: false}
	NULL_OBJ  = &NullObj{}
)

// ============================================================
// Environment
// ============================================================

type Environment struct {
	store map[string]Object
	outer *Environment
}

func NewEnvironment() *Environment {
	return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.store[name]
	if !ok && e.outer != nil {
		val, ok = e.outer.Get(name)
	}
	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// ============================================================
// Evaluator — literals only
// ============================================================

func Eval(node Node, env *Environment) Object {
	switch n := node.(type) {
	case *NumberLiteral:
		return &NumberObj{Value: n.Value}
	case *StringLiteral:
		return &StringObj{Value: n.Value}
	case *BoolLiteral:
		if n.Value {
			return TRUE_OBJ
		}
		return FALSE_OBJ
	}
	return NULL_OBJ
}

func main() {
	env := NewEnvironment()
	tokens := Tokenize("42")
	prog := NewParser(tokens).ParseProgram()
	if len(prog.Statements) > 0 {
		fmt.Println(Eval(prog.Statements[0], env).Inspect())
	}
}
