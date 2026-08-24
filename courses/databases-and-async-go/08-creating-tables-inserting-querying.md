# Chapter 08: Creating Tables, Inserting Data, and Querying

You have learned the SQL language in the previous chapter. Now you are going to go deep. You are going to learn how to build the structure of a database, fill it with data, and retrieve exactly what you need — with precision and safety. By the end of this chapter, you will have a fully working Go application that stores and retrieves blog posts from a real PostgreSQL database.

## Table of Contents

- [The Spreadsheet Analogy — Tables, Rows, and Columns](#the-spreadsheet-analogy)
- [CREATE TABLE: Building the Structure](#create-table-building-the-structure)
- [INSERT INTO: Adding Data](#insert-into-adding-data)
- [SELECT: Retrieving Data](#select-retrieving-data)
- [WHERE: Filtering Rows](#where-filtering-rows)
- [ORDER BY and LIMIT: Sorting and Pagination](#order-by-and-limit)
- [UPDATE: Changing Existing Data](#update-changing-existing-data)
- [DELETE: Removing Data (and Why You Should Think Twice)](#delete-removing-data)
- [Soft Deletes: A Safer Approach](#soft-deletes-a-safer-approach)
- [Prepared Statements: Your Defense Against SQL Injection](#prepared-statements-your-defense-against-sql-injection)
- [Building a Go CRUD Application with PostgreSQL](#building-a-go-crud-application-with-postgresql)
- [Mini Project: Blog Post Storage System](#mini-project-blog-post-storage-system)
- [Exercises](#exercises)
- [Summary](#summary)

---

## The Spreadsheet Analogy

Before we write a single line of SQL, let us build the right mental model.

Imagine a spreadsheet. You know what a spreadsheet looks like — Google Sheets, Microsoft Excel. It has rows and columns. The top row has column names (Name, Age, Email). Every row after that is one record — one person, one order, one blog post.

A database table is exactly this, but with rules. In a spreadsheet, you can type anything in any cell — a number in the Name column, a date in the Age column. Nobody stops you. A database table does not let you do this. You define the structure upfront: "column Age must be a whole number, column Name must be text, column Email cannot be empty." The database enforces these rules every single time you add or change data.

This enforcement is one of the most powerful things about relational databases. Your data stays clean and consistent because the database itself rejects anything that does not fit the rules. You do not need to remember to check — the database checks for you.

Now let us learn how to build these tables and rules.

---

## CREATE TABLE: Building the Structure

### Setting Up PostgreSQL

This chapter uses PostgreSQL. If you have not installed it yet, the quickest way is:

```bash
# On macOS with Homebrew:
brew install postgresql@16
brew services start postgresql@16

# On Ubuntu/Debian:
sudo apt install postgresql
sudo systemctl start postgresql

# Then create a database to work with:
createdb learningdb
```

To verify it is running, open a terminal and type:

```bash
psql learningdb
```

You should see the `learningdb=#` prompt. Type `\q` to exit.

### Your First CREATE TABLE

Think of `CREATE TABLE` as drawing the blueprint for a spreadsheet before you start using it. You are saying: "Here is the name of my table. Here are its columns. Here are the rules each column must follow."

The syntax looks like this:

```sql
CREATE TABLE table_name (
    column_name  data_type  constraints,
    column_name  data_type  constraints,
    ...
);
```

Let us create a table for blog posts:

```sql
CREATE TABLE posts (
    id          SERIAL PRIMARY KEY,
    title       TEXT NOT NULL,
    content     TEXT NOT NULL,
    author      TEXT NOT NULL,
    published   BOOLEAN NOT NULL DEFAULT false,
    created_at  TIMESTAMP NOT NULL DEFAULT NOW()
);
```

Let us read this line by line.

**`id SERIAL PRIMARY KEY`**

- `id` is the column name. It will hold a unique number for each row.
- `SERIAL` is a PostgreSQL data type that means "auto-incrementing integer." Every time you insert a row without providing an `id`, PostgreSQL assigns the next available number automatically. The first row gets 1, the second gets 2, and so on. You never have to manage this yourself.
- `PRIMARY KEY` is a constraint. It means two things: (1) no two rows can have the same `id`, and (2) the column cannot be empty. Every table should have a primary key. It is how you identify any specific row in the table — like a passport number that belongs to exactly one person.

**`title TEXT NOT NULL`**

- `title` is the column name.
- `TEXT` is the data type. In PostgreSQL, `TEXT` holds a string of any length.
- `NOT NULL` is a constraint. It means this column cannot be empty. If you try to insert a row without providing a title, PostgreSQL will refuse and give you an error. This prevents "ghost rows" with missing data from ending up in your database.

**`published BOOLEAN NOT NULL DEFAULT false`**

- `BOOLEAN` stores either `true` or `false` — nothing else.
- `DEFAULT false` means if you do not provide a value for `published` when inserting a row, PostgreSQL automatically uses `false`. This is useful for columns that have a sensible initial value.

**`created_at TIMESTAMP NOT NULL DEFAULT NOW()`**

- `TIMESTAMP` stores a date and time (for example: `2024-03-15 14:32:00`).
- `DEFAULT NOW()` means if you do not provide a value, PostgreSQL uses the current date and time. Your rows automatically get a creation timestamp without you doing anything special.

### Common PostgreSQL Data Types

Here is a reference table of the types you will use most often:

| Type | What It Stores | Example |
|---|---|---|
| `INTEGER` or `INT` | Whole numbers | `42`, `-7`, `0` |
| `SERIAL` | Auto-incrementing integer | Assigned automatically |
| `BIGINT` | Very large whole numbers | Used for IDs in huge tables |
| `TEXT` | Text of any length | `'Hello, world'` |
| `VARCHAR(n)` | Text up to n characters | `VARCHAR(255)` |
| `BOOLEAN` | True or false | `true`, `false` |
| `DECIMAL(p,s)` | Exact decimal numbers | `9.99`, used for money |
| `REAL` | Approximate decimal | `3.14159` |
| `TIMESTAMP` | Date and time | `2024-03-15 14:32:00` |
| `DATE` | Date only | `2024-03-15` |
| `UUID` | Globally unique identifier | `a3b8...` |

### Constraints: Rules That Protect Your Data

Constraints are rules the database enforces automatically. Here are the main ones:

| Constraint | What It Does |
|---|---|
| `PRIMARY KEY` | Unique, not null. Identifies each row. |
| `NOT NULL` | Column must always have a value. |
| `UNIQUE` | No two rows can have the same value in this column. |
| `DEFAULT value` | Used as the value when none is provided. |
| `CHECK (condition)` | Column value must satisfy this condition. |
| `REFERENCES table(col)` | Foreign key — value must exist in another table. |

Example using `UNIQUE` and `CHECK`:

```sql
CREATE TABLE users (
    id       SERIAL PRIMARY KEY,
    email    TEXT NOT NULL UNIQUE,
    age      INTEGER NOT NULL CHECK (age >= 0 AND age <= 150),
    username TEXT NOT NULL UNIQUE
);
```

The `UNIQUE` on `email` means no two users can register with the same email address. PostgreSQL will refuse the insert.

The `CHECK (age >= 0 AND age <= 150)` means nobody can insert an age of -5 or 999. The database enforces that the age is a sensible number.

### Quick Check

> 1. What does `PRIMARY KEY` mean? Why should every table have one?
> 2. What is the difference between `TEXT` and `VARCHAR(255)`?
> 3. If you create a column with `DEFAULT false`, what value does it get when you insert a row without specifying that column?

---

## INSERT INTO: Adding Data

Creating a table gives you an empty spreadsheet. `INSERT INTO` is how you add rows — how you put data in.

### Inserting a Single Row

```sql
INSERT INTO posts (title, content, author)
VALUES ('My First Post', 'This is the content of my first post.', 'Alice');
```

Let us break this down:

- `INSERT INTO posts` — we are adding a row to the `posts` table.
- `(title, content, author)` — we are providing values for these three columns. We are skipping `id`, `published`, and `created_at` because they have defaults (`SERIAL`, `DEFAULT false`, and `DEFAULT NOW()` respectively).
- `VALUES ('My First Post', ...)` — the actual values, in the same order as the column list.

After running this, PostgreSQL automatically assigns `id = 1`, `published = false`, and `created_at = <current time>`.

### Inserting Multiple Rows at Once

You do not have to insert one row at a time. You can insert many rows in a single statement:

```sql
INSERT INTO posts (title, content, author)
VALUES
    ('Getting Started with Go', 'Go is a statically typed language...', 'Bob'),
    ('Understanding Databases', 'A database is a structured...', 'Alice'),
    ('Web Servers in Go', 'The net/http package...', 'Carol');
```

This is more efficient than three separate `INSERT` statements because PostgreSQL processes all three rows as a single operation.

### Getting the Auto-Generated ID Back

Often after inserting a row, you need to know what `id` was assigned. PostgreSQL provides `RETURNING`:

```sql
INSERT INTO posts (title, content, author)
VALUES ('A New Post', 'Content here.', 'Dave')
RETURNING id;
```

PostgreSQL executes the insert and then returns the `id` of the newly created row. You will see this pattern constantly in Go code.

### Quick Check

> 1. What happens to columns with `DEFAULT` values if you do not include them in your `INSERT INTO` statement?
> 2. How would you insert three users into a `users` table in a single statement?
> 3. What does `RETURNING id` do after an `INSERT`?

---

## SELECT: Retrieving Data

`SELECT` is the most used SQL command. It reads data from a table. In everyday English, it means "give me this data."

### Selecting Everything

```sql
SELECT * FROM posts;
```

The `*` is a wildcard that means "all columns." This returns every row and every column in the `posts` table.

In production code, avoid `SELECT *`. Always name the columns you need:

```sql
SELECT id, title, author, created_at FROM posts;
```

Why? Because if someone adds a new column to the table later, `SELECT *` suddenly returns that column too, which might break your Go code that expects a specific number of columns.

### Selecting Specific Columns

```sql
SELECT title, author FROM posts;
```

This returns only the `title` and `author` columns for every row.

---

## WHERE: Filtering Rows

`WHERE` is how you narrow down results to only the rows you care about. Think of it as a filter on your spreadsheet.

### Basic Filtering

```sql
SELECT title, author FROM posts WHERE author = 'Alice';
```

This returns only the rows where the `author` column equals `'Alice'`. Every other row is ignored.

### Comparison Operators

| Operator | Meaning | Example |
|---|---|---|
| `=` | Equal to | `author = 'Alice'` |
| `!=` or `<>` | Not equal to | `published != true` |
| `>` | Greater than | `age > 18` |
| `<` | Less than | `age < 65` |
| `>=` | Greater than or equal | `score >= 90` |
| `<=` | Less than or equal | `price <= 100` |
| `LIKE` | Pattern match | `title LIKE 'Go%'` |
| `IS NULL` | Column has no value | `deleted_at IS NULL` |
| `IS NOT NULL` | Column has a value | `published_at IS NOT NULL` |

### AND: Both Conditions Must Be True

Using `AND` means a row is only returned if all conditions are true. It is like saying "I want rows where this is true AND that is also true."

```sql
SELECT title, author
FROM posts
WHERE author = 'Alice' AND published = true;
```

This returns only posts by Alice that are published. A post by Alice that is not published would not appear. A published post by Bob would not appear either.

### OR: At Least One Condition Must Be True

`OR` means a row is returned if any one of the conditions is true.

```sql
SELECT title, author
FROM posts
WHERE author = 'Alice' OR author = 'Bob';
```

This returns all posts by either Alice or Bob.

### NOT: Inverting a Condition

`NOT` reverses a condition. Instead of "rows where this is true," you get "rows where this is NOT true."

```sql
SELECT title, author
FROM posts
WHERE NOT published = false;
```

This is equivalent to `WHERE published = true`. `NOT` is most useful when combined with `LIKE` or `IN`:

```sql
-- All posts NOT by Alice
SELECT title FROM posts WHERE NOT author = 'Alice';

-- All posts whose title does NOT contain "Go"
SELECT title FROM posts WHERE title NOT LIKE '%Go%';
```

### Combining AND, OR, and NOT with Parentheses

When you mix `AND` and `OR`, use parentheses to make the logic explicit. Without parentheses, `AND` is evaluated before `OR` — just like multiplication before addition in math — and this can cause unexpected results.

```sql
-- Posts published by Alice, OR any post by Bob (regardless of published status)
-- Without parentheses, AND binds tighter than OR, so this is:
-- (author = 'Alice' AND published = true) OR (author = 'Bob')
SELECT title FROM posts WHERE author = 'Alice' AND published = true OR author = 'Bob';

-- The same query with explicit parentheses — clearer and correct:
SELECT title FROM posts WHERE (author = 'Alice' AND published = true) OR author = 'Bob';
```

Make it a habit to use parentheses whenever you combine more than one logical operator.

### LIKE: Pattern Matching

`LIKE` lets you match text against a pattern. There are two wildcard characters:
- `%` matches any sequence of zero or more characters.
- `_` matches exactly one character.

```sql
-- Titles that start with "Go"
SELECT title FROM posts WHERE title LIKE 'Go%';

-- Titles that contain the word "database" anywhere
SELECT title FROM posts WHERE title LIKE '%database%';

-- Titles where the second character is 'o'
SELECT title FROM posts WHERE title LIKE '_o%';
```

`LIKE` is case-sensitive in PostgreSQL. Use `ILIKE` if you want case-insensitive matching:

```sql
SELECT title FROM posts WHERE title ILIKE '%database%';
```

This would match "database", "Database", and "DATABASE".

### Quick Check

> 1. Write a `SELECT` that returns all posts where `published` is `true` AND `author` is `'Carol'`.
> 2. What does `%` mean in a `LIKE` pattern?
> 3. Why should you use parentheses when combining `AND` and `OR`?

---

## ORDER BY and LIMIT

### ORDER BY: Sorting Results

Imagine getting a list of 1,000 blog posts in no particular order. That is hard to work with. `ORDER BY` sorts the results.

```sql
-- Newest posts first
SELECT title, created_at FROM posts ORDER BY created_at DESC;

-- Alphabetical by title
SELECT title, author FROM posts ORDER BY title ASC;
```

- `ASC` means ascending: A to Z, smallest to largest, oldest to newest. This is the default if you do not specify.
- `DESC` means descending: Z to A, largest to smallest, newest to oldest.

You can sort by multiple columns. PostgreSQL sorts by the first column, and then uses the second column to break ties:

```sql
-- Sort by author A-Z, then within each author by newest post first
SELECT title, author, created_at
FROM posts
ORDER BY author ASC, created_at DESC;
```

### LIMIT: Restricting the Number of Results

`LIMIT` caps how many rows are returned.

```sql
-- Return only the 5 most recent posts
SELECT title, created_at FROM posts ORDER BY created_at DESC LIMIT 5;
```

This is extremely important for real applications. Imagine a posts table with 10 million rows. Without `LIMIT`, `SELECT * FROM posts` would try to send all 10 million rows to your Go program. Your server would run out of memory. Always use `LIMIT` when you do not need every row.

### OFFSET: Skipping Rows (Pagination)

`OFFSET` tells PostgreSQL to skip the first N rows before starting to return results. Combined with `LIMIT`, this enables pagination — showing results in pages.

```sql
-- Page 1: posts 1-10
SELECT title FROM posts ORDER BY created_at DESC LIMIT 10 OFFSET 0;

-- Page 2: posts 11-20
SELECT title FROM posts ORDER BY created_at DESC LIMIT 10 OFFSET 10;

-- Page 3: posts 21-30
SELECT title FROM posts ORDER BY created_at DESC LIMIT 10 OFFSET 20;
```

The pattern: for page number `p` with `n` items per page, the offset is `(p - 1) * n`.

Note: `OFFSET` pagination has performance problems on very large tables because PostgreSQL must still scan and discard all the skipped rows. For large datasets, cursor-based pagination (using `WHERE id > last_seen_id`) is faster, but for most applications, `OFFSET` is fine.

### Quick Check

> 1. Write a query that returns the 3 most recently created posts.
> 2. What is the `OFFSET` for page 4 of results when displaying 10 items per page?
> 3. What is the difference between `ASC` and `DESC` in `ORDER BY`?

---

## UPDATE: Changing Existing Data

`UPDATE` modifies rows that already exist in the table.

### Basic UPDATE

```sql
UPDATE posts
SET published = true
WHERE id = 1;
```

This finds the row where `id = 1` and changes `published` to `true`. The `WHERE` clause is crucial — without it, you update every row in the table, which is almost never what you want.

### Updating Multiple Columns at Once

```sql
UPDATE posts
SET title = 'New Title Here', content = 'Updated content.'
WHERE id = 3;
```

Separate multiple column updates with commas inside the `SET` clause.

### Updating Multiple Rows at Once

`WHERE` is a filter, so any rows matching the condition will be updated:

```sql
-- Publish all posts by Alice
UPDATE posts
SET published = true
WHERE author = 'Alice';
```

### Getting Updated Data Back with RETURNING

Like `INSERT`, `UPDATE` can return data from the modified rows:

```sql
UPDATE posts
SET published = true
WHERE id = 5
RETURNING id, title, published;
```

This is useful when your Go code needs to confirm what changed.

### The Golden Rule of UPDATE

**Always include a `WHERE` clause unless you truly mean to update every single row.** This mistake is easy to make and can be catastrophic:

```sql
-- DANGEROUS: This updates EVERY post in the table
UPDATE posts SET published = true;

-- SAFE: This updates only the post with id 5
UPDATE posts SET published = true WHERE id = 5;
```

---

## DELETE: Removing Data

`DELETE` removes rows from a table permanently.

```sql
DELETE FROM posts WHERE id = 7;
```

This removes the row where `id = 7`. It is gone. You cannot undo it (unless you are inside a transaction you have not committed yet).

### The Same Golden Rule Applies

```sql
-- CATASTROPHIC: This deletes every row in the entire table
DELETE FROM posts;

-- SAFE: This deletes only the post with id 7
DELETE FROM posts WHERE id = 7;
```

### Why You Should Almost Never Delete

Here is something experienced engineers learn over time: deleting data is usually a mistake.

Consider these scenarios:

1. A user deletes their account. Six months later, they call customer support asking about a charge. Their account is gone. You cannot investigate.
2. You delete old "inactive" posts. A journalist asks about a post that was published three years ago. It is gone.
3. You delete a row that other rows reference via a foreign key. The database either refuses the delete (if you have constraints set up correctly) or allows it and leaves orphaned data.
4. A bug in your code accidentally deletes the wrong rows. You have no way to recover them.

Real production systems almost never use `DELETE`. Instead, they use a technique called soft deletes.

---

## Soft Deletes: A Safer Approach

### The Idea

Instead of removing a row, you mark it as deleted. You add a column called `deleted_at` (a timestamp). When you want to "delete" a row, you set `deleted_at` to the current time. When you query for active rows, you filter to only rows where `deleted_at IS NULL`.

The row still exists in the database. You can recover it. You can audit it. You can answer questions about what existed in the past.

### Adding deleted_at to the Posts Table

```sql
ALTER TABLE posts ADD COLUMN deleted_at TIMESTAMP;
```

`ALTER TABLE` modifies an existing table. Here we are adding a new column. Since we did not specify `NOT NULL`, the column defaults to `NULL` for all existing rows and all new rows — meaning nothing is deleted yet.

### Soft Deleting a Row

Instead of `DELETE FROM posts WHERE id = 7`, you do:

```sql
UPDATE posts
SET deleted_at = NOW()
WHERE id = 7;
```

The row is still in the table. `deleted_at` now has a timestamp showing when it was "deleted."

### Querying Only Active (Non-Deleted) Rows

Every query that should only see active posts must include `WHERE deleted_at IS NULL`:

```sql
SELECT id, title, author FROM posts WHERE deleted_at IS NULL;

-- Only published, non-deleted posts by Alice:
SELECT title
FROM posts
WHERE author = 'Alice'
  AND published = true
  AND deleted_at IS NULL;
```

### Recovering a Soft-Deleted Row

Recovering a soft-deleted row is just another `UPDATE`:

```sql
UPDATE posts
SET deleted_at = NULL
WHERE id = 7;
```

The post is "undeleted" instantly. Compare this to trying to recover a row that was hard-deleted — that requires a database backup restore, which is slow, expensive, and might lose other recent data.

### Quick Check

> 1. What is the difference between hard delete and soft delete?
> 2. How do you query for only non-deleted rows when using soft deletes?
> 3. Name one real-world scenario where soft deletes would have prevented a serious problem.

---

## Prepared Statements: Your Defense Against SQL Injection

### What Is SQL Injection?

SQL injection is one of the most common and dangerous security vulnerabilities in software. Let us see what it looks like.

Imagine you have a Go web server that takes a username from a form and queries the database:

```go
// DANGEROUS — DO NOT DO THIS
username := r.FormValue("username")
query := "SELECT * FROM users WHERE username = '" + username + "'"
rows, err := db.Query(query)
```

This looks reasonable. If `username` is `"alice"`, the query becomes:

```sql
SELECT * FROM users WHERE username = 'alice'
```

But what if someone submits this as their username?

```
' OR '1'='1
```

Your query becomes:

```sql
SELECT * FROM users WHERE username = '' OR '1'='1'
```

The condition `'1'='1'` is always true. This query now returns every user in the database. The attacker just bypassed your authentication.

It gets worse. What if they submit:

```
'; DROP TABLE users; --
```

Your query becomes:

```sql
SELECT * FROM users WHERE username = ''; DROP TABLE users; --'
```

Your entire users table is deleted. This is SQL injection.

### The Solution: Prepared Statements

A prepared statement separates the SQL structure from the data. You write the query with placeholders, and you pass the data separately. The database driver ensures the data can never be interpreted as SQL.

In Go's `database/sql` package and the `pgx` PostgreSQL driver, placeholders are written as `$1`, `$2`, `$3`, etc.:

```go
// SAFE — always do this
username := r.FormValue("username")
rows, err := db.Query("SELECT * FROM users WHERE username = $1", username)
```

Even if the attacker submits `' OR '1'='1`, the database treats it as a literal string — a username that literally equals `' OR '1'='1'` — not as SQL code. The injection is harmless.

**The rule is simple: never build SQL queries by concatenating strings. Always use placeholders.**

---

## Building a Go CRUD Application with PostgreSQL

Now let us bring everything together. We will build a complete Go application that connects to PostgreSQL and performs all four CRUD operations: Create, Read, Update, Delete (soft delete).

CRUD stands for:
- **C**reate — INSERT new rows
- **R**ead — SELECT existing rows
- **U**pdate — UPDATE existing rows
- **D**elete — DELETE or soft-delete rows

### Setting Up the Go Project

```bash
mkdir blog-app
cd blog-app
go mod init blog-app
go get github.com/jackc/pgx/v5
```

We use `pgx/v5` — the most popular and well-maintained PostgreSQL driver for Go.

### The Complete Application

Here is the complete, runnable Go application. Read each comment carefully.

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

// Post represents one row from the posts table.
// Each field corresponds to a column.
// The `db` struct tags would be used by ORMs, but pgx uses positional Scan,
// so we just need the field names and types to match the columns we select.
type Post struct {
	ID        int
	Title     string
	Content   string
	Author    string
	Published bool
	CreatedAt time.Time
	DeletedAt *time.Time // pointer because this can be NULL
}

func main() {
	// context.Background() is a standard Go context — think of it as a background
	// task holder. pgx requires a context for all operations so that you can cancel
	// long-running queries if needed.
	ctx := context.Background()

	// Read the database connection string from an environment variable.
	// A connection string tells pgx where PostgreSQL is and how to log in.
	// Format: postgres://username:password@host:port/database
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		// Default for local development — adjust for your setup.
		connStr = "postgres://postgres:postgres@localhost:5432/learningdb"
	}

	// pgx.Connect opens a connection to PostgreSQL.
	// If the connection fails (wrong password, PostgreSQL not running, etc.),
	// err will be non-nil and we log.Fatal to exit the program with an error message.
	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Unable to connect to database: %v\n", err)
	}
	// defer conn.Close runs when main() returns, closing the connection cleanly.
	defer conn.Close(ctx)

	fmt.Println("Connected to PostgreSQL successfully.")

	// --- Step 1: Create the table ---
	// We use CREATE TABLE IF NOT EXISTS so this is safe to run multiple times.
	// Without IF NOT EXISTS, running this twice would give an error
	// because the table already exists.
	_, err = conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id          SERIAL PRIMARY KEY,
			title       TEXT NOT NULL,
			content     TEXT NOT NULL,
			author      TEXT NOT NULL,
			published   BOOLEAN NOT NULL DEFAULT false,
			created_at  TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at  TIMESTAMP
		)
	`)
	if err != nil {
		log.Fatalf("Failed to create table: %v\n", err)
	}
	fmt.Println("Table 'posts' is ready.")

	// --- Step 2: Insert some posts ---
	// We call a helper function to keep main() clean.
	id1, err := createPost(ctx, conn, "Getting Started with Go", "Go is a statically typed, compiled language...", "Alice")
	if err != nil {
		log.Fatalf("Failed to create post: %v\n", err)
	}
	fmt.Printf("Created post with ID: %d\n", id1)

	id2, err := createPost(ctx, conn, "Understanding PostgreSQL", "PostgreSQL is a powerful open-source database...", "Bob")
	if err != nil {
		log.Fatalf("Failed to create post: %v\n", err)
	}
	fmt.Printf("Created post with ID: %d\n", id2)

	id3, err := createPost(ctx, conn, "Web Servers in Go", "The net/http package makes building servers easy...", "Alice")
	if err != nil {
		log.Fatalf("Failed to create post: %v\n", err)
	}
	fmt.Printf("Created post with ID: %d\n", id3)

	// --- Step 3: Read all posts ---
	fmt.Println("\n--- All Posts ---")
	posts, err := getAllPosts(ctx, conn)
	if err != nil {
		log.Fatalf("Failed to get posts: %v\n", err)
	}
	for _, p := range posts {
		fmt.Printf("[%d] %s by %s (published: %v)\n", p.ID, p.Title, p.Author, p.Published)
	}

	// --- Step 4: Read posts by a specific author ---
	fmt.Println("\n--- Posts by Alice ---")
	alicePosts, err := getPostsByAuthor(ctx, conn, "Alice")
	if err != nil {
		log.Fatalf("Failed to get posts by author: %v\n", err)
	}
	for _, p := range alicePosts {
		fmt.Printf("[%d] %s\n", p.ID, p.Title)
	}

	// --- Step 5: Publish the first post ---
	err = publishPost(ctx, conn, id1)
	if err != nil {
		log.Fatalf("Failed to publish post: %v\n", err)
	}
	fmt.Printf("\nPublished post %d\n", id1)

	// --- Step 6: Soft delete the second post ---
	err = softDeletePost(ctx, conn, id2)
	if err != nil {
		log.Fatalf("Failed to soft-delete post: %v\n", err)
	}
	fmt.Printf("Soft-deleted post %d\n", id2)

	// --- Step 7: Read only active (non-deleted) posts ---
	fmt.Println("\n--- Active Posts (after soft delete) ---")
	activePosts, err := getActivePosts(ctx, conn)
	if err != nil {
		log.Fatalf("Failed to get active posts: %v\n", err)
	}
	for _, p := range activePosts {
		fmt.Printf("[%d] %s by %s (published: %v)\n", p.ID, p.Title, p.Author, p.Published)
	}
}

// createPost inserts a new post and returns the auto-assigned ID.
// Note the $1, $2, $3 placeholders — these prevent SQL injection.
// The actual values (title, content, author) are passed as separate arguments.
func createPost(ctx context.Context, conn *pgx.Conn, title, content, author string) (int, error) {
	// QueryRow executes a query that returns exactly one row.
	// RETURNING id tells PostgreSQL to give us back the auto-generated ID.
	var id int
	err := conn.QueryRow(ctx,
		"INSERT INTO posts (title, content, author) VALUES ($1, $2, $3) RETURNING id",
		title, content, author,
	).Scan(&id) // Scan reads the returned value into our id variable
	if err != nil {
		return 0, fmt.Errorf("createPost: %w", err)
	}
	return id, nil
}

// getAllPosts returns all posts that have not been soft-deleted,
// ordered by newest first.
func getAllPosts(ctx context.Context, conn *pgx.Conn) ([]Post, error) {
	// conn.Query executes a query that returns multiple rows.
	rows, err := conn.Query(ctx,
		"SELECT id, title, content, author, published, created_at FROM posts WHERE deleted_at IS NULL ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("getAllPosts: %w", err)
	}
	// defer rows.Close ensures we release resources when we are done.
	// If you forget this, the connection stays busy and you get errors.
	defer rows.Close()

	// We will collect results into this slice.
	var posts []Post

	// rows.Next() advances to the next row and returns true if there is one.
	// When there are no more rows, it returns false and the loop ends.
	for rows.Next() {
		var p Post
		// rows.Scan reads the current row's column values into Go variables.
		// The order must match exactly the order of columns in the SELECT.
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Published, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("getAllPosts scan: %w", err)
		}
		posts = append(posts, p)
	}

	// rows.Err() checks if the loop stopped due to an error (rather than
	// running out of rows normally). Always check this after the loop.
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("getAllPosts rows error: %w", err)
	}

	return posts, nil
}

// getPostsByAuthor returns all non-deleted posts by a specific author.
// $1 is the placeholder for the author argument.
func getPostsByAuthor(ctx context.Context, conn *pgx.Conn, author string) ([]Post, error) {
	rows, err := conn.Query(ctx,
		"SELECT id, title, content, author, published, created_at FROM posts WHERE author = $1 AND deleted_at IS NULL ORDER BY created_at DESC",
		author, // This value replaces $1 safely — no SQL injection possible
	)
	if err != nil {
		return nil, fmt.Errorf("getPostsByAuthor: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Published, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("getPostsByAuthor scan: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// publishPost sets published = true for a given post ID.
func publishPost(ctx context.Context, conn *pgx.Conn, id int) error {
	// conn.Exec is used for queries that do not return rows (UPDATE, DELETE, etc.)
	// commandTag tells us how many rows were affected.
	commandTag, err := conn.Exec(ctx,
		"UPDATE posts SET published = true WHERE id = $1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return fmt.Errorf("publishPost: %w", err)
	}
	// RowsAffected() returns 0 if no row with this id was found.
	// This helps catch bugs where you try to update a non-existent row.
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("publishPost: no post found with id %d", id)
	}
	return nil
}

// softDeletePost marks a post as deleted by setting deleted_at to the current time.
// The actual row is never removed from the database.
func softDeletePost(ctx context.Context, conn *pgx.Conn, id int) error {
	commandTag, err := conn.Exec(ctx,
		"UPDATE posts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL",
		id,
	)
	if err != nil {
		return fmt.Errorf("softDeletePost: %w", err)
	}
	if commandTag.RowsAffected() == 0 {
		return fmt.Errorf("softDeletePost: no post found with id %d", id)
	}
	return nil
}

// getActivePosts returns all posts where deleted_at IS NULL.
// "Active" means not soft-deleted.
func getActivePosts(ctx context.Context, conn *pgx.Conn) ([]Post, error) {
	rows, err := conn.Query(ctx,
		"SELECT id, title, content, author, published, created_at FROM posts WHERE deleted_at IS NULL ORDER BY created_at DESC",
	)
	if err != nil {
		return nil, fmt.Errorf("getActivePosts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		err := rows.Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Published, &p.CreatedAt)
		if err != nil {
			return nil, fmt.Errorf("getActivePosts scan: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
```

To run this:

```bash
export DATABASE_URL="postgres://postgres:postgres@localhost:5432/learningdb"
go run main.go
```

You should see output like:

```
Connected to PostgreSQL successfully.
Table 'posts' is ready.
Created post with ID: 1
Created post with ID: 2
Created post with ID: 3

--- All Posts ---
[3] Web Servers in Go by Alice (published: false)
[2] Understanding PostgreSQL by Bob (published: false)
[1] Getting Started with Go by Alice (published: false)

--- Posts by Alice ---
[3] Web Servers in Go
[1] Getting Started with Go

Published post 1
Soft-deleted post 2

--- Active Posts (after soft delete) ---
[3] Web Servers in Go by Alice (published: false)
[1] Getting Started with Go by Alice (published: true)
```

Post 2 (Bob's PostgreSQL post) no longer appears because it was soft-deleted.

---

## Mini Project: Blog Post Storage System

Now let us build something more complete — a command-line blog system that lets you add posts, list them, publish them, and soft-delete them, all from your terminal.

### Project Structure

```
blog-cli/
├── go.mod
├── main.go
└── db/
    └── posts.go
```

```bash
mkdir blog-cli
cd blog-cli
go mod init blog-cli
go get github.com/jackc/pgx/v5
mkdir db
```

### db/posts.go — The Data Layer

This file contains all database logic. It is completely separate from the user interface code in `main.go`. This separation is called "layered architecture" and it makes your code easier to test and change.

```go
package db

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5"
)

// Post represents a single blog post.
type Post struct {
	ID        int
	Title     string
	Content   string
	Author    string
	Published bool
	CreatedAt time.Time
	DeletedAt *time.Time
}

// SetupSchema creates the posts table if it does not already exist.
// Safe to call every time the application starts.
func SetupSchema(ctx context.Context, conn *pgx.Conn) error {
	_, err := conn.Exec(ctx, `
		CREATE TABLE IF NOT EXISTS posts (
			id         SERIAL PRIMARY KEY,
			title      TEXT NOT NULL,
			content    TEXT NOT NULL,
			author     TEXT NOT NULL,
			published  BOOLEAN NOT NULL DEFAULT false,
			created_at TIMESTAMP NOT NULL DEFAULT NOW(),
			deleted_at TIMESTAMP
		)
	`)
	if err != nil {
		return fmt.Errorf("SetupSchema: %w", err)
	}
	return nil
}

// CreatePost inserts a new post and returns its assigned ID.
func CreatePost(ctx context.Context, conn *pgx.Conn, title, content, author string) (int, error) {
	var id int
	err := conn.QueryRow(ctx,
		`INSERT INTO posts (title, content, author)
		 VALUES ($1, $2, $3)
		 RETURNING id`,
		title, content, author,
	).Scan(&id)
	if err != nil {
		return 0, fmt.Errorf("CreatePost: %w", err)
	}
	return id, nil
}

// ListActivePosts returns all non-deleted posts, newest first.
func ListActivePosts(ctx context.Context, conn *pgx.Conn) ([]Post, error) {
	rows, err := conn.Query(ctx,
		`SELECT id, title, author, published, created_at
		 FROM posts
		 WHERE deleted_at IS NULL
		 ORDER BY created_at DESC`,
	)
	if err != nil {
		return nil, fmt.Errorf("ListActivePosts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		// We only select 5 columns here, so we only scan 5 fields.
		// The order must match the SELECT column order exactly.
		if err := rows.Scan(&p.ID, &p.Title, &p.Author, &p.Published, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("ListActivePosts scan: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}

// GetPost retrieves a single post by its ID.
// Returns an error if the post does not exist or has been soft-deleted.
func GetPost(ctx context.Context, conn *pgx.Conn, id int) (Post, error) {
	var p Post
	err := conn.QueryRow(ctx,
		`SELECT id, title, content, author, published, created_at
		 FROM posts
		 WHERE id = $1 AND deleted_at IS NULL`,
		id,
	).Scan(&p.ID, &p.Title, &p.Content, &p.Author, &p.Published, &p.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			// pgx.ErrNoRows is returned when QueryRow finds no matching row.
			// We return a clear error message rather than a generic database error.
			return Post{}, fmt.Errorf("GetPost: post with id %d not found", id)
		}
		return Post{}, fmt.Errorf("GetPost: %w", err)
	}
	return p, nil
}

// PublishPost sets published = true for a given post.
func PublishPost(ctx context.Context, conn *pgx.Conn, id int) error {
	tag, err := conn.Exec(ctx,
		`UPDATE posts SET published = true WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("PublishPost: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("PublishPost: post %d not found or already deleted", id)
	}
	return nil
}

// UnpublishPost sets published = false for a given post.
func UnpublishPost(ctx context.Context, conn *pgx.Conn, id int) error {
	tag, err := conn.Exec(ctx,
		`UPDATE posts SET published = false WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("UnpublishPost: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("UnpublishPost: post %d not found or already deleted", id)
	}
	return nil
}

// SoftDelete marks a post as deleted without removing it from the database.
func SoftDelete(ctx context.Context, conn *pgx.Conn, id int) error {
	tag, err := conn.Exec(ctx,
		`UPDATE posts SET deleted_at = NOW() WHERE id = $1 AND deleted_at IS NULL`,
		id,
	)
	if err != nil {
		return fmt.Errorf("SoftDelete: %w", err)
	}
	if tag.RowsAffected() == 0 {
		return fmt.Errorf("SoftDelete: post %d not found or already deleted", id)
	}
	return nil
}

// SearchPosts returns non-deleted posts where the title contains the search term.
// ILIKE is case-insensitive LIKE in PostgreSQL.
// The %% around the term are literal percent signs (% must be doubled in fmt.Sprintf
// to avoid being treated as a format verb, but here we pass it as a query argument
// so we just use the string directly).
func SearchPosts(ctx context.Context, conn *pgx.Conn, term string) ([]Post, error) {
	// We wrap the term with % wildcards to do a "contains" search.
	// This is safe from SQL injection because it goes through a placeholder.
	searchPattern := "%" + term + "%"

	rows, err := conn.Query(ctx,
		`SELECT id, title, author, published, created_at
		 FROM posts
		 WHERE title ILIKE $1 AND deleted_at IS NULL
		 ORDER BY created_at DESC`,
		searchPattern,
	)
	if err != nil {
		return nil, fmt.Errorf("SearchPosts: %w", err)
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Title, &p.Author, &p.Published, &p.CreatedAt); err != nil {
			return nil, fmt.Errorf("SearchPosts scan: %w", err)
		}
		posts = append(posts, p)
	}
	return posts, rows.Err()
}
```

### main.go — The CLI Interface

```go
package main

import (
	"context"
	"fmt"
	"log"
	"os"

	"blog-cli/db"

	"github.com/jackc/pgx/v5"
)

func main() {
	ctx := context.Background()

	// Get connection string from environment, or use a local default.
	connStr := os.Getenv("DATABASE_URL")
	if connStr == "" {
		connStr = "postgres://postgres:postgres@localhost:5432/learningdb"
	}

	conn, err := pgx.Connect(ctx, connStr)
	if err != nil {
		log.Fatalf("Cannot connect to database: %v\n", err)
	}
	defer conn.Close(ctx)

	// Ensure the table exists before doing anything.
	if err := db.SetupSchema(ctx, conn); err != nil {
		log.Fatalf("Schema setup failed: %v\n", err)
	}

	// The program expects a command as the first argument.
	// Usage examples:
	//   go run main.go add "Title Here" "Content here" "Author Name"
	//   go run main.go list
	//   go run main.go get 3
	//   go run main.go publish 3
	//   go run main.go delete 3
	//   go run main.go search "Go"
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	command := os.Args[1]

	switch command {

	case "add":
		// Requires: title, content, author
		if len(os.Args) < 5 {
			fmt.Println("Usage: add <title> <content> <author>")
			os.Exit(1)
		}
		title := os.Args[2]
		content := os.Args[3]
		author := os.Args[4]

		id, err := db.CreatePost(ctx, conn, title, content, author)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		fmt.Printf("Created post with ID %d\n", id)

	case "list":
		posts, err := db.ListActivePosts(ctx, conn)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		if len(posts) == 0 {
			fmt.Println("No posts found.")
			return
		}
		fmt.Printf("%-5s  %-35s  %-15s  %-10s  %s\n", "ID", "Title", "Author", "Published", "Created")
		fmt.Println("-----  -----------------------------------  ---------------  ----------  -------------------")
		for _, p := range posts {
			published := "no"
			if p.Published {
				published = "yes"
			}
			fmt.Printf("%-5d  %-35s  %-15s  %-10s  %s\n",
				p.ID,
				truncate(p.Title, 35),
				truncate(p.Author, 15),
				published,
				p.CreatedAt.Format("2006-01-02 15:04"),
			)
		}

	case "get":
		if len(os.Args) < 3 {
			fmt.Println("Usage: get <id>")
			os.Exit(1)
		}
		var id int
		fmt.Sscanf(os.Args[2], "%d", &id)

		post, err := db.GetPost(ctx, conn, id)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		published := "Draft"
		if post.Published {
			published = "Published"
		}
		fmt.Printf("ID:      %d\n", post.ID)
		fmt.Printf("Title:   %s\n", post.Title)
		fmt.Printf("Author:  %s\n", post.Author)
		fmt.Printf("Status:  %s\n", published)
		fmt.Printf("Created: %s\n", post.CreatedAt.Format("2006-01-02 15:04:05"))
		fmt.Printf("\n%s\n", post.Content)

	case "publish":
		if len(os.Args) < 3 {
			fmt.Println("Usage: publish <id>")
			os.Exit(1)
		}
		var id int
		fmt.Sscanf(os.Args[2], "%d", &id)

		if err := db.PublishPost(ctx, conn, id); err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		fmt.Printf("Post %d is now published.\n", id)

	case "delete":
		if len(os.Args) < 3 {
			fmt.Println("Usage: delete <id>")
			os.Exit(1)
		}
		var id int
		fmt.Sscanf(os.Args[2], "%d", &id)

		if err := db.SoftDelete(ctx, conn, id); err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		fmt.Printf("Post %d has been soft-deleted. It can be recovered.\n", id)

	case "search":
		if len(os.Args) < 3 {
			fmt.Println("Usage: search <term>")
			os.Exit(1)
		}
		term := os.Args[2]

		posts, err := db.SearchPosts(ctx, conn, term)
		if err != nil {
			log.Fatalf("Error: %v\n", err)
		}
		if len(posts) == 0 {
			fmt.Printf("No posts found matching '%s'.\n", term)
			return
		}
		fmt.Printf("Results for '%s':\n", term)
		for _, p := range posts {
			fmt.Printf("  [%d] %s by %s\n", p.ID, p.Title, p.Author)
		}

	default:
		fmt.Printf("Unknown command: %s\n", command)
		printUsage()
		os.Exit(1)
	}
}

// truncate shortens a string to maxLen characters for clean table display.
func truncate(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	return s[:maxLen-3] + "..."
}

func printUsage() {
	fmt.Println("Blog CLI — a simple blog post manager")
	fmt.Println()
	fmt.Println("Commands:")
	fmt.Println("  add <title> <content> <author>  Create a new post")
	fmt.Println("  list                            List all active posts")
	fmt.Println("  get <id>                        Show a single post")
	fmt.Println("  publish <id>                    Publish a post")
	fmt.Println("  delete <id>                     Soft-delete a post")
	fmt.Println("  search <term>                   Search posts by title")
}
```

### Using the Blog CLI

```bash
# Add some posts
go run main.go add "Hello World" "My first post content." "Alice"
go run main.go add "Learning Go" "Go is a fantastic language for systems programming." "Bob"
go run main.go add "PostgreSQL Tips" "Here are my top PostgreSQL tips." "Alice"

# List all posts
go run main.go list

# See a single post's full content
go run main.go get 2

# Publish post 1
go run main.go publish 1

# Search for posts about Go
go run main.go search "Go"

# Soft-delete post 2
go run main.go delete 2

# List again — post 2 is gone from the list
go run main.go list

# But the data is still in the database — verify with psql:
# psql learningdb -c "SELECT id, title, deleted_at FROM posts;"
```

This is a complete, working blog post storage system. It uses every concept from this chapter: `CREATE TABLE`, `INSERT`, `SELECT`, `WHERE`, `ORDER BY`, `UPDATE`, soft deletes, and prepared statements throughout.

---

## Exercises

### Easy

1. Add a `views` column to the `posts` table that stores an integer (how many times the post was viewed). Give it a `DEFAULT 0`.

2. Write a single SQL statement that returns all posts ordered by `title` alphabetically (A to Z).

3. Write a `SELECT` that returns only posts where `published = true` AND the `author` is either `'Alice'` OR `'Bob'`.

4. In the blog CLI project, add a `go run main.go unpublish <id>` command that calls the existing `UnpublishPost` function.

### Medium

5. Write a SQL `UPDATE` that increments the `views` column by 1 for a specific post. Hint: you can use `SET views = views + 1`.

6. Add a `getPublishedPosts` function to `db/posts.go` that returns only posts where `published = true` and `deleted_at IS NULL`, with the most-viewed posts first.

7. Write a SQL query that returns the author name and the number of posts they have written, for all non-deleted posts. Order by post count descending. Hint: this uses `GROUP BY` and `COUNT(*)`, which you learned in Chapter 10.

8. Modify the `SearchPosts` function in `db/posts.go` to accept a `limit` and `offset` integer parameter for pagination, and include those in the SQL query.

### Hard

9. Add a `tags` feature to the blog system. Create a new table called `post_tags` with columns `post_id` (integer, references `posts`) and `tag` (text). Write Go functions to: (a) add a tag to a post, (b) remove a tag from a post, (c) list all posts with a given tag.

10. Implement a `recoverPost` command in the CLI that takes a post ID and sets `deleted_at` back to `NULL`, making a soft-deleted post active again. Add the database function in `db/posts.go` and wire it up in `main.go`.

11. The current `ListActivePosts` function always returns every post. Add an optional `author` filter: if an `author` string is provided, filter by that author; if it is an empty string, return all authors. Do this without building SQL strings by concatenation — hint: think about how you can use `$1` together with a condition like `($1 = '' OR author = $1)`.

---

## Summary

- `CREATE TABLE` defines the structure of a table, including column names, data types, and constraints like `PRIMARY KEY`, `NOT NULL`, `UNIQUE`, and `DEFAULT`.
- `INSERT INTO` adds new rows to a table. Multiple rows can be inserted in a single statement. `RETURNING` gives back auto-generated values like IDs.
- `SELECT` retrieves data, and `WHERE` filters rows using conditions. `AND` requires all conditions to be true. `OR` requires at least one. `NOT` inverts a condition. Parentheses control the order of evaluation.
- `ORDER BY` sorts results (`ASC` or `DESC`). `LIMIT` caps the number of rows returned. `OFFSET` skips rows for pagination.
- `UPDATE` changes existing rows. Always include a `WHERE` clause to avoid updating every row in the table.
- Hard `DELETE` is permanent and dangerous in production. Soft deletes — adding a `deleted_at` column and filtering with `WHERE deleted_at IS NULL` — preserve data and allow recovery.
- Prepared statements (using `$1`, `$2` placeholders in Go's `pgx` driver) are mandatory. Never build SQL by string concatenation. SQL injection can destroy your database or expose all your data.
