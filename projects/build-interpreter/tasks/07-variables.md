# Task 07 — Variables

## Concept

Variables need a place to live. That place is called the **environment** (or symbol table) — a mapping from names to values. When you write `let x = 5`, the name `"x"` is bound to the value `NumberObj(5)` in the current environment. When you later write `x + 3`, the evaluator looks up `"x"` in the environment to find its value.

**Lexical scoping** means that a lookup searches the current environment first, then walks outward through parent environments until the name is found or the outermost environment is exhausted. You implemented the `Environment` struct in Task 05 — now you will put it to use.

## What to Implement

### Evaluate LetStatement

```go
case *LetStatement:
    val := Eval(n.Value, env)
    if isError(val) { return val }
    env.Set(n.Name, val)
    return val  // or return NULL — both are fine
```

### Evaluate Identifier

```go
case *Identifier:
    if val, ok := env.Get(n.Name); ok {
        return val
    }
    return newError("identifier not found: %s", n.Name)
```

That is the entire implementation. The power comes from the environment chain: `env.Get` already walks up through outer scopes (you built that in Task 05).

## Multi-Statement Programs

When you evaluate a `Program` or `BlockStatement`, you call `evalStatements` which evaluates each statement in sequence. The environment is shared across all statements in the same block, so a `let` in an earlier statement is visible to a later one:

```
let x = 5
let y = x + 3   // x is already in env
y               // lookup returns 8
```

The final statement's value is returned as the program result.

## Shadowing

If you declare the same name twice in the same scope, the second `let` overwrites the first:

```
let x = 1
let x = 2
x        // 2
```

This is fine — `env.Set` simply overwrites.

## Testing Scope Isolation

Enclosed environments will be used in Task 09 (functions). For now verify that the outer environment works correctly:

```go
outer := NewEnvironment()
outer.Set("x", &NumberObj{Value: 10})

inner := NewEnclosedEnvironment(outer)
inner.Set("y", &NumberObj{Value: 5})

// inner can see x from outer
val, _ := inner.Get("x")   // NumberObj(10)

// outer cannot see y from inner
_, ok := outer.Get("y")    // false
```

`NewEnclosedEnvironment` returns a new environment whose `outer` field points to the given parent:

```go
func NewEnclosedEnvironment(outer *Environment) *Environment {
    env := NewEnvironment()
    env.outer = outer
    return env
}
```

## Examples

| Program | Result |
|---------|--------|
| `let x = 5; x` | `NumberObj(5)` |
| `let x = 5; let y = x + 3; y` | `NumberObj(8)` |
| `let a = 2; let b = 3; a * b` | `NumberObj(6)` |
| `z` (undeclared) | `ErrorObj("identifier not found: z")` |
| `let x = 10; let x = 20; x` | `NumberObj(20)` |

## Tests to Pass

1. `"let x = 5; x"` → `NumberObj(5)`
2. `"let x = 5; let y = x + 3; y"` → `NumberObj(8)`
3. `"let a = 2; let b = 3; a * b"` → `NumberObj(6)`
4. `"let s = \"hello\"; s"` → `StringObj("hello")`
5. Evaluating an undeclared identifier returns `ErrorObj`
6. Outer environment lookup: inner env can see outer bindings
7. Inner env bindings do not leak into outer env
