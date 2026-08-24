# Chapter 53: Designing the Astra Language — The Complete Specification

> "A language that doesn't affect the way you think about programming is not worth knowing." — Alan Perlis

---

## Overview

You have spent the previous seven volumes learning how programming languages work from the inside out. You understand grammars, parsing, abstract syntax trees, type systems, intermediate representations, and code generation. Now it is time to put all of that knowledge to work.

In this chapter we design Astra — a real, compilable programming language. This is not a toy. By the end of Volume 8, you will have built a working compiler called `astrac` that takes `.as` source files and produces native executables that run on your machine.

This chapter is the **Astra Language Specification** — the single authoritative reference document for everything the language can do. Every decision a compiler writer makes traces back to a document like this one. When your type checker asks "can I add an int to a float?" the answer comes from the specification. When your parser asks "does `else if` require braces?" the specification decides.

Read this chapter carefully. Return to it often. The remaining chapters of Volume 8 will implement every rule written here.

---

## Table of Contents

1. Design Philosophy
2. Language Goals
3. Lexical Structure
4. Types
5. Declarations
6. Statements
7. Expressions
8. Standard Library Overview
9. Module System
10. Naming Conventions
11. Keywords
12. Design Tradeoffs
13. Formal EBNF Grammar

---

## 1. Design Philosophy

Astra is built on four pillars:

**Simplicity.** Every feature in Astra must earn its place. If a feature can be expressed clearly using existing features, it does not get its own syntax. Astra has no operator overloading, no implicit conversions, no hidden control flow.

**Readability.** Code is read far more often than it is written. Astra's syntax is designed so that a programmer who has never seen Astra can understand what a function does in thirty seconds. There are no arcane symbols, no context-dependent meanings for punctuation.

**Safety.** Astra eliminates entire classes of bugs by construction. There is no null — use `Option<T>`. There are no unchecked exceptions — use `Result<T, E>`. Arrays are bounds-checked at runtime. Integer overflow is a compile-time error in debug mode.

**Performance.** Astra compiles to native machine code. There is no interpreter, no virtual machine. The compiler optimizes aggressively. Garbage collection uses a modern tri-color mark-and-sweep algorithm with generational hints.

---

## 2. Language Goals

Astra v1.0 is designed to be:

| Property | Description |
|---|---|
| **Compiled** | Source files compile to native executables |
| **Statically typed** | All types known at compile time |
| **Garbage collected** | No manual memory management |
| **Expressive** | Common tasks require little boilerplate |
| **Fast to compile** | Under 1 second for programs under 10,000 lines |
| **Fast to run** | Competitive with C for compute-heavy code |

What Astra v1.0 does **not** include (deferred to v2.0):
- Generics (type parameters)
- Async/await
- Closures that capture mutable state
- Macros

These omissions are intentional. A smaller, correct language beats a larger, broken one.

---

## 3. Lexical Structure

### 3.1 Character Set

Astra source files are UTF-8 encoded. Identifiers may contain Unicode letters and digits, but keywords and built-in type names are ASCII-only.

### 3.2 Comments

```astra
// This is a single-line comment. Everything to end of line is ignored.

/* This is a block comment.
   It can span multiple lines.
   Block comments do NOT nest. */
```

### 3.3 Whitespace

Spaces, tabs, and newlines are whitespace. Whitespace separates tokens but has no other semantic meaning. Astra does **not** use significant indentation (unlike Python).

### 3.4 Identifiers

An identifier begins with a letter or underscore, followed by any number of letters, digits, or underscores.

```
identifier = letter { letter | digit | "_" }
letter     = "A".."Z" | "a".."z" | "_"
digit      = "0".."9"
```

Valid identifiers: `x`, `count`, `snake_case`, `_private`, `myVar42`
Invalid identifiers: `42var`, `my-var`, `fn` (reserved keyword)

### 3.5 Integer Literals

```astra
let decimal = 42
let hex     = 0xFF        // hexadecimal
let binary  = 0b1010_1010 // binary with _ separator
let big     = 1_000_000   // underscore separators allowed
```

### 3.6 Float Literals

```astra
let pi    = 3.14159
let e     = 2.718_281
let small = 0.001
```

Float literals require at least one digit before and after the decimal point.

### 3.7 String Literals

Strings are delimited by double quotes. The following escape sequences are recognized:

| Escape | Meaning |
|---|---|
| `\n` | Newline (LF) |
| `\t` | Tab |
| `\r` | Carriage return |
| `\\` | Backslash |
| `\"` | Double quote |
| `\0` | Null byte |

```astra
let greeting = "Hello, World!\n"
let path     = "C:\\Users\\alice"
let quote    = "She said \"hello\""
```

### 3.8 Boolean Literals

The keywords `true` and `false` are the only boolean literals.

### 3.9 Operators and Punctuation

```
+   -   *   /   %           arithmetic
==  !=  <   >   <=  >=      comparison
&&  ||  !                   logical
=                           assignment
->                          return type arrow
..                          range (inclusive of start, exclusive of end)
::                          path separator (module access)
.                           field access / method call
,                           separator
;                           statement terminator (optional in Astra)
:                           type annotation
{   }                       block delimiters
(   )                       grouping / call arguments
[   ]                       index access
```

---

## 4. Types

### 4.1 Primitive Types

| Type | Description | Size |
|---|---|---|
| `int` | 64-bit signed integer | 8 bytes |
| `float` | 64-bit IEEE 754 floating point | 8 bytes |
| `bool` | Boolean true/false | 1 byte |
| `string` | UTF-8 immutable string | pointer + length |
| `void` | No value (function return only) | 0 bytes |

### 4.2 Compound Types

**Struct types** are user-defined records:

```astra
struct Point {
    x: float
    y: float
}
```

**Array types** (fixed size, planned for v1.1):

```astra
let arr: [int; 10] = [0; 10]
```

**Function types:**

```astra
let f: fn(int, int) -> int = add
```

### 4.3 Generic Types (Standard Library Only)

Two generic types are provided by the standard library for safety:

**`Option<T>`** — a value that may or may not be present:

```astra
let x: Option<int> = Some(42)
let y: Option<int> = None
```

**`Result<T, E>`** — a value that is either success or error:

```astra
let r: Result<int, string> = Ok(42)
let e: Result<int, string> = Err("something went wrong")
```

### 4.4 Type Compatibility

Astra has **no implicit type conversion**. Every conversion must be explicit:

```astra
let n: int   = 42
let f: float = n as float        // explicit cast required
let s: string = n.to_string()   // method call required
```

---

## 5. Declarations

### 5.1 Function Declarations

```astra
fn function_name(param1: Type1, param2: Type2) -> ReturnType {
    // body
}
```

Functions with no return value omit the `->` clause:

```astra
fn greet(name: string) {
    println("Hello, " + name)
}
```

Functions **must** have explicit parameter types. Return types are required when the function returns a value.

The entry point is always a function named `main` with no parameters and no return type:

```astra
fn main() {
    // program starts here
}
```

### 5.2 Variable Declarations

```astra
let name = "Alice"          // type inferred from right-hand side
let age: int = 30           // explicit type annotation
let pi: float = 3.14159     // explicit type annotation
```

Variables declared with `let` are **immutable by default** in Astra v2.0. In v1.0, all variables are mutable for simplicity.

Constants use `const` and must be known at compile time:

```astra
const MAX_SIZE: int = 1024
const PI: float = 3.14159
```

### 5.3 Struct Declarations

```astra
struct Person {
    name: string
    age:  int
}
```

Fields are listed one per line. Each field has a name and a type, separated by `:`.

Struct literals create instances:

```astra
let p = Person { name: "Alice", age: 30 }
```

### 5.4 Impl Blocks

Methods are defined in `impl` blocks separate from the struct definition:

```astra
impl Person {
    fn new(name: string, age: int) -> Person {
        return Person { name: name, age: age }
    }

    fn greet(self) {
        println("Hello, I am " + self.name)
    }

    fn birthday(self) {
        self.age = self.age + 1
    }
}
```

The first parameter `self` refers to the instance. It is always passed by reference (the compiler handles this). Omitting `self` makes the method static (called on the type, not an instance).

### 5.5 Import Declarations

```astra
import math
import "io/file"
import math::sqrt
```

Imports must appear at the top of the file, before any other declarations.

---

## 6. Statements

### 6.1 Expression Statements

Any expression followed by a newline (or optional semicolon) is a statement:

```astra
print("hello")
x = x + 1
do_work()
```

### 6.2 Variable Declaration Statements

```astra
let x = 42
let name: string = "Alice"
```

### 6.3 If / Else

```astra
if condition {
    // then branch
}

if condition {
    // then branch
} else {
    // else branch
}

if condition1 {
    // ...
} else if condition2 {
    // ...
} else {
    // ...
}
```

The condition must be a `bool` expression. Braces are **required** — no single-statement shorthand.

### 6.4 For Loops

```astra
for i in 0..10 {
    print(i.to_string())
}
```

`0..10` is a range from 0 (inclusive) to 10 (exclusive). The loop variable `i` is bound fresh in each iteration and has type `int`.

### 6.5 While Loops

```astra
while condition {
    // body
}
```

### 6.6 Return Statements

```astra
return
return value
return expression + 1
```

`return` with no value is used in void functions. `return value` is used in value-returning functions. The type of `value` must match the function's declared return type.

### 6.7 Block Statements

A block `{ ... }` introduces a new lexical scope. Variables declared inside are not visible outside.

---

## 7. Expressions

### 7.1 Operator Precedence

From highest to lowest:

| Level | Operators | Associativity |
|---|---|---|
| 7 | `()` `.` `[]` (call, field, index) | Left |
| 6 | `!` `-` (unary) | Right |
| 5 | `*` `/` `%` | Left |
| 4 | `+` `-` | Left |
| 3 | `<` `>` `<=` `>=` | Left |
| 2 | `==` `!=` | Left |
| 1 | `&&` | Left |
| 0 | `\|\|` | Left |

Assignment `=` has the lowest precedence and is right-associative.

### 7.2 Arithmetic

```astra
let sum  = a + b
let diff = a - b
let prod = a * b
let quot = a / b    // integer division if both are int
let rem  = a % b
```

Integer division truncates toward zero. Division by zero is a runtime panic.

### 7.3 Comparison

```astra
let eq  = a == b
let neq = a != b
let lt  = a < b
let gt  = a > b
let lte = a <= b
let gte = a >= b
```

Comparison operators return `bool`. Both operands must have the same type.

### 7.4 Logical

```astra
let and = a && b    // short-circuit: b not evaluated if a is false
let or  = a || b    // short-circuit: b not evaluated if a is true
let not = !a
```

### 7.5 Function Calls

```astra
print("hello")
add(3, 4)
person.greet()
Person::new("Alice", 30)
```

Arguments are evaluated left-to-right. The number and types of arguments must exactly match the function's parameter list.

### 7.6 Field Access

```astra
let name = person.name
let len  = str.length
```

Field access on a struct returns the value of the named field. The type checker verifies the field exists on the struct type.

### 7.7 Method Calls

```astra
person.greet()
42.to_string()
"hello".length()
```

Methods are functions defined in `impl` blocks. When called on a value, the value is passed as `self`.

### 7.8 Range Expressions

```astra
let r = 0..10    // range from 0 to 9 (exclusive end)
```

Ranges are only valid as the right-hand side of `for ... in`.

---

## 8. Standard Library Overview

The Astra standard library is available without any import for the following:

### Built-in Functions

```astra
print(s: string)          // print without newline
println(s: string)        // print with newline appended
```

### string Methods

```astra
s.length() -> int
s.contains(sub: string) -> bool
s.starts_with(prefix: string) -> bool
s.ends_with(suffix: string) -> bool
s.to_upper() -> string
s.to_lower() -> string
s.trim() -> string
s.split(sep: string) -> [string]
s.replace(old: string, new: string) -> string
```

### int Methods

```astra
n.to_string() -> string
n.to_float() -> float
n.abs() -> int
```

### float Methods

```atml
f.to_string() -> string
f.to_int() -> int       // truncates
f.abs() -> float
f.floor() -> float
f.ceil() -> float
f.round() -> float
```

### Imported Modules (v1.0 plan)

| Module | Purpose |
|---|---|
| `math` | sqrt, pow, sin, cos, log, etc. |
| `time` | current time, sleep, duration |
| `file` | read/write files |
| `json` | encode/decode JSON |
| `http` | simple HTTP client |

---

## 9. Module System

Each `.as` file is a module. The module name is the filename without the extension.

```astra
// File: math_utils.as
fn square(x: int) -> int {
    return x * x
}
```

```astra
// File: main.as
import math_utils

fn main() {
    let s = math_utils::square(5)
    println(s.to_string())
}
```

Functions are private by default. Use `pub` to export:

```astra
pub fn square(x: int) -> int {
    return x * x
}
```

---

## 10. Naming Conventions

These are enforced by the compiler as warnings (future: errors):

| Entity | Convention | Example |
|---|---|---|
| Variables | `snake_case` | `user_count`, `file_path` |
| Functions | `snake_case` | `parse_input()`, `compute_sum()` |
| Constants | `SCREAMING_SNAKE_CASE` | `MAX_SIZE`, `PI` |
| Structs | `PascalCase` | `Person`, `HttpRequest` |
| Impl methods | `snake_case` | `person.get_name()` |
| Modules | `snake_case` | `import math_utils` |

---

## 11. Keywords

The following words are reserved and cannot be used as identifiers:

```
fn        let       const     if        else
for       while     in        return    struct
impl      import    true      false     match
enum      trait     pub       self      as
and       or        not       break     continue
```

Reserved for future use (v2.0+):

```
async     await     yield     type      where
unsafe    extern    macro     use       mod
```

---

## 12. Design Tradeoffs

### Why No Null?

Null references are famously described as the "billion-dollar mistake." In Astra, variables cannot hold a null value. When a value might be absent, use `Option<T>`:

```astra
fn find_user(id: int) -> Option<Person> {
    // ...
}

let result = find_user(42)
// Must explicitly handle both cases:
match result {
    Some(person) -> person.greet()
    None         -> println("User not found")
}
```

This forces the programmer to handle the absent case. The compiler guarantees that if you have a `Person`, it is a real `Person`.

### Why No Exceptions?

Exception handling creates hidden control flow. A function call might or might not throw — the type signature does not tell you. Astra uses `Result<T, E>` for operations that can fail:

```astra
fn read_file(path: string) -> Result<string, string> {
    // ...
}

let content = read_file("data.txt")
match content {
    Ok(text)  -> println(text)
    Err(msg)  -> println("Error: " + msg)
}
```

Errors are values. They appear in type signatures. They cannot be accidentally ignored.

### Why No Inheritance?

Inheritance creates tight coupling between types. Adding a field to a parent class can break all subclasses. Astra uses **traits** for polymorphism:

```astra
trait Greetable {
    fn greet(self)
}

impl Greetable for Person {
    fn greet(self) {
        println("Hi, I'm " + self.name)
    }
}

impl Greetable for Robot {
    fn greet(self) {
        println("BEEP BOOP I AM " + self.designation)
    }
}

fn greet_all(things: [Greetable]) {
    for thing in things {
        thing.greet()
    }
}
```

Traits are like interfaces. A type implements a trait by providing the required methods in an `impl` block. There is no inheritance hierarchy, no fragile base class problem.

### Why Garbage Collection Over Manual Memory?

Manual memory management (like C's `malloc`/`free`) is error-prone. Use-after-free and double-free bugs are the source of most security vulnerabilities. Astra uses garbage collection to eliminate these entire classes of bugs.

The tradeoff is occasional GC pauses. For Astra's target use cases (web servers, command-line tools, data processing), these pauses are acceptable. Real-time applications requiring sub-millisecond guarantees should use a different language.

---

## 13. Formal EBNF Grammar

This is the complete, formal grammar of Astra v1.0. The parser in Chapter 55 implements exactly this grammar.

```ebnf
(* Top Level *)
program         = { import_decl } { declaration } ;

import_decl     = "import" ( string_lit | identifier { "::" identifier } ) ;

declaration     = fn_decl
                | struct_decl
                | impl_decl
                | const_decl ;

(* Function Declaration *)
fn_decl         = [ "pub" ] "fn" identifier "(" [ param_list ] ")" [ "->" type ] block ;
param_list      = param { "," param } ;
param           = identifier ":" type ;

(* Struct Declaration *)
struct_decl     = [ "pub" ] "struct" identifier "{" { field_decl } "}" ;
field_decl      = identifier ":" type ;

(* Impl Block *)
impl_decl       = "impl" identifier "{" { fn_decl } "}" ;

(* Const Declaration *)
const_decl      = "const" identifier ":" type "=" expr ;

(* Types *)
type            = "int"
                | "float"
                | "bool"
                | "string"
                | "void"
                | identifier
                | "fn" "(" [ type_list ] ")" "->" type
                | "Option" "<" type ">"
                | "Result" "<" type "," type ">"
                | "[" type "]"
                | "[" type ";" integer_lit "]" ;

type_list       = type { "," type } ;

(* Statements *)
block           = "{" { statement } "}" ;

statement       = let_stmt
                | return_stmt
                | if_stmt
                | for_stmt
                | while_stmt
                | break_stmt
                | continue_stmt
                | expr_stmt
                | block ;

let_stmt        = "let" identifier [ ":" type ] "=" expr ;

return_stmt     = "return" [ expr ] ;

if_stmt         = "if" expr block [ "else" ( if_stmt | block ) ] ;

for_stmt        = "for" identifier "in" expr block ;

while_stmt      = "while" expr block ;

break_stmt      = "break" ;

continue_stmt   = "continue" ;

expr_stmt       = expr ;

(* Expressions — Pratt precedence *)
expr            = assignment ;

assignment      = ( call_expr "." identifier "=" assignment )
                | ( identifier "=" assignment )
                | logical_or ;

logical_or      = logical_and { "||" logical_and } ;

logical_and     = equality { "&&" equality } ;

equality        = comparison { ( "==" | "!=" ) comparison } ;

comparison      = addition { ( "<" | ">" | "<=" | ">=" ) addition } ;

addition        = multiplication { ( "+" | "-" ) multiplication } ;

multiplication  = unary { ( "*" | "/" | "%" ) unary } ;

unary           = ( "!" | "-" ) unary | postfix ;

postfix         = primary { postfix_op } ;

postfix_op      = "(" [ arg_list ] ")"        (* function call *)
                | "." identifier              (* field access *)
                | "." identifier "(" [ arg_list ] ")"  (* method call *)
                | "[" expr "]"               (* index *)
                | ".." expr                  (* range — only in for context *) ;

arg_list        = expr { "," expr } ;

primary         = integer_lit
                | float_lit
                | string_lit
                | "true"
                | "false"
                | identifier
                | identifier "{" field_init_list "}"  (* struct literal *)
                | identifier "::" identifier          (* path expression *)
                | "(" expr ")"
                | range_expr ;

range_expr      = expr ".." expr ;

field_init_list = field_init { "," field_init } ;
field_init      = identifier ":" expr ;

(* Literals *)
integer_lit     = decimal_lit | hex_lit | binary_lit ;
decimal_lit     = digit { digit | "_" } ;
hex_lit         = "0x" hex_digit { hex_digit | "_" } ;
binary_lit      = "0b" bin_digit { bin_digit | "_" } ;

float_lit       = decimal_lit "." decimal_lit ;

string_lit      = '"' { string_char } '"' ;
string_char     = any_char_except_dquote_and_backslash
                | "\\" escape_seq ;
escape_seq      = "n" | "t" | "r" | "\\" | '"' | "0" ;

(* Lexical primitives *)
identifier      = letter { letter | digit | "_" } ;
letter          = "A".."Z" | "a".."z" | "_" ;
digit           = "0".."9" ;
hex_digit       = digit | "A".."F" | "a".."f" ;
bin_digit       = "0" | "1" ;
```

---

## Summary Table

| Feature | Status in v1.0 | Notes |
|---|---|---|
| Functions | Yes | `fn name(params) -> type { }` |
| Variables | Yes | `let x = value` |
| Constants | Yes | `const X: T = value` |
| Structs | Yes | `struct Name { field: type }` |
| Impl blocks | Yes | Methods on structs |
| If/else | Yes | Condition must be bool |
| For loops | Yes | `for i in range` only |
| While loops | Yes | Standard while |
| String interpolation | No | Use concatenation for v1.0 |
| Generics | No | Deferred to v2.0 |
| Closures | No | Deferred to v2.0 |
| Async/await | No | Deferred to v2.0 |
| Pattern matching | Partial | `match` on Option/Result |
| Traits | Partial | Defined, impl supported |
| Enums | Partial | Basic enum values |
| Modules | Yes | File-based, `import` keyword |
| Null safety | Yes | `Option<T>` instead of null |
| Error handling | Yes | `Result<T, E>` |

---

## Exercises

1. **Grammar Extension**: The current grammar does not support `break` with a label (like `break 'outer` in Rust). Write the EBNF rules that would add labeled break support to Astra. What complications does this introduce for the parser?

2. **Type System Analysis**: Astra has no implicit conversions. Write a table showing every pair of primitive types (int, float, bool, string) and whether they can be combined with `+`, `-`, `*`, `/`. For each "no" entry, explain why.

3. **Spec Reading Exercise**: According to the EBNF grammar, is `let x = if true { 1 } else { 2 }` valid Astra syntax? Trace through the grammar rules to find out. What change would you need to make to the grammar to allow if-expressions (if as a value, not just a statement)?

4. **Design Comparison**: Compare Astra's error handling approach (Result<T, E>) with Go's approach (multiple return values: `value, err := doSomething()`). What are the advantages and disadvantages of each? Which would you prefer as a programmer?

5. **Standard Library Design**: Design the API for an Astra `http` module. Write the function signatures (not implementations) for: making a GET request, making a POST request with a body, and handling the response. Use Result<T, E> appropriately.

6. **Naming Convention Enforcement**: The spec says naming conventions are enforced as compiler warnings. Write the algorithm (in pseudocode) that the compiler would use to check whether a function name follows snake_case convention. What edge cases must you handle?

7. **Grammar Ambiguity Hunt**: Find a potential ambiguity in the EBNF grammar above. Hint: look at struct literal syntax (`Name { field: value }`) and how the parser might confuse it with a block statement after an expression. How would you resolve this ambiguity?

8. **Module System Design**: The current module system is file-based. Design a package system where multiple files can belong to the same package (like Go's `package main`). Write the specification for how this would work, including the `package` keyword syntax and how imports would change.

---

This specification is the foundation everything else in Volume 8 is built on. In the next chapter, we begin building the Astra compiler by implementing the lexer — the first stage that turns raw text into the structured tokens the parser needs.
