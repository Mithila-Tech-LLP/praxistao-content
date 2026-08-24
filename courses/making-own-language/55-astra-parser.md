# Chapter 55: Building the Astra Parser — Turning Tokens Into Structure

> "Parsing is the art of finding meaning in a stream of symbols. It is the bridge between what was written and what was meant." — Unknown

---

## Overview

The lexer transformed raw characters into tokens. Now the parser takes those tokens and builds an **Abstract Syntax Tree (AST)** — a hierarchical data structure that represents the grammatical structure of the program.

Consider this Astra code:

```astra
let result = 3 + 4 * 2
```

The lexer produces: `LET IDENT ASSIGN INT PLUS INT STAR INT EOF`

The parser's job is to understand that `4 * 2` must be computed before `3 + ...` because multiplication has higher precedence. It produces:

```
LetDeclaration
  name: "result"
  value:
    BinaryExpr(+)
      left:  IntLiteral(3)
      right: BinaryExpr(*)
                left:  IntLiteral(4)
                right: IntLiteral(2)
```

This structure is correct regardless of how it looks on paper. The tree shape encodes the precedence — multiplication's subtree is nested deeper than addition's.

In this chapter we build the complete Astra parser using **Pratt parsing** (also known as Top-Down Operator Precedence parsing). Pratt parsing is elegant for expressions with operator precedence. We combine it with a straightforward recursive descent approach for statements and declarations.

---

## Table of Contents

1. Parser Architecture
2. Helper Methods
3. Parsing Declarations
4. Parsing Statements
5. Pratt Parsing for Expressions
6. Error Recovery
7. Complete Implementation
8. Testing the Parser

---

## 1. Parser Architecture

```
ASCII Diagram: Parser Input/Output

Token Stream (from lexer):
  [FN] [IDENT "add"] [LPAREN] [IDENT "a"] [COLON] [IDENT "int"] ...

                        ↓ Parser ↓

Abstract Syntax Tree:
  FnDeclaration
  ├── name: "add"
  ├── params: [Param{name:"a", type:IntType}]
  ├── returnType: IntType
  └── body: Block
        └── ReturnStatement
              └── BinaryExpr(+)
                    ├── Identifier("a")
                    └── Identifier("b")
```

The parser struct holds:
- The complete token slice (produced by the lexer)
- A `current` index pointing to the next token to consume
- A slice of parse errors (we collect errors and continue)

```
ASCII Diagram: Parser's "sliding window" over token stream

Tokens: [FN] [IDENT "add"] [LPAREN] [IDENT "a"] [COLON] [IDENT "int"] [RPAREN] [EOF]
              ↑
              current (pointing to IDENT "add")

peek()    → tokens[current]   → IDENT "add"
peekNext() → tokens[current+1] → LPAREN
advance() → returns tokens[current], increments current
```

---

## 2. Complete Implementation

### `ast/ast.go` (Minimal version needed by parser)

We need the AST node types before building the parser. Here is a minimal but complete version:

```go
// ast/ast.go
// AST node definitions — the data model for parsed Astra programs.
// (Full version in Chapter 56; this is the interface-level view)

package ast

import "github.com/astra-lang/astrac/lexer"

// Node is the base interface for all AST nodes.
type Node interface {
	nodeType() string
	Pos() lexer.Position
}

// Declaration nodes appear at the top level of a program.
type Declaration interface {
	Node
	declNode()
}

// Statement nodes appear inside function bodies.
type Statement interface {
	Node
	stmtNode()
}

// Expression nodes appear inside statements and produce values.
type Expression interface {
	Node
	exprNode()
}

// TypeExpr represents a type annotation in the source.
type TypeExpr interface {
	Node
	typeNode()
}
```

### `parser/parser.go`

```go
// parser/parser.go
// Recursive descent + Pratt parser for the Astra language.

package parser

import (
	"fmt"
	"strconv"

	"github.com/astra-lang/astrac/ast"
	"github.com/astra-lang/astrac/lexer"
)

// ParseError represents a parsing error.
type ParseError struct {
	Pos     lexer.Position
	Message string
}

func (e *ParseError) Error() string {
	return fmt.Sprintf("parse error at %s: %s", e.Pos, e.Message)
}

// Parser converts a token stream into an AST.
type Parser struct {
	tokens  []lexer.Token
	current int
	errors  []*ParseError
}

// New creates a new Parser for the given token slice.
func New(tokens []lexer.Token) *Parser {
	return &Parser{
		tokens:  tokens,
		current: 0,
		errors:  nil,
	}
}

// Errors returns all parse errors collected during parsing.
func (p *Parser) Errors() []*ParseError {
	return p.errors
}

// ─── Navigation helpers ────────────────────────────────────────────────────────

// peek returns the current token (the next one to be consumed).
func (p *Parser) peek() lexer.Token {
	return p.tokens[p.current]
}

// peekType returns the type of the current token.
func (p *Parser) peekType() lexer.TokenType {
	return p.peek().Type
}

// previous returns the most recently consumed token.
func (p *Parser) previous() lexer.Token {
	return p.tokens[p.current-1]
}

// isAtEnd returns true when the current token is EOF.
func (p *Parser) isAtEnd() bool {
	return p.peekType() == lexer.EOF
}

// check returns true if the current token has the given type (without consuming).
func (p *Parser) check(t lexer.TokenType) bool {
	if p.isAtEnd() {
		return t == lexer.EOF
	}
	return p.peekType() == t
}

// advance consumes and returns the current token.
func (p *Parser) advance() lexer.Token {
	if !p.isAtEnd() {
		p.current++
	}
	return p.previous()
}

// match consumes the current token if it has any of the given types.
// Returns true if it matched, false otherwise.
func (p *Parser) match(types ...lexer.TokenType) bool {
	for _, t := range types {
		if p.check(t) {
			p.advance()
			return true
		}
	}
	return false
}

// expect consumes the token if it matches the expected type.
// If it does not match, records an error and returns the current token without consuming.
func (p *Parser) expect(t lexer.TokenType, msg string) lexer.Token {
	if p.check(t) {
		return p.advance()
	}
	p.addError(p.peek().Pos, msg)
	return p.peek() // don't consume — try to recover
}

// addError records a parse error.
func (p *Parser) addError(pos lexer.Position, msg string) {
	p.errors = append(p.errors, &ParseError{Pos: pos, Message: msg})
}

// synchronize discards tokens until we find a statement boundary.
// Used for error recovery so one error doesn't cascade into dozens.
func (p *Parser) synchronize() {
	p.advance()
	for !p.isAtEnd() {
		// After a closing brace, we're at a clean boundary
		if p.previous().Type == lexer.RBRACE {
			return
		}
		// These tokens start new statements or declarations
		switch p.peekType() {
		case lexer.FN, lexer.LET, lexer.CONST, lexer.STRUCT,
			lexer.IMPL, lexer.IMPORT, lexer.RETURN,
			lexer.IF, lexer.FOR, lexer.WHILE:
			return
		}
		p.advance()
	}
}

// ─── Top-Level Parsing ─────────────────────────────────────────────────────────

// ParseProgram parses the entire program, returning all top-level declarations.
func (p *Parser) ParseProgram() *ast.Program {
	prog := &ast.Program{}

	// Collect imports first
	for p.check(lexer.IMPORT) {
		decl := p.parseImportDeclaration()
		if decl != nil {
			prog.Declarations = append(prog.Declarations, decl)
		}
	}

	// Then all other declarations
	for !p.isAtEnd() {
		decl := p.parseDeclaration()
		if decl != nil {
			prog.Declarations = append(prog.Declarations, decl)
		}
	}

	return prog
}

// parseDeclaration dispatches to the appropriate top-level parser.
func (p *Parser) parseDeclaration() ast.Declaration {
	defer func() {
		if r := recover(); r != nil {
			// If a panic occurred (from a fatal parse error), synchronize
			p.synchronize()
		}
	}()

	pub := false
	if p.match(lexer.PUB) {
		pub = true
	}

	switch p.peekType() {
	case lexer.FN:
		return p.parseFnDeclaration(pub)
	case lexer.STRUCT:
		return p.parseStructDeclaration(pub)
	case lexer.IMPL:
		return p.parseImplDeclaration()
	case lexer.CONST:
		return p.parseConstDeclaration(pub)
	default:
		p.addError(p.peek().Pos, fmt.Sprintf(
			"expected declaration (fn/struct/impl/const), got %s", p.peek().Type))
		p.synchronize()
		return nil
	}
}

// ─── Import ────────────────────────────────────────────────────────────────────

func (p *Parser) parseImportDeclaration() *ast.ImportDeclaration {
	pos := p.peek().Pos
	p.expect(lexer.IMPORT, "expected 'import'")

	var path string
	if p.check(lexer.STRING_LIT) {
		tok := p.advance()
		path = tok.Literal.(string)
	} else if p.check(lexer.IDENT) {
		// identifier optionally followed by :: identifier chains
		path = p.advance().Lexeme
		for p.match(lexer.COLON_COLON) {
			name := p.expect(lexer.IDENT, "expected identifier after '::'")
			path += "::" + name.Lexeme
		}
	} else {
		p.addError(p.peek().Pos, "expected module path after 'import'")
		return nil
	}

	return &ast.ImportDeclaration{Path: path, StartPos: pos}
}

// ─── Function Declaration ──────────────────────────────────────────────────────

func (p *Parser) parseFnDeclaration(pub bool) *ast.FnDeclaration {
	pos := p.peek().Pos
	p.expect(lexer.FN, "expected 'fn'")
	name := p.expect(lexer.IDENT, "expected function name")

	p.expect(lexer.LPAREN, "expected '(' after function name")
	params := p.parseParamList()
	p.expect(lexer.RPAREN, "expected ')' after parameters")

	var returnType ast.TypeExpr
	if p.match(lexer.ARROW) {
		returnType = p.parseTypeExpr()
	}

	body := p.parseBlock()

	return &ast.FnDeclaration{
		Pub:        pub,
		Name:       name.Lexeme,
		Params:     params,
		ReturnType: returnType,
		Body:       body,
		StartPos:   pos,
	}
}

func (p *Parser) parseParamList() []ast.Param {
	var params []ast.Param
	if p.check(lexer.RPAREN) {
		return params // empty parameter list
	}

	params = append(params, p.parseParam())
	for p.match(lexer.COMMA) {
		if p.check(lexer.RPAREN) {
			break // trailing comma allowed
		}
		params = append(params, p.parseParam())
	}
	return params
}

func (p *Parser) parseParam() ast.Param {
	pos := p.peek().Pos

	// Handle "self" parameter
	if p.match(lexer.SELF) {
		return ast.Param{Name: "self", Type: &ast.NamedType{Name: "Self"}, StartPos: pos}
	}

	name := p.expect(lexer.IDENT, "expected parameter name")
	p.expect(lexer.COLON, fmt.Sprintf("expected ':' after parameter name %q", name.Lexeme))
	typ := p.parseTypeExpr()
	return ast.Param{Name: name.Lexeme, Type: typ, StartPos: pos}
}

// ─── Struct Declaration ────────────────────────────────────────────────────────

func (p *Parser) parseStructDeclaration(pub bool) *ast.StructDeclaration {
	pos := p.peek().Pos
	p.expect(lexer.STRUCT, "expected 'struct'")
	name := p.expect(lexer.IDENT, "expected struct name")
	p.expect(lexer.LBRACE, "expected '{' after struct name")

	var fields []ast.FieldDecl
	for !p.check(lexer.RBRACE) && !p.isAtEnd() {
		fieldPos := p.peek().Pos
		fieldName := p.expect(lexer.IDENT, "expected field name")
		p.expect(lexer.COLON, "expected ':' after field name")
		fieldType := p.parseTypeExpr()
		fields = append(fields, ast.FieldDecl{
			Name:     fieldName.Lexeme,
			Type:     fieldType,
			StartPos: fieldPos,
		})
	}

	p.expect(lexer.RBRACE, "expected '}' after struct fields")
	return &ast.StructDeclaration{Pub: pub, Name: name.Lexeme, Fields: fields, StartPos: pos}
}

// ─── Impl Block ────────────────────────────────────────────────────────────────

func (p *Parser) parseImplDeclaration() *ast.ImplDeclaration {
	pos := p.peek().Pos
	p.expect(lexer.IMPL, "expected 'impl'")
	typeName := p.expect(lexer.IDENT, "expected type name after 'impl'")
	p.expect(lexer.LBRACE, "expected '{' after type name")

	var methods []*ast.FnDeclaration
	for !p.check(lexer.RBRACE) && !p.isAtEnd() {
		method := p.parseFnDeclaration(false)
		if method != nil {
			methods = append(methods, method)
		}
	}

	p.expect(lexer.RBRACE, "expected '}' after impl block")
	return &ast.ImplDeclaration{TypeName: typeName.Lexeme, Methods: methods, StartPos: pos}
}

// ─── Const Declaration ─────────────────────────────────────────────────────────

func (p *Parser) parseConstDeclaration(pub bool) *ast.ConstDeclaration {
	pos := p.peek().Pos
	p.expect(lexer.CONST, "expected 'const'")
	name := p.expect(lexer.IDENT, "expected constant name")
	p.expect(lexer.COLON, "expected ':' after constant name")
	typ := p.parseTypeExpr()
	p.expect(lexer.ASSIGN, "expected '=' after constant type")
	value := p.parseExpression()
	return &ast.ConstDeclaration{Pub: pub, Name: name.Lexeme, Type: typ, Value: value, StartPos: pos}
}

// ─── Block and Statements ──────────────────────────────────────────────────────

func (p *Parser) parseBlock() *ast.Block {
	pos := p.peek().Pos
	p.expect(lexer.LBRACE, "expected '{'")
	var stmts []ast.Statement
	for !p.check(lexer.RBRACE) && !p.isAtEnd() {
		stmt := p.parseStatement()
		if stmt != nil {
			stmts = append(stmts, stmt)
		}
	}
	p.expect(lexer.RBRACE, "expected '}'")
	return &ast.Block{Stmts: stmts, StartPos: pos}
}

func (p *Parser) parseStatement() ast.Statement {
	switch p.peekType() {
	case lexer.LET:
		return p.parseLetStatement()
	case lexer.CONST:
		return p.parseConstStatement()
	case lexer.RETURN:
		return p.parseReturnStatement()
	case lexer.IF:
		return p.parseIfStatement()
	case lexer.FOR:
		return p.parseForStatement()
	case lexer.WHILE:
		return p.parseWhileStatement()
	case lexer.BREAK:
		pos := p.advance().Pos
		return &ast.BreakStatement{StartPos: pos}
	case lexer.CONTINUE:
		pos := p.advance().Pos
		return &ast.ContinueStatement{StartPos: pos}
	case lexer.LBRACE:
		block := p.parseBlock()
		return &ast.BlockStatement{Block: block, StartPos: block.StartPos}
	default:
		return p.parseExpressionStatement()
	}
}

func (p *Parser) parseLetStatement() *ast.LetStatement {
	pos := p.peek().Pos
	p.expect(lexer.LET, "expected 'let'")
	name := p.expect(lexer.IDENT, "expected variable name")

	var typeAnnotation ast.TypeExpr
	if p.match(lexer.COLON) {
		typeAnnotation = p.parseTypeExpr()
	}

	p.expect(lexer.ASSIGN, "expected '=' after variable name")
	value := p.parseExpression()

	return &ast.LetStatement{
		Name:     name.Lexeme,
		Type:     typeAnnotation,
		Value:    value,
		StartPos: pos,
	}
}

func (p *Parser) parseConstStatement() *ast.LetStatement {
	// Inside a function, const is treated like let for simplicity
	pos := p.peek().Pos
	p.advance() // consume 'const'
	name := p.expect(lexer.IDENT, "expected constant name")
	p.expect(lexer.COLON, "expected ':' after constant name")
	typ := p.parseTypeExpr()
	p.expect(lexer.ASSIGN, "expected '='")
	value := p.parseExpression()
	return &ast.LetStatement{Name: name.Lexeme, Type: typ, Value: value, StartPos: pos}
}

func (p *Parser) parseReturnStatement() *ast.ReturnStatement {
	pos := p.advance().Pos // consume 'return'

	// If the next token starts a statement or is '}', the return has no value
	var value ast.Expression
	if !p.check(lexer.RBRACE) && !p.isAtEnd() {
		value = p.parseExpression()
	}

	return &ast.ReturnStatement{Value: value, StartPos: pos}
}

func (p *Parser) parseIfStatement() *ast.IfStatement {
	pos := p.advance().Pos // consume 'if'
	condition := p.parseExpression()
	thenBranch := p.parseBlock()

	var elseBranch ast.Statement
	if p.match(lexer.ELSE) {
		if p.check(lexer.IF) {
			elseBranch = p.parseIfStatement()
		} else {
			block := p.parseBlock()
			elseBranch = &ast.BlockStatement{Block: block, StartPos: block.StartPos}
		}
	}

	return &ast.IfStatement{
		Condition:  condition,
		ThenBranch: thenBranch,
		ElseBranch: elseBranch,
		StartPos:   pos,
	}
}

func (p *Parser) parseForStatement() *ast.ForStatement {
	pos := p.advance().Pos // consume 'for'
	varName := p.expect(lexer.IDENT, "expected loop variable name")
	p.expect(lexer.IN, "expected 'in' after loop variable")
	rangeExpr := p.parseExpression()
	body := p.parseBlock()

	return &ast.ForStatement{
		VarName:   varName.Lexeme,
		RangeExpr: rangeExpr,
		Body:      body,
		StartPos:  pos,
	}
}

func (p *Parser) parseWhileStatement() *ast.WhileStatement {
	pos := p.advance().Pos // consume 'while'
	condition := p.parseExpression()
	body := p.parseBlock()

	return &ast.WhileStatement{Condition: condition, Body: body, StartPos: pos}
}

func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	pos := p.peek().Pos
	expr := p.parseExpression()
	// optional semicolon
	p.match(lexer.SEMICOLON)
	return &ast.ExpressionStatement{Expr: expr, StartPos: pos}
}

// ─── Type Expressions ──────────────────────────────────────────────────────────

func (p *Parser) parseTypeExpr() ast.TypeExpr {
	pos := p.peek().Pos

	switch p.peekType() {
	case lexer.IDENT:
		name := p.advance().Lexeme
		// Check for generic types: Option<T>, Result<T, E>
		// (simplified: just return named type for now)
		return &ast.NamedType{Name: name, StartPos: pos}

	case lexer.FN:
		return p.parseFnType()

	case lexer.LBRACKET:
		return p.parseArrayType()

	default:
		p.addError(pos, fmt.Sprintf("expected type, got %s", p.peek().Type))
		return &ast.NamedType{Name: "<error>", StartPos: pos}
	}
}

func (p *Parser) parseFnType() *ast.FnType {
	pos := p.advance().Pos // consume 'fn'
	p.expect(lexer.LPAREN, "expected '(' in function type")
	var params []ast.TypeExpr
	if !p.check(lexer.RPAREN) {
		params = append(params, p.parseTypeExpr())
		for p.match(lexer.COMMA) {
			params = append(params, p.parseTypeExpr())
		}
	}
	p.expect(lexer.RPAREN, "expected ')' in function type")
	p.expect(lexer.ARROW, "expected '->' in function type")
	ret := p.parseTypeExpr()
	return &ast.FnType{Params: params, Return: ret, StartPos: pos}
}

func (p *Parser) parseArrayType() *ast.ArrayType {
	pos := p.advance().Pos // consume '['
	elem := p.parseTypeExpr()
	p.expect(lexer.RBRACKET, "expected ']'")
	return &ast.ArrayType{Element: elem, StartPos: pos}
}

// ─── Pratt Parser for Expressions ─────────────────────────────────────────────
//
// Pratt parsing assigns each token type a "binding power" (precedence).
// - nud (null denotation): how a token starts an expression (prefix position)
// - led (left denotation): how a token continues an expression (infix position)
//
// The main loop: parse the left side, then keep absorbing operators as long
// as their binding power is greater than our current minimum.

// bindingPower returns the left binding power for an infix operator.
// Higher number = higher precedence.
func bindingPower(t lexer.TokenType) int {
	switch t {
	case lexer.ASSIGN:
		return 5
	case lexer.OR:
		return 10
	case lexer.AND:
		return 20
	case lexer.EQ, lexer.NEQ:
		return 30
	case lexer.LT, lexer.GT, lexer.LTE, lexer.GTE:
		return 40
	case lexer.PLUS, lexer.MINUS:
		return 50
	case lexer.STAR, lexer.SLASH, lexer.PERCENT:
		return 60
	case lexer.DOTDOT:
		return 70 // range operator
	case lexer.DOT, lexer.LPAREN, lexer.LBRACKET:
		return 80 // field access, call, index
	default:
		return 0
	}
}

// parseExpression is the entry point for parsing an expression.
func (p *Parser) parseExpression() ast.Expression {
	return p.parseExprWithPrecedence(0)
}

// parseExprWithPrecedence implements the core Pratt parsing loop.
func (p *Parser) parseExprWithPrecedence(minBP int) ast.Expression {
	// Parse the "nud" (prefix/atom) part
	left := p.parsePrefix()

	// Keep consuming infix operators while their precedence exceeds minBP
	for !p.isAtEnd() && bindingPower(p.peekType()) > minBP {
		left = p.parseInfix(left)
	}

	return left
}

// parsePrefix handles tokens that appear at the start of an expression.
func (p *Parser) parsePrefix() ast.Expression {
	tok := p.peek()

	switch tok.Type {
	case lexer.INT_LIT:
		p.advance()
		return &ast.IntLiteral{Value: tok.Literal.(int64), StartPos: tok.Pos}

	case lexer.FLOAT_LIT:
		p.advance()
		return &ast.FloatLiteral{Value: tok.Literal.(float64), StartPos: tok.Pos}

	case lexer.STRING_LIT:
		p.advance()
		return &ast.StringLiteral{Value: tok.Literal.(string), StartPos: tok.Pos}

	case lexer.TRUE:
		p.advance()
		return &ast.BoolLiteral{Value: true, StartPos: tok.Pos}

	case lexer.FALSE:
		p.advance()
		return &ast.BoolLiteral{Value: false, StartPos: tok.Pos}

	case lexer.IDENT:
		p.advance()
		// Check for struct literal: Identifier { field: value, ... }
		// We must be careful: this only applies when the identifier is followed
		// by '{' and then a field name and ':'. This disambiguates from blocks.
		if p.check(lexer.LBRACE) && p.isStructLiteralStart() {
			return p.parseStructLiteral(tok)
		}
		// Check for path expression: Module::function
		if p.check(lexer.COLON_COLON) {
			return p.parsePathExpr(tok)
		}
		return &ast.Identifier{Name: tok.Lexeme, StartPos: tok.Pos}

	case lexer.MINUS:
		p.advance()
		operand := p.parseExprWithPrecedence(65) // higher than any binary op
		return &ast.UnaryExpr{Op: "-", Operand: operand, StartPos: tok.Pos}

	case lexer.BANG:
		p.advance()
		operand := p.parseExprWithPrecedence(65)
		return &ast.UnaryExpr{Op: "!", Operand: operand, StartPos: tok.Pos}

	case lexer.LPAREN:
		p.advance() // consume '('
		expr := p.parseExpression()
		p.expect(lexer.RPAREN, "expected ')' after grouped expression")
		return expr

	default:
		p.addError(tok.Pos, fmt.Sprintf("unexpected token in expression: %s %q", tok.Type, tok.Lexeme))
		p.advance() // consume the problematic token
		return &ast.ErrorExpr{StartPos: tok.Pos}
	}
}

// isStructLiteralStart peeks ahead to determine if we're looking at a struct literal.
// A struct literal looks like: Ident { field: expr, ... }
// A block looks like: { stmt; stmt; ... }
func (p *Parser) isStructLiteralStart() bool {
	// Save current position
	savedCurrent := p.current
	// We've already consumed the identifier, peeking at '{'
	// Look ahead: next should be IDENT ':' (field name colon)
	p.advance() // consume '{'
	result := p.check(lexer.IDENT)
	if result {
		p.advance() // consume IDENT
		result = p.check(lexer.COLON)
	}
	// Restore position
	p.current = savedCurrent
	return result
}

func (p *Parser) parseStructLiteral(nameTok lexer.Token) *ast.StructLiteral {
	p.expect(lexer.LBRACE, "expected '{'")
	var fields []ast.FieldInit
	for !p.check(lexer.RBRACE) && !p.isAtEnd() {
		fieldPos := p.peek().Pos
		fieldName := p.expect(lexer.IDENT, "expected field name")
		p.expect(lexer.COLON, "expected ':' after field name")
		fieldVal := p.parseExpression()
		fields = append(fields, ast.FieldInit{
			Name:     fieldName.Lexeme,
			Value:    fieldVal,
			StartPos: fieldPos,
		})
		if !p.check(lexer.RBRACE) {
			p.expect(lexer.COMMA, "expected ',' or '}' after field value")
		}
	}
	p.expect(lexer.RBRACE, "expected '}' after struct fields")
	return &ast.StructLiteral{TypeName: nameTok.Lexeme, Fields: fields, StartPos: nameTok.Pos}
}

func (p *Parser) parsePathExpr(first lexer.Token) *ast.PathExpr {
	p.advance() // consume '::'
	second := p.expect(lexer.IDENT, "expected identifier after '::'")
	return &ast.PathExpr{
		Module: first.Lexeme,
		Name:   second.Lexeme,
		StartPos: first.Pos,
	}
}

// parseInfix handles operators that appear between two expressions.
func (p *Parser) parseInfix(left ast.Expression) ast.Expression {
	tok := p.advance()

	switch tok.Type {
	// ── Binary arithmetic and comparison ────────────────────────────────────
	case lexer.PLUS, lexer.MINUS, lexer.STAR, lexer.SLASH, lexer.PERCENT,
		lexer.EQ, lexer.NEQ, lexer.LT, lexer.GT, lexer.LTE, lexer.GTE,
		lexer.AND, lexer.OR:
		right := p.parseExprWithPrecedence(bindingPower(tok.Type))
		return &ast.BinaryExpr{Op: tok.Lexeme, Left: left, Right: right, StartPos: tok.Pos}

	// ── Assignment ────────────────────────────────────────────────────────
	case lexer.ASSIGN:
		// Assignment is right-associative: parse right side with same BP
		right := p.parseExprWithPrecedence(bindingPower(tok.Type) - 1)
		return &ast.AssignExpr{Target: left, Value: right, StartPos: tok.Pos}

	// ── Range ─────────────────────────────────────────────────────────────
	case lexer.DOTDOT:
		right := p.parseExprWithPrecedence(bindingPower(tok.Type))
		return &ast.RangeExpr{Start: left, End: right, StartPos: tok.Pos}

	// ── Field access and method call ──────────────────────────────────────
	case lexer.DOT:
		fieldName := p.expect(lexer.IDENT, "expected field name after '.'")
		if p.check(lexer.LPAREN) {
			// Method call: expr.method(args)
			p.advance() // consume '('
			args := p.parseArgList()
			p.expect(lexer.RPAREN, "expected ')' after method arguments")
			return &ast.MethodCall{
				Object:   left,
				Method:   fieldName.Lexeme,
				Args:     args,
				StartPos: tok.Pos,
			}
		}
		// Field access: expr.field
		return &ast.FieldAccess{Object: left, Field: fieldName.Lexeme, StartPos: tok.Pos}

	// ── Function call ─────────────────────────────────────────────────────
	case lexer.LPAREN:
		args := p.parseArgList()
		p.expect(lexer.RPAREN, "expected ')' after arguments")
		return &ast.CallExpr{Callee: left, Args: args, StartPos: tok.Pos}

	// ── Index access ──────────────────────────────────────────────────────
	case lexer.LBRACKET:
		index := p.parseExpression()
		p.expect(lexer.RBRACKET, "expected ']' after index")
		return &ast.IndexExpr{Object: left, Index: index, StartPos: tok.Pos}

	default:
		p.addError(tok.Pos, fmt.Sprintf("unexpected infix operator: %s", tok.Type))
		return left
	}
}

func (p *Parser) parseArgList() []ast.Expression {
	var args []ast.Expression
	if p.check(lexer.RPAREN) {
		return args
	}
	args = append(args, p.parseExpression())
	for p.match(lexer.COMMA) {
		if p.check(lexer.RPAREN) {
			break // trailing comma
		}
		args = append(args, p.parseExpression())
	}
	return args
}

// ─── Unused import suppression ─────────────────────────────────────────────────
var _ = strconv.Itoa
```

---

## 3. The AST Node Types (referenced by parser)

```go
// ast/nodes.go — all concrete node types

package ast

import "github.com/astra-lang/astrac/lexer"

// ─── Program ───────────────────────────────────────────────────────────────────

type Program struct {
	Declarations []Declaration
}

// ─── Declarations ──────────────────────────────────────────────────────────────

type ImportDeclaration struct {
	Path     string
	StartPos lexer.Position
}
func (n *ImportDeclaration) nodeType() string     { return "ImportDeclaration" }
func (n *ImportDeclaration) Pos() lexer.Position  { return n.StartPos }
func (n *ImportDeclaration) declNode()            {}

type FnDeclaration struct {
	Pub        bool
	Name       string
	Params     []Param
	ReturnType TypeExpr
	Body       *Block
	StartPos   lexer.Position
}
func (n *FnDeclaration) nodeType() string    { return "FnDeclaration" }
func (n *FnDeclaration) Pos() lexer.Position { return n.StartPos }
func (n *FnDeclaration) declNode()           {}

type StructDeclaration struct {
	Pub      bool
	Name     string
	Fields   []FieldDecl
	StartPos lexer.Position
}
func (n *StructDeclaration) nodeType() string    { return "StructDeclaration" }
func (n *StructDeclaration) Pos() lexer.Position { return n.StartPos }
func (n *StructDeclaration) declNode()           {}

type ImplDeclaration struct {
	TypeName string
	Methods  []*FnDeclaration
	StartPos lexer.Position
}
func (n *ImplDeclaration) nodeType() string    { return "ImplDeclaration" }
func (n *ImplDeclaration) Pos() lexer.Position { return n.StartPos }
func (n *ImplDeclaration) declNode()           {}

type ConstDeclaration struct {
	Pub      bool
	Name     string
	Type     TypeExpr
	Value    Expression
	StartPos lexer.Position
}
func (n *ConstDeclaration) nodeType() string    { return "ConstDeclaration" }
func (n *ConstDeclaration) Pos() lexer.Position { return n.StartPos }
func (n *ConstDeclaration) declNode()           {}

// ─── Statements ────────────────────────────────────────────────────────────────

type Block struct {
	Stmts    []Statement
	StartPos lexer.Position
}

type LetStatement struct {
	Name     string
	Type     TypeExpr
	Value    Expression
	StartPos lexer.Position
}
func (n *LetStatement) nodeType() string    { return "LetStatement" }
func (n *LetStatement) Pos() lexer.Position { return n.StartPos }
func (n *LetStatement) stmtNode()           {}

type ReturnStatement struct {
	Value    Expression
	StartPos lexer.Position
}
func (n *ReturnStatement) nodeType() string    { return "ReturnStatement" }
func (n *ReturnStatement) Pos() lexer.Position { return n.StartPos }
func (n *ReturnStatement) stmtNode()           {}

type IfStatement struct {
	Condition  Expression
	ThenBranch *Block
	ElseBranch Statement
	StartPos   lexer.Position
}
func (n *IfStatement) nodeType() string    { return "IfStatement" }
func (n *IfStatement) Pos() lexer.Position { return n.StartPos }
func (n *IfStatement) stmtNode()           {}

type ForStatement struct {
	VarName   string
	RangeExpr Expression
	Body      *Block
	StartPos  lexer.Position
}
func (n *ForStatement) nodeType() string    { return "ForStatement" }
func (n *ForStatement) Pos() lexer.Position { return n.StartPos }
func (n *ForStatement) stmtNode()           {}

type WhileStatement struct {
	Condition Expression
	Body      *Block
	StartPos  lexer.Position
}
func (n *WhileStatement) nodeType() string    { return "WhileStatement" }
func (n *WhileStatement) Pos() lexer.Position { return n.StartPos }
func (n *WhileStatement) stmtNode()           {}

type BreakStatement    struct { StartPos lexer.Position }
type ContinueStatement struct { StartPos lexer.Position }
func (n *BreakStatement) nodeType() string    { return "BreakStatement" }
func (n *BreakStatement) Pos() lexer.Position { return n.StartPos }
func (n *BreakStatement) stmtNode()           {}
func (n *ContinueStatement) nodeType() string    { return "ContinueStatement" }
func (n *ContinueStatement) Pos() lexer.Position { return n.StartPos }
func (n *ContinueStatement) stmtNode()           {}

type BlockStatement struct {
	Block    *Block
	StartPos lexer.Position
}
func (n *BlockStatement) nodeType() string    { return "BlockStatement" }
func (n *BlockStatement) Pos() lexer.Position { return n.StartPos }
func (n *BlockStatement) stmtNode()           {}

type ExpressionStatement struct {
	Expr     Expression
	StartPos lexer.Position
}
func (n *ExpressionStatement) nodeType() string    { return "ExpressionStatement" }
func (n *ExpressionStatement) Pos() lexer.Position { return n.StartPos }
func (n *ExpressionStatement) stmtNode()           {}

// ─── Expressions ───────────────────────────────────────────────────────────────

type IntLiteral struct {
	Value    int64
	StartPos lexer.Position
}
func (n *IntLiteral) nodeType() string    { return "IntLiteral" }
func (n *IntLiteral) Pos() lexer.Position { return n.StartPos }
func (n *IntLiteral) exprNode()           {}

type FloatLiteral struct {
	Value    float64
	StartPos lexer.Position
}
func (n *FloatLiteral) nodeType() string    { return "FloatLiteral" }
func (n *FloatLiteral) Pos() lexer.Position { return n.StartPos }
func (n *FloatLiteral) exprNode()           {}

type StringLiteral struct {
	Value    string
	StartPos lexer.Position
}
func (n *StringLiteral) nodeType() string    { return "StringLiteral" }
func (n *StringLiteral) Pos() lexer.Position { return n.StartPos }
func (n *StringLiteral) exprNode()           {}

type BoolLiteral struct {
	Value    bool
	StartPos lexer.Position
}
func (n *BoolLiteral) nodeType() string    { return "BoolLiteral" }
func (n *BoolLiteral) Pos() lexer.Position { return n.StartPos }
func (n *BoolLiteral) exprNode()           {}

type Identifier struct {
	Name     string
	StartPos lexer.Position
}
func (n *Identifier) nodeType() string    { return "Identifier" }
func (n *Identifier) Pos() lexer.Position { return n.StartPos }
func (n *Identifier) exprNode()           {}

type BinaryExpr struct {
	Op       string
	Left     Expression
	Right    Expression
	StartPos lexer.Position
}
func (n *BinaryExpr) nodeType() string    { return "BinaryExpr" }
func (n *BinaryExpr) Pos() lexer.Position { return n.StartPos }
func (n *BinaryExpr) exprNode()           {}

type UnaryExpr struct {
	Op       string
	Operand  Expression
	StartPos lexer.Position
}
func (n *UnaryExpr) nodeType() string    { return "UnaryExpr" }
func (n *UnaryExpr) Pos() lexer.Position { return n.StartPos }
func (n *UnaryExpr) exprNode()           {}

type AssignExpr struct {
	Target   Expression
	Value    Expression
	StartPos lexer.Position
}
func (n *AssignExpr) nodeType() string    { return "AssignExpr" }
func (n *AssignExpr) Pos() lexer.Position { return n.StartPos }
func (n *AssignExpr) exprNode()           {}

type CallExpr struct {
	Callee   Expression
	Args     []Expression
	StartPos lexer.Position
}
func (n *CallExpr) nodeType() string    { return "CallExpr" }
func (n *CallExpr) Pos() lexer.Position { return n.StartPos }
func (n *CallExpr) exprNode()           {}

type MethodCall struct {
	Object   Expression
	Method   string
	Args     []Expression
	StartPos lexer.Position
}
func (n *MethodCall) nodeType() string    { return "MethodCall" }
func (n *MethodCall) Pos() lexer.Position { return n.StartPos }
func (n *MethodCall) exprNode()           {}

type FieldAccess struct {
	Object   Expression
	Field    string
	StartPos lexer.Position
}
func (n *FieldAccess) nodeType() string    { return "FieldAccess" }
func (n *FieldAccess) Pos() lexer.Position { return n.StartPos }
func (n *FieldAccess) exprNode()           {}

type IndexExpr struct {
	Object   Expression
	Index    Expression
	StartPos lexer.Position
}
func (n *IndexExpr) nodeType() string    { return "IndexExpr" }
func (n *IndexExpr) Pos() lexer.Position { return n.StartPos }
func (n *IndexExpr) exprNode()           {}

type RangeExpr struct {
	Start    Expression
	End      Expression
	StartPos lexer.Position
}
func (n *RangeExpr) nodeType() string    { return "RangeExpr" }
func (n *RangeExpr) Pos() lexer.Position { return n.StartPos }
func (n *RangeExpr) exprNode()           {}

type StructLiteral struct {
	TypeName string
	Fields   []FieldInit
	StartPos lexer.Position
}
func (n *StructLiteral) nodeType() string    { return "StructLiteral" }
func (n *StructLiteral) Pos() lexer.Position { return n.StartPos }
func (n *StructLiteral) exprNode()           {}

type PathExpr struct {
	Module   string
	Name     string
	StartPos lexer.Position
}
func (n *PathExpr) nodeType() string    { return "PathExpr" }
func (n *PathExpr) Pos() lexer.Position { return n.StartPos }
func (n *PathExpr) exprNode()           {}

type ErrorExpr struct {
	StartPos lexer.Position
}
func (n *ErrorExpr) nodeType() string    { return "ErrorExpr" }
func (n *ErrorExpr) Pos() lexer.Position { return n.StartPos }
func (n *ErrorExpr) exprNode()           {}

// ─── Type Expressions ──────────────────────────────────────────────────────────

type NamedType struct {
	Name     string
	StartPos lexer.Position
}
func (n *NamedType) nodeType() string    { return "NamedType" }
func (n *NamedType) Pos() lexer.Position { return n.StartPos }
func (n *NamedType) typeNode()           {}

type FnType struct {
	Params   []TypeExpr
	Return   TypeExpr
	StartPos lexer.Position
}
func (n *FnType) nodeType() string    { return "FnType" }
func (n *FnType) Pos() lexer.Position { return n.StartPos }
func (n *FnType) typeNode()           {}

type ArrayType struct {
	Element  TypeExpr
	StartPos lexer.Position
}
func (n *ArrayType) nodeType() string    { return "ArrayType" }
func (n *ArrayType) Pos() lexer.Position { return n.StartPos }
func (n *ArrayType) typeNode()           {}

// ─── Auxiliary Types ────────────────────────────────────────────────────────────

type Param struct {
	Name     string
	Type     TypeExpr
	StartPos lexer.Position
}

type FieldDecl struct {
	Name     string
	Type     TypeExpr
	StartPos lexer.Position
}

type FieldInit struct {
	Name     string
	Value    Expression
	StartPos lexer.Position
}
```

---

## 4. Testing the Parser

```go
// parser/parser_test.go
package parser

import (
	"testing"

	"github.com/astra-lang/astrac/ast"
	"github.com/astra-lang/astrac/lexer"
)

func parseSource(t *testing.T, source string) *ast.Program {
	t.Helper()
	l := lexer.New(source)
	tokens := l.ScanAll()
	if errs := l.Errors(); len(errs) > 0 {
		t.Fatalf("lexer errors: %v", errs)
	}
	p := New(tokens)
	prog := p.ParseProgram()
	if errs := p.Errors(); len(errs) > 0 {
		t.Fatalf("parser errors: %v", errs)
	}
	return prog
}

func TestParseFunctionDeclaration(t *testing.T) {
	prog := parseSource(t, `fn add(a: int, b: int) -> int { return a + b }`)
	if len(prog.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(prog.Declarations))
	}
	fn, ok := prog.Declarations[0].(*ast.FnDeclaration)
	if !ok {
		t.Fatalf("expected FnDeclaration, got %T", prog.Declarations[0])
	}
	if fn.Name != "add" {
		t.Errorf("expected name 'add', got %q", fn.Name)
	}
	if len(fn.Params) != 2 {
		t.Errorf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.ReturnType == nil {
		t.Error("expected return type, got nil")
	}
}

func TestParseLetStatement(t *testing.T) {
	prog := parseSource(t, `fn main() { let x: int = 42 }`)
	fn := prog.Declarations[0].(*ast.FnDeclaration)
	letStmt, ok := fn.Body.Stmts[0].(*ast.LetStatement)
	if !ok {
		t.Fatalf("expected LetStatement, got %T", fn.Body.Stmts[0])
	}
	if letStmt.Name != "x" {
		t.Errorf("expected name 'x', got %q", letStmt.Name)
	}
	lit, ok := letStmt.Value.(*ast.IntLiteral)
	if !ok || lit.Value != 42 {
		t.Errorf("expected IntLiteral(42), got %T", letStmt.Value)
	}
}

func TestParseIfStatement(t *testing.T) {
	prog := parseSource(t, `fn main() { if x > 0 { return x } else { return 0 } }`)
	fn := prog.Declarations[0].(*ast.FnDeclaration)
	ifStmt, ok := fn.Body.Stmts[0].(*ast.IfStatement)
	if !ok {
		t.Fatalf("expected IfStatement, got %T", fn.Body.Stmts[0])
	}
	if ifStmt.ElseBranch == nil {
		t.Error("expected else branch, got nil")
	}
}

func TestParseForLoop(t *testing.T) {
	prog := parseSource(t, `fn main() { for i in 0..10 { print(i.to_string()) } }`)
	fn := prog.Declarations[0].(*ast.FnDeclaration)
	forStmt, ok := fn.Body.Stmts[0].(*ast.ForStatement)
	if !ok {
		t.Fatalf("expected ForStatement, got %T", fn.Body.Stmts[0])
	}
	if forStmt.VarName != "i" {
		t.Errorf("expected loop var 'i', got %q", forStmt.VarName)
	}
	rangeExpr, ok := forStmt.RangeExpr.(*ast.RangeExpr)
	if !ok {
		t.Fatalf("expected RangeExpr, got %T", forStmt.RangeExpr)
	}
	start, _ := rangeExpr.Start.(*ast.IntLiteral)
	end, _ := rangeExpr.End.(*ast.IntLiteral)
	if start == nil || start.Value != 0 {
		t.Error("expected range start 0")
	}
	if end == nil || end.Value != 10 {
		t.Error("expected range end 10")
	}
}

func TestParseStructDeclaration(t *testing.T) {
	prog := parseSource(t, `struct Person { name: string age: int }`)
	if len(prog.Declarations) != 1 {
		t.Fatalf("expected 1 declaration, got %d", len(prog.Declarations))
	}
	structDecl, ok := prog.Declarations[0].(*ast.StructDeclaration)
	if !ok {
		t.Fatalf("expected StructDeclaration, got %T", prog.Declarations[0])
	}
	if structDecl.Name != "Person" {
		t.Errorf("expected 'Person', got %q", structDecl.Name)
	}
	if len(structDecl.Fields) != 2 {
		t.Errorf("expected 2 fields, got %d", len(structDecl.Fields))
	}
}

func TestParseErrors(t *testing.T) {
	l := lexer.New(`fn () { }`) // missing function name
	tokens := l.ScanAll()
	p := New(tokens)
	p.ParseProgram()
	if len(p.Errors()) == 0 {
		t.Error("expected parse error for missing function name")
	}
}
```

---

## Summary Table

| Parser Method | What It Parses | Returns |
|---|---|---|
| `ParseProgram()` | Entire source file | `*ast.Program` |
| `parseDeclaration()` | One top-level declaration | `ast.Declaration` |
| `parseFnDeclaration()` | `fn name(params) -> type { }` | `*ast.FnDeclaration` |
| `parseStructDeclaration()` | `struct Name { fields }` | `*ast.StructDeclaration` |
| `parseImplDeclaration()` | `impl Type { methods }` | `*ast.ImplDeclaration` |
| `parseBlock()` | `{ statements }` | `*ast.Block` |
| `parseStatement()` | One statement | `ast.Statement` |
| `parseExpression()` | One expression | `ast.Expression` |
| `parsePrefix()` | Literal, identifier, unary | `ast.Expression` |
| `parseInfix(left)` | Binary op, call, field access | `ast.Expression` |
| `parseTypeExpr()` | A type annotation | `ast.TypeExpr` |

---

## Exercises

1. **Parse Chained Methods**: The expression `"hello".to_upper().trim()` should parse as `MethodCall(MethodCall(StringLiteral, "to_upper"), "trim")`. Trace through the Pratt parser to verify this works with the current implementation.

2. **Trailing Commas**: Astra allows trailing commas in argument lists. Add a test case that verifies `add(1, 2,)` parses correctly (comma after last argument).

3. **Error Recovery Test**: Write a test for the parser's error recovery. Give it a program with two syntax errors in different functions. Verify that (a) both errors are reported, and (b) the parser still produces the correct AST structure where possible.

4. **Precedence Verification**: Write test cases that verify operator precedence. The expression `2 + 3 * 4` should produce `BinaryExpr(+, 2, BinaryExpr(*, 3, 4))`. Write similar tests for `&&` vs `||`, and comparison vs arithmetic.

5. **While Loop Extension**: Astra's grammar allows `break` and `continue` inside loops. Modify `parseBreakStatement` to optionally accept a label (`break 'outer`). Write the new grammar rule and implement it.

6. **Struct Literal Ambiguity**: The parser uses `isStructLiteralStart()` to distinguish struct literals from blocks. This lookahead is O(1) but fragile. Find a case where this heuristic could fail and describe a more robust solution.

7. **Full Program Parse**: Write a parser test that parses the complete Astra sample from Chapter 53 (the one with main, add, Person struct, and impl block). Verify the declaration count and structure.

8. **Pretty Printer Integration**: Connect the parser output to the AST pretty printer from Chapter 56. Write a round-trip test: parse source → pretty print AST → verify the AST structure matches what you expect.
