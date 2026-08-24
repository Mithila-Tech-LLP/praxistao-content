# Chapter 02: NumPy — The Engine of Numerical Computing

> **"NumPy is the bedrock on which the entire scientific Python ecosystem is built. Without it, there is no pandas, no scikit-learn, no PyTorch."**

---

## Table of Contents
1. [Why NumPy? The Speed Problem](#1-why-numpy-the-speed-problem)
2. [The ndarray: Core Data Structure](#2-the-ndarray-core-data-structure)
3. [Creating Arrays](#3-creating-arrays)
4. [dtypes in Depth](#4-dtypes-in-depth)
5. [Indexing and Slicing](#5-indexing-and-slicing)
6. [Boolean and Fancy Indexing](#6-boolean-and-fancy-indexing)
7. [Broadcasting](#7-broadcasting)
8. [Vectorized Operations vs Loops](#8-vectorized-operations-vs-loops)
9. [Mathematical Operations](#9-mathematical-operations)
10. [Aggregations with Axis](#10-aggregations-with-axis)
11. [Linear Algebra](#11-linear-algebra)
12. [Reshaping and Transposing](#12-reshaping-and-transposing)
13. [Stacking and Splitting](#13-stacking-and-splitting)
14. [Useful Utility Functions](#14-useful-utility-functions)
15. [Random Module and Reproducibility](#15-random-module-and-reproducibility)
16. [Memory Layout](#16-memory-layout)
17. [Saving and Loading](#17-saving-and-loading)
18. [Mini Project: Linear Regression from Scratch](#18-mini-project-linear-regression-from-scratch)
19. [Summary](#19-summary)
20. [Exercises](#20-exercises)

---

## 1. Why NumPy? The Speed Problem

Python is a wonderfully expressive language, but it has a fundamental performance problem for numerical computing: **everything is a Python object**.

When you write `x = 5`, Python doesn't store just the number 5 in a memory location. It creates a Python object with a reference count, type information, and then the actual value. Every single integer, every float, every loop iteration has this overhead.

```
Python list element:              NumPy array element:
──────────────────────────────    ──────────────────────────────
┌──────────────────────────┐      ┌──────────────┐
│ Python Object            │      │ 4-byte float │  ← just the data!
│  - Reference count (8B)  │      └──────────────┘
│  - Type pointer (8B)     │
│  - Value (8B)            │      Array of 1M floats:
└──────────────────────────┘      1M × 4 bytes = 4 MB
  24 bytes per float!
  Array of 1M floats:
  1M × 24 bytes = 24 MB
```

NumPy arrays store data in **contiguous blocks of raw memory** — just like C arrays. Operations on them call highly optimized C and FORTRAN code (Intel MKL, OpenBLAS) that uses SIMD instructions to process multiple elements simultaneously.

```python
import numpy as np
import time

n = 10_000_000

# Pure Python loop
python_list = list(range(n))
start = time.time()
result_py = [x * 2 for x in python_list]
py_time = time.time() - start

# NumPy vectorized
numpy_array = np.arange(n)
start = time.time()
result_np = numpy_array * 2
np_time = time.time() - start

print(f"Python: {py_time:.3f}s")
print(f"NumPy:  {np_time:.3f}s")
print(f"Speedup: {py_time / np_time:.0f}x faster")
# Typical output:
# Python: 0.832s
# NumPy:  0.012s
# Speedup: 69x faster
```

---

## 2. The ndarray: Core Data Structure

The **ndarray** (n-dimensional array) is NumPy's central data structure. It is:
- A fixed-size block of memory
- All elements are the same dtype
- Described by a shape, stride, and dtype

```
ndarray anatomy:
─────────────────────────────────────────────────────
Shape:   (3, 4)      ← dimensions in each axis
Dtype:   float32     ← type of each element
Ndim:    2           ← number of dimensions (axes)
Size:    12          ← total number of elements
Itemsize: 4          ← bytes per element
Nbytes:  48          ← total bytes = size × itemsize
Strides: (16, 4)     ← bytes to skip to move 1 step along each axis
─────────────────────────────────────────────────────
```

```python
import numpy as np

arr = np.array([[1.0, 2.0, 3.0, 4.0],
                [5.0, 6.0, 7.0, 8.0],
                [9.0, 10., 11., 12.]], dtype=np.float32)

print(f"shape:    {arr.shape}")      # (3, 4)
print(f"dtype:    {arr.dtype}")      # float32
print(f"ndim:     {arr.ndim}")       # 2
print(f"size:     {arr.size}")       # 12
print(f"itemsize: {arr.itemsize}")   # 4 bytes
print(f"nbytes:   {arr.nbytes}")     # 48 bytes
print(f"strides:  {arr.strides}")    # (16, 4)
```

### Understanding Strides

Strides are how NumPy navigates memory. For a (3, 4) float32 array:
- Moving along axis 1 (columns): jump 4 bytes (1 float32)
- Moving along axis 0 (rows): jump 16 bytes (4 floats = 1 full row)

```
Memory layout (row-major / C-order):
──────────────────────────────────────────────────────
Index: [0,0] [0,1] [0,2] [0,3] [1,0] [1,1] [1,2] [1,3] ...
Value:  1.0   2.0   3.0   4.0   5.0   6.0   7.0   8.0  ...
Byte:   0     4     8    12    16    20    24    28   ...
       ←── row 0 ──────────→ ←── row 1 ──────────→
```

---

## 3. Creating Arrays

```python
import numpy as np

# ── From Python data ────────────────────────────────────────────────────
a = np.array([1, 2, 3, 4])                    # 1D from list
b = np.array([[1, 2], [3, 4]], dtype=float)   # 2D, specify dtype
c = np.array([[[1,2],[3,4]],[[5,6],[7,8]]])   # 3D

# ── Filled arrays ───────────────────────────────────────────────────────
zeros  = np.zeros((3, 4))         # all zeros, float64 by default
ones   = np.ones((2, 3), dtype=np.float32)
full   = np.full((3, 3), 7.0)    # fill with specific value
empty  = np.empty((100, 100))    # UNINITIALIZED memory — fastest, but unpredictable values!
# Note: np.empty is used when you'll overwrite every element anyway

# ── Like another array ──────────────────────────────────────────────────
template = np.array([1.0, 2.0, 3.0])
z_like   = np.zeros_like(template)   # same shape and dtype, filled with 0
o_like   = np.ones_like(template)    # same shape and dtype, filled with 1

# ── Range arrays ────────────────────────────────────────────────────────
# np.arange: like Python range(), returns ndarray
r1 = np.arange(10)             # [0, 1, 2, ..., 9]
r2 = np.arange(0, 1, 0.1)     # [0., 0.1, 0.2, ..., 0.9]  (step=0.1)
r3 = np.arange(10, 0, -2)     # [10, 8, 6, 4, 2]

# np.linspace: evenly spaced values (inclusive of endpoint)
#   USE LINSPACE when you want N points in a range
#   USE ARANGE when you want a specific step size
ls = np.linspace(0, 1, 11)     # [0.0, 0.1, 0.2, ..., 1.0] — 11 points
ls2 = np.linspace(0, 2*np.pi, 100)   # 100 points from 0 to 2π

# ── Special matrices ────────────────────────────────────────────────────
eye  = np.eye(4)                # 4×4 identity matrix
diag = np.diag([1, 2, 3, 4])   # diagonal matrix from values
diag_from = np.diag(eye)       # extract diagonal from matrix

# ── Random arrays ───────────────────────────────────────────────────────
# (covered in depth in Section 15)
rng = np.random.default_rng(42)   # modern way to create a random number generator
uniform  = rng.random((3, 3))     # uniform [0, 1)
normal   = rng.standard_normal((3, 3))  # N(0, 1)
integers = rng.integers(0, 10, (3, 3))  # random integers
```

---

## 4. dtypes in Depth

Every NumPy array has a dtype — the data type of its elements. Choosing the right dtype is not just theoretical; it directly affects memory usage, speed, and numerical stability.

```
dtype summary:
──────────────────────────────────────────────────────────────────────
Type         Bits   Range                          ML Use
──────────────────────────────────────────────────────────────────────
float64      64     ±1.8×10^308, ~15 decimal dig   Default, CPU training
float32      32     ±3.4×10^38,  ~7 decimal dig    GPU training (PREFERRED)
float16      16     ±65504,      ~3 decimal dig    Mixed-precision training
bfloat16     16     ±3.4×10^38,  ~2 decimal dig    Google TPU, modern GPUs
int64        64     -9.2×10^18 to 9.2×10^18        Large indices, IDs
int32        32     -2.1×10^9  to 2.1×10^9         Labels, most indices
int16        16     -32768 to 32767                Compact storage
int8         8      -128 to 127                    Quantized models
uint8        8      0 to 255                       Image pixel values
bool         8      True/False                     Masks, conditions
──────────────────────────────────────────────────────────────────────
```

```python
import numpy as np

# Memory comparison: same data, different dtypes
data = np.random.randn(1000, 1000)   # 1M elements

f64 = data.astype(np.float64)
f32 = data.astype(np.float32)
f16 = data.astype(np.float16)

print(f"float64: {f64.nbytes / 1e6:.0f} MB")  # 8 MB
print(f"float32: {f32.nbytes / 1e6:.0f} MB")  # 4 MB
print(f"float16: {f16.nbytes / 1e6:.0f} MB")  # 2 MB

# For neural networks with 1 BILLION parameters:
params = 1_000_000_000
print(f"\n1B parameter model:")
print(f"float64: {params * 8 / 1e9:.0f} GB")   # 8 GB
print(f"float32: {params * 4 / 1e9:.0f} GB")   # 4 GB
print(f"float16: {params * 2 / 1e9:.0f} GB")   # 2 GB

# Dtype conversions
x = np.array([1.7, 2.3, 3.9])
print(x.astype(np.int32))     # [1, 2, 3]  — truncates, not rounds!
print(x.astype(np.bool_))     # [True, True, True]  — nonzero → True

# Checking and setting dtype
arr = np.zeros((3, 4))
print(arr.dtype)              # float64 (default)
arr32 = arr.astype(np.float32)
print(arr32.dtype)            # float32
print(arr32.itemsize)         # 4 bytes

# Type check
print(np.issubdtype(arr.dtype, np.floating))   # True
print(np.issubdtype(arr.dtype, np.integer))    # False
```

---

## 5. Indexing and Slicing

Indexing is how you access elements. In ML you'll do this constantly — selecting features, extracting batches, masking data.

### 1D Indexing

```
Array:   [10, 20, 30, 40, 50]
Index:     0   1   2   3   4
Neg idx:  -5  -4  -3  -2  -1

arr[0]   → 10      arr[-1]  → 50
arr[1:3] → [20,30] arr[::2] → [10,30,50]
arr[::-1]→ [50,40,30,20,10]  (reversed)
```

```python
arr = np.array([10, 20, 30, 40, 50])

# Basic indexing
print(arr[0])      # 10
print(arr[-1])     # 50  (last element)
print(arr[-2])     # 40  (second to last)

# Slicing: [start:stop:step]  — stop is EXCLUSIVE
print(arr[1:4])    # [20, 30, 40]  — indices 1, 2, 3
print(arr[::2])    # [10, 30, 50]  — every other element
print(arr[::-1])   # [50, 40, 30, 20, 10]  — reversed

# Important: NumPy slices return VIEWS, not copies!
view = arr[1:3]
view[0] = 999
print(arr)         # [10, 999, 30, 40, 50]  — original changed!
copy = arr[1:3].copy()   # use .copy() to avoid this
```

### 2D Indexing

```
2D array (3 rows, 4 columns):
       col 0  col 1  col 2  col 3
row 0 [  1,    2,    3,    4  ]
row 1 [  5,    6,    7,    8  ]
row 2 [  9,   10,   11,   12  ]

arr[row, col]
arr[1, 2]     → 7
arr[:, 1]     → [2, 6, 10]   (all rows, col 1)
arr[0, :]     → [1, 2, 3, 4] (row 0, all cols)
arr[0:2, 1:3] → [[2,3],[6,7]] (rows 0-1, cols 1-2)
```

```python
arr = np.array([[1, 2, 3, 4],
                [5, 6, 7, 8],
                [9,10,11,12]])

# Single element
print(arr[1, 2])         # 7

# Entire row / column
print(arr[0, :])         # [1 2 3 4]  — row 0
print(arr[:, 1])         # [2 6 10]   — column 1

# Sub-matrix
print(arr[0:2, 1:3])     # [[2,3],[6,7]]

# Last column
print(arr[:, -1])        # [4, 8, 12]

# Every other row
print(arr[::2, :])       # [[1,2,3,4],[9,10,11,12]]

# First two rows, all columns
X_train = arr[:2, :]     # [[1,2,3,4],[5,6,7,8]]
```

### 3D Indexing (Batches of Images)

```
3D array shape: (batch_size, height, width)
                    (4,         3,      3  )

Think of it as 4 images, each 3×3:

Axis 0 = batch index     (which image)
Axis 1 = row             (which row in image)
Axis 2 = column          (which column in image)

arr[0]         → first image (3×3)
arr[0, 1]      → first image, second row (1D of width 3)
arr[0, 1, 2]   → first image, row 1, col 2 (scalar)
arr[:, 0, :]   → first row of ALL images (shape: 4,3)
```

```python
# Typical ML: batch of images
batch = np.random.rand(4, 28, 28)   # 4 images of 28×28
print(batch.shape)   # (4, 28, 28)

first_image   = batch[0]            # shape (28, 28)
pixel         = batch[0, 14, 14]    # center pixel of first image
top_rows      = batch[:, :5, :]     # top 5 rows of all images: (4, 5, 28)
all_centers   = batch[:, 14, 14]    # center pixel of each image: (4,)
```

---

## 6. Boolean and Fancy Indexing

### Boolean Indexing

Boolean indexing lets you select elements based on conditions. This is one of the most powerful and frequently used features in ML data processing.

```python
import numpy as np

ages = np.array([15, 25, 17, 30, 16, 42, 21])
incomes = np.array([0, 45000, 0, 75000, 0, 120000, 35000])

# Create boolean mask
is_adult = ages >= 18                    # [F, T, F, T, F, T, T]
print(is_adult)

# Apply mask — returns only True positions
adult_ages    = ages[is_adult]           # [25, 30, 42, 21]
adult_incomes = incomes[is_adult]        # [45000, 75000, 120000, 35000]

# Compound conditions
rich_adults = incomes[(ages >= 18) & (incomes > 50000)]  # & for and
any_either  = incomes[(ages < 18) | (incomes == 0)]      # | for or
not_adult   = ages[~is_adult]                             # ~ for not

# Modify elements matching a condition (in-place)
data = np.array([1.0, -2.0, 3.0, -4.0, 5.0])
data[data < 0] = 0.0     # ReLU activation! (clip negatives to zero)
print(data)              # [1. 0. 3. 0. 5.]

# Count elements matching condition
n_negatives = np.sum(data < 0)
print(n_negatives)       # 0 (after clipping)
```

### Fancy Indexing

Fancy indexing uses arrays of indices to select non-contiguous elements.

```python
arr = np.array([10, 20, 30, 40, 50])

# Select specific indices
idx = np.array([0, 2, 4])
print(arr[idx])          # [10, 30, 50]

# 2D fancy indexing — select specific rows
X = np.arange(16).reshape(4, 4)
#  [[ 0  1  2  3]
#   [ 4  5  6  7]
#   [ 8  9 10 11]
#   [12 13 14 15]]

row_indices = np.array([0, 2, 3])
print(X[row_indices])    # rows 0, 2, 3

# ML use case: gather a mini-batch by indices
dataset = np.random.randn(10000, 128)  # 10000 samples, 128 features
batch_indices = np.array([42, 91, 7, 305, 1024])   # random batch
batch = dataset[batch_indices]          # shape: (5, 128)

# np.random.choice — sample random indices (without replacement)
random_batch_idx = np.random.choice(len(dataset), size=32, replace=False)
random_batch = dataset[random_batch_idx]
```

---

## 7. Broadcasting

Broadcasting is the rule by which NumPy performs operations on arrays of **different but compatible shapes**. This trips almost everyone up at first, but it's one of the most powerful features.

### The Rules

```
Broadcasting rules (applied left to right, padding with 1s on left):
───────────────────────────────────────────────────────────────────
Rule 1: If arrays have different ndims, pad the smaller shape
        on the LEFT with 1s until ndims match.
        
Rule 2: Arrays are compatible along each dimension if:
        (a) sizes are equal, OR
        (b) one of them is 1.
        
Rule 3: Arrays with size 1 in a dimension are "stretched"
        to match the other array's size in that dimension.
───────────────────────────────────────────────────────────────────
```

### Visual Examples

**Example 1: scalar × array**
```
(3, 4) × ()
After Rule 1:  (3, 4) × (1, 1)
After Rule 3:  (3, 4) × (3, 4) → result: (3, 4)

[[1,2,3,4],      ×  5  =  [[5,10,15,20],
 [5,6,7,8],               [25,30,35,40],
 [9,10,11,12]]             [45,50,55,60]]
```

**Example 2: (3,1) + (1,4)**
```
Shape A: (3, 1)     Shape B: (1, 4)
Rule 2:  compatible (one is 1 in each dim)
Rule 3:  stretch A to (3,4), stretch B to (3,4)
Result:  (3, 4)

[[1],       [[10, 20, 30, 40]]     [[11,21,31,41],
 [2],    +                    =     [12,22,32,42],
 [3]]                               [13,23,33,43]]
```

**Example 3: INCOMPATIBLE shapes**
```
Shape A: (3, 4)     Shape B: (3,)
After Rule 1:  (3, 4)  vs  (1, 3)
Dim 1: 4 vs 3 — neither is 1 → ERROR!

Fix: reshape B to (3, 1) or (1, 4)
```

```python
import numpy as np

# Example 1: Add bias to matrix (common in neural networks)
X = np.array([[1, 2, 3],
              [4, 5, 6],
              [7, 8, 9]])      # shape (3, 3)
bias = np.array([10, 20, 30])  # shape (3,) → broadcasts as (1,3) → (3,3)
result = X + bias
print(result)
# [[11, 22, 33],
#  [14, 25, 36],
#  [17, 28, 39]]

# Example 2: Normalize each row (subtract row mean)
# X shape: (100, 5) — 100 samples, 5 features
X = np.random.randn(100, 5)
row_means = X.mean(axis=1)       # shape (100,)
row_means = row_means[:, np.newaxis]  # shape (100, 1) — must reshape!
X_centered = X - row_means       # (100, 5) - (100, 1) → (100, 5)

# Example 3: Pairwise distances (classic ML operation)
# Compute distance from each of 5 query points to each of 3 database points
queries = np.array([[1.0, 2.0],
                    [3.0, 4.0],
                    [5.0, 6.0],
                    [7.0, 8.0],
                    [9.0, 10.0]])   # (5, 2)

db = np.array([[0.0, 0.0],
               [1.0, 1.0],
               [8.0, 8.0]])         # (3, 2)

# queries[:, np.newaxis] → (5, 1, 2)
# db[np.newaxis, :]       → (1, 3, 2)
# subtract broadcasts    → (5, 3, 2)
diff = queries[:, np.newaxis] - db[np.newaxis, :]
distances = np.sqrt((diff ** 2).sum(axis=2))  # (5, 3)
print(distances.shape)   # (5, 3) — distance from each query to each db point

# Example 4: Common gotcha — wrong shape
X = np.random.randn(100, 5)
col_means = X.mean(axis=0)    # shape (5,) — this is fine
X_norm = X - col_means        # (100,5) - (5,) works! NumPy pads left: (1,5) → (100,5)
```

---

## 8. Vectorized Operations vs Loops

The golden rule: **never write a Python for-loop over array elements when you can vectorize**.

```python
import numpy as np
import time

n = 1_000_000
a = np.random.randn(n)
b = np.random.randn(n)

# Method 1: Pure Python loop
def dot_loop(a, b):
    result = 0.0
    for i in range(len(a)):
        result += a[i] * b[i]
    return result

# Method 2: NumPy vectorized
def dot_numpy(a, b):
    return np.dot(a, b)

# Timing
start = time.time()
r1 = dot_loop(a, b)
loop_time = time.time() - start

start = time.time()
r2 = dot_numpy(a, b)
numpy_time = time.time() - start

print(f"Loop:  {loop_time:.3f}s  result: {r1:.4f}")
print(f"NumPy: {numpy_time:.4f}s  result: {r2:.4f}")
print(f"Speedup: {loop_time / numpy_time:.0f}x")
# Typical: Loop: 0.891s, NumPy: 0.002s, Speedup: ~400x

# The vectorization mindset: think in arrays, not elements
# BAD:
def sigmoid_loop(x):
    result = np.zeros_like(x)
    for i in range(len(x)):
        result[i] = 1.0 / (1.0 + np.exp(-x[i]))
    return result

# GOOD:
def sigmoid_vectorized(x):
    return 1.0 / (1.0 + np.exp(-x))

# These produce identical results, but vectorized is ~100x faster
z = np.linspace(-5, 5, 1000)
y = sigmoid_vectorized(z)   # applies to ALL elements at once
```

---

## 9. Mathematical Operations

```python
import numpy as np

x = np.array([1.0, 4.0, 9.0, 16.0, 25.0])

# Element-wise arithmetic (standard operators)
print(x + 2)       # [3., 6., 11., 18., 27.]
print(x * 3)       # [3., 12., 27., 48., 75.]
print(x ** 0.5)    # [1., 2., 3., 4., 5.]   — square root
print(x / 2)       # [0.5, 2., 4.5, 8., 12.5]

# Mathematical functions
print(np.sqrt(x))       # [1., 2., 3., 4., 5.]
print(np.exp(x))        # e^x for each element
print(np.log(x))        # natural log (ln)
print(np.log2(x))       # log base 2
print(np.log10(x))      # log base 10
print(np.log1p(x))      # log(1+x) — numerically stable for small x
print(np.abs([-1, 2, -3]))  # [1, 2, 3]

# Trigonometric
angles = np.linspace(0, 2*np.pi, 5)
print(np.sin(angles))
print(np.cos(angles))

# Rounding
y = np.array([1.4, 1.5, 1.6, 2.5, 3.7])
print(np.floor(y))     # [1., 1., 1., 2., 3.]  — round down
print(np.ceil(y))      # [2., 2., 2., 3., 4.]  — round up
print(np.round(y))     # [1., 2., 2., 2., 4.]  — round to nearest
print(np.clip(y, 1.5, 3.0))  # [1.5, 1.5, 1.6, 2.5, 3.0] — clamp values

# Numerically stable operations
print(np.expm1(0.001))   # e^0.001 - 1, more accurate than exp(x)-1 near 0
print(np.log1p(0.001))   # log(1.001), more accurate than log(1+x) near 0

# Common ML activations implemented with NumPy
def relu(x):
    return np.maximum(0, x)     # max of 0 and x, element-wise

def sigmoid(x):
    return 1.0 / (1.0 + np.exp(-x))

def softmax(x):
    """Numerically stable softmax."""
    e = np.exp(x - np.max(x))  # subtract max for stability
    return e / e.sum()

z = np.array([2.0, 1.0, 0.1])
print(softmax(z))    # [0.659, 0.242, 0.099] — sums to 1.0
```

---

## 10. Aggregations with Axis

The `axis` parameter tells NumPy *which dimension to collapse*. This is used constantly when computing statistics over batches or features.

```
Understanding axis:

Array shape: (3, 4)   — 3 rows, 4 columns

axis=None:  collapse ALL dimensions → scalar
axis=0:     collapse along rows → result shape (4,)   [per-column stats]
axis=1:     collapse along cols → result shape (3,)   [per-row stats]

Visualization:
          col0  col1  col2  col3
row0  →  [ 1,    2,    3,    4 ]   ← axis=1 collapses → scalar per row
row1  →  [ 5,    6,    7,    8 ]
row2  →  [ 9,   10,   11,   12]
          ↓     ↓     ↓     ↓
         axis=0 collapses → scalar per column
```

```python
import numpy as np

arr = np.array([[1, 2, 3, 4],
                [5, 6, 7, 8],
                [9,10,11,12]], dtype=float)

# Sum
print(np.sum(arr))          # 78.0   — sum of all elements
print(np.sum(arr, axis=0))  # [15. 18. 21. 24.]  — sum of each column
print(np.sum(arr, axis=1))  # [10. 26. 42.]      — sum of each row

# Mean
print(np.mean(arr, axis=0)) # [5., 6., 7., 8.]  — mean of each column

# Standard deviation
print(np.std(arr, axis=0))  # std of each column

# Min / Max
print(np.min(arr, axis=1))   # [1., 5., 9.]  — min of each row
print(np.max(arr, axis=0))   # [9.,10.,11.,12.]  — max of each column

# Argmin / Argmax — INDEX of min/max (not the value!)
print(np.argmax(arr, axis=1))   # [3, 3, 3]  — max is always in last column
print(np.argmin(arr, axis=0))   # [0, 0, 0, 0]  — min is always in first row

# keepdims=True — keep dimension for broadcasting
col_means = np.mean(arr, axis=0, keepdims=True)   # shape (1, 4) not (4,)
row_means  = np.mean(arr, axis=1, keepdims=True)  # shape (3, 1) not (3,)
centered   = arr - row_means   # (3,4) - (3,1) broadcasts correctly

# ML application: batch normalization (simplified)
def batch_normalize(X, eps=1e-8):
    """Normalize each feature to have mean=0, std=1 across the batch."""
    mean = X.mean(axis=0, keepdims=True)  # (1, n_features)
    std  = X.std(axis=0, keepdims=True)   # (1, n_features)
    return (X - mean) / (std + eps)

X = np.random.randn(32, 10) * 5 + 3    # batch of 32, 10 features
X_norm = batch_normalize(X)
print(f"Before: mean={X.mean(axis=0)[:3]}, std={X.std(axis=0)[:3]}")
print(f"After:  mean={X_norm.mean(axis=0)[:3].round(10)}, "
      f"std={X_norm.std(axis=0)[:3].round(4)}")
```

---

## 11. Linear Algebra

Linear algebra is the backbone of ML. Every neural network forward pass is matrix multiplication. Understanding NumPy's linear algebra operations is non-negotiable.

```python
import numpy as np

# ── Dot product (1D) ────────────────────────────────────────────────────
a = np.array([1, 2, 3])
b = np.array([4, 5, 6])
print(np.dot(a, b))     # 1*4 + 2*5 + 3*6 = 32

# ── Matrix multiplication ────────────────────────────────────────────────
A = np.array([[1, 2],
              [3, 4]])   # (2,2)
B = np.array([[5, 6],
              [7, 8]])   # (2,2)

# Three equivalent ways:
C1 = np.dot(A, B)        # classic
C2 = np.matmul(A, B)    # explicit matrix multiply
C3 = A @ B              # @ operator (Python 3.5+, preferred in modern code)
print(C3)
# [[19, 22],
#  [43, 50]]

# Shape rules: (m,n) @ (n,p) → (m,p)
# The inner dimensions MUST match!
X  = np.random.randn(100, 20)   # 100 samples, 20 features
W  = np.random.randn(20, 10)    # weight matrix
b  = np.random.randn(10)        # bias
Z  = X @ W + b                  # (100,20)@(20,10) + (10,) → (100,10)
print(Z.shape)   # (100, 10)

# ── Norms ────────────────────────────────────────────────────────────────
v = np.array([3.0, 4.0])

l1_norm = np.linalg.norm(v, ord=1)   # |3| + |4| = 7   (L1 / Manhattan)
l2_norm = np.linalg.norm(v, ord=2)   # sqrt(9+16) = 5  (L2 / Euclidean, default)
print(f"L1: {l1_norm}, L2: {l2_norm}")

# Matrix Frobenius norm (sum of squared elements)
W = np.random.randn(3, 3)
frob = np.linalg.norm(W)  # default for matrices is Frobenius
# This is used in L2 regularization: loss += lambda * ||W||_F^2

# ── Matrix inverse ───────────────────────────────────────────────────────
A = np.array([[2., 1.],
              [5., 3.]])
A_inv = np.linalg.inv(A)
print(A @ A_inv)   # should be identity (up to floating point error)
print(np.allclose(A @ A_inv, np.eye(2)))  # True

# ── Determinant ──────────────────────────────────────────────────────────
print(np.linalg.det(A))   # 2*3 - 1*5 = 1.0

# ── Solving linear systems: Ax = b ──────────────────────────────────────
A = np.array([[2., 1.],
              [1., 3.]])
b = np.array([5., 10.])
x = np.linalg.solve(A, b)   # more numerically stable than A_inv @ b
print(x)                     # solution vector

# ── Eigenvalues and Eigenvectors ─────────────────────────────────────────
A = np.array([[4., 2.],
              [1., 3.]])
eigenvalues, eigenvectors = np.linalg.eig(A)
print("Eigenvalues:", eigenvalues)    # [5., 2.]
print("Eigenvectors (columns):")
print(eigenvectors)

# Verify: A @ v = λ @ v (for each eigenvector)
for i in range(len(eigenvalues)):
    v = eigenvectors[:, i]
    lam = eigenvalues[i]
    print(np.allclose(A @ v, lam * v))  # True

# ── SVD (Singular Value Decomposition) ──────────────────────────────────
# For any matrix A (not just square), A = U @ Σ @ V^T
A = np.array([[1., 2., 3.],
              [4., 5., 6.]])  # (2, 3)

U, s, Vt = np.linalg.svd(A, full_matrices=False)
print(f"U shape: {U.shape}")   # (2, 2)  — left singular vectors
print(f"s shape: {s.shape}")   # (2,)    — singular values (always positive!)
print(f"Vt shape: {Vt.shape}") # (2, 3)  — right singular vectors (transposed)

# Reconstruct A from its SVD
S = np.diag(s)
A_reconstructed = U @ S @ Vt
print(np.allclose(A, A_reconstructed))   # True

# ── Rank and Condition Number ────────────────────────────────────────────
A = np.array([[1., 2.], [2., 4.]])   # rank-deficient (rows are multiples)
print(np.linalg.matrix_rank(A))      # 1 (not 2!)
print(np.linalg.cond(A))             # very large → ill-conditioned
```

---

## 12. Reshaping and Transposing

Reshaping is how you change the *view* of data without copying it. In ML, you'll reshape constantly — flattening images for dense layers, adding batch dimensions, etc.

```python
import numpy as np

arr = np.arange(24)     # [0, 1, 2, ..., 23]  shape: (24,)

# ── reshape ──────────────────────────────────────────────────────────────
# Total elements must be the same!
A = arr.reshape(4, 6)        # (24,) → (4, 6)
B = arr.reshape(2, 3, 4)     # (24,) → (2, 3, 4)
C = arr.reshape(2, -1)       # -1 means "infer this dim" → (2, 12)
D = arr.reshape(-1, 6)       # → (4, 6)

# reshape usually returns a VIEW (no copy), but not always!
A[0, 0] = 999
print(arr[0])    # 999 — it was a view!

# ── flatten vs ravel ──────────────────────────────────────────────────────
mat = np.array([[1, 2, 3], [4, 5, 6]])
print(mat.flatten())    # [1 2 3 4 5 6]  — always returns a COPY
print(mat.ravel())      # [1 2 3 4 5 6]  — returns VIEW if possible

# ── transpose ─────────────────────────────────────────────────────────────
#
# Before:  (3, 4)         After .T:  (4, 3)
# ┌───────────────┐       ┌─────────────┐
# │ 1  2  3  4   │       │ 1  5  9    │
# │ 5  6  7  8   │  .T → │ 2  6  10   │
# │ 9  10 11 12  │       │ 3  7  11   │
# └───────────────┘       │ 4  8  12   │
#                         └─────────────┘

mat = np.arange(12).reshape(3, 4)
print(mat.T.shape)       # (4, 3)
print(mat.T)

# 3D transpose — specify axis order
images = np.random.rand(32, 3, 224, 224)   # (batch, channels, H, W) — PyTorch format
# Convert to TensorFlow format (batch, H, W, channels)
images_tf = np.transpose(images, (0, 2, 3, 1))
print(images_tf.shape)  # (32, 224, 224, 3)

# ── np.newaxis (adds a dimension of size 1) ───────────────────────────────
#
# This is essential for broadcasting!
v = np.array([1, 2, 3])       # shape (3,)
col = v[:, np.newaxis]         # shape (3, 1) — column vector
row = v[np.newaxis, :]         # shape (1, 3) — row vector

print(col.shape)  # (3, 1)
print(row.shape)  # (1, 3)

# Visualization:
# v         = [1  2  3]          shape (3,)
# col       = [[1],              shape (3, 1)
#              [2],
#              [3]]
# row       = [[1  2  3]]        shape (1, 3)

# Common use: add batch dimension to single sample
single_image = np.random.rand(28, 28)      # (28, 28)
batch_of_one = single_image[np.newaxis]    # (1, 28, 28) — "batch" of 1 image
# or equivalently:
batch_of_one = np.expand_dims(single_image, axis=0)
```

---

## 13. Stacking and Splitting

```python
import numpy as np

a = np.array([1, 2, 3])
b = np.array([4, 5, 6])

# ── np.concatenate ────────────────────────────────────────────────────────
# Joins arrays along EXISTING axis — shapes must match on all other axes
print(np.concatenate([a, b]))           # [1 2 3 4 5 6]  — 1D concat

A = np.array([[1, 2], [3, 4]])   # (2,2)
B = np.array([[5, 6]])           # (1,2)
print(np.concatenate([A, B], axis=0))   # (3,2) — stack vertically
# [[1,2],[3,4],[5,6]]

# ── np.vstack / np.hstack ─────────────────────────────────────────────────
# Convenience wrappers around concatenate
print(np.vstack([A, B]))         # same as concatenate axis=0 — adds ROWS
                                 # (2,2) + (1,2) → (3,2)
C = np.array([[7], [8]])         # (2,1)
print(np.hstack([A, C]))         # (2,2) + (2,1) → (2,3) — adds COLUMNS
# [[1,2,7],[3,4,8]]

# ── np.stack ─────────────────────────────────────────────────────────────
# Creates a NEW axis — all arrays must have identical shape
x = np.array([1, 2, 3])    # (3,)
y = np.array([4, 5, 6])    # (3,)
print(np.stack([x, y], axis=0))    # (2, 3) — new axis at position 0
# [[1,2,3],
#  [4,5,6]]
print(np.stack([x, y], axis=1))    # (3, 2) — new axis at position 1
# [[1,4],[2,5],[3,6]]

# ML use: building a batch from individual samples
sample1 = np.random.rand(28, 28)
sample2 = np.random.rand(28, 28)
sample3 = np.random.rand(28, 28)
batch = np.stack([sample1, sample2, sample3])  # (3, 28, 28)

# ── Splitting ─────────────────────────────────────────────────────────────
data = np.arange(24).reshape(6, 4)
parts = np.split(data, 3)           # split into 3 equal parts along axis 0
print(len(parts), parts[0].shape)  # 3 parts, each (2, 4)

train, val, test = np.split(data, [4, 5])   # split at indices 4 and 5
print(train.shape, val.shape, test.shape)   # (4,4), (1,4), (1,4)
```

---

## 14. Useful Utility Functions

```python
import numpy as np

# ── np.where ──────────────────────────────────────────────────────────────
# np.where(condition, value_if_true, value_if_false)
x = np.array([-2, -1, 0, 1, 2])
relu = np.where(x > 0, x, 0)          # ReLU!
print(relu)  # [0, 0, 0, 1, 2]

# More complex: replace outliers
data = np.array([1, 2, 100, 3, -50, 4])
cleaned = np.where(np.abs(data) < 10, data, np.median(data))
print(cleaned)  # [1, 2, 3, 3, 3, 4]  — outliers replaced with median

# ── np.argsort ────────────────────────────────────────────────────────────
scores = np.array([0.1, 0.9, 0.3, 0.7, 0.5])
sorted_idx = np.argsort(scores)            # indices that would sort the array
print(sorted_idx)    # [0, 2, 4, 3, 1] — ascending
print(scores[sorted_idx])  # [0.1, 0.3, 0.5, 0.7, 0.9]

# Top-k predictions (e.g., top-5 classification)
top5_idx = np.argsort(scores)[-5:][::-1]   # top 5, descending
print(top5_idx)    # [1, 3, 4, 2, 0]

# ── np.argmax / np.argmin ─────────────────────────────────────────────────
probs = np.array([[0.1, 0.7, 0.2],     # sample 1: class 1 is most likely
                  [0.8, 0.1, 0.1],     # sample 2: class 0 is most likely
                  [0.2, 0.3, 0.5]])    # sample 3: class 2 is most likely

predictions = np.argmax(probs, axis=1)
print(predictions)   # [1, 0, 2] — predicted class for each sample

# ── np.unique ─────────────────────────────────────────────────────────────
labels = np.array([2, 1, 0, 2, 1, 0, 0, 1, 2])
unique = np.unique(labels)
print(unique)   # [0, 1, 2]

# With return_counts
unique, counts = np.unique(labels, return_counts=True)
for cls, cnt in zip(unique, counts):
    print(f"Class {cls}: {cnt} samples ({cnt/len(labels):.1%})")

# ── np.bincount ──────────────────────────────────────────────────────────
# Faster than unique for counting integers 0 to max_value
counts = np.bincount(labels)
print(counts)   # [3, 3, 3] — count of 0s, 1s, 2s

# ── np.cumsum / np.cumprod ────────────────────────────────────────────────
x = np.array([1, 2, 3, 4, 5])
print(np.cumsum(x))    # [1, 3, 6, 10, 15]  — running sum
print(np.cumprod(x))   # [1, 2, 6, 24, 120] — running product

# ── np.clip ───────────────────────────────────────────────────────────────
predictions = np.array([-0.1, 0.3, 0.7, 1.2])
clipped = np.clip(predictions, 0, 1)   # ensure probabilities stay in [0,1]
print(clipped)   # [0.0, 0.3, 0.7, 1.0]
```

---

## 15. Random Module and Reproducibility

Reproducibility is critical in ML: you need to be able to re-run experiments and get the same results. NumPy's random module must be seeded properly.

```python
import numpy as np

# ── Old API (still works but deprecated style) ────────────────────────────
np.random.seed(42)           # global seed — affects all subsequent calls
x = np.random.randn(5)
y = np.random.randint(0, 10, 5)

# Problem: different libraries share the global state
# Order of function calls determines results → fragile

# ── Modern API (NumPy 1.17+): use Generator objects ──────────────────────
rng = np.random.default_rng(seed=42)   # self-contained generator

# This generator is independent of the global state
x = rng.standard_normal(5)             # N(0,1)
y = rng.integers(0, 10, 5)             # integers in [0, 10)
z = rng.random(5)                      # uniform [0, 1)
w = rng.choice(100, size=10, replace=False)  # random sample without replacement

# Reproducibility check
rng1 = np.random.default_rng(42)
rng2 = np.random.default_rng(42)
a = rng1.random(5)
b = rng2.random(5)
print(np.allclose(a, b))    # True — same seed, same results

# ── Distribution reference ────────────────────────────────────────────────
rng = np.random.default_rng(42)

# Uniform [low, high)
u = rng.uniform(low=0.0, high=1.0, size=(3, 3))

# Normal with mean and std
n = rng.normal(loc=0.0, scale=1.0, size=1000)

# Binomial (n trials, probability p)
b = rng.binomial(n=10, p=0.5, size=100)

# Exponential (for inter-arrival times)
e = rng.exponential(scale=1.0, size=100)

# Shuffle (in-place)
arr = np.arange(10)
rng.shuffle(arr)
print(arr)  # shuffled

# Permutation (returns new shuffled array)
perm = rng.permutation(10)   # shuffled array of 0..9

# ML use: train/val/test split
n_samples = 1000
indices = rng.permutation(n_samples)
train_idx = indices[:700]    # 70%
val_idx   = indices[700:850] # 15%
test_idx  = indices[850:]    # 15%

# Weight initialization (Kaiming/He initialization for ReLU layers)
fan_in = 512
W = rng.normal(0, np.sqrt(2.0 / fan_in), size=(fan_in, 256))
```

---

## 16. Memory Layout

Understanding memory layout becomes important when performance is critical and when interfacing with libraries like PyTorch.

```python
import numpy as np

# C-order (row-major): default in NumPy — rows stored contiguously
# Fortran-order (col-major): columns stored contiguously

A = np.array([[1, 2, 3],
              [4, 5, 6]])

C = np.ascontiguousarray(A)      # ensure C-order
F = np.asfortranarray(A)         # convert to F-order

print(C.flags['C_CONTIGUOUS'])   # True
print(C.flags['F_CONTIGUOUS'])   # False
print(F.flags['C_CONTIGUOUS'])   # False
print(F.flags['F_CONTIGUOUS'])   # True

# Transpose is a view with swapped strides — might not be contiguous!
A_T = A.T
print(A_T.flags['C_CONTIGUOUS'])  # False — .T swaps strides, doesn't copy

# If you need a contiguous copy after transpose:
A_T_c = np.ascontiguousarray(A.T)   # now it's C-contiguous

# Memory layout affects performance in tight loops
# Rule: access elements in memory order for best cache performance
large = np.random.rand(1000, 1000)

# Summing rows vs columns
import time
start = time.time()
for _ in range(100):
    s = large.sum(axis=1)   # sum along columns (row-major: rows are fast)
row_time = time.time() - start

start = time.time()
for _ in range(100):
    s = large.sum(axis=0)   # sum along rows (slightly less cache-friendly)
col_time = time.time() - start
print(f"Row sum: {row_time:.3f}s, Col sum: {col_time:.3f}s")
```

---

## 17. Saving and Loading

```python
import numpy as np

X = np.random.randn(1000, 50)
y = np.random.randint(0, 10, 1000)

# ── .npy — single array (binary, fast) ───────────────────────────────────
np.save("features.npy", X)          # saves to features.npy
X_loaded = np.load("features.npy")  # loads it back
print(np.allclose(X, X_loaded))     # True

# ── .npz — multiple arrays (compressed) ──────────────────────────────────
np.savez("dataset.npz", X=X, y=y)                    # uncompressed
np.savez_compressed("dataset_c.npz", X=X, y=y)       # compressed (smaller file)

data = np.load("dataset.npz")
X_back = data["X"]     # access by name
y_back = data["y"]

# ── Text formats ──────────────────────────────────────────────────────────
np.savetxt("data.csv", X[:10], delimiter=",", fmt="%.6f")
X_text = np.loadtxt("data.csv", delimiter=",")

# ── Best practices ────────────────────────────────────────────────────────
# For small datasets:   .npy or .npz
# For large datasets:   use HDF5 (h5py library) or memory-mapped arrays
# For sharing data:     CSV (human-readable) or Parquet
# For deep learning:    frameworks have their own serialization (torch.save)

# Memory-mapped arrays — access huge arrays without loading into RAM
fp = np.memmap("large_array.dat", dtype='float32', mode='w+', shape=(1000000, 128))
fp[:100] = np.random.randn(100, 128)   # write 100 rows
del fp  # flush to disk
fp_read = np.memmap("large_array.dat", dtype='float32', mode='r', shape=(1000000, 128))
print(fp_read[0, :5])   # reads only this piece from disk
```

---

## 18. Mini Project: Linear Regression from Scratch

Now let's put everything together and implement linear regression using only NumPy. This is the foundation of all supervised learning — understanding it deeply will pay dividends later.

```
The goal: find weights W and bias b such that:
    ŷ = X @ W + b  ≈  y

We minimize Mean Squared Error:
    MSE = (1/n) * Σ (y_i - ŷ_i)²

Using gradient descent:
    W := W - α * (∂MSE/∂W)
    b := b - α * (∂MSE/∂b)

Gradients:
    ∂MSE/∂W = (2/n) * X^T @ (ŷ - y)
    ∂MSE/∂b = (2/n) * sum(ŷ - y)
```

```python
import numpy as np
import matplotlib.pyplot as plt


class LinearRegressionNumPy:
    """
    Linear Regression implemented from scratch using only NumPy.
    
    Model:  ŷ = X @ W + b
    Loss:   MSE = mean((y - ŷ)²)
    Update: gradient descent
    """
    
    def __init__(self, learning_rate: float = 0.01, n_iterations: int = 1000):
        self.lr = learning_rate
        self.n_iter = n_iterations
        self.W = None
        self.b = None
        self.loss_history = []
    
    def _initialize_weights(self, n_features: int):
        """Small random initialization."""
        rng = np.random.default_rng(42)
        self.W = rng.normal(0, 0.01, size=(n_features, 1))
        self.b = np.zeros(1)
    
    def _forward(self, X: np.ndarray) -> np.ndarray:
        """Compute predictions: ŷ = X @ W + b"""
        return X @ self.W + self.b     # (n, p) @ (p, 1) + (1,) → (n, 1)
    
    def _compute_loss(self, y_true: np.ndarray, y_pred: np.ndarray) -> float:
        """MSE loss."""
        residuals = y_true - y_pred    # (n, 1)
        return float(np.mean(residuals ** 2))
    
    def _compute_gradients(self, X, y_true, y_pred):
        """
        Gradients of MSE with respect to W and b.
        
        MSE = (1/n) * Σ (y - ŷ)²
        ∂MSE/∂W = -(2/n) * X^T @ (y - ŷ)  [shape: (p, 1)]
        ∂MSE/∂b = -(2/n) * sum(y - ŷ)     [scalar]
        """
        n = len(y_true)
        error = y_pred - y_true          # (n, 1) — note: reversed for sign
        
        dW = (2 / n) * (X.T @ error)    # (p, n) @ (n, 1) → (p, 1)
        db = (2 / n) * np.sum(error)    # scalar
        return dW, db
    
    def fit(self, X: np.ndarray, y: np.ndarray) -> "LinearRegressionNumPy":
        """Fit using gradient descent."""
        X = np.asarray(X, dtype=np.float64)
        y = np.asarray(y, dtype=np.float64).reshape(-1, 1)  # ensure column vector
        
        n_samples, n_features = X.shape
        self._initialize_weights(n_features)
        self.loss_history = []
        
        for iteration in range(self.n_iter):
            # Forward pass
            y_pred = self._forward(X)
            
            # Compute loss
            loss = self._compute_loss(y, y_pred)
            self.loss_history.append(loss)
            
            # Compute gradients
            dW, db = self._compute_gradients(X, y, y_pred)
            
            # Update weights (gradient descent step)
            self.W -= self.lr * dW
            self.b -= self.lr * db
            
            # Log progress
            if iteration % 100 == 0:
                print(f"Iteration {iteration:4d}: MSE = {loss:.6f}")
        
        return self
    
    def predict(self, X: np.ndarray) -> np.ndarray:
        """Generate predictions."""
        X = np.asarray(X, dtype=np.float64)
        return self._forward(X).flatten()
    
    def score(self, X: np.ndarray, y: np.ndarray) -> float:
        """R² score: coefficient of determination."""
        y_pred = self.predict(X)
        y = np.asarray(y)
        ss_res = np.sum((y - y_pred) ** 2)
        ss_tot = np.sum((y - y.mean()) ** 2)
        return 1 - (ss_res / ss_tot)


# ─── Generate synthetic data ────────────────────────────────────────────────
np.random.seed(42)
n = 200
X = np.random.randn(n, 3)      # 3 features
true_W = np.array([2.5, -1.3, 0.8])
true_b = 4.0
noise  = np.random.randn(n) * 0.5
y = X @ true_W + true_b + noise

# Split
X_train, X_test = X[:160], X[160:]
y_train, y_test = y[:160], y[160:]

# Feature normalization (critical for gradient descent!)
mean = X_train.mean(axis=0)
std  = X_train.std(axis=0)
X_train_n = (X_train - mean) / std
X_test_n  = (X_test  - mean) / std

# ─── Train ──────────────────────────────────────────────────────────────────
model = LinearRegressionNumPy(learning_rate=0.1, n_iterations=500)
model.fit(X_train_n, y_train)

# ─── Evaluate ────────────────────────────────────────────────────────────────
train_r2 = model.score(X_train_n, y_train)
test_r2  = model.score(X_test_n,  y_test)
print(f"\nTrain R²: {train_r2:.4f}")
print(f"Test  R²: {test_r2:.4f}")
print(f"\nLearned W: {model.W.flatten()}")
print(f"Learned b: {model.b[0]:.4f}")

# ─── Closed-form solution (for comparison) ──────────────────────────────────
# For linear regression: W* = (X^T X)^(-1) X^T y
X_aug = np.column_stack([X_train_n, np.ones(len(X_train_n))])  # add bias column
W_exact = np.linalg.lstsq(X_aug, y_train, rcond=None)[0]
print(f"\nExact solution (lstsq): {W_exact}")
```

---

## 19. Summary

```
NumPy Fundamentals — Core Concepts
───────────────────────────────────────────────────────────────────────────
ndarray    Homogeneous typed, contiguous memory, shape/stride/dtype
dtypes     float32 for GPU, float64 for CPU, uint8 for images
indexing   [i, j], [start:stop:step], [mask], [indices_array]
broadcast  Shapes compatible if sizes equal or one is 1
vectorize  NEVER loop over array elements; use ufuncs instead
axis       axis=0 collapses rows (per-column), axis=1 collapses columns
linalg     @, np.linalg.inv/det/eig/svd/solve — foundation of all ML
reshape    .reshape(), .T, np.newaxis, .flatten() — view, rarely copy
random     np.random.default_rng(seed) — always seed for reproducibility
───────────────────────────────────────────────────────────────────────────
```

### Key Operations at a Glance

| Operation | Code | Notes |
|-----------|------|-------|
| Element-wise ops | `a + b`, `a * b`, `np.exp(a)` | All broadcasting-aware |
| Matrix multiply | `A @ B` | Shape: (m,n)@(n,p)→(m,p) |
| Dot product | `np.dot(a, b)` | 1D: scalar; 2D: matmul |
| Transpose | `A.T` | Returns view, not copy |
| Reshape | `A.reshape(m, n)` | -1 to infer dim |
| Select with mask | `arr[arr > 0]` | Returns 1D array of matches |
| Add dimension | `arr[:, np.newaxis]` | For broadcasting |
| Aggregate | `arr.mean(axis=0)` | axis=0→per-col, axis=1→per-row |

---

## 20. Exercises

**Exercise 1: Broadcasting Challenge**
Without running the code, predict the output shape for each of these operations:
- `np.zeros((3,)) + np.zeros((4, 3))`
- `np.zeros((5, 1)) + np.zeros((1, 7))`
- `np.zeros((3, 4, 5)) + np.zeros((4, 1))`
- `np.zeros((2, 3)) + np.zeros((2,))`
Then run them to verify.

*Hint: Align shapes from the right, pad with 1s on the left where needed.*

**Exercise 2: Vectorized Operations**
Implement the following operations using only NumPy (no loops):
1. Given a matrix `X` of shape (n, m), compute pairwise squared Euclidean distances between all rows. Result shape should be (n, n).
2. Compute the softmax of each row of a matrix X (shape n×k): each row should sum to 1.
3. Given binary predictions `y_pred` and true labels `y_true` (both 1D integer arrays), compute accuracy, precision, recall, and F1 score.

*Hint for (1): ||a-b||² = ||a||² + ||b||² - 2a·b. Use broadcasting.*

**Exercise 3: Numerical Stability**
Write a function `log_sum_exp(x)` that computes `log(sum(exp(x_i)))` in a numerically stable way (without overflow). Test it with `x = [1000, 1001, 1002]` — the naive version will overflow, yours should not.

*Hint: Factor out `max(x)`: log(Σexp(x_i)) = max(x) + log(Σexp(x_i - max(x)))*

**Exercise 4: Data Preprocessing Pipeline**
Write a function `preprocess(X_train, X_test)` that:
1. Removes constant columns (std = 0) — they carry no information
2. Standardizes remaining features using train statistics
3. Returns `(X_train_clean, X_test_clean, kept_indices)` where `kept_indices` are the original column indices that were kept

*Hint: Use `np.std(X_train, axis=0)` and boolean indexing on columns.*

**Exercise 5: Gradient Descent Visualization**
Extend the `LinearRegressionNumPy` class to support mini-batch gradient descent: instead of using all samples per iteration, randomly sample `batch_size` rows. Compare convergence curves (loss vs iterations) between batch sizes: 1 (SGD), 32 (mini-batch), and full dataset (batch GD). Plot all three curves.

*Hint: Use `np.random.choice(n_samples, batch_size, replace=False)` to get batch indices.*

---

**What's Next →** [Chapter 03: Pandas — Data Wrangling and Analysis](./03-pandas-data-wrangling.md)

*NumPy is the engine, but data in the real world comes in tables — CSVs, databases, Excel files. Pandas is the library that handles tabular data, and you'll use it every single day for loading, cleaning, and exploring data before any ML training begins.*
