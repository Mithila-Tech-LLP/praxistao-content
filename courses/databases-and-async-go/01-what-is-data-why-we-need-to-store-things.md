# Chapter 01: What Is Data? Why We Need to Store Things

Every app you have ever used — Instagram, your school's grade portal, a banking app — is, at its core, a machine for storing and retrieving data. Before we write a single line of Go that touches a database, we need to understand what data actually is, why losing it is a catastrophe, and why a simple text file is not enough when you have a million users depending on you.

## Table of Contents

1. [What Is Data?](#1-what-is-data)
2. [How Humans Stored Data Before Computers](#2-how-humans-stored-data-before-computers)
3. [Data in Computer Memory — And Why It Disappears](#3-data-in-computer-memory--and-why-it-disappears)
4. [Persistent Storage: Hard Drives and SSDs](#4-persistent-storage-hard-drives-and-ssds)
5. [The Problem of Scale](#5-the-problem-of-scale)
6. [Files vs Databases](#6-files-vs-databases)
7. [Mini Project: A Name List in a File](#7-mini-project-a-name-list-in-a-file)
8. [Summary](#summary)
9. [Exercises](#exercises)

---

## 1. What Is Data?

Imagine you walk into a coffee shop and the barista asks for your name so they can write it on your cup. Your name — "Priya", "Carlos", "Sam" — is a piece of **data**. It is a fact about you that someone needed to write down and use later.

Data is any piece of information that can be recorded and used. That definition is intentionally broad, because data really is everywhere.

### Real-World Examples of Data

Here are examples you interact with every day:

| What you do | The data being created |
|---|---|
| Score 8200 points in a game | A number linked to your username |
| Deposit $50 into a bank account | A new balance: old amount + 50 |
| Sign up for a website | Your email, password, and the date you signed up |
| Search "how to make pasta" | The words you typed, the time, your location |
| Like a photo | A connection between your account ID and that photo's ID |

Notice a pattern: data is always **structured**. It has a *what* (the value: "Priya", 8200, $50.00) and a *context* (what that value means: a name, a score, a bank balance). A number on its own means nothing. The number 42 sitting alone is not data in any useful sense. But "user ID 7's high score is 42" — that is data, because it has meaning.

### The Three Flavors of Data You Will See Most

**Numbers** — ages, prices, distances, scores. Computers handle these very efficiently.

**Text** — names, messages, addresses, search queries. Humans think in text; computers have to work a bit harder to handle it.

**Relationships** — "this user owns these orders", "this author wrote these books". Relationships between pieces of data are often the most valuable thing of all. We will spend a lot of time on these later in the course.

### Quick Check

> 1. Name three pieces of data that get created when you send a text message.
> 2. Is the color "blue" data? What would make it useful data?
> 3. What is the difference between the number `1000` and the data point "Ana's bank balance is $1000"?

---

## 2. How Humans Stored Data Before Computers

Long before anyone had heard of a hard drive or a database, humans were already solving the exact same problem: how do you remember important information reliably, and how do you find it again quickly?

### The Library Card Catalog

Walk into an old library and you would find a wall of small wooden drawers. Each drawer held index cards, one per book. Every card told you the book's title, its author, and — crucially — its location on the shelves. The cards were sorted alphabetically so that a librarian could find any book's location in under a minute.

This is a **manual database**. The cards are records. The drawers are organized storage. The alphabetical sorting is an **index** — a way to find things quickly without reading every single card. These are concepts we will return to throughout this course.

### Ledgers and Accounting Books

A bank in 1920 did not have computers. It had large paper books called **ledgers**. Each page tracked a customer's account: every deposit written in, every withdrawal written out, and the current balance calculated by hand. Multiple clerks might update different ledgers at the same time, and at the end of the day a senior accountant would check that all the numbers added up.

This introduced a problem that still exists today in software: **concurrency**. What happens if two clerks try to update the same account at the same time? A good bank had strict rules: only one clerk could hold a ledger at a time. We will see that databases have the exact same concept, called **locking**.

### Filing Cabinets

A doctor's office kept one paper folder per patient, stuffed with their medical history, test results, and prescriptions. Finding the folder for a patient named "Monica Chen" meant opening the "C" drawer, flipping through dozens of folders, and pulling the right one. Adding a new test result meant pulling the folder, clipping in the new page, and filing the folder back.

This works fine for a hundred patients. But imagine doing it for ten million patients across a national hospital network. The filing cabinet breaks down. The speed of searching, the risk of misfiling, the sheer physical space required — it all becomes unmanageable.

This is the fundamental problem that databases solve. Everything else in this course is a more sophisticated answer to the question: "How do we store, find, and update information reliably, even when the amount of information is enormous?"

### Quick Check

> 1. What is the library card catalog an analogy for in databases?
> 2. What is the "concurrency" problem, and how did banks solve it before computers?
> 3. Why does the filing cabinet approach fail at large scale?

---

## 3. Data in Computer Memory — And Why It Disappears

Now let us talk about how computers handle data. To understand why we need databases, we first need to understand the difference between two types of storage inside a computer.

### RAM: The Whiteboard

Think of your computer's **RAM** (Random Access Memory) as a whiteboard in a classroom. You can write on it fast. You can read from it fast. You can erase and rewrite in an instant.

But here is the catch: the whiteboard only works while the lights are on. The moment you turn off the power, everything written on the whiteboard is gone. Tomorrow morning, the whiteboard is blank.

RAM is the same way. It is **volatile memory** — "volatile" means it does not survive without power. Everything your Go program stores in variables, slices, and maps lives in RAM. The moment your program stops — either intentionally or because of a crash or power cut — that data vanishes.

Here is a tiny Go program that demonstrates this. It stores some names in a slice (a list, in Go):

```go
package main

import "fmt"

func main() {
    // This creates a slice (a list) of strings in RAM.
    // A string is just text data.
    names := []string{"Alice", "Bob", "Carol"}

    // We can add a new name to the list.
    names = append(names, "David")

    // We can print the list to see it.
    fmt.Println("Names stored in memory:", names)

    // When this function ends, the program exits.
    // The 'names' slice is gone forever.
    // The next time you run this program, it starts from scratch.
}
```

Run this program and you will see:

```
Names stored in memory: [Alice Bob Carol David]
```

Now run it a second time. The output is exactly the same. "David" was added in the first run, but the list is back to its original state. The computer did not remember what happened last time. Every run starts fresh because RAM was cleared when the program exited.

This is not a bug. RAM is designed to be fast and temporary. Your web browser uses RAM to store the webpage you are looking at right now. When you close the browser, it does not need to remember that page — you can just load it again. But a bank cannot "just load it again" when you close the app. Your balance needs to survive.

### What Lives in RAM in a Go Program

When you write a Go program:

- Variables like `name := "Alice"` live in RAM.
- Slices like `scores := []int{95, 87, 62}` live in RAM.
- Maps like `users := map[string]int{"alice": 1, "bob": 2}` live in RAM.

All of these are gone the moment your program stops. If your program crashes at 2am while processing a payment, and the data was only in RAM, the payment is lost. This is unacceptable for almost any real application.

### Quick Check

> 1. What does "volatile" mean when we talk about computer memory?
> 2. In the Go code above, where does the `names` slice live while the program runs?
> 3. Why is it okay for a web browser to use RAM for a loaded page, but not okay for a bank to use RAM for your balance?

---

## 4. Persistent Storage: Hard Drives and SSDs

The solution to volatile memory is **persistent storage** — storage that survives power loss.

### The Hard Drive: A Library That Never Forgets

Think of a hard drive or SSD (Solid State Drive) as a library. Unlike the whiteboard (RAM), the library does not empty out overnight. Books placed on the shelves stay there indefinitely, surviving power cuts, reboots, and even the librarian going on vacation.

**Hard drives** (HDD) store data on spinning magnetic disks, similar to how a vinyl record works. A tiny read/write head moves across the spinning disk to find and write data. This mechanical movement makes hard drives relatively slow but very cheap for large amounts of storage.

**SSDs** (Solid State Drives) have no moving parts. They store data as electrical charges in flash memory chips — similar to a USB drive, just faster and more robust. SSDs are significantly faster than hard drives.

The key property both share is **persistence**: data written to them survives indefinitely, with or without power.

When a Go program writes data to a file on disk, that data stays there even after the program exits. When the program runs again, it can read the data back.

This is the foundation of all databases. Every database — whether it is SQLite, PostgreSQL, or MongoDB — ultimately stores its data as files on disk. A database is, at its heart, a program that manages files in a very smart and structured way.

### How Fast Is "Slow"?

Understanding speed differences matters in this course, so here is a rough comparison:

| Storage Type | Approximate Read Speed |
|---|---|
| CPU Cache (inside the CPU) | ~1 nanosecond |
| RAM | ~100 nanoseconds |
| SSD | ~100 microseconds (1,000x slower than RAM) |
| Hard Drive | ~10 milliseconds (100,000x slower than RAM) |

This is why databases spend enormous effort keeping frequently accessed data in RAM (called a **cache** or **buffer pool**) and only going to disk when necessary. We will explore this concept in depth in later chapters.

### Quick Check

> 1. What is the difference between a hard drive and an SSD?
> 2. What does "persistent" mean in the context of storage?
> 3. Roughly how much slower is an SSD compared to RAM?

---

## 5. The Problem of Scale

So far, the solution sounds simple: write data to disk and read it back. Let us explore why that alone is not enough once your application grows.

### Start Small: One User

Imagine you are building a leaderboard for a game. On day one, you have 10 players. You store their names and scores in a file. When a player beats their high score, you rewrite the file. Simple, elegant, fast.

### Grow to a Thousand Users

With a thousand players, the file gets larger. When you want to find the top 10 scores, your program has to read the entire file, load all 1,000 records into RAM, sort them, and return the top 10. It is still fast enough that users do not notice.

But what if two players beat their high scores at exactly the same moment? Both of their scores need to be written to the file simultaneously. If you are not careful, the two write operations collide, and one of the updates overwrites the other. One player's new high score is lost. This is the same concurrency problem the bank clerks had with ledgers.

### Grow to One Million Users

Now things get difficult. A file with a million records might be hundreds of megabytes. Every time you want to find a single user's score, your program has to:

1. Open the file (a disk operation — slow).
2. Read through it line by line until it finds the right user (reading potentially millions of lines).
3. Update that one line.
4. Write the entire file back to disk.

That single lookup might take several seconds. With a million users hitting your app at the same time, each waiting seconds for a response, your server grinds to a halt. This is called a **bottleneck**.

### The Specific Problems That Emerge at Scale

**Searching is slow.** A file has no structure. To find "Carlos" among a million names, you have to read every name one by one. A database uses data structures called **indexes** (remember the library card catalog?) that let it find any record in milliseconds regardless of how many records there are.

**Updates are dangerous.** If your program crashes halfway through rewriting a file, you can end up with a corrupt, half-written file. A database uses a technique called **transactions** with a **write-ahead log** to ensure that an update either fully completes or fully rolls back — never half-completes.

**Concurrent access breaks things.** Two users updating the same file at the same moment causes corruption. Databases handle this with sophisticated **locking** and **isolation** mechanisms.

**No relationships.** A file can store a list of users. Another file can store a list of orders. But how do you link "this order belongs to this user"? In a file, you have to build this linking logic yourself, and it gets complicated fast. A relational database handles this natively.

**No queries.** With a file, you write code to answer questions like "find all users who scored over 9000". With a database, you use a query language (SQL, which we will learn) that lets you ask complex questions in a single line.

These are the reasons databases exist. They are not magic — they are programs written by very clever people to solve exactly these problems, so that you do not have to solve them yourself for every application you build.

### Quick Check

> 1. What is a "bottleneck" in the context of reading a large file?
> 2. What can go wrong if two users try to update the same file at the same time?
> 3. Name three specific limitations of using plain files that databases are designed to solve.

---

## 6. Files vs Databases

Let us make the comparison explicit before we dive into the mini project.

### A Plain File

A plain file is like a notebook. You write things in it, you read them back. It is simple and it works for small, simple data.

**When files are fine:**
- Configuration files (your app's settings — read once at startup).
- Log files (you only ever append new lines; you rarely search through them in real time).
- Exporting data for a human to read in a spreadsheet.
- Storing a single user's preferences in a desktop app.

**When files break down:**
- Multiple users or processes need to read and write at the same time.
- You need to search, filter, or sort data quickly.
- You cannot afford to lose any data, even if the program crashes mid-write.
- The data has relationships (users own orders, orders contain products).
- The data grows beyond a few thousand records.

### A Database

A database is a program whose entire purpose is to store, organize, and retrieve data reliably and efficiently. It sits between your Go program and the disk, handling all of the hard problems:

```
Your Go Program
      |
      |  (sends queries: "give me user 42's score")
      v
 [Database]   <-- handles concurrency, indexing, crash recovery
      |
      |  (reads/writes files in a structured format)
      v
 [Disk Files]
```

There are several types of databases, each designed for different problems:

**Relational Databases (e.g., PostgreSQL, MySQL, SQLite)** — Store data in tables, like spreadsheets. Excellent for structured data with relationships. Use SQL as their query language. These are what most applications use, and they are the focus of this course.

**Key-Value Stores (e.g., Redis)** — Store data as simple key-value pairs, like a giant dictionary. Blazingly fast. Often used as a cache (keeping frequently accessed data in RAM). We will use Redis later in this course for async job queues.

**Document Databases (e.g., MongoDB)** — Store data as documents (similar to JSON objects). Flexible structure. Good when your data does not have a rigid shape.

For this course, we will primarily use **PostgreSQL** (a powerful, open-source relational database) and **Redis** (for async/queue patterns). But before we get there, let us build something with a plain file so you can feel, firsthand, exactly why files are not enough.

---

## 7. Mini Project: A Name List in a File

Let us write a real Go program that stores a list of names in a file. We will build it step by step, then walk through what its limitations are. This is your first real Go program in this course — read every comment.

### What We Are Building

A small command-line program that:
1. Reads a list of names from a file (if the file exists).
2. Lets you add a new name.
3. Saves the updated list back to the file.
4. Shows you all stored names.

The file will simply store one name per line, like this:

```
Alice
Bob
Carol
```

### Project Setup

Create a new directory for this project and create a file called `main.go` inside it:

```
mkdir namelist
cd namelist
touch main.go
```

### The Complete Program

```go
package main

// We import three packages (libraries) that our program needs.
// "bufio"   - helps us read a file line by line efficiently
// "fmt"     - lets us print output and format strings
// "os"      - gives us tools to work with files and the operating system
import (
    "bufio"
    "fmt"
    "os"
    "strings"
)

// The name of the file where we will store our names.
// We define it as a constant so it is easy to change in one place.
const filename = "names.txt"

// loadNames reads the file and returns a slice (list) of names.
// If the file does not exist yet, it returns an empty list — not an error.
func loadNames() ([]string, error) {
    // os.Open tries to open the file for reading.
    // It returns a file handle (f) and an error (err).
    f, err := os.Open(filename)
    if err != nil {
        // os.IsNotExist checks if the error is specifically "file not found".
        // If so, this is not a real error — the file just has not been created yet.
        if os.IsNotExist(err) {
            return []string{}, nil // Return an empty list, no error.
        }
        // Any other error (e.g., permission denied) is a real problem.
        return nil, fmt.Errorf("could not open file: %w", err)
    }
    // defer means "run this line when the function exits, no matter what".
    // Closing the file when we are done is important — otherwise the operating
    // system keeps a lock on it.
    defer f.Close()

    // bufio.NewScanner creates a scanner that reads the file one line at a time.
    // Reading line by line is more memory-efficient than loading the whole file at once.
    scanner := bufio.NewScanner(f)

    // names will collect all the names we find in the file.
    var names []string

    // scanner.Scan() returns true each time it successfully reads a line.
    // When there are no more lines, it returns false and the loop ends.
    for scanner.Scan() {
        line := scanner.Text() // scanner.Text() gives us the current line as a string.

        // strings.TrimSpace removes any extra spaces or newline characters
        // from the beginning and end of the line.
        line = strings.TrimSpace(line)

        // Skip empty lines — they are not names.
        if line != "" {
            names = append(names, line)
        }
    }

    // scanner.Err() returns any error that occurred while scanning.
    // (scanner.Scan() itself does not return errors — they are collected here.)
    if err := scanner.Err(); err != nil {
        return nil, fmt.Errorf("error reading file: %w", err)
    }

    return names, nil
}

// saveNames writes a slice of names to the file, one name per line.
// It overwrites the entire file each time.
func saveNames(names []string) error {
    // os.Create opens the file for writing. If the file already exists,
    // it truncates it (empties it out). If it does not exist, it creates it.
    f, err := os.Create(filename)
    if err != nil {
        return fmt.Errorf("could not create file: %w", err)
    }
    defer f.Close()

    // bufio.NewWriter wraps the file in a buffered writer.
    // A buffered writer collects small writes in memory and sends them
    // to disk in larger, more efficient chunks.
    writer := bufio.NewWriter(f)

    // Write each name as a line in the file.
    for _, name := range names {
        // fmt.Fprintln writes the name followed by a newline character (\n)
        // into the writer (and eventually onto disk).
        _, err := fmt.Fprintln(writer, name)
        if err != nil {
            return fmt.Errorf("error writing name: %w", err)
        }
    }

    // writer.Flush() forces any data still sitting in the buffer to be
    // written to disk. Without this, some names might never reach the file.
    return writer.Flush()
}

func main() {
    // Step 1: Load whatever names are already stored in the file.
    names, err := loadNames()
    if err != nil {
        // If we cannot even read the file, something is seriously wrong.
        // fmt.Println prints to the screen. We then exit the program.
        fmt.Println("Error loading names:", err)
        os.Exit(1)
    }

    fmt.Printf("Loaded %d name(s) from file.\n", len(names))

    // Step 2: Ask the user to type a new name.
    fmt.Print("Enter a new name to add: ")

    // bufio.NewReader(os.Stdin) lets us read text that the user types.
    // os.Stdin is the "standard input" — basically, the keyboard.
    reader := bufio.NewReader(os.Stdin)
    newName, _ := reader.ReadString('\n') // Read until the user presses Enter.
    newName = strings.TrimSpace(newName)  // Remove the trailing newline.

    if newName == "" {
        fmt.Println("No name entered. Exiting without changes.")
        return
    }

    // Step 3: Add the new name to our list.
    names = append(names, newName)

    // Step 4: Save the updated list back to the file.
    if err := saveNames(names); err != nil {
        fmt.Println("Error saving names:", err)
        os.Exit(1)
    }

    // Step 5: Print all the names we now have stored.
    fmt.Println("\nAll stored names:")
    for i, name := range names {
        // %d is the format verb for integers (the index number).
        // %s is the format verb for strings (the name).
        fmt.Printf("  %d. %s\n", i+1, name)
    }

    fmt.Println("\nNames saved to", filename)
}
```

### Running the Program

To run this program, open your terminal in the `namelist` directory and type:

```
go run main.go
```

The first time you run it, no file exists yet, so you will see:

```
Loaded 0 name(s) from file.
Enter a new name to add: Alice
All stored names:
  1. Alice

Names saved to names.txt
```

Run it again and add another name:

```
Loaded 1 name(s) from file.
Enter a new name to add: Bob
All stored names:
  1. Alice
  2. Bob

Names saved to names.txt
```

Notice that Alice is still there! The data survived between runs because we saved it to disk. This is persistence in action.

### What the File Looks Like

Open `names.txt` in any text editor and you will see:

```
Alice
Bob
```

Plain text, one name per line. Simple and readable.

### Now, Let Us Find the Limitations

This program works. But let us think carefully about where it starts to break.

**Problem 1: Searching is slow.**

Suppose instead of names, this file had 500,000 user records. To find whether "Zara" is in the list, the program would have to read every line from line 1 to line 500,000. This is called a **linear scan** — the more data you have, the longer it takes, in a straight line. A database with an index could find Zara in a fraction of a millisecond regardless of how many users exist.

**Problem 2: We rewrite the entire file every time.**

Our `saveNames` function rewrites the whole file even if only one name changed. With 500,000 names, that means writing megabytes of data to disk for every single update. A database only writes the changed record.

**Problem 3: Two users at the same time would corrupt the file.**

Imagine two people running this program at the exact same time. Both call `loadNames()` and get a list back. Both add a name. Both call `saveNames()` with their updated list. The second save overwrites the first — one name is silently lost. This is a **race condition**, and it is a critical bug in any multi-user system.

**Problem 4: A crash during save corrupts the file.**

Our save process is: open the file, erase everything, write new data. If the program crashes after "erase everything" but before "write new data completes", the file is empty. All data is lost. A database uses a **write-ahead log** to prevent this: it records what it is about to do before doing it, so that it can recover if something goes wrong.

**Problem 5: No way to query.**

What if you want "find all names that start with the letter A"? You would have to write your own searching code. With SQL, you would just write `SELECT * FROM names WHERE name LIKE 'A%'` and the database handles it.

These five limitations are exactly what databases are engineered to solve. The rest of this course is the journey from this simple file program to building real, production-quality data systems in Go.

---

## Summary

- **Data** is any piece of information paired with context that makes it meaningful — names, balances, scores, relationships between things.
- Humans solved data storage problems long before computers, using ledgers, card catalogs, and filing cabinets. The same fundamental challenges (searching, concurrency, scale) exist in software.
- **RAM is volatile** — all data stored in Go variables, slices, and maps disappears the moment the program stops. This is by design, but it means RAM alone is not enough for any application that needs to remember things.
- **Persistent storage** (hard drives, SSDs) survives power loss. Databases use disk as their permanent home, with RAM as a fast cache for frequently accessed data.
- **Files work for small, simple data**, but they fail at scale because linear searches are slow, concurrent writes cause corruption, crashes can destroy data, and there is no built-in way to query or relate data.
- **Databases** exist to solve exactly these problems — they are programs that manage disk files in smart, safe, efficient ways so that you do not have to rebuild those solutions yourself.

---

## Exercises

### Easy

1. Modify the mini project so that after loading the names, the program checks whether the new name the user entered already exists in the list. If it does, print "Name already exists" and do not add a duplicate.

2. Add a function called `countNames()` that reads `names.txt` and prints how many names are stored, without loading them all into a slice. Can you do it by counting lines?

3. Right now the names are stored in the order they were added. Add a feature to print the names sorted alphabetically. (Hint: look up the `sort.Strings` function in Go's standard library.)

### Medium

4. Extend the program to store not just names, but names with scores, like `Alice,9200`. Update `loadNames` to parse this format and return a slice of a struct you define:

   ```go
   type Player struct {
       Name  string
       Score int
   }
   ```

   Print the players sorted by score from highest to lowest.

5. The current program rewrites the entire file on every save. Write a new function `appendName(name string) error` that opens the file in **append mode** and adds just the new name to the end, without reading or rewriting the whole file. When would this approach be better than the current one? When would it be worse?

6. Simulate the race condition described in the limitations section. Write two goroutines (concurrent Go functions) that both try to add a name to the file at the same time. Run the program several times and observe what happens to the file contents. Document what you observe.

### Hard

7. Implement a simple crash-safety mechanism for `saveNames`. Instead of writing directly to `names.txt`, write to a temporary file called `names.txt.tmp` first. Once the temporary file is fully written and flushed, rename it to `names.txt` (replacing the old file). Research why a file rename is typically an **atomic** operation on most operating systems, and why this makes the approach safer than writing directly.

8. Build a minimal search index. After saving names to the file, create a second file called `index.txt` that stores each name alongside its **byte offset** in `names.txt` (the position in the file where that name starts), sorted alphabetically. Then write a `findName(query string) (int64, error)` function that uses the index to jump directly to the right position in `names.txt` using `file.Seek()`, instead of scanning from the beginning. Measure how much faster this is for large files.
