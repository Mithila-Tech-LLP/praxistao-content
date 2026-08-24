# Chapter 01: Python — Your First Programming Language

*This chapter assumes you have never written a line of code. If you have, skim Section 1-4 and start from Section 5.*

---

## Table of Contents

1. What Is a Programming Language?
2. Why Python?
3. Installing Python — Step by Step
4. Your Very First Program
5. Variables — Giving Things Names
6. Numbers and Maths
7. Text (Strings)
8. True or False (Booleans)
9. Making Decisions — If/Else
10. Doing Things Over and Over — Loops
11. Organising Code — Functions
12. Lists — Groups of Things
13. Dictionaries — Labelled Boxes
14. Classes — Your Own Types
15. Reading and Writing Files
16. Installing Libraries
17. Virtual Environments
18. Jupyter Notebooks
19. Common Mistakes and How to Fix Them
20. Summary and Exercises

---

## 1. What Is a Programming Language?

A programming language is a way to give instructions to a computer.

Computers only understand one language natively: electricity. A tiny transistor is either on (1) or off (0). Everything a computer does comes down to sequences of 1s and 0s.

Writing instructions in 1s and 0s is impossible for humans. So, decades ago, people invented programming languages: a way to write instructions that looks more like English, which a special program (called a **compiler** or **interpreter**) then translates into 1s and 0s.

Think of it like this: your recipe book is written in English. A robot that can follow recipes only understands robot-language. The interpreter is the person who stands in the middle, reading your English recipe aloud and translating it into robot-language in real time.

Python is one of hundreds of programming languages. But it is the one used almost universally in AI and machine learning, and it is widely considered the easiest one to learn.

---

## 2. Why Python?

Here is why Python became the language of AI:

**It reads like English.** Compare these two programs that do the same thing:

```
JAVA (harder to read):
for (int i = 0; i < 10; i++) {
    System.out.println(i);
}

PYTHON (easier to read):
for i in range(10):
    print(i)
```

Python is designed to look like pseudocode — the informal descriptions programmers write on whiteboards. This makes it faster to write and easier to read.

**It has the best libraries for AI.** A library is a collection of pre-written code you can use. NumPy, PyTorch, TensorFlow, scikit-learn — every major AI tool is built for Python. Other languages have weaker or no equivalents.

**It is free and runs everywhere.** Windows, Mac, Linux — all work. No licence fees, no subscriptions.

---

## 3. Installing Python — Step by Step

### Step 1: Download Python

Go to **python.org** and click the big "Download Python" button. As of 2025, download Python 3.11 or newer.

### Step 2: Install

**On Windows:**
1. Run the downloaded `.exe` file
2. **IMPORTANT:** Tick "Add Python to PATH" before clicking Install
3. Click "Install Now"
4. Wait for it to finish

**On Mac:**
1. Run the downloaded `.pkg` file
2. Follow the prompts (click Continue, Agree, Install)

**On Linux (Ubuntu/Debian):**
```bash
sudo apt update
sudo apt install python3 python3-pip
```

### Step 3: Verify it worked

Open your **terminal** (called "Command Prompt" on Windows, "Terminal" on Mac/Linux). Type:

```bash
python --version
```

You should see something like:
```
Python 3.11.7
```

If you see an error, try:
```bash
python3 --version
```

### Step 4: Open the interactive Python shell

Type `python` (or `python3`) in your terminal and press Enter. You will see:

```
Python 3.11.7 (default, ...)
>>>
```

The `>>>` is Python waiting for you to type something. This is called the **REPL** (Read-Eval-Print Loop). Type things, press Enter, see results. Try:

```python
>>> 2 + 2
4
>>> print("Hello!")
Hello!
```

Type `exit()` to leave.

---

## 4. Your Very First Program

Create a new file called `hello.py`. (The `.py` extension tells your computer it is a Python file.)

You can use any text editor. On Windows, Notepad works. On Mac, TextEdit works (but set it to plain text mode). Better options: VS Code (free, download at code.visualstudio.com) or PyCharm.

Type this into the file:

```python
print("Hello, world!")
print("I am learning Python.")
print("One day I will build my own AI.")
```

Save the file. Open your terminal, navigate to the folder where you saved it, and run:

```bash
python hello.py
```

You should see:
```
Hello, world!
I am learning Python.
One day I will build my own AI.
```

**Congratulations.** You just ran your first program.

`print()` is a **function** — a command that tells Python to display something. Whatever you put inside the parentheses gets shown on screen.

---

## 5. Variables — Giving Things Names

A variable is a named container for a piece of information.

```python
name = "Alice"
age = 15
height = 1.65
is_student = True
```

Think of `name = "Alice"` as: "create a box, label it `name`, put the value `"Alice"` inside."

```
  name        age       height    is_student
 ┌──────┐   ┌─────┐   ┌──────┐   ┌──────┐
 │"Alice"│   │ 15  │   │ 1.65 │   │ True │
 └──────┘   └─────┘   └──────┘   └──────┘
```

You can use variables anywhere you'd use the value:

```python
name = "Alice"
age = 15
print(name)          # prints: Alice
print(age + 1)       # prints: 16
print("My name is", name)   # prints: My name is Alice
```

Variable naming rules:
- Must start with a letter or underscore: `name`, `_count`, `my_data` ✓
- Cannot start with a number: `1thing` ✗
- No spaces: `my name` ✗, use `my_name` instead
- Case-sensitive: `Name` and `name` are different variables

---

## 6. Numbers and Maths

Python has two main types of numbers:

```python
# Integers (whole numbers)
x = 42
y = -7
big = 1_000_000    # underscores make big numbers readable

# Floats (decimal numbers)
pi = 3.14159
temperature = 98.6
```

Maths operations:

```python
print(10 + 3)    # 13   (addition)
print(10 - 3)    # 7    (subtraction)
print(10 * 3)    # 30   (multiplication)
print(10 / 3)    # 3.333...  (division — always gives a float)
print(10 // 3)   # 3    (integer division — rounds down)
print(10 % 3)    # 1    (remainder, called "modulo")
print(2 ** 10)   # 1024 (2 to the power 10)
```

This matters for ML because neural networks are essentially doing billions of these operations.

---

## 7. Text (Strings)

A string is a sequence of characters — any text.

```python
greeting = "Hello"
name = 'Alice'      # single quotes work too

# Combining strings
message = greeting + ", " + name + "!"
print(message)      # Hello, Alice!

# A better way: f-strings (f for "format")
message = f"Hello, {name}! You are {age} years old."
print(message)      # Hello, Alice! You are 15 years old.

# String methods (built-in tools for strings)
sentence = "The quick brown fox"
print(sentence.upper())     # THE QUICK BROWN FOX
print(sentence.lower())     # the quick brown fox
print(sentence.split())     # ['The', 'quick', 'brown', 'fox']
print(len(sentence))        # 19  (number of characters)
print(sentence[0])          # T   (first character)
print(sentence[4:9])        # quic  (characters 4 to 8)
```

Strings are enormously important in AI because all text data — the training data for language models, the prompts you send to ChatGPT, the responses you get back — is made of strings.

---

## 8. True or False (Booleans)

A boolean is a value that is either `True` or `False`. Named after mathematician George Boole.

```python
is_raining = True
is_sunny = False

# Comparisons produce booleans
print(5 > 3)     # True
print(5 < 3)     # False
print(5 == 5)    # True  (== means "is equal to")
print(5 != 3)    # True  (!= means "is not equal to")
print(5 >= 5)    # True  (>= means "greater than or equal to")
```

Combining booleans with `and`, `or`, `not`:

```python
age = 17
has_id = False

can_buy_alcohol = age >= 18 and has_id
print(can_buy_alcohol)   # False (both must be True for AND)

can_vote = age >= 16 or age >= 18
print(can_vote)          # True (only one needs to be True for OR)

not_raining = not is_raining
print(not_raining)       # False
```

In ML, we constantly check conditions: "if the model's accuracy is above 95%, stop training."

---

## 9. Making Decisions — If/Else

```python
temperature = 25

if temperature > 30:
    print("It is hot!")
elif temperature > 20:
    print("It is warm.")
elif temperature > 10:
    print("It is cool.")
else:
    print("It is cold!")

# Output: It is warm.
```

**Important:** Python uses **indentation** (spaces at the start of a line) to show which code belongs inside an `if`. The standard is 4 spaces. This is not just style — Python will crash if your indentation is wrong.

```python
# WRONG:
if True:
print("this will crash")   # not indented — Python error!

# RIGHT:
if True:
    print("this works")    # 4 spaces of indentation
```

---

## 10. Doing Things Over and Over — Loops

### The `for` loop

A `for` loop repeats code a set number of times, or for each item in a group.

```python
# Count from 0 to 4
for i in range(5):
    print(i)
# Output: 0  1  2  3  4  (each on its own line)

# Loop over a list of items
fruits = ["apple", "banana", "cherry"]
for fruit in fruits:
    print(f"I like {fruit}")
# Output:
# I like apple
# I like banana
# I like cherry
```

### The `while` loop

A `while` loop repeats code as long as a condition is True.

```python
count = 0
while count < 5:
    print(count)
    count = count + 1   # or: count += 1
# Output: 0  1  2  3  4
```

**Warning:** If the condition never becomes False, the loop runs forever. This is called an **infinite loop**. If your program freezes, press Ctrl+C.

Loops are fundamental to ML — training a model means looping over your data thousands of times.

---

## 11. Organising Code — Functions

A function is a named block of code you can run whenever you want, as many times as you want.

```python
# Defining a function
def greet(name):
    message = f"Hello, {name}!"
    return message

# Calling a function
result = greet("Alice")
print(result)         # Hello, Alice!
print(greet("Bob"))   # Hello, Bob!
```

Functions can take multiple inputs (called **parameters** or **arguments**):

```python
def add(a, b):
    return a + b

print(add(3, 4))    # 7
print(add(10, 5))   # 15
```

Functions can have default values:

```python
def greet(name, greeting="Hello"):
    return f"{greeting}, {name}!"

print(greet("Alice"))            # Hello, Alice!
print(greet("Alice", "Hi"))      # Hi, Alice!
```

Functions are crucial in ML for reusing code — "train the model", "evaluate the model", "plot the results" are all functions you will write.

---

## 12. Lists — Groups of Things

A list is an ordered collection of items, all in one variable.

```python
scores = [95, 87, 92, 78, 88]

print(scores[0])    # 95   (first item — Python counts from 0!)
print(scores[-1])   # 88   (last item)
print(scores[1:3])  # [87, 92]  (items at index 1 and 2)
print(len(scores))  # 5   (how many items)

scores.append(91)   # add 91 to the end
print(scores)       # [95, 87, 92, 78, 88, 91]

scores.sort()
print(scores)       # [78, 87, 88, 91, 92, 95]
```

A **list comprehension** creates a new list by processing each item:

```python
numbers = [1, 2, 3, 4, 5]
squares = [n ** 2 for n in numbers]
print(squares)   # [1, 4, 9, 16, 25]

# with a filter
evens = [n for n in numbers if n % 2 == 0]
print(evens)     # [2, 4]
```

In ML, datasets are essentially lists of examples.

---

## 13. Dictionaries — Labelled Boxes

A dictionary stores key-value pairs. Instead of looking things up by position number (like a list), you look things up by name (the key).

```python
person = {
    "name": "Alice",
    "age": 15,
    "city": "London"
}

print(person["name"])   # Alice
print(person["age"])    # 15

person["email"] = "alice@example.com"   # add a new key
print(person)
# {'name': 'Alice', 'age': 15, 'city': 'London', 'email': 'alice@example.com'}

# Check if a key exists
if "age" in person:
    print("Age found:", person["age"])

# Loop through a dictionary
for key, value in person.items():
    print(f"{key}: {value}")
```

In ML, you will use dictionaries constantly — for storing model configuration, metrics, hyperparameters.

---

## 14. Classes — Your Own Types

A class is a blueprint for creating objects. Think of a class like a cookie-cutter — it defines the shape. Objects are the actual cookies.

```python
class Dog:
    def __init__(self, name, breed):
        self.name = name
        self.breed = breed

    def bark(self):
        return f"{self.name} says: Woof!"

    def describe(self):
        return f"{self.name} is a {self.breed}"

# Create objects (instances) from the class
dog1 = Dog("Rex", "German Shepherd")
dog2 = Dog("Bella", "Labrador")

print(dog1.bark())        # Rex says: Woof!
print(dog2.describe())    # Bella is a Labrador
print(dog1.name)          # Rex
```

The `__init__` method (with double underscores on each side) runs automatically when you create a new object. It sets up the object's starting state.

`self` refers to the object itself — it is Python's way of saying "this specific dog."

In ML, every model is a class. When you write `model = LinearRegression()`, you are creating an object from the `LinearRegression` class.

---

## 15. Reading and Writing Files

```python
# Writing to a file
with open("data.txt", "w") as f:
    f.write("Line 1\n")
    f.write("Line 2\n")
    f.write("Line 3\n")

# Reading from a file
with open("data.txt", "r") as f:
    contents = f.read()
    print(contents)

# Reading line by line
with open("data.txt", "r") as f:
    for line in f:
        print(line.strip())   # .strip() removes the newline character at the end
```

The `with` keyword ensures the file is properly closed when you are done, even if an error occurs.

---

## 16. Installing Libraries

Python comes with many built-in tools, but for ML you need extra libraries. Install them with `pip` (Python's package installer):

```bash
# Basic ML libraries
pip install numpy pandas matplotlib seaborn scikit-learn

# Deep learning
pip install torch

# AI APIs
pip install anthropic openai
```

Then use them in your code:

```python
import numpy as np
import pandas as pd
import matplotlib.pyplot as plt
```

The `as` part gives the library a short nickname. `np` for NumPy and `pd` for Pandas are standard conventions — everyone uses them, so use them too.

---

## 17. Virtual Environments

A virtual environment is a separate Python installation just for one project. This prevents different projects from interfering with each other.

```bash
# Create a virtual environment
python -m venv myproject

# Activate it (Mac/Linux)
source myproject/bin/activate

# Activate it (Windows)
myproject\Scripts\activate

# Your terminal prompt will change, showing you're inside the env
# Install packages — they only go into this environment
pip install numpy pandas torch

# Deactivate when done
deactivate
```

**Always use a virtual environment for each ML project.** This is not optional advice — it will save you enormous headaches later.

---

## 18. Jupyter Notebooks

A Jupyter notebook is an interactive document where you can mix code, text, and charts. It is the most popular tool for ML exploration.

```bash
# Install Jupyter
pip install jupyter notebook

# Start Jupyter
jupyter notebook
```

Your browser will open a file browser. Click "New" → "Python 3" to create a notebook.

A notebook is made of **cells**. Each cell contains either code or text. Press **Shift+Enter** to run a cell.

```
┌────────────────────────────────────────┐
│ [1]: # This is a code cell              │
│       x = 5                             │
│       print(x)                          │
├────────────────────────────────────────┤
│ 5                                       │  ← output appears here
├────────────────────────────────────────┤
│ ## This is a text cell (Markdown)       │
│ You can write **bold** and _italic_     │
└────────────────────────────────────────┘
```

Jupyter is where you will do all your ML experiments in this course.

---

## 19. Common Mistakes and How to Fix Them

### Mistake 1: Indentation errors

```python
# WRONG
def greet(name):
print("Hello", name)   # IndentationError

# RIGHT
def greet(name):
    print("Hello", name)  # 4 spaces
```

### Mistake 2: Forgetting colons

```python
# WRONG
if x > 5
    print(x)   # SyntaxError: missing colon

# RIGHT
if x > 5:
    print(x)
```

### Mistake 3: Off-by-one errors with indexing

```python
items = ["a", "b", "c"]
print(items[3])   # IndexError: list index out of range
# Lists go from index 0 to len-1. Last item is items[2], not items[3].
```

### Mistake 4: == vs =

```python
x = 5     # assignment: set x to 5
x == 5    # comparison: is x equal to 5? (produces True or False)

# WRONG (common mistake)
if x = 5:  # SyntaxError! You used = instead of ==
    ...

# RIGHT
if x == 5:
    ...
```

### Reading error messages

Python error messages always tell you what went wrong and where. Read them from the bottom up — the last line is the most important.

```
Traceback (most recent call last):
  File "hello.py", line 5, in <module>    ← look here: line 5
    print(namee)                           ← the exact line
NameError: name 'namee' is not defined    ← the error: typo in variable name
```

---

## Summary

You now know:
- **Variables** store data
- **Strings** are text
- **Numbers** are integers or floats
- **Booleans** are True or False
- **If/else** makes decisions
- **Loops** repeat actions
- **Functions** reuse code
- **Lists** hold ordered collections
- **Dictionaries** hold key-value pairs
- **Classes** create custom objects
- **pip** installs libraries
- **Virtual environments** isolate projects
- **Jupyter** is where ML happens

---

## Exercises

**Easy:**

1. Write a program that prints the numbers 1 to 100. (Hint: `range(1, 101)`)

2. Write a function `is_even(n)` that returns `True` if `n` is even, `False` otherwise. Test it on 4, 7, 0, -3.

3. Create a dictionary for a book: title, author, year, pages. Print each value using a for loop.

**Medium:**

4. Write a program that counts how many times each word appears in this sentence: `"the cat sat on the mat the cat"`. Store the counts in a dictionary. (Hint: use `sentence.split()` and `dict.get(word, 0)`)

5. Create a `BankAccount` class with methods `deposit(amount)`, `withdraw(amount)`, and `get_balance()`. Make sure withdrawals fail if there is not enough money.

6. Write a function that takes a list of numbers and returns: the minimum, maximum, sum, and average. Return them all as a dictionary.

**Hard:**

7. Write a program that reads a text file (create one with any content), counts the words, and prints the 5 most common words and their counts. (Hint: you will need the `sorted()` function with a `key` argument.)

8. Create a class `SimpleDatabase` that can `insert(key, value)`, `get(key)`, and `delete(key)`. Make it persist data by writing to and reading from a file (you can use JSON format with `import json`).
