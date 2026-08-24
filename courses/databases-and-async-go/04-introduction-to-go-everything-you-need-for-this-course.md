# Chapter 04: Introduction to Go — Everything You Need for This Course

Imagine you needed to build a racing car from scratch, but you also got to design your own engine. Most people reach for familiar tools — the ones that exist already — even if those tools were not built for the job. Go is different: it was designed from the ground up by engineers at Google who were frustrated that the existing languages were too slow, too complex, or too bad at handling many things at once. In this chapter you will learn everything about Go that you need for the rest of this course — from installing it on your machine all the way to writing a small working application.

## Table of Contents

- [Why Go?](#why-go)
- [Installing Go and Setting Up VS Code](#installing-go-and-setting-up-vs-code)
- [Your First Go Program: Hello World, Line by Line](#your-first-go-program-hello-world-line-by-line)
- [Variables, Types, Functions, and Structs](#variables-types-functions-and-structs)
- [Error Handling in Go](#error-handling-in-go)
- [Goroutines and Channels — A Brief Intro](#goroutines-and-channels--a-brief-intro)
- [Packages and Modules](#packages-and-modules)
- [The database/sql Package](#the-databasesql-package)
- [Writing Clean Go Code](#writing-clean-go-code)
- [Mini Project: Command-Line Note-Taking App](#mini-project-command-line-note-taking-app)
- [Summary](#summary)
- [Exercises](#exercises)

---

## Why Go?

Think about the tools a surgeon uses. A scalpel is not trying to be a hammer. It is designed for exactly one job and it does that job extraordinarily well. Go (sometimes called Golang) is a programming language designed for exactly the kind of work we do in this course: writing programs that talk to databases, handle many users at once, and run reliably for months or years without crashing.

Here is why Go is the right language for this course.

### Speed

Go compiles directly to machine code — the raw instructions your computer's processor executes. There is no middleman. When you run a Go program, your computer runs it directly. This makes Go programs fast: typically within 2 to 5 times the speed of C or C++ programs, and many times faster than Python or JavaScript for CPU-intensive work.

For database engines and message brokers — which is what we build in this course — speed matters. Every microsecond saved when writing a page to disk or delivering a message is a win.

### Simplicity

Go has a deliberately small language specification. There are very few keywords (25 total), no inheritance, no operator overloading, no magic. When you read someone else's Go code, you can understand it quickly because there is usually only one idiomatic way to do something. This is not an accident — the Go designers considered "readability" a first-class feature of the language.

### Built-in Concurrency

Concurrency means doing multiple things at the same time. A database server must handle requests from hundreds of users simultaneously. A message broker must read from disk and send to multiple subscribers at once.

Most languages bolt concurrency on as an afterthought. Go was designed with it from the start. Goroutines (Go's version of lightweight threads) and channels (the communication mechanism between goroutines) are built into the language itself — not a library you import, but a core part of Go syntax.

### Great Database Libraries

Go's standard library includes `database/sql` — a universal interface for working with any SQL database. Whether you connect to PostgreSQL, MySQL, or SQLite, your code looks the same. The driver — the code that knows how to speak the specific database's language — is swapped out underneath, but your application code stays identical. This design is elegant and practical.

### Who Uses Go?

Go powers some of the most demanding software in the world: Docker (the containerization system), Kubernetes (the cluster orchestration system), CockroachDB (a distributed SQL database), etcd (a distributed key-value store used by Kubernetes), InfluxDB (a time-series database), and much more. The systems we build in this course — VaultDB and StreamFlow — are in excellent company.

---

### Quick Check

1. What does "compiles to machine code" mean, and why does it make Go fast?
2. Name two Go features that make it well-suited for building database systems.
3. What is the difference between Go's `database/sql` package and a database driver?

---

## Installing Go and Setting Up VS Code

### Installing Go

Think of installing Go like setting up a workshop. Before you can build furniture, you need a workbench, your tools laid out, and the sawdust somewhere other than the living room. Installing Go sets up that workshop on your computer.

Go to the official download page: **https://go.dev/dl/**

Pick the installer for your operating system:

- On macOS: download the `.pkg` file and run it. It installs everything automatically.
- On Windows: download the `.msi` file and run it.
- On Linux: download the `.tar.gz`, extract it, and follow the instructions on the download page to add Go to your PATH.

After installation, open a terminal (on macOS/Linux that is the Terminal app; on Windows it is PowerShell or Command Prompt) and type:

```
go version
```

You should see something like:

```
go version go1.23.0 darwin/amd64
```

The exact version number may differ — anything from Go 1.21 onward will work fine for this course. If you see an error saying the command was not found, close and reopen your terminal and try again. If it still fails, the Go installer may not have updated your PATH — check the official installation instructions for your operating system.

### Installing VS Code

VS Code (Visual Studio Code) is a free code editor made by Microsoft. It is lightweight, fast, and has excellent Go support. Download it from **https://code.visualstudio.com/**.

After installing VS Code, open it and install the Go extension:

1. Click the Extensions icon in the left sidebar (it looks like four squares).
2. Search for "Go" and look for the extension published by Google.
3. Click Install.

The first time you open a `.go` file, VS Code will offer to install a set of Go tools (`gopls`, `dlv`, and others). Click "Install All". These tools provide autocompletion, error highlighting, and debugging support.

### Creating Your Workspace

Create a folder somewhere on your computer for this course. A good location might be `~/learning/go-databases/` (on macOS/Linux) or `C:\learning\go-databases\` (on Windows). All the code for this course will live inside this folder.

---

## Your First Go Program: Hello World, Line by Line

In the same tradition as every other programming language, we start with a program that prints "Hello, World!" to the screen. Create a new file called `hello.go` inside a new subfolder called `chapter04/`:

```go
package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")
}
```

To run it, open a terminal, navigate to the folder containing `hello.go`, and type:

```
go run hello.go
```

You should see:

```
Hello, World!
```

Now let us examine every single line.

### Line 1: `package main`

Every Go file starts with a package declaration. A package is a collection of related Go files that work together — think of it like a chapter in a book. All the files in the same folder must declare the same package name.

The name `main` is special. It tells the Go compiler: "this is an executable program — something a user can run directly." Any other package name (like `database` or `util`) creates a library that other programs can import, but cannot be run on its own.

Rule of thumb: if your file has a `main` function (we will get to that in a moment), its package must be named `main`.

### Line 2: (blank line)

Go does not care about blank lines — they are just for readability. Blank lines between code sections make the structure clearer, like paragraph breaks in a book.

### Line 3: `import "fmt"`

This line imports the `fmt` package. The word "import" means "I want to use code that lives elsewhere." `fmt` is a package included in Go's standard library — the collection of packages that comes bundled with every Go installation. `fmt` provides functions for formatted input and output (fmt is short for "format").

If you need to import multiple packages, you can list them like this:

```go
import (
    "fmt"
    "os"
    "strings"
)
```

The parentheses group the imports together. This is the conventional style when you have more than one.

### Line 5: `func main() {`

`func` is the keyword for declaring a function. A function is a named block of code that does a specific job. Think of a function like a recipe: it has a name (like "make chocolate cake"), may take some inputs (ingredients), and produces an output (the cake).

`main` is the special function name that Go looks for when running an executable program. When you type `go run hello.go`, Go compiles your code and then calls the `main` function. Everything starts here.

The `()` after `main` is where you would list inputs (called "parameters"). Our `main` function takes no inputs, so the parentheses are empty.

The `{` opens the function body — all the code that runs when this function is called.

### Line 6: `fmt.Println("Hello, World!")`

This line calls the `Println` function from the `fmt` package. The dot (`.`) is how you access a function from a package: `packagename.FunctionName(...)`.

`Println` stands for "print line." It prints its argument to the terminal and then moves to a new line (the "ln" at the end).

`"Hello, World!"` is a string — a piece of text enclosed in double quotes. This is the argument we are passing to `Println`.

### Line 7: `}`

The closing brace ends the `main` function body. Every `{` must have a matching `}`.

---

### Quick Check

1. What does `package main` tell the Go compiler?
2. What happens if you remove the `import "fmt"` line and try to run the program?
3. Why is the function named `main` and not something else?

---

## Variables, Types, Functions, and Structs

### Variables

A variable is a named container that holds a value. Think of a variable like a labeled jar in your kitchen: a jar labeled "sugar" holds sugar, a jar labeled "salt" holds salt. The label is the name; what is inside is the value.

In Go, you declare variables using the `var` keyword or the short declaration operator `:=`:

```go
package main

import "fmt"

func main() {
    // Long form: declare a variable with an explicit type.
    // "var" says we are declaring a variable.
    // "name" is the variable's label (its name).
    // "string" is the type — this jar can only hold text.
    // "= "Alice"" sets the initial value.
    var name string = "Alice"

    // Short form: Go figures out the type from the value on the right.
    // This is the most common style in Go.
    // ":=" means "declare and assign at the same time."
    age := 17

    // fmt.Printf is like fmt.Println but lets you format the output.
    // %s means "insert a string here", %d means "insert an integer here".
    fmt.Printf("%s is %d years old\n", name, age)
}
```

Run this and you will see:

```
Alice is 17 years old
```

The `:=` operator only works inside a function. For package-level variables (declared outside any function), you must use `var`.

### Types

A type defines what kind of value a variable can hold, and what operations you can do with it. Go is a statically typed language — every variable has a type that is decided at compile time and cannot change. This is like labeling a jar: once it says "sugar," you cannot pour engine oil into it without something going wrong.

The most common types in Go:

| Type | What it holds | Example value |
|------|---------------|---------------|
| `string` | Text (any sequence of characters) | `"hello"` |
| `int` | Whole numbers (positive, negative, zero) | `42`, `-7`, `0` |
| `float64` | Decimal numbers | `3.14`, `-0.5` |
| `bool` | True or false | `true`, `false` |
| `[]string` | A list (slice) of strings | `[]string{"a", "b", "c"}` |
| `map[string]int` | A dictionary mapping strings to ints | `map[string]int{"apples": 5}` |

```go
package main

import "fmt"

func main() {
    // A slice is an ordered list of values.
    // Think of it like a numbered row of boxes.
    fruits := []string{"apple", "banana", "cherry"}

    // range iterates over the slice.
    // "i" is the index (position): 0, 1, 2...
    // "fruit" is the value at that position.
    for i, fruit := range fruits {
        fmt.Printf("Position %d: %s\n", i, fruit)
    }

    // A map is a dictionary: keys map to values.
    // Here, a fruit name (string) maps to a count (int).
    inventory := map[string]int{
        "apple":  5,
        "banana": 3,
    }

    // Look up a value by its key.
    // The second return value ("ok") is true if the key existed.
    count, ok := inventory["apple"]
    if ok {
        fmt.Printf("We have %d apples\n", count)
    }
}
```

### Functions

A function is a reusable block of code with a name. Here is a function that takes two numbers and returns their sum:

```go
package main

import "fmt"

// "add" is the function name.
// It takes two parameters: "a" and "b", both of type int.
// The "int" after the closing parenthesis is the return type —
// the type of value this function sends back to its caller.
func add(a int, b int) int {
    // "return" sends a value back to whoever called this function.
    return a + b
}

func main() {
    // Call the add function with arguments 3 and 4.
    // The result (7) is stored in "result".
    result := add(3, 4)
    fmt.Println(result) // Prints: 7
}
```

Go functions can return multiple values. This is unusual in most languages and becomes very important in the next section on error handling:

```go
package main

import "fmt"

// This function returns two values: a string and a bool.
// The two return types are listed inside parentheses.
func greet(name string) (string, bool) {
    if name == "" {
        // Return an empty string and false to signal: "this did not work."
        return "", false
    }
    // Return the greeting and true to signal: "this worked."
    return "Hello, " + name + "!", true
}

func main() {
    // Capture both return values.
    message, ok := greet("Carlos")
    if ok {
        fmt.Println(message) // Prints: Hello, Carlos!
    }

    // Try it with an empty name.
    message, ok = greet("")
    if !ok {
        fmt.Println("Name was empty, could not greet")
    }
}
```

### Structs

A struct is a custom type that groups related fields together. Think of a struct like a form: a library card has fields for name, member ID, and expiry date. A struct in Go lets you define that form and then create many filled-out copies of it.

```go
package main

import "fmt"

// Define a struct called "Person".
// It has three fields: Name (string), Age (int), and Email (string).
// By convention in Go, struct field names start with a capital letter
// when they should be accessible from other packages.
type Person struct {
    Name  string
    Age   int
    Email string
}

// A method is a function attached to a type.
// "(p Person)" is called the receiver — it means this method belongs to Person.
// You call it like: somePerson.Greet()
func (p Person) Greet() string {
    return fmt.Sprintf("Hi, I am %s and I am %d years old.", p.Name, p.Age)
}

func main() {
    // Create a Person value using field names.
    // This is called a struct literal.
    alice := Person{
        Name:  "Alice",
        Age:   17,
        Email: "alice@example.com",
    }

    // Access a field using the dot operator.
    fmt.Println(alice.Name)  // Prints: Alice

    // Call the method.
    fmt.Println(alice.Greet()) // Prints: Hi, I am Alice and I am 17 years old.

    // You can also use a pointer to a struct.
    // A pointer holds the memory address of a value rather than the value itself.
    // Think of it like writing down the address of a house rather than
    // physically carrying the house around.
    bob := &Person{Name: "Bob", Age: 19, Email: "bob@example.com"}
    bob.Age = 20 // Modify through the pointer.
    fmt.Println(bob.Age) // Prints: 20
}
```

Structs are used everywhere in Go code. In this course, we will use them to represent database rows, query results, connection configurations, and more.

---

### Quick Check

1. What is the difference between `var x int = 5` and `x := 5`?
2. What does a map do, and what is the real-world analogy for it?
3. What is a struct, and how does it differ from a regular variable?

---

## Error Handling in Go

### The Error Interface

In most programming languages, when something goes wrong a program throws an exception — a dramatic signal that interrupts normal execution. Go takes a different approach. In Go, errors are just values. A function that can fail returns two things: the result it was trying to compute, and an error value. If everything went well, the error is `nil` (Go's way of saying "nothing"). If something went wrong, the error describes what happened.

This design is like a waiter who always brings your order and a status note. If the kitchen has your dish, the note says "all good." If they are out of it, the note says "sorry, we are out of the salmon." You as the customer always receive both: the dish (or nothing) and the status note. You cannot ignore the status note — it is right there in front of you.

The `error` type in Go is an interface — a contract that says "anything that has an `Error() string` method counts as an error." The standard library provides a simple way to create errors:

```go
package main

import (
    "errors"
    "fmt"
)

// divide returns the result of a/b.
// If b is zero, it returns an error because dividing by zero is undefined.
func divide(a, b float64) (float64, error) {
    if b == 0 {
        // errors.New creates a new error value with the given message.
        return 0, errors.New("cannot divide by zero")
    }
    return a / b, nil // nil means "no error"
}

func main() {
    // Call divide and capture both return values.
    result, err := divide(10, 2)

    // Always check the error before using the result.
    // This is the most important habit in Go programming.
    if err != nil {
        fmt.Println("Error:", err)
        return // Stop here — there is no valid result to use.
    }
    fmt.Printf("10 / 2 = %.1f\n", result) // Prints: 10 / 2 = 5.0

    // Now try dividing by zero.
    result, err = divide(10, 0)
    if err != nil {
        fmt.Println("Error:", err) // Prints: Error: cannot divide by zero
        return
    }
    fmt.Printf("10 / 0 = %.1f\n", result) // This line never runs.
}
```

### The `fmt.Errorf` Function

Often you want to include context in an error message — not just "something went wrong" but "opening the file config.json failed because: file not found." The `fmt.Errorf` function lets you build error messages with formatting, just like `fmt.Sprintf`:

```go
package main

import (
    "fmt"
    "os"
)

// openConfig tries to open a configuration file.
// If it fails, it returns an error that includes the filename.
func openConfig(filename string) error {
    _, err := os.Open(filename)
    if err != nil {
        // %w wraps the original error so callers can inspect it later.
        // This is the recommended way to add context to an error.
        return fmt.Errorf("openConfig: could not open %s: %w", filename, err)
    }
    return nil
}

func main() {
    err := openConfig("config.json")
    if err != nil {
        // Prints something like:
        // openConfig: could not open config.json: open config.json: no such file or directory
        fmt.Println(err)
    }
}
```

### The Error Handling Habit

In Go, you will write `if err != nil` dozens of times. This might seem repetitive at first, but it has a real benefit: you can never accidentally ignore an error. Every error is explicitly handled or explicitly ignored (using the blank identifier `_` to discard it, which is a conscious choice, not an oversight). This habit is one of the most important things you can learn from Go, and it will serve you well throughout this entire course.

```go
// The blank identifier _ discards a value you don't need.
// Use it when you are deliberately ignoring a return value.
// Be careful: never ignore errors from operations that can actually fail.
result, _ := divide(10, 2) // OK here because we know 2 != 0
```

---

### Quick Check

1. What does it mean when a function returns `nil` as its error value?
2. What is the `error` interface, and what method must a type implement to satisfy it?
3. Why does Go prefer returning errors as values instead of using exceptions?

---

## Goroutines and Channels — A Brief Intro

### Goroutines

Imagine a restaurant kitchen. There is one chef who does everything: takes the order, chops vegetables, cooks the steak, plates the food, and hands it to the waiter. One task at a time. This is how synchronous programs work — one thing after another.

Now imagine the same kitchen with six chefs. While one is chopping vegetables, another is cooking the steak, and a third is plating a different dish. Everything happens at the same time. This is concurrency — multiple tasks progressing simultaneously.

A goroutine is a lightweight concurrent task in Go. You start a goroutine by writing the keyword `go` before a function call:

```go
package main

import (
    "fmt"
    "time"
)

func sayHello(name string) {
    fmt.Printf("Hello from %s!\n", name)
}

func main() {
    // Start three goroutines. Each one runs sayHello concurrently.
    // The "go" keyword says: start this function in a new goroutine
    // and move on immediately without waiting for it to finish.
    go sayHello("Alice")
    go sayHello("Bob")
    go sayHello("Carol")

    // Without this sleep, the main function would exit before the
    // goroutines have a chance to run.
    // In real programs, we use channels or sync.WaitGroup instead of sleeping.
    time.Sleep(100 * time.Millisecond)
    fmt.Println("All done.")
}
```

Goroutines are extremely lightweight — you can have hundreds of thousands of goroutines running in a single Go program. Under the hood, Go's runtime schedules them across your computer's CPU cores automatically.

### Channels

Channels are the way goroutines communicate. Think of a channel like a pipe: one goroutine pushes a value in one end, and another goroutine pulls it out the other end. This is safe because only one goroutine touches the value at a time.

```go
package main

import "fmt"

func sum(numbers []int, result chan int) {
    total := 0
    for _, n := range numbers {
        total += n
    }
    // Send the result into the channel.
    // The arrow <- means "send this value into the channel on the left."
    result <- total
}

func main() {
    numbers := []int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

    // make creates a new channel.
    // chan int means this channel carries int values.
    ch := make(chan int)

    // Split the work: one goroutine sums the first half,
    // another goroutine sums the second half.
    go sum(numbers[:5], ch)  // numbers[0] through numbers[4]
    go sum(numbers[5:], ch)  // numbers[5] through numbers[9]

    // Receive two results from the channel.
    // The arrow <- on its own (right-hand side of an assignment)
    // means "receive a value from this channel."
    // This blocks until a value is available.
    firstHalf := <-ch
    secondHalf := <-ch

    fmt.Println("Total:", firstHalf+secondHalf) // Prints: Total: 55
}
```

We will go much deeper into goroutines and channels throughout this course — especially in Part 2 when we build StreamFlow, a message broker that needs to handle thousands of concurrent connections. For now, just know: goroutines start concurrent work, channels let that concurrent work communicate.

---

### Quick Check

1. What does the `go` keyword do when placed before a function call?
2. Describe a real-world analogy for a channel between two goroutines.
3. Why do we need to receive from a channel in `main` after starting goroutines with `go sum(...)`?

---

## Packages and Modules

### Packages

A package is a collection of related Go files in the same directory. The standard library is a large collection of packages: `fmt` for formatted I/O, `os` for operating system operations, `strings` for string manipulation, `net/http` for building web servers, and many more. Throughout this course we will import packages from the standard library and from third-party developers.

The package name used in the `package` declaration at the top of a file is usually the same as the directory name. For example, a file at `database/page.go` would start with `package database`.

### Modules

A module is the top-level unit of Go code organization — it is a collection of packages that are versioned and distributed together. Think of a module like a published book: a module has a title (its module path), a version number, and a list of other modules it depends on.

Every Go project starts with a `go.mod` file. This file records:
- The module path (a unique name for your project, usually a GitHub URL like `github.com/you/myproject`)
- The Go version your code targets
- All external dependencies your code uses

To create a new module, navigate to your project directory in a terminal and run:

```
go mod init github.com/yourname/chapter04
```

This creates a `go.mod` file that looks like this:

```
module github.com/yourname/chapter04

go 1.23
```

### The go.sum File

When you add external dependencies, Go also creates a `go.sum` file. This file records the cryptographic checksums of every version of every dependency your project uses. It is a security feature: the next time someone builds your project, Go verifies that the downloaded dependencies exactly match the checksums recorded in `go.sum`. If they do not match, Go refuses to build. This prevents a type of attack where a package maintainer replaces a library with malicious code after publishing it.

You should commit both `go.mod` and `go.sum` to version control. Never edit `go.sum` by hand — Go manages it automatically.

### Adding Dependencies with `go get`

To add an external package to your project, use the `go get` command:

```
go get github.com/mattn/go-sqlite3
```

This downloads the package, adds it to `go.mod` as a dependency, and updates `go.sum`. You can then import it in your code:

```go
import _ "github.com/mattn/go-sqlite3"
```

The underscore `_` in front of the import path is special — it means "import this package for its side effects only." For database drivers, importing them this way registers the driver with Go's `database/sql` package without you needing to call any functions yourself. We will see this in action in the next section.

### Keeping Dependencies Tidy

Over time, as you add and remove imports, your `go.mod` file may accumulate dependencies that your code no longer uses. Running:

```
go mod tidy
```

cleans up: it removes unused dependencies and adds any missing ones that your imports require. It is a good habit to run `go mod tidy` before committing code.

---

### Quick Check

1. What is the purpose of the `go.mod` file?
2. What does `go.sum` protect against?
3. What does `go mod tidy` do, and when should you run it?

---

## The database/sql Package

### The Universal Database Interface

Go's standard library includes a package called `database/sql`. Think of it like an electrical outlet: no matter where you live — whether the appliance was built in Japan, Germany, or Brazil — as long as you have the right adapter (driver), you can plug in and get power. The `database/sql` package is the outlet; the driver is the adapter for your specific database.

The `database/sql` package defines:
- `sql.Open` — opens a database connection
- `db.Query` — runs a SELECT query and returns rows
- `db.QueryRow` — runs a SELECT query that returns exactly one row
- `db.Exec` — runs an INSERT, UPDATE, DELETE, or CREATE TABLE
- `rows.Scan` — reads a row's columns into Go variables
- `db.Begin` / `tx.Commit` / `tx.Rollback` — transactions

Here is a complete example using SQLite (the simplest database — no server needed, just a file):

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    // The underscore import registers the SQLite driver
    // with the database/sql package as a side effect.
    // We do not call anything from this package directly.
    _ "github.com/mattn/go-sqlite3"
)

func main() {
    // Open a connection to a SQLite database file called "notes.db".
    // If the file does not exist, SQLite creates it.
    // "sqlite3" is the driver name registered by the import above.
    db, err := sql.Open("sqlite3", "notes.db")
    if err != nil {
        // log.Fatal prints the error and then exits the program.
        // Use this for errors that mean "we cannot continue at all."
        log.Fatal(err)
    }
    // defer schedules a function call to happen when the surrounding
    // function returns — no matter how it returns (normally or via error).
    // This guarantees db.Close() is always called.
    defer db.Close()

    // Create a table if it does not already exist.
    // The backtick (`) starts and ends a raw string literal in Go —
    // convenient for multi-line strings like SQL queries.
    _, err = db.Exec(`
        CREATE TABLE IF NOT EXISTS notes (
            id      INTEGER PRIMARY KEY AUTOINCREMENT,
            content TEXT NOT NULL,
            created TEXT NOT NULL
        )
    `)
    if err != nil {
        log.Fatal(err)
    }

    // Insert a row.
    // The ? marks are placeholders. The actual values come after the query string.
    // This prevents SQL injection — a security vulnerability where
    // malicious input in the values could modify the query itself.
    _, err = db.Exec(
        `INSERT INTO notes (content, created) VALUES (?, ?)`,
        "Remember to buy milk",
        "2026-06-14",
    )
    if err != nil {
        log.Fatal(err)
    }

    // Query rows.
    // db.Query returns a *sql.Rows, which is a cursor over the result set.
    rows, err := db.Query(`SELECT id, content, created FROM notes`)
    if err != nil {
        log.Fatal(err)
    }
    // Always close rows when done. defer handles this even if we return early.
    defer rows.Close()

    // rows.Next() advances to the next row and returns true if there is one.
    // When there are no more rows, it returns false and the loop ends.
    for rows.Next() {
        var id int
        var content, created string

        // Scan reads the current row's columns into our variables.
        // The order of arguments to Scan must match the order of
        // columns in the SELECT statement.
        err := rows.Scan(&id, &content, &created)
        if err != nil {
            log.Fatal(err)
        }
        fmt.Printf("Note %d [%s]: %s\n", id, created, content)
    }

    // After iterating, check if rows.Next() exited due to an error.
    if err := rows.Err(); err != nil {
        log.Fatal(err)
    }
}
```

Before you can run this, you need to add the SQLite driver to your project:

```
go get github.com/mattn/go-sqlite3
```

This pattern — `sql.Open`, `db.Exec`, `db.Query`, `rows.Scan`, `rows.Close` — is the backbone of every database operation in this course. We will use it hundreds of times. Once it feels natural, working with any SQL database in Go will feel easy.

---

### Quick Check

1. Why does the `database/sql` package not include its own SQLite or PostgreSQL implementation?
2. What does `rows.Scan` do, and why does it take `&id` (with the ampersand) instead of just `id`?
3. What is SQL injection, and how do placeholders (`?`) prevent it?

---

## Writing Clean Go Code

### gofmt

One of Go's greatest gifts to programmers is `gofmt` — a tool that automatically formats your Go code according to a single, official style. There is no argument about where to put braces, how to indent, or whether to use spaces or tabs. `gofmt` decides, and everyone uses it. Every Go file in the world is formatted the same way.

Run it manually with:

```
gofmt -w yourfile.go
```

The `-w` flag tells gofmt to write the changes directly to the file instead of just printing the formatted version.

If you installed the VS Code Go extension, gofmt runs automatically every time you save a file. You never need to think about formatting.

### Naming Conventions

Go has clear, community-wide naming conventions:

- **Variable and function names** use camelCase: `userName`, `openDatabase`, `maxRetries`.
- **Exported names** (accessible from other packages) start with a capital letter: `Person`, `OpenDatabase`, `MaxRetries`.
- **Unexported names** (private to the current package) start with a lowercase letter: `person`, `openDatabase`, `maxRetries`.
- **Acronyms** are written in all caps or all lowercase, not mixed: `userID` not `userId`, `HTTPServer` not `HttpServer`.
- **Short variable names** are fine for short-lived variables: `i` for a loop index, `err` for an error, `db` for a database connection.
- **Single-letter receiver names** are conventional: if your type is `Page`, use `p` as the receiver name in methods: `func (p *Page) Write() error`.

### Comments

Go uses `//` for single-line comments and `/* ... */` for multi-line comments. Comments for exported functions, types, and variables should start with the name of the thing being documented:

```go
// Person represents a user of the system.
type Person struct {
    Name string
    Age  int
}

// Greet returns a friendly greeting string for this person.
func (p Person) Greet() string {
    return "Hello, " + p.Name + "!"
}
```

This style means Go's documentation tool (`go doc`) can extract and display your comments in a readable format automatically.

### Keep Functions Small

A function that is longer than about 30 lines is usually doing too many things. Break large functions into smaller, named helpers. Each function should do one thing and do it well. This makes your code easier to read, test, and debug.

---

## Mini Project: Command-Line Note-Taking App

Now let us put everything from this chapter together. We will build a command-line app that lets you add notes and list them. Notes are saved to a text file so they persist between runs.

This project uses: variables, functions, structs, error handling, file I/O, and command-line argument parsing — all the core Go skills from this chapter.

### Project Setup

Create a new directory and initialize a module:

```
mkdir chapter04/notes-app
cd chapter04/notes-app
go mod init github.com/yourname/notes-app
```

Create a file called `main.go`:

```go
package main

import (
    "bufio"
    "fmt"
    "os"
    "strings"
    "time"
)

// Note represents a single note with its content and the time it was created.
type Note struct {
    Content   string
    CreatedAt string
}

// notesFilePath is the path to the file where notes are stored.
// We keep this as a constant so it is easy to change in one place.
const notesFilePath = "notes.txt"

// saveNote appends a single note to the notes file.
// Each note is saved as one line: "YYYY-MM-DD HH:MM:SS | content"
func saveNote(note Note) error {
    // os.OpenFile opens a file with specific options.
    // os.O_APPEND means "add to the end" rather than overwriting.
    // os.O_CREATE means "create the file if it does not exist."
    // os.O_WRONLY means "open for writing only."
    // 0644 is the file permission: owner can read+write, others can only read.
    f, err := os.OpenFile(notesFilePath, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return fmt.Errorf("saveNote: could not open file: %w", err)
    }
    defer f.Close()

    // Format the note as a single line and write it to the file.
    line := fmt.Sprintf("%s | %s\n", note.CreatedAt, note.Content)
    _, err = f.WriteString(line)
    if err != nil {
        return fmt.Errorf("saveNote: could not write note: %w", err)
    }

    return nil
}

// loadNotes reads all notes from the notes file and returns them as a slice.
func loadNotes() ([]Note, error) {
    // os.Open opens a file for reading.
    f, err := os.Open(notesFilePath)
    if err != nil {
        // If the file does not exist yet, that is fine — return an empty slice.
        // os.IsNotExist checks whether an error is a "file not found" error.
        if os.IsNotExist(err) {
            return []Note{}, nil
        }
        return nil, fmt.Errorf("loadNotes: could not open file: %w", err)
    }
    defer f.Close()

    var notes []Note

    // bufio.NewScanner reads a file line by line efficiently.
    // Without bufio, reading line by line requires more complex code.
    scanner := bufio.NewScanner(f)
    for scanner.Scan() {
        line := scanner.Text() // Get the current line as a string.

        // Each line is "timestamp | content". Split on " | " to get both parts.
        parts := strings.SplitN(line, " | ", 2)
        if len(parts) != 2 {
            // Skip malformed lines rather than crashing.
            continue
        }

        notes = append(notes, Note{
            CreatedAt: parts[0],
            Content:   parts[1],
        })
    }

    // Check if the scanner stopped due to an error (not just end of file).
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("loadNotes: error reading file: %w", err)
    }

    return notes, nil
}

func main() {
    // os.Args is a slice of strings containing the command-line arguments.
    // os.Args[0] is always the program name itself.
    // os.Args[1] would be the first argument the user typed.
    if len(os.Args) < 2 {
        fmt.Println("Usage:")
        fmt.Println("  notes-app add <your note text>")
        fmt.Println("  notes-app list")
        os.Exit(1) // Exit with a non-zero code to signal an error.
    }

    // The first argument is the command: "add" or "list".
    command := os.Args[1]

    switch command {
    case "add":
        // The rest of the arguments form the note content.
        // strings.Join glues them back together with spaces between.
        if len(os.Args) < 3 {
            fmt.Println("Error: please provide the note text after 'add'")
            os.Exit(1)
        }
        content := strings.Join(os.Args[2:], " ")

        // time.Now() returns the current time.
        // .Format with a layout string produces a human-readable timestamp.
        // The layout "2006-01-02 15:04:05" is Go's reference time — a fixed
        // date Go uses as a template for time formatting. (Unusual but memorable!)
        createdAt := time.Now().Format("2006-01-02 15:04:05")

        note := Note{
            Content:   content,
            CreatedAt: createdAt,
        }

        err := saveNote(note)
        if err != nil {
            fmt.Println("Error saving note:", err)
            os.Exit(1)
        }

        fmt.Printf("Note saved: \"%s\"\n", content)

    case "list":
        notes, err := loadNotes()
        if err != nil {
            fmt.Println("Error loading notes:", err)
            os.Exit(1)
        }

        if len(notes) == 0 {
            fmt.Println("No notes yet. Add one with: notes-app add <text>")
            return
        }

        fmt.Printf("You have %d note(s):\n\n", len(notes))
        for i, note := range notes {
            // Printf with %3d formats the number in a field 3 characters wide
            // so the columns line up neatly even with double-digit note numbers.
            fmt.Printf("  %3d. [%s] %s\n", i+1, note.CreatedAt, note.Content)
        }

    default:
        fmt.Printf("Unknown command: %q\n", command)
        fmt.Println("Use 'add' or 'list'.")
        os.Exit(1)
    }
}
```

### Running the App

Build and run the app:

```
go run main.go add "Buy groceries"
go run main.go add "Read chapter 5 tonight"
go run main.go add "Call dentist on Monday"
go run main.go list
```

You should see output like:

```
Note saved: "Buy groceries"
Note saved: "Read chapter 5 tonight"
Note saved: "Call dentist on Monday"
You have 3 note(s):

    1. [2026-06-14 10:23:45] Buy groceries
    2. [2026-06-14 10:23:47] Read chapter 5 tonight
    3. [2026-06-14 10:23:49] Call dentist on Monday
```

The notes are saved to `notes.txt` in the same directory. Open that file and look at it — it is just plain text, one note per line. If you run the program again, the existing notes are still there because we open the file in append mode.

### What You Just Built

Look at what this small program demonstrates:

- A custom `Note` struct to represent data
- Two functions (`saveNote`, `loadNotes`) that each do one thing and return errors
- Proper error handling with `fmt.Errorf` for wrapping errors with context
- File I/O with `os.OpenFile` and `bufio.Scanner`
- Command-line argument parsing with `os.Args`
- A `switch` statement for routing to different behaviors
- Time formatting with `time.Now().Format`

These are the exact same patterns — just at a much smaller scale — that we will use throughout this course when building VaultDB and StreamFlow.

---

## Summary

- Go was designed for speed, simplicity, and concurrency — making it ideal for building database systems and message brokers. Its `database/sql` package provides a universal interface for any SQL database.

- Go's type system, structs, and multi-return functions give you expressive but readable code. The short declaration operator `:=` and `range` make everyday code concise.

- Error handling in Go treats errors as ordinary values. Every function that can fail returns an `error` as its last return value. Always checking `if err != nil` is the most important habit in Go.

- Goroutines and channels give Go first-class concurrency support. A goroutine is a lightweight concurrent task started with `go`; a channel is a typed pipe for communicating between goroutines.

- Modules and packages (`go.mod`, `go.sum`, `go get`, `go mod tidy`) manage dependencies in a reproducible and secure way. The `gofmt` tool enforces a single universal formatting style.

---

## Exercises

### Easy

1. Write a Go program that declares variables for your name, age, and favorite programming language, then prints them using `fmt.Printf` in a single formatted sentence.

2. Write a function called `isEven` that takes an `int` and returns a `bool` indicating whether the number is even. Write a `main` function that calls `isEven` for the numbers 1 through 10 and prints each result.

3. Create a struct called `Rectangle` with fields `Width` and `Height` (both `float64`). Add a method called `Area` that returns the rectangle's area. In `main`, create two rectangles and print both their dimensions and area.

### Medium

4. Write a function called `readLines` that takes a filename (as a `string`) and returns a `[]string` containing all the lines in the file, and an `error`. If the file does not exist, return an empty slice with no error. Test it by reading a file you create manually.

5. Write a program that starts three goroutines. Each goroutine should receive a number from 1 to 3 and print "Goroutine N is done" after sleeping for N seconds. The main function should wait for all three goroutines to finish before printing "All goroutines complete." (Hint: use a channel to signal completion.)

6. Extend the note-taking app from the Mini Project to support a `delete` command that takes a note number and removes that note from the file. For example: `notes-app delete 2` should remove the second note. Remember to handle the case where the note number is out of range.

### Hard

7. Write a function called `wordCount` that takes a string and returns a `map[string]int` where each key is a unique word and each value is how many times that word appears. Ignore punctuation and treat uppercase and lowercase versions of the same word as identical. Print the top 5 most frequent words in a sample paragraph.

8. Rewrite the note-taking app to use SQLite (`database/sql` + `github.com/mattn/go-sqlite3`) instead of a text file. Create a `notes` table with `id`, `content`, and `created_at` columns. The `add` command inserts a row, the `list` command queries all rows ordered by `created_at`, and the `delete` command deletes a row by its ID.

9. Write a program that demonstrates a producer-consumer pattern using goroutines and channels. The producer goroutine generates 100 numbers (1 through 100) and sends them to a channel. Three consumer goroutines each read from the same channel and print which consumer processed each number and what the number was. The main function waits for all work to be done before exiting.
