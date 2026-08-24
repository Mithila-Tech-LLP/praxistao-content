package main

import "testing"

func evalSource(src string) Object {
	tokens := Tokenize(src)
	prog := NewParser(tokens).ParseProgram()
	env := NewEnvironment()
	if len(prog.Statements) == 0 {
		return NULL_OBJ
	}
	return Eval(prog.Statements[0], env)
}

func TestEval_NumberLiteral(t *testing.T) {
	obj := evalSource("42")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 42 {
		t.Errorf("expected 42, got %v", num.Value)
	}
}

func TestEval_FloatLiteral(t *testing.T) {
	obj := evalSource("3.14")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 3.14 {
		t.Errorf("expected 3.14, got %v", num.Value)
	}
}

func TestEval_StringLiteral(t *testing.T) {
	obj := evalSource(`"hello"`)
	str, ok := obj.(*StringObj)
	if !ok {
		t.Fatalf("expected *StringObj, got %T", obj)
	}
	if str.Value != "hello" {
		t.Errorf("expected hello, got %q", str.Value)
	}
}

func TestEval_BoolLiteralTrue(t *testing.T) {
	obj := evalSource("true")
	b, ok := obj.(*BoolObj)
	if !ok {
		t.Fatalf("expected *BoolObj, got %T", obj)
	}
	if !b.Value {
		t.Error("expected true")
	}
}

func TestEval_BoolLiteralFalse(t *testing.T) {
	obj := evalSource("false")
	b, ok := obj.(*BoolObj)
	if !ok {
		t.Fatalf("expected *BoolObj, got %T", obj)
	}
	if b.Value {
		t.Error("expected false")
	}
}

func TestEval_UnknownNodeReturnsNull(t *testing.T) {
	// An identifier by itself (unresolved) should return NullObj
	// since the evaluator only handles literals in this task.
	obj := Eval(&Identifier{Name: "x"}, NewEnvironment())
	if _, ok := obj.(*NullObj); !ok {
		t.Fatalf("expected *NullObj for unknown node, got %T", obj)
	}
}

func TestEval_ObjectTypes(t *testing.T) {
	if evalSource("42").Type() != "NUMBER" {
		t.Error("NumberObj.Type() should be NUMBER")
	}
	if evalSource(`"x"`).Type() != "STRING" {
		t.Error("StringObj.Type() should be STRING")
	}
	if evalSource("true").Type() != "BOOL" {
		t.Error("BoolObj.Type() should be BOOL")
	}
}
