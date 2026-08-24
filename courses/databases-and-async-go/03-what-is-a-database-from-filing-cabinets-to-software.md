# Chapter 03: What Is a Database? From Filing Cabinets to Software

Before computers existed, companies hired entire rooms full of people to manage paper records. Banks had clerks who spent their whole careers doing nothing but filing, finding, and updating paper cards. Today, a single database can do all of that work — for millions of records — in milliseconds. Understanding what a database actually is will change how you think about every app you have ever used.

## Table of Contents

1. [What Is a Database, Really?](#what-is-a-database-really)
2. [The Database Management System (DBMS)](#the-database-management-system-dbms)
3. [Why Databases Beat Plain Files](#why-databases-beat-plain-files)
4. [Types of Databases](#types-of-databases)
5. [The Client-Server Model](#the-client-server-model)
6. [Overview: Databases We Will Cover in This Course](#overview-databases-we-will-cover-in-this-course)
7. [Real-World: Which Databases Power Netflix, Instagram, and Twitter?](#real-world-which-databases-power-netflix-instagram-and-twitter)
8. [Your First Database Interaction: SQLite in Go](#your-first-database-interaction-sqlite-in-go)
9. [Exercises](#exercises)
10. [Summary](#summary)

---

## What Is a Database, Really?

Imagine a library. Not a pile of books thrown in a room — a real library with shelves, labels, a card catalog (or today, a computer system), and a librarian who knows exactly where everything is. You can walk in and say "I want every science-fiction book published after 2010 with more than 400 pages," and within minutes you have a stack of books.

A database is exactly that — for data.

More precisely: **a database is an organized collection of data that is stored so it can be quickly retrieved, updated, and managed.** The word "organized" is the key part. A random pile of papers is not a database. A filing cabinet with labeled folders, sorted alphabetically, with an index at the front — that is a database. The computer version just does everything faster and at a much larger scale.

Here are some everyday databases you interact with every day without thinking about it:

- When you search on Google, a database finds matching web pages.
- When you log into Instagram, a database checks if your username and password match.
- When Netflix recommends a show, a database is storing what you have watched and what people with similar taste enjoyed.
- When you buy something on Amazon, a database records your order, updates the inventory count, and stores your shipping address.

Every single one of those actions happens in milliseconds, across billions of records. That is the power of a well-designed database.

### The Two Parts of "Database"

People often use the word "database" loosely to mean two different things:

1. **The data itself** — the actual records, tables, and information being stored.
2. **The software that manages the data** — the program that lets you store, find, and change that data.

The software part has its own name: the **Database Management System**, or DBMS. We will cover that in the next section.

### Quick Check

> 1. In your own words, what makes something a "database" rather than just a pile of files?
> 2. Name two apps on your phone that almost certainly use a database behind the scenes.
> 3. What is the difference between "the database" (the data) and a "DBMS"?

---

## The Database Management System (DBMS)

Think of a DBMS like the librarian, not the library itself.

The books (data) are the library. The librarian is the DBMS — the software that knows where everything is, lets you search for things, controls who is allowed to see what, makes sure two people do not accidentally mess up the same record at the same time, and keeps backups in case something goes wrong.

A **Database Management System (DBMS)** is the software that:

- Stores data on disk or in memory in an organized format
- Lets you add, read, update, and delete data
- Handles many users reading and writing at the same time
- Protects data from corruption
- Controls who has permission to access what
- Provides a language (usually SQL) for talking to the data

When developers say "we use PostgreSQL" or "we use MongoDB," they are naming their DBMS — the software managing their data.

### The DBMS as a Middleman

Your Go application never touches the raw data files directly. Instead, it talks to the DBMS, and the DBMS handles everything. This is important for several reasons:

```
Your Go App  -->  DBMS (e.g., PostgreSQL)  -->  Data files on disk
```

The DBMS acts as a gatekeeper. It ensures that when your app asks "give me all users named Alice," it gets a correct, consistent answer — even if 10,000 other requests are happening at the same moment.

### Quick Check

> 1. What does DBMS stand for?
> 2. Why does your application talk to the DBMS instead of directly reading files?
> 3. Name one thing a DBMS does that a simple folder of files cannot do.

---

## Why Databases Beat Plain Files

You might wonder: why not just store data in plain text files? A file with one user per line would work, right?

Let's think through why that breaks down quickly.

### The Plain File Problem

Imagine you run a small school and store student records in a text file like this:

```
Alice,17,Grade 11,Math:90,Science:85
Bob,16,Grade 10,Math:78,Science:92
Charlie,17,Grade 11,Math:95,Science:88
```

This works fine when you have 30 students. But what happens when:

- You have 30,000 students and need to find everyone in Grade 11 with a Math score above 80? Your program has to read every single line.
- Two teachers try to update the same student record at the same moment? One of their changes will be lost.
- Your computer crashes halfway through writing a new student record? Now you have a half-written, corrupted file.
- You want to ask "which students are in both the chess club AND the debate team?" — those are in separate files. Linking them manually is painful.

Databases were invented to solve exactly these problems. The four core guarantees databases provide are summarized by the acronym **ACID**.

### ACID: The Four Guarantees

**ACID** stands for Atomicity, Consistency, Isolation, and Durability. These are promises that a serious database makes to you about your data. Do not worry if they sound abstract — each one maps to a real problem.

#### Atomicity: All or Nothing

Imagine you are transferring $100 from your bank account to a friend. This involves two steps:
1. Subtract $100 from your account.
2. Add $100 to your friend's account.

What if the computer crashes after step 1 but before step 2? The $100 disappeared from your account but never arrived. That is a disaster.

**Atomicity** means a group of operations either ALL succeed or ALL fail together. There is no "halfway done." In database terms, this group of operations is called a **transaction**.

#### Consistency: Rules Are Always Enforced

A database lets you define rules. For example: "every student must have a name" or "an order cannot have a negative quantity." Consistency means the database will reject any operation that would break one of these rules. The data always stays in a valid, meaningful state.

#### Isolation: Operations Do Not Step on Each Other

When 1,000 users are hitting your app at the same time, each database operation appears to run as if it were the only one. They do not interfere with each other. Think of it like each person getting their own private copy of the library to work in, with changes merged in an orderly way.

#### Durability: Committed Data Survives Crashes

Once the database says "yes, I saved that," it means it. Even if the server loses power one second later, when it comes back up, your data is there. The DBMS uses techniques like writing to a log file before making changes, so it can always recover.

### Indexing: Finding Data Without Reading Everything

Back to the plain file problem: to find Grade 11 students, you had to read every line. In a database with 30,000 students, that could mean reading 30,000 records just to find 500.

A database can create an **index** on a column. An index is like the index at the back of a textbook — instead of reading every page to find where "photosynthesis" appears, you flip to the index, find "photosynthesis," and it tells you the exact page numbers.

A database index works the same way: instead of scanning every row, it maintains a separate sorted structure that points directly to the rows you want. Searches that took seconds now take milliseconds.

### Querying: Asking Questions in a Language

Plain files give you data. Databases give you a language to ask complex questions about data.

SQL (Structured Query Language, pronounced "sequel" or "S-Q-L") is the most common such language. You can ask things like:

```sql
SELECT name, grade
FROM students
WHERE grade = 11 AND math_score > 80
ORDER BY math_score DESC;
```

This reads as plain English: "Give me the names and grades of all Grade 11 students with a math score above 80, sorted from highest to lowest score." The database figures out the fastest way to answer that — you just describe what you want.

### Concurrency: Many Users at Once

A plain file has no built-in way to handle two people changing it at the same time. Databases are designed from the ground up for this. They use techniques like locking (like "I am editing this row, please wait") and multi-version concurrency control to let thousands of reads and writes happen simultaneously without data getting scrambled.

### Quick Check

> 1. What does ACID stand for? Describe each letter in one sentence.
> 2. What is a database index, and why does it make searches faster?
> 3. What problem does "Atomicity" solve in a bank transfer example?

---

## Types of Databases

Databases come in several fundamentally different shapes. Choosing the right type for your problem is one of the most important architectural decisions you will make as a developer.

### Relational Databases (SQL)

Think of a relational database like a collection of spreadsheets that are linked together.

Each **table** has rows (records) and columns (fields). A users table might look like this:

| id | name    | email              | age |
|----|---------|---------------------|-----|
| 1  | Alice   | alice@example.com   | 17  |
| 2  | Bob     | bob@example.com     | 16  |
| 3  | Charlie | charlie@example.com | 17  |

The power comes from **relationships** between tables. An orders table can reference the users table by storing the user's `id`. This is called a **foreign key**. You can then ask: "give me all orders placed by users in Grade 11" — and the database joins the two tables together automatically.

Relational databases use **SQL** to query data. Examples include:

- **PostgreSQL** — powerful, open-source, used by many startups and large companies
- **MySQL / MariaDB** — very popular for web applications
- **SQLite** — a self-contained database that lives in a single file; perfect for learning and small apps
- **Microsoft SQL Server** — common in enterprise and corporate environments

**When to use relational databases:** Your data has clear structure and relationships. You need complex queries. Data integrity rules are important (finance, healthcare, inventory).

### Document Databases (NoSQL)

Imagine instead of a spreadsheet, you store each record as a JSON document — a self-contained blob of data. Each document can have different fields.

```json
{
  "id": "user_001",
  "name": "Alice",
  "age": 17,
  "hobbies": ["chess", "coding", "running"],
  "address": {
    "city": "Toronto",
    "country": "Canada"
  }
}
```

Notice that `hobbies` is a list, and `address` is a nested object. In a relational database, these would need extra tables. In a document database, everything about Alice lives in one document.

The most famous document database is **MongoDB**. Others include CouchDB and Firestore (used by Firebase/Google).

**When to use document databases:** Your data is hierarchical or varies a lot between records. You are building apps that evolve quickly and the schema (structure) changes often. Content management systems, user profiles, product catalogs.

### Key-Value Stores (NoSQL)

This is the simplest type of database. Think of it like a dictionary (or a Go `map`): you have a key, and you have a value. That is it.

```
"session:abc123"  -->  "{user_id: 42, expires: 1718000000}"
"cache:homepage"  -->  "<html>...</html>"
"counter:pageviews"  -->  "1847293"
```

You look things up by their exact key. There is no querying by content — just "give me the thing with this key."

The most famous key-value store is **Redis**. It stores data in memory (RAM) rather than on disk, making it extraordinarily fast — hundreds of thousands of operations per second.

**When to use key-value stores:** Caching (storing results of expensive operations so you do not repeat them), session management (storing who is logged in), rate limiting, real-time leaderboards and counters.

### Column-Family Databases (NoSQL)

Imagine a spreadsheet but designed for billions of rows, where most cells are empty, and you only ever need to read specific columns — not full rows.

These databases store data column by column rather than row by row. This sounds odd until you realize: if you want the average age of all your users, and you have a billion users, you only need the "age" column — not names, emails, or addresses. Column storage lets you read just what you need.

Examples: **Apache Cassandra**, **HBase**, **Google Bigtable** (the system that inspired many of these).

**When to use column-family databases:** Massive-scale time-series data (logs, metrics, sensor readings), where you always query by the same column patterns and horizontal scalability across many servers is essential.

### Graph Databases (NoSQL)

Think of a social network. Alice is friends with Bob, Bob is friends with Charlie, Charlie follows Alice. These relationships form a web of connections — a **graph**.

A graph database stores data as **nodes** (things) and **edges** (relationships between things). Questions like "find all people within 3 degrees of connection from Alice" are nearly impossible to do efficiently in a relational database but are native to a graph database.

Examples: **Neo4j**, **Amazon Neptune**.

**When to use graph databases:** Social networks, recommendation engines ("people who liked X also liked Y"), fraud detection (finding suspicious patterns of transactions), knowledge graphs.

### A Quick Comparison

| Type        | Best For                       | Example Products         |
|-------------|--------------------------------|--------------------------|
| Relational  | Structured data, complex queries, integrity | PostgreSQL, MySQL, SQLite |
| Document    | Flexible schema, nested data   | MongoDB, Firestore        |
| Key-Value   | Speed, caching, simple lookups | Redis, DynamoDB           |
| Column      | Massive scale, analytics       | Cassandra, BigTable       |
| Graph       | Relationships, networks        | Neo4j, Neptune            |

### Quick Check

> 1. What is the main difference between a relational database and a document database?
> 2. Why would you choose Redis over PostgreSQL for storing user session data?
> 3. What type of database would you use to build a social network's "friend-of-a-friend" feature?

---

## The Client-Server Model

Almost every database you will use in professional software works as a **server** — a separate program running somewhere (on the same computer or on a different machine entirely) that listens for requests.

Your Go application is the **client** — it sends requests to the database server over a network connection.

```
[Your Go App]  ---network connection--->  [Database Server]
  (client)           (TCP/IP)              (e.g., PostgreSQL)
                                                  |
                                           [Data Files on Disk]
```

Here is how a typical interaction works:

1. Your Go app starts and **opens a connection** to the database server (providing an address, username, and password).
2. The database authenticates your app — "yes, this client is allowed in."
3. Your app sends a **query** (e.g., `SELECT * FROM users WHERE id = 5`).
4. The database server processes the query, finds the data, and sends back a **result set**.
5. Your app reads the results and uses them (e.g., displays them on a web page).
6. Eventually, your app **closes the connection** (or keeps it open for the next request).

### Connection Pools

Opening a new network connection is slow — it can take tens of milliseconds. For a web server handling hundreds of requests per second, opening a new database connection for each request would be far too slow.

The solution is a **connection pool**: a group of connections that are kept open and reused. When your app needs to talk to the database, it borrows a connection from the pool, uses it, and returns it. This is so important that Go's standard library `database/sql` package manages connection pools automatically.

### SQLite: The Exception

**SQLite** is different from most databases. It is **not** a server. The entire database lives in a single `.db` file on disk, and your program reads and writes that file directly — no network, no separate server process.

This makes SQLite perfect for:
- Learning and experimenting (no setup required)
- Mobile apps (each phone has its own local database)
- Desktop apps
- Small websites with low traffic
- Embedded devices

Every Android and iOS phone has SQLite built in. Every Chrome browser uses SQLite to store your browsing history. It is the most widely deployed database engine in the world.

For the rest of this chapter, we will use SQLite because it requires zero installation — you just write Go code and it works.

### Quick Check

> 1. In the client-server model, what role does your Go application play?
> 2. What is a connection pool and why is it important?
> 3. Why is SQLite different from databases like PostgreSQL or MySQL?

---

## Overview: Databases We Will Cover in This Course

Here is a map of what is coming in this course, and the key question each database answers:

**PostgreSQL** — "I need a reliable, powerful relational database for a real application. It should handle complex queries and enforce data integrity."

**SQLite** — "I need a simple, zero-setup database for learning, testing, or a small application."

**Redis** — "I need to store and retrieve data extremely fast. I am building a cache, a session store, or a real-time feature."

**MongoDB** — "My data does not fit neatly into tables. I need flexible, document-based storage."

We will build real projects with each of these, and by the end of this course you will know not just how to use each one, but when to choose each one.

---

## Real-World: Which Databases Power Netflix, Instagram, and Twitter?

Let us look at how some of the most famous apps in the world actually store data. This will make the abstract types we just discussed feel very concrete.

### Netflix

Netflix serves over 230 million subscribers in 190 countries. They use multiple databases for different purposes:

- **MySQL** (relational) — billing data, user account information. This data needs to be precise and consistent. You cannot accidentally double-charge a user.
- **Cassandra** (column-family) — storing what you have watched, your viewing history. Cassandra can handle millions of writes per second across many servers — perfect for tracking billions of viewing events.
- **Redis** (key-value) — caching. When you load the Netflix homepage, the list of recommended shows is pre-computed and stored in Redis so it loads in milliseconds instead of running a complex algorithm every time.
- **EVCache** (custom key-value, built on memcached) — Netflix's own high-performance caching layer.

**The lesson:** No company uses just one database. Different parts of the system have different needs.

### Instagram

Instagram handles over 100 million photos uploaded per day and 500 million daily active users. Their database choices:

- **PostgreSQL** — the primary database for user data, posts, followers, and likes. Instagram famously pushed PostgreSQL to its limits before eventually adding sharding (splitting data across many database servers).
- **Cassandra** — for feeds and activity data (who followed who, recent activity).
- **Redis** — session management (keeping you logged in) and caching hot data.

**The lesson:** Relational databases can scale further than many people think, especially with good engineering.

### Twitter (now X)

Twitter's data is fundamentally about relationships: who follows whom, and what did they post. This is challenging because a tweet from a celebrity can instantly need to appear in the feeds of 50 million followers.

- **MySQL** (relational) — core user and tweet data storage, heavily customized.
- **Redis** — caching timelines. Your Twitter feed is often pre-computed and stored in Redis so it loads instantly.
- **Blobstore** — for media files (images, videos).
- **FlockDB** — Twitter actually built their own graph database to store the follower/following relationships because none of the existing options handled their scale.

**The lesson:** At extreme scale, companies sometimes build custom database solutions. But this is very rare — most companies never need this.

---

## Your First Database Interaction: SQLite in Go

Now let us write real code. We are going to use SQLite from Go to create a database, store some data, and read it back.

### Setting Up

SQLite support in Go requires one external package. In your terminal, create a new directory for this project and initialize a Go module:

```bash
mkdir chapter03-sqlite
cd chapter03-sqlite
go mod init chapter03
```

Now install the SQLite driver. We will use the most common pure-Go SQLite driver:

```bash
go get modernc.org/sqlite
```

This package is a pure Go implementation of SQLite — no C compiler required, which makes it easy to get running on any system.

### Creating Your First Database

Here is a complete, runnable Go program that creates a database, adds a table, inserts some data, and queries it back:

```go
package main

import (
	"database/sql"  // The standard library package for talking to SQL databases
	"fmt"
	"log"

	_ "modernc.org/sqlite" // The SQLite driver — the underscore means we import
	                        // it only for its side effects (registering the driver)
)

func main() {
	// Open (or create) a SQLite database file called "school.db"
	// If the file does not exist, SQLite creates it automatically
	db, err := sql.Open("sqlite", "school.db")
	if err != nil {
		log.Fatal("Failed to open database:", err)
	}
	// defer means "run this when the function returns" — always close your DB connection
	defer db.Close()

	// Verify the connection actually works by pinging the database
	err = db.Ping()
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	fmt.Println("Connected to the database successfully!")

	// Create a table called "students" if it does not already exist
	// A table is like a spreadsheet with defined columns
	createTableSQL := `
		CREATE TABLE IF NOT EXISTS students (
			id      INTEGER PRIMARY KEY AUTOINCREMENT,
			name    TEXT    NOT NULL,
			grade   INTEGER NOT NULL,
			score   REAL    NOT NULL
		);
	`
	// Exec runs SQL that does not return rows (CREATE, INSERT, UPDATE, DELETE)
	_, err = db.Exec(createTableSQL)
	if err != nil {
		log.Fatal("Failed to create table:", err)
	}
	fmt.Println("Table created (or already exists)!")

	// Insert three student records into the table
	// The ? symbols are placeholders — we fill them in with actual values below
	// This is called a "prepared statement" and protects against SQL injection attacks
	insertSQL := `INSERT INTO students (name, grade, score) VALUES (?, ?, ?)`

	students := []struct {
		name  string
		grade int
		score float64
	}{
		{"Alice", 11, 92.5},
		{"Bob", 10, 78.0},
		{"Charlie", 11, 88.5},
		{"Diana", 10, 95.0},
		{"Eve", 11, 76.5},
	}

	for _, s := range students {
		_, err = db.Exec(insertSQL, s.name, s.grade, s.score)
		if err != nil {
			log.Fatal("Failed to insert student:", err)
		}
	}
	fmt.Println("Inserted 5 students!")

	// Query: find all Grade 11 students with a score above 80
	// Rows is like a cursor — it points to each result row one at a time
	rows, err := db.Query(
		"SELECT id, name, score FROM students WHERE grade = ? AND score > ? ORDER BY score DESC",
		11,   // grade to search for
		80.0, // minimum score
	)
	if err != nil {
		log.Fatal("Failed to query database:", err)
	}
	// Always close rows when done — this releases database resources
	defer rows.Close()

	fmt.Println("\nGrade 11 students with score above 80:")
	fmt.Println("----------------------------------------")

	// Loop through each result row
	for rows.Next() {
		var id int
		var name string
		var score float64

		// Scan reads the values from the current row into our variables
		err = rows.Scan(&id, &name, &score)
		if err != nil {
			log.Fatal("Failed to read row:", err)
		}
		fmt.Printf("ID: %d | Name: %-10s | Score: %.1f\n", id, name, score)
	}

	// Check if looping through rows caused any error
	if err = rows.Err(); err != nil {
		log.Fatal("Error after reading rows:", err)
	}
}
```

### Running the Program

```bash
go run main.go
```

You should see:

```
Connected to the database successfully!
Table created (or already exists)!
Inserted 5 students!

Grade 11 students with score above 80:
----------------------------------------
ID: 1 | Name: Alice      | Score: 92.5
ID: 3 | Name: Charlie    | Score: 88.5
```

Notice that Eve (Grade 11, score 76.5) does not appear — she is in Grade 11 but her score is below 80. Bob and Diana are Grade 10, so they are filtered out too. The database handled all of this filtering for us.

Run the program a second time. You will notice the students get inserted again, giving you duplicates. This is intentional for now — we will learn how to prevent duplicates with unique constraints in a later chapter.

### Understanding the Code Line by Line

Let us break down the most important parts:

**`sql.Open("sqlite", "school.db")`** — This opens a connection to a SQLite database stored in a file called `school.db`. The first argument `"sqlite"` is the driver name (registered by the `modernc.org/sqlite` package). The second argument is the data source name — for SQLite, this is just a file path.

**`_ "modernc.org/sqlite"`** — The underscore import is a Go convention for packages you import only for their side effects. When the SQLite package is imported, it registers itself with Go's `database/sql` package so you can use `"sqlite"` as a driver name. You never call functions from this package directly.

**`CREATE TABLE IF NOT EXISTS`** — This SQL command creates a new table, but only if no table with that name already exists. This means you can run the program multiple times without errors.

**`INTEGER PRIMARY KEY AUTOINCREMENT`** — The `id` column is special. `PRIMARY KEY` means it uniquely identifies each row. `AUTOINCREMENT` means the database automatically assigns the next available number (1, 2, 3, ...) — you do not have to provide it.

**`TEXT NOT NULL`** — `TEXT` is the data type (a string). `NOT NULL` is a constraint — the database will reject any attempt to insert a student without a name.

**`?` placeholders** — Never build SQL queries by concatenating user input into strings. That is how SQL injection attacks work (a hacker types `'; DROP TABLE students; --` as their name). Placeholders let the database handle values safely, completely separated from the SQL command structure.

**`rows.Next()`** — This advances the cursor to the next row. It returns `true` if there was a row, `false` when you have read all rows or if there was an error.

**`rows.Scan(&id, &name, &score)`** — This reads the values from the current row into your Go variables. The `&` means "the address of" this variable — Scan needs to know where to write the values.

### A Note on `database/sql`

Notice that our Go code uses `database/sql` from the standard library — not anything from the `modernc.org/sqlite` package directly. This is intentional and powerful: Go's `database/sql` package defines a **universal interface** for all SQL databases. The same code structure works for PostgreSQL, MySQL, SQLite, or almost any other SQL database — you just change the driver import and the connection string. This will save you enormous amounts of work as we move through this course.

### Quick Check

> 1. What does `defer db.Close()` do, and why is it important?
> 2. What is a SQL placeholder (`?`) and why should you use it instead of building query strings with string concatenation?
> 3. What does `rows.Scan()` do?

---

## Mini Project: Student Grade Tracker

Extend the program above to build a simple student grade tracker with the following features:

1. A `students` table with name, grade level, and a score.
2. A function `addStudent(db *sql.DB, name string, grade int, score float64) error` that inserts a new student.
3. A function `getTopStudents(db *sql.DB, minScore float64) ([]string, error)` that returns the names of all students with a score above the given minimum.
4. A `main` function that adds at least 6 students and then prints the top students with a score above 85.

This project will test whether you understand the full read-write cycle with a database.

---

## Exercises

### Easy

1. Modify the program from this chapter to also print the grade level of each student in the results.
2. Add a new SQL query that counts how many students are in each grade. The SQL keyword `COUNT(*)` will help you. Hint: `SELECT grade, COUNT(*) FROM students GROUP BY grade`.
3. What would happen if you ran the program five times? How many rows would be in the database? How could you check? (Hint: look for a SQLite command-line tool, or add a query that counts all rows.)

### Medium

4. Add a new table called `teachers` with columns for `id`, `name`, and `subject`. Insert three teachers and query them back.
5. Modify the query to find the student with the highest score in each grade. Research the SQL `MAX()` function and `GROUP BY` clause.
6. Add error handling: if inserting a student fails, print a useful error message but do not crash the program. Continue inserting the remaining students.

### Hard

7. Refactor the program so that all the database operations happen inside separate functions: `createSchema`, `seedData`, and `queryTopStudents`. Each function should take a `*sql.DB` as its first argument.
8. Research the concept of a SQL **transaction**. Wrap the five INSERT statements in a single transaction so that either all five succeed or none of them do. In Go's `database/sql`, look for `db.Begin()`, `tx.Commit()`, and `tx.Rollback()`.
9. The program currently creates duplicate rows if you run it multiple times. Fix this by checking whether the students table already has rows before inserting. (Hint: use `SELECT COUNT(*) FROM students` and only insert if the count is 0.)

---

## Summary

- A **database** is an organized collection of data designed for fast storage, retrieval, and updating. A **DBMS** (Database Management System) is the software that manages it — examples include PostgreSQL, MySQL, SQLite, MongoDB, and Redis.
- Databases beat plain files because they provide **ACID guarantees** (Atomicity, Consistency, Isolation, Durability), **indexing** for fast lookups, **SQL** for expressive queries, and **concurrency control** so many users can work simultaneously without corrupting data.
- The five main database types are: **relational (SQL)**, **document**, **key-value**, **column-family**, and **graph**. Each has strengths suited to specific problems — real companies like Netflix and Instagram use multiple types at once.
- Most databases follow a **client-server model** where your application sends queries over a network connection. **SQLite** is the notable exception: a self-contained file-based database with no separate server process.
- Go's **`database/sql`** standard library package provides a universal interface for SQL databases. You import a specific driver (like `modernc.org/sqlite`) and use the same `db.Exec`, `db.Query`, and `rows.Scan` pattern regardless of which database you are using.
