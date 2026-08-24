# Task 06 — Evaluator: Arithmetic and Comparison

## Concept

With literals working, you can now extend `Eval` to handle **expressions** — prefix and infix operations. This covers arithmetic, string concatenation, comparisons, and logical negation.

The key insight is that you evaluate the operands first, then apply the operation to the resulting objects. This is called **eager evaluation** (as opposed to lazy evaluation, where operands are computed only when needed).

## Prefix Expressions

Extend the `Eval` type switch:

```go
case *PrefixExpr:
    right := Eval(n.Right, env)
    if isError(right) { return right }
    return evalPrefixExpr(n.Op, right)
```

### evalPrefixExpr

| Op | Operand type | Result |
|----|-------------|--------|
| `-` | `*NumberObj` | `*NumberObj{Value: -v}` |
| `-` | anything else | `ErrorObj{"operand must be a number"}` |
| `!` | `*BoolObj` | flipped boolean singleton |
| `!` | `*NullObj` | `TRUE` (null is falsy) |
| `!` | anything else | `FALSE` (all other values are truthy) |

## Infix Expressions

```go
case *InfixExpr:
    left := Eval(n.Left, env)
    if isError(left) { return left }
    right := Eval(n.Right, env)
    if isError(right) { return right }
    return evalInfixExpr(n.Op, left, right)
```

### evalInfixExpr

Dispatch on operand types:

```go
func evalInfixExpr(op string, left, right Object) Object {
    switch {
    case left.Type() == OBJ_NUMBER && right.Type() == OBJ_NUMBER:
        return evalNumberInfix(op, left.(*NumberObj).Value, right.(*NumberObj).Value)
    case left.Type() == OBJ_STRING && right.Type() == OBJ_STRING:
        return evalStringInfix(op, left.(*StringObj).Value, right.(*StringObj).Value)
    case op == "==":
        return nativeBoolObj(left == right)  // pointer equality for singletons
    case op == "!=":
        return nativeBoolObj(left != right)
    default:
        return newError("type mismatch: %s %s %s", left.Type(), op, right.Type())
    }
}
```

`nativeBoolObj(b bool) *BoolObj` returns the `TRUE` or `FALSE` singleton.

### evalNumberInfix

| Op | Result |
|----|--------|
| `+` | sum |
| `-` | difference |
| `*` | product |
| `/` | quotient — **return error if right == 0** |
| `<` | bool |
| `>` | bool |
| `<=` | bool |
| `>=` | bool |
| `==` | bool |
| `!=` | bool |
| other | `ErrorObj{"unknown operator: NUMBER <op> NUMBER"}` |

### evalStringInfix

| Op | Result |
|----|--------|
| `+` | concatenated `StringObj` |
| `==` | bool (value equality) |
| `!=` | bool (value equality negated) |
| other | `ErrorObj{"unknown operator: STRING <op> STRING"}` |

## Division by Zero

```go
if right == 0 {
    return newError("division by zero")
}
```

This should return an `ErrorObj`, not panic.

## Error Propagation Pattern

Always check `isError` immediately after evaluating a sub-expression and short-circuit:

```go
left := Eval(n.Left, env)
if isError(left) { return left }
```

This ensures an error anywhere in an expression bubbles up cleanly.

## Examples

| Expression | Result |
|-----------|--------|
| `1 + 2` | `NumberObj(3)` |
| `10 - 3 * 2` | `NumberObj(4)` |
| `10 / 2 + 3 * 4 - 1` | `NumberObj(16)` |
| `"hello" + " world"` | `StringObj("hello world")` |
| `1 < 2` | `TRUE` |
| `5 >= 10` | `FALSE` |
| `true == true` | `TRUE` |
| `true == false` | `FALSE` |
| `-5` | `NumberObj(-5)` |
| `!true` | `FALSE` |
| `!false` | `TRUE` |
| `!!true` | `TRUE` |
| `1 / 0` | `ErrorObj("division by zero")` |

## Tests to Pass

1. All examples in the table above.
2. `"hello" + " " + "world"` → `StringObj("hello world")`
3. `"hello" - "world"` → `ErrorObj` with type mismatch message
4. `1 + "a"` → `ErrorObj` with type mismatch message
5. `-"hello"` → `ErrorObj` (operand must be a number)
6. Division by zero returns `ErrorObj`, not a panic
7. `2 * 3 + 4 * 5` → `NumberObj(26)` (precedence correct)
