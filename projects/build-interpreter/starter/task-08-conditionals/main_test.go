package main

import "testing"

func evalCond(src string) Object {
	env := NewEnvironment()
	prog := NewParser(Tokenize(src)).ParseProgram()
	return evalProgram(prog, env)
}

func TestIf_TrueBranch(t *testing.T) {
	obj := evalCond("if (1 < 2) { 10 } else { 20 }")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 10 {
		t.Errorf("expected 10, got %v", num.Value)
	}
}

func TestIf_FalseBranch(t *testing.T) {
	obj := evalCond("if (2 < 1) { 10 } else { 20 }")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 20 {
		t.Errorf("expected 20, got %v", num.Value)
	}
}

func TestIf_FalseConditionNoElse(t *testing.T) {
	obj := evalCond("if (false) { 1 }")
	if _, ok := obj.(*NullObj); !ok {
		t.Fatalf("expected *NullObj when condition is false and no else, got %T", obj)
	}
}

func TestIf_TrueConditionWithExpr(t *testing.T) {
	obj := evalCond("if (true) { 3+3 }")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 6 {
		t.Errorf("expected 6, got %v", num.Value)
	}
}

func TestIf_NumberIsTruthy(t *testing.T) {
	// Non-zero numbers are truthy
	obj := evalCond("if (42) { 1 } else { 2 }")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 1 {
		t.Errorf("expected 1 (number is truthy), got %v", num.Value)
	}
}

func TestIf_NullIsFalsy(t *testing.T) {
	// if we pass a null-producing expression, it should take the else branch
	// We'll use a variable that evaluates to null: evaluate let and use result
	// Actually let's just test isTruthy directly
	if isTruthy(NULL_OBJ) {
		t.Error("null should be falsy")
	}
	if !isTruthy(TRUE_OBJ) {
		t.Error("true should be truthy")
	}
	if isTruthy(FALSE_OBJ) {
		t.Error("false should be falsy")
	}
	if !isTruthy(&NumberObj{Value: 1}) {
		t.Error("number should be truthy")
	}
}

func TestBlock_ReturnsLastValue(t *testing.T) {
	obj := evalCond("if (true) { 1; 2; 3 }")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 3 {
		t.Errorf("expected 3 (last value in block), got %v", num.Value)
	}
}

func TestIf_ConditionWithVariable(t *testing.T) {
	obj := evalCond("let x = 5; if (x > 3) { x * 2 } else { 0 }")
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 10 {
		t.Errorf("expected 10, got %v", num.Value)
	}
}
