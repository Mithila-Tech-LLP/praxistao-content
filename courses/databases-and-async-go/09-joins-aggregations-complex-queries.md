# Chapter 09: Joins, Aggregations, and Complex Queries

Imagine you walk into a library and ask: "Who are the top five students who have checked out the most books this year, and what is their average grade?" To answer that, you would need to look at three separate places — the student list, the checkout records, and the grade book — and combine the information. That is exactly what this chapter is about. Real data lives in multiple places, and SQL gives you powerful tools to pull it all together.

## Table of Contents

1. [Relational Data: Why We Split Things Up](#1-relational-data-why-we-split-things-up)
2. [Foreign Keys: Linking Tables Together](#2-foreign-keys-linking-tables-together)
3. [Joins: Combining Tables](#3-joins-combining-tables)
4. [Aggregations: COUNT, SUM, AVG, MIN, MAX](#4-aggregations-count-sum-avg-min-max)
5. [GROUP BY and HAVING](#5-group-by-and-having)
6. [Subqueries: Queries Inside Queries](#6-subqueries-queries-inside-queries)
7. [Common Table Expressions (CTEs) with WITH](#7-common-table-expressions-ctes-with-with)
8. [Window Functions: ROW_NUMBER, RANK, LAG, LEAD](#8-window-functions-row_number-rank-lag-lead)
9. [Building Complex Queries in Go with sqlx and pgx](#9-building-complex-queries-in-go-with-sqlx-and-pgx)
10. [Mini Project: Analytics Dashboard for an E-Commerce Database](#10-mini-project-analytics-dashboard-for-an-e-commerce-database)
11. [Summary](#summary)
12. [Exercises](#exercises)

---

## 1. Relational Data: Why We Split Things Up

Think about a school notebook where you record every student's name alongside every grade they ever receive. For a student named "Aisha" who takes 10 classes, you write "Aisha" ten times. If Aisha moves and changes her last name, you have to find and fix it in ten places. Miss one, and your records are inconsistent.

Databases solve this with a concept called **normalization** — the practice of organizing data so that each piece of information lives in exactly one place.

Here is the naive way to store an e-commerce order in a single table:

```
order_id | customer_name | customer_email        | product_name | price | quantity
---------|---------------|----------------------|--------------|-------|--------
1        | Alice Chen    | alice@example.com     | Widget A     | 9.99  | 2
2        | Alice Chen    | alice@example.com     | Gadget B     | 24.99 | 1
3        | Bob Smith     | bob@example.com       | Widget A     | 9.99  | 5
```

Notice the problems:
- "Alice Chen" and "alice@example.com" are repeated. If Alice changes her email, we must update multiple rows and hope we find them all.
- "Widget A" and its price appear twice. If the price changes, we might update one row but forget the other.

The normalized version splits this into three tables:

```
customers               products                 orders
---------               --------                 ------
id | name   | email     id | name     | price    id | customer_id | product_id | quantity
---+--------+------     ---+----------+------    ---+-------------+------------+---------
1  | Alice  | alice@    1  | Widget A | 9.99     1  | 1           | 1          | 2
2  | Bob    | bob@      2  | Gadget B | 24.99    2  | 1           | 2          | 1
                                                 3  | 2           | 1          | 5
```

Now:
- Alice's email appears once. Change it once, fixed everywhere.
- Widget A's price appears once. Change it once, fixed everywhere.
- The `orders` table only stores what is unique to each order: which customer, which product, how many.

The `orders` table connects to `customers` and `products` using numbers called **foreign keys**.

---

## 2. Foreign Keys: Linking Tables Together

A **foreign key** is a column in one table that refers to the primary key of another table. It is the "link" that connects related rows across tables.

Think of it like a library catalog card. The card does not contain the full book — it just contains the book's ID number, which tells you where the real book lives.

In our e-commerce example:
- `orders.customer_id` is a foreign key that references `customers.id`
- `orders.product_id` is a foreign key that references `products.id`

Here is how to define this in SQL using PostgreSQL:

```sql
CREATE TABLE customers (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    email TEXT NOT NULL UNIQUE
);

CREATE TABLE products (
    id    SERIAL PRIMARY KEY,
    name  TEXT NOT NULL,
    price NUMERIC(10, 2) NOT NULL
);

CREATE TABLE orders (
    id          SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES customers(id),
    product_id  INTEGER NOT NULL REFERENCES products(id),
    quantity    INTEGER NOT NULL,
    created_at  TIMESTAMP DEFAULT NOW()
);
```

The `REFERENCES customers(id)` part is the foreign key declaration. It tells PostgreSQL: "the value in `customer_id` must exist as an `id` in the `customers` table."

This enforcement is called **referential integrity**. If you try to insert an order with `customer_id = 999` but no customer with id 999 exists, PostgreSQL will reject the insert with an error. This is a safety net — you cannot have orders pointing to customers that do not exist.

Let us also insert some sample data that we will use throughout this chapter:

```sql
INSERT INTO customers (name, email) VALUES
    ('Alice Chen',   'alice@example.com'),
    ('Bob Smith',    'bob@example.com'),
    ('Carol Davis',  'carol@example.com');

INSERT INTO products (name, price) VALUES
    ('Widget A',  9.99),
    ('Gadget B', 24.99),
    ('Doohickey C', 4.99);

INSERT INTO orders (customer_id, product_id, quantity, created_at) VALUES
    (1, 1, 2,  '2024-01-10'),
    (1, 2, 1,  '2024-01-15'),
    (2, 1, 5,  '2024-01-20'),
    (2, 3, 10, '2024-02-01'),
    (3, 2, 3,  '2024-02-05');
-- Note: customer 3 (Carol) has one order, and no one ordered Gadget B in February
```

### Quick Check

1. Why do we split data into multiple tables instead of keeping everything in one big table?
2. What does a foreign key do?
3. What happens in PostgreSQL if you try to insert a row with a foreign key value that does not exist in the referenced table?

---

## 3. Joins: Combining Tables

Now that our data lives in separate tables, how do we pull it back together to answer questions? The answer is the **JOIN**.

A JOIN tells the database: "combine rows from these two tables wherever a certain condition is true."

### The Venn Diagram Mental Model

Imagine two overlapping circles. The left circle is Table A. The right circle is Table B. The overlap in the middle represents rows that exist in both tables (matched by the join condition).

Different types of joins return different parts of this picture:

```
    Table A only    |   Both A and B   |   Table B only
    [  A  |  (A∩B)  |  B  ]

INNER JOIN  →  returns only  (A∩B)      — just the overlap
LEFT JOIN   →  returns       A + (A∩B)  — left circle entirely
RIGHT JOIN  →  returns       (A∩B) + B  — right circle entirely
FULL JOIN   →  returns  A + (A∩B) + B   — everything
```

### INNER JOIN: Only Matching Rows

An INNER JOIN returns rows only when there is a match in both tables. Rows with no match on either side are excluded.

**Real-world analogy:** You have two guest lists for a party. An INNER JOIN gives you only the names that appear on both lists — the confirmed attendees.

```sql
SELECT
    customers.name  AS customer_name,
    products.name   AS product_name,
    orders.quantity
FROM orders
INNER JOIN customers ON orders.customer_id = customers.id
INNER JOIN products  ON orders.product_id  = products.id;
```

Result:
```
customer_name | product_name | quantity
--------------+--------------+---------
Alice Chen    | Widget A     | 2
Alice Chen    | Gadget B     | 1
Bob Smith     | Widget A     | 5
Bob Smith     | Doohickey C  | 10
Carol Davis   | Gadget B     | 3
```

Every row in the result has a real customer and a real product. Notice that all five orders appear because every order has a valid customer and product.

The `ON orders.customer_id = customers.id` part is the **join condition** — it says which columns to use when matching rows between tables.

### LEFT JOIN: Keep All Rows from the Left Table

A LEFT JOIN returns every row from the left table. If there is a matching row in the right table, it fills in those columns. If there is no match, those columns appear as NULL (meaning "no value").

**Real-world analogy:** A teacher's class roster (left table) joined with an attendance sheet (right table). A LEFT JOIN shows every student from the roster, with attendance data filled in where available. Students who were absent show NULL attendance columns — they still appear in the result.

```sql
-- Find all customers and their orders.
-- Customers with no orders should still appear (with NULL order columns).
SELECT
    customers.name,
    orders.id        AS order_id,
    orders.quantity
FROM customers
LEFT JOIN orders ON customers.id = orders.customer_id;
```

If we had a customer "Dave" with no orders, the result would include a row:
```
name     | order_id | quantity
---------+----------+---------
Alice    | 1        | 2
Alice    | 2        | 1
Bob      | 3        | 5
Bob      | 4        | 10
Carol    | 5        | 3
Dave     | NULL     | NULL   ← Dave has no orders, but he still appears
```

The LEFT JOIN is extremely useful for finding things that are missing. Want to find all customers who have never placed an order?

```sql
SELECT customers.name
FROM customers
LEFT JOIN orders ON customers.id = orders.customer_id
WHERE orders.id IS NULL;
```

The WHERE `orders.id IS NULL` filter keeps only the rows where the right side had no match — which is precisely the customers with no orders.

### RIGHT JOIN: Keep All Rows from the Right Table

A RIGHT JOIN is the mirror image of LEFT JOIN. It keeps all rows from the right table, filling in NULL for unmatched left rows.

```sql
-- Find all products and any orders for them.
-- Products with no orders still appear.
SELECT
    products.name  AS product_name,
    orders.id      AS order_id
FROM orders
RIGHT JOIN products ON orders.product_id = products.id;
```

In practice, RIGHT JOIN is rarely used because you can always rewrite it as a LEFT JOIN by switching the table order. Most developers prefer LEFT JOIN for consistency.

### FULL JOIN: Keep Everything

A FULL JOIN (also called FULL OUTER JOIN) returns all rows from both tables. Unmatched rows from either side get NULL for the other side's columns.

```sql
-- Show all customers and all products together, matching where orders connect them.
SELECT
    customers.name AS customer_name,
    products.name  AS product_name
FROM customers
FULL JOIN orders   ON customers.id   = orders.customer_id
FULL JOIN products ON orders.product_id = products.id;
```

FULL JOINs are less common in everyday queries but are invaluable for reconciliation — finding all the discrepancies between two tables.

### Joining More Than Two Tables

You can chain as many JOINs as you need. Each JOIN adds another table to the result:

```sql
-- A readable query with table aliases to keep things short.
-- "c", "o", "p" are aliases — short names we give to tables in this query.
SELECT
    c.name          AS customer,
    p.name          AS product,
    o.quantity,
    o.quantity * p.price AS total_value
FROM orders AS o
JOIN customers AS c ON o.customer_id = c.id
JOIN products  AS p ON o.product_id  = p.id
ORDER BY total_value DESC;
```

The `AS o`, `AS c`, `AS p` parts are **table aliases** — temporary short names we assign to tables within a query so we do not have to type the full table name every time.

`o.quantity * p.price` is a computed column — we calculate it on the fly and give it an alias `total_value`.

### Quick Check

1. What is the difference between an INNER JOIN and a LEFT JOIN?
2. How would you find all products that have never been ordered?
3. What does `NULL` mean in the context of a LEFT JOIN result?

---

## 4. Aggregations: COUNT, SUM, AVG, MIN, MAX

So far, our queries return one row per matching record. But often we want a summary: how many orders were placed? What was the total revenue? What is the average order size? SQL provides **aggregate functions** that collapse many rows into a single summary value.

Think of aggregate functions like a cashier at a checkout counter. Instead of handing you a list of every item's price, they total them up and tell you one number: $47.83.

### COUNT: How Many Rows?

`COUNT(*)` counts all rows. `COUNT(column)` counts only rows where that column is not NULL.

```sql
-- How many orders do we have in total?
SELECT COUNT(*) AS total_orders FROM orders;
-- Result: 5

-- How many distinct customers have placed at least one order?
SELECT COUNT(DISTINCT customer_id) AS active_customers FROM orders;
-- Result: 3
```

`DISTINCT` inside COUNT counts each unique value only once.

### SUM: Add Them All Up

```sql
-- What is the total number of items sold across all orders?
SELECT SUM(quantity) AS total_items_sold FROM orders;
-- Result: 21  (2 + 1 + 5 + 10 + 3)
```

### AVG: The Average

```sql
-- What is the average order quantity?
SELECT AVG(quantity) AS avg_quantity FROM orders;
-- Result: 4.2  (21 / 5)
```

AVG returns a decimal even if all input values are integers.

### MIN and MAX: Smallest and Largest

```sql
-- What is the smallest and largest single order quantity?
SELECT
    MIN(quantity) AS smallest_order,
    MAX(quantity) AS largest_order
FROM orders;
-- Result: smallest = 1, largest = 10
```

### Combining Aggregates

You can use multiple aggregates in a single query:

```sql
SELECT
    COUNT(*)          AS total_orders,
    SUM(quantity)     AS total_items,
    AVG(quantity)     AS avg_items_per_order,
    MIN(quantity)     AS min_order,
    MAX(quantity)     AS max_order
FROM orders;
```

### Quick Check

1. What is the difference between `COUNT(*)` and `COUNT(customer_id)`?
2. Write a SQL query to find the total revenue across all orders (hint: you need to join with products to get the price).
3. What does `AVG` return if you apply it to an empty table?

---

## 5. GROUP BY and HAVING

Aggregate functions on their own collapse the entire table into one row. Usually, you want a summary per category: revenue per customer, order count per product, sales per month. This is where `GROUP BY` comes in.

**Real-world analogy:** Imagine a spreadsheet with all your bank transactions. You want to know how much you spent at each store. You "group" the transactions by store name and sum the amounts in each group. `GROUP BY` does exactly this.

### GROUP BY: Summarize by Category

```sql
-- How many orders has each customer placed?
SELECT
    customer_id,
    COUNT(*) AS order_count
FROM orders
GROUP BY customer_id;
```

Result:
```
customer_id | order_count
------------+------------
1           | 2
2           | 2
3           | 1
```

`GROUP BY customer_id` tells the database: "divide the rows into buckets where all rows in a bucket share the same `customer_id`, then apply the aggregate function to each bucket separately."

You can join and then group to get human-readable names:

```sql
SELECT
    c.name           AS customer_name,
    COUNT(o.id)      AS order_count,
    SUM(o.quantity)  AS total_items
FROM customers AS c
LEFT JOIN orders AS o ON c.id = o.customer_id
GROUP BY c.id, c.name
ORDER BY total_items DESC NULLS LAST;
```

Result:
```
customer_name | order_count | total_items
--------------+-------------+------------
Bob Smith     | 2           | 15
Carol Davis   | 1           | 3
Alice Chen    | 2           | 3
```

Notice `GROUP BY c.id, c.name` — when using GROUP BY with a JOIN, you must group by all non-aggregated columns in the SELECT. Here both `c.id` and `c.name` are non-aggregated, so both go in GROUP BY.

### Grouping by Multiple Columns

You can group by more than one column to create finer-grained summaries:

```sql
-- Total items ordered per customer per product
SELECT
    c.name     AS customer_name,
    p.name     AS product_name,
    SUM(o.quantity) AS total_qty
FROM orders AS o
JOIN customers AS c ON o.customer_id = c.id
JOIN products  AS p ON o.product_id  = p.id
GROUP BY c.id, c.name, p.id, p.name
ORDER BY c.name, total_qty DESC;
```

### HAVING: Filtering Groups

`WHERE` filters individual rows before grouping. `HAVING` filters groups after aggregation.

The key rule to remember:
- `WHERE` asks "which rows should I include?"
- `HAVING` asks "which groups should I keep?"

**Real-world analogy:** In your bank transaction summary, WHERE removes transactions before you sum them (e.g., "ignore transactions before January"). HAVING filters the final summaries (e.g., "only show stores where I spent more than $100").

```sql
-- Find only the customers who have placed more than one order
SELECT
    c.name           AS customer_name,
    COUNT(o.id)      AS order_count
FROM customers AS c
LEFT JOIN orders AS o ON c.id = o.customer_id
GROUP BY c.id, c.name
HAVING COUNT(o.id) > 1;
```

Result:
```
customer_name | order_count
--------------+------------
Alice Chen    | 2
Bob Smith     | 2
```

Carol has only 1 order and is excluded. Notice we cannot write `HAVING order_count > 1` — we must repeat the aggregate expression `COUNT(o.id) > 1` because `order_count` is just an alias that does not exist yet when HAVING is evaluated.

### WHERE and HAVING Together

```sql
-- Among orders placed in January 2024,
-- find customers who ordered more than 3 items total
SELECT
    c.name            AS customer_name,
    SUM(o.quantity)   AS jan_items
FROM orders AS o
JOIN customers AS c ON o.customer_id = c.id
WHERE o.created_at >= '2024-01-01'
  AND o.created_at  < '2024-02-01'
GROUP BY c.id, c.name
HAVING SUM(o.quantity) > 3;
```

The WHERE clause first narrows the orders to January. Then GROUP BY groups the remaining rows. Then HAVING filters the groups to those with more than 3 items.

### Quick Check

1. What is the key difference between WHERE and HAVING?
2. If you GROUP BY two columns, what does each resulting row represent?
3. Can you use a column alias defined in SELECT inside a HAVING clause in PostgreSQL?

---

## 6. Subqueries: Queries Inside Queries

A **subquery** is a query nested inside another query. The inner query runs first and its result is used by the outer query.

**Real-world analogy:** "Give me the names of all students whose grade is above the class average." To answer this, you first compute the average (inner query), then find students above it (outer query). Subqueries let you express this kind of two-step reasoning in a single SQL statement.

### Subquery in WHERE

```sql
-- Find all orders where the quantity is above the average quantity
SELECT customer_id, product_id, quantity
FROM orders
WHERE quantity > (SELECT AVG(quantity) FROM orders);
```

The database runs `SELECT AVG(quantity) FROM orders` first — this returns `4.2`. Then it runs the outer query: `WHERE quantity > 4.2`. Only the order with quantity 5 and the order with quantity 10 match.

### Subquery in FROM (Derived Tables)

You can use a subquery as if it were a table in the FROM clause:

```sql
-- Find the product with the highest total quantity ordered
SELECT product_name, total_qty
FROM (
    SELECT
        p.name        AS product_name,
        SUM(o.quantity) AS total_qty
    FROM orders AS o
    JOIN products AS p ON o.product_id = p.id
    GROUP BY p.id, p.name
) AS product_totals
ORDER BY total_qty DESC
LIMIT 1;
```

The inner query (everything between the outer parentheses) produces a temporary table named `product_totals`. The outer query then runs against that temporary table as if it were a real table.

### EXISTS Subquery

`EXISTS` checks whether a subquery returns at least one row. It does not return any data from the subquery — just true or false.

```sql
-- Find customers who have placed at least one order
SELECT name
FROM customers
WHERE EXISTS (
    SELECT 1
    FROM orders
    WHERE orders.customer_id = customers.id
);
```

`SELECT 1` inside EXISTS is a convention — since we only care whether any rows exist, we return the constant 1 to avoid fetching real data. The `customers.id` in the subquery refers to the outer query's current row — this is called a **correlated subquery** because the inner query depends on the outer query's current row.

### IN Subquery

```sql
-- Find customers who have ordered product id 1 (Widget A)
SELECT name
FROM customers
WHERE id IN (
    SELECT DISTINCT customer_id
    FROM orders
    WHERE product_id = 1
);
```

The `IN` subquery returns a list of values, and the outer WHERE checks if `id` is in that list.

---

## 7. Common Table Expressions (CTEs) with WITH

As subqueries get more complex, queries become hard to read. You end up with multiple levels of parentheses and it is difficult to follow what is happening. **Common Table Expressions (CTEs)** solve this by letting you name subqueries and refer to them by name.

**Real-world analogy:** When writing a long report, instead of writing a complex formula every time you need the quarterly average, you define it once at the top ("let Q_AVG = the average of all quarterly sales") and then reference Q_AVG throughout. CTEs are exactly this — define it once, use it everywhere.

### Basic CTE Syntax

```sql
WITH cte_name AS (
    -- your subquery here
)
SELECT ...
FROM cte_name
...;
```

Let us rewrite our product totals query using a CTE:

```sql
WITH product_totals AS (
    SELECT
        p.id          AS product_id,
        p.name        AS product_name,
        SUM(o.quantity) AS total_qty,
        SUM(o.quantity * p.price) AS total_revenue
    FROM orders AS o
    JOIN products AS p ON o.product_id = p.id
    GROUP BY p.id, p.name
)
SELECT
    product_name,
    total_qty,
    total_revenue
FROM product_totals
ORDER BY total_revenue DESC;
```

This is much more readable. The CTE `product_totals` acts like a temporary named table that only exists for the duration of this query.

### Multiple CTEs

You can define multiple CTEs by separating them with commas:

```sql
WITH
customer_stats AS (
    SELECT
        customer_id,
        COUNT(*)        AS order_count,
        SUM(quantity)   AS total_items
    FROM orders
    GROUP BY customer_id
),
high_value_customers AS (
    SELECT customer_id
    FROM customer_stats
    WHERE total_items >= 5
)
SELECT
    c.name,
    cs.order_count,
    cs.total_items
FROM customers AS c
JOIN customer_stats        AS cs  ON c.id = cs.customer_id
JOIN high_value_customers  AS hvc ON c.id = hvc.customer_id;
```

The second CTE (`high_value_customers`) can reference the first CTE (`customer_stats`). The main query can reference both.

### Recursive CTEs

CTEs can even reference themselves — this is called a **recursive CTE** and it is used for hierarchical data like organization charts, file system trees, or networks.

```sql
-- Imagine a table: employees(id, name, manager_id)
-- Find all employees under manager id 1, to any depth.
WITH RECURSIVE org_tree AS (
    -- Base case: start with the manager themselves
    SELECT id, name, manager_id, 0 AS depth
    FROM employees
    WHERE id = 1

    UNION ALL

    -- Recursive case: find direct reports of everyone we already found
    SELECT e.id, e.name, e.manager_id, ot.depth + 1
    FROM employees AS e
    JOIN org_tree   AS ot ON e.manager_id = ot.id
)
SELECT name, depth
FROM org_tree
ORDER BY depth, name;
```

`UNION ALL` combines the base case and the recursive case. The database keeps applying the recursive step until no new rows are added.

---

## 8. Window Functions: ROW_NUMBER, RANK, LAG, LEAD

Aggregate functions collapse many rows into one. But what if you want to compute a value per row that takes other rows into account — without collapsing anything?

**Real-world analogy:** A teacher wants to give each student their grade AND their rank within the class — all in a single table. An aggregate function would collapse the class to one row (the average). A window function keeps every student row but adds a rank column computed from all the rows.

This is what **window functions** do. The "window" is the set of rows each calculation can "see."

### The OVER Clause

Every window function uses an `OVER(...)` clause that defines the window:

```sql
function_name() OVER (
    PARTITION BY some_column   -- divide rows into groups (optional)
    ORDER BY     some_column   -- order within each group (optional)
)
```

### ROW_NUMBER: Sequential Number Within a Group

```sql
-- Number each customer's orders chronologically
SELECT
    c.name             AS customer_name,
    o.created_at,
    o.quantity,
    ROW_NUMBER() OVER (
        PARTITION BY o.customer_id
        ORDER BY o.created_at
    ) AS order_number
FROM orders AS o
JOIN customers AS c ON o.customer_id = c.id;
```

Result:
```
customer_name | created_at | quantity | order_number
--------------+------------+----------+-------------
Alice Chen    | 2024-01-10 | 2        | 1
Alice Chen    | 2024-01-15 | 1        | 2
Bob Smith     | 2024-01-20 | 5        | 1
Bob Smith     | 2024-02-01 | 10       | 2
Carol Davis   | 2024-02-05 | 3        | 1
```

`PARTITION BY o.customer_id` tells the database: restart the row number counter for each customer. Order 1 and order 2 for Alice are numbered 1 and 2. Bob's orders independently start from 1 again.

Without `PARTITION BY`, `ROW_NUMBER()` would number all rows sequentially from 1 to 5.

### RANK: Rank with Ties

`RANK()` is like ROW_NUMBER but handles ties — two rows with the same value get the same rank, and the next rank skips ahead.

```sql
-- Rank products by total quantity sold (ties get same rank)
WITH product_totals AS (
    SELECT
        p.name          AS product_name,
        SUM(o.quantity) AS total_qty
    FROM orders AS o
    JOIN products AS p ON o.product_id = p.id
    GROUP BY p.id, p.name
)
SELECT
    product_name,
    total_qty,
    RANK() OVER (ORDER BY total_qty DESC) AS rank
FROM product_totals;
```

If two products had the same total quantity, they would both get rank 2, and the next product would get rank 4 (skipping rank 3). `DENSE_RANK()` is similar but does not skip — it would give rank 2, 2, 3 instead of 2, 2, 4.

### LAG and LEAD: Look at Previous and Next Rows

`LAG(column, n)` returns the value of `column` from `n` rows behind the current row.
`LEAD(column, n)` returns the value from `n` rows ahead.

**Real-world analogy:** In a list of monthly sales figures, for each month you want to know: what was last month's revenue (LAG) and what will next month's be (LEAD)? This lets you calculate month-over-month growth.

```sql
-- For each order per customer, show the previous order's quantity for comparison
SELECT
    c.name           AS customer_name,
    o.created_at,
    o.quantity       AS this_order,
    LAG(o.quantity, 1) OVER (
        PARTITION BY o.customer_id
        ORDER BY o.created_at
    )                AS previous_order
FROM orders AS o
JOIN customers AS c ON o.customer_id = c.id;
```

Result:
```
customer_name | created_at | this_order | previous_order
--------------+------------+------------+---------------
Alice Chen    | 2024-01-10 | 2          | NULL
Alice Chen    | 2024-01-15 | 1          | 2
Bob Smith     | 2024-01-20 | 5          | NULL
Bob Smith     | 2024-02-01 | 10         | 5
Carol Davis   | 2024-02-05 | 3          | NULL
```

The first order for each customer shows NULL because there is no previous order to look at.

### Running Totals with SUM OVER

Window functions also work with standard aggregates like SUM:

```sql
-- Show a running total of items ordered over time
SELECT
    created_at,
    quantity,
    SUM(quantity) OVER (ORDER BY created_at) AS running_total
FROM orders;
```

Result:
```
created_at | quantity | running_total
-----------+----------+--------------
2024-01-10 | 2        | 2
2024-01-15 | 1        | 3
2024-01-20 | 5        | 8
2024-02-01 | 10       | 18
2024-02-05 | 3        | 21
```

Each row shows the cumulative total up to and including that row. Notice the rows are not collapsed — each original row is preserved with its running total appended.

### Quick Check

1. What is the key difference between an aggregate function and a window function?
2. What does `PARTITION BY` do inside an OVER clause?
3. What does LAG return for the very first row in a partition?

---

## 9. Building Complex Queries in Go with sqlx and pgx

Now let us bring all these SQL concepts into Go. We will look at two popular libraries: `sqlx` (a thin wrapper over the standard `database/sql` that adds conveniences) and `pgx` (a full-featured PostgreSQL driver written specifically for Go).

### Setting Up the Project

```bash
mkdir ecommerce-queries
cd ecommerce-queries
go mod init ecommerce-queries
go get github.com/jmoiron/sqlx
go get github.com/lib/pq
go get github.com/jackc/pgx/v5
```

### Approach 1: Using sqlx

`sqlx` adds one very useful feature over the standard library: it can scan query results directly into Go structs using field tags, saving you from writing long `rows.Scan(...)` calls.

```go
package main

import (
	"fmt"
	"log"

	// sqlx wraps the standard database/sql with extra convenience methods.
	"github.com/jmoiron/sqlx"

	// The pq package is the PostgreSQL driver.
	// The underscore import runs its init() function, which registers "postgres"
	// as a driver name with database/sql. We never call pq functions directly.
	_ "github.com/lib/pq"
)

// OrderSummary holds one row from our analytics query.
// The `db:"..."` struct tags tell sqlx which SQL column maps to which field.
// The field names in Go are capitalized (exported); the SQL columns are lowercase.
type OrderSummary struct {
	CustomerName string  `db:"customer_name"`
	OrderCount   int     `db:"order_count"`
	TotalItems   int     `db:"total_items"`
	TotalRevenue float64 `db:"total_revenue"`
}

func main() {
	// Connect to PostgreSQL.
	// The connection string format is: "postgres://user:password@host:port/dbname?sslmode=..."
	// sslmode=disable means we are not using TLS — fine for local development.
	dsn := "postgres://postgres:secret@localhost:5432/ecommerce?sslmode=disable"

	// sqlx.Open is like sql.Open but returns a *sqlx.DB instead.
	// It does not actually connect yet — it just sets up the driver.
	db, err := sqlx.Open("postgres", dsn)
	if err != nil {
		log.Fatalf("failed to open database: %v", err)
	}
	// defer schedules db.Close() to run when main() returns.
	// This ensures the connection is always closed, even if we return early.
	defer db.Close()

	// Ping actually opens a connection and checks the database is reachable.
	if err := db.Ping(); err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}
	fmt.Println("Connected to PostgreSQL.")

	// Our CTE-powered analytics query.
	// The query is written as a raw string literal (backtick string) for readability.
	// It uses a CTE to compute per-customer stats, then joins with the customers table.
	query := `
		WITH customer_stats AS (
			SELECT
				o.customer_id,
				COUNT(o.id)                       AS order_count,
				SUM(o.quantity)                    AS total_items,
				SUM(o.quantity * p.price)          AS total_revenue
			FROM orders AS o
			JOIN products AS p ON o.product_id = p.id
			GROUP BY o.customer_id
		)
		SELECT
			c.name                AS customer_name,
			cs.order_count,
			cs.total_items,
			COALESCE(cs.total_revenue, 0) AS total_revenue
		FROM customers AS c
		LEFT JOIN customer_stats AS cs ON c.id = cs.customer_id
		ORDER BY cs.total_revenue DESC NULLS LAST
	`

	// db.Select scans all rows from the query into a slice of structs.
	// This is a sqlx feature — standard database/sql does not have this.
	// It reads every row, calls Scan on each one, and appends to the slice.
	var summaries []OrderSummary
	if err := db.Select(&summaries, query); err != nil {
		log.Fatalf("query failed: %v", err)
	}

	// Print the results in a formatted table.
	fmt.Printf("%-15s %10s %12s %14s\n",
		"Customer", "Orders", "Items Sold", "Revenue ($)")
	fmt.Println("----------------------------------------------------")
	for _, s := range summaries {
		fmt.Printf("%-15s %10d %12d %14.2f\n",
			s.CustomerName, s.OrderCount, s.TotalItems, s.TotalRevenue)
	}
}
```

Let us walk through what is happening line by line:

1. We define a `OrderSummary` struct with `db:` tags. When `sqlx.Select` runs, it matches each column from the SQL result to a struct field by looking at these tags.
2. `sqlx.Open` creates a database handle. No connection is made yet.
3. `db.Ping()` actually connects and verifies the database is available.
4. The SQL query uses a CTE named `customer_stats` to compute per-customer aggregates. The outer query joins these stats back to the customers table.
5. `COALESCE(cs.total_revenue, 0)` handles the LEFT JOIN case — if a customer has no orders, `cs.total_revenue` is NULL. COALESCE returns its first non-NULL argument, so this returns 0 for customers with no orders.
6. `db.Select(&summaries, query)` runs the query and populates our slice. The `&` passes a pointer so `db.Select` can write into our variable.

### Approach 2: Using pgx Directly

`pgx` is a more powerful PostgreSQL-specific driver. It supports all PostgreSQL features that `database/sql` does not — like streaming large result sets, COPY, and PostgreSQL-specific types.

```go
package main

import (
	"context"
	"fmt"
	"log"

	// pgx/v5 is the PostgreSQL driver.
	// pgxpool manages a pool of reusable connections.
	"github.com/jackc/pgx/v5/pgxpool"
)

// WindowResult holds one row from our window function query.
type WindowResult struct {
	CustomerName  string
	OrderDate     string
	Quantity      int
	OrderNumber   int
	RunningTotal  int
}

func main() {
	// context.Background() creates a root context with no deadline.
	// Almost all pgx operations require a context — this lets you cancel
	// long-running queries if needed (e.g., if an HTTP request is cancelled).
	ctx := context.Background()

	// pgxpool.New creates a connection pool.
	// A pool keeps multiple connections open and reuses them across goroutines.
	// This is more efficient than opening a new connection for every query.
	dsn := "postgres://postgres:secret@localhost:5432/ecommerce?sslmode=disable"
	pool, err := pgxpool.New(ctx, dsn)
	if err != nil {
		log.Fatalf("failed to create connection pool: %v", err)
	}
	defer pool.Close()

	// pool.Ping verifies at least one connection can be established.
	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("failed to connect: %v", err)
	}
	fmt.Println("Connected via pgx pool.")

	// A window function query: number each customer's orders and show running totals.
	query := `
		SELECT
			c.name                              AS customer_name,
			o.created_at::DATE::TEXT            AS order_date,
			o.quantity,
			ROW_NUMBER() OVER (
				PARTITION BY o.customer_id
				ORDER BY o.created_at
			)::INT                              AS order_number,
			SUM(o.quantity) OVER (
				PARTITION BY o.customer_id
				ORDER BY o.created_at
			)::INT                              AS running_total
		FROM orders AS o
		JOIN customers AS c ON o.customer_id = c.id
		ORDER BY c.name, o.created_at
	`

	// pool.Query executes the query and returns a Rows cursor.
	// The cursor is lazy — it does not fetch all rows at once.
	// We iterate through rows one by one with rows.Next().
	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}
	// rows.Close() releases the connection back to the pool.
	// defer ensures this always happens even if we return early.
	defer rows.Close()

	fmt.Printf("%-12s %-12s %8s %12s %13s\n",
		"Customer", "Date", "Qty", "Order#", "RunningTotal")
	fmt.Println("--------------------------------------------------------------")

	// rows.Next() advances the cursor to the next row.
	// It returns true if there is a row to read, false when all rows are consumed.
	for rows.Next() {
		var r WindowResult

		// rows.Scan reads the current row's columns into Go variables in order.
		// The order of arguments must match the order of columns in the SELECT.
		err := rows.Scan(
			&r.CustomerName,
			&r.OrderDate,
			&r.Quantity,
			&r.OrderNumber,
			&r.RunningTotal,
		)
		if err != nil {
			log.Fatalf("scan failed: %v", err)
		}

		fmt.Printf("%-12s %-12s %8d %12d %13d\n",
			r.CustomerName, r.OrderDate, r.Quantity, r.OrderNumber, r.RunningTotal)
	}

	// rows.Err() returns any error that occurred during iteration.
	// Always check this after the loop — errors during Next() do not panic.
	if err := rows.Err(); err != nil {
		log.Fatalf("rows iteration error: %v", err)
	}
}
```

### Parameterized Queries: Preventing SQL Injection

**SQL injection** is a security vulnerability where user-provided input is inserted directly into a SQL string, allowing an attacker to run arbitrary SQL. Always use parameterized queries.

```go
package main

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5/pgxpool"
)

func main() {
	ctx := context.Background()
	pool, _ := pgxpool.New(ctx, "postgres://postgres:secret@localhost:5432/ecommerce?sslmode=disable")
	defer pool.Close()

	// WRONG — Never do this. If customerName came from user input like
	// "'; DROP TABLE orders; --", this would destroy the orders table.
	// customerName := "Alice"
	// badQuery := "SELECT * FROM customers WHERE name = '" + customerName + "'"

	// CORRECT — Use $1, $2, ... placeholders.
	// pgx replaces the placeholder with the actual value safely,
	// ensuring the input can never be interpreted as SQL.
	customerName := "Alice Chen"
	query := `
		SELECT
			c.name,
			COUNT(o.id)       AS order_count,
			SUM(o.quantity)   AS total_items
		FROM customers AS c
		LEFT JOIN orders AS o ON c.id = o.customer_id
		WHERE c.name = $1
		GROUP BY c.id, c.name
	`

	// The extra arguments after the query string are the parameter values.
	// pgx substitutes $1 with customerName, safely escaped.
	rows, err := pool.Query(ctx, query, customerName)
	if err != nil {
		log.Fatalf("query failed: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var name string
		var orderCount, totalItems int
		if err := rows.Scan(&name, &orderCount, &totalItems); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%s: %d orders, %d items\n", name, orderCount, totalItems)
	}
}
```

`$1` is PostgreSQL's placeholder syntax. `database/sql` and `sqlx` use `?` for SQLite and MySQL, and `$1` for PostgreSQL. `pgx` always uses `$1`, `$2`, etc.

---

## 10. Mini Project: Analytics Dashboard for an E-Commerce Database

Let us put everything together. We will build a Go program that runs an analytics dashboard for an e-commerce database. The program connects to PostgreSQL and prints five analytics reports, each using different SQL techniques from this chapter.

### The Schema

First, run this SQL to set up a richer test database in PostgreSQL:

```sql
-- Drop and recreate for a clean slate
DROP TABLE IF EXISTS order_items;
DROP TABLE IF EXISTS orders;
DROP TABLE IF EXISTS products;
DROP TABLE IF EXISTS categories;
DROP TABLE IF EXISTS customers;

CREATE TABLE categories (
    id   SERIAL PRIMARY KEY,
    name TEXT NOT NULL
);

CREATE TABLE products (
    id          SERIAL PRIMARY KEY,
    name        TEXT NOT NULL,
    price       NUMERIC(10, 2) NOT NULL,
    category_id INTEGER REFERENCES categories(id)
);

CREATE TABLE customers (
    id         SERIAL PRIMARY KEY,
    name       TEXT NOT NULL,
    email      TEXT NOT NULL UNIQUE,
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE TABLE orders (
    id          SERIAL PRIMARY KEY,
    customer_id INTEGER NOT NULL REFERENCES customers(id),
    created_at  TIMESTAMP DEFAULT NOW(),
    status      TEXT NOT NULL DEFAULT 'completed'
);

CREATE TABLE order_items (
    id         SERIAL PRIMARY KEY,
    order_id   INTEGER NOT NULL REFERENCES orders(id),
    product_id INTEGER NOT NULL REFERENCES products(id),
    quantity   INTEGER NOT NULL,
    unit_price NUMERIC(10, 2) NOT NULL  -- price at time of purchase
);

-- Seed data
INSERT INTO categories (name) VALUES ('Electronics'), ('Kitchen'), ('Books');

INSERT INTO products (name, price, category_id) VALUES
    ('Wireless Headphones', 79.99, 1),
    ('Phone Stand',          9.99, 1),
    ('Coffee Maker',        49.99, 2),
    ('Blender',             34.99, 2),
    ('Go Programming',      39.99, 3),
    ('Database Internals',  44.99, 3);

INSERT INTO customers (name, email) VALUES
    ('Alice Chen',   'alice@example.com'),
    ('Bob Smith',    'bob@example.com'),
    ('Carol Davis',  'carol@example.com'),
    ('Dave Wilson',  'dave@example.com');

-- Orders with items
INSERT INTO orders (customer_id, created_at, status) VALUES
    (1, '2024-01-05', 'completed'),
    (1, '2024-02-10', 'completed'),
    (2, '2024-01-20', 'completed'),
    (2, '2024-03-15', 'completed'),
    (3, '2024-02-28', 'completed'),
    (3, '2024-03-20', 'completed'),
    (4, '2024-01-10', 'completed');

INSERT INTO order_items (order_id, product_id, quantity, unit_price) VALUES
    (1, 1, 1, 79.99),  -- Alice: Headphones
    (1, 2, 2,  9.99),  -- Alice: Phone Stand x2
    (2, 5, 1, 39.99),  -- Alice: Go book
    (3, 3, 1, 49.99),  -- Bob: Coffee Maker
    (3, 4, 1, 34.99),  -- Bob: Blender
    (4, 1, 2, 79.99),  -- Bob: Headphones x2
    (5, 6, 1, 44.99),  -- Carol: Database book
    (5, 5, 1, 39.99),  -- Carol: Go book
    (6, 2, 3,  9.99),  -- Carol: Phone Stand x3
    (7, 3, 1, 49.99);  -- Dave: Coffee Maker
```

### The Analytics Dashboard in Go

```go
package main

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/jackc/pgx/v5/pgxpool"
)

// --------------------------------------------------------------------------
// Data structures for each report
// --------------------------------------------------------------------------

// RevenueByCategory holds one row of the revenue-by-category report.
type RevenueByCategory struct {
	CategoryName string
	ProductCount int
	TotalOrders  int
	TotalRevenue float64
}

// TopCustomer holds one row of the top-customers report.
type TopCustomer struct {
	Rank         int
	CustomerName string
	OrderCount   int
	TotalSpent   float64
}

// ProductRank holds one row of the product-ranking-within-category report.
type ProductRank struct {
	CategoryName string
	ProductName  string
	Revenue      float64
	RankInCat    int
}

// MonthlyTrend holds one row of the month-over-month revenue trend.
type MonthlyTrend struct {
	Month        string
	Revenue      float64
	PrevRevenue  float64
	GrowthPct    float64
}

// CustomerSegment holds one row of the customer segmentation report.
type CustomerSegment struct {
	CustomerName string
	TotalSpent   float64
	Segment      string
}

// --------------------------------------------------------------------------
// Main program
// --------------------------------------------------------------------------

func main() {
	ctx := context.Background()

	pool, err := pgxpool.New(ctx, "postgres://postgres:secret@localhost:5432/ecommerce?sslmode=disable")
	if err != nil {
		log.Fatalf("failed to create pool: %v", err)
	}
	defer pool.Close()

	if err := pool.Ping(ctx); err != nil {
		log.Fatalf("cannot reach database: %v", err)
	}

	printHeader("ECOMMERCE ANALYTICS DASHBOARD")

	// Run each report in sequence.
	// Each function handles its own query and printing.
	reportRevenueByCategory(ctx, pool)
	reportTopCustomers(ctx, pool)
	reportProductRanksWithinCategory(ctx, pool)
	reportMonthlyTrend(ctx, pool)
	reportCustomerSegments(ctx, pool)
}

// --------------------------------------------------------------------------
// Report 1: Revenue by Category (GROUP BY + JOIN + aggregate)
// --------------------------------------------------------------------------

func reportRevenueByCategory(ctx context.Context, pool *pgxpool.Pool) {
	printHeader("Report 1: Revenue by Product Category")

	// This query joins four tables:
	// order_items -> products -> categories
	// It uses GROUP BY to aggregate per category.
	query := `
		SELECT
			cat.name                              AS category_name,
			COUNT(DISTINCT p.id)                  AS product_count,
			COUNT(DISTINCT oi.order_id)           AS total_orders,
			SUM(oi.quantity * oi.unit_price)      AS total_revenue
		FROM order_items AS oi
		JOIN products    AS p   ON oi.product_id  = p.id
		JOIN categories  AS cat ON p.category_id  = cat.id
		GROUP BY cat.id, cat.name
		ORDER BY total_revenue DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("report 1 query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-15s %14s %13s %14s\n",
		"Category", "Products", "Orders", "Revenue ($)")
	fmt.Println(strings.Repeat("-", 60))

	for rows.Next() {
		var r RevenueByCategory
		if err := rows.Scan(&r.CategoryName, &r.ProductCount,
			&r.TotalOrders, &r.TotalRevenue); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-15s %14d %13d %14.2f\n",
			r.CategoryName, r.ProductCount, r.TotalOrders, r.TotalRevenue)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

// --------------------------------------------------------------------------
// Report 2: Top Customers by Spending (CTE + RANK window function)
// --------------------------------------------------------------------------

func reportTopCustomers(ctx context.Context, pool *pgxpool.Pool) {
	printHeader("Report 2: Top Customers by Total Spending")

	// This query uses a CTE to compute total spend per customer,
	// then uses the RANK() window function to rank them.
	query := `
		WITH customer_spend AS (
			SELECT
				c.id,
				c.name                              AS customer_name,
				COUNT(DISTINCT o.id)                AS order_count,
				SUM(oi.quantity * oi.unit_price)    AS total_spent
			FROM customers AS c
			JOIN orders      AS o  ON c.id         = o.customer_id
			JOIN order_items AS oi ON o.id          = oi.order_id
			GROUP BY c.id, c.name
		)
		SELECT
			RANK() OVER (ORDER BY total_spent DESC)::INT AS rank,
			customer_name,
			order_count,
			total_spent
		FROM customer_spend
		ORDER BY rank
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("report 2 query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%6s %-15s %10s %14s\n",
		"Rank", "Customer", "Orders", "Total Spent ($)")
	fmt.Println(strings.Repeat("-", 50))

	for rows.Next() {
		var r TopCustomer
		if err := rows.Scan(&r.Rank, &r.CustomerName,
			&r.OrderCount, &r.TotalSpent); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%6d %-15s %10d %14.2f\n",
			r.Rank, r.CustomerName, r.OrderCount, r.TotalSpent)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

// --------------------------------------------------------------------------
// Report 3: Product Rankings Within Each Category (RANK PARTITION BY)
// --------------------------------------------------------------------------

func reportProductRanksWithinCategory(ctx context.Context, pool *pgxpool.Pool) {
	printHeader("Report 3: Product Rankings Within Each Category")

	// RANK() OVER (PARTITION BY ...) restarts the rank counter for each category.
	// This lets us say "the #1 product in Electronics, the #1 in Kitchen, etc."
	query := `
		WITH product_revenue AS (
			SELECT
				cat.name                            AS category_name,
				p.name                              AS product_name,
				SUM(oi.quantity * oi.unit_price)    AS revenue
			FROM order_items AS oi
			JOIN products    AS p   ON oi.product_id = p.id
			JOIN categories  AS cat ON p.category_id = cat.id
			GROUP BY cat.id, cat.name, p.id, p.name
		)
		SELECT
			category_name,
			product_name,
			revenue,
			RANK() OVER (
				PARTITION BY category_name
				ORDER BY revenue DESC
			)::INT AS rank_in_category
		FROM product_revenue
		ORDER BY category_name, rank_in_category
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("report 3 query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-15s %-25s %12s %6s\n",
		"Category", "Product", "Revenue ($)", "Rank")
	fmt.Println(strings.Repeat("-", 62))

	for rows.Next() {
		var r ProductRank
		if err := rows.Scan(&r.CategoryName, &r.ProductName,
			&r.Revenue, &r.RankInCat); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-15s %-25s %12.2f %6d\n",
			r.CategoryName, r.ProductName, r.Revenue, r.RankInCat)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

// --------------------------------------------------------------------------
// Report 4: Month-over-Month Revenue Trend (LAG window function)
// --------------------------------------------------------------------------

func reportMonthlyTrend(ctx context.Context, pool *pgxpool.Pool) {
	printHeader("Report 4: Month-over-Month Revenue Trend")

	// TO_CHAR formats a timestamp as a string like "2024-01".
	// LAG(revenue, 1) looks back one row to get the previous month's revenue.
	// NULLIF prevents division-by-zero when prev_revenue is 0.
	// ROUND rounds the growth percentage to 1 decimal place.
	query := `
		WITH monthly AS (
			SELECT
				TO_CHAR(o.created_at, 'YYYY-MM')       AS month,
				SUM(oi.quantity * oi.unit_price)        AS revenue
			FROM orders      AS o
			JOIN order_items AS oi ON o.id = oi.order_id
			GROUP BY TO_CHAR(o.created_at, 'YYYY-MM')
			ORDER BY month
		)
		SELECT
			month,
			revenue,
			COALESCE(LAG(revenue, 1) OVER (ORDER BY month), 0)  AS prev_revenue,
			COALESCE(
				ROUND(
					(revenue - LAG(revenue, 1) OVER (ORDER BY month))
					/ NULLIF(LAG(revenue, 1) OVER (ORDER BY month), 0)
					* 100,
					1
				),
				0
			)                                                    AS growth_pct
		FROM monthly
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("report 4 query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-10s %14s %16s %12s\n",
		"Month", "Revenue ($)", "Prev Month ($)", "Growth %")
	fmt.Println(strings.Repeat("-", 55))

	for rows.Next() {
		var r MonthlyTrend
		if err := rows.Scan(&r.Month, &r.Revenue,
			&r.PrevRevenue, &r.GrowthPct); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-10s %14.2f %16.2f %11.1f%%\n",
			r.Month, r.Revenue, r.PrevRevenue, r.GrowthPct)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

// --------------------------------------------------------------------------
// Report 5: Customer Segmentation (subquery + CASE)
// --------------------------------------------------------------------------

func reportCustomerSegments(ctx context.Context, pool *pgxpool.Pool) {
	printHeader("Report 5: Customer Segmentation by Spending")

	// CASE WHEN is SQL's if/else expression — it returns different values
	// based on conditions, like a switch statement in Go.
	// The subquery computes total spend per customer, then the outer query
	// assigns a segment label based on spending thresholds.
	query := `
		SELECT
			customer_name,
			total_spent,
			CASE
				WHEN total_spent >= 150 THEN 'VIP'
				WHEN total_spent >= 75  THEN 'Regular'
				ELSE                         'Occasional'
			END AS segment
		FROM (
			SELECT
				c.name                              AS customer_name,
				COALESCE(SUM(oi.quantity * oi.unit_price), 0) AS total_spent
			FROM customers   AS c
			LEFT JOIN orders      AS o  ON c.id  = o.customer_id
			LEFT JOIN order_items AS oi ON o.id   = oi.order_id
			GROUP BY c.id, c.name
		) AS spending
		ORDER BY total_spent DESC
	`

	rows, err := pool.Query(ctx, query)
	if err != nil {
		log.Fatalf("report 5 query failed: %v", err)
	}
	defer rows.Close()

	fmt.Printf("%-15s %14s %12s\n",
		"Customer", "Total Spent ($)", "Segment")
	fmt.Println(strings.Repeat("-", 44))

	for rows.Next() {
		var r CustomerSegment
		if err := rows.Scan(&r.CustomerName, &r.TotalSpent,
			&r.Segment); err != nil {
			log.Fatal(err)
		}
		fmt.Printf("%-15s %14.2f %12s\n",
			r.CustomerName, r.TotalSpent, r.Segment)
	}
	if err := rows.Err(); err != nil {
		log.Fatal(err)
	}
	fmt.Println()
}

// --------------------------------------------------------------------------
// Helper: print a section header
// --------------------------------------------------------------------------

func printHeader(title string) {
	fmt.Println(strings.Repeat("=", 62))
	fmt.Printf("  %s\n", title)
	fmt.Println(strings.Repeat("=", 62))
}
```

### Expected Output

When you run `go run main.go` against the seeded database, you will see output like:

```
==============================================================
  ECOMMERCE ANALYTICS DASHBOARD
==============================================================
==============================================================
  Report 1: Revenue by Product Category
==============================================================
Category         Products       Orders     Revenue ($)
------------------------------------------------------------
Electronics             2            3          219.93
Kitchen                 2            3          184.96
Books                   2            3          124.97

==============================================================
  Report 2: Top Customers by Total Spending
==============================================================
  Rank Customer           Orders Total Spent ($)
--------------------------------------------------
     1 Bob Smith               2         224.95
     2 Alice Chen              2         109.96
     3 Carol Davis             2          94.97
     4 Dave Wilson             1          49.99

==============================================================
  Report 3: Product Rankings Within Each Category
==============================================================
Category        Product                   Revenue ($)   Rank
--------------------------------------------------------------
Books           Go Programming               79.98        1
Books           Database Internals           44.99        2
Electronics     Wireless Headphones         239.97        1
Electronics     Phone Stand                  29.97        2
Kitchen         Coffee Maker                 99.98        1
Kitchen         Blender                      34.99        2

==============================================================
  Report 4: Month-over-Month Revenue Trend
==============================================================
Month       Revenue ($)   Prev Month ($)   Growth %
-------------------------------------------------------
2024-01          209.93           0.00         0.0%
2024-02          174.94         209.93       -16.7%
2024-03          144.97         174.94       -17.1%

==============================================================
  Report 5: Customer Segmentation by Spending
==============================================================
Customer    Total Spent ($)      Segment
--------------------------------------------
Bob Smith            224.95          VIP
Alice Chen           109.96      Regular
Carol Davis           94.97      Regular
Dave Wilson           49.99   Occasional
```

Every report in this dashboard uses a different combination of techniques from this chapter: JOINs, aggregations, GROUP BY, CTEs, window functions, and subqueries.

---

## Summary

- **Normalization** splits data into multiple tables to eliminate redundancy. Each fact is stored in exactly one place, with **foreign keys** linking related rows across tables.

- **JOINs** combine rows from multiple tables. INNER JOIN returns only matched rows; LEFT JOIN keeps all rows from the left table (with NULLs for unmatched right rows); RIGHT JOIN mirrors this; FULL JOIN keeps everything.

- **Aggregate functions** (COUNT, SUM, AVG, MIN, MAX) collapse many rows into a single summary value. Combined with **GROUP BY**, they produce one summary per group. **HAVING** filters which groups appear in the final result — unlike WHERE, which filters individual rows before grouping.

- **Subqueries** nest one query inside another, allowing two-step reasoning in a single statement. **Common Table Expressions (CTEs)** with the `WITH` keyword give subqueries readable names, making complex queries much easier to understand.

- **Window functions** (ROW_NUMBER, RANK, LAG, LEAD, SUM OVER) compute values that depend on multiple rows without collapsing them. The `OVER (PARTITION BY ... ORDER BY ...)` clause defines which rows each calculation sees.

---

## Exercises

### Easy

1. Using the e-commerce schema from the Mini Project, write a query that returns the name and price of every product in the "Books" category. Use a JOIN between `products` and `categories`.

2. Write a query that counts how many orders each customer has placed. Include customers who have placed zero orders (hint: LEFT JOIN). Order the results by order count descending.

3. Using aggregate functions, find the most expensive product, the cheapest product, and the average price across all products in a single query.

### Medium

4. Write a query that finds the top three best-selling products by total units sold (quantity). Use a JOIN between `order_items` and `products`, GROUP BY the product, and LIMIT 3.

5. Rewrite the following subquery using a CTE instead. The original query finds all customers who have spent more than $100 total:
   ```sql
   SELECT name
   FROM customers
   WHERE id IN (
       SELECT customer_id
       FROM orders AS o
       JOIN order_items AS oi ON o.id = oi.order_id
       GROUP BY customer_id
       HAVING SUM(oi.quantity * oi.unit_price) > 100
   );
   ```

6. Write a window function query that, for each order, shows the order's revenue AND the running total revenue for that customer up to and including that order (ordered by `created_at`). Each row should show: customer name, order date, order revenue, and the running total.

### Hard

7. Write a query that returns each product along with the date of its most recent order and how many days ago that was (as of today). Products that have never been ordered should show NULL for the date and NULL for days ago. Use `CURRENT_DATE` to get today's date in PostgreSQL.

8. Build a cohort retention query: for each calendar month (cohort), show how many customers placed their first-ever order in that month, and how many of those customers also placed at least one order in the following month. (Hint: use window functions to find each customer's first order month, then LEFT JOIN the cohort data against orders in the next month.)

9. Write a full Go program using `pgx` that accepts a product category name as a command-line argument (using `os.Args`) and prints a report for that category only: each product in the category with its total revenue, rank within the category, and percentage of the category's total revenue. Use parameterized queries to prevent SQL injection. The output should be formatted as a table.
