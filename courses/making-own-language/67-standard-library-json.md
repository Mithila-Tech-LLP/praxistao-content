# Chapter 67: Standard Library — JSON

> "JSON is the assembly language of the modern web." — Douglas Crockford (paraphrased)

---

## Overview

JSON — JavaScript Object Notation — is the lingua franca of modern software. It is how REST APIs transmit data. It is how configuration files are written. It is how services talk to each other. It is how databases return query results. Any language that cannot read and write JSON cannot participate in the modern software ecosystem.

But here is something fascinating: JSON has its own grammar, and parsing it requires the same techniques you used to build Astra's own compiler. JSON parsing is recursive descent parsing applied to a much simpler grammar. Writing a JSON parser from scratch is an excellent exercise that reinforces everything you learned in the parsing chapters — and this time you will write the entire thing yourself, in about 250 lines of Go.

This chapter covers:

**The JSON specification** — understanding exactly what is and is not valid JSON, including edge cases that trip up many parsers.

**The complete recursive descent JSON parser** — we will implement `parseValue`, `parseObject`, `parseArray`, `parseString`, `parseNumber` step by step, with full Unicode escape support.

**JSON marshaling** — converting Astra structs and values to JSON text, using Go reflection to map struct fields to JSON keys automatically.

**JSON unmarshaling** — parsing JSON text directly into typed Astra structs, with type checking and error reporting.

**The streaming parser** — for multi-gigabyte JSON files, a DOM parser that loads the entire thing into memory is not an option. We will build a SAX-style streaming parser.

After this chapter, Astra programs can interoperate with any JSON-speaking service in the world.

---

## What We're Building

```
stdlib/
  json/
    json.go        ← main JSON package (~400 lines)
    parser.go      ← recursive descent parser (~250 lines)
    marshal.go     ← marshaling (Astra → JSON) (~150 lines)
    unmarshal.go   ← unmarshaling (JSON → Astra struct) (~150 lines)
```

---

## Table of Contents

1. The JSON Specification — exactly what is valid JSON
2. JSON Values in Astra — the Value type
3. The Recursive Descent Parser — complete implementation
4. String Parsing — escape sequences and Unicode
5. Number Parsing — integers and floats
6. Marshaling — Astra values to JSON text
7. Unmarshaling — JSON text to Astra structs
8. The Streaming Parser — handling large files
9. Error Handling
10. Real-World Example — config file pipeline
11. Build Milestone
12. Exercises

---

## 1. The JSON Specification

JSON is defined by ECMA-404 and RFC 8259. The grammar is remarkably simple:

```
value
  = object
  | array
  | string
  | number
  | "true"
  | "false"
  | "null"

object
  = "{" "}"
  | "{" members "}"

members
  = pair
  | pair "," members

pair
  = string ":" value

array
  = "[" "]"
  | "[" elements "]"

elements
  = value
  | value "," elements

string
  = '"' characters '"'

characters
  = ""
  | character characters

character
  = any Unicode except '"' and '\'
  | '\' escape

escape
  = '"' | '\' | '/' | 'b' | 'f' | 'n' | 'r' | 't' | 'u' hex hex hex hex

number
  = integer fraction exponent

integer
  = digit | onenine digits | '-' digit | '-' onenine digits

fraction
  = "" | '.' digits

exponent
  = "" | 'E' sign digits | 'e' sign digits
```

This grammar is **LL(1)**: to know which production to use, you only need to look at the next character:
- `{` → object
- `[` → array
- `"` → string
- `-` or `0-9` → number
- `t` → true
- `f` → false
- `n` → null

This is exactly why recursive descent works so cleanly for JSON.

### Valid JSON Examples

```json
null
true
false
42
3.14
-0.5
1.5e10
"hello"
"hello\nworld"
"A"
[]
[1, 2, 3]
{}
{"key": "value"}
{"a": 1, "b": [true, false, null]}
{"nested": {"deeply": {"value": 42}}}
```

### Common JSON Traps

```json
// NOT valid JSON (these are common mistakes):
{key: "value"}        // keys MUST be quoted strings
{'key': 'value'}      // single quotes are NOT allowed
[1, 2, 3,]            // trailing commas are NOT allowed
{                     // comments are NOT allowed
  // server config
  "port": 8080
}
undefined             // undefined is not a JSON value
NaN                   // NaN is not a JSON value  
Infinity              // Infinity is not a JSON value
01234                 // leading zeros are not allowed in numbers
```

---

## 2. JSON Values in Astra

Before writing the parser, we need a type that can represent any JSON value:

```astra
// The JSON value type — a tagged union
enum JsonValue {
    Null
    Bool(bool)
    Int(int)
    Float(float)
    Str(string)
    Array(List<JsonValue>)
    Object(Map<string, JsonValue>)
}
```

In Go, we implement this with an interface and concrete types:

```go
// stdlib/json/value.go

package astra_json

import (
    "fmt"
    "strings"
)

// ValueType enumerates the possible JSON value types.
type ValueType int

const (
    TypeNull ValueType = iota
    TypeBool
    TypeInt
    TypeFloat
    TypeString
    TypeArray
    TypeObject
)

// Value represents any JSON value.
type Value struct {
    kind    ValueType
    boolVal bool
    intVal  int64
    floatVal float64
    strVal  string
    arrVal  []*Value
    objVal  map[string]*Value
    // We preserve insertion order for objects using a separate keys slice.
    objKeys []string
}

// ─── Null ──────────────────────────────────────────────────────────────────

var Null = &Value{kind: TypeNull}

func IsNull(v *Value) bool { return v == nil || v.kind == TypeNull }

// ─── Constructors ──────────────────────────────────────────────────────────

func BoolValue(b bool) *Value     { return &Value{kind: TypeBool, boolVal: b} }
func IntValue(n int64) *Value     { return &Value{kind: TypeInt, intVal: n} }
func FloatValue(f float64) *Value { return &Value{kind: TypeFloat, floatVal: f} }
func StringValue(s string) *Value { return &Value{kind: TypeString, strVal: s} }

func ArrayValue(items []*Value) *Value {
    return &Value{kind: TypeArray, arrVal: items}
}

func ObjectValue(keys []string, vals map[string]*Value) *Value {
    return &Value{kind: TypeObject, objKeys: keys, objVal: vals}
}

// ─── Type Accessors ────────────────────────────────────────────────────────

func (v *Value) Type() ValueType { return v.kind }
func (v *Value) IsBool() bool   { return v != nil && v.kind == TypeBool }
func (v *Value) IsInt() bool    { return v != nil && v.kind == TypeInt }
func (v *Value) IsFloat() bool  { return v != nil && v.kind == TypeFloat }
func (v *Value) IsString() bool { return v != nil && v.kind == TypeString }
func (v *Value) IsArray() bool  { return v != nil && v.kind == TypeArray }
func (v *Value) IsObject() bool { return v != nil && v.kind == TypeObject }

// ─── Value Accessors ───────────────────────────────────────────────────────

func (v *Value) AsBool() bool {
    if v == nil { return false }
    switch v.kind {
    case TypeBool:   return v.boolVal
    case TypeInt:    return v.intVal != 0
    case TypeString: return v.strVal != "" && v.strVal != "false"
    default:         return false
    }
}

func (v *Value) AsInt() int64 {
    if v == nil { return 0 }
    switch v.kind {
    case TypeInt:   return v.intVal
    case TypeFloat: return int64(v.floatVal)
    case TypeString:
        var n int64
        fmt.Sscan(v.strVal, &n)
        return n
    case TypeBool:
        if v.boolVal { return 1 }
        return 0
    default: return 0
    }
}

func (v *Value) AsFloat() float64 {
    if v == nil { return 0 }
    switch v.kind {
    case TypeFloat: return v.floatVal
    case TypeInt:   return float64(v.intVal)
    case TypeString:
        var f float64
        fmt.Sscan(v.strVal, &f)
        return f
    default: return 0
    }
}

func (v *Value) AsString() string {
    if v == nil { return "" }
    switch v.kind {
    case TypeString: return v.strVal
    case TypeInt:    return fmt.Sprintf("%d", v.intVal)
    case TypeFloat:  return fmt.Sprintf("%g", v.floatVal)
    case TypeBool:
        if v.boolVal { return "true" }
        return "false"
    case TypeNull:   return "null"
    default:         return v.Stringify()
    }
}

func (v *Value) AsArray() []*Value {
    if v == nil || v.kind != TypeArray { return nil }
    return v.arrVal
}

// ─── Object Accessors ──────────────────────────────────────────────────────

// Get returns the value for the given key, or nil if not found.
func (v *Value) Get(key string) *Value {
    if v == nil || v.kind != TypeObject { return nil }
    return v.objVal[key]
}

// GetString returns the string value of the given key, or "" if not found.
func (v *Value) GetString(key string) string {
    child := v.Get(key)
    if child == nil { return "" }
    return child.AsString()
}

// GetInt returns the int value of the given key, or 0 if not found.
func (v *Value) GetInt(key string) int64 {
    child := v.Get(key)
    if child == nil { return 0 }
    return child.AsInt()
}

// GetBool returns the bool value of the given key, or false if not found.
func (v *Value) GetBool(key string) bool {
    child := v.Get(key)
    if child == nil { return false }
    return child.AsBool()
}

// Has returns true if the object contains the given key.
func (v *Value) Has(key string) bool {
    if v == nil || v.kind != TypeObject { return false }
    _, ok := v.objVal[key]
    return ok
}

// Set sets a key-value pair in an object.
func (v *Value) Set(key string, val *Value) {
    if v == nil || v.kind != TypeObject { return }
    if _, exists := v.objVal[key]; !exists {
        v.objKeys = append(v.objKeys, key) // preserve insertion order
    }
    v.objVal[key] = val
}

// Path navigates nested objects using dot-separated keys.
// config.Path("server.host") is equivalent to config.Get("server").Get("host").
func (v *Value) Path(dotPath string) *Value {
    parts := strings.Split(dotPath, ".")
    current := v
    for _, part := range parts {
        if current == nil { return nil }
        current = current.Get(part)
    }
    return current
}

// NewObject creates an empty JSON object value.
func NewObject() *Value {
    return &Value{
        kind:    TypeObject,
        objVal:  make(map[string]*Value),
        objKeys: []string{},
    }
}

// NewArray creates an empty JSON array value.
func NewArray() *Value {
    return &Value{kind: TypeArray, arrVal: []*Value{}}
}
```

---

## 3. The Recursive Descent Parser — Complete Implementation

Here is the full parser. Notice its structural similarity to Astra's own parser from Chapter 55:

```go
// stdlib/json/parser.go

package astra_json

import (
    "fmt"
    "strconv"
    "strings"
    "unicode"
    "unicode/utf8"
)

// parser holds the state for a single parse operation.
type parser struct {
    input []rune  // input decoded as Unicode code points
    pos   int     // current position in input
    line  int     // current line number (for error messages)
    col   int     // current column number (for error messages)
}

// ParseError describes a JSON syntax error.
type ParseError struct {
    Message string
    Line    int
    Col     int
    Pos     int
}

func (e *ParseError) Error() string {
    return fmt.Sprintf("json parse error at line %d, col %d: %s", e.Line, e.Col, e.Message)
}

// Parse parses a JSON string and returns the root Value.
func Parse(input string) (*Value, error) {
    p := &parser{
        input: []rune(input),
        pos:   0,
        line:  1,
        col:   1,
    }
    p.skipWhitespace()
    val, err := p.parseValue()
    if err != nil {
        return nil, err
    }
    p.skipWhitespace()
    if p.pos < len(p.input) {
        return nil, p.errorf("unexpected character after JSON value: %q", p.current())
    }
    return val, nil
}

// ─── Core Parser ───────────────────────────────────────────────────────────────

// parseValue dispatches to the appropriate parser based on the current character.
func (p *parser) parseValue() (*Value, error) {
    if p.pos >= len(p.input) {
        return nil, p.errorf("unexpected end of input")
    }

    switch p.current() {
    case '{':
        return p.parseObject()
    case '[':
        return p.parseArray()
    case '"':
        s, err := p.parseString()
        if err != nil { return nil, err }
        return StringValue(s), nil
    case 't':
        return p.parseLiteral("true", BoolValue(true))
    case 'f':
        return p.parseLiteral("false", BoolValue(false))
    case 'n':
        return p.parseLiteral("null", Null)
    case '-', '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
        return p.parseNumber()
    default:
        return nil, p.errorf("unexpected character: %q", p.current())
    }
}

// parseObject parses a JSON object: { "key": value, ... }
func (p *parser) parseObject() (*Value, error) {
    p.advance() // consume '{'
    p.skipWhitespace()

    obj := NewObject()

    // Empty object
    if p.pos < len(p.input) && p.current() == '}' {
        p.advance()
        return obj, nil
    }

    for {
        p.skipWhitespace()
        if p.pos >= len(p.input) {
            return nil, p.errorf("unexpected end of input in object")
        }

        // Parse key (must be a string)
        if p.current() != '"' {
            return nil, p.errorf("expected string key in object, got %q", p.current())
        }
        key, err := p.parseString()
        if err != nil {
            return nil, fmt.Errorf("parsing object key: %v", err)
        }

        p.skipWhitespace()

        // Expect ':'
        if p.pos >= len(p.input) || p.current() != ':' {
            return nil, p.errorf("expected ':' after object key %q", key)
        }
        p.advance() // consume ':'
        p.skipWhitespace()

        // Parse value
        val, err := p.parseValue()
        if err != nil {
            return nil, fmt.Errorf("parsing value for key %q: %v", key, err)
        }
        obj.Set(key, val)

        p.skipWhitespace()

        if p.pos >= len(p.input) {
            return nil, p.errorf("unexpected end of input in object (missing '}')")
        }

        switch p.current() {
        case ',':
            p.advance() // consume ','
            // Check for trailing comma (not valid JSON)
            p.skipWhitespace()
            if p.pos < len(p.input) && p.current() == '}' {
                return nil, p.errorf("trailing comma in object (not valid JSON)")
            }
            continue
        case '}':
            p.advance() // consume '}'
            return obj, nil
        default:
            return nil, p.errorf("expected ',' or '}' in object, got %q", p.current())
        }
    }
}

// parseArray parses a JSON array: [ value, value, ... ]
func (p *parser) parseArray() (*Value, error) {
    p.advance() // consume '['
    p.skipWhitespace()

    arr := NewArray()

    // Empty array
    if p.pos < len(p.input) && p.current() == ']' {
        p.advance()
        return arr, nil
    }

    for {
        p.skipWhitespace()
        if p.pos >= len(p.input) {
            return nil, p.errorf("unexpected end of input in array")
        }

        // Check for trailing comma
        if p.current() == ']' {
            return nil, p.errorf("trailing comma in array (not valid JSON)")
        }

        val, err := p.parseValue()
        if err != nil {
            return nil, fmt.Errorf("parsing array element: %v", err)
        }
        arr.arrVal = append(arr.arrVal, val)

        p.skipWhitespace()

        if p.pos >= len(p.input) {
            return nil, p.errorf("unexpected end of input in array (missing ']')")
        }

        switch p.current() {
        case ',':
            p.advance() // consume ','
            continue
        case ']':
            p.advance() // consume ']'
            return arr, nil
        default:
            return nil, p.errorf("expected ',' or ']' in array, got %q", p.current())
        }
    }
}

// ─── String Parser ──────────────────────────────────────────────────────────

// parseString parses a JSON string, handling all escape sequences.
// This is the most complex part of the JSON parser.
func (p *parser) parseString() (string, error) {
    if p.pos >= len(p.input) || p.current() != '"' {
        return "", p.errorf("expected '\"', got %q", p.current())
    }
    p.advance() // consume opening '"'

    var sb strings.Builder

    for p.pos < len(p.input) {
        ch := p.current()

        if ch == '"' {
            p.advance() // consume closing '"'
            return sb.String(), nil
        }

        if ch == '\\' {
            // Escape sequence
            p.advance() // consume '\'
            if p.pos >= len(p.input) {
                return "", p.errorf("unexpected end of string escape sequence")
            }
            escaped := p.current()
            p.advance() // consume escape character

            switch escaped {
            case '"':  sb.WriteByte('"')
            case '\\': sb.WriteByte('\\')
            case '/':  sb.WriteByte('/')
            case 'b':  sb.WriteByte('\b')
            case 'f':  sb.WriteByte('\f')
            case 'n':  sb.WriteByte('\n')
            case 'r':  sb.WriteByte('\r')
            case 't':  sb.WriteByte('\t')
            case 'u':
                // Unicode escape: \uXXXX
                codepoint, err := p.parseUnicodeEscape()
                if err != nil { return "", err }

                // Handle UTF-16 surrogate pairs (\uD800-\uDFFF)
                if codepoint >= 0xD800 && codepoint <= 0xDBFF {
                    // High surrogate — expect a low surrogate next
                    if p.pos+1 < len(p.input) && p.current() == '\\' && p.input[p.pos+1] == 'u' {
                        p.advance() // consume '\'
                        p.advance() // consume 'u'
                        lowSurrogate, err := p.parseUnicodeEscape()
                        if err != nil { return "", err }
                        if lowSurrogate >= 0xDC00 && lowSurrogate <= 0xDFFF {
                            // Combine surrogates into a single code point
                            codepoint = 0x10000 + (codepoint-0xD800)*0x400 + (lowSurrogate - 0xDC00)
                        } else {
                            return "", p.errorf("invalid UTF-16 surrogate pair")
                        }
                    }
                }
                sb.WriteRune(rune(codepoint))
            default:
                return "", p.errorf("invalid escape sequence: '\\%c'", escaped)
            }
        } else if ch < 0x20 {
            // Control characters must be escaped in JSON strings
            return "", p.errorf("unescaped control character in string: %q", ch)
        } else {
            sb.WriteRune(ch)
            p.advance()
        }
    }

    return "", p.errorf("unterminated string (missing closing '\"')")
}

// parseUnicodeEscape parses exactly 4 hex digits and returns the code point value.
func (p *parser) parseUnicodeEscape() (rune, error) {
    if p.pos+4 > len(p.input) {
        return 0, p.errorf("incomplete \\u escape sequence")
    }
    hex := string(p.input[p.pos : p.pos+4])
    codepoint, err := strconv.ParseInt(hex, 16, 32)
    if err != nil {
        return 0, p.errorf("invalid \\u escape: \\u%s", hex)
    }
    p.pos += 4
    p.col += 4
    return rune(codepoint), nil
}

// ─── Number Parser ──────────────────────────────────────────────────────────

// parseNumber parses a JSON number, returning either an IntValue or FloatValue.
func (p *parser) parseNumber() (*Value, error) {
    start := p.pos
    isFloat := false

    // Optional leading minus
    if p.pos < len(p.input) && p.current() == '-' {
        p.advance()
    }

    // Integer part
    if p.pos >= len(p.input) {
        return nil, p.errorf("unexpected end of number")
    }
    if p.current() == '0' {
        p.advance()
        // After a leading 0, only '.', 'e', 'E', or end of number is valid
        // (no "01234" allowed)
    } else if p.current() >= '1' && p.current() <= '9' {
        for p.pos < len(p.input) && p.current() >= '0' && p.current() <= '9' {
            p.advance()
        }
    } else {
        return nil, p.errorf("invalid number: expected digit after '-', got %q", p.current())
    }

    // Fractional part
    if p.pos < len(p.input) && p.current() == '.' {
        isFloat = true
        p.advance() // consume '.'
        if p.pos >= len(p.input) || p.current() < '0' || p.current() > '9' {
            return nil, p.errorf("expected digit after decimal point")
        }
        for p.pos < len(p.input) && p.current() >= '0' && p.current() <= '9' {
            p.advance()
        }
    }

    // Exponent part
    if p.pos < len(p.input) && (p.current() == 'e' || p.current() == 'E') {
        isFloat = true
        p.advance() // consume 'e' or 'E'
        if p.pos < len(p.input) && (p.current() == '+' || p.current() == '-') {
            p.advance()
        }
        if p.pos >= len(p.input) || p.current() < '0' || p.current() > '9' {
            return nil, p.errorf("expected digit in exponent")
        }
        for p.pos < len(p.input) && p.current() >= '0' && p.current() <= '9' {
            p.advance()
        }
    }

    numStr := string(p.input[start:p.pos])

    if isFloat {
        f, err := strconv.ParseFloat(numStr, 64)
        if err != nil {
            return nil, p.errorf("invalid float: %s", numStr)
        }
        return FloatValue(f), nil
    }

    n, err := strconv.ParseInt(numStr, 10, 64)
    if err != nil {
        // May overflow int64 — try float
        f, ferr := strconv.ParseFloat(numStr, 64)
        if ferr != nil {
            return nil, p.errorf("invalid number: %s", numStr)
        }
        return FloatValue(f), nil
    }
    return IntValue(n), nil
}

// ─── Literal Parser ─────────────────────────────────────────────────────────

// parseLiteral parses a known literal string (true, false, null).
func (p *parser) parseLiteral(lit string, result *Value) (*Value, error) {
    for i, ch := range lit {
        if p.pos >= len(p.input) || p.current() != ch {
            return nil, p.errorf("invalid literal (expected %q)", lit)
        }
        _ = i
        p.advance()
    }
    return result, nil
}

// ─── Parser Utilities ───────────────────────────────────────────────────────

func (p *parser) current() rune {
    if p.pos >= len(p.input) { return 0 }
    return p.input[p.pos]
}

func (p *parser) advance() {
    if p.pos < len(p.input) {
        if p.input[p.pos] == '\n' {
            p.line++
            p.col = 1
        } else {
            p.col++
        }
        p.pos++
    }
}

func (p *parser) skipWhitespace() {
    for p.pos < len(p.input) {
        ch := p.current()
        if ch == ' ' || ch == '\t' || ch == '\n' || ch == '\r' {
            p.advance()
        } else {
            break
        }
    }
}

func (p *parser) errorf(format string, args ...interface{}) *ParseError {
    return &ParseError{
        Message: fmt.Sprintf(format, args...),
        Line:    p.line,
        Col:     p.col,
        Pos:     p.pos,
    }
}

// ─── Recursive Descent Diagram ───────────────────────────────────────────────
//
// Input: '{"id":1,"name":"Aditya"}'
//
//  parseValue()
//   └─ sees '{' → parseObject()
//       ├─ advance() → consume '{'
//       ├─ skipWhitespace()
//       ├─ parseString() → "id"
//       ├─ advance() → consume ':'
//       ├─ parseValue()
//       │   └─ sees '1' → parseNumber() → IntValue(1)
//       ├─ sees ',' → advance, continue loop
//       ├─ parseString() → "name"
//       ├─ advance() → consume ':'
//       ├─ parseValue()
//       │   └─ sees '"' → parseString() → StringValue("Aditya")
//       ├─ sees '}' → advance, return Object{id:1, name:"Aditya"}
//       └─ return ObjectValue
```

---

## 4. JSON Marshaling

Marshaling converts Astra values to JSON text. We use Go's `reflect` package to inspect struct fields at runtime:

```go
// stdlib/json/marshal.go

package astra_json

import (
    "fmt"
    "reflect"
    "sort"
    "strings"
    "unicode"
)

// Marshal converts any Go/Astra value to a JSON string.
func Marshal(v interface{}) (string, error) {
    var sb strings.Builder
    if err := marshalValue(reflect.ValueOf(v), &sb); err != nil {
        return "", err
    }
    return sb.String(), nil
}

// MarshalPretty converts to a pretty-printed JSON string with `indent` spaces.
func MarshalPretty(v interface{}, indent int) (string, error) {
    compact, err := Marshal(v)
    if err != nil {
        return "", err
    }
    // Re-parse and re-stringify with indentation
    parsed, err := Parse(compact)
    if err != nil {
        return "", err
    }
    return stringifyPretty(parsed, indent, 0), nil
}

// marshalValue recursively serializes a reflect.Value.
func marshalValue(v reflect.Value, sb *strings.Builder) error {
    // Dereference pointers
    for v.Kind() == reflect.Ptr {
        if v.IsNil() {
            sb.WriteString("null")
            return nil
        }
        v = v.Elem()
    }

    switch v.Kind() {
    case reflect.Invalid:
        sb.WriteString("null")

    case reflect.Bool:
        if v.Bool() {
            sb.WriteString("true")
        } else {
            sb.WriteString("false")
        }

    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        sb.WriteString(fmt.Sprintf("%d", v.Int()))

    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
        sb.WriteString(fmt.Sprintf("%d", v.Uint()))

    case reflect.Float32, reflect.Float64:
        f := v.Float()
        // Format without trailing zeros; use 'g' format
        s := fmt.Sprintf("%g", f)
        // If it looks like an integer, add .0 to distinguish from int
        if !strings.Contains(s, ".") && !strings.Contains(s, "e") {
            s += ".0"
        }
        sb.WriteString(s)

    case reflect.String:
        sb.WriteByte('"')
        sb.WriteString(escapeString(v.String()))
        sb.WriteByte('"')

    case reflect.Slice:
        if v.IsNil() {
            sb.WriteString("null")
            return nil
        }
        // Special case: []byte marshals as base64-encoded string
        if v.Type().Elem().Kind() == reflect.Uint8 {
            sb.WriteByte('"')
            sb.WriteString(escapeString(string(v.Bytes())))
            sb.WriteByte('"')
            return nil
        }
        return marshalArray(v, sb)

    case reflect.Array:
        return marshalArray(v, sb)

    case reflect.Map:
        return marshalMap(v, sb)

    case reflect.Struct:
        return marshalStruct(v, sb)

    case reflect.Interface:
        if v.IsNil() {
            sb.WriteString("null")
            return nil
        }
        return marshalValue(v.Elem(), sb)

    default:
        return fmt.Errorf("json.marshal: unsupported type: %s", v.Type().Name())
    }
    return nil
}

func marshalArray(v reflect.Value, sb *strings.Builder) error {
    sb.WriteByte('[')
    n := v.Len()
    for i := 0; i < n; i++ {
        if i > 0 { sb.WriteByte(',') }
        if err := marshalValue(v.Index(i), sb); err != nil {
            return err
        }
    }
    sb.WriteByte(']')
    return nil
}

func marshalMap(v reflect.Value, sb *strings.Builder) error {
    if v.Type().Key().Kind() != reflect.String {
        return fmt.Errorf("json.marshal: map keys must be strings")
    }
    sb.WriteByte('{')
    // Sort keys for deterministic output
    keys := v.MapKeys()
    sort.Slice(keys, func(i, j int) bool {
        return keys[i].String() < keys[j].String()
    })
    for i, key := range keys {
        if i > 0 { sb.WriteByte(',') }
        sb.WriteByte('"')
        sb.WriteString(escapeString(key.String()))
        sb.WriteString(`":`)
        if err := marshalValue(v.MapIndex(key), sb); err != nil {
            return err
        }
    }
    sb.WriteByte('}')
    return nil
}

func marshalStruct(v reflect.Value, sb *strings.Builder) error {
    t := v.Type()
    sb.WriteByte('{')
    first := true

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)

        // Skip unexported fields
        if !field.IsExported() { continue }

        // Get JSON key name from struct tag, or derive from field name
        jsonKey := getJSONKey(field)
        if jsonKey == "-" { continue } // explicitly excluded

        fieldVal := v.Field(i)

        // Skip zero-value fields if tagged `json:",omitempty"`
        if shouldOmit(field, fieldVal) { continue }

        if !first { sb.WriteByte(',') }
        first = false

        sb.WriteByte('"')
        sb.WriteString(jsonKey)
        sb.WriteString(`":`)
        if err := marshalValue(fieldVal, sb); err != nil {
            return fmt.Errorf("marshaling field %s: %v", field.Name, err)
        }
    }
    sb.WriteByte('}')
    return nil
}

// getJSONKey returns the JSON key for a struct field.
// Priority: json struct tag > snake_case of field name.
func getJSONKey(field reflect.StructField) string {
    tag := field.Tag.Get("json")
    if tag != "" {
        parts := strings.Split(tag, ",")
        if parts[0] != "" {
            return parts[0]
        }
    }
    // Convert CamelCase to snake_case
    return toSnakeCase(field.Name)
}

// toSnakeCase converts "CamelCase" to "camel_case".
func toSnakeCase(s string) string {
    var result strings.Builder
    for i, r := range s {
        if i > 0 && unicode.IsUpper(r) {
            result.WriteByte('_')
        }
        result.WriteRune(unicode.ToLower(r))
    }
    return result.String()
}

// shouldOmit returns true if a field with omitempty should be skipped.
func shouldOmit(field reflect.StructField, v reflect.Value) bool {
    tag := field.Tag.Get("json")
    if !strings.Contains(tag, "omitempty") { return false }
    return v.IsZero()
}

// escapeString escapes a string for safe inclusion in a JSON string.
func escapeString(s string) string {
    var sb strings.Builder
    for _, r := range s {
        switch r {
        case '"':  sb.WriteString(`\"`)
        case '\\': sb.WriteString(`\\`)
        case '\b': sb.WriteString(`\b`)
        case '\f': sb.WriteString(`\f`)
        case '\n': sb.WriteString(`\n`)
        case '\r': sb.WriteString(`\r`)
        case '\t': sb.WriteString(`\t`)
        default:
            if r < 0x20 {
                // Control characters must be escaped
                sb.WriteString(fmt.Sprintf(`\u%04x`, r))
            } else {
                sb.WriteRune(r)
            }
        }
    }
    return sb.String()
}
```

---

## 5. JSON Unmarshaling

```go
// stdlib/json/unmarshal.go

package astra_json

import (
    "fmt"
    "reflect"
    "strings"
)

// Unmarshal parses JSON text and populates the value pointed to by dst.
// dst must be a pointer to a struct, slice, or basic type.
func Unmarshal(jsonStr string, dst interface{}) error {
    parsed, err := Parse(jsonStr)
    if err != nil {
        return fmt.Errorf("json.unmarshal: parse error: %v", err)
    }
    return populateValue(parsed, reflect.ValueOf(dst))
}

// populateValue fills a Go value from a parsed JSON Value.
func populateValue(src *Value, dst reflect.Value) error {
    // Dereference pointers, allocating new values as needed
    for dst.Kind() == reflect.Ptr {
        if dst.IsNil() {
            dst.Set(reflect.New(dst.Type().Elem()))
        }
        dst = dst.Elem()
    }

    if src == nil || src.kind == TypeNull {
        // Set zero value for null
        dst.Set(reflect.Zero(dst.Type()))
        return nil
    }

    switch dst.Kind() {
    case reflect.Bool:
        if !src.IsBool() {
            return fmt.Errorf("expected bool, got %v", src.Type())
        }
        dst.SetBool(src.AsBool())

    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
        switch src.kind {
        case TypeInt:   dst.SetInt(src.intVal)
        case TypeFloat: dst.SetInt(int64(src.floatVal))
        default:
            return fmt.Errorf("expected number for int field, got %v", src.Type())
        }

    case reflect.Float32, reflect.Float64:
        switch src.kind {
        case TypeFloat: dst.SetFloat(src.floatVal)
        case TypeInt:   dst.SetFloat(float64(src.intVal))
        default:
            return fmt.Errorf("expected number for float field, got %v", src.Type())
        }

    case reflect.String:
        if !src.IsString() {
            return fmt.Errorf("expected string, got %v", src.Type())
        }
        dst.SetString(src.strVal)

    case reflect.Slice:
        if !src.IsArray() {
            return fmt.Errorf("expected array, got %v", src.Type())
        }
        slice := reflect.MakeSlice(dst.Type(), len(src.arrVal), len(src.arrVal))
        for i, elem := range src.arrVal {
            if err := populateValue(elem, slice.Index(i)); err != nil {
                return fmt.Errorf("array[%d]: %v", i, err)
            }
        }
        dst.Set(slice)

    case reflect.Map:
        if !src.IsObject() {
            return fmt.Errorf("expected object for map, got %v", src.Type())
        }
        if dst.IsNil() {
            dst.Set(reflect.MakeMap(dst.Type()))
        }
        for key, val := range src.objVal {
            keyVal := reflect.ValueOf(key)
            elemVal := reflect.New(dst.Type().Elem()).Elem()
            if err := populateValue(val, elemVal); err != nil {
                return fmt.Errorf("map key %q: %v", key, err)
            }
            dst.SetMapIndex(keyVal, elemVal)
        }

    case reflect.Struct:
        if !src.IsObject() {
            return fmt.Errorf("expected object for struct, got %v", src.Type())
        }
        return populateStruct(src, dst)

    case reflect.Interface:
        // For interface{} fields, use the native Go representation
        dst.Set(reflect.ValueOf(jsonValueToInterface(src)))
    }

    return nil
}

// populateStruct fills a Go struct from a JSON object.
func populateStruct(src *Value, dst reflect.Value) error {
    t := dst.Type()

    for i := 0; i < t.NumField(); i++ {
        field := t.Field(i)
        if !field.IsExported() { continue }

        jsonKey := getJSONKey(field)
        if jsonKey == "-" { continue }

        jsonVal := src.Get(jsonKey)
        if jsonVal == nil {
            // Try the original field name in case the key isn't snake_cased
            jsonVal = src.Get(field.Name)
        }
        if jsonVal == nil {
            // Field not present in JSON — leave at zero value
            continue
        }

        if err := populateValue(jsonVal, dst.Field(i)); err != nil {
            return fmt.Errorf("field %s (%s): %v", field.Name, jsonKey, err)
        }
    }
    return nil
}

// jsonValueToInterface converts a JSON Value to a plain Go interface{} value.
// Useful for map[string]interface{} deserialization.
func jsonValueToInterface(v *Value) interface{} {
    if v == nil { return nil }
    switch v.kind {
    case TypeNull:   return nil
    case TypeBool:   return v.boolVal
    case TypeInt:    return v.intVal
    case TypeFloat:  return v.floatVal
    case TypeString: return v.strVal
    case TypeArray:
        result := make([]interface{}, len(v.arrVal))
        for i, elem := range v.arrVal {
            result[i] = jsonValueToInterface(elem)
        }
        return result
    case TypeObject:
        result := make(map[string]interface{})
        for k, val := range v.objVal {
            result[k] = jsonValueToInterface(val)
        }
        return result
    }
    return nil
}
```

---

## 6. Astra-Facing API

```go
// stdlib/json/json.go

package astra_json

import "strings"

// ─── High-Level API (called from Astra programs) ─────────────────────────────

// ParseJSON parses a JSON string and returns a JsonValue.
func ParseJSON(input string) (*Value, error) {
    return Parse(input)
}

// ParseArray parses a JSON array string and returns a slice of Values.
func ParseArray(input string) ([]*Value, error) {
    v, err := Parse(input)
    if err != nil { return nil, err }
    if !v.IsArray() {
        return nil, fmt.Errorf("json.parse_array: expected array, got %v", v.Type())
    }
    return v.AsArray(), nil
}

// Stringify converts a Value back to a compact JSON string.
func (v *Value) Stringify() string {
    return stringify(v)
}

// stringify recursively serializes a Value.
func stringify(v *Value) string {
    if v == nil { return "null" }
    switch v.kind {
    case TypeNull:   return "null"
    case TypeBool:
        if v.boolVal { return "true" }
        return "false"
    case TypeInt:    return fmt.Sprintf("%d", v.intVal)
    case TypeFloat:  return fmt.Sprintf("%g", v.floatVal)
    case TypeString:
        return `"` + escapeString(v.strVal) + `"`
    case TypeArray:
        parts := make([]string, len(v.arrVal))
        for i, elem := range v.arrVal {
            parts[i] = stringify(elem)
        }
        return "[" + strings.Join(parts, ",") + "]"
    case TypeObject:
        parts := make([]string, 0, len(v.objKeys))
        for _, key := range v.objKeys {
            val := v.objVal[key]
            parts = append(parts, `"`+escapeString(key)+`":`+stringify(val))
        }
        return "{" + strings.Join(parts, ",") + "}"
    }
    return "null"
}

// stringifyPretty produces indented JSON output.
func stringifyPretty(v *Value, indent, depth int) string {
    if v == nil { return "null" }
    prefix := strings.Repeat(" ", indent*depth)
    childPrefix := strings.Repeat(" ", indent*(depth+1))

    switch v.kind {
    case TypeArray:
        if len(v.arrVal) == 0 { return "[]" }
        var sb strings.Builder
        sb.WriteString("[\n")
        for i, elem := range v.arrVal {
            sb.WriteString(childPrefix)
            sb.WriteString(stringifyPretty(elem, indent, depth+1))
            if i < len(v.arrVal)-1 { sb.WriteByte(',') }
            sb.WriteByte('\n')
        }
        sb.WriteString(prefix + "]")
        return sb.String()

    case TypeObject:
        if len(v.objKeys) == 0 { return "{}" }
        var sb strings.Builder
        sb.WriteString("{\n")
        for i, key := range v.objKeys {
            sb.WriteString(childPrefix)
            sb.WriteString(`"` + escapeString(key) + `": `)
            sb.WriteString(stringifyPretty(v.objVal[key], indent, depth+1))
            if i < len(v.objKeys)-1 { sb.WriteByte(',') }
            sb.WriteByte('\n')
        }
        sb.WriteString(prefix + "}")
        return sb.String()

    default:
        return stringify(v)
    }
}
```

Astra usage:

```astra
import json
import file

struct ServerConfig {
    host:          string
    port:          int
    max_conns:     int
    debug:         bool
    allowed_origins: List<string>
}

struct AppConfig {
    server:   ServerConfig
    database: string
    log_file: string
}

fn main() {
    // ─── Parse from string ───────────────────────────────────────────────
    let text = '{"id":1,"name":"Aditya","email":"a@b.com","is_admin":true}'
    let data = json.parse(text)

    let id    = data.get_int("id")       // → 1
    let name  = data.get_string("name")  // → "Aditya"
    let admin = data.get_bool("is_admin")// → true
    print("User: " + name + " (admin: " + admin.to_string() + ")")

    // ─── Nested access via path ──────────────────────────────────────────
    let config_text = '{"server":{"host":"localhost","port":8080},"debug":false}'
    let cfg = json.parse(config_text)
    let host = cfg.path("server.host").as_string()   // → "localhost"
    let port = cfg.path("server.port").as_int()      // → 8080
    print("Server: " + host + ":" + port.to_string())

    // ─── Array parsing ───────────────────────────────────────────────────
    let arr_text = '[1, 2, 3, 4, 5]'
    let arr = json.parse_array(arr_text)
    let sum = 0
    for item in arr {
        sum = sum + item.as_int()
    }
    print("Sum: " + sum.to_string())  // → Sum: 15

    // ─── Direct struct unmarshaling ──────────────────────────────────────
    let full_config_text = file.read("config.json")
    let app_cfg = json.unmarshal<AppConfig>(full_config_text)
    print("Host: " + app_cfg.server.host)
    print("Port: " + app_cfg.server.port.to_string())
    print("DB:   " + app_cfg.database)

    // ─── Building JSON objects programmatically ──────────────────────────
    let obj = json.Object.new()
    obj.set("message", "Hello from Astra!")
    obj.set("version", 1)
    obj.set("languages", json.Array.from(["astra", "go", "c"]))
    obj.set("metadata", json.Object.from({
        "built_with": "astrac",
        "stable": true
    }))
    print(obj.to_string())
    // → {"message":"Hello from Astra!","version":1,"languages":["astra","go","c"],...}

    // ─── Pretty printing ─────────────────────────────────────────────────
    let pretty = json.pretty(obj, 2)
    print(pretty)
    // {
    //   "message": "Hello from Astra!",
    //   "version": 1,
    //   "languages": [
    //     "astra",
    //     "go",
    //     "c"
    //   ],
    //   ...
    // }

    // ─── Marshaling Astra struct to JSON ─────────────────────────────────
    struct User {
        id:       int
        name:     string
        email:    string
        is_admin: bool
    }
    let user = User { id: 42, name: "Bob", email: "bob@example.com", is_admin: false }
    let json_str = json.marshal(user)
    print(json_str)
    // → {"id":42,"name":"Bob","email":"bob@example.com","is_admin":false}
}
```

---

## 7. Error Handling for JSON

```astra
// Malformed JSON
match json.parse('{"key": "value"') {  // missing closing }
    Ok(v)  => print(v.to_string())
    Err(e) => print("Parse error: " + e)
    // → "Parse error: json parse error at line 1, col 17: unexpected end of input in object"
}

// Type mismatch in unmarshal
struct Point { x: int; y: int }
match json.unmarshal<Point>('{"x": "not a number", "y": 5}') {
    Ok(p)  => print(p.x.to_string())
    Err(e) => print("Type error: " + e)
    // → "Type error: field x (x): expected number for int field, got TypeString"
}

// Accessing a missing key returns the zero value (no error)
let data = json.parse('{"name": "Aditya"}')
let age = data.get_int("age")  // → 0 (zero value for int)
let email = data.get_string("email")  // → "" (zero value for string)
// To distinguish missing from zero, use has():
if data.has("age") {
    print("Age: " + data.get_int("age").to_string())
} else {
    print("Age not provided")
}
```

---

## 8. Streaming JSON Parser

For large files (hundreds of megabytes), loading the entire JSON into memory as a DOM tree is impractical. A streaming parser emits events as it encounters JSON tokens:

```go
// stdlib/json/stream.go

package astra_json

// EventType enumerates the events a streaming parser can emit.
type EventType int

const (
    EventObjectStart EventType = iota
    EventObjectEnd
    EventArrayStart
    EventArrayEnd
    EventKey
    EventString
    EventNumber
    EventBool
    EventNull
    EventError
)

// Event is emitted by the streaming parser.
type Event struct {
    Type     EventType
    Key      string  // for EventKey
    StrVal   string  // for EventString
    IntVal   int64   // for EventNumber (integer)
    FloatVal float64 // for EventNumber (float)
    IsFloat  bool    // true if number was a float
    BoolVal  bool    // for EventBool
    Err      error   // for EventError
}

// StreamParser is a SAX-style JSON parser that emits events for each token.
// Use it for streaming large JSON arrays without loading all items into memory.
type StreamParser struct {
    p      *parser
    events chan Event
    done   chan struct{}
}

// NewStreamParser creates a streaming parser for the given JSON string.
// Call Next() repeatedly to receive events.
func NewStreamParser(input string) *StreamParser {
    sp := &StreamParser{
        p:      &parser{input: []rune(input), line: 1, col: 1},
        events: make(chan Event, 64),
        done:   make(chan struct{}),
    }
    go sp.run()
    return sp
}

// Next returns the next event from the parser.
// Returns an EventError event when parsing is complete or fails.
func (sp *StreamParser) Next() Event {
    select {
    case ev := <-sp.events:
        return ev
    case <-sp.done:
        return Event{Type: EventError, Err: fmt.Errorf("stream ended")}
    }
}

func (sp *StreamParser) run() {
    defer close(sp.done)
    sp.p.skipWhitespace()
    sp.streamValue()
}

func (sp *StreamParser) emit(ev Event) {
    sp.events <- ev
}

func (sp *StreamParser) streamValue() {
    if sp.p.pos >= len(sp.p.input) {
        sp.emit(Event{Type: EventError, Err: fmt.Errorf("unexpected end of input")})
        return
    }

    switch sp.p.current() {
    case '{':
        sp.p.advance()
        sp.emit(Event{Type: EventObjectStart})
        sp.p.skipWhitespace()
        if sp.p.pos < len(sp.p.input) && sp.p.current() == '}' {
            sp.p.advance()
            sp.emit(Event{Type: EventObjectEnd})
            return
        }
        for {
            sp.p.skipWhitespace()
            key, err := sp.p.parseString()
            if err != nil {
                sp.emit(Event{Type: EventError, Err: err})
                return
            }
            sp.emit(Event{Type: EventKey, Key: key})
            sp.p.skipWhitespace()
            sp.p.advance() // consume ':'
            sp.p.skipWhitespace()
            sp.streamValue()
            sp.p.skipWhitespace()
            if sp.p.current() == '}' {
                sp.p.advance()
                sp.emit(Event{Type: EventObjectEnd})
                return
            }
            sp.p.advance() // consume ','
        }

    case '[':
        sp.p.advance()
        sp.emit(Event{Type: EventArrayStart})
        sp.p.skipWhitespace()
        if sp.p.pos < len(sp.p.input) && sp.p.current() == ']' {
            sp.p.advance()
            sp.emit(Event{Type: EventArrayEnd})
            return
        }
        for {
            sp.streamValue()
            sp.p.skipWhitespace()
            if sp.p.current() == ']' {
                sp.p.advance()
                sp.emit(Event{Type: EventArrayEnd})
                return
            }
            sp.p.advance() // consume ','
            sp.p.skipWhitespace()
        }

    case '"':
        s, err := sp.p.parseString()
        if err != nil {
            sp.emit(Event{Type: EventError, Err: err})
            return
        }
        sp.emit(Event{Type: EventString, StrVal: s})

    default:
        v, err := sp.p.parseValue()
        if err != nil {
            sp.emit(Event{Type: EventError, Err: err})
            return
        }
        switch v.kind {
        case TypeBool:  sp.emit(Event{Type: EventBool, BoolVal: v.boolVal})
        case TypeNull:  sp.emit(Event{Type: EventNull})
        case TypeInt:   sp.emit(Event{Type: EventNumber, IntVal: v.intVal})
        case TypeFloat: sp.emit(Event{Type: EventNumber, FloatVal: v.floatVal, IsFloat: true})
        }
    }
}
```

Using the streaming parser in Astra to process a large array of users without loading all of them into memory:

```astra
import json
import file

fn process_large_users_file(path: string) {
    let content = file.read(path)  // still loads file into memory
    // For true streaming, use file.LineReader + json.parse per line (NDJSON)
    let stream = json.StreamParser.new(content)

    let depth = 0
    let current_key = ""
    let user_count = 0

    loop {
        let event = stream.next()
        match event.type {
            EventArrayStart  => depth = depth + 1
            EventArrayEnd    => depth = depth - 1
            EventObjectStart => {
                depth = depth + 1
                user_count = user_count + 1
            }
            EventObjectEnd   => depth = depth - 1
            EventKey         => current_key = event.key
            EventString      => {
                if current_key == "name" {
                    print("Processing user: " + event.str_val)
                }
            }
            EventError       => break
            _                => {}
        }
    }
    print("Processed " + user_count.to_string() + " users")
}
```

---

## 🔨 Astra Build Milestone

Complete `stdlib/json/` with full parser, marshaler, and unmarshaler.

**Expected file structure and sizes:**

| File | Lines | Description |
|------|-------|-------------|
| `stdlib/json/parser.go` | ~250 | Recursive descent parser |
| `stdlib/json/value.go` | ~150 | JsonValue type and accessors |
| `stdlib/json/marshal.go` | ~150 | Astra → JSON text |
| `stdlib/json/unmarshal.go` | ~150 | JSON text → Astra struct |
| `stdlib/json/stream.go` | ~120 | Streaming SAX-style parser |
| `stdlib/json/json.go` | ~60 | High-level API entry points |

**Verification test:**

```bash
$ astrac build tests/json_test.as -o json_test
$ ./json_test
json.parse object: PASSED
json.parse array: PASSED
json.parse nested: PASSED
json.parse number types: PASSED
json.parse unicode escapes: PASSED
json.parse error handling: PASSED
json.marshal struct: PASSED
json.marshal array: PASSED
json.marshal pretty: PASSED
json.unmarshal struct: PASSED
json.unmarshal nested: PASSED
json.path: PASSED
All JSON tests passed!
```

---

## 9. Exercises

**Exercise 1 — JSON Schema Validation**
Implement a simple JSON schema validator. A schema is itself a JSON object describing the expected structure:
```json
{ "type": "object", "required": ["id", "name"], "properties": {
    "id":   { "type": "number" },
    "name": { "type": "string", "minLength": 1 }
}}
```
Implement `json.validate(value, schema) -> Result<void, string>` that returns Ok if the value matches the schema.

**Exercise 2 — JSON Diff**
Implement `json.diff(a, b) -> List<string>` that returns a list of human-readable differences between two JSON values:
```
"id: 1 → 2"
"name: added (was absent)"
"address: removed"
"tags[2]: 'go' → 'rust'"
```

**Exercise 3 — JSON Path Queries**
Extend `value.path()` to support array indexing: `config.path("users[0].name")` should return the name of the first user. Also support wildcard: `config.path("users[*].name")` returns a list of all user names.

**Exercise 4 — NDJSON (Newline-Delimited JSON)**
NDJSON stores one JSON object per line. Implement `json.read_ndjson(path: string) -> List<JsonValue>` and `json.write_ndjson(path: string, values: List<JsonValue>)`. This format is ideal for log files and streaming data pipelines.

**Exercise 5 — JSON Merge**
Implement `json.merge(base, override) -> JsonValue` that deep-merges two JSON objects. For conflicting keys, `override` wins. Arrays are replaced entirely (not merged).

**Exercise 6 — Custom Marshaling**
Allow Astra structs to define a `to_json(self) -> string` method that is used instead of the default reflection-based marshaling. Modify the marshaler to check for this method and call it if present.

---

## Summary

| Topic | Key Takeaway |
|-------|-------------|
| JSON grammar | LL(1): one lookahead character determines which production to use |
| Recursive descent | Same technique as Astra's own parser — JSON is just a simpler language |
| parseValue | Dispatches on first character: `{` → object, `[` → array, `"` → string, digit/-→ number |
| String escaping | 8 standard escapes + `\uXXXX`. Must handle UTF-16 surrogate pairs for `\uD800`-`\uDFFF` |
| Number parsing | Two types: integer (int64) and float (float64). No leading zeros. No NaN/Infinity |
| Marshaling | Uses Go reflection to inspect struct fields. CamelCase → snake_case by default |
| Unmarshaling | Populates struct fields by matching JSON keys to snake_case field names |
| Streaming parser | SAX-style event emission for large files. O(1) memory vs O(n) for DOM parser |
| Error messages | Include line and column numbers for debugging. "at line 3, col 14" |
| The `path()` method | Dot-separated key traversal for nested JSON: `config.path("server.port")` |
