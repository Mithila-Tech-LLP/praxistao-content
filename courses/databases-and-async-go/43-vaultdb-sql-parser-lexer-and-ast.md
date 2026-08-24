# Chapter 43: VaultDB — SQL Parser: Lexer and AST

SQL is text. Databases work with structured data. The parser bridges these two worlds: it turns a string like `"SELECT name FROM users WHERE age > 18"` into a tree of Go structs that the executor can work with.

## Table of Contents

1. The Two Steps: Lexing and Parsing
2. The Lexer (Tokenizer)
3. The AST (Abstract Syntax Tree)
4. The Parser
5. Testing the Parser
6. Exercises

---

## 1. The Two Steps: Lexing and Parsing

**Step 1 — Lexing:** Split the raw SQL string into tokens (words, symbols, numbers).

```
"SELECT name FROM users WHERE age > 18"
→ [SELECT] [name] [FROM] [users] [WHERE] [age] [>] [18]
```

**Step 2 — Parsing:** Arrange tokens into a tree that represents the query's meaning.

```
SELECT Statement
├── Columns: ["name"]
├── From: "users"
└── Where:
    └── Greater Than
        ├── Column: "age"
        └── Integer: 18
```

This tree is called an **AST (Abstract Syntax Tree)**. The executor walks this tree to actually run the query.

---

## 2. The Lexer (Tokenizer)

```go
// query/lexer.go
package query

import (
    "fmt"
    "strings"
    "unicode"
)

type TokenType int

const (
    // Keywords
    TokenSelect TokenType = iota + 1
    TokenInsert
    TokenUpdate
    TokenDelete
    TokenCreate
    TokenDrop
    TokenTable
    TokenFrom
    TokenWhere
    TokenSet
    TokenInto
    TokenValues
    TokenAnd
    TokenOr
    TokenNot
    TokenNull
    TokenTrue
    TokenFalse
    TokenOrder
    TokenBy
    TokenLimit
    TokenOffset

    // Literals
    TokenIdent   // column or table name
    TokenInt     // 42
    TokenFloat   // 3.14
    TokenString  // 'hello'

    // Operators
    TokenEq   // =
    TokenNeq  // !=  or  <>
    TokenLt   // <
    TokenLte  // <=
    TokenGt   // >
    TokenGte  // >=

    // Punctuation
    TokenComma    // ,
    TokenStar     // *
    TokenLParen   // (
    TokenRParen   // )
    TokenSemicolon // ;

    // Special
    TokenEOF
)

var keywords = map[string]TokenType{
    "select": TokenSelect,
    "insert": TokenInsert,
    "update": TokenUpdate,
    "delete": TokenDelete,
    "create": TokenCreate,
    "drop":   TokenDrop,
    "table":  TokenTable,
    "from":   TokenFrom,
    "where":  TokenWhere,
    "set":    TokenSet,
    "into":   TokenInto,
    "values": TokenValues,
    "and":    TokenAnd,
    "or":     TokenOr,
    "not":    TokenNot,
    "null":   TokenNull,
    "true":   TokenTrue,
    "false":  TokenFalse,
    "order":  TokenOrder,
    "by":     TokenBy,
    "limit":  TokenLimit,
    "offset": TokenOffset,
}

type Token struct {
    Type    TokenType
    Literal string
    Pos     int
}

type Lexer struct {
    input []rune
    pos   int
}

func NewLexer(input string) *Lexer {
    return &Lexer{input: []rune(input)}
}

func (l *Lexer) NextToken() Token {
    l.skipWhitespace()

    if l.pos >= len(l.input) {
        return Token{Type: TokenEOF, Pos: l.pos}
    }

    ch := l.input[l.pos]

    switch {
    case ch == '\'':
        return l.readString()
    case unicode.IsDigit(ch) || (ch == '-' && l.pos+1 < len(l.input) && unicode.IsDigit(l.input[l.pos+1])):
        return l.readNumber()
    case unicode.IsLetter(ch) || ch == '_':
        return l.readIdent()
    case ch == '=':
        l.pos++
        return Token{Type: TokenEq, Literal: "="}
    case ch == '!':
        if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
            l.pos += 2
            return Token{Type: TokenNeq, Literal: "!="}
        }
    case ch == '<':
        if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
            l.pos += 2
            return Token{Type: TokenLte, Literal: "<="}
        }
        l.pos++
        return Token{Type: TokenLt, Literal: "<"}
    case ch == '>':
        if l.pos+1 < len(l.input) && l.input[l.pos+1] == '=' {
            l.pos += 2
            return Token{Type: TokenGte, Literal: ">="}
        }
        l.pos++
        return Token{Type: TokenGt, Literal: ">"}
    case ch == ',':
        l.pos++
        return Token{Type: TokenComma, Literal: ","}
    case ch == '*':
        l.pos++
        return Token{Type: TokenStar, Literal: "*"}
    case ch == '(':
        l.pos++
        return Token{Type: TokenLParen, Literal: "("}
    case ch == ')':
        l.pos++
        return Token{Type: TokenRParen, Literal: ")"}
    case ch == ';':
        l.pos++
        return Token{Type: TokenSemicolon, Literal: ";"}
    }

    l.pos++
    return Token{Type: TokenEOF, Literal: fmt.Sprintf("unexpected char: %c", ch)}
}

func (l *Lexer) readIdent() Token {
    start := l.pos
    for l.pos < len(l.input) && (unicode.IsLetter(l.input[l.pos]) || unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '_') {
        l.pos++
    }
    word := string(l.input[start:l.pos])
    lower := strings.ToLower(word)
    if tt, ok := keywords[lower]; ok {
        return Token{Type: tt, Literal: word, Pos: start}
    }
    return Token{Type: TokenIdent, Literal: word, Pos: start}
}

func (l *Lexer) readNumber() Token {
    start := l.pos
    isFloat := false
    if l.input[l.pos] == '-' {
        l.pos++
    }
    for l.pos < len(l.input) && (unicode.IsDigit(l.input[l.pos]) || l.input[l.pos] == '.') {
        if l.input[l.pos] == '.' {
            isFloat = true
        }
        l.pos++
    }
    lit := string(l.input[start:l.pos])
    if isFloat {
        return Token{Type: TokenFloat, Literal: lit, Pos: start}
    }
    return Token{Type: TokenInt, Literal: lit, Pos: start}
}

func (l *Lexer) readString() Token {
    l.pos++ // skip opening '
    start := l.pos
    for l.pos < len(l.input) && l.input[l.pos] != '\'' {
        l.pos++
    }
    lit := string(l.input[start:l.pos])
    if l.pos < len(l.input) {
        l.pos++ // skip closing '
    }
    return Token{Type: TokenString, Literal: lit, Pos: start}
}

func (l *Lexer) skipWhitespace() {
    for l.pos < len(l.input) && unicode.IsSpace(l.input[l.pos]) {
        l.pos++
    }
}

// Tokenize returns all tokens
func Tokenize(sql string) []Token {
    l := NewLexer(sql)
    var tokens []Token
    for {
        t := l.NextToken()
        tokens = append(tokens, t)
        if t.Type == TokenEOF {
            break
        }
    }
    return tokens
}
```

---

## 3. The AST (Abstract Syntax Tree)

```go
// query/ast.go
package query

import "github.com/yourname/vaultdb/storage"

// Statement is the top-level node
type Statement interface {
    statementNode()
}

// SELECT name, age FROM users WHERE age > 18 ORDER BY name LIMIT 10
type SelectStmt struct {
    Columns   []string  // ["name", "age"] or ["*"]
    Table     string    // "users"
    Where     Expr      // nil if no WHERE
    OrderBy   string    // "" if no ORDER BY
    OrderDesc bool
    Limit     int       // 0 = no limit
}

func (s *SelectStmt) statementNode() {}

// INSERT INTO users (name, age) VALUES ('Alice', 25)
type InsertStmt struct {
    Table   string
    Columns []string
    Values  []Expr
}

func (s *InsertStmt) statementNode() {}

// UPDATE users SET age = 26 WHERE name = 'Alice'
type UpdateStmt struct {
    Table  string
    Sets   []Assignment
    Where  Expr
}

type Assignment struct {
    Column string
    Value  Expr
}

func (s *UpdateStmt) statementNode() {}

// DELETE FROM users WHERE age < 18
type DeleteStmt struct {
    Table string
    Where Expr
}

func (s *DeleteStmt) statementNode() {}

// CREATE TABLE users (id INT, name VARCHAR, age INT)
type CreateTableStmt struct {
    Table   string
    Columns []ColDef
}

type ColDef struct {
    Name   string
    TypeID storage.TypeID
}

func (s *CreateTableStmt) statementNode() {}

// Expr is any expression: literals, column references, binary ops
type Expr interface {
    exprNode()
}

type IntLiteral   struct{ Value int64 }
type FloatLiteral struct{ Value float64 }
type StrLiteral   struct{ Value string }
type BoolLiteral  struct{ Value bool }
type NullLiteral  struct{}
type ColumnRef    struct{ Name string }

func (e *IntLiteral)   exprNode() {}
func (e *FloatLiteral) exprNode() {}
func (e *StrLiteral)   exprNode() {}
func (e *BoolLiteral)  exprNode() {}
func (e *NullLiteral)  exprNode() {}
func (e *ColumnRef)    exprNode() {}

type BinaryOp struct {
    Op    string // "=", "!=", "<", "<=", ">", ">=", "AND", "OR"
    Left  Expr
    Right Expr
}

func (e *BinaryOp) exprNode() {}
```

---

## 4. The Parser

```go
// query/parser.go
package query

import (
    "fmt"
    "strconv"
    "strings"

    "github.com/yourname/vaultdb/storage"
)

type Parser struct {
    tokens []Token
    pos    int
}

func NewParser(sql string) *Parser {
    return &Parser{tokens: Tokenize(sql)}
}

func Parse(sql string) (Statement, error) {
    p := NewParser(sql)
    return p.parseStatement()
}

func (p *Parser) parseStatement() (Statement, error) {
    tok := p.peek()
    switch tok.Type {
    case TokenSelect:
        return p.parseSelect()
    case TokenInsert:
        return p.parseInsert()
    case TokenUpdate:
        return p.parseUpdate()
    case TokenDelete:
        return p.parseDelete()
    case TokenCreate:
        return p.parseCreateTable()
    }
    return nil, fmt.Errorf("unexpected token: %s", tok.Literal)
}

func (p *Parser) parseSelect() (*SelectStmt, error) {
    p.consume(TokenSelect)
    stmt := &SelectStmt{}

    // Parse column list or *
    tok := p.peek()
    if tok.Type == TokenStar {
        p.advance()
        stmt.Columns = []string{"*"}
    } else {
        for {
            col, err := p.expectIdent()
            if err != nil {
                return nil, err
            }
            stmt.Columns = append(stmt.Columns, col)
            if p.peek().Type != TokenComma {
                break
            }
            p.advance()
        }
    }

    // FROM table
    if err := p.consume(TokenFrom); err != nil {
        return nil, err
    }
    table, err := p.expectIdent()
    if err != nil {
        return nil, err
    }
    stmt.Table = table

    // Optional WHERE
    if p.peek().Type == TokenWhere {
        p.advance()
        where, err := p.parseExpr()
        if err != nil {
            return nil, err
        }
        stmt.Where = where
    }

    // Optional ORDER BY
    if p.peek().Type == TokenOrder {
        p.advance()
        p.consume(TokenBy)
        col, err := p.expectIdent()
        if err != nil {
            return nil, err
        }
        stmt.OrderBy = col
    }

    // Optional LIMIT
    if p.peek().Type == TokenLimit {
        p.advance()
        lit, err := p.expectInt()
        if err != nil {
            return nil, err
        }
        stmt.Limit = int(lit)
    }

    return stmt, nil
}

func (p *Parser) parseInsert() (*InsertStmt, error) {
    p.consume(TokenInsert)
    p.consume(TokenInto)
    stmt := &InsertStmt{}

    table, err := p.expectIdent()
    if err != nil {
        return nil, err
    }
    stmt.Table = table

    // Optional column list
    if p.peek().Type == TokenLParen {
        p.advance()
        for {
            col, err := p.expectIdent()
            if err != nil {
                return nil, err
            }
            stmt.Columns = append(stmt.Columns, col)
            if p.peek().Type == TokenRParen {
                break
            }
            p.consume(TokenComma)
        }
        p.consume(TokenRParen)
    }

    p.consume(TokenValues)
    p.consume(TokenLParen)
    for {
        val, err := p.parsePrimary()
        if err != nil {
            return nil, err
        }
        stmt.Values = append(stmt.Values, val)
        if p.peek().Type == TokenRParen {
            break
        }
        p.consume(TokenComma)
    }
    p.consume(TokenRParen)

    return stmt, nil
}

func (p *Parser) parseUpdate() (*UpdateStmt, error) {
    p.consume(TokenUpdate)
    stmt := &UpdateStmt{}
    table, err := p.expectIdent()
    if err != nil {
        return nil, err
    }
    stmt.Table = table
    p.consume(TokenSet)

    for {
        col, err := p.expectIdent()
        if err != nil {
            return nil, err
        }
        p.consume(TokenEq)
        val, err := p.parsePrimary()
        if err != nil {
            return nil, err
        }
        stmt.Sets = append(stmt.Sets, Assignment{Column: col, Value: val})
        if p.peek().Type != TokenComma {
            break
        }
        p.advance()
    }

    if p.peek().Type == TokenWhere {
        p.advance()
        where, err := p.parseExpr()
        if err != nil {
            return nil, err
        }
        stmt.Where = where
    }
    return stmt, nil
}

func (p *Parser) parseDelete() (*DeleteStmt, error) {
    p.consume(TokenDelete)
    p.consume(TokenFrom)
    stmt := &DeleteStmt{}
    table, err := p.expectIdent()
    if err != nil {
        return nil, err
    }
    stmt.Table = table
    if p.peek().Type == TokenWhere {
        p.advance()
        where, err := p.parseExpr()
        if err != nil {
            return nil, err
        }
        stmt.Where = where
    }
    return stmt, nil
}

func (p *Parser) parseCreateTable() (*CreateTableStmt, error) {
    p.consume(TokenCreate)
    p.consume(TokenTable)
    stmt := &CreateTableStmt{}
    table, err := p.expectIdent()
    if err != nil {
        return nil, err
    }
    stmt.Table = table
    p.consume(TokenLParen)
    for {
        col, err := p.expectIdent()
        if err != nil {
            return nil, err
        }
        typeStr, err := p.expectIdent()
        if err != nil {
            return nil, err
        }
        stmt.Columns = append(stmt.Columns, ColDef{
            Name:   col,
            TypeID: parseType(typeStr),
        })
        if p.peek().Type == TokenRParen {
            break
        }
        p.consume(TokenComma)
    }
    p.consume(TokenRParen)
    return stmt, nil
}

// parseExpr handles AND/OR expressions
func (p *Parser) parseExpr() (Expr, error) {
    left, err := p.parseComparison()
    if err != nil {
        return nil, err
    }
    for p.peek().Type == TokenAnd || p.peek().Type == TokenOr {
        op := p.advance().Literal
        right, err := p.parseComparison()
        if err != nil {
            return nil, err
        }
        left = &BinaryOp{Op: strings.ToUpper(op), Left: left, Right: right}
    }
    return left, nil
}

func (p *Parser) parseComparison() (Expr, error) {
    left, err := p.parsePrimary()
    if err != nil {
        return nil, err
    }
    tok := p.peek()
    switch tok.Type {
    case TokenEq, TokenNeq, TokenLt, TokenLte, TokenGt, TokenGte:
        op := p.advance().Literal
        right, err := p.parsePrimary()
        if err != nil {
            return nil, err
        }
        return &BinaryOp{Op: op, Left: left, Right: right}, nil
    }
    return left, nil
}

func (p *Parser) parsePrimary() (Expr, error) {
    tok := p.advance()
    switch tok.Type {
    case TokenInt:
        n, err := strconv.ParseInt(tok.Literal, 10, 64)
        if err != nil {
            return nil, err
        }
        return &IntLiteral{Value: n}, nil
    case TokenFloat:
        f, err := strconv.ParseFloat(tok.Literal, 64)
        if err != nil {
            return nil, err
        }
        return &FloatLiteral{Value: f}, nil
    case TokenString:
        return &StrLiteral{Value: tok.Literal}, nil
    case TokenTrue:
        return &BoolLiteral{Value: true}, nil
    case TokenFalse:
        return &BoolLiteral{Value: false}, nil
    case TokenNull:
        return &NullLiteral{}, nil
    case TokenIdent:
        return &ColumnRef{Name: tok.Literal}, nil
    }
    return nil, fmt.Errorf("unexpected token in expression: %s", tok.Literal)
}

func (p *Parser) peek() Token {
    if p.pos >= len(p.tokens) {
        return Token{Type: TokenEOF}
    }
    return p.tokens[p.pos]
}

func (p *Parser) advance() Token {
    t := p.peek()
    p.pos++
    return t
}

func (p *Parser) consume(expected TokenType) error {
    t := p.advance()
    if t.Type != expected {
        return fmt.Errorf("expected token type %d, got %d (%q)", expected, t.Type, t.Literal)
    }
    return nil
}

func (p *Parser) expectIdent() (string, error) {
    t := p.advance()
    if t.Type != TokenIdent {
        return "", fmt.Errorf("expected identifier, got %q", t.Literal)
    }
    return t.Literal, nil
}

func (p *Parser) expectInt() (int64, error) {
    t := p.advance()
    if t.Type != TokenInt {
        return 0, fmt.Errorf("expected integer, got %q", t.Literal)
    }
    return strconv.ParseInt(t.Literal, 10, 64)
}

func parseType(s string) storage.TypeID {
    switch strings.ToLower(s) {
    case "int", "integer", "bigint":
        return storage.TypeInt
    case "float", "real", "double":
        return storage.TypeFloat
    case "bool", "boolean":
        return storage.TypeBool
    case "varchar", "text", "string", "char":
        return storage.TypeString
    }
    return storage.TypeString
}
```

---

## 5. Testing the Parser

```go
package query_test

import (
    "testing"
    "github.com/yourname/vaultdb/query"
)

func TestParseSelect(t *testing.T) {
    stmt, err := query.Parse("SELECT name, age FROM users WHERE age > 18")
    if err != nil {
        t.Fatal(err)
    }
    sel, ok := stmt.(*query.SelectStmt)
    if !ok {
        t.Fatal("expected SelectStmt")
    }
    if sel.Table != "users" {
        t.Errorf("table: want 'users', got %q", sel.Table)
    }
    if len(sel.Columns) != 2 || sel.Columns[0] != "name" {
        t.Errorf("columns: %v", sel.Columns)
    }
    if sel.Where == nil {
        t.Error("expected WHERE clause")
    }
}

func TestParseInsert(t *testing.T) {
    stmt, err := query.Parse("INSERT INTO users (name, age) VALUES ('Alice', 25)")
    if err != nil {
        t.Fatal(err)
    }
    ins, ok := stmt.(*query.InsertStmt)
    if !ok {
        t.Fatal("expected InsertStmt")
    }
    if len(ins.Values) != 2 {
        t.Errorf("expected 2 values, got %d", len(ins.Values))
    }
}
```

---

## Summary

- The lexer splits SQL into tokens (keywords, identifiers, literals, operators).
- The parser builds an AST (Abstract Syntax Tree) — a Go struct tree representing the query.
- Recursive descent parsing: one function per grammar rule, each calls others as needed.
- The AST is the contract between the parser and the executor — the executor never sees raw SQL.

### Exercises

**Easy:** Add support for the `IN` operator: `WHERE status IN ('active', 'pending')`. Add a new `InExpr` node to the AST and handle it in the parser.

**Medium:** Add support for `COUNT(*)`, `SUM(col)`, and `AVG(col)` aggregate functions. Add an `AggregateExpr` AST node and parse them in the column list of SELECT.

**Hard:** Add support for aliases: `SELECT name AS username, age AS years_old FROM users`. Update `SelectStmt` to track column aliases and add the alias to `ColumnRef`.
