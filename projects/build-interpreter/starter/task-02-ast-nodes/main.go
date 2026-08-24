package main

import (
	"fmt"
	"strings"
)

// Node is the base interface for all AST nodes.
type Node interface {
	TokenLiteral() string
	String() string
}

// NumberLiteral represents a numeric value like 42 or 3.14.
type NumberLiteral struct {
	Value float64
}

func (n *NumberLiteral) TokenLiteral() string { return fmt.Sprintf("%g", n.Value) }
func (n *NumberLiteral) String() string       { return fmt.Sprintf("%g", n.Value) }

// StringLiteral represents a quoted string value.
type StringLiteral struct {
	Value string
}

func (s *StringLiteral) TokenLiteral() string { return s.Value }
func (s *StringLiteral) String() string       { return fmt.Sprintf("%q", s.Value) }

// BoolLiteral represents true or false.
type BoolLiteral struct {
	Value bool
}

func (b *BoolLiteral) TokenLiteral() string { return fmt.Sprintf("%v", b.Value) }
func (b *BoolLiteral) String() string       { return fmt.Sprintf("%v", b.Value) }

// Identifier represents a variable name.
type Identifier struct {
	Name string
}

func (i *Identifier) TokenLiteral() string { return i.Name }
func (i *Identifier) String() string       { return i.Name }

// PrefixExpr represents a prefix operator applied to an expression (-x, !x).
type PrefixExpr struct {
	Op    string
	Right Node
}

func (p *PrefixExpr) TokenLiteral() string { return p.Op }
func (p *PrefixExpr) String() string {
	return fmt.Sprintf("(%s%s)", p.Op, p.Right.String())
}

// InfixExpr represents a binary operator expression (left op right).
type InfixExpr struct {
	Left  Node
	Op    string
	Right Node
}

func (i *InfixExpr) TokenLiteral() string { return i.Op }
func (i *InfixExpr) String() string {
	return fmt.Sprintf("(%s %s %s)", i.Left.String(), i.Op, i.Right.String())
}

// LetStatement represents "let name = value".
type LetStatement struct {
	Name  string
	Value Node
}

func (l *LetStatement) TokenLiteral() string { return "let" }
func (l *LetStatement) String() string {
	return fmt.Sprintf("let %s = %s", l.Name, l.Value.String())
}

// ReturnStatement represents "return value".
type ReturnStatement struct {
	Value Node
}

func (r *ReturnStatement) TokenLiteral() string { return "return" }
func (r *ReturnStatement) String() string {
	return fmt.Sprintf("return %s", r.Value.String())
}

// BlockStatement represents a sequence of statements enclosed in braces.
type BlockStatement struct {
	Stmts []Node
}

func (b *BlockStatement) TokenLiteral() string { return "{" }
func (b *BlockStatement) String() string {
	parts := make([]string, len(b.Stmts))
	for i, s := range b.Stmts {
		parts[i] = s.String()
	}
	return fmt.Sprintf("{ %s }", strings.Join(parts, "; "))
}

// Program is the root AST node containing all top-level statements.
type Program struct {
	Statements []Node
}

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

func main() {
	prog := &Program{
		Statements: []Node{
			&LetStatement{
				Name:  "x",
				Value: &InfixExpr{Left: &NumberLiteral{Value: 1}, Op: "+", Right: &NumberLiteral{Value: 2}},
			},
		},
	}
	fmt.Println(prog.String())
}
