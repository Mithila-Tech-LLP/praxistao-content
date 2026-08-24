# Task 08 — Conditionals

## Concept

An `if` expression evaluates a condition and, depending on whether it is **truthy**, evaluates either the `then` branch or the `else` branch. Unlike some languages where `if` is a statement, in our language `if` is an **expression** — it produces a value that can be used directly:

```
let result = if (x > 0) { "positive" } else { "non-positive" }
```

This task also completes the evaluation of `BlockStatement` with proper value propagation.

## Truthiness

In our language only two values are falsy: `false` (the boolean) and `null`. Everything else — including `0`, `""`, and any number or string — is truthy.

```go
func isTruthy(obj Object) bool {
    switch o := obj.(type) {
    case *NullObj:
        return false
    case *BoolObj:
        return o.Value
    default:
        return true
    }
}
```

This is a deliberate design choice. You could make `0` and `""` falsy, but it tends to cause surprising bugs. Explicit is better.

## Evaluating IfExpr

```go
case *IfExpr:
    cond := Eval(n.Cond, env)
    if isError(cond) { return cond }

    if isTruthy(cond) {
        return Eval(n.Then, env)
    } else if n.Else != nil {
        return Eval(n.Else, env)
    } else {
        return NULL
    }
```

When there is no else branch and the condition is falsy, the expression evaluates to `NULL`.

## Evaluating BlockStatement

A block is a sequence of statements enclosed in `{ }`. Evaluate each statement in order and return the value of the last one:

```go
case *BlockStatement:
    var result Object
    for _, stmt := range n.Stmts {
        result = Eval(stmt, env)
        if isError(result) { return result }
        // ReturnValue propagation handled in Task 09
    }
    if result == nil { return NULL }
    return result
```

## Program vs Block

`evalStatements` for `*Program` and `*BlockStatement` behave slightly differently once `return` is added in Task 09. For now they are the same — evaluate all statements, return the last value.

## Examples

| Expression | Result |
|-----------|--------|
| `if (true) { 10 }` | `NumberObj(10)` |
| `if (false) { 10 }` | `NULL` |
| `if (1 < 2) { 10 } else { 20 }` | `NumberObj(10)` |
| `if (1 > 2) { 10 } else { 20 }` | `NumberObj(20)` |
| `if (1) { "yes" }` | `StringObj("yes")` |
| `if (0) { "yes" }` | `StringObj("yes")` (0 is truthy!) |
| `if (null) { "yes" }` | `NULL` |

The `0 is truthy` case is the intentional design described above — make sure your tests explicitly cover it.

## Nested Conditionals

```
let x = 5
if (x > 3) {
    if (x > 10) {
        "big"
    } else {
        "medium"
    }
} else {
    "small"
}
```

Expected result: `StringObj("medium")`.

## Block as Last Value

```
let x = if (true) { let a = 1; let b = 2; a + b }
x
```

Expected: `NumberObj(3)`. The block evaluates all its statements and returns the last value (`a + b`).

## Tests to Pass

1. `"if (true) { 10 }"` → `NumberObj(10)`
2. `"if (false) { 10 }"` → `NULL`
3. `"if (false) { 10 } else { 20 }"` → `NumberObj(20)`
4. `"if (1 < 2) { 10 } else { 20 }"` → `NumberObj(10)`
5. `"if (0) { 42 }"` → `NumberObj(42)` (0 is truthy)
6. `"if (null) { 42 }"` — but we have no `null` literal yet; skip or test via a missing branch
7. Nested if/else resolves to the correct branch
8. Block returns last value: `"if (true) { 1; 2; 3 }"` → `NumberObj(3)`
9. Error in condition propagates: `"if (1 + \"a\") { 10 }"` → `ErrorObj`
