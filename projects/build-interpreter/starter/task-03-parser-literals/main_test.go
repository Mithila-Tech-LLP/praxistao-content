package main

import "testing"

func parseSource(src string) *Program {
	tokens := Tokenize(src)
	p := NewParser(tokens)
	return p.ParseProgram()
}

func TestParse_LetStatement(t *testing.T) {
	prog := parseSource("let x = 42")
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	ls, ok := prog.Statements[0].(*LetStatement)
	if !ok {
		t.Fatalf("expected *LetStatement, got %T", prog.Statements[0])
	}
	if ls.Name != "x" {
		t.Errorf("expected name=x, got %q", ls.Name)
	}
	num, ok := ls.Value.(*NumberLiteral)
	if !ok {
		t.Fatalf("expected *NumberLiteral for value, got %T", ls.Value)
	}
	if num.Value != 42 {
		t.Errorf("expected value=42, got %v", num.Value)
	}
}

func TestParse_BoolLiteralTrue(t *testing.T) {
	prog := parseSource("true")
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	b, ok := prog.Statements[0].(*BoolLiteral)
	if !ok {
		t.Fatalf("expected *BoolLiteral, got %T", prog.Statements[0])
	}
	if !b.Value {
		t.Error("expected true")
	}
}

func TestParse_BoolLiteralFalse(t *testing.T) {
	prog := parseSource("false")
	b, ok := prog.Statements[0].(*BoolLiteral)
	if !ok {
		t.Fatalf("expected *BoolLiteral, got %T", prog.Statements[0])
	}
	if b.Value {
		t.Error("expected false")
	}
}

func TestParse_NumberLiteral(t *testing.T) {
	prog := parseSource("3.14")
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	n, ok := prog.Statements[0].(*NumberLiteral)
	if !ok {
		t.Fatalf("expected *NumberLiteral, got %T", prog.Statements[0])
	}
	if n.Value != 3.14 {
		t.Errorf("expected 3.14, got %v", n.Value)
	}
}

func TestParse_Identifier(t *testing.T) {
	prog := parseSource("hello")
	if len(prog.Statements) != 1 {
		t.Fatalf("expected 1 statement, got %d", len(prog.Statements))
	}
	id, ok := prog.Statements[0].(*Identifier)
	if !ok {
		t.Fatalf("expected *Identifier, got %T", prog.Statements[0])
	}
	if id.Name != "hello" {
		t.Errorf("expected name=hello, got %q", id.Name)
	}
}

func TestParse_LetWithString(t *testing.T) {
	prog := parseSource(`let msg = "hello"`)
	ls, ok := prog.Statements[0].(*LetStatement)
	if !ok {
		t.Fatalf("expected *LetStatement, got %T", prog.Statements[0])
	}
	sv, ok := ls.Value.(*StringLiteral)
	if !ok {
		t.Fatalf("expected *StringLiteral, got %T", ls.Value)
	}
	if sv.Value != "hello" {
		t.Errorf("expected hello, got %q", sv.Value)
	}
}

func TestParse_MultipleStatements(t *testing.T) {
	prog := parseSource("let a = 1\nlet b = 2")
	if len(prog.Statements) != 2 {
		t.Fatalf("expected 2 statements, got %d", len(prog.Statements))
	}
}
