# Chapter 07: Introduction to SQL — Talking to Databases in Plain English

Every application you have ever used — Instagram, Spotify, your school's grade portal — stores its data somewhere. That "somewhere" is almost always a database, and the language those databases understand is SQL. In this chapter, you will learn to speak that language, set up a real database on your own computer, and connect to it from Go.

## Table of Contents

1. [What is SQL and Why Does It Exist?](#1-what-is-sql-and-why-does-it-exist)
2. [SQL is Declarative: Say What, Not How](#2-sql-is-declarative-say-what-not-how)
3. [Tables: The Spreadsheets of Databases](#3-tables-the-spreadsheets-of-databases)
4. [Data Types in SQL](#4-data-types-in-sql)
5. [NULL: The Value That Means "I Don't Know"](#5-null-the-value-that-means-i-dont-know)
6. [The Five Fundamental SQL Statements](#6-the-five-fundamental-sql-statements)
7. [Your First Real Queries: A Library Database](#7-your-first-real-queries-a-library-database)
8. [SQL Clients: Exploring Your Database Visually](#8-sql-clients-exploring-your-database-visually)
9. [Setting Up PostgreSQL Locally with Docker](#9-setting-up-postgresql-locally-with-docker)
10. [Connecting from Go: database/sql and pgx](#10-connecting-from-go-databasesql-and-pgx)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. What is SQL and Why Does It Exist?

Imagine you work at a library with a million books. Each book has a title, an author, a publication year, and a shelf location. If someone asks "Do you have any science fiction books published after 2010?", you need a system to find that answer quickly. You cannot check every shelf by hand.

In the 1970s, engineers at IBM faced the same problem — except their "books" were business records: customer orders, employee salaries, inventory counts. They invented a way to store data in a structured format and a language to ask questions about it. That language is **SQL** (pronounced either "S-Q-L" or "sequel" — both are acceptable).

SQL stands for **Structured Query Language**. The word "query" just means "question" or "request." When you write SQL, you are asking a question of your database: "Show me all customers who spent more than $100 last month."

SQL was standardized in 1986 by the American National Standards Institute (ANSI), which means that once you learn SQL, you can work with almost any database system — PostgreSQL, MySQL, SQLite, Microsoft SQL Server — because they all understand the same core language. Minor differences exist, but the fundamentals transfer completely.

This is why SQL survived for over 50 years and shows no signs of dying. It is the universal language of stored data.

### Quick Check

1. What does SQL stand for?
2. Why is the fact that SQL is standardized useful for a developer?
3. Name two database systems that understand SQL.

---

## 2. SQL is Declarative: Say What, Not How

This is the most important idea in this chapter. Read it carefully.

Think about two ways you could ask someone to make you a cup of tea:

**The "how" way (step by step):** "Fill the kettle with water. Plug it in. Wait for it to boil. Place a teabag in a mug. Pour the boiling water over the teabag. Wait three minutes. Remove the teabag. Add milk if desired."

**The "what" way:** "Can I have a cup of tea, please?"

SQL works like the second approach. You describe the result you want, and the database figures out how to get it. This style is called **declarative programming**.

Compare this to Go, which is **imperative** — you write the exact steps: loop through a slice, check each element, collect matches into a new slice. In SQL you say:

```sql
SELECT title FROM books WHERE genre = 'science fiction' AND year > 2010;
```

You did not say "scan every row, compare the genre field, check the year." You just described what you want. The database engine — a sophisticated piece of software — decides the most efficient way to retrieve it.

This matters because databases are optimized specifically for this kind of work. They maintain internal structures called **indexes** (like the index at the back of a textbook) that let them find data in milliseconds even across millions of rows. You benefit from all of that optimization simply by describing what you want.

### Quick Check

1. What does "declarative" mean in the context of SQL?
2. In Go, you write a loop to find elements in a slice. What is the SQL equivalent approach?
3. What is an index in a database, and what everyday object is it similar to?

---

## 3. Tables: The Spreadsheets of Databases

A database organizes data into **tables**. Think of a table exactly like a spreadsheet:

- Each **column** has a name and holds one type of data (like "title" or "year").
- Each **row** is one complete record (like one book).

Here is what a `books` table might look like:

| id | title                        | author              | year | genre           |
|----|------------------------------|---------------------|------|-----------------|
| 1  | The Left Hand of Darkness    | Ursula K. Le Guin   | 1969 | science fiction |
| 2  | Dune                         | Frank Herbert       | 1965 | science fiction |
| 3  | Pride and Prejudice          | Jane Austen         | 1813 | romance         |
| 4  | The Name of the Wind         | Patrick Rothfuss    | 2007 | fantasy         |

Every row in the table is a book. Every column is a piece of information about that book. The `id` column is special — it is a unique number assigned to each book so we can refer to it unambiguously. This is called a **primary key**.

A real database is a collection of many tables. Our library database might have a `books` table, a `members` table (for people who have library cards), and a `loans` table (recording who borrowed which book and when).

---

## 4. Data Types in SQL

When you create a table column, you must tell the database what kind of data it will hold. This is the column's **data type**. The database uses the data type to store data efficiently and to prevent mistakes — you would not want someone storing "banana" in a column meant to hold a year.

Here are the core data types you will use constantly:

### INTEGER

Whole numbers with no decimal point. Use this for counts, ages, quantities, and IDs.

```sql
year INTEGER,
page_count INTEGER
```

### TEXT (also called VARCHAR)

A string of characters — letters, numbers, symbols. Use this for names, titles, descriptions, email addresses. `VARCHAR(n)` means "text up to n characters long." In PostgreSQL, `TEXT` means unlimited length.

```sql
title TEXT,
email VARCHAR(255)
```

### REAL (also called FLOAT or NUMERIC)

Numbers with decimal points. Use this for prices, measurements, coordinates. For money, `NUMERIC(10, 2)` is preferred — it stores exactly 10 digits with 2 after the decimal, avoiding floating-point rounding errors.

```sql
price NUMERIC(10, 2),
latitude REAL
```

### BOOLEAN

True or false. Use this for flags and yes/no answers.

```sql
is_available BOOLEAN,
has_ebook BOOLEAN
```

### DATE

A calendar date: year, month, day. No time information.

```sql
published_date DATE,
birth_date DATE
```

### TIMESTAMP

A specific moment in time: date plus time, down to microseconds. Often stored with timezone information as `TIMESTAMPTZ` in PostgreSQL.

```sql
borrowed_at TIMESTAMP,
created_at TIMESTAMPTZ
```

### A Quick Type Comparison

| SQL Type      | Go equivalent      | Example value         |
|---------------|--------------------|-----------------------|
| INTEGER       | int, int64         | 42                    |
| TEXT          | string             | "The Hobbit"          |
| REAL          | float64            | 3.14                  |
| BOOLEAN       | bool               | true                  |
| DATE          | time.Time          | 2024-03-15            |
| TIMESTAMP     | time.Time          | 2024-03-15 14:30:00   |

### Quick Check

1. Which SQL data type would you use to store a person's age?
2. Why is `NUMERIC(10, 2)` preferred over `REAL` for storing prices?
3. What is the difference between `DATE` and `TIMESTAMP`?

---

## 5. NULL: The Value That Means "I Don't Know"

Here is a concept that trips up almost everyone when they first encounter it. SQL has a special value called **NULL**. NULL does not mean zero. It does not mean an empty string. It means **the value is unknown or missing**.

Think of a library membership form. It has a field for "phone number." Some members fill it in. Some leave it blank because they do not have a phone, or they just prefer not to share it. In the database, those blank phone numbers are stored as NULL — not as an empty string, but as "we do not have this information."

This creates some surprising behavior.

### NULL is not equal to anything, including itself

In Go, you can check `x == 0` or `x == ""`. In SQL, you cannot check `phone_number = NULL`. That comparison always returns false, even if the value really is NULL. Instead, you must write:

```sql
-- Correct: find members with no phone number
SELECT name FROM members WHERE phone_number IS NULL;

-- Also correct: find members who have a phone number
SELECT name FROM members WHERE phone_number IS NOT NULL;
```

### NULL in arithmetic

If you add NULL to a number, the result is NULL. If any part of an expression is unknown, the whole result is unknown. This makes sense when you think about it: if you do not know someone's age, you cannot calculate when they were born.

```sql
-- If salary is NULL, this returns NULL, not 0
SELECT salary + 1000 FROM employees;
```

### NULL in boolean logic

This is where NULL becomes genuinely tricky. SQL uses three-valued logic: TRUE, FALSE, and NULL (which means "unknown"). Consider:

- `TRUE AND NULL` = NULL (we do not know the answer)
- `FALSE AND NULL` = FALSE (one side is definitely false, so the whole thing is false)
- `TRUE OR NULL` = TRUE (one side is definitely true, so the whole thing is true)

### Handling NULL in Go

When a database column can contain NULL, you cannot scan it into a regular Go `string` or `int` — your program would crash. Instead, use the special nullable types from the `database/sql` package:

```go
var phoneNumber sql.NullString  // can hold a string or NULL
var age sql.NullInt64           // can hold an integer or NULL
```

We will see these in action when we write Go code later in this chapter.

### Quick Check

1. What does NULL mean in SQL? Is it the same as zero or an empty string?
2. How do you check if a column contains NULL in a SQL query?
3. Why can you not scan a nullable database column into a regular Go `string`?

---

## 6. The Five Fundamental SQL Statements

SQL has dozens of features, but five statements cover 95% of what you will ever do:

### CREATE — Build a New Table

Before you can store data, you need a place to store it. `CREATE TABLE` defines a new table with its columns and data types.

```sql
CREATE TABLE books (
    id        SERIAL PRIMARY KEY,
    title     TEXT NOT NULL,
    author    TEXT NOT NULL,
    year      INTEGER,
    genre     TEXT,
    available BOOLEAN NOT NULL DEFAULT true
);
```

Let us read this line by line:

- `id SERIAL PRIMARY KEY` — `SERIAL` means "automatically assign the next available integer." `PRIMARY KEY` means this column uniquely identifies each row.
- `title TEXT NOT NULL` — The title is text, and it cannot be NULL. Every book must have a title.
- `year INTEGER` — The year is optional (no `NOT NULL`), so it can be NULL if we do not know it.
- `available BOOLEAN NOT NULL DEFAULT true` — Cannot be NULL, and if we do not specify a value when inserting, it defaults to `true`.

### INSERT — Add a New Row

```sql
INSERT INTO books (title, author, year, genre)
VALUES ('Dune', 'Frank Herbert', 1965, 'science fiction');
```

We list the column names in parentheses, then provide the matching values. We did not include `id` (it is automatic) or `available` (it defaults to `true`).

### SELECT — Read Data

`SELECT` is the workhorse of SQL. It retrieves rows from a table.

```sql
-- Get everything from the books table
SELECT * FROM books;

-- Get only specific columns
SELECT title, author FROM books;

-- Get books with a condition (the WHERE clause)
SELECT title, author FROM books WHERE genre = 'science fiction';

-- Sort the results
SELECT title, year FROM books ORDER BY year DESC;

-- Limit how many rows you get back
SELECT title FROM books LIMIT 10;
```

The `*` means "all columns." Use it for quick exploration, but in real applications, always list the specific columns you need — it is clearer and faster.

### UPDATE — Change Existing Data

```sql
-- Mark a specific book as unavailable
UPDATE books SET available = false WHERE id = 2;

-- Fix a typo in an author's name
UPDATE books SET author = 'Ursula K. Le Guin' WHERE id = 7;
```

**Warning:** Always include a `WHERE` clause in UPDATE. Without it, you update every single row in the table.

### DELETE — Remove a Row

```sql
-- Remove one specific book
DELETE FROM books WHERE id = 5;

-- Remove all unavailable books
DELETE FROM books WHERE available = false;
```

**Warning:** Same as UPDATE — always use `WHERE` unless you actually want to delete everything.

### Quick Check

1. What does `PRIMARY KEY` mean on a column?
2. What does `NOT NULL` mean, and why would you use it?
3. What dangerous thing happens if you run `UPDATE books SET available = false` without a `WHERE` clause?

---

## 7. Your First Real Queries: A Library Database

Let us build a small but realistic library database and practice writing queries against it.

### The Schema

A **schema** is the structure of your database — the tables and their columns. Here is our library schema:

```sql
CREATE TABLE members (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    joined_on  DATE NOT NULL DEFAULT CURRENT_DATE,
    phone      TEXT
);

CREATE TABLE books (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    author      TEXT NOT NULL,
    year        INTEGER,
    genre       TEXT,
    available   BOOLEAN NOT NULL DEFAULT true
);

CREATE TABLE loans (
    id          SERIAL PRIMARY KEY,
    book_id     INTEGER NOT NULL REFERENCES books(id),
    member_id   INTEGER NOT NULL REFERENCES members(id),
    borrowed_on DATE NOT NULL DEFAULT CURRENT_DATE,
    due_on      DATE NOT NULL,
    returned_on DATE
);
```

The `REFERENCES` keyword creates a **foreign key** — a link between tables. `book_id INTEGER REFERENCES books(id)` means "this column must contain a value that exists in the `id` column of the `books` table." This prevents orphaned records (a loan for a book that does not exist).

### Inserting Sample Data

```sql
INSERT INTO members (name, email) VALUES
    ('Alice Chen', 'alice@example.com'),
    ('Bob Kumar', 'bob@example.com'),
    ('Carmen Silva', 'carmen@example.com');

INSERT INTO books (title, author, year, genre) VALUES
    ('Dune', 'Frank Herbert', 1965, 'science fiction'),
    ('The Left Hand of Darkness', 'Ursula K. Le Guin', 1969, 'science fiction'),
    ('Pride and Prejudice', 'Jane Austen', 1813, 'romance'),
    ('The Name of the Wind', 'Patrick Rothfuss', 2007, 'fantasy'),
    ('Project Hail Mary', 'Andy Weir', 2021, 'science fiction');

INSERT INTO loans (book_id, member_id, due_on) VALUES
    (1, 1, CURRENT_DATE + INTERVAL '14 days'),
    (3, 2, CURRENT_DATE + INTERVAL '14 days');

-- Mark those books as unavailable
UPDATE books SET available = false WHERE id IN (1, 3);
```

### Useful Queries

Now let us ask interesting questions:

```sql
-- Which books are currently available?
SELECT title, author FROM books WHERE available = true;

-- Which member borrowed "Dune"?
SELECT members.name, loans.borrowed_on, loans.due_on
FROM loans
JOIN members ON loans.member_id = members.id
JOIN books ON loans.book_id = books.id
WHERE books.title = 'Dune';

-- How many books does each genre have?
SELECT genre, COUNT(*) AS total
FROM books
GROUP BY genre
ORDER BY total DESC;

-- Which loans are overdue (no return date and past due date)?
SELECT books.title, members.name, loans.due_on
FROM loans
JOIN books ON loans.book_id = books.id
JOIN members ON loans.member_id = members.id
WHERE loans.returned_on IS NULL
  AND loans.due_on < CURRENT_DATE;
```

The `JOIN` keyword combines rows from two tables based on a matching condition. `COUNT(*)` counts the number of rows in each group. `GROUP BY` groups rows with the same genre together so `COUNT` can work on each group separately.

---

## 8. SQL Clients: Exploring Your Database Visually

Before writing Go code, it helps to explore your database using a dedicated tool. These are called **SQL clients** or **database clients**.

### psql — The Command-Line Client

`psql` ships with PostgreSQL. It is a text-based interface where you type SQL directly. It is fast, available everywhere, and essential to know.

```
$ psql -h localhost -U postgres -d library

library=# \dt          -- list all tables
library=# \d books     -- describe the books table (show columns)
library=# SELECT * FROM books;
library=# \q           -- quit
```

### DBeaver — Free, Cross-Platform GUI

DBeaver (dbeaver.io) is a free application that shows your database as a visual tree. You can click on a table to see its data, write queries in an editor, and view results in a spreadsheet-like grid. It works on Mac, Windows, and Linux and supports PostgreSQL, MySQL, SQLite, and dozens of others.

### TablePlus — Mac/Windows GUI

TablePlus (tableplus.com) is a polished, modern database client popular on macOS. It has a free tier that covers learning and small projects. The interface is clean and fast.

For this course, any of these three tools works. Use psql if you like the command line; use DBeaver or TablePlus if you prefer a visual interface.

---

## 9. Setting Up PostgreSQL Locally with Docker

Rather than installing PostgreSQL directly onto your computer (which involves several setup steps and can conflict with other software), we will use **Docker** to run it in a container.

Think of a Docker container as a small, self-contained computer running inside your computer. It has its own operating system, its own software, and when you delete it, every trace of it disappears. This makes it perfect for development.

### Installing Docker

Download Docker Desktop from docker.com and install it. Once installed, you will have the `docker` command available in your terminal.

### Starting a PostgreSQL Container

Run this single command:

```bash
docker run --name library-db \
  -e POSTGRES_USER=postgres \
  -e POSTGRES_PASSWORD=secret \
  -e POSTGRES_DB=library \
  -p 5432:5432 \
  -d postgres:16
```

Let us understand each part:

- `--name library-db` — Give the container a friendly name.
- `-e POSTGRES_USER=postgres` — Set the database username to `postgres`.
- `-e POSTGRES_PASSWORD=secret` — Set the password to `secret`.
- `-e POSTGRES_DB=library` — Create a database called `library` when the container starts.
- `-p 5432:5432` — Forward port 5432 on your computer to port 5432 inside the container. PostgreSQL listens on port 5432 by default.
- `-d` — Run in the background (detached mode).
- `postgres:16` — Use the official PostgreSQL version 16 image.

### Verifying It Works

```bash
# Check that the container is running
docker ps

# Connect with psql
psql -h localhost -U postgres -d library
# Enter the password: secret
```

If you see the `library=#` prompt, PostgreSQL is running and you are connected.

### Stopping and Starting the Container

```bash
# Stop the container (data is preserved)
docker stop library-db

# Start it again
docker start library-db

# Remove it entirely (data is lost!)
docker rm -f library-db
```

### Creating the Schema

Connect with psql and run the CREATE TABLE statements from section 7. Or save them to a file and run:

```bash
psql -h localhost -U postgres -d library -f schema.sql
```

---

## 10. Connecting from Go: database/sql and pgx

Now the exciting part — writing Go code that talks to our PostgreSQL database.

### The database/sql Package

Go's standard library includes `database/sql`. This package defines a common interface for all SQL databases. It does not know how to talk to PostgreSQL specifically — that is the job of a **driver**. The driver is a separate package that implements the `database/sql` interface for a specific database.

The most popular PostgreSQL driver for Go is **pgx** (github.com/jackc/pgx).

### Setting Up the Project

```bash
mkdir library-app && cd library-app
go mod init library-app
go get github.com/jackc/pgx/v5/stdlib
```

The `stdlib` sub-package of pgx provides a `database/sql`-compatible driver, which means we can use it with the standard `database/sql` interface.

### Connecting to the Database

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib" // Import for its side effect: registering the driver
)

func main() {
    // The connection string tells database/sql how to find our database.
    // Format: postgres://username:password@host:port/database_name
    connStr := "postgres://postgres:secret@localhost:5432/library"

    // sql.Open does not actually connect yet. It just validates the arguments
    // and prepares the connection pool.
    db, err := sql.Open("pgx", connStr)
    if err != nil {
        log.Fatal("failed to open database:", err)
    }
    defer db.Close() // Close the connection pool when main() returns

    // db.Ping() sends a real request to the database. This is where
    // a connection error would actually surface.
    if err := db.Ping(); err != nil {
        log.Fatal("failed to ping database:", err)
    }

    fmt.Println("Connected to the library database!")
}
```

Notice the `import _ "github.com/jackc/pgx/v5/stdlib"`. The underscore `_` means we are importing this package only for its **side effect** — when it loads, it calls `sql.Register("pgx", ...)` to register itself as a driver. We never call functions from this package directly.

### Inserting a Row

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib"
)

func main() {
    db, err := sql.Open("pgx", "postgres://postgres:secret@localhost:5432/library")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal(err)
    }

    // db.Exec runs a statement that does not return rows.
    // The $1 is a placeholder — the actual value is passed as the second argument.
    // This prevents SQL injection attacks.
    result, err := db.Exec(
        "INSERT INTO books (title, author, year, genre) VALUES ($1, $2, $3, $4)",
        "Project Hail Mary", "Andy Weir", 2021, "science fiction",
    )
    if err != nil {
        log.Fatal("failed to insert book:", err)
    }

    // RowsAffected tells us how many rows were changed.
    rowsAffected, _ := result.RowsAffected()
    fmt.Printf("Inserted %d book(s)\n", rowsAffected)
}
```

The `$1`, `$2`, `$3`, `$4` are **placeholders**. The database driver replaces them with the actual values you pass as arguments. This is critically important for security: it prevents **SQL injection**, where an attacker puts SQL code inside user input to manipulate your database.

Never build SQL queries by concatenating strings from user input:

```go
// DANGEROUS. Never do this.
query := "SELECT * FROM users WHERE email = '" + userInput + "'"

// Safe. Always do this.
query := "SELECT * FROM users WHERE email = $1"
db.Query(query, userInput)
```

### Querying Multiple Rows

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib"
)

// Book represents one row from the books table.
// We define a struct to hold the data — one field per column we are selecting.
type Book struct {
    ID        int
    Title     string
    Author    string
    Year      sql.NullInt64 // Year can be NULL
    Genre     sql.NullString // Genre can be NULL
    Available bool
}

func main() {
    db, err := sql.Open("pgx", "postgres://postgres:secret@localhost:5432/library")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    // db.Query runs a SELECT and returns a *sql.Rows — a cursor over the results.
    rows, err := db.Query("SELECT id, title, author, year, genre, available FROM books ORDER BY id")
    if err != nil {
        log.Fatal("query failed:", err)
    }
    // rows.Close() releases the database connection back to the pool.
    // Always defer this immediately after checking the error.
    defer rows.Close()

    // rows.Next() advances to the next row. It returns false when there are
    // no more rows, or when an error occurred.
    for rows.Next() {
        var b Book

        // rows.Scan reads the columns of the current row into our variables.
        // The order of arguments must match the order of columns in the SELECT.
        err := rows.Scan(&b.ID, &b.Title, &b.Author, &b.Year, &b.Genre, &b.Available)
        if err != nil {
            log.Fatal("scan failed:", err)
        }

        // sql.NullString has two fields: String (the value) and Valid (true if not NULL)
        genre := "unknown"
        if b.Genre.Valid {
            genre = b.Genre.String
        }

        fmt.Printf("[%d] %s by %s (%s) — available: %v\n",
            b.ID, b.Title, b.Author, genre, b.Available)
    }

    // After the loop, check if Next() stopped due to an error.
    if err := rows.Err(); err != nil {
        log.Fatal("rows error:", err)
    }
}
```

### Querying a Single Row

When you expect exactly one row (like fetching a book by its ID), use `db.QueryRow`:

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib"
)

type Book struct {
    ID        int
    Title     string
    Author    string
    Available bool
}

func getBookByID(db *sql.DB, id int) (Book, error) {
    var b Book

    // QueryRow returns a *sql.Row. Scanning it runs the query.
    // If no row is found, Scan returns sql.ErrNoRows.
    err := db.QueryRow(
        "SELECT id, title, author, available FROM books WHERE id = $1",
        id,
    ).Scan(&b.ID, &b.Title, &b.Author, &b.Available)

    if err == sql.ErrNoRows {
        return Book{}, fmt.Errorf("no book with id %d", id)
    }
    if err != nil {
        return Book{}, fmt.Errorf("query error: %w", err)
    }

    return b, nil
}

func main() {
    db, err := sql.Open("pgx", "postgres://postgres:secret@localhost:5432/library")
    if err != nil {
        log.Fatal(err)
    }
    defer db.Close()

    book, err := getBookByID(db, 1)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("Found: %s by %s\n", book.Title, book.Author)
}
```

### A Complete Mini Project: Library Query Tool

Let us put everything together into a small but complete program that creates the schema, inserts books, and answers queries.

```go
package main

import (
    "database/sql"
    "fmt"
    "log"

    _ "github.com/jackc/pgx/v5/stdlib"
)

const schema = `
CREATE TABLE IF NOT EXISTS books (
    id        SERIAL PRIMARY KEY,
    title     TEXT NOT NULL,
    author    TEXT NOT NULL,
    year      INTEGER,
    genre     TEXT,
    available BOOLEAN NOT NULL DEFAULT true
);
`

type Book struct {
    ID        int
    Title     string
    Author    string
    Year      sql.NullInt64
    Genre     sql.NullString
    Available bool
}

func main() {
    db, err := sql.Open("pgx", "postgres://postgres:secret@localhost:5432/library")
    if err != nil {
        log.Fatal("open:", err)
    }
    defer db.Close()

    if err := db.Ping(); err != nil {
        log.Fatal("ping:", err)
    }

    // Create the table (IF NOT EXISTS means it is safe to run multiple times)
    if _, err := db.Exec(schema); err != nil {
        log.Fatal("create schema:", err)
    }
    fmt.Println("Schema ready.")

    // Insert some books
    books := []struct {
        title, author, genre string
        year                 int
    }{
        {"Dune", "Frank Herbert", "science fiction", 1965},
        {"The Left Hand of Darkness", "Ursula K. Le Guin", "science fiction", 1969},
        {"The Name of the Wind", "Patrick Rothfuss", "fantasy", 2007},
    }

    for _, b := range books {
        _, err := db.Exec(
            "INSERT INTO books (title, author, year, genre) VALUES ($1, $2, $3, $4) ON CONFLICT DO NOTHING",
            b.title, b.author, b.year, b.genre,
        )
        if err != nil {
            log.Printf("insert %q: %v", b.title, err)
        }
    }
    fmt.Println("Books inserted.")

    // Count books per genre
    fmt.Println("\nBooks by genre:")
    rows, err := db.Query("SELECT genre, COUNT(*) FROM books GROUP BY genre ORDER BY COUNT(*) DESC")
    if err != nil {
        log.Fatal("genre query:", err)
    }
    defer rows.Close()

    for rows.Next() {
        var genre sql.NullString
        var count int
        if err := rows.Scan(&genre, &count); err != nil {
            log.Fatal("scan:", err)
        }
        g := "unknown"
        if genre.Valid {
            g = genre.String
        }
        fmt.Printf("  %-20s %d book(s)\n", g, count)
    }
    if err := rows.Err(); err != nil {
        log.Fatal("rows:", err)
    }
}
```

### Quick Check

1. What is the purpose of the blank import `_ "github.com/jackc/pgx/v5/stdlib"`?
2. Why do we use `$1` placeholders instead of building SQL strings by concatenation?
3. What is the difference between `db.Query` and `db.QueryRow`?

---

## Summary

- SQL (Structured Query Language) is the universal language for working with relational databases. It has been standardized since 1986 and works across PostgreSQL, MySQL, SQLite, and many other systems.
- SQL is **declarative**: you describe the data you want, and the database engine decides the most efficient way to retrieve it. You do not write loops or algorithms.
- Data is organized into **tables** with typed columns. Key data types include `INTEGER`, `TEXT`, `REAL`, `BOOLEAN`, `DATE`, and `TIMESTAMP`. The special value `NULL` means "unknown" and behaves differently from zero or empty string.
- The five core SQL statements are `CREATE` (define structure), `INSERT` (add data), `SELECT` (read data), `UPDATE` (change data), and `DELETE` (remove data). Always use a `WHERE` clause with `UPDATE` and `DELETE`.
- In Go, you connect to PostgreSQL using the `database/sql` standard library with the `pgx` driver. Always use parameter placeholders (`$1`, `$2`, ...) to prevent SQL injection. Use `sql.NullString`, `sql.NullInt64`, and similar types to handle nullable columns safely.

---

## Exercises

### Easy

1. Write a SQL statement that creates a `members` table with columns: `id` (auto-incrementing primary key), `name` (required text), `email` (required, unique text), and `joined_on` (date, defaults to today).

2. Write a SQL query that selects the `title` and `author` from a `books` table, but only for books in the `'fantasy'` genre, sorted alphabetically by title.

3. Write a SQL statement that updates the `available` column to `false` for the book with `id = 3`.

### Medium

4. Create the `books` and `members` tables from section 7 using Docker and psql. Insert at least five books and three members. Write a query that returns the names of members who have not borrowed any books (hint: use `LEFT JOIN` and check for NULL in the loans table).

5. Write a complete Go program that connects to your local PostgreSQL instance and prints the count of available and unavailable books as two separate lines.

6. The `loans` table has a `returned_on` column that is NULL until the book is returned. Write a SQL query that finds all books currently on loan (not yet returned) along with the name of the member who borrowed them and the due date.

### Hard

7. Write a Go function with this signature:

```go
func searchBooks(db *sql.DB, genre string, availableOnly bool) ([]Book, error)
```

The function should return all books matching the genre (case-insensitive). If `availableOnly` is true, only return books where `available = true`. Handle NULLs correctly. Write a `main` function that calls this and prints the results.

8. **Mini Project:** Build a command-line library tool that accepts a subcommand:
   - `add-book <title> <author> <year> <genre>` — inserts a book
   - `list-books` — prints all books in a formatted table
   - `checkout <book-id> <member-id>` — creates a loan record and marks the book unavailable
   - `return <loan-id>` — sets `returned_on` to today and marks the book available

   Use `os.Args` to read the subcommand and arguments. Handle errors gracefully — if a book is already checked out, print a clear message rather than crashing.
