# Task 09 — Functions

## Concept

Functions are the heart of any real language. In our language, functions are **first-class values**: you can assign them to variables, pass them as arguments, and return them from other functions. They are created with a `fn` literal and called like any other expression.

The critical detail is **closures**: a function captures the environment that existed when it was created, not the one that exists when it is called. This is what makes `makeAdder` work.

## The Function Object

```go
type FunctionObj struct {
    Params []string
    Body   *BlockStatement
    Env    *Environment  // captured environment — this is the closure
}

func (f *FunctionObj) Type() ObjectType { return OBJ_FUNCTION }
func (f *FunctionObj) Inspect() string  { return "fn(...){...}" }
```

## Evaluating FnLiteral

```go
case *FnLiteral:
    return &FunctionObj{
        Params: n.Params,
        Body:   n.Body,
        Env:    env,   // capture the current environment
    }
```

That is all. The closure is captured by simply storing the current `env` pointer.

## Evaluating CallExpr

```go
case *CallExpr:
    fn := Eval(n.Fn, env)
    if isError(fn) { return fn }

    args := evalArgs(n.Args, env)
    if len(args) == 1 && isError(args[0]) { return args[0] }

    return applyFunction(fn, args)
```

### evalArgs

Evaluate each argument expression in the current env, propagating errors:

```go
func evalArgs(args []Expression, env *Environment) []Object {
    var result []Object
    for _, a := range args {
        val := Eval(a, env)
        if isError(val) { return []Object{val} }
        result = append(result, val)
    }
    return result
}
```

### applyFunction

```go
func applyFunction(fn Object, args []Object) Object {
    f, ok := fn.(*FunctionObj)
    if !ok {
        return newError("not a function: %s", fn.Type())
    }

    extEnv := extendFunctionEnv(f, args)
    evaluated := Eval(f.Body, extEnv)
    return unwrapReturnValue(evaluated)
}
```

### extendFunctionEnv

Create a new enclosed environment from the **function's captured env** (not the caller's env). Bind each parameter name to the corresponding argument value:

```go
func extendFunctionEnv(fn *FunctionObj, args []Object) *Environment {
    env := NewEnclosedEnvironment(fn.Env)  // fn.Env, not caller's env!
    for i, param := range fn.Params {
        env.Set(param, args[i])
    }
    return env
}
```

## Return Statements

A `return` statement exits a function early with a value. Implement it with a **sentinel object** that travels up the call stack until it hits a function boundary, where it is unwrapped.

```go
type ReturnValue struct{ Value Object }
func (r *ReturnValue) Type() ObjectType { return "RETURN_VALUE" }
func (r *ReturnValue) Inspect() string  { return r.Value.Inspect() }
```

Evaluate `ReturnStatement`:

```go
case *ReturnStatement:
    val := Eval(n.Value, env)
    if isError(val) { return val }
    return &ReturnValue{Value: val}
```

In `BlockStatement`, stop as soon as you see a `ReturnValue` or error:

```go
case *BlockStatement:
    var result Object
    for _, stmt := range n.Stmts {
        result = Eval(stmt, env)
        if result != nil {
            rt := result.Type()
            if rt == "RETURN_VALUE" || rt == OBJ_ERROR {
                return result  // propagate without unwrapping
            }
        }
    }
    return result
```

`unwrapReturnValue` extracts the inner value at the function call site:

```go
func unwrapReturnValue(obj Object) Object {
    if rv, ok := obj.(*ReturnValue); ok {
        return rv.Value
    }
    return obj
}
```

For `*Program` evaluation, also unwrap return values (a top-level `return` ends the program):

```go
case *Program:
    return evalProgram(n.Statements, env)

func evalProgram(stmts []Statement, env *Environment) Object {
    var result Object
    for _, stmt := range stmts {
        result = Eval(stmt, env)
        switch r := result.(type) {
        case *ReturnValue:
            return r.Value  // unwrap here
        case *ErrorObj:
            return result
        }
    }
    return result
}
```

## Examples

```
let add = fn(a, b) { a + b }
add(3, 4)                           // 7

let identity = fn(x) { x }
identity(42)                        // 42

let applyTwice = fn(f, x) { f(f(x)) }
let double = fn(x) { x * 2 }
applyTwice(double, 3)               // 12

let factorial = fn(n) {
    if (n < 2) { return 1 }
    n * factorial(n - 1)
}
factorial(5)                        // 120
```

## Tests to Pass

1. `"let add = fn(a, b) { a + b }; add(3, 4)"` → `NumberObj(7)`
2. `"let identity = fn(x) { x }; identity(42)"` → `NumberObj(42)`
3. Calling a non-function → `ErrorObj`
4. Wrong number of args: covered by Go panicking on slice index (or add explicit check)
5. `"let f = fn(x) { return x * 2 }; f(5)"` → `NumberObj(10)` (explicit return)
6. Early return: `"fn(x) { return 1; 99 }(0)"` → `NumberObj(1)`
7. Higher-order: `applyTwice(double, 3)` → `NumberObj(12)`
8. Recursive: `factorial(5)` → `NumberObj(120)`
