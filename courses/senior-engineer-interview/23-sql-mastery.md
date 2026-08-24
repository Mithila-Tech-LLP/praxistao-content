# Chapter 23: SQL Mastery — Complex Joins, Window Functions & CTEs

SQL is tested in every senior backend interview. The expectation is that you can write complex queries without looking things up, explain execution plans, and discuss performance tradeoffs. This chapter covers the SQL patterns that appear in interviews at top companies.

## Table of Contents

1. [JOIN Types Explained](#1-join-types)
2. [Aggregation and GROUP BY](#2-aggregation)
3. [Window Functions](#3-window-functions)
4. [Common Table Expressions (CTEs)](#4-common-table-expressions-ctes)
5. [Subqueries vs JOINs vs CTEs](#5-subqueries-vs-joins-vs-ctes)
6. [10 Interview-Level SQL Problems](#6-interview-level-sql-problems)
7. [Summary](#summary)

---

## 1. JOIN Types

```sql
-- Sample schema:
-- employees(id, name, department_id, salary, manager_id)
-- departments(id, name, budget)
-- projects(id, name)
-- project_assignments(employee_id, project_id, start_date)

-- INNER JOIN: only rows with matches in BOTH tables
SELECT e.name, d.name as department
FROM employees e
INNER JOIN departments d ON e.department_id = d.id;
-- Excludes employees with no department, and departments with no employees

-- LEFT JOIN: all rows from left table, matching rows from right (NULL if no match)
SELECT e.name, d.name as department
FROM employees e
LEFT JOIN departments d ON e.department_id = d.id;
-- Includes employees with no department (department columns will be NULL)

-- RIGHT JOIN: all rows from right table (rarely used — just swap tables and use LEFT JOIN)

-- FULL OUTER JOIN: all rows from both tables
SELECT e.name, d.name as department
FROM employees e
FULL OUTER JOIN departments d ON e.department_id = d.id;
-- Includes employees with no dept AND departments with no employees

-- CROSS JOIN: cartesian product (every row of A × every row of B)
SELECT e.name, p.name
FROM employees e
CROSS JOIN projects p;
-- Every employee paired with every project

-- SELF JOIN: join a table with itself
-- Find each employee and their manager's name
SELECT e.name as employee, m.name as manager
FROM employees e
LEFT JOIN employees m ON e.manager_id = m.id;
```

### Finding Records That DON'T Match

```sql
-- Employees with no department assignment (LEFT JOIN + IS NULL)
SELECT e.name
FROM employees e
LEFT JOIN departments d ON e.department_id = d.id
WHERE d.id IS NULL;

-- Same result with NOT EXISTS (often faster):
SELECT e.name
FROM employees e
WHERE NOT EXISTS (
    SELECT 1 FROM departments d WHERE d.id = e.department_id
);

-- Same result with NOT IN (careful with NULLs!):
-- NOT IN with NULL values behaves unexpectedly — use NOT EXISTS instead
```

---

## 2. Aggregation

```sql
-- COUNT(*) vs COUNT(column): COUNT(*) counts all rows, COUNT(col) skips NULLs
SELECT 
    COUNT(*) as total_employees,
    COUNT(department_id) as employees_with_dept,  -- excludes NULLs
    COUNT(DISTINCT department_id) as unique_depts
FROM employees;

-- GROUP BY: one output row per unique combination of grouped columns
SELECT department_id, COUNT(*) as emp_count, AVG(salary) as avg_salary
FROM employees
GROUP BY department_id;

-- HAVING: filter groups (like WHERE but for aggregated results)
SELECT department_id, COUNT(*) as emp_count
FROM employees
GROUP BY department_id
HAVING COUNT(*) > 5;  -- only departments with more than 5 employees

-- ORDER OF OPERATIONS: FROM → WHERE → GROUP BY → HAVING → SELECT → ORDER BY
-- You CANNOT use a SELECT alias in WHERE or HAVING (they execute before SELECT)
```

---

## 3. Window Functions

Window functions compute values across a "window" of related rows without collapsing them into a single row like GROUP BY does.

```sql
-- Syntax: function() OVER (PARTITION BY ... ORDER BY ...)

-- ROW_NUMBER: unique sequential number within each partition
SELECT 
    name, 
    salary,
    department_id,
    ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC) as rank_in_dept
FROM employees;
-- Row 1 in each department is the highest-paid employee

-- RANK vs DENSE_RANK vs ROW_NUMBER:
-- Employees: [100k, 100k, 80k]
-- ROW_NUMBER: 1, 2, 3 (always unique)
-- RANK:       1, 1, 3 (gap after ties)
-- DENSE_RANK: 1, 1, 2 (no gap after ties)

SELECT name, salary,
    ROW_NUMBER() OVER (ORDER BY salary DESC) as row_num,
    RANK()       OVER (ORDER BY salary DESC) as rank,
    DENSE_RANK() OVER (ORDER BY salary DESC) as dense_rank
FROM employees;

-- LAG/LEAD: access previous/next row's values
SELECT 
    name,
    salary,
    LAG(salary)  OVER (ORDER BY hire_date) as prev_salary,  -- previous row's salary
    LEAD(salary) OVER (ORDER BY hire_date) as next_salary    -- next row's salary
FROM employees;

-- Running total (SUM with window)
SELECT 
    name,
    salary,
    SUM(salary) OVER (ORDER BY hire_date ROWS UNBOUNDED PRECEDING) as running_total
FROM employees;

-- Moving average (SUM over last 3 rows)
SELECT 
    name,
    salary,
    AVG(salary) OVER (ORDER BY hire_date ROWS BETWEEN 2 PRECEDING AND CURRENT ROW) as moving_avg_3
FROM employees;

-- NTILE: divide into N buckets
SELECT name, salary,
    NTILE(4) OVER (ORDER BY salary) as quartile
FROM employees;
-- Quartile 1 = bottom 25%, Quartile 4 = top 25%
```

---

## 4. Common Table Expressions (CTEs)

CTEs make complex queries readable by giving subqueries a name.

```sql
-- Basic CTE: name a subquery
WITH high_earners AS (
    SELECT * FROM employees WHERE salary > 100000
),
dept_counts AS (
    SELECT department_id, COUNT(*) as cnt FROM high_earners GROUP BY department_id
)
SELECT d.name, dc.cnt
FROM departments d
JOIN dept_counts dc ON d.id = dc.department_id
ORDER BY dc.cnt DESC;

-- Recursive CTE: for hierarchical data
-- Find all reports in a management chain
WITH RECURSIVE org_chart AS (
    -- Base case: start with top-level managers (no manager)
    SELECT id, name, manager_id, 0 as level
    FROM employees
    WHERE manager_id IS NULL

    UNION ALL

    -- Recursive case: find direct reports of employees already in the CTE
    SELECT e.id, e.name, e.manager_id, oc.level + 1
    FROM employees e
    JOIN org_chart oc ON e.manager_id = oc.id
)
SELECT name, level FROM org_chart ORDER BY level, name;
```

---

## 5. Subqueries vs JOINs vs CTEs

```sql
-- Same result, different styles:

-- Subquery in WHERE:
SELECT name FROM employees
WHERE department_id IN (SELECT id FROM departments WHERE budget > 1000000);

-- JOIN equivalent:
SELECT DISTINCT e.name FROM employees e
JOIN departments d ON e.department_id = d.id
WHERE d.budget > 1000000;

-- CTE equivalent:
WITH rich_depts AS (SELECT id FROM departments WHERE budget > 1000000)
SELECT e.name FROM employees e
WHERE e.department_id IN (SELECT id FROM rich_depts);
```

**When to use which:**
- **Subquery in SELECT:** scalar value per row (correlated subquery). Can be slow if executed per row.
- **Subquery in WHERE/FROM:** when you need to filter or compute before joining.
- **JOIN:** when you need columns from both tables, or when the optimizer can plan better.
- **CTE:** when you reference the same subquery multiple times, or for readability in complex queries.

---

## 6. Interview-Level SQL Problems

### Problem 1: Second Highest Salary

```sql
-- Use OFFSET to skip the first row
SELECT MAX(salary) as second_highest
FROM employees
WHERE salary < (SELECT MAX(salary) FROM employees);

-- Or with window function (handles ties correctly):
SELECT salary
FROM (
    SELECT salary, DENSE_RANK() OVER (ORDER BY salary DESC) as rnk
    FROM employees
) ranked
WHERE rnk = 2
LIMIT 1;
```

### Problem 2: Employees Who Earn More Than Their Manager

```sql
SELECT e.name as employee
FROM employees e
JOIN employees m ON e.manager_id = m.id
WHERE e.salary > m.salary;
```

### Problem 3: Consecutive Rows with Same Value (Gaps and Islands)

```sql
-- Find consecutive login days for each user
WITH ranked AS (
    SELECT user_id, login_date,
        ROW_NUMBER() OVER (PARTITION BY user_id ORDER BY login_date) as rn
    FROM logins
),
groups AS (
    SELECT user_id, login_date,
        DATE_SUB(login_date, INTERVAL rn DAY) as grp  -- same value = consecutive days
    FROM ranked
)
SELECT user_id, MIN(login_date) as start_date, MAX(login_date) as end_date,
    COUNT(*) as consecutive_days
FROM groups
GROUP BY user_id, grp
HAVING COUNT(*) >= 3
ORDER BY user_id, start_date;
```

### Problem 4: Top N Per Group

```sql
-- Top 3 highest-paid employees per department
SELECT name, salary, department_id
FROM (
    SELECT name, salary, department_id,
        ROW_NUMBER() OVER (PARTITION BY department_id ORDER BY salary DESC) as rn
    FROM employees
) ranked
WHERE rn <= 3;
```

### Problem 5: Running Total with Reset

```sql
-- Running total of daily sales, reset each month
SELECT 
    sale_date,
    amount,
    SUM(amount) OVER (
        PARTITION BY DATE_TRUNC('month', sale_date)  -- reset each month
        ORDER BY sale_date
        ROWS UNBOUNDED PRECEDING
    ) as monthly_running_total
FROM sales;
```

---

## Summary

- **INNER JOIN:** rows matching in both tables. **LEFT JOIN:** all left rows + matched right rows.
- **CROSS JOIN:** cartesian product. **SELF JOIN:** join table to itself (hierarchies, comparisons).
- For rows with NO match: use `LEFT JOIN ... WHERE right.col IS NULL` or `NOT EXISTS`.
- **Window functions:** compute across a window of rows without collapsing. Key functions: ROW_NUMBER, RANK, DENSE_RANK, LAG, LEAD, SUM, AVG with OVER clause.
- **CTEs:** name a subquery with `WITH`. Use recursive CTEs for hierarchical data.
- WHERE filters rows before grouping; HAVING filters groups after aggregation.
