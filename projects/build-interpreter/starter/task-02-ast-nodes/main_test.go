package main

import "testing"

func TestNumberLiteral_String(t *testing.T) {
	n := &NumberLiteral{Value: 3.14}
	if n.String() != "3.14" {
		t.Errorf("expected 3.14, got %q", n.String())
	}
	n2 := &NumberLiteral{Value: 42}
	if n2.String() != "42" {
		t.Errorf("expected 42, got %q", n2.String())
	}
}

func TestStringLiteral_String(t *testing.T) {
	s := &StringLiteral{Value: "hello"}
	if s.String() != `"hello"` {
		t.Errorf("expected %q, got %q", `"hello"`, s.String())
	}
}

func TestBoolLiteral_String(t *testing.T) {
	b := &BoolLiteral{Value: true}
	if b.String() != "true" {
		t.Errorf("expected true, got %q", b.String())
	}
	b2 := &BoolLiteral{Value: false}
	if b2.String() != "false" {
		t.Errorf("expected false, got %q", b2.String())
	}
}

func TestIdentifier_String(t *testing.T) {
	id := &Identifier{Name: "foo"}
	if id.String() != "foo" {
		t.Errorf("expected foo, got %q", id.String())
	}
}

func TestPrefixExpr_String(t *testing.T) {
	e := &PrefixExpr{Op: "-", Right: &NumberLiteral{Value: 5}}
	if e.String() != "(-5)" {
		t.Errorf("expected (-5), got %q", e.String())
	}
	e2 := &PrefixExpr{Op: "!", Right: &BoolLiteral{Value: true}}
	if e2.String() != "(!true)" {
		t.Errorf("expected (!true), got %q", e2.String())
	}
}

func TestInfixExpr_String(t *testing.T) {
	e := &InfixExpr{
		Left:  &Identifier{Name: "a"},
		Op:    "+",
		Right: &Identifier{Name: "b"},
	}
	if e.String() != "(a + b)" {
		t.Errorf("expected (a + b), got %q", e.String())
	}
}

func TestLetStatement_String(t *testing.T) {
	ls := &LetStatement{Name: "x", Value: &NumberLiteral{Value: 42}}
	if ls.String() != "let x = 42" {
		t.Errorf("expected 'let x = 42', got %q", ls.String())
	}
}

func TestReturnStatement_String(t *testing.T) {
	rs := &ReturnStatement{Value: &Identifier{Name: "x"}}
	if rs.String() != "return x" {
		t.Errorf("expected 'return x', got %q", rs.String())
	}
}

func TestBlockStatement_String(t *testing.T) {
	bs := &BlockStatement{
		Stmts: []Node{
			&Identifier{Name: "a"},
			&Identifier{Name: "b"},
		},
	}
	s := bs.String()
	if s == "" {
		t.Error("BlockStatement.String() should not be empty")
	}
	// Should contain both identifiers
	for _, name := range []string{"a", "b"} {
		found := false
		for i := 0; i < len(s)-len(name)+1; i++ {
			if s[i:i+len(name)] == name {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("BlockStatement.String() should contain %q, got %q", name, s)
		}
	}
}

func TestProgram_String(t *testing.T) {
	p := &Program{
		Statements: []Node{
			&LetStatement{Name: "x", Value: &NumberLiteral{Value: 1}},
			&ReturnStatement{Value: &Identifier{Name: "x"}},
		},
	}
	s := p.String()
	if s == "" {
		t.Error("Program.String() should not be empty")
	}
}

func TestNode_Interface(t *testing.T) {
	// All node types must satisfy the Node interface
	var _ Node = &NumberLiteral{}
	var _ Node = &StringLiteral{}
	var _ Node = &BoolLiteral{}
	var _ Node = &Identifier{}
	var _ Node = &PrefixExpr{Op: "-", Right: &NumberLiteral{}}
	var _ Node = &InfixExpr{Left: &NumberLiteral{}, Op: "+", Right: &NumberLiteral{}}
	var _ Node = &LetStatement{Value: &NumberLiteral{}}
	var _ Node = &ReturnStatement{Value: &NumberLiteral{}}
	var _ Node = &BlockStatement{}
	var _ Node = &Program{}
}
