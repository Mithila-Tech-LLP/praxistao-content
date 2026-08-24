# Task 10 — Closures and Built-ins

## Concept

### Closures

A **closure** is a function that carries its enclosing environment with it. You already implemented this in Task 09 by storing `fn.Env` in `FunctionObj`. This task verifies that complex closure patterns work correctly and introduces a new pattern: **closures as factories**.

When `makeAdder(5)` is called, it returns a new `FunctionObj` whose captured environment contains `n = 5`. When the returned function is later called with `3`, it creates a new enclosed environment on top of `n = 5`'s environment, binds `x = 3`, and evaluates `x + n` — finding `x` in the inner env and `n` in the outer captured env.

### Built-in Functions

Built-ins are native Go functions exposed to the language. They are stored in the global environment before any user code runs and implement things like `len` and `puts` that cannot be written in the language itself.

## Closure Patterns to Verify

### makeAdder

```
let makeAdder = fn(n) { fn(x) { x + n } }
let add5 = makeAdder(5)
add5(3)        // 8
add5(10)       // 15
```

`add5` closes over `n = 5`. Each call creates a fresh inner env with `x` bound to the argument.

### Counter (mutable state via closures)

```
let makeCounter = fn() {
    let count = 0
    fn() { let count = count + 1; count }
}
let counter = makeCounter()
counter()   // 1
counter()   // 2
counter()   // 3
```

Note: in our language `let` always creates a new binding, so `let count = count + 1` reads the old `count` from the outer env and creates a new binding in the inner env shadowing it. This means each call produces the right increment but the outer `count` is not mutated. The result is still correct for the test.

## Built-in Object Type

```go
type BuiltinFn func(args ...Object) Object

type BuiltinObj struct{ Fn BuiltinFn }

func (b *BuiltinObj) Type() ObjectType { return "BUILTIN" }
func (b *BuiltinObj) Inspect() string  { return "builtin function" }
```

## Registering Built-ins

Create a map of built-in functions and populate the global environment with them before evaluating any user code:

```go
var builtins = map[string]*BuiltinObj{
    "len":  {Fn: builtinLen},
    "puts": {Fn: builtinPuts},
}

func NewGlobalEnvironment() *Environment {
    env := NewEnvironment()
    for name, fn := range builtins {
        env.Set(name, fn)
    }
    return env
}
```

### len

```go
func builtinLen(args ...Object) Object {
    if len(args) != 1 {
        return newError("wrong number of arguments: len takes 1, got %d", len(args))
    }
    switch a := args[0].(type) {
    case *StringObj:
        return &NumberObj{Value: float64(len(a.Value))}
    default:
        return newError("argument to len not supported: %s", args[0].Type())
    }
}
```

### puts

```go
func builtinPuts(args ...Object) Object {
    for _, a := range args {
        fmt.Println(a.Inspect())
    }
    return NULL
}
```

## Calling Built-ins

Extend `applyFunction` to handle `*BuiltinObj`:

```go
func applyFunction(fn Object, args []Object) Object {
    switch f := fn.(type) {
    case *FunctionObj:
        extEnv := extendFunctionEnv(f, args)
        evaluated := Eval(f.Body, extEnv)
        return unwrapReturnValue(evaluated)
    case *BuiltinObj:
        return f.Fn(args...)
    default:
        return newError("not a function: %s", fn.Type())
    }
}
```

## Examples

| Expression | Result |
|-----------|--------|
| `makeAdder(5)(3)` | `NumberObj(8)` |
| `makeAdder(10)(10)` | `NumberObj(20)` |
| `len("hello")` | `NumberObj(5)` |
| `len("")` | `NumberObj(0)` |
| `len(42)` | `ErrorObj` |
| `len("a", "b")` | `ErrorObj` (wrong arg count) |
| `puts("hi")` | prints `hi`, returns `NULL` |

## Tests to Pass

1. `makeAdder(5)(3)` → `NumberObj(8)`
2. `makeAdder(10)(10)` → `NumberObj(20)`
3. Two separate adders do not share state:
   - `let add3 = makeAdder(3); let add7 = makeAdder(7); add3(10) + add7(10)` → `NumberObj(30)`
4. Counter closure produces sequential values (see pattern above)
5. `len("hello")` → `NumberObj(5)`
6. `len("")` → `NumberObj(0)`
7. `len(42)` → `ErrorObj`
8. `puts("hello")` returns `NULL` (and prints to stdout, verified separately)
9. Built-in called as first-class value: `let myLen = len; myLen("abc")` → `NumberObj(3)`

## You Now Have a Complete Interpreter

With Task 10 done, your interpreter handles:
- Number, string, and boolean literals
- Arithmetic, comparison, and logical operators
- Variables with lexical scoping
- If/else as an expression
- First-class functions with closures
- Explicit return statements
- Built-in functions

From here you could add: arrays, hash maps, while loops, more built-ins (`push`, `map`, `filter`), or a REPL (read-eval-print loop). The architecture you have built scales naturally to all of these.
