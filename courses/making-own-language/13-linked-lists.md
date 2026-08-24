# Chapter 13: Linked Lists — Nodes Connected by Pointers

> "To understand recursion, you must first understand recursion. To understand linked lists, you must first understand pointers." — Every CS professor, ever

---

## Overview

Arrays are great, but they have one fundamental limitation: their size is fixed at allocation time. If you allocate an array of 100 elements and need a 101st, you must allocate a completely new array, copy everything over, and throw away the old one. This is expensive.

Linked lists solve this problem by storing data in individual **nodes**, each of which holds both the data itself and a pointer to the next node. The nodes do not need to be contiguous in memory — they can be scattered anywhere in the heap, connected only by the chain of pointers.

This chapter covers:
- What a node is and how pointers connect them
- Singly linked lists: insert, delete, traverse
- Doubly linked lists: navigating backwards
- Circular linked lists: the last node points back to the first
- Memory trade-offs vs arrays
- When linked lists beat arrays — and when they do not
- A generic linked list in Go using `interface{}` and generics
- **Astra Build Milestone**: How the Astra lexer's token list is traversed by the parser

---

## What We're Building

By the end of this chapter, you will have built a full generic linked list in Go and you will understand how the Astra compiler's parser walks through a flat slice of tokens — which is conceptually identical to traversing a linked list — to build meaning from raw source code.

---

## Table of Contents

1. What Is a Node?
2. Singly Linked Lists
3. Insertion: Head, Tail, and Middle
4. Deletion: Three Cases
5. Traversal: Walking the Chain
6. Doubly Linked Lists
7. Circular Linked Lists
8. Memory Overhead vs Arrays
9. Cache Unfriendliness — Why It Matters
10. When to Use Linked Lists
11. Generic Linked List in Go
12. Astra Build Milestone: The Token List

---

## 1. What Is a Node?

A **node** is a small container that holds two things:

1. **Data** — the value being stored (an integer, a string, whatever you need)
2. **A pointer to the next node** — the "link" in the list

```
+--------+------+    +--------+------+    +--------+------+
|  data  | next |--->|  data  | next |--->|  data  | next |---> nil
|   42   |  *   |    |   17   |  *   |    |   99   | nil  |
+--------+------+    +--------+------+    +--------+------+
  Node 1               Node 2               Node 3
```

The very first node is called the **head**. The very last node is called the **tail**, and its `next` pointer is `nil` (or `null` in many languages) to signal the end of the list.

In Go, a node for integers looks like this:

```go
type IntNode struct {
    Data int
    Next *IntNode
}
```

Notice the `Next` field is `*IntNode` — a **pointer** to another `IntNode`. This self-referential definition is how linked lists work.

---

## 2. Singly Linked Lists

A **singly linked list** means each node has exactly one pointer: `Next`, pointing forward. You can only traverse the list in one direction — from head to tail.

Here is a complete singly linked list struct in Go:

```go
package main

import "fmt"

// Node holds data and a pointer to the next node.
type Node struct {
    Data int
    Next *Node
}

// LinkedList keeps track of the head (first node) and size.
type LinkedList struct {
    Head *Node
    Size int
}

// NewLinkedList creates an empty list.
func NewLinkedList() *LinkedList {
    return &LinkedList{}
}
```

An empty list looks like this:

```
Head --> nil
Size = 0
```

After inserting 10, 20, and 30:

```
Head --> [10|*] --> [20|*] --> [30|nil]
Size = 3
```

---

## 3. Insertion: Head, Tail, and Middle

There are three natural places to insert a new node.

### Inserting at the Head — O(1)

This is the fastest operation. Create a new node, point it at the current head, and make the new node the new head.

```go
func (l *LinkedList) PrependHead(data int) {
    newNode := &Node{Data: data, Next: l.Head}
    l.Head = newNode
    l.Size++
}
```

Before: `Head --> [20] --> [30] --> nil`
After inserting 10 at head: `Head --> [10] --> [20] --> [30] --> nil`

### Inserting at the Tail — O(n)

We must walk the entire list to find the last node, then attach the new node.

```go
func (l *LinkedList) AppendTail(data int) {
    newNode := &Node{Data: data}
    if l.Head == nil {
        l.Head = newNode
        l.Size++
        return
    }
    current := l.Head
    for current.Next != nil {
        current = current.Next
    }
    current.Next = newNode
    l.Size++
}
```

This is O(n) because we traverse all n nodes. We can make it O(1) by keeping a `Tail` pointer in the list struct — a common optimization.

### Inserting at a Position — O(n)

To insert at position k, walk to position k-1, then rewire pointers.

```go
func (l *LinkedList) InsertAt(pos int, data int) error {
    if pos < 0 || pos > l.Size {
        return fmt.Errorf("position %d out of range [0, %d]", pos, l.Size)
    }
    if pos == 0 {
        l.PrependHead(data)
        return nil
    }
    newNode := &Node{Data: data}
    current := l.Head
    for i := 0; i < pos-1; i++ {
        current = current.Next
    }
    newNode.Next = current.Next
    current.Next = newNode
    l.Size++
    return nil
}
```

The pointer rewiring looks like this:

```
Before:  ... --> [A] --> [C] --> ...
                  ^
                  insert B here

Step 1: newNode.Next = current.Next   (B points to C)
         ... --> [A] --> [C] --> ...
                  ^
                  [B] ------^

Step 2: current.Next = newNode        (A points to B)
         ... --> [A] --> [B] --> [C] --> ...
```

Order matters here: if you do step 2 before step 1, you lose the reference to C forever.

---

## 4. Deletion: Three Cases

Deleting a node requires rewiring the previous node's `Next` to skip over the deleted node.

```go
func (l *LinkedList) Delete(data int) bool {
    if l.Head == nil {
        return false
    }
    // Case 1: deleting the head node
    if l.Head.Data == data {
        l.Head = l.Head.Next
        l.Size--
        return true
    }
    // Case 2: deleting a middle or tail node
    current := l.Head
    for current.Next != nil {
        if current.Next.Data == data {
            current.Next = current.Next.Next   // skip the target node
            l.Size--
            return true
        }
        current = current.Next
    }
    return false  // Case 3: not found
}
```

Visualizing deletion of node [20]:

```
Before: [10] --> [20] --> [30] --> nil
         ^
         current

current.Next.Data == 20, so:
current.Next = current.Next.Next

After:  [10] --> [30] --> nil
```

The deleted node [20] is now unreachable. Go's garbage collector will reclaim its memory automatically.

---

## 5. Traversal: Walking the Chain

To visit every node, start at the head and follow `Next` pointers until you reach `nil`.

```go
func (l *LinkedList) Print() {
    current := l.Head
    for current != nil {
        fmt.Printf("%d", current.Data)
        if current.Next != nil {
            fmt.Print(" --> ")
        }
        current = current.Next
    }
    fmt.Println(" --> nil")
}

func (l *LinkedList) Contains(data int) bool {
    current := l.Head
    for current != nil {
        if current.Data == data {
            return true
        }
        current = current.Next
    }
    return false
}

func (l *LinkedList) ToSlice() []int {
    result := make([]int, 0, l.Size)
    current := l.Head
    for current != nil {
        result = append(result, current.Data)
        current = current.Next
    }
    return result
}
```

Traversal is always O(n) — you must visit every node.

---

## 6. Doubly Linked Lists

A **doubly linked list** adds a `Prev` pointer to each node, enabling backward traversal.

```go
type DNode struct {
    Data int
    Prev *DNode
    Next *DNode
}

type DoublyLinkedList struct {
    Head *DNode
    Tail *DNode
    Size int
}
```

The structure looks like this:

```
nil <-- [10] <--> [20] <--> [30] --> nil
         ^                    ^
        Head                 Tail
```

Each node knows its predecessor. This makes backwards traversal, and deletion of a known node, O(1) — because you already have the `Prev` pointer and don't need to walk from the head to find the predecessor.

```go
func (l *DoublyLinkedList) AppendTail(data int) {
    newNode := &DNode{Data: data, Prev: l.Tail}
    if l.Tail != nil {
        l.Tail.Next = newNode
    } else {
        l.Head = newNode
    }
    l.Tail = newNode
    l.Size++
}

func (l *DoublyLinkedList) DeleteNode(node *DNode) {
    if node.Prev != nil {
        node.Prev.Next = node.Next
    } else {
        l.Head = node.Next   // deleting the head
    }
    if node.Next != nil {
        node.Next.Prev = node.Prev
    } else {
        l.Tail = node.Prev   // deleting the tail
    }
    l.Size--
}

func (l *DoublyLinkedList) PrintForward() {
    current := l.Head
    for current != nil {
        fmt.Printf("[%d] ", current.Data)
        current = current.Next
    }
    fmt.Println()
}

func (l *DoublyLinkedList) PrintBackward() {
    current := l.Tail
    for current != nil {
        fmt.Printf("[%d] ", current.Data)
        current = current.Prev
    }
    fmt.Println()
}
```

Doubly linked lists are used in:
- **Browser history** (back and forward navigation)
- **LRU cache** (most recently used at tail, evict from head)
- **Text editors** (cursor movement)

---

## 7. Circular Linked Lists

In a **circular linked list**, the tail's `Next` pointer points back to the head instead of `nil`. There is no true "end" — the list wraps around.

```
     +------------------------------------------+
     |                                          |
     v                                          |
    [10] --> [20] --> [30] --> [40] --> [50] ---+
     ^
    Head
```

This is useful for:
- **Round-robin scheduling**: cycle through a list of processes, wrapping around endlessly
- **Ring buffers**: fixed-capacity circular queues
- **Josephus problem**: a classic algorithm puzzle

```go
type CNode struct {
    Data int
    Next *CNode
}

type CircularList struct {
    Head *CNode
    Size int
}

func (l *CircularList) Append(data int) {
    newNode := &CNode{Data: data}
    if l.Head == nil {
        newNode.Next = newNode  // points to itself
        l.Head = newNode
        l.Size++
        return
    }
    // find tail (the node whose Next == Head)
    current := l.Head
    for current.Next != l.Head {
        current = current.Next
    }
    current.Next = newNode
    newNode.Next = l.Head
    l.Size++
}

func (l *CircularList) Print(n int) {
    // print first n elements to avoid infinite loop
    current := l.Head
    for i := 0; i < n; i++ {
        fmt.Printf("%d ", current.Data)
        current = current.Next
    }
    fmt.Println()
}
```

---

## 8. Memory Overhead vs Arrays

Let's be honest: linked lists use significantly more memory than arrays.

For a list of n integers:

```
Array:
+----+----+----+----+----+
| 10 | 20 | 30 | 40 | 50 |   ← 5 × 8 bytes = 40 bytes (on 64-bit)
+----+----+----+----+----+

Linked List (singly):
+--------+------+    +--------+------+    +--------+------+
|   10   |  *   |--->|   20   |  *   |--->|   30   | nil  |
+--------+------+    +--------+------+    +--------+------+
  16 bytes            16 bytes              16 bytes
  (8 data + 8 ptr)
```

Each node stores:
- The data itself (8 bytes for int64)
- One pointer for `Next` (8 bytes on 64-bit systems)
- Go's allocator metadata overhead (typically 8-16 bytes)

So a linked list of integers takes **roughly 3x the memory** of an equivalent array. For a doubly linked list, it's even worse — another 8 bytes per node for `Prev`.

| Operation       | Array  | Linked List |
|-----------------|--------|-------------|
| Access by index | O(1)   | O(n)        |
| Insert at head  | O(n)   | O(1)        |
| Insert at tail  | O(1)*  | O(n)**      |
| Delete at head  | O(n)   | O(1)        |
| Delete at tail  | O(1)*  | O(n)**      |
| Search          | O(n)   | O(n)        |
| Memory          | Low    | High        |

*Amortized for dynamic arrays
**O(1) if you maintain a tail pointer

---

## 9. Cache Unfriendliness — Why It Matters

Modern CPUs are much faster than RAM. To compensate, they use **cache lines** — when you access one memory address, the CPU automatically loads the surrounding 64 bytes into the fast L1 cache, betting that you'll need nearby data soon.

Arrays exploit this beautifully. Elements are contiguous, so when you access element 0, elements 1 through 7 are already in cache when you need them.

Linked list nodes are **scattered randomly** across the heap:

```
Memory addresses (conceptual):
0x1000: Node [10] --> 0x5F30
0x5F30: Node [20] --> 0xA210
0xA210: Node [30] --> 0x2B80
0x2B80: Node [40] --> nil

Every pointer dereference = potential cache miss = CPU stalls waiting for RAM
```

This is called **pointer chasing** and it is one of the most common performance pitfalls in low-level programming. On modern hardware, an array traversal can be **10-50x faster** than a linked list traversal due purely to cache effects, even though both are "O(n)."

Big-O notation measures algorithmic complexity but ignores constants. Cache behavior is a hidden constant that can dominate real-world performance.

---

## 10. When to Use Linked Lists

**Use a linked list when:**
- You insert or delete frequently at the head or tail (O(1) operations)
- You never need random access by index
- The list size changes dramatically and unpredictably
- You're implementing a queue or deque with a fixed API
- You hold a pointer to a specific node and need to delete it in O(1)

**Use an array (slice) when:**
- You need random access: `list[i]` — this is O(n) in a linked list
- You iterate more than you insert or delete
- Memory is tight
- Cache performance matters (embedded systems, game engines, HPC)
- You don't know what data structure to use — arrays are the default choice

In practice, modern language runtimes and CPUs are so optimized for sequential memory access that slices beat linked lists in most benchmarks, even for workloads where linked lists theoretically win. Go's `container/list` package exists but is rarely used in idiomatic Go code.

---

## 11. Generic Linked List in Go

Let's build a proper generic linked list using Go 1.18+ generics:

```go
package linkedlist

import (
    "fmt"
    "strings"
)

// Node[T] is a generic node.
type Node[T any] struct {
    Data T
    Next *Node[T]
}

// LinkedList[T] is a generic singly linked list.
type LinkedList[T comparable] struct {
    head *Node[T]
    tail *Node[T]
    size int
}

// New creates an empty linked list.
func New[T comparable]() *LinkedList[T] {
    return &LinkedList[T]{}
}

// Len returns the number of elements.
func (l *LinkedList[T]) Len() int { return l.size }

// IsEmpty returns true if the list has no elements.
func (l *LinkedList[T]) IsEmpty() bool { return l.size == 0 }

// Prepend inserts data at the head — O(1).
func (l *LinkedList[T]) Prepend(data T) {
    newNode := &Node[T]{Data: data, Next: l.head}
    l.head = newNode
    if l.tail == nil {
        l.tail = newNode
    }
    l.size++
}

// Append inserts data at the tail — O(1) because we track tail.
func (l *LinkedList[T]) Append(data T) {
    newNode := &Node[T]{Data: data}
    if l.tail != nil {
        l.tail.Next = newNode
    } else {
        l.head = newNode
    }
    l.tail = newNode
    l.size++
}

// Get returns the element at position i — O(n).
func (l *LinkedList[T]) Get(i int) (T, error) {
    var zero T
    if i < 0 || i >= l.size {
        return zero, fmt.Errorf("index %d out of bounds (size %d)", i, l.size)
    }
    current := l.head
    for k := 0; k < i; k++ {
        current = current.Next
    }
    return current.Data, nil
}

// Delete removes the first occurrence of data — O(n).
func (l *LinkedList[T]) Delete(data T) bool {
    if l.head == nil {
        return false
    }
    if l.head.Data == data {
        l.head = l.head.Next
        if l.head == nil {
            l.tail = nil
        }
        l.size--
        return true
    }
    current := l.head
    for current.Next != nil {
        if current.Next.Data == data {
            if current.Next == l.tail {
                l.tail = current
            }
            current.Next = current.Next.Next
            l.size--
            return true
        }
        current = current.Next
    }
    return false
}

// Contains returns true if the list contains data — O(n).
func (l *LinkedList[T]) Contains(data T) bool {
    current := l.head
    for current != nil {
        if current.Data == data {
            return true
        }
        current = current.Next
    }
    return false
}

// ForEach calls fn on each element in order.
func (l *LinkedList[T]) ForEach(fn func(T)) {
    current := l.head
    for current != nil {
        fn(current.Data)
        current = current.Next
    }
}

// Reverse reverses the list in-place — O(n).
func (l *LinkedList[T]) Reverse() {
    var prev *Node[T]
    current := l.head
    l.tail = l.head
    for current != nil {
        next := current.Next
        current.Next = prev
        prev = current
        current = next
    }
    l.head = prev
}

// String returns a human-readable representation.
func (l *LinkedList[T]) String() string {
    parts := make([]string, 0, l.size)
    current := l.head
    for current != nil {
        parts = append(parts, fmt.Sprintf("%v", current.Data))
        current = current.Next
    }
    return "[" + strings.Join(parts, " -> ") + "]"
}
```

And a usage example:

```go
package main

import (
    "fmt"
    ll "your-module/linkedlist"
)

func main() {
    list := ll.New[int]()
    list.Append(10)
    list.Append(20)
    list.Append(30)
    list.Prepend(5)

    fmt.Println(list)          // [5 -> 10 -> 20 -> 30]
    fmt.Println(list.Len())    // 4

    list.Delete(20)
    fmt.Println(list)          // [5 -> 10 -> 30]

    list.Reverse()
    fmt.Println(list)          // [30 -> 10 -> 5]

    val, _ := list.Get(1)
    fmt.Println(val)           // 10

    // Generic: also works with strings
    words := ll.New[string]()
    words.Append("hello")
    words.Append("world")
    fmt.Println(words)         // [hello -> world]
}
```

---

## 12. Astra Build Milestone: The Token List

In the Astra compiler, the **lexer** reads source code character by character and produces a flat list of **tokens**. Each token represents one meaningful unit of the language: a keyword, an identifier, a number, an operator, a brace.

```
Source code:
    fn add(a: int, b: int) -> int { return a + b }

Tokens:
    FN "fn"
    IDENTIFIER "add"
    LPAREN "("
    IDENTIFIER "a"
    COLON ":"
    INT_TYPE "int"
    COMMA ","
    IDENTIFIER "b"
    COLON ":"
    INT_TYPE "int"
    RPAREN ")"
    ARROW "->"
    INT_TYPE "int"
    LBRACE "{"
    RETURN "return"
    IDENTIFIER "a"
    PLUS "+"
    IDENTIFIER "b"
    RBRACE "}"
    EOF ""
```

The token list is the "linked list" connecting lexing to parsing. In the Astra compiler we actually use a Go slice — which is contiguous in memory and cache-friendly — but the parser's traversal API looks exactly like traversing a linked list: advance one token at a time, peek ahead, check for the end.

```go
// lexer/token.go

package lexer

// TokenType identifies what kind of token this is.
type TokenType int

const (
    // Literals
    INTEGER    TokenType = iota
    FLOAT
    STRING_LIT
    TRUE
    FALSE

    // Identifiers
    IDENTIFIER

    // Keywords
    FN
    LET
    CONST
    IF
    ELSE
    FOR
    WHILE
    RETURN
    STRUCT
    IMPL
    IMPORT
    IN
    MATCH
    ENUM
    TRAIT
    PUB

    // Types
    INT_TYPE
    FLOAT_TYPE
    BOOL_TYPE
    STRING_TYPE

    // Operators
    PLUS        // +
    MINUS       // -
    STAR        // *
    SLASH       // /
    PERCENT     // %
    EQ          // ==
    NEQ         // !=
    LT          // <
    GT          // >
    LEQ         // <=
    GEQ         // >=
    AND         // &&
    OR          // ||
    NOT         // !
    ASSIGN      // =
    PLUS_ASSIGN // +=
    MINUS_ASSIGN// -=
    DOTDOT      // ..  (range)

    // Delimiters
    LPAREN    // (
    RPAREN    // )
    LBRACE    // {
    RBRACE    // }
    LBRACKET  // [
    RBRACKET  // ]
    COMMA     // ,
    COLON     // :
    SEMICOLON // ;
    DOT       // .
    ARROW     // ->

    // Special
    EOF
    ILLEGAL
)

// Token is a single unit produced by the lexer.
type Token struct {
    Type    TokenType
    Lexeme  string  // the raw text from source, e.g. "42" or "fn"
    Line    int     // source line (1-indexed, for error messages)
    Column  int     // source column (1-indexed)
}

func (t Token) String() string {
    return fmt.Sprintf("Token(%d, %q, %d:%d)", t.Type, t.Lexeme, t.Line, t.Column)
}
```

Now the parser. It holds a **slice** of tokens (produced by the lexer) and a `current` index. Its API is identical to a linked list cursor:

```go
// parser/parser.go

package parser

import (
    "fmt"
    "your-module/lexer"
)

// Parser walks through the token list and builds an AST.
type Parser struct {
    tokens  []lexer.Token
    current int
    errors  []string
}

// NewParser creates a parser from the token slice.
func NewParser(tokens []lexer.Token) *Parser {
    return &Parser{tokens: tokens, current: 0}
}

// peek returns the current token without consuming it.
// Like reading the "head" of the remaining list.
func (p *Parser) peek() lexer.Token {
    return p.tokens[p.current]
}

// peekNext returns the token after current (lookahead).
func (p *Parser) peekNext() lexer.Token {
    if p.current+1 < len(p.tokens) {
        return p.tokens[p.current+1]
    }
    return p.tokens[len(p.tokens)-1] // return EOF
}

// advance consumes and returns the current token.
// Like calling Next() on a linked list iterator.
func (p *Parser) advance() lexer.Token {
    t := p.peek()
    if !p.isAtEnd() {
        p.current++
    }
    return t
}

// isAtEnd returns true when we've consumed all tokens.
func (p *Parser) isAtEnd() bool {
    return p.peek().Type == lexer.EOF
}

// check returns true if the current token has the given type.
func (p *Parser) check(t lexer.TokenType) bool {
    if p.isAtEnd() {
        return false
    }
    return p.peek().Type == t
}

// match consumes the current token if it matches any of the given types.
// Returns true if a match was found.
func (p *Parser) match(types ...lexer.TokenType) bool {
    for _, t := range types {
        if p.check(t) {
            p.advance()
            return true
        }
    }
    return false
}

// expect consumes the current token if it matches, or records an error.
func (p *Parser) expect(t lexer.TokenType, msg string) (lexer.Token, error) {
    if p.check(t) {
        return p.advance(), nil
    }
    tok := p.peek()
    err := fmt.Errorf("line %d:%d: %s, got %q", tok.Line, tok.Column, msg, tok.Lexeme)
    p.errors = append(p.errors, err.Error())
    return tok, err
}
```

The parser traversal mirrors linked list traversal exactly:

```
Linked List Traversal:       Parser Token Traversal:

current = head               current = 0
while current != nil:        while !isAtEnd():
    process(current.Data)        process(peek())
    current = current.Next       advance()
```

Here is how the parser uses these primitives to parse a `let` statement:

```go
// parseLet parses:  let name = expression
func (p *Parser) parseLet() (*LetStmt, error) {
    // We already consumed "let" via match()
    nameToken, err := p.expect(lexer.IDENTIFIER, "expected variable name after 'let'")
    if err != nil {
        return nil, err
    }

    _, err = p.expect(lexer.ASSIGN, "expected '=' after variable name")
    if err != nil {
        return nil, err
    }

    value, err := p.parseExpression()
    if err != nil {
        return nil, err
    }

    return &LetStmt{Name: nameToken.Lexeme, Value: value}, nil
}
```

Each call to `expect` or `match` is like calling `advance()` on a linked list — we consume one token and move to the next.

The beautiful thing about the token list approach is that it separates **lexing** (turning characters into tokens) from **parsing** (turning tokens into meaning). The lexer produces the list; the parser is a cursor walking through it. They never need to run at the same time, and the parser can look ahead as many tokens as it needs without the lexer needing to know.

---

## Astra Build Milestone: Putting It Together

```go
package main

import (
    "fmt"
    "your-module/lexer"
    "your-module/parser"
)

func main() {
    source := `fn add(a: int, b: int) -> int { return a + b }`

    // Lex: characters --> tokens
    l := lexer.New(source)
    tokens := l.Scan()  // returns []Token

    fmt.Println("=== Tokens ===")
    for i, tok := range tokens {
        fmt.Printf("[%2d] %s\n", i, tok)
    }

    // Parse: tokens --> AST
    p := parser.NewParser(tokens)
    ast, err := p.ParseProgram()
    if err != nil {
        fmt.Println("Parse errors:", p.Errors())
        return
    }

    fmt.Println("\n=== AST ===")
    ast.Print(0)
}
```

Output (conceptual):
```
=== Tokens ===
[ 0] Token(FN, "fn", 1:1)
[ 1] Token(IDENTIFIER, "add", 1:4)
[ 2] Token(LPAREN, "(", 1:7)
...
[17] Token(EOF, "", 1:47)

=== AST ===
Program
  FunctionDecl: add
    Param: a (int)
    Param: b (int)
    ReturnType: int
    Body:
      ReturnStmt
        BinaryExpr: +
          Identifier: a
          Identifier: b
```

The token list is the bridge between the raw text of a program and its structural meaning. Without it — without understanding how to walk a sequence of items one at a time — the compiler could not exist.

---

## Exercises

1. **Reverse a linked list** without allocating any new nodes. Only rewire the `Next` pointers. Test with [1 -> 2 -> 3 -> 4 -> 5] to get [5 -> 4 -> 3 -> 2 -> 1].

2. **Find the middle node** of a linked list in a single pass (without knowing the size). Hint: use two pointers, one moving twice as fast as the other.

3. **Detect a cycle** in a linked list. If the tail points back to some earlier node instead of nil, the list has a cycle. Detect this in O(n) time and O(1) space.

4. **Merge two sorted linked lists** into one sorted linked list. E.g. [1->3->5] + [2->4->6] = [1->2->3->4->5->6].

5. **Implement a doubly linked list** with these operations: `AddFront`, `AddBack`, `RemoveFront`, `RemoveBack`, `PrintForward`, `PrintBackward`.

6. **Remove nth node from end** in one pass. Given n=2 and list [1->2->3->4->5], remove the 2nd from end to get [1->2->3->5].

7. **Implement an LRU cache** using a doubly linked list + hash map. Operations: `Get(key)` and `Put(key, value)`, both O(1).

8. **Astra extension**: Add a `previous()` method to the Parser that moves `current` backward by one. When would this be useful in a recursive descent parser?

---

## Summary

| Concept             | Key Point                                           |
|---------------------|-----------------------------------------------------|
| Node                | Data + pointer(s) to neighboring nodes             |
| Singly linked list  | One `Next` pointer, forward traversal only         |
| Doubly linked list  | `Prev` + `Next` pointers, bidirectional traversal  |
| Circular list       | Tail's `Next` points back to head                  |
| Insert at head      | O(1) — just rewire the head pointer                |
| Insert at tail      | O(n) without tail pointer, O(1) with tail pointer  |
| Delete              | O(n) to find, O(1) to rewire                       |
| Access by index     | O(n) — must walk from head                         |
| Memory overhead     | ~3x arrays (data + pointer + allocator metadata)   |
| Cache performance   | Poor — nodes scattered in heap (pointer chasing)   |
| Best use case       | Frequent head/tail insert/delete, no random access |
| Astra usage         | Token list: produced by lexer, consumed by parser   |
