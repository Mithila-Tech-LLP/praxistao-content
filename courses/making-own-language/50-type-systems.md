# Chapter 50: Type Systems — Catching Bugs at Compile Time

> "A type system is a tractable syntactic method for proving the absence of certain program behaviors by classifying phrases according to the kinds of values they compute."
> — Benjamin C. Pierce, *Types and Programming Languages*

---

When you declare `let age: int = "hello"` in Astra, the compiler rejects it before the program ever runs. That rejection is the work of the **type system** — a formal set of rules that classifies every expression in a program as belonging to a particular *type*, and then verifies that every operation in the program is applied only to values of appropriate types. Type systems are one of the most powerful tools in programming language design: they transform entire categories of runtime bugs into compile-time errors, they document code intent directly in the source, and they enable compilers to generate more efficient machine code.

In this chapter we build Astra's type checker — the phase that runs immediately after semantic analysis and annotates every node in the AST with a resolved type. We will explore the theory behind type systems (from formal type rules to the Hindley-Milner algorithm), understand why Astra makes the specific design choices it does (strong typing, local type inference, monomorphized generics, trait-based polymorphism), and implement a complete `sema/typechecker.go` module. When this chapter is done, the Astra compiler can detect and report every type error before a single line of machine code is generated.

Understanding type systems deeply will make you a better programmer in any language, not just a better compiler writer. Every time you see a type error in Rust or TypeScript, the exact rules we study here are running inside the compiler, doing exactly what we are about to implement.

---

## What We're Building

A complete static type checker for Astra. This module reads the resolved AST (from Chapter 49) and annotates every expression with its type. It validates all operations, function calls, assignments, and return statements for type correctness, and reports precise, helpful error messages.

## Table of Contents

1. What Is a Type?
2. Why Types Matter
3. Static vs Dynamic Typing
4. Strong vs Weak Typing
5. Type Inference
6. Type Rules: The Formal Language of Type Systems
7. Subtyping and Variance
8. Generics and Parametric Polymorphism
9. Ad-Hoc Polymorphism: Traits
10. Type Errors and Great Error Messages
11. Implementation: The Type Checker in Go
12. Astra Build Milestone

---

## 1. What Is a Type?

A **type** is a set of values together with the operations that are valid on those values.

Consider `int` in Astra:
- **The set of values:** all integers from −9,223,372,036,854,775,808 to 9,223,372,036,854,775,807 (64-bit signed integers)
- **The valid operations:** addition `+`, subtraction `-`, multiplication `*`, division `/`, modulo `%`, comparison `<`, `>`, `==`, `!=`, `<=`, `>=`, bitwise AND `&`, bitwise OR `|`, bitwise XOR `^`, bit shifts `<<`, `>>`

Consider `string` in Astra:
- **The set of values:** all finite sequences of Unicode code points
- **The valid operations:** concatenation `+`, indexing `[]`, length `.len()`, slicing

Consider `bool` in Astra:
- **The set of values:** `{true, false}`
- **The valid operations:** logical AND `&&`, logical OR `||`, logical NOT `!`, equality `==`, `!=`

The type system's job is to verify that every operation in a program is drawn from the valid operations of its operands' types. `"hello" - 1` is rejected because subtraction is not a valid operation on strings. `true * false` is rejected because multiplication is not valid on booleans. This sounds simple, but scaling it to a real language with functions, generics, structs, and traits requires substantial machinery.

A crucial insight: types are a compile-time abstraction. At runtime, a value is just bytes in memory. Types only exist in the compiler. By the time code generation is done, most type information has been discarded — the compiler has already checked everything and generated code that handles the right bytes in the right way.

---

## 2. Why Types Matter

**Catching bugs at compile time.** The most obvious benefit: a type error that would cause a runtime crash (`null pointer exception`, `TypeError: cannot read property of undefined`, `AttributeError`) is caught before the program ever runs. In a dynamically typed language, you might not discover the bug until a specific code path is executed — perhaps only in production under a specific set of user inputs. A static type system catches it on every developer's machine, every build.

Consider a real-world scenario. A web API handler receives a JSON body and extracts a field. In a dynamically typed language, if the field is missing or of the wrong type, you get a runtime exception. In Astra (statically typed), the compiler forces you to handle the case where the field might not be present, before the code ships.

**Enabling optimization.** When a compiler knows a value is an `int`, it can store it in a CPU register and use native integer arithmetic instructions. Without type information, the compiler must either box all values (store them behind pointers with runtime type tags) or dynamically dispatch every operation (check the type tag at runtime and then call the right operation). Static types allow the compiler to generate the most efficient possible machine code.

**Documentation that cannot go stale.** Type annotations in function signatures are a form of documentation that the compiler enforces. Unlike comments, which can become outdated and misleading, type annotations must always match the actual behavior of the code — the compiler guarantees it. When you see:

```astra
fn processPayment(amount: int, currency: string, userId: int) -> bool
```

You know exactly what this function expects and what it returns, with no ambiguity.

**Enabling IDE tooling.** Language servers (the backends that power features like autocomplete, go-to-definition, and find-all-references in editors like VS Code) rely entirely on type information. When you type `order.`, the IDE shows you the fields and methods available on the `Order` type because the type checker has already resolved what type `order` has. Without a static type system, these features are either impossible or unreliable guesses.

---

## 3. Static vs Dynamic Typing

**Static typing** means types are checked at compile time, before the program runs. If there is a type error, the compiler rejects the program and produces an error message. The program never executes. Examples: C, C++, Go, Rust, Java, Kotlin, Swift, and Astra.

**Dynamic typing** means types are checked at runtime, as the program executes. A function can receive any value; if it calls an operation that is not valid for that value's type, an exception is raised at that moment. Examples: Python, JavaScript, Ruby, PHP, Lua.

**Gradually typed** languages allow a mix: some parts of the code have static types, other parts are dynamically typed. Examples: TypeScript (compiles to JavaScript but adds static types), Python with mypy or Pyright type checking, Dart's optional types.

The trade-offs:

| Property | Static | Dynamic |
|---|---|---|
| Bug detection | Compile time (early) | Runtime (late, possibly in production) |
| Performance | Faster (no runtime type checks) | Slower (boxing, dispatch overhead) |
| Flexibility | Less (must declare types up front) | More (duck typing, late binding) |
| IDE tooling | Excellent | Limited |
| Refactoring safety | High (compiler validates all changes) | Low (may break at runtime) |
| Learning curve | Higher initially | Lower initially |

Astra is statically typed because its goal is to be a safe, performant compiled language. These goals are fundamentally better served by static typing.

### The Cost of Dynamic Typing: A Concrete Example

In Python:
```python
def double(x):
    return x * 2

double(5)        # Works: returns 10
double("hello")  # Works: returns "hellohello"
double([1, 2])   # Works: returns [1, 2, 1, 2]
```

This flexibility sounds nice. But now consider:
```python
def compute_area(width, height):
    return width * height

# Called correctly:
compute_area(10, 20)         # 200 ✓

# Called incorrectly (discovered at runtime):
compute_area("10", 20)       # TypeError: can't multiply sequence by non-int
compute_area(10, "twenty")   # Same error
compute_area(None, 20)       # TypeError: unsupported operand type(s)
```

In Astra:
```astra
fn computeArea(width: int, height: int) -> int {
    return width * height
}

computeArea("10", 20)    // Compile error: expected int, got string
computeArea(10, "twenty") // Compile error: expected int, got string
```

The errors are caught at compile time, with precise locations, before any code runs.

---

## 4. Strong vs Weak Typing

Strong vs weak typing is a separate axis from static vs dynamic typing. It describes whether a language performs **implicit type conversions**.

**Strong typing:** The compiler/runtime never silently converts a value from one type to another. If you want to use an integer where a string is expected, you must explicitly convert it. Languages with strong typing: Python (dynamically but strongly typed), Haskell, Rust, and Astra.

**Weak typing:** The compiler/runtime freely converts between types based on context, often in surprising ways. Languages with weak typing: C, JavaScript, PHP.

JavaScript's weak typing produces famous surprises:

```javascript
"5" + 1         // "51"  (1 coerced to string, then concatenated)
"5" - 1         // 4     (string coerced to number, then subtracted)
[] + []         // ""    (two arrays? becomes empty string)
[] + {}         // "[object Object]"
{} + []         // 0
true + true     // 2
false + []      // "false"
null + 1        // 1
undefined + 1   // NaN
```

These implicit conversions make programs hard to reason about. A function that receives the wrong type might produce a nonsensical result rather than an error.

Astra is **strongly typed**: there are no implicit conversions. To combine an integer and a string, you must explicitly convert:

```astra
let n = 42
let s = "The answer is " + n.toString()   // explicit conversion ✓

let s2 = "The answer is " + n            // compile error:
                                          // cannot use int as string in +
```

This strictness catches bugs and makes programs more predictable.

---

## 5. Type Inference

Writing a type annotation for every single expression would be tedious:

```astra
let x: int = 1
let y: int = 2
let z: int = x + y
let greeting: string = "hello"
let doubled: int = x * 2
```

**Type inference** allows the compiler to deduce types from context, so you don't have to write them everywhere:

```astra
let x = 1           // inferred: int
let y = 2           // inferred: int
let z = x + y       // inferred: int (because int + int = int)
let greeting = "hello"  // inferred: string
let doubled = x * 2     // inferred: int
```

### Local Type Inference (Astra's Approach)

Astra uses **local type inference**: the type of a `let` binding is inferred from its initializer expression. This is simple to implement, produces clear error messages, and covers the vast majority of cases where type annotations would be redundant.

```go
// Pseudocode: type inference for let statement
func inferLetType(stmt *ast.LetStmt) ast.Type {
    if stmt.TypeAnnotation != nil {
        return stmt.TypeAnnotation  // explicit annotation wins
    }
    if stmt.Initializer == nil {
        // error: cannot infer type without initializer
        return PoisonType
    }
    return checkExpr(stmt.Initializer)  // infer from RHS
}
```

### Hindley-Milner Type Inference (For Deep Understanding)

The gold standard of type inference is the **Hindley-Milner (HM)** algorithm, used in Haskell, ML languages, and partially in Rust. HM can infer types for *any* expression without *any* type annotations at all, even for polymorphic functions.

HM works with **type variables** — unknowns that get unified as the algorithm discovers constraints. Here is a simplified example:

```
fn add(a, b) = a + b

Step 1: Assign type variables:
  a : α
  b : β
  add : α → β → γ

Step 2: The expression `a + b` requires that + applies.
  + has type: (int, int) → int
  Constraint: α = int, β = int, γ = int

Step 3: Unify:
  a : int, b : int, add : int → int → int
```

HM produces maximally general types. The function `fn identity(x) = x` gets type `∀α. α → α` — it works for any type.

The downside of HM is error messages. When unification fails across many constraints accumulated from many locations, the error message says "cannot unify type T inferred from line 4 with type U inferred from line 12 via line 8" — confusing even for experienced programmers.

### Bidirectional Type Checking

**Bidirectional type checking** is a modern alternative used by languages like Swift and Scala. It combines two modes:

- **Checking mode (⇐):** Given a term and an *expected* type, verify the term has that type. Used when the expected type is known from context.
- **Synthesis mode (⇒):** Given a term alone, *compute* its type. Used when no expected type is known.

Bidirectional checking produces better error messages because the expected type is always known at the point where an error occurs. Astra's type checker uses a simplified form of bidirectional checking.

---

## 6. Type Rules: The Formal Language of Type Systems

Type theorists use a notation called **typing judgments** to express type rules precisely. Reading these is important for understanding academic literature on type systems.

The notation `Γ ⊢ e : T` reads: "In context Γ (the type environment, mapping variables to their types), expression `e` has type `T`."

A **typing rule** looks like:

```
    premises
────────────────  [rule name]
    conclusion
```

If the premises are true, then the conclusion is true.

### Basic Rules

**Integer literals always have type int:**

```
──────────────── [INT]
Γ ⊢ n : int
```

(No premises needed — this is always true.)

**Boolean literals have type bool:**

```
─────────────────── [BOOL-TRUE]      ──────────────────── [BOOL-FALSE]
Γ ⊢ true : bool                      Γ ⊢ false : bool
```

**Variable lookup (reads from the context):**

```
x : T ∈ Γ
────────── [VAR]
Γ ⊢ x : T
```

Reads: "If x is bound to type T in context Γ, then x has type T."

**Binary arithmetic operations:**

```
Γ ⊢ e₁ : int      Γ ⊢ e₂ : int
─────────────────────────────── [ADD]
Γ ⊢ e₁ + e₂ : int
```

**If-else expression (both branches must have the same type):**

```
Γ ⊢ cond : bool      Γ ⊢ e₁ : T      Γ ⊢ e₂ : T
──────────────────────────────────────────────────── [IF]
Γ ⊢ if cond { e₁ } else { e₂ } : T
```

This rule captures why `if condition { 1 } else { "hello" }` is a type error: the then-branch has type `int` and the else-branch has type `string`. They must be the same type `T`, and there is no single T that satisfies both.

**Function declaration:**

```
Γ, x₁:T₁, x₂:T₂, ..., xₙ:Tₙ ⊢ body : Tᵣ
──────────────────────────────────────────────────────────── [FN]
Γ ⊢ fn f(x₁:T₁, ..., xₙ:Tₙ) -> Tᵣ { body } : T₁→...→Tₙ→Tᵣ
```

The premise says: in the context extended with the parameter types, the body must have the return type `Tᵣ`.

**Function call:**

```
Γ ⊢ f : T₁ → T₂    Γ ⊢ arg : T₁
────────────────────────────────── [CALL]
Γ ⊢ f(arg) : T₂
```

The function must have a function type, and the argument must match the parameter type. The result is the return type.

### Reading a Type Derivation

A full type derivation is a tree of rule applications. For `let z = 1 + 2`:

```
──────────    ──────────
Γ ⊢ 1 : int   Γ ⊢ 2 : int
─────────────────────────────  [ADD]
Γ ⊢ 1 + 2 : int
─────────────────────────────  [LET]
Γ ⊢ let z = 1 + 2 : void
  (z : int added to Γ)
```

Every node in the derivation tree is justified by a typing rule. If no valid derivation can be constructed, the expression is ill-typed.

---

## 7. Subtyping and Variance

**Subtyping** asks: when can a value of type B be used where type A is expected? We say "B is a subtype of A" (written B <: A) if any value of type B can safely be used as a value of type A.

Classic OOP example: if `Dog` is a subclass of `Animal`, then a `Dog` can be used wherever an `Animal` is expected. `Dog <: Animal`.

### Variance: The Subtlety That Trips Up Java Programmers

Variance describes how subtyping relationships on simple types extend to composite types like `List<T>` or function types `T → U`.

**Covariance:** If `Dog <: Animal`, then `List<Dog> <: List<Animal>`. A list of dogs "is" a list of animals — seems intuitive! But Java allows this (for arrays, not for `ArrayList`) and it creates bugs:

```java
// Java (UNSOUND — this compiles but fails at runtime!)
Dog[] dogs = new Dog[1];
Animal[] animals = dogs;     // covariant array assignment, allowed
animals[0] = new Cat();      // ArrayStoreException at runtime!
```

A `Cat` is an `Animal`, so we can put it in `Animal[]` — but this `Animal[]` is actually a `Dog[]`, so storing a `Cat` breaks the type guarantee.

**Invariance:** `List<Dog>` and `List<Animal>` are unrelated types, even if `Dog <: Animal`. You cannot assign one to the other. This is what Java's generics (`ArrayList`) do, and it is safe.

**Contravariance:** Function parameters are contravariant. If `Animal <: Dog` (wait, reversed!), then a function `Animal → bool` can be used where a `Dog → bool` is expected. Reason: if you pass a `Dog` to a function expecting an `Animal`, that's fine because `Dog` is at least as capable as `Animal`.

Astra does not have subtyping at all. Instead of class hierarchies, Astra uses **traits** (similar to Rust). This completely eliminates variance problems. The trait system provides polymorphism without the pitfalls of subtype polymorphism.

---

## 8. Generics and Parametric Polymorphism

**Generics** allow you to write code that works for multiple types, parameterized over a type variable:

```astra
fn max<T: Comparable>(a: T, b: T) -> T {
    if a.cmp(b) > 0 { return a }
    return b
}

let biggerInt    = max(3, 5)        // T instantiated to int
let biggerString = max("a", "b")    // T instantiated to string
let biggerFloat  = max(1.5, 2.7)    // T instantiated to float
```

The type variable `T` stands for any type that implements the `Comparable` trait. The constraint `T: Comparable` is called a **type bound**.

### Monomorphization vs Type Erasure

There are two fundamentally different implementation strategies for generics:

**Monomorphization (Astra, Rust, C++ templates):** For each concrete instantiation of a generic function, the compiler generates a completely separate copy of the function specialized for that type.

```
// From:
fn max<T: Comparable>(a: T, b: T) -> T { ... }

// Compiler generates:
fn max_int(a: int, b: int) -> int { ... }    // for max(3, 5)
fn max_string(a: string, b: string) -> string { ... } // for max("a","b")
fn max_float(a: float, b: float) -> float { ... } // for max(1.5, 2.7)
```

Advantages: zero runtime overhead (no boxing, no vtable), full compiler optimization of each specialization.

Disadvantages: larger binary size (code duplication), longer compile times when many instantiations exist.

**Type erasure (Java, Go with interfaces, pre-Generics Java, Haskell):** The generic function is compiled once, with all type parameters replaced by a generic representation (a pointer/interface). Runtime type casts are inserted where needed.

```
// Compiled once:
fn max(a: any, b: any) -> any {
    if (a as Comparable).cmp(b) > 0 { return a }
    return b
}
```

Advantages: smaller binary, faster compile times.

Disadvantages: runtime overhead (boxing, dynamic dispatch, type casts), cannot optimize for specific types.

Astra uses **monomorphization** for the same reason Rust does: performance is a primary goal, and zero-cost abstractions require that generic code compiles to the same machine code as hand-specialized code.

---

## 9. Ad-Hoc Polymorphism: Traits

**Ad-hoc polymorphism** (also called **overloading**) means a single function name has different implementations for different types. Astra's trait system provides this:

```astra
trait Display {
    fn display(self) -> string
}

impl Display for int {
    fn display(self) -> string {
        return intToString(self)
    }
}

impl Display for bool {
    fn display(self) -> string {
        if self { return "true" }
        return "false"
    }
}

struct Point { x: int, y: int }

impl Display for Point {
    fn display(self) -> string {
        return "(" + self.x.display() + ", " + self.y.display() + ")"
    }
}
```

Now `display()` works for `int`, `bool`, and `Point`, with each type providing its own implementation.

### Trait Bounds in Generic Functions

Traits become powerful when combined with generics:

```astra
fn printAll<T: Display>(items: [T]) {
    for item in items {
        println(item.display())
    }
}

printAll([1, 2, 3])           // T = int, uses Display for int
printAll([true, false])        // T = bool, uses Display for bool
printAll([Point{1,2}, Point{3,4}]) // T = Point, uses Display for Point
```

### Static Dispatch (Monomorphization)

When the concrete type is known at compile time (which it is for generic functions in Astra), trait method calls are **statically dispatched**: the compiler inserts a direct call to the specific implementation, with no runtime overhead.

```
printAll([1, 2, 3]):
  Compiler knows T = int
  Generates: for each element, call display_for_int(element)
  No vtable, no indirection — just a direct function call
```

### Dynamic Dispatch (Trait Objects)

Sometimes you need a collection of mixed types that all implement a trait, without knowing the concrete types at compile time. This requires a **vtable** (virtual method table):

```astra
// Conceptual Astra (future feature):
let displayables: [dyn Display] = [1, true, Point{3,4}]
// Each element has a vtable pointer to its Display implementation
```

Astra initially uses only static dispatch (simpler, faster). Dynamic dispatch is a future extension.

---

## 10. Type Errors and Great Error Messages

The quality of a compiler's error messages is a major factor in developer experience. Rust is famous for having the best error messages of any systems language. Astra aims to match this quality.

Here are ten distinct type error messages Astra produces, each with context, description, and a suggestion:

**Error 1: Type Mismatch in Assignment**
```
[main.astra:4:14] type error: cannot assign value of type 'string' to variable of type 'int'
  4 | let age: int = "twenty"
    |                ^^^^^^^^ expected int, found string
  hint: try `age = 20` or convert with `int.parse("twenty")`
```

**Error 2: Binary Operation Type Mismatch**
```
[main.astra:7:18] type error: operator '+' cannot be applied to 'string' and 'int'
  7 | let result = "hello" + 42
    |              ─────── ^ ── int
    |              string
  hint: to concatenate, convert int to string: `"hello" + 42.display()`
```

**Error 3: Wrong Return Type**
```
[main.astra:12:12] type error: function 'square' declares return type 'int' but returns 'string'
  12|     return "not an int"
    |            ^^^^^^^^^^^^ found string
  note: function declared to return 'int' at line 11
```

**Error 4: Calling Non-Function**
```
[main.astra:15:1] type error: 'x' is not callable (it has type 'int', not a function type)
  15| x(5)
    | ^ 'x' declared as int at line 3
  hint: did you mean to call a function? Check if 'x' should be a function.
```

**Error 5: Condition Not Bool**
```
[main.astra:18:4] type error: condition of 'if' must be 'bool', found 'int'
  18|     if x + 1 {
    |        ^^^^^ this has type int
  hint: use a comparison: `if x + 1 > 0 {`
```

**Error 6: Mismatched If-Else Branch Types**
```
[main.astra:22:9] type error: if-else branches have incompatible types
  22|     let val = if flag { 1 } else { "one" }
    |                         ^ int    ^ string
  note: the then-branch has type 'int', the else-branch has type 'string'
  hint: make both branches the same type
```

**Error 7: Wrong Argument Type in Call**
```
[main.astra:26:14] type error: argument 2 of 'add' has wrong type
  26|     let sum = add(1, "two")
    |                      ^^^^^ expected int, found string
  note: 'add' declared at line 1: fn add(a: int, b: int) -> int
```

**Error 8: Struct Field Type Mismatch**
```
[main.astra:30:22] type error: field 'x' of struct 'Point' has type 'int', but assigned 'float'
  30|     let p = Point { x: 1.5, y: 2 }
    |                        ^^^ expected int, found float
  hint: use an integer literal: `x: 1` or `x: 2`
```

**Error 9: Trait Bound Not Satisfied**
```
[main.astra:35:12] type error: type 'Point' does not implement trait 'Comparable'
  35|     let bigger = max(p1, p2)
    |                  ^^^ T = Point, but Point: Comparable required
  hint: add `impl Comparable for Point { ... }` to your code
```

**Error 10: Array Element Type Mismatch**
```
[main.astra:40:22] type error: array elements must all have the same type
  40|     let arr = [1, 2, "three", 4]
    |                       ^^^^^^^ expected int (from element 0), found string
  hint: make all elements the same type, or use a union type
```

Each error message follows the same pattern: file, line, column; what was expected; what was found; a note with additional context; and a concrete suggestion for how to fix it.

---

## 11. Implementation: The Type Checker in Go

```go
// sema/typechecker.go

package sema

import (
    "fmt"
    "astra/ast"
    "astra/token"
)

// TypeChecker annotates every expression in the AST with a resolved type.
// It runs after the Resolver (Chapter 49).
type TypeChecker struct {
    errors      []SemanticError
    funcStack   []*ast.FunctionDecl // stack of enclosing functions
    globalScope *Scope
    current     *Scope
}

func NewTypeChecker(globalScope *Scope) *TypeChecker {
    return &TypeChecker{
        globalScope: globalScope,
        current:     globalScope,
    }
}

func (tc *TypeChecker) error(pos token.Position, format string, args ...interface{}) {
    tc.errors = append(tc.errors, SemanticError{
        Pos:     pos,
        Message: fmt.Sprintf(format, args...),
    })
}

func (tc *TypeChecker) Errors() []SemanticError { return tc.errors }
func (tc *TypeChecker) HasErrors() bool          { return len(tc.errors) > 0 }

// currentFunction returns the innermost enclosing function, or nil if at top level.
func (tc *TypeChecker) currentFunction() *ast.FunctionDecl {
    if len(tc.funcStack) == 0 {
        return nil
    }
    return tc.funcStack[len(tc.funcStack)-1]
}

// ─────────────────────────────────────────────────────────────────────
// Type representation
// ─────────────────────────────────────────────────────────────────────

// We use ast.Type values to represent types. For built-in types we use
// singleton pointers so equality checks work with ==.
var (
    TypeInt    = &ast.NamedType{Name: "int"}
    TypeFloat  = &ast.NamedType{Name: "float"}
    TypeString = &ast.NamedType{Name: "string"}
    TypeBool   = &ast.NamedType{Name: "bool"}
    TypeVoid   = &ast.NamedType{Name: "void"}
    TypeError  = &ast.NamedType{Name: "<error>"} // poison type
)

func isError(t ast.TypeExpr) bool {
    if nt, ok := t.(*ast.NamedType); ok {
        return nt.Name == "<error>"
    }
    return false
}

func typesEqual(a, b ast.TypeExpr) bool {
    if isError(a) || isError(b) {
        return true // avoid cascading errors
    }
    switch at := a.(type) {
    case *ast.NamedType:
        if bt, ok := b.(*ast.NamedType); ok {
            return at.Name == bt.Name
        }
        return false
    case *ast.ArrayType:
        if bt, ok := b.(*ast.ArrayType); ok {
            return typesEqual(at.ElementType, bt.ElementType)
        }
        return false
    case *ast.FuncType:
        bt, ok := b.(*ast.FuncType)
        if !ok || len(at.ParamTypes) != len(bt.ParamTypes) {
            return false
        }
        for i := range at.ParamTypes {
            if !typesEqual(at.ParamTypes[i], bt.ParamTypes[i]) {
                return false
            }
        }
        return typesEqual(at.ReturnType, bt.ReturnType)
    }
    return false
}

func typeString(t ast.TypeExpr) string {
    switch ty := t.(type) {
    case *ast.NamedType:
        return ty.Name
    case *ast.ArrayType:
        return "[" + typeString(ty.ElementType) + "]"
    case *ast.FuncType:
        return "fn(...) -> " + typeString(ty.ReturnType)
    }
    return "<unknown>"
}

// ─────────────────────────────────────────────────────────────────────
// Program-level checking
// ─────────────────────────────────────────────────────────────────────

func (tc *TypeChecker) CheckProgram(prog *ast.Program) {
    for _, decl := range prog.Declarations {
        tc.checkDeclaration(decl)
    }
}

func (tc *TypeChecker) checkDeclaration(decl ast.Declaration) {
    switch d := decl.(type) {
    case *ast.FunctionDecl:
        tc.checkFn(d)
    case *ast.StructDecl:
        tc.checkStruct(d)
    case *ast.VarDecl:
        tc.checkGlobalVar(d)
    }
}

// ─────────────────────────────────────────────────────────────────────
// Function checking
// ─────────────────────────────────────────────────────────────────────

func (tc *TypeChecker) checkFn(d *ast.FunctionDecl) {
    tc.funcStack = append(tc.funcStack, d)
    defer func() { tc.funcStack = tc.funcStack[:len(tc.funcStack)-1] }()

    // Check each statement in the body
    for _, stmt := range d.Body {
        tc.checkStmt(stmt)
    }

    // Check that non-void functions have a reachable return
    if d.ReturnType != nil && !typesEqual(d.ReturnType, TypeVoid) {
        if !tc.hasReturn(d.Body) {
            tc.error(d.Pos,
                "function '%s' may not return a value on all paths", d.Name)
        }
    }
}

// hasReturn is a simplified check: does the body always reach a return?
func (tc *TypeChecker) hasReturn(stmts []ast.Statement) bool {
    for _, stmt := range stmts {
        switch s := stmt.(type) {
        case *ast.ReturnStmt:
            return true
        case *ast.IfStmt:
            if s.ElseBlock != nil &&
               tc.hasReturn(s.ThenBlock) &&
               tc.hasReturn(s.ElseBlock) {
                return true
            }
        }
    }
    return false
}

func (tc *TypeChecker) checkStruct(d *ast.StructDecl) {
    // Validate field types are all known types
    for _, field := range d.Fields {
        if !tc.isValidType(field.Type) {
            tc.error(field.Pos, "unknown type '%s' for field '%s'",
                typeString(field.Type), field.Name)
        }
    }
}

func (tc *TypeChecker) checkGlobalVar(d *ast.VarDecl) {
    if d.TypeAnnotation != nil && d.Initializer != nil {
        initType := tc.checkExpr(d.Initializer)
        if !isError(initType) && !typesEqual(d.TypeAnnotation, initType) {
            tc.error(d.Pos,
                "cannot assign value of type '%s' to variable '%s' of type '%s'",
                typeString(initType), d.Name, typeString(d.TypeAnnotation))
        }
    } else if d.Initializer != nil {
        tc.checkExpr(d.Initializer) // infer type for annotation later
    }
}

// ─────────────────────────────────────────────────────────────────────
// Statement checking
// ─────────────────────────────────────────────────────────────────────

func (tc *TypeChecker) checkStmt(stmt ast.Statement) {
    switch s := stmt.(type) {
    case *ast.LetStmt:
        tc.checkLetStmt(s)
    case *ast.AssignStmt:
        tc.checkAssignStmt(s)
    case *ast.ReturnStmt:
        tc.checkReturnStmt(s)
    case *ast.IfStmt:
        tc.checkIfStmt(s)
    case *ast.WhileStmt:
        tc.checkWhileStmt(s)
    case *ast.ForStmt:
        tc.checkForStmt(s)
    case *ast.ExprStmt:
        tc.checkExpr(s.Expr)
    case *ast.BlockStmt:
        for _, inner := range s.Stmts {
            tc.checkStmt(inner)
        }
    }
}

func (tc *TypeChecker) checkLetStmt(s *ast.LetStmt) {
    var initType ast.TypeExpr
    if s.Initializer != nil {
        initType = tc.checkExpr(s.Initializer)
    }

    if s.TypeAnnotation != nil && initType != nil {
        if !isError(initType) && !typesEqual(s.TypeAnnotation, initType) {
            tc.error(s.Pos,
                "cannot assign value of type '%s' to variable '%s' of type '%s'",
                typeString(initType), s.Name, typeString(s.TypeAnnotation))
        }
        s.ResolvedType = s.TypeAnnotation
    } else if initType != nil {
        s.ResolvedType = initType // type inference
    } else if s.TypeAnnotation != nil {
        s.ResolvedType = s.TypeAnnotation
    } else {
        tc.error(s.Pos,
            "cannot determine type of '%s': no annotation or initializer", s.Name)
        s.ResolvedType = TypeError
    }

    // Update symbol table entry with resolved type
    if sym := tc.current.LookupLocal(s.Name); sym != nil {
        sym.Type = s.ResolvedType
    }
}

func (tc *TypeChecker) checkAssignStmt(s *ast.AssignStmt) {
    targetSym := tc.current.Lookup(s.Target)
    valueType := tc.checkExpr(s.Value)

    if targetSym != nil && !isError(valueType) {
        if !typesEqual(targetSym.Type, valueType) {
            tc.error(s.Pos,
                "cannot assign value of type '%s' to '%s' (type '%s')",
                typeString(valueType), s.Target, typeString(targetSym.Type))
        }
    }
}

func (tc *TypeChecker) checkReturnStmt(s *ast.ReturnStmt) {
    fn := tc.currentFunction()
    if fn == nil {
        return // caught by resolver
    }

    var returnedType ast.TypeExpr = TypeVoid
    if s.Value != nil {
        returnedType = tc.checkExpr(s.Value)
    }

    expectedReturn := fn.ReturnType
    if expectedReturn == nil {
        expectedReturn = TypeVoid
    }

    if !isError(returnedType) && !typesEqual(expectedReturn, returnedType) {
        tc.error(s.Pos,
            "function '%s' declares return type '%s' but returns '%s'",
            fn.Name, typeString(expectedReturn), typeString(returnedType))
    }
}

func (tc *TypeChecker) checkIfStmt(s *ast.IfStmt) {
    condType := tc.checkExpr(s.Condition)
    if !isError(condType) && !typesEqual(condType, TypeBool) {
        tc.error(s.Pos,
            "condition of 'if' must be 'bool', found '%s'",
            typeString(condType))
    }
    for _, stmt := range s.ThenBlock {
        tc.checkStmt(stmt)
    }
    if s.ElseBlock != nil {
        for _, stmt := range s.ElseBlock {
            tc.checkStmt(stmt)
        }
    }
}

func (tc *TypeChecker) checkWhileStmt(s *ast.WhileStmt) {
    condType := tc.checkExpr(s.Condition)
    if !isError(condType) && !typesEqual(condType, TypeBool) {
        tc.error(s.Pos,
            "condition of 'while' must be 'bool', found '%s'",
            typeString(condType))
    }
    for _, stmt := range s.Body {
        tc.checkStmt(stmt)
    }
}

func (tc *TypeChecker) checkForStmt(s *ast.ForStmt) {
    iterType := tc.checkExpr(s.Iterable)
    // iterable must be an array or range type
    if !isError(iterType) {
        if _, isArray := iterType.(*ast.ArrayType); !isArray {
            if !typesEqual(iterType, &ast.NamedType{Name: "range"}) {
                tc.error(s.Pos,
                    "'for' iterable must be an array or range, found '%s'",
                    typeString(iterType))
            }
        }
    }
    for _, stmt := range s.Body {
        tc.checkStmt(stmt)
    }
}

// ─────────────────────────────────────────────────────────────────────
// Expression type checking — the heart of the type checker
// ─────────────────────────────────────────────────────────────────────

// checkExpr computes and returns the type of an expression.
// It also stores the resolved type on the AST node for later passes.
func (tc *TypeChecker) checkExpr(expr ast.Expression) ast.TypeExpr {
    var t ast.TypeExpr
    switch e := expr.(type) {
    case *ast.IntLiteral:
        t = TypeInt
    case *ast.FloatLiteral:
        t = TypeFloat
    case *ast.StringLiteral:
        t = TypeString
    case *ast.BoolLiteral:
        t = TypeBool
    case *ast.Identifier:
        t = tc.checkIdentifier(e)
    case *ast.BinaryExpr:
        t = tc.checkBinaryExpr(e)
    case *ast.UnaryExpr:
        t = tc.checkUnaryExpr(e)
    case *ast.CallExpr:
        t = tc.checkCallExpr(e)
    case *ast.FieldAccess:
        t = tc.checkFieldAccess(e)
    case *ast.IndexExpr:
        t = tc.checkIndexExpr(e)
    case *ast.ArrayLiteral:
        t = tc.checkArrayLiteral(e)
    case *ast.StructLiteral:
        t = tc.checkStructLiteral(e)
    default:
        t = TypeError
    }

    // Attach the resolved type to the AST node
    expr.SetType(t)
    return t
}

func (tc *TypeChecker) checkIdentifier(e *ast.Identifier) ast.TypeExpr {
    if e.ResolvedSymbol == nil {
        return TypeError // resolver already reported the error
    }
    if e.ResolvedSymbol.Type == nil {
        // Type not yet resolved (might happen for forward-declared items)
        return TypeError
    }
    return e.ResolvedSymbol.Type
}

// checkBinaryExpr implements the type rules for binary operations.
// This is essentially the [ADD], [SUB], [LT], [EQ], etc. rules.
func (tc *TypeChecker) checkBinaryExpr(e *ast.BinaryExpr) ast.TypeExpr {
    leftType  := tc.checkExpr(e.Left)
    rightType := tc.checkExpr(e.Right)

    if isError(leftType) || isError(rightType) {
        return TypeError
    }

    switch e.Op {
    case "+", "-", "*", "/", "%":
        // Arithmetic: both operands must be numeric, result is same type
        if typesEqual(leftType, TypeInt) && typesEqual(rightType, TypeInt) {
            return TypeInt
        }
        if typesEqual(leftType, TypeFloat) && typesEqual(rightType, TypeFloat) {
            return TypeFloat
        }
        // Special case: string + string = string (concatenation)
        if e.Op == "+" &&
           typesEqual(leftType, TypeString) &&
           typesEqual(rightType, TypeString) {
            return TypeString
        }
        tc.error(e.Pos,
            "operator '%s' cannot be applied to '%s' and '%s'",
            e.Op, typeString(leftType), typeString(rightType))
        return TypeError

    case "<", ">", "<=", ">=":
        // Comparison: both operands must be same ordered type, result is bool
        if typesEqual(leftType, rightType) &&
           (typesEqual(leftType, TypeInt) || typesEqual(leftType, TypeFloat) ||
            typesEqual(leftType, TypeString)) {
            return TypeBool
        }
        tc.error(e.Pos,
            "operator '%s' cannot compare '%s' and '%s'",
            e.Op, typeString(leftType), typeString(rightType))
        return TypeError

    case "==", "!=":
        // Equality: both operands must be same type, result is bool
        if typesEqual(leftType, rightType) {
            return TypeBool
        }
        tc.error(e.Pos,
            "cannot compare '%s' and '%s' with '%s'",
            typeString(leftType), typeString(rightType), e.Op)
        return TypeError

    case "&&", "||":
        // Logical: both operands must be bool, result is bool
        if !typesEqual(leftType, TypeBool) {
            tc.error(e.Pos,
                "left operand of '%s' must be bool, found '%s'",
                e.Op, typeString(leftType))
            return TypeError
        }
        if !typesEqual(rightType, TypeBool) {
            tc.error(e.Pos,
                "right operand of '%s' must be bool, found '%s'",
                e.Op, typeString(rightType))
            return TypeError
        }
        return TypeBool
    }

    tc.error(e.Pos, "unknown binary operator '%s'", e.Op)
    return TypeError
}

func (tc *TypeChecker) checkUnaryExpr(e *ast.UnaryExpr) ast.TypeExpr {
    operandType := tc.checkExpr(e.Operand)
    if isError(operandType) {
        return TypeError
    }
    switch e.Op {
    case "-":
        if typesEqual(operandType, TypeInt) || typesEqual(operandType, TypeFloat) {
            return operandType
        }
        tc.error(e.Pos, "unary '-' cannot be applied to '%s'", typeString(operandType))
        return TypeError
    case "!":
        if typesEqual(operandType, TypeBool) {
            return TypeBool
        }
        tc.error(e.Pos, "unary '!' cannot be applied to '%s'", typeString(operandType))
        return TypeError
    }
    return TypeError
}

func (tc *TypeChecker) checkCallExpr(e *ast.CallExpr) ast.TypeExpr {
    sym := tc.current.Lookup(e.FuncName)
    if sym == nil {
        return TypeError // resolver caught this
    }

    // Check argument types match parameter types
    for i, arg := range e.Args {
        argType := tc.checkExpr(arg)
        if i < len(sym.ParamTypes) && !isError(argType) {
            expected := sym.ParamTypes[i]
            if !typesEqual(expected, argType) {
                tc.error(arg.GetPos(),
                    "argument %d of '%s' has wrong type: expected '%s', found '%s'",
                    i+1, e.FuncName, typeString(expected), typeString(argType))
            }
        }
    }

    if sym.ReturnType == nil {
        return TypeVoid
    }
    return sym.ReturnType
}

func (tc *TypeChecker) checkFieldAccess(e *ast.FieldAccess) ast.TypeExpr {
    objType := tc.checkExpr(e.Object)
    if isError(objType) {
        return TypeError
    }
    namedType, ok := objType.(*ast.NamedType)
    if !ok {
        tc.error(e.Pos, "field access on non-struct type '%s'", typeString(objType))
        return TypeError
    }
    structSym := tc.current.Lookup(namedType.Name)
    if structSym == nil || structSym.Kind != SymKindStruct {
        tc.error(e.Pos, "'%s' is not a struct type", namedType.Name)
        return TypeError
    }
    fieldSym, exists := structSym.Fields[e.Field]
    if !exists {
        tc.error(e.Pos, "struct '%s' has no field '%s'", namedType.Name, e.Field)
        return TypeError
    }
    return fieldSym.Type
}

func (tc *TypeChecker) checkIndexExpr(e *ast.IndexExpr) ast.TypeExpr {
    objType   := tc.checkExpr(e.Object)
    indexType := tc.checkExpr(e.Index)

    if isError(objType) || isError(indexType) {
        return TypeError
    }
    if !typesEqual(indexType, TypeInt) {
        tc.error(e.Pos,
            "array index must be 'int', found '%s'", typeString(indexType))
        return TypeError
    }
    if arrType, ok := objType.(*ast.ArrayType); ok {
        return arrType.ElementType
    }
    tc.error(e.Pos, "cannot index into non-array type '%s'", typeString(objType))
    return TypeError
}

func (tc *TypeChecker) checkArrayLiteral(e *ast.ArrayLiteral) ast.TypeExpr {
    if len(e.Elements) == 0 {
        // Empty array: type cannot be inferred without annotation
        if e.TypeAnnotation == nil {
            tc.error(e.Pos, "cannot infer type of empty array without type annotation")
            return TypeError
        }
        return &ast.ArrayType{ElementType: e.TypeAnnotation}
    }

    firstType := tc.checkExpr(e.Elements[0])
    for i, elem := range e.Elements[1:] {
        elemType := tc.checkExpr(elem)
        if !isError(elemType) && !typesEqual(firstType, elemType) {
            tc.error(elem.GetPos(),
                "array element %d has type '%s', expected '%s' (from element 0)",
                i+1, typeString(elemType), typeString(firstType))
            return TypeError
        }
    }
    return &ast.ArrayType{ElementType: firstType}
}

func (tc *TypeChecker) checkStructLiteral(e *ast.StructLiteral) ast.TypeExpr {
    structSym := tc.current.Lookup(e.StructName)
    if structSym == nil || structSym.Kind != SymKindStruct {
        tc.error(e.Pos, "unknown struct type '%s'", e.StructName)
        return TypeError
    }
    for _, field := range e.Fields {
        fieldSym, exists := structSym.Fields[field.Name]
        if !exists {
            tc.error(field.Pos, "struct '%s' has no field '%s'",
                e.StructName, field.Name)
            continue
        }
        valueType := tc.checkExpr(field.Value)
        if !isError(valueType) && !typesEqual(fieldSym.Type, valueType) {
            tc.error(field.Pos,
                "field '%s' of struct '%s' has type '%s', but assigned '%s'",
                field.Name, e.StructName,
                typeString(fieldSym.Type), typeString(valueType))
        }
    }
    return &ast.NamedType{Name: e.StructName}
}

func (tc *TypeChecker) isValidType(t ast.TypeExpr) bool {
    switch ty := t.(type) {
    case *ast.NamedType:
        return isBuiltinType(ty.Name) || tc.current.Lookup(ty.Name) != nil
    case *ast.ArrayType:
        return tc.isValidType(ty.ElementType)
    }
    return false
}
```

---

## 🔨 Astra Build Milestone

### Type Checking the Power Function

Given this Astra program:

```astra
fn power(base: int, exp: int) -> int {
    let result = 1
    while exp > 0 {
        result = result * base
        exp = exp - 1
    }
    return result
}

fn main() {
    let answer = power(2, 10)
    println(answer.display())
}
```

The type checker produces these annotations:

```
power:
  base → int (parameter)
  exp  → int (parameter)
  result → int (inferred from literal 1)
  exp > 0 → bool (int > int = bool) ✓
  result * base → int (int * int = int) ✓
  exp - 1 → int (int - int = int) ✓
  return result → int (matches declared return type int) ✓

main:
  power(2, 10) → int (matches return type of power) ✓
  answer → int (inferred from call result)
  answer.display() → string (Display trait method for int) ✓
```

All types check out. No errors.

---

## Exercises

1. **Implement Numeric Promotion:** Add implicit widening from `int` to `float` when one operand is `int` and the other is `float` in arithmetic operations. This is the one implicit conversion Astra decides to allow (as a pragmatic choice). How does this affect the `checkBinaryExpr` logic? What does the resulting type become?

2. **Void Function Return Check:** Currently, a `void` function (no return type annotation) that contains `return 42` would not produce an error. Extend `checkReturnStmt` to detect and report this: a void function must not return a value.

3. **Array Type Inference from Context:** In `let xs: [int] = []`, the empty array literal should be accepted because the context (the type annotation `[int]`) provides the element type. Add support for this by threading an optional "expected type" through `checkExpr` as a second parameter (the bidirectional approach).

4. **Struct Completeness Check:** When creating a struct literal, all fields must be provided. Extend `checkStructLiteral` to detect missing fields and report an error listing which fields were not initialized.

5. **Type Rule Derivation:** Write out the full type derivation tree (in the Γ ⊢ e : T notation) for the expression `if x > 0 { x * 2 } else { 0 }` where `x: int`. Which rules apply at each step? What type does the whole expression have?

6. **Generic Function Type Checking:** Describe (in pseudocode or Go) how you would extend the type checker to handle a simple generic function:
   ```astra
   fn first<T>(arr: [T]) -> T { return arr[0] }
   ```
   At each call site, how do you determine the concrete type `T`? What constraint does `arr[0]` impose on `T`?

7. **Trait Implementation Verification:** When `impl Display for Point` appears, verify that the implementation provides all methods declared in the `Display` trait, and that each method's signature matches exactly. What data structures do you need? Where does this check run?

8. **Union Types:** Astra does not currently have union types, but many modern languages do. Describe how you would add a simple union type `int | string | bool` to the type system. How does type checking for binary operations change? What new rules are needed?

---

## Summary

| Concept | Key Idea |
|---|---|
| Type | A set of values plus the operations valid on them |
| Static typing | Types checked at compile time; bugs caught before runtime |
| Strong typing | No implicit conversions; operations must match types exactly |
| Type inference | Compiler deduces types from context; reduces annotation burden |
| Hindley-Milner | Full inference algorithm; powerful but complex error messages |
| Bidirectional checking | Mix of synthesis (compute type) and checking (verify type) modes |
| Type rules (Γ ⊢ e : T) | Formal notation for type system axioms and derivation rules |
| Subtyping | When B <: A, a B can be used where A is expected |
| Variance | How subtyping extends through composite types (covariant/contravariant/invariant) |
| Generics | Type-parameterized code; Astra uses monomorphization for zero overhead |
| Traits | Ad-hoc polymorphism; multiple types implement a shared interface |
| Monomorphization | Compiler generates one copy per concrete type instantiation |
| Poison type | Error type that propagates without causing cascading false errors |
