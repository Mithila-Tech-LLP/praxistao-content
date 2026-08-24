# Task 05 — Evaluator: Literals

## Concept

The **evaluator** (also called the interpreter) walks the AST and computes a value for each node. This style — walking the tree directly without compiling to bytecode — is called a **tree-walking interpreter**. It is the simplest and most readable approach, and it is exactly how early versions of many languages (Ruby, PHP) worked.

This task wires up the object system and handles the simplest case: evaluating literal values.

## Object System

Every value in the language is represented as a Go struct that implements the `Object` interface:

```go
type ObjectType string

const (
    OBJ_NUMBER = ObjectType("NUMBER")
    OBJ_STRING = ObjectType("STRING")
    OBJ_BOOL   = ObjectType("BOOL")
    OBJ_NULL   = ObjectType("NULL")
    OBJ_ERROR  = ObjectType("ERROR")
)

type Object interface {
    Type() ObjectType
    Inspect() string  // human-readable, used for printing
}
```

### Concrete Object Types

```go
type NumberObj struct{ Value float64 }
func (n *NumberObj) Type() ObjectType { return OBJ_NUMBER }
func (n *NumberObj) Inspect() string  { return strconv.FormatFloat(n.Value, 'f', -1, 64) }

type StringObj struct{ Value string }
func (s *StringObj) Type() ObjectType { return OBJ_STRING }
func (s *StringObj) Inspect() string  { return s.Value }

type BoolObj struct{ Value bool }
func (b *BoolObj) Type() ObjectType { return OBJ_BOOL }
func (b *BoolObj) Inspect() string  { return strconv.FormatBool(b.Value) }

type NullObj struct{}
func (n *NullObj) Type() ObjectType { return OBJ_NULL }
func (n *NullObj) Inspect() string  { return "null" }

type ErrorObj struct{ Message string }
func (e *ErrorObj) Type() ObjectType { return OBJ_ERROR }
func (e *ErrorObj) Inspect() string  { return "ERROR: " + e.Message }
```

Tip: create package-level singletons for the two boolean values and null so you are not allocating fresh objects on every evaluation:

```go
var (
    TRUE  = &BoolObj{Value: true}
    FALSE = &BoolObj{Value: false}
    NULL  = &NullObj{}
)
```

## Environment

The environment stores variable bindings. Define it now even though variables come in Task 07:

```go
type Environment struct {
    store map[string]Object
    outer *Environment
}

func NewEnvironment() *Environment {
    return &Environment{store: make(map[string]Object)}
}

func (e *Environment) Get(name string) (Object, bool) {
    val, ok := e.store[name]
    if !ok && e.outer != nil {
        return e.outer.Get(name)
    }
    return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
    e.store[name] = val
    return val
}
```

## Eval

```go
func Eval(node Node, env *Environment) Object
```

Use a type switch on `node`:

```go
switch n := node.(type) {
case *Program:
    return evalStatements(n.Statements, env)
case *ExpressionStatement:
    return Eval(n.Expr, env)
case *NumberLiteral:
    return &NumberObj{Value: n.Value}
case *StringLiteral:
    return &StringObj{Value: n.Value}
case *BoolLiteral:
    if n.Value { return TRUE }
    return FALSE
case *BlockStatement:
    return evalStatements(n.Stmts, env)
default:
    return newError("unknown node type: %T", node)
}
```

`evalStatements` evaluates each statement in order and returns the value of the last one.

`newError` is a helper: `func newError(format string, args ...any) *ErrorObj`.

## Helper: isError

```go
func isError(obj Object) bool {
    return obj != nil && obj.Type() == OBJ_ERROR
}
```

Propagate errors early: if `Eval` returns an `ErrorObj`, stop and return it immediately rather than continuing to evaluate the rest of the program.

## Example

```go
tokens := Tokenize(`42`)
parser := NewParser(tokens)
program := parser.ParseProgram()
env := NewEnvironment()
result := Eval(program, env)
// result is *NumberObj{Value: 42}
fmt.Println(result.Inspect()) // "42"
```

## Tests to Pass

1. `Eval(parse("42"))` → `*NumberObj{Value: 42}`
2. `Eval(parse("3.14"))` → `*NumberObj{Value: 3.14}`
3. `Eval(parse(`"hello"`))` → `*StringObj{Value: "hello"}`
4. `Eval(parse("true"))` → pointer equal to the `TRUE` singleton
5. `Eval(parse("false"))` → pointer equal to the `FALSE` singleton
6. Multiple statements: `Eval(parse("1; 2; 3"))` → `*NumberObj{Value: 3}` (last value)
7. `result.Type()` returns the correct `ObjectType` for each case
