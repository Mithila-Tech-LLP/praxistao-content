package main

import "testing"

func evalProg(src string) Object {
	env := NewEnvironment()
	prog := NewParser(Tokenize(src)).ParseProgram()
	return evalProgram(prog, env)
}

func TestVariables_LetAndLookup(t *testing.T) {
	obj := evalProg("let x = 5; x")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 5 {
		t.Errorf("expected 5, got %v", num.Value)
	}
}

func TestVariables_MultiStep(t *testing.T) {
	// let x = 5; let y = x + 3; y  → 8
	obj := evalProg("let x = 5; let y = x + 3; y")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 8 {
		t.Errorf("expected 8, got %v", num.Value)
	}
}

func TestVariables_UndefinedVariable(t *testing.T) {
	obj := evalProg("foobar")
	if _, ok := obj.(*ErrorObj); !ok {
		t.Fatalf("expected *ErrorObj for undefined variable, got %T", obj)
	}
}

func TestVariables_ReferencingOtherVar(t *testing.T) {
	obj := evalProg("let a = 10; let b = a * 2; b")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 20 {
		t.Errorf("expected 20, got %v", num.Value)
	}
}

func TestVariables_LetReturnsNull(t *testing.T) {
	// LetStatement itself evaluates to NullObj
	env := NewEnvironment()
	prog := NewParser(Tokenize("let x = 1")).ParseProgram()
	obj := Eval(prog.Statements[0], env)
	if _, ok := obj.(*NullObj); !ok {
		t.Fatalf("expected *NullObj from LetStatement eval, got %T", obj)
	}
}

func TestVariables_ErrorPropagates(t *testing.T) {
	obj := evalProg("let x = 1/0; x")
	if _, ok := obj.(*ErrorObj); !ok {
		t.Fatalf("expected *ErrorObj to propagate from bad let value, got %T", obj)
	}
}
