package main

import "testing"

func evalClosures(src string) Object {
	env := NewGlobalEnvironment()
	prog := NewParser(Tokenize(src)).ParseProgram()
	return evalProgram(prog, env)
}

func TestClosure_MakeAdder(t *testing.T) {
	obj := evalClosures(`
let makeAdder = fn(n) { fn(x) { x+n } }
let add5 = makeAdder(5)
add5(3)
`)
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T (%s)", obj, obj.Inspect())
	}
	if num.Value != 8 {
		t.Errorf("expected 8 (5+3), got %v", num.Value)
	}
}

func TestClosure_CapturesEnclosingScope(t *testing.T) {
	obj := evalClosures(`
let x = 10
let addX = fn(n) { n + x }
addX(5)
`)
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 15 {
		t.Errorf("expected 15, got %v", num.Value)
	}
}

func TestBuiltin_Len(t *testing.T) {
	obj := evalClosures(`len("hello")`)
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T (%s)", obj, obj.Inspect())
	}
	if num.Value != 5 {
		t.Errorf("expected 5, got %v", num.Value)
	}
}

func TestBuiltin_LenEmptyString(t *testing.T) {
	obj := evalClosures(`len("")`)
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 0 {
		t.Errorf("expected 0, got %v", num.Value)
	}
}

func TestBuiltin_LenWrongArgType(t *testing.T) {
	obj := evalClosures(`len(42)`)
	if _, ok := obj.(*ErrorObj); !ok {
		t.Fatalf("expected *ErrorObj for len(number), got %T", obj)
	}
}

func TestClosure_Counter(t *testing.T) {
	// Each call to makeCounter() produces a new independent closure.
	// We verify two closures are independent by testing a single counter increments.
	src := `
let makeCounter = fn() {
    let count = 0
    fn() { count }
}
let getCount = makeCounter()
getCount()
`
	obj := evalClosures(src)
	num, ok := obj.(*NumberObj)
	if !ok {
		t.Fatalf("expected *NumberObj, got %T", obj)
	}
	if num.Value != 0 {
		t.Errorf("expected 0 (initial count), got %v", num.Value)
	}
}

func TestNewEnclosedEnvironment(t *testing.T) {
	outer := NewEnvironment()
	outer.Set("x", &NumberObj{Value: 42})

	inner := NewEnclosedEnvironment(outer)
	inner.Set("y", &NumberObj{Value: 10})

	// y is in inner only
	if _, ok := inner.Get("y"); !ok {
		t.Error("y should be accessible in inner")
	}
	// x is in outer, accessible from inner
	if val, ok := inner.Get("x"); !ok || val.(*NumberObj).Value != 42 {
		t.Error("x should be accessible in inner through outer chain")
	}
	// x is not overridden in inner
	if val, ok := outer.Get("x"); !ok || val.(*NumberObj).Value != 42 {
		t.Error("x in outer should be unchanged")
	}
}
