package main

import "testing"

func evalArith(src string) Object {
	env := NewEnvironment()
	prog := NewParser(Tokenize(src)).ParseProgram()
	if len(prog.Statements) == 0 {
		return NULL_OBJ
	}
	return Eval(prog.Statements[0], env)
}

func assertNumber(t *testing.T, obj Object, expected float64) {
	t.Helper()
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T (%s)", obj, obj.Inspect())
	}
	if num.Value != expected {
		t.Errorf("expected %v, got %v", expected, num.Value)
	}
}

func assertBool(t *testing.T, obj Object, expected bool) {
	t.Helper()
	b, ok := obj.(*BoolObj)
	if !ok {
		t.Fatalf("expected *BoolObj, got %T", obj)
	}
	if b.Value != expected {
		t.Errorf("expected %v, got %v", expected, b.Value)
	}
}

func TestEval_Addition(t *testing.T) {
	assertNumber(t, evalArith("1+2"), 3)
}

func TestEval_OrderOfOperations(t *testing.T) {
	// 10/2 + 3*4 - 1 = 5 + 12 - 1 = 16
	assertNumber(t, evalArith("10/2+3*4-1"), 16)
}

func TestEval_StringConcatenation(t *testing.T) {
	obj := evalArith(`"a"+"b"`)
	s, ok := obj.(*StringObj)
	if !ok {
		t.Fatalf("expected *StringObj, got %T", obj)
	}
	if s.Value != "ab" {
		t.Errorf("expected ab, got %q", s.Value)
	}
}

func TestEval_DivisionByZero(t *testing.T) {
	obj := evalArith("10/0")
	if _, ok := obj.(*ErrorObj); !ok {
		t.Fatalf("expected *ErrorObj for division by zero, got %T", obj)
	}
}

func TestEval_PrefixMinus(t *testing.T) {
	assertNumber(t, evalArith("-5"), -5)
}

func TestEval_PrefixBangTrue(t *testing.T) {
	assertBool(t, evalArith("!true"), false)
}

func TestEval_PrefixBangFalse(t *testing.T) {
	assertBool(t, evalArith("!false"), true)
}

func TestEval_Comparisons(t *testing.T) {
	assertBool(t, evalArith("1 < 2"), true)
	assertBool(t, evalArith("2 > 3"), false)
	assertBool(t, evalArith("1 == 1"), true)
	assertBool(t, evalArith("1 != 2"), true)
	assertBool(t, evalArith("3 <= 3"), true)
	assertBool(t, evalArith("4 >= 5"), false)
}

func TestEval_Subtraction(t *testing.T) {
	assertNumber(t, evalArith("10-3"), 7)
}

func TestEval_Multiplication(t *testing.T) {
	assertNumber(t, evalArith("4*5"), 20)
}

func TestEval_Division(t *testing.T) {
	assertNumber(t, evalArith("10/2"), 5)
}
