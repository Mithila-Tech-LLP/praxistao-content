package main

import "testing"

func parseExpr(src string) Node {
	tokens := Tokenize(src)
	p := NewParser(tokens)
	prog := p.ParseProgram()
	if len(prog.Statements) == 0 {
		return nil
	}
	return prog.Statements[0]
}

func TestParse_InfixPrecedence(t *testing.T) {
	// 1+2*3 should parse as (1 + (2 * 3)) due to * binding tighter
	node := parseExpr("1+2*3")
	infix, ok := node.(*InfixExpr)
	if !ok {
		t.Fatalf("expected *InfixExpr at root, got %T", node)
	}
	if infix.Op != "+" {
		t.Errorf("expected root op=+, got %q", infix.Op)
	}
	// Right side should be (2 * 3)
	right, ok := infix.Right.(*InfixExpr)
	if !ok {
		t.Fatalf("expected *InfixExpr on right, got %T", infix.Right)
	}
	if right.Op != "*" {
		t.Errorf("expected right op=*, got %q", right.Op)
	}
}

func TestParse_PrefixMinus(t *testing.T) {
	node := parseExpr("-5")
	pre, ok := node.(*PrefixExpr)
	if !ok {
		t.Fatalf("expected *PrefixExpr, got %T", node)
	}
	if pre.Op != "-" {
		t.Errorf("expected op=-, got %q", pre.Op)
	}
	num, ok := pre.Right.(*NumberLiteral)
	if !ok {
		t.Fatalf("expected *NumberLiteral for right, got %T", pre.Right)
	}
	if num.Value != 5 {
		t.Errorf("expected 5, got %v", num.Value)
	}
}

func TestParse_PrefixBang(t *testing.T) {
	node := parseExpr("!true")
	pre, ok := node.(*PrefixExpr)
	if !ok {
		t.Fatalf("expected *PrefixExpr, got %T", node)
	}
	if pre.Op != "!" {
		t.Errorf("expected op=!, got %q", pre.Op)
	}
}

func TestParse_GroupingOverridesPrecedence(t *testing.T) {
	// (1+2)*3 should parse as ((1+2) * 3)
	node := parseExpr("(1+2)*3")
	infix, ok := node.(*InfixExpr)
	if !ok {
		t.Fatalf("expected *InfixExpr at root, got %T", node)
	}
	if infix.Op != "*" {
		t.Errorf("expected root op=*, got %q (grouping should override precedence)", infix.Op)
	}
}

func TestParse_FnLiteral(t *testing.T) {
	node := parseExpr("fn(a, b) { a + b }")
	fn, ok := node.(*FnLiteral)
	if !ok {
		t.Fatalf("expected *FnLiteral, got %T", node)
	}
	if len(fn.Params) != 2 {
		t.Fatalf("expected 2 params, got %d", len(fn.Params))
	}
	if fn.Params[0] != "a" || fn.Params[1] != "b" {
		t.Errorf("expected params [a b], got %v", fn.Params)
	}
	if fn.Body == nil || len(fn.Body.Stmts) == 0 {
		t.Error("expected non-empty body")
	}
}

func TestParse_FnLiteralNoParams(t *testing.T) {
	node := parseExpr("fn() { 42 }")
	fn, ok := node.(*FnLiteral)
	if !ok {
		t.Fatalf("expected *FnLiteral, got %T", node)
	}
	if len(fn.Params) != 0 {
		t.Errorf("expected 0 params, got %d", len(fn.Params))
	}
}

func TestParse_CallExpr(t *testing.T) {
	node := parseExpr("add(1, 2)")
	call, ok := node.(*CallExpr)
	if !ok {
		t.Fatalf("expected *CallExpr, got %T", node)
	}
	id, ok := call.Fn.(*Identifier)
	if !ok {
		t.Fatalf("expected *Identifier for fn, got %T", call.Fn)
	}
	if id.Name != "add" {
		t.Errorf("expected fn=add, got %q", id.Name)
	}
	if len(call.Args) != 2 {
		t.Errorf("expected 2 args, got %d", len(call.Args))
	}
}

func TestParse_IfExpr(t *testing.T) {
	node := parseExpr("if (x < 10) { x } else { 10 }")
	ie, ok := node.(*IfExpr)
	if !ok {
		t.Fatalf("expected *IfExpr, got %T", node)
	}
	if ie.Condition == nil {
		t.Error("expected non-nil condition")
	}
	if ie.Then == nil {
		t.Error("expected non-nil then branch")
	}
	if ie.Else == nil {
		t.Error("expected non-nil else branch")
	}
}

func TestParse_IfExprNoElse(t *testing.T) {
	node := parseExpr("if (true) { 1 }")
	ie, ok := node.(*IfExpr)
	if !ok {
		t.Fatalf("expected *IfExpr, got %T", node)
	}
	if ie.Else != nil {
		t.Error("expected nil else branch when no else clause")
	}
}

func TestParse_ComparisonOps(t *testing.T) {
	for _, src := range []string{"a == b", "a != b", "a < b", "a > b", "a <= b", "a >= b"} {
		node := parseExpr(src)
		if _, ok := node.(*InfixExpr); !ok {
			t.Errorf("%q: expected *InfixExpr, got %T", src, node)
		}
	}
}
