---
title: Core Language & Data Structures
---
With the basics down, go deeper into the parts of Go that separate "I can write a script" from "I can write software other people rely on" — proper error handling, testing, and the data structures every backend system is built out of.

### Error Handling the Go Way
Go has no exceptions. Errors are values, returned and checked explicitly. This one design decision shapes almost every function signature you'll write.

**Resources:**
- [Error Handling](course:go-programming#14-error-handling)
- [Defer, Panic, Recover](course:go-programming#15-defer-panic-recover)

### Pointers and Memory
Understand when Go copies a value and when it shares one — this matters the moment you start passing large structs or structs with mutable state around.

**Resources:**
- [Pointers](course:go-programming#13-pointers)

### Testing in Go
Go's testing story is built into the language and toolchain (`go test`), not bolted on. Every project from here on assumes you can write and read a test file.

**Resources:**
- [Testing in Go](course:go-programming#16-testing-in-go)

### Data Structures & Algorithms
Arrays and slices, hash tables, trees, and graphs — the structures that show up in interviews and in the actual internals of the databases and caches you'll use later.

**Resources:**
- [Arrays and Slices](course:go-programming#08-arrays-and-slices)
- [Maps](course:go-programming#09-maps)
- [Slices Deep Dive](course:go-programming#30-slices-deep-dive)
- [Trees and BST](course:go-programming#35-trees-and-bst)
- [Graphs](course:go-programming#39-graphs)

### Complexity Analysis
> optional

Big-O notation isn't academic trivia — it's how you reason about whether your API will still be fast when the table has ten million rows instead of ten.

**Resources:**
- [Complexity Analysis](course:go-programming#29-complexity-analysis)

### Practice: Implement Core Data Structures
> branches-from: Data Structures & Algorithms

Implement the structures yourself — a linked list, a stack, a hash table — instead of only reading about them. This is what makes them stick.

**Resources:**
- [Data Structures project](project:data-structures)
