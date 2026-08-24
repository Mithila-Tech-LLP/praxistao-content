package main

import "testing"

func evalFn(src string) Object {
	env := NewEnvironment()
	prog := NewParser(Tokenize(src)).ParseProgram()
	return evalProgram(prog, env)
}

func TestFunction_CallWithArgs(t *testing.T) {
	obj := evalFn("let add = fn(a,b){ a+b }; add(3,4)")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T (%s)", obj, obj.Inspect())
	}
	if num.Value != 7 {
		t.Errorf("expected 7, got %v", num.Value)
	}
}

func TestFunction_EarlyReturn(t *testing.T) {
	// return x*2 should fire before x+1 is reached
	obj := evalFn("let f = fn(x){ return x*2; x+1 }; f(5)")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T (%s)", obj, obj.Inspect())
	}
	if num.Value != 10 {
		t.Errorf("expected 10 (early return x*2), got %v", num.Value)
	}
}

func TestFunction_NoReturnUsesLastExpr(t *testing.T) {
	obj := evalFn("let f = fn(x){ x * x }; f(4)")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 16 {
		t.Errorf("expected 16, got %v", num.Value)
	}
}

func TestFunction_FnLiteralCreatesObject(t *testing.T) {
	env := NewEnvironment()
	prog := NewParser(Tokenize("fn(x){ x }")).ParseProgram()
	obj := evalProgram(prog, env)
	if _, ok := obj.(*FunctionObj); !ok {
		t.Fatalf("expected *FunctionObj, got %T", obj)
	}
}

func TestFunction_ReturnFromNested(t *testing.T) {
	src := `
let f = fn(x) {
    if (x > 0) {
        return x * 2
    }
    return 0
}
f(5)
`
	obj := evalFn(src)
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 10 {
		t.Errorf("expected 10, got %v", num.Value)
	}
}

func TestFunction_CallWithNoArgs(t *testing.T) {
	obj := evalFn("let f = fn() { 42 }; f()")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 42 {
		t.Errorf("expected 42, got %v", num.Value)
	}
}
